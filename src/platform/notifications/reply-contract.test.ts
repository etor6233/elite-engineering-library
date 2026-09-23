import { describe,it,expect } from "vitest";
import { replyCommand,responseText,type ReplyView } from "./reply-contract";
const id="wa-reply:"+"a".repeat(64),sha="b".repeat(64);
describe("human WhatsApp response contract",()=>{
 it("binds approval to exact stored identity/hash and excludes editable text",()=>{expect(replyCommand.safeParse({action:"decision",request_id:id,payload_sha256:sha,approved:true,reason:"Reviewed"}).success).toBe(true);expect(replyCommand.safeParse({action:"decision",request_id:id,payload_sha256:sha,approved:true,reason:"Reviewed",text:"Replacement"}).success).toBe(false)});
 it("rejects null decision and path escape",()=>{expect(replyCommand.safeParse({action:"decision",request_id:id,payload_sha256:sha,approved:null,reason:"Reviewed"}).success).toBe(false);expect(replyCommand.safeParse({action:"send",request_id:"../another",payload_sha256:sha}).success).toBe(false)});
 it("recovery has no provider ID or override payload",()=>{expect(replyCommand.safeParse({action:"recover",request_id:id,payload_sha256:sha}).success).toBe(true);expect(replyCommand.safeParse({action:"recover",request_id:id,payload_sha256:sha,provider_message_id:"forged"}).success).toBe(false)});
 it("displays exact response and rejects recipient divergence",()=>{const value={context:{message:{ExternalID:"5491112345678",Text:JSON.stringify({kind:"text_reply",recipient:"5491112345678",text:"Texto <exacto>"})}}} as ReplyView;expect(responseText(value)).toBe("Texto <exacto>");value.context.message.ExternalID="other";expect(()=>responseText(value)).toThrow("RECIPIENT_MISMATCH")});
});
