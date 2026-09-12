# Go OpenAI Responses Tool Adapter — evidencia V228

## Resultado

```yaml
pack: GO-OPENAI-RESPONSES-TOOL-ADAPTER
version: 0.1.0
authority: SUPPORTED_REFERENCE
implementation: REBUILD_VERIFIED
admission: CONDITIONED
verified_at: 2026-09-04
go: 1.26.7
tests: 4/4 PASS
materialization: 2/2 PASS
gofmt_idempotence: 2/2 PASS
full_go_packages: 47 PASS
go_vet: PASS
go_build: PASS
```

## Contrato oficial aplicado

La referencia oficial vigente define `POST /responses`, `tools` con funciones tipadas, outputs `function_call` con `call_id/name/arguments`, y devolución `function_call_output`. La guía oficial recomienda entregar tools mediante el campo API en vez de inyectar descripciones en prompt y parsear texto libre.

El adapter local demuestra con HTTP falso inspeccionable:

- tool top-level `type=function`, schema, descripción y `strict=true`;
- `parallel_tool_calls=false`, modelo exacto, storage aprobado y cache key;
- validación local del schema y de cada argumento antes de exponer la llamada;
- rechazo de función no ofrecida y propiedades desconocidas;
- continuación por `previous_response_id` + `function_call_output`/`call_id`;
- `tool_choice=none` y rechazo de una segunda llamada en la fase final;
- contabilización de tokens informada por provider.

## Procedencia

Los dos archivos son `ADAPTED` desde la documentación pública oficial; no contienen SDK ni código fuente OpenAI y no se presentan como OpenAI-authored.

## Hashes

| Archivo | SHA-256 |
|---|---|
| `internal/llmopenai/responses_tools.go` | `4f8969035730e44f7b440d153f77e36042958e17973d6eb20c1fd84549c4eb65` |
| `internal/llmopenai/responses_tools_test.go` | `9f50beaaf84c2243d1f56a21ca3fbcf30fc859506d0f82959fac12d59a19eee7` |

Pack Markdown SHA-256: `96c2b290e374c9a670b530b49be4635a8a793dffd59f4cde845e384d6a01757b`.

## Límites retenidos

- No se usó cuenta ni modelo live; disponibilidad, costo, latencia y comportamiento se prueban por proyecto.
- `store=true` sólo se admite si la retención/provider policy del proyecto lo aprueba.
- El adapter no autoriza ni ejecuta el efecto: eso corresponde al runtime de aplicación y al dominio.
- Falta el wiring completo canal→identidad→presupuesto→Responses→approval/tool→persistencia→respuesta.
