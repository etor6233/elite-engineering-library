"use client";
import { useId, useRef, useState, type FormEvent } from "react";
import { publicText as t } from "./messages";
import { PublicIcon } from "./public-ui";
import s from "./public-web.module.css";

type Submission = { modelId:string; name:string; email:string; consentGranted:true };
type SubmissionState = "idle"|"pending"|"error"|"uncertain"|"success";
// AUTHORED presentation adapter over the existing BFF. No new commercial rule.
// An uncertain write retains both key and exact payload. No automatic retry.
export function PublicLeadForm({modelId,locale}:{modelId:string;locale:string}) {
  const formId=useId();
  const request=useRef<{key:string;payload:Submission}|null>(null);
  const inFlight=useRef(false);
  const [state,setState]=useState<SubmissionState>("idle");
  const [leadId,setLeadId]=useState("");
  const [error,setError]=useState<"invalid"|"failed"|"unavailable">("failed");
  async function submit(event:FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if(inFlight.current||state==="success")return;
    const values=new FormData(event.currentTarget);
    if(state!=="uncertain") {
      request.current={key:crypto.randomUUID(),payload:{modelId,name:String(values.get("name")??""),email:String(values.get("email")??""),consentGranted:true}};
      if(values.get("consent")!=="on"){setError("invalid");setState("error");request.current=null;return;}
    }
    const attempt=request.current;
    if(!attempt)return;
    inFlight.current=true;setState("pending");
    try {
      const response=await fetch("/api/enterprise/leads",{method:"POST",headers:{"content-type":"application/json","idempotency-key":attempt.key},body:JSON.stringify(attempt.payload),signal:AbortSignal.timeout(12000)});
      if(!response.ok){
        // A gateway failure or conflict may follow a committed operation.
        if(response.status>=500||[408,409,429].includes(response.status)){setState("uncertain");return;}
        request.current=null;setError(response.status===400||response.status===422?"invalid":response.status===401||response.status===403?"unavailable":"failed");setState("error");return;
      }
      const body:unknown=await response.json();
      const receipt=body&&typeof body==="object"&&"lead_id"in body?body.lead_id:null;
      if(typeof receipt!=="string"||!/^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$/u.test(receipt)){setState("uncertain");return;}
      setLeadId(receipt);setState("success");
    }catch{setState("uncertain");}finally{inFlight.current=false;}
  }
  return <form className={s.leadForm} onSubmit={submit} aria-describedby={`${formId}-result`} data-public-state={state} aria-busy={state==="pending"}>
    {state!=="success"&&<fieldset disabled={state==="pending"||state==="uncertain"}>
      <label htmlFor={`${formId}-name`}>{t(locale,"name")}<input id={`${formId}-name`} name="name" required minLength={2} maxLength={120} autoComplete="name" /></label>
      <label htmlFor={`${formId}-email`}>{t(locale,"emailLabel")}<input id={`${formId}-email`} name="email" type="email" required maxLength={254} autoComplete="email" /></label>
      <label className={s.consent} htmlFor={`${formId}-consent`}><input id={`${formId}-consent`} name="consent" type="checkbox" required/><span>{t(locale,"consent")}</span></label>
    </fieldset>}
    <div id={`${formId}-result`} role="status" aria-live="polite" aria-atomic="true" className={state==="success"?s.success:state==="error"?s.error:state==="uncertain"?s.uncertain:s.result}>
      {state==="pending"&&<p>{t(locale,"sending")}</p>}
      {state==="error"&&<p>{t(locale,error)}</p>}
      {state==="uncertain"&&<><p><strong>{t(locale,"uncertain")}</strong></p><p>{t(locale,"uncertainBody")}</p></>}
      {state==="success"&&<><PublicIcon name="check"/><h3>{t(locale,"received")}</h3><p>{t(locale,"receivedBody")}</p></>}
    </div>
    {state==="success"?<a className={s.primary} href={`/locations?lead_id=${encodeURIComponent(leadId)}&model_id=${encodeURIComponent(modelId)}`}>{t(locale,"next")}<PublicIcon name="arrow"/></a>:<button type="submit" disabled={state==="pending"} className={s.primary}>{t(locale,state==="pending"?"sending":state==="uncertain"?"check":"send")}<PublicIcon name="arrow"/></button>}
  </form>;
}
