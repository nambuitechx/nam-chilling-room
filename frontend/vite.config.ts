import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,      // 👈 run on http://localhost:3000
    open: true       // optional: auto-open browser
  }
})
