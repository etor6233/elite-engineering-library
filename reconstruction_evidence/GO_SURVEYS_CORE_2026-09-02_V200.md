# Go Surveys Core — V200

## Resultado estrecho

V200 materializa `GO-SURVEYS-CORE 0.1.0`, las encuestas/NPS: respuestas con score 0-10, una por cliente por encuesta, tenant-scoped y NPS estándar (promotores 9-10, detractores 0-6, pasivos 7-8).

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — `AUTHORED`.

## Archivos materializados (2)

| Archivo | SHA-256 |
|---|---|
| internal/surveys/surveys.go | 21f086b106402edd388879711a49ffb95311b86c17de2650e9e76a7086a6c989 |
| internal/surveys/surveys_test.go | 25bf6ed15fe3a251b4e1a5905755dfe83759e4dfb1a55a68aee31b81468602b0 |

SHA-256 del pack: `e5786433b98612a5ff24b86336d610fe8818ce98dc9d84bb877bfd9a31565b9e`.

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).

## Verificación ejecutada

- `go test ./internal/surveys/ -count=1`: **5/5 PASS** (NPS cálculo, vacío 0, score inválido, duplicado, aislamiento tenant).
- `go test ./... -count=1` (22 paquetes): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 2/2 bloques reproducen byte a byte.

## Cobertura de invariantes probada

1. score 0-10.
2. Una respuesta por cliente por encuesta.
3. NPS = (promotores - detractores) / total * 100.
4. tenant-scoped.

## Condiciones residuales

- Envío de encuesta post-compra (trigger): composición del proyecto.

V200 cierra el hueco de encuestas/NPS; el trigger de envío es composición.
