import { createHash } from "node:crypto";
import { loadPublishedCatalog } from "@/platform/catalog/load";
import { readPublishedCatalogPNG } from "@/platform/backend/catalog-client";
export const dynamic = "force-dynamic";
export async function GET(_request: Request, { params }: { params: Promise<{ sha: string }> }) {
  const { sha } = await params;
  if (!/^[0-9a-f]{64}$/u.test(sha)) return new Response(null, { status: 404 });
  const catalog = await loadPublishedCatalog();
  if (!catalog?.models.some(model => model.media.sha256 === sha)) return new Response(null, { status: 404 });
  const response = await readPublishedCatalogPNG(sha);
  if (!response.ok || response.headers.get("content-type") !== "image/png" || !response.body) {
    if (response.body) void response.body.cancel().catch(() => {});
    return new Response(null, { status: response.status === 404 ? 404 : 502 });
  }
  const reader = response.body.getReader(), parts: Uint8Array[] = [];
  let size = 0, done = false;
  try {
    while (true) {
      const item = await reader.read(); if (item.done) { done = true; break; }
      size += item.value.byteLength; if (size > 4_194_304) return new Response(null, { status: 502 });
      parts.push(item.value);
    }
    const body = Buffer.concat(parts);
    if (createHash("sha256").update(body).digest("hex") !== sha) return new Response(null, { status: 502 });
    return new Response(new Uint8Array(body), { headers: { "content-type": "image/png", "x-content-type-options": "nosniff", "cache-control": "public, max-age=0, must-revalidate", etag: `"${sha}"` } });
  } finally {
    if (!done) void reader.cancel().catch(() => {});
    reader.releaseLock();
  }
}
