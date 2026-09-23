# Mercado Libre question outbound

`MercadoLibreQuestionSender` is an authored adapter governed by Mercado Libre's official Questions & Answers contract (last update 2026-01-15). It uses `POST /answers` with integer `question_id` and UTF-8 `text` limited to 2,000 characters, then confirms the exact seller/question/text/status through `GET /questions/{id}?api_version=4`.

The sender plugs into the existing `outbounddelivery.Channel`; it does not create another ledger. A durable, answer-hash-bound operator approval is required before the provider call. Network ambiguity becomes `unknown`, never an automatic retry. `ReconcileMercadoLibreQuestion` performs a GET-only reconciliation and transitions the existing PostgreSQL store only when the result is proved accepted or terminal.

No access token, answer body or provider response is written by this package. Target-account authorization, terms, quota, moderation behavior, privacy, callback origin, load, security and operational acceptance remain project gates.

Version 0.1.1 validates canonical positive int64 identifiers with Go's `strconv.ParseInt` and a lossless decimal round-trip. No overflowing or aliased identifier reaches HTTP. Only HTTP 400 carrying the documented `invalid_question` or `invalid_post_body` error is a proven rejection. Other non-2xx responses remain uncertain, with GET-only reconciliation and no automatic second POST. Reconciliation requires the exact seller and question; incomplete responses, unknown states or another identity cannot close the ledger.

Method authorities: https://pkg.go.dev/strconv#ParseInt defines numeric range errors; https://aws.amazon.com/builders-library/making-retries-safe-with-idempotent-APIs/ explains uncertain side effects and reconciliation. These corrections are AUTHORED and are not copied provider implementation. Use an HTTP transport with bounded timeouts and no automatic redirects or POST retries; the sender's single invocation does not itself prove transport-level exactly-once delivery.
