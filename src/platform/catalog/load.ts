import "server-only";
import { cache } from "react";
import { readPublishedCatalogSource } from "@/platform/backend/catalog-client";
import { readPublicIndexing } from "@/platform/seo/public-indexing";
import { parsePublishedCatalog } from "./schema";

// Request-local React memoization keeps metadata and page on the same snapshot.
export const loadPublishedCatalog = cache(async () => {
  const source = await readPublishedCatalogSource();
  return source === null ? null : parsePublishedCatalog(source.value, source.tenant, readPublicIndexing().origin);
});
