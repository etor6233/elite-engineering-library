# Microsoft Kiota OpenAPI Client Gate — 2026-09-05 V247

## Decisión

Microsoft Kiota 1.35.0 queda admitido como `ELITE_REFERENCE / REBUILD_VERIFIED / CONDITIONED` para un claim estrecho: adquirir el CLI oficial autocontenido Windows x64 y generar clientes Go desde OpenAPI local hash-locked. No se atribuye a Microsoft el acquisition/generator/verifier glue `AUTHORED` de Elite.

## Identidad oficial

- repositorio: `microsoft/kiota`;
- release: `v1.35.0`, publicada 2026-09-04;
- commit: `114aa7ee609262d892fd9ceb02b2d9f7ecb84190`, firma GitHub verificada;
- source ZIP: 3.580.264 bytes, SHA-256 `9a9856310e11e5c0877b8b3a732dee78e051d67e5809875c560b46e6c5541248`;
- árbol: 1.073 archivos/8.798.877 bytes, 267 archivos bajo tests;
- MIT LICENSE SHA-256 `c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383`;
- asset oficial `win-x64.zip`: 37.140.722 bytes, SHA-256 `66b5547b948f7be724fa5e0ddddebda8f7b662574ae58eb9d3b8c51c609c6271`;
- `kiota.exe`: 89.879.984 bytes, SHA-256 `dffa90d51f5068c0fdd51dd9b539bce2a4c028335c09f449d9b2aa1f5f3d5d71`;
- versión ejecutada: `1.35.0+114aa7ee609262d892fd9ceb02b2d9f7ecb84190`.

## Demostración reproducible

El fixture oficial `tests/Kiota.Builder.IntegrationTests/ToDoApi.yaml` quedó fijado en 1.655 bytes/SHA-256 `9141520b53b6f9fcf544006290e912f7e02ffee9eab79b4e9907d694b8b18b49`.

Dos generaciones independientes produjeron los mismos seis archivos de código, 28.750 bytes, con hash de manifest `d6dcb4092c1b83ecd84ae864729f9a0df4ab1114aca8b3f54508d1a4d3cd2a53`. `kiota-lock.json` fue validado semánticamente aparte porque su `descriptionLocation` conserva legítimamente el path relativo del input; ningún archivo de código fue excluido.

```text
KIOTA_ACQUISITION_PASS version=1.35.0+114aa7ee609262d892fd9ceb02b2d9f7ecb84190
KIOTA_GENERATION_PASS files=7 code_files=6 code_tree_sha256=d6dcb409…
go test ./... PASS (3 packages)
go vet ./... PASS
OSV Scanner 2.5.1 packages=14 findings=0
negative_cases=2: altered OpenAPI hash; occupied output
MICROSOFT_KIOTA_OPENAPI_GATE_PASS
VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=154 upstream_sources=121 execution_state_controls=1 microsoft_kiota_openapi_client_gates=1 microsoft_playwright_browser_gates=1 google_lighthouse_web_quality_gates=1 provider_adapters=16 document_orchestrators=1 evidence_logs=1
```

El grafo Go se congeló mediante `go.mod`/`go.sum`; `go test` compila los tres packages generados y `go vet` pasó. OSV consultó 14 paquetes sin vulnerabilidades conocidas en esta ejecución. Un resultado vacío no garantiza ausencia de vulnerabilidades futuras.

## Fallo que mejoró el gate

El primer hash dorado incluyó `kiota-lock.json` y cambió entre dos raíces porque Kiota registra `descriptionLocation`. El gate no ocultó el delta: demostró que el único archivo distinto era el lock, verificó que los seis archivos Go eran byte-idénticos y separó hash portable de código de validación semántica del lock.

## Límites conservados

No quedan demostrados por Kiota: autenticación/autorización, secretos, deadlines, retry budget, idempotencia, webhooks, reconciliación, reglas de negocio, compatibilidad con un proveedor real, carga, seguridad ofensiva, despliegue o producción. El proyecto debe fijar su propio contrato, revisar el diff, congelar dependencias y ejecutar contract/E2E/live gates. Otros lenguajes de Kiota quedan fuera de esta admisión.

Fuentes oficiales: <https://github.com/microsoft/kiota>, <https://github.com/microsoft/kiota/releases/tag/v1.35.0> y <https://learn.microsoft.com/openapi/kiota/overview>.
