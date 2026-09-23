import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import Link from "next/link";
import type { Route } from "next";
import { notFound, redirect } from "next/navigation";
import { readSession } from "@/platform/auth/session";
import { loadBusinessConfig } from "@/platform/config/load";
import { roleGuide, sectionsForRole } from "@/platform/roles/role-visibility";
import { RoleDashboard } from "@/components/role-dashboard";
export default async function RoleGuidePage({params}:{params:Promise<{role:string}>}) {
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

  const config=await loadBusinessConfig();if(config.features.role_workspace!==true)notFound();
  const session=await readSession();if(!session)redirect("/api/auth/login?return_to=/guide" as Route);
  const {role}=await params;const guide=roleGuide(role);if(!guide)notFound();
  return <><div className="eyebrow">{t("p0100")}</div><h1 className="pageTitle">{t("p0110")} {controlled(guide.label)}</h1>
    <p className="lede">{controlled(guide.blurb)} {t("p0111")}</p>
    <RoleDashboard sections={sectionsForRole(role,session.permissions,config)}/>
    <p>{t("p0112")}</p>
    <p><Link href="/help">{t("p0113")}</Link></p><p><Link href="/guide">{t("p0114")}</Link></p></>;
}
