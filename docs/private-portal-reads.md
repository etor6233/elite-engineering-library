# Private portal read views — V382

The administrative route requires admin:read. Overview, orders and service cases use that
grant; its optional opportunities section additionally requires lead:read (or exact wildcard).
The BFF does not request that endpoint without its permission. Roles and encrypted frontend
claims never enlarge the actual bearer grant enforced by Go. A session that overclaims lead
access while its bearer lacks it still fails closed; no implicit fallback or permission grant.

The first organization in the current verified session selects the existing query scope.
The real HTTP layer verifies tenant, permission and organization; customer reads also bind
the bearer subject. This change adds no organization switcher, business role or IdP lifecycle.

Admin orders and customer order/quote summaries now use the exact existing quote amount
algorithm shared under platform/i18n/money.ts. It checks safe nonnegative minor-unit integers
and supported currency codes, derives currency digits from the pinned Intl runtime, and uses
BigInt plus formatToParts to retain the remainder without floating-point division. Existing
Spanish presentation and invalid-amount wording are preserved. USD100minor units displays
USD1,00; JPY123 displaysJPY123; KWD1234 displaysKWD1,234. No price, FX, tax, acceptance,
business rounding or stored amount changes. The quote-page extraction is exactly reversible.

Qualification:11private-view regressions and13literal amount/invalid-input cases;260full web
tests pass with one pre-existing explicit connected skip, typecheck and build. The actual
Next/JWE to Go/RS256/JWKS to PostgreSQL fixture uses53unchanged migrations, five independent
HTTP negatives and four browser projects with retries0. It tests admin-only, combined/wildcard,
case-sensitive grants, other organization/tenant, non-admin denial, factory and customer reads,
customer quote display, permission removal on reload, backend rejection of overclaimed frontend
permissions and subsequent recovery. No browser/API writes; snapshots across six domain tables
are unchanged. The repository fixture is durable and synthetic, not an external IdP/provider.

Run Go TestAdministrativeReadBrowserPostgres with ELITE_ADMIN_READ_E2E=1, absolute ELITE_WEB_ROOT,
an owned loopback elite_confirmation_* TEST_DATABASE_URL, the composed web build and installed
admitted Playwright gate. Enable role_workspace/customer_portal/factory_portal/public_catalog
and the catalog/crm/procurement modules in the isolated business config. The fixture reuses the
selected journey test issuer/clock, starts its owned API/TLS/Next, runs test:admin-reads and
stops its processes. The caller owns migration/setup and PostgreSQL shutdown. Never use a
project database, real identities, production keys or the synthetic issuer as deployment.

Rollback must keep a compatible composition. Preserve the old source for diagnosis, but do not
promote the known admin permission coupling or the misleading minor-unit labels as a rollback.
On regression, block promotion and repair the candidate. Authorization errors, IdP revocation,
arbitrary analytics/KPIs, full operational workflows, monetary policies, accessibility audit,
native licence/security, release and target acceptance retain their own gates. This read boundary
does not close any of those by implication.

V383 adds fixed-text Next error boundaries for admin/customer/factory. They never read or display the thrown error, stack, digest or bearer. An ordinary anchor to the fixed portal root performs a new GET only after the user clicks; no reset hook, automatic retry or operation replay. The message does not claim success or failure of a previous write. The existing page rechecks the current session and the backend still enforces its bearer. A segment error can cover descendants too: recovery deliberately returns to the portal root, not an arbitrary failed URL. Root-layout failures and external IdP recovery remain separate.

Verified against the pinned Next16.3.4 local ErrorComponent contract, actual production build and Go/JWKS/PostgreSQL in four browser projects. For all three portal roots an overclaimed fixture session with an insufficient bearer fails without domain records, an unchanged reconsult still fails, then changing to a valid session and clicking the same link issues GET and restores scoped data. No browser/API writes and six durable-table snapshots unchanged.263web tests plus one inherited explicit skip, typecheck/build and canonical reconstruction pass. This proves the exercised permission-rejection recovery, not an injected database outage or complete workflow/security/release admission.

V387 connects the existing cursor API for customer orders/service cases and factory units. Fixed25row requests retain session organization and bearer; orders_after/cases_after are independent, units_after belongs to factory. Ordinary GET links preserve the other list, reload retains position and first-page recovery clears only the chosen cursor. Opaque cursor values are URL-encoded and bounded to512characters; duplicates/excessive values enter existing fixed error recovery. Unknown query inputs cannot replace organization or page size. Final/empty pages have no next link. No write, offset, count, cross-page snapshot, permission, currency or warranty-policy guarantee is added. Concurrent data changes retain the existing keyset semantics. Admin lists and the separate journey array projection are outside this correction.

The connected fixture holds27records per list:81unique records per browser, no omissions/duplicates across first/last pages; four projects also exercise independent cursor retention, reload, first-page navigation, existing scopes, fixed read recovery and unchanged six-table snapshot. No provider or production access. The qualification is independent of the V386 deferred native vulnerability research and cannot close TEST03/07.

V388 supersedes the V387 admin-list exclusion: admin orders/service cases/leads now use the same cursor navigator. Three independent query cursors retain the other granted lists. The lead cursor is read/preserved only with lead:read; revoking that permission removes its list and cursor from subsequent links.25row requests and session organization remain fixed. Four added tests and the extended actual browser/Go/PostgreSQL fixture traverse28admin orders,27cases and27leads (82rows across these lists), while all81customer/factory observations remain tested. Total163list records per browser, not163distinct underlying business entities. Reload/first-page recovery and no writes/snapshot preservation pass. No domain API, business rule, source graph or dependency change.

V389: customer account and appointment management render the same stored appointment instant using the admitted default locale and market timeZone, shared with public booking. Each time element retains the original ISO datetime and shows the explicit market zone. The interactive cancellation component receives the server-formatted label, so browser locale/timeZone cannot change it during hydration. Date validity is checked by the existing formatter. Six added rendering regressions cover day boundaries, New York winter/summer offsets, Tokyo, invalid dates and access checks. Two actual business configurations, each in four browser projects with Honolulu/de-DE clients and UTC server, preserve two appointment labels across account/management/reload. All seven table snapshots remain unchanged; zero writes. Cancellation command, authorization, timestamps, versions and eligibility rules are unchanged. Other private portal dates and translation catalogs remain outside this narrow repair.

V390 closes the existing customer cancellation path in the isolated reference. Account links to appointment management; no-results state is explicit. A synchronous in-flight guard allows one command; the client confirms only a receipt bound to appointment ID, organization, cancelled state and exact successor version. A lost/unverifiable/error response leaves the outcome unknown and commands disabled; an explicit GET reload resolves current state. There is no automatic mutation retry. The existing10second bounded request pattern is reused. Four actual browser projects prove normal cancellation/double click, committed response loss, unbound HTTP200 response, two concurrent requests and stale-page recovery. Sixteen appointments produce exactly16audits/16outbox events;28API writes include rejected duplicate/stale/foreign-resource attempts. Four browser/BFF boundary negatives per project and five inherited direct HTTP negatives remain enforced. All other durable data and appointment immutable fields are unchanged. The prior six-list/time-zone read suites also pass in eight browser/configuration runs with zero writes. No backend domain code, migration, cancellation eligibility, identity grant, provider or dependency changed. This is narrow reference closure, not full TEST02/production admission.
