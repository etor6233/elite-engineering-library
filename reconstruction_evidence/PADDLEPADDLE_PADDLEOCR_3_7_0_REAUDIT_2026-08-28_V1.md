# PaddlePaddle PaddleOCR 3.7.0 re-audit — 2026-08-28 V1

## Resultado gobernante

`PaddlePaddle/PaddleOCR` 3.7.0 y su dependencia `PaddlePaddle/PaddleX` 3.7.2 quedan admitidos como `PINNED_CANDIDATE_CONDITIONED`, no como `REUSABLE_PACK`, implementación productiva ni autoridad de almacenamiento automático.

El código real oficial, los paquetes PyPI y dos modelos PP-OCRv6 se fijaron por revisión y SHA-256. La suite oficial PaddleOCR pasó completa en el target publicado para Python distinto de 3.8: 220 tests ejecutados, 220 PASS, 5 deselected y 1 warning. Una inferencia CPU local con los dos directorios de modelo auditados produjo JSON real. Nada de eso demuestra exactitud sobre invoices, proformas, packing lists u otros documentos del negocio: el fixture upstream no publica ground truth y la confianza del modelo no equivale a corrección.

## Identidad de código

| Componente | Identidad inmutable | Distribución | Licencia |
|---|---|---|---|
| PaddleOCR | release `v3.7.0`; commit firmado `b03f46425e8ff4442b268ce449e3eef758146cd4`; tree `ece94ffce26900737058146515e63bbd3d79f374` | ZIP 118.686.571 bytes; 2.488 files + 513 directory entries; SHA-256 `449760f5d29842ab6f9cb02cbfbc1f343a5b40fbfd7105333835d72d6a5c741e` | Apache-2.0; `LICENSE` 11.376 bytes/SHA-256 `3840c5c0c61c294264d2dd77b8777be6ddd90121ef4e0e64abcd22edea581d6e` |
| PaddleX | release `v3.7.2`; commit unsigned/unverified `ffb64904d23708863ff5b8da312a5cbd52a7f462`; tree `12ae9b4e36f2b388a12115b7cfb688487fa9f7fb` | ZIP 7.321.711 bytes; 3.038 files + 757 directory entries; SHA-256 `59a787a9abe2ed43d96cfefc9bb1179a8f221cbe7bea19b23594501c1e4492b5` | Apache-2.0; `LICENSE` 11.325 bytes/SHA-256 `accc36c817ac5ede473d05732e14afc11cb6e55ef5919bc826b653f631165043` |

Fuentes oficiales: [PaddleOCR commit](https://github.com/PaddlePaddle/PaddleOCR/tree/b03f46425e8ff4442b268ce449e3eef758146cd4), [PaddleX commit](https://github.com/PaddlePaddle/PaddleX/tree/ffb64904d23708863ff5b8da312a5cbd52a7f462).

El tag PaddleX no se presenta como firmado. Además, el workflow PaddleOCR instala PaddleX desde la rama móvil `release/3.7`. La auditoría no reprodujo esa fuente móvil: resolvió el paquete estable exacto 3.7.2 y conservó la diferencia como condición.

## Paquetes ejecutados

| Artefacto oficial | Bytes | SHA-256 |
|---|---:|---|
| `paddleocr-3.7.0-py3-none-any.whl` | 146.750 | `c0f0a81ad4112727f30c6fcf986ac0ef6a120d31ee0991a01fae0357ee32d338` |
| `paddleocr-3.7.0.tar.gz` | 3.082.862 | `9b81eb62cf37e590d4873e60262ff56f968ff9de6ce1f565749760ceb1c422e2` |
| `paddlex-3.7.2-py3-none-any.whl` | 2.239.708 | `f1678bf650bbaccfd8f0d4e49d0ae631b4685c829fdae6e802ccd90d4fcb9a7f` |
| `paddlex-3.7.2.tar.gz` | 4.415.130 | `6a60b4595e2d51460f7f51326672adca86b89ff9ecfd02ac348fb69739e0093c` |
| `paddlepaddle-3.1.0-cp312-cp312-win_amd64.whl` | 99.784.181 | `3cb6d98eece900e34c05fa0428ccc32836525e72af25cc8ad10a48d4046c4639` |

Los metadatos y URLs se consultaron en la API pública de [PyPI PaddleOCR 3.7.0](https://pypi.org/project/paddleocr/3.7.0/), [PyPI PaddleX 3.7.2](https://pypi.org/project/paddlex/3.7.2/) y [PyPI PaddlePaddle 3.1.0](https://pypi.org/project/paddlepaddle/3.1.0/).

### Runtime mínimo document parser

- input: `paddlepaddle==3.1.0`, `paddleocr[doc-parser]==3.7.0`, pytest y pytest-cov;
- lock: 104 paquetes, 239.480 bytes, SHA-256 `37671931bc10d65fd42974bbb161d6076b53797d3565434c11c842b458e5c599`;
- sync con `--require-hashes`: PASS;
- imports: PaddlePaddle 3.1.0, PaddleOCR 3.7.0 y PaddleX 3.7.2 PASS;
- `uv pip check`: 105 paquetes instalados compatibles;
- `pip-audit 2.10.1 --require-hashes --disable-pip`: 104 dependencias auditadas, 0 paquetes afectados, 0 vulnerabilidades conocidas; JSON 5.892 bytes/SHA-256 `41253fe57d5369eabf45a11f10fe660f4ddd24cd725db0dd7580e5b295a21f44`.

### Entorno oficial de tests

El workflow oficial instala `requirements.txt`, pytest, `.[all]` y PaddleX. El lock reproducible Windows/Python 3.12 usó las mismas superficies, sustituyendo sólo la rama móvil PaddleX por la versión PyPI exacta 3.7.2:

- lock: 140 paquetes, 335.242 bytes, SHA-256 `8c16e681a289ab72221a18d113166229282a08f539047bf9af558eb608cfebd6`;
- sync hash-locked, `uv pip check` e imports de Paddle/PaddleOCR/PaddleX/scikit-image: PASS;
- comando oficial: `pytest --verbose tests/`;
- resultado: 220 PASS, 5 deselected, 1 warning, 824,45 segundos;
- `pip-audit 2.10.1`: 7 registros conocidos en 4 paquetes opcionales del grafo `all`/test: `langchain`, `langchain-core`, `langchain-openai` y `langchain-text-splitters`;
- IDs únicos conservados: `PYSEC-2026-2192`, `PYSEC-2026-2193`, `PYSEC-2026-2562`, `PYSEC-2026-76`, `PYSEC-2026-77`;
- JSON de auditoría: 19.717 bytes/SHA-256 `f685e14799b18cec5f745bc7c302ef064b667cf885cd693eda899f5b672853c0`.

La separación es obligatoria: “runtime mínimo sin findings conocidos” no borra los findings del extra `all` ni autoriza que un proyecto los instale.

## Modelos exactos ejecutados

Las páginas oficiales de Hugging Face declaran ambos repositorios Apache-2.0. Se siguieron revisiones verificadas, nunca el alias móvil `main`.

### PP-OCRv6 medium detection

- repositorio: [PaddlePaddle/PP-OCRv6_medium_det](https://huggingface.co/PaddlePaddle/PP-OCRv6_medium_det/tree/8e0f56fb2ef86b461d99cfc7ac5c137738985f61);
- revisión verificada: `8e0f56fb2ef86b461d99cfc7ac5c137738985f61`;
- `README.md`: 23.247 bytes/SHA-256 `59f37b08a9a49aa6c1275e7b0363991353582f4dc7676539a5d2b05aab57e744`;
- `inference.json`: 312.150 bytes/SHA-256 `0f1a7ec35da36173529c7a60238b7f7919e3831929c3f700ad90ad4896adecd5`;
- `inference.pdiparams`: 61.960.476 bytes/SHA-256 `85218d2e3d98f5a21c58b4220627be923a97aee5db3cc71f39536ab31ac53960`;
- `inference.yml`: 886 bytes/SHA-256 `7298d5ead546584af2504d03355f881ac7a7bc0eb1e282d3e159277c1d0af871`.

### PP-OCRv6 medium recognition

- repositorio: [PaddlePaddle/PP-OCRv6_medium_rec](https://huggingface.co/PaddlePaddle/PP-OCRv6_medium_rec/tree/e5a92bcbc5cc1b494628e458d267778f0704fd7c);
- revisión verificada: `e5a92bcbc5cc1b494628e458d267778f0704fd7c`;
- `README.md`: 23.474 bytes/SHA-256 `75ff9c4853ed171f36224127805062f9bcd1cfbbcacae3d577344866df74f6d3`;
- `inference.json`: 221.814 bytes/SHA-256 `0b2e25e990bd072f1bf77d59d67d508bce6c4bd44af6624e0fb27d6da2cd00e8`;
- `inference.pdiparams`: 76.465.087 bytes/SHA-256 `1b01c79a914587933f615569e75de54f2e638ebb5d3f3b3c1b38c24ede8c7319`;
- `inference.yml`: 150.580 bytes/SHA-256 `991b700facf5b50a7de193468207d5f4255b538dde0d312ae3b7c7a9b6873129`.

Total ejecutado: 8 files, 139.157.714 bytes. Los modelos se pasaron por `text_detection_model_dir` y `text_recognition_model_dir`; se deshabilitaron orientación de documento, unwarping y orientación de línea. Por lo tanto no hubo descarga de modelos implícita.

## Inferencia local observada

- fixture upstream: `tests/test_files/table.jpg`, 16.208 bytes, SHA-256 `acd113bb3a89b488941ee0962776a28e45897fa2802cd306ec3bb68d9043115c`;
- inicialización CPU: 1,565 s;
- predicción CPU: 1,223 s;
- resultados: 1 documento, 12 textos y 12 scores;
- score mínimo: 0,991830; score medio: 0,999201;
- JSON: 10.198 bytes/SHA-256 `7eddd0f2241d3b48bc979250268737de88f959051f27cfa7cd85eb30d55674f1`.

El repositorio usa ese archivo como fixture de smoke para OCR y tablas, pero no entrega una anotación ground-truth junto a él. No se calculó exactitud, precision, recall, F1 ni error por carácter. Los scores no autorizan persistencia.

## Artefactos focales de código

| Path PaddleOCR | SHA-256 |
|---|---|
| `pyproject.toml` | `63e8bfa19e197e47649a95dad422bcb09c1d2959367cafd40b0a7f19e4452412` |
| `requirements.txt` | `d33da5c48e908fc328b48569098e52c4a9df4529a177fe09920c1ccb4531a9ef` |
| `.github/workflows/tests.yml` | `e45309175bf7a94f369c4cfadc8eab5113c67eb3b6567e53e93930aed4b7b69c` |
| `paddleocr/_pipelines/ocr.py` | `fc870f7a7f1b33756020d75c63c6dc0dc04981df1ff46627a040336ead9a9b88` |
| `paddleocr/_pipelines/pp_structurev3.py` | `c4f2bf78774489843da5a9632507c4453b9531275cb851da6e937cce55e6ccb9` |
| `tests/pipelines/test_ocr.py` | `936b9f99f88af16d17d4aed53e29c3c1992d9df3d00cc1fa8f292d39f8dac0a9` |
| `tests/pipelines/test_pp_structurev3.py` | `628dd6290f23edfb0db5dc523f96495a1231153d65feb6769584c4a02f91d1cc` |
| `tests/security/test_latexocr_pickle.py` | `c5f32872c6cf4f90ba130a570c0ffd5a17ccc9915df979629ab6cd53969f0eb2` |
| `tests/security/test_lmdb_pickle.py` | `27ea359401aeff402bd422b5d1f22b0e248e1507568ec12a2641b48ed674e574` |
| `docs/version3.x/pipeline_usage/OCR.en.md` | `57b3df0656e1829c747d180bbb43934551f71bb9e19e80aa17c8477281fb7881` |
| `docs/version3.x/pipeline_usage/PP-StructureV3.en.md` | `8e7e1a7dfdfa1f0f4c13bbe50c495140e04b827e6a07686b787b238e3862edae` |

PaddleX conserva además hashes de `setup.py`, su pipeline/result OCR y la arquitectura/configuración PP-OCRv6 en el source lock canónico.

## Decisión y condiciones de uso

El agente puede adquirir estos sources como candidato local oficial, estudiar y reutilizar el código bajo Apache-2.0 y reproducir la inferencia sólo después de verificar paquetes, modelos y locks exactos. No puede:

- usar `main`, `release/3.7`, nombres de modelo móviles o descargas implícitas;
- instalar el extra `all` ignorando sus findings;
- equiparar score/confidence con exactitud;
- afirmar que invoices, proformas, packing lists u otra clase están soportadas por este smoke;
- almacenar automáticamente campos de negocio sin corpus anotado, schema cerrado, reglas cruzadas, umbrales, rechazo/revisión, evidencia y autorización explícita;
- presentar PaddleOCR/PaddleX como sistema empresarial completo.

La promoción requiere un proyecto concreto, documentos autorizados y anotados por clase/variante/idioma/proveedor, evaluación de campo y documento, security/load/soak/recovery, modelo/runtime reauditable, proceso de revisión y decisión del propietario sobre persistencia.

## Integración canónica

- `OFFICIAL-UPSTREAM-ACQUISITION-CORE` 0.4.53: 100 sources;
- `document-intelligence-leaders`: 45 sources;
- tests del adquiridor: 9 negativos PASS;
- perfiles: 14 válidos, 5 negativos y 1 positivo PASS;
- failure learning: `LIB-FAIL-676` a `LIB-FAIL-690`, sin borrar intentos fallidos ni atribuirlos al upstream cuando correspondían al proceso de auditoría.

