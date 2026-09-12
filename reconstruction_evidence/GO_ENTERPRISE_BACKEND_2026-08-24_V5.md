# Go Enterprise Backend — Reconstruction Evidence V5

```yaml
evidence_id: "GO-EB-20260824-V5"
pack_id: "GO-ENTERPRISE-BACKEND"
pack_version: "0.3.0"
pack_sha256: "bace6b83561ddfcd1fd121a58393f605414bc415e1ca97600bcebd9af099748a"
materialized_file_count: 17
ordered_file_hash_aggregate_sha256: "d32039be7bf8fda6391b89f874afe003a14a43cfd257820f43162e443b3fc7fe"
toolchain: "go1.26.7 windows/amd64"
database: "PostgreSQL 18.6 windows/amd64"
result: "PASS_CONDITIONED"
verified_at: "2026-08-24"
```

## Gates ejecutados

```text
materialize_markdown_pack.ps1 PASS (17 hashes)
gofmt -d                      PASS
go test -count=1 ./...        PASS
go vet ./...                  PASS
go build ./cmd/api            PASS
```

Los tests PostgreSQL de V4 continúan pasando. V5 añadió dos eventos listos y demostró que workers distintos reciben claims disjuntos mediante `FOR UPDATE SKIP LOCKED`, que un worker ajeno no puede marcar el evento, que el owner puede confirmarlo y que release permite reclaim controlado. El claim incrementa attempts y usa lease vencible; publish externo sigue siendo al menos una vez.

## Límite del claim

El pack aporta la primitiva durable, no un provider publisher completo. Faltan retry budget/jitter por proveedor, dead-letter policy, backpressure, inbox/jobs workers, pruebas de crash en cada frontera, carreras de idempotencia, issuer real, roles mínimos, load/security y recovery. Continúa `CONDITIONED`.

