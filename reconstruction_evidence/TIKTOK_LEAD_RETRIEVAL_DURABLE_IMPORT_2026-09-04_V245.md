# Evidencia V245 — TikTok Lead retrieval e import durable

Fecha: 2026-09-04

## Decisión estrecha

TikTok Lead Generation dispone ahora de una cadena reconstruible `API autenticada → respuesta oficial preservada → candidato/rechazo + receipt hash-linked → verificación estricta → PostgreSQL/outbox compartido`. Esta evidencia no atribuye el wrapper a TikTok, no autoriza persistir señales webhook no autenticadas y no afirma acceso live.

## Procedencia exacta

- distribución oficial: `tiktok-business-api-sdk-official==1.1.3`;
- wheel: `tiktok_business_api_sdk_official-1.1.3-py3-none-any.whl`;
- SHA-256: `663b4a1f2585c4b386144d968646f695befd733c464418bbaeeac41befff33b7`;
- repositorio oficial fijado: `https://github.com/tiktok/tiktok-business-api-sdk` commit `f809c396520df2d7b201a9ccc5378d822b728ed3`;
- contrato GET: `https://business-api.tiktok.com/portal/docs/get-an-instant-form-lead-or-a-direct-message-lead/v1.3`, `/open_api/v1.3/lead/get/`;
- contrato subscription: `https://business-api.tiktok.com/portal/docs/create-a-subscription/v1.3`, `/open_api/v1.3/subscription/subscribe/`.

El wheel exacto contiene una divergencia upstream preservada: metadata de distribución `1.1.3` y `business_api_client.__version__ == 1.2.1`. El lock y el test gobiernan ambas declaraciones; la identidad ejecutable se fija por URL, filename, bytes y SHA-256. El repositorio/SDK revisado no contiene un cliente Lead generado por TikTok. Por ello `adapter.py` es `ADAPTED`; los demás diez archivos Python y los cinco archivos Go son `AUTHORED`.

## Resultado reconstruido

| Gate | Resultado observado |
|---|---|
| `VERIFY_LIBRARY.ps1` previo al endurecimiento | `PASS`: 153 packs, 1.356 bloques, 665 Markdown, 45 perfiles; TikTok Lead profile 51 archivos y franquicia 63/675 |
| composición desde destino vacío | `PASS`: 63 packs, 675 archivos de producto + `MATERIALIZATION_RECORD.md` |
| Python frozen install | cinco wheels hash-pinned, `pip check` PASS |
| Python adapter/evidence/webhook | 16/16 tests PASS |
| Go completo | `go test ./...`, `go vet ./...`, `go build ./...` PASS |
| PostgreSQL 18.6 desde cluster nuevo | 49 migraciones aplicadas; test TikTok focal PASS |
| package PostgreSQL real | 38/38 tests PASS, cero skips, incluidos TikTok y provider webhook |
| root + Audit posteriores a documentación | `VERIFY_LIBRARY_PASS`: 153 packs/1.356 bloques/667 Markdown/45 perfiles; `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: 153 packs/121 fuentes/16 adapters |

## Invariantes demostradas

- el profile `PROVEN` es obligatorio al construir el cliente;
- método, path, header y body coinciden con el lock del contrato v1.3;
- errores no filtran token;
- los tres artefactos se escriben atómicamente, no se sobrescriben y se ligan por SHA-256;
- valores dinámicos no-string/null/vacíos se rechazan sin perder la respuesta fuente;
- hashes de selector de cuenta malformados, dos selectores o ningún selector se rechazan antes del store;
- metadatos opcionales mal tipados se rechazan antes del store;
- candidato y receipt deben coincidir con respuesta, request ID, modo, fuente, cuenta, página, campos y cardinalidad;
- el candidato permanece `pending_policy`; no se concede consentimiento ni contacto;
- replay exacto es idempotente y divergencia usa el owner durable ya existente;
- el receipt no contiene campos del lead.

## Condiciones que no se convierten en PASS

- cuenta Business, app, permisos Lead, advertiser/library/page y términos reales;
- suscripción, callback TLS, test lead, cuotas/costo y reconciliación live;
- mecanismo oficial de autenticidad del webhook: no fue demostrado en las fuentes públicas revisadas;
- política jurídica/comercial de consentimiento, retención, borrado, contacto y promoción CRM;
- despliegue, observabilidad, carga, seguridad y aceptación en el proyecto target.

Mientras esas evidencias no existan, una señal webhook TikTok sólo puede producir `UNAUTHENTICATED_PROVIDER_SIGNAL`; la recuperación autenticada y la reconciliación son obligatorias antes de cualquier persistencia empresarial.
