import {test,expect} from '@playwright/test';
import {createHash} from 'node:crypto';
import {createRequire} from 'node:module';
import {resolve} from 'node:path';
import {pathToFileURL} from 'node:url';
test('sold terms and connected warranty through role forms',async({page,context},info)=>{
 if(process.env.ELITE_WARRANTY_ROLE_BROWSER!=='1')throw new Error('explicit fixture required');
 const base=process.env.ELITE_BASE_URL;expect(new URL(base).protocol).toBe('https:');expect(new URL(base).hostname).toBe('127.0.0.1');page.setDefaultTimeout(12000);
 const require=createRequire(resolve(process.env.ELITE_WEB_ROOT,'package.json')),{EncryptJWT}=await import(pathToFileURL(require.resolve('jose')).href),identities=JSON.parse(process.env.ELITE_WARRANTY_IDENTITIES),inputs=JSON.parse(process.env.ELITE_WARRANTY_INPUTS);
 async function identity(name){const jwt=await new EncryptJWT(identities[name]).setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('300s').encrypt(createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest());await context.clearCookies();await context.addCookies([{name:'__Host-elite_session',value:jwt,url:base+'/',secure:true,httpOnly:true,sameSite:'Lax'}])}
 const errors=[],posts=[];page.on('pageerror',e=>errors.push(e.message));page.on('request',r=>{if(r.method()==='POST'&&r.url().includes('/api/enterprise/warranty'))posts.push(r.url())});
 const textbox=name=>page.getByRole('textbox',{name,exact:true});
 const button=name=>page.getByRole('button',{name,exact:true});
 async function evidence(suffix){await page.getByLabel('Archivo de evidencia',{exact:true}).setInputFiles({name:'fixture.txt',mimeType:'text/plain',buffer:Buffer.from('Authorized synthetic evidence '+suffix)});await expect(page.getByText('Evidencia preparada.',{exact:true})).toBeVisible();await textbox('Motivo de la decisión').fill('Revisión de la evidencia sintética '+suffix)}
 async function open(name,params={}){await identity(name);await page.goto('/warranty?'+new URLSearchParams(params));await expect(page.getByRole('heading',{name:'Garantía y reparación',exact:true})).toBeVisible();if(params.case){await button('Consultar caso').click();await expect(page.getByRole('status')).toContainText('Consulta actualizada.')}await evidence(name)}
 async function click(name){await button(name).click();await expect(page.getByRole('status')).toContainText('Operación registrada.');await expect(button('Consultar caso')).toBeEnabled()}
 async function consult(name){await button(name).click();await expect(page.getByRole('status')).toContainText('Consulta actualizada.');await expect(button(name)).toBeEnabled()}
 async function drop(action){let dropped=false;await page.route('**/api/enterprise/warranty**',async route=>{const r=route.request();if(dropped||r.method()!=='POST'||r.postDataJSON().action!==action){await route.continue();return}dropped=true;const response=await route.fetch();expect(response.status()).toBe(200);await route.abort('failed')})}
 async function recover(){await expect(page.getByRole('status')).toContainText('Resultado sin confirmar');const n=posts.length;await page.unrouteAll();await page.reload();await button('Consultar resultado pendiente').click();await expect(page.getByRole('status')).toContainText('Resultado recuperado');await expect(button('Consultar caso')).toBeEnabled();expect(posts.length).toBe(n)}
 if(process.env.ELITE_WARRANTY_PHASE==='terms'){
  await identity('unprivileged');await page.goto('/warranty');await expect(page.getByText('No tenés permiso para consultar este espacio.',{exact:true})).toBeVisible();
  await identity('reader');expect((await context.request.post(base+'/api/enterprise/warranty',{headers:{origin:base,'content-type':'application/json'},data:{}})).status()).toBe(403);
  await open('operator',{quote:inputs.quote});await consult('Consultar cotización para términos');await consult('Consultar política configurada');await drop('offer');await button('Vincular términos revisados').click();await recover();await expect(page.getByText('Recepción de términos pendiente.',{exact:true})).toBeVisible();
  await open('customer',{quote:inputs.quote});await consult('Consultar términos ofrecidos');await expect(button('Registrar recepción de términos')).toBeDisabled();await page.getByRole('checkbox',{name:'Leí los términos mostrados de esta cotización',exact:true}).check();await click('Registrar recepción de términos');await expect(page.getByText('Recepción de términos registrada.',{exact:true})).toBeVisible();
  expect(posts).toHaveLength(2);
 }else{
  expect(process.env.ELITE_WARRANTY_PHASE).toBe('claims');
  await open('operator',{handover:inputs.handover});await click('Activar garantía de entrega aceptada');await expect(page.getByRole('heading',{name:'Garantía activa',exact:true})).toBeVisible();
  async function claim(appointment){await open('customer',{handover:inputs.handover});await consult('Consultar garantía activa');await textbox('Referencia de la atención completada').fill(appointment);await textbox('Descripción del problema').fill('Problema sintético de freno');await click('Abrir reclamo');const id=await textbox('Referencia del caso').inputValue();expect(id).toMatch(/^[a-f0-9-]{36}$/);return id}
  async function diagnosis(id,fault){await open('operator',{case:id});await textbox('Código de falla').fill(fault);await textbox('Descripción del diagnóstico').fill('Inspección sintética de freno');await click('Registrar diagnóstico')}
  async function plan(withParts){await textbox('Trabajo propuesto').fill('Ajuste sintético de freno');if(withParts){await button('Agregar repuesto').click();for(const [name,value]of [['Artículo','warranty-part'],['Ubicación de stock','warranty-bin'],['Lote opcional','warranty-lot'],['Cantidad de repuesto','7']])await textbox(name).fill(value)}await click('Proponer reparación');await expect(button('Aprobar reparación')).toBeDisabled()}
  const main=await claim(inputs.main);await diagnosis(main,'fixture-wear');await plan(true);
  await open('reviewer',{case:main});await expect(page.getByRole('heading',{name:'Plan de reparación',exact:true})).toBeVisible();await expect(page.getByText('warranty-part: 7 · warranty-bin · lote warranty-lot',{exact:true})).toBeVisible();await click('Aprobar reparación');
  await open('operator',{case:main});await drop('work');await button('Registrar trabajo realizado').click();await recover();await evidence('after-work');await expect(button('Aprobar calidad del trabajo')).toBeDisabled();
  await open('quality',{case:main});await click('Registrar calidad fallida');await expect(button('Aprobar calidad del trabajo')).toBeDisabled();await page.getByLabel('Evidencia de corrección',{exact:true}).setInputFiles({name:'correction.txt',mimeType:'text/plain',buffer:Buffer.from('Synthetic corrected brake')});await expect(button('Aprobar calidad del trabajo')).toBeEnabled();await click('Aprobar calidad del trabajo');
  await open('customer',{case:main});await click('Aceptar reparación recibida');await expect(page.getByText('Reparación aceptada por el cliente.',{exact:true})).toBeVisible();
  await open('factory',{case:main});await click('Conciliar servicio de garantía');await expect(page.getByText('Cerrado',{exact:true})).toBeVisible();await expect(page.getByText('Costo de inventario registrado: 740.0000.',{exact:true})).toBeVisible();
  await page.screenshot({path:info.outputPath('warranty-closed-desktop.png'),fullPage:true});await page.setViewportSize({width:390,height:844});expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth)).toBe(true);await page.screenshot({path:info.outputPath('warranty-closed-mobile.png'),fullPage:true});
  const cancelled=await claim(inputs.cancel);await open('operator',{case:cancelled});await click('Cancelar reparación sin trabajo emitido');await expect(page.getByText('Cancelado',{exact:true})).toBeVisible();
  const excluded=await claim(inputs.excluded);await diagnosis(excluded,'fixture-exclusion');await plan(false);await open('reviewer',{case:excluded});await expect(button('Aprobar reparación')).toBeDisabled();await click('Rechazar reparación');await expect(page.getByText('Cancelado',{exact:true})).toBeVisible();
  await identity('other-customer');expect((await context.request.get(base+'/api/enterprise/warranty?'+new URLSearchParams({kind:'claim',organization_id:'store',id:main,surface:'customer'}))).status()).toBe(409);
  await identity('foreign');expect((await context.request.get(base+'/api/enterprise/warranty?'+new URLSearchParams({kind:'claim',organization_id:'foreign',id:main,surface:'franchise'}))).status()).toBe(409);
  expect(posts).toHaveLength(16);
 }
 expect(errors).toEqual([]);
});
