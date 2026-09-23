# Roadmap — elite-engineering-library (tip pin)

**Tip SHA:** `0d9c258cc66a6748cd86dd8b38050b1cb5c2a0b6`  
**Scope:** biblioteca / infrastructure / local fixtures. **No** production. **No** REVESTEX product code.

Re-sign Nightly: **READY_RECONFIRMED** for `LIBRARY_INFRASTRUCTURE` / local-fixtures only (`qualification/FINAL_LIBRARY_READY_V402.json`). Not production. Not “every orphan in selected JSON”.

Nightly signed sources cited below: `qualification/FINAL_LIBRARY_READY_V402.json`, `markdown_system/LIBRARY_HEALTH_CHECK.md`, `markdown_system/PACK_PER_CLAIM_INDEX.md`, `markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md` (selected **116**; **GO-QR-CORE** still **not** in packId list), `markdown_system/FRANCHISE_COMPLETE_EXTENSION_PACK_PLAN_V403.md` (+4 → **120**), `START_FRANCHISE.md`, `START_REFERENCE_V402.md`, `LIBRARY_VS_PRODUCT_GATE_V402.md` (stub → `reconstruction_evidence/LIBRARY_VS_PRODUCT_GATE_V402.md`), `unified-experience-next/markdown_system/SUCCESSOR_AUTHORITY_ROUTER.md`, and `implementation_packs/*.md` on disk (catalog ≠ selected).

See also: [`docs/FRANCHISE_PLAYBOOK.md`](FRANCHISE_PLAYBOOK.md).

## Matrix (Nightly)

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

## Honest bounds

- Catalog pack markdown on disk ≠ selected 116 ≠ V403 extension 120.
- Lead Form ≠ GTM Tag Manager.
- `qr_capture/` ≠ selected admitted pack; identity uses `GO_QR_CORE` (CONDITIONED / not in selected JSON).
- Do not invent a named SDR pack or GTM pack.
- Stop line for product work: see playbook — no REVESTEX product until N L says **ESTAMOS LISTOS**.