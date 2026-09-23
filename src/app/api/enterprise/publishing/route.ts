import { NextResponse,type NextRequest } from "next/server";
import { z } from "zod";
import { allowed,readSession,type PortalSession } from "@/platform/auth/session";
import { applicationBaseUrl } from "@/platform/auth/oidc-client";
import { loadBusinessConfig } from "@/platform/config/load";
import { boundedCommandBody } from "@/platform/backend/bounded-command";
import { backendFetch } from "@/platform/backend/cloud-run-transport";
import { readBackendResponse } from "@/platform/backend/public-client";
import { publicationQuery,publicationsPage,publicationDetail,publishingCommand,preparedPublication } from "@/platform/publishing/contract";
const reply=(data:unknown,status=200)=>NextResponse.json(data,{status,headers:{"cache-control":"no-store"}});
const problem=(code:string,status:number)=>reply({code},status);
function socialBase(){const raw=process.env.SOCIAL_API_BASE_URL;if(!raw||!URL.canParse(raw))throw new Error("SOCIAL_API_BASE_URL is invalid");const value=new URL(raw);if(!["http:","https:"].includes(value.protocol)||value.username||value.password||value.search||value.hash)throw new Error("social backend URL is unsafe");return value}
async function socialGet<T>(session:PortalSession,path:string,query:Record<string,string|undefined>):Promise<T>{if(!path.startsWith("/v1/social/")||path.includes(".."))throw new Error("social path invalid");const target=new URL(path,socialBase());for(const [name,value] of Object.entries(query))if(value)target.searchParams.set(name,value);const response=await backendFetch(target,{cache:"no-store",redirect:"error",signal:AbortSignal.timeout(10000),headers:{accept:"application/json",authorization:`Bearer ${session.accessToken}`}});return readBackendResponse<T>(response)}
async function socialPost<T>(session:PortalSession,path:string,command:unknown):Promise<T>{if(!path.startsWith("/v1/social/")||path.includes(".."))throw new Error("social path invalid");const response=await backendFetch(new URL(path,socialBase()),{method:"POST",cache:"no-store",redirect:"error",signal:AbortSignal.timeout(10000),headers:{accept:"application/json","content-type":"application/json",authorization:`Bearer ${session.accessToken}`},body:JSON.stringify(command)});return readBackendResponse<T>(response)}
export async function GET(request:NextRequest){
 if((await loadBusinessConfig()).features.social_publishing!==true)return problem("NOT_FOUND",404);
 const session=await readSession();if(!session)return problem("UNAUTHENTICATED",401);if(!allowed(session,"social:read"))return problem("FORBIDDEN",403);
 const params=new URL(request.url).searchParams;const keys=[...params.keys()];if(new Set(keys).size!==keys.length)return problem("INVALID_QUERY",400);
 const q=publicationQuery.safeParse(Object.fromEntries(params));if(!q.success)return problem("INVALID_QUERY",400);
 try{
  const data="id"in q.data?publicationDetail.parse(await socialGet(session,`/v1/social/workspace/${encodeURIComponent(q.data.id)}`,{})):publicationsPage.parse(await socialGet(session,"/v1/social/workspace",q.data));
  if(!session.organizations.includes(data.organization_id)||("status"in data&&data.status.payload.intent.tenant_id!==session.tenantId))return problem("FORBIDDEN",403);
  return reply(data);
 }catch{return problem("CONSULTATION_UNAVAILABLE",503)}
}
export async function POST(request:NextRequest){
 if((await loadBusinessConfig()).features.social_publishing!==true)return problem("NOT_FOUND",404);
 try{if(request.headers.get("origin")!==applicationBaseUrl().origin)return problem("CROSS_ORIGIN_REJECTED",403)}catch{return problem("UNAVAILABLE",503)}
 if(new URL(request.url).search)return problem("INVALID_QUERY",400);
 if(request.headers.get("content-type")!=="application/json")return problem("UNSUPPORTED_MEDIA_TYPE",415);
 const session=await readSession();if(!session)return problem("UNAUTHENTICATED",401);if(!allowed(session,"social:read"))return problem("FORBIDDEN",403);
 let command:z.infer<typeof publishingCommand>;
 try{command=publishingCommand.parse(JSON.parse(await boundedCommandBody(request,32768)))}catch{return problem("INVALID_COMMAND",400)}
 const permission=command.action==="decide"?"social:approve":command.action==="reconcile"?"social:reconcile":"social:request";
 if(!allowed(session,permission))return problem("FORBIDDEN",403);
 try{
  if(command.action==="prepare"){
   const result=preparedPublication.parse(await socialPost(session,"/v1/social/workspace/prepare",command.draft));if(result.request.intent.tenant_id!==session.tenantId)return problem("FORBIDDEN",403);return reply(result);
  }
  if(command.action==="submit"){
   if(command.prepared.request.intent.tenant_id!==session.tenantId)return problem("FORBIDDEN",403);
   const result=z.object({approval_id:z.string(),request_sha256:z.string(),replayed:z.boolean()}).strict().parse(await socialPost(session,"/v1/social/requests",command.prepared.request));
   if(result.approval_id!==command.prepared.request.intent.approval_id||result.request_sha256!==command.prepared.request_sha256)return problem("CONSULT_STATE_BEFORE_RETRY",409);return reply(result);
  }
  if(command.action==="decide"){
   const result=z.object({state:z.enum(["approved","rejected"])}).strict().parse(await socialPost(session,`/v1/social/requests/${encodeURIComponent(command.approval_id)}/decision`,{request_sha256:command.request_sha256,approve:command.approve,reason:command.reason}));return reply(result);
  }
  const result=z.object({state:z.literal("reconciled")}).strict().parse(await socialPost(session,`/v1/social/requests/${encodeURIComponent(command.approval_id)}/reconcile`,{}));return reply(result);
 }catch{return problem("CONSULT_STATE_BEFORE_RETRY",409)}
}
