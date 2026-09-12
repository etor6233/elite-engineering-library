# Go Enterprise Backend — Reconstruction Evidence V4

```yaml
evidence_id: "GO-EB-20260824-V4"
pack_id: "GO-ENTERPRISE-BACKEND"
pack_version: "0.2.1"
pack_sha256: "5b54d24cea9845ee8c7ed44a6d16b550c3697fc7b26c37b73225dc03eb0a9efe"
materialized_file_count: 15
ordered_file_hash_aggregate_sha256: "5fbc6f473c76724d2468341a2bf6babdb5e761908d5a006144194a76aa0c9c26"
toolchain: "go1.26.7 windows/amd64"
database: "PostgreSQL 18.6 windows/amd64"
result: "PASS_CONDITIONED"
verified_at: "2026-08-24"
```

## Gates ejecutados

```text
materialize_markdown_pack.ps1                   PASS (15 hashes)
gofmt -d                                        PASS (sin diff)
go mod download                                 PASS
go test -count=1 ./...                          PASS
go vet ./...                                    PASS
go build ./cmd/api                              PASS
PostgreSQL foundation/order up-test-down-up-test PASS
```

`TEST_DATABASE_URL` activó el test de integración, que creó fixtures multi-tenant aislados, comprobó create, replay idempotente, conflicto por reutilizar key con hash distinto, lectura con tenant/organización y transición optimista con outbox. El test limpia sus filas al inicio y al final.

## Seguridad cubierta

La ruta de creación rechaza ausencia de bearer, deriva tenant/subject del principal verificado y exige `order:create`. El adapter OIDC fija issuer, audience y RS256 a través de coreos/go-oidc 3.20.0.

## Límite del claim

Quedan fuera: issuer/JWKS real y rotación/revocación, pruebas de autorización por recurso adicionales, carreras multi-session, race/load/abuse, roles DB mínimos, RLS opcional, backup/PITR/restore, packaging y deployment. Por eso continúa `CONDITIONED`, aunque persistencia real ya fue demostrada.

