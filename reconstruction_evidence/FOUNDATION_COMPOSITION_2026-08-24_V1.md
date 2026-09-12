# Foundation Composition — Integrated Evidence V1

```yaml
evidence_id: "FOUNDATION-COMPOSITION-20260824-V1"
composition_version: "1.0"
plan_sha256: "6481fa11a2f0c3103f479629641c98e5d0113d1ed2266cc5f5ffdf02cabba2c3"
materialization_record_sha256: "e7fae509a2dadbf445e6541ede117e42a000795a0a66b04a00691ce8cebc6bcb"
selected_pack_count: 4
materialized_file_count: 35
record_file_count: 1
result: "PASS_CONDITIONED"
verified_at: "2026-08-24"
```

## Packs exactos

| Pack | Versión | Pack SHA-256 | Estado |
|---|---:|---|---|
| `PBC-CORE` | 0.1.0 | `ff3166eef7645a87f43b8fc3ffb7a452a3d3e6b80d2646b8cbe01bb8f32ab1de` | REBUILD_VERIFIED / CONDITIONED |
| `PG-TX-FOUNDATION` | 0.1.0 | `861717d65c5f551d7e7ad67c29ef1300be43b474e802e0691e092780eab87e4d` | REBUILD_VERIFIED / CONDITIONED |
| `GO-ENTERPRISE-BACKEND` | 0.3.0 | `bace6b83561ddfcd1fd121a58393f605414bc415e1ca97600bcebd9af099748a` | REBUILD_VERIFIED / CONDITIONED |
| `EXECUTION-VALIDATOR` | 1.1.0 | `10f02173e2b2ddbee564df18d50617378be88a7c21dd995f0087b87bf6520967` | REBUILD_VERIFIED / CONDITIONED |

## Gates integrados

```text
compositor preflight, hashes, collisions y lineage          PASS
PBC schema/example materializados                           PASS
PBC Python semantic tests                                   PASS (3)
PBC Go semantic tests                                       PASS
PostgreSQL 18.6: 0001 + 0002 up y tests con rollback        PASS
Go: unit + HTTP + repository + outbox integration           PASS
Go: vet + build                                              PASS
Execution Validator hermético con authority root externo     PASS (7 tests + plan)
```

La primera ejecución detectó y rechazó metadata CRLF; se corrigió el compositor, se añadió la regresión y recién entonces se repitió la composición. Un error posterior en la URL de test fue del comando PowerShell (`${db}`), no del código materializado; corregido el harness, los tests de repository/outbox pasaron.

## Qué demuestra y qué no

Demuestra que el agente puede componer desde Markdown una foundation stack-neutral, configuración validada, PostgreSQL, backend Go con OIDC boundary, outbox concurrente y contract validator. No demuestra todavía frontend E2E, todos los módulos empresariales, issuer real, provider integrations, CI/CD, observabilidad, supply chain, secrets/PKI, load/security ni recovery. Por eso el sistema completo continúa `NOT_READY`.

