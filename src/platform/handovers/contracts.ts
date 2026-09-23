// AUTHORED DTO and recovery glue; all business decisions belong to the Go owner.
import { z } from "zod";
export const handoverID=z.string().min(1).max(128).regex(/^[A-Za-z0-9_-]+$/);
const hash=z.string().regex(/^[a-f0-9]{64}$/),version=z.number().int().positive().max(Number.MAX_SAFE_INTEGER);
const date=z.iso.datetime({offset:true});
const fundingFields={payment_attempt_id:z.union([handoverID,z.literal("")]),funding_receipt_id:handoverID.optional()};
const oneFunding=(v:{payment_attempt_id:string;funding_receipt_id?:string|undefined})=>Boolean(v.payment_attempt_id)!==Boolean(v.funding_receipt_id);
export const handoverViewSchema=z.object({id:handoverID,organization_id:handoverID,order_id:handoverID,state:z.enum(["prepared","presented","accepted","rejected"]),version,checklist_id:z.string().max(128).optional(),checklist_version:z.number().int().nonnegative().optional(),checklist_completed_at:date.optional(),customer_accepted_at:date.optional()});
export const handoverContextSchema=z.object({organization_id:handoverID,order_id:handoverID,handover:handoverViewSchema.optional(),can_prepare:z.boolean(),release_effect:z.enum(["READ_ONLY_ELIGIBILITY","COMMIT_COMMERCIAL_RELEASE_RECEIPT"]),evaluated_at:date}).superRefine((v,ctx)=>{
 if(v.can_prepare===Boolean(v.handover)||v.handover&&(v.handover.organization_id!==v.organization_id||v.handover.order_id!==v.order_id))ctx.addIssue({code:"custom",message:"inconsistent context"});
});
export const backendHandoverContextSchema=handoverContextSchema.safeExtend({order_line_id:handoverID,...fundingFields,observation_sha256:hash}).refine(oneFunding);
export const preparationSchema=z.object({handover:handoverViewSchema,order_line_id:handoverID,reservation_id:handoverID,...fundingFields,observation_sha256:hash,contract_id:z.string().min(1).max(128),contract_sha256:hash,prepared_by:z.string().min(1).max(256),prepared_at:date}).refine(oneFunding);
export const commercialReceiptSchema=z.object({id:handoverID,organization_id:handoverID,handover_id:handoverID,order_id:handoverID,...fundingFields,observation_sha256:hash,observation_generation:version,handover_version:version,acceptance_sha256:hash,checklist_id:handoverID,checklist_version:version,contract_id:z.string().min(1).max(128),contract_sha256:hash,effect:z.literal("COMMIT_COMMERCIAL_RELEASE_RECEIPT"),released_by:z.string().min(1).max(256),recorded_at:date,valid_until:date}).refine(oneFunding).refine(v=>Date.parse(v.valid_until)>Date.parse(v.recorded_at));
export const currentReleaseSchema=z.object({receipt:commercialReceiptSchema,evaluated_at:date,current:z.boolean()}).refine(v=>!v.current||Date.parse(v.evaluated_at)<Date.parse(v.receipt.valid_until));
export type HandoverContext=z.infer<typeof handoverContextSchema>;
export type CommercialReceipt=z.infer<typeof commercialReceiptSchema>;
export type CurrentRelease=z.infer<typeof currentReleaseSchema>;
export const operationMarkerSchema=z.discriminatedUnion("action",[
 z.object({action:z.literal("prepare"),orderId:handoverID,requestKey:z.uuid()}).strict(),
 z.object({action:z.literal("release"),orderId:handoverID,handoverId:handoverID,requestKey:z.uuid()}).strict(),
]);
export type OperationMarker=z.infer<typeof operationMarkerSchema>;
export function retainOperation(storage:Pick<Storage,"getItem"|"setItem">,key:string,value:OperationMarker):OperationMarker{
 const parsed=operationMarkerSchema.parse(value),existing=storage.getItem(key);
 if(existing!==null){const old=operationMarkerSchema.parse(JSON.parse(existing));if(old.action!==value.action||old.orderId!==value.orderId||old.action==="release"&&value.action==="release"&&old.handoverId!==value.handoverId)throw new Error("operation scope mismatch");return old;}
 storage.setItem(key,JSON.stringify(parsed));if(storage.getItem(key)!==JSON.stringify(parsed))throw new Error("operation reference not retained");return parsed;
}
