import "server-only";
import {applicationBaseUrl} from "@/platform/auth/oidc-client";
import {backendFetch} from "@/platform/backend/cloud-run-transport";
import type {CaptureTransport} from "./gateway";
export function captureConfig(){if(process.env.CAPTURE_ENABLED!=="true")throw new Error("Capture not enabled");return {origin:applicationBaseUrl().origin}}
export const captureTransport:CaptureTransport=async(session,action,body)=>{const base=new URL(process.env.ENTERPRISE_API_BASE_URL??"");if(!["http:","https:"].includes(base.protocol)||base.username||base.password||base.search||base.hash||!["decode","resolve"].includes(action))throw new Error("Invalid capture endpoint");return backendFetch(new URL(`/v1/capture/${action}`,base),{method:"POST",headers:{authorization:`Bearer ${session.accessToken}`,"content-type":action==="decode"?"application/octet-stream":"application/json",accept:"application/json"},body:body instanceof Uint8Array?new Uint8Array(body).buffer:body,cache:"no-store",redirect:"error",signal:AbortSignal.timeout(action==="decode"?15000:7000)})};
