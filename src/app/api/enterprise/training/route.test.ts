import{expect,it,vi,afterEach}from"vitest";
const m=vi.hoisted(()=>({get:vi.fn(),post:vi.fn(),session:vi.fn()}));
vi.mock("@/platform/auth/session",()=>({readSession:m.session,allowed:(s:{permissions:string[]},p:string)=>s.permissions.includes(p)}));
vi.mock("@/platform/auth/oidc-client",()=>({applicationBaseUrl:()=>new URL("https://portal.example.test")}));
vi.mock("@/platform/backend/protected-client",()=>({protectedGet:m.get,protectedPost:m.post}));
vi.mock("@/platform/config/load",()=>({loadBusinessConfig:async()=>({features:{training_portal:true}})}));
import{POST}from"./route";
afterEach(()=>vi.clearAllMocks());
it("bounds decoded chunked training input before downstream and cancels the stream",async()=>{
 m.session.mockResolvedValue({subject:"learner",permissions:["training:learn"]});
 let pulls=0,cancelled=false;
 const body=new ReadableStream<Uint8Array>({pull(c){pulls++;c.enqueue(new TextEncoder().encode("é".repeat(8193)))},cancel(){cancelled=true}});
 const r=new Request("https://portal.example.test/api/enterprise/training",{method:"POST",headers:{origin:"https://portal.example.test","content-type":"application/json"},body,duplex:"half"}as RequestInit);
 const reply=await POST(r as never);expect(reply.status).toBe(413);expect(cancelled).toBe(true);expect(pulls).toBeLessThanOrEqual(5);
 expect(m.post).not.toHaveBeenCalled();expect(m.get).not.toHaveBeenCalled();
});
it("requires an authenticated role before consuming the body",async()=>{
 m.session.mockResolvedValue(null);
 const r=new Request("https://portal.example.test/api/enterprise/training",{method:"POST",headers:{origin:"https://portal.example.test","content-type":"application/json"},body:"{}"});
 expect((await POST(r as never)).status).toBe(401);expect(r.bodyUsed).toBe(false);expect(m.post).not.toHaveBeenCalled();
});
