import { fileURLToPath, URL } from 'node:url';

import tailwindcss from '@tailwindcss/vite';
import vue from '@vitejs/plugin-vue';
import { defineConfig } from 'vite';
import vueDevTools from 'vite-plugin-vue-devtools';

export default defineConfig({
  plugins: [vue(), vueDevTools(), tailwindcss()],
  envPrefix: 'VITE_',
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    port: 3100,
    strictPort: true,
    proxy: {
      '/v1': 'http://localhost:8080',
      '/health': 'http://localhost:8080',
    },
  },
});
