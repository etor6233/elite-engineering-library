import "server-only";
import { z } from "zod";
import { readBackendResponse } from "@/platform/backend/public-client";
import { portalServiceAccessToken } from "./portal-protocol";
import type { PortalProfile } from "./portal-profile";
const record=z.object({version:z.number().int().positive(),state:z.enum(["active","refreshing","reauth_required","revoked"]),operation_id:z.string().max(43),ciphertext:z.string().max(32768),access_expires_at:z.string().datetime({offset:true}),absolute_expires_at:z.string().datetime({offset:true}),revocation_pending:z.boolean()}).strict();
export type PortalStoredSession=z.infer<typeof record>;
export type PortalSessionAction="create"|"read"|"claim"|"commit"|"abort"|"revoke"|"ack-revocation";
export type PortalSessionCommand={session_id:string;operation_id?:string;version?:number;ciphertext?:string;access_expires_at?:string;absolute_expires_at?:string};
export async function portalSessionBridge(p:PortalProfile,action:PortalSessionAction,command:PortalSessionCommand):Promise<PortalStoredSession>{
 if(!/^[a-zA-Z0-9_-]{43}$/.test(command.session_id))throw new Error("PORTAL_SESSION_REJECTED");
 const body=JSON.stringify(command);if(Buffer.byteLength(body)>40000)throw new Error("PORTAL_SESSION_REJECTED");
 const response=await fetch(new URL("/v1/private/portal-sessions/"+action,p.bridge_url),{method:"POST",cache:"no-store",redirect:"error",signal:AbortSignal.timeout(5000),headers:{authorization:`Bearer ${await portalServiceAccessToken(p)}`,"content-type":"application/json","X-Portal-Profile-SHA256":p.documentSHA256},body});
 return record.parse(await readBackendResponse<unknown>(response));
}
