# Document Intelligence Multi-Vendor Reaudit — V81

Fecha: 2026-08-28  
Estado: `RECONSTRUCTED / CONDITIONED / EXPANDED_STANDARD_NOT_READY`

## Objetivo

Determinar qué código documental público y licenciable de Microsoft, OpenAI, Oracle y SAP puede incorporarse sin inventar una implementación, y ampliar el lock ejecutable sólo cuando identidad, licencia, bytes e interfaces sean demostrables.

## Resultado admitido

`OFFICIAL-DOCUMENT-SDK-ARTIFACT-CORE` sube a `0.3.0`. Conserva seis archivos materializables y ahora fija doce artefactos oficiales: wheel+sdist de Azure Document Intelligence 1.0.2, Azure Content Understanding 1.1.0, Google Document AI 3.15.0, OpenAI 3.3.1, Oracle OCI 2.185.0 y el par Toolbox 0.17.3 condicionado.

Oracle OCI 2.185.0 se admite únicamente como `INTEGRATION_ONLY`:

- repositorio oficial [`oracle/oci-python-sdk`](https://github.com/oracle/oci-python-sdk/tree/e988c91dcc9963718454cb1e215bd44540524881), release/commit firmado `e988c91dcc9963718454cb1e215bd44540524881`;
- wheel PyPI oficial: 36.716.800 bytes, SHA-256 `93556f08c6270e2e174d59d94fcd49fa794ba3f4d60630af173bc9b6204b085c`;
- sdist PyPI oficial: 18.098.448 bytes, SHA-256 `375eec7c8568b0066bd5e6ff48cf406a4a1af34437c29abd7981f1ab78fac9c2`;
- metadata `Development Status :: 5 - Production/Stable`; licencia `UPL-1.0 OR Apache-2.0`;
- 115 archivos Python bajo `oci.ai_document` en el source fijado;
- Python 3.12.13: import de `AIServiceDocumentClient`, `AnalyzeDocumentDetails`, `DocumentKeyValueExtractionFeature`, `InvoiceProcessorConfig`, `BoundingPolygon`, `FieldValue` y `KeyValueDetectionConfidenceEntry` PASS;
- `pip check` PASS; manifest completo 23 paquetes SHA-256 `4e5e2d0ccc9dc9f7a577803beeadcd8245b536e5e4b0df8871726c83375a6dec`;
- pip-audit 2.10.1: 23 dependencias, 0 registros, 0 IDs y 0 skips; reporte SHA-256 `8986bb362e06a9aeb101996c863eb8af6c2a6434c11d8e9ba51f43a12638df86`.

Esto prueba adquisición e interfaces oficiales, no precisión de documentos ni un flujo de negocio completo. El SDK no aporta en su repositorio un sample AI Document ejecutable para redistribuir; no se creó uno local y no se atribuyó ninguno a Oracle.

## Hallazgos oficiales no admitidos como código inmediato

### Microsoft

La documentación oficial de [prebuilt analyzers de Azure Content Understanding](https://learn.microsoft.com/en-us/azure/ai-services/content-understanding/concepts/prebuilt-analyzers) cubre procurement mediante invoice, receipt y purchase order y permite copiar un analyzer prebuilt. El pack existente `MICROSOFT_AZURE_CONTENT_UNDERSTANDING_OFFICIAL_ANALYZER_COPY_SAMPLE` ya conserva código Microsoft exacto para crear/copiar analyzer, `GenerationMethod.EXTRACT`, fuente/confianza y geometría. No se encontró en la lista oficial un analyzer específico publicado para packing list, bill of lading o aduana; no se inventaron esos schemas.

El repositorio oficial [`Azure-Samples/data-extraction-using-azure-content-understanding`](https://github.com/Azure-Samples/data-extraction-using-azure-content-understanding/tree/0461a5549104ca769b8ec082c05894468997f300) fue rechazado como solución reusable inmediata: ZIP SHA-256 `23fb059d172cc9ebbe76cebc5461cf9150222fbc91830f35c0519ddaf2b7f7cc`, schema principal inválido, Functions `ANONYMOUS`, requirements sin lock y dependencia dev omitida, 81 findings flake8 y 11 registros/9 IDs de vulnerabilidad en tres paquetes. Al excluir la prueba bloqueada por la dependencia omitida, 165 tests pasaron; ese resultado parcial no revierte el rechazo.

### OpenAI

La documentación oficial vigente de [file inputs](https://developers.openai.com/api/docs/guides/pdf-files) y [Structured Outputs](https://developers.openai.com/api/docs/guides/structured-outputs) permite entrada de PDF/archivo y salida restringida por schema. El sdist oficial OpenAI 3.3.1, 1.282.113 bytes/SHA-256 `6f22807de1a976c932cecda620e8172a8c3fdbaeed29c7f21564e0c2410edf56`, contiene el ejemplo oficial `examples/responses/structured_outputs.py`, 1.050 bytes/SHA-256 `6318ddd61cbec436f0cf9b161a3b2cafbc7f4f98c5bd99bc4b0a43d7173cd6fd`. No contiene un ejemplo/test que combine `input_file` con Structured Outputs ni un engine de evidencia por campo. Se conserva `INTEGRATION_ONLY`; no se fabricó la unión ni se llamó código OpenAI.

### SAP y Oracle DevRel

[`SAP/business-document-processing`](https://github.com/SAP/business-document-processing/tree/5b17b6da8130721f950b6497c68021e75e712d95) está archivado, sin tags/releases y con último commit 2025-09-11: referencia histórica, no base actual.

[`oracle-devrel/oci-ai-invoice-handling`](https://github.com/oracle-devrel/oci-ai-invoice-handling/tree/f1f52bb99cdb91c4e4e8332caba2c051fccd638f) usa UPL pero no tiene tag/release, el commit no está firmado, el último push fue 2024-07-23 y el contenido es un export Oracle Integration limitado a invoice. Permanece `REJECTED_FOR_IMMEDIATE_ADOPTION`; no se confunde con el SDK OCI estable.

## Reconstrucción y gates

- pack 0.3.0 materializado desde Markdown: 6/6 archivos;
- tres bloques modificados comparados byte a byte contra staging: 3/3 hashes iguales;
- `DOCUMENT_SDK_ACQUISITION_TEST_PASS negatives=4 offline_verified=1`;
- `DOCUMENT_SDK_LOCK_VALID artifacts=12 selected=2` para wheel+sdist OCI;
- cinco perfiles consumidores alineados exactamente a `0.3.0`;
- `VERIFY_LIBRARY_PASS`: 59 packs, 548 archivos materializables, 391 Markdown antes de esta evidencia y 32 perfiles;
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: 59 packs, 110 fuentes upstream y `document_sdk_artifacts=12`.

## Fallos y límites retenidos

Las condiciones upstream 165–169 y los fallos locales 896–920 están en `markdown_system/LIBRARY_FAILURE_LEARNING_LEDGER.md`. Los gates no autorizados por red/cuenta siguen marcados `SKIPPED`, nunca convertidos en PASS.

El estándar ampliado continúa `NOT_READY`: no existe todavía código oficial admitido que cierre de punta a punta clasificación, schema, extracción, evidencia, revisión y persistencia para proforma, packing list, bill of lading, aduana y todas las demás clases empresariales. V81 agrega un SDK Oracle real y verificable; no afirma que esa brecha desapareció.
