# All implementation packs — materialización V49

Fecha: 2026-08-27. Este snapshot sucede a V48 e incorpora la adquisición oficial condicionada de Red Hat/Debezium para outbox relacional. V48 y anteriores quedan como evidencia histórica, no como estado vigente.

## Estado canónico

- 46 implementation packs.
- 415 archivos materializables desde bloques Markdown con manifest y SHA-256.
- 308 Markdown antes de esta evidencia; 309 incluyéndola.
- 23 perfiles de composición.
- 87 fuentes oficiales fijadas y 12 perfiles source.
- 418 fallos locales retenidos y 118 condiciones upstream retenidas.
- `engineering_execution_kit/`, `node_modules/`, caches, builds y temporales no forman parte de la fuente canónica.

## Delta Debezium

`OFFICIAL-UPSTREAM-ACQUISITION-CORE` avanzó de 0.4.39 a 0.4.40 y de 20 a 21 archivos. El lock pasó de 84 a 87 fuentes:

1. `debezium/debezium` 3.6.1.Final exacto — core Event Router, 48 tests focales PASS;
2. `debezium/debezium-quarkus` 3.6.1.Final exacto — common/runtime y outbox/runtime compilados en dos JAR reales;
3. `debezium/debezium-examples` commit exacto — sample outbox productor/consumidor empaquetado, sin tests propios.

El perfil nuevo `relational-outbox-debezium.json` selecciona esas tres fuentes y no autoriza adquisición hasta responder runtime, PostgreSQL, logical replication, Kafka Connect/Engine, esquema, partición/orden, consumer idempotency, tenancy, SLO, costes, seguridad, recuperación y rollback. Preserva explícitamente que Debezium no garantiza exactamente-una-vez del efecto empresarial ni reemplaza un audit ledger.

La evidencia técnica y de supply chain está en `DEBEZIUM_RELATIONAL_OUTBOX_2026-08-27_V1.md`. Las condiciones abiertas se retienen como `UP-FAIL-116/117/118`.

## Gates ejecutados

```text
materialize_markdown_pack.ps1 OFFICIAL_UPSTREAM_ACQUISITION_CORE
Materialized 21 files

test_upstream_acquisition.ps1
UPSTREAM_ACQUISITION_TEST_PASS negatives=7

test_source_profiles.ps1
SOURCE_PROFILE_TEST_PASS valid=12 negatives=5 positives=1
SOURCE_PROFILE_VALID profile=relational-outbox-debezium selected=3

VERIFY_LIBRARY.ps1
VERIFY_LIBRARY_PASS
packs=46 materialized_files=415 markdown_files=308

VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Audit
VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=46 upstream_sources=87 document_sdk_artifacts=5 business_central_artifacts=1 provider_adapters=9 document_orchestrators=1 evidence_logs=1
```

Los 23 perfiles compusieron. Los perfiles que incluyen el pack de adquisición aumentaron un archivo por el nuevo source profile; backend 95, web/BFF 44, Mercado Libre 10, payment 7 y Business Central 6 permanecieron iguales según su selección.

## Condición de cierre

La biblioteca pasa sus gates estructurales y ejecutables de auditoría. Debezium queda disponible inmediatamente como fuente oficial exacta y guiada, pero `CONDITIONED`: no se declara integración productiva sin Docker/PostgreSQL/Kafka, pruebas commit/rollback, restart/redelivery/replay/order/concurrency/load, lock/SBOM del proyecto, seguridad, observabilidad, backup/restore, canary y rollback reales.
