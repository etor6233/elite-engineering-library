# Official return refund worker

## 0.1.4 — observable host report failure

The real run loop now checks log delivery before claiming another job. A failed,
short or invalid step-log write returns static REFUND_REPORT_FAILED; the owning
main exits1 via its existing private REFUND_HOST_FAILED boundary. Each step log
makes one Write, with no retry of the uncertain record and no formatting of the
underlying step/sink error. The terminal host message is a separate best-effort
write; exit status remains available even when stderr is unusable.

Successful reports preserve the configured logger prefix/flags and fixed
REFUND_STEP_FAILED message. Nil/ErrNoWork still emit nothing; ordinary step
errors continue under the existing polling policy when their report succeeds.
An already-cancelled loop does not generate a token or start another step.

This is AUTHORED runtime wiring, not a provider/financial-policy change. Durable
refund reconciliation, idempotency, amounts, currencies and results are unchanged.
The owner must supervise the exit status and diagnose durable results before a
restart; this patch does not install a supervisor, endpoint, alert or retry daemon.
A blocking writer still needs target-owned process supervision. Complete local
Write acceptance does not prove retention, remote acknowledgement or alerting.

Tests cover six failed-write modes against exact0.1.3 before the fix, the actual
host loop, a real closed file, privacy, healthy continuation/idle/cancellation,
and a subprocess exiting1 with an unusable report sink. Finite fuzz covers write
counts/error combinations. No provider credentials, payment calls or live effects
are used. Source contracts: https://pkg.go.dev/log#Logger.Output and
https://pkg.go.dev/io#Writer; executed toolchain stays Go1.26.7.

This standalone Go worker claims only `refund/payment` return effects from the enterprise PostgreSQL schema. It resolves the exact allocated order line and one unambiguous captured payment, creates a partial refund through the exact official Stripe or Mercado Pago SDK, persists immutable provider observations, retrieves non-final refunds, and completes the effect only after an exact final provider result.

Required runtime values are `DATABASE_URL` and `REFUND_WORKER_ID`; add `STRIPE_SECRET_KEY` and/or `MERCADO_PAGO_ACCESS_TOKEN` only through the deployment secret manager. A missing provider credential blocks that provider's effect without exposing a secret. The Mercado Pago adapter is deliberately ARS/two-decimal until another market has an admitted currency-exponent contract.

Production admission still requires provider account/country enablement, sandbox and live contractual tests, webhook delivery into the durable inbox, reconciliation alerts, business approval of line-level refund allocation, and load/security/recovery evidence for the selected target.

Host stderr reports only REFUND_HOST_FAILED (startup exit1) or REFUND_STEP_FAILED (step failure; inspect the durable result). Arbitrary database/provider errors are not formatted in host logs. Idle results do not emit failure messages. This does not prove target log delivery, retention or alerts.


## Optional closed-outcome operational telemetry

Set both REFUND_TELEMETRY_CONFIG and REFUND_TELEMETRY_CONFIG_SHA256 to an exact
local JSON config accepted by operationaltelemetry.Config. Pin CA, client
certificate and private-key bytes individually. TLS1.3 mutual authentication,
no proxy or redirects, response and request byte caps, and one shared finite
observation deadline are required. With neither variable, the previous local
logger behavior remains. A required-observability project must configure and
verify this path explicitly before its own admission; absence is not readiness.

Outcomes are started, idle, handled, error and stopped. handled only means Step
returned nil; the database and provider receipts own the refund result. No DSN,
customer, claim token, provider response or free-form error enters OTLP. Export
uses cumulative metrics and a correlated log/span. An unknown delivery stops
new claims without retrying the observation or changing a committed refund.
The HTTP acknowledgement is not an fsync or business-commit acknowledgement.

The inherited ELITE_STOP_EVENT_HANDLE protocol is optional and Windows-only.
Only a trusted launcher may supply the manual-reset event. Native protocol
failure remains a host failure; the supervisor receipt owns process status.
This does not isolate hostile same-user workers or install a production service.

cmd/telemetry-reference supplies a finite synthetic load probe and ephemeral
certificate fixture for reference_telemetry/run_reference_runtime.py. It does
not call the database, provider or refund processor. Certificate keys are for
that short-lived reference only; the CA private key is never saved. The full
reference runner requires a hash-bound runtime lock, exact migration identities
and a fresh owned PostgreSQL cluster. See the Windows reference telemetry pack.
