// AUTHORED read contract. Money remains an exact decimal minor-unit string.
import {z} from "zod";
const id=z.string().min(1).max(128),sha=z.string().regex(/^[0-9a-f]{64}$/);
const owner=z.object({state:z.string().min(1).max(40),reference:z.string().max(255),amount_minor_units:z.string().regex(/^[1-9][0-9]{0,18}$/).refine(v=>BigInt(v)<=9223372036854775807n),currency:z.string().regex(/^[A-Z]{3}$/),handover_state:z.string().min(1).max(40).optional()}).strict();
export const returnOutcome=z.object({authorization_id:id,organization_id:id,order_id:id,disposition_id:id.optional(),remedy:z.enum(["refund","exchange"]).optional(),observed_at:z.string().datetime({offset:true}),stages:z.array(z.object({request_id:id,kind:z.enum(["inventory","refund","exchange","accounting","fiscal"]),status:z.enum(["requested","claimed","retry","blocked","failed","succeeded"]),attempts:z.number().int().min(0),updated_at:z.string().datetime({offset:true}),error_code:z.string().regex(/^[A-Z][A-Z0-9_]{1,63}$/).optional(),result_sha256:sha.optional(),owner:owner.optional()}).strict()).max(4)}).strict().superRefine((v,ctx)=>{
 const fail=()=>ctx.addIssue({code:"custom",message:"Incoherent return outcome"});
 if(!v.disposition_id){if(v.stages.length||v.remedy)fail();return}
 const kinds=v.remedy==="refund"?["inventory","refund","accounting","fiscal"]:v.remedy==="exchange"?["inventory","exchange","accounting"]:[];
 if(!kinds.length||v.stages.length!==kinds.length||new Set(v.stages.map(s=>s.kind)).size!==kinds.length||new Set(v.stages.map(s=>s.request_id)).size!==v.stages.length)fail();
 for(const s of v.stages){
  if(!kinds.includes(s.kind))fail();
  if(s.status==="succeeded"&&(!s.result_sha256||s.error_code))fail();
  if(["retry","blocked","failed"].includes(s.status)&&!s.error_code)fail();
  if(!s.owner){if(s.status==="succeeded"&&s.kind!=="inventory")fail();continue}
  const o=s.owner;const states:Record<string,string[]>={refund:["prepared","pending","requires_action","blocked","succeeded","failed","canceled"],exchange:["prepared"],accounting:["posted"],fiscal:["queued","authorizing","reconcile_required","authorized","rejected"]};
  if(!states[s.kind]?.includes(o.state))fail();
  if(!(s.kind==="refund"&&o.state==="prepared")&&!o.reference)fail();
  if(s.kind==="exchange"&&!o.handover_state)fail();
  if(s.status==="succeeded"&&o.state!==({refund:"succeeded",exchange:"prepared",accounting:"posted",fiscal:"authorized"}as Record<string,string>)[s.kind])fail();
 }
});
export type ReturnOutcome=z.infer<typeof returnOutcome>;
export function formatReturnMoney(minor:string,currency:string,locale:string){
 const amount=BigInt(minor),format=new Intl.NumberFormat(locale,{style:"currency",currency,currencyDisplay:"code"});const digits=format.resolvedOptions().maximumFractionDigits??2,scale=10n**BigInt(digits);
 return format.formatToParts(amount/scale).map(p=>p.type==="fraction"?(amount%scale).toString().padStart(digits,"0"):p.value).join("");
}
export function returnStageLabel(stage:ReturnOutcome["stages"][number],en=false){
 const successful:Record<ReturnOutcome["stages"][number]["kind"],[string,string]>={inventory:["Stock actualizado","Stock updated"],refund:["Reintegro confirmado","Refund confirmed"],exchange:["Cambio preparado","Replacement prepared"],accounting:["Asiento registrado","Entry posted"],fiscal:["Comprobante autorizado","Credit note authorized"]};
 if(stage.status==="succeeded")return successful[stage.kind][en?1:0];
 const status:Record<Exclude<ReturnOutcome["stages"][number]["status"],"succeeded">,[string,string]>={requested:["En espera","Waiting"],claimed:["Procesando","Processing"],retry:["Por conciliar","Awaiting reconciliation"],blocked:["Requiere revisión","Needs review"],failed:["No completado","Not completed"]};return status[stage.status][en?1:0];
}
