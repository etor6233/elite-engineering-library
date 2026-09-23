"use client";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

// AUTHORED UI glue; an observed state never claims attribution to a lost request.
import{useEffect,useRef,useState}from"react";
import{factoryUnitSchema,type FactoryUnit}from"@/platform/factory/contracts";
export function FactoryUnitActions({unit,canWrite,storageScope}:{unit:FactoryUnit;canWrite:boolean;storageScope:string}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

 const[value,setValue]=useState(unit),[pending,setPending]=useState(false),[message,setMessage]=useState(""),[busy,setBusy]=useState(false);const active=useRef(false);
 const key=`elite:factory:${storageScope}:${unit.id}`;
 useEffect(()=>{try{setPending(localStorage.getItem(key)==="unknown")}catch{setPending(true);setMessage(t("p0299"))}},[key]);
 async function act(write:boolean){if(active.current)return;active.current=true;setBusy(true);setMessage("");
 try{
  if(write){localStorage.setItem(key,"unknown");setPending(true)}
  const q=new URLSearchParams({organizationId:unit.organization_id,unitId:unit.id});
  const response=await fetch(`/api/enterprise/factory${write?"":"?"+q}`,write?{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify({organizationId:unit.organization_id,unitId:unit.id,action:"start-assembly",expectedState:"planned"})}:{cache:"no-store"});
  if(!response.ok)throw new Error("unconfirmed");const body:unknown=await response.json();
  if(write){if(!(typeof body==="object"&&body!==null&&"status"in body&&body.status==="accepted"))throw new Error("invalid response");setMessage(t("p0300"))}
  else{const parsed=factoryUnitSchema.parse((body as {unit?:unknown}).unit);if(parsed.id!==unit.id||parsed.organization_id!==unit.organization_id)throw new Error("mismatch");setValue(parsed);localStorage.removeItem(key);setPending(false);setMessage(t("p0301"))}
 }catch{setMessage(write?t("p0302"):t("p0303"))}finally{active.current=false;setBusy(false)}}
 return <section aria-label={`Seguimiento de ${unit.serial_number}`}><p>{t("p0304")} <strong>{value.state}</strong></p>{pending?<p>{t("p0305")}</p>:null}
 {canWrite&&value.state==="planned"?<button className="button" type="button" disabled={busy||pending} onClick={()=>void act(true)}>{t("p0306")}</button>:null}
 <button className="button" type="button" disabled={busy} onClick={()=>void act(false)}>{t("p0307")}</button>
 {message?<p role="status" aria-live="polite">{message}</p>:null}</section>
}
