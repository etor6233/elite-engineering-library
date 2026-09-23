# Analytics / dashboards / CRM segmentation — compose plan (not a HECHO pack)

**Revision:** `ANALYTICS-DASHBOARDS-PACK-PLAN-0.1.0`  
**Tip pin:** `5bc34a1fcf485895283ca530b2ccdc27dd8fa443`  
**Status:** **PLAN ONLY** — no new `implementation_packs/GO_ANALYTICS_*` or `GO_DASHBOARD_*` HECHO body, not added to selected 116/120 JSON, not `REBUILD_VERIFIED` as a composed pack, not production authorization.

## Architecture contract (read-only cites @ tip `5bc34a1`)

| Door | Path | Analytics row |
| --- | --- | --- |
| Domain architecture | [`docs/FRANCHISE_ARCHITECTURE.md`](../docs/FRANCHISE_ARCHITECTURE.md) § **PARCIAL → Analytics** | **PARCIAL** — selected cores + disk-only dashboards alternative |
| Build-order matrix | [`docs/ROADMAP.md`](../docs/ROADMAP.md) § **FRANCHISE / COMPANY DOMAINS** → Analytics / dashboards / CRM segmentation | **PARCIAL** — cores in selected 116; `GO_DASHBOARDS_CORE.md` disk **FALTA** select; segmentation **FALTA** |
| Tangible slice | [`docs/slices/ANALYTICS_DASHBOARDS_COMPOSE.md`](../docs/slices/ANALYTICS_DASHBOARDS_COMPOSE.md) | LOCAL_FIXTURES compose path |
| Slice evidence | [`docs/slices/ANALYTICS_DASHBOARDS_FIXTURE_EVIDENCE.json`](../docs/slices/ANALYTICS_DASHBOARDS_FIXTURE_EVIDENCE.json) | machine-readable fixture index |
| Selected composition | [`markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md`](FRANCHISE_COMPLETE_PACK_PLAN.md) | authoritative `packId` / `path` / `version` |

**Proposed addition (this file):** compose recipe over existing selected analytics cores, connected role-metrics proof, and adjacent CRM read-model plan. Does **not** close the architecture **PARCIAL** row to **HECHO**; it documents honest assembly until gap-gated packs exist.

## Claim (narrow)

Define an **analytics / dashboards / CRM-segmentation compose lane** as:

1. **Semantic ingest + metric definitions** — `GO-DATA-ANALYTICS-CORE` landing, versioned metric definitions, freshness, reconciliation (`DATA-INGEST` + `ANALYTICS-BI` contract surfaces).
2. **Operational dashboards (role-scoped)** — `GO-CONNECTED-ROLE-METRICS-PROOF` + `GO-ENTERPRISE-QUERY-API` read models exposed at `GET /v1/reporting/operations/{kind}`; portal `/dashboard` gated by `features.role_workspace=true` (`docs/ROLE_METRICS_REFERENCE.md`).
3. **CRM segmentation adjacency** — query composition documented in [`markdown_system/CUSTOMER_360_READ_MODEL_PACK_PLAN.md`](CUSTOMER_360_READ_MODEL_PACK_PLAN.md) (**PLAN ONLY**); no standalone segmentation pack, no Customer-360 HECHO endpoint.

No new durable write owner. No named `*SDR*`, `*GTM*`, or `*COTIZADOR*` pack. No promotion of disk-only `GO-DASHBOARDS-CORE` into selected 116 JSON by this plan.

## Selected packs (primary composition paths)

Exact filenames from [`markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md`](FRANCHISE_COMPLETE_PACK_PLAN.md):

| packId | version | Selected JSON `path` | Compose lane |
| --- | --- | --- | --- |
| `GO-DATA-ANALYTICS-CORE` | `0.1.0` | `implementation_packs/GO_DATA_ANALYTICS_CORE.md` | Ingest landing + semantic metrics (`analytics.*` schema); warehouse/lakehouse **SPEC_ONLY** |
| `GO-CONNECTED-ROLE-METRICS-PROOF` | `0.1.0` | `implementation_packs/GO_CONNECTED_ROLE_METRICS_PROOF.md` | Twelve operational metric families; exact-string wire; T2804 PROVEN_LOCAL proof |
| `GO-HTTP-METRICS-REFERENCE` | `0.1.16` | `implementation_packs/GO_HTTP_METRICS_REFERENCE.md` | HTTP observability reference profile (adjacent to analytics row in architecture) |
| `GO-ENTERPRISE-QUERY-API` | `0.3.0` | `implementation_packs/GO_ENTERPRISE_QUERY_API.md` | Host for `RoleMetrics` reader + reporting routes |
| `GO-ELECTROMOBILITY-APPLICATION` | `1.22.3` | `implementation_packs/GO_ELECTROMOBILITY_APPLICATION.md` | Runtime wiring (`cmd/electromobility-api`) |
| `TS-OIDC-PORTAL-ADAPTER` | `0.3.1` | `implementation_packs/TYPESCRIPT_OIDC_PORTAL_ADAPTER.md` | BFF protected GET to reporting endpoints |
| `TS-GO-API-WEB-BRIDGE` | `0.15.2` | `implementation_packs/TYPESCRIPT_GO_API_WEB_BRIDGE.md` | Typed backend client for dashboard panel |
| `TS-FRANCHISE-JOURNEY-PORTALS` | `0.21.0` | `implementation_packs/TYPESCRIPT_FRANCHISE_JOURNEY_PORTALS.md` | `/dashboard` + `RoleMetrics` UI glue |

### Foundation rows (required for honest LOCAL_FIXTURES)

| packId | version | Selected JSON `path` | Role |
| --- | --- | --- | --- |
| `PG-TX-FOUNDATION` | `0.1.0` | `implementation_packs/POSTGRES_TRANSACTIONAL_FOUNDATION.md` | tenant/org/platform idempotency |
| `GO-ENTERPRISE-BACKEND` | `0.4.8` | `implementation_packs/GO_ENTERPRISE_BACKEND_CORE.md` | HTTP/OIDC/identity shell |

Resolve filenames via selected JSON `path` only — do **not** substitute shorthand aliases (see [`docs/slices/ANALYTICS_DASHBOARDS_COMPOSE.md`](../docs/slices/ANALYTICS_DASHBOARDS_COMPOSE.md) §2).

## Disk-only honesty — `GO_DASHBOARDS_CORE.md`

| Item | Status @ `5bc34a1` | Cite |
| --- | --- | --- |
| Pack on disk | **YES** — `implementation_packs/GO_DASHBOARDS_CORE.md` (`GO-DASHBOARDS-CORE` `0.1.2`, `REBUILD_VERIFIED`) | [`reconstruction_evidence/GO_DASHBOARDS_CORE_2026-09-02_V207.md`](../reconstruction_evidence/GO_DASHBOARDS_CORE_2026-09-02_V207.md) |
| Selected 116/120 JSON | **FALTA** — absent from [`markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md`](FRANCHISE_COMPLETE_PACK_PLAN.md) | [`docs/ROADMAP.md`](../docs/ROADMAP.md) analytics row |
| Promotion by this plan | **NO** | [`docs/ROLE_METRICS_REFERENCE.md`](../docs/ROLE_METRICS_REFERENCE.md): isolated pure KPI calculator; role metrics use domain owners instead |
| Isolated scope | Pure caller-supplied KPI aggregation in `internal/dashboards/` | No domain authorization; not wired as franchise dashboard owner |

Consumers needing monthly KPI snapshots from caller-supplied inputs may materialize `GO-DASHBOARDS-CORE` in a **separate** gap-gated expediente. This plan composes **selected** cores only.

## CRM segmentation adjacency (do not fake 360 HECHO)

| Artifact | Status | Boundary |
| --- | --- | --- |
| [`markdown_system/CUSTOMER_360_READ_MODEL_PACK_PLAN.md`](CUSTOMER_360_READ_MODEL_PACK_PLAN.md) | **PLAN ONLY** | BFF parallel GETs over journey + CRM APIs; not in selected 116/120 JSON |
| Unified `GET /v1/customer/360` | **FALTA** | Compose documented; not shipped as elite HTTP |
| Standalone segmentation pack | **FALTA** | [`docs/ROADMAP.md`](../docs/ROADMAP.md) CRM + analytics rows |
| Named `*SDR*` / `*GTM*` | **FALTA** | [`markdown_system/PACK_PER_CLAIM_INDEX.md`](PACK_PER_CLAIM_INDEX.md) — ingress only via lead packs |

Segmentation for dashboards in this plan means **org-scoped operational metrics** and **customer-scoped journey reads** composed in project BFF — not a marketing-list engine.

## Owner map

```text
┌──────────────────────────────────────────────────────────────────────────┐
│ Portal BFF / /dashboard (project-owned compose — NOT an elite pack body) │
└──────────────────────────────────────────────────────────────────────────┘
         │ GET /v1/reporting/operations/{kind}?organization_id=<org>
         │ (optional CRM facet: GET /v1/customer/journey — see 360 plan)
         ▼
┌─────────────────────────┐  ┌──────────────────────────┐  ┌─────────────────────────┐
│ GO-ENTERPRISE-QUERY-API │  │ GO-DATA-ANALYTICS-CORE   │  │ CRM journey owners      │
│ RoleMetrics reader      │  │ analytics landing/metric │  │ (360 plan adjacency)    │
└──────────┬──────────────┘  └────────────┬─────────────┘  └────────────┬────────────┘
           │                              │                              │
           ▼                              ▼                              ▼
   domain tables (sales, crm,      analytics.landing_record      crm.appointment_transition
   inventory, service_ops, ...)    analytics.metric_definition   + journey projection
                                   analytics.metric_snapshot
```

| Facet | Authoritative owner | Local HTTP / schema | Dashboard compose note |
| --- | --- | --- | --- |
| Operational KPIs (12 families) | Domain owners via `GO-CONNECTED-ROLE-METRICS-PROOF` | `GET /v1/reporting/operations/{kind}` | Exact strings; max 100 groups; unavailable ≠ zero |
| Semantic metrics / ingest | `GO-DATA-ANALYTICS-CORE` | `analytics.*` DDL + Go Fresh/Reconcile | Not merged into role-metrics wire by default |
| Portal dashboard shell | `TS-FRANCHISE-JOURNEY-PORTALS` + fixture | `/dashboard` when `features.role_workspace=true` | BFF never receives SQL or client-calculated totals |
| Pure KPI calculator (disk) | `GO-DASHBOARDS-CORE` | `internal/dashboards.Build` | **Not** selected; caller-owned inputs only |
| Customer timeline (segmentation adjacency) | `GO-FRANCHISE-CUSTOMER-JOURNEY-API` | `GET /v1/customer/journey` | See 360 plan — **PLAN ONLY** merge |
| Admin ops overlay | `TYPESCRIPT_ADMIN_OPS_EXPERIENCE_V403.md` via plan | disk; **not** ∈ 116 JSON | [`markdown_system/ADMIN_OPS_EXPERIENCE_PACK_PLAN_V403.md`](ADMIN_OPS_EXPERIENCE_PACK_PLAN_V403.md) |

## Compose recipe

1. **Pin** tip `5bc34a1fcf485895283ca530b2ccdc27dd8fa443`; walk franchise doors per [`START_FRANCHISE.md`](../START_FRANCHISE.md).
2. **Select** primary rows from [`markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md`](FRANCHISE_COMPLETE_PACK_PLAN.md) and foundation rows in [`docs/slices/ANALYTICS_DASHBOARDS_COMPOSE.md`](../docs/slices/ANALYTICS_DASHBOARDS_COMPOSE.md) §2.
3. **Apply migrations** through analytics owners at minimum:
   - `0001` — `POSTGRES_TRANSACTIONAL_FOUNDATION.md`
   - `0002` — `GO_ENTERPRISE_BACKEND_CORE.md`
   - `0041` — `GO_DATA_ANALYTICS_CORE.md` (`analytics.*`)
   - `0087` — additive immutability guards (`docs/analytics/immutability-and-retention.md`)
   - domain migrations required by metric families (sales, crm, inventory, service_ops, logistics, stored value, surveys — as activated)
4. **Configure** runtime: `DATABASE_URL`, OIDC issuer/audience; enable role metrics fixture per `config/role.metrics.fixture.json` for local panel only.
5. **Expose reporting** — register `RoleMetricsModule` on enterprise query host; twelve kinds documented in [`docs/ROLE_METRICS_REFERENCE.md`](../docs/ROLE_METRICS_REFERENCE.md).
6. **Wire portal** — `/dashboard` requires session + `features.role_workspace=true`; BFF passes only `kind`, `organization_id`, and survey program id when applicable.
7. **Optional CRM facet** — for operator/customer segmentation views, follow [`markdown_system/CUSTOMER_360_READ_MODEL_PACK_PLAN.md`](CUSTOMER_360_READ_MODEL_PACK_PLAN.md) compose recipe (**PLAN ONLY** — no unified elite endpoint).
8. **Do not** materialize `GO-DASHBOARDS-CORE` into the selected 116 profile unless a separate gap-gated expediente promotes it.

## Verify steps (LOCAL_FIXTURES only)

Evidence index: [`docs/slices/ANALYTICS_DASHBOARDS_FIXTURE_EVIDENCE.json`](../docs/slices/ANALYTICS_DASHBOARDS_FIXTURE_EVIDENCE.json).

### SQL invariants (disposable `TEST_DATABASE_URL`)

```bash
psql "$TEST_DATABASE_URL" -v ON_ERROR_STOP=1 -f db/tests/0041_data_analytics.test.sql
psql "$TEST_DATABASE_URL" -v ON_ERROR_STOP=1 -f db/tests/0087_analytics_immutability.test.sql
```

### Go unit / integration

```bash
go test ./internal/analytics/... -count=1
go test ./internal/rolemetrics/... -count=1
go test ./internal/platform/postgres -run 'RoleMetrics|TestRoleMetrics' -count=1 -timeout=5m
go test ./internal/platform/httpapi -run RoleMetrics -count=1 -timeout=5m
```

### Opt-in browser connected proof

Requires `ELITE_ROLE_METRICS_BROWSER=1` and materialized web tree — see [`reconstruction_evidence/ROLE_METRICS_RELEASE_V402.md`](../reconstruction_evidence/ROLE_METRICS_RELEASE_V402.md).

```bash
export ELITE_ROLE_METRICS_BROWSER=1
go test ./internal/platform/postgres -run TestRoleMetricsBrowser -count=1 -timeout=5m
unset ELITE_ROLE_METRICS_BROWSER
```

### Library preflight (optional)

```powershell
./VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Preflight
```

Pass criteria for **this plan**:

- Analytics DDL + immutability fixtures PASS.
- Role-metrics connected scope/precision tests PASS.
- No single elite endpoint returns a merged “analytics dashboard + segmentation” JSON — compose is **documented**, not **shipped** as one HECHO pack.
- `GO-DASHBOARDS-CORE` remains disk-only relative to selected 116.

## Gaps — public patterns (cite only)

No code copy into new packs. Remaining **FALTA** informed by official guidance:

| Pattern | Official reference | Local status @ `5bc34a1` |
| --- | --- | --- |
| Monitoring / SLO dashboards | [Google SRE Book — Monitoring distributed systems](https://sre.google/sre-book/monitoring-distributed-systems/) | Role-metrics operational reads PROVEN_LOCAL; observability backend/dashboards delivery **CONDITIONED** |
| Server-side tag / web analytics | [Google server-side Tag Manager](https://developers.google.com/tag-platform/tag-manager/server-side/intro) | **FALTA** `*GTM*` pack — [`markdown_system/PACK_PER_CLAIM_INDEX.md`](PACK_PER_CLAIM_INDEX.md) |
| Marketing segmentation / lists | ROADMAP analytics + CRM rows | **FALTA** standalone segmentation pack |
| Customer 360 unified read | [`markdown_system/CUSTOMER_360_READ_MODEL_PACK_PLAN.md`](CUSTOMER_360_READ_MODEL_PACK_PLAN.md) | **PLAN ONLY** — not HECHO |
| Warehouse / lakehouse / CDC | `GO-DATA-ANALYTICS-CORE` metadata | **SPEC_ONLY** — project decision |
| Admin ops dashboard chrome | [`markdown_system/ADMIN_OPS_EXPERIENCE_PACK_PLAN_V403.md`](ADMIN_OPS_EXPERIENCE_PACK_PLAN_V403.md) | disk overlay; visual **PENDING**; not ∈ 116 JSON |
| Named SDR workflow | [`markdown_system/PACK_PER_CLAIM_INDEX.md`](PACK_PER_CLAIM_INDEX.md) | **FALTA** `*SDR*` — lead ingress only |

## Explicit non-claims

- This file is a **pack plan**, not an `implementation_packs/*` HECHO body and not added to selected 116/120 JSON.
- Does **not** invent `GO-ANALYTICS-DASHBOARDS-*`, `*SDR*`, `*GTM*`, or `*COTIZADOR*` packs.
- Does **not** claim dashboards **HECHO** — architecture row stays **PARCIAL** (`docs/ROADMAP.md`, `docs/FRANCHISE_ARCHITECTURE.md`).
- Does **not** promote `GO-DASHBOARDS-CORE` into selected JSON or invent selected-116 `packId` rows.
- Does **not** claim Customer 360 or CRM segmentation **HECHO** — adjacency via **PLAN ONLY** 360 doc.
- PASS on cited fixtures ≠ cloud CI, live IdP, provider accounts, operator BI warehouse, or production analytics.
- `production_authorized: false` per `qualification/FINAL_LIBRARY_READY_V402.json`.
- Did **not** edit `docs/ROADMAP.md` or `docs/FRANCHISE_ARCHITECTURE.md`.

## Evidence cross-links

| Receipt | Path |
| --- | --- |
| Data analytics core V179 | `reconstruction_evidence/GO_DATA_ANALYTICS_CORE_2026-09-02_V179.md` |
| Dashboards core V207 (disk-only) | `reconstruction_evidence/GO_DASHBOARDS_CORE_2026-09-02_V207.md` |
| Role metrics T2804 release | `reconstruction_evidence/ROLE_METRICS_RELEASE_V402.md`, `reconstruction_evidence/ROLE_METRICS_RELEASE_V402.json` |
| Dashboard/help boundaries | `reconstruction_evidence/DASHBOARD_HELP_BOUNDARIES_V362.md` |
| Analytics immutability guards | `docs/analytics/immutability-and-retention.md` |
| Role metrics reference | `docs/ROLE_METRICS_REFERENCE.md` |
| CRM / 360 adjacency | `markdown_system/CUSTOMER_360_READ_MODEL_PACK_PLAN.md`, `docs/slices/CRM_HISTORY_COMPOSE.md` |
| Capability catalog row | `markdown_system/CAPABILITY_CATALOG.md` (data/analytics/search) |

## Nightly verify commands (library PR checklist)

```bash
psql "$TEST_DATABASE_URL" -v ON_ERROR_STOP=1 -f db/tests/0041_data_analytics.test.sql
psql "$TEST_DATABASE_URL" -v ON_ERROR_STOP=1 -f db/tests/0087_analytics_immutability.test.sql
go test ./internal/analytics/... -count=1
go test ./internal/platform/postgres -run RoleMetricsConnected -count=1 -timeout=5m
```

```powershell
./VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Preflight
```
