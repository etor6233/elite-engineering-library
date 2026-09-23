import {privateTranslator} from "@/platform/i18n/private-catalog";
import type {PrivateLocale} from "@/platform/i18n/private-locale";
import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import { minorAmountPresentation } from "@/platform/i18n/money";
import { redirect } from "next/navigation";
import type { Route } from "next";
import { allowed, readSession } from "@/platform/auth/session";
import { protectedGet, type CustomerJourney } from "@/platform/backend/protected-client";
import { CustomerQuoteActions, QuoteAcceptanceHelp } from "@/components/customer-quote-actions";
import Link from "next/link";

// AUTHORED presentation only. ISO currency digits from the pinned runtime's Intl;
// no pricing, exchange, tax or rounding rule is introduced. Compute on server
// once so browser ICU/timezone differences do not cause hydration mismatches.
function quotePresentation(quote: CustomerJourney["quotes"][number], asOf: number, locale:PrivateLocale) {
 const t=privateTranslator(locale.language);

  const { amountLabel, amountValid } = minorAmountPresentation(quote.total_minor_units, quote.currency, locale.locale);
  const until = Date.parse(quote.valid_until);
  const dateValid = Number.isFinite(until);
  const expired = dateValid && until <= asOf;
  const states: Record<string, string> = { issued: expired ? t("p0053") : t("p0054"), accepted: t("p0055"), cancelled: t("p0056"), draft: t("p0057"), expired: t("p0053") };
  return { ...quote, amountLabel, statusLabel: states[quote.state] ?? t("p0058"),
    validUntilLabel: dateValid ? new Intl.DateTimeFormat(locale.locale, { dateStyle: "medium", timeStyle: "short", timeZone: locale.timeZone }).format(until) + " (" + locale.timeZone + ")" : t("p0059"),
    acceptanceAvailable: quote.state === "issued" && amountValid && dateValid && !expired };
}

export default async function CustomerQuotesPage() {
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

  const session = await readSession();
  if (!session) redirect("/api/auth/login?return_to=/customer" as Route);
  if (!allowed(session, "customer:self")) return <><h1 className="pageTitle">{t("p0006")}</h1><div className="notice">{t("p0007")} <code>customer:self</code>{t("p0008")}</div></>;
  const organization = session.organizations[0]!;
  let journey: CustomerJourney;
  try {
    journey = await protectedGet<CustomerJourney>(session, "/v1/customer/journey", { organization_id: organization });
    if (!Array.isArray(journey.quotes)) throw new Error("INVALID_JOURNEY");
  } catch {
    return <>
      <h1 className="pageTitle">{t("p0052")}</h1>
      <p role="alert">{t("p0060")}</p>
      <p><a className="button" href="/customer/quotes">{t("p0033")}</a></p>
      <QuoteAcceptanceHelp/>
    </>;
  }
  return <>
    <div className="eyebrow">{t("p0034")} {organization}</div>
    <h1 className="pageTitle">{t("p0052")}</h1>
    <p className="lede">{t("p0061")}</p>
    <p><Link href={"/customer/appointments" as Route}>{t("p0062")}</Link></p>
    <CustomerQuoteActions organizationId={organization} quotes={journey.quotes.map(quote => quotePresentation(quote, Date.now(), privateLocale))}/>
  </>;
}
