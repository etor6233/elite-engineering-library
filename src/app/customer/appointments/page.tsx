import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import { publicAppointmentTime } from "@/platform/i18n/public-catalog";
import Link from "next/link";
import { redirect } from "next/navigation";
import type { Route } from "next";
import { allowed, readSession } from "@/platform/auth/session";
import { protectedGet, type CustomerJourney } from "@/platform/backend/protected-client";
import { CustomerAppointmentActions } from "@/components/customer-appointment-actions";

export default async function CustomerAppointmentsPage() {
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

  const session = await readSession();
  if (!session) redirect("/api/auth/login?return_to=/customer/appointments" as Route);
  if (!allowed(session, "customer:self")) return <><h1 className="pageTitle">{t("p0006")}</h1><div className="notice">{t("p0007")} <code>customer:self</code>{t("p0008")}</div></>;
  const organization = session.organizations[0]!;
  const journey = await protectedGet<CustomerJourney>(session, "/v1/customer/journey", { organization_id: organization });
  const locale = privateLocale;
  const appointments = journey.appointments.map((item) => ({ ...item, startsAtLabel: `${publicAppointmentTime(locale, item.starts_at)} (${locale.timeZone})` }));
  return <>
    <div className="eyebrow">{t("p0034")} {organization}</div><h1 className="pageTitle">{t("p0035")}</h1>
    <p className="lede">{t("p0036")}</p>
    <p><Link href={"/customer/quotes" as Route}>{t("p0037")}</Link></p>
    <p><Link href={"/customer/handovers" as Route}>{t("p0038")}</Link></p>
    <CustomerAppointmentActions organizationId={organization} appointments={appointments}/>
  </>;
}
