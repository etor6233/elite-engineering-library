# Microsoft Content Processing Solution Accelerator main — reauditoría 2026-08-28 V1

## Veredicto

El commit oficial actual aporta código real y valioso para procesamiento multiarchivo, pero queda `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`. No se copió, corrigió ni presentó como implementación Elite propia. El source lock permite adquirir sus bytes exactos sólo para investigación aprobada y mantiene la release estable v2.1.2 como entrada separada.

## Identidad, licencia y archivo

- Repositorio oficial activo: `microsoft/content-processing-solution-accelerator`.
- Branch `main`; commit firmado/verificado `659eaa1f503dd08b1e1aea1c72eab11c7c191d00`; tree `a30887cd0de524d5f16db345af47ffcfa2fd9c6f`; fecha 2026-08-25.
- El commit es posterior a la última release estable observada, v2.1.2 (`b47cec48475cfedd7109debf2f407ef0c5530311`, publicada 2026-06-16); por lo tanto se registra como rama no publicada, no como release.
- Archive exacto: 18.361.248 bytes; SHA-256 `5f94e4b5e052ed736abd759e500fee42857f0b21e0d7228cc82ccbffd2026565`; 701 archivos y 177 entradas de directorio.
- MIT: `LICENSE`, 1.140 bytes, SHA-256 `d9a1b1e30d633d5732ea18e3cba9538d293ebc53e1a9e4e96ab739e0c5c4f1cb`.

## Toolchains y reproducción cerrada

- Node 22.23.2 exacto: ZIP 35.683.585 bytes, SHA-256 `1177b4137ba5adaa56354ae40f1080c7450e8ae09cecb47da459d1c52ac99f97`; satisface `>=22.22.0 <23`.
- pnpm 10.28.2 exacto vía Corepack y `pnpm-lock.yaml` frozen.
- uv 0.12.7 exacto: ZIP 16.979.508 bytes, SHA-256 `bf1518af459a3915511a11fdc6e2f43ef9a2afa138b9d498eeb9642fe9d85218`; Python 3.12.14 y los tres `uv.lock` usados sin elevar versiones.
- Se siguieron los comandos oficiales de `.github/workflows/test.yml`, no un `pytest` inventado desde otro cwd.

## Gates ejecutados

| Componente | Resultado exacto |
|---|---|
| ContentProcessor | 244 PASS; cobertura 86,37%; 7 warnings |
| ContentProcessorAPI | 42 PASS; cobertura 99,53% |
| ContentProcessorWorkflow | 268 PASS; 3 warnings |
| total Python | 554 PASS |
| ContentProcessorWeb install | frozen install PASS |
| ContentProcessorWeb tests | 9 suites PASS, 3 suites FAIL; 107 tests PASS; el fallo es carga ESM de `react-router` desde Jest/CommonJS |
| ContentProcessorWeb build | PASS; JS gzip 428,06 kB, CSS gzip 5,46 kB; deprecaciones PostCSS y caniuse-lite 7 meses desactualizado |

El primer rerun web con un separador `--` adicional queda descartado: produjo argumentos Jest inválidos y una imagen falsa de 12 suites fallidas. La repetición canónica, `pnpm test --watchAll=false --runInBand`, gobierna los números anteriores.

## SCA actual

- Web producción: 1.421 dependencias; 5 vulnerabilidades — 1 low y 4 moderate; IDs npm `1121360`, `1123965`, `1123966`, `1123977`, `1130709`.
- Processor: 205 paquetes auditados, 5 findings en click/h2/idna; `GHSA-65pc-fj4g-8rjx`, `GHSA-6hr6-w5qg-qmwg`, `PYSEC-2026-2132`, `PYSEC-2026-215`, `PYSEC-2026-3628`.
- API: 122 paquetes, 3 findings en click/idna; `GHSA-65pc-fj4g-8rjx`, `PYSEC-2026-2132`, `PYSEC-2026-215`.
- Workflow: 254 paquetes, 6 findings en h2/kafka-python/mem0ai; `GHSA-6hr6-w5qg-qmwg`, `GHSA-xqxw-r767-67m7`, `PYSEC-2026-2190`, `PYSEC-2026-2191`, `PYSEC-2026-2636`, `PYSEC-2026-3628`.

No se aplicó `audit fix`, no se elevaron locks y no se atribuyó una corrección local a Microsoft.

## Código oficial útil y límite humano demostrado

El source contiene Extract→Map→Evaluate→Save, schemas, scores, Blob/Cosmos, colas/workflows, retries/DLQ, gap analysis, API/UI, resultados editables, comentarios y fixtures multiarchivo. Cinco fixtures de reclamos quedaron fijados por path, bytes y hash en el catálogo.

La superficie de corrección actual no satisface el contrato empresarial requerido:

- `PUT /processed/{process_id}` acepta un resultado reemplazado o comentario; los modelos sólo contienen `process_id` más `modified_result` o `comment`.
- `CosmosContentProcess.update_process_result` y `update_process_comment` filtran únicamente por `process_id`; `update_document_by_query` ejecuta `update_one(query, {$set: update})`.
- `last_modified_by` se fija al literal `user`; no se persisten identidad EasyAuth, motivo obligatorio ni historial inmutable.
- El middleware lee headers EasyAuth y coloca `enduser.id` en telemetría, con fallback `anonymous`, pero esa identidad no llega al write model.
- El comentario de claim carga el agregado, cambia `process_comment` y llama `update_async` sin versión del caller. `sas-cosmosdb==0.1.5` usa Mongo `update_one({id,...predicate}, {$set:...})` y Cosmos SQL `replace_item(item=id, body=document)` sin ETag/match condition; la aplicación no entrega predicate/version.
- No existe lease/expiry de asignación ni dual control en estas rutas.

Por eso la edición visual es una primitive oficial útil, no una autorización para almacenar automáticamente datos empresariales corregidos.

## Condición de uso

`UP-FAIL-141` gobierna la admisión. El agente no puede componer ni desplegar esta rama, parchearla y atribuir el patch a Microsoft, ni afirmar exactitud productiva. Una release oficial posterior debe repetir identidad/licencia, frozen installs, suites, build, SCA, autenticación, reviewer/reason/corrections, lease, ETag/version conflict, request binding, evidencia durable, aislamiento, recovery y rollback. Azure AI/OpenAI, Blob, Cosmos y deployment requieren cuentas, regiones, políticas, presupuesto y autorización del usuario; no se realizaron llamadas live ni despliegues.
