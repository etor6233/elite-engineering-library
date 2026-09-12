# Go Agent Domain Binding — evidencia V233

## Resultado

```yaml
pack: GO-AGENT-DOMAIN-BINDING
version: 0.3.0
implementation: REBUILD_VERIFIED
admission: CONDITIONED
verified_at: 2026-09-04
materialization: 3/3 PASS
gofmt_idempotence: 3/3 PASS
full_go_packages: 47 PASS
pack_sha256: 1e78d38a2bc9ce315304a882833716d82a7358cb038254f26cb45ea2163d06e3
```

## Corrección multi-contacto

`Scope` lleva organización y lead resueltos fuera del modelo. Los métodos `BookAppointmentFor`, `CreateQuoteFor`, `OrderStatusFor` y `RequestReturnFor` usan ese scope en paths/bodies; tenant cruzado o scope incompleto falla antes de I/O. Una regresión ejecuta dos contactos secuenciales y observa `north/lead-a` y `south/lead-b` sin contaminación.

Se conservan la identidad estable de comando y el proveedor rotatorio de Bearer demostrados en V230. Los defaults estáticos siguen disponibles sólo para composición single-scope condicionada; el runtime omnicanal debe usar los métodos `*For`.

## Hashes

| Archivo | SHA-256 |
|---|---|
| `internal/domainbind/config.go` | `ee62027dac5f728385c0d8314d4f7d81d751a6568e9df8f782a3783f141f5a8c` |
| `internal/domainbind/gateway.go` | `7829c722655985d90932b9dd7fa3bf0098fd8db2ae0ce33fd1c61a7fe5418c23` |
| `internal/domainbind/gateway_test.go` | `2b8255e283880c774fbd383acc00dd2ef250b9fb433889e00ac40123d525081c` |

## Límite

El resolver que vincula identidad de proveedor con tenant/org/lead es una dependencia obligatoria del runtime y debe probar autorización real por proyecto.
