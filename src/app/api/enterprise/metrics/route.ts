import{NextResponse,type NextRequest}from"next/server";
import{readSession,allowed}from"@/platform/auth/session";
import{loadBusinessConfig}from"@/platform/config/load";
import{protectedGet}from"@/platform/backend/protected-client";
import{metricID,metricFamilies,parseMetric}from"@/platform/metrics/contract";
const response=(value:unknown,status=200)=>NextResponse.json(value,{status,headers:{"cache-control":"private, no-store",vary:"Cookie"}});
export async function GET(r:NextRequest){
 if((await loadBusinessConfig()).features.role_workspace!==true)return response({code:"NOT_FOUND"},404);
 const s=await readSession();if(!s)return response({code:"UNAUTHENTICATED"},401);
 const q=r.nextUrl.searchParams,kind=q.get("kind"),org=q.get("organization_id"),id=q.get("id")??"",family=metricFamilies.find(f=>f.id===kind);
 if(r.nextUrl.search.length>600||q.size!==(kind==="stored-value"||kind==="survey"?3:2)||["kind","organization_id"].some(k=>q.getAll(k).length!==1)||q.getAll("id").length>1||!family||!metricID.safeParse(org).success||(kind==="stored-value"||kind==="survey")&&!metricID.safeParse(id).success||Array.from(q.keys()).some(k=>!["kind","organization_id","id"].includes(k)))return response({code:"INVALID_QUERY"},400);
 if(!allowed(s,family.permission)||!s.organizations.includes(org!))return response({code:"FORBIDDEN"},403);
 try{
  const path=kind==="stored-value"?"/v1/franchise/stored-value/metrics":kind==="survey"?"/v1/reporting/surveys/"+encodeURIComponent(id):"/v1/reporting/operations/"+family.id;
  const value=await protectedGet(s,path,{organization_id:org!,...(kind==="stored-value"?{program_id:id}:{})});
  return response(parseMetric(family.id,value,org!,id));
 }catch{return response({code:"METRIC_UNAVAILABLE"},503)}
}
