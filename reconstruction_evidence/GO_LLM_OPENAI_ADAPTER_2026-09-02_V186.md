# Go OpenAI-Compatible LLM Adapter — V186

## Resultado estrecho

V186 materializa `GO-LLM-OPENAI-ADAPTER 0.1.0`, el adapter que implementa `aifoundation.LLMProvider` sobre un endpoint Chat Completions OpenAI-compatible. La API key se inyecta del entorno en runtime (nunca en código); los tests usan `httptest` (servidor HTTP falso) y no requieren red ni cuenta.

Demuestra el principio correcto: **el código del adapter se escribe ahora; las credenciales se inyectan cuando arrancás el proyecto.**

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — `AUTHORED` sobre el contrato `aifoundation.LLMProvider`.
- Referencia del protocolo: OpenAI Chat Completions API (`https://platform.openai.com/docs/api-reference/chat`).

## Archivos materializados (3)

| Archivo | SHA-256 |
|---|---|
| internal/llmopenai/config.go | e361c4e7cb32ecd28da64843dec667b254acc700babf120eb9e8d46067a0d28e |
| internal/llmopenai/adapter.go | 32f42b87994db1a225decc753989e543a5d38a588b69be1cc4fecca36c77be72 |
| internal/llmopenai/adapter_test.go | b01956c6264ce79144b99128bbddf1799f0e6bfa5c5f283f3102a1ac397f65ab |

SHA-256 del pack: `98e99f787b764773d58e85e33bc45b76f0f96ee2912a963643723ebd0d758bb7`.

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).

## Verificación ejecutada

- `go test ./internal/llmopenai/ -count=1`: **4/4 PASS** (config válida/negativos, round-trip con servidor falso verificando path/Bearer/modelo/`response_format`, fail-closed en 401, rechazo de contenido vacío).
- `go test ./... -count=1` (10 paquetes): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 3/3 bloques reproducen byte a byte; `go test` sobre el árbol materializado (llmopenai + aifoundation) PASS.

## Cobertura de invariantes probada

1. URL absoluta http(s) y modelo exacto (no `latest`/`auto`).
2. API key sólo del entorno (Bearer header); nunca en código.
3. Status no-2xx, contenido vacío o payload inesperado → error (fail-closed).
4. `response_format: json_object` sólo si hay schema.

## Condiciones residuales

- Cuenta/provider real y su costo: runtime del proyecto (la key se inyecta ahí).
- Tiering por proveedor más barato: composición sobre `GO-LLM-ECONOMY-CORE`.
- Adapters de canales (Meta Graph/SMTP/SMS) y binding del `Domain`: pendientes.

V186 ofrece el adapter provider, pero por sí solo no conecta el agente al LLM. V228 agrega tool calling nativo Responses; el runtime integral continúa gobernado por `LIB-FAIL-1743` hasta demostrar el recorrido completo. Ninguno ejecuta contra cuenta sin credenciales.

## Revalidación 2026-09-04

`adapter.go` y `config.go` fueron normalizados con `gofmt` sin cambio conductual. Una materialización fresca 3/3 quedó `gofmt`-estable y la composición completa posterior pasó 47 paquetes, `go vet ./...` y `go build ./...`.
