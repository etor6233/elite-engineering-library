// AUTHORED bounded form/receipt correspondence; existing Go owners decide effects.
import {z} from "zod";
import {catalogCanonical,catalogSHA} from "@/platform/catalog/authoring";
export const warrantyID=z.string().regex(/^[A-Za-z0-9][A-Za-z0-9._:-]{0,119}$/u);
export const warrantySHA=z.string().regex(/^[0-9a-f]{64}$/u);
const positive=z.string().regex(/^[1-9][0-9]{0,18}$/u).refine(v=>BigInt(v)<=9223372036854775807n);
const text=(max:number,min=1)=>z.string().refine(v=>v.trim().length>=min&&new TextEncoder().encode(v).length<=max&&new TextDecoder().decode(new TextEncoder().encode(v))===v);
export const warrantySurface=z.enum(["franchise","customer","factory"]);
const base={organization_id:warrantyID,surface:warrantySurface};
const command={case_id:warrantyID,command_id:warrantyID,expected_version:positive};
const part=z.object({line_id:warrantyID,item_id:warrantyID,bin_id:warrantyID,lot_id:warrantyID.optional(),quantity:z.string().regex(/^[1-9][0-9]{0,5}$/u),specific_receipt_entry:warrantyID.optional()}).strict();
const plan=z.object({...command,labor_work:z.string().refine(v=>new TextEncoder().encode(v).length<=4000),parts:z.array(part).max(16).refine(v=>new Set(v.map(p=>p.line_id)).size===v.length)}).strict().refine(v=>v.parts.length>0||v.labor_work.length>0);
export const warrantyCommand=z.discriminatedUnion("action",[
 z.object({action:z.literal("offer"),...base,quote_id:warrantyID,quote_version:positive,profile_sha256:warrantySHA}).strict(),
 z.object({action:z.literal("acknowledge"),...base,quote_id:warrantyID,quote_version:positive,profile_sha256:warrantySHA,evidence_sha256:warrantySHA}).strict(),
 z.object({action:z.literal("activate"),...base,handover_id:warrantyID}).strict(),
 z.object({action:z.literal("open"),...base,case_id:warrantyID,handover_id:warrantyID,appointment_id:warrantyID,severity:z.enum(["low","medium","high","safety"]),description:text(4000,3),evidence_sha256:warrantySHA}).strict(),
 z.object({action:z.literal("diagnose"),...base,...command,fault_code:warrantyID,description:text(4000,3),evidence_sha256:warrantySHA}).strict(),
 z.object({action:z.literal("plan"),...base,...plan.shape}).strict().refine(v=>v.parts.length>0||v.labor_work.length>0),
 z.object({action:z.literal("decide"),...base,...command,payload_sha256:warrantySHA,approved:z.boolean(),reason:text(2048,3)}).strict(),
 z.object({action:z.literal("work"),...base,...command,evidence_sha256:warrantySHA}).strict(),
 z.object({action:z.literal("quality"),...base,...command,passed:z.boolean(),work_evidence_sha256:warrantySHA,evidence_sha256:warrantySHA,correction_evidence_sha256:warrantySHA.optional()}).strict(),
 z.object({action:z.literal("accept"),...base,...command,quality_sha256:warrantySHA,evidence_sha256:warrantySHA}).strict(),
 z.object({action:z.literal("reconcile"),...base,...command,acceptance_sha256:warrantySHA,evidence_sha256:warrantySHA}).strict(),
 z.object({action:z.literal("cancel"),...base,...command,reason:text(2048,3),evidence_sha256:warrantySHA}).strict()
]);
export type WarrantyCommand=z.infer<typeof warrantyCommand>;
export const warrantyProfile=z.object({schema:z.literal("elite-warranty-profile/v1"),scope:z.literal("MATERIALIZED_PROFILE"),algorithm:z.literal("bc-inclusive-fixed-terms"),algorithm_revision:z.literal(1),tenant_id:z.uuid(),organization_id:warrantyID,factory_organization_id:warrantyID,policy_id:warrantyID,terms_version:warrantyID,terms_text:text(16000),business_time_zone:z.string().min(1).max(128),parts_duration_days:z.number().int().min(1).max(3652500),labor_duration_days:z.number().int().min(1).max(3652500),work_reservation_seconds:z.number().int().min(60).max(86400),fault_exclusions:z.array(warrantyID).max(64),settlement:z.literal("INTERNAL_WARRANTY_SERVICE_ACKNOWLEDGEMENT"),authority_reference:warrantyID,decision_reference:warrantyID}).strict();
export const warrantyQuote=z.object({quote_id:warrantyID,organization_id:warrantyID,quote_version:positive,customer_subject:z.string().max(256),state:z.string().min(1).max(64),current:z.boolean(),currency:z.string().regex(/^[A-Z]{3}$/u),total_minor_units:z.string().regex(/^[0-9]{1,19}$/u)}).strict();
export const warrantyProfileView=z.object({profile_sha256:warrantySHA,profile:warrantyProfile}).strict();
export const warrantyOffer=z.object({quote_id:warrantyID,quote_version:positive,customer_subject:z.string().min(1).max(256),profile_sha256:warrantySHA,profile:warrantyProfile,offered_by:z.string().min(1).max(256),offered_at:z.iso.datetime({offset:true}),acknowledged:z.boolean(),evidence_sha256:warrantySHA.optional(),replay:z.boolean()}).strict();
export type WarrantyOffer=z.infer<typeof warrantyOffer>;
export const warrantyActivation=z.object({warranty_id:warrantyID,handover_id:warrantyID,quote_id:warrantyID,order_id:warrantyID,stock_unit_id:warrantyID,organization_id:warrantyID,customer_subject:z.string().min(1).max(256),profile_sha256:warrantySHA,terms_version:warrantyID,dates:z.object({parts_start:z.iso.date(),parts_end:z.iso.date(),labor_start:z.iso.date(),labor_end:z.iso.date()}).strict(),accepted_at:z.iso.datetime({offset:true}),activated_by:z.string().min(1).max(256),activated_at:z.iso.datetime({offset:true}),replay:z.boolean()}).strict();
const kinds=z.enum(["opened","diagnosed","planned","decision","work","quality","accepted","reconciled","cancelled"]);
const states=z.enum(["opened","diagnosis","repair","quality","closed","cancelled"]);
export const warrantyStep=z.object({case_id:warrantyID,command_id:warrantyID,version:positive,kind:kinds,state:states,actor:z.string().min(1).max(256),request_sha256:warrantySHA,payload_sha256:warrantySHA,payload:z.unknown(),recorded_at:z.iso.datetime({offset:true}),replay:z.boolean()}).strict();
export type WarrantyStep=z.infer<typeof warrantyStep>;
const roleContext=z.object({description:text(4000,3),severity:z.enum(["low","medium","high","safety"]),diagnosis:z.object({fault_code:warrantyID,description:text(4000,3),excluded:z.boolean(),evidence_sha256:warrantySHA}).strict().optional(),plan:z.object({request:plan,requester:z.string().min(1).max(256),payload_sha256:warrantySHA,excluded:z.boolean(),expires_at:z.iso.datetime({offset:true}),approval_state:z.enum(["pending","approved","rejected"])}).strict().optional(),work:z.object({actor:z.string().min(1).max(256),payload_sha256:warrantySHA,evidence_sha256:warrantySHA}).strict().optional(),quality:z.object({passed:z.boolean(),payload_sha256:warrantySHA,evidence_sha256:warrantySHA}).strict().optional(),acceptance:z.object({actor:z.string().min(1).max(256),payload_sha256:warrantySHA}).strict().optional(),reconciliation:z.object({settlement:z.literal("INTERNAL_WARRANTY_SERVICE_ACKNOWLEDGEMENT"),recorded_inventory_cost:z.string().regex(/^[0-9]+\.[0-9]{4}$/u)}).strict().optional()}).strict();
export const warrantyClaim=z.object({case_id:warrantyID,warranty_id:warrantyID,handover_id:warrantyID,appointment_id:warrantyID,organization_id:warrantyID,factory_organization_id:warrantyID,customer_subject:z.string().min(1).max(256),stock_unit_id:warrantyID,service_date:z.iso.date(),state:states,version:positive,profile_sha256:warrantySHA,parts_covered:z.boolean(),labor_covered:z.boolean(),latest:warrantyStep,context:roleContext}).strict();
export type WarrantyClaim=z.infer<typeof warrantyClaim>;
export const pendingWarrantyStep=z.object({family:z.literal("claim"),case_id:warrantyID,command_id:warrantyID,kind:kinds,version:positive,request_sha256:warrantySHA}).strict();
export type PendingWarrantyStep=z.infer<typeof pendingWarrantyStep>;
export function warrantyPayload(c:WarrantyCommand){const{action:_,organization_id:__,surface:___,...payload}=c;return payload}
export function warrantyPermission(c:WarrantyCommand){
 const action=c.action;
 if(action==="open")return c.surface==="customer"?"warranty:self":c.surface==="franchise"?"warranty:request":null;
 if(action==="acknowledge"||action==="accept")return c.surface==="customer"?"warranty:self":null;
 if(action==="reconcile")return c.surface==="factory"?"warranty:reconcile":null;
 if(c.surface!=="franchise")return null;
 return {offer:"warranty:offer",activate:"warranty:activate",diagnose:"warranty:diagnose",plan:"warranty:plan",decide:"warranty:approve",work:"warranty:work",quality:"warranty:quality",cancel:"warranty:cancel"}[action];
}
export const warrantyReadPermission=(surface:z.infer<typeof warrantySurface>)=>surface==="customer"?"warranty:self":surface==="factory"?"warranty:factory-read":"warranty:read";
export function warrantyPath(c:WarrantyCommand){const base=`/v1/${c.surface}/warranty/`;
 if(c.action==="offer"||c.action==="acknowledge")return base+`quotes/${encodeURIComponent(c.quote_id)}/${c.action==="offer"?"terms":"acknowledgements"}`;
 if(c.action==="activate")return base+`handovers/${encodeURIComponent(c.handover_id)}/activation`;
 if(c.action==="open")return base+"claims";
 return base+`claims/${encodeURIComponent(c.case_id)}/`+{diagnose:"diagnosis",plan:"plan",decide:"decision",work:"work",quality:"quality",accept:"acceptance",reconcile:"reconciliation",cancel:"cancellation"}[c.action];
}
export async function warrantyStepReference(c:WarrantyCommand):Promise<PendingWarrantyStep>{
 if(!("case_id"in c))throw new Error("claim command required");
 const kind={open:"opened",diagnose:"diagnosed",plan:"planned",decide:"decision",work:"work",quality:"quality",accept:"accepted",reconcile:"reconciled",cancel:"cancelled"}[c.action];
 return pendingWarrantyStep.parse({family:"claim",case_id:c.case_id,command_id:c.action==="open"?"open":c.command_id,kind,version:c.action==="open"?"1":(BigInt(c.expected_version)+1n).toString(),request_sha256:await catalogSHA(catalogCanonical(warrantyPayload(c)))});
}
export function warrantyStepMatches(v:WarrantyStep,p:PendingWarrantyStep,subject:string){return v.case_id===p.case_id&&v.command_id===p.command_id&&v.kind===p.kind&&v.version===p.version&&v.actor===subject&&v.request_sha256===p.request_sha256}
