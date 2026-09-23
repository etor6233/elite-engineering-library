# CRM / CUSTOMERS + HISTORY — tangible slice compose

Tip pin: `227c6e0033eaab96cb42fa809ac6c37b758f84f4`

Status: **PARCIAL → tangible slice** (local fixtures + selected composition only). **Not** production, **not** REVESTEX, **not** a named SDR or GTM pack.

## 0. Architecture contract (read-only cites @ tip `227c6e0`)

This slice doc does **not** edit the architecture doors; it aligns with them:

| Door | Path @ `227c6e0` | CRM / history row |
| --- | --- | --- |
| Domain architecture | [`docs/FRANCHISE_ARCHITECTURE.md`](FRANCHISE_ARCHITECTURE.md) § **3. Customers + history / CRM** | **PARCIAL** — compose + fixtures; no standalone segmentation / Customer-360 pack |
| Build-order matrix | [`docs/ROADMAP.md`](ROADMAP.md) § **FRANCHISE / COMPANY DOMAINS** → Customers + history / CRM | **PARCIAL** — cite selected packs with full `implementation_packs/<NAME>.md` paths below; segmentation/360 **FALTA** |
| Selected JSON | [`markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md`](../markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md) | Exact `packId` + `path` + `version` per row — no shorthand filenames |
| Customer 360 compose plan | [`markdown_system/CUSTOMER_360_READ_MODEL_PACK_PLAN.md`](../markdown_system/CUSTOMER_360_READ_MODEL_PACK_PLAN.md) | **PLAN ONLY** — BFF query composition over journey + CRM APIs; not a HECHO pack |

Lead ingress remains adjacent per [`docs/FRANCHISE_ARCHITECTURE.md`](FRANCHISE_ARCHITECTURE.md) layer **4 — CRM + leads**; no named SDR pack (`markdown_system/PACK_PER_CLAIM_INDEX.md`).

## 1. Claim (narrow)

Prove a **single CRM owner** in PostgreSQL with:

- durable **customer** identity (`crm.customer_profile`);
- **lead** + **consent** evidence (`crm.lead`, `crm.consent_evidence`);
- **append-only operational history** (`crm.appointment_transition`, lead transition enforcement);
- a **customer-scoped timeline read model** (`GET /v1/customer/journey` → appointments, quotes, handovers, delivery exceptions).

Lead **ingress** (Google/Meta/TikTok → candidate → promotion) is adjacent wiring only; it does **not** substitute an SDR pack.

## 2. Pack map (selected composition + schema owner)

Authoritative selector: [`markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md`](../markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md) (116 packs). Rows below are the **minimum** subset for this slice; `packId`, `path`, and `version` match the selected JSON exactly — use the **`path` column filename**, not shorthand aliases.

### PackId → selected JSON filename (alias normalization)

| Selected JSON `packId` | Exact selected JSON `path` | Do **not** substitute |
| --- | --- | --- |
| `PG-TX-FOUNDATION` | `implementation_packs/POSTGRES_TRANSACTIONAL_FOUNDATION.md` | `POSTGRES-TRANSACTIONAL-FOUNDATION.md`, `PG-TX-FOUNDATION.md`, bare `POSTGRES_TRANSACTIONAL_FOUNDATION.md` |
| `GO-ENTERPRISE-BACKEND` | `implementation_packs/GO_ENTERPRISE_BACKEND_CORE.md` | `GO_ENTERPRISE_BACKEND.md`, `GO-ENTERPRISE-BACKEND-CORE.md`, bare `GO_ENTERPRISE_BACKEND_CORE.md` |

| Role | packId | version | Selected JSON `path` | What it owns here |
| --- | --- | --- | --- | --- |
| PostgreSQL foundation | `PG-TX-FOUNDATION` | `0.1.0` | `implementation_packs/POSTGRES_TRANSACTIONAL_FOUNDATION.md` | tenant, org, platform idempotency/outbox |
| Backend shell | `GO-ENTERPRISE-BACKEND` | `0.4.8` | `implementation_packs/GO_ENTERPRISE_BACKEND_CORE.md` | HTTP/OIDC/identity wiring |
| **CRM schema owner** | `ELECTROMOBILITY-FRANCHISE-MODULES` | `0.1.2` | `implementation_packs/ELECTROMOBILITY_FRANCHISE_MODULES.md` | `crm` schema: `customer_profile`, `lead`, `consent_evidence`; catalog/sales FK graph |
| Public CRM capture | `GO-ELECTROMOBILITY-PUBLIC-CRM-API` | `0.4.0` | `implementation_packs/GO_ELECTROMOBILITY_PUBLIC_CRM_API.md` | public catalog + consented lead capture + durable idempotency |
| Customer journey + **history** | `GO-FRANCHISE-CUSTOMER-JOURNEY-API` | `0.10.21` | `implementation_packs/GO_FRANCHISE_CUSTOMER_JOURNEY_API.md` | appointments, quotes, handovers; `crm.appointment_transition`; `CustomerJourney` projection; `GET /v1/customer/journey` |
| Runtime wiring | `GO-ELECTROMOBILITY-APPLICATION` | `1.22.3` | `implementation_packs/GO_ELECTROMOBILITY_APPLICATION.md` | `cmd/electromobility-api` route activation |
| Portal BFF (timeline client) | `TS-OIDC-PORTAL-ADAPTER` | `0.3.1` | `implementation_packs/TYPESCRIPT_OIDC_PORTAL_ADAPTER.md` | protected GET to `/v1/customer/journey` |
| Web bridge | `TS-GO-API-WEB-BRIDGE` | `0.15.2` | `implementation_packs/TYPESCRIPT_GO_API_WEB_BRIDGE.md` | typed backend client |
| Journey UI (optional local) | `TS-FRANCHISE-JOURNEY-PORTALS` | `0.21.0` | `implementation_packs/TYPESCRIPT_FRANCHISE_JOURNEY_PORTALS.md` | customer portal surfaces |

### Schema materialization (ELECTROMOBILITY-FRANCHISE-MODULES)

Exact manifest from the pack:

```text
db/migrations/0003_electromobility_franchise_modules.up.sql
db/migrations/0003_electromobility_franchise_modules.down.sql
db/tests/0003_electromobility_franchise_modules.test.sql
```

CRM tables introduced there (excerpt):

- `crm.customer_profile` — principal customer row, normalized email/phone checks
- `crm.lead` — organization-scoped lead, optional link to `customer_principal_id`
- `crm.consent_evidence` — purpose/policy/decision + evidence hash

Journey/history extensions arrive in later migrations owned by `GO-FRANCHISE-CUSTOMER-JOURNEY-API` (notably `0005`, `0015` for `crm.appointment_transition` immutability).

### Adjacent ingress (cite only — **not** SDR)

| packId | version | Selected JSON `path` | Boundary |
| --- | --- | --- | --- |
| `GO-OMNICHANNEL-LEAD-INGRESS` | `0.2.0` | `implementation_packs/GO_OMNICHANNEL_LEAD_INGRESS.md` | durable provider candidate + outbox; **no** CRM promotion by receipt alone |
| `GO-LEAD-CANDIDATE-PROMOTION` | `0.1.0` | `implementation_packs/GO_LEAD_CANDIDATE_PROMOTION.md` | explicit contact decision → `crm.lead` + `crm.consent_evidence`; **not** a named `*SDR*` pack |
| `GO-PG-CONTACT-CHANNEL-IDENTITY` | `0.1.1` | `implementation_packs/GO_PG_CONTACT_CHANNEL_IDENTITY.md` | channel identity binding to existing CRM lead |

See `markdown_system/PACK_PER_CLAIM_INDEX.md` (SDR row: **PARCIAL**, **FALTA** pack nombrado `*SDR*`).

## 3. Compose recipe

1. **Pin** tip `227c6e0033eaab96cb42fa809ac6c37b758f84f4`; walk doors 1–7 in [`docs/FRANCHISE_PLAYBOOK.md`](../FRANCHISE_PLAYBOOK.md) before composing.
2. **Select** the pack rows in §2 from [`markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md`](../markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md); do not add undeclared packIds; resolve filenames via the alias table (§2), not shorthand.
3. **Materialize** with the library compositor / franchise protocol in [`START_FRANCHISE.md`](../START_FRANCHISE.md) into an empty consumer tree (or use this repo’s already-materialized paths under `internal/`, `db/`, `cmd/electromobility-api/`).
4. **Apply migrations** in order through journey/history owners at minimum:
   - `0001` — `implementation_packs/POSTGRES_TRANSACTIONAL_FOUNDATION.md` (`PG-TX-FOUNDATION`)
   - `0002` — `implementation_packs/GO_ENTERPRISE_BACKEND_CORE.md` (`GO-ENTERPRISE-BACKEND`) enterprise order core
   - `0003` — `implementation_packs/ELECTROMOBILITY_FRANCHISE_MODULES.md` (**CRM schema**)
   - `0005` — `implementation_packs/GO_FRANCHISE_CUSTOMER_JOURNEY_API.md`
   - `0015` — `implementation_packs/GO_FRANCHISE_CUSTOMER_JOURNEY_API.md` (**`crm.appointment_transition`**)
   - further migrations only if the consumer activates quotes, handovers, ingress, or surveys
5. **Configure** runtime per `docs/electromobility-api-runtime.md`: `DATABASE_URL`, `OIDC_ISSUER`, `OIDC_AUDIENCE`; synthetic RS256 fixtures for local gates only.
6. **Expose** customer timeline: `GET /v1/customer/journey?organization_id=<org>` with token scope `customer:self` and subject bound to `customer_principal_id`.

Data flow (local slice):

```text
visitor lead (public CRM API) ──► crm.lead + crm.consent_evidence
        │ optional ingress ──► integration.* ──► GO-LEAD-CANDIDATE-PROMOTION ──► crm.lead
        ▼
appointment / quote / handover mutations (franchise journey API)
        ▼
crm.appointment_transition (append-only) + sales.* durable rows
        ▼
GET /v1/customer/journey (scoped projection, max 100 rows/collection)
```

## 4. Fixtures — customer timeline

| Fixture | Path | Proves |
| --- | --- | --- |
| SQL journey seed | `db/tests/0005_franchise_customer_journey.test.sql` | customer profile, lead, appointment, quote, handover rows coexist |
| SQL history immutability | `db/tests/0015_franchise_availability_and_appointment_audit.test.sql` | `crm.appointment_transition` + immutability trigger |
| SQL timeline projection | `db/tests/0089_crm_customer_timeline_fixture.test.sql` | customer-scoped counts, transition ledger, org isolation |
| SQL lead promotion (ingress) | `db/tests/0045_lead_candidate_promotion.test.sql` | promotion row immutability; adjacent to CRM, not SDR |
| Go repository journey | `internal/platform/postgres/franchisejourney_integration_test.go` | `CustomerJourney` repeatable-read snapshot; conflict on incoherent handover links |
| Go HTTP redaction | `internal/platform/httpapi/franchisejourney_test.go` (`TestCustomerJourneyReadFailureIsRedactedAndRecoverable`) | failed projection → redacted 503, no partial leak |
| BFF confirmation + timeline | `internal/platform/httpapi/franchisejourney_test.go` (`TestAppointmentConfirmationBFFPostgres`) + `reconstruction_evidence/APPOINTMENT_CONFIRMATION_BFF_POSTGRES_V262.md` | concurrent confirm → one durable state; owner reads appointment in timeline; stranger reads zero |
| Evidence manifest | `docs/slices/CRM_HISTORY_FIXTURE_EVIDENCE.json` | machine-readable index of the above |

Synthetic fixture rules (from pack contracts):

- Tokens are **RS256/JWKS loopback**, not live IdP login.
- Timeline caps: **100 rows** per collection; HTTP body **1 MiB** max (`GO-FRANCHISE-CUSTOMER-JOURNEY-API` performance budget).
- Customer reads bind **verified subject** + **organization_id**; cross-scope returns empty or 403/409 as documented.

## 5. Verify steps (local only)

Run against a **disposable** PostgreSQL 18 database (`TEST_DATABASE_URL`); never a production DSN.

### 5.1 SQL invariants

```bash
# After all required migrations on TEST_DATABASE_URL
for f in db/tests/0003_electromobility_franchise_modules.test.sql \
         db/tests/0005_franchise_customer_journey.test.sql \
         db/tests/0015_franchise_availability_and_appointment_audit.test.sql \
         db/tests/0089_crm_customer_timeline_fixture.test.sql; do
  psql "$TEST_DATABASE_URL" -v ON_ERROR_STOP=1 -f "$f"
done
```

Optional ingress fixture (adjacent, not SDR):

```bash
psql "$TEST_DATABASE_URL" -v ON_ERROR_STOP=1 -f db/tests/0045_lead_candidate_promotion.test.sql
```

### 5.2 Go unit/integration

```bash
go test ./internal/franchisejourney/... -count=1
go test ./internal/platform/postgres -run CustomerJourney -count=1 -timeout=5m
go test ./internal/platform/httpapi -run 'TestCustomerJourneyReadFailure|TestAppointmentConfirmation' -count=1 -timeout=5m
```

### 5.3 Opt-in connected confirmation gate (timeline readback)

Requires materialized web tree + `ELITE_WEB_ROOT` absolute path. See `reconstruction_evidence/APPOINTMENT_CONFIRMATION_BFF_POSTGRES_V262.md`.

```bash
export ELITE_CONFIRMATION_E2E=1
go test ./internal/platform/httpapi -run '^TestAppointmentConfirmationBFFPostgres$' -v -count=1 -timeout=3m
unset ELITE_CONFIRMATION_E2E
```

Pass criteria:

- `crm.appointment_transition` rows are **immutable** (update/delete rejected).
- Owner customer timeline includes confirmed appointment; unrelated subject returns **zero** appointments.
- Inconsistent handover/order/stock linkage → **`ErrConflict`**, not partial JSON.

### 5.4 Library preflight (optional)

```powershell
./VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Preflight
```

Records local receipts only; does **not** authorize production.

## 6. Gaps — public CRM patterns (cite only)

These Microsoft Dynamics 365 **guidance** pages inform gaps; they do **not** admit Dynamics code, select Business Central, or certify this slice for production:

| Pattern | Official reference | Local gap |
| --- | --- | --- |
| Prospect → quote | [Prospect to quote overview](https://learn.microsoft.com/en-us/dynamics365/guidance/business-processes/prospect-to-quote-overview) | Quote acceptance, taxes, financing, and operator UX remain **CONDITIONED** beyond fixture quotes |
| Case → resolution | [Case to resolution introduction](https://learn.microsoft.com/en-us/dynamics365/guidance/business-processes/case-to-resolution-introduction) | No standalone customer-service case owner; delivery exceptions cover a narrow rejection/resolution lane only |
| Customer Service hub | [Customer Service overview](https://learn.microsoft.com/en-us/dynamics365/customer-service/overview) | No omnichannel agent desktop, SLA engine, or knowledge base product surface |
| Sales lead entity (concept) | [Create or edit leads](https://learn.microsoft.com/en-us/dynamics365/sales/create-edit-leads) | Local `crm.lead` exists; scoring, assignment rules, and live marketing connectors are project gates |
| Activity timeline (concept) | [Activity feed](https://learn.microsoft.com/en-us/dynamics365/customerengagement/on-premises/basics/set-up-interaction-timeline) | Local timeline is **scoped SQL projection**, not a unified activity stream across email/phone/tasks |
| Named SDR workflow | `markdown_system/PACK_PER_CLAIM_INDEX.md` | **FALTA** pack `*SDR*` — use ingress + promotion packs only |

Also see `markdown_system/BUSINESS_FUNCTION_OPERATING_CONTRACT_V403.md` and `roles/method-sources.v403.json` for verified Dynamics guidance queries (prospect-to-quote, case-to-resolution) with explicit “no BC selection” limits.

## 7. Explicit non-claims

- No named **SDR** or **GTM Tag Manager** pack was invented.
- `GO-OMNICHANNEL-LEAD-INGRESS` / `GO-LEAD-CANDIDATE-PROMOTION` are **ingress**, not SDR.
- `HISTORY-MODEL-TRAINING-PIPELINE` / chat corpus training is a **different** “history” contract (`markdown_system/PROJECT_HISTORY_MODEL_TRAINING_CONTRACT.md`) — out of scope here.
- A PASS on these fixtures is **not** cloud CI, live IdP, provider accounts, or production CRM.

## 8. Customer 360 read model (compose plan — not HECHO)

Architecture proposes [`markdown_system/CUSTOMER_360_READ_MODEL_PACK_PLAN.md`](../markdown_system/CUSTOMER_360_READ_MODEL_PACK_PLAN.md) as the honest closure path for the **segmentation/360 FALTA** row. That document is a **plan only**:

- **Spine (PROVEN_LOCAL):** `GET /v1/customer/journey` from `GO-FRANCHISE-CUSTOMER-JOURNEY-API` — fixtures in §4–§5 above.
- **Capture lane:** `GO-ELECTROMOBILITY-PUBLIC-CRM-API` → `crm.lead` + `crm.consent_evidence`.
- **Channel facet:** `GO-PG-CONTACT-CHANNEL-IDENTITY` → `communication.contact_channel_binding` — verify `db/tests/0047_contact_channel_identity.test.sql`; **FALTA** customer-facing merged HTTP.
- **Adjacent reads (optional BFF merge):** `GET /v1/customer/orders`, `GET /v1/customer/service-cases` (`GO-ENTERPRISE-QUERY-API`).

No elite `GET /v1/customer/360` exists @ `227c6e0`. Consumer BFF merges parallel GETs per the pack plan; Microsoft [prospect-to-quote](https://learn.microsoft.com/en-us/dynamics365/guidance/business-processes/prospect-to-quote-overview) and [case-to-resolution](https://learn.microsoft.com/en-us/dynamics365/guidance/business-processes/case-to-resolution-introduction) cite remaining gaps only.

## 9. Evidence cross-links

- `reconstruction_evidence/GO_ELECTROMOBILITY_PUBLIC_CRM_API_2026-08-24_V2.md` — public lead idempotency
- `reconstruction_evidence/FRANCHISE_CUSTOMER_JOURNEY_2026-08-29_V114.md` — journey vertical slice
- `reconstruction_evidence/GO_PG_CONTACT_CHANNEL_IDENTITY_2026-09-04_V237.md` — channel identity binding
- `reconstruction_evidence/APPOINTMENT_CONFIRMATION_BFF_POSTGRES_V262.md` — timeline readback after confirmation
- `docs/electromobility-api-runtime.md` — runtime route surface including `/v1/customer/journey`
- `markdown_system/CUSTOMER_360_READ_MODEL_PACK_PLAN.md` — Customer 360 read-model compose recipe (plan only)
