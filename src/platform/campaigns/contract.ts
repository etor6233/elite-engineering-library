// AUTHORED projection of the existing Go campaign owner; no client send authority.
import {z} from "zod";
export const campaignID=z.string().regex(/^[A-Za-z0-9_-]{16,128}$/);
const leadID=z.string().regex(/^[A-Za-z0-9][A-Za-z0-9_.:-]{0,127}$/);
const digest=z.string().regex(/^[a-f0-9]{64}$/);
const step=z.object({template_name:z.string().min(1).max(128),language_code:z.string().min(1).max(32),body_parameters:z.array(z.string().min(1).max(1024)).max(20),not_before:z.iso.datetime({offset:true}),expires_at:z.iso.datetime({offset:true})}).strict();
export const campaignSelection=z.object({campaign_id:campaignID,lead_ids:z.array(leadID).min(1).max(20),steps:z.array(step).min(1).max(3)}).strict();
export const campaignSnapshot=z.object({schema:z.literal("elite-whatsapp-campaign/v1"),request:z.object({campaign_id:campaignID,sources:z.array(z.string()).max(8),states:z.array(z.string()).max(3),members:z.array(z.object({lead_id:leadID,recipient:z.string().max(32)})).max(20),steps:z.array(step).max(3)}),requester:z.string().min(1).max(128),tenant_id:z.string().min(1),organization_id:z.string().min(1),connection_id:z.string().min(1),profile_sha256:digest,policy:z.object({schema:z.string(),consent_purpose:z.string(),policy_version:z.string(),max_members:z.number().int().min(1).max(20),max_steps:z.number().int().min(1).max(3)}),members:z.array(z.object({lead_id:leadID,recipient:z.string(),source:z.string(),lifecycle:z.string(),lead_updated_at:z.string(),subject_id:z.string(),external_id_hmac:digest,binding_version:z.string(),binding_policy:z.string(),consent_id:z.string(),consent_sha256:digest})).max(20)}).strict();
export const campaignCommand=z.discriminatedUnion("action",[
 z.object({action:z.literal("prepare"),selection:campaignSelection}).strict(),
 z.object({action:z.literal("create"),selection:campaignSelection,campaign_sha256:digest}).strict(),
 z.object({action:z.literal("resume"),campaign_id:campaignID,campaign_sha256:digest}).strict(),
 z.object({action:z.literal("decision"),campaign_id:campaignID,campaign_sha256:digest,approve:z.boolean(),reason:z.string().trim().min(1).max(2048)}).strict(),
 z.object({action:z.literal("stop"),campaign_id:campaignID,campaign_sha256:digest,reason:z.string().trim().min(1).max(2048)}).strict()
]);
export type CampaignCommand=z.infer<typeof campaignCommand>;
export type CampaignSnapshot=z.infer<typeof campaignSnapshot>;
export type CampaignSelection=z.infer<typeof campaignSelection>;
export const preparedCampaign=z.object({snapshot:campaignSnapshot,campaign_sha256:digest}).strict();
export type PreparedCampaign=z.infer<typeof preparedCampaign>;
export type CampaignStep=z.infer<typeof step>;
export const campaignPage=z.object({kind:z.enum(["campaigns","audience"]),campaigns:z.array(z.object({campaign_id:campaignID,template_name:z.string(),language_code:z.string(),members:z.number().int().nonnegative(),steps:z.number().int().nonnegative(),created_at:z.iso.datetime({offset:true}),stopped:z.boolean()})).max(25).optional(),audience:z.array(z.object({lead_id:leadID,name:z.string().max(256),contact_hint:z.string().max(32),source:z.string(),state:z.string()})).max(25).optional(),next_cursor:leadID.optional()});
export type CampaignPage=z.infer<typeof campaignPage>;
export const campaignTemplates=z.object({templates:z.array(z.object({name:z.string(),language_code:z.string(),body_parameter_count:z.number().int().min(0).max(20)})).max(100),max_members:z.number().int().min(1).max(20),max_steps:z.number().int().min(1).max(3)});
export type CampaignTemplates=z.infer<typeof campaignTemplates>;
export const campaignStatus=z.object({campaign_sha256:digest,snapshot:campaignSnapshot,stopped:z.boolean(),items:z.array(z.object({lead_id:leadID,step:z.number().int(),delivery_key:z.string(),preparation:z.string(),schedule:z.object({approval_state:z.string(),delivery_state:z.string(),outcome:z.string(),cancelled:z.boolean(),job_completed:z.boolean()}).passthrough().optional()})).max(60),conversions:z.array(z.object({lead_id:leadID,quotation_id:z.string(),order_id:z.string(),order_state:z.string(),accepted_at:z.string(),evidence_sha256:digest})).max(20),conversion_meaning:z.string()});
export type CampaignStatus=z.infer<typeof campaignStatus>;
export const campaignBatch=z.object({campaign_id:campaignID,campaign_sha256:digest,complete:z.boolean(),items:z.array(z.object({lead_id:leadID,step:z.number().int(),delivery_key:z.string(),request_sha256:digest.optional(),state:z.string(),code:z.string().optional()})).max(60)});
export function campaignPermissions(permissions:readonly string[]){
 const has=(...names:string[])=>names.every(n=>permissions.includes(n)||permissions.includes("*"));
 return {read:has("marketing:read","notification:read"),prepare:has("marketing:request","notification:request","lead:read"),review:has("marketing:approve","notification:approve"),stop:has("marketing:cancel")};
}
export function campaignActionAllowed(action:CampaignCommand["action"],permissions:readonly string[]){const p=campaignPermissions(permissions);return p.read&&(action==="decision"?p.review:action==="stop"?p.stop:p.prepare)}
export function campaignHumanState(item:CampaignStatus["items"][number]){
 const s=item.schedule;if(!s)return "Preparación pendiente";
 if(s.cancelled||s.outcome==="SUPPRESSED")return "Detenido";
 if(s.outcome==="RECONCILIATION_REQUIRED"||s.delivery_state==="uncertain")return "Resultado por confirmar";
 if(s.notification && typeof s.notification==="object" && "delivery_status" in s.notification){const p=s.notification.delivery_status;if(p==="read")return "Leído";if(p==="delivered")return "Entregado";if(p==="failed")return "No entregado"}
 if(s.delivery_state==="accepted"||s.outcome==="ACCEPTED")return "Aceptado por WhatsApp";
 if(s.approval_state==="rejected")return "Rechazado";
 if(s.approval_state==="approved")return "Programado";
 return "En revisión";
}
