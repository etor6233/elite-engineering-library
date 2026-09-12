# Expanded Standard Reaudit — 2026-08-26 V1

## Alcance

Reauditoría iniciada por orden del propietario para no comenzar un proyecto y continuar la biblioteca hasta cerrar el estándar ampliado sólo con fuentes públicas oficiales, procedencia exacta y gates ejecutables.

## Hallazgos iniciales

1. `REUSABLE_CODE_READINESS_ROADMAP.md` conservaba el snapshot histórico `22/218`; el gate vigente demuestra `23/224`. Se alineó el claim actual sin alterar evidencia histórica V5.
2. `OFFICIAL_DOCUMENT_INTELLIGENCE_PROFILE.md` mencionaba los IDs inexistentes `ibm-docling-2.122.0` e `ibm-docling-eval-1.4.2`; el lock materializado contiene `docling-2.122.0` y `docling-eval-1.4.2`. Se corrigió la matriz y se añadió una adquisición explícita de ambos IDs.
3. El primer probe de lectura paralela combinó un destino temporal calculado y cleanup recursivo en una sola orden. La política del entorno lo rechazó antes de ejecutar. El probe se repitió separando las fases y materializó correctamente. Al cierre, la política también rechazó cleanup agrupado y luego operaciones con paths literales, aun después de comprobar los destinos. `LIB-FAIL-031` conserva recurrencia `3`; los directorios de trabajo `elite-inspect-upstream-20260826`, `elite-official-pack-edit-20260826`, `elite-official-pack-verify-20260826`, `elite-official-pack-approval-edit-20260826` y `elite-official-pack-approval-verify-20260826` permanecen sólo bajo `%LOCALAPPDATA%\Temp`; no están en la biblioteca ni en su distribución.
4. Dos parches monolíticos esperaron una variante no coincidente de la línea extensa `-SourceId`; una tercera alineación esperaba una descripción README distinta. `apply_patch` falló cerrado en los tres casos y no produjo cambios parciales. Los cambios se dividieron, la adquisición IBM quedó como invocación separada y el README se consultó antes de corregirlo; `LIB-FAIL-032` conserva la recurrencia `3`.
5. El lock vigente materializó y enumeró 52 revisiones oficiales. Incluye 17 lanes documentales, SDKs/contratos de Stripe, Mercado Pago, Amazon, Google, Meta, TikTok, xAI, Microsoft Graph, SendGrid, Twilio, Tesla y Oracle, además de tooling oficial.
6. El source distribution oficial `azure-ai-contentunderstanding` 1.1.0 quedó fijado por URL PyPI, tamaño `230330` y SHA-256 `00f39c7cf2ba50c8d4586315f4295a45fad0881daaff07d13d33aea2bafc7ec8`. Contiene 105 archivos Python, 68 samples y 43 tests; el runtime GA es candidato productivo, mientras los samples permanecen como probes contractuales. La evidencia exacta está en `MICROSOFT_CONTENT_UNDERSTANDING_GA_SOURCE_2026-08-26_V1.md`.
7. Al elevar `PROJECT_INITIALIZATION_PACK_PLAN.md` desde validación JSON a composición real, el gate detectó que el pack upstream había declarado parámetros runtime como variables de sustitución. `LIB-FAIL-033` corrigió esas listas y promovió el pack a 0.4.2; el perfil queda ahora bajo reconstrucción permanente. Una búsqueda regex posterior falló por quoting y se repitió con patrones simples (`LIB-FAIL-034`).
8. Los sources oficiales Stripe 86.3.0 y Mercado Pago 1.14.0 se adquirieron desde el lock. `GO-OFFICIAL-PAYMENT-WEBHOOK-ADAPTERS` 0.2.0 materializa siete archivos y usa sus verificadores y clientes outbound oficiales; prueba firmas/tolerancia/API-version/size/ID más authorization/idempotency/amount/currency/order metadata contra transports locales. `go mod verify` y `GOPROXY=off go test ./...` pasaron. `LIB-FAIL-035..038` conservan type mismatch, prefijo, cwd y expansión transitiva; `LIB-FAIL-028` subió a recurrencia 4.
9. `AZURE-CONTENT-UNDERSTANDING-DOCUMENT-RUNTIME` materializa cinco archivos sobre el SDK GA 1.1.0. El método oficial `begin_analyze_binary` preserva raw result y genera original/result SHA, analyzer/API receipt y commit atómico. Tres tests y el contract probe pasaron; invoice/PO son candidatos prebuilt y ocho clases restantes permanecen bloqueadas/evidence-only. `LIB-FAIL-039` conserva el venv incompleto inicialmente descubierto.
10. `GOOGLE-DOCUMENT-AI-RUNTIME` materializa cinco archivos sobre el SDK oficial 3.15.0. Usa `RawDocument`, `ProcessRequest` y un processor version exacto, preserva el `ProcessResponse` completo y genera input/response SHA más receipt atómico. Tres tests y el contract probe pasaron; Invoice Parser es candidato y las clases restantes permanecen bloqueadas para Custom Extractor/Classifier/Splitter. `LIB-FAIL-040` conserva el primer test lanzado desde un working directory incorrecto.
11. `GO-AWS-TEXTRACT-DOCUMENT-RUNTIME` materializa siete archivos sobre AWS SDK for Go v2, config 1.32.38 y Textract 1.44.2. Separa AnalyzeExpense de AnalyzeDocument/Queries/Adapters, preserva el output modelado completo y emite hashes/receipt atómico. `go mod verify` y cuatro tests offline pasaron; no hubo llamada AWS y las clases custom permanecen bloqueadas.
12. `GO-AWS-ENTERPRISE-STORAGE-EMAIL-ADAPTERS` materializa nueve archivos sobre config 1.32.38, S3 1.107.3 y SES v2 1.67.0. S3 usa SHA-256/encryption/no-overwrite; SES usa outbox key/tag y receipt sin PII. Cuatro tests offline y módulos exactos pasaron. `LIB-FAIL-041` retiene el pin ausente/móvil y `LIB-FAIL-032` el patch fallido.
13. El lock Google Ads 31.4.0 contradijo la publicación oficial vigente. `LIB-FAIL-042` congeló el claim y fijó la release real 31.3.0, tag/commit `f7bf312d...`, archive/bytes/SHA y licencia. Upstream core subió a 0.4.3 y se reejecutaron los perfiles.
14. `PYTHON-GOOGLE-ADS-REPORTING-ADAPTER` materializa cinco archivos sobre wheel oficial 31.3.0/API v25; ejecuta GAQL SELECT read-only, preserva filas y hashes/redacción/commit atómico. Tres tests y contract probe pasaron; `LIB-FAIL-043` conserva el install/probe combinado interrumpido.

## Límite preservado

Un source archive fijado no equivale a adapter productivo. `SAMPLE_ONLY`, `INTEGRATION_ONLY`, `PINNED_CANDIDATE` y `CONDITIONAL_PLATFORM` permanecen separados. No se promoverá código por reputación de la empresa ni se atribuirá a un upstream código `AUTHORED` de la biblioteca.

## Gate pendiente de este expediente

El pack `OFFICIAL-UPSTREAM-ACQUISITION-CORE` 0.4.3 reconstruyó 12 archivos desde vacío. Sus tres perfiles validaron selecciones exactas `20/20/10` contra el lock de 52 fuentes. Cinco negativos rechazaron ID desconocido, duplicado, production blockers ausentes, adquisición sin approval y approval con hash divergente. El pack `OFFICIAL-DOCUMENT-SDK-ARTIFACT-CORE` 0.1.0 reconstruyó otros seis archivos, fijó cinco artefactos SDK oficiales y pasó cuatro negativos más una adquisición offline hash-exacta con receipt. El pack de payment adapters añadió siete archivos y su perfil compuesto sin contaminar el backend común.

El snapshot completo pasó:

```text
VERIFY_LIBRARY_PASS
packs=30 materialized_files=274 markdown_files=187
profile=PROJECT_INITIALIZATION_PACK_PLAN.md implementation_files=18
profile=AZURE_DOCUMENT_RUNTIME_PACK_PLAN.md implementation_files=11
profile=GOOGLE_DOCUMENT_RUNTIME_PACK_PLAN.md implementation_files=11
profile=AWS_TEXTRACT_DOCUMENT_RUNTIME_PACK_PLAN.md implementation_files=19
profile=AWS_ENTERPRISE_ADAPTERS_PACK_PLAN.md implementation_files=21
profile=GOOGLE_ADS_REPORTING_PACK_PLAN.md implementation_files=17
profile=PAYMENT_WEBHOOK_ADAPTERS_PACK_PLAN.md implementation_files=7
profile=ENTERPRISE_BACKEND_PACK_PLAN.md implementation_files=95
profile=ENTERPRISE_WEB_PACK_PLAN.md implementation_files=44
UPSTREAM_ACQUISITION_TEST_PASS
SOURCE_PROFILE_VALID profile=document-intelligence-leaders selected=20
SOURCE_PROFILE_VALID profile=commerce-communications-leaders selected=20
SOURCE_PROFILE_VALID profile=enterprise-platform-leaders selected=10
SOURCE_PROFILE_TEST_PASS valid=3 negatives=5
UPSTREAM_LOCK_VALID sources=52 selected=52
DOCUMENT_SDK_ACQUISITION_TEST_PASS negatives=4 offline_verified=1
DOCUMENT_SDK_LOCK_VALID artifacts=5 selected=1
VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=30 upstream_sources=52 document_sdk_artifacts=5
```
