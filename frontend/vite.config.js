import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import viteCompression from 'vite-plugin-compression'

export default defineConfig({
  plugins: [
    vue(),
    viteCompression({
      algorithm: 'gzip',
      ext: '.gz',
    }),
    viteCompression({
      algorithm: 'brotliCompress',
      ext: '.br',
    }),
  ],
  server: {
    strictPort: true,
  },
  build: {
    outDir: 'dist',
    // Code splitting
    rollupOptions: {
      output: {
        manualChunks: {
          'vendor': ['vue'],
        },
      },
    },
  },
})
