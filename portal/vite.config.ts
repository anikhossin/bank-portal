import { svelte } from '@sveltejs/vite-plugin-svelte'
import tailwindcss from '@tailwindcss/vite'
import { defineConfig } from 'vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [tailwindcss(), svelte()],
  server: {
    port: 5173,
    allowedHosts: [
      'meine-deutsche-bank-treuhand.de',
      'www.meine-deutsche-bank-treuhand.de',
    ],
  },
  preview: {
    allowedHosts: [
      'meine-deutsche-bank-treuhand.de',
      'www.meine-deutsche-bank-treuhand.de',
    ],
  },
})
