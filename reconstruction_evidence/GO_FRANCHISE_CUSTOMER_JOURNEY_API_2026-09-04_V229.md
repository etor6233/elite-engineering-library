# Go Franchise Customer Journey API — evidencia V229

## Resultado

```yaml
pack: GO-FRANCHISE-CUSTOMER-JOURNEY-API
version: 0.9.0
implementation: REBUILD_VERIFIED
admission: CONDITIONED
verified_at: 2026-09-04
materialization: 33/33 PASS
gofmt_idempotence: 9/9 PASS
postgres: 18.6 real
focused_postgres_runs: 2/2 PASS
full_go_packages: 47 PASS
go_vet: PASS
go_build: PASS
pack_sha256: 8710a8774588dc5565ccc44ae729505a22b8eccf3522a9c8593d6d5c2b5ca53c
```

## Corrección demostrada

Cotización dejó de ser un write no reintentable. HTTP calcula SHA-256 del body exacto y exige `Idempotency-Key`; servicio y repositorio conservan ambos. PostgreSQL usa el owner común `platform.idempotency_record` dentro de la misma transacción que precio server-side, `sales.quotation` y `quotation.issued`.

Las regresiones reales demuestran:

- misma clave + mismo hash devuelve la cotización original;
- misma clave + hash distinto devuelve conflicto;
- dos requests concurrentes producen una creación, un replay, una cotización y un evento;
- un fallo de lead/precio no deja idempotencia parcial;
- la suite completa con PostgreSQL, `go vet` y `go build` pasa.

## Hashes modificados

| Archivo | SHA-256 |
|---|---|
| `internal/franchisejourney/service.go` | `d434e2e32e3bc0c30038cb6443701347421ff68524145bb4de2b15a10c0c70de` |
| `internal/franchisejourney/service_test.go` | `eed8ed48abec452ada7c605b6081b46e751c486b875585e5657156d6fa76dd1d` |
| `internal/platform/httpapi/franchisejourney.go` | `6cff8e3f194c77f35f2fdc95ff3ede779ebb589c824eac3861bb2de4d74184c4` |
| `internal/platform/httpapi/franchisejourney_test.go` | `72729e2560f0d3cb9ba8ea550f1df5d7c37d1414202405e59e5d34eb7f422a2c` |
| `internal/platform/postgres/franchisejourney.go` | `2e6adfc79317b01b77f7e78340109cae56a3d3acccd5397fc32d183d3e0bd5be` |
| `internal/platform/postgres/franchisejourney_integration_test.go` | `8ed981f86591c73e6d906bbf53440508d76bd337776277ce3fbcdb13bfbe72ee` |

## Procedencia y límite

Los cambios son `AUTHORED`; no se atribuyen a PostgreSQL, Microsoft ni Google. Aplican el contrato de transacción/bloqueo ya gobernado por las fuentes oficiales fijadas en el pack. La idempotencia de creación queda demostrada; identidad live, términos comerciales, fiscalidad y producción siguen siendo condiciones del proyecto.
