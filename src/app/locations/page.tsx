import Link from "next/link";
import type {Metadata} from "next";
import {loadPublicLocale} from "@/platform/i18n/load-public-locale";
import {listLocations,listAppointmentSlots,type AppointmentKind,type AppointmentSlot} from "@/platform/backend/public-client";
import {readPublicIndexing} from "@/platform/seo/public-indexing";
import {AppointmentForm} from "./appointment-form";
import {publicText as t} from "@/components/public-web/messages";
import {PublicIcon,PublicEmpty} from "@/components/public-web/public-ui";
import s from "@/components/public-web/public-web.module.css";
export async function generateMetadata():Promise<Metadata>{
  const locale=await loadPublicLocale();
  const {origin}=readPublicIndexing();
  // Locations is outside the approved indexing allowlist. Keep it noindex;
  // receipt/query references never enter the canonical URL or title.
  return {title:t(locale.locale,"locations"),robots:{index:false,follow:false},...(origin?{alternates:{canonical:new URL("/locations",origin).href}}:{})};
}
export default async function LocationsPage({searchParams}:{searchParams:Promise<{lead_id?:string;model_id?:string}>}){
  const [locale,query]=await Promise.all([loadPublicLocale(),searchParams]);
  const locationResult=await listLocations().then(items=>({items,failed:false}),()=>({items:[],failed:true}));
  let slots:AppointmentSlot[]=[];let slotsFailed=false;
  // No appointment calls until a lead receipt supplies context, as this page
  // cannot request a valid appointment without that existing owner reference.
  if(query.lead_id){
    const from=new Date(Date.now()+60000).toISOString();const to=new Date(Date.now()+30*24*60*60000).toISOString();
    const kinds:AppointmentKind[]=["consultation","test-drive","delivery","service"];
    try{slots=(await Promise.all(kinds.map(kind=>listAppointmentSlots(kind,from,to)))).flat().sort((a,b)=>a.starts_at.localeCompare(b.starts_at));}catch{slotsFailed=true;}
  }
  const retry=new URLSearchParams();if(query.lead_id)retry.set("lead_id",query.lead_id);if(query.model_id)retry.set("model_id",query.model_id);
  return <div className={s.content} lang={locale.locale}><section className={s.pageHeading}><p className={s.eyebrow}>{t(locale.locale,"locations")}</p><h1>{t(locale.locale,"visit")}</h1><p>{t(locale.locale,"visitBody")}</p></section>
    {locationResult.failed?<PublicEmpty title={t(locale.locale,"locationError")} locale={locale.locale} retryHref={`/locations${retry.size?"?"+retry.toString():""}`}/>:locationResult.items.length?<section className={s.locations} aria-label={t(locale.locale,"locations")}>{locationResult.items.map(location=><article className={s.location} key={location.organization_id}><PublicIcon name="pin"/><h2>{location.name}</h2><p>{[location.city,location.region,location.country].filter(Boolean).join(", ")}</p><div className={s.contactLinks}>{location.contact_phone&&<a href={`tel:${location.contact_phone}`}><PublicIcon name="phone"/>{t(locale.locale,"call")}</a>}{location.contact_email&&<a href={`mailto:${location.contact_email}`}><PublicIcon name="mail"/>{t(locale.locale,"email")}</a>}</div></article>)}</section>:<div className={s.empty} data-public-state="empty"><h2>{t(locale.locale,"noLocations")}</h2><p>{t(locale.locale,"noLocationsBody")}</p><Link href="/models" className={s.secondary}>{t(locale.locale,"explore")}</Link></div>}
    {query.lead_id?<section className={s.booking}><h2>{t(locale.locale,"book")}</h2><p>{t(locale.locale,"bookBody")}</p>{slotsFailed?<div className={s.uncertain} role="status" data-public-state="error"><p>{t(locale.locale,"slotsError")}</p><a className={s.secondary} href={`/locations?${retry.toString()}`}>{t(locale.locale,"retry")}</a></div>:<AppointmentForm leadId={query.lead_id} slots={slots} locale={locale} {...(query.model_id?{modelId:query.model_id}:{})}/>}</section>:<section className={s.startInquiry}><h2>{t(locale.locale,"startInquiry")}</h2><p>{t(locale.locale,"startInquiryBody")}</p><Link href="/models" className={s.secondary}>{t(locale.locale,"explore")}<PublicIcon name="arrow"/></Link></section>}
  </div>;
}
