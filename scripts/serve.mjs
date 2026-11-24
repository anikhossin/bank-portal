import { createReadStream, existsSync, statSync } from 'node:fs'
import { createServer, request as httpRequest } from 'node:http'
import { request as httpsRequest } from 'node:https'
import { extname, join, resolve, sep } from 'node:path'

const TYPES = {
  '.html': 'text/html; charset=utf-8',
  '.js': 'text/javascript; charset=utf-8',
  '.mjs': 'text/javascript; charset=utf-8',
  '.css': 'text/css; charset=utf-8',
  '.json': 'application/json; charset=utf-8',
  '.svg': 'image/svg+xml',
  '.png': 'image/png',
  '.jpg': 'image/jpeg',
  '.jpeg': 'image/jpeg',
  '.webp': 'image/webp',
  '.gif': 'image/gif',
  '.ico': 'image/x-icon',
  '.woff': 'font/woff',
  '.woff2': 'font/woff2',
  '.txt': 'text/plain; charset=utf-8',
  '.map': 'application/json; charset=utf-8',
  '.webmanifest': 'application/manifest+json',
}

const HOP_BY_HOP = new Set([
  'connection',
  'keep-alive',
  'proxy-authenticate',
  'proxy-authorization',
  'te',
  'trailers',
  'transfer-encoding',
  'upgrade',
  'host',
])

/**
 * Serves a Vite `dist` folder and forwards `/api` to the Go API.
 * `root` is the app directory that contains `dist`.
 */
export function start({ root }) {
  const dist = resolve(root, 'dist')
  const indexFile = join(dist, 'index.html')
  if (!existsSync(indexFile)) {
    console.error(`Missing ${indexFile}. Run npm run build in this app first.`)
    process.exit(1)
  }

  const port = Number(process.env.PORT || 3000)
  const host = process.env.HOST || '0.0.0.0'
  const upstream = (process.env.API_UPSTREAM || 'http://127.0.0.1:8080').replace(/\/$/, '')

  const server = createServer((req, res) => {
    const url = new URL(req.url || '/', 'http://localhost')
    if (url.pathname === '/api' || url.pathname.startsWith('/api/')) {
      proxy(req, res, upstream)
      return
    }
    serveStatic(req, res, dist, url.pathname)
  })

  server.on('error', (err) => {
    console.error(err.message)
    process.exit(1)
  })

  server.listen(port, host, () => {
    console.log(`listening on http://${host}:${port} (api ${upstream})`)
  })

  for (const signal of ['SIGTERM', 'SIGINT']) {
    process.on(signal, () => {
      server.close(() => process.exit(0))
    })
  }
}

function serveStatic(req, res, dist, pathname) {
  let rel
  try {
    rel = decodeURIComponent(pathname)
  } catch {
    send(res, 400, 'Bad request')
    return
  }

  const file = safeFile(dist, rel)
  if (file && isFile(file)) {
    sendFile(req, res, file, extname(file), rel.startsWith('/assets/'))
    return
  }

  const ext = extname(rel)
  if (req.method !== 'GET' && req.method !== 'HEAD') {
    send(res, 405, 'Method not allowed')
    return
  }
  if (ext && ext !== '.html') {
    send(res, 404, 'Not found')
    return
  }

  sendFile(req, res, join(dist, 'index.html'), '.html', false)
}

function safeFile(dist, pathname) {
  const rel = pathname.replace(/^\/+/, '')
  const file = resolve(dist, rel)
  const root = resolve(dist)
  if (file !== root && !file.startsWith(root + sep)) return null
  return file
}

function isFile(file) {
  try {
    return statSync(file).isFile()
  } catch {
    return false
  }
}

function sendFile(req, res, file, ext, immutable) {
  const headers = {
    'content-type': TYPES[ext] || 'application/octet-stream',
    'cache-control': immutable ? 'public, max-age=31536000, immutable' : 'no-cache',
    'x-content-type-options': 'nosniff',
  }
  res.writeHead(200, headers)
  if (req.method === 'HEAD') {
    res.end()
    return
  }
  createReadStream(file).pipe(res)
}

function send(res, status, body) {
  res.writeHead(status, { 'content-type': 'text/plain; charset=utf-8' })
  res.end(body)
}

function proxy(req, res, upstream) {
  let target
  try {
    target = new URL(req.url || '/', upstream.endsWith('/') ? upstream : `${upstream}/`)
  } catch {
    send(res, 500, 'Bad API upstream')
    return
  }

  const headers = { ...req.headers, host: target.host }
  for (const name of HOP_BY_HOP) delete headers[name]

  const request = target.protocol === 'https:' ? httpsRequest : httpRequest
  const upstreamReq = request(target, { method: req.method, headers }, (upstreamRes) => {
    const out = { ...upstreamRes.headers }
    for (const name of HOP_BY_HOP) delete out[name]
    res.writeHead(upstreamRes.statusCode || 502, out)
    upstreamRes.pipe(res)
  })

  upstreamReq.on('error', () => {
    if (!res.headersSent) send(res, 502, 'API unavailable')
  })
  req.pipe(upstreamReq)
}
