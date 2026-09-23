import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import { createHash } from "node:crypto";
import { redirect,notFound } from "next/navigation";
import type { Route } from "next";
import { allowed,readSession } from "@/platform/auth/session";
import { protectedGet } from "@/platform/backend/protected-client";
import { loadBusinessConfig } from "@/platform/config/load";
import { courseView,assessmentSchema } from "@/platform/training/contract";
import { TrainingWorkspace } from "@/components/training-workspace";
export default async function TrainingPage(){
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

 if((await loadBusinessConfig()).features.training_portal!==true)notFound();
 const session=await readSession();if(!session)redirect("/api/auth/login?return_to=/guide/training" as Route);
 const learn=allowed(session,"training:learn"),review=allowed(session,"training:review");
 if(!learn&&!review)return <><h1 className="pageTitle">{t("p0105")}</h1><p>{t("p0106")}</p></>;
 try{
  const [courses,assessments]=await Promise.all([protectedGet(session,"/v1/training/courses",{}),protectedGet(session,"/v1/training/assessments",{})]);
  const views=courseView.array().max(8).parse(courses),records=assessmentSchema.array().max(50).parse(assessments);
  if(records.some(x=>!session.organizations.includes(x.payload.organization_id)||(!review&&x.payload.learner_subject!==session.subject)))throw new Error("SCOPE");
  const scope=createHash("sha256").update(JSON.stringify([session.tenantId,session.subject,session.organizations])).digest("hex");
  return <><h1 className="pageTitle">{t("p0107")}</h1><p>{t("p0108")}</p><TrainingWorkspace courses={views} assessments={records} canLearn={learn} canReview={review} scope={scope}/></>
 }catch{return <><h1 className="pageTitle">{t("p0105")}</h1><p role="alert">{t("p0109")}</p><a className="button" href="/guide/training">{t("p0033")}</a></>}
}
