import Link from "next/link";
import { headers } from "next/headers";
import { loadBusinessConfig } from "@/platform/config/load";
import { loadPublicLocale } from "@/platform/i18n/load-public-locale";
import { publicPageMetadata, websiteJsonLd } from "@/platform/seo/public-indexing";
import { PublicIcon, PublicEmpty, PublicVisit } from "@/components/public-web/public-ui";
import { PublicProductCard, readPublicModels } from "@/components/public-web/public-catalogue";
import { publicText as t } from "@/components/public-web/messages";
import s from "@/components/public-web/public-web.module.css";

export function generateMetadata(){return publicPageMetadata("/");}
export default async function HomePage(){
  const [config,locale,result]=await Promise.all([loadBusinessConfig(),loadPublicLocale(),readPublicModels()]);
  const structuredData=websiteJsonLd(config.business.name,locale.locale);
  const nonce=(await headers()).get("x-nonce")??undefined;
  return <div lang={locale.locale}>
    {structuredData&&<script type="application/ld+json" nonce={nonce} dangerouslySetInnerHTML={{__html:structuredData}}/>}
    <section className={s.hero}>
      <img className={s.heroImage} src="/public-web/urban-hero.jpg" alt={t(locale.locale,"heroAlt")} width="1672" height="941" fetchPriority="high"/>
      <div className={s.heroContent}><p className={s.eyebrow}>{t(locale.locale,"eyebrow")}</p><h1>{t(locale.locale,"homeTitle")}</h1><p>{t(locale.locale,"homeBody")}</p><div className={s.actions}><Link href="/models" className={s.primary}>{t(locale.locale,"explore")}<PublicIcon name="arrow"/></Link><Link href="/locations" className={s.secondary}>{t(locale.locale,"findUs")}</Link></div></div>
      <span className={s.heroCaption}>{t(locale.locale,"conceptual")}</span>
    </section>
    <div className={s.content}><section className={s.collection}><div className={s.sectionTitle}><div><p className={s.eyebrow}>{t(locale.locale,"collection")}</p><h2>{t(locale.locale,"catalogueTitle")}</h2></div><Link href="/models" className={s.textLink}>{t(locale.locale,"explore")}<PublicIcon name="arrow"/></Link></div>
      {result.failed?<PublicEmpty title={t(locale.locale,"catalogueError")} locale={locale.locale} retryHref="/">{t(locale.locale,"catalogueErrorBody")}</PublicEmpty>:result.models.length?<div className={s.productGrid}>{result.models.slice(0,2).map(model=><PublicProductCard model={model} catalog={result.catalog} locale={locale.locale} key={model.id}/>)}</div>:<PublicEmpty title={t(locale.locale,"noModels")} locale={locale.locale}>{t(locale.locale,"noModelsBody")}</PublicEmpty>}
    </section><PublicVisit locale={locale.locale}/></div>
  </div>;
}
