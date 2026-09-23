# Roadmap — elite-engineering-library (tip pin)

**Tip SHA:** `477358d0d7090f2c39467abe5634d81f38afae7f`  
**Scope:** biblioteca / infrastructure / franchise architecture documentation. **No** production. **No** REVESTEX product code.

**Course correction:** The sections below separate **library local-fixtures** readiness from **full franchise / company domain** coverage. The thin CAN list in [`docs/FRANCHISE_PLAYBOOK.md`](FRANCHISE_PLAYBOOK.md) is **not** the franchise product scope — see [`docs/FRANCHISE_ARCHITECTURE.md`](FRANCHISE_ARCHITECTURE.md).

---

## FRANCHISE / COMPANY DOMAINS (primary)

Full domain map, dependencies, gap research, and proposed elite additions: **[`docs/FRANCHISE_ARCHITECTURE.md`](FRANCHISE_ARCHITECTURE.md)**.

### Domain status matrix @ `477358d`

Legend: **HECHO** = selected 116/120 packs + findable reference; **PARCIAL** = conditioned, disk-only, or live gates open; **NO** = intentional GAP / no admitted pack.

| Domain | Status | Cite (paths @ tip) |
|---|---|---|
| Identity / security / authz / audit | **PARCIAL** | Selected: `implementation_packs/GO_OIDC_SERVICE_TOKEN_BROKER.md`, `GO_OIDC_PORTAL_SESSION.md`, `TYPESCRIPT_OIDC_PORTAL_ADAPTER.md`, `MICROSOFT_DEVSKIM_ADAPTED_SAST_GATE.md`, `SECURE_OPERATIONS_DELIVERY_CORE.md`, `GO_HUMAN_APPROVAL_CORE.md`; evidence `reconstruction_evidence/IDENTITY_J5_RELEASE_V402.md`; disk-only: `GO_PCI_DSS_SCOPE_CORE.md`, `GO_GDPR_CONSENT_ERASURE_CORE.md` |
| Catalog / stock / PI-PL / costs / fiscal (ARCA) | **HECHO** / **PARCIAL** live | Selected: `GO_SUPPLY_FACTORY_INVENTORY_API.md`, `GO_CONNECTED_CATALOG_*`, `GO_EXACT_FX_SNAPSHOT_ACCOUNTING.md`, `GO_FINOPS_CORE.md`, ARCA stack (`MICROSOFT_ARCA_*`, `GO_ARCA_FISCAL_ISSUANCE_API.md`, `ARCA_WSFE_UDS_WORKER.md`); `markdown_system/FRANCHISE_GAP_MAP.md` (`ARCA_INFRA` PROVEN_LOCAL); live creds **CONDITIONED** |
| Customers + history / CRM | **HECHO** | `GO_ELECTROMOBILITY_PUBLIC_CRM_API.md`, `GO_FRANCHISE_CUSTOMER_JOURNEY_API.md`, `GO_PG_CONTACT_CHANNEL_IDENTITY.md`, `GO_CUSTOMER_SURVEY_API.md` |
| Sales / quotes / orders / payments | **HECHO** / **PARCIAL** live | `GO_COMMERCE_PRICING_PAYMENT_API.md`, `GO_FRANCHISE_CUSTOMER_JOURNEY_API.md`, `GO_PAYMENT_CHECKOUT_RUNTIME.md`, `GO_OFFICIAL_PAYMENT_WEBHOOK_ADAPTERS.md`, `GO_INITIAL_HANDOVER_API.md`, `GO_RETURN_*` workers; `reconstruction_evidence/FRANCHISE_QUOTE_TO_ORDER_2026-08-29_V116.md` |
| Analytics / reporting | **PARCIAL** | Selected: `GO_DATA_ANALYTICS_CORE.md`, `GO_HTTP_METRICS_REFERENCE.md`, `PYTHON_*_ADS_REPORTING_ADAPTER.md`, `GO_CONNECTED_ROLE_METRICS_PROOF.md`; disk-only: `GO_DASHBOARDS_CORE.md`, `GO_SLO_CORE.md` |
| Multi-site / multi-sede | **HECHO** | `ELECTROMOBILITY_FRANCHISE_MODULES.md`, `GO_NETWORK_ROLE_COMPOSITION.md`, `GO_FRANCHISE_ROYALTY_SETTLEMENT_API.md`, `TYPESCRIPT_NETWORK_ROLE_PORTAL.md` |
| Channels — WhatsApp | **HECHO** / **PARCIAL** live | `PYTHON_META_WHATSAPP_CLOUD_ADAPTER.md`, `GO_CONNECTED_WHATSAPP_*`, `GO_CONNECTED_SCHEDULED_WHATSAPP.md` |
| Channels — Messenger / Instagram | **PARCIAL** | **Disk only (not selected 116):** `GO_MESSAGING_CHANNEL_ADAPTERS.md`; `internal/msgchannels/metagraph.go` (no outbound IG/Messenger claim) |
| Channels — ads / email / web leads | **PARCIAL** | Leads: `GO_OMNICHANNEL_LEAD_INGRESS.md`, Meta/TikTok lead packs; ads reporting: `PYTHON_*_ADS_REPORTING_*`; email: `GO_CHANNELS_CORE.md`, `GO_PG_OUTBOUND_DELIVERY_FENCE.md`; disk: `AWS_SES_IMMUTABLE_EMAIL_RECEIVER.md` |
| Channels — GTM Tag Manager | **NO** | GAP — no `*GTM*` pack (`markdown_system/LIBRARY_HEALTH_CHECK.md`, `PACK_PER_CLAIM_INDEX.md`) |
| Logistics / fulfillment / receiving | **HECHO** | `GO_FULFILLMENT_SERVICE_FRANCHISE_API.md`, `GO_CONNECTED_SUPPLY_CREATION.md`, `GO_CONNECTED_SERIAL_SUPPLY.md`; `docs/inventory/MICROSOFT_BC_*_DERIVATION.md` |
| Brand / public web / publish gates | **PARCIAL** | Selected storefront + gates; V403: `TYPESCRIPT_PUBLIC_WEB_EXPERIENCE_V403.md`, `PUBLIC_WEB_EXPERIENCE_PACK_PLAN_V403.md` (visual **PENDING**); disk SEO: `GO_SEO_CORE.md` |
| Admin / ops / governance / agent doors | **HECHO** / **PARCIAL** V403 UI | `START_FRANCHISE.md`, `FRANCHISE_PROJECT_OPERATING_PROTOCOL.md`, `BUSINESS_FUNCTION_OPERATING_V403.md`, `ADMIN_OPS_EXPERIENCE_PACK_PLAN_V403.md`; 48 surfaces = inventory not execution (`TOTAL_SYSTEM_CAPABILITY_CONTRACT.md`) |
| Lead ingress / promotion (no invented SDR pack) | **HECHO/PARCIAL** | `GO_OMNICHANNEL_LEAD_INGRESS.md`, `GO_LEAD_CANDIDATE_PROMOTION.md`; SDR **function** via `BUSINESS_FUNCTION_OPERATING_CONTRACT_V403.md`; named `*SDR*` pack **NO** |
| QR identity | **PARCIAL** | `implementation_packs/GO_QR_CORE.md`, `internal/qr/` — **not** in selected 116 JSON |
| QR capture | **PARCIAL** / pack **NO** | `qr_capture/` reference only; **FALTA** admitted capture pack |
| Galaxy fleet orchestration | **PARCIAL** how-tos / **NO** e2e | `docs/FRANCHISE_PLAYBOOK.md` §5 — `etor6233/grok-bot-galaxy@f8546c3`; 0 `franchise` hits |
| Conversational agent runtime | **HECHO** | `GO_CONVERSATIONAL_AGENT.md`, `GO_CONNECTED_CONVERSATION_RUNTIME.md`, `GO_APP_WIRING.md` |
| Document intelligence | **PARCIAL** | `GO_CONNECTED_DOCUMENT_REFERENCE.md`, Textract/Azure samples, T2806 PROVEN_LOCAL |
| Production / franchise live | **NO** | `qualification/FINAL_LIBRARY_READY_V402.json` (`production_authorized: false`) |

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

Nightly signed sources cited below: `qualification/FINAL_LIBRARY_READY_V402.json`, `markdown_system/LIBRARY_HEALTH_CHECK.md`, `markdown_system/PACK_PER_CLAIM_INDEX.md`, `markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md` (selected **116**; **GO-QR-CORE** still **not** in packId list), `markdown_system/FRANCHISE_COMPLETE_EXTENSION_PACK_PLAN_V403.md` (+4 → **120**), `START_FRANCHISE.md`, `START_REFERENCE_V402.md`, `LIBRARY_VS_PRODUCT_GATE_V402.md` (stub → `reconstruction_evidence/LIBRARY_VS_PRODUCT_GATE_V402.md`), `unified-experience-next/markdown_system/SUCCESSOR_AUTHORITY_ROUTER.md`, and `implementation_packs/*.md` on disk (catalog ≠ selected).

See also: [`docs/FRANCHISE_PLAYBOOK.md`](FRANCHISE_PLAYBOOK.md).

### Matrix (Nightly — infrastructure wiring only)

| Área | Estado | Cite |
|---|---|---|
| Selected composition 116→120 on disk | **HECHO** | `markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md` (116 packIds) + `markdown_system/FRANCHISE_COMPLETE_EXTENSION_PACK_PLAN_V403.md` (120); catalog `implementation_packs/*.md` on disk (do not confuse with selected 116/120) |
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
