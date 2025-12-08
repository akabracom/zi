import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  base: '/zi/', // Для GitHub Pages: имя репозитория - zi
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      // Проксируем backend, чтобы избежать CORS-проблем
      '/api': {
        target: 'http://localhost:3001',
        changeOrigin: true,
        secure: false,
      },
      '/images': {
        target: 'http://localhost:3001',
        changeOrigin: true,
        secure: false,
        ws: false,
      },
    },
  },
})
