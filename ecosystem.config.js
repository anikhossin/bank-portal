// Start both Svelte apps: pm2 start ecosystem.config.js
// Stop: pm2 stop ecosystem.config.js
// Logs: pm2 logs
module.exports = {
  apps: [
    {
      name: 'portal',
      cwd: './portal',
      script: 'node_modules/vite/bin/vite.js',
      args: '--host --port 3000',
      interpreter: 'node',
      watch: false,
      autorestart: true,
    },
    {
      name: 'admin',
      cwd: './admin',
      script: 'node_modules/vite/bin/vite.js',
      args: '--host --port 3002',
      interpreter: 'node',
      watch: false,
      autorestart: true,
    },
  ],
};
