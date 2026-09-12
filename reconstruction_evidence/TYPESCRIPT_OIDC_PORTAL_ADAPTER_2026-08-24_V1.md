# TypeScript OIDC Portal Adapter — reconstruction evidence V1

Date: 2026-08-24  
Status: `REBUILD_VERIFIED / CONDITIONED`

## Result

The optional web profile now authenticates administrative, customer and factory portals with OpenID Connect Authorization Code, PKCE S256, state and nonce. The access token remains server-side inside a short encrypted JWE cookie; browser-visible session data contains only bounded identity and authorization claims.

The callback reconstructs its registered URL from `APP_BASE_URL` and only copies the authorization response query. It therefore does not trust an inbound proxy host when calculating `redirect_uri`. Production requires HTTPS and `__Host-` cookies; loopback HTTP and unprefixed cookies are allowed only outside production for local verification.

## Reconstruction

- `TS-OIDC-PORTAL-ADAPTER 0.1.0`: 10 files independently materializable with declared SHA-256 digests.
- `TS-ENTERPRISE-WEB 0.2.0`: administrative, customer and factory server components.
- `TS-GO-API-WEB-BRIDGE 0.1.0`: bounded public and protected Go API clients.
- compatible web profile: 3 packs, 80 implementation files plus `MATERIALIZATION_RECORD.md`.
- clean target: `${TEMP}/elite-web-profile-20260824-v6`.

## Gates

- offline frozen installation from the pinned lockfile: PASS, 92 packages reused and 0 downloaded.
- TypeScript `tsc --noEmit`: PASS.
- Vitest: 12 files and 31 tests PASS.
- Next.js 16.3.2 optimized production build: PASS, 13 static pages generated and protected routes server-rendered.
- tampered JWE and weak session secret rejection: PASS.
- access-token and cookie-size budget rejection: PASS.
- unsafe return target rejection and literal backend-route admission: PASS.
- real HTTP redirect chain against a local OIDC issuer: PASS.
- PKCE verifier, state, nonce, registered callback and client authentication at the token endpoint: PASS.
- `HttpOnly`, `SameSite=Lax`, path `/` and bounded session lifetime: PASS; `Secure`/`__Host-` are production-conditioned.
- authenticated `/admin`: PASS and exposed the seeded `order-portal` and `case-portal` records.
- authenticated `/customer`: PASS and exposed only the seeded subject-owned order and service case.
- authenticated `/factory`: PASS and exposed `PORTAL-UNIT-001` through its authorized organization.
- `/api/auth/session`: PASS and returned only subject, tenant, permissions and organizations; no access token.
- Go API re-verification and PostgreSQL reads during the portal flow: PASS.

## Recovery evidence

The first local exchange was rejected because the development server normalized the inbound host to `localhost` while the registered callback used `127.0.0.1`. The secure correction reconstructs the authorization response URL from the configured application base instead of trusting the request host. The flow was repeated from a new authorization code and passed.

A token-tampering test originally changed a Base64URL character whose unused trailing bits could decode to identical bytes. The test now mutates authenticated ciphertext in the middle and reliably proves rejection.

## Conditions

Production admission still requires the selected identity provider, registered HTTPS redirect/logout URLs, secret-manager delivery, key rotation policy, identity lifecycle and revocation tests, representative load evidence, CSP/header policy at the selected ingress, and provider/deployment observability. This adapter does not make TypeScript the backend base.
