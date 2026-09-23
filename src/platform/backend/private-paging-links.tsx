"use client";
import{usePrivateI18n}from"@/platform/i18n/private-provider";
export function PrivatePagingLinks({label,nextHref,firstHref}:{label:string;nextHref:string|undefined;firstHref:string|undefined}){
 const{t}=usePrivateI18n();if(!nextHref&&!firstHref)return null;
 return <nav aria-label={label}>{firstHref?<p><a href={firstHref}>{t("p1197")}</a></p>:null}{nextHref?<p><a href={nextHref}>{t("p1198")}</a></p>:null}</nav>;
}
