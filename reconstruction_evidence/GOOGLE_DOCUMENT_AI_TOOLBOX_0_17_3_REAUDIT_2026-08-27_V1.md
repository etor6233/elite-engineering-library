# Google Document AI Toolbox 0.17.3 — reauditoría V1

Fecha de corte: 2026-08-27. Clasificación: `PINNED_CANDIDATE_CONDITIONED`. No es `REUSABLE_PACK`, no se compone automáticamente y no autoriza persistencia productiva.

## 1. Autoridad e identidad

- propietario: Google LLC / `googleapis`;
- release oficial: [`google-cloud-documentai-toolbox-v0.17.3`](https://github.com/googleapis/google-cloud-python/releases/tag/google-cloud-documentai-toolbox-v0.17.3);
- commit directo: [`3e29682282ac03d929fbfb93b5f94736fd420249`](https://github.com/googleapis/google-cloud-python/commit/3e29682282ac03d929fbfb93b5f94736fd420249), 2026-08-24T21:08:08Z, firma presente y verificación GitHub `valid`;
- source archive oficial: 209,622,796 bytes, SHA-256 `2e1a72b40df935e0fb2055a42464eb89a9b289270d82926e749424a160c81bb7`;
- licencia raíz y package: Apache-2.0, SHA-256 `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30`;
- PyPI oficial: [`google-cloud-documentai-toolbox 0.17.3`](https://pypi.org/project/google-cloud-documentai-toolbox/0.17.3/), publicado por Trusted Publishing de Google;
- wheel: 43,670 bytes, SHA-256 `46f37a42ab98a1255ea078fe4ce85084cc7ecf127b6f5ccc3c546599262da014`;
- sdist: 24,314,548 bytes, SHA-256 `4e19e04abc9fcb4d3b131697b91c3486d03c38cbf01d06581f2bda43050f0f46`.

El repositorio standalone anterior está archivado y Google indica que el source se movió al monorepo `google-cloud-python`; no se usó una copia comunitaria.

## 2. Fidelidad source ↔ sdist ↔ wheel

El subtree oficial contiene 128 entradas: 97 archivos regulares, 29 directorios, un symlink POSIX y la entrada raíz. `docs/CHANGELOG.md` es symlink modo `0120777` con target `../CHANGELOG.md`; Windows no puede representarlo mediante `tar.exe`, por lo que no se fingió una extracción fiel. La proyección auditada extrajo sólo archivos regulares con validación de traversal y declaró el enlace por separado.

Comparación commit/sdist:

- 72 paths comunes;
- 72/72 SHA-256 idénticos;
- cero mismatch;
- 25 archivos sólo-repo: configuración, documentación, `noxfile.py`, `pytest.ini` y constraints;
- ocho archivos sólo-sdist: metadata `PKG-INFO`/egg-info y `setup.cfg` generado.

El wheel y el sdist se fijaron como artefactos externos; no se redistribuyen dentro de Elite.

## 3. Superficie oficial reutilizable

El código de Google implementa wrappers y utilidades para:

- cargar resultados Document AI locales/GCS;
- combinar shards y documentos;
- páginas, tokens, líneas, párrafos, form fields, tablas y entidades;
- splitter/classifier outputs y separación de PDF;
- imágenes de entidades;
- hOCR con template autoescaped;
- conversiones y utilidades GCS/BigQuery.

Los fixtures oficiales incluyen invoice, procurement multi-document, packing/splitter, classifier, form parser, layout parser, tablas, imágenes y documentos vacíos. Esto prueba capacidades de postprocesado del output Document AI; no prueba OCR/model inference ni exactitud de datos del negocio.

## 4. Ejecución reproducida

Entornos aislados, sin credenciales ni llamadas de proveedor:

- Windows CPython 3.12.13, wheel exacto y `testing/constraints-3.12.txt` SHA-256 `712160110a975b0eba66e0afe62bd0cc6b821556d918e0dd51b73811b337110d`;
- Linux/WSL Python 3.14.4, uv 0.11.16 portable: archive SHA-256 `74947fe2c03315cf07e82ab3acc703eddef01aba4d5232a98e4c6825ec116131`, checksum manifest SHA-256 `8ef7fe76d67be3330e18e8d6ecbbb68f7a1ae46fe31198008170e911ad025c6a`;
- `pip check`: PASS en ambos entornos.

Resultados oficiales:

- Linux/protobuf `python`: 164/164 PASS; JUnit SHA-256 `bd53cfb59005dfdd69d54fa9e12ce909a97cb5ef595058d7590fd84fb2ee8418`;
- Linux/protobuf `upb`: 164/164 PASS; JUnit SHA-256 `44d607648cbe34dfc120b9d2ae4f9106671bb07b684ea2d41bba749a67051dc9`;
- Windows/protobuf `python`: 160 PASS y cuatro FAIL exclusivamente porque pandas emite CRLF y los tests comparan LF literal; JUnit SHA-256 `642628dde5073b5641477c4fe7553a7da406829e0d957196d7c82910f98877f1`;
- regresiones focales Windows en ambas implementaciones: traversal de `split_pdf`, escape de contenido hOCR y escape de title, 6/6 PASS;
- probe combinado de los cinco wheels documentales exactos: `DOCUMENT_SDK_PROBE_PASS`, versiones e interfaces Microsoft/Google/OpenAI verificadas; `pip check` PASS.

No se modificó una línea upstream ni se normalizaron tests para ocultar el fallo Windows.

## 5. Dependencias y seguridad

- freeze Windows: 78 líneas, SHA-256 `adfc8369c026b14751938a88a6d2ae61998eec3bbc29b7528a8476c1bed1db0a`;
- freeze Linux: 77 distribuciones, SHA-256 `07e0e2037e26846cbc682fec2b3ffaa75fdbb62fb0a119daa1060521bc5d7a78`;
- Google OSV-Scanner 2.5.1 exacto, commit `c84fa4568f2526d0333e9a914ea8a0a5f74ad68b`, binario SHA-256 `25e42f5ef6711fd8c0fb45390972205891dd44c6bd02ac93f0f63e8e98d9bfb6`;
- scan Windows: 77 packages detectados, reporte SHA-256 `9d9bcf7478249810518d0b83be2d42f15d9ac0d6bbf9c59833352453c1314b1c`;
- scan Linux: 76 packages consultables de 77 distribuciones, reporte SHA-256 `15c7c7ebeb86512c8bdd4f9cb8d8b51d94b9ce23cf97f6e2c9f90b9603ee5219`.

Ambos scans detectaron un único advisory bajo aliases `CVE-2026-25087`, `GHSA-rgxp-2hwp-jwgg` y `PYSEC-2026-113` para `pyarrow==22.0.0`. El [`advisory revisado`](https://github.com/advisories/GHSA-rgxp-2hwp-jwgg) y la [release oficial Apache Arrow 23.0.1](https://arrow.apache.org/blog/2026/02/16/23.0.1-release/) fijan la corrección en 23.0.1. Apache aclara que el camino vulnerable es el pre-buffer del lector IPC C++ y no está expuesto por los bindings Python. Sin embargo, Toolbox 0.17.3 declara `pyarrow>=15.0.0,<23.0.0`; no existe una versión corregida que satisfaga su metadata. No se forzó 23.0.1, no se silenció OSV y la condición permanece `UPSTREAM_OPEN`.

## 6. Por qué no está listo

Google describe Toolbox como experimental/work-in-progress, con cambios incompatibles probables; PyPI lo clasifica Alpha y la documentación de Document AI lo presenta como Pre-GA. Además faltan para un proyecto real:

- cuenta, región, processor/model versions, IAM, cuotas y costos;
- corpus representativo por clase documental/campo/idioma/calidad;
- ground truth, evidencia por campo, umbrales y revisión;
- reconciliación, retención, privacidad, observabilidad y manejo de fallos;
- carga, concurrencia, recovery, canary y rollback;
- release oficial que acepte PyArrow corregido o excepción de reachability formal, acotada y con vencimiento.

Por eso los artefactos quedan `PINNED_CANDIDATE`, la fuente `PINNED_CANDIDATE_CONDITIONED` y ningún perfil actual la selecciona.

## 7. Integración Elite V52

- `OFFICIAL-UPSTREAM-ACQUISITION-CORE` 0.4.44: 90 sources, 23 archivos, 14 perfiles, cinco negativos de perfiles, un positivo y nueve negativos de lock/acquirer;
- source acquisition amplía symlink handling de forma fail-closed para `directory|file`, mantiene restricción no-Windows y rechaza cualquier kind diferente;
- `OFFICIAL-DOCUMENT-SDK-ARTIFACT-CORE` 0.1.1: siete artefactos, cuatro negativos y una adquisición offline hash-exacta;
- wheel/sdist Toolbox validan individualmente sin descarga en `ValidateOnly`;
- el probe incluye `DocumentAIToolboxDocument` y las cinco versiones exactas.

Este expediente demuestra identidad, licencia, fidelidad, imports y tests locales. No demuestra exactitud documental productiva ni autoriza uso automático.
