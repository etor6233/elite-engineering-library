# AWS Receipt Mobile Ingestion Candidate — 2026-08-28 V1

## Resultado

`aws-samples/sample-ai-receipt-processing-methods` aporta código oficial y licenciado para PWA/cámara, presigned POST, S3 event, Step Functions y cuatro métodos de extracción. No está listo para adopción inmediata. Clasificación: `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`.

No se copió código al producto ni se escribió una corrección local atribuida a AWS.

## Identidad exacta

- repo: `https://github.com/aws-samples/sample-ai-receipt-processing-methods`;
- default branch: `main`;
- head unsigned: `6eb72bf0060e4fb38074d86fe9242e781db9df97`;
- fecha del commit: `2025-11-18T21:47:41Z`;
- archive: 8.640.517 bytes;
- archive SHA-256: `844add2f964397361331d3554cc1ac9bd0ef5d5c3b8b07812b1c693740253674`;
- inventario: 132 archivos y 35 directorios;
- tags reales: 0; releases reales: 0;
- licencia MIT-0: 947 bytes, SHA-256 `7cb750713252efd1d578837ba8785b61319109906f0c10b87adab7cf4badfc42`.

## Archivos focales

| archivo | SHA-256 |
|---|---|
| `README.md` | `b9894d5bf1303049795cf947e1e07ca67196732e7b5228b27f613c0bde6609b2` |
| `package-lock.json` | `f476e6dd357fd82d4ff24a5c96999e4f6674097fc68d622cad6c439d84db5864` |
| `src/frontend/package-lock.json` | `7abfdcc3df130b4117a3cffbe971269b788005f5e6214f08118aec992e856284` |
| `src/frontend/src/components/Camera/CameraComponent.tsx` | `3b05327b43cb77e7cd3de52404f0783d35a3a7916189529cc8c44c05ee2d8111` |
| `src/frontend/src/utils/s3-upload.ts` | `861cb5af1ebd5c3d2af5b525431a07b383f6514ffefeff67ea134c97f8b477f3` |
| `src/lambda/receipt-presigned-url/index.js` | `37291281bd31b45f6495941cc4758a0494e56c13880d0112a1dbc187e2324e6b` |
| `src/lambda/receipt-status/index.js` | `21b066d8830b32727c64c597b6f4842d2ba30570b3f0787180fb5c20d07a86f0` |
| `src/stepfunctions/receipt-processing-workflow.json` | `99741254a6d147e74b56c10b8a3d1aa74af7b7d41e8571f9669221c330fb0b2b` |
| `infrastructure/lib/stacks/storage-stack.ts` | `c1a961083256f9d470ab330822e56f967f90e3b54dab8d0eb46d019be701ea0c` |
| `infrastructure/config/environment-config.ts` | `5b11d76c3a24eb84387646a065963ad8725785f89ea152dfd992f1ab8396c477` |

## Ejecución reproducida

- 19 JavaScript pasan `node --check`;
- 8 Python pasan `py_compile`;
- `npm ci --ignore-scripts` instaló 598 paquetes raíz y 493 frontend;
- Vite produjo PWA y el TypeScript de infraestructura compiló;
- `npm test -- --runInBand` falla: el script apunta a `tests/integration/`, pero el archive no contiene `tests/`;
- cero workflows CI;
- SCA productiva fechada: raíz 14 vulnerabilidades conocidas (6 moderadas/8 altas) y frontend 13 (2 moderadas/11 altas). No se ejecutó `audit fix`.

## Brechas materiales observadas

1. `receipt-status` acepta `?test=true` y `?debug=true` antes de exigir claims; el debug expone configuración/estructura de headers y el handler registra el header Authorization completo.
2. El ID usa tiempo+`Math.random`; el registro DynamoDB se escribe antes de completar el POST. No existe receipt byte/hash/version que vincule solicitud, objeto y procesamiento.
3. La política del POST limita tamaño/MIME declarado y SSE-S3, pero no exige checksum; tampoco valida magic bytes antes de activar OCR.
4. El entorno demo desactiva versionado S3 y PITR DynamoDB. No hay Object Lock, quarantine/malware ni storage authority.
5. El bridge S3→Step Functions usa `receiptId+eventTime`, ignora `sequencer/versionId/eTag`, captura cada fallo de `StartExecution` y termina normalmente para evitar retries. AWS documenta que S3 entrega al menos una vez, puede duplicar/desordenar y recomienda comparar `sequencer` por key.
6. El workflow actualiza estados sin `ConditionExpression`, persiste automáticamente cuando una confianza global es mayor a 50 y almacena respuesta raw; no existe decisión por campo, corpus/threshold por clase, revisión humana, CAS ni evidencia durable.
7. El frontend marca upload como completo al recibir 2xx de S3; `processingTime` se calcula como `Date.now() - Date.now()`. Reintentar el POST no demuestra que se reusa el mismo receipt/object.

## Fuentes AWS primarias de contraste

- source exacto: `https://github.com/aws-samples/sample-ai-receipt-processing-methods/tree/6eb72bf0060e4fb38074d86fe9242e781db9df97`;
- S3 event notifications at-least-once/duplicadas/desordenadas: `https://docs.aws.amazon.com/AmazonS3/latest/userguide/notification-how-to-event-types-and-destinations.html`;
- `sequencer` y `versionId`: `https://docs.aws.amazon.com/AmazonS3/latest/userguide/notification-content-structure.html`;
- integridad/checksums de objeto: `https://docs.aws.amazon.com/AmazonS3/latest/userguide/checking-object-integrity-upload.html`;
- idempotencia limitada de `StartExecution`: `https://docs.aws.amazon.com/step-functions/latest/apireference/API_StartExecution.html`.

## Admisión

Conservar sólo como referencia exacta de UX cámara/PWA, composición CDK y comparación de métodos OCR. No desplegar, no persistir resultados y no copiar fragmentos al perfil reusable hasta que una revisión oficial o una derivación claramente propia demuestre: tests reales/CI/SCA limpia, auth sin bypass/log de token, checksum+magic bytes+quarantine, version binding, event idempotency/retry/DLQ/reconciliation, CAS de estados, evidencia por campo, revisión y storage authority.

