# Debezium relational outbox — evidencia V1

Fecha de verificación: 2026-08-27. Alcance: código público oficial de Red Hat/Debezium para escritura transaccional outbox, transformación CDC y deduplicación de consumidor. Esta evidencia no afirma exactamente-una-vez, integración productiva ni seguridad futura.

## Fuentes oficiales fijadas

| Fuente | Identidad exacta | Archive oficial | Licencia |
|---|---|---|---|
| [`debezium/debezium`](https://github.com/debezium/debezium/tree/63371830ceff75437f95c775b91eeeaf55a056c1) | tag anotado `v3.6.1.Final` object `862e2c8cb6d5bb6704a3b910b9250d1be0be8318`; commit unsigned `63371830ceff75437f95c775b91eeeaf55a056c1`; tree `0bbe32a864a79a587abe5148ce16cb270b99c915` | 15.633.061 bytes; SHA-256 `3cc98e463d57b86029d991480893df4adf05577925137b4f465315ee0c94a938` | Apache-2.0; `LICENSE.txt` SHA-256 `58d1e17ffe5109a7ae296caafcadfdbe6a7d176f0bc4ab01e12a689b0499d8bd` |
| [`debezium/debezium-quarkus`](https://github.com/debezium/debezium-quarkus/tree/73f1c42da7400af500d6be8d9f524d3a9dee81bf) | tag anotado `v3.6.1.Final` object `69667c9eb90fd1192bc18370f0dd0ba2f04567bf`; commit unsigned `73f1c42da7400af500d6be8d9f524d3a9dee81bf`; tree `d611739c25ec45302353103233bce987ae7cbc88` | 740.306 bytes; SHA-256 `53c5f193aff2a845bc9b2bba239a5b9e292679c22e0829fa7a0fa0b2931df28e` | Apache-2.0; licencia SHA-256 `58d1e17ffe5109a7ae296caafcadfdbe6a7d176f0bc4ab01e12a689b0499d8bd` |
| [`debezium/debezium-examples`](https://github.com/debezium/debezium-examples/tree/7b0d765a02cf66ef29de2a01b67ec763e9acedf5) | main exacto unsigned `7b0d765a02cf66ef29de2a01b67ec763e9acedf5`; tree `54f8e5edff7b21430b047f825d9c04fbaecf33c8` | 9.292.212 bytes; SHA-256 `0fba61c33bdbdb8f32eeab790335dcec412826dd98fa0bbab282afa2130f51e4` | Apache-2.0; licencia SHA-256 `58d1e17ffe5109a7ae296caafcadfdbe6a7d176f0bc4ab01e12a689b0499d8bd` |

Los repositorios están activos. Debezium publicó 3.6.1.Final el 2026-08-04. Las identidades Git no tienen firma verificable; la adquisición depende del commit exacto, bytes y SHA-256 del archive.

## Código focal verificable

Core Event Router:

- `EventRouter.java` — `7142e407d2fedba4dac5760a74e4707737574f6590ae44f2aed985f886f3d735`.
- `EventRouterDelegate.java` — `f09f89070c286155cb99e1cb042ade6fd59b181f9335ac746197e31efac4e380`.
- `EventRouterConfigDefinition.java` — `811fe582b793a160264ff317aea0ef614b77329b210cdd1c3ccd6e38cc370b03`.
- `AdditionalFieldsValidator.java` — `a287c91382d40c2e6e8fac471d974545844e55c2cbbaf50c6cb5fdb32ffa647f`.
- `EventRouterTest.java` — `d2d345d24c4f80ab4f6344257c820c1567a53b057582e370e4c9ae74cb701266`.
- `JsonSchemaDataTest.java` — `a7a723fad9365588336d025ff4bb12e7a59452f88aaed02f19f1295c538b489a`.
- PostgreSQL `OutboxEventRouterIT.java` — `39b38dd17140ae7ac487ead6f73853e77ef0d25ebe6bfe3d64673efb45a27df8`.
- documentación oficial — `cccc4f94659e89c88fdd22ad08a47f6ba66ffa23d3a0083946216a65f0e2b86b`.

Quarkus writer:

- `ExportedEvent.java` — `db4c356cdc93325fe8566cfaa902e465ca9ae3d13a9f08fa1cd962fe9b368e5d`.
- `AbstractEventWriter.java` — `fe7b3c82bd518ab1ca5eca4eedaa7a3700f72679630173504fcd5797e53b91b5`.
- `DefaultEventDispatcher.java` — `08f307cacd38246e9ed5bfa9a6200e048d9ba14607b28e97ac74ca8a24daa980`.
- transacción ejemplo `MyService.java` — `801cf3305baf0546d7a2d69859e10b1f2085316b7638d4a2d8d1398c2dd24817`.
- integración `OutboxTestIT.java` — `2c063e70b35c8bdb23b3f3b347047133576ae137a5faeffefd38c6f5ab606ea4`.
- fixture PostgreSQL `DatabaseTestResource.java` — `45d6c7fffff3631df3d3027048fbe31e31672651100611faa339c7c636d77032`.

Ejemplo productor/consumidor:

- `OrderService.java` — `3290a16295a722967fda6d21e59cc3fe2538ba206f5ff0ab2165b1b2d8c644a5`.
- `OrderEventHandler.java` — `ccd1de2cb73a6630a08355fddcbba979bea7b8f3650e6ddb16d2590f0b1198ab`.
- `MessageLog.java` — `b60bc4cf5fd7a4dd0d783dad04e2595fa039f862af3e1f6dede931af9b74e55b`.
- compose — `188ce83f869cf31a48fb7aede6270acaaaa55ce77eb0618f35911ed213851cbd`.

## Toolchain fijado y pruebas

Microsoft OpenJDK 21.0.12.1 Windows x64: 201.096.952 bytes, SHA-256 oficial verificado `192441a9d27da813bada974bb88b4cf64d37a9589ed37f204374d411ca5ce07f`. Maven Wrapper 3.9.12: `mvnw.cmd` SHA-256 `4a361e1374a3e5ad6d03e18e9adc0cf181ac5058ac6203b76f0ba3b456b56481`; properties `e67de487071199117b7c9808dca5f56c73a44b656c925140ec68761947bdd384`.

Core, comando reproducido desde `<AUDIT_ROOT>/debezium`:

```text
mvnw.cmd -B -pl debezium-connect-plugins -am -Dtest=io.debezium.transforms.outbox.EventRouterTest,io.debezium.transforms.outbox.JsonSchemaDataTest -Dsurefire.failIfNoSpecifiedTests=false -DskipITs -Dformat.skip=true -Dcheckstyle.skip=true -Drevapi.skip=true test
```

Resultado: reactor de 15 módulos `BUILD SUCCESS` en 14:37; `EventRouterTest` 44 y `JsonSchemaDataTest` 4, total 48, cero fallos, errores u omisiones. Reports SHA-256: `5d08af065fe220eaf06c0880220b5f4438c8646ffc5afb843a7f738e80c33217` y `f733cb6b893b0dd2d5d32ae2c01af73e98665af02284f765cd88131ac33bd8b5`.

Quarkus runtime:

```text
mvn.cmd -B -f <QUARKUS_ROOT>/pom.xml -pl debezium-quarkus-outbox/runtime -am -DskipTests -DskipITs -Dformat.skip=true -Dcheckstyle.skip=true -Drevapi.skip=true package
```

Resultado: cinco módulos `BUILD SUCCESS` en 04:22. Common runtime JAR: 12.090 bytes, 19 entradas, contiene `ExportedEvent`, SHA-256 `84a8d845ae3222b3380c1d1d65927f867da5c3dbb9d486898ddf5dae671977a2`. Outbox runtime JAR: 11.435 bytes, 27 entradas, contiene `DefaultEventDispatcher`, SHA-256 `b5f54d9d18ef569aba86feb952f3ec3a828fe47a40d1ace23fc7350e6e9facda`. Los tests se omitieron: Docker no está disponible y la integración oficial necesita Testcontainers PostgreSQL.

Ejemplo:

```text
mvn.cmd -B -f <EXAMPLES_ROOT>/outbox/pom.xml -DskipTests -DskipITs -Dformat.skip=true -Dcheckstyle.skip=true package
```

Resultado: tres módulos `BUILD SUCCESS` en 06:29. Order JAR: 24.334 bytes, SHA-256 `8fffc5892248eb9445beb2955c2de76b4f4b086b638556543ba9fb9a0c3b28b5`. Shipment JAR: 13.980 bytes, SHA-256 `b811c3d296991e61e3135fdbd0419405c3105a5c0ff81a7ae9162d9ba704e1f6`. El ejemplo no trae tests automatizados propios.

## Supply chain fechada

Google OSV Scanner 2.5.1 oficial fue verificado con SHA-256 `25e42f5ef6711fd8c0fb45390972205891dd44c6bd02ac93f0f63e8e98d9bfb6`.

- ejemplo outbox: cero matches conocidos; JSON 124 bytes, SHA-256 `54345d5194d238b6017d0f3181f392109ce35023cc6b4bdfa8c134943dcb578e`;
- scans focales connect-plugins, Quarkus common y outbox: exit cero y cero matches conocidos en el alcance resuelto; mismo SHA del JSON vacío;
- core completo: 79 advisories únicos en 32 package-version de módulos no focales; JSON 1.739.298 bytes, SHA-256 `e40e0005d093ec18501a98480c4f8ebbda65f579fb7a66f1c6811942953550b3`; dos POM no pudieron resolver `kafka-connect-avro-converter:7.0.16`;
- Quarkus completo: 61 advisories únicos en 24 package-version; JSON 11.491.040 bytes, SHA-256 `60bdc5739be1a0cb4148cbdddff736c454b7ae724f3d5c6212226ed1c9370a03`; 52 advisories aparecen en el módulo separado `outbox-reactive/integration-tests`, no en los paths del outbox no reactivo seleccionado.

Cero matches no es garantía. Cada proyecto debe generar su lock/SBOM real, escanearlo, evaluar alcance y volver a probar antes de promoción.

## Admisión honesta

Estado de las tres fuentes: `CONDITIONED`. El core probado implementa el Event Router; Quarkus implementa persistencia dentro de la transacción actual; el ejemplo demuestra un consumidor que consulta y persiste el ID procesado en el mismo handler transaccional. Ninguna pieza garantiza por sí sola exactamente-una-vez, efecto empresarial idempotente, orden global, auditoría inmutable o recuperación.

Antes de producción siguen obligatorios: elección de runtime, esquema/event key, PostgreSQL logical replication, Kafka Connect o Engine, privilegios y WAL, compatibilidad, consumer idempotency, poison events, restart/redelivery/replay/order/concurrency/load, TLS/secrets, PII, observabilidad, costes, backup/restore, canary y rollback. El compose de ejemplo usa credenciales de demostración e imágenes mutables y no se copia como configuración productiva.
