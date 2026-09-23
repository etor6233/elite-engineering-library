import { NextResponse,type NextRequest } from "next/server";
import { allowed,readSession } from "@/platform/auth/session";
import { applicationBaseUrl } from "@/platform/auth/oidc-client";
import { protectedGet,protectedPost } from "@/platform/backend/protected-client";
import { loadBusinessConfig } from "@/platform/config/load";
import { attemptSchema,trainingCommand } from "@/platform/training/contract";
import { boundedCommandBody } from "@/platform/backend/bounded-command";
const response=(code:string,status:number)=>NextResponse.json({code},{status,headers:{"cache-control":"no-store"}});
async function enabled(){return (await loadBusinessConfig()).features.training_portal===true}
export async function GET(request:NextRequest){
 if(!await enabled())return response("NOT_FOUND",404);const session=await readSession();if(!session)return response("UNAUTHENTICATED",401);
 if(!allowed(session,"training:learn"))return response("FORBIDDEN",403);
 const query=new URL(request.url).searchParams;if([...query.keys()].some(x=>x!=="attempt_id")||query.getAll("attempt_id").length!==1)return response("INVALID_QUERY",400);
 const id=query.get("attempt_id")!;if(!/^[0-9a-f]{8}-[0-9a-f]{4}-[1-8][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/.test(id))return response("INVALID_QUERY",400);
 try{const value=attemptSchema.parse(await protectedGet(session,`/v1/training/attempts/${id}`,{}));if(!session.organizations.includes(value.organization_id)||value.learner_subject!==session.subject)return response("NOT_FOUND",404);return NextResponse.json(value,{headers:{"cache-control":"no-store"}})}catch{return response("CONSULTATION_UNAVAILABLE",409)}
}
export async function POST(request:NextRequest){
 if(!await enabled())return response("NOT_FOUND",404);let origin:string;try{origin=applicationBaseUrl().origin}catch{return response("UNAVAILABLE",503)};
 if(request.headers.get("origin")!==origin)return response("CROSS_ORIGIN_REJECTED",403);
 if(request.headers.get("content-type")!=="application/json")return response("UNSUPPORTED_MEDIA_TYPE",415);
 const session=await readSession();if(!session)return response("UNAUTHENTICATED",401);
 if(!allowed(session,"training:learn")&&!allowed(session,"training:review"))return response("FORBIDDEN",403);
 let raw:string;try{raw=await boundedCommandBody(request,65536)}catch(error){return response(error instanceof Error&&error.message==="BODY_TOO_LARGE"?"COMMAND_TOO_LARGE":"INVALID_COMMAND",error instanceof Error&&error.message==="BODY_TOO_LARGE"?413:400)}
 let body:unknown;try{body=JSON.parse(raw)}catch{return response("INVALID_COMMAND",400)};
 const parsed=trainingCommand.safeParse(body);if(!parsed.success)return response("INVALID_COMMAND",400);
 const c=parsed.data;
 if(!allowed(session,c.action==="assess"?"training:review":"training:learn"))return response("FORBIDDEN",403);
 let path:string,payload:unknown;
 if(c.action==="start"){path="/v1/training/attempts";payload={attempt_id:c.attempt_id,course_id:c.course_id,profile_sha256:c.profile_sha256}}
 else if(c.action==="acknowledge"){path=`/v1/training/attempts/${c.attempt_id}/acknowledgements`;payload={lesson_id:c.lesson_id,profile_sha256:c.profile_sha256}}
 else if(c.action==="submit"){path=`/v1/training/attempts/${c.attempt_id}/submissions`;payload={answers:c.answers,profile_sha256:c.profile_sha256}}
 else{path=`/v1/training/assessments/${encodeURIComponent(c.request_id)}/decision`;payload={payload_sha256:c.payload_sha256,approved:c.approved,reason:c.reason}}
 try{return NextResponse.json(await protectedPost(session,path,payload),{headers:{"cache-control":"no-store"}})}catch{return response("CONSULT_RECORDED_STATE",409)}
}
