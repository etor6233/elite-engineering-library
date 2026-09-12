# Go Conversational Agent — V178

## Resultado estrecho

V178 materializa `GO-CONVERSATIONAL-AGENT 0.1.0`, el núcleo determinista del agente conversacional 24/7 (superficie `RAG-AGENTS`): ruteo de intención acotado, autorización de tools por intención, recuperación de conocimiento tenant-scoped, guardrails fail-closed y handoff humano con evidencia append-only.

Cubre atención al cliente, citas/reservas y ventas a nivel de núcleo: el agente decide, autoriza, recupera y deriva. No incorpora un LLM, un canal ni reglas de negocio: el ruteo es determinista y el LLM, el canal (web/WhatsApp) y las tools de dominio son adapters `CONDITIONED`. El dinero, el inventario y la agenda nunca los decide el agente; delega en los endpoints Go existentes vía el contrato `Tool`.

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — núcleo `AUTHORED`.
- `internal/aifoundation 0.1.0` (guardrails `SafetyGate`) y `internal/search 0.1.0` (recuperación FAQ) — de este repositorio.
- LLM provider, canal y tools de dominio: adapters `CONDITIONED` (WhatsApp ya existe como `PYTHON_META_WHATSAPP_CLOUD_ADAPTER`).

## Archivos materializados (9)

| Archivo | SHA-256 |
|---|---|
| internal/agent/session.go | 0b1b178493505f8282838f74a0140de3350d7ff555b4c36d162ae7c16e3c02c7 |
| internal/agent/router.go | 139567afadddba98dc7dda418c9750bc6437a12fb071af7f85ddb5fef524809c |
| internal/agent/tool.go | 54b3c15615b9a342c62f0e347c4688d3b13c9dbda3b614571891c015b0cb5e1d |
| internal/agent/knowledge.go | 7bda6e8318b933d51967565100335b38f7f3faf172b24df1db2a23caa0cff60e |
| internal/agent/agent.go | d2e2bd5d5953a0a92a8ace066522349a5bbb90347de0e330c281e6a2f7b21314 |
| internal/agent/router_test.go | 68b0e2bc4d2444c2c9638cc7b6b0fdcd7a09a924596a2d3a094a93c73ffd9cbb |
| internal/agent/tool_test.go | 50e5e28185f4901d56f8a1863fc87882caea5173ccf093a474bb9e535dbd7843 |
| internal/agent/knowledge_test.go | 3cfb191f87f69458b857958a8d303589e418a8944bb05cc52fbce4a68df174ed |
| internal/agent/agent_test.go | 6cb6b8b284ea5e162d192cce52dc07e0685f674f6b17fd6f1a79ff859b0a2107 |

SHA-256 del pack: `992c2e2f8cbad0ae6b6d16b7c7087a4edfd051f5fc4458fab63a4657ae8be1ab` (27.907 bytes).

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).
- Sin PostgreSQL para este pack (núcleo en memoria; la persistencia de sesión y las tools de dominio son integración del proyecto).

## Verificación ejecutada

- `go test ./internal/agent/ -count=1`: **13/13 PASS** (ruteo de 10 intenciones + precedencia handoff, registro de tools con duplicado/nil/empty rechazados, conocimiento tenant-scoped sin fuga cross-tenant, saludo, fallback→handoff, tool éxito/necesita-info/fallo, FAQ recupera/vacía→handoff, guardrails pre/post→handoff, y conectividad agente↔`search.Store`).
- `go test ./... -count=1` (aifoundation + search + cache + agent): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 9/9 bloques reproducen byte a byte los archivos probados; `go test ./...` sobre el árbol materializado (4 paquetes) PASS.

## Cobertura de invariantes probada

1. `tenant_id` obligatorio en sesión y recuperación (sin fuga cross-tenant).
2. Una intención autoriza a lo sumo una tool; intención sin tool deriva a humano.
3. Guardrails pre/post fail-closed: injection o salida insegura → handoff.
4. Conocimiento vacío no se inventa → handoff.
5. Tool con `ErrNeedsInfo` mantiene la sesión abierta (colección de datos en turnos siguientes); fallo de tool → handoff con error registrado.
6. Evidencia de turnos/tools append-only.

## Condiciones residuales

- LLM provider real (OpenAI/Anthropic/vLLM) y ruteo/answer asistido por LLM: pendiente (adapters `CONDITIONED`).
- Canal web/WhatsApp y su autenticación/antiabuse: pendiente (WhatsApp adapter ya materializado).
- Tools de dominio reales (agenda→`GO-FRANCHISE-CUSTOMER-JOURNEY-API`, venta→`GO-COMMERCE`, devoluciones→`GO-RETURN-*`): adapters que bindean al contrato `Tool`.
- Colección multi-turno de argumentos con validación por campo: pendiente.
- Persistencia durable de sesión/evidencia (cache/outbox/Tessera): integración del proyecto.
- Evaluaciones de conversación (golden set de intenciones/respuestas): pendiente.

V178 añade el núcleo del agente conversacional; no declara el chatbot en producción (LLM, canales y reglas de negocio son inputs del proyecto).
