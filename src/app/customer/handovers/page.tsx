import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import { publicAppointmentTime } from "@/platform/i18n/public-catalog";
import { OperationalGuide } from "@/components/operational-guide";
import { OPERATIONAL_GUIDES } from "@/platform/help/content";
import Link from "next/link";
import { redirect } from "next/navigation";
import type { Route } from "next";
import { allowed, readSession } from "@/platform/auth/session";
import { protectedGet, type CustomerJourney } from "@/platform/backend/protected-client";
import { CustomerHandoverActions } from "@/components/customer-handover-actions";

type DeliveryHandover = CustomerJourney["handovers"][number] & { checklist_id?: string; checklist_version?: number; checklist_title?: string; checklist_items: { id: string; ordinal: number; prompt: string; response_type: string; required: boolean }[]; checklist_completed_at?: string };
type DeliveryException = { id: string; handover_id: string; reason_code: string; details: string; state: string; version: number; resolution_action?: string; successor_handover_id?: string; return_authorization_id?: string };
type DeliveryJourney = Omit<CustomerJourney, "handovers"> & { handovers: DeliveryHandover[]; delivery_exceptions: DeliveryException[] };

export default async function CustomerHandoversPage() {
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

  const session = await readSession();
  if (!session) redirect("/api/auth/login?return_to=/customer/handovers" as Route);
  if (!allowed(session, "customer:self")) return <><h1 className="pageTitle">{t("p0006")}</h1><div className="notice">{t("p0007")} <code>customer:self</code>{t("p0008")}</div></>;
  const organization = session.organizations[0]!;
  const journey = await protectedGet<DeliveryJourney>(session, "/v1/customer/journey", { organization_id: organization }).catch(() => null);
  const locale = privateLocale;
  const handovers = journey?.handovers.map(item => ({ ...item, ...(item.checklist_completed_at ? { checklistCompletedLabel: `${publicAppointmentTime(locale, item.checklist_completed_at)} (${locale.timeZone})` } : {}), ...(item.customer_accepted_at ? { acceptedLabel: `${publicAppointmentTime(locale, item.customer_accepted_at)} (${locale.timeZone})` } : {}) })) ?? [];
  return <><div className="eyebrow">{t("p0034")} {organization}</div><h1 className="pageTitle">{t("p0039")}</h1><p className="lede">{t("p0040")}</p><p><Link href={"/customer/appointments" as Route}>{t("p0041")}</Link></p>
    {journey ? <>{journey.handovers.length === 0 ? <p>{t("p0042")}</p> : null}<CustomerHandoverActions organizationId={organization} handovers={handovers} exceptions={journey.delivery_exceptions}/></> : <><div className="notice" role="alert">{t("p0043")}</div><p><a href="/customer/handovers">{t("p0044")}</a></p></>}
    <OperationalGuide guide={OPERATIONAL_GUIDES["handover-read-view"]}/>
  </>;
}
