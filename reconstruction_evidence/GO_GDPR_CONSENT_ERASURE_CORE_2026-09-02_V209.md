# Go GDPR Consent & Erasure Core — V209

## Resultado estrecho

V209 materializa `GO-GDPR-CONSENT-ERASURE-CORE 0.1.0`, el gobierno GDPR de consentimiento (Art. 6/7) y derecho de supresión (Art. 17) con tenant-scoping y fail-closed.

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — `AUTHORED`, gobernado por GDPR Art. 6/7/17 (https://gdpr-info.eu).

## Archivos materializados (2)

| Archivo | SHA-256 |
|---|---|
| internal/gdpr/gdpr.go | 295da7e159b6de92dfaf25498802fdfcc5a0366c1ee47a8ecd9e431a373c6da8 |
| internal/gdpr/gdpr_test.go | 722b7ddd9debd48757793ea10129c1e88cee7b5ca851611baf67d269bcdbbec5 |

SHA-256 del pack: `80d5971afc600c3fdf570c37156d2cb302ecd3a9e07d35e0999412ff64a24bf3`.

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).

## Verificación ejecutada

- `go test ./internal/gdpr/ -count=1`: **7/7 PASS**.
- `go test ./... -count=1` (31 paquetes): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 2/2 bloques reproducen byte a byte.

## Cobertura de invariantes probada

1. Consentimiento específico e informado (finalidad + versión de política).
2. Retirada tan fácil como otorgarlo; re-otorgamiento.
3. Borrado happy path requested→in_progress→completed.
4. Retención legal activa bloquea el borrado; al liberar, completa.
5. Rechazo exige base legal.
6. Transición inválida y duplicado rechazados.
7. Aislamiento de tenant.

## Condiciones residuales

- La ejecución real del borrado en los almacenes de dominio (cascada a pedidos/mensajes/copias) es composición del proyecto.
