import { defineConfig } from 'cypress';

export default defineConfig({
  video: false,
  screenshotOnRunFailure: false,
  retries: 2,
  e2e: {
    baseUrl: 'http://localhost:3000',
  },
  defaultCommandTimeout: 8000,
});
