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
