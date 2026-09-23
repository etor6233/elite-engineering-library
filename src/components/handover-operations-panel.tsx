"use client";
import {usePrivateI18n} from "@/platform/i18n/private-provider";


// AUTHORED orchestration UI; uses existing checklist and customer acceptance.
import { useEffect,useRef,useState } from "react";
import { OperationsList } from "@/components/admin-ops/operations-list";
import { OperationDetail } from "@/components/admin-ops/operation-detail";
import { adminText, adminState } from "@/components/admin-ops/messages";
import styles from "@/components/admin-ops/operations.module.css";
import { ChecklistCompletionPanel } from "@/components/franchise-command-panel";
import { commercialReceiptSchema,currentReleaseSchema,handoverContextSchema,handoverViewSchema,operationMarkerSchema,retainOperation,type CommercialReceipt,type CurrentRelease,type HandoverContext,type OperationMarker } from "@/platform/handovers/contracts";

const endpoint="/api/enterprise/handovers";
type DeliveryOrder = { id: string; state?: string; lines?: { name: string; quantity: number }[]; handovers?: { id: string; state: string }[] };
function orderTitle(order: DeliveryOrder, language: string) { return order.lines?.length ? order.lines.map(line => line.name).join(" · ") : `${adminText("referenceOrder", language)} ${order.id.length > 18 ? order.id.slice(0, 8) + "…" + order.id.slice(-6) : order.id}`; }
export function HandoverOperationsPanel({orders,organization,scope}:{orders:DeliveryOrder[]|null;organization:string;scope:string}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();
 const [selected,setSelected]=useState<string|null>(null),[visited,setVisited]=useState<string[]>([]);
 const listRef=useRef<HTMLDivElement>(null);
 const language=privateLocale.language;
 const choose=(id:string)=>{if(!orders?.some(order=>order.id===id))return;setSelected(id);setVisited(previous=>previous.includes(id)?previous:[...previous,id]);};
 const back=()=>{const prior=selected;setSelected(null);requestAnimationFrame(()=>{const row=Array.from(listRef.current?.querySelectorAll<HTMLButtonElement>("button[data-record-id]")??[]).find(item=>item.dataset.recordId===prior);row?.focus({preventScroll:true});});};
 return <section className={styles.workspace} aria-label={t("p0617")} data-testid="handover-workspace"><h2 className={styles.srOnly}>{adminText("deliveries",language)}</h2>
 {orders===null?<div className={styles.contextState} role="alert"><p>{adminText("unavailableList",language)}</p></div>:orders.length===0?<div className={styles.empty} role="status"><p>{adminText("noDeliveries",language)}</p></div>:<div className={styles.split} data-detail-open={Boolean(selected)}>
 <div className={styles.listColumn} ref={listRef}><OperationsList items={orders.map(order=>({id:order.id,title:orderTitle(order,language),subtitle:`${adminText("referenceOrder",language)} ${order.id}`,state:adminState(order.handovers?.[0]?.state??order.state,language)}))} selected={selected} onSelect={choose} language={language} label={adminText("deliveries",language)} testId="handover-list"/></div>
 <div>{selected===null?<div className={styles.placeholder}><h3>{adminText("chooseDelivery",language)}</h3><p>{adminText("oneSelection",language)}</p></div>:null}
 {orders.filter(order=>visited.includes(order.id)).map(order=><OperationDetail key={`${scope}:${order.id}`} active={selected===order.id} title={orderTitle(order,language)} reference={order.id} language={language} onBack={back} testId="handover-detail"><HandoverOrder orderId={order.id} organization={organization} scope={scope}/></OperationDetail>)}</div>
 </div>}
 </section>;
}

function HandoverOrder({orderId,organization,scope}:{orderId:string;organization:string;scope:string}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

 const [context,setContext]=useState<HandoverContext|null>(null),[message,setMessage]=useState(""),[ready,setReady]=useState(false),[busy,setBusy]=useState(false),[storageOK,setStorageOK]=useState(false);
 const [contextState,setContextState]=useState<"loading"|"ready"|"error"|"forbidden">("loading");
 const [markers,setMarkers]=useState<Partial<Record<"prepare"|"release",OperationMarker>>>({}),[receipt,setReceipt]=useState<CommercialReceipt|null>(null),[current,setCurrent]=useState<CurrentRelease|null>(null);
 const operations=useRef<Partial<Record<"prepare"|"release",OperationMarker>>>({}),fence=useRef(false);
 const storageKey=(action:string)=>`elite-handover:${scope}:${orderId}:${action}`;
 async function read(kind:string,extra:Record<string,string>={}){
  const response=await fetch(endpoint+"?"+new URLSearchParams({kind,organizationId:organization,orderId,...extra}),{cache:"no-store",signal:AbortSignal.timeout(10000)});
  if(!response.ok)throw new Error(response.status===401||response.status===403?"forbidden":"unavailable");return response.json() as Promise<unknown>;
 }
 async function refresh(){
  setCurrent(null);setContextState("loading");
  try{const v=handoverContextSchema.parse(await read("context"));if(v.organization_id!==organization||v.order_id!==orderId)throw new Error("scope");setContext(v);setContextState("ready");setMessage(t("p0622"))}
  catch(error){setContext(null);setContextState(error instanceof Error&&error.message==="forbidden"?"forbidden":"error");setMessage(t("p0623"))}
 }
 useEffect(()=>{
  let mounted=true;
  try{
   if(!/^[a-f0-9]{64}$/.test(scope))throw new Error("scope");const values:Partial<Record<"prepare"|"release",OperationMarker>>={};
   for(const action of ["prepare","release"] as const){const raw=sessionStorage.getItem(storageKey(action));if(raw!==null){if(raw.length>1024)throw new Error("reference");const v=operationMarkerSchema.parse(JSON.parse(raw));if(v.action!==action||v.orderId!==orderId)throw new Error("scope");values[action]=v}}
   operations.current=values;setMarkers(values);setStorageOK(true);
  }catch{setStorageOK(false);setMessage(t("p0624"))}
  setReady(true);
  void read("context").then(raw=>{if(!mounted)return;const v=handoverContextSchema.parse(raw);if(v.order_id!==orderId||v.organization_id!==organization)throw new Error("scope");setContext(v);setContextState("ready")}).catch(error=>{if(mounted){setContext(null);setContextState(error instanceof Error&&error.message==="forbidden"?"forbidden":"error")}});
  return()=>{mounted=false};
 // Identity changes remount this keyed component; no background polling.
 // eslint-disable-next-line react-hooks/exhaustive-deps
 },[scope,orderId,organization]);
 async function mutate(action:"prepare"|"release"){
  if(!ready||!storageOK||fence.current||operations.current[action]||!context)return;
  if(action==="prepare"?!context.can_prepare:context.handover?.state!=="accepted"||context.release_effect!=="COMMIT_COMMERCIAL_RELEASE_RECEIPT")return;
  fence.current=true;setBusy(true);setCurrent(null);
  let marker:OperationMarker;
  try{
   const requested:OperationMarker=action==="prepare"?{action,orderId,requestKey:crypto.randomUUID()}:{action,orderId,handoverId:context.handover!.id,requestKey:crypto.randomUUID()};
   if(sessionStorage.getItem(storageKey(action))!==null)throw new Error("existing operation");
   marker=retainOperation(sessionStorage,storageKey(action),requested);operations.current={...operations.current,[action]:marker};setMarkers(operations.current);
  }catch{setStorageOK(false);setBusy(false);fence.current=false;setMessage(t("p0625"));return}
  try{
   const response=await fetch(endpoint,{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify({action,organizationId:organization,orderId,requestKey:marker.requestKey}),signal:AbortSignal.timeout(10000)}),raw=await response.json();
   if(!response.ok)throw new Error("unconfirmed");
   if(action==="prepare"){
    const v=handoverViewSchema.parse(raw.handover);if(v.organization_id!==organization||v.order_id!==orderId)throw new Error("scope");
    setContext({...context,handover:v,can_prepare:false});setMessage(t("p0626"));
   }else{
    const v=commercialReceiptSchema.parse(raw.receipt);if(v.organization_id!==organization||v.order_id!==orderId||marker.action!=="release"||v.handover_id!==marker.handoverId)throw new Error("scope");setReceipt(v);setMessage(t("p0627"));
   }
  }catch{setMessage(t("p0628"))}
  finally{setBusy(false);fence.current=false}
 }
 async function recover(action:"prepare"|"release"){
  const marker=operations.current[action];if(!marker||fence.current)return;fence.current=true;setBusy(true);setCurrent(null);
  try{
   const raw=await read(action+"-result",{requestKey:marker.requestKey,...(marker.action==="release"?{handoverId:marker.handoverId}:{})}) as {handover?:unknown;receipt?:unknown};
   if(action==="prepare"){const h=handoverViewSchema.parse(raw.handover);if(h.order_id!==orderId||h.organization_id!==organization)throw new Error("scope");await refresh();setMessage(t("p0629", {handover_id:(h.id),handover_state:(h.state)}))}
   else{const r=commercialReceiptSchema.parse(raw.receipt);if(r.order_id!==orderId||r.organization_id!==organization||marker.action!=="release"||r.handover_id!==marker.handoverId)throw new Error("scope");setReceipt(r);setMessage(t("p0630"))}
  }catch{setMessage(t("p0631"))}
  finally{fence.current=false;setBusy(false)}
 }
 async function validate(){
  const id=receipt?.handover_id??context?.handover?.id??(operations.current.release?.action==="release"?operations.current.release.handoverId:undefined);if(!id||fence.current)return;
  fence.current=true;setBusy(true);setCurrent(null);
  try{const v=currentReleaseSchema.parse(await read("release-current",{handoverId:id}));if(v.receipt.organization_id!==organization||v.receipt.order_id!==orderId||v.receipt.handover_id!==id)throw new Error("scope");setReceipt(v.receipt);setCurrent(v)}catch{setMessage(t("p0632"))}
  finally{fence.current=false;setBusy(false)}
 }
 const language=privateLocale.language;
 const step=receipt?3:context?.handover?.state==="accepted"?3:context?.handover?.state==="presented"?2:context?.handover?1:0;
 return <article aria-label={`${adminText("deliveries",language)} ${orderId}`} data-testid="handover-context" data-state={contextState} aria-busy={contextState==="loading"||busy}>
 <ol className={styles.steps} aria-label={adminText("deliveries",language)}>{(["readyStep","reviewStep","customerStep","releaseStep"] as const).map((key,index)=><li key={key} aria-current={index===step?"step":undefined}>{index+1}. {adminText(key,language)}</li>)}</ol>
 {!context?<div className={styles.contextState} data-state={contextState} role={contextState==="loading"?"status":"alert"}><p>{adminText(contextState==="loading"?"loading":contextState==="forbidden"?"forbidden":"unavailable",language)}</p>{contextState!=="loading"?<button className={styles.secondary} type="button" disabled={!ready||busy} onClick={()=>void refresh()}>{adminText("retry",language)}</button>:null}</div>:<>
 <div className={styles.actions}>{context.handover?<span className={styles.state}>{adminState(context.handover.state,language)}</span>:<span className={styles.state}>{adminText("pendingPrepare",language)}</span>}<button className={styles.secondary} type="button" disabled={!ready||busy||contextState==="loading"} onClick={()=>void refresh()} data-testid="handover-refresh">{adminText("refresh",language)}</button></div>
 {context.handover?<><ChecklistCompletionPanel key={context.handover.id} organization={organization} scope={scope} initialHandover={{id:context.handover.id,version:context.handover.version}} recoveryOnly={context.handover.state!=="prepared"}/>{context.handover.state==="presented"?<p>{adminText("pendingAcceptance",language)}</p>:null}</>:<div className={styles.formActions}>{context.can_prepare?<button className="button" type="button" disabled={!ready||busy||!storageOK||Boolean(markers.prepare)} onClick={()=>void mutate("prepare")} data-testid="handover-prepare">{t("p0636")}</button>:<p>{adminText("noPermissionPrepare",language)}</p>}</div>}
 </>}
 <div className={styles.formActions}>
 {markers.prepare?<button className={styles.secondary} type="button" disabled={!ready||busy} onClick={()=>void recover("prepare")} data-testid="handover-recover-prepare">{t("p0637")}</button>:null}
 {context?.handover?.state==="accepted"&&context.release_effect==="COMMIT_COMMERCIAL_RELEASE_RECEIPT"?<button className="button" type="button" disabled={!ready||busy||!storageOK||Boolean(markers.release)} onClick={()=>void mutate("release")} data-testid="handover-release">{t("p0638")}</button>:context?.release_effect==="READ_ONLY_ELIGIBILITY"?<p>{t("p0639")}</p>:null}
 {markers.release?<button className={styles.secondary} type="button" disabled={!ready||busy} onClick={()=>void recover("release")} data-testid="handover-recover-release">{t("p0640")}</button>:null}
 {receipt||context?.handover?.state==="accepted"?<button className={styles.secondary} type="button" disabled={!ready||busy} onClick={()=>void validate()} data-testid="handover-validate">{t("p0641")}</button>:null}
 </div>
 {receipt?<details data-testid="handover-receipt"><summary>{adminText("receipt",language)}</summary><p>{t("p0642")} {receipt.id} {t("p0643")} {receipt.recorded_at}</p><p>{t("p0644")}</p></details>:null}
 {current?<p role="status">{current.current?t("p0645", {evaluated_at:(current.evaluated_at),valid_until:(current.receipt.valid_until)}):t("p0646", {evaluated_at:(current.evaluated_at)})}</p>:receipt?<p>{t("p0647")}</p>:null}
 <p className={styles.message} role="status" aria-live="polite" data-testid="handover-message">{message}</p>
 </article>;
}
