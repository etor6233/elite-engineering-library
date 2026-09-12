# Public Web → Go → PostgreSQL Browser E2E — Evidence V1

## Runtime

- Browser: Codex in-app Chromium automation
- Web: Next 16.3.2 production server on loopback
- API: materialized `cmd/electromobility-api` on loopback
- Identity discovery: disposable local OIDC discovery server; no protected token used
- Database: fresh profile schema on PostgreSQL 18.6

## Journey

1. Seed one active tenant, public store and public vehicle model.
2. Render `/connected` through Next server-side data retrieval from the Go public catalog API.
3. Verify accessible DOM roles for heading, model, labels, checkbox, button and live status.
4. Submit synthetic name/email and explicit consent through the browser.
5. Next BFF forwards a stable `Idempotency-Key` to Go.
6. Go commits lead, consent evidence, idempotency record and outbox event.
7. UI displays the success live-region message and resets the form.
8. Database query proves exactly one consent, one event and one idempotency record for the lead.
9. Browser console contains zero warnings/errors after the corrected run.

## Error-recovery proof

The first browser run exposed a React async-event lifetime bug at `event.currentTarget.reset()`. The canonical bridge captured the form element before `await`, refreshed its hash, reconstructed all 70 web files, repeated offline install/typecheck/26 tests/build and reran the browser journey successfully.

## Remaining conditions

This closes the public catalog→lead reference slice. Protected admin/customer/factory journeys remain open until resource-scoped OIDC/session integration and read/query APIs are composed. Accessibility tooling, responsive/browser matrix, performance and antiabuse load remain separate gates.
