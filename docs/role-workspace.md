# Role workspace — V381

This pack supplies session-filtered navigation and four curated entry views, not a role grant,
financial dashboard, employee activation system or training curriculum. Selecting owner, admin,
employee or customer changes only the displayed subset of links. The existing verified session
and destination page/backend remain authoritative. Public catalog/help/locations stay public.

Compose TS-MULTIROLE-ONBOARDING0.1.8 with TS-GO-API-WEB-BRIDGE0.5.16,
TS-OIDC-PORTAL-ADAPTER0.2.x and TS-FRANCHISE-JOURNEY-PORTALS0.14.19. Reserve
/dashboard, /guide and /guide/[role] before composition. Then explicitly set
features.role_workspace=true in the consumer business configuration. The default absent/false
flag hides the menu and returns404 for all three routes, including wildcard sessions. Do not
enable it in a profile that does not select this pack. Existing module/feature flags still apply.

An empty authenticated permission set retains only common public sections and the guide. A guest
gets the public navigation and the customer login entry but no workspace. The private customer,
factory and admin links require customer:self, factory:read and admin:read respectively. Franchise
uses the actual destination's any-permission predicate: inventory:allocate, payment:create,
handover:manage, admin:read, lead:read, resource:manage, availability:read, availability:manage,
appointment:manage. Exact wildcard is supported; wrong case, quote:write alone and role labels
grant no private link. Features/modules constrain wildcard too. Hiding links is not an API ACL.

Operational instructions remain under the existing versioned /help owner. No unsupported payment
approval, royalty or high-value policy remains in the curated role labels. Reading a view does
not train a model, enroll an employee, change permissions or record learning progress.

Qualification:22 policy tests;236 full web tests pass plus one pre-existing connected skip;
typecheck/build;15 explicit permission profiles across five actual pages per profile in each of
four browsers (300 page states), with exact header/panel/view link sets. Guest/expired/tampered
sessions require the actual307/Location using the same cookie context with maxRedirects:0; no
IdP follows. Real encrypted cookies are synthetic fixture identities, not proof of external IdP.
Four disabled-flag tests return404. Unknown/prototype-like guide names return404, customer direct
/admin is denied, owner view cannot widen customer access, and links reach the existing help.
No browser writes, subject/tenant/token disclosure or mobile horizontal overflow observed.

Run test:workspace twice in the admitted Playwright gate: ELITE_WORKSPACE_E2E=enabled and
disabled, with owned HTTPS loopback production Next, ELITE_WEB_ROOT, synthetic AUTH_SESSION_SECRET
and the corresponding runtime config. Use the existing help-loopback-proxy fixture; never deploy
it. Retries are zero. Retain traces and stop both owned processes. No dependency lock changed.

Rollback: set features.role_workspace=false; routes become404 without state migration. When
removing the selected pack, leave the flag off and remove only its seven materialized files.
No database, provider, confidential document, external identity or training dataset is required
for this bounded reference qualification. Consumer authorization/IdP/revocation, business
onboarding, assessed progress, localization and target acceptance remain separate requirements.

V382 refreshes these compatible versions after private-read qualification. The V381 workspace implementation and its fifteen-profile policy remain byte-identical; see private-portal-reads.md for the new administrative/customer read and amount evidence.

V383 refreshes compatible selected versions for private-read recovery; the role source and permission policy remain unchanged.

V387 preserves role navigation and updates compatible portal/query/browser revisions for cursor traversal. No role or domain-policy change.

V388 aligns selected versions after extending cursor traversal to admin lists; no new role or grant.

V389 aligns selected versions after correcting customer appointment presentation to the existing business time zone. No role or grant changes.

V390 aligns selected versions for the existing customer appointment management link and durable cancellation recovery; roles and grants unchanged.

V393 aligns the composed BFF version after omission of its unused optional image optimizer dependency. Role implementation and permission policy remain unchanged.
