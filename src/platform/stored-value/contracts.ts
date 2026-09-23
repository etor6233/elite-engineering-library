// AUTHORED transport validation. Amounts remain decimal strings; the selected
// backend/source profile owns every points, discount and eligibility decision.
import { z } from "zod";
export const storedID=z.string().regex(/^[A-Za-z0-9][A-Za-z0-9_.:-]{0,127}$/);
export const storedHash=z.string().regex(/^[0-9a-f]{64}$/);
const integer=z.string().regex(/^(0|[1-9][0-9]{0,18})$/).refine(v=>BigInt(v)<=9223372036854775807n);
const version=integer.refine(v=>v!=="0");
const signedInteger=z.string().regex(/^-?(0|[1-9][0-9]{0,18})$/).refine(v=>BigInt(v)>=-9223372036854775808n&&BigInt(v)<=9223372036854775807n);
const document=z.string().min(1).max(65536);
export const storedProfile=z.object({organization_id:storedID,profile_sha256:storedHash,profile_json:document}).strict();
export const storedOrder=z.object({organization_id:storedID,order_id:storedID,order_version:version,currency:z.literal("ARS"),gross_minor_units:integer,gift_minor_units:integer,discount_minor_units:integer,provider_due_minor_units:integer,approvals:z.object({ids:z.array(storedID).max(50),next_cursor:storedID.optional()}).strict()}).strict().refine(v=>BigInt(v.gross_minor_units)===BigInt(v.gift_minor_units)+BigInt(v.discount_minor_units)+BigInt(v.provider_due_minor_units));
export const storedApproval=z.object({review:z.object({operation_id:storedID,program_id:storedID,original_operation_id:z.union([storedID,z.literal("")]),currency:z.literal("ARS"),gross_minor_units:integer,gift_minor_units:integer,discount_minor_units:integer,entries:z.array(z.object({account_id:storedID,points_delta:z.string().regex(/^-?(0|[1-9][0-9]*)\.[0-9]{6}$/),applied_minor_units:signedInteger}).strict()).min(1).max(40)}).strict(),approval_id:storedID,payload_sha256:storedHash,state:z.enum(["pending","approved","rejected"]),replay:z.boolean(),payload_json:document,receipt_json:document,order_id:storedID,organization_id:storedID,operation:z.enum(["issue","accrue","redeem","reverse"]),requester:z.string().min(1).max(128),expires_at:z.iso.datetime({offset:true})}).strict();
export const storedFunding=z.object({funding_receipt_id:storedID,receipt_sha256:storedHash,receipt_json:document,organization_id:storedID,order_id:storedID,request_key:storedID}).strict();
export type StoredOrder=z.infer<typeof storedOrder>;export type StoredApproval=z.infer<typeof storedApproval>;export type StoredFunding=z.infer<typeof storedFunding>;
const base={organization_id:storedID};
export const storedCommand=z.discriminatedUnion("action",[
 z.object({...base,action:z.literal("propose"),operation_id:z.uuid(),operation:z.enum(["issue","accrue","redeem","reverse"]),program_id:storedID,order_id:storedID,expected_order_version:version,account_id:z.union([storedID,z.literal("")]),original_operation_id:z.union([storedID,z.literal("")]),profile_sha256:storedHash}).strict(),
 z.object({...base,action:z.literal("decision"),approval_id:storedID,payload_sha256:storedHash,approved:z.boolean(),reason:z.string().trim().min(1).max(2048)}).strict(),
 z.object({...base,action:z.literal("fund"),order_id:storedID,expected_order_version:version,request_key:z.uuid()}).strict(),
]);
export type StoredCommand=z.infer<typeof storedCommand>;
export const storedQuery=z.discriminatedUnion("kind",[
 z.object({...base,kind:z.literal("order"),order_id:storedID,after:storedID.optional()}).strict(),
 z.object({...base,kind:z.literal("approval"),approval_id:storedID}).strict(),
 z.object({...base,kind:z.literal("proposal-result"),operation_id:z.uuid()}).strict(),
 z.object({...base,kind:z.literal("funding-result"),request_key:z.uuid()}).strict(),
]);
export function formatStoredMoney(value:string){const v=integer.parse(value);const cents=v.padStart(3,"0");return `${cents.slice(0,-2)},${cents.slice(-2)} ARS`;}

export function formatStoredAdjustment(value:string){const v=signedInteger.parse(value);return v.startsWith("-")?"−"+formatStoredMoney(v.slice(1)):formatStoredMoney(v)}
