# SALES / QUOTES / ORDERS — compose slice

## Tip pin

Elite library tip: `477358d0d7090f2c39467abe5634d81f38afae7f`

This slice documents how the **selected 116** franchise composition materializes the sales journey from CRM lead through server-priced quotation to a durable `placed` customer order. It is **HECHO-first local fixture evidence** — not production authorization, not live payments, not REVESTEX product code.

## Scope

| In scope | Out of scope (sibling agents / later slices) |
|---|---|
| `sales.quotation` issuance with server-side price | Payment capture (`GO-COMMERCE-PRICING-PAYMENT-API` payment-request lane) |
| Quote acceptance → `sales.customer_order` + line | Stock allocation / reservation (`ORDER_STOCK_CONNECTED_V296`) |
| Immutable `sales.quotation_acceptance` + outbox | Fiscal issuance, ARCA, tax jurisdiction |
| BC-derived quote→order field transfer | SDR/GTM packs (honest GAP — none admitted) |
| Local PostgreSQL fixtures | Live Salesforce / Dynamics connectors |

## Pack map (selected 116)

All `packId` / `version` values below are taken from `markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md` at the tip pin. Do not add undeclared packs.

| packId | version | elite path | role in this slice |
|---|---|---|---|
| `ELECTROMOBILITY-FRANCHISE-MODULES` | `0.1.2` | `implementation_packs/ELECTROMOBILITY_FRANCHISE_MODULES.md` | PostgreSQL schemas `catalog`, `crm`, `pricing`, `sales` (migration `db/migrations/0003_electromobility_franchise_modules.up.sql`) |
| `GO-ENTERPRISE-BACKEND` | `0.4.8` | `implementation_packs/GO_ENTERPRISE_BACKEND_CORE.md` | Tenant/org/platform foundation consumed by journey + commerce owners |
| `GO-ELECTROMOBILITY-PUBLIC-CRM-API` | `0.4.0` | `implementation_packs/GO_ELECTROMOBILITY_PUBLIC_CRM_API.md` | Upstream lead/customer capture into `crm.*` |
| `GO-BC-EXACT-AMOUNT-ADAPTER` | `0.1.0` | `implementation_packs/GO_BC_EXACT_AMOUNT_ADAPTER.md` | Required MIT notice dependency for BC sales derivation |
| `GO-BC-SALES-CONTRACT-ADAPTER` | `0.1.0` | `implementation_packs/GO_BC_SALES_CONTRACT_ADAPTER.md` | Active price-book eligibility + `bcsales.TransferQuoteToOrder` |
| `GO-COMMERCE-PRICING-PAYMENT-API` | `0.8.0` | `implementation_packs/GO_COMMERCE_PRICING_PAYMENT_API.md` | Canonical order owner; pricing server-side (payment lane deferred here) |
| `GO-FRANCHISE-CUSTOMER-JOURNEY-API` | `0.10.21` | `implementation_packs/GO_FRANCHISE_CUSTOMER_JOURNEY_API.md` | Quote issue/accept HTTP + PostgreSQL transaction (`internal/franchisejourney`, `internal/platform/postgres/franchisejourney.go`) |
| `TS-FRANCHISE-JOURNEY-PORTALS` | `0.21.0` | `implementation_packs/TYPESCRIPT_FRANCHISE_JOURNEY_PORTALS.md` | Same-origin BFF + `/customer/quotes` UI (browser gate sibling; not required for this backend-only slice) |

Related reconstruction evidence (not additional packIds):

| evidence | elite path | claim |
|---|---|---|
| Quote→order V116 | `reconstruction_evidence/FRANCHISE_QUOTE_TO_ORDER_2026-08-29_V116.md` | Atomic acceptance creates one `placed` order, one line at server price, acceptance row, two outbox events |
| Quote acceptance connected V294 | `reconstruction_evidence/QUOTE_ACCEPTANCE_CONNECTED_V294.md` | Portal/BFF/Go/PG connected acceptance + recovery harness |
| Quote creation recovery V306 | `reconstruction_evidence/QUOTE_CREATION_RECOVERY_V306.md` | Idempotent quote issuance + scoped recovery |
| BC sales adaptation V402 | `reconstruction_evidence/BC_SALES_CONTRACT_ADAPTATION_V402.md` | Narrow BCApps quote/order field transfer gate |
| Order stock connected V296 | `reconstruction_evidence/ORDER_STOCK_CONNECTED_V296.md` | Downstream stock/payment sibling (not in this slice) |

## Compose recipe

Materialize from the selected 116 plan (`markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md`) using the franchise operating protocol (`markdown_system/FRANCHISE_PROJECT_OPERATING_PROTOCOL.md` via `START_FRANCHISE.md`). For this slice, the minimum connected subgraph is:

```text
ELECTROMOBILITY-FRANCHISE-MODULES (schemas)
  → GO-ELECTROMOBILITY-PUBLIC-CRM-API (lead/customer)
  → GO-BC-EXACT-AMOUNT-ADAPTER + GO-BC-SALES-CONTRACT-ADAPTER (price + quote→order copy)
  → GO-FRANCHISE-CUSTOMER-JOURNEY-API (issue/accept transaction)
  → GO-COMMERCE-PRICING-PAYMENT-API (order owner; payment deferred)
```

### Transaction flow (quote → order)

1. **Lead exists** in `crm.lead` scoped to `org.organization` (`ELECTROMOBILITY-FRANCHISE-MODULES` + CRM API).
2. **Issue quote** — `CreateQuote` / `CreateQuoteAs`:
   - Resolves `currency` + `total_minor_units` from `pricing.price_book` + `pricing.price_book_entry` using the BC eligibility predicate (`GO-BC-SALES-CONTRACT-ADAPTER` → `internal/platform/postgres/bc_price_sql.go`).
   - Persists `sales.quotation` state `issued` and emits `quotation.issued` outbox event.
   - Idempotency: `platform.idempotency_record` scope `franchise-quote` (`QUOTE_CREATION_RECOVERY_V306`).
3. **Accept quote** — `AcceptQuote`:
   - Locks issued, unexpired, customer-scoped quotation.
   - Calls `bcsales.TransferQuoteToOrder` (`internal/bcsales/quote.go`, derived from pinned BCApps `SalesQuotetoOrder.Codeunit.al` per `docs/provenance/BC_SALES_DERIVATION.md`).
   - Inserts `sales.customer_order` (`placed`), `sales.customer_order_line`, `sales.quotation_acceptance`.
   - Updates quotation `state=accepted`, binds `order_id`.
   - Emits `quotation.accepted` + `customer-order.placed` outbox events in the same transaction (`FRANCHISE_QUOTE_TO_ORDER_2026-08-29_V116`).
4. **Stop** — no `payment.payment_attempt`, no stock allocation, no handover in this slice.

Migration anchor for acceptance schema: `db/migrations/0006_quote_order_conversion.up.sql`.

## Fixtures and harness

| artifact | elite path | what it proves |
|---|---|---|
| **Slice harness (new)** | `internal/platform/postgres/sales_orders_compose_integration_test.go` | Minimal tenant seed → issue quote (server price `123456 ARS`) → accept → `placed` order + acceptance + outbox; stale accept rejected |
| Journey persistence gate | `internal/platform/postgres/franchisejourney_integration_test.go` (`TestFranchiseJourneyPersistenceIsolationAndReplay`) | Full quote issue/accept/replay/rollback matrix on disposable PG |
| BC price eligibility | `internal/platform/postgres/bc_price_integration_test.go` | 180-vector BC predicate vs local half-open contract |
| BC quote transfer unit | `internal/bcsales/quote_test.go` | `TransferQuoteToOrder` conservation without database |
| SQL invariant | `db/tests/0005_franchise_customer_journey.test.sql` | Schema-level quotation/order/handover presence |
| Payment-connected seed (downstream) | `internal/platform/postgres/payment_connected_integration_test.go` (`connectedSeedOrder`) | Reuses quote→order before payment lane (sibling; requires `PAYMENT_CONNECTED_DB_URL`) |

### Evidence log line

When the slice harness passes against a migrated disposable database, it emits:

```text
SALES_ORDERS_COMPOSE_PASS tenant=<uuid> quote=quote order=order placed=true acceptance=true outbox=true server_price=true payment=false
```

## Local verify steps

Prerequisites: PostgreSQL 18.6 loopback, all migrations applied to a **disposable** database named `elite_confirmation_<unique>` (same guard as other journey/commerce integration tests).

```bash
# 1. Create disposable DB and apply migrations (project-specific tooling; see START_REFERENCE_V402.md)
export TEST_DATABASE_URL='postgresql://user@127.0.0.1/elite_confirmation_sales_compose?sslmode=disable'

# 2. Slice harness (this PR)
go test ./internal/platform/postgres -run TestSalesOrdersComposeQuoteToOrder -count=1 -v

# 3. Unit gate — BC quote transfer (no database)
go test ./internal/bcsales -run TransferQuoteToOrder -count=1 -v

# 4. Optional deeper gates already in library
go test ./internal/platform/postgres -run TestFranchiseJourneyPersistenceIsolationAndReplay -count=1 -v
go test ./internal/bcsales ./internal/franchisejourney -count=1

# 5. SQL invariant (requires psql + migrated DB)
psql "$TEST_DATABASE_URL" -v ON_ERROR_STOP=1 -f db/tests/0005_franchise_customer_journey.test.sql
```

Without `TEST_DATABASE_URL`, integration tests **skip** (not PASS). A skip is not production evidence.

Preflight only (no product materialization): `VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Preflight` from library root per `START_FRANCHISE.md`.

## Honest boundaries

- Proves **local** quote→order composition on PostgreSQL fixtures.
- Does **not** prove legal acceptance wording, tax/fiscal rules, financing, target IdP, CDN/WAF, provider live traffic, or business acceptance (`FRANCHISE_QUOTE_TO_ORDER_2026-08-29_V116` §Honest boundary).
- Payment capture is intentionally deferred to the payments sibling agent (`GO-COMMERCE-PRICING-PAYMENT-API` V298+ / `payment_connected_integration_test.go`).
- No SDR/GTM pack exists in selected 116/120 (`markdown_system/LIBRARY_HEALTH_CHECK.md` honest GAP).

## Gaps — public patterns only (cite, do not implement)

These are **reference patterns** for consumer projects. The library does not ship Salesforce or Dynamics connectors for this slice.

| pattern | public source | note vs this library |
|---|---|---|
| Salesforce Order Management lifecycle | [Order Management overview](https://help.salesforce.com/s/articleView?id=sf.order_mgmt_overview.htm) | External orchestration of order fulfillment; library uses PostgreSQL `sales.customer_order` as canonical owner |
| Salesforce quote management | [Salesforce Quotes](https://help.salesforce.com/s/articleView?id=sf.quotes.htm) | Client-submitted totals are anti-pattern here; server-owned pricing is mandatory (`GO-COMMERCE-PRICING-PAYMENT-API` §Architecture contract) |
| Dynamics 365 Sales — create sales quotes | [Create quotes](https://learn.microsoft.com/en-us/dynamics365/sales/create-edit-quote-sales) | BCApps quote→order copy is narrowly adapted via `GO-BC-SALES-CONTRACT-ADAPTER`; consent/expiry/retention remain AUTHORED |
| Dynamics 365 — post sales documents | [Post sales documents](https://learn.microsoft.com/en-us/dynamics365/business-central/ui-post-sales) | Posting/ledger consequences are out of scope; local state stops at `placed` |
| Dynamics 365 Commerce order exchanges | [Order exchanges](https://learn.microsoft.com/en-us/dynamics365/commerce/orderexchanges) | Return/exchange is a separate worker slice (`GO-RETURN-EXCHANGE-FULFILLMENT-WORKER`) |

## Related doors

1. `docs/FRANCHISE_PLAYBOOK.md` — franchise assembly context (do not edit `docs/ROADMAP.md`; PR #13 owns it)
2. `markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md` — selected 116 packIds
3. `markdown_system/LIBRARY_HEALTH_CHECK.md` — honest wiring GAPs
4. `reconstruction_evidence/LIBRARY_VS_PRODUCT_GATE_V402.md` — library vs product gate
