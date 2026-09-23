import { createHash } from "node:crypto";
import { redirect,notFound } from "next/navigation";
import type { Route } from "next";
import { readSession,allowed } from "@/platform/auth/session";
import { loadBusinessConfig } from "@/platform/config/load";
import { loadPrivateI18n } from "@/platform/i18n/load-private-i18n";
import { AIActivityWorkspace } from "@/components/ai-activity-workspace";
import { aiMessages } from "@/platform/intelligence/activity";
export const dynamic="force-dynamic";
export default async function AIActivityPage(){
 if((await loadBusinessConfig()).features.ai_activity!==true)notFound();
 const session=await readSession();if(!session)redirect("/api/auth/login?return_to=/intelligence/activity" as Route);
 const {locale}=await loadPrivateI18n();const language=locale.language==="en"?"en":"es";const t=aiMessages[language];
 if(!allowed(session,"whatsapp:approve"))return <section><h1>{t.title}</h1><p role="alert">{t.denied}</p></section>;
 const scope=createHash("sha256").update(JSON.stringify([session.tenantId,session.subject,[...session.organizations].sort()])).digest("hex");
 return <AIActivityWorkspace scope={scope} language={language} locale={locale.locale} timeZone={locale.timeZone}/>;
}
