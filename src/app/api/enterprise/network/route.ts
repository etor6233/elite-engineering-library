import {NextResponse,type NextRequest} from "next/server";
import {allowed,readSession,type PortalSession} from "@/platform/auth/session";
import {applicationBaseUrl} from "@/platform/auth/oidc-client";
import {protectedGet,protectedPost} from "@/platform/backend/protected-client";
import {boundedCommandBody} from "@/platform/backend/bounded-command";
import {loadBusinessConfig} from "@/platform/config/load";
import {networkID,networkCommand,networkReceipt,networkEntity,networkPermission,networkPayload,networkReference,networkMatches} from "@/platform/network/contract";
const fail=(code:string,status:number)=>NextResponse.json({code},{status,headers:{"cache-control":"no-store"}});
const ok=(value:unknown)=>NextResponse.json(value,{headers:{"cache-control":"no-store"}});
const scoped=(s:PortalSession,org:string)=>allowed(s,"*")||org!==""&&s.organizations.includes(org);
const anyPermission=(s:PortalSession)=>allowed(s,"network:admin")||allowed(s,"franchise:write");
export async function GET(r:NextRequest){
 if((await loadBusinessConfig()).features.network_portal!==true)return fail("NOT_FOUND",404);
 const s=await readSession();if(!s)return fail("UNAUTHENTICATED",401);if(!anyPermission(s))return fail("FORBIDDEN",403);
 const q=r.nextUrl.searchParams,kind=q.get("kind"),id=q.get("id"),org=q.get("scope_organization_id");
 if(q.size!==3||["kind","id","scope_organization_id"].some(k=>q.getAll(k).length!==1)||!["result","organization","agreement"].includes(kind??"")||!networkID.safeParse(id).success||org===null||org!==""&&!networkID.safeParse(org).success||!scoped(s,org))return fail("INVALID_SCOPE",403);
 try{
  if(kind==="result"){const v=networkReceipt.parse(await protectedGet(s,`/v1/franchise/network/commands/${encodeURIComponent(id!)}`,{scope_organization_id:org}));if(v.command_id!==id||v.scope_organization_id!==org||v.actor!==s.subject||!allowed(s,networkPermission(v.action)))return fail("SCOPE_MISMATCH",409);return ok(v)}
  if(!allowed(s,networkPermission("transition-"+kind)))return fail("FORBIDDEN",403);
  const v=networkEntity.parse(await protectedGet(s,`/v1/franchise/network/entities/${kind}/${encodeURIComponent(id!)}`,{scope_organization_id:org}));if(v.id!==id||v.kind!==kind||v.organization_id!==org)return fail("SCOPE_MISMATCH",409);return ok(v);
 }catch{return fail("CONSULTATION_UNAVAILABLE",409)}
}
export async function POST(r:NextRequest){
 if((await loadBusinessConfig()).features.network_portal!==true)return fail("NOT_FOUND",404);
 try{if(r.headers.get("origin")!==applicationBaseUrl().origin)return fail("CROSS_ORIGIN_REJECTED",403)}catch{return fail("UNAVAILABLE",503)}
 const s=await readSession();if(!s)return fail("UNAUTHENTICATED",401);if(!anyPermission(s))return fail("FORBIDDEN",403);
 if(r.headers.get("content-type")!=="application/json")return fail("UNSUPPORTED_MEDIA_TYPE",415);if(r.nextUrl.searchParams.size)return fail("INVALID_COMMAND",400);
 try{
  let raw:unknown;try{raw=JSON.parse(await boundedCommandBody(r,32768))}catch(e){if(e instanceof Error&&e.message==="BODY_TOO_LARGE")throw e;return fail("INVALID_COMMAND",400)}
  const parsed=networkCommand.safeParse(raw);if(!parsed.success)return fail("INVALID_COMMAND",400);const c=parsed.data;
  if(!allowed(s,networkPermission(c.action))||!scoped(s,c.scope_organization_id))return fail("FORBIDDEN",403);
  const expected=await networkReference(c),v=networkReceipt.parse(await protectedPost(s,"/v1/franchise/network/commands",networkPayload(c)));
  if(!networkMatches(v,expected,s.subject))return fail("UNCONFIRMED",409);return ok(v);
 }catch(e){return fail(e instanceof Error&&e.message==="BODY_TOO_LARGE"?"COMMAND_TOO_LARGE":"UNCONFIRMED",e instanceof Error&&e.message==="BODY_TOO_LARGE"?413:409)}
}
