import { beforeEach, expect, test, vi } from "vitest";
import { renderToStaticMarkup } from "react-dom/server";
const mocks=vi.hoisted(()=>({readSession:vi.fn(),get:vi.fn()}));
vi.mock("@/platform/auth/session",()=>({readSession:mocks.readSession,allowed:(s:{permissions:string[]},p:string)=>s.permissions.includes(p)||s.permissions.includes("*")}));
vi.mock("@/platform/backend/protected-client",()=>({protectedGet:mocks.get}));
vi.mock("next/navigation",()=>({redirect:(path:string)=>{throw new Error("REDIRECT "+path)}}));
import AdminPage from "./page";
const session=(permissions:string[])=>({subject:"synthetic-admin",tenantId:"synthetic-tenant",organizations:["org-a"],permissions,accessToken:"private-fixture-token"});
beforeEach(()=>{vi.resetAllMocks();mocks.readSession.mockResolvedValue(session(["admin:read"]));mocks.get.mockImplementation(async(s:{permissions:string[]},path:string)=>{
 if(path==="/v1/franchise/leads"){
  if(!s.permissions.includes("lead:read")&&!s.permissions.includes("*"))throw new Error("FORBIDDEN lead:read");
  return {items:[{id:"lead-visible",state:"new",source_code:"fixture",assigned_subject:""}]};
 }
 if(path==="/v1/admin/overview")return {orders:2,stock_available:1,open_cases:1,active_shipments:0};
 if(path==="/v1/admin/orders")return {items:[{id:"order-visible",state:"draft",currency:"USD",total_minor_units:100}]};
 if(path==="/v1/admin/service-cases")return {items:[{id:"case-visible",state:"opened",severity:"medium"}]};
 throw new Error("UNEXPECTED PATH "+path);
});});
test.each([["admin:read"],["admin:read","Lead:Read"]].map(permissions=>({permissions})))("admin without exact lead access keeps its permitted data: $permissions",async({permissions})=>{
 mocks.readSession.mockResolvedValue(session(permissions));const html=renderToStaticMarkup(await AdminPage());
 expect(html).toContain("order-visible");expect(html).toContain("case-visible");expect(html).not.toContain("Oportunidades");expect(html).not.toContain("lead-visible");
 expect(mocks.get.mock.calls.map(c=>c[1]).sort()).toEqual(["/v1/admin/overview","/v1/admin/orders","/v1/admin/service-cases"].sort());
 expect(html).not.toContain("private-fixture-token");
});
test.each([["admin:read","lead:read"],["*"]].map(permissions=>({permissions})))("authorized lead view remains available: $permissions",async({permissions})=>{
 mocks.readSession.mockResolvedValue(session(permissions));const html=renderToStaticMarkup(await AdminPage());expect(html).toContain("Oportunidades");expect(html).toContain("lead-visible");expect(mocks.get).toHaveBeenCalledTimes(4);
 for(const call of mocks.get.mock.calls){expect(call[0].organizations).toEqual(["org-a"]);expect(call[2].organization_id).toBe("org-a");}
});
test.each([["lead:read"],["customer:self"],["factory:read"],[]].map(permissions=>({permissions})))("non-admin session cannot issue admin reads: $permissions",async({permissions})=>{
 mocks.readSession.mockResolvedValue(session(permissions));const html=renderToStaticMarkup(await AdminPage());expect(html).toContain("Acceso denegado");expect(mocks.get).not.toHaveBeenCalled();
});
test("missing session redirects before backend reads",async()=>{mocks.readSession.mockResolvedValue(null);await expect(AdminPage()).rejects.toThrow("REDIRECT /api/auth/login?return_to=/admin");expect(mocks.get).not.toHaveBeenCalled();});

 test("admin renders minor units as the actual currency amount",async()=>{const html=renderToStaticMarkup(await AdminPage()).replace(/\u00a0/g," ");expect(html).toContain("USD 1,00");expect(html).not.toContain("USD 100");});

test("customer order and quote summaries use the same verified minor-unit display",async()=>{
 const {default:CustomerPage}=await import("../customer/page");mocks.readSession.mockResolvedValue(session(["customer:self"]));
 mocks.get.mockImplementation(async(_session:unknown,path:string)=>{
  if(path==="/v1/customer/orders")return {items:[{id:"customer-order",state:"draft",currency:"USD",total_minor_units:100}]};
  if(path==="/v1/customer/service-cases")return {items:[]};
  if(path==="/v1/customer/journey")return {appointments:[],handovers:[],quotes:[{id:"customer-quote",state:"issued",currency:"KWD",total_minor_units:1234}]};
  throw new Error("UNEXPECTED CUSTOMER PATH "+path);
 });
 const html=renderToStaticMarkup(await CustomerPage()).replace(/\u00a0/g," ");expect(html).toContain("USD 1,00");expect(html).toContain("KWD 1,234");expect(html).not.toContain("USD 100");expect(html).not.toContain("KWD 1234");
});

test("admin pages expose backend cursors for all granted lists",async()=>{
 mocks.readSession.mockResolvedValue(session(["admin:read","lead:read"]));const old=mocks.get.getMockImplementation()!;
 mocks.get.mockImplementation(async(...args:unknown[])=>({...await old(...args),next_cursor:String(args[1]).endsWith("orders")?"order-next":String(args[1]).endsWith("service-cases")?"case-next":"lead-next"}));
 const html=renderToStaticMarkup(await AdminPage());
 expect(html).toContain('href="/admin?orders_after=order-next"');expect(html).toContain('href="/admin?cases_after=case-next"');expect(html).toContain('href="/admin?leads_after=lead-next"');
});

test("admin keeps all granted cursors independent and fixed session scope",async()=>{
 mocks.readSession.mockResolvedValue(session(["admin:read","lead:read"]));
 await AdminPage({searchParams:Promise.resolve({orders_after:"o",cases_after:"c",leads_after:"l",organization_id:"ignored",limit:"1000"})});
 for(const [path,after] of [["/v1/admin/orders","o"],["/v1/admin/service-cases","c"],["/v1/franchise/leads","l"]])expect(mocks.get).toHaveBeenCalledWith(expect.anything(),path,{organization_id:"org-a",limit:"25",after});
});
test("admin without lead permission neither reads nor preserves lead cursor",async()=>{
 const html=renderToStaticMarkup(await AdminPage({searchParams:Promise.resolve({orders_after:"o",leads_after:["ignored","duplicate"]})}));
 expect(html).not.toContain("leads_after");expect(html).not.toContain("Oportunidades");expect(mocks.get.mock.calls.some(c=>c[1]==="/v1/franchise/leads")).toBe(false);
});
test("admin final-page links retain the other lists",async()=>{mocks.readSession.mockResolvedValue(session(["admin:read","lead:read"]));const html=renderToStaticMarkup(await AdminPage({searchParams:Promise.resolve({orders_after:"o",cases_after:"c",leads_after:"l"})}));expect(html).toContain('href="/admin?cases_after=c&amp;leads_after=l"');expect(html).toContain('href="/admin?orders_after=o&amp;leads_after=l"');expect(html).toContain('href="/admin?orders_after=o&amp;cases_after=c"');expect(html).not.toContain("Ver siguientes");});
