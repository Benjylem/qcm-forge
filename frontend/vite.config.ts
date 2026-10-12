import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    // Évite le CORS en dev : le navigateur appelle le frontend sur
    // 5173, et /health, /api/... partent transparemment vers l'API Go.
    proxy: {
      '/health': 'http://localhost:8080',
      '/signup': 'http://localhost:8080',
      '/login': 'http://localhost:8080',
      '/logout': 'http://localhost:8080',
      '/me': 'http://localhost:8080',
      '/api': 'http://localhost:8080',
    },
  },
})
