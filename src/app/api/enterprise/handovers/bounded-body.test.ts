import{expect,it,vi}from"vitest";
const m=vi.hoisted(()=>({get:vi.fn(),post:vi.fn()}));
vi.mock("@/platform/auth/session",()=>({readSession:async()=>({subject:"operator",tenantId:"tenant",organizations:["store"],permissions:["handover:manage"]}),allowed:()=>true}));
vi.mock("@/platform/auth/oidc-client",()=>({applicationBaseUrl:()=>new URL("https://portal.example.test")}));
vi.mock("@/platform/backend/protected-client",()=>({protectedGet:m.get,protectedPost:m.post}));
import{POST}from"./route";
it("cuts off an oversized chunked UTF-8 request without Content-Length before downstream",async()=>{
 let pulls=0,cancelled=false;
 const body=new ReadableStream<Uint8Array>({pull(controller){pulls++;controller.enqueue(new TextEncoder().encode("é".repeat(1025)))},cancel(){cancelled=true}});
 const r=new Request("https://portal.example.test/api/enterprise/handovers",{method:"POST",headers:{origin:"https://portal.example.test","content-type":"application/json"},body,duplex:"half"}as RequestInit);
 expect(r.headers.get("content-length")).toBeNull();const reply=await POST(r);
 expect(reply.status).toBe(413);expect(await reply.json()).toEqual({code:"COMMAND_TOO_LARGE"});expect(cancelled).toBe(true);expect(pulls).toBeLessThanOrEqual(3);expect(m.get).not.toHaveBeenCalled();expect(m.post).not.toHaveBeenCalled();
});
