# Go Enterprise Backend — Reconstruction Evidence V2

```yaml
evidence_id: "GO-EB-20260824-V2"
pack_id: "GO-ENTERPRISE-BACKEND"
pack_version: "0.1.0"
pack_sha256: "d8271e9e702890ad8c5233955c5291a957b97055af0cb983368edf6a94aa8922"
materialized_file_count: 13
ordered_file_hash_aggregate_sha256: "3191e36e1ae42e15f62413e032be190a3938ac8de929ad0efe62c9ab2f933c4d"
toolchain: "go1.26.7 windows/amd64"
toolchain_archive_sha256: "f4f534a486e4bc3387fa18f08208f2f854b7aaea8a08f2a2d829a914a05abb11"
result: "PASS_CONDITIONED"
verified_at: "2026-08-24"
```

## Reconstrucción

El pack se materializó en un directorio temporal inicialmente vacío mediante `materialize_markdown_pack.ps1`. El materializador verificó el SHA-256 declarado de los 13 archivos antes de escribirlos. El agregado se calculó sobre las rutas relativas y hashes ordenados de esos 13 archivos; el binario generado por `go build` no forma parte del agregado.

## Gates ejecutados

```text
gofmt -d <todos los archivos .go>  PASS (sin diff)
go mod download                    PASS
go test ./...                      PASS
go vet ./...                       PASS
go build ./cmd/api                 PASS
```

El archivo oficial `go1.26.7.windows-amd64.zip` se descargó desde `go.dev/dl` y se rechazaba si su digest no coincidía con el publicado por Go.

## Diferencia material respecto de V1

- dominio, servicio y repositorio propagan `tenant_id`;
- idempotencia usa `request_sha256_hex` y la clave compuesta de `PG-TX-FOUNDATION`;
- outbox aporta `tenant_id`, `occurred_at` y versiones positivas;
- lecturas y optimistic updates filtran por tenant y organización;
- se materializan migrations y tests SQL para organización y órdenes.

## Límite del claim

Esta evidencia autoriza reconstrucción, formato, dependencias, unit/HTTP contract, vet y build. No autoriza persistencia ni producción: las migrations todavía deben ejecutarse con `PG-TX-FOUNDATION` sobre PostgreSQL 18.6; falta identidad OIDC/policy integrada que derive el tenant de credenciales verificadas, además de tests de concurrencia, roles mínimos, race/load/security, backup/restore y deployment.

