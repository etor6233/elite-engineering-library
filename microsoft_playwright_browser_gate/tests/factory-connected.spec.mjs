import {test,expect} from '@playwright/test';
import {createHash} from 'node:crypto';
import {createRequire} from 'node:module';
import {resolve} from 'node:path';
import {pathToFileURL} from 'node:url';

test('factory operator records progress and observes concurrent state through real BFF API and PostgreSQL',async({page,context},info)=>{
 if(process.env.ELITE_FACTORY_BROWSER!=='1')throw new Error('explicit local factory fixture required');
 const base=process.env.ELITE_BASE_URL;expect(new URL(base).protocol).toBe('https:');expect(new URL(base).hostname).toBe('127.0.0.1');
 const identities=JSON.parse(process.env.ELITE_FACTORY_IDENTITIES);
 const require=createRequire(resolve(process.env.ELITE_WEB_ROOT,'package.json'));
 const {EncryptJWT}=await import(pathToFileURL(require.resolve('jose')).href);
 async function identity(name){
  const jwt=await new EncryptJWT(identities[name]).setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('300s').encrypt(createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest());
  await context.clearCookies();await context.addCookies([{name:'__Host-elite_session',value:jwt,url:base+'/',secure:true,httpOnly:true,sameSite:'Lax'}]);
 }

 const errors=[];page.on('pageerror',error=>errors.push(error.message));let pagePosts=0;page.on('request',request=>{if(request.method()==='POST'&&request.url().includes('/api/enterprise/factory'))pagePosts++});
 const body={organizationId:'store',unitId:'unit',action:'start-assembly',expectedState:'planned'};
 async function post(data=body){return context.request.post(base+'/api/enterprise/factory',{headers:{origin:base,'content-type':'application/json'},data})}
 await identity('reader');await page.goto('/factory');await expect(page.getByRole('button',{name:'Registrar inicio de ensamblaje'})).toHaveCount(0);expect((await post()).status()).toBe(403);
 await identity('foreign-org');await page.goto('/factory');await expect(page.getByText('SERIAL-FACTORY-FIXTURE',{exact:true})).toHaveCount(0);expect((await post()).status()).toBe(403);expect((await context.request.get(base+'/api/enterprise/factory?organizationId=other&unitId=unit')).status()).toBe(404);
 await identity('foreign-tenant');expect((await context.request.get(base+'/api/enterprise/factory?organizationId=store&unitId=unit')).status()).toBe(404);
 await identity('operator');await page.goto('/factory');const panel=page.getByRole('region',{name:'Seguimiento de SERIAL-FACTORY-FIXTURE'});await expect(panel.getByText('planned',{exact:true})).toBeVisible();
 await page.route('**/api/enterprise/factory',async route=>{if(route.request().method()!=='POST'){await route.continue();return};const response=await route.fetch();expect(response.status()).toBe(200);await route.abort('failed')});
 await panel.getByRole('button',{name:'Registrar inicio de ensamblaje'}).evaluate(button=>{button.click();button.click()});await expect(panel.getByText('La respuesta no quedó confirmada.',{exact:false})).toBeVisible();expect(pagePosts).toBe(1);
 // A separate actor advances the same unit before recovery: state alone cannot attribute the lost request.
 expect((await context.request.post(process.env.ELITE_FACTORY_CONTROL,{headers:{'X-Fixture-Token':process.env.ELITE_FACTORY_CONTROL_TOKEN}})).status()).toBe(200);
 await page.unrouteAll();await page.reload();await expect(panel.getByText('Hay una respuesta pendiente de comprobar.')).toBeVisible();await panel.getByRole('button',{name:'Consultar estado actual'}).click();await expect(panel.getByText('quality',{exact:true})).toBeVisible();await expect(panel.getByText('Estado actual consultado. No identifica la solicitud que lo produjo.')).toBeVisible();expect(pagePosts).toBe(1);
 expect((await post()).status()).toBe(409);await expect(panel.getByRole('button',{name:'Registrar inicio de ensamblaje'})).toHaveCount(0);
 const observed=await context.request.get(base+'/api/enterprise/factory?organizationId=store&unitId=unit');expect(observed.status()).toBe(200);expect(await observed.json()).toMatchObject({unit:{state:'quality'},observationOnly:true});
 await page.setViewportSize({width:390,height:844});expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth)).toBe(true);await page.screenshot({path:info.outputPath('factory-recovered-mobile.png'),fullPage:true});expect(errors).toEqual([]);
});
