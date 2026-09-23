// AUTHORED typed projection validation; no local price or business rules.
import { z } from "zod";
const id = z.string().regex(/^[A-Za-z0-9][A-Za-z0-9._:-]{0,63}$/u);
const code = z.string().regex(/^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$/u);
const digest = z.string().regex(/^[0-9a-f]{64}$/u);
const unsigned = z.string().regex(/^(0|[1-9][0-9]{0,18})$/u).refine(value => BigInt(value) <= 9223372036854775807n);
export const publishedCatalogSchema = z.object({
  tenant_code: code, generation: unsigned.refine(value => BigInt(value) > 0n),
  source_sha256: digest, effective_price_book_id: id,
  market: z.string().regex(/^[A-Z]{2}$/u), currency: z.string().regex(/^[A-Z]{3}$/u),
  models: z.array(z.object({
    id, code, displayName: z.string().min(2).max(160), vehicleClass: z.enum(["motorcycle", "bicycle", "scooter", "utility", "other"]),
    specification: z.record(z.string(), z.unknown()),
    media: z.object({ id, original_sha256: digest, sha256: digest, width: z.number().int().min(1).max(2048), height: z.number().int().min(1).max(2048) }).strict(),
    canonical_url: z.string().url()
  }).strict()).min(1).max(32),
  variants: z.array(z.object({
    id, model_id: id, code: z.string().min(1).max(128), display_name: z.string().min(1).max(160),
    battery_specification: z.record(z.string(), z.unknown()),
    homologation_state: z.enum(["unknown", "pending", "approved", "rejected", "expired"]),
    amount_minor_units: unsigned, tax_mode: z.enum(["inclusive", "exclusive", "not-applicable"])
  }).strict()).min(1).max(128)
}).strict();
export type PublishedCatalog = z.infer<typeof publishedCatalogSchema>;
export function parsePublishedCatalog(value: unknown, tenant: string, origin: string | null): PublishedCatalog {
  const result = publishedCatalogSchema.parse(value);
  if (result.tenant_code !== tenant) throw new Error("Catalog tenant binding mismatch");
  const ids = new Set<string>(), codes = new Set<string>(), variants = new Set<string>();
  for (const model of result.models) {
    if (ids.has(model.id) || codes.has(model.code) || model.media.width * model.media.height > 1_048_576) throw new Error("Invalid catalog model identity");
    ids.add(model.id); codes.add(model.code);
    const uri = new URL(model.canonical_url);
    if (uri.protocol !== "https:" || uri.username || uri.password || uri.search || uri.hash ||
        uri.pathname !== "/models/" + model.code || (origin !== null && uri.origin !== origin)) throw new Error("Catalog canonical binding mismatch");
  }
  for (const variant of result.variants) {
    if (!ids.has(variant.model_id) || variants.has(variant.id)) throw new Error("Invalid catalog variant binding");
    variants.add(variant.id);
  }
  return result;
}
