import {NextResponse,type NextRequest} from "next/server";
import {readSession,allowed} from "@/platform/auth/session";
import {applicationBaseUrl} from "@/platform/auth/oidc-client";
import {protectedGet,protectedPost} from "@/platform/backend/protected-client";
import {boundedCommandBody} from "@/platform/backend/bounded-command";
import {loadBusinessConfig} from "@/platform/config/load";
import {warrantyID,warrantySurface,warrantyCommand,warrantyPermission,warrantyReadPermission,warrantyPath,warrantyPayload,warrantyOffer,warrantyActivation,warrantyClaim,warrantyStep,warrantyStepReference,warrantyStepMatches,warrantyProfileView,warrantyQuote} from "@/platform/warranty/contract";
const fail=(code:string,status:number)=>NextResponse.json({code},{status,headers:{"cache-control":"no-store"}});
const ok=(value:unknown)=>NextResponse.json(value,{headers:{"cache-control":"no-store"}});
const writePermissions=["warranty:self","warranty:offer","warranty:activate","warranty:request","warranty:diagnose","warranty:plan","warranty:approve","warranty:work","warranty:quality","warranty:reconcile","warranty:cancel"];
export async function GET(r:NextRequest){
 if((await loadBusinessConfig()).features.warranty_portal!==true)return fail("NOT_FOUND",404);
 const s=await readSession();if(!s)return fail("UNAUTHENTICATED",401);
 const q=r.nextUrl.searchParams,kind=q.get("kind"),org=q.get("organization_id"),id=q.get("id"),command=q.get("command_id"),surface=warrantySurface.safeParse(q.get("surface"));
 const keys=["kind","organization_id","surface",...(kind!=="profile"?["id"]:[]),...(kind==="command"?["command_id"]:[])];
 if(!surface.success||!["profile","quote","offer","activation","claim","command"].includes(kind??"")||q.size!==keys.length||keys.some(k=>q.getAll(k).length!==1)||!warrantyKeys(q,keys)||!warrantyID.safeParse(org).success||!s.organizations.includes(org!)||kind!=="profile"&&!warrantyID.safeParse(id).success||kind==="command"&&!warrantyID.safeParse(command).success)return fail("INVALID_SCOPE",403);
 const role=surface.data;if(kind==="quote"&&(role!=="franchise"||!allowed(s,"warranty:offer")))return fail("FORBIDDEN",403);if(!allowed(s,warrantyReadPermission(role))||kind==="profile"&&role!=="franchise"||["offer","activation"].includes(kind!)&&role==="factory")return fail("FORBIDDEN",403);
 try{
  const tail=kind==="quote"?`quotes/${encodeURIComponent(id!)}/source`:kind==="profile"?"profile":kind==="offer"?`quotes/${encodeURIComponent(id!)}`:kind==="activation"?`handovers/${encodeURIComponent(id!)}/activation`:`claims/${encodeURIComponent(id!)}`;
  const raw=await protectedGet(s,`/v1/${role}/warranty/${tail}`,{organization_id:org!,command_id:command??undefined});
  if(kind==="quote"){const v=warrantyQuote.parse(raw);if(v.quote_id!==id||v.organization_id!==org)return fail("SCOPE_MISMATCH",409);return ok(v)}
  if(kind==="profile"){const v=warrantyProfileView.parse(raw);if(v.profile.tenant_id!==s.tenantId||v.profile.organization_id!==org)return fail("SCOPE_MISMATCH",409);return ok(v)}
  if(kind==="offer"){const v=warrantyOffer.parse(raw);if(v.quote_id!==id||v.profile.tenant_id!==s.tenantId||v.profile.organization_id!==org||role==="customer"&&v.customer_subject!==s.subject)return fail("SCOPE_MISMATCH",409);return ok(v)}
  if(kind==="activation"){const v=warrantyActivation.parse(raw);if(v.handover_id!==id||v.organization_id!==org||role==="customer"&&v.customer_subject!==s.subject)return fail("SCOPE_MISMATCH",409);return ok(v)}
  if(kind==="command"){const v=warrantyStep.parse(raw);if(v.case_id!==id||v.command_id!==command||v.actor!==s.subject)return fail("SCOPE_MISMATCH",409);return ok(v)}
  const v=warrantyClaim.parse(raw);if(v.case_id!==id||(role==="factory"?v.factory_organization_id:v.organization_id)!==org||role==="customer"&&v.customer_subject!==s.subject)return fail("SCOPE_MISMATCH",409);return ok(v);
 }catch{return fail("CONSULTATION_UNAVAILABLE",409)}
}
function warrantyKeys(q:URLSearchParams,keys:string[]){return [...q.keys()].every(k=>keys.includes(k))}
export async function POST(r:NextRequest){
 if((await loadBusinessConfig()).features.warranty_portal!==true)return fail("NOT_FOUND",404);
 try{if(r.headers.get("origin")!==applicationBaseUrl().origin)return fail("CROSS_ORIGIN_REJECTED",403)}catch{return fail("UNAVAILABLE",503)}
 const s=await readSession();if(!s)return fail("UNAUTHENTICATED",401);if(!writePermissions.some(p=>allowed(s,p)))return fail("FORBIDDEN",403);
 if(r.headers.get("content-type")!=="application/json")return fail("UNSUPPORTED_MEDIA_TYPE",415);if(r.nextUrl.searchParams.size)return fail("INVALID_COMMAND",400);
 try{
  let raw:unknown;try{raw=JSON.parse(await boundedCommandBody(r,32768))}catch(e){if(e instanceof Error&&e.message==="BODY_TOO_LARGE")throw e;return fail("INVALID_COMMAND",400)}
  const parsed=warrantyCommand.safeParse(raw);if(!parsed.success)return fail("INVALID_COMMAND",400);const c=parsed.data,permission=warrantyPermission(c);
  if(!permission||!s.organizations.includes(c.organization_id)||!allowed(s,permission)||!allowed(s,warrantyReadPermission(c.surface)))return fail("FORBIDDEN",403);
  const result=await protectedPost(s,warrantyPath(c)+"?"+new URLSearchParams({organization_id:c.organization_id}),warrantyPayload(c));
  if(c.action==="offer"||c.action==="acknowledge"){
   const v=warrantyOffer.parse(result);if(v.quote_id!==c.quote_id||v.profile_sha256!==c.profile_sha256||v.profile.tenant_id!==s.tenantId||v.profile.organization_id!==c.organization_id||v.quote_version!==(c.action==="offer"?(BigInt(c.quote_version)+1n).toString():c.quote_version))return fail("UNCONFIRMED",409);
   if(c.action==="offer"?v.offered_by!==s.subject:v.customer_subject!==s.subject||!v.acknowledged||v.evidence_sha256!==c.evidence_sha256)return fail("UNCONFIRMED",409);return ok(v);
  }
  if(c.action==="activate"){const v=warrantyActivation.parse(result);if(v.handover_id!==c.handover_id||v.organization_id!==c.organization_id||v.activated_by!==s.subject)return fail("UNCONFIRMED",409);return ok(v)}
  const v=warrantyStep.parse(result);if(!warrantyStepMatches(v,await warrantyStepReference(c),s.subject))return fail("UNCONFIRMED",409);return ok(v);
 }catch(e){return fail(e instanceof Error&&e.message==="BODY_TOO_LARGE"?"COMMAND_TOO_LARGE":"UNCONFIRMED",e instanceof Error&&e.message==="BODY_TOO_LARGE"?413:409)}
}
