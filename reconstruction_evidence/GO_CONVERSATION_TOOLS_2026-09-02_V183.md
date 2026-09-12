# Go Conversation Tools — V183

## Resultado estrecho

V183 materializa `GO-CONVERSATION-TOOLS 0.1.0`, el contrato de negocio agnóstico del chatbot (citas, cotización/venta, estado de pedido, devolución) y las tools que lo bindean al agente por intención. El bot delega en el `Domain`; nunca toca dinero/inventario directamente.

Hace al chatbot adaptable a cualquier negocio: cualquier franquicia implementa el contrato `Domain` (cuatro operaciones) y el agente queda conectado. El parser `clave: valor` es el camino determinista; convertir texto libre a estructura es trabajo del LLM (CONDITIONED).

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — core `AUTHORED`.
- `internal/agent 0.1.0` (intents/tools/registry) — de este repositorio.

## Archivos materializados (5)

| Archivo | SHA-256 |
|---|---|
| internal/agenttools/domain.go | b4e134e2843f3202f22a5c66531f0dc69db729445400490c9640086f4f264f0c |
| internal/agenttools/parse.go | 12fb9fc6aefd5cbaf884683bf543c310330cc00ba5bec50f408af8da8c8fa3a2 |
| internal/agenttools/tools.go | 4eca0fc21c16602dae6f778b0ca78d280ed310d158ce2b7b6b458cbfefacea7a |
| internal/agenttools/toolkit.go | aa1a9d9c53e8c914a9ca6089bdebe1844ffc93dd3dcb13eb815f38de88cb4b43 |
| internal/agenttools/agenttools_test.go | 3dac1f760ff224e33afe826b823426892a29cf849a2ce90d58b69b9eebad439e |

SHA-256 del pack: `ad58fda99bd5ca2ff13eac0491f20608848d67f25330d06c52893e63cee7ad57`.

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).

## Verificación ejecutada

- `go test ./internal/agenttools/ -count=1`: **6/6 PASS** (toolkit con 4 intents, extracción y llamada de cita, `ErrNeedsInfo` por campo ausente, cantidad parseada e inválida, estado de pedido y devolución).
- `go test ./... -count=1` (8 paquetes): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 5/5 bloques reproducen byte a byte; `go test` sobre el árbol materializado (agent + agenttools + aifoundation + search) PASS.

## Cobertura de invariantes probada

1. Una intención → una tool → un método del `Domain`.
2. Campo requerido ausente → `ErrNeedsInfo` (nunca se inventa el dato).
3. Cantidad inválida → `ErrNeedsInfo`.
4. El bot delega, no ejecuta efectos de dinero/inventario.

## Condiciones residuales

- Implementación real del `Domain` con los endpoints Go existentes (agenda→journey, venta→commerce, devoluciones→return workers): composición del proyecto.
- Conversión texto-libre → estructura por LLM: CONDITIONED.
- Canales y frontend multi-rol: pendientes.

V183 conecta el chatbot al dominio de negocio de forma agnóstica; no ejecuta efectos productivos.
