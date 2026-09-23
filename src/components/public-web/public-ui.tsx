import type { ReactNode } from "react";
import Link from "next/link";
import { publicText as t } from "./messages";
import s from "./public-web.module.css";

// AUTHORED small geometric pictograms. No upstream icon assets or brand marks.
export function PublicIcon({ name }: { name: "arrow" | "pin" | "menu" | "close" | "check" | "mail" | "phone" | "model" }) {
  const paths = { arrow:"M5 12h14 M13 6l6 6-6 6", pin:"M12 21s7-7 7-12a7 7 0 0 0-14 0c0 5 7 12 7 12Z M12 6a3 3 0 1 0 0 6 3 3 0 0 0 0-6", menu:"M4 7h16 M4 12h16 M4 17h16",close:"M6 6l12 12 M18 6 6 18", check:"m5 12 4 4 10-10",mail:"M3 5h18v14H3z m0 0 9 8 9-8",phone:"M7 3H4c-1 8 9 18 17 17v-4l-5-2-2 2-6-6 2-2-3-5Z",model:"M4 6h16v12H4z M4 10h16 M10 10v8" };
  return <svg aria-hidden="true" viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"><path d={paths[name]} /></svg>;
}
export function PublicEmpty({ title, children, locale, retryHref }: { title: string; children?: ReactNode; locale: string; retryHref?: string }) {
  return <div className={s.empty} data-public-state={retryHref ? "error" : "empty"}><span className={s.emptyIcon}><PublicIcon name="model" /></span><h2>{title}</h2>{children && <p>{children}</p>}<a className={s.secondary} href={retryHref ?? "/locations"}>{t(locale,retryHref ? "retry" : "findUs")}<PublicIcon name="arrow" /></a></div>;
}
export function PublicVisit({ locale }: { locale: string }) { return <section className={s.visit}><div><p className={s.eyebrow}>{t(locale,"locations")}</p><h2>{t(locale,"closer")}</h2><p>{t(locale,"closerBody")}</p></div><Link href="/locations" className={s.secondary}><PublicIcon name="pin" />{t(locale,"findUs")}</Link></section>; }
