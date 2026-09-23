"use client";
import { GuidedRecordSelect } from "@/components/guided-record-select";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

import {OperationalGuide} from "@/components/operational-guide";
import {RELEASE_GUIDES} from "@/platform/help/content";
import Link from "next/link";
import {useEffect,useRef,useState,type FormEvent} from "react";
import {networkCommand,networkEntity,networkReceipt,networkReference,networkMatches,pendingNetwork,networkTargets,type NetworkEntity,type PendingNetwork} from "@/platform/network/contract";
const endpoint="/api/enterprise/network";


export function NetworkWorkspace({scope,subject,permissions,organizations}:{scope:string;subject:string;permissions:string[];organizations:string[]}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();
const states:Record<string,string>={provisioning:t("p0057"),active:t("p0705"),suspended:t("p0706"),closed:t("p0707"),draft:t("p0648"),terminated:t("p0708"),expired:t("p0709")};
const actions:Record<string,string>={active:t("p0710"),suspended:t("p0711"),closed:t("p0712"),terminated:t("p0713"),expired:t("p0714")};

 const can=(p:string)=>permissions.includes("*")||permissions.includes(p),canOrg=(id:string)=>permissions.includes("*")||organizations.includes(id);
 const held=useRef(false),saved=useRef<PendingNetwork|null>(null);const[ready,setReady]=useState(false),[busy,setBusy]=useState(false),[notice,setNotice]=useState(""),[pending,setPending]=useState<PendingNetwork|null>(null);
 const[kind,setKind]=useState("organization"),[id,setID]=useState(""),[org,setOrg]=useState(organizations[0]??""),[entity,setEntity]=useState<NetworkEntity|null>(null),[current,setCurrent]=useState(false);
 const[parent,setParent]=useState(organizations[0]??""),[type,setType]=useState("franchisee"),[code,setCode]=useState(""),[name,setName]=useState("");
 const[agreementOrg,setAgreementOrg]=useState(organizations[0]??""),[territory,setTerritory]=useState(""),[terms,setTerms]=useState(""),[start,setStart]=useState(""),[end,setEnd]=useState("");
 const key=`elite-network:${scope}`;
 const [selectionRevision,setSelectionRevision]=useState(0);
 function persist(v:PendingNetwork|null){if(v){const raw=JSON.stringify(pendingNetwork.parse(v));sessionStorage.setItem(key,raw);if(sessionStorage.getItem(key)!==raw)throw new Error("storage")}else sessionStorage.removeItem(key);saved.current=v;setPending(v)}
 useEffect(()=>{try{const raw=sessionStorage.getItem(key);if(raw){const v=pendingNetwork.parse(JSON.parse(raw));saved.current=v;setPending(v)}setReady(true)}catch{setNotice(t("p0651"))}},[key]);
 const blocked=!ready||busy||pending!==null;
 async function get(kind:string,id:string,org:string){const r=await fetch(endpoint+"?"+new URLSearchParams({kind,id,scope_organization_id:org}),{cache:"no-store",signal:AbortSignal.timeout(10000)});if(!r.ok)throw new Error("unavailable");return r.json()as Promise<unknown>}
 function show(v:NetworkEntity,isCurrent:boolean){setEntity(v);setKind(v.kind);setID(v.id);setOrg(v.organization_id);setCurrent(isCurrent)}
 async function consult(){if(held.current)return;held.current=true;setBusy(true);setCurrent(false);setEntity(null);try{const v=networkEntity.parse(await get(kind,id,kind==="organization"?id:org));if(v.id!==id||v.kind!==kind||v.organization_id!==(kind==="organization"?id:org))throw new Error("scope");show(v,true);setNotice(t("p0715"))}catch{setNotice(t("p0716"))}finally{held.current=false;setBusy(false)}}
 function accept(raw:unknown,p:PendingNetwork){const v=networkReceipt.parse(raw);if(!networkMatches(v,p,subject))throw new Error("receipt");persist(null);show(v.entity,false);setSelectionRevision(value=>value+1);return v.entity}
 async function send(input:unknown){
  if(held.current||!ready||saved.current)return;held.current=true;setBusy(true);let sent=false;
  try{const c=networkCommand.parse(input),p=await networkReference(c);persist(p);sent=true;
   const r=await fetch(endpoint,{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify(c),signal:AbortSignal.timeout(10000)});
   if(!r.ok){if([400,401,403,413,415].includes(r.status)){persist(null);sent=false}throw new Error("unconfirmed")}
   accept(await r.json(),p);sent=false;setNotice(t("p0717"));
  }catch{setNotice(sent?t("p0660"):t("p0718"))}finally{held.current=false;setBusy(false)}
 }
 async function recover(){const p=saved.current;if(!p||held.current)return;held.current=true;setBusy(true);try{accept(await get("result",p.command_id,p.scope_organization_id),p);setNotice(t("p0719"))}catch{setNotice(t("p0663"))}finally{held.current=false;setBusy(false)}}
 const createOrganization=(e:FormEvent)=>{e.preventDefault();void send({command_id:crypto.randomUUID(),action:"create-organization",entity_id:crypto.randomUUID(),scope_organization_id:["enterprise","franchisor"].includes(type)?"":parent,code,display_name:name,type})};
 const createAgreement=(e:FormEvent)=>{e.preventDefault();void send({command_id:crypto.randomUUID(),action:"create-agreement",entity_id:crypto.randomUUID(),scope_organization_id:agreementOrg,territory_code:territory,terms_version:terms,starts_on:start,...(end?{ends_on:end}:{})})};
 return <>
 <p role="status" aria-live="polite">{notice}</p>
 <OperationalGuide className="card" guide={RELEASE_GUIDES["network-role-view"]}/>
 {pending&&<section className="card"><h2>{t("p0166")}</h2><p>{t("p0167")} <code>{pending.command_id}</code></p><button className="button" disabled={busy} onClick={()=>void recover()}>{t("p0168")}</button></section>}
 {can("network:admin")&&<section className="card"><h2>{t("p0720")}</h2><form onSubmit={createOrganization}><fieldset disabled={blocked}><label>{t("p0721")}<select value={type} onChange={e=>setType(e.target.value)}>{(can("network:bootstrap")?["enterprise","franchisor","franchisee","factory","warehouse","store","service_center"]:["franchisee","factory","warehouse","store","service_center"]).map(v=><option key={v} value={v}>{({enterprise:t("p0722"),franchisor:t("p0723"),franchisee:t("p0724"),factory:t("p0121"),warehouse:t("p0725"),store:t("p0726"),service_center:t("p0727")}as Record<string,string>)[v]}</option>)}</select></label>
 {!["enterprise","franchisor"].includes(type)&&<GuidedRecordSelect key={selectionRevision+"reference:parent"} organization={organizations[0]??""} group="networkOrganizations" name="reference:parent" label="Organización superior" disabled={busy} value={parent} onChange={nextValue=>setParent(nextValue)}/>}
 <label>{t("p0729")}<input value={code} onChange={e=>setCode(e.target.value)} maxLength={128} pattern="[a-z][a-z0-9]*(-[a-z0-9]+)*" required/></label><p>{t("p0730")}</p><label>{t("p0731")}<input value={name} onChange={e=>setName(e.target.value)} maxLength={100} required/></label>
 <button className="button" disabled={!["enterprise","franchisor"].includes(type)&&!canOrg(parent)}>{t("p0720")}</button></fieldset></form></section>}
 {can("franchise:write")&&<section className="card"><h2>{t("p0732")}</h2><form onSubmit={createAgreement}><fieldset disabled={blocked}>
 <GuidedRecordSelect key={selectionRevision+"reference:agreementOrg"} organization={organizations[0]??""} group="networkOrganizations" name="reference:agreementOrg" label="Franquicia del acuerdo" disabled={busy} value={agreementOrg} onChange={nextValue=>setAgreementOrg(nextValue)}/><GuidedRecordSelect key={selectionRevision+"reference:territory"} organization={organizations[0]??""} group="networkTerritories" name="reference:territory" label="Territorio" disabled={busy} value={territory} onChange={nextValue=>setTerritory(nextValue)}/><GuidedRecordSelect key={selectionRevision+"reference:terms"} organization={organizations[0]??""} group="networkTerms" name="reference:terms" label="Condiciones del acuerdo" disabled={busy} value={terms} onChange={nextValue=>setTerms(nextValue)}/><label>{t("p0736")}<input type="date" value={start} onChange={e=>setStart(e.target.value)} required/></label><label>{t("p0737")}<input type="date" value={end} onChange={e=>setEnd(e.target.value)}/></label><button className="button" disabled={!canOrg(agreementOrg)}>{t("p0738")}</button>
 </fieldset></form></section>}
 <section className="card"><h2>{t("p0739")}</h2><form onSubmit={e=>{e.preventDefault();void consult()}}><fieldset disabled={busy}><label>{t("p0740")}<select value={kind} onChange={e=>{setKind(e.target.value);setEntity(null);setCurrent(false)}}><option value="organization">{t("p0091")}</option><option value="agreement">{t("p0741")}</option></select></label><GuidedRecordSelect key={selectionRevision+"reference:id"} organization={organizations[0]??""} group="networkEntities" name="reference:id" label="Registro de la red" disabled={busy} value={id} onChange={nextValue=>{setID(nextValue);setEntity(null);setCurrent(false)}}/>{kind==="agreement"&&<GuidedRecordSelect key={selectionRevision+"reference:org"} organization={organizations[0]??""} group="networkOrganizations" name="reference:org" label="Organización del registro" disabled={busy} value={org} onChange={nextValue=>{setOrg(nextValue);setEntity(null);setCurrent(false)}}/>}<button className="button">{t("p0307")}</button></fieldset></form>
 {entity&&<article><h3>{entity.kind==="organization"?t("p0744"):t("p0745")}</h3><p>{t("p0167")} <code>{entity.id}</code>{t("p0008")}</p><p>{t("p0247")} {states[entity.state]??entity.state}{t("p0674")} {entity.version}{t("p0008")}</p>{entity.display_name&&<p>{entity.display_name} {t("p0016")} {entity.code}</p>}{entity.territory_code&&<p>{t("p0734")} {entity.territory_code}{t("p0746")} {entity.terms_version}{t("p0747")} {entity.starts_on} {t("p0748")} {entity.ends_on??t("p0749")}{t("p0008")}</p>}
 {!current&&<p>{t("p0750")}</p>}
 {networkTargets(entity).map(target=><button className="button" key={target} disabled={blocked||!current||!canOrg(entity.organization_id)||!can(entity.kind==="organization"?"network:admin":"franchise:write")} onClick={()=>void send({command_id:crypto.randomUUID(),action:"transition-"+entity.kind,entity_id:entity.id,scope_organization_id:entity.organization_id,current:entity.state,target,version:entity.version})}>{actions[target]} {entity.kind==="organization"?t("p0751"):"acuerdo"}</button>)}
 </article>}
 </section>
 </>;
}
