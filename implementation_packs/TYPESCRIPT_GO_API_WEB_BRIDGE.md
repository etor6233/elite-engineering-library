# Enterprise TypeScript Web BFF

## 1. Metadata

```yaml
pack_id: "TS-GO-API-WEB-BRIDGE"
pack_version: "0.7.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un frontend Next.js y BFF autónomo del golden path histórico, con web pública configurable, catálogo, leads y lectura acotada de capacidad de turnos, portales cliente/operaciones/fábrica, clientes acotados hacia APIs Go, primitivas Google SafeValues contra DOM XSS y CSP estricta con nonce por request evaluada por Google CSP Evaluator; no contiene dominio, SQL ni persistencia."
stacks: ["Node 24", "TypeScript 7", "Next 16", "React 19", "Vitest 4", "Google SafeValues 1.2.0", "Google CSP Evaluator 1.1.8"]
compatible_with: ["TS-OIDC-PORTAL-ADAPTER 0.2.4", "MICROSOFT-PLAYWRIGHT-BROWSER-GATE 0.1.37", "GO-ELECTROMOBILITY-PUBLIC-CRM-API 0.2.3", "GO-ENTERPRISE-QUERY-API 0.1.6", "GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.10.19"]
incompatible_with: ["TS-ENTERPRISE-WEB as a co-installed foundation", "direct browser access to protected backend bearer tokens", "business rules or persistence owned by TypeScript"]
license_expression: "LicenseRef-Workspace-Owner AND MIT AND Apache-2.0 (complete remaining dependency obligations remain gated)"
upstream_sources: ["https://nextjs.org/docs/app/guides/backend-for-frontend", "https://nextjs.org/docs/app/guides/data-security", "https://nextjs.org/docs/app/guides/content-security-policy", "https://github.com/google/safevalues/tree/0a900a2dc1d3dce6f28ab28aca8d0a5fd63b1c5b", "https://www.npmjs.com/package/safevalues/v/1.2.0", "https://github.com/google/csp-evaluator/tree/2e00e419d0be0c204a91ce08f7cc4cb838a2b75d", "https://www.npmjs.com/package/csp_evaluator/v/1.1.8"]
verified_at: "2026-09-11"
```

This is the selectable frontend/BFF foundation, never the backend or domain foundation. It has no runtime or build dependency on `TS-ENTERPRISE-WEB`; protected portals are composed with `TS-OIDC-PORTAL-ADAPTER 0.2.4`. V269 revalidates web compilation/tests, not the historical provider or browser gates.

Current V393 profile excludes the entire optional Sharp closure and disables runtime image optimization. The historical Sharp LGPL scope below does not describe an installed dependency of this profile. Remaining bundled/native dependency and redistribution gates still apply; re-admit an image pipeline before enabling it.

Historical V376 repairs the exact Next16.3.2/Sharp0.35.3 advisory exposure with Next16.3.4/Sharp0.35.4. Preserve the signed artifact receipts and full native/vendor scope limits in `reconstruction_evidence/NEXT_SHARP_SECURITY_REPAIR_V376.md`. The declared npm graph is not a complete native or vendored SBOM. Redistribution must include the exact Next MIT companion license for @next/env/SWC and resolve all Sharp/libvips native component notices/source obligations. Whole TEST03/07 and production admission remain blocked.

## 2. Applicability

Use when a project needs a configurable public web, lead capture and customer, administration and factory portals in front of admitted Go APIs. Select it with the OIDC overlay for protected journeys. Reject it when the frontend must own transactional invariants, write directly to a database, expose backend bearer tokens to the browser or coexist with the historical TypeScript full-stack golden path.

V377 adds actual Next metadata/robots/sitemap integration for the two existing public pages, a configured HTTPS canonical origin, inert escaped WebSite JSON-LD and default noindex for private pages. Indexing is opt-in; no domain, translated URL, model detail route, publication date, offer or Google admission is inferred. See `docs/public-indexing.md` and `reconstruction_evidence/PUBLIC_WEB_METADATA_INTEGRATION_V377.md`.

V378 integrates es/en catalogs, regional number/plural formatting and configured market timeZone through the actual public lead/appointment journey. Unsupported language falls back to Spanish with an honest lang; private portals remain Spanish. No translation CMS, locale negotiation, provider approval or general CLDR equivalence. See `docs/public-localization.md`.

## 3. Architecture contract

- Next.js owns presentation, accessibility baseline, configuration loading and narrow BFF routes; Go owns domain rules, authorization enforcement, workflows, persistence and integrations.
- The public backend client accepts only a safe base URL, fixed path templates, bounded JSON responses, five-second timeout and no redirects.
- V269 uses one shared response reader for public and protected calls: at most 1,048,576 decompressed body bytes retained before UTF-8 decoding, cancellation on overflow or rejected media type, bounded sanitized error codes and no retry. Native Node fetch tests prove chunked cancellation, gzip expansion and deadline during body reads. Domain schema validation remains with consumers; this is not a process RSS or network-buffer guarantee. Implementation is AUTHORED, not Microsoft code. Method: Microsoft HttpCompletionOption content timeout/buffer warning and WHATWG Streams read/cancel; evidence `reconstruction_evidence/BFF_STREAM_RESPONSE_BUDGET_V269.md`.
- Lead submission requires consent plus a stable idempotency key; the BFF forwards neither arbitrary upstream paths nor browser credentials.
- Protected portals call literal server-only routes supplied by the OIDC overlay. The Go API revalidates token, tenant, permission, organization and object scope.
- Business configuration is schema-versioned, bounded and read from a basename inside `config/`; invalid dependencies, identifiers, workflows or references fail closed.
- The frontend contains no SQL client, migrations, outbox, worker or source of truth for business state.
- React text interpolation remains the default. When code must cross a non-React DOM sink, the admitted primitive is Google SafeValues 1.2.0; raw HTML still requires explicit product policy, provenance and review.
- `src/proxy.ts` follows the official Next.js nonce pattern, overwrites untrusted inbound nonce/CSP values, forces dynamic rendering and emits a fresh strict CSP without production `unsafe-inline`/`unsafe-eval`; Google CSP Evaluator rejects the previous incomplete header and admits the generated production policy.
- The local four-browser nonce proof is mandatory, but a CDN/WAF may alter request headers or caching. Project gates still own HTTPS edge repetition, cache isolation, provider allowlists, antiabuse, TLS proxy trust, observability and production browser testing.

## 4. Exact file manifest

```text
CREATE src/platform/backend/public-client.ts
CREATE src/platform/backend/public-client.test.ts
CREATE src/platform/backend/response-boundary.test.ts
CREATE src/app/api/enterprise/models/route.ts
CREATE src/app/api/enterprise/leads/route.ts
CREATE src/app/connected/lead-form.tsx
CREATE src/app/connected/page.tsx
CREATE .env.example
CREATE .gitignore
CREATE config/business.example.json
CREATE next-env.d.ts
CREATE next.config.ts
CREATE package.json
CREATE pnpm-lock.yaml
CREATE pnpm-workspace.yaml
CREATE README.md
CREATE THIRD_PARTY_NOTICES.md
CREATE tsconfig.json
CREATE vitest.config.ts
CREATE src/app/admin/page.tsx
CREATE src/app/customer/page.tsx
CREATE src/app/factory/page.tsx
CREATE src/app/globals.css
CREATE src/app/icon.svg
CREATE src/app/layout.tsx
CREATE src/app/models/page.tsx
CREATE src/app/page.tsx
CREATE src/platform/config/load.ts
CREATE src/platform/config/registry.ts
CREATE src/platform/config/schema.test.ts
CREATE src/platform/config/schema.ts
CREATE src/platform/http/problem.test.ts
CREATE src/platform/http/problem.ts
CREATE src/platform/security/google-safevalues.test.ts
CREATE src/platform/security/google-safevalues.ts
CREATE src/platform/security/csp-policy.ts
CREATE src/platform/security/csp-policy.test.ts
CREATE src/proxy.ts
CREATE src/platform/seo/public-indexing.ts
CREATE src/platform/seo/public-indexing.test.ts
CREATE src/app/robots.ts
CREATE src/app/sitemap.ts
CREATE docs/public-indexing.md
CREATE src/platform/i18n/public-catalog.ts
CREATE src/platform/i18n/load-public-locale.ts
CREATE src/platform/i18n/public-catalog.test.ts
CREATE docs/public-localization.md
CREATE src/app/admin/page.test.ts
CREATE src/platform/i18n/money.ts
CREATE src/platform/i18n/money.test.ts
CREATE docs/private-portal-reads.md
CREATE src/platform/backend/private-read-failure.tsx
CREATE src/platform/backend/private-read-failure.test.ts
CREATE src/app/admin/error.tsx
CREATE src/app/customer/error.tsx
CREATE src/app/factory/error.tsx
CREATE src/platform/backend/portal-paging.tsx
CREATE src/app/customer/page.test.ts
CREATE docs/factory-operator-flow.md
CREATE microsoft_playwright_browser_gate/tests/factory-connected.spec.mjs
CREATE src/app/api/enterprise/factory/route.test.ts
CREATE src/app/api/enterprise/factory/route.ts
CREATE src/components/factory-unit-actions.tsx
CREATE src/platform/factory/contracts.ts
CREATE src/app/api/enterprise/franchise/whatsapp/replies/route.test.ts
CREATE src/app/api/enterprise/franchise/whatsapp/replies/route.ts
CREATE src/platform/notifications/reply-contract.test.ts
CREATE src/platform/notifications/reply-contract.ts
```

## 5. Materialization blocks

### FILE: `src/platform/backend/response-boundary.test.ts`

```yaml
block_id: "TS-GO-API-WEB-BRIDGE:response-boundary-test:v1"
operation: CREATE
provenance: AUTHORED
source: "Local regression tests for shared WHATWG/Node response reading; not copied from an upstream company"
license: "LicenseRef-Workspace-Owner"
sha256: "4ae89980e12a768e92503b6b72c58e54a45f5bdb07bc6df16ebaaca8d5071e03"
variables: []
secrets_allowed: false
```

````typescript
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { listModels } from "./public-client";
import { createServer, type RequestListener } from "node:http";
import { gzipSync } from "node:zlib";

const limit = 1_048_576;
const utf8 = new TextEncoder();
beforeEach(() => {
  vi.stubEnv("ENTERPRISE_API_BASE_URL", "https://api.example.test");
  vi.stubEnv("ENTERPRISE_TENANT_CODE", "tenant");
  vi.stubEnv("ENTERPRISE_ORGANIZATION_CODE", "store");
});

async function withHTTP(handler: RequestListener, run: () => Promise<void>) {
  const server = createServer(handler);
  await new Promise<void>((resolve, reject) => { server.once("error", reject); server.listen(0, "127.0.0.1", resolve); });
  try {
    const address = server.address();
    if (!address || typeof address === "string") throw new Error("loopback address missing");
    vi.stubEnv("ENTERPRISE_API_BASE_URL", `http://127.0.0.1:${address.port}`);
    await run();
  } finally {
    server.closeAllConnections();
    await new Promise<void>((resolve, reject) => server.close((error) => error ? reject(error) : resolve()));
  }
}

describe("native fetch against a real loopback HTTP server", () => {
  it("limits decompressed bytes rather than trusting compressed Content-Length", async () => {
    const compressed = gzipSync(JSON.stringify({ models: [], padding: "é".repeat(limit / 2) }));
    expect(compressed.byteLength).toBeLessThan(limit);
    let requests = 0;
    await withHTTP((_request, response) => {
      requests++;
      response.writeHead(200, { "content-type": "application/json", "content-encoding": "gzip", "content-length": compressed.byteLength });
      response.end(compressed);
    }, async () => {
      await expect(listModels()).rejects.toMatchObject({ code: "RESPONSE_TOO_LARGE" });
      expect(requests).toBe(1);
    });
  });
  it("cancels an oversized chunked response before the server finishes it", async () => {
    let closed!: () => void;
    const closure = new Promise<void>((resolve) => { closed = resolve; });
    await withHTTP((_request, response) => {
      response.on("close", closed);
      response.writeHead(200, { "content-type": "application/json" });
      response.write(" ".repeat(limit + 1));
      // No end: the client must reject and cancel rather than await EOF.
    }, async () => {
      await expect(listModels()).rejects.toMatchObject({ code: "RESPONSE_TOO_LARGE" });
      await closure;
    });
  }, 2000);
  it("retains the five-second deadline while waiting for body bytes", async () => {
    let requests = 0;
    await withHTTP((_request, response) => {
      requests++;
      response.writeHead(200, { "content-type": "application/json" });
      response.write('{"models":[');
    }, async () => {
      await expect(listModels()).rejects.toMatchObject({ code: "INVALID_RESPONSE" });
      expect(requests).toBe(1);
    });
  }, 7000);
});
afterEach(() => { vi.unstubAllGlobals(); vi.unstubAllEnvs(); });

function reply(chunks: Uint8Array[], contentType = "application/json", status = 200, cancel = vi.fn()) {
  let count = 0;
  const stream = new ReadableStream<Uint8Array>({
    pull(controller) { const chunk = chunks[count++]; if (chunk) controller.enqueue(chunk); else controller.close(); },
    cancel
  }, { highWaterMark: 0 });
  const fetch = vi.fn(async () => new Response(stream, { status, headers: { "content-type": contentType, "content-length": "1" } }));
  vi.stubGlobal("fetch", fetch);
  return { stream, cancel, fetch, pulls: () => count };
}

describe("bounded backend response transport", () => {
  it("rejects UTF-8 bytes above budget even when the string is shorter and Content-Length lies", async () => {
    reply([utf8.encode(JSON.stringify({ models: [], padding: "é".repeat(limit / 2) }))]);
    await expect(listModels()).rejects.toMatchObject({ status: 502, code: "RESPONSE_TOO_LARGE" });
  });
  it("cancels at the first excessive chunk without draining the rest or retrying", async () => {
    const fixture = reply([utf8.encode(" ".repeat(limit)), utf8.encode(" "), utf8.encode('{"models":[]}')]);
    await expect(listModels()).rejects.toMatchObject({ status: 502, code: "RESPONSE_TOO_LARGE" });
    expect(fixture.pulls()).toBe(2);
    expect(fixture.cancel).toHaveBeenCalledOnce();
    expect(fixture.fetch).toHaveBeenCalledOnce();
    expect(fixture.stream.locked).toBe(false);
  });
  it("accepts exactly the byte budget", async () => {
    const json = '{"models":[]}';
    reply([utf8.encode(json + " ".repeat(limit - utf8.encode(json).byteLength))]);
    await expect(listModels()).resolves.toEqual([]);
  });
  it("preserves Unicode split across byte chunks", async () => {
    const bytes = utf8.encode('{"models":[{"displayName":"Eléctrica 🚲"}]}');
    reply(Array.from(bytes, (byte) => new Uint8Array([byte])));
    await expect(listModels()).resolves.toEqual([{ displayName: "Eléctrica 🚲" }]);
  });
  it("cancels invalid media types without reading the body", async () => {
    const fixture = reply([utf8.encode("private body")], "text/html");
    await expect(listModels()).rejects.toMatchObject({ status: 502, code: "INVALID_CONTENT_TYPE" });
    expect(fixture.pulls()).toBe(0);
    expect(fixture.cancel).toHaveBeenCalledOnce();
  });
  it("does not wait for a stalled cancellation promise", async () => {
    const fixture = reply([utf8.encode(" ".repeat(limit + 1)), utf8.encode("{}")], "application/json", 200, vi.fn(() => new Promise<void>(() => {})));
    await expect(listModels()).rejects.toMatchObject({ code: "RESPONSE_TOO_LARGE" });
    expect(fixture.cancel).toHaveBeenCalledOnce();
    expect(fixture.stream.locked).toBe(false);
  }, 1000);
  for (const bytes of [utf8.encode('{"secret":"private-token",'), new Uint8Array([123, 34, 120, 34, 58, 34, 255, 34, 125])]) {
    it(`rejects malformed JSON/UTF-8 without propagating payload (${bytes.byteLength})`, async () => {
      reply([bytes]);
      await expect(listModels()).rejects.toMatchObject({ status: 502, code: "INVALID_RESPONSE", message: "backend request failed: INVALID_RESPONSE" });
    });
  }
  it("does not manufacture a successful response after a stream failure", async () => {
    vi.stubGlobal("fetch", vi.fn(async () => new Response(new ReadableStream({ start(controller) { controller.error(new Error("private-token")); } }), { headers: { "content-type": "application/json" } })));
    await expect(listModels()).rejects.toMatchObject({ status: 502, code: "INVALID_RESPONSE", message: "backend request failed: INVALID_RESPONSE" });
  });
  for (const body of ["null", '{"code":{"secret":"private-token"}}']) {
    it(`preserves error status but rejects unusable upstream codes (${body})`, async () => {
      reply([utf8.encode(body)], "application/problem+json", 403);
      await expect(listModels()).rejects.toMatchObject({ status: 403, code: "UPSTREAM_ERROR" });
    });
  }
});
````

### FILE: `src/platform/backend/public-client.ts`

```yaml
block_id: "TS-GO-BRIDGE:client:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "68f3b45038d22330b0d81aa7f3bf2795b7c847e0faa68d672302b5b009f08463"
variables: []
secrets_allowed: false
```

````typescript
const codePattern = /^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$/;
export type PublicModel = { id: string; code: string; displayName: string; vehicleClass: string; specification: Record<string, unknown> };
export type LeadInput = { modelId?: string; sourceCode: string; contact: Record<string, unknown>; consentGranted: true };
export type PublicLocation = { organization_id: string; code: string; name: string; city: string; region: string; country: string; contact_phone?: string; contact_email?: string };
export type AppointmentKind = "consultation" | "test-drive" | "delivery" | "service";
export type AppointmentSlot = { id: string; organization_id: string; kind: AppointmentKind; starts_at: string; ends_at: string; capacity: number; booked: number; state: "open"; version: number };
export type AppointmentInput = { leadId: string; modelId?: string; kind: AppointmentKind; startsAt: string };
export class BackendProblem extends Error { constructor(readonly status: number, readonly code: string) { super(`backend request failed: ${code}`); } }

// One transport boundary for public calls and server-only protected calls.
// Count decoded HTTP body bytes, not UTF-16 characters or Content-Length.
// This bounds our retained payload, not buffers owned by fetch/the network stack.
export async function readBackendResponse<T>(response: Response): Promise<T> {
  const contentType = (response.headers.get("content-type") ?? "").split(";", 1)[0]?.trim().toLowerCase();
  if (contentType !== "application/json" && !(!response.ok && contentType === "application/problem+json")) {
    if (response.body) void response.body.cancel().catch(() => {});
    throw new BackendProblem(502, "INVALID_CONTENT_TYPE");
  }
  if (!response.body) throw new BackendProblem(502, "INVALID_RESPONSE");
  const reader = response.body.getReader();
  const decoder = new TextDecoder("utf-8", { fatal: true });
  let received = 0;
  let body = "";
  let complete = false;
  try {
    while (true) {
      const chunk = await reader.read();
      if (chunk.done) { complete = true; break; }
      received += chunk.value.byteLength;
      if (received > 1_048_576) throw new BackendProblem(502, "RESPONSE_TOO_LARGE");
      body += decoder.decode(chunk.value, { stream: true });
    }
    body += decoder.decode();
    const value: unknown = JSON.parse(body);
    if (!response.ok) {
      const code = value !== null && typeof value === "object" && "code" in value ? value.code : undefined;
      throw new BackendProblem(response.status, typeof code === "string" && /^[A-Z][A-Z0-9_]{0,95}$/.test(code) ? code : "UPSTREAM_ERROR");
    }
    // Domain schemas remain the responsibility of each typed consumer.
    return value as T;
  } catch (error) {
    if (error instanceof BackendProblem) throw error;
    throw new BackendProblem(502, "INVALID_RESPONSE");
  } finally {
    // Cancellation may itself reject or never settle; do not delay the error
    // or retry a write whose durable outcome can be unknown.
    if (!complete) void reader.cancel().catch(() => {});
    reader.releaseLock();
  }
}
function settings() {
  const base = process.env.ENTERPRISE_API_BASE_URL;
  const tenant = process.env.ENTERPRISE_TENANT_CODE;
  const organization = process.env.ENTERPRISE_ORGANIZATION_CODE;
  if (!base || !URL.canParse(base) || !tenant || !organization || !codePattern.test(tenant) || !codePattern.test(organization)) throw new Error("enterprise backend configuration is invalid");
  const url = new URL(base); if (!['http:', 'https:'].includes(url.protocol) || url.username || url.password || url.search || url.hash) throw new Error("enterprise backend base URL is unsafe");
  return { base: url, tenant, organization };
}
async function call<T>(path: string, init?: RequestInit): Promise<{ value: T; status: number; replayed: boolean }> {
  const { base } = settings(); const target = new URL(path, base);
  const response = await fetch(target, { ...init, cache: "no-store", redirect: "error", signal: AbortSignal.timeout(5000), headers: { accept: "application/json", ...init?.headers } });
  const value = await readBackendResponse<T>(response);
  return { value, status: response.status, replayed: response.headers.get("idempotency-replayed") === "true" };
}
export async function listModels(): Promise<PublicModel[]> { const { tenant } = settings(); const result = await call<{ models: PublicModel[] }>(`/v1/public/${encodeURIComponent(tenant)}/models`); return result.value.models; }
export async function captureLead(input: LeadInput, idempotencyKey: string) { if (idempotencyKey.length < 16 || idempotencyKey.length > 128) throw new Error("invalid idempotency key"); const { tenant, organization } = settings(); return call<{ lead_id: string; status: string }>(`/v1/public/${encodeURIComponent(tenant)}/${encodeURIComponent(organization)}/leads`, { method: "POST", headers: { "content-type": "application/json", "idempotency-key": idempotencyKey }, body: JSON.stringify({ model_id: input.modelId ?? "", source_code: input.sourceCode, contact: input.contact, consent_granted: input.consentGranted }) }); }
export async function listLocations(): Promise<PublicLocation[]> { const { tenant } = settings(); const result = await call<{ items: PublicLocation[] }>(`/v1/public/${encodeURIComponent(tenant)}/locations`); return result.value.items; }
export async function listAppointmentSlots(kind: AppointmentKind, from: string, to: string): Promise<AppointmentSlot[]> { if (!Number.isFinite(Date.parse(from)) || !Number.isFinite(Date.parse(to)) || Date.parse(to) <= Date.parse(from)) throw new Error("invalid appointment slot range"); const { tenant, organization } = settings(); const query = new URLSearchParams({ kind, from, to }); const result = await call<{ items: AppointmentSlot[] }>(`/v1/public/${encodeURIComponent(tenant)}/${encodeURIComponent(organization)}/appointment-slots?${query}`); return result.value.items; }
export async function requestAppointment(input: AppointmentInput, idempotencyKey: string) { if (idempotencyKey.length < 16 || idempotencyKey.length > 128) throw new Error("invalid idempotency key"); const { tenant, organization } = settings(); return call<{ id: string; state: string; starts_at: string }>(`/v1/public/${encodeURIComponent(tenant)}/${encodeURIComponent(organization)}/appointments`, { method: "POST", headers: { "content-type": "application/json", "idempotency-key": idempotencyKey }, body: JSON.stringify({ lead_id: input.leadId, model_id: input.modelId ?? "", kind: input.kind, starts_at: input.startsAt }) }); }
````

### FILE: `src/platform/backend/public-client.test.ts`

```yaml
block_id: "TS-GO-BRIDGE:client-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "0ffeb7e4fc48101e2cb2710265ee8129dbd323ced09020aef2764d8aadd6c5c5"
variables: []
secrets_allowed: false
```

````typescript
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { BackendProblem, captureLead, listAppointmentSlots, listLocations, listModels, requestAppointment } from "./public-client";
beforeEach(() => { process.env.ENTERPRISE_API_BASE_URL = "https://api.example.test/"; process.env.ENTERPRISE_TENANT_CODE = "elite-mobility"; process.env.ENTERPRISE_ORGANIZATION_CODE = "central-store"; });
afterEach(() => { vi.unstubAllGlobals(); });
describe("public backend client", () => {
  it("preserves a Go problem response without accepting it as successful data", async () => {
    vi.stubGlobal("fetch", vi.fn().mockResolvedValue(new Response(JSON.stringify({ code: "IDEMPOTENCY_CONFLICT", status: 409 }), { status: 409, headers: { "content-type": "application/problem+json; charset=utf-8" } })));
    await expect(captureLead({ sourceCode: "public-web", contact: { email: "a@example.test" }, consentGranted: true }, "lead-request-problem-0001")).rejects.toEqual(new BackendProblem(409, "IDEMPOTENCY_CONFLICT"));
  });
  it("rejects a problem document on a successful HTTP response", async () => {
    vi.stubGlobal("fetch", vi.fn().mockResolvedValue(new Response(JSON.stringify({ code: "BROKEN" }), { status: 200, headers: { "content-type": "application/problem+json" } })));
    await expect(listModels()).rejects.toEqual(new BackendProblem(502, "INVALID_CONTENT_TYPE"));
  });
  it("uses the fixed tenant path and parses models", async () => { const fetch = vi.fn().mockResolvedValue(new Response(JSON.stringify({ models: [{ id: "m", code: "urban", displayName: "Urban", vehicleClass: "bicycle", specification: {} }] }), { status: 200, headers: { "content-type": "application/json" } })); vi.stubGlobal("fetch", fetch); expect((await listModels())[0]?.id).toBe("m"); expect(String(fetch.mock.calls[0]?.[0])).toBe("https://api.example.test/v1/public/elite-mobility/models"); });
  it("forwards the same idempotency key", async () => { const fetch = vi.fn().mockResolvedValue(new Response(JSON.stringify({ lead_id: "l", status: "accepted" }), { status: 200, headers: { "content-type": "application/json", "idempotency-replayed": "true" } })); vi.stubGlobal("fetch", fetch); const result = await captureLead({ modelId: "m", sourceCode: "public-web", contact: { email: "a@example.test" }, consentGranted: true }, "lead-request-00000001"); expect(result.replayed).toBe(true); expect((fetch.mock.calls[0]?.[1] as RequestInit).headers).toMatchObject({ "idempotency-key": "lead-request-00000001" }); });
  it("uses the same public scope for locations and appointment replay", async () => { const fetch = vi.fn().mockResolvedValueOnce(new Response(JSON.stringify({ items: [{ organization_id: "store", code: "central-store", name: "Central", city: "Cordoba", region: "Cordoba", country: "AR" }] }), { status: 200, headers: { "content-type": "application/json" } })).mockResolvedValueOnce(new Response(JSON.stringify({ id: "appointment", state: "requested", starts_at: "2026-09-01T15:00:00Z" }), { status: 200, headers: { "content-type": "application/json", "idempotency-replayed": "true" } })); vi.stubGlobal("fetch", fetch); expect((await listLocations())[0]?.code).toBe("central-store"); const result = await requestAppointment({ leadId: "lead", modelId: "model", kind: "test-drive", startsAt: "2026-09-01T15:00:00Z" }, "appointment-key-0001"); expect(result.replayed).toBe(true); expect(String(fetch.mock.calls[1]?.[0])).toContain("/v1/public/elite-mobility/central-store/appointments"); });
  it("queries bounded server-owned appointment capacity", async () => { const fetch = vi.fn().mockResolvedValue(new Response(JSON.stringify({ items: [{ id: "slot", organization_id: "store", kind: "service", starts_at: "2026-09-01T15:00:00Z", ends_at: "2026-09-01T16:00:00Z", capacity: 2, booked: 1, state: "open", version: 1 }] }), { status: 200, headers: { "content-type": "application/json" } })); vi.stubGlobal("fetch", fetch); const slots = await listAppointmentSlots("service", "2026-09-01T00:00:00Z", "2026-09-02T00:00:00Z"); expect(slots[0]?.booked).toBe(1); const target = String(fetch.mock.calls[0]?.[0]); expect(target).toContain("/v1/public/elite-mobility/central-store/appointment-slots?"); expect(target).toContain("kind=service"); });
  it("rejects non-json upstream responses", async () => { vi.stubGlobal("fetch", vi.fn().mockResolvedValue(new Response("gateway", { status: 502, headers: { "content-type": "text/plain" } }))); await expect(listModels()).rejects.toEqual(new BackendProblem(502, "INVALID_CONTENT_TYPE")); });
  it("rejects credentials in the base URL", async () => { process.env.ENTERPRISE_API_BASE_URL = "https://user:pass@example.test"; await expect(listModels()).rejects.toThrow("unsafe"); });
});
````

### FILE: `src/app/api/enterprise/models/route.ts`

```yaml
block_id: "TS-GO-BRIDGE:models-route:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "543d404d4553c0120bbe70308437c5701d1fb018538c7b91b0f3752033b28b80"
variables: []
secrets_allowed: false
```

````typescript
import { listModels } from "@/platform/backend/public-client";
import { errorResponse } from "@/platform/http/problem";
export const dynamic = "force-dynamic";
export async function GET() { try { return Response.json({ models: await listModels() }, { headers: { "cache-control": "public, max-age=30, stale-while-revalidate=120" } }); } catch (error) { return errorResponse(error); } }
````

### FILE: `src/app/api/enterprise/leads/route.ts`

```yaml
block_id: "TS-GO-BRIDGE:leads-route:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "32427962f65b957ef86a09d774cba9d715b1eff99ec478dae5287a158e1fa20d"
variables: []
secrets_allowed: false
```

````typescript
import { z } from "zod";
import { captureLead } from "@/platform/backend/public-client";
import { errorResponse, problem } from "@/platform/http/problem";
const schema = z.object({ modelId: z.string().min(1).max(128).optional(), name: z.string().trim().min(2).max(120), email: z.email().max(254), consentGranted: z.literal(true) }).strict();
export async function POST(request: Request) { const key = request.headers.get("idempotency-key"); if (!key || key.length < 16 || key.length > 128) return problem(400, "Falta idempotencia", "Envíe una clave estable para este intento.", "IDEMPOTENCY_KEY_REQUIRED"); try { const input = schema.parse(await request.json()); const result = await captureLead({ ...(input.modelId ? { modelId: input.modelId } : {}), sourceCode: "public-web", contact: { name: input.name, email: input.email.toLowerCase() }, consentGranted: true }, key); return Response.json(result.value, { status: result.status, headers: { "cache-control": "no-store", ...(result.replayed ? { "idempotency-replayed": "true" } : {}) } }); } catch (error) { return errorResponse(error); } }
````

### FILE: `src/app/connected/lead-form.tsx`

```yaml
block_id: "TS-GO-BRIDGE:lead-form:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "19e31840eadf7bb275956580365d6faf6913d70cbbcd1fabf36e19844d2c48b7"
variables: []
secrets_allowed: false
```

````tsx
"use client";
import { publicMessage as message } from "@/platform/i18n/public-catalog";
import { FormEvent, useRef, useState } from "react";
export function LeadForm({ modelId, locale }: { modelId: string; locale: string }) { const key = useRef(crypto.randomUUID()); const [statusMessage, setMessage] = useState(""); const [leadId, setLeadId] = useState(""); const [pending, setPending] = useState(false); async function submit(event: FormEvent<HTMLFormElement>) { event.preventDefault(); const formElement = event.currentTarget; setPending(true); setMessage(""); setLeadId(""); const form = new FormData(formElement); try { const response = await fetch("/api/enterprise/leads", { method: "POST", headers: { "content-type": "application/json", "idempotency-key": key.current }, body: JSON.stringify({ modelId, name: form.get("name"), email: form.get("email"), consentGranted: form.get("consent") === "on" }) }); if (!response.ok) { setMessage(message(locale, "lead.failed")); return; } const body = await response.json() as { lead_id?: string }; if (!body.lead_id) { setMessage(message(locale, "lead.invalid")); return; } setLeadId(body.lead_id); setMessage(message(locale, "lead.received")); key.current = crypto.randomUUID(); formElement.reset(); } catch { setMessage(message(locale, "lead.failed")); } finally { setPending(false); } } return <form lang={locale} onSubmit={submit} aria-describedby={`result-${modelId}`}><label>{message(locale, "lead.name")}<input name="name" required minLength={2} maxLength={120} autoComplete="name" /></label><label>{message(locale, "lead.email")}<input name="email" type="email" required maxLength={254} autoComplete="email" /></label><label><input name="consent" type="checkbox" required /> {message(locale, "lead.consent")}</label><button disabled={pending} type="submit">{pending ? message(locale, "lead.sending") : message(locale, "lead.submit")}</button><p id={`result-${modelId}`} role="status" aria-live="polite">{statusMessage}</p>{leadId ? <a className="primaryAction" href={`/locations?lead_id=${encodeURIComponent(leadId)}&model_id=${encodeURIComponent(modelId)}`}>{message(locale, "lead.next")}</a> : null}</form>; }
````

### FILE: `src/app/connected/page.tsx`

```yaml
block_id: "TS-GO-BRIDGE:page:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "e68a3a1fdb1d7483dabc328723620cda9cd87b9704b59bb1ffa05d94252402ab"
variables: []
secrets_allowed: false
```

````tsx
import { publicMessage as message } from "@/platform/i18n/public-catalog";
import { loadPublicLocale } from "@/platform/i18n/load-public-locale";
import { publicCount } from "@/platform/i18n/public-catalog";
import { listModels } from "@/platform/backend/public-client";
import { LeadForm } from "./lead-form";
export const dynamic = "force-dynamic";
export default async function ConnectedCatalogPage() { const [models, locale] = await Promise.all([listModels(), loadPublicLocale()]); return <section lang={locale.locale}><div className="eyebrow">{message(locale.locale, "models.eyebrow")}</div><h1 className="pageTitle">{message(locale.locale, "models.title")}</h1><p className="lede">{message(locale.locale, "models.description")}</p><p>{publicCount(locale.locale, "models", models.length)}</p><div className="grid">{models.map(model => <article className="card" key={model.id}><h2>{model.displayName}</h2><p>{model.vehicleClass}</p><LeadForm modelId={model.id} locale={locale.locale} /></article>)}</div></section>; }
````

### FILE: `.env.example`

```yaml
block_id: "TS-GO-BRIDGE:env-example:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified standalone BFF"
license: "LicenseRef-Workspace-Owner"
sha256: "43fae1630727b68baf1c7a7e2e25a1b889e07028ab3adf5c00e81fd1c52cdf7f"
variables: []
secrets_allowed: false
```

````text
# Non-secret business presentation configuration. The file name is resolved
# inside ./config; arbitrary paths are rejected.
BUSINESS_CONFIG_FILE=business.example.json

# API business policy: both absent selects the embedded hash-locked reference.
# A custom file requires its exact SHA-256 from the reviewed deployment lock.
# These are non-secret fields; incomplete or invalid configuration stops startup.
BUSINESS_POLICY_PROFILE_FILE=
BUSINESS_POLICY_PROFILE_SHA256=

# Public origin registered at the identity provider. Production requires HTTPS.
APP_BASE_URL=http://localhost:3000
AUTH_SESSION_SECRET=replace-with-at-least-32-random-characters
OIDC_ISSUER=https://identity.example.com
OIDC_CLIENT_ID=enterprise-web
OIDC_CLIENT_SECRET=inject-from-secret-manager

# The BFF is the only browser-facing caller. The Go service remains the
# transactional/domain backend and validates every access token again.
ENTERPRISE_API_BASE_URL=https://api.example.com
ENTERPRISE_TENANT_CODE=example-tenant
ENTERPRISE_ORGANIZATION_CODE=example-store

# Hosted payment infrastructure: disabled until its complete account scope is
# configured. No value below is a credential or authorization for a live charge.
PAYMENT_CHECKOUT_ENABLED=false
PAYMENT_TENANT_ID=
PAYMENT_ORGANIZATION_ID=
PAYMENT_CONNECTION_ID=
# stripe | mercadopago (Argentina Checkout Pro lane)
PAYMENT_REQUEST_PROVIDER=
PAYMENT_ACCOUNT_REF=
PAYMENT_CURRENCY=
PAYMENT_MINOR_UNIT_EXPONENT=
# Explicit true | false when enabled; false selects provider sandbox/test mode.
PAYMENT_LIVE_MODE=
PAYMENT_SUCCESS_URL=
PAYMENT_CANCEL_URL=
# MercadoPago only: https://<public-api-origin>/v1/payment-provider/webhook
PAYMENT_NOTIFICATION_URL=
PAYMENT_DISPLAY_NAME=
PAYMENT_WORKER_ID=
# Inject through the environment/secret manager when the owner chooses to start.
PAYMENT_PROVIDER_SECRET=
PAYMENT_WEBHOOK_SECRET=
PAYMENT_OUTBOUND_HMAC_KEY_BASE64=

# Initial handover: explicit, non-secret materialized profile activation.
# These values must match the profile receipt and the active payment scope.
HANDOVER_ENABLED=false
HANDOVER_PROFILE_FILE=
HANDOVER_PROFILE_ID=
HANDOVER_PROFILE_REVISION=
HANDOVER_PROFILE_SHA256=
# Tenant, organization, provider, connection, account and mode are taken from
# the prepared payment runtime and must match the hash-locked profile document.

# Optional exact supplied-snapshot FX conversion receipts; no journal posting.
FX_ENABLED=false
FX_PROFILE_FILE=
FX_RATES_FILE=
FX_PROFILE_ID=
FX_PROFILE_REVISION=
FX_PROFILE_SHA256=
FX_TENANT_ID=
FX_ORGANIZATION_ID=
FX_LOCAL_CURRENCY=

# Connected WhatsApp host: configure the exact private profile and external secret files.
WHATSAPP_ENABLED=false
WHATSAPP_HOST_PROFILE_FILE=
WHATSAPP_HOST_PROFILE_SHA256=
WHATSAPP_ACCESS_TOKEN_FILE=
WHATSAPP_APP_SECRET_FILE=
WHATSAPP_VERIFY_TOKEN_FILE=
WHATSAPP_SERVICE_TOKEN_FILE=
# Choose the token file above OR an OIDC client-credentials broker below.
WHATSAPP_SERVICE_IDENTITY_PROFILE_FILE=
WHATSAPP_SERVICE_IDENTITY_PROFILE_SHA256=
WHATSAPP_SERVICE_CLIENT_SECRET_FILE=
WHATSAPP_CONTACT_HMAC_KEY_HEX_FILE=
LLM_API_KEY_FILE=
WHATSAPP_COLLECTOR_CLIENT_CERT_FILE=
WHATSAPP_COLLECTOR_CLIENT_KEY_FILE=
````

### FILE: `.gitignore`

```yaml
block_id: "TS-GO-BRIDGE:gitignore:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified standalone BFF"
license: "LicenseRef-Workspace-Owner"
sha256: "1cf6cf78f1a64047af43a0372ccae8ebb70abf5c99d80b37925c8f0aadad6286"
variables: []
secrets_allowed: false
```

````text
node_modules/
.next/
.data/
coverage/
.env
.env.local
*.log
*.tsbuildinfo
````

### FILE: `config/business.example.json`

```yaml
block_id: "TS-GO-BRIDGE:business-config:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified standalone BFF"
license: "LicenseRef-Workspace-Owner"
sha256: "d25541768786b8b20e534b013e523ef464170ea045c39d51a5223ed1f80f6941"
variables: []
secrets_allowed: false
```

````text
{
  "schemaVersion": "1.0.0",
  "business": {
    "id": "electric-mobility-network",
    "name": "Electric Mobility Network",
    "defaultLocale": "es-AR",
    "defaultMarket": "AR",
    "supportEmail": "soporte@example.invalid"
  },
  "markets": [
    {
      "code": "AR",
      "name": "Argentina",
      "currency": "ARS",
      "locales": ["es-AR"],
      "timeZone": "America/Argentina/Buenos_Aires",
      "taxMode": "external"
    }
  ],
  "organizationTypes": [
    {"id": "hq", "label": "Casa central", "allowedParents": []},
    {"id": "franchise", "label": "Franquicia", "allowedParents": ["hq"]},
    {"id": "branch", "label": "Sucursal", "allowedParents": ["franchise", "hq"]},
    {"id": "factory", "label": "Fábrica", "allowedParents": ["hq"]},
    {"id": "supplier", "label": "Proveedor", "allowedParents": ["hq"]}
  ],
  "roles": [
    {"id": "hq_admin", "label": "Administrador central", "permissions": ["*"]},
    {"id": "branch_manager", "label": "Responsable de sucursal", "permissions": ["admin:read", "lead:read", "lead:assign", "lead:update", "quote:write", "catalog:write", "pricing:write", "order:create", "order:write", "inventory:allocate", "payment:create", "payment:write", "service:write", "communication:write"]},
    {"id": "sales", "label": "Ventas", "permissions": ["admin:read", "lead:read", "lead:assign", "lead:update", "quote:write", "catalog:write", "order:create", "order:write", "payment:create"]},
    {"id": "factory_operator", "label": "Operador de fábrica", "permissions": ["factory:read", "procurement:write", "factory:write", "inventory:write", "logistics:write"]},
    {"id": "customer", "label": "Cliente", "permissions": ["customer:self"]}
  ],
  "modules": {
    "catalog": {"enabled": true},
    "crm": {"enabled": true},
    "procurement": {"enabled": true},
    "inventory": {"enabled": true},
    "orders": {"enabled": true},
    "payments": {"enabled": true},
    "fulfillment": {"enabled": true},
    "service": {"enabled": true},
    "documents": {"enabled": true},
    "integrations": {"enabled": true}
  },
  "workflows": {
    "lead": {
      "initial": "new",
      "states": ["new", "contacted", "qualified", "converted", "lost"],
      "transitions": [
        {"from": "new", "to": "contacted", "permission": "lead:update"},
        {"from": "contacted", "to": "qualified", "permission": "lead:update"},
        {"from": "qualified", "to": "converted", "permission": "lead:update"},
        {"from": "contacted", "to": "lost", "permission": "lead:update"},
        {"from": "qualified", "to": "lost", "permission": "lead:update"}
      ]
    },
    "order": {
      "initial": "draft",
      "states": ["draft", "placed", "confirmed", "paid", "allocated", "delivered", "cancelled"],
      "transitions": [
        {"from": "draft", "to": "placed", "permission": "order:create"},
        {"from": "placed", "to": "confirmed", "permission": "order:transition"},
        {"from": "confirmed", "to": "paid", "permission": "payment:reconcile"},
        {"from": "paid", "to": "allocated", "permission": "inventory:reserve"},
        {"from": "allocated", "to": "delivered", "permission": "order:transition"},
        {"from": "draft", "to": "cancelled", "permission": "order:transition"},
        {"from": "placed", "to": "cancelled", "permission": "order:transition"}
      ]
    }
  },
  "customFields": {
    "lead": [
      {"id": "preferred_vehicle_use", "label": "Uso principal", "type": "select", "required": false, "options": ["urban", "delivery", "recreation", "fleet"]}
    ],
    "catalog_model": [
      {"id": "estimated_range_km", "label": "Autonomía estimada (km)", "type": "number", "required": true}
    ]
  },
  "integrations": [
    {"id": "mercado_pago", "provider": "mercado_pago", "enabled": false, "mode": "sandbox", "capabilities": ["payments"], "credentialRefEnv": "MERCADO_PAGO_CREDENTIAL_REF"},
    {"id": "amazon_sp_api", "provider": "amazon_sp_api", "enabled": false, "mode": "sandbox", "capabilities": ["catalog", "orders", "fulfillment"], "credentialRefEnv": "AMAZON_SP_API_CREDENTIAL_REF"},
    {"id": "mercado_libre", "provider": "mercado_libre", "enabled": false, "mode": "sandbox", "capabilities": ["catalog", "orders", "fulfillment"], "credentialRefEnv": "MERCADO_LIBRE_CREDENTIAL_REF"},
    {"id": "google_ads", "provider": "google_ads", "enabled": false, "mode": "sandbox", "capabilities": ["ads", "conversions"], "credentialRefEnv": "GOOGLE_ADS_CREDENTIAL_REF"},
    {"id": "meta_ads", "provider": "meta_ads", "enabled": false, "mode": "sandbox", "capabilities": ["ads", "conversions"], "credentialRefEnv": "META_ADS_CREDENTIAL_REF"}
  ],
  "features": {
    "public_catalog": true,
    "lead_capture": true,
    "customer_portal": true,
    "factory_portal": true,
    "vehicle_telemetry": false
  }
}
````

### FILE: `next-env.d.ts`

```yaml
block_id: "TS-GO-BRIDGE:next-env:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified standalone BFF"
license: "LicenseRef-Workspace-Owner"
sha256: "f2b3bca04d1bfe583daae1e1f798c92ec24bb6693bd88d0a09ba6802dee362a8"
variables: []
secrets_allowed: false
```

````text
/// <reference types="next" />
/// <reference types="next/image-types/global" />

// NOTE: This file should not be edited
// see https://nextjs.org/docs/app/api-reference/config/typescript for more information.
````

### FILE: `next.config.ts`

```yaml
block_id: "TS-GO-BRIDGE:next-config:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified standalone BFF"
license: "LicenseRef-Workspace-Owner"
sha256: "80b28ef70097951bd8106ec3f6a9cf03373f0f22624687a1916ab0fa2c11c657"
variables: []
secrets_allowed: false
```

````text
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  poweredByHeader: false,
  typedRoutes: true,
  // This reference does not transform images at runtime. Re-admit an image
  // pipeline and its dependencies before enabling the built-in optimizer.
  images: { unoptimized: true },
  async headers() {
    return [
      {
        source: "/(.*)",
        headers: [
          { key: "X-Content-Type-Options", value: "nosniff" },
          { key: "Referrer-Policy", value: "strict-origin-when-cross-origin" },
          { key: "Permissions-Policy", value: "camera=(), microphone=(), geolocation=()" },
          { key: "X-Frame-Options", value: "DENY" },
          { key: "Cross-Origin-Opener-Policy", value: "same-origin" }
        ]
      }
    ];
  }
};

export default nextConfig;
````

### FILE: `package.json`

```yaml
block_id: "TS-GO-BRIDGE:package:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified standalone BFF"
license: "LicenseRef-Workspace-Owner"
sha256: "acc4e4bd62c13cf48e071e933266d82d067cab849f92474d71a533ebf6508951"
variables: []
secrets_allowed: false
```

````text
{
  "name": "elite-enterprise-web-bff",
  "version": "0.4.0",
  "private": true,
  "type": "module",
  "packageManager": "pnpm@11.25.0",
  "engines": {
    "node": ">=24.0.0"
  },
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "typecheck": "next typegen && tsc --noEmit",
    "test": "vitest run",
    "test:watch": "vitest",
    "licenses:report": "pnpm licenses list --prod --json",
    "verify": "pnpm typecheck && pnpm test && pnpm build"
  },
  "dependencies": {
    "jose": "6.2.10",
    "next": "16.3.4",
    "openid-client": "6.8.5",
    "react": "19.2.8",
    "react-dom": "19.2.8",
    "safevalues": "1.2.0",
    "server-only": "0.0.1",
    "zod": "4.4.3"
  },
  "devDependencies": {
    "@types/node": "26.2.0",
    "@types/react": "19.2.18",
    "@types/react-dom": "19.2.5",
    "csp_evaluator": "1.1.8",
    "typescript": "7.0.2",
    "vitest": "4.1.11"
  }
}
````

### FILE: `pnpm-lock.yaml`

```yaml
block_id: "TS-GO-BRIDGE:pnpm-lock:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified standalone BFF"
license: "LicenseRef-Workspace-Owner"
sha256: "d6c73eb82a56c7e8d21ee89ffb402d0adb2f0c15db9c33bb8ec7035775119928"
variables: []
secrets_allowed: false
```

````text
lockfileVersion: '9.0'

settings:
  autoInstallPeers: true
  excludeLinksFromLockfile: false

importers:

  .:
    dependencies:
      jose:
        specifier: 6.2.10
        version: 6.2.10
      next:
        specifier: 16.3.4
        version: 16.3.4(@types/node@26.2.0)(react-dom@19.2.8(react@19.2.8))(react@19.2.8)
      openid-client:
        specifier: 6.8.5
        version: 6.8.5
      react:
        specifier: 19.2.8
        version: 19.2.8
      react-dom:
        specifier: 19.2.8
        version: 19.2.8(react@19.2.8)
      safevalues:
        specifier: 1.2.0
        version: 1.2.0
      server-only:
        specifier: 0.0.1
        version: 0.0.1
      zod:
        specifier: 4.4.3
        version: 4.4.3
    devDependencies:
      '@types/node':
        specifier: 26.2.0
        version: 26.2.0
      '@types/react':
        specifier: 19.2.18
        version: 19.2.18
      '@types/react-dom':
        specifier: 19.2.5
        version: 19.2.5(@types/react@19.2.18)
      csp_evaluator:
        specifier: 1.1.8
        version: 1.1.8
      typescript:
        specifier: 7.0.2
        version: 7.0.2
      vitest:
        specifier: 4.1.11
        version: 4.1.11(@types/node@26.2.0)(vite@8.2.2(@types/node@26.2.0))

packages:

  '@jridgewell/sourcemap-codec@1.5.5':
    resolution: {integrity: sha512-cYQ9310grqxueWbl+WuIUIaiUaDcj7WOq5fVhEljNVgRfOUhY9fy2zTvfoqWsnebh8Sl70VScFbICvJnLKB0Og==}

  '@next/env@16.3.4':
    resolution: {integrity: sha512-cjWZnUUa6jZq2kFaNe/ZyJdZonOZ/QoN0Zka2nz/FLOrfx14pQuM9c5RaSVkWMqgdt4ksgPAMWPyHSs/CyV48Q==}

  '@next/swc-darwin-arm64@16.3.4':
    resolution: {integrity: sha512-iBr3I5LZNk5/bgl5//iTgD2tcym14MX0Xo7fD//u9dYAEgGzza1y9oywluPtf74YnOswVdH1908aK9xVz7zQTw==}
    engines: {node: '>= 10'}
    cpu: [arm64]
    os: [darwin]

  '@next/swc-darwin-x64@16.3.4':
    resolution: {integrity: sha512-2dpiSyl2Jw/NrBPaU2MAKGSa+2MR82pJIn4Sm5Rjr+gxAeuh0z158Su3Z2O8zn7UNNq+ej4bToed6RcRN/Lydg==}
    engines: {node: '>= 10'}
    cpu: [x64]
    os: [darwin]

  '@next/swc-linux-arm64-gnu@16.3.4':
    resolution: {integrity: sha512-+t+U8HZT+fApePCS5h89CSH3datz29MkzyfCn+6fpsZBG/oiEOhINcb9rtkv6sdpToLGFn2e6146NzaKCXkqrA==}
    engines: {node: '>= 10'}
    cpu: [arm64]
    os: [linux]
    libc: [glibc]

  '@next/swc-linux-arm64-musl@16.3.4':
    resolution: {integrity: sha512-mx03GNs1ocQA5JQ4FxDMmIsNkdrZh8cuezKCrId28e5/gIPU/l7Kcy2+vmCCzdjnnmXJy+iOAu+7K0QppO6Urg==}
    engines: {node: '>= 10'}
    cpu: [arm64]
    os: [linux]
    libc: [musl]

  '@next/swc-linux-x64-gnu@16.3.4':
    resolution: {integrity: sha512-YIhGY6fSMfha52bnVxnzc9zaVBzJg+cqQTOD8tXIBSx4fuv0pVMxQTE0PaS59YhnMOiYiG09IMwxJAf/CFm/Dw==}
    engines: {node: '>= 10'}
    cpu: [x64]
    os: [linux]
    libc: [glibc]

  '@next/swc-linux-x64-musl@16.3.4':
    resolution: {integrity: sha512-+eaaX6axpDb0yF1GCpiERe6njplvdC+nks/fKfcHu3XPGRrald8P3/X7yv7QLdjA51knnxwl9pxdIJsg+w1L+Q==}
    engines: {node: '>= 10'}
    cpu: [x64]
    os: [linux]
    libc: [musl]

  '@next/swc-win32-arm64-msvc@16.3.4':
    resolution: {integrity: sha512-0jcXW7Xs/uzICrmgV3MhDYDeRy++1CqnpDIerlPIqYO4bhzB4WNbX/aRnQclustsAyTkFKB0z6rbcjmNg5tR8A==}
    engines: {node: '>= 10'}
    cpu: [arm64]
    os: [win32]

  '@next/swc-win32-x64-msvc@16.3.4':
    resolution: {integrity: sha512-vvBzwu1pYQCp92maZCFCIw/XgOTMR5tur9GjakwIo2cmwRTMKajRZZDS9+e4KsUZWKu1E007WUeAFXRRjZeuzw==}
    engines: {node: '>= 10'}
    cpu: [x64]
    os: [win32]

  '@oxc-project/types@0.146.0':
    resolution: {integrity: sha512-XC0QsnnhVe7sLIWmYmdPw7x5P0h4W8vUU3Nv1ySgWXtvCz8NizoAEpGXA0sOYoJQV2Rl13LgURAHQ5cI5ILCSA==}

  '@rolldown/binding-android-arm-eabi@1.2.5':
    resolution: {integrity: sha512-DLe/i+l8ynIBY7XEQ191TeZvCoowIGa18R+dIV30GW7DiOtp74i/xX8hs8GUjW5ARV7VZuie3d6AumSmCwbeRA==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [arm]
    os: [android]

  '@rolldown/binding-android-arm64@1.2.5':
    resolution: {integrity: sha512-zXcwKlQApYAOELHd8PwKDFkagYF9Wy4e0RJ+0qnzl9Pjnpj75TEG8ufv40p2J7kCEfwZAsNiuzRIyNNMWT38ig==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [arm64]
    os: [android]

  '@rolldown/binding-darwin-arm64@1.2.5':
    resolution: {integrity: sha512-dK4QakI42nzWgJT5sm4y4y/O//D4OxM75/cH28RLV+nzIN9AY+YsbuUVrUTjlLjXR6vpyxFbSsbmNuJ6BP9sww==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [arm64]
    os: [darwin]

  '@rolldown/binding-darwin-x64@1.2.5':
    resolution: {integrity: sha512-fqSALaUu1Wjd1nK2uW2kJDWdLCc8lx1IcY+MTY26Aurfdx19anlzhqXOgCFbBFQnlFDTn4TC1/7Nz4Bl2mLP3A==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [x64]
    os: [darwin]

  '@rolldown/binding-freebsd-x64@1.2.5':
    resolution: {integrity: sha512-/vCnNxlkxs9tKxNDcyWUePpJ/PgTzxIaVhoM5SmG8UV+GR/IcPam4VYxi7GIMo7PSDuNqlJqvprqii9NqqVCMw==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [x64]
    os: [freebsd]

  '@rolldown/binding-linux-arm-gnueabihf@1.2.5':
    resolution: {integrity: sha512-abk0NLA519LxRCszmbE0jYKuQ9YPocOXTiOXOo6Yr+YAT95VH+PtqYAjOJvGKt3viEd/x4qzabAlwd5bHOOARg==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [arm]
    os: [linux]

  '@rolldown/binding-linux-arm64-gnu@1.2.5':
    resolution: {integrity: sha512-Y7eALiJ8lr0M2HH103Js+g7V34wf6snlpZLAsHI90uLhr3PVlNsbFVAXJC9d/V6BnPyKtpSwI+NcB/RLxsQxuA==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [arm64]
    os: [linux]
    libc: [glibc]

  '@rolldown/binding-linux-arm64-musl@1.2.5':
    resolution: {integrity: sha512-xMvZgnbZg4YVnR/AX2b3oOPDTFYJvUVaJg5FedA/LuvexAtXibZQej4cnTkw3rjsJ/ggUROB64TdtETiim+FYA==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [arm64]
    os: [linux]
    libc: [musl]

  '@rolldown/binding-linux-ppc64-gnu@1.2.5':
    resolution: {integrity: sha512-GRjeqTUDHTo5GwntsLaAMcBahG3nlpjftXWZLN73HiYQlhwEowvarFgQnRnQZtIp4keXX7quXFbG38uPZBa2EA==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [ppc64]
    os: [linux]
    libc: [glibc]

  '@rolldown/binding-linux-s390x-gnu@1.2.5':
    resolution: {integrity: sha512-vLNTR45F2Uwc8AufkNXPmB4VliaXs+FvcheEogIzOXzO4l+LzieXF5A/TWxLy5HtqpsRCHUfd0lPVrrdgXdLHQ==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [s390x]
    os: [linux]
    libc: [glibc]

  '@rolldown/binding-linux-x64-gnu@1.2.5':
    resolution: {integrity: sha512-Mgj59/HTuYeK9Gz2MA+mBWKnHsAgkBSec15ZMb1st3oIfFbX7gCjOae7GydHhzcyQi9Z/7M1QuN9bR3oFqF0jQ==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [x64]
    os: [linux]
    libc: [glibc]

  '@rolldown/binding-linux-x64-musl@1.2.5':
    resolution: {integrity: sha512-mY8AP0/ichsbhAxGnLa3d3+MwV0EfgrPND2bplI3Ym8T6R2pJ0N87bvrKVwNXmdy3jnr6eQBecdqx/HMknBmpA==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [x64]
    os: [linux]
    libc: [musl]

  '@rolldown/binding-openharmony-arm64@1.2.5':
    resolution: {integrity: sha512-8SLssA2oweAxyRgDp789ACfRb/3P+zNRJpzZxSizxF9m8NUDQ4+3xjo8ttjhVGGw6Qxb70oZiEtIjaKikCO7Yw==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [arm64]
    os: [openharmony]

  '@rolldown/binding-win32-arm64-msvc@1.2.5':
    resolution: {integrity: sha512-vGbruD5zquhoc8D9SViXgN2FBJtNdTyQ4DtG+SWiEGlJiAzoKcZ2xp+xuXCffhubVdt0NJlTZqkeRuERy7g8Cw==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [arm64]
    os: [win32]

  '@rolldown/binding-win32-x64-msvc@1.2.5':
    resolution: {integrity: sha512-e/SXpgISz+IoqVcSSI0rx/d/he8zqLex+/rCWpnHpmVfmPIUjag9H6P7zotf0gJHwPUhQxZ/mF8tr6acebT9yw==}
    engines: {node: ^20.19.0 || >=22.12.0}
    cpu: [x64]
    os: [win32]

  '@rolldown/pluginutils@1.0.1':
    resolution: {integrity: sha512-2j9bGt5Jh8hj+vPtgzPtl72j0yRxHAyumoo6TNfAjsLB04UtpSvPbPcDcBMxz7n+9CYB0c1GxQFxYRg2jimqGw==}

  '@standard-schema/spec@1.1.0':
    resolution: {integrity: sha512-l2aFy5jALhniG5HgqrD6jXLi/rUWrKvqN/qJx6yoJsgKhblVd+iqqU4RCXavm/jPityDo5TCvKMnpjKnOriy0w==}

  '@swc/helpers@0.5.23':
    resolution: {integrity: sha512-5lSsMOTXURePglDfvuAQUqkGek9Hg2kksOYay2m0+XR++b2NWYL/4sWyuvVBIs8oKnJaxkdi9whaL/sqN13afw==}

  '@types/chai@5.2.3':
    resolution: {integrity: sha512-Mw558oeA9fFbv65/y4mHtXDs9bPnFMZAL/jxdPFUpOHHIXX91mcgEHbS5Lahr+pwZFR8A7GQleRWeI6cGFC2UA==}

  '@types/deep-eql@4.0.2':
    resolution: {integrity: sha512-c9h9dVVMigMPc4bwTvC5dxqtqJZwQPePsWjPlpSOnojbor6pGqdk541lfA7AqFQr5pB1BRdq0juY9db81BwyFw==}

  '@types/estree@1.0.9':
    resolution: {integrity: sha512-GhdPgy1el4/ImP05X05Uw4cw2/M93BCUmnEvWZNStlCzEKME4Fkk+YpoA5OiHNQmoS7Cafb8Xa3Pya8m1Qrzeg==}

  '@types/node@26.2.0':
    resolution: {integrity: sha512-5IviulTZeRNp2vAJ514cc/HUlY5nZ9fCbq9DMyC52BrhFZACo3nI0R7qBxhQmo/d27NFe96ur/b7Wwxklda+kg==}

  '@types/react-dom@19.2.5':
    resolution: {integrity: sha512-fMPwH9v7r/pp43yUd2/Mbiex5KouJwwR3dzHkhLREUC6764VyDsqxhAxv6OFEYR1RhjOyD1naqba8ECDBe7ZQg==}
    peerDependencies:
      '@types/react': ^19.2.0

  '@types/react@19.2.18':
    resolution: {integrity: sha512-AnzbBERsrLKtk2XSfTbYRLjQPdy116Sty4q+T+Bp3IC4l6jNBvreVPAHmpq9qhXQM7CXZPjLVmGMw9sy+hxQ3w==}

  '@typescript/typescript-aix-ppc64@7.0.2':
    resolution: {integrity: sha512-MTKKkWB7p/0E9xi1d1tHtZ5PiLkGEMIq88pK2CubZjOsLtYTLqhgIgi6zepFa+9GHZ6h05NMCkQxGKiPXMxXtQ==}
    engines: {node: '>=16.20.0'}
    cpu: [ppc64]
    os: [aix]

  '@typescript/typescript-darwin-arm64@7.0.2':
    resolution: {integrity: sha512-gowzar9MwS/aRWp6f3a4KUqzRjAZjOsmGNCM6LcTgXum+dBfgsBVMN+AgvOCCbguXyick6LJhpBszxMebJ8syA==}
    engines: {node: '>=16.20.0'}
    cpu: [arm64]
    os: [darwin]

  '@typescript/typescript-darwin-x64@7.0.2':
    resolution: {integrity: sha512-SZ9xZInqApNlNGc9s0W1VSsktYSOe9cFqNOIqmN1Gs8SmkjKZYFt017G4VwPxASInODuAdbTW7sXiFUf893RgA==}
    engines: {node: '>=16.20.0'}
    cpu: [x64]
    os: [darwin]

  '@typescript/typescript-freebsd-arm64@7.0.2':
    resolution: {integrity: sha512-W5NH4y/J0plIIS5b2xvTEkU7JFxyqdMAOgf+Ilhl0vHQXKO5dZoxd+C/jEtq56c4F3wk71RB4BMRQ2XdI+bwYQ==}
    engines: {node: '>=16.20.0'}
    cpu: [arm64]
    os: [freebsd]

  '@typescript/typescript-freebsd-x64@7.0.2':
    resolution: {integrity: sha512-UMGDx5sTpzNw3WiPebH7l90IWfJggEd+egHt/q6p7/Cm3zqoV7VxkGXt+3DxPIw8CcmvAB0j3sVVfbhX+M4Tpw==}
    engines: {node: '>=16.20.0'}
    cpu: [x64]
    os: [freebsd]

  '@typescript/typescript-linux-arm64@7.0.2':
    resolution: {integrity: sha512-Qh4eU4/y3yDjnfjjyPYihMj5/ODIlmt+Bzu17OI+fiSRDW57QmU5SiN63exPRNJPKUzcc1INa1NXdrJ+MqHjUQ==}
    engines: {node: '>=16.20.0'}
    cpu: [arm64]
    os: [linux]

  '@typescript/typescript-linux-arm@7.0.2':
    resolution: {integrity: sha512-gffT3xPz9sR7j/YJExkyPntrI0P2EP9XbOyWzth2/Gs0RstK+90RBcO0ncXoXy/beYll1SXw846Nf2zdnEz0QQ==}
    engines: {node: '>=16.20.0'}
    cpu: [arm]
    os: [linux]

  '@typescript/typescript-linux-loong64@7.0.2':
    resolution: {integrity: sha512-uEHck9i8hoAzXPiYRib1O7miOnz23SxIeVl6F4LXox+qov1K35jHcEW6VHKvZI+pyvl7fZEP4MCU5LYvIq1GuQ==}
    engines: {node: '>=16.20.0'}
    cpu: [loong64]
    os: [linux]

  '@typescript/typescript-linux-mips64el@7.0.2':
    resolution: {integrity: sha512-R4KvAMnE43W5Qeqb0Ly56O3mWMWIAgsMyz36DCaycd5nbg/9kzm0liw3JocfRqyJY0KPmzFjbswozXyW0DnIYA==}
    engines: {node: '>=16.20.0'}
    cpu: [mips64el]
    os: [linux]

  '@typescript/typescript-linux-ppc64@7.0.2':
    resolution: {integrity: sha512-DORx5b3sd/4S7eayxm4FQv+A7CrkUIGRaHiwI8oiHTAI1fAPWhF4J0vAlkC8biAlHSVVwxMQ3tjZ2/DVbnQiiA==}
    engines: {node: '>=16.20.0'}
    cpu: [ppc64]
    os: [linux]

  '@typescript/typescript-linux-riscv64@7.0.2':
    resolution: {integrity: sha512-wf0jqEDOjrPRnKwYRyyJDRo11KMbvMFrU+q4zqKyChODBzvlkbhNQfKvLxQCcwTpdDaXSHZTVuh0JoCrKCUMHQ==}
    engines: {node: '>=16.20.0'}
    cpu: [riscv64]
    os: [linux]

  '@typescript/typescript-linux-s390x@7.0.2':
    resolution: {integrity: sha512-IkwJc3L7yhytWd/ewjyxNDfOmswCm9GWMJT/ue/dU4aZNbwZeYAetq42VyLmsmSjvoX7z74X6ZaYCtzAr0EuGw==}
    engines: {node: '>=16.20.0'}
    cpu: [s390x]
    os: [linux]

  '@typescript/typescript-linux-x64@7.0.2':
    resolution: {integrity: sha512-EYdf2cNg7rgCWJnxCdJ+F3V39O8ihb37eHAu1LK8oAFizgTQbPOK7zHHXbPt8rX24COqODXeI3sIf0fCXG7H/A==}
    engines: {node: '>=16.20.0'}
    cpu: [x64]
    os: [linux]

  '@typescript/typescript-netbsd-arm64@7.0.2':
    resolution: {integrity: sha512-+polYF4MF04aPpO5FTkHran9yUQDSXqy5GiSDKpsll5jy3l3+g9QLhpf39T+ePtefhXLOGrLl0QIjkQP6VnelA==}
    engines: {node: '>=16.20.0'}
    cpu: [arm64]
    os: [netbsd]

  '@typescript/typescript-netbsd-x64@7.0.2':
    resolution: {integrity: sha512-8YIT0EHM/3dq10ZOVF/A7pc/YSMtbcecct4rWtexrnSCHOPcpC2KTLXfTCR6vDpnSiY12heNb1GiN/wu+T/FyA==}
    engines: {node: '>=16.20.0'}
    cpu: [x64]
    os: [netbsd]

  '@typescript/typescript-openbsd-arm64@7.0.2':
    resolution: {integrity: sha512-APT8+ClYnuYm1u9+kgGXoMj2VzWzcymwh2gNSQVySHfkRDGOTVkoWLjCmOQSaO+PoqQ57B0flRp9SA+7GnnkzQ==}
    engines: {node: '>=16.20.0'}
    cpu: [arm64]
    os: [openbsd]

  '@typescript/typescript-openbsd-x64@7.0.2':
    resolution: {integrity: sha512-yX7s+Q0Dln0Dt9tEzZsAjXXR/+ytBM7AlglaqyeMPxQszJ1JhlJdZ6jLA+IzldHtflX81em7lDao1xXu+aRRkg==}
    engines: {node: '>=16.20.0'}
    cpu: [x64]
    os: [openbsd]

  '@typescript/typescript-sunos-x64@7.0.2':
    resolution: {integrity: sha512-dLJDGaLZ1D4HPQn62u1n8mBDkJREwMsAkCdkwd4Ieqw+x3TUyTsqY0YiBCtE6H6OzzgGk3iuZ3vFWRS+E8/d1g==}
    engines: {node: '>=16.20.0'}
    cpu: [x64]
    os: [sunos]

  '@typescript/typescript-win32-arm64@7.0.2':
    resolution: {integrity: sha512-Gyl1Vy6OsWesLzmq+EP0Fb7b4Nid5232AvcA2SFcdYreldpNtYFFofPjnt62y9hQy7VTaZp65ICJjuAQRaVcIQ==}
    engines: {node: '>=16.20.0'}
    cpu: [arm64]
    os: [win32]

  '@typescript/typescript-win32-x64@7.0.2':
    resolution: {integrity: sha512-0BQ3HkAHHlKLSp1qRvf3SUhGpGsDuhB/jgFw75guyqbxJqEaS0Cw/VFO8i2nHglJUzQCRtMMR/IBAKE3ETMC4g==}
    engines: {node: '>=16.20.0'}
    cpu: [x64]
    os: [win32]

  '@vitest/expect@4.1.11':
    resolution: {integrity: sha512-VX2x5vNJXET47KAFzwERI+KRMtTTCSWTfSMKsW7JsUsXV4psq++e3DvZpuTDOpHcxytiDs6p2nhVb2tVDiiUYw==}

  '@vitest/mocker@4.1.11':
    resolution: {integrity: sha512-2XJVD55d1o5AZous5CCGKS74g/riOj9odEt2bQpCVZeblHyHdnMeFl4jl0XjU21stf4mbjUkew2eXQZt65g5CQ==}
    peerDependencies:
      msw: ^2.4.9
      vite: ^6.0.0 || ^7.0.0 || ^8.0.0
    peerDependenciesMeta:
      msw:
        optional: true
      vite:
        optional: true

  '@vitest/pretty-format@4.1.11':
    resolution: {integrity: sha512-yiZzPbGTS9Sr/JpFl8zHrcIkAofNbFV6k21vIgQN/cY/oxZeXhJv5sc/MBJ5jFKWmWs+oJHw0UXLZjmf931+Vw==}

  '@vitest/runner@4.1.11':
    resolution: {integrity: sha512-LztvUgdwMNJMIkj3hQnnxiC2Xy1zNxq928W/xhjCLaNCzqTZOudjwbQf6v9IntZGPw132i2Lq2rgTRZHD3JHNw==}

  '@vitest/snapshot@4.1.11':
    resolution: {integrity: sha512-pN7ikn1ON7h8ee4gIAp4AzyK+zBtJPzVbqOgu5LCEh4VaJVbPQcgYQYJIMGQPXVeJJq1fnfazis7a5pFNPahog==}

  '@vitest/spy@4.1.11':
    resolution: {integrity: sha512-apNa/prQy2qCeywhnixOHPRCgGNhvg7T4Dapfl1GahLp/R+uhBm5cPyFoNVyqsNd2h1nJxL6BqqdIjiABL60YA==}

  '@vitest/utils@4.1.11':
    resolution: {integrity: sha512-zTCVGpyFsGWBhllOyKlTw/vnr6D9qxsfSDyfbyZmTyjHw5N/VuvzHpHoQjm2ZJzn4RJgx5w4r7V0er69CmLgPQ==}

  assertion-error@2.0.1:
    resolution: {integrity: sha512-Izi8RQcffqCeNVgFigKli1ssklIbpHnCYc6AknXGYoB6grJqyeby7jv12JUQgmTAnIDnbck1uxksT4dzN3PWBA==}
    engines: {node: '>=12'}

  baseline-browser-mapping@2.11.18:
    resolution: {integrity: sha512-1iEmLEYSiE1SeBoAfPo/Mnx3PzfzHUkDK61ASkCpuk3YXugYLH5DYK1SzqV55F8FMI6s0F+/tCP7Polz1QRjxw==}
    engines: {node: '>=6.0.0'}
    hasBin: true

  caniuse-lite@1.0.30001809:
    resolution: {integrity: sha512-xxWVywk6a6Arlk+hymeycyn/VgqEfLDxupvhH/xiY5SJ/18kmi9o6MiO320DCUzypORHLtvh0I4i04tUhCNHNQ==}

  chai@6.2.2:
    resolution: {integrity: sha512-NUPRluOfOiTKBKvWPtSD4PhFvWCqOi0BGStNWs57X9js7XGTprSmFoz5F0tWhR4WPjNeR9jXqdC7/UpSJTnlRg==}
    engines: {node: '>=18'}

  client-only@0.0.1:
    resolution: {integrity: sha512-IV3Ou0jSMzZrd3pZ48nLkT9DA7Ag1pnPzaiQhpW7c3RbcqqzvzzVu+L8gfqMp/8IM2MQtSiqaCxrrcfu8I8rMA==}

  convert-source-map@2.0.0:
    resolution: {integrity: sha512-Kvp459HrV2FEJ1CAsi1Ku+MY3kasH19TFykTz2xWmMeq6bk2NU3XXvfJ+Q61m0xktWwt+1HSYf3JZsTms3aRJg==}

  csp_evaluator@1.1.8:
    resolution: {integrity: sha512-EwOnfYuNbTytvbMKsLixTrRgnjOa0WZCxGy8A9nnSYAicrdwn+T/epU/yjgymmOxlgKnvH+8wXt+7p/8ak5Feg==}

  csstype@3.2.3:
    resolution: {integrity: sha512-z1HGKcYy2xA8AGQfwrn0PAy+PB7X/GSj3UVJW9qKyn43xWa+gl5nXmU4qqLMRzWVLFC8KusUX8T/0kCiOYpAIQ==}

  detect-libc@2.1.2:
    resolution: {integrity: sha512-Btj2BOOO83o3WyH59e8MgXsxEQVcarkUOpEYrubB0urwnN10yQ364rsiByU11nZlqWYZm05i/of7io4mzihBtQ==}
    engines: {node: '>=8'}

  es-module-lexer@2.3.2:
    resolution: {integrity: sha512-poHGpORABojJJucnV9KbOavETW8lBVnphkW77ER5/BQ5Fz7oXSoCNek7IH3vR5nRjdsEz926ibFYX8KtLQmdyw==}

  estree-walker@3.0.3:
    resolution: {integrity: sha512-7RUKfXgSMMkzt6ZuXmqapOurLGPPfgj6l9uRZ7lRGolvk0y2yocc35LdcxKC5PQZdn2DMqioAQ2NoWcrTKmm6g==}

  expect-type@1.4.0:
    resolution: {integrity: sha512-KfYbmpRm0VbLjEvVa9yGwCi9GI34xvi7A/HXYWQO65CSD2u3MczUJSuwXKFIxlGsgBQizV9q5J9NHj4VG0n+pA==}
    engines: {node: '>=12.0.0'}

  fdir@6.5.0:
    resolution: {integrity: sha512-tIbYtZbucOs0BRGqPJkshJUYdL+SDH7dVM8gjy+ERp3WAUjLEFJE+02kanyHtwjWOnwrKYBiwAmM0p4kLJAnXg==}
    engines: {node: '>=12.0.0'}
    peerDependencies:
      picomatch: ^3 || ^4
    peerDependenciesMeta:
      picomatch:
        optional: true

  fsevents@2.3.3:
    resolution: {integrity: sha512-5xoDfX+fL7faATnagmWPpbFtwh/R77WmMMqqHGS65C3vvB0YHrgF+B1YmZ3441tMj5n63k0212XNoJwzlhffQw==}
    engines: {node: ^8.16.0 || ^10.6.0 || >=11.0.0}
    os: [darwin]

  jose@6.2.10:
    resolution: {integrity: sha512-iiW7J9qRFlGxvCOIBDBDxFePQSn7ZMAnrYGhrrOo6siO/MIqwfyilLR27pkfDgUk+raLuzADS8A3S/KLBisc0g==}

  lightningcss-android-arm64@1.33.0:
    resolution: {integrity: sha512-gEpRTalKdosp4Bb8qWtc2iOgE5SeIHlpS1up9bFq2wAyYhl1UdTObYiHe98zEM9SQvSoqQZ1IQD0JNpg3Ml5pg==}
    engines: {node: '>= 12.0.0'}
    cpu: [arm64]
    os: [android]

  lightningcss-darwin-arm64@1.33.0:
    resolution: {integrity: sha512-Sciaz8eenNTKn9b3t7+xr0ipTp9YxKQY4npwQ3mrRuL0BAVHBLyZxofhaKBAVtzmtRZ/zTyo0/to4B1uWG/Djg==}
    engines: {node: '>= 12.0.0'}
    cpu: [arm64]
    os: [darwin]

  lightningcss-darwin-x64@1.33.0:
    resolution: {integrity: sha512-Z5UPAxzrjlWNNyGy6i65cJzzvgJ5D3T6wMvs+gWpY9d7qRhANrxqAp6LhxIgZhWEw18RfJTGcRxjuLIBr+m8XQ==}
    engines: {node: '>= 12.0.0'}
    cpu: [x64]
    os: [darwin]

  lightningcss-freebsd-x64@1.33.0:
    resolution: {integrity: sha512-QQM/Ti/hQajJwCY+RiWuCZ9sdtI/XQk7nDK5vC8kkdwixezOlDgvDx7+RT+QjK6FcFT4MpsuoBnHIo/O3StRRg==}
    engines: {node: '>= 12.0.0'}
    cpu: [x64]
    os: [freebsd]

  lightningcss-linux-arm-gnueabihf@1.33.0:
    resolution: {integrity: sha512-N7FVBe6iS24MlM6R/4RBTxGhQheZGs7tiQ9U32UtF75NzP5Q7xWPRqLBCKxlRQRk3rY1jCIPLzx7WzOhuUIRLQ==}
    engines: {node: '>= 12.0.0'}
    cpu: [arm]
    os: [linux]

  lightningcss-linux-arm64-gnu@1.33.0:
    resolution: {integrity: sha512-j2v/itmy4HlNxlc6voKXYgBqNi0Ng2LShg4z7GufpEgs05P+2suBVyi9I6YHq5uoVFx9ETin3eCEhLVyXGQnKg==}
    engines: {node: '>= 12.0.0'}
    cpu: [arm64]
    os: [linux]
    libc: [glibc]

  lightningcss-linux-arm64-musl@1.33.0:
    resolution: {integrity: sha512-yiO5ROMuYQgXbC60yjZU5CYSFZGKXL0HFATXt9mHJn1+zW55oCtMI9NfcVhYLMFDL7gV7oBPon/EmMMGg2OvtQ==}
    engines: {node: '>= 12.0.0'}
    cpu: [arm64]
    os: [linux]
    libc: [musl]

  lightningcss-linux-x64-gnu@1.33.0:
    resolution: {integrity: sha512-ar+Ju7LmcN0Jo4FpL4hpFybwNG9/3A/Br5KW2n2jyODg3MEZXaDYADdemoNS+BDNfMgKvylJLj4S5tyRActuAg==}
    engines: {node: '>= 12.0.0'}
    cpu: [x64]
    os: [linux]
    libc: [glibc]

  lightningcss-linux-x64-musl@1.33.0:
    resolution: {integrity: sha512-RYiYbkokw0trfKqqzfF55lginwEPrD3OJDfTuJzFs1MK6iFnDenaz1fqLLtX4ITG3OktJQXOeTaw1awrBAlZPw==}
    engines: {node: '>= 12.0.0'}
    cpu: [x64]
    os: [linux]
    libc: [musl]

  lightningcss-win32-arm64-msvc@1.33.0:
    resolution: {integrity: sha512-1K+MPfLSFVpphzpdbfkhlWk6wBrTObBzS2T6db10PNOZgR9GoVsAWzwNyuhUYYbTp23j+4RrncfujZ4uAzXvwA==}
    engines: {node: '>= 12.0.0'}
    cpu: [arm64]
    os: [win32]

  lightningcss-win32-x64-msvc@1.33.0:
    resolution: {integrity: sha512-OlEICDx/Xl0FqSp4bry8zFnCvGpig3Gl4gCquvYwHuqJKEC1+n9NgDniFvqHGmMv1ZkqDJrDqKKSykTDX+ehuA==}
    engines: {node: '>= 12.0.0'}
    cpu: [x64]
    os: [win32]

  lightningcss@1.33.0:
    resolution: {integrity: sha512-WkUDrojuJs0xkgGf2udWxa3yGBRxPtxUkB79i6aCZLRgc7PM8fZe9TosfPDcvEpQZbuFASnHYmRLBLUbmLOIIA==}
    engines: {node: '>= 12.0.0'}

  magic-string@0.30.21:
    resolution: {integrity: sha512-vd2F4YUyEXKGcLHoq+TEyCjxueSeHnFxyyjNp80yg0XV4vUhnDer/lvvlqM/arB5bXQN5K2/3oinyCRyx8T2CQ==}

  nanoid@3.3.18:
    resolution: {integrity: sha512-DTg4MJbGMWkfi6VZFdNt2/caMbQy4Ou+Op/hJQvGEWcnVfoA1QA+xzRKAzw9jD6+GVOOeYr/mIcuDSdug6F6+w==}
    engines: {node: ^10 || ^12 || ^13.7 || ^14 || >=15.0.1}
    hasBin: true

  next@16.3.4:
    resolution: {integrity: sha512-/Ztf6CeRH+ejEXUrYtqI4gkS66eFIHuSwqi60RgcpWKodxFZx2/dqVCMKBwILfAHXQ+F1b1vAudgj3mnxqtoIA==}
    engines: {node: '>=20.9.0'}
    hasBin: true
    peerDependencies:
      '@opentelemetry/api': ^1.1.0
      '@playwright/test': ^1.51.1
      babel-plugin-react-compiler: '*'
      react: ^18.2.0 || 19.0.0-rc-de68d2f4-20241204 || ^19.0.0
      react-dom: ^18.2.0 || 19.0.0-rc-de68d2f4-20241204 || ^19.0.0
      sass: ^1.3.0
    peerDependenciesMeta:
      '@opentelemetry/api':
        optional: true
      '@playwright/test':
        optional: true
      babel-plugin-react-compiler:
        optional: true
      sass:
        optional: true

  oauth4webapi@3.8.7:
    resolution: {integrity: sha512-4RxcKxXjuItDFZ20RRPf4YTw3kpeXJyCgJFxVzJ068A7PNJ18st2Dg90tlC1LkSDS0GecroagCLHYEIVUhCAkw==}

  obug@2.1.4:
    resolution: {integrity: sha512-4a+OsYv9UktOJKE+l1A4OufDgdRF9PifWj+tJnHURo/P+WOxpG4GzUFL9qCalmWauao6ogiG+QvnCovwPoyAWA==}
    engines: {node: '>=12.20.0'}

  openid-client@6.8.5:
    resolution: {integrity: sha512-jNGC/5wnTYwCcEUe2ss0IRUmVRQcgxM0A1nLb3eX/9llqNbMWOQd2xd+qDAgfVCpA5Qh96Y1cdnkfbva6+bSdA==}

  pathe@2.0.3:
    resolution: {integrity: sha512-WUjGcAqP1gQacoQe+OBJsFA7Ld4DyXuUIjZ5cc75cLHvJ7dtNsTugphxIADwspS+AraAUePCKrSVtPLFj/F88w==}

  picocolors@1.1.1:
    resolution: {integrity: sha512-xceH2snhtb5M9liqDsmEw56le376mTZkEX/jEb/RxNFyegNul7eNslCXP9FDj/Lcu0X8KEyMceP2ntpaHrDEVA==}

  picomatch@4.0.5:
    resolution: {integrity: sha512-RvwwcruNjI1ncT5xRakeyS9Lf8lcItv34KD+aif+VH9kduAyfYBipGh12274xtenIPZ119/R9BdTBa8gAwSh0A==}
    engines: {node: '>=12'}

  postcss@8.5.23:
    resolution: {integrity: sha512-g50586zr4bZmwFiTlflMu8E0bDTb5I5gertgwAKmsdUlTQIhZtunzUlD1WSzwcVWPoAVpsrA6vlfCD7oXvRwgg==}
    engines: {node: ^10 || ^12 || >=14}

  postcss@8.5.26:
    resolution: {integrity: sha512-u82N74LFzG8ca+dD8puPnplTXoGH4fTPpVGuIbt36G3qvNlkvfD0lEAZSxaly3KX8TS/L1A1gsCEmvKmBcVbkQ==}
    engines: {node: ^10 || ^12 || >=14}

  react-dom@19.2.8:
    resolution: {integrity: sha512-rVprimfGBG3DR+Tq0IQG2DT5PxKth1WIGDmj5yPmlzr4YBe7uyE+Du4oVqTDXZSHGGGXRtTJEGSSePyQCMBglQ==}
    peerDependencies:
      react: ^19.2.8

  react@19.2.8:
    resolution: {integrity: sha512-PWaYA1L/q9u2u7xYQi+Y3L3Yfnie7XyLeaJICV1MGD6LprsBxcAqGjYyr0eY3p+QdsA+x/Irkt4Qif8D63+Sbw==}
    engines: {node: '>=0.10.0'}

  rolldown@1.2.5:
    resolution: {integrity: sha512-VD2IE5PUG4Oj8zz2VGykiYd5wbnjdIiSsNQb8Qu5B+noEp+A78mu2iVvpp27g8es14Tk9rofNs5Tku9iQCS4fA==}
    engines: {node: ^20.19.0 || >=22.12.0}
    hasBin: true

  safevalues@1.2.0:
    resolution: {integrity: sha512-zIsuhjYvJCjfsfjoim2ab6gLKFYAnTiDSJGh0cC3T44L/4kNLL90hBG2BzrXPrHA3f8Ms8FSJ1mljKH5dVR1cw==}

  scheduler@0.27.0:
    resolution: {integrity: sha512-eNv+WrVbKu1f3vbYJT/xtiF5syA5HPIMtf9IgY/nKg0sWqzAUEvqY/xm7OcZc/qafLx/iO9FgOmeSAp4v5ti/Q==}

  server-only@0.0.1:
    resolution: {integrity: sha512-qepMx2JxAa5jjfzxG79yPPq+8BuFToHd1hm7kI+Z4zAq1ftQiP7HcxMhDDItrbtwVeLg/cY2JnKnrcFkmiswNA==}

  siginfo@2.0.0:
    resolution: {integrity: sha512-ybx0WO1/8bSBLEWXZvEd7gMW3Sn3JFlW3TvX1nREbDLRNQNaeNN8WK0meBwPdAaOI7TtRRRJn/Es1zhrrCHu7g==}

  source-map-js@1.2.1:
    resolution: {integrity: sha512-UXWMKhLOwVKb728IUtQPXxfYU+usdybtUrK/8uGE8CQMvrhOpwvzDBwj0QhSL7MQc7vIsISBG8VQ8+IDQxpfQA==}
    engines: {node: '>=0.10.0'}

  stackback@0.0.2:
    resolution: {integrity: sha512-1XMJE5fQo1jGH6Y/7ebnwPOBEkIEnT4QF32d5R1+VXdXveM0IBMJt8zfaxX1P3QhVwrYe+576+jkANtSS2mBbw==}

  std-env@4.2.0:
    resolution: {integrity: sha512-oCUKSupKTHX53EyjDtuZQ64pjLJ6yYCtpmEw0goYxtjG9KpbRe8KAsl2tBUGU9DyMcJ0RwJ8GqJAFzMXcXW1Rw==}

  styled-jsx@5.1.6:
    resolution: {integrity: sha512-qSVyDTeMotdvQYoHWLNGwRFJHC+i+ZvdBRYosOFgC+Wg1vx4frN2/RG/NA7SYqqvKNLf39P2LSRA2pu6n0XYZA==}
    engines: {node: '>= 12.0.0'}
    peerDependencies:
      '@babel/core': '*'
      babel-plugin-macros: '*'
      react: '>= 16.8.0 || 17.x.x || ^18.0.0-0 || ^19.0.0-0'
    peerDependenciesMeta:
      '@babel/core':
        optional: true
      babel-plugin-macros:
        optional: true

  tinybench@2.9.0:
    resolution: {integrity: sha512-0+DUvqWMValLmha6lr4kD8iAMK1HzV0/aKnCtWb9v9641TnP/MFb7Pc2bxoxQjTXAErryXVgUOfv2YqNllqGeg==}

  tinyexec@1.3.0:
    resolution: {integrity: sha512-QKAl9m8gWWGHV8jZcPeym6j+XULi6tOf1mT83WYJ4Lk2ytW/uwAWkrP0uFsdoYMdueVJ0qs26wZ+23xeB4ibNQ==}
    engines: {node: '>=18'}

  tinyglobby@0.2.17:
    resolution: {integrity: sha512-wXR/dYpcqKmfWpEdZjiKJOwCNFndD0DMnrW/cYjVGttEkBfVgcLFHoNrlj47mjOVic9yyNu65alsgF4NQyTa2g==}
    engines: {node: '>=12.0.0'}

  tinyrainbow@3.1.1:
    resolution: {integrity: sha512-yau8yJdTt989Mm0Bd/236QnzEiPf2xLLTqUZRUJOo/3CB078LSwzei343DgtJVmfJKJE3TMINY1u42SQsP6mXw==}
    engines: {node: '>=14.0.0'}

  tslib@2.8.1:
    resolution: {integrity: sha512-oJFu94HQb+KVduSUQL7wnpmqnfmLsOA/nAh6b6EH0wCEoK0/mPeXU6c3wKDV83MkOuHPRHtSXKKU99IBazS/2w==}

  typescript@7.0.2:
    resolution: {integrity: sha512-8FYau96o3NKOhbjKi/qNvG/W5jhzxkbdm5sj9AbZ/5T5sWqn3hJgLfGx27sRKZWTvyzCP8dLRBTf5tBTSRVUNA==}
    engines: {node: '>=16.20.0'}
    hasBin: true

  undici-types@8.3.0:
    resolution: {integrity: sha512-j375ScV60dom+YkPFIfTLcOiPxkN/buHz5GobjLhixFuANaNs3C9l4GmrWqejgXWJ7BbJcFYpTEUkS1Ge8bpZQ==}

  vite@8.2.2:
    resolution: {integrity: sha512-cFKLV/PRgAUlIRm5WjMjJ86jrftzpqcgH+Us+DS8mI3CDNiH30Whrz8uHL3+MOLPAgqbMBAqWdAHAphOAM+z/Q==}
    engines: {node: ^20.19.0 || >=22.12.0}
    hasBin: true
    peerDependencies:
      '@types/node': ^20.19.0 || >=22.12.0
      '@vitejs/devtools': ^0.4.0 || ^0.5.0
      esbuild: ^0.27.0 || ^0.28.0
      jiti: '>=1.21.0'
      less: ^4.0.0
      sass: ^1.70.0
      sass-embedded: ^1.70.0
      stylus: '>=0.54.8'
      sugarss: ^5.0.0
      terser: ^5.16.0
      tsx: ^4.8.1
      yaml: ^2.4.2
    peerDependenciesMeta:
      '@types/node':
        optional: true
      '@vitejs/devtools':
        optional: true
      esbuild:
        optional: true
      jiti:
        optional: true
      less:
        optional: true
      sass:
        optional: true
      sass-embedded:
        optional: true
      stylus:
        optional: true
      sugarss:
        optional: true
      terser:
        optional: true
      tsx:
        optional: true
      yaml:
        optional: true

  vitest@4.1.11:
    resolution: {integrity: sha512-fhACrNXUidIbGSBr5FlbuBkO7VWC1ZyLl0DO4CU2DrQoAPxX84Ysxs+HeGQpii5lZWV1Q4gBZTTu49mF+A6Edw==}
    engines: {node: ^20.0.0 || ^22.0.0 || >=24.0.0}
    hasBin: true
    peerDependencies:
      '@edge-runtime/vm': '*'
      '@opentelemetry/api': ^1.9.0
      '@types/node': ^20.0.0 || ^22.0.0 || >=24.0.0
      '@vitest/browser-playwright': 4.1.11
      '@vitest/browser-preview': 4.1.11
      '@vitest/browser-webdriverio': 4.1.11
      '@vitest/coverage-istanbul': 4.1.11
      '@vitest/coverage-v8': 4.1.11
      '@vitest/ui': 4.1.11
      happy-dom: '*'
      jsdom: '*'
      vite: ^6.0.0 || ^7.0.0 || ^8.0.0
    peerDependenciesMeta:
      '@edge-runtime/vm':
        optional: true
      '@opentelemetry/api':
        optional: true
      '@types/node':
        optional: true
      '@vitest/browser-playwright':
        optional: true
      '@vitest/browser-preview':
        optional: true
      '@vitest/browser-webdriverio':
        optional: true
      '@vitest/coverage-istanbul':
        optional: true
      '@vitest/coverage-v8':
        optional: true
      '@vitest/ui':
        optional: true
      happy-dom:
        optional: true
      jsdom:
        optional: true

  why-is-node-running@2.3.0:
    resolution: {integrity: sha512-hUrmaWBdVDcxvYqnyh09zunKzROWjbZTiNy8dBEjkS7ehEDQibXJ7XvlmtbwuTclUiIyN+CyXQD4Vmko8fNm8w==}
    engines: {node: '>=8'}
    hasBin: true

  zod@4.4.3:
    resolution: {integrity: sha512-ytENFjIJFl2UwYglde2jchW2Hwm4GJFLDiSXWdTrJQBIN9Fcyp7n4DhxJEiWNAJMV1/BqWfW/kkg71UDcHJyTQ==}

ignoredOptionalDependencies:
  - sharp

snapshots:

  '@jridgewell/sourcemap-codec@1.5.5': {}

  '@next/env@16.3.4': {}

  '@next/swc-darwin-arm64@16.3.4':
    optional: true

  '@next/swc-darwin-x64@16.3.4':
    optional: true

  '@next/swc-linux-arm64-gnu@16.3.4':
    optional: true

  '@next/swc-linux-arm64-musl@16.3.4':
    optional: true

  '@next/swc-linux-x64-gnu@16.3.4':
    optional: true

  '@next/swc-linux-x64-musl@16.3.4':
    optional: true

  '@next/swc-win32-arm64-msvc@16.3.4':
    optional: true

  '@next/swc-win32-x64-msvc@16.3.4':
    optional: true

  '@oxc-project/types@0.146.0': {}

  '@rolldown/binding-android-arm-eabi@1.2.5':
    optional: true

  '@rolldown/binding-android-arm64@1.2.5':
    optional: true

  '@rolldown/binding-darwin-arm64@1.2.5':
    optional: true

  '@rolldown/binding-darwin-x64@1.2.5':
    optional: true

  '@rolldown/binding-freebsd-x64@1.2.5':
    optional: true

  '@rolldown/binding-linux-arm-gnueabihf@1.2.5':
    optional: true

  '@rolldown/binding-linux-arm64-gnu@1.2.5':
    optional: true

  '@rolldown/binding-linux-arm64-musl@1.2.5':
    optional: true

  '@rolldown/binding-linux-ppc64-gnu@1.2.5':
    optional: true

  '@rolldown/binding-linux-s390x-gnu@1.2.5':
    optional: true

  '@rolldown/binding-linux-x64-gnu@1.2.5':
    optional: true

  '@rolldown/binding-linux-x64-musl@1.2.5':
    optional: true

  '@rolldown/binding-openharmony-arm64@1.2.5':
    optional: true

  '@rolldown/binding-win32-arm64-msvc@1.2.5':
    optional: true

  '@rolldown/binding-win32-x64-msvc@1.2.5':
    optional: true

  '@rolldown/pluginutils@1.0.1': {}

  '@standard-schema/spec@1.1.0': {}

  '@swc/helpers@0.5.23':
    dependencies:
      tslib: 2.8.1

  '@types/chai@5.2.3':
    dependencies:
      '@types/deep-eql': 4.0.2
      assertion-error: 2.0.1

  '@types/deep-eql@4.0.2': {}

  '@types/estree@1.0.9': {}

  '@types/node@26.2.0':
    dependencies:
      undici-types: 8.3.0

  '@types/react-dom@19.2.5(@types/react@19.2.18)':
    dependencies:
      '@types/react': 19.2.18

  '@types/react@19.2.18':
    dependencies:
      csstype: 3.2.3

  '@typescript/typescript-aix-ppc64@7.0.2':
    optional: true

  '@typescript/typescript-darwin-arm64@7.0.2':
    optional: true

  '@typescript/typescript-darwin-x64@7.0.2':
    optional: true

  '@typescript/typescript-freebsd-arm64@7.0.2':
    optional: true

  '@typescript/typescript-freebsd-x64@7.0.2':
    optional: true

  '@typescript/typescript-linux-arm64@7.0.2':
    optional: true

  '@typescript/typescript-linux-arm@7.0.2':
    optional: true

  '@typescript/typescript-linux-loong64@7.0.2':
    optional: true

  '@typescript/typescript-linux-mips64el@7.0.2':
    optional: true

  '@typescript/typescript-linux-ppc64@7.0.2':
    optional: true

  '@typescript/typescript-linux-riscv64@7.0.2':
    optional: true

  '@typescript/typescript-linux-s390x@7.0.2':
    optional: true

  '@typescript/typescript-linux-x64@7.0.2':
    optional: true

  '@typescript/typescript-netbsd-arm64@7.0.2':
    optional: true

  '@typescript/typescript-netbsd-x64@7.0.2':
    optional: true

  '@typescript/typescript-openbsd-arm64@7.0.2':
    optional: true

  '@typescript/typescript-openbsd-x64@7.0.2':
    optional: true

  '@typescript/typescript-sunos-x64@7.0.2':
    optional: true

  '@typescript/typescript-win32-arm64@7.0.2':
    optional: true

  '@typescript/typescript-win32-x64@7.0.2':
    optional: true

  '@vitest/expect@4.1.11':
    dependencies:
      '@standard-schema/spec': 1.1.0
      '@types/chai': 5.2.3
      '@vitest/spy': 4.1.11
      '@vitest/utils': 4.1.11
      chai: 6.2.2
      tinyrainbow: 3.1.1

  '@vitest/mocker@4.1.11(vite@8.2.2(@types/node@26.2.0))':
    dependencies:
      '@vitest/spy': 4.1.11
      estree-walker: 3.0.3
      magic-string: 0.30.21
    optionalDependencies:
      vite: 8.2.2(@types/node@26.2.0)

  '@vitest/pretty-format@4.1.11':
    dependencies:
      tinyrainbow: 3.1.1

  '@vitest/runner@4.1.11':
    dependencies:
      '@vitest/utils': 4.1.11
      pathe: 2.0.3

  '@vitest/snapshot@4.1.11':
    dependencies:
      '@vitest/pretty-format': 4.1.11
      '@vitest/utils': 4.1.11
      magic-string: 0.30.21
      pathe: 2.0.3

  '@vitest/spy@4.1.11': {}

  '@vitest/utils@4.1.11':
    dependencies:
      '@vitest/pretty-format': 4.1.11
      convert-source-map: 2.0.0
      tinyrainbow: 3.1.1

  assertion-error@2.0.1: {}

  baseline-browser-mapping@2.11.18: {}

  caniuse-lite@1.0.30001809: {}

  chai@6.2.2: {}

  client-only@0.0.1: {}

  convert-source-map@2.0.0: {}

  csp_evaluator@1.1.8: {}

  csstype@3.2.3: {}

  detect-libc@2.1.2: {}

  es-module-lexer@2.3.2: {}

  estree-walker@3.0.3:
    dependencies:
      '@types/estree': 1.0.9

  expect-type@1.4.0: {}

  fdir@6.5.0(picomatch@4.0.5):
    optionalDependencies:
      picomatch: 4.0.5

  fsevents@2.3.3:
    optional: true

  jose@6.2.10: {}

  lightningcss-android-arm64@1.33.0:
    optional: true

  lightningcss-darwin-arm64@1.33.0:
    optional: true

  lightningcss-darwin-x64@1.33.0:
    optional: true

  lightningcss-freebsd-x64@1.33.0:
    optional: true

  lightningcss-linux-arm-gnueabihf@1.33.0:
    optional: true

  lightningcss-linux-arm64-gnu@1.33.0:
    optional: true

  lightningcss-linux-arm64-musl@1.33.0:
    optional: true

  lightningcss-linux-x64-gnu@1.33.0:
    optional: true

  lightningcss-linux-x64-musl@1.33.0:
    optional: true

  lightningcss-win32-arm64-msvc@1.33.0:
    optional: true

  lightningcss-win32-x64-msvc@1.33.0:
    optional: true

  lightningcss@1.33.0:
    dependencies:
      detect-libc: 2.1.2
    optionalDependencies:
      lightningcss-android-arm64: 1.33.0
      lightningcss-darwin-arm64: 1.33.0
      lightningcss-darwin-x64: 1.33.0
      lightningcss-freebsd-x64: 1.33.0
      lightningcss-linux-arm-gnueabihf: 1.33.0
      lightningcss-linux-arm64-gnu: 1.33.0
      lightningcss-linux-arm64-musl: 1.33.0
      lightningcss-linux-x64-gnu: 1.33.0
      lightningcss-linux-x64-musl: 1.33.0
      lightningcss-win32-arm64-msvc: 1.33.0
      lightningcss-win32-x64-msvc: 1.33.0

  magic-string@0.30.21:
    dependencies:
      '@jridgewell/sourcemap-codec': 1.5.5

  nanoid@3.3.18: {}

  next@16.3.4(@types/node@26.2.0)(react-dom@19.2.8(react@19.2.8))(react@19.2.8):
    dependencies:
      '@next/env': 16.3.4
      '@swc/helpers': 0.5.23
      baseline-browser-mapping: 2.11.18
      caniuse-lite: 1.0.30001809
      postcss: 8.5.23
      react: 19.2.8
      react-dom: 19.2.8(react@19.2.8)
      styled-jsx: 5.1.6(react@19.2.8)
    optionalDependencies:
      '@next/swc-darwin-arm64': 16.3.4
      '@next/swc-darwin-x64': 16.3.4
      '@next/swc-linux-arm64-gnu': 16.3.4
      '@next/swc-linux-arm64-musl': 16.3.4
      '@next/swc-linux-x64-gnu': 16.3.4
      '@next/swc-linux-x64-musl': 16.3.4
      '@next/swc-win32-arm64-msvc': 16.3.4
      '@next/swc-win32-x64-msvc': 16.3.4
    transitivePeerDependencies:
      - '@babel/core'
      - '@types/node'
      - babel-plugin-macros

  oauth4webapi@3.8.7: {}

  obug@2.1.4: {}

  openid-client@6.8.5:
    dependencies:
      jose: 6.2.10
      oauth4webapi: 3.8.7

  pathe@2.0.3: {}

  picocolors@1.1.1: {}

  picomatch@4.0.5: {}

  postcss@8.5.23:
    dependencies:
      nanoid: 3.3.18
      picocolors: 1.1.1
      source-map-js: 1.2.1

  postcss@8.5.26:
    dependencies:
      nanoid: 3.3.18
      picocolors: 1.1.1
      source-map-js: 1.2.1

  react-dom@19.2.8(react@19.2.8):
    dependencies:
      react: 19.2.8
      scheduler: 0.27.0

  react@19.2.8: {}

  rolldown@1.2.5:
    dependencies:
      '@oxc-project/types': 0.146.0
      '@rolldown/pluginutils': 1.0.1
    optionalDependencies:
      '@rolldown/binding-android-arm-eabi': 1.2.5
      '@rolldown/binding-android-arm64': 1.2.5
      '@rolldown/binding-darwin-arm64': 1.2.5
      '@rolldown/binding-darwin-x64': 1.2.5
      '@rolldown/binding-freebsd-x64': 1.2.5
      '@rolldown/binding-linux-arm-gnueabihf': 1.2.5
      '@rolldown/binding-linux-arm64-gnu': 1.2.5
      '@rolldown/binding-linux-arm64-musl': 1.2.5
      '@rolldown/binding-linux-ppc64-gnu': 1.2.5
      '@rolldown/binding-linux-s390x-gnu': 1.2.5
      '@rolldown/binding-linux-x64-gnu': 1.2.5
      '@rolldown/binding-linux-x64-musl': 1.2.5
      '@rolldown/binding-openharmony-arm64': 1.2.5
      '@rolldown/binding-win32-arm64-msvc': 1.2.5
      '@rolldown/binding-win32-x64-msvc': 1.2.5

  safevalues@1.2.0: {}

  scheduler@0.27.0: {}

  server-only@0.0.1: {}

  siginfo@2.0.0: {}

  source-map-js@1.2.1: {}

  stackback@0.0.2: {}

  std-env@4.2.0: {}

  styled-jsx@5.1.6(react@19.2.8):
    dependencies:
      client-only: 0.0.1
      react: 19.2.8

  tinybench@2.9.0: {}

  tinyexec@1.3.0: {}

  tinyglobby@0.2.17:
    dependencies:
      fdir: 6.5.0(picomatch@4.0.5)
      picomatch: 4.0.5

  tinyrainbow@3.1.1: {}

  tslib@2.8.1: {}

  typescript@7.0.2:
    optionalDependencies:
      '@typescript/typescript-aix-ppc64': 7.0.2
      '@typescript/typescript-darwin-arm64': 7.0.2
      '@typescript/typescript-darwin-x64': 7.0.2
      '@typescript/typescript-freebsd-arm64': 7.0.2
      '@typescript/typescript-freebsd-x64': 7.0.2
      '@typescript/typescript-linux-arm': 7.0.2
      '@typescript/typescript-linux-arm64': 7.0.2
      '@typescript/typescript-linux-loong64': 7.0.2
      '@typescript/typescript-linux-mips64el': 7.0.2
      '@typescript/typescript-linux-ppc64': 7.0.2
      '@typescript/typescript-linux-riscv64': 7.0.2
      '@typescript/typescript-linux-s390x': 7.0.2
      '@typescript/typescript-linux-x64': 7.0.2
      '@typescript/typescript-netbsd-arm64': 7.0.2
      '@typescript/typescript-netbsd-x64': 7.0.2
      '@typescript/typescript-openbsd-arm64': 7.0.2
      '@typescript/typescript-openbsd-x64': 7.0.2
      '@typescript/typescript-sunos-x64': 7.0.2
      '@typescript/typescript-win32-arm64': 7.0.2
      '@typescript/typescript-win32-x64': 7.0.2

  undici-types@8.3.0: {}

  vite@8.2.2(@types/node@26.2.0):
    dependencies:
      lightningcss: 1.33.0
      picomatch: 4.0.5
      postcss: 8.5.26
      rolldown: 1.2.5
      tinyglobby: 0.2.17
    optionalDependencies:
      '@types/node': 26.2.0
      fsevents: 2.3.3

  vitest@4.1.11(@types/node@26.2.0)(vite@8.2.2(@types/node@26.2.0)):
    dependencies:
      '@vitest/expect': 4.1.11
      '@vitest/mocker': 4.1.11(vite@8.2.2(@types/node@26.2.0))
      '@vitest/pretty-format': 4.1.11
      '@vitest/runner': 4.1.11
      '@vitest/snapshot': 4.1.11
      '@vitest/spy': 4.1.11
      '@vitest/utils': 4.1.11
      es-module-lexer: 2.3.2
      expect-type: 1.4.0
      magic-string: 0.30.21
      obug: 2.1.4
      pathe: 2.0.3
      picomatch: 4.0.5
      std-env: 4.2.0
      tinybench: 2.9.0
      tinyexec: 1.3.0
      tinyglobby: 0.2.17
      tinyrainbow: 3.1.1
      vite: 8.2.2(@types/node@26.2.0)
      why-is-node-running: 2.3.0
    optionalDependencies:
      '@types/node': 26.2.0
    transitivePeerDependencies:
      - msw

  why-is-node-running@2.3.0:
    dependencies:
      siginfo: 2.0.0
      stackback: 0.0.2

  zod@4.4.3: {}
````

### FILE: `pnpm-workspace.yaml`

```yaml
block_id: "TS-GO-BRIDGE:pnpm-workspace:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified standalone BFF"
license: "LicenseRef-Workspace-Owner"
sha256: "f52887318fe556dbdb3574622f5c6248483e843edc3b61ad98b6ae3b351571ba"
variables: []
secrets_allowed: false
```

````text
allowBuilds:
  # Vitest/tsx depend on esbuild. This exact package is intentionally allowlisted.
  esbuild: true
minimumReleaseAgeExclude:
  - '@types/react-dom@19.2.5'
ignoredOptionalDependencies:
  # Next declares sharp optional; runtime image optimization is disabled.
  - sharp
````

### FILE: `README.md`

```yaml
block_id: "TS-GO-BRIDGE:readme:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified standalone BFF"
license: "LicenseRef-Workspace-Owner"
sha256: "5af210bc0c6dd811224b79783f1188544a64878637eedff6261585573e218fae"
variables: []
secrets_allowed: false
```

````text
# Enterprise Web BFF

Frontend y BFF opcional para el perfil empresarial Go/PostgreSQL. No contiene migrations, acceso SQL, persistencia, workers ni reglas transaccionales.

## Composición

Se materializa con `TS-OIDC-PORTAL-ADAPTER`; el backend debe exponer las APIs públicas y las queries protegidas declaradas por el perfil empresarial.

Configuración no secreta: `BUSINESS_CONFIG_FILE`, `APP_BASE_URL`, `OIDC_ISSUER`, `OIDC_CLIENT_ID`, `ENTERPRISE_API_BASE_URL`, `ENTERPRISE_TENANT_CODE` y `ENTERPRISE_ORGANIZATION_CODE`. `AUTH_SESSION_SECRET` y `OIDC_CLIENT_SECRET` se entregan mediante el mecanismo de secretos elegido, nunca en el repositorio.

## Verificación

```powershell
pnpm install --frozen-lockfile --offline
pnpm typecheck
pnpm test
pnpm build
pnpm licenses:report
```

La composición genera un nonce impredecible por request mediante `src/proxy.ts`, fuerza render dinámico y aplica una política CSP sin `unsafe-inline` ni `unsafe-eval` en producción. Google CSP Evaluator 1.1.8 verifica la política y Microsoft Playwright comprueba en cuatro navegadores que cada script lleva el nonce de su respuesta y que el siguiente request recibe otro valor.

Producción sigue condicionada al IdP, HTTPS/edge, comportamiento del CDN/WAF, protección CSRF donde corresponda, antiabuso distribuido, accesibilidad con navegador/AT, carga representativa, observabilidad y rollout/rollback del proyecto. El gate local no sustituye la repetición sobre el edge productivo.

## Imágenes

Esta referencia no transforma imágenes en runtime. `next.config.ts` establece `images.unoptimized: true`; `pnpm-workspace.yaml` excluye la dependencia opcional `sharp` y el lock congelado conserva esa selección. `/_next/image` responde404, comprobado en el recorrido de agenda con cuatro navegadores.

Antes de agregar optimización de imágenes, admitir la canalización elegida, versiones exactas, componentes nativos, licencias, avisos de seguridad y pruebas de rendimiento. Quitar estas dos opciones no equivale a admitir el grafo anterior. La omisión de Sharp no sustituye los gates del resto de dependencias ni la aceptación del proyecto.
````

### FILE: `THIRD_PARTY_NOTICES.md`

```yaml
block_id: "TS-GO-BRIDGE:third-party-notices:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified standalone BFF"
license: "LicenseRef-Workspace-Owner"
sha256: "4ed2b3256d843fbbdaa4dc5818d588231d88534288eece94693069079ad703ab"
variables: []
secrets_allowed: false
```

````text
# Dependencias y licencias

Las dependencias directas están fijadas en `package.json` y el grafo completo en `pnpm-lock.yaml`.

| Dependencia directa | Versión | Licencia declarada por el paquete |
|---|---:|---|
| jose | 6.2.10 | MIT |
| Next.js | 16.3.2 | MIT |
| openid-client | 6.8.5 | MIT |
| React / React DOM | 19.2.8 | MIT |
| Google SafeValues | 1.2.0 | Apache-2.0 |
| Google CSP Evaluator | 1.1.8 | Apache-2.0 |
| server-only | 0.0.1 | MIT |
| Zod | 4.4.3 | MIT |
| TypeScript | 7.0.2 | Apache-2.0 |
| Vitest | 4.1.11 | MIT |

Google CSP Evaluator se usa sólo como dependencia de desarrollo y su propio README declara que no es un producto oficial de Google ni ofrece garantía. Regenerar el inventario efectivo con `pnpm licenses:report`. El grafo transitivo observado también contiene Apache-2.0, MIT, ISC, BSD-3-Clause, 0BSD, CC-BY-4.0, MPL-2.0 y componentes binarios de imagen con términos adicionales. Este resumen no reemplaza conservar notices, revisar artefactos realmente distribuidos ni una revisión legal para el modo de entrega elegido.

No se asignó una licencia de distribución al código propio del adapter: esa decisión pertenece al propietario del repositorio. El pack se compone con el overlay OIDC y el adaptador de licencia para producir el inventario legal efectivo.

## FX direct-base adaptation

`internal/bcfx/exchange.go` and `selection.go` adapt Microsoft BCApps commit `2eae56d704a1fd035d104f333602aea7091b7749` under MIT. The exact copyright/license is `licenses/Microsoft-BCApps-MIT.txt`; derivation and source hashes are in `docs/provenance/BC_FX_DERIVATION.md` and `BC_FX_SOURCE_LOCK.json`. `rounding.go` and snapshot/accounting/API/host glue are AUTHORED; no AL built-in equivalence or upstream runtime execution is claimed.
````

### FILE: `tsconfig.json`

```yaml
block_id: "TS-GO-BRIDGE:tsconfig:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified standalone BFF"
license: "LicenseRef-Workspace-Owner"
sha256: "4c67dbbd61d014a203345b85086dc93ce3046260f229c733eafeeffe34808c0e"
variables: []
secrets_allowed: false
```

````text
{
  "compilerOptions": {
    "target": "ES2023",
    "lib": [
      "dom",
      "dom.iterable",
      "es2023"
    ],
    "allowJs": false,
    "skipLibCheck": true,
    "strict": true,
    "noUncheckedIndexedAccess": true,
    "exactOptionalPropertyTypes": true,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "bundler",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "react-jsx",
    "incremental": true,
    "plugins": [
      {
        "name": "next"
      }
    ],
    "paths": {
      "@/*": [
        "./src/*"
      ]
    }
  },
  "include": [
    "next-env.d.ts",
    "**/*.ts",
    "**/*.tsx",
    ".next/types/**/*.ts",
    ".next/dev/types/**/*.ts"
  ],
  "exclude": [
    "node_modules"
  ]
}
````

### FILE: `vitest.config.ts`

```yaml
block_id: "TS-GO-BRIDGE:vitest-config:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified standalone BFF"
license: "LicenseRef-Workspace-Owner"
sha256: "3a7bad10216bb4f319b026ca349a9f343bcc5e1f3b8a4b09892a4c9d91ef6a28"
variables: []
secrets_allowed: false
```

````text
import { defineConfig } from "vitest/config";
import { fileURLToPath } from "node:url";

export default defineConfig({
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
      "server-only": fileURLToPath(new URL("./src/test/server-only.ts", import.meta.url))
    }
  },
  test: {
    environment: "node",
    include: ["src/**/*.test.ts"],
    coverage: { reporter: ["text", "json", "html"] }
  }
});
````

### FILE: `src/app/admin/page.tsx`

```yaml
block_id: "TS-GO-BRIDGE:admin-page:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified standalone BFF"
license: "LicenseRef-Workspace-Owner"
sha256: "bdfd769fe1ce515beb10acb3a3f4e36efff13e9c418a10560914d1ad0544354b"
variables: []
secrets_allowed: false
```

````text
import type { ReactElement } from "react";
import { PortalPaging, portalCursor, portalPageHref, type PortalPageProps } from "@/platform/backend/portal-paging";
import { minorAmountPresentation } from "@/platform/i18n/money";
import { redirect } from "next/navigation";
import type { Route } from "next";
import { allowed, readSession } from "@/platform/auth/session";
import { protectedGet, type LeadSummary, type OrderSummary, type Overview, type Page, type ServiceCaseSummary } from "@/platform/backend/protected-client";

// AUTHORED framework signature glue: Next inspects the required props overload;
// the zero-argument overload preserves existing server-component test callers.
export default function AdminPage(): Promise<ReactElement>;
export default function AdminPage(props: PortalPageProps): Promise<ReactElement>;
export default async function AdminPage({ searchParams }: PortalPageProps = {}) {
  const session = await readSession();
  if (!session) redirect("/api/auth/login?return_to=/admin" as Route);
  if (!allowed(session, "admin:read")) return <><h1 className="pageTitle">Acceso denegado</h1><div className="notice">La sesión no posee <code>admin:read</code>.</div></>;
  const organization = session.organizations[0]!;
  const query = await searchParams ?? {};
  const ordersAfter=portalCursor(query.orders_after), casesAfter=portalCursor(query.cases_after);
  const leadsAfter=allowed(session,"lead:read") ? portalCursor(query.leads_after) : undefined;
  const [overview, orders, cases, leads] = await Promise.all([
    protectedGet<Overview>(session, "/v1/admin/overview", { organization_id: organization }),
    protectedGet<Page<OrderSummary>>(session, "/v1/admin/orders", { organization_id: organization, limit: "25", after: ordersAfter }),
    protectedGet<Page<ServiceCaseSummary>>(session, "/v1/admin/service-cases", { organization_id: organization, limit: "25", after: casesAfter }),
    allowed(session, "lead:read") ? protectedGet<Page<LeadSummary>>(session, "/v1/franchise/leads", { organization_id: organization, limit: "25", after: leadsAfter }) : Promise.resolve(null),
  ]);
  return <>
    <div className="eyebrow">Portal operativo · {organization}</div><h1 className="pageTitle">Operación</h1>
    <div className="grid">
      <article className="card"><h2>Pedidos</h2><p>{overview.orders}</p></article>
      <article className="card"><h2>Stock disponible</h2><p>{overview.stock_available}</p></article>
      <article className="card"><h2>Casos abiertos</h2><p>{overview.open_cases}</p></article>
      <article className="card"><h2>Envíos activos</h2><p>{overview.active_shipments}</p></article>
    </div>
    <section aria-label="Pedidos recientes"><h2>Pedidos recientes</h2><div className="grid">{orders.items.map((order) => <article className="card" key={order.id}><h3>{order.id}</h3><p>{order.state} · {minorAmountPresentation(order.total_minor_units, order.currency).amountLabel}</p></article>)}</div>
      {orders.items.length===0 ? <p>No hay registros en esta página.</p> : null}
      <PortalPaging label="Páginas de pedidos recientes" nextHref={orders.next_cursor ? portalPageHref("/admin",{orders_after:orders.next_cursor,cases_after:casesAfter,leads_after:leadsAfter}) : undefined} firstHref={ordersAfter ? portalPageHref("/admin",{cases_after:casesAfter,leads_after:leadsAfter}) : undefined}/></section>
    <section aria-label="Servicio"><h2>Servicio</h2><div className="grid">{cases.items.map((item) => <article className="card" key={item.id}><h3>{item.id}</h3><p>{item.state} · {item.severity}</p></article>)}</div>
      {cases.items.length===0 ? <p>No hay registros en esta página.</p> : null}
      <PortalPaging label="Páginas de servicio" nextHref={cases.next_cursor ? portalPageHref("/admin",{cases_after:cases.next_cursor,orders_after:ordersAfter,leads_after:leadsAfter}) : undefined} firstHref={casesAfter ? portalPageHref("/admin",{orders_after:ordersAfter,leads_after:leadsAfter}) : undefined}/></section>
    {leads && <section aria-label="Oportunidades"><h2>Oportunidades</h2><div className="grid">{leads.items.map((item) => <article className="card" key={item.id}><h3>{item.id}</h3><p>{item.state} · {item.source_code} · {item.assigned_subject || "sin asignar"}</p></article>)}</div>
      {leads.items.length===0 ? <p>No hay registros en esta página.</p> : null}
      <PortalPaging label="Páginas de oportunidades" nextHref={leads.next_cursor ? portalPageHref("/admin",{leads_after:leads.next_cursor,orders_after:ordersAfter,cases_after:casesAfter}) : undefined} firstHref={leadsAfter ? portalPageHref("/admin",{orders_after:ordersAfter,cases_after:casesAfter}) : undefined}/></section>}
  </>;
}
````

### FILE: `src/app/customer/page.tsx`

```yaml
block_id: "TS-GO-BRIDGE:customer-page:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified standalone BFF"
license: "LicenseRef-Workspace-Owner"
sha256: "00e6a23cf3b7cca2588912e7ed7ad541312784ea05a57870a0b1dca1829ba8cb"
variables: []
secrets_allowed: false
```

````text
import type { ReactElement } from "react";
import Link from "next/link";
import { CustomerCheckoutActions } from "@/components/customer-checkout-actions";
import { loadPublicLocale } from "@/platform/i18n/load-public-locale";
import { publicAppointmentTime } from "@/platform/i18n/public-catalog";
import { PortalPaging, portalCursor, portalPageHref, type PortalPageProps } from "@/platform/backend/portal-paging";
import { minorAmountPresentation } from "@/platform/i18n/money";
import { redirect } from "next/navigation";
import type { Route } from "next";
import { allowed, readSession } from "@/platform/auth/session";
import { protectedGet, type CustomerJourney, type OrderSummary, type Page, type ServiceCaseSummary } from "@/platform/backend/protected-client";

// AUTHORED framework signature glue: Next inspects the required props overload;
// the zero-argument overload preserves existing server-component test callers.
export default function CustomerPage(): Promise<ReactElement>;
export default function CustomerPage(props: PortalPageProps): Promise<ReactElement>;
export default async function CustomerPage({ searchParams }: PortalPageProps = {}) {
  const session = await readSession();
  if (!session) redirect("/api/auth/login?return_to=/customer" as Route);
  if (!allowed(session, "customer:self")) return <><h1 className="pageTitle">Acceso denegado</h1><div className="notice">La sesión no posee <code>customer:self</code>.</div></>;
  const organization = session.organizations[0]!;
  const query = await searchParams ?? {};
  const ordersAfter = portalCursor(query.orders_after), casesAfter = portalCursor(query.cases_after);
  const [orders, cases, journey, locale] = await Promise.all([
    protectedGet<Page<OrderSummary>>(session, "/v1/customer/orders", { organization_id: organization, limit: "25", after: ordersAfter }),
    protectedGet<Page<ServiceCaseSummary>>(session, "/v1/customer/service-cases", { organization_id: organization, limit: "25", after: casesAfter }),
    protectedGet<CustomerJourney>(session, "/v1/customer/journey", { organization_id: organization }),
    loadPublicLocale(),
  ]);
  return <><div className="eyebrow">Portal de cliente</div><h1 className="pageTitle">Mi cuenta</h1>
    <section aria-label="Mis pedidos"><h2>Mis pedidos</h2><div className="grid">{orders.items.map((order) => <article className="card" key={order.id}><h3>{order.id}</h3><p>{order.state} · {minorAmountPresentation(order.total_minor_units, order.currency).amountLabel}</p>{["placed","confirmed","allocated"].includes(order.state)?<CustomerCheckoutActions orderId={order.id} organizationId={order.organization_id}/>:null}</article>)}</div>
      {orders.items.length === 0 ? <p>No hay registros en esta página.</p> : null}
<PortalPaging label="Páginas de mis pedidos" nextHref={orders.next_cursor ? portalPageHref("/customer", { orders_after: orders.next_cursor, cases_after: casesAfter }) : undefined} firstHref={ordersAfter ? portalPageHref("/customer", { cases_after: casesAfter }) : undefined} /></section>
    <section aria-label="Mis casos de servicio"><h2>Mis casos de servicio</h2><div className="grid">{cases.items.map((item) => <article className="card" key={item.id}><h3>{item.id}</h3><p>{item.state} · {item.description}</p></article>)}</div>
      {cases.items.length === 0 ? <p>No hay registros en esta página.</p> : null}
<PortalPaging label="Páginas de mis casos de servicio" nextHref={cases.next_cursor ? portalPageHref("/customer", { cases_after: cases.next_cursor, orders_after: ordersAfter }) : undefined} firstHref={casesAfter ? portalPageHref("/customer", { orders_after: ordersAfter }) : undefined} /></section>
    <section><h2>Mis turnos</h2><p><Link href={"/customer/appointments" as Route}>Gestionar mis turnos</Link></p><div className="grid">{journey.appointments.map((item) => <article className="card" key={item.id}><h3>{item.kind}</h3><p>{item.state} · <time dateTime={item.starts_at}>{publicAppointmentTime(locale, item.starts_at)} ({locale.timeZone})</time></p></article>)}</div></section>
    <section><h2>Mis cotizaciones</h2><div className="grid">{journey.quotes.map((item) => <article className="card" key={item.id}><h3>{item.id}</h3><p>{item.state} · {minorAmountPresentation(item.total_minor_units, item.currency).amountLabel}</p></article>)}</div></section>
    <section><h2>Mis entregas</h2><div className="grid">{journey.handovers.map((item) => <article className="card" key={item.id}><h3>{item.order_id}</h3><p>{item.state}</p></article>)}</div></section>
  </>;
}
````

### FILE: `src/app/factory/page.tsx`

```yaml
block_id: "TS-GO-BRIDGE:factory-page:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified standalone BFF"
license: "LicenseRef-Workspace-Owner"
sha256: "927f3b474d0de5e72d76e0ea329f9efbce2ccb0286523d4ac46036307b51c512"
variables: []
secrets_allowed: false
```

````text
import { createHash } from "node:crypto";
import { FactoryUnitActions } from "@/components/factory-unit-actions";
import { factoryUnitSchema } from "@/platform/factory/contracts";
import type { ReactElement } from "react";
import { PortalPaging, portalCursor, portalPageHref, type PortalPageProps } from "@/platform/backend/portal-paging";
import { redirect } from "next/navigation";
import type { Route } from "next";
import { allowed, readSession } from "@/platform/auth/session";
import { protectedGet, type FactoryUnitSummary, type Page } from "@/platform/backend/protected-client";

// AUTHORED framework signature glue: Next inspects the required props overload;
// the zero-argument overload preserves existing server-component test callers.
export default function FactoryPage(): Promise<ReactElement>;
export default function FactoryPage(props: PortalPageProps): Promise<ReactElement>;
export default async function FactoryPage({ searchParams }: PortalPageProps = {}) {
  const session = await readSession();
  if (!session) redirect("/api/auth/login?return_to=/factory" as Route);
  if (!allowed(session, "factory:read")) return <><h1 className="pageTitle">Acceso denegado</h1><div className="notice">La sesión no posee <code>factory:read</code>.</div></>;
  const organization = session.organizations[0]!;
  const storageScope=createHash("sha256").update(JSON.stringify([session.tenantId,session.subject,organization])).digest("hex");
  const query = await searchParams ?? {};
  const unitsAfter = portalCursor(query.units_after);
  const units = await protectedGet<Page<FactoryUnitSummary>>(session, "/v1/factory/units", { organization_id: organization, limit: "25", after: unitsAfter });
  return <><div className="eyebrow">Portal de fábrica · {organization}</div><h1 className="pageTitle">Abastecimiento y seguimiento</h1>
    <div className="grid">{units.items.map((unit) => <article className="card" key={unit.id}><h2>{unit.serial_number}</h2><p>Orden {unit.purchase_order_id}</p><FactoryUnitActions unit={factoryUnitSchema.parse(unit)} canWrite={allowed(session,"factory:write")} storageScope={storageScope}/></article>)}</div>
    {units.items.length === 0 ? <p>No hay registros en esta página.</p> : null}
    <PortalPaging label="Páginas de unidades" nextHref={units.next_cursor ? portalPageHref("/factory", { units_after: units.next_cursor }) : undefined} firstHref={unitsAfter ? "/factory" : undefined} />
  </>;
}
````

### FILE: `src/app/globals.css`

```yaml
block_id: "TS-GO-BRIDGE:global-styles:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified standalone BFF"
license: "LicenseRef-Workspace-Owner"
sha256: "b7964d15c39cd463b57f522409e07b4592a19fdcd900da05909eae2bb2c796aa"
variables: []
secrets_allowed: false
```

````text
:root {
  color-scheme: light;
  --ink: #12221c;
  --muted: #5d6f67;
  --surface: #f4f7f5;
  --panel: #ffffff;
  --line: #d9e2dd;
  --accent: #006b4f;
  --accent-strong: #004c39;
  --warning: #7a4d00;
  font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
}

* { box-sizing: border-box; }
body { margin: 0; color: var(--ink); background: var(--surface); }
a { color: inherit; }
.skipLink { position: fixed; z-index: 100; left: 1rem; top: 1rem; padding: .7rem 1rem; border-radius: 8px; background: var(--ink); color: white; transform: translateY(-180%); }
.skipLink:focus { transform: translateY(0); }
header { background: var(--panel); border-bottom: 1px solid var(--line); }
.shell { width: min(1120px, calc(100% - 2rem)); margin: 0 auto; }
.headerRow { min-height: 68px; display: flex; align-items: center; justify-content: space-between; gap: 1.5rem; }
.brand { font-weight: 800; text-decoration: none; letter-spacing: -0.03em; }
nav { display: flex; flex-wrap: wrap; gap: 1rem; }
nav a { color: var(--muted); text-decoration: none; font-weight: 650; }
nav a:hover, nav a:focus-visible { color: var(--accent); text-decoration: underline; }
main { padding: 3rem 0 5rem; }
.hero { display: grid; grid-template-columns: minmax(0, 1.5fr) minmax(280px, .8fr); gap: 2rem; align-items: center; }
.eyebrow { color: var(--accent); font-weight: 800; text-transform: uppercase; font-size: .78rem; letter-spacing: .12em; }
h1 { font-size: clamp(2.4rem, 7vw, 5.5rem); line-height: .95; letter-spacing: -.06em; margin: .6rem 0 1.2rem; }
.pageTitle { font-size: clamp(2.2rem, 5vw, 4rem); }
h2 { font-size: clamp(1.5rem, 3vw, 2.4rem); letter-spacing: -.035em; }
p { line-height: 1.65; }
.lede { font-size: 1.18rem; color: var(--muted); max-width: 68ch; }
.panel, .card { background: var(--panel); border: 1px solid var(--line); border-radius: 18px; padding: 1.4rem; box-shadow: 0 14px 35px rgba(18, 34, 28, .06); }
.grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); gap: 1rem; margin-top: 1.5rem; }
.metric { font-size: 2rem; font-weight: 850; color: var(--accent); }
.muted { color: var(--muted); }
.button { display: inline-flex; align-items: center; justify-content: center; padding: .8rem 1rem; border-radius: 10px; border: 0; background: var(--accent); color: white; font-weight: 750; text-decoration: none; cursor: pointer; }
.button:hover, .button:focus-visible { background: var(--accent-strong); }
.actions { display: flex; flex-wrap: wrap; gap: .75rem; margin-top: 1.5rem; }
.primaryAction { display: inline-flex; padding: .85rem 1.1rem; border-radius: 10px; background: var(--accent); color: white; font-weight: 750; text-decoration: none; }
.primaryAction:hover, .primaryAction:focus-visible { background: var(--accent-strong); }
label { display: grid; gap: .4rem; font-weight: 700; }
input, select { width: 100%; border: 1px solid #aebdb5; border-radius: 9px; padding: .75rem; font: inherit; background: white; }
input:focus, select:focus { outline: 3px solid rgba(0,107,79,.25); border-color: var(--accent); }
form { display: grid; gap: 1rem; }
.status { padding: .8rem; border-radius: 8px; background: #e7f6ef; color: var(--accent-strong); }
.notice { margin: 1rem 0 2rem; padding: 1rem; border: 1px solid #e3c98e; border-radius: 10px; background: #fff8e8; color: var(--warning); line-height: 1.6; }
.error { background: #fff0ed; color: #8a2717; }
footer { padding: 2rem 0; border-top: 1px solid var(--line); color: var(--muted); background: var(--panel); }
code { background: #e7ece9; border-radius: 5px; padding: .12rem .3rem; }
.notificationStatus { border-top: 1px solid var(--line); margin-top: 1rem; padding-top: 1rem; overflow-wrap: anywhere; }
.notificationStatus dd { margin-inline-start: 0; }
@media (max-width: 760px) { .hero { grid-template-columns: 1fr; } .headerRow { align-items: flex-start; flex-direction: column; padding: 1rem 0; } }
@media (prefers-reduced-motion: reduce) { *, *::before, *::after { scroll-behavior: auto !important; transition: none !important; } }
````

### FILE: `src/app/icon.svg`

```yaml
block_id: "TS-GO-BRIDGE:app-icon:v1"
operation: CREATE
provenance: AUTHORED
source: "local accessible application icon"
license: "LicenseRef-Workspace-Owner"
sha256: "9c6fbf50693d37bf84e24f0385eee69818bf8f24a2012e6cce09f4975958bba4"
variables: []
secrets_allowed: false
```

````xml
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64">
  <rect width="64" height="64" rx="14" fill="#0b3d2e"/>
  <path d="M35 8 17 35h13l-2 21 19-30H34z" fill="#86efac"/>
</svg>
````

### FILE: `src/app/layout.tsx`

```yaml
block_id: "TS-GO-BRIDGE:root-layout:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified standalone BFF"
license: "LicenseRef-Workspace-Owner"
sha256: "a59ea5a809717e84b780f649ede51693f68008fbf977eb32d2a03ea5c401491e"
variables: []
secrets_allowed: false
```

````text
import { publicMessage as message } from "@/platform/i18n/public-catalog";
import { loadPublicLocale } from "@/platform/i18n/load-public-locale";
import type { Metadata, Route } from "next";
import Link from "next/link";
import { loadBusinessConfig } from "@/platform/config/load";
import { readSession } from "@/platform/auth/session";
import { enabledNavigation } from "@/platform/config/registry";
import { connection } from "next/server";
import "./globals.css";

export async function generateMetadata(): Promise<Metadata> {
  const config = await loadBusinessConfig();
  const locale = await loadPublicLocale();
  return {
    title: { default: config.business.name, template: `%s | ${config.business.name}` },
    description: message(locale.locale, "site.description"),
    robots: { index: false, follow: false }
  };
}

export default async function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  await connection();
  const config = await loadBusinessConfig();
  const locale = await loadPublicLocale();
  const session = await readSession();
  const navigation = enabledNavigation(config, session?.permissions ?? null);
  return (
    <html lang={locale.locale}>
      <body>
        <a className="skipLink" href="#main">{message(locale.locale, "site.skip")}</a>
        <header>
          <div className="shell headerRow">
            <Link className="brand" href="/">{config.business.name}</Link>
            <nav aria-label={message(locale.locale, "site.navigation")}>
              {navigation.map((item) => <Link key={item.id} href={item.href as Route}>{message(locale.locale, "nav." + item.id)}</Link>)}
            </nav>
          </div>
        </header>
        <main id="main" className="shell" lang="es">{children}</main>
        <footer><div className="shell">{message(locale.locale, "site.contact")} <a href={`mailto:${config.business.supportEmail}`}>{config.business.supportEmail}</a></div></footer>
      </body>
    </html>
  );
}
````

### FILE: `src/app/models/page.tsx`

```yaml
block_id: "TS-GO-BRIDGE:models-page:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified standalone BFF"
license: "LicenseRef-Workspace-Owner"
sha256: "4d6d2f5d015c04863da4dcd2a8f1650f86aae410a7299b31d0a16c75a4ee996d"
variables: []
secrets_allowed: false
```

````text
import { publicMessage as message } from "@/platform/i18n/public-catalog";
import { loadPublicLocale } from "@/platform/i18n/load-public-locale";
import { publicCount } from "@/platform/i18n/public-catalog";
import { listModels } from "@/platform/backend/public-client";
import { LeadForm } from "@/app/connected/lead-form";
import { publicPageMetadata } from "@/platform/seo/public-indexing";

export async function generateMetadata() { const locale = await loadPublicLocale(); return publicPageMetadata("/models", message(locale.locale, "models.short")); }
export const dynamic = "force-dynamic";

export default async function ModelsPage() {
  const [models, locale] = await Promise.all([listModels(), loadPublicLocale()]);
  return (
    <section lang={locale.locale}>
      <div className="eyebrow">{message(locale.locale, "models.eyebrow")}</div>
      <h1 className="pageTitle">{message(locale.locale, "models.title")}</h1>
      <p className="lede">{message(locale.locale, "models.description")}</p>
      <p>{publicCount(locale.locale, "models", models.length)}</p>
      <div className="grid">
        {models.map((model) => (
          <article className="card" key={model.id}>
            <h2>{model.displayName}</h2>
            <p>{model.vehicleClass}</p>
            <LeadForm modelId={model.id} locale={locale.locale} />
          </article>
        ))}
      </div>
    </section>
  );
}
````

### FILE: `src/app/page.tsx`

```yaml
block_id: "TS-GO-BRIDGE:home-page:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified standalone BFF"
license: "LicenseRef-Workspace-Owner"
sha256: "29442ab784362d7bdbe2f460c6c531fbd1568154d41bbab34d7ee91b59f36a26"
variables: []
secrets_allowed: false
```

````text
import { publicMessage as message } from "@/platform/i18n/public-catalog";
import { loadPublicLocale } from "@/platform/i18n/load-public-locale";
import Link from "next/link";
import { loadBusinessConfig } from "@/platform/config/load";
import { headers } from "next/headers";
import { publicPageMetadata, websiteJsonLd } from "@/platform/seo/public-indexing";

export function generateMetadata() { return publicPageMetadata("/"); }

export default async function HomePage() {
  const config = await loadBusinessConfig();
  const locale = await loadPublicLocale();
  const structuredData = websiteJsonLd(config.business.name, locale.locale);
  const nonce = (await headers()).get("x-nonce") ?? undefined;
  return (
    <section lang={locale.locale}>
      {structuredData && <script type="application/ld+json" nonce={nonce} dangerouslySetInnerHTML={{ __html: structuredData }} />}
      <div className="eyebrow">{message(locale.locale, "home.eyebrow")}</div>
      <h1 className="pageTitle">{config.business.name}</h1>
      <p className="lede">{message(locale.locale, "home.description")}</p>
      <div className="actions"><Link className="primaryAction" href="/models">{message(locale.locale, "home.action")}</Link></div>
    </section>
  );
}
````

### FILE: `src/platform/config/load.ts`

```yaml
block_id: "TS-GO-BRIDGE:config-load:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified standalone BFF"
license: "LicenseRef-Workspace-Owner"
sha256: "b61d038dc12ccb77e9002a0fedc8c6310c3cdd1e7b6b5ba03c660b5258642920"
variables: []
secrets_allowed: false
```

````text
import { readFile } from "node:fs/promises";
import { basename, join } from "node:path";
import { businessConfigSchema, type BusinessConfig } from "./schema";

let cached: BusinessConfig | undefined;

export async function loadBusinessConfig(options?: { path?: string; bypassCache?: boolean }): Promise<BusinessConfig> {
  if (cached && !options?.bypassCache) return cached;

  const configFile = basename(options?.path ?? process.env.BUSINESS_CONFIG_FILE ?? "business.example.json");
  const configPath = join(process.cwd(), "config", configFile);
  const raw = JSON.parse(await readFile(configPath, "utf8")) as unknown;
  const parsed = businessConfigSchema.parse(raw);
  if (!options?.bypassCache) cached = parsed;
  return parsed;
}

export function clearBusinessConfigCache(): void {
  cached = undefined;
}
````

### FILE: `src/platform/config/registry.ts`

```yaml
block_id: "TS-GO-BRIDGE:config-registry:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified standalone BFF"
license: "LicenseRef-Workspace-Owner"
sha256: "a13463dc975803a7789c584b96df248a0d5305fbf8553055c69fb2681538a907"
variables: []
secrets_allowed: false
```

````text
import type { BusinessConfig } from "./schema";

export interface NavigationItem {
  id: string;
  label: string;
  href: string;
  module?: string;
  feature?: string;
  permission?: string;
  anyPermissions?: readonly string[];
  guestEntry?: boolean;
}

const navigation: NavigationItem[] = [
  { id: "dashboard", label: "Mi panel", href: "/dashboard", feature: "role_workspace" },
  { id: "help", label: "Ayuda", href: "/help" },
  { id: "public_catalog", label: "Modelos", href: "/models", module: "catalog", feature: "public_catalog" },
  { id: "locations", label: "Dónde estamos", href: "/locations", module: "crm" },
  { id: "customer", label: "Mi cuenta", href: "/customer", feature: "customer_portal", permission: "customer:self", guestEntry: true },
  { id: "admin", label: "Operación", href: "/admin", permission: "admin:read" },
  { id: "franchise", label: "Franquicia", href: "/franchise", module: "crm", anyPermissions: ["inventory:allocate", "payment:create", "handover:manage", "admin:read", "lead:read", "resource:manage", "availability:read", "availability:manage", "appointment:manage"] },
  { id: "factory", label: "Fábrica", href: "/factory", module: "procurement", feature: "factory_portal", permission: "factory:read" }
];

export function enabledNavigation(config: BusinessConfig, permissions: readonly string[] | null = null): NavigationItem[] {
  return navigation.filter((item) => {
    if (item.module && !config.modules[item.module]?.enabled) return false;
    if (item.feature && !config.features[item.feature]) return false;
    const can = (permission: string) => permissions !== null && (permissions.includes("*") || permissions.includes(permission));
    if (item.permission && !(permissions === null && item.guestEntry) && !can(item.permission)) return false;
    if (item.anyPermissions && !item.anyPermissions.some(can)) return false;
    return true;
  });
}

export function enabledModules(config: BusinessConfig): string[] {
  return Object.entries(config.modules).filter(([, value]) => value.enabled).map(([name]) => name);
}
````

### FILE: `src/platform/config/schema.test.ts`

```yaml
block_id: "TS-GO-BRIDGE:config-schema-test:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified standalone BFF"
license: "LicenseRef-Workspace-Owner"
sha256: "66e81ed47b5c7694ada84bc355d30cc3a6347106d8c6d9cc0c09cc12d10bf425"
variables: []
secrets_allowed: false
```

````text
import { afterEach, describe, expect, it } from "vitest";
import { clearBusinessConfigCache, loadBusinessConfig } from "./load";
import { businessConfigSchema } from "./schema";

afterEach(() => clearBusinessConfigCache());

describe("business presentation configuration", () => {
  it("loads the bounded config file from ./config", async () => {
    const config = await loadBusinessConfig({ bypassCache: true });
    expect(config.business.id).toBe("electric-mobility-network");
    expect(config.workflows.lead?.states).toEqual(["new", "contacted", "qualified", "converted", "lost"]);
    expect(config.roles.find((role) => role.id === "customer")?.permissions).toEqual(["customer:self"]);
  });

  it("keeps the example lead workflow inside the PostgreSQL state contract", async () => {
    const config = await loadBusinessConfig({ bypassCache: true });
    const invalid = structuredClone(config);
    invalid.workflows.lead!.states.push("assigned");
    expect(invalid.workflows.lead!.states).not.toEqual(config.workflows.lead!.states);
    expect(config.workflows.lead!.states).not.toContain("assigned");
  });

  it("rejects an invalid module dependency", async () => {
    const config = await loadBusinessConfig({ bypassCache: true });
    const invalid = structuredClone(config);
    invalid.modules.orders = { enabled: false };
    expect(() => businessConfigSchema.parse(invalid)).toThrow(/requires enabled module orders/);
  });

  it("does not allow config paths to escape ./config", async () => {
    await expect(loadBusinessConfig({ path: "../outside.json", bypassCache: true })).rejects.toThrow();
  });
});
````

### FILE: `src/platform/config/schema.ts`

```yaml
block_id: "TS-GO-BRIDGE:config-schema:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified standalone BFF"
license: "LicenseRef-Workspace-Owner"
sha256: "c8ad5a8d05122ca180fdb1be4098d04096f78b3039a672c241df703223584892"
variables: []
secrets_allowed: false
```

````text
import { z } from "zod";

const identifier = z.string().regex(/^[a-z][a-z0-9_]*$/, "use snake_case identifiers");
const permission = z.string().regex(/^\*$|^[a-z][a-z0-9_]*(?::[a-z][a-z0-9_]*){1,2}$/);

const fieldSchema = z.discriminatedUnion("type", [
  z.object({ id: identifier, label: z.string().min(1), type: z.literal("string"), required: z.boolean() }),
  z.object({ id: identifier, label: z.string().min(1), type: z.literal("number"), required: z.boolean() }),
  z.object({ id: identifier, label: z.string().min(1), type: z.literal("boolean"), required: z.boolean() }),
  z.object({ id: identifier, label: z.string().min(1), type: z.literal("date"), required: z.boolean() }),
  z.object({ id: identifier, label: z.string().min(1), type: z.literal("select"), required: z.boolean(), options: z.array(identifier).min(1) })
]);

const workflowSchema = z.object({
  initial: identifier,
  states: z.array(identifier).min(1),
  transitions: z.array(z.object({ from: identifier, to: identifier, permission }))
});

export const businessConfigSchema = z.object({
  schemaVersion: z.literal("1.0.0"),
  business: z.object({
    id: z.string().regex(/^[a-z][a-z0-9-]*$/),
    name: z.string().min(1),
    defaultLocale: z.string().min(2),
    defaultMarket: z.string().length(2),
    supportEmail: z.email()
  }),
  markets: z.array(z.object({
    code: z.string().length(2),
    name: z.string().min(1),
    currency: z.string().length(3),
    locales: z.array(z.string().min(2)).min(1),
    timeZone: z.string().min(1),
    taxMode: z.enum(["internal", "external"])
  })).min(1),
  organizationTypes: z.array(z.object({
    id: identifier,
    label: z.string().min(1),
    allowedParents: z.array(identifier)
  })).min(1),
  roles: z.array(z.object({
    id: identifier,
    label: z.string().min(1),
    permissions: z.array(permission).min(1)
  })).min(1),
  modules: z.record(identifier, z.object({ enabled: z.boolean() })),
  workflows: z.record(identifier, workflowSchema),
  customFields: z.record(identifier, z.array(fieldSchema)),
  integrations: z.array(z.object({
    id: identifier,
    provider: identifier,
    enabled: z.boolean(),
    mode: z.enum(["sandbox", "production"]),
    capabilities: z.array(identifier).min(1),
    credentialRefEnv: z.string().regex(/^[A-Z][A-Z0-9_]+$/)
  })),
  features: z.record(identifier, z.boolean())
}).superRefine((config, context) => {
  const unique = (values: string[], path: (string | number)[], label: string) => {
    if (new Set(values).size !== values.length) {
      context.addIssue({ code: "custom", message: `${label} must be unique`, path });
    }
  };

  unique(config.markets.map((item) => item.code), ["markets"], "market codes");
  unique(config.organizationTypes.map((item) => item.id), ["organizationTypes"], "organization type ids");
  unique(config.roles.map((item) => item.id), ["roles"], "role ids");
  unique(config.integrations.map((item) => item.id), ["integrations"], "integration ids");

  const market = config.markets.find((item) => item.code === config.business.defaultMarket);
  if (!market) {
    context.addIssue({ code: "custom", message: "defaultMarket must reference an existing market", path: ["business", "defaultMarket"] });
  } else if (!market.locales.includes(config.business.defaultLocale)) {
    context.addIssue({ code: "custom", message: "defaultLocale must belong to defaultMarket", path: ["business", "defaultLocale"] });
  }

  const organizationTypeIds = new Set(config.organizationTypes.map((item) => item.id));
  config.organizationTypes.forEach((type, index) => {
    type.allowedParents.forEach((parent) => {
      if (!organizationTypeIds.has(parent)) {
        context.addIssue({ code: "custom", message: `unknown parent organization type: ${parent}`, path: ["organizationTypes", index, "allowedParents"] });
      }
      if (parent === type.id) {
        context.addIssue({ code: "custom", message: "organization type cannot parent itself", path: ["organizationTypes", index, "allowedParents"] });
      }
    });
  });

  for (const [name, workflow] of Object.entries(config.workflows)) {
    unique(workflow.states, ["workflows", name, "states"], `${name} states`);
    const states = new Set(workflow.states);
    if (!states.has(workflow.initial)) {
      context.addIssue({ code: "custom", message: "initial state must be declared", path: ["workflows", name, "initial"] });
    }
    workflow.transitions.forEach((transition, index) => {
      if (!states.has(transition.from) || !states.has(transition.to)) {
        context.addIssue({ code: "custom", message: "transition references an unknown state", path: ["workflows", name, "transitions", index] });
      }
      if (transition.from === transition.to) {
        context.addIssue({ code: "custom", message: "self transitions are not allowed", path: ["workflows", name, "transitions", index] });
      }
    });
  }

  for (const [entity, fields] of Object.entries(config.customFields)) {
    unique(fields.map((field) => field.id), ["customFields", entity], `${entity} custom field ids`);
  }

  const moduleDependencies: Record<string, string[]> = {
    payments: ["orders"],
    fulfillment: ["orders", "inventory"],
    service: ["catalog"],
    procurement: ["inventory"],
    integrations: ["catalog"]
  };
  for (const [moduleName, dependencies] of Object.entries(moduleDependencies)) {
    if (config.modules[moduleName]?.enabled) {
      dependencies.forEach((dependency) => {
        if (!config.modules[dependency]?.enabled) {
          context.addIssue({ code: "custom", message: `${moduleName} requires enabled module ${dependency}`, path: ["modules", moduleName] });
        }
      });
    }
  }
});

export type BusinessConfig = z.infer<typeof businessConfigSchema>;
export type WorkflowConfig = z.infer<typeof workflowSchema>;
export type CustomField = z.infer<typeof fieldSchema>;
````

### FILE: `src/platform/http/problem.test.ts`

```yaml
block_id: "TS-GO-BRIDGE:http-problem-test:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified standalone BFF"
license: "LicenseRef-Workspace-Owner"
sha256: "9df9682db5bf731686285a64869487b416f1bcc630fc170bbfd006e358d67247"
variables: []
secrets_allowed: false
```

````text
import { describe, expect, it } from "vitest";
import { z } from "zod";
import { BackendProblem } from "@/platform/backend/public-client";
import { errorResponse } from "./problem";

describe("BFF problem responses", () => {
  it("does not expose validation input", async () => {
    let failure: unknown;
    try { z.object({ email: z.email() }).parse({ email: "secret-invalid" }); } catch (error) { failure = error; }
    const response = errorResponse(failure);
    expect(response.status).toBe(400);
    expect(await response.text()).not.toContain("secret-invalid");
  });

  it("maps upstream server failures to a bounded gateway response", async () => {
    const response = errorResponse(new BackendProblem(503, "UPSTREAM_OVERLOAD"));
    expect(response.status).toBe(502);
    expect(await response.json()).toMatchObject({ code: "UPSTREAM_OVERLOAD" });
  });
});
````

### FILE: `src/platform/http/problem.ts`

```yaml
block_id: "TS-GO-BRIDGE:http-problem:v2"
operation: CREATE
provenance: AUTHORED
source: "local verified standalone BFF"
license: "LicenseRef-Workspace-Owner"
sha256: "dc748632c8465c9e4f7bab4cb9249fc76bfb1cb76e04ac7ad5a48dd60c118b55"
variables: []
secrets_allowed: false
```

````text
import { NextResponse } from "next/server";
import { ZodError } from "zod";
import { BackendProblem } from "@/platform/backend/public-client";

export function problem(status: number, title: string, detail: string, code: string): NextResponse {
  return NextResponse.json({ type: "about:blank", title, status, detail, code }, {
    status,
    headers: { "cache-control": "no-store", "content-type": "application/problem+json" }
  });
}

export function errorResponse(error: unknown): NextResponse {
  if (error instanceof ZodError) return problem(400, "Solicitud inválida", "Revise los datos enviados.", "VALIDATION_FAILED");
  if (error instanceof BackendProblem) {
    const status = error.status >= 400 && error.status < 500 ? error.status : 502;
    return problem(status, "Backend no disponible", "No pudimos completar la operación.", error.code);
  }
  return problem(500, "Error interno", "La operación no pudo completarse.", "INTERNAL_ERROR");
}
````

### FILE: `src/platform/security/google-safevalues.test.ts`

```yaml
block_id: "TS-GO-BRIDGE:google-safevalues-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local adapter test exercising safevalues 1.2.0 from Google"
license: "LicenseRef-Workspace-Owner"
sha256: "d0078fda42a116582bd863efd16a3269bb130f5c1c190b32e689d3cd4fe4dbb3"
variables: []
secrets_allowed: false
```

````typescript
import { describe, expect, it, vi } from "vitest";
import {
  htmlEscape,
  setAnchorHref,
  unwrapHtml,
} from "./google-safevalues";

describe("Google SafeValues DOM-XSS primitives", () => {
  it("escapes attacker-controlled HTML instead of creating markup", () => {
    const escaped = unwrapHtml(htmlEscape('<img src=x onerror="alert(1)">'));
    expect(String(escaped)).toBe("&lt;img src=x onerror=&quot;alert(1)&quot;&gt;");
  });

  it("refuses a javascript URL without overwriting the current anchor", () => {
    const error = vi.spyOn(console, "error").mockImplementation(() => undefined);
    const anchor = { href: "https://example.invalid/safe" } as HTMLAnchorElement;
    setAnchorHref(anchor, "javascript:alert(document.domain)");
    expect(anchor.href).toBe("https://example.invalid/safe");
    expect(error).toHaveBeenCalledOnce();
    error.mockRestore();
  });

  it("preserves an ordinary HTTPS URL", () => {
    const anchor = { href: "" } as HTMLAnchorElement;
    setAnchorHref(anchor, "https://example.invalid/catalog?model=e-bike");
    expect(anchor.href).toBe("https://example.invalid/catalog?model=e-bike");
  });
});
````

### FILE: `src/platform/security/google-safevalues.ts`

```yaml
block_id: "TS-GO-BRIDGE:google-safevalues-adapter:v1"
operation: CREATE
provenance: AUTHORED
source: "thin exports over Google safevalues 1.2.0"
license: "LicenseRef-Workspace-Owner"
sha256: "56cdffcd67130be327d0b842fd88e115784c890ef566d4a3e272410f4168b9b3"
variables: []
secrets_allowed: false
```

````typescript
/**
 * Google SafeValues is the admitted DOM-XSS primitive for non-React DOM sinks.
 * React text interpolation remains preferred; raw HTML still requires an explicit
 * product policy, provenance and review before calling a sanitizer.
 */
export {
  htmlEscape,
  isHtml,
  sanitizeHtml,
  sanitizeHtmlAssertUnchanged,
  unwrapHtml,
} from "safevalues";

export {
  setAnchorHref,
  setElementInnerHtml,
  setIframeSrcdoc,
  setLocationHref,
  setScriptSrc,
} from "safevalues/dom";
````

### FILE: `src/platform/security/csp-policy.ts`

```yaml
block_id: "TS-GO-BRIDGE:csp-policy:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed policy builder evaluated by Google CSP Evaluator 1.1.8"
license: "LicenseRef-Workspace-Owner"
sha256: "48f48d7293133ae9095c544eac29ce9b9b8f53e0c8a5101e5af433c45fa4492a"
variables: []
secrets_allowed: false
```

````typescript
const noncePattern = /^[A-Za-z0-9+/]+={0,2}$/;

export function buildContentSecurityPolicy(nonce: string, development: boolean): string {
  if (nonce.length < 16 || !noncePattern.test(nonce)) throw new Error("CSP nonce must be base64 and at least 16 characters");
  return `
    default-src 'self';
    script-src 'self' 'nonce-${nonce}' 'strict-dynamic'${development ? " 'unsafe-eval'" : ""};
    style-src 'self' 'nonce-${nonce}';
    img-src 'self' blob: data:;
    font-src 'self';
    connect-src 'self';
    object-src 'none';
    base-uri 'self';
    form-action 'self';
    frame-ancestors 'none';
    upgrade-insecure-requests;
  `.replace(/\s{2,}/g, " ").trim();
}
````

### FILE: `src/platform/security/csp-policy.test.ts`

```yaml
block_id: "TS-GO-BRIDGE:csp-policy-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local positive/negative regression over the exact Google evaluator package"
license: "LicenseRef-Workspace-Owner"
sha256: "60ddab62cc214a7bcb9ebc419f1c824f71ce372abccf9bf3be5aeb9ad32bee5f"
variables: []
secrets_allowed: false
```

````typescript
import { describe, expect, it } from "vitest";
import { CspEvaluator } from "csp_evaluator/dist/evaluator.js";
import { CspParser } from "csp_evaluator/dist/parser.js";
import { Severity } from "csp_evaluator/dist/finding.js";
import { buildContentSecurityPolicy } from "./csp-policy";

const nonce = "dGVzdC1ub25jZS0xMjM0NTY=";

describe("strict Content Security Policy", () => {
  it("has no actionable finding in Google's CSP Evaluator", () => {
    const policy = buildContentSecurityPolicy(nonce, false);
    const findings = new CspEvaluator(new CspParser(policy).csp).evaluate();
    expect(findings.filter((finding) => finding.severity < Severity.NONE)).toEqual([]);
    expect(policy).not.toContain("'unsafe-inline'");
    expect(policy).not.toContain("'unsafe-eval'");
  });

  it("keeps unsafe-eval restricted to development", () => {
    expect(buildContentSecurityPolicy(nonce, true)).toContain("'unsafe-eval'");
    expect(buildContentSecurityPolicy(nonce, false)).not.toContain("'unsafe-eval'");
  });

  it("fails closed for malformed or short nonces", () => {
    expect(() => buildContentSecurityPolicy("short", false)).toThrow(/nonce/);
    expect(() => buildContentSecurityPolicy("not-a-valid-base64-token!", false)).toThrow(/nonce/);
  });

  it("proves the previous header was not acceptable", () => {
    const previous = "base-uri 'self'; object-src 'none'; frame-ancestors 'none'; form-action 'self'; upgrade-insecure-requests";
    const findings = new CspEvaluator(new CspParser(previous).csp).evaluate();
    expect(findings.some((finding) => finding.severity === Severity.HIGH && finding.directive === "script-src")).toBe(true);
  });
});
````

### FILE: `src/proxy.ts`

```yaml
block_id: "TS-GO-BRIDGE:strict-csp-proxy:v1"
operation: CREATE
provenance: ADAPTED
source: "adapted from the official Next.js 16 Content Security Policy guide; local module boundary and policy builder"
license: "LicenseRef-Workspace-Owner"
sha256: "40392fc078888a9998b63e8026b280ed1e38c18ec4dd916d24c13ca5a2324fd0"
variables: []
secrets_allowed: false
```

````typescript
import { NextRequest, NextResponse } from "next/server";
import { buildContentSecurityPolicy } from "@/platform/security/csp-policy";

export function proxy(request: NextRequest) {
  const nonce = Buffer.from(crypto.randomUUID()).toString("base64");
  const csp = buildContentSecurityPolicy(nonce, process.env.NODE_ENV === "development");
  const requestHeaders = new Headers(request.headers);
  requestHeaders.set("x-nonce", nonce);
  requestHeaders.set("Content-Security-Policy", csp);
  const response = NextResponse.next({ request: { headers: requestHeaders } });
  response.headers.set("Content-Security-Policy", csp);
  return response;
}

export const config = {
  matcher: [{
    source: "/((?!api|_next/static|_next/image|favicon.ico).*)",
    missing: [
      { type: "header", key: "next-router-prefetch" },
      { type: "header", key: "purpose", value: "prefetch" }
    ]
  }]
};
````

### FILE: `src/platform/seo/public-indexing.ts`
```yaml
block_id: "TS-GO-WEB-SEO:public-indexing.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "AUTHORED integration of exact Next16.3.4 metadata APIs; PUBLIC_WEB_METADATA_INTEGRATION_V377.md records official authority and executed evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "503cec4f8c6d025a746b66b0fb16bb3c703d3f0b30a9357ff82c52b158694ea9"
variables: []
secrets_allowed: false
```
````typescript
import type { Metadata, MetadataRoute } from "next";

// Only these existing public pages are approved by this library reference.
// Protected portals and dynamic business objects are never inferred from routes.
const publicPaths = ["/", "/models"] as const;
export type PublicPage = typeof publicPaths[number];
export type PublicIndexing = Readonly<{ enabled: boolean; origin: string | null }>;

export function readPublicIndexing(input: {
  PUBLIC_SITE_ORIGIN?: string | undefined;
  PUBLIC_INDEXING_ENABLED?: string | undefined;
} = { PUBLIC_SITE_ORIGIN: process.env.PUBLIC_SITE_ORIGIN,
  PUBLIC_INDEXING_ENABLED: process.env.PUBLIC_INDEXING_ENABLED }): PublicIndexing {
  const flag = input.PUBLIC_INDEXING_ENABLED;
  if (flag !== undefined && flag !== "" && flag !== "0" && flag !== "1") {
    throw new Error("Invalid public indexing configuration");
  }
  const enabled = flag === "1";
  const raw = input.PUBLIC_SITE_ORIGIN;
  if (!raw) {
    if (enabled) throw new Error("Public indexing requires an explicit HTTPS origin");
    return { enabled: false, origin: null };
  }
  // This value is deployment configuration, never Host/X-Forwarded-Host input.
  // Refuse URL normalization that could hide credentials, controls or paths.
  if (raw.length > 2048 || /[\s\\%?#]/u.test(raw)) throw new Error("Invalid public site origin");
  let url: URL;
  try { url = new URL(raw); } catch { throw new Error("Invalid public site origin"); }
  if (url.protocol !== "https:" || !url.hostname || url.username || url.password ||
      url.pathname !== "/" || url.port || !/^https:\/\/[^/]+\/?$/u.test(raw)) {
    throw new Error("Invalid public site origin");
  }
  return { enabled, origin: url.origin };
}

export function publicPageMetadata(page: PublicPage, title?: string): Metadata {
  if (!publicPaths.includes(page)) throw new Error("Unknown public page");
  const settings = readPublicIndexing();
  return {
    ...(title ? { title } : {}),
    robots: { index: settings.enabled, follow: settings.enabled },
    ...(settings.origin ? { alternates: { canonical: new URL(page, settings.origin).href } } : {})
  };
}

export function publicRobots(settings = readPublicIndexing()): MetadataRoute.Robots {
  if (!settings.enabled || !settings.origin) return { rules: { userAgent: "*", disallow: "/" } };
  // Exact public-path matches: adding a portal does not expose it to crawlers.
  // This is a crawl instruction, not access control or removal from an index.
  return { rules: { userAgent: "*", disallow: "/", allow: [...publicPaths.map(path => path + "$"), "/_next/static/", "/icon.svg$"] },
    sitemap: new URL("/sitemap.xml", settings.origin).href };
}

export function publicSitemap(settings = readPublicIndexing()): MetadataRoute.Sitemap {
  if (!settings.enabled || !settings.origin) return [];
  // Do not invent lastModified dates, model detail URLs, prices or translations.
  return publicPaths.map(path => ({ url: new URL(path, settings.origin!).href }));
}

export function websiteJsonLd(name: string, language: string): string | null {
  const settings = readPublicIndexing();
  if (!settings.enabled || !settings.origin) return null;
  const data = { "@context": "https://schema.org", "@type": "WebSite", name,
    url: new URL("/", settings.origin).href, inLanguage: language };
  // A JSON string embedded in an HTML script must not contain a closing tag.
  return JSON.stringify(data).replace(/[<>&\u2028\u2029]/gu, character =>
    "\\u" + character.charCodeAt(0).toString(16).padStart(4, "0"));
}
````

### FILE: `src/platform/seo/public-indexing.test.ts`
```yaml
block_id: "TS-GO-WEB-SEO:public-indexing.test.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "AUTHORED integration of exact Next16.3.4 metadata APIs; PUBLIC_WEB_METADATA_INTEGRATION_V377.md records official authority and executed evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "26c7d509dc719b3f96b6f89dd59ca3f034476dfe87ab5cf2f29568b1c4c4eb8e"
variables: []
secrets_allowed: false
```
````typescript
import { afterEach, describe, expect, it, vi } from "vitest";
import { publicPageMetadata, publicRobots, publicSitemap, readPublicIndexing, websiteJsonLd } from "./public-indexing";

afterEach(() => vi.unstubAllEnvs());
function enable(origin = "https://public.example.invalid") {
  vi.stubEnv("PUBLIC_INDEXING_ENABLED", "1"); vi.stubEnv("PUBLIC_SITE_ORIGIN", origin);
}

describe("public indexing boundary", () => {
  it("defaults to no public crawl or fabricated URLs", () => {
    const settings = readPublicIndexing({});
    expect(publicSitemap(settings)).toEqual([]);
    expect(publicRobots(settings)).toEqual({ rules: { userAgent: "*", disallow: "/" } });
  });
  it.each(["true", "yes", "2", " 1"])("rejects ambiguous enable flag %s", flag => {
    expect(() => readPublicIndexing({ PUBLIC_INDEXING_ENABLED: flag })).toThrow("Invalid public indexing configuration");
  });
  it("requires explicit configured origin for indexing", () => {
    expect(() => readPublicIndexing({ PUBLIC_INDEXING_ENABLED: "1" })).toThrow("explicit HTTPS origin");
  });
  it.each(["http://public.example.invalid", "https://user:PRIVATE_PASSWORD@example.invalid", "https://example.invalid/private", "https://example.invalid/?secret=value", "https://example.invalid/#fragment", "https://example.invalid:444", " https://example.invalid", "https://example.invalid\n", "https://example.invalid\\private", "https://example.invalid/%2e%2e", "https://example.invalid/path/.."])("rejects non-origin input without reflecting it: %s", origin => {
    let message = "";
    try { readPublicIndexing({ PUBLIC_SITE_ORIGIN: origin }); } catch (error) { message = (error as Error).message; }
    expect(message).toBe("Invalid public site origin"); expect(message).not.toContain(origin);
  });
  it("normalizes HTTPS hostname without taking an untrusted request host", () => {
    expect(readPublicIndexing({ PUBLIC_SITE_ORIGIN: "https://PUBLIC.example.invalid:443/", PUBLIC_INDEXING_ENABLED: "1" })).toEqual({ enabled: true, origin: "https://public.example.invalid" });
  });
  it("publishes only two real pages with no invented dates or private identifiers", () => {
    enable();
    expect(publicSitemap()).toEqual([{ url: "https://public.example.invalid/" }, { url: "https://public.example.invalid/models" }]);
    expect(publicRobots()).toEqual({ rules: { userAgent: "*", disallow: "/", allow: ["/$", "/models$", "/_next/static/", "/icon.svg$"] }, sitemap: "https://public.example.invalid/sitemap.xml" });
    expect(publicPageMetadata("/models", "Modelos").alternates).toEqual({ canonical: "https://public.example.invalid/models" });
  });
  it("keeps an explicit disabled deployment out of the sitemap", () => {
    enable(); vi.stubEnv("PUBLIC_INDEXING_ENABLED", "0");
    expect(publicSitemap()).toEqual([]); expect(websiteJsonLd("Fixture", "es-AR")).toBeNull();
    expect(publicPageMetadata("/").robots).toEqual({ index: false, follow: false });
  });
  it("rejects caller attempts to canonize a protected page", () => {
    enable(); expect(() => publicPageMetadata("/admin" as "/")).toThrow("Unknown public page");
  });
  it("preserves structured text but prevents HTML script termination", () => {
    enable(); const name = '</script><script>alert("PRIVATE")</script>&\u2028\u2029';
    const encoded = websiteJsonLd(name, "es-AR")!;
    expect(encoded).not.toMatch(/[<>&\u2028\u2029]/u);
    expect(JSON.parse(encoded)).toEqual({ "@context": "https://schema.org", "@type": "WebSite", name, url: "https://public.example.invalid/", inLanguage: "es-AR" });
  });
});
````

### FILE: `src/app/robots.ts`
```yaml
block_id: "TS-GO-WEB-SEO:robots.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "AUTHORED integration of exact Next16.3.4 metadata APIs; PUBLIC_WEB_METADATA_INTEGRATION_V377.md records official authority and executed evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "5be7286d5581089c19425723f20c1ed246da0dda39c310d67d72ac06e8eafebd"
variables: []
secrets_allowed: false
```
````typescript
import { publicRobots } from "@/platform/seo/public-indexing";

export const dynamic = "force-dynamic";
export default function robots() { return publicRobots(); }
````

### FILE: `src/app/sitemap.ts`
```yaml
block_id: "TS-GO-WEB-SEO:sitemap.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "AUTHORED integration of exact Next16.3.4 metadata APIs; PUBLIC_WEB_METADATA_INTEGRATION_V377.md records official authority and executed evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "bb271c962bbe202e440746011c0e2ba13002d5a4be3ef945e4d0f6cf32344c38"
variables: []
secrets_allowed: false
```
````typescript
import { publicSitemap } from "@/platform/seo/public-indexing";

export const dynamic = "force-dynamic";
export default function sitemap() { return publicSitemap(); }
````

### FILE: `docs/public-indexing.md`
```yaml
block_id: "TS-GO-WEB-SEO:public-indexing.md:v1"
operation: CREATE
provenance: AUTHORED
source: "AUTHORED integration of exact Next16.3.4 metadata APIs; PUBLIC_WEB_METADATA_INTEGRATION_V377.md records official authority and executed evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "aaf390ab393a499c51d15741ffb0e87082c5e5e7d9584b2b726d0a52453febc0"
variables: []
secrets_allowed: false
```
````markdown
# Public indexing

This BFF uses the exact installed Next metadata/robots/sitemap APIs. It publishes
only the existing home and model-list pages; model detail URLs do not exist in
this composition. Private portals inherit noindex,nofollow. Robots directives
are crawler advice and do not replace authentication, remove indexed pages or
guarantee indexing, ranking or search-engine support.

The approved pages and their public Next static assets/icon remain crawlable;
otherwise a crawler could see the HTML but be unable to render CSS/JavaScript.

Indexing is disabled by default: no sitemap URLs, no WebSite JSON-LD, noindex on
pages and Disallow: / in robots.txt. The operator must explicitly set both
PUBLIC_INDEXING_ENABLED=1 and PUBLIC_SITE_ORIGIN to the public HTTPS origin for
the approved deployment, for example https://www.example.com. A malformed
origin or flag is a configuration error; errors never echo its value.

PUBLIC_SITE_ORIGIN is independent of the identity provider APP_BASE_URL.
Only an origin is accepted: no credentials, path prefix, non-default port,
query, fragment, whitespace or percent encoding. Request Host and forwarded
headers never define canonical URLs. Runtime routes read the current setting;
do not bake a staging host into a production sitemap. Deployment must confirm
domain ownership, TLS, redirects, canonical origin, cache behavior and target
crawler acceptance before enabling indexing.

Sitemap entries do not invent publication/modification dates, translated routes,
offers, prices or review ratings. Homepage JSON-LD contains configured business
name and locale only, is escaped for HTML and uses the existing CSP nonce.
This is a single configured site, not a domain resolver for multiple tenants.
Full localization, CMS publishing and per-model SEO require their own routes,
sources and integration evidence before extending the exact public allowlist.
````

### FILE: `src/platform/i18n/public-catalog.ts`
```yaml
block_id: "TS-GO-WEB-I18N:public-catalog.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "Local integration of existing Node/Intl/React runtime; exact qualification in PUBLIC_JOURNEY_LOCALIZATION_V378.md"
license: "LicenseRef-Workspace-Owner"
sha256: "86cb935a0f2b8d3a33076e66f2bc28858fbfd8fedd47b71cfb6da19e06f5d0e9"
variables: []
secrets_allowed: false
```
````typescript
// AUTHORED public-journey catalog. Runtime Intl supplies locale/date/plural semantics.
const es = {
  "site.description": "Catálogo y portales de una red de movilidad eléctrica.",
  "site.skip": "Saltar al contenido", "site.navigation": "Principal", "site.contact": "Contacto:",
  "nav.help": "Ayuda", "nav.dashboard": "Mi panel",
  "nav.public_catalog": "Modelos", "nav.locations": "Dónde estamos", "nav.customer": "Mi cuenta",
  "nav.admin": "Operación", "nav.franchise": "Franquicia", "nav.factory": "Fábrica",
  "home.eyebrow": "Movilidad eléctrica configurable",
  "home.description": "Explorá modelos publicados y enviá una consulta con tu consentimiento.",
  "home.action": "Ver modelos", "models.title": "Modelos eléctricos", "models.short": "Modelos",
  "models.eyebrow": "Catálogo público", "models.description": "Elegí un modelo para solicitar información.",
  "models.one": "modelo disponible", "models.other": "modelos disponibles",
  "lead.name": "Nombre", "lead.email": "Email", "lead.consent": "Acepto ser contactado sobre este modelo.",
  "lead.sending": "Enviando…", "lead.submit": "Solicitar información",
  "lead.failed": "No pudimos registrar la solicitud.",
  "lead.invalid": "La respuesta no incluyó la referencia de la solicitud.",
  "lead.received": "Solicitud recibida. Ya podés reservar una consulta o prueba.",
  "lead.next": "Reservar turno",
  "locations.eyebrow": "Red de atención", "locations.title": "Encontrá tu franquicia",
  "locations.book": "Reservar atención",
  "locations.start": "Elegí un modelo y enviá tu consulta para reservar un turno asociado.",
  "appointment.select": "Elegí un turno disponible.",
  "appointment.failed": "No pudimos reservar el turno.",
  "appointment.invalid": "No pudimos verificar la referencia del turno. Reintentá sin cambiar la selección.",
  "appointment.received": "Turno solicitado. Referencia:",
  "appointment.confirmation": "La franquicia confirmará la disponibilidad.",
  "appointment.slot": "Turno disponible", "appointment.option": "Seleccioná una opción",
  "appointment.empty": "No hay turnos publicados en los próximos 30 días.",
  "appointment.sending": "Reservando…", "appointment.submit": "Solicitar turno",
  "appointment.one": "lugar", "appointment.other": "lugares",
  "kind.consultation": "Consulta", "kind.test-drive": "Prueba de manejo",
  "kind.delivery": "Entrega", "kind.service": "Servicio"
} as const;
type MessageKey = keyof typeof es;
const en: Record<MessageKey, string> = {
  "site.description": "Catalog and portals for an electric mobility network.",
  "site.skip": "Skip to content", "site.navigation": "Main", "site.contact": "Contact:",
  "nav.help": "Help", "nav.dashboard": "My workspace",
  "nav.public_catalog": "Models", "nav.locations": "Locations", "nav.customer": "My account",
  "nav.admin": "Operations", "nav.franchise": "Franchise", "nav.factory": "Factory",
  "home.eyebrow": "Configurable electric mobility",
  "home.description": "Explore published models and send an inquiry with your consent.",
  "home.action": "View models", "models.title": "Electric models", "models.short": "Models",
  "models.eyebrow": "Public catalog", "models.description": "Choose a model to request information.",
  "models.one": "model available", "models.other": "models available",
  "lead.name": "Name", "lead.email": "Email", "lead.consent": "I agree to be contacted about this model.",
  "lead.sending": "Sending…", "lead.submit": "Request information",
  "lead.failed": "We could not register your request.",
  "lead.invalid": "The response did not include a request reference.",
  "lead.received": "Request received. You can now book a consultation or test drive.",
  "lead.next": "Book an appointment",
  "locations.eyebrow": "Service network", "locations.title": "Find your franchise",
  "locations.book": "Book an appointment",
  "locations.start": "Choose a model and send your inquiry to book a related appointment.",
  "appointment.select": "Choose an available appointment.",
  "appointment.failed": "We could not request the appointment.",
  "appointment.invalid": "We could not verify the appointment reference. Retry without changing your selection.",
  "appointment.received": "Appointment requested. Reference:",
  "appointment.confirmation": "The franchise will confirm availability.",
  "appointment.slot": "Available appointment", "appointment.option": "Select an option",
  "appointment.empty": "No appointments have been published for the next 30 days.",
  "appointment.sending": "Booking…", "appointment.submit": "Request appointment",
  "appointment.one": "place", "appointment.other": "places",
  "kind.consultation": "Consultation", "kind.test-drive": "Test drive",
  "kind.delivery": "Delivery", "kind.service": "Service"
};
export type PublicLanguage = "es" | "en";
export interface PublicLocale { language: PublicLanguage; locale: string; timeZone: string }

// Only trusted business configuration selects the locale; never headers/query/cookies.
export function resolvePublicLocale(requested: string, timeZone: string): PublicLocale {
  let locale = "es";
  try {
    const canonical = Intl.getCanonicalLocales(requested)[0];
    if (canonical && ["es", "en"].includes(new Intl.Locale(canonical).language)) locale = canonical;
  } catch { /* Unsupported or malformed configuration falls back to the actual Spanish catalog. */ }
  const language = new Intl.Locale(locale).language as PublicLanguage;
  // Invalid time zone fails configuration explicitly; do not silently book in another zone.
  new Intl.DateTimeFormat(locale, { timeZone }).format(0);
  return { language, locale, timeZone };
}

export function publicMessage(locale: string, key: string): string {
  const language = resolvePublicLocale(locale, "UTC").language;
  const catalog: Readonly<Record<string, string>> = language === "en" ? en : es;
  return Object.hasOwn(catalog, key) ? catalog[key]! : key;
}

export function publicCount(locale: string, kind: "models" | "appointment", count: number): string {
  if (!Number.isSafeInteger(count) || count < 0) throw new RangeError("Invalid public count");
  const resolved = resolvePublicLocale(locale, "UTC");
  const category = new Intl.PluralRules(resolved.locale).select(count);
  const suffix = publicMessage(resolved.locale, kind + (category === "one" ? ".one" : ".other"));
  return new Intl.NumberFormat(resolved.locale).format(count) + " " + suffix;
}

export function publicAppointmentTime(config: PublicLocale, startsAt: string): string {
  const date = new Date(startsAt);
  if (!Number.isFinite(date.getTime())) throw new RangeError("Invalid appointment date");
  return new Intl.DateTimeFormat(config.locale, {
    dateStyle: "medium", timeStyle: "short", timeZone: config.timeZone
  }).format(date);
}

````

### FILE: `src/platform/i18n/load-public-locale.ts`
```yaml
block_id: "TS-GO-WEB-I18N:load-public-locale.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "Local integration of existing Node/Intl/React runtime; exact qualification in PUBLIC_JOURNEY_LOCALIZATION_V378.md"
license: "LicenseRef-Workspace-Owner"
sha256: "3920d0098d710223ed2b975d63f05a81ecbf83d7c1c54a1f1b52864ecb42acf3"
variables: []
secrets_allowed: false
```
````typescript
import { loadBusinessConfig } from "@/platform/config/load";
import { resolvePublicLocale } from "./public-catalog";

export async function loadPublicLocale() {
  const config = await loadBusinessConfig();
  const market = config.markets.find((item) => item.code === config.business.defaultMarket);
  if (!market) throw new Error("Configured market is missing");
  return resolvePublicLocale(config.business.defaultLocale, market.timeZone);
}

````

### FILE: `src/platform/i18n/public-catalog.test.ts`
```yaml
block_id: "TS-GO-WEB-I18N:public-catalog.test.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "Local integration of existing Node/Intl/React runtime; exact qualification in PUBLIC_JOURNEY_LOCALIZATION_V378.md"
license: "LicenseRef-Workspace-Owner"
sha256: "e3c58a0c95e9f889c55d0701e8a3dba0e66a7cb5e98fb5e827e83b8959ec197a"
variables: []
secrets_allowed: false
```
````typescript
import { describe, expect, it } from "vitest";
import { publicAppointmentTime, publicCount, publicMessage, resolvePublicLocale } from "./public-catalog";

describe("public journey localization", () => {
  it.each(["en-US", "en-GB", "EN-au"])("uses the English catalog for %s", (locale) => {
    expect(publicMessage(locale, "lead.name")).toBe("Name");
    expect(resolvePublicLocale(locale, "UTC").language).toBe("en");
  });
  it.each(["es-AR", "es-ES", "ES-mx"])("uses the Spanish catalog for %s", (locale) => {
    expect(publicMessage(locale, "lead.name")).toBe("Nombre");
    expect(resolvePublicLocale(locale, "UTC").language).toBe("es");
  });
  it.each(["fr-FR", "xx", "", "not_a_locale", "__proto__", "<script>"])("falls back honestly for %s", (locale) => {
    expect(resolvePublicLocale(locale, "UTC")).toEqual({ language: "es", locale: "es", timeZone: "UTC" });
    expect(publicMessage(locale, "lead.name")).toBe("Nombre");
  });
  it("returns unknown keys literally without reading prototypes", () => {
    expect(publicMessage("en-US", "missing")).toBe("missing");
    expect(publicMessage("en-US", "__proto__")).toBe("__proto__");
    expect(publicMessage("en-US", "constructor")).toBe("constructor");
  });
  it.each([
    ["en-US", 0, "0 models available"], ["en-US", 1, "1 model available"],
    ["en-US", 2, "2 models available"], ["es-AR", 0, "0 modelos disponibles"],
    ["es-AR", 1, "1 modelo disponible"], ["es-AR", 2, "2 modelos disponibles"]
  ])("formats %s count %s", (locale, count, expected) => {
    expect(publicCount(String(locale), "models", Number(count))).toBe(expected);
  });
  it.each([-1, 1.5, NaN, Infinity, Number.MAX_SAFE_INTEGER + 1])("rejects invalid count %s", (count) => {
    expect(() => publicCount("en", "models", count)).toThrow(RangeError);
  });
  it("keeps explicit configured time zones and crosses daylight saving", () => {
    const config = resolvePublicLocale("en-US", "America/New_York");
    const winter = publicAppointmentTime(config, "2026-01-15T15:00:00Z");
    const summer = publicAppointmentTime(config, "2026-07-15T15:00:00Z");
    expect(winter).toContain("10:00");
    expect(summer).toContain("11:00");
    expect(publicAppointmentTime(resolvePublicLocale("es-AR", "America/Argentina/Buenos_Aires"), "2026-01-15T15:00:00Z")).toContain("12:00");
  });
  it("rejects invalid dates and time zones without a silent fallback", () => {
    expect(() => resolvePublicLocale("en", "Mars/Olympus")).toThrow(RangeError);
    expect(() => publicAppointmentTime(resolvePublicLocale("en", "UTC"), "bad")).toThrow(RangeError);
  });
  it("does not share mutable locale state across concurrent renders", async () => {
    const values = await Promise.all(Array.from({ length: 100 }, (_, i) =>
      Promise.resolve(publicMessage(i % 2 ? "en-GB" : "es-AR", "lead.name"))));
    expect(values.filter((v) => v === "Name")).toHaveLength(50);
    expect(values.filter((v) => v === "Nombre")).toHaveLength(50);
  });
});

````

### FILE: `docs/public-localization.md`
```yaml
block_id: "TS-GO-WEB-I18N:public-localization.md:v1"
operation: CREATE
provenance: AUTHORED
source: "Local integration of existing Node/Intl/React runtime; exact qualification in PUBLIC_JOURNEY_LOCALIZATION_V378.md"
license: "LicenseRef-Workspace-Owner"
sha256: "916e523873609e540f7a547565d8e78c144442909cd53e0875479c446bb56024"
variables: []
secrets_allowed: false
```
````markdown
# Public journey localization

AUTHORED integration under the existing BFF and journey-portal owners. No new dependency.

The configured business.defaultLocale selects Spanish or English, including regional Intl formatting.
Unsupported or malformed tags fall back to Spanish and are labeled as Spanish. Unknown message keys remain
literal keys. The fixed catalogs contain complete Spanish and English messages; they are not an arbitrary
translation CMS. Intl.PluralRules and NumberFormat handle integer counts, with one/other messages for these
two languages. No claim of general CLDR translation coverage or interpolation.

The default market's explicit timeZone controls displayed appointment time on both server and client.
Invalid time zones fail configuration; no guessed business time zone. The backend's starts_at and kind,
lead and model IDs, consent boolean, endpoints and idempotency keys remain unchanged.
The displayed localized time is never parsed back into a request.

The localized reference journey covers home, models, connected catalog, lead request, locations and appointment
request, including failures, empty state, pending state and durable receipt. Product names and vehicle-class
content come from the backend unchanged. Navigation, footer, metadata and JSON-LD use the resolved language.
Private portals retain their existing Spanish content under main lang=es; public sections override that
language explicitly. Customer/operator portal translation, user language negotiation, URL language variants,
hreflang and multilingual content storage require their own integration and target review.
Consent wording is a translation of the existing technical fixture, not a legal approval.

References: TC39 ECMA-402 Intl specification, Node internationalization documentation, W3C declaring
language in HTML. Exact runtime behavior and connected tests are recorded by V378. Existing native/runtime
admission conditions are retained; using Intl does not close their security or distribution gates.

````

### FILE: `src/app/admin/page.test.ts`
```yaml
block_id: "TS-GO-API-WEB-BRIDGE:PRIVATE-READ:src-app-admin-page.test.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "Local private-read permission repair and exact existing quote-format extraction; see PRIVATE_PORTAL_READS_V382.md"
license: "LicenseRef-Workspace-Owner"
sha256: "ba60870ccf806ecb0555f42e06346a371fbe1cad0ee0bc6f559cb488b7c23e8c"
variables: []
secrets_allowed: false
```
````typescript
import { beforeEach, expect, test, vi } from "vitest";
import { renderToStaticMarkup } from "react-dom/server";
const mocks=vi.hoisted(()=>({readSession:vi.fn(),get:vi.fn()}));
vi.mock("@/platform/auth/session",()=>({readSession:mocks.readSession,allowed:(s:{permissions:string[]},p:string)=>s.permissions.includes(p)||s.permissions.includes("*")}));
vi.mock("@/platform/backend/protected-client",()=>({protectedGet:mocks.get}));
vi.mock("next/navigation",()=>({redirect:(path:string)=>{throw new Error("REDIRECT "+path)}}));
import AdminPage from "./page";
const session=(permissions:string[])=>({subject:"synthetic-admin",tenantId:"synthetic-tenant",organizations:["org-a"],permissions,accessToken:"private-fixture-token"});
beforeEach(()=>{vi.resetAllMocks();mocks.readSession.mockResolvedValue(session(["admin:read"]));mocks.get.mockImplementation(async(s:{permissions:string[]},path:string)=>{
 if(path==="/v1/franchise/leads"){
  if(!s.permissions.includes("lead:read")&&!s.permissions.includes("*"))throw new Error("FORBIDDEN lead:read");
  return {items:[{id:"lead-visible",state:"new",source_code:"fixture",assigned_subject:""}]};
 }
 if(path==="/v1/admin/overview")return {orders:2,stock_available:1,open_cases:1,active_shipments:0};
 if(path==="/v1/admin/orders")return {items:[{id:"order-visible",state:"draft",currency:"USD",total_minor_units:100}]};
 if(path==="/v1/admin/service-cases")return {items:[{id:"case-visible",state:"opened",severity:"medium"}]};
 throw new Error("UNEXPECTED PATH "+path);
});});
test.each([["admin:read"],["admin:read","Lead:Read"]].map(permissions=>({permissions})))("admin without exact lead access keeps its permitted data: $permissions",async({permissions})=>{
 mocks.readSession.mockResolvedValue(session(permissions));const html=renderToStaticMarkup(await AdminPage());
 expect(html).toContain("order-visible");expect(html).toContain("case-visible");expect(html).not.toContain("Oportunidades");expect(html).not.toContain("lead-visible");
 expect(mocks.get.mock.calls.map(c=>c[1]).sort()).toEqual(["/v1/admin/overview","/v1/admin/orders","/v1/admin/service-cases"].sort());
 expect(html).not.toContain("private-fixture-token");
});
test.each([["admin:read","lead:read"],["*"]].map(permissions=>({permissions})))("authorized lead view remains available: $permissions",async({permissions})=>{
 mocks.readSession.mockResolvedValue(session(permissions));const html=renderToStaticMarkup(await AdminPage());expect(html).toContain("Oportunidades");expect(html).toContain("lead-visible");expect(mocks.get).toHaveBeenCalledTimes(4);
 for(const call of mocks.get.mock.calls){expect(call[0].organizations).toEqual(["org-a"]);expect(call[2].organization_id).toBe("org-a");}
});
test.each([["lead:read"],["customer:self"],["factory:read"],[]].map(permissions=>({permissions})))("non-admin session cannot issue admin reads: $permissions",async({permissions})=>{
 mocks.readSession.mockResolvedValue(session(permissions));const html=renderToStaticMarkup(await AdminPage());expect(html).toContain("Acceso denegado");expect(mocks.get).not.toHaveBeenCalled();
});
test("missing session redirects before backend reads",async()=>{mocks.readSession.mockResolvedValue(null);await expect(AdminPage()).rejects.toThrow("REDIRECT /api/auth/login?return_to=/admin");expect(mocks.get).not.toHaveBeenCalled();});

 test("admin renders minor units as the actual currency amount",async()=>{const html=renderToStaticMarkup(await AdminPage()).replace(/\u00a0/g," ");expect(html).toContain("USD 1,00");expect(html).not.toContain("USD 100");});

test("customer order and quote summaries use the same verified minor-unit display",async()=>{
 const {default:CustomerPage}=await import("../customer/page");mocks.readSession.mockResolvedValue(session(["customer:self"]));
 mocks.get.mockImplementation(async(_session:unknown,path:string)=>{
  if(path==="/v1/customer/orders")return {items:[{id:"customer-order",state:"draft",currency:"USD",total_minor_units:100}]};
  if(path==="/v1/customer/service-cases")return {items:[]};
  if(path==="/v1/customer/journey")return {appointments:[],handovers:[],quotes:[{id:"customer-quote",state:"issued",currency:"KWD",total_minor_units:1234}]};
  throw new Error("UNEXPECTED CUSTOMER PATH "+path);
 });
 const html=renderToStaticMarkup(await CustomerPage()).replace(/\u00a0/g," ");expect(html).toContain("USD 1,00");expect(html).toContain("KWD 1,234");expect(html).not.toContain("USD 100");expect(html).not.toContain("KWD 1234");
});

test("admin pages expose backend cursors for all granted lists",async()=>{
 mocks.readSession.mockResolvedValue(session(["admin:read","lead:read"]));const old=mocks.get.getMockImplementation()!;
 mocks.get.mockImplementation(async(...args:unknown[])=>({...await old(...args),next_cursor:String(args[1]).endsWith("orders")?"order-next":String(args[1]).endsWith("service-cases")?"case-next":"lead-next"}));
 const html=renderToStaticMarkup(await AdminPage());
 expect(html).toContain('href="/admin?orders_after=order-next"');expect(html).toContain('href="/admin?cases_after=case-next"');expect(html).toContain('href="/admin?leads_after=lead-next"');
});

test("admin keeps all granted cursors independent and fixed session scope",async()=>{
 mocks.readSession.mockResolvedValue(session(["admin:read","lead:read"]));
 await AdminPage({searchParams:Promise.resolve({orders_after:"o",cases_after:"c",leads_after:"l",organization_id:"ignored",limit:"1000"})});
 for(const [path,after] of [["/v1/admin/orders","o"],["/v1/admin/service-cases","c"],["/v1/franchise/leads","l"]])expect(mocks.get).toHaveBeenCalledWith(expect.anything(),path,{organization_id:"org-a",limit:"25",after});
});
test("admin without lead permission neither reads nor preserves lead cursor",async()=>{
 const html=renderToStaticMarkup(await AdminPage({searchParams:Promise.resolve({orders_after:"o",leads_after:["ignored","duplicate"]})}));
 expect(html).not.toContain("leads_after");expect(html).not.toContain("Oportunidades");expect(mocks.get.mock.calls.some(c=>c[1]==="/v1/franchise/leads")).toBe(false);
});
test("admin final-page links retain the other lists",async()=>{mocks.readSession.mockResolvedValue(session(["admin:read","lead:read"]));const html=renderToStaticMarkup(await AdminPage({searchParams:Promise.resolve({orders_after:"o",cases_after:"c",leads_after:"l"})}));expect(html).toContain('href="/admin?cases_after=c&amp;leads_after=l"');expect(html).toContain('href="/admin?orders_after=o&amp;leads_after=l"');expect(html).toContain('href="/admin?orders_after=o&amp;cases_after=c"');expect(html).not.toContain("Ver siguientes");});
````

### FILE: `src/platform/i18n/money.ts`
```yaml
block_id: "TS-GO-API-WEB-BRIDGE:PRIVATE-READ:src-platform-i18n-money.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "Local private-read permission repair and exact existing quote-format extraction; see PRIVATE_PORTAL_READS_V382.md"
license: "LicenseRef-Workspace-Owner"
sha256: "bbdb7473773c9b7d1d15c39be9060cfa3f9926b93b057048465b34685a42806b"
variables: []
secrets_allowed: false
```
````typescript
// AUTHORED exact extraction of the existing quote display, shared by private read views.
// Pinned Intl currency digits and BigInt parts preserve all safe minor-unit integers.
// No price, FX, tax, business rounding or acceptance policy is introduced.
export function minorAmountPresentation(minorUnits: number, currency: string) {
  let amountLabel = "Importe no verificable; consultá al soporte.";
  let amountValid = false;
  if (Number.isSafeInteger(minorUnits) && minorUnits >= 0 &&
      Intl.supportedValuesOf("currency").includes(currency)) {
    const formatter = new Intl.NumberFormat("es-AR", { style: "currency", currency: currency, currencyDisplay: "code" });
    const digits = formatter.resolvedOptions().maximumFractionDigits;
    if (typeof digits === "number" && Number.isInteger(digits) && digits >= 0 && digits <= 20) {
      const minor = BigInt(minorUnits);
      const scale = 10n ** BigInt(digits);
      const fraction = (minor % scale).toString().padStart(digits, "0");
      amountLabel = formatter.formatToParts(minor / scale).map(part => part.type === "fraction" ? fraction : part.value).join("");
      amountValid = true;
    }
  }
  return { amountLabel, amountValid };
}
````

### FILE: `src/platform/i18n/money.test.ts`
```yaml
block_id: "TS-GO-API-WEB-BRIDGE:PRIVATE-READ:src-platform-i18n-money.test.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "Local private-read permission repair and exact existing quote-format extraction; see PRIVATE_PORTAL_READS_V382.md"
license: "LicenseRef-Workspace-Owner"
sha256: "01eb0d0e1734bf70fa55e1ad62a0eee5599bdc189ef8e4b64132225c65bc2697"
variables: []
secrets_allowed: false
```
````typescript
import { expect, test } from "vitest";
import { minorAmountPresentation } from "./money";
test.each([
 {minor:100,currency:"USD",label:"USD 1,00"},
 {minor:0,currency:"USD",label:"USD 0,00"},
 {minor:123,currency:"JPY",label:"JPY 123"},
 {minor:1234,currency:"KWD",label:"KWD 1,234"},
 {minor:9007199254740991,currency:"USD",label:"USD 90.071.992.547.409,91"},
 {minor:100,currency:"CLP",label:"CLP 100"},
])("exact existing quote representation: $currency $minor",({minor,currency,label})=>{const result=minorAmountPresentation(minor,currency);expect(result.amountLabel.replace(/\u00a0/g," ")).toBe(label);expect(result.amountValid).toBe(true);});
test.each([
 {minor:-1,currency:"USD"},{minor:1.5,currency:"USD"},{minor:NaN,currency:"USD"},
 {minor:Infinity,currency:"USD"},{minor:9007199254740992,currency:"USD"},
 {minor:100,currency:"unknown"},{minor:100,currency:"usd"},
])("invalid display cannot imply a verified amount: $currency $minor",({minor,currency})=>{expect(minorAmountPresentation(minor,currency)).toEqual({amountLabel:"Importe no verificable; consultá al soporte.",amountValid:false});});
````

### FILE: `docs/private-portal-reads.md`
```yaml
block_id: "TS-GO-API-WEB-BRIDGE:PRIVATE-READ:docs-private-portal-reads.md:v1"
operation: CREATE
provenance: AUTHORED
source: "Local private-read permission repair and exact existing quote-format extraction; see PRIVATE_PORTAL_READS_V382.md"
license: "LicenseRef-Workspace-Owner"
sha256: "2f606a960122c8c7c28d4cbc518d694309cd73240479b7227cd1d460a8f57566"
variables: []
secrets_allowed: false
```
````markdown
# Private portal read views — V382

The administrative route requires admin:read. Overview, orders and service cases use that
grant; its optional opportunities section additionally requires lead:read (or exact wildcard).
The BFF does not request that endpoint without its permission. Roles and encrypted frontend
claims never enlarge the actual bearer grant enforced by Go. A session that overclaims lead
access while its bearer lacks it still fails closed; no implicit fallback or permission grant.

The first organization in the current verified session selects the existing query scope.
The real HTTP layer verifies tenant, permission and organization; customer reads also bind
the bearer subject. This change adds no organization switcher, business role or IdP lifecycle.

Admin orders and customer order/quote summaries now use the exact existing quote amount
algorithm shared under platform/i18n/money.ts. It checks safe nonnegative minor-unit integers
and supported currency codes, derives currency digits from the pinned Intl runtime, and uses
BigInt plus formatToParts to retain the remainder without floating-point division. Existing
Spanish presentation and invalid-amount wording are preserved. USD100minor units displays
USD1,00; JPY123 displaysJPY123; KWD1234 displaysKWD1,234. No price, FX, tax, acceptance,
business rounding or stored amount changes. The quote-page extraction is exactly reversible.

Qualification:11private-view regressions and13literal amount/invalid-input cases;260full web
tests pass with one pre-existing explicit connected skip, typecheck and build. The actual
Next/JWE to Go/RS256/JWKS to PostgreSQL fixture uses53unchanged migrations, five independent
HTTP negatives and four browser projects with retries0. It tests admin-only, combined/wildcard,
case-sensitive grants, other organization/tenant, non-admin denial, factory and customer reads,
customer quote display, permission removal on reload, backend rejection of overclaimed frontend
permissions and subsequent recovery. No browser/API writes; snapshots across six domain tables
are unchanged. The repository fixture is durable and synthetic, not an external IdP/provider.

Run Go TestAdministrativeReadBrowserPostgres with ELITE_ADMIN_READ_E2E=1, absolute ELITE_WEB_ROOT,
an owned loopback elite_confirmation_* TEST_DATABASE_URL, the composed web build and installed
admitted Playwright gate. Enable role_workspace/customer_portal/factory_portal/public_catalog
and the catalog/crm/procurement modules in the isolated business config. The fixture reuses the
selected journey test issuer/clock, starts its owned API/TLS/Next, runs test:admin-reads and
stops its processes. The caller owns migration/setup and PostgreSQL shutdown. Never use a
project database, real identities, production keys or the synthetic issuer as deployment.

Rollback must keep a compatible composition. Preserve the old source for diagnosis, but do not
promote the known admin permission coupling or the misleading minor-unit labels as a rollback.
On regression, block promotion and repair the candidate. Authorization errors, IdP revocation,
arbitrary analytics/KPIs, full operational workflows, monetary policies, accessibility audit,
native licence/security, release and target acceptance retain their own gates. This read boundary
does not close any of those by implication.

V383 adds fixed-text Next error boundaries for admin/customer/factory. They never read or display the thrown error, stack, digest or bearer. An ordinary anchor to the fixed portal root performs a new GET only after the user clicks; no reset hook, automatic retry or operation replay. The message does not claim success or failure of a previous write. The existing page rechecks the current session and the backend still enforces its bearer. A segment error can cover descendants too: recovery deliberately returns to the portal root, not an arbitrary failed URL. Root-layout failures and external IdP recovery remain separate.

Verified against the pinned Next16.3.4 local ErrorComponent contract, actual production build and Go/JWKS/PostgreSQL in four browser projects. For all three portal roots an overclaimed fixture session with an insufficient bearer fails without domain records, an unchanged reconsult still fails, then changing to a valid session and clicking the same link issues GET and restores scoped data. No browser/API writes and six durable-table snapshots unchanged.263web tests plus one inherited explicit skip, typecheck/build and canonical reconstruction pass. This proves the exercised permission-rejection recovery, not an injected database outage or complete workflow/security/release admission.

V387 connects the existing cursor API for customer orders/service cases and factory units. Fixed25row requests retain session organization and bearer; orders_after/cases_after are independent, units_after belongs to factory. Ordinary GET links preserve the other list, reload retains position and first-page recovery clears only the chosen cursor. Opaque cursor values are URL-encoded and bounded to512characters; duplicates/excessive values enter existing fixed error recovery. Unknown query inputs cannot replace organization or page size. Final/empty pages have no next link. No write, offset, count, cross-page snapshot, permission, currency or warranty-policy guarantee is added. Concurrent data changes retain the existing keyset semantics. Admin lists and the separate journey array projection are outside this correction.

The connected fixture holds27records per list:81unique records per browser, no omissions/duplicates across first/last pages; four projects also exercise independent cursor retention, reload, first-page navigation, existing scopes, fixed read recovery and unchanged six-table snapshot. No provider or production access. The qualification is independent of the V386 deferred native vulnerability research and cannot close TEST03/07.

V388 supersedes the V387 admin-list exclusion: admin orders/service cases/leads now use the same cursor navigator. Three independent query cursors retain the other granted lists. The lead cursor is read/preserved only with lead:read; revoking that permission removes its list and cursor from subsequent links.25row requests and session organization remain fixed. Four added tests and the extended actual browser/Go/PostgreSQL fixture traverse28admin orders,27cases and27leads (82rows across these lists), while all81customer/factory observations remain tested. Total163list records per browser, not163distinct underlying business entities. Reload/first-page recovery and no writes/snapshot preservation pass. No domain API, business rule, source graph or dependency change.

V389: customer account and appointment management render the same stored appointment instant using the admitted default locale and market timeZone, shared with public booking. Each time element retains the original ISO datetime and shows the explicit market zone. The interactive cancellation component receives the server-formatted label, so browser locale/timeZone cannot change it during hydration. Date validity is checked by the existing formatter. Six added rendering regressions cover day boundaries, New York winter/summer offsets, Tokyo, invalid dates and access checks. Two actual business configurations, each in four browser projects with Honolulu/de-DE clients and UTC server, preserve two appointment labels across account/management/reload. All seven table snapshots remain unchanged; zero writes. Cancellation command, authorization, timestamps, versions and eligibility rules are unchanged. Other private portal dates and translation catalogs remain outside this narrow repair.

V390 closes the existing customer cancellation path in the isolated reference. Account links to appointment management; no-results state is explicit. A synchronous in-flight guard allows one command; the client confirms only a receipt bound to appointment ID, organization, cancelled state and exact successor version. A lost/unverifiable/error response leaves the outcome unknown and commands disabled; an explicit GET reload resolves current state. There is no automatic mutation retry. The existing10second bounded request pattern is reused. Four actual browser projects prove normal cancellation/double click, committed response loss, unbound HTTP200 response, two concurrent requests and stale-page recovery. Sixteen appointments produce exactly16audits/16outbox events;28API writes include rejected duplicate/stale/foreign-resource attempts. Four browser/BFF boundary negatives per project and five inherited direct HTTP negatives remain enforced. All other durable data and appointment immutable fields are unchanged. The prior six-list/time-zone read suites also pass in eight browser/configuration runs with zero writes. No backend domain code, migration, cancellation eligibility, identity grant, provider or dependency changed. This is narrow reference closure, not full TEST02/production admission.
````

### FILE: `src/platform/backend/private-read-failure.tsx`
```yaml
block_id: "TS-GO-API-WEB-BRIDGE:PRIVATE-READ:src-platform-backend-private-read-failure.tsx:v1"
operation: CREATE
provenance: AUTHORED
source: "Local private-read recovery using fixed Next error boundaries and explicit GET; see PRIVATE_READ_RECOVERY_V383.md"
license: "LicenseRef-Workspace-Owner"
sha256: "ab6da4af4dffdd3c4ebd5d365511cc892faf620ca7f663807604e93e2380583a"
variables: []
secrets_allowed: false
```
````typescript
type PrivateReadPath = "/admin" | "/customer" | "/factory";

export function PrivateReadFailure({ path }: { path: PrivateReadPath }) {
  return <section aria-labelledby="private-read-title">
    <h1 id="private-read-title" className="pageTitle">No pudimos cargar esta información</h1>
    <p className="notice" role="alert">La consulta no se completó. Podés volver a consultar con tu sesión actual.</p>
    <p>Este mensaje no confirma el resultado de una operación anterior. Consultá su estado antes de repetirla.</p>
    <a className="button" href={path}>Volver a consultar</a>
  </section>;
}
````

### FILE: `src/platform/backend/private-read-failure.test.ts`
```yaml
block_id: "TS-GO-API-WEB-BRIDGE:PRIVATE-READ:src-platform-backend-private-read-failure.test.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "Local private-read recovery using fixed Next error boundaries and explicit GET; see PRIVATE_READ_RECOVERY_V383.md"
license: "LicenseRef-Workspace-Owner"
sha256: "4859cf3c3aedde3d243e59a3932a007b2e80f2b0479cff2d1e1a6ee7e884d316"
variables: []
secrets_allowed: false
```
````typescript
import { createElement } from "react";
import { renderToStaticMarkup } from "react-dom/server";
import { expect, test, vi } from "vitest";
import AdminError from "@/app/admin/error";
import CustomerError from "@/app/customer/error";
import FactoryError from "@/app/factory/error";

test.each([
  { path: "/admin", component: AdminError },
  { path: "/customer", component: CustomerError },
  { path: "/factory", component: FactoryError }
])("$path recovery never displays the supplied failure or calls a retry", ({ path, component }) => {
  const reset = vi.fn(); const retry = vi.fn();
  const error = Object.assign(new Error("private-bearer customer-secret database-host"), { digest: "private-digest" });
  const html = renderToStaticMarkup(createElement(component, { ...{ error, reset, retry } }));
  expect(html).toContain('role="alert"');
  expect(html).toContain('href="'+path+'"');
  expect(html).toContain("Volver a consultar");
  expect(html).toContain("no confirma el resultado de una operación anterior");
  expect(html).not.toMatch(/private-bearer|customer-secret|database-host|private-digest|<form|<button/);
  expect(reset).not.toHaveBeenCalled(); expect(retry).not.toHaveBeenCalled();
});
````

### FILE: `src/app/admin/error.tsx`
```yaml
block_id: "TS-GO-API-WEB-BRIDGE:PRIVATE-READ:src-app-admin-error.tsx:v1"
operation: CREATE
provenance: AUTHORED
source: "Local private-read recovery using fixed Next error boundaries and explicit GET; see PRIVATE_READ_RECOVERY_V383.md"
license: "LicenseRef-Workspace-Owner"
sha256: "819318819fd95ef0c6a5fdcd8d8ec0576d60f4b1f87abf09df4a80feac455dc5"
variables: []
secrets_allowed: false
```
````typescript
"use client";

import { PrivateReadFailure } from "@/platform/backend/private-read-failure";

export default function PrivatePortalError() {
  return <PrivateReadFailure path="/admin" />;
}
````

### FILE: `src/app/customer/error.tsx`
```yaml
block_id: "TS-GO-API-WEB-BRIDGE:PRIVATE-READ:src-app-customer-error.tsx:v1"
operation: CREATE
provenance: AUTHORED
source: "Local private-read recovery using fixed Next error boundaries and explicit GET; see PRIVATE_READ_RECOVERY_V383.md"
license: "LicenseRef-Workspace-Owner"
sha256: "63c610078bf34614b0a53dad36bac92ce40c62b21f2f1fef62dcaa8ef6127c14"
variables: []
secrets_allowed: false
```
````typescript
"use client";

import { PrivateReadFailure } from "@/platform/backend/private-read-failure";

export default function PrivatePortalError() {
  return <PrivateReadFailure path="/customer" />;
}
````

### FILE: `src/app/factory/error.tsx`
```yaml
block_id: "TS-GO-API-WEB-BRIDGE:PRIVATE-READ:src-app-factory-error.tsx:v1"
operation: CREATE
provenance: AUTHORED
source: "Local private-read recovery using fixed Next error boundaries and explicit GET; see PRIVATE_READ_RECOVERY_V383.md"
license: "LicenseRef-Workspace-Owner"
sha256: "25624fb9b41f0b14c000baf2905e4dc90722eb27d1731cdd26771a18b9494bf7"
variables: []
secrets_allowed: false
```
````typescript
"use client";

import { PrivateReadFailure } from "@/platform/backend/private-read-failure";

export default function PrivatePortalError() {
  return <PrivateReadFailure path="/factory" />;
}
````

### FILE: `src/platform/backend/portal-paging.tsx`
```yaml
block_id: "TS-GO-API-WEB-BRIDGE:PAGING:src-platform-backend-portal-paging.tsx:v1"
operation: CREATE
provenance: AUTHORED
source: "Local cursor navigation for existing scoped API; PRIVATE_PORTAL_PAGINATION_V387.md"
license: "LicenseRef-Workspace-Owner"
sha256: "a35a512b99ce85e0d41706b0ee96a0b1a61015155ca2fcfabf59031e9512d999"
variables: []
secrets_allowed: false
```
````typescript
export type PortalQuery = Record<string, string | string[] | undefined>;
export type PortalPageProps = { searchParams?: Promise<PortalQuery> };
type CursorName = "orders_after" | "cases_after" | "units_after" | "leads_after";

export function portalCursor(value: string | string[] | undefined): string | undefined {
  if (value === undefined || value === "") return undefined;
  if (typeof value !== "string" || value.length > 512) throw new Error("Invalid page cursor");
  return value;
}

export function portalPageHref(path: "/customer" | "/factory" | "/admin", cursors: Partial<Record<CursorName, string | undefined>>): string {
  const query = new URLSearchParams();
  for (const name of path === "/admin" ? ["orders_after", "cases_after", "leads_after"] as const : path === "/customer" ? ["orders_after", "cases_after"] as const : ["units_after"] as const) {
    const value = portalCursor(cursors[name]);
    if (value !== undefined) query.set(name, value);
  }
  return path + (query.size ? "?" + query.toString() : "");
}

export function PortalPaging({label,nextHref,firstHref}:{label:string;nextHref:string|undefined;firstHref:string|undefined}) {
  if (!nextHref && !firstHref) return null;
  return <nav aria-label={label}>{firstHref ? <p><a href={firstHref}>Volver al inicio</a></p> : null}{nextHref ? <p><a href={nextHref}>Ver siguientes</a></p> : null}</nav>;
}
````

### FILE: `src/app/customer/page.test.ts`
```yaml
block_id: "TS-GO-API-WEB-BRIDGE:PAGING:src-app-customer-page.test.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "Local cursor navigation for existing scoped API; PRIVATE_PORTAL_PAGINATION_V387.md"
license: "LicenseRef-Workspace-Owner"
sha256: "f7df56048956ca0c6eb1ba6bfe541802abdc756dd4705ef633b77f1bad441a32"
variables: []
secrets_allowed: false
```
````typescript
import {beforeEach,expect,test,vi} from "vitest";
import {renderToStaticMarkup} from "react-dom/server";
const mocks=vi.hoisted(()=>({session:vi.fn(),get:vi.fn(),locale:vi.fn()}));
vi.mock("@/platform/auth/session",()=>({readSession:mocks.session,allowed:(s:{permissions:string[]},p:string)=>s.permissions.includes(p)}));
vi.mock("@/platform/backend/protected-client",()=>({protectedGet:mocks.get}));
vi.mock("next/navigation",()=>({redirect:(p:string)=>{throw new Error("REDIRECT "+p)}}));
vi.mock("@/platform/i18n/load-public-locale",()=>({loadPublicLocale:mocks.locale}));
import AppointmentsPage from "./appointments/page";
import CustomerPage from "./page";
import FactoryPage from "../factory/page";
beforeEach(()=>{vi.resetAllMocks();mocks.locale.mockResolvedValue({language:"es",locale:"es-AR",timeZone:"America/Argentina/Buenos_Aires"});mocks.session.mockResolvedValue({subject:"customer",tenantId:"tenant",organizations:["org-a"],permissions:["customer:self","factory:read"],accessToken:"test-only"});mocks.get.mockImplementation(async(_s:unknown,p:string)=>p.endsWith("journey")?{appointments:[],quotes:[],handovers:[]}:{items:[],next_cursor:p.endsWith("orders")?"order-025":p.endsWith("service-cases")?"case-025":"unit-025"});});
test("customer exposes each backend next cursor independently",async()=>{
 const html=renderToStaticMarkup(await CustomerPage());
 expect(html).toContain('href="/customer?orders_after=order-025"');expect(html).toContain('href="/customer?cases_after=case-025"');
});
test("factory exposes backend next cursor",async()=>{expect(renderToStaticMarkup(await FactoryPage())).toContain('href="/factory?units_after=unit-025"');});

test("customer forwards opaque cursors separately and preserves the other list",async()=>{
 const html=renderToStaticMarkup(await CustomerPage({searchParams:Promise.resolve({orders_after:"old order/&",cases_after:"old-case",organization_id:"ignored",limit:"1000"})}));
 const calls=mocks.get.mock.calls;expect(calls.find(c=>c[1].endsWith("orders"))?.[2]).toEqual({organization_id:"org-a",limit:"25",after:"old order/&"});
 expect(calls.find(c=>c[1].endsWith("service-cases"))?.[2]).toEqual({organization_id:"org-a",limit:"25",after:"old-case"});
 expect(html).toContain('/customer?orders_after=order-025&amp;cases_after=old-case');expect(html).toContain('/customer?orders_after=old+order%2F%26&amp;cases_after=case-025');
 expect(html).toContain('href="/customer?cases_after=old-case"');expect(html).toContain('href="/customer?orders_after=old+order%2F%26"');
});
test("last and empty pages have no next link and retain first-page recovery",async()=>{
 mocks.get.mockImplementation(async(_s:unknown,p:string)=>p.endsWith("journey")?{appointments:[],quotes:[],handovers:[]}:{items:[]});
 const html=renderToStaticMarkup(await CustomerPage({searchParams:Promise.resolve({orders_after:"last"})}));expect(html).not.toContain('Ver siguientes');expect(html).toContain('href="/customer"');expect(html).toContain('No hay registros en esta página.');
 const factory=renderToStaticMarkup(await FactoryPage({searchParams:Promise.resolve({units_after:"last"})}));expect(factory).not.toContain('Ver siguientes');expect(factory).toContain('href="/factory"');
});
test("factory reads the session organization and fixed limit",async()=>{await FactoryPage({searchParams:Promise.resolve({units_after:"unit-025",organization_id:"ignored",limit:"1000"})});expect(mocks.get).toHaveBeenCalledWith(expect.anything(),"/v1/factory/units",{organization_id:"org-a",limit:"25",after:"unit-025"});});
test.each([{value:["duplicate","cursor"]},{value:"x".repeat(513)}])("ambiguous or excessive cursors stop before reads",async({value})=>{await expect(CustomerPage({searchParams:Promise.resolve({orders_after:value})})).rejects.toThrow("Invalid page cursor");expect(mocks.get).not.toHaveBeenCalled();});
test("missing session still redirects before querying",async()=>{mocks.session.mockResolvedValue(null);await expect(CustomerPage({searchParams:Promise.resolve({orders_after:"cursor"})})).rejects.toThrow("REDIRECT");expect(mocks.get).not.toHaveBeenCalled();});
test("existing permission checks run before paging",async()=>{mocks.session.mockResolvedValue({organizations:["org-a"],permissions:[]});const html=renderToStaticMarkup(await FactoryPage({searchParams:Promise.resolve({units_after:"cursor"})}));expect(html).toContain("Acceso denegado");expect(mocks.get).not.toHaveBeenCalled();});

const appointmentDates = [
 {locale:"es-AR",language:"es",timeZone:"America/Argentina/Buenos_Aires",date:"2035-01-02T01:30:00Z",dateLabel:"1 ene",hour:"10:30 p. m."},
 {locale:"en-US",language:"en",timeZone:"America/New_York",date:"2035-01-02T01:30:00Z",dateLabel:"Jan 1",hour:"8:30 PM"},
 {locale:"en-US",language:"en",timeZone:"America/New_York",date:"2035-07-02T01:30:00Z",dateLabel:"Jul 1",hour:"9:30 PM"},
 {locale:"en-US",language:"en",timeZone:"Asia/Tokyo",date:"2035-01-02T01:30:00Z",dateLabel:"Jan 2",hour:"10:30 AM"},
];
test.each(appointmentDates)("both customer appointment views use configured $timeZone / $date",async(row)=>{
 mocks.locale.mockResolvedValue(row);
 mocks.get.mockImplementation(async(_s:unknown,p:string)=>p.endsWith("journey")?{appointments:[{id:"appointment",kind:"consultation",starts_at:row.date,state:"requested",version:1}],quotes:[],handovers:[]}:{items:[]});
 for(const view of [CustomerPage,AppointmentsPage]){
  const html=renderToStaticMarkup(await view()).replace(/\u00a0|\u202f/g," ");
  expect(html).toContain('dateTime="'+row.date+'"');expect(html).toContain(row.dateLabel);expect(html).toContain(row.hour);expect(html).toContain(row.timeZone);
 }
});
test("unverifiable appointment time fails instead of showing Invalid Date",async()=>{
 mocks.get.mockImplementation(async(_s:unknown,p:string)=>p.endsWith("journey")?{appointments:[{id:"appointment",kind:"consultation",starts_at:"not-a-date",state:"requested",version:1}],quotes:[],handovers:[]}:{items:[]});
 for(const view of [CustomerPage,AppointmentsPage]) await expect(view()).rejects.toThrow("Invalid appointment date");
});
test("appointment page denies access before loading time configuration or journey",async()=>{
 mocks.session.mockResolvedValue({organizations:["org-a"],permissions:[]});expect(renderToStaticMarkup(await AppointmentsPage())).toContain("Acceso denegado");expect(mocks.get).not.toHaveBeenCalled();expect(mocks.locale).not.toHaveBeenCalled();
});

import {createElement} from "react";
import {isCancellationReceipt,CustomerAppointmentActions} from "@/components/customer-appointment-actions";
const cancellationReceipt={id:"appointment",organization_id:"org-a",state:"cancelled",version:2};
test.each([
 {name:"exact receipt",value:cancellationReceipt,valid:true},
 {name:"null",value:null,valid:false},
 {name:"array",value:[],valid:false},
 {name:"different appointment",value:{...cancellationReceipt,id:"other"},valid:false},
 {name:"different organization",value:{...cancellationReceipt,organization_id:"org-b"},valid:false},
 {name:"unconfirmed state",value:{...cancellationReceipt,state:"requested"},valid:false},
 {name:"stale version",value:{...cancellationReceipt,version:1},valid:false},
 {name:"future version",value:{...cancellationReceipt,version:3},valid:false},
 {name:"string version",value:{...cancellationReceipt,version:"2"},valid:false},
 {name:"nonfinite version",value:{...cancellationReceipt,version:NaN},valid:false},
])("cancellation response binding: $name",({value,valid})=>{expect(isCancellationReceipt(value,"org-a",{id:"appointment",version:1})).toBe(valid)});
test("customer can reach appointment management from their account",async()=>{expect(renderToStaticMarkup(await CustomerPage())).toContain('href="/customer/appointments"');expect(renderToStaticMarkup(await CustomerPage())).toContain('Gestionar mis turnos')});
test("empty appointment view provides an explicit state",()=>{expect(renderToStaticMarkup(createElement(CustomerAppointmentActions,{organizationId:"org-a",appointments:[]}))).toContain("No hay turnos disponibles")});
````

### FILE: `docs/factory-operator-flow.md`

```yaml
block_id: "FACTORY-DELTA-TYPESCRIPT_GO_API_WEB_BRIDGE:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "5b9fb61ec11c58a5e1bb5d61216e773b8f00d9d0c98773e16063f435465f43c2"
variables: []
secrets_allowed: false
```

````markdown
# Factory operational progress

The factory portal now lets a role with both `factory:read` and `factory:write` record the existing `planned` → `assembly` transition for an authorized unit. Its source of truth remains `operations.Service` and the existing transactional PostgreSQL writer. The BFF only maps the named UI command to that owner; it does not add a transition graph, QC policy, inventory effect or production authorization.

`GET /v1/factory/units/{id}` reuses the factory list projection with exact tenant, purchase-order destination organization and unit predicates. The existing EnterpriseQueryModule mounts this read automatically; no additional host module is needed. The page derives organization and storage scope from the server session. The BFF validates permission, origin, strict body/query, streamed 4096-byte limit, scoped backend projection and expected state. The writer rechecks current state and scope inside its atomic state/outbox transaction, closing the read/write race.

The browser preserves an unresolved marker before POST and suppresses duplicate clicks and automatic mutation retries. Its recovery action performs only a GET. This is a current-state observation, not an idempotency receipt. Another actor may have advanced the unit after the lost response; the UI explicitly says the state does not identify the request that produced it. A successful HTTP response acknowledges the existing owner operation; an uncertain response stays uncertain until the operator reads the current state and decides a supported next action.

The focused browser gate uses an isolated PostgreSQL database, actual encrypted role sessions, RS256/JWKS bearer verification, production Next BFF and Go API. It starts with one planned fixture unit, loses the real assembly response after commit, performs a separate competing assembly→quality advancement through the same owner, reloads and observes quality without claiming request attribution. Read-only role, foreign organization, foreign tenant and stale command attempts fail. PostgreSQL contains one assembly event and one separate competing quality event; the new UI produced exactly one backend transition POST. No stock or money mutation is introduced.

Use the admitted runtime and dependency locks. Build with the direct Node Next CLI (`build --webpack` for the isolated external node_modules junction). Set `ELITE_FACTORY_BROWSER=1`, `ELITE_WEB_ROOT`, the absolute `ELITE_NODE_BIN`, and `FACTORY_BROWSER_DB_URL` for a fresh loopback database named `elite_payment_connected_*` with the composition migrations. The fixture has an exact environment/database guard and must not target business data. With module downloads disabled, invoke the admitted Go executable:

```text
go test -mod=readonly ./internal/platform/postgres -run ^TestFactoryBrowserPostgres$ -count=1 -v -timeout=4m
```

The fixture invokes the installed Playwright CLI directly, without package-manager execution or downloads. It requires the reference composition's `handoverBrowserIssuer` test helper, reused for local generated identity credentials; it does not rerun the handover gate. Chromium desktop plus a 390×844 viewport check passed with zero page errors and no horizontal overflow. This is not a four-browser claim. Omitted opt-in produces an explicit skip; acceptance must reject skips. The self-signed HTTPS exception remains restricted to the existing explicit loopback fixture configuration.

All delta code is AUTHORED read/transport/UI/test glue. The existing factory state graph is also AUTHORED and is not attributed to a vendor. No dependency, license or notice is added. This gate does not close general person management, training/evaluation, factory registration or all procurement/inventory UI capabilities.
````

### FILE: `microsoft_playwright_browser_gate/tests/factory-connected.spec.mjs`

```yaml
block_id: "FACTORY-DELTA-TYPESCRIPT_GO_API_WEB_BRIDGE:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "d456d1752f1eace33cbf5d34063d29e123bd272cb09b083b3f3dbf484dd52fdf"
variables: []
secrets_allowed: false
```

````text
import {test,expect} from '@playwright/test';
import {createHash} from 'node:crypto';
import {createRequire} from 'node:module';
import {resolve} from 'node:path';
import {pathToFileURL} from 'node:url';

test('factory operator records progress and observes concurrent state through real BFF API and PostgreSQL',async({page,context},info)=>{
 if(process.env.ELITE_FACTORY_BROWSER!=='1')throw new Error('explicit local factory fixture required');
 const base=process.env.ELITE_BASE_URL;expect(new URL(base).protocol).toBe('https:');expect(new URL(base).hostname).toBe('127.0.0.1');
 const identities=JSON.parse(process.env.ELITE_FACTORY_IDENTITIES);
 const require=createRequire(resolve(process.env.ELITE_WEB_ROOT,'package.json'));
 const {EncryptJWT}=await import(pathToFileURL(require.resolve('jose')).href);
 async function identity(name){
  const jwt=await new EncryptJWT(identities[name]).setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('300s').encrypt(createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest());
  await context.clearCookies();await context.addCookies([{name:'__Host-elite_session',value:jwt,url:base+'/',secure:true,httpOnly:true,sameSite:'Lax'}]);
 }

 const errors=[];page.on('pageerror',error=>errors.push(error.message));let pagePosts=0;page.on('request',request=>{if(request.method()==='POST'&&request.url().includes('/api/enterprise/factory'))pagePosts++});
 const body={organizationId:'store',unitId:'unit',action:'start-assembly',expectedState:'planned'};
 async function post(data=body){return context.request.post(base+'/api/enterprise/factory',{headers:{origin:base,'content-type':'application/json'},data})}
 await identity('reader');await page.goto('/factory');await expect(page.getByRole('button',{name:'Registrar inicio de ensamblaje'})).toHaveCount(0);expect((await post()).status()).toBe(403);
 await identity('foreign-org');await page.goto('/factory');await expect(page.getByText('SERIAL-FACTORY-FIXTURE',{exact:true})).toHaveCount(0);expect((await post()).status()).toBe(403);expect((await context.request.get(base+'/api/enterprise/factory?organizationId=other&unitId=unit')).status()).toBe(404);
 await identity('foreign-tenant');expect((await context.request.get(base+'/api/enterprise/factory?organizationId=store&unitId=unit')).status()).toBe(404);
 await identity('operator');await page.goto('/factory');const panel=page.getByRole('region',{name:'Seguimiento de SERIAL-FACTORY-FIXTURE'});await expect(panel.getByText('planned',{exact:true})).toBeVisible();
 await page.route('**/api/enterprise/factory',async route=>{if(route.request().method()!=='POST'){await route.continue();return};const response=await route.fetch();expect(response.status()).toBe(200);await route.abort('failed')});
 await panel.getByRole('button',{name:'Registrar inicio de ensamblaje'}).evaluate(button=>{button.click();button.click()});await expect(panel.getByText('La respuesta no quedó confirmada.',{exact:false})).toBeVisible();expect(pagePosts).toBe(1);
 // A separate actor advances the same unit before recovery: state alone cannot attribute the lost request.
 expect((await context.request.post(process.env.ELITE_FACTORY_CONTROL,{headers:{'X-Fixture-Token':process.env.ELITE_FACTORY_CONTROL_TOKEN}})).status()).toBe(200);
 await page.unrouteAll();await page.reload();await expect(panel.getByText('Hay una respuesta pendiente de comprobar.')).toBeVisible();await panel.getByRole('button',{name:'Consultar estado actual'}).click();await expect(panel.getByText('quality',{exact:true})).toBeVisible();await expect(panel.getByText('Estado actual consultado. No identifica la solicitud que lo produjo.')).toBeVisible();expect(pagePosts).toBe(1);
 expect((await post()).status()).toBe(409);await expect(panel.getByRole('button',{name:'Registrar inicio de ensamblaje'})).toHaveCount(0);
 const observed=await context.request.get(base+'/api/enterprise/factory?organizationId=store&unitId=unit');expect(observed.status()).toBe(200);expect(await observed.json()).toMatchObject({unit:{state:'quality'},observationOnly:true});
 await page.setViewportSize({width:390,height:844});expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth)).toBe(true);await page.screenshot({path:info.outputPath('factory-recovered-mobile.png'),fullPage:true});expect(errors).toEqual([]);
});
````

### FILE: `src/app/api/enterprise/factory/route.test.ts`

```yaml
block_id: "FACTORY-DELTA-TYPESCRIPT_GO_API_WEB_BRIDGE:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "7f581709cea725507924d04870e5d81fef0a6339b08531f03fbb648aac66b65f"
variables: []
secrets_allowed: false
```

````typescript
import{beforeEach,expect,it,vi}from"vitest";
const m=vi.hoisted(()=>({session:vi.fn(),get:vi.fn(),post:vi.fn()}));
vi.mock("@/platform/auth/session",()=>({readSession:m.session,allowed:(s:{permissions:string[]},p:string)=>s.permissions.includes(p)}));
vi.mock("@/platform/auth/oidc-client",()=>({applicationBaseUrl:()=>new URL("https://portal.example.test")}));
vi.mock("@/platform/backend/protected-client",()=>({protectedGet:m.get,protectedPost:m.post}));
import{GET,POST}from"./route";
const session={subject:"operator",tenantId:"tenant",organizations:["store"],permissions:["factory:read","factory:write"],accessToken:"private-token"};
const unit={id:"unit",organization_id:"store",purchase_order_id:"po",variant_id:"v",serial_number:"SERIAL",state:"planned"};
const url="https://portal.example.test/api/enterprise/factory";
const request=(extra:Record<string,unknown>={})=>new Request(url,{method:"POST",headers:{origin:"https://portal.example.test","content-type":"application/json"},body:JSON.stringify({organizationId:"store",unitId:"unit",action:"start-assembly",expectedState:"planned",...extra})});
beforeEach(()=>{vi.clearAllMocks();m.session.mockResolvedValue(session);m.get.mockResolvedValue(unit);m.post.mockResolvedValue({status:"accepted"})});
it("forwards only the existing transition after scoped read, without inventing a receipt",async()=>{const r=await POST(request());expect(r.status).toBe(200);expect(m.post).toHaveBeenCalledWith(session,"/v1/factory/units/unit/transitions",{organization_id:"store",current:"planned",target:"assembly"});expect(await r.json()).toEqual({status:"accepted"});expect(m.post).toHaveBeenCalledTimes(1)});
it.each(["missing","readonly","no-read","foreign-org","origin"])("rejects %s before downstream",async mode=>{let req=request();if(mode==="missing")m.session.mockResolvedValue(null);if(mode==="readonly")m.session.mockResolvedValue({...session,permissions:["factory:read"]});if(mode==="no-read")m.session.mockResolvedValue({...session,permissions:["factory:write"]});if(mode==="foreign-org")req=request({organizationId:"other"});if(mode==="origin")req=new Request(req,{headers:{origin:"https://evil.test","content-type":"application/json"}});expect([401,403]).toContain((await POST(req)).status);expect(m.get).not.toHaveBeenCalled();expect(m.post).not.toHaveBeenCalled()});
it.each([{target:"released"},{expectedState:"assembly"},{unitId:"../foreign"},{tenantId:"foreign"},{requestReceipt:"invented"}])("rejects unsupported body %j",async body=>{expect((await POST(request(body))).status).toBe(400);expect(m.post).not.toHaveBeenCalled()});
it("rejects stale observation and scoped response mismatch",async()=>{m.get.mockResolvedValue({...unit,state:"assembly"});expect((await POST(request())).status).toBe(409);m.get.mockResolvedValue({...unit,organization_id:"other"});expect((await POST(request())).status).toBe(502);expect(m.post).not.toHaveBeenCalled()});
it("current state is read-only and does not attribute actor or request",async()=>{m.get.mockResolvedValue({...unit,state:"quality",private:"not-returned"});const r=await GET(new Request(url+"?organizationId=store&unitId=unit"));expect(r.status).toBe(200);expect(await r.json()).toEqual({unit:{...unit,state:"quality"},observationOnly:true});expect(m.post).not.toHaveBeenCalled()});
it.each(["organizationId=other&unitId=unit","organizationId=store&unitId=unit&unitId=unit","organizationId=store&unitId=unit&extra=x"])("rejects invalid scope query %s",async q=>{expect([400,403]).toContain((await GET(new Request(url+"?"+q))).status);expect(m.get).not.toHaveBeenCalled()});
it("rejects oversized chunked body before downstream",async()=>{let cancelled=false;const stream=new ReadableStream<Uint8Array>({start(controller){controller.enqueue(new TextEncoder().encode("x".repeat(4097)))},cancel(){cancelled=true}});const req=new Request(url,{method:"POST",headers:{origin:"https://portal.example.test","content-type":"application/json"},body:stream,duplex:"half"}as RequestInit);expect((await POST(req)).status).toBe(413);expect(cancelled).toBe(true);expect(m.get).not.toHaveBeenCalled();expect(m.post).not.toHaveBeenCalled()});
it("hides backend details and does not automatically retry",async()=>{m.post.mockRejectedValue(new Error("PRIVATE"));const r=await POST(request());expect(r.status).toBe(503);expect(await r.text()).not.toContain("PRIVATE");expect(m.post).toHaveBeenCalledTimes(1)});
````

### FILE: `src/app/api/enterprise/factory/route.ts`

```yaml
block_id: "FACTORY-DELTA-TYPESCRIPT_GO_API_WEB_BRIDGE:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "7f3d5a6be9739efcfc2b65985428cf89b8de5acc621aba71cda21834205c8b36"
variables: []
secrets_allowed: false
```

````typescript
// AUTHORED BFF glue. Existing operations.Service remains the transition authority.
import {NextResponse}from"next/server";
import{z}from"zod";
import{allowed,readSession}from"@/platform/auth/session";
import{applicationBaseUrl}from"@/platform/auth/oidc-client";
import{protectedGet,protectedPost}from"@/platform/backend/protected-client";
import{BackendProblem}from"@/platform/backend/public-client";
import{factoryID,factoryUnitSchema}from"@/platform/factory/contracts";
const scope={organizationId:factoryID,unitId:factoryID};
const query=z.object(scope).strict();
const command=z.object({...scope,action:z.literal("start-assembly"),expectedState:z.literal("planned")}).strict();
const reply=(body:unknown,status=200)=>NextResponse.json(body,{status,headers:{"Cache-Control":"no-store","X-Content-Type-Options":"nosniff"}});
function problem(e:unknown){return reply({code:"FACTORY_UNAVAILABLE"},e instanceof BackendProblem&&[400,401,403,404,409,503].includes(e.status)?e.status:503)}
async function boundedBody(request: Request) {
 if (!request.body) throw new Error("invalid factory body");
 const reader = request.body.getReader(); const chunks: Uint8Array[] = []; let bytes = 0; let timer: ReturnType<typeof setTimeout> | undefined;
 const deadline = new Promise<never>((_, reject) => { timer = setTimeout(() => reject(new Error("body deadline")), 2500); });
 try {
  while (true) {
   const part = await Promise.race([reader.read(), deadline]); if (part.done) break;
   bytes += part.value.byteLength; if (bytes > 4096) throw new Error("BODY_TOO_LARGE"); chunks.push(part.value);
  }
  const data = new Uint8Array(bytes); let offset = 0; for (const chunk of chunks) { data.set(chunk, offset); offset += chunk.byteLength; }
  return new TextDecoder("utf-8", { fatal: true }).decode(data);
 } finally { clearTimeout(timer); void reader.cancel().catch(() => {}); reader.releaseLock(); }
}

export async function GET(request:Request){
 if(["cross-site","none"].includes(request.headers.get("sec-fetch-site")??""))return reply({code:"CROSS_ORIGIN_REJECTED"},403);
 const session=await readSession();if(!session)return reply({code:"UNAUTHENTICATED"},401);if(!allowed(session,"factory:read"))return reply({code:"FORBIDDEN"},403);
 const search=new URL(request.url).searchParams;if(request.url.length>2048||[...search.keys()].some(k=>search.getAll(k).length!==1))return reply({code:"INVALID_QUERY"},400);
 const parsed=query.safeParse(Object.fromEntries(search));if(!parsed.success)return reply({code:"INVALID_QUERY"},400);const q=parsed.data;if(!session.organizations.includes(q.organizationId))return reply({code:"FORBIDDEN"},403);
 try{const unit=factoryUnitSchema.parse(await protectedGet<unknown>(session,`/v1/factory/units/${encodeURIComponent(q.unitId)}`,{organization_id:q.organizationId}));if(unit.id!==q.unitId||unit.organization_id!==q.organizationId)return reply({code:"CONTEXT_MISMATCH"},502);return reply({unit,observationOnly:true})}catch(e){return problem(e)}
}
export async function POST(request:Request){
 try{if(request.headers.get("origin")!==applicationBaseUrl().origin)return reply({code:"CROSS_ORIGIN_REJECTED"},403)}catch{return reply({code:"CROSS_ORIGIN_REJECTED"},403)}
 if(request.headers.get("content-type")!=="application/json")return reply({code:"UNSUPPORTED_MEDIA_TYPE"},415);
 const session=await readSession();if(!session)return reply({code:"UNAUTHENTICATED"},401);if(!allowed(session,"factory:read")||!allowed(session,"factory:write"))return reply({code:"FORBIDDEN"},403);
 let body:unknown;try{body=JSON.parse(await boundedBody(request))}catch(e){return reply({code:e instanceof Error&&e.message==="BODY_TOO_LARGE"?"COMMAND_TOO_LARGE":"INVALID_COMMAND"},e instanceof Error&&e.message==="BODY_TOO_LARGE"?413:400)}
 const parsed=command.safeParse(body);if(!parsed.success)return reply({code:"INVALID_COMMAND"},400);const q=parsed.data;if(!session.organizations.includes(q.organizationId))return reply({code:"FORBIDDEN"},403);
 try{const unit=factoryUnitSchema.parse(await protectedGet<unknown>(session,`/v1/factory/units/${encodeURIComponent(q.unitId)}`,{organization_id:q.organizationId}));if(unit.id!==q.unitId||unit.organization_id!==q.organizationId)return reply({code:"CONTEXT_MISMATCH"},502);if(unit.state!==q.expectedState)return reply({code:"STATE_CHANGED"},409);
 const ack=z.object({status:z.literal("accepted")}).strict().parse(await protectedPost<unknown>(session,`/v1/factory/units/${encodeURIComponent(q.unitId)}/transitions`,{organization_id:q.organizationId,current:q.expectedState,target:"assembly"}));return reply(ack)
 }catch(e){return problem(e)}
}
````

### FILE: `src/components/factory-unit-actions.tsx`

```yaml
block_id: "FACTORY-DELTA-TYPESCRIPT_GO_API_WEB_BRIDGE:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "a58e5e091eb784f58c78d21f7b9b3c1e4664f73f07184a94558086d8c2dfaa32"
variables: []
secrets_allowed: false
```

````tsx
"use client";
// AUTHORED UI glue; an observed state never claims attribution to a lost request.
import{useEffect,useRef,useState}from"react";
import{factoryUnitSchema,type FactoryUnit}from"@/platform/factory/contracts";
export function FactoryUnitActions({unit,canWrite,storageScope}:{unit:FactoryUnit;canWrite:boolean;storageScope:string}){
 const[value,setValue]=useState(unit),[pending,setPending]=useState(false),[message,setMessage]=useState(""),[busy,setBusy]=useState(false);const active=useRef(false);
 const key=`elite:factory:${storageScope}:${unit.id}`;
 useEffect(()=>{try{setPending(localStorage.getItem(key)==="unknown")}catch{setPending(true);setMessage("No se pudo consultar el registro local. Consultá el estado antes de avanzar.")}},[key]);
 async function act(write:boolean){if(active.current)return;active.current=true;setBusy(true);setMessage("");
 try{
  if(write){localStorage.setItem(key,"unknown");setPending(true)}
  const q=new URLSearchParams({organizationId:unit.organization_id,unitId:unit.id});
  const response=await fetch(`/api/enterprise/factory${write?"":"?"+q}`,write?{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify({organizationId:unit.organization_id,unitId:unit.id,action:"start-assembly",expectedState:"planned"})}:{cache:"no-store"});
  if(!response.ok)throw new Error("unconfirmed");const body:unknown=await response.json();
  if(write){if(!(typeof body==="object"&&body!==null&&"status"in body&&body.status==="accepted"))throw new Error("invalid response");setMessage("El servidor aceptó el avance. Consultá el estado actual.")}
  else{const parsed=factoryUnitSchema.parse((body as {unit?:unknown}).unit);if(parsed.id!==unit.id||parsed.organization_id!==unit.organization_id)throw new Error("mismatch");setValue(parsed);localStorage.removeItem(key);setPending(false);setMessage("Estado actual consultado. No identifica la solicitud que lo produjo.")}
 }catch{setMessage(write?"La respuesta no quedó confirmada. Consultá el estado antes de volver a operar.":"No se pudo consultar el estado actual.")}finally{active.current=false;setBusy(false)}}
 return <section aria-label={`Seguimiento de ${unit.serial_number}`}><p>Estado actual: <strong>{value.state}</strong></p>{pending?<p>Hay una respuesta pendiente de comprobar.</p>:null}
 {canWrite&&value.state==="planned"?<button className="button" type="button" disabled={busy||pending} onClick={()=>void act(true)}>Registrar inicio de ensamblaje</button>:null}
 <button className="button" type="button" disabled={busy} onClick={()=>void act(false)}>Consultar estado actual</button>
 {message?<p role="status" aria-live="polite">{message}</p>:null}</section>
}
````

### FILE: `src/platform/factory/contracts.ts`

```yaml
block_id: "FACTORY-DELTA-TYPESCRIPT_GO_API_WEB_BRIDGE:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "5c8bba603a068e010b5eb280fe60043a6b012b31caeae547006a40788854306e"
variables: []
secrets_allowed: false
```

````typescript
// AUTHORED transport contract; state names mirror the existing operations owner.
import { z } from "zod";
export const factoryID=z.string().regex(/^[A-Za-z0-9][A-Za-z0-9._:-]{0,199}$/);
export const factoryUnitSchema=z.object({id:factoryID,organization_id:factoryID,purchase_order_id:factoryID,variant_id:factoryID,serial_number:z.string().min(1).max(200),state:z.enum(["planned","assembly","quality","released","shipped","received","rejected"])});
export type FactoryUnit=z.infer<typeof factoryUnitSchema>;
````

### FILE: `src/app/api/enterprise/franchise/whatsapp/replies/route.test.ts`

```yaml
block_id: "COMMUNICATIONS-DELTA-TYPESCRIPT_GO_API_WEB_BRIDGE:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "350ddc206345927ac6752a0d453a639e3ae4fa4cd7cf86953c098814e93c5ae4"
variables: []
secrets_allowed: false
```

````typescript
import { beforeEach,describe,it,expect,vi } from "vitest";
import { NextRequest } from "next/server";
const mocks=vi.hoisted(()=>({session:{subject:"human",tenantId:"tenant",organizations:["org"],permissions:["whatsapp:approve","whatsapp:send"],accessToken:"private-server-only"},post:vi.fn()}));
vi.mock("@/platform/auth/session",()=>({readSession:async()=>mocks.session,allowed:(s:{permissions:string[]},p:string)=>s.permissions.includes(p)}));
vi.mock("@/platform/auth/oidc-client",()=>({applicationBaseUrl:()=>new URL("https://portal.example")}));
vi.mock("@/platform/backend/protected-client",()=>({protectedPost:mocks.post}));
import { POST } from "./route";
const id="wa-reply:"+"a".repeat(64),sha="b".repeat(64);
function request(body:unknown,origin="https://portal.example"){return new NextRequest("https://portal.example/api/enterprise/franchise/whatsapp/replies",{method:"POST",headers:{origin,"content-type":"application/json"},body:JSON.stringify(body)})}
beforeEach(()=>{mocks.post.mockReset();mocks.session.permissions=["whatsapp:approve","whatsapp:send"]});
describe("WhatsApp review BFF",()=>{
 it("rejects foreign origin before backend or effects",async()=>{expect((await POST(request({action:"send",request_id:id,payload_sha256:sha},"https://foreign.example"))).status).toBe(403);expect(mocks.post).not.toHaveBeenCalled()});
 it("requires send permission and cannot substitute editable text",async()=>{mocks.session.permissions=["whatsapp:approve"];expect((await POST(request({action:"send",request_id:id,payload_sha256:sha}))).status).toBe(403);expect((await POST(request({action:"decision",request_id:id,payload_sha256:sha,approved:true,reason:"Reviewed",text:"Replacement"}))).status).toBe(400);expect(mocks.post).not.toHaveBeenCalled()});
 it("uses only the exact stored identity/hash, and uncertainty directs a read",async()=>{mocks.post.mockRejectedValue(new Error("uncertain"));const result=await POST(request({action:"recover",request_id:id,payload_sha256:sha}));expect(result.status).toBe(409);expect(mocks.post).toHaveBeenCalledWith(mocks.session,`/v1/franchise/whatsapp/replies/${encodeURIComponent(id)}/recover`,{payload_sha256:sha});expect(await result.json()).toEqual({code:"CONSULT_STATE_BEFORE_RETRY"})});
});
````

### FILE: `src/app/api/enterprise/franchise/whatsapp/replies/route.ts`

```yaml
block_id: "COMMUNICATIONS-DELTA-TYPESCRIPT_GO_API_WEB_BRIDGE:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "63d81a322738e0ec378e6dc5b07237625390ee57e35280b55624468a4191cce5"
variables: []
secrets_allowed: false
```

````typescript
import { NextResponse,type NextRequest } from "next/server";
import { allowed,readSession } from "@/platform/auth/session";
import { applicationBaseUrl } from "@/platform/auth/oidc-client";
import { protectedPost } from "@/platform/backend/protected-client";
import { replyCommand } from "@/platform/notifications/reply-contract";
export async function POST(request:NextRequest){
 const response=(code:string,status:number)=>NextResponse.json({code},{status,headers:{"cache-control":"no-store"}});
 let expected:string;try{expected=applicationBaseUrl().origin}catch{return response("UNAVAILABLE",503)}
 if(request.headers.get("origin")!==expected)return response("CROSS_ORIGIN_REJECTED",403);
 if(request.headers.get("content-type")!=="application/json")return response("UNSUPPORTED_MEDIA_TYPE",415);
 const text=await request.text();if(text.length>4096)return response("COMMAND_TOO_LARGE",413);
 let input:unknown;try{input=JSON.parse(text)}catch{return response("INVALID_COMMAND",400)}
 const parsed=replyCommand.safeParse(input);if(!parsed.success)return response("INVALID_COMMAND",400);
 const session=await readSession();if(!session)return response("UNAUTHENTICATED",401);
 const command=parsed.data;if(!allowed(session,"whatsapp:approve")||(command.action!=="decision"&&!allowed(session,"whatsapp:send")))return response("FORBIDDEN",403);
 const body=command.action==="decision"?{payload_sha256:command.payload_sha256,approved:command.approved,reason:command.reason}:{payload_sha256:command.payload_sha256};
 try{const result=await protectedPost<unknown>(session,`/v1/franchise/whatsapp/replies/${encodeURIComponent(command.request_id)}/${command.action}`,body);return NextResponse.json(result,{headers:{"cache-control":"no-store"}})}catch{return response("CONSULT_STATE_BEFORE_RETRY",409)}
}
````

### FILE: `src/platform/notifications/reply-contract.test.ts`

```yaml
block_id: "COMMUNICATIONS-DELTA-TYPESCRIPT_GO_API_WEB_BRIDGE:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "de11c28f90432603b1b396724ad225a1acb793b81c59e307d060440c72bc337f"
variables: []
secrets_allowed: false
```

````typescript
import { describe,it,expect } from "vitest";
import { replyCommand,responseText,type ReplyView } from "./reply-contract";
const id="wa-reply:"+"a".repeat(64),sha="b".repeat(64);
describe("human WhatsApp response contract",()=>{
 it("binds approval to exact stored identity/hash and excludes editable text",()=>{expect(replyCommand.safeParse({action:"decision",request_id:id,payload_sha256:sha,approved:true,reason:"Reviewed"}).success).toBe(true);expect(replyCommand.safeParse({action:"decision",request_id:id,payload_sha256:sha,approved:true,reason:"Reviewed",text:"Replacement"}).success).toBe(false)});
 it("rejects null decision and path escape",()=>{expect(replyCommand.safeParse({action:"decision",request_id:id,payload_sha256:sha,approved:null,reason:"Reviewed"}).success).toBe(false);expect(replyCommand.safeParse({action:"send",request_id:"../another",payload_sha256:sha}).success).toBe(false)});
 it("recovery has no provider ID or override payload",()=>{expect(replyCommand.safeParse({action:"recover",request_id:id,payload_sha256:sha}).success).toBe(true);expect(replyCommand.safeParse({action:"recover",request_id:id,payload_sha256:sha,provider_message_id:"forged"}).success).toBe(false)});
 it("displays exact response and rejects recipient divergence",()=>{const value={context:{message:{ExternalID:"5491112345678",Text:JSON.stringify({kind:"text_reply",recipient:"5491112345678",text:"Texto <exacto>"})}}} as ReplyView;expect(responseText(value)).toBe("Texto <exacto>");value.context.message.ExternalID="other";expect(()=>responseText(value)).toThrow("RECIPIENT_MISMATCH")});
});
````

### FILE: `src/platform/notifications/reply-contract.ts`

```yaml
block_id: "COMMUNICATIONS-DELTA-TYPESCRIPT_GO_API_WEB_BRIDGE:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "c691336e44b02403a910d9d6c30200e95b866214083ca8c50616bf4280ce4bc5"
variables: []
secrets_allowed: false
```

````typescript
import { z } from "zod";
// AUTHORED UI/API projection; the Go owner makes every approval/send decision.
const sha = z.string().regex(/^[0-9a-f]{64}$/);
export const replyCommand = z.discriminatedUnion("action", [
 z.object({action:z.literal("decision"),request_id:z.string().regex(/^wa-reply:[0-9a-f]{64}$/),payload_sha256:sha,approved:z.boolean(),reason:z.string().max(2048)}).strict(),
 z.object({action:z.enum(["send","recover"]),request_id:z.string().regex(/^wa-reply:[0-9a-f]{64}$/),payload_sha256:sha}).strict(),
]);
export type ReplyView = { request_id:string; state:"pending"|"approved"|"rejected"; payload_sha256:string; context:{organization_id:string;expires_at:string;message:{Text:string;ExternalID:string}};status?:{fence_state:string;delivery_status:string;reconciliation_required:boolean;approval_expired:boolean} };
export function responseText(value:ReplyView):string {
 const request:unknown=JSON.parse(value.context.message.Text);
 const parsed=z.object({kind:z.literal("text_reply"),recipient:z.string(),text:z.string().min(1).max(4096)}).passthrough().parse(request);
 if(parsed.recipient!==value.context.message.ExternalID)throw new Error("RECIPIENT_MISMATCH");
 return parsed.text;
}
````

## 6. Configuration surface

| Variable/file | Type/default | Secret | Validation/effect |
|---|---|---|---|
| `BUSINESS_CONFIG_FILE` | basename / `business.example.json` | no | selected only inside `config/`; Zod schema and cross-reference checks fail closed |
| `APP_BASE_URL` | absolute HTTPS production URL / none | no | OIDC callback/origin input supplied to the overlay |
| `OIDC_ISSUER`, `OIDC_CLIENT_ID` | issuer URL/non-empty / none | no | consumed by the OIDC overlay |
| `AUTH_SESSION_SECRET`, `OIDC_CLIENT_SECRET` | external values / none | yes | never committed; consumed server-only by the OIDC overlay |
| `ENTERPRISE_API_BASE_URL` | safe HTTP(S) server URL / none | no | fixed public and protected upstream paths only |
| `ENTERPRISE_TENANT_CODE`, `ENTERPRISE_ORGANIZATION_CODE` | lowercase bounded code / none | no | encoded into admitted public API paths |
| `config/business.example.json` | schema `1.0.0` | no | business label, markets, roles, modules, workflows, custom fields, integrations and feature flags |

PUBLIC_SITE_ORIGIN and PUBLIC_INDEXING_ENABLED are non-secret server deployment configuration. Origin must be an explicit HTTPS origin; indexing defaults disabled and requires flag1 plus origin. They are independent of APP_BASE_URL and never derived from request headers.

## 7. Dependency bill

| Dependency | Pin | Use | License | Scope | Source |
|---|---|---|---|---|---|
| Node.js / pnpm | `>=24` / `11.25.0` | build and BFF runtime | project terms/MIT | build/runtime | official projects |
| Next.js / React | `16.3.4` / `19.2.8` | web and BFF | MIT | runtime | official projects |
| TypeScript | `7.0.2` | strict compile-time checks | Apache-2.0 | build | official project |
| Zod | `4.4.3` | request and business-config validation | MIT | runtime | `github.com/colinhacks/zod` |
| Vitest | `4.1.11` | unit tests | MIT | test | official project |
| Google SafeValues | `1.2.0` | typed XSS-safe builders and DOM wrappers | Apache-2.0 | runtime | exact official npm package |
| Google CSP Evaluator | `1.1.8` | fail-closed CSP analysis | Apache-2.0 | test | exact official npm package; upstream says not an official Google product |
| `jose` / `openid-client` / `server-only` | `6.2.10` / `6.8.5` / `0.0.1` | exact dependency surface used by OIDC overlay | MIT | runtime | official packages |

The lockfile is part of the pack. `THIRD_PARTY_NOTICES.md` is a baseline; the composed project must regenerate the effective production inventory with the license adapter and review the distributed artifact.

## 8. Apply order

Select this pack as the web foundation, then compose `TS-OIDC-PORTAL-ADAPTER 0.2.x`, the journey portal and the license adapter. Inject target configuration and secrets outside source control. Run frozen offline-capable install, typecheck, tests and production build before connecting the admitted Go public/query APIs. Existing frontends require an explicit route, cookie, DNS and rollback migration; never overlay this pack on `TS-ENTERPRISE-WEB` implicitly.

## 9. Verification

Clean composition must pass `pnpm install --frozen-lockfile --offline`, `pnpm typecheck`, `pnpm test` and `pnpm build`. Unit suites cover public client path/idempotency/content-type and URL safety, business-config invariants, HTTP problem mapping, encrypted session behavior, protected-client constraints and Google CSP evaluation. Microsoft Playwright must additionally prove a new response nonce and exact nonce propagation to every rendered script across four projects. Production still requires edge/CDN repetition, browser E2E for public lead plus every role, selected IdP sandbox, API contract tests, distributed antiabuse, target CSRF/provider allowlists, accessibility with browser/AT, representative performance and failure tests, observability, rollout and rollback evidence.

## 10. Reconstruction evidence

Version `0.5.0` was reconstructed with `TS-OIDC-PORTAL-ADAPTER 0.2.1` and `TS-FRANCHISE-JOURNEY-PORTALS 0.4.0` on 2026-08-30. Frozen offline install downloaded zero packages; strict typecheck, nine Vitest files/32 integrated tests and the Next.js 16.3.2 production build passed. The single public client now reads bounded server-owned appointment slots; it does not add a calendar or backend owner. Playwright 4+8 and Lighthouse five-run evidence were repeated against the exact composition. See `FRANCHISE_APPOINTMENT_CAPACITY_2026-08-30_V117.md`; production remains conditioned by section 9.

## V321 — official package-manager security update

pnpm11.25.0 replaces vulnerable11.19.0 for this consumer; dependency lock bytes remain unchanged. Published bundle473packages/0advisories, registry ECDSA and publish/SLSA attestations verified, exact commit6d90c71efdffbc909b499490b64c66badc720327. Node24.20.0 standalone observed;115webPASS/1existingSKIP, build,92browser phases+agenda4,Playwright runtime4,Lighthouse policy5, local installation comparison. Eight changed materialized files across three owners; no upstream matcher or browser version changed. See reconstruction_evidence/PNPM_BUNDLE_SECURITY_V321.md. Native6files unchanged and transitive native SCA remains separate; product admission/monitor/target still blocked.

V379 adds the /help navigation label in es/en for the selected TS-FRANCHISE-JOURNEY-PORTALS0.14.12. Private operational guides remain Spanish. No locale, backend, dependency or action contract changes.

V381 integrates TS-MULTIROLE-ONBOARDING0.1.1 under opt-in role_workspace. BFF header filters exact effective session permissions and current feature/module flags using its canonical registry; guest customer entry remains. Browser test:workspace proves fifteen permission profiles across four browsers plus disabled flag, local307/session rejection, direct admin denial and help navigation. No domain handler, dependency lock or training policy changes. See ROLE_WORKSPACE_INTEGRATION_V381.md.

V382: exact selected-composition compatibility refreshed. Administrative lead reads are conditional on the existing lead grant; shared quote-format code corrects admin/customer minor-unit labels without repricing. TestAdministrativeReadBrowserPostgres uses real Go/JWKS/PostgreSQL and four browsers with synthetic identities; all production Go/SQL remains unchanged. The role implementation and quote acceptance/date/state logic retain exact source parity. See docs/private-portal-reads.md and reconstruction_evidence/PRIVATE_PORTAL_READS_V382.md. No full business, IdP, security, release or target admission.

V383: private read failures expose a fixed alert and explicit GET reconsultation in admin/customer/factory, without error details or command replay. Production build plus Go/JWKS/PostgreSQL and four browser projects prove persistent denial with the insufficient bearer and recovery after session correction; zero writes and unchanged durable snapshots.263web tests and canonical parity. Existing role/portal/Go business sources unchanged. See reconstruction_evidence/PRIVATE_READ_RECOVERY_V383.md; no full-control or target promotion.

V387: customer/factory cursor traversal connects three existing lists beyond25rows.81durable records per browser, four projects, reload/independent first-page recovery, zero writes and unchanged snapshots verified. Source is AUTHORED, business/API/SQL/dependencies unchanged. Supporting TEST02 progress, no whole-control or native security/release promotion. See reconstruction_evidence/PRIVATE_PORTAL_PAGINATION_V387.md.

V388: admin orders/cases/leads now traverse existing scoped keyset APIs using the shared navigator.276web tests and four browser projects exercise all six lists,163list records, no omissions/duplicates per list, permission removal, reload/first-page recovery and unchanged durable snapshots. All code remains AUTHORED/CONDITIONED; no full TEST02, native security or release promotion. See reconstruction_evidence/ADMIN_PORTAL_PAGINATION_V388.md.

V389: customer appointment times share trusted business locale/zone across SSR and interactive management, preserving original instants and cancellation behavior.282web tests plus eight connected browser/configuration runs, two appointment instants per view and seven-table unchanged snapshots. AUTHORED/CONDITIONED; no whole TEST02, security or release promotion. See reconstruction_evidence/CUSTOMER_APPOINTMENT_TIMEZONE_V389.md.

V390: customer cancellation receipt binding, synchronous single-submit and explicit GET recovery close the existing reference journey.294web tests;4actual cancellation browser projects,16unique durable cancellations/audits/outbox and28API writes;8read/time regression runs remain read-only. AUTHORED/CONDITIONED, no whole TEST02/security/release promotion. See reconstruction_evidence/CUSTOMER_CANCELLATION_RECOVERY_V390.md.

V393: the reference does not use runtime image transformation. Next image optimization is explicitly disabled and pnpm11.25.0 ignores the optional sharp dependency; frozen lock removes only its exclusive closure. This is dependency omission, not a repaired/admitted Sharp binary. Before enabling runtime image optimization, re-admit a compatible image pipeline, its exact dependencies, licenses/native provenance/security and web performance. Existing V386 vulnerability reproduction remains deferred. A connected real HTTP regression requires /_next/image to return404. No new version or business policy; AUTHORED/CONDITIONED. See reconstruction_evidence/OPTIONAL_IMAGE_DEPENDENCY_CONTAINMENT_V393.md.

V402 composed delta: Payment/initial-handover composition. New behavior and tests belong to the explicit runtime/portal packs; source and library release claims remain bounded to their evidence. Existing source provenance is preserved.

Canonical V402 integration: selected by the current profile with exact dependencies and caller overlays. Metadata promotion records byte reconstruction, not closure of every admission/release gate. Payload provenance is unchanged.

V402 composed delta: Connected handover browser/BFF/Go/PostgreSQL gate; bounded body, stable server date formatting, Next PageProps signatures. AUTHORED integration glue; existing domain and fixed upstreams unchanged.

V402 browser evidence: reconstruction_evidence/HANDOVER_BROWSER_V402.md. Production Webpack build includes TypeScript checking;42 direct-call and23 focused BFF/date tests PASS. Chromium desktop and mobile viewport through real BFF/API/PG cover lost-response recovery and callback invalidation. Hosted IdP/live payment/physical shipment not claimed.

V402 composed delta: Connected factory tracking browser/BFF/Go/PostgreSQL gate; exact scoped unit read and planned-to-assembly call to existing owner; recovery observes current state without attributing an uncertain effect. AUTHORED glue; existing domain and fixed upstreams unchanged.

V402 factory browser evidence: reconstruction_evidence/FACTORY_BROWSER_V402.md.18BFF tests, Next production build with TypeScript, Go vet and real browser/API/PG PASS; current-state recovery after another operator advances. No idempotent receipt claim or domain/writer change.

V402 composed delta: Typed WhatsApp BFF/contracts and external host/identity configuration. Exact FX source notice preserved. BFF uses existing authenticated permission and downstream error contracts.
