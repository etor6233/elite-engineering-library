import {NextResponse,type NextRequest} from "next/server";
import {readSession} from "@/platform/auth/session";
import {applicationBaseUrl} from "@/platform/auth/oidc-client";
import {boundedCommandBody} from "@/platform/backend/bounded-command";
import {protectedGet,protectedPost} from "@/platform/backend/protected-client";
import {BackendProblem} from "@/platform/backend/public-client";
import {campaignActionAllowed,campaignBatch,campaignCommand,campaignID,campaignPage,campaignPermissions,preparedCampaign,campaignStatus,campaignTemplates} from "@/platform/campaigns/contract";
const root="/v1/franchise/marketing/campaigns";
const reply=(value:unknown,status=200)=>NextResponse.json(value,{status,headers:{"cache-control":"no-store"}});
export async function GET(request:NextRequest){
 if(process.env.CAMPAIGN_WORKSPACE_ENABLED!=="true")return reply({code:"NOT_INSTALLED"},404);
 const session=await readSession();if(!session)return reply({code:"UNAUTHENTICATED"},401);
 const p=campaignPermissions(session.permissions);if(!p.read)return reply({code:"FORBIDDEN"},403);
 const q=request.nextUrl.searchParams;if([...q.keys()].some(k=>!["kind","after","id"].includes(k)||q.getAll(k).length!==1))return reply({code:"INVALID_QUERY"},400);
 const kind=q.get("kind")??"campaigns",id=q.get("id"),after=q.get("after");
 if(!["campaigns","audience","templates","status"].includes(kind)||kind==="status"&&!campaignID.safeParse(id).success||after!==null&&!/^[A-Za-z0-9][A-Za-z0-9_.:-]{0,127}$/.test(after)||id!==null&&kind!=="status"||after!==null&&!["campaigns","audience"].includes(kind))return reply({code:"INVALID_QUERY"},400);
 if((kind==="audience"||kind==="templates")&&!p.prepare)return reply({code:"FORBIDDEN"},403);
 try{
  if(kind==="status")return reply(campaignStatus.parse(await protectedGet(session,`${root}/workspace/result/${encodeURIComponent(id!)}`,{})));
  if(kind==="templates")return reply(campaignTemplates.parse(await protectedGet(session,root+"/workspace/templates",{})));
  return reply(campaignPage.parse(await protectedGet(session,root+"/workspace",{kind,after:after??undefined})));
 }catch(e){if(e instanceof BackendProblem&&e.status===404)return reply({code:"NOT_FOUND"},404);return reply({code:"READ_UNAVAILABLE"},503)}
}
export async function POST(request:NextRequest){
 if(process.env.CAMPAIGN_WORKSPACE_ENABLED!=="true")return reply({code:"NOT_INSTALLED"},404);
 let origin:string;try{origin=applicationBaseUrl().origin}catch{return reply({code:"UNAVAILABLE"},503)}
 if(request.headers.get("origin")!==origin)return reply({code:"CROSS_ORIGIN_REJECTED"},403);
 if(request.headers.get("content-type")!=="application/json")return reply({code:"UNSUPPORTED_MEDIA_TYPE"},415);
 const session=await readSession();if(!session)return reply({code:"UNAUTHENTICATED"},401);
 let body:unknown;try{body=JSON.parse(await boundedCommandBody(request,32768))}catch(e){return reply({code:e instanceof Error&&e.message==="BODY_TOO_LARGE"?"BODY_TOO_LARGE":"INVALID_BODY"},e instanceof Error&&e.message==="BODY_TOO_LARGE"?413:400)}
 const parsed=campaignCommand.safeParse(body);if(!parsed.success)return reply({code:"INVALID_COMMAND"},400);const c=parsed.data;
 if(!campaignActionAllowed(c.action,session.permissions))return reply({code:"FORBIDDEN"},403);
 try{
  if(c.action==="prepare")return reply(preparedCampaign.parse(await protectedPost(session,root+"/workspace/prepare",c.selection)));
  if(c.action==="create"){
   return reply(campaignBatch.parse(await protectedPost(session,root+"/workspace/create",{selection:c.selection,campaign_sha256:c.campaign_sha256})));
  }
  const payload=c.action==="decision"?{campaign_sha256:c.campaign_sha256,approve:c.approve,reason:c.reason}:c.action==="stop"?{campaign_sha256:c.campaign_sha256,reason:c.reason}:{campaign_sha256:c.campaign_sha256};
  const value=await protectedPost(session,`${root}/${encodeURIComponent(c.campaign_id)}/${c.action}`,payload);
  if(c.action==="stop"){if(!value||typeof value!=="object"||!("stopped"in value)||value.stopped!==true)throw Error("invalid stop response");return reply({stopped:true})}
  return reply(campaignBatch.parse(value));
 }catch{return reply({code:c.action==="prepare"?"SOURCE_NOT_CURRENT":"CONSULT_STATE_BEFORE_RETRY"},409)}
}
