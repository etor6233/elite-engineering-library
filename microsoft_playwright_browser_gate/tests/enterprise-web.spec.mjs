import { test, expect } from '@playwright/test';
import { createHash } from 'node:crypto';
import { createRequire } from 'node:module';
import { resolve } from 'node:path';
import { pathToFileURL } from 'node:url';

test('operator agenda confirms without typing identifiers or versions', async ({ page, context }) => {
  test.skip(process.env.ELITE_OPERATOR_AGENDA_E2E !== '1', 'Requires isolated Go/PG and synthetic signed identity');
  const require = createRequire(resolve(process.env.ELITE_WEB_ROOT, 'package.json'));
  const { EncryptJWT } = await import(pathToFileURL(require.resolve('jose')).href);
  const key = createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest();
  const session = await new EncryptJWT({ subject: 'operator', tenantId: process.env.ELITE_CONFIRMATION_TENANT, permissions: ['appointment:manage'], organizations: ['store'], accessToken: process.env.ELITE_CONFIRMATION_OPERATOR })
    .setProtectedHeader({ alg: 'dir', enc: 'A256GCM', typ: 'JWT' }).setIssuedAt().setExpirationTime('300s').encrypt(key);
  const base=process.env.ELITE_BASE_URL;
  expect(new URL(base).protocol).toBe('https:');
  expect(new URL(base).hostname).toBe('127.0.0.1');
  // Runtime image optimization is deliberately absent from this reference.
  expect((await page.request.get(base+'/_next/image')).status()).toBe(404);
  await context.addCookies([{ name: '__Host-elite_session', value: session, url: base+'/', httpOnly: true, secure: true, sameSite: 'Lax' }]);
  expect((await context.cookies(base)).map(cookie=>cookie.name)).toContain('__Host-elite_session');
  await page.goto(base+'/franchise?day=' + process.env.ELITE_CONFIRMATION_DAY);
  const agenda = page.getByRole('region', { name: 'Agenda de turnos' });
  await expect(agenda).toBeVisible();
  await expect(agenda.getByText('Estado: Solicitado', { exact: true })).toBeVisible();
  await expect(agenda.locator('input[name="appointmentId"],input[name="version"],input[name="resourceId"]')).toHaveCount(0);
  await agenda.getByRole('combobox', { name: 'Recurso disponible para evaluar' }).selectOption({ label: 'Synthetic technician' });
  await agenda.getByRole('button', { name: 'Asignar recurso', exact: true }).click();
  await expect(agenda.getByText('Recurso: Synthetic technician', { exact: true })).toBeVisible();
  const loseResponse=test.info().project.name!=='chromium-desktop';
  if(loseResponse) {
    await page.route('**/api/enterprise/franchise/commands',async route=>{
      const response=await route.fetch();
      expect(response.status()).toBe(200);
      expect((await response.json()).state).toBe('confirmed');
      await route.abort('failed');
    },{times:1});
  }
  await agenda.getByRole('button', { name: 'Confirmar turno', exact: true }).click();
  if(loseResponse) {
    await expect(agenda.getByRole('status', { name: 'Resultado de la agenda', exact: true })).toContainText('No pudimos comprobar el resultado');
    await expect(agenda.getByRole('button',{name:'Confirmar turno',exact:true})).toBeDisabled();
    await agenda.getByRole('button',{name:'Actualizar agenda',exact:true}).click();
  }
  await expect(agenda.getByText('Estado: Confirmado', { exact: true })).toBeVisible();
  await expect(agenda.getByRole('button', { name: 'Confirmar turno', exact: true })).toHaveCount(0);
  await page.reload();
  await expect(agenda.getByText('Estado: Confirmado', { exact: true })).toBeVisible();
  expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth)).toBe(true);
  await page.screenshot({ path: test.info().outputPath('confirmed-agenda.png'), fullPage: true });
});

// Independent expected translations for the connected public journey; do not import the implementation catalog.
const publicLocale = process.env.ELITE_PUBLIC_LOCALE || 'es-AR';
if (!['es-AR', 'en-US', 'fr-FR'].includes(publicLocale)) throw new Error('Unsupported public test locale');
const publicEnglish = publicLocale === 'en-US';
const publicLanguage = publicEnglish ? 'en-US' : publicLocale === 'fr-FR' ? 'es' : 'es-AR';
const publicText = publicEnglish ? {
  name: 'Name', submit: 'Request information', received: 'Request received', next: 'Book an appointment',
  slot: 'Available appointment', request: 'Request appointment', requested: 'Appointment requested',
  modelTitle: 'Electric models', modelCount: '1 model available', place: '1 place', kind: 'Consultation',
  home: 'Configurable electric mobility', action: 'View models'
} : {
  name: 'Nombre', submit: 'Solicitar información', received: 'Solicitud recibida', next: 'Reservar turno',
  slot: 'Turno disponible', request: 'Solicitar turno', requested: 'Turno solicitado',
  modelTitle: 'Modelos eléctricos', modelCount: '1 modelo disponible', place: '1 lugar', kind: 'Consulta',
  home: 'Movilidad eléctrica configurable', action: 'Ver modelos'
};

async function capturePublicLead({ page, browserName }, testInfo) {
  test.skip(process.env.ELITE_PUBLIC_LEAD_E2E !== '1', 'Requires the Go disposable-database harness');
  if (browserName === 'webkit') {
    await page.route(/^https:\/\/127\.0\.0\.1:4173\//, async route => {
      const response = await route.fetch({ url: route.request().url().replace(/^https:/, 'http:') });
      await route.fulfill({ response });
    });
  }
  await page.goto('/');
  await expect(page.getByText(publicText.home, { exact: true })).toBeVisible();
  await expect(page.locator('html')).toHaveAttribute('lang', publicLanguage);
  await page.getByRole('link', { name: publicText.action, exact: true }).click();
  await expect(page.getByRole('heading', { name: publicText.modelTitle, exact: true })).toBeVisible();
  await expect(page.getByText(publicText.modelCount, { exact: true })).toBeVisible();
  await expect(page.locator('form')).toHaveAttribute('lang', publicLanguage);
  await expect(page.getByRole('heading', { name: 'Browser Fixture Model', exact: true })).toBeVisible();
  await page.getByRole('textbox', { name: publicText.name, exact: true }).fill('Browser ' + testInfo.project.name);
  await page.getByRole('textbox', { name: 'Email', exact: true }).fill(testInfo.project.name + '@example.invalid');
  await page.getByRole('checkbox').check();
  const submitted = page.waitForResponse(response => response.url().endsWith('/api/enterprise/leads') && response.request().method() === 'POST');
  await page.getByRole('button', { name: publicText.submit, exact: true }).click();
  const response = await submitted;
  expect(response.status()).toBe(202);
  const receipt = await response.json();
  expect(receipt.lead_id).toMatch(/^[0-9a-f-]{36}$/);
  await expect(page.getByRole('status')).toContainText(publicText.received);
  await expect(page.getByRole('link', { name: publicText.next, exact: true })).toHaveAttribute('href', `/locations?lead_id=${receipt.lead_id}&model_id=browser-model`);
  const body = response.request().postDataJSON();
  const key = response.request().headers()['idempotency-key'];
  const send = (data, idempotencyKey = key) => page.request.post('/api/enterprise/leads', {
    data, headers: { 'idempotency-key': idempotencyKey }, maxRedirects: 0,
  });
  const replay = await send(body);
  expect(replay.status()).toBe(200);
  expect(replay.headers()['idempotency-replayed']).toBe('true');
  expect(await replay.json()).toEqual(receipt);
  const conflict = await send({ ...body, name: 'Changed request' });
  expect(conflict.status()).toBe(409);
  expect((await conflict.json()).code).toBe('IDEMPOTENCY_CONFLICT');
  const noConsent = await send({ ...body, consentGranted: false }, `${key}-no-consent`);
  expect(noConsent.status()).toBe(400);
  const missingKey = await page.request.post('/api/enterprise/leads', { data: body, maxRedirects: 0 });
  expect(missingKey.status()).toBe(400);
  return receipt;
}

test('public lead crosses browser BFF Go and durable storage', async ({ page, browserName }, testInfo) => {
  await capturePublicLead({ page, browserName }, testInfo);
});

test('public appointment follows captured lead without overbooking', async ({ page, browserName }, testInfo) => {
  const lead = await capturePublicLead({ page, browserName }, testInfo);
  await page.getByRole('link', { name: publicText.next, exact: true }).click();
  await expect(page.getByRole('heading', { name: 'Browser Store', exact: true })).toBeVisible();
  const slotId = 'browser-slot-' + testInfo.project.name;
  await page.getByRole('combobox', { name: publicText.slot, exact: true }).selectOption(slotId);
  const selectedLabel = await page.getByRole('combobox', { name: publicText.slot, exact: true }).locator('option:checked').textContent();
  expect(selectedLabel).toContain(publicText.kind);
  expect(selectedLabel).toContain(publicText.place);
  await expect(page.locator('form')).toHaveAttribute('lang', publicLanguage);
  let committedReceipt;
  let originalKey;
  await page.route('**/api/enterprise/appointments', async route => {
    const localURL = route.request().url().replace(/^https:\/\/127\.0\.0\.1:4173\//, 'http://127.0.0.1:4173/');
    const committed = await route.fetch({ url: localURL, maxRedirects: 0, maxRetries: 0 });
    expect(committed.status()).toBe(202);
    committedReceipt = await committed.json();
    originalKey = route.request().headers()['idempotency-key'];
    // Lose only the browser response AFTER the real Go/PG commit, not the write.
    await route.abort('failed');
  }, { times: 1 });
  await page.getByRole('button', { name: publicText.request, exact: true }).click();
  await expect(page.getByRole('status')).not.toHaveText('');
  await expect(page.getByRole('status')).not.toContainText(publicText.requested);
  await expect(page.getByRole('button', { name: publicText.request, exact: true })).toBeEnabled();
  const submitted = page.waitForResponse(response => response.url().endsWith('/api/enterprise/appointments') && response.request().method() === 'POST');
  await page.getByRole('button', { name: publicText.request, exact: true }).click();
  const response = await submitted;
  expect(response.status()).toBe(200);
  const appointment = await response.json();
  expect(appointment).toEqual(committedReceipt);
  expect(response.request().headers()['idempotency-key']).toBe(originalKey);
  expect(appointment.id).toMatch(/^[0-9a-f-]{36}$/);
  expect(appointment.lead_id).toBe(lead.lead_id);
  expect(appointment.state).toBe('requested');
  expect(appointment.slot_id).toBe(slotId);
  // The actual submitted instant comes from the unchanged backend slot, not a formatted UI string.
  const posted = response.request().postDataJSON();
  expect(posted.kind).toBe('consultation');
  expect(posted.startsAt).toMatch(/^\d{4}-\d{2}-\d{2}T/);
  const zone = publicLocale === 'es-AR' ? 'America/Argentina/Buenos_Aires' : 'UTC';
  const expectedTime = new Intl.DateTimeFormat(publicLanguage, { dateStyle: 'medium', timeStyle: 'short', timeZone: zone }).format(new Date(posted.startsAt));
  // Intl data may differ in non-breaking space characters across admitted engines.
  // Preserve all date/time digits, punctuation and language; normalize only the two space variants.
  const comparableSpaces = value => value.replace(/[\u00a0\u202f]/g, ' ');
  expect(comparableSpaces(selectedLabel)).toContain(comparableSpaces(expectedTime));
  expect(posted.startsAt).toBe(appointment.starts_at);
  await expect(page.getByRole('status')).toContainText(publicText.requested);
  await expect(page.getByRole('status')).toContainText(appointment.id);
  await expect(page.getByRole('button', { name: publicText.request, exact: true })).toBeDisabled();
  const body = response.request().postDataJSON();
  const key = response.request().headers()['idempotency-key'];
  const send = (data, requestKey = key) => page.request.post('/api/enterprise/appointments', { data, headers: { 'idempotency-key': requestKey }, maxRedirects: 0 });
  const replay = await send(body);
  expect(replay.status()).toBe(200);
  expect(replay.headers()['idempotency-replayed']).toBe('true');
  expect(await replay.json()).toEqual(appointment);
  const divergent = await send({ ...body, kind: 'service' });
  expect(divergent.status()).toBe(409);
  expect((await divergent.json()).code).toBe('APPOINTMENT_CONFLICT');
  const full = await send(body, key + '-second');
  expect(full.status()).toBe(409);
  const noKey = await page.request.post('/api/enterprise/appointments', { data: body });
  expect(noKey.status()).toBe(400);
  const injected = await send({ ...body, state: 'confirmed' }, key + '-injected');
  expect(injected.status()).toBe(400);
  await page.reload();
  await expect(page.getByRole('combobox', { name: publicText.slot, exact: true }).locator(`option[value="${slotId}"]`)).toHaveCount(0);
});

test('public home is semantic, hardened and responsive', async ({ page, browserName }) => {
  // WebKit applies the production upgrade-insecure-requests directive even
  // to loopback subresources. Keep the header intact and proxy only those
  // exact local upgraded requests back to the same loopback HTTP server.
  if (browserName === 'webkit') {
    await page.route(/^https:\/\/127\.0\.0\.1:4173\//, async route => {
      const localURL = route.request().url().replace(/^https:/, 'http:');
      const response = await route.fetch({ url: localURL });
      await route.fulfill({ response });
    });
  }
  const errors = [];
  page.on('pageerror', error => errors.push(`pageerror:${error.message}`));
  page.on('console', message => { if (message.type() === 'error') errors.push(`console:${message.text()}`); });
  const response = await page.goto('/', { waitUntil: 'networkidle' });
  expect(response?.status()).toBe(200);
  await expect(page.locator('main')).toBeVisible();
  await expect(page.getByRole('heading', { level: 1 })).toHaveCount(1);
  await expect(page.getByRole('heading', { level: 1 })).not.toHaveText('');
  await expect(page.getByRole('link', { name: 'Ver modelos' })).toHaveAttribute('href', '/models');
  await expect(page.locator('html')).toHaveAttribute('lang', /^[a-z]{2}(?:-[A-Z]{2})?$/);
  expect(response?.headers()['x-content-type-options']).toBe('nosniff');
  expect(response?.headers()['x-frame-options']).toBe('DENY');
  expect(response?.headers()['referrer-policy']).toBe('strict-origin-when-cross-origin');
  expect(response?.headers()['content-security-policy']).toContain("frame-ancestors 'none'");
  const horizontalOverflow = await page.evaluate(() => document.documentElement.scrollWidth > document.documentElement.clientWidth);
  expect(horizontalOverflow).toBe(false);
  expect(errors).toEqual([]);
});

test('operator stock continues existing orders through authorized HTTP and durable recovery', async ({page,context})=>{
  test.skip(!process.env.ELITE_QUOTE_PHASE?.startsWith('operation'),'requires disposable connected operation fixture');
  test.setTimeout(60000);
  const {createHash}=await import('node:crypto');
  const {createRequire}=await import('node:module');
  const {resolve}=await import('node:path');
  const {pathToFileURL}=await import('node:url');
  const require=createRequire(resolve(process.env.ELITE_WEB_ROOT,'package.json'));
  const {EncryptJWT}=await import(pathToFileURL(require.resolve('jose')).href);
  const base=process.env.ELITE_BASE_URL;
  const setSession=async(kind='OPERATOR',permissions=['inventory:allocate'],organizations=['store'])=>{
    await context.clearCookies();
    const value=await new EncryptJWT({subject:kind.toLowerCase(),tenantId:process.env.ELITE_QUOTE_TENANT,permissions,organizations,accessToken:process.env['ELITE_QUOTE_'+kind]})
      .setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('300s').encrypt(createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest());
    await context.addCookies([{name:'__Host-elite_session',value,url:base+'/',httpOnly:true,secure:true,sameSite:'Lax'}]);
  };


  if(process.env.ELITE_QUOTE_PHASE==='operation-create-resource') {
    await setSession('RESOURCE',['resource:manage']);await page.goto(base+'/franchise');
    const section=()=>page.locator('section').filter({has:page.getByRole('heading',{name:'Recursos y turnos',exact:true})});
    const status=()=>section().getByRole('status',{name:'Resultado de recurso'});
    let posts=0;page.on('request',r=>{if(r.method()==='POST'&&r.url().endsWith('/api/enterprise/franchise/commands'))posts++});
    const keys=[];
    for(let index=0;index<3;index++){
      if(index>0){await section().getByRole('button',{name:'Preparar otro recurso',exact:true}).click();await expect(section().getByRole('button',{name:'Crear recurso',exact:true})).toBeEnabled();expect(posts).toBe(index);}
      const kind=['service-bay','employee','vehicle'][index];
      const fill=async()=>{await section().getByLabel('Nombre',{exact:true}).fill('Synthetic resource '+index);await section().getByRole('combobox',{name:'Tipo',exact:true}).selectOption(kind);if(index===1)await section().getByLabel('Identidad del proveedor de acceso',{exact:true}).fill('synthetic-employee');await section().getByRole('checkbox',{name:'service',exact:true}).check();await section().getByRole('checkbox',{name:'delivery',exact:true}).check();};
      await fill();
      if(index===2){await page.evaluate(()=>{window.originalSet=Storage.prototype.setItem;Storage.prototype.setItem=function(k,v){if(k.startsWith('elite-resource-create:'))throw new Error('synthetic storage');return window.originalSet.call(this,k,v)}});await section().getByRole('button',{name:'Crear recurso',exact:true}).click();await expect(status()).toContainText('No se envió');expect(posts).toBe(2);await page.reload();await fill();}
      let key,command,id;
      await page.route('**/api/enterprise/franchise/commands',async route=>{
        command=route.request().postDataJSON();key=route.request().headers()['idempotency-key'];expect(key).toMatch(/^resource-[a-f0-9-]{36}$/);expect(keys).not.toContain(key);keys.push(key);
        const responses=index===2?await Promise.all([route.fetch(),route.fetch()]):[await route.fetch()];
        for(const response of responses)expect(response.status()).toBe(201);
        const body=await responses[0].json();id=body.id;expect(body).toMatchObject({organization_id:'store',kind,status:'active',version:1});
        if(index===2)expect((await responses[1].json()).id).toBe(id);
        if(index===0)await route.abort('failed');else if(index===1)await route.fulfill({status:201,contentType:'application/json',body:'{}'});else await route.fulfill({response:responses[0]});
      },{times:1});
      await section().getByRole('button',{name:'Crear recurso',exact:true}).click();await expect(status()).toContainText(index<2?'No pudimos comprobar':'Recurso registrado');expect(posts).toBe(index+1);
      await expect(section().getByRole('button',{name:'Crear recurso',exact:true})).toBeDisabled();await expect(section().getByRole('button',{name:'Preparar otro recurso',exact:true})).toBeDisabled();
      expect(await page.evaluate(()=>Object.entries(sessionStorage).filter(([k])=>k.startsWith('elite-resource-create:')).map(([,v])=>JSON.parse(v)))).toEqual([{requestKey:key}]);
      await page.reload();
      if(index===0){await page.route('**/api/enterprise/franchise/commands?*',r=>r.fulfill({status:503,contentType:'application/json',body:'{"private":"synthetic-detail"}'}),{times:1});await section().getByRole('button',{name:'Consultar recurso creado',exact:true}).click();await expect(status()).toContainText('No pudimos recuperar');await expect(page.getByText('synthetic-detail')).toHaveCount(0);await expect(section().getByRole('button',{name:'Preparar otro recurso',exact:true})).toBeDisabled();}
      await section().getByRole('button',{name:'Consultar recurso creado',exact:true}).click();await expect(status()).toContainText('Consulta recuperada');await expect(section()).toContainText(id);expect(posts).toBe(index+1);
      const api=process.env.ENTERPRISE_API_BASE_URL,headers={authorization:'Bearer '+process.env.ELITE_QUOTE_RESOURCE,'idempotency-key':key};
      const data={organization_id:'store',principal_subject:command.principalSubject??'',display_name:command.displayName,kind:command.kind,skills:command.skills};
      const replay=await page.request.post(api+'/v1/franchise/resources',{headers,data});expect(replay.status()).toBe(200);expect((await replay.json()).id).toBe(id);
      expect((await page.request.post(api+'/v1/franchise/resources',{headers,data:{...data,display_name:'Divergent name'}})).status()).toBe(409);
      expect((await page.request.post(api+'/v1/franchise/resources',{headers,data:{...data,actor_subject:'forged'}})).status()).toBe(400);
      expect((await page.request.post(api+'/v1/franchise/resources',{headers:{authorization:headers.authorization},data})).status()).toBe(400);
      const lookup=api+'/v1/franchise/resources/result?'+new URLSearchParams({organization_id:'store',request_key:key});
      expect((await page.request.get(lookup,{headers:{authorization:'Bearer '+process.env.ELITE_QUOTE_READER}})).status()).toBe(403);
      expect((await page.request.get(lookup.replace('organization_id=store','organization_id=other'),{headers})).status()).toBe(403);
      expect((await page.request.get(lookup.replace(key,'unknown-resource-key'),{headers})).status()).toBe(409);
      expect(await page.evaluate(()=>document.documentElement.scrollWidth<=document.documentElement.clientWidth)).toBe(true);
    }
    await setSession('READER',['admin:read']);await page.goto(base+'/franchise');await expect(page.getByRole('button',{name:'Crear recurso',exact:true})).toHaveCount(0);
    return;
  }

  if(process.env.ELITE_QUOTE_PHASE==='operation-returns'&&process.env.ELITE_RETURN_MULTITAB_E2E==='1') {
    await setSession('HANDOVER',['handover:manage']);await exerciseReturnTabs(context,base,setSession);return;
  }

  if(process.env.ELITE_QUOTE_PHASE==='operation-returns') {
    await setSession('HANDOVER',['handover:manage']);await page.goto(base+'/franchise');
    const panel=()=>page.getByRole('region',{name:'Operaciones de devoluciones',exact:true}),status=()=>panel().getByRole('status',{name:'Resultado de devolución'}),consult=()=>panel().getByRole('button',{name:'Consultar operación de devolución',exact:true}),next=()=>panel().getByRole('button',{name:'Continuar operando',exact:true});
    const api=process.env.ENTERPRISE_API_BASE_URL,headers={authorization:'Bearer '+process.env.ELITE_QUOTE_HANDOVER},serial=process.env.ELITE_DELIVERY_SERIAL;
    let posts=0;page.on('request',r=>{if(r.method()==='POST'&&r.url().endsWith('/api/enterprise/franchise/commands'))posts++});
    for(let index=0;index<3;index++){
      const id=['ui-return-loss','ui-exchange-invalid','ui-return-race'][index],card=()=>panel().getByRole('article',{name:'Caso '+id,exact:true});let receiptId='';
      for(const action of ['receive-return','decide-return']){
        const before=posts,label=action==='receive-return'?'Registrar recepción':'Registrar decisión';
        const fill=async()=>{if(action==='receive-return'){await card().getByLabel('Serie recibida',{exact:true}).fill(serial);await card().getByLabel('Inspección/observaciones',{exact:true}).fill('  Synthetic receipt '+index+'  ')}else await card().getByLabel('Fundamento',{exact:true}).fill('  Synthetic decision '+index+'  ')};
        await fill();
        if(index===2){await page.evaluate(()=>{window.originalSet=Storage.prototype.setItem;Storage.prototype.setItem=function(k,v){if(k.startsWith('elite-return-operation:'))throw new Error('synthetic storage');return window.originalSet.call(this,k,v)}});await card().getByRole('button',{name:label,exact:true}).click();await expect(status()).toContainText('No se envió');expect(posts).toBe(before);await page.reload();await fill();}
        let command;
        await page.route('**/api/enterprise/franchise/commands',async route=>{
          command=route.request().postDataJSON();expect(command.action).toBe(action);
          const responses=index===2?await Promise.all([route.fetch(),route.fetch()]):[await route.fetch()];expect(responses.map(r=>r.status()).sort()).toEqual(index===2?[200,409]:[200]);
          const success=responses.find(r=>r.status()===200),body=await success.json();expect(body.id).toBeTruthy();if(action==='receive-return'){receiptId=body.id;expect(body.received_by_subject).toBe('handover')}else{expect(body.receipt_id).toBe(receiptId);expect(body.effect_requests).toHaveLength(index===1?3:4)}
          if(index===0)await route.abort('failed');else if(index===1)await route.fulfill({status:200,contentType:'application/json',body:'{}'});else await route.fulfill({response:success});
        },{times:1});
        await card().getByRole('button',{name:label,exact:true}).click();await expect(status()).toContainText(index<2?'No pudimos comprobar':'Respuesta recibida');expect(posts).toBe(before+1);await expect(next()).toBeDisabled();await expect(card().getByRole('button',{name:label,exact:true})).toBeDisabled();
        const markers=await page.evaluate(()=>Object.entries(sessionStorage).filter(([k])=>k.startsWith('elite-return-operation:')).map(([,v])=>JSON.parse(v)));expect(markers).toHaveLength(1);expect(Object.keys(markers[0]).sort()).toEqual(['action','authorizationId','digest','receiptId']);expect(markers[0]).toMatchObject({action,authorizationId:id,receiptId:action==='receive-return'?'':receiptId});expect(markers[0].digest).toMatch(/^[a-f0-9]{64}$/);expect(JSON.stringify(markers)).not.toContain(serial);expect(JSON.stringify(markers)).not.toContain(command.notes.trim());
        await page.reload();
        if(index===0&&action==='receive-return'){await expect(panel()).toContainText('Lista de devoluciones no disponible');await expect(card()).toHaveCount(0);await expect(page.getByText('PRIVATE-SYNTHETIC-RETURN-DETAIL')).toHaveCount(0)}
        if(index===0){
          await page.route('**/api/enterprise/franchise/commands?*',r=>r.fulfill({status:503,contentType:'application/json',body:'{"private":"synthetic-detail"}'}),{times:1});await consult().click();await expect(status()).toContainText('Conservamos el bloqueo');await expect(page.getByText('synthetic-detail')).toHaveCount(0);
          await page.route('**/api/enterprise/franchise/commands?*',async r=>{const response=await r.fetch(),body=await response.json();body.returnCase[action==='receive-return'?'receipt':'disposition'].notes='Different';await r.fulfill({status:200,contentType:'application/json',body:JSON.stringify(body)})},{times:1});await consult().click();await expect(status()).toContainText('Conservamos el bloqueo');await expect(next()).toBeDisabled();
        }
        await consult().click();await expect(status()).toContainText('Caso recuperado');await expect(panel().getByRole('region',{name:'Evidencia recuperada de devolución'})).toContainText(action==='receive-return'?'recibido por handover':'registrada por handover');expect(posts).toBe(before+1);
        const lookup=api+'/v1/franchise/returns/result?'+new URLSearchParams({organization_id:'store',authorization_id:id});const response=await page.request.get(lookup,{headers});expect(response.status()).toBe(200);expect(response.headers()['cache-control']).toContain('no-store');const value=await response.json();expect(value.receipt.id).toBe(receiptId);if(action==='decide-return'){expect(value.disposition.effect_requests).toHaveLength(index===1?3:4);expect(value.disposition.effect_requests.every(e=>e.state==='requested')).toBe(true);if(index===1)expect(value.disposition.effect_requests.some(e=>e.owner_context==='fiscal')).toBe(false)}
        expect((await page.request.get(lookup,{headers:{authorization:'Bearer '+process.env.ELITE_QUOTE_READER}})).status()).toBe(403);expect((await page.request.get(lookup.replace('organization_id=store','organization_id=other'),{headers})).status()).toBe(403);expect((await page.request.get(lookup.replace(id,'unknown'),{headers})).status()).toBe(409);
        const path=action==='receive-return'?'/v1/franchise/return-authorizations/'+id+'/receive':'/v1/franchise/return-receipts/'+receiptId+'/decide';const data=action==='receive-return'?{organization_id:'store',serial_number:serial,condition_code:command.conditionCode,notes:command.notes.trim(),evidence_sha256:value.receipt.evidence_sha256}:{organization_id:'store',inventory_action:command.inventoryAction,notes:command.notes.trim()};expect((await page.request.post(api+path,{headers,data})).status()).toBe(409);expect((await page.request.post(api+path,{headers,data:{...data,actor_subject:'forged'}})).status()).toBe(400);
        await next().click();await expect(status()).toContainText('Caso actualizado');expect(posts).toBe(before+1);expect(await page.evaluate(()=>Object.keys(sessionStorage).filter(k=>k.startsWith('elite-return-operation:')))).toEqual([]);
        expect(await page.evaluate(()=>document.documentElement.scrollWidth<=document.documentElement.clientWidth)).toBe(true);
      }
      // Finish the operator's refresh before changing cookies to the reader session.
      const refreshed=page.waitForResponse(response=>{
        const request=response.request(),url=new URL(response.url()),headers=request.headers();
        return request.method()==='GET'&&url.origin===new URL(base).origin&&url.pathname==='/franchise'&&headers.rsc==='1'&&headers['next-router-prefetch']!=='1';
      });
      await panel().getByRole('button',{name:'Actualizar casos',exact:true}).click();
      // Headers prove the authenticated GET was sent; an RSC body may remain streaming.
      const refreshResponse=await refreshed;expect(refreshResponse.ok()).toBe(true);
      if(index<2)await expect(panel().getByRole('article',{name:'Caso '+['ui-return-loss','ui-exchange-invalid','ui-return-race'][index+1],exact:true})).toBeVisible();expect(posts).toBe(2*(index+1));
    }
    await setSession('READER',['admin:read']);await page.goto(base+'/franchise');await expect(page.getByRole('region',{name:'Operaciones de devoluciones',exact:true})).toHaveCount(0);return;
  }
  if(process.env.ELITE_QUOTE_PHASE==='operation-complete-checklist') {
    await setSession('HANDOVER',['handover:manage']);await page.goto(base+'/franchise');
    const panel=()=>page.getByRole('region',{name:'Presentación de entrega',exact:true}),status=()=>panel().getByRole('status',{name:'Resultado de presentación'});
    let posts=0;page.on('request',r=>{if(r.method()==='POST'&&r.url().endsWith('/api/enterprise/franchise/commands'))posts++});
    const api=process.env.ENTERPRISE_API_BASE_URL,headers={authorization:'Bearer '+process.env.ELITE_QUOTE_HANDOVER},serial=process.env.ELITE_DELIVERY_SERIAL;
    for(let index=0;index<3;index++){
      if(index>0){await panel().getByRole('button',{name:'Preparar otra presentación',exact:true}).click();await expect(panel().getByRole('button',{name:'Completar y presentar',exact:true})).toBeEnabled();expect(posts).toBe(index);}
      const id=['ui-completion-loss','ui-completion-invalid','ui-completion-race'][index];
      const fill=async()=>{await panel().getByLabel('Entrega',{exact:true}).fill(id);await panel().getByLabel('Versión actual de entrega',{exact:true}).fill('1');await panel().getByLabel('ID del checklist',{exact:true}).fill('ui-completion-checklist');await panel().getByLabel('Versión del checklist',{exact:true}).fill('1');await panel().getByLabel('Respuestas, una por línea',{exact:true}).fill('  serial = '+serial+'  \n confirmed = confirmed  ');};
      await fill();
      if(index===2){await page.evaluate(()=>{window.originalSet=Storage.prototype.setItem;Storage.prototype.setItem=function(k,v){if(k.startsWith('elite-checklist-complete:'))throw new Error('synthetic storage');return window.originalSet.call(this,k,v)}});await panel().getByRole('button',{name:'Completar y presentar',exact:true}).click();await expect(status()).toContainText('No se envió');expect(posts).toBe(2);await page.reload();await fill();}
      let command;
      await page.route('**/api/enterprise/franchise/commands',async route=>{
        command=route.request().postDataJSON();const responses=index===2?await Promise.all([route.fetch(),route.fetch()]):[await route.fetch()];expect(responses.map(r=>r.status()).sort()).toEqual(index===2?[200,409]:[200]);
        const success=responses.find(r=>r.status()===200),body=await success.json();expect(body).toMatchObject({id,organization_id:'store',version:2,state:'presented',checklist_id:'ui-completion-checklist',checklist_version:1});
        if(index===0)await route.abort('failed');else if(index===1)await route.fulfill({status:200,contentType:'application/json',body:'{}'});else await route.fulfill({response:success});
      },{times:1});
      await panel().getByRole('button',{name:'Completar y presentar',exact:true}).click();await expect(status()).toContainText(index<2?'No pudimos comprobar':'Respuesta recibida');expect(posts).toBe(index+1);
      await expect(panel().getByRole('button',{name:'Completar y presentar',exact:true})).toBeDisabled();await expect(panel().getByRole('button',{name:'Preparar otra presentación',exact:true})).toBeDisabled();
      const markers=await page.evaluate(()=>Object.entries(sessionStorage).filter(([k])=>k.startsWith('elite-checklist-complete:')).map(([,v])=>JSON.parse(v)));expect(markers).toHaveLength(1);expect(Object.keys(markers[0]).sort()).toEqual(['checklistId','checklistVersion','digest','handoverId','version']);expect(markers[0]).toMatchObject({handoverId:id,version:1,checklistId:'ui-completion-checklist',checklistVersion:1});expect(markers[0].digest).toMatch(/^[a-f0-9]{64}$/);expect(JSON.stringify(markers)).not.toContain(serial);
      if(index===0){const accepted=await page.request.post(api+'/v1/customer/handovers/'+id+'/accept',{headers:{authorization:'Bearer '+process.env.ELITE_QUOTE_CUSTOMER},data:{organization_id:'store',version:2,confirmed_received:true,serial_number:serial,checklist_id:'ui-completion-checklist',checklist_version:1}});expect(accepted.status()).toBe(200)}
      if(index===1){const rejected=await page.request.post(api+'/v1/customer/handovers/'+id+'/reject',{headers:{authorization:'Bearer '+process.env.ELITE_QUOTE_CUSTOMER},data:{organization_id:'store',version:2,reason_code:'serial-mismatch',details:'Synthetic rejection'}});expect(rejected.status()).toBe(200)}
      await page.reload();
      if(index===0){
        await page.route('**/api/enterprise/franchise/commands?*',r=>r.fulfill({status:503,contentType:'application/json',body:'{"private":"synthetic-detail"}'}),{times:1});await panel().getByRole('button',{name:'Consultar presentación',exact:true}).click();await expect(status()).toContainText('Conservamos el bloqueo');await expect(page.getByText('synthetic-detail')).toHaveCount(0);
        await page.route('**/api/enterprise/franchise/commands?*',async r=>{const response=await r.fetch(),body=await response.json();body.completion.responses[0].response_text='Different';await r.fulfill({status:200,contentType:'application/json',body:JSON.stringify(body)})},{times:1});await panel().getByRole('button',{name:'Consultar presentación',exact:true}).click();await expect(status()).toContainText('Conservamos el bloqueo');await expect(panel().getByRole('button',{name:'Preparar otra presentación',exact:true})).toBeDisabled();
      }
      await panel().getByRole('button',{name:'Consultar presentación',exact:true}).click();await expect(status()).toContainText('Presentación recuperada');await expect(panel()).toContainText(['accepted','rejected','presented'][index]);await expect(panel()).toContainText('completado por handover');expect(posts).toBe(index+1);
      const lookup=api+'/v1/franchise/handovers/'+id+'/checklist-result?organization_id=store';const recovered=await page.request.get(lookup,{headers});expect(recovered.status()).toBe(200);expect(recovered.headers()['cache-control']).toContain('no-store');const value=await recovered.json();expect(value.actor_subject).toBe('handover');expect(value.responses.map(r=>r.item_id)).toEqual(['confirmed','serial']);
      expect((await page.request.get(lookup,{headers:{authorization:'Bearer '+process.env.ELITE_QUOTE_READER}})).status()).toBe(403);expect((await page.request.get(lookup.replace('organization_id=store','organization_id=other'),{headers})).status()).toBe(403);expect((await page.request.get(lookup.replace(id,'unknown'),{headers})).status()).toBe(409);
      const data={organization_id:'store',version:1,checklist_id:command.checklistId,checklist_version:command.checklistVersion,responses:command.responses};expect((await page.request.post(api+'/v1/franchise/handovers/'+id+'/complete-checklist',{headers,data})).status()).toBe(409);expect((await page.request.post(api+'/v1/franchise/handovers/'+id+'/complete-checklist',{headers,data:{...data,actor_subject:'forged'}})).status()).toBe(400);
      expect(await page.evaluate(()=>document.documentElement.scrollWidth<=document.documentElement.clientWidth)).toBe(true);
    }
    await setSession('READER',['admin:read']);await page.goto(base+'/franchise');await expect(page.getByRole('button',{name:'Completar y presentar',exact:true})).toHaveCount(0);return;
  }
  if(process.env.ELITE_QUOTE_PHASE==='operation-publish-checklist') {
    await setSession('HANDOVER',['handover:manage']);await page.goto(base+'/franchise');
    const panel=()=>page.getByRole('heading',{name:'Publicar versión de checklist',exact:true}).locator('..');
    const status=()=>panel().getByRole('status',{name:'Resultado de publicación de checklist'});
    let posts=0;page.on('request',r=>{if(r.method()==='POST'&&r.url().endsWith('/api/enterprise/franchise/commands'))posts++});
    const api=process.env.ENTERPRISE_API_BASE_URL,headers={authorization:'Bearer '+process.env.ELITE_QUOTE_HANDOVER};
    for(let index=0;index<3;index++){
      if(index>0){await panel().getByRole('button',{name:'Preparar otra versión',exact:true}).click();await expect(panel().getByRole('button',{name:'Publicar versión',exact:true})).toBeEnabled();expect(posts).toBe(index);}
      const id=index===2?'ui-publication-concurrent':'ui-publication-versioned',version=index===1?2:1,title='Título sintético '+index;
      const fill=async()=>{await panel().getByLabel('ID del checklist',{exact:true}).fill(id);await panel().getByLabel('Versión',{exact:true}).fill(String(version));await panel().getByLabel('Título',{exact:true}).fill('  '+title+'  ');await panel().getByLabel('ID',{exact:true}).first().fill('serial');await panel().getByLabel('Pregunta o control',{exact:true}).first().fill('  Serie sintética  ');await panel().getByRole('combobox',{name:'Respuesta',exact:true}).first().selectOption('serial');await panel().getByRole('button',{name:'Agregar ítem',exact:true}).click();await panel().getByLabel('ID',{exact:true}).nth(1).fill('confirmed');await panel().getByLabel('Pregunta o control',{exact:true}).nth(1).fill('Confirmación sintética');};
      await fill();
      if(index===2){await page.evaluate(()=>{window.originalSet=Storage.prototype.setItem;Storage.prototype.setItem=function(k,v){if(k.startsWith('elite-checklist-publish:'))throw new Error('synthetic storage');return window.originalSet.call(this,k,v)}});await panel().getByRole('button',{name:'Publicar versión',exact:true}).click();await expect(status()).toContainText('No se envió');expect(posts).toBe(2);await page.reload();await fill();}
      let command;
      await page.route('**/api/enterprise/franchise/commands',async route=>{
        command=route.request().postDataJSON();expect(route.request().headers()['idempotency-key']).toBeUndefined();
        const responses=index===2?await Promise.all([route.fetch(),route.fetch()]):[await route.fetch()];
        expect(responses.map(r=>r.status()).sort()).toEqual(index===2?[201,409]:[201]);
        const success=responses.find(r=>r.status()===201),body=await success.json();expect(body).toMatchObject({id,organization_id:'store',version,title,state:'published'});expect(body.items.map(i=>i.ordinal)).toEqual([1,2]);
        if(index===0)await route.abort('failed');else if(index===1)await route.fulfill({status:201,contentType:'application/json',body:'{}'});else await route.fulfill({response:success});
      },{times:1});
      await panel().getByRole('button',{name:'Publicar versión',exact:true}).click();await expect(status()).toContainText(index<2?'No pudimos comprobar':'Versión publicada');expect(posts).toBe(index+1);
      await expect(panel().getByRole('button',{name:'Publicar versión',exact:true})).toBeDisabled();await expect(panel().getByRole('button',{name:'Preparar otra versión',exact:true})).toBeDisabled();
      const markers=await page.evaluate(()=>Object.entries(sessionStorage).filter(([k])=>k.startsWith('elite-checklist-publish:')).map(([,v])=>JSON.parse(v)));expect(markers).toHaveLength(1);expect(Object.keys(markers[0]).sort()).toEqual(['checklistId','digest','version']);expect(markers[0]).toMatchObject({checklistId:id,version});expect(markers[0].digest).toMatch(/^[a-f0-9]{64}$/);
      await page.reload();
      if(index===0){
        await page.route('**/api/enterprise/franchise/commands?*',r=>r.fulfill({status:503,contentType:'application/json',body:'{"private":"synthetic-detail"}'}),{times:1});await panel().getByRole('button',{name:'Consultar versión publicada',exact:true}).click();await expect(status()).toContainText('Conservamos el bloqueo');await expect(page.getByText('synthetic-detail')).toHaveCount(0);
        await page.route('**/api/enterprise/franchise/commands?*',async r=>{const response=await r.fetch(),body=await response.json();body.checklist.items[0].prompt='Distinto';await r.fulfill({status:200,contentType:'application/json',body:JSON.stringify(body)})},{times:1});await panel().getByRole('button',{name:'Consultar versión publicada',exact:true}).click();await expect(status()).toContainText('Conservamos el bloqueo');await expect(panel().getByRole('button',{name:'Preparar otra versión',exact:true})).toBeDisabled();
      }
      await panel().getByRole('button',{name:'Consultar versión publicada',exact:true}).click();await expect(status()).toContainText('Versión recuperada');await expect(panel()).toContainText(title);expect(posts).toBe(index+1);
      const lookup=api+'/v1/franchise/delivery-checklists/result?'+new URLSearchParams({organization_id:'store',checklist_id:id,version:String(version)});
      const recovered=await page.request.get(lookup,{headers});expect(recovered.status()).toBe(200);expect(recovered.headers()['cache-control']).toContain('no-store');expect((await recovered.json()).title).toBe(title);
      expect((await page.request.get(lookup,{headers:{authorization:'Bearer '+process.env.ELITE_QUOTE_READER}})).status()).toBe(403);
      expect((await page.request.get(lookup.replace('organization_id=store','organization_id=other'),{headers})).status()).toBe(403);
      expect((await page.request.get(lookup.replace('version='+version,'version=999'),{headers})).status()).toBe(409);
      expect((await page.request.get(lookup.replace('version='+version,'version=0'),{headers})).status()).toBe(400);
      const data={organization_id:'store',checklist_id:id,version,title:command.title,items:command.items};expect((await page.request.post(api+'/v1/franchise/delivery-checklists',{headers,data})).status()).toBe(409);
      expect((await page.request.post(api+'/v1/franchise/delivery-checklists',{headers,data:{...data,title:'Different'}})).status()).toBe(409);
      expect((await page.request.post(api+'/v1/franchise/delivery-checklists',{headers,data:{...data,actor_subject:'forged'}})).status()).toBe(400);
      if(index===1){const first=await page.request.get(lookup.replace('version=2','version=1'),{headers});expect(first.status()).toBe(200);expect((await first.json()).title).toBe('Título sintético 0');}
      expect(await page.evaluate(()=>document.documentElement.scrollWidth<=document.documentElement.clientWidth)).toBe(true);
    }
    await setSession('READER',['admin:read']);await page.goto(base+'/franchise');await expect(page.getByRole('button',{name:'Publicar versión',exact:true})).toHaveCount(0);
    return;
  }
  if(process.env.ELITE_QUOTE_PHASE==='operation-create-slot') {
    await setSession('PLANNER',['appointment:manage']);await page.goto(base+'/franchise');
    const section=()=>page.locator('section').filter({has:page.getByRole('heading',{name:'Publicar capacidad',exact:true})});
    const status=()=>section().getByRole('status',{name:'Resultado de turno'});
    let posts=0;page.on('request',r=>{if(r.method()==='POST'&&r.url().endsWith('/api/enterprise/franchise/commands'))posts++});
    const keys=[];
    for(let index=0;index<3;index++){
      if(index>0){await section().getByRole('button',{name:'Preparar otro turno',exact:true}).click();await expect(section().getByRole('button',{name:'Publicar turno',exact:true})).toBeEnabled();expect(posts).toBe(index);}
      const days=index===2?200:30+index;const kind=["service","delivery","test-drive"][index];
      const fill=async()=>{const local=await page.evaluate(days=>[0,3600000].map(extra=>{const d=new Date(Date.now()+days*86400000+extra);return new Date(d.getTime()-d.getTimezoneOffset()*60000).toISOString().slice(0,16)}),days);await section().getByLabel('Desde',{exact:true}).fill(local[0]);await section().getByLabel('Hasta',{exact:true}).fill(local[1]);await section().getByRole('combobox',{name:'Tipo',exact:true}).selectOption(kind);await section().getByLabel('Cupos',{exact:true}).fill('2');};
      await fill();
      if(index===2){await page.evaluate(()=>{window.originalSet=Storage.prototype.setItem;Storage.prototype.setItem=function(k,v){if(k.startsWith('elite-slot-create:'))throw new Error('synthetic storage');return window.originalSet.call(this,k,v)}});await section().getByRole('button',{name:'Publicar turno',exact:true}).click();await expect(status()).toContainText('No se envió');expect(posts).toBe(2);await page.reload();await fill();}
      let key,command,id;
      await page.route('**/api/enterprise/franchise/commands',async route=>{
        command=route.request().postDataJSON();key=route.request().headers()['idempotency-key'];expect(key).toMatch(/^slot-[a-f0-9-]{36}$/);expect(keys).not.toContain(key);keys.push(key);
        const responses=index===2?await Promise.all([route.fetch(),route.fetch()]):[await route.fetch()];
        for(const response of responses)expect(response.status()).toBe(201);
        const body=await responses[0].json();id=body.id;expect(body).toMatchObject({organization_id:'store',kind,capacity:2,booked:0,state:'open',version:1});
        if(index===2)expect((await responses[1].json()).id).toBe(id);
        if(index===0)await route.abort('failed');else if(index===1)await route.fulfill({status:201,contentType:'application/json',body:'{}'});else await route.fulfill({response:responses[0]});
      },{times:1});
      await section().getByRole('button',{name:'Publicar turno',exact:true}).click();await expect(status()).toContainText(index<2?'No pudimos comprobar':'Turno registrado');expect(posts).toBe(index+1);
      await expect(section().getByRole('button',{name:'Publicar turno',exact:true})).toBeDisabled();await expect(section().getByRole('button',{name:'Preparar otro turno',exact:true})).toBeDisabled();
      expect(await page.evaluate(()=>Object.entries(sessionStorage).filter(([k])=>k.startsWith('elite-slot-create:')).map(([,v])=>JSON.parse(v)))).toEqual([{requestKey:key}]);
      await page.reload();
      if(index===0){await page.route('**/api/enterprise/franchise/commands?*',r=>r.fulfill({status:503,contentType:'application/json',body:'{"private":"synthetic-detail"}'}),{times:1});await section().getByRole('button',{name:'Consultar turno creado',exact:true}).click();await expect(status()).toContainText('No pudimos recuperar');await expect(page.getByText('synthetic-detail')).toHaveCount(0);await expect(section().getByRole('button',{name:'Preparar otro turno',exact:true})).toBeDisabled();}
      await section().getByRole('button',{name:'Consultar turno creado',exact:true}).click();await expect(status()).toContainText('Consulta recuperada');await expect(section()).toContainText(id);expect(posts).toBe(index+1);
      const api=process.env.ENTERPRISE_API_BASE_URL,headers={authorization:'Bearer '+process.env.ELITE_QUOTE_PLANNER,'idempotency-key':key};
      const reordered={organization_id:'store',kind:command.kind,capacity:command.capacity,starts_at:command.startsAt,ends_at:command.endsAt};expect((await page.request.post(api+'/v1/franchise/appointment-slots',{headers,data:reordered})).status()).toBe(409);const data={organization_id:'store',kind:command.kind,starts_at:command.startsAt,ends_at:command.endsAt,capacity:command.capacity};
      const replay=await page.request.post(api+'/v1/franchise/appointment-slots',{headers,data});expect(replay.status()).toBe(200);expect((await replay.json()).id).toBe(id);
      expect((await page.request.post(api+'/v1/franchise/appointment-slots',{headers,data:{...data,ends_at:new Date(Date.parse(data.ends_at)+60000).toISOString()}})).status()).toBe(409);
      expect((await page.request.post(api+'/v1/franchise/appointment-slots',{headers,data:{...data,actor_subject:'forged'}})).status()).toBe(400);
      expect((await page.request.post(api+'/v1/franchise/appointment-slots',{headers:{authorization:headers.authorization},data})).status()).toBe(400);
      const lookup=api+'/v1/franchise/appointment-slots/result?'+new URLSearchParams({organization_id:'store',request_key:key});
      expect((await page.request.get(lookup,{headers:{authorization:'Bearer '+process.env.ELITE_QUOTE_READER}})).status()).toBe(403);
      expect((await page.request.get(lookup.replace('organization_id=store','organization_id=other'),{headers})).status()).toBe(403);
      expect((await page.request.get(lookup.replace(key,'unknown-slot-key'),{headers})).status()).toBe(409);
      expect(await page.evaluate(()=>document.documentElement.scrollWidth<=document.documentElement.clientWidth)).toBe(true);
    }
    await setSession('READER',['admin:read']);await page.goto(base+'/franchise');await expect(page.getByRole('button',{name:'Publicar turno',exact:true})).toHaveCount(0);
    return;
  }
  if(process.env.ELITE_QUOTE_PHASE==='operation-create-availability') {
    await setSession('SCHEDULER',['availability:read','availability:manage','inventory:allocate']);await page.goto(base+'/franchise');
    const section=()=>page.locator('section').filter({has:page.getByRole('heading',{name:'Jornadas y ausencias',exact:true})});
    const status=()=>section().getByRole('status',{name:'Resultado de intervalo'});
    let posts=0;page.on('request',r=>{if(r.method()==='POST'&&r.url().endsWith('/api/enterprise/franchise/commands'))posts++});
    const keys=[];
    for(let index=0;index<3;index++){
      if(index>0){await section().getByRole('button',{name:'Preparar otro intervalo',exact:true}).click();await expect(section().getByRole('button',{name:'Registrar intervalo',exact:true})).toBeEnabled();expect(posts).toBe(index);}
      const days=index===2?180:20+index;
      const fill=async()=>{await section().getByLabel('Desde',{exact:true}).fill(new Date(Date.now()+days*86400000).toISOString().slice(0,16));await section().getByLabel('Hasta',{exact:true}).fill(new Date(Date.now()+days*86400000+3600000).toISOString().slice(0,16));await section().getByRole('combobox',{name:'Tipo',exact:true}).selectOption(index===1?'unavailable':'working');if(index===1)await section().getByLabel('Motivo',{exact:true}).fill('synthetic-absence');};
      await fill();
      if(index===2){await page.evaluate(()=>{window.originalSet=Storage.prototype.setItem;Storage.prototype.setItem=function(k,v){if(k.startsWith('elite-availability-create:'))throw new Error('synthetic storage');return window.originalSet.call(this,k,v)}});await section().getByRole('button',{name:'Registrar intervalo',exact:true}).click();await expect(status()).toContainText('No se envió');expect(posts).toBe(2);await page.reload();await fill();}
      let key,command,id;
      await page.route('**/api/enterprise/franchise/commands',async route=>{
        command=route.request().postDataJSON();key=route.request().headers()['idempotency-key'];expect(key).toMatch(/^availability-[a-f0-9-]{36}$/);expect(keys).not.toContain(key);keys.push(key);
        const responses=index===2?await Promise.all([route.fetch(),route.fetch()]):[await route.fetch()];
        for(const response of responses)expect(response.status()).toBe(201);
        const body=await responses[0].json();id=body.id;expect(body).toMatchObject({organization_id:'store',entry_type:index===1?'unavailable':'working',state:'active',version:1});
        if(index===2)expect((await responses[1].json()).id).toBe(id);
        if(index===0)await route.abort('failed');else if(index===1)await route.fulfill({status:201,contentType:'application/json',body:'{}'});else await route.fulfill({response:responses[0]});
      },{times:1});
      await section().getByRole('button',{name:'Registrar intervalo',exact:true}).click();await expect(status()).toContainText(index<2?'No pudimos comprobar':'Intervalo registrado');expect(posts).toBe(index+1);
      await expect(section().getByRole('button',{name:'Registrar intervalo',exact:true})).toBeDisabled();await expect(section().getByRole('button',{name:'Preparar otro intervalo',exact:true})).toBeDisabled();
      expect(await page.evaluate(()=>Object.entries(sessionStorage).filter(([k])=>k.startsWith('elite-availability-create:')).map(([,v])=>JSON.parse(v)))).toEqual([{requestKey:key}]);
      await page.reload();
      if(index===0){await page.route('**/api/enterprise/franchise/commands?*',r=>r.fulfill({status:503,contentType:'application/json',body:'{"private":"synthetic-detail"}'}),{times:1});await section().getByRole('button',{name:'Consultar intervalo creado',exact:true}).click();await expect(status()).toContainText('No pudimos recuperar');await expect(page.getByText('synthetic-detail')).toHaveCount(0);await expect(section().getByRole('button',{name:'Preparar otro intervalo',exact:true})).toBeDisabled();}
      await section().getByRole('button',{name:'Consultar intervalo creado',exact:true}).click();await expect(status()).toContainText('Consulta recuperada');await expect(section()).toContainText(id);expect(posts).toBe(index+1);
      const api=process.env.ENTERPRISE_API_BASE_URL,headers={authorization:'Bearer '+process.env.ELITE_QUOTE_SCHEDULER,'idempotency-key':key};
      const data={organization_id:'store',resource_id:command.resourceId??'',entry_type:command.entryType,reason_code:command.reasonCode??'',starts_at:command.startsAt,ends_at:command.endsAt};
      const replay=await page.request.post(api+'/v1/franchise/availability',{headers,data});expect(replay.status()).toBe(200);expect((await replay.json()).id).toBe(id);
      expect((await page.request.post(api+'/v1/franchise/availability',{headers,data:{...data,ends_at:new Date(Date.parse(data.ends_at)+60000).toISOString()}})).status()).toBe(409);
      expect((await page.request.post(api+'/v1/franchise/availability',{headers,data:{...data,actor_subject:'forged'}})).status()).toBe(400);
      expect((await page.request.post(api+'/v1/franchise/availability',{headers:{authorization:headers.authorization},data})).status()).toBe(400);
      const lookup=api+'/v1/franchise/availability/result?'+new URLSearchParams({organization_id:'store',request_key:key});
      expect((await page.request.get(lookup,{headers:{authorization:'Bearer '+process.env.ELITE_QUOTE_AVAILABILITY}})).status()).toBe(403);
      expect((await page.request.get(lookup.replace('organization_id=store','organization_id=other'),{headers})).status()).toBe(403);
      expect((await page.request.get(lookup.replace(key,'unknown-availability-key'),{headers})).status()).toBe(409);
      if(index===2){
        const cancelled=await page.request.post(api+'/v1/franchise/availability/'+encodeURIComponent(id)+'/cancel',{headers,data:{organization_id:'store',version:1,reason_code:'schedule-correction'}});expect(cancelled.status()).toBe(200);
        await section().getByRole('button',{name:'Consultar intervalo creado',exact:true}).click();await expect(status()).toContainText('Consulta recuperada');await expect(section()).toContainText('cancelled');
        const after=await page.request.post(api+'/v1/franchise/availability',{headers,data});expect(after.status()).toBe(200);expect(await after.json()).toMatchObject({id,state:'cancelled',version:2});
      }
      expect(await page.evaluate(()=>document.documentElement.scrollWidth<=document.documentElement.clientWidth)).toBe(true);
    }
    await setSession('AVAILABILITY',['availability:read']);await page.goto(base+'/franchise');await expect(page.getByRole('button',{name:'Registrar intervalo',exact:true})).toHaveCount(0);
    return;
  }
  if(process.env.ELITE_QUOTE_PHASE==='operation-quote') {
    await setSession('WRITER',['lead:read','quote:write']);await page.goto(base+'/franchise');
    let posts=0;
    page.on('request',r=>{if(r.method()==='POST'&&r.url().endsWith('/api/enterprise/franchise/commands'))posts++});
    for(let index=0;index<3;index++) {
      const card=()=>page.locator('article').filter({has:page.getByRole('heading',{name:'quote-create-'+index,exact:true})});
      await card().getByLabel('Variante',{exact:true}).fill('variant');await card().getByLabel('Lista de precios',{exact:true}).fill('book');
      await card().getByLabel('Válida hasta',{exact:true}).fill(new Date(Date.now()+86400000).toISOString().slice(0,16));
      if(index===2){await page.evaluate(()=>{window.originalSet=Storage.prototype.setItem;Storage.prototype.setItem=function(k,v){if(k.startsWith('elite-quote-create:'))throw new Error('synthetic');return window.originalSet.call(this,k,v)}});await card().getByRole('button',{name:'Emitir cotización',exact:true}).click();await expect(card().getByRole('status',{name:'Resultado de cotización'})).toContainText('No se envió');expect(posts).toBe(2);await page.reload();await card().getByLabel('Variante',{exact:true}).fill('variant');await card().getByLabel('Lista de precios',{exact:true}).fill('book');await card().getByLabel('Válida hasta',{exact:true}).fill(new Date(Date.now()+86400000).toISOString().slice(0,16));}
      let key,command,id;
      await page.route('**/api/enterprise/franchise/commands',async route=>{
        command=route.request().postDataJSON();key=route.request().headers()['idempotency-key'];expect(key).toMatch(/^quote-[a-f0-9-]{36}$/);
        const responses=index===2?await Promise.all([route.fetch(),route.fetch()]):[await route.fetch()];
        for(const response of responses)expect(response.status()).toBe(201);
        const value=await responses[0].json();id=value.id;expect(value).toMatchObject({lead_id:'quote-create-'+index,organization_id:'store',currency:'ARS',total_minor_units:250000,state:'issued'});
        if(index===2)expect((await responses[1].json()).id).toBe(id);
        if(index===0)await route.abort('failed');else if(index===1)await route.fulfill({status:201,contentType:'application/json',body:'{}'});else await route.fulfill({response:responses[0]});
      },{times:1});
      await card().getByRole('button',{name:'Emitir cotización',exact:true}).click();
      await expect(card().getByRole('status',{name:'Resultado de cotización'})).toContainText(index<2?'No pudimos comprobar':'Cotización registrada');
      await expect(card().getByRole('button',{name:'Emitir cotización',exact:true})).toBeDisabled();expect(posts).toBe(index+1);
      const marker=await page.evaluate(i=>Object.entries(sessionStorage).filter(([k])=>k.startsWith('elite-quote-create:')&&k.endsWith(':quote-create-'+i)).map(([,v])=>JSON.parse(v)),index);expect(marker).toEqual([{requestKey:key}]);
      await page.reload();
      if(index===0){await page.route('**/api/enterprise/franchise/commands?*',route=>route.fulfill({status:503,contentType:'application/json',body:'{"internal":"synthetic-private"}'}),{times:1});await card().getByRole('button',{name:'Consultar cotización',exact:true}).click();await expect(card().getByRole('status',{name:'Resultado de cotización'})).toContainText('No pudimos recuperar');await expect(page.getByText('synthetic-private')).toHaveCount(0);}
      await card().getByRole('button',{name:'Consultar cotización',exact:true}).click();await expect(card().getByRole('status',{name:'Resultado de cotización'})).toContainText('Consulta recuperada');await expect(card()).toContainText(id);expect(posts).toBe(index+1);
      const lookup='/api/enterprise/franchise/commands?'+new URLSearchParams({organizationId:'store',leadId:'quote-create-'+index,requestKey:key});
      expect((await page.request.get(lookup.replace('organizationId=store','organizationId=other'))).status()).toBe(403);
      expect((await page.request.get(lookup.replace('leadId=quote-create-'+index,'leadId=absent'))).status()).toBe(409);
      expect((await page.request.get(lookup.replace(key,'quote-00000000-0000-4000-8000-000000000000'))).status()).toBe(409);
      const api=process.env.ENTERPRISE_API_BASE_URL;
      const readPath=api+'/v1/franchise/quotes/result?'+new URLSearchParams({organization_id:'store',lead_id:'quote-create-'+index,request_key:key});
      expect((await page.request.get(readPath,{headers:{authorization:'Bearer '+process.env.ELITE_QUOTE_OBSERVER}})).status()).toBe(403);
      const replay=await page.request.post(api+'/v1/franchise/quotes',{headers:{authorization:'Bearer '+process.env.ELITE_QUOTE_WRITER,'idempotency-key':key},data:{organization_id:'store',lead_id:command.leadId,variant_id:'variant',price_book_id:'book',valid_until:command.validUntil}});expect(replay.status()).toBe(200);expect((await replay.json()).id).toBe(id);
      const divergent=await page.request.post(api+'/v1/franchise/quotes',{headers:{authorization:'Bearer '+process.env.ELITE_QUOTE_WRITER,'idempotency-key':key},data:{organization_id:'store',lead_id:command.leadId,variant_id:'variant',price_book_id:'other',valid_until:command.validUntil}});expect(divergent.status()).toBe(409);
    }
    await setSession('OBSERVER',['lead:read']);await page.goto(base+'/franchise');await expect(page.getByRole('button',{name:'Emitir cotización',exact:true})).toHaveCount(0);
    return;
  }
  // Test-runner-only direct HTTP reads verify Go authorization independently of BFF.
  const read=async(kind='OPERATOR',organization='store')=>fetch(process.env.ENTERPRISE_API_BASE_URL+'/v1/commerce/orders?organization_id='+organization,{headers:{authorization:'Bearer '+process.env['ELITE_QUOTE_'+kind]}});
  const command=(order,stock,overrides={})=>({action:'allocate-order-stock',organizationId:'store',orderId:order.id,lineId:order.lines[0].id,stockUnitId:stock.id,orderVersion:order.version,stockVersion:stock.version,...overrides});
  const send=(body,origin=base)=>page.request.post(base+'/api/enterprise/franchise/commands',{headers:{origin},data:body,maxRedirects:0});
  if(process.env.ELITE_QUOTE_PHASE==='operation-availability') {
    const perms=['availability:read','availability:manage','inventory:allocate'];
    await setSession('SCHEDULER',perms);await page.goto(base+'/franchise');
    const section=page.locator('section').filter({has:page.getByRole('heading',{name:'Jornadas y ausencias',exact:true})});
    const publicBase=process.env.ENTERPRISE_API_BASE_URL+'/v1/public/quote-'+process.env.ELITE_QUOTE_TENANT.replaceAll('-','')+'/store';
    const slots=async()=>{
      const now=Date.now();
      const query=new URLSearchParams({kind:'consultation',from:new Date(now+3600000).toISOString(),to:new Date(now+14*86400000).toISOString()});
      const response=await fetch(publicBase+'/appointment-slots?'+query);expect(response.status).toBe(200);return (await response.json()).items;
    };
    const beforeSlots=await slots();expect(beforeSlots.map(s=>s.id).sort()).toEqual(['window-slot-0','window-slot-2','window-slot-4']);
    let posts=0;page.on('request',r=>{if(r.method()==='POST'&&r.url().includes('/api/enterprise/franchise/commands'))posts++});
    for(let index=0;index<4;index++) {
      const card=section.locator('article').nth(index),button=card.getByRole('button',{name:'Cancelar intervalo',exact:true});
      await expect(card).toHaveCount(1);let original;
      if(index===2) {
        await page.evaluate(()=>{Storage.prototype.setItem=function(){throw new Error('synthetic storage failure')}});
        await button.click();await expect(card.getByRole('status')).toContainText('No se envió');expect(posts).toBe(2);
        await page.reload();
      }
      await page.route('**/api/enterprise/franchise/commands',async route=>{
        original=route.request().postDataJSON();
        const responses=index===2?await Promise.all([route.fetch(),route.fetch()]):[await route.fetch()];
        if(index===2)expect(responses.map(r=>r.status()).sort()).toEqual([200,409]);
        const committed=responses.find(r=>r.status()===200);expect(committed).toBeTruthy();
        const body=await committed.json();expect(body.id).toBe('cancel-window-'+index);expect(body.version).toBe(2);expect(body.state).toBe('cancelled');
        if(index===0)await route.abort('failed');else if(index===1)await route.fulfill({status:200,contentType:'application/json',body:'{}'});else await route.fulfill({response:committed});
      },{times:1});
      await expect(button).toBeEnabled();await button.click();
      await expect(card.getByRole('status')).toContainText(index<2?'No pudimos comprobar la cancelación':'Cancelación registrada');
      await expect(button).toBeDisabled();expect(posts).toBe(index+1);
      expect(await page.evaluate(()=>Object.keys(sessionStorage).filter(k=>k.startsWith('elite-availability-cancel:')).map(k=>JSON.parse(sessionStorage.getItem(k))))).toEqual([{version:1}]);
      await card.getByRole('button',{name:'Consultar intervalo',exact:true}).click();
      if(index===0) {
        await expect(page.getByRole('main').getByRole('alert')).toContainText('No pudimos consultar disponibilidad');
        await expect(page.getByRole('main')).not.toContainText('PRIVATE-SYNTHETIC');
        await expect(page.getByRole('heading',{name:'Pedidos: stock, pago y entrega',exact:true})).toBeVisible();
        expect(await page.evaluate(()=>Object.keys(sessionStorage).filter(k=>k.startsWith('elite-availability-cancel:')).length)).toBe(1);
        await page.getByRole('link',{name:'Volver a consultar operación',exact:true}).click();
      }
      await expect(card.getByRole('status')).toContainText('Consulta recuperada: el intervalo figura cancelado');
      await expect(card.getByText('Cancelado',{exact:true})).toBeVisible();await expect(button).toBeDisabled();
      if(index%2===0) {
        expect((await slots()).some(s=>s.id==='window-slot-'+index)).toBe(false);
        const originalSlot=beforeSlots.find(s=>s.id==='window-slot-'+index);
        const attempt=await fetch(publicBase+'/appointments',{method:'POST',headers:{'content-type':'application/json','Idempotency-Key':'cancelled-slot-attempt-'+index},body:JSON.stringify({lead_id:'lead',kind:'consultation',starts_at:originalSlot.starts_at})});
        expect(attempt.status).toBe(409);
      }
      expect(posts).toBe(index+1);
      expect(await page.evaluate(()=>Object.keys(sessionStorage).filter(k=>k.startsWith('elite-availability-cancel:')).length)).toBe(0);
      expect((await send(original)).status()).toBe(409);
      expect((await send({...original,organizationId:'other'})).status()).toBe(403);
      expect((await send({...original,actorSubject:'forged'})).status()).toBe(400);
      await setSession('AVAILABILITY',['availability:read','availability:manage']);expect((await send(original)).status()).toBe(403);
      await setSession('SCHEDULER',perms);await page.reload();await expect(button).toBeDisabled();
    }
    const blocked=section.locator('article').nth(4);
    await blocked.getByRole('button',{name:'Cancelar intervalo',exact:true}).click();
    await expect(blocked.getByRole('status')).toContainText('El intervalo cambió o tiene restricciones');
    expect(posts).toBe(5);await expect(blocked.getByRole('button',{name:'Cancelar intervalo',exact:true})).toBeDisabled();
    await blocked.getByRole('button',{name:'Consultar intervalo',exact:true}).click();
    await expect(blocked.getByRole('status')).toContainText('Hay una cancelación pendiente de comprobar');
    await expect(blocked.getByText('Activo',{exact:true})).toBeVisible();
    await expect(blocked.getByRole('button',{name:'Cancelar intervalo',exact:true})).toBeDisabled();expect(posts).toBe(5);
    expect(await page.evaluate(()=>Object.keys(sessionStorage).filter(k=>k.startsWith('elite-availability-cancel:')).map(k=>JSON.parse(sessionStorage.getItem(k))))).toEqual([{version:1}]);
    await setSession('AVAILABILITY',['availability:read']);await page.reload();
    await expect(section.getByRole('button',{name:'Cancelar intervalo',exact:true})).toHaveCount(0);
    await section.locator('article').first().getByText('Ayuda para cancelar un intervalo',{exact:true}).click();
    await expect(section).toContainText('availability-cancel-view/1.0.0');
    expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth)).toBe(true);return;
  }
  if(process.env.ELITE_QUOTE_PHASE==='operation-lead') {
    await setSession('LEAD',['lead:read','lead:assign','lead:update']);await page.goto(base+'/franchise');
    const card=page.locator('article').filter({has:page.getByRole('heading',{name:'operator-lead',exact:true})});
    let posts=0;page.on('request',r=>{if(r.method()==='POST'&&r.url().includes('/api/enterprise/franchise/commands'))posts++});
    for(const [index,target] of ['sales-owner','contacted','qualified','lost'].entries()) {
      await expect(card).toHaveCount(1);
      const action=index===0?'assign-lead':'transition-lead';
      const button=card.getByRole('button',{name:index===0?'Asignar':'Cambiar estado',exact:true});
      const fill=async()=>{if(index===0)await card.getByRole('textbox',{name:'Responsable',exact:true}).fill(target);else await card.getByRole('combobox',{name:'Nuevo estado',exact:true}).selectOption(target)};
      await fill();let original;
      if(index===2) {
        await page.evaluate(()=>{Storage.prototype.setItem=function(){throw new Error('synthetic storage failure')}});
        await button.click();await expect(card.getByRole('status')).toContainText('No se envió');expect(posts).toBe(2);
        await page.reload();await fill();
      }
      await page.route('**/api/enterprise/franchise/commands',async route=>{
        original=route.request().postDataJSON();
        const responses=index===2?await Promise.all([route.fetch(),route.fetch()]):[await route.fetch()];
        if(index===2)expect(responses.map(r=>r.status()).sort()).toEqual([200,409]);
        const committed=responses.find(r=>r.status()===200);expect(committed).toBeTruthy();const body=await committed.json();
        expect(body.id).toBe('operator-lead');expect(body.version).toBe(index+2);expect(index===0?body.assigned_subject:body.state).toBe(target);
        if(index===0)await route.abort('failed');else if(index===1)await route.fulfill({status:200,contentType:'application/json',body:'{}'});else await route.fulfill({response:committed});
      },{times:1});
      await expect(button).toBeEnabled();await button.click();
      await expect(card.getByRole('status')).toContainText(index<2?'No pudimos comprobar el resultado':'Cambio registrado');
      await expect(button).toBeDisabled();expect(posts).toBe(index+1);
      const markers=await page.evaluate(()=>Object.keys(sessionStorage).filter(k=>k.startsWith('elite-lead-command:')).map(k=>JSON.parse(sessionStorage.getItem(k))));
      expect(markers).toEqual([{version:index+1,action}]);
      await card.getByRole('button',{name:'Consultar lead',exact:true}).click();
      if(index===0) {
        await expect(page.getByRole('main').getByRole('alert')).toContainText('No pudimos consultar leads');
        await expect(page.getByRole('main')).not.toContainText('PRIVATE-SYNTHETIC');
        expect(await page.evaluate(()=>Object.keys(sessionStorage).filter(k=>k.startsWith('elite-lead-command:')).length)).toBe(1);
        await page.getByRole('link',{name:'Volver a consultar operación',exact:true}).click();
      }
      await expect(card.getByRole('status')).toContainText('Consulta recuperada');
      await expect(card).toContainText(index===0?'Responsable: sales-owner':target+' · fixture');
      expect(posts).toBe(index+1);
      expect(await page.evaluate(()=>Object.keys(sessionStorage).filter(k=>k.startsWith('elite-lead-command:')).length)).toBe(0);
      expect((await send(original)).status()).toBe(409);
      expect((await send({...original,organizationId:'other'})).status()).toBe(403);
      expect((await send({...original,actorSubject:'forged'})).status()).toBe(400);
      await setSession('OBSERVER',['lead:read','lead:assign','lead:update']);expect((await send(original)).status()).toBe(403);
      await setSession('LEAD',['lead:read','lead:assign','lead:update']);
      await page.reload();await expect(card).toContainText(index===0?'Responsable: sales-owner':target+' · fixture');
    }
    await expect(card.getByRole('button',{name:'Asignar',exact:true})).toBeDisabled();
    await expect(card.getByRole('button',{name:'Cambiar estado',exact:true})).toHaveCount(0);
    await card.getByText('Ayuda para actualizar un lead',{exact:true}).click();await expect(card).toContainText('lead-command-view/1.0.0');
    expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth)).toBe(true);return;
  }
  if(process.env.ELITE_QUOTE_PHASE==='operation-resolution') {
    // Delivery recovery remains independent from the lead command lane.
    await setSession('HANDOVER',['handover:manage']);await page.goto(base+'/franchise');
    let posts=0;page.on('request',r=>{if(r.method()==='POST'&&r.url().includes('/api/enterprise/franchise/commands'))posts++});
    for(const [index,action] of ['correct-and-represent','return','exchange'].entries()) {
      const id=(index===0?'reject':action)+'-fixture-handover';
      const card=page.locator('article').filter({hasText:'Entrega '+id});
      await expect(card).toHaveCount(1);
      const fill=async()=>{await card.getByRole('combobox').selectOption(action);await card.getByRole('textbox',{name:'Notas',exact:true}).fill('Synthetic authorized correction')};
      await fill();
      const button=card.getByRole('button',{name:'Registrar resolución',exact:true});
      if(index===2) {
        await page.evaluate(()=>{Storage.prototype.setItem=function(){throw new Error('synthetic storage outage')}});
        await button.click();await expect(card.getByRole('status')).toContainText('No se envió');
        expect(posts).toBe(2);await page.reload();await fill();
      }
      let original;
      await page.route('**/api/enterprise/franchise/commands',async route=>{
        original=route.request().postDataJSON();
        const responses=index===2?await Promise.all([route.fetch(),route.fetch()]):[await route.fetch()];
        if(index===2)expect(responses.map(r=>r.status()).sort()).toEqual([200,409]);
        const committed=responses.find(r=>r.status()===200);expect(committed).toBeTruthy();
        const body=await committed.json();expect(body.exception.state).toBe('resolved');expect(body.exception.resolution_action).toBe(action);
        if(index===0) {expect(body.successor_handover.state).toBe('prepared');await route.abort('failed')}
        else if(index===1) {expect(body.return_authorization_id).toBeTruthy();await route.fulfill({status:200,contentType:'application/json',body:'{}'})}
        else {expect(body.return_authorization_id).toBeTruthy();await route.fulfill({response:committed})}
      },{times:1});
      await expect(button).toBeEnabled();await button.click();
      await expect(card.getByRole('status')).toContainText(index===2?'Resolución registrada':'No pudimos comprobar el resultado');
      await expect(button).toBeDisabled();expect(posts).toBe(index+1);
      const markers=await page.evaluate(()=>Object.keys(sessionStorage).filter(k=>k.startsWith('elite-delivery-resolution:')).map(k=>JSON.parse(sessionStorage.getItem(k))));
      expect(markers).toHaveLength(1);expect(markers[0]).toEqual({version:1,action});
      await card.getByRole('button',{name:'Consultar resolución',exact:true}).click();
      if(index===0) {
        await expect(page.getByRole('main').getByRole('alert')).toContainText('No pudimos consultar discrepancias de entrega');
        await expect(page.getByRole('main')).not.toContainText('PRIVATE-SYNTHETIC');
        expect(await page.evaluate(()=>Object.keys(sessionStorage).filter(k=>k.startsWith('elite-delivery-resolution:')).length)).toBe(1);
        expect(posts).toBe(1);
        await page.getByRole('link',{name:'Volver a consultar operación',exact:true}).click();
      }
      await expect(card).toContainText('Resolución: '+action);
      await expect(card.getByRole('status')).toContainText('Consulta recuperada');
      await expect(card.getByRole('button',{name:'Registrar resolución',exact:true})).toHaveCount(0);
      expect(posts).toBe(index+1);
      expect(await page.evaluate(()=>Object.keys(sessionStorage).filter(k=>k.startsWith('elite-delivery-resolution:')).length)).toBe(0);
      expect((await send(original)).status()).toBe(409);
      expect((await send({...original,organizationId:'other'})).status()).toBe(403);
      await setSession('OBSERVER',['handover:manage']);
      expect((await send(original)).status()).toBe(403);
      await setSession('HANDOVER',['handover:manage']);
      await card.getByText('Ayuda para resolver una discrepancia',{exact:true}).click();
      await expect(card).toContainText('delivery-resolution-view/1.0.0');
      await page.reload();await expect(card).toContainText('Resolución: '+action);
    }
    expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth)).toBe(true);
    return;
  }
  if(process.env.ELITE_QUOTE_PHASE==='operation-sections') {
    let posts=0;page.on('request',r=>{if(r.method()==='POST'&&r.url().includes('/api/enterprise/franchise/commands'))posts++});
    await setSession('OBSERVER',['lead:read']);await page.goto(base+'/franchise');
    await expect(page.getByRole('heading',{name:'Operación comercial, agenda y entregas',exact:true})).toBeVisible();
    for(const name of ['Preparación versionada de entregas','Recursos y turnos','Publicar capacidad','Jornadas y ausencias']) await expect(page.getByRole('heading',{name,exact:true})).toHaveCount(0);
    for(const name of ['Asignar','Cambiar estado','Emitir cotización']) await expect(page.getByRole('button',{name,exact:true})).toHaveCount(0);
    await setSession('HANDOVER',['handover:manage']);await page.goto(base+'/franchise');
    await expect(page.getByRole('heading',{name:'Pedidos: stock, pago y entrega',exact:true})).toBeVisible();
    await expect(page.getByRole('main').getByRole('alert')).toContainText('No pudimos consultar discrepancias de entrega');
    await expect(page.getByRole('main')).not.toContainText('PRIVATE-SYNTHETIC');
    await expect(page.getByRole('button',{name:'Registrar resolución',exact:true})).toHaveCount(0);
    await expect(page.getByRole('heading',{name:'Preparación versionada de entregas',exact:true})).toBeVisible();
    await page.getByRole('link',{name:'Volver a consultar operación',exact:true}).click();
    await expect(page.getByRole('heading',{name:'Discrepancias de entrega',exact:true})).toBeVisible();
    await expect(page.getByRole('main').getByRole('alert')).toHaveCount(0);
    await page.getByText('Ayuda para recuperar una sección',{exact:true}).click();
    await expect(page.getByText('operation-sections-view/1.0.0',{exact:false})).toBeVisible();
    await setSession('RESOURCE',['resource:manage']);await page.goto(base+'/franchise');
    await expect(page.getByRole('heading',{name:'Recursos y turnos',exact:true})).toBeVisible();
    await expect(page.getByRole('heading',{name:'Preparación versionada de entregas',exact:true})).toHaveCount(0);
    await setSession('AVAILABILITY',['availability:read']);await page.goto(base+'/franchise');
    await expect(page.getByRole('heading',{name:'Jornadas y ausencias',exact:true})).toBeVisible();
    await expect(page.getByRole('button',{name:'Registrar intervalo',exact:true})).toHaveCount(0);
    expect(posts).toBe(0);
    await setSession('OBSERVER',['handover:manage']);
    const denied=await send({action:'publish-delivery-checklist',organizationId:'store',checklistId:'denied-probe',version:1,title:'Synthetic probe',items:[{id:'serial',prompt:'Serial',response_type:'serial',required:true}]});
    expect(denied.status()).toBe(403);
    const readDenied=await fetch(process.env.ENTERPRISE_API_BASE_URL+'/v1/franchise/delivery-exceptions?organization_id=store',{headers:{authorization:'Bearer '+process.env.ELITE_QUOTE_OBSERVER}});
    expect(readDenied.status).toBe(403);
    expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth)).toBe(true);
    return;
  }
  await setSession();await page.goto(base+'/franchise');
  await expect(page.getByRole('heading',{name:'Pedidos: stock, pago y entrega',exact:true})).toBeVisible();
  if(process.env.ELITE_QUOTE_PHASE==='operation-read-failure') {
    await expect(page.getByRole('main').getByRole('alert')).toContainText('No pudimos consultar los pedidos');
    await expect(page.getByRole('button',{name:'Reservar unidad',exact:true})).toHaveCount(0);
    await expect(page.getByRole('button',{name:'Actualizar pedidos'})).toBeVisible();return;
  }
  const response=await read();expect(response.status).toBe(200);const snapshot=await response.json();
  expect(snapshot.orders).toHaveLength(3);expect(snapshot.truncated).toBe(false);
  if(process.env.ELITE_QUOTE_PHASE==='operation-restart') {
    expect(snapshot.orders.flatMap(o=>o.lines).filter(l=>l.stock_id)).toHaveLength(2);
    expect(snapshot.stock).toHaveLength(0);
    await expect(page.getByText('Unidad reservada:',{exact:false})).toHaveCount(2);
    await expect(page.getByRole('button',{name:'Sin stock disponible',exact:true})).toBeDisabled();return;
  }
  const [first,second,third]=snapshot.orders;
  const lost=snapshot.stock.find(s=>s.id==='stock-lost'), race=snapshot.stock.find(s=>s.id==='stock-race');
  expect((await send(command(first,lost,{organizationId:'other'}))).status()).toBe(403);
  expect((await send(command(first,lost), 'https://other.invalid')).status()).toBe(403);
  expect((await send(command(first,lost,{orderVersion:99}))).status()).toBe(409);
  expect((await send(command(first,lost,{lineId:second.lines[0].id}))).status()).toBe(409);
  expect((await send(command(first,lost,{stockUnitId:'missing'}))).status()).toBe(409);
  expect((await read('CUSTOMER')).status).toBe(403);
  expect((await read('OPERATOR','other')).status).toBe(403);
  // Backend denial must survive a forged BFF permission claim.
  await setSession('CUSTOMER',['inventory:allocate']);
  expect((await send(command(first,lost))).status()).toBe(403);
  await setSession('READER',['admin:read']);await page.reload();
  await expect(page.getByRole('button',{name:'Reservar unidad',exact:true})).toHaveCount(0);
  expect((await send(command(first,lost))).status()).toBe(403);
  await context.clearCookies();expect((await send(command(first,lost))).status()).toBe(401);
  await setSession();await page.reload();
  const card=page.getByRole('article',{name:'Pedido '+first.id,exact:true});
  await card.getByRole('combobox').selectOption('stock-lost');
  await page.route('**/api/enterprise/franchise/commands',async route=>{
    const committed=await route.fetch();expect(committed.status()).toBe(200);await route.abort('failed');
  },{times:1});
  await card.getByRole('button',{name:'Reservar unidad'}).focus();await page.keyboard.press('Enter');
  await expect(page.getByRole('status',{name:'Resultado de reserva'})).toContainText('No pudimos comprobar el resultado');
  await expect(card.getByRole('button',{name:'Reservar unidad'})).toBeDisabled();
  await page.getByRole('button',{name:'Actualizar pedidos'}).click();
  await expect(card).toContainText('Unidad reservada: stock-lost');
  expect((await send(command(first,lost))).status()).toBe(409);
  const results=await Promise.all([send(command(second,race)),send(command(third,race))]);
  expect(results.map(r=>r.status()).sort()).toEqual([200,409]);
  await page.reload();await expect(page.getByText('Unidad reservada:',{exact:false})).toHaveCount(2);
  const final=await(await read()).json();expect(final.stock).toHaveLength(0);
  expect(final.orders.flatMap(o=>o.payments)).toHaveLength(0);expect(final.orders.flatMap(o=>o.handovers)).toHaveLength(0);
  await page.getByText('Ayuda para continuar un pedido',{exact:true}).click();
  await expect(page.getByText('order-operations-view/1.1.0',{exact:false})).toBeVisible();
  expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth)).toBe(true);
  await page.screenshot({path:test.info().outputPath('order-stock.png'),fullPage:true});
});

test('operator payment request connects frontend BFF identity and durable intent',async({page,context})=>{
  test.skip(process.env.ELITE_QUOTE_E2E!=='1','explicit connected fixture required');test.setTimeout(60000);
  const {createHash,randomUUID}=await import('node:crypto');const {createRequire}=await import('node:module');const {resolve}=await import('node:path');const {pathToFileURL}=await import('node:url');
  const require=createRequire(resolve(process.env.ELITE_WEB_ROOT,'package.json'));const {EncryptJWT}=await import(pathToFileURL(require.resolve('jose')).href);
  const base=process.env.ELITE_BASE_URL;
  const session=async(kind='PAYMENT',permissions=['payment:create'])=>{
    await context.clearCookies();const value=await new EncryptJWT({subject:kind.toLowerCase(),tenantId:process.env.ELITE_QUOTE_TENANT,permissions,organizations:['store'],accessToken:process.env['ELITE_QUOTE_'+kind]}).setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('300s').encrypt(createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest());
    await context.addCookies([{name:'__Host-elite_session',value,url:base+'/',httpOnly:true,secure:true,sameSite:'Lax'}]);
  };
  const read=async()=>{const r=await fetch(process.env.ENTERPRISE_API_BASE_URL+'/v1/commerce/orders?organization_id=store',{headers:{authorization:'Bearer '+process.env.ELITE_QUOTE_PAYMENT}});expect(r.status).toBe(200);return r.json()};
  const command=(order,key=randomUUID())=>({action:'request-order-payment',organizationId:'store',orderId:order.id,requestKey:key});
  const send=(body,origin=base)=>page.request.post(base+'/api/enterprise/franchise/commands',{headers:{origin},data:body,maxRedirects:0});
  await session();await page.goto(base+'/franchise');const snapshot=await read();expect(snapshot.orders).toHaveLength(3);
  if(process.env.ELITE_QUOTE_PHASE==='payment-disabled'){
    expect(snapshot.payment_provider).toBe('');await expect(page.getByRole('button',{name:'Registrar solicitud de pago'})).toHaveCount(0);
    const unused=snapshot.orders.find(o=>!o.payments.length);expect((await send(command(unused))).status()).toBe(503);
    expect((await read()).orders.flatMap(o=>o.payments)).toHaveLength(2);return;
  }
  expect(snapshot.payment_provider).toBe('stripe');
  if(process.env.ELITE_QUOTE_PHASE==='payment-restart'){
    expect(snapshot.orders.flatMap(o=>o.payments)).toHaveLength(2);await expect(page.getByRole('button',{name:'Registrar solicitud de pago'})).toHaveCount(1);
    const paid=snapshot.orders.find(o=>o.payments.length);expect((await send(command(paid))).status()).toBe(409);return;
  }
  const [first,second,third]=snapshot.orders;const card=page.getByRole('article',{name:'Pedido '+first.id,exact:true});
  expect((await send({...command(third),amount:1})).status()).toBe(400);
  expect((await send({...command(third),provider:'mercadopago'})).status()).toBe(400);
  expect((await send({...command(third),organizationId:'other'})).status()).toBe(403);
  expect((await send(command(third),'https://other.invalid')).status()).toBe(403);
  await session('READER',['admin:read']);await page.reload();await expect(page.getByRole('button',{name:'Registrar solicitud de pago'})).toHaveCount(0);
  expect((await send(command(third))).status()).toBe(403);
  await session('CUSTOMER',['payment:create']);expect((await send(command(third))).status()).toBe(403);
  await session('OTHER_TENANT',['payment:create']);expect((await send(command(third))).status()).toBe(403);
  await context.clearCookies();expect((await send(command(third))).status()).toBe(401);
  await session();await page.reload();
  let retainedKey;let committedID;
  // Storage failure must make zero requests, not lose the replay identity.
  let attempts=0;const countRequests=req=>{if(req.url().includes('/api/enterprise/franchise/commands'))attempts++};page.on('request',countRequests);
  await page.evaluate(()=>{Storage.prototype.setItem=function(){throw new Error('synthetic denied')}});
  await card.getByRole('button',{name:'Registrar solicitud de pago'}).click();await expect(card.getByRole('status')).toContainText('No se envió la solicitud');expect(attempts).toBe(0);
  page.off('request',countRequests);await page.reload();
  await page.route('**/api/enterprise/franchise/commands',async route=>{
    const input=route.request().postDataJSON();retainedKey=input.requestKey;
    const committed=await route.fetch();expect(committed.status()).toBe(200);const receipt=await committed.json();expect(receipt.amount_minor_units).toBe(250000);expect(receipt.currency).toBe('ARS');committedID=receipt.id;await route.abort('failed');
  },{times:1});
  const paymentButton=card.getByRole('button',{name:'Registrar solicitud de pago'});
  await expect(paymentButton).toBeEnabled();await paymentButton.focus();await expect(paymentButton).toBeFocused();await page.keyboard.press('Enter');
  await expect(card.getByRole('status')).toContainText('No pudimos comprobar la solicitud');await expect(card.getByRole('button',{name:'Registrar solicitud de pago'})).toBeDisabled();
  await page.getByRole('button',{name:'Actualizar pedidos'}).click();await expect(card).toContainText(committedID);await expect(card.getByRole('button',{name:'Registrar solicitud de pago'})).toHaveCount(0);
  const recovered=await send(command(first,retainedKey));expect(recovered.status()).toBe(200);expect((await recovered.json()).id).toBe(committedID);
  expect((await send(command(first))).status()).toBe(409);
  const keys=Array.from({length:8},()=>randomUUID());const results=await Promise.all(keys.map(key=>send(command(second,key))));
  expect(results.filter(r=>r.status()===200)).toHaveLength(1);expect(results.filter(r=>r.status()===409)).toHaveLength(7);
  const winner=results.findIndex(r=>r.status()===200);const original=await results[winner].json();
  const replay=await send(command(second,keys[winner]));expect(replay.status()).toBe(200);expect((await replay.json()).id).toBe(original.id);
  expect((await send(command(second))).status()).toBe(409);
  const final=await read();expect(final.orders.flatMap(o=>o.payments)).toHaveLength(2);
  await page.reload();await expect(page.getByRole('button',{name:'Registrar solicitud de pago'})).toHaveCount(1);
  await page.getByText('Ayuda para continuar un pedido',{exact:true}).click();await expect(page.getByText('Pago: registrar una solicitud guarda una intención del total',{exact:false})).toBeVisible();
  expect(await page.evaluate(()=>document.documentElement.scrollWidth<=innerWidth)).toBe(true);
  await page.screenshot({path:test.info().outputPath('payment-intent.png'),fullPage:true});
});

test('customer quote acceptance connects portal identity durable order and recovery', async ({ page, context }) => {
  test.skip(process.env.ELITE_QUOTE_E2E !== '1', 'Requires isolated Go/PostgreSQL quote harness');
  test.setTimeout(60000);
  const { createHash } = await import('node:crypto');
  const { createRequire } = await import('node:module');
  const { resolve } = await import('node:path');
  const { pathToFileURL } = await import('node:url');
  const require = createRequire(resolve(process.env.ELITE_WEB_ROOT, 'package.json'));
  const { EncryptJWT } = await import(pathToFileURL(require.resolve('jose')).href);
  const base = process.env.ELITE_BASE_URL;
  expect(new URL(base).hostname).toBe('127.0.0.1');
  expect(new URL(base).protocol).toBe('https:');
  const setSession = async (kind='CUSTOMER', permissions=['customer:self'], organizations=['store']) => {
    await context.clearCookies();
    const value = await new EncryptJWT({subject:kind==='STRANGER'?'stranger':'customer',
      tenantId:process.env.ELITE_QUOTE_TENANT, permissions, organizations,
      accessToken:process.env['ELITE_QUOTE_'+kind]})
      .setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('300s')
      .encrypt(createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest());
    await context.addCookies([{name:'__Host-elite_session',value,url:base+'/',httpOnly:true,secure:true,sameSite:'Lax'}]);
  };
  const card = id => page.getByRole('article').filter({has:page.getByRole('heading',{name:id,exact:true})});
  const send = (id, version=1, organizationId='store') => page.request.post(base+'/api/enterprise/franchise/commands',
    {headers:{origin:base},data:{action:'accept-quote',organizationId,quoteId:id,version},maxRedirects:0});
  await setSession();
  if(process.env.ELITE_QUOTE_PHASE==='delivery-action-loss') {
    const commands=[];page.on('request',request=>{if(request.method()==='POST'&&request.url().includes('/api/enterprise/franchise/commands'))commands.push(request.url());});
    const delivery=id=>page.getByRole('article').filter({has:page.getByRole('heading',{name:/^Pedido /})}).nth(id==='read-fixture-handover'?1:0);
    await page.goto(base+'/customer/handovers');
    for(const [id,action,state] of [['read-fixture-handover','accept-handover','accepted'],['reject-fixture-handover','reject-handover','open']]) {
      await page.route('**/api/enterprise/franchise/commands',async route=>{
        expect(route.request().postDataJSON().action).toBe(action);
        const response=await route.fetch();expect(response.status()).toBe(200);
        expect((await response.json()).state).toBe(state);
        await route.abort('failed');
      },{times:1});
      if(action==='accept-handover') {
        await delivery(id).getByLabel('Número de serie observado').fill(process.env.ELITE_DELIVERY_SERIAL);
        await delivery(id).getByRole('checkbox').check();
        await delivery(id).getByRole('button',{name:'Registrar recepción',exact:true}).click();
      }else{
        await delivery(id).getByLabel('Código de motivo').fill('visible-damage');
        await delivery(id).getByLabel('Detalle',{exact:true}).fill('Synthetic discrepancy');
        await delivery(id).getByRole('button',{name:'Rechazar esta presentación',exact:true}).click();
      }
      await expect(page.getByRole('status')).toContainText('No pudimos comprobar el resultado');
      for(const button of await page.getByRole('button',{name:/registrar recepción|rechazar/i}).all())await expect(button).toBeDisabled();
      await page.getByRole('link',{name:'Consultar estado de entregas',exact:true}).click();
      await expect(delivery(id)).toContainText(action==='accept-handover'?'Estado: accepted':'Estado: rejected');
      await expect(delivery(id).getByRole('button')).toHaveCount(0);
    }
    expect(commands).toHaveLength(2);
    await expect(page.getByText('Synthetic discrepancy',{exact:true})).toBeVisible();
    return;
  }
  if(process.env.ELITE_QUOTE_PHASE?.startsWith('delivery-read-')) {
    const commands=[];page.on('request',request=>{if(request.method()==='POST'&&request.url().includes('/api/enterprise/franchise/commands'))commands.push(request.url());});
    await page.goto(base+'/customer/handovers');
    await expect(page.getByRole('heading',{name:'Mis entregas',exact:true})).toBeVisible();
    if(process.env.ELITE_QUOTE_PHASE==='delivery-read-corrupt') {
      await expect(page.getByRole('main').getByRole('alert')).toContainText('No pudimos verificar tus entregas');
      await expect(page.getByRole('article')).toHaveCount(0);
      await page.getByRole('link',{name:'Volver a consultar entregas',exact:true}).click();
      await expect(page.getByRole('main').getByRole('alert')).toContainText('No pudimos verificar tus entregas');
    } else {
      await expect(page.getByRole('main').getByRole('alert')).toHaveCount(0);
      await expect(page.getByRole('article')).toHaveCount(1);
      await expect(page.getByRole('article')).toContainText('Sin aceptación disponible hasta completar la preparación exacta.');
    }
    await expect(page.getByText('PRIVATE-FOREIGN-STOCK',{exact:false})).toHaveCount(0);
    await expect(page.getByRole('button',{name:/registrar recepción|rechazar/i})).toHaveCount(0);
    expect(commands).toEqual([]);
    expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth)).toBe(true);
    return;
  }
  await page.goto(base+'/customer/quotes');
  await expect(page.getByRole('heading',{name:'Mis cotizaciones',exact:true})).toBeVisible();
  if(process.env.ELITE_QUOTE_PHASE !== 'read-failure') {
    await expect(card('quote-expired').getByRole('button',{name:'Aceptar y crear pedido'})).toHaveCount(0);
    await expect(card('quote-expired')).toContainText('Vigencia finalizada');
    await expect(card('quote-success')).toContainText(/ARS\s*2[.]500,00/);
    await expect(card('quote-success')).toContainText('UTC');
    await expect(page.getByText(/unidades menores/)).toHaveCount(0);
  }
  if(process.env.ELITE_QUOTE_PHASE==='read-failure') {
    await expect(page.getByRole('main').getByRole('alert')).toContainText('No pudimos consultar tus cotizaciones');
    await expect(page.getByRole('article')).toHaveCount(0);
    await expect(page.getByRole('link',{name:'Volver a consultar',exact:true})).toBeVisible();
    await page.getByText('Ayuda para aceptar una cotización',{exact:true}).click();
    await expect(page.getByText('quote-acceptance-view/1.0.0',{exact:false})).toBeVisible();
    const commands=[];page.on('request',request=>{if(request.method()==='POST'&&request.url().includes('/api/enterprise/franchise/commands'))commands.push(request.url());});
    await page.goto(base+'/customer/handovers');
    await expect(page.getByRole('heading',{name:'Mis entregas',exact:true})).toBeVisible();
    await expect(page.getByRole('main').getByRole('alert')).toContainText('No pudimos verificar tus entregas');
    await expect(page.getByRole('link',{name:'Volver a consultar entregas',exact:true})).toBeVisible();
    await expect(page.getByRole('button',{name:/aceptar|rechazar/i})).toHaveCount(0);
    await page.getByText('Ayuda para consultar entregas',{exact:true}).click();
    await expect(page.getByText('handover-read-view/1.0.0',{exact:false})).toBeVisible();
    await page.getByRole('link',{name:'Volver a consultar entregas',exact:true}).click();
    await expect(page.getByRole('main').getByRole('alert')).toContainText('No pudimos verificar tus entregas');
    expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth)).toBe(true);
    expect(commands).toEqual([]);
    return;
  }
  if(process.env.ELITE_QUOTE_PHASE==='restart') {
    for(const id of ['quote-lost','quote-success','quote-race']) {
      await expect(card(id)).toContainText('Pedido creado:');
      await expect(card(id).getByRole('button',{name:'Aceptar y crear pedido',exact:true})).toHaveCount(0);
    }
    expect((await send('quote-lost')).status()).toBe(409);
    await page.goto(base+'/customer/handovers');
    await expect(page.getByRole('heading',{name:'Mis entregas',exact:true})).toBeVisible();
    await expect(page.getByRole('main').getByRole('alert')).toHaveCount(0);
    await expect(page.getByText('Por ahora no hay entregas disponibles para consultar.',{exact:true})).toBeVisible();
    await expect(page.getByRole('button',{name:/aceptar|rechazar/i})).toHaveCount(0);
    return;
  }
  // Lose only the response AFTER real BFF/API/PostgreSQL commit.
  await page.route('**/api/enterprise/franchise/commands',async route=>{
    const response=await route.fetch();
    expect(response.status()).toBe(200);
    const committed=await response.json();
    expect(committed.state).toBe('accepted');
    expect(committed.order_id).toMatch(/^[0-9a-f-]{36}$/);
    await route.abort('failed');
  },{times:1});
  await card('quote-lost').getByRole('button',{name:'Aceptar y crear pedido',exact:true}).focus();
  await page.keyboard.press('Enter');
  await expect(page.getByRole('status')).toContainText('No pudimos comprobar el resultado');
  await expect(card('quote-lost').getByRole('button',{name:'Aceptar y crear pedido',exact:true})).toBeDisabled();
  await page.getByRole('button',{name:'Actualizar estado',exact:true}).click();
  await expect(card('quote-lost')).toContainText('Pedido creado:');
  await expect(card('quote-lost').getByRole('button',{name:'Aceptar y crear pedido',exact:true})).toHaveCount(0);
  expect((await send('quote-lost')).status()).toBe(409);

  await card('quote-success').getByRole('button',{name:'Aceptar y crear pedido',exact:true}).click();
  await expect(card('quote-success')).toContainText('Pedido creado:');
  const results=await Promise.all([send('quote-race'),send('quote-race')]);
  expect(results.map(r=>r.status()).sort()).toEqual([200,409]);
  const success=await results.find(r=>r.status()===200).json();
  expect(success.order_id).toMatch(/^[0-9a-f-]{36}$/);
  await page.reload();
  await expect(card('quote-race')).toContainText(success.order_id);

  expect((await send('quote-denied',99)).status()).toBe(409);
  expect((await send('quote-expired')).status()).toBe(409);
  expect((await send('quote-denied',1,'other')).status()).toBe(403);
  await setSession('OBSERVER', ['lead:read']);
  expect((await send('quote-denied')).status()).toBe(403);
  await page.goto(base+'/customer/quotes');
  await expect(page.getByRole('heading',{name:'Acceso denegado',exact:true})).toBeVisible();
  await setSession('STRANGER');
  expect((await send('quote-denied')).status()).toBe(409);
  await page.goto(base+'/customer/quotes');
  await expect(page.getByRole('article')).toHaveCount(0);
  await setSession('OTHER_TENANT');
  expect((await send('quote-denied')).status()).toBe(409);
  await context.clearCookies();
  expect((await send('quote-denied')).status()).toBe(401);
  await setSession();
  await page.goto(base+'/customer/quotes');
  await page.getByText('Ayuda para aceptar una cotización',{exact:true}).click();
  await expect(page.getByText('quote-acceptance-view/1.0.0',{exact:false})).toBeVisible();
  await expect(page.getByText('Aceptar una cotización no confirma el pago, el stock ni la entrega.',{exact:true})).toBeVisible();
  expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth)).toBe(true);
  await page.screenshot({path:test.info().outputPath('quote-acceptance.png'),fullPage:true});
});

// AUTHORED regression: real pages, same session, existing immutable return commands.
async function exerciseReturnTabs(context, base, setSession) {
  const ids=['ui-return-loss','ui-exchange-invalid','ui-return-race'];
  const serial=process.env.ELITE_DELIVERY_SERIAL;
  const api=process.env.ENTERPRISE_API_BASE_URL;
  const headers={authorization:'Bearer '+process.env.ELITE_QUOTE_HANDOVER};
  const region=p=>p.getByRole('region',{name:'Operaciones de devoluciones',exact:true});
  const status=p=>region(p).getByRole('status',{name:'Resultado de devolución'});
  const consult=p=>region(p).getByRole('button',{name:'Consultar operación de devolución',exact:true});
  const next=p=>region(p).getByRole('button',{name:'Continuar operando',exact:true});
  const references=p=>p.evaluate(()=>Object.entries(sessionStorage).filter(([k])=>k.startsWith('elite-return-operation:')).map(([key,raw])=>({key,value:JSON.parse(raw)})));
  const evidence=[],sessionEvidence=[];let posts=0,losses=0,popupCopies=0;
  const observer=request=>{if(request.method()==='POST'&&request.url().endsWith('/api/enterprise/franchise/commands'))posts++};
  context.on('request',observer);
  try {
    for(const [caseIndex,id] of ids.entries()) {
      for(const action of ['receive-return','decide-return']) {
        // A new pair has separate sessionStorage and the same authenticated cookie.
        const pages=[await context.newPage(),await context.newPage()];
        const different=caseIndex===1||(caseIndex===2&&action==='decide-return');
        const loseWinner=caseIndex===0;
        const label=action==='receive-return'?'Registrar recepción':'Registrar decisión';
        const arrived=[],responses=[],commands=[];let release;
        const barrier=new Promise(resolve=>{release=resolve});
        const before=posts;
        try {
          await Promise.all(pages.map(p=>p.goto(base+'/franchise')));
          for(const [index,p] of pages.entries()) {
            expect(await references(p)).toEqual([]);
            const card=region(p).getByRole('article',{name:'Caso '+id,exact:true});
            const notes='Synthetic multi-tab '+caseIndex+' '+action+(different?' tab-'+index:' identical');
            if(action==='receive-return') {
              await card.getByLabel('Serie recibida',{exact:true}).fill(serial);
              await card.getByLabel('Inspección/observaciones',{exact:true}).fill(notes);
            } else {
              await card.getByLabel('Fundamento',{exact:true}).fill(notes);
              if(different&&index===1)await card.getByRole('combobox',{name:'Acción de inventario solicitada',exact:true}).selectOption('repair');
            }
            await p.route('**/api/enterprise/franchise/commands',async route=>{
              commands[index]=route.request().postDataJSON();arrived.push(index);
              await barrier;
              const response=await route.fetch();
              const body=await response.json();
              responses.push({index,status:response.status(),body});
              if(response.status()===200&&loseWinner){losses++;await route.abort('failed')}
              else await route.fulfill({response});
            },{times:1});
          }
          await Promise.all(pages.map(p=>region(p).getByRole('article',{name:'Caso '+id,exact:true}).getByRole('button',{name:label,exact:true}).click()));
          await expect.poll(()=>arrived.length).toBe(2);
          // Neither request reaches the backend until both UI submissions arrive.
          expect(responses).toHaveLength(0);release();
          await expect.poll(()=>responses.length).toBe(2);
          expect(responses.map(r=>r.status).sort()).toEqual([200,409]);
          const winner=responses.find(r=>r.status===200).index,loser=1-winner;
          await expect(status(pages[loser])).toContainText('No pudimos comprobar');
          await expect(status(pages[winner])).toContainText(loseWinner?'No pudimos comprobar':'Respuesta recibida');
          const saved=await Promise.all(pages.map(references));
          expect(saved.map(x=>x.length)).toEqual([1,1]);
          expect(saved[0][0].key).toBe(saved[1][0].key);
          for(const [index,items] of saved.entries()) {
            expect(Object.keys(items[0].value).sort()).toEqual(['action','authorizationId','digest','receiptId']);
            expect(items[0].value).toMatchObject({action,authorizationId:id});
            expect(items[0].value.digest).toMatch(/^[a-f0-9]{64}$/);
            expect(JSON.stringify(items)).not.toContain(serial);
            expect(JSON.stringify(items)).not.toContain(commands[index].notes);
            await expect(next(pages[index])).toBeDisabled();
            await expect(region(pages[index]).getByRole('article',{name:'Caso '+id,exact:true}).getByRole('button',{name:label,exact:true})).toBeDisabled();
          }
          expect(saved[0][0].value.digest===saved[1][0].value.digest).toBe(!different);
          if(process.env.ELITE_RETURN_SESSION_E2E==='1'&&caseIndex===0) {
            sessionEvidence.push(await exerciseReturnSession(context,pages[winner],base,setSession,commands[winner],saved[winner]));
          }
          await Promise.all(pages.map(p=>p.reload()));
          expect(await Promise.all(pages.map(references))).toEqual(saved);
          await Promise.all(pages.map(p=>consult(p).click()));
          await expect(status(pages[winner])).toContainText('Caso recuperado');
          await expect(next(pages[winner])).toBeEnabled();
          if(different) {
            await expect(status(pages[loser])).toContainText('Conservamos el bloqueo');
            await expect(next(pages[loser])).toBeDisabled();
            await expect(region(pages[loser]).getByRole('region',{name:'Evidencia recuperada de devolución'})).toHaveCount(0);
          } else {
            await expect(status(pages[loser])).toContainText('Caso recuperado');
            await expect(next(pages[loser])).toBeEnabled();
          }
          const url=api+'/v1/franchise/returns/result?'+new URLSearchParams({organization_id:'store',authorization_id:id});
          const response=await context.request.get(url,{headers});expect(response.status()).toBe(200);
          const value=await response.json(),committed=action==='receive-return'?value.receipt:value.disposition;
          expect(committed.notes).toBe(commands[winner].notes);
          expect(action==='receive-return'?committed.received_by_subject:committed.decided_by_subject).toBe('handover');
          if(action==='decide-return') {
            expect(value.disposition.inventory_action).toBe(commands[winner].inventoryAction);
            expect(value.disposition.effect_requests).toHaveLength(caseIndex===1?3:4);
            expect(value.disposition.effect_requests.every(e=>e.state==='requested')).toBe(true);
          }
          if(caseIndex===0&&action==='receive-return') {
            // The opener copy is independent too: clearing its parent cannot clear it.
            const popupPromise=pages[winner].waitForEvent('popup');
            await pages[winner].evaluate(()=>{window.open(location.href,'_blank')});
            const popup=await popupPromise;
            try {
              await popup.waitForLoadState('domcontentloaded');
              expect(await references(popup)).toEqual(saved[winner]);
              await expect(next(popup)).toBeDisabled();
              await next(pages[winner]).click();
              expect(await references(pages[winner])).toEqual([]);
              expect(await references(popup)).toEqual(saved[winner]);
              await consult(popup).click();await expect(status(popup)).toContainText('Caso recuperado');
              await next(popup).click();expect(await references(popup)).toEqual([]);popupCopies++;
            } finally {await popup.close()}
          } else await next(pages[winner]).click();
          expect(await references(pages[winner])).toEqual([]);
          expect(await references(pages[loser])).toEqual(saved[loser]);
          if(!different) {await next(pages[loser]).click();expect(await references(pages[loser])).toEqual([])}
          else {
            // The losing tab's retained digest continues to reject the other intent.
            await pages[loser].reload();await consult(pages[loser]).click();
            await expect(status(pages[loser])).toContainText('Conservamos el bloqueo');
            await expect(next(pages[loser])).toBeDisabled();
          }
          expect(posts).toBe(before+2);
          const final=await context.request.get(url,{headers});expect(await final.json()).toEqual(value);
          evidence.push({authorization:id,action,different,loseWinner,winner,statuses:responses.map(r=>({tab:r.index,status:r.status})),posts:posts-before,exactDurableState:true});
        } finally {release();await Promise.all(pages.map(p=>p.close()))}
      }
    }
    expect(evidence).toHaveLength(6);expect(posts).toBe(12);expect(losses).toBe(2);expect(popupCopies).toBe(1);
    await test.info().attach('return-multitab-evidence',{body:JSON.stringify({evidence,posts,losses,popupCopies},null,2),contentType:'application/json'});
    console.log('RETURN_MULTITAB_BROWSER_PASS races=6 posts=12 wins=6 conflicts=6 losses=2 opener_copy=1');
    if(process.env.ELITE_RETURN_SESSION_E2E==='1') {
      expect(sessionEvidence).toHaveLength(2);
      await test.info().attach('return-session-evidence',{body:JSON.stringify(sessionEvidence,null,2),contentType:'application/json'});
      console.log('RETURN_SESSION_BOUNDARY_PASS transitions=8 denied=12 authorized_reads=2 restored=2');
    }
  } finally {context.off('request',observer)}
}

// AUTHORED: session changes do not turn retained client references into authority.
async function exerciseReturnSession(context,page,base,setSession,command,saved) {
  const region=()=>page.getByRole('region',{name:'Operaciones de devoluciones',exact:true});
  const consult=()=>region().getByRole('button',{name:'Consultar operación de devolución',exact:true});
  const next=()=>region().getByRole('button',{name:'Continuar operando',exact:true});
  const references=()=>page.evaluate(()=>Object.entries(sessionStorage).filter(([k])=>k.startsWith('elite-return-operation:')).map(([key,raw])=>({key,value:JSON.parse(raw)})));
  const url=base+'/api/enterprise/franchise/commands?'+new URLSearchParams({kind:'return',organizationId:'store',authorizationId:saved[0].value.authorizationId});
  // Observe real storage reads after navigation without changing their result.
  await page.addInitScript(()=>{
    window.returnScopeReads=[];const original=Storage.prototype.getItem;
    Storage.prototype.getItem=function(key){if(this===sessionStorage&&key.startsWith('elite-return-operation:'))window.returnScopeReads.push(key);return original.call(this,key)};
  });
  const scopeRead=async(subject,organization)=>{
    const {createHash}=await import('node:crypto');
    const key='elite-return-operation:'+createHash('sha256').update(JSON.stringify([process.env.ELITE_QUOTE_TENANT,subject,organization])).digest('hex');
    await expect.poll(()=>page.evaluate(()=>window.returnScopeReads??[])).toContain(key);
    if(subject!=='handover'||organization!=='store')expect(await page.evaluate(()=>window.returnScopeReads)).not.toContain(saved[0].key);
  };
  let denied=0;
  const reject=async(response,status,code)=>{
    expect(response.status()).toBe(status);
    const value=await response.json();
    if(code)expect(value.code).toBe(code);
    expect(JSON.stringify(value)).not.toContain(command.notes);
    expect(value).not.toHaveProperty('receipt');expect(value).not.toHaveProperty('disposition');
    denied++;
  };
  const post=()=>context.request.post(base+'/api/enterprise/franchise/commands',{headers:{origin:base,'content-type':'application/json'},data:command});
  // Existing UI remains open while its cookie changes to a less privileged user.
  await setSession('READER',['admin:read']);
  const response=page.waitForResponse(r=>r.request().method()==='GET'&&r.url().includes('/api/enterprise/franchise/commands?'));
  await consult().click();await reject(await response,403,'FORBIDDEN');
  await expect(next()).toBeDisabled();expect(await references()).toEqual(saved);
  await reject(await post(),403,'FORBIDDEN');
  await page.goto(base+'/franchise');await expect(region()).toHaveCount(0);
  expect(await references()).toEqual(saved);
  await context.clearCookies();
  await reject(await context.request.get(url),401,'UNAUTHENTICATED');
  await reject(await post(),401,'UNAUTHENTICATED');
  expect(await references()).toEqual(saved);
  // The fixture session selects another organization; the store command is denied
  // by BFF membership, and the actual API token remains limited to store.
  await setSession('HANDOVER',['handover:manage'],['other']);
  await reject(await context.request.get(url),403,'FORBIDDEN');
  await reject(await post(),403,'ORGANIZATION_FORBIDDEN');
  await page.goto(base+'/franchise');await scopeRead('handover','other');await expect(consult()).toBeDisabled();
  expect(await references()).toEqual(saved);
  // A second authorized operator may read shared organization evidence, but does
  // not inherit the first operator's local pending intent or rewrite its actor.
  await setSession('HANDOVER_OTHER',['handover:manage']);
  await page.goto(base+'/franchise');await scopeRead('handover_other','store');await expect(region().getByRole('article',{name:'Caso ui-return-race',exact:true}).getByRole('button',{name:'Registrar recepción',exact:true})).toBeEnabled();await expect(consult()).toBeDisabled();
  expect(await references()).toEqual(saved);
  const allowed=await context.request.get(url);expect(allowed.status()).toBe(200);
  const {returnCase:result}=await allowed.json();
  const committed=command.action==='receive-return'?result.receipt:result.disposition;
  expect(committed.notes).toBe(command.notes);
  expect(command.action==='receive-return'?committed.received_by_subject:committed.decided_by_subject).toBe('handover');
  await setSession('HANDOVER',['handover:manage']);await page.goto(base+'/franchise');await scopeRead('handover','store');
  expect(await references()).toEqual(saved);await expect(consult()).toBeEnabled();
  await expect(next()).toBeDisabled();expect(denied).toBe(6);
  return {action:command.action,transitions:4,denied,authorizedReads:1,restored:true};
}
