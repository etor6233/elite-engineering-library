import {NextRequest} from "next/server";
import {expect,it,vi,afterEach} from "vitest";
const m=vi.hoisted(()=>({get:vi.fn(),post:vi.fn(),png:vi.fn(),session:vi.fn()}));
vi.mock("@/platform/auth/session",()=>({readSession:m.session,allowed:(s:{permissions:string[]},p:string)=>s.permissions.includes(p)}));
vi.mock("@/platform/auth/oidc-client",()=>({applicationBaseUrl:()=>new URL("https://portal.example.test")}));
vi.mock("@/platform/backend/protected-client",()=>({protectedGet:m.get,protectedPost:m.post,protectedPostPNG:m.png}));
vi.mock("@/platform/config/load",()=>({loadBusinessConfig:async()=>({features:{catalog_editor:true}})}));
import {POST,GET} from "./route";
afterEach(()=>vi.clearAllMocks());
const req=(body:BodyInit,tail="",type="application/json")=>{
 const options={method:"POST",headers:{origin:"https://portal.example.test","content-type":type},body,duplex:"half" as const};
 return new NextRequest("https://portal.example.test/api/enterprise/catalog"+tail,options);
};
it("rejects unauthenticated media before body consumption",async()=>{
 m.session.mockResolvedValue(null);const r=req(new Uint8Array([1,2,3]),"?organization_id=org&command_id=media","image/png");
 expect((await POST(r)).status).toBe(401);expect(r.bodyUsed).toBe(false);expect(m.png).not.toHaveBeenCalled();
});
it("bounds chunked binary media and cancels before downstream",async()=>{
 m.session.mockResolvedValue({subject:"maker",organizations:["org"],permissions:["catalog:draft","catalog:read"]});
 let cancelled=false,pulls=0;const stream=new ReadableStream<Uint8Array>({pull(c){pulls++;c.enqueue(new Uint8Array(262145))},cancel(){cancelled=true}});
 const r=req(stream,"?organization_id=org&command_id=media","image/png");
 expect((await POST(r)).status).toBe(413);expect(cancelled).toBe(true);expect(pulls).toBeLessThanOrEqual(5);expect(m.png).not.toHaveBeenCalled();
});
it("rejects duplicate and foreign scope before upload",async()=>{
 m.session.mockResolvedValue({subject:"maker",organizations:["org"],permissions:["catalog:draft","catalog:read"]});
 for(const tail of ["?organization_id=other&command_id=media","?organization_id=org&organization_id=org&command_id=media"]){const r=req(new Uint8Array([1]),tail,"image/png");expect((await POST(r)).status).toBe(403);expect(r.bodyUsed).toBe(false)}
 expect(m.png).not.toHaveBeenCalled();
});
it("uses private authorized current publication and rejects foreign draft projection",async()=>{
 m.session.mockResolvedValue({subject:"reader",tenantId:"50f38793-8a22-4f6b-983f-81dd0fca8208",organizations:["org"],permissions:["catalog:read"]});
 m.get.mockRejectedValue(new Error("unavailable"));
 expect((await GET(new NextRequest("https://portal.example.test/api/enterprise/catalog?kind=current&organization_id=org"))).status).toBe(409);
 expect(m.get).toHaveBeenCalledWith(expect.anything(),"/v1/admin/catalog/current",{organization_id:"org"});
 expect(m.post).not.toHaveBeenCalled();
});
