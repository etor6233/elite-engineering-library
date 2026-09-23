import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import {createHash} from "node:crypto";
import {redirect,notFound} from "next/navigation";
import type {Route} from "next";
import {readSession,allowed} from "@/platform/auth/session";
import {loadBusinessConfig} from "@/platform/config/load";
import {catalogID} from "@/platform/catalog/authoring";
import {CatalogAuthoring} from "@/components/catalog-authoring";
export default async function CatalogAuthoringPage({searchParams}:{searchParams:Promise<Record<string,string|string[]|undefined>>}){
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

 if((await loadBusinessConfig()).features.catalog_editor!==true)notFound();
 const session=await readSession();if(!session)redirect("/api/auth/login?return_to=/admin/catalog" as Route);
 if(!allowed(session,"catalog:read"))return <><h1>{t("p0001")}</h1><p>{t("p0002")}</p></>;
 const q=await searchParams,requested=q.organization;
 const organization=typeof requested==="string"&&session.organizations.includes(requested)?requested:session.organizations[0];
 if(!organization)return <><h1>{t("p0001")}</h1><p>{t("p0003")}</p></>;
 const initialDraft=typeof q.draft==="string"&&catalogID.safeParse(q.draft).success?q.draft:"";
 const scope=createHash("sha256").update(JSON.stringify([session.tenantId,organization,session.subject])).digest("hex");
 return <><h1 className="pageTitle">{t("p0004")}</h1><p>{t("p0005")} {organization}</p><CatalogAuthoring key={scope} scope={scope} organization={organization} tenant={session.tenantId} subject={session.subject} permissions={session.permissions} initialDraft={initialDraft}/></>;
}
