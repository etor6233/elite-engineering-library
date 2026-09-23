// ADAPTED local V403 fixed-route adapter, versioned for class catalog and inbox; backend remains the document/tenant/reviewer owner.
import { createHash } from "node:crypto";
import { commandSchema, documentId, hash, classId, schemaVersion, catalogSchema, inboxSchema, MAX_DOCUMENT_BYTES, viewSchema } from "./contract";
import type { PortalSession } from "@/platform/auth/session";
import { allowed } from "@/platform/auth/session";
import { boundedCommandBody } from "@/platform/backend/bounded-command";
import { BackendProblem } from "@/platform/backend/public-client";

export type DocumentConfig = { origin: string; tenant: string; organization: string; profile: string; mode: "FIXTURE" | "PROVIDER" | "TYPED_FIXTURE" };
export type DocumentTransport = (session: PortalSession, path: string, init: { method: string; body?: Uint8Array | string; headers?: Record<string, string> }) => Promise<Response>;
type Dependencies = { session: () => Promise<PortalSession | null>; config: () => DocumentConfig; transport: DocumentTransport };
const json = (body: unknown, status = 200) => Response.json(body, { status, headers: { "cache-control": "no-store", "x-content-type-options": "nosniff" } });
const noEffect = (code: string, status: number) => json({ code, effect: "NOT_ATTEMPTED" }, status);
const sha = (bytes: Uint8Array) => createHash("sha256").update(bytes).digest("hex");

// ADAPTED from the existing boundedRequestBytes owner: same finite read/cancel
// behavior, widened only to this owner's documented 2MiB original limit.
export async function boundedDocumentBytes(request: Request | Response, limit: number): Promise<Uint8Array> {
  if (!request.body) throw new Error("INVALID_BODY");
  const reader = request.body.getReader(); const chunks: Uint8Array[] = []; let length = 0;
  let timer: ReturnType<typeof setTimeout> | undefined;
  const deadline = new Promise<never>((_, reject) => { timer = setTimeout(() => reject(new Error("BODY_DEADLINE")), 5000); });
  try {
    for (;;) {
      const part = await Promise.race([reader.read(), deadline]); if (part.done) break;
      length += part.value.byteLength; if (length > limit) throw new Error("BODY_TOO_LARGE"); chunks.push(part.value);
    }
    const result = new Uint8Array(length); let offset = 0;
    for (const chunk of chunks) { result.set(chunk, offset); offset += chunk.length; }
    return result;
  } finally { clearTimeout(timer); void reader.cancel().catch(() => {}); reader.releaseLock(); }
}

export function createDocumentGateway(deps: Dependencies) {
  return async function handle(request: Request): Promise<Response> {
    const session = await deps.session(); if (!session) return noEffect("UNAUTHENTICATED", 401);
    const can = (p: string) => allowed(session, `documents:${p}`);
    if (!["read", "write", "process", "review"].some(can)) return noEffect("FORBIDDEN", 403);
    let config: DocumentConfig;
    try { config = deps.config(); } catch { return noEffect("DOCUMENTS_NOT_CONFIGURED", 503); }
    if (session.tenantId !== config.tenant || !session.organizations.includes(config.organization)) return noEffect("FORBIDDEN", 403);
    const url = new URL(request.url);
    if (request.method !== "GET" && request.headers.get("origin") !== config.origin) return noEffect("CROSS_ORIGIN_REJECTED", 403);
    for (const key of url.searchParams.keys()) if (!["id", "part", "view", "limit", "cursor"].includes(key) || url.searchParams.getAll(key).length !== 1) return noEffect("INVALID_QUERY", 400);
    if (request.method === "GET" && !url.search) {
      return json({ scope: sha(new TextEncoder().encode(JSON.stringify([session.tenantId, config.organization, session.subject, config.profile]))), subject: session.subject, canRead: ["read","write","review","process"].some(can), profile_sha256: config.profile, canWrite: can("write"), canProcess: can("process"), canReview: can("review"), mode: config.mode });
    }
    let id = "", action = "read", method = request.method, body: string | Uint8Array | undefined, headers: Record<string,string> = {}, path: string;
    try {
      if (request.method === "POST") {
        if (url.search || request.headers.get("content-type") !== "application/json") return noEffect("INVALID_COMMAND", 400);
        const command = commandSchema.parse(JSON.parse(await boundedCommandBody(request, 32768)));
        id = command.id; action = command.action;
        if (!can(action === "decision" ? "review" : action === "review" ? "write" : "process")) return noEffect("FORBIDDEN", 403);
        const { id: _id, action: _action, ...payload } = command; body = JSON.stringify(payload); headers["content-type"] = "application/json";
        path = `/v1/documents/${id}/${action}`;
      } else {
        const view = url.searchParams.get("view");
        if (request.method === "GET" && view) {
          if (url.searchParams.has("id") || url.searchParams.has("part") || !["classes","inbox"].includes(view)) return noEffect("INVALID_QUERY",400);
          if (view === "classes") {
            if (url.searchParams.has("limit") || url.searchParams.has("cursor")) return noEffect("INVALID_QUERY",400);
            path = "/v1/documents/classes";
          } else {
            const rawLimit=url.searchParams.get("limit")??"25", cursor=url.searchParams.get("cursor");
            if (!/^(?:[1-9]|[1-4][0-9]|50)$/u.test(rawLimit) || cursor!==null && (!cursor || cursor.length>2048 || /[\u0000-\u0020\u007f]/u.test(cursor))) return noEffect("INVALID_QUERY",400);
            path = "/v1/documents?limit="+rawLimit+(cursor?"&cursor="+encodeURIComponent(cursor):"");
          }
        } else {
        if (view || url.searchParams.has("limit") || url.searchParams.has("cursor")) return noEffect("INVALID_QUERY",400);
        id = documentId.parse(url.searchParams.get("id"));
        if (request.method === "PUT") {
          if (!can("write")) return noEffect("FORBIDDEN", 403);
          if (url.searchParams.has("part") || request.headers.get("content-type") !== "application/octet-stream") return noEffect("INVALID_UPLOAD", 400);
          const name = request.headers.get("x-document-name") ?? "";
          // Names travel percent-encoded ASCII so UTF-8 is lossless in HTTP headers.
          if (!/^[\x21-\x7e]{1,128}$/u.test(name) || /[\\/:]/u.test(name) || !/\.(pdf|jpe?g)$/iu.test(name)) return noEffect("INVALID_FILENAME", 400);
          try { const decoded = decodeURIComponent(name); if (/[\\/:\x00-\x1f\x7f]/u.test(decoded) || decoded.trim() !== decoded) return noEffect("INVALID_FILENAME", 400); } catch { return noEffect("INVALID_FILENAME", 400); }
          const expected = hash.parse(request.headers.get("x-document-sha256"));
          body = await boundedDocumentBytes(request, MAX_DOCUMENT_BYTES);
          if (body.length === 0 || sha(body) !== expected) return noEffect("INVALID_UPLOAD", 400);
          const jpeg = body[0] === 255 && body[1] === 216 && body[2] === 255;
          const pdf = new TextDecoder().decode(body.subarray(0, 5)) === "%PDF-";
          if (!(jpeg && /\.jpe?g$/iu.test(name) || pdf && /\.pdf$/iu.test(name))) return noEffect("INVALID_UPLOAD", 400);
          // This sniff is an input contract, never malware clearance or OCR admission.
          headers = { "content-type": "application/octet-stream", "X-Document-Name": name, "X-Document-SHA256": expected, "X-Document-Profile-SHA256": config.profile,
            "X-Document-Class": classId.parse(request.headers.get("x-document-class")), "X-Document-Schema-Version": schemaVersion.parse(request.headers.get("x-document-schema-version")) };
          path = `/v1/documents/${id}/original`; action = "receive";
        } else if (request.method === "GET") {
          const part = url.searchParams.get("part");
          if (part && !["original", "security", "provider", "analysis"].includes(part)) return noEffect("INVALID_QUERY", 400);
          path = `/v1/documents/${id}${part ? part === "original" ? "/original" : `/evidence/${part}` : ""}`;
        } else return noEffect("METHOD_NOT_ALLOWED", 405);
        }
      }
    } catch (e) { return noEffect(e instanceof Error && e.message === "BODY_TOO_LARGE" ? "DOCUMENT_TOO_LARGE" : "INVALID_COMMAND", e instanceof Error && e.message === "BODY_TOO_LARGE" ? 413 : 400); }
    try {
      const upstream = await deps.transport(session, path, { method, headers, ...(body === undefined ? {} : { body }) });
      const bytes = await boundedDocumentBytes(upstream, MAX_DOCUMENT_BYTES);
      if (!upstream.ok) {
        const status = [400, 401, 403, 404, 409, 422, 503].includes(upstream.status) ? upstream.status : 502;
        // Receive validates the original before opening its transaction. Other
        // provider failures may have recorded attempts; keep their uncertain fence.
        return json({ code: status === 409 ? "CONSULT_RECORDED_STATE" : status === 404 ? "NOT_FOUND" : "DOCUMENT_OPERATION_FAILED", effect: action === "receive" && status === 400 ? "NOT_ATTEMPTED" : "UNCONFIRMED" }, status);
      }
      const part = url.searchParams.get("part");
      if (method === "GET" && part) {
        const expected = upstream.headers.get(part === "original" ? "X-Document-SHA256" : "X-Document-Evidence-SHA256");
        if (!hash.safeParse(expected).success || sha(bytes) !== expected) throw new Error("EVIDENCE_HASH_MISMATCH");
        return new Response(bytes.buffer as ArrayBuffer, { headers: { "content-type": "application/octet-stream", "content-disposition": `attachment; filename="${part === "original" ? "original.bin" : "evidence.json"}"`, "cache-control": "no-store", "x-content-type-options": "nosniff", "X-Document-SHA256": expected! } });
      }
      const payload = JSON.parse(new TextDecoder("utf-8", { fatal: true }).decode(bytes));
      if (url.searchParams.get("view") === "classes") {
        const value=catalogSchema.parse(payload);
        if(value.profile_sha256!==config.profile || value.mode!==config.mode)throw new Error("CATALOG_SCOPE_MISMATCH");
        return json(value,upstream.status);
      }
      const matchesScope=(value:ReturnType<typeof viewSchema.parse>)=>value.profile_sha256===config.profile&&value.mode===config.mode&&(!value.proposal||value.proposal.organization_id===config.organization);
      if (url.searchParams.get("view") === "inbox") {
        const value=inboxSchema.parse(payload);if(!value.items.every(matchesScope))throw new Error("INBOX_SCOPE_MISMATCH");return json(value,upstream.status);
      }
      const value=viewSchema.parse(payload);
      if(value.document_id!==id||!matchesScope(value))throw new Error("DOCUMENT_SCOPE_MISMATCH");
      return json(value, upstream.status);
    } catch (e) { return json({ code: e instanceof BackendProblem ? e.code : "DOCUMENT_RESULT_UNCONFIRMED", effect: "UNCONFIRMED" }, 503); }
  };
}
