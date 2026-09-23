"use client";
import { useEffect,useRef,useState,type FormEvent } from "react";
import { publicationDetail,publicationsPage,preparedPublication,publicState,type PublicationDetail,type PublicationsPage,type PreparedPublication } from "@/platform/publishing/contract";
import { publishingMessages } from "@/platform/publishing/messages";
import styles from "./publishing-workspace.module.css";
type Props={scope:string;language:"es"|"en";locale:string;timeZone:string;canRequest:boolean;canApprove:boolean;canReconcile:boolean};
type Draft={id:string;message:string;when:string;expires:string;operation:"publish"|"revoke";original:string};
const localTime=()=>{const d=new Date();return new Date(d.getTime()-d.getTimezoneOffset()*60000).toISOString().slice(0,16)};
function Icon({kind}:{kind:"add"|"send"|"back"|"refresh"}){return <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.7" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">{kind==="add"?<path d="M12 5v14M5 12h14"/>:kind==="send"?<><path d="m21 3-7 18-4-7-7-4 18-7Z"/><path d="m10 14 6-6"/></>:kind==="back"?<path d="m14 6-6 6 6 6"/>:<><path d="M20 7v5h-5M4 17v-5h5"/><path d="M6.1 7a7 7 0 0 1 11.5-1L20 9M4 15l2.4 3A7 7 0 0 0 18 17"/></>}</svg>}
export function PublishingWorkspace(props:Props){return <Workspace key={props.scope} {...props}/>}
function Workspace({scope,language,locale,timeZone,canRequest,canApprove,canReconcile}:Props){
 const t=publishingMessages[language],storageKey=`elite.publishing.draft.${scope}`;
 const [page,setPage]=useState<PublicationsPage|null>(null),[selected,setSelected]=useState<string|null>(null),[detail,setDetail]=useState<PublicationDetail|null>(null);
 const [loading,setLoading]=useState(true),[detailLoading,setDetailLoading]=useState(false),[busy,setBusy]=useState(false),[error,setError]=useState(""),[notice,setNotice]=useState("");
 const [search,setSearch]=useState(""),[filter,setFilter]=useState("all"),[draft,setDraft]=useState<Draft|null>(null),[prepared,setPrepared]=useState<PreparedPublication|null>(null),[uncertain,setUncertain]=useState(false),[reason,setReason]=useState("");
 const heading=useRef<HTMLHeadingElement>(null),messageInput=useRef<HTMLTextAreaElement>(null),generation=useRef(0),alive=useRef(true),commandFence=useRef(false);
 const date=(v:string)=>new Intl.DateTimeFormat(locale,{dateStyle:"medium",timeStyle:"short",timeZone}).format(new Date(v));
 async function read(path:string){const response=await fetch(`/api/enterprise/publishing${path}`,{cache:"no-store"});if(!response.ok)throw new Error("unavailable");return response.json()}
 async function list(after=""){
  setLoading(true);setError("");try{const next=publicationsPage.parse(await read(after?`?after=${encodeURIComponent(after)}`:""));if(!alive.current)return;
   setPage(old=>after&&old&&old.organization_id===next.organization_id?{...next,items:[...old.items,...next.items]}:next);
  }catch{if(alive.current){setError(t.error);setPage(null);setDetail(null);setSelected(null);generation.current++}}finally{if(alive.current)setLoading(false)}
 }
 async function open(id:string){
  const current=++generation.current;setSelected(id);setDetail(null);setReason("");setDetailLoading(true);setError("");
  try{const data=publicationDetail.parse(await read(`?id=${encodeURIComponent(id)}`));if(alive.current&&current===generation.current){setDetail(data);if(prepared?.request.intent.approval_id===id){setPrepared(null);setDraft(null);setUncertain(false);try{sessionStorage.removeItem(storageKey)}catch{}}}}
  catch{if(alive.current&&current===generation.current)setError(t.error)}finally{if(alive.current&&current===generation.current){setDetailLoading(false);requestAnimationFrame(()=>heading.current?.focus())}}
 }
 useEffect(()=>{alive.current=true;void list();try{const raw=sessionStorage.getItem(storageKey);if(raw){const restored=preparedPublication.safeParse(JSON.parse(raw));if(restored.success){setPrepared(restored.data);setUncertain(true);setNotice(t.draftRecovered)}}}catch{}return()=>{alive.current=false;generation.current++}},[scope]);
 function retain(value:PreparedPublication){try{sessionStorage.setItem(storageKey,JSON.stringify(value))}catch{/* In-memory exact replay remains available; no false persistence claim. */}}
 function create(operation:"publish"|"revoke"="publish",original=""){
  if(uncertain){setNotice(t.cannotChange);return}setError("");setNotice("");setPrepared(null);setDraft({id:crypto.randomUUID(),message:"",when:localTime(),expires:"",operation,original});setSelected(null);setDetail(null);generation.current++;requestAnimationFrame(()=>messageInput.current?.focus());
 }
 async function post(command:unknown){const r=await fetch("/api/enterprise/publishing",{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify(command)});if(!r.ok)throw new Error("unconfirmed");return r.json()}
 async function command(fn:()=>Promise<void>){if(commandFence.current)return;commandFence.current=true;setBusy(true);setError("");setNotice("");try{await fn()}catch{if(alive.current)setError(t.commandError)}finally{commandFence.current=false;if(alive.current)setBusy(false)}}
 function preview(event:FormEvent){event.preventDefault();if(!draft||!canRequest)return;
  const scheduled=new Date(draft.when),expires=new Date(draft.expires);if(!Number.isFinite(+scheduled)||!Number.isFinite(+expires)||+expires<=+scheduled){setError(t.invalid);return}
  void command(async()=>{const result=preparedPublication.parse(await post({action:"prepare",draft:{approval_id:draft.id,operation:draft.operation,message:draft.operation==="publish"?draft.message:"",original_approval_id:draft.original,scheduled_at:scheduled.toISOString(),expires_at:expires.toISOString()}}));if(!alive.current)return;setPrepared(result);requestAnimationFrame(()=>heading.current?.focus())});
 }
 function submit(){if(!prepared||!canRequest)return;const exact=prepared;retain(exact);setUncertain(true);
  void command(async()=>{await post({action:"submit",prepared:exact});if(!alive.current)return;setNotice(t.saved);setPrepared(null);setDraft(null);setUncertain(false);try{sessionStorage.removeItem(storageKey)}catch{}await list();await open(exact.request.intent.approval_id)});
 }
 function decision(approve:boolean){if(!detail||!canApprove||detail.requested_by_self)return;const current=detail;
  void command(async()=>{await post({action:"decide",approval_id:current.status.approval_id,request_sha256:current.status.request_sha256,approve,reason});if(!alive.current)return;setNotice(t.decisionSaved);await list();await open(current.status.approval_id)});
 }
 function reconcile(){if(!detail||!canReconcile)return;const id=detail.status.approval_id;void command(async()=>{await post({action:"reconcile",approval_id:id});if(!alive.current)return;setNotice(t.checked);await open(id)})}
 function cancel(){if(uncertain)return;setDraft(null);setPrepared(null);setError("");try{sessionStorage.removeItem(storageKey)}catch{}}
 const visible=page?.items.filter(row=>(filter==="all"||row.approval_state===filter)&&row.message.toLocaleLowerCase(locale).includes(search.toLocaleLowerCase(locale)))??[];
 const composing=!!draft||!!prepared,hasDetail=!!selected||composing;
 const state=detail?publicState(detail):null;
 return <section className={styles.root} aria-labelledby="publishing-title" data-workspace="publishing">
  <div className={styles.top}><div><p className={styles.eyebrow}>{language==="es"?"Comunicación":"Communication"}</p><h1 id="publishing-title">{t.title}</h1></div>{canRequest&&page?.publish_enabled&&!composing&&!selected&&<button className={styles.primary} onClick={()=>create()} disabled={busy}><Icon kind="add"/>{t.new}</button>}</div>
  <p className={styles.intro}>{t.intro}</p>
  {notice&&<p role="status" className={styles.notice}>{notice}</p>}{error&&<p role="alert" className={styles.error}>{error}</p>}
  <div className={styles.columns} data-detail={hasDetail}>
   <section className={styles.list} aria-label={t.list}>
    <div className={styles.listTools}><label className={styles.search}><span className={styles.sr}>{t.search}</span><input type="search" placeholder={t.search} value={search} onChange={e=>setSearch(e.target.value)}/></label><button className={styles.iconButton} aria-label={t.refresh} disabled={loading||busy} onClick={()=>void list()}><Icon kind="refresh"/></button></div>
    <div className={styles.filters} role="group" aria-label={t.list}>{(["all","pending","approved","rejected"]as const).map(value=><button key={value} aria-pressed={filter===value} onClick={()=>setFilter(value)}>{t[value]}</button>)}</div>
    {loading&&!page?<p role="status" className={styles.quiet}>{t.loading}</p>:!visible.length?<p className={styles.quiet}>{page?.items.length?t.noMatches:t.empty}</p>:<ul className={styles.rows}>{visible.map(row=><li key={row.approval_id}><button className={styles.row} aria-current={selected===row.approval_id?"true":undefined} onClick={()=>{if(uncertain){setNotice(t.cannotChange);return}cancel();void open(row.approval_id)}}><span className={styles.rowIcon}><Icon kind="send"/></span><span className={styles.rowBody}><strong>{row.message}</strong><small>{row.operation==="revoke"?t.remove:t.publish} · {date(row.scheduled_at)}</small></span><span className={styles.badge}>{t[row.approval_state]}</span></button></li>)}</ul>}
    {page?.next_cursor&&<button className={styles.textButton} disabled={loading||busy} onClick={()=>void list(page.next_cursor)}>{t.more}</button>}
   </section>
   <section className={styles.detail} aria-label={t.detail}>
    {hasDetail&&<button className={`${styles.textButton} ${styles.back}`} disabled={busy} onClick={()=>{if(uncertain){setNotice(t.cannotChange);return}cancel();setSelected(null);setDetail(null);generation.current++}}><Icon kind="back"/>{t.back}</button>}
    {composing?prepared?<div className={styles.composer}><p className={styles.eyebrow}>{prepared.request.intent.operation==="revoke"?t.remove:t.publish}</p><h2 ref={heading} tabIndex={-1}>{t.preview}</h2><p className={styles.content}>{prepared.request.intent.message}</p><dl className={styles.dates}><div><dt>{prepared.request.intent.operation==="revoke"?t.removeWhen:t.when}</dt><dd>{date(prepared.request.scheduled_at)}</dd></div><div><dt>{t.expiry}</dt><dd>{date(prepared.request.expires_at)}</dd></div></dl>{uncertain&&<p role="status" className={styles.quiet}>{t.draftUncertain}</p>}<div className={styles.actions}>{canRequest&&<button className={styles.primary} disabled={busy} onClick={submit}>{busy?t.working:uncertain?t.retryExact:t.review}</button>}{uncertain?<button className={styles.secondary} disabled={busy} onClick={()=>void open(prepared.request.intent.approval_id)}>{t.reconcile}</button>:<button className={styles.textButton} disabled={busy} onClick={()=>setPrepared(null)}>{t.edit}</button>}</div></div>:draft&&<form className={styles.composer} onSubmit={preview}><h2 ref={heading} tabIndex={-1}>{draft.operation==="publish"?t.new:t.revoke}</h2>{draft.operation==="publish"?<label>{t.message}<textarea ref={messageInput} value={draft.message} onChange={e=>setDraft({...draft,message:e.target.value})} rows={6} required maxLength={16384} disabled={busy||!canRequest}/></label>:<p className={styles.quiet}>{t.revokeNote}</p>}<div className={styles.dateFields}><label>{draft.operation==="publish"?t.when:t.removeWhen}<input type="datetime-local" value={draft.when} onChange={e=>setDraft({...draft,when:e.target.value})} required disabled={busy||!canRequest}/></label><label>{t.expiry}<input type="datetime-local" value={draft.expires} onChange={e=>setDraft({...draft,expires:e.target.value})} required disabled={busy||!canRequest}/></label></div><p className={styles.quiet}>{t.timezone}</p><div className={styles.actions}><button className={styles.primary} disabled={busy||!canRequest} type="submit">{busy?t.working:t.prepare}</button><button type="button" className={styles.textButton} onClick={cancel} disabled={busy}>{t.cancel}</button></div></form>
    :detailLoading?<p role="status" className={styles.quiet}>{t.loading}</p>:detail&&state?<div className={styles.composer}><div className={styles.detailTitle}><h2 ref={heading} tabIndex={-1}>{detail.status.payload.intent.operation==="revoke"?t.remove:t.publish}</h2><span className={styles.badge} data-state={state}>{t[state]}</span></div><p className={styles.content}>{detail.status.payload.intent.message}</p><dl className={styles.dates}><div><dt>{detail.status.payload.intent.operation==="revoke"?t.removeWhen:t.when}</dt><dd>{date(detail.status.payload.scheduled_at)}</dd></div><div><dt>{t.expiry}</dt><dd>{date(detail.status.payload.expires_at)}</dd></div></dl>
     {state==="pending"?(detail.requested_by_self?<p className={styles.quiet}>{t.self}</p>:canApprove?<div className={styles.review}><label>{t.reason}<textarea value={reason} onChange={e=>setReason(e.target.value)} rows={2} maxLength={1000} disabled={busy}/></label><p className={styles.quiet}>{t.confirmApprove}</p><div className={styles.actions}><button className={styles.primary} disabled={busy} onClick={()=>decision(true)}>{t.approve}</button><button className={styles.secondary} disabled={busy} onClick={()=>decision(false)}>{t.reject}</button></div></div>:<p className={styles.quiet}>{t.waiting}</p>):null}
     <div className={styles.actions}>{state==="uncertain"&&canReconcile?<button className={styles.primary} onClick={reconcile} disabled={busy}>{t.reconcile}</button>:<button className={styles.textButton} onClick={()=>void open(detail.status.approval_id)} disabled={busy}><Icon kind="refresh"/>{t.refresh}</button>}{state==="published"&&canRequest&&page?.revoke_enabled&&<button className={styles.secondary} disabled={busy} onClick={()=>create("revoke",detail.status.approval_id)}>{t.revoke}</button>}</div><details className={styles.help}><summary>{language==="es"?"Cómo continúa":"What happens next"}</summary><p>{t.noAuto}</p><p>{t.expiryNote}</p></details>
    </div>:<div className={styles.placeholder}><Icon kind="send"/><h2>{t.choose}</h2></div>}
   </section>
  </div>
 </section>
}
