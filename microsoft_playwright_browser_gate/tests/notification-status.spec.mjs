import { test, expect } from '@playwright/test';
import { createHash } from 'node:crypto';
import { createRequire } from 'node:module';
import { resolve } from 'node:path';
import { pathToFileURL } from 'node:url';

test('notification status connects scoped history, recovery and help without sends', async ({ page, context }) => {
  test.skip(process.env.ELITE_NOTIFICATION_BROWSER_E2E !== '1', 'Requires isolated Go/PG/TLS and synthetic signed identity');
  const require = createRequire(resolve(process.env.ELITE_WEB_ROOT, 'package.json'));
  const { EncryptJWT } = await import(pathToFileURL(require.resolve('jose')).href);
  const key = createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest();
  const base = process.env.ELITE_BASE_URL;
  expect(new URL(base).protocol).toBe('https:'); expect(new URL(base).hostname).toBe('127.0.0.1');
  async function session(permissions) {
    const value = await new EncryptJWT({ subject: 'fixture-operator', tenantId: process.env.ELITE_NOTIFICATION_TENANT, permissions, organizations: [process.env.ELITE_NOTIFICATION_ORGANIZATION], accessToken: process.env.ELITE_NOTIFICATION_TOKEN })
      .setProtectedHeader({ alg: 'dir', enc: 'A256GCM', typ: 'JWT' }).setIssuedAt().setExpirationTime('300s').encrypt(key);
    await context.addCookies([{ name: '__Host-elite_session', value, url: base+'/', httpOnly: true, secure: true, sameSite: 'Lax' }]);
  }
  await session(['appointment:manage']);
  let reads=0, writes=0;
  page.on('request',r=>{if(r.url().includes('/api/enterprise/franchise/notifications'))reads++;if(r.method()==='POST' && r.url().includes('/api/enterprise/'))writes++;});
  await page.goto(base+'/franchise?day='+process.env.ELITE_NOTIFICATION_DAY);
  const panel=page.getByRole('region',{name:'Estado de notificaciones WhatsApp'});
  await expect(panel).toBeVisible();
  await expect(panel.getByRole('status', { name: 'Resultado de consulta de WhatsApp', exact: true })).toBeVisible();
  await expect(panel.locator('input')).toHaveCount(0); expect(reads).toBe(0);
  const button=panel.getByRole('button',{name:'Consultar estado de WhatsApp',exact:true});
  await expect(button).toBeEnabled();
  // Lose only a completed real read response; retry is GET-only, no new send.
  await page.route('**/api/enterprise/franchise/notifications?**',async route=>{const response=await route.fetch();expect(response.status()).toBe(200);await route.abort('failed');},{times:1});
  await button.focus();await page.keyboard.press('Enter');
  await expect(panel.getByRole('status')).toContainText('No pudimos verificar');
  await button.click();
  await expect(panel.getByRole('heading',{name:'El proveedor informó entrega',exact:true})).toBeVisible();
  await expect(panel.getByText('Fecha del aviso del proveedor:', { exact: false })).toBeVisible();
  await expect(panel.locator('time').nth(1)).toHaveAttribute('datetime', '2020-10-19T05:45:14.000Z');
  await expect(panel.getByText('Este aviso no acredita una venta ni la aceptación del turno por el cliente. No vuelve a enviar mensajes.',{exact:true})).toBeVisible();
  await panel.getByText('Ayuda y referencia para soporte',{exact:true}).click();
  await expect(panel.getByText(process.env.ELITE_NOTIFICATION_EVENT,{exact:true})).toBeVisible();
  await expect(panel.getByText('whatsapp-status-view/1.0.0',{exact:true})).toBeVisible();
  expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth)).toBe(true);
  await page.screenshot({path:test.info().outputPath('notification-status-help.png'),fullPage:true});
  // A foreign response must replace success with an error, never stale success.
  await page.route('**/api/enterprise/franchise/notifications?**',async route=>{const response=await route.fetch();const body=await response.json();body.organization_id='foreign';await route.fulfill({response,json:body});},{times:1});
  await button.click();await expect(panel.getByRole('status')).toContainText('No pudimos verificar');
  await expect(panel.getByRole('heading',{name:'El proveedor informó entrega',exact:true})).toHaveCount(0);
  await button.click();await expect(panel.getByRole('heading',{name:'El proveedor informó entrega',exact:true})).toBeVisible();
  const query=new URLSearchParams({organizationId:process.env.ELITE_NOTIFICATION_ORGANIZATION,appointmentId:process.env.ELITE_NOTIFICATION_APPOINTMENT});
  const response=await page.request.get(base+'/api/enterprise/franchise/notifications?'+query);expect(response.status()).toBe(200);expect(response.headers()['cache-control']).toBe('no-store');
  const text=await response.text();for(const forbidden of ['5491112345678','wamid.synthetic','test-app-secret',process.env.ELITE_NOTIFICATION_TOKEN])expect(text).not.toContain(forbidden);
  await session([]);
  expect((await page.request.get(base+'/api/enterprise/franchise/notifications?'+query)).status()).toBe(403);
  await context.clearCookies();expect((await page.request.get(base+'/api/enterprise/franchise/notifications?'+query)).status()).toBe(401);
  expect(writes).toBe(0);
});
