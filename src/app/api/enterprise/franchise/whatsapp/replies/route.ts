import { boundedCommandBody } from "@/platform/backend/bounded-command";
import { NextResponse,type NextRequest } from "next/server";
import { allowed,readSession } from "@/platform/auth/session";
import { applicationBaseUrl } from "@/platform/auth/oidc-client";
import { protectedPost } from "@/platform/backend/protected-client";
import { replyCommand } from "@/platform/notifications/reply-contract";
export async function POST(request:NextRequest){
 const response=(code:string,status:number)=>NextResponse.json({code},{status,headers:{"cache-control":"no-store"}});
 let expected:string;try{expected=applicationBaseUrl().origin}catch{return response("UNAVAILABLE",503)}
 if(request.headers.get("origin")!==expected)return response("CROSS_ORIGIN_REJECTED",403);
 if(request.headers.get("content-type")!=="application/json")return response("UNSUPPORTED_MEDIA_TYPE",415);
 let text:string;try{text=await boundedCommandBody(request,4096)}catch(error){return response(error instanceof Error&&error.message==="BODY_TOO_LARGE"?"COMMAND_TOO_LARGE":"INVALID_BODY",error instanceof Error&&error.message==="BODY_TOO_LARGE"?413:400)}
 let input:unknown;try{input=JSON.parse(text)}catch{return response("INVALID_COMMAND",400)}
 const parsed=replyCommand.safeParse(input);if(!parsed.success)return response("INVALID_COMMAND",400);
 const session=await readSession();if(!session)return response("UNAUTHENTICATED",401);
 const command=parsed.data;if(!allowed(session,"whatsapp:approve")||(command.action!=="decision"&&!allowed(session,"whatsapp:send")))return response("FORBIDDEN",403);
 const body=command.action==="decision"?{payload_sha256:command.payload_sha256,approved:command.approved,reason:command.reason}:{payload_sha256:command.payload_sha256};
 try{const result=await protectedPost<unknown>(session,`/v1/franchise/whatsapp/replies/${encodeURIComponent(command.request_id)}/${command.action}`,body);return NextResponse.json(result,{headers:{"cache-control":"no-store"}})}catch{return response("CONSULT_STATE_BEFORE_RETRY",409)}
}
