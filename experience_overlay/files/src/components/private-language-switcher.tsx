"use client";
import{useRef,useState,type FormEvent}from"react";
import{usePrivateI18n}from"@/platform/i18n/private-provider";
export function PrivateLanguageSwitcher({compact=false}:{compact?:boolean}){
 const{locale}=usePrivateI18n(),en=locale.language==="en";
 const[selected,setSelected]=useState(locale.language),[busy,setBusy]=useState(false),[failed,setFailed]=useState(false);const fence=useRef(false);
 async function save(event:FormEvent){event.preventDefault();if(fence.current||selected===locale.language)return;fence.current=true;setBusy(true);setFailed(false);
  try{const r=await fetch("/api/enterprise/locale",{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify({language:selected}),cache:"no-store",redirect:"error",signal:AbortSignal.timeout(7000)});if(!r.ok)throw new Error();window.location.reload()}
  catch{setFailed(true);setBusy(false);fence.current=false}
 }
 return <details><summary>{compact?(en?"Language":"Idioma"):(en?"Display language":"Idioma de la interfaz")}</summary><form onSubmit={save}>
 <label>{en?"Language":"Idioma"}<select value={selected} disabled={busy} onChange={e=>setSelected(e.target.value==="en"?"en":"es")}><option value="es" lang="es">Español</option><option value="en" lang="en">English</option></select></label>
 <p>{en?"Saving reloads this page. Finish or save any unsaved form first.":"Guardar recarga esta página. Terminá o guardá primero cualquier formulario sin guardar."}</p>
 <button className="button" disabled={busy||selected===locale.language}>{en?"Save language and reload":"Guardar idioma y recargar"}</button>
 {failed&&<p role="alert">{en?"Language save was not confirmed. No business operation was sent. Retry saving the preference or reload to check it.":"No se confirmó el guardado del idioma. No se envió una operación de negocio. Reintentá guardar la preferencia o recargá para consultarla."}</p>}
 </form></details>
}
