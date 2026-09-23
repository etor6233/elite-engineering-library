// AUTHORED bounded projection of the versioned document owner. Decimal values stay strings.
import {z} from "zod";
export const MAX_DOCUMENT_BYTES=2*1024*1024;
export const hash=z.string().regex(/^[a-f0-9]{64}$/u);
export const documentId=z.uuid().refine(value=>value!=="00000000-0000-0000-0000-000000000000");
export const classId=z.string().regex(/^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$/u).max(100);
export const schemaVersion=z.string().regex(/^[0-9]+(?:\.[0-9]+)*$/u).max(20);
const text=(max:number)=>z.string().min(1).refine(value=>value.trim()===value&&!value.includes("\0")&&new TextEncoder().encode(value).length<=max);
export const mode=z.enum(["FIXTURE","PROVIDER","TYPED_FIXTURE"]);
export const fieldValues=z.record(z.string().regex(/^[a-z][a-z0-9_]*$/u).max(80),z.string().max(4096)).refine(value=>Object.keys(value).length>0&&Object.keys(value).length<=50);
export const classSchema=z.object({class_id:classId,schema_version:schemaVersion,label_es:text(180),task_id:text(180),fields:z.array(z.object({key:z.string().regex(/^[a-z][a-z0-9_]*$/u),label_es:text(180),type:z.enum(["text","date","decimal","integer","currency"]),max_length:z.number().int().min(1).max(4096)}).strict()).min(1).max(50),extraction_method:z.string().max(128),automatic_storage_authorized:z.literal(false),fixture_sha256:hash.optional()}).strict().refine(value=>new Set(value.fields.map(f=>f.key)).size===value.fields.length);
export const catalogSchema=z.object({schema:z.literal("document-class-catalog/v1"),mode,profile_sha256:hash,classes:z.array(classSchema).min(1).max(21)}).strict().refine(value=>new Set(value.classes.map(c=>c.class_id)).size===value.classes.length);
const proposal=z.object({schema:z.enum(["document-review/v1","document-review/v2"]),document_id:documentId,organization_id:text(128),original_sha256:hash,profile_sha256:hash,evidence_sha256:hash,mode,class_id:classId.optional(),schema_version:schemaVersion.optional(),fields:fieldValues}).strict();
export const viewSchema=z.object({document_id:documentId,name:text(128),uploader:text(128),original_sha256:hash,profile_sha256:hash,mode,state:z.enum(["QUARANTINED","QUARANTINE_TERMINAL","REVIEW_REQUIRED","REVIEW_PENDING","REJECTED","PERSISTED"]),class_id:classId,schema_version:schemaVersion,evidence_sha256:hash.optional(),suggested:fieldValues.optional(),proposal:proposal.optional(),payload_sha256:hash.optional(),reviewer:text(128).optional(),reason:z.string().max(2048).optional()}).strict().superRefine((view,ctx)=>{
 const p=view.proposal;
 if(p&&(p.document_id!==view.document_id||p.original_sha256!==view.original_sha256||p.profile_sha256!==view.profile_sha256||p.evidence_sha256!==view.evidence_sha256||p.mode!==view.mode||!view.payload_sha256||p.schema==="document-review/v2"&&(p.class_id!==view.class_id||p.schema_version!==view.schema_version)))ctx.addIssue({code:"custom",message:"Proposal binding mismatch"});
 if(["REVIEW_PENDING","PERSISTED","REJECTED"].includes(view.state)&&!p)ctx.addIssue({code:"custom",message:"Proposal required"});
 if(["PERSISTED","REJECTED"].includes(view.state)&&(!view.reviewer||view.reviewer===view.uploader))ctx.addIssue({code:"custom",message:"Independent decision evidence required"});
});
export const inboxSchema=z.object({schema:z.literal("document-inbox/v1"),items:z.array(viewSchema).max(50),next_cursor:z.string().max(2048).nullable()}).strict();
export const accessSchema=z.object({scope:hash,subject:text(128),canRead:z.boolean(),canWrite:z.boolean(),canProcess:z.boolean(),canReview:z.boolean(),mode,profile_sha256:hash}).strict();
export const commandSchema=z.discriminatedUnion("action",[
 z.object({id:documentId,action:z.literal("process")}).strict(),
 z.object({id:documentId,action:z.literal("review"),evidence_sha256:hash,fields:fieldValues}).strict(),
 z.object({id:documentId,action:z.literal("decision"),payload_sha256:hash,approved:z.boolean(),reason:text(2048)}).strict()
]);
export type DocumentClass=z.infer<typeof classSchema>;export type DocumentCatalog=z.infer<typeof catalogSchema>;export type DocumentView=z.infer<typeof viewSchema>;export type DocumentAccess=z.infer<typeof accessSchema>;export type DocumentInbox=z.infer<typeof inboxSchema>;export type DocumentCommand=z.infer<typeof commandSchema>;export type DocumentFields=Record<string,string>;
export const pendingSchema=z.object({id:documentId,action:z.enum(["receive","process","review","decision"]),expected:hash,original_sha256:hash,class_id:classId,schema_version:schemaVersion,evidence_sha256:hash.optional(),approved:z.boolean().optional(),metadata_sha256:hash.optional(),reason_sha256:hash.optional()}).strict().refine(p=>p.action==="decision"?typeof p.approved==="boolean"&&Boolean(p.reason_sha256):p.approved===undefined&&p.reason_sha256===undefined).refine(p=>p.action!=="receive"||Boolean(p.metadata_sha256));
export type PendingDocument=z.infer<typeof pendingSchema>;
export async function sha256(bytes:Uint8Array){return Array.from(new Uint8Array(await crypto.subtle.digest("SHA-256",new Uint8Array(bytes).buffer)),b=>b.toString(16).padStart(2,"0")).join("")}
export const textHash=(value:string)=>sha256(new TextEncoder().encode(value));
export const fieldsHash=(value:DocumentFields)=>textHash(JSON.stringify(Object.fromEntries(Object.keys(value).sort().map(key=>[key,value[key]]))));
export function validFields(values:DocumentFields,kind:DocumentClass){return Object.keys(values).length===kind.fields.length&&kind.fields.every(field=>{const value=values[field.key];if(typeof value!=="string"||!value||value.trim()!==value||new TextEncoder().encode(value).length>field.max_length||/[\u0000-\u001f\u007f]/u.test(value))return false;switch(field.type){case"currency":return /^[A-Z]{3}$/u.test(value);case"integer":return /^\d+$/u.test(value);case"decimal":return /^-?\d+(?:\.\d+)?$/u.test(value);case"date":return /^\d{4}-\d{2}-\d{2}$/u.test(value)&&Number.isFinite(new Date(value+"T00:00:00Z").getTime())&&new Date(value+"T00:00:00Z").toISOString().slice(0,10)===value;default:return true}})}
export async function reconciles(p:PendingDocument,view:DocumentView){if(p.id!==view.document_id||p.original_sha256!==view.original_sha256||p.class_id!==view.class_id||p.schema_version!==view.schema_version)return false;
 if(p.action==="receive")return p.expected===view.original_sha256;
 if(p.action==="process")return view.state!=="QUARANTINED"&&p.expected===view.original_sha256;
 if(p.action==="review")return Boolean(view.proposal&&view.evidence_sha256===p.evidence_sha256&&await fieldsHash(view.proposal.fields)===p.expected);
 return view.payload_sha256===p.expected&&await textHash(view.reason??"")===p.reason_sha256&&(p.approved?view.state==="PERSISTED":view.state==="REJECTED");
}
