import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import type { ReactElement } from "react";
import Link from "next/link";
import { CustomerCheckoutActions } from "@/components/customer-checkout-actions";
import { publicAppointmentTime } from "@/platform/i18n/public-catalog";
import { PortalPaging, portalCursor, portalPageHref, type PortalPageProps } from "@/platform/backend/portal-paging";
import { minorAmountPresentation } from "@/platform/i18n/money";
import { redirect } from "next/navigation";
import type { Route } from "next";
import { allowed, readSession } from "@/platform/auth/session";
import { protectedGet, type CustomerJourney, type OrderSummary, type Page, type ServiceCaseSummary } from "@/platform/backend/protected-client";

// AUTHORED framework signature glue: Next inspects the required props overload;
// the zero-argument overload preserves existing server-component test callers.
export default function CustomerPage(): Promise<ReactElement>;
export default function CustomerPage(props: PortalPageProps): Promise<ReactElement>;
export default async function CustomerPage({ searchParams }: PortalPageProps = {}) {
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

  const session = await readSession();
  if (!session) redirect("/api/auth/login?return_to=/customer" as Route);
  if (!allowed(session, "customer:self")) return <><h1 className="pageTitle">{t("p0006")}</h1><div className="notice">{t("p0007")} <code>customer:self</code>{t("p0008")}</div></>;
  const organization = session.organizations[0]!;
  const query = await searchParams ?? {};
  const ordersAfter = portalCursor(query.orders_after), casesAfter = portalCursor(query.cases_after);
  const [orders, cases, journey, locale] = await Promise.all([
    protectedGet<Page<OrderSummary>>(session, "/v1/customer/orders", { organization_id: organization, limit: "25", after: ordersAfter }),
    protectedGet<Page<ServiceCaseSummary>>(session, "/v1/customer/service-cases", { organization_id: organization, limit: "25", after: casesAfter }),
    protectedGet<CustomerJourney>(session, "/v1/customer/journey", { organization_id: organization }),
    Promise.resolve(privateLocale),
  ]);
  return <><div className="eyebrow">{t("p0045")}</div><h1 className="pageTitle">{t("p0046")}</h1>
    <section aria-label={t("p0047")}><h2>{t("p0047")}</h2><div className="grid">{orders.items.map((order) => <article className="card" key={order.id}><h3>{order.id}</h3><p>{order.state} {t("p0016")} {minorAmountPresentation(order.total_minor_units, order.currency, privateLocale.locale).amountLabel}</p>{["placed","confirmed","allocated"].includes(order.state)?<CustomerCheckoutActions orderId={order.id} organizationId={order.organization_id}/>:null}</article>)}</div>
      {orders.items.length === 0 ? <p>{t("p0017")}</p> : null}
<PortalPaging label="Páginas de mis pedidos" nextHref={orders.next_cursor ? portalPageHref("/customer", { orders_after: orders.next_cursor, cases_after: casesAfter }) : undefined} firstHref={ordersAfter ? portalPageHref("/customer", { cases_after: casesAfter }) : undefined} /></section>
    <section aria-label={t("p0048")}><h2>{t("p0048")}</h2><div className="grid">{cases.items.map((item) => <article className="card" key={item.id}><h3>{item.id}</h3><p>{item.state} {t("p0016")} {item.description}</p></article>)}</div>
      {cases.items.length === 0 ? <p>{t("p0017")}</p> : null}
<PortalPaging label="Páginas de mis casos de servicio" nextHref={cases.next_cursor ? portalPageHref("/customer", { cases_after: cases.next_cursor, orders_after: ordersAfter }) : undefined} firstHref={casesAfter ? portalPageHref("/customer", { orders_after: ordersAfter }) : undefined} /></section>
    <section><h2>{t("p0035")}</h2><p><Link href={"/customer/appointments" as Route}>{t("p0049")}</Link></p><div className="grid">{journey.appointments.map((item) => <article className="card" key={item.id}><h3>{item.kind}</h3><p>{item.state} {t("p0016")} <time dateTime={item.starts_at}>{publicAppointmentTime(locale, item.starts_at)} {t("p0050")}{locale.timeZone}{t("p0051")}</time></p></article>)}</div></section>
    <section><h2>{t("p0052")}</h2><div className="grid">{journey.quotes.map((item) => <article className="card" key={item.id}><h3>{item.id}</h3><p>{item.state} {t("p0016")} {minorAmountPresentation(item.total_minor_units, item.currency, privateLocale.locale).amountLabel}</p></article>)}</div></section>
    <section><h2>{t("p0039")}</h2><div className="grid">{journey.handovers.map((item) => <article className="card" key={item.id}><h3>{item.order_id}</h3><p>{item.state}</p></article>)}</div></section>
  </>;
}
