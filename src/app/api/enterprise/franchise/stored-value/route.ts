// AUTHORED scoped BFF. Exact JSON payload strings are hashed and forwarded as
// text; no amount is converted to Number or sent back as a business decision.
import { NextResponse } from "next/server";
import { createHash } from "node:crypto";
import { allowed,readSession } from "@/platform/auth/session";
import { applicationBaseUrl } from "@/platform/auth/oidc-client";
import { protectedGet,protectedPost } from "@/platform/backend/protected-client";
import { BackendProblem } from "@/platform/backend/public-client";
import { boundedCommandBody } from "@/platform/backend/bounded-command";
import { storedCommand,storedQuery,storedOrder,storedApproval,storedFunding } from "@/platform/stored-value/contracts";
const base="/v1/franchise/stored-value";
const reply=(value:unknown,status=200)=>NextResponse.json(value,{status,headers:{"cache-control":"no-store","x-content-type-options":"nosniff"}});
const digest=(s:string)=>createHash("sha256").update(s,"utf8").digest("hex");
function failure(e:unknown){return reply({code:"STORED_VALUE_UNCONFIRMED"},e instanceof BackendProblem&&[400,401,403,404,409,503].includes(e.status)?e.status:503)}
function approval(raw:unknown,org:string){const v=storedApproval.parse(raw);if(v.organization_id!==org||digest(v.payload_json)!==v.payload_sha256)throw new Error("exact approval mismatch");return v}
function funding(raw:unknown,org:string){const v=storedFunding.parse(raw);if(v.organization_id!==org||digest(v.receipt_json)!==v.receipt_sha256)throw new Error("exact receipt mismatch");return v}
export async function POST(request:Request){
 try{if(request.headers.get("origin")!==applicationBaseUrl().origin)return reply({code:"CROSS_ORIGIN_REJECTED"},403)}catch{return reply({code:"CROSS_ORIGIN_REJECTED"},403)}
 if(request.headers.get("content-type")!=="application/json")return reply({code:"UNSUPPORTED_MEDIA_TYPE"},415);
 const session=await readSession();if(!session)return reply({code:"UNAUTHENTICATED"},401);
 let parsed;try{parsed=storedCommand.safeParse(JSON.parse(await boundedCommandBody(request)))}catch{return reply({code:"INVALID_COMMAND"},400)}
 if(!parsed.success)return reply({code:"INVALID_COMMAND"},400);const c=parsed.data;
 const permission=c.action==="propose"?"stored_value:request":c.action==="decision"?"stored_value:approve":"stored_value:fund";
 if(!session.organizations.includes(c.organization_id)||!allowed(session,permission))return reply({code:"FORBIDDEN"},403);
 try{
  if(c.action==="propose"){
   const {action,...body}=c;void action;
   const v=approval(await protectedPost<unknown>(session,`${base}/proposals`,body),c.organization_id);
   if(v.order_id!==c.order_id||v.approval_id!=="svapproval_"+digest(c.operation_id))throw new Error("proposal reference mismatch");return reply(v,201);
  }
  if(c.action==="decision"){
   const v=approval(await protectedPost<unknown>(session,`${base}/approvals/${encodeURIComponent(c.approval_id)}/decision`,{organization_id:c.organization_id,payload_sha256:c.payload_sha256,approved:c.approved,reason:c.reason}),c.organization_id);
   if(v.approval_id!==c.approval_id||v.payload_sha256!==c.payload_sha256)throw new Error("decision reference mismatch");return reply(v);
  }
  const v=funding(await protectedPost<unknown>(session,`${base}/orders/${encodeURIComponent(c.order_id)}/funding`,{organization_id:c.organization_id,expected_order_version:c.expected_order_version},c.request_key),c.organization_id);
  if(v.order_id!==c.order_id||v.request_key!==c.request_key)throw new Error("funding reference mismatch");return reply(v,201);
 }catch(e){return failure(e)}
}
export async function GET(request:Request){
 if(["cross-site","none"].includes(request.headers.get("sec-fetch-site")??""))return reply({code:"CROSS_ORIGIN_REJECTED"},403);
 const session=await readSession();if(!session)return reply({code:"UNAUTHENTICATED"},401);if(!allowed(session,"stored_value:read"))return reply({code:"FORBIDDEN"},403);
 const search=new URL(request.url).searchParams;if(request.url.length>2048||[...search.keys()].some(k=>search.getAll(k).length!==1))return reply({code:"INVALID_QUERY"},400);
 const parsed=storedQuery.safeParse(Object.fromEntries(search));if(!parsed.success)return reply({code:"INVALID_QUERY"},400);const q=parsed.data;
 if(!session.organizations.includes(q.organization_id))return reply({code:"FORBIDDEN"},403);
 try{
  if(q.kind==="order"){
   const v=storedOrder.parse(await protectedGet<unknown>(session,`${base}/orders/${encodeURIComponent(q.order_id)}`,{organization_id:q.organization_id,after:q.after}));
   if(v.organization_id!==q.organization_id||v.order_id!==q.order_id)throw new Error("order scope mismatch");return reply(v);
  }
  if(q.kind==="funding-result"){
   const v=funding(await protectedGet<unknown>(session,`${base}/funding-result`,{organization_id:q.organization_id},q.request_key),q.organization_id);if(v.request_key!==q.request_key)throw new Error("funding reference mismatch");return reply(v);
  }
  const id=q.kind==="approval"?q.approval_id:"svapproval_"+digest(q.operation_id);
  const v=approval(await protectedGet<unknown>(session,`${base}/approvals/${encodeURIComponent(id)}`,{organization_id:q.organization_id}),q.organization_id);if(v.approval_id!==id)throw new Error("approval reference mismatch");return reply(v);
 }catch(e){return failure(e)}
}
