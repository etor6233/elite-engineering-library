# PostgreSQL Transactional Foundation — Reconstruction Evidence V1

```yaml
evidence_id: "PG-TX-20260824-V1"
pack_id: "PG-TX-FOUNDATION"
pack_version: "0.1.0"
pack_sha256: "99bffa5965df0091a8a70cd5ef62f5009e2732c7d5d7b5d6981f7c9d00b24f23"
materialized_file_count: 3
ordered_file_hash_aggregate_sha256: "c74c40e3e6978056c8e3db44c531f92bcddcdc8bd5927b0ef12e7b2c6b4fd266"
database: "PostgreSQL 18.6 windows/amd64, page checksums enabled"
binary_archive_sha256_observed: "fbe23da234ee31547bf8a36d29dfd81e82b849df2d2b78d2eecb43d360252f8c"
result: "PASS_CONDITIONED"
verified_at: "2026-08-24"
```

## Entorno

Se usó el archive binario 18.6 enlazado por la página oficial de PostgreSQL para Windows y alojado por EDB. El digest anterior es el observado localmente; no se presenta como digest publicado por PostgreSQL/EDB. El cluster efímero se inicializó con UTF-8, locale `C`, page checksums y bind exclusivo a `127.0.0.1`.

## Gates ejecutados

```text
materialización y SHA de 3 archivos             PASS
0001 up                                         PASS
0001 tests con ROLLBACK                         PASS
0002 organization/order up                     PASS
0002 tests con ROLLBACK                         PASS
0002 down → 0001 down                           PASS
0001 up → 0002 up → ambos tests nuevamente      PASS
repository Go: create/replay/hash conflict/get/transition PASS
```

Los tests ejercitaron perfil activo único, formato de hashes, outbox versionado, auditoría append-only, FKs multi-tenant, estados de orden y optimistic concurrency.

## Límite del claim

No se ejecutaron todavía pruebas multi-session de carrera, workers outbox/inbox/jobs, grants mínimos/RLS, volumen representativo, backup/PITR/restore, failover ni upgrade con locks. El archive no tenía un SHA publicado visible en la fuente consultada; para CI/producción se exige imagen por digest firmado o artifact interno verificado.

