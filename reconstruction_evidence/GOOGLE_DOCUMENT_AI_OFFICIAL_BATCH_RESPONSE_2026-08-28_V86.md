# Google Document AI Official Batch and Response — V86

Fecha: 2026-08-28  
Alcance: código oficial Google para batch processing y lectura de respuestas Document AI, artefactos runtime necesarios y composición reproducible.

## Fuentes oficiales fijadas

- repositorio: [`GoogleCloudPlatform/python-docs-samples`](https://github.com/GoogleCloudPlatform/python-docs-samples/tree/dc0eecc2187791fea70f48201241288e25fc9620/documentai/snippets), commit firmado/verificado `dc0eecc2187791fea70f48201241288e25fc9620`, Apache-2.0;
- documentación batch: [`Send a batch process documents request`](https://docs.cloud.google.com/document-ai/docs/samples/documentai-batch-process-document);
- API: [`processors.batchProcess`](https://docs.cloud.google.com/document-ai/docs/reference/rest/v1/projects.locations.processors/batchProcess);
- wheel/sdist Google Cloud Storage 3.13.1: [`PyPI`](https://pypi.org/project/google-cloud-storage/3.13.1/), Apache-2.0.

## Bytes Google incorporados sin modificación

| Archivo | Bytes | Git blob | SHA-256 |
|---|---:|---|---|
| `batch_process_documents_sample.py` | 7.031 | `6d15d8713a7aa96d6012f8b36e2fea2333ac94d0` | `3ccd1f828172c76f9268451091a90225dbd173cd7979755d5a66adebcea41759` |
| `batch_process_documents_sample_test.py` | 2.568 | `168fe3ac3ee003580a63d2eb4c7aa1eb26eb6253` | `674a525239061de0a8d4c7075e846f69d91713aff27940f70d3f68ca3a19f6ec` |
| `handle_response_sample.py` | 20.533 | `58bbb1debe09c6fad3814789f3db78af4679e061` | `8fc1adaceba3ac8ad4d871f8003b2ce43fbadfd9cc7b80599cefcc6edc8bfb91` |
| `handle_response_sample_test.py` | 7.859 | `b7c65834ccafbcb8fbae6f46807d0a7184ff1048` | `a60d4d9a53aa3c6383d60bb2cfc22489fde622d09ac994d68161ee771b1fea45` |

El código batch admite documento GCS individual o prefijo, construye `BatchProcessRequest`, espera la operación, exige estado `SUCCEEDED`, recorre `individual_process_statuses`, lista sólo JSON y reconstruye `documentai.Document`. El response sample contiene funciones Google para OCR/calidad/layout, forms/tables, entidades especializadas, splitter, layout chunks, Custom Extractor, normalización y text anchors. No se reescribió lógica ni se atribuyó a Google el lock/verifier Elite.

## Artefactos runtime

- `OFFICIAL-DOCUMENT-SDK-ARTIFACT-CORE` 0.3.1 fija 14 artefactos, incluidos wheel y sdist Storage 3.13.1: wheel 341.486 bytes/SHA-256 `98208de6c21e85cecd3eb44551894efff33d98365500e178867d4305854a770a`; sdist 17.341.051 bytes/SHA-256 `a80bf8cac2794808aa61c50c5f769ecbbe2d10331bacd0d69d30e59b14b346b2`;
- `GOOGLE-DOCUMENT-AI-RUNTIME` 0.1.3 fija un lock CPython 3.12/Windows x86-64 de 23 wheels para Document AI 3.15.0 + Storage 3.13.1;
- instalación `--require-hashes`, `pip check` e imports de ambos clientes: PASS;
- consulta oficial Google OSV por los 23 nombres/versiones: cero matches conocidos el 2026-08-28. Es un resultado fechado, no garantía futura.

## Reconstrucción y pruebas

- lifecycle 0.2.0 materializó 21/21 archivos;
- cuatro round-trips Google coincidieron por SHA-256 con la adquisición temporal;
- verifier estático: `upstream_files=18 python_files=17 functions=10 static_only=1` PASS;
- verifier completo con Document AI 3.15.0 + Storage 3.13.1: PASS;
- cinco tests lifecycle Google offline-safe: `5 passed`;
- `VERIFY_LIBRARY_PASS`: 60 packs, 572 archivos, 398 Markdown antes de esta evidencia; lifecycle profile 3/40 y Google runtime 9/102;
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: 60 packs, 111 upstreams, 14 SDK artifacts y un lifecycle Google materializado/estáticamente verificado.

## Límite honesto

Los módulos set-default, batch y response de Google conservan tests live que requieren GCP, ADC, processor, Cloud Storage, IAM, cuota y costo. No fueron sustituidos por mocks propios para fingir precisión. El sample oficial no demuestra completitud de todos los shards, exactitud por campo, reconciliación empresarial ni autorización de almacenamiento para un corpus desconocido. `UP-FAIL-174` obliga a conservar esas condiciones. El código ya es inmediatamente materializable; sus efectos externos y la persistencia permanecen bloqueados hasta la evidencia del proyecto.
