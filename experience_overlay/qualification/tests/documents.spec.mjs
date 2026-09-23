// AUTHORED browser/BFF fixture journey. No OCR, real document corpus or Go/PG claim.
import {target,evidence} from '../paths.mjs';
import {syntheticJpeg} from '../document-fixture.mjs';
import {createRequire} from 'node:module';
import {resolve} from 'node:path';
import {pathToFileURL} from 'node:url';
import {createHash} from 'node:crypto';
import {readFileSync} from 'node:fs';
const require=createRequire(process.env.PLAYWRIGHT_TEST_PACKAGE||resolve(target,'microsoft_playwright_browser_gate/package.json'));
const {test,expect}=require('@playwright/test');
const dep=createRequire(resolve(target,'package.json'));
const {EncryptJWT}=await import(pathToFileURL(dep.resolve('jose')).href);
const base='https://127.0.0.1:4313',backend='http://127.0.0.1:4314',api='/api/enterprise/documents-v403';
const writer={subject:'document-writer',tenantId:'tenant',organizations:['store'],permissions:['documents:write','documents:process'],accessToken:'fixture-document-writer'};
const reviewer={subject:'document-reviewer',tenantId:'tenant',organizations:['store'],permissions:['documents:review'],accessToken:'fixture-document-reviewer'};
async function login(context,identity=writer) {
  const token=await new EncryptJWT(identity).setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('10m').encrypt(createHash('sha256').update(readFileSync(resolve(evidence,'ephemeral-session.key'),'utf8')).digest());
  await context.clearCookies();await context.addCookies([{name:'__Host-elite_session',value:token,url:base+'/',secure:true,httpOnly:true,sameSite:'Lax'}]);
}
const workspace=page=>page.getByRole('region',{name:'Facturas',exact:true});
const stats=async context=>(await context.request.get(backend+'/fixture/documents/stats')).json();
async function start(page,context) {await login(context);await page.goto('/experience/documents');await expect(workspace(page).getByRole('button',{name:'Cargar factura',exact:true})).toBeVisible();}
async function selectFile(page) {await workspace(page).getByLabel('Archivo',{exact:true}).setInputFiles({name:'synthetic-invoice.jpg',mimeType:'image/jpeg',buffer:syntheticJpeg});}
async function upload(page,context) {await selectFile(page);await workspace(page).getByRole('button',{name:'Cargar factura',exact:true}).click();await expect(workspace(page).getByText('Recibida',{exact:true})).toBeVisible();return (await stats(context)).documents[0].document_id;}
async function processDocument(page) {await workspace(page).getByRole('button',{name:'Leer factura',exact:true}).click();await expect(workspace(page).getByLabel('Número de factura',{exact:true})).toHaveValue('SYNTHETIC-001');}
async function editFields(page) {for(const [label,value] of [['Número de factura','REVIEWED-002'],['Proveedor','Reviewed synthetic supplier'],['Total','124.50'],['Moneda','USD']])await workspace(page).getByLabel(label,{exact:true}).fill(value);}
async function submitReview(page) {await workspace(page).getByRole('button',{name:'Enviar a revisión',exact:true}).click();await expect(workspace(page).getByText('Esperando revisión',{exact:true})).toBeVisible();}
async function openReviewer(page,context,id) {await login(context,reviewer);await page.goto('/experience/documents?document='+id);await expect(workspace(page).getByLabel('Motivo de la decisión',{exact:true})).toBeVisible();await workspace(page).getByLabel('Motivo de la decisión',{exact:true}).fill('Reviewed synthetic fixture against supplied evidence');}
async function pendingEntries(page) {return page.evaluate(()=>Object.entries(sessionStorage).filter(([key])=>key.startsWith('elite-document-v403:')&&!key.endsWith(':current')));}
test.beforeEach(async({context})=>{expect((await context.request.post(backend+'/fixture/documents/reset')).ok()).toBe(true);});

test('documents upload, four-field correction, separate reviewer and persisted record',async({page,context})=>{
  await start(page,context);const id=await upload(page,context);await processDocument(page);await editFields(page);await submitReview(page);
  await expect(workspace(page).getByText('Otra persona debe revisar y confirmar esta factura.',{exact:true})).toBeVisible();
  await expect(workspace(page).getByRole('button',{name:'Confirmar factura',exact:true})).toHaveCount(0);
  const proposal=(await stats(context)).documents[0];
  const self=await context.request.post(base+api,{headers:{origin:base,'content-type':'application/json'},data:{action:'decision',id,payload_sha256:proposal.payload_sha256,approved:true,reason:'Forbidden writer approval'}});
  expect(self.status()).toBe(403);expect((await self.json()).effect).toBe('NOT_ATTEMPTED');
  await openReviewer(page,context,id);await expect(workspace(page).getByText('Reviewed synthetic supplier',{exact:true})).toBeVisible();
  const original=await context.request.get(base+api+'?id='+id+'&part=original');expect(original.status()).toBe(200);expect(Buffer.from(await original.body()).equals(syntheticJpeg)).toBe(true);expect(original.headers()['content-disposition']).toContain('attachment');
  for(const part of ['security','provider','analysis']){const response=await context.request.get(base+api+'?id='+id+'&part='+part);expect(response.status()).toBe(200);expect(JSON.parse((await response.body()).toString()).ocr_truth).toBe(false);}
  await workspace(page).getByRole('button',{name:'Confirmar factura',exact:true}).click();await expect(workspace(page).getByText('Confirmada',{exact:true})).toBeVisible();
  const result=await stats(context);expect(result.effects.map(x=>x.action)).toEqual(['receive','process','review','decision']);expect(result.effects[3].actor).toBe('document-reviewer');expect(result.documents[0].proposal.fields).toEqual({invoice_number:'REVIEWED-002',vendor:'Reviewed synthetic supplier',total:'124.50',currency:'USD'});
  expect(result.provider_executed).toBe(false);expect(result.database_executed).toBe(false);expect(await pendingEntries(page)).toEqual([]);
  await page.reload();await expect(workspace(page).getByText('Confirmada',{exact:true})).toBeVisible();expect((await stats(context)).effects).toHaveLength(4);
  await page.setViewportSize({width:390,height:844});await expect(workspace(page).getByText('Confirmada',{exact:true})).toBeVisible();expect(await page.evaluate(()=>document.documentElement.scrollWidth<=innerWidth+1)).toBe(true);
});

test('documents unauthorized, wrong tenant and wrong organization never reach backend',async({context})=>{
  await context.clearCookies();expect((await context.request.get(base+api)).status()).toBe(401);
  for(const delta of [{tenantId:'other-tenant'},{organizations:['other-store']},{permissions:[]}]){
    await login(context,{...writer,...delta});const response=await context.request.get(base+api);expect(response.status()).toBe(403);expect((await response.json()).effect).toBe('NOT_ATTEMPTED');
  }
  const observed=await stats(context);expect(observed.requests).toEqual([]);expect(observed.effects).toEqual([]);
});

test('documents refresh preserves edited fields without storing document contents in browser recovery',async({page,context})=>{
  await start(page,context);await upload(page,context);await processDocument(page);await editFields(page);
  const response=page.waitForResponse(r=>r.url().includes(api+'?id=')&&r.request().method()==='GET');await workspace(page).getByRole('button',{name:'Consultar estado',exact:true}).click();expect((await response).status()).toBe(200);
  await expect(workspace(page).getByLabel('Proveedor',{exact:true})).toBeEnabled();await expect(workspace(page).getByLabel('Proveedor',{exact:true})).toHaveValue('Reviewed synthetic supplier');await expect(workspace(page).getByLabel('Total',{exact:true})).toHaveValue('124.50');
  const storage=await page.evaluate(()=>JSON.stringify(Object.entries(sessionStorage)));expect(storage).not.toContain('Reviewed synthetic supplier');expect(storage).not.toContain('REVIEWED-002');expect(storage).not.toContain('124.50');
  expect((await stats(context)).effects.map(x=>x.action)).toEqual(['receive','process']);
});

async function installHashLatch(page) {
  await page.addInitScript(()=>{
    const original=crypto.subtle.digest.bind(crypto.subtle);
    crypto.subtle.digest=async function(algorithm,data){const bytes=new Uint8Array(data instanceof ArrayBuffer?data:data.buffer,data.byteOffset??0,data.byteLength);const kind=window.__documentHashKind;
      const match=kind==='upload'?bytes[0]===255&&bytes[1]===216:kind==='review'&&new TextDecoder().decode(bytes).startsWith('{"invoice_number":');
      if(match){window.__documentHashKind=null;window.__documentHashStarted=true;await new Promise(resolve=>{window.__releaseDocumentHash=resolve;});}
      return original(algorithm,data);
    };
  });
}

for(const action of ['upload','review'])test('documents double submit during '+action+' hash creates only one effect',async({page,context})=>{
  await installHashLatch(page);await start(page,context);
  if(action==='upload')await selectFile(page);else{await upload(page,context);await processDocument(page);await editFields(page);}
  await page.evaluate(kind=>{window.__documentHashKind=kind;window.__documentHashStarted=false;},action);
  const button=workspace(page).getByRole('button',{name:action==='upload'?'Cargar factura':'Enviar a revisión',exact:true});
  await button.evaluate(node=>{node.click();node.click();});await expect.poll(()=>page.evaluate(()=>window.__documentHashStarted)).toBe(true);
  expect((await stats(context)).effects.filter(x=>x.action===(action==='upload'?'receive':'review'))).toHaveLength(0);
  await page.evaluate(()=>window.__releaseDocumentHash());await expect(workspace(page).getByText(action==='upload'?'Recibida':'Esperando revisión',{exact:true})).toBeVisible();
  const observed=await stats(context);expect(observed.effects.filter(x=>x.action===(action==='upload'?'receive':'review'))).toHaveLength(1);expect(observed.documents).toHaveLength(1);
});

for(const action of ['receive','process','review','decision'])test('documents lost '+action+' response recovers by read after reload without resending',async({page,context})=>{
  await start(page,context);let id;
  if(action==='receive')await selectFile(page);else id=await upload(page,context);
  if(['review','decision'].includes(action)){await processDocument(page);await editFields(page);}
  if(action==='decision'){await submitReview(page);await openReviewer(page,context,id);}
  let lost=false;
  await page.route('**/api/enterprise/documents-v403*',async route=>{
    const request=route.request(),selected=action==='receive'?request.method()==='PUT':request.method()==='POST'&&request.postDataJSON()?.action===action;
    if(selected&&!lost){lost=true;const upstream=await route.fetch();expect(upstream.ok()).toBe(true);await route.abort('failed');}else await route.continue();
  });
  const labels={receive:'Cargar factura',process:'Leer factura',review:'Enviar a revisión',decision:'Confirmar factura'};
  await workspace(page).getByRole('button',{name:labels[action],exact:true}).click();await expect(workspace(page).getByText('Resultado pendiente. Consultá el estado; no vuelvas a cargar ni confirmar.',{exact:true})).toBeVisible();
  expect(lost).toBe(true);expect(await pendingEntries(page)).toHaveLength(1);
  const before=await stats(context);expect(before.effects.filter(x=>x.action===action)).toHaveLength(1);const mutations=before.requests.filter(x=>x.method!=='GET').length;
  await page.reload();const expected={receive:'Recibida',process:'Revisar campos',review:'Esperando revisión',decision:'Confirmada'};
  await expect(workspace(page).getByText(expected[action],{exact:true})).toBeVisible();await expect.poll(()=>pendingEntries(page)).toEqual([]);
  const after=await stats(context);expect(after.requests.filter(x=>x.method!=='GET')).toHaveLength(mutations);expect(after.requests.filter(x=>x.method==='GET').length).toBeGreaterThan(before.requests.filter(x=>x.method==='GET').length);expect(after.effects).toEqual(before.effects);
});
