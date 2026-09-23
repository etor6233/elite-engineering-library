"use client";
// AUTHORED. Authorized download is hashed before bytes enter an opaque sandbox.
import {useEffect,useRef,useState} from "react";
import {boundedDocumentBytes} from "@/platform/documents-vnext/preview-bytes";
import {sha256} from "@/platform/documents-vnext/contract";
import s from "./documents.module.css";
export function DocumentPreview({id,sha,name,language="es"}:{id:string;sha:string;name:string;language?:string}){
 const frame=useRef<HTMLIFrameElement>(null),[opened,setOpened]=useState(false),[nonce,setNonce]=useState(""),[state,setState]=useState("closed");
 const en=language==="en";
 useEffect(()=>{
  if(!opened||!nonce)return;const abort=new AbortController();let received=false;
  const timer=setTimeout(()=>{abort.abort();setState("rejected");setOpened(false)},20000);
  const listen=async(event:MessageEvent)=>{
   if(event.source!==frame.current?.contentWindow||event.origin!=="null"||event.data?.nonce!==nonce)return;
   if(event.data.kind==="ready"&&!received){received=true;try{
    const response=await fetch(`/api/enterprise/documents-vnext?id=${id}&part=original`,{cache:"no-store",redirect:"error",signal:abort.signal});
    if(!response.ok||response.headers.get("x-document-sha256")!==sha)throw Error("Original unavailable");
    const bytes=await boundedDocumentBytes(response,2097152);if(await sha256(bytes)!==sha||abort.signal.aborted)throw Error("Original binding");
    const buffer=new Uint8Array(bytes).buffer;frame.current?.contentWindow?.postMessage({kind:"document-bytes",nonce,bytes:buffer},"*",[buffer]);
   }catch{clearTimeout(timer);setState("rejected");setOpened(false)}}
   if(event.data.kind==="rendered"){clearTimeout(timer);setState("rendered")}
   if(event.data.kind==="rejected"){clearTimeout(timer);setState("rejected");setOpened(false)}
  };
  addEventListener("message",listen);return()=>{abort.abort();clearTimeout(timer);removeEventListener("message",listen)};
 },[opened,nonce,id,sha]);
 if(!/\.pdf$/iu.test(name))return null;
 return <section className={s.section} data-testid="document-preview" data-preview-state={state}>
  <button type="button" className={s.secondary} onClick={()=>{if(opened){setOpened(false);setState("closed")}else{setNonce(crypto.randomUUID());setOpened(true);setState("loading")}}}>{opened?(en?"Close preview":"Cerrar vista"):(en?"View original":"Ver original")}</button>
  {state==="loading"?<p role="status">{en?"Preparing preview…":"Preparando vista…"}</p>:null}
  {state==="rejected"?<p role="status">{en?"Preview unavailable. Download the original to review it.":"Vista no disponible. Descargá el original para revisarlo."}</p>:null}
  {opened?<iframe ref={frame} key={nonce} title={en?"Original document":"Documento original"} src={`/api/document-preview#${nonce}`} sandbox="allow-scripts" referrerPolicy="no-referrer" width="100%" height="560"/>:null}
 </section>;
}
