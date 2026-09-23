import { publicMessage as message } from "@/platform/i18n/public-catalog";
import { loadPublicLocale } from "@/platform/i18n/load-public-locale";
import { publicCount } from "@/platform/i18n/public-catalog";
import { listModels } from "@/platform/backend/public-client";
import { LeadForm } from "@/app/connected/lead-form";
import { publicPageMetadata } from "@/platform/seo/public-indexing";
import Link from "next/link";
import { loadPublishedCatalog } from "@/platform/catalog/load";

export async function generateMetadata() { const locale = await loadPublicLocale(); return publicPageMetadata("/models", message(locale.locale, "models.short")); }
export const dynamic = "force-dynamic";

export default async function ModelsPage() {
  const [catalog, locale] = await Promise.all([loadPublishedCatalog(), loadPublicLocale()]);
  // An enabled publication profile never falls back to mutable draft data.
  const models = catalog?.models ?? (process.env.CATALOG_RELEASE_ENABLED === "true" ? [] : await listModels());
  return (
    <section className="publicCatalogue" lang={locale.locale}>
      <div className="eyebrow">{message(locale.locale, "models.eyebrow")}</div>
      <h1 className="pageTitle">{message(locale.locale, "models.title")}</h1>
      <p className="lede">{message(locale.locale, "models.description")}</p>
      <p>{publicCount(locale.locale, "models", models.length)}</p>
      <div className="grid">
        {models.map((model) => (
          <article className="card" key={model.id}>
            <h2>{catalog ? <Link href={`/models/${model.code}`}>{model.displayName}</Link> : model.displayName}</h2>
            <p>{({bicycle:locale.locale.startsWith("es")?"Bicicleta":"Bicycle",motorcycle:locale.locale.startsWith("es")?"Motocicleta":"Motorcycle",car:locale.locale.startsWith("es")?"Automóvil":"Car"} as Record<string,string>)[model.vehicleClass]??model.vehicleClass}</p>
            <details className="modelInquiry"><summary>{message(locale.locale,"lead.submit")}</summary><LeadForm modelId={model.id} locale={locale.locale} /></details>
          </article>
        ))}
      </div>
    </section>
  );
}
