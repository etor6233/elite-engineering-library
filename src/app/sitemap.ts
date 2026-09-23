import { publicSitemap, readPublicIndexing } from "@/platform/seo/public-indexing";
import { loadPublishedCatalog } from "@/platform/catalog/load";

export const dynamic = "force-dynamic";
export default async function sitemap() {
  const settings = readPublicIndexing();
  if (!settings.enabled) return publicSitemap(settings);
  return publicSitemap(settings, (await loadPublishedCatalog())?.models.map(model => model.code) ?? []);
}
