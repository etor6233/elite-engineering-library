import "server-only";
import { createHash } from "node:crypto";
import { cookies } from "next/headers";
import { EncryptJWT, jwtDecrypt, type JWTPayload } from "jose";
import { portalLifecycleEnabled, portalProfile, portalSecret } from "./portal-profile";
import { readPortalSession } from "./portal-lifecycle";

export const FLOW_COOKIE = process.env.NODE_ENV === "production" ? "__Host-elite_oidc_flow" : "elite_oidc_flow";
export const SESSION_COOKIE = process.env.NODE_ENV === "production" ? "__Host-elite_session" : "elite_session";

export type AuthFlow = {
  codeVerifier: string;
  state: string;
  nonce: string;
  returnTo: string;
};

export type PortalSession = {
  subject: string;
  tenantId: string;
  permissions: string[];
  organizations: string[];
  accessToken: string;
};

const MAX_ACCESS_TOKEN_CHARACTERS = 2_000;
const MAX_SESSION_COOKIE_CHARACTERS = 3_800;

function key() {
  const secret = portalLifecycleEnabled() ? portalSecret("OIDC_PORTAL_SESSION_KEY_FILE") : process.env.AUTH_SESSION_SECRET;
  if (!secret || secret.length < 32) throw new Error("AUTH_SESSION_SECRET must contain at least 32 characters");
  return createHash("sha256").update(secret, "utf8").digest();
}

async function seal(payload: JWTPayload, ttlSeconds: number) {
  return new EncryptJWT(payload)
    .setProtectedHeader({ alg: "dir", enc: "A256GCM", typ: "JWT" })
    .setIssuedAt()
    .setExpirationTime(`${ttlSeconds}s`)
    .encrypt(key());
}

async function open(token: string): Promise<JWTPayload> {
  const result = await jwtDecrypt(token, key(), { keyManagementAlgorithms: ["dir"], contentEncryptionAlgorithms: ["A256GCM"] });
  return result.payload;
}

function strings(value: unknown): string[] | null {
  if (!Array.isArray(value) || value.length > 100 || value.some((item) => typeof item !== "string" || item.length < 1 || item.length > 200)) return null;
  return [...new Set(value as string[])];
}

export async function sealFlow(flow: AuthFlow) { return seal(flow, 600); }

export async function openFlow(token: string): Promise<AuthFlow> {
  const value = await open(token);
  if (typeof value.codeVerifier !== "string" || typeof value.state !== "string" || typeof value.nonce !== "string" || typeof value.returnTo !== "string") throw new Error("invalid OIDC flow cookie");
  return { codeVerifier: value.codeVerifier, state: value.state, nonce: value.nonce, returnTo: value.returnTo };
}

export async function sealSession(session: PortalSession, ttlSeconds: number) {
  if (ttlSeconds < 1 || ttlSeconds > 3600) throw new Error("invalid session lifetime");
  if (session.accessToken.length < 1 || session.accessToken.length > MAX_ACCESS_TOKEN_CHARACTERS) throw new Error("invalid access token size");
  const token = await seal(session, ttlSeconds);
  if (token.length > MAX_SESSION_COOKIE_CHARACTERS) throw new Error("portal session exceeds cookie budget");
  return token;
}

export async function openSession(token: string): Promise<PortalSession> {
  const value = await open(token);
  const permissions = strings(value.permissions);
  const organizations = strings(value.organizations);
  if (typeof value.subject !== "string" || typeof value.tenantId !== "string" || typeof value.accessToken !== "string" || !permissions || !organizations || organizations.length === 0) throw new Error("invalid portal session");
  return { subject: value.subject, tenantId: value.tenantId, accessToken: value.accessToken, permissions, organizations };
}

export async function readSession(): Promise<PortalSession | null> {
  const token = (await cookies()).get(SESSION_COOKIE)?.value;
  if (!token) return null;
  if (portalLifecycleEnabled()) return readPortalSession(portalProfile(), token);
  try { return await openSession(token); } catch { return null; }
}

export function allowed(session: PortalSession, permission: string) {
  return session.permissions.includes("*") || session.permissions.includes(permission);
}

export function cookieOptions(maxAge: number) {
  return { httpOnly: true, secure: process.env.NODE_ENV === "production", sameSite: "lax" as const, path: "/", maxAge };
}
