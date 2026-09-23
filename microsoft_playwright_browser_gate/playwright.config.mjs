import { defineConfig, devices } from '@playwright/test';
import { resolve } from 'node:path';

const webRoot = resolve(process.env.ELITE_WEB_ROOT || '..');
const suppliedBaseURL = process.env.ELITE_BASE_URL;
const runtimeOnly = process.env.ELITE_RUNTIME_ONLY === '1';
const baseURL = suppliedBaseURL || 'http://127.0.0.1:4173';
const parsed = new URL(baseURL);
const loopback = parsed.hostname === '127.0.0.1' || parsed.hostname === 'localhost' || parsed.hostname === '::1';
if (parsed.username || parsed.password || (parsed.protocol !== 'https:' && !(parsed.protocol === 'http:' && loopback)))
  throw new Error('ELITE_BASE_URL must be HTTPS or loopback HTTP and contain no credentials');

export default defineConfig({
  testDir: './tests',
  fullyParallel: false,
  forbidOnly: true,
  retries: 0,
  workers: 1,
  reporter: [['list']],
  outputDir: 'test-results',
  use: {
    baseURL,
    // Self-signed certificate only in the explicit loopback operator fixture.
    ignoreHTTPSErrors: loopback && (process.env.ELITE_OPERATOR_AGENDA_E2E === '1' || process.env.ELITE_NOTIFICATION_BROWSER_E2E === '1' || process.env.ELITE_QUOTE_E2E === '1' || process.env.ELITE_HELP_E2E === '1' || ['enabled','disabled'].includes(process.env.ELITE_WORKSPACE_E2E)),
    trace: 'retain-on-failure',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
  },
  webServer: runtimeOnly || suppliedBaseURL ? undefined : {
    command: 'pnpm exec next start --hostname 127.0.0.1 --port 4173',
    cwd: webRoot,
    url: baseURL,
    reuseExistingServer: false,
    timeout: 120_000,
  },
  projects: [
    { name: 'chromium-desktop', use: { ...devices['Desktop Chrome'] } },
    { name: 'chromium-mobile', use: { ...devices['Pixel 7'] } },
    { name: 'firefox-desktop', use: { ...devices['Desktop Firefox'] } },
    { name: 'webkit-desktop', use: { ...devices['Desktop Safari'] } },
  ],
});
