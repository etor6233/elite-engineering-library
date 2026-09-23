import * as oidc from "openid-client";
import { portalLifecycleEnabled, portalProfile } from "./portal-profile";
import { portalConfiguration } from "./portal-protocol";

let cached: Promise<oidc.Configuration> | undefined;

function required(name: string) {
  const value = process.env[name];
  if (!value) throw new Error(`${name} is required`);
  return value;
}

function secureUrl(name: string) {
  const value = new URL(required(name));
  const localDevelopment = process.env.NODE_ENV !== "production" && ["localhost", "127.0.0.1"].includes(value.hostname);
  if (value.username || value.password || value.search || value.hash || (value.protocol !== "https:" && !localDevelopment)) throw new Error(`${name} must be a safe HTTPS URL`);
  return value;
}

function localInsecure(url: URL) {
  return process.env.NODE_ENV !== "production" && url.protocol === "http:" && ["localhost", "127.0.0.1"].includes(url.hostname);
}

export function applicationBaseUrl() { return portalLifecycleEnabled() ? new URL(portalProfile().post_logout_url) : secureUrl("APP_BASE_URL"); }
export function callbackUrl() { return portalLifecycleEnabled() ? portalProfile().callback_url : new URL("/api/auth/callback", applicationBaseUrl()).href; }

export function oidcConfiguration() {
  if (portalLifecycleEnabled()) return portalConfiguration(portalProfile());
  const issuer = secureUrl("OIDC_ISSUER");
  cached ??= oidc.discovery(issuer, required("OIDC_CLIENT_ID"), required("OIDC_CLIENT_SECRET"), undefined, localInsecure(issuer) ? { execute: [oidc.allowInsecureRequests] } : undefined);
  return cached;
}

export const oidcClient = oidc;
