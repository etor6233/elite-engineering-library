# Go App Wiring — evidencia de corrección V231

## Resultado

```yaml
pack: GO-APP-WIRING
version: 0.1.1
implementation: RECONSTRUCTIBLE
admission: CONDITIONED
verified_at: 2026-09-04
tests: PASS
materialization: 3/3 PASS
gofmt_idempotence: 3/3 PASS
full_go_packages: 47 PASS
go_vet: PASS
go_build: PASS
pack_sha256: c5552921231139321b15892ca22578a5f12e7a8d880e5143de8cd6ed841cdb74
```

## Alcance exacto

La fábrica ahora exige `TENANT_ID` y un `DOMAIN_TOKEN_PROVIDER`, y los entrega al gateway 0.2.0. Esto corrige el ensamblado de identidad del borde protegido, pero no cierra `LIB-FAIL-1743`: `Agent.Process` todavía no usa Responses, economía, aprobación, canal ni persistencia durable, y sus tools heredadas no poseen identidad de comando.

## Hashes

| Archivo | SHA-256 |
|---|---|
| `internal/app/config.go` | `28737f8c5505526ea365fc5699cffbe6852a5dec48870c096cfe19c1a484d26a` |
| `internal/app/app.go` | `0c2c037c027f57c80f60bd389255f2793ebcb471d3f8ee7c7cd4a26d821fe6b2` |
| `internal/app/app_test.go` | `10bc05aecc541286408dc004f0f227a8999821b9a73de04578d8135b54001fc7` |

## Verdad retenida

No se promueve a runtime completo. La compilación y la construcción de objetos son evidencia necesaria, no evidencia E2E.
