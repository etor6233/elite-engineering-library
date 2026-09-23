import {test,expect} from '@playwright/test';
import {createHash} from 'node:crypto';
import {createRequire} from 'node:module';
import {resolve} from 'node:path';
import {pathToFileURL} from 'node:url';
test('catalog role source to publication and rollback preserves command identities after response loss',async({page,context},info)=>{
 if(process.env.ELITE_CATALOG_ROLE_BROWSER!=='1')throw new Error('explicit local fixture required');
 page.setDefaultTimeout(12000);
 const base=process.env.ELITE_BASE_URL;expect(new URL(base).protocol).toBe('https:');expect(new URL(base).hostname).toBe('127.0.0.1');
 const require=createRequire(resolve(process.env.ELITE_WEB_ROOT,'package.json')),{EncryptJWT}=await import(pathToFileURL(require.resolve('jose')).href);
 const identities=JSON.parse(process.env.ELITE_CATALOG_IDENTITIES);
 async function identity(name){const jwt=await new EncryptJWT(identities[name]).setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('300s').encrypt(createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest());await context.clearCookies();await context.addCookies([{name:'__Host-elite_session',value:jwt,url:base+'/',secure:true,httpOnly:true,sameSite:'Lax'}])}
 const errors=[];page.on('pageerror',e=>errors.push(e.message));const posts=[];page.on('request',r=>{if(r.method()==='POST'&&r.url().includes('/api/enterprise/catalog'))posts.push(r.url())});
 async function editor(draft=''){await page.goto('/admin/catalog'+(draft?'?draft='+encodeURIComponent(draft):''));await expect(page.getByRole('heading',{name:'Edición y publicación del catálogo'})).toBeVisible()}
 async function dropOnce(action){
  let dropped=false;
  await page.route('**/api/enterprise/catalog**',async route=>{
   const r=route.request();if(dropped||r.method()!=='POST'||r.headers()['content-type']!=='application/json'||r.postDataJSON().action!==action){await route.continue();return}
   dropped=true;const response=await route.fetch();expect(response.status()).toBe(200);await route.abort('failed');
  });
 }
 async function recover(){await expect(page.getByRole('status')).toContainText('Resultado sin confirmar');const before=posts.length;await page.unrouteAll();await page.reload();await page.getByRole('button',{name:'Consultar resultado pendiente',exact:true}).click();await expect(page.getByRole('status')).toContainText('Resultado recuperado');expect(posts.length).toBe(before)}
 await identity('unprivileged');await page.goto('/admin/catalog');await expect(page.getByText('No tenés permiso para consultar este espacio.')).toBeVisible();
 await identity('foreign-org');await editor();await page.getByRole('button',{name:'Consultar publicación vigente',exact:true}).click();await expect(page.getByRole('status')).toContainText('No pudimos confirmar');await expect(page.getByText('Todavía no hay una publicación.',{exact:true})).toHaveCount(0);
 await identity('reader');expect((await context.request.post(base+'/api/enterprise/catalog',{headers:{origin:base,'content-type':'application/json'},data:{}})).status()).toBe(403);
 async function createVersion(suffix,name,amount,loseSource,upload){
  await identity('maker');await editor();
  await page.getByLabel('Código del modelo',{exact:true}).fill('role-'+suffix);
  await page.getByLabel('Nombre del modelo',{exact:true}).fill(name);
  await page.getByRole('combobox',{name:'Clase de vehículo',exact:true}).selectOption('bicycle');
  await page.getByLabel('Descripción técnica',{exact:true}).fill('Ficha técnica de referencia <verificada> & sintética.');
  await page.getByLabel('Vigencia desde (UTC)',{exact:true}).fill(process.env.ELITE_CATALOG_VALID_FROM);
  await page.getByLabel('Código de variante',{exact:true}).fill('variant-'+suffix);
  await page.getByLabel('Nombre de variante',{exact:true}).fill('Variante '+suffix);
  await page.getByLabel('Especificación de batería',{exact:true}).fill('Batería de referencia, sin datos privados.');
  await page.getByLabel('Importe en unidades menores',{exact:true}).fill(amount);
  await page.getByRole('combobox',{name:'Tratamiento del precio',exact:true}).selectOption('not-applicable');
  if(loseSource)await dropOnce('source');
  await page.getByRole('button',{name:'Guardar modelo y precios',exact:true}).click();
  if(loseSource)await recover();else await expect(page.getByRole('status')).toContainText('Modelo, variantes y precios guardados');
  expect(await page.getByLabel('Referencia del modelo',{exact:true}).inputValue()).not.toBe('');
  if(upload){
   await page.getByLabel('Imagen PNG',{exact:true}).setInputFiles({name:'reference.png',mimeType:'image/png',buffer:Buffer.from(process.env.ELITE_CATALOG_PNG,'base64')});
   await page.getByRole('button',{name:'Guardar imagen',exact:true}).click();await expect(page.getByRole('status')).toContainText('Imagen recibida');
  }
  expect(await page.getByLabel('Referencia de la imagen',{exact:true}).inputValue()).not.toBe('');
  await page.getByRole('button',{name:'Guardar versión para revisión',exact:true}).click();await expect(page.getByRole('status')).toContainText('Versión guardada');
  const id=await page.getByLabel('Referencia del borrador',{exact:true}).inputValue();expect(id).not.toBe('');
  await page.getByRole('button',{name:'Consultar borrador',exact:true}).click();await expect(page.getByRole('heading',{name,exact:true})).toBeVisible();
  await expect(page.getByText('Otra persona autorizada debe revisar y publicar esta versión.',{exact:true})).toBeVisible();await expect(page.getByRole('button',{name:'Aprobar etapa revisada',exact:true})).toHaveCount(0);
  return id;
 }
 async function review(id,expectedName){
  await identity('reviewer');await editor(id);
  await page.getByRole('button',{name:'Consultar borrador',exact:true}).click();await expect(page.getByRole('heading',{name:expectedName,exact:true})).toBeVisible();
  await page.getByLabel('Evidencia de la revisión',{exact:true}).setInputFiles({name:'review.txt',mimeType:'text/plain',buffer:Buffer.from('Synthetic evidence of explicit reference review.')});
  for(const stage of ['legal','technical','media','publication']){
   await page.getByRole('combobox',{name:'Etapa de revisión',exact:true}).selectOption(stage);
   await page.getByRole('textbox',{name:'Motivo de la decisión',exact:true}).fill('Revisé el contenido sintético para la etapa '+stage+'.');
   await page.getByRole('button',{name:'Aprobar etapa revisada',exact:true}).click();
   const label={legal:'Legal',technical:'Técnica',media:'Imagen',publication:'Publicación'}[stage];
   await expect(page.getByLabel('Estado de revisiones').getByText(label+': approved',{exact:true})).toBeVisible();
  }
 }
 async function publish(name,generation,lose=false){
  await page.getByRole('button',{name:'Consultar publicación vigente',exact:true}).click();await expect(page.getByRole('status')).toContainText('Publicación vigente consultada');
  await page.getByRole('textbox',{name:'Motivo de publicación',exact:true}).fill('Publicar la versión revisada de referencia.');
  if(lose)await dropOnce('publish');
  await page.getByRole('button',{name:'Publicar versión revisada',exact:true}).click();
  if(lose)await recover();else await expect(page.getByRole('status')).toContainText('Publicación registrada');
  await page.getByRole('button',{name:'Consultar publicación vigente',exact:true}).click();await expect(page.getByText('Versión publicada '+generation+'.',{exact:true})).toBeVisible();
  await page.goto('/models');await expect(page.getByRole('heading',{name,exact:true})).toBeVisible();
 }
 const first=await createVersion('one','Modelo Rol Uno','123456',true,true);
 await review(first,'Modelo Rol Uno');await publish('Modelo Rol Uno','1',true);
 const second=await createVersion('two','Modelo Rol Dos','234567',false,false);
 await review(second,'Modelo Rol Dos');await publish('Modelo Rol Dos','2');
 await editor(first);await page.getByRole('button',{name:'Consultar borrador',exact:true}).click();await expect(page.getByRole('heading',{name:'Modelo Rol Uno',exact:true})).toBeVisible();
 await publish('Modelo Rol Uno','3');await expect(page.getByRole('heading',{name:'Modelo Rol Dos',exact:true})).toHaveCount(0);
 await page.screenshot({path:info.outputPath('catalog-role-published-desktop.png'),fullPage:true});
 await editor(first);await page.getByRole('button',{name:'Consultar borrador',exact:true}).click();await page.setViewportSize({width:390,height:844});
 expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth)).toBe(true);
 await page.screenshot({path:info.outputPath('catalog-role-review-mobile.png'),fullPage:true});expect(errors).toEqual([]);
 // 16 successful writes; the deliberate reader POST was rejected by the BFF.
 expect(posts).toHaveLength(16);
});
