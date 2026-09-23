import { NextResponse } from "next/server";
import { allowed, readSession } from "@/platform/auth/session";
import { protectedGet } from "@/platform/backend/protected-client";
import { BackendProblem } from "@/platform/backend/public-client";
import { notificationHistorySchema } from "@/platform/notifications/status-contract";
import { loadBusinessConfig } from "@/platform/config/load";

const reply = (value: unknown, status = 200) => NextResponse.json(value, { status, headers: { "Cache-Control": "no-store", "X-Content-Type-Options": "nosniff" } });
const safeID = (value: string | null): value is string => value !== null && /^[a-zA-Z0-9][a-zA-Z0-9_-]{0,127}$/.test(value);

// Read-only same-site BFF. Token/tenant are resolved server-side, never returned.
export async function GET(request: Request) {
  try {
    if ((await loadBusinessConfig()).features.whatsapp_status_history !== true) return reply({ code: "CAPABILITY_NOT_ENABLED" }, 404);
  } catch { return reply({ code: "NOTIFICATION_READ_UNAVAILABLE" }, 503); }
  if (["cross-site", "none"].includes(request.headers.get("sec-fetch-site") ?? "")) return reply({ code: "CROSS_SITE_REJECTED" }, 403);
  const session = await readSession();
  if (!session) return reply({ code: "UNAUTHENTICATED" }, 401);
  if (!allowed(session, "appointment:manage")) return reply({ code: "FORBIDDEN" }, 403);
  const query = new URL(request.url).searchParams;
  const organization = query.get("organizationId"), appointment = query.get("appointmentId");
  if (request.url.length > 2048 || [...query.keys()].length !== 2 || query.getAll("organizationId").length !== 1 || query.getAll("appointmentId").length !== 1 || !safeID(organization) || !safeID(appointment)) return reply({ code: "INVALID_QUERY" }, 400);
  if (!session.organizations.includes(organization) && !session.organizations.includes("*")) return reply({ code: "ORGANIZATION_FORBIDDEN" }, 403);
  try {
    const raw = await protectedGet<unknown>(session, `/v1/franchise/appointments/${encodeURIComponent(appointment)}/whatsapp-confirmations`, { organization_id: organization });
    const parsed = notificationHistorySchema.safeParse(raw);
    if (!parsed.success || parsed.data.organization_id !== organization || parsed.data.appointment_id !== appointment) return reply({ code: "INVALID_NOTIFICATION_RESPONSE" }, 502);
    return reply(parsed.data);
  } catch (error) {
    const permitted = new Set(["NOTIFICATION_NOT_FOUND", "NOTIFICATION_EVIDENCE_MISMATCH", "NOTIFICATION_HISTORY_LIMIT", "NOTIFICATION_READ_UNAVAILABLE", "UNAUTHENTICATED", "FORBIDDEN", "ORGANIZATION_FORBIDDEN"]);
    if (error instanceof BackendProblem && permitted.has(error.code) && [401, 403, 404, 409, 503].includes(error.status)) return reply({ code: error.code }, error.status);
    return reply({ code: "NOTIFICATION_READ_UNAVAILABLE" }, 503);
  }
}
