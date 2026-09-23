import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import { notFound, redirect } from "next/navigation";
import type { Route } from "next";
import { readSession } from "@/platform/auth/session";
import { loadBusinessConfig } from "@/platform/config/load";
import { visibleSections } from "@/platform/roles/role-visibility";
import { RoleDashboard } from "@/components/role-dashboard";
import { RoleMetrics } from "@/components/role-metrics";
export default async function DashboardPage() {
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

  const config=await loadBusinessConfig();if(config.features.role_workspace!==true)notFound();
  const session=await readSession();if(!session)redirect("/api/auth/login?return_to=/dashboard" as Route);
  return <><div className="eyebrow">{t("p0068")}</div><h1 className="pageTitle">{t("p0069")}</h1>
    <p className="lede">{t("p0070")}</p>
    <RoleDashboard sections={visibleSections(session.permissions,config)}/><RoleMetrics locale={privateLocale.language} organizations={session.organizations} permissions={session.permissions}/></>;
}
