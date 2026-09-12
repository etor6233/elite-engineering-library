# All Implementation Packs Materialization — 2026-08-27 V41

## Alcance del snapshot

Este snapshot sucede a V40 sin reescribirlo. Conserva 46 packs, 413 archivos materializables y 23 perfiles; agrega al lock una tercera implementación oficial de interfaz/procesamiento documental únicamente después de demostrar que tampoco supera admisión inmediata.

Estado canónico al cierre:

- 46 implementation packs;
- 413 bloques `CREATE` materializables: 401 `AUTHORED`, 3 `ADAPTED`, 9 `VERBATIM`;
- 23 perfiles de composición;
- 79 upstream source archives fijados por commit, bytes, SHA-256 y licencia;
- 10 source profiles; ninguno selecciona AWS RAPID, AWS Assess Workbench ni IBM Document Extraction Toolkit;
- 5 artefactos SDK documentales y 1 artefacto Business Central;
- 9 provider adapters, 1 orquestador documental y 1 log de evidencia;
- 338 fallos locales y 97 condiciones upstream preservados al corte;
- 293 Markdown y 8 scripts PowerShell canónicos después de incorporar este expediente.

## Cambio material

`OFFICIAL-UPSTREAM-ACQUISITION-CORE` sube a `0.4.32`. Su lock agrega `ibm-document-extraction-toolkit-d9144ec`, código oficial de IBM, sólo como `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`.

El archive exacto del commit unsigned `d9144ec82194f173c5a73cc29fc08b757f6eced1` mide 2.174.373 bytes y tiene SHA-256 `53dc40c6f5fcdbc729cc45abee487ff2476bcccd79464ca3a042097228c1f73d`. Licencia Apache-2.0, README, requirements, manifests y locks quedan identificados. El test del adquiridor verifica identidad, clasificación, licencia y hashes.

IBM lo publica como prototipo, declara roles no funcionales y no ofrece release. El cliente exacto instaló 1.731 paquetes, pero Jest falló 2/2 suites antes de ejecutar tests y el build falló por un import no resoluble. Los audits npm reportaron 5/77/33 vulnerabilidades en sus tres locks de aplicación, con 1/5/5 críticas; OSV reportó 130 filas afectadas y 108 paquete-versiones distintos. La evidencia completa está en `IBM_DOCUMENT_EXTRACTION_TOOLKIT_REAUDIT_2026-08-27_V1.md` y `UP-FAIL-095` a `UP-FAIL-097`.

## Regresiones ejecutadas

Desde materialización limpia del pack de adquisición:

```text
Materialized 19 files
UPSTREAM_ACQUISITION_TEST_PASS negatives=6
UPSTREAM_LOCK_VALID sources=79 selected=79
SOURCE_PROFILE_TEST_PASS valid=10 negatives=5 positives=1
```

Los diez perfiles conservaron sus selecciones 6, 1, 10, 7, 16, 35, 1, 25, 12 y 1. Ninguno incorpora las tres referencias rechazadas.

La verificación estructural global pasó después de añadir este expediente:

```text
VERIFY_LIBRARY_PASS
packs=46 materialized_files=413 markdown_files=293
```

El audit ejecutable global pasó con Go oficial 1.26.7 fijado:

```text
VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=46 upstream_sources=79 document_sdk_artifacts=5 business_central_artifacts=1 provider_adapters=9 document_orchestrators=1 evidence_logs=1
```

El PASS posterior demuestra el conteo final de 293 Markdown; V41 gobierna este estado hasta que un snapshot nuevo lo suceda.

## Estado honesto

La biblioteca continúa incompleta globalmente. RAPID, Assess Workbench e IBM Document Extraction Toolkit son referencias oficiales de evaluación/revisión ya descartadas, no código listo. Su lock evita búsquedas repetidas y selección accidental; no cubre aún la revisión humana documental productiva exigida.

Continúan abiertos revisión humana ejecutable, handoff idempotente al negocio, browser/a11y/performance, telemetry/SBOM/provenance/rollback y promoción por claim. No existe autorización para Git, ZIP final, cuentas, credenciales, despliegue ni servicios pagos.

## Fuentes oficiales

- <https://github.com/IBM/document-extraction-toolkit/tree/d9144ec82194f173c5a73cc29fc08b757f6eced1>
- <https://github.com/aws-samples/sample-assess-workbench/tree/79a57b557fce417bb2bbbcc0794a8f7a6ce9245f>
- <https://github.com/aws-samples/review-and-assessment-powered-by-intelligent-documentation/tree/68fa528c040f4184316a7b8fd214a8113656f9f1>
