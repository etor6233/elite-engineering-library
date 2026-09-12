# Go Loyalty Core — V194

## Resultado estrecho

V194 materializa `GO-LOYALTY-CORE 0.1.0`, el ledger de puntos de fidelización: earn/burn idempotente, historia inmutable, tenant-scoped y fail-closed (saldo nunca negativo, nunca sobre-canje).

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — `AUTHORED` sobre el patrón estándar de ledger de fidelización.

## Archivos materializados (2)

| Archivo | SHA-256 |
|---|---|
| internal/loyalty/loyalty.go | a80dc8016ba9af64ec1e604ab3386c50fba7bf79b03e5632f6867a10a34e1c9c |
| internal/loyalty/loyalty_test.go | 4208c298902f7e5413d099c82d647d457f3cd4506311ac0fe36078064ba89fa3 |

SHA-256 del pack: `b2f3e07c4b127be207351b82d77be2a977c708b0aaed9aa200ed2103699cbcc8`.

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).

## Verificación ejecutada

- `go test ./internal/loyalty/ -count=1`: **6/6 PASS** (earn/burn balance, burn insuficiente fail-closed, duplicado rechazado, puntos inválidos, aislamiento tenant, historia inmutable ordenada).
- `go test ./... -count=1` (16 paquetes): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 2/2 bloques reproducen byte a byte.

## Cobertura de invariantes probada

1. points > 0.
2. burn ≤ balance (fail-closed).
3. Idempotente por entry id.
4. Historia inmutable y tenant-scoped (sin fuga cross-tenant).

## Condiciones residuales

- Earn automático en compra (binding al pedido): composición del proyecto.
- Canje por recompensas (catálogo de premios): composición.

V194 cierra el hueco de fidelización; el binding al pedido es del proyecto.
