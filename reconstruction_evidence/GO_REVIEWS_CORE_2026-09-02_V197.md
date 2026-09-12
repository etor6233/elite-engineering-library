# Go Reviews Core — V197

## Resultado estrecho

V197 materializa `GO-REVIEWS-CORE 0.1.0`, las reseñas/ratings tenant-scoped con moderación: rating 1-5, una por autor por sujeto, estados pending→approved/rejected y promedio sólo sobre aprobadas.

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — `AUTHORED`.

## Archivos materializados (2)

| Archivo | SHA-256 |
|---|---|
| internal/reviews/reviews.go | cb2ed137af85eef222159c3cf3c5f0eadd027a179c7e09603258fa065f56d208 |
| internal/reviews/reviews_test.go | f07b1aa8d7017346fd870ddc3966e579d69e176bc1d5099915352dff9f7033ca |

SHA-256 del pack: `80456619b92b15ee479b57fe760c81d2cb5d5d280be3ae9e7ccf6fa599c16b6e`.

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).

## Verificación ejecutada

- `go test ./internal/reviews/ -count=1`: **5/5 PASS** (moderación y promedio, rating inválido, duplicado, transición inválida, aislamiento tenant).
- `go test ./... -count=1` (19 paquetes): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 2/2 bloques reproducen byte a byte.

## Cobertura de invariantes probada

1. rating 1-5.
2. Una reseña por autor por sujeto.
3. Promedio sólo sobre aprobadas (no-aprobadas no cuentan).
4. tenant-scoped.

## Condiciones residuales

- Publicación en el portal de producto: composición del proyecto.

V197 cierra el hueco de reseñas; la UI de publicación es del proyecto.
