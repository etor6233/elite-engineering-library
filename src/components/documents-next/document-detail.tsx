"use client";
import {useEffect,useRef,useState} from "react";
import {commandSchema,fieldsHash,reconciles,textHash,validFields,viewSchema,type DocumentAccess,type DocumentClass,type DocumentCommand,type DocumentFields,type DocumentView,type PendingDocument} from "@/platform/documents-vnext/contract";
import {clearMarker,readMarker,retainMarker} from "@/platform/documents-vnext/recovery";
import {DocumentPreview} from "./document-preview";
import {DocumentUpload} from "./document-upload";
import {DocumentFieldsForm,DocumentIcon,DocumentNotice,documentName,documentState} from "./document-ui";
import s from "./documents.module.css";
const endpoint="/api/enterprise/documents-vnext";
export function DocumentDetail({id,access,classes,language="es",onBack,onChanged,onMarkersChange}:{id:string;access:DocumentAccess;classes:DocumentClass[];language?:string;onBack:()=>void;onChanged:(view:DocumentView)=>void;onMarkersChange:()=>void}){
 const en=language==="en",tx=(es:string,enText:string)=>en?enText:es;
 const [view,setView]=useState<DocumentView|null>(null),[pending,setPending]=useState<PendingDocument|null>(null),[fields,setFields]=useState<DocumentFields>({}),[reason,setReason]=useState(""),[busy,setBusy]=useState(false),[message,setMessage]=useState(""),[readState,setReadState]=useState<"loading"|"ready"|"error"|"forbidden"|"missing">("loading");
 const generation=useRef(0);
 const activity=useRef(false),current=useRef<DocumentView|null>(null),edited=useRef(false),alive=useRef(true),abort=useRef(new AbortController()),pendingRef=useRef<PendingDocument|null>(null);
 const kind=classes.find(c=>c.class_id===view?.class_id&&c.schema_version===view?.schema_version);
 const signal=(ms:number)=>AbortSignal.any([abort.current.signal,AbortSignal.timeout(ms)]);
 function storePending(marker:PendingDocument|null){pendingRef.current=marker;setPending(marker)}
 async function accept(value:unknown,marker:PendingDocument|null){
  const next=viewSchema.parse(value);
  if(next.document_id!==id||next.profile_sha256!==access.profile_sha256||next.mode!==access.mode||!classes.some(c=>c.class_id===next.class_id&&c.schema_version===next.schema_version))throw new Error("View binding");
  if(!alive.current)return;
  const keep=edited.current&&current.current?.document_id===next.document_id&&next.state==="REVIEW_REQUIRED"&&!next.proposal;
  current.current=next;setView(next);setReadState("ready");
  if(!keep){setFields(next.proposal?.fields??next.suggested??{});edited.current=false}
  if(marker&&await reconciles(marker,next)){clearMarker(access.scope,id);storePending(null);setMessage("");onMarkersChange()}
  else if(marker)setMessage(tx("El resultado sigue pendiente. Consultá el estado o reintentá exactamente la misma operación.","The result is still pending. Check its status or retry the exact same operation."));
  onChanged(next);
 }
 async function refresh(known?:PendingDocument|null){
  if(activity.current)return;const requestGeneration=generation.current;activity.current=true;setBusy(true);
  try{const response=await fetch(`${endpoint}?id=${id}`,{cache:"no-store",redirect:"error",signal:signal(10000)});
   if(requestGeneration!==generation.current)return;
   if(response.status===401||response.status===403){setReadState("forbidden");setMessage(tx("No tenés acceso a este documento.","You do not have access to this document."));return}
   if(response.status===404){setReadState("missing");setMessage(tx("No encontramos un registro confirmado de este documento.","No confirmed record was found for this document."));return}
   if(!response.ok)throw new Error("Read failed");const value=await response.json();if(requestGeneration!==generation.current)return;await accept(value,known===undefined?pendingRef.current:known);
  }catch{if(alive.current&&requestGeneration===generation.current){setReadState("error");setMessage(tx("No pudimos consultar el documento. Tus datos no se descartan.","The document could not be loaded. Your changes are retained."))}}
  finally{if(requestGeneration===generation.current){activity.current=false;if(alive.current)setBusy(false)}}
 }
 useEffect(()=>{generation.current++;activity.current=false;alive.current=true;abort.current=new AbortController();try{const marker=readMarker(access.scope,id);storePending(marker);void refresh(marker)}catch{setReadState("error");setMessage(tx("No se pudo recuperar la operación. No se enviaron cambios.","The operation could not be recovered. No changes were sent."))}return()=>{generation.current++;activity.current=false;alive.current=false;abort.current.abort()}},[id,access.scope]);
 async function act(action:"process"|"review"|"decision",approved?:boolean){
  const record=current.current;if(activity.current||!record||!kind||readState!=="ready")return;
  const permission=action==="process"?access.canProcess:action==="review"?access.canWrite&&access.subject===record.uploader:access.canReview&&access.subject!==record.uploader;
  const sourceState=action==="process"?"QUARANTINED":action==="review"?"REVIEW_REQUIRED":"REVIEW_PENDING";
  if(!permission||record.state!==sourceState)return;
  activity.current=true;setBusy(true);setMessage("");
  try{
   if(action==="review"&&!validFields(fields,kind))throw new Error("FIELDS");
   const command:DocumentCommand=commandSchema.parse(action==="process"?{action,id}:action==="review"?{action,id,evidence_sha256:record.evidence_sha256,fields}:{action,id,payload_sha256:record.payload_sha256,approved,reason});
   const common={id,original_sha256:record.original_sha256,class_id:record.class_id,schema_version:record.schema_version};
   const marker:PendingDocument=action==="process"?{...common,action,expected:record.original_sha256}:action==="review"?{...common,action,expected:await fieldsHash(fields),evidence_sha256:record.evidence_sha256!}:{...common,action,expected:record.payload_sha256!,approved:Boolean(approved),reason_sha256:await textHash(reason)};
   retainMarker(access.scope,marker);storePending(marker);onMarkersChange();
   const response=await fetch(endpoint,{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify(command),cache:"no-store",redirect:"error",signal:signal(action==="process"?135000:15000)});
   const value=await response.json();
   if(!response.ok){if(value.effect==="NOT_ATTEMPTED"){clearMarker(access.scope,id);storePending(null);onMarkersChange();setMessage(tx("No se enviaron cambios. Revisá los campos y el acceso.","No changes were sent. Check the fields and your access."));return}throw new Error("UNCERTAIN")}
   await accept(value,marker);
  }catch{setMessage(pendingRef.current?tx("No se confirmó el resultado. Conservamos la operación para consultarla.","The result was not confirmed. The operation is retained for recovery."):tx("Revisá los campos antes de confirmar.","Check the fields before confirming."))}
  finally{activity.current=false;if(alive.current)setBusy(false)}
 }
 if(readState==="forbidden")return <article className={s.detail} data-testid="document-detail" data-read-state="forbidden"><button className={s.back} type="button" onClick={onBack}>{tx("Volver a documentos","Back to documents")}</button><DocumentNotice error>{message}</DocumentNotice><button className={s.secondary} type="button" disabled={busy} onClick={()=>void refresh()}>{tx("Volver a consultar","Try again")}</button></article>;
 if(readState==="missing"&&pending?.action==="receive")return <><DocumentNotice>{message}</DocumentNotice><DocumentUpload access={access} classes={classes} recovery={pending} language={language} onCancel={onBack} onMarkersChange={onMarkersChange} onReceived={next=>{void accept(next,pending)}}/><button className={s.secondary} type="button" disabled={busy} onClick={()=>void refresh()}>{tx("Consultar estado","Check status")}</button></>;
 return <article className={s.detail} data-testid="document-detail" data-read-state={readState} aria-busy={busy}>
  <button className={s.back} type="button" onClick={onBack} data-testid="document-back"><svg aria-hidden="true" viewBox="0 0 24 24"><path d="m14 5-7 7 7 7M7 12h14"/></svg>{tx("Volver a documentos","Back to documents")}</button>
  {!view?readState==="loading"?<DocumentNotice>{tx("Cargando documento…","Loading document…")}</DocumentNotice>:<><DocumentNotice error>{message}</DocumentNotice><button type="button" className={s.secondary} disabled={busy} onClick={()=>void refresh()}>{tx("Volver a consultar","Try again")}</button></>:<>
   <div className={s.detailHeader}><DocumentIcon/><div><h2>{documentName(view.name)}</h2><p>{kind?.label_es}</p></div></div>
   <DocumentPreview id={id} sha={view.original_sha256} name={view.name} language={language}/><span className={s.status} role="status">{documentState(view.state,en)}</span>
   <div className={`${s.actions} ${s.section}`}><a className={s.secondary} href={`${endpoint}?id=${id}&part=original`}>{tx("Descargar original","Download original")}</a><button className={s.secondary} type="button" disabled={busy} onClick={()=>void refresh()}>{busy?tx("Consultando…","Checking…"):tx("Consultar estado","Check status")}</button></div>
   {message?<DocumentNotice error={readState==="error"}>{message}</DocumentNotice>:null}
   {view.state==="QUARANTINED"&&access.canProcess?<div className={s.section}><button type="button" className={s.primary} disabled={busy||Boolean(pending&&pending.action!=="process")} onClick={()=>void act("process")}>{busy?tx("Leyendo documento…","Reading document…"):pending?tx("Reintentar lectura","Retry reading"):tx("Leer documento","Read document")}</button></div>:null}
   {view.state==="REVIEW_REQUIRED"&&kind&&access.canWrite&&access.subject===view.uploader?<form className={`${s.section} ${s.form}`} onSubmit={event=>{event.preventDefault();void act("review")}}><h3>{tx("Revisá los campos","Check the fields")}</h3><DocumentFieldsForm kind={kind} values={fields} disabled={busy||Boolean(pending&&pending.action!=="review")} onChange={(key,value)=>{edited.current=true;setFields(previous=>({...previous,[key]:value}))}}/><button className={s.primary} disabled={busy||Boolean(pending&&pending.action!=="review")}>{pending?tx("Reintentar la misma propuesta","Retry the same proposal"):tx("Enviar a revisión","Submit for review")}</button></form>:null}
   {view.proposal?<section className={s.section}><h3>{tx("Datos revisados","Reviewed fields")}</h3><dl className={s.summary}>{kind?.fields.map(field=><div key={field.key}><dt>{field.label_es}</dt><dd>{view.proposal!.fields[field.key]}</dd></div>)}</dl></section>:null}
   {view.state==="REVIEW_PENDING"?<section className={s.section}>{access.canReview&&access.subject!==view.uploader?<div className={s.decision}><h3>{tx("Decisión de revisión","Review decision")}</h3><label>{tx("Motivo","Reason")}<textarea value={reason} disabled={busy||Boolean(pending&&pending.action!=="decision")} maxLength={2048} onChange={event=>setReason(event.target.value)} required/></label><div className={s.actions}><button className={s.primary} type="button" disabled={busy||!reason.trim()||Boolean(pending&&(pending.action!=="decision"||pending.approved!==true))} onClick={()=>void act("decision",true)}>{pending?tx("Reintentar confirmación","Retry confirmation"):tx("Aprobar y confirmar","Approve and confirm")}</button><button className={s.secondary} type="button" disabled={busy||!reason.trim()||Boolean(pending&&(pending.action!=="decision"||pending.approved!==false))} onClick={()=>void act("decision",false)}>{tx("Rechazar","Reject")}</button></div></div>:<DocumentNotice>{tx("Otra persona debe revisar y confirmar este documento.","Another person must review and confirm this document.")}</DocumentNotice>}
   <button type="button" className={s.secondary} onClick={()=>{const url=new URL("/documents",window.location.origin);url.searchParams.set("document",id);void navigator.clipboard.writeText(url.href).then(()=>setMessage(tx("Enlace de revisión copiado.","Review link copied."))).catch(()=>setMessage(tx("No se pudo copiar el enlace.","The link could not be copied.")))}}>{tx("Copiar enlace de revisión","Copy review link")}</button></section>:null}
   {view.reviewer?<section className={s.section}><h3>{tx("Revisión registrada","Review recorded")}</h3><p>{view.reason}</p></section>:null}
   {view.evidence_sha256?<details className={s.receipt}><summary>{tx("Ver registro de revisión","View review record")}</summary>{["security","provider","analysis"].map(part=><div key={part}><a href={`${endpoint}?id=${id}&part=${part}`}>{part==="security"?tx("Comprobaciones del archivo","File checks"):part==="provider"?tx("Resultado de lectura","Reading result"):tx("Registro de extracción","Extraction record")}</a></div>)}</details>:null}
  </>}
 </article>
}
