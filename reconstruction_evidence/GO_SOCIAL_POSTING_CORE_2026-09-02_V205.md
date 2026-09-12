# Go Social Posting Core — V205

## Resultado estrecho

V205 materializa `GO-SOCIAL-POSTING-CORE 0.1.0`, la publicación programada a redes (Instagram/Facebook/TikTok/LinkedIn): schedule, estados scheduled→posted/failed y fail-closed, tenant-scoped.

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — `AUTHORED`.

## Archivos materializados (2)

| Archivo | SHA-256 |
|---|---|
| internal/social/social.go | 9955ff37e1c971974fbe54eaff0c502cc839cb98b39fc1fd011342a9fe2d3ede |
| internal/social/social_test.go | 3ea62bc54249337cffd654964894e2a6e035596b49fb8a0c19d26e28a97ff8cc |

SHA-256 del pack: `04791bba47e6173fb72098c8f690f8ffa62b7c2535ecd677da40b77e71e15330`.

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).

## Verificación ejecutada

- `go test ./internal/social/ -count=1`: **5/5 PASS** (schedule/due/posted, failed, pasado/inválido, duplicado, aislamiento tenant).
- `go test ./... -count=1` (27 paquetes): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 2/2 bloques reproducen byte a byte.

## Cobertura de invariantes probada

1. Canal admitido, texto no vacío.
2. Schedule futuro.
3. Estados scheduled→posted/failed.
4. tenant-scoped.

## Condiciones residuales

- Publicación real por red (adapter Meta Graph/otros): composición (CONDITIONED).

V205 cierra el hueco de social posting; el adapter de red real es composición.
