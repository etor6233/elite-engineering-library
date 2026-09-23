import { test, expect } from '@playwright/test';
import { createHash } from 'node:crypto';
import { createRequire } from 'node:module';
import { resolve } from 'node:path';
import { pathToFileURL } from 'node:url';

test('versioned operational help authorizes reads and recovers without operations', async ({ page, context }) => {
 test.skip(process.env.ELITE_HELP_E2E !== '1', 'Requires synthetic JWE and production Next on HTTPS loopback');
 const base=process.env.ELITE_BASE_URL; expect(new URL(base).protocol).toBe('https:'); expect(new URL(base).hostname).toBe('127.0.0.1');
 const require=createRequire(resolve(process.env.ELITE_WEB_ROOT,'package.json'));
 const {EncryptJWT}=await import(pathToFileURL(require.resolve('jose')).href);
 const key=createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest();
 async function identity(permissions, expiry='300s', tamper=false) {
  const jwt=await new EncryptJWT({subject:'synthetic-help-user',tenantId:'synthetic-help-tenant',organizations:['synthetic-org'],permissions,accessToken:'synthetic-unused-token'})
   .setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime(expiry).encrypt(key);
  await context.clearCookies(); await context.addCookies([{name:'__Host-elite_session',value:tamper?'invalid.'+jwt:jwt,url:base+'/',httpOnly:true,secure:true,sameSite:'Lax'}]);
 }
 const writes=[];page.on('request',r=>{if(!['GET','HEAD'].includes(r.method()))writes.push(r.method()+' '+new URL(r.url()).pathname);});
 await page.goto('/help');const panel=page.getByRole('region',{name:'Ayuda de los recorridos'});const status=panel.getByRole('status');
 await expect(status).toHaveText('Iniciá sesión para consultar la ayuda.'); await expect(panel.locator('article')).toHaveCount(0);
 await identity(['customer:self']);await panel.getByRole('button',{name:'Actualizar ayuda'}).click();
 await expect(panel.getByRole('heading',{name:'Aceptar una cotización'})).toBeVisible();await expect(panel.locator('article')).toHaveCount(2);
 await expect(panel).toContainText('Guía quote-acceptance-view/1.0.0');await expect(panel).toContainText('Aceptar una cotización no confirma el pago, el stock ni la entrega.');
 const input=panel.getByRole('searchbox',{name:'Buscar en las guías'});await input.fill('COTIZACIÓN');await expect(panel.locator('article')).toHaveCount(0);await panel.getByRole('button',{name:'Buscar',exact:true}).click();await expect(panel.locator('article')).toHaveCount(1);
 await panel.getByRole('link',{name:'Enlace a esta versión'}).click();await expect(page).toHaveURL(/article=quote-acceptance-view&version=1.0.0$/);await expect(panel.locator('article')).toHaveCount(1);
 await page.goto('/help?article=quote-acceptance-view&version=9.0.0');await expect(status).toContainText('Esta guía o versión no está disponible');await expect(panel.locator('article')).toHaveCount(0);
 await page.goto('/help?article=whatsapp-status-view&version=1.0.0');await expect(status).toContainText('Esta guía o versión no está disponible');
 await identity(['appointment:manage']);await panel.getByRole('button',{name:'Actualizar ayuda'}).click();await expect(panel.getByRole('heading',{name:'Consultar el estado de WhatsApp'})).toBeVisible();await expect(panel).toContainText('Guía whatsapp-status-view/1.0.0');await expect(panel).toContainText('no reenvíes ni borres el historial');
 let r=await page.request.get('/api/enterprise/help');expect(r.status()).toBe(200);expect(r.headers()['cache-control']).toBe('no-store');const body=await r.text();expect(body).not.toMatch(/synthetic-help|synthetic-org|synthetic-unused-token/);
 r=await page.request.post('/api/enterprise/help',{data:{}});expect(r.status()).toBe(405);
 r=await page.request.get('/api/enterprise/help?role=*');expect(r.status()).toBe(400);
 await page.route('**/api/enterprise/help*',route=>route.fulfill({status:503,contentType:'application/json',body:'{"code":"HELP_UNAVAILABLE"}'}),{times:1});await panel.getByRole('button',{name:'Actualizar ayuda'}).click();await expect(status).toContainText('No pudimos consultar la ayuda');await expect(panel.locator('article')).toHaveCount(0);
 await panel.getByRole('button',{name:'Actualizar ayuda'}).click();await expect(panel.locator('article')).toHaveCount(1);
 await page.route('**/api/enterprise/help*',route=>route.fulfill({status:200,contentType:'application/json',body:'{"schema":"bad","articles":[]}'}),{times:1});await panel.getByRole('button',{name:'Actualizar ayuda'}).click();await expect(status).toContainText('No pudimos consultar la ayuda');await expect(panel.locator('article')).toHaveCount(0);
 // Hold an actual authorized response, then supersede it. Late data must not repopulate the UI.
 let release;const held=new Promise(resolve=>{release=resolve});let arrived;const seen=new Promise(resolve=>{arrived=resolve});
 await page.route('**/api/enterprise/help*',async route=>{const response=await route.fetch();arrived();await held;await route.fulfill({response}).catch(()=>{});},{times:1});
 await panel.getByRole('button',{name:'Actualizar ayuda'}).click();await seen;
 await input.fill('nothing matches');await panel.getByRole('button',{name:'Buscar',exact:true}).click();await expect(status).toContainText('No encontramos guías');release();await expect(panel.locator('article')).toHaveCount(0);
 await identity(['factory:read']);await panel.getByRole('button',{name:'Actualizar ayuda'}).click();await expect(status).toContainText('Tu sesión no tiene guías disponibles');await expect(panel.locator('article')).toHaveCount(0);
 for(const [expiry,tamper] of [[Math.floor(Date.now()/1000)-60,false],['300s',true]]) {await identity(['customer:self'],expiry,tamper);await panel.getByRole('button',{name:'Actualizar ayuda'}).click();await expect(status).toContainText('Iniciá sesión');await expect(panel.locator('article')).toHaveCount(0);}
 await identity(['*']);await page.goto('/help');await expect(panel.locator('article')).toHaveCount(15);
 await input.fill('<script>alert(1)</script>');await panel.getByRole('button',{name:'Buscar',exact:true}).click();await expect(status).toContainText('No encontramos guías');expect(await page.locator('script').filter({hasText:'alert(1)'}).count()).toBe(0);
 expect(writes).toEqual([]);expect(await page.evaluate(()=>document.documentElement.scrollWidth<=innerWidth)).toBe(true);
 await page.screenshot({path:test.info().outputPath('help-recovery.png'),fullPage:true});
});

const guideVersions = {
  "return-operations-view": "1.0.0",
  "checklist-completion-view": "1.0.0",
  "checklist-publication-view": "1.0.0",
  "slot-create-view": "1.0.0",
  "resource-create-view": "1.0.0",
  "availability-create-view": "1.0.0",
  "quote-create-view": "1.0.0",
  "delivery-resolution-view": "1.0.0",
  "lead-command-view": "1.0.0",
  "availability-cancel-view": "1.0.0",
  "order-operations-view": "1.1.0",
  "operation-sections-view": "1.0.0",
  "handover-read-view": "1.0.0",
  "quote-acceptance-view": "1.0.0",
  "whatsapp-status-view": "1.0.0"
};
const roleMatrix = [
  {
    "role": "none",
    "permissions": [],
    "expected": []
  },
  {
    "role": "factory",
    "permissions": [
      "factory:read"
    ],
    "expected": []
  },
  {
    "role": "customer",
    "permissions": [
      "customer:self"
    ],
    "expected": [
      "handover-read-view",
      "quote-acceptance-view"
    ]
  },
  {
    "role": "lead",
    "permissions": [
      "lead:read"
    ],
    "expected": [
      "lead-command-view",
      "operation-sections-view"
    ]
  },
  {
    "role": "quote-without-lead-read",
    "permissions": [
      "quote:write"
    ],
    "expected": []
  },
  {
    "role": "quote",
    "permissions": [
      "lead:read",
      "quote:write"
    ],
    "expected": [
      "lead-command-view",
      "operation-sections-view",
      "quote-create-view"
    ]
  },
  {
    "role": "availability-read",
    "permissions": [
      "availability:read"
    ],
    "expected": [
      "availability-cancel-view",
      "operation-sections-view"
    ]
  },
  {
    "role": "availability-manage",
    "permissions": [
      "availability:manage"
    ],
    "expected": [
      "availability-create-view",
      "operation-sections-view"
    ]
  },
  {
    "role": "resource",
    "permissions": [
      "resource:manage"
    ],
    "expected": [
      "operation-sections-view",
      "resource-create-view"
    ]
  },
  {
    "role": "appointment",
    "permissions": [
      "appointment:manage"
    ],
    "expected": [
      "operation-sections-view",
      "slot-create-view",
      "whatsapp-status-view"
    ]
  },
  {
    "role": "inventory",
    "permissions": [
      "inventory:allocate"
    ],
    "expected": [
      "operation-sections-view",
      "order-operations-view"
    ]
  },
  {
    "role": "payment",
    "permissions": [
      "payment:create"
    ],
    "expected": [
      "operation-sections-view",
      "order-operations-view"
    ]
  },
  {
    "role": "handover",
    "permissions": [
      "handover:manage"
    ],
    "expected": [
      "checklist-completion-view",
      "checklist-publication-view",
      "delivery-resolution-view",
      "operation-sections-view",
      "order-operations-view",
      "return-operations-view"
    ]
  },
  {
    "role": "admin",
    "permissions": [
      "admin:read"
    ],
    "expected": [
      "operation-sections-view",
      "order-operations-view"
    ]
  },
  {
    "role": "wildcard",
    "permissions": [
      "*"
    ],
    "expected": [
      "availability-cancel-view",
      "availability-create-view",
      "checklist-completion-view",
      "checklist-publication-view",
      "delivery-resolution-view",
      "handover-read-view",
      "lead-command-view",
      "operation-sections-view",
      "order-operations-view",
      "quote-acceptance-view",
      "quote-create-view",
      "resource-create-view",
      "return-operations-view",
      "slot-create-view",
      "whatsapp-status-view"
    ]
  }
];
test('all existing versioned guides have exact role coverage and usable links', async ({page,context}) => {
 test.skip(process.env.ELITE_HELP_E2E !== '1', 'Requires isolated production Next with synthetic JWE');
 const base=process.env.ELITE_BASE_URL;expect(new URL(base).protocol).toBe('https:');expect(new URL(base).hostname).toBe('127.0.0.1');
 const require=createRequire(resolve(process.env.ELITE_WEB_ROOT,'package.json'));const {EncryptJWT}=await import(pathToFileURL(require.resolve('jose')).href);
 const key=createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest();
 const writes=[];page.on('request',r=>{if(!['GET','HEAD'].includes(r.method()))writes.push(r.method()+' '+new URL(r.url()).pathname);});
 for(const row of roleMatrix){
  const jwt=await new EncryptJWT({subject:'synthetic-'+row.role,tenantId:'synthetic-coverage',organizations:['synthetic-org'],permissions:row.permissions,accessToken:'unused-synthetic-token'})
    .setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('300s').encrypt(key);
  await context.clearCookies();await context.addCookies([{name:'__Host-elite_session',value:jwt,url:base+'/',httpOnly:true,secure:true,sameSite:'Lax'}]);
  const list=await page.request.get('/api/enterprise/help');expect(list.status(),row.role).toBe(row.expected.length?200:403);
  if(row.expected.length){const data=await list.json();expect(data.articles.map(g=>g.id).sort(),row.role).toEqual(row.expected);expect(JSON.stringify(data)).not.toMatch(/synthetic-coverage|synthetic-org|unused-synthetic-token/);}
  await page.goto('/help');const panel=page.getByRole('region',{name:'Ayuda de los recorridos'});
  if(row.expected.length){await expect(panel.getByRole('status')).toContainText('Ayuda actualizada');await expect(panel.locator('article')).toHaveCount(row.expected.length);}
  else {await expect(panel.getByRole('status')).toContainText('Tu sesión no tiene guías');await expect(panel.locator('article')).toHaveCount(0);}
  for(const [id,version] of Object.entries(guideVersions)){
   const endpoint='/api/enterprise/help?'+new URLSearchParams({article:id,version});const r=await page.request.get(endpoint);
   const included=row.expected.includes(id);expect(r.status(),row.role+'/'+id).toBe(included?200:row.expected.length?404:403);expect(r.headers()['cache-control']).toBe('no-store');
   if(included){const data=await r.json();expect(data.articles).toHaveLength(1);expect(data.articles[0].id).toBe(id);expect(data.articles[0].version).toBe(version);const a=panel.locator(`a[href="/help?article=${id}&version=${version}"]`);await expect(a).toBeVisible();}
  }
 }
 // The last profile is wildcard. Follow every version link through the actual browser UI.
 for(const [id,version] of Object.entries(guideVersions)){
  await page.goto('/help');const panel=page.getByRole('region',{name:'Ayuda de los recorridos'});
  await panel.locator(`a[href="/help?article=${id}&version=${version}"]`).click();await expect(panel.locator('article')).toHaveCount(1);
  await expect(panel.locator('article')).toContainText('Guía '+id+'/'+version);
 }
 expect(writes).toEqual([]);expect(await page.evaluate(()=>document.documentElement.scrollWidth<=innerWidth)).toBe(true);
 await page.screenshot({path:test.info().outputPath('all-guides-covered.png'),fullPage:true});
});
