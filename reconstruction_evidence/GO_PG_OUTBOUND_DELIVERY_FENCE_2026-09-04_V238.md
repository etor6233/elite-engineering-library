# Evidencia V238 — fence durable de entrega outbound

Fecha: 2026-09-04  
Pack: `GO-PG-OUTBOUND-DELIVERY-FENCE 0.1.0`  
Pack SHA-256: `efbc71a1e6400a8325ff94821a3e8af2fdfc80d4f34e25872c36e0727f640412`  
Plan completo SHA-256 tras selección: `1ff2cc2da7b7d4601c18eedd7a400bea6c68852c67b770e917e42ac54eaab35c`

## Claim demostrado

Once archivos `AUTHORED` eliminan el retry ciego entre respuesta durable y proveedor externo:

- claim PostgreSQL antes de invocar al sender;
- una única invocación automática por tenant/canal/DeliveryKey;
- replay aceptado sin segunda llamada;
- payload divergente rechazado;
- timeout, receipt inválido, commit incierto o lease vencida quedan `unknown` y no se reenvían;
- cierre `unknown→accepted` o `unknown→failed_terminal` únicamente por reconciliación explícita con evidencia;
- recipient y provider message ID sólo como HMAC tenant/canal; el texto no se almacena en el ledger;
- secuencia de eventos append-only;
- `app.NewProduction` rechaza registry vacío o cualquier canal sin el wrapper durable;
- E2E real local: lead/binding → conversación → cotización → reply → receipt PostgreSQL; dos dispatch producen una llamada LLM de tool + continuación, un efecto de dominio y una sola llamada al sender.

## Autoridad y procedencia

- AWS Prescriptive Guidance documenta el problema dual-write, transactional outbox, duplicados e idempotencia: <https://docs.aws.amazon.com/prescriptive-guidance/latest/cloud-design-patterns/transactional-outbox.html>.
- PostgreSQL 18 gobierna aislamiento serializable y manejo de concurrencia: <https://www.postgresql.org/docs/18/mvcc.html>.
- Los once archivos son `AUTHORED`; no son código textual de AWS, PostgreSQL, Meta, Google ni OpenAI.

## Gates reproducidos

- round-trip: 11/11 SHA-256 idénticos;
- `gofmt` sobre ocho archivos Go: no-op;
- PostgreSQL 18.6, base limpia `elite_outbound_v238`: 48/48 migraciones;
- SQL focal 0048: PASS con ledger inmutable y attempt count fijo en uno;
- tests focales outbound/channels/app/PostgreSQL: PASS;
- lease vencida: `unknown`, attempt count uno;
- reconciliación positiva y negativa: PASS;
- full suite: 50 paquetes `go test -count=1 ./...` PASS;
- `go vet ./...` y `go build ./...`: PASS.
- `VERIFY_LIBRARY_PASS`: 150 packs, 1.322 archivos materializables, 658 Markdown y 44 perfiles; perfil completo 58 packs/631 archivos;
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: 150 packs y 121 fuentes upstream; Playwright/Lighthouse y componentes focales ejecutables pasaron, mientras cada gate live no autorizado quedó explícitamente omitido por el propio runner.

## Semántica conservadora

Esto demuestra **at-most-one automatic provider invocation**, no exactly-once delivery. La aceptación o entrega real depende del receipt y webhook/polling del proveedor exacto. El proyecto todavía debe probar cuenta, template/recipient, status events, rate/costo, opt-out, reconciliación y recuperación live; un estado `unknown` no autoriza otro envío.
