# All Implementation Packs Materialization — 2026-08-27 V40

## Alcance del snapshot

Este snapshot sucede a V39 sin reescribirlo. Conserva 46 packs, 413 archivos materializables y 23 perfiles de composición; agrega una segunda implementación oficial AWS de evaluación/revisión documental al lock de investigación únicamente después de auditarla y rechazar su adopción inmediata.

Estado canónico al cierre:

- 46 implementation packs;
- 413 bloques `CREATE` materializables: 401 `AUTHORED`, 3 `ADAPTED`, 9 `VERBATIM`;
- 23 perfiles de composición;
- 78 upstream source archives fijados por commit, bytes, SHA-256 y licencia;
- 10 source profiles; ninguno selecciona AWS RAPID ni AWS Assess Workbench;
- 5 artefactos SDK documentales y 1 artefacto Business Central;
- 9 provider adapters, 1 orquestador documental y 1 log de evidencia;
- 337 fallos locales y 94 condiciones upstream preservados al corte;
- 291 Markdown y 8 scripts PowerShell canónicos después de incorporar este expediente.

## Cambio material

`OFFICIAL-UPSTREAM-ACQUISITION-CORE` sube a `0.4.31`. Su lock agrega `aws-assess-workbench-79a57b5`, código oficial de `aws-samples/sample-assess-workbench`, sólo como `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`.

El archive exacto del commit unsigned `79a57b557fce417bb2bbbcc0794a8f7a6ce9245f` mide 1.665.545 bytes y tiene SHA-256 `cab872e5919eadffac962527f9347fc098badbf875343aa179462c31467a966a`. La licencia MIT-0, README, manifiestos y locks Python/npm quedan identificados; el test del adquiridor vuelve a verificar identidad, clasificación, licencia y hashes.

El frontend exacto pasó lint, Knip limpio, 47 tests y build. Eso no alcanza para promoverlo: la suite Python quedó en 459 PASS, 17 FAIL y 6 skip, y OSV-Scanner 2.5.1 reportó 43 registros en 9 paquetes de los locks. No se usaron fixes automáticos, credenciales, Terraform, despliegue AWS ni llamadas con costo. La decisión y sus límites están en `AWS_ASSESS_WORKBENCH_REAUDIT_2026-08-27_V1.md` y `UP-FAIL-092` a `UP-FAIL-094`.

## Regresiones ejecutadas

Desde materialización limpia del pack de adquisición:

```text
Materialized 19 files
UPSTREAM_ACQUISITION_TEST_PASS negatives=6
UPSTREAM_LOCK_VALID sources=78 selected=78
SOURCE_PROFILE_TEST_PASS valid=10 negatives=5 positives=1
```

Los diez perfiles se validaron otra vez contra el lock exacto y conservaron sus selecciones: 6, 1, 10, 7, 16, 35, 1, 25, 12 y 1 fuentes. RAPID y Assess Workbench no pertenecen a ninguna selección ejecutable.

La verificación estructural global pasó:

```text
VERIFY_LIBRARY_PASS
packs=46 materialized_files=413 markdown_files=291
```

Los 23 perfiles compusieron con sus conteos vigentes: 25, 23, 47, 64, 64, 59, 36, 25, 24, 30, 28, 24, 27, 27, 30, 28, 10, 26, 27, 7, 6, 95 y 44 archivos.

El audit ejecutable global pasó con Go oficial 1.26.7 fijado:

```text
VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=46 upstream_sources=78 document_sdk_artifacts=5 business_central_artifacts=1 provider_adapters=9 document_orchestrators=1 evidence_logs=1
```

El audit volvió a ejecutar adquisición/perfiles, contratos Tessera, lock de SDKs, Business Central offline, MarkItDown, routing, secure-file, adapters AWS/GCS, WhatsApp, y materialización de Durable/Azure/Google/strict-field. Los runtimes que requieren red, cuenta o Linux conservaron su `SKIPPED` explícito; no se convirtió una omisión autorizada en PASS productivo.

## Estado honesto

La biblioteca continúa incompleta globalmente. Esta ampliación mejora la capacidad del agente para descubrir y descartar rápidamente dos referencias oficiales AWS de revisión humana ya auditadas; no entrega todavía una revisión humana productiva admitida. Assess Workbench no es un sustituto de una implementación promovida y no puede componer, desplegar ni justificar garantías AWS.

El siguiente trabajo conserva los blockers reales: revisión humana documental ejecutable, handoff idempotente al negocio, browser/a11y/performance, telemetry/SBOM/provenance/rollback y promoción por claim. No existe autorización para Git, ZIP final, cuentas, credenciales, despliegue ni servicios pagos.

## Fuentes oficiales

- <https://github.com/aws-samples/sample-assess-workbench/tree/79a57b557fce417bb2bbbcc0794a8f7a6ce9245f>
- <https://github.com/aws-samples/review-and-assessment-powered-by-intelligent-documentation/tree/68fa528c040f4184316a7b8fd214a8113656f9f1>
- <https://github.com/microsoft/content-processing-solution-accelerator/tree/659eaa1f503dd08b1e1aea1c72eab11c7c191d00>
- <https://github.com/SAP-samples/btp-cap-dox-invoice-validation/tree/990f4d8650c12d39c9a84974091c1b248a549ab6>
