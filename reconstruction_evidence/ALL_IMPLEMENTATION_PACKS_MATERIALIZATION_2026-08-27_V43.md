# All Implementation Packs Materialization — 2026-08-27 V43

## Alcance del snapshot

Este snapshot sucede a V42 sin reescribirlo. Conserva 46 packs, 413 archivos materializables, 23 perfiles y 80 archives; profundiza el archive Microsoft ya fijado para auditar Approvals Box y fortalece el adquiridor con colecciones de artefactos adicionales exactos, sin duplicar source ni descarga.

Estado canónico al cierre:

- 46 implementation packs;
- 413 bloques `CREATE` materializables: 401 `AUTHORED`, 3 `ADAPTED`, 9 `VERBATIM`;
- 23 perfiles de composición;
- 80 upstream source archives fijados por commit, bytes, SHA-256 y licencia;
- 10 source profiles; ninguno selecciona AWS RAPID, AWS Assess Workbench, IBM Document Extraction Toolkit ni los samples Microsoft Expense/Approvals;
- 5 artefactos SDK documentales y 1 artefacto Business Central;
- 9 provider adapters, 1 orquestador documental y 1 log de evidencia;
- 342 fallos locales y 103 condiciones upstream preservados al corte;
- 297 Markdown y 8 scripts PowerShell canónicos después de incorporar este expediente.

## Cambio material

`OFFICIAL-UPSTREAM-ACQUISITION-CORE` sube a `0.4.34`. El lock sigue teniendo 80 archives: `microsoft/mcp-interactiveUI-samples` ya estaba fijado en el commit firmado `8c2cb6eed8d916dd4d8af355c55133be98da95fd`, por lo que Approvals Box no se cuenta ni se descarga como una fuente distinta. Nueve artefactos adicionales del componente quedaron registrados por path relativo seguro, propósito y SHA-256.

El adquiridor ahora:

- valida paths relativos canónicos de licencia, artefactos opcionales y adicionales;
- rechaza roots, segmentos vacíos, `.` y `..` antes de red;
- exige colecciones adicionales no vacías, paths únicos, hash lowercase y propósito;
- verifica cada artefacto adicional dentro del root autenticado después de expandir;
- registra cuántos artefactos adicionales fueron verificados en el receipt.

La regresión nueva cambia el resumen a siete negativos y demuestra que `../outside.md` es rechazado antes de adquisición.

## Resultado de Approvals Box

El componente oficial ofrece widgets y tools de lista, detalle, riesgo, comentarios, approve/reject/bulk/create. No supera admisión:

- no publica lock ni tests;
- `npm install` resolvió 144 paquetes de forma mutable; npm audit y OSV quedaron en cero para esa resolución, sin convertirla en grafo oficial;
- `npm run build` y `npx tsc --noEmit` expusieron 32 errores; el script publicado usa `tsc ... || true`, que oculta fallo en POSIX y no es portable a Windows;
- la frontera MCP no autentica y usa el gerente demo si no resuelve actor;
- `/stats` no exige autorización;
- approve/reject separan lectura, update y auditoría sin transacción, CAS/version ni idempotency key;
- SQLite local, bulk parcial y adjuntos no demuestran tenant, outbox, retención, backup/recovery ni handoff empresarial.

La evidencia completa está en `MICROSOFT_APPROVALS_BOX_REAUDIT_2026-08-27_V1.md` y `UP-FAIL-101` a `UP-FAIL-103`.

## Regresiones ejecutadas

Desde materialización limpia del pack de adquisición:

```text
Materialized 19 files
UPSTREAM_ACQUISITION_TEST_PASS negatives=7
UPSTREAM_LOCK_VALID sources=80 selected=80
SOURCE_PROFILE_TEST_PASS valid=10 negatives=5 positives=1
```

Los diez perfiles conservaron sus selecciones 6, 1, 10, 7, 16, 35, 1, 25, 12 y 1.

La verificación previa a este expediente produjo:

```text
VERIFY_LIBRARY_PASS
packs=46 materialized_files=413 markdown_files=296
VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=46 upstream_sources=80 document_sdk_artifacts=5 business_central_artifacts=1 provider_adapters=9 document_orchestrators=1 evidence_logs=1
```

La verificación estructural final posterior produjo exactamente:

```text
VERIFY_LIBRARY_PASS
packs=46 materialized_files=413 markdown_files=297
```

Este PASS posterior demuestra el conteo final; V43 gobierna este estado hasta que un snapshot nuevo lo suceda.

## Estado honesto

La biblioteca continúa incompleta globalmente. Ninguno de los cuatro componentes oficiales de revisión/handoff evaluados satisface los gates productivos. Sus fuentes exactas quedan disponibles para investigación y para evitar repetir callejones sin salida, no para composición automática.

Continúan abiertos revisión humana documental durable, handoff idempotente al negocio, browser/a11y/performance, telemetry/SBOM/provenance/rollback y promoción por claim. No existe autorización para Git, ZIP final, cuentas, credenciales, despliegue ni servicios pagos.

## Fuente oficial

- <https://github.com/microsoft/mcp-interactiveUI-samples/tree/8c2cb6eed8d916dd4d8af355c55133be98da95fd/oai-apps-sdk/approvals-box/node>
