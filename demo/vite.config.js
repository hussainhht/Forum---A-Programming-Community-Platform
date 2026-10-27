import { defineConfig } from 'vite';
import { fileURLToPath } from 'node:url';

export default defineConfig({
  base: '/Forum---A-Programming-Community-Platform/',
  server: {
    port: 5173,
    strictPort: true,
    fs: { allow: [fileURLToPath(new URL('..', import.meta.url))] },
  },
  preview: { port: 4173, strictPort: true },
});
