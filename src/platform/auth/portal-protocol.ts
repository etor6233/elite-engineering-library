import "server-only";
import * as oidc from "openid-client";
import { createHash } from "node:crypto";
import { portalSecret, type PortalProfile } from "./portal-profile";

const configurations=new Map<string,Promise<oidc.Configuration>>();
function boundedFetch(p:PortalProfile): oidc.CustomFetch {
 const endpoints=new Set([p.issuer.replace(/\/$/,"")+"/.well-known/openid-configuration",p.token_endpoint,p.jwks_endpoint,p.revocation_endpoint]);
 return async(input,init)=>{
  const url=String(input);if(!endpoints.has(url))throw new Error("OIDC_ENDPOINT_REJECTED");
  const body=init.body instanceof Uint8Array?new Uint8Array(init.body):init.body;
  const response=await fetch(input,{method:init.method,headers:init.headers,...(body!==undefined?{body}:{}),...(init.duplex?{duplex:init.duplex}:{}),redirect:"error",signal:AbortSignal.any([AbortSignal.timeout(10000),...(init.signal?[init.signal]:[])])});
  if(!response.body)return response;
  const reader=response.body.getReader();let received=0;const chunks:Uint8Array[]=[];let complete=false;
  try{while(true){const chunk=await reader.read();if(chunk.done){complete=true;break}received+=chunk.value.byteLength;if(received>262144)throw new Error("OIDC_RESPONSE_TOO_LARGE");chunks.push(chunk.value)};const bytes=new Uint8Array(received);let position=0;for(const chunk of chunks){bytes.set(chunk,position);position+=chunk.byteLength};return new Response(bytes,{status:response.status,statusText:response.statusText,headers:response.headers});}
  finally{if(!complete)void reader.cancel().catch(()=>{});reader.releaseLock()}
 };
}
export async function portalConfiguration(p:PortalProfile,service=false):Promise<oidc.Configuration>{
 const secret=portalSecret(service?"OIDC_PORTAL_SERVICE_CLIENT_SECRET_FILE":"OIDC_PORTAL_CLIENT_SECRET_FILE");
 const cacheKey=p.documentSHA256+":"+service+":"+createHash("sha256").update(secret).digest("hex");
 let config=configurations.get(cacheKey);
 if(!config){if(configurations.size>=2)configurations.clear();config=(async()=>{
  const c=await oidc.discovery(new URL(p.issuer),service?p.service_client_id:p.client_id,{id_token_signed_response_alg:"RS256"},oidc.ClientSecretBasic(secret),{[oidc.customFetch]:boundedFetch(p),timeout:10,execute:[...(p.transport==="LOOPBACK_FIXTURE"?[oidc.allowInsecureRequests]:[]),oidc.enableNonRepudiationChecks]});
  const m=c.serverMetadata();
  if(m.issuer!==p.issuer||m.token_endpoint!==p.token_endpoint||m.jwks_uri!==p.jwks_endpoint||m.authorization_endpoint!==p.authorization_endpoint||m.revocation_endpoint!==p.revocation_endpoint||m.end_session_endpoint!==p.end_session_endpoint||!m.id_token_signing_alg_values_supported?.includes("RS256")||(m.token_endpoint_auth_methods_supported&&!m.token_endpoint_auth_methods_supported.includes("client_secret_basic")))throw new Error("OIDC_METADATA_REJECTED");
  return c;
 })();configurations.set(cacheKey,config);void config.catch(()=>configurations.delete(cacheKey));}
 return config;
}
let serviceToken:{key:string;value:string;expires:number}|undefined;
const serviceFlights=new Map<string,Promise<string>>();
export async function portalServiceAccessToken(p:PortalProfile):Promise<string>{
 if(serviceToken?.key===p.documentSHA256&&serviceToken.expires>Date.now()+10000)return serviceToken.value;
 const existing=serviceFlights.get(p.documentSHA256);if(existing)return existing;
 const flight=(async()=>{serviceToken=undefined;try{const config=await portalConfiguration(p,true);const result=await oidc.clientCredentialsGrant(config,{scope:p.service_scopes.join(" "),...(p.service_audience_parameter?{audience:p.service_audience_parameter}:{})});if(!result.access_token||result.access_token.length>16384||result.token_type.toLowerCase()!=="bearer"||!result.expires_in||result.expires_in<20||result.expires_in>3600)throw new Error("OIDC_SERVICE_TOKEN_REJECTED");serviceToken={key:p.documentSHA256,value:result.access_token,expires:Date.now()+result.expires_in*1000};return result.access_token;}catch{throw new Error("OIDC_SERVICE_UNAVAILABLE")}finally{serviceFlights.delete(p.documentSHA256)}})();serviceFlights.set(p.documentSHA256,flight);return flight;
}
export {oidc as portalOIDC};
