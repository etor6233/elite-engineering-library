# Customer 360 read model — compose plan (not a HECHO pack)

**Revision:** `CUSTOMER-360-READ-MODEL-PLAN-0.1.0`  
**Tip pin:** `227c6e0033eaab96cb42fa809ac6c37b758f84f4`  
**Status:** **PLAN ONLY** — no `implementation_packs/GO_CUSTOMER_360_*.md`, not in selected 116/120 JSON, not `REBUILD_VERIFIED`, not production authorization.

## Architecture contract (read-only cites @ tip `227c6e0`)

| Door | Path | CRM row |
| --- | --- | --- |
| Domain architecture | [`docs/FRANCHISE_ARCHITECTURE.md`](../docs/FRANCHISE_ARCHITECTURE.md) § **3. Customers + history / CRM** | **PARCIAL** — compose + fixtures; no standalone segmentation / Customer-360 pack |
| Build-order matrix | [`docs/ROADMAP.md`](../docs/ROADMAP.md) § **FRANCHISE / COMPANY DOMAINS** → Customers + history / CRM | **PARCIAL** — selected packs below; segmentation/360 **FALTA** as admitted pack |
| Tangible slice | [`docs/slices/CRM_HISTORY_COMPOSE.md`](../docs/slices/CRM_HISTORY_COMPOSE.md) | LOCAL_FIXTURES timeline + history |
| Slice evidence | [`docs/slices/CRM_HISTORY_FIXTURE_EVIDENCE.json`](../docs/slices/CRM_HISTORY_FIXTURE_EVIDENCE.json) | machine-readable fixture index |
| Selected composition | [`markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md`](FRANCHISE_COMPLETE_PACK_PLAN.md) | authoritative `packId` / `path` / `version` |

**Proposed addition (this file):** query-composition recipe over existing journey + CRM APIs. Does **not** close the architecture **PARCIAL** row to **HECHO**; it documents how a consumer composes reads honestly until a future admitted pack is gap-gated.

## Claim (narrow)

Define a **Customer 360 read model** as a **BFF/query composition** over admitted owners:

1. **Identity + capture** — `crm.customer_profile`, `crm.lead`, `crm.consent_evidence` (public CRM API + franchise modules schema).
2. **Operational history + timeline** — `crm.appointment_transition` ledger + `GET /v1/customer/journey` projection (`GO-FRANCHISE-CUSTOMER-JOURNEY-API`).
3. **Channel identity bindings** — `communication.contact_channel_binding` (+ decision ledger) resolved via `GO-PG-CONTACT-CHANNEL-IDENTITY` (no raw external id in clear).

No new durable write owner. No named `*SDR*` or `*GTM*` pack. No Dynamics 365 code copy.

## Selected packs (primary composition paths)

Exact filenames from [`markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md`](FRANCHISE_COMPLETE_PACK_PLAN.md):

| packId | version | Selected JSON `path` | Read-model lane |
| --- | --- | --- | --- |
| `GO-ELECTROMOBILITY-PUBLIC-CRM-API` | `0.4.0` | `implementation_packs/GO_ELECTROMOBILITY_PUBLIC_CRM_API.md` | Public lead capture (`POST /v1/public/{tenantCode}/{organizationCode}/leads`); admin catalog surfaces; seeds `crm.lead` + `crm.consent_evidence` |
| `GO-FRANCHISE-CUSTOMER-JOURNEY-API` | `0.10.21` | `implementation_packs/GO_FRANCHISE_CUSTOMER_JOURNEY_API.md` | Customer-scoped timeline (`GET /v1/customer/journey`); append-only `crm.appointment_transition`; quotes/handovers/delivery exceptions |
| `GO-PG-CONTACT-CHANNEL-IDENTITY` | `0.1.1` | `implementation_packs/GO_PG_CONTACT_CHANNEL_IDENTITY.md` | Versioned channel→lead binding; HMAC external id; resolver for conversation runtime — **no** customer-facing unified GET today |

### Supporting owners (compose context — not 360-specific packs)

| packId | Selected JSON `path` | Role in 360 compose |
| --- | --- | --- |
| `ELECTROMOBILITY-FRANCHISE-MODULES` | `implementation_packs/ELECTROMOBILITY_FRANCHISE_MODULES.md` | `crm.*` schema owner (`customer_profile`, `lead`, `consent_evidence`) |
| `GO-ENTERPRISE-QUERY-API` | `implementation_packs/GO_ENTERPRISE_QUERY_API.md` | Adjacent reads: `GET /v1/customer/orders`, `GET /v1/customer/service-cases` (pagination; not merged into journey JSON) |
| `GO-CUSTOMER-SURVEY-API` | `implementation_packs/GO_CUSTOMER_SURVEY_API.md` | Post-interaction feedback lane (`GET/POST /v1/customer/surveys/...`) — orthogonal to journey timeline |
| `TS-OIDC-PORTAL-ADAPTER` | `implementation_packs/TYPESCRIPT_OIDC_PORTAL_ADAPTER.md` | BFF client pattern for protected GETs (`protectedGet` → `/v1/customer/journey`) |

### Adjacent ingress (**not** SDR)

| packId | Selected JSON `path` | Boundary |
| --- | --- | --- |
| `GO-OMNICHANNEL-LEAD-INGRESS` | `implementation_packs/GO_OMNICHANNEL_LEAD_INGRESS.md` | Durable provider candidate + outbox only |
| `GO-LEAD-CANDIDATE-PROMOTION` | `implementation_packs/GO_LEAD_CANDIDATE_PROMOTION.md` | Explicit promotion → `crm.lead`; **FALTA** named `*SDR*` pack (`markdown_system/PACK_PER_CLAIM_INDEX.md`) |

## Owner map (who owns what in the read model)

```text
┌─────────────────────────────────────────────────────────────────────────┐
│ Consumer BFF / portal (project-owned compose — NOT an elite pack body) │
└─────────────────────────────────────────────────────────────────────────┘
         │ parallel GETs (scoped by verified subject + organization_id)
         ▼
┌──────────────────────┐  ┌──────────────────────────┐  ┌─────────────────────────┐
│ GO-ELECTROMOBILITY-  │  │ GO-FRANCHISE-CUSTOMER-   │  │ GO-PG-CONTACT-CHANNEL-  │
│ PUBLIC-CRM-API       │  │ JOURNEY-API              │  │ IDENTITY                │
│ (capture + catalog)  │  │ (timeline projection)    │  │ (binding resolver)      │
└──────────┬───────────┘  └────────────┬─────────────┘  └────────────┬────────────┘
           │                           │                              │
           ▼                           ▼                              ▼
   crm.customer_profile         crm.appointment                 communication.
   crm.lead                     crm.appointment_transition      contact_channel_binding
   crm.consent_evidence         sales.quotation / handover*     contact_channel_binding_decision
                                delivery exception rows*
   * journey projection tables — see GO-FRANCHISE-CUSTOMER-JOURNEY-API migrations
```

| Data facet | Authoritative owner | Exposed HTTP (local reference) | 360 compose note |
| --- | --- | --- | --- |
| Customer principal | `ELECTROMOBILITY-FRANCHISE-MODULES` → `crm.customer_profile` | Indirect via journey scope | BFF may join profile display fields from SQL or future admin read — **no** dedicated `GET /v1/customer/profile` in selected packs |
| Leads + consent | `GO-ELECTROMOBILITY-PUBLIC-CRM-API` | `POST /v1/public/.../leads` (write) | Read back via SQL/admin paths or journey lead linkage — not a public 360 aggregate |
| Appointment history | `GO-FRANCHISE-CUSTOMER-JOURNEY-API` | `GET /v1/customer/journey` → `appointments` + implicit transition ledger | **PROVEN_LOCAL** via `db/tests/0089_crm_customer_timeline_fixture.test.sql` |
| Quotes / handovers / delivery exceptions | `GO-FRANCHISE-CUSTOMER-JOURNEY-API` | same journey GET collections | Caps: 100 rows/collection; 1 MiB body budget |
| Channel identities | `GO-PG-CONTACT-CHANNEL-IDENTITY` | Internal `ContactResolver` + binding command path | Operator 360 may list bindings by `lead_id` in SQL; customer token must not receive raw channel ids |
| Orders (commercial) | `GO-ENTERPRISE-QUERY-API` | `GET /v1/customer/orders` | Separate pagination contract — merge in BFF only |
| Service cases | `GO-ENTERPRISE-QUERY-API` | `GET /v1/customer/service-cases` | Narrow rejection/resolution lane; not full case management |
| Surveys | `GO-CUSTOMER-SURVEY-API` | `/v1/customer/surveys/{survey}` | Optional satisfaction facet |

## Compose recipe (read model assembly)

1. **Pin** tip `227c6e0033eaab96cb42fa809ac6c37b758f84f4`; walk franchise doors per [`START_FRANCHISE.md`](../START_FRANCHISE.md).
2. **Select** primary three packs from [`markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md`](FRANCHISE_COMPLETE_PACK_PLAN.md) plus foundation rows in [`docs/slices/CRM_HISTORY_COMPOSE.md`](../docs/slices/CRM_HISTORY_COMPOSE.md) §2.
3. **Apply migrations** through CRM + journey + contact identity owners (minimum): `0003`, `0005`, `0015`, `0047` — paths listed in slice compose doc.
4. **Authenticate** customer principal with OIDC scope `customer:self` (`GO-OIDC-PORTAL-SESSION` + `TS-OIDC-PORTAL-ADAPTER`); subject must match `customer_principal_id`.
5. **Fetch timeline core** — `GET /v1/customer/journey?organization_id=<org>` (mandatory 360 spine).
6. **Fetch adjacent facets** (optional, parallel):
   - `GET /v1/customer/orders?organization_id=<org>&limit=25`
   - `GET /v1/customer/service-cases?organization_id=<org>&limit=25`
   - Survey definition/response endpoints when feedback is in scope.
7. **Resolve channel bindings** (operator/admin 360 only): query `communication.contact_channel_binding` by `lead_id` through authorized admin SQL/API — never expose HMAC preimage or raw external id to customer tokens.
8. **Merge in BFF** into a single view model document; enforce org isolation on every sub-query; fail closed on partial errors (mirror `CustomerJourney` no-partial-json rule in `GO-FRANCHISE-CUSTOMER-JOURNEY-API`).

Example BFF merge shape (illustrative — not shipped as elite code):

```json
{
  "customer_principal_id": "customer-a",
  "organization_id": "store",
  "timeline": { "appointments": [], "quotes": [], "handovers": [], "delivery_exceptions": [] },
  "orders_page": { "items": [], "next": null },
  "service_cases_page": { "items": [], "next": null },
  "channel_bindings": null
}
```

`channel_bindings` remains **FALTA** for customer-scoped HTTP until a future gap-gated read pack; operator compose may populate from SQL.

## Verify steps (LOCAL_FIXTURES only)

Evidence index: [`docs/slices/CRM_HISTORY_FIXTURE_EVIDENCE.json`](../docs/slices/CRM_HISTORY_FIXTURE_EVIDENCE.json).

### SQL invariants (disposable `TEST_DATABASE_URL`)

```bash
for f in db/tests/0003_electromobility_franchise_modules.test.sql \
         db/tests/0005_franchise_customer_journey.test.sql \
         db/tests/0015_franchise_availability_and_appointment_audit.test.sql \
         db/tests/0089_crm_customer_timeline_fixture.test.sql \
         db/tests/0047_contact_channel_identity.test.sql; do
  psql "$TEST_DATABASE_URL" -v ON_ERROR_STOP=1 -f "$f"
done
```

| Fixture | Proves for 360 compose |
| --- | --- |
| `db/tests/0089_crm_customer_timeline_fixture.test.sql` | Customer-scoped timeline counts + org isolation (journey spine) |
| `db/tests/0047_contact_channel_identity.test.sql` | Binding immutability + no raw channel id leak (channel facet) |
| `db/tests/0005_franchise_customer_journey.test.sql` | Lead/profile/quote graph coexistence |
| `db/tests/0015_franchise_availability_and_appointment_audit.test.sql` | Append-only `crm.appointment_transition` |

### Go integration / HTTP

```bash
go test ./internal/platform/postgres -run 'CustomerJourney|ContactIdentity' -count=1 -timeout=5m
go test ./internal/platform/httpapi -run 'TestCustomerJourneyReadFailure|TestAppointmentConfirmation' -count=1 -timeout=5m
go test ./internal/platform/postgres -run ContactIdentity -count=1 -timeout=5m
```

### Library preflight (optional)

```powershell
./VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Preflight
```

Pass criteria for **this plan**:

- Journey spine fixtures PASS.
- Contact-channel identity fixtures PASS.
- No single elite endpoint returns a merged “360” JSON today — compose is **documented**, not **shipped**.

## Gaps — Microsoft patterns (cite only)

Official Dynamics 365 **guidance** informs remaining gaps; no Dynamics code admission, no BC product selection:

| Pattern | Official reference | Local status @ `227c6e0` |
| --- | --- | --- |
| Prospect → quote | [Prospect to quote overview](https://learn.microsoft.com/en-us/dynamics365/guidance/business-processes/prospect-to-quote-overview) | Fixture quotes via journey API; financing/tax/operator UX **CONDITIONED** |
| Case → resolution | [Case to resolution introduction](https://learn.microsoft.com/en-us/dynamics365/guidance/business-processes/case-to-resolution-introduction) | Delivery exceptions + narrow `service-cases` read — **FALTA** full case SLA/knowledge hub |
| Customer Service hub | [Customer Service overview](https://learn.microsoft.com/en-us/dynamics365/customer-service/overview) | **FALTA** omnichannel agent desktop |
| Unified activity timeline | [Activity feed (on-premises)](https://learn.microsoft.com/en-us/dynamics365/customerengagement/on-premises/basics/set-up-interaction-timeline) | Journey projection only — **FALTA** cross-channel activity stream |
| Segmentation / marketing lists | ROADMAP analytics row | **FALTA** standalone segmentation pack |
| Named SDR workflow | `markdown_system/PACK_PER_CLAIM_INDEX.md` | **FALTA** `*SDR*` pack — ingress + promotion only |

Also: `markdown_system/BUSINESS_FUNCTION_OPERATING_CONTRACT_V403.md` and `roles/method-sources.v403.json` (prospect-to-quote, case-to-resolution queries with explicit no-BC-selection limits).

## Explicit non-claims

- This file is a **pack plan**, not an `implementation_packs/*` HECHO body and not added to selected 116/120 JSON.
- Does **not** invent `GO-CUSTOMER-360-*`, `*SDR*`, or `*GTM*` packs.
- `GO-OMNICHANNEL-LEAD-INGRESS` / `GO-LEAD-CANDIDATE-PROMOTION` are **ingress**, not SDR or 360 owners.
- PASS on cited fixtures ≠ cloud CI, live IdP, provider accounts, operator 360 UI, or production CRM.
- `production_authorized: false` per `qualification/FINAL_LIBRARY_READY_V402.json`.
- `HISTORY-MODEL-TRAINING-PIPELINE` “history” is a different contract (`markdown_system/PROJECT_HISTORY_MODEL_TRAINING_CONTRACT.md`).

## Evidence cross-links

| Receipt | Path |
| --- | --- |
| Public CRM lead idempotency | `reconstruction_evidence/GO_ELECTROMOBILITY_PUBLIC_CRM_API_2026-08-24_V2.md` |
| Journey vertical slice | `reconstruction_evidence/FRANCHISE_CUSTOMER_JOURNEY_2026-08-29_V114.md` |
| Contact channel identity | `reconstruction_evidence/GO_PG_CONTACT_CHANNEL_IDENTITY_2026-09-04_V237.md` |
| Timeline BFF confirmation | `reconstruction_evidence/APPOINTMENT_CONFIRMATION_BFF_POSTGRES_V262.md` |
| Runtime routes | `docs/electromobility-api-runtime.md` |
| CRM history slice | `docs/slices/CRM_HISTORY_COMPOSE.md` |

## Nightly verify commands (library PR checklist)

```bash
# CRM / 360 compose spine
psql "$TEST_DATABASE_URL" -v ON_ERROR_STOP=1 -f db/tests/0089_crm_customer_timeline_fixture.test.sql
psql "$TEST_DATABASE_URL" -v ON_ERROR_STOP=1 -f db/tests/0047_contact_channel_identity.test.sql
go test ./internal/platform/postgres -run CustomerJourney -count=1 -timeout=5m
go test ./internal/platform/httpapi -run TestCustomerJourneyReadFailure -count=1 -timeout=5m
```

```powershell
./VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Preflight
```
