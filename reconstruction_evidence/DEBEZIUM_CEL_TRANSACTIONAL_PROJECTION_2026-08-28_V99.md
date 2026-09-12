# Debezium + CEL transactional document projection — evidencia V99

Fecha: 2026-08-28  
Estado: `REBUILD_VERIFIED / CONDITIONED`

## Autoridades exactas

- Red Hat/Debezium Examples commit `7b0d765a02cf66ef29de2a01b67ec763e9acedf5`, archive SHA-256 `0fba61c33bdbdb8f32eeab790335dcec412826dd98fa0bbab282afa2130f51e4`, Apache-2.0; cinco Java `VERBATIM` y licencia `ADAPTED` sólo por LF final;
- Google-origin CEL-Java `v0.13.0`, commit `c57d9d02259795c761e82f8870303b93a7a672c8`, JAR Maven Central SHA-256 `a3bc22d78affe749f81b17cce2131a3807275436157ba4d6fe3f10d8f541ea94`, Apache-2.0;
- Quarkus 3.39.1, Jackson Databind 2.22.1, Java 21.0.8, Maven 3.9.16, PostgreSQL 18.6, CycloneDX 2.9.1 y OSV Scanner 2.5.1.

El código de integración es `AUTHORED`: no se atribuye a Red Hat ni Google y no copia source CEL. Las condiciones upstream V97/V98 siguen visibles, incluidos los cinco fallos `ErrorsTest` Windows de CEL-Java.

## Unión material demostrada

`DEBEZIUM-POSTGRES-INBOX-CONSUMER` 0.2.0 materializa 24 archivos. Tras validar evento/topic/key/headers, abre una transacción PostgreSQL, inserta inbox, exige que el `domainType` del evento coincida exactamente con la clase aprobada del mapping, lee el `document_payload` revisado bajo tenant/RLS, ejecuta una expresión CEL single-line con macros vacías, valida claves exactas y persiste la proyección con mapping ID/clase/version/expression SHA, input/output schema ID/SHA, source payload SHA y mapped payload SHA. Sólo después del commit se ejecuta ACK Kafka.

Una redelivery idéntica vuelve a calcular y comparar toda la autoridad de mapping. Evento cross-class, mapping inválido, input ausente, hash/schema/config divergente, evento divergente o fallo de proyección revierte inbox y proyección; no se reconoce el mensaje. El template nace con clase, mapping, schemas, owners, approvals y proofs vacíos/false.

## Gates ejecutados desde Markdown

- round-trip fuente→pack→materialización: 24/24 hashes exactos;
- static/source/profile/PowerShell: PASS;
- Java/PostgreSQL real: 18/18 PASS —6 mapping, 5 envelope, 7 transacción—;
- cuatro carreras de ocho entregas: 1 `APPLIED` + 7 `DUPLICATE` por ronda;
- clase documental cruzada y mapping inválido revierten inbox y proyección; offset reutilizado, intake ausente y duplicate divergente: PASS;
- los runners PostgreSQL y SCA fueron invocados desde fuera del componente y resolvieron el `pom.xml` desde `$PSScriptRoot`: PASS;
- CycloneDX 1.6: 170 componentes; OSV Scanner 2.5.1: cero findings; identidades Quarkus/Jackson/CEL verificadas;
- `VERIFY_LIBRARY_PASS`: 71 packs/682 archivos/426 Markdown, email 14/140 y persistence 5/56; `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit` cerró 71 packs/111 fuentes con consumer y mapper contados;
- pack Markdown: 116.990 bytes, SHA-256 `f044c0b63ee3dd5624c6cb8de742710284dfe8e178e3897b06bcd52007efff05`.

Hashes focales: mapper `ec2f6a2c27d5cc1383553139daaf934e281e1733cf73e2f956910bbe6c60e356`; processor `fb571ec0d38dd20ee0db2a10f07ff845aad9c320d0db5f08ca76728c3ea27400`; migration `ec8e041bccd303361a2ef93722f2d6c7f8ad77314f793ec9d92d49df5061dd4c`; tests mapping/JDBC `9686d84705e76d4f1c26b779a8357437cdd491bce8cdea648c3577b5191fe663`/`23b0f4d867ccaa9d88b9daceb51a2aecdb95e488ae2c49e3a2acdc2640541a86`; profile cerrado `c37930a730631f997df56cd5b6da80b22a34ca4893ef66d102e297abf0c27445`; runners PostgreSQL/SCA `4411659a3e536e7b3006f5c8089ed15bca64fb4211f11c2e47b7731e798c0c7b`/`41eb917976c219fa5575887ba696381ece3f7cd4287cb05e5a1a7404ea197ae2`.

## Decisión y límites

La unión local mapper→inbox/proyección→commit ya no es sólo composición de carpetas: es un recorrido ejecutable y probado. Continúa `CONDITIONED` porque un proyecto debe aportar y aprobar clase, schemas, expresión, corpus/ground truth, límites, rollback y mapping ERP/CRM; además faltan Kafka/Connect/broker/TLS/auth/rebalance/crash entre commit/ACK, load/recovery/restore/observabilidad/canary live. No se afirma exactly-once global ni exactitud documental universal.

Memoria al cierre focal: 1.097 fallos locales + 185 condiciones upstream = 1.282 IDs únicos, cero abiertos locales.
