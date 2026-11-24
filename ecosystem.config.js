// Build, then start both sites:
//   npm run build --prefix portal
//   npm run build --prefix admin
//   pm2 start ecosystem.config.js
// Stop: pm2 stop ecosystem.config.js
// Logs: pm2 logs
module.exports = {
  apps: [
    {
      name: 'portal',
      cwd: './portal',
      script: 'server.js',
      interpreter: 'node',
      env: {
        NODE_ENV: 'production',
        PORT: '3000',
        HOST: '0.0.0.0',
        API_UPSTREAM: 'http://127.0.0.1:8080',
      },
      watch: false,
      autorestart: true,
    },
    {
      name: 'admin',
      cwd: './admin',
      script: 'server.js',
      interpreter: 'node',
      env: {
        NODE_ENV: 'production',
        PORT: '3002',
        HOST: '0.0.0.0',
        API_UPSTREAM: 'http://127.0.0.1:8080',
      },
      watch: false,
      autorestart: true,
    },
  ],
};
