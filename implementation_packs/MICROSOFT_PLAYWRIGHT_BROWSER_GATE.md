# Microsoft Playwright Browser Gate

V317: el harness espera la respuesta RSC exitosa (headers) de Actualizar casos antes de
cambiar la sesión operador/lector. Conserva seis POST y el rechazo del lector;
sin sleeps, retries ni cambios TLS. Evidencia y límites:
reconstruction_evidence/TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md.

V312: corrección AUTHORED de recepción y decisión de devoluciones.
Consulta exacta por autorización, scope de grafo bloqueado en escritura,
lectura atómica, actor y evidencia inmutables. UI conserva huella y consulta
sin reenviar aun con lista inaccesible. Refund4/exchange3 solicitudes.
Evidencia reconstruction_evidence/RETURN_OPERATIONS_RECOVERY_V312.md.
Sin efectos downstream, cambio comercial ni promoción integral.

V311: corrección AUTHORED de recuperación de checklist completada por
identidad natural de entrega, actor y respuestas inmutables. Consulta scoped
preserva estado vigente accepted/rejected/presented; no reenvía ni reabre.
Evidencia reconstruction_evidence/CHECKLIST_COMPLETION_RECOVERY_V311.md.
Sin cambio de reglas comerciales, upstreams ni promoción integral.

V310: corrección AUTHORED de publicación de checklist con consulta exacta
por identidad natural org/id/versión, actor atómico y recuperación sin POST.
Evidencia reconstruction_evidence/CHECKLIST_PUBLICATION_RECOVERY_V310.md.
No altera inmutabilidad, reglas de entrega, upstreams ni promoción integral.

V309: corrección AUTHORED de publicación/consulta de turnos con identidad
retenida, Idempotency-Key y recuperación scoped sin reenvío. Pruebas y límites
en reconstruction_evidence/SLOT_CREATION_RECOVERY_V309.md. No cambio de reglas de agenda,
upstreams ni promoción integral; FAIL457 conserva los otros comandos pendientes.

V308: corrección AUTHORED de alta/consulta de recursos con identidad
retenida, Idempotency-Key y recuperación scoped sin reenvío. Pruebas y límites
en reconstruction_evidence/RESOURCE_CREATION_RECOVERY_V308.md. No cambio de reglas de agenda,
upstreams ni promoción integral; FAIL457 conserva los otros comandos pendientes.

V307: corrección AUTHORED de alta/consulta de disponibilidad con identidad
retenida, Idempotency-Key y recuperación scoped sin reenvío. Pruebas y límites
en reconstruction_evidence/AVAILABILITY_CREATION_RECOVERY_V307.md. No cambio de reglas de agenda,
upstreams ni promoción integral; FAIL457 conserva los otros comandos pendientes.

V306: corrección AUTHORED de emisión/consulta de cotización con identidad
retenida, Idempotency-Key y recuperación scoped sin reenvío. Pruebas y límites
en reconstruction_evidence/QUOTE_CREATION_RECOVERY_V306.md. No cambio de pricing,
upstreams ni promoción integral; FAIL457 conserva los otros comandos pendientes.

V294 / 0.1.7 incorpora el recorrido opt-in ELITE_QUOTE_E2E=1 en el spec existente.
Sólo mediante TestQuoteAcceptanceBrowserPostgres: datos PostgreSQL sintéticos,
sesión/JWKS reales de fixture, loopback TLS y build del web exacto. Rechazos,
concurrencia, pérdida de respuesta, consulta fallida, teclado y recuperación
tras reiniciar Next/recrear API-pool. No login externo, pagos ni TLS productivo.
Los artefactos usan directorio distinto por ejecución; nunca borrar el fallo
anterior mediante un retry. No existe retry automático configurado.
El harness/UI/tests son AUTHORED; Playwright oficial conserva pin/licencia.
Evidencia y hallazgo Firefox teardown: reconstruction_evidence/QUOTE_ACCEPTANCE_CONNECTED_V294.md.

## 1. Metadata

V263 adds an opt-in operator agenda gate on the existing Microsoft runtime: synthetic encrypted session, real RS256/JWKS verifier, built Next BFF, Go API and PostgreSQL. Four browser projects exercise name-based assignment, confirmation and durable reload; three lose a post-commit response deliberately. A local Go TLS proxy replaces the failed browser routing shim. Self-signed certificate errors are allowed only for this explicit loopback fixture; production PKI and OIDC login remain unproven. Config/tests/harness remain AUTHORED, with no new dependencies or copied upstream code. Evidence: `reconstruction_evidence/OPERATOR_AGENDA_BROWSER_POSTGRES_V263.md`.

```yaml
pack_id: "MICROSOFT-PLAYWRIGHT-BROWSER-GATE"
pack_version: "0.1.37"
status:
  authority: ELITE_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un gate de navegador separado que instala @playwright/test 1.62.1 oficial de Microsoft con lock exacto, prueba Chromium desktop/mobile, Firefox y WebKit desktop, y verifica semántica, headers, overflow y nonce CSP fresco/propagado del home público del perfil web Elite sin afirmar roles, backend, accesibilidad AT ni producción."
stacks: ["Microsoft Playwright 1.62.1", "Chromium 151.0.7922.34 rev 1234", "Firefox 153.0 rev 1538", "WebKit 26.5 rev 2336", "Node.js >=20", "pnpm 11.25.0", "PowerShell 7"]
compatible_with: ["TS-GO-API-WEB-BRIDGE 0.5.16", "TS-OIDC-PORTAL-ADAPTER 0.2.4", "GO-ELECTROMOBILITY-PUBLIC-CRM-API 0.2.3 for opt-in connected lead gate", "GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.10.19 and TS-FRANCHISE-JOURNEY-PORTALS 0.14.19 for appointment/agenda/quote gates", "PYTHON-META-WHATSAPP-CLOUD-ADAPTER 0.13.0 for opt-in notification history", "MARKDOWN-COMPOSITOR-CORE 0.2.x"]
incompatible_with: ["credentials in URLs", "non-loopback HTTP targets", "production claims from the local smoke", "node_modules committed as source"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0"
upstream_sources: ["https://github.com/microsoft/playwright/tree/26a9e470a7b3c7822084b09fb7f13902c5f37b51", "https://registry.npmjs.org/@playwright/test/-/test-1.62.1.tgz"]
verified_at: "2026-09-08"
```

## 2. Applicability

V302 / 0.1.12 añade ELITE_OPERATOR_SECTIONS_E2E=1 al harness de pedidos; requiere
ELITE_ORDER_E2E=1. Prueba cuatro permisos de operador, consulta503 recuperable,
pedidos independientes y rechazo Go aunque la sesión del fixture simule más
permisos. Tests AUTHORED, runtime Microsoft sin cambios de pin/licencia.
Ver reconstruction_evidence/OPERATOR_SECTIONS_RECOVERY_V302.md.

V300 / 0.1.11 extiende el harness existente con ELITE_DELIVERY_READ_E2E=1
y ELITE_DELIVERY_ACTION_E2E=1, dependientes de los gates previos. Comprueba
lectura incoherente desde PostgreSQL, recuperación GET y pérdida de respuestas
después de aceptar/rechazar realmente en el backend local. Actas iniciales son
fixtures explícitos; no prueban el creador ordinario pendiente ni liberación
financiera. No cambiar el pin oficial de Playwright ni atribuir estos tests
AUTHORED a Microsoft. Ver reconstruction_evidence/DELIVERY_READ_RECOVERY_V300.md.

V298 añade regresión AUTHORED de solicitud de pago conectada al gate existente:
browser/BFF/Go/OIDC/PG, respuesta perdida post-commit, clave retenida, ocho claves
concurrentes para el mismo pedido, permisos y campos inyectados, storage rechazado,
desactivación y relectura tras reiniciar procesos. Playwright oficial permanece
en su pin existente; los tests locales no se atribuyen a Microsoft. Ver V298.

V278 agrega un gate opt-in de historial de notificación sobre Next construido, Go/PostgreSQL reales y credenciales sintéticas. Cuatro proyectos verifican recuperación GET, scope, ayuda, responsive acotado y cero envíos nuevos. El harness preexige build y CLI aislado; no demuestra Meta ni IdP productivos. Evidencia: `reconstruction_evidence/WHATSAPP_HOST_OPERATOR_HISTORY_V278.md`. Runtime, pin, licencia y dependencias Microsoft permanecen sin cambios; integración AUTHORED.

V260 añade un gate opt-in separado del smoke histórico, ejecutado por `TestPublicLeadBrowserPostgres`: navegador → BFF → Go → PostgreSQL, con captura pública, replay, conflicto, rechazos y efectos durables. El README explica el comando y las condiciones descartables. No amplía el claim a turnos confirmados, usuarios autenticados, operadores, proveedores ni producción. La revisión/firma/licencia del runtime Microsoft no cambia; el harness y las pruebas nuevas son `AUTHORED`. Evidencia: `reconstruction_evidence/PUBLIC_LEAD_BROWSER_POSTGRES_V260.md`.

Use after composing the current enterprise web profile when a project needs a deterministic browser smoke before provider-specific authenticated journeys. The official Microsoft runtime is acquired as an exact npm dependency; all configuration, target tests and verification glue are Elite `AUTHORED` and never attributed to Microsoft.

Reject the default home smoke as proof of backend effects; those require the separate opt-in harness and its database assertions. Neither gate proves authenticated portals, assistive-technology interoperability, performance or production readiness. Remote home targets must use HTTPS. Local HTTP is limited to loopback, and URLs containing credentials fail during configuration load. The connected fixture gate never targets a remote project.

## 3. Architecture contract

The browser gate remains in its own module and never modifies the web application's `package.json` or lock. A frozen four-package pnpm graph installs Playwright Test 1.62.1, Playwright/Core 1.62.1 and the platform-conditional fsevents entry. The exact signed Microsoft release commit, GitHub archive, npm tarball, Apache license, Chromium revision and dated OSV result are retained in `source-lock.json`.

The runtime smoke uses isolated `page.setContent` and proves the real browser API in four projects: Chromium desktop/mobile, Firefox desktop and WebKit desktop. The bounded target tests start a previously built enterprise web app on loopback: four check one H1, main landmark, catalog link, locale, four security-header properties, no horizontal overflow and no browser/console errors; four require production CSP without unsafe inline/eval, a valid fresh nonce per response and exact nonce propagation to every rendered script. Traces, screenshots and video are retained only on failure. For local WebKit only, exact loopback subresources upgraded by the production CSP are proxied back to the same loopback HTTP server; the header is preserved and asserted. No secret, token, customer datum or external provider is required.

## 4. Exact file manifest

```text
CREATE microsoft_playwright_browser_gate/tests/notification-status.spec.mjs
CREATE microsoft_playwright_browser_gate/package.json
CREATE microsoft_playwright_browser_gate/pnpm-lock.yaml
CREATE microsoft_playwright_browser_gate/source-lock.json
CREATE microsoft_playwright_browser_gate/LICENSE.playwright.txt
CREATE microsoft_playwright_browser_gate/playwright.config.mjs
CREATE microsoft_playwright_browser_gate/tests/runtime.spec.mjs
CREATE microsoft_playwright_browser_gate/tests/enterprise-web.spec.mjs
CREATE microsoft_playwright_browser_gate/tests/nonce-csp.spec.mjs
CREATE microsoft_playwright_browser_gate/verify_contract.ps1
CREATE microsoft_playwright_browser_gate/README.md
CREATE microsoft_playwright_browser_gate/tests/journey-help.spec.mjs
CREATE microsoft_playwright_browser_gate/fixtures/help-loopback-proxy.go
CREATE microsoft_playwright_browser_gate/tests/role-workspace.spec.mjs
```

## 5. Materialization blocks

### FILE: `microsoft_playwright_browser_gate/tests/notification-status.spec.mjs`

```yaml
block_id: "MICROSOFT-PLAYWRIGHT-BROWSER-GATE:v278:1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "43c1dd61f0df5ddd3e6c1a06c36104170256da54e9d0d4243a0222b771f00c4a"
variables: []
secrets_allowed: false
```

````javascript
import { test, expect } from '@playwright/test';
import { createHash } from 'node:crypto';
import { createRequire } from 'node:module';
import { resolve } from 'node:path';
import { pathToFileURL } from 'node:url';

test('notification status connects scoped history, recovery and help without sends', async ({ page, context }) => {
  test.skip(process.env.ELITE_NOTIFICATION_BROWSER_E2E !== '1', 'Requires isolated Go/PG/TLS and synthetic signed identity');
  const require = createRequire(resolve(process.env.ELITE_WEB_ROOT, 'package.json'));
  const { EncryptJWT } = await import(pathToFileURL(require.resolve('jose')).href);
  const key = createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest();
  const base = process.env.ELITE_BASE_URL;
  expect(new URL(base).protocol).toBe('https:'); expect(new URL(base).hostname).toBe('127.0.0.1');
  async function session(permissions) {
    const value = await new EncryptJWT({ subject: 'fixture-operator', tenantId: process.env.ELITE_NOTIFICATION_TENANT, permissions, organizations: [process.env.ELITE_NOTIFICATION_ORGANIZATION], accessToken: process.env.ELITE_NOTIFICATION_TOKEN })
      .setProtectedHeader({ alg: 'dir', enc: 'A256GCM', typ: 'JWT' }).setIssuedAt().setExpirationTime('300s').encrypt(key);
    await context.addCookies([{ name: '__Host-elite_session', value, url: base+'/', httpOnly: true, secure: true, sameSite: 'Lax' }]);
  }
  await session(['appointment:manage']);
  let reads=0, writes=0;
  page.on('request',r=>{if(r.url().includes('/api/enterprise/franchise/notifications'))reads++;if(r.method()==='POST' && r.url().includes('/api/enterprise/'))writes++;});
  await page.goto(base+'/franchise?day='+process.env.ELITE_NOTIFICATION_DAY);
  const panel=page.getByRole('region',{name:'Estado de notificaciones WhatsApp'});
  await expect(panel).toBeVisible();
  await expect(panel.getByRole('status', { name: 'Resultado de consulta de WhatsApp', exact: true })).toBeVisible();
  await expect(panel.locator('input')).toHaveCount(0); expect(reads).toBe(0);
  const button=panel.getByRole('button',{name:'Consultar estado de WhatsApp',exact:true});
  await expect(button).toBeEnabled();
  // Lose only a completed real read response; retry is GET-only, no new send.
  await page.route('**/api/enterprise/franchise/notifications?**',async route=>{const response=await route.fetch();expect(response.status()).toBe(200);await route.abort('failed');},{times:1});
  await button.focus();await page.keyboard.press('Enter');
  await expect(panel.getByRole('status')).toContainText('No pudimos verificar');
  await button.click();
  await expect(panel.getByRole('heading',{name:'El proveedor informó entrega',exact:true})).toBeVisible();
  await expect(panel.getByText('Fecha del aviso del proveedor:', { exact: false })).toBeVisible();
  await expect(panel.locator('time').nth(1)).toHaveAttribute('datetime', '2020-10-19T05:45:14.000Z');
  await expect(panel.getByText('Este aviso no acredita una venta ni la aceptación del turno por el cliente. No vuelve a enviar mensajes.',{exact:true})).toBeVisible();
  await panel.getByText('Ayuda y referencia para soporte',{exact:true}).click();
  await expect(panel.getByText(process.env.ELITE_NOTIFICATION_EVENT,{exact:true})).toBeVisible();
  await expect(panel.getByText('whatsapp-status-view/1.0.0',{exact:true})).toBeVisible();
  expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth)).toBe(true);
  await page.screenshot({path:test.info().outputPath('notification-status-help.png'),fullPage:true});
  // A foreign response must replace success with an error, never stale success.
  await page.route('**/api/enterprise/franchise/notifications?**',async route=>{const response=await route.fetch();const body=await response.json();body.organization_id='foreign';await route.fulfill({response,json:body});},{times:1});
  await button.click();await expect(panel.getByRole('status')).toContainText('No pudimos verificar');
  await expect(panel.getByRole('heading',{name:'El proveedor informó entrega',exact:true})).toHaveCount(0);
  await button.click();await expect(panel.getByRole('heading',{name:'El proveedor informó entrega',exact:true})).toBeVisible();
  const query=new URLSearchParams({organizationId:process.env.ELITE_NOTIFICATION_ORGANIZATION,appointmentId:process.env.ELITE_NOTIFICATION_APPOINTMENT});
  const response=await page.request.get(base+'/api/enterprise/franchise/notifications?'+query);expect(response.status()).toBe(200);expect(response.headers()['cache-control']).toBe('no-store');
  const text=await response.text();for(const forbidden of ['5491112345678','wamid.synthetic','test-app-secret',process.env.ELITE_NOTIFICATION_TOKEN])expect(text).not.toContain(forbidden);
  await session([]);
  expect((await page.request.get(base+'/api/enterprise/franchise/notifications?'+query)).status()).toBe(403);
  await context.clearCookies();expect((await page.request.get(base+'/api/enterprise/franchise/notifications?'+query)).status()).toBe(401);
  expect(writes).toBe(0);
});
````


### FILE: `microsoft_playwright_browser_gate/package.json`
```yaml
block_id: "MICROSOFT-PLAYWRIGHT-BROWSER-GATE:package-json:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite isolated dependency module"
license: "LicenseRef-Workspace-Owner"
sha256: "eabb3e705b59759748724b2be3e8720184d26b570ab51e6c8a2ee8a0bdf41b90"
variables: []
secrets_allowed: false
```
````json
{
  "name": "elite-microsoft-playwright-browser-gate",
  "version": "0.1.0",
  "private": true,
  "type": "module",
  "packageManager": "pnpm@11.25.0",
  "engines": {
    "node": ">=20"
  },
  "scripts": {
    "test:runtime": "playwright test tests/runtime.spec.mjs",
    "test:enterprise-web": "playwright test tests/enterprise-web.spec.mjs tests/nonce-csp.spec.mjs",
    "test:help": "playwright test tests/journey-help.spec.mjs",
    "test:workspace": "playwright test tests/role-workspace.spec.mjs",
    "test:admin-reads": "playwright test tests/role-workspace.spec.mjs -g \"admin reads preserve independent grants\""
  },
  "devDependencies": {
    "@playwright/test": "1.62.1"
  }
}
````

### FILE: `microsoft_playwright_browser_gate/pnpm-lock.yaml`
```yaml
block_id: "MICROSOFT-PLAYWRIGHT-BROWSER-GATE:pnpm-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "pnpm 11.19.0 exact resolution of Microsoft Playwright 1.62.1"
license: "LicenseRef-Workspace-Owner"
sha256: "63eec1e3d5bac29965e850bf751269c7a8b317929b570205a1e223544f80a470"
variables: []
secrets_allowed: false
```
````yaml
lockfileVersion: '9.0'

settings:
  autoInstallPeers: true
  excludeLinksFromLockfile: false

importers:

  .:
    devDependencies:
      '@playwright/test':
        specifier: 1.62.1
        version: 1.62.1

packages:

  '@playwright/test@1.62.1':
    resolution: {integrity: sha512-DTcUc8qii+cpHvtOwggMtBRMjKZHXYWdw8syRYu2vtzuq4Wxphqq4NfCs5Zt44L6mA8rfDfj+PHnxFc/FeK6mQ==}
    engines: {node: '>=20'}
    hasBin: true

  fsevents@2.3.2:
    resolution: {integrity: sha512-xiqMQR4xAeHTuB9uWm+fFRcIOgKBMiOBP+eXiyT7jsgVCq1bkVygt00oASowB7EdtpOHaaPgKt812P9ab+DDKA==}
    engines: {node: ^8.16.0 || ^10.6.0 || >=11.0.0}
    os: [darwin]

  playwright-core@1.62.1:
    resolution: {integrity: sha512-wPYSwEBJY9GHraISXqyqtx0na0LpO3XEX7jNDhntbex7tzUS7kLnZsOlFruFJB4Hi/rhDMjXGqHewDZ68nYZVw==}
    engines: {node: '>=20'}
    hasBin: true

  playwright@1.62.1:
    resolution: {integrity: sha512-0M+L3LAD8/nm554LOla9Ayx0j0tmFZ0FBcoQ7F1VuVHpM/XpiC8RcDzBQB8W5+hA8L22THxELzeF+2WcUzvcLg==}
    engines: {node: '>=20'}
    hasBin: true

snapshots:

  '@playwright/test@1.62.1':
    dependencies:
      playwright: 1.62.1

  fsevents@2.3.2:
    optional: true

  playwright-core@1.62.1: {}

  playwright@1.62.1:
    dependencies:
      playwright-core: 1.62.1
    optionalDependencies:
      fsevents: 2.3.2
````

### FILE: `microsoft_playwright_browser_gate/source-lock.json`
```yaml
block_id: "MICROSOFT-PLAYWRIGHT-BROWSER-GATE:source-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite receipt from Microsoft GitHub release and npm registry artifacts"
license: "LicenseRef-Workspace-Owner"
sha256: "e55e8f9d756a65d9c67b30572f419b0776c103120a75e4b74ca5d365128db5ce"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-microsoft-playwright-browser-gate-lock/v1",
  "verifiedAt": "2026-08-29",
  "source": {
    "owner": "Microsoft",
    "repository": "microsoft/playwright",
    "release": "v1.62.1",
    "commit": "26a9e470a7b3c7822084b09fb7f13902c5f37b51",
    "commitSignatureVerified": true,
    "archiveUrl": "https://github.com/microsoft/playwright/archive/26a9e470a7b3c7822084b09fb7f13902c5f37b51.zip",
    "archiveBytes": 42927948,
    "archiveSha256": "29377110a042e60e8c8192eae145e8e8554b58b31073bdcdf9e25b47d019b186",
    "license": "Apache-2.0",
    "licensePath": "LICENSE",
    "licenseBytes": 11601,
    "licenseSha256": "45873d00a0dd243596deb4aa23b2493b3d1f0671921bf2538ea431d7380220eb",
    "packagedLicenseNormalization": "CRLF_TO_LF_ONLY",
    "packagedLicenseBytes": 11399,
    "packagedLicenseSha256": "7fab1461b41970ff376f1c9303a637076bfaaeb71cd12dd3a1c44aaf59a1a2b9"
  },
  "npm": {
    "package": "@playwright/test",
    "version": "1.62.1",
    "tarball": "https://registry.npmjs.org/@playwright/test/-/test-1.62.1.tgz",
    "tarballBytes": 8750,
    "tarballSha256": "009534220efd98c0361d8c4ee7e3db1ed510ab88a23e98b5081ef1c8fed64965",
    "integrity": "sha512-DTcUc8qii+cpHvtOwggMtBRMjKZHXYWdw8syRYu2vtzuq4Wxphqq4NfCs5Zt44L6mA8rfDfj+PHnxFc/FeK6mQ==",
    "node": ">=20",
    "pnpm": "11.25.0",
    "graphPackages": 4,
    "datedOsvFindings": 0
  },
  "browsers": [
    { "name": "Chromium", "revision": "1234", "version": "151.0.7922.34", "projects": 2 },
    { "name": "Firefox", "revision": "1538", "version": "153.0", "projects": 1 },
    { "name": "WebKit", "revision": "2336", "version": "26.5", "projects": 1 }
  ],
  "admission": {
    "runtimeSmoke": "4/4 PASS on Windows x86-64",
    "targetClaim": "public home semantics, security headers and responsive overflow only",
    "notProven": [
      "authenticated role journeys",
      "backend catalog or lead effects",
      "assistive technology interoperability",
      "performance, load, production edge, canary or rollback"
    ]
  }
}
````

### FILE: `microsoft_playwright_browser_gate/LICENSE.playwright.txt`
```yaml
block_id: "MICROSOFT-PLAYWRIGHT-BROWSER-GATE:upstream-license:v1"
operation: CREATE
provenance: ADAPTED
source: "https://github.com/microsoft/playwright/blob/26a9e470a7b3c7822084b09fb7f13902c5f37b51/LICENSE; CRLF_TO_LF_ONLY"
license: "Apache-2.0"
sha256: "7fab1461b41970ff376f1c9303a637076bfaaeb71cd12dd3a1c44aaf59a1a2b9"
variables: []
secrets_allowed: false
```
````text
                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "[]"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Portions Copyright (c) Microsoft Corporation.
   Portions Copyright 2017 Google Inc.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
````

### FILE: `microsoft_playwright_browser_gate/playwright.config.mjs`
```yaml
block_id: "MICROSOFT-PLAYWRIGHT-BROWSER-GATE:config:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite bounded configuration using official Microsoft Playwright API"
license: "LicenseRef-Workspace-Owner"
sha256: "23abf8c48048c142c4fc232d496444d526903e711aff8bef1ab1f8ba17b8e621"
variables: []
secrets_allowed: false
```
````javascript
import { defineConfig, devices } from '@playwright/test';
import { resolve } from 'node:path';

const webRoot = resolve(process.env.ELITE_WEB_ROOT || '..');
const suppliedBaseURL = process.env.ELITE_BASE_URL;
const runtimeOnly = process.env.ELITE_RUNTIME_ONLY === '1';
const baseURL = suppliedBaseURL || 'http://127.0.0.1:4173';
const parsed = new URL(baseURL);
const loopback = parsed.hostname === '127.0.0.1' || parsed.hostname === 'localhost' || parsed.hostname === '::1';
if (parsed.username || parsed.password || (parsed.protocol !== 'https:' && !(parsed.protocol === 'http:' && loopback)))
  throw new Error('ELITE_BASE_URL must be HTTPS or loopback HTTP and contain no credentials');

export default defineConfig({
  testDir: './tests',
  fullyParallel: false,
  forbidOnly: true,
  retries: 0,
  workers: 1,
  reporter: [['list']],
  outputDir: 'test-results',
  use: {
    baseURL,
    // Self-signed certificate only in the explicit loopback operator fixture.
    ignoreHTTPSErrors: loopback && (process.env.ELITE_OPERATOR_AGENDA_E2E === '1' || process.env.ELITE_NOTIFICATION_BROWSER_E2E === '1' || process.env.ELITE_QUOTE_E2E === '1' || process.env.ELITE_HELP_E2E === '1' || ['enabled','disabled'].includes(process.env.ELITE_WORKSPACE_E2E)),
    trace: 'retain-on-failure',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
  },
  webServer: runtimeOnly || suppliedBaseURL ? undefined : {
    command: 'pnpm exec next start --hostname 127.0.0.1 --port 4173',
    cwd: webRoot,
    url: baseURL,
    reuseExistingServer: false,
    timeout: 120_000,
  },
  projects: [
    { name: 'chromium-desktop', use: { ...devices['Desktop Chrome'] } },
    { name: 'chromium-mobile', use: { ...devices['Pixel 7'] } },
    { name: 'firefox-desktop', use: { ...devices['Desktop Firefox'] } },
    { name: 'webkit-desktop', use: { ...devices['Desktop Safari'] } },
  ],
});
````

### FILE: `microsoft_playwright_browser_gate/tests/runtime.spec.mjs`
```yaml
block_id: "MICROSOFT-PLAYWRIGHT-BROWSER-GATE:runtime-test:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite runtime smoke using official Microsoft Playwright API"
license: "LicenseRef-Workspace-Owner"
sha256: "67b87e6498145e481440924dfe11bad5baa6ce203fcc4771e567fdbd85e1784b"
variables: []
secrets_allowed: false
```
````javascript
import { test, expect } from '@playwright/test';

test('Microsoft Playwright executes isolated semantic browser assertions', async ({ page }) => {
  await page.setContent(`<!doctype html><html lang="en"><head><title>Elite Browser Gate</title></head>
    <body><main><h1>Browser gate</h1><form><label>Email <input type="email" name="email"></label>
    <button type="submit">Send request</button></form></main></body></html>`);
  await expect(page).toHaveTitle('Elite Browser Gate');
  await expect(page.getByRole('heading', { name: 'Browser gate' })).toBeVisible();
  await page.getByRole('textbox', { name: 'Email' }).fill('browser-gate@example.invalid');
  await expect(page.getByRole('button', { name: 'Send request' })).toBeEnabled();
});
````

### FILE: `microsoft_playwright_browser_gate/tests/enterprise-web.spec.mjs`
```yaml
block_id: "MICROSOFT-PLAYWRIGHT-BROWSER-GATE:enterprise-web-test:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite current enterprise-web public-home contract"
license: "LicenseRef-Workspace-Owner"
sha256: "bb5e5803c09f0f412c1db93fcae52b0f8e3bc877c5e2016ef662c7c33410145d"
variables: []
secrets_allowed: false
```
````javascript
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
````

### FILE: `microsoft_playwright_browser_gate/tests/nonce-csp.spec.mjs`
```yaml
block_id: "MICROSOFT-PLAYWRIGHT-BROWSER-GATE:nonce-csp-test:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite production response nonce and script propagation regression"
license: "LicenseRef-Workspace-Owner"
sha256: "a32564ed2fc449e9d2d413615ca9c6e62ea6e36415a0071ca756761c82c78ccc"
variables: []
secrets_allowed: false
```
````javascript
import { test, expect } from '@playwright/test';

function nonceFrom(policy) {
  const match = /'nonce-([^']+)'/.exec(policy || '');
  if (!match) throw new Error('CSP nonce is missing');
  return match[1];
}

test('strict CSP has a fresh nonce and every rendered script carries it', async ({ page, browserName }) => {
  if (browserName === 'webkit') {
    await page.route(/^https:\/\/127\.0\.0\.1:4173\//, async route => {
      const response = await route.fetch({ url: route.request().url().replace(/^https:/, 'http:') });
      await route.fulfill({ response });
    });
  }
  const first = await page.goto('/', { waitUntil: 'networkidle' });
  expect(first?.status()).toBe(200);
  const firstPolicy = first?.headers()['content-security-policy'];
  expect(firstPolicy).not.toContain("'unsafe-inline'");
  expect(firstPolicy).not.toContain("'unsafe-eval'");
  expect(firstPolicy).toContain("'strict-dynamic'");
  const firstNonce = nonceFrom(firstPolicy);
  expect(firstNonce).toMatch(/^[A-Za-z0-9+/]+=*$/);
  expect(firstNonce.length).toBeGreaterThanOrEqual(16);
  const scriptNonces = await page.locator('script').evaluateAll(nodes => nodes.map(node => node.nonce));
  expect(scriptNonces.length).toBeGreaterThan(0);
  expect(scriptNonces.every(value => value === firstNonce)).toBe(true);
  const second = await page.reload({ waitUntil: 'networkidle' });
  const secondNonce = nonceFrom(second?.headers()['content-security-policy']);
  expect(secondNonce).not.toBe(firstNonce);
  const secondScriptNonces = await page.locator('script').evaluateAll(nodes => nodes.map(node => node.nonce));
  expect(secondScriptNonces.every(value => value === secondNonce)).toBe(true);
});
````

### FILE: `microsoft_playwright_browser_gate/verify_contract.ps1`
```yaml
block_id: "MICROSOFT-PLAYWRIGHT-BROWSER-GATE:verify:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite source, lock, runtime and bounded target verification"
license: "LicenseRef-Workspace-Owner"
sha256: "7a53ced1799ae27cdef8ea91c86435c8d4b2a05652d8a28ec2aff6cb49c998de"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0

[CmdletBinding()]
param(
  [string]$ComponentRoot = $PSScriptRoot,
  [switch]$RunEnterpriseWeb,
  [string]$WebRoot
)

$ErrorActionPreference = 'Stop'
$root = (Resolve-Path -LiteralPath $ComponentRoot).Path
$lock = Get-Content -Raw -LiteralPath (Join-Path $root 'source-lock.json') | ConvertFrom-Json -Depth 30
$package = Get-Content -Raw -LiteralPath (Join-Path $root 'package.json') | ConvertFrom-Json -Depth 20
if ($lock.source.repository -cne 'microsoft/playwright' -or $lock.source.release -cne 'v1.62.1' -or $lock.source.commit -cne '26a9e470a7b3c7822084b09fb7f13902c5f37b51' -or $lock.source.commitSignatureVerified -ne $true) { throw 'Microsoft source identity drifted' }
if ($lock.npm.package -cne '@playwright/test' -or $lock.npm.version -cne '1.62.1' -or $lock.npm.graphPackages -ne 4 -or $lock.npm.datedOsvFindings -ne 0) { throw 'Playwright artifact/runtime evidence drifted' }
$expectedBrowsers = @(
  @{ name='Chromium'; revision='1234'; version='151.0.7922.34'; projects=2 },
  @{ name='Firefox'; revision='1538'; version='153.0'; projects=1 },
  @{ name='WebKit'; revision='2336'; version='26.5'; projects=1 }
)
if (@($lock.browsers).Count -ne 3) { throw 'Playwright browser inventory drifted' }
for ($index = 0; $index -lt $expectedBrowsers.Count; $index++) {
  $actual = $lock.browsers[$index]; $expected = $expectedBrowsers[$index]
  if ($actual.name -cne $expected.name -or $actual.revision -cne $expected.revision -or $actual.version -cne $expected.version -or $actual.projects -ne $expected.projects) { throw "Playwright browser identity drifted: $($expected.name)" }
}
if ($package.devDependencies.'@playwright/test' -cne '1.62.1' -or $package.packageManager -cne 'pnpm@11.25.0') { throw 'Package pins drifted' }
$license = Join-Path $root 'LICENSE.playwright.txt'
if ((Get-Item -LiteralPath $license).Length -ne 11399 -or (Get-FileHash -LiteralPath $license -Algorithm SHA256).Hash.ToLowerInvariant() -cne '7fab1461b41970ff376f1c9303a637076bfaaeb71cd12dd3a1c44aaf59a1a2b9' -or $lock.source.packagedLicenseNormalization -cne 'CRLF_TO_LF_ONLY') { throw 'Packaged Microsoft license normalization drifted' }
$lockHash = (Get-FileHash -LiteralPath (Join-Path $root 'pnpm-lock.yaml') -Algorithm SHA256).Hash.ToLowerInvariant()
if ($lockHash -cne '63eec1e3d5bac29965e850bf751269c7a8b317929b570205a1e223544f80a470') { throw 'pnpm lock drifted' }
if (-not (Test-Path -LiteralPath (Join-Path $root 'node_modules') -PathType Container)) { throw 'Install with pnpm --ignore-workspace --frozen-lockfile before runtime verification' }
Push-Location $root
try {
  $version = (& pnpm --ignore-workspace exec playwright --version | Out-String).Trim()
  if ($LASTEXITCODE -ne 0 -or $version -cne 'Version 1.62.1') { throw "Unexpected Playwright runtime: $version" }
  $env:ELITE_RUNTIME_ONLY = '1'
  & pnpm --ignore-workspace exec playwright test tests/runtime.spec.mjs
  Remove-Item Env:ELITE_RUNTIME_ONLY -ErrorAction SilentlyContinue
  if ($LASTEXITCODE -ne 0) { throw 'Microsoft Playwright runtime smoke failed' }
  if ($RunEnterpriseWeb) {
    if ([string]::IsNullOrWhiteSpace($WebRoot)) { throw 'WebRoot is required with RunEnterpriseWeb' }
    $resolvedWeb = (Resolve-Path -LiteralPath $WebRoot).Path
    foreach ($required in @('package.json','.next/BUILD_ID','config/business.example.json')) {
      if (-not (Test-Path -LiteralPath (Join-Path $resolvedWeb $required) -PathType Leaf)) { throw "Enterprise web target missing: $required" }
    }
    $env:ELITE_WEB_ROOT = $resolvedWeb
    Remove-Item Env:ELITE_BASE_URL -ErrorAction SilentlyContinue
    & pnpm --ignore-workspace exec playwright test tests/enterprise-web.spec.mjs tests/nonce-csp.spec.mjs
    if ($LASTEXITCODE -ne 0) { throw 'Enterprise web browser gate failed' }
  }
} finally {
  Pop-Location
}
'MICROSOFT_PLAYWRIGHT_BROWSER_GATE_PASS runtime=4 target=' + ($(if ($RunEnterpriseWeb) {'8'} else {'SKIPPED'}))
````

### FILE: `microsoft_playwright_browser_gate/README.md`
```yaml
block_id: "MICROSOFT-PLAYWRIGHT-BROWSER-GATE:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite usage and non-claim contract"
license: "LicenseRef-Workspace-Owner"
sha256: "2d5b65ac7746ac5b964c355f9fc2200f87d5caade476be1295c02602abd1ba97"
variables: []
secrets_allowed: false
```
````markdown
# Microsoft Playwright browser gate

## Opt-in notification history and recovery gate (V278)

Build the composed Next application first. In this isolated gate directory run
pnpm install --ignore-workspace --frozen-lockfile (add --offline only if the exact
locked artifacts are already cached). Installing from the parent workspace does
not prove that this separate CLI is available.

Set ELITE_NOTIFICATION_BROWSER_E2E=1, ELITE_WEB_ROOT to the absolute built
composition, ELITE_WHATSAPP_PYTHON to Python and
ELITE_WHATSAPP_TEST_DATABASE_URL to the dedicated loopback test database.
From the composition run:
go test ./internal/whatsappbridge -run '^TestNotificationStatusBrowserPostgres$' -v -count=1 -timeout=5m

Four browser projects read the actual Next BFF / Go / PostgreSQL status history,
using a uniquely created synthetic business configuration that enables the
existing features.whatsapp_status_history flag. The default profile keeps the
function off; BFF tests prove missing/false/config-failure without backend calls.
lose a read response, recover explicitly and reject a foreign-scope projection.
Database postconditions retain one completed status job, one observation and one
outbound attempt. The fixture pre-seeds an accepted outbound and runs the actual
status worker; it does NOT demonstrate a browser send to Meta.
Self-signed TLS relaxation is confined to this explicit loopback fixture.
Each project retains its own output directory and help-view screenshot.
No automatic polling or POST is permitted by the test.
Configuration, harness and scenarios are AUTHORED, not Microsoft source.

The current section above governs the V278 local claim. The versioned sections
below retain their historical scope; earlier host/UI-pending statements are
superseded only for this tested local integration, never for target deployment.

## Operator agenda, V263 (opt-in)

Compose the franchise profile, install the web and this gate's exact locks and build the web first. Apply all 49 migrations to a **disposable** PostgreSQL database on 127.0.0.1 named `elite_confirmation_*`. With `TEST_DATABASE_URL`, absolute `ELITE_WEB_ROOT` and `ELITE_CONFIRMATION_E2E=1`, run:

```powershell
go test ./internal/platform/httpapi -run '^TestAppointment(AgendaBrowserPostgres|ConfirmationBFFPostgres|AgendaAuthorizationAndRange)$' -v -count=1 -timeout=3m
```

The Go harness launches the actual built Next server on port 4173 (no existing server may own it), a loopback TLS reverse proxy with an ephemeral self-signed certificate, real API/PostgreSQL and RS256/JWKS identity verification. Browser sessions are synthetic JWE fixtures, not completed OIDC logins. Certificate errors are ignored **only** for the explicitly selected loopback operator fixture; this is not deployed TLS/certificate evidence. The earlier HTTPS-to-HTTP browser routing shim is not used for this gate.

Each browser has its own tenant. Chromium desktop proves normal confirmation; mobile Chromium, Firefox and WebKit lose a response **after the real commit** and reload the durable state without a second confirmation POST. All select resources by name, use server-owned versions, display confirmed state after reload and check horizontal overflow. The gate preserves Secure/HttpOnly cookies, verifies configured APP_BASE_URL origin and never replaces successful API responses with mocks. Audit rows remain in the dedicated database until explicit disposal; never disable immutable-audit triggers. Runtime/login/credentials/PKI, notifications and project production gates remain separate.

## Connected public lead gate (opt-in)

### V261: continue to a durable appointment request

With `GO-FRANCHISE-CUSTOMER-JOURNEY-API` 0.9.1, `GO-ELECTROMOBILITY-PUBLIC-CRM-API` 0.2.3 and `TS-FRANCHISE-JOURNEY-PORTALS` 0.9.1, the same disposable harness also runs `TestPublicAppointmentBrowserPostgres`. It reuses lead capture, then follows the actual appointment link, reads published location/availability and submits a capacity-one slot. It deliberately drops the browser response **after a real commit**, retries with the same key, checks the same receipt and prevents another submission after success. Divergent replay, full capacity, missing key and browser-supplied state are rejected; refresh removes the full slot. SQL independently checks four requested appointments joined to their leads/slots, four events, four completed receipts and no capacity mismatch.

Run both connected gates from the backend after the prerequisites below:

```powershell
$env:ELITE_PUBLIC_LEAD_E2E = '1'
try {
  go test ./internal/platform/httpapi -run '^TestPublic(Lead|Appointment)BrowserPostgres$' -v -count=1 -timeout=5m
  if ($LASTEXITCODE -ne 0) { throw 'Connected journey gate failed' }
} finally {
  Remove-Item Env:ELITE_PUBLIC_LEAD_E2E -ErrorAction SilentlyContinue
}
```

The harness is shared, not a second backend. Each test owns a random tenant, synthetic work window and disposable database rows. It fails fast on the first browser failure, keeps retries at zero and never counts unexecuted projects as PASS. The UI disables interactions before hydration following Microsoft's guidance: https://playwright.dev/docs/navigations#hydration. Method for recovery/idempotency: https://aws.amazon.com/builders-library/making-retries-safe-with-idempotent-APIs/ (checked 2026-09-06). Runtime/dependency pins are unchanged; glue, tests and corrections are `AUTHORED`.

The durable state is **requested**, not **confirmed** by an authenticated operator. This does not prove operator/resource assignment, business confirmation, notifications, sale, concurrent load, antiabuse, legal policy or production. The key recovery demonstrated is within the mounted form, not across browser restart. A no-opt-in home smoke skips both connected cases per browser explicitly.

The canonical Go test `TestPublicLeadBrowserPostgres` starts the real public Go handler/repository, injects its loopback address into the production-built BFF and runs four real browser projects. It checks catalog visibility, form submission, receipt and appointment link, identical replay, conflicting replay, missing consent and missing idempotency key. Afterwards it queries PostgreSQL for exactly four leads, four consents, four outbox events and four completed idempotency records. No catalog/lead response is mocked. WebKit retains the existing loopback-only transport shim; this does not prove deployment HTTPS.

Compose the current backend and web profiles, apply their migrations to a **dedicated disposable** loopback database named `elite_browser_*`, install the frozen web/Playwright dependencies and build the web. From the backend root, with Go and Node on PATH:

```powershell
$env:TEST_DATABASE_URL = '<dedicated loopback PostgreSQL URL; database elite_browser_*>'
$env:ELITE_WEB_ROOT = '<absolute path to the built web composition>'
$env:ELITE_PUBLIC_LEAD_E2E = '1'
try {
  go test ./internal/platform/httpapi -run '^TestPublicLeadBrowserPostgres$' -v -count=1 -timeout=5m
  if ($LASTEXITCODE -ne 0) { throw 'Connected public lead gate failed' }
} finally {
  Remove-Item Env:ELITE_PUBLIC_LEAD_E2E -ErrorAction SilentlyContinue
}
```

The harness rejects non-loopback/project databases, creates a unique synthetic tenant, removes only that tenant's fixture rows, rejects every protected authentication attempt and passes no database URL to the browser runner. A missing opt-in skips this additional test during the historical home smoke; a skipped connected test is not evidence of connectivity. The opt-in path fails on missing dependencies or persistence evidence. No appointment reservation, operator response, authenticated journey, antiabuse, production consent policy, external provider or production readiness is claimed. Tests and harness are `AUTHORED`, using the locked Microsoft Playwright runtime and the admitted Go/PostgreSQL components, not code copied from Microsoft.

Method authority checked 2026-09-06: https://playwright.dev/docs/best-practices (user-visible assertions, isolation, controlled test data) and https://playwright.dev/docs/api-testing (API postconditions combined with browser actions).

This component pins the official Microsoft `@playwright/test` 1.62.1 package and Chromium revision 1234, Firefox revision 1538 and WebKit revision 2336. `LICENSE.playwright.txt` preserves the official license text with the Markdown system's declared CRLF→LF-only normalization; both upstream and packaged hashes are locked. Configuration, tests, lock receipt and PowerShell verification are Elite `AUTHORED`; they are not presented as Microsoft source.

Install and verify the isolated Microsoft runtime:

```powershell
pnpm install --ignore-workspace --frozen-lockfile
pwsh -NoProfile -File ./verify_contract.ps1
```

After composing, installing and building `ENTERPRISE_WEB_PACK_PLAN.md`, run the bounded public-home gate:

```powershell
pwsh -NoProfile -File ./verify_contract.ps1 -RunEnterpriseWeb -WebRoot ../enterprise-web
```

The default public-home gate executes Chromium desktop/mobile, Firefox desktop and WebKit desktop projects. Four tests check one semantic H1, the catalog link, locale, security headers, horizontal overflow and browser/console errors; four more require a fresh production CSP nonce per response and prove that every rendered script carries exactly its response nonce. That default smoke does not claim authenticated roles, backend effects, accessibility with real assistive technology, performance, production edge/CDN behavior, canary or rollback. Those remain project gates. Remote execution accepts only HTTPS through `ELITE_BASE_URL`; local HTTP is loopback-only and credentials in URLs are rejected.

## Frozen offline installation with pnpm 11.25.0

Use the exact admitted pnpm artifact and observed Node runtime from the project tool lock. An offline store needs both package bytes and the registry metadata used by supply-chain policy checks. ERR_PNPM_NO_OFFLINE_META is a failed gate even if node_modules was linked. In an authorized public-registry preparation step, run the same frozen install online to populate metadata, then repeat --offline --frozen-lockfile --ignore-workspace. Preserve the lock hash and both receipts; never disable trust or release-age policy to force PASS. V321 verified this sequence with unchanged dependency lockfiles.
````

### FILE: `microsoft_playwright_browser_gate/tests/journey-help.spec.mjs`
```yaml
block_id: "MICROSOFT-PLAYWRIGHT-BROWSER-GATE:HELP:microsoft_playwright_browser_gate-tests-journey-help.spec.mjs:v1"
operation: CREATE
provenance: AUTHORED
source: "Local read-only integration of existing public operational guides and admitted session/runtime; see VERSIONED_JOURNEY_HELP_V379.md"
license: "LicenseRef-Workspace-Owner"
sha256: "8570942d1f0a8b0dae136e9f70d8e89cfd4e246309a9420efc5859dec1397436"
variables: []
secrets_allowed: false
```
````javascript
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
````

### FILE: `microsoft_playwright_browser_gate/fixtures/help-loopback-proxy.go`
```yaml
block_id: "MICROSOFT-PLAYWRIGHT-BROWSER-GATE:HELP:microsoft_playwright_browser_gate-fixtures-help-loopback-proxy.go:v1"
operation: CREATE
provenance: AUTHORED
source: "Local read-only integration of existing public operational guides and admitted session/runtime; see VERSIONED_JOURNEY_HELP_V379.md"
license: "LicenseRef-Workspace-Owner"
sha256: "54d0954b628054fabf6c6e5f82329bfd3ae5056269ffc032e96c3ef18c1c07e5"
variables: []
secrets_allowed: false
```
````go
// AUTHORED test-only TLS bridge for production Next/JWE help tests. Never deploy.
package main

import (
	"fmt"
	"io"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 { panic("expected explicit loopback Next URL") }
	u, err := url.Parse(os.Args[1])
	if err != nil || u.Scheme != "http" || u.Hostname() != "127.0.0.1" || u.Port() == "" || u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" { panic("only loopback HTTP target allowed") }
	server := httptest.NewTLSServer(httputil.NewSingleHostReverseProxy(u))
	defer server.Close()
	fmt.Println(server.URL)
	done := make(chan struct{})
	go func() { _, _ = io.Copy(io.Discard, os.Stdin); close(done) }()
	select { case <-done: case <-time.After(10 * time.Minute): }
}
````

### FILE: `microsoft_playwright_browser_gate/tests/role-workspace.spec.mjs`
```yaml
block_id: "MICROSOFT-PLAYWRIGHT-BROWSER-GATE:WORKSPACE:microsoft_playwright_browser_gate-tests-role-workspace.spec.mjs:v1"
operation: CREATE
provenance: AUTHORED
source: "Local permission-filtered workspace integration; see ROLE_WORKSPACE_INTEGRATION_V381.md"
license: "LicenseRef-Workspace-Owner"
sha256: "07524a4d8c6f2c4c0fcab1c97eac90bbb2a09bea535d14cafdb2a074c0853b8d"
variables: []
secrets_allowed: false
```
````javascript
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
````

## 6. Configuration surface

- `ELITE_WEB_ROOT`: absolute or caller-resolved path to a built enterprise-web composition; defaults to the parent of the gate directory.
- `ELITE_BASE_URL`: optional HTTPS target. When absent, Playwright starts the local built target at loopback `127.0.0.1:4173`.
- `RunEnterpriseWeb`/`WebRoot`: explicit switch/path for the bounded target gate.

No credential, cookie, token, account or provider endpoint belongs in this pack. Authenticated state and secrets require a project-specific Playwright project and external secret handling after readiness approval.

## 7. Dependency bill

| Package/tool | Exact pin | Purpose | License | Evidence |
|---|---|---|---|---|
| Microsoft `@playwright/test` | 1.62.1 | browser test runner/API | Apache-2.0 | signed commit, archive/npm/license hashes |
| Playwright + Core | 1.62.1 | browser automation/runtime | Apache-2.0 | pnpm integrity lock |
| Chromium | 151.0.7922.34 rev 1234 | desktop/mobile smoke | upstream browser terms | Playwright `browsers.json`; binary not redistributed |
| Firefox | 153.0 rev 1538 | desktop smoke | upstream browser terms | Playwright `browsers.json`; binary not redistributed |
| WebKit | 26.5 rev 2336 | desktop smoke | upstream browser terms | Playwright `browsers.json`; binary not redistributed |
| Node.js | >=20 | runtime | project installation | version checked through package engine |
| pnpm | 11.25.0 | frozen install | MIT | exact packageManager/lock |

The pack redistributes only the exact Apache-2.0 license text, not Microsoft package or browser binaries. `pnpm install --ignore-workspace --frozen-lockfile` obtains the official package graph; `playwright install chromium firefox webkit` is an explicit network/capacity action when the fixed revisions are not already cached.

## 8. Apply order

1. Materialize the pack into an empty target and verify all ten hashes.
2. Run `pnpm install --ignore-workspace --frozen-lockfile`; this keeps the gate isolated even when composed below the web pnpm workspace. Install Chromium 1234, Firefox 1538 and WebKit 2336 only if absent and network/capacity are approved.
3. Run `verify_contract.ps1` for the isolated 2-project runtime smoke.
4. Compose/install/typecheck/test/build `ENTERPRISE_WEB_PACK_PLAN.md` in a separate directory.
5. Run `verify_contract.ps1 -RunEnterpriseWeb -WebRoot <built-web>`; require eight target tests: four semantic/header/overflow and four nonce/CSP.
6. Add project-specific role/provider/accessibility/performance journeys only after readiness supplies their real authority and access.

Rollback removes only the separate gate module, its downloaded dependencies/browser cache if explicitly desired, and generated Playwright failure artifacts. It never modifies the application composition.

## 9. Verification

- source commit signature `verified=true`; GitHub archive 42,927,948 bytes/SHA-256 exact;
- npm tarball 8,750 bytes/SHA-256 and registry SHA-512 integrity exact;
- four-package pnpm graph; `pnpm audit` and dated official OSV batch: zero known findings;
- frozen offline reinstall PASS from the audited cache;
- Playwright version exactly 1.62.1; Chromium 151.0.7922.34 rev 1234, Firefox 153.0 rev 1538 and WebKit 26.5 rev 2336;
- runtime matrix: 4/4 PASS;
- freshly composed web: frozen install, typecheck, 14 unit tests and production build PASS;
- public-home Chromium desktop/mobile + Firefox/WebKit desktop: semantic/header/overflow 4/4 and fresh nonce/exact script propagation 4/4 PASS.

These results prove the bounded claim only. Authenticated roles, real backend/provider effects, AT accessibility, performance and production deployment remain conditioned.

## 10. Reconstruction evidence

V107 staging on Windows x86-64 used Node 24.14.1, pnpm 11.19.0 and the three official browser revisions. Runtime passed 4/4; the target passed semantic/header/overflow 4/4 plus fresh nonce and exact rendered-script propagation 4/4 after preserving CSP and handling only WebKit's loopback upgrade inside the local test route. No Microsoft bytes were modified. Evidence distinguishes that transport shim from production HTTPS and leaves CDN/WAF repetition conditioned.

## V321 — official package-manager security update

pnpm11.25.0 replaces vulnerable11.19.0 for this consumer; dependency lock bytes remain unchanged. Published bundle473packages/0advisories, registry ECDSA and publish/SLSA attestations verified, exact commit6d90c71efdffbc909b499490b64c66badc720327. Node24.20.0 standalone observed;115webPASS/1existingSKIP, build,92browser phases+agenda4,Playwright runtime4,Lighthouse policy5, local installation comparison. Eight changed materialized files across three owners; no upstream matcher or browser version changed. See reconstruction_evidence/PNPM_BUNDLE_SECURITY_V321.md. Native6files unchanged and transitive native SCA remains separate; product admission/monitor/target still blocked.

## Delta V330 — verificación multi-tab del flujo existente

Cambio AUTHORED sólo de pruebas; no amplía lógica de producto ni admisión.
ELITE_RETURN_MULTITAB_E2E=1 requiere fixtures ORDER/DELIVERY_READ/RETURN_RECOVERY.
Seis carreras entre páginas con cookie compartida y sessionStorage independiente;
200/409, efectos durables únicos y recuperación por digest. Informe
reconstruction_evidence/RETURN_MULTITAB_VERIFICATION_V330.md. Conserva condiciones
de target, providers, seguridad, operación y aceptación aún no demostradas.

## Delta V331 — fronteras de sesión de devolución

Dos archivos de pruebas AUTHORED; producto y dependencias intactos.
ELITE_RETURN_SESSION_E2E=1 requiere ELITE_RETURN_MULTITAB_E2E y sus fixtures.
Comprueba lector, cookie ausente, organización y segundo operador autorizado;
restaura referencia original sin reenviar. No prueba revocación IdP ni borrado
de información ya entregada al cliente. Evidencia RETURN_SESSION_BOUNDARIES_V331.md
en reconstruction_evidence; condiciones productivas permanecen abiertas.

V378 extends the same public lead→appointment browser test with independent expected text for `ELITE_PUBLIC_LOCALE=es-AR|en-US|fr-FR`; fr-FR proves explicit Spanish fallback. Configure matching BUSINESS_CONFIG_FILE and default-market timeZone (es-AR Buenos Aires, en-US/fr-FR UTC). Keep all prior consent, replay, divergent-intent, response-loss, capacity and SQL oracles. Date comparison equates only U+00A0/U+202F and U+0020, preserving all numbers/punctuation and the exact backend UTC instant. This does not suppress a browser lifecycle failure.

V379: test:help runs journey-help.spec.mjs only with ELITE_HELP_E2E=1, ELITE_WEB_ROOT, synthetic AUTH_SESSION_SECRET and ELITE_BASE_URL=https://127.0.0.1:<port>. Enable whatsapp_status_history only in the isolated fixture configuration. Compile fixtures/help-loopback-proxy.go with the admitted Go toolchain, run it against the owned production Next HTTP loopback URL and keep stdin open; it prints an ephemeral self-signed HTTPS URL and stops on stdin EOF or10minutes. Stop the owned Next process afterward. This test-only stdlib proxy must never be deployed. No database/provider is required for static guidance. Four browsers prove exact-version filtering, encrypted-cookie authorization, expiry/tampering, superseded reads, error recovery and no operation writes. Synthetic cookies do not certify an external IdP.

V380 extends test:help with an independent15profile/15guide access matrix and all15exact-version UI links in four browsers. Eight tests total when ELITE_HELP_E2E=1, no retries. Run with --timeout=120000 for the complete role matrix. The TLS fixture, JWE mechanism and dependency lock are unchanged; no provider/database required for public static guidance.

V381 integrates TS-MULTIROLE-ONBOARDING0.1.1 under opt-in role_workspace. BFF header filters exact effective session permissions and current feature/module flags using its canonical registry; guest customer entry remains. Browser test:workspace proves fifteen permission profiles across four browsers plus disabled flag, local307/session rejection, direct admin denial and help navigation. No domain handler, dependency lock or training policy changes. See ROLE_WORKSPACE_INTEGRATION_V381.md.

V382: exact selected-composition compatibility refreshed. Administrative lead reads are conditional on the existing lead grant; shared quote-format code corrects admin/customer minor-unit labels without repricing. TestAdministrativeReadBrowserPostgres uses real Go/JWKS/PostgreSQL and four browsers with synthetic identities; all production Go/SQL remains unchanged. The role implementation and quote acceptance/date/state logic retain exact source parity. See docs/private-portal-reads.md and reconstruction_evidence/PRIVATE_PORTAL_READS_V382.md. No full business, IdP, security, release or target admission.

V383: private read failures expose a fixed alert and explicit GET reconsultation in admin/customer/factory, without error details or command replay. Production build plus Go/JWKS/PostgreSQL and four browser projects prove persistent denial with the insufficient bearer and recovery after session correction; zero writes and unchanged durable snapshots.263web tests and canonical parity. Existing role/portal/Go business sources unchanged. See reconstruction_evidence/PRIVATE_READ_RECOVERY_V383.md; no full-control or target promotion.

V387: customer/factory cursor traversal connects three existing lists beyond25rows.81durable records per browser, four projects, reload/independent first-page recovery, zero writes and unchanged snapshots verified. Source is AUTHORED, business/API/SQL/dependencies unchanged. Supporting TEST02 progress, no whole-control or native security/release promotion. See reconstruction_evidence/PRIVATE_PORTAL_PAGINATION_V387.md.

V388: admin orders/cases/leads now traverse existing scoped keyset APIs using the shared navigator.276web tests and four browser projects exercise all six lists,163list records, no omissions/duplicates per list, permission removal, reload/first-page recovery and unchanged durable snapshots. All code remains AUTHORED/CONDITIONED; no full TEST02, native security or release promotion. See reconstruction_evidence/ADMIN_PORTAL_PAGINATION_V388.md.

V389: customer appointment times share trusted business locale/zone across SSR and interactive management, preserving original instants and cancellation behavior.282web tests plus eight connected browser/configuration runs, two appointment instants per view and seven-table unchanged snapshots. AUTHORED/CONDITIONED; no whole TEST02, security or release promotion. See reconstruction_evidence/CUSTOMER_APPOINTMENT_TIMEZONE_V389.md.

V390: customer cancellation receipt binding, synchronous single-submit and explicit GET recovery close the existing reference journey.294web tests;4actual cancellation browser projects,16unique durable cancellations/audits/outbox and28API writes;8read/time regression runs remain read-only. AUTHORED/CONDITIONED, no whole TEST02/security/release promotion. See reconstruction_evidence/CUSTOMER_CANCELLATION_RECOVERY_V390.md.

V393: the reference does not use runtime image transformation. Next image optimization is explicitly disabled and pnpm11.25.0 ignores the optional sharp dependency; frozen lock removes only its exclusive closure. This is dependency omission, not a repaired/admitted Sharp binary. Before enabling runtime image optimization, re-admit a compatible image pipeline, its exact dependencies, licenses/native provenance/security and web performance. Existing V386 vulnerability reproduction remains deferred. A connected real HTTP regression requires /_next/image to return404. No new version or business policy; AUTHORED/CONDITIONED. See reconstruction_evidence/OPTIONAL_IMAGE_DEPENDENCY_CONTAINMENT_V393.md.
