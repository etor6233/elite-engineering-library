"use client";
import{usePrivateI18n}from"@/platform/i18n/private-provider";
type PrivateReadPath = "/admin" | "/customer" | "/factory";

export function PrivateReadFailure({ path }: { path: PrivateReadPath }) {
  const{t}=usePrivateI18n();
  return <section aria-labelledby="private-read-title">
    <h1 id="private-read-title" className="pageTitle">{t("p1194")}</h1>
    <p className="notice" role="alert">{t("p1195")}</p>
    <p>{t("p1196")}</p>
    <a className="button" href={path}>{t("p0033")}</a>
  </section>;
}
