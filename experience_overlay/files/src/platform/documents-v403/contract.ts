// AUTHORED UI/BFF projection of internal/documentbridge and httpapi/document.go.
// No new document class, money interpretation, scanner or approval policy.
import { z } from "zod";

export const MAX_DOCUMENT_BYTES = 2 * 1024 * 1024;
export const documentId = z.uuid().refine(v => v !== "00000000-0000-0000-0000-000000000000");
export const hash = z.string().regex(/^[a-f0-9]{64}$/u);
const text = (bytes: number) => z.string().refine(v => v.length > 0 && v.trim() === v && !v.includes("\0") && new TextDecoder().decode(new TextEncoder().encode(v)) === v && new TextEncoder().encode(v).length <= bytes);
export const fields = z.object({ invoice_number: text(128), vendor: text(512), total: text(64), currency: z.string().regex(/^[A-Z]{3}$/u) }).strict();
const suggested = z.object({ invoice_number: z.string(), vendor: z.string(), total: z.string(), currency: z.string() }).strict();
export const viewSchema = z.object({
  document_id: documentId, name: text(128), uploader: text(128), original_sha256: hash, profile_sha256: hash,
  mode: z.enum(["FIXTURE", "PROVIDER"]), state: z.enum(["QUARANTINED", "REVIEW_REQUIRED", "REVIEW_PENDING", "REJECTED", "PERSISTED", "QUARANTINE_TERMINAL"]),
  evidence_sha256: hash.optional(), suggested: suggested.optional(),
  proposal: z.object({ schema: z.literal("document-review/v1"), document_id: documentId, organization_id: text(128), original_sha256: hash, profile_sha256: hash, evidence_sha256: hash, mode: z.enum(["FIXTURE", "PROVIDER"]), fields }).strict().optional(),
  payload_sha256: hash.optional(), reviewer: text(128).optional(), reason: z.string().optional()
}).strict().superRefine((v, ctx) => {
  if (v.proposal && (v.proposal.document_id !== v.document_id || v.proposal.original_sha256 !== v.original_sha256 || v.proposal.profile_sha256 !== v.profile_sha256 || v.proposal.evidence_sha256 !== v.evidence_sha256 || v.proposal.mode !== v.mode || !v.payload_sha256)) ctx.addIssue({ code: "custom", message: "Proposal binding mismatch" });
  if (["REVIEW_PENDING", "PERSISTED", "REJECTED"].includes(v.state) && !v.proposal) ctx.addIssue({ code: "custom", message: "Proposal required" });
  if (["PERSISTED", "REJECTED"].includes(v.state) && !v.reviewer) ctx.addIssue({ code: "custom", message: "Decision evidence required" });
});
export const commandSchema = z.discriminatedUnion("action", [
  z.object({ id: documentId, action: z.literal("process") }).strict(),
  z.object({ id: documentId, action: z.literal("review"), evidence_sha256: hash, fields }).strict(),
  z.object({ id: documentId, action: z.literal("decision"), payload_sha256: hash, approved: z.boolean(), reason: text(2048) }).strict()
]);
export type DocumentView = z.infer<typeof viewSchema>;
export type DocumentFields = z.infer<typeof fields>;
export type DocumentCommand = z.infer<typeof commandSchema>;
export type DocumentAccess = { scope: string; subject: string; canWrite: boolean; canProcess: boolean; canReview: boolean; mode: "FIXTURE" | "PROVIDER" };

export async function sha256(bytes: Uint8Array): Promise<string> {
  return Array.from(new Uint8Array(await crypto.subtle.digest("SHA-256", new Uint8Array(bytes).buffer)), b => b.toString(16).padStart(2, "0")).join("");
}
export async function fieldsHash(value: DocumentFields): Promise<string> {
  return sha256(new TextEncoder().encode(JSON.stringify({ invoice_number: value.invoice_number, vendor: value.vendor, total: value.total, currency: value.currency })));
}

export type PendingDocument = { id: string; action: "receive" | "process" | "review" | "decision"; expected: string; approved?: boolean };
export async function reconciles(pending: PendingDocument, view: DocumentView): Promise<boolean> {
  if (pending.id !== view.document_id) return false;
  if (pending.action === "receive") return view.original_sha256 === pending.expected;
  if (pending.action === "process") return view.state !== "QUARANTINED" && view.original_sha256 === pending.expected;
  if (pending.action === "review") return !!view.proposal && await fieldsHash(view.proposal.fields) === pending.expected;
  return view.payload_sha256 === pending.expected && (pending.approved === true ? view.state === "PERSISTED" : pending.approved === false && view.state === "REJECTED");
}

export function parsePending(raw: string): PendingDocument {
  const parsed = z.object({ id: documentId, action: z.enum(["receive", "process", "review", "decision"]), expected: hash, approved: z.boolean().optional() }).strict()
    .refine(v => v.action === "decision" ? typeof v.approved === "boolean" : v.approved === undefined).parse(JSON.parse(raw));
  return { id: parsed.id, action: parsed.action, expected: parsed.expected, ...(typeof parsed.approved === "boolean" ? { approved: parsed.approved } : {}) };
}
