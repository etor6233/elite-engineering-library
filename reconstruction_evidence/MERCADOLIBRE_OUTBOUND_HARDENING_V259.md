# Evidencia V259 — correcciones verificadas del outbound Mercado Libre

Registro estructurado; los resultados detallados corresponden a los comandos ejecutados en esta revisión.

```json
{
  "revision": "V259",
  "date": "2026-09-05",
  "classification": "LIBRARY_MAINTENANCE",
  "pack": "GO-MERCADOLIBRE-QUESTION-OUTBOUND",
  "version": "0.1.1",
  "provenance": "AUTHORED",
  "admission": "CONDITIONED",
  "files": [
    {
      "Path": "C:\\Users\\NL\\Desktop\\Public Elite Codes\\implementation_packs\\GO_MERCADOLIBRE_QUESTION_OUTBOUND.md",
      "Hash": "45F39471A47DB8EB03C4E7ED865B4A4CB9CA73C53696753598A15CA92EDB364F"
    },
    {
      "Path": "C:\\Users\\NL\\Desktop\\Public Elite Codes\\markdown_system\\FRANCHISE_COMPLETE_PACK_PLAN.md",
      "Hash": "C230DB882D53997B4B2C1B567766E09FE7ECE83F94932F5A0B7A7E86037E5C63"
    },
    {
      "Path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-v259-ecad69896f04440fb45b03ac3f87d17a\\roundtrip\\internal\\outbounddelivery\\mercadolibre_question.go",
      "Hash": "5C30E9937BAD08DB395D98243198740E17F186A1DA63C64AD49065E35587C6F0"
    },
    {
      "Path": "C:\\Users\\NL\\AppData\\Local\\Temp\\elite-v259-ecad69896f04440fb45b03ac3f87d17a\\roundtrip\\internal\\outbounddelivery\\mercadolibre_question_test.go",
      "Hash": "3243D098D830C574F11F1495C6C107939EBE3B2AD2DCEA0414778A845E05BB28"
    }
  ],
  "authorities": [
    {
      "url": "https://pkg.go.dev/strconv#ParseInt",
      "claim": "ParseInt reports range errors; checked conversion replaces unchecked accumulation."
    },
    {
      "url": "https://aws.amazon.com/builders-library/making-retries-safe-with-idempotent-APIs/",
      "claim": "A failed request can have uncertain side effects; reconciliation and idempotency need explicit contracts."
    },
    {
      "url": "https://developers.mercadolibre.com.ar/en_us/categories-and-attributes/manage-questions-and-answers",
      "claim": "Existing V258 official contract: POST answers, integer question_id, GET v4 and validation errors.",
      "evidence": "reconstruction_evidence/MERCADOLIBRE_QUESTION_OUTBOUND_2026-09-05_V258.md"
    }
  ],
  "before": {
    "command": "go test ./internal/outbounddelivery -run TestMercadoLibre(Numeric|Uncertain|Untrusted) -count=1",
    "exit_code": 1,
    "failures": [
      "Five invalid numeric IDs reached HTTP.",
      "Nine uncertain HTTP statuses changed state to failed_terminal.",
      "Five incomplete/mismatched lookups changed store to failed_terminal."
    ]
  },
  "after": {
    "fresh_composition": {
      "packs": 67,
      "files": 704
    },
    "outbound_top_level_tests": 14,
    "full_go_tests": "PASS (database suites repeated separately with live loopback PostgreSQL)",
    "go_vet": "PASS",
    "go_build": "PASS",
    "postgres": {
      "version": "18.6",
      "migration_count": 49,
      "packages": [
        "internal/platform/postgres",
        "internal/app",
        "internal/outbounddelivery"
      ],
      "result": "PASS",
      "stopped": true
    },
    "roundtrip": "3/3 identical",
    "target_directory": "C:/Users/NL/AppData/Local/Temp/elite-v259-ecad69896f04440fb45b03ac3f87d17a/after"
  },
  "production_ready": false,
  "library_verifier": {"result": "VERIFY_LIBRARY_PASS", "exit_code": 0, "packs": 160, "materialized_files": 1395, "markdown_files": 692},
  "live_provider_calls": false,
  "progress_percentage": null,
  "percentage_correction": "85% and 25/41 were derived from textual labels, not workload or verified end-to-end closure; withdrawn.",
  "remaining": [
    "Connected browser-to-backend-to-database reference journey",
    "Release-linked help/training/support",
    "Journey telemetry and privacy proof",
    "Compact bootstrap artifact generation",
    "Required provider mutations and notifications reconciliation",
    "Project accounts, corpus, fiscal rules and production acceptance"
  ],
  "limitations": [
    "HTTPDoer must be configured with bounded deadlines and without redirects or automatic POST retries.",
    "Runtime caller still supplies real approvals, credentials and provider conditions.",
    "No claimed equivalence to all enterprise code, and no guarantee of production readiness."
  ]
}
```
