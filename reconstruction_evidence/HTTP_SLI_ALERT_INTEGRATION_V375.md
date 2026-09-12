# HTTP SLI and real alert integration — V375

2026-09-10. Library maintenance; bounded trusted-local reference, synthetic data only.
45/48 controls remain closed. TEST02/03/07 and the broader FAIL385 integration
remain open. This closes the HTTP request → metrics → canonical alert → repair
reference boundary; it does not assert all-service SLOs or production acceptance.

The actual library order handler, service, repository and complete migrations run
against a fresh PostgreSQL18.6 cluster. Official otelhttp0.70.0 instruments actual
requests; SDK1.46.0 and exporter0.67.0 expose a closed method/status duration
histogram. Prometheus3.14.0 scrapes those series. The unchanged SECUREOPS1.1.4
rule retains its five-minute window, one-percent threshold, ten-minute hold and
thirty-second evaluation. Renaming only the owned synthetic order table causes
real HTTP500 responses. Restoring it and replaying the original order clears
the alert without another order, idempotency record or outbox event.

No samples are injected into Prometheus, clock changed or alert hold shortened.
The fixed principal belongs exclusively to the reference main; the production
OIDC main is unchanged. API/metrics use loopback HTTP on the trusted account;
Prometheus queries use ephemeral mutual TLS1.3. This does not prove production
identity, telemetry exposure, all PII, business SLI selection, alert routing,
human response, power-loss recovery, all-service coverage or security clearance.

## Official source and dependency admission

Official APIs and documentation inspected: versioned otelhttp package documentation
(https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.70.0),
HTTP metric semantics (https://opentelemetry.io/docs/specs/semconv/http/http-metrics/)
and Go instrumentation libraries (https://opentelemetry.io/docs/languages/go/libraries/).
Exact source ZIP/cache parity and immutable origins are embedded below and in
source-lock.json. No sample implementation is copied. The eight wrapper/lock/test
files are AUTHORED; runtime middleware/exporter sources and test assertions are
unchanged official code. Tests select the consumer SDK1.46.0; the exporter test
copy removes six monorepo-relative replaces, preserving the original manifest.
238 middleware plus 142 exporter tests pass, zero failures/skips. Two authored
HTTP privacy/isolation tests, vet, build and module verification pass.

The reference's complete public graph has 62 unique version queries and zero
OSV findings at the recorded scan. Thirty compiled upstream modules have exact
license/notice file hashes; two local modules remain AUTHORED. Four NOTICE files
are attribution documents, not unrecognized licenses. Preserve all notices on
redistribution. The reused V372 Prometheus/Collector runtime dispositions remain
separate: three graph advisories have package-absent compiled-closure evidence;
this is not a zero-finding claim for their entire source graphs or TEST03.

## Rebuild and prior failures

FAIL719 repeated guessed read paths; the manifest is the authority for paths.
FAIL720 copied readonly Go manifests; only staged attributes changed.
FAIL721 recovered two pinned official proxy metadata records, no source update.
FAIL722 initdb environment omitted COMSPEC; the exact proven environment restores
real cluster startup. FAIL723 preserves exporter standalone test-layout failure.
FAIL724 completes the unchanged parent command for source parity.
FAIL725 preserves a 758-second failed run: the alert fired, but three increments
inside one sliding window broke the legacy-ratio oracle. Recovery did not run.
The successor uses one failure per four minutes, retaining every rule criterion.
FAIL726 preserves the minimal-parent binary mismatch: full composed packages,
not eighteen backend files alone, are the actual compiler inputs. A fresh
67-pack/755-file reconstruction matches all 441 original parent files and builds
the exact byte-identical tested executable. No rerun is inferred across changed
binary bytes. The final lock adds the parent plan digest without altering Go code.

## Consumption and remaining conditions

GO-HTTP-METRICS-REFERENCE0.1.0 is REBUILD_VERIFIED / CONDITIONED, selected only by
HTTP_METRICS_REFERENCE_PACK_PLAN (68 packs / 763 files). Ordinary franchise remains
67 / 755. Reconstruct the full profile, run the embedded Go suites, vet and build,
bind exact executable/rule/runtime-lock hashes, then execute run_reference.py.
Use an absent destination and a private owned work root. Runtime binaries stay
external, licensed and admitted; this pack never installs or deploys them.
The final structural/Preflight receipt is appended after execution, not assumed.

Raw receipts remain in the retained local stage with SHA identities below.
Displayed local paths are explicitly normalized to $V375/$LOCAL_USER; original
bytes were not rewritten. This report contains synthetic identifiers only.

## Actual runtime result

```json
{
  "state": "PASS",
  "rules_sha256": "cc718b7bbe35de3f230eeb47c854c270580dbb488101ee5f30c0b758e1d967b0",
  "host_sha256": "4643654625c71aaf5b45367bd3a198381257588d073ae0349fc90c1bc1b5ab85",
  "requests": {
    "201": 1,
    "200": 501,
    "401": 1,
    "400": 1,
    "500": 4
  },
  "synthetic_only": true,
  "firing_alert": [
    {
      "labels": {
        "alertname": "PlatformHighErrorRate",
        "severity": "page"
      },
      "annotations": {
        "runbook_url": "https://replace.invalid/runbooks/api-availability",
        "summary": "User-visible server error rate exceeds one percent"
      },
      "state": "firing",
      "activeAt": "2026-09-10T16:53:14.950346156Z",
      "value": "1e+00"
    }
  ],
  "observed_ratios": {
    "corrected": {
      "status": "success",
      "data": {
        "resultType": "vector",
        "result": [
          {
            "metric": {},
            "value": [
              1789059795.524,
              "1"
            ]
          }
        ]
      }
    },
    "legacy": {
      "status": "success",
      "data": {
        "resultType": "vector",
        "result": [
          {
            "metric": {},
            "value": [
              1789059795.542,
              "0.003344403308283753"
            ]
          }
        ]
      }
    }
  },
  "actual_sparse_failures": 4,
  "pending_to_firing_seconds": 599.904,
  "fault_to_firing_seconds": 856.427,
  "same_order_after_recovery": true,
  "order_count": 1,
  "outbox_count": 1,
  "idempotency_count": 1,
  "children_exited": true,
  "elapsed_seconds": 911.211,
  "event_log_sha256": "18077245c7ee248577d1d79a0bc9da8216d629c1dca7f78d7e74dd3842ef5bf1"
}
```

## Actual event sequence

```json
[
  {
    "seconds": 24.796,
    "kind": "fault_injected",
    "method": "rename owned synthetic order table",
    "interval_seconds": 240,
    "rule_hold_seconds": 600
  },
  {
    "seconds": 24.824,
    "kind": "actual_http_500",
    "count": 1
  },
  {
    "seconds": 24.842,
    "kind": "alert_state",
    "state": "inactive"
  },
  {
    "seconds": 265.076,
    "kind": "actual_http_500",
    "count": 2
  },
  {
    "seconds": 281.319,
    "kind": "alert_state",
    "state": "pending"
  },
  {
    "seconds": 505.25,
    "kind": "actual_http_500",
    "count": 3
  },
  {
    "seconds": 745.258,
    "kind": "actual_http_500",
    "count": 4
  },
  {
    "seconds": 881.222,
    "kind": "alert_state",
    "state": "firing"
  },
  {
    "seconds": 911.054,
    "kind": "recovered",
    "order_count": 1,
    "outbox_count": 1,
    "idempotency_count": 1
  }
]
```

## source-identity.json

```json
{
  "module": "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp",
  "version": "v0.70.0",
  "origin": {
    "VCS": "git",
    "URL": "https://github.com/open-telemetry/opentelemetry-go-contrib",
    "Subdir": "instrumentation/net/http/otelhttp",
    "Hash": "c8a87a60ba1b3374fd16df11fc3eeae6c41abbc9",
    "Ref": "refs/tags/instrumentation/net/http/otelhttp/v0.70.0"
  },
  "zip_sha256": "5afae012c182ed95c9068c4b4de772677fbfecb9fbe41605046f907abba8bc92",
  "go_sum": "h1:LMuyCAyfalSjDyjdC65nK6N0zoTT63+E/u95X0JovZI=",
  "files_verified_against_zip": 36,
  "license_sha256": "1ae07514be1d7bb33f0698f8d91fb51b8b9fe1463157ec1c72081a49b9bc6f40",
  "license_expression": "Apache-2.0 AND BSD-3-Clause",
  "runtime_source_changed": false,
  "cached_source_reuse": true
}
```

## exporter-source-identity.json

```json
{
  "module": "go.opentelemetry.io/otel/exporters/prometheus",
  "version": "v0.67.0",
  "origin": {
    "VCS": "git",
    "URL": "https://github.com/open-telemetry/opentelemetry-go",
    "Subdir": "exporters/prometheus",
    "Hash": "93a693edeed0e07ce5ebd1dfe67af42d1e2055d8",
    "Ref": "refs/tags/exporters/prometheus/v0.67.0"
  },
  "zip_sha256": "ce725796160333dcee033dc9e6a2a81ea8d6fb24846f9ddbdc714c5bc86d6d43",
  "go_sum": "h1:7IefDa35e6V3NoiqIeLDMDxMFyZDk5qcoC0Ax4cC16E=",
  "files_verified_against_zip": 63,
  "license_sha256": "1ae07514be1d7bb33f0698f8d91fb51b8b9fe1463157ec1c72081a49b9bc6f40",
  "license_expression": "Apache-2.0",
  "runtime_source_changed": false
}
```

## official-tests-receipt.json

```json
{
  "counts": {
    "pass": 238,
    "fail": 0,
    "skip": 0
  },
  "log_sha256": "ce8c8e3780fde47552a503ed8f75ceb47fc47a02c85c0cb2c15c1b36fa8cb949",
  "go_mod_sha256": "ae40d21c6b98aa51f6395bf9be63a7c8af628a601127066e054c773c2720bc44",
  "go_sum_sha256": "8898d3ecdd0ab3263d94535e36064dea461f23c45257697f788dbac67758d95d"
}
```

## exporter-tests-receipt.json

```json
{
  "counts": {
    "pass": 142,
    "fail": 0,
    "skip": 0
  },
  "log_sha256": "0f507a2a77e324f8b4aa56b0f58bee77c93e32d8afc3577c6a7c21e7b621c27f",
  "go_mod_sha256": "be1bf0405faf33b77ab20f5a8bbf64c31d1a074b1cb258e87f0718b2093762c0",
  "go_sum_sha256": "e6611d60d4a7661cf7b633e8e6d265624de4c61f80e338592a4711533c3d64f7",
  "removed_monorepo_replaces": [
    "replace go.opentelemetry.io/otel => ../..",
    "replace go.opentelemetry.io/otel/sdk => ../../sdk",
    "replace go.opentelemetry.io/otel/sdk/metric => ../../sdk/metric",
    "replace go.opentelemetry.io/otel/trace => ../../trace",
    "replace go.opentelemetry.io/otel/metric => ../../metric",
    "replace go.opentelemetry.io/otel/metric/x => ../../metric/x"
  ],
  "runtime_and_test_source_modified": false
}
```

## reference-public-security-summary.json

```json
{
  "unique_public_version_queries": 62,
  "findings": [],
  "authored_replacement_excluded": true,
  "compiled_modules": 32,
  "compiled_upstream_modules": 30,
  "unknown_license_texts": [
    [
      "github.com/coreos/go-oidc/v3",
      "NOTICE"
    ],
    [
      "github.com/prometheus/client_golang",
      "NOTICE"
    ],
    [
      "github.com/prometheus/client_model",
      "NOTICE"
    ],
    [
      "github.com/prometheus/common",
      "NOTICE"
    ]
  ]
}
```

## http-full-rebuild-receipt.json

```json
{
  "candidate_files_rebuilt_identically": 8,
  "parent_consumer_files_identical": 441,
  "parent_profile_files": 755,
  "parent_profile_packs": 67,
  "binary_builds_byte_identical": true,
  "binary_sha256": "4643654625c71aaf5b45367bd3a198381257588d073ae0349fc90c1bc1b5ab85",
  "candidate_pack_sha256": "146a3d6dc6d3e7fc003a208805d247b9dd2d6a48a836cf8b0ed0dea2fcbe222a",
  "reference_tests": 2
}
```

## Preserved original receipt identities

```json
[
  {
    "local_stage_ref": "source-identity.json",
    "bytes": 778,
    "sha256": "fa86c0cbeda90591e78ed14708fac14455b698a25cff0248d052ba0de1cb78a1"
  },
  {
    "local_stage_ref": "exporter-source-identity.json",
    "bytes": 680,
    "sha256": "80b34e1b507085bfc1c089285eeea09611b23b2ad29bb14c12c343d9c631778b"
  },
  {
    "local_stage_ref": "official-tests.log",
    "bytes": 214142,
    "sha256": "ce8c8e3780fde47552a503ed8f75ceb47fc47a02c85c0cb2c15c1b36fa8cb949"
  },
  {
    "local_stage_ref": "exporter-official-tests-fixed.log",
    "bytes": 128557,
    "sha256": "0f507a2a77e324f8b4aa56b0f58bee77c93e32d8afc3577c6a7c21e7b621c27f"
  },
  {
    "local_stage_ref": "reference-test-isolation.log",
    "bytes": 269,
    "sha256": "2288885b3f4714447f950bbef825ca9685e75a39402bbc89f8479539b3e6aa39"
  },
  {
    "local_stage_ref": "reference-vet.log",
    "bytes": 0,
    "sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
  },
  {
    "local_stage_ref": "reference-build.log",
    "bytes": 0,
    "sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
  },
  {
    "local_stage_ref": "module-verify.log",
    "bytes": 21,
    "sha256": "b4537ed75f533f993f371954de47e42a793b8e5b0587577de7e27fb3e50696bd"
  },
  {
    "local_stage_ref": "compiled-imports.log",
    "bytes": 1153628,
    "sha256": "cf4de44d926785023da05bccb4348b9911850ac62a211fcf55d8e12f19d9f56c"
  },
  {
    "local_stage_ref": "compiled-license-inventory.json",
    "bytes": 14764,
    "sha256": "faae82739dfc3a7e395b7df6b7f34eb61c052097d21d07bbb603e498bd703998"
  },
  {
    "local_stage_ref": "reference-public-osv-request.json",
    "bytes": 6015,
    "sha256": "25f5e1fcd00e615359c8b978e4b2edc62d9b051385a1ae4fe86e62b4e431a0be"
  },
  {
    "local_stage_ref": "reference-public-osv-response.json",
    "bytes": 199,
    "sha256": "0398800ce15d818ab0ce4dffcef48a6b9054b03fee991473ef750bb210f742c2"
  },
  {
    "local_stage_ref": "reference-public-security-summary.json",
    "bytes": 472,
    "sha256": "54068c1258abf4e9b7b0be9e2d2fa6c3e18061833b57d36ef5881724b2f1d632"
  },
  {
    "local_stage_ref": "http-full-rebuild-receipt.json",
    "bytes": 396,
    "sha256": "6c7730ab00f8df895357b4b5b625ecd799d218312f56f45949476e94b134bd94"
  },
  {
    "local_stage_ref": "full-rebuild-tests.log",
    "bytes": 269,
    "sha256": "f4d01c91ad4600cfbf7799f2607e9d01a9c060b473128802f242740223c78263"
  },
  {
    "local_stage_ref": "full-rebuild-build.log",
    "bytes": 0,
    "sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
  },
  {
    "local_stage_ref": "http-runtime.log",
    "bytes": 1265,
    "sha256": "0ef122fb255e39f5074ce0eb9bf359d8016f7155020980fd18d1dfbdd3755a14"
  },
  {
    "local_stage_ref": "http-runtime2.log",
    "bytes": 2336,
    "sha256": "cbbf8473bb0f831c7f06c3cde697403cb31d23773d808c5cb260db8ae5e76be9"
  },
  {
    "local_stage_ref": "http-runtime3.log",
    "bytes": 2092,
    "sha256": "5b5f88dcfbe68744cccceebfc68a47726eb93621d2970ac501ef4ff745084415"
  },
  {
    "local_stage_ref": "run-http-alert-reference.py",
    "bytes": 12007,
    "sha256": "e23057fc1a71ced1dddfbefc51886ce66508d483effe64e87a933527956f31ab"
  },
  {
    "local_stage_ref": "runtime-command3.json",
    "bytes": 1086,
    "sha256": "e6fa678813d8829d78d075c5c340f9fed8da8af4af634e2e160f7e9cef69f0ad"
  },
  {
    "local_stage_ref": "runtime-runs/http-df4b1109bfa24ed78b89d8817624edfc/result.json",
    "bytes": 1679,
    "sha256": "ae360c505f8620340896addc936144b4a86420495abf442d9588d4e0ba7e785c"
  },
  {
    "local_stage_ref": "runtime-runs/http-df4b1109bfa24ed78b89d8817624edfc/progress.json",
    "bytes": 867,
    "sha256": "18077245c7ee248577d1d79a0bc9da8216d629c1dca7f78d7e74dd3842ef5bf1"
  },
  {
    "local_stage_ref": "runtime-runs/http-df4b1109bfa24ed78b89d8817624edfc/healthy.prom",
    "bytes": 5554,
    "sha256": "06007803e3ad2e30f3635e66d40d0aaeab0ab6a845c48704a5ca0bfff065cea5"
  },
  {
    "local_stage_ref": "runtime-runs/http-df4b1109bfa24ed78b89d8817624edfc/recovered.prom",
    "bytes": 6967,
    "sha256": "b294da804c8b67795a5e3cd2ee4e7a569df095adcd5d4904b3904e1cde761dac"
  },
  {
    "local_stage_ref": "runtime-runs/http-df4b1109bfa24ed78b89d8817624edfc/postgres-stop.log",
    "bytes": 58,
    "sha256": "3589d2bd01f8295cc7c9ec3f5e07625fcd347ccb6f33f2d80edafe5f5217c1fb"
  }
]
```

## G0–G8 scoped admission

G0 scope: required local HTTP signal/alert reference; G1 authority: official versioned APIs and semantics; G2 exact ZIP/commit identities; G3 source license/notice evidence; G4 380 unchanged official tests plus two authored integration/isolation tests; G5 exact62-version SCA and closed telemetry labels; G6 actual bounded fault/hold/repair and owned-process teardown; G7 complete current consumer/migration compatibility; G8 canonical reconstruction and byte-identical executable. All PASS applies only to this fixed reference. The pack stays CONDITIONED and the resolver must not return implementation_ready=true for a new target.

```json
{
  "schema_version": 1,
  "capability_id": "HTTP-SLI-ALERT-REFERENCE",
  "requirement_ref": "markdown_system/FRANCHISE_GAP_MAP.md",
  "observed_at": "2026-09-10T17:04:45.229912Z",
  "max_age_days": 30,
  "authority_refs": [
    "SECURITY_SRE_CLOUD_INFRASTRUCTURE.md",
    "reconstruction_evidence/HTTP_SLI_ALERT_INTEGRATION_V375.md"
  ],
  "source_searches": [
    {
      "search_id": "official-http-instrumentation",
      "query": "Versioned otelhttp and exporter APIs, official HTTP metric semantics; exact existing source ZIPs and unchanged tests inspected",
      "official_domains": [
        "github.com",
        "opentelemetry.io",
        "go.dev"
      ],
      "executed_at": "2026-09-10T17:04:45.229912Z",
      "evidence_refs": [
        "reconstruction_evidence/HTTP_SLI_ALERT_INTEGRATION_V375.md"
      ]
    }
  ],
  "candidates": [
    {
      "candidate_id": "authored-http-reference",
      "origin_kind": "LOCAL_AUTHORED",
      "provenance": "AUTHORED",
      "source_url": null,
      "source_ref": "implementation_packs/GO_HTTP_METRICS_REFERENCE.md",
      "claim": "Reconstruct the exact trusted synthetic HTTP metrics and fault/alert/recovery reference using pinned unchanged official dependencies",
      "non_claims": [
        "Production identity or telemetry deployment",
        "Whole TEST02/03/07 closure",
        "Official-vendor authorship of reference glue"
      ],
      "status": "CONDITIONED",
      "immutable_revision": "GO-HTTP-METRICS-REFERENCE@0.1.0",
      "artifact_sha256": "fa955025970d4a795c9f2eb547ea78040840bb5ecbc5576434061b95f8f3d7fb",
      "artifact_ref": "implementation_packs/GO_HTTP_METRICS_REFERENCE.md",
      "license_expression": "LicenseRef-Workspace-Owner",
      "license_evidence_ref": "implementation_packs/GO_HTTP_METRICS_REFERENCE.md",
      "gates": {
        "G0": "PASS",
        "G1": "PASS",
        "G2": "PASS",
        "G3": "PASS",
        "G4": "PASS",
        "G5": "PASS",
        "G6": "PASS",
        "G7": "PASS",
        "G8": "PASS"
      },
      "governing_authorities": [
        "SECURITY_SRE_CLOUD_INFRASTRUCTURE.md",
        "reconstruction_evidence/HTTP_SLI_ALERT_INTEGRATION_V375.md"
      ],
      "no_source_search_ids": [
        "official-http-instrumentation"
      ]
    }
  ],
  "decision": {
    "state": "USE_CONDITIONED_PACK",
    "selected_candidate_id": "authored-http-reference",
    "reason": "All nine gates evidenced for the exact trusted local synthetic reference. New target identity, runtime, data, security and SLI policy are not inherited; no production implementation_ready claim.",
    "canonical_updates": [
      "implementation_packs/GO_HTTP_METRICS_REFERENCE.md",
      "markdown_system/HTTP_METRICS_REFERENCE_PACK_PLAN.md",
      "reconstruction_evidence/HTTP_SLI_ALERT_INTEGRATION_V375.md",
      "markdown_system/CAPABILITY_CATALOG.md"
    ],
    "blockers": [
      "New consumer must satisfy all declared exact runtime/source/security/identity/SLI conditions; never deploy the synthetic reference main"
    ],
    "exhaustion": null
  }
}
```

## Final canonical consumer verification

```json
{
  "status": "PASS",
  "packs": 68,
  "files": 763,
  "authored_reference_files": 8,
  "parent_parity_files": 441,
  "exact_canonical_rule": true,
  "tests": 2,
  "vet": "PASS",
  "build": "PASS",
  "byte_identical_to_actual_runtime": true,
  "binary_sha256": "4643654625c71aaf5b45367bd3a198381257588d073ae0349fc90c1bc1b5ab85",
  "canonical_pack_sha256": "fa955025970d4a795c9f2eb547ea78040840bb5ecbc5576434061b95f8f3d7fb",
  "canonical_plan_sha256": "1b605baec500f7f20ac1c80a4bd7b9c7aadb9ce72ad264bbc0b4c2e49fa17571",
  "gap_receipt": {
    "authority_count": 2,
    "candidate_count": 1,
    "capability_id": "HTTP-SLI-ALERT-REFERENCE",
    "decision_state": "USE_CONDITIONED_PACK",
    "gate": "CAPABILITY-GAP-RESOLUTION-GATE",
    "implementation_ready": false,
    "observed_at": "2026-09-10T17:04:45.229912Z",
    "official_domain_count": 3,
    "outcome": "BLOCKED",
    "record_sha256": "d4eac2a7e591b1938192b8056cebe90a96c1ed3bdd7c02a4d62eee95096ed655",
    "requirement_ref": "markdown_system/FRANCHISE_GAP_MAP.md",
    "schema_version": 1,
    "search_count": 1,
    "selected_candidate_id": "authored-http-reference"
  },
  "gap_receipt_sha256": "ed4486b90e919abad602a28298511f62d8ec1ea44735015829d9109c979ca134"
}
```

## Canonical documentation successor after Preflight210

FAIL727 rejected an ambiguous edit anchor before writes. FAIL728: Preflight210 rejected missing canonical Apply order/Verification/Reconstruction evidence headings. Only prose headings/content were repaired; all eight file blocks retain their exact hashes, so the actual runtime executable is unchanged. FAIL729 was a shell-wrapper parse rejection before launch. Structural rerun accepted the section contract but correctly rejected stale ledger hashes (FAIL730); checkpoint210 precedes the next full verifier. Original receipts and earlier pack digest remain historical.

Current pack SHA-256: d16dddeb3aeae76b8efb89750c148df3ed4f9b94d1b997dc3df7fd90569ea3ee

```json
{
  "authority_count": 2,
  "candidate_count": 1,
  "capability_id": "HTTP-SLI-ALERT-REFERENCE",
  "decision_state": "USE_CONDITIONED_PACK",
  "gate": "CAPABILITY-GAP-RESOLUTION-GATE",
  "implementation_ready": false,
  "observed_at": "2026-09-10T17:04:45.229912Z",
  "official_domain_count": 3,
  "outcome": "BLOCKED",
  "record_sha256": "69575d19ad524dae7c453319333216a0fe7927334461e4218ef2961737aca451",
  "requirement_ref": "markdown_system/FRANCHISE_GAP_MAP.md",
  "schema_version": 1,
  "search_count": 1,
  "selected_candidate_id": "authored-http-reference"
}
```

## Final Preflight212 and checkpoint212

```json
{
  "executed_steps": 164,
  "all_executed_steps": "PASS",
  "profiles_composed": 56,
  "packs": 165,
  "materializable_files": 1515,
  "distributable_markdown": 823,
  "provenance": {
    "AUTHORED": 1268,
    "ADAPTED": 140,
    "VERBATIM": 107
  },
  "integral_profile": {
    "packs": 67,
    "files": 755
  },
  "http_reference_profile": {
    "packs": 68,
    "files": 763
  },
  "availability_status": "BLOCKED",
  "missing_tools": [
    {
      "id": "docker",
      "required_for": "PostgreSQL integration/recovery and optional Business Central Windows-container gates",
      "path": null,
      "version": "",
      "available": false
    }
  ],
  "json_sha256": "1f51f7c3ba55a6fecd718ccf4911e5cd95698d4949ce186b7e4b8e620e54ba35",
  "log_sha256": "4762395d12753043e815fbfb24ed60bbf2db61c1c4affd0044cad086a72aa525",
  "production_admission": false
}
```

FAIL728 and731 are regression-proven by the unchanged full gate. All164 executed steps PASS; availability remains BLOCKED solely for Docker. This is structural/local verification, not SCA for every runtime. During subsequent official web research, FAIL732 identifies Next16.3.2 Windows RCE and related image-processing advisories; web promotion and TEST03/07 remain blocked. The actual V375 runtime contains Go/PostgreSQL/Prometheus, does not execute Next, and retains its narrow proven scope.45/48 unchanged.
