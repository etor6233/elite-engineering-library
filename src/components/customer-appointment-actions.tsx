"use client";
import {usePrivateI18n} from "@/platform/i18n/private-provider";


import { useRef, useState } from "react";

type Appointment = { id: string; kind: string; starts_at: string; startsAtLabel: string; state: string; version: number };

// The response must identify this exact mutation, not merely be an HTTP 200.
export function isCancellationReceipt(value: unknown, organizationId: string, appointment: { id: string; version: number }): boolean {
  if (typeof value !== "object" || value === null || Array.isArray(value)) return false;
  const receipt = value as Record<string, unknown>;
  return receipt.id === appointment.id && receipt.organization_id === organizationId &&
    receipt.state === "cancelled" && Number.isSafeInteger(receipt.version) &&
    Number.isSafeInteger(appointment.version + 1) && receipt.version === appointment.version + 1;
}

export function CustomerAppointmentActions({ organizationId, appointments }: { organizationId: string; appointments: Appointment[] }) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

  const [pending, setPending] = useState("");
  const [message, setMessage] = useState("");
  const [uncertain, setUncertain] = useState(false);
  const inFlight = useRef(false);

  async function cancel(appointment: Appointment) {
    if (inFlight.current) return;
    inFlight.current = true;
    setPending(appointment.id);
    setMessage(t("p0231"));
    const controller = new AbortController();
    const timer = setTimeout(() => controller.abort(), 10000);
    try {
      const response = await fetch("/api/enterprise/franchise/commands", { method: "POST", headers: { "content-type": "application/json" }, signal: controller.signal, body: JSON.stringify({ action: "cancel-customer-appointment", organizationId, appointmentId: appointment.id, version: appointment.version, reasonCode: "customer-request" }) });
      const result: unknown = await response.json();
      if (!response.ok || !isCancellationReceipt(result, organizationId, appointment)) throw new Error("RESULT_NOT_CONFIRMED");
      setMessage(t("p0232"));
      window.location.reload();
    } catch {
      setUncertain(true);
      setMessage(t("p0233"));
      // Keep commands disabled; recover with an explicit GET, never an automatic POST retry.
    } finally {
      clearTimeout(timer);
    }
  }

  return <>
    {message && <p className="status" role="status" aria-live="polite">{message}</p>}
    {uncertain && <p><a href="/customer/appointments">{t("p0234")}</a></p>}
    {appointments.length === 0 && <p className="notice">{t("p0235")}</p>}
    <div className="grid">{appointments.map((appointment) => <article className="card" key={appointment.id}>
      <h2>{appointment.kind}</h2><p>{appointment.state} {t("p0016")} <time dateTime={appointment.starts_at}>{appointment.startsAtLabel}</time></p>
      {["requested", "confirmed"].includes(appointment.state) && new Date(appointment.starts_at).getTime() > Date.now() ? <button disabled={pending !== ""} onClick={() => void cancel(appointment)}>{t("p0236")}</button> : null}
    </article>)}</div>
  </>;
}
