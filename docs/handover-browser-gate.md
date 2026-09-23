# Connected handover browser gate

This gate proves the admitted handover operator/customer workflow through an actual Chromium page, Next BFF, Go API and PostgreSQL. It uses only generated local identity credentials and the existing official Stripe SDK HTTP fixture. It does not perform hosted IdP login, live payment, physical shipment or a fiscal effect.

## Required composition

Materialize the connected payment, initial handover, commercial release, handover context/UI, identity session, Next and Playwright packs together. The Go test reuses `newConnectedRun` from `payment_connected_integration_test.go`; its database guard requires an isolated loopback database whose name starts with `elite_payment_connected_`. Apply the composition's migrations before running the test. Do not point the fixture at business data.

Use the admitted Node, Go, Playwright and Chromium versions recorded in the accompanying runtime/evidence locks. Existing dependency caches are sufficient; this gate has no dependency installer. Build the frontend with its standard Next CLI and all TypeScript checks enabled. For a node_modules junction outside the workspace, the checked configuration uses `next build --webpack`; Turbopack rejects that junction. This is a builder choice, not a type-check exception.

Set `ELITE_HANDOVER_BROWSER=1`, `ELITE_WEB_ROOT` to the materialized frontend root, `ELITE_NODE_BIN` to the absolute admitted Node executable, and `PAYMENT_CONNECTED_DB_URL` to the isolated fixture database. The reference run also uses `BUSINESS_CONFIG_FILE=business.example.json`, `NEXT_TELEMETRY_DISABLED=1`, `GOTOOLCHAIN=local`, `GOPROXY=off` and `GOSUMDB=off`.

Run the admitted Go executable directly:

```text
go test -mod=readonly ./internal/platform/postgres -run ^TestHandoverBrowserPostgres$ -count=1 -v -timeout=4m
```

The Go test starts the already-built Next server through the explicit Node executable, creates a loopback HTTPS fixture and invokes the existing Playwright CLI directly for `tests/handover-connected.spec.mjs`, project `chromium-desktop`. It supplies the base URL so Playwright's package-manager webServer path is not used. An omitted fixture flag/database causes an explicit skip; a release gate must reject skipped execution.

## Assertions and effects

- Real encrypted role sessions and RS256/JWKS bearer verification authorize operator/customer identities; foreign organization and unrelated customer attempts fail.
- The payment fixture executes official SDK calls and signed callback reconciliation before preparation. The policy is an explicitly activated, hash-bound materialized profile with exact tenant, organization and provider/account/connection bindings.
- The browser prepares a handover, loses the response after commit, recovers it by GET, completes the existing checklist and obtains the real customer's acceptance.
- The browser records a commercial release receipt, loses the response after commit, then recovers the immutable receipt after a signed refund callback has already placed the payment on hold.
- Current authorization is false both before and after the official GET refund reconciliation. Recovery never converts historical receipt existence into current authorization.
- Double-click preparation produces one POST. PostgreSQL contains one preparation, one acceptance, one commercial receipt and one release event. There are two customer acceptance HTTP requests by design: one forbidden stranger attempt followed by one successful customer request.
- The final run requires zero React page errors and no horizontal overflow at a 390 by 844 mobile viewport. This is one Chromium project with a mobile viewport check, not a four-browser matrix.
- The commercial receipt creates no stock or money journal mutation and does not claim physical dispatch.

The local self-signed HTTPS fixture is allowed only by the existing explicit loopback Playwright configuration. Production TLS verification is not changed. Fixture control endpoints exist only in Go test code and use a generated per-run token.

## Source classification and delta

All new composition code, tests and frontend glue are AUTHORED. The timestamp fix reuses the admitted server locale/timezone formatter; the BFF bounded reader reuses the existing surveys reader with a 4096-byte limit. New PageProps overloads satisfy Next's generated type contract while preserving zero-argument direct-call tests. No upstream business logic, dependency or license is added or attributed to a vendor by these changes.

The evidence manifest records source prehashes, final source hashes, build/test/browser receipts and retained failure logs. Runtime directories, node_modules, .next, videos and screenshots are evidence or generated artifacts; they are not product source files.
