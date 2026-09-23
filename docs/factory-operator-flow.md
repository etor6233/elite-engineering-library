# Factory operational progress

The factory portal now lets a role with both `factory:read` and `factory:write` record the existing `planned` → `assembly` transition for an authorized unit. Its source of truth remains `operations.Service` and the existing transactional PostgreSQL writer. The BFF only maps the named UI command to that owner; it does not add a transition graph, QC policy, inventory effect or production authorization.

`GET /v1/factory/units/{id}` reuses the factory list projection with exact tenant, purchase-order destination organization and unit predicates. The existing EnterpriseQueryModule mounts this read automatically; no additional host module is needed. The page derives organization and storage scope from the server session. The BFF validates permission, origin, strict body/query, streamed 4096-byte limit, scoped backend projection and expected state. The writer rechecks current state and scope inside its atomic state/outbox transaction, closing the read/write race.

The browser preserves an unresolved marker before POST and suppresses duplicate clicks and automatic mutation retries. Its recovery action performs only a GET. This is a current-state observation, not an idempotency receipt. Another actor may have advanced the unit after the lost response; the UI explicitly says the state does not identify the request that produced it. A successful HTTP response acknowledges the existing owner operation; an uncertain response stays uncertain until the operator reads the current state and decides a supported next action.

The focused browser gate uses an isolated PostgreSQL database, actual encrypted role sessions, RS256/JWKS bearer verification, production Next BFF and Go API. It starts with one planned fixture unit, loses the real assembly response after commit, performs a separate competing assembly→quality advancement through the same owner, reloads and observes quality without claiming request attribution. Read-only role, foreign organization, foreign tenant and stale command attempts fail. PostgreSQL contains one assembly event and one separate competing quality event; the new UI produced exactly one backend transition POST. No stock or money mutation is introduced.

Use the admitted runtime and dependency locks. Build with the direct Node Next CLI (`build --webpack` for the isolated external node_modules junction). Set `ELITE_FACTORY_BROWSER=1`, `ELITE_WEB_ROOT`, the absolute `ELITE_NODE_BIN`, and `FACTORY_BROWSER_DB_URL` for a fresh loopback database named `elite_payment_connected_*` with the composition migrations. The fixture has an exact environment/database guard and must not target business data. With module downloads disabled, invoke the admitted Go executable:

```text
go test -mod=readonly ./internal/platform/postgres -run ^TestFactoryBrowserPostgres$ -count=1 -v -timeout=4m
```

The fixture invokes the installed Playwright CLI directly, without package-manager execution or downloads. It requires the reference composition's `handoverBrowserIssuer` test helper, reused for local generated identity credentials; it does not rerun the handover gate. Chromium desktop plus a 390×844 viewport check passed with zero page errors and no horizontal overflow. This is not a four-browser claim. Omitted opt-in produces an explicit skip; acceptance must reject skips. The self-signed HTTPS exception remains restricted to the existing explicit loopback fixture configuration.

All delta code is AUTHORED read/transport/UI/test glue. The existing factory state graph is also AUTHORED and is not attributed to a vendor. No dependency, license or notice is added. This gate does not close general person management, training/evaluation, factory registration or all procurement/inventory UI capabilities.
