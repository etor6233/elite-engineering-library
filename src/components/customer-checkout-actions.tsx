"use client";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

import { useState } from "react";
import { checkoutSchema } from "@/platform/payments/checkout";

export function CustomerCheckoutActions({orderId,organizationId}:{orderId:string;organizationId:string}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

 const [pending,setPending]=useState(false),[message,setMessage]=useState("");
 async function openCheckout(){
  setPending(true);setMessage("");
  try{
   const query=new URLSearchParams({orderId,organizationId});
   const response=await fetch("/api/enterprise/checkout?"+query,{cache:"no-store",redirect:"error",credentials:"same-origin",headers:{accept:"application/json"},signal:AbortSignal.timeout(7000)});
   if(response.status===404){setMessage(t("p0237"));return;}
   if(!response.ok)throw new Error("checkout unavailable");
   const value=checkoutSchema.parse(await response.json());
   if(value.order_id!==orderId)throw new Error("checkout does not match order");
   window.location.assign(value.url);
  }catch{setMessage(t("p0238"));}
  finally{setPending(false);}
 }
 return <div><button type="button" disabled={pending} onClick={openCheckout}>{pending?t("p0239"):t("p0240")}</button><p role="status" aria-live="polite">{message}</p></div>;
}
