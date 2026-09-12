# Franchise authenticated commands V115

## Scope and ownership

V115 extends the V114 web profile without introducing another backend, CRM, pricing service, identity store or browser token path.

- `TS-OIDC-PORTAL-ADAPTER 0.2.1` adds a bounded server-only POST primitive beside its existing protected GET primitive.
- `TS-FRANCHISE-JOURNEY-PORTALS 0.2.0` adds one same-origin command route, its security/mapping tests and one operator component.
- The Go APIs remain the command authority for lead assignment, lifecycle transition and quotation creation.
- The browser never receives the access token and cannot send quotation price or currency. The Go repository resolves those values from the active price book.

All three added files are `AUTHORED`. Next.js, React, Zod, Microsoft BCApps and Google Online Boutique remain declared dependencies/references; no local line is attributed to those projects.

## Exact reconstruction

The current inventory is 75 packs / 759 materialization blocks / 36 profiles / 121 exact upstream sources. Provenance is 622 `AUTHORED`, 32 `ADAPTED` and 105 `VERBATIM`. The web plan composes 6 packs / 78 files.

The compositor rebuilt the web plan into a new empty directory. Node 24.14.1 and pnpm 11.19.0 performed a frozen offline install with zero downloads. The exact result passed:

- 9 Vitest files / 29 tests;
- Next route type generation;
- TypeScript 7 strict typecheck;
- Next.js 16.3.2 production build;
- Microsoft Playwright 1.62.1 runtime 4/4 and exact target 8/8 across Chromium desktop/mobile, Firefox and WebKit;
- Google Lighthouse 13.4.1 contract 5/5 and five exact target runs: performance median 0.99, accessibility 1.00, best practices 1.00 and SEO 1.00.

The new regressions prove that a cross-origin request is rejected before session lookup, an admitted command is bound to session organization and permission, assignment maps only to the existing Go endpoint, and a quotation command cannot carry client-provided price/currency.

## Limits

This closes the web control gap for lead assignment, allowed lead transitions and quote issuance. It does not prove a real IdP, role effects against a deployed backend, quote acceptance/order conversion, customer handover evidence UX, scheduling capacity, complete branch administration, royalties/settlement, accounting/tax/fiscal behavior, financing or production edge behavior.

The project remains `PROJECT_CONDITIONED`. Production still requires all target-bound provider, identity, edge, recovery, load, offensive-security, deploy/rollback and business-acceptance receipts on the exact packaged artefact.

## Result

`V115_REBUILD_VERIFIED / PROJECT_CONDITIONED`.
