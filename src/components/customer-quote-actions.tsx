"use client";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

import { isQuoteAcceptanceReceipt } from "@/platform/payments/quote-receipt";
import { QUOTE_GUIDE } from "@/platform/help/content";

import { useState } from "react";

type Quote = { id: string; currency: string; total_minor_units: number; valid_until: string; state: string; version: number; order_id?: string; amountLabel: string; validUntilLabel: string; statusLabel: string; acceptanceAvailable: boolean };

export function CustomerQuoteActions({ organizationId, quotes }: { organizationId: string; quotes: Quote[] }) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

  const [pending, setPending] = useState("");
  const [message, setMessage] = useState("");
  const [uncertain, setUncertain] = useState(false);

  async function accept(quote: Quote) {
    if (pending || uncertain || !quote.acceptanceAvailable) return;
    if (Date.parse(quote.valid_until) <= Date.now()) {
      setUncertain(true);
      setMessage(t("p0269"));
      return;
    }
    setPending(quote.id);
    setMessage(t("p0270"));
    const controller = new AbortController();
    const timer = setTimeout(() => controller.abort(), 10000);
    try {
      const response = await fetch("/api/enterprise/franchise/commands", {
        method: "POST",
        headers: { "content-type": "application/json" },
        signal: controller.signal,
        body: JSON.stringify({ action: "accept-quote", organizationId, quoteId: quote.id, version: quote.version })
      });
      const result: unknown = await response.json();
      if (!response.ok || !isQuoteAcceptanceReceipt(result, organizationId, quote)) {
        throw new Error("RESULT_NOT_CONFIRMED");
      }
      setMessage(t("p0271"));
      window.location.reload();
    } catch {
      setUncertain(true);
      setMessage(t("p0272"));
      // Keep acceptance disabled. GET recovery, never an automatic POST retry.
    } finally {
      clearTimeout(timer);
    }
  }

  return <>
    {message && <p className="status" role="status" aria-live="polite">{message}</p>}
    {uncertain && <button className="button" onClick={() => window.location.reload()}>{t("p0273")}</button>}
    {quotes.length === 0 && <p className="notice">{t("p0274")}</p>}
    <div className="grid">{quotes.map((quote) => <article className="card" key={quote.id}>
      <h2>{quote.id}</h2>
      <p>{t("p0275")} {quote.amountLabel}</p>
      <p>{t("p0247")} {quote.statusLabel}</p>
      <p>{t("p0276")} {quote.validUntilLabel}</p>
      {quote.state === "accepted" && quote.order_id ? <p role="note">{t("p0277")} <strong>{quote.order_id}</strong></p> : null}
      {quote.acceptanceAvailable ? <button className="button" disabled={pending !== "" || uncertain} onClick={() => void accept(quote)}>{t("p0278")}</button> : null}
    </article>)}</div>
    <QuoteAcceptanceHelp/>
  </>;
}

export function QuoteAcceptanceHelp() {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

  return <details className="card">
    <summary>{t("p0279")}</summary>
    <p>{t("p0142")} {QUOTE_GUIDE.id}{t("p0280")}{QUOTE_GUIDE.version}</p>
    {QUOTE_GUIDE.paragraphs.map(text => <p key={text}>{controlled(text)}</p>)}
    <a href={`/help?article=${QUOTE_GUIDE.id}&version=${QUOTE_GUIDE.version}`}>{t("p0143")}</a>
  </details>;
}
