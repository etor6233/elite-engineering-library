# Go Promotions Core — V195

## Resultado estrecho

V195 materializa `GO-PROMOTIONS-CORE 0.1.0`, los cupones/descuentos tenant-scoped: percent o fixed, con expiración, monto mínimo, límite de uso y redención fail-closed.

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — `AUTHORED`.

## Archivos materializados (2)

| Archivo | SHA-256 |
|---|---|
| internal/promotions/promotions.go | 072dd84ef5e882621aec66ab483c5783c946c26aae0a608a2c0aa0b93e873e77 |
| internal/promotions/promotions_test.go | b300559e2a9abf44b4431715a64a64119ec641b272b48be8d9f216b277d61883 |

SHA-256 del pack: `ef721b1c53b55fbc39ee251a6e15184350ca01bf0dfef105aa2bc043eea91780`.

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).

## Verificación ejecutada

- `go test ./internal/promotions/ -count=1`: **6/6 PASS** (percent, fixed con tope, límite agotado, expirado/inactivo/mínimo, valores inválidos, aislamiento tenant).
- `go test ./... -count=1` (17 paquetes): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 2/2 bloques reproducen byte a byte.

## Cobertura de invariantes probada

1. percent ≤ 10000 basis points, fixed > 0.
2. Expiración obligatoria.
3. uso ≤ MaxUses.
4. pedido ≥ mínimo.
5. tenant-scoped (cross-tenant → ErrNotFound).

## Condiciones residuales

- Aplicación del descuento al pedido final: composición del dominio (commerce).

V195 cierra el hueco de cupones; la aplicación al precio final es del proyecto.
