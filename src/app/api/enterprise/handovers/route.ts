import { boundedCommandBody as boundedBody } from "@/platform/backend/bounded-command";
// AUTHORED scoped BFF; obtains observation bindings server-to-server and forwards
// only supported commands. No browser-entered amount, customer or digest.
import { NextResponse } from "next/server";
import { z } from "zod";
import { allowed,readSession } from "@/platform/auth/session";
import { applicationBaseUrl } from "@/platform/auth/oidc-client";
import { protectedGet,protectedPost } from "@/platform/backend/protected-client";
import { BackendProblem } from "@/platform/backend/public-client";
import { backendHandoverContextSchema,handoverContextSchema,handoverID,preparationSchema,commercialReceiptSchema,currentReleaseSchema } from "@/platform/handovers/contracts";
const reply=(body:unknown,status=200)=>NextResponse.json(body,{status,headers:{"Cache-Control":"no-store","X-Content-Type-Options":"nosniff"}});
const base={organizationId:handoverID,orderId:handoverID};
const command=z.object({...base,action:z.enum(["prepare","release"]),requestKey:z.uuid()}).strict();
const query=z.discriminatedUnion("kind",[
 z.object({...base,kind:z.literal("context")}).strict(),
 z.object({...base,kind:z.literal("prepare-result"),requestKey:z.uuid()}).strict(),
 z.object({...base,kind:z.literal("release-result"),handoverId:handoverID,requestKey:z.uuid()}).strict(),
 z.object({...base,kind:z.literal("release-current"),handoverId:handoverID}).strict(),
]);
function problem(error:unknown){return reply({code:"HANDOVER_UNAVAILABLE"},error instanceof BackendProblem&&[400,401,403,404,409,503].includes(error.status)?error.status:503);}
export async function POST(request:Request){
 try{if(request.headers.get("origin")!==applicationBaseUrl().origin)return reply({code:"CROSS_ORIGIN_REJECTED"},403)}catch{return reply({code:"CROSS_ORIGIN_REJECTED"},403)}
 if(request.headers.get("content-type")!=="application/json")return reply({code:"UNSUPPORTED_MEDIA_TYPE"},415);
 const session=await readSession();if(!session)return reply({code:"UNAUTHENTICATED"},401);if(!allowed(session,"handover:manage"))return reply({code:"FORBIDDEN"},403);
 let body:unknown;try{body=JSON.parse(await boundedBody(request))}catch(error){return error instanceof Error&&error.message==="BODY_TOO_LARGE"?reply({code:"COMMAND_TOO_LARGE"},413):reply({code:"INVALID_COMMAND"},400)}
 const parsed=command.safeParse(body);if(!parsed.success)return reply({code:"INVALID_COMMAND"},400);const c=parsed.data;
 if(!session.organizations.includes(c.organizationId))return reply({code:"FORBIDDEN"},403);
 try{
  const v=backendHandoverContextSchema.parse(await protectedGet<unknown>(session,`/v1/franchise/orders/${encodeURIComponent(c.orderId)}/handover-context`,{organization_id:c.organizationId}));
  if(v.organization_id!==c.organizationId||v.order_id!==c.orderId)return reply({code:"CONTEXT_MISMATCH"},502);
  if(c.action==="prepare"){
   if(!v.can_prepare)return reply({code:"HANDOVER_ALREADY_PREPARED"},409);
   const result=preparationSchema.parse(await protectedPost<unknown>(session,`/v1/franchise/orders/${encodeURIComponent(c.orderId)}/handover`,{organization_id:c.organizationId,order_line_id:v.order_line_id,payment_attempt_id:v.payment_attempt_id,...(v.funding_receipt_id?{funding_receipt_id:v.funding_receipt_id}:{}),observation_sha256:v.observation_sha256},c.requestKey));
   if(result.handover.organization_id!==c.organizationId||result.handover.order_id!==c.orderId||result.order_line_id!==v.order_line_id||result.payment_attempt_id!==v.payment_attempt_id||result.funding_receipt_id!==v.funding_receipt_id||result.observation_sha256!==v.observation_sha256)return reply({code:"RECEIPT_MISMATCH"},502);
   return reply({handover:result.handover},201);
  }
  if(v.release_effect!=="COMMIT_COMMERCIAL_RELEASE_RECEIPT"||v.handover?.state!=="accepted")return reply({code:"RELEASE_NOT_ENABLED"},409);
  const receipt=commercialReceiptSchema.parse(await protectedPost<unknown>(session,`/v1/franchise/handovers/${encodeURIComponent(v.handover.id)}/commercial-release`,{organization_id:c.organizationId,observation_sha256:v.observation_sha256},c.requestKey));
  if(receipt.organization_id!==c.organizationId||receipt.order_id!==c.orderId||receipt.handover_id!==v.handover.id||receipt.observation_sha256!==v.observation_sha256)return reply({code:"RECEIPT_MISMATCH"},502);
  return reply({receipt},201);
 }catch(error){return problem(error)}
}
export async function GET(request:Request){
 if(["cross-site","none"].includes(request.headers.get("sec-fetch-site")??""))return reply({code:"CROSS_ORIGIN_REJECTED"},403);
 const session=await readSession();if(!session)return reply({code:"UNAUTHENTICATED"},401);if(!allowed(session,"handover:manage"))return reply({code:"FORBIDDEN"},403);
 const search=new URL(request.url).searchParams;if(request.url.length>2048||[...search.keys()].some(k=>search.getAll(k).length!==1))return reply({code:"INVALID_QUERY"},400);
 const parsed=query.safeParse(Object.fromEntries(search));if(!parsed.success)return reply({code:"INVALID_QUERY"},400);const q=parsed.data;
 if(!session.organizations.includes(q.organizationId))return reply({code:"FORBIDDEN"},403);
 try{
  if(q.kind==="context"){
   const value=handoverContextSchema.parse(await protectedGet<unknown>(session,`/v1/franchise/orders/${encodeURIComponent(q.orderId)}/handover-context`,{organization_id:q.organizationId}));
   if(value.organization_id!==q.organizationId||value.order_id!==q.orderId)return reply({code:"CONTEXT_MISMATCH"},502);return reply(value);
  }
  if(q.kind==="prepare-result"){
   const value=preparationSchema.parse(await protectedGet<unknown>(session,`/v1/franchise/orders/${encodeURIComponent(q.orderId)}/handover-result`,{organization_id:q.organizationId},q.requestKey));
   if(value.handover.organization_id!==q.organizationId||value.handover.order_id!==q.orderId)return reply({code:"RECEIPT_MISMATCH"},502);return reply({handover:value.handover});
  }
  const raw=await protectedGet<unknown>(session,`/v1/franchise/handovers/${encodeURIComponent(q.handoverId)}/${q.kind==="release-result"?"commercial-release-result":"commercial-release-current"}`,{organization_id:q.organizationId},q.kind==="release-result"?q.requestKey:undefined);
  const value=q.kind==="release-result"?{receipt:commercialReceiptSchema.parse(raw)}:currentReleaseSchema.parse(raw);
  if(value.receipt.organization_id!==q.organizationId||value.receipt.order_id!==q.orderId||value.receipt.handover_id!==q.handoverId)return reply({code:"RECEIPT_MISMATCH"},502);return reply(value);
 }catch(error){return problem(error)}
}
