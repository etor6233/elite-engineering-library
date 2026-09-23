// AUTHORED bounded transport binding; existing Go kernel owns article lifecycle.
import{z}from"zod";
import{catalogCanonical,catalogSHA}from"@/platform/catalog/authoring";
export const helpID=z.string().regex(/^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$/u);
const version=z.string().regex(/^[1-9][0-9]{0,18}$/u).refine(v=>BigInt(v)<=9223372036854775807n);
const expected=version.refine(v=>BigInt(v)<9223372036854775807n);
const text=(max:number,multiline=false)=>z.string().refine(v=>v.trim()!==""&&new TextEncoder().encode(v).length<=max&&new TextDecoder().decode(new TextEncoder().encode(v))===v&&!(multiline?/[\u0000-\u0008\u000b\u000c\u000e-\u001f\u007f-\u009f]/u:/[\u0000-\u001f\u007f-\u009f]/u).test(v));
const base={command_id:helpID,article_id:helpID,organization_id:helpID};
const locale=z.enum(["es","en"]),category=z.string().regex(/^[a-z][a-z0-9_-]{0,63}$/u),sha=z.string().regex(/^[a-f0-9]{64}$/u);
export const helpCommand=z.discriminatedUnion("action",[
 z.object({...base,action:z.literal("create"),locale,category,title:text(200),body:text(16384,true)}).strict(),
 z.object({...base,action:z.literal("update"),version:expected,title:text(200),body:text(16384,true)}).strict(),
 z.object({...base,action:z.literal("publish"),version:expected}).strict(),
 z.object({...base,action:z.literal("archive"),version:expected}).strict()
]).refine(c=>new TextEncoder().encode(catalogCanonical(c)).length<=24576);
export type HelpCommand=z.infer<typeof helpCommand>;
const state=z.enum(["draft","published","archived"]);
export const helpArticle=z.object({id:helpID,organization_id:helpID,locale,category,title:text(200),body:text(16384,true),state,version,updated_at:z.iso.datetime({offset:true})}).strict();
export type HelpArticle=z.infer<typeof helpArticle>;
export const helpReceipt=z.object({command_id:helpID,action:z.enum(["create","update","publish","archive"]),actor:z.string().min(1).max(256),request_sha256:sha,article_sha256:sha,article:helpArticle,replay:z.boolean()}).strict();
export const helpPage=z.object({items:z.array(helpArticle.omit({organization_id:true,body:true,updated_at:true})).max(50),next:helpID.optional()}).strict();
export const pendingHelp=z.object({command_id:helpID,action:helpReceipt.shape.action,organization_id:helpID,article_id:helpID,version,state,request_sha256:sha}).strict();
export type PendingHelp=z.infer<typeof pendingHelp>;
export const helpPermission=(action:string)=>["publish","archive"].includes(action)?"help:publish":"help:write";
export async function helpReference(c:HelpCommand):Promise<PendingHelp>{return pendingHelp.parse({command_id:c.command_id,action:c.action,organization_id:c.organization_id,article_id:c.article_id,version:c.action==="create"?"1":(BigInt(c.version)+1n).toString(),state:c.action==="publish"?"published":c.action==="archive"?"archived":"draft",request_sha256:await catalogSHA(catalogCanonical(c))})}
export async function helpMatches(v:z.infer<typeof helpReceipt>,p:PendingHelp,subject:string){return v.command_id===p.command_id&&v.action===p.action&&v.actor===subject&&v.request_sha256===p.request_sha256&&v.article.id===p.article_id&&v.article.organization_id===p.organization_id&&v.article.version===p.version&&v.article.state===p.state&&v.article_sha256===await catalogSHA(catalogCanonical(v.article))}
