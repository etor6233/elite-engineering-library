import{NextResponse,type NextRequest}from"next/server";
import{allowed,readSession,type PortalSession}from"@/platform/auth/session";
import{applicationBaseUrl}from"@/platform/auth/oidc-client";
import{protectedGet,protectedPost}from"@/platform/backend/protected-client";
import{boundedCommandBody}from"@/platform/backend/bounded-command";
import{loadBusinessConfig}from"@/platform/config/load";
import{helpID,helpCommand,helpReceipt,helpArticle,helpPage,helpPermission,helpReference,helpMatches}from"@/platform/help/cms-contract";
const fail=(code:string,status:number)=>NextResponse.json({code},{status,headers:{"cache-control":"private, no-store"}});
const ok=(value:unknown)=>NextResponse.json(value,{headers:{"cache-control":"private, no-store"}});
const scoped=(s:PortalSession,org:string)=>allowed(s,"*")||s.organizations.includes(org);
const editor=(s:PortalSession)=>allowed(s,"help:write")||allowed(s,"help:publish");
export async function GET(r:NextRequest){
 if((await loadBusinessConfig()).features.help_cms!==true)return fail("NOT_FOUND",404);
 const s=await readSession();if(!s)return fail("UNAUTHENTICATED",401);if(!editor(s)&&!allowed(s,"help:read"))return fail("FORBIDDEN",403);
 const q=r.nextUrl.searchParams,kind=q.get("kind"),org=q.get("organization_id"),id=q.get("id");
 const keys=kind==="list"?["kind","organization_id","locale","q","after"]:kind==="article"?["kind","organization_id","id","version"]:kind==="history"?["kind","organization_id","id","before"]:["kind","organization_id","id"];
 if(r.url.length>2048||!["list","article","history","result"].includes(kind??"")||[...q.keys()].some(k=>!keys.includes(k)||q.getAll(k).length!==1)||!helpID.safeParse(org).success||!scoped(s,org!))return fail("INVALID_SCOPE",403);
 if(kind!=="list"&&!helpID.safeParse(id).success)return fail("INVALID_QUERY",400);
 if(["history","result"].includes(kind!)&&!editor(s))return fail("FORBIDDEN",403);
 const query:Record<string,string>={organization_id:org!};for(const k of keys.filter(k=>!["kind","organization_id","id"].includes(k))){if(q.has(k))query[k]=q.get(k)!}
 try{
  if(kind==="result"){
   const v=helpReceipt.parse(await protectedGet(s,`/v1/help/cms/commands/${encodeURIComponent(id!)}`,query));
   if(v.command_id!==id||v.article.organization_id!==org||v.actor!==s.subject||!allowed(s,helpPermission(v.action)))return fail("SCOPE_MISMATCH",409);return ok(v)
  }
  if(kind==="list"||kind==="history"){
   const path=kind==="list"?"/v1/help/cms/articles":`/v1/help/cms/articles/${encodeURIComponent(id!)}/history`;
   const v=helpPage.parse(await protectedGet(s,path,query));if(kind==="history"&&v.items.some(a=>a.id!==id)||kind==="list"&&v.items.some(a=>a.locale!==query.locale)||!editor(s)&&v.items.some(a=>a.state!=="published"))return fail("SCOPE_MISMATCH",409);return ok(v)
  }
  const v=helpArticle.parse(await protectedGet(s,`/v1/help/cms/articles/${encodeURIComponent(id!)}`,query));if(v.id!==id||v.organization_id!==org||query.version&&v.version!==query.version||!editor(s)&&v.state!=="published")return fail("SCOPE_MISMATCH",409);return ok(v)
 }catch{return fail("CONSULTATION_UNAVAILABLE",409)}
}
export async function POST(r:NextRequest){
 if((await loadBusinessConfig()).features.help_cms!==true)return fail("NOT_FOUND",404);
 try{if(r.headers.get("origin")!==applicationBaseUrl().origin)return fail("CROSS_ORIGIN_REJECTED",403)}catch{return fail("UNAVAILABLE",503)}
 const s=await readSession();if(!s)return fail("UNAUTHENTICATED",401);if(!editor(s))return fail("FORBIDDEN",403);
 if(r.headers.get("content-type")!=="application/json")return fail("UNSUPPORTED_MEDIA_TYPE",415);if(r.nextUrl.searchParams.size)return fail("INVALID_COMMAND",400);
 try{
  let raw:unknown;try{raw=JSON.parse(await boundedCommandBody(r,32768))}catch(e){if(e instanceof Error&&e.message==="BODY_TOO_LARGE")throw e;return fail("INVALID_COMMAND",400)}
  const parsed=helpCommand.safeParse(raw);if(!parsed.success)return fail("INVALID_COMMAND",400);const c=parsed.data;if(!allowed(s,helpPermission(c.action))||!scoped(s,c.organization_id))return fail("FORBIDDEN",403);
  const p=await helpReference(c),v=helpReceipt.parse(await protectedPost(s,"/v1/help/cms/commands",c));if(!await helpMatches(v,p,s.subject))return fail("UNCONFIRMED",409);return ok(v)
 }catch(e){return fail(e instanceof Error&&e.message==="BODY_TOO_LARGE"?"COMMAND_TOO_LARGE":"UNCONFIRMED",e instanceof Error&&e.message==="BODY_TOO_LARGE"?413:409)}
}
