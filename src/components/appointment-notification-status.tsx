"use client";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

import { NOTIFICATION_GUIDE } from "@/platform/help/content";

import { useEffect, useRef, useState } from "react";
import { readBackendResponse } from "@/platform/backend/public-client";
import { notificationExplanation, notificationHistorySchema, STATUS_VIEW_VERSION, type NotificationHistory } from "@/platform/notifications/status-contract";



// No automatic polling/N+1 page load or write action. Abort stale reads when
// scope changes/unmounts, following React's effect cleanup contract.
export function AppointmentNotificationStatus({ organization, appointment }: { organization: string; appointment: string }) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();
 const displayTime=(value:string|number)=>new Date(value).toLocaleString(privateLocale.locale,{timeZone:privateLocale.timeZone,dateStyle:"short",timeStyle:"medium"})+" ("+privateLocale.timeZone+")";

  const [ready, setReady] = useState(false), [busy, setBusy] = useState(false);
  const [history, setHistory] = useState<NotificationHistory | null>(null), [error, setError] = useState("");
  const request = useRef<AbortController | null>(null);
  useEffect(() => {
    setReady(true); setHistory(null); setError(""); setBusy(false);
    return () => { request.current?.abort(); request.current = null; };
  }, [organization, appointment]);
  async function refresh() {
    if (request.current) return;
    const controller = new AbortController(); request.current = controller;
    setBusy(true); setError(""); setHistory(null);
    try {
      const query = new URLSearchParams({ organizationId: organization, appointmentId: appointment });
      const response = await fetch(`/api/enterprise/franchise/notifications?${query}`, { cache: "no-store", redirect: "error", signal: AbortSignal.any([controller.signal, AbortSignal.timeout(7000)]) });
      const value = notificationHistorySchema.parse(await readBackendResponse<unknown>(response));
      if (value.organization_id !== organization || value.appointment_id !== appointment) throw new Error("scope mismatch");
      if (!controller.signal.aborted) setHistory(value);
    } catch {
      if (!controller.signal.aborted) setError(t("p0126"));
    } finally { if (request.current === controller) { request.current = null; setBusy(false); } }
  }
  const visible = history?.organization_id === organization && history.appointment_id === appointment ? history : null;
  return <section aria-label={t("p0127")} className="notificationStatus">
    <h4>{t("p0128")}</h4>
    <button type="button" className="button" disabled={!ready || busy} onClick={() => void refresh()}>{busy ? t("p0129") : t("p0130")}</button>
    <p role="status" aria-label={t("p0131")} aria-live="polite">{error || (busy ? t("p0132") : visible ? t("p0133") : t("p0134"))}</p>
    {visible?.items.length === 0 ? <p>{t("p0135")}</p> : null}
    {visible?.items.map(item => {
      const explanation = notificationExplanation(item.status);
      return <article key={item.confirmation_event_id}>
        <h5>{controlled(explanation.title)}</h5><p>{controlled(explanation.guidance)}</p>
        <p>{t("p0136")} <time dateTime={item.status.observed_at}>{displayTime(item.status.observed_at)}</time></p>
        {item.status.provider_timestamp ? <p>{t("p0137")} <time dateTime={new Date(item.status.provider_timestamp * 1000).toISOString()}>{displayTime(item.status.provider_timestamp * 1000)}</time></p> : null}
        {item.status.approval_expired ? <p>{t("p0138")}</p> : null}
        <details><summary>{t("p0139")}</summary>
          <p>{controlled(NOTIFICATION_GUIDE.paragraphs[0])}</p>
          <dl><dt>{t("p0140")}</dt><dd>{appointment}</dd><dt>{t("p0141")}</dt><dd>{item.confirmation_event_id}</dd><dt>{t("p0142")}</dt><dd>{STATUS_VIEW_VERSION}</dd></dl>
          <p>{controlled(NOTIFICATION_GUIDE.paragraphs[1])}</p>
        <a href={`/help?article=${NOTIFICATION_GUIDE.id}&version=${NOTIFICATION_GUIDE.version}`}>{t("p0143")}</a>
        </details>
      </article>;
    })}
  </section>;
}
