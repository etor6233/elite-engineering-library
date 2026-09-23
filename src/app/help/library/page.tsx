import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import{createHash}from"node:crypto";
import{notFound,redirect}from"next/navigation";
import type{Route}from"next";
import{readSession,allowed}from"@/platform/auth/session";
import{loadBusinessConfig}from"@/platform/config/load";
import{HelpCMSWorkspace}from"@/components/help-cms-workspace";
export default async function HelpLibraryPage(){
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

 if((await loadBusinessConfig()).features.help_cms!==true)notFound();const s=await readSession();if(!s)redirect("/api/auth/login?return_to=/help/library"as Route);
 if(!["help:read","help:write","help:publish"].some(p=>allowed(s,p)))return <><h1>{t("p0115")}</h1><p>{t("p0002")}</p></>;
 const scope=createHash("sha256").update(JSON.stringify([s.tenantId,s.subject])).digest("hex");
 return <><h1 className="pageTitle">{t("p0115")}</h1><HelpCMSWorkspace scope={scope} subject={s.subject} permissions={s.permissions} organizations={s.organizations}/></>
}
