"use client";
import {usePrivateI18n} from "@/platform/i18n/private-provider";
import Link from "next/link";
import type { Route } from "next";
import type { Section } from "@/platform/roles/role-visibility";
export function RoleDashboard({sections}:{sections:readonly Section[]}) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

  return <section className="card" aria-label={t("p0752")}>
    <h2>{t("p0753")}</h2><p>{t("p0754")}</p>
    <nav aria-label={t("p0755")}><ul className="grid">{sections.map(section=><li key={section.id}><Link href={section.href as Route} className="button">{controlled(section.label)}</Link></li>)}</ul></nav>
    {sections.length===0?<p>{t("p0756")}</p>:null}
  </section>;
}
