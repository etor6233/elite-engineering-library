import { createHash } from "node:crypto";
import { redirect,notFound } from "next/navigation";
import type { Route } from "next";
import { readSession,allowed } from "@/platform/auth/session";
import { loadBusinessConfig } from "@/platform/config/load";
import { loadPrivateI18n } from "@/platform/i18n/load-private-i18n";
import { PublishingWorkspace } from "@/components/publishing-workspace";
import { publishingMessages } from "@/platform/publishing/messages";
export const dynamic="force-dynamic";
export default async function PublishingPage(){
 if((await loadBusinessConfig()).features.social_publishing!==true)notFound();
 const session=await readSession();if(!session)redirect("/api/auth/login?return_to=/admin/publishing" as Route);
 const {locale}=await loadPrivateI18n();const language=locale.language==="en"?"en":"es";const t=publishingMessages[language];
 if(!allowed(session,"social:read"))return <section><h1>{t.title}</h1><p role="alert">{t.denied}</p></section>;
 const scope=createHash("sha256").update(JSON.stringify([session.tenantId,session.subject,[...session.organizations].sort()])).digest("hex");
 return <PublishingWorkspace key={scope} scope={scope} language={language} locale={locale.locale} timeZone={locale.timeZone} canRequest={allowed(session,"social:request")} canApprove={allowed(session,"social:approve")} canReconcile={allowed(session,"social:reconcile")}/>;
}
