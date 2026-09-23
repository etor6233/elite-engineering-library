import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import { redirect } from "next/navigation";
import type { Route } from "next";
import { allowed, readSession } from "@/platform/auth/session";
import { protectedGet } from "@/platform/backend/protected-client";
import { responseText, type ReplyView } from "@/platform/notifications/reply-contract";
import { WhatsAppReplyReview } from "@/components/whatsapp-reply-review";

export default async function WhatsAppReplyPage(){
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

 const session=await readSession();if(!session)redirect("/api/auth/login?return_to=/franchise/whatsapp" as Route);
 if(!allowed(session,"whatsapp:approve"))return <><h1 className="pageTitle">{t("p0096")}</h1><p>{t("p0097")}</p></>;
 try{
  const replies=await protectedGet<ReplyView[]>(session,"/v1/franchise/whatsapp/replies",{});
  if(!Array.isArray(replies)||replies.length>50||replies.some(v=>!session.organizations.includes(v.context.organization_id)))throw new Error("INVALID_SCOPE");
  const rows=replies.map(value=>({...value,body:responseText(value),windowLabel:new Intl.DateTimeFormat(privateLocale.locale,{dateStyle:"medium",timeStyle:"short",timeZone:privateLocale.timeZone}).format(new Date(value.context.expires_at))+" ("+privateLocale.timeZone+")"}));
  return <><h1 className="pageTitle">{t("p0096")}</h1><p>{t("p0098")}</p><WhatsAppReplyReview replies={rows} canSend={allowed(session,"whatsapp:send")}/></>;
 }catch{return <><h1 className="pageTitle">{t("p0096")}</h1><p role="alert">{t("p0099")}</p><a className="button" href="/franchise/whatsapp">{t("p0095")}</a></>}
}
