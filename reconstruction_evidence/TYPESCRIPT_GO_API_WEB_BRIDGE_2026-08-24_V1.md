# TypeScript Web to Go API Bridge — Reconstruction Evidence V1

## Scope

- Base: optional TypeScript web V5 (64 files)
- Add-on: Go API bridge (6 files)
- Runtime: Node 24.14.1, pnpm 11.19.0, TypeScript 7.0.2, Next 16.3.2, Vitest 4.1.11

## Results

| Gate | Result |
|---|---:|
| 70-file clean composition and hashes | PASS |
| frozen offline dependency installation | PASS |
| strict TypeScript typecheck | PASS |
| 10 test files / 26 tests | PASS |
| fixed-path API client, unsafe URL and non-JSON negative tests | PASS |
| stable idempotency-key forwarding and replay parsing | PASS |
| Next production build and connected/BFF route generation | PASS |

The first typecheck rejected an explicitly `undefined` optional field under `exactOptionalPropertyTypes`. The BFF now omits the property when absent, and a second clean reconstruction passed all gates.

## Conditions

The public browser E2E is now recorded in `PUBLIC_WEB_GO_POSTGRES_BROWSER_E2E_2026-08-24_V1.md`. Protected admin/customer/factory journeys still require the selected OIDC/session and resource-authorization adapter. Production requires CSP/CSRF, distributed antiabuse, broader accessibility and performance evidence, and trusted edge/proxy policy.
