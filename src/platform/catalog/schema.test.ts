import { afterEach, describe, expect, it, vi } from "vitest";
import { parsePublishedCatalog } from "./schema";
import { publicModelMetadata, publicRobots, publicSitemap } from "../seo/public-indexing";
const origin = "https://catalog.example.invalid", hash = "a".repeat(64);
function publication() {
  return { tenant_code: "catalog-j3", generation: "2", source_sha256: hash,
    effective_price_book_id: "effective-2", market: "AR", currency: "ARS",
    models: [{ id: "one", code: "bicycle-one", displayName: "Approved", vehicleClass: "bicycle",
      specification: {}, media: { id: "png", original_sha256: hash, sha256: hash, width: 2, height: 2 },
      canonical_url: origin + "/models/bicycle-one" }],
    variants: [{ id: "v-one", model_id: "one", code: "v-one", display_name: "Approved variant",
      battery_specification: {}, homologation_state: "unknown", amount_minor_units: "120001", tax_mode: "inclusive" }] };
}
afterEach(() => vi.unstubAllEnvs());
describe("approved publication storefront boundary", () => {
  it("preserves quoted int64 money and binds the configured tenant and canonical origin", () => {
    const p = publication(); p.variants[0]!.amount_minor_units = "9223372036854775807";
    expect(parsePublishedCatalog(p, "catalog-j3", origin)).toEqual(p);
  });
  it("rejects another tenant, origin, draft-only field and a variant outside the snapshot", () => {
    expect(() => parsePublishedCatalog(publication(), "other-tenant", origin)).toThrow("tenant binding");
    expect(() => parsePublishedCatalog(publication(), "catalog-j3", "https://other.invalid")).toThrow("canonical binding");
    expect(() => parsePublishedCatalog({ ...publication(), profile: {} }, "catalog-j3", origin)).toThrow();
    const p = publication(); p.variants[0]!.model_id = "absent";
    expect(() => parsePublishedCatalog(p, "catalog-j3", origin)).toThrow("variant binding");
  });
  it.each(["/admin", "/models/other", "/models/bicycle-one?secret=1", "/models/bicycle-one#private"])("rejects mismatched canonical path %s", path => {
    const p = publication(); p.models[0]!.canonical_url = origin + path;
    expect(() => parsePublishedCatalog(p, "catalog-j3", origin)).toThrow();
  });
  it.each(["-1", "01", "1.1", "9223372036854775808"])("rejects money outside the source int64 contract %s", value => {
    const p = publication(); p.variants[0]!.amount_minor_units = value;
    expect(() => parsePublishedCatalog(p, "catalog-j3", origin)).toThrow();
  });
  it("exposes only approved real model paths and keeps the disabled profile out of crawling", () => {
    const settings = { enabled: true, origin };
    expect(publicSitemap(settings, ["bicycle-one"])).toEqual([{url: origin + "/"}, {url: origin + "/models"}, {url: origin + "/models/bicycle-one"}]);
    expect(publicRobots(settings, ["bicycle-one"]).rules).toEqual({userAgent:"*",disallow:"/",allow:["/$","/models$","/models/bicycle-one$","/_next/static/","/icon.svg$"]});
    expect(publicSitemap({ ...settings, enabled: false }, ["bicycle-one"])).toEqual([]);
    expect(() => publicSitemap(settings, ["../../admin"])).toThrow();
    vi.stubEnv("PUBLIC_SITE_ORIGIN", origin); vi.stubEnv("PUBLIC_INDEXING_ENABLED", "1");
    expect(publicModelMetadata("bicycle-one", "Approved", origin + "/models/bicycle-one")).toEqual({
      title: "Approved", robots: {index:true,follow:true}, alternates:{canonical:origin+"/models/bicycle-one"} });
    expect(() => publicModelMetadata("bicycle-one", "Approved", origin + "/admin")).toThrow();
  });
});
