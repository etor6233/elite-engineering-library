import {test,expect} from '@playwright/test';
import {createHash} from 'node:crypto';
import {createRequire} from 'node:module';
import {resolve} from 'node:path';
import {pathToFileURL} from 'node:url';

test('human reviews WhatsApp proposals and reconciles an uncertain send through real BFF Go and PostgreSQL',async({page,context},info)=>{
 if(process.env.ELITE_WHATSAPP_BROWSER!=='1')throw new Error('explicit local WhatsApp fixture required');
 const base=process.env.ELITE_BASE_URL;expect(new URL(base).protocol).toBe('https:');expect(new URL(base).hostname).toBe('127.0.0.1');
 const identities=JSON.parse(process.env.ELITE_WHATSAPP_IDENTITIES);
 const require=createRequire(resolve(process.env.ELITE_WEB_ROOT,'package.json'));
 const {EncryptJWT}=await import(pathToFileURL(require.resolve('jose')).href);
 async function identity(name){
  const jwt=await new EncryptJWT(identities[name]).setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('300s').encrypt(createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest());
  await context.clearCookies();await context.addCookies([{name:'__Host-elite_session',value:jwt,url:base+'/',secure:true,httpOnly:true,sameSite:'Lax'}]);
 }
 const errors=[];page.on('pageerror',error=>errors.push(error.message));let sends=0,recoveries=0;
 page.on('request',r=>{if(r.method()==='POST'&&r.url().includes('/api/enterprise/franchise/whatsapp/replies')){const b=r.postDataJSON();if(b.action==='send')sends++;if(b.action==='recover')recoveries++}});
 await identity('reader');await page.goto('/franchise');await expect(page.getByRole('link',{name:'Revisar respuestas de WhatsApp'})).toHaveCount(0);
 await page.goto('/franchise/whatsapp');await expect(page.getByText('Tu sesión no permite revisar respuestas.')).toBeVisible();
 await identity('foreign-org');await page.goto('/franchise/whatsapp');await expect(page.getByRole('alert').filter({hasText:'No pudimos consultar'})).toBeVisible();await expect(page.locator('article')).toHaveCount(0);
 await identity('human');await page.goto('/franchise');await page.getByRole('link',{name:'Revisar respuestas de WhatsApp'}).click();
 await expect(page.getByRole('heading',{name:'Respuestas de WhatsApp',exact:true})).toBeVisible();await expect(page.locator('article')).toHaveCount(3);
 const rejected=page.locator('article').nth(0),lost=page.locator('article').nth(1),normal=page.locator('article').nth(2);
 await expect(normal).toContainText('Cotización preparada para revisar.');await expect(normal).toContainText('5491112345678');
 await rejected.getByRole('button',{name:'Rechazar',exact:true}).click();await expect(rejected).toContainText('Revisión: Rechazada');await expect(rejected.getByRole('button')).toHaveCount(0);
 await normal.getByRole('button',{name:'Aprobar este texto'}).click();await expect(normal).toContainText('Revisión: Aprobada');expect(sends).toBe(0);
 await identity('review-only');await page.reload();await expect(normal.getByRole('button',{name:'Enviar respuesta aprobada'})).toHaveCount(0);
 await identity('human');await page.reload();await normal.getByRole('button',{name:'Enviar respuesta aprobada'}).click();await expect(normal).toContainText('Envío: Aceptado por WhatsApp');expect(sends).toBe(1);
 expect((await context.request.post(process.env.ELITE_WHATSAPP_CONTROL,{headers:{'X-Fixture-Token':process.env.ELITE_WHATSAPP_CONTROL_TOKEN},data:{index:1}})).status()).toBe(200);
 await page.getByRole('link',{name:'Consultar estado',exact:true}).click();await expect(normal).toContainText('Entregado según WhatsApp');
 await lost.getByRole('button',{name:'Aprobar este texto'}).click();await expect(lost).toContainText('Revisión: Aprobada');
 await lost.getByRole('button',{name:'Enviar respuesta aprobada'}).click();await expect(page.getByRole('status')).toContainText('La acción no está confirmada');expect(sends).toBe(2);
 await expect(lost.getByRole('button',{name:'Enviar respuesta aprobada'})).toBeDisabled();
 await page.getByRole('link',{name:'Consultar estado',exact:true}).click();await expect(lost).toContainText('El resultado del envío necesita reconciliación.');await expect(lost.getByRole('button',{name:'Enviar respuesta aprobada'})).toHaveCount(0);
 await page.screenshot({path:info.outputPath('whatsapp-reconciliation-desktop.png'),fullPage:true});
 await lost.getByRole('button',{name:'Verificar comprobante guardado'}).click();await expect(lost).toContainText('Envío: Aceptado por WhatsApp');expect(recoveries).toBe(1);expect(sends).toBe(2);
 expect((await context.request.post(process.env.ELITE_WHATSAPP_CONTROL,{headers:{'X-Fixture-Token':process.env.ELITE_WHATSAPP_CONTROL_TOKEN},data:{index:2}})).status()).toBe(200);
 await page.getByRole('link',{name:'Consultar estado',exact:true}).click();await expect(lost).toContainText('Entregado según WhatsApp');await expect(normal).toContainText('Entregado según WhatsApp');
 await page.setViewportSize({width:390,height:844});expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth)).toBe(true);await page.screenshot({path:info.outputPath('whatsapp-recovered-mobile.png'),fullPage:true});expect(errors).toEqual([]);
});
