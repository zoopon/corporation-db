import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    host: '0.0.0.0',
    port: 3000,
    allowedHosts: ['frontend', 'localhost', '127.0.0.1'],
    watch: {
      usePolling: true,
    },
    hmr: {
      port: 3000,
    },
  },
  preview: {
    host: '0.0.0.0',
    port: 3000,
  },
})
