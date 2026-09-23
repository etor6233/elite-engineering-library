"use client";
import { AuthorizedUnitCodePicker } from "@/components/authorized-unit-code-picker";
import { GuidedRecordSelect } from "@/components/guided-record-select";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

import {OperationalGuide} from "@/components/operational-guide";
import {RELEASE_GUIDES} from "@/platform/help/content";
import {useEffect,useRef,useState,type FormEvent} from "react";
import {supplyCommand,supplyReceipt,supplyPlan,pendingSupply,supplyDigest,supplyReference,supplyMatches,supplyPolicy,type SupplyCommand,type SupplyPlan,type PendingSupply} from "@/platform/supply/contract";
const endpoint="/api/enterprise/supply";

export function SupplyWorkspace({organization,subject,permissions,scope,surface,initial}:{organization:string;subject:string;permissions:string[];scope:string;surface:"factory"|"franchise";initial:string}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();
const labels:Record<string,string>={draft:t("p0648"),submitted:t("p0822"),accepted:t("p0823"),"in-production":t("p0824"),shipped:t("p0825"),received:t("p0826"),cancelled:t("p0056"),planned:t("p0827"),assembly:t("p0828"),quality:t("p0829"),released:t("p0830"),rejected:t("p0337"),"in-transit":t("p0831"),quarantine:t("p0832"),available:t("p0833"),retired:t("p0834")};

 const can=(p:string)=>permissions.includes("*")||permissions.includes(p);
 const held=useRef(false),saved=useRef<PendingSupply|null>(null);
 const [ready,setReady]=useState(false),[busy,setBusy]=useState(false),[pending,setPending]=useState<PendingSupply|null>(null),[notice,setNotice]=useState("");
 const [id,setID]=useState(initial),[plan,setPlan]=useState<SupplyPlan|null>(null),[evidence,setEvidence]=useState(""),[reason,setReason]=useState("");
 const [supplier,setSupplier]=useState(""),[factory,setFactory]=useState(""),[currency,setCurrency]=useState(""),[total,setTotal]=useState(""),[demand,setDemand]=useState("");
 const [lines,setLines]=useState([{id:"line-1",variant_id:"",quantity:1}]),[line,setLine]=useState(""),[serial,setSerial]=useState(""),[vin,setVIN]=useState(""),[battery,setBattery]=useState("");
 const [completeUnitPage,setCompleteUnitPage]=useState(false);
 const [selected,setSelected]=useState<string[]>([]),[shipment,setShipment]=useState("");
 const key=`elite-supply:${scope}`;
 function persist(value:PendingSupply|null){if(value){const raw=JSON.stringify(pendingSupply.parse(value));sessionStorage.setItem(key,raw);if(sessionStorage.getItem(key)!==raw)throw new Error("storage")}else sessionStorage.removeItem(key);saved.current=value;setPending(value)}
 useEffect(()=>{try{const raw=sessionStorage.getItem(key);if(raw){const v=pendingSupply.parse(JSON.parse(raw));saved.current=v;setPending(v);setID(v.purchase_order_id)}setReady(true)}catch{setNotice(t("p0835"))}},[key]);
 async function get(order:string,extra:Record<string,string>={}){const r=await fetch(endpoint+"?"+new URLSearchParams({organization_id:organization,purchase_order_id:order,surface,...extra}),{cache:"no-store",signal:AbortSignal.timeout(10000)});if(!r.ok)throw new Error("unavailable");return r.json() as Promise<unknown>}
 async function readPlan(order:string,after=""){const v=supplyPlan.parse(await get(order,after?{after_unit:after}:{}));if(v.purchase_order_id!==order||(surface==="factory"?v.factory_organization_id:v.destination_organization_id)!==organization)throw new Error("scope");setPlan(v);setCompleteUnitPage(!after&&!v.next_unit_id);setID(order);setSelected([]);setLine(v.lines[0]?.id??"")}
 async function consult(after=""){if(held.current)return;held.current=true;setBusy(true);setPlan(null);try{await readPlan(id,after);setNotice(t("p0836"))}catch{setNotice(t("p0837"))}finally{held.current=false;setBusy(false)}}
 async function send(input:unknown){
  if(held.current||!ready||saved.current)return;held.current=true;setBusy(true);let started=false;
  try{
   const c=supplyCommand.parse(input),ref=await supplyReference(c);persist(ref);started=true;
   const response=await fetch(endpoint,{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify(c),signal:AbortSignal.timeout(10000)});
   if(!response.ok){if([400,401,403,413,415].includes(response.status)){persist(null);started=false}throw new Error("unconfirmed")}
   const value=supplyReceipt.parse(await response.json());if(!supplyMatches(value,ref,subject))throw new Error("identity");
   persist(null);started=false;setID(value.purchase_order_id);setPlan(null);setNotice(t("p0838"));
   try{await readPlan(value.purchase_order_id)}catch{setNotice(t("p0839"))}
  }catch{setNotice(started?t("p0660"):t("p0840"))}
  finally{held.current=false;setBusy(false)}
 }
 async function recover(){const p=saved.current;if(!p||held.current)return;held.current=true;setBusy(true);try{const value=supplyReceipt.parse(await get(p.purchase_order_id,{command_id:p.command_id}));if(!supplyMatches(value,p,subject))throw new Error("identity");persist(null);setPlan(null);setID(value.purchase_order_id);setNotice(t("p0841"));try{await readPlan(value.purchase_order_id)}catch{setNotice(t("p0842"))}}catch{setNotice(t("p0843"))}finally{held.current=false;setBusy(false)}}
 async function selectEvidence(file:File|null){setEvidence("");if(!file||file.size<1||file.size>1048576){setNotice(t("p0844"));return}setEvidence(await supplyDigest(new Uint8Array(await file.arrayBuffer())))}
 const blocked=busy||!ready||pending!==null;
 const command=(kind:SupplyCommand["kind"],extra:Record<string,unknown>={})=>{if(plan)void send({kind,organization_id:organization,purchase_order_id:plan.purchase_order_id,command_id:crypto.randomUUID(),evidence_sha256:evidence,expected_version:plan.version,...extra})};
 function create(e:FormEvent){e.preventDefault();void send({kind:"create",organization_id:organization,destination_organization_id:organization,purchase_order_id:crypto.randomUUID(),command_id:crypto.randomUUID(),supplier_id:supplier,factory_organization_id:factory,currency,total_minor_units:total,demand_reference:demand,policy_code:supplyPolicy,evidence_sha256:evidence,lines})}
 return <>
 <p role="status" aria-live="polite">{notice}</p>
 <OperationalGuide className="card" guide={RELEASE_GUIDES["supply-role-view"]}/>
 {pending&&<section className="card" aria-label={t("p0166")}><h2>{t("p0166")}</h2><p>{t("p0845")} <code>{pending.purchase_order_id}</code></p><p>{t("p0846")} <code>{pending.command_id}</code></p><button className="button" disabled={busy} onClick={()=>void recover()}>{t("p0168")}</button></section>}
 <section className="card"><h2>{t("p0847")}</h2><label>{t("p0848")}<input type="file" disabled={blocked} onChange={e=>void selectEvidence(e.target.files?.[0]??null)}/></label><p>{t("p0849")}</p>{evidence&&<p>{t("p0850")}</p>}<label>{t("p0851")}<textarea value={reason} onChange={e=>setReason(e.target.value)} maxLength={2000} disabled={blocked}/></label></section>
 {surface==="franchise"&&can("supply:plan")&&<section className="card"><h2>{t("p0852")}</h2><form onSubmit={create}><fieldset disabled={blocked}>
 <GuidedRecordSelect organization={organization} group="supplySuppliers" name="reference:supplier" label="Proveedor" disabled={blocked} value={supplier} onChange={nextValue=>setSupplier(nextValue)}/>
 <GuidedRecordSelect organization={organization} group="supplyFactories" name="reference:factory" label="Fábrica" disabled={blocked} value={factory} onChange={nextValue=>setFactory(nextValue)}/>
 <label>{t("p0855")}<input value={currency} onChange={e=>setCurrency(e.target.value)} pattern="[A-Z]{3}" maxLength={3} required/></label>
 <label>{t("p0856")}<input value={total} onChange={e=>setTotal(e.target.value)} inputMode="numeric" pattern="[0-9]+" maxLength={19} required/></label>
 <label>{t("p0857")}<input value={demand} onChange={e=>setDemand(e.target.value)} maxLength={1024} required/></label>
 {lines.map((v,i)=><fieldset key={v.id}><legend>{t("p0858")} {i+1}</legend><GuidedRecordSelect organization={organization} group="supplyVariants" name="reference:v.variant_id" label="Producto de la línea" disabled={blocked} value={v.variant_id} onChange={nextValue=>setLines(old=>old.map((x,j)=>i===j?{...x,variant_id:nextValue}:x))}/><label>{t("p0859")}<input type="number" min={1} max={1000} step={1} value={v.quantity} onChange={e=>setLines(old=>old.map((x,j)=>i===j?{...x,quantity:Number(e.target.value)}:x))} required/></label>{lines.length>1&&<button type="button" onClick={()=>setLines(old=>old.filter((_,j)=>j!==i))}>{t("p0860")} {i+1}</button>}</fieldset>)}
 <button type="button" className="button" disabled={lines.length>=32} onClick={()=>setLines(old=>[...old,{id:crypto.randomUUID(),variant_id:"",quantity:1}])}>{t("p0861")}</button>
 <button className="button" disabled={!evidence}>{t("p0862")}</button></fieldset></form></section>}
 <section className="card"><h2>{t("p0863")}</h2><form onSubmit={e=>{e.preventDefault();void consult()}}><GuidedRecordSelect organization={organization} group="supplyOrders" name="reference:id" label="Orden de abastecimiento" disabled={blocked} value={id} onChange={nextValue=>{setID(nextValue);setPlan(null)}}/><button className="button" disabled={busy}>{t("p0863")}</button></form>
 {plan&&<><p>{t("p0865")} <strong>{labels[plan.state]}</strong>{t("p0674")} {plan.version}{t("p0008")}</p><p>{t("p0866")} {plan.supplier_id}{t("p0867")} {plan.factory_organization_id}{t("p0868")} {plan.destination_organization_id}{t("p0008")}</p><p>{t("p0275")} {plan.total_minor_units} {t("p0209")} {plan.currency}{t("p0869")} {plan.demand_reference}{t("p0008")}</p><p><a href={`/supply?order=${encodeURIComponent(plan.purchase_order_id)}&surface=factory`}>{t("p0870")}</a></p><ul>{plan.lines.map(l=><li key={l.id}>{l.variant_id}{t("p0213")} {l.quantity} {t("p0871")} {l.id}</li>)}</ul>
 {surface==="franchise"&&can("supply:plan")&&<><button className="button" disabled={blocked||!evidence||plan.state!=="draft"} onClick={()=>command("submit")}>{t("p0872")}</button><button className="button" disabled={blocked||!evidence||!["draft","submitted","accepted"].includes(plan.state)} onClick={()=>command("cancel")}>{t("p0873")}</button></>}
 {surface==="factory"&&can("supply:factory")&&<><button className="button" disabled={blocked||!evidence||plan.state!=="submitted"} onClick={()=>command("confirm")}>{t("p0874")}</button><button className="button" disabled={blocked||!evidence||plan.state!=="accepted"} onClick={()=>command("start")}>{t("p0875")}</button>
 <form onSubmit={e=>{e.preventDefault();command("register",{line_id:line,serial_number:serial,...(vin?{vin}:{}),...(battery?{battery_serial_number:battery}:{})});}}><fieldset disabled={blocked||plan.state!=="in-production"}><legend>{t("p0876")}</legend><label>{t("p0877")}<select value={line} onChange={e=>setLine(e.target.value)}>{plan.lines.map(l=><option key={l.id} value={l.id}>{l.variant_id} {t("p0016")} {l.id}</option>)}</select></label><label>{t("p0878")}<input value={serial} onChange={e=>setSerial(e.target.value)} maxLength={128} required/></label><label>{t("p0879")}<input value={vin} onChange={e=>setVIN(e.target.value)} maxLength={128}/></label><label>{t("p0880")}<input value={battery} onChange={e=>setBattery(e.target.value)} maxLength={128}/></label><button className="button" disabled={!evidence}>{t("p0881")}</button></fieldset></form></>}
 <AuthorizedUnitCodePicker units={plan.units.map(unit=>({id:unit.id,serial_number:unit.serial_number,product:plan.lines.find(item=>item.id===unit.line_id)?.variant_id??"",organization,selectable:(surface==="factory"&&can("supply:factory")&&unit.state==="released")||(surface==="franchise"&&can("supply:receive")&&unit.state==="shipped")}))} complete={completeUnitPage} selectedIds={selected} disabled={blocked} onChoose={unit=>setSelected(previous=>previous.includes(unit)?previous:[...previous,unit])}/>
 <div className="grid">{plan.units.map(u=><article className="card" key={u.id} aria-label={`Serie ${u.serial_number}`}><h3>{u.serial_number}</h3><p>{t("p0882")} {labels[u.state]}{t("p0008")} {u.stock_state?`Stock: ${labels[u.stock_state]??u.stock_state}.`:""}</p>{u.receipt_review_state&&<p>{t("p0883")} {u.receipt_review_state==="pending"?t("p0797"):u.receipt_review_state==="approved"?t("p0798"):t("p0337")}{t("p0008")}</p>}{u.shipment_id&&<p>{t("p0884")} {u.shipment_id}</p>}
 {((surface==="factory"&&can("supply:factory")&&u.state==="released")||(surface==="franchise"&&can("supply:receive")&&u.state==="shipped"))&&<label><input type="checkbox" disabled={blocked} checked={selected.includes(u.id)} onChange={e=>setSelected(old=>e.target.checked?[...old,u.id]:old.filter(x=>x!==u.id))}/>{t("p0885")} {u.serial_number}</label>}
 {surface==="factory"&&can("supply:factory")&&<>{u.state==="planned"&&<button disabled={blocked||!evidence} onClick={()=>command("milestone",{unit_id:u.id,target_state:"assembly"})}>{t("p0886")}</button>}{u.state==="assembly"&&<button disabled={blocked||!evidence} onClick={()=>command("milestone",{unit_id:u.id,target_state:"quality"})}>{t("p0887")}</button>}{u.state==="quality"&&<><p>{t("p0888")}</p><button disabled={blocked||!evidence||!reason.trim()||u.factory_review_state!=="pending"||u.factory_requester===subject} onClick={()=>command("milestone",{unit_id:u.id,target_state:"released",reason})}>{t("p0889")}</button><button disabled={blocked||!evidence||!reason.trim()||u.factory_review_state!=="pending"||u.factory_requester===subject} onClick={()=>command("milestone",{unit_id:u.id,target_state:"rejected",reason})}>{t("p0890")}</button></>}</>}
 {surface==="franchise"&&u.stock_state==="quarantine"&&<><p>{t("p0891")}</p>{can("supply:release")&&<><button disabled={blocked||!evidence||!reason.trim()||u.receipt_review_state!=="pending"||u.receipt_requester===subject} onClick={()=>command("quality",{unit_id:u.id,reason})}>{t("p0892")}</button><button disabled={blocked||!evidence||!reason.trim()||u.receipt_review_state!=="pending"||u.receipt_requester===subject} onClick={()=>command("quality-reject",{unit_id:u.id,reason})}>{t("p0893")}</button></>}{can("supply:inspect")&&<button disabled={blocked||!evidence||!reason.trim()||u.receipt_review_state!=="rejected"} onClick={()=>command("reinspect",{unit_id:u.id,reason})}>{t("p0894")}</button>}</>}
 </article>)}</div>
 {plan.next_unit_id&&<button className="button" disabled={busy} onClick={()=>void consult(plan.next_unit_id)}>{t("p0895")}</button>}<button className="button" disabled={busy} onClick={()=>void consult()}>{t("p0896")}</button>
 {((surface==="factory"&&can("supply:factory"))||(surface==="franchise"&&can("supply:receive")))&&<form onSubmit={e=>{e.preventDefault();command(surface==="factory"?"ship":"receive",{shipment_id:shipment,units:selected})}}><fieldset disabled={blocked}><legend>{surface==="factory"?t("p0897"):t("p0898")}</legend><label>{t("p0899")}<input value={shipment} onChange={e=>setShipment(e.target.value)} maxLength={128} required/></label><p>{selected.length} {t("p0900")}</p><button className="button" disabled={!evidence||selected.length===0}>{surface==="factory"?t("p0901"):t("p0902")}</button></fieldset></form>}
 </>}
 </section>
 </>;
}
