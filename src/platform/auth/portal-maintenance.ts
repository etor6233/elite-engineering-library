import "server-only";
import {randomBytes} from "node:crypto";
import {z} from "zod";
import {readBackendResponse} from "@/platform/backend/public-client";
import type {PortalProfile} from "./portal-profile";
import {revokePortalSweepItem} from "./portal-lifecycle";
const item=z.object({session_id_sha256:z.string().regex(/^[a-f0-9]{64}$/),version:z.number().int().positive(),state:z.literal("revoked"),operation_id:z.string().max(43),ciphertext:z.string().max(32768),access_expires_at:z.string().datetime({offset:true}),absolute_expires_at:z.string().datetime({offset:true}),revocation_pending:z.literal(true)}).strict();
const result=z.object({items:z.array(item).max(2),purged:z.number().int().min(0).max(100),unconfirmed_purged:z.number().int().min(0).max(100)}).strict();
export async function maintainPortalSessions(p:PortalProfile,serviceBearer:string){
 if(!/^Bearer [^\s]{1,16384}$/.test(serviceBearer))throw new Error("PORTAL_SERVICE_REQUIRED");
 // Go verifies this inbound service bearer before exposing/claiming any record.
 const response=await fetch(new URL("/v1/private/portal-sessions/sweep",p.bridge_url),{method:"POST",redirect:"error",cache:"no-store",signal:AbortSignal.timeout(5000),headers:{authorization:serviceBearer,"content-type":"application/json","X-Portal-Profile-SHA256":p.documentSHA256},body:JSON.stringify({operation_id:randomBytes(32).toString("base64url")})});
 const batch=result.parse(await readBackendResponse<unknown>(response));const completed=await Promise.all(batch.items.map(row=>revokePortalSweepItem(p,row,row.session_id_sha256)));
 return {claimed:batch.items.length,confirmed:completed.filter(Boolean).length,pending:completed.filter(v=>!v).length,purged:batch.purged,unconfirmed_purged:batch.unconfirmed_purged};
}
