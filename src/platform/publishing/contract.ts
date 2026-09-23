// AUTHORED UI contract for the existing, admitted Page text publishing owner.
import { z } from "zod";
const id=z.string().regex(/^[A-Za-z0-9_-]{1,128}$/), sha=z.string().regex(/^[a-f0-9]{64}$/);
const instant=z.string().datetime({offset:true});
const message=z.string().refine(v=>v.trim().length>0&&new TextEncoder().encode(v).length<=16384,"bounded text required");
export const publicationRequest=z.object({intent:z.object({tenant_id:id,page_id:z.string().regex(/^[1-9][0-9]{0,31}$/),profile_sha256:sha,approval_id:id,delivery_key:z.string().regex(/^social_[a-f0-9]{64}$/),operation:z.enum(["publish","revoke"]),message,content_sha256:sha}).strict(),scheduled_at:instant,expires_at:instant,original_approval_id:z.string().max(128),provider_reference:z.string().max(100)}).strict().refine(v=>Date.parse(v.expires_at)>Date.parse(v.scheduled_at),"expiry after schedule");
export const preparedPublication=z.object({request:publicationRequest,request_sha256:sha}).strict();
export type PreparedPublication=z.infer<typeof preparedPublication>;
export const publicationReceipt=z.object({tenant_id:id,page_id:z.string(),profile_sha256:sha,approval_id:id,delivery_key:z.string(),provider_reference:z.string(),state:z.enum(["published","revoked"]),content_sha256:sha,evidence_sha256:sha}).strict();
const approvalState=z.enum(["pending","approved","rejected"]);
export const publicationRow=z.object({approval_id:id,message,operation:z.enum(["publish","revoke"]),requested_by_self:z.boolean(),approval_state:approvalState,scheduled_at:instant,expires_at:instant}).strict();
export const publicationsPage=z.object({organization_id:id,items:z.array(publicationRow).max(25),next_cursor:id.optional(),publish_enabled:z.boolean(),revoke_enabled:z.boolean()}).strict();
export type PublicationsPage=z.infer<typeof publicationsPage>;
export const publicationDetail=z.object({organization_id:id,requested_by_self:z.boolean(),status:z.object({approval_id:id,approval_state:approvalState,request_sha256:sha,payload:publicationRequest,delivery_state:z.enum(["","sending","unknown","accepted","failed_terminal"]),receipt:publicationReceipt.optional(),terminal_code:z.string().max(128).optional()}).strict()}).strict().superRefine((value,ctx)=>{
 const s=value.status,i=s.payload.intent,r=s.receipt;
 if(s.approval_id!==i.approval_id)ctx.addIssue({code:"custom",message:"request binding"});
 if(s.delivery_state==="accepted"&&!r)ctx.addIssue({code:"custom",message:"accepted receipt required"});
 if(r&&(r.tenant_id!==i.tenant_id||r.page_id!==i.page_id||r.profile_sha256!==i.profile_sha256||r.approval_id!==i.approval_id||r.delivery_key!==i.delivery_key||r.content_sha256!==i.content_sha256||r.state!==(i.operation==="publish"?"published":"revoked")||!new RegExp(`^${i.page_id}_[1-9][0-9]{0,63}$`).test(r.provider_reference)))ctx.addIssue({code:"custom",message:"receipt binding"});
});
export type PublicationDetail=z.infer<typeof publicationDetail>;
const draft=z.object({approval_id:z.string().regex(/^[A-Za-z0-9_-]{16,128}$/),operation:z.enum(["publish","revoke"]),message:z.string().max(16384),original_approval_id:z.string().max(128),scheduled_at:instant,expires_at:instant}).strict().superRefine((v,ctx)=>{
 if(Date.parse(v.expires_at)<=Date.parse(v.scheduled_at))ctx.addIssue({code:"custom",message:"expiry after schedule"});
 if(v.operation==="publish"&&(!message.safeParse(v.message).success||v.original_approval_id!==""))ctx.addIssue({code:"custom",message:"publish draft"});
 if(v.operation==="revoke"&&(v.message!==""||!id.safeParse(v.original_approval_id).success))ctx.addIssue({code:"custom",message:"revoke draft"});
});
export const publishingCommand=z.discriminatedUnion("action",[
 z.object({action:z.literal("prepare"),draft}).strict(),
 z.object({action:z.literal("submit"),prepared:preparedPublication}).strict(),
 z.object({action:z.literal("decide"),approval_id:id,request_sha256:sha,approve:z.boolean(),reason:z.string().max(1000)}).strict(),
 z.object({action:z.literal("reconcile"),approval_id:id}).strict(),
]);
export const publicationQuery=z.union([z.object({id}).strict(),z.object({after:id.optional()}).strict()]);
export function publicState(detail:PublicationDetail){const s=detail.status;if(s.delivery_state==="accepted"&&s.receipt)return s.receipt.state;if(s.delivery_state==="unknown")return "uncertain";if(s.delivery_state==="failed_terminal")return "stopped";if(s.approval_state==="rejected")return "rejected";if(s.approval_state==="pending")return "pending";return "scheduled"}
