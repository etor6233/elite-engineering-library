# Architecture Pack — Observability and Secure Delivery Baseline

> **Estado:** `CANDIDATE_PACK`.  
> **Fuentes:** OpenTelemetry Collector, Prometheus, Kubernetes/OpenTofu condicionados, Playwright, supply-chain authorities.  
> **Claim:** baseline de señales, CI evidence, artifact provenance, rollout, rollback y recovery.

## 1. Telemetry contract

Todo servicio/app/worker emite:

```text
service.name / service.version / deployment.environment
trace_id / span_id / request_id / correlation_id
tenant pseudonymous key only when allowed
operation / outcome / stable error_code
duration / retry_count / queue_age
release/config/schema/policy model versions
```

No emitir passwords, tokens, cookies, authorization headers, card data, documents, full webhook bodies ni PII libre. Allowlist de attributes antes de exportar.

## 2. OTel pipeline baseline

```text
applications → OTLP local/agent or gateway
→ memory limiter
→ resource normalization
→ attribute allowlist/redaction
→ sampling decision
→ batching + bounded queue
→ metrics/traces/log backends
```

- health/readiness del collector;
- queue/drop/export failure metrics;
- TLS/mTLS según boundary;
- configs versionadas y validated;
- processor bloqueante no comparte pipeline sin entender backpressure;
- telemetry outage no derriba el business path salvo requerimiento de audit separado.

Audit evidence no depende únicamente de logs/telemetry; usa storage/contract dedicado.

## 3. Métricas

### Technical RED/USE

- request/job rate, errors, duration;
- CPU, memory, GC, thread/event loop, connection pools;
- DB acquisition/query/transaction/locks;
- queue depth/oldest age/attempts/dead letters;
- provider latency/errors/quota/auth;
- cache hit/miss/eviction/stampede y search lag.

### Business correctness

- lead acceptance/dedup/assignment lag;
- stock reservation conflicts and aged reservations;
- orders stuck by state;
- payment uncertain/reconciliation diffs;
- catalog projection/feed lag;
- cross-system inventory/order mismatch;
- warranty cases past SLA;
- privilege/grant changes and denied high-risk actions.

Prometheus no se usa como ledger/billing exacto.

## 4. Tracing

- ingress creates/accepts trace context after validation;
- spans around DB, queue, authorization and provider calls;
- async links preserve causal relation from outbox to consumer;
- sampling retains errors/high latency/high-risk operations;
- payloads are events/attributes allowlisted, not bodies;
- trace IDs included in provider attempt/audit evidence where safe.

## 5. Logging

Structured event:

```json
{
  "timestamp": "...",
  "level": "...",
  "service": "...",
  "version": "...",
  "operation": "...",
  "outcome": "...",
  "error_code": "...",
  "trace_id": "...",
  "resource_type": "...",
  "resource_id_hash": "..."
}
```

Stack traces sólo server-side, fingerprinted y access-controlled. Log volume tiene budgets/rate limits. Security/audit events tienen retention/immutability distintos.

## 6. SLO y alerting

Cada SLO declara:

```text
user journey and failure definition
SLI query/source
target/window
dependencies/exclusions
error budget policy
alert thresholds and burn rates
owner/runbook
```

Alerts son accionables: page por burn/correctness/data loss/security; ticket por tendencia/capacity. No alertar por CPU aislada sin impacto o riesgo demostrado.

## 7. CI pipeline

```text
source checkout at immutable revision
→ formatting/static/type/unit
→ secret/license/dependency/SAST checks
→ real DB integration + contract tests
→ Playwright critical journeys
→ migration compatibility
→ build once
→ SBOM + provenance + vulnerability scan + signature
→ policy gate
→ publish immutable artifact
```

Generated code y lockfiles se verifican contra source. CI usa short-lived/workload identity; no long-lived cloud keys. Pull requests no confiables no reciben secrets.

## 8. Artifact y provenance

- source revision, clean/dirty state y builder identity;
- dependency lock and checksums;
- base image digest, not tag only;
- SBOM including OS/runtime/application;
- license evidence/notices;
- build parameters and generated contracts;
- signature/attestation stored with artifact;
- same artifact promoted; no rebuild per environment.

## 9. Deployment

Preflight:

- config/secret references resolve;
- schema compatibility and migration lock risk;
- capacity/quota/dependency health;
- rollback compatibility;
- synthetic identity and critical provider sandbox/ping where allowed.

Rollout:

```text
deploy no-traffic → readiness/smoke
→ internal/pilot tenant
→ canary percentage/market
→ observe technical + business invariants
→ promote progressively
```

Feature/data migrations se desacoplan del binary rollout. Health endpoint no llama todas las dependencias; readiness refleja capacidad de servir sin cascading failure.

## 10. Rollback y roll-forward

- binary rollback documentado y probado;
- schema expand/contract mantiene compatibilidad;
- irreversible migration usa backup/checkpoint y roll-forward plan;
- bad catalog/policy/config tiene version activation rollback;
- provider mappings/version canary by connection;
- mobile client backward compatibility window explícita.

No prometer rollback instantáneo si hubo external side effects; allí se necesita compensation/reconciliation.

## 11. Recovery

Drills:

- restore PostgreSQL to point-in-time;
- rebuild cache/search/projections from source;
- replay outbox/inbox safely;
- rotate identity/provider credentials;
- operate with provider/telemetry outage;
- recover deleted/misconfigured policy/catalog version;
- regional/cloud exit proportional al risk.

Evidence incluye timestamps, RPO/RTO observado, missing/corrupt records checks, business invariant queries y follow-ups.

## 12. Security gates

- threat model updated for changed boundary;
- authn/authz negative tests;
- dependency and container vulnerabilities triaged by reachability/exposure;
- secret scan and runtime secret exposure tests;
- egress allowlist/proxy where required;
- CSP/security headers/mobile secure storage;
- upload/webhook/provider adversarial tests;
- admin high-risk action audit;
- incident owner, severity model and notification obligations.

## 13. Release evidence gate

- [ ] requirements ↔ tests ↔ telemetry;
- [ ] all critical suites passed with artifact links;
- [ ] migration and rollback evidence;
- [ ] SBOM/provenance/signature/license bundle;
- [ ] canary plan and abort thresholds;
- [ ] dashboards/alerts/runbooks live;
- [ ] backup freshness and latest restore drill within policy;
- [ ] OPEN risks accepted by named authority;
- [ ] post-release verification and ownership window.
