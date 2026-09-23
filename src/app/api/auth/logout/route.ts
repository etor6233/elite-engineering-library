import { NextResponse } from "next/server";
import { applicationBaseUrl } from "@/platform/auth/oidc-client";
import { cookieOptions, SESSION_COOKIE } from "@/platform/auth/session";
import { portalLifecycleEnabled, portalProfile } from "@/platform/auth/portal-profile";
import { logoutPortalSession } from "@/platform/auth/portal-lifecycle";

export async function POST(request: Request) {
  const origin = request.headers.get("origin");
  if (origin !== applicationBaseUrl().origin) return NextResponse.json({ code: "ORIGIN_REJECTED" }, { status: 403 });
  if (portalLifecycleEnabled()) {
    const cookie=request.headers.get("cookie")?.split(";").map(v=>v.trim()).find(v=>v.startsWith(SESSION_COOKIE+"="))?.slice(SESSION_COOKIE.length+1);
    if(cookie){
      try {
        const result=await logoutPortalSession(portalProfile(),decodeURIComponent(cookie));
        if(!result.providerRevocationComplete) return NextResponse.json({code:"PROVIDER_REVOCATION_PENDING",local_session_revoked:true,retry:"POST /api/auth/logout"},{status:503,headers:{"cache-control":"no-store"}});
        const response=NextResponse.redirect(result.redirect,303);response.cookies.set(SESSION_COOKIE,"",cookieOptions(0));return response;
      } catch { return NextResponse.json({code:"LOGOUT_NOT_CONFIRMED",retry:"POST /api/auth/logout"},{status:503,headers:{"cache-control":"no-store"}}); }
    }
  }
  const response = NextResponse.redirect(new URL("/", applicationBaseUrl()), 303);
  response.cookies.set(SESSION_COOKIE, "", cookieOptions(0));
  return response;
}
