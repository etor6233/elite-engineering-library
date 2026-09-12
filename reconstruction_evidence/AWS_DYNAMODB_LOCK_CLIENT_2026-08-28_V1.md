# AWS DynamoDB Lock Client 1.5.0 — evidencia V1

Fecha de auditoría: 2026-08-28 (`America/Argentina/Buenos_Aires`).

## Resultado

Clasificación: **`PINNED_CANDIDATE_CONDITIONED`**.

Es código oficial de AWS utilizable como componente Java de lock/lease con ownership y compare-and-set sobre DynamoDB. No es un worker object-event completo y no se admite con su grafo publicado sin resolver vulnerabilidades. Su `recordVersionNumber` es un UUID opaco que cambia con cada heartbeat: protege las mutaciones de la fila de lock, pero no es un fencing token monotónico para efectos externos.

## Identidad oficial y distribución

- repositorio: `awslabs/amazon-dynamodb-lock-client`;
- versión Maven Central: `com.amazonaws:dynamodb-lock-client:1.5.0`;
- commit de versión: `914934ee227b06e5e421c2a6674d7c7ec6340058`, tree `dd947680398701de6603ec3cfa39884bc7d41456`, unsigned;
- rama protegida actual observada: `master` en `c85c65e7e19d8d4e893f11fc894a03f9ddb271b9`, un merge firmado posterior que conserva SDK `2.32.10` y Log4j `2.17.1`;
- no existe release GitHub para 1.5.0 y el único tag enumerado es `1.1.0`; Maven Central declara 1.5.0 como latest/release;
- archive del commit: 118.137 bytes, SHA-256 `5eb0a7dffc199a5b434de5b1bcc3cb9b2360c4141e8dd40cde9150abcffdce70`;
- JAR oficial: 62.587 bytes, SHA-256 `277f350b1f3b274a4a69838ddd0bd72b9168e1cebf049be72f20ac212759b6d4`;
- sources JAR: 54.258 bytes, SHA-256 `f8a2706f924b920efee417d7d54deb83f83f2cb8dfd1d1a13d9e164730812a8c`;
- POM oficial: 27.799 bytes, SHA-256 `eac080bd5c3fe86a02b0f79ecbfe2429bef41639c854c99424a7178a1d01c6e0`;
- source POM y POM Maven tienen el mismo SHA-256;
- las 18 fuentes Java de `src/main/java` coinciden byte a byte con las 18 del sources JAR;
- las firmas detached de JAR, sources JAR y POM producen `GOODSIG` y `VALIDSIG` con fingerprint `79D6089D38724F977A7DFD922B60A375F5B66ABF`, clave `Dylan Wilson <dwlsn@amazon.com>`; la clave fue obtenida por fingerprint y no está certificada por un trust root local;
- licencia: Apache-2.0, `LICENSE.txt` SHA-256 `7382b1cf711e7a7c7af9816e0e07b49b91e7149b7c6d347225f98800fcb1f1c2`;
- `NOTICE.txt` SHA-256 `43e948d7cf913af1c2f50a64e6feac9c45b7078455ce5e7e431d044c61c7440d`.

## Archivos focales exactos

- `README.md`: `11c1f2677f095990fd66b74a394b28b55758dd65bbc2f4fc24828958f46a6f27`;
- `src/main/java/com/amazonaws/services/dynamodbv2/AmazonDynamoDBLockClient.java`: `46e95542eb149d34d6dc6aea78f52cab07f8e1afcc4cfb7082c8094abc30cb47`;
- `src/main/java/com/amazonaws/services/dynamodbv2/LockItem.java`: `e2247abb9d910893ef64a34153ac98898201d94ddd354fa5bcc4ea7f418ebcdf`;
- `src/test/java/com/amazonaws/services/dynamodbv2/AmazonDynamoDBLockClientTest.java`: `23e1baecd92d8e0309d858598e00652b38cf88bb7c14792875d343b4aa23e141`.

## Ejecución reproducida

Tooling aislado:

- Apache Maven 3.9.16, ZIP 9.395.475 bytes, SHA-512 `ed41650d42485cfc243fad22158caf9cbb5dc408ce7a09ddb94dd42a019de929ca43065bfa450612cf12bf78b5cafa3884b96c090de326ff590448c933454af3`;
- Eclipse Temurin JDK 8u504-b01, ZIP 106.457.258 bytes, SHA-256 `ea43d46ede95b51e44a12c66711706cddc762e0a766c54bccea18954e902b2aa`;
- Google OSV-Scanner 2.5.1, 58.970.112 bytes, SHA-256 `25e42f5ef6711fd8c0fb45390972205891dd44c6bd02ac93f0f63e8e98d9bfb6`.

Resultados del source oficial sin cambios:

- 18 fuentes main y 18 fuentes test compilaron;
- 12 reportes Surefire: 112 tests, 0 failures, 0 errors y 0 skips;
- `clean package` produjo un JAR local de 62.773 bytes/SHA-256 `09c32bf3e8d0316f52cbb7395d315d1761f8bb7bcfb539fc7ee70b0c632cd069`;
- la suite no ejecuta una integración real con una cuenta DynamoDB;
- bajo JDK 21 el código main compila, pero PowerMock 2.0.0 falla antes de ejercer gran parte de la suite por acceso modular a `java.lang`; no existe PASS de suite JDK 21;
- los 30 class paths del JAR oficial y local coinciden, pero 8 class files no son byte-idénticos y cuatro difieren también en el texto `javap`; el build local no reproduce el JAR firmado, por lo que la distribución gobernante es el JAR oficial.

## Vulnerabilidades del grafo publicado

Reporte OSV source: 210.836 bytes, SHA-256 `08af373f568cd96a69855a306634caf767ef56f950176622fae8abfecc9b7cf4`, exit 1.

Después de deduplicar la representación del reporte hay 15 advisories únicos en cuatro package-version:

- runtime, `io.netty:netty-codec:4.1.118.Final`: 3 IDs, dos HIGH y uno MODERATE;
- runtime, `io.netty:netty-codec-http2:4.1.118.Final`: 8 IDs, cinco HIGH y tres MODERATE;
- test, `org.apache.logging.log4j:log4j-core:2.17.1`: 3 MODERATE;
- test, `org.apache.logging.log4j:log4j-api:2.17.1`: 1 MODERATE.

Los fixes OSV observados exigen como mínimo Netty 4.1.136.Final para el conjunto y Log4j 2.25.5 para la rama 2.x correspondiente. El AWS SDK Java v2 oficial 2.54.6, publicado 2026-08-27, declara Netty 4.1.137.Final. El lock client 1.5.0 y su master observado no incorporan esa actualización.

## Update candidate no atribuido a upstream

En una copia temporal rotulada `AUDIT_CANDIDATE_NOT_UPSTREAM` se cambiaron únicamente:

- AWS SDK Java v2 `2.32.10` → `2.54.6`;
- Log4j `2.17.1` → `2.25.5`.

El effective POM confirmó ambas versiones; la misma suite terminó 112/112 PASS y OSV-Scanner produjo exit 0, dos resultados, 58 package entries y cero vulnerability objects. Reporte: 11.389 bytes, SHA-256 `39c6909992e1cfc1f1abb8f8df53a72423824079c3b73e96f6e3a2fb6832c7a2`.

Esto prueba un candidato de actualización con releases oficiales, no una release AWS Lock Client ni compatibilidad productiva universal. El proyecto debe fijar el grafo completo, repetir tests/SBOM/licencias/SCA y conservar rollback; no se modifica el artefacto JAR firmado.

## Contrato de locking y límite de fencing

El código oficial sí implementa:

- owner name, lease duration y heartbeat;
- adquisición condicional de lock nuevo, liberado o expirado;
- release y heartbeat condicionados por primary key, owner y `recordVersionNumber`;
- cambio de `recordVersionNumber` en cada heartbeat;
- detección de lock robado y detención del heartbeat;
- session monitor y opción de no bloquear cuando el lock está ocupado.

El código no demuestra por sí solo:

- fencing monotónico: `generateRecordVersionNumber()` devuelve `UUID.randomUUID().toString()`;
- rechazo de un holder vencido por un storage/servicio externo que no participa en la condición DynamoDB;
- transacción atómica entre lock y efecto empresarial fuera de la misma frontera DynamoDB;
- ordering por S3 sequencer/version;
- SQS partial batch, visibility timeout, DLQ, replay y reconciliation;
- idempotencia del efecto de negocio;
- IAM/IaC, alarmas, backup/restore, load/soak, outage y E2E AWS del target.

## Decisión de uso

El agente puede adquirir inmediatamente el JAR/sources/POM oficiales por hashes y usar el API de lock como componente condicionado. Debe bloquear producción si el grafo conserva los 15 advisories, si usa `holdLockOnServiceUnavailable=true` sin aceptar explícitamente seguridad degradada, o si un efecto externo no valida un fencing/ownership token vigente.

Para el worker documental/object-event aún se requiere una composición separada y declarada: Powertools para idempotency/partial batch, lock/lease o transacción adecuada, ordering, DLQ/replay/reconciliation e IaC. Cualquier glue o política creada por Elite/proyecto se etiqueta `AUTHORED`; nunca se atribuye a AWS.
