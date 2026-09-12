# Evidencia V258 — Mercado Libre Question Outbound conectado

## Resultado

```yaml
date: 2026-09-05
decision: REBUILD_VERIFIED_CONDITIONED
production_ready_project: false
pack: GO-MERCADOLIBRE-QUESTION-OUTBOUND@0.1.0
outbound_owner: GO-PG-OUTBOUND-DELIVERY-FENCE@0.2.0
franchise_profile: 67 packs / 704 files
library: 160 packs / 1395 materializable files / 121 locked upstream sources
```

La franquicia ya no termina Mercado Libre Questions en un candidato aislado. El recorrido materializado es:

```text
notification questions
→ GET question v4
→ normalized durable candidate
→ conversation/policy
→ durable approval bound to tenant + DeliveryKey + question + SHA-256(answer)
→ existing PostgreSQL outbound claim
→ POST /answers (once)
→ GET /questions/{id}?api_version=4 exact confirmation
→ accepted receipt | terminal evidence | unknown awaiting GET-only reconciliation
```

No se creó otra tabla, cola, store ni runtime. Los tres archivos nuevos viven en el mismo package `internal/outbounddelivery`; el owner PostgreSQL existente conserva intento, receipt, HMAC, eventos append-only y estado.

## Autoridad oficial observada

La URL oficial [Mercado Libre — Manage questions & answers](https://developers.mercadolibre.com.ar/en_us/categories-and-attributes/manage-questions-and-answers), visible el 2026-09-05 y marcada `Last update 15/01/2026`, documenta:

- `POST https://api.mercadolibre.com/answers`;
- payload JSON exacto con `question_id` entero y `text`;
- Bearer token y `Content-Type: application/json`;
- UTF-8 para caracteres especiales;
- máximo 2.000 caracteres;
- recomendación semi-automática: sugerencias recibidas por operadores;
- `GET /questions/{id}?api_version=4`, seller/question/status/answer;
- `invalid_question` e `invalid_post_body` como rechazos documentados;
- notificaciones del topic `questions` para flujo en tiempo real.

El fetch automatizado del portal recibió 403 y se registró como `LIB-FAIL-2060`; la misma URL fue inspeccionada mediante navegador. No se usó una fuente secundaria ni se inventó el contrato.

## Procedencia y diseño

Los cambios son `AUTHORED`, no código copiado ni atribuido a Mercado Libre, AWS o PostgreSQL. El contrato HTTP estrecho procede de Mercado Libre. El control de dual-write/ambigüedad continúa gobernado por las autoridades AWS Transactional Outbox y PostgreSQL MVCC ya fijadas en el owner.

Decisiones comprobables:

- el modelo no autoriza el envío: un resolver durable debe devolver una aprobación exacta, firmemente ligada al hash del texto;
- un 2xx no basta: el adapter hace read-after-write y exige seller/question/text exactos, question `ANSWERED` y answer `ACTIVE`;
- un 4xx demostrado produce `TerminalFailure` con código/evidencia y cierra `failed_terminal`;
- timeout/error de transporte produce `unknown` y el replay no vuelve a llamar al proveedor;
- la reconciliación sólo hace GET y nunca reenvía;
- `UNANSWERED` después de una ambigüedad permanece `unknown`; divergencia o estado terminal cierra sin fingir aceptación;
- token, texto y body provider no se persisten en el ledger; sólo hashes/HMAC y receipts.

## Artefactos

| Artefacto | SHA-256 |
|---|---|
| `implementation_packs/GO_MERCADOLIBRE_QUESTION_OUTBOUND.md` | `1f91fbf1ba7590e676899b67411d9dad3340ad4d1bff142228beae39a1d099fd` |
| `implementation_packs/GO_PG_OUTBOUND_DELIVERY_FENCE.md` | `53095201b99ed2746e487a2f33a13023ffb438e5024dfd2c576537be6b2a48e5` |
| `markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md` | `5af49f88692800c40d71573a62e5d751098edabeadeff7bbff01cd805aa873cf` |
| Audit JSON externo | `113902030ece068c9927ca352d84284994ec22d62d8accd87ba5dd8c0349528d` |

## Gates ejecutados

- toolchain oficial: `go1.26.7 windows/amd64`, ejecutable SHA-256 `5463fe58fa999d74420f00ee1b36d31c3da90a57ff204159f159859567ad61fc`;
- pack nuevo: materialización/round-trip byte-idéntico 3/3;
- composición canónica: 67 packs/704 archivos, cero colisiones;
- package outbound: 10 tests top-level más subcasos; POST/GET/headers/body, approvals, 2.000 runas, seller/text/status, timeout, terminal y tres reconciliaciones;
- grafo Go completo compuesto: `go test ./...`, `go vet ./...`, `go build ./...` PASS;
- PostgreSQL 18.6 limpio con data checksums: 49 migraciones; `outbounddelivery`, `platform/postgres` y `app` PASS con `TEST_DATABASE_URL`; servidor detenido normalmente;
- `VERIFY_LIBRARY_PASS`: 160 packs, 1.395 archivos materializables, 690 Markdown antes de añadir este expediente, 51 perfiles; franquicia 67/704;
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: 160 packs, 121 upstreams, 16 provider adapters y todos los gates offline aplicables. OpenGrep sigue bloqueado por Cosign SCA y fuera del perfil.

## Condiciones retenidas

No se contactó una cuenta real ni se envió una respuesta. Antes de producción el proyecto debe demostrar owner/app/seller, token store/rotation, scopes, términos, privacidad/retención, política y autoridad de aprobación, moderación/BANNED, cuotas/costo/rate limits, callbacks reales, carga, observabilidad sin PII, seguridad ofensiva, recovery/rollback y aceptación empresarial. El PASS de biblioteca no convierte esos inputs externos en `YES`.
