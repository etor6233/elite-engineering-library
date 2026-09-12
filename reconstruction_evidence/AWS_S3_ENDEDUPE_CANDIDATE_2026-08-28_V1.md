# AWS S3 Endedupe Candidate — 2026-08-28 V1

## Resultado

`aws-samples/amazon-s3-endedupe` contiene una implementación oficial de conditional locking para eventos S3 duplicados, desordenados y concurrentes. Sus 17 unit tests pasan, pero la revisión exacta no es reutilizable inmediatamente. Clasificación: `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`.

No se copió código al producto ni se escribió una corrección local atribuida a AWS.

## Identidad exacta

- repo: `https://github.com/aws-samples/amazon-s3-endedupe`;
- default branch: `main`;
- head firmado/verificado: `7209a02a3cf7a613fb6d70e42caf0522b3eb0d7a`;
- fecha: `2024-09-27T14:46:39Z`;
- archive: 487.825 bytes;
- archive SHA-256: `1827d3401a01034252623e4101ac02c989c68dcb0028c46470a6ee6206af4e0a`;
- inventario: 18 archivos y 4 directorios;
- tags reales: 0; releases reales: 0;
- licencia MIT-0: 947 bytes, SHA-256 `7cb750713252efd1d578837ba8785b61319109906f0c10b87adab7cf4badfc42`.

## Archivos focales

| archivo | SHA-256 |
|---|---|
| `README.md` | `c21cb8cab2a4a30c7fa65fd916ee04064f30c845cf157d32130e5d2ad94027b3` |
| `template.yaml` | `27a39c7b0e00e02a09430bf82ef729f1ba60a66c7b10bb786a5f0b4b7c21b01f` |
| `endedupe/app.py` | `f5e6dd37fa00fb24ffc4f5743acf4bca59ff7e64fc46b6e2eb1b19a246637858` |
| `endedupe/s3index.py` | `f85fffa8d2ebcf179e7f6592ed0036a9b91f612c4a153c8d4c1091adaf6d570e` |
| `endedupe/requirements.txt` | `863a99f52a9cb5777b260e8e09d7ab3bd38681e3e8ed226e4cfcf515acaf1aea` |
| `endedupe/tests/test_app.py` | `bc26cb50b9688f73f119b237de0dc18ee89e1fd4c132bd9cf212921277b35291` |
| `endedupe/tests/test_s3index.py` | `e1a340f520979c46c73e2bc6274d5f9f99b67333d00c495b666ce1316932bb57` |
| `endedupe_integ_test/requirements.txt` | `29534986bb5b518cfa83d417bdc136f7bac1d3af52520d5724aca288b17d3559` |
| `endedupe_integ_test/test_app.py` | `68ef0f70d654e29cbafaa102f418db584fe82398bdf032a04106ecbe8ef80763` |

## Ejecución reproducida

- Python 3.12 aislado; 7 archivos compilan;
- requirements runtime exactos instalaron y `pip check` pasó;
- 17 unit tests PASS;
- cobertura medida: `s3index.py` 92%; `app.py` no fue importado. El claim README “full coverage” no quedó demostrado;
- los integration tests existen, pero exigen desplegar recursos facturables en una cuenta AWS; no fueron ejecutados sin autorización/cuenta;
- `cfn-lint 1.55.1`: 0 errores, 1 warning; runtime `python3.9` deprecado desde 2025-12-15;
- `pip-audit 2.10.1`: runtime e integración producen 30 ocurrencias/23 IDs únicos; paquetes afectados focales: Pillow 10.1.0 (23), pytest 7.4.3 (1) y urllib3 2.0.7 (6). No se ejecutó fix.

## Brechas materiales reproducidas

1. AWS exige left-padding al comparar sequencers hexadecimales de distinta longitud. El código usa comparación string directa. Harness: `old=f`, `new=10` devolvió `OUTCOME_OUT_OF_DATE`, descartando el evento numéricamente más nuevo.
2. `unlock()` y rollback no condicionan por `updated_by`; el source conserva el TODO de comprobar propiedad del lock. Una ejecución tardía podría liberar/revertir estado que ya no le pertenece.
3. No hay expiry/lease del lock; README reconoce que puede quedar tomado al expirar Lambda. El loop usa backoff dentro de Lambda sin budget explícito.
4. Template sin DLQ de target, retry policy explícita, alarmas, versionado, KMS, public-access block, PITR, TTL ni retention. EventBridge documenta que, agotados retries, descarta si no existe DLQ.
5. Runtime Lambda Python 3.9 ya está deprecado; las dependencias exactas tienen advisories.
6. La aplicación ejemplo procesa imágenes; no aporta receipt inmutable, quarantine/malware, checksum/version authority, storage release ni reconciliación del pipeline documental.

## Fuentes AWS primarias de contraste

- source exacto: `https://github.com/aws-samples/amazon-s3-endedupe/tree/7209a02a3cf7a613fb6d70e42caf0522b3eb0d7a`;
- ordering/sequencer: `https://docs.aws.amazon.com/AmazonS3/latest/userguide/notification-content-structure.html`;
- entrega duplicada/desordenada: `https://docs.aws.amazon.com/AmazonS3/latest/userguide/notification-how-to-event-types-and-destinations.html`;
- retry/DLQ EventBridge: `https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-rule-retry-policy.html`;
- lifecycle runtimes Lambda: `https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html`.

## Admisión

Conservar sólo como referencia exacta de conditional write y rollback testeados. No desplegar ni copiar como worker productivo hasta una revisión oficial que cierre: normalize/compare sequencer, lock ownership+lease fencing, timeout/retry/DLQ/replay/reconciliation, runtime/dependencies actuales, security/retention de recursos, integración real y gates de cuarentena/storage.

