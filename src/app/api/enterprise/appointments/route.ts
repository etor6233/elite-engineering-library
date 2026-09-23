import { boundedCommandBody } from "@/platform/backend/bounded-command";
import { z } from "zod";
import { requestAppointment } from "@/platform/backend/public-client";
import { applicationBaseUrl } from "@/platform/auth/oidc-client";
import { errorResponse, problem } from "@/platform/http/problem";

const schema = z.object({
  leadId: z.string().min(1).max(128),
  modelId: z.string().min(1).max(128).optional(),
  kind: z.enum(["consultation", "test-drive", "delivery", "service"]),
  startsAt: z.iso.datetime({ offset: true })
}).strict();

export async function POST(request: Request) {
  // Same configured-origin guard as the public lead stage; no effect precedes it.
  try { if (request.headers.get("origin") !== applicationBaseUrl().origin) return problem(403,"Origen no permitido","El turno debe solicitarse desde esta aplicación.","CROSS_ORIGIN_REJECTED"); }
  catch { return problem(403,"Origen no permitido","La reserva no está disponible.","CROSS_ORIGIN_REJECTED"); }
  if (request.headers.get("content-type") !== "application/json") return problem(415,"Formato no permitido","Envíe JSON.","UNSUPPORTED_MEDIA_TYPE");
  const key = request.headers.get("idempotency-key");
  if (!key || key.length < 16 || key.length > 128) return problem(400, "Falta idempotencia", "Envíe una clave estable para este intento.", "IDEMPOTENCY_KEY_REQUIRED");
  try {
    let raw: string;
    try { raw = await boundedCommandBody(request,4096); } catch (error) { return problem(error instanceof Error && error.message==="BODY_TOO_LARGE" ? 413 : 400,"Solicitud inválida","No se procesó la solicitud.",error instanceof Error && error.message==="BODY_TOO_LARGE" ? "COMMAND_TOO_LARGE" : "INVALID_BODY"); }
    let value: unknown;
    try { value = JSON.parse(raw); } catch { return problem(400,"Solicitud inválida","Revise los datos enviados.","INVALID_JSON"); }
    const input = schema.parse(value);
    const result = await requestAppointment({ leadId: input.leadId, kind: input.kind, startsAt: input.startsAt, ...(input.modelId ? { modelId: input.modelId } : {}) }, key);
    return Response.json(result.value, { status: result.status, headers: { "cache-control": "no-store", ...(result.replayed ? { "idempotency-replayed": "true" } : {}) } });
  } catch (error) { return errorResponse(error); }
}
