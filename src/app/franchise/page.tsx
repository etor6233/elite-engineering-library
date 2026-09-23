import {portalCursor} from "@/platform/backend/portal-paging";
import {experienceText} from "@/platform/experience/messages";
import {organizationLabel} from "@/platform/experience/options";
import {TaskWorkspace} from "@/components/task-workspace";
import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import { HandoverOperationsPanel } from "@/components/handover-operations-panel";
import { OperationalGuide } from "@/components/operational-guide";
import { OPERATIONAL_GUIDES } from "@/platform/help/content";
import { redirect } from "next/navigation";
import { createHash } from "node:crypto";
import type { Route } from "next";
import { allowed, readSession } from "@/platform/auth/session";
import { protectedGet, type LeadSummary, type Page } from "@/platform/backend/protected-client";
import { AppointmentAgendaPanel, FranchiseCommandPanel, OrderOperationsPanel, type OrderOperations, type AppointmentAgenda } from "@/components/franchise-command-panel";
import { loadBusinessConfig } from "@/platform/config/load";
import styles from "@/components/admin-ops/admin-ops.module.css";

type Availability = { id: string; organization_id: string; resource_id?: string; entry_type: "working" | "unavailable"; reason_code?: string; starts_at: string; ends_at: string; state: string; version: number };
type DeliveryException = { id: string; handover_id: string; customer_subject: string; reason_code: string; details: string; state: string; version: number; resolution_action?: string };
type ReturnEffect = { id: string; effect_kind: "inventory" | "refund" | "exchange" | "accounting" | "fiscal"; owner_context: string; state: "requested" };
type ReturnCase = { authorization_id: string; order_id: string; stock_unit_id: string; customer_subject: string; authorized_action: "return" | "exchange"; receipt?: { id: string; received_serial_number: string; condition_code: string; received_at: string }; disposition?: { id: string; inventory_action: string; customer_remedy: string; effect_requests: ReturnEffect[] } };

export default async function FranchisePage({ searchParams }: { searchParams: Promise<{ day?: string; task?: string; leads_after?: string | string[] }> }) {
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

  const session = await readSession();
  if (!session) redirect("/api/auth/login?return_to=/franchise" as Route);
  const canOperate=["inventory:allocate","payment:create","handover:manage","admin:read"].some(permission=>allowed(session,permission));
  const canPanel=["lead:read","handover:manage","resource:manage","availability:read","availability:manage","appointment:manage"].some(permission=>allowed(session,permission));
  const canReviewWhatsApp=allowed(session,"whatsapp:approve");
  const canStoredValue=allowed(session,"stored_value:read");
  if (!canOperate && !canPanel && !canReviewWhatsApp && !canStoredValue) return <><h1 className="pageTitle">{t("p0006")}</h1><div className="notice">{t("p0074")}</div></>;
  const organization = session.organizations[0]!;
  const branchName=await organizationLabel(session,organization);
  const notificationStatusEnabled = (await loadBusinessConfig()).features.whatsapp_status_history === true;
  const from = new Date().toISOString();
  const to = new Date(Date.now() + 90 * 24 * 60 * 60_000).toISOString();
  const requestedQuery=await searchParams;
  const leadsAfter=allowed(session,"lead:read")?portalCursor(requestedQuery.leads_after):undefined;
  const requestedDay=requestedQuery.day ?? new Date().toISOString().slice(0,10);
  const parsedDay=new Date(requestedDay+"T00:00:00Z");
  if (!/^\d{4}-\d{2}-\d{2}$/.test(requestedDay) || !Number.isFinite(parsedDay.getTime()) || parsedDay.toISOString().slice(0,10)!==requestedDay) return <><h1>{t("p0075")}</h1><a href="/franchise">{t("p0076")}</a></>;
  const agendaFrom=parsedDay.toISOString(), agendaTo=new Date(parsedDay.getTime()+24*60*60_000).toISOString();
  // Each owner may fail independently; null means unavailable, never an empty result.
  const [operationSnapshot,leads,availability,exceptions,returnCases,agenda]=await Promise.all([
    canOperate ? protectedGet<OrderOperations>(session,"/v1/commerce/orders",{organization_id:organization}).catch(()=>null) : Promise.resolve(null),
    allowed(session,"lead:read") ? protectedGet<Page<LeadSummary>>(session,"/v1/franchise/leads",{organization_id:organization,limit:"50",after:leadsAfter}).catch(()=>null) : Promise.resolve<Page<LeadSummary>>({items:[]}),
    allowed(session,"availability:read") ? protectedGet<{items:Availability[]}>(session,"/v1/franchise/availability",{organization_id:organization,from,to}).catch(()=>null) : Promise.resolve({items:[]}),
    allowed(session,"handover:manage") ? protectedGet<{items:DeliveryException[]}>(session,"/v1/franchise/delivery-exceptions",{organization_id:organization}).catch(()=>null) : Promise.resolve({items:[]}),
    allowed(session,"handover:manage") ? protectedGet<{items:ReturnCase[]}>(session,"/v1/franchise/returns",{organization_id:organization}).catch(()=>null) : Promise.resolve({items:[]}),
    allowed(session,"appointment:manage") ? protectedGet<AppointmentAgenda>(session,"/v1/franchise/agenda",{organization_id:organization,from:agendaFrom,to:agendaTo}).catch(()=>null) : Promise.resolve(null)
  ]);
  const failed=[
    leads===null ? "leads" : null,
    availability===null ? "disponibilidad" : null,
    exceptions===null ? t("p0077") : null,
    returnCases===null ? "devoluciones" : null,
    allowed(session,"appointment:manage") && agenda===null ? "agenda" : null
  ].filter((item):item is string=>item!==null);
  const requestedTask = requestedQuery.task;
  const contactPageHref=(after?:string)=>"/franchise?"+new URLSearchParams({day:requestedDay,task:"daily",...(after?{leads_after:portalCursor(after)!}:{})}).toString();
  const taskLabels: Record<string,string> = {agenda:"Agenda",deliveries:"Entregas",orders:"Pedidos",daily:"Gestión diaria"};
  return <div className={styles.page} data-admin-root="franchise">
    <div className={styles.branch}><svg aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5"><path d="M4 10h16v11H4zM3 10l2-7h14l2 7M9 21v-7h6v7"/></svg><span>{branchName??experienceText("Sucursal: ",privateLocale.language)+organization}</span></div>
    {failed.length>0 ? <div className={styles.ownerNotice} role="alert">{privateLocale.language==="en"?"Some information is unavailable. Open Retry to recover it.":"Hay información no disponible. Abrí Reintentar para recuperarla."}</div> : null}
    <TaskWorkspace sidebar title="Área de trabajo" initialLabel={requestedTask&&taskLabels[requestedTask] ? taskLabels[requestedTask] : canOperate?"Pedidos":"Gestión diaria"}>
    {canStoredValue ? <section data-task-label="Saldos"><p><a className="button" href="/franchise/stored-value">{t("p0081")}</a></p></section> : null}
    {canReviewWhatsApp ? <section data-task-label="Mensajes"><p><a className="button" href="/franchise/whatsapp">{t("p0082")}</a></p></section> : null}
    {canOperate ? <section data-task-label="Pedidos"><OrderOperationsPanel snapshot={operationSnapshot} organization={organization} canAllocate={allowed(session,"inventory:allocate")} canRequestPayment={allowed(session,"payment:create")} paymentScope={createHash("sha256").update(JSON.stringify([session.tenantId,session.subject,organization])).digest("hex")}/></section> : null}
    {allowed(session,"handover:manage") ? <section data-task-label="Entregas"><HandoverOperationsPanel orders={operationSnapshot && !operationSnapshot.truncated ? operationSnapshot.orders : null} organization={organization} scope={createHash("sha256").update(JSON.stringify([session.tenantId,session.subject,organization])).digest("hex")}/></section> : null}
    {failed.length>0 ? <section data-task-label="Reintentar" role="alert">{failed.map(name=><p key={name}>{t("p0083")} {name}{t("p0084")}</p>)}<a href={("/franchise?"+new URLSearchParams({day:requestedDay,task:requestedTask&&["agenda","deliveries","orders","daily"].includes(requestedTask)?requestedTask:"daily",...(leadsAfter?{leads_after:leadsAfter}:{})})) as Route}>{t("p0085")}</a></section> : null}
    {allowed(session,"appointment:manage") ? <section data-task-label="Agenda"><form action="/franchise" method="get"><input type="hidden" name="task" value="agenda"/><label>{t("p0086")}<input type="date" name="day" defaultValue={requestedDay} required /></label><button className="button">{t("p0087")}</button></form>{agenda ? <AppointmentAgendaPanel agenda={agenda} organization={organization} notificationStatusEnabled={notificationStatusEnabled}/> : null}</section> : null}
    {canPanel ? <section data-task-label="Gestión diaria"><FranchiseCommandPanel commandScope={createHash("sha256").update(JSON.stringify([session.tenantId,session.subject,organization])).digest("hex")} permissions={session.permissions} leads={leads?.items??[]} availability={availability?.items??[]} exceptions={exceptions?.items??null} returnCases={returnCases?.items??null} defaultAssignee={session.subject} organization={organization}/>{leads&&(leads.next_cursor||leadsAfter)?<nav aria-label={privateLocale.language==="en"?"Contact pages":"Páginas de contactos"}>{leadsAfter?<a className="button" href={contactPageHref() as Route}>{privateLocale.language==="en"?"First contacts":"Primeros contactos"}</a>:null}{leads.next_cursor?<a className="button" href={contactPageHref(leads.next_cursor) as Route}>{privateLocale.language==="en"?"More contacts":"Más contactos"}</a>:null}</nav>:null}</section> : null}
    </TaskWorkspace>
  </div>;
}
