# Azure Content Understanding Document Runtime — 2026-08-26 V1

## Código oficial usado

- `azure-ai-contentunderstanding==1.1.0`, wheel Microsoft GA hash `d1d6bdeffe02f5c8cc5309f0e120a1338f51af04638605f9617b7b501e2d1dd6`.
- API `2025-11-01`; `ContentUnderstandingClient.begin_analyze_binary` con `analyzer_id`, bytes locales y `content_type`.
- Source distribution oficial ya fijado: 230330 bytes, SHA-256 `00f39c7cf2ba50c8d4586315f4295a45fad0881daaff07d13d33aea2bafc7ec8`.

El runner es glue `AUTHORED`; SDK, LRO, modelos, field confidence/source y resultado son Microsoft.

## Resultado

```text
Materialized 5 files into ...\elite-azure-document-runtime-verify-20260826
Ran 3 tests in 0.024s — OK
AZURE_CU_SDK_CONTRACT_PASS version=1.1.0
```

El positivo envió un PDF fixture como bytes al contrato oficial inyectado, preservó resultado completo y escribió original/result SHA, analyzer/API version y receipt. Los negativos rechazaron analyzer divergente, contents vacío, output ocupado y extensión no permitida sin dejar destino parcial.

## Clases

El template incluye invoice y purchase order con los prebuilt IDs oficiales como candidatos. Proforma, packing list, BOL, delivery note, certificates y customs quedan `BLOCKED_CUSTOM_ANALYZER_REQUIRED`; commercial invoice queda baseline-only hasta evaluación. Todas conservan `automatic_storage=false`.

## Límite

No se contactó Azure ni se usaron secretos/costos. No hay claim de precisión sobre REVESTEX ni promoción de storage: faltan cuenta/región/analyzer real, corpus/ground truth, evaluación por campo, privacy/cost/load, review/correction y rollback. `LIB-FAIL-039` registra que un venv encontrado por nombre no contenía el SDK y que el probe se repitió sólo después de comprobar import/version.

Fecha/revisor: 2026-08-26 / Codex.
