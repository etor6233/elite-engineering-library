import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import Link from "next/link";
import type { Route } from "next";
import { redirect, notFound } from "next/navigation";
import { readSession, allowed } from "@/platform/auth/session";
import { loadBusinessConfig } from "@/platform/config/load";
import { protectedGet } from "@/platform/backend/protected-client";
import { BackendProblem } from "@/platform/backend/public-client";
import { definitionSchema, answerSchema, surveyID, type SurveyAnswer } from "@/platform/surveys/contract";
import { CustomerSurvey } from "@/components/customer-survey";

export default async function CustomerSurveyPage({ searchParams }: { searchParams: Promise<Record<string, string | string[] | undefined>> }) {
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

 if ((await loadBusinessConfig()).features.customer_surveys !== true) notFound();
 const query = await searchParams; const id = query.surveyId, org = query.organizationId;
 if (Object.keys(query).length !== 2 || !surveyID.safeParse(id).success || !surveyID.safeParse(org).success) return <section lang="es"><h1 className="pageTitle">{t("p0063")}</h1><p>{t("p0064")}</p></section>;
 const survey = id as string, organization = org as string;
 const href = "/customer/surveys?surveyId=" + encodeURIComponent(survey) + "&organizationId=" + encodeURIComponent(organization);
 const session = await readSession(); if (!session) redirect(("/api/auth/login?return_to=" + encodeURIComponent(href)) as Route);
 if (!allowed(session, "customer:self") || (!session.organizations.includes(organization) && !session.organizations.includes("*"))) return <section lang="es"><h1 className="pageTitle">{t("p0006")}</h1><p>{t("p0065")}</p></section>;
 try {
  const path = "/v1/customer/surveys/" + encodeURIComponent(survey);
  const definition = definitionSchema.parse(await protectedGet<unknown>(session, path, { organization_id: organization }));
  if (definition.id !== survey) throw new Error("survey mismatch");
  let answer: SurveyAnswer | null = null;
  try { answer = answerSchema.parse(await protectedGet<unknown>(session, path + "/response", { organization_id: organization })); }
  catch (error) { if (!(error instanceof BackendProblem && error.status === 404 && error.code === "SURVEY_NOT_AVAILABLE")) throw error; }
  return <CustomerSurvey organization={organization} definition={definition} initialAnswer={answer} />;
 } catch {
  return <section lang="es"><h1 className="pageTitle">{t("p0066")}</h1><p>{t("p0067")}</p><Link href={href as Route}>{t("p0033")}</Link></section>;
 }
}

