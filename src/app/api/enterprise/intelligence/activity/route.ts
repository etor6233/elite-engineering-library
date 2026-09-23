import { NextResponse,type NextRequest } from "next/server";
import { allowed,readSession } from "@/platform/auth/session";
import { protectedGet } from "@/platform/backend/protected-client";
import { loadBusinessConfig } from "@/platform/config/load";
import { aiActivity } from "@/platform/intelligence/activity";
const problem=(code:string,status:number)=>NextResponse.json({code},{status,headers:{"cache-control":"no-store"}});
export async function GET(request:NextRequest){
 if((await loadBusinessConfig()).features.ai_activity!==true)return problem("NOT_FOUND",404);
 const session=await readSession();if(!session)return problem("UNAUTHENTICATED",401);
 if(!allowed(session,"whatsapp:approve"))return problem("FORBIDDEN",403);
 if(new URL(request.url).search!="")return problem("INVALID_QUERY",400);
 try{const data=aiActivity.parse(await protectedGet(session,"/v1/franchise/whatsapp/ai-activity",{}));
  if(!session.organizations.includes(data.organization_id))return problem("FORBIDDEN",403);
  return NextResponse.json(data,{headers:{"cache-control":"no-store"}})
 }catch{return problem("CONSULTATION_UNAVAILABLE",503)}
}
