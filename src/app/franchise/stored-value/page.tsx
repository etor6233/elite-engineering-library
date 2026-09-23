import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import { redirect } from "next/navigation";
import type { Route } from "next";
import { createHash } from "node:crypto";
import { z } from "zod";
import { readSession,allowed } from "@/platform/auth/session";
import { protectedGet } from "@/platform/backend/protected-client";
import { storedProfile,storedID } from "@/platform/stored-value/contracts";
import { StoredValueOperations } from "@/components/stored-value-operations";
const configuration=z.object({programs:z.array(z.object({calculation:z.object({id:storedID,kind:z.enum(["gift_card","loyalty"])})})).min(1).max(8)});
export default async function StoredValuePage({searchParams}:{searchParams:Promise<{organization?:string}>}){
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

 const session=await readSession();if(!session)redirect("/api/auth/login?return_to=/franchise/stored-value" as Route);
 if(!allowed(session,"stored_value:read"))return <><h1>{t("p0088")}</h1><p>{t("p0089")}</p></>;
 const organization=(await searchParams).organization??session.organizations[0];
 if(!organization||!session.organizations.includes(organization))return <><h1>{t("p0088")}</h1><p>{t("p0090")}</p></>;
 const choose=<form action="/franchise/stored-value" method="get"><label>{t("p0091")}<select name="organization" defaultValue={organization}>{session.organizations.map(o=><option key={o}>{o}</option>)}</select></label><button className="button">{t("p0092")}</button></form>;
 try{
  const p=storedProfile.parse(await protectedGet<unknown>(session,"/v1/franchise/stored-value/profile",{organization_id:organization}));
  if(p.organization_id!==organization||createHash("sha256").update(p.profile_json,"utf8").digest("hex")!==p.profile_sha256)throw new Error("profile mismatch");
  const programs=configuration.parse(JSON.parse(p.profile_json)).programs.map(p=>p.calculation);
  return <><h1 className="pageTitle">{t("p0088")}</h1><p>{t("p0093")}</p>{choose}<StoredValueOperations organization={organization} profileHash={p.profile_sha256} programs={programs} permissions={session.permissions} subject={session.subject} scope={createHash("sha256").update(JSON.stringify([session.tenantId,session.subject,organization])).digest("hex")}/></>;
 }catch{return <><h1>{t("p0088")}</h1>{choose}<p role="alert">{t("p0094")}</p><a className="button" href={(`/franchise/stored-value?organization=${encodeURIComponent(organization)}`) as Route}>{t("p0095")}</a></>}
}
