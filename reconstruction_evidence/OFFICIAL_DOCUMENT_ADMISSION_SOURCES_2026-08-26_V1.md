# Official Document Admission Sources — 2026-08-26 V1

## Resultado

Se auditaron y fijaron dos fuentes oficiales adicionales de empresas líderes para cerrar partes concretas del flujo documental que los tres runtimes de extracción no implementan por sí solos. No hubo llamadas a proveedores, credenciales ni costo. La evidencia admite código fuente exacto y pruebas locales; no afirma precisión sobre documentos de un proyecto ni autoriza almacenamiento automático.

## Google Cloud Document Intake Accelerator

- repositorio oficial: `GoogleCloudPlatform/document-intake-accelerator`;
- commit: `956e3cc7900338d1146d82c5ee2bd4b2551f70bb`, fechado 2022-11-08;
- archive: 43.108.717 bytes, SHA-256 `38204b7365ed2df6b00633c63943d588fcb085d876e5decd7cd96d99ca2309bb`;
- licencia Apache-2.0, SHA-256 `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30`;
- `microservices/adp_ui/package-lock.json` SHA-256 `722a9cee3d481c92f51fa3ef7d1335b8babc3bbca0545c82670e6ee71716b703`;
- 596 archivos, incluidos 163 Python; los 163 parsearon con Python 3.12 usando rutas largas sin escribir bytecode;
- código observado para clasificación, extracción, validación, profile matching, autoapproval configurable, estados Accept/Reject/Need Review, cola y edición humana.

Es `SAMPLE_ONLY / CONDITIONED`: el snapshot no es reciente, sólo el frontend tiene lock integral observado y el flujo requiere GCP, GKE, Firestore, BigQuery, processors, identidades, billing y pruebas del target.

## AWS OCR Evaluation Workbench

- repositorio oficial: `aws-samples/ocr-with-aws-ai-services`;
- commit: `5341910c3aee2da872f2966caed0fe0246b8a55c`, fechado 2026-08-10;
- archive: 4.127.057 bytes, SHA-256 `0e419c085a1957a37de791050a72e004b3d54b91f392ccf0311a4a5be1d2ac8f`;
- licencia MIT-0, SHA-256 `7cb750713252efd1d578837ba8785b61319109906f0c10b87adab7cf4badfc42`;
- 95 archivos, 55 Python y 26 bajo tests;
- Python 3.12 con los pins oficiales `gradio==5.31.0`, `Pillow==11.2.1` y `numpy==2.2.6`: 667 tests PASS; una prueba adicional no pudo crear un symlink por privilegio Windows y quedó conservada, no omitida del relato;
- tres bundles sintéticos oficiales coincidieron exactamente con su generador;
- código observado para JSON Schema, truth por documento, comparación campo por campo, reporte de truncamiento como error, seguridad de preview paths, costo/tiempo/accuracy e historial JSON/JSONL por ejecución.

Es `SAMPLE_ONLY / CONDITIONED`: `requirements.txt` no congela todo el grafo, Ruff no está fijado y Ruff 0.16.4 reportó 633 observaciones. Las llamadas Textract, Bedrock o BDA requieren cuenta, IAM, región, buckets, presupuesto y datos autorizados.

## Integración en Elite

`OFFICIAL-UPSTREAM-ACQUISITION-CORE` 0.4.9 fija 59 sources. Su perfil `document-intelligence-leaders` selecciona 22, incluidas estas dos nuevas fuentes, y pasa validación fail-closed. La adquisición continúa exigiendo respuestas del usuario, aprobación enlazada a los hashes exactos y reconocimiento de blockers; no usa Git en la biblioteca ni convierte samples en componentes productivos.

## Verificaciones

```text
UPSTREAM_ACQUISITION_TEST_PASS
SOURCE_PROFILE_TEST_PASS valid=3 negatives=5
UPSTREAM_LOCK_VALID sources=59 selected=59
SOURCE_PROFILE_VALID profile=document-intelligence-leaders selected=22
AWS pytest: 667 passed, 1 deselected por privilegio symlink Windows
AWS synthetic fixtures: 3 bundles exactos
Google Python AST: 163 files, 0 failures
```

Fuentes oficiales: https://github.com/GoogleCloudPlatform/document-intake-accelerator/tree/956e3cc7900338d1146d82c5ee2bd4b2551f70bb y https://github.com/aws-samples/ocr-with-aws-ai-services/tree/5341910c3aee2da872f2966caed0fe0246b8a55c.
