import { PublicCatalogue } from "@/components/public-web/public-catalogue";
import PublicLoading from "@/components/public-web/public-loading";
import {Suspense} from "react";
import { loadPublicLocale } from "@/platform/i18n/load-public-locale";
import { publicMessage as message } from "@/platform/i18n/public-catalog";
import { publicPageMetadata } from "@/platform/seo/public-indexing";
export async function generateMetadata(){const locale=await loadPublicLocale();return publicPageMetadata("/models",message(locale.locale,"models.short"));}
export const dynamic="force-dynamic";
// Keep loading local to the listing. A segment-wide loading boundary would
// stream HTTP 200 before /models/[code] can determine that a model is absent.
export default function ModelsPage(){return <Suspense fallback={<PublicLoading/>}><PublicCatalogue/></Suspense>;}
