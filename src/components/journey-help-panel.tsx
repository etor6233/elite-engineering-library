"use client";
import {parseHelpQuery} from "@/platform/help/query";
import {guideDisplay} from "@/platform/i18n/private-guide-display";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

import { useCallback, useEffect, useRef, useState } from "react";
import { helpResponseSchema, type HelpResponse } from "@/platform/help/contract";

export function JourneyHelpPanel() {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

  const [query, setQuery] = useState("");
  const [result, setResult] = useState<HelpResponse | null>(null);
  const [status, setStatus] = useState(t("p0690"));
  const [busy, setBusy] = useState(false);
  const request = useRef<AbortController | null>(null);
  const currentQuery = useRef("");
  const clear = useCallback(() => { request.current?.abort(); request.current = null; setResult(null); setBusy(false); }, []);
  const load = useCallback(async (params: string) => {
    clear(); currentQuery.current = params;
    const controller = new AbortController(); request.current = controller;
    setBusy(true); setStatus(t("p0690"));
    try {
      const validated=parseHelpQuery(new URL("/help?"+params,window.location.origin).href);if(!validated)throw new Error("INVALID_HELP_QUERY");
      const parameters=new URLSearchParams(params),search=parameters.get("q")??"";parameters.delete("q");
      const response = await fetch(`/api/enterprise/help${parameters.size ? "?" + parameters.toString() : ""}`, { cache: "no-store", redirect: "error", signal: AbortSignal.any([controller.signal, AbortSignal.timeout(7000)]) });
      if (!response.ok) {
        const text = response.status === 401 ? t("p0691") : response.status === 403 ? t("p0692") : response.status === 404 ? t("p0693") : t("p0694");
        if (request.current === controller) setStatus(text);
        return;
      }
      const data = helpResponseSchema.parse(await response.json());
      data.articles=data.articles.filter(article=>{const view=guideDisplay(article,privateLocale.language);return [view.title,...view.paragraphs].join(" ").normalize("NFC").toLocaleLowerCase(privateLocale.locale).includes(search.trim().normalize("NFC").toLocaleLowerCase(privateLocale.locale))});
      if (request.current === controller && !controller.signal.aborted) { setResult(data); setStatus(data.articles.length ? t("p0695") : t("p0696")); }
    } catch { if (request.current === controller && !controller.signal.aborted) setStatus(t("p0694")); }
    finally { if (request.current === controller) { request.current = null; setBusy(false); } }
  }, [clear,privateLocale]);
  useEffect(() => {
    const initial = window.location.search.slice(1);
    setQuery(new URLSearchParams(initial).get("q")?.slice(0, 160) ?? "");
    void load(initial);
    const visible = () => { if (document.visibilityState === "hidden") { clear(); setStatus(t("p0697")); } else void load(currentQuery.current); };
    const focus = () => { if (document.visibilityState === "visible") void load(currentQuery.current); };
    document.addEventListener("visibilitychange", visible); window.addEventListener("focus", focus);
    return () => { request.current?.abort(); request.current = null; document.removeEventListener("visibilitychange", visible); window.removeEventListener("focus", focus); };
  }, [clear, load]);
  return <section aria-label={t("p0698")}>
    <h1>{t("p0699")}</h1>
    <p>{t("p0700")}</p>
    <form onSubmit={event => { event.preventDefault(); void load(new URLSearchParams({ q: query }).toString()); }}>
      <label>{t("p0701")}<input type="search" maxLength={160} value={query} onChange={event => { clear(); setQuery(event.target.value); currentQuery.current = new URLSearchParams({ q: event.target.value }).toString(); setStatus(t("p1193")); }}/></label>
      <button className="button" type="submit">{t("p0702")}</button>
    </form>
    <button className="button" type="button" disabled={busy} onClick={() => void load(currentQuery.current)}>{t("p0703")}</button>
    <p role="status" aria-live="polite">{status}</p>
    <div aria-busy={busy}>{result?.articles.map(source => {const article=guideDisplay(source,privateLocale.language);return <article lang={article.displayLanguage} className="card" key={article.id + "/" + article.version}>
      <h2>{article.title}</h2><p>{t("p0142")} {article.id}{t("p0280")}{article.version}</p>
      {article.paragraphs.map(text => <p key={text}>{text}</p>)}
      <a href={`/help?article=${encodeURIComponent(article.id)}&version=${encodeURIComponent(article.version)}`}>{t("p0704")}</a>
    </article>})}</div>
  </section>;
}
