import type { Metadata, MetadataRoute } from "next";

// Only these existing public pages are approved by this library reference.
// Protected portals and dynamic business objects are never inferred from routes.
const publicPaths = ["/", "/models"] as const;
export type PublicPage = typeof publicPaths[number];
export type PublicIndexing = Readonly<{ enabled: boolean; origin: string | null }>;

export function readPublicIndexing(input: {
  PUBLIC_SITE_ORIGIN?: string | undefined;
  PUBLIC_INDEXING_ENABLED?: string | undefined;
} = { PUBLIC_SITE_ORIGIN: process.env.PUBLIC_SITE_ORIGIN,
  PUBLIC_INDEXING_ENABLED: process.env.PUBLIC_INDEXING_ENABLED }): PublicIndexing {
  const flag = input.PUBLIC_INDEXING_ENABLED;
  if (flag !== undefined && flag !== "" && flag !== "0" && flag !== "1") {
    throw new Error("Invalid public indexing configuration");
  }
  const enabled = flag === "1";
  const raw = input.PUBLIC_SITE_ORIGIN;
  if (!raw) {
    if (enabled) throw new Error("Public indexing requires an explicit HTTPS origin");
    return { enabled: false, origin: null };
  }
  // This value is deployment configuration, never Host/X-Forwarded-Host input.
  // Refuse URL normalization that could hide credentials, controls or paths.
  if (raw.length > 2048 || /[\s\\%?#]/u.test(raw)) throw new Error("Invalid public site origin");
  let url: URL;
  try { url = new URL(raw); } catch { throw new Error("Invalid public site origin"); }
  if (url.protocol !== "https:" || !url.hostname || url.username || url.password ||
      url.pathname !== "/" || url.port || !/^https:\/\/[^/]+\/?$/u.test(raw)) {
    throw new Error("Invalid public site origin");
  }
  return { enabled, origin: url.origin };
}

export function publicPageMetadata(page: PublicPage, title?: string): Metadata {
  if (!publicPaths.includes(page)) throw new Error("Unknown public page");
  const settings = readPublicIndexing();
  return {
    ...(title ? { title } : {}),
    robots: { index: settings.enabled, follow: settings.enabled },
    ...(settings.origin ? { alternates: { canonical: new URL(page, settings.origin).href } } : {})
  };
}

export function publicModelPath(code: string): string {
  if (!/^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$/u.test(code) || code.length > 128) throw new Error("Invalid public model code");
  return "/models/" + code;
}
export function publicModelMetadata(code: string, title: string, canonical: string): Metadata {
  const path = publicModelPath(code), settings = readPublicIndexing(), uri = new URL(canonical);
  if (uri.protocol !== "https:" || uri.username || uri.password || uri.search || uri.hash ||
      uri.pathname !== path || (settings.origin !== null && uri.origin !== settings.origin)) throw new Error("Invalid public model canonical");
  return { title, robots: { index: settings.enabled, follow: settings.enabled }, alternates: { canonical: uri.href } };
}
function approvedPublicPaths(modelCodes: readonly string[]) {
  if (modelCodes.length > 32) throw new Error("Too many published model paths");
  return [...publicPaths, ...Array.from(new Set(modelCodes.map(publicModelPath))).sort()];
}
export function publicRobots(settings = readPublicIndexing(), modelCodes: readonly string[] = []): MetadataRoute.Robots {
  if (!settings.enabled || !settings.origin) return { rules: { userAgent: "*", disallow: "/" } };
  // Exact public-path matches: adding a portal does not expose it to crawlers.
  // This is a crawl instruction, not access control or removal from an index.
  return { rules: { userAgent: "*", disallow: "/", allow: [...approvedPublicPaths(modelCodes).map(path => path + "$"), "/_next/static/", "/icon.svg$"] },
    sitemap: new URL("/sitemap.xml", settings.origin).href };
}

export function publicSitemap(settings = readPublicIndexing(), modelCodes: readonly string[] = []): MetadataRoute.Sitemap {
  if (!settings.enabled || !settings.origin) return [];
  // Callers supply only codes from the validated current approved publication.
  // No inferred private routes or invented lastModified dates.
  return approvedPublicPaths(modelCodes).map(path => ({ url: new URL(path, settings.origin!).href }));
}

export function websiteJsonLd(name: string, language: string): string | null {
  const settings = readPublicIndexing();
  if (!settings.enabled || !settings.origin) return null;
  const data = { "@context": "https://schema.org", "@type": "WebSite", name,
    url: new URL("/", settings.origin).href, inLanguage: language };
  // A JSON string embedded in an HTML script must not contain a closing tag.
  return JSON.stringify(data).replace(/[<>&\u2028\u2029]/gu, character =>
    "\\u" + character.charCodeAt(0).toString(16).padStart(4, "0"));
}
