# Go Agent Domain Binding — V188

## Resultado estrecho

V188 materializa `GO-AGENT-DOMAIN-BINDING 0.1.0`, el binding del agente conversacional a los endpoints reales del backend: cita (POST público), cotización (POST `/v1/franchise/quotes`), estado de pedido (GET `/v1/customer/orders`) y devolución (rechazo de entrega vía `/v1/customer/handovers/{id}/reject`). Los mapeos de negocio se inyectan por `Config`; valores desconocidos fallan cerrado.

La devolución no ejecuta dinero/inventario: rechaza la entrega y deja el efecto gobernado por la cola de aprobación y los workers de devolución.

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — `AUTHORED` sobre el contrato `agenttools.Domain` y los endpoints reales.

## Archivos materializados (3)

| Archivo | SHA-256 |
|---|---|
| internal/domainbind/config.go | 69d8e37e46240016e45cac28c1fdd2098cf7fe2ed1018e9253ae15e58e8ddf44 |
| internal/domainbind/gateway.go | b2f2ac5e8309fc091de270560161747b81d6e4dbfe2e2ba56f8ddbd6980f4cb7 |
| internal/domainbind/gateway_test.go | 30fd310aa3982f52623bdc6394fcfffab8ef09c5c843aa0b8a1f5a4e6d541a31 |

SHA-256 del pack: `b95024fa9224ce2e3220be0f54528e1e68c2e9814c0c1b89c5a96d87e2e8c242`.

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).

## Verificación ejecutada

- `go test ./internal/domainbind/ -count=1`: **6/6 PASS** (config y resolución fail-closed, cita con servidor falso verificando path/kind/lead, servicio desconocido, cotización con mapeo variant/price_book, estado de pedido, devolución vía rechazo de entrega).
- `go test ./... -count=1` (12 paquetes): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 3/3 bloques reproducen byte a byte; `go test` sobre el árbol materializado (domainbind + agenttools + agent + aifoundation + search) PASS.

## Cobertura de invariantes probada

1. Servicio/producto desconocido → error (fail-closed).
2. Lead ausente → error.
3. `when` debe ser futuro (ISO o hoy/mañana).
4. La devolución rechaza la entrega; no ejecuta efectos.

## Condiciones residuales

- Mapeos de negocio reales (service→kind, product→variant, price book): config del proyecto.
- Canal SMS (Twilio-compatible): pendiente (mismo patrón HTTP).
- Gate pnpm del frontend: entorno.

V188 conecta el chatbot a las operaciones reales; los mapeos de negocio se inyectan al arrancar el proyecto.
