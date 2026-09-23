import { NextResponse } from "next/server";
import { callbackUrl, oidcClient, oidcConfiguration } from "@/platform/auth/oidc-client";
import { cookieOptions, FLOW_COOKIE, sealFlow } from "@/platform/auth/session";
import { portalLifecycleEnabled, portalProfile } from "@/platform/auth/portal-profile";

const destinations = new Set(["/admin", "/customer", "/factory", "/franchise"]);

export async function GET(request: Request) {
  const requested = new URL(request.url).searchParams.get("return_to") ?? "/customer";
  const returnTo = destinations.has(requested) ? requested : "/customer";
  const configuration = await oidcConfiguration();
  const codeVerifier = oidcClient.randomPKCECodeVerifier();
  const state = oidcClient.randomState();
  const nonce = oidcClient.randomNonce();
  const redirect = oidcClient.buildAuthorizationUrl(configuration, {
    redirect_uri: callbackUrl(), scope: portalLifecycleEnabled() ? portalProfile().scopes.join(" ") : "openid profile email", response_type: "code",
    code_challenge: await oidcClient.calculatePKCECodeChallenge(codeVerifier), code_challenge_method: "S256", state, nonce,
  });
  const response = NextResponse.redirect(redirect, 303);
  response.cookies.set(FLOW_COOKIE, await sealFlow({ codeVerifier, state, nonce, returnTo }), cookieOptions(600));
  return response;
}
