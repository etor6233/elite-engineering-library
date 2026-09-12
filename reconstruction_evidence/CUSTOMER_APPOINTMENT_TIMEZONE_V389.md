# V389 — customer appointment time-zone consistency

Library maintenance / TEST02 owner integration. Checkpoint239 resume passed;
checkpoint240 captured the actual presentation defect and baseline tests.
V386 native vulnerability research remains deferred by the user's instruction.

## Corrected behavior

The customer account rendered appointment time using the server default locale
and zone; interactive appointment management used its own host defaults. Both
ignored the already configured market zone used by public booking. These could
show different dates/hours for the same appointment and change during hydration.

Both pages now load the existing trusted business locale/market zone and reuse
the existing public appointment formatter. Each time element preserves the
original API ISO datetime and displays the explicit zone. The management page
passes its server-formatted label to the client component. This prevents a device
locale/zone from changing the appointment display. Invalid dates fail through
the existing formatter instead of showing Invalid Date. No new dependency, rule,
API, stored instant, cancellation command, permission or eligibility change.

## Evidence

Baseline: five failing rendering cases, ten passing. Corrected final web suite:
282passed and one inherited opt-in skip; typecheck/build passed. Six added cases
exercise configured Spanish/Argentina, English/New York winter and summer,
Tokyo, invalid timestamps and permission checks. The skip is not counted as a
connected pass; the actual connected harness executes separately.

Two fresh PostgreSQL fixtures with all53unchanged migrations and actual
Next/JWE → Go/RS256/JWKS → PostgreSQL reads run four browsers each, eight runs in
total. Configuration A uses es-AR/Buenos Aires; configuration B en-US/New York.
Node server TZ=UTC, browser locale de-DE and zone Pacific/Honolulu deliberately
differ. Two winter/summer appointment instants retain their original API datetime
and show the independently expected local date/hour/zone in account, management
and after reload. Cancellation buttons remain present for the two future active
appointments; no cancellation effect is claimed by this read-only test.

Each run also preserves the prior six-list traversal, role/grant checks, fixed
read-error recovery and five direct HTTP negatives. There are zero browser/API
writes and identical content snapshots of all seven observed durable tables,
including appointments. All owned PostgreSQL/Next processes stopped. The mobile
view has no horizontal overflow; screenshot inspection is recorded below.

Six canonical packs reconstruct to the tested source. Final69pack/805file
reference and68pack/797file parent retain counts; no extra runtime source is
introduced. The Go HTTP executable is byte-identical to V375. Dependency locks
and unchanged installed artifacts retain the previous admission conditions;
this work makes no new dependency or native-security claim.

## Failures and scope limits

FAIL772 is the implicit host-zone defect. FAIL773 preserves the mistaken24hour
es-AR display oracle: pinned Intl emits equivalent10:30p.m. FAIL774 preserves
four failed first browsers whose test selected only the literal Z string; traces
demonstrate the API returned equivalent -03:00. The corrected oracle compares
exact epoch instants plus independently expected labels and retains raw API
datetime equality across views/reload. No product or database constraint was
weakened to satisfy these tests.

FAIL775 preserved WebKit prefetch errors during forced navigation/reload in the
self-signed-loopback fixture. A temporary error classifier was rejected; the
final test instead waits for prior requests to settle before observing the date
segment and before reload. All eight final runs require zero page errors and
zero React hydration errors in that segment. This qualifies fixture sequencing,
not an upstream browser/TLS repair or arbitrary rapid navigation.

This closes customer appointment display consistency only. Other private portal
dates, localized interface catalogs, broader functional equivalences and actual
target/business/provider conditions remain open. TEST02,TEST03,TEST07 and all48
acceptance definitions/statuses remain unchanged:45/48. No release approval.

## Local receipts

The stage resolves through Temp/elite-v389-current.txt; tools.json identifies the
reused installed candidate. V388's final805source tree and failed V389 receipts
remain preserved. This stage has its own independent final reconstruction.

| Receipt | SHA-256 |
|---|---|
| baseline.json | 3be7f05ca91fbbeed216e565809c2d0814c36efedcafaac6e694d9b9ef064b10 |
| web-qualified-results.json | 8788bd2d787e63d1bf8a1e6e126ecc5acb234d1d28fbf47ed8352b4fe0844846 |
| time-trace-observed.json | 8e31f7c996fa89357f000e30f6f9177376f22fe324da1b0fe239f71a9b8dfcff |
| artifact-parity.json | 7563f34b2a058f7ce4c0f799cb29ad7b97d517326579e335c400497c8723cb43 |
| canonical-pending-changes.json | 6f9a4d6af4ac0f0020ce47bf34f717614028a37a053838e4e8a63d5b469816a1 |
| final-parity.json | 718d9d1bf6a1bd16a47cb4bd8d1d506573f1152c34fce0382b0a266556d217e7 |
| ar/result.json | c11edaaffc3f77eab48ad5086355e99a4145d2de9b6d9a2fcdaf110834519294 |
| ar/admin-reads.log | f48b48832b3457d86a988455d56a6b39beceaa83e34f47409f8ed77fed1930a7 |
| ny/result.json | 46799ada6e1078f8606a2ff5fc2f2e277815e605f171124805ef7c8642bd8556 |
| ny/admin-reads.log | 7326c260f23c142d530fb9825019da0f348e1d51bdb7ab01766fe79dcb773168 |

Full Preflight241 pending.

## Final checkpoint242

V389 cierre242: turnos cliente muestran fecha/hora/zona del negocio coherentes en cuenta/gestion/recarga, sin depender del dispositivo.282tests/1skip previo,8runs navegador(2config x4),0errores en segmento de fechas,0escrituras/snapshot7tablas idéntico.6packs/805files reconstruidos,8589artefactos yHTTPbinario idénticos. Preflight241164pasos/56composiciones PASS; Docker ausente.165packs/1561files/837Markdown;68/797.45/48; TEST02integral pendiente, V386native diferido yTEST07dependiente.

Mobile Argentina appointment screenshot inspected: both dates/zones legible, action buttons visible, no horizontal overflow. Final browser source removes the temporary prefetch exception classifier entirely. Quiescing the previous page and waiting for prefetched requests before forced reload passes all eight runs with zero page errors and zero hydration errors in the measured date segment. This qualifies the isolated test sequence, not arbitrary rapid navigation or a WebKit/TLS upstream repair. All failed candidates remain preserved.

```json
{
  "checkpoint_under_test": 241,
  "steps": 164,
  "step_status": "ALL_PASS",
  "compositions": 56,
  "inventory": {
    "packs": 165,
    "files": 1561,
    "markdown": 837
  },
  "franchise_files": 797,
  "http_reference_files": 805,
  "web_tests": 282,
  "inherited_opt_in_skip": 1,
  "browser_runs": 8,
  "business_configurations": 2,
  "appointment_instants": 2,
  "rendered_views": [
    "customer",
    "customer/appointments",
    "appointments_reload"
  ],
  "date_segment_page_errors": 0,
  "hydration_errors": 0,
  "writes": 0,
  "snapshot_tables": 7,
  "missing_toolchains": [
    "docker"
  ],
  "controls_passed": 45,
  "controls_total": 48,
  "acceptance_unchanged": true,
  "preflight_report_sha256": "9a6bee8c6b134878bca3a96ac8834e9ac7678371efd7109487bb187fe0a1fbd6",
  "final_parity_sha256": "718d9d1bf6a1bd16a47cb4bd8d1d506573f1152c34fce0382b0a266556d217e7",
  "dependency_file_parity_sha256": "7563f34b2a058f7ce4c0f799cb29ad7b97d517326579e335c400497c8723cb43",
  "behavior_parity_sha256": "ff9312c41189ff03aafb64a8d9799f1364e0e01d2aa517becbc7c77e7b4ddaf1"
}
```
