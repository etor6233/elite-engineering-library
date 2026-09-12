# Franchise Availability Web Portals — V135

Date: 2026-08-30  
Decision: `REBUILD_VERIFIED / CONDITIONED`

## Exact claim

`TS-FRANCHISE-JOURNEY-PORTALS 0.5.0` exposes the existing Go journey owner through the canonical server-side OIDC BFF. It adds organization/resource availability administration, resource creation and skills, resource assignment, reason-bound appointment transitions and customer-owned future cancellation. It adds no database, token-in-browser path, recurrence engine, labor policy or customer identifier supplied by the browser.

All twelve pack files are local `AUTHORED` code governed by the official Microsoft Business Central calendar/absence, Next.js and Google SRE references in pack metadata. No file is represented as Microsoft, Google or Vercel product source.

## Reconstructed identity

- Pack SHA-256: `628efa9340a42daf1155dcd20a1cfe39996d197c12eb482bfcfbf24669ad798b`.
- The pack materialized 12/12 files from Markdown into an empty directory.
- The enterprise web plan composed 6 packs and 82 implementation files into an empty directory, plus `MATERIALIZATION_RECORD.md`.
- Frozen pnpm installation completed offline with zero downloads from the existing exact lock.

## Executed gates

1. Nine Vitest files and 34 tests passed.
2. `next typegen` and strict `tsc --noEmit` passed.
3. Next.js 16.3.2 production build passed and emitted `/customer/appointments`, `/customer/quotes`, `/franchise` and the same-origin command route.
4. The license policy suite passed two tests.
5. Strict Zod mapping tests proved organization scope, permission selection, server-owned-field rejection, reason-bound no-show mapping and customer cancellation without a browser-provided customer subject.
6. The final gates ran on the fresh 82-file Markdown composition, not the authoring tree.

## Preserved conditions

The project still must prove real IdP roles, backend effects, approved timezone/recurrence/workforce policy, CDN/WAF/CSRF behavior, assisted accessibility, browser acceptance of protected journeys, RUM/load, offensive security, deploy/rollback and business acceptance. V135 accelerates those journeys; it does not fabricate their production evidence.
