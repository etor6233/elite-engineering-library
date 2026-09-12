# Go Help Center Core — V203

## Resultado estrecho

V203 materializa `GO-HELP-CENTER-CORE 0.1.0`, la base de conocimiento autoservicio: artículos con categorías, versionado (concurrencia optimista) y estados draft→published→archived, tenant-scoped. La búsqueda delega en `GO-SEARCH-CORE`.

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — `AUTHORED`.

## Archivos materializados (2)

| Archivo | SHA-256 |
|---|---|
| internal/helpcenter/helpcenter.go | 9907c8d98702558fd874b9187c119d3306efce5ab2a89ee556ee90c0618684e5 |
| internal/helpcenter/helpcenter_test.go | a433cb6b34a351ca6ec0501a36827992b5d0829c8145c173d946fd9daba538c8 |

SHA-256 del pack: `20fc7f9ab1ea1e41102c7fc16bd1a352b95df2b561194d1ff9abbab343b11216`.

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).

## Verificación ejecutada

- `go test ./internal/helpcenter/ -count=1`: **5/5 PASS** (estados, concurrencia optimista, transición inválida, búsqueda sólo publicados, inválido y aislamiento tenant).
- `go test ./... -count=1` (25 paquetes): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 2/2 bloques reproducen byte a byte.

## Cobertura de invariantes probada

1. Título/cuerpo no vacíos, categoría segura.
2. Versionado con concurrencia optimista.
3. Estados draft→published→archived.
4. Búsqueda sólo publicados.
5. tenant-scoped.

## Condiciones residuales

- Indexado full-text en `GO-SEARCH-CORE`: composición del proyecto.

V203 cierra el hueco de help-center; el indexado full-text es composición.
