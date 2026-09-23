import {test,expect} from '@playwright/test';
test.use({ locale: 'es-AR' });
import {createHash} from 'node:crypto';
import {createRequire} from 'node:module';
import {resolve} from 'node:path';
import {pathToFileURL} from 'node:url';
test('network organization agreement and branch lifecycle through actual forms',async({page,context},info)=>{
 if(process.env.ELITE_NETWORK_ROLE_BROWSER!=='1')throw new Error('explicit fixture required');
 const base=process.env.ELITE_BASE_URL;expect(new URL(base).protocol).toBe('https:');expect(new URL(base).hostname).toBe('127.0.0.1');page.setDefaultTimeout(12000);
 const require=createRequire(resolve(process.env.ELITE_WEB_ROOT,'package.json')),{EncryptJWT}=await import(pathToFileURL(require.resolve('jose')).href),identities=JSON.parse(process.env.ELITE_NETWORK_IDENTITIES);
 async function identity(name){const jwt=await new EncryptJWT(identities[name]).setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('300s').encrypt(createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest());await context.clearCookies();await context.addCookies([{name:'__Host-elite_session',value:jwt,url:base+'/',secure:true,httpOnly:true,sameSite:'Lax'}])}
 const errors=[],posts=[];page.on('pageerror',e=>errors.push(e.message));page.on('request',r=>{if(r.method()==='POST'&&r.url().includes('/api/enterprise/network'))posts.push(r.url())});
 const box=name=>page.getByRole('textbox',{name,exact:true}),button=name=>page.getByRole('button',{name,exact:true}),select=name=>page.getByRole('combobox',{name,exact:true});
 async function click(name){await button(name).click();await expect(page.locator('main [role="status"][aria-live="polite"]')).toContainText('Operación registrada.');await expect(button('Consultar estado actual')).toBeEnabled()}
 async function consult(){await button('Consultar estado actual').click();await expect(page.locator('main [role="status"][aria-live="polite"]')).toContainText('Estado actual consultado.');await expect(button('Consultar estado actual')).toBeEnabled()}
 async function drop(action){let dropped=false;await page.route('**/api/enterprise/network**',async route=>{const r=route.request();if(dropped||r.method()!=='POST'||r.postDataJSON().action!==action){await route.continue();return}dropped=true;const response=await route.fetch();expect(response.status()).toBe(200);await route.abort('failed')})}
 async function recover(){await expect(page.locator('main [role="status"][aria-live="polite"]')).toContainText('Resultado sin confirmar');const n=posts.length;await page.unrouteAll();await page.reload();await button('Consultar resultado pendiente').click();await expect(page.locator('main [role="status"][aria-live="polite"]')).toContainText('Resultado recuperado');expect(posts.length).toBe(n)}
 async function organization(type,code,parent,loss=false){await select('Tipo de organización').selectOption(type);if(parent)await select('Organización superior').selectOption(parent);await box('Código de organización').fill(code);await box('Nombre visible').fill('Organización sintética '+code);if(loss){await drop('create-organization');await button('Registrar organización').click();await recover()}else await click('Registrar organización');const id=await select('Registro de la red').inputValue();expect(id).toMatch(/^[a-f0-9-]{36}$/);await expect(button('Activar organización')).toBeDisabled();await consult();await click('Activar organización');await consult();return id}
 await identity('unprivileged');await page.goto('/network');await expect(page.getByText('No tenés permiso para consultar este espacio.',{exact:true})).toBeVisible();expect((await context.request.post(base+'/api/enterprise/network',{headers:{origin:base,'content-type':'application/json'},data:{}})).status()).toBe(403);
 await identity('bootstrap');await page.goto('/network');await expect(page.getByRole('heading',{name:'Red de franquicia',exact:true})).toBeVisible();
 const root=await organization('franchisor','root',null,true);
 const franchise=await organization('franchisee','franchise',root);
 await select('Franquicia del acuerdo').selectOption(franchise);await select('Territorio').selectOption('SYNTHETIC');await select('Condiciones del acuerdo').selectOption('fixture-v1');
 await page.getByLabel('Fecha de inicio',{exact:true}).fill('2026-01-01');await page.getByLabel('Fecha de fin opcional',{exact:true}).fill('2027-01-01');await click('Registrar acuerdo en borrador');
 const agreement=await select('Registro de la red').inputValue();await consult();await click('Activar acuerdo');await consult();
 const branch=await organization('store','branch',franchise);await click('Suspender organización');await consult();await click('Activar organización');await consult();
 await select('Tipo de consulta').selectOption('agreement');await select('Registro de la red').selectOption(agreement);await select('Organización del registro').selectOption(franchise);await consult();await click('Terminar acuerdo');await consult();
 await select('Tipo de consulta').selectOption('organization');await select('Registro de la red').selectOption(branch);await consult();await drop('transition-organization');await button('Cerrar organización').click();await recover();await consult();await expect(page.getByText('Estado: Cerrada. Versión 5.',{exact:true})).toBeVisible();
 await select('Registro de la red').selectOption(franchise);await consult();await click('Cerrar organización');await consult();await expect(page.getByText('Estado: Cerrada. Versión 3.',{exact:true})).toBeVisible();
 await page.screenshot({path:info.outputPath('network-closed-desktop.png'),fullPage:true});await page.setViewportSize({width:390,height:844});expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth)).toBe(true);await page.screenshot({path:info.outputPath('network-closed-mobile.png'),fullPage:true});
 await identity('limited');await page.goto('/network');await expect(page.getByRole('heading',{name:'Registrar acuerdo',exact:true})).toHaveCount(0);await expect(select('Organización superior')).toBeVisible();expect(await select('Organización superior').locator('option[value="'+root+'"]').count()).toBe(0);
 expect((await context.request.get(base+'/api/enterprise/network?'+new URLSearchParams({kind:'organization',id:root,scope_organization_id:root}))).status()).toBe(403);
 expect(posts).toHaveLength(13);expect(errors).toEqual([]);
});
