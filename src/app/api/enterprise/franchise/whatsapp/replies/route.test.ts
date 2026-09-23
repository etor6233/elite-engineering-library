import { beforeEach,describe,it,expect,vi } from "vitest";
import { NextRequest } from "next/server";
const mocks=vi.hoisted(()=>({session:{subject:"human",tenantId:"tenant",organizations:["org"],permissions:["whatsapp:approve","whatsapp:send"],accessToken:"private-server-only"},post:vi.fn()}));
vi.mock("@/platform/auth/session",()=>({readSession:async()=>mocks.session,allowed:(s:{permissions:string[]},p:string)=>s.permissions.includes(p)}));
vi.mock("@/platform/auth/oidc-client",()=>({applicationBaseUrl:()=>new URL("https://portal.example")}));
vi.mock("@/platform/backend/protected-client",()=>({protectedPost:mocks.post}));
import { POST } from "./route";
const id="wa-reply:"+"a".repeat(64),sha="b".repeat(64);
function request(body:unknown,origin="https://portal.example"){return new NextRequest("https://portal.example/api/enterprise/franchise/whatsapp/replies",{method:"POST",headers:{origin,"content-type":"application/json"},body:JSON.stringify(body)})}
beforeEach(()=>{mocks.post.mockReset();mocks.session.permissions=["whatsapp:approve","whatsapp:send"]});
describe("WhatsApp review BFF",()=>{
 it("rejects foreign origin before backend or effects",async()=>{expect((await POST(request({action:"send",request_id:id,payload_sha256:sha},"https://foreign.example"))).status).toBe(403);expect(mocks.post).not.toHaveBeenCalled()});
 it("requires send permission and cannot substitute editable text",async()=>{mocks.session.permissions=["whatsapp:approve"];expect((await POST(request({action:"send",request_id:id,payload_sha256:sha}))).status).toBe(403);expect((await POST(request({action:"decision",request_id:id,payload_sha256:sha,approved:true,reason:"Reviewed",text:"Replacement"}))).status).toBe(400);expect(mocks.post).not.toHaveBeenCalled()});
 it("uses only the exact stored identity/hash, and uncertainty directs a read",async()=>{mocks.post.mockRejectedValue(new Error("uncertain"));const result=await POST(request({action:"recover",request_id:id,payload_sha256:sha}));expect(result.status).toBe(409);expect(mocks.post).toHaveBeenCalledWith(mocks.session,`/v1/franchise/whatsapp/replies/${encodeURIComponent(id)}/recover`,{payload_sha256:sha});expect(await result.json()).toEqual({code:"CONSULT_STATE_BEFORE_RETRY"})});
});
