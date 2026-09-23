// AUTHORED narrow transport glue over Google's documented metadata/Cloud Run
// protocol. No Google SDK/source code copied. Server-only callers own authZ.
const METADATA = "http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/identity";
let cached: { audience: string; value: string; expires: number } | undefined;
let pending: { audience: string; request: Promise<string> } | undefined;

function serviceUrl(value: string | undefined): URL {
  if (!value) throw new Error("Cloud service identity configuration missing");
  const url = new URL(value);
  if (url.protocol !== "https:" || url.username || url.password || url.port || url.search || url.hash || !/^[a-z0-9][a-z0-9.-]*\.run\.app$/.test(url.hostname) || url.pathname !== "/") throw new Error("Cloud service identity URL invalid");
  return url;
}

async function metadataToken(audience: string): Promise<string> {
  const now = Math.floor(Date.now() / 1000);
  if (cached?.audience === audience && cached.expires > now) return cached.value;
  if (pending?.audience === audience) return pending.request;
  const request = (async () => {
    const url = new URL(METADATA);
    url.searchParams.set("audience", audience);
    url.searchParams.set("format", "full");
    const response = await fetch(url, { method: "GET", headers: { "Metadata-Flavor": "Google" }, cache: "no-store", redirect: "error", signal: AbortSignal.timeout(1500) });
    if (!response.ok || response.headers.get("Metadata-Flavor") !== "Google" || !response.body) throw new Error("Cloud service identity unavailable");
    const reader = response.body.getReader();
    const parts: Uint8Array[] = []; let length = 0;
    try {
      while (true) {
        const part = await reader.read(); if (part.done) break;
        length += part.value.byteLength;
        if (length > 16384) throw new Error("Cloud service identity response invalid");
        parts.push(part.value);
      }
    } catch {
      void reader.cancel().catch(() => {});
      throw new Error("Cloud service identity response invalid");
    } finally { reader.releaseLock(); }
    const bytes = new Uint8Array(length); let offset = 0;
    for (const part of parts) { bytes.set(part, offset); offset += part.byteLength; }
    const token = new TextDecoder("utf-8", { fatal: true }).decode(bytes);
    if (!/^[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+$/.test(token)) throw new Error("Cloud service identity response invalid");
    const payload: unknown = JSON.parse(Buffer.from(token.split(".")[1]!, "base64url").toString("utf8"));
    if (!payload || typeof payload !== "object" || !("aud" in payload) || payload.aud !== audience || !("exp" in payload) || typeof payload.exp !== "number" || !Number.isFinite(payload.exp) || payload.exp <= now + 60 || payload.exp > now + 7200) throw new Error("Cloud service identity claims invalid");
    // This is cache/transport validation only. Cloud Run validates the signature.
    cached = { audience, value: token, expires: Math.min(payload.exp - 60, now + 300) };
    return token;
  })();
  pending = { audience, request };
  try { return await request; } finally { if (pending?.request === request) pending = undefined; }
}

export async function backendFetch(input: URL, init: RequestInit = {}): Promise<Response> {
  const mode = process.env.ELITE_CLOUD_RUN_IDENTITY;
  if (mode === undefined || mode === "disabled") return fetch(input, init);
  if (mode !== "metadata" || typeof window !== "undefined") throw new Error("Cloud service identity mode invalid");
  const destination = serviceUrl(process.env.ENTERPRISE_API_BASE_URL);
  const audience = serviceUrl(process.env.ELITE_CLOUD_RUN_AUDIENCE).origin;
  if (input.origin !== destination.origin || input.username || input.password) throw new Error("Cloud service identity destination mismatch");
  const token = await metadataToken(audience);
  const headers = new Headers(init.headers);
  // Preserve the caller's end-user OIDC Authorization. This second header is
  // consumed by Cloud Run IAM and never grants a product tenant/role permission.
  headers.set("X-Serverless-Authorization", `Bearer ${token}`);
  return fetch(input, { ...init, headers, cache: "no-store", redirect: "error", signal: init.signal ?? AbortSignal.timeout(5000) });
}

export function clearCloudIdentityCacheForTests() { cached = undefined; pending = undefined; }
