# TypeScript OIDC Portal Adapter

V306: corrección AUTHORED de emisión/consulta de cotización con identidad
retenida, Idempotency-Key y recuperación scoped sin reenvío. Pruebas y límites
en reconstruction_evidence/QUOTE_CREATION_RECOVERY_V306.md. No cambio de pricing,
upstreams ni promoción integral; FAIL457 conserva los otros comandos pendientes.

## 1. Metadata

```yaml
pack_id: "TS-OIDC-PORTAL-ADAPTER"
pack_version: "0.2.5"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Abre portales Next.js mediante Authorization Code + PKCE + state + nonce, sesión JWE HttpOnly y consultas server-only al backend Go."
stacks: ["Node 24", "TypeScript 7", "Next 16", "React 19", "jose 6.2.10", "openid-client 6.8.5"]
compatible_with: ["TS-GO-API-WEB-BRIDGE 0.5.2", "GO-ENTERPRISE-QUERY-API 0.1.1"]
incompatible_with: ["edge runtimes without the required Node.js cryptography and server-only APIs"]
license_expression: "LicenseRef-Workspace-Owner AND MIT dependencies"
upstream_sources: ["https://nextjs.org/docs/app/guides/authentication", "https://nextjs.org/docs/app/guides/data-security", "https://github.com/panva/openid-client"]
verified_at: "2026-09-07"
```

Adapter web opcional. TypeScript no gobierna el backend ni las reglas de negocio. El navegador nunca recibe el access token: vive en una sesión JWE corta, `HttpOnly`, `SameSite=Lax`, marcada `Secure` en producción. El BFF sólo llama rutas literales del backend y éste vuelve a verificar token, tenant, permiso, organización y subject.

## 2. Applicability

V262 preserves Go Problem Details status/code for protected GET and POST, while rejecting invalid media-type prefixes and successful Problem Details. Twelve additional isolated regressions cover both methods. The opt-in `connected appointment confirmation` test is launched by GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.9.2 against its disposable loopback API/PostgreSQL fixture: no mocked successful fetch or verifier. It proves the server-side client boundary, not browser cookies, Authorization Code login, live IdP or production. Fixture tokens are ephemeral test inputs, never committed or promoted to application credentials. Evidence: `reconstruction_evidence/APPOINTMENT_CONFIRMATION_BFF_POSTGRES_V262.md`.

Use on the standalone Next.js BFF when staff, factory and customer portals require OIDC Authorization Code, short encrypted server session and protected server-to-server calls to Go. Reject it for edge runtimes lacking required Node cryptography/server APIs or issuers that cannot satisfy the claim/redirect contract.

## 3. Architecture contract

- V269 delegates GET and POST response parsing to the existing public transport owner, eliminating duplicate readers. Decompressed response bytes are limited before decoding, malformed bodies cannot confirm writes, and errors never trigger an automatic retry. Rebuild together with TS-GO-API-WEB-BRIDGE 0.5.2; evidence `reconstruction_evidence/BFF_STREAM_RESPONSE_BUDGET_V269.md`. Historical TS-ENTERPRISE-WEB combinations are not revalidated by this release.

- OIDC discovery requires a safe issuer, fixed client ID/secret and registered callback.
- Every login generates PKCE S256 verifier/challenge, state and nonce; the callback validates all three.
- Required claims are `sub`, `tenant_id`, `organization_ids`, `permissions`; missing or malformed values fail closed.
- Flow state expires after 10 minutes. Portal session expiry never exceeds the access token lifetime or one hour.
- Session encryption uses `dir` + `A256GCM` with a key derived from a secret of at least 32 characters.
- Login return targets use a fixed allowlist. Logout is POST and requires same origin.
- Server-only modules are guarded by Next.js `server-only`; the browser never receives the access token.
- No refresh token is retained. Expiry causes reauthentication until a selected revocation/rotation design is admitted.

## 4. Exact file manifest

```text
CREATE src/platform/auth/session.ts
CREATE src/platform/auth/oidc-client.ts
CREATE src/platform/auth/session.test.ts
CREATE src/platform/backend/protected-client.ts
CREATE src/platform/backend/protected-client.test.ts
CREATE src/test/server-only.ts
CREATE src/app/api/auth/login/route.ts
CREATE src/app/api/auth/callback/route.ts
CREATE src/app/api/auth/logout/route.ts
CREATE src/app/api/auth/session/route.ts
```

## 5. Materialization blocks

### FILE: `src/platform/auth/session.ts`

```yaml
block_id: "TS-OIDC-PORTAL:src-platform-auth-session-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "b15257fddb724d1a7a86891b9d90754bf8b53dbe153063f35122820fd7c49f83"
variables: []
secrets_allowed: false
```

````typescript
import "server-only";
import { createHash } from "node:crypto";
import { cookies } from "next/headers";
import { EncryptJWT, jwtDecrypt, type JWTPayload } from "jose";

export const FLOW_COOKIE = process.env.NODE_ENV === "production" ? "__Host-elite_oidc_flow" : "elite_oidc_flow";
export const SESSION_COOKIE = process.env.NODE_ENV === "production" ? "__Host-elite_session" : "elite_session";

export type AuthFlow = {
  codeVerifier: string;
  state: string;
  nonce: string;
  returnTo: string;
};

export type PortalSession = {
  subject: string;
  tenantId: string;
  permissions: string[];
  organizations: string[];
  accessToken: string;
};

const MAX_ACCESS_TOKEN_CHARACTERS = 2_000;
const MAX_SESSION_COOKIE_CHARACTERS = 3_800;

function key() {
  const secret = process.env.AUTH_SESSION_SECRET;
  if (!secret || secret.length < 32) throw new Error("AUTH_SESSION_SECRET must contain at least 32 characters");
  return createHash("sha256").update(secret, "utf8").digest();
}

async function seal(payload: JWTPayload, ttlSeconds: number) {
  return new EncryptJWT(payload)
    .setProtectedHeader({ alg: "dir", enc: "A256GCM", typ: "JWT" })
    .setIssuedAt()
    .setExpirationTime(`${ttlSeconds}s`)
    .encrypt(key());
}

async function open(token: string): Promise<JWTPayload> {
  const result = await jwtDecrypt(token, key(), { keyManagementAlgorithms: ["dir"], contentEncryptionAlgorithms: ["A256GCM"] });
  return result.payload;
}

function strings(value: unknown): string[] | null {
  if (!Array.isArray(value) || value.length > 100 || value.some((item) => typeof item !== "string" || item.length < 1 || item.length > 200)) return null;
  return [...new Set(value as string[])];
}

export async function sealFlow(flow: AuthFlow) { return seal(flow, 600); }

export async function openFlow(token: string): Promise<AuthFlow> {
  const value = await open(token);
  if (typeof value.codeVerifier !== "string" || typeof value.state !== "string" || typeof value.nonce !== "string" || typeof value.returnTo !== "string") throw new Error("invalid OIDC flow cookie");
  return { codeVerifier: value.codeVerifier, state: value.state, nonce: value.nonce, returnTo: value.returnTo };
}

export async function sealSession(session: PortalSession, ttlSeconds: number) {
  if (ttlSeconds < 1 || ttlSeconds > 3600) throw new Error("invalid session lifetime");
  if (session.accessToken.length < 1 || session.accessToken.length > MAX_ACCESS_TOKEN_CHARACTERS) throw new Error("invalid access token size");
  const token = await seal(session, ttlSeconds);
  if (token.length > MAX_SESSION_COOKIE_CHARACTERS) throw new Error("portal session exceeds cookie budget");
  return token;
}

export async function openSession(token: string): Promise<PortalSession> {
  const value = await open(token);
  const permissions = strings(value.permissions);
  const organizations = strings(value.organizations);
  if (typeof value.subject !== "string" || typeof value.tenantId !== "string" || typeof value.accessToken !== "string" || !permissions || !organizations || organizations.length === 0) throw new Error("invalid portal session");
  return { subject: value.subject, tenantId: value.tenantId, accessToken: value.accessToken, permissions, organizations };
}

export async function readSession(): Promise<PortalSession | null> {
  const token = (await cookies()).get(SESSION_COOKIE)?.value;
  if (!token) return null;
  try { return await openSession(token); } catch { return null; }
}

export function allowed(session: PortalSession, permission: string) {
  return session.permissions.includes("*") || session.permissions.includes(permission);
}

export function cookieOptions(maxAge: number) {
  return { httpOnly: true, secure: process.env.NODE_ENV === "production", sameSite: "lax" as const, path: "/", maxAge };
}
````

### FILE: `src/platform/auth/oidc-client.ts`

```yaml
block_id: "TS-OIDC-PORTAL:src-platform-auth-oidc-client-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "6004dcacb7c828af14cc5fffe1a8fadc74e8e5e44b02f6517ae5d263297baa4e"
variables: []
secrets_allowed: false
```

````typescript
import * as oidc from "openid-client";

let cached: Promise<oidc.Configuration> | undefined;

function required(name: string) {
  const value = process.env[name];
  if (!value) throw new Error(`${name} is required`);
  return value;
}

function secureUrl(name: string) {
  const value = new URL(required(name));
  const localDevelopment = process.env.NODE_ENV !== "production" && ["localhost", "127.0.0.1"].includes(value.hostname);
  if (value.username || value.password || value.search || value.hash || (value.protocol !== "https:" && !localDevelopment)) throw new Error(`${name} must be a safe HTTPS URL`);
  return value;
}

function localInsecure(url: URL) {
  return process.env.NODE_ENV !== "production" && url.protocol === "http:" && ["localhost", "127.0.0.1"].includes(url.hostname);
}

export function applicationBaseUrl() { return secureUrl("APP_BASE_URL"); }
export function callbackUrl() { return new URL("/api/auth/callback", applicationBaseUrl()).href; }

export function oidcConfiguration() {
  const issuer = secureUrl("OIDC_ISSUER");
  cached ??= oidc.discovery(issuer, required("OIDC_CLIENT_ID"), required("OIDC_CLIENT_SECRET"), undefined, localInsecure(issuer) ? { execute: [oidc.allowInsecureRequests] } : undefined);
  return cached;
}

export const oidcClient = oidc;
````

### FILE: `src/platform/auth/session.test.ts`

```yaml
block_id: "TS-OIDC-PORTAL:src-platform-auth-session-test-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "1de790cc570d37a9b2d663ae1ac5108f8a88fc20cec0a8b938a4c57d9dda8120"
variables: []
secrets_allowed: false
```

````typescript
import { afterEach, describe, expect, it } from "vitest";
import { allowed, openFlow, openSession, sealFlow, sealSession } from "./session";

const original = process.env.AUTH_SESSION_SECRET;
afterEach(() => { process.env.AUTH_SESSION_SECRET = original; });

describe("encrypted portal session", () => {
  it("round trips bounded flow and session values", async () => {
    process.env.AUTH_SESSION_SECRET = "a-secure-test-secret-with-more-than-thirty-two-characters";
    const flow = { codeVerifier: "verifier", state: "state", nonce: "nonce", returnTo: "/admin" };
    await expect(openFlow(await sealFlow(flow))).resolves.toEqual(flow);
    const session = { subject: "user", tenantId: "tenant", permissions: ["admin:read"], organizations: ["org"], accessToken: "access" };
    const opened = await openSession(await sealSession(session, 60));
    expect(opened).toEqual(session);
    expect(allowed(opened, "admin:read")).toBe(true);
    expect(allowed(opened, "factory:read")).toBe(false);
  });

  it("rejects tampering and weak secrets", async () => {
    process.env.AUTH_SESSION_SECRET = "short";
    await expect(sealFlow({ codeVerifier: "v", state: "s", nonce: "n", returnTo: "/" })).rejects.toThrow(/32/);
    process.env.AUTH_SESSION_SECRET = "a-secure-test-secret-with-more-than-thirty-two-characters";
    const token = await sealFlow({ codeVerifier: "v", state: "s", nonce: "n", returnTo: "/" });
    const position = Math.floor(token.length / 2);
    const tampered = token.slice(0, position) + (token[position] === "a" ? "b" : "a") + token.slice(position + 1);
    await expect(openFlow(tampered)).rejects.toThrow();
  });

  it("rejects sessions that cannot fit safely in one cookie", async () => {
    process.env.AUTH_SESSION_SECRET = "a-secure-test-secret-with-more-than-thirty-two-characters";
    const session = { subject: "user", tenantId: "tenant", permissions: ["admin:read"], organizations: ["org"], accessToken: "x".repeat(2_001) };
    await expect(sealSession(session, 60)).rejects.toThrow(/access token size/);
  });
});
````

### FILE: `src/platform/backend/protected-client.ts`

```yaml
block_id: "TS-OIDC-PORTAL:src-platform-backend-protected-client-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "ab7531e3d6dc3a7311f546655abe4d89d482380392852b1accc976b923fa97ea"
variables: []
secrets_allowed: false
```

````typescript
import "server-only";
import { readBackendResponse } from "./public-client";
import type { PortalSession } from "@/platform/auth/session";

export type Overview = { organization_id: string; orders: number; open_leads: number; stock_available: number; open_cases: number; active_shipments: number };
export type OrderSummary = { id: string; organization_id: string; customer_subject: string; state: string; currency: string; total_minor_units: number; version: number };
export type FactoryUnitSummary = { id: string; organization_id: string; purchase_order_id: string; variant_id: string; serial_number: string; state: string };
export type ServiceCaseSummary = { id: string; organization_id: string; stock_unit_id: string; state: string; severity: string; description: string; version: number };
export type LeadSummary = { id: string; organization_id: string; model_id?: string; state: string; source_code: string; assigned_subject?: string; created_at: string; version: number };
export type AppointmentSummary = { id: string; organization_id: string; lead_id: string; model_id?: string; kind: string; starts_at: string; state: string; version: number };
export type QuoteSummary = { id: string; organization_id: string; lead_id: string; variant_id: string; price_book_id: string; currency: string; total_minor_units: number; valid_until: string; state: string; version: number };
export type HandoverSummary = { id: string; organization_id: string; order_id: string; stock_unit_id: string; state: string; version: number; customer_accepted_at?: string };
export type CustomerJourney = { appointments: AppointmentSummary[]; quotes: QuoteSummary[]; handovers: HandoverSummary[] };
export type Page<T> = { items: T[]; next_cursor?: string };

function baseUrl() {
  const raw = process.env.ENTERPRISE_API_BASE_URL;
  if (!raw || !URL.canParse(raw)) throw new Error("ENTERPRISE_API_BASE_URL is invalid");
  const value = new URL(raw);
  if (!['http:', 'https:'].includes(value.protocol) || value.username || value.password || value.search || value.hash) throw new Error("enterprise backend URL is unsafe");
  return value;
}

export async function protectedGet<T>(session: PortalSession, path: string, query: Record<string, string | undefined>, idempotencyKey?: string): Promise<T> {
  if (!path.startsWith("/v1/") || path.includes("..")) throw new Error("protected backend path is invalid");
  if (idempotencyKey !== undefined && !/^[A-Za-z0-9_-]{16,128}$/.test(idempotencyKey)) throw new Error("invalid idempotency key");
  const target = new URL(path, baseUrl());
  for (const [name, value] of Object.entries(query)) if (value) target.searchParams.set(name, value);
  const response = await fetch(target, { cache: "no-store", redirect: "error", signal: AbortSignal.timeout(5000), headers: { accept: "application/json", authorization: `Bearer ${session.accessToken}`, ...(idempotencyKey ? { "Idempotency-Key": idempotencyKey } : {}) } });
  return readBackendResponse<T>(response);
}

export async function protectedPost<T>(session: PortalSession, path: string, command: unknown, idempotencyKey?: string): Promise<T> {
  if (!path.startsWith("/v1/") || path.includes("..")) throw new Error("protected backend path is invalid");
  if (idempotencyKey !== undefined && !/^[A-Za-z0-9_-]{16,128}$/.test(idempotencyKey)) throw new Error("invalid idempotency key");
  const body = JSON.stringify(command);
  if (body.length > 1_048_576) throw new Error("protected backend command is too large");
  const response = await fetch(new URL(path, baseUrl()), {
    method: "POST",
    cache: "no-store",
    redirect: "error",
    signal: AbortSignal.timeout(5000),
    headers: { accept: "application/json", authorization: `Bearer ${session.accessToken}`, "content-type": "application/json", ...(idempotencyKey ? { "Idempotency-Key": idempotencyKey } : {}) },
    body
  });
  return readBackendResponse<T>(response);
}
````

### FILE: `src/platform/backend/protected-client.test.ts`

```yaml
block_id: "TS-OIDC-PORTAL:src-platform-backend-protected-client-test-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "5da9a69bf52d697b1411cde59502b821e95d48a987c81df1ec46a933dc2fc9c5"
variables: []
secrets_allowed: false
```

````typescript
import { afterEach, describe, expect, it, vi } from "vitest";
import { protectedGet, protectedPost } from "./protected-client";

const original = process.env.ENTERPRISE_API_BASE_URL;
afterEach(() => { process.env.ENTERPRISE_API_BASE_URL = original; vi.unstubAllGlobals(); });

it.skipIf(process.env.ELITE_CONFIRMATION_E2E !== "1")("connected appointment confirmation commits once and is readable only by its customer", async () => {
  const base = new URL(process.env.ENTERPRISE_API_BASE_URL!);
  if (base.protocol !== "http:" || base.hostname !== "127.0.0.1") throw new Error("local fixture API required");
  const id = process.env.ELITE_CONFIRMATION_ID!;
  const tenantId = process.env.ELITE_CONFIRMATION_TENANT!;
  if (!id || !tenantId) throw new Error("missing fixture identity");
  const operator = { subject: "operator", tenantId, permissions: ["appointment:manage"], organizations: ["store"], accessToken: process.env.ELITE_CONFIRMATION_OPERATOR! };
  const path = `/v1/franchise/appointments/${encodeURIComponent(id)}/transitions`;
  const command = { organization_id: "store", current: "requested", target: "confirmed", version: 2 };
  const results = await Promise.allSettled([protectedPost<{ id: string; state: string; version: number }>(operator, path, command), protectedPost<{ id: string; state: string; version: number }>(operator, path, command)]);
  const successes = results.filter((value) => value.status === "fulfilled");
  const rejected = results.filter((value) => value.status === "rejected");
  expect(successes).toHaveLength(1);
  expect(rejected).toHaveLength(1);
  expect(successes[0]!.value).toMatchObject({ id, state: "confirmed", version: 3 });
  expect(rejected[0]!.reason).toMatchObject({ status: 409, code: "JOURNEY_CONFLICT" });
  for (const [subject, token, expected] of [["customer", process.env.ELITE_CONFIRMATION_CUSTOMER!, 1], ["stranger", process.env.ELITE_CONFIRMATION_STRANGER!, 0]] as const) {
    const session = { subject, tenantId, permissions: ["customer:self"], organizations: ["store"], accessToken: token };
    const timeline = await protectedGet<{ appointments: { id: string; state: string; version: number }[] }>(session, "/v1/customer/journey", { organization_id: "store" });
    expect(timeline.appointments).toHaveLength(expected);
    if (expected) expect(timeline.appointments[0]).toMatchObject({ id, state: "confirmed", version: 3 });
  }
}, 20_000);

describe("protected backend client", () => {
  const session = { subject: "operator", tenantId: "tenant", permissions: ["appointment:manage"], organizations: ["store"], accessToken: "synthetic-token" };
  for (const method of ["GET", "POST"] as const) {
    it(`${method} rejects excess UTF-8 response bytes with a single fetch`, async () => {
      process.env.ENTERPRISE_API_BASE_URL = "https://api.example.test";
      const fetch = vi.fn(async () => new Response(JSON.stringify({ padding: "é".repeat(524288) }), { headers: { "content-type": "application/json" } }));
      vi.stubGlobal("fetch", fetch);
      const result = method === "GET" ? protectedGet(session, "/v1/customer/journey", {}) : protectedPost(session, "/v1/franchise/appointments/a/transitions", {});
      await expect(result).rejects.toMatchObject({ status: 502, code: "RESPONSE_TOO_LARGE" });
      expect(fetch).toHaveBeenCalledOnce();
    });
    for (const [status, code] of [[401, "UNAUTHENTICATED"], [403, "ORGANIZATION_FORBIDDEN"], [409, "JOURNEY_CONFLICT"]] as const) {
      it(`${method} preserves ${status} Problem Details from the real API contract`, async () => {
        process.env.ENTERPRISE_API_BASE_URL = "https://api.example.test";
        vi.stubGlobal("fetch", vi.fn(async () => new Response(JSON.stringify({ code }), { status, headers: { "content-type": "application/problem+json; charset=utf-8" } })));
        const result = method === "GET" ? protectedGet(session, "/v1/customer/journey", {}) : protectedPost(session, "/v1/franchise/appointments/a/transitions", {});
        await expect(result).rejects.toMatchObject({ status, code });
      });
    }
    for (const contentType of ["application/jsonp", "text/html", "application/problem+json"]) {
      it(`${method} rejects successful ${contentType} instead of confirming it`, async () => {
        process.env.ENTERPRISE_API_BASE_URL = "https://api.example.test";
        vi.stubGlobal("fetch", vi.fn(async () => new Response('{"state":"confirmed"}', { status: 200, headers: { "content-type": contentType } })));
        const result = method === "GET" ? protectedGet(session, "/v1/customer/journey", {}) : protectedPost(session, "/v1/franchise/appointments/a/transitions", {});
        await expect(result).rejects.toMatchObject({ status: 502, code: "INVALID_CONTENT_TYPE" });
      });
    }
  }
  it("keeps the access token server-side and forwards bounded query", async () => {
    process.env.ENTERPRISE_API_BASE_URL = "https://api.example.test";
    const fetchMock = vi.fn(async (_url: URL, init: RequestInit) => new Response('{"items":[]}', { status: 200, headers: { "content-type": "application/json" } }));
    vi.stubGlobal("fetch", fetchMock);
    const session = { subject: "u", tenantId: "t", permissions: ["customer:self"], organizations: ["o"], accessToken: "secret-access-token" };
    await protectedGet(session, "/v1/customer/orders", { organization_id: "o", limit: "25" });
    const [url, init] = fetchMock.mock.calls[0]!;
    expect(String(url)).toContain("organization_id=o");
    expect((init.headers as Record<string, string>).authorization).toBe("Bearer secret-access-token");
  });

  it("rejects unsafe paths", async () => {
    process.env.ENTERPRISE_API_BASE_URL = "https://api.example.test";
    const session = { subject: "u", tenantId: "t", permissions: [], organizations: ["o"], accessToken: "token" };
    await expect(protectedGet(session, "/v1/../admin", {})).rejects.toThrow(/path/);
  });
});

it("carries the exact quote identity and rejects invalid keys without sending",async()=>{
 process.env.ENTERPRISE_API_BASE_URL="https://api.example.test";
 const fetch=vi.fn(async(_input:RequestInfo|URL,_init?:RequestInit)=>new Response("{}",{headers:{"content-type":"application/json"}}));vi.stubGlobal("fetch",fetch);
 const session={subject:"operator",tenantId:"tenant",permissions:["quote:write"],organizations:["store"],accessToken:"synthetic-token"};
 await protectedPost(session,"/v1/franchise/quotes",{},"quote-request-key-0001");
 expect(fetch.mock.calls[0]?.[1]).toMatchObject({headers:{"Idempotency-Key":"quote-request-key-0001"}});
 for(const key of ["", "short", "a".repeat(129), "valid-length-but space"])await expect(protectedPost(session,"/v1/franchise/quotes",{},key)).rejects.toThrow();
 expect(fetch).toHaveBeenCalledOnce();
});
````

### FILE: `src/test/server-only.ts`

```yaml
block_id: "TS-OIDC-PORTAL:src-test-server-only-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "8e609bb71c20b858c77f0e9f90bb1319db8477b13f9f965f1a1e18524bf50881"
variables: []
secrets_allowed: false
```

````typescript
export {};
````

### FILE: `src/app/api/auth/login/route.ts`

```yaml
block_id: "TS-OIDC-PORTAL:src-app-api-auth-login-route-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "ad195223494dd0d2de770639802faae03f56de337e7b902565bf3d03250e4630"
variables: []
secrets_allowed: false
```

````typescript
import { NextResponse } from "next/server";
import { callbackUrl, oidcClient, oidcConfiguration } from "@/platform/auth/oidc-client";
import { cookieOptions, FLOW_COOKIE, sealFlow } from "@/platform/auth/session";

const destinations = new Set(["/admin", "/customer", "/factory", "/franchise"]);

export async function GET(request: Request) {
  const requested = new URL(request.url).searchParams.get("return_to") ?? "/customer";
  const returnTo = destinations.has(requested) ? requested : "/customer";
  const configuration = await oidcConfiguration();
  const codeVerifier = oidcClient.randomPKCECodeVerifier();
  const state = oidcClient.randomState();
  const nonce = oidcClient.randomNonce();
  const redirect = oidcClient.buildAuthorizationUrl(configuration, {
    redirect_uri: callbackUrl(), scope: "openid profile email", response_type: "code",
    code_challenge: await oidcClient.calculatePKCECodeChallenge(codeVerifier), code_challenge_method: "S256", state, nonce,
  });
  const response = NextResponse.redirect(redirect, 303);
  response.cookies.set(FLOW_COOKIE, await sealFlow({ codeVerifier, state, nonce, returnTo }), cookieOptions(600));
  return response;
}
````

### FILE: `src/app/api/auth/callback/route.ts`

```yaml
block_id: "TS-OIDC-PORTAL:src-app-api-auth-callback-route-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "ae9eb07b85fc4016b01c357e18b707a8b888e693983efad5b8b860d3bc818f85"
variables: []
secrets_allowed: false
```

````typescript
import { NextResponse } from "next/server";
import { applicationBaseUrl, callbackUrl, oidcClient, oidcConfiguration } from "@/platform/auth/oidc-client";
import { cookieOptions, FLOW_COOKIE, openFlow, sealSession, SESSION_COOKIE } from "@/platform/auth/session";

function stringArray(value: unknown): string[] | null {
  if (!Array.isArray(value) || value.length > 100 || value.some((item) => typeof item !== "string" || item.length < 1 || item.length > 200)) return null;
  return [...new Set(value as string[])];
}

export async function GET(request: Request) {
  const flowCookie = request.headers.get("cookie")?.split(";").map((part) => part.trim()).find((part) => part.startsWith(`${FLOW_COOKIE}=`))?.slice(FLOW_COOKIE.length + 1);
  if (!flowCookie) return NextResponse.json({ code: "OIDC_FLOW_MISSING" }, { status: 400 });
  try {
    const flow = await openFlow(decodeURIComponent(flowCookie));
    const configuration = await oidcConfiguration();
    const authorizationResponse = new URL(callbackUrl());
    authorizationResponse.search = new URL(request.url).search;
    const tokens = await oidcClient.authorizationCodeGrant(configuration, authorizationResponse, { pkceCodeVerifier: flow.codeVerifier, expectedState: flow.state, expectedNonce: flow.nonce });
    const claims = tokens.claims();
    const organizations = stringArray(claims?.organization_ids);
    const permissions = stringArray(claims?.permissions);
    if (!tokens.access_token || typeof claims?.sub !== "string" || typeof claims.tenant_id !== "string" || !organizations?.length || !permissions) throw new Error("required identity claims missing");
    const ttl = Math.max(1, Math.min(tokens.expires_in ?? 300, 3600));
    const response = NextResponse.redirect(new URL(flow.returnTo, applicationBaseUrl()), 303);
    response.cookies.set(SESSION_COOKIE, await sealSession({ subject: claims.sub, tenantId: claims.tenant_id, permissions, organizations, accessToken: tokens.access_token }, ttl), cookieOptions(ttl));
    response.cookies.set(FLOW_COOKIE, "", cookieOptions(0));
    return response;
  } catch (error) {
    console.error("OIDC callback rejected", { name: error instanceof Error ? error.name : "UnknownError", message: error instanceof Error ? error.message : "unknown failure" });
    const response = NextResponse.json({ code: "OIDC_CALLBACK_REJECTED" }, { status: 400 });
    response.cookies.set(FLOW_COOKIE, "", cookieOptions(0));
    return response;
  }
}
````

### FILE: `src/app/api/auth/logout/route.ts`

```yaml
block_id: "TS-OIDC-PORTAL:src-app-api-auth-logout-route-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "074fdd5ec0123b2123527121a001a8320fb0194f51f89160563b219f15ab848c"
variables: []
secrets_allowed: false
```

````typescript
import { NextResponse } from "next/server";
import { applicationBaseUrl } from "@/platform/auth/oidc-client";
import { cookieOptions, SESSION_COOKIE } from "@/platform/auth/session";

export async function POST(request: Request) {
  const origin = request.headers.get("origin");
  if (origin !== applicationBaseUrl().origin) return NextResponse.json({ code: "ORIGIN_REJECTED" }, { status: 403 });
  const response = NextResponse.redirect(new URL("/", applicationBaseUrl()), 303);
  response.cookies.set(SESSION_COOKIE, "", cookieOptions(0));
  return response;
}
````

### FILE: `src/app/api/auth/session/route.ts`

```yaml
block_id: "TS-OIDC-PORTAL:src-app-api-auth-session-route-ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "48c314b7c45fb1b407bf89b6425ee2b7183036342032825ee45be1d8bca3447d"
variables: []
secrets_allowed: false
```

````typescript
import { NextResponse } from "next/server";
import { readSession } from "@/platform/auth/session";

export async function GET() {
  const session = await readSession();
  if (!session) return NextResponse.json({ authenticated: false }, { status: 401, headers: { "cache-control": "no-store" } });
  return NextResponse.json({ authenticated: true, subject: session.subject, tenant_id: session.tenantId, permissions: session.permissions, organization_ids: session.organizations }, { headers: { "cache-control": "no-store" } });
}
````

## 6. Configuration surface

| Variable | Type/default | Secret | Validation/effect |
|---|---|---|---|
| OIDC issuer/client ID | HTTPS URL/non-empty / none | no | discovery and audience/client binding |
| OIDC client secret | minimum policy / none | yes | server-only external injection |
| session encryption secret | at least 32 chars / none | yes | derives JWE key; weak/missing fails |
| public BFF origin/callback | allowlisted URL / none | no | controls redirects and same-origin logout |
| backend URL | safe server URL / none | no | literal protected paths only |
| cookie/expiry settings | bounded by implementation | no | `HttpOnly`, `SameSite=Lax`, production `Secure` |

## 7. Dependency bill

| Dependency | Pin | Use | License | Scope | Source |
|---|---|---|---|---|---|
| Node.js | `24` baseline | BFF runtime | MIT | build/runtime | `nodejs.org` |
| Next/React/TypeScript | `16` / `19` / `7` baseline | portal/BFF | MIT/Apache-2.0 | build/runtime | official projects |
| `jose` | `6.2.10` | JWE/JWT crypto | MIT | runtime | `github.com/panva/jose` |
| `openid-client` | `6.8.5` | OIDC protocol | MIT | runtime | `github.com/panva/openid-client` |

## 8. Apply order

Compose over `TS-GO-API-WEB-BRIDGE 0.2.x`, inject target configuration outside source control, register exact issuer redirects, run frozen install/tests/typecheck/build, then execute role/session browser E2E against the Go backend. Existing auth must use an explicit migration/cookie cutover. Rollback invalidates the new session cookie/key and restores the prior compatible BFF artifact.

## 9. Verification

Frozen install, strict typecheck, all Vitest suites and Next production build must pass. Unit tests cover encrypted round trip, tamper rejection, weak-secret rejection, permission behavior, unsafe backend path and server-side bearer forwarding. Production additionally requires a selected issuer sandbox, registered redirect/logout URIs, key/secret rotation, revocation/logout policy, CSP nonce strategy, edge rate limits, browser E2E for each role and cookie behavior behind the selected TLS proxy.

## 10. Reconstruction evidence

The earlier cryptographic/authorization evidence remains applicable to unchanged session/OIDC code. Version 0.2.1 retains the journey response types and adds one bounded server-only POST primitive with the same URL, token, content-type, redirect, timeout and response-size constraints as protected reads. The integrated 2026-08-29 web run passed 29 tests, strict typecheck, production build, Playwright and Lighthouse. Real issuer, roles, sessions and revocation remain target conditions.

V402 composed delta: Connected handover browser/BFF/Go/PostgreSQL gate; bounded body, stable server date formatting, Next PageProps signatures. AUTHORED integration glue; existing domain and fixed upstreams unchanged.

V402 browser evidence: reconstruction_evidence/HANDOVER_BROWSER_V402.md. Production Webpack build includes TypeScript checking;42 direct-call and23 focused BFF/date tests PASS. Chromium desktop and mobile viewport through real BFF/API/PG cover lost-response recovery and callback invalidation. Hosted IdP/live payment/physical shipment not claimed.
