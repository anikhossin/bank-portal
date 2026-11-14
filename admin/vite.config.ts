import { svelte } from "@sveltejs/vite-plugin-svelte";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";

// https://vite.dev/config/
const apiProxy = {
  "/api": "http://127.0.0.1:8080",
};

export default defineConfig({
  plugins: [tailwindcss(), svelte()],
  server: {
    port: 3002,
    allowedHosts: ["admin02243.meine-deutsche-bank-treuhand.de"],
    proxy: apiProxy,
  },
  preview: {
    port: 3002,
    allowedHosts: ["admin02243.meine-deutsche-bank-treuhand.de"],
    proxy: apiProxy,
  },
});
