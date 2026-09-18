# Public web V403 — T2804 additive composition

Revision: `V403-PUBLIC-WEB-0.1.0-r3`. Pack `TS-PUBLIC-WEB-EXPERIENCE-V403` version `0.1.0`. Phase 2 authorized; public visual approval remains **PENDING_USER_APPROVAL**.


## Existing owner and preconditions

The complete base remains [FRANCHISE_COMPLETE_EXTENSION_PACK_PLAN_V403.md](FRANCHISE_COMPLETE_EXTENSION_PACK_PLAN_V403.md). Its UI 0.3.0 must be applied first. This is one additional bundle, not a second franchise roadmap or design system. No protected337/V402 file or immutable ZIP is replaced.
Source compatibility is enforced by `public_web_overlay/overlay-manifest.json`; package and unaffected private source hashes must match. Apply only in a NEW independent reference.

## Exact additional selection

```json
{
  "compositionVersion": "1.0",
  "secretVariables": {},
  "packs": [
    {
      "path": "implementation_packs/TYPESCRIPT_PUBLIC_WEB_EXPERIENCE_V403.md",
      "packId": "TS-PUBLIC-WEB-EXPERIENCE-V403",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    }
  ]
}
```

## Declared path overlap and application

The standard compositor creates only `public_web_overlay/**`, disjoint from existing product files. The later explicit overlay replaces exactly the public target paths listed in its manifest after all before hashes match.
`src/app/layout.tsx` is the permitted shared rendering bridge, reviewed with private screenshot/behavior regression. The exact `src/app/api/enterprise/leads/route.ts` and `src/app/api/enterprise/appointments/route.ts` are the only public-effect API exceptions, narrowly adding Origin enforcement discovered by the negative public journey. All other API files remain protected. `ApplicationChrome`, private navigation/sidebar/drawer, global CSS and active design tokens are not modified. Browser evidence must prove private behavior independently of static file preservation; separate inquiry/appointment negatives prove Origin/consent/idempotency boundaries.

Use `python public_web_overlay/apply_public_overlay.py --target <new-reference> --report <new-evidence>/apply.json`, then repeat with `--verify --report <new-evidence>/verify.json`.
Build and run the public qualification entry point described in the pack. Never regenerate a manifest from a drifted consumer merely to make hashes pass.

## Honest release dimensions

| Dimension | Rule |
|---|---|
| Reconstruction | Requires standard materializer + exact overlay receipt for this revision |
| Public behavior | Requires local public journey/negative-state tests actually executed |
| Admin preservation | Requires original-private source guards AND targeted browser regression |
| Public visual candidates | PENDING_USER_APPROVAL; no inheritance from admin or Brave references |
| Screen reader / physical phone / cloud | NOT_RUN unless independent exact evidence exists |
| Production | Not authorized |

Dependencies remain fixed by the inherited lock. New components, pictograms, tooling and fixture integration are AUTHORED; method references are not source provenance. Own generated media are declared in PROVENANCE.json and must not be read as real-product photography.
