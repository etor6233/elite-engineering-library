import {z} from "zod";
import {catalogCanonical,catalogSHA} from "@/platform/catalog/authoring";
export const networkID=z.string().regex(/^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$/u);
const scope=networkID.or(z.literal(""));
const version=z.string().regex(/^[1-9][0-9]{0,18}$/u).refine(v=>BigInt(v)<9223372036854775807n);
const text=(max:number)=>z.string().min(2).refine(v=>v.trim()===v&&!/[\u0000-\u001f\u007f-\u009f]/u.test(v)&&new TextEncoder().encode(v).length<=max);
const date=z.string().regex(/^[0-9]{4}-[0-9]{2}-[0-9]{2}$/u).refine(v=>v.slice(0,4)!=="0000"&&Number.isFinite(Date.parse(v+"T00:00:00Z"))&&new Date(v+"T00:00:00Z").toISOString().slice(0,10)===v);
const base={command_id:networkID,scope_organization_id:scope,entity_id:networkID};
export const networkCommand=z.discriminatedUnion("action",[
 z.object({...base,action:z.literal("create-organization"),code:z.string().max(128).regex(/^[a-z][a-z0-9]*(-[a-z0-9]+)*$/u),display_name:text(100),type:z.enum(["enterprise","franchisor","franchisee","factory","warehouse","store","service_center"])}).strict(),
 z.object({...base,action:z.literal("create-agreement"),scope_organization_id:networkID,territory_code:networkID,terms_version:networkID,starts_on:date,ends_on:date.optional()}).strict(),
 z.object({...base,action:z.literal("transition-organization"),scope_organization_id:networkID,current:networkID,target:networkID,version}).strict(),
 z.object({...base,action:z.literal("transition-agreement"),scope_organization_id:networkID,current:networkID,target:networkID,version}).strict()
]).refine(c=>c.action!=="transition-organization"||c.scope_organization_id===c.entity_id)
 .refine(c=>c.action!=="create-organization"||(["enterprise","franchisor"].includes(c.type)?c.scope_organization_id==="":c.scope_organization_id!==""))
 .refine(c=>c.action!=="create-agreement"||!c.ends_on||c.ends_on>c.starts_on);
export type NetworkCommand=z.infer<typeof networkCommand>;
export const networkEntity=z.object({
 kind:z.enum(["organization","agreement"]),id:networkID,organization_id:networkID,version:z.string().regex(/^[1-9][0-9]{0,18}$/u).refine(v=>BigInt(v)<=9223372036854775807n),
 state:networkID,code:networkID.optional(),display_name:text(100).optional(),type:networkID.optional(),parent_organization_id:networkID.optional(),
 territory_code:networkID.optional(),terms_version:networkID.optional(),starts_on:date.optional(),ends_on:date.optional()
}).strict();
export type NetworkEntity=z.infer<typeof networkEntity>;
export const networkReceipt=z.object({command_id:networkID,action:z.enum(["create-organization","create-agreement","transition-organization","transition-agreement"]),scope_organization_id:scope,actor:z.string().min(1).max(256),request_sha256:z.string().regex(/^[a-f0-9]{64}$/u),entity:networkEntity,recorded_at:z.string().datetime({offset:true}),replay:z.boolean()}).strict();
export const pendingNetwork=z.object({command_id:networkID,action:networkReceipt.shape.action,scope_organization_id:scope,entity_id:networkID,version:networkEntity.shape.version,state:networkID,request_sha256:networkReceipt.shape.request_sha256}).strict();
export type PendingNetwork=z.infer<typeof pendingNetwork>;
export function networkPermission(action:string){return action.endsWith("-organization")?"network:admin":"franchise:write"}
export function networkPayload(c:NetworkCommand){return Object.fromEntries(Object.entries(c).filter(([k,v])=>v!==undefined&&(v!==""||k==="scope_organization_id")))}
export async function networkReference(c:NetworkCommand):Promise<PendingNetwork>{return pendingNetwork.parse({command_id:c.command_id,action:c.action,scope_organization_id:c.scope_organization_id,entity_id:c.entity_id,version:"version"in c?(BigInt(c.version)+1n).toString():"1",state:"target"in c?c.target:c.action==="create-organization"?"provisioning":"draft",request_sha256:await catalogSHA(catalogCanonical(networkPayload(c)))})}
export function networkMatches(v:z.infer<typeof networkReceipt>,p:PendingNetwork,subject:string){
 return v.command_id===p.command_id&&v.action===p.action&&v.scope_organization_id===p.scope_organization_id&&v.actor===subject&&v.request_sha256===p.request_sha256&&v.entity.id===p.entity_id&&v.entity.version===p.version&&v.entity.state===p.state&&v.entity.kind===(p.action.endsWith("-organization")?"organization":"agreement")&&v.entity.organization_id===(v.entity.kind==="organization"?p.entity_id:p.scope_organization_id);
}
// Existing fulfillment.Service transition tables, surfaced for forms; Go remains authority.
export function networkTargets(v:NetworkEntity){const table:Record<string,string[]>=v.kind==="organization"?{provisioning:["active","closed"],active:["suspended","closed"],suspended:["active","closed"]}:{draft:["active","terminated"],active:["suspended","terminated","expired"],suspended:["active","terminated","expired"]};return table[v.state]??[]}
