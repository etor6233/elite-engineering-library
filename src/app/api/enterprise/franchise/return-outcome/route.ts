import {NextRequest,NextResponse} from "next/server";
import {z} from "zod";
import {readSession,allowed} from "@/platform/auth/session";
import {protectedGet} from "@/platform/backend/protected-client";
import {BackendProblem} from "@/platform/backend/public-client";
import {returnOutcome} from "@/platform/returns/outcome";
const query=z.object({organizationId:z.string().min(1).max(128),authorizationId:z.string().min(1).max(128)}).strict();
const reply=(body:unknown,status=200)=>NextResponse.json(body,{status,headers:{"cache-control":"no-store"}});
export async function GET(request:NextRequest){
 const session=await readSession();if(!session)return reply({code:"UNAUTHENTICATED"},401);
 const raw=new URL(request.url).searchParams;for(const key of raw.keys())if(raw.getAll(key).length!==1)return reply({code:"INVALID_QUERY"},400);
 const parsed=query.safeParse(Object.fromEntries(raw));if(!parsed.success)return reply({code:"INVALID_QUERY"},400);
 const q=parsed.data;if(!session.organizations.includes(q.organizationId)||!allowed(session,"handover:manage"))return reply({code:"FORBIDDEN"},403);
 try{
  const value=returnOutcome.safeParse(await protectedGet<unknown>(session,`/v1/franchise/returns/${encodeURIComponent(q.authorizationId)}/outcome`,{organization_id:q.organizationId}));
  if(!value.success||value.data.organization_id!==q.organizationId||value.data.authorization_id!==q.authorizationId)return reply({code:"RETURN_OUTCOME_UNCONFIRMED"},502);
  return reply(value.data);
 }catch(error){if(error instanceof BackendProblem)return reply({code:error.code},error.status);return reply({code:"RETURN_OUTCOME_UNAVAILABLE"},503)}
}
