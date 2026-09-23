import Link from "next/link";
import { notFound } from "next/navigation";
import { loadPublishedCatalog } from "@/platform/catalog/load";
import { publicModelMetadata } from "@/platform/seo/public-indexing";
import { loadPublicLocale } from "@/platform/i18n/load-public-locale";
import { modelClass, modelDescription, publicSpecifications, startingPrice, variantPrice } from "@/components/public-web/catalog-presentation";
import { PublicIcon } from "@/components/public-web/public-ui";
import { PublicLeadForm } from "@/components/public-web/public-lead-form";
import { PublicMedia } from "@/components/public-web/public-media";
import { publicText as t } from "@/components/public-web/messages";
import s from "@/components/public-web/public-web.module.css";

export const dynamic="force-dynamic";
async function selected(code:string){const catalog=await loadPublishedCatalog();const model=catalog?.models.find(item=>item.code===code);if(!catalog||!model)notFound();return{catalog,model};}
export async function generateMetadata({params}:{params:Promise<{code:string}>}){const{model}=await selected((await params).code);return publicModelMetadata(model.code,model.displayName,model.canonical_url);}
export default async function ModelPage({params}:{params:Promise<{code:string}>}){
  const [{catalog,model},locale]=await Promise.all([selected((await params).code),loadPublicLocale()]);
  const variants=catalog.variants.filter(item=>item.model_id===model.id);
  const price=startingPrice(catalog,model.id,locale.locale);
  const description=modelDescription(model.specification);
  const specifications=publicSpecifications(model.specification,locale.locale);
  return <div className={s.content} lang={locale.locale}>
    <Link href="/models" className={s.backLink}><PublicIcon name="arrow"/>{t(locale.locale,"backModels")}</Link>
    <section className={s.detailHero}><div className={s.detailIntro}><p className={s.eyebrow}>{modelClass(model.vehicleClass,locale.locale)}</p><h1>{model.displayName}</h1>{description&&<p>{description}</p>}{price&&<div className={s.price}><p>{price.prefix} {price.amount}</p><span>{price.tax}</span></div>}<div className={s.actions}><a href="#consulta" className={s.primary}>{t(locale.locale,"inquire")}<PublicIcon name="arrow"/></a><a href="#versions" className={s.secondary}>{t(locale.locale,"variants")}</a></div></div>
      <figure className={s.detailImage}><PublicMedia src={`/api/public/catalog/media/${model.media.sha256}`} alt={model.displayName} width={model.media.width} height={model.media.height} locale={locale.locale} priority/></figure>
    </section>
    <section className={s.detailsSection} id="versions"><h2>{t(locale.locale,"variants")}</h2><div className={s.variants}>{variants.map(variant=>{const value=variantPrice(variant,catalog.currency,locale.locale);const battery=modelDescription(variant.battery_specification);return<article className={s.variant} key={variant.id}><h3>{variant.display_name}</h3><div className={s.price}><p>{value.amount}</p><span>{value.tax}</span></div>{battery&&<p>{battery}</p>}</article>;})}</div></section>
    {specifications.length>0&&<section className={s.detailsSection}><h2>{t(locale.locale,"specifications")}</h2><dl className={s.specifications}>{specifications.map(item=><div key={item.label}><dt>{item.label}</dt><dd>{item.value}</dd></div>)}</dl></section>}
    <section className={s.inquiry} id="consulta"><div className={s.inquiryIntro}><p className={s.eyebrow}>{model.displayName}</p><h2>{t(locale.locale,"talk")}</h2><p>{t(locale.locale,"inquiryBody")}</p></div><PublicLeadForm modelId={model.id} locale={locale.locale}/></section>
  </div>;
}
