import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import Link from "next/link";
import type { Route } from "next";
import { notFound, redirect } from "next/navigation";
import { loadBusinessConfig } from "@/platform/config/load";
import { readSession, allowed } from "@/platform/auth/session";
import { protectedGet } from "@/platform/backend/protected-client";
import { summarySchema, surveyID } from "@/platform/surveys/contract";
export default async function SurveySummaryPage({ searchParams }: { searchParams: Promise<Record<string, string | string[] | undefined>> }) {
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

 if ((await loadBusinessConfig()).features.customer_surveys !== true) notFound();
 const query = await searchParams;
 if (Object.keys(query).length !== 2 || !surveyID.safeParse(query.surveyId).success || !surveyID.safeParse(query.organizationId).success) return <section lang="es"><h1 className="pageTitle">{t("p0021")}</h1><p>{t("p0022")}</p></section>;
 const id = query.surveyId as string, organization = query.organizationId as string;
 const href = "/admin/surveys?surveyId=" + encodeURIComponent(id) + "&organizationId=" + encodeURIComponent(organization);
 const session = await readSession(); if (!session) redirect(("/api/auth/login?return_to=" + encodeURIComponent(href)) as Route);
 if (!allowed(session, "surveys:read") || (!session.organizations.includes(organization) && !session.organizations.includes("*"))) return <section lang="es"><h1 className="pageTitle">{t("p0006")}</h1><p>{t("p0023")}</p></section>;
 try {
  const value = summarySchema.parse(await protectedGet<unknown>(session, "/v1/admin/surveys/" + encodeURIComponent(id) + "/summary", { organization_id: organization }));
  return <section lang="es"><h1 className="pageTitle">{t("p0024")}</h1><p>{t("p0025")} {value.responses}</p>{value.available ? <p>{t("p0026")} <strong data-testid="survey-nps">{value.nps!.toFixed(2)}</strong></p> : <p>{t("p0027")}</p>}<p>{t("p0028")}</p><details><summary>{t("p0029")}</summary><p>{t("p0030")}</p><p>{t("p0031")}</p></details></section>;
 } catch { return <section lang="es"><h1 className="pageTitle">{t("p0032")}</h1><Link href={href as Route}>{t("p0033")}</Link></section>; }
}

