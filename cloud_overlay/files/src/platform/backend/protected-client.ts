import { backendFetch } from "./cloud-run-transport";
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
  const response = await backendFetch(target, { cache: "no-store", redirect: "error", signal: AbortSignal.timeout(5000), headers: { accept: "application/json", authorization: `Bearer ${session.accessToken}`, ...(idempotencyKey ? { "Idempotency-Key": idempotencyKey } : {}) } });
  return readBackendResponse<T>(response);
}

export async function protectedPost<T>(session: PortalSession, path: string, command: unknown, idempotencyKey?: string): Promise<T> {
  if (!path.startsWith("/v1/") || path.includes("..")) throw new Error("protected backend path is invalid");
  if (idempotencyKey !== undefined && !/^[A-Za-z0-9_-]{16,128}$/.test(idempotencyKey)) throw new Error("invalid idempotency key");
  const body = JSON.stringify(command);
  if (body.length > 1_048_576) throw new Error("protected backend command is too large");
  const response = await backendFetch(new URL(path, baseUrl()), {
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
 const response=await backendFetch(new URL(path,baseUrl()),{method:"POST",cache:"no-store",redirect:"error",signal:AbortSignal.timeout(5000),headers:{accept:"application/json",authorization:`Bearer ${session.accessToken}`,"content-type":"image/png"},body:new Uint8Array(bytes).buffer});
 return readBackendResponse<T>(response);
}
