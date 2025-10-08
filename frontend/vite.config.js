import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'https://super-duper-space-doodle-wrjw7w7q9q5gc5xxw-3000.app.github.dev',
        changeOrigin: true,
        secure: false
      }
    }
  }
});