// AUTHORED BFF composition over fixed, existing owners and the scoped read module.
import {createHash} from "node:crypto";
import {allowed,type PortalSession} from "@/platform/auth/session";
import {boundedCommandBody} from "@/platform/backend/bounded-command";
import {datasets,idSchema,pageSchema,parseCommand} from "./contract";
type Transport=(session:PortalSession,path:string,method:string,body?:string)=>Promise<Response>;
type Deps={session:()=>Promise<PortalSession|null>;origin:()=>string;transport:Transport};
const json=(body:unknown,status=200)=>Response.json(body,{status,headers:{"cache-control":"no-store","x-content-type-options":"nosniff"}});
const no=(code:string,status:number)=>json({code,effect:"NOT_ATTEMPTED"},status);
async function boundedResponse(response:Response){
 if(!response.body)throw new Error("NO_BODY");const reader=response.body.getReader();let n=0;const chunks:Uint8Array[]=[];let timer:ReturnType<typeof setTimeout>|undefined;
 const deadline=new Promise<never>((_,reject)=>{timer=setTimeout(()=>reject(new Error("BODY_DEADLINE")),5000)});
 try{for(;;){const part=await Promise.race([reader.read(),deadline]);if(part.done)break;n+=part.value.byteLength;if(n>2097152)throw new Error("BODY_BUDGET");chunks.push(part.value)}const bytes=new Uint8Array(n);let at=0;for(const c of chunks){bytes.set(c,at);at+=c.length}return JSON.parse(new TextDecoder("utf-8",{fatal:true}).decode(bytes)) as unknown}finally{clearTimeout(timer);void reader.cancel().catch(()=>{});reader.releaseLock()}
}
export function createWarehouseGateway(d:Deps){return async(request:Request)=>{
 const session=await d.session();if(!session)return no("UNAUTHENTICATED",401);if(!allowed(session,"inventory:read"))return no("FORBIDDEN",403);
 const url=new URL(request.url);const authorized=(org:string)=>allowed(session,"*")||session.organizations.includes(org);
 if(request.method==="GET"&&!url.search)return json({schema:"warehouse-access/v1",scope:createHash("sha256").update(JSON.stringify([session.tenantId,session.subject,session.organizations])).digest("hex"),organizations:session.organizations,canWrite:allowed(session,"inventory:write")});
 let path:string,body:string|undefined,read=true;
 try{
  if(request.method==="POST"){
   if(!allowed(session,"inventory:write"))return no("FORBIDDEN",403);if(request.headers.get("origin")!==d.origin())return no("CROSS_ORIGIN_REJECTED",403);
   if(url.search||request.headers.get("content-type")!=="application/json")return no("INVALID_COMMAND",400);
   const c=parseCommand(JSON.parse(await boundedCommandBody(request,16384)));if(!authorized(c.organization))return no("ORGANIZATION_FORBIDDEN",403);
   for(const k of ["from_organization_id","to_organization_id"])if(c.payload[k]!==undefined&&!authorized(String(c.payload[k])))return no("ORGANIZATION_FORBIDDEN",403);
   path=c.actionDefinition.path.replace("{id}",encodeURIComponent(c.id??""));body=JSON.stringify(c.payload);read=false;
  }else if(request.method==="GET"){
   const allowedKeys=["dataset","organization_id","search","item_id","parent_id","id","request_id","cursor","limit","reservation_request"];
   for(const key of url.searchParams.keys())if(!allowedKeys.includes(key)||url.searchParams.getAll(key).length!==1)return no("INVALID_QUERY",400);
   const org=idSchema.parse(url.searchParams.get("organization_id"));if(!authorized(org))return no("ORGANIZATION_FORBIDDEN",403);
   if(url.searchParams.has("reservation_request")){
    if([...url.searchParams.keys()].some(k=>!["organization_id","reservation_request"].includes(k)))return no("INVALID_QUERY",400);
    path="/v1/inventory/warehouse-workspace/reservation-requests/"+encodeURIComponent(idSchema.parse(url.searchParams.get("reservation_request")))+"?organization_id="+encodeURIComponent(org);
   }else{
    const dataset=url.searchParams.get("dataset");if(!datasets.includes(dataset as typeof datasets[number]))return no("INVALID_QUERY",400);
    const params=new URLSearchParams({organization_id:org});
    for(const key of allowedKeys){if(["dataset","organization_id","reservation_request"].includes(key))continue;const v=url.searchParams.get(key);if(v===null)continue;if(v.length>(key==="cursor"?4096:key==="search"?100:128)||/[\x00-\x1f\x7f]/u.test(v))return no("INVALID_QUERY",400);if(key==="limit"&&!/^(?:[1-9]|[1-4][0-9]|50)$/u.test(v))return no("INVALID_QUERY",400);params.set(key,v)}
    path="/v1/inventory/workspace/"+dataset+"?"+params;
   }
  }else return no("METHOD_NOT_ALLOWED",405);
 }catch{return no("INVALID_REQUEST",400)}
 try{
  const response=await d.transport(session,path,read?"GET":"POST",body);const value=await boundedResponse(response);
  if(!response.ok)return json({code:response.status===409?"STATE_CHANGED":"WAREHOUSE_REQUEST_FAILED",effect:read?"READ_ONLY":response.status===409?"REJECTED":"UNKNOWN"},response.status>=500?503:response.status);
  if(read&&path.startsWith("/v1/inventory/workspace/")){const page=pageSchema.parse(value),expected=new URL(path,'http://warehouse.invalid');if(page.organization_id!==expected.searchParams.get('organization_id')||page.dataset!==expected.pathname.split('/').at(-1))throw new Error('RESPONSE_SCOPE_MISMATCH');return json(page)}
  if(value===null||typeof value!=='object'||Array.isArray(value))throw new Error('INVALID_OWNER_RESPONSE');
  if(path.includes('/warehouse-workspace/')){const v=value as Record<string,unknown>;if(v.schema!=='warehouse-reservation-receipt/v1'||typeof v.request_id!=='string'||typeof v.payload_sha256!=='string'||!/^[a-f0-9]{64}$/u.test(v.payload_sha256)||!v.reservation||typeof v.reservation!=='object')throw new Error('INVALID_RESERVATION_RECEIPT')}
  // Mutation shape belongs to its exact backend owner; its successful response is not a cloud/prod qualification.
  return json({schema:"warehouse-result/v1",result:value,effect:read?"READ_ONLY":"CONFIRMED"});
 }catch{return json({code:read?"WAREHOUSE_UNAVAILABLE":"RESULT_UNCERTAIN",effect:read?"READ_ONLY":"UNKNOWN"},503)}
}}
