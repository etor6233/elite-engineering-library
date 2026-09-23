# Private display locale

AUTHORED presentation glue, not code attributed to a company. Existing fixed Next,
React, Node Intl and original session/BFF owners supply runtime behavior.

Authenticated private displays use per-user/tenant cookie, then bounded browser
language, then configured es/en. The preference is not a permission or business
policy. Saving explicitly reloads the page with a warning about unsaved forms.
Anonymous public pages retain their configured locale.

Money uses the original integer/BigInt algorithm and currency digits. Display
dates use the configured timezone. UTC form command semantics remain unchanged.
Wire status codes, IDs, audit reasons, content and assessment answers are data.

Twenty-one same-release guides translate only when id/version/title/paragraphs
match original source. Five course translations additionally require the exact
profile/content hashes and course fields. Unknown or historical content remains
literal with lang=und; the original training profile/bundle and evidence stay intact.
Tenant CMS content retains its explicit source language. English help search
filters only guides returned by the authorized endpoint, after the original
strict query parser; no permission or API query-boundary bypass.

Reference proof: private-locale-connected.spec.mjs, TestPrivateLocaleBrowser,
private-display.test.ts, private-locale.test.ts, locale/route.test.ts.
Use the existing fixed local runtimes and an owned migrated PostgreSQL fixture.
The connected browser tests a lost create response, explicit language switch,
GET-only recovery and update/publication: exactly three durable writes. No live
account is needed. Operational readiness and provider credentials are separate.
