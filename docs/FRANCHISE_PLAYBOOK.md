# Franchise playbook — assemble from this library

## 1. Objective / tip pin

Assemble a **franchise consumer** from **elite-engineering-library** at tip:

`477358d0d7090f2c39467abe5634d81f38afae7f`

This library is infrastructure + local fixtures evidence — **not** a live franchise, **not** production authorization, **not** REVESTEX product code.

**Scope split (course correction):**

| Doc | Scope |
|---|---|
| [`docs/FRANCHISE_ARCHITECTURE.md`](FRANCHISE_ARCHITECTURE.md) | **Full franchise / company domains** — identity, catalog, CRM, sales, channels, logistics, brand, admin, leads, QR, Galaxy; HECHO/PARCIAL/NO with path cites |
| [`docs/ROADMAP.md`](ROADMAP.md) | Same domains in matrix form + **`LIBRARY_INFRASTRUCTURE`** Nightly wiring (local fixtures only) |
| This playbook §3 CAN/CANNOT | **Library infrastructure / local-fixtures** — **not** the franchise product architecture surface |

Roadmap matrix: [`docs/ROADMAP.md`](ROADMAP.md). Domain architecture: [`docs/FRANCHISE_ARCHITECTURE.md`](FRANCHISE_ARCHITECTURE.md).

## 2. Door order

Follow doors in order (same table as README / AGENTS):

1. `README.md` — library vs product
2. `AGENTS.md` — compact agent router
3. `markdown_system/USE_LIBRARY_V403.md` — use / search / resume
4. `markdown_system/LIBRARY_HEALTH_CHECK.md` — health, wiring, honest GAPs
5. `START_FRANCHISE.md` — new franchise / consumer from chosen copy
6. Recon gate: root stub `LIBRARY_VS_PRODUCT_GATE_V402.md` → `reconstruction_evidence/LIBRARY_VS_PRODUCT_GATE_V402.md`
7. Receipt: `qualification/FINAL_LIBRARY_READY_V402.json` (evidence only; does not open production)

After 1–7, optional operators: `START_REFERENCE_V402.md` (ZIP / local reference / portable firma+verify).

## 3. CAN / CANNOT

> **Not franchise product scope.** This list answers “what can the **library** compose and verify locally?” For end-to-end business domains (catalog, CRM, payments, multi-site, channels, fiscal, …), use [`docs/FRANCHISE_ARCHITECTURE.md`](FRANCHISE_ARCHITECTURE.md) — not this section alone.

### CAN

- Compose **selected 116** from `markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md`, or **116+4 → 120** via `markdown_system/FRANCHISE_COMPLETE_EXTENSION_PACK_PLAN_V403.md`, using packs on disk under `implementation_packs/`.
- Bridge → `START_FRANCHISE.md` → verify **local** receipts/fixtures (`VERIFY_LIBRARY.ps1`, `VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Preflight`).
- Google Lead Form via `GO-OMNICHANNEL-LEAD-INGRESS` (`implementation_packs/GO_OMNICHANNEL_LEAD_INGRESS.md`) + Meta/TikTok + `GO-LEAD-CANDIDATE-PROMOTION`.
- QR **identity** via `GO-QR-CORE` + `internal/qr/` (**CONDITIONED**; pack exists on disk but **GO-QR-CORE is not** in the selected 116 packId JSON).
- QR capture **reference** in `qr_capture/` (not an admitted selected pack).
- ARCA infra/fixtures **PROVEN_LOCAL** (user creds required for live).
- Portable firma/verify + ZIPs guided by `START_REFERENCE_V402.md`.
- Treat V403 SUCCESSOR as **IN_PROGRESS** (`unified-experience-next/markdown_system/SUCCESSOR_AUTHORITY_ROUTER.md`) — **do not inherit** V402 READY onto NEXT.

### CANNOT

- Authorize production / live franchise.
- Invent a named **SDR** pack or **GTM Tag Manager** pack.
- Treat Lead Form as GTM.
- Treat `qr_capture/` as a selected admitted pack.
- Assume Desktop ***V3*** recoverable from Git (`UNKNOWN_NEVER_PUSHED`).
- Confuse catalog pack count on disk with selected **116** / extension **120**.
- Treat V403 NEXT as closed library-ready.
- Treat Galaxy how-tos / Duplicate / Marketplace as franchise e2e or installable fleet packs.

## 4. Assemble steps

1. Pin this tip SHA; install agent bridge if the consumer root ≠ library root (`INSTALL_AGENT_BRIDGE.ps1`).
2. Walk door order §2; read HEALTH gaps before composing.
3. Select composition: base 116 **or** V403 extension 120; cite exact `packId`/`version` from the plan JSON — do not add undeclared packs.
4. Materialize with the library compositor / franchise protocol named in `START_FRANCHISE.md` and the extension plan.
5. Verify local receipts/fixtures only; record CONDITIONED live gates (creds, providers, IdP) as open — not PASS.
6. Keep consumer blueprint / state / secrets in the **consumer**; reusable owners stay in the library.

## 5. Galaxy — fleet orchestration (survey truth)

Fleet orchestration surface pin (full tip):

`etor6233/grok-bot-galaxy@f8546c3990a165be8d17a178880a364aeeb3ac80`

(short: `f8546c3`)

**What it is:** a historical field guide / KB for operating bot fleets. Reading `llms.txt` (or any Galaxy how-to) **≠** installing bots, connectors, or routines.

**Exists at tip (cite doors, do not invent members):**

- `ROUTES` / `INDEX` (navigation)
- how-tos / demos: `why-a-fleet-of-bots`, `flylo-engineering-fleet`, `operate-a-bot-team`, `share-duplicate-marketplace`, `agent-to-agent`
- `validate.yml` = **KB QA only** (not franchise e2e, not provisioning)

**Gaps (honest):**

- `franchise` string hits = **0** at tip
- no installable packs / provisioning IaC
- no multi-tenant franchise fleet product
- Desktop **Elite Franchise V402** is **unlinked** from galaxy (do not assume a wire)
- Duplicate / Marketplace share patterns ≠ franchise e2e

If a local/galaxy survey is not available in this checkout: cite that pin and treat **fleet evidence = repo doors at that pin**. **Do not invent bots**, agents, or fleet members.

Optional ops artifact (Desktop only; **not** required in this library repo): `REVESTEX/galaxy-franchise-fleet-survey.md`.

## 6. Stop line

**No REVESTEX product** work from this library until N L says **ESTAMOS LISTOS**.

READY_WITH_CAVEATS applies to library infrastructure / local fixtures only — never to production admission.