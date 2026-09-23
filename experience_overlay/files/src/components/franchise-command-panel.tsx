"use client";
import {experienceText} from "@/platform/experience/messages";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

import { definitiveSelectionRejection } from "@/platform/experience/rejection";
import {TaskWorkspace} from "@/components/task-workspace";
import { GuidedRecordSelect } from "@/components/guided-record-select";
import { GuidedChecklistFields } from "@/components/guided-checklist-fields";
import { checklistResponses, checklistSchema } from "@/platform/experience/checklists";

import { OperationalGuide } from "@/components/operational-guide";
import { OPERATIONAL_GUIDES } from "@/platform/help/content";


import { useEffect, useRef, useState, useTransition, type FormEvent } from "react";
import { useRouter } from "next/navigation";
import { AppointmentNotificationStatus } from "./appointment-notification-status";

type Lead = { id: string; organization_id: string; created_at?: string; state: string; source_code: string; assigned_subject?: string; version: number };
type Availability = { id: string; resource_id?: string; entry_type: "working" | "unavailable"; reason_code?: string; starts_at: string; ends_at: string; state: string; version: number };
type DeliveryException = { id: string; handover_id: string; customer_subject: string; reason_code: string; details: string; state: string; version: number; resolution_action?: string };
type ReturnEffect = { id: string; effect_kind: "inventory" | "refund" | "exchange" | "accounting" | "fiscal"; owner_context: string; state: "requested" };
type ReturnCase = { authorization_id: string; order_id: string; stock_unit_id: string; customer_subject: string; authorized_action: "return" | "exchange"; receipt?: { id: string; received_serial_number: string; condition_code: string; received_at: string }; disposition?: { id: string; inventory_action: string; customer_remedy: string; effect_requests: ReturnEffect[] } };
type Command = Record<string, unknown> & { action: string };
type ChecklistDraftItem = { key: number; id: string; prompt: string; responseType: "confirmation" | "text" | "serial"; required: boolean };

const transitions: Record<string, string[]> = { new: ["contacted", "lost"], contacted: ["qualified", "lost"], qualified: ["converted", "lost"] };

export type OrderOperations = {
  payment_provider?: string;
  orders: { id: string; state: string; version: number; lines: { id: string; variant_id: string; name: string; quantity: number; stock_id: string }[]; payments: { id: string; state: string }[]; handovers: { id: string; state: string }[] }[];
  stock: { id: string; variant_id: string; serial: string; version: number }[];
  truncated: boolean;
};

function OrderPaymentRequest({order,organization,provider,scope,disabled}:{order:OrderOperations["orders"][number];organization:string;provider:string;scope:string;disabled:boolean}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

  const [message,setMessage]=useState("");const [locked,setLocked]=useState(false);const inFlight=useRef(false);
  async function requestPayment(){
    if(disabled||locked||inFlight.current)return;
    inFlight.current=true;setLocked(true);setMessage(t("p0308"));
    let key:string;
    try{
      const storageKey=`elite-payment-request:${scope}:${order.id}`;
      key=window.sessionStorage.getItem(storageKey)??crypto.randomUUID();
      if(!/^[a-f0-9-]{36}$/i.test(key))throw new Error("invalid retained key");
      window.sessionStorage.setItem(storageKey,key);
    }catch{setMessage(t("p0309"));return}
    try{
      const response=await fetch("/api/enterprise/franchise/commands",{method:"POST",headers:{"content-type":"application/json"},signal:AbortSignal.timeout(10000),body:JSON.stringify({action:"request-order-payment",organizationId:organization,orderId:order.id,requestKey:key})});
      const value=await response.json() as {id?:unknown;order_id?:unknown;organization_id?:unknown;provider_code?:unknown;state?:unknown;version?:unknown};
      if(response.ok&&typeof value.id==="string"&&value.id.length>0&&value.order_id===order.id&&value.organization_id===organization&&value.provider_code===provider&&typeof value.state==="string"&&["created","pending","authorized","captured","failed","refunded","disputed"].includes(value.state)&&typeof value.version==="number"&&Number.isSafeInteger(value.version)&&value.version>0){
        setMessage(t("p0310", {request_id:(value.id)}));
      }else if(response.status===409){setMessage(t("p0311"))}
      else if(response.status===401||response.status===403){setMessage(t("p0312"))}
      else throw new Error("unconfirmed");
    }catch{setMessage(t("p0313"))}
  }
  return <div><button className="button" type="button" disabled={disabled||locked} onClick={()=>void requestPayment()}>{t("p0314")}</button><p>{t("p0315")} {provider}{t("p0316")}</p><p role="status" aria-label={`Resultado de pago ${order.id}`} aria-live="polite">{message}</p></div>;
}

export function OrderOperationsPanel({ snapshot, organization, canAllocate, canRequestPayment=false, paymentScope="" }: { snapshot: OrderOperations | null; organization: string; canAllocate: boolean;canRequestPayment?:boolean;paymentScope?:string }) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

  const [ready,setReady]=useState(false);
  const [message,setMessage]=useState("");
  const [locked,setLocked]=useState(false);
  const inFlight=useRef(false);
  useEffect(()=>setReady(true),[]);
  async function allocate(event: FormEvent<HTMLFormElement>, order: OrderOperations["orders"][number], line: OrderOperations["orders"][number]["lines"][number]) {
    event.preventDefault();
    if(inFlight.current || locked || !snapshot || snapshot.truncated || !canAllocate) return;
    const stock=snapshot.stock.find(item=>item.id===new FormData(event.currentTarget).get("stock") && item.variant_id===line.variant_id);
    if(!stock) {setMessage(t("p0317"));return;}
    inFlight.current=true;setLocked(true);setMessage(t("p0318"));
    try {
      const response=await fetch("/api/enterprise/franchise/commands",{method:"POST",headers:{"content-type":"application/json"},signal:AbortSignal.timeout(10000),body:JSON.stringify({action:"allocate-order-stock",organizationId:organization,orderId:order.id,lineId:line.id,stockUnitId:stock.id,orderVersion:order.version,stockVersion:stock.version})});
      const result=await response.json() as {status?:unknown};
      if(response.ok && result.status==="accepted") setMessage(t("p0319"));
      else if(response.status===409) setMessage(t("p0320"));
      else if(response.status===401 || response.status===403) setMessage(t("p0321"));
      else throw new Error("unconfirmed");
    } catch {setMessage(t("p0322"));}
    // Remain fenced after every outcome. Recovery is a fresh server-side GET,
    // never an automatic replay of a possibly committed command.
  }
  const state=(value:string)=>({placed:t("p0323"),confirmed:t("p0324"),delivered:t("p0325"),cancelled:t("p0326"),created:t("p0327"),pending:t("p0328"),authorized:t("p0329"),captured:t("p0330"),failed:t("p0331"),refunded:t("p0332"),disputed:t("p0333"),prepared:t("p0334"),presented:t("p0335"),accepted:t("p0336"),rejected:t("p0337")}[value] ?? t("p0338"));
  return <section className="card" aria-labelledby="order-operations-title">
    <h2 id="order-operations-title">{t("p0339")}</h2>
    <p>{t("p0340")}</p>
    <button className="button" type="button" disabled={!ready} onClick={()=>window.location.reload()}>{t("p0341")}</button>
    <p role="status" aria-label={t("p0342")} aria-live="polite">{message}</p>
    {!snapshot ? <p role="alert">{t("p0343")}</p> : <>
      {snapshot.truncated ? <p role="alert">{t("p0344")}</p> : null}
      {snapshot.orders.length===0 ? <p>{t("p0345")}</p> : snapshot.orders.map(order=><article className="card" key={order.id} aria-label={`Pedido ${order.id}`}>
        <h3>{t("p0251")} {order.id}</h3><p>{state(order.state)}</p>
        {order.lines.map(line=><div key={line.id}><h4>{line.name} {t("p0346")} {line.quantity}</h4>
          {line.stock_id ? <p>{t("p0347")} <span>{line.stock_id}</span></p> : <><p>{t("p0348")}</p>
          {canAllocate && ["placed","confirmed"].includes(order.state) && line.quantity===1 ? <form onSubmit={event=>void allocate(event,order,line)}>
            <label>{t("p0349")} {line.name}<select name="stock" required defaultValue="" disabled={!ready || locked || snapshot.truncated}><option value="" disabled>{t("p0350")}</option>{snapshot.stock.filter(item=>item.variant_id===line.variant_id).map(item=><option key={item.id} value={item.id}>{item.serial}</option>)}</select></label>
            {!snapshot.stock.some(item=>item.variant_id===line.variant_id) ? <p>{t("p0351")}</p> : null}
            <button className="button" disabled={!ready || locked || snapshot.truncated || !snapshot.stock.some(item=>item.variant_id===line.variant_id)}>{snapshot.stock.some(item=>item.variant_id===line.variant_id) ? t("p0352") : t("p0353")}</button>
          </form> : <p>{t("p0354")}</p>}</>}
        </div>)}
        <p>{t("p0355")} {order.payments.length ? order.payments.map(item=>`${state(item.state)} (${item.id})`).join("; ") : t("p0356")}</p>
        {!snapshot.payment_provider ? <p>{t("p0357")}</p> : canRequestPayment&&paymentScope&&order.payments.length===0&&["placed","confirmed"].includes(order.state) ? <OrderPaymentRequest order={order} organization={organization} provider={snapshot.payment_provider} scope={paymentScope} disabled={!ready||locked||snapshot.truncated}/> : null}
        <p>{t("p0358")} {order.handovers.length ? order.handovers.map(item=>`${state(item.state)} (${item.id})`).join("; ") : t("p0359")}</p>
      </article>)}
    </>}
    <OperationalGuide guide={OPERATIONAL_GUIDES["order-operations-view"]}/>
  </section>;
}

export type AppointmentAgenda = {
  appointments: { id: string; organization_id: string; lead_id: string; kind: string; starts_at: string; ends_at: string; state: string; version: number; resource_id?: string }[];
  resources: { id: string; display_name: string; skills: string[] }[];
  truncated: boolean;
};

export function AppointmentAgendaPanel({ agenda, organization, notificationStatusEnabled = false }: { agenda: AppointmentAgenda; organization: string; notificationStatusEnabled?: boolean }) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

  const router = useRouter();
  const [ready, setReady] = useState(false);
  const [busy, setBusy] = useState(false);
  const [refreshing, startRefresh] = useTransition();
  const [uncertain, setUncertain] = useState(false);
  const [message, setMessage] = useState("");
  useEffect(() => setReady(true), []);
  const disabled = !ready || busy || refreshing || uncertain || agenda.truncated;
  async function apply(appointment: AppointmentAgenda["appointments"][number], command: Command, expectedState: string) {
    setBusy(true); setMessage("");
    try {
      const response = await fetch("/api/enterprise/franchise/commands", { method: "POST", headers: { "content-type": "application/json" }, body: JSON.stringify({ ...command, organizationId: organization, appointmentId: appointment.id, version: appointment.version }) });
      const receipt = await response.json() as { id?: unknown; organization_id?: unknown; state?: unknown; version?: unknown; resource_id?: unknown };
      if (!response.ok) throw new Error("command not confirmed");
      if (receipt.id !== appointment.id || receipt.organization_id !== organization || receipt.state !== expectedState || receipt.version !== appointment.version + 1 || (command.action === "assign-appointment-resource" && receipt.resource_id !== command.resourceId)) throw new Error("receipt differs");
      setMessage(expectedState === "confirmed" ? t("p0360") : t("p0361"));
      startRefresh(() => router.refresh());
    } catch {
      setUncertain(true);
      setMessage(t("p0362"));
    } finally { setBusy(false); }
  }
  const names: Record<string, string> = { requested: t("p0363"), confirmed: t("p0324"), completed: t("p0364"), cancelled: t("p0326"), "no-show": t("p0365") };
  const kinds: Record<string,string> = { consultation: t("p0366"), "test-drive": t("p0367"), delivery: t("p0368"), service: t("p0018") };
  return <section className="card" aria-labelledby="appointment-agenda-title">
    <h2 id="appointment-agenda-title">{t("p0369")}</h2>
    <p>{t("p0370")}</p>
    <button className="button" type="button" disabled={!ready || busy || refreshing} onClick={() => window.location.reload()}>{t("p0371")}</button>
    <p role="status" aria-label={t("p0372")} aria-live="polite">{message}</p>
    {agenda.truncated ? <p role="alert">{t("p0373")}</p> : null}
    {!agenda.appointments.length ? <p>{t("p0374")}</p> : <div className="grid">{agenda.appointments.map((appointment) => <article className="card" key={appointment.id} aria-label={`Turno ${appointment.kind} ${appointment.starts_at}`}>
      <h3>{kinds[appointment.kind] ?? appointment.kind} {t("p0016")} <time dateTime={appointment.starts_at}>{new Date(appointment.starts_at).toISOString().slice(0,16).replace("T"," ")} {t("p0375")}</time></h3>
      <p>{t("p0247")} {names[appointment.state] ?? appointment.state}</p>
      <p>{t("p0167")} {appointment.id}</p>
      <p>{t("p0376")} {agenda.resources.find((item) => item.id === appointment.resource_id)?.display_name ?? (appointment.resource_id ? t("p0377") : t("p0378"))}</p>
      {appointment.state === "requested" && !appointment.resource_id ? <form onSubmit={(event) => { event.preventDefault(); const data = new FormData(event.currentTarget); void apply(appointment, { action: "assign-appointment-resource", resourceId: String(data.get("resourceId") ?? "") }, "requested"); }}>
        <label>{t("p0379")}<select name="resourceId" required defaultValue="" disabled={disabled}><option value="" disabled>{t("p0380")}</option>{agenda.resources.filter((item) => item.skills.includes(appointment.kind)).map((item) => <option key={item.id} value={item.id}>{item.display_name}</option>)}</select></label>
        <button className="button" disabled={disabled}>{t("p0381")}</button>
      </form> : null}
      {appointment.state === "requested" && appointment.resource_id ? <button className="button" disabled={disabled} onClick={() => void apply(appointment, { action: "transition-appointment", current: appointment.state, target: "confirmed" }, "confirmed")}>{t("p0382")}</button> : null}
      {["requested","confirmed"].includes(appointment.state) ? <form onSubmit={(event) => { event.preventDefault(); const data=new FormData(event.currentTarget); const target=String(data.get("target")); const reasonCode=String(data.get("reasonCode") ?? "").trim(); void apply(appointment,{action:"transition-appointment",current:appointment.state,target,...(reasonCode?{reasonCode}:{})},target); }}>
        <label>{t("p0383")}<select name="target" disabled={disabled}><option value="cancelled">{t("p0384")}</option>{appointment.state==="confirmed"?<><option value="completed">{t("p0385")}</option><option value="no-show">{t("p0386")}</option></>:null}</select></label>
        <GuidedRecordSelect organization={organization} group="appointmentReasons" name="reasonCode" label="Motivo del cambio" disabled={disabled} optional/><button className="button" disabled={disabled}>{t("p0388")}</button>
      </form> : null}
      {notificationStatusEnabled ? <AppointmentNotificationStatus key={`${organization}:${appointment.id}`} organization={organization} appointment={appointment.id}/> : null}
    </article>)}</div>}
  </section>;
}

function AvailabilityCancellation({item,organization,scope,canCancel}:{item:Availability;organization:string;scope:string;canCancel:boolean}) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

  const [ready,setReady]=useState(false),[locked,setLocked]=useState(true),[message,setMessage]=useState("");
  const inFlight=useRef(false);
  const key=`elite-availability-cancel:${scope}:${item.id}`;
  useEffect(()=>{
    try {
      if(!/^[a-f0-9]{64}$/.test(scope)||!Number.isSafeInteger(item.version)||item.version<1)throw new Error("invalid scope/version");
      const raw=window.sessionStorage.getItem(key);
      if(raw!==null) {
        const marker=JSON.parse(raw) as {version?:unknown}|null;
        if(!marker||typeof marker!=="object"||Object.keys(marker).join(",")!=="version"||!Number.isSafeInteger(marker.version)||Number(marker.version)<1)throw new Error("invalid marker");
        if(item.state==="cancelled"&&item.version>Number(marker.version)) {
          window.sessionStorage.removeItem(key);
          setMessage(t("p0389"));
        } else setMessage(t("p0390"));
        setLocked(true);
      } else setLocked(item.state!=="active");
      setReady(true);
    } catch {setReady(true);setLocked(true);setMessage(t("p0391"))}
  },[key,scope,item.version,item.state]);
  async function cancel() {
    if(!ready||locked||inFlight.current||!canCancel||item.state!=="active")return;
    inFlight.current=true;setLocked(true);
    try {if(window.sessionStorage.getItem(key)!==null)throw new Error("unresolved");window.sessionStorage.setItem(key,JSON.stringify({version:item.version}))}
    catch {setMessage(t("p0392"));return}
    setMessage(t("p0393"));
    try {
      const response=await fetch("/api/enterprise/franchise/commands",{method:"POST",headers:{"content-type":"application/json"},signal:AbortSignal.timeout(10000),body:JSON.stringify({action:"cancel-availability",organizationId:organization,availabilityId:item.id,version:item.version,reasonCode:"schedule-correction"})});
      const receipt=await response.json() as (Partial<Availability>&{organization_id?:unknown})|null;
      if(response.ok&&receipt?.id===item.id&&receipt.organization_id===organization&&receipt.state==="cancelled"&&receipt.version===item.version+1&&receipt.entry_type===item.entry_type&&receipt.resource_id===item.resource_id&&receipt.starts_at===item.starts_at&&receipt.ends_at===item.ends_at) {
        setMessage(t("p0394"));
      } else if(response.status===401||response.status===403)setMessage(t("p0395"));
      else if(response.status===409)setMessage(t("p0396"));
      else throw new Error("unconfirmed");
    } catch {setMessage(t("p0397"))}
  }
  return <>
    <p>{item.state==="active"?t("p0398"):item.state==="cancelled"?t("p0326"):t("p0338")}</p>
    {canCancel?<button type="button" disabled={!ready||locked} onClick={()=>void cancel()}>{t("p0399")}</button>:null}
    <p role="status" aria-label={t("p0400")} aria-live="polite">{message}</p>
    <button type="button" disabled={!ready} onClick={()=>window.location.reload()}>{t("p0401")}</button>
    <OperationalGuide guide={OPERATIONAL_GUIDES["availability-cancel-view"]}/>
  </>;
}

function LeadActions({lead,scope,defaultAssignee,canAssign,canTransition}:{lead:Lead;scope:string;defaultAssignee:string;canAssign:boolean;canTransition:boolean}) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

 const [selections,setSelections]=useState<Record<string,boolean>>({});
  const [ready,setReady]=useState(false),[locked,setLocked]=useState(true),[message,setMessage]=useState("");
  const inFlight=useRef(false);
  const key=`elite-lead-command:${scope}:${lead.id}`;
  const terminal=lead.state==="converted"||lead.state==="lost";
  useEffect(()=>{
    try {
      if(!/^[a-f0-9]{64}$/.test(scope)||!Number.isSafeInteger(lead.version)||lead.version<1)throw new Error("invalid scope/version");
      const raw=window.sessionStorage.getItem(key);
      if(raw!==null) {
        const marker=JSON.parse(raw) as {version?:unknown;action?:unknown};
        if(!marker||typeof marker!=="object"||Object.keys(marker).sort().join(",")!=="action,version"||!Number.isSafeInteger(marker.version)||Number(marker.version)<1||!["assign-lead","transition-lead"].includes(String(marker.action)))throw new Error("invalid marker");
        if(lead.version>Number(marker.version)) {
          window.sessionStorage.removeItem(key);inFlight.current=false;setLocked(terminal);
          setMessage(t("p0402"));
        } else {setLocked(true);setMessage(t("p0403"))}
      } else {setLocked(terminal)}
      setReady(true);
    } catch {setReady(true);setLocked(true);setMessage(t("p0391"))}
  },[key,scope,lead.version,terminal]);
  async function submit(event:FormEvent<HTMLFormElement>,action:"assign-lead"|"transition-lead") {
    event.preventDefault();if(!ready||locked||inFlight.current||terminal)return;
    if(action==="assign-lead"?!canAssign:!canTransition)return;
    const data=new FormData(event.currentTarget),value=String(data.get(action==="assign-lead"?"assignedSubject":"target")??"").trim();
    if(action==="assign-lead"?(!value||value.length>256):!transitions[lead.state]?.includes(value))return;
    inFlight.current=true;setLocked(true);
    try {if(window.sessionStorage.getItem(key)!==null)throw new Error("unresolved");window.sessionStorage.setItem(key,JSON.stringify({version:lead.version,action}))}
    catch {setMessage(t("p0392"));return}
    setMessage(t("p0404"));
    const command={action,organizationId:lead.organization_id,leadId:lead.id,version:lead.version,...(action==="assign-lead"?{assignedSubject:value}:{current:lead.state,target:value})};
    try {
      const response=await fetch("/api/enterprise/franchise/commands",{method:"POST",headers:{"content-type":"application/json"},signal:AbortSignal.timeout(10000),body:JSON.stringify(command)});
      const receipt=await response.json() as Partial<Lead>|null;
      if(definitiveSelectionRejection(response,receipt)){sessionStorage.removeItem(key);inFlight.current=false;setLocked(false);setMessage(experienceText("Las opciones cambiaron antes del envío. No se intentó la operación; actualizá la página para elegir una opción vigente.",privateLocale.language));return;}
      if(response.ok&&receipt?.id===lead.id&&receipt.organization_id===lead.organization_id&&receipt.version===lead.version+1&&(action==="assign-lead"?receipt.assigned_subject===value&&receipt.state===lead.state:receipt.state===value)) {
        setMessage(t("p0405"));
      } else if(response.status===401||response.status===403) {setMessage(t("p0406"))}
      else if(response.status===409) {setMessage(t("p0407"))}
      else throw new Error("unconfirmed");
    } catch {setMessage(t("p0408"))}
  }
  return <>
    <TaskWorkspace title="Trabajo con el contacto">
    {canAssign?<section data-task-label="Responsable"><form onSubmit={event=>void submit(event,"assign-lead")}><GuidedRecordSelect organization={lead.organization_id} group="assignees" onReady={value=>setSelections(old=>({...old,"assignees":value}))} name="assignedSubject" label="Responsable" disabled={!ready||locked}/><button className="primaryAction" disabled={!ready||locked||!selections.assignees}>{t("p0410")}</button></form></section>:null}
    {canTransition&&transitions[lead.state]?.length?<section data-task-label="Estado"><form onSubmit={event=>void submit(event,"transition-lead")}><label>{t("p0411")}<select name="target" disabled={!ready||locked}>{transitions[lead.state]!.map(target=><option key={target} value={target}>{privateLocale.language==="en"?({contacted:"Contacted",qualified:"Qualified",converted:"Converted",lost:"Not converted"} as Record<string,string>)[target]??target:({contacted:"Contactado",qualified:"Calificado",converted:"Convertido",lost:"No concretado"} as Record<string,string>)[target]??target}</option>)}</select></label><button disabled={!ready||locked}>{t("p0412")}</button></form></section>:null}
    </TaskWorkspace>
    <p role="status" aria-label={t("p0413")} aria-live="polite">{message}</p>
    {message ? <button type="button" disabled={!ready} onClick={()=>window.location.reload()}>{experienceText("Actualizar contacto",privateLocale.language)}</button> : null}
    <OperationalGuide guide={OPERATIONAL_GUIDES["lead-command-view"]} summaryLabel="Ayuda"/>
  </>;
}

function DeliveryExceptionResolution({item,organization,scope}:{item:DeliveryException;organization:string;scope:string}) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

  // This owner resolves delivery exceptions, not CRM commands.
  const [ready,setReady]=useState(false);
  const [locked,setLocked]=useState(true);
  const [message,setMessage]=useState("");
  const inFlight=useRef(false);
  const key=`elite-delivery-resolution:${scope}:${item.id}`;
  useEffect(()=>{
    try {
      if(!/^[a-f0-9]{64}$/.test(scope))throw new Error("missing scope");
      const raw=window.sessionStorage.getItem(key);
      if(raw!==null) {
        const marker=JSON.parse(raw) as {version?:unknown;action?:unknown};
        if(!marker||typeof marker!=="object"||Object.keys(marker).sort().join(",")!=="action,version"||!Number.isSafeInteger(marker.version)||Number(marker.version)<1||!["correct-and-represent","return","exchange"].includes(String(marker.action)))throw new Error("invalid marker");
        if(item.state==="resolved"&&Number.isSafeInteger(item.version)&&item.version>Number(marker.version)&&["correct-and-represent","return","exchange"].includes(item.resolution_action??"")) {
          setMessage(item.resolution_action===marker.action?t("p0415"):t("p0416"));
          window.sessionStorage.removeItem(key);
        } else {setMessage(t("p0417"))}
        setLocked(true);
      } else {setLocked(item.state!=="open")}
      setReady(true);
    } catch {setLocked(true);setReady(true);setMessage(t("p0391"))}
  },[key,scope,item.state,item.version,item.resolution_action]);
  async function submit(event:FormEvent<HTMLFormElement>) {
    event.preventDefault();if(!ready||locked||inFlight.current||item.state!=="open")return;
    const data=new FormData(event.currentTarget), action=String(data.get("resolutionAction")??""), notes=String(data.get("notes")??"");
    if(!["correct-and-represent","return","exchange"].includes(action)||!notes.trim()||notes.length>1000)return;
    inFlight.current=true;setLocked(true);
    try {
      if(window.sessionStorage.getItem(key)!==null)throw new Error("unresolved");
      window.sessionStorage.setItem(key,JSON.stringify({version:item.version,action}));
    } catch {setMessage(t("p0392"));return}
    setMessage(t("p0418"));
    try {
      const response=await fetch("/api/enterprise/franchise/commands",{method:"POST",headers:{"content-type":"application/json"},signal:AbortSignal.timeout(10000),body:JSON.stringify({action:"resolve-delivery-exception",organizationId:organization,exceptionId:item.id,version:item.version,resolutionAction:action,notes})});
      const value=await response.json() as {exception?:Record<string,unknown>;successor_handover?:Record<string,unknown>;return_authorization_id?:unknown;disposition?:unknown};
      const e=value?.exception;
      const linked=action==="correct-and-represent" ? typeof e?.successor_handover_id==="string"&&e.successor_handover_id!==""&&value.successor_handover?.id===e.successor_handover_id&&value.successor_handover?.supersedes_handover_id===item.handover_id&&value.successor_handover?.state==="prepared" : typeof e?.return_authorization_id==="string"&&e.return_authorization_id!==""&&value.return_authorization_id===e.return_authorization_id&&value.disposition===action;
      if(response.ok&&e?.id===item.id&&e.organization_id===organization&&e.handover_id===item.handover_id&&e.state==="resolved"&&e.version===item.version+1&&e.resolution_action===action&&linked) {
        setMessage(t("p0419"));
      } else if(response.status===401||response.status===403) {setMessage(t("p0420"))}
      else if(response.status===409) {setMessage(t("p0421"))}
      else throw new Error("unconfirmed");
    } catch {setMessage(t("p0422"))}
  }
  return <>
    {item.state==="open"?<form onSubmit={event=>void submit(event)}><label>{t("p0423")}<select name="resolutionAction" required disabled={!ready||locked}><option value="correct-and-represent">{t("p0424")}</option><option value="return">{t("p0425")}</option><option value="exchange">{t("p0426")}</option></select></label><label>{t("p0427")}<textarea name="notes" required maxLength={1000} disabled={!ready||locked}/></label><button disabled={!ready||locked}>{t("p0428")}</button></form>:<p>{t("p0429")} {item.resolution_action}</p>}
    <p role="status" aria-label={t("p0430")} aria-live="polite">{message}</p>
    <button type="button" disabled={!ready} onClick={()=>window.location.reload()}>{t("p0431")}</button>
    <OperationalGuide guide={OPERATIONAL_GUIDES["delivery-resolution-view"]}/>
  </>;
}

export function FranchiseCommandPanel({ permissions = [], leads, availability, exceptions, returnCases, defaultAssignee, organization, commandScope="" }: { permissions?: readonly string[]; leads: Lead[]; availability: Availability[]; exceptions: DeliveryException[] | null; returnCases: ReturnCase[] | null; defaultAssignee: string; organization: string;commandScope?:string }) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

  const can=(permission:string)=>permissions.includes("*")||permissions.includes(permission);
  return <TaskWorkspace title="Trabajo de gestión" initialLabel="Contactos">
    {(can("availability:read")||can("availability:manage")) ? <section data-task-label="Disponibilidad" className="card"><h2>{t("p0432")}</h2>{can("availability:manage") ? <AvailabilityCreation organization={organization} scope={commandScope}/> : null}<div className="grid">{availability.map((item) => <article className="card" key={item.id}><h3>{item.resource_id || t("p0434")}</h3><p>{item.entry_type} {t("p0016")} {new Date(item.starts_at).toLocaleString(privateLocale.locale, {timeZone:privateLocale.timeZone})} {t("p0435")} {new Date(item.ends_at).toLocaleString(privateLocale.locale, {timeZone:privateLocale.timeZone})}</p>{item.reason_code ? <p>{t("p0436")} {item.reason_code}</p> : null}<AvailabilityCancellation item={item} organization={organization} scope={commandScope} canCancel={can("availability:manage")}/></article>)}</div></section> : null}
    {can("resource:manage") ? <section data-task-label="Recursos" className="card"><h2>{t("p0437")}</h2><ResourceCreation organization={organization} scope={commandScope}/></section> : null}
    {can("appointment:manage") ? <section data-task-label="Turnos" className="card"><h2>{t("p0438")}</h2><SlotCreation organization={organization} scope={commandScope}/></section> : null}
    {can("handover:manage") ? <section data-task-label="Configuración" className="card"><h2>{experienceText("Listas de entrega",privateLocale.language)}</h2><ChecklistPublication key={commandScope} organization={organization} scope={commandScope}/></section> : null}
    {can("handover:manage") && exceptions!==null ? <section data-task-label="Incidencias" className="card"><h2>{t("p0246")}</h2>{exceptions.length ? <div className="grid">{exceptions.map((item) => <article className="card" key={item.id}><h3>{item.reason_code}</h3><p>{item.details}</p><p>{t("p0368")} {item.handover_id} {t("p0442")} {item.customer_subject} {t("p0443")} {item.state}</p><DeliveryExceptionResolution item={item} organization={organization} scope={commandScope}/></article>)}</div> : <p>{t("p0444")}</p>}</section> : null}
    {can("handover:manage") ? <section data-task-label="Devoluciones"><ReturnOperationsPanel key={commandScope} organization={organization} scope={commandScope} cases={returnCases}/></section> : null}
    <div data-task-label="Contactos" className="grid">{leads.map((lead) => <article className="card" key={lead.id}>
      <h2>{lead.created_at&&Number.isFinite(Date.parse(lead.created_at))?experienceText("Contacto del ",privateLocale.language)+new Date(lead.created_at).toLocaleDateString(privateLocale.locale,{day:"numeric",month:"short",timeZone:privateLocale.timeZone}):experienceText("Contacto",privateLocale.language)}</h2><p>{experienceText(({new:"Nuevo",contacted:"Contactado",qualified:"Calificado",converted:"Convertido",lost:"No concretado",closed:"Cerrado"} as Record<string,string>)[lead.state]??lead.state,privateLocale.language)} · {experienceText(lead.source_code==="website"?"Sitio web":lead.source_code,privateLocale.language)}</p><details className="referenceDetails"><summary>{experienceText("Referencia del contacto",privateLocale.language)}</summary><code>{lead.id}</code></details><p>{lead.assigned_subject ? `Responsable: ${lead.assigned_subject}` : t("p0378")}</p>
      <details className="contactWork" name="franchise-contact"><summary>{experienceText("Gestionar contacto",privateLocale.language)}</summary><TaskWorkspace title="Trabajo con el contacto"><section data-task-label="Seguimiento"><LeadActions lead={lead} scope={commandScope} defaultAssignee={defaultAssignee} canAssign={can("lead:assign")} canTransition={can("lead:update")}/></section>
      {can("quote:write") ? <section data-task-label="Cotización"><QuoteCreation lead={lead} scope={commandScope}/></section> : null}</TaskWorkspace></details>
    </article>)}</div>
  </TaskWorkspace>;
}

function QuoteCreation({lead,scope}:{lead:Lead;scope:string}) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

 const [selections,setSelections]=useState<Record<string,boolean>>({});
 const [ready,setReady]=useState(false),[locked,setLocked]=useState(true),[message,setMessage]=useState("");
 const [result,setResult]=useState<{id:string;state:string;currency:string;total_minor_units:number}|null>(null);
 const sending=useRef(false),reading=useRef(false),operation=useRef<string|null>(null);
 const storageKey=`elite-quote-create:${scope}:${lead.id}`;
 useEffect(()=>{
   try {
     if(!/^[a-f0-9]{64}$/.test(scope))throw new Error("scope");
     const raw=sessionStorage.getItem(storageKey);
     if(raw!==null){const marker=JSON.parse(raw);if(!marker||Object.keys(marker).join(",")!=="requestKey"||!/^quote-[a-f0-9-]{36}$/.test(marker.requestKey))throw new Error("marker");operation.current=marker.requestKey;setLocked(true);setMessage(t("p0445"));}
     else {operation.current=null;setLocked(false)}
     setReady(true);
   }catch{setReady(true);setLocked(true);setMessage(t("p0446"))}
 },[scope,storageKey]);
 function valid(value:unknown): value is {id:string;organization_id:string;lead_id:string;variant_id:string;price_book_id:string;valid_until:string;state:string;currency:string;total_minor_units:number;version:number} {
   if(!value||typeof value!=="object")return false;
   const v=value as Record<string,unknown>;
   return typeof v.id==="string"&&v.id.length>0&&v.id.length<=128&&v.organization_id===lead.organization_id&&v.lead_id===lead.id&&typeof v.variant_id==="string"&&typeof v.price_book_id==="string"&&typeof v.valid_until==="string"&&Number.isFinite(Date.parse(v.valid_until))&&["issued","accepted","expired","cancelled"].includes(String(v.state))&&typeof v.currency==="string"&&/^[A-Z]{3}$/.test(v.currency)&&Number.isSafeInteger(v.total_minor_units)&&Number(v.total_minor_units)>=0&&Number.isSafeInteger(v.version)&&Number(v.version)>0;
 }
 async function submit(event:FormEvent<HTMLFormElement>){
   event.preventDefault();if(!ready||locked||sending.current||!selections.variants||!selections.priceBooks)return;
   const data=new FormData(event.currentTarget),variantId=String(data.get("variantId")??""),priceBookId=String(data.get("priceBookId")??""),until=new Date(String(data.get("validUntil")??""));
   if(!variantId||!priceBookId||!Number.isFinite(until.getTime())||until<=new Date())return;
   sending.current=true;setLocked(true);
   try{if(sessionStorage.getItem(storageKey)!==null)throw new Error("unresolved");const requestKey="quote-"+crypto.randomUUID();sessionStorage.setItem(storageKey,JSON.stringify({requestKey}));operation.current=requestKey;}
   catch{setMessage(t("p0447"));return}
   setMessage(t("p0448"));
   try{
     const response=await fetch("/api/enterprise/franchise/commands",{method:"POST",headers:{"content-type":"application/json","idempotency-key":operation.current!},signal:AbortSignal.timeout(10000),body:JSON.stringify({action:"create-quote",organizationId:lead.organization_id,leadId:lead.id,variantId,priceBookId,validUntil:until.toISOString()})});
     const receipt=await response.json();
     if(definitiveSelectionRejection(response,receipt)){sessionStorage.removeItem(storageKey);operation.current=null;sending.current=false;setLocked(false);setMessage(experienceText("Las opciones cambiaron antes del envío. No se intentó la operación; actualizá la página para elegir una opción vigente.",privateLocale.language));return;}
     if(!response.ok||!valid(receipt)||receipt.variant_id!==variantId||receipt.price_book_id!==priceBookId||new Date(receipt.valid_until).getTime()!==until.getTime()||receipt.state!=="issued"||receipt.version!==1)throw new Error("unconfirmed");
     setResult(receipt);setMessage(t("p0449"));
   }catch{setMessage(t("p0450"))}
 }
 async function consult(){
   if(!ready||!operation.current||reading.current)return;reading.current=true;
   try{
     const query=new URLSearchParams({organizationId:lead.organization_id,leadId:lead.id,requestKey:operation.current});
     const response=await fetch("/api/enterprise/franchise/commands?"+query,{cache:"no-store",signal:AbortSignal.timeout(10000)});
     const receipt=await response.json();
     if(!response.ok||receipt?.request_key!==operation.current||!valid(receipt.quote))throw new Error("unconfirmed");
     setResult(receipt.quote);setMessage(t("p0451"));
   }catch{setMessage(t("p0452"))}
   finally{reading.current=false}
 }
 return <div>
  <form onSubmit={event=>void submit(event)}>
   <GuidedRecordSelect organization={lead.organization_id} group="variants" onReady={value=>setSelections(old=>({...old,"variants":value}))} name="variantId" label="Producto" disabled={!ready||locked}/>
   <GuidedRecordSelect organization={lead.organization_id} group="priceBooks" onReady={value=>setSelections(old=>({...old,"priceBooks":value}))} name="priceBookId" label="Lista de precios" disabled={!ready||locked}/>
   <label>{t("p0454")}<input name="validUntil" type="datetime-local" required disabled={!ready||locked}/></label>
   <button disabled={!ready||locked||!selections.variants||!selections.priceBooks}>{t("p0455")}</button>
  </form>
  <p role="status" aria-label={t("p0456")} aria-live="polite">{message}</p>
  {result?<p>{t("p0457")} {result.id} {t("p0016")} {result.state} {t("p0016")} {result.currency}</p>:null}
  <button type="button" disabled={!ready||!operation.current} onClick={()=>void consult()}>{t("p0458")}</button>
  <OperationalGuide guide={OPERATIONAL_GUIDES["quote-create-view"]}/>
 </div>;
}

function AvailabilityCreation({organization,scope}:{organization:string;scope:string}) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

 const [selections,setSelections]=useState<Record<string,boolean>>({});
 const [ready,setReady]=useState(false),[locked,setLocked]=useState(true),[message,setMessage]=useState("");
 const [confirmed,setConfirmed]=useState(false),[posting,setPosting]=useState(false),[result,setResult]=useState<Availability|null>(null);
 const fence=useRef(false),reading=useRef(false),operation=useRef<string|null>(null),form=useRef<HTMLFormElement>(null);
 const storageKey=`elite-availability-create:${scope}`;
 useEffect(()=>{
  try{
   if(!/^[a-f0-9]{64}$/.test(scope))throw new Error("scope");
   const raw=sessionStorage.getItem(storageKey);
   if(raw!==null){const marker=JSON.parse(raw);if(!marker||Object.keys(marker).join(",")!=="requestKey"||!/^availability-[a-f0-9-]{36}$/.test(marker.requestKey))throw new Error("marker");operation.current=marker.requestKey;setLocked(true);setMessage(t("p0459"));}
   else{operation.current=null;setLocked(false)}setReady(true);
  }catch{setReady(true);setLocked(true);setMessage(t("p0446"))}
 },[scope,storageKey]);
 function valid(value:unknown):value is Availability & {organization_id:string}{
  if(!value||typeof value!=="object")return false;const v=value as Record<string,unknown>;
  return typeof v.id==="string"&&v.id.length>0&&v.id.length<=128&&v.organization_id===organization&&["working","unavailable"].includes(String(v.entry_type))&&["active","cancelled"].includes(String(v.state))&&Number.isSafeInteger(v.version)&&Number(v.version)>0&&typeof v.starts_at==="string"&&typeof v.ends_at==="string"&&Number.isFinite(Date.parse(v.starts_at))&&Date.parse(v.ends_at)>Date.parse(v.starts_at)&&(v.resource_id===undefined||typeof v.resource_id==="string")&&(v.reason_code===undefined||typeof v.reason_code==="string");
 }
 async function submit(event:FormEvent<HTMLFormElement>){
  event.preventDefault();if(!ready||locked||fence.current||!selections.resources||!selections.availabilityReasons)return;
  const data=new FormData(event.currentTarget),from=new Date(String(data.get("startsAt")??"")),to=new Date(String(data.get("endsAt")??""));
  if(!Number.isFinite(from.getTime())||!Number.isFinite(to.getTime())||to<=from){setMessage(t("p0460"));return}
  const resourceId=String(data.get("resourceId")??"").trim(),entryType=String(data.get("entryType")??"working"),reasonCode=String(data.get("reasonCode")??"").trim();
  fence.current=true;setLocked(true);setConfirmed(false);
  try{if(sessionStorage.getItem(storageKey)!==null)throw new Error("unresolved");const key="availability-"+crypto.randomUUID();sessionStorage.setItem(storageKey,JSON.stringify({requestKey:key}));operation.current=key;}
  catch{setMessage(t("p0447"));return}
  const key=operation.current;setPosting(true);setMessage(t("p0461"));
  try{
   const response=await fetch("/api/enterprise/franchise/commands",{method:"POST",headers:{"content-type":"application/json","idempotency-key":key!},signal:AbortSignal.timeout(10000),body:JSON.stringify({action:"create-availability",organizationId:organization,...(resourceId?{resourceId}:{}),entryType,...(reasonCode?{reasonCode}:{}),startsAt:from.toISOString(),endsAt:to.toISOString()})});
   const receipt=await response.json();if(operation.current!==key)return;
   if(definitiveSelectionRejection(response,receipt)){sessionStorage.removeItem(storageKey);operation.current=null;fence.current=false;setLocked(false);setMessage(experienceText("Las opciones cambiaron antes del envío. No se intentó la operación; actualizá la página para elegir una opción vigente.",privateLocale.language));return;}
   if(!response.ok||!valid(receipt)||(receipt.resource_id??"")!==resourceId||receipt.entry_type!==entryType||(receipt.reason_code??"")!==reasonCode||Date.parse(receipt.starts_at)!==from.getTime()||Date.parse(receipt.ends_at)!==to.getTime()||receipt.state!=="active"||receipt.version!==1)throw new Error("unconfirmed");
   setResult(receipt);setMessage(t("p0462"));
  }catch{if(operation.current===key)setMessage(t("p0463"))}
  finally{setPosting(false)}
 }
 async function consult(){
  if(!ready||!operation.current||reading.current)return;reading.current=true;const key=operation.current;
  try{
   const query=new URLSearchParams({kind:"availability",organizationId:organization,requestKey:key});
   const response=await fetch("/api/enterprise/franchise/commands?"+query,{cache:"no-store",signal:AbortSignal.timeout(10000)}),receipt=await response.json();
   if(operation.current!==key)return;
   if(!response.ok||receipt?.request_key!==key||!valid(receipt.availability))throw new Error("unconfirmed");
   setResult(receipt.availability);setConfirmed(true);setMessage(t("p0464"));
  }catch{if(operation.current===key){setConfirmed(false);setMessage(t("p0465"))}}
  finally{reading.current=false}
 }
 function prepare(){
  if(!confirmed||posting||reading.current)return;
  try{sessionStorage.removeItem(storageKey);operation.current=null;fence.current=false;setLocked(false);setConfirmed(false);setResult(null);form.current?.reset();setMessage(t("p0466"))}
  catch{setMessage(t("p0467"))}
 }
 return <div>
  <form ref={form} onSubmit={event=>void submit(event)}>
   <GuidedRecordSelect organization={organization} group="resources" onReady={value=>setSelections(old=>({...old,"resources":value}))} name="resourceId" label="Persona o recurso" disabled={!ready||locked} optional/>
   <label>{t("p0470")}<select name="entryType" defaultValue="working" disabled={!ready||locked}><option value="working">{t("p0471")}</option><option value="unavailable">{t("p0472")}</option></select></label>
   <GuidedRecordSelect organization={organization} group="availabilityReasons" onReady={value=>setSelections(old=>({...old,"availabilityReasons":value}))} name="reasonCode" label="Motivo" disabled={!ready||locked} optional/>
   <label>{t("p0475")}<input name="startsAt" type="datetime-local" required disabled={!ready||locked}/></label><label>{t("p0476")}<input name="endsAt" type="datetime-local" required disabled={!ready||locked}/></label>
   <button disabled={!ready||locked||!selections.resources||!selections.availabilityReasons}>{t("p0477")}</button>
  </form>
  <p role="status" aria-label={t("p0478")} aria-live="polite">{message}</p>
  {result?<p>{t("p0479")} {result.id} {t("p0016")} {result.state} {t("p0016")} {new Date(result.starts_at).toLocaleString(privateLocale.locale, {timeZone:privateLocale.timeZone})} {t("p0435")} {new Date(result.ends_at).toLocaleString(privateLocale.locale, {timeZone:privateLocale.timeZone})}</p>:null}
  <button type="button" disabled={!ready||!operation.current} onClick={()=>void consult()}>{t("p0480")}</button>
  <button type="button" disabled={!confirmed||posting} onClick={prepare}>{t("p0481")}</button>
  <OperationalGuide guide={OPERATIONAL_GUIDES["availability-create-view"]}/>
 </div>;
}

type ResourceReceipt={id:string;organization_id:string;display_name:string;principal_subject?:string;kind:string;status:string;version:number;skills:string[]};
function ResourceCreation({organization,scope}:{organization:string;scope:string}) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

 const [ready,setReady]=useState(false),[locked,setLocked]=useState(true),[message,setMessage]=useState("");
 const [confirmed,setConfirmed]=useState(false),[posting,setPosting]=useState(false),[result,setResult]=useState<ResourceReceipt|null>(null);
 const fence=useRef(false),reading=useRef(false),operation=useRef<string|null>(null),form=useRef<HTMLFormElement>(null);
 const storageKey=`elite-resource-create:${scope}`;
 useEffect(()=>{
  try{
   if(!/^[a-f0-9]{64}$/.test(scope))throw new Error("scope");
   const raw=sessionStorage.getItem(storageKey);
   if(raw!==null){const marker=JSON.parse(raw);if(!marker||Object.keys(marker).join(",")!=="requestKey"||!/^resource-[a-f0-9-]{36}$/.test(marker.requestKey))throw new Error("marker");operation.current=marker.requestKey;setLocked(true);setMessage(t("p0482"));}
   else{operation.current=null;setLocked(false)}setReady(true);
  }catch{setReady(true);setLocked(true);setMessage(t("p0446"))}
 },[scope,storageKey]);
 function valid(value:unknown):value is ResourceReceipt {
  if(!value||typeof value!=="object")return false;const v=value as Record<string,unknown>;
  return typeof v.id==="string"&&v.id.length>0&&v.id.length<=128&&v.organization_id===organization&&typeof v.display_name==="string"&&v.display_name.length>0&&v.display_name.length<=160&&["employee","contractor","service-bay","vehicle"].includes(String(v.kind))&&["active","inactive"].includes(String(v.status))&&Number.isSafeInteger(v.version)&&Number(v.version)>0&&(v.principal_subject===undefined||typeof v.principal_subject==="string")&&Array.isArray(v.skills)&&v.skills.length>0&&v.skills.length<=4&&new Set(v.skills).size===v.skills.length&&v.skills.every(skill=>typeof skill==="string"&&["consultation","test-drive","delivery","service"].includes(skill));
 }
 async function submit(event:FormEvent<HTMLFormElement>){
  event.preventDefault();if(!ready||locked||fence.current)return;
  const data=new FormData(event.currentTarget),displayName=String(data.get("displayName")??"").trim(),principalSubject=String(data.get("principalSubject")??"").trim(),kind=String(data.get("kind")??"employee"),skills=data.getAll("skills").map(String);
  if(!displayName||!skills.length||((kind==="employee"||kind==="contractor")!==!!principalSubject)){setMessage(t("p0483"));return}
  fence.current=true;setLocked(true);setConfirmed(false);
  try{if(sessionStorage.getItem(storageKey)!==null)throw new Error("unresolved");const key="resource-"+crypto.randomUUID();sessionStorage.setItem(storageKey,JSON.stringify({requestKey:key}));operation.current=key;}
  catch{setMessage(t("p0447"));return}
  const key=operation.current;setPosting(true);setMessage(t("p0484"));
  try{
   const response=await fetch("/api/enterprise/franchise/commands",{method:"POST",headers:{"content-type":"application/json","idempotency-key":key!},signal:AbortSignal.timeout(10000),body:JSON.stringify({action:"create-resource",organizationId:organization,displayName,...(principalSubject?{principalSubject}:{}),kind,skills})});
   const receipt=await response.json();if(operation.current!==key)return;
   if(!response.ok||!valid(receipt)||receipt.display_name!==displayName||(receipt.principal_subject??"")!==principalSubject||receipt.kind!==kind||JSON.stringify([...receipt.skills].sort())!==JSON.stringify([...skills].sort())||receipt.status!=="active"||receipt.version!==1)throw new Error("unconfirmed");
   setResult(receipt);setMessage(t("p0485"));
  }catch{if(operation.current===key)setMessage(t("p0486"))}
  finally{setPosting(false)}
 }
 async function consult(){
  if(!ready||!operation.current||reading.current)return;reading.current=true;const key=operation.current;
  try{
   const query=new URLSearchParams({kind:"resource",organizationId:organization,requestKey:key});
   const response=await fetch("/api/enterprise/franchise/commands?"+query,{cache:"no-store",signal:AbortSignal.timeout(10000)}),receipt=await response.json();
   if(operation.current!==key)return;
   if(!response.ok||receipt?.request_key!==key||!valid(receipt.resource))throw new Error("unconfirmed");
   setResult(receipt.resource);setConfirmed(true);setMessage(t("p0487"));
  }catch{if(operation.current===key){setConfirmed(false);setMessage(t("p0488"))}}
  finally{reading.current=false}
 }
 function prepare(){
  if(!confirmed||posting||reading.current)return;
  try{sessionStorage.removeItem(storageKey);operation.current=null;fence.current=false;setLocked(false);setConfirmed(false);setResult(null);form.current?.reset();setMessage(t("p0489"))}
  catch{setMessage(t("p0467"))}
 }
 return <div>
  <form ref={form} onSubmit={event=>void submit(event)}>
   <label>{t("p0490")}<input name="displayName" required maxLength={160} disabled={!ready||locked}/></label>
   <GuidedRecordSelect organization={organization} group="resourcePrincipals" name="principalSubject" label="Persona asociada al recurso" disabled={!ready||locked} optional/>
   <label>{t("p0470")}<select name="kind" defaultValue="employee" disabled={!ready||locked}><option value="employee">{t("p0492")}</option><option value="contractor">{t("p0493")}</option><option value="service-bay">{t("p0494")}</option><option value="vehicle">{t("p0495")}</option></select></label>
   <fieldset disabled={!ready||locked}><legend>{t("p0496")}</legend>{["consultation","test-drive","delivery","service"].map(skill=><label key={skill}><input name="skills" type="checkbox" value={skill}/>{skill}</label>)}</fieldset>
   <button disabled={!ready||locked}>{t("p0497")}</button>
  </form>
  <p role="status" aria-label={t("p0498")} aria-live="polite">{message}</p>
  {result?<p>{t("p0499")} {result.id} {t("p0016")} {result.display_name} {t("p0016")} {result.kind} {t("p0016")} {result.status}</p>:null}
  <button type="button" disabled={!ready||!operation.current} onClick={()=>void consult()}>{t("p0500")}</button>
  <button type="button" disabled={!confirmed||posting} onClick={prepare}>{t("p0501")}</button>
  <OperationalGuide guide={OPERATIONAL_GUIDES["resource-create-view"]}/>
 </div>;
}

type SlotReceipt={id:string;organization_id:string;kind:string;starts_at:string;ends_at:string;capacity:number;booked:number;state:string;version:number};
function SlotCreation({organization,scope}:{organization:string;scope:string}) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

 const [ready,setReady]=useState(false),[locked,setLocked]=useState(true),[message,setMessage]=useState("");
 const [confirmed,setConfirmed]=useState(false),[posting,setPosting]=useState(false),[result,setResult]=useState<SlotReceipt|null>(null);
 const fence=useRef(false),reading=useRef(false),operation=useRef<string|null>(null),form=useRef<HTMLFormElement>(null);
 const storageKey=`elite-slot-create:${scope}`;
 useEffect(()=>{
  try{
   if(!/^[a-f0-9]{64}$/.test(scope))throw new Error("scope");
   const raw=sessionStorage.getItem(storageKey);
   if(raw!==null){const marker=JSON.parse(raw);if(!marker||Object.keys(marker).join(",")!=="requestKey"||!/^slot-[a-f0-9-]{36}$/.test(marker.requestKey))throw new Error("marker");operation.current=marker.requestKey;setLocked(true);setMessage(t("p0502"));}
   else{operation.current=null;setLocked(false)}setReady(true);
  }catch{setReady(true);setLocked(true);setMessage(t("p0446"))}
 },[scope,storageKey]);
 function valid(value:unknown):value is SlotReceipt {
  if(!value||typeof value!=="object")return false;const v=value as Record<string,unknown>;
  return typeof v.id==="string"&&v.id.length>0&&v.id.length<=128&&v.organization_id===organization&&["consultation","test-drive","delivery","service"].includes(String(v.kind))&&["open","closed"].includes(String(v.state))&&Number.isSafeInteger(v.version)&&Number(v.version)>0&&typeof v.starts_at==="string"&&typeof v.ends_at==="string"&&Number.isFinite(Date.parse(v.starts_at))&&Date.parse(v.ends_at)>Date.parse(v.starts_at)&&Number.isSafeInteger(v.capacity)&&Number(v.capacity)>=1&&Number(v.capacity)<=100&&Number.isSafeInteger(v.booked)&&Number(v.booked)>=0&&Number(v.booked)<=Number(v.capacity);
 }
 async function submit(event:FormEvent<HTMLFormElement>){
  event.preventDefault();if(!ready||locked||fence.current)return;
  const data=new FormData(event.currentTarget),from=new Date(String(data.get("startsAt")??"")),to=new Date(String(data.get("endsAt")??""));
  if(!Number.isFinite(from.getTime())||!Number.isFinite(to.getTime())||to<=from){setMessage(t("p0503"));return}
  const kind=String(data.get("kind")??"test-drive"),capacity=Number(data.get("capacity"));
  if(!Number.isInteger(capacity)||capacity<1||capacity>100){setMessage(t("p0504"));return}
  fence.current=true;setLocked(true);setConfirmed(false);
  try{if(sessionStorage.getItem(storageKey)!==null)throw new Error("unresolved");const key="slot-"+crypto.randomUUID();sessionStorage.setItem(storageKey,JSON.stringify({requestKey:key}));operation.current=key;}
  catch{setMessage(t("p0447"));return}
  const key=operation.current;setPosting(true);setMessage(t("p0505"));
  try{
   const response=await fetch("/api/enterprise/franchise/commands",{method:"POST",headers:{"content-type":"application/json","idempotency-key":key!},signal:AbortSignal.timeout(10000),body:JSON.stringify({action:"create-appointment-slot",organizationId:organization,kind,capacity,startsAt:from.toISOString(),endsAt:to.toISOString()})});
   const receipt=await response.json();if(operation.current!==key)return;
   if(!response.ok||!valid(receipt)||receipt.kind!==kind||receipt.capacity!==capacity||Date.parse(receipt.starts_at)!==from.getTime()||Date.parse(receipt.ends_at)!==to.getTime()||receipt.state!=="open"||receipt.version!==1)throw new Error("unconfirmed");
   setResult(receipt);setMessage(t("p0506"));
  }catch{if(operation.current===key)setMessage(t("p0507"))}
  finally{setPosting(false)}
 }
 async function consult(){
  if(!ready||!operation.current||reading.current)return;reading.current=true;const key=operation.current;
  try{
   const query=new URLSearchParams({kind:"slot",organizationId:organization,requestKey:key});
   const response=await fetch("/api/enterprise/franchise/commands?"+query,{cache:"no-store",signal:AbortSignal.timeout(10000)}),receipt=await response.json();
   if(operation.current!==key)return;
   if(!response.ok||receipt?.request_key!==key||!valid(receipt.slot))throw new Error("unconfirmed");
   setResult(receipt.slot);setConfirmed(true);setMessage(t("p0508"));
  }catch{if(operation.current===key){setConfirmed(false);setMessage(t("p0509"))}}
  finally{reading.current=false}
 }
 function prepare(){
  if(!confirmed||posting||reading.current)return;
  try{sessionStorage.removeItem(storageKey);operation.current=null;fence.current=false;setLocked(false);setConfirmed(false);setResult(null);form.current?.reset();setMessage(t("p0510"))}
  catch{setMessage(t("p0467"))}
 }
 return <div>
  <form ref={form} onSubmit={event=>void submit(event)}>
   <label>{t("p0470")}<select name="kind" defaultValue="test-drive" disabled={!ready||locked}><option value="test-drive">{t("p0367")}</option><option value="consultation">{t("p0366")}</option><option value="delivery">{t("p0368")}</option><option value="service">{t("p0018")}</option></select></label>
   <label>{t("p0475")}<input name="startsAt" type="datetime-local" required disabled={!ready||locked}/></label><label>{t("p0476")}<input name="endsAt" type="datetime-local" required disabled={!ready||locked}/></label>
   <label>{t("p0511")}<input name="capacity" type="number" min={1} max={100} defaultValue={1} required disabled={!ready||locked}/></label>
   <button disabled={!ready||locked}>{t("p0512")}</button>
  </form>
  <p role="status" aria-label={t("p0513")} aria-live="polite">{message}</p>
  {result?<p>{t("p0140")} {result.id} {t("p0016")} {result.state} {t("p0016")} {result.booked}{t("p0280")}{result.capacity} {t("p0514")} {new Date(result.starts_at).toLocaleString(privateLocale.locale, {timeZone:privateLocale.timeZone})} {t("p0435")} {new Date(result.ends_at).toLocaleString(privateLocale.locale, {timeZone:privateLocale.timeZone})}</p>:null}
  <button type="button" disabled={!ready||!operation.current} onClick={()=>void consult()}>{t("p0515")}</button>
  <button type="button" disabled={!confirmed||posting} onClick={prepare}>{t("p0516")}</button>
  <OperationalGuide guide={OPERATIONAL_GUIDES["slot-create-view"]}/>
 </div>;
}



type PublishedChecklist={id:string;organization_id:string;version:number;title:string;state:string;items:{id:string;ordinal:number;prompt:string;response_type:string;required:boolean}[]};
type ChecklistMarker={checklistId:string;version:number;digest:string};
function ChecklistPublication({organization,scope}:{organization:string;scope:string}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

 const [ready,setReady]=useState(false),[locked,setLocked]=useState(true),[message,setMessage]=useState("");
 const [posting,setPosting]=useState(false),[querying,setQuerying]=useState(false),[confirmed,setConfirmed]=useState(false),[result,setResult]=useState<PublishedChecklist|null>(null);
 const [checklistItems,setChecklistItems]=useState<ChecklistDraftItem[]>([{key:1,id:"",prompt:"",responseType:"confirmation",required:true}]);
 const form=useRef<HTMLFormElement>(null),fence=useRef(false),reading=useRef(false),operation=useRef<ChecklistMarker|null>(null);
 const storageKey=`elite-checklist-publish:${scope}`;
 const code=(v:unknown):v is string=>typeof v==="string"&&v.length<=64&&/^[a-z][a-z0-9]*(-[a-z0-9]+)*$/.test(v);
 const textWithin=(v:unknown,max:number):v is string=>typeof v==="string"&&v.trim().length>0&&new TextEncoder().encode(v).length<=max;
 function valid(v:unknown):v is PublishedChecklist{
  if(!v||typeof v!=="object")return false;const value=v as PublishedChecklist;
  return code(value.id)&&value.organization_id===organization&&Number.isSafeInteger(value.version)&&value.version>0&&value.state==="published"&&textWithin(value.title,160)&&Array.isArray(value.items)&&value.items.length>0&&value.items.length<=64&&new Set(value.items.map(i=>i?.id)).size===value.items.length&&value.items.every((i,index)=>i&&code(i.id)&&i.ordinal===index+1&&textWithin(i.prompt,500)&&["confirmation","text","serial","evidence"].includes(i.response_type)&&typeof i.required==="boolean");
 }
 async function digest(value:PublishedChecklist){
  const bytes=new TextEncoder().encode(JSON.stringify({id:value.id,organization_id:value.organization_id,version:value.version,title:value.title,items:value.items.map(i=>({id:i.id,prompt:i.prompt,response_type:i.response_type,required:i.required}))}));
  return Array.from(new Uint8Array(await crypto.subtle.digest("SHA-256",bytes)),b=>b.toString(16).padStart(2,"0")).join("");
 }
 useEffect(()=>{
  try{
   if(!/^[a-f0-9]{64}$/.test(scope))throw new Error("scope");
   const raw=sessionStorage.getItem(storageKey);
   if(raw!==null){const v=JSON.parse(raw);if(!v||Object.keys(v).sort().join(",")!=="checklistId,digest,version"||!code(v.checklistId)||!Number.isSafeInteger(v.version)||v.version<1||!/^[a-f0-9]{64}$/.test(v.digest))throw new Error("marker");operation.current=v;setLocked(true);setMessage(t("p0517"));}
   else{operation.current=null;setLocked(false)}setReady(true);
  }catch{setReady(true);setLocked(true);setMessage(t("p0518"))}
 },[scope,storageKey]);
 async function submit(event:FormEvent<HTMLFormElement>){
  event.preventDefault();if(!ready||locked||fence.current)return;
  const data=new FormData(event.currentTarget);const value:PublishedChecklist={id:String(data.get("checklistId")??""),organization_id:organization,version:Number(data.get("checklistVersion")),title:String(data.get("title")??"").trim(),state:"published",items:checklistItems.map((i,index)=>({id:i.id,ordinal:index+1,prompt:i.prompt.trim(),response_type:i.responseType,required:i.required}))};
  if(!valid(value)){setMessage(t("p0519"));return}
  fence.current=true;setLocked(true);setConfirmed(false);setPosting(true);
  let marker:ChecklistMarker;
  try{if(sessionStorage.getItem(storageKey)!==null)throw new Error("unresolved");marker={checklistId:value.id,version:value.version,digest:await digest(value)};sessionStorage.setItem(storageKey,JSON.stringify(marker));operation.current=marker;}
  catch{setPosting(false);setMessage(t("p0520"));return}
  setMessage(t("p0521"));
  try{
   const response=await fetch("/api/enterprise/franchise/commands",{method:"POST",headers:{"content-type":"application/json"},signal:AbortSignal.timeout(10000),body:JSON.stringify({action:"publish-delivery-checklist",organizationId:organization,checklistId:value.id,version:value.version,title:value.title,items:value.items.map(i=>({id:i.id,prompt:i.prompt,response_type:i.response_type,required:i.required}))})});
   const receipt=await response.json();if(operation.current!==marker)return;
   if(!response.ok||!valid(receipt)||await digest(receipt)!==marker.digest)throw new Error("unconfirmed");
   setResult(receipt);setMessage(t("p0522"));
  }catch{if(operation.current===marker)setMessage(t("p0523"))}
  finally{setPosting(false)}
 }
 async function consult(){
  const marker=operation.current;if(!ready||!marker||reading.current)return;reading.current=true;setQuerying(true);
  try{
   const q=new URLSearchParams({kind:"checklist",organizationId:organization,checklistId:marker.checklistId,version:String(marker.version)});
   const response=await fetch("/api/enterprise/franchise/commands?"+q,{cache:"no-store",signal:AbortSignal.timeout(10000)}),body=await response.json();
   if(operation.current!==marker)return;
   if(!response.ok||!valid(body?.checklist)||body.checklist.id!==marker.checklistId||body.checklist.version!==marker.version||await digest(body.checklist)!==marker.digest)throw new Error("unconfirmed");
   setResult(body.checklist);setConfirmed(true);setMessage(t("p0524"));
  }catch{if(operation.current===marker){setConfirmed(false);setMessage(t("p0525"))}}
  finally{reading.current=false;setQuerying(false)}
 }
 function prepare(){
  if(!confirmed||posting||reading.current)return;
  try{sessionStorage.removeItem(storageKey);operation.current=null;fence.current=false;setLocked(false);setConfirmed(false);setResult(null);form.current?.reset();setChecklistItems([{key:1,id:"",prompt:"",responseType:"confirmation",required:true}]);setMessage(t("p0526"))}
  catch{setMessage(t("p0527"))}
 }
 return <div><h3>{t("p0528")}</h3><form ref={form} onSubmit={event=>void submit(event)}><fieldset disabled={!ready||locked}>
 
      <label>{t("p0529")}<input name="checklistId" required maxLength={64} pattern="[a-z][a-z0-9]*(-[a-z0-9]+)*" /></label><label>{t("p0530")}<input name="checklistVersion" type="number" min={1} required /></label><label>{t("p0531")}<input name="title" required maxLength={160} /></label>
      {checklistItems.map((item, index) => <fieldset key={item.key}><legend>{t("p0532")} {index + 1}</legend><label>{t("p0533")}<input required value={item.id} maxLength={64} pattern="[a-z][a-z0-9]*(-[a-z0-9]+)*" onChange={(event) => setChecklistItems((items) => items.map((candidate) => candidate.key === item.key ? { ...candidate, id: event.target.value } : candidate))}/></label><label>{t("p0534")}<input required value={item.prompt} maxLength={500} onChange={(event) => setChecklistItems((items) => items.map((candidate) => candidate.key === item.key ? { ...candidate, prompt: event.target.value } : candidate))}/></label><label>{t("p0535")}<select value={item.responseType} onChange={(event) => setChecklistItems((items) => items.map((candidate) => candidate.key === item.key ? { ...candidate, responseType: event.target.value as ChecklistDraftItem["responseType"] } : candidate))}><option value="confirmation">{t("p0141")}</option><option value="text">{t("p0536")}</option><option value="serial">{t("p0537")}</option></select></label><label><input type="checkbox" checked={item.required} onChange={(event) => setChecklistItems((items) => items.map((candidate) => candidate.key === item.key ? { ...candidate, required: event.target.checked } : candidate))}/> {t("p0538")}</label>{checklistItems.length > 1 ? <button type="button" onClick={() => setChecklistItems((items) => items.filter((candidate) => candidate.key !== item.key))}>{t("p0539")}</button> : null}</fieldset>)}
      <button type="button" disabled={checklistItems.length >= 64} onClick={() => setChecklistItems((items) => [...items, { key: Math.max(...items.map((item) => item.key)) + 1, id: "", prompt: "", responseType: "confirmation", required: true }])}>{t("p0540")}</button><button disabled={!ready||locked}>{t("p0541")}</button>
    
 </fieldset></form>
 <p role="status" aria-label={t("p0542")} aria-live="polite">{message}</p>
 {result&&result.organization_id===organization?<div><p>{result.id} {t("p0543")} {result.version} {t("p0016")} {result.title} {t("p0016")} {result.state}</p><ol>{result.items.map(item=><li key={item.id}>{item.id}{t("p0213")} {item.prompt} {t("p0016")} {item.response_type} {t("p0016")} {item.required?t("p0538"):t("p0544")}</li>)}</ol></div>:null}
 <button type="button" disabled={!ready||!operation.current||querying} onClick={()=>void consult()}>{t("p0545")}</button>
 <button type="button" disabled={!confirmed||posting||querying} onClick={prepare}>{t("p0546")}</button>
 <OperationalGuide guide={OPERATIONAL_GUIDES["checklist-publication-view"]}/>
 </div>;
}


type CompletionReference={handoverId:string;version:number;checklistId:string;checklistVersion:number;digest:string};
type CompletionResponse={item_id:string;response_text:string;evidence_sha256?:string};
type CompletionResult={handover_id:string;organization_id:string;state:string;version:number;checklist_id:string;checklist_version:number;completed_at:string;actor_subject:string;responses:CompletionResponse[]};
export function ChecklistCompletionPanel({organization,scope,initialHandover}:{organization:string;scope:string;initialHandover?:{id:string;version:number}}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

 const [guidedReady,setGuidedReady]=useState(false);
 const [ready,setReady]=useState(false),[locked,setLocked]=useState(true),[message,setMessage]=useState("");
 const [posting,setPosting]=useState(false),[querying,setQuerying]=useState(false),[confirmed,setConfirmed]=useState(false),[result,setResult]=useState<CompletionResult|null>(null);
 const form=useRef<HTMLFormElement>(null),fence=useRef(false),reading=useRef(false),operation=useRef<CompletionReference|null>(null);
 const storageKey=`elite-checklist-complete:${scope}${initialHandover ? ":"+initialHandover.id : ""}`;
 const id=(v:unknown):v is string=>typeof v==="string"&&/^[A-Za-z0-9_-]{1,128}$/.test(v);
 const code=(v:unknown):v is string=>typeof v==="string"&&v.length<=64&&/^[a-z][a-z0-9]*(-[a-z0-9]+)*$/.test(v);
 const version=(v:unknown):v is number=>typeof v==="number"&&Number.isSafeInteger(v)&&v>=1&&v<Number.MAX_SAFE_INTEGER;
 const bounded=(v:unknown,max:number):v is string=>typeof v==="string"&&v.trim().length>0&&new TextEncoder().encode(v).length<=max;
 function validResponses(value:unknown):value is CompletionResponse[]{return Array.isArray(value)&&value.length>=1&&value.length<=64&&new Set(value.map(r=>r?.item_id)).size===value.length&&value.every(r=>r&&code(r.item_id)&&bounded(r.response_text,2048)&&(r.evidence_sha256===undefined||r.evidence_sha256===""||typeof r.evidence_sha256==="string"&&/^[a-f0-9]{64}$/.test(r.evidence_sha256)))}
 function valid(v:unknown):v is CompletionResult{if(!v||typeof v!=="object")return false;const value=v as CompletionResult;return id(value.handover_id)&&value.organization_id===organization&&version(value.version)&&value.version>=2&&["presented","accepted","rejected"].includes(value.state)&&code(value.checklist_id)&&version(value.checklist_version)&&typeof value.completed_at==="string"&&Number.isFinite(Date.parse(value.completed_at))&&bounded(value.actor_subject,255)&&validResponses(value.responses)}
 async function digest(value:Pick<CompletionResult,"handover_id"|"organization_id"|"checklist_id"|"checklist_version"|"responses">){
  const responses=[...value.responses].sort((a,b)=>a.item_id<b.item_id?-1:a.item_id>b.item_id?1:0).map(r=>({item_id:r.item_id,response_text:r.response_text,evidence_sha256:r.evidence_sha256??""}));
  const bytes=new TextEncoder().encode(JSON.stringify({handover_id:value.handover_id,organization_id:value.organization_id,checklist_id:value.checklist_id,checklist_version:value.checklist_version,responses}));
  return Array.from(new Uint8Array(await crypto.subtle.digest("SHA-256",bytes)),b=>b.toString(16).padStart(2,"0")).join("");
 }
 useEffect(()=>{try{
  if(!/^[a-f0-9]{64}$/.test(scope))throw new Error("scope");const raw=sessionStorage.getItem(storageKey);
  if(raw!==null){const v=JSON.parse(raw);if(!v||Object.keys(v).sort().join(",")!=="checklistId,checklistVersion,digest,handoverId,version"||!id(v.handoverId)||!version(v.version)||!code(v.checklistId)||!version(v.checklistVersion)||!/^[a-f0-9]{64}$/.test(v.digest))throw new Error("reference");operation.current=v;setLocked(true);setMessage(t("p0547"))}
  else{operation.current=null;setLocked(false)}setReady(true);
 }catch{setReady(true);setLocked(true);setMessage(t("p0518"))}},[scope,storageKey]);
 async function submit(event:FormEvent<HTMLFormElement>){
  event.preventDefault();if(!ready||locked||fence.current||!guidedReady)return;const data=new FormData(event.currentTarget);
  let responses:CompletionResponse[];
  try{const checklist=checklistSchema.parse(JSON.parse(String(data.get("selectedChecklist")??"")));if(checklist.organization_id!==organization||checklist.id!==data.get("checklistId")||checklist.version!==Number(data.get("checklistVersion")))throw new Error("scope");responses=checklistResponses(checklist,data);}
  catch{setMessage("Completá los campos obligatorios de la lista publicada.");return;}
  const value={handover_id:String(data.get("handoverId")??""),organization_id:organization,checklist_id:String(data.get("checklistId")??""),checklist_version:Number(data.get("checklistVersion")),responses};const currentVersion=Number(data.get("handoverVersion"));
  if(!id(value.handover_id)||!code(value.checklist_id)||!version(value.checklist_version)||!version(currentVersion)||!validResponses(responses)){setMessage(t("p0548"));return}
  fence.current=true;setLocked(true);setPosting(true);setConfirmed(false);let reference:CompletionReference;
  try{if(sessionStorage.getItem(storageKey)!==null)throw new Error("unresolved");reference={handoverId:value.handover_id,version:currentVersion,checklistId:value.checklist_id,checklistVersion:value.checklist_version,digest:await digest(value)};sessionStorage.setItem(storageKey,JSON.stringify(reference));operation.current=reference}
  catch{setPosting(false);setMessage(t("p0520"));return}
  setMessage(t("p0549"));
  try{
   const response=await fetch("/api/enterprise/franchise/commands",{method:"POST",headers:{"content-type":"application/json"},signal:AbortSignal.timeout(10000),body:JSON.stringify({action:"complete-delivery-checklist",organizationId:organization,handoverId:value.handover_id,version:currentVersion,checklistId:value.checklist_id,checklistVersion:value.checklist_version,responses})});
   const receipt=await response.json();if(operation.current!==reference)return;
   if(definitiveSelectionRejection(response,receipt)){sessionStorage.removeItem(storageKey);operation.current=null;fence.current=false;setLocked(false);setMessage(experienceText("Las opciones cambiaron antes del envío. No se intentó la operación; actualizá la página para elegir una opción vigente.",privateLocale.language));return;}
   if(!response.ok||receipt?.id!==value.handover_id||receipt.organization_id!==organization||receipt.state!=="presented"||receipt.version!==currentVersion+1||receipt.checklist_id!==value.checklist_id||receipt.checklist_version!==value.checklist_version||typeof receipt.checklist_completed_at!=="string"||!Number.isFinite(Date.parse(receipt.checklist_completed_at)))throw new Error("unconfirmed");
   setMessage(t("p0550"));
  }catch{if(operation.current===reference)setMessage(t("p0551"))}
  finally{setPosting(false)}
 }
 async function consult(){
  const reference=operation.current;if(!ready||!reference||reading.current)return;reading.current=true;setQuerying(true);
  try{
   const response=await fetch("/api/enterprise/franchise/commands?"+new URLSearchParams({kind:"checklist-completion",organizationId:organization,handoverId:reference.handoverId}),{cache:"no-store",signal:AbortSignal.timeout(10000)}),body=await response.json();if(operation.current!==reference)return;
   if(!response.ok||!valid(body?.completion)||body.completion.handover_id!==reference.handoverId||body.completion.version<reference.version+1||body.completion.checklist_id!==reference.checklistId||body.completion.checklist_version!==reference.checklistVersion||await digest(body.completion)!==reference.digest)throw new Error("unconfirmed");
   setResult(body.completion);setConfirmed(true);setMessage(t("p0552"));
  }catch{if(operation.current===reference){setConfirmed(false);setMessage(t("p0553"))}}
  finally{reading.current=false;setQuerying(false)}
 }
 function prepare(){if(!confirmed||posting||reading.current)return;try{sessionStorage.removeItem(storageKey);operation.current=null;fence.current=false;setLocked(false);setConfirmed(false);setResult(null);form.current?.reset();setMessage(t("p0554"))}catch{setMessage(t("p0527"))}}
 return <div role="region" aria-label={t("p0555")}><form ref={form} onSubmit={event=>void submit(event)}><fieldset disabled={!ready||locked}><legend>{t("p0556")}</legend><GuidedChecklistFields organization={organization} handover={initialHandover} onReady={setGuidedReady}/><button disabled={!ready||locked||!initialHandover||!guidedReady}>{experienceText("Confirmar revisión",privateLocale.language)}</button></fieldset></form>
 <p role="status" aria-label={t("p0562")} aria-live="polite">{message}</p>
 {result&&result.organization_id===organization?<details><summary>{experienceText("Ver comprobante",privateLocale.language)}</summary><div><p>{result.handover_id} {t("p0563")} {result.state} {t("p0543")} {result.version}</p><p>{t("p0564")} {result.checklist_id} {t("p0543")} {result.checklist_version} {t("p0565")} {result.actor_subject} {t("p0016")} {result.completed_at}</p><ol>{result.responses.map(r=><li key={r.item_id}>{r.item_id}{t("p0213")} {r.response_text}{r.evidence_sha256?` · evidencia ${r.evidence_sha256}`:""}</li>)}</ol></div></details>:null}
 {operation.current?<button type="button" disabled={!ready||querying} onClick={()=>void consult()}>{experienceText("Consultar estado",privateLocale.language)}</button>:null}
 {confirmed?<button type="button" disabled={posting||querying} onClick={prepare}>{t("p0567")}</button>:null}
 <details><summary>{experienceText("Ayuda",privateLocale.language)}</summary><OperationalGuide guide={OPERATIONAL_GUIDES["checklist-completion-view"]}/></details>
 </div>;
}


type ReturnReference={action:"receive-return"|"decide-return";authorizationId:string;receiptId:string;digest:string};
type ReturnCommand={action:"receive-return";organizationId:string;authorizationId:string;serialNumber:string;conditionCode:string;notes:string}|{action:"decide-return";organizationId:string;receiptId:string;inventoryAction:string;notes:string};
type ReceivedReturn={id:string;authorization_id:string;organization_id:string;order_id:string;stock_unit_id:string;customer_subject:string;received_serial_number:string;condition_code:string;notes:string;evidence_sha256:string;received_by_subject:string;received_at:string};
type DecidedReturn={id:string;receipt_id:string;inventory_action:string;customer_remedy:string;notes:string;decided_by_subject:string;decided_at:string;effect_requests:(ReturnEffect&{idempotency_key:string;requested_at:string})[]};
type RecoveredReturnCase=Omit<ReturnCase,"receipt"|"disposition">&{organization_id:string;authorized_at:string;receipt?:ReceivedReturn;disposition?:DecidedReturn};
function ReturnOperationsPanel({organization,scope,cases}:{organization:string;scope:string;cases:ReturnCase[]|null}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

 const router=useRouter();const [caseList,setCaseList]=useState<ReturnCase[]|null>(cases),[ready,setReady]=useState(false),[locked,setLocked]=useState(true),[message,setMessage]=useState("");
 const [posting,setPosting]=useState(false),[querying,setQuerying]=useState(false),[confirmed,setConfirmed]=useState(false),[result,setResult]=useState<RecoveredReturnCase|null>(null);
 const fence=useRef(false),reading=useRef(false),operation=useRef<ReturnReference|null>(null);const storageKey=`elite-return-operation:${scope}`;
 const text=(v:unknown,max:number):v is string=>typeof v==="string"&&new TextEncoder().encode(v).length>=1&&new TextEncoder().encode(v).length<=max;
 const bounded=(v:unknown,max:number):v is string=>text(v,max)&&v.trim().length>0;
 const date=(v:unknown):v is string=>typeof v==="string"&&Number.isFinite(Date.parse(v));
 const hash=(v:unknown):v is string=>typeof v==="string"&&/^[a-f0-9]{64}$/.test(v);
 const effects:Record<string,string>={inventory:t("p0568"),refund:t("p0569"),exchange:t("p0570"),accounting:t("p0571"),fiscal:t("p0572")};
 function validCase(value:unknown):value is RecoveredReturnCase{
  if(!value||typeof value!=="object")return false;const v=value as RecoveredReturnCase;
  if(!text(v.authorization_id,128)||v.organization_id!==organization||!bounded(v.order_id,128)||!bounded(v.stock_unit_id,128)||!bounded(v.customer_subject,255)||!["return","exchange"].includes(v.authorized_action)||!date(v.authorized_at))return false;
  if(v.receipt===undefined)return v.disposition===undefined;const r=v.receipt;
  if(!r||!text(r.id,128)||r.authorization_id!==v.authorization_id||r.organization_id!==organization||r.order_id!==v.order_id||r.stock_unit_id!==v.stock_unit_id||r.customer_subject!==v.customer_subject||!text(r.received_serial_number,128)||!["sealed","opened","damaged","incomplete"].includes(r.condition_code)||!bounded(r.notes,1000)||!hash(r.evidence_sha256)||!bounded(r.received_by_subject,255)||!date(r.received_at))return false;
  if(v.disposition===undefined)return true;const d=v.disposition,remedy=v.authorized_action==="return"?"refund":"exchange";
  if(!d||!text(d.id,128)||d.receipt_id!==r.id||!["quarantine","restock","repair","scrap"].includes(d.inventory_action)||d.customer_remedy!==remedy||!bounded(d.notes,1000)||!bounded(d.decided_by_subject,255)||!date(d.decided_at)||!Array.isArray(d.effect_requests))return false;
  const owners:Record<string,string>={inventory:"inventory",accounting:"accounting",[remedy]:remedy==="refund"?"payment":"fulfillment"};if(remedy==="refund")owners.fiscal="fiscal";
  return d.effect_requests.length===Object.keys(owners).length&&new Set(d.effect_requests.map(e=>e?.id)).size===d.effect_requests.length&&new Set(d.effect_requests.map(e=>e?.effect_kind)).size===d.effect_requests.length&&d.effect_requests.every(e=>e&&text(e.id,128)&&Object.hasOwn(owners,e.effect_kind)&&owners[e.effect_kind]===e.owner_context&&e.state==="requested"&&text(e.idempotency_key,128)&&e.idempotency_key.length>=16&&date(e.requested_at));
 }
 async function digest(command:ReturnCommand){const bytes=new TextEncoder().encode(JSON.stringify(command));return Array.from(new Uint8Array(await crypto.subtle.digest("SHA-256",bytes)),b=>b.toString(16).padStart(2,"0")).join("")}
 async function matches(value:RecoveredReturnCase,reference:ReturnReference){
  if(value.authorization_id!==reference.authorizationId||!value.receipt)return false;
  if(reference.action==="receive-return"){
   const r=value.receipt;const command:ReturnCommand={action:"receive-return",organizationId:organization,authorizationId:reference.authorizationId,serialNumber:r.received_serial_number,conditionCode:r.condition_code,notes:r.notes};
   return r.evidence_sha256===reference.digest&&await digest(command)===reference.digest;
  }
  if(value.receipt.id!==reference.receiptId||!value.disposition)return false;const d=value.disposition;
  return await digest({action:"decide-return",organizationId:organization,receiptId:reference.receiptId,inventoryAction:d.inventory_action,notes:d.notes})===reference.digest;
 }
 useEffect(()=>setCaseList(cases),[cases]);
 useEffect(()=>{try{
  if(!/^[a-f0-9]{64}$/.test(scope))throw new Error("scope");const raw=sessionStorage.getItem(storageKey);
  if(raw!==null){if(raw.length>2048)throw new Error("reference");const v=JSON.parse(raw);if(!v||Object.keys(v).sort().join(",")!=="action,authorizationId,digest,receiptId"||!["receive-return","decide-return"].includes(v.action)||!text(v.authorizationId,128)||!hash(v.digest)||(v.action==="receive-return"?v.receiptId!=="":!text(v.receiptId,128)))throw new Error("reference");operation.current=v;setLocked(true);setMessage(t("p0573"))}
  else{operation.current=null;setLocked(false)}setReady(true);
 }catch{setReady(true);setLocked(true);setMessage(t("p0518"))}},[scope,storageKey]);
 async function submit(event:FormEvent<HTMLFormElement>,item:ReturnCase,action:ReturnReference["action"]){
  event.preventDefault();if(!ready||locked||fence.current)return;const data=new FormData(event.currentTarget);
  const command:ReturnCommand=action==="receive-return"?{action:"receive-return",organizationId:organization,authorizationId:item.authorization_id,serialNumber:String(data.get("serialNumber")??""),conditionCode:String(data.get("conditionCode")??""),notes:String(data.get("notes")??"").trim()}:{action:"decide-return",organizationId:organization,receiptId:item.receipt?.id??"",inventoryAction:String(data.get("inventoryAction")??""),notes:String(data.get("notes")??"").trim()};
  if(!text(item.authorization_id,128)||!bounded(command.notes,1000)||(command.action==="receive-return"?(!text(command.serialNumber,128)||!["sealed","opened","damaged","incomplete"].includes(command.conditionCode)):(!text(command.receiptId,128)||!["quarantine","restock","repair","scrap"].includes(command.inventoryAction)))){setMessage(t("p0574"));return}
  fence.current=true;setLocked(true);setPosting(true);setConfirmed(false);setResult(null);let reference:ReturnReference;
  try{if(sessionStorage.getItem(storageKey)!==null)throw new Error("unresolved");reference={action,authorizationId:item.authorization_id,receiptId:command.action==="decide-return"?command.receiptId:"",digest:await digest(command)};sessionStorage.setItem(storageKey,JSON.stringify(reference));operation.current=reference}
  catch{setPosting(false);setMessage(t("p0520"));return}
  setMessage(action==="receive-return"?t("p0575"):t("p0576"));
  try{
   const response=await fetch("/api/enterprise/franchise/commands",{method:"POST",headers:{"content-type":"application/json"},signal:AbortSignal.timeout(10000),body:JSON.stringify(command)}),body=await response.json();if(operation.current!==reference)return;
   if(!response.ok||!text(body?.id,128)||(action==="receive-return"?(body.authorization_id!==reference.authorizationId||body.organization_id!==organization||body.evidence_sha256!==reference.digest):body.receipt_id!==reference.receiptId))throw new Error("unconfirmed");
   setMessage(t("p0577"));
  }catch{if(operation.current===reference)setMessage(t("p0578"))}
  finally{setPosting(false)}
 }
 async function consult(){
  const reference=operation.current;if(!ready||!reference||reading.current)return;reading.current=true;setQuerying(true);
  try{
   const response=await fetch("/api/enterprise/franchise/commands?"+new URLSearchParams({kind:"return",organizationId:organization,authorizationId:reference.authorizationId}),{cache:"no-store",signal:AbortSignal.timeout(10000)}),body=await response.json();if(operation.current!==reference)return;
   if(!response.ok||!validCase(body?.returnCase)||!await matches(body.returnCase,reference))throw new Error("unconfirmed");
   setResult(body.returnCase);setConfirmed(true);setMessage(t("p0579"));
  }catch{if(operation.current===reference){setConfirmed(false);setMessage(t("p0580"))}}
  finally{reading.current=false;setQuerying(false)}
 }
 function prepare(){
  if(!confirmed||!result||posting||reading.current)return;
  try{sessionStorage.removeItem(storageKey);const recovered=result;setCaseList(previous=>previous?.some(item=>item.authorization_id===recovered.authorization_id)?previous.map(item=>item.authorization_id===recovered.authorization_id?recovered:item):[recovered,...(previous??[])]);operation.current=null;fence.current=false;setLocked(false);setConfirmed(false);setResult(null);setMessage(t("p0581"))}
  catch{setMessage(t("p0527"))}
 }
 return <section className="card" aria-label={t("p0582")}><h2>{t("p0583")}</h2><p>{t("p0584")}</p>
 <p role="status" aria-label={t("p0585")} aria-live="polite">{message}</p>
 <button type="button" disabled={!ready||!operation.current||querying} onClick={()=>void consult()}>{t("p0586")}</button>
 <button type="button" disabled={!confirmed||posting||querying} onClick={prepare}>{t("p0587")}</button>
 <button type="button" disabled={posting||querying} onClick={()=>router.refresh()}>{t("p0588")}</button>
 {result?<div role="region" aria-label={t("p0589")}><h3>{t("p0590")} {result.authorization_id}</h3>{result.receipt?<><p>{t("p0591")} {result.receipt.id} {t("p0592")} {result.receipt.received_by_subject} {t("p0016")} {result.receipt.received_at}</p><p>{t("p0593")} {result.receipt.received_serial_number} {t("p0016")} {result.receipt.condition_code}</p><p>{result.receipt.notes}</p></>:null}{result.disposition?<><p>{t("p0594")} {result.disposition.id} {t("p0595")} {result.disposition.decided_by_subject} {t("p0016")} {result.disposition.decided_at}</p><p>{result.disposition.inventory_action} {t("p0016")} {result.disposition.customer_remedy} {t("p0016")} {result.disposition.notes}</p><ul>{result.disposition.effect_requests.map(e=><li key={e.id}>{effects[e.effect_kind]}{t("p0596")} {e.id}</li>)}</ul></>:null}</div>:null}
 {cases===null?<p role="alert">{t("p0597")}</p>:null}
 {caseList?.length?<div className="grid">{caseList.map(item=><article className="card" key={item.authorization_id} aria-label={`Caso ${item.authorization_id}`}><h3>{item.authorized_action==="return"?t("p0598"):t("p0570")} {t("p0016")} {item.authorization_id}</h3><p>{t("p0251")} {item.order_id} {t("p0599")} {item.stock_unit_id}</p>
 {!item.receipt?<form onSubmit={event=>void submit(event,item,"receive-return")}><fieldset disabled={!ready||locked}><label>{t("p0600")}<input name="serialNumber" required maxLength={128}/></label><label>{t("p0601")}<select name="conditionCode" required><option value="sealed">{t("p0602")}</option><option value="opened">{t("p0603")}</option><option value="damaged">{t("p0604")}</option><option value="incomplete">{t("p0605")}</option></select></label><label>{t("p0606")}<textarea name="notes" required maxLength={1000}/></label><button disabled={!ready||locked}>{t("p0261")}</button></fieldset></form>:!item.disposition?<form onSubmit={event=>void submit(event,item,"decide-return")}><fieldset disabled={!ready||locked}><p>{t("p0607")} {item.receipt.received_serial_number} {t("p0016")} {item.receipt.condition_code}</p><label>{t("p0608")}<select name="inventoryAction" required><option value="quarantine">{t("p0609")}</option><option value="restock">{t("p0610")}</option><option value="repair">{t("p0611")}</option><option value="scrap">{t("p0612")}</option></select></label><label>{t("p0613")}<textarea name="notes" required maxLength={1000}/></label><button disabled={!ready||locked}>{t("p0614")}</button></fieldset></form>:<><p>{t("p0591")} {item.receipt.id} {t("p0615")} {item.disposition.id}</p><p>{item.disposition.inventory_action} {t("p0016")} {item.disposition.customer_remedy}</p><ul>{item.disposition.effect_requests.map(e=><li key={e.id}>{effects[e.effect_kind]}{t("p0596")} {e.id}</li>)}</ul></>}
 </article>)}</div>:cases!==null?<p>{t("p0616")}</p>:null}
 <OperationalGuide guide={OPERATIONAL_GUIDES["return-operations-view"]}/>
 </section>;
}
