# Go ML/AI Foundation — V175

## Resultado estrecho

V175 materializa `GO-ML-AI-FOUNDATION 0.1.0`, un contrato proveedor-agnóstico para servir modelos generativos de forma gobernada. No incorpora ningún SDK, modelo ni proveedor: es composición propia `AUTHORED` sobre la biblioteca estándar de Go, gobernada por el corpus de autoridad de IA de la biblioteca.

Cubre: identidad exacta de revisión (pin), registro de versiones con un único activo y auditoría inmutable, salida validada contra un schema objeto cerrado (`additionalProperties:false`), gates de seguridad fail-closed y decisión de release por evaluaciones. Es el cimiento que consumirán la búsqueda semántica y el agente conversacional; no declara producción ni sustituye clasificador de contenido, recognizer de PII, serving real, RAG ni agente.

## Autoridad y procedencia

- Código `AUTHORED`, sin copia de upstream. Autoridad de lenguaje: Go 1.26.7 (`https://go.dev`), sin dependencias de terceros.
- Gobernanza conceptual: `ML_PRODUCTION_LLMOPS_EVALUATION.md` (ciclo/evals/rollback), `AI_SECURITY_GOVERNANCE_PRIVACY.md` (guardrails/injection/PII), `AGENTIC_AI_SOFTWARE_ENGINEERING_CODEX.md` (agentes) y `NLP_RAG_RETRIEVAL_DATA.md` (retrieval), que son el corpus `SPEC_ONLY` de la biblioteca. Ninguno se copia como código.

## Archivos materializados (10)

| Archivo | SHA-256 |
|---|---|
| internal/aifoundation/model.go | 1c16cf9a11653f54ff23d045cd2583ef5963cf34b100c3fe71aec6fde5dab53b |
| internal/aifoundation/contract.go | bfd1278bf3a0f323656bc93152aded01c3df652459006adcd62a7b715bf93894 |
| internal/aifoundation/safety.go | f9657572c87e7a90857b8352c54b8d9830c2e3cec15ca3c4f6cd5b658bdee865 |
| internal/aifoundation/eval.go | 4e1b8f5fb9359e1407e914f70329f194692523f39cc9cbf79fcfda9c3652bf2b |
| internal/aifoundation/gateway.go | 8b57e18afdc17e49aa7179b42dd64d94a358f5d4c6e90bb252f58b784acf4791 |
| internal/aifoundation/model_test.go | 5bb768ffb73a4f45c05c753834452d6bb685f393641ea37e1534ed5e9fb782a7 |
| internal/aifoundation/contract_test.go | b2cab77c694656a34596a6113a9f9b7787974630329d3d24ca052b8a9f1a5cf2 |
| internal/aifoundation/safety_test.go | dbd8aa1dc61d8aac794ca7e85c2c61935c2a0063448b1f4aadc720831ef741fa |
| internal/aifoundation/eval_test.go | 0f892c4113e37f6c0a1497668050edbf7ffb5a9cfe615322ff0f3190cb5ee370 |
| internal/aifoundation/gateway_test.go | 17ed54ed4ae5fa890eea86a3dace1c627fef09a5549b8ab480c39575176259c6 |

SHA-256 del pack: `18e5330aa9e3f15d1de56adcfa4330a7fcfb55add7bb0ed7a74080e6bbfcf7a4` (36.861 bytes).

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado, fuera de PATH).
- Sin PostgreSQL para este pack (paquete en memoria, sin migraciones).

## Verificación ejecutada

- `go test ./... -count=1` sobre el árbol de trabajo: **27/27 PASS**.
- `go test ./... -count=1` sobre el árbol **materializado desde el Markdown** (round-trip): **PASS** (`ok elite.local/enterprise/internal/aifoundation`).
- `go vet ./...`: exit 0 (en ambos árboles).
- `go build ./...`: exit 0.
- Round-trip de hashes: 10/10 bloques del pack reproducen byte a byte los archivos probados (SHA-256 igual al declarado).
- `go test ./... -race`: no ejecutado en este host (requiere cgo no disponible); la carrera de promoción se cubre con mutex y un test de concurrencia determinista (8 candidatos → exactamente un activo).

## Cobertura de invariantes probada

1. Pin exacto: sólo revisions con digest sha256 (64 hex) se registran; `latest`/`auto` y digest inválido se rechazan.
2. Registro: a lo sumo una versión `active`; promoción exige gate superado; rollback no reactiva versiones; duplicados rechazados.
3. Schema cerrado: salida valida contra `additionalProperties:false`; clave desconocida, requerida ausente, tipo incorrecto, schema abierto y raíz no-objeto fallan cerrado.
4. Seguridad fail-closed: valor cero bloquea todo; injection bloqueado; PII sin allowlist bloqueado; refusal respetado; frase de salida bloqueada.
5. Gateway: pin activo, pre-flight, post-flight y structured output se aplican en cadena; 5 caminos negativos/positivos probados.
6. Evals: gate de release bloquea promoción bajo el umbral; suite vacía y errores de proveedor no se silencian.

## Condiciones residuales

- Adapter de proveedor real (OpenAI/Anthropic/vLLM/…) y sus contract tests: pendientes.
- Clasificador de contenido y recognizer de PII admitidos: gates del proyecto (las heurísticas no son garantía).
- Persistencia durable del registro de versiones (hoy en memoria): paso de proyecto.
- SEARCH, CACHE y RAG-AGENTS: consumidores aún no materializados.
- Serving, evals con datos reales, observabilidad, carga y rollback live: condicionados.

V175 añade la fundación de IA gobernada; no declara terminada la superficie `ML-AI` ni el chatbot.
