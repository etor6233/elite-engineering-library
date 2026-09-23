import "server-only";
import {applicationBaseUrl} from "@/platform/auth/oidc-client";
import {backendFetch} from "@/platform/backend/cloud-run-transport";
import {readSession} from "@/platform/auth/session";
import {createWarehouseGateway} from "./gateway";
export const warehouseGateway=createWarehouseGateway({session:readSession,origin:()=>applicationBaseUrl().origin,transport:async(session,path,method,body)=>{
 const base=new URL(process.env.ENTERPRISE_API_BASE_URL??"");if(!["https:","http:"].includes(base.protocol)||base.username||base.password||base.search||base.hash||!path.startsWith("/v1/inventory/")||/[\r\n]/u.test(path))throw new Error("Invalid warehouse transport");
 return backendFetch(new URL(path,base),{method,headers:{accept:"application/json",authorization:`Bearer ${session.accessToken}`,...(body?{"content-type":"application/json"}:{})},...(body?{body}:{}),cache:"no-store",redirect:"error",signal:AbortSignal.timeout(10000)});
}});
