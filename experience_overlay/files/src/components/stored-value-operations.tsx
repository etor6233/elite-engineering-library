"use client";
import { GuidedRecordSelect } from "@/components/guided-record-select";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

// AUTHORED operator interaction and recovery glue. Browser submits references
// and selected operations; it never supplies an amount or approves itself.
import { useEffect,useState } from "react";
import { storedCommand,storedOrder,storedApproval,storedFunding,formatStoredMoney,formatStoredAdjustment,type StoredCommand,type StoredOrder,type StoredApproval,type StoredFunding } from "@/platform/stored-value/contracts";
type Program={id:string;kind:"gift_card"|"loyalty"};
export function StoredValueOperations({organization,profileHash,programs,permissions,subject,scope}:{organization:string;profileHash:string;programs:Program[];permissions:string[];subject:string;scope:string}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

 const [orderID,setOrderID]=useState("");const [order,setOrder]=useState<StoredOrder|null>(null);const [view,setView]=useState<StoredApproval|null>(null);const [funding,setFunding]=useState<StoredFunding|null>(null);
 const [program,setProgram]=useState(programs[0]?.id??"");const [operation,setOperation]=useState<"issue"|"accrue"|"redeem"|"reverse">("redeem");const [account,setAccount]=useState("");const [original,setOriginal]=useState("");const [reason,setReason]=useState("");
 const [busy,setBusy]=useState(false);const [notice,setNotice]=useState("");const [pending,setPending]=useState<StoredCommand|null>(null);const [storageReady,setStorageReady]=useState(false);
 const storageKey=`stored-value-request:${scope}`;
 useEffect(()=>{try{const raw=localStorage.getItem(storageKey);if(raw!==null){const v=storedCommand.parse(JSON.parse(raw));if(v.organization_id!==organization)throw new Error("scope");setPending(v)}setStorageReady(true)}catch{setNotice(t("p0757"))}},[storageKey,organization]);
 const can=(p:string)=>permissions.includes(p);
 async function read(query:Record<string,string>){const r=await fetch("/api/enterprise/franchise/stored-value?"+new URLSearchParams({organization_id:organization,...query}),{cache:"no-store"});if(!r.ok)throw new Error("unconfirmed");return r.json() as Promise<unknown>}
 async function loadOrder(after?:string){if(busy)return;setBusy(true);setNotice("");setOrder(null);try{setOrder(storedOrder.parse(await read({kind:"order",order_id:orderID,...(after?{after}:{})})))}catch{setNotice(t("p0758"))}finally{setBusy(false)}}
 async function loadApproval(id:string){if(busy)return;setBusy(true);setNotice("");setView(null);try{setView(storedApproval.parse(await read({kind:"approval",approval_id:id})))}catch{setNotice(t("p0759"))}finally{setBusy(false)}}
 function keep(c:StoredCommand){const v=storedCommand.parse(c);const old=localStorage.getItem(storageKey);const raw=JSON.stringify(v);if(old!==null&&old!==raw)throw new Error("pending request");localStorage.setItem(storageKey,raw);if(localStorage.getItem(storageKey)!==raw)throw new Error("not retained");setPending(v)}
 function clear(){localStorage.removeItem(storageKey);setPending(null)}
 async function send(c:StoredCommand){
  if(busy||!storageReady)return;setBusy(true);setNotice("");
  try{keep(c);const r=await fetch("/api/enterprise/franchise/stored-value",{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify(c)});if(!r.ok)throw new Error("unconfirmed");const raw:unknown=await r.json();if(c.action==="fund")setFunding(storedFunding.parse(raw));else setView(storedApproval.parse(raw));clear();setOrder(null);setNotice(c.action==="propose"?t("p0760"):c.action==="fund"?t("p0761"):t("p0762"))}
  catch{setNotice(t("p0763"))}
  finally{setBusy(false)}
 }
 async function recover(){
  if(!pending||busy)return;setBusy(true);setNotice("");
  try{
   if(pending.action==="fund")setFunding(storedFunding.parse(await read({kind:"funding-result",request_key:pending.request_key})));
   else {const v=storedApproval.parse(await read(pending.action==="propose"?{kind:"proposal-result",operation_id:pending.operation_id}:{kind:"approval",approval_id:pending.approval_id}));setView(v);if(pending.action==="decision"&&v.state==="pending"){setNotice(t("p0764"));return}}
   clear();setOrder(null);setNotice(t("p0765"));
  }catch{setNotice(t("p0766"))}finally{setBusy(false)}
 }
 const blocked=busy||!storageReady||pending!==null;
 const submit=()=>{if(!order)return;void send({action:"propose",organization_id:organization,operation_id:crypto.randomUUID(),operation,program_id:program,order_id:order.order_id,expected_order_version:order.order_version,account_id:operation==="redeem"?account:"",original_operation_id:operation==="reverse"?original:"",profile_sha256:profileHash})};
 return <>
  <p role="status" aria-live="polite">{notice}</p>
  {pending&&<section className="card" aria-label={t("p0767")}><h2>{t("p0767")}</h2><p>{t("p0768")}</p><div style={{display:"flex",gap:12,flexWrap:"wrap"}}><button className="button" disabled={busy} onClick={()=>void recover()}>{t("p0168")}</button><button className="button" disabled={busy} onClick={()=>void send(pending)}>{t("p0769")}</button></div></section>}
  <section className="card"><h2>{t("p0770")}</h2><form onSubmit={e=>{e.preventDefault();void loadOrder()}}><GuidedRecordSelect organization={organization} group="storedOrders" name="order" label="Pedido" disabled={!storageReady||busy||Boolean(pending)} value={orderID} onChange={nextValue=>{setOrderID(nextValue);setOrder(null);setView(null);setFunding(null)}}/><button className="button" disabled={busy}>{t("p0770")}</button></form>
   {order&&<><dl><dt>{t("p0772")}</dt><dd>{formatStoredMoney(order.gross_minor_units)}</dd><dt>{t("p0773")}</dt><dd>{formatStoredMoney(order.gift_minor_units)}</dd><dt>{t("p0774")}</dt><dd>{formatStoredMoney(order.discount_minor_units)}</dd><dt>{t("p0775")}</dt><dd>{formatStoredMoney(order.provider_due_minor_units)}</dd></dl>
    <h3>{t("p0776")}</h3>{order.approvals.ids.length===0?<p>{t("p0777")}</p>:<ul>{order.approvals.ids.map(id=><li key={id}><button className="button" disabled={busy} onClick={()=>void loadApproval(id)}>{t("p0778")} {id}</button></li>)}</ul>}
    {order.approvals.next_cursor&&<button className="button" disabled={busy} onClick={()=>void loadOrder(order.approvals.next_cursor)}>{t("p0779")}</button>}
    {can("stored_value:fund")&&order.provider_due_minor_units==="0"&&<button className="button" disabled={blocked} onClick={()=>void send({action:"fund",organization_id:organization,order_id:order.order_id,expected_order_version:order.order_version,request_key:crypto.randomUUID()})}>{t("p0780")}</button>}
   </>}
  </section>
  {can("stored_value:request")&&<section className="card"><h2>{t("p0781")}</h2><p>{t("p0782")}</p><form onSubmit={e=>{e.preventDefault();submit()}}>
   <label>{t("p0783")}<select value={program} onChange={e=>setProgram(e.target.value)}>{programs.map(p=><option key={p.id} value={p.id}>{p.kind==="gift_card"?t("p0784"):t("p0785")} {t("p0016")} {p.id}</option>)}</select></label>
   <label>{t("p0010")}<select value={operation} onChange={e=>setOperation(e.target.value as typeof operation)}><option value="redeem">{t("p0786")}</option><option value="issue">{t("p0787")}</option><option value="accrue">{t("p0788")}</option><option value="reverse">{t("p0789")}</option></select></label>
   {operation==="redeem"&&<GuidedRecordSelect organization={organization} group="storedAccounts" name="reference:account" label="Cuenta del programa" disabled={!storageReady||busy||Boolean(pending)} value={account} onChange={nextValue=>setAccount(nextValue)}/>}
   {operation==="reverse"&&<GuidedRecordSelect organization={organization} group="storedOperations" name="reference:original" label="Operación original" disabled={!storageReady||busy||Boolean(pending)} value={original} onChange={nextValue=>setOriginal(nextValue)}/>}
   <button className="button" disabled={blocked||!order}>{t("p0792")}</button>
  </form></section>}
  {view&&<section className="card" style={{overflowWrap:"anywhere"}}><h2>{t("p0793")}</h2><p>{t("p0794")} {view.order_id}{t("p0795")} {view.requester}{t("p0796")} {view.state==="pending"?t("p0797"):view.state==="approved"?t("p0798"):t("p0337")}{t("p0008")}</p><p>{t("p0799")} {view.expires_at}</p><p>{t("p0800")}</p><dl><dt>{t("p0783")}</dt><dd>{view.review.program_id}</dd><dt>{t("p0010")}</dt><dd>{{issue:t("p0801"),accrue:t("p0802"),redeem:t("p0803"),reverse:t("p0804")}[view.operation]}</dd><dt>{t("p0805")}</dt><dd>{view.review.operation_id}</dd><dt>{t("p0806")}</dt><dd>{formatStoredMoney(view.review.gross_minor_units)}</dd>{view.review.original_operation_id&&<><dt>{t("p0807")}</dt><dd>{view.review.original_operation_id}</dd></>}</dl>
   {view.review.entries.map(entry=><article key={entry.account_id}><h3>{t("p0808")} {entry.account_id}</h3><p>{t("p0809")} {entry.points_delta}</p><p>{t("p0810")} {formatStoredAdjustment(entry.applied_minor_units)}</p></article>)}
   <details><summary>{t("p0811")}</summary><pre style={{whiteSpace:"pre-wrap",overflowWrap:"anywhere",maxHeight:400,overflow:"auto"}} aria-label={t("p0812")}>{view.payload_json}</pre></details>
   {can("stored_value:approve")&&view.state==="pending"&&view.requester!==subject&&<><label>{t("p0220")}<textarea value={reason} onChange={e=>setReason(e.target.value)} maxLength={2048} required/></label><div style={{display:"flex",gap:12,flexWrap:"wrap"}}>{[true,false].map(approved=><button className="button" key={String(approved)} disabled={blocked||!reason.trim()} onClick={()=>void send({action:"decision",organization_id:organization,approval_id:view.approval_id,payload_sha256:view.payload_sha256,approved,reason})}>{approved?t("p0813"):t("p0814")}</button>)}</div></>}
   {view.requester===subject&&view.state==="pending"&&<p>{t("p0815")}</p>}
   <button className="button" disabled={busy} onClick={()=>void loadApproval(view.approval_id)}>{t("p0816")}</button>
   {view.receipt_json!=="null"&&<details><summary>{t("p0817")}</summary><pre style={{whiteSpace:"pre-wrap",overflowWrap:"anywhere"}}>{view.receipt_json}</pre></details>}
  </section>}
  {funding&&<section className="card" style={{overflowWrap:"anywhere"}}><h2>{t("p0818")}</h2><p>{t("p0251")} {funding.order_id}{t("p0819")} {funding.funding_receipt_id}</p><a className="button" href="/franchise">{t("p0820")}</a><details><summary>{t("p0821")}</summary><pre style={{whiteSpace:"pre-wrap",overflowWrap:"anywhere"}}>{funding.receipt_json}</pre></details></section>}
 </>;
}
