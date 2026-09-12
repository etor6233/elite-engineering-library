# Reconstruction Evidence — Go Enterprise Backend V1

```yaml
evidence_id: "GO-EB-20260824-V1"
pack_sha256: "db390900cc2e3227be7ea72d061cfb4b2e4bdd9cdb41d69652e0c00cdb0a1e1e"
materialized_file_count: 10
ordered_file_hash_aggregate_sha256: "9eb91ddbf78198a36e4e46aa4103e8284acfcb391eabf9b08129af493f4f09bc"
toolchain: "go1.26.7 windows/amd64"
toolchain_archive_sha256: "f4f534a486e4bc3387fa18f08208f2f854b7aaea8a08f2a2d829a914a05abb11"
result: "PASS"
```

El target comenzó vacío y recibió los 10 archivos exclusivamente desde `GO_ENTERPRISE_BACKEND_CORE.md`. Cada hash fue verificado antes de escribir. `go mod download`, `go test ./...`, `go vet ./...` y `go build ./cmd/api` aprobaron.

Cobertura: dominio de pedidos y state machine, permisos, idempotencia, servicio de aplicación, adapter PostgreSQL transaccional con optimistic concurrency/outbox, contrato HTTP estricto, límites, problem JSON, timeouts y graceful shutdown.

No se simularon PostgreSQL, OIDC ni deployment reales; permanecen conditions explícitas de composición productiva.

