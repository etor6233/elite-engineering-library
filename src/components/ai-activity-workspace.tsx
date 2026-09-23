"use client";
// AUTHORED UI glue: GET only. Response review/send stays in its existing owner.
import { useEffect,useRef,useState } from "react";
import { aiActivity,aiMessages,visibleActivity,type AIActivity,type ActivityState } from "@/platform/intelligence/activity";
import styles from "./ai-activity-workspace.module.css";
export function AIActivityWorkspace({scope,language="es",locale="es-AR",timeZone="UTC"}:{scope:string;language?:"es"|"en";locale?:string;timeZone?:string}){
 const t=aiMessages[language];const [stored,setStored]=useState<{scope:string;value:AIActivity}|null>(null);const [status,setStatus]=useState<"loading"|"ready"|"error"|"denied">("loading");const [filter,setFilter]=useState<ActivityState>("all");
 const generation=useRef(0),controller=useRef<AbortController|null>(null);
 const data=stored?.scope===scope?stored.value:null;
 async function refresh(){const seq=++generation.current;controller.current?.abort();const c=new AbortController();controller.current=c;setStatus("loading");setStored(null);
  try{const res=await fetch("/api/enterprise/intelligence/activity",{cache:"no-store",credentials:"same-origin",signal:c.signal});if(seq!==generation.current)return;
   if(res.status===401||res.status===403){setStatus("denied");return}if(!res.ok)throw Error("UNAVAILABLE");const value=aiActivity.parse(await res.json());if(seq!==generation.current)return;setStored({scope,value});setStatus("ready");
  }catch{if(seq===generation.current&&!c.signal.aborted){setStored(null);setStatus("error")}}
 }
 useEffect(()=>{setStored(null);setFilter("all");void refresh();return()=>{++generation.current;controller.current?.abort()}},[scope]);
 const rows=data?visibleActivity(data,filter):[];
 const date=(s:string)=>new Intl.DateTimeFormat(locale,{dateStyle:"medium",timeStyle:"short",timeZone}).format(new Date(s));
 return <section className={styles.workspace} aria-label={t.title}>
  <div className={styles.heading}><div><h1>{t.title}</h1><p>{t.subtitle}</p></div><button type="button" onClick={()=>void refresh()} disabled={status==="loading"}>{t.refresh}</button></div>
  <div className={styles.filters} role="group" aria-label={t.title}>{(["all","pending","handed_off"]as const).map(x=><button type="button" key={x} aria-pressed={filter===x} onClick={()=>setFilter(x)}>{x==="all"?t.all:x==="pending"?t.pending:t.handoff}</button>)}</div>
  <div aria-live="polite" aria-busy={status==="loading"}>
   {status==="loading"&&<p className={styles.message}>{t.loading}</p>}
   {(status==="error"||status==="denied")&&<p className={styles.message} role="alert">{status==="denied"?t.denied:t.error}</p>}
   {status==="ready"&&rows.length===0&&<p className={styles.message}>{data?.items.length?t.filteredEmpty:t.empty}</p>}
   {status==="ready"&&rows.map((item,index)=><article className={styles.card} key={item.request_id}>
    <div className={styles.row}><span className={styles.pictogram} aria-hidden="true"><svg viewBox="0 0 24 24" focusable="false">{item.state==="handed_off"?<path d="M6 18 18 6M7 6h11v11"/>:<path d="M5 5h14v11H9l-4 4V5Zm4 4h6m-6 3h4"/>}</svg></span><div><h2>{t.tool[item.tool]}</h2><time dateTime={item.occurred_at}>{date(item.occurred_at)}</time></div><span className={styles.state} data-tone={item.state==="handed_off"||item.review_state==="pending"?"pending":"quiet"}>{item.state==="handed_off"?t.handoff:item.review_state==="pending"?t.pending:item.review_state==="approved"?t.approved:t.rejected}</span></div>
    <details><summary>{t.details}<span className={styles.visuallyHidden}> {index+1}</span></summary><dl><dt>{t.reviewState}</dt><dd>{item.review_state==="pending"?t.pending:item.review_state==="approved"?t.approved:t.rejected}</dd><dt>{t.completed}</dt><dd>{item.state==="completed"?t.yes:t.handoff}</dd><dt>{t.attempts}</dt><dd>{item.attempts}</dd><dt>{t.reserved}</dt><dd>{BigInt(item.reserved_tokens).toLocaleString(locale)}</dd><dt>{t.input}</dt><dd>{BigInt(item.input_tokens).toLocaleString(locale)}</dd><dt>{t.output}</dt><dd>{BigInt(item.output_tokens).toLocaleString(locale)}</dd><dt>{t.intent}</dt><dd>{item.intent_pinned?t.yes:t.no}</dd></dl><p>{t.note}</p></details>
   </article>)}
  </div>
  {data?.has_more&&<p>{t.more}</p>}
  <a className={styles.review} href="/franchise/whatsapp">{t.review}<span aria-hidden="true"> →</span></a>
 </section>
}
