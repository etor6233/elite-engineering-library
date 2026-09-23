import Link from "next/link";
import { loadPublicLocale } from "@/platform/i18n/load-public-locale";
import { loadPublishedCatalog } from "@/platform/catalog/load";
import { listModels, type PublicModel } from "@/platform/backend/public-client";
import type { PublishedCatalog } from "@/platform/catalog/schema";
import { publicText as t } from "./messages";
import { modelClass, modelDescription, startingPrice } from "./catalog-presentation";
import { PublicIcon, PublicEmpty, PublicVisit } from "./public-ui";
import { PublicLeadForm } from "./public-lead-form";
import { PublicMedia } from "./public-media";
import s from "./public-web.module.css";

export async function readPublicModels():Promise<{catalog:PublishedCatalog|null;models:PublicModel[];failed:boolean}>{
  try{const catalog=await loadPublishedCatalog();return{catalog,models:catalog?.models??(process.env.CATALOG_RELEASE_ENABLED==="true"?[]:await listModels()),failed:false};}
  catch{return{catalog:null,models:[],failed:true};}
}
export function PublicProductCard({model,catalog,locale}:{model:PublicModel;catalog:PublishedCatalog|null;locale:string}){
  const published=catalog?.models.find(item=>item.id===model.id);
  const price=catalog?startingPrice(catalog,model.id,locale):null;
  const description=modelDescription(model.specification);
  return <article className={s.productCard}>
    <div className={s.productMedia}>{published?<PublicMedia src={`/api/public/catalog/media/${published.media.sha256}`} alt={model.displayName} width={published.media.width} height={published.media.height} locale={locale}/>:<span className={s.noMedia}>{t(locale,"noMedia")}</span>}</div>
    <div className={s.productInfo}><div><p className={s.eyebrow}>{modelClass(model.vehicleClass,locale)}</p><h2>{model.displayName}</h2>{description&&<p className={s.cardDescription}>{description}</p>}</div>
      {price&&<div className={s.price}><p>{price.prefix} {price.amount}</p><span>{price.tax}</span></div>}
      {published?<Link className={s.cardAction} aria-label={`${t(locale,"detail")}: ${model.displayName}`} href={`/models/${model.code}`}>{t(locale,"detail")}<PublicIcon name="arrow"/></Link>:<details className={s.inquiryDisclosure}><summary>{t(locale,"inquire")}</summary><PublicLeadForm modelId={model.id} locale={locale}/></details>}
    </div>
  </article>;
}
export async function PublicCatalogue(){const [result,locale]=await Promise.all([readPublicModels(),loadPublicLocale()]);
  return <div className={s.content} lang={locale.locale}><section className={s.pageHeading}><p className={s.eyebrow}>{t(locale.locale,"collection")}</p><h1>{t(locale.locale,"catalogueTitle")}</h1><p>{t(locale.locale,"catalogueBody")}</p></section>
    {result.failed?<PublicEmpty title={t(locale.locale,"catalogueError")} locale={locale.locale} retryHref="/models">{t(locale.locale,"catalogueErrorBody")}</PublicEmpty>:result.models.length?<><p className={s.count}>{result.models.length} {t(locale.locale,result.models.length===1?"modelsOne":"modelsMany")}</p><section className={s.productGrid} aria-label={t(locale.locale,"models")}>{result.models.map(model=><PublicProductCard key={model.id} model={model} catalog={result.catalog} locale={locale.locale}/>)}</section></>:<PublicEmpty title={t(locale.locale,"noModels")} locale={locale.locale}>{t(locale.locale,"noModelsBody")}</PublicEmpty>}
    <PublicVisit locale={locale.locale}/></div>;
}
