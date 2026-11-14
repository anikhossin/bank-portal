import { svelte } from "@sveltejs/vite-plugin-svelte";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";

// https://vite.dev/config/
const apiProxy = {
  "/api": "http://127.0.0.1:8080",
};

const hosts = [
  "meine-Deutsche-bank-treuhand.de",
  "www.meine-Deutsche-bank-treuhand.de",
];

export default defineConfig({
  plugins: [tailwindcss(), svelte()],
  server: {
    port: 3000,
    allowedHosts: hosts,
    proxy: apiProxy,
  },
  preview: {
    port: 3000,
    allowedHosts: hosts,
    proxy: apiProxy,
  },
});
