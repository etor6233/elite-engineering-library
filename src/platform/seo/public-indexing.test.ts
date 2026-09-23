import { afterEach, describe, expect, it, vi } from "vitest";
import { publicPageMetadata, publicRobots, publicSitemap, readPublicIndexing, websiteJsonLd } from "./public-indexing";

afterEach(() => vi.unstubAllEnvs());
function enable(origin = "https://public.example.invalid") {
  vi.stubEnv("PUBLIC_INDEXING_ENABLED", "1"); vi.stubEnv("PUBLIC_SITE_ORIGIN", origin);
}

describe("public indexing boundary", () => {
  it("defaults to no public crawl or fabricated URLs", () => {
    const settings = readPublicIndexing({});
    expect(publicSitemap(settings)).toEqual([]);
    expect(publicRobots(settings)).toEqual({ rules: { userAgent: "*", disallow: "/" } });
  });
  it.each(["true", "yes", "2", " 1"])("rejects ambiguous enable flag %s", flag => {
    expect(() => readPublicIndexing({ PUBLIC_INDEXING_ENABLED: flag })).toThrow("Invalid public indexing configuration");
  });
  it("requires explicit configured origin for indexing", () => {
    expect(() => readPublicIndexing({ PUBLIC_INDEXING_ENABLED: "1" })).toThrow("explicit HTTPS origin");
  });
  it.each(["http://public.example.invalid", "https://user:PRIVATE_PASSWORD@example.invalid", "https://example.invalid/private", "https://example.invalid/?secret=value", "https://example.invalid/#fragment", "https://example.invalid:444", " https://example.invalid", "https://example.invalid\n", "https://example.invalid\\private", "https://example.invalid/%2e%2e", "https://example.invalid/path/.."])("rejects non-origin input without reflecting it: %s", origin => {
    let message = "";
    try { readPublicIndexing({ PUBLIC_SITE_ORIGIN: origin }); } catch (error) { message = (error as Error).message; }
    expect(message).toBe("Invalid public site origin"); expect(message).not.toContain(origin);
  });
  it("normalizes HTTPS hostname without taking an untrusted request host", () => {
    expect(readPublicIndexing({ PUBLIC_SITE_ORIGIN: "https://PUBLIC.example.invalid:443/", PUBLIC_INDEXING_ENABLED: "1" })).toEqual({ enabled: true, origin: "https://public.example.invalid" });
  });
  it("publishes only two real pages with no invented dates or private identifiers", () => {
    enable();
    expect(publicSitemap()).toEqual([{ url: "https://public.example.invalid/" }, { url: "https://public.example.invalid/models" }]);
    expect(publicRobots()).toEqual({ rules: { userAgent: "*", disallow: "/", allow: ["/$", "/models$", "/_next/static/", "/icon.svg$"] }, sitemap: "https://public.example.invalid/sitemap.xml" });
    expect(publicPageMetadata("/models", "Modelos").alternates).toEqual({ canonical: "https://public.example.invalid/models" });
  });
  it("keeps an explicit disabled deployment out of the sitemap", () => {
    enable(); vi.stubEnv("PUBLIC_INDEXING_ENABLED", "0");
    expect(publicSitemap()).toEqual([]); expect(websiteJsonLd("Fixture", "es-AR")).toBeNull();
    expect(publicPageMetadata("/").robots).toEqual({ index: false, follow: false });
  });
  it("rejects caller attempts to canonize a protected page", () => {
    enable(); expect(() => publicPageMetadata("/admin" as "/")).toThrow("Unknown public page");
  });
  it("preserves structured text but prevents HTML script termination", () => {
    enable(); const name = '</script><script>alert("PRIVATE")</script>&\u2028\u2029';
    const encoded = websiteJsonLd(name, "es-AR")!;
    expect(encoded).not.toMatch(/[<>&\u2028\u2029]/u);
    expect(JSON.parse(encoded)).toEqual({ "@context": "https://schema.org", "@type": "WebSite", name, url: "https://public.example.invalid/", inLanguage: "es-AR" });
  });
});
