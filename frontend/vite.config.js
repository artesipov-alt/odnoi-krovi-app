import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import path from 'path'; // Добавлено: для алиасов

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'), // Добавлено: поддержка алиаса @ для src
    },
  },
  server: {
    port: 5173,
    host: '0.0.0.0',
    proxy: {
      '/api': {
        target: 'https://super-duper-space-doodle-wrjw7w7q9q5gc5xxw-3000.app.github.dev', // Проверьте доступность
        changeOrigin: true,
        secure: false
      }
    }
  }
});