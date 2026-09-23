import { backendFetch } from "./cloud-run-transport";
import "server-only";
import { publishedCatalogTransportSettings, readBackendResponse, BackendProblem } from "./public-client";
export async function readPublishedCatalogSource(): Promise<{ value: unknown; tenant: string } | null> {
  if (process.env.CATALOG_RELEASE_ENABLED !== "true") return null;
  const { base, tenant } = publishedCatalogTransportSettings();
  const response = await backendFetch(new URL("/v1/public/catalog", base), { cache: "no-store", redirect: "error", signal: AbortSignal.timeout(5000), headers: { accept: "application/json" } });
  if (response.status === 404) { if (response.body) void response.body.cancel().catch(() => {}); return null; }
  return { value: await readBackendResponse<unknown>(response), tenant };
}
export async function readPublishedCatalogPNG(sha: string): Promise<Response> {
  if (process.env.CATALOG_RELEASE_ENABLED !== "true" || !/^[0-9a-f]{64}$/u.test(sha)) throw new BackendProblem(404, "CATALOG_NOT_FOUND");
  const { base } = publishedCatalogTransportSettings();
  return backendFetch(new URL("/v1/public/catalog/media/" + sha, base), { cache: "no-store", redirect: "error", signal: AbortSignal.timeout(5000), headers: { accept: "image/png" } });
}
