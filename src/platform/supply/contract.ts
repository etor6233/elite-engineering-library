// AUTHORED transport/form binding to the existing Go serial supply owner.
import {z} from "zod";
export const supplyID=z.string().regex(/^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$/u);
export const supplySHA=z.string().regex(/^[0-9a-f]{64}$/u);
const unsigned=z.string().regex(/^(0|[1-9][0-9]{0,18})$/u).refine(v=>BigInt(v)<=9223372036854775807n);
const positive=unsigned.refine(v=>v!=="0");
const text=(max:number)=>z.string().min(1).refine(v=>new TextEncoder().encode(v).length<=max&&v.trim()===v&&!/[\p{Cc}]/u.test(v)&&new TextDecoder().decode(new TextEncoder().encode(v))===v);
export const supplyPolicy="strict-serial-reference/v1";
export const supplyKinds=["submit","confirm","start","cancel","register","milestone","ship","receive","quality","quality-reject","reinspect"] as const;
const line=z.object({id:supplyID,variant_id:supplyID,quantity:z.number().int().min(1).max(1000)}).strict();
const planFields={factory_organization_id:supplyID,demand_reference:text(1024),policy_code:z.literal(supplyPolicy),lines:z.array(line).min(1).max(32).refine(v=>new Set(v.map(l=>l.id)).size===v.length&&new Set(v.map(l=>l.variant_id)).size===v.length&&v.reduce((n,l)=>n+l.quantity,0)<=1000)};
const common={organization_id:supplyID,purchase_order_id:supplyID,command_id:supplyID,evidence_sha256:supplySHA};
const base={...common,expected_version:positive};
const simple=(kind:"submit"|"confirm"|"start"|"cancel")=>z.object({kind:z.literal(kind),...base}).strict();
export const supplyCommand=z.discriminatedUnion("kind",[
 z.object({kind:z.literal("create"),...common,...planFields,destination_organization_id:supplyID,supplier_id:supplyID,currency:z.string().regex(/^[A-Z]{3}$/u),total_minor_units:unsigned}).strict().refine(v=>v.organization_id===v.destination_organization_id),
 z.object({kind:z.literal("plan"),...common,...planFields}).strict(),
 simple("submit"),simple("confirm"),simple("start"),simple("cancel"),
 z.object({kind:z.literal("register"),...base,line_id:supplyID,serial_number:text(128),vin:text(128).optional(),battery_serial_number:text(128).optional()}).strict(),
 z.object({kind:z.literal("milestone"),...base,unit_id:supplyID,target_state:z.enum(["assembly","quality","released","rejected"]),reason:text(2000).optional()}).strict().refine(v=>["released","rejected"].includes(v.target_state)?v.reason!==undefined:v.reason===undefined),
 ...(["ship","receive"] as const).map(kind=>z.object({kind:z.literal(kind),...base,shipment_id:supplyID,units:z.array(supplyID).min(1).max(100).refine(v=>new Set(v).size===v.length)}).strict()),
 ...(["quality","quality-reject","reinspect"] as const).map(kind=>z.object({kind:z.literal(kind),...base,unit_id:supplyID,reason:text(2000)}).strict())
]);
export type SupplyCommand=z.infer<typeof supplyCommand>;
export const supplyReceipt=z.object({purchase_order_id:supplyID,command_id:supplyID,version:positive,kind:z.enum(["planned",...supplyKinds]),actor:z.string().min(1).max(256),request_sha256:supplySHA,payload_sha256:supplySHA,payload:z.unknown(),recorded_at:z.iso.datetime({offset:true}),replay:z.boolean()}).strict();
export type SupplyReceipt=z.infer<typeof supplyReceipt>;
const unit=z.object({factory_review_state:z.enum(["pending","approved","rejected"]).optional(),factory_requester:z.string().min(1).max(256).optional(),receipt_review_state:z.enum(["pending","approved","rejected"]).optional(),receipt_requester:z.string().min(1).max(256).optional(),id:supplyID,line_id:supplyID,state:z.enum(["planned","assembly","quality","released","rejected","shipped","received"]),serial_number:text(128),stock_unit_id:supplyID.optional(),stock_state:z.enum(["in-transit","quarantine","available","reserved","sold","service","retired"]).optional(),shipment_id:supplyID.optional()}).strict();
export const supplyPlan=z.object({purchase_order_id:supplyID,destination_organization_id:supplyID,factory_organization_id:supplyID,supplier_id:supplyID,state:z.enum(["draft","submitted","accepted","in-production","shipped","received","cancelled"]),purchase_version:positive,version:positive,currency:z.string().regex(/^[A-Z]{3}$/u),total_minor_units:unsigned,policy_code:z.literal(supplyPolicy),demand_reference:text(1024),lines:z.array(line).min(1).max(32),units:z.array(unit).max(100),latest:supplyReceipt.optional(),next_unit_id:supplyID.optional()}).strict();
export type SupplyPlan=z.infer<typeof supplyPlan>;
export const pendingSupply=z.object({purchase_order_id:supplyID,command_id:supplyID,kind:supplyReceipt.shape.kind,version:positive,request_sha256:supplySHA}).strict();
export type PendingSupply=z.infer<typeof pendingSupply>;
export function supplyPermission(kind:SupplyCommand["kind"]){return ["confirm","start","register","milestone","ship"].includes(kind)?"supply:factory":kind==="receive"?"supply:receive":["quality","quality-reject"].includes(kind)?"supply:release":kind==="reinspect"?"supply:inspect":"supply:plan"}
export function supplySurface(kind:SupplyCommand["kind"]){return supplyPermission(kind)==="supply:factory"?"factory":"franchise"}
export function supplyPayload(c:SupplyCommand){const{organization_id:_,...body}=c;if(body.kind==="create"||body.kind==="plan"){const{kind:__,...plan}=body;return plan}return body}
// Narrow Go encoding/json correspondence: ASCII keys, validated strings and
// whole quantities1..1000 only. Not a general-purpose JSON canonicalizer.
export function supplyCanonical(value:unknown):string{
 const literal=(v:string)=>JSON.stringify(v).replace(/[<>&\u2028\u2029]/gu,x=>({"<":"\\u003c",">":"\\u003e","&":"\\u0026","\u2028":"\\u2028","\u2029":"\\u2029"}[x]!));
 if(value===null)return"null";if(typeof value==="string")return literal(value);
 if(typeof value==="number"&&Number.isInteger(value)&&value>=1&&value<=1000)return String(value);
 if(Array.isArray(value))return"["+value.map(supplyCanonical).join(",")+"]";
 if(typeof value==="object"){const o=value as Record<string,unknown>,keys=Object.keys(o).sort();if(keys.some(k=>!/^[a-z_][a-z0-9_]*$/u.test(k)))throw new Error("unsupported key");return"{"+keys.map(k=>literal(k)+":"+supplyCanonical(o[k])).join(",")+"}"}
 throw new Error("unsupported supply input");
}
export async function supplyDigest(value:string|Uint8Array){const bytes=typeof value==="string"?new TextEncoder().encode(value):value;return Array.from(new Uint8Array(await crypto.subtle.digest("SHA-256",new Uint8Array(bytes).buffer)),v=>v.toString(16).padStart(2,"0")).join("")}
export async function supplyReference(c:SupplyCommand):Promise<PendingSupply>{return pendingSupply.parse({purchase_order_id:c.purchase_order_id,command_id:c.command_id,kind:c.kind==="create"||c.kind==="plan"?"planned":c.kind,version:"expected_version"in c?(BigInt(c.expected_version)+1n).toString():"1",request_sha256:await supplyDigest(supplyCanonical(supplyPayload(c)))})}
export function supplyMatches(v:SupplyReceipt,p:PendingSupply,subject:string){return v.purchase_order_id===p.purchase_order_id&&v.command_id===p.command_id&&v.actor===subject&&v.kind===p.kind&&v.version===p.version&&v.request_sha256===p.request_sha256}
