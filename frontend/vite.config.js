import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [tailwindcss(), svelte()],
  build: {
    outDir: '../web/static',
    emptyOutDir: true,
  },
  server: {
    proxy: {
      '/rooms': {
        target: 'http://localhost:8080',
        bypass(req) {
          // Let Vite serve its own index.html for page-load GETs (/rooms/{code})
          // Only proxy actual API calls (/rooms/{code}/action or SSE)
          const segments = req.url.split('/').filter(Boolean);
          if (req.method === 'GET' && segments.length <= 2) return req.url;
        },
      },
    },
  },
})
