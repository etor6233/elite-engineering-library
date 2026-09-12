# Go Electromobility Public CRM API — Reconstruction Evidence V2

## Change

Version 0.2.0 adds durable public-lead idempotency using `platform.idempotency_record`: required key, SHA-256 request binding, completed response/resource replay and conflict on same key with different payload.

## Results

| Gate | Result |
|---|---:|
| five-file rematerialization and hashes | PASS |
| first request creates lead, consent and outbox atomically | PASS |
| same key and same request returns the original lead without duplicate rows/events | PASS |
| same key with changed request hash returns conflict | PASS |
| HTTP requires a bounded idempotency key and exposes replay | PASS |
| clean 13-pack profile composition, all Go tests and vet | PASS |

Distributed antiabuse, CAPTCHA/risk decisions and edge rate limiting remain project conditions; idempotency is not a substitute for abuse prevention.
