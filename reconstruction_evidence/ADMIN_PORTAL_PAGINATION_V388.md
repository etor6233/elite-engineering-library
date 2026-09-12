# V388 — admin cursor traversal and combined private lists

Library maintenance / TEST02 actual owner integration. Continues V387; V386 native
vulnerability validation remains ACCESS_BLOCKED by user instruction. The platform
restriction is not bypassed or treated as security acceptance.

Admin orders, service cases and leads had the same unreachable-after25defect. One
baseline rendering test failed; the shared cursor navigator now serves all six
existing lists across customer/factory/admin. No new domain/API/SQL/dependency.
Admin keeps three independent cursors; lead:read still gates the lead request and
cursor. Removing that grant removes the lead section and its cursor from links.
Scope/page size remain server-selected. Cursor values are encoded, bounded and
checked using the V387 helper. Four new regressions also cover independent scopes,
final-page reset links and an irrelevant lead cursor without the lead grant.

276web tests PASS with the same inherited opt-in skip; typecheck/build PASS.
The actual connected suite is executed separately in four browser projects against
production Next/JWE, Go/JWKS and disposable PostgreSQL with53unchanged migrations.
Each browser traverses27customer orders,27customer cases,27factory units,28admin
orders,27admin cases and27leads:163records across lists, with overlaps between
roles. The exact per-list sequences prove no omission/duplication. Reload retains
positions, returning one list to first preserves the others, and all lists return
to their root URL after resets. Existing role/read-error cases remain in the suite.
Five direct HTTP negatives; zero browser/API writes and unchanged content digest
of six durable tables. All owned fixture processes stopped.

Six canonical packs reconstruct and the full69pack/805file reference matches the
tested source; parent68/797. Only declared page/helper/test/docs/source-lock changes;
Go HTTP binary is byte-identical to V375. V387's8589authenticated dependency hashes
are reused for the identical installation; no new install or artifact substitution.
No test or binary-identity claim closes the deferred native security admission.

FAIL770 is repaired by connected evidence. Preflight236 had correctly rejected the
stale FRANCHISE_PREFLIGHT_GAP counts (FAIL771); the canonical inventory was repaired,
and the final full Preflight is pending below. Failed236 and original source/test
logs remain preserved. V387 correction is carried by this combined final gate.

TEST02 remains blocked for broader functional equivalence/integration; TEST03/07
retain their blockers.45/48 unchanged. No new business policy, source admission,
training, provider effect, deployment or final release is claimed.

## Receipts

Raw local artifacts resolve through elite-v388-current.txt. Reused candidate root
is recorded in tools.json; V387's final805source tree and earlier browser artifacts
are preserved separately. Final independent reconstruction is under this stage.

| Receipt | SHA-256 |
|---|---|
| baseline.json | 045ea1bd7eef4bd7db05616db112a72b08142e9dfbc7c684a5c6f18fb389318c |
| web-pagination-qualified-results.json | 4c02c4df463998a0b64b8d109b7860a78f6825bf479cc26fb3c3629f73adae77 |
| canonical-pending-changes.json | 32d37daa225cec158dbd2ed9a54a66815ee288eb5b39c9dc7a87fc4805f198db |
| final-parity.json | d84a667fc60016246ab7fb463f243efc24d57a935b4299b9d4dd30462b41b7c2 |
| connected result.json | cd9f245c68b1ee6849b5593110b33d32d3d072fd9fc32695476eae78d984e856 |
| connected admin-reads.log | 53c2387f21589984820047a40d6ded80ffcdf660370f7c2aa44b99ef824d79d8 |

## Final closure239

V388 cierre239: seis listados cliente/fabrica/admin conectados al cursor real;163registros de lista por navegador,4navegadores PASS,0escrituras/snapshot6tablas idéntico.276tests/1skip previo,typecheck/build;6packs y805files reconstruidos,HTTPbinario idéntico. Preflight238164pasos/56perfiles PASS; Docker ausente.165packs/1561files/836Markdown; integral68/797.45/48,TEST02integral pendiente,TEST03 ACCESS_BLOCKED por restriccion reportada/usuario,TEST07dependiente.

```json
{
  "checkpoint_under_test": 238,
  "steps": 164,
  "step_status": "ALL_PASS",
  "compositions": 56,
  "inventory": {
    "packs": 165,
    "materialized_files": 1561,
    "markdown_files": 836
  },
  "franchise": {
    "packs": 68,
    "files": 797
  },
  "http_reference": {
    "packs": 69,
    "files": 805
  },
  "web_tests": {
    "passed": 276,
    "inherited_opt_in_skipped": 1
  },
  "connected_browser_projects": 4,
  "lists_per_browser": 6,
  "list_records_per_browser": 163,
  "writes": 0,
  "durable_snapshot_unchanged": true,
  "missing_toolchains": [
    "docker"
  ],
  "native_security_research": "ACCESS_BLOCKED_BY_REPORTED_PLATFORM_RESTRICTION_AND_USER_DEFERRAL",
  "controls_passed": 45,
  "controls_total": 48,
  "acceptance_unchanged": true,
  "preflight_log_sha256": "2bf79102c4f80044aec22a3d40c7cc26f219a9c27a5ef78510c99e9586097bb8",
  "preflight_report_sha256": "be49e665e69e3762ad95a7bb83dd5dd0088ae8752455b036f51205c6d86fe44d",
  "canonical_parity_sha256": "d84a667fc60016246ab7fb463f243efc24d57a935b4299b9d4dd30462b41b7c2"
}
```

FAIL771 REGRESSION_PROVEN by full unchanged verifier on corrected inventory. Preflight236 failed report remains in V387 staging. This final gate includes both V387 customer/factory and V388 admin changes.
