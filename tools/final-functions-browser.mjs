// AUTHORED qualification: no route interception; actual Next BFF and Go/PG owner.
import {createHash} from 'node:crypto';
import {writeFile} from 'node:fs/promises';
import {pathToFileURL} from 'node:url';
import {EncryptJWT} from 'jose';
const {chromium,expect}=await import(pathToFileURL(process.env.ELITE_PLAYWRIGHT_MODULE).href);
const kind=process.env.FUNCTIONS_BROWSER_KIND,base=process.env.FUNCTIONS_BROWSER_ORIGIN,out=process.env.FUNCTIONS_BROWSER_OUTPUT;
const identity=JSON.parse(process.env.FUNCTIONS_IDENTITY),key=createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest();
const seal=async(p)=>new EncryptJWT(p).setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('3600s').encrypt(key);
const browser=await chromium.launch({headless:true});const context=await browser.newContext({viewport:{width:1440,height:1000}}),page=await context.newPage();
const checks=[],responses=[];page.on('response',r=>{if(r.url().includes('/api/enterprise/'))responses.push({url:new URL(r.url()).pathname,status:r.status()})});
const route=kind==='ai'?'/intelligence/activity':'/franchise/campaigns';
async function cookie(p){await context.addCookies([{name:'elite_session',value:await seal(p),url:base,httpOnly:true,sameSite:'Lax'}]);}
try{
 await cookie(identity);await page.goto(base+route);await expect(page.getByRole('heading',{name:kind==='ai'?'Actividad de IA':'Campañas',exact:true})).toBeVisible({timeout:90000});
 if(kind==='ai'){
  await expect(page.getByRole('heading',{name:'Cotización',exact:true})).toBeVisible({timeout:20000});checks.push('runtime_generated_proposal_read_from_postgres');
  await page.getByText('Ver detalle',{exact:false}).first().click();await expect(page.getByText('La reserva es un límite interno en tokens; no es una factura ni confirma un envío.',{exact:true})).toBeVisible();checks.push('progressive_detail_internal_budget_boundary');
  await page.getByRole('button',{name:'Intervención humana',exact:true}).click();await expect(page.getByText('No hay propuestas en este filtro.',{exact:true})).toBeVisible();
  await page.getByRole('button',{name:'Todas',exact:true}).click();await page.reload();await expect(page.getByRole('heading',{name:'Cotización',exact:true})).toBeVisible();checks.push('filter_and_reload_preserve_owner_projection');
 }else{
  const list=page.getByRole('region',{name:'Campañas disponibles'}),row=list.getByRole('button').filter({has:page.locator('strong')}).first();await expect(row).toBeVisible({timeout:20000});await row.click();await expect(page.getByText('Contenido y horarios',{exact:true})).toBeVisible({timeout:20000});checks.push('campaign_list_to_durable_detail');await page.reload();await expect(row).toBeVisible();
 }
 await page.screenshot({path:out+'/desktop.png',fullPage:true});await page.setViewportSize({width:390,height:844});await expect(page.getByRole('heading',{name:kind==='ai'?'Actividad de IA':'Campañas',exact:true})).toBeVisible();
 if(await page.evaluate(()=>document.documentElement.scrollWidth>window.innerWidth+1))throw Error('mobile page overflow');
 await page.screenshot({path:out+'/mobile.png',fullPage:true});checks.push('390_viewport_no_document_overflow');
 const api=kind==='ai'?'/api/enterprise/intelligence/activity':'/api/enterprise/franchise/campaigns?kind=campaigns';
 let response=await context.request.get(base+api);if(response.status()!==200)throw Error('authorized BFF '+response.status());checks.push('actual_bff_authorized_200');
 await cookie({...identity,permissions:[]});response=await context.request.get(base+api);if(response.status()!==403)throw Error('BFF missing permission not 403');checks.push('bff_missing_permission_403');
 await context.clearCookies();response=await context.request.get(base+api);if(response.status()!==401)throw Error('BFF guest not 401');checks.push('bff_guest_401');
 await writeFile(out+'/RESULT.json',JSON.stringify({result:'PASS',kind,checks,responses,method:'Chromium headless -> actual Next page/BFF -> Go owner -> owned PostgreSQL; no browser interception',visual_approval:'PENDING_USER_APPROVAL',physical_device:'NOT_RUN',production_authorized:false},null,2)+'\n');
}finally{await browser.close();}
