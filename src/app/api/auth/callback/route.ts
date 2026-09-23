import { NextResponse } from "next/server";
import { applicationBaseUrl, callbackUrl, oidcClient, oidcConfiguration } from "@/platform/auth/oidc-client";
import { cookieOptions, FLOW_COOKIE, openFlow, sealSession, SESSION_COOKIE } from "@/platform/auth/session";
import { portalLifecycleEnabled, portalProfile } from "@/platform/auth/portal-profile";
import { createPortalSession } from "@/platform/auth/portal-lifecycle";

function stringArray(value: unknown): string[] | null {
  if (!Array.isArray(value) || value.length > 100 || value.some((item) => typeof item !== "string" || item.length < 1 || item.length > 200)) return null;
  return [...new Set(value as string[])];
}

export async function GET(request: Request) {
  const flowCookie = request.headers.get("cookie")?.split(";").map((part) => part.trim()).find((part) => part.startsWith(`${FLOW_COOKIE}=`))?.slice(FLOW_COOKIE.length + 1);
  if (!flowCookie) return NextResponse.json({ code: "OIDC_FLOW_MISSING" }, { status: 400 });
  try {
    const flow = await openFlow(decodeURIComponent(flowCookie));
    const configuration = await oidcConfiguration();
    const authorizationResponse = new URL(callbackUrl());
    authorizationResponse.search = new URL(request.url).search;
    if (portalLifecycleEnabled()) {
      const stored = await createPortalSession(portalProfile(), flow, authorizationResponse);
      const response = NextResponse.redirect(new URL(flow.returnTo, applicationBaseUrl()), 303);
      response.cookies.set(SESSION_COOKIE, stored.cookie, cookieOptions(stored.maxAge));
      response.cookies.set(FLOW_COOKIE, "", cookieOptions(0));
      return response;
    }
    const tokens = await oidcClient.authorizationCodeGrant(configuration, authorizationResponse, { pkceCodeVerifier: flow.codeVerifier, expectedState: flow.state, expectedNonce: flow.nonce });
    const claims = tokens.claims();
    const organizations = stringArray(claims?.organization_ids);
    const permissions = stringArray(claims?.permissions);
    if (!tokens.access_token || typeof claims?.sub !== "string" || typeof claims.tenant_id !== "string" || !organizations?.length || !permissions) throw new Error("required identity claims missing");
    const ttl = Math.max(1, Math.min(tokens.expires_in ?? 300, 3600));
    const response = NextResponse.redirect(new URL(flow.returnTo, applicationBaseUrl()), 303);
    response.cookies.set(SESSION_COOKIE, await sealSession({ subject: claims.sub, tenantId: claims.tenant_id, permissions, organizations, accessToken: tokens.access_token }, ttl), cookieOptions(ttl));
    response.cookies.set(FLOW_COOKIE, "", cookieOptions(0));
    return response;
  } catch {
    const response = NextResponse.json({ code: "OIDC_CALLBACK_REJECTED" }, { status: 400 });
    response.cookies.set(FLOW_COOKIE, "", cookieOptions(0));
    return response;
  }
}
