"use client";
import { GuidedRecordSelect } from "@/components/guided-record-select";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

import {OperationalGuide} from "@/components/operational-guide";
import {RELEASE_GUIDES} from "@/platform/help/content";
import Link from"next/link";
import{useEffect,useRef,useState,type FormEvent}from"react";
import{helpCommand,helpArticle,helpPage,helpReceipt,pendingHelp,helpReference,helpMatches,type HelpArticle,type PendingHelp}from"@/platform/help/cms-contract";
const endpoint="/api/enterprise/help/cms";
export function HelpCMSWorkspace({scope,subject,permissions,organizations}:{scope:string;subject:string;permissions:string[];organizations:string[]}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();
const labels:Record<string,string>={draft:t("p0648"),published:t("p0649"),archived:t("p0650")};

 const can=(p:string)=>permissions.includes("*")||permissions.includes(p),editor=can("help:write")||can("help:publish"),canOrg=(v:string)=>can("*")||organizations.includes(v);
 const held=useRef(false),saved=useRef<PendingHelp|null>(null);const[ready,setReady]=useState(false),[busy,setBusy]=useState(false),[notice,setNotice]=useState(""),[pending,setPending]=useState<PendingHelp|null>(null);
 const[org,setOrg]=useState(organizations[0]??""),[locale,setLocale]=useState("es"),[query,setQuery]=useState(""),[items,setItems]=useState<ReturnType<typeof helpPage.parse>|null>(null);
 const[id,setID]=useState(""),[article,setArticle]=useState<HelpArticle|null>(null),[current,setCurrent]=useState(false),[history,setHistory]=useState<ReturnType<typeof helpPage.parse>|null>(null);
 const[category,setCategory]=useState("operations"),[title,setTitle]=useState(""),[body,setBody]=useState(""),key=`elite-help-cms:${scope}`;
 function persist(v:PendingHelp|null){if(v){const raw=JSON.stringify(pendingHelp.parse(v));sessionStorage.setItem(key,raw);if(sessionStorage.getItem(key)!==raw)throw new Error("storage")}else sessionStorage.removeItem(key);saved.current=v;setPending(v)}
 useEffect(()=>{try{const raw=sessionStorage.getItem(key);if(raw){const v=pendingHelp.parse(JSON.parse(raw));saved.current=v;setPending(v)}setReady(true)}catch{setNotice(t("p0651"))}},[key]);
 const blocked=!ready||busy||pending!==null;
 async function get(params:Record<string,string>){const r=await fetch(endpoint+"?"+new URLSearchParams(params),{cache:"no-store",signal:AbortSignal.timeout(10000)});if(!r.ok)throw new Error("unavailable");return r.json()as Promise<unknown>}
 function show(v:HelpArticle,isCurrent:boolean){setArticle(v);setOrg(v.organization_id);setID(v.id);setCurrent(isCurrent);setTitle(v.title);setBody(v.body);setCategory(v.category);setLocale(v.locale);setHistory(null)}
 async function read(selected=id,version=""){if(held.current)return;held.current=true;setBusy(true);setArticle(null);setCurrent(false);try{const v=helpArticle.parse(await get({kind:"article",organization_id:org,id:selected,...(version?{version}:{})}));if(v.id!==selected||v.organization_id!==org||version&&v.version!==version)throw new Error("scope");show(v,!version);setNotice(version?t("p0652"):t("p0653"))}catch{setNotice(t("p0654"))}finally{held.current=false;setBusy(false)}}
 async function search(after=""){if(held.current)return;held.current=true;setBusy(true);try{const v=helpPage.parse(await get({kind:"list",organization_id:org,locale,q:query,...(after?{after}:{})}));setItems(v);setNotice(t("p0655"))}catch{setItems(null);setNotice(t("p0656"))}finally{held.current=false;setBusy(false)}}
 async function versions(before=""){if(!article||held.current)return;held.current=true;setBusy(true);try{setHistory(helpPage.parse(await get({kind:"history",organization_id:article.organization_id,id:article.id,...(before?{before}:{})})));setNotice(t("p0657"))}catch{setHistory(null);setNotice(t("p0658"))}finally{held.current=false;setBusy(false)}}
 async function accept(raw:unknown,p:PendingHelp){const v=helpReceipt.parse(raw);if(!await helpMatches(v,p,subject))throw new Error("receipt");persist(null);show(v.article,false)}
 async function send(input:unknown){if(held.current||!ready||saved.current)return;held.current=true;setBusy(true);let sent=false;try{const c=helpCommand.parse(input),p=await helpReference(c);persist(p);sent=true;const r=await fetch(endpoint,{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify(c),signal:AbortSignal.timeout(10000)});if(!r.ok){if([400,401,403,413,415].includes(r.status)){persist(null);sent=false}throw new Error("unconfirmed")};await accept(await r.json(),p);sent=false;setNotice(t("p0659"))}catch{setNotice(sent?t("p0660"):t("p0661"))}finally{held.current=false;setBusy(false)}}
 async function recover(){const p=saved.current;if(!p||held.current)return;held.current=true;setBusy(true);try{await accept(await get({kind:"result",organization_id:p.organization_id,id:p.command_id}),p);setNotice(t("p0662"))}catch{setNotice(t("p0663"))}finally{held.current=false;setBusy(false)}}
 function save(e:FormEvent){e.preventDefault();if(article){void send({command_id:crypto.randomUUID(),action:"update",article_id:article.id,organization_id:article.organization_id,version:article.version,title,body})}else void send({command_id:crypto.randomUUID(),action:"create",article_id:crypto.randomUUID(),organization_id:org,locale,category,title,body})}
 return <>
 <p role="status" aria-live="polite">{notice}</p>
 <OperationalGuide className="card" guide={RELEASE_GUIDES["help-cms-view"]}/>
 {pending&&<section className="card"><h2>{t("p0166")}</h2><p>{t("p0167")} <code>{pending.command_id}</code></p><button className="button" disabled={busy} onClick={()=>void recover()}>{t("p0168")}</button></section>}
 <section className="card"><h2>{t("p0664")}</h2><fieldset disabled={busy}>
 <GuidedRecordSelect organization={organizations[0]??""} group="helpOrganizations" name="reference:org" label="Organización" disabled={busy} value={org} onChange={nextValue=>{setOrg(nextValue);setArticle(null);setCurrent(false);setItems(null);setHistory(null);setTitle("");setBody("")}} optional/>
 <label>{t("p0666")}<select value={locale} onChange={e=>{setLocale(e.target.value);setItems(null);setArticle(null);setID("");setCurrent(false);setHistory(null);setTitle("");setBody("")}}><option value="es">{t("p0667")}</option><option value="en">{t("p0668")}</option></select></label>
 <form onSubmit={e=>{e.preventDefault();void search()}}><label>{t("p0669")}<input value={query} onChange={e=>{setQuery(e.target.value);setItems(null)}} maxLength={160}/></label><button className="button" disabled={!canOrg(org)}>{t("p0670")}</button></form>
 {items&&<><ul>{items.items.map(a=><li key={a.id}><button className="button" onClick={()=>void read(a.id)}>{a.title} {t("p0016")} {labels[a.state]}</button></li>)}</ul>{items.next&&<button className="button" onClick={()=>void search(items.next)}>{t("p0671")}</button>}</>}
 <form onSubmit={e=>{e.preventDefault();void read()}}><label>{t("p0672")}<input value={id} onChange={e=>{setID(e.target.value);setArticle(null);setCurrent(false);setHistory(null)}} maxLength={128} required/></label><button className="button" disabled={!canOrg(org)}>{t("p0673")}</button></form>
 </fieldset></section>
 {article&&<article className="card"><h2 lang={article.locale}>{article.title}</h2><p>{t("p0247")} {labels[article.state]}{t("p0674")} {article.version}{t("p0675")} {article.locale}{t("p0008")}</p><p>{t("p0167")} <code>{article.id}</code>{t("p0008")}</p><p lang={article.locale} style={{whiteSpace:"pre-wrap",overflowWrap:"anywhere"}}>{article.body}</p>{!current&&<p>{t("p0676")}</p>}
 {can("help:publish")&&["draft","published"].includes(article.state)&&<button className="button" disabled={blocked||!current} onClick={()=>void send({command_id:crypto.randomUUID(),action:article.state==="draft"?"publish":"archive",article_id:article.id,organization_id:article.organization_id,version:article.version})}>{article.state==="draft"?t("p0677"):t("p0678")}</button>}
 {editor&&<button className="button" disabled={busy} onClick={()=>void versions()}>{t("p0679")}</button>}
 {history&&<><ul>{history.items.map(a=><li key={a.version}><button className="button" disabled={busy} onClick={()=>void read(a.id,a.version)}>{t("p0680")} {a.version} {t("p0016")} {labels[a.state]}</button></li>)}</ul>{history.next&&<button className="button" disabled={busy} onClick={()=>void versions(history.next)}>{t("p0681")}</button>}</>}
 </article>}
 {can("help:write")&&<section className="card"><h2>{article?t("p0682"):t("p0683")}</h2><button className="button" disabled={blocked} onClick={()=>{setArticle(null);setID("");setTitle("");setBody("");setCurrent(false);setHistory(null)}}>{t("p0684")}</button>
 <form onSubmit={save}><fieldset disabled={blocked||article!==null&&(!current||article.state!=="draft")}>
 {!article&&<label>{t("p0685")}<input value={category} onChange={e=>setCategory(e.target.value)} pattern="[a-z][a-z0-9_-]*" maxLength={64} required/></label>}
 <label>{t("p0686")}<input value={title} onChange={e=>setTitle(e.target.value)} maxLength={200} required/></label><label>{t("p0687")}<textarea value={body} onChange={e=>setBody(e.target.value)} maxLength={16384} rows={8} required/></label><button className="button" disabled={!canOrg(org)}>{article?t("p0688"):t("p0689")}</button></fieldset></form></section>}
 </>
}
