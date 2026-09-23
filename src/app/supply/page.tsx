import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import {createHash} from "node:crypto";
import {redirect,notFound} from "next/navigation";
import type {Route} from "next";
import {allowed,readSession} from "@/platform/auth/session";
import {loadBusinessConfig} from "@/platform/config/load";
import {supplyID} from "@/platform/supply/contract";
import {SupplyWorkspace} from "@/components/supply-workspace";
export default async function SupplyPage({searchParams}:{searchParams:Promise<Record<string,string|string[]|undefined>>}){
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

 if((await loadBusinessConfig()).features.supply_portal!==true)notFound();
 const s=await readSession();if(!s)redirect("/api/auth/login?return_to=/supply" as Route);
 if(!allowed(s,"supply:read")&&!allowed(s,"supply:factory-read"))return <><h1>{t("p0118")}</h1><p>{t("p0002")}</p></>;
 const q=await searchParams,organization=typeof q.organization==="string"&&s.organizations.includes(q.organization)?q.organization:s.organizations[0];
 if(!organization)return <p>{t("p0003")}</p>;
 const surface=q.surface==="factory"&&allowed(s,"supply:factory-read")?"factory":allowed(s,"supply:read")?"franchise":"factory";
 const initial=typeof q.order==="string"&&supplyID.safeParse(q.order).success?q.order:"";
 const scope=createHash("sha256").update(JSON.stringify([s.tenantId,organization,s.subject,surface])).digest("hex");
 return <><h1 className="pageTitle">{t("p0119")}</h1><p>{t("p0005")} {organization}{t("p0120")} {surface==="factory"?t("p0121"):t("p0122")}{t("p0008")}</p><SupplyWorkspace key={scope} scope={scope} organization={organization} subject={s.subject} permissions={s.permissions} surface={surface} initial={initial}/></>;
}
