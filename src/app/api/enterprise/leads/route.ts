import { boundedCommandBody } from "@/platform/backend/bounded-command";
import { z } from "zod";
import { captureLead } from "@/platform/backend/public-client";
import { applicationBaseUrl } from "@/platform/auth/oidc-client";
import { errorResponse, problem } from "@/platform/http/problem";
const schema = z.object({ modelId: z.string().min(1).max(128).optional(), name: z.string().trim().min(2).max(120), email: z.email().max(254), consentGranted: z.literal(true) }).strict();
export async function POST(request: Request) {
  // AUTHORED narrow correction using the existing configured-origin owner.
  // Reject before parsing the body, idempotency transport or any business effect.
  // Neither the request Host nor forwarding headers supply the trusted origin.
  try { if (request.headers.get("origin") !== applicationBaseUrl().origin) return problem(403,"Origen no permitido","La consulta debe enviarse desde esta aplicación.","CROSS_ORIGIN_REJECTED"); }
  catch { return problem(403,"Origen no permitido","La consulta no está disponible.","CROSS_ORIGIN_REJECTED"); }
  if (request.headers.get("content-type") !== "application/json") return problem(415,"Formato no permitido","Envíe JSON.","UNSUPPORTED_MEDIA_TYPE");
  const key = request.headers.get("idempotency-key");
  if (!key || key.length < 16 || key.length > 128) return problem(400, "Falta idempotencia", "Envíe una clave estable para este intento.", "IDEMPOTENCY_KEY_REQUIRED");
  try {
    let raw: string;
    try { raw = await boundedCommandBody(request,4096); } catch (error) { return problem(error instanceof Error && error.message==="BODY_TOO_LARGE" ? 413 : 400,"Solicitud inválida","No se procesó la solicitud.",error instanceof Error && error.message==="BODY_TOO_LARGE" ? "COMMAND_TOO_LARGE" : "INVALID_BODY"); }
    let value: unknown;
    try { value = JSON.parse(raw); } catch { return problem(400,"Solicitud inválida","Revise los datos enviados.","INVALID_JSON"); }
    const input = schema.parse(value);
    const result = await captureLead({ ...(input.modelId ? { modelId: input.modelId } : {}), sourceCode: "public-web", contact: { name: input.name, email: input.email.toLowerCase() }, consentGranted: true }, key);
    return Response.json(result.value, { status: result.status, headers: { "cache-control": "no-store", ...(result.replayed ? { "idempotency-replayed": "true" } : {}) } });
  } catch (error) { return errorResponse(error); }
}
