import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import Link from "next/link";
import type { Route } from "next";
import { notFound, redirect } from "next/navigation";
import { readSession } from "@/platform/auth/session";
import { loadBusinessConfig } from "@/platform/config/load";
import { ROLE_GUIDES } from "@/platform/roles/role-visibility";
export default async function GuidePage() {
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

  const config=await loadBusinessConfig();if(config.features.role_workspace!==true)notFound();
  const session=await readSession();if(!session)redirect("/api/auth/login?return_to=/guide" as Route);
  return <><div className="eyebrow">{t("p0100")}</div><h1 className="pageTitle">{t("p0101")}</h1>
    <p className="lede">{t("p0102")}</p>
    <div className="grid">{ROLE_GUIDES.map(guide=><article className="card" key={guide.id}><h2>{controlled(guide.label)}</h2><p>{controlled(guide.blurb)}</p><Link href={`/guide/${guide.id}` as Route} className="button">{t("p0103")} {controlled(guide.label)}</Link></article>)}</div>
    <p><Link href="/dashboard">{t("p0104")}</Link></p></>;
}
