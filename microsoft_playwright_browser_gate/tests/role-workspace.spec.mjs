import {test,expect} from '@playwright/test';
import {createHash} from 'node:crypto';
import {createRequire} from 'node:module';
import {resolve} from 'node:path';
import {pathToFileURL} from 'node:url';
async function identity(context,permissions,expiry='300s',tamper=false){
 const base=process.env.ELITE_BASE_URL;expect(new URL(base).protocol).toBe('https:');expect(new URL(base).hostname).toBe('127.0.0.1');
 const require=createRequire(resolve(process.env.ELITE_WEB_ROOT,'package.json'));const {EncryptJWT}=await import(pathToFileURL(require.resolve('jose')).href);
 const jwt=await new EncryptJWT({subject:'private-workspace-subject',tenantId:'private-workspace-tenant',organizations:['private-org'],permissions,accessToken:'private-unused-token'})
  .setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime(expiry).encrypt(createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest());
 await context.clearCookies();await context.addCookies([{name:'__Host-elite_session',value:tamper?'invalid.'+jwt:jwt,url:base+'/',httpOnly:true,secure:true,sameSite:'Lax'}]);
}
const matrix=[
  {
    "name": "none",
    "permissions": [],
    "privateIds": []
  },
  {
    "name": "customer",
    "permissions": [
      "customer:self"
    ],
    "privateIds": [
      "customer"
    ]
  },
  {
    "name": "factory",
    "permissions": [
      "factory:read"
    ],
    "privateIds": [
      "factory"
    ]
  },
  {
    "name": "admin",
    "permissions": [
      "admin:read"
    ],
    "privateIds": [
      "admin",
      "franchise"
    ]
  },
  {
    "name": "lead",
    "permissions": [
      "lead:read"
    ],
    "privateIds": [
      "franchise"
    ]
  },
  {
    "name": "inventory",
    "permissions": [
      "inventory:allocate"
    ],
    "privateIds": [
      "franchise"
    ]
  },
  {
    "name": "payment",
    "permissions": [
      "payment:create"
    ],
    "privateIds": [
      "franchise"
    ]
  },
  {
    "name": "handover",
    "permissions": [
      "handover:manage"
    ],
    "privateIds": [
      "franchise"
    ]
  },
  {
    "name": "resource",
    "permissions": [
      "resource:manage"
    ],
    "privateIds": [
      "franchise"
    ]
  },
  {
    "name": "availability-read",
    "permissions": [
      "availability:read"
    ],
    "privateIds": [
      "franchise"
    ]
  },
  {
    "name": "availability-manage",
    "permissions": [
      "availability:manage"
    ],
    "privateIds": [
      "franchise"
    ]
  },
  {
    "name": "appointment",
    "permissions": [
      "appointment:manage"
    ],
    "privateIds": [
      "franchise"
    ]
  },
  {
    "name": "quote-only",
    "permissions": [
      "quote:write"
    ],
    "privateIds": []
  },
  {
    "name": "wrong-case",
    "permissions": [
      "Admin:Read",
      "owner",
      "lead:reader",
      "* "
    ],
    "privateIds": []
  },
  {
    "name": "wildcard",
    "permissions": [
      "*"
    ],
    "privateIds": [
      "admin",
      "customer",
      "factory",
      "franchise"
    ]
  }
];
const paths={dashboard:'/dashboard',help:'/help',locations:'/locations',public_catalog:'/models',customer:'/customer',factory:'/factory',admin:'/admin',franchise:'/franchise'};
const rolePaths={owner:['/dashboard','/admin','/franchise','/factory','/models','/help'],admin:['/dashboard','/admin','/franchise','/models','/help'],employee:['/dashboard','/franchise','/models','/help'],customer:['/dashboard','/customer','/models','/locations','/help']};
test('workspace navigation and role views preserve effective permissions',async({page,context})=>{
 test.skip(process.env.ELITE_WORKSPACE_E2E!=='enabled','Requires the enabled isolated workspace fixture');
 const writes=[];page.on('request',r=>{if(!['GET','HEAD'].includes(r.method()))writes.push(r.method()+' '+new URL(r.url()).pathname);});
 // Inspect the protected route response without following it to an identity provider.
 async function expectLogin(path,returnTo){
  const response=await context.request.get(path,{maxRedirects:0});
  expect(response.status()).toBe(307);
  const actual=new URL(response.headers().location,process.env.ELITE_BASE_URL);
  expect(actual.origin).toBe(new URL(process.env.ELITE_BASE_URL).origin);
  expect(actual.pathname).toBe('/api/auth/login');expect(actual.searchParams.get('return_to')).toBe(returnTo);
 }
 await expectLogin('/dashboard','/dashboard');
 for(const row of matrix){
  await identity(context,row.permissions);await page.goto('/dashboard');await expect(page.getByRole('heading',{name:'Tu espacio de trabajo'})).toBeVisible();
  const allowed=['dashboard','help','locations','public_catalog',...row.privateIds].map(id=>paths[id]).sort();
  const header=await page.locator('header nav a').evaluateAll(nodes=>nodes.map(n=>n.getAttribute('href')).sort());expect(header,row.name).toEqual(allowed);
  const panel=page.getByRole('region',{name:'Accesos de tu sesión'});expect(await panel.locator('a').evaluateAll(nodes=>nodes.map(n=>n.getAttribute('href')).sort()),row.name).toEqual([...allowed,'/guide'].sort());
  expect(await page.locator('main').innerText()).not.toMatch(/private-workspace|private-org|private-unused-token|alto monto|alto valor|regalías/);
  for(const role of Object.keys(rolePaths)){
   await page.goto('/guide/'+role);await expect(page.getByRole('heading',{name:/Recorrido:/})).toBeVisible();
   const expected=allowed.filter(path=>rolePaths[role].includes(path));const links=await page.getByRole('region',{name:'Accesos de tu sesión'}).locator('a').evaluateAll(nodes=>nodes.map(n=>n.getAttribute('href')).sort());expect(links,row.name+'/'+role).toEqual(expected);
  }
 }
 await identity(context,['customer:self']);await page.goto('/guide');await page.getByRole('link',{name:'Ver recorrido de Dueño'}).click();
 await expect(page.getByRole('region',{name:'Accesos de tu sesión'}).locator('a[href="/franchise"],a[href="/admin"],a[href="/factory"]')).toHaveCount(0);
 await page.getByRole('link',{name:'Abrir ayuda versionada'}).click();await expect(page.locator('main article')).toHaveCount(2);
 const denied=await page.request.get('/admin');expect(await denied.text()).toContain('Acceso denegado'); // Existing page denies before any backend read.
 for(const path of ['/guide/unknown','/guide/__proto__','/guide/constructor']){const r=await page.request.get(path);expect(r.status()).toBe(404);}
 for(const [expiry,tamper] of [[Math.floor(Date.now()/1000)-60,false],['300s',true]]){await identity(context,['*'],expiry,tamper);await expectLogin('/guide/owner','/guide');}
 await identity(context,['resource:manage']);await page.goto('/dashboard');await page.screenshot({path:test.info().outputPath('workspace-resource.png'),fullPage:true});
 expect(await page.evaluate(()=>document.documentElement.scrollWidth<=innerWidth)).toBe(true);expect(writes).toEqual([]);
});
test('workspace disabled configuration mounts no workspace route',async({page,context})=>{
 test.skip(process.env.ELITE_WORKSPACE_E2E!=='disabled','Requires the disabled isolated workspace fixture');
 await identity(context,['*']);for(const path of ['/dashboard','/guide','/guide/owner']){const response=await page.goto(path);expect(response.status()).toBe(404);await expect(page.locator('header nav a[href="/dashboard"]')).toHaveCount(0);}
});


test.describe('private read presentation across host locales',()=>{
 test.use({timezoneId:'Pacific/Honolulu',locale:'de-DE'});
test('admin reads preserve independent grants and durable scopes',async({page,context})=>{
 test.skip(process.env.ELITE_ADMIN_READ_E2E!=='1','Requires the actual isolated Go/PostgreSQL admin-read fixture');
 const identities=JSON.parse(process.env.ELITE_ADMIN_IDENTITIES);const writes=[];
 page.on('request',r=>{if(!['GET','HEAD'].includes(r.method()))writes.push(r.method()+' '+new URL(r.url()).pathname)});
 async function setIdentity(name,permissions){
  const require=createRequire(resolve(process.env.ELITE_WEB_ROOT,'package.json'));const {EncryptJWT}=await import(pathToFileURL(require.resolve('jose')).href);
  const payload={...identities[name]};if(permissions)payload.permissions=permissions;
  const jwt=await new EncryptJWT(payload).setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('300s').encrypt(createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest());
  await context.clearCookies();await context.addCookies([{name:'__Host-elite_session',value:jwt,url:process.env.ELITE_BASE_URL+'/',httpOnly:true,secure:true,sameSite:'Lax'}]);
 }
 for(const [name,leads] of [['admin',false],['admin-case',false],['admin-lead',true],['wildcard',true]]){
  await setIdentity(name);await page.goto('/dashboard');await page.getByRole('region',{name:'Accesos de tu sesión'}).locator('a[href="/admin"]').click();
  await expect(page.getByRole('heading',{name:'Operación',exact:true})).toBeVisible();
  await expect(page.getByRole('heading',{name:'order-a',exact:true})).toBeVisible();await expect(page.getByRole('heading',{name:'order-b',exact:true})).toBeVisible();await expect(page.getByRole('heading',{name:'case-a',exact:true})).toBeVisible();
  await expect(page.getByRole('heading',{name:'Oportunidades',exact:true})).toHaveCount(leads?1:0);await expect(page.getByRole('heading',{name:'lead-visible',exact:true})).toHaveCount(leads?1:0);
  expect(await page.locator('main').innerText()).not.toMatch(/order-c|lead-foreign-org/);
  const order=page.locator('article').filter({has:page.getByRole('heading',{name:'order-a',exact:true})});expect((await order.innerText()).replace(/\u00a0/g,' ')).toContain('USD 1,00');
 }
 // Changing the session on the same path must remove the now-ungranted section.
 await setIdentity('admin');await page.reload();await expect(page.getByRole('heading',{name:'Oportunidades',exact:true})).toHaveCount(0);await expect(page.getByRole('heading',{name:'order-a',exact:true})).toBeVisible();
 await setIdentity('other-org');await page.goto('/admin');await expect(page.getByRole('heading',{name:'order-c',exact:true})).toBeVisible();await expect(page.getByRole('heading',{name:'lead-foreign-org',exact:true})).toBeVisible();expect(await page.locator('main').innerText()).not.toMatch(/order-a|order-b|case-a|lead-visible/);
 await setIdentity('other-tenant');await page.goto('/admin');await expect(page.getByRole('heading',{name:'Operación',exact:true})).toBeVisible();expect(await page.locator('main').innerText()).not.toMatch(/order-a|order-b|order-c|case-a|lead-visible|lead-foreign-org/);
 await setIdentity('lead');await page.goto('/admin');await expect(page.getByRole('heading',{name:'Acceso denegado',exact:true})).toBeVisible();
 await setIdentity('factory');await page.goto('/dashboard');await page.getByRole('region',{name:'Accesos de tu sesión'}).locator('a[href="/factory"]').click();await expect(page.getByRole('heading',{name:'UNIT-QUERY',exact:true})).toBeVisible();
 await setIdentity('customer');await page.goto('/dashboard');await page.getByRole('region',{name:'Accesos de tu sesión'}).locator('a[href="/customer"]').click();
 await expect(page.getByRole('heading',{name:'Mi cuenta',exact:true})).toBeVisible();
 for(const id of ['order-a','quote-visible']){const card=page.locator('article').filter({has:page.getByRole('heading',{name:id,exact:true})});expect((await card.innerText()).replace(/\u00a0/g,' ')).toContain('USD 1,00');}
 expect(await page.locator('main').innerText()).not.toMatch(/order-b|order-c/);

 const dateFixtures = JSON.parse(process.env.ELITE_APPOINTMENT_TIME_EXPECTATIONS);
 const timeText = async () => {
  const times = page.locator('main time');await expect(times).toHaveCount(2);
  const observed=[];
  for(const row of dateFixtures){
   const entries=await times.evaluateAll(nodes=>nodes.map(n=>({instant:n.getAttribute('datetime'),label:n.textContent})));
   const matches=entries.filter(item=>Date.parse(item.instant)===Date.parse(row.instant));expect(matches).toHaveLength(1);
   const value=matches[0].label.replace(/[\u00a0\u202f]/g,' ');expect(value).toContain(row.day);expect(value).toContain(row.hour);expect(value).toContain(row.zone);observed.push({instant:matches[0].instant,label:value});
  }
  expect(await page.evaluate(()=>Intl.DateTimeFormat().resolvedOptions().timeZone)).toBe('Pacific/Honolulu');
  return observed;
 };
 await page.waitForLoadState("networkidle");
 const accountTimes=await timeText();
 const pageErrors=[], hydrationErrors=[];page.on('pageerror',error=>pageErrors.push(error.message));
 page.on('console',message=>{if(message.type()==='error' && /hydration|Minified React error #(418|419|421|422|423|425)/i.test(message.text()))hydrationErrors.push(message.text())});
 await page.goto('/customer/appointments');await expect(page.getByRole('heading',{name:'Mis turnos',exact:true})).toBeVisible();
 expect(await timeText()).toEqual(accountTimes);await expect(page.getByRole('button',{name:'Cancelar mi turno',exact:true})).toHaveCount(2);
 await page.waitForLoadState("networkidle");
 await page.reload();await page.waitForLoadState("networkidle");expect(await timeText()).toEqual(accountTimes);expect(hydrationErrors).toEqual([]);
 await test.info().attach('appointment-page-errors',{body:JSON.stringify({pageErrors,hydrationErrors},null,2),contentType:'application/json'});
 expect(pageErrors).toEqual([]);
 await page.screenshot({path:test.info().outputPath('customer-appointment-time.png'),fullPage:true});
 expect(await page.evaluate(()=>document.documentElement.scrollWidth<=innerWidth)).toBe(true);
 await page.goto('/customer/quotes');await expect(page.getByRole('heading',{name:'Mis cotizaciones',exact:true})).toBeVisible();expect((await page.locator('main').innerText()).replace(/\u00a0/g,' ')).toContain('USD 1,00');
 // Encrypted fixture claims cannot widen the real bearer. Recovery makes a new GET.
 for(const row of [
  {path:'/admin',permissions:['admin:read','lead:read'],valid:'admin',heading:'order-a'},
  {path:'/customer',permissions:['customer:self'],valid:'customer',heading:'order-a'},
  {path:'/factory',permissions:['factory:read'],valid:'factory',heading:'UNIT-QUERY'}
 ]){
  await setIdentity('forged-admin',row.permissions);await page.goto(row.path);
  await expect(page.getByRole('heading',{name:'No pudimos cargar esta información',exact:true})).toBeVisible();
  await expect(page.locator('main').getByRole('alert')).toContainText('La consulta no se completó');
  expect(await page.locator('main').innerText()).not.toMatch(/order-a|order-b|lead-visible|quote-visible|UNIT-QUERY|FORBIDDEN|Bearer|digest/);
  const reconsult=page.getByRole('link',{name:'Volver a consultar',exact:true});await expect(reconsult).toHaveAttribute('href',row.path);
  // An unchanged insufficient bearer must still fail; no implicit escalation or empty success.
  await reconsult.click();await expect(page.locator('main').getByRole('alert')).toBeVisible();
  await page.screenshot({path:test.info().outputPath('read-recovery-'+row.valid+'.png'),fullPage:true});
  await setIdentity(row.valid);
  const requestPromise=page.waitForRequest(r=>r.isNavigationRequest()&&new URL(r.url()).pathname===row.path);
  await page.getByRole('link',{name:'Volver a consultar',exact:true}).click();expect((await requestPromise).method()).toBe('GET');
  await expect(page.getByRole('heading',{name:row.heading,exact:true})).toBeVisible();await expect(page.locator('main').getByRole('alert')).toHaveCount(0);
 }

 // Traverse the real PostgreSQL pages, retaining independent list cursors.
 await setIdentity('customer');await page.goto('/customer');
 const orders=page.getByRole('region',{name:'Mis pedidos',exact:true});
 const cases=page.getByRole('region',{name:'Mis casos de servicio',exact:true});
 const ids=Array.from({length:26},(_,i)=>String(i+1).padStart(3,'0'));
 await expect(orders.locator('article')).toHaveCount(25);await expect(cases.locator('article')).toHaveCount(25);
 const orderIds=await orders.locator('article h3').allTextContents();const caseIds=await cases.locator('article h3').allTextContents();
 await orders.getByRole('link',{name:'Ver siguientes',exact:true}).click();await expect(orders.locator('article')).toHaveCount(2);await expect(cases.locator('article')).toHaveCount(25);
 orderIds.push(...await orders.locator('article h3').allTextContents());expect(orderIds).toEqual(['order-a',...ids.map(x=>'zz-order-'+x)]);
 await expect(orders.getByRole('link',{name:'Ver siguientes',exact:true})).toHaveCount(0);
 await cases.getByRole('link',{name:'Ver siguientes',exact:true}).click();await expect(cases.locator('article')).toHaveCount(2);await expect(orders.locator('article')).toHaveCount(2);
 caseIds.push(...await cases.locator('article h3').allTextContents());expect(caseIds).toEqual(['case-a',...ids.map(x=>'zz-case-'+x)]);
 expect(new URL(page.url()).searchParams.get('orders_after')).toBe('zz-order-024');expect(new URL(page.url()).searchParams.get('cases_after')).toBe('zz-case-024');
 await page.reload();await expect(orders.locator('article')).toHaveCount(2);await expect(cases.locator('article')).toHaveCount(2);
 await orders.getByRole('link',{name:'Volver al inicio',exact:true}).click();await expect(orders.locator('article')).toHaveCount(25);await expect(cases.locator('article')).toHaveCount(2);
 await cases.getByRole('link',{name:'Volver al inicio',exact:true}).click();await expect(cases.locator('article')).toHaveCount(25);expect(new URL(page.url()).search).toBe('');
 await setIdentity('factory');await page.goto('/factory');await expect(page.locator('main article')).toHaveCount(25);
 const unitIds=await page.locator('main article h2').allTextContents();await page.getByRole('link',{name:'Ver siguientes',exact:true}).click();await expect(page.locator('main article')).toHaveCount(2);
 unitIds.push(...await page.locator('main article h2').allTextContents());expect(unitIds).toEqual(['UNIT-QUERY',...ids.map(x=>'PAGE-UNIT-'+x)]);
 await expect(page.getByRole('link',{name:'Ver siguientes',exact:true})).toHaveCount(0);await page.reload();await expect(page.locator('main article')).toHaveCount(2);
 expect(await page.evaluate(()=>document.documentElement.scrollWidth<=innerWidth)).toBe(true);await page.screenshot({path:test.info().outputPath('factory-last-page.png'),fullPage:true});
 await page.getByRole('link',{name:'Volver al inicio',exact:true}).click();await expect(page.locator('main article')).toHaveCount(25);

 await setIdentity('admin-lead');await page.goto('/admin');
 const adminLists=[{name:'Pedidos recientes',ids:['order-a','order-b',...ids.map(x=>'zz-order-'+x)]},{name:'Servicio',ids:['case-a',...ids.map(x=>'zz-case-'+x)]},{name:'Oportunidades',ids:['lead-visible',...ids.map(x=>'zz-lead-'+x)]}];
 for(const list of adminLists){const section=page.getByRole('region',{name:list.name,exact:true});await expect(section.locator('article')).toHaveCount(25);const values=await section.locator('article h3').allTextContents();await section.getByRole('link',{name:'Ver siguientes',exact:true}).click();await expect(section.locator('article')).toHaveCount(list.ids.length-25);values.push(...await section.locator('article h3').allTextContents());expect(values).toEqual(list.ids);await expect(section.getByRole('link',{name:'Ver siguientes',exact:true})).toHaveCount(0);}
 await page.reload();for(const list of adminLists)await expect(page.getByRole('region',{name:list.name,exact:true}).locator('article')).toHaveCount(list.ids.length-25);
 await setIdentity('admin');await page.reload();await expect(page.getByRole('region',{name:'Oportunidades',exact:true})).toHaveCount(0);for(const href of await page.locator('main nav a').evaluateAll(items=>items.map(x=>x.getAttribute('href'))))expect(href).not.toContain('leads_after');
 await setIdentity('admin-lead');await page.reload();for(const list of adminLists){const section=page.getByRole('region',{name:list.name,exact:true});await section.getByRole('link',{name:'Volver al inicio',exact:true}).click();await expect(section.locator('article')).toHaveCount(25);}expect(new URL(page.url()).search).toBe('');
 await setIdentity('admin');await page.goto('/admin');
 expect(await page.evaluate(()=>document.documentElement.scrollWidth<=innerWidth)).toBe(true);expect(writes).toEqual([]);
 await page.screenshot({path:test.info().outputPath('admin-permission-read.png'),fullPage:true});
});

});

test.describe('customer cancellation result recovery',()=>{
 test.use({timezoneId:'Pacific/Honolulu',locale:'de-DE'});
 test('customer cancels with durable result recovery',async({page,context})=>{
  test.skip(process.env.ELITE_CUSTOMER_CANCEL_E2E!=='1','Requires the isolated actual cancellation fixture');
  const identities=JSON.parse(process.env.ELITE_ADMIN_IDENTITIES);
  const fixtures=JSON.parse(process.env.ELITE_CUSTOMER_CANCEL_FIXTURES)[test.info().project.name];
  const endpoint='/api/enterprise/franchise/commands';const pattern='**'+endpoint;
  async function identity(name){
   const require=createRequire(resolve(process.env.ELITE_WEB_ROOT,'package.json'));const {EncryptJWT}=await import(pathToFileURL(require.resolve('jose')).href);
   const jwt=await new EncryptJWT(identities[name]).setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('300s').encrypt(createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest());
   await context.clearCookies();await context.addCookies([{name:'__Host-elite_session',value:jwt,url:process.env.ELITE_BASE_URL+'/',httpOnly:true,secure:true,sameSite:'Lax'}]);
  }
  const command=row=>({action:'cancel-customer-appointment',organizationId:'org-a',appointmentId:row.id,version:1,reasonCode:'customer-request'});
  const send=body=>context.request.post(endpoint,{data:body,headers:{origin:process.env.ELITE_BASE_URL},maxRedirects:0});
  async function card(row){
   const dates=await page.locator('main time').evaluateAll(nodes=>nodes.map(n=>n.getAttribute('datetime')));
   const matches=dates.filter(x=>Date.parse(x)===Date.parse(row.instant));expect(matches).toHaveLength(1);
   return page.locator('article').filter({has:page.locator('time[datetime="'+matches[0]+'"]')});
  }
  async function open(){await page.goto('/customer/appointments');await expect(page.getByRole('heading',{name:'Mis turnos',exact:true})).toBeVisible();await page.waitForLoadState('networkidle');}
  async function checkState(row,state){const item=await card(row);await expect(item).toContainText(state);await expect(item.getByRole('button',{name:'Cancelar mi turno',exact:true})).toHaveCount(state==='cancelled'?0:1);}
  async function uncertain(){
   await expect(page.getByRole('status')).toContainText('No pudimos comprobar el resultado');
   await expect(page.getByRole('status')).not.toContainText('No se canceló');
   const buttons=page.getByRole('button',{name:'Cancelar mi turno',exact:true});for(const button of await buttons.all())await expect(button).toBeDisabled();
   await expect(page.getByRole('link',{name:'Volver a consultar mis turnos',exact:true})).toHaveAttribute('href','/customer/appointments');
  }
  async function recover(){await page.waitForLoadState('networkidle');await page.getByRole('link',{name:'Volver a consultar mis turnos',exact:true}).click();await expect(page.getByRole('status')).toHaveCount(0);await page.waitForLoadState('networkidle');}
  await identity('customer');
  expect((await send({...command(fixtures.normal),appointmentId:'time-foreign'})).status()).toBe(409);
  expect((await send({...command(fixtures.normal),organizationId:'org-b'})).status()).toBe(403);
  await identity('factory');expect((await send(command(fixtures.normal))).status()).toBe(403);
  await context.clearCookies();expect((await send(command(fixtures.normal))).status()).toBe(401);
  await identity('customer');await open();
  // Normal result, including a same-event double click: only one mutation may leave the client.
  let normalPosts=0;
  await page.route(pattern,async route=>{normalPosts++;await route.continue()});
  const normalReload=page.waitForEvent("domcontentloaded");
  const normalResponse=page.waitForResponse(r=>r.url().endsWith(endpoint)&&r.request().method()==='POST');
  await (await card(fixtures.normal)).getByRole('button',{name:'Cancelar mi turno'}).evaluate(button=>{button.click();button.click()});
  expect((await normalResponse).status()).toBe(200);await normalReload;expect(normalPosts).toBe(1);
  await page.waitForLoadState('networkidle');await expect.poll(async()=>await (await card(fixtures.normal)).innerText()).toContain('cancelled');await checkState(fixtures.normal,'cancelled');await page.unroute(pattern);
  // Forward the real write, then lose only its response after PostgreSQL committed.
  let lostReceipt, lostPosts=0;
  await page.route(pattern,async route=>{lostPosts++;const committed=await route.fetch({maxRetries:0,maxRedirects:0});expect(committed.status()).toBe(200);lostReceipt=await committed.json();await route.abort('failed')});
  await (await card(fixtures.lost)).getByRole('button',{name:'Cancelar mi turno'}).click();await uncertain();expect(lostPosts).toBe(1);expect(lostReceipt).toMatchObject({id:fixtures.lost.id,state:'cancelled',version:2});
  await page.unroute(pattern);await recover();await checkState(fixtures.lost,'cancelled');expect(lostPosts).toBe(1);
  // A syntactically valid but unbound success receipt cannot claim cancellation.
  await page.route(pattern,route=>route.fulfill({status:200,contentType:'application/json',body:JSON.stringify({id:'different-appointment',organization_id:'org-a',state:'cancelled',version:2})}));
  await (await card(fixtures.malformed)).getByRole('button',{name:'Cancelar mi turno'}).click();await uncertain();await page.unroute(pattern);await recover();await checkState(fixtures.malformed,'requested');
  const recoveredReload=page.waitForEvent("domcontentloaded");
  const recoveredResponse=page.waitForResponse(r=>r.url().endsWith(endpoint)&&r.request().method()==='POST');await (await card(fixtures.malformed)).getByRole('button',{name:'Cancelar mi turno'}).click();expect((await recoveredResponse).status()).toBe(200);await recoveredReload;await page.waitForLoadState('networkidle');await expect.poll(async()=>await (await card(fixtures.malformed)).innerText()).toContain('cancelled');
  // Two sessions/requests race with the same version. The stale page then recovers by GET.
  await page.goto('/customer');await page.getByRole('link',{name:'Gestionar mis turnos',exact:true}).click();await page.waitForLoadState('networkidle');await checkState(fixtures.stale,'requested');
  const concurrent=await Promise.all([send(command(fixtures.stale)),send(command(fixtures.stale))]);expect(concurrent.map(x=>x.status()).sort()).toEqual([200,409]);
  await (await card(fixtures.stale)).getByRole('button',{name:'Cancelar mi turno'}).click();await uncertain();await recover();await checkState(fixtures.stale,'cancelled');
  await page.reload();await page.waitForLoadState('networkidle');for(const row of Object.values(fixtures))await checkState(row,'cancelled');
  expect(await page.evaluate(()=>document.documentElement.scrollWidth<=innerWidth)).toBe(true);
  await page.screenshot({path:test.info().outputPath('customer-cancellation-recovered.png'),fullPage:true});
 });
});
