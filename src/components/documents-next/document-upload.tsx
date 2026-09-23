"use client";
import {useRef,useState} from "react";
import {MAX_DOCUMENT_BYTES,sha256,textHash,viewSchema,type DocumentAccess,type DocumentClass,type DocumentView,type PendingDocument} from "@/platform/documents-vnext/contract";
import {retainMarker,clearMarker} from "@/platform/documents-vnext/recovery";
import {DocumentIcon,DocumentNotice} from "./document-ui";
import s from "./documents.module.css";
const endpoint="/api/enterprise/documents-vnext";
export function DocumentUpload({access,classes,recovery,onReceived,onMarkersChange,onCancel,language="es"}:{access:DocumentAccess;classes:DocumentClass[];recovery?:PendingDocument;onReceived:(view:DocumentView)=>void;onMarkersChange:()=>void;onCancel:()=>void;language?:string}){
 const en=language==="en",tx=(es:string,enText:string)=>en?enText:es;
 const [classID,setClassID]=useState(recovery?.class_id??classes[0]?.class_id??""),[file,setFile]=useState<File|null>(null),[busy,setBusy]=useState(false),[message,setMessage]=useState("");
 const fence=useRef(false),retained=useRef<PendingDocument|null>(recovery??null);
 const kind=classes.find(value=>value.class_id===classID),locked=busy||!access.canWrite;
 async function upload(){
  if(fence.current||!file||!kind||!access.canWrite)return;fence.current=true;setBusy(true);setMessage("");
  try{
   const name=encodeURIComponent(file.name);
   if(!/^[\x21-\x7e]{1,128}$/u.test(name)||! /\.(pdf|jpe?g)$/iu.test(name)||file.size<1||file.size>MAX_DOCUMENT_BYTES)throw new Error("FILE");
   const bytes=new Uint8Array(await file.arrayBuffer()),digest=await sha256(bytes),metadata=await textHash(JSON.stringify([name,kind.class_id,kind.schema_version,access.profile_sha256]));
   const previous=retained.current;
   if(previous&&(previous.action!=="receive"||previous.expected!==digest||previous.metadata_sha256!==metadata||previous.class_id!==kind.class_id||previous.schema_version!==kind.schema_version)){setMessage(tx("Elegí el mismo archivo y tipo de la carga pendiente.","Choose the same file and type as the pending upload."));return}
   const marker:PendingDocument=previous??{id:crypto.randomUUID(),action:"receive",expected:digest,original_sha256:digest,class_id:kind.class_id,schema_version:kind.schema_version,metadata_sha256:metadata};
   retainMarker(access.scope,marker);retained.current=marker;onMarkersChange();
   const response=await fetch(`${endpoint}?id=${marker.id}`,{method:"PUT",headers:{"content-type":"application/octet-stream","x-document-name":name,"x-document-sha256":digest,"x-document-class":kind.class_id,"x-document-schema-version":kind.schema_version},body:new Uint8Array(bytes).buffer,cache:"no-store",redirect:"error",signal:AbortSignal.timeout(15000)});
   const value=await response.json();
   if(!response.ok){if(value.effect==="NOT_ATTEMPTED"){clearMarker(access.scope,marker.id);retained.current=null;onMarkersChange();setMessage(tx("No se cargó el archivo. Revisá el tipo y el tamaño.","The file was not uploaded. Check its type and size."));return}throw new Error("UNCERTAIN")}
   const view=viewSchema.parse(value);
   if(view.document_id!==marker.id||view.original_sha256!==digest||view.profile_sha256!==access.profile_sha256||view.mode!==access.mode||view.class_id!==kind.class_id||view.schema_version!==kind.schema_version)throw new Error("BINDING");
   clearMarker(access.scope,marker.id);retained.current=null;onMarkersChange();setFile(null);onReceived(view);
  }catch{setMessage(retained.current?tx("No se confirmó la carga. Consultá su estado antes de reintentar.","Upload was not confirmed. Check its status before retrying."):tx("No pudimos leer el archivo. Usá un PDF o JPEG de hasta 2 MiB.","The file could not be read. Use a PDF or JPEG up to 2 MiB."))}
  finally{fence.current=false;setBusy(false)}
 }
 return <section className={s.upload} aria-label={tx("Cargar documento","Upload document")} data-testid="document-upload">
  <div className={s.actions}><button type="button" className={s.back} disabled={busy} onClick={onCancel}>← {tx("Volver","Back")}</button></div>
  <h2>{recovery?tx("Recuperar carga","Recover upload"):tx("Nuevo documento","New document")}</h2>
  <label>{tx("Tipo de documento","Document type")}<select aria-label={tx("Tipo de documento","Document type")} value={classID} disabled={locked||Boolean(recovery)||Boolean(retained.current)} onChange={event=>setClassID(event.target.value)}>{classes.map(c=><option key={c.class_id} value={c.class_id}>{c.label_es}</option>)}</select></label>
  <label className={s.file}><DocumentIcon/><strong>{tx("Seleccioná el archivo","Choose the file")}</strong><input type="file" accept="application/pdf,image/jpeg,.pdf,.jpg,.jpeg" disabled={locked} onChange={event=>{setFile(event.target.files?.[0]??null);setMessage("")}}/><p>{file?`${file.name} · ${Math.ceil(file.size/1024)} KiB`:tx("PDF o JPEG · hasta 2 MiB","PDF or JPEG · up to 2 MiB")}</p></label>
  {message?<DocumentNotice error>{message}</DocumentNotice>:null}
  <button className={s.primary} type="button" disabled={locked||!file||!kind||Boolean(retained.current&&!recovery)} onClick={()=>void upload()}>{busy?tx("Cargando…","Uploading…"):recovery?tx("Reintentar la misma carga","Retry the same upload"):tx("Cargar documento","Upload document")}</button>
 </section>
}
