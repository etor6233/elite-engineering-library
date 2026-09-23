import{afterEach,expect,it,vi}from"vitest";
import{protectedGet}from"./protected-client";
const original=process.env.ENTERPRISE_API_BASE_URL;
afterEach(()=>{process.env.ENTERPRISE_API_BASE_URL=original;vi.unstubAllGlobals()});
it("forwards the recovery key in a header through the existing bounded client",async()=>{
 process.env.ENTERPRISE_API_BASE_URL="https://api.example.test";
 const fetch=vi.fn(async(_url:unknown,_init:RequestInit)=>new Response("{}",{headers:{"content-type":"application/json"}}));vi.stubGlobal("fetch",fetch);
 const s={subject:"operator",tenantId:"tenant",permissions:["handover:manage"],organizations:["store"],accessToken:"fixture-server-only"};
 await protectedGet(s,"/v1/franchise/orders/order/handover-result",{organization_id:"store"},"a03c7850-13ad-4b38-81d0-120000000001");
 expect(fetch.mock.calls[0]?.[1]).toMatchObject({cache:"no-store",redirect:"error",headers:{"Idempotency-Key":"a03c7850-13ad-4b38-81d0-120000000001",authorization:"Bearer fixture-server-only"}});
 expect(String(fetch.mock.calls[0]?.[0])).not.toContain("a03c7850");
 await expect(protectedGet(s,"/v1/franchise/orders/order/handover-result",{},"injected\r\nvalue")).rejects.toThrow();expect(fetch).toHaveBeenCalledTimes(1);
});
