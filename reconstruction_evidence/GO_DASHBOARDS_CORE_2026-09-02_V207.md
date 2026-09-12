# Go Dashboards Core — V207

## Resultado estrecho

V207 materializa `GO-DASHBOARDS-CORE 0.1.0`, el snapshot de KPIs por franquicia (ventas, pedidos, puntos, reseñas, NPS, citas) con agregación tenant-scoped y fail-closed.

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — `AUTHORED`.

## Archivos materializados (2)

| Archivo | SHA-256 |
|---|---|
| internal/dashboards/dashboards.go | e041541596d98e2749653bbad7177593ed6fc43bd0c4f1bff5c018e813aa930f |
| internal/dashboards/dashboards_test.go | 463bc8d65d131b6bf75f8fd56992d5d40e4bd6645fdfe4f3338a5c9207a6ca0c |

SHA-256 del pack: `6ef647ef1e3084a717a6f00521663a64c96ffa70caa37cccad94cffd1f418085`.

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).

## Verificación ejecutada

- `go test ./internal/dashboards/ -count=1`: **3/3 PASS** (snapshot con reviews/NPS, vacíos, inválido).
- `go test ./... -count=1` (29 paquetes): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 2/2 bloques reproducen byte a byte.

## Cobertura de invariantes probada

1. Tenant y periodo válidos.
2. Conteos ≥ 0.
3. Reviews 1-5, NPS 0-10.
4. Agregación pura y fail-closed.

## Condiciones residuales

- Render visual (gráficos) en el frontend: composición del proyecto.

V207 cierra el último hueco de negocio; el render visual es composición.
