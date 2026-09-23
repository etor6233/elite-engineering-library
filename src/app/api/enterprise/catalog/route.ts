import {NextResponse,type NextRequest} from "next/server";
import {readSession,allowed} from "@/platform/auth/session";
import {applicationBaseUrl} from "@/platform/auth/oidc-client";
import {protectedGet,protectedPost,protectedPostPNG} from "@/platform/backend/protected-client";
import {BackendProblem} from "@/platform/backend/public-client";
import {boundedCommandBody,boundedRequestBytes} from "@/platform/backend/bounded-command";
import {loadBusinessConfig} from "@/platform/config/load";
import {catalogCommand,catalogDraft,catalogReceipt,catalogID,commandPayload} from "@/platform/catalog/authoring";
import {publishedCatalogSchema} from "@/platform/catalog/schema";
const response=(code:string,status:number)=>NextResponse.json({code},{status,headers:{"cache-control":"no-store"}});
const result=(v:unknown)=>NextResponse.json(v,{headers:{"cache-control":"no-store"}});
const permits=["catalog:draft","catalog:publish","catalog:review:legal","catalog:review:technical","catalog:review:media","catalog:review:publication"];
async function enabled(){return(await loadBusinessConfig()).features.catalog_editor===true}
export async function GET(r:NextRequest){
 if(!await enabled())return response("NOT_FOUND",404);
 const session=await readSession();if(!session)return response("UNAUTHENTICATED",401);
 if(!allowed(session,"catalog:read"))return response("FORBIDDEN",403);
 const q=r.nextUrl.searchParams,kind=q.get("kind"),org=q.get("organization_id"),id=q.get("id");
 const keys=kind==="current"?["kind","organization_id"]:["kind","organization_id","id"];
 if([...q.keys()].length!==keys.length||keys.some(k=>q.getAll(k).length!==1)||![...q.keys()].every(k=>keys.includes(k))||!catalogID.safeParse(org).success||!session.organizations.includes(org!))return response("INVALID_SCOPE",403);
 if(!["current","draft","command"].includes(kind??"")||kind!=="current"&&!catalogID.safeParse(id).success)return response("INVALID_REFERENCE",400);
 try{
  if(kind==="current"){
   try{return result({publication:publishedCatalogSchema.parse(await protectedGet(session,"/v1/admin/catalog/current",{organization_id:org!}))})}
   catch(e){if(e instanceof BackendProblem&&e.status===404)return result({publication:null});throw e}
  }
  const raw=await protectedGet(session,`/v1/admin/catalog/${kind==="draft"?"drafts":"commands"}/${encodeURIComponent(id!)}`,{organization_id:org!});
  if(kind==="draft"){const value=catalogDraft.parse(raw);if(value.id!==id||value.snapshot.profile.tenant_id!==session.tenantId||value.snapshot.profile.organization_id!==org)return response("SCOPE_MISMATCH",409);return result(value)}
  const value=catalogReceipt.parse(raw);if(value.command_id!==id||value.actor!==session.subject)return response("SCOPE_MISMATCH",409);return result(value);
 }catch{return response("CONSULTATION_UNAVAILABLE",409)}
}
export async function POST(r:NextRequest){
 if(!await enabled())return response("NOT_FOUND",404);
 try{if(r.headers.get("origin")!==applicationBaseUrl().origin)return response("CROSS_ORIGIN_REJECTED",403)}catch{return response("UNAVAILABLE",503)}
 const session=await readSession();if(!session)return response("UNAUTHENTICATED",401);
 if(!allowed(session,"catalog:read")||!permits.some(p=>allowed(session,p)))return response("FORBIDDEN",403);
 try{
  const type=r.headers.get("content-type");let value:unknown,commandID:string;
  if(type==="image/png"){
   if(!allowed(session,"catalog:draft"))return response("FORBIDDEN",403);
   const q=r.nextUrl.searchParams,org=q.get("organization_id"),id=q.get("command_id");
   if([...q.keys()].length!==2||q.getAll("organization_id").length!==1||q.getAll("command_id").length!==1||!catalogID.safeParse(org).success||!catalogID.safeParse(id).success||!session.organizations.includes(org!))return response("INVALID_SCOPE",403);
   commandID=id!;const bytes=await boundedRequestBytes(r,1048576);
   value=await protectedPostPNG(session,`/v1/admin/catalog/media/${encodeURIComponent(id!)}?${new URLSearchParams({organization_id:org!})}`,bytes);
  }else{
   if(type!=="application/json")return response("UNSUPPORTED_MEDIA_TYPE",415);
   if(r.nextUrl.searchParams.size!==0)return response("INVALID_COMMAND",400);
   let body:unknown;try{body=JSON.parse(await boundedCommandBody(r,32768))}catch(e){if(e instanceof Error&&e.message==="BODY_TOO_LARGE")throw e;return response("INVALID_COMMAND",400)}
   const parsed=catalogCommand.safeParse(body);if(!parsed.success)return response("INVALID_COMMAND",400);
   const c=parsed.data;commandID=c.command_id;
   if(!session.organizations.includes(c.organization_id))return response("INVALID_SCOPE",403);
   const permission=c.action==="review"?`catalog:review:${c.stage}`:c.action==="publish"?"catalog:publish":"catalog:draft";
   if(!allowed(session,permission))return response("FORBIDDEN",403);
   const tail=c.action==="source"?"sources":c.action==="draft"?"drafts":`drafts/${encodeURIComponent(c.draft_id)}/${c.action==="review"?"review/"+c.stage:"publish"}`;
   value=await protectedPost(session,`/v1/admin/catalog/${tail}?${new URLSearchParams({organization_id:c.organization_id})}`,commandPayload(c));
  }
  const v=catalogReceipt.parse(value);
  if(v.command_id!==commandID||v.actor!==session.subject)return response("UNCONFIRMED",409);
  return result(v);
 }catch(e){return response(e instanceof Error&&e.message==="BODY_TOO_LARGE"?"COMMAND_TOO_LARGE":"UNCONFIRMED",e instanceof Error&&e.message==="BODY_TOO_LARGE"?413:409)}
}
