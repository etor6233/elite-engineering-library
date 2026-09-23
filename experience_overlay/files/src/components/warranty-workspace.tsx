"use client";
import { GuidedRecordSelect } from "@/components/guided-record-select";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

import {OperationalGuide} from "@/components/operational-guide";
import {RELEASE_GUIDES} from "@/platform/help/content";
import {useEffect,useRef,useState,type FormEvent} from "react";
import {z} from "zod";
import {catalogSHA} from "@/platform/catalog/authoring";
import {warrantyID,warrantySHA,warrantyCommand,warrantyOffer,warrantyActivation,warrantyClaim,warrantyStep,warrantyStepReference,warrantyStepMatches,pendingWarrantyStep,warrantyProfileView,warrantyQuote,type WarrantyCommand,type WarrantyClaim,type WarrantyOffer} from "@/platform/warranty/contract";
const endpoint="/api/enterprise/warranty";
const offerPending=z.object({family:z.literal("offer"),action:z.enum(["offer","acknowledge"]),quote_id:warrantyID,quote_version:z.string().regex(/^[1-9][0-9]{0,18}$/u),profile_sha256:warrantySHA,evidence_sha256:warrantySHA.optional()}).strict();
const activationPending=z.object({family:z.literal("activation"),handover_id:warrantyID}).strict();
const pendingSchema=z.discriminatedUnion("family",[pendingWarrantyStep,offerPending,activationPending]);
type Pending=z.infer<typeof pendingSchema>;
type PartInput={line_id:string;item_id:string;bin_id:string;lot_id:string;quantity:string;specific_receipt_entry:string};
const newPart=():PartInput=>({line_id:crypto.randomUUID(),item_id:"",bin_id:"",lot_id:"",quantity:"1",specific_receipt_entry:""});

export function WarrantyWorkspace({scope,organization,subject,tenant,permissions,surface,initialClaim,initialQuote,initialHandover}:{scope:string;organization:string;subject:string;tenant:string;permissions:string[];surface:"franchise"|"customer"|"factory";initialClaim:string;initialQuote:string;initialHandover:string}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();
const statusLabel={opened:t("p0603"),diagnosis:t("p0946"),repair:t("p0611"),quality:t("p0947"),closed:t("p0948"),cancelled:t("p0326")};

 const can=(p:string)=>permissions.includes("*")||permissions.includes(p);
 const held=useRef(false),saved=useRef<Pending|null>(null);const[ready,setReady]=useState(false),[busy,setBusy]=useState(false),[pending,setPending]=useState<Pending|null>(null),[notice,setNotice]=useState("");
 const[id,setID]=useState(initialClaim),[claim,setClaim]=useState<WarrantyClaim|null>(null),[quoteID,setQuoteID]=useState(initialQuote),[quote,setQuote]=useState<z.infer<typeof warrantyQuote>|null>(null),[offer,setOffer]=useState<WarrantyOffer|null>(null),[profile,setProfile]=useState<z.infer<typeof warrantyProfileView>|null>(null),[ack,setAck]=useState(false);
 const[handoverID,setHandoverID]=useState(initialHandover),[activation,setActivation]=useState<z.infer<typeof warrantyActivation>|null>(null),[appointment,setAppointment]=useState(""),[severity,setSeverity]=useState("medium"),[description,setDescription]=useState(""),[fault,setFault]=useState(""),[diagnosis,setDiagnosis]=useState(""),[labor,setLabor]=useState(""),[parts,setParts]=useState<PartInput[]>([]),[reason,setReason]=useState("");
 const[evidence,setEvidence]=useState(""),[correction,setCorrection]=useState("");const key=`elite-warranty:${scope}`;
 function persist(value:Pending|null){if(value){const raw=JSON.stringify(pendingSchema.parse(value));sessionStorage.setItem(key,raw);if(sessionStorage.getItem(key)!==raw)throw new Error("storage")}else sessionStorage.removeItem(key);saved.current=value;setPending(value)}
 useEffect(()=>{try{const raw=sessionStorage.getItem(key);if(raw){const p=pendingSchema.parse(JSON.parse(raw));saved.current=p;setPending(p);if(p.family==="claim")setID(p.case_id);if(p.family==="offer")setQuoteID(p.quote_id);if(p.family==="activation")setHandoverID(p.handover_id)}setReady(true)}catch{setNotice(t("p0835"))}},[key]);
 async function get(kind:string,id?:string,command?:string){const r=await fetch(endpoint+"?"+new URLSearchParams({kind,organization_id:organization,surface,...(id?{id}:{}),...(command?{command_id:command}:{})}),{cache:"no-store",signal:AbortSignal.timeout(10000)});if(!r.ok)throw new Error("unavailable");return r.json() as Promise<unknown>}
 async function readClaim(id:string){const v=warrantyClaim.parse(await get("claim",id));if(v.case_id!==id||(surface==="factory"?v.factory_organization_id:v.organization_id)!==organization||surface==="customer"&&v.customer_subject!==subject)throw new Error("scope");setClaim(v);setID(id)}
 async function consult(kind:"claim"|"quote"|"offer"|"profile"|"activation"){
  if(held.current)return;held.current=true;setBusy(true);if(kind==="claim")setClaim(null);if(kind==="quote")setQuote(null);if(kind==="offer"){setOffer(null);setAck(false)}
  try{
   if(kind==="claim")await readClaim(id);
   if(kind==="quote")setQuote(warrantyQuote.parse(await get(kind,quoteID)));
   if(kind==="offer")setOffer(warrantyOffer.parse(await get(kind,quoteID)));
   if(kind==="profile")setProfile(warrantyProfileView.parse(await get(kind)));
   if(kind==="activation")setActivation(warrantyActivation.parse(await get(kind,handoverID)));
   setNotice(t("p0949"));
  }catch{setNotice(t("p0950"))}finally{held.current=false;setBusy(false)}
 }
 async function reference(c:WarrantyCommand):Promise<Pending>{
  if(c.action==="offer"||c.action==="acknowledge")return offerPending.parse({family:"offer",action:c.action,quote_id:c.quote_id,quote_version:c.action==="offer"?(BigInt(c.quote_version)+1n).toString():c.quote_version,profile_sha256:c.profile_sha256,...(c.action==="acknowledge"?{evidence_sha256:c.evidence_sha256}:{})});
  if(c.action==="activate")return{family:"activation",handover_id:c.handover_id};return warrantyStepReference(c);
 }
 function acceptReceipt(raw:unknown,p:Pending){
  if(p.family==="claim"){const v=warrantyStep.parse(raw);if(!warrantyStepMatches(v,p,subject))throw new Error("identity");persist(null);setClaim(null);setID(v.case_id);return v.case_id}
  if(p.family==="offer"){const v=warrantyOffer.parse(raw);if(v.quote_id!==p.quote_id||v.quote_version!==p.quote_version||v.profile_sha256!==p.profile_sha256||v.profile.tenant_id!==tenant||v.profile.organization_id!==organization||(p.action==="offer"?v.offered_by!==subject:!v.acknowledged||v.customer_subject!==subject||v.evidence_sha256!==p.evidence_sha256))throw new Error("identity");persist(null);setOffer(v);setQuote(null);setAck(false);return null}
  const v=warrantyActivation.parse(raw);if(v.handover_id!==p.handover_id||v.organization_id!==organization||v.activated_by!==subject)throw new Error("identity");persist(null);setActivation(v);return null;
 }
 async function send(input:unknown){
  if(held.current||!ready||saved.current)return;held.current=true;setBusy(true);let started=false;
  try{const c=warrantyCommand.parse(input),p=await reference(c);persist(p);started=true;
   const r=await fetch(endpoint,{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify(c),signal:AbortSignal.timeout(10000)});
   if(!r.ok){if([400,401,403,413,415].includes(r.status)){persist(null);started=false}throw new Error("unconfirmed")}
   const caseID=acceptReceipt(await r.json(),p);started=false;setNotice(t("p0838"));if(caseID){try{await readClaim(caseID)}catch{setNotice(t("p0951"))}}
  }catch{setNotice(started?t("p0660"):t("p0718"))}finally{held.current=false;setBusy(false)}
 }
 async function recover(){const p=saved.current;if(!p||held.current)return;held.current=true;setBusy(true);try{const raw=p.family==="claim"?await get("command",p.case_id,p.command_id):p.family==="offer"?await get("offer",p.quote_id):await get("activation",p.handover_id);const caseID=acceptReceipt(raw,p);setNotice(t("p0841"));if(caseID){try{await readClaim(caseID)}catch{setNotice(t("p0952"))}}}catch{setNotice(t("p0843"))}finally{held.current=false;setBusy(false)}}
 async function choose(file:File|null,correct=false){const set=correct?setCorrection:setEvidence;set("");if(!file||file.size<1||file.size>1048576){setNotice(t("p0953"));return}set(await catalogSHA(new Uint8Array(await file.arrayBuffer())))}
 const blocked=busy||!ready||pending!==null;
 const base={organization_id:organization,surface};
 const command=(action:WarrantyCommand["action"],extra:Record<string,unknown>={})=>{if(claim)void send({...base,action,case_id:claim.case_id,command_id:crypto.randomUUID(),expected_version:claim.version,...extra})};
 function open(e:FormEvent){e.preventDefault();void send({...base,action:"open",case_id:crypto.randomUUID(),handover_id:handoverID,appointment_id:appointment,severity,description,evidence_sha256:evidence})}
 function planRepair(e:FormEvent){e.preventDefault();command("plan",{labor_work:labor,parts:parts.map(({lot_id,specific_receipt_entry,...p})=>({...p,...(lot_id?{lot_id}:{}),...(specific_receipt_entry?{specific_receipt_entry}:{})}))})}
 function partChange(i:number,k:keyof PartInput,v:string){setParts(old=>old.map((p,j)=>j===i?{...p,[k]:v}:p))}
 const ctx=claim?.context;
 return <>
 <p role="status" aria-live="polite">{notice}</p>
 <OperationalGuide className="card" guide={RELEASE_GUIDES["warranty-role-view"]}/>
 {pending&&<section className="card" aria-label={t("p0166")}><h2>{t("p0166")}</h2><p>{t("p0167")} <code>{pending.family==="claim"?pending.case_id:pending.family==="offer"?pending.quote_id:pending.handover_id}</code></p><button className="button" disabled={busy} onClick={()=>void recover()}>{t("p0168")}</button></section>}
 <section className="card"><h2>{t("p0954")}</h2><label>{t("p0848")}<input type="file" disabled={blocked} onChange={e=>void choose(e.target.files?.[0]??null)}/></label><p>{t("p0955")}</p>{evidence&&<p>{t("p0850")}</p>}<label>{t("p0220")}<textarea value={reason} onChange={e=>setReason(e.target.value)} maxLength={2048} disabled={blocked}/></label></section>
 {surface!=="factory"&&<section className="card"><h2>{t("p0956")}</h2><GuidedRecordSelect organization={organization} group="warrantyQuotes" name="reference:quoteID" label="Cotización" disabled={blocked} value={quoteID} onChange={nextValue=>{setQuoteID(nextValue);setQuote(null);setOffer(null);setAck(false)}} optional/>
 {surface==="franchise"&&can("warranty:offer")&&<><button className="button" disabled={busy||!quoteID} onClick={()=>void consult("quote")}>{t("p0958")}</button><button className="button" disabled={busy} onClick={()=>void consult("profile")}>{t("p0959")}</button>
 {quote&&<p>{t("p0960")} {quote.total_minor_units} {t("p0209")} {quote.currency}{t("p0008")} {quote.current?t("p0961"):t("p0962")}{t("p0008")}</p>}{profile&&<article><h3>{t("p0963")} {profile.profile.terms_version}</h3><p style={{whiteSpace:"pre-wrap"}}>{profile.profile.terms_text}</p><p>{t("p0964")} {profile.profile.parts_duration_days} {t("p0965")} {profile.profile.labor_duration_days} {t("p0966")}</p></article>}
 <button className="button" disabled={blocked||!quote||!profile||quote.state!=="issued"||!quote.current||!quote.customer_subject} onClick={()=>{if(quote&&profile)void send({...base,action:"offer",quote_id:quote.quote_id,quote_version:quote.quote_version,profile_sha256:profile.profile_sha256})}}>{t("p0967")}</button></>}
 <button className="button" disabled={busy||!quoteID} onClick={()=>void consult("offer")}>{t("p0968")}</button>
 {offer&&<article><h3>{t("p0969")} {offer.profile.terms_version}</h3><p style={{whiteSpace:"pre-wrap"}}>{offer.profile.terms_text}</p><p>{offer.acknowledged?t("p0970"):t("p0971")}</p>{surface==="customer"&&!offer.acknowledged&&<><label><input type="checkbox" checked={ack} onChange={e=>setAck(e.target.checked)} disabled={blocked}/>{t("p0972")}</label><button className="button" disabled={blocked||!ack||!evidence} onClick={()=>void send({...base,action:"acknowledge",quote_id:offer.quote_id,quote_version:offer.quote_version,profile_sha256:offer.profile_sha256,evidence_sha256:evidence})}>{t("p0973")}</button></>}</article>}
 </section>}
 {surface!=="factory"&&<section className="card"><h2>{t("p0974")}</h2><GuidedRecordSelect organization={organization} group="warrantyHandovers" name="reference:handoverID" label="Entrega" disabled={blocked} value={handoverID} onChange={nextValue=>{setHandoverID(nextValue);setActivation(null)}} optional/><button className="button" disabled={busy||!handoverID} onClick={()=>void consult("activation")}>{t("p0976")}</button>{surface==="franchise"&&can("warranty:activate")&&<button className="button" disabled={blocked||!handoverID} onClick={()=>void send({...base,action:"activate",handover_id:handoverID})}>{t("p0977")}</button>}
 {activation&&<article><h3>{t("p0978")}</h3><p>{t("p0167")} {activation.warranty_id}{t("p0979")} {activation.terms_version}{t("p0008")}</p><p>{t("p0980")} {activation.dates.parts_start} {t("p0748")} {activation.dates.parts_end}{t("p0981")} {activation.dates.labor_start} {t("p0748")} {activation.dates.labor_end}{t("p0008")}</p></article>}
 {(surface==="customer"||can("warranty:request"))&&<form onSubmit={open}><fieldset disabled={blocked}><legend>{t("p0982")}</legend><GuidedRecordSelect organization={organization} group="warrantyAppointments" name="reference:appointment" label="Cita de atención" disabled={blocked} value={appointment} onChange={nextValue=>setAppointment(nextValue)}/><label>{t("p0984")}<select value={severity} onChange={e=>setSeverity(e.target.value)}><option value="low">{t("p0612")}</option><option value="medium">{t("p0985")}</option><option value="high">{t("p0986")}</option><option value="safety">{t("p0987")}</option></select></label><label>{t("p0988")}<textarea value={description} onChange={e=>setDescription(e.target.value)} maxLength={4000} required/></label><button className="button" disabled={!evidence||!handoverID}>{t("p0989")}</button></fieldset></form>}
 </section>}
 <section className="card"><h2>{t("p0990")}</h2><form onSubmit={e=>{e.preventDefault();void consult("claim")}}><GuidedRecordSelect organization={organization} group="warrantyClaims" name="reference:id" label="Caso de garantía" disabled={blocked} value={id} onChange={nextValue=>{setID(nextValue);setClaim(null)}}/><button className="button" disabled={busy}>{t("p0990")}</button></form>
 {claim&&ctx&&<><p>{t("p0992")} <strong>{statusLabel[claim.state]}</strong>{t("p0674")} {claim.version}{t("p0008")}</p><p>{ctx.description}</p><p>{t("p0993")} {claim.parts_covered?t("p0994"):t("p0995")}{t("p0996")} {claim.labor_covered?t("p0994"):t("p0995")}{t("p0008")}</p>{ctx.diagnosis&&<article><h3>{t("p0997")}</h3><p>{ctx.diagnosis.fault_code}{t("p0213")} {ctx.diagnosis.description}</p><p>{ctx.diagnosis.excluded?t("p0998"):t("p0999")}</p></article>}
 {surface==="franchise"&&can("warranty:diagnose")&&claim.state==="opened"&&<form onSubmit={e=>{e.preventDefault();command("diagnose",{fault_code:fault,description:diagnosis,evidence_sha256:evidence})}}><fieldset disabled={blocked}><legend>{t("p1000")}</legend><GuidedRecordSelect organization={organization} group="warrantyFaults" name="reference:fault" label="Diagnóstico identificado" disabled={blocked} value={fault} onChange={nextValue=>setFault(nextValue)}/><label>{t("p1002")}<textarea value={diagnosis} onChange={e=>setDiagnosis(e.target.value)} maxLength={4000} required/></label><button className="button" disabled={!evidence}>{t("p1000")}</button></fieldset></form>}
 {surface==="franchise"&&can("warranty:plan")&&claim.state==="diagnosis"&&!ctx.plan&&<form onSubmit={planRepair}><fieldset disabled={blocked}><legend>{t("p1003")}</legend><label>{t("p1004")}<textarea value={labor} onChange={e=>setLabor(e.target.value)} maxLength={4000}/></label>{parts.map((p,i)=><fieldset key={p.line_id}><legend>{t("p1005")} {i+1}</legend>{([['item_id',t("p1006")],['bin_id',t("p1007")],['lot_id',t("p1008")],['quantity',t("p1009")],['specific_receipt_entry',t("p1010")]] as const).map(([k,label])=>k==="quantity"?<label key={k}>{label}<input value={p[k]} onChange={e=>partChange(i,k,e.target.value)} maxLength={6} required/></label>:<GuidedRecordSelect key={k} organization={organization} group={{item_id:"warrantyParts",bin_id:"warrantyBins",lot_id:"warrantyLots",specific_receipt_entry:"warrantyReceipts"}[k]} name={"part:"+i+":"+k} label={{item_id:"Repuesto",bin_id:"Depósito",lot_id:"Lote",specific_receipt_entry:"Recepción de origen"}[k]} disabled={blocked} value={p[k]} onChange={nextValue=>partChange(i,k,nextValue)} optional={k!=="item_id"&&k!=="bin_id"}/>)}<button type="button" onClick={()=>setParts(old=>old.filter((_,j)=>i!==j))}>{t("p1011")} {i+1}</button></fieldset>)}<button className="button" type="button" disabled={parts.length>=16} onClick={()=>setParts(old=>[...old,newPart()])}>{t("p1012")}</button><button className="button" disabled={parts.length===0&&!labor}>{t("p1003")}</button></fieldset></form>}
 {ctx.plan&&<article><h3>{t("p1013")}</h3><p>{t("p1014")} {ctx.plan.requester}{t("p1015")} {ctx.plan.approval_state==="pending"?t("p0797"):ctx.plan.approval_state==="approved"?t("p0798"):t("p0337")}{t("p0008")}</p><p>{ctx.plan.request.labor_work}</p><ul>{ctx.plan.request.parts.map(p=><li key={p.line_id}>{p.item_id}{t("p0213")} {p.quantity} {t("p0016")} {p.bin_id}{p.lot_id?` · lote ${p.lot_id}`:""}</li>)}</ul>
 {surface==="franchise"&&can("warranty:approve")&&claim.state==="diagnosis"&&ctx.plan.approval_state==="pending"&&<><p>{t("p1016")} {ctx.plan.expires_at}{t("p0008")}</p>{[true,false].map(approved=><button key={String(approved)} className="button" disabled={blocked||reason.trim().length<3||ctx.plan!.requester===subject||approved&&(ctx.plan!.excluded||!claim.parts_covered&&!claim.labor_covered)} onClick={()=>command("decide",{payload_sha256:ctx.plan!.payload_sha256,approved,reason})}>{approved?t("p1017"):t("p1018")}</button>)}</>}
 </article>}
 {surface==="franchise"&&can("warranty:work")&&claim.state==="repair"&&<button className="button" disabled={blocked||!evidence} onClick={()=>command("work",{evidence_sha256:evidence})}>{t("p1019")}</button>}
 {ctx.work&&<p>{t("p1020")} {ctx.work.actor}{t("p0008")}</p>}
 {surface==="franchise"&&can("warranty:quality")&&claim.state==="quality"&&ctx.work&&!ctx.acceptance&&!ctx.quality?.passed&&<article><h3>{t("p1021")}</h3>{ctx.quality&&!ctx.quality.passed&&<><p>{t("p1022")}</p><label>{t("p1023")}<input type="file" disabled={blocked} onChange={e=>void choose(e.target.files?.[0]??null,true)}/></label></>}{[true,false].map(passed=><button className="button" key={String(passed)} disabled={blocked||!evidence||ctx.work!.actor===subject||!!ctx.quality&&!correction} onClick={()=>command("quality",{passed,work_evidence_sha256:ctx.work!.payload_sha256,evidence_sha256:evidence,...(correction?{correction_evidence_sha256:correction}:{})})}>{passed?t("p1024"):t("p1025")}</button>)}</article>}
 {ctx.quality&&<p>{t("p1026")} {ctx.quality.passed?t("p0798"):t("p1027")}{t("p0008")}</p>}
 {surface==="customer"&&claim.state==="quality"&&ctx.quality?.passed&&!ctx.acceptance&&<button className="button" disabled={blocked||!evidence} onClick={()=>command("accept",{quality_sha256:ctx.quality!.payload_sha256,evidence_sha256:evidence})}>{t("p1028")}</button>}
 {ctx.acceptance&&<p>{t("p1029")}</p>}
 {surface==="factory"&&can("warranty:reconcile")&&claim.state==="quality"&&ctx.acceptance&&<button className="button" disabled={blocked||!evidence} onClick={()=>command("reconcile",{acceptance_sha256:ctx.acceptance!.payload_sha256,evidence_sha256:evidence})}>{t("p1030")}</button>}
 {ctx.reconciliation&&<article><h3>{t("p1031")}</h3><p>{t("p1032")} {ctx.reconciliation.recorded_inventory_cost}{t("p0008")}</p><p>{t("p1033")}</p></article>}
 {surface==="franchise"&&can("warranty:cancel")&&["opened","diagnosis","repair"].includes(claim.state)&&<button className="button" disabled={blocked||!evidence||reason.trim().length<3||ctx.plan?.approval_state==="pending"&&ctx.plan.requester===subject} onClick={()=>command("cancel",{reason,evidence_sha256:evidence})}>{t("p1034")}</button>}
 </>}
 </section>
 </>;
}
