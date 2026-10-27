import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: './tests/browser',
  timeout: 45000,
  expect: { timeout: 7000 },
  fullyParallel: true,
  workers: 2,
  reporter: 'list',
  use: {
    baseURL: 'http://127.0.0.1:4173/Forum---A-Programming-Community-Platform/',
    viewport: { width: 1440, height: 1000 },
    trace: 'retain-on-failure',
    screenshot: 'only-on-failure',
  },
  webServer: {
    command: 'npm run preview',
    url: 'http://127.0.0.1:4173/Forum---A-Programming-Community-Platform/',
    reuseExistingServer: false,
    timeout: 30000,
  },
});
