import { describe, it, expect, vi } from "vitest";
import { createHash } from "node:crypto";
import { createDocumentGateway, type DocumentConfig, type DocumentTransport } from "./gateway";
import { commandSchema, fields, fieldsHash, parsePending, reconciles, viewSchema, type DocumentView } from "./contract";
import type { PortalSession } from "@/platform/auth/session";

const id = "11111111-1111-4111-8111-111111111111", other = "22222222-2222-4222-8222-222222222222";
const bytes = new Uint8Array([255,216,255,1,2,3]); // synthetic transport bytes, not an OCR fixture
const sha = (v: Uint8Array) => createHash("sha256").update(v).digest("hex");
const original = sha(bytes), profile = "a".repeat(64), evidence = "b".repeat(64), payload = "c".repeat(64);
const sample = { invoice_number: "000045", vendor: "Fixture supplier", total: "1.234,50", currency: "ARS" };
const config: DocumentConfig = { origin: "https://app.example.test", tenant: "tenant-fixture", organization: "store", profile, mode: "FIXTURE" };
const writer: PortalSession = { subject: "writer", tenantId: config.tenant, organizations: [config.organization], permissions: ["documents:write", "documents:process"], accessToken: "fixture-not-a-provider-token" };
const base: DocumentView = { document_id: id, name: "invoice.jpg", uploader: "writer", original_sha256: original, profile_sha256: profile, mode: "FIXTURE", state: "QUARANTINED" };
const proposal = { schema: "document-review/v1" as const, document_id: id, organization_id: "store", original_sha256: original, profile_sha256: profile, evidence_sha256: evidence, mode: "FIXTURE" as const, fields: sample };
const pendingView: DocumentView = { ...base, state: "REVIEW_PENDING", suggested: sample, evidence_sha256: evidence, proposal, payload_sha256: payload };
const confirmed: DocumentView = { ...pendingView, state: "PERSISTED", reviewer: "reviewer", reason: "Compared with original" };
const reply = (v: unknown, status = 200) => Response.json(v, { status });
const post = (value: unknown) => new Request(config.origin + "/api/documents", { method: "POST", headers: { origin: config.origin, "content-type": "application/json" }, body: JSON.stringify(value) });
const put = (body = bytes, digest = sha(body), name = "invoice.jpg") => new Request(config.origin + "/api/documents?id=" + id, { method: "PUT", headers: { origin: config.origin, "content-type": "application/octet-stream", "x-document-name": name, "x-document-sha256": digest }, body: new Uint8Array(body).buffer });
function setup(session: PortalSession | null = writer, transport = vi.fn<DocumentTransport>(async () => reply(base))) {
  return { transport, handle: createDocumentGateway({ session: async () => session, config: () => config, transport }) };
}

describe("document adapter: existing owner contracts, simulated provider only", () => {
  it("receives exact bytes/profile and preserves the generated document identity", async () => {
    const s = setup(); const response = await s.handle(put()); expect(response.status).toBe(200);
    const [, path, request] = s.transport.mock.calls[0]!;
    expect(path).toBe(`/v1/documents/${id}/original`); expect(request.method).toBe("PUT"); expect(request.body).toEqual(bytes);
    expect(request.headers?.["X-Document-Profile-SHA256"]).toBe(profile);
  });
  it("denies unauthenticated, other organization and other tenant before transport", async () => {
    for (const session of [null, { ...writer, organizations: ["foreign"] }, { ...writer, tenantId: "foreign" }]) {
      const s = setup(session); expect((await s.handle(put())).status).toBe(session ? 403 : 401); expect(s.transport).not.toHaveBeenCalled();
    }
  });
  it("does not grant reviewer or document writer permissions via UI", async () => {
    const s = setup(); expect((await s.handle(post({ action: "decision", id, payload_sha256: payload, approved: true, reason: "review" }))).status).toBe(403); expect(s.transport).not.toHaveBeenCalled();
    const reader = setup({ ...writer, permissions: ["documents:review"] }); expect((await reader.handle(put())).status).toBe(403); expect(reader.transport).not.toHaveBeenCalled();
  });
  it("rejects cross-origin and query/path injection before sending", async () => {
    const s = setup(); const request = put(); request.headers.set("origin", "https://foreign.example.test"); expect((await s.handle(request)).status).toBe(403);
    for (const query of ["?id=../original", `?id=${id}&id=${other}`, `?id=${id}&part=../../users`, `?id=${id}&tenant=foreign`]) expect((await s.handle(new Request(config.origin + query))).status).toBe(400);
    expect(s.transport).not.toHaveBeenCalled();
  });
  it("checks streamed byte limit and digest, not just Content-Length", async () => {
    const s = setup(); expect((await s.handle(put(bytes, "f".repeat(64)))).status).toBe(400);
    const huge = put(new Uint8Array(2 * 1024 * 1024 + 1)); huge.headers.set("content-length", "1"); expect((await s.handle(huge)).status).toBe(413); expect(s.transport).not.toHaveBeenCalled();
  });
  it("preserves encoded Unicode names, rejects paths and unsupported file content", async () => {
    const s = setup(); expect((await s.handle(put(bytes, original, encodeURIComponent("factura-á.jpg")))).ok).toBe(true);
    for (const name of ["..%2finvoice.jpg", "folder%5cfile.jpg", "invoice.png", "invalid%FF.jpg"]) expect((await s.handle(put(bytes, original, name))).status).toBe(400);
    expect((await s.handle(put(new TextEncoder().encode("not a jpeg")))).status).toBe(400); expect(s.transport).toHaveBeenCalledTimes(1);
  });
  it("keeps amount text and leading zeroes; never manufactures a missing field", () => {
    expect(fields.parse(sample)).toEqual(sample); expect(fields.safeParse({ ...sample, currency: "ars" }).success).toBe(false);
    expect(fields.safeParse({ ...sample, invoice_number: "" }).success).toBe(false); expect(fields.safeParse({ ...sample, vendor: "ñ".repeat(257) }).success).toBe(false);
    expect(commandSchema.safeParse({ action: "decision", id, payload_sha256: payload, approved: "true", reason: "review" }).success).toBe(false);
  });
  it("connects receive/process/proposal/separate decision without an extra commit route", async () => {
    let view = { ...base }; let writes = 0;
    const transport = vi.fn<DocumentTransport>(async (session, path, request) => {
      if (request.method === "GET") return reply(view);
      writes++;
      if (path.endsWith("/process")) view = { ...base, state: "REVIEW_REQUIRED", suggested: sample, evidence_sha256: evidence };
      if (path.endsWith("/review")) view = pendingView;
      if (path.endsWith("/decision")) {
        if (session.subject === view.uploader) return reply({ code: "CONSULT_RECORDED_STATE" }, 409);
        view = confirmed;
      }
      return reply(view);
    });
    const s = setup(writer, transport); await s.handle(put()); await s.handle(post({ action: "process", id }));
    expect((await s.handle(post({ action: "review", id, evidence_sha256: evidence, fields: sample }))).ok).toBe(true);
    const reviewer = setup({ ...writer, subject: "reviewer", permissions: ["documents:review"] }, transport);
    const result = await reviewer.handle(post({ action: "decision", id, payload_sha256: payload, approved: true, reason: "Compared with original" }));
    expect((await result.json()).state).toBe("PERSISTED"); expect(writes).toBe(4);
    expect(transport.mock.calls.some(([,path]) => path.endsWith("/commit"))).toBe(false);
  });
  it("passes owner separation rejection through; it cannot become an approval", async () => {
    const transport = vi.fn<DocumentTransport>(async () => reply({ code: "CONSULT_RECORDED_STATE" }, 409));
    const s = setup({ ...writer, permissions: ["documents:review"] }, transport);
    const response = await s.handle(post({ action: "decision", id, payload_sha256: payload, approved: true, reason: "attempt" }));
    expect(response.status).toBe(409); expect((await response.json()).effect).toBe("UNCONFIRMED");
  });
  it("does not mark an ambiguous write or malformed upstream receipt successful", async () => {
    for (const transport of [vi.fn<DocumentTransport>(async () => { throw new Error("response lost"); }), vi.fn<DocumentTransport>(async () => reply({ ...base, document_id: other })), vi.fn<DocumentTransport>(async () => reply({ ...pendingView, proposal: { ...proposal, profile_sha256: "f".repeat(64) } }))]) {
      const s = setup(writer, transport); const response = await s.handle(put()); expect(response.status).toBe(503); expect((await response.json()).effect).toBe("UNCONFIRMED"); expect(transport).toHaveBeenCalledTimes(1);
    }
  });
  it("returns original/evidence as hash-checked attachments, never active HTML", async () => {
    const transport = vi.fn<DocumentTransport>(async () => new Response(bytes, { headers: { "X-Document-SHA256": original } })); const s = setup(writer, transport);
    const response = await s.handle(new Request(config.origin + `?id=${id}&part=original`));
    expect(response.headers.get("content-type")).toBe("application/octet-stream"); expect(response.headers.get("content-disposition")).toContain("attachment"); expect(new Uint8Array(await response.arrayBuffer())).toEqual(bytes);
    transport.mockImplementation(async () => new Response(bytes, { headers: { "X-Document-SHA256": "d".repeat(64) } }));
    expect((await s.handle(new Request(config.origin + `?id=${id}&part=original`))).status).toBe(503);
  });
  it("resumes only from matching recorded state after a response is lost", async () => {
    expect(await reconciles({ id, action: "receive", expected: original }, base)).toBe(true);
    expect(await reconciles({ id, action: "process", expected: original }, base)).toBe(false);
    expect(await reconciles({ id, action: "process", expected: original }, pendingView)).toBe(true);
    expect(await reconciles({ id, action: "review", expected: await fieldsHash(sample) }, pendingView)).toBe(true);
    expect(await reconciles({ id, action: "decision", expected: payload, approved: true }, pendingView)).toBe(false);
    expect(await reconciles({ id, action: "decision", expected: payload, approved: true }, confirmed)).toBe(true);
    expect(await reconciles({ id, action: "decision", expected: payload, approved: false }, confirmed)).toBe(false);
    expect(() => parsePending(JSON.stringify({ id, action: "review", expected: original, vendor: "private" }))).toThrow();
    expect(viewSchema.safeParse({ ...confirmed, reviewer: undefined }).success).toBe(false);
  });
});
