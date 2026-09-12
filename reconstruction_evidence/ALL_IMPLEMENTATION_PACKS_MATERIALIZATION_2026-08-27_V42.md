# All Implementation Packs Materialization — 2026-08-27 V42

## Alcance del snapshot

Este snapshot sucede a V41 sin reescribirlo. Conserva 46 packs, 413 archivos materializables y 23 perfiles; agrega al lock un sample oficial firmado de Microsoft únicamente después de demostrar que no supera la admisión inmediata.

Estado canónico al cierre:

- 46 implementation packs;
- 413 bloques `CREATE` materializables: 401 `AUTHORED`, 3 `ADAPTED`, 9 `VERBATIM`;
- 23 perfiles de composición;
- 80 upstream source archives fijados por commit, bytes, SHA-256 y licencia;
- 10 source profiles; ninguno selecciona AWS RAPID, AWS Assess Workbench, IBM Document Extraction Toolkit ni Microsoft Expense Submission MCP;
- 5 artefactos SDK documentales y 1 artefacto Business Central;
- 9 provider adapters, 1 orquestador documental y 1 log de evidencia;
- 341 fallos locales y 100 condiciones upstream preservados al corte;
- 295 Markdown y 8 scripts PowerShell canónicos después de incorporar este expediente.

## Cambio material

`OFFICIAL-UPSTREAM-ACQUISITION-CORE` sube a `0.4.33`. Su lock agrega `microsoft-expense-submission-mcp-8c2cb6e`, código oficial de Microsoft, sólo como `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`.

El archive exacto del commit firmado/verificado `8c2cb6eed8d916dd4d8af355c55133be98da95fd` mide 316.870.113 bytes y tiene SHA-256 `6894c595e969603a8c20c33221dde040304d2a59c2130c6f05aa1ebfcbfde19a`. La licencia MIT tiene SHA-256 `275b4dd619de4e16a017b10d0beec72abbbbf14ee8a2fc68f8bdb398e821f623`. README, package, server, database y MCP tools quedaron identificados por hash en el lock y el test del adquiridor verifica identidad, clasificación, licencia y ausencia explícita de lock de dependencias.

La instalación mutable de desarrollo resolvió 519 paquetes, compiló el servidor y compiló los tres widgets oficiales. Eso no constituye un grafo reproducible: el sample no publica lock, tests ni lint, y la resolución actual informó cuatro vulnerabilidades moderadas. La revisión de código encontró una API key de desarrollo fija, cache OBO sólo en memoria, claims y URLs sensibles en logs, descargas de archivos sin límites ni scan, comprobantes en filesystem local y un único borrador global. El submit genera un identificador temporal pero no persiste reporte/aprobación y elimina borrador/archivos; no demuestra aislamiento, idempotencia, concurrencia ni integración ERP. La evidencia completa está en `MICROSOFT_EXPENSE_SUBMISSION_MCP_REAUDIT_2026-08-27_V1.md` y `UP-FAIL-098` a `UP-FAIL-100`.

## Regresiones ejecutadas

Desde materialización limpia del pack de adquisición:

```text
Materialized 19 files
UPSTREAM_ACQUISITION_TEST_PASS negatives=6
UPSTREAM_LOCK_VALID sources=80 selected=80
SOURCE_PROFILE_TEST_PASS valid=10 negatives=5 positives=1
```

Los diez perfiles conservaron sus selecciones 6, 1, 10, 7, 16, 35, 1, 25, 12 y 1. Ninguno incorpora las cuatro referencias rechazadas.

El audit ejecutable global pasó con Go oficial 1.26.7 fijado:

```text
VERIFY_LIBRARY_PASS
packs=46 materialized_files=413 markdown_files=294
VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=46 upstream_sources=80 document_sdk_artifacts=5 business_central_artifacts=1 provider_adapters=9 document_orchestrators=1 evidence_logs=1
```

El primer conteo de 294 Markdown precedió a este expediente. La verificación estructural final posterior produjo exactamente:

```text
VERIFY_LIBRARY_PASS
packs=46 materialized_files=413 markdown_files=295
```

Este PASS posterior demuestra el conteo final; V42 gobierna este estado hasta que un snapshot nuevo lo suceda.

## Fallo local convertido en aprendizaje

`LIB-FAIL-341` preserva el rechazo preventivo de un comando que intentaba limpiar recursivamente una ruta temporal calculada antes de resolverla. La regresión creó una raíz exacta ausente, resolvió y comprobó que su padre era el directorio temporal permitido, ejecutó materialización y suites, inventarió 21 entradas y eliminó exclusivamente ese literal. No hubo mutación del workspace durante el intento rechazado.

## Estado honesto

La biblioteca continúa incompleta globalmente. RAPID, Assess Workbench, IBM Document Extraction Toolkit y Microsoft Expense Submission MCP son referencias oficiales ya descartadas para adopción inmediata, no código listo. Su lock evita búsquedas repetidas, deriva adquisición exacta y bloquea selección accidental.

Continúan abiertos revisión humana documental ejecutable y durable, handoff idempotente al negocio, browser/a11y/performance, telemetry/SBOM/provenance/rollback y promoción por claim. No existe autorización para Git, ZIP final, cuentas, credenciales, despliegue ni servicios pagos.

## Fuente oficial

- <https://github.com/microsoft/mcp-interactiveUI-samples/tree/8c2cb6eed8d916dd4d8af355c55133be98da95fd/mcp-apps/expense-submission/node>
