# Go Gift Cards Core — V198

## Resultado estrecho

V198 materializa `GO-GIFT-CARDS-CORE 0.1.0`, el crédito de tienda (gift cards): issue/redeem idempotente, saldo nunca negativo, tenant-scoped y fail-closed.

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — `AUTHORED`.

## Archivos materializados (2)

| Archivo | SHA-256 |
|---|---|
| internal/giftcards/giftcards.go | e6339a5da5e58c5232d726e1c1d5fad3e380d1ea91e0b94df7d1a8b8969b9d29 |
| internal/giftcards/giftcards_test.go | 66eba48567fcb54d8a3c4161919d1bb024fbe5e2cae74077caf73d752bf0661a |

SHA-256 del pack: `0ce8666f978d918c076652ec11a99fb1ba2e0477ac44ab4ce3c96b8b0490ad2d`.

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).

## Verificación ejecutada

- `go test ./internal/giftcards/ -count=1`: **5/5 PASS** (issue/redeem balance, insuficiente fail-closed, idempotencia, duplicado/inválido, aislamiento tenant).
- `go test ./... -count=1` (20 paquetes): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 2/2 bloques reproducen byte a byte.

## Cobertura de invariantes probada

1. amount > 0.
2. redeem ≤ balance (fail-closed).
3. Idempotente por redeem id.
4. tenant-scoped.

## Condiciones residuales

- Pago con gift card en el checkout: composición del dominio (commerce).

V198 cierra el hueco de gift cards; el pago en checkout es del proyecto.
