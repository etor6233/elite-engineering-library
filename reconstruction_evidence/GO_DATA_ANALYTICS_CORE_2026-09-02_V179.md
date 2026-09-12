# Go Data Analytics Core — V179

## Resultado estrecho

V179 materializa `GO-DATA-ANALYTICS-CORE 0.1.0`, el contrato de ingesta (`DATA-INGEST`) y de métricas semánticas (`ANALYTICS-BI`): landing inmutable con replay idempotente, late-data bound y ownership de fuente; definiciones de métrica versionadas con freshness y reconciliación; todo tenant-scoped.

No afirma warehouse/lakehouse ni CDC productivo: define el landing y sus invariantes, y el versionado/freshness/reconciliación de métricas. El motor de agregación sobre el source of truth y la plataforma analítica son decisión del proyecto.

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — core `AUTHORED`.
- PostgreSQL 18.6 (PostgreSQL License) — migración `0041` `AUTHORED`.

## Archivos materializados (7)

| Archivo | SHA-256 |
|---|---|
| internal/analytics/metric.go | 9fd7359576bc1f14e1a91aacb90ea9fba8ea3e74aae1e9e3a22066a40a900637 |
| internal/analytics/snapshot.go | 59e1d536351cccb279961a5ff6ab099dc9a64f36575d6f123f30a0e73f6673a8 |
| internal/analytics/ingest.go | 5623eec2192192bada127e1492f5e7bac6d75855fd8cd039f88f5f202a4ec8a7 |
| internal/analytics/analytics_test.go | 1682eea431dc85c0309305cb00df809647c12344fc1afa882da8936288809ba3 |
| db/migrations/0041_data_analytics.up.sql | bf98c3f6cceda9883c15c0e72615dc124d0f13e3a8401052a84a3d395e587f9e |
| db/migrations/0041_data_analytics.down.sql | 427c5b8c10fec82410d31f1b335188dbb47fa9e20c3f284ef4ec648b4ca379aa |
| db/tests/0041_data_analytics.test.sql | 2e52a4823898877095946c4467e646c2975ba66499e5595ed4da2f2e8dc36191 |

SHA-256 del pack: `ff2fbb2147c984b9b5eb20323cd76e5406d0d4c12491bde5f5228f29c37c3acb` (21.147 bytes).

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).
- PostgreSQL 18.6 x86-64 (cluster aislado, puerto 55442).

## Verificación ejecutada

- `go test ./internal/analytics/ -count=1`: **5/5 PASS** (métrica con negativos de agregación/dimensión duplicada/TTL cero, freshness fail-closed, reconciliación exacta/tolerada/negativa, ingesta con late-data/recibido-antes-evento/sha inválida y replay key).
- `go test ./... -count=1` (5 paquetes): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 7/7 bloques reproducen byte a byte; `go test` sobre el árbol materializado PASS.
- PostgreSQL 18.6 real (initdb → up → test → down): `ON_ERROR_STOP=1` exit 0. Aserciones DO: replay dedup (`INSERT 0 0` en duplicado), rechazo received-before-event por `check_violation`, métrica versionada + snapshot; down deja `data_ns=true analytics_ns=true`.

## Cobertura de invariantes probada

1. `tenant_id` obligatorio en ingesta y métricas.
2. Definición de métrica versionada (code+version únicos) con agregación acotada.
3. Replay idempotente por (source, external_id, payload sha).
4. `received_at >= event_time`; late-data bound por `MaxLateness`.
5. Freshness fail-closed (TTL no positivo nunca es fresco); reconciliación no pasa con tolerancia negativa.

## Condiciones residuales

- Motor de agregación real sobre el source of truth (SQL de métricas): decisión del proyecto.
- CDC/batch productivo (fuera del outbox documental ya existente): plataforma del proyecto.
- Warehouse/lakehouse y su costo/freshness: `SPEC_ONLY` (se elige por workload).
- Acelerador y perfiles por vertical: pendientes.

V179 cierra las superficies `DATA-INGEST` y `ANALYTICS-BI` a nivel de contrato; no declara plataforma analítica productiva.
