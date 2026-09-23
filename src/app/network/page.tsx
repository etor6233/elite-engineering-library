import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import {createHash} from "node:crypto";
import {notFound,redirect} from "next/navigation";
import type {Route} from "next";
import {readSession,allowed} from "@/platform/auth/session";
import {loadBusinessConfig} from "@/platform/config/load";
import {NetworkWorkspace} from "@/components/network-workspace";
export default async function NetworkPage(){
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

 if((await loadBusinessConfig()).features.network_portal!==true)notFound();
 const s=await readSession();if(!s)redirect("/api/auth/login?return_to=/network"as Route);
 if(!allowed(s,"network:admin")&&!allowed(s,"franchise:write"))return <><h1>{t("p0116")}</h1><p>{t("p0002")}</p></>;
 const scope=createHash("sha256").update(JSON.stringify([s.tenantId,s.subject])).digest("hex");
 return <><h1 className="pageTitle">{t("p0116")}</h1><p>{t("p0117")}</p><NetworkWorkspace scope={scope} subject={s.subject} permissions={s.permissions} organizations={s.organizations}/></>;
}
