import { backendFetch } from "./cloud-run-transport";
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
  const response = await backendFetch(target, { ...init, cache: "no-store", redirect: "error", signal: AbortSignal.timeout(5000), headers: { accept: "application/json", ...init?.headers } });
  const value = await readBackendResponse<T>(response);
  return { value, status: response.status, replayed: response.headers.get("idempotency-replayed") === "true" };
}
export async function listModels(): Promise<PublicModel[]> { const { tenant } = settings(); const result = await call<{ models: PublicModel[] }>(`/v1/public/${encodeURIComponent(tenant)}/models`); return result.value.models; }
export async function captureLead(input: LeadInput, idempotencyKey: string) { if (idempotencyKey.length < 16 || idempotencyKey.length > 128) throw new Error("invalid idempotency key"); const { tenant, organization } = settings(); return call<{ lead_id: string; status: string }>(`/v1/public/${encodeURIComponent(tenant)}/${encodeURIComponent(organization)}/leads`, { method: "POST", headers: { "content-type": "application/json", "idempotency-key": idempotencyKey }, body: JSON.stringify({ model_id: input.modelId ?? "", source_code: input.sourceCode, contact: input.contact, consent_granted: input.consentGranted }) }); }
export async function listLocations(): Promise<PublicLocation[]> { const { tenant } = settings(); const result = await call<{ items: PublicLocation[] }>(`/v1/public/${encodeURIComponent(tenant)}/locations`); return result.value.items; }
export async function listAppointmentSlots(kind: AppointmentKind, from: string, to: string): Promise<AppointmentSlot[]> { if (!Number.isFinite(Date.parse(from)) || !Number.isFinite(Date.parse(to)) || Date.parse(to) <= Date.parse(from)) throw new Error("invalid appointment slot range"); const { tenant, organization } = settings(); const query = new URLSearchParams({ kind, from, to }); const result = await call<{ items: AppointmentSlot[] }>(`/v1/public/${encodeURIComponent(tenant)}/${encodeURIComponent(organization)}/appointment-slots?${query}`); return result.value.items; }
export async function requestAppointment(input: AppointmentInput, idempotencyKey: string) { if (idempotencyKey.length < 16 || idempotencyKey.length > 128) throw new Error("invalid idempotency key"); const { tenant, organization } = settings(); return call<{ id: string; state: string; starts_at: string }>(`/v1/public/${encodeURIComponent(tenant)}/${encodeURIComponent(organization)}/appointments`, { method: "POST", headers: { "content-type": "application/json", "idempotency-key": idempotencyKey }, body: JSON.stringify({ lead_id: input.leadId, model_id: input.modelId ?? "", kind: input.kind, starts_at: input.startsAt }) }); }
