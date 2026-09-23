"use client";
import { GuidedRecordSelect } from "@/components/guided-record-select";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

import {OperationalGuide} from "@/components/operational-guide";
import {RELEASE_GUIDES} from "@/platform/help/content";
import {useEffect,useRef,useState,type FormEvent} from "react";
import {z} from "zod";
import {catalogCommand,catalogDraft,catalogReceipt,pendingCatalog,commandReference,catalogCanonical,catalogSHA,receiptMatches,type CatalogCommand,type CatalogDraft,type CatalogReceipt,type PendingCatalog} from "@/platform/catalog/authoring";
import {publishedCatalogSchema,type PublishedCatalog} from "@/platform/catalog/schema";
const endpoint="/api/enterprise/catalog";
const capsule=z.object({pending:pendingCatalog.nullable(),source:catalogReceipt.nullable(),media:catalogReceipt.nullable()}).strict();
type Capsule=z.infer<typeof capsule>;
type VariantInput={code:string;name:string;battery:string;amount:string;tax:"inclusive"|"exclusive"|"not-applicable"};
const emptyVariant=():VariantInput=>({code:"",name:"",battery:"",amount:"",tax:"not-applicable"});
const stages=["legal","technical","media","publication"] as const;

export function CatalogAuthoring({organization,tenant,subject,permissions,scope,initialDraft}:{organization:string;tenant:string;subject:string;permissions:string[];scope:string;initialDraft:string}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();
const stageLabel={legal:t("p0144"),technical:t("p0145"),media:t("p0146"),publication:t("p0147")};

 const can=(p:string)=>permissions.includes("*")||permissions.includes(p);
 const held=useRef(false),saved=useRef<Capsule>({pending:null,source:null,media:null});
 const [storageReady,setStorageReady]=useState(false),[busy,setBusy]=useState(false),[pending,setPending]=useState<PendingCatalog|null>(null),[notice,setNotice]=useState("");
 const [modelCode,setModelCode]=useState(""),[modelName,setModelName]=useState(""),[vehicleClass,setVehicleClass]=useState("bicycle"),[description,setDescription]=useState("");
 const [variants,setVariants]=useState<VariantInput[]>([emptyVariant()]),[from,setFrom]=useState(""),[until,setUntil]=useState("");
 const [modelID,setModelID]=useState(""),[mediaID,setMediaID]=useState(""),[bookID,setBookID]=useState(""),[file,setFile]=useState<File|null>(null);
 const [draftID,setDraftID]=useState(initialDraft),[draft,setDraft]=useState<CatalogDraft|null>(null),[current,setCurrent]=useState<PublishedCatalog|null>(null),[currentKnown,setCurrentKnown]=useState(false);
 const [stage,setStage]=useState<typeof stages[number]>("legal"),[reason,setReason]=useState(""),[evidence,setEvidence]=useState(""),[publishReason,setPublishReason]=useState("");
 const key=`elite-catalog:${scope}`;
 function persist(value:Capsule){const raw=JSON.stringify(capsule.parse(value));sessionStorage.setItem(key,raw);if(sessionStorage.getItem(key)!==raw)throw new Error("storage");saved.current=value;setPending(value.pending)}
 function accepted(v:CatalogReceipt){
  const next={...saved.current,pending:null,...(v.kind==="source"?{source:v}:{}),...(v.kind==="media"?{media:v}:{})};persist(next);
  if(v.kind==="source"){setModelID(v.resource_id);setBookID(v.source_price_book_id!)}
  if(v.kind==="media")setMediaID(v.resource_id);
  if(v.kind==="draft"){setDraftID(v.resource_id);setDraft(null)}
 }
 useEffect(()=>{try{
  const raw=sessionStorage.getItem(key);if(raw){const v=capsule.parse(JSON.parse(raw));if(v.source&&v.source.actor!==subject||v.media&&v.media.actor!==subject)throw new Error("scope");saved.current=v;setPending(v.pending);if(v.source){setModelID(v.source.resource_id);setBookID(v.source.source_price_book_id!)}if(v.media)setMediaID(v.media.resource_id)}
  setStorageReady(true);
 }catch{setNotice(t("p0148"))}},[key,subject]);
 async function get(kind:string,id?:string){const r=await fetch(endpoint+"?"+new URLSearchParams({kind,organization_id:organization,...(id?{id}:{})}),{cache:"no-store"});if(!r.ok)throw new Error("unavailable");return r.json() as Promise<unknown>}
 async function readDraftValue(id:string){const value=catalogDraft.parse(await get("draft",id));if(value.id!==id||value.snapshot.profile.tenant_id!==tenant||value.snapshot.profile.organization_id!==organization)throw new Error("scope");return value}
 async function loadDraft(){
  if(held.current)return;held.current=true;setBusy(true);setDraft(null);
  try{setDraft(await readDraftValue(draftID));setNotice(t("p0149"))}catch{setNotice(t("p0150"))}finally{held.current=false;setBusy(false)}
 }
 async function loadCurrent(){
  if(held.current)return;held.current=true;setBusy(true);setCurrentKnown(false);
  try{const v=z.object({publication:publishedCatalogSchema.nullable()}).strict().parse(await get("current"));setCurrent(v.publication);setCurrentKnown(true);setNotice(t("p0151"))}catch{setNotice(t("p0152"))}finally{held.current=false;setBusy(false)}
 }
 async function send(c:CatalogCommand|{action:"media";command_id:string;bytes:Uint8Array}){
  if(held.current||!storageReady||saved.current.pending)return;held.current=true;setBusy(true);setNotice("");
  let started=false;
  try{
   let ref:PendingCatalog,url=endpoint,body:BodyInit,type:string;
   if(c.action==="media"){
    const original=await catalogSHA(c.bytes);ref={command_id:c.command_id,kind:"media",request_sha256:await catalogSHA(catalogCanonical({command_id:c.command_id,original_sha256:original}))};
    url+="?"+new URLSearchParams({organization_id:organization,command_id:c.command_id});body=new Uint8Array(c.bytes).buffer;type="image/png";
   }else{const parsed=catalogCommand.parse(c);ref=await commandReference(parsed);body=JSON.stringify(parsed);type="application/json"}
   persist({...saved.current,pending:ref});started=true;
   const r=await fetch(url,{method:"POST",headers:{"content-type":type},body});
   if(!r.ok){
    if([400,401,403,413,415].includes(r.status)){persist({...saved.current,pending:null});started=false;throw new Error("rejected-before-write")}
    throw new Error("unconfirmed");
   }
   const v=catalogReceipt.parse(await r.json());if(!receiptMatches(v,ref,subject))throw new Error("receipt identity");
   accepted(v);started=false;setNotice(v.kind==="source"?t("p0153"):v.kind==="media"?t("p0154"):v.kind==="draft"?t("p0155"):v.kind==="review"?t("p0156"):t("p0157"));
   if(v.kind==="review"){try{setDraft(await readDraftValue(v.resource_id))}catch{setDraft(null);setNotice(t("p0158"))}}
   if(v.kind==="publish"){setCurrentKnown(false);setCurrent(null)}
  }catch{
   setNotice(started?t("p0159"):t("p0160"));
  }finally{held.current=false;setBusy(false)}
 }
 async function recover(){
  const p=saved.current.pending;if(!p||held.current)return;held.current=true;setBusy(true);
  try{const value=catalogReceipt.parse(await get("command",p.command_id));if(!receiptMatches(value,p,subject))throw new Error("identity");accepted(value);setNotice(t("p0161"));if(value.kind==="review"){setDraft(null);setDraftID(value.resource_id)}if(value.kind==="publish"){setCurrentKnown(false);setCurrent(null)}}
  catch{setNotice(t("p0162"))}finally{held.current=false;setBusy(false)}
 }
 const blocked=busy||!storageReady||pending!==null;
 const prevent=(e:FormEvent)=>e.preventDefault();
 const source=async(e:FormEvent)=>{prevent(e);try{
  const value=catalogCommand.parse({action:"source",organization_id:organization,command_id:crypto.randomUUID(),model:{id:"",code:modelCode,displayName:modelName,vehicleClass,specification:{description}},variants:variants.map(v=>({code:v.code,display_name:v.name,battery_specification:{description:v.battery},amount_minor_units:v.amount,tax_mode:v.tax})),valid_from:new Date(from+"Z").toISOString(),valid_until:until?new Date(until+"Z").toISOString():null});
  await send(value);
 }catch{setNotice(t("p0163"))}};
 function changeVariant(i:number,key:keyof VariantInput,value:string){setVariants(old=>old.map((v,j)=>j===i?{...v,[key]:value}:v))}
 const draftAction=(e:FormEvent)=>{prevent(e);void send({action:"draft",organization_id:organization,command_id:crypto.randomUUID(),price_book_id:bookID,models:[{model_id:modelID,media_id:mediaID}]})};
 async function upload(e:FormEvent){prevent(e);if(!file||file.type!=="image/png"||file.size<1||file.size>1048576){setNotice(t("p0164"));return}await send({action:"media",command_id:crypto.randomUUID(),bytes:new Uint8Array(await file.arrayBuffer())})}
 async function evidenceFile(v:File|null){setEvidence("");if(!v||v.size<1||v.size>1048576){setNotice(t("p0165"));return}setEvidence(await catalogSHA(new Uint8Array(await v.arrayBuffer())))}
 return <>
 <OperationalGuide className="card" guide={RELEASE_GUIDES["catalog-role-view"]}/>
 <p role="status" aria-live="polite">{notice}</p>
 {pending&&<section className="card" aria-label={t("p0166")}><h2>{t("p0166")}</h2><p>{t("p0167")} <code>{pending.command_id}</code></p><button className="button" disabled={busy} onClick={()=>void recover()}>{t("p0168")}</button></section>}
 {can("catalog:draft")&&<section className="card"><h2>{t("p0169")}</h2><p>{t("p0170")}</p>
 <form onSubmit={source}><fieldset disabled={blocked}>
 <label>{t("p0171")}<input value={modelCode} onChange={e=>setModelCode(e.target.value)} required maxLength={64}/></label>
 <label>{t("p0172")}<input value={modelName} onChange={e=>setModelName(e.target.value)} required maxLength={160}/></label>
 <label>{t("p0173")}<select value={vehicleClass} onChange={e=>setVehicleClass(e.target.value)}><option value="bicycle">{t("p0174")}</option><option value="motorcycle">{t("p0175")}</option><option value="scooter">{t("p0176")}</option><option value="utility">{t("p0177")}</option><option value="other">{t("p0178")}</option></select></label>
 <label>{t("p0179")}<textarea value={description} onChange={e=>setDescription(e.target.value)} required maxLength={4000}/></label>
 <label>{t("p0180")}<input type="datetime-local" value={from} onChange={e=>setFrom(e.target.value)} required/></label>
 <label>{t("p0181")}<input type="datetime-local" value={until} onChange={e=>setUntil(e.target.value)}/></label>
 {variants.map((v,i)=><fieldset key={i}><legend>{t("p0182")} {i+1}</legend>
 <label>{t("p0183")}<input value={v.code} onChange={e=>changeVariant(i,"code",e.target.value)} required maxLength={64}/></label>
 <label>{t("p0184")}<input value={v.name} onChange={e=>changeVariant(i,"name",e.target.value)} required maxLength={160}/></label>
 <label>{t("p0185")}<textarea value={v.battery} onChange={e=>changeVariant(i,"battery",e.target.value)} required maxLength={4000}/></label>
 <label>{t("p0186")}<input inputMode="numeric" pattern="[0-9]+" value={v.amount} onChange={e=>changeVariant(i,"amount",e.target.value)} required maxLength={19}/></label>
 <label>{t("p0187")}<select value={v.tax} onChange={e=>changeVariant(i,"tax",e.target.value)}><option value="not-applicable">{t("p0188")}</option><option value="inclusive">{t("p0189")}</option><option value="exclusive">{t("p0190")}</option></select></label>
 {variants.length>1&&<button type="button" onClick={()=>setVariants(old=>old.filter((_,j)=>i!==j))}>{t("p0191")} {i+1}</button>}
 </fieldset>)}
 <button type="button" className="button" disabled={variants.length>=16} onClick={()=>setVariants(old=>[...old,emptyVariant()])}>{t("p0192")}</button>
 <button className="button">{t("p0193")}</button>
 </fieldset></form></section>}
 {can("catalog:draft")&&<section className="card"><h2>{t("p0194")}</h2><form onSubmit={upload}><label>{t("p0195")}<input type="file" accept="image/png" disabled={blocked} onChange={e=>setFile(e.target.files?.[0]??null)} required/></label><button className="button" disabled={blocked||!file}>{t("p0196")}</button></form>
 <form onSubmit={draftAction}><fieldset disabled={blocked}>
 <GuidedRecordSelect organization={organization} group="catalogModels" name="reference:modelID" label="Producto del catálogo" disabled={blocked} value={modelID} onChange={nextValue=>setModelID(nextValue)}/>
 <GuidedRecordSelect organization={organization} group="catalogMedia" name="reference:mediaID" label="Imagen aprobada" disabled={blocked} value={mediaID} onChange={nextValue=>setMediaID(nextValue)}/>
 <GuidedRecordSelect organization={organization} group="catalogBooks" name="reference:bookID" label="Lista de precios" disabled={blocked} value={bookID} onChange={nextValue=>setBookID(nextValue)}/>
 <button className="button">{t("p0200")}</button></fieldset></form></section>}
 <section className="card"><h2>{t("p0201")}</h2>
 <form onSubmit={e=>{prevent(e);void loadDraft()}}><GuidedRecordSelect organization={organization} group="catalogDrafts" name="reference:draftID" label="Borrador del catálogo" disabled={blocked} value={draftID} onChange={nextValue=>{setDraftID(nextValue);setDraft(null)}}/><button className="button" disabled={busy}>{t("p0203")}</button></form>
 {draft&&<><p>{t("p0204")} {draft.maker}{t("p0008")}</p><a href={`/admin/catalog?draft=${encodeURIComponent(draft.id)}`}>{t("p0205")}</a>
 {draft.snapshot.models.map(m=><article key={m.id}><h3>{m.displayName}</h3><p>{t("p0206")} {m.code}{t("p0207")} {m.vehicleClass}{t("p0008")}</p><pre style={{whiteSpace:"pre-wrap",overflowWrap:"anywhere"}}>{JSON.stringify(m.specification,null,2)}</pre></article>)}
 {draft.snapshot.variants.map(v=><article key={v.id}><h3>{v.display_name}</h3><p>{t("p0208")} {v.amount_minor_units} {t("p0209")} {draft.snapshot.profile.currency}{t("p0210")} {v.tax_mode}{t("p0211")} {v.homologation_state}{t("p0008")}</p></article>)}
 <ul aria-label={t("p0212")}>{stages.map(s=><li key={s}>{stageLabel[s]}{t("p0213")} {draft.reviews[s]}</li>)}</ul>
 <details><summary>{t("p0214")}</summary><pre style={{whiteSpace:"pre-wrap",overflowWrap:"anywhere"}}>{JSON.stringify(draft.snapshot,null,2)}</pre><p>{t("p0215")} <code>{draft.sha256}</code></p></details>
 {draft.maker===subject?<p>{t("p0216")}</p>:<><label>{t("p0217")}<select value={stage} onChange={e=>setStage(e.target.value as typeof stage)}>{stages.filter(s=>can("catalog:review:"+s)).map(s=><option key={s} value={s}>{stageLabel[s]}</option>)}</select></label>
 <label>{t("p0218")}<input type="file" disabled={blocked} onChange={e=>void evidenceFile(e.target.files?.[0]??null)}/></label><p>{t("p0219")}</p>
 <label>{t("p0220")}<textarea value={reason} onChange={e=>setReason(e.target.value)} maxLength={2000}/></label>
 {[true,false].map(approved=><button key={String(approved)} className="button" disabled={blocked||!can("catalog:review:"+stage)||!evidence||!reason.trim()||draft.reviews[stage]!=="pending"||approved&&stage==="publication"&&stages.some(s=>s!=="publication"&&draft.reviews[s]!=="approved")} onClick={()=>void send({action:"review",organization_id:organization,command_id:crypto.randomUUID(),draft_id:draft.id,stage,snapshot_sha256:draft.sha256,approved,reason,evidence_sha256:evidence})}>{approved?t("p0221"):t("p0222")}</button>)}
 </>}
 </>}
 </section>
 <section className="card"><h2>{t("p0223")}</h2><button className="button" disabled={busy} onClick={()=>void loadCurrent()}>{t("p0224")}</button>
 {currentKnown&&(current?<><p>{t("p0225")} {current.generation}{t("p0008")}</p><ul>{current.models.map(m=><li key={m.id}>{m.displayName}</li>)}</ul><a href="/models">{t("p0226")}</a></>:<p>{t("p0227")}</p>)}
 {draft&&can("catalog:publish")&&draft.maker!==subject&&<><p>{t("p0228")}</p>
 <label>{t("p0229")}<textarea value={publishReason} onChange={e=>setPublishReason(e.target.value)} maxLength={2000}/></label>
 <button className="button" disabled={blocked||!currentKnown||!publishReason.trim()||stages.some(s=>draft.reviews[s]!=="approved")} onClick={()=>void send({action:"publish",organization_id:organization,command_id:crypto.randomUUID(),draft_id:draft.id,snapshot_sha256:draft.sha256,expected_generation:current?.generation??"0",reason:publishReason})}>{t("p0230")}</button>
 </>}
 </section>
 </>;
}
