import {test,expect} from '@playwright/test';
import {createHash} from 'node:crypto';
import {createRequire} from 'node:module';
import {resolve} from 'node:path';
import {pathToFileURL} from 'node:url';

test('source-backed gift issue, distinct review and redemption fund actual delivery without a provider',async({page,context},info)=>{
 if(process.env.ELITE_STORED_VALUE_BROWSER!=='1')throw new Error('explicit local handover fixture required');
 const base=process.env.ELITE_BASE_URL;expect(new URL(base).protocol).toBe('https:');expect(new URL(base).hostname).toBe('127.0.0.1');
 const identities=JSON.parse(process.env.ELITE_HANDOVER_IDENTITIES);
 const require=createRequire(resolve(process.env.ELITE_WEB_ROOT,'package.json'));
 const {EncryptJWT}=await import(pathToFileURL(require.resolve('jose')).href);
 async function identity(name){
  const jwt=await new EncryptJWT(identities[name]).setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('300s').encrypt(createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest());
  await context.clearCookies();await context.addCookies([{name:'__Host-elite_session',value:jwt,url:base+'/',secure:true,httpOnly:true,sameSite:'Lax'}]);
 }
 page.setDefaultTimeout(12000);
 const errors=[];page.on('pageerror',error=>errors.push(error.message));
 const writes=[];page.on('request',request=>{if(request.method()==='POST'&&request.url().includes('/api/enterprise/'))writes.push(request.postDataJSON()?.action)});

 async function consult(id){await page.getByLabel('Referencia del pedido',{exact:true}).fill(id);await page.getByRole('button',{name:'Consultar pedido',exact:true}).click();await expect(page.getByText('Total del pedido',{exact:true})).toBeVisible()}
 async function propose(op,order,account){
  await identity('operator');await page.goto('/franchise/stored-value');await consult(order);
  await page.getByRole('combobox',{name:'Operación',exact:true}).selectOption(op);
  if(account)await page.getByLabel('Cuenta de gift card o fidelidad',{exact:true}).fill(account);
  await page.getByRole('button',{name:'Guardar solicitud para revisión',exact:true}).click();
  await expect(page.getByText('Solicitud guardada.',{exact:false})).toBeVisible();
  await expect(page.getByRole('button',{name:'Aprobar contenido revisado',exact:true})).toHaveCount(0);
 }
 async function approve(order){
  await identity('reviewer');await page.goto('/franchise/stored-value');await consult(order);
  await page.getByRole('button',{name:/^Revisar svapproval_/}).click();
  await page.getByLabel('Motivo de la decisión',{exact:true}).fill('Pedido, cuenta y valores exactos revisados en la composición de referencia');
  await page.getByRole('button',{name:'Aprobar contenido revisado',exact:true}).click();
  await expect(page.getByText('Decisión guardada.',{exact:false})).toBeVisible();
 }
 await propose('issue','source');await approve('source');
 const account=(await page.getByRole('heading',{name:/^Cuenta gift_/}).innerText()).replace('Cuenta ','');
 await propose('redeem','order',account);await approve('order');
 await identity('operator');await page.goto('/franchise/stored-value');await consult('order');
 await expect(page.getByRole('button',{name:'Registrar cobertura total para entrega',exact:true})).toBeEnabled();
 await page.route('**/api/enterprise/franchise/stored-value',async route=>{
  if(route.request().method()!=='POST'||route.request().postDataJSON()?.action!=='fund'){await route.continue();return}
  const response=await route.fetch();expect(response.status()).toBe(201);await route.abort('failed');
 });
 await page.getByRole('button',{name:'Registrar cobertura total para entrega',exact:true}).click();
 await expect(page.getByText('Resultado sin confirmar.',{exact:false})).toBeVisible();
 await page.unrouteAll();await page.reload();
 await page.getByRole('button',{name:'Consultar resultado pendiente',exact:true}).click();
 await expect(page.getByRole('heading',{name:'Cobertura total registrada',exact:true})).toBeVisible();
 expect(writes.filter(x=>x==='fund')).toHaveLength(1);
 await page.setViewportSize({width:390,height:844});expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth+1)).toBe(true);
 await page.screenshot({path:info.outputPath('stored-value-recovered-mobile.png'),fullPage:true});await page.setViewportSize({width:1280,height:900});
 await page.getByRole('link',{name:'Continuar con la entrega',exact:true}).click();
 const panel=page.getByRole('article',{name:'Entrega del pedido order',exact:true});
 await expect(panel.getByRole('button',{name:'Preparar entrega',exact:true})).toBeEnabled();
 // Discard only the browser response after the actual Go/PG transaction.
 await page.route('**/api/enterprise/handovers',async route=>{
  const action=route.request().postDataJSON()?.action;
  if(route.request().method()!=='POST'||action!=='prepare'){await route.continue();return}
  const response=await route.fetch();expect(response.status()).toBe(201);await route.abort('failed');
 });
 await panel.getByRole('button',{name:'Preparar entrega',exact:true}).evaluate(button=>{button.click();button.click()});
 await expect(panel.getByText('El resultado quedó sin comprobar.',{exact:false})).toBeVisible();
 expect(writes.filter(x=>x==='prepare')).toHaveLength(1);
 await panel.getByRole('button',{name:'Consultar resultado de preparación',exact:true}).click();
 await expect(panel.getByText('Preparación recuperada:',{exact:false})).toBeVisible();
 await page.unrouteAll();await page.reload();
 const presentation=panel.getByRole('region',{name:'Presentación de entrega',exact:true});
 await expect(presentation.getByLabel('Entrega',{exact:true})).toHaveAttribute('readonly','');
 const handoverID=await presentation.getByLabel('Entrega',{exact:true}).inputValue();
 await expect(presentation.getByLabel('Versión actual de entrega')).toHaveValue('1');
 await presentation.getByLabel('ID del checklist',{exact:true}).fill('browser-checklist');
 await presentation.getByLabel('Versión del checklist',{exact:true}).fill('1');
 await presentation.getByLabel('Respuestas, una por línea').fill('serial=SERIAL-SYNTHETIC');
 await presentation.getByRole('button',{name:'Completar y presentar',exact:true}).click();
 await expect(presentation.getByText('Respuesta recibida.',{exact:false})).toBeVisible();
 await presentation.getByRole('button',{name:'Consultar presentación',exact:true}).click();
 await expect(presentation.getByText('Presentación recuperada:',{exact:false})).toBeVisible();
 await identity('stranger');await page.goto('/customer/handovers');
 await expect(page.getByRole('button',{name:'Registrar recepción',exact:true})).toHaveCount(0);
 const unauthorized=await context.request.post(base+'/api/enterprise/franchise/commands',{headers:{origin:base,'content-type':'application/json'},data:{action:'accept-handover',organizationId:'store',handoverId:handoverID,version:2,confirmedReceived:true,serialNumber:'SERIAL-SYNTHETIC',checklistId:'browser-checklist',checklistVersion:1}});
 expect([404,409]).toContain(unauthorized.status());
 await identity('customer');await page.goto('/customer/handovers');
 expect((await context.request.get(base+'/api/enterprise/handovers?kind=context&organizationId=store&orderId=order')).status()).toBe(403);
 await page.getByLabel('Número de serie observado').fill('SERIAL-SYNTHETIC');
 await page.getByRole('checkbox',{name:'Confirmo que recibí el activo identificado y revisé la preparación indicada'}).check();
 await page.getByRole('button',{name:'Registrar recepción',exact:true}).click();
 await expect(page.getByText('Aceptado:',{exact:false})).toBeVisible();
 await expect(page.getByRole('button',{name:'Registrar recepción',exact:true})).toHaveCount(0);
 await identity('operator');await page.goto('/franchise');
 await expect(panel.getByRole('button',{name:'Registrar cierre comercial',exact:true})).toBeEnabled();
 await page.route('**/api/enterprise/handovers',async route=>{
  if(route.request().method()!=='POST'||route.request().postDataJSON()?.action!=='release'){await route.continue();return}
  const response=await route.fetch();expect(response.status()).toBe(201);await route.abort('failed');
 });
 await panel.getByRole('button',{name:'Registrar cierre comercial',exact:true}).click();
 await expect(panel.getByText('El resultado quedó sin comprobar.',{exact:false})).toBeVisible();
 await page.unrouteAll();

 await panel.getByRole('button',{name:'Consultar resultado comercial',exact:true}).click();
 await expect(panel.getByText('Recibo histórico recuperado.',{exact:false})).toBeVisible();
 await panel.getByRole('button',{name:'Consultar vigencia del recibo',exact:true}).click();
 expect(writes.filter(x=>x==='propose')).toHaveLength(2);expect(writes.filter(x=>x==='decision')).toHaveLength(2);expect(writes.filter(x=>x==='prepare')).toHaveLength(1);expect(writes.filter(x=>x==='release')).toHaveLength(1);
 await identity('customer');expect((await context.request.get(base+'/api/enterprise/franchise/stored-value?kind=order&organization_id=store&order_id=order')).status()).toBe(403);
 await identity('foreign-org');expect((await context.request.get(base+'/api/enterprise/franchise/stored-value?kind=order&organization_id=store&order_id=order')).status()).toBe(403);
 expect(errors).toEqual([]);
 console.log('STORED_VALUE_BROWSER_PASS issue=1 redeem=1 distinct_decisions=2 funding=1 prepare=1 acceptance=1 release=1 recovered_lost_responses=true zero_provider=true');
});
