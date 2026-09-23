import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import {createHash} from "node:crypto";
import {redirect,notFound} from "next/navigation";
import type {Route} from "next";
import {readSession,allowed} from "@/platform/auth/session";
import {loadBusinessConfig} from "@/platform/config/load";
import {warrantyID,warrantySurface,warrantyReadPermission} from "@/platform/warranty/contract";
import {WarrantyWorkspace} from "@/components/warranty-workspace";
export default async function WarrantyPage({searchParams}:{searchParams:Promise<Record<string,string|string[]|undefined>>}){
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

 if((await loadBusinessConfig()).features.warranty_portal!==true)notFound();const s=await readSession();if(!s)redirect("/api/auth/login?return_to=/warranty" as Route);
 const available=(["franchise","customer","factory"] as const).filter(role=>allowed(s,warrantyReadPermission(role)));if(!available.length)return <><h1>{t("p0123")}</h1><p>{t("p0002")}</p></>;
 const q=await searchParams,requested=warrantySurface.safeParse(q.surface),surface=requested.success&&available.includes(requested.data)?requested.data:available[0]!;
 const organization=typeof q.organization==="string"&&s.organizations.includes(q.organization)?q.organization:s.organizations[0];if(!organization)return <p>{t("p0003")}</p>;
 const initial=(key:string)=>typeof q[key]==="string"&&warrantyID.safeParse(q[key]).success?q[key] as string:"";
 const scope=createHash("sha256").update(JSON.stringify([s.tenantId,organization,s.subject,surface])).digest("hex");
 return <><h1 className="pageTitle">{t("p0124")}</h1><p>{t("p0005")} {organization}{t("p0120")} {surface==="franchise"?t("p0010"):surface==="customer"?t("p0125"):t("p0121")}{t("p0008")}</p><WarrantyWorkspace key={scope} scope={scope} organization={organization} tenant={s.tenantId} subject={s.subject} permissions={s.permissions} surface={surface} initialClaim={initial("case")} initialQuote={initial("quote")} initialHandover={initial("handover")}/></>;
}
