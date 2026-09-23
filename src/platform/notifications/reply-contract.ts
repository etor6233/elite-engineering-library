import { z } from "zod";
// AUTHORED UI/API projection; the Go owner makes every approval/send decision.
const sha = z.string().regex(/^[0-9a-f]{64}$/);
export const replyCommand = z.discriminatedUnion("action", [
 z.object({action:z.literal("decision"),request_id:z.string().regex(/^wa-reply:[0-9a-f]{64}$/),payload_sha256:sha,approved:z.boolean(),reason:z.string().max(2048)}).strict(),
 z.object({action:z.enum(["send","recover"]),request_id:z.string().regex(/^wa-reply:[0-9a-f]{64}$/),payload_sha256:sha}).strict(),
]);
export type ReplyView = { request_id:string; state:"pending"|"approved"|"rejected"; payload_sha256:string; context:{organization_id:string;expires_at:string;message:{Text:string;ExternalID:string}};status?:{fence_state:string;delivery_status:string;reconciliation_required:boolean;approval_expired:boolean} };
export function responseText(value:ReplyView):string {
 const request:unknown=JSON.parse(value.context.message.Text);
 const parsed=z.object({kind:z.literal("text_reply"),recipient:z.string(),text:z.string().min(1).max(4096)}).passthrough().parse(request);
 if(parsed.recipient!==value.context.message.ExternalID)throw new Error("RECIPIENT_MISMATCH");
 return parsed.text;
}
