// AUTHORED BFF glue. Existing operations.Service remains the transition authority.
import {NextResponse}from"next/server";
import{z}from"zod";
import{allowed,readSession}from"@/platform/auth/session";
import{applicationBaseUrl}from"@/platform/auth/oidc-client";
import{protectedGet,protectedPost}from"@/platform/backend/protected-client";
import{BackendProblem}from"@/platform/backend/public-client";
import{factoryID,factoryUnitSchema}from"@/platform/factory/contracts";
const scope={organizationId:factoryID,unitId:factoryID};
const query=z.object(scope).strict();
const command=z.object({...scope,action:z.literal("start-assembly"),expectedState:z.literal("planned")}).strict();
const reply=(body:unknown,status=200)=>NextResponse.json(body,{status,headers:{"Cache-Control":"no-store","X-Content-Type-Options":"nosniff"}});
function problem(e:unknown){return reply({code:"FACTORY_UNAVAILABLE"},e instanceof BackendProblem&&[400,401,403,404,409,503].includes(e.status)?e.status:503)}
async function boundedBody(request: Request) {
 if (!request.body) throw new Error("invalid factory body");
 const reader = request.body.getReader(); const chunks: Uint8Array[] = []; let bytes = 0; let timer: ReturnType<typeof setTimeout> | undefined;
 const deadline = new Promise<never>((_, reject) => { timer = setTimeout(() => reject(new Error("body deadline")), 2500); });
 try {
  while (true) {
   const part = await Promise.race([reader.read(), deadline]); if (part.done) break;
   bytes += part.value.byteLength; if (bytes > 4096) throw new Error("BODY_TOO_LARGE"); chunks.push(part.value);
  }
  const data = new Uint8Array(bytes); let offset = 0; for (const chunk of chunks) { data.set(chunk, offset); offset += chunk.byteLength; }
  return new TextDecoder("utf-8", { fatal: true }).decode(data);
 } finally { clearTimeout(timer); void reader.cancel().catch(() => {}); reader.releaseLock(); }
}

export async function GET(request:Request){
 if(["cross-site","none"].includes(request.headers.get("sec-fetch-site")??""))return reply({code:"CROSS_ORIGIN_REJECTED"},403);
 const session=await readSession();if(!session)return reply({code:"UNAUTHENTICATED"},401);if(!allowed(session,"factory:read"))return reply({code:"FORBIDDEN"},403);
 const search=new URL(request.url).searchParams;if(request.url.length>2048||[...search.keys()].some(k=>search.getAll(k).length!==1))return reply({code:"INVALID_QUERY"},400);
 const parsed=query.safeParse(Object.fromEntries(search));if(!parsed.success)return reply({code:"INVALID_QUERY"},400);const q=parsed.data;if(!session.organizations.includes(q.organizationId))return reply({code:"FORBIDDEN"},403);
 try{const unit=factoryUnitSchema.parse(await protectedGet<unknown>(session,`/v1/factory/units/${encodeURIComponent(q.unitId)}`,{organization_id:q.organizationId}));if(unit.id!==q.unitId||unit.organization_id!==q.organizationId)return reply({code:"CONTEXT_MISMATCH"},502);return reply({unit,observationOnly:true})}catch(e){return problem(e)}
}
export async function POST(request:Request){
 try{if(request.headers.get("origin")!==applicationBaseUrl().origin)return reply({code:"CROSS_ORIGIN_REJECTED"},403)}catch{return reply({code:"CROSS_ORIGIN_REJECTED"},403)}
 if(request.headers.get("content-type")!=="application/json")return reply({code:"UNSUPPORTED_MEDIA_TYPE"},415);
 const session=await readSession();if(!session)return reply({code:"UNAUTHENTICATED"},401);if(!allowed(session,"factory:read")||!allowed(session,"factory:write"))return reply({code:"FORBIDDEN"},403);
 let body:unknown;try{body=JSON.parse(await boundedBody(request))}catch(e){return reply({code:e instanceof Error&&e.message==="BODY_TOO_LARGE"?"COMMAND_TOO_LARGE":"INVALID_COMMAND"},e instanceof Error&&e.message==="BODY_TOO_LARGE"?413:400)}
 const parsed=command.safeParse(body);if(!parsed.success)return reply({code:"INVALID_COMMAND"},400);const q=parsed.data;if(!session.organizations.includes(q.organizationId))return reply({code:"FORBIDDEN"},403);
 try{const unit=factoryUnitSchema.parse(await protectedGet<unknown>(session,`/v1/factory/units/${encodeURIComponent(q.unitId)}`,{organization_id:q.organizationId}));if(unit.id!==q.unitId||unit.organization_id!==q.organizationId)return reply({code:"CONTEXT_MISMATCH"},502);if(unit.state!==q.expectedState)return reply({code:"STATE_CHANGED"},409);
 const ack=z.object({status:z.literal("accepted")}).strict().parse(await protectedPost<unknown>(session,`/v1/factory/units/${encodeURIComponent(q.unitId)}/transitions`,{organization_id:q.organizationId,current:q.expectedState,target:"assembly"}));return reply(ack)
 }catch(e){return problem(e)}
}
