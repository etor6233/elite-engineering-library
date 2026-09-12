# Google Document AI Runtime — reconstruction evidence V1

## Alcance

Reconstrucción local del pack `GOOGLE-DOCUMENT-AI-RUNTIME` 0.1.0. El runtime es glue `AUTHORED`; usa interfaces públicas del SDK oficial Apache-2.0 `google-cloud-documentai==3.15.0` y no copia su implementación. No se realizó una llamada a Google Cloud ni se afirmó precisión productiva.

## Entorno y artefacto oficial

- fecha: 2026-08-26;
- host de auditoría: Windows, PowerShell 7, Python 3.14.2;
- distribución importada: `google-cloud-documentai 3.15.0`;
- wheel directo bloqueado: `google_cloud_documentai-3.15.0-py3-none-any.whl`;
- SHA-256 oficial verificado por el adquiridor: `f040f4f9db43411184197a808b11fde52b580723a9fca336f43ea7c6a885bfd8`;
- firma observada: `DocumentProcessorServiceClient.process_document(request: ProcessRequest | dict | None, ...) -> ProcessResponse`;
- serialización observada: `Document.to_dict` y `Document.to_json` de proto-plus.

## Reconstrucción y pruebas

El materializador reconstruyó 5/5 archivos en un directorio nuevo y validó SHA-256 por bloque. Desde `google_document_runtime/`:

```text
test_fail_closed_and_atomic ... ok
test_official_sdk_contract_and_evidence ... ok
test_profile_keeps_unproven_classes_blocked ... ok
Ran 3 tests in 0.013s
OK
GOOGLE_DOCUMENT_AI_SDK_CONTRACT_PASS version=3.15.0
```

Se verificaron `RawDocument` con bytes/MIME, `ProcessRequest` con recurso de processor version exacto, preservación completa de `ProcessResponse`, hashes input/response, receipt, commit atómico, rechazo de response vacío, output existente, formato no admitido y processor default mutable. Todas las clases no probadas siguen `BLOCKED_*` y `automatic_storage=false`.

## Fallo retenido

La primera ejecución del test después de materializar se lanzó desde la raíz Elite y falló con `ModuleNotFoundError: test_document_ai_runtime`. No hubo fallo del código materializado. Se registró `LIB-FAIL-040`, se fijó el working directory exacto y la repetición fue verde.

## Condiciones abiertas

No existían cuenta/proyecto GCP autorizado, billing approval, ADC, IAM, processor desplegado, corpus representativo ni ground truth. Faltan llamadas sandbox, evaluación por clase/campo/ambigüedad, revisión, residencia/privacidad, costo/cuota, carga/fallos y rollback real. Por eso la admisión es `CONDITIONED`, no `REUSABLE_PACK`.

## Fuentes oficiales

- Google Cloud Document AI, pretrained overview e Invoice Parser: https://cloud.google.com/document-ai/docs/pretrained-overview
- Google Cloud Document AI overview, Custom Extractor/Classifier/Splitter: https://cloud.google.com/document-ai/docs/overview
- SDK Python 3.15.0, `DocumentProcessorServiceClient`: https://docs.cloud.google.com/python/docs/reference/documentai/3.15.0/google.cloud.documentai_v1.services.document_processor_service.DocumentProcessorServiceClient
- API resource format and processor types: https://docs.cloud.google.com/document-ai/docs/reference/rest/v1/projects.locations.processorTypes
