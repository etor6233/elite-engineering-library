import { test, expect } from '@playwright/test';

test('Microsoft Playwright executes isolated semantic browser assertions', async ({ page }) => {
  await page.setContent(`<!doctype html><html lang="en"><head><title>Elite Browser Gate</title></head>
    <body><main><h1>Browser gate</h1><form><label>Email <input type="email" name="email"></label>
    <button type="submit">Send request</button></form></main></body></html>`);
  await expect(page).toHaveTitle('Elite Browser Gate');
  await expect(page.getByRole('heading', { name: 'Browser gate' })).toBeVisible();
  await page.getByRole('textbox', { name: 'Email' }).fill('browser-gate@example.invalid');
  await expect(page.getByRole('button', { name: 'Send request' })).toBeEnabled();
});
