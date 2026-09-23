import {NextResponse} from "next/server";
import {portalLifecycleEnabled,portalProfile} from "@/platform/auth/portal-profile";
import {maintainPortalSessions} from "@/platform/auth/portal-maintenance";
import {BackendProblem} from "@/platform/backend/public-client";
export const maxDuration=45;
export async function POST(request:Request){
 if(!portalLifecycleEnabled())return NextResponse.json({code:"OIDC_LIFECYCLE_DISABLED"},{status:404});
 if(request.body){const reader=request.body.getReader();try{const chunk=await reader.read();if(!chunk.done){void reader.cancel().catch(()=>{});return NextResponse.json({code:"EMPTY_BODY_REQUIRED"},{status:400})}}finally{reader.releaseLock()}}
 const bearer=request.headers.get("authorization")??"";if(!/^Bearer [^\s]{1,16384}$/.test(bearer))return NextResponse.json({code:"SERVICE_IDENTITY_REQUIRED"},{status:401});
 try{const result=await maintainPortalSessions(portalProfile(),bearer);return NextResponse.json(result,{status:result.pending?503:200,headers:{"cache-control":"no-store"}})}catch(error){return NextResponse.json({code:"SESSION_MAINTENANCE_UNAVAILABLE"},{status:error instanceof BackendProblem&&[401,403].includes(error.status)?error.status:503,headers:{"cache-control":"no-store"}})}
}
