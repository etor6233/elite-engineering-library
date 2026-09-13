# Enterprise TypeScript Web BFF

## 1. Metadata

```yaml
pack_id: "TS-GO-API-WEB-BRIDGE"
pack_version: "0.15.2"
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
CREATE src/components/private-language-switcher.tsx
CREATE src/platform/backend/private-paging-links.tsx
CREATE src/platform/i18n/load-private-i18n.ts
CREATE src/platform/i18n/load-private-locale.ts
CREATE src/platform/i18n/private-catalog.ts
CREATE src/platform/i18n/private-locale.test.ts
CREATE src/platform/i18n/private-locale.ts
CREATE src/platform/i18n/private-messages.json
CREATE src/platform/i18n/private-provider.tsx
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
sha256: "c8372a0245b5822fbf2094c71374c12397ffa0fd49dd6102e2a260ce386b72ee"
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
// Narrow publication adapter shares this existing validated server configuration.
export function publishedCatalogTransportSettings() { return settings(); }
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
sha256: "d62b7497f655dfe5a271631093aedf61c68dbaf1b1309ddb6e710344c905a2b9"
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
      "locales": [
        "es-AR"
      ],
      "timeZone": "America/Argentina/Buenos_Aires",
      "taxMode": "external"
    }
  ],
  "organizationTypes": [
    {
      "id": "hq",
      "label": "Casa central",
      "allowedParents": []
    },
    {
      "id": "franchise",
      "label": "Franquicia",
      "allowedParents": [
        "hq"
      ]
    },
    {
      "id": "branch",
      "label": "Sucursal",
      "allowedParents": [
        "franchise",
        "hq"
      ]
    },
    {
      "id": "factory",
      "label": "Fábrica",
      "allowedParents": [
        "hq"
      ]
    },
    {
      "id": "supplier",
      "label": "Proveedor",
      "allowedParents": [
        "hq"
      ]
    }
  ],
  "roles": [
    {
      "id": "hq_admin",
      "label": "Administrador central",
      "permissions": [
        "*"
      ]
    },
    {
      "id": "branch_manager",
      "label": "Responsable de sucursal",
      "permissions": [
        "admin:read",
        "lead:read",
        "lead:assign",
        "lead:update",
        "quote:write",
        "catalog:write",
        "pricing:write",
        "order:create",
        "order:write",
        "inventory:allocate",
        "payment:create",
        "payment:write",
        "service:write",
        "communication:write"
      ]
    },
    {
      "id": "sales",
      "label": "Ventas",
      "permissions": [
        "admin:read",
        "lead:read",
        "lead:assign",
        "lead:update",
        "quote:write",
        "catalog:write",
        "order:create",
        "order:write",
        "payment:create"
      ]
    },
    {
      "id": "factory_operator",
      "label": "Operador de fábrica",
      "permissions": [
        "factory:read",
        "procurement:write",
        "factory:write",
        "inventory:write",
        "logistics:write"
      ]
    },
    {
      "id": "customer",
      "label": "Cliente",
      "permissions": [
        "customer:self"
      ]
    }
  ],
  "modules": {
    "catalog": {
      "enabled": true
    },
    "crm": {
      "enabled": true
    },
    "procurement": {
      "enabled": true
    },
    "inventory": {
      "enabled": true
    },
    "orders": {
      "enabled": true
    },
    "payments": {
      "enabled": true
    },
    "fulfillment": {
      "enabled": true
    },
    "service": {
      "enabled": true
    },
    "documents": {
      "enabled": true
    },
    "integrations": {
      "enabled": true
    }
  },
  "workflows": {
    "lead": {
      "initial": "new",
      "states": [
        "new",
        "contacted",
        "qualified",
        "converted",
        "lost"
      ],
      "transitions": [
        {
          "from": "new",
          "to": "contacted",
          "permission": "lead:update"
        },
        {
          "from": "contacted",
          "to": "qualified",
          "permission": "lead:update"
        },
        {
          "from": "qualified",
          "to": "converted",
          "permission": "lead:update"
        },
        {
          "from": "contacted",
          "to": "lost",
          "permission": "lead:update"
        },
        {
          "from": "qualified",
          "to": "lost",
          "permission": "lead:update"
        }
      ]
    },
    "order": {
      "initial": "draft",
      "states": [
        "draft",
        "placed",
        "confirmed",
        "paid",
        "allocated",
        "delivered",
        "cancelled"
      ],
      "transitions": [
        {
          "from": "draft",
          "to": "placed",
          "permission": "order:create"
        },
        {
          "from": "placed",
          "to": "confirmed",
          "permission": "order:transition"
        },
        {
          "from": "confirmed",
          "to": "paid",
          "permission": "payment:reconcile"
        },
        {
          "from": "paid",
          "to": "allocated",
          "permission": "inventory:reserve"
        },
        {
          "from": "allocated",
          "to": "delivered",
          "permission": "order:transition"
        },
        {
          "from": "draft",
          "to": "cancelled",
          "permission": "order:transition"
        },
        {
          "from": "placed",
          "to": "cancelled",
          "permission": "order:transition"
        }
      ]
    }
  },
  "customFields": {
    "lead": [
      {
        "id": "preferred_vehicle_use",
        "label": "Uso principal",
        "type": "select",
        "required": false,
        "options": [
          "urban",
          "delivery",
          "recreation",
          "fleet"
        ]
      }
    ],
    "catalog_model": [
      {
        "id": "estimated_range_km",
        "label": "Autonomía estimada (km)",
        "type": "number",
        "required": true
      }
    ]
  },
  "integrations": [
    {
      "id": "mercado_pago",
      "provider": "mercado_pago",
      "enabled": false,
      "mode": "sandbox",
      "capabilities": [
        "payments"
      ],
      "credentialRefEnv": "MERCADO_PAGO_CREDENTIAL_REF"
    },
    {
      "id": "amazon_sp_api",
      "provider": "amazon_sp_api",
      "enabled": false,
      "mode": "sandbox",
      "capabilities": [
        "catalog",
        "orders",
        "fulfillment"
      ],
      "credentialRefEnv": "AMAZON_SP_API_CREDENTIAL_REF"
    },
    {
      "id": "mercado_libre",
      "provider": "mercado_libre",
      "enabled": false,
      "mode": "sandbox",
      "capabilities": [
        "catalog",
        "orders",
        "fulfillment"
      ],
      "credentialRefEnv": "MERCADO_LIBRE_CREDENTIAL_REF"
    },
    {
      "id": "google_ads",
      "provider": "google_ads",
      "enabled": false,
      "mode": "sandbox",
      "capabilities": [
        "ads",
        "conversions"
      ],
      "credentialRefEnv": "GOOGLE_ADS_CREDENTIAL_REF"
    },
    {
      "id": "meta_ads",
      "provider": "meta_ads",
      "enabled": false,
      "mode": "sandbox",
      "capabilities": [
        "ads",
        "conversions"
      ],
      "credentialRefEnv": "META_ADS_CREDENTIAL_REF"
    }
  ],
  "features": {
    "public_catalog": true,
    "lead_capture": true,
    "customer_portal": true,
    "factory_portal": true,
    "vehicle_telemetry": false,
    "training_portal": false,
    "catalog_editor": false,
    "supply_portal": false,
    "warranty_portal": false,
    "network_portal": false
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
sha256: "311db84ef23a12a0c56191b21ee782224e782e71e2b25f04a77f435924d11a7a"
variables: []
secrets_allowed: false
```

````text
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  generateBuildId: async () => {
    const id = process.env.ELITE_SOURCE_SHA256;
    if (!id) return null;
    if (!/^[0-9a-f]{64}$/.test(id)) throw new Error("ELITE_SOURCE_SHA256 must bind the source inventory");
    return id;
  },
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
sha256: "9c7de890641f7098867d39e82f3cdc852c1627797e2e536270c851864b7a7f0e"
variables: []
secrets_allowed: false
```

````text
# Dependencias y licencias

Las dependencias directas están fijadas en `package.json` y el grafo completo en `pnpm-lock.yaml`.

| Dependencia directa | Versión | Licencia declarada por el paquete |
|---|---:|---|
| jose | 6.2.10 | MIT |
| Next.js | 16.3.4 | MIT |
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

## Odoo stored-value calculation

The optional gift-card/loyalty module adapts selected Odoo Community19.0 code at99edb6dd82b7b560930c00b03b694ba700785370 under LGPL-3.0-only. Complete license, copyright, source, local changes and replacement instructions are preserved in odoo_loyalty/ and docs/provenance/ODOO_STORED_VALUE_NOTICES.md. Go/SQL/portal glue is AUTHORED with its own declared provenance. No complete Odoo runtime or corporate authorship for that glue is claimed. The Next.js direct-version notice above is aligned to the existing16.3.4 package/lock; no dependency update was performed by this notice correction.

## Connected document reference

When the AWS Textract/document capability is selected (these dependencies are not implied in a web-only profile), AWS SDK for Go v2 Textract1.45.0 and its15fixed modules: Apache-2.0; complete original LICENSE/NOTICE and per-module hashes are in `aws_textract_runtime/licenses/selected-notices.json`. Preserve all four original texts on redistribution. The Microsoft public invoice fixture uses its exact MIT notice at `azure_document_intelligence_official_invoice/upstream/LICENSE.txt`. Original fixture bytes are acquired by fixed commit/hash, not included here. Document pipeline/test/configuration glue is AUTHORED locally and is not attributed to AWS, Microsoft or Google.
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
sha256: "9126740c5fc287683977736b1114fd5acbbc2bdeef1dc710ee343c9daebd1474"
variables: []
secrets_allowed: false
```

````text
import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
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
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

  const session = await readSession();
  if (!session) redirect("/api/auth/login?return_to=/admin" as Route);
  if (!allowed(session, "admin:read")) return <><h1 className="pageTitle">{t("p0006")}</h1><div className="notice">{t("p0007")} <code>admin:read</code>{t("p0008")}</div></>;
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
    <div className="eyebrow">{t("p0009")} {organization}</div><h1 className="pageTitle">{t("p0010")}</h1>
    <div className="grid">
      <article className="card"><h2>{t("p0011")}</h2><p>{overview.orders}</p></article>
      <article className="card"><h2>{t("p0012")}</h2><p>{overview.stock_available}</p></article>
      <article className="card"><h2>{t("p0013")}</h2><p>{overview.open_cases}</p></article>
      <article className="card"><h2>{t("p0014")}</h2><p>{overview.active_shipments}</p></article>
    </div>
    <section aria-label={t("p0015")}><h2>{t("p0015")}</h2><div className="grid">{orders.items.map((order) => <article className="card" key={order.id}><h3>{order.id}</h3><p>{order.state} {t("p0016")} {minorAmountPresentation(order.total_minor_units, order.currency, privateLocale.locale).amountLabel}</p></article>)}</div>
      {orders.items.length===0 ? <p>{t("p0017")}</p> : null}
      <PortalPaging label="Páginas de pedidos recientes" nextHref={orders.next_cursor ? portalPageHref("/admin",{orders_after:orders.next_cursor,cases_after:casesAfter,leads_after:leadsAfter}) : undefined} firstHref={ordersAfter ? portalPageHref("/admin",{cases_after:casesAfter,leads_after:leadsAfter}) : undefined}/></section>
    <section aria-label={t("p0018")}><h2>{t("p0018")}</h2><div className="grid">{cases.items.map((item) => <article className="card" key={item.id}><h3>{item.id}</h3><p>{item.state} {t("p0016")} {item.severity}</p></article>)}</div>
      {cases.items.length===0 ? <p>{t("p0017")}</p> : null}
      <PortalPaging label="Páginas de servicio" nextHref={cases.next_cursor ? portalPageHref("/admin",{cases_after:cases.next_cursor,orders_after:ordersAfter,leads_after:leadsAfter}) : undefined} firstHref={casesAfter ? portalPageHref("/admin",{orders_after:ordersAfter,leads_after:leadsAfter}) : undefined}/></section>
    {leads && <section aria-label={t("p0019")}><h2>{t("p0019")}</h2><div className="grid">{leads.items.map((item) => <article className="card" key={item.id}><h3>{item.id}</h3><p>{item.state} {t("p0016")} {item.source_code} {t("p0016")} {item.assigned_subject || t("p0020")}</p></article>)}</div>
      {leads.items.length===0 ? <p>{t("p0017")}</p> : null}
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
sha256: "b17fa71135c84b037519bfc4dd49641458d5e79a24ddbf13bbc7755e95890ca6"
variables: []
secrets_allowed: false
```

````text
import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import type { ReactElement } from "react";
import Link from "next/link";
import { CustomerCheckoutActions } from "@/components/customer-checkout-actions";
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
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

  const session = await readSession();
  if (!session) redirect("/api/auth/login?return_to=/customer" as Route);
  if (!allowed(session, "customer:self")) return <><h1 className="pageTitle">{t("p0006")}</h1><div className="notice">{t("p0007")} <code>customer:self</code>{t("p0008")}</div></>;
  const organization = session.organizations[0]!;
  const query = await searchParams ?? {};
  const ordersAfter = portalCursor(query.orders_after), casesAfter = portalCursor(query.cases_after);
  const [orders, cases, journey, locale] = await Promise.all([
    protectedGet<Page<OrderSummary>>(session, "/v1/customer/orders", { organization_id: organization, limit: "25", after: ordersAfter }),
    protectedGet<Page<ServiceCaseSummary>>(session, "/v1/customer/service-cases", { organization_id: organization, limit: "25", after: casesAfter }),
    protectedGet<CustomerJourney>(session, "/v1/customer/journey", { organization_id: organization }),
    Promise.resolve(privateLocale),
  ]);
  return <><div className="eyebrow">{t("p0045")}</div><h1 className="pageTitle">{t("p0046")}</h1>
    <section aria-label={t("p0047")}><h2>{t("p0047")}</h2><div className="grid">{orders.items.map((order) => <article className="card" key={order.id}><h3>{order.id}</h3><p>{order.state} {t("p0016")} {minorAmountPresentation(order.total_minor_units, order.currency, privateLocale.locale).amountLabel}</p>{["placed","confirmed","allocated"].includes(order.state)?<CustomerCheckoutActions orderId={order.id} organizationId={order.organization_id}/>:null}</article>)}</div>
      {orders.items.length === 0 ? <p>{t("p0017")}</p> : null}
<PortalPaging label="Páginas de mis pedidos" nextHref={orders.next_cursor ? portalPageHref("/customer", { orders_after: orders.next_cursor, cases_after: casesAfter }) : undefined} firstHref={ordersAfter ? portalPageHref("/customer", { cases_after: casesAfter }) : undefined} /></section>
    <section aria-label={t("p0048")}><h2>{t("p0048")}</h2><div className="grid">{cases.items.map((item) => <article className="card" key={item.id}><h3>{item.id}</h3><p>{item.state} {t("p0016")} {item.description}</p></article>)}</div>
      {cases.items.length === 0 ? <p>{t("p0017")}</p> : null}
<PortalPaging label="Páginas de mis casos de servicio" nextHref={cases.next_cursor ? portalPageHref("/customer", { cases_after: cases.next_cursor, orders_after: ordersAfter }) : undefined} firstHref={casesAfter ? portalPageHref("/customer", { orders_after: ordersAfter }) : undefined} /></section>
    <section><h2>{t("p0035")}</h2><p><Link href={"/customer/appointments" as Route}>{t("p0049")}</Link></p><div className="grid">{journey.appointments.map((item) => <article className="card" key={item.id}><h3>{item.kind}</h3><p>{item.state} {t("p0016")} <time dateTime={item.starts_at}>{publicAppointmentTime(locale, item.starts_at)} {t("p0050")}{locale.timeZone}{t("p0051")}</time></p></article>)}</div></section>
    <section><h2>{t("p0052")}</h2><div className="grid">{journey.quotes.map((item) => <article className="card" key={item.id}><h3>{item.id}</h3><p>{item.state} {t("p0016")} {minorAmountPresentation(item.total_minor_units, item.currency, privateLocale.locale).amountLabel}</p></article>)}</div></section>
    <section><h2>{t("p0039")}</h2><div className="grid">{journey.handovers.map((item) => <article className="card" key={item.id}><h3>{item.order_id}</h3><p>{item.state}</p></article>)}</div></section>
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
sha256: "855fa37ea641fd94ce184d5af38be365009170a8cfac63f1d4fab69e3b1c993e"
variables: []
secrets_allowed: false
```

````text
import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
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
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

  const session = await readSession();
  if (!session) redirect("/api/auth/login?return_to=/factory" as Route);
  if (!allowed(session, "factory:read")) return <><h1 className="pageTitle">{t("p0006")}</h1><div className="notice">{t("p0007")} <code>factory:read</code>{t("p0008")}</div></>;
  const organization = session.organizations[0]!;
  const storageScope=createHash("sha256").update(JSON.stringify([session.tenantId,session.subject,organization])).digest("hex");
  const query = await searchParams ?? {};
  const unitsAfter = portalCursor(query.units_after);
  const units = await protectedGet<Page<FactoryUnitSummary>>(session, "/v1/factory/units", { organization_id: organization, limit: "25", after: unitsAfter });
  return <><div className="eyebrow">{t("p0071")} {organization}</div><h1 className="pageTitle">{t("p0072")}</h1>
    <div className="grid">{units.items.map((unit) => <article className="card" key={unit.id}><h2>{unit.serial_number}</h2><p>{t("p0073")} {unit.purchase_order_id}</p><FactoryUnitActions unit={factoryUnitSchema.parse(unit)} canWrite={allowed(session,"factory:write")} storageScope={storageScope}/></article>)}</div>
    {units.items.length === 0 ? <p>{t("p0017")}</p> : null}
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
sha256: "e585bdfd5439d8f97af94afb7018cb2853f1322b93753e36f0b137a39bef974a"
variables: []
secrets_allowed: false
```

````text
import { publicMessage as message } from "@/platform/i18n/public-catalog";
import {PrivateLocaleProvider} from "@/platform/i18n/private-provider";
import {PrivateLanguageSwitcher} from "@/components/private-language-switcher";
import {loadPrivateLocale} from "@/platform/i18n/load-private-locale";
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
  const locale = await loadPrivateLocale();
  const session = await readSession();
  const navigation = enabledNavigation(config, session?.permissions ?? null);
  return (
    <html lang={locale.locale}>
      <body><PrivateLocaleProvider locale={locale}>
        <a className="skipLink" href="#main">{message(locale.locale, "site.skip")}</a>
        <header>
          <div className="shell headerRow">
            <Link className="brand" href="/">{config.business.name}</Link>
            <nav aria-label={message(locale.locale, "site.navigation")}>
              {navigation.map((item) => <Link key={item.id} href={item.href as Route}>{message(locale.locale, "nav." + item.id)}</Link>)}
            </nav>{session?<PrivateLanguageSwitcher/>:null}
          </div>
        </header>
        <main id="main" className="shell" lang={locale.locale}>{children}</main>
        <footer><div className="shell">{message(locale.locale, "site.contact")} <a href={`mailto:${config.business.supportEmail}`}>{config.business.supportEmail}</a></div></footer>
      </PrivateLocaleProvider></body>
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
sha256: "d715b58d7c1be02aae1051dc5d91d619499be5fd35972658cc67d01895c9d573"
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
import Link from "next/link";
import { loadPublishedCatalog } from "@/platform/catalog/load";

export async function generateMetadata() { const locale = await loadPublicLocale(); return publicPageMetadata("/models", message(locale.locale, "models.short")); }
export const dynamic = "force-dynamic";

export default async function ModelsPage() {
  const [catalog, locale] = await Promise.all([loadPublishedCatalog(), loadPublicLocale()]);
  // An enabled publication profile never falls back to mutable draft data.
  const models = catalog?.models ?? (process.env.CATALOG_RELEASE_ENABLED === "true" ? [] : await listModels());
  return (
    <section lang={locale.locale}>
      <div className="eyebrow">{message(locale.locale, "models.eyebrow")}</div>
      <h1 className="pageTitle">{message(locale.locale, "models.title")}</h1>
      <p className="lede">{message(locale.locale, "models.description")}</p>
      <p>{publicCount(locale.locale, "models", models.length)}</p>
      <div className="grid">
        {models.map((model) => (
          <article className="card" key={model.id}>
            <h2>{catalog ? <Link href={`/models/${model.code}`}>{model.displayName}</Link> : model.displayName}</h2>
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
sha256: "d6524e3e9dd1da55eff7862db9d3c2d2a32e0b70136b6b8dbc19a3806f70c665"
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
  {id:"help_content",label:"Contenido de ayuda",href:"/help/library",feature:"help_cms",anyPermissions:["help:read","help:write","help:publish"]},
  {id:"network",label:"Red de franquicia",href:"/network",feature:"network_portal",anyPermissions:["network:admin","franchise:write"]},
  {id:"warranty",label:"Garantía",href:"/warranty",feature:"warranty_portal",anyPermissions:["warranty:read","warranty:self","warranty:factory-read"]},
  {id:"supply",label:"Suministro",href:"/supply",feature:"supply_portal",anyPermissions:["supply:read","supply:factory-read"]},
  { id: "dashboard", label: "Mi panel", href: "/dashboard", feature: "role_workspace" },
  { id: "help", label: "Ayuda", href: "/help" },
  { id:"catalog_editor",label:"Editar catálogo",href:"/admin/catalog",feature:"catalog_editor",permission:"catalog:read" },
  { id: "training", label: "Capacitación", href: "/guide/training", feature: "training_portal", anyPermissions: ["training:learn","training:review"] },
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
sha256: "9ee34d3e4f99a5ef3ee81de3e14bcfd4bfdab0bf34625adbe73a4849128678da"
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

export function publicModelPath(code: string): string {
  if (!/^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$/u.test(code) || code.length > 128) throw new Error("Invalid public model code");
  return "/models/" + code;
}
export function publicModelMetadata(code: string, title: string, canonical: string): Metadata {
  const path = publicModelPath(code), settings = readPublicIndexing(), uri = new URL(canonical);
  if (uri.protocol !== "https:" || uri.username || uri.password || uri.search || uri.hash ||
      uri.pathname !== path || (settings.origin !== null && uri.origin !== settings.origin)) throw new Error("Invalid public model canonical");
  return { title, robots: { index: settings.enabled, follow: settings.enabled }, alternates: { canonical: uri.href } };
}
function approvedPublicPaths(modelCodes: readonly string[]) {
  if (modelCodes.length > 32) throw new Error("Too many published model paths");
  return [...publicPaths, ...Array.from(new Set(modelCodes.map(publicModelPath))).sort()];
}
export function publicRobots(settings = readPublicIndexing(), modelCodes: readonly string[] = []): MetadataRoute.Robots {
  if (!settings.enabled || !settings.origin) return { rules: { userAgent: "*", disallow: "/" } };
  // Exact public-path matches: adding a portal does not expose it to crawlers.
  // This is a crawl instruction, not access control or removal from an index.
  return { rules: { userAgent: "*", disallow: "/", allow: [...approvedPublicPaths(modelCodes).map(path => path + "$"), "/_next/static/", "/icon.svg$"] },
    sitemap: new URL("/sitemap.xml", settings.origin).href };
}

export function publicSitemap(settings = readPublicIndexing(), modelCodes: readonly string[] = []): MetadataRoute.Sitemap {
  if (!settings.enabled || !settings.origin) return [];
  // Callers supply only codes from the validated current approved publication.
  // No inferred private routes or invented lastModified dates.
  return approvedPublicPaths(modelCodes).map(path => ({ url: new URL(path, settings.origin!).href }));
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
sha256: "72a53423d3bc704fee36fbbda2c1864be814137d36557d6461fb78ee57ddf5e9"
variables: []
secrets_allowed: false
```
````typescript
import { publicRobots, readPublicIndexing } from "@/platform/seo/public-indexing";
import { loadPublishedCatalog } from "@/platform/catalog/load";

export const dynamic = "force-dynamic";
export default async function robots() {
  const settings = readPublicIndexing();
  if (!settings.enabled) return publicRobots(settings);
  return publicRobots(settings, (await loadPublishedCatalog())?.models.map(model => model.code) ?? []);
}
````

### FILE: `src/app/sitemap.ts`
```yaml
block_id: "TS-GO-WEB-SEO:sitemap.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "AUTHORED integration of exact Next16.3.4 metadata APIs; PUBLIC_WEB_METADATA_INTEGRATION_V377.md records official authority and executed evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "0a974a040a4942bff64dc4e48d74bd9bd5409995436f336c6d83465b34f9861a"
variables: []
secrets_allowed: false
```
````typescript
import { publicSitemap, readPublicIndexing } from "@/platform/seo/public-indexing";
import { loadPublishedCatalog } from "@/platform/catalog/load";

export const dynamic = "force-dynamic";
export default async function sitemap() {
  const settings = readPublicIndexing();
  if (!settings.enabled) return publicSitemap(settings);
  return publicSitemap(settings, (await loadPublishedCatalog())?.models.map(model => model.code) ?? []);
}
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
sha256: "e58b71d7852ae61519df18389c2f97e6351b3fbafe61a28f78d372809afca793"
variables: []
secrets_allowed: false
```
````typescript
// AUTHORED public-journey catalog. Runtime Intl supplies locale/date/plural semantics.
const es = {
  "nav.network": "Red de franquicia",
    "nav.help_content": "Contenido de ayuda",
  "site.description": "Catálogo y portales de una red de movilidad eléctrica.",
  "site.skip": "Saltar al contenido", "site.navigation": "Principal", "site.contact": "Contacto:",
  "nav.warranty": "Garantía", "nav.supply": "Suministro", "nav.catalog_editor": "Editar catálogo", "nav.training": "Capacitación", "nav.help": "Ayuda", "nav.dashboard": "Mi panel",
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
  "nav.network": "Franchise network",
    "nav.help_content": "Help content",
  "site.description": "Catalog and portals for an electric mobility network.",
  "site.skip": "Skip to content", "site.navigation": "Main", "site.contact": "Contact:",
  "nav.warranty": "Warranty", "nav.supply": "Supply", "nav.catalog_editor": "Edit catalog", "nav.training": "Training", "nav.help": "Help", "nav.dashboard": "My workspace",
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
sha256: "e53dd7d3ea10c5e5bae92388218b810cbdc9d69cc31940e253b719dc9e099441"
variables: []
secrets_allowed: false
```
````typescript
// AUTHORED exact extraction of the existing quote display, shared by private read views.
// Pinned Intl currency digits and BigInt parts preserve all safe minor-unit integers.
// No price, FX, tax, business rounding or acceptance policy is introduced.
export function minorAmountPresentation(minorUnits: number, currency: string, displayLocale = "es-AR") {
  let amountLabel = new Intl.Locale(displayLocale).language === "en" ? "Amount cannot be verified; contact support." : "Importe no verificable; consultá al soporte.";
  let amountValid = false;
  if (Number.isSafeInteger(minorUnits) && minorUnits >= 0 &&
      Intl.supportedValuesOf("currency").includes(currency)) {
    const formatter = new Intl.NumberFormat(displayLocale, { style: "currency", currency: currency, currencyDisplay: "code" });
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
sha256: "675badc36ec0bee912159db33ad579e9c470dd9f2d1fd22bfbd286d9514bd96b"
variables: []
secrets_allowed: false
```
````typescript
"use client";
import{usePrivateI18n}from"@/platform/i18n/private-provider";
type PrivateReadPath = "/admin" | "/customer" | "/factory";

export function PrivateReadFailure({ path }: { path: PrivateReadPath }) {
  const{t}=usePrivateI18n();
  return <section aria-labelledby="private-read-title">
    <h1 id="private-read-title" className="pageTitle">{t("p1194")}</h1>
    <p className="notice" role="alert">{t("p1195")}</p>
    <p>{t("p1196")}</p>
    <a className="button" href={path}>{t("p0033")}</a>
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
sha256: "67fdc2cac6bf208c60af27399ce82fdc99d8662d3b039df668c75dae5c1b9f1e"
variables: []
secrets_allowed: false
```
````typescript
import{PrivatePagingLinks}from"./private-paging-links";
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

export function PortalPaging(props:{label:string;nextHref:string|undefined;firstHref:string|undefined}){return <PrivatePagingLinks {...props}/>;}
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
sha256: "4cfb50f5698807e1f129b0b9919060363831524903a45687bf7f4c2324f5c70f"
variables: []
secrets_allowed: false
```

````tsx
"use client";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

// AUTHORED UI glue; an observed state never claims attribution to a lost request.
import{useEffect,useRef,useState}from"react";
import{factoryUnitSchema,type FactoryUnit}from"@/platform/factory/contracts";
export function FactoryUnitActions({unit,canWrite,storageScope}:{unit:FactoryUnit;canWrite:boolean;storageScope:string}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

 const[value,setValue]=useState(unit),[pending,setPending]=useState(false),[message,setMessage]=useState(""),[busy,setBusy]=useState(false);const active=useRef(false);
 const key=`elite:factory:${storageScope}:${unit.id}`;
 useEffect(()=>{try{setPending(localStorage.getItem(key)==="unknown")}catch{setPending(true);setMessage(t("p0299"))}},[key]);
 async function act(write:boolean){if(active.current)return;active.current=true;setBusy(true);setMessage("");
 try{
  if(write){localStorage.setItem(key,"unknown");setPending(true)}
  const q=new URLSearchParams({organizationId:unit.organization_id,unitId:unit.id});
  const response=await fetch(`/api/enterprise/factory${write?"":"?"+q}`,write?{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify({organizationId:unit.organization_id,unitId:unit.id,action:"start-assembly",expectedState:"planned"})}:{cache:"no-store"});
  if(!response.ok)throw new Error("unconfirmed");const body:unknown=await response.json();
  if(write){if(!(typeof body==="object"&&body!==null&&"status"in body&&body.status==="accepted"))throw new Error("invalid response");setMessage(t("p0300"))}
  else{const parsed=factoryUnitSchema.parse((body as {unit?:unknown}).unit);if(parsed.id!==unit.id||parsed.organization_id!==unit.organization_id)throw new Error("mismatch");setValue(parsed);localStorage.removeItem(key);setPending(false);setMessage(t("p0301"))}
 }catch{setMessage(write?t("p0302"):t("p0303"))}finally{active.current=false;setBusy(false)}}
 return <section aria-label={`Seguimiento de ${unit.serial_number}`}><p>{t("p0304")} <strong>{value.state}</strong></p>{pending?<p>{t("p0305")}</p>:null}
 {canWrite&&value.state==="planned"?<button className="button" type="button" disabled={busy||pending} onClick={()=>void act(true)}>{t("p0306")}</button>:null}
 <button className="button" type="button" disabled={busy} onClick={()=>void act(false)}>{t("p0307")}</button>
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

V402 composed delta: V402 source-backed stored-value integration: exact remaining provider due, explicit payment/funding XOR, shared approval, bounded browser transport and optional host. See STORED_VALUE_OPERATOR_FLOW_V402.md; source/pack admission successor governs final claim. Existing provider-only behavior retained.

V402 composed delta: J3 immutable approved catalog publication reuses original Commerce SQL/shared approval and optional host/public model owner; existing Next storefront consumes a validated published projection. No new dependencies or corporate attribution. CATALOG_CONNECTED_RELEASE_V402.md.

V402 composed delta: T2804 connected training: existing versioned help/audit/shared approvals/outbox/BFF; opt-in host and navigation; bounded body retains original default. No new dependency or automatic grant. TRAINING_CONNECTED_RELEASE_V402.md.

V402 composed delta: T2804 catalog role source/edit/review/publication transport reuses original model/price writers and catalog owner. Bounded PNG and text defaults retained. No new dependency. CATALOG_ROLE_AUTHORING_RELEASE_V402.md.

V402 composed delta: T2804 connected supply role: original Operations/BindPlan composed atomically; bounded forms, latest review/requester projection, nav and GET recovery. No migration/dependency/state-machine change. SUPPLY_ROLE_RELEASE_V402.md.

V402 composed delta: T2804 warranty roles: same authorized transaction projects existing immutable diagnosis/plan/work/quality/acceptance; quote version read, bounded forms and durable GET recovery. No migration/dependency/domain-rule change. WARRANTY_ROLE_RELEASE_V402.md.

V402 composed delta: T2804 network role: original four fulfillment SQL bodies extracted unchanged into one transaction with immutable result; optional host, forms and GET recovery. Migration0077, no dependency/domain-rule change. NETWORK_ROLE_RELEASE_V402.md.

V402 composed delta: T2804 help CMS and21same-release guides;15oldguides unchanged,5bounded training curricula revision2, no automatic grants. New optional host, shared inline text, actual browser/PG proofs. HELP_CMS_RELEASE_V402.md.

V402 composed delta: T2804 private es/en display, per-user/tenant preference and browser negotiation;1198messages,21exact guides,5hash-bound curricula; original commands/content/policies preserved. PRIVATE_LOCALE_RELEASE_V402.md/json. AUTHORED glue; no new dependency or corporate attribution.

### FILE: `src/components/private-language-switcher.tsx`

```yaml
block_id: "TS-GO-API-WEB-BRIDGE-LOCALE-DELTA:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "ed07036719291f88bc03dcc89cc468853fd71c2c84075c6b51f63a7138a788ab"
variables: []
secrets_allowed: false
```

````tsx
"use client";
import{useRef,useState,type FormEvent}from"react";
import{usePrivateI18n}from"@/platform/i18n/private-provider";
export function PrivateLanguageSwitcher(){
 const{locale}=usePrivateI18n(),en=locale.language==="en";
 const[selected,setSelected]=useState(locale.language),[busy,setBusy]=useState(false),[failed,setFailed]=useState(false);const fence=useRef(false);
 async function save(event:FormEvent){event.preventDefault();if(fence.current||selected===locale.language)return;fence.current=true;setBusy(true);setFailed(false);
  try{const r=await fetch("/api/enterprise/locale",{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify({language:selected}),cache:"no-store",redirect:"error",signal:AbortSignal.timeout(7000)});if(!r.ok)throw new Error();window.location.reload()}
  catch{setFailed(true);setBusy(false);fence.current=false}
 }
 return <details><summary>{en?"Display language":"Idioma de la interfaz"}</summary><form onSubmit={save}>
 <label>{en?"Language":"Idioma"}<select value={selected} disabled={busy} onChange={e=>setSelected(e.target.value==="en"?"en":"es")}><option value="es" lang="es">Español</option><option value="en" lang="en">English</option></select></label>
 <p>{en?"Saving reloads this page. Finish or save any unsaved form first.":"Guardar recarga esta página. Terminá o guardá primero cualquier formulario sin guardar."}</p>
 <button className="button" disabled={busy||selected===locale.language}>{en?"Save language and reload":"Guardar idioma y recargar"}</button>
 {failed&&<p role="alert">{en?"Language save was not confirmed. No business operation was sent. Retry saving the preference or reload to check it.":"No se confirmó el guardado del idioma. No se envió una operación de negocio. Reintentá guardar la preferencia o recargá para consultarla."}</p>}
 </form></details>
}
````

### FILE: `src/platform/backend/private-paging-links.tsx`

```yaml
block_id: "TS-GO-API-WEB-BRIDGE-LOCALE-DELTA:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f1469176e541f7a13e0ae06c19402e6f3ef98e42ddfd30d6a35bd33e83db6937"
variables: []
secrets_allowed: false
```

````tsx
"use client";
import{usePrivateI18n}from"@/platform/i18n/private-provider";
export function PrivatePagingLinks({label,nextHref,firstHref}:{label:string;nextHref:string|undefined;firstHref:string|undefined}){
 const{t}=usePrivateI18n();if(!nextHref&&!firstHref)return null;
 return <nav aria-label={label}>{firstHref?<p><a href={firstHref}>{t("p1197")}</a></p>:null}{nextHref?<p><a href={nextHref}>{t("p1198")}</a></p>:null}</nav>;
}
````

### FILE: `src/platform/i18n/load-private-i18n.ts`

```yaml
block_id: "TS-GO-API-WEB-BRIDGE-LOCALE-DELTA:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "a46669fe5a82a75e70614bd4ac485a1f3267bbcf1871b842e7a74c4f8a08e970"
variables: []
secrets_allowed: false
```

````typescript
import"server-only";
import{loadPrivateLocale}from"./load-private-locale";
import{privateTranslator,controlledPrivateLabel}from"./private-catalog";
export async function loadPrivateI18n(){const locale=await loadPrivateLocale();return{locale,t:privateTranslator(locale.language),controlled:(value:string)=>controlledPrivateLabel(locale.language,value)}}
````

### FILE: `src/platform/i18n/load-private-locale.ts`

```yaml
block_id: "TS-GO-API-WEB-BRIDGE-LOCALE-DELTA:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "7ee2c9cfae57b8d9b87f3776400c222a2a33944917d10cf2cab5cca20ff48e29"
variables: []
secrets_allowed: false
```

````typescript
import"server-only";
import{cache}from"react";
import{createHash}from"node:crypto";
import{cookies,headers}from"next/headers";
import{loadBusinessConfig}from"@/platform/config/load";
import{readSession,type PortalSession}from"@/platform/auth/session";
import{resolvePrivateLocale}from"./private-locale";
// Non-secret display preference is namespaced to the authenticated user/tenant.
// It never contributes to authorization or to the signed session.
export function privateLocaleCookieName(s:Pick<PortalSession,"tenantId"|"subject">){
 return(process.env.NODE_ENV==="production"?"__Host-":"")+"elite_locale_"+createHash("sha256").update(JSON.stringify([s.tenantId,s.subject])).digest("hex");
}
export const loadPrivateLocale=cache(async()=>{
 const [config,session]=await Promise.all([loadBusinessConfig(),readSession()]);
 const market=config.markets.find(x=>x.code===config.business.defaultMarket);if(!market)throw new Error("configured market missing");
 // Anonymous public pages keep their configured public language.
 if(!session)return resolvePrivateLocale(config.business.defaultLocale,market.timeZone,null,null);
 const [jar,requestHeaders]=await Promise.all([cookies(),headers()]);
 return resolvePrivateLocale(config.business.defaultLocale,market.timeZone,jar.get(privateLocaleCookieName(session))?.value,requestHeaders.get("accept-language"));
});
````

### FILE: `src/platform/i18n/private-catalog.ts`

```yaml
block_id: "TS-GO-API-WEB-BRIDGE-LOCALE-DELTA:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "83f9644322e911f88aea0b7c155e1a93f3b5964b526832b4a15674e94ff18532"
variables: []
secrets_allowed: false
```

````typescript
// AUTHORED es/en presentation catalog for local glue; no corporate code attribution.
import messages from "./private-messages.json";
import type {PrivateLanguage} from "./private-locale";
export type PrivateMessageKey=keyof typeof messages;
type Parameters={
 p0248:{resolution_action:string};
 p0255:{completed_at:string};
 p0267:{accepted_at:string};
 p0310:{request_id:string};
 p0629:{handover_id:string;handover_state:string};
 p0645:{evaluated_at:string;valid_until:string};
 p0646:{evaluated_at:string};
};
type Arguments<K extends PrivateMessageKey>=K extends keyof Parameters?[values:Parameters[K]]:[];
export function privateMessage<K extends PrivateMessageKey>(language:PrivateLanguage,key:K,...args:Arguments<K>):string{
 const value=messages[key][language];
 if(args.length===0)return value;
 const values=args[0] as Record<string,string>;
 return value.replace(/\{([a-z_]+)\}/g,(_,name:string)=>{if(typeof values[name]!=="string")throw new Error("missing display parameter");return values[name]});
}
export const privateTranslator=(language:PrivateLanguage)=><K extends PrivateMessageKey>(key:K,...args:Arguments<K>)=>privateMessage(language,key,...args);
export type PrivateTranslator=ReturnType<typeof privateTranslator>;
const known=new Map<string,PrivateMessageKey>(Object.entries(messages).filter(([,m])=>!m.es.includes("{")).map(([key,m])=>[m.es,key as PrivateMessageKey]));
// Call ONLY for locally controlled display labels, never CMS/provider/user content.
export function controlledPrivateLabel(language:PrivateLanguage,value:string):string{
 const key=known.get(value);return key===undefined?value:privateMessage(language,key as Exclude<PrivateMessageKey,keyof Parameters>);
}
````

### FILE: `src/platform/i18n/private-locale.test.ts`

```yaml
block_id: "TS-GO-API-WEB-BRIDGE-LOCALE-DELTA:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "3d24bb7c3bfa3d78919f9c4fe98d97ae18e30a8f20504e3d82246ace8316fc77"
variables: []
secrets_allowed: false
```

````typescript
import{it,expect}from"vitest";
import{browserLanguage,resolvePrivateLocale}from"./private-locale";
it("prefers validated explicit choice and preserves configured timezone",()=>{expect(resolvePrivateLocale("es-AR","America/Argentina/Buenos_Aires","en","es")).toEqual({language:"en",locale:"en-US",timeZone:"America/Argentina/Buenos_Aires",source:"preference"});expect(resolvePrivateLocale("es-AR","UTC","es","en")).toMatchObject({locale:"es-AR",source:"preference"})});
it("orders supported browser choices by quality and honors zero",()=>{expect(browserLanguage("es;q=0.4,en-GB;q=0.9")).toBe("en");expect(browserLanguage("en;q=0, en-US;q=1,es;q=0.5")).toBe("es");expect(browserLanguage("en;q=0.5,es;q=0.5")).toBe("en")});
it("malformed excessive and unsupported hints never select a policy",()=>{for(const raw of["fr,zh","en;q=5","en;q=-1","en;q=0.1234","*","x".repeat(513),Array(17).fill("en").join(",")])expect(browserLanguage(raw)).toBe(null);expect(resolvePrivateLocale("es-AR","UTC",{language:"en",tenant:"other"},"fr")).toMatchObject({language:"es",source:"configuration"});expect(()=>resolvePrivateLocale("en-US","Wrong/Zone","en",null)).toThrow()});
````

### FILE: `src/platform/i18n/private-locale.ts`

```yaml
block_id: "TS-GO-API-WEB-BRIDGE-LOCALE-DELTA:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "befb0c0e4d1714087a99058b130bec60a370789a12e0b57703a58e3e60182c97"
variables: []
secrets_allowed: false
```

````typescript
// AUTHORED bounded es/en display preference; Intl remains the locale authority.
import{resolvePublicLocale,type PublicLanguage}from"./public-catalog";
export type PrivateLanguage=PublicLanguage;
export type PrivateLocale={language:PrivateLanguage;locale:string;timeZone:string;source:"preference"|"browser"|"configuration"};
export function privateLanguage(v:unknown):PrivateLanguage|null{return v==="es"||v==="en"?v:null}
// Narrow two-catalog preference, not a general RFC4647 matcher. Invalid ranges
// are ignored; q=0 excludes their base language. Oversize/unfulfillable headers
// are disregarded and the explicit configured catalog is used (RFC9110 12.5.4).
export function browserLanguage(header:string|null):PrivateLanguage|null{
 if(!header||header.length>512)return null;
 const ranges=header.split(",");if(ranges.length>16)return null;
 const candidates:{language:PrivateLanguage;quality:number;order:number}[]=[],excluded=new Set<PrivateLanguage>();
 for(const [order,raw]of ranges.entries()){
  const match=/^\s*([A-Za-z]{1,8}(?:-[A-Za-z0-9]{1,8})*)(?:\s*;\s*q=(0(?:\.[0-9]{0,3})?|1(?:\.0{0,3})?))?\s*$/i.exec(raw);if(!match)continue;
  try{const canonical=Intl.getCanonicalLocales(match[1]!)[0];if(!canonical)continue;const language=privateLanguage(new Intl.Locale(canonical).language);if(!language)continue;const quality=match[2]===undefined?1:Number(match[2]);if(quality===0)excluded.add(language);else candidates.push({language,quality,order})}catch{/* malformed tag has no authority */}
 }
 return candidates.filter(x=>!excluded.has(x.language)).sort((a,b)=>b.quality-a.quality||a.order-b.order)[0]?.language??null;
}
export function resolvePrivateLocale(configured:string,timeZone:string,preference:unknown,header:string|null):PrivateLocale{
 const baseline=resolvePublicLocale(configured,timeZone),chosen=privateLanguage(preference),browser=chosen?null:browserLanguage(header),language=chosen??browser??baseline.language;
 return{language,locale:language===baseline.language?baseline.locale:language==="en"?"en-US":"es-AR",timeZone:baseline.timeZone,source:chosen?"preference":browser?"browser":"configuration"};
}
````

### FILE: `src/platform/i18n/private-messages.json`

```yaml
block_id: "TS-GO-API-WEB-BRIDGE-LOCALE-DELTA:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "0842b07852354e6ea5b0aba93624d6bac143e1f5f2955ad97f9cf94fcee947fd"
variables: []
secrets_allowed: false
```

````json
{
  "p0001": {
    "es": "Catálogo",
    "en": "Catalog"
  },
  "p0002": {
    "es": "No tenés permiso para consultar este espacio.",
    "en": "You do not have permission to view this workspace."
  },
  "p0003": {
    "es": "La sesión no tiene una organización disponible.",
    "en": "No organization is available in this session."
  },
  "p0004": {
    "es": "Edición y publicación del catálogo",
    "en": "Catalog editing and publication"
  },
  "p0005": {
    "es": "Organización:",
    "en": "Organization:"
  },
  "p0006": {
    "es": "Acceso denegado",
    "en": "Access denied"
  },
  "p0007": {
    "es": "La sesión no posee",
    "en": "This session does not have"
  },
  "p0008": {
    "es": ".",
    "en": "."
  },
  "p0009": {
    "es": "Portal operativo ·",
    "en": "Operations portal ·"
  },
  "p0010": {
    "es": "Operación",
    "en": "Operations"
  },
  "p0011": {
    "es": "Pedidos",
    "en": "Orders"
  },
  "p0012": {
    "es": "Stock disponible",
    "en": "Available stock"
  },
  "p0013": {
    "es": "Casos abiertos",
    "en": "Open cases"
  },
  "p0014": {
    "es": "Envíos activos",
    "en": "Active shipments"
  },
  "p0015": {
    "es": "Pedidos recientes",
    "en": "Recent orders"
  },
  "p0016": {
    "es": "·",
    "en": "·"
  },
  "p0017": {
    "es": "No hay registros en esta página.",
    "en": "There are no records on this page."
  },
  "p0018": {
    "es": "Servicio",
    "en": "Service"
  },
  "p0019": {
    "es": "Oportunidades",
    "en": "Leads"
  },
  "p0020": {
    "es": "sin asignar",
    "en": "unassigned"
  },
  "p0021": {
    "es": "Resultados de encuestas",
    "en": "Survey results"
  },
  "p0022": {
    "es": "Abrí el enlace de resultados de la encuesta configurada.",
    "en": "Open the results link for the configured survey."
  },
  "p0023": {
    "es": "No tenés acceso a estos resultados.",
    "en": "You do not have access to these results."
  },
  "p0024": {
    "es": "Resultados de la encuesta",
    "en": "Survey results"
  },
  "p0025": {
    "es": "Respuestas:",
    "en": "Responses:"
  },
  "p0026": {
    "es": "NPS:",
    "en": "NPS:"
  },
  "p0027": {
    "es": "El resultado todavía no está disponible según el mínimo de respuestas configurado.",
    "en": "The result is not yet available under the configured minimum response count."
  },
  "p0028": {
    "es": "El resultado describe estas respuestas. No demuestra que la muestra represente a todos los clientes.",
    "en": "The result describes these responses. It does not establish that the sample represents all customers."
  },
  "p0029": {
    "es": "Ayuda de resultados",
    "en": "Results help"
  },
  "p0030": {
    "es": "El cálculo resta el porcentaje de puntuaciones de 0 a 6 al de puntuaciones de 9 a 10. Las de 7 y 8 se incluyen en el total.",
    "en": "The calculation subtracts the percentage of scores from 0 to 6 from the percentage of scores from 9 to 10. Scores of 7 and 8 are included in the total."
  },
  "p0031": {
    "es": "Ayuda de encuestas, versión 1.0.0.",
    "en": "Survey help, version 1.0.0."
  },
  "p0032": {
    "es": "Resultados no disponibles",
    "en": "Results unavailable"
  },
  "p0033": {
    "es": "Volver a consultar",
    "en": "Check again"
  },
  "p0034": {
    "es": "Portal de cliente ·",
    "en": "Customer portal ·"
  },
  "p0035": {
    "es": "Mis turnos",
    "en": "My appointments"
  },
  "p0036": {
    "es": "La cancelación se limita a tus turnos futuros activos, exige la versión vigente y deja motivo auditable.",
    "en": "Cancellation is limited to your active future appointments, requires the current version, and records an auditable reason."
  },
  "p0037": {
    "es": "Ver cotizaciones",
    "en": "View quotes"
  },
  "p0038": {
    "es": "Ver entregas",
    "en": "View deliveries"
  },
  "p0039": {
    "es": "Mis entregas",
    "en": "My deliveries"
  },
  "p0040": {
    "es": "Antes de aceptar verás la preparación completada, sus controles y la versión exacta. También podés rechazar esa presentación y seguir su corrección, devolución o cambio sin perder el historial.",
    "en": "Before accepting, you will see the completed preparation, its checks, and the exact version. You can also reject that handover and follow its correction, return, or exchange while preserving its history."
  },
  "p0041": {
    "es": "Ver turnos",
    "en": "View appointments"
  },
  "p0042": {
    "es": "Por ahora no hay entregas disponibles para consultar.",
    "en": "There are currently no deliveries available to view."
  },
  "p0043": {
    "es": "No pudimos verificar tus entregas. No se confirmó ni modificó ninguna entrega desde esta consulta.",
    "en": "We could not verify your deliveries. This lookup did not confirm or change any delivery."
  },
  "p0044": {
    "es": "Volver a consultar entregas",
    "en": "Check deliveries again"
  },
  "p0045": {
    "es": "Portal de cliente",
    "en": "Customer portal"
  },
  "p0046": {
    "es": "Mi cuenta",
    "en": "My account"
  },
  "p0047": {
    "es": "Mis pedidos",
    "en": "My orders"
  },
  "p0048": {
    "es": "Mis casos de servicio",
    "en": "My service cases"
  },
  "p0049": {
    "es": "Gestionar mis turnos",
    "en": "Manage my appointments"
  },
  "p0050": {
    "es": "(",
    "en": "("
  },
  "p0051": {
    "es": ")",
    "en": ")"
  },
  "p0052": {
    "es": "Mis cotizaciones",
    "en": "My quotes"
  },
  "p0053": {
    "es": "Vigencia finalizada",
    "en": "Expired"
  },
  "p0054": {
    "es": "Disponible para aceptar",
    "en": "Available to accept"
  },
  "p0055": {
    "es": "Aceptada",
    "en": "Accepted"
  },
  "p0056": {
    "es": "Cancelada",
    "en": "Cancelled"
  },
  "p0057": {
    "es": "En preparación",
    "en": "In preparation"
  },
  "p0058": {
    "es": "Estado pendiente de verificar",
    "en": "Status awaiting verification"
  },
  "p0059": {
    "es": "Fecha no verificable; consultá al soporte.",
    "en": "Date cannot be verified; contact support."
  },
  "p0060": {
    "es": "No pudimos consultar tus cotizaciones. Esto no indica que la aceptación haya fallado. No repitas la operación sin consultar su estado.",
    "en": "We could not retrieve your quotes. This does not mean acceptance failed. Do not repeat the operation without checking its status."
  },
  "p0061": {
    "es": "El precio y la moneda provienen del servidor. Aceptar crea un único pedido en la misma transacción.",
    "en": "The price and currency come from the server. Acceptance creates a single order in the same transaction."
  },
  "p0062": {
    "es": "Ver y cancelar turnos",
    "en": "View and cancel appointments"
  },
  "p0063": {
    "es": "Encuestas",
    "en": "Surveys"
  },
  "p0064": {
    "es": "Abrí el enlace de la encuesta que recibiste. Si el enlace está incompleto, solicitá uno nuevo.",
    "en": "Open the survey link you received. If the link is incomplete, request a new one."
  },
  "p0065": {
    "es": "No tenés acceso a esta encuesta.",
    "en": "You do not have access to this survey."
  },
  "p0066": {
    "es": "Encuesta no disponible",
    "en": "Survey unavailable"
  },
  "p0067": {
    "es": "No pudimos cargar la encuesta. Tu respuesta no se envió desde esta página.",
    "en": "We could not load the survey. Your response was not submitted from this page."
  },
  "p0068": {
    "es": "Panel",
    "en": "Workspace"
  },
  "p0069": {
    "es": "Tu espacio de trabajo",
    "en": "Your workspace"
  },
  "p0070": {
    "es": "Elegí un acceso para continuar. Consultá la ayuda de cada recorrido antes de actuar.",
    "en": "Choose an available section to continue. Read the help for each workflow before taking action."
  },
  "p0071": {
    "es": "Portal de fábrica ·",
    "en": "Factory portal ·"
  },
  "p0072": {
    "es": "Abastecimiento y seguimiento",
    "en": "Supply and tracking"
  },
  "p0073": {
    "es": "Orden",
    "en": "Order"
  },
  "p0074": {
    "es": "La sesión no posee permisos para esta operación.",
    "en": "This session does not have permission for this operation."
  },
  "p0075": {
    "es": "Fecha inválida",
    "en": "Invalid date"
  },
  "p0076": {
    "es": "Volver a la agenda",
    "en": "Back to schedule"
  },
  "p0077": {
    "es": "discrepancias de entrega",
    "en": "delivery discrepancies"
  },
  "p0078": {
    "es": "Franquicia ·",
    "en": "Franchise ·"
  },
  "p0079": {
    "es": "Operación comercial, agenda y entregas",
    "en": "Sales operations, scheduling, and deliveries"
  },
  "p0080": {
    "es": "Consultá pedidos y agenda. Cada acción requiere sus permisos y conserva su estado.",
    "en": "View orders and appointments. Each action requires its permissions and preserves its state."
  },
  "p0081": {
    "es": "Operar gift cards y fidelidad",
    "en": "Manage gift cards and loyalty"
  },
  "p0082": {
    "es": "Revisar respuestas de WhatsApp",
    "en": "Review WhatsApp replies"
  },
  "p0083": {
    "es": "No pudimos consultar",
    "en": "We could not retrieve"
  },
  "p0084": {
    "es": ". No se muestran datos de esa sección.",
    "en": ". No data from that section is shown."
  },
  "p0085": {
    "es": "Volver a consultar operación",
    "en": "Check operations again"
  },
  "p0086": {
    "es": "Fecha de agenda (UTC)",
    "en": "Schedule date (UTC)"
  },
  "p0087": {
    "es": "Ver fecha",
    "en": "View date"
  },
  "p0088": {
    "es": "Gift cards y fidelidad",
    "en": "Gift cards and loyalty"
  },
  "p0089": {
    "es": "Tu sesión no permite consultar estas operaciones.",
    "en": "Your session does not allow you to view these operations."
  },
  "p0090": {
    "es": "Elegí una organización autorizada.",
    "en": "Choose an authorized organization."
  },
  "p0091": {
    "es": "Organización",
    "en": "Organization"
  },
  "p0092": {
    "es": "Cambiar organización",
    "en": "Change organization"
  },
  "p0093": {
    "es": "Operación asistida con revisión de otra persona. Los canjes se aplican al pedido antes de solicitar su pago.",
    "en": "Assisted operation with review by another person. Redemptions are applied to the order before requesting payment."
  },
  "p0094": {
    "es": "No pudimos cargar el programa habilitado para esta organización. Consultá de nuevo antes de operar.",
    "en": "We could not load the program enabled for this organization. Check again before proceeding."
  },
  "p0095": {
    "es": "Consultar de nuevo",
    "en": "Check again"
  },
  "p0096": {
    "es": "Respuestas de WhatsApp",
    "en": "WhatsApp replies"
  },
  "p0097": {
    "es": "Tu sesión no permite revisar respuestas.",
    "en": "Your session does not allow you to review replies."
  },
  "p0098": {
    "es": "Revisá el destinatario y el texto completo. Aprobar registra tu decisión; enviar es una acción separada.",
    "en": "Review the recipient and the complete text. Approval records your decision; sending is a separate action."
  },
  "p0099": {
    "es": "No pudimos consultar las respuestas. Volvé a consultar su estado antes de intentar un envío.",
    "en": "We could not retrieve the replies. Check their status again before attempting to send."
  },
  "p0100": {
    "es": "Guía de uso",
    "en": "User guide"
  },
  "p0101": {
    "es": "Elegí un enfoque",
    "en": "Choose a perspective"
  },
  "p0102": {
    "es": "Podés consultar el recorrido de un rol. Tus accesos no cambian.",
    "en": "You can view a role's workflow. Your access remains unchanged."
  },
  "p0103": {
    "es": "Ver recorrido de",
    "en": "View workflow for"
  },
  "p0104": {
    "es": "Volver a mi panel",
    "en": "Back to my workspace"
  },
  "p0105": {
    "es": "Capacitación",
    "en": "Training"
  },
  "p0106": {
    "es": "Tu sesión no permite acceder a capacitación.",
    "en": "Your session does not allow access to training."
  },
  "p0107": {
    "es": "Capacitación y evaluación",
    "en": "Training and assessment"
  },
  "p0108": {
    "es": "Leé la versión indicada, registrá tu lectura y presentá tus respuestas para revisión humana.",
    "en": "Read the specified version, record that you have read it, and submit your answers for human review."
  },
  "p0109": {
    "es": "No pudimos consultar el contenido o las evaluaciones. Conservá tu referencia y consultá de nuevo.",
    "en": "We could not retrieve the content or assessments. Keep your reference and check again."
  },
  "p0110": {
    "es": "Recorrido:",
    "en": "Workflow:"
  },
  "p0111": {
    "es": "Elegir este enfoque no cambia tus accesos.",
    "en": "Choosing this perspective does not change your access."
  },
  "p0112": {
    "es": "Consultá la ayuda de los recorridos disponibles. Leer esta página no acredita capacitación ni habilita operaciones.",
    "en": "Read the help for the available workflows. Reading this page does not certify training or enable operations."
  },
  "p0113": {
    "es": "Abrir ayuda versionada",
    "en": "Open versioned help"
  },
  "p0114": {
    "es": "Volver a la guía",
    "en": "Back to guide"
  },
  "p0115": {
    "es": "Contenido de ayuda",
    "en": "Help content"
  },
  "p0116": {
    "es": "Red de franquicia",
    "en": "Franchise network"
  },
  "p0117": {
    "es": "Organizaciones, sucursales y acuerdos.",
    "en": "Organizations, branches, and agreements."
  },
  "p0118": {
    "es": "Suministro",
    "en": "Supply"
  },
  "p0119": {
    "es": "Suministro por serie",
    "en": "Serial supply"
  },
  "p0120": {
    "es": ". Vista:",
    "en": ". View:"
  },
  "p0121": {
    "es": "Fábrica",
    "en": "Factory"
  },
  "p0122": {
    "es": "Compras y recepción",
    "en": "Purchasing and receiving"
  },
  "p0123": {
    "es": "Garantía",
    "en": "Warranty"
  },
  "p0124": {
    "es": "Garantía y reparación",
    "en": "Warranty and repair"
  },
  "p0125": {
    "es": "Cliente",
    "en": "Customer"
  },
  "p0126": {
    "es": "No pudimos verificar el estado. No se envió ningún mensaje desde esta consulta. Volvé a consultar o pedí revisión a soporte.",
    "en": "We could not verify the status. This lookup did not send any message. Check again or request a support review."
  },
  "p0127": {
    "es": "Estado de notificaciones WhatsApp",
    "en": "WhatsApp notification status"
  },
  "p0128": {
    "es": "Notificación de este turno",
    "en": "Notification for this appointment"
  },
  "p0129": {
    "es": "Consultando notificación…",
    "en": "Checking notification…"
  },
  "p0130": {
    "es": "Consultar estado de WhatsApp",
    "en": "Check WhatsApp status"
  },
  "p0131": {
    "es": "Resultado de consulta de WhatsApp",
    "en": "WhatsApp lookup result"
  },
  "p0132": {
    "es": "Consultando el registro guardado…",
    "en": "Checking the saved record…"
  },
  "p0133": {
    "es": "Consulta actualizada. No se enviaron mensajes.",
    "en": "Lookup updated. No messages were sent."
  },
  "p0134": {
    "es": "Consulta de sólo lectura, sin envíos automáticos.",
    "en": "Read-only lookup, with no automatic sending."
  },
  "p0135": {
    "es": "No hay aprobaciones de notificación registradas para este turno. Esto no significa que el cliente haya recibido un mensaje.",
    "en": "No notification approvals are recorded for this appointment. This does not mean the customer has received a message."
  },
  "p0136": {
    "es": "Consultado:",
    "en": "Checked:"
  },
  "p0137": {
    "es": "Fecha del aviso del proveedor:",
    "en": "Provider notification date:"
  },
  "p0138": {
    "es": "La aprobación de envío está vencida. Su historial sigue visible; no autoriza otro envío.",
    "en": "The sending approval has expired. Its history remains visible; it does not authorize another send."
  },
  "p0139": {
    "es": "Ayuda y referencia para soporte",
    "en": "Help and support reference"
  },
  "p0140": {
    "es": "Turno",
    "en": "Appointment"
  },
  "p0141": {
    "es": "Confirmación",
    "en": "Confirmation"
  },
  "p0142": {
    "es": "Guía",
    "en": "Guide"
  },
  "p0143": {
    "es": "Abrir esta guía en Ayuda",
    "en": "Open this guide in Help"
  },
  "p0144": {
    "es": "Legal",
    "en": "Legal"
  },
  "p0145": {
    "es": "Técnica",
    "en": "Technical"
  },
  "p0146": {
    "es": "Imagen",
    "en": "Image"
  },
  "p0147": {
    "es": "Publicación",
    "en": "Publication"
  },
  "p0148": {
    "es": "No pudimos leer las referencias guardadas. Conservá esta pestaña y pedí revisión antes de operar.",
    "en": "We could not read the saved references. Keep this tab open and request a review before proceeding."
  },
  "p0149": {
    "es": "Borrador y revisiones consultados.",
    "en": "Draft and reviews retrieved."
  },
  "p0150": {
    "es": "No pudimos consultar ese borrador. Revisá la referencia y volvé a consultar.",
    "en": "We could not retrieve that draft. Check the reference and try again."
  },
  "p0151": {
    "es": "Publicación vigente consultada.",
    "en": "Current publication retrieved."
  },
  "p0152": {
    "es": "No pudimos confirmar la publicación vigente.",
    "en": "We could not confirm the current publication."
  },
  "p0153": {
    "es": "Modelo, variantes y precios guardados como borrador.",
    "en": "Model, variants, and prices saved as a draft."
  },
  "p0154": {
    "es": "Imagen recibida y normalizada.",
    "en": "Image received and normalized."
  },
  "p0155": {
    "es": "Versión guardada. Compartí la referencia con las personas revisoras.",
    "en": "Version saved. Share the reference with the reviewers."
  },
  "p0156": {
    "es": "Decisión registrada.",
    "en": "Decision recorded."
  },
  "p0157": {
    "es": "Publicación registrada. Consultá su vigencia.",
    "en": "Publication recorded. Check whether it is current."
  },
  "p0158": {
    "es": "Decisión registrada. Volvé a consultar el borrador para continuar.",
    "en": "Decision recorded. Retrieve the draft again to continue."
  },
  "p0159": {
    "es": "Resultado sin confirmar. Conservamos su referencia: consultá el resultado pendiente antes de otra escritura.",
    "en": "Result unconfirmed. We have kept its reference: check the pending result before another write."
  },
  "p0160": {
    "es": "No se envió una escritura nueva, o fue rechazada antes de llegar al servicio. Revisá los campos y permisos.",
    "en": "No new write was sent, or it was rejected before reaching the service. Check the fields and permissions."
  },
  "p0161": {
    "es": "Resultado recuperado desde el registro, sin repetir la escritura.",
    "en": "Result recovered from the record without repeating the write."
  },
  "p0162": {
    "es": "El resultado todavía no está confirmado. Conservá la referencia y volvé a consultar.",
    "en": "The result is still unconfirmed. Keep the reference and check again."
  },
  "p0163": {
    "es": "Revisá los datos del modelo, los precios y la vigencia UTC.",
    "en": "Check the model data, prices, and UTC validity period."
  },
  "p0164": {
    "es": "Seleccioná una imagen PNG de hasta 1 MiB.",
    "en": "Select a PNG image up to 1 MiB."
  },
  "p0165": {
    "es": "La evidencia debe ser un archivo de hasta 1 MiB.",
    "en": "Evidence must be a file up to 1 MiB."
  },
  "p0166": {
    "es": "Resultado pendiente",
    "en": "Pending result"
  },
  "p0167": {
    "es": "Referencia:",
    "en": "Reference:"
  },
  "p0168": {
    "es": "Consultar resultado pendiente",
    "en": "Check pending result"
  },
  "p0169": {
    "es": "Crear modelo y precios en borrador",
    "en": "Create draft model and prices"
  },
  "p0170": {
    "es": "Estos datos necesitan revisión antes de publicarse. Las fechas se ingresan en UTC; los importes, en unidades menores de la moneda configurada.",
    "en": "These details require review before publication. Dates are entered in UTC; amounts are in minor units of the configured currency."
  },
  "p0171": {
    "es": "Código del modelo",
    "en": "Model code"
  },
  "p0172": {
    "es": "Nombre del modelo",
    "en": "Model name"
  },
  "p0173": {
    "es": "Clase de vehículo",
    "en": "Vehicle class"
  },
  "p0174": {
    "es": "Bicicleta",
    "en": "Bicycle"
  },
  "p0175": {
    "es": "Motocicleta",
    "en": "Motorcycle"
  },
  "p0176": {
    "es": "Scooter",
    "en": "Scooter"
  },
  "p0177": {
    "es": "Utilitario",
    "en": "Utility vehicle"
  },
  "p0178": {
    "es": "Otro",
    "en": "Other"
  },
  "p0179": {
    "es": "Descripción técnica",
    "en": "Technical description"
  },
  "p0180": {
    "es": "Vigencia desde (UTC)",
    "en": "Valid from (UTC)"
  },
  "p0181": {
    "es": "Vigencia hasta (UTC, opcional)",
    "en": "Valid until (UTC, optional)"
  },
  "p0182": {
    "es": "Variante",
    "en": "Variant"
  },
  "p0183": {
    "es": "Código de variante",
    "en": "Variant code"
  },
  "p0184": {
    "es": "Nombre de variante",
    "en": "Variant name"
  },
  "p0185": {
    "es": "Especificación de batería",
    "en": "Battery specification"
  },
  "p0186": {
    "es": "Importe en unidades menores",
    "en": "Amount in minor units"
  },
  "p0187": {
    "es": "Tratamiento del precio",
    "en": "Price tax treatment"
  },
  "p0188": {
    "es": "No aplica impuesto en este precio",
    "en": "No tax applies to this price"
  },
  "p0189": {
    "es": "Impuesto incluido",
    "en": "Tax included"
  },
  "p0190": {
    "es": "Impuesto excluido",
    "en": "Tax excluded"
  },
  "p0191": {
    "es": "Quitar variante",
    "en": "Remove variant"
  },
  "p0192": {
    "es": "Agregar variante",
    "en": "Add variant"
  },
  "p0193": {
    "es": "Guardar modelo y precios",
    "en": "Save model and prices"
  },
  "p0194": {
    "es": "Imagen y versión para revisar",
    "en": "Image and version for review"
  },
  "p0195": {
    "es": "Imagen PNG",
    "en": "PNG image"
  },
  "p0196": {
    "es": "Guardar imagen",
    "en": "Save image"
  },
  "p0197": {
    "es": "Referencia del modelo",
    "en": "Model reference"
  },
  "p0198": {
    "es": "Referencia de la imagen",
    "en": "Image reference"
  },
  "p0199": {
    "es": "Referencia del libro de precios",
    "en": "Price book reference"
  },
  "p0200": {
    "es": "Guardar versión para revisión",
    "en": "Save version for review"
  },
  "p0201": {
    "es": "Consultar y revisar versión",
    "en": "Retrieve and review version"
  },
  "p0202": {
    "es": "Referencia del borrador",
    "en": "Draft reference"
  },
  "p0203": {
    "es": "Consultar borrador",
    "en": "Retrieve draft"
  },
  "p0204": {
    "es": "Creado por",
    "en": "Created by"
  },
  "p0205": {
    "es": "Enlace de esta versión para otra persona autorizada",
    "en": "Link to this version for another authorized person"
  },
  "p0206": {
    "es": "Código:",
    "en": "Code:"
  },
  "p0207": {
    "es": ". Clase:",
    "en": ". Class:"
  },
  "p0208": {
    "es": "Precio:",
    "en": "Price:"
  },
  "p0209": {
    "es": "unidades menores de",
    "en": "minor units of"
  },
  "p0210": {
    "es": "; tratamiento:",
    "en": "; treatment:"
  },
  "p0211": {
    "es": ". Homologación registrada:",
    "en": ". Recorded type approval:"
  },
  "p0212": {
    "es": "Estado de revisiones",
    "en": "Review status"
  },
  "p0213": {
    "es": ":",
    "en": ":"
  },
  "p0214": {
    "es": "Contenido exacto de la versión",
    "en": "Exact version content"
  },
  "p0215": {
    "es": "SHA-256:",
    "en": "SHA-256:"
  },
  "p0216": {
    "es": "Otra persona autorizada debe revisar y publicar esta versión.",
    "en": "Another authorized person must review and publish this version."
  },
  "p0217": {
    "es": "Etapa de revisión",
    "en": "Review stage"
  },
  "p0218": {
    "es": "Evidencia de la revisión",
    "en": "Review evidence"
  },
  "p0219": {
    "es": "Se registra la huella del archivo, no su contenido.",
    "en": "The file's hash is recorded, not its content."
  },
  "p0220": {
    "es": "Motivo de la decisión",
    "en": "Decision reason"
  },
  "p0221": {
    "es": "Aprobar etapa revisada",
    "en": "Approve reviewed stage"
  },
  "p0222": {
    "es": "Rechazar etapa revisada",
    "en": "Reject reviewed stage"
  },
  "p0223": {
    "es": "Publicación vigente",
    "en": "Current publication"
  },
  "p0224": {
    "es": "Consultar publicación vigente",
    "en": "Retrieve current publication"
  },
  "p0225": {
    "es": "Versión publicada",
    "en": "Published version"
  },
  "p0226": {
    "es": "Ver catálogo público",
    "en": "View public catalog"
  },
  "p0227": {
    "es": "Todavía no hay una publicación.",
    "en": "There is no publication yet."
  },
  "p0228": {
    "es": "Publicar esta versión reemplaza la vigente. Para restaurar una anterior, consultá su borrador aprobado y publicalo como una nueva versión.",
    "en": "Publishing this version replaces the current one. To restore an earlier version, retrieve its approved draft and publish it as a new version."
  },
  "p0229": {
    "es": "Motivo de publicación",
    "en": "Publication reason"
  },
  "p0230": {
    "es": "Publicar versión revisada",
    "en": "Publish reviewed version"
  },
  "p0231": {
    "es": "Estamos comprobando la cancelación. Esperá el resultado.",
    "en": "We are checking the cancellation. Wait for the result."
  },
  "p0232": {
    "es": "Turno cancelado y auditado. Actualizando…",
    "en": "Appointment cancelled and audited. Updating…"
  },
  "p0233": {
    "es": "No pudimos comprobar el resultado. La cancelación podría haberse registrado. Volvé a consultar el estado antes de intentar otra acción.",
    "en": "We could not verify the result. The cancellation may have been recorded. Check the status again before taking another action."
  },
  "p0234": {
    "es": "Volver a consultar mis turnos",
    "en": "Check my appointments again"
  },
  "p0235": {
    "es": "No hay turnos disponibles para tu cuenta en esta organización.",
    "en": "There are no appointments available for your account in this organization."
  },
  "p0236": {
    "es": "Cancelar mi turno",
    "en": "Cancel my appointment"
  },
  "p0237": {
    "es": "El pago todavía no está disponible para este pedido.",
    "en": "Payment is not yet available for this order."
  },
  "p0238": {
    "es": "No pudimos abrir el pago. Volvé a intentarlo más tarde.",
    "en": "We could not open payment. Try again later."
  },
  "p0239": {
    "es": "Abriendo pago…",
    "en": "Opening payment…"
  },
  "p0240": {
    "es": "Continuar con el pago",
    "en": "Continue to payment"
  },
  "p0241": {
    "es": "Recepción registrada con identidad, serie, versión y evidencia del servidor. Actualizando…",
    "en": "Receipt recorded with identity, serial number, version, and server evidence. Updating…"
  },
  "p0242": {
    "es": "No pudimos comprobar el resultado. La recepción pudo haberse registrado. Consultá el estado antes de volver a actuar.",
    "en": "We could not verify the result. Receipt may have been recorded. Check the status before taking another action."
  },
  "p0243": {
    "es": "Discrepancia registrada con tu identidad y la versión exacta. Actualizando…",
    "en": "Discrepancy recorded with your identity and the exact version. Updating…"
  },
  "p0244": {
    "es": "No pudimos comprobar el resultado. La discrepancia pudo haberse registrado. Consultá el estado antes de volver a actuar.",
    "en": "We could not verify the result. The discrepancy may have been recorded. Check the status before taking another action."
  },
  "p0245": {
    "es": "Consultar estado de entregas",
    "en": "Check delivery status"
  },
  "p0246": {
    "es": "Discrepancias de entrega",
    "en": "Delivery discrepancies"
  },
  "p0247": {
    "es": "Estado:",
    "en": "Status:"
  },
  "p0248": {
    "es": " · resolución: {resolution_action}",
    "en": " · resolution: {resolution_action}"
  },
  "p0249": {
    "es": "Nueva preparación:",
    "en": "New preparation:"
  },
  "p0250": {
    "es": "Autorización:",
    "en": "Authorization:"
  },
  "p0251": {
    "es": "Pedido",
    "en": "Order"
  },
  "p0252": {
    "es": "Preparación verificada",
    "en": "Verified preparation"
  },
  "p0253": {
    "es": "Checklist de entrega",
    "en": "Delivery checklist"
  },
  "p0254": {
    "es": "· v",
    "en": "· v"
  },
  "p0255": {
    "es": "Completado por la franquicia: {completed_at}",
    "en": "Completed by the franchise: {completed_at}"
  },
  "p0256": {
    "es": "Fecha pendiente de validación",
    "en": "Date awaiting validation"
  },
  "p0257": {
    "es": "La preparación todavía no está completa.",
    "en": "Preparation is not yet complete."
  },
  "p0258": {
    "es": "La franquicia todavía no vinculó una versión de preparación.",
    "en": "The franchise has not yet linked a preparation version."
  },
  "p0259": {
    "es": "Número de serie observado",
    "en": "Observed serial number"
  },
  "p0260": {
    "es": "Confirmo que recibí el activo identificado y revisé la preparación indicada",
    "en": "I confirm that I received the identified asset and reviewed the specified preparation"
  },
  "p0261": {
    "es": "Registrar recepción",
    "en": "Record receipt"
  },
  "p0262": {
    "es": "Informar una discrepancia",
    "en": "Report a discrepancy"
  },
  "p0263": {
    "es": "Código de motivo",
    "en": "Reason code"
  },
  "p0264": {
    "es": "Código aprobado por la empresa",
    "en": "Company-approved code"
  },
  "p0265": {
    "es": "Detalle",
    "en": "Details"
  },
  "p0266": {
    "es": "Rechazar esta presentación",
    "en": "Reject this handover"
  },
  "p0267": {
    "es": "Aceptado: {accepted_at}",
    "en": "Accepted: {accepted_at}"
  },
  "p0268": {
    "es": "Sin aceptación disponible hasta completar la preparación exacta.",
    "en": "Acceptance is unavailable until the exact preparation is completed."
  },
  "p0269": {
    "es": "La vigencia terminó. Actualizá el estado para consultar una cotización vigente.",
    "en": "The validity period has ended. Refresh the status to find a current quote."
  },
  "p0270": {
    "es": "Estamos comprobando la aceptación. No cierres esta pantalla.",
    "en": "We are checking acceptance. Do not close this screen."
  },
  "p0271": {
    "es": "Aceptación confirmada por el servidor. Consultando el pedido…",
    "en": "Acceptance confirmed by the server. Retrieving the order…"
  },
  "p0272": {
    "es": "No pudimos comprobar el resultado. El pedido podría haberse creado. Actualizá el estado antes de intentar otra acción.",
    "en": "We could not verify the result. The order may have been created. Refresh the status before taking another action."
  },
  "p0273": {
    "es": "Actualizar estado",
    "en": "Refresh status"
  },
  "p0274": {
    "es": "No hay cotizaciones disponibles para tu cuenta en esta organización.",
    "en": "There are no quotes available for your account in this organization."
  },
  "p0275": {
    "es": "Importe:",
    "en": "Amount:"
  },
  "p0276": {
    "es": "Vigencia:",
    "en": "Valid until:"
  },
  "p0277": {
    "es": "Pedido creado:",
    "en": "Order created:"
  },
  "p0278": {
    "es": "Aceptar y crear pedido",
    "en": "Accept and create order"
  },
  "p0279": {
    "es": "Ayuda para aceptar una cotización",
    "en": "Help with accepting a quote"
  },
  "p0280": {
    "es": "/",
    "en": "/"
  },
  "p0281": {
    "es": "Tu respuesta quedó guardada.",
    "en": "Your response has been saved."
  },
  "p0282": {
    "es": "No pudimos confirmar el resultado. Consultá la respuesta guardada antes de realizar otra acción.",
    "en": "We could not confirm the result. Check the saved response before taking another action."
  },
  "p0283": {
    "es": "Todavía no hay una respuesta disponible. Podés volver a consultar.",
    "en": "No response is available yet. You can check again."
  },
  "p0284": {
    "es": "Recuperamos tu respuesta guardada.",
    "en": "We recovered your saved response."
  },
  "p0285": {
    "es": "No pudimos consultar tu respuesta. Volvé a consultar cuando se restablezca la conexión.",
    "en": "We could not retrieve your response. Check again when the connection is restored."
  },
  "p0286": {
    "es": "Respuesta guardada:",
    "en": "Saved response:"
  },
  "p0287": {
    "es": "de 10.",
    "en": "out of 10."
  },
  "p0288": {
    "es": "No se enviará una nueva respuesta.",
    "en": "No new response will be sent."
  },
  "p0289": {
    "es": "Puntuación",
    "en": "Score"
  },
  "p0290": {
    "es": "Elegí una puntuación",
    "en": "Choose a score"
  },
  "p0291": {
    "es": "Leí y acepto el aviso mostrado para esta encuesta.",
    "en": "I have read and accept the notice shown for this survey."
  },
  "p0292": {
    "es": "Procesando…",
    "en": "Processing…"
  },
  "p0293": {
    "es": "Enviar respuesta",
    "en": "Send response"
  },
  "p0294": {
    "es": "Esta encuesta no admite nuevas respuestas.",
    "en": "This survey is not accepting new responses."
  },
  "p0295": {
    "es": "Consultar respuesta guardada",
    "en": "Check saved response"
  },
  "p0296": {
    "es": "Ayuda con esta encuesta",
    "en": "Help with this survey"
  },
  "p0297": {
    "es": "Elegí una puntuación de 0 a 10 y revisá el aviso antes de enviar. Se guarda una respuesta por persona y encuesta.",
    "en": "Choose a score from 0 to 10 and review the notice before sending. One response is saved per person and survey."
  },
  "p0298": {
    "es": "Si se interrumpe la conexión, consultá la respuesta guardada. La consulta no vuelve a enviar el formulario.",
    "en": "If the connection is interrupted, check the saved response. The lookup does not resubmit the form."
  },
  "p0299": {
    "es": "No se pudo consultar el registro local. Consultá el estado antes de avanzar.",
    "en": "The local record could not be read. Check the status before proceeding."
  },
  "p0300": {
    "es": "El servidor aceptó el avance. Consultá el estado actual.",
    "en": "The server accepted the update. Check the current status."
  },
  "p0301": {
    "es": "Estado actual consultado. No identifica la solicitud que lo produjo.",
    "en": "Current status retrieved. It does not identify the request that produced it."
  },
  "p0302": {
    "es": "La respuesta no quedó confirmada. Consultá el estado antes de volver a operar.",
    "en": "The response is unconfirmed. Check the status before proceeding again."
  },
  "p0303": {
    "es": "No se pudo consultar el estado actual.",
    "en": "The current status could not be retrieved."
  },
  "p0304": {
    "es": "Estado actual:",
    "en": "Current status:"
  },
  "p0305": {
    "es": "Hay una respuesta pendiente de comprobar.",
    "en": "A response is awaiting verification."
  },
  "p0306": {
    "es": "Registrar inicio de ensamblaje",
    "en": "Record assembly start"
  },
  "p0307": {
    "es": "Consultar estado actual",
    "en": "Check current status"
  },
  "p0308": {
    "es": "Registrando la solicitud local…",
    "en": "Recording the local request…"
  },
  "p0309": {
    "es": "No se envió la solicitud: el navegador no pudo conservar su clave. Consultá soporte antes de continuar.",
    "en": "The request was not sent: the browser could not retain its key. Contact support before continuing."
  },
  "p0310": {
    "es": "Solicitud registrada: {request_id}. Esto no confirma un cobro. Actualizá los pedidos para consultar su estado.",
    "en": "Request recorded: {request_id}. This does not confirm a charge. Refresh the orders to check its status."
  },
  "p0311": {
    "es": "Ya existe una solicitud o el pedido cambió. Actualizá los pedidos; no generes otra clave ni otro pedido.",
    "en": "A request already exists or the order has changed. Refresh the orders; do not generate another key or order."
  },
  "p0312": {
    "es": "Tu sesión no permite solicitar este pago. Consultá los accesos con el responsable.",
    "en": "Your session does not allow you to request this payment. Ask the responsible person to review your access."
  },
  "p0313": {
    "es": "No pudimos comprobar la solicitud. Actualizá los pedidos antes de continuar. Conservamos la clave de esta sesión; no borres los datos del navegador para forzar un reintento.",
    "en": "We could not verify the request. Refresh the orders before continuing. We kept this session's key; do not clear browser data to force a retry."
  },
  "p0314": {
    "es": "Registrar solicitud de pago",
    "en": "Record payment request"
  },
  "p0315": {
    "es": "Solicitud local para",
    "en": "Local request for"
  },
  "p0316": {
    "es": "; el servidor toma el total del pedido. No se ejecuta un cobro desde esta pantalla.",
    "en": "; the server uses the order total. This screen does not execute a charge."
  },
  "p0317": {
    "es": "Seleccioná una unidad de la lista disponible.",
    "en": "Select a unit from the available list."
  },
  "p0318": {
    "es": "Consultando el resultado de la reserva…",
    "en": "Checking the reservation result…"
  },
  "p0319": {
    "es": "Reserva guardada. Actualizá los pedidos para consultar la unidad asignada.",
    "en": "Reservation saved. Refresh the orders to check the allocated unit."
  },
  "p0320": {
    "es": "El pedido o la unidad cambió. Actualizá los pedidos antes de volver a actuar.",
    "en": "The order or unit has changed. Refresh the orders before taking another action."
  },
  "p0321": {
    "es": "Tu sesión no permite esta reserva. Consultá los accesos con el responsable.",
    "en": "Your session does not allow this reservation. Ask the responsible person to review your access."
  },
  "p0322": {
    "es": "No pudimos comprobar el resultado. No repitas la reserva: actualizá los pedidos para consultar el estado guardado.",
    "en": "We could not verify the result. Do not repeat the reservation: refresh the orders to check the saved status."
  },
  "p0323": {
    "es": "Pedido recibido",
    "en": "Order received"
  },
  "p0324": {
    "es": "Confirmado",
    "en": "Confirmed"
  },
  "p0325": {
    "es": "Entregado",
    "en": "Delivered"
  },
  "p0326": {
    "es": "Cancelado",
    "en": "Cancelled"
  },
  "p0327": {
    "es": "Solicitud creada: no confirma cobro",
    "en": "Request created: charge not confirmed"
  },
  "p0328": {
    "es": "Pendiente de confirmación",
    "en": "Awaiting confirmation"
  },
  "p0329": {
    "es": "Autorizado: no confirma cobro",
    "en": "Authorized: charge not confirmed"
  },
  "p0330": {
    "es": "Cobro registrado",
    "en": "Charge recorded"
  },
  "p0331": {
    "es": "Fallido",
    "en": "Failed"
  },
  "p0332": {
    "es": "Reembolso registrado",
    "en": "Refund recorded"
  },
  "p0333": {
    "es": "En disputa",
    "en": "Disputed"
  },
  "p0334": {
    "es": "Entrega preparada",
    "en": "Delivery prepared"
  },
  "p0335": {
    "es": "Presentada al cliente",
    "en": "Presented to customer"
  },
  "p0336": {
    "es": "Aceptada por el cliente",
    "en": "Accepted by customer"
  },
  "p0337": {
    "es": "Rechazada",
    "en": "Rejected"
  },
  "p0338": {
    "es": "Estado no reconocido: consultar soporte",
    "en": "Unrecognized status: contact support"
  },
  "p0339": {
    "es": "Pedidos: stock, pago y entrega",
    "en": "Orders: stock, payment, and delivery"
  },
  "p0340": {
    "es": "Los estados se consultan en el registro empresarial. Reservar una unidad no cobra ni confirma su entrega.",
    "en": "Statuses are read from the business record. Reserving a unit does not charge for it or confirm its delivery."
  },
  "p0341": {
    "es": "Actualizar pedidos",
    "en": "Refresh orders"
  },
  "p0342": {
    "es": "Resultado de reserva",
    "en": "Reservation result"
  },
  "p0343": {
    "es": "No pudimos consultar los pedidos. Las acciones están bloqueadas; volvé a consultar.",
    "en": "We could not retrieve the orders. Actions are blocked; check again."
  },
  "p0344": {
    "es": "Vista parcial: supera 25 pedidos, 100 líneas o estados por pedido, o 200 unidades. Las reservas están bloqueadas; solicitá una consulta acotada al responsable.",
    "en": "Partial view: exceeds 25 orders, 100 lines or statuses per order, or 200 units. Reservations are blocked; ask the responsible person for a narrower query."
  },
  "p0345": {
    "es": "No hay pedidos en esta organización.",
    "en": "There are no orders in this organization."
  },
  "p0346": {
    "es": "· cantidad",
    "en": "· quantity"
  },
  "p0347": {
    "es": "Unidad reservada:",
    "en": "Reserved unit:"
  },
  "p0348": {
    "es": "Sin unidad reservada.",
    "en": "No unit reserved."
  },
  "p0349": {
    "es": "Unidad disponible ·",
    "en": "Available unit ·"
  },
  "p0350": {
    "es": "Elegir número de serie",
    "en": "Choose serial number"
  },
  "p0351": {
    "es": "No hay unidades disponibles de esta variante. Actualizá después de la recepción de stock.",
    "en": "No units of this variant are available. Refresh after stock is received."
  },
  "p0352": {
    "es": "Reservar unidad",
    "en": "Reserve unit"
  },
  "p0353": {
    "es": "Sin stock disponible",
    "en": "No stock available"
  },
  "p0354": {
    "es": "Reserva no habilitada para tu permiso, cantidad o estado actual.",
    "en": "Reservation is not enabled for your permission, quantity, or current state."
  },
  "p0355": {
    "es": "Pago:",
    "en": "Payment:"
  },
  "p0356": {
    "es": "Sin solicitud registrada. No hay cobro confirmado.",
    "en": "No request recorded. No charge is confirmed."
  },
  "p0357": {
    "es": "Solicitud de pago desactivada: falta seleccionar y configurar el proveedor.",
    "en": "Payment requests are disabled: a provider must be selected and configured."
  },
  "p0358": {
    "es": "Entrega:",
    "en": "Delivery:"
  },
  "p0359": {
    "es": "Sin entrega preparada.",
    "en": "No delivery prepared."
  },
  "p0360": {
    "es": "Turno confirmado. El estado quedó guardado.",
    "en": "Appointment confirmed. The status has been saved."
  },
  "p0361": {
    "es": "Cambio guardado. Actualizando agenda.",
    "en": "Change saved. Refreshing schedule."
  },
  "p0362": {
    "es": "No pudimos comprobar el resultado. Actualizá la agenda para consultar el estado guardado antes de volver a actuar.",
    "en": "We could not verify the result. Refresh the schedule to check the saved status before taking another action."
  },
  "p0363": {
    "es": "Solicitado",
    "en": "Requested"
  },
  "p0364": {
    "es": "Completado",
    "en": "Completed"
  },
  "p0365": {
    "es": "Ausente",
    "en": "No-show"
  },
  "p0366": {
    "es": "Consulta",
    "en": "Consultation"
  },
  "p0367": {
    "es": "Prueba de manejo",
    "en": "Test drive"
  },
  "p0368": {
    "es": "Entrega",
    "en": "Delivery"
  },
  "p0369": {
    "es": "Agenda de turnos",
    "en": "Appointment schedule"
  },
  "p0370": {
    "es": "Elegí el recurso por su nombre. El servidor valida jornada, habilidad y conflictos al asignar. Los horarios se muestran en UTC.",
    "en": "Choose the resource by name. The server checks working hours, skills, and conflicts when assigning it. Times are shown in UTC."
  },
  "p0371": {
    "es": "Actualizar agenda",
    "en": "Refresh schedule"
  },
  "p0372": {
    "es": "Resultado de la agenda",
    "en": "Schedule result"
  },
  "p0373": {
    "es": "La consulta supera el límite de 200 turnos o recursos. No se muestra una agenda completa y las acciones están bloqueadas; acotá el rango o solicitá una consulta paginada.",
    "en": "The query exceeds the limit of 200 appointments or resources. The complete schedule is not shown and actions are blocked; narrow the range or request a paginated query."
  },
  "p0374": {
    "es": "No hay turnos para esta fecha.",
    "en": "There are no appointments for this date."
  },
  "p0375": {
    "es": "UTC",
    "en": "UTC"
  },
  "p0376": {
    "es": "Recurso:",
    "en": "Resource:"
  },
  "p0377": {
    "es": "Recurso asignado",
    "en": "Resource assigned"
  },
  "p0378": {
    "es": "Sin asignar",
    "en": "Unassigned"
  },
  "p0379": {
    "es": "Recurso disponible para evaluar",
    "en": "Available resource to evaluate"
  },
  "p0380": {
    "es": "Seleccionar recurso",
    "en": "Select resource"
  },
  "p0381": {
    "es": "Asignar recurso",
    "en": "Assign resource"
  },
  "p0382": {
    "es": "Confirmar turno",
    "en": "Confirm appointment"
  },
  "p0383": {
    "es": "Otra acción",
    "en": "Another action"
  },
  "p0384": {
    "es": "Cancelar",
    "en": "Cancel"
  },
  "p0385": {
    "es": "Completar",
    "en": "Complete"
  },
  "p0386": {
    "es": "Registrar ausencia",
    "en": "Record no-show"
  },
  "p0387": {
    "es": "Motivo para cancelar o registrar ausencia",
    "en": "Reason for cancellation or no-show"
  },
  "p0388": {
    "es": "Aplicar acción",
    "en": "Apply action"
  },
  "p0389": {
    "es": "Consulta recuperada: el intervalo figura cancelado. No reenviamos la acción.",
    "en": "Lookup recovered: the interval is recorded as cancelled. We did not resend the action."
  },
  "p0390": {
    "es": "Hay una cancelación pendiente de comprobar. Consultá el intervalo antes de continuar.",
    "en": "A cancellation is awaiting verification. Check the interval before continuing."
  },
  "p0391": {
    "es": "No pudimos leer la referencia de recuperación. No enviaremos acciones; consultá soporte sin borrar los datos del navegador.",
    "en": "We could not read the recovery reference. No actions will be sent; contact support without clearing browser data."
  },
  "p0392": {
    "es": "No se envió: no pudimos conservar la referencia de recuperación. Consultá soporte antes de continuar.",
    "en": "Not sent: we could not retain the recovery reference. Contact support before continuing."
  },
  "p0393": {
    "es": "Registrando cancelación…",
    "en": "Recording cancellation…"
  },
  "p0394": {
    "es": "Cancelación registrada. Consultá el intervalo para ver su estado guardado.",
    "en": "Cancellation recorded. Check the interval to see its saved status."
  },
  "p0395": {
    "es": "Tu sesión no autoriza esta cancelación. Revisá tus accesos y consultá el intervalo.",
    "en": "Your session does not authorize this cancellation. Review your access and check the interval."
  },
  "p0396": {
    "es": "El intervalo cambió o tiene restricciones, por ejemplo citas activas. Consultá su estado y pedí revisión; no repetimos la acción.",
    "en": "The interval has changed or has restrictions, such as active appointments. Check its status and request a review; the action is not repeated."
  },
  "p0397": {
    "es": "No pudimos comprobar la cancelación. Puede haberse registrado; consultá el intervalo antes de continuar.",
    "en": "We could not verify the cancellation. It may have been recorded; check the interval before continuing."
  },
  "p0398": {
    "es": "Activo",
    "en": "Active"
  },
  "p0399": {
    "es": "Cancelar intervalo",
    "en": "Cancel interval"
  },
  "p0400": {
    "es": "Resultado de cancelación de intervalo",
    "en": "Interval cancellation result"
  },
  "p0401": {
    "es": "Consultar intervalo",
    "en": "Check interval"
  },
  "p0402": {
    "es": "Consulta recuperada: el lead tiene una versión posterior. Revisá el responsable y el estado actuales antes de continuar.",
    "en": "Lookup recovered: the lead has a newer version. Review the current owner and status before continuing."
  },
  "p0403": {
    "es": "Hay un cambio pendiente de comprobar. Consultá el lead; no repitas la acción ni borres la referencia.",
    "en": "A change is awaiting verification. Check the lead; do not repeat the action or delete the reference."
  },
  "p0404": {
    "es": "Registrando cambio…",
    "en": "Recording change…"
  },
  "p0405": {
    "es": "Cambio registrado. Consultá el lead para trabajar sobre su estado actual.",
    "en": "Change recorded. Retrieve the lead to work with its current state."
  },
  "p0406": {
    "es": "La sesión no autoriza este cambio. Revisá tus accesos y consultá el lead; conservamos la referencia.",
    "en": "This session does not authorize the change. Review your access and check the lead; we kept the reference."
  },
  "p0407": {
    "es": "El lead cambió o la operación no admite esa versión. Consultá su estado antes de continuar.",
    "en": "The lead has changed or the operation does not accept that version. Check its status before continuing."
  },
  "p0408": {
    "es": "No pudimos comprobar el resultado. El cambio puede haberse registrado; consultá el lead antes de continuar.",
    "en": "We could not verify the result. The change may have been recorded; check the lead before continuing."
  },
  "p0409": {
    "es": "Responsable",
    "en": "Owner"
  },
  "p0410": {
    "es": "Asignar",
    "en": "Assign"
  },
  "p0411": {
    "es": "Nuevo estado",
    "en": "New status"
  },
  "p0412": {
    "es": "Cambiar estado",
    "en": "Change status"
  },
  "p0413": {
    "es": "Resultado del lead",
    "en": "Lead result"
  },
  "p0414": {
    "es": "Consultar lead",
    "en": "Check lead"
  },
  "p0415": {
    "es": "Consulta recuperada: la discrepancia figura resuelta. No reenviamos la acción.",
    "en": "Lookup recovered: the discrepancy is recorded as resolved. We did not resend the action."
  },
  "p0416": {
    "es": "La discrepancia figura resuelta con otra decisión. Revisá el resultado con el responsable.",
    "en": "The discrepancy was resolved with a different decision. Review the result with the responsible person."
  },
  "p0417": {
    "es": "Hay una resolución pendiente de comprobar en esta sesión. Consultá su estado; no repitas la acción.",
    "en": "A resolution is awaiting verification in this session. Check its status; do not repeat the action."
  },
  "p0418": {
    "es": "Registrando resolución…",
    "en": "Recording resolution…"
  },
  "p0419": {
    "es": "Resolución registrada. Consultá el estado para continuar; esto no confirma entrega, reembolso ni cambio completados.",
    "en": "Resolution recorded. Check the status to continue; this does not confirm a completed delivery, refund, or exchange."
  },
  "p0420": {
    "es": "La sesión no autoriza esta resolución. Conservamos la referencia; revisá tus accesos y consultá el estado.",
    "en": "This session does not authorize the resolution. We kept the reference; review your access and check the status."
  },
  "p0421": {
    "es": "La discrepancia cambió o ya fue resuelta. Consultá el estado antes de continuar.",
    "en": "The discrepancy has changed or has already been resolved. Check the status before continuing."
  },
  "p0422": {
    "es": "No pudimos comprobar el resultado. La resolución puede haberse registrado; consultá su estado antes de continuar.",
    "en": "We could not verify the result. The resolution may have been recorded; check its status before continuing."
  },
  "p0423": {
    "es": "Resolución",
    "en": "Resolution"
  },
  "p0424": {
    "es": "Corregir y volver a presentar",
    "en": "Correct and present again"
  },
  "p0425": {
    "es": "Autorizar devolución",
    "en": "Authorize return"
  },
  "p0426": {
    "es": "Autorizar cambio",
    "en": "Authorize exchange"
  },
  "p0427": {
    "es": "Notas",
    "en": "Notes"
  },
  "p0428": {
    "es": "Registrar resolución",
    "en": "Record resolution"
  },
  "p0429": {
    "es": "Resolución:",
    "en": "Resolution:"
  },
  "p0430": {
    "es": "Resultado de resolución",
    "en": "Resolution result"
  },
  "p0431": {
    "es": "Consultar resolución",
    "en": "Check resolution"
  },
  "p0432": {
    "es": "Jornadas y ausencias",
    "en": "Working hours and absences"
  },
  "p0433": {
    "es": "Sin jornada explícita el backend no publica cupos ni asigna recursos.",
    "en": "Without explicit working hours, the backend does not publish capacity or assign resources."
  },
  "p0434": {
    "es": "Franquicia",
    "en": "Franchise"
  },
  "p0435": {
    "es": "—",
    "en": "—"
  },
  "p0436": {
    "es": "Motivo:",
    "en": "Reason:"
  },
  "p0437": {
    "es": "Recursos y turnos",
    "en": "Resources and appointments"
  },
  "p0438": {
    "es": "Publicar capacidad",
    "en": "Publish capacity"
  },
  "p0439": {
    "es": "Preparación versionada de entregas",
    "en": "Versioned delivery preparation"
  },
  "p0440": {
    "es": "Publicá una versión inmutable. Para confirmaciones, la respuesta operativa exacta es",
    "en": "Publish an immutable version. For confirmations, the exact operational response is"
  },
  "p0441": {
    "es": "; el servidor impide presentar la entrega si falta un ítem obligatorio.",
    "en": "; the server prevents presenting the delivery if a required item is missing."
  },
  "p0442": {
    "es": "· cliente",
    "en": "· customer"
  },
  "p0443": {
    "es": "· estado",
    "en": "· status"
  },
  "p0444": {
    "es": "No hay discrepancias registradas.",
    "en": "No discrepancies are recorded."
  },
  "p0445": {
    "es": "Hay una cotización pendiente de comprobar. Consultá el resultado sin repetir la emisión.",
    "en": "A quote is awaiting verification. Check the result without issuing it again."
  },
  "p0446": {
    "es": "No pudimos leer la referencia. Consultá soporte sin borrar los datos del navegador.",
    "en": "We could not read the reference. Contact support without clearing browser data."
  },
  "p0447": {
    "es": "No se envió: no pudimos conservar la referencia. Consultá soporte antes de continuar.",
    "en": "Not sent: we could not retain the reference. Contact support before continuing."
  },
  "p0448": {
    "es": "Emitiendo cotización…",
    "en": "Issuing quote…"
  },
  "p0449": {
    "es": "Cotización registrada. Consultá el resultado persistente antes de continuar.",
    "en": "Quote recorded. Check the persistent result before continuing."
  },
  "p0450": {
    "es": "No pudimos comprobar el resultado. La cotización puede haberse registrado; consultá sin volver a emitir.",
    "en": "We could not verify the result. The quote may have been recorded; check without issuing it again."
  },
  "p0451": {
    "es": "Consulta recuperada: esta referencia corresponde a la cotización mostrada. La emisión permanece cerrada para evitar duplicados.",
    "en": "Recovered result: this reference matches the displayed quote. Issuing is still locked to prevent duplicates."
  },
  "p0452": {
    "es": "No pudimos recuperar la cotización. Conservamos la referencia y el bloqueo; consultá nuevamente o pedí revisión autorizada.",
    "en": "We could not retrieve the quote. The reference and lock are retained; check again or request an authorized review."
  },
  "p0453": {
    "es": "Lista de precios",
    "en": "Price list"
  },
  "p0454": {
    "es": "Válida hasta",
    "en": "Valid until"
  },
  "p0455": {
    "es": "Emitir cotización",
    "en": "Issue quote"
  },
  "p0456": {
    "es": "Resultado de cotización",
    "en": "Quote result"
  },
  "p0457": {
    "es": "Cotización",
    "en": "Quote"
  },
  "p0458": {
    "es": "Consultar cotización",
    "en": "Check quote"
  },
  "p0459": {
    "es": "Hay un intervalo pendiente de comprobar. Consultá el resultado sin volver a registrarlo.",
    "en": "An interval is awaiting verification. Check the result without registering it again."
  },
  "p0460": {
    "es": "El intervalo no es válido.",
    "en": "The interval is invalid."
  },
  "p0461": {
    "es": "Registrando intervalo…",
    "en": "Registering interval…"
  },
  "p0462": {
    "es": "Intervalo registrado. Consultá el resultado antes de preparar otro.",
    "en": "Interval registered. Check the result before preparing another."
  },
  "p0463": {
    "es": "No pudimos comprobar el resultado. El intervalo puede haberse registrado; consultá sin volver a enviarlo.",
    "en": "We could not verify the result. The interval may have been registered; check without submitting it again."
  },
  "p0464": {
    "es": "Consulta recuperada: esta referencia corresponde al intervalo mostrado. Revisá su estado antes de preparar otro.",
    "en": "Recovered result: this reference matches the displayed interval. Review its status before preparing another."
  },
  "p0465": {
    "es": "No pudimos recuperar el intervalo. Conservamos la referencia y el bloqueo; consultá nuevamente o pedí revisión autorizada.",
    "en": "We could not retrieve the interval. The reference and lock are retained; check again or request an authorized review."
  },
  "p0466": {
    "es": "Formulario nuevo preparado. El intervalo anterior permanece registrado.",
    "en": "New form ready. The previous interval remains registered."
  },
  "p0467": {
    "es": "No pudimos actualizar la referencia. Conservamos el bloqueo; consultá soporte.",
    "en": "We could not update the reference. The lock is retained; contact support."
  },
  "p0468": {
    "es": "Alcance",
    "en": "Scope"
  },
  "p0469": {
    "es": "Vacío = toda la franquicia",
    "en": "Empty = entire franchise"
  },
  "p0470": {
    "es": "Tipo",
    "en": "Type"
  },
  "p0471": {
    "es": "Jornada activa",
    "en": "Working hours"
  },
  "p0472": {
    "es": "No disponible",
    "en": "Unavailable"
  },
  "p0473": {
    "es": "Motivo",
    "en": "Reason"
  },
  "p0474": {
    "es": "Obligatorio para ausencia",
    "en": "Required for absence"
  },
  "p0475": {
    "es": "Desde",
    "en": "From"
  },
  "p0476": {
    "es": "Hasta",
    "en": "Until"
  },
  "p0477": {
    "es": "Registrar intervalo",
    "en": "Register interval"
  },
  "p0478": {
    "es": "Resultado de intervalo",
    "en": "Interval result"
  },
  "p0479": {
    "es": "Intervalo",
    "en": "Interval"
  },
  "p0480": {
    "es": "Consultar intervalo creado",
    "en": "Check created interval"
  },
  "p0481": {
    "es": "Preparar otro intervalo",
    "en": "Prepare another interval"
  },
  "p0482": {
    "es": "Hay un recurso pendiente de comprobar. Consultá el resultado sin volver a registrarlo.",
    "en": "A resource is awaiting verification. Check the result without registering it again."
  },
  "p0483": {
    "es": "Indicá nombre y habilidades. Empleados y contratistas requieren identidad; bahías y vehículos no la usan.",
    "en": "Enter a name and skills. Employees and contractors require an identity; bays and vehicles do not use one."
  },
  "p0484": {
    "es": "Registrando recurso…",
    "en": "Registering resource…"
  },
  "p0485": {
    "es": "Recurso registrado. Consultá el resultado antes de preparar otro.",
    "en": "Resource registered. Check the result before preparing another."
  },
  "p0486": {
    "es": "No pudimos comprobar el resultado. El recurso puede haberse registrado; consultá sin volver a enviarlo.",
    "en": "We could not verify the result. The resource may have been registered; check without submitting it again."
  },
  "p0487": {
    "es": "Consulta recuperada: esta referencia corresponde al recurso mostrado. Revisá su estado antes de preparar otro.",
    "en": "Recovered result: this reference matches the displayed resource. Review its status before preparing another."
  },
  "p0488": {
    "es": "No pudimos recuperar el recurso. Conservamos la referencia y el bloqueo; consultá nuevamente o pedí revisión autorizada.",
    "en": "We could not retrieve the resource. The reference and lock are retained; check again or request an authorized review."
  },
  "p0489": {
    "es": "Formulario nuevo preparado. El recurso anterior permanece registrado.",
    "en": "New form ready. The previous resource remains registered."
  },
  "p0490": {
    "es": "Nombre",
    "en": "Name"
  },
  "p0491": {
    "es": "Identidad del proveedor de acceso",
    "en": "Access provider identity"
  },
  "p0492": {
    "es": "Empleado",
    "en": "Employee"
  },
  "p0493": {
    "es": "Contratista",
    "en": "Contractor"
  },
  "p0494": {
    "es": "Bahía",
    "en": "Bay"
  },
  "p0495": {
    "es": "Vehículo",
    "en": "Vehicle"
  },
  "p0496": {
    "es": "Habilidades",
    "en": "Skills"
  },
  "p0497": {
    "es": "Crear recurso",
    "en": "Create resource"
  },
  "p0498": {
    "es": "Resultado de recurso",
    "en": "Resource result"
  },
  "p0499": {
    "es": "Recurso",
    "en": "Resource"
  },
  "p0500": {
    "es": "Consultar recurso creado",
    "en": "Check created resource"
  },
  "p0501": {
    "es": "Preparar otro recurso",
    "en": "Prepare another resource"
  },
  "p0502": {
    "es": "Hay un turno pendiente de comprobar. Consultá el resultado sin volver a registrarlo.",
    "en": "A slot is awaiting verification. Check the result without registering it again."
  },
  "p0503": {
    "es": "El turno no es válido.",
    "en": "The slot is invalid."
  },
  "p0504": {
    "es": "La capacidad debe ser un entero entre 1 y 100.",
    "en": "Capacity must be an integer from 1 to 100."
  },
  "p0505": {
    "es": "Registrando turno…",
    "en": "Registering slot…"
  },
  "p0506": {
    "es": "Turno registrado. Consultá el resultado antes de preparar otro.",
    "en": "Slot registered. Check the result before preparing another."
  },
  "p0507": {
    "es": "No pudimos comprobar el resultado. El turno puede haberse registrado; consultá sin volver a enviarlo.",
    "en": "We could not verify the result. The slot may have been registered; check without submitting it again."
  },
  "p0508": {
    "es": "Consulta recuperada: esta referencia corresponde al turno mostrado. Revisá su estado antes de preparar otro.",
    "en": "Recovered result: this reference matches the displayed slot. Review its status before preparing another."
  },
  "p0509": {
    "es": "No pudimos recuperar el turno. Conservamos la referencia y el bloqueo; consultá nuevamente o pedí revisión autorizada.",
    "en": "We could not retrieve the slot. The reference and lock are retained; check again or request an authorized review."
  },
  "p0510": {
    "es": "Formulario nuevo preparado. El turno anterior permanece registrado.",
    "en": "New form ready. The previous slot remains registered."
  },
  "p0511": {
    "es": "Cupos",
    "en": "Capacity"
  },
  "p0512": {
    "es": "Publicar turno",
    "en": "Publish slot"
  },
  "p0513": {
    "es": "Resultado de turno",
    "en": "Slot result"
  },
  "p0514": {
    "es": "ocupados ·",
    "en": "occupied ·"
  },
  "p0515": {
    "es": "Consultar turno creado",
    "en": "Check created slot"
  },
  "p0516": {
    "es": "Preparar otro turno",
    "en": "Prepare another slot"
  },
  "p0517": {
    "es": "Hay una versión pendiente de comprobar. Consultá su contenido sin volver a publicarla.",
    "en": "A version is awaiting verification. Check its content without publishing it again."
  },
  "p0518": {
    "es": "No pudimos leer la referencia. Pedí revisión sin borrar los datos del navegador.",
    "en": "We could not read the reference. Request a review without clearing browser data."
  },
  "p0519": {
    "es": "Revisá ID, versión, título y contenido. Cada ítem debe tener un ID distinto y un texto válido.",
    "en": "Review the ID, version, title and content. Each item must have a distinct ID and valid text."
  },
  "p0520": {
    "es": "No se envió: no pudimos conservar la referencia. Pedí revisión antes de continuar.",
    "en": "Not sent: we could not retain the reference. Request a review before continuing."
  },
  "p0521": {
    "es": "Publicando versión…",
    "en": "Publishing version…"
  },
  "p0522": {
    "es": "Versión publicada. Consultá el resultado antes de preparar otra.",
    "en": "Version published. Check the result before preparing another."
  },
  "p0523": {
    "es": "No pudimos comprobar el resultado. La versión puede estar publicada; consultá sin reenviarla.",
    "en": "We could not verify the result. The version may be published; check without resubmitting it."
  },
  "p0524": {
    "es": "Versión recuperada: su contenido coincide con la referencia conservada. Revisalo antes de preparar otra.",
    "en": "Version retrieved: its content matches the retained reference. Review it before preparing another."
  },
  "p0525": {
    "es": "No pudimos comprobar esa versión y su contenido. Conservamos el bloqueo; consultá nuevamente o pedí revisión autorizada.",
    "en": "We could not verify that version and its content. The lock is retained; check again or request an authorized review."
  },
  "p0526": {
    "es": "Formulario nuevo preparado. La versión anterior sigue publicada e inmutable.",
    "en": "New form ready. The previous version remains published and immutable."
  },
  "p0527": {
    "es": "No pudimos actualizar la referencia. Conservamos el bloqueo; pedí revisión.",
    "en": "We could not update the reference. The lock is retained; request a review."
  },
  "p0528": {
    "es": "Publicar versión de checklist",
    "en": "Publish checklist version"
  },
  "p0529": {
    "es": "ID del checklist",
    "en": "Checklist ID"
  },
  "p0530": {
    "es": "Versión",
    "en": "Version"
  },
  "p0531": {
    "es": "Título",
    "en": "Title"
  },
  "p0532": {
    "es": "Ítem",
    "en": "Item"
  },
  "p0533": {
    "es": "ID",
    "en": "ID"
  },
  "p0534": {
    "es": "Pregunta o control",
    "en": "Question or check"
  },
  "p0535": {
    "es": "Respuesta",
    "en": "Answer"
  },
  "p0536": {
    "es": "Texto",
    "en": "Text"
  },
  "p0537": {
    "es": "Serie",
    "en": "Serial number"
  },
  "p0538": {
    "es": "Obligatorio",
    "en": "Required"
  },
  "p0539": {
    "es": "Quitar ítem",
    "en": "Remove item"
  },
  "p0540": {
    "es": "Agregar ítem",
    "en": "Add item"
  },
  "p0541": {
    "es": "Publicar versión",
    "en": "Publish version"
  },
  "p0542": {
    "es": "Resultado de publicación de checklist",
    "en": "Checklist publication result"
  },
  "p0543": {
    "es": "· versión",
    "en": "· version"
  },
  "p0544": {
    "es": "Opcional",
    "en": "Optional"
  },
  "p0545": {
    "es": "Consultar versión publicada",
    "en": "Check published version"
  },
  "p0546": {
    "es": "Preparar otra versión",
    "en": "Prepare another version"
  },
  "p0547": {
    "es": "Hay una presentación pendiente de comprobar. Consultá sin volver a enviarla.",
    "en": "A submission is awaiting verification. Check without sending it again."
  },
  "p0548": {
    "es": "Revisá entrega, versiones y respuestas únicas con formato item-id=respuesta.",
    "en": "Review the handover, versions and unique answers in item-id=answer format."
  },
  "p0549": {
    "es": "Completando y presentando…",
    "en": "Completing and submitting…"
  },
  "p0550": {
    "es": "Respuesta recibida. Consultá la presentación y sus respuestas antes de preparar otra.",
    "en": "Response received. Check the submission and its answers before preparing another."
  },
  "p0551": {
    "es": "No pudimos comprobar el resultado. La entrega puede estar presentada; consultá sin reenviarla.",
    "en": "We could not verify the result. The handover may be submitted; check without resubmitting it."
  },
  "p0552": {
    "es": "Presentación recuperada: checklist y respuestas coinciden. Revisá el estado actual y el actor antes de preparar otra.",
    "en": "Submission retrieved: the checklist and answers match. Review the current status and actor before preparing another."
  },
  "p0553": {
    "es": "No pudimos comprobar esa presentación y sus respuestas. Conservamos el bloqueo; consultá nuevamente o pedí revisión autorizada.",
    "en": "We could not verify that submission and its answers. The lock is retained; check again or request an authorized review."
  },
  "p0554": {
    "es": "Formulario nuevo preparado. La entrega anterior conserva su estado y sus respuestas.",
    "en": "New form ready. The previous handover retains its status and answers."
  },
  "p0555": {
    "es": "Presentación de entrega",
    "en": "Handover submission"
  },
  "p0556": {
    "es": "Completar y presentar una entrega preparada",
    "en": "Complete and submit a prepared handover"
  },
  "p0557": {
    "es": "Versión actual de entrega",
    "en": "Current handover version"
  },
  "p0558": {
    "es": "Versión del checklist",
    "en": "Checklist version"
  },
  "p0559": {
    "es": "Respuestas, una por línea",
    "en": "Answers, one per line"
  },
  "p0560": {
    "es": "serial-observed=SERIAL-123\nasset-condition=confirmed",
    "en": "serial-observed=SERIAL-123\nasset-condition=confirmed"
  },
  "p0561": {
    "es": "Completar y presentar",
    "en": "Complete and submit"
  },
  "p0562": {
    "es": "Resultado de presentación",
    "en": "Submission result"
  },
  "p0563": {
    "es": "· estado actual:",
    "en": "· current status:"
  },
  "p0564": {
    "es": "Checklist",
    "en": "Checklist"
  },
  "p0565": {
    "es": "· completado por",
    "en": "· completed by"
  },
  "p0566": {
    "es": "Consultar presentación",
    "en": "Check submission"
  },
  "p0567": {
    "es": "Preparar otra presentación",
    "en": "Prepare another submission"
  },
  "p0568": {
    "es": "Inventario",
    "en": "Inventory"
  },
  "p0569": {
    "es": "Reembolso",
    "en": "Refund"
  },
  "p0570": {
    "es": "Cambio",
    "en": "Exchange"
  },
  "p0571": {
    "es": "Contabilidad",
    "en": "Accounting"
  },
  "p0572": {
    "es": "Fiscal",
    "en": "Tax"
  },
  "p0573": {
    "es": "Hay una operación pendiente de comprobar. Consultá el caso sin volver a enviarla.",
    "en": "An operation is awaiting verification. Check the case without sending it again."
  },
  "p0574": {
    "es": "Revisá el caso, la serie, la opción elegida y las observaciones.",
    "en": "Review the case, serial number, selected option and notes."
  },
  "p0575": {
    "es": "Registrando recepción…",
    "en": "Recording receipt…"
  },
  "p0576": {
    "es": "Registrando decisión…",
    "en": "Recording decision…"
  },
  "p0577": {
    "es": "Respuesta recibida. Consultá el caso y su evidencia antes de continuar operando.",
    "en": "Response received. Check the case and its evidence before continuing operations."
  },
  "p0578": {
    "es": "No pudimos comprobar el resultado. La operación puede estar registrada; consultá sin reenviarla.",
    "en": "We could not verify the result. The operation may be recorded; check without resubmitting it."
  },
  "p0579": {
    "es": "Caso recuperado: la evidencia coincide con la operación conservada. Revisá el actor y las solicitudes antes de continuar.",
    "en": "Case retrieved: the evidence matches the retained operation. Review the actor and requests before continuing."
  },
  "p0580": {
    "es": "No pudimos comprobar el caso y su evidencia. Conservamos el bloqueo; consultá nuevamente o pedí revisión autorizada.",
    "en": "We could not verify the case and its evidence. The lock is retained; check again or request an authorized review."
  },
  "p0581": {
    "es": "Caso actualizado. Elegí la siguiente acción; la operación anterior conserva su evidencia.",
    "en": "Case updated. Choose the next action; the previous operation retains its evidence."
  },
  "p0582": {
    "es": "Operaciones de devoluciones",
    "en": "Return operations"
  },
  "p0583": {
    "es": "Recepción y disposición de devoluciones",
    "en": "Return receipt and disposition"
  },
  "p0584": {
    "es": "Las decisiones registran solicitudes de trabajo. Su ejecución se verifica en los procesos de inventario, reembolso o cambio, contabilidad y fiscalidad correspondientes.",
    "en": "Decisions record work requests. Their execution is verified in the corresponding inventory, refund or exchange, accounting and tax processes."
  },
  "p0585": {
    "es": "Resultado de devolución",
    "en": "Return result"
  },
  "p0586": {
    "es": "Consultar operación de devolución",
    "en": "Check return operation"
  },
  "p0587": {
    "es": "Continuar operando",
    "en": "Continue operations"
  },
  "p0588": {
    "es": "Actualizar casos",
    "en": "Refresh cases"
  },
  "p0589": {
    "es": "Evidencia recuperada de devolución",
    "en": "Retrieved return evidence"
  },
  "p0590": {
    "es": "Autorización",
    "en": "Authorization"
  },
  "p0591": {
    "es": "Recibo",
    "en": "Receipt"
  },
  "p0592": {
    "es": "· recibido por",
    "en": "· received by"
  },
  "p0593": {
    "es": "Serie:",
    "en": "Serial number:"
  },
  "p0594": {
    "es": "Decisión",
    "en": "Decision"
  },
  "p0595": {
    "es": "· registrada por",
    "en": "· recorded by"
  },
  "p0596": {
    "es": ": solicitada ·",
    "en": ": requested ·"
  },
  "p0597": {
    "es": "Lista de devoluciones no disponible. La consulta de una operación conservada sigue disponible.",
    "en": "Return list unavailable. You can still check a retained operation."
  },
  "p0598": {
    "es": "Devolución",
    "en": "Return"
  },
  "p0599": {
    "es": "· unidad",
    "en": "· unit"
  },
  "p0600": {
    "es": "Serie recibida",
    "en": "Received serial number"
  },
  "p0601": {
    "es": "Condición",
    "en": "Condition"
  },
  "p0602": {
    "es": "Sellado",
    "en": "Sealed"
  },
  "p0603": {
    "es": "Abierto",
    "en": "Opened"
  },
  "p0604": {
    "es": "Dañado",
    "en": "Damaged"
  },
  "p0605": {
    "es": "Incompleto",
    "en": "Incomplete"
  },
  "p0606": {
    "es": "Inspección/observaciones",
    "en": "Inspection/notes"
  },
  "p0607": {
    "es": "Recibido: serie",
    "en": "Received: serial number"
  },
  "p0608": {
    "es": "Acción de inventario solicitada",
    "en": "Requested inventory action"
  },
  "p0609": {
    "es": "Cuarentena",
    "en": "Quarantine"
  },
  "p0610": {
    "es": "Reingreso",
    "en": "Restock"
  },
  "p0611": {
    "es": "Reparación",
    "en": "Repair"
  },
  "p0612": {
    "es": "Baja",
    "en": "Write-off"
  },
  "p0613": {
    "es": "Fundamento",
    "en": "Rationale"
  },
  "p0614": {
    "es": "Registrar decisión",
    "en": "Record decision"
  },
  "p0615": {
    "es": "· decisión",
    "en": "· decision"
  },
  "p0616": {
    "es": "No hay casos en la lista actual. Una referencia conservada se consulta de forma independiente.",
    "en": "No cases in the current list. A retained reference can be checked independently."
  },
  "p0617": {
    "es": "Preparación y liberación comercial",
    "en": "Preparation and commercial release"
  },
  "p0618": {
    "es": "Entregas: preparación y cierre comercial",
    "en": "Handovers: preparation and commercial closure"
  },
  "p0619": {
    "es": "La franquicia prepara la entrega y completa la checklist. El cliente registra la recepción desde su portal. Después se puede registrar el cierre comercial permitido por el perfil activo.",
    "en": "The franchise prepares the handover and completes the checklist. The customer records receipt in their portal. The commercial closure allowed by the active profile can then be recorded."
  },
  "p0620": {
    "es": "La lista de pedidos no está disponible o es parcial. Actualizá la consulta antes de operar.",
    "en": "The order list is unavailable or incomplete. Refresh the query before operating."
  },
  "p0621": {
    "es": "No hay pedidos para consultar.",
    "en": "No orders to check."
  },
  "p0622": {
    "es": "Estado consultado. Cada operación vuelve a validar sus condiciones en el servidor.",
    "en": "Status retrieved. Each operation revalidates its conditions on the server."
  },
  "p0623": {
    "es": "Todavía no pudimos verificar una preparación habilitada para este pedido. Actualizá después de confirmar pago, reserva y perfil. Las referencias guardadas siguen disponibles para consultar.",
    "en": "We have not yet verified an eligible preparation for this order. Refresh after confirming payment, reservation and profile. Saved references remain available to check."
  },
  "p0624": {
    "es": "No pudimos leer las referencias guardadas. Consultá el estado y pedí revisión antes de operar.",
    "en": "We could not read the saved references. Check the status and request a review before operating."
  },
  "p0625": {
    "es": "No se envió la operación porque no pudimos guardar su referencia. Conservá los datos y pedí revisión.",
    "en": "The operation was not sent because we could not save its reference. Keep the data and request a review."
  },
  "p0626": {
    "es": "Entrega preparada. Completá la checklist correspondiente y consultá el estado.",
    "en": "Handover prepared. Complete the corresponding checklist and check the status."
  },
  "p0627": {
    "es": "Recibo comercial registrado. Consultá su vigencia antes de continuar con una operación posterior.",
    "en": "Commercial receipt recorded. Check its validity before continuing with a subsequent operation."
  },
  "p0628": {
    "es": "El resultado quedó sin comprobar. La operación pudo guardarse: usá Consultar resultado, que recupera la misma referencia sin reenviar.",
    "en": "The result is unverified. The operation may have been saved: use Check result to retrieve the same reference without resubmitting."
  },
  "p0629": {
    "es": "Preparación recuperada: {handover_id}. Estado {handover_state}.",
    "en": "Preparation retrieved: {handover_id}. Status {handover_state}."
  },
  "p0630": {
    "es": "Recibo histórico recuperado. Esto no confirma su vigencia actual.",
    "en": "Historical receipt retrieved. This does not confirm its current validity."
  },
  "p0631": {
    "es": "No pudimos recuperar un resultado confirmado. Conservamos la referencia y el bloqueo de reenvío; consultá nuevamente o pedí revisión.",
    "en": "We could not retrieve a confirmed result. The reference and resubmission lock are retained; check again or request a review."
  },
  "p0632": {
    "es": "No pudimos verificar la vigencia. El recibo histórico permanece disponible, pero no confirma autorización actual.",
    "en": "We could not verify validity. The historical receipt remains available but does not confirm current authorization."
  },
  "p0633": {
    "es": "Consultar estado de entrega",
    "en": "Check handover status"
  },
  "p0634": {
    "es": "No hay una preparación habilitada verificada en esta consulta.",
    "en": "No eligible preparation has been verified in this query."
  },
  "p0635": {
    "es": "Checklist presentada. Falta que el cliente registre la recepción desde su portal de entregas.",
    "en": "Checklist submitted. The customer still needs to record receipt in their handover portal."
  },
  "p0636": {
    "es": "Preparar entrega",
    "en": "Prepare handover"
  },
  "p0637": {
    "es": "Consultar resultado de preparación",
    "en": "Check preparation result"
  },
  "p0638": {
    "es": "Registrar cierre comercial",
    "en": "Record commercial closure"
  },
  "p0639": {
    "es": "El perfil activo permite consultar elegibilidad; no habilita un recibo comercial.",
    "en": "The active profile allows eligibility checks; it does not enable a commercial receipt."
  },
  "p0640": {
    "es": "Consultar resultado comercial",
    "en": "Check commercial result"
  },
  "p0641": {
    "es": "Consultar vigencia del recibo",
    "en": "Check receipt validity"
  },
  "p0642": {
    "es": "Recibo comercial",
    "en": "Commercial receipt"
  },
  "p0643": {
    "es": "· registrado",
    "en": "· recorded"
  },
  "p0644": {
    "es": "El registro no modifica stock ni acredita despacho físico.",
    "en": "This record does not change stock or certify physical dispatch."
  },
  "p0645": {
    "es": "Condiciones verificadas en la consulta de {evaluated_at}, con límite {valid_until}. Una operación posterior debe volver a verificarlas.",
    "en": "Conditions verified in the query at {evaluated_at}, with a limit of {valid_until}. A subsequent operation must reverify them."
  },
  "p0646": {
    "es": "El recibo no tiene vigencia comprobada en la consulta de {evaluated_at}. Su historial se conserva.",
    "en": "The receipt has no verified validity in the query at {evaluated_at}. Its history is retained."
  },
  "p0647": {
    "es": "Vigencia pendiente de consulta.",
    "en": "Validity awaiting query."
  },
  "p0648": {
    "es": "Borrador",
    "en": "Draft"
  },
  "p0649": {
    "es": "Publicado",
    "en": "Published"
  },
  "p0650": {
    "es": "Archivado",
    "en": "Archived"
  },
  "p0651": {
    "es": "No pudimos leer la referencia guardada. Conservá esta pestaña y pedí revisión autorizada.",
    "en": "We could not read the saved reference. Keep this tab open and request an authorized review."
  },
  "p0652": {
    "es": "Versión histórica consultada.",
    "en": "Historical version retrieved."
  },
  "p0653": {
    "es": "Artículo actual consultado.",
    "en": "Current article retrieved."
  },
  "p0654": {
    "es": "Artículo no disponible. Revisá organización, permisos y publicación.",
    "en": "Article unavailable. Review the organization, permissions and publication."
  },
  "p0655": {
    "es": "Resultados consultados.",
    "en": "Results retrieved."
  },
  "p0656": {
    "es": "No pudimos consultar los artículos de esta organización.",
    "en": "We could not retrieve articles for this organization."
  },
  "p0657": {
    "es": "Historial consultado.",
    "en": "History retrieved."
  },
  "p0658": {
    "es": "Historial no disponible para este rol.",
    "en": "History unavailable for this role."
  },
  "p0659": {
    "es": "Operación registrada. Consultá el artículo actual antes de otro cambio.",
    "en": "Operation recorded. Check the current article before another change."
  },
  "p0660": {
    "es": "Resultado sin confirmar. Consultá el resultado pendiente antes de otra escritura.",
    "en": "Unconfirmed result. Check the pending result before another write."
  },
  "p0661": {
    "es": "No se envió una escritura nueva, o fue rechazada antes del servicio. Revisá campos y permisos.",
    "en": "No new write was sent, or it was rejected before reaching the service. Review fields and permissions."
  },
  "p0662": {
    "es": "Resultado recuperado sin repetir la escritura. Consultá el artículo actual para continuar.",
    "en": "Result retrieved without repeating the write. Check the current article to continue."
  },
  "p0663": {
    "es": "Resultado aún incierto. Conservá la referencia y pedí revisión autorizada.",
    "en": "Result still uncertain. Retain the reference and request an authorized review."
  },
  "p0664": {
    "es": "Consultar artículos",
    "en": "Browse articles"
  },
  "p0665": {
    "es": "Organización del contenido",
    "en": "Content organization"
  },
  "p0666": {
    "es": "Idioma del contenido",
    "en": "Content language"
  },
  "p0667": {
    "es": "Español",
    "en": "Español"
  },
  "p0668": {
    "es": "English",
    "en": "English"
  },
  "p0669": {
    "es": "Buscar en título y texto",
    "en": "Search title and text"
  },
  "p0670": {
    "es": "Buscar artículos",
    "en": "Search articles"
  },
  "p0671": {
    "es": "Siguiente página de artículos",
    "en": "Next article page"
  },
  "p0672": {
    "es": "Referencia del artículo",
    "en": "Article reference"
  },
  "p0673": {
    "es": "Consultar artículo actual",
    "en": "Check current article"
  },
  "p0674": {
    "es": ". Versión",
    "en": ". Version"
  },
  "p0675": {
    "es": ". Idioma",
    "en": ". Language"
  },
  "p0676": {
    "es": "Versión histórica. Consultá el artículo actual antes de modificarlo.",
    "en": "Historical version. Check the current article before changing it."
  },
  "p0677": {
    "es": "Publicar artículo",
    "en": "Publish article"
  },
  "p0678": {
    "es": "Archivar artículo",
    "en": "Archive article"
  },
  "p0679": {
    "es": "Consultar historial",
    "en": "View history"
  },
  "p0680": {
    "es": "Ver versión",
    "en": "View version"
  },
  "p0681": {
    "es": "Versiones anteriores",
    "en": "Previous versions"
  },
  "p0682": {
    "es": "Editar borrador",
    "en": "Edit draft"
  },
  "p0683": {
    "es": "Nuevo artículo",
    "en": "New article"
  },
  "p0684": {
    "es": "Preparar artículo nuevo",
    "en": "Prepare new article"
  },
  "p0685": {
    "es": "Categoría",
    "en": "Category"
  },
  "p0686": {
    "es": "Título del artículo",
    "en": "Article title"
  },
  "p0687": {
    "es": "Texto del artículo",
    "en": "Article text"
  },
  "p0688": {
    "es": "Guardar revisión",
    "en": "Save revision"
  },
  "p0689": {
    "es": "Guardar borrador",
    "en": "Save draft"
  },
  "p0690": {
    "es": "Consultando ayuda…",
    "en": "Loading help…"
  },
  "p0691": {
    "es": "Iniciá sesión para consultar la ayuda.",
    "en": "Sign in to view help."
  },
  "p0692": {
    "es": "Tu sesión no tiene guías disponibles.",
    "en": "No guides are available for your session."
  },
  "p0693": {
    "es": "Esta guía o versión no está disponible para tu sesión.",
    "en": "This guide or version is unavailable for your session."
  },
  "p0694": {
    "es": "No pudimos consultar la ayuda. Podés volver a intentar.",
    "en": "We could not load help. You can try again."
  },
  "p0695": {
    "es": "Ayuda actualizada. Consulta de sólo lectura.",
    "en": "Help refreshed. Read-only query."
  },
  "p0696": {
    "es": "No encontramos guías para esa búsqueda.",
    "en": "No guides match that search."
  },
  "p0697": {
    "es": "Volvé a consultar para actualizar la ayuda.",
    "en": "Query again to refresh help."
  },
  "p0698": {
    "es": "Ayuda de los recorridos",
    "en": "Journey help"
  },
  "p0699": {
    "es": "Ayuda",
    "en": "Help"
  },
  "p0700": {
    "es": "Guías de uso en español, según los permisos de tu sesión. Leerlas no ejecuta operaciones ni acredita una capacitación.",
    "en": "Usage guides according to your session permissions. Reading them does not perform operations or certify training."
  },
  "p0701": {
    "es": "Buscar en las guías",
    "en": "Search guides"
  },
  "p0702": {
    "es": "Buscar",
    "en": "Search"
  },
  "p0703": {
    "es": "Actualizar ayuda",
    "en": "Refresh help"
  },
  "p0704": {
    "es": "Enlace a esta versión",
    "en": "Link to this version"
  },
  "p0705": {
    "es": "Activa",
    "en": "Active"
  },
  "p0706": {
    "es": "Suspendida",
    "en": "Suspended"
  },
  "p0707": {
    "es": "Cerrada",
    "en": "Closed"
  },
  "p0708": {
    "es": "Terminado",
    "en": "Terminated"
  },
  "p0709": {
    "es": "Vencido",
    "en": "Expired"
  },
  "p0710": {
    "es": "Activar",
    "en": "Activate"
  },
  "p0711": {
    "es": "Suspender",
    "en": "Suspend"
  },
  "p0712": {
    "es": "Cerrar",
    "en": "Close"
  },
  "p0713": {
    "es": "Terminar",
    "en": "Terminate"
  },
  "p0714": {
    "es": "Marcar vencido",
    "en": "Mark expired"
  },
  "p0715": {
    "es": "Estado actual consultado.",
    "en": "Current status retrieved."
  },
  "p0716": {
    "es": "No pudimos consultar esta referencia. Revisá organización y permisos.",
    "en": "We could not retrieve this reference. Review the organization and permissions."
  },
  "p0717": {
    "es": "Operación registrada. Consultá el estado actual antes de otra transición.",
    "en": "Operation recorded. Check the current status before another transition."
  },
  "p0718": {
    "es": "No se envió una escritura nueva, o fue rechazada antes del servicio. Revisá los campos y permisos.",
    "en": "No new write was sent, or it was rejected before reaching the service. Review the fields and permissions."
  },
  "p0719": {
    "es": "Resultado recuperado sin repetir la escritura. Consultá el estado actual para continuar.",
    "en": "Result retrieved without repeating the write. Check the current status to continue."
  },
  "p0720": {
    "es": "Registrar organización",
    "en": "Register organization"
  },
  "p0721": {
    "es": "Tipo de organización",
    "en": "Organization type"
  },
  "p0722": {
    "es": "Empresa",
    "en": "Company"
  },
  "p0723": {
    "es": "Franquiciante",
    "en": "Franchisor"
  },
  "p0724": {
    "es": "Franquiciado",
    "en": "Franchisee"
  },
  "p0725": {
    "es": "Depósito",
    "en": "Warehouse"
  },
  "p0726": {
    "es": "Sucursal",
    "en": "Branch"
  },
  "p0727": {
    "es": "Centro de servicio",
    "en": "Service center"
  },
  "p0728": {
    "es": "Organización padre",
    "en": "Parent organization"
  },
  "p0729": {
    "es": "Código de organización",
    "en": "Organization code"
  },
  "p0730": {
    "es": "Usá minúsculas, números y guiones.",
    "en": "Use lowercase letters, numbers and hyphens."
  },
  "p0731": {
    "es": "Nombre visible",
    "en": "Display name"
  },
  "p0732": {
    "es": "Registrar acuerdo",
    "en": "Register agreement"
  },
  "p0733": {
    "es": "Organización del franquiciado",
    "en": "Franchisee organization"
  },
  "p0734": {
    "es": "Territorio",
    "en": "Territory"
  },
  "p0735": {
    "es": "Versión de los términos",
    "en": "Terms version"
  },
  "p0736": {
    "es": "Fecha de inicio",
    "en": "Start date"
  },
  "p0737": {
    "es": "Fecha de fin opcional",
    "en": "Optional end date"
  },
  "p0738": {
    "es": "Registrar acuerdo en borrador",
    "en": "Register draft agreement"
  },
  "p0739": {
    "es": "Consultar estructura o acuerdo",
    "en": "Check structure or agreement"
  },
  "p0740": {
    "es": "Tipo de consulta",
    "en": "Query type"
  },
  "p0741": {
    "es": "Acuerdo",
    "en": "Agreement"
  },
  "p0742": {
    "es": "Referencia a consultar",
    "en": "Reference to check"
  },
  "p0743": {
    "es": "Organización de la consulta",
    "en": "Query organization"
  },
  "p0744": {
    "es": "Organización registrada",
    "en": "Registered organization"
  },
  "p0745": {
    "es": "Acuerdo registrado",
    "en": "Registered agreement"
  },
  "p0746": {
    "es": ". Términos",
    "en": ". Terms"
  },
  "p0747": {
    "es": ". Vigencia:",
    "en": ". Validity:"
  },
  "p0748": {
    "es": "a",
    "en": "to"
  },
  "p0749": {
    "es": "sin fin especificado",
    "en": "no end specified"
  },
  "p0750": {
    "es": "Resultado histórico de la operación. Consultá el estado actual antes de otra transición.",
    "en": "Historical operation result. Check the current status before another transition."
  },
  "p0751": {
    "es": "organización",
    "en": "organization"
  },
  "p0752": {
    "es": "Accesos de tu sesión",
    "en": "Your session access"
  },
  "p0753": {
    "es": "Accesos disponibles",
    "en": "Available access"
  },
  "p0754": {
    "es": "Las opciones dependen de tus accesos y de las funciones habilitadas.",
    "en": "Options depend on your access and enabled features."
  },
  "p0755": {
    "es": "Secciones disponibles",
    "en": "Available sections"
  },
  "p0756": {
    "es": "No hay accesos disponibles para esta vista.",
    "en": "No access is available for this view."
  },
  "p0757": {
    "es": "No pudimos conservar o leer la referencia de recuperación. Revisá el almacenamiento del navegador antes de operar.",
    "en": "We could not retain or read the recovery reference. Check browser storage before operating."
  },
  "p0758": {
    "es": "No pudimos consultar ese pedido. Revisá la organización y la referencia.",
    "en": "We could not retrieve that order. Review the organization and reference."
  },
  "p0759": {
    "es": "La aprobación no está disponible. Volvé a consultar su estado.",
    "en": "Approval is unavailable. Check its status again."
  },
  "p0760": {
    "es": "Solicitud guardada. Otra persona debe revisar el contenido exacto.",
    "en": "Request saved. Another person must review the exact content."
  },
  "p0761": {
    "es": "Cobertura total registrada. Ya podés continuar con la preparación de la entrega.",
    "en": "Full coverage recorded. You can now continue preparing the handover."
  },
  "p0762": {
    "es": "Decisión guardada. Consultá el pedido para ver el saldo actualizado.",
    "en": "Decision saved. Check the order to see the updated balance."
  },
  "p0763": {
    "es": "Resultado sin confirmar. Conservamos la misma referencia: consultá el resultado o reintentá esa misma solicitud.",
    "en": "Unconfirmed result. The same reference is retained: check the result or retry that same request."
  },
  "p0764": {
    "es": "La propuesta sigue pendiente. Podés reintentar la misma decisión.",
    "en": "The proposal is still pending. You can retry the same decision."
  },
  "p0765": {
    "es": "Resultado recuperado desde su registro.",
    "en": "Result retrieved from its record."
  },
  "p0766": {
    "es": "Todavía no hay un resultado confirmado. Conservá la referencia y consultá de nuevo o reintentá la misma solicitud.",
    "en": "There is no confirmed result yet. Retain the reference and check again or retry the same request."
  },
  "p0767": {
    "es": "Solicitud por confirmar",
    "en": "Request awaiting confirmation"
  },
  "p0768": {
    "es": "La referencia permanece guardada en este navegador.",
    "en": "The reference remains saved in this browser."
  },
  "p0769": {
    "es": "Reintentar la misma solicitud",
    "en": "Retry the same request"
  },
  "p0770": {
    "es": "Consultar pedido",
    "en": "Check order"
  },
  "p0771": {
    "es": "Referencia del pedido",
    "en": "Order reference"
  },
  "p0772": {
    "es": "Total del pedido",
    "en": "Order total"
  },
  "p0773": {
    "es": "Gift cards aplicadas",
    "en": "Applied gift cards"
  },
  "p0774": {
    "es": "Descuento de fidelidad",
    "en": "Loyalty discount"
  },
  "p0775": {
    "es": "Importe pendiente de pago",
    "en": "Amount remaining to pay"
  },
  "p0776": {
    "es": "Solicitudes registradas para este pedido",
    "en": "Requests recorded for this order"
  },
  "p0777": {
    "es": "No hay solicitudes en esta página.",
    "en": "No requests on this page."
  },
  "p0778": {
    "es": "Revisar",
    "en": "Review"
  },
  "p0779": {
    "es": "Ver más solicitudes",
    "en": "View more requests"
  },
  "p0780": {
    "es": "Registrar cobertura total para entrega",
    "en": "Record full coverage for handover"
  },
  "p0781": {
    "es": "Solicitar una operación",
    "en": "Request an operation"
  },
  "p0782": {
    "es": "Consultá primero el pedido. Los importes se calculan con el programa habilitado y quedan pendientes de revisión.",
    "en": "Check the order first. Amounts are calculated with the enabled program and await review."
  },
  "p0783": {
    "es": "Programa",
    "en": "Program"
  },
  "p0784": {
    "es": "Gift card",
    "en": "Gift card"
  },
  "p0785": {
    "es": "Fidelidad",
    "en": "Loyalty"
  },
  "p0786": {
    "es": "Canjear saldo o puntos",
    "en": "Redeem balance or points"
  },
  "p0787": {
    "es": "Emitir gift card desde una venta confirmada",
    "en": "Issue a gift card from a confirmed sale"
  },
  "p0788": {
    "es": "Acreditar puntos desde una venta confirmada",
    "en": "Credit points from a confirmed sale"
  },
  "p0789": {
    "es": "Revertir una operación de un pedido cancelado",
    "en": "Reverse an operation on a cancelled order"
  },
  "p0790": {
    "es": "Cuenta de gift card o fidelidad",
    "en": "Gift card or loyalty account"
  },
  "p0791": {
    "es": "Referencia de la operación original",
    "en": "Original operation reference"
  },
  "p0792": {
    "es": "Guardar solicitud para revisión",
    "en": "Save request for review"
  },
  "p0793": {
    "es": "Revisión de operación",
    "en": "Operation review"
  },
  "p0794": {
    "es": "Pedido:",
    "en": "Order:"
  },
  "p0795": {
    "es": ". Solicitante:",
    "en": ". Requester:"
  },
  "p0796": {
    "es": ". Estado:",
    "en": ". Status:"
  },
  "p0797": {
    "es": "Pendiente",
    "en": "Pending"
  },
  "p0798": {
    "es": "Aprobada",
    "en": "Approved"
  },
  "p0799": {
    "es": "Vigencia registrada:",
    "en": "Recorded validity:"
  },
  "p0800": {
    "es": "Revisá el pedido, la cuenta, los puntos y el importe aplicable en el contenido exacto. Aprobar confirma este contenido.",
    "en": "Review the order, account, points and applicable amount in the exact content. Approval confirms this content."
  },
  "p0801": {
    "es": "Emisión de gift card",
    "en": "Gift card issuance"
  },
  "p0802": {
    "es": "Acreditación de puntos",
    "en": "Points credit"
  },
  "p0803": {
    "es": "Canje",
    "en": "Redemption"
  },
  "p0804": {
    "es": "Reversión",
    "en": "Reversal"
  },
  "p0805": {
    "es": "Referencia",
    "en": "Reference"
  },
  "p0806": {
    "es": "Total del pedido revisado",
    "en": "Reviewed order total"
  },
  "p0807": {
    "es": "Operación que se revierte",
    "en": "Operation to reverse"
  },
  "p0808": {
    "es": "Cuenta",
    "en": "Account"
  },
  "p0809": {
    "es": "Cambio en saldo o puntos:",
    "en": "Balance or points change:"
  },
  "p0810": {
    "es": "Cambio en importe aplicado al pedido:",
    "en": "Change in amount applied to the order:"
  },
  "p0811": {
    "es": "Ver contenido completo de la solicitud",
    "en": "View full request content"
  },
  "p0812": {
    "es": "Contenido exacto a revisar",
    "en": "Exact content to review"
  },
  "p0813": {
    "es": "Aprobar contenido revisado",
    "en": "Approve reviewed content"
  },
  "p0814": {
    "es": "Rechazar solicitud",
    "en": "Reject request"
  },
  "p0815": {
    "es": "La solicitud necesita la revisión de otra persona autorizada.",
    "en": "The request requires review by another authorized person."
  },
  "p0816": {
    "es": "Actualizar estado de revisión",
    "en": "Refresh review status"
  },
  "p0817": {
    "es": "Ver comprobante y referencias de las cuentas",
    "en": "View receipt and account references"
  },
  "p0818": {
    "es": "Cobertura total registrada",
    "en": "Full coverage recorded"
  },
  "p0819": {
    "es": ". Referencia:",
    "en": ". Reference:"
  },
  "p0820": {
    "es": "Continuar con la entrega",
    "en": "Continue with handover"
  },
  "p0821": {
    "es": "Ver comprobante exacto",
    "en": "View exact receipt"
  },
  "p0822": {
    "es": "Enviada a fábrica",
    "en": "Sent to factory"
  },
  "p0823": {
    "es": "Aceptada por fábrica",
    "en": "Accepted by factory"
  },
  "p0824": {
    "es": "En producción",
    "en": "In production"
  },
  "p0825": {
    "es": "Despachada",
    "en": "Dispatched"
  },
  "p0826": {
    "es": "Recibida",
    "en": "Received"
  },
  "p0827": {
    "es": "Planificada",
    "en": "Planned"
  },
  "p0828": {
    "es": "En montaje",
    "en": "In assembly"
  },
  "p0829": {
    "es": "En revisión",
    "en": "Under review"
  },
  "p0830": {
    "es": "Liberada por fábrica",
    "en": "Released by factory"
  },
  "p0831": {
    "es": "En tránsito",
    "en": "In transit"
  },
  "p0832": {
    "es": "En cuarentena",
    "en": "Quarantined"
  },
  "p0833": {
    "es": "Disponible",
    "en": "Available"
  },
  "p0834": {
    "es": "Retirada",
    "en": "Withdrawn"
  },
  "p0835": {
    "es": "No pudimos leer la referencia guardada. Conservá esta pestaña y pedí revisión antes de operar.",
    "en": "We could not read the saved reference. Keep this tab open and request a review before operating."
  },
  "p0836": {
    "es": "Orden y versión consultadas.",
    "en": "Order and version retrieved."
  },
  "p0837": {
    "es": "No pudimos consultar la orden. Revisá su referencia y organización.",
    "en": "We could not retrieve the order. Review its reference and organization."
  },
  "p0838": {
    "es": "Operación registrada.",
    "en": "Operation recorded."
  },
  "p0839": {
    "es": "Operación registrada. Volvé a consultar la orden para continuar.",
    "en": "Operation recorded. Check the order again to continue."
  },
  "p0840": {
    "es": "No se envió una escritura nueva, o fue rechazada antes del servicio. Revisá campos, evidencia y permisos.",
    "en": "No new write was sent, or it was rejected before reaching the service. Review fields, evidence and permissions."
  },
  "p0841": {
    "es": "Resultado recuperado sin repetir la escritura.",
    "en": "Result retrieved without repeating the write."
  },
  "p0842": {
    "es": "Resultado recuperado. Volvé a consultar la orden.",
    "en": "Result retrieved. Check the order again."
  },
  "p0843": {
    "es": "El resultado sigue sin confirmar. Conservá la referencia y pedí revisión autorizada si persiste.",
    "en": "The result remains unconfirmed. Retain the reference and request an authorized review if this persists."
  },
  "p0844": {
    "es": "Seleccioná evidencia de hasta 1 MiB.",
    "en": "Select evidence up to 1 MiB."
  },
  "p0845": {
    "es": "Orden:",
    "en": "Order:"
  },
  "p0846": {
    "es": "Comando:",
    "en": "Command:"
  },
  "p0847": {
    "es": "Evidencia de la operación",
    "en": "Operation evidence"
  },
  "p0848": {
    "es": "Archivo de evidencia",
    "en": "Evidence file"
  },
  "p0849": {
    "es": "Se registra la huella del archivo, no su contenido. Seleccioná la evidencia correspondiente a cada decisión.",
    "en": "The file fingerprint is recorded, not its content. Select the evidence corresponding to each decision."
  },
  "p0850": {
    "es": "Evidencia preparada.",
    "en": "Evidence ready."
  },
  "p0851": {
    "es": "Motivo de revisión",
    "en": "Review reason"
  },
  "p0852": {
    "es": "Crear orden y plan",
    "en": "Create order and plan"
  },
  "p0853": {
    "es": "Proveedor",
    "en": "Supplier"
  },
  "p0854": {
    "es": "Organización de fábrica",
    "en": "Factory organization"
  },
  "p0855": {
    "es": "Moneda",
    "en": "Currency"
  },
  "p0856": {
    "es": "Importe total en unidades menores",
    "en": "Total amount in minor units"
  },
  "p0857": {
    "es": "Referencia de la demanda",
    "en": "Demand reference"
  },
  "p0858": {
    "es": "Línea",
    "en": "Line"
  },
  "p0859": {
    "es": "Cantidad",
    "en": "Quantity"
  },
  "p0860": {
    "es": "Quitar línea",
    "en": "Remove line"
  },
  "p0861": {
    "es": "Agregar línea",
    "en": "Add line"
  },
  "p0862": {
    "es": "Crear orden con cantidades",
    "en": "Create order with quantities"
  },
  "p0863": {
    "es": "Consultar orden",
    "en": "Check order"
  },
  "p0864": {
    "es": "Referencia de la orden",
    "en": "Order reference"
  },
  "p0865": {
    "es": "Estado de la orden:",
    "en": "Order status:"
  },
  "p0866": {
    "es": "Proveedor:",
    "en": "Supplier:"
  },
  "p0867": {
    "es": ". Fábrica:",
    "en": ". Factory:"
  },
  "p0868": {
    "es": ". Destino:",
    "en": ". Destination:"
  },
  "p0869": {
    "es": ". Demanda:",
    "en": ". Demand:"
  },
  "p0870": {
    "es": "Abrir esta orden en la vista de fábrica",
    "en": "Open this order in the factory view"
  },
  "p0871": {
    "es": "unidades · línea",
    "en": "units · line"
  },
  "p0872": {
    "es": "Enviar a fábrica",
    "en": "Send to factory"
  },
  "p0873": {
    "es": "Cancelar orden",
    "en": "Cancel order"
  },
  "p0874": {
    "es": "Confirmar orden en fábrica",
    "en": "Confirm order at factory"
  },
  "p0875": {
    "es": "Iniciar producción",
    "en": "Start production"
  },
  "p0876": {
    "es": "Registrar una serie",
    "en": "Register a serial number"
  },
  "p0877": {
    "es": "Línea del plan",
    "en": "Plan line"
  },
  "p0878": {
    "es": "Número de serie",
    "en": "Serial number"
  },
  "p0879": {
    "es": "VIN opcional",
    "en": "Optional VIN"
  },
  "p0880": {
    "es": "Serie de batería opcional",
    "en": "Optional battery serial number"
  },
  "p0881": {
    "es": "Registrar serie",
    "en": "Register serial number"
  },
  "p0882": {
    "es": "Producción:",
    "en": "Production:"
  },
  "p0883": {
    "es": "Revisión de recepción:",
    "en": "Receipt review:"
  },
  "p0884": {
    "es": "Envío:",
    "en": "Shipment:"
  },
  "p0885": {
    "es": "Seleccionar",
    "en": "Select"
  },
  "p0886": {
    "es": "Registrar montaje",
    "en": "Record assembly"
  },
  "p0887": {
    "es": "Solicitar revisión de fábrica",
    "en": "Request factory review"
  },
  "p0888": {
    "es": "Debe decidir una persona diferente de quien solicitó la revisión.",
    "en": "A different person from the review requester must decide."
  },
  "p0889": {
    "es": "Liberar serie en fábrica",
    "en": "Release serial number at factory"
  },
  "p0890": {
    "es": "Rechazar serie en fábrica",
    "en": "Reject serial number at factory"
  },
  "p0891": {
    "es": "La liberación requiere una persona diferente de quien recibió o pidió la reinspección.",
    "en": "Release requires a different person from whoever received the stock or requested reinspection."
  },
  "p0892": {
    "es": "Liberar stock recibido",
    "en": "Release received stock"
  },
  "p0893": {
    "es": "Rechazar calidad recibida",
    "en": "Reject received quality"
  },
  "p0894": {
    "es": "Solicitar reinspección",
    "en": "Request reinspection"
  },
  "p0895": {
    "es": "Siguiente página de series",
    "en": "Next serial numbers page"
  },
  "p0896": {
    "es": "Actualizar desde primera página",
    "en": "Refresh from first page"
  },
  "p0897": {
    "es": "Despachar series seleccionadas",
    "en": "Dispatch selected serial numbers"
  },
  "p0898": {
    "es": "Recibir series seleccionadas",
    "en": "Receive selected serial numbers"
  },
  "p0899": {
    "es": "Referencia del envío",
    "en": "Shipment reference"
  },
  "p0900": {
    "es": "series seleccionadas de esta página.",
    "en": "serial numbers selected on this page."
  },
  "p0901": {
    "es": "Registrar despacho",
    "en": "Record dispatch"
  },
  "p0902": {
    "es": "Registrar recepción en cuarentena",
    "en": "Record receipt into quarantine"
  },
  "p0903": {
    "es": "Dueño",
    "en": "Owner"
  },
  "p0904": {
    "es": "Administrador",
    "en": "Administrator"
  },
  "p0905": {
    "es": "Evaluación favorable",
    "en": "Favorable assessment"
  },
  "p0906": {
    "es": "Requiere otra práctica",
    "en": "Further practice required"
  },
  "p0907": {
    "es": "Pendiente de revisión humana",
    "en": "Awaiting human review"
  },
  "p0908": {
    "es": "Hay un intento guardado. Consultá su estado antes de continuar.",
    "en": "A saved attempt exists. Check its status before continuing."
  },
  "p0909": {
    "es": "No pudimos leer la referencia local. Conservá la pestaña y pedí revisión autorizada.",
    "en": "We could not read the local reference. Keep the tab open and request an authorized review."
  },
  "p0910": {
    "es": "Estado consultado. La lectura no repite ninguna operación.",
    "en": "Status retrieved. Reading does not repeat any operation."
  },
  "p0911": {
    "es": "No pudimos confirmar el intento. Conservamos su referencia; volvé a consultar.",
    "en": "We could not confirm the attempt. Its reference is retained; check again."
  },
  "p0912": {
    "es": "Registro confirmado.",
    "en": "Record confirmed."
  },
  "p0913": {
    "es": "La respuesta no quedó confirmada. Consultá el estado guardado antes de volver a actuar.",
    "en": "The response was not confirmed. Check the saved status before acting again."
  },
  "p0914": {
    "es": "No pudimos guardar la referencia. El inicio no se envió.",
    "en": "We could not save the reference. The start request was not sent."
  },
  "p0915": {
    "es": "Podés preparar otra práctica. El resultado anterior conserva su evidencia.",
    "en": "You can prepare another practice. The previous result retains its evidence."
  },
  "p0916": {
    "es": "No pudimos actualizar la referencia. Conservamos el intento anterior.",
    "en": "We could not update the reference. The previous attempt is retained."
  },
  "p0917": {
    "es": "Completá todas las respuestas con textos de hasta 2048 bytes.",
    "en": "Complete all answers with text up to 2048 bytes each."
  },
  "p0918": {
    "es": "Mi capacitación",
    "en": "My training"
  },
  "p0919": {
    "es": "Guía de rol",
    "en": "Role guide"
  },
  "p0920": {
    "es": "Elegir una guía organiza el contenido disponible; no cambia los permisos de la sesión.",
    "en": "Choosing a guide organizes the available content; it does not change session permissions."
  },
  "p0921": {
    "es": "Iniciar práctica",
    "en": "Start practice"
  },
  "p0922": {
    "es": "Consultar intento guardado",
    "en": "Check saved attempt"
  },
  "p0923": {
    "es": "Reanudar este inicio",
    "en": "Resume this start"
  },
  "p0924": {
    "es": "Contenido",
    "en": "Content"
  },
  "p0925": {
    "es": ", revisión",
    "en": ", revision"
  },
  "p0926": {
    "es": "· guía de",
    "en": "· guide for"
  },
  "p0927": {
    "es": "Este intento conserva una versión anterior del contenido. Puede consultarse; prepará una práctica nueva para la versión activa.",
    "en": "This attempt retains an earlier content version. You can view it; prepare a new practice for the active version."
  },
  "p0928": {
    "es": "Lectura registrada",
    "en": "Reading recorded"
  },
  "p0929": {
    "es": "Registrar lectura",
    "en": "Record reading"
  },
  "p0930": {
    "es": "Respuestas para revisión humana",
    "en": "Answers for human review"
  },
  "p0931": {
    "es": "Presentar respuestas",
    "en": "Submit answers"
  },
  "p0932": {
    "es": "Preparar otra práctica",
    "en": "Prepare another practice"
  },
  "p0933": {
    "es": "Evaluaciones registradas",
    "en": "Recorded assessments"
  },
  "p0934": {
    "es": "Evaluaciones para revisar",
    "en": "Assessments to review"
  },
  "p0935": {
    "es": "Mis evaluaciones",
    "en": "My assessments"
  },
  "p0936": {
    "es": "Consultar evaluaciones",
    "en": "Browse assessments"
  },
  "p0937": {
    "es": "No hay evaluaciones registradas.",
    "en": "No assessments recorded."
  },
  "p0938": {
    "es": "Participante:",
    "en": "Participant:"
  },
  "p0939": {
    "es": "Versión presentada:",
    "en": "Submitted version:"
  },
  "p0940": {
    "es": "Fundamento de la evaluación",
    "en": "Assessment rationale"
  },
  "p0941": {
    "es": "Resultado",
    "en": "Outcome"
  },
  "p0942": {
    "es": "Registrar evaluación humana",
    "en": "Record human assessment"
  },
  "p0943": {
    "es": "Resultado de evaluación",
    "en": "Assessment outcome"
  },
  "p0944": {
    "es": "Revisó:",
    "en": "Reviewed by:"
  },
  "p0945": {
    "es": "Esta constancia registra la revisión del contenido y las respuestas; no modifica los permisos de operación.",
    "en": "This record documents review of the content and answers; it does not change operating permissions."
  },
  "p0946": {
    "es": "Diagnóstico",
    "en": "Diagnosis"
  },
  "p0947": {
    "es": "Calidad y aceptación",
    "en": "Quality and acceptance"
  },
  "p0948": {
    "es": "Cerrado",
    "en": "Closed"
  },
  "p0949": {
    "es": "Consulta actualizada.",
    "en": "Query refreshed."
  },
  "p0950": {
    "es": "No pudimos confirmar esta consulta. Revisá la referencia y sus permisos.",
    "en": "We could not confirm this query. Review the reference and its permissions."
  },
  "p0951": {
    "es": "Operación registrada. Volvé a consultar el caso para continuar.",
    "en": "Operation recorded. Check the case again to continue."
  },
  "p0952": {
    "es": "Resultado recuperado. Volvé a consultar el caso.",
    "en": "Result retrieved. Check the case again."
  },
  "p0953": {
    "es": "Seleccioná un archivo de evidencia de hasta 1 MiB.",
    "en": "Select an evidence file up to 1 MiB."
  },
  "p0954": {
    "es": "Evidencia y decisión",
    "en": "Evidence and decision"
  },
  "p0955": {
    "es": "Se registra su huella; el archivo y los motivos no se guardan en la referencia del navegador.",
    "en": "Its fingerprint is recorded; the file and reasons are not saved in the browser reference."
  },
  "p0956": {
    "es": "Términos de la cotización",
    "en": "Quote terms"
  },
  "p0957": {
    "es": "Referencia de la cotización",
    "en": "Quote reference"
  },
  "p0958": {
    "es": "Consultar cotización para términos",
    "en": "Check quote for terms"
  },
  "p0959": {
    "es": "Consultar política configurada",
    "en": "Check configured policy"
  },
  "p0960": {
    "es": "Cotización:",
    "en": "Quote:"
  },
  "p0961": {
    "es": "Vigente",
    "en": "Valid"
  },
  "p0962": {
    "es": "Sin vigencia",
    "en": "Not valid"
  },
  "p0963": {
    "es": "Términos",
    "en": "Terms"
  },
  "p0964": {
    "es": "Períodos configurados:",
    "en": "Configured periods:"
  },
  "p0965": {
    "es": "días para piezas y",
    "en": "days for parts and"
  },
  "p0966": {
    "es": "para trabajo.",
    "en": "for labor."
  },
  "p0967": {
    "es": "Vincular términos revisados",
    "en": "Bind reviewed terms"
  },
  "p0968": {
    "es": "Consultar términos ofrecidos",
    "en": "Check offered terms"
  },
  "p0969": {
    "es": "Términos ofrecidos",
    "en": "Offered terms"
  },
  "p0970": {
    "es": "Recepción de términos registrada.",
    "en": "Receipt of terms recorded."
  },
  "p0971": {
    "es": "Recepción de términos pendiente.",
    "en": "Receipt of terms pending."
  },
  "p0972": {
    "es": "Leí los términos mostrados de esta cotización",
    "en": "I have read the displayed terms for this quote"
  },
  "p0973": {
    "es": "Registrar recepción de términos",
    "en": "Record receipt of terms"
  },
  "p0974": {
    "es": "Garantía de la entrega",
    "en": "Handover warranty"
  },
  "p0975": {
    "es": "Referencia de la entrega",
    "en": "Handover reference"
  },
  "p0976": {
    "es": "Consultar garantía activa",
    "en": "Check active warranty"
  },
  "p0977": {
    "es": "Activar garantía de entrega aceptada",
    "en": "Activate warranty for accepted handover"
  },
  "p0978": {
    "es": "Garantía activa",
    "en": "Active warranty"
  },
  "p0979": {
    "es": ". Términos:",
    "en": ". Terms:"
  },
  "p0980": {
    "es": "Piezas:",
    "en": "Parts:"
  },
  "p0981": {
    "es": ". Trabajo:",
    "en": ". Labor:"
  },
  "p0982": {
    "es": "Abrir reclamo de la entrega",
    "en": "Open handover claim"
  },
  "p0983": {
    "es": "Referencia de la atención completada",
    "en": "Completed service reference"
  },
  "p0984": {
    "es": "Prioridad",
    "en": "Priority"
  },
  "p0985": {
    "es": "Media",
    "en": "Medium"
  },
  "p0986": {
    "es": "Alta",
    "en": "High"
  },
  "p0987": {
    "es": "Seguridad",
    "en": "Safety"
  },
  "p0988": {
    "es": "Descripción del problema",
    "en": "Problem description"
  },
  "p0989": {
    "es": "Abrir reclamo",
    "en": "Open claim"
  },
  "p0990": {
    "es": "Consultar caso",
    "en": "Check case"
  },
  "p0991": {
    "es": "Referencia del caso",
    "en": "Case reference"
  },
  "p0992": {
    "es": "Estado del caso:",
    "en": "Case status:"
  },
  "p0993": {
    "es": "Piezas cubiertas por fecha:",
    "en": "Parts covered by date:"
  },
  "p0994": {
    "es": "Sí",
    "en": "Yes"
  },
  "p0995": {
    "es": "No",
    "en": "No"
  },
  "p0996": {
    "es": ". Trabajo cubierto por fecha:",
    "en": ". Labor covered by date:"
  },
  "p0997": {
    "es": "Diagnóstico registrado",
    "en": "Recorded diagnosis"
  },
  "p0998": {
    "es": "La falla está excluida por los términos vendidos.",
    "en": "The fault is excluded by the sold terms."
  },
  "p0999": {
    "es": "La falla no figura entre las exclusiones configuradas.",
    "en": "The fault is not among the configured exclusions."
  },
  "p1000": {
    "es": "Registrar diagnóstico",
    "en": "Record diagnosis"
  },
  "p1001": {
    "es": "Código de falla",
    "en": "Fault code"
  },
  "p1002": {
    "es": "Descripción del diagnóstico",
    "en": "Diagnosis description"
  },
  "p1003": {
    "es": "Proponer reparación",
    "en": "Propose repair"
  },
  "p1004": {
    "es": "Trabajo propuesto",
    "en": "Proposed work"
  },
  "p1005": {
    "es": "Repuesto",
    "en": "Spare part"
  },
  "p1006": {
    "es": "Artículo",
    "en": "Item"
  },
  "p1007": {
    "es": "Ubicación de stock",
    "en": "Stock location"
  },
  "p1008": {
    "es": "Lote opcional",
    "en": "Optional lot"
  },
  "p1009": {
    "es": "Cantidad de repuesto",
    "en": "Spare part quantity"
  },
  "p1010": {
    "es": "Entrada específica opcional",
    "en": "Optional specific entry"
  },
  "p1011": {
    "es": "Quitar repuesto",
    "en": "Remove spare part"
  },
  "p1012": {
    "es": "Agregar repuesto",
    "en": "Add spare part"
  },
  "p1013": {
    "es": "Plan de reparación",
    "en": "Repair plan"
  },
  "p1014": {
    "es": "Propuesto por",
    "en": "Proposed by"
  },
  "p1015": {
    "es": ". Aprobación:",
    "en": ". Approval:"
  },
  "p1016": {
    "es": "Debe decidir una persona diferente de quien propuso el plan. La reserva pendiente vence:",
    "en": "A different person from the plan proposer must decide. The pending reservation expires:"
  },
  "p1017": {
    "es": "Aprobar reparación",
    "en": "Approve repair"
  },
  "p1018": {
    "es": "Rechazar reparación",
    "en": "Reject repair"
  },
  "p1019": {
    "es": "Registrar trabajo realizado",
    "en": "Record completed work"
  },
  "p1020": {
    "es": "Trabajo registrado por",
    "en": "Work recorded by"
  },
  "p1021": {
    "es": "Verificar calidad",
    "en": "Verify quality"
  },
  "p1022": {
    "es": "La revisión anterior falló. Adjuntá evidencia de la corrección antes de registrar otra.",
    "en": "The previous review failed. Attach evidence of the correction before recording another."
  },
  "p1023": {
    "es": "Evidencia de corrección",
    "en": "Correction evidence"
  },
  "p1024": {
    "es": "Aprobar calidad del trabajo",
    "en": "Approve work quality"
  },
  "p1025": {
    "es": "Registrar calidad fallida",
    "en": "Record failed quality"
  },
  "p1026": {
    "es": "Calidad del trabajo:",
    "en": "Work quality:"
  },
  "p1027": {
    "es": "Fallida",
    "en": "Failed"
  },
  "p1028": {
    "es": "Aceptar reparación recibida",
    "en": "Accept received repair"
  },
  "p1029": {
    "es": "Reparación aceptada por el cliente.",
    "en": "Repair accepted by the customer."
  },
  "p1030": {
    "es": "Conciliar servicio de garantía",
    "en": "Reconcile warranty service"
  },
  "p1031": {
    "es": "Servicio conciliado",
    "en": "Reconciled service"
  },
  "p1032": {
    "es": "Costo de inventario registrado:",
    "en": "Recorded inventory cost:"
  },
  "p1033": {
    "es": "Acuse interno del servicio; no representa un pago o documento fiscal.",
    "en": "Internal service acknowledgment; it does not represent a payment or tax document."
  },
  "p1034": {
    "es": "Cancelar reparación sin trabajo emitido",
    "en": "Cancel repair with no work issued"
  },
  "p1035": {
    "es": "Sin iniciar",
    "en": "Not started"
  },
  "p1036": {
    "es": "En proceso",
    "en": "In progress"
  },
  "p1037": {
    "es": "Aceptado por WhatsApp",
    "en": "Accepted by WhatsApp"
  },
  "p1038": {
    "es": "Resultado por confirmar",
    "en": "Result awaiting confirmation"
  },
  "p1039": {
    "es": "No aceptado",
    "en": "Not accepted"
  },
  "p1040": {
    "es": "Aún no hay confirmación del proveedor",
    "en": "No provider confirmation yet"
  },
  "p1041": {
    "es": "Enviado según WhatsApp",
    "en": "Sent according to WhatsApp"
  },
  "p1042": {
    "es": "Entregado según WhatsApp",
    "en": "Delivered according to WhatsApp"
  },
  "p1043": {
    "es": "Leído según WhatsApp",
    "en": "Read according to WhatsApp"
  },
  "p1044": {
    "es": "WhatsApp informó un fallo",
    "en": "WhatsApp reported a failure"
  },
  "p1045": {
    "es": "Hay informes que requieren revisión",
    "en": "Reports require review"
  },
  "p1046": {
    "es": "Texto y destinatario revisados en el portal",
    "en": "Text and recipient reviewed in the portal"
  },
  "p1047": {
    "es": "Respuesta rechazada desde el portal",
    "en": "Reply rejected from the portal"
  },
  "p1048": {
    "es": "La acción no está confirmada. Consultá el estado; un resultado incierto requiere reconciliación y no otro envío.",
    "en": "The action is unconfirmed. Check the status; an uncertain result requires reconciliation, not another send."
  },
  "p1049": {
    "es": "Consultar estado",
    "en": "Check status"
  },
  "p1050": {
    "es": "No hay respuestas para revisar.",
    "en": "No replies to review."
  },
  "p1051": {
    "es": "Respuesta para",
    "en": "Reply for"
  },
  "p1052": {
    "es": "Ventana de respuesta hasta",
    "en": "Reply window until"
  },
  "p1053": {
    "es": "Revisión:",
    "en": "Review:"
  },
  "p1054": {
    "es": ". Envío:",
    "en": ". Send:"
  },
  "p1055": {
    "es": "Requiere revisión",
    "en": "Review required"
  },
  "p1056": {
    "es": ". Estado informado:",
    "en": ". Reported status:"
  },
  "p1057": {
    "es": "Sin observación del proveedor",
    "en": "No provider observation"
  },
  "p1058": {
    "es": "Aprobar este texto",
    "en": "Approve this text"
  },
  "p1059": {
    "es": "Rechazar",
    "en": "Reject"
  },
  "p1060": {
    "es": "Enviar respuesta aprobada",
    "en": "Send approved reply"
  },
  "p1061": {
    "es": "El resultado del envío necesita reconciliación. La recuperación sólo consulta el comprobante guardado.",
    "en": "The send result needs reconciliation. Recovery only retrieves the saved receipt."
  },
  "p1062": {
    "es": "Verificar comprobante guardado",
    "en": "Verify saved receipt"
  },
  "p1063": {
    "es": "Aceptar una cotización",
    "en": "Accept a quote"
  },
  "p1064": {
    "es": "Revisá el importe, la moneda y la vigencia. Aceptar solicita crear un pedido con esos datos del servidor.",
    "en": "Review the amount, currency and validity. Acceptance requests creation of an order using those server values."
  },
  "p1065": {
    "es": "Aceptar una cotización no confirma el pago, el stock ni la entrega.",
    "en": "Accepting a quote does not confirm payment, stock or handover."
  },
  "p1066": {
    "es": "Si se corta la conexión, usá Actualizar estado o volvé a abrir Mis cotizaciones. No repitas la aceptación hasta comprobar el estado.",
    "en": "If the connection drops, use Refresh status or reopen My quotes. Do not repeat acceptance until you have checked the status."
  },
  "p1067": {
    "es": "Un pedido se confirma aquí cuando la lectura muestra su referencia. Si el problema continúa, contactá al soporte de la empresa con la referencia de cotización y organización, sólo por un canal autorizado; no compartas contraseñas ni tokens.",
    "en": "An order is confirmed here when the query shows its reference. If the problem continues, contact company support with the quote and organization references through an authorized channel only; do not share passwords or tokens."
  },
  "p1068": {
    "es": "Práctica de esta versión en el entorno de capacitación: reconocer un resultado incierto, recuperar la lectura y ubicar el pedido sin reenviar. No practiques creando pedidos en producción.",
    "en": "Practice this version in the training environment: recognize an uncertain result, recover the query and locate the order without resubmitting. Do not practice by creating production orders."
  },
  "p1069": {
    "es": "Consultar el estado de WhatsApp",
    "en": "Check WhatsApp status"
  },
  "p1070": {
    "es": "Compartí esta referencia sólo con soporte autorizado de tu organización. No adjuntes teléfonos, tokens ni el contenido del mensaje.",
    "en": "Share this reference only with authorized support for your organization. Do not attach phone numbers, tokens or message content."
  },
  "p1071": {
    "es": "Si el resultado es incierto o contradictorio: preservá el intento, consultá su evidencia y no reenvíes ni borres el historial. La entrega informada no confirma una venta ni la aceptación del turno.",
    "en": "If the result is uncertain or contradictory, preserve the attempt, check its evidence and do not resend or erase history. Reported delivery does not confirm a sale or appointment acceptance."
  },
  "p1072": {
    "es": "Cancelar un intervalo",
    "en": "Cancel an interval"
  },
  "p1073": {
    "es": "Ayuda para cancelar un intervalo",
    "en": "Help with cancelling an interval"
  },
  "p1074": {
    "es": "Cancelar este registro no cancela ni reprograma citas de clientes. El servidor comprueba los permisos, la versión y las restricciones existentes. Revisá recurso, fechas y tipo antes de actuar.",
    "en": "Cancelling this record does not cancel or reschedule customer appointments. The server checks permissions, version and existing constraints. Review the resource, dates and type before acting."
  },
  "p1075": {
    "es": "Si se pierde la respuesta, Consultar intervalo sólo lee el servidor. Una cancelación consultada no demuestra quién la solicitó; la auditoría autorizada conserva al actor. La referencia de esta pestaña guarda sólo la versión, no motivos, horarios ni datos personales.",
    "en": "If the response is lost, Check interval only reads the server. A retrieved cancellation does not establish who requested it; authorized audit retains the actor. This tab's reference stores only the version, not reasons, times or personal data."
  },
  "p1076": {
    "es": "Si la consulta falla, el intervalo no aparece en la ventana consultada o permanece activo, conservamos el bloqueo: solicitá al soporte autorizado revisar la referencia, hora y estado. No borres la referencia ni abras otra pestaña para forzar un reintento.",
    "en": "If the query fails, the interval is outside the queried window or remains active, the lock is retained: ask authorized support to review the reference, time and status. Do not erase the reference or open another tab to force a retry."
  },
  "p1077": {
    "es": "Práctica en pruebas: cancelar un intervalo sin citas, perder la respuesta y consultar Cancelado sin reenviar. No practicar con la agenda real. Reactivar o crear otro intervalo es una operación distinta.",
    "en": "Practice in tests: cancel an interval with no appointments, lose the response and retrieve Cancelled without resubmitting. Do not practice with the live schedule. Reactivating or creating another interval is a separate operation."
  },
  "p1078": {
    "es": "Registrar un intervalo",
    "en": "Register an interval"
  },
  "p1079": {
    "es": "Ayuda para registrar un intervalo",
    "en": "Help with registering an interval"
  },
  "p1080": {
    "es": "Conservamos sólo la referencia de operación. Una respuesta perdida no demuestra rechazo: consultá antes de volver a actuar. La consulta muestra el registro aunque quede fuera del rango de la agenda. Preparar otro intervalo no reenvía ni cancela el anterior y sólo se habilita tras recuperarlo. Ante referencia ilegible o consulta fallida, pedí revisión autorizada sin borrar datos. Practicá con datos sintéticos el envío, pérdida de respuesta, consulta y preparación explícita de otro intervalo.",
    "en": "Only the operation reference is retained. A lost response does not establish rejection: check before acting again. The query shows the record even outside the schedule range. Preparing another interval does not resubmit or cancel the previous one and is enabled only after retrieval. If the reference is unreadable or the query fails, request an authorized review without erasing data. Practice submission, response loss, retrieval and explicit preparation of another interval with synthetic data."
  },
  "p1081": {
    "es": "Presentar una entrega",
    "en": "Submit a handover"
  },
  "p1082": {
    "es": "Ayuda para presentar una entrega",
    "en": "Help with submitting a handover"
  },
  "p1083": {
    "es": "La presentación conserva respuestas inmutables. Guardamos sólo identidades, versiones y una huella; no guardamos las respuestas ni evidencias en la referencia del navegador. Una respuesta perdida puede ocultar un cambio correcto: consultá sin reenviarlo. El estado actual puede avanzar a aceptado o rechazado; consultar o preparar otro formulario no lo cambia. Si el contenido no coincide, pedí revisión. Practicá con datos sintéticos antes de operar.",
    "en": "Submission retains immutable answers. Only identities, versions and a fingerprint are saved; answers and evidence are not stored in the browser reference. A lost response can hide a successful change: check without resubmitting. The current status may advance to accepted or rejected; querying or preparing another form does not change it. If the content does not match, request a review. Practice with synthetic data before operating."
  },
  "p1084": {
    "es": "Publicar una checklist",
    "en": "Publish a checklist"
  },
  "p1085": {
    "es": "Ayuda para publicar una checklist",
    "en": "Help with publishing a checklist"
  },
  "p1086": {
    "es": "Una versión publicada es inmutable. Conservamos sólo su ID, versión y huella de contenido; no guardamos el título ni los textos en la referencia del navegador. Una respuesta perdida puede ocultar una publicación correcta: consultá sin volver a enviarla. Si no coincide el contenido o no se puede recuperar, pedí revisión sin borrar la referencia. Preparar otra versión no cambia ni vuelve a publicar la anterior. Practicá este recorrido con datos sintéticos antes de usarlo en operación.",
    "en": "A published version is immutable. Only its ID, version and content fingerprint are retained; the title and text are not stored in the browser reference. A lost response can hide a successful publication: check without resubmitting. If the content does not match or cannot be retrieved, request a review without erasing the reference. Preparing another version does not change or republish the previous one. Practice this journey with synthetic data before operating."
  },
  "p1087": {
    "es": "Resolver una discrepancia",
    "en": "Resolve a discrepancy"
  },
  "p1088": {
    "es": "Ayuda para resolver una discrepancia",
    "en": "Help with resolving a discrepancy"
  },
  "p1089": {
    "es": "Una respuesta perdida no demuestra rechazo. Conservamos en esta pestaña sólo la versión y la decisión pendiente, sin notas ni credenciales. Consultar vuelve a leer el servidor sin reenviar la acción.",
    "en": "A lost response does not establish rejection. This tab retains only the version and pending decision, without notes or credentials. Checking reads the server again without resubmitting the action."
  },
  "p1090": {
    "es": "Si la consulta falla, el registro no aparece o hay otra decisión, pedí revisión al responsable autorizado. No borres la referencia ni abras otra pestaña para forzar un reintento. Informá la referencia de la entrega y la hora, sin compartir notas privadas ni tokens.",
    "en": "If the query fails, the record is absent or another decision exists, request a review by an authorized owner. Do not erase the reference or open another tab to force a retry. Provide the handover reference and time without sharing private notes or tokens."
  },
  "p1091": {
    "es": "Práctica: en el entorno de prueba, recuperá una resolución tras perder la respuesta y comprobá su estado antes del siguiente paso. Preparada o autorizada no significa entrega física, reembolso ni cambio finalizados.",
    "en": "Practice: in the test environment, retrieve a resolution after losing the response and check its status before the next step. Prepared or authorized does not mean physical handover, refund or exchange is complete."
  },
  "p1092": {
    "es": "Consultar entregas",
    "en": "Browse handovers"
  },
  "p1093": {
    "es": "Ayuda para consultar entregas",
    "en": "Help with browsing handovers"
  },
  "p1094": {
    "es": "Si la información no está disponible, volvé a consultar antes de aceptar o rechazar. No crees otra solicitud para recuperar una lectura.",
    "en": "If information is unavailable, query again before accepting or rejecting. Do not create another request to recover a read."
  },
  "p1095": {
    "es": "Si el problema continúa, contactá al soporte autorizado e indicá la acción y la hora. No envíes contraseñas, tokens ni información de otras personas.",
    "en": "If the problem continues, contact authorized support and specify the action and time. Do not send passwords, tokens or other people's information."
  },
  "p1096": {
    "es": "Práctica en pruebas: interrumpí la consulta, verificá que no haya acciones de entrega disponibles y restablecé la conexión. La próxima consulta debe recuperar el estado sin enviar una aceptación automática.",
    "en": "Practice in tests: interrupt the query, verify that no handover actions are available and restore the connection. The next query must retrieve the status without automatically sending an acceptance."
  },
  "p1097": {
    "es": "Actualizar un lead",
    "en": "Update a lead"
  },
  "p1098": {
    "es": "Ayuda para actualizar un lead",
    "en": "Help with updating a lead"
  },
  "p1099": {
    "es": "Una respuesta perdida no demuestra que el cambio haya sido rechazado. Consultar lee el servidor sin reenviar comandos. Una versión posterior permite revisar el estado real; no prueba que haya sido modificado exclusivamente por tu solicitud.",
    "en": "A lost response does not establish that the change was rejected. Checking reads the server without resubmitting commands. A later version lets you review the actual status; it does not prove it was changed exclusively by your request."
  },
  "p1100": {
    "es": "Conservamos sólo versión y tipo de acción en esta pestaña, no responsables ni datos de contacto. Si falta el lead, falla su consulta o la versión no avanzó, pedí revisión al soporte autorizado; no borres la referencia ni abras otra pestaña para forzar un reintento.",
    "en": "This tab retains only the version and action type, not owners or contact details. If the lead is absent, its query fails or its version has not advanced, request a review by authorized support; do not erase the reference or open another tab to force a retry."
  },
  "p1101": {
    "es": "Práctica: recuperá un cambio en el entorno de prueba, comprobá responsable y estado y continuá con la versión consultada. No asumas una venta o un cliente convertido por una asignación.",
    "en": "Practice: retrieve a change in the test environment, check the owner and status and continue with the retrieved version. Do not infer a sale or converted customer from an assignment."
  },
  "p1102": {
    "es": "Recuperar una sección",
    "en": "Recover a section"
  },
  "p1103": {
    "es": "Ayuda para recuperar una sección",
    "en": "Help with recovering a section"
  },
  "p1104": {
    "es": "Una consulta fallida no significa que no existan registros. Volvé a consultar antes de decidir sobre esa sección; las secciones independientes conservan sus propios permisos y estados.",
    "en": "A failed query does not mean no records exist. Query again before making decisions about that section; independent sections retain their own permissions and states."
  },
  "p1105": {
    "es": "La consulta no reenvía acciones. Si el problema persiste, informá al soporte autorizado la sección, la hora y la operación esperada; no compartas tokens, contraseñas ni datos personales en capturas.",
    "en": "The query does not resubmit actions. If the problem persists, tell authorized support the section, time and expected operation; do not share tokens, passwords or personal data in screenshots."
  },
  "p1106": {
    "es": "Práctica: ante una consulta fallida, identificá la sección, recuperala mediante consulta y verificá su estado antes de actuar.",
    "en": "Practice: after a failed query, identify the section, recover it through a query and verify its status before acting."
  },
  "p1107": {
    "es": "Continuar un pedido",
    "en": "Continue an order"
  },
  "p1108": {
    "es": "Ayuda para continuar un pedido",
    "en": "Help with continuing an order"
  },
  "p1109": {
    "es": "Elegí una unidad por su serie y reservá una sola vez. El servidor vuelve a comprobar organización, variante, disponibilidad y versiones.",
    "en": "Choose a unit by its serial number and reserve once. The server rechecks the organization, variant, availability and versions."
  },
  "p1110": {
    "es": "Si se pierde la respuesta o alguien actuó primero, usá Actualizar pedidos. No crees otro pedido ni cambies identificadores para forzar la operación.",
    "en": "If the response is lost or someone acted first, use Refresh orders. Do not create another order or change identifiers to force the operation."
  },
  "p1111": {
    "es": "Pago: registrar una solicitud guarda una intención del total; no confirma dinero recibido ni libera entrega. Una solicitud existente se consulta, no se duplica. Si falla, no borres la clave del navegador: consultá el estado y escalá al responsable antes de intentar otro pago.",
    "en": "Payment: registering a request saves an intention for the total; it does not confirm receipt of money or release a handover. An existing request is queried, not duplicated. If it fails, do not erase the browser key: check the status and escalate to the owner before attempting another payment."
  },
  "p1112": {
    "es": "Práctica segura: en pruebas, dos operadores solicitan el pago del mismo pedido; sólo debe existir una intención inicial. Simulá respuesta perdida y recuperá el registro con Actualizar pedidos. No practiques con pedidos reales.",
    "en": "Safe practice: in tests, two operators request payment for the same order; only one initial intention should exist. Simulate a lost response and retrieve the record with Refresh orders. Do not practice with real orders."
  },
  "p1113": {
    "es": "Para soporte, comunicá la referencia del pedido o solicitud, acción y hora al responsable autorizado. No adjuntes tokens ni datos de clientes. Cobro, conciliación y entrega mantienen sus gates separados.",
    "en": "For support, provide the order or request reference, action and time to the authorized owner. Do not attach tokens or customer data. Collection, reconciliation and handover retain separate gates."
  },
  "p1114": {
    "es": "Emitir una cotización",
    "en": "Issue a quote"
  },
  "p1115": {
    "es": "Ayuda para emitir una cotización",
    "en": "Help with issuing a quote"
  },
  "p1116": {
    "es": "El servidor determina precio y moneda. La referencia conserva sólo un identificador de operación. Una respuesta perdida no permite emitir otra: consultá el resultado. Si no aparece, conservá la referencia y pedí revisión al responsable autorizado. Practicá emisión, respuesta perdida y consulta con datos sintéticos; no borres la referencia para repetir una cotización.",
    "en": "The server determines price and currency. The reference retains only an operation identifier. A lost response does not permit issuing another quote: check the result. If it is absent, retain the reference and request a review by the authorized owner. Practice issuance, response loss and retrieval with synthetic data; do not erase the reference to repeat a quote."
  },
  "p1117": {
    "es": "Registrar un recurso",
    "en": "Register a resource"
  },
  "p1118": {
    "es": "Ayuda para registrar un recurso",
    "en": "Help with registering a resource"
  },
  "p1119": {
    "es": "Conservamos sólo la referencia de operación. Una respuesta perdida no demuestra rechazo: consultá antes de volver a actuar. La consulta identifica el registro por la referencia de esta operación, sin buscar por nombre. Preparar otro recurso no reenvía ni cancela el anterior y sólo se habilita tras recuperarlo. Ante referencia ilegible o consulta fallida, pedí revisión autorizada sin borrar datos. Practicá con datos sintéticos el envío, pérdida de respuesta, consulta y preparación explícita de otro recurso.",
    "en": "Only the operation reference is retained. A lost response does not establish rejection: check before acting again. The query identifies the record by this operation's reference, without searching by name. Preparing another resource does not resubmit or cancel the previous one and is enabled only after retrieval. If the reference is unreadable or the query fails, request an authorized review without erasing data. Practice submission, response loss, retrieval and explicit preparation of another resource with synthetic data."
  },
  "p1120": {
    "es": "Recibir y decidir devoluciones",
    "en": "Receive and decide returns"
  },
  "p1121": {
    "es": "Ayuda para recibir y decidir devoluciones",
    "en": "Help with receiving and deciding returns"
  },
  "p1122": {
    "es": "Recepción, decisión y solicitudes conservan evidencia inmutable. La referencia del navegador guarda sólo operación, IDs y huella; no guarda serie, notas ni tokens. Una respuesta perdida puede ocultar un cambio correcto. Consultá el caso exacto aunque la lista falle o no lo muestre. Si la evidencia no coincide, pedí revisión. Continuar operando actualiza el caso consultado y habilita una acción nueva; no reenvía la anterior. Una devolución con reembolso crea cuatro solicitudes; un cambio crea tres y no solicita emisión fiscal. Practicá con datos sintéticos antes de operar.",
    "en": "Receipt, decision and requests retain immutable evidence. The browser reference stores only the operation, IDs and fingerprint; it does not store the serial number, notes or tokens. A lost response can hide a successful change. Check the exact case even if the list fails or omits it. If the evidence does not match, request a review. Continuing operations refreshes the queried case and enables a new action; it does not resubmit the previous one. A return with a refund creates four requests; an exchange creates three and does not request tax issuance. Practice with synthetic data before operating."
  },
  "p1123": {
    "es": "Registrar un turno",
    "en": "Register a slot"
  },
  "p1124": {
    "es": "Ayuda para registrar un turno",
    "en": "Help with registering a slot"
  },
  "p1125": {
    "es": "Conservamos sólo la referencia de operación. Una respuesta perdida no demuestra rechazo: consultá antes de volver a actuar. La consulta muestra el turno por su referencia, aunque esté cerrado, completo o fuera del listado público. Publicar no confirma una reserva; siguen vigentes la jornada y las restricciones de agenda. Preparar otro turno no reenvía ni cancela el anterior y sólo se habilita tras recuperarlo. Ante referencia ilegible o consulta fallida, pedí revisión autorizada sin borrar datos. Practicá con datos sintéticos el envío, pérdida de respuesta, consulta y preparación explícita de otro turno.",
    "en": "Only the operation reference is retained. A lost response does not establish rejection: check before acting again. The query shows the slot by its reference even when closed, full or outside the public listing. Publishing does not confirm a booking; working hours and schedule constraints still apply. Preparing another slot does not resubmit or cancel the previous one and is enabled only after retrieval. If the reference is unreadable or the query fails, request an authorized review without erasing data. Practice submission, response loss, retrieval and explicit preparation of another slot with synthetic data."
  },
  "p1126": {
    "es": "Guía de suministro",
    "en": "Supply guide"
  },
  "p1127": {
    "es": "Guía de suministro · versión 1.0.0",
    "en": "Supply guide · version 1.0.0"
  },
  "p1128": {
    "es": "Compras registra proveedor, destino, fábrica, importe y cantidades. La fábrica confirma la orden, registra cada serie y pide revisión de calidad. Otra persona autorizada decide su liberación o rechazo con evidencia.",
    "en": "Purchasing records the supplier, destination, factory, amount and quantities. The factory confirms the order, registers each serial number and requests quality review. Another authorized person decides release or rejection with evidence."
  },
  "p1129": {
    "es": "El despacho identifica un envío y sus series. Recepción registra sólo lo recibido; queda en cuarentena hasta la decisión de otra persona autorizada. Las cantidades y versiones se comprueban en cada operación.",
    "en": "Dispatch identifies a shipment and its serial numbers. Receiving records only what arrived; it remains quarantined until another authorized person decides. Quantities and versions are checked in every operation."
  },
  "p1130": {
    "es": "Si se pierde una respuesta, consultá el resultado pendiente. Esa consulta no vuelve a enviar. No borres referencias para forzar un reintento. En capacitación, practicá con datos sintéticos el recorrido y la recuperación.",
    "en": "If a response is lost, check the pending result. That query does not resend. Do not erase references to force a retry. In training, practice the journey and recovery with synthetic data."
  },
  "p1131": {
    "es": "Guía de garantía",
    "en": "Warranty guide"
  },
  "p1132": {
    "es": "Guía de garantía · versión 1.0.0",
    "en": "Warranty guide · version 1.0.0"
  },
  "p1133": {
    "es": "Leé los términos vinculados a la cotización antes de acusar su recepción. La garantía activa conserva esos términos vendidos. Un reclamo necesita entrega y atención registradas.",
    "en": "Read the terms bound to the quote before acknowledging receipt. The active warranty retains those sold terms. A claim requires recorded handover and service."
  },
  "p1134": {
    "es": "Diagnóstico, plan y aprobación preceden al trabajo. La persona revisora es diferente de quien propone; calidad es diferente de quien ejecuta. El cliente acepta y la fábrica concilia el servicio. La conciliación es un acuse interno, no un pago.",
    "en": "Diagnosis, plan and approval precede work. The reviewer differs from the proposer; quality review differs from execution. The customer accepts and the factory reconciles the service. Reconciliation is an internal acknowledgment, not a payment."
  },
  "p1135": {
    "es": "Si se pierde una respuesta, consultá el resultado pendiente: no vuelve a enviar. Conservá su referencia y pedí revisión autorizada si sigue incierto. Practicá con datos sintéticos en capacitación.",
    "en": "If a response is lost, check the pending result: it does not resend. Retain its reference and request an authorized review if it remains uncertain. Practice with synthetic data in training."
  },
  "p1136": {
    "es": "Guía de red y acuerdos",
    "en": "Network and agreements guide"
  },
  "p1137": {
    "es": "Guía de red y acuerdos · versión 1.0.0",
    "en": "Network and agreements guide · version 1.0.0"
  },
  "p1138": {
    "es": "Registrá la organización bajo un padre autorizado y usá la referencia devuelta para consultar su estado. Crear una organización no concede accesos a personas. Los nuevos accesos se asignan mediante el proveedor de identidad autorizado.",
    "en": "Register the organization under an authorized parent and use the returned reference to check its status. Creating an organization does not grant people access. New access is assigned through the authorized identity provider."
  },
  "p1139": {
    "es": "La estructura activa permite preparar acuerdos y sucursales. No equivale a autorización de despliegue o producción. Territorio, versión de términos y fechas deben proceder del acuerdo revisado; no se generan condiciones comerciales aquí.",
    "en": "The active structure allows preparation of agreements and branches. It does not authorize deployment or production. Territory, terms version and dates must come from the reviewed agreement; no commercial terms are generated here."
  },
  "p1140": {
    "es": "Antes de operar, completá la capacitación de tu rol. La evaluación no concede permisos. Si una respuesta se pierde, recuperá su resultado y luego consultá el estado actual. No repitas una escritura incierta.",
    "en": "Before operating, complete your role training. Assessment does not grant permissions. If a response is lost, recover its result and then check the current status. Do not repeat an uncertain write."
  },
  "p1141": {
    "es": "Para cerrar una organización, revisá primero sus sucursales y acuerdos activos. No se eliminan los registros ni se revocan sesiones desde este formulario.",
    "en": "Before closing an organization, review its active branches and agreements. This form does not delete records or revoke sessions."
  },
  "p1142": {
    "es": "Guía de contenidos",
    "en": "Content guide"
  },
  "p1143": {
    "es": "Guía de contenidos · versión 1.0.0",
    "en": "Content guide · version 1.0.0"
  },
  "p1144": {
    "es": "El contenido pertenece a una organización y un idioma. Guardá un borrador, revisalo y publicalo con el permiso correspondiente. Archivar retira el artículo de la consulta de lectores y conserva su historial.",
    "en": "Content belongs to an organization and a language. Save a draft, review it and publish it with the corresponding permission. Archiving removes the article from reader queries and preserves its history."
  },
  "p1145": {
    "es": "Las versiones publicadas no se editan. Para reemplazarlas, prepará un artículo nuevo y archivá el anterior después de revisar la sustitución. El texto se muestra literalmente; no ejecuta HTML.",
    "en": "Published versions are not edited. To replace one, prepare a new article and archive the previous one after reviewing the replacement. Text is displayed literally; it does not execute HTML."
  },
  "p1146": {
    "es": "Publicar aquí no modifica cursos aprobados ni concede permisos. La incorporación a capacitación requiere revisión y una nueva versión del curso.",
    "en": "Publishing here does not change approved courses or grant permissions. Inclusion in training requires review and a new course version."
  },
  "p1147": {
    "es": "Guía de catálogo y publicación",
    "en": "Catalog and publication guide"
  },
  "p1148": {
    "es": "Guía de catálogo · versión 1.0.0",
    "en": "Catalog guide · version 1.0.0"
  },
  "p1149": {
    "es": "Creá modelo, variantes y precios con importes y vigencias explícitos. Conservá las referencias devueltas, incorporá la imagen PNG y prepará una versión del catálogo.",
    "en": "Create a model, variants and prices with explicit amounts and validity periods. Retain the returned references, add the PNG image and prepare a catalog version."
  },
  "p1150": {
    "es": "La versión requiere las revisiones legal, técnica, de imagen y de publicación. Cada decisión refiere al contenido exacto y su evidencia. Publicar comprueba la generación vigente; una revisión desactualizada no autoriza contenido nuevo.",
    "en": "The version requires legal, technical, image and publication reviews. Each decision refers to the exact content and its evidence. Publishing checks the current generation; an outdated review does not authorize new content."
  },
  "p1151": {
    "es": "Consultar un resultado pendiente recupera la operación sin repetirla. Antes de publicar otra versión o volver a una anterior, consultá la publicación vigente. La publicación del catálogo no confirma stock, cobro ni habilitación comercial. Practicá únicamente con datos sintéticos.",
    "en": "Checking a pending result retrieves the operation without repeating it. Before publishing another version or restoring a previous one, check the current publication. Catalog publication does not confirm stock, collection or commercial authorization. Practice only with synthetic data."
  },
  "p1152": {
    "es": "Guía de capacitación y evaluación",
    "en": "Training and assessment guide"
  },
  "p1153": {
    "es": "Guía de capacitación · versión 1.0.0",
    "en": "Training guide · version 1.0.0"
  },
  "p1154": {
    "es": "Elegí la guía del rol e iniciá una práctica. Registrá la lectura de cada lección y presentá respuestas para revisión humana. Elegir una guía no cambia los permisos de tu sesión.",
    "en": "Choose the role guide and start a practice. Record reading each lesson and submit answers for human review. Choosing a guide does not change your session permissions."
  },
  "p1155": {
    "es": "Una persona autorizada diferente revisa las respuestas contra la versión presentada. Una evaluación favorable conserva evidencia, pero no asigna accesos ni reemplaza la aceptación del responsable de la operación.",
    "en": "A different authorized person reviews the answers against the submitted version. A favorable assessment retains evidence but does not assign access or replace acceptance by the operations owner."
  },
  "p1156": {
    "es": "Ante una respuesta incierta, consultá el intento guardado. Un curso nuevo no cambia la evidencia de intentos anteriores; prepará otra práctica para la versión activa. Los artículos privados del CMS requieren incorporación explícita a un curso revisado.",
    "en": "After an uncertain response, check the saved attempt. A new course does not change evidence from previous attempts; prepare another practice for the active version. Private CMS articles require explicit inclusion in a reviewed course."
  },
  "p1157": {
    "es": "Accesos de supervisión disponibles en tu sesión.",
    "en": "Supervision access available in your session."
  },
  "p1158": {
    "es": "Accesos disponibles para la operación diaria.",
    "en": "Access available for daily operations."
  },
  "p1159": {
    "es": "Accesos disponibles para atención y operación.",
    "en": "Access available for service and operations."
  },
  "p1160": {
    "es": "Accesos disponibles para tus consultas y entregas.",
    "en": "Access available for your queries and handovers."
  },
  "p1161": {
    "es": "inconsistent observations",
    "en": "inconsistent observations"
  },
  "p1162": {
    "es": "missing acceptance timestamp",
    "en": "missing acceptance timestamp"
  },
  "p1163": {
    "es": "unknown outcome requires reconciliation",
    "en": "unknown outcome requires reconciliation"
  },
  "p1164": {
    "es": "duplicate confirmation",
    "en": "duplicate confirmation"
  },
  "p1165": {
    "es": "Resultado por conciliar",
    "en": "Result awaiting reconciliation"
  },
  "p1166": {
    "es": "No reenvíes el mensaje. Soporte debe revisar el intento y su evidencia antes de autorizar otra acción.",
    "en": "Do not resend the message. Support must review the attempt and its evidence before authorizing another action."
  },
  "p1167": {
    "es": "Estados contradictorios",
    "en": "Contradictory statuses"
  },
  "p1168": {
    "es": "Hay avisos distintos con la misma fecha del proveedor. Conservá la referencia y pedí revisión; no elijas uno por orden de llegada.",
    "en": "Different notices have the same provider timestamp. Retain the reference and request a review; do not choose one by arrival order."
  },
  "p1169": {
    "es": "El proveedor informó envío",
    "en": "The provider reported sending"
  },
  "p1170": {
    "es": "El proveedor informó entrega",
    "en": "The provider reported delivery"
  },
  "p1171": {
    "es": "El proveedor informó lectura",
    "en": "The provider reported reading"
  },
  "p1172": {
    "es": "El proveedor informó un fallo",
    "en": "The provider reported a failure"
  },
  "p1173": {
    "es": "El proveedor informó eliminación",
    "en": "The provider reported deletion"
  },
  "p1174": {
    "es": "Este aviso no acredita una venta ni la aceptación del turno por el cliente. No vuelve a enviar mensajes.",
    "en": "This notice does not certify a sale or the customer's appointment acceptance. It does not resend messages."
  },
  "p1175": {
    "es": "Aceptado por el proveedor; entrega sin verificar",
    "en": "Accepted by the provider; delivery unverified"
  },
  "p1176": {
    "es": "El envío fue aceptado, pero todavía no hay un aviso verificado de entrega. Actualizá la consulta; no reenvíes por falta de confirmación.",
    "en": "The send was accepted, but there is no verified delivery notice yet. Refresh the query; do not resend because confirmation is absent."
  },
  "p1177": {
    "es": "Envío en proceso",
    "en": "Send in progress"
  },
  "p1178": {
    "es": "Esperá y actualizá la consulta. No ejecutes un segundo envío mientras el intento siga abierto.",
    "en": "Wait and refresh the query. Do not perform a second send while the attempt remains open."
  },
  "p1179": {
    "es": "Intento cerrado con fallo",
    "en": "Attempt closed with failure"
  },
  "p1180": {
    "es": "Revisá la evidencia con soporte. Esta pantalla no habilita reintentos ni cambia el historial.",
    "en": "Review the evidence with support. This screen does not enable retries or change history."
  },
  "p1181": {
    "es": "Sin intento de envío registrado",
    "en": "No send attempt recorded"
  },
  "p1182": {
    "es": "Una aprobación no implica un envío. Esta consulta es sólo de lectura; no inicia mensajes ni captura consentimiento.",
    "en": "Approval does not imply sending. This query is read-only; it does not initiate messages or capture consent."
  },
  "p1183": {
    "es": "Solicitud inválida",
    "en": "Invalid request"
  },
  "p1184": {
    "es": "Revise los datos enviados.",
    "en": "Review the submitted data."
  },
  "p1185": {
    "es": "Backend no disponible",
    "en": "Backend unavailable"
  },
  "p1186": {
    "es": "No pudimos completar la operación.",
    "en": "We could not complete the operation."
  },
  "p1187": {
    "es": "Error interno",
    "en": "Internal error"
  },
  "p1188": {
    "es": "La operación no pudo completarse.",
    "en": "The operation could not be completed."
  },
  "p1189": {
    "es": "Mi panel",
    "en": "My dashboard"
  },
  "p1190": {
    "es": "Editar catálogo",
    "en": "Edit catalog"
  },
  "p1191": {
    "es": "Modelos",
    "en": "Models"
  },
  "p1192": {
    "es": "Dónde estamos",
    "en": "Locations"
  },
  "p1193": {
    "es": "Presioná Buscar para consultar.",
    "en": "Press Search to query."
  },
  "p1194": {
    "es": "No pudimos cargar esta información",
    "en": "We could not load this information"
  },
  "p1195": {
    "es": "La consulta no se completó. Podés volver a consultar con tu sesión actual.",
    "en": "The query did not complete. You can query again with your current session."
  },
  "p1196": {
    "es": "Este mensaje no confirma el resultado de una operación anterior. Consultá su estado antes de repetirla.",
    "en": "This message does not confirm a previous operation result. Check its status before repeating it."
  },
  "p1197": {
    "es": "Volver al inicio",
    "en": "Back to first page"
  },
  "p1198": {
    "es": "Ver siguientes",
    "en": "View next page"
  }
}
````

### FILE: `src/platform/i18n/private-provider.tsx`

```yaml
block_id: "TS-GO-API-WEB-BRIDGE-LOCALE-DELTA:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "282df2429a29b65502b39825e24ac0b05a36c394be86c6df8ab96781e4e99363"
variables: []
secrets_allowed: false
```

````tsx
"use client";
import{createContext,useContext,useMemo,type ReactNode}from"react";
import{privateTranslator,controlledPrivateLabel}from"./private-catalog";
import type{PrivateLocale}from"./private-locale";
const fallback:PrivateLocale={language:"es",locale:"es-AR",timeZone:"UTC",source:"configuration"};
const Context=createContext<PrivateLocale>(fallback);
export function PrivateLocaleProvider({locale,children}:{locale:PrivateLocale;children:ReactNode}){return <Context.Provider value={locale}>{children}</Context.Provider>}
export function usePrivateI18n(){const locale=useContext(Context);return useMemo(()=>({locale,t:privateTranslator(locale.language),controlled:(value:string)=>controlledPrivateLabel(locale.language,value)}),[locale])}
````


V402 composed delta: Conditional document/AWS original notices; no AWS dependency claim in web-only profiles. Both alternative notice owners use the same reviewed text.

V402 composed delta: Local delivery316: native stop owner reused, Next standalone exact build identity, container template without mutable defaults and complete local module context. Local fixture qualification only; docs/LOCAL_REFERENCE_DELIVERY.md.
