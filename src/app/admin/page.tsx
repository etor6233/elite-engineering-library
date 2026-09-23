import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import type { ReactElement } from "react";
import { PortalPaging, portalCursor, portalPageHref, type PortalPageProps } from "@/platform/backend/portal-paging";
import { minorAmountPresentation } from "@/platform/i18n/money";
import { redirect } from "next/navigation";
import type { Route } from "next";
import { allowed, readSession } from "@/platform/auth/session";
import { protectedGet, type LeadSummary, type OrderSummary, type Overview, type Page, type ServiceCaseSummary } from "@/platform/backend/protected-client";

// AUTHORED framework signature glue: Next inspects the required props overload;
// the zero-argument overload preserves existing server-component test callers.
export default function AdminPage(): Promise<ReactElement>;
export default function AdminPage(props: PortalPageProps): Promise<ReactElement>;
export default async function AdminPage({ searchParams }: PortalPageProps = {}) {
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

  const session = await readSession();
  if (!session) redirect("/api/auth/login?return_to=/admin" as Route);
  if (!allowed(session, "admin:read")) return <><h1 className="pageTitle">{t("p0006")}</h1><div className="notice">{t("p0007")} <code>admin:read</code>{t("p0008")}</div></>;
  const organization = session.organizations[0]!;
  const query = await searchParams ?? {};
  const ordersAfter=portalCursor(query.orders_after), casesAfter=portalCursor(query.cases_after);
  const leadsAfter=allowed(session,"lead:read") ? portalCursor(query.leads_after) : undefined;
  const [overview, orders, cases, leads] = await Promise.all([
    protectedGet<Overview>(session, "/v1/admin/overview", { organization_id: organization }),
    protectedGet<Page<OrderSummary>>(session, "/v1/admin/orders", { organization_id: organization, limit: "25", after: ordersAfter }),
    protectedGet<Page<ServiceCaseSummary>>(session, "/v1/admin/service-cases", { organization_id: organization, limit: "25", after: casesAfter }),
    allowed(session, "lead:read") ? protectedGet<Page<LeadSummary>>(session, "/v1/franchise/leads", { organization_id: organization, limit: "25", after: leadsAfter }) : Promise.resolve(null),
  ]);
  return <>
    <div className="eyebrow">{t("p0009")} {organization}</div><h1 className="pageTitle">{t("p0010")}</h1>
    <div className="grid">
      <article className="card"><h2>{t("p0011")}</h2><p>{overview.orders}</p></article>
      <article className="card"><h2>{t("p0012")}</h2><p>{overview.stock_available}</p></article>
      <article className="card"><h2>{t("p0013")}</h2><p>{overview.open_cases}</p></article>
      <article className="card"><h2>{t("p0014")}</h2><p>{overview.active_shipments}</p></article>
    </div>
    <section aria-label={t("p0015")}><h2>{t("p0015")}</h2><div className="grid">{orders.items.map((order) => <article className="card" key={order.id}><h3>{order.id}</h3><p>{order.state} {t("p0016")} {minorAmountPresentation(order.total_minor_units, order.currency, privateLocale.locale).amountLabel}</p></article>)}</div>
      {orders.items.length===0 ? <p>{t("p0017")}</p> : null}
      <PortalPaging label="Páginas de pedidos recientes" nextHref={orders.next_cursor ? portalPageHref("/admin",{orders_after:orders.next_cursor,cases_after:casesAfter,leads_after:leadsAfter}) : undefined} firstHref={ordersAfter ? portalPageHref("/admin",{cases_after:casesAfter,leads_after:leadsAfter}) : undefined}/></section>
    <section aria-label={t("p0018")}><h2>{t("p0018")}</h2><div className="grid">{cases.items.map((item) => <article className="card" key={item.id}><h3>{item.id}</h3><p>{item.state} {t("p0016")} {item.severity}</p></article>)}</div>
      {cases.items.length===0 ? <p>{t("p0017")}</p> : null}
      <PortalPaging label="Páginas de servicio" nextHref={cases.next_cursor ? portalPageHref("/admin",{cases_after:cases.next_cursor,orders_after:ordersAfter,leads_after:leadsAfter}) : undefined} firstHref={casesAfter ? portalPageHref("/admin",{orders_after:ordersAfter,leads_after:leadsAfter}) : undefined}/></section>
    {leads && <section aria-label={t("p0019")}><h2>{t("p0019")}</h2><div className="grid">{leads.items.map((item) => <article className="card" key={item.id}><h3>{item.id}</h3><p>{item.state} {t("p0016")} {item.source_code} {t("p0016")} {item.assigned_subject || t("p0020")}</p></article>)}</div>
      {leads.items.length===0 ? <p>{t("p0017")}</p> : null}
      <PortalPaging label="Páginas de oportunidades" nextHref={leads.next_cursor ? portalPageHref("/admin",{leads_after:leads.next_cursor,orders_after:ordersAfter,cases_after:casesAfter}) : undefined} firstHref={leadsAfter ? portalPageHref("/admin",{orders_after:ordersAfter,cases_after:casesAfter}) : undefined}/></section>}
  </>;
}
