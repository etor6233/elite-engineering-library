"use client";
import { GuidedRecordSelect } from "@/components/guided-record-select";
import {usePrivateI18n} from "@/platform/i18n/private-provider";


import { useRef, useState, type FormEvent } from "react";

type ChecklistItem = { id: string; ordinal: number; prompt: string; response_type: string; required: boolean };
type Handover = { id: string; order_id: string; stock_unit_id: string; state: string; version: number; customer_accepted_at?: string; checklist_id?: string; checklist_version?: number; checklist_title?: string; checklist_items: ChecklistItem[]; checklist_completed_at?: string; checklistCompletedLabel?: string; acceptedLabel?: string };
type DeliveryException = { id: string; handover_id: string; reason_code: string; details: string; state: string; version: number; resolution_action?: string; successor_handover_id?: string; return_authorization_id?: string };

export function CustomerHandoverActions({ organizationId, handovers, exceptions }: { organizationId: string; handovers: Handover[]; exceptions: DeliveryException[] }) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

  const [pending, setPending] = useState("");
  const [message, setMessage] = useState("");
  const inFlight = useRef(false);

  async function send(body: Record<string, unknown>) {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), 10_000);
    try {
      const response = await fetch("/api/enterprise/franchise/commands", { method: "POST", headers: { "content-type": "application/json" }, body: JSON.stringify(body), signal: controller.signal });
      const result: unknown = await response.json();
      if (!response.ok || !result || typeof result !== "object" || Array.isArray(result)) throw new Error("UNVERIFIED_RECEIPT");
      return result as Record<string, unknown>;
    } finally { clearTimeout(timeout); }
  }

  async function accept(event: FormEvent<HTMLFormElement>, handover: Handover) {
    event.preventDefault();
    if (inFlight.current) return;
    inFlight.current = true;
    const data = new FormData(event.currentTarget);
    setPending(handover.id);
    setMessage("");
    try {
      const result = await send({ action: "accept-handover", organizationId, handoverId: handover.id, version: handover.version, confirmedReceived: data.get("confirmedReceived") === "yes", serialNumber: String(data.get("serialNumber") ?? ""), checklistId: handover.checklist_id, checklistVersion: handover.checklist_version });
      if (result.id !== handover.id || result.state !== "accepted" || result.version !== handover.version + 1) throw new Error("UNVERIFIED_RECEIPT");
      setMessage(t("p0241"));
      window.location.reload();
    } catch {
      setMessage(t("p0242"));
    }
  }

  async function reject(event: FormEvent<HTMLFormElement>, handover: Handover) {
    event.preventDefault();
    if (inFlight.current) return;
    inFlight.current = true;
    const data = new FormData(event.currentTarget);
    setPending(handover.id);
    setMessage("");
    try {
      const result = await send({ action: "reject-handover", organizationId, handoverId: handover.id, version: handover.version, reasonCode: String(data.get("reasonCode") ?? ""), details: String(data.get("details") ?? "") });
      if (typeof result.id !== "string" || !result.id || result.handover_id !== handover.id || result.state !== "open" || result.version !== 1) throw new Error("UNVERIFIED_RECEIPT");
      setMessage(t("p0243"));
      window.location.reload();
    } catch {
      setMessage(t("p0244"));
    }
  }

  return <>{message && <><p className="status" role="status">{message}</p><p><a href="/customer/handovers">{t("p0245")}</a></p></>}{exceptions.length ? <section className="card"><h2>{t("p0246")}</h2>{exceptions.map((item) => <article key={item.id}><h3>{item.reason_code}</h3><p>{item.details}</p><p>{t("p0247")} {item.state}{item.resolution_action ? t("p0248", {resolution_action:(item.resolution_action)}) : ""}</p>{item.successor_handover_id ? <p>{t("p0249")} {item.successor_handover_id}</p> : null}{item.return_authorization_id ? <p>{t("p0250")} {item.return_authorization_id}</p> : null}</article>)}</section> : null}<div className="grid">{handovers.map((handover) => {
    const checklistReady = Boolean(handover.checklist_id && handover.checklist_version && handover.checklist_completed_at && handover.checklist_items.length);
    return <article className="card" key={handover.id}><h2>{t("p0251")} {handover.order_id}</h2><p>{t("p0247")} {handover.state}</p>{handover.checklist_id ? <section aria-label={t("p0252")}><h3>{handover.checklist_title || t("p0253")} {t("p0254")}{handover.checklist_version}</h3><ol>{handover.checklist_items.map((item) => <li key={item.id}>{item.prompt}{item.required ? " · obligatorio" : ""}</li>)}</ol><p>{handover.checklist_completed_at ? t("p0255", {completed_at:(handover.checklistCompletedLabel ?? t("p0256"))}) : t("p0257")}</p></section> : <p>{t("p0258")}</p>}{handover.state === "presented" && checklistReady ? <><form onSubmit={(event) => void accept(event, handover)}><label>{t("p0259")}<input name="serialNumber" required maxLength={128}/></label><label><input name="confirmedReceived" type="checkbox" value="yes" required/> {t("p0260")}</label><button disabled={pending !== ""}>{t("p0261")}</button></form><form onSubmit={(event) => void reject(event, handover)}><h3>{t("p0262")}</h3><GuidedRecordSelect organization={organizationId} group="customerReasons" name="reasonCode" label="Motivo de la observación" disabled={pending!==""}/><label>{t("p0265")}<textarea name="details" required maxLength={1000}/></label><button disabled={pending !== ""}>{t("p0266")}</button></form></> : <p>{handover.customer_accepted_at ? t("p0267", {accepted_at:(handover.acceptedLabel ?? t("p0256"))}) : t("p0268")}</p>}</article>;
  })}</div></>;
}
