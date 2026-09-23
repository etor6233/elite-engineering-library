import "server-only";
import { createHash } from "node:crypto";
import { openSync, closeSync, readSync, fstatSync } from "node:fs";
import { z } from "zod";

const atom = z.string().min(1).max(200).regex(/^[^\s\x00]+$/).refine(v => v !== "*");
const set = z.array(atom).min(1).max(100).refine(v => new Set(v).size === v.length);
const schema = z.object({
  schema: z.literal("elite.oidc.portal-lifecycle.v1"), profile_id: atom, revision: z.literal(1), transport: z.enum(["TLS", "LOOPBACK_FIXTURE"]),
  issuer: z.string(), authorization_endpoint: z.string(), token_endpoint: z.string(), jwks_endpoint: z.string(), revocation_endpoint: z.string(), end_session_endpoint: z.string(),
  client_id: atom, callback_url: z.string(), post_logout_url: z.string(), bridge_url: z.string(), backend_audience: atom,
  service_client_id: atom, service_subject: atom, service_scopes: set, service_audience_parameter: z.string(),
  tenant_id: atom, organization_ids: set, allowed_permissions: set, scopes: set,
  maximum_session_seconds: z.number().int().min(300).max(86400), refresh_before_seconds: z.number().int().min(5).max(300),
  retention_seconds:z.number().int().min(3600).max(604800),
}).strict();
export type PortalProfile = z.infer<typeof schema> & { documentSHA256: string };
export function boundedFile(path: string, max: number): Buffer {
  const fd = openSync(path, "r");
  try { if (!fstatSync(fd).isFile()) throw new Error("OIDC_CONFIGURATION_REJECTED"); const bytes = Buffer.alloc(max + 1); let used = 0; while (used < bytes.length) { const n = readSync(fd, bytes, used, bytes.length - used, null); if (!n) break; used += n; } if (used > max) throw new Error("OIDC_CONFIGURATION_REJECTED"); return bytes.subarray(0, used); }
  finally { closeSync(fd); }
}
export function portalSecret(name: string): string {
  const path = process.env[name]; if (!path) throw new Error("OIDC_CONFIGURATION_REJECTED");
  const value = new TextDecoder("utf-8", { fatal: true }).decode(boundedFile(path, 4096)).replace(/\r?\n$/, "");
  if (!value || /[\x00\r\n]/.test(value)) throw new Error("OIDC_CONFIGURATION_REJECTED"); return value;
}
export function parsePortalProfile(raw: Uint8Array, digest: string): PortalProfile {
  if (!raw.length || raw.length > 16384 || !/^[a-f0-9]{64}$/.test(digest) || createHash("sha256").update(raw).digest("hex") !== digest) throw new Error("OIDC_CONFIGURATION_REJECTED");
  const text = new TextDecoder("utf-8", { fatal: true }).decode(raw);
  // Schema contains only scalar/array members. JSON token walk rejects duplicate root names.
  const keys: string[] = []; let depth = 0; const tokens = text.match(/"(?:[^"\\]|\\.)*"|[{}\[\]:,]|[^\s{}\[\]:,]+/g) ?? [];
  for (let i=0;i<tokens.length;i++) { const token=tokens[i]; if (token==="{"||token==="[") depth++; else if(token==="}"||token==="]") depth--; else if(depth===1 && token?.startsWith('"') && tokens[i+1]===":") keys.push(JSON.parse(token) as string); }
  if (new Set(keys).size !== keys.length) throw new Error("OIDC_CONFIGURATION_REJECTED");
  const p=schema.parse(JSON.parse(text));
  for (const rawURL of [p.issuer,p.authorization_endpoint,p.token_endpoint,p.jwks_endpoint,p.revocation_endpoint,p.end_session_endpoint,p.callback_url,p.post_logout_url,p.bridge_url]) {
    const u=new URL(rawURL); const local=["127.0.0.1","[::1]"].includes(u.hostname);
    if(u.username||u.password||u.search||u.hash||rawURL.length>2048||!(p.transport==="TLS" ? u.protocol==="https:" : local&&["http:","https:"].includes(u.protocol))) throw new Error("OIDC_CONFIGURATION_REJECTED");
  }
  if(!p.callback_url.endsWith("/api/auth/callback")||p.callback_url.slice(0,-18)+"/"!==p.post_logout_url||p.bridge_url.endsWith("/")||new URL(p.bridge_url).pathname!=="/"||!p.scopes.includes("openid")||!p.scopes.includes("offline_access")||p.refresh_before_seconds*2>=p.maximum_session_seconds||(p.service_audience_parameter!==""&&p.service_audience_parameter!==p.backend_audience)) throw new Error("OIDC_CONFIGURATION_REJECTED");
  return Object.freeze({...p, organization_ids:Object.freeze([...p.organization_ids]) as unknown as string[], allowed_permissions:Object.freeze([...p.allowed_permissions]) as unknown as string[], scopes:Object.freeze([...p.scopes]) as unknown as string[], service_scopes:Object.freeze([...p.service_scopes]) as unknown as string[],documentSHA256:digest});
}
export function portalLifecycleEnabled() { const value=process.env.OIDC_PORTAL_LIFECYCLE_ENABLED; if(value!==undefined&&value!=="false"&&value!=="true")throw new Error("OIDC_CONFIGURATION_REJECTED");return value==="true"; }
export function portalProfile(): PortalProfile {
  if(!portalLifecycleEnabled())throw new Error("OIDC_LIFECYCLE_DISABLED");
  const path=process.env.OIDC_PORTAL_PROFILE_FILE, hash=process.env.OIDC_PORTAL_PROFILE_SHA256;
  if(!path||!hash)throw new Error("OIDC_CONFIGURATION_REJECTED");
  const p=parsePortalProfile(boundedFile(path,16384),hash);
  if(process.env.NODE_ENV==="production"&&p.transport!=="TLS")throw new Error("OIDC_CONFIGURATION_REJECTED");
  if(process.env.APP_BASE_URL&&new URL(process.env.APP_BASE_URL).origin!==new URL(p.callback_url).origin)throw new Error("OIDC_CONFIGURATION_REJECTED");
  return p;
}
