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
