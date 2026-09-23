import {NextResponse,type NextRequest} from "next/server";
import {allowed,readSession} from "@/platform/auth/session";
import {applicationBaseUrl} from "@/platform/auth/oidc-client";
import {protectedGet,protectedPost} from "@/platform/backend/protected-client";
import {boundedCommandBody} from "@/platform/backend/bounded-command";
import {loadBusinessConfig} from "@/platform/config/load";
import {supplyID,supplyCommand,supplyReceipt,supplyPlan,supplyPermission,supplySurface,supplyPayload,supplyReference,supplyMatches} from "@/platform/supply/contract";
const fail=(code:string,status:number)=>NextResponse.json({code},{status,headers:{"cache-control":"no-store"}});
const ok=(value:unknown)=>NextResponse.json(value,{headers:{"cache-control":"no-store"}});
const permissions=["supply:plan","supply:factory","supply:receive","supply:release","supply:inspect"];
export async function GET(r:NextRequest){
 if((await loadBusinessConfig()).features.supply_portal!==true)return fail("NOT_FOUND",404);
 const s=await readSession();if(!s)return fail("UNAUTHENTICATED",401);
 const q=r.nextUrl.searchParams,org=q.get("organization_id"),id=q.get("purchase_order_id"),surface=q.get("surface"),command=q.get("command_id"),after=q.get("after_unit");
 const keys=["organization_id","purchase_order_id","surface",...(command!==null?["command_id"]:[]),...(after!==null?["after_unit"]:[])];
 if(command!==null&&after!==null||q.size!==keys.length||keys.some(k=>q.getAll(k).length!==1)||!supplyID.safeParse(org).success||!s.organizations.includes(org!)||!supplyID.safeParse(id).success||!["factory","franchise"].includes(surface??"")||command!==null&&!supplyID.safeParse(command).success||after!==null&&!supplyID.safeParse(after).success)return fail("INVALID_SCOPE",403);
 if(!allowed(s,surface==="factory"?"supply:factory-read":"supply:read"))return fail("FORBIDDEN",403);
 try{
  const raw=await protectedGet(s,`/v1/${surface}/supply/orders/${encodeURIComponent(id!)}`,{organization_id:org!,command_id:command??undefined,after_unit:after??undefined});
  if(command!==null){const v=supplyReceipt.parse(raw);if(v.command_id!==command||v.purchase_order_id!==id||v.actor!==s.subject)return fail("SCOPE_MISMATCH",409);return ok(v)}
  const v=supplyPlan.parse(raw);if(v.purchase_order_id!==id||(surface==="factory"?v.factory_organization_id:v.destination_organization_id)!==org)return fail("SCOPE_MISMATCH",409);return ok(v);
 }catch{return fail("CONSULTATION_UNAVAILABLE",409)}
}
export async function POST(r:NextRequest){
 if((await loadBusinessConfig()).features.supply_portal!==true)return fail("NOT_FOUND",404);
 try{if(r.headers.get("origin")!==applicationBaseUrl().origin)return fail("CROSS_ORIGIN_REJECTED",403)}catch{return fail("UNAVAILABLE",503)}
 const s=await readSession();if(!s)return fail("UNAUTHENTICATED",401);
 if(!permissions.some(p=>allowed(s,p)))return fail("FORBIDDEN",403);
 if(r.headers.get("content-type")!=="application/json")return fail("UNSUPPORTED_MEDIA_TYPE",415);
 if(r.nextUrl.searchParams.size)return fail("INVALID_COMMAND",400);
 try{
  let raw:unknown;try{raw=JSON.parse(await boundedCommandBody(r,32768))}catch(e){if(e instanceof Error&&e.message==="BODY_TOO_LARGE")throw e;return fail("INVALID_COMMAND",400)}
  const parsed=supplyCommand.safeParse(raw);if(!parsed.success)return fail("INVALID_COMMAND",400);const c=parsed.data,surface=supplySurface(c.kind);
  if(!s.organizations.includes(c.organization_id)||!allowed(s,supplyPermission(c.kind))||!allowed(s,surface==="factory"?"supply:factory-read":"supply:read"))return fail("FORBIDDEN",403);
  const expected=await supplyReference(c);
  const v=supplyReceipt.parse(await protectedPost(s,`/v1/${surface}/supply/orders/${encodeURIComponent(c.purchase_order_id)}/${c.kind}?${new URLSearchParams({organization_id:c.organization_id})}`,supplyPayload(c)));
  if(!supplyMatches(v,expected,s.subject))return fail("UNCONFIRMED",409);return ok(v);
 }catch(e){return fail(e instanceof Error&&e.message==="BODY_TOO_LARGE"?"COMMAND_TOO_LARGE":"UNCONFIRMED",e instanceof Error&&e.message==="BODY_TOO_LARGE"?413:409)}
}
