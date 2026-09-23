"use client";
import {useCallback,useEffect,useRef,useState} from "react";
import {OperationsList} from "@/components/admin-ops/operations-list";
import {accessSchema,catalogSchema,inboxSchema,type DocumentAccess,type DocumentCatalog,type DocumentInbox,type DocumentView,type PendingDocument} from "@/platform/documents-vnext/contract";
import {readMarkers} from "@/platform/documents-vnext/recovery";
import {DocumentDetail} from "./document-detail";
import {DocumentUpload} from "./document-upload";
import {DocumentIcon,DocumentNotice,documentName,documentState} from "./document-ui";
import s from "./documents.module.css";
const endpoint="/api/enterprise/documents-vnext";
export function DocumentWorkspace({initialDocumentId,language="es"}:{initialDocumentId?:string;language?:string}){
 const en=language==="en",tx=(es:string,enText:string)=>en?enText:es;
 const [access,setAccess]=useState<DocumentAccess|null>(null),[catalog,setCatalog]=useState<DocumentCatalog|null>(null),[inbox,setInbox]=useState<DocumentInbox|null>(null),[selected,setSelected]=useState<string|null>(initialDocumentId??null),[visited,setVisited]=useState<string[]>(initialDocumentId?[initialDocumentId]:[]),[upload,setUpload]=useState(false),[markers,setMarkers]=useState<PendingDocument[]>([]),[message,setMessage]=useState(""),[state,setState]=useState<"loading"|"ready"|"error"|"forbidden">("loading");
 const [cursor,setCursor]=useState<string|null>(null),[trail,setTrail]=useState<(string|null)[]>([]),[readBusy,setReadBusy]=useState(false);
 const accessRef=useRef<DocumentAccess|null>(null),listRef=useRef<HTMLDivElement>(null),listAbort=useRef<AbortController|null>(null);
 const markersChanged=useCallback(()=>{try{if(accessRef.current)setMarkers(readMarkers(accessRef.current.scope))}catch{setMessage("No se pudo recuperar una operación pendiente.")}},[]);
 async function loadPage(nextCursor:string|null){
  listAbort.current?.abort();const controller=new AbortController();listAbort.current=controller;setReadBusy(true);
  try{const response=await fetch(endpoint+"?view=inbox&limit=25"+(nextCursor?"&cursor="+encodeURIComponent(nextCursor):""),{cache:"no-store",redirect:"error",signal:AbortSignal.any([controller.signal,AbortSignal.timeout(10000)])});
   if([401,403].includes(response.status)){setState("forbidden");setInbox(null);throw new Error("FORBIDDEN")}
   if(!response.ok)throw new Error("INBOX");const page=inboxSchema.parse(await response.json());const a=accessRef.current;
   if(!a||page.items.some(view=>view.mode!==a.mode||view.profile_sha256!==a.profile_sha256))throw new Error("SCOPE");
   if(!controller.signal.aborted){setInbox(page);setCursor(nextCursor);setState("ready");setMessage("");return true}
  }catch(error){if(!controller.signal.aborted){if(!(error instanceof Error&&error.message==="FORBIDDEN"))setState("error");setMessage(tx("No pudimos cargar los documentos. Volvé a consultar.","Documents could not be loaded. Try again."))}return false}
  finally{if(!controller.signal.aborted)setReadBusy(false)}
 }
 useEffect(()=>{let active=true;const controller=new AbortController();void(async()=>{
  try{const response=await fetch(endpoint,{cache:"no-store",redirect:"error",signal:AbortSignal.any([controller.signal,AbortSignal.timeout(10000)])});if(!response.ok){if([401,403].includes(response.status))setState("forbidden");throw new Error("ACCESS")};const a=accessSchema.parse(await response.json());
   const classesResponse=await fetch(endpoint+"?view=classes",{cache:"no-store",redirect:"error",signal:AbortSignal.any([controller.signal,AbortSignal.timeout(10000)])});if(!classesResponse.ok)throw new Error("CATALOG");const c=catalogSchema.parse(await classesResponse.json());
   if(c.mode!==a.mode||c.profile_sha256!==a.profile_sha256)throw new Error("SCOPE");if(!active)return;
   accessRef.current=a;setAccess(a);setCatalog(c);setMarkers(readMarkers(a.scope));await loadPage(null);
  }catch{if(active){setState(previous=>previous==="forbidden"?previous:"error");setMessage(tx("Documentos no está disponible para esta sesión.","Documents is unavailable for this session."))}}
 })();return()=>{active=false;controller.abort();listAbort.current?.abort()}},[]);
 function choose(id:string){setSelected(id);setVisited(previous=>previous.includes(id)?previous:[...previous,id]);setUpload(false);const url=new URL(window.location.href);url.searchParams.set("document",id);window.history.replaceState(null,"",url.pathname+url.search)}
 function back(){const prior=selected;setSelected(null);const url=new URL(window.location.href);url.searchParams.delete("document");window.history.replaceState(null,"",url.pathname+url.search);requestAnimationFrame(()=>Array.from(listRef.current?.querySelectorAll<HTMLButtonElement>('button[data-record-id]')??[]).find(row=>row.dataset.recordId===prior)?.focus({preventScroll:true}))}
 function changed(view:DocumentView){setInbox(previous=>previous?{...previous,items:previous.items.map(item=>item.document_id===view.document_id?view:item)}:previous)}
 function received(view:DocumentView){changed(view);setUpload(false);choose(view.document_id);markersChanged();setTrail([]);void loadPage(null)}
 if(!access||!catalog||state==="forbidden")return <section className={s.root}><DocumentNotice error={Boolean(message)}>{message||tx("Cargando documentos…","Loading documents…")}</DocumentNotice>{message?<button type="button" className={s.secondary} onClick={()=>window.location.reload()}>{tx("Volver a abrir","Reopen")}</button>:null}</section>;
 return <section className={s.root} aria-label={tx("Documentos","Documents")} data-testid="document-workspace" data-state={state}>
  <header className={s.header}><h1>{tx("Documentos","Documents")}</h1><div className={s.actions}><button type="button" className={s.secondary} disabled={readBusy} onClick={()=>void loadPage(cursor)}>{readBusy?tx("Consultando…","Checking…"):tx("Actualizar","Refresh")}</button>{access.canWrite?<button type="button" className={selected?s.secondary:s.primary} onClick={()=>setUpload(true)}><DocumentIcon/>{tx("Nuevo documento","New document")}</button>:null}</div></header>
  {markers.length?<div className={s.recovery} role="status"><span>{tx("Hay operaciones por confirmar.","Some operations need confirmation.")}</span>{markers.map((marker,index)=><button type="button" key={marker.id} onClick={()=>choose(marker.id)}>{tx("Consultar operación","Check operation")}{markers.length>1?` ${index+1}`:""}</button>)}</div>:null}
  {message?<DocumentNotice error>{message}</DocumentNotice>:null}
  <div hidden={!upload}>{access.canWrite?<DocumentUpload access={access} classes={catalog.classes} language={language} onCancel={()=>setUpload(false)} onReceived={received} onMarkersChange={markersChanged}/>:null}</div>
  <div className={s.split} data-detail-open={Boolean(selected)} hidden={upload}>
   <div className={s.list} ref={listRef}>{inbox?<OperationsList items={inbox.items.map(view=>({id:view.document_id,title:documentName(view.name),subtitle:catalog.classes.find(c=>c.class_id===view.class_id)?.label_es??view.class_id,state:documentState(view.state,en)}))} selected={selected} onSelect={choose} language={language} label={tx("Documentos de esta página","Documents on this page")} testId="document-list"/>:<DocumentNotice error={state==="error"}>{state==="error"?tx("La lista no está disponible.","The list is unavailable."):tx("Cargando documentos…","Loading documents…")}</DocumentNotice>}
    <div className={s.pager}>{trail.length?<button className={s.secondary} disabled={readBusy} type="button" onClick={()=>{const previous=trail[trail.length-1]??null;void loadPage(previous).then(ok=>{if(ok)setTrail(value=>value.slice(0,-1))})}}>{tx("Anterior","Previous")}</button>:null}{inbox?.next_cursor?<button className={s.secondary} disabled={readBusy} type="button" onClick={()=>{const previous=cursor;void loadPage(inbox.next_cursor).then(ok=>{if(ok)setTrail(value=>[...value,previous])})}}>{tx("Siguiente","Next")}</button>:null}</div>
   </div>
   <div className={s.content}>{!selected?<div className={s.placeholder}><DocumentIcon/><p>{tx("Elegí un documento para revisarlo","Choose a document to review")}</p></div>:null}{visited.map(id=><div className={s.panel} key={`${access.scope}:${id}`} hidden={selected!==id}><DocumentDetail id={id} access={access} classes={catalog.classes} language={language} onBack={back} onChanged={changed} onMarkersChange={markersChanged}/></div>)}</div>
  </div>
 </section>
}
