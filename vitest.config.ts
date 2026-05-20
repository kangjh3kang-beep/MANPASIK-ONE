import { defineConfig } from 'vitest/config';
import path from 'path';
import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));

export default defineConfig({
  test: {
    globals: true,
    environment: 'jsdom',
    include: [
      'packages/*/src/**/*.test.{ts,tsx}',
    ],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'html'],
      include: ['packages/*/src/**/*.{ts,tsx}'],
      exclude: ['**/*.d.ts', '**/*.test.{ts,tsx}', '**/index.ts'],
    },
  },
  resolve: {
    alias: {
      '@mmup/api-client': path.resolve(__dirname, 'packages/api-client/src'),
      '@mmup/store': path.resolve(__dirname, 'packages/store/src'),
      '@mmup/auth': path.resolve(__dirname, 'packages/auth/src'),
    },
  },
});
