# Google Document AI classify/split re-audit — 2026-08-28 V1

## 1. Alcance y autoridad

Auditoría read-only del repositorio oficial `GoogleCloudPlatform/document-ai-samples` fijado en el commit `001ba391ab4a2f40d001cc0387618cb3c3699523` y del archive exacto:

- URL: `https://github.com/GoogleCloudPlatform/document-ai-samples/archive/001ba391ab4a2f40d001cc0387618cb3c3699523.zip`;
- bytes: `155461201`;
- SHA-256: `3131ff0604685967932cad1b689a64f7e50af1dc6b45524d99c5389342d34e8f`;
- licencia raíz: Apache-2.0, SHA-256 `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30`.

Se revisaron por separado `classify-split-extract-workflow/` y `pdf-splitter-python/`. La documentación oficial vigente consultada fue:

- `https://docs.cloud.google.com/document-ai/docs/splitters`;
- `https://docs.cloud.google.com/document-ai/docs/custom-splitter`;
- `https://docs.cloud.google.com/python/docs/reference/documentai-toolbox/latest/google.cloud.documentai_toolbox.wrappers.document.Document`;
- `https://github.com/googleapis/google-cloud-python/tree/main/packages/google-cloud-documentai-toolbox`.

## 2. Workflow oficial completo

`classify-split-extract-workflow/` contiene 43 archivos, nueve Python, Cloud Run Job, Workflows, GCS, Document AI y BigQuery. Sus nueve archivos Python compilaron con `python -m compileall` sin modificar source.

Archivos focales exactos:

| Path | Bytes | SHA-256 |
|---|---:|---|
| `README.md` | 20.063 | `9d2bbc133a7aafaed81de23d2546fe97965f97f00ee6b13d9d5dfda43e16e051` |
| `classify-job/main.py` | 2.675 | `06c8c070b3de2ea9989338e2240cf5a3f610fad8939ab764674d7739ecc443ec` |
| `classify-job/split_and_classify.py` | 14.582 | `4d737a2b518f900d9a614dedbba9c0eb5c674cd5f26689860ee1661d0f2fad2e` |
| `classify-job/main_test.py` | 1.200 | `9e54fd7a8d487225d22e4ed9beee6d84f5262c1b0137d2ce2af30396d8f3b23a` |
| `classify-job/requirements.txt` | 215 | `11dea956092d71cfa442d2948ef073f158f9345084f75a613b427f77cb094ac7` |
| `classify-job/Dockerfile` | 307 | `b9ada06bbe08353e5a01816e7ac7be1b40ab3cecf7c5b217771de25f39708297` |
| `setup.sh` | 12.621 | `94754f3ac8d06d4772fe245636b88e8f7ff01495cfb6e2b7e67858c3d7416860` |
| `deploy.sh` | 2.668 | `91fb056337bcc27167b6043c5cc9bb30bb1da678f06b08007e303c0450e8ba78` |
| `classify-extract.yaml` | 10.509 | `8aeac06af37c140cd1b78060d6bc269f35319d3d4e8221220c55d473dd1f59cc` |

La única ruta con nombre de test, `classify-job/main_test.py`, es un helper manual que muta `config.INPUT_FILE` y llama al servicio; no contiene casos/aserciones automatizados. No existe una suite focal que demuestre clasificación, separación, persistencia, reintentos, concurrencia o IAM.

El propio README conserva como trabajo futuro: revisión humana por confianza de clasificación y extracción, control de cuotas, soporte de más de 15 páginas, otros MIME types, UI/HITL, logging correcto y mejoras de splitting. Además:

- `config/config.json` usa umbral `0.50` y `classification_default_class=generic_form`;
- baja confianza se enruta a la clase por defecto, no a rechazo/revisión;
- `main.py` captura `Exception`, registra el texto y continúa al callback; Workflows contiene dos `except:` sin contrato tipado;
- `setup.sh` concede `roles/storage.admin`, `roles/bigquery.admin` y `roles/run.admin` a nivel de proyecto o bucket;
- `deploy.sh` elimina y recrea el job en cada despliegue;
- el Dockerfile parte de `python:3.10-buster`, instala desde requirements sin hashes y su `CMD` apunta a `classify.py`, archivo que no existe en el directorio;
- no hay intake de tipo/malware/cuarentena, persistencia idempotente, CAS/reconciliación, expediente de evidencia por campo, retención inmutable ni rollback demostrado.

Conclusión: `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`. El diseño oficial informa la topología, pero no se copia, despliega ni presenta como pipeline empresarial listo.

## 3. Splitter histórico independiente

`pdf-splitter-python/` sí trae dos unit tests con mocks y un PDF oficial de diez páginas; el source compila. Identidad focal:

| Path | Bytes | SHA-256 |
|---|---:|---|
| `README.md` | 3.812 | `f70c30ea0d6bc526e9e5a9dd0afc10a174b4927ebbd664ce1326a535508a4c47` |
| `main.py` | 6.949 | `2e5a87d271b46072b12efbe046a6393fac2508839f9ba0c12fb47f749caeb96e` |
| `main_test.py` | 8.469 | `08c339b02af7044ef522f31f0509f4e8aa39e0d0f68db83ec4adb4d2315bd800` |
| `requirements.txt` | 46 | `e66374c66dda97533e4c1d35bcfc3f7da4e826d6d788f1b4363fd27c788bd2bd` |
| `multi_document.pdf` | 3.161.090 | `c850758d883db365eb7303359b7b2ac39e36c842869673655b14c356a69e766c` |

Google marca este sample explícitamente `deprecated`, fija `google-cloud-documentai==3.0.0`, declara `pikepdf>=6.2.9` sin techo y remite al Document AI Toolbox. No se ejecutó ni promovió como baseline vigente.

## 4. Implementación vigente y condición

La documentación Google actual establece que el splitter predice límites y clases pero no separa físicamente, y recomienda revisión humana porque un corte incorrecto daña dos documentos. El `split_pdf` vigente pertenece a Document AI Toolbox.

La biblioteca ya reaudita Toolbox `0.17.3` exacto en `GOOGLE_DOCUMENT_AI_TOOLBOX_0_17_3_REAUDIT_2026-08-27_V1.md`: 164/164 tests oficiales pasan en Linux, pero Google lo marca experimental/Alpha y su metadata exige `pyarrow<23`, mientras CVE-2026-25087 se corrige en 23.0.1. No se fuerza una combinación fuera del contrato upstream ni se extrae una función como fork local atribuido a Google.

## 5. Decisión

- conservar el process sample y la fixture packing-list ya admitidos sólo para smoke/contrato;
- rechazar el workflow grande y el splitter deprecated para adopción inmediata;
- mantener Toolbox como candidato condicionado, no como persistencia automática;
- no afirmar cobertura de proforma, commercial invoice, packing list o bill of lading sin schema, corpus representativo, ground truth por campo, evaluación y revisión del proyecto;
- reabrir sólo ante una release oficial corregida o un pipeline oficial nuevo que cierre tests, seguridad, evidencia y operación.

