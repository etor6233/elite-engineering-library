import { test, expect } from '@playwright/test';

function nonceFrom(policy) {
  const match = /'nonce-([^']+)'/.exec(policy || '');
  if (!match) throw new Error('CSP nonce is missing');
  return match[1];
}

test('strict CSP has a fresh nonce and every rendered script carries it', async ({ page, browserName }) => {
  if (browserName === 'webkit') {
    await page.route(/^https:\/\/127\.0\.0\.1:4173\//, async route => {
      const response = await route.fetch({ url: route.request().url().replace(/^https:/, 'http:') });
      await route.fulfill({ response });
    });
  }
  const first = await page.goto('/', { waitUntil: 'networkidle' });
  expect(first?.status()).toBe(200);
  const firstPolicy = first?.headers()['content-security-policy'];
  expect(firstPolicy).not.toContain("'unsafe-inline'");
  expect(firstPolicy).not.toContain("'unsafe-eval'");
  expect(firstPolicy).toContain("'strict-dynamic'");
  const firstNonce = nonceFrom(firstPolicy);
  expect(firstNonce).toMatch(/^[A-Za-z0-9+/]+=*$/);
  expect(firstNonce.length).toBeGreaterThanOrEqual(16);
  const scriptNonces = await page.locator('script').evaluateAll(nodes => nodes.map(node => node.nonce));
  expect(scriptNonces.length).toBeGreaterThan(0);
  expect(scriptNonces.every(value => value === firstNonce)).toBe(true);
  const second = await page.reload({ waitUntil: 'networkidle' });
  const secondNonce = nonceFrom(second?.headers()['content-security-policy']);
  expect(secondNonce).not.toBe(firstNonce);
  const secondScriptNonces = await page.locator('script').evaluateAll(nodes => nodes.map(node => node.nonce));
  expect(secondScriptNonces.every(value => value === secondNonce)).toBe(true);
});
