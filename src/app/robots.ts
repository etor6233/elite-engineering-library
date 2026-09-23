import { publicRobots, readPublicIndexing } from "@/platform/seo/public-indexing";
import { loadPublishedCatalog } from "@/platform/catalog/load";

export const dynamic = "force-dynamic";
export default async function robots() {
  const settings = readPublicIndexing();
  if (!settings.enabled) return publicRobots(settings);
  return publicRobots(settings, (await loadPublishedCatalog())?.models.map(model => model.code) ?? []);
}
