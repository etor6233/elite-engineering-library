# Índice de tareas de mantenimiento

No contiene estados ni un backlog paralelo. El owner único de orden, dependencias y
finalización es REUSABLE_CODE_READINESS_ROADMAP.md §10.

| Referencia | Responsabilidad |
|---|---|
| T2801 | preparación, requisitos, autoridades y assurance |
| T2802 | recorridos transaccionales y equivalencias |
| T2803 | identidad, seguridad y supply chain |
| T2804 | frontend, ayuda, capacitación y soporte |
| T2805 | integraciones y comunicación |
| T2806 | documentos y datos |
| T2807 | IA gobernada y evals |
| T2808 | build, entornos y delivery |
| T2809 | observabilidad, carga y recuperación |
| T2810 | validación integrada y archivo final |

PROJECT_EXECUTION_STATE.json continúa enlazando el roadmap, no esta tabla, como tasks.
Consultar spec.md para requisitos LIB-R01–LIB-R09 y evidencia necesaria; no repetir
trabajos históricos sin cambio de inputs, contradicción o falta de evidencia aplicable.

## V402333 — retired unsupported comparative claim

```json
{
  "schema": "elite-retired-unaccepted-claim/v1",
  "id": "BENCH-01",
  "state": "NONE_WITH_REASON_FOR_LIBRARY_INFRASTRUCTURE",
  "reason": "The user acceptance requires ready local infrastructure, exact build and measured portable NEW/EXISTING. It does not require a comparative LLM token-saving claim. No measured token baseline exists; no saving or less-than-one-week claim is accepted. Keep measured BENCH02 and finite local workload as performance proof. Removing this retired hypothesis from active benchmarks does not remove any REQUIRED capability or executed test.",
  "original_benchmark": {
    "id": "BENCH-01",
    "hypothesis": "La entrada por mapas y materialización selectiva reduce trabajo repetido sin perder controles",
    "environment": {
      "os": "Windows local",
      "toolchains": "PowerShell 7.6.5/Python 3.14.4 en V287; registrar exactos en nueva medición",
      "hardware": "No caracterizado; pendiente para comparación"
    },
    "workload": {
      "baseline": "Ensayo comparable pendiente",
      "candidate": "NEW/EXISTING de maintenance-entry más contexto efectivo de agente"
    },
    "correctness_gate": "TEST-01 debe pasar antes de comparar; no extrapolar tiempos de procesos a tiempo de construcción del negocio",
    "metrics": [
      "Tiempo total y por etapa",
      "Lecturas/bytes de contexto",
      "Tokens sólo si medibles",
      "Reintentos y errores"
    ],
    "baseline": "OPEN: falta comparación controlada",
    "candidate": "V287 observó 9507/8997 ms de procesos; no constituye benchmark comparativo ni medición de tokens; V326/BENCH-02 mide latencia del tooling, no prueba todavía reducción de trabajo/tokens del agente. Final scope decision: historical comparative agent/token-saving hypothesis is NOT_REQUIRED for LIBRARY_INFRASTRUCTURE; no saving/speedup claim accepted. Keep planned historical record, exclude from required performance verification; BENCH02 and finite local load remain required and proven.",
    "evidence_ids": [],
    "status": "planned"
  },
  "prior_contract_sha256": "63e0771e147a53ecdf344367add97137d563716bf97437778e01dcdc2d8124d4"
}
```

## V402336 — final library infrastructure acceptance

```json
{
  "schema": "elite-library-final-distribution-acceptance/v1",
  "state": "PASS",
  "scope": "LIBRARY_INFRASTRUCTURE",
  "source_revision": 330,
  "metadata_revision": 331,
  "final_execution_revision": 336,
  "production_authorized": false,
  "distribution": {
    "state": "PASS_FINAL_PREFLIGHT_ARCHIVE_PORTABLE_NEW_EXISTING",
    "steps": [
      {
        "step": "final-preflight335",
        "exit_code": 0,
        "seconds": 468.767,
        "log_sha256": "ffbdb0cab27822230686291e0289d999cd3a3f36f1f4ff5ee7edbb25638951e2"
      },
      {
        "step": "final-library-archive335",
        "exit_code": 0,
        "seconds": 170.585,
        "log_sha256": "323b6a0b23aace29223062322601a8682e165c3d3b89f4aac27e538e13d06009"
      },
      {
        "step": "portable-new-bridge335",
        "exit_code": 0,
        "seconds": 1.042,
        "log_sha256": "489caa96ad668a6657d9ebe5e4f0eed73f0eafc8fc994bc43892919cd9967e92"
      },
      {
        "step": "portable-new-compositor335",
        "exit_code": 0,
        "seconds": 0.967,
        "log_sha256": "e8104dba50d467848d0b04f141786eb5fe368ca77e17b513280035daf73d8e7b"
      },
      {
        "step": "portable-new-profile335",
        "exit_code": 0,
        "seconds": 12.735,
        "log_sha256": "b254e7d63f2b178f79a3ed718e3bb39ed50629760f92dee96ff6aae59966b821"
      },
      {
        "step": "portable-existing-destination-rejection335",
        "exit_code": 1,
        "seconds": 10.958,
        "log_sha256": "cc7c483139ebf021459203e91237daf60c0bd319c9a3f83e564085827cd7bbf5"
      }
    ],
    "source_revision": 330,
    "metadata_revision": 331,
    "accounts_used": false,
    "production_authorized": false,
    "archive": {
      "path": "C:\\Users\\NL\\Desktop\\Elite Franchise Reference V402\\elite-library-v402-source330-meta331.zip",
      "sha256": "507c9f40442233226708a49aea98e9bc1bcb59ad5e3948bf8f0f0653d0033768",
      "bytes": 58881092,
      "files": 1264,
      "manifest_sha256": "ee2f6be6256342a4ba2cc4163375a34abfbe21cb5f33f988484741186a10d76c",
      "extracted_root": "C:\\Users\\NL\\Desktop\\Elite Franchise Reference V402\\portable-v402-smoke\\Elite Engineering Library",
      "no_local_approvals_or_secrets": true
    },
    "portable_reference": {
      "path": "C:\\Users\\NL\\Desktop\\Elite Franchise Reference V402\\portable-v402-smoke\\New Franchise\\reference-staging",
      "source_files": 1653,
      "exact": true,
      "existing_rejected_without_mutation": true
    },
    "canonical_guard_sha256": {
      "PROJECT_EXECUTION_STATE.json": "9aa0f6b7ea27848006be34aa3f7ce58ff479d9654579eb7297e7fc87a6f279b3",
      "PROJECT_EXECUTION_EVENTS.jsonl": "13a4f88591311501bea53a38eadc810ddb1085fdae6b219eb41bbe20f781773e",
      "PROJECT_ENGINEERING_CONTRACT.json": "c4153fabe7a040270d5370d935606480e8c91bec6d12151bd0e441e765e64266",
      "PROJECT_LIBRARY_READINESS_GATE.json": "7cec77a8a53cce4e9ed2711b060d175bbcbe2ad750796c62f4a8f9086e528da4",
      "PROJECT_FAILURE_LESSONS.md": "6b35fb80fd4d44ef5fd578e0210676532652ae5efed07e0faef279149e408496",
      "markdown_system/PUBLIC_RELEASE_INPUT_POLICY.md": "91d3b886b9dcaf014353faa1dd3538533272151fa30cdf9ee59a903504f3a19d",
      "markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md": "4aaa1fb70a613dace854645c2d053097c6d8b3173cb8c2f6db80aedb9e04e672"
    },
    "seconds": 684.934
  },
  "preflight_scope_assessment": {
    "status": "PASS_SELECTED_LIBRARY_INFRASTRUCTURE",
    "generic_preflight_status": "BLOCKED",
    "executed_checks_passed": 164,
    "executed_checks_failed": 0,
    "missing_or_rejected_global_tools": {
      "pnpm-11.25.0": {
        "id": "pnpm-11.25.0",
        "required_for": "web/BFF frozen install",
        "path": "C:\\Users\\NL\\AppData\\Roaming\\npm\\pnpm.ps1",
        "version": null,
        "available": false,
        "admission": "BLOCKED",
        "reason": "V401: published pnpm affected by GHSA-vwc7-r8mq-g2x9; ZIP-disabled candidate only has scoped evidence. Exact runtime admission required."
      },
      "docker": {
        "id": "docker",
        "required_for": "PostgreSQL integration/recovery and optional Business Central Windows-container gates",
        "path": null,
        "version": "",
        "available": false
      }
    },
    "pnpm_resolution": "General installed pnpm remains security-admission BLOCKED. The exact restricted private installer is separately admitted and was used in both successful current builds; no general pnpm admission or redistribution is inferred.",
    "docker_resolution": "Docker is absent. This selected local reference uses the admitted native PostgreSQL runtime, with actual migration/replay/activation/rollback/recovery receipts; optional container targets are outside the selected blueprint.",
    "qualification_receipt_sha256": {
      "durable-inputs326.json": "0e20b1a355e46fb6cf072a45d13e34df4a75bd1499d513352a4570a029060004",
      "durable-operational-inputs329.json": "6331ad79cf4b120d18a91d220f58582dfce67bd0a2f0110586dbca36de5b721a",
      "signed-release330-result.json": "3714dac97da24c4cd103430e4b2b39028543b562a0101634e8c328df11907e54",
      "independent-verification330.json": "ee1369eb6586d3892c2fbb5e2964a0696a951a64880ad34eaf25a84bf6b951d9",
      "activation330/result.json": "0f5bc23405e58e4a0dc1d1be680522f2a829a8f7b0e1fa5101554152a1c51032"
    },
    "generic_preflight_is_not_a_PASS": true,
    "no_check_relaxed": true
  },
  "preflight_receipt_sha256": "73fe2d11d1817f4af2ddc436d0c9ff71b0699d7b247bcd717e4b0a5a9b86d416",
  "actual_signed_reference_receipt_sha256": "ee1369eb6586d3892c2fbb5e2964a0696a951a64880ad34eaf25a84bf6b951d9",
  "all_public_archive_bytes_match_canonical": true,
  "no_required_library_work_remaining": true,
  "retired_comparative_token_claim": "BENCH01 remains unaccepted and recorded above, not a claimed PASS.",
  "deferred": "Only configured provider credentials and the exact user-deferred Daybreak/libxml2 branch. Target permanent operations/regulation retain their explicit acceptance scope."
}
```

T2801–T2810 and ARCA_INFRA: PROVEN_LOCAL. TEST02/03/07: passed in the selected local scope.48/48controls. Daybreak last dossier remains ACCESS_BLOCKED, no investigation. Library infrastructure work complete; no production authorization.

### Final checkpoint336 correction — FAIL986

The preceding archive335 acceptance is a valid frozen historical qualification, not a completed execution checkpoint. Its checkpoint was refused before event336 by two stale roadmap checkboxes. Both are now synchronized with the already proven local controls. Public roadmap/failure metadata changed after archive335, so generate a new distribution336 and repeat only archive verification/portable materialization. Existing executable335 checks, source330, builds/signature and runtime proof remain byte-identical.

FAIL-20260913-987 RESOLVED_CURSOR_FORMAT: same336 checkpoint refused a bare64hex release digest; schema requires sha256: prefix. Actual archive hash unchanged. Correct the cursor representation and finalizer template; no new artifact or test claim. Failed logs preserved; no event appended until validation succeeds.

## V402337 — final library infrastructure acceptance

```json
{
  "schema": "elite-library-final-distribution-acceptance/v1",
  "state": "PASS",
  "scope": "LIBRARY_INFRASTRUCTURE",
  "source_revision": 330,
  "metadata_revision": 331,
  "final_execution_revision": 337,
  "production_authorized": false,
  "distribution": {
    "state": "PASS_FINAL_PREFLIGHT_ARCHIVE_PORTABLE_NEW_EXISTING",
    "steps": [
      {
        "step": "final-library-archive336",
        "exit_code": 0,
        "seconds": 167.659,
        "log_sha256": "9fd2594d9c63748ac56f7a047f1902e102c8e28729a5b5988fc22ff6ee592aab"
      },
      {
        "step": "portable-new-bridge336",
        "exit_code": 0,
        "seconds": 1.191,
        "log_sha256": "58b7d2557358375f7b0425ffd6dc67cf6bbe058439b2211fb6aa4e5ea5deddde"
      },
      {
        "step": "portable-new-compositor336",
        "exit_code": 0,
        "seconds": 0.945,
        "log_sha256": "72cafb9eac6428a4b88fe3ece7cb2d00edab31a4afb600810acb1d7e3e5da311"
      },
      {
        "step": "portable-new-profile336",
        "exit_code": 0,
        "seconds": 13.238,
        "log_sha256": "bf2c897f9e25e8a44504a5484053c0803c00428ba7bda08f5a8f5034b2bc557e"
      },
      {
        "step": "portable-existing-destination-rejection336",
        "exit_code": 1,
        "seconds": 8.38,
        "log_sha256": "5d85d01d607ef1d8fe98025ec86d5dc8218068597a7d39a5689a8eea7341e333"
      }
    ],
    "source_revision": 330,
    "metadata_revision": 331,
    "accounts_used": false,
    "production_authorized": false,
    "reused_executable_checks": {
      "receipt": "final-preflight335.json",
      "sha256": "73fe2d11d1817f4af2ddc436d0c9ff71b0699d7b247bcd717e4b0a5a9b86d416",
      "checks": 164,
      "generic_tool_diagnostic": "BLOCKED",
      "public_delta_only": [
        "PROJECT_FAILURE_LESSONS.md",
        "REUSABLE_CODE_READINESS_ROADMAP.md"
      ],
      "all_other_previous_public_payload_hashes_unchanged": true
    },
    "archive": {
      "path": "C:\\Users\\NL\\Desktop\\Elite Franchise Reference V402\\elite-library-v402-source330-meta331-distribution336.zip",
      "sha256": "cd24757744d01dbc37b666495d659f76824a8a680015eee469187e975542591f",
      "bytes": 58881955,
      "files": 1264,
      "manifest_sha256": "8b831fe6ec0dec14c05992e361efa62979d69c05ff3615b08d192458caf33935",
      "extracted_root": "C:\\Users\\NL\\Desktop\\Elite Franchise Reference V402\\portable-v402-smoke336\\Elite Engineering Library",
      "no_local_approvals_or_secrets": true
    },
    "portable_reference": {
      "path": "C:\\Users\\NL\\Desktop\\Elite Franchise Reference V402\\portable-v402-smoke336\\New Franchise\\reference-staging",
      "source_files": 1653,
      "exact": true,
      "existing_rejected_without_mutation": true
    },
    "canonical_guard_sha256": {
      "PROJECT_EXECUTION_STATE.json": "89a85cef826c6d15657dc80d2bb04d20db1e67b4e7210810a8f582be95c13d0d",
      "PROJECT_EXECUTION_EVENTS.jsonl": "c210f0d0fcc593c110803850bc86540381144dcae382781f120c602f5f891cbd",
      "PROJECT_ENGINEERING_CONTRACT.json": "ddfbd504b827c732a2f7d87d302958ec95388b31259338e9cccea04beb7f35ac",
      "PROJECT_LIBRARY_READINESS_GATE.json": "7cec77a8a53cce4e9ed2711b060d175bbcbe2ad750796c62f4a8f9086e528da4",
      "PROJECT_FAILURE_LESSONS.md": "ac4cecaec7b0e71aa546e77aebbf06d3053cb8186065a94f8af092957e107b9e",
      "markdown_system/PUBLIC_RELEASE_INPUT_POLICY.md": "91d3b886b9dcaf014353faa1dd3538533272151fa30cdf9ee59a903504f3a19d",
      "markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md": "4aaa1fb70a613dace854645c2d053097c6d8b3173cb8c2f6db80aedb9e04e672"
    },
    "seconds": 207.494
  },
  "preflight_scope_assessment": {
    "status": "PASS_SELECTED_LIBRARY_INFRASTRUCTURE",
    "generic_preflight_status": "BLOCKED",
    "executed_checks_passed": 164,
    "executed_checks_failed": 0,
    "missing_or_rejected_global_tools": {
      "pnpm-11.25.0": {
        "id": "pnpm-11.25.0",
        "required_for": "web/BFF frozen install",
        "path": "C:\\Users\\NL\\AppData\\Roaming\\npm\\pnpm.ps1",
        "version": null,
        "available": false,
        "admission": "BLOCKED",
        "reason": "V401: published pnpm affected by GHSA-vwc7-r8mq-g2x9; ZIP-disabled candidate only has scoped evidence. Exact runtime admission required."
      },
      "docker": {
        "id": "docker",
        "required_for": "PostgreSQL integration/recovery and optional Business Central Windows-container gates",
        "path": null,
        "version": "",
        "available": false
      }
    },
    "pnpm_resolution": "General installed pnpm remains security-admission BLOCKED. The exact restricted private installer is separately admitted and was used in both successful current builds; no general pnpm admission or redistribution is inferred.",
    "docker_resolution": "Docker is absent. This selected local reference uses the admitted native PostgreSQL runtime, with actual migration/replay/activation/rollback/recovery receipts; optional container targets are outside the selected blueprint.",
    "qualification_receipt_sha256": {
      "durable-inputs326.json": "0e20b1a355e46fb6cf072a45d13e34df4a75bd1499d513352a4570a029060004",
      "durable-operational-inputs329.json": "6331ad79cf4b120d18a91d220f58582dfce67bd0a2f0110586dbca36de5b721a",
      "signed-release330-result.json": "3714dac97da24c4cd103430e4b2b39028543b562a0101634e8c328df11907e54",
      "independent-verification330.json": "ee1369eb6586d3892c2fbb5e2964a0696a951a64880ad34eaf25a84bf6b951d9",
      "activation330/result.json": "0f5bc23405e58e4a0dc1d1be680522f2a829a8f7b0e1fa5101554152a1c51032"
    },
    "generic_preflight_is_not_a_PASS": true,
    "no_check_relaxed": true
  },
  "preflight_receipt_sha256": "73fe2d11d1817f4af2ddc436d0c9ff71b0699d7b247bcd717e4b0a5a9b86d416",
  "actual_signed_reference_receipt_sha256": "ee1369eb6586d3892c2fbb5e2964a0696a951a64880ad34eaf25a84bf6b951d9",
  "all_public_archive_bytes_match_canonical": true,
  "no_required_library_work_remaining": true,
  "retired_comparative_token_claim": "BENCH01 remains unaccepted and recorded above, not a claimed PASS.",
  "deferred": "Only configured provider credentials and the exact user-deferred Daybreak/libxml2 branch. Target permanent operations/regulation retain their explicit acceptance scope."
}
```

T2801–T2810 and ARCA_INFRA: PROVEN_LOCAL. TEST02/03/07: passed in the selected local scope.48/48controls. Daybreak last dossier remains ACCESS_BLOCKED, no investigation. Library infrastructure work complete; no production authorization.
