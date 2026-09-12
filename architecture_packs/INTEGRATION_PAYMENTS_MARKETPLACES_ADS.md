# Architecture Pack — Payments, Marketplaces and Ads Integrations

> **Estado:** `CANDIDATE_PACK`.  
> **Fuentes verificadas:** Stripe SDK MIT, Mercado Pago SDK MIT, Amazon SP-API models Apache-2.0, Google Ads clients Apache-2.0, Meta Business SDK bajo Platform License.  
> **Regla:** licencia del SDK y autorización para usar la API son controles distintos.

## 1. Anti-corruption layer

```text
domain command/event
→ canonical integration contract
→ provider adapter
→ official SDK/generated client
→ provider API
→ durable evidence/cursor/mapping
→ reconciliation back to domain
```

El core no importa tipos del SDK. Los provider IDs y raw payloads no reemplazan IDs/estados internos.

## 2. Interface mínima

```ts
interface ProviderAdapter<Command, Result> {
  validate(command: Command): ValidationResult;
  execute(command: Command, context: {
    connectionId: string;
    idempotencyKey: string;
    deadline: Date;
    traceId: string;
  }): Promise<Result>;
  classify(error: unknown): 'retryable' | 'terminal' | 'auth' | 'quota' | 'unknown';
  reconcile(reference: ProviderReference): Promise<CanonicalState>;
}
```

La implementación concreta agrega SDK version/API version y no loguea credentials, tokens, card/customer payloads ni audiences.

## 3. Persistencia de integración

Entidades:

```text
connection        tenant/provider/market/scopes/status/credential_ref
mapping           internal_type/id ↔ provider/type/id + mapping_version
sync_cursor       connection/resource/cursor/watermark
sync_run          direction/range/counts/status/error summary
delivery_attempt  operation/idempotency/provider request ID/outcome
webhook_event     provider/event ID/signature metadata/raw hash/state
reconciliation    expected/observed/diff/resolution
```

Credential material vive en secret manager/KMS; DB guarda referencia, version y rotation metadata.

## 4. Webhook ingestion

Secuencia obligatoria:

1. aplicar body/headers/size/deadline limits;
2. conservar raw body para verificación;
3. verificar firma, timestamp y endpoint secret activo;
4. calcular event identity/dedup key;
5. persistir evento `received` antes de responder;
6. responder rápidamente dentro del contrato del proveedor;
7. procesar async con inbox/idempotencia;
8. obtener estado autoritativo del proveedor cuando el evento no sea suficiente;
9. aplicar transición de dominio válida;
10. marcar processed/terminal y reconciliar pendientes.

Rotación soporta overlap de secrets limitado. Eventos inválidos se cuentan y retienen sólo lo necesario; no se reintentan internamente.

## 5. Payments

### Modelo canónico

```text
PaymentIntent → PaymentAttempt* → Authorization/Capture/Failure
                              ↘ Refund* / Dispute* / SettlementLink
```

Invariantes:

- dinero en minor units + currency;
- order total snapshot no se recalcula desde catálogo actual;
- browser redirect nunca confirma pago;
- create/capture/refund usa idempotency key estable;
- state uncertain inicia reconciliation, no retry ciego;
- webhook duplicado/fuera de orden no retrocede state machine;
- no almacenar PAN/CVV; hosted/tokenized collection reduce PCI scope;
- ledger contable, si existe, es módulo separado con reglas contables explícitas.

Stripe aporta webhook signature helpers e idempotency metadata; Mercado Pago aporta SDK oficial. Ninguno elimina country/product onboarding, PCI, chargeback, settlement ni reconciliation.

## 6. Marketplaces

Capabilities se separan:

```text
catalog/listing
price
inventory availability
orders/acknowledgement
shipment/tracking
returns/cancellations
fees/settlements
messages
```

No crear un adapter gigante. Cada capability declara direction, authority, frequency, quota, freshness SLO y conflict policy.

### Amazon SP-API

- generar/validar clients desde modelos fijados;
- Least-privilege roles/scopes y Restricted Data Tokens cuando aplique;
- reports/feeds son async: submit → poll/notification → download → validate;
- rate limits por operation/account; token bucket y retry-after;
- sample code es educativo; production hardening propio.

### Mercado Libre

- SDKs históricos oficiales están archivados/deprecados: no depender de ellos;
- usar documentación/contratos vigentes y HTTP client propio/generado;
- OAuth refresh, seller/site/country boundaries;
- items/orders/shipments/questions se modelan en adapters separados;
- cambios de API se detectan con contract tests y sandbox/canary seller.

## 7. Ads y attribution

Modelo canónico:

```text
campaign → ad group/set → creative/ad → targeting reference
budget + schedule + status
performance observation(window, attribution model, currency)
conversion event(consent, purpose, source, dedup ID)
```

- plataforma de ads no es source of truth de customer/order/revenue;
- conversions usan mínima data y consent/purpose comprobable;
- hashing no vuelve anónimo un identificador por sí solo;
- audiences tienen provenance, expiry y deletion propagation;
- metrics se guardan con attribution window/model y timezone;
- budgets requieren guardrails, approval y spend anomaly alert;
- API deprecations se calendarizan como operational risk.

Google Ads tiene clientes oficiales Apache-2.0 y versiones frecuentes. Meta Business SDK está limitado a los servicios/APIs Meta y Platform Policy: se usa sólo dentro de esos términos.

## 8. Retry y quota policy

```text
connect timeout < request deadline < job lease
retry only classified transient errors
exponential backoff + full jitter + cap
respect Retry-After/provider quota headers
per-connection and global concurrency limits
circuit/open state does not discard durable work
```

Autenticación/permission/validation son terminales hasta intervención/config change. `unknown` no se reintenta infinitamente.

## 9. Reconciliation

Jobs mínimos:

- payment intents no terminales;
- captures/refunds/disputes versus orders;
- marketplace order/fulfillment status;
- listing/price/inventory diffs;
- ads campaign status/spend/conversions;
- stale cursor/no events heartbeat;
- credential/scope health.

Cada diff tiene severity, owner, auto-fix policy y evidence. Reconciliation no modifica dominio silenciosamente si viola invariantes.

## 10. Contract tests

- captured fixtures sanitizadas + schema validation;
- provider sandbox tests separados de mocks;
- signature/replay/clock-skew webhooks;
- duplicate/out-of-order/missing events;
- rate limit and pagination/cursor boundaries;
- token expiry/rotation/scope loss;
- version deprecation fixture;
- unknown field tolerance y removed-required-field failure;
- monetary rounding/currency/timezone;
- reconciliation converges after injected failures.

## 11. Go-live gate por conexión

- [ ] terms/DPA/permissions/account review completos;
- [ ] production credentials in secret manager;
- [ ] exact scopes and owners;
- [ ] API/SDK/model version pinned;
- [ ] idempotency, webhook y reconciliation tests;
- [ ] quotas/load budget y backfill plan;
- [ ] alert/runbook/provider support path;
- [ ] data classification, retention y deletion propagation;
- [ ] sandbox → pilot account → staged rollout;
- [ ] disconnect/export/exit plan.
