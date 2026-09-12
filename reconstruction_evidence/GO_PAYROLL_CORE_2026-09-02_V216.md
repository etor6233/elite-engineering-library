# Go Payroll Core — V216

## Resultado estrecho

V216 materializa `GO-PAYROLL-CORE 0.1.0`: cálculo de nómina tenant-scoped config-driven (deducciones fijas o porcentuales inyectadas en runtime), net fail-closed e idempotente. Sin tablas impositivas hardcodeadas.

## Archivos materializados (2)

| Archivo | SHA-256 |
|---|---|
| internal/payroll/payroll.go | 8b9a284052d930f50075d3644a0d27426bf24267c80bd8ab14f100198f18b690 |
| internal/payroll/payroll_test.go | 7862084eebcbdaa1fd2a27f1d1d7e51b5067a67320a55eb2a1f4cd506c839809 |

SHA-256 del pack: `b23d1ca807194768b2bf095bde037e186338c07dc528fe83993b8d7d30e3804e`.

## Verificación

- `go test ./internal/payroll/ -count=1`: 5/5 PASS.
- `go test ./... -count=1` (37 paquetes): PASS. `go vet`: exit 0.
- Round-trip: 2/2 bloques reproducen byte a byte.

## Condición residual

Las reglas impositivas/legales por jurisdicción se inyectan por run; no se inventan tablas.
