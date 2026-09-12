# Go Official Payment Webhook Adapters — 2026-08-26 V1

## Upstreams oficiales exactos

- Stripe `stripe/stripe-go` 86.3.0, commit `a2df585a800a97fe8ec4ebf551b4449bdb3d90a1`, MIT.
- Mercado Pago `mercadopago/sdk-go` 1.14.0, commit `f910ee53fbb6819e435eaf3d0f800cb1fe74ae09`, MIT.

Ambos sources fueron adquiridos con `OFFICIAL-UPSTREAM-ACQUISITION-CORE` y sus archives/hash/licencias exactos. El módulo usa sus paquetes webhook y clientes oficiales PaymentIntent/Payment; no reimplementa ni atribuye a Elite la verificación, serialización o autenticación upstream.

## Reconstrucción

```text
Go: go1.26.7 windows/amd64, distribución oficial verificada
Materialized 7 files into ...\elite-official-payment-adapters-verify-v2-20260826
all modules verified
GOPROXY=off go test ./...
ok example.com/elite/official-payment-webhooks/officialpayments
```

Tests: signature y API-version Stripe; signature/tolerance Mercado Pago; body/firma alterados; ID query/body divergente; input faltante; body >1 MiB. Outbound Create se ejecutó a través de ambos clientes oficiales contra transports locales controlados y verificó authorization, idempotency, amount/currency, order metadata y respuestas normalizadas. El test Stripe genera su header con el helper oficial del SDK.

## Fallos aprendidos

- `LIB-FAIL-028` recurrencia 4: dos IDs enviados mediante `pwsh -File`; corregido con call operator.
- `LIB-FAIL-035`: `stripe.EventType` no es `string`; conversión explícita en la frontera.
- `LIB-FAIL-036`: prefijo del árbol draft divergente del manifest; updater abortó sin cambio.
- `LIB-FAIL-037`: primer `go mod verify` se lanzó desde cwd incorrecto; repetido desde el módulo.
- `LIB-FAIL-038`: el cliente Mercado Pago amplió el grafo a `google/uuid`; `go.sum` se regeneró y el módulo volvió a pasar offline.

## Límites

Estado `REBUILD_VERIFIED / CONDITIONED`. La request outbound fue probada sin red real a través del SDK oficial, pero no se probó una cuenta sandbox ni un cobro/refund/capture, settlement, disputes, OAuth, país/producto, PCI o reconciliación real. El core durable gobierna deduplicación/job/reconciliation. QR Mercado Pago no se admite porque el propio SDK declara que no lleva esta firma.

Fecha/revisor: 2026-08-26 / Codex.
