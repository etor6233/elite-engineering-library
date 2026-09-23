"use client";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

import { useState } from "react";
import type { ReplyView } from "@/platform/notifications/reply-contract";
type Row=ReplyView&{body:string;windowLabel:string};


export function WhatsAppReplyReview({replies,canSend}:{replies:Row[];canSend:boolean}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();
const sendLabels:Record<string,string>={not_started:t("p1035"),sending:t("p1036"),accepted:t("p1037"),unknown:t("p1038"),failed:t("p1039")};
const deliveryLabels:Record<string,string>={not_observed_by_this_reader:t("p1040"),observed_sent:t("p1041"),observed_delivered:t("p1042"),observed_read:t("p1043"),observed_failed:t("p1044"),ambiguous_latest_timestamp:t("p1045")};

 const [busy,setBusy]=useState<string|null>(null);const [notice,setNotice]=useState("");
 async function act(row:Row,action:"decision"|"send"|"recover",approved?:boolean){
  if(busy)return;setBusy(row.request_id);setNotice("");
  // Existing audit reasons are wire data, not translated display messages.
  const body={action,request_id:row.request_id,payload_sha256:row.payload_sha256,...(action==="decision"?{approved,reason:approved?"Texto y destinatario revisados en el portal":"Respuesta rechazada desde el portal"}:{})};
  try{const r=await fetch("/api/enterprise/franchise/whatsapp/replies",{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify(body)});if(!r.ok)throw new Error("UNCONFIRMED");window.location.reload()}
  catch{setNotice(t("p1048"))}
  // Stay disabled after uncertainty until the user reloads the authoritative state.
 }
 return <><p><a className="button" href="/franchise/whatsapp">{t("p1049")}</a></p><p role="status" aria-live="polite">{notice}</p>{replies.length===0?<p>{t("p1050")}</p>:replies.map(row=><article className="card" style={{marginTop:16,overflowWrap:"anywhere"}} key={row.request_id}>
  <h2>{t("p1051")} {row.context.message.ExternalID}</h2><p>{t("p0005")} {row.context.organization_id}</p><p>{t("p1052")} {row.windowLabel}</p><p style={{whiteSpace:"pre-wrap",overflowWrap:"anywhere"}}>{row.body}</p>
  <p>{t("p1053")} {row.state==="pending"?t("p0797"):row.state==="approved"?t("p0798"):t("p0337")}{t("p1054")} {row.status ? sendLabels[row.status.fence_state]??t("p1055") : t("p1035")}{t("p1056")} {row.status ? deliveryLabels[row.status.delivery_status]??t("p1055") : t("p1057")}{t("p0008")}</p>
  {row.state==="pending"&&<div style={{display:"flex",gap:12,flexWrap:"wrap"}}><button className="button" disabled={busy!==null} onClick={()=>void act(row,"decision",true)}>{t("p1058")}</button><button className="button" disabled={busy!==null} onClick={()=>void act(row,"decision",false)}>{t("p1059")}</button></div>}
  {row.state==="approved"&&canSend&&row.status?.fence_state==="not_started"&&!row.status.approval_expired&&<button className="button" disabled={busy!==null} onClick={()=>void act(row,"send")}>{t("p1060")}</button>}
  {row.state==="approved"&&canSend&&row.status?.reconciliation_required&&<><p>{t("p1061")}</p><button className="button" disabled={busy!==null} onClick={()=>void act(row,"recover")}>{t("p1062")}</button></>}
 </article>)}</>;
}
