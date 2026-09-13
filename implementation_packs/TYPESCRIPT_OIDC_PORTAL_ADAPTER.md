# TypeScript OIDC Portal Adapter

V306: corrección AUTHORED de emisión/consulta de cotización con identidad
retenida, Idempotency-Key y recuperación scoped sin reenvío. Pruebas y límites
en reconstruction_evidence/QUOTE_CREATION_RECOVERY_V306.md. No cambio de pricing,
upstreams ni promoción integral; FAIL457 conserva los otros comandos pendientes.

## 1. Metadata

```yaml
pack_id: "TS-OIDC-PORTAL-ADAPTER"
pack_version: "0.3.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Hash-bound portal lifecycle: official OIDC code/PKCE/refresh/revocation, encrypted random browser handle and bound token vault, PostgreSQL CAS/logout/revocation sweep and optional joined host. AUTHORED glue; no IdP/password implementation or live certification."
stacks: ["Node 24", "TypeScript 7", "Next 16", "React 19", "jose 6.2.10", "openid-client 6.8.5"]
compatible_with: ["TS-GO-API-WEB-BRIDGE 0.5.2", "GO-ENTERPRISE-QUERY-API 0.1.1"]
incompatible_with: ["edge runtimes without the required Node.js cryptography and server-only APIs"]
license_expression: "LicenseRef-Workspace-Owner AND MIT dependencies"
upstream_sources: ["https://nextjs.org/docs/app/guides/authentication", "https://nextjs.org/docs/app/guides/data-security", "https://github.com/panva/openid-client"]
verified_at: "2026-09-13"
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
- Default compatibility mode retains no refresh token. When the selected hash-bound lifecycle is explicitly enabled, refresh credentials remain only in the server-side encrypted PostgreSQL vault; browser receives a random encrypted handle. See the admitted lifecycle contract and exact Go composition.

## 4. Exact file manifest

```text
CREATE src/platform/auth/portal-access-review.connected.test.ts
CREATE src/app/api/internal/oidc/session-maintenance/route.ts
CREATE src/platform/auth/portal-lifecycle.connected.test.ts
CREATE src/platform/auth/portal-lifecycle.ts
CREATE src/platform/auth/portal-maintenance.ts
CREATE src/platform/auth/portal-profile.test.ts
CREATE src/platform/auth/portal-profile.ts
CREATE src/platform/auth/portal-protocol.test.ts
CREATE src/platform/auth/portal-protocol.ts
CREATE src/platform/auth/portal-session-bridge.ts
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
sha256: "f409de96a6db5aaf09ec9bae79234d883aa6563ad94687b9d81c8cff1ec32321"
variables: []
secrets_allowed: false
```

````typescript
import "server-only";
import { createHash } from "node:crypto";
import { cookies } from "next/headers";
import { EncryptJWT, jwtDecrypt, type JWTPayload } from "jose";
import { portalLifecycleEnabled, portalProfile, portalSecret } from "./portal-profile";
import { readPortalSession } from "./portal-lifecycle";

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
  const secret = portalLifecycleEnabled() ? portalSecret("OIDC_PORTAL_SESSION_KEY_FILE") : process.env.AUTH_SESSION_SECRET;
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
  if (portalLifecycleEnabled()) return readPortalSession(portalProfile(), token);
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
sha256: "4bf24f3334e3bc02d13d8ca0f12c1c2b9f84223bda27ff644e2530f5ef8c88b7"
variables: []
secrets_allowed: false
```

````typescript
import * as oidc from "openid-client";
import { portalLifecycleEnabled, portalProfile } from "./portal-profile";
import { portalConfiguration } from "./portal-protocol";

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

export function applicationBaseUrl() { return portalLifecycleEnabled() ? new URL(portalProfile().post_logout_url) : secureUrl("APP_BASE_URL"); }
export function callbackUrl() { return portalLifecycleEnabled() ? portalProfile().callback_url : new URL("/api/auth/callback", applicationBaseUrl()).href; }

export function oidcConfiguration() {
  if (portalLifecycleEnabled()) return portalConfiguration(portalProfile());
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
sha256: "4f9a6dcb7ccc7d9becd161fd2feb8b8951c5b61d804a18dbd74fdd2d2482b90c"
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

// AUTHORED binary transport for the existing catalog PNG owner; fixed BFF route.
export async function protectedPostPNG<T>(session:PortalSession,path:string,bytes:Uint8Array):Promise<T>{
 if(!path.startsWith("/v1/admin/catalog/media/")||path.includes("..")||bytes.byteLength<1||bytes.byteLength>1048576)throw new Error("invalid bounded media request");
 const response=await fetch(new URL(path,baseUrl()),{method:"POST",cache:"no-store",redirect:"error",signal:AbortSignal.timeout(5000),headers:{accept:"application/json",authorization:`Bearer ${session.accessToken}`,"content-type":"image/png"},body:new Uint8Array(bytes).buffer});
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
sha256: "f3f3d257c290d360ab8e4526befcdbf285fd7a9bad12135e60cb68ad43e3a98d"
variables: []
secrets_allowed: false
```

````typescript
import { NextResponse } from "next/server";
import { callbackUrl, oidcClient, oidcConfiguration } from "@/platform/auth/oidc-client";
import { cookieOptions, FLOW_COOKIE, sealFlow } from "@/platform/auth/session";
import { portalLifecycleEnabled, portalProfile } from "@/platform/auth/portal-profile";

const destinations = new Set(["/admin", "/customer", "/factory", "/franchise"]);

export async function GET(request: Request) {
  const requested = new URL(request.url).searchParams.get("return_to") ?? "/customer";
  const returnTo = destinations.has(requested) ? requested : "/customer";
  const configuration = await oidcConfiguration();
  const codeVerifier = oidcClient.randomPKCECodeVerifier();
  const state = oidcClient.randomState();
  const nonce = oidcClient.randomNonce();
  const redirect = oidcClient.buildAuthorizationUrl(configuration, {
    redirect_uri: callbackUrl(), scope: portalLifecycleEnabled() ? portalProfile().scopes.join(" ") : "openid profile email", response_type: "code",
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
sha256: "e8d696b1ff6f607b7484bfda8e71a626ad2a33ebc751d3b0bd9542d0c4ec7c24"
variables: []
secrets_allowed: false
```

````typescript
import { NextResponse } from "next/server";
import { applicationBaseUrl, callbackUrl, oidcClient, oidcConfiguration } from "@/platform/auth/oidc-client";
import { cookieOptions, FLOW_COOKIE, openFlow, sealSession, SESSION_COOKIE } from "@/platform/auth/session";
import { portalLifecycleEnabled, portalProfile } from "@/platform/auth/portal-profile";
import { createPortalSession } from "@/platform/auth/portal-lifecycle";

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
    if (portalLifecycleEnabled()) {
      const stored = await createPortalSession(portalProfile(), flow, authorizationResponse);
      const response = NextResponse.redirect(new URL(flow.returnTo, applicationBaseUrl()), 303);
      response.cookies.set(SESSION_COOKIE, stored.cookie, cookieOptions(stored.maxAge));
      response.cookies.set(FLOW_COOKIE, "", cookieOptions(0));
      return response;
    }
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
  } catch {
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
sha256: "f971b6cd0da6f0c27731c5fa26b8cdee103be8a8606209cd0c8c8791b8a82ef5"
variables: []
secrets_allowed: false
```

````typescript
import { NextResponse } from "next/server";
import { applicationBaseUrl } from "@/platform/auth/oidc-client";
import { cookieOptions, SESSION_COOKIE } from "@/platform/auth/session";
import { portalLifecycleEnabled, portalProfile } from "@/platform/auth/portal-profile";
import { logoutPortalSession } from "@/platform/auth/portal-lifecycle";

export async function POST(request: Request) {
  const origin = request.headers.get("origin");
  if (origin !== applicationBaseUrl().origin) return NextResponse.json({ code: "ORIGIN_REJECTED" }, { status: 403 });
  if (portalLifecycleEnabled()) {
    const cookie=request.headers.get("cookie")?.split(";").map(v=>v.trim()).find(v=>v.startsWith(SESSION_COOKIE+"="))?.slice(SESSION_COOKIE.length+1);
    if(cookie){
      try {
        const result=await logoutPortalSession(portalProfile(),decodeURIComponent(cookie));
        if(!result.providerRevocationComplete) return NextResponse.json({code:"PROVIDER_REVOCATION_PENDING",local_session_revoked:true,retry:"POST /api/auth/logout"},{status:503,headers:{"cache-control":"no-store"}});
        const response=NextResponse.redirect(result.redirect,303);response.cookies.set(SESSION_COOKIE,"",cookieOptions(0));return response;
      } catch { return NextResponse.json({code:"LOGOUT_NOT_CONFIRMED",retry:"POST /api/auth/logout"},{status:503,headers:{"cache-control":"no-store"}}); }
    }
  }
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

V402 composed delta: T2804 catalog role source/edit/review/publication transport reuses original model/price writers and catalog owner. Bounded PNG and text defaults retained. No new dependency. CATALOG_ROLE_AUTHORING_RELEASE_V402.md.

### FILE: `src/app/api/internal/oidc/session-maintenance/route.ts`

```yaml
block_id: "PORTAL-EXTENSION-9:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "7cb4be95905c5eb9061fd210811ec6291387b51df64d3f9254d7bda4218a5ab1"
variables: []
secrets_allowed: false
```

````typescript
import {NextResponse} from "next/server";
import {portalLifecycleEnabled,portalProfile} from "@/platform/auth/portal-profile";
import {maintainPortalSessions} from "@/platform/auth/portal-maintenance";
import {BackendProblem} from "@/platform/backend/public-client";
export const maxDuration=45;
export async function POST(request:Request){
 if(!portalLifecycleEnabled())return NextResponse.json({code:"OIDC_LIFECYCLE_DISABLED"},{status:404});
 if(request.body){const reader=request.body.getReader();try{const chunk=await reader.read();if(!chunk.done){void reader.cancel().catch(()=>{});return NextResponse.json({code:"EMPTY_BODY_REQUIRED"},{status:400})}}finally{reader.releaseLock()}}
 const bearer=request.headers.get("authorization")??"";if(!/^Bearer [^\s]{1,16384}$/.test(bearer))return NextResponse.json({code:"SERVICE_IDENTITY_REQUIRED"},{status:401});
 try{const result=await maintainPortalSessions(portalProfile(),bearer);return NextResponse.json(result,{status:result.pending?503:200,headers:{"cache-control":"no-store"}})}catch(error){return NextResponse.json({code:"SESSION_MAINTENANCE_UNAVAILABLE"},{status:error instanceof BackendProblem&&[401,403].includes(error.status)?error.status:503,headers:{"cache-control":"no-store"}})}
}
````

### FILE: `src/platform/auth/portal-lifecycle.connected.test.ts`

```yaml
block_id: "PORTAL-EXTENSION-9:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "9d53285c79e9cdaf71d6357d241e3b0f3b5f86619ceeca54e816513c142d0e19"
variables: []
secrets_allowed: false
```

````typescript
import {describe,it,expect,vi,afterEach} from "vitest";
import {GET as login} from "@/app/api/auth/login/route";
import {GET as callback} from "@/app/api/auth/callback/route";
import {POST as logout} from "@/app/api/auth/logout/route";
import {SESSION_COOKIE,FLOW_COOKIE,openFlow} from "./session";
import {portalProfile} from "./portal-profile";
import {readPortalSession,createPortalSession} from "./portal-lifecycle";
import {portalServiceAccessToken} from "./portal-protocol";
import {POST as maintenance} from "@/app/api/internal/oidc/session-maintenance/route";

const enabled=!!process.env.PORTAL_FIXTURE_CONTROL_URL;
const controls=process.env.PORTAL_FIXTURE_CONTROL_URL??"http://127.0.0.1:9";
async function control(value:unknown){const response=await fetch(controls+"/fixture/control",{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify(value)});expect(response.status).toBe(200)}
async function stats(){return await(await fetch(controls+"/fixture/stats")).json() as Record<string,number>}
async function begin(){const response=await login(new Request("http://127.0.0.1:4567/api/auth/login?return_to=/customer"));expect(response.status).toBe(303);const flow=response.cookies.get(FLOW_COOKIE)?.value;expect(flow).toBeTruthy();const authorize=await fetch(response.headers.get("location")!,{redirect:"manual"});expect(authorize.status).toBe(303);return{flow:flow!,url:authorize.headers.get("location")!}}
async function finish(flow:{flow:string;url:string}){return callback(new Request(flow.url,{headers:{cookie:`${FLOW_COOKIE}=${flow.flow}`}}))}
async function loginCookie(){const response=await finish(await begin());if(response.status!==303){const diagnostic=await begin();try{await createPortalSession(portalProfile(),await openFlow(diagnostic.flow),new URL(diagnostic.url))}catch(error){const e=error as {name?:string;code?:string;message?:string;cause?:{code?:string;message?:string}};throw new Error(JSON.stringify({name:e.name,code:e.code,message:/^OIDC_[A-Z_]+$/.test(e.message??"")?e.message:"SDK_FAILURE",causeCode:e.cause?.code,causeMessage:e.cause?.message?.replace(/[A-Za-z0-9_-]{60,}/g,"REDACTED"),stats:await stats()}))}}expect(response.status).toBe(303);const cookie=response.cookies.get(SESSION_COOKIE)?.value;expect(cookie).toBeTruthy();expect(cookie).not.toContain("refresh");return cookie!}
afterEach(()=>vi.useRealTimers());
describe.skipIf(!enabled)("official OIDC SDK → BFF route → authenticated Go bridge → PostgreSQL",()=>{
 it("performs code/PKCE login, one concurrent refresh, lost commit recovery and durable logout",async()=>{
  const p=portalProfile(),cookie=await loginCookie();const session=await readPortalSession(p,cookie);expect(session?.subject).toBe("person");expect(session?.permissions).toEqual(["customer:read"]);
  const before=await stats();await control({expire:true,mode:"commit_lost"});const sessions=await Promise.all(Array.from({length:20},()=>readPortalSession(p,cookie)));expect(sessions.every(s=>s?.subject==="person")).toBe(true);expect((await stats()).refresh_grants).toBe(before.refresh_grants!+1);
  const denied=await logout(new Request(p.post_logout_url+"api/auth/logout",{method:"POST",headers:{origin:"https://foreign.invalid",cookie:`${SESSION_COOKIE}=${cookie}`}}));expect(denied.status).toBe(403);expect(await readPortalSession(p,cookie)).not.toBeNull();
  const done=await logout(new Request(p.post_logout_url+"api/auth/logout",{method:"POST",headers:{origin:new URL(p.post_logout_url).origin,cookie:`${SESSION_COOKIE}=${cookie}`}}));expect(done.status).toBe(303);expect(done.cookies.get(SESSION_COOKIE)?.value).toBe("");expect(await readPortalSession(p,cookie)).toBeNull();expect((await stats()).revocations).toBeGreaterThanOrEqual(2);
 });
 it("never reuses consumed refresh after an ambiguous provider response",async()=>{
  const p=portalProfile(),cookie=await loginCookie();await control({expire:true,mode:"grant_lost"});const before=await stats();expect(await readPortalSession(p,cookie)).toBeNull();expect(await readPortalSession(p,cookie)).toBeNull();const after=await stats();expect(after.refresh_grants).toBe(before.refresh_grants!+1);expect(after.reauth_required).toBeGreaterThan(0);
 });
 it("invalidates locally before failed remote revocation and recovers explicitly",async()=>{
  const p=portalProfile(),cookie=await loginCookie();const request=()=>new Request(p.post_logout_url+"api/auth/logout",{method:"POST",headers:{origin:new URL(p.post_logout_url).origin,cookie:`${SESSION_COOKIE}=${cookie}`}});
  await control({mode:"revoke_unavailable"});const first=await logout(request());expect(first.status).toBe(503);expect(await first.json()).toMatchObject({local_session_revoked:true});expect(await readPortalSession(p,cookie)).toBeNull();await control({mode:""});expect((await logout(request())).status).toBe(303);
 });
 it("rejects bad state, replayed code, issuer/audience/expiry/nonce/tenant and excess permissions",async()=>{
  const original=await begin();const wrong=new URL(original.url);wrong.searchParams.set("state","foreign");expect((await finish({...original,url:wrong.href})).status).toBe(400);
  const response=await finish(original);expect(response.status).toBe(303);expect((await finish(original)).status).toBe(400);
  for(const mode of ["wrong_issuer","wrong_audience","expired","wrong_nonce","foreign_tenant","extra_permission"]){await control({mode});expect((await finish(await begin())).status).toBe(400)}await control({mode:""});
 });
 it("rejects unsigned/tampered handles and accepts SDK JWKS rotation after its cache interval",async()=>{
  const p=portalProfile(),cookie=await loginCookie();expect(await readPortalSession(p,cookie.slice(0,-4)+"bad!")).toBeNull();const before=await stats();await control({expire:true,rotate:true});vi.useFakeTimers({toFake:["Date"]});vi.setSystemTime(Date.now()+61000);expect((await readPortalSession(p,cookie))?.subject).toBe("person");expect((await stats()).jwks).toBeGreaterThan(before.jwks!);
 });
 it("rejects a valid ciphertext swapped between two sessions of the same profile",async()=>{
  const p=portalProfile(),first=await loginCookie(),second=await loginCookie();expect(await readPortalSession(p,first)).not.toBeNull();await control({swap:true});expect(await readPortalSession(p,first)).toBeNull();expect(await readPortalSession(p,second)).toBeNull();await control({swap:true});expect(await readPortalSession(p,first)).not.toBeNull();
 });
 it("recovers provider revocation without cookies and sweeps bounded retention with durable unconfirmed counts",async()=>{
  const p=portalProfile();await loginCookie();await control({expire_absolute:true});const request=(bearer:string)=>new Request(p.post_logout_url+"api/internal/oidc/session-maintenance",{method:"POST",headers:{authorization:bearer}});
  expect((await maintenance(request("Bearer invalid"))).status).toBe(401);
  const bearer="Bearer "+await portalServiceAccessToken(p);let confirmed=0;
  for(let i=0;i<8;i++){const response=await maintenance(request(bearer));expect(response.status).toBe(200);const value=await response.json() as {claimed:number;confirmed:number};expect(value.claimed).toBeLessThanOrEqual(2);confirmed+=value.confirmed;if(!value.claimed)break}expect(confirmed).toBeGreaterThan(0);
  await loginCookie();await control({retention:true});const response=await maintenance(request(bearer));expect(response.status).toBe(200);expect(await response.json()).toMatchObject({unconfirmed_purged:1});
 });
});
````

### FILE: `src/platform/auth/portal-lifecycle.ts`

```yaml
block_id: "PORTAL-EXTENSION-9:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "1826dc6450e90e912fad6c391276da9c8313eb7d54a4e912bd02463dace73a3e"
variables: []
secrets_allowed: false
```

````typescript
import "server-only";
import { hkdfSync,randomBytes,createHash } from "node:crypto";
import { EncryptJWT,jwtDecrypt,type JWTPayload } from "jose";
import { z } from "zod";
import { portalSecret,type PortalProfile } from "./portal-profile";
import { portalConfiguration,portalOIDC } from "./portal-protocol";
import { portalSessionBridge,type PortalStoredSession } from "./portal-session-bridge";
import type { PortalSession,AuthFlow } from "./session";

const atom=z.string().min(1).max(200);const values=z.array(atom).max(100).refine(a=>new Set(a).size===a.length);
const tokenSchema=z.object({subject:atom,tenantId:atom,permissions:values,organizations:values.min(1),accessToken:z.string().min(1).max(16384),refreshToken:z.string().min(1).max(8192),idToken:z.string().min(1).max(16384),absoluteExpires:z.number().int().positive()});
type StoredTokens=z.infer<typeof tokenSchema>;
const flows=new Map<string,Promise<PortalSession|null>>();
function key(p:PortalProfile,purpose:string){const secret=portalSecret("OIDC_PORTAL_SESSION_KEY_FILE");if(secret.length<32)throw new Error("OIDC_SESSION_KEY_REJECTED");return new Uint8Array(hkdfSync("sha256",secret,p.documentSHA256,purpose,32));}
async function seal(p:PortalProfile,purpose:string,payload:JWTPayload,expires:number){let encrypted=new EncryptJWT(payload).setProtectedHeader({alg:"dir",enc:"A256GCM",typ:"JWT"}).setIssuer(p.profile_id).setAudience(purpose).setIssuedAt();if(purpose!=="portal-token-vault-v1")encrypted=encrypted.setExpirationTime(expires);return encrypted.encrypt(key(p,purpose));}
async function open(p:PortalProfile,purpose:string,value:string){return (await jwtDecrypt(value,key(p,purpose),{issuer:p.profile_id,audience:purpose,keyManagementAlgorithms:["dir"],contentEncryptionAlgorithms:["A256GCM"]})).payload;}
function safeSession(t:StoredTokens):PortalSession{return{subject:t.subject,tenantId:t.tenantId,permissions:t.permissions,organizations:t.organizations,accessToken:t.accessToken}}
function checkedTokens(p:PortalProfile,result:Awaited<ReturnType<typeof portalOIDC.authorizationCodeGrant>>,absolute:number,previous?:StoredTokens):{tokens:StoredTokens;expiry:string}{
 const claims=result.claims();if(!claims||!result.id_token||!result.refresh_token||!result.expires_in||result.expires_in<1||result.expires_in>3600||claims.iss!==p.issuer||claims.aud!==p.client_id||claims.tenant_id!==p.tenant_id||typeof claims.exp!=="number"||claims.exp*1000<=Date.now()||typeof claims.iat!=="number"||claims.iat*1000>Date.now()+30000)throw new Error("OIDC_IDENTITY_REJECTED");
 const tokens=tokenSchema.parse({subject:claims.sub,tenantId:claims.tenant_id,permissions:claims.permissions,organizations:claims.organization_ids,accessToken:result.access_token,refreshToken:result.refresh_token,idToken:result.id_token,absoluteExpires:absolute});
 if(tokens.organizations.some(o=>!p.organization_ids.includes(o))||tokens.permissions.some(permission=>!p.allowed_permissions.includes(permission))||(previous&&(tokens.subject!==previous.subject||tokens.tenantId!==previous.tenantId||tokens.refreshToken===previous.refreshToken)))throw new Error("OIDC_IDENTITY_REJECTED");
 const expiry=Math.min(Date.now()+result.expires_in*1000,claims.exp*1000,absolute*1000);if(expiry<=Date.now())throw new Error("OIDC_SESSION_EXPIRED");return{tokens,expiry:new Date(expiry).toISOString()};
}
async function encryptTokens(p:PortalProfile,sid:string,t:StoredTokens){return seal(p,"portal-token-vault-v1",{...t,sessionID:sid,profile:p.documentSHA256},t.absoluteExpires)}
async function decryptTokens(p:PortalProfile,sid:string,row:PortalStoredSession){const value=await open(p,"portal-token-vault-v1",row.ciphertext);if(value.sessionID!==sid||value.profile!==p.documentSHA256)throw new Error("OIDC_SESSION_BINDING_REJECTED");return tokenSchema.parse(value)}
export async function createPortalSession(p:PortalProfile,flow:AuthFlow,currentURL:URL):Promise<{cookie:string;maxAge:number}>{
 const c=await portalConfiguration(p);const result=await portalOIDC.authorizationCodeGrant(c,currentURL,{pkceCodeVerifier:flow.codeVerifier,expectedState:flow.state,expectedNonce:flow.nonce});
 const absolute=Math.floor(Date.now()/1000)+p.maximum_session_seconds;const checked=checkedTokens(p,result,absolute);const sid=randomBytes(32).toString("base64url");const ciphertext=await encryptTokens(p,sid,checked.tokens);
 const command={session_id:sid,ciphertext,access_expires_at:checked.expiry,absolute_expires_at:new Date(absolute*1000).toISOString()};
 try{await portalSessionBridge(p,"create",command)}catch{const recovered=await portalSessionBridge(p,"read",{session_id:sid});if(recovered.state!=="active"||recovered.version!==1||recovered.ciphertext!==ciphertext)throw new Error("OIDC_SESSION_STORE_UNAVAILABLE")}
 return{cookie:await seal(p,"portal-browser-handle-v1",{sessionID:sid,profile:p.documentSHA256},absolute),maxAge:p.maximum_session_seconds};
}
async function sessionID(p:PortalProfile,cookie:string){const value=await open(p,"portal-browser-handle-v1",cookie);if(typeof value.sessionID!=="string"||!/^[A-Za-z0-9_-]{43}$/.test(value.sessionID)||value.profile!==p.documentSHA256)throw new Error("OIDC_SESSION_REJECTED");return value.sessionID}
async function revokeStored(p:PortalProfile,sid:string,row:PortalStoredSession):Promise<boolean>{
 if(!row.revocation_pending)return true;
  try{const tokens=await decryptTokens(p,sid,row);const c=await portalConfiguration(p);await portalOIDC.tokenRevocation(c,tokens.refreshToken,{token_type_hint:"refresh_token"});await portalOIDC.tokenRevocation(c,tokens.accessToken,{token_type_hint:"access_token"});await portalSessionBridge(p,"ack-revocation",{session_id:sid,version:row.version,operation_id:row.operation_id});return true}catch{return false}
}
async function resolve(p:PortalProfile,sid:string):Promise<PortalSession|null>{
 let row=await portalSessionBridge(p,"read",{session_id:sid});if(row.state!=="active"||Date.parse(row.absolute_expires_at)<=Date.now())return null;
 let tokens=await decryptTokens(p,sid,row);if(tokens.absoluteExpires*1000<=Date.now())return null;if(Date.parse(row.access_expires_at)>Date.now()+p.refresh_before_seconds*1000)return safeSession(tokens);
 const operation=randomBytes(32).toString("base64url"),version=row.version;
 try{row=await portalSessionBridge(p,"claim",{session_id:sid,version,operation_id:operation})}catch{
  row=await portalSessionBridge(p,"read",{session_id:sid});
  if(row.state==="active"&&row.version>version)return safeSession(await decryptTokens(p,sid,row));
  if(row.state!=="refreshing"||row.version!==version||row.operation_id!==operation)return null;
 }
 if(row.state!=="refreshing"||row.operation_id!==operation||row.version!==version)return null;
 try{
  const result=await portalOIDC.refreshTokenGrant(await portalConfiguration(p),tokens.refreshToken);
  const checked=checkedTokens(p,result,tokens.absoluteExpires,tokens);tokens=checked.tokens;const ciphertext=await encryptTokens(p,sid,tokens);
  try{row=await portalSessionBridge(p,"commit",{session_id:sid,version,operation_id:operation,ciphertext,access_expires_at:checked.expiry})}catch{
   row=await portalSessionBridge(p,"read",{session_id:sid});
   if(row.version!==version+1||row.operation_id!==operation||row.ciphertext!==ciphertext)throw new Error("OIDC_REFRESH_COMMIT_UNKNOWN");
  }
  if(row.state!=="active"){await revokeStored(p,sid,row);return null}
  return safeSession(tokens);
 }catch{
  try{await portalSessionBridge(p,"abort",{session_id:sid,version,operation_id:operation})}catch{/* Durable refreshing state already rejects a second grant. */}
  return null;
 }
}
export async function readPortalSession(p:PortalProfile,cookie:string):Promise<PortalSession|null>{
 try{const sid=await sessionID(p,cookie);const key=p.documentSHA256+":"+sid;let pending=flows.get(key);if(!pending){pending=resolve(p,sid).catch(()=>null);flows.set(key,pending);void pending.finally(()=>flows.delete(key))}return await pending}catch{return null}
}
export async function logoutPortalSession(p:PortalProfile,cookie:string):Promise<{revoked:boolean;providerRevocationComplete:boolean;redirect:string}>{
 const sid=await sessionID(p,cookie);let row:PortalStoredSession;
 try{row=await portalSessionBridge(p,"revoke",{session_id:sid})}catch{row=await portalSessionBridge(p,"read",{session_id:sid});if(row.state!=="revoked")throw new Error("OIDC_LOGOUT_STORE_UNAVAILABLE")}
 const done=await revokeStored(p,sid,row);const c=await portalConfiguration(p);const redirect=portalOIDC.buildEndSessionUrl(c,{client_id:p.client_id,post_logout_redirect_uri:p.post_logout_url}).toString();return{revoked:true,providerRevocationComplete:done,redirect};
}
// A worker may open only the vault whose authenticated SID hashes to the leased row.
export async function revokePortalSweepItem(p:PortalProfile,row:PortalStoredSession,expectedSessionHash:string):Promise<boolean>{
 try{const value=await open(p,"portal-token-vault-v1",row.ciphertext);if(typeof value.sessionID!=="string"||!/^[A-Za-z0-9_-]{43}$/.test(value.sessionID)||value.profile!==p.documentSHA256||createHash("sha256").update(value.sessionID).digest("hex")!==expectedSessionHash)return false;return revokeStored(p,value.sessionID,row)}catch{return false}
}
````

### FILE: `src/platform/auth/portal-maintenance.ts`

```yaml
block_id: "PORTAL-EXTENSION-9:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "ffe715308a01aa64a0bffcf14beda5eb7cb8925205c758e01eb20149c0181143"
variables: []
secrets_allowed: false
```

````typescript
import "server-only";
import {randomBytes} from "node:crypto";
import {z} from "zod";
import {readBackendResponse} from "@/platform/backend/public-client";
import type {PortalProfile} from "./portal-profile";
import {revokePortalSweepItem} from "./portal-lifecycle";
const item=z.object({session_id_sha256:z.string().regex(/^[a-f0-9]{64}$/),version:z.number().int().positive(),state:z.literal("revoked"),operation_id:z.string().max(43),ciphertext:z.string().max(32768),access_expires_at:z.string().datetime({offset:true}),absolute_expires_at:z.string().datetime({offset:true}),revocation_pending:z.literal(true)}).strict();
const result=z.object({items:z.array(item).max(2),purged:z.number().int().min(0).max(100),unconfirmed_purged:z.number().int().min(0).max(100)}).strict();
export async function maintainPortalSessions(p:PortalProfile,serviceBearer:string){
 if(!/^Bearer [^\s]{1,16384}$/.test(serviceBearer))throw new Error("PORTAL_SERVICE_REQUIRED");
 // Go verifies this inbound service bearer before exposing/claiming any record.
 const response=await fetch(new URL("/v1/private/portal-sessions/sweep",p.bridge_url),{method:"POST",redirect:"error",cache:"no-store",signal:AbortSignal.timeout(5000),headers:{authorization:serviceBearer,"content-type":"application/json","X-Portal-Profile-SHA256":p.documentSHA256},body:JSON.stringify({operation_id:randomBytes(32).toString("base64url")})});
 const batch=result.parse(await readBackendResponse<unknown>(response));const completed=await Promise.all(batch.items.map(row=>revokePortalSweepItem(p,row,row.session_id_sha256)));
 return {claimed:batch.items.length,confirmed:completed.filter(Boolean).length,pending:completed.filter(v=>!v).length,purged:batch.purged,unconfirmed_purged:batch.unconfirmed_purged};
}
````

### FILE: `src/platform/auth/portal-profile.test.ts`

```yaml
block_id: "PORTAL-EXTENSION-9:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "ab23be4ca0cca49c4402fba1c028fb9e650b2ee38bafbc470e934c24abca91ee"
variables: []
secrets_allowed: false
```

````typescript
import {describe,it,expect} from "vitest";
import {createHash} from "node:crypto";
import {parsePortalProfile} from "./portal-profile";
function fixture(){return{schema:"elite.oidc.portal-lifecycle.v1",profile_id:"fixture",revision:1,transport:"TLS",issuer:"https://issuer.invalid",authorization_endpoint:"https://issuer.invalid/authorize",token_endpoint:"https://issuer.invalid/token",jwks_endpoint:"https://issuer.invalid/jwks",revocation_endpoint:"https://issuer.invalid/revoke",end_session_endpoint:"https://issuer.invalid/logout",client_id:"portal",callback_url:"https://portal.invalid/api/auth/callback",post_logout_url:"https://portal.invalid/",bridge_url:"https://backend.invalid",backend_audience:"api",service_client_id:"service",service_subject:"service",service_scopes:["portal-session:manage"],service_audience_parameter:"api",tenant_id:"tenant",organization_ids:["org"],allowed_permissions:["customer:read"],scopes:["openid","offline_access"],maximum_session_seconds:3600,refresh_before_seconds:30,retention_seconds:86400}}
function parse(value:unknown){const bytes=Buffer.from(JSON.stringify(value));return parsePortalProfile(bytes,createHash("sha256").update(bytes).digest("hex"))}
describe("hash-bound portal profile",()=>{
 it("accepts explicit TLS profile and prevents mutation",()=>{const p=parse(fixture());expect(p.tenant_id).toBe("tenant");expect(()=>p.organization_ids.push("foreign")).toThrow()});
 it.each([['transport',{transport:"LOOPBACK_FIXTURE"}],['HTTP',{issuer:"http://issuer.invalid"}],['secret URL',{token_endpoint:"https://user:secret@issuer.invalid/token"}],['wildcard',{allowed_permissions:["*"]}],['duplicate',{organization_ids:["org","org"]}],['missing refresh',{scopes:["openid"]}],['callback origin',{post_logout_url:"https://foreign.invalid/"}],['bridge path',{bridge_url:"https://backend.invalid/path"}],['audience',{service_audience_parameter:"foreign"}],['refresh margin',{refresh_before_seconds:300,maximum_session_seconds:300}],['unknown',{unknown:true}]])("rejects %s",(_name,change)=>{expect(()=>parse({...fixture(),...(change as object)})).toThrow()});
 it("rejects changed bytes and duplicate root names",()=>{const bytes=Buffer.from(JSON.stringify(fixture()));expect(()=>parsePortalProfile(bytes,"0".repeat(64))).toThrow();const duplicate=Buffer.from('{"profile_id":"other",'+bytes.toString().slice(1));expect(()=>parsePortalProfile(duplicate,createHash("sha256").update(duplicate).digest("hex"))).toThrow()});
});
````

### FILE: `src/platform/auth/portal-profile.ts`

```yaml
block_id: "PORTAL-EXTENSION-9:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "17f05c39141020c2434bf758b445b77fc00313df3c67fa1972fa39de24f2a2e6"
variables: []
secrets_allowed: false
```

````typescript
import "server-only";
import { createHash } from "node:crypto";
import { openSync, closeSync, readSync, fstatSync } from "node:fs";
import { z } from "zod";

const atom = z.string().min(1).max(200).regex(/^[^\s\x00]+$/).refine(v => v !== "*");
const set = z.array(atom).min(1).max(100).refine(v => new Set(v).size === v.length);
const schema = z.object({
  schema: z.literal("elite.oidc.portal-lifecycle.v1"), profile_id: atom, revision: z.literal(1), transport: z.enum(["TLS", "LOOPBACK_FIXTURE"]),
  issuer: z.string(), authorization_endpoint: z.string(), token_endpoint: z.string(), jwks_endpoint: z.string(), revocation_endpoint: z.string(), end_session_endpoint: z.string(),
  client_id: atom, callback_url: z.string(), post_logout_url: z.string(), bridge_url: z.string(), backend_audience: atom,
  service_client_id: atom, service_subject: atom, service_scopes: set, service_audience_parameter: z.string(),
  tenant_id: atom, organization_ids: set, allowed_permissions: set, scopes: set,
  maximum_session_seconds: z.number().int().min(300).max(86400), refresh_before_seconds: z.number().int().min(5).max(300),
  retention_seconds:z.number().int().min(3600).max(604800),
}).strict();
export type PortalProfile = z.infer<typeof schema> & { documentSHA256: string };
export function boundedFile(path: string, max: number): Buffer {
  const fd = openSync(path, "r");
  try { if (!fstatSync(fd).isFile()) throw new Error("OIDC_CONFIGURATION_REJECTED"); const bytes = Buffer.alloc(max + 1); let used = 0; while (used < bytes.length) { const n = readSync(fd, bytes, used, bytes.length - used, null); if (!n) break; used += n; } if (used > max) throw new Error("OIDC_CONFIGURATION_REJECTED"); return bytes.subarray(0, used); }
  finally { closeSync(fd); }
}
export function portalSecret(name: string): string {
  const path = process.env[name]; if (!path) throw new Error("OIDC_CONFIGURATION_REJECTED");
  const value = new TextDecoder("utf-8", { fatal: true }).decode(boundedFile(path, 4096)).replace(/\r?\n$/, "");
  if (!value || /[\x00\r\n]/.test(value)) throw new Error("OIDC_CONFIGURATION_REJECTED"); return value;
}
export function parsePortalProfile(raw: Uint8Array, digest: string): PortalProfile {
  if (!raw.length || raw.length > 16384 || !/^[a-f0-9]{64}$/.test(digest) || createHash("sha256").update(raw).digest("hex") !== digest) throw new Error("OIDC_CONFIGURATION_REJECTED");
  const text = new TextDecoder("utf-8", { fatal: true }).decode(raw);
  // Schema contains only scalar/array members. JSON token walk rejects duplicate root names.
  const keys: string[] = []; let depth = 0; const tokens = text.match(/"(?:[^"\\]|\\.)*"|[{}\[\]:,]|[^\s{}\[\]:,]+/g) ?? [];
  for (let i=0;i<tokens.length;i++) { const token=tokens[i]; if (token==="{"||token==="[") depth++; else if(token==="}"||token==="]") depth--; else if(depth===1 && token?.startsWith('"') && tokens[i+1]===":") keys.push(JSON.parse(token) as string); }
  if (new Set(keys).size !== keys.length) throw new Error("OIDC_CONFIGURATION_REJECTED");
  const p=schema.parse(JSON.parse(text));
  for (const rawURL of [p.issuer,p.authorization_endpoint,p.token_endpoint,p.jwks_endpoint,p.revocation_endpoint,p.end_session_endpoint,p.callback_url,p.post_logout_url,p.bridge_url]) {
    const u=new URL(rawURL); const local=["127.0.0.1","[::1]"].includes(u.hostname);
    if(u.username||u.password||u.search||u.hash||rawURL.length>2048||!(p.transport==="TLS" ? u.protocol==="https:" : local&&["http:","https:"].includes(u.protocol))) throw new Error("OIDC_CONFIGURATION_REJECTED");
  }
  if(!p.callback_url.endsWith("/api/auth/callback")||p.callback_url.slice(0,-18)+"/"!==p.post_logout_url||p.bridge_url.endsWith("/")||new URL(p.bridge_url).pathname!=="/"||!p.scopes.includes("openid")||!p.scopes.includes("offline_access")||p.refresh_before_seconds*2>=p.maximum_session_seconds||(p.service_audience_parameter!==""&&p.service_audience_parameter!==p.backend_audience)) throw new Error("OIDC_CONFIGURATION_REJECTED");
  return Object.freeze({...p, organization_ids:Object.freeze([...p.organization_ids]) as unknown as string[], allowed_permissions:Object.freeze([...p.allowed_permissions]) as unknown as string[], scopes:Object.freeze([...p.scopes]) as unknown as string[], service_scopes:Object.freeze([...p.service_scopes]) as unknown as string[],documentSHA256:digest});
}
export function portalLifecycleEnabled() { const value=process.env.OIDC_PORTAL_LIFECYCLE_ENABLED; if(value!==undefined&&value!=="false"&&value!=="true")throw new Error("OIDC_CONFIGURATION_REJECTED");return value==="true"; }
export function portalProfile(): PortalProfile {
  if(!portalLifecycleEnabled())throw new Error("OIDC_LIFECYCLE_DISABLED");
  const path=process.env.OIDC_PORTAL_PROFILE_FILE, hash=process.env.OIDC_PORTAL_PROFILE_SHA256;
  if(!path||!hash)throw new Error("OIDC_CONFIGURATION_REJECTED");
  const p=parsePortalProfile(boundedFile(path,16384),hash);
  if(process.env.NODE_ENV==="production"&&p.transport!=="TLS")throw new Error("OIDC_CONFIGURATION_REJECTED");
  if(process.env.APP_BASE_URL&&new URL(process.env.APP_BASE_URL).origin!==new URL(p.callback_url).origin)throw new Error("OIDC_CONFIGURATION_REJECTED");
  return p;
}
````

### FILE: `src/platform/auth/portal-protocol.test.ts`

```yaml
block_id: "PORTAL-EXTENSION-9:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "a8e8e668d3c560bced9ad1e543f523436efb194c2c25ada24f908ac756cea8ae"
variables: []
secrets_allowed: false
```

````typescript
import {describe,it,expect,vi,afterEach} from "vitest";
import {mkdtempSync,writeFileSync,rmSync} from "node:fs";
import {tmpdir} from "node:os";
import {join} from "node:path";
import {createHash} from "node:crypto";
import {parsePortalProfile} from "./portal-profile";
import {portalConfiguration} from "./portal-protocol";
let serial=0;
const directories:string[]=[];
afterEach(()=>{vi.unstubAllGlobals();vi.unstubAllEnvs();for(const directory of directories.splice(0))rmSync(directory,{recursive:true,force:true})});
function configured(){const directory=mkdtempSync(join(tmpdir(),"elite-portal-protocol-fixture-"));directories.push(directory);const secret=join(directory,"synthetic-secret");writeFileSync(secret,"synthetic-fixture-credential");vi.stubEnv("OIDC_PORTAL_CLIENT_SECRET_FILE",secret);const example={"schema":"elite.oidc.portal-lifecycle.v1","profile_id":"reference-portal","revision":1,"transport":"TLS","issuer":"https://issuer.invalid","authorization_endpoint":"https://issuer.invalid/authorize","token_endpoint":"https://issuer.invalid/token","jwks_endpoint":"https://issuer.invalid/jwks","revocation_endpoint":"https://issuer.invalid/revoke","end_session_endpoint":"https://issuer.invalid/logout","client_id":"reference-portal","callback_url":"https://portal.invalid/api/auth/callback","post_logout_url":"https://portal.invalid/","bridge_url":"https://backend.invalid","backend_audience":"reference-api","service_client_id":"reference-portal-maintenance","service_subject":"reference-portal-maintenance","service_scopes":["portal-session:manage"],"service_audience_parameter":"reference-api","tenant_id":"00000000-0000-4000-8000-000000000001","organization_ids":["00000000-0000-4000-8000-000000000002"],"allowed_permissions":["customer:read"],"scopes":["openid","offline_access"],"maximum_session_seconds":3600,"refresh_before_seconds":30,"retention_seconds":86400} as Record<string,unknown>;const bytes=Buffer.from(JSON.stringify({...example,profile_id:"protocol-boundary-"+(++serial)}));return parsePortalProfile(bytes,createHash("sha256").update(bytes).digest("hex"))}
describe("official SDK transport response boundary",()=>{
 it("cancels a streamed discovery body over262144bytes with one request",async()=>{const p=configured();let cancelled=false,calls=0;vi.stubGlobal("fetch",vi.fn(async()=>{calls++;return new Response(new ReadableStream<Uint8Array>({start(controller){controller.enqueue(new Uint8Array(262145))},cancel(){cancelled=true}}),{headers:{"content-type":"application/json"}})}));await expect(portalConfiguration(p)).rejects.toThrow();expect(calls).toBe(1);expect(cancelled).toBe(true)});
 it("passes the bounded valid discovery through the official SDK and refuses metadata endpoint drift",async()=>{const p=configured();const metadata={issuer:p.issuer,authorization_endpoint:p.authorization_endpoint,token_endpoint:p.token_endpoint,jwks_uri:p.jwks_endpoint,revocation_endpoint:p.revocation_endpoint,end_session_endpoint:p.end_session_endpoint,response_types_supported:["code"],subject_types_supported:["public"],id_token_signing_alg_values_supported:["RS256"],token_endpoint_auth_methods_supported:["client_secret_basic"]};const seen:string[]=[];vi.stubGlobal("fetch",vi.fn(async(input:unknown,init:RequestInit)=>{seen.push(String(input));expect(init.redirect).toBe("error");expect(init.signal).toBeDefined();return new Response(JSON.stringify(metadata),{headers:{"content-type":"application/json"}})}));expect((await portalConfiguration(p)).serverMetadata().issuer).toBe(p.issuer);expect(seen).toEqual([p.issuer+"/.well-known/openid-configuration"]);const changed=configured();vi.stubGlobal("fetch",vi.fn(async()=>new Response(JSON.stringify({...metadata,token_endpoint:"https://foreign.invalid/token"}),{headers:{"content-type":"application/json"}})));await expect(portalConfiguration(changed)).rejects.toThrow("OIDC_METADATA_REJECTED")});
});
````

### FILE: `src/platform/auth/portal-protocol.ts`

```yaml
block_id: "PORTAL-EXTENSION-9:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "8ee530b6ae65f7a24370071ee89c96284d27c4a18de4822131d374ffa5242ee2"
variables: []
secrets_allowed: false
```

````typescript
import "server-only";
import * as oidc from "openid-client";
import { createHash } from "node:crypto";
import { portalSecret, type PortalProfile } from "./portal-profile";

const configurations=new Map<string,Promise<oidc.Configuration>>();
function boundedFetch(p:PortalProfile): oidc.CustomFetch {
 const endpoints=new Set([p.issuer.replace(/\/$/,"")+"/.well-known/openid-configuration",p.token_endpoint,p.jwks_endpoint,p.revocation_endpoint]);
 return async(input,init)=>{
  const url=String(input);if(!endpoints.has(url))throw new Error("OIDC_ENDPOINT_REJECTED");
  const body=init.body instanceof Uint8Array?new Uint8Array(init.body):init.body;
  const response=await fetch(input,{method:init.method,headers:init.headers,...(body!==undefined?{body}:{}),...(init.duplex?{duplex:init.duplex}:{}),redirect:"error",signal:AbortSignal.any([AbortSignal.timeout(10000),...(init.signal?[init.signal]:[])])});
  if(!response.body)return response;
  const reader=response.body.getReader();let received=0;const chunks:Uint8Array[]=[];let complete=false;
  try{while(true){const chunk=await reader.read();if(chunk.done){complete=true;break}received+=chunk.value.byteLength;if(received>262144)throw new Error("OIDC_RESPONSE_TOO_LARGE");chunks.push(chunk.value)};const bytes=new Uint8Array(received);let position=0;for(const chunk of chunks){bytes.set(chunk,position);position+=chunk.byteLength};return new Response(bytes,{status:response.status,statusText:response.statusText,headers:response.headers});}
  finally{if(!complete)void reader.cancel().catch(()=>{});reader.releaseLock()}
 };
}
export async function portalConfiguration(p:PortalProfile,service=false):Promise<oidc.Configuration>{
 const secret=portalSecret(service?"OIDC_PORTAL_SERVICE_CLIENT_SECRET_FILE":"OIDC_PORTAL_CLIENT_SECRET_FILE");
 const cacheKey=p.documentSHA256+":"+service+":"+createHash("sha256").update(secret).digest("hex");
 let config=configurations.get(cacheKey);
 if(!config){if(configurations.size>=2)configurations.clear();config=(async()=>{
  const c=await oidc.discovery(new URL(p.issuer),service?p.service_client_id:p.client_id,{id_token_signed_response_alg:"RS256"},oidc.ClientSecretBasic(secret),{[oidc.customFetch]:boundedFetch(p),timeout:10,execute:[...(p.transport==="LOOPBACK_FIXTURE"?[oidc.allowInsecureRequests]:[]),oidc.enableNonRepudiationChecks]});
  const m=c.serverMetadata();
  if(m.issuer!==p.issuer||m.token_endpoint!==p.token_endpoint||m.jwks_uri!==p.jwks_endpoint||m.authorization_endpoint!==p.authorization_endpoint||m.revocation_endpoint!==p.revocation_endpoint||m.end_session_endpoint!==p.end_session_endpoint||!m.id_token_signing_alg_values_supported?.includes("RS256")||(m.token_endpoint_auth_methods_supported&&!m.token_endpoint_auth_methods_supported.includes("client_secret_basic")))throw new Error("OIDC_METADATA_REJECTED");
  return c;
 })();configurations.set(cacheKey,config);void config.catch(()=>configurations.delete(cacheKey));}
 return config;
}
let serviceToken:{key:string;value:string;expires:number}|undefined;
const serviceFlights=new Map<string,Promise<string>>();
export async function portalServiceAccessToken(p:PortalProfile):Promise<string>{
 if(serviceToken?.key===p.documentSHA256&&serviceToken.expires>Date.now()+10000)return serviceToken.value;
 const existing=serviceFlights.get(p.documentSHA256);if(existing)return existing;
 const flight=(async()=>{serviceToken=undefined;try{const config=await portalConfiguration(p,true);const result=await oidc.clientCredentialsGrant(config,{scope:p.service_scopes.join(" "),...(p.service_audience_parameter?{audience:p.service_audience_parameter}:{})});if(!result.access_token||result.access_token.length>16384||result.token_type.toLowerCase()!=="bearer"||!result.expires_in||result.expires_in<20||result.expires_in>3600)throw new Error("OIDC_SERVICE_TOKEN_REJECTED");serviceToken={key:p.documentSHA256,value:result.access_token,expires:Date.now()+result.expires_in*1000};return result.access_token;}catch{throw new Error("OIDC_SERVICE_UNAVAILABLE")}finally{serviceFlights.delete(p.documentSHA256)}})();serviceFlights.set(p.documentSHA256,flight);return flight;
}
export {oidc as portalOIDC};
````

### FILE: `src/platform/auth/portal-session-bridge.ts`

```yaml
block_id: "PORTAL-EXTENSION-9:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "dea8594a2ce386046e1da92f4f38647b3556851eefeb4b896b86a3a56197e9a9"
variables: []
secrets_allowed: false
```

````typescript
import "server-only";
import { z } from "zod";
import { readBackendResponse } from "@/platform/backend/public-client";
import { portalServiceAccessToken } from "./portal-protocol";
import type { PortalProfile } from "./portal-profile";
const record=z.object({version:z.number().int().positive(),state:z.enum(["active","refreshing","reauth_required","revoked"]),operation_id:z.string().max(43),ciphertext:z.string().max(32768),access_expires_at:z.string().datetime({offset:true}),absolute_expires_at:z.string().datetime({offset:true}),revocation_pending:z.boolean()}).strict();
export type PortalStoredSession=z.infer<typeof record>;
export type PortalSessionAction="create"|"read"|"claim"|"commit"|"abort"|"revoke"|"ack-revocation";
export type PortalSessionCommand={session_id:string;operation_id?:string;version?:number;ciphertext?:string;access_expires_at?:string;absolute_expires_at?:string};
export async function portalSessionBridge(p:PortalProfile,action:PortalSessionAction,command:PortalSessionCommand):Promise<PortalStoredSession>{
 if(!/^[a-zA-Z0-9_-]{43}$/.test(command.session_id))throw new Error("PORTAL_SESSION_REJECTED");
 const body=JSON.stringify(command);if(Buffer.byteLength(body)>40000)throw new Error("PORTAL_SESSION_REJECTED");
 const response=await fetch(new URL("/v1/private/portal-sessions/"+action,p.bridge_url),{method:"POST",cache:"no-store",redirect:"error",signal:AbortSignal.timeout(5000),headers:{authorization:`Bearer ${await portalServiceAccessToken(p)}`,"content-type":"application/json","X-Portal-Profile-SHA256":p.documentSHA256},body});
 return record.parse(await readBackendResponse<unknown>(response));
}
````


V402 composed delta: Portal lifecycle optional host/BFF integration; IDENTITY_PORTAL_RELEASE_V402.md/json; no change to closed business journeys.

### FILE: `src/platform/auth/portal-access-review.connected.test.ts`

```yaml
block_id: "J5-EXTENSION-TYPESCRIPT_OIDC_PORTAL_ADAPTER:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "1b9a92995af387e975ecdf5498dd6b1b239a89a3dc22179b9cd7ea4f49252c75"
variables: []
secrets_allowed: false
```

````typescript
import {describe,it,expect,vi,afterEach} from "vitest";
import {GET as login} from "@/app/api/auth/login/route";
import {GET as callback} from "@/app/api/auth/callback/route";
import {POST as logout} from "@/app/api/auth/logout/route";
import {SESSION_COOKIE,FLOW_COOKIE,openFlow} from "./session";
import {portalProfile} from "./portal-profile";
import {readPortalSession,createPortalSession} from "./portal-lifecycle";
import {portalServiceAccessToken} from "./portal-protocol";
import {POST as maintenance} from "@/app/api/internal/oidc/session-maintenance/route";

const enabled=!!process.env.PORTAL_FIXTURE_CONTROL_URL;
const controls=process.env.PORTAL_FIXTURE_CONTROL_URL??"http://127.0.0.1:9";
async function control(value:unknown){const response=await fetch(controls+"/fixture/control",{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify(value)});expect(response.status).toBe(200)}
async function stats(){return await(await fetch(controls+"/fixture/stats")).json() as Record<string,number>}
async function begin(){const response=await login(new Request("http://127.0.0.1:4567/api/auth/login?return_to=/customer"));expect(response.status).toBe(303);const flow=response.cookies.get(FLOW_COOKIE)?.value;expect(flow).toBeTruthy();const authorize=await fetch(response.headers.get("location")!,{redirect:"manual"});expect(authorize.status).toBe(303);return{flow:flow!,url:authorize.headers.get("location")!}}
async function finish(flow:{flow:string;url:string}){return callback(new Request(flow.url,{headers:{cookie:`${FLOW_COOKIE}=${flow.flow}`}}))}
async function loginCookie(){const response=await finish(await begin());if(response.status!==303){const diagnostic=await begin();try{await createPortalSession(portalProfile(),await openFlow(diagnostic.flow),new URL(diagnostic.url))}catch(error){const e=error as {name?:string;code?:string;message?:string;cause?:{code?:string;message?:string}};throw new Error(JSON.stringify({name:e.name,code:e.code,message:/^OIDC_[A-Z_]+$/.test(e.message??"")?e.message:"SDK_FAILURE",causeCode:e.cause?.code,causeMessage:e.cause?.message?.replace(/[A-Za-z0-9_-]{60,}/g,"REDACTED"),stats:await stats()}))}}expect(response.status).toBe(303);const cookie=response.cookies.get(SESSION_COOKIE)?.value;expect(cookie).toBeTruthy();expect(cookie).not.toContain("refresh");return cookie!}
afterEach(()=>vi.useRealTimers());
describe.skipIf(!enabled)("IdP access review through official SDK and durable portal",()=>{
 it("applies an explicit permission withdrawal on refresh without granting a replacement role",async()=>{const p=portalProfile(),cookie=await loginCookie();expect((await readPortalSession(p,cookie))?.permissions).toEqual(["customer:read"]);await control({expire:true,mode:"permission_withdrawn"});const changed=await readPortalSession(p,cookie);expect(changed?.subject).toBe("person");expect(changed?.permissions).toEqual([]);await control({mode:""})});
 it("rejects a changed subject during refresh and never links it to the original session",async()=>{const p=portalProfile(),cookie=await loginCookie();await control({expire:true,mode:"subject_changed"});expect(await readPortalSession(p,cookie)).toBeNull();expect(await readPortalSession(p,cookie)).toBeNull();await control({mode:""})});
 it("requires reauthentication after IdP deprovisioning and rejects a fresh authorization code",async()=>{const p=portalProfile(),cookie=await loginCookie();await control({expire:true,mode:"account_disabled"});expect(await readPortalSession(p,cookie)).toBeNull();expect(await readPortalSession(p,cookie)).toBeNull();expect((await finish(await begin())).status).toBe(400);await control({mode:""})});
});
````


V402 J5 original SDK access-review fixture; IDENTITY_J5_RELEASE_V402.md/json.
