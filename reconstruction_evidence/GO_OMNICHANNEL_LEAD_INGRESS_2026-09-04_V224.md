# Go Omnichannel Lead Ingress — V224

## Resultado estrecho

V224 materializa `GO-OMNICHANNEL-LEAD-INGRESS 0.1.0`: una frontera provider-neutral y un adapter exacto para Google Lead Form. Autentica antes de persistir, conserva el SHA-256 del body fuente, elimina `google_key` del payload durable, registra el perfil de redacción, deduplica por `lead_id`, rechaza identidad divergente y transacciona raw redactado+candidato `pending_policy`+outbox.

No es una afirmación de adapter Meta/TikTok, contacto autorizado, cuenta Google live, exactly-once, cita o venta terminada.

## Autoridades exactas consultadas

- Google Lead Form Webhook, observado 2026-09-04 y actualizado por Google 2026-05-05: schema, `google_key` confidencial, int64, unknown fields, 200/4xx/5xx, duplicados y `lead_id`: <https://developers.google.com/google-ads/webhook/docs/implementation>.
- Sample oficial Google para `WebhookDelivery`: <https://developers.google.com/google-ads/api/samples/add-lead-form-asset>.
- CloudEvents 1.0.2, identidad `source`+`id`: <https://github.com/cloudevents/spec/blob/ce@v1.0.2/cloudevents/spec.md>.
- AWS Lambda/SQS, entrega at-least-once e idempotencia: <https://docs.aws.amazon.com/lambda/latest/dg/with-sqs.html>.
- Google Cloud Pub/Sub, retry/dead-letter: <https://docs.cloud.google.com/pubsub/docs/subscription-retry-policy>.
- PostgreSQL 18, aislamiento y transacciones: <https://www.postgresql.org/docs/18/transaction-iso.html>.
- OWASP Logging Cheat Sheet, exclusión/masking/hashing/encryption de secretos y datos sensibles: <https://cheatsheetseries.owasp.org/cheatsheets/Logging_Cheat_Sheet.html>.

Los diez bloques son código local: seis `ADAPTED` y cuatro `AUTHORED`; cero `VERBATIM`. Ninguno se presenta como source publicado por Google, Amazon, OWASP o PostgreSQL.

## Archivos materializados (10)

| Archivo | Procedencia | SHA-256 |
|---|---|---|
| `internal/leadstream/event.go` | ADAPTED | `60fee6649a1dc6f2cc65d5a11bd6e58bb6b4196d98c734ef8c248c189d547178` |
| `internal/leadstream/google_ads.go` | ADAPTED | `42d83974aeeca63375d0c5a91506a7189e9040584f84832261c1262466280ad1` |
| `internal/leadstream/store.go` | ADAPTED | `241b5a7a617910ac56326a8649c3cbf32fe6a5be0f763db72f4cc5f57070ae16` |
| `internal/leadstream/handler.go` | ADAPTED | `bb2eed2d0393407968bdb302f29aff868b28c4d1147ae0bc49c208230b640185` |
| `internal/leadstream/leadstream_test.go` | AUTHORED | `80c5bc3e9970536a1b453cfe62cb2abbd7151c7aaa7ba251e684d2f535202d94` |
| `internal/platform/postgres/lead_ingress.go` | ADAPTED | `f84ea9d5940ab5bac04b14d71dfa2a8fa4e4d8041181162fc4acce278a04a842` |
| `internal/platform/postgres/lead_ingress_integration_test.go` | AUTHORED | `993a73808c130bf0da7a293f918181bd8e3778818a585ef3250c9d8adba1556d` |
| `db/migrations/0044_omnichannel_lead_ingress.up.sql` | ADAPTED | `761153e744342c8bec1a77d5e93d81ee794189037276cbb051f216fdcd2a91cf` |
| `db/migrations/0044_omnichannel_lead_ingress.down.sql` | AUTHORED | `1cd0fb4673bdcb0d21004fd16677b64f3713c47144bb49b67b0155e7813f5167` |
| `db/tests/0044_omnichannel_lead_ingress.test.sql` | AUTHORED | `b4f1f2d8933424df91215abf69479c61e5a61bd6ffa0105996a26460085e7c10` |

SHA-256 del pack: `fae1d4df9c5aa031a8184b6ef708f89127fb370bcd46044cb0dbe0156a94c950`.

## Verificación ejecutada

- Go 1.26.7 oficial fijado: `go test ./... -count=1` PASS sobre 45 paquetes con `TEST_DATABASE_URL`/`DATABASE_URL` reales; `go vet ./...` y `go build ./...` PASS.
- HTTP/parser: credencial errónea no persiste; límite 1 MiB; JSON trailing/malformado; unknown fields; int64; replay exacto; identidad divergente; missing lead ID retenido como rechazo; `google_key` ausente del payload durable.
- Concurrencia: 32 entregas exactas en `MemoryStore` y 16 en PostgreSQL producen una sola identidad durable; las restantes son duplicate.
- PostgreSQL 18.6 Windows x86-64, checksums on y loopback: 44 migraciones y 30 tests SQL PASS desde base vacía.
- Repositorio real: raw redactado, hashes fuente/stored, perfil de redacción, candidato y outbox atómicos; divergencia rollback; rechazo sin candidato.
- 0044 down/up/test PASS.
- Round-trip Markdown: diez archivos, cardinalidad exacta y SHA-256 source=materializado PASS bajo errores terminantes.
- Lifecycle: arranque, health, gates, stop fast y cero procesos PostgreSQL remanentes PASS.

## Fallos que mejoraron el resultado

- `LIB-FAIL-1766/1767`: fixture append-only repetible y aislamiento de la cola global.
- `LIB-FAIL-1769`: se invalidó un falso PASS de round-trip y se reemplazó por comparación fail-closed de diez hashes.
- `LIB-FAIL-1771`: se bloqueó la persistencia clara de `google_key`; el modelo ahora separa hash fuente de payload durable redactado.
- `LIB-FAIL-1772`: se descartó investigación truncada y se abrieron las autoridades Google/OWASP exactas.

## Condiciones residuales

- `go test -race` no puede ejecutarse en este host sin CGO/GCC; sigue `BLOCKED_EXTERNAL` y debe pasar antes de promover el pack a `REUSABLE_PACK`.
- Cuenta, formulario, URL TLS, secreto rotado, delivery/retry real, WAF/rate limit, métricas y reconciliación Google siguen siendo gates de proyecto.
- PII requiere cifrado, acceso, residencia y retención aprobados; este pack no inventa esas decisiones.
- Meta y TikTok requieren adapters separados basados en sus schemas/auth oficiales exactos.
- El paso candidato→CRM/contacto→conversación→cita/venta y el composition root permanecen fuera de este claim.
