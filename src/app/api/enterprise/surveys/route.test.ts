import { beforeEach, describe, expect, it, vi } from "vitest";
const mock = vi.hoisted(() => ({ session:vi.fn(), config:vi.fn(), get:vi.fn(), post:vi.fn() }));
vi.mock("@/platform/auth/session",()=>({readSession:mock.session,allowed:(s:{permissions:string[]},p:string)=>s.permissions.includes(p)||s.permissions.includes("*")}));
vi.mock("@/platform/auth/oidc-client",()=>({applicationBaseUrl:()=>new URL("https://portal.example.test")}));
vi.mock("@/platform/config/load",()=>({loadBusinessConfig:mock.config}));
vi.mock("@/platform/backend/protected-client",()=>({protectedGet:mock.get,protectedPost:mock.post}));
vi.mock("@/platform/backend/public-client",()=>({BackendProblem:class extends Error {status=503;code="SURVEY_UNAVAILABLE";}}));
import { GET, POST } from "./route";
const input = {score:0,consent:true,consent_version:"v1"};
const answer = {score:0,consent_version:"v1",received_at:"2035-01-01T12:00:00Z"};
const url="https://portal.example.test/api/enterprise/surveys?organizationId=org-a&surveyId=survey";
const request = (body=JSON.stringify(input),headers:Record<string,string>={},address=url) => new Request(address,{method:"POST",headers:{"origin":"https://portal.example.test","content-type":"application/json",...headers},body});
beforeEach(()=>{
 vi.clearAllMocks();
 mock.config.mockResolvedValue({features:{customer_surveys:true}});
 mock.session.mockResolvedValue({subject:"alice",tenantId:"tenant",organizations:["org-a"],permissions:["customer:self"],accessToken:"server-only"});
 mock.post.mockResolvedValue({answer,replay:false});
});
describe("survey BFF boundary",()=>{
 it("uses the session scope and exact zero score; returns no token",async()=>{
  const response=await POST(request());expect(response.status).toBe(201);
  expect(await response.json()).toEqual({answer,replay:false});
  expect(mock.post).toHaveBeenCalledWith(expect.objectContaining({subject:"alice"}),"/v1/customer/surveys/survey/response?organization_id=org-a",input);
  expect(response.headers.get("cache-control")).toBe("no-store");
 });
 it("keeps a replay a successful read of the original receipt",async()=>{
  mock.post.mockResolvedValue({answer,replay:true});expect((await POST(request())).status).toBe(200);
 });
 it("does not activate by possession of a session",async()=>{
  mock.config.mockResolvedValue({features:{}});expect((await POST(request())).status).toBe(404);expect(mock.session).not.toHaveBeenCalled();expect(mock.post).not.toHaveBeenCalled();
 });
 it.each([undefined,"https://evil.example.test","null"])("rejects wrong or missing origin %s",async origin=>{
  const r=request();if(origin===undefined)r.headers.delete("origin");else r.headers.set("origin",origin);
  expect((await POST(r)).status).toBe(403);expect(mock.post).not.toHaveBeenCalled();
 });
 it("rejects a cross-site read",async()=>{
  expect((await GET(new Request(url+"&view=response",{headers:{"sec-fetch-site":"cross-site"}}))).status).toBe(403);expect(mock.get).not.toHaveBeenCalled();
 });
 it("requires a session",async()=>{
  mock.session.mockResolvedValue(null);expect((await POST(request())).status).toBe(401);expect(mock.post).not.toHaveBeenCalled();
 });
 it.each(["organizationId=org-b&surveyId=survey","organizationId=org-a&surveyId=survey&organizationId=org-a","organizationId=org-a&surveyId=survey&customerId=bob"])("rejects foreign/ambiguous query %s",async query=>{
  const response=await POST(request(undefined,{},url.split("?")[0]+"?"+query));expect([400,403]).toContain(response.status);expect(mock.post).not.toHaveBeenCalled();
 });
 it("requires separate aggregate permission",async()=>{
  expect((await GET(new Request(url+"&view=summary"))).status).toBe(403);expect(mock.get).not.toHaveBeenCalled();
 });
 it.each([
  '{"score":0,"score":9,"consent":true,"consent_version":"v1"}',
  '{"score":0,"consent":false,"consent_version":"v1"}',
  '{"score":0,"consent":true,"consent_version":"v1","customer_id":"bob"}',
  " ".repeat(2049)
 ])("rejects body without making a write: %s",async body=>{
  expect((await POST(request(body))).status).toBe(400);expect(mock.post).not.toHaveBeenCalled();
 });
 it("rejects encoded input",async()=>{
  expect((await POST(request(undefined,{"content-encoding":"gzip"}))).status).toBe(415);expect(mock.post).not.toHaveBeenCalled();
 });
 it("never retries a failed write",async()=>{
  mock.post.mockRejectedValue(new Error("opaque upstream detail"));
  const response=await POST(request());expect(response.status).toBe(503);expect(await response.text()).not.toContain("opaque");expect(mock.post).toHaveBeenCalledTimes(1);
 });
 it("rejects a mismatched successful receipt",async()=>{
  mock.post.mockResolvedValue({answer:{...answer,score:9},replay:false});expect((await POST(request())).status).toBe(502);expect(mock.post).toHaveBeenCalledTimes(1);
 });
 it("rejects a mismatched definition",async()=>{
  mock.get.mockResolvedValue({id:"other",prompt:"Q",consent_version:"v1",consent_notice:"C",accepting:true});
  expect((await GET(new Request(url+"&view=definition"))).status).toBe(502);
 });
 it("permits authorized aggregate reads without a customer grant",async()=>{
  mock.session.mockResolvedValue({organizations:["org-a"],permissions:["surveys:read"]});
  mock.get.mockResolvedValue({responses:1,available:false,nps:null});
  const response=await GET(new Request(url+"&view=summary"));expect(response.status).toBe(200);expect(await response.json()).toEqual({responses:1,available:false,nps:null});
 });
});
