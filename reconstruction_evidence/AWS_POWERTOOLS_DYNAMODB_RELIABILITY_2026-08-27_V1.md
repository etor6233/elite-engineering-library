# AWS Powertools TypeScript 2.35.0 — batch parcial e idempotencia

## 1. Decisión

`aws-powertools/powertools-lambda-typescript` release `v2.35.0` queda fijado como `PINNED_CANDIDATE_CONDITIONED` para dos claims estrechos:

- partial batch response para SQS, Kinesis y DynamoDB Streams;
- idempotencia con persistencia DynamoDB, estados `INPROGRESS`/`COMPLETE`, retry y exclusión de ejecuciones idénticas concurrentes.

No se admite como implementación completa de outbox, auditoría durable ni exactly-once. AWS advierte que partial batch reduce pero no elimina el procesamiento duplicado, y el release contiene un ejemplo DynamoDB Streams para batch y otro ejemplo SQS para batch+idempotency; no contiene una composición oficial DynamoDB Streams+idempotency+side effect+audit durable. `UP-FAIL-109` permanece abierto.

## 2. Autoridad e identidad exacta

- repositorio oficial: <https://github.com/aws-powertools/powertools-lambda-typescript>;
- documentación AWS batch: <https://docs.aws.amazon.com/powertools/typescript/latest/features/batch/>;
- documentación AWS idempotency: <https://docs.aws.amazon.com/powertools/typescript/latest/features/idempotency/>;
- comportamiento DynamoDB/Lambda: <https://docs.aws.amazon.com/lambda/latest/dg/services-ddb-batchfailurereporting.html>;
- release: `v2.35.0`, publicada `2026-08-18T16:41:12Z`;
- tag object: `54a898b0fb1713ed6180ef0b5f323a82fa9a6957`, anotado pero unsigned;
- commit target verificado: `7bcc27b1574493f9452688673658f52b80c53847`;
- tree: `a4f3fc3b4146c63234b6fbb2b799b56640f41d3d`;
- ZIP commit: 3.617.192 bytes, SHA-256 `6c2a59a54f2f964bf14465740f5d33b0fc704b3961a0992e93d71fdc3b2a4805`;
- contenido: 1.395 archivos, 7.608.014 bytes sin comprimir;
- licencia MIT-0: `LICENSE` SHA-256 `5cf66a4554832601b0277c7a539b8d8e8cf5912f53d91e4f83ecf6cb67f1cce6`;
- `NOTICE` SHA-256 `d4290ed64c2edd0fce1d84e3f9dfb2881240fe534def76b8cd29ed6af683e287`;
- `package-lock.json` SHA-256 `b164477bcca8442eb97609859fc38b277fc4443de26cb8e1c4d0b3c024606965`.

Paquetes oficiales npm:

| Paquete | Bytes | SHA-1 registry | SHA-256 observado | Integrity registry |
|---|---:|---|---|---|
| `@aws-lambda-powertools/batch@2.35.0` | 31.198 | `900db7fd39379e64a5c422678d05c9f81c477603` | `248829873612091b3061a7f0ac5ed8f7c368a8ba4b2751d45b6205429facdef0` | `sha512-RdM7gs6YmblfK/llGBt1/u32XgQjGMa1sr9zlhULRzMdumAnmnPGYOqXqEGTu4gjfZAnMdYldNntYVSGkGQ8rQ==` |
| `@aws-lambda-powertools/idempotency@2.35.0` | 42.520 | `e0e840d79076926f03463aa93c7be3bbc0ebfb56` | `620a7093261b1204da5392f2400dc9e3df8babef43093aae477f541637b1a502` | `sha512-l4WJCWgf+K6CFJCfKtxgORAJH1x9Y+WlytQqEwFHu+sZqqNKlj3YJiK1VaWbpwBeU2z8Wx/ePX4r6hR5FLtAbQ==` |

El source lock conserva quince artefactos adicionales —implementación, tests, manifests, documentación y ejemplos— con path y SHA-256 exactos.

## 3. Rebuild soportado

Entorno aislado WSL2/Linux:

- kernel `6.18.33.1-microsoft-standard-WSL2` x86_64;
- Node `v22.22.3`, dentro del engine oficial `>=22`;
- npm `10.9.8`;
- install `npm ci --ignore-scripts --no-audit --no-fund` desde el lock exacto;
- builds ESM/CJS, en orden de workspace, de commons, jmespath, parser, testing-utils, logger, idempotency y batch;
- `build:tests` de batch e idempotency PASS;
- Biome `lint:ci`: 37 archivos batch y 48 idempotency PASS.

Suites focales oficiales:

- batch: 9 archivos, 115 tests PASS; threshold instrumentado 100% statements/branches/functions/lines para los sources incluidos;
- idempotency: 9 archivos, 130 tests PASS; mismo threshold instrumentado 100%; incluye 20 tests de `DynamoDbPersistenceLayer` y 31 de `makeIdempotent`.

No se ejecutó E2E AWS porque requiere cuenta, IAM, recursos potencialmente facturables y autorización de cleanup. El PASS local no demuestra Lambda event source mapping, DynamoDB, concurrency distribuida, DLQ/failure destination, replay, restore ni costo reales.

## 4. Supply chain

Google OSV-Scanner oficial 2.5.1, commit `c84fa4568f2526d0333e9a914ea8a0a5f74ad68b`, binario Windows SHA-256 `25e42f5ef6711fd8c0fb45390972205891dd44c6bd02ac93f0f63e8e98d9bfb6`:

- lock del monorepo: 590 paquetes; exit no-cero; 10 advisories en cinco package-version: `brace-expansion@2.0.3`, `brace-expansion@5.0.6`, `fast-uri@3.1.2`, `nanoid@3.3.15`, `postcss@8.5.16`;
- npm también advierte por `glob@10.5.0`, alcanzado desde testing-utils→AWS CDK;
- lock mínimo batch: 4 paquetes, 0 findings;
- lock mínimo idempotency + `@aws-sdk/client-dynamodb@3.1120.0` + `@aws-sdk/lib-dynamodb@3.1120.0`: 34 paquetes, 0 findings.

El primer lock idempotency con client/lib `3.1106.0` produjo un peer override porque `util-dynamodb@3.996.9` requiere client `^3.1111.0`. El candidato se corrigió a client/lib `3.1120.0`, versión oficial publicada el 2026-08-27, y resolvió sin warning ni `--force`. Estos locks de consumidor son evidencia fechada y no sustituyen el lock que debe generar, preservar y reescanear cada proyecto.

## 5. Límites de corrección

AWS Lambda establece que:

- `ReportBatchItemFailures` debe habilitarse en el event source mapping; devolver la respuesta desde código sin esa configuración no tiene efecto;
- con múltiples fallos de stream, Lambda usa el menor sequence number como checkpoint y puede reprocesar registros posteriores ya exitosos;
- partial batch no elimina la posibilidad de reintentar un registro exitoso;
- `BisectBatchOnFunctionError`, retry, maximum record age y destination cambian la recuperación y deben configurarse por proyecto.

Powertools idempotency aporta exclusión concurrente y retry para una clave/payload, pero no vuelve atómico un side effect externo con el cambio de estado de idempotencia. Un key incorrecto, TTL insuficiente, timeout, respuesta no serializable, persistencia no disponible o efecto remoto no idempotente conserva riesgo. No se atribuye a AWS ninguna adaptación local para resolver esas fronteras.

## 6. Resultado reusable y blocker restante

Listo para adquisición inmediata, no para adopción ciega:

- source archive exacto y dos tarballs npm exactas;
- licencia/notice/manifests/15 artefactos fijados;
- build, typecheck, lint y 245 tests focales reproducidos en POSIX;
- locks runtime mínimos con consulta OSV cero;
- perfil `authenticated-human-approval-aws` ahora selecciona la aprobación AWS y Powertools como dos fuentes condicionadas.

Para cerrar audit durable aún se requiere uno de dos caminos legítimos:

1. localizar una implementación oficial mantenida que ya componga DynamoDB Streams, idempotencia y ledger/side effect durable con las pruebas necesarias; o
2. si el propietario autoriza adaptación, crear código local con provenance `AUTHORED_ADAPTER`, sin atribuirlo a AWS, y demostrar event key, atomicidad/reconciliación, carreras, retries, duplicate delivery, poison record, DLQ/replay, restore, IAM, observabilidad y live sandbox.

V48 admite código oficial de AWS por claims reales y mantiene el hueco abierto donde AWS no publicó una solución completa.
