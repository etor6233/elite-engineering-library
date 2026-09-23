# Franchise playbook — assemble from this library

## 1. Objective / tip pin

Assemble a **franchise consumer** from **elite-engineering-library** at tip:

`0d9c258cc66a6748cd86dd8b38050b1cb5c2a0b6`

This library is infrastructure + local fixtures evidence — **not** a live franchise, **not** production authorization, **not** REVESTEX product code.

Roadmap matrix: [`docs/ROADMAP.md`](ROADMAP.md).

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

## 4. Assemble steps

1. Pin this tip SHA; install agent bridge if the consumer root ≠ library root (`INSTALL_AGENT_BRIDGE.ps1`).
2. Walk door order §2; read HEALTH gaps before composing.
3. Select composition: base 116 **or** V403 extension 120; cite exact `packId`/`version` from the plan JSON — do not add undeclared packs.
4. Materialize with the library compositor / franchise protocol named in `START_FRANCHISE.md` and the extension plan.
5. Verify local receipts/fixtures only; record CONDITIONED live gates (creds, providers, IdP) as open — not PASS.
6. Keep consumer blueprint / state / secrets in the **consumer**; reusable owners stay in the library.

## 5. Galaxy fleet pin (orchestration only)

Fleet orchestration surface pin:

`etor6233/grok-bot-galaxy@f8546c3`

If a local/galaxy survey is not available in this checkout: cite that pin and treat **fleet evidence = repo doors at that pin**. **Do not invent bots**, agents, or fleet members.

## 6. Stop line

**No REVESTEX product** work from this library until N L says **ESTAMOS LISTOS**.

READY_RECONFIRMED applies to library infrastructure / local fixtures only — never to production admission.