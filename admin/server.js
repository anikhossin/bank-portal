import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { start } from '../scripts/serve.mjs'

start({ root: resolve(dirname(fileURLToPath(import.meta.url))) })
