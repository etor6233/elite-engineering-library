# Roadmap — elite-engineering-library (tip pin)

**Tip SHA:** `477358d0d7090f2c39467abe5634d81f38afae7f`  
**Scope:** biblioteca / infrastructure / franchise architecture documentation. **No** production. **No** REVESTEX product code.

**Course correction:** The sections below separate **library `LOCAL_FIXTURES`** readiness from **full franchise / company domain** coverage. Selected **116→120** = **full franchise-domain library compose** — not a thin trio (compose / Lead / QR). The CAN list in [`docs/FRANCHISE_PLAYBOOK.md`](FRANCHISE_PLAYBOOK.md) is **not** the franchise product scope — see [`docs/FRANCHISE_ARCHITECTURE.md`](FRANCHISE_ARCHITECTURE.md).

**Stop line:** Architecture PR = library compose claims only. **No REVESTEX** until N L says **ESTAMOS LISTOS**. `production_authorized: false` (`qualification/FINAL_LIBRARY_READY_V402.json`).

---

## FRANCHISE / COMPANY DOMAINS (primary)

Full domain map, dependencies, gap research, and proposed elite additions: **[`docs/FRANCHISE_ARCHITECTURE.md`](FRANCHISE_ARCHITECTURE.md)**.

### Nightly Domain Matrix READY @ `477358d` (PM confirmation)

Compose profile **116→120** = full franchise-domain `LOCAL_FIXTURES` — not three features. Detail: [`docs/FRANCHISE_ARCHITECTURE.md`](FRANCHISE_ARCHITECTURE.md#nightly-domain-matrix-ready--477358d-pm-confirmation).

| Status | Domains |
|---|---|
| **HECHO** | Catalog/stock · Sales/orders · Payments (+ BC exact) · Channels (compose) · Signed release · AI/conversation |
| **PARCIAL** | Identity/security · CRM · Analytics · Multi-site · Logistics · Brand/DAM · Admin UX · ARCA live |
| **NO** | GTM · Named `*SDR*` pack · Production live · V403 successor READY · Desktop V3 from Git · Galaxy franchise e2e |
| **QR** | `GO-QR-CORE` door-wired (`GO_QR_CORE.md`, `internal/qr/`) — **absent** from selected `packId` JSON; capture `qr_capture/` ref-only |

### Nightly Audit Engineer inventory @ `477358d`

| # | Status | Claim | Cite |
|---|---|---|---|
| 1 | **HECHO** | Composition 116→120; catalog ≠ selected | `markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md`, `FRANCHISE_COMPLETE_EXTENSION_PACK_PLAN_V403.md`; **214** markdown files under `implementation_packs/` (directory listing) |
| 2 | **HECHO** | V402 READY + READY_WITH_CAVEATS + `production_authorized: false` | `qualification/FINAL_LIBRARY_READY_V402.json` |
| 3 | **HECHO** | Multi-domain transactional model | `implementation_packs/ELECTROMOBILITY_FRANCHISE_MODULES.md`, `db/migrations/0003_electromobility_franchise_modules.up.sql` |
| 4 | **HECHO** | Payments | `implementation_packs/GO_BC_EXACT_AMOUNT_ADAPTER.md`, `GO_PAYMENT_CHECKOUT_RUNTIME.md`, `GO_OFFICIAL_PAYMENT_WEBHOOK_ADAPTERS.md`, `GO_EXACT_FX_SNAPSHOT_ACCOUNTING.md` |
| 5 | **HECHO** | Channels: Lead Form + Meta/TikTok + WA + ML (compose in 116; live creds **CONDITIONED**) | `implementation_packs/GO_OMNICHANNEL_LEAD_INGRESS.md`, `implementation_packs/GO_META_LEAD_WEBHOOK_SIGNAL.md`, `implementation_packs/GO_META_LEAD_EVIDENCE_IMPORT.md`, `implementation_packs/PYTHON_META_LEAD_RECONCILIATION_ADAPTER.md`, `implementation_packs/PYTHON_TIKTOK_LEAD_ADAPTER.md`, `implementation_packs/GO_TIKTOK_LEAD_DURABLE_IMPORT.md`, `implementation_packs/PYTHON_META_WHATSAPP_CLOUD_ADAPTER.md`, `implementation_packs/GO_CONNECTED_WHATSAPP_HOST.md`, `implementation_packs/GO_CONNECTED_SCHEDULED_WHATSAPP.md`, `implementation_packs/GO_CONNECTED_WHATSAPP_CAMPAIGNS.md`, `implementation_packs/GO_MERCADOLIBRE_MARKETPLACE_ADAPTER.md`, `implementation_packs/GO_MERCADOLIBRE_QUESTION_OUTBOUND.md`, `implementation_packs/GO_CONNECTED_MARKETPLACE_MUTATION.md`; `FRANCHISE_GAP_MAP.md` T2805 |
| 6 | **HECHO** | Signed release / verify | `implementation_packs/PORTABLE_SIGNED_RELEASE_EVIDENCE_GATE.md`, `START_REFERENCE_V402.md` |
| 7 | **HECHO** / **PARCIAL** keys | AI / conversation cluster | `implementation_packs/GO_CONVERSATIONAL_AGENT.md`, `GO_CONVERSATION_TOOLS.md`, `GO_CONNECTED_CONVERSATION_RUNTIME.md`, `GO_APP_WIRING.md`, `GO_AGENT_DOMAIN_BINDING.md`, `GO_OPENAI_RESPONSES_TOOL_ADAPTER.md` |
| 8 | **PARCIAL** | QR identity door-wired; absent from selected JSON | `implementation_packs/GO_QR_CORE.md`, `internal/qr/` |
| 9 | **PARCIAL** / **FALTA** pack | QR capture ref-only | `qr_capture/` |
| 10 | **PARCIAL** | Dashboards / CRM segmentation / brand DAM / admin UX | `implementation_packs/GO_DATA_ANALYTICS_CORE.md`, `GO_CONNECTED_ROLE_METRICS_PROOF.md`, `GO_DASHBOARDS_CORE.md`, `TYPESCRIPT_ADMIN_OPS_EXPERIENCE_V403.md`, `markdown_system/ADMIN_OPS_EXPERIENCE_PACK_PLAN_V403.md`, `TYPESCRIPT_PUBLISHED_CATALOG_STOREFRONT.md`; **additive-only** `TYPESCRIPT_PUBLIC_WEB_EXPERIENCE_V403.md` via `markdown_system/PUBLIC_WEB_EXPERIENCE_PACK_PLAN_V403.md` (not ∈ extension +4); no dedicated DAM pack |
| 11 | **PARCIAL** | ARCA PROVEN_LOCAL; live creds pending | `implementation_packs/MICROSOFT_ARCA_WSAA_CREDENTIAL_CORE.md`, `MICROSOFT_ARCA_WSFE_GENERATED_CLIENT.md`, `MICROSOFT_ARCA_WSFE_SOAP_ADAPTER.md`, `GO_ARCA_FISCAL_ISSUANCE_API.md`, `ARCA_WSFE_UDS_WORKER.md`; `FRANCHISE_GAP_MAP.md` (`ARCA_INFRA`) |
| 12 | **PARCIAL** | Multi-site compose; site CMS/IaC thin; Galaxy = field guide | `implementation_packs/GO_NETWORK_ROLE_COMPOSITION.md`, `GO_FRANCHISE_ROYALTY_SETTLEMENT_API.md`; Galaxy §5 playbook |
| 13 | **NO** | GTM Tag Manager | No GTM Tag Manager pack on disk; Lead Form ≠ GTM |
| 14 | **NO** | Named SDR workflow pack | Function: `implementation_packs/GO_OMNICHANNEL_LEAD_INGRESS.md`, `GO_LEAD_CANDIDATE_PROMOTION.md`, `markdown_system/BUSINESS_FUNCTION_OPERATING_CONTRACT_V403.md` |
| 15 | **NO** | Production live / V403 successor READY / Desktop V3 / Galaxy e2e | `FINAL_LIBRARY_READY_V402.json`, `SUCCESSOR_AUTHORITY_ROUTER.md`, `docs/ROADMAP.md` |

**PARCIAL public patterns (cite only):** [NIST 800-63B](https://pages.nist.gov/800-63-3/sp800-63b.html); [OIDC Core](https://openid.net/specs/openid-connect-core-1_0.html); [Google sGTM](https://developers.google.com/tag-platform/tag-manager/server-side/intro); [Google SRE Book](https://sre.google/sre-book/monitoring-distributed-systems/); [AFIP WSFE](https://www.afip.gob.ar/ws/documentacion/ws-factura-electronica.asp); [Amazon SP-API Fulfillment](https://developer-docs.amazon.com/sp-api/docs/fulfillment-outbound-api); [AWS Control Tower](https://docs.aws.amazon.com/controltower/latest/userguide/what-is-control-tower.html); Material / Fluent / Spectrum design systems.

### Domain status matrix @ `477358d`

Legend: **HECHO** = selected 116/120 packs + findable reference; **PARCIAL** = conditioned, disk-only, or live gates open; **NO** = intentional GAP / no admitted pack.

| Domain | Status | Cite (paths @ tip) |
|---|---|---|
| Identity / security / authz / audit | **PARCIAL** | Selected: `implementation_packs/GO_OIDC_SERVICE_TOKEN_BROKER.md`, `GO_OIDC_PORTAL_SESSION.md`, `TYPESCRIPT_OIDC_PORTAL_ADAPTER.md`, `MICROSOFT_DEVSKIM_ADAPTED_SAST_GATE.md`, `SECURE_OPERATIONS_DELIVERY_CORE.md`, `GO_HUMAN_APPROVAL_CORE.md`; evidence `reconstruction_evidence/IDENTITY_J5_RELEASE_V402.md`; disk-only: `GO_PCI_DSS_SCOPE_CORE.md`, `GO_GDPR_CONSENT_ERASURE_CORE.md` |
| Fiscal (ARCA live) | **PARCIAL** | ARCA stack in selected 116; `FRANCHISE_GAP_MAP.md` (`ARCA_INFRA` PROVEN_LOCAL); live creds pending — [AFIP WSFE](https://www.afip.gob.ar/ws/documentacion/ws-factura-electronica.asp) |
| Catalog / stock | **HECHO** | `implementation_packs/ELECTROMOBILITY_FRANCHISE_MODULES.md`, `GO_SUPPLY_FACTORY_INVENTORY_API.md`, `GO_CONNECTED_CATALOG_AUTHORING.md`, `GO_CONNECTED_CATALOG_PUBLICATION.md`, `GO_CONNECTED_SERIAL_SUPPLY.md` |
| Customers + history / CRM | **PARCIAL** | `GO_ELECTROMOBILITY_PUBLIC_CRM_API.md`, `GO_FRANCHISE_CUSTOMER_JOURNEY_API.md`, `GO_PG_CONTACT_CHANNEL_IDENTITY.md` — compose; segmentation/360 **FALTA** |
| Sales / quotes / orders | **HECHO** | `GO_FRANCHISE_CUSTOMER_JOURNEY_API.md`, `GO_COMMERCE_PRICING_PAYMENT_API.md`, `GO_INITIAL_HANDOVER_API.md`, `GO_RETURN_EFFECT_EXECUTION_WORKER.md`, `GO_OFFICIAL_RETURN_REFUND_WORKER.md`, `GO_RETURN_EXCHANGE_FULFILLMENT_WORKER.md`, `GO_RETURN_ACCOUNTING_REVERSAL_WORKER.md`, `GO_RETURN_FISCAL_CREDIT_NOTE_WORKER.md`; `reconstruction_evidence/FRANCHISE_QUOTE_TO_ORDER_2026-08-29_V116.md` |
| Payments (+ BC exact) | **HECHO** | `GO_BC_EXACT_AMOUNT_ADAPTER.md`, `GO_PAYMENT_CHECKOUT_RUNTIME.md`, `GO_OFFICIAL_PAYMENT_WEBHOOK_ADAPTERS.md`, `GO_EXACT_FX_SNAPSHOT_ACCOUNTING.md` |
| Analytics / dashboards / CRM segmentation | **PARCIAL** | Cores: `GO_DATA_ANALYTICS_CORE.md`, `GO_CONNECTED_ROLE_METRICS_PROOF.md`; disk **FALTA** select: `GO_DASHBOARDS_CORE.md`; segmentation **FALTA**; tangible slice @ `4503540`: [`docs/slices/ANALYTICS_DASHBOARDS_COMPOSE.md`](slices/ANALYTICS_DASHBOARDS_COMPOSE.md), [`docs/slices/ANALYTICS_DASHBOARDS_FIXTURE_EVIDENCE.json`](slices/ANALYTICS_DASHBOARDS_FIXTURE_EVIDENCE.json) |
| Brand / DAM | **PARCIAL** / DAM **FALTA** | `TYPESCRIPT_PUBLISHED_CATALOG_STOREFRONT.md` (selected 116); **additive** `TYPESCRIPT_PUBLIC_WEB_EXPERIENCE_V403.md` via `markdown_system/PUBLIC_WEB_EXPERIENCE_PACK_PLAN_V403.md` (not extension +4); **FALTA** dedicated DAM pack |
| Admin UX | **PARCIAL** | `TYPESCRIPT_ADMIN_OPS_EXPERIENCE_V403.md` + plan on disk; **not** in 116/120 JSON — gap-gate |
| Multi-site / multi-sede | **PARCIAL** | Hierarchy/royalty compose: `GO_NETWORK_ROLE_COMPOSITION.md`, `GO_FRANCHISE_ROYALTY_SETTLEMENT_API.md`; site CMS/IaC thin; Galaxy = field guide only |
| Channels (WA / leads / Meta / TikTok / ML) | **HECHO** compose / **PARCIAL** live | `implementation_packs/GO_OMNICHANNEL_LEAD_INGRESS.md`, `GO_META_LEAD_WEBHOOK_SIGNAL.md`, `GO_META_LEAD_EVIDENCE_IMPORT.md`, `PYTHON_META_WHATSAPP_CLOUD_ADAPTER.md`, `GO_CONNECTED_WHATSAPP_HOST.md`, `GO_CONNECTED_SCHEDULED_WHATSAPP.md`, `GO_CONNECTED_WHATSAPP_CAMPAIGNS.md`, `GO_MERCADOLIBRE_MARKETPLACE_ADAPTER.md`, `GO_MERCADOLIBRE_QUESTION_OUTBOUND.md`, `GO_CONNECTED_MARKETPLACE_MUTATION.md`; live creds **CONDITIONED** |
| Channels — Messenger / Instagram | **PARCIAL** | **Disk only (not selected 116):** `GO_MESSAGING_CHANNEL_ADAPTERS.md`; `internal/msgchannels/metagraph.go` (no outbound IG/Messenger claim) |
| Channels — ads / email / web leads | **PARCIAL** | Leads: `GO_OMNICHANNEL_LEAD_INGRESS.md`, `GO_META_LEAD_WEBHOOK_SIGNAL.md`, `PYTHON_TIKTOK_LEAD_ADAPTER.md`; ads reporting: `PYTHON_GOOGLE_ADS_REPORTING_ADAPTER.md`, `PYTHON_META_ADS_REPORTING_ADAPTER.md`, `PYTHON_TIKTOK_ADS_REPORTING_ADAPTER.md`; email: `GO_CHANNELS_CORE.md`, `GO_PG_OUTBOUND_DELIVERY_FENCE.md`; disk: `AWS_SES_IMMUTABLE_EMAIL_RECEIVER.md` |
| Channels — GTM Tag Manager | **NO** | GAP — no `*GTM*` pack (`markdown_system/LIBRARY_HEALTH_CHECK.md`, `PACK_PER_CLAIM_INDEX.md`) |
| Logistics / fulfillment / receiving | **PARCIAL** | `GO_FULFILLMENT_SERVICE_FRANCHISE_API.md`, `GO_CONNECTED_SUPPLY_CREATION.md`, `GO_CONNECTED_SERIAL_SUPPLY.md`; compose PROVEN_LOCAL; DOM/live carrier thin; tangible slice @ `0397b2f`: [`docs/slices/LOGISTICS_FULFILLMENT_COMPOSE.md`](slices/LOGISTICS_FULFILLMENT_COMPOSE.md), [`docs/slices/LOGISTICS_FULFILLMENT_FIXTURE_EVIDENCE.json`](slices/LOGISTICS_FULFILLMENT_FIXTURE_EVIDENCE.json) |
| Brand / public web / publish gates | **PARCIAL** | Selected 116: `TYPESCRIPT_PUBLISHED_CATALOG_STOREFRONT.md`, `GOOGLE_LIGHTHOUSE_WEB_QUALITY_GATE.md`, `MICROSOFT_PLAYWRIGHT_BROWSER_GATE.md`; extension +4: `TYPESCRIPT_FRANCHISE_EXPERIENCE_V403.md`, `TS_DESIGN_SYSTEM_V403.md`; **additive-only** `TYPESCRIPT_PUBLIC_WEB_EXPERIENCE_V403.md` via `markdown_system/PUBLIC_WEB_EXPERIENCE_PACK_PLAN_V403.md` (not ∈120 +4; visual **PENDING**); disk SEO: `GO_SEO_CORE.md` |
| Admin / ops / governance / agent doors | **HECHO** doors / **PARCIAL** admin UX | `START_FRANCHISE.md`, `FRANCHISE_PROJECT_OPERATING_PROTOCOL.md`, `BUSINESS_FUNCTION_OPERATING_V403.md` (in 120); admin overlay separate expediente |
| Lead ingress / promotion | **PARCIAL** / **NO** named SDR pack | `GO_OMNICHANNEL_LEAD_INGRESS.md`, `GO_LEAD_CANDIDATE_PROMOTION.md` (channels compose); SDR function via `BUSINESS_FUNCTION_OPERATING_CONTRACT_V403.md`; named `*SDR*` pack **FALTA** |
| QR identity | **PARCIAL** | `implementation_packs/GO_QR_CORE.md`, `internal/qr/` — **not** in selected 116 JSON |
| QR capture | **PARCIAL** / pack **NO** | `qr_capture/` reference only; **FALTA** admitted capture pack |
| Galaxy fleet orchestration | **PARCIAL** how-tos / **NO** e2e | `docs/FRANCHISE_PLAYBOOK.md` §5 — `etor6233/grok-bot-galaxy@f8546c3`; 0 `franchise` hits |
| AI / conversation | **HECHO** | `GO_CONVERSATIONAL_AGENT.md`, `GO_CONNECTED_CONVERSATION_RUNTIME.md`, `GO_APP_WIRING.md` (+ tools); keys **CONDITIONED** |
| Signed release | **HECHO** | `PORTABLE_SIGNED_RELEASE_EVIDENCE_GATE.md`, `START_REFERENCE_V402.md` |
| Document intelligence | **PARCIAL** | `GO_CONNECTED_DOCUMENT_REFERENCE.md`, Textract/Azure samples, T2806 PROVEN_LOCAL |
| Production live | **NO** | `qualification/FINAL_LIBRARY_READY_V402.json` (`production_authorized: false`) |
| V403 successor READY | **NO** | `unified-experience-next/markdown_system/SUCCESSOR_AUTHORITY_ROUTER.md` (IN_PROGRESS) |
| Desktop V3 from Git | **NO** | `docs/ROADMAP.md` (UNKNOWN_NEVER_PUSHED) |
| Galaxy franchise e2e | **NO** | `docs/FRANCHISE_PLAYBOOK.md` §5 — 0 `franchise` hits at pin |

### Recommended build order (summary)

| Phase | Focus |
|---|---|
| 0 | Consumer gates — `START_FRANCHISE` bridge, pack plan 116 or 120 |
| 1 | Platform — PG + backend + OIDC + secure ops + metrics |
| 2 | Org / network — franchise modules, roles, policy |
| 3 | Catalog + supply — inventory, authoring, publication |
| 4 | CRM + leads — ingress, promotion, channel identity |
| 5 | Commercial — quote, order, pay, handover |
| 6 | Fulfillment + returns |
| 7 | Fiscal (ARCA infra → live creds) |
| 8 | Channels live (WA, Meta, ads) |
| 9 | Brand / public web (V403 overlay when approved) |
| 10 | Analytics, training, help |
| 11 | Explicit gaps — GTM, QR capture, named SDR, Galaxy |

Detail per phase: [`docs/FRANCHISE_ARCHITECTURE.md`](FRANCHISE_ARCHITECTURE.md#recommended-build-order-franchise-consumer).

### Honest bounds (franchise scope)

- Selected **116→120** = full franchise-domain **LOCAL_FIXTURES** compose — not “three features”.
- Catalog on disk (**214** packs) ≠ selected **116** ≠ V403 **120**.
- Lead Form (`GO-OMNICHANNEL-LEAD-INGRESS`) ≠ GTM Tag Manager.
- SDR **function** ≠ named `*SDR*` pack — do not invent.
- `qr_capture/` ≠ admitted pack; `GO_QR_CORE` ≠ selected 116.
- Galaxy pin = orchestration KB — not franchise e2e or installable IaC.
- `READY_FOR_LIBRARY_USE` ≠ franchise product readiness.

---

## LIBRARY_INFRASTRUCTURE (Nightly local-fixtures)

Re-sign Nightly: **READY_WITH_CAVEATS** for `LIBRARY_INFRASTRUCTURE` / local-fixtures only (`qualification/FINAL_LIBRARY_READY_V402.json`), tip `477358d`. Not production. Not “every orphan in selected JSON”.

Nightly caveats (delta vs READY_RECONFIRMED wording):

1. QR identity is door-wired (`GO-QR-CORE` + `internal/qr/`) but **not** in selected 116/120 packId JSON.
2. Intentional **FALTA**: GTM Tag Manager pack, named SDR pack, QR capture admitted pack.
3. `production_authorized: false` in FINAL_LIBRARY_READY + generic preflight remains **BLOCKED** (local fixtures only).

Nightly signed sources cited below: `qualification/FINAL_LIBRARY_READY_V402.json`, `markdown_system/LIBRARY_HEALTH_CHECK.md`, `markdown_system/PACK_PER_CLAIM_INDEX.md`, `markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md` (selected **116**; **GO-QR-CORE** still **not** in packId list), `markdown_system/FRANCHISE_COMPLETE_EXTENSION_PACK_PLAN_V403.md` (+4 → **120**; extension +4 only: `BUSINESS_FUNCTION_OPERATING_V403.md`, `TS_DESIGN_SYSTEM_V403.md`, `TYPESCRIPT_FRANCHISE_EXPERIENCE_V403.md`, `FRANCHISE_CLOUD_EXECUTION_V403.md`; **additive** `TYPESCRIPT_PUBLIC_WEB_EXPERIENCE_V403.md` via `markdown_system/PUBLIC_WEB_EXPERIENCE_PACK_PLAN_V403.md` is **not** one of those four), `START_FRANCHISE.md`, `START_REFERENCE_V402.md`, `LIBRARY_VS_PRODUCT_GATE_V402.md` (stub → `reconstruction_evidence/LIBRARY_VS_PRODUCT_GATE_V402.md`), `unified-experience-next/markdown_system/SUCCESSOR_AUTHORITY_ROUTER.md`, and **214** markdown files under `implementation_packs/` (directory listing; catalog ≠ selected).

See also: [`docs/FRANCHISE_PLAYBOOK.md`](FRANCHISE_PLAYBOOK.md).

### Matrix (Nightly — infrastructure wiring only)

| Área | Estado | Cite |
|---|---|---|
| Selected composition 116→120 on disk | **HECHO** | `markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md` (116 packIds) + `markdown_system/FRANCHISE_COMPLETE_EXTENSION_PACK_PLAN_V403.md` (120); **214** markdown files under `implementation_packs/` (directory listing; do not confuse with selected 116/120) |
| V402 library READY receipt | **HECHO** | `qualification/FINAL_LIBRARY_READY_V402.json` (`READY_FOR_LIBRARY_USE` / `LIBRARY_INFRASTRUCTURE`; `production_authorized: false`) |
| Gates/doors START_* + LIBRARY_VS_PRODUCT | **HECHO** | `START_FRANCHISE.md`, `START_REFERENCE_V402.md`, root stub `LIBRARY_VS_PRODUCT_GATE_V402.md` → `reconstruction_evidence/LIBRARY_VS_PRODUCT_GATE_V402.md` |
| SUCCESSOR authority door path | **PARCIAL** | `unified-experience-next/markdown_system/SUCCESSOR_AUTHORITY_ROUTER.md` (V403 NEXT = **IN_PROGRESS**; do not inherit READY) |
| Google Lead Form | **HECHO** | `implementation_packs/GO_OMNICHANNEL_LEAD_INGRESS.md` (+ selected composition) |
| Lead promotion / Meta / TikTok | **HECHO/PARCIAL** | Selected packs e.g. `GO_LEAD_CANDIDATE_PROMOTION.md`, Meta/TikTok lead packs; live creds **CONDITIONED** |
| QR identity | **PARCIAL** | `implementation_packs/GO_QR_CORE.md` + `internal/qr/`; **not** in selected pack JSON (GO-QR-CORE absent from FRANCHISE_COMPLETE packId list) |
| QR capture | **PARCIAL** | `qr_capture/` ref-only; **FALTA** admitted capture pack |
| GTM Tag Manager | **NO** | GAP; fuera — no `*GTM*` pack (`LIBRARY_HEALTH_CHECK.md`) |
| Named SDR pack | **NO** | **FALTA**; función **PARCIAL** vía lead packs — no inventar |
| ARCA fiscal infra | **PARCIAL** | PROVEN_LOCAL infra/fixtures; live creds pending |
| Signed release / verify | **HECHO** | `START_REFERENCE_V402.md` + portable release / verify gates |
| Production / franchise live | **NO** | `production_authorized: false` in FINAL_LIBRARY_READY receipt |
| Desktop V3 via Git | **NO** | UNKNOWN_NEVER_PUSHED — do not assume recoverable from Git |
| V403 successor closed READY | **NO** | SUCCESSOR = **IN_PROGRESS** |
| Galaxy franchise fleet orchestration | **PARCIAL** (how-tos / demo fleets) / **NO** for franchise e2e | `etor6233/grok-bot-galaxy@f8546c3` — 0 `franchise` hits; no provisioning IaC; Duplicate/Marketplace ≠ e2e |

### Honest bounds (infrastructure)

- Catalog pack markdown on disk ≠ selected 116 ≠ V403 extension 120.
- Lead Form ≠ GTM Tag Manager.
- `qr_capture/` ≠ selected admitted pack; identity uses `GO_QR_CORE` (CONDITIONED / not in selected JSON).
- Do not invent a named SDR pack or GTM pack.
- Galaxy pin = orchestration field guide only — not franchise e2e, not installable packs/IaC.
- Stop line for product work: see playbook — no REVESTEX product until N L says **ESTAMOS LISTOS**.
