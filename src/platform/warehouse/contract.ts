// AUTHORED presentation contract; all quantities remain exact decimal strings.
import {z} from "zod";
import rawActions from "./actions.json";
export type Field={key:string;label:string;kind:string;source?:string;value?:string|number;choices?:string[];optional?:boolean;hidden?:boolean;parentField?:string;organizationField?:string};
export type Action={id:string;label:string;group:string;path:string;fields:Field[];selection?:string;activityType?:string;recovery:string;requestRecovery?:boolean;terminal?:string|null};
export const actions=rawActions as Action[];
export const datasets=["organizations","items","bins","policies","uoms","balances","reservations","activities","activity-lines","transfers","crossdock","entries","orders","order-lines","bindings","shipments","packaging","transfer-shipments","transfer-receipts"] as const;
export type Dataset=typeof datasets[number];
export const idSchema=z.string().min(1).max(128).regex(/^[^\x00-\x20\x7f/\\?#]+$/u);
const text=z.string().min(1).max(128).regex(/^[^\x00-\x1f\x7f]+$/u);
export const rowSchema=z.object({id:z.string().min(1).max(384),label:z.string().max(640),state:z.string().max(64),version:z.string().max(20).regex(/^[1-9]\d*$/u),data:z.record(z.string(),z.unknown())}).strict();
export const pageSchema=z.object({schema:z.literal("warehouse-workspace/v1"),dataset:z.enum(datasets),organization_id:idSchema,rows:z.array(rowSchema).max(50),next_cursor:z.string().max(4096).nullable()}).strict();
export type Row=z.infer<typeof rowSchema>;export type Page=z.infer<typeof pageSchema>;
function fieldSchema(f:Field):z.ZodType{
 let s:z.ZodType;
 if(f.kind==="constant")s=z.literal(f.value as string|number);
 else if(f.kind==="boolean")s=z.boolean();
 else if(f.kind==="integer")s=z.number().int().min(0).max(f.key==="due_date_days"?365:1000000);
 else if(f.kind==="version")s=z.number().int().min(1).max(Number.MAX_SAFE_INTEGER);
 else if(f.kind==="decimal")s=z.string().regex(/^(0|[1-9]\d{0,13})(\.\d{1,6})?$/u).refine(v=>f.key==='min_quantity'||/[1-9]/u.test(v));
 else if(f.kind==="cost")s=z.string().regex(/^(0|[1-9]\d{0,15})(\.\d{1,4})?$/u).refine(v=>/[1-9]/u.test(v));
 else if(f.kind==="date")s=z.string().datetime({offset:true});
 else if(f.kind==="choice")s=z.enum(f.choices as [string,...string[]]);
 else s=text;
 return f.optional?s.optional():s;
}
export function parseCommand(value:unknown){
 const outer=z.object({action:z.string(),organization:idSchema,id:idSchema.optional(),payload:z.record(z.string(),z.unknown())}).strict().parse(value);
 const action=actions.find(a=>a.id===outer.action);if(!action)throw new Error("UNKNOWN_ACTION");
 if(Boolean(action.selection)!==Boolean(outer.id))throw new Error("INVALID_SELECTION");
 const shape:Record<string,z.ZodType>={};for(const field of action.fields)shape[field.key]=fieldSchema(field);
 const payload=z.object(shape).strict().parse(outer.payload);
 for(const key of ["organization_id","from_organization_id"])if(payload[key]!==undefined&&key!=="from_organization_id"&&payload[key]!==outer.organization)throw new Error("WRONG_ORGANIZATION");
 if(action.id==="transfer-create"&&payload.from_organization_id!==outer.organization)throw new Error("WRONG_ORGANIZATION");
 return {...outer,actionDefinition:action,payload};
}
