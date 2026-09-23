import { NextResponse } from "next/server";
import { allowed, readSession } from "@/platform/auth/session";
import { applicationBaseUrl } from "@/platform/auth/oidc-client";
import { loadBusinessConfig } from "@/platform/config/load";
import { protectedGet, protectedPost } from "@/platform/backend/protected-client";
import { BackendProblem } from "@/platform/backend/public-client";
import { answerSchema, definitionSchema, resultSchema, summarySchema, surveyID, parseSubmission } from "@/platform/surveys/contract";

const reply = (value: unknown, status = 200) => NextResponse.json(value, { status, headers: { "Cache-Control": "no-store", "X-Content-Type-Options": "nosniff" } });
async function scope(request: Request, writing: boolean) {
 try { if ((await loadBusinessConfig()).features.customer_surveys !== true) return reply({ code: "CAPABILITY_NOT_ENABLED" }, 404); }
 catch { return reply({ code: "SURVEY_UNAVAILABLE" }, 503); }
 if (["cross-site", "none"].includes(request.headers.get("sec-fetch-site") ?? "")) return reply({ code: "CROSS_SITE_REJECTED" }, 403);
 if (writing) {
  try { if (request.headers.get("origin") !== applicationBaseUrl().origin) return reply({ code: "CROSS_ORIGIN_REJECTED" }, 403); }
  catch { return reply({ code: "CROSS_ORIGIN_REJECTED" }, 403); }
 }
 const session = await readSession(); if (!session) return reply({ code: "UNAUTHENTICATED" }, 401);
 const query = new URL(request.url).searchParams;
 const organization = query.get("organizationId"), id = query.get("surveyId"), view = query.get("view");
 const keys = writing ? ["organizationId", "surveyId"] : ["organizationId", "surveyId", "view"];
 if (request.url.length > 2048 || [...query.keys()].length !== keys.length || keys.some(key => query.getAll(key).length !== 1) ||
     !surveyID.safeParse(organization).success || !surveyID.safeParse(id).success || (!writing && !["definition", "response", "summary"].includes(view ?? ""))) return reply({ code: "INVALID_QUERY" }, 400);
 if (!allowed(session, view === "summary" ? "surveys:read" : "customer:self")) return reply({ code: "FORBIDDEN" }, 403);
 if (!session.organizations.includes(organization!) && !session.organizations.includes("*")) return reply({ code: "ORGANIZATION_FORBIDDEN" }, 403);
 return { session, organization: organization!, id: id!, view };
}
function failure(error: unknown) {
 if (error instanceof BackendProblem && [400, 401, 403, 404, 409, 503].includes(error.status) &&
     ["INVALID_SURVEY_REQUEST", "SURVEY_RESPONSE_CONFLICT", "SURVEY_NOT_AVAILABLE", "SURVEY_UNAVAILABLE", "UNAUTHENTICATED", "FORBIDDEN", "ORGANIZATION_FORBIDDEN"].includes(error.code)) return reply({ code: error.code }, error.status);
 return reply({ code: "SURVEY_UNAVAILABLE" }, 503);
}
export async function GET(request: Request) {
 const value = await scope(request, false); if (value instanceof Response) return value;
 const { session, organization, id, view } = value;
 const path = view === "summary" ? "/v1/admin/surveys/" + encodeURIComponent(id) + "/summary" :
  "/v1/customer/surveys/" + encodeURIComponent(id) + (view === "response" ? "/response" : "");
 try {
  const raw = await protectedGet<unknown>(session, path, { organization_id: organization });
  const parsed = (view === "summary" ? summarySchema : view === "response" ? answerSchema : definitionSchema).safeParse(raw);
  if (!parsed.success || (view === "definition" && (!("id" in parsed.data) || parsed.data.id !== id))) return reply({ code: "INVALID_SURVEY_RESPONSE" }, 502);
  return reply(parsed.data);
 } catch (error) { return failure(error); }
}
async function boundedBody(request: Request) {
 if (!request.body) throw new Error("invalid survey body");
 const reader = request.body.getReader(); const chunks: Uint8Array[] = []; let bytes = 0; let timer: ReturnType<typeof setTimeout> | undefined;
 const deadline = new Promise<never>((_, reject) => { timer = setTimeout(() => reject(new Error("body deadline")), 2500); });
 try {
  while (true) {
   const part = await Promise.race([reader.read(), deadline]); if (part.done) break;
   bytes += part.value.byteLength; if (bytes > 2048) throw new Error("body limit"); chunks.push(part.value);
  }
  const data = new Uint8Array(bytes); let offset = 0; for (const chunk of chunks) { data.set(chunk, offset); offset += chunk.byteLength; }
  return new TextDecoder("utf-8", { fatal: true }).decode(data);
 } finally { clearTimeout(timer); void reader.cancel().catch(() => {}); reader.releaseLock(); }
}
export async function POST(request: Request) {
 const value = await scope(request, true); if (value instanceof Response) return value;
 const type = (request.headers.get("content-type") ?? "").split(";", 1)[0]?.trim().toLowerCase();
 if (type !== "application/json" || !["", "identity"].includes(request.headers.get("content-encoding") ?? "")) return reply({ code: "INVALID_CONTENT_TYPE" }, 415);
 let input;
 try { input = parseSubmission(await boundedBody(request)); } catch { return reply({ code: "INVALID_SURVEY_REQUEST" }, 400); }
 const { session, organization, id } = value;
 try {
  const raw = await protectedPost<unknown>(session, "/v1/customer/surveys/" + encodeURIComponent(id) + "/response?organization_id=" + encodeURIComponent(organization), input);
  const parsed = resultSchema.safeParse(raw);
  if (!parsed.success || parsed.data.answer.score !== input.score || parsed.data.answer.consent_version !== input.consent_version) return reply({ code: "INVALID_SURVEY_RESPONSE" }, 502);
  return reply(parsed.data, parsed.data.replay ? 200 : 201);
 } catch (error) { return failure(error); }
}

