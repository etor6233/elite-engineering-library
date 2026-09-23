// AUTHORED typed form/receipt binding. Existing Go owners govern all effects.
import {z} from "zod";
import {publishedCatalogSchema} from "./schema";
export const catalogID=z.string().regex(/^[A-Za-z0-9][A-Za-z0-9._:-]{0,63}$/u);
const sha=z.string().regex(/^[0-9a-f]{64}$/u);
const code=z.string().max(64).regex(/^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$/u);
const text=(max:number)=>z.string().min(1).max(max).refine(v=>v.trim()===v&&!/[\x00-\x1f\x7f]/u.test(v)&&new TextDecoder().decode(new TextEncoder().encode(v))===v);
const unsigned=z.string().regex(/^(0|[1-9][0-9]{0,18})$/u).refine(v=>BigInt(v)<=9223372036854775807n);
const instant=z.iso.datetime({offset:true}).transform(v=>new Date(v).toISOString().replace(/\.000Z$/u,"Z").replace(/(\.\d*?[1-9])0+Z$/u,"$1Z"));
export const reviewStage=z.enum(["legal","technical","media","publication"]);
const common={organization_id:catalogID,command_id:catalogID};
export const sourceForm=z.object({
 action:z.literal("source"),...common,
 model:z.object({id:z.literal(""),code,displayName:text(160).refine(v=>v.length>=2),vehicleClass:z.enum(["motorcycle","bicycle","scooter","utility","other"]),specification:z.object({description:text(4000)}).strict()}).strict(),
 variants:z.array(z.object({code,display_name:text(160),battery_specification:z.object({description:text(4000)}).strict(),amount_minor_units:unsigned,tax_mode:z.enum(["inclusive","exclusive","not-applicable"])}).strict()).min(1).max(16),
 valid_from:instant,valid_until:instant.nullable()
}).strict().refine(v=>v.valid_until===null||Date.parse(v.valid_until)>Date.parse(v.valid_from)).refine(v=>new Set(v.variants.map(x=>x.code)).size===v.variants.length);
export const catalogCommand=z.discriminatedUnion("action",[
 sourceForm,
 z.object({action:z.literal("draft"),...common,price_book_id:catalogID,models:z.array(z.object({model_id:catalogID,media_id:catalogID}).strict()).min(1).max(32)}).strict(),
 z.object({action:z.literal("review"),...common,draft_id:catalogID,stage:reviewStage,snapshot_sha256:sha,approved:z.boolean(),reason:text(2000),evidence_sha256:sha}).strict(),
 z.object({action:z.literal("publish"),...common,draft_id:catalogID,snapshot_sha256:sha,expected_generation:unsigned,reason:text(2000)}).strict()
]);
export type CatalogCommand=z.infer<typeof catalogCommand>;
export const catalogReceipt=z.object({command_id:catalogID,actor:z.string().min(1).max(256),kind:z.enum(["media","source","draft","review","publish"]),resource_id:catalogID,request_sha256:sha,snapshot_sha256:sha.optional(),generation:unsigned,replay:z.boolean(),effective_price_book_id:catalogID.optional(),source_price_book_id:catalogID.optional(),source_variant_ids:z.array(catalogID).min(1).max(16).optional()}).strict().superRefine((v,c)=>{
 if(v.kind==="source"&&(!v.source_price_book_id||!v.source_variant_ids))c.addIssue({code:"custom",message:"missing source references"});
 if(["draft","review","publish"].includes(v.kind)&&!v.snapshot_sha256)c.addIssue({code:"custom",message:"missing snapshot"});
 if(v.kind==="publish"&&(!v.effective_price_book_id||v.generation==="0"))c.addIssue({code:"custom",message:"missing publication"});
});
export type CatalogReceipt=z.infer<typeof catalogReceipt>;
const snapshot=z.object({tenant_code:code,schema:z.literal("elite-catalog-snapshot/v1"),profile:z.object({schema:z.literal("elite-catalog-publication/v1"),tenant_id:z.uuid(),organization_id:catalogID,market:z.string().regex(/^[A-Z]{2}$/u),currency:z.string().regex(/^[A-Z]{3}$/u),origin:z.url(),policy_code:z.literal("catalog-reference-v1")}).strict(),price_book_id:catalogID,valid_from:z.iso.datetime({offset:true}),valid_until:z.iso.datetime({offset:true}).nullable(),models:publishedCatalogSchema.shape.models,variants:publishedCatalogSchema.shape.variants}).strict();
export const catalogDraft=z.object({id:catalogID,maker:z.string().min(1).max(256),sha256:sha,snapshot,reviews:z.object({legal:z.enum(["pending","approved","rejected"]),technical:z.enum(["pending","approved","rejected"]),media:z.enum(["pending","approved","rejected"]),publication:z.enum(["pending","approved","rejected"])}).strict()}).strict();
export type CatalogDraft=z.infer<typeof catalogDraft>;
export const pendingCatalog=z.object({command_id:catalogID,kind:catalogReceipt.shape.kind,request_sha256:sha,resource_id:catalogID.optional(),snapshot_sha256:sha.optional()}).strict();
export type PendingCatalog=z.infer<typeof pendingCatalog>;
export function commandPayload(c:CatalogCommand){const{action:_,organization_id:__,...body}=c;return body}
// Transport correspondence with Go encoding/json for this explicitly bounded
// string/boolean/null schema and ASCII object keys; not a generic JSON/JCS engine.
export function catalogCanonical(value:unknown):string{
 const literal=(v:string)=>JSON.stringify(v).replace(/[<>&\u2028\u2029]/gu,x=>({"<":"\\u003c",">":"\\u003e","&":"\\u0026","\u2028":"\\u2028","\u2029":"\\u2029"}[x]!));
 if(value===null)return"null";
 if(typeof value==="string")return literal(value);
 if(typeof value==="boolean")return String(value);
 if(Array.isArray(value))return"["+value.map(catalogCanonical).join(",")+"]";
 if(typeof value==="object"){
  const o=value as Record<string,unknown>;const keys=Object.keys(o).sort();
  if(keys.some(k=>!/^[A-Za-z_][A-Za-z0-9_]*$/u.test(k)))throw new Error("unsupported catalog key");
  return"{"+keys.map(k=>literal(k)+":"+catalogCanonical(o[k])).join(",")+"}";
 }
 throw new Error("unsupported catalog value");
}
export async function catalogSHA(value:Uint8Array|string){const bytes=typeof value==="string"?new TextEncoder().encode(value):value;return Array.from(new Uint8Array(await crypto.subtle.digest("SHA-256",new Uint8Array(bytes).buffer)),x=>x.toString(16).padStart(2,"0")).join("")}
export async function commandReference(c:CatalogCommand):Promise<PendingCatalog>{const body=commandPayload(c);return pendingCatalog.parse({command_id:c.command_id,kind:c.action,request_sha256:await catalogSHA(catalogCanonical(body)),...("draft_id"in c?{resource_id:c.draft_id,snapshot_sha256:c.snapshot_sha256}:{})})}
export function receiptMatches(v:CatalogReceipt,p:PendingCatalog,subject:string){return v.command_id===p.command_id&&v.actor===subject&&v.kind===p.kind&&v.request_sha256===p.request_sha256&&(!p.resource_id||p.resource_id===v.resource_id)&&(!p.snapshot_sha256||p.snapshot_sha256===v.snapshot_sha256)}
