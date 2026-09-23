import { test, expect } from '@playwright/test';
import { createHash } from 'node:crypto';
import { createRequire } from 'node:module';
import { resolve } from 'node:path';
import { pathToFileURL } from 'node:url';

test('customer survey persists and recovers without repeated writes',async({page,context},info)=>{
 test.skip(process.env.ELITE_FEEDBACK_E2E!=='1','Requires disposable survey API and PostgreSQL');
 const identities=JSON.parse(process.env.ELITE_FEEDBACK_IDENTITIES);
 const require=createRequire(resolve(process.env.ELITE_WEB_ROOT,'package.json'));
 const {EncryptJWT}=await import(pathToFileURL(require.resolve('jose')).href);
 const id=info.project.name;
 async function identity(name) {
  const jwt=await new EncryptJWT(identities[name]).setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('300s').encrypt(createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest());
  await context.clearCookies();
  await context.addCookies([{name:'__Host-elite_session',value:jwt,url:process.env.ELITE_BASE_URL,secure:true,httpOnly:true,sameSite:'Lax'}]);
 }
 const customer='/customer/surveys?organizationId=org-a&surveyId='+id;
 const admin='/admin/surveys?organizationId=org-a&surveyId='+id;
 const api='/api/enterprise/surveys?organizationId=org-a&surveyId='+id;
 let writes=0;page.on('request',r=>{if(r.method()==='POST'&&r.url().includes('/api/enterprise/surveys'))writes++;});
 await identity('admin');await page.goto(admin);await expect(page.getByText('El resultado todavía no está disponible',{exact:false})).toBeVisible();
 await identity('alice');await page.goto(customer);
 await expect(page.getByRole('heading',{name:'Synthetic survey <img src=x onerror=alert(1)>'})).toBeVisible();
 await expect(page.locator('section img')).toHaveCount(0);
 await expect(page.getByRole('button',{name:'Enviar respuesta'})).toBeDisabled();
 await page.getByRole('combobox',{name:'Puntuación'}).selectOption('10');
 await expect(page.getByRole('button',{name:'Enviar respuesta'})).toBeDisabled();
 await page.getByRole('checkbox').check();
 // Two synchronous submissions before React can paint a new disabled state.
 await page.locator('form').evaluate(form=>{form.requestSubmit();form.requestSubmit();});
 await expect(page.getByTestId('survey-stored-score')).toHaveText('10');expect(writes).toBe(1);
 await page.reload();await expect(page.getByTestId('survey-stored-score')).toHaveText('10');expect(writes).toBe(1);
 expect((await context.request.get(api+'&view=summary')).status()).toBe(403);
 await page.goto(admin);await expect(page.getByRole('heading',{name:'Acceso denegado'})).toBeVisible();
 await identity('admin');await page.goto(admin);await expect(page.getByText('El resultado todavía no está disponible',{exact:false})).toBeVisible();
 await identity('bob');await page.goto(customer);
 await page.getByRole('combobox',{name:'Puntuación'}).selectOption('0');await page.getByRole('checkbox').check();
 // Complete the real upstream POST, then discard only its response at the browser.
 await page.route('**/api/enterprise/surveys?*',async route=>{
  if(route.request().method()!=='POST'){await route.continue();return;}
  const response=await route.fetch();expect(response.status()).toBe(201);await route.abort('failed');
 });
 await page.getByRole('button',{name:'Enviar respuesta'}).click();
 await expect(page.getByRole('button',{name:'Consultar respuesta guardada'})).toBeVisible();expect(writes).toBe(2);
 await expect(page.getByRole('button',{name:'Enviar respuesta'})).toBeDisabled();
 await page.getByRole('button',{name:'Consultar respuesta guardada'}).click();
 await expect(page.getByTestId('survey-stored-score')).toHaveText('0');expect(writes).toBe(2);
 await page.unrouteAll();await page.reload();await expect(page.getByTestId('survey-stored-score')).toHaveText('0');expect(writes).toBe(2);
 await identity('admin');await page.goto(admin);await expect(page.getByTestId('survey-nps')).toHaveText('0.00');await expect(page.getByText('Respuestas: 2',{exact:true})).toBeVisible();
 expect((await context.request.get(api.replace('org-a','org-b')+'&view=summary')).status()).toBe(403);
 await page.screenshot({path:info.outputPath('survey-results.png'),fullPage:true});
});
