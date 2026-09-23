# Library health check — agent & franchise wiring

**Commit audited:** `57f77b03773fbbaeb160ec7ee68bf19ea0f9befe` (`fix: restore library index as repo front door`)

**Scope:** entrypoints, indexes, cross-links, and domain findability. Not a production or cloud certification.

## How to use this doc

| Need | Start here |
|---|---|
| Agent router (all hosts) | `AGENTS.md` |
| One-page use / resume | `markdown_system/USE_LIBRARY_V403.md` |
| Persona (3 graphs) | `markdown_system/LIBRARY_HUMAN_GRAPH.md` |
| Find any path | `rg -n -F "<término>" markdown_system/LIBRARY_SEARCH_INDEX.md` |
| One pack per claim | `markdown_system/PACK_PER_CLAIM_INDEX.md` |
| New franchise from library copy | `START_FRANCHISE.md` → `markdown_system/FRANCHISE_PROJECT_OPERATING_PROTOCOL.md` |
| New any project (portable prompt) | `START_ANY_PROJECT.md` |
| Franchise gaps / preflight | `markdown_system/FRANCHISE_GAP_MAP.md`, `markdown_system/FRANCHISE_PREFLIGHT_GAP.md` |
| Rebuild local V403 reference | `markdown_system/START_V403_LOCAL.md` |

## Entrypoint inventory (@ 7dbf39c)

| Surface | Path | Role |
|---|---|---|
| Universal agent index | `AGENTS.md` | Compact router; authority table; mandatory project routing |
| Historical full router | `AGENT_SYSTEM_START.md` | Section router — load slices on demand, not whole file |
| Claude | `CLAUDE.md` | → `AGENTS.md` + `USE_LIBRARY_V403.md` |
| Grok | `GROK_AGENT_ENTRY.md` | Progressive disclosure → skill + `AGENTS.md` |
| Codex skill (bridge) | `.agents/skills/elite-engineering-library/SKILL.md` | Pointer to canonical Grok skill |
| Grok skill (canonical) | `.grok/skills/elite-engineering-library/SKILL.md` | Progressive entry workflow |
| Project bridge installer | `INSTALL_AGENT_BRIDGE.ps1` | Installs managed blocks + skills when library ≠ project root |
| Franchise start | `START_FRANCHISE.md` | Library vs product separation; operating bridge materialization |
| Any-project prompt | `START_ANY_PROJECT.md` | Copy-paste first instruction for consumer projects |
| Pack-per-claim | `markdown_system/PACK_PER_CLAIM_INDEX.md` | One pack per bounded claim |
| Search index | `markdown_system/LIBRARY_SEARCH_INDEX.md` | Every indexed path; `rg` target |
| Human graph | `markdown_system/LIBRARY_HUMAN_GRAPH.md` | Three mermaid graphs for people |
| Public agent methods | `markdown_system/ELITE_PUBLIC_AGENT_METHODS.md` | Spec Kit + AI-DLC + Skills comparison |
| Architecture packs index | `architecture_packs/README.md` | Six `CANDIDATE_PACK` architecture docs |
| Implementation packs index | `implementation_packs/README.md` | Materializable packs by domain |
| This health report | `markdown_system/LIBRARY_HEALTH_CHECK.md` | Wiring + domain matrix + honest GAPs |

## Cross-link audit (navigation layer)

Checked relative links in entrypoints and navigation files listed above.

| Check | Result |
|---|---|
| Entrypoint → index → pack paths | **PASS** — 31 relative links, 0 broken in navigation layer |
| `START_FRANCHISE` → protocol / V403 extension | **PASS** |
| `USE_LIBRARY_V403` → franchise / capability tables | **PASS** |
| `PACK_PER_CLAIM_INDEX` → `implementation_packs/*.md` targets | **PASS** for listed rows |
| Embedded code inside materialized packs | **Not audited** — PowerShell/Dapr/Hugo link syntax inside pack bodies is not navigation |

### Orphan / weak-index notes

| Item | Status |
|---|---|
| `MATERIALIZATION_RECORD.md` (repo root) | Listed in search index prose count but not as a dedicated line before this polish — add if agents need it |
| `README.md` (repo root) | **Library index** — agents and visitors start at root README → `AGENTS.md` / `USE_LIBRARY_V403` / this file |
| `docs/ENTERPRISE_WEB_BFF.md` | BFF operational readme relocated from root; linked from library README |
| `AGENTS.md` @ 7dbf39c (pre-fix) | **Corrupted** (mojibake, single line) — **restored** in this PR from `1ce96be` + health wiring |
| `GROK_AGENT_ENTRY.md` @ 7dbf39c (pre-fix) | **Truncated** — **restored** persona/search lines from `1ce96be` |
| `qr_capture/`, `db/migrations/`, `src/platform/seo/` | Were **missing from search index** before this polish — added under **Código de referencia local** |
| `implementation_packs/GO_QR_CORE.md`, `internal/qr/` | **Orphan vs selected doors** — pack + Go owner existed on tip but HEALTH/PACK_PER_CLAIM still marked QR as GAP citing only `qr_capture/`; wired into domain matrix + claim index (identity vs capture split) |
| `implementation_packs/GO_LEAD_CANDIDATE_PROMOTION.md`, `GO_OMNICHANNEL_LEAD_INGRESS.md` | Present on tip and in search index; cited under Latency/SDR *function* labeling (still **FALTA** a pack named `*SDR*`) |
| `markdown_system/USE_LIBRARY_V403.md`, `LIBRARY_HUMAN_GRAPH.md`, `START_V403_LOCAL.md` | Referenced in prose but **not listed** in index sections — added |

## Domain coverage matrix (@ 7dbf39c)

Legend: **HECHO** = admitted pack/plan + code or contract path findable; **PARCIAL** = conditioned, fixture-only, or function contract without dedicated pack; **GAP** = no admitted pack filename or explicit blocker.

| Domain | Status | Primary paths | Notes |
|---|---|---|---|
| **Backend** (Go enterprise, PostgreSQL, workers) | **HECHO** | `implementation_packs/GO_ENTERPRISE_BACKEND_CORE.md`, `GO_RELIABLE_ASYNC_WORKERS.md`, `POSTGRES_TRANSACTIONAL_FOUNDATION.md`, `db/migrations/` (178 SQL files), `markdown_system/ENTERPRISE_BACKEND_PACK_PLAN.md` | `REBUILD_VERIFIED / CONDITIONED`; live target still required |
| **Security** (SAST, secrets, secure delivery) | **HECHO** | `implementation_packs/MICROSOFT_DEVSKIM_ADAPTED_SAST_GATE.md`, `SECURE_OPERATIONS_DELIVERY_CORE.md`, `markdown_system/ZERO_COST_VULNERABILITY_MONITORING_PROFILE.md`, `PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD.md` | OpenGrep/Cosign path conditioned; not full offensive security |
| **Latency / SDR** (SLI, metrics, sales-dev response) | **PARCIAL** | Latency: `implementation_packs/GO_HTTP_METRICS_REFERENCE.md`, `GO_OBSERVABILITY_CORE.md`, `COMPUTER_SYSTEMS_PERFORMANCE_LOW_LATENCY.md`, `markdown_system/HTTP_METRICS_REFERENCE_PACK_PLAN.md`, `reconstruction_evidence/ENTRY_LATENCY_AND_INTERRUPTION_V326.md`. SDR *function* (not a `*SDR*` pack): `markdown_system/BUSINESS_FUNCTION_OPERATING_CONTRACT_V403.md` → `implementation_packs/GO_LEAD_CANDIDATE_PROMOTION.md`, `implementation_packs/GO_OMNICHANNEL_LEAD_INGRESS.md` | HTTP SLI reference is synthetic/local; **FALTA** dedicated `*SDR*` workflow pack — use lead ingress + promotion packs already on tip; gap gate if a named SDR pack is REQUIRED |
| **SEO** | **PARCIAL** | `implementation_packs/GO_SEO_CORE.md`, `src/app/sitemap.ts`, `src/app/robots.ts`, `src/platform/seo/public-indexing.ts`, `implementation_packs/GOOGLE_LIGHTHOUSE_WEB_QUALITY_GATE.md` | Go SEO core explicitly excludes sitemap XML certification and Google Search Console; TS web bridge composes public metadata |
| **GTM** (Google Tag Manager / web tags) | **GAP** | `markdown_system/PACK_PER_CLAIM_INDEX.md` (no row), ads/reporting only: `GOOGLE_ADS_REPORTING_PACK_PLAN.md`, `META_ADS_REPORTING_PACK_PLAN.md`, `TIKTOK_ADS_REPORTING_PACK_PLAN.md` | **No `*GTM*` pack plan or implementation pack** at this commit; do not invent — use gap gate if REQUIRED |
| **Frontend** (BFF, portals, CSP, Playwright) | **HECHO** | `implementation_packs/TYPESCRIPT_GO_API_WEB_BRIDGE.md`, `TYPESCRIPT_FRANCHISE_JOURNEY_PORTALS.md`, `MICROSOFT_PLAYWRIGHT_BROWSER_GATE.md`, `FRONTEND_PRODUCT_ENGINEERING_UX.md`, `src/` | Production IdP/edge/CDN remain project gates |
| **Agents** (bridge, skills, conversational runtime) | **HECHO** | `INSTALL_AGENT_BRIDGE.ps1`, `GROK_AGENT_ENTRY.md`, `.grok/skills/elite-engineering-library/SKILL.md`, `implementation_packs/GO_CONVERSATIONAL_AGENT.md`, `GO_CONNECTED_CONVERSATION_RUNTIME.md`, `markdown_system/ELITE_PUBLIC_AGENT_METHODS.md`, `markdown_system/AGENT_AUTONOMY_CONTRACT.md` | Platform execution `NOT_TESTED` unless proven on that host |
| **QR / mobile capture** | **PARCIAL** | **Identity (admitted pack):** `implementation_packs/GO_QR_CORE.md` + `internal/qr/` (payload/checksum/tenant Verify). **Capture (reference code only):** `qr_capture/` (`worker.py`, `client.py`, `browser-capture.mjs`, …) | Split already on tip: wire identity via `GO_QR_CORE`; treat `qr_capture/` as local reference until a capture pack is admitted. Do not invent a second QR pack. `START_ANY_PROJECT.md` still documents rejected AWS mobile receipt |
| **Franchise planning** | **HECHO** | `START_FRANCHISE.md`, `FRANCHISE_ACCELERATOR.md`, `markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md`, `FRANCHISE_COMPLETE_EXTENSION_PACK_PLAN_V403.md`, `FRANCHISE_GAP_MAP.md`, `FRANCHISE_PREFLIGHT_GAP.md`, `FRANCHISE_PROJECT_OPERATING_PROTOCOL.md` | V402 ZIP bytes ≠ live Git tree; cloud extension is separate expediente |

### Backend ↔ security ↔ latency ↔ frontend ↔ agents (wiring)

```text
AGENTS.md / skills
  → PACK_PER_CLAIM_INDEX (one claim)
  → implementation_packs/* (materialize)
  → db/migrations + src/ (reference tree when using V403 body)
  → SECURE_OPERATIONS_DELIVERY_CORE + DEVSKIM gate (security lint)
  → GO_HTTP_METRICS_REFERENCE (latency SLI — conditioned)
  → TYPESCRIPT_GO_API_WEB_BRIDGE + Playwright/Lighthouse gates (frontend)
  → GO_CONVERSATIONAL_AGENT / CONNECTED_CONVERSATION_RUNTIME (agents)
```

Franchise path adds: `START_FRANCHISE` → `PROJECT_OPERATING_CONNECTION` pack → consumer `PROJECT_AGENT_ENTRY.md` → same pack-per-claim loop.

## Index maintenance (@ 7dbf39c)

| Index | Staleness found | Action |
|---|---|---|
| `LIBRARY_SEARCH_INDEX.md` | Missing agent-use docs + local code modules | **Updated in this PR** — added `USE_LIBRARY_V403`, `LIBRARY_HUMAN_GRAPH`, `START_V403_LOCAL`, `LIBRARY_HEALTH_CHECK`, and **Código de referencia local** |
| `LIBRARY_HUMAN_GRAPH.md` | No health/wiring pointer | **Updated in this PR** — link to this file |
| `PACK_PER_CLAIM_INDEX.md` | SEO / latency / SDR / GTM / QR rows present; QR row lagged tip assets | **Wired** — QR identity → `GO_QR_CORE` + `internal/qr/`; capture stays `qr_capture/` reference; GTM remains GAP; SDR remains PARCIAL via lead packs |

### Regeneration commands (do not invent automation)

| Tool | Command | When |
|---|---|---|
| XERJ search node | `xerj --insecure --data-dir <XERJ_DATA_DIR>` (outside repo) then `xerj autoindex .` | After add/remove/edit under library tree; incremental |
| Pack body from source tree | `pwsh -File markdown_system/update_pack_from_tree.ps1 -PackFile <pack.md> -SourceRoot <dir> -Files <rel...>` | Refresh embedded FILE blocks in a pack |
| Library V403 verifier | `pwsh -File ./markdown_system/VERIFY_LIBRARY_V403.ps1 -LibraryRoot .` | Preflight/extended library checks (requires PowerShell 7+) |
| Search index | **Manual** — append paths to `markdown_system/LIBRARY_SEARCH_INDEX.md`; no in-repo generator at 7dbf39c | When new top-level surfaces or code modules must be findable |

## Remaining GAPs (honest, selected doors ~116–120)

1. **GTM pack** — **FALTA** (intentional): no `markdown_system/*GTM*` or `implementation_packs/*GTM*` on tip; reporting adapters are read-only ads, not tag management. Do not invent.
2. **SDR filename** — **FALTA** dedicated `*SDR*` pack (intentional): use `BUSINESS_FUNCTION_OPERATING_CONTRACT_V403.md` + existing `GO_LEAD_CANDIDATE_PROMOTION.md` / `GO_OMNICHANNEL_LEAD_INGRESS.md`; gap gate if a named SDR workflow pack is REQUIRED.
3. **QR capture pack** — **identity wired** via existing `GO_QR_CORE.md` + `internal/qr/`; **capture** (`qr_capture/`) remains reference-only (no admitted capture pack). Do not invent a capture pack body.
4. **Mobile native** — `PACK_PER_CLAIM_INDEX` already marks `SPEC_ONLY` / gap gate.
5. **Production / cloud / live providers** — conditioned across all domains; `FRANCHISE_GAP_MAP.md` T2805+ live credentials still pending.
6. **CI** — no `.github/workflows` at this commit; local verification via PowerShell scripts only.

## Verification performed for this report

- Entrypoint relative-link scan (navigation files): 31 links, 0 broken.
- `rg` spot-check: `GO_SEO_CORE`, `GO_HTTP_METRICS_REFERENCE`, `GO_CONVERSATIONAL_AGENT`, `GO_QR_CORE`, `qr_capture`, `GO_LEAD_CANDIDATE_PROMOTION`, `GO_OMNICHANNEL_LEAD_INGRESS`, `db/migrations` exist on disk.
- Index lines added for previously unlisted agent entry docs and reference code paths.
- Orphan wiring (@ tip after #10): selected doors cite `GO_QR_CORE` / lead packs already on tip; no new pack bodies.
