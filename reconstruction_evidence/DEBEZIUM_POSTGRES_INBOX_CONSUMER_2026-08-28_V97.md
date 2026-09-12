# Debezium PostgreSQL inbox consumer — evidencia V97

Fecha: 2026-08-28  
Estado: `REBUILD_VERIFIED / CONDITIONED`

## Autoridad y procedencia

- Red Hat/Debezium `debezium-examples` commit `7b0d765a02cf66ef29de2a01b67ec763e9acedf5`, 2026-08-25, repo activo Apache-2.0; commit unsigned.
- Archive oficial: 9.292.212 bytes, SHA-256 `0fba61c33bdbdb8f32eeab790335dcec412826dd98fa0bbab282afa2130f51e4`.
- El shipment consumer oficial compiló seis fuentes con Java 21/Quarkus 3.33.2 y ejecutó cero tests. Se preservan cinco fuentes Java `VERBATIM`; la licencia oficial conserva hash raw `58d1e17…9d8bd` y se marca `ADAPTED` sólo por el LF final requerido por el fence.
- La adaptación local no se atribuye a Red Hat. Inventario nuevo: 11 `AUTHORED`, 6 `ADAPTED`, 5 `VERBATIM`.

## Rechazos y runtime corregido

El sample oficial usa check-then-insert, `CompletableFuture.runAsync`, spans mutables, logs con payload, evento desconocido procesado y cero tests (`UP-FAIL-182`). Su BOM Quarkus 3.33.2 produjo SBOM 154 componentes/42 IDs OSV y se rechazó (`UP-FAIL-183`). Quarkus oficial 3.39.1 redujo el resultado a dos advisories de Jackson 2.22.0; se aplicó el fix oficial Maven Central `jackson-databind` 2.22.1 (`UP-FAIL-184`). El grafo final recompiló y OSV Scanner 2.5.1 terminó 154/0.

## Pack y hashes focales

- `DEBEZIUM-POSTGRES-INBOX-CONSUMER` 0.1.0, 22 archivos, SHA-256 del Markdown `3e6602298ef3464462a9c06dadf54a72d45bd33d60be215a6efd01710bc48936`.
- `pom.xml`: `d8415efbf3d2ea432b880785d4b30f255596b2c54ddc943b0c5f97b8a4377f70`.
- `source-lock.json`: `f3e0684e06744fb8fc4b093bee69b0d6d01afda2cde26d10dae170b7fa7e9552`.
- migration: `e098f2e1723e3ab37b43c478af6fd7ad0c39d992b9ba3f8923eb758744b18067`.
- JDBC processor: `fc2f239b01af68726f3f1810d01e39f5baf2a7bcdd223ec0b2d4db1d5df7bd9f`.
- Kafka consumer: `df239f0a5dd7a147bb6e0752dd2a64dae232760a2db2245e834e41cbb6ca015f`.
- PostgreSQL tests: `b70304c1d6e0c21208422976fcd6d3073fb7dddbe2842978d809fd83b5fbd9f1`.
- SCA receipt: `9df8929961da0579e078e570d54a988295819db2eaaf7b188d13b1338579a531`.

## Gates ejecutados

- materialización y comparación: 22/22, cero mismatch;
- static/source/profile/PowerShell: PASS;
- Java 21.0.8 + Maven 3.9.16 + Quarkus 3.39.1 + Jackson 2.22.1: BUILD SUCCESS;
- PostgreSQL 18.6 real: 10/10 tests PASS;
- cuatro carreras de ocho entregas: en cada una 1 `APPLIED` + 7 `DUPLICATE`;
- duplicate divergente, autoridad intake ausente, reuse de offset con event ID distinto y rollback: PASS;
- CycloneDX 2.9.1: 154 componentes; OSV Scanner 2.5.1: cero findings fechados;
- memoria: 1.081 fallos locales + 184 condiciones upstream = 1.265 IDs únicos, cero abiertos locales y cero duplicados; incluye la corrección de la descripción de licencia normalizada por LF.

## Límites

No se lanzó Kafka/broker/Connect ni se registró el consumer contra infraestructura real. TLS/auth, schema registry, poison/DLQ policy, consumer group rebalance, crash entre commit/ack, restart/redelivery/order/load, observabilidad, backup/restore, canary/rollback y mappings ERP/CRM del proyecto siguen condicionados. El código demuestra at-least-once con inbox y efecto local atómico; no promete exactly-once global.

El cierre V97 confirmó `VERIFY_LIBRARY_PASS` sobre 70 packs/673 archivos/423 Markdown/36 perfiles; email 14/138 y persistence focal 5/54. `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit` reconstruyó los 70 packs/111 fuentes y emitió `debezium_postgres_inbox_consumers=1`.
