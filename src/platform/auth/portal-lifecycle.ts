import "server-only";
import { hkdfSync,randomBytes,createHash } from "node:crypto";
import { EncryptJWT,jwtDecrypt,type JWTPayload } from "jose";
import { z } from "zod";
import { portalSecret,type PortalProfile } from "./portal-profile";
import { portalConfiguration,portalOIDC } from "./portal-protocol";
import { portalSessionBridge,type PortalStoredSession } from "./portal-session-bridge";
import type { PortalSession,AuthFlow } from "./session";

const atom=z.string().min(1).max(200);const values=z.array(atom).max(100).refine(a=>new Set(a).size===a.length);
const tokenSchema=z.object({subject:atom,tenantId:atom,permissions:values,organizations:values.min(1),accessToken:z.string().min(1).max(16384),refreshToken:z.string().min(1).max(8192),idToken:z.string().min(1).max(16384),absoluteExpires:z.number().int().positive()});
type StoredTokens=z.infer<typeof tokenSchema>;
const flows=new Map<string,Promise<PortalSession|null>>();
function key(p:PortalProfile,purpose:string){const secret=portalSecret("OIDC_PORTAL_SESSION_KEY_FILE");if(secret.length<32)throw new Error("OIDC_SESSION_KEY_REJECTED");return new Uint8Array(hkdfSync("sha256",secret,p.documentSHA256,purpose,32));}
async function seal(p:PortalProfile,purpose:string,payload:JWTPayload,expires:number){let encrypted=new EncryptJWT(payload).setProtectedHeader({alg:"dir",enc:"A256GCM",typ:"JWT"}).setIssuer(p.profile_id).setAudience(purpose).setIssuedAt();if(purpose!=="portal-token-vault-v1")encrypted=encrypted.setExpirationTime(expires);return encrypted.encrypt(key(p,purpose));}
async function open(p:PortalProfile,purpose:string,value:string){return (await jwtDecrypt(value,key(p,purpose),{issuer:p.profile_id,audience:purpose,keyManagementAlgorithms:["dir"],contentEncryptionAlgorithms:["A256GCM"]})).payload;}
function safeSession(t:StoredTokens):PortalSession{return{subject:t.subject,tenantId:t.tenantId,permissions:t.permissions,organizations:t.organizations,accessToken:t.accessToken}}
function checkedTokens(p:PortalProfile,result:Awaited<ReturnType<typeof portalOIDC.authorizationCodeGrant>>,absolute:number,previous?:StoredTokens):{tokens:StoredTokens;expiry:string}{
 const claims=result.claims();if(!claims||!result.id_token||!result.refresh_token||!result.expires_in||result.expires_in<1||result.expires_in>3600||claims.iss!==p.issuer||claims.aud!==p.client_id||claims.tenant_id!==p.tenant_id||typeof claims.exp!=="number"||claims.exp*1000<=Date.now()||typeof claims.iat!=="number"||claims.iat*1000>Date.now()+30000)throw new Error("OIDC_IDENTITY_REJECTED");
 const tokens=tokenSchema.parse({subject:claims.sub,tenantId:claims.tenant_id,permissions:claims.permissions,organizations:claims.organization_ids,accessToken:result.access_token,refreshToken:result.refresh_token,idToken:result.id_token,absoluteExpires:absolute});
 if(tokens.organizations.some(o=>!p.organization_ids.includes(o))||tokens.permissions.some(permission=>!p.allowed_permissions.includes(permission))||(previous&&(tokens.subject!==previous.subject||tokens.tenantId!==previous.tenantId||tokens.refreshToken===previous.refreshToken)))throw new Error("OIDC_IDENTITY_REJECTED");
 const expiry=Math.min(Date.now()+result.expires_in*1000,claims.exp*1000,absolute*1000);if(expiry<=Date.now())throw new Error("OIDC_SESSION_EXPIRED");return{tokens,expiry:new Date(expiry).toISOString()};
}
async function encryptTokens(p:PortalProfile,sid:string,t:StoredTokens){return seal(p,"portal-token-vault-v1",{...t,sessionID:sid,profile:p.documentSHA256},t.absoluteExpires)}
async function decryptTokens(p:PortalProfile,sid:string,row:PortalStoredSession){const value=await open(p,"portal-token-vault-v1",row.ciphertext);if(value.sessionID!==sid||value.profile!==p.documentSHA256)throw new Error("OIDC_SESSION_BINDING_REJECTED");return tokenSchema.parse(value)}
export async function createPortalSession(p:PortalProfile,flow:AuthFlow,currentURL:URL):Promise<{cookie:string;maxAge:number}>{
 const c=await portalConfiguration(p);const result=await portalOIDC.authorizationCodeGrant(c,currentURL,{pkceCodeVerifier:flow.codeVerifier,expectedState:flow.state,expectedNonce:flow.nonce});
 const absolute=Math.floor(Date.now()/1000)+p.maximum_session_seconds;const checked=checkedTokens(p,result,absolute);const sid=randomBytes(32).toString("base64url");const ciphertext=await encryptTokens(p,sid,checked.tokens);
 const command={session_id:sid,ciphertext,access_expires_at:checked.expiry,absolute_expires_at:new Date(absolute*1000).toISOString()};
 try{await portalSessionBridge(p,"create",command)}catch{const recovered=await portalSessionBridge(p,"read",{session_id:sid});if(recovered.state!=="active"||recovered.version!==1||recovered.ciphertext!==ciphertext)throw new Error("OIDC_SESSION_STORE_UNAVAILABLE")}
 return{cookie:await seal(p,"portal-browser-handle-v1",{sessionID:sid,profile:p.documentSHA256},absolute),maxAge:p.maximum_session_seconds};
}
async function sessionID(p:PortalProfile,cookie:string){const value=await open(p,"portal-browser-handle-v1",cookie);if(typeof value.sessionID!=="string"||!/^[A-Za-z0-9_-]{43}$/.test(value.sessionID)||value.profile!==p.documentSHA256)throw new Error("OIDC_SESSION_REJECTED");return value.sessionID}
async function revokeStored(p:PortalProfile,sid:string,row:PortalStoredSession):Promise<boolean>{
 if(!row.revocation_pending)return true;
  try{const tokens=await decryptTokens(p,sid,row);const c=await portalConfiguration(p);await portalOIDC.tokenRevocation(c,tokens.refreshToken,{token_type_hint:"refresh_token"});await portalOIDC.tokenRevocation(c,tokens.accessToken,{token_type_hint:"access_token"});await portalSessionBridge(p,"ack-revocation",{session_id:sid,version:row.version,operation_id:row.operation_id});return true}catch{return false}
}
async function resolve(p:PortalProfile,sid:string):Promise<PortalSession|null>{
 let row=await portalSessionBridge(p,"read",{session_id:sid});if(row.state!=="active"||Date.parse(row.absolute_expires_at)<=Date.now())return null;
 let tokens=await decryptTokens(p,sid,row);if(tokens.absoluteExpires*1000<=Date.now())return null;if(Date.parse(row.access_expires_at)>Date.now()+p.refresh_before_seconds*1000)return safeSession(tokens);
 const operation=randomBytes(32).toString("base64url"),version=row.version;
 try{row=await portalSessionBridge(p,"claim",{session_id:sid,version,operation_id:operation})}catch{
  row=await portalSessionBridge(p,"read",{session_id:sid});
  if(row.state==="active"&&row.version>version)return safeSession(await decryptTokens(p,sid,row));
  if(row.state!=="refreshing"||row.version!==version||row.operation_id!==operation)return null;
 }
 if(row.state!=="refreshing"||row.operation_id!==operation||row.version!==version)return null;
 try{
  const result=await portalOIDC.refreshTokenGrant(await portalConfiguration(p),tokens.refreshToken);
  const checked=checkedTokens(p,result,tokens.absoluteExpires,tokens);tokens=checked.tokens;const ciphertext=await encryptTokens(p,sid,tokens);
  try{row=await portalSessionBridge(p,"commit",{session_id:sid,version,operation_id:operation,ciphertext,access_expires_at:checked.expiry})}catch{
   row=await portalSessionBridge(p,"read",{session_id:sid});
   if(row.version!==version+1||row.operation_id!==operation||row.ciphertext!==ciphertext)throw new Error("OIDC_REFRESH_COMMIT_UNKNOWN");
  }
  if(row.state!=="active"){await revokeStored(p,sid,row);return null}
  return safeSession(tokens);
 }catch{
  try{await portalSessionBridge(p,"abort",{session_id:sid,version,operation_id:operation})}catch{/* Durable refreshing state already rejects a second grant. */}
  return null;
 }
}
export async function readPortalSession(p:PortalProfile,cookie:string):Promise<PortalSession|null>{
 try{const sid=await sessionID(p,cookie);const key=p.documentSHA256+":"+sid;let pending=flows.get(key);if(!pending){pending=resolve(p,sid).catch(()=>null);flows.set(key,pending);void pending.finally(()=>flows.delete(key))}return await pending}catch{return null}
}
export async function logoutPortalSession(p:PortalProfile,cookie:string):Promise<{revoked:boolean;providerRevocationComplete:boolean;redirect:string}>{
 const sid=await sessionID(p,cookie);let row:PortalStoredSession;
 try{row=await portalSessionBridge(p,"revoke",{session_id:sid})}catch{row=await portalSessionBridge(p,"read",{session_id:sid});if(row.state!=="revoked")throw new Error("OIDC_LOGOUT_STORE_UNAVAILABLE")}
 const done=await revokeStored(p,sid,row);const c=await portalConfiguration(p);const redirect=portalOIDC.buildEndSessionUrl(c,{client_id:p.client_id,post_logout_redirect_uri:p.post_logout_url}).toString();return{revoked:true,providerRevocationComplete:done,redirect};
}
// A worker may open only the vault whose authenticated SID hashes to the leased row.
export async function revokePortalSweepItem(p:PortalProfile,row:PortalStoredSession,expectedSessionHash:string):Promise<boolean>{
 try{const value=await open(p,"portal-token-vault-v1",row.ciphertext);if(typeof value.sessionID!=="string"||!/^[A-Za-z0-9_-]{43}$/.test(value.sessionID)||value.profile!==p.documentSHA256||createHash("sha256").update(value.sessionID).digest("hex")!==expectedSessionHash)return false;return revokeStored(p,value.sessionID,row)}catch{return false}
}
