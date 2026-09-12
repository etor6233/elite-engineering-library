# Google Cloud Document AI Official Process Sample — Evidence V1

Date: 2026-08-28  
Scope: código/test/licencias oficiales Google, invoice del test, packing-list input/output y reconstrucción offline.  
No cubre: llamada live, cuenta GCP, exactitud de negocio, persistencia productiva ni otras clases documentales.

## Authorities

1. `GoogleCloudPlatform/python-docs-samples` commit `841379c28828404c6d3944e3da9a8a0f46b4cd0e`, 2026-08-26T15:27:46Z, firma GitHub `verified=true/reason=valid`.
2. `GoogleCloudPlatform/document-ai-samples` commit `001ba391ab4a2f40d001cc0387618cb3c3699523`, 2026-01-05T16:01:34Z, firma GitHub `verified=true/reason=valid`.

Ambos son Google Cloud Platform y Apache-2.0. El segundo repositorio declara que sus samples son demostrativos y no un producto oficialmente soportado; esa condición se conserva.

## Exact embedded upstream files

| Archivo | Bytes | SHA-256 | Provenance |
|---|---:|---|---|
| `process_document_sample.py` | 3.596 | `7e384dc2c38ebdcbc2c268a63bb9be9794576c1408034e84d53fd8ac5dc504ee` | `VERBATIM`, Google Apache-2.0 |
| `process_document_sample_test.py` | 1.697 | `5b39bbd2d309efcc4791bd739fa38961639fa33d7a4e58b739050d3cfd16fcd7` | `VERBATIM`, Google Apache-2.0 |
| python-docs-samples `LICENSE` | 11.357 | `c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4` | `VERBATIM` |

## Exact acquired fixtures

| Artefacto | Bytes | SHA-256 |
|---|---:|---|
| Google official test `invoice.pdf` | 58.980 | `8dd79a9bfc49c8bc79180686b19a223a52177221460f5279e1c2edcada214ce5` |
| packing-list PNG | 69.350 | `b45f86a82194a6210ec3f2b149a9c486934b01f6087cf7e6f859f32c5f709e7b` |
| cached Document AI JSON | 361.762 | `509cdc9e763f4c37253c617fba800452473c7a7e2bc5545ad6f6c5db956b91fb` |
| document-ai-samples `LICENSE` | 11.358 | `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30` |

Adquisición real offline desde cache exacta: lock `sources=2 fixtures=4 selected=4`, receipt y 5 archivos finales contando el receipt. El validador reconstruido devolvió:

```text
GOOGLE_DOCAI_PACKING_LIST_OUTPUT_VALIDATION_PASS pages=1 entities=7 barcodes=1 storage_authorized=false
```

El JSON oficial contiene barcode CODE_39 `E-58702865`; entity types invoice, line item, purchase order, ship-to y supplier. Confidencias publicadas: purchase order `0.5372485`, ship-to name `0.55553776` y supplier name `0.30161396`. Estas cifras impiden presentar el output como extracción perfecta o autorizar almacenamiento automático.

## Rebuild/tests

El pack 0.1.0 materializó diez archivos desde Markdown. Resultados:

```text
python test_official_process_document_sample.py
Ran 3 tests ... OK

GOOGLE_DOCAI_OFFICIAL_FIXTURE_TEST_PASS positives=2 negatives=4
py_compile con `cfile` temporal corto ... PASS
```

La regresión ejecuta el sample Google real con dobles sólo del SDK y prueba endpoint regional, bytes/MIME, processor y processor-version paths, field mask y selector de página. Los negativos prueban ID desconocido, destino ocupado, cache alterada y licencia no aceptada.

La primera integración global reveló una ruta de exactamente 260 caracteres: PowerShell veía el archivo, pero Python 3.14 no podía abrirlo. El layout materializado se acortó de `upstream/google_python_docs_samples/...` a `upstream/google/...` sin modificar los bytes Google. Una reconstrucción nueva bajo un nombre temporal largo volvió a pasar los tres tests y 2 positivos/4 negativos. `py_compile` también escribe el `.pyc` en un temporal corto explícito para no reintroducir el límite mediante `__pycache__`.

## Admission

`SUPPORTED_REFERENCE / REBUILD_VERIFIED / CONDITIONED`. El código está listo para materializar y conectar a un processor autorizado; no está listo para persistir automáticamente. GCP/ADC/processor/version/costo, corpus, ground truth, schema, campo requerido, revisión, drift, reconciliación y retención siguen siendo condiciones del proyecto. El requirements histórico de los samples no gobierna: el runtime seleccionado es el artifact oficial Google Document AI 3.15.0 con lock transitivo del target.
