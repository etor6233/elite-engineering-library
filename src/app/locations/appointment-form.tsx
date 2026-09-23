"use client";

import { publicAppointmentTime, publicCount, publicMessage as message, type PublicLocale } from "@/platform/i18n/public-catalog";
import { FormEvent, useEffect, useRef, useState } from "react";
import type { AppointmentSlot } from "@/platform/backend/public-client";

export function AppointmentForm({ leadId, modelId, slots, locale }: { leadId: string; modelId?: string; slots: AppointmentSlot[]; locale: PublicLocale }) {
  const key = useRef(crypto.randomUUID());
  const [pending, setPending] = useState(false);
  const [statusMessage, setMessage] = useState("");
  const [receiptId, setReceiptId] = useState("");
  const [ready, setReady] = useState(false);
  useEffect(() => { setReady(true); }, []);
  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault(); setPending(true); setMessage("");
    const form = new FormData(event.currentTarget);
    const slot = slots.find((item) => item.id === String(form.get("slotId") ?? ""));
    if (!slot) { setMessage(message(locale.locale, "appointment.select")); setPending(false); return; }
    try {
      const response = await fetch("/api/enterprise/appointments", { method: "POST", headers: { "content-type": "application/json", "idempotency-key": key.current }, body: JSON.stringify({ leadId, ...(modelId ? { modelId } : {}), kind: slot.kind, startsAt: slot.starts_at }) });
      if (!response.ok) { setMessage(message(locale.locale, "appointment.failed")); return; }
      const receipt = await response.json() as { id?: unknown; state?: unknown };
      if (typeof receipt.id !== "string" || !/^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(receipt.id) || receipt.state !== "requested") { setMessage(message(locale.locale, "appointment.invalid")); return; }
      setReceiptId(receipt.id);
      setMessage(`${message(locale.locale, "appointment.received")} ${receipt.id}. ${message(locale.locale, "appointment.confirmation")}`);
    } catch { setMessage(message(locale.locale, "appointment.failed")); }
    finally { setPending(false); }
  }
  return <form lang={locale.locale} onSubmit={submit} aria-describedby="appointment-result">
    <label>{message(locale.locale, "appointment.slot")}<select name="slotId" required defaultValue="" disabled={!ready || pending || receiptId !== ""}><option value="" disabled>{message(locale.locale, "appointment.option")}</option>{slots.map((slot) => <option value={slot.id} key={slot.id}>{message(locale.locale, "kind." + slot.kind)} · {publicAppointmentTime(locale, slot.starts_at)} · {publicCount(locale.locale, "appointment", slot.capacity - slot.booked)}</option>)}</select></label>
    {slots.length === 0 ? <p className="notice">{message(locale.locale, "appointment.empty")}</p> : null}
    <button type="submit" disabled={!ready || pending || slots.length === 0 || receiptId !== ""}>{pending ? message(locale.locale, "appointment.sending") : message(locale.locale, "appointment.submit")}</button>
    <p id="appointment-result" role="status" aria-live="polite">{statusMessage}</p>
  </form>;
}
