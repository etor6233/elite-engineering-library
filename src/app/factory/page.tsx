import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import { createHash } from "node:crypto";
import { FactoryUnitActions } from "@/components/factory-unit-actions";
import { factoryUnitSchema } from "@/platform/factory/contracts";
import type { ReactElement } from "react";
import { PortalPaging, portalCursor, portalPageHref, type PortalPageProps } from "@/platform/backend/portal-paging";
import { redirect } from "next/navigation";
import type { Route } from "next";
import { allowed, readSession } from "@/platform/auth/session";
import { protectedGet, type FactoryUnitSummary, type Page } from "@/platform/backend/protected-client";

// AUTHORED framework signature glue: Next inspects the required props overload;
// the zero-argument overload preserves existing server-component test callers.
export default function FactoryPage(): Promise<ReactElement>;
export default function FactoryPage(props: PortalPageProps): Promise<ReactElement>;
export default async function FactoryPage({ searchParams }: PortalPageProps = {}) {
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

  const session = await readSession();
  if (!session) redirect("/api/auth/login?return_to=/factory" as Route);
  if (!allowed(session, "factory:read")) return <><h1 className="pageTitle">{t("p0006")}</h1><div className="notice">{t("p0007")} <code>factory:read</code>{t("p0008")}</div></>;
  const organization = session.organizations[0]!;
  const storageScope=createHash("sha256").update(JSON.stringify([session.tenantId,session.subject,organization])).digest("hex");
  const query = await searchParams ?? {};
  const unitsAfter = portalCursor(query.units_after);
  const units = await protectedGet<Page<FactoryUnitSummary>>(session, "/v1/factory/units", { organization_id: organization, limit: "25", after: unitsAfter });
  return <><div className="eyebrow">{t("p0071")} {organization}</div><h1 className="pageTitle">{t("p0072")}</h1>
    <div className="grid">{units.items.map((unit) => <article className="card" key={unit.id}><h2>{unit.serial_number}</h2><p>{t("p0073")} {unit.purchase_order_id}</p><FactoryUnitActions unit={factoryUnitSchema.parse(unit)} canWrite={allowed(session,"factory:write")} storageScope={storageScope}/></article>)}</div>
    {units.items.length === 0 ? <p>{t("p0017")}</p> : null}
    <PortalPaging label="Páginas de unidades" nextHref={units.next_cursor ? portalPageHref("/factory", { units_after: units.next_cursor }) : undefined} firstHref={unitsAfter ? "/factory" : undefined} />
  </>;
}
