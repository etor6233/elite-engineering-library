import "server-only";
import { applicationBaseUrl } from "@/platform/auth/oidc-client";
import { backendFetch } from "@/platform/backend/cloud-run-transport";
import type { DocumentConfig, DocumentTransport } from "./gateway";

export function documentConfig(): DocumentConfig {
  const env = process.env;
  if (env.DOCUMENTS_ENABLED !== "true" || !/^[a-f0-9]{64}$/u.test(env.DOCUMENTS_PROFILE_SHA256 ?? "") || !env.DOCUMENTS_TENANT_ID || !env.DOCUMENTS_ORGANIZATION_ID || !["FIXTURE", "PROVIDER"].includes(env.DOCUMENTS_MODE ?? "")) throw new Error("Document scope not configured");
  return { origin: applicationBaseUrl().origin, tenant: env.DOCUMENTS_TENANT_ID, organization: env.DOCUMENTS_ORGANIZATION_ID, profile: env.DOCUMENTS_PROFILE_SHA256!, mode: env.DOCUMENTS_MODE as "FIXTURE" | "PROVIDER" };
}

export const documentTransport: DocumentTransport = async (session, path, request) => {
  const base = new URL(process.env.ENTERPRISE_API_BASE_URL ?? "");
  if (!["http:", "https:"].includes(base.protocol) || base.username || base.password || base.search || base.hash || !/^\/v1\/documents\/[a-f0-9-]{36}(?:\/(?:original|process|review|decision|evidence\/(?:security|provider|analysis)))?$/u.test(path)) throw new Error("Invalid document backend route");
  const body = request.body instanceof Uint8Array ? new Uint8Array(request.body).buffer : request.body;
  return backendFetch(new URL(path, base), { method: request.method, ...(body === undefined ? {} : { body }), headers: { accept: "application/json", authorization: `Bearer ${session.accessToken}`, ...request.headers }, cache: "no-store", redirect: "error", signal: AbortSignal.timeout(path.endsWith("/process") ? 130000 : 5000) });
};
