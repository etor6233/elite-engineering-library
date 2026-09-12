# Go Enterprise Backend — Reconstruction Evidence V3

```yaml
evidence_id: "GO-EB-20260824-V3"
pack_id: "GO-ENTERPRISE-BACKEND"
pack_version: "0.2.0"
pack_sha256: "365caac81527e971d3d6e4f022a162f43ba7f8ebdded77d38126b17c5ffd5e48"
materialized_file_count: 14
ordered_file_hash_aggregate_sha256: "0af19ee5bd22f969b56dae923f46f4c7278180efa8914ec893cf5ca88e45d0b0"
toolchain: "go1.26.7 windows/amd64"
toolchain_archive_sha256: "f4f534a486e4bc3387fa18f08208f2f854b7aaea8a08f2a2d829a914a05abb11"
result: "PASS_CONDITIONED"
verified_at: "2026-08-24"
```

## Gates ejecutados

```text
materialize_markdown_pack.ps1 PASS (14 hashes)
gofmt -d                      PASS (sin diff)
go mod download               PASS
go test ./...                 PASS
go vet ./...                  PASS
go build ./cmd/api            PASS
```

El archive oficial de Go fue comprobado contra el SHA-256 publicado antes de usarlo. `go.mod` y `go.sum` fijan pgx 5.10.0, coreos/go-oidc 3.20.0 y sus dependencias transitivas.

## Cambio de seguridad respecto de V2

- el arranque exige `OIDC_ISSUER` y `OIDC_AUDIENCE` y falla cerrado si discovery no funciona;
- `go-oidc` valida issuer, audience, expiración, firma y algoritmos RS256 admitidos;
- `tenant_id`, subject y permissions provienen del token verificado;
- crear una orden exige `order:create`; tenant y cliente dejaron de aceptarse desde el body;
- el contrato HTTP prueba el rechazo de requests sin bearer.

## Límite del claim

Esta evidencia prueba reconstrucción, formato, dependencias, tests, vet y build. No prueba un issuer real, rotación/revocación de JWKS, policy migration, PostgreSQL, concurrencia, roles mínimos, race/load/security, backup/restore ni deployment. Las migrations SQL se reconstruyen, pero deben ejecutarse junto a `PG-TX-FOUNDATION` sobre PostgreSQL 18.6 antes de autorizar persistencia.

