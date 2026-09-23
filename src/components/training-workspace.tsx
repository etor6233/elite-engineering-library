"use client";
import {guideDisplay} from "@/platform/i18n/private-guide-display";
import {trainingCourseDisplay} from "@/platform/i18n/private-training-display";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

import {OperationalGuide} from "@/components/operational-guide";
import {RELEASE_GUIDES} from "@/platform/help/content";
import { useEffect,useRef,useState,type FormEvent } from "react";
import { assessmentSchema,attemptSchema,trainingReference,type CourseView,type Assessment,type Attempt,type TrainingReference } from "@/platform/training/contract";


const endpoint="/api/enterprise/training";
export function TrainingWorkspace({courses,assessments,canLearn,canReview,scope}:{courses:CourseView[];assessments:Assessment[];canLearn:boolean;canReview:boolean;scope:string}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();
const labels:Record<string,string>={owner:t("p0903"),admin:t("p0904"),employee:t("p0492"),customer:t("p0125")};

 const [selected,setSelected]=useState(courses[0]?.course.id??""),[reference,setReference]=useState<TrainingReference|null>(null),[attempt,setAttempt]=useState<Attempt|null>(null);
 const [ready,setReady]=useState(false),[busy,setBusy]=useState(false),[reading,setReading]=useState(false),[notice,setNotice]=useState("");
 const fence=useRef(false),posting=useRef(false),readingNow=useRef(false);const[transportPending,setTransportPending]=useState(false);const key=`elite-training:${scope}`;
 useEffect(()=>{try{const raw=sessionStorage.getItem(key);if(raw){const saved=trainingReference.parse(JSON.parse(raw));setReference(saved);setSelected(saved.course_id);setNotice(t("p0908"))};setReady(true)}catch{fence.current=true;setNotice(t("p0909"))}},[key]);
 async function consult(ref=reference){if(!ref||readingNow.current||posting.current)return;readingNow.current=true;setReading(true);try{const r=await fetch(`${endpoint}?attempt_id=${encodeURIComponent(ref.attempt_id)}`,{cache:"no-store"});if(!r.ok)throw new Error();const value=attemptSchema.parse(await r.json());if(value.attempt_id!==ref.attempt_id||value.content.profile_sha256!==ref.profile_sha256)throw new Error();setAttempt(value);fence.current=false;setBusy(false);setNotice(t("p0910"))}catch{setNotice(t("p0911"))}finally{readingNow.current=false;setReading(false)}}
 async function act(body:Record<string,unknown>){if(fence.current||readingNow.current)return;fence.current=true;posting.current=true;setTransportPending(true);setBusy(true);setNotice("");try{const r=await fetch(endpoint,{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify(body)});if(!r.ok)throw new Error();const value:unknown=await r.json();posting.current=false;setTransportPending(false);if(body.action==="assess"){assessmentSchema.parse(value);window.location.reload();return};if(body.action==="submit"){assessmentSchema.parse(value);await consult();return};const current=attemptSchema.parse(value);setAttempt(current);fence.current=false;setBusy(false);setNotice(t("p0912"))}catch{setNotice(t("p0913"))}finally{posting.current=false;setTransportPending(false)}}
 function start(existing?:TrainingReference){const view=courses.find(c=>c.course.id===selected);if(!view||!ready)return;const ref=existing??{attempt_id:crypto.randomUUID(),course_id:view.course.id,profile_sha256:view.profile_sha256};try{sessionStorage.setItem(key,JSON.stringify(ref));setReference(ref)}catch{setNotice(t("p0914"));return};void act({action:"start",...ref})}
 function nextAttempt(){if(!attempt||(!attempt.assessment&&attempt.current_profile))return;try{sessionStorage.removeItem(key);setReference(null);setAttempt(null);fence.current=false;setBusy(false);setNotice(t("p0915"))}catch{setNotice(t("p0916"))}}
 function submit(event:FormEvent<HTMLFormElement>){event.preventDefault();if(!attempt)return;const form=new FormData(event.currentTarget),answers:Record<string,string>={};for(const prompt of attempt.content.course.prompts){const value=String(form.get(prompt.id)??"");if(!value.trim()||new TextEncoder().encode(value).length>2048){setNotice(t("p0917"));return};answers[prompt.id]=value};void act({action:"submit",attempt_id:attempt.attempt_id,profile_sha256:attempt.content.profile_sha256,answers})}
 const course=attempt?.content??courses.find(c=>c.course.id===selected);
 return <>
 <OperationalGuide className="card" guide={RELEASE_GUIDES["training-role-view"]}/>
  <p role="status" aria-live="polite">{notice}</p>
  {canLearn&&<section className="card" aria-label={t("p0918")} style={{overflowWrap:"anywhere"}}><h2>{t("p0918")}</h2>
   {!reference&&<><label>{t("p0919")}<select value={selected} disabled={!ready||busy} onChange={event=>setSelected(event.target.value)}>{courses.map(c=><option key={c.course.id} value={c.course.id}>{labels[c.course.role]} {t("p0016")} {trainingCourseDisplay(c,privateLocale.language).title}</option>)}</select></label><p>{t("p0920")}</p><button className="button" disabled={!ready||busy} onClick={()=>start()}>{t("p0921")}</button></>}
   {reference&&<p><button className="button" disabled={reading||transportPending} onClick={()=>void consult()}>{t("p0922")}</button>{!attempt&&!busy&&<button className="button" onClick={()=>start(reference)}>{t("p0923")}</button>}</p>}
   {course&&<><h3>{trainingCourseDisplay(course,privateLocale.language).title}</h3><p>{t("p0924")} {course.profile_id}{t("p0925")} {course.profile_revision} {t("p0926")} {labels[course.course.role]}</p></>}
   {attempt&&!attempt.current_profile&&<p role="alert">{t("p0927")}</p>}
   {attempt&&course?.articles.map(source=>{const article=guideDisplay(source,privateLocale.language);return <section lang={article.displayLanguage} key={article.id} aria-label={article.title}><h3>{article.title} {t("p0543")} {article.version}</h3>{article.paragraphs.map((text,index)=><p key={index}>{text}</p>)}{attempt.read_lessons.includes(article.id)?<p>{t("p0928")}</p>:<button className="button" disabled={busy||!attempt.current_profile} onClick={()=>void act({action:"acknowledge",attempt_id:attempt.attempt_id,lesson_id:article.id,profile_sha256:attempt.content.profile_sha256})}>{t("p0929")}</button>}</section>})}
   {attempt&&!attempt.assessment&&attempt.current_profile&&attempt.read_lessons.length===attempt.content.course.lessons.length&&<form onSubmit={submit}><h3>{t("p0930")}</h3>{trainingCourseDisplay(attempt.content,privateLocale.language).prompts.map(p=><label key={p.id} style={{display:"block",marginBottom:16}}>{p.text}<textarea name={p.id} required maxLength={2048} disabled={busy} style={{display:"block",width:"100%",minHeight:100}}/></label>)}<button className="button" disabled={busy}>{t("p0931")}</button></form>}
   {attempt?.assessment&&<AssessmentResult value={attempt.assessment}/>}
   {attempt&&(attempt.assessment?.state==="approved"||attempt.assessment?.state==="rejected"||!attempt.current_profile)&&<p><button className="button" disabled={busy} onClick={nextAttempt}>{t("p0932")}</button></p>}
  </section>}
  <section aria-label={t("p0933")}><h2>{canReview?t("p0934"):t("p0935")}</h2><p><a className="button" href="/guide/training">{t("p0936")}</a></p>
   {assessments.length===0?<p>{t("p0937")}</p>:assessments.map(value=><article className="card" key={value.request_id} style={{marginTop:16,overflowWrap:"anywhere"}}><h3>{trainingCourseDisplay(value.payload.content,privateLocale.language).title}</h3><p>{t("p0938")} {value.payload.learner_subject}</p><p>{t("p0939")} {value.payload.content.profile_revision}</p>
    {value.payload.content.articles.map(source=>{const article=guideDisplay(source,privateLocale.language);return <details lang={article.displayLanguage} key={article.id}><summary>{article.title} {t("p0543")} {article.version}</summary>{article.paragraphs.map((text,i)=><p key={i}>{text}</p>)}</details>})}
    <AssessmentResult value={value}/>
    {canReview&&value.state==="pending"&&<form onSubmit={event=>{event.preventDefault();const form=new FormData(event.currentTarget);void act({action:"assess",request_id:value.request_id,payload_sha256:value.payload_sha256,approved:form.get("decision")==="approve",reason:String(form.get("reason")??"")})}}><label>{t("p0940")}<textarea name="reason" required maxLength={2048} disabled={busy} style={{display:"block",width:"100%",minHeight:90}}/></label><label>{t("p0941")}<select name="decision" disabled={busy}><option value="reject">{t("p0906")}</option><option value="approve">{t("p0905")}</option></select></label><p><button className="button" disabled={busy}>{t("p0942")}</button></p></form>}
   </article>)}
  </section>
 </>;
}
function AssessmentResult({value}:{value:Assessment}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();
const result=(state:string)=>state==="approved"?t("p0905"):state==="rejected"?t("p0906"):t("p0907");
return <section aria-label={t("p0943")}><p><strong>{result(value.state)}</strong></p>{trainingCourseDisplay(value.payload.content,privateLocale.language).prompts.map(p=><div key={p.id}><p>{p.text}</p><p lang="und" style={{whiteSpace:"pre-wrap",overflowWrap:"anywhere"}}>{value.payload.answers[p.id]}</p></div>)}{value.reviewer&&<p>{t("p0944")} {value.reviewer}{t("p0008")} {value.reason}</p>}<p>{t("p0945")}</p></section>}
