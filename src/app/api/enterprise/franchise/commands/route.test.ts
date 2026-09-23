vi.mock("@/platform/experience/options",()=>({validateGuidedSelection:vi.fn().mockResolvedValue(true)}));
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { NextRequest } from "next/server";

const mocks = vi.hoisted(() => ({ readSession: vi.fn(), protectedPost: vi.fn(), protectedGet: vi.fn() }));
vi.mock("@/platform/auth/session", () => ({ readSession: mocks.readSession, allowed: (session: { permissions: string[] }, permission: string) => session.permissions.includes(permission) }));
vi.mock("@/platform/backend/protected-client", () => ({ protectedPost: mocks.protectedPost, protectedGet:mocks.protectedGet }));

import { GET, POST } from "./route";

const session = { subject: "operator", tenantId: "tenant", organizations: ["store"], permissions: ["lead:assign", "lead:update", "quote:write", "appointment:manage", "resource:manage", "availability:manage", "handover:manage", "customer:self"], accessToken: "token" };

function request(body: unknown, origin = "https://portal.example") {
  return new NextRequest("https://portal.example/api/enterprise/franchise/commands", { method: "POST", headers: { origin, "content-type": "application/json" }, body: JSON.stringify(body) });
}

describe("franchise command BFF", () => {
  beforeEach(() => { vi.clearAllMocks(); vi.stubEnv("APP_BASE_URL","https://portal.example"); mocks.readSession.mockResolvedValue(session); mocks.protectedPost.mockResolvedValue({ version: 2 }); });
  afterEach(() => vi.unstubAllEnvs());

  it("recovers a published checklist by its natural scoped version without POST",async()=>{
    const url="https://portal.example/api/enterprise/franchise/commands?kind=checklist&organizationId=store&checklistId=fixture-checklist&version=1";mocks.protectedGet.mockResolvedValue({id:"fixture-checklist",version:1});
    expect((await GET(new NextRequest(url))).status).toBe(200);expect(mocks.protectedGet).toHaveBeenCalledWith(session,"/v1/franchise/delivery-checklists/result",{organization_id:"store",checklist_id:"fixture-checklist",version:"1"});expect(mocks.protectedPost).not.toHaveBeenCalled();
    mocks.protectedGet.mockClear();for(const bad of [url+"&requestKey=unknown-key",url.replace("version=1","version=0"),url.replace("version=1","version=9007199254740993")])expect((await GET(new NextRequest(bad))).status).toBe(400);expect(mocks.protectedGet).not.toHaveBeenCalled();
    mocks.readSession.mockResolvedValue({...session,permissions:["resource:manage"]});expect((await GET(new NextRequest(url))).status).toBe(403);expect(mocks.protectedGet).not.toHaveBeenCalled();
  });

  it("retains slot identity in POST and uses scoped GET recovery",async()=>{
    const command={action:"create-appointment-slot",organizationId:"store",kind:"service",startsAt:"2030-01-01T10:00:00Z",endsAt:"2030-01-01T11:00:00Z",capacity:2};
    expect((await POST(request(command))).status).toBe(400);expect(mocks.protectedPost).not.toHaveBeenCalled();
    const keyed=request(command);keyed.headers.set("idempotency-key","slot-fixture-key");expect((await POST(keyed)).status).toBe(201);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session,"/v1/franchise/appointment-slots",{organization_id:"store",kind:"service",starts_at:command.startsAt,ends_at:command.endsAt,capacity:2},"slot-fixture-key");
    mocks.protectedPost.mockClear();mocks.protectedGet.mockResolvedValue({id:"resource"});
    const url="https://portal.example/api/enterprise/franchise/commands?kind=slot&organizationId=store&requestKey=slot-fixture-key";
    expect((await GET(new NextRequest(url))).status).toBe(200);expect(mocks.protectedGet).toHaveBeenCalledWith(session,"/v1/franchise/appointment-slots/result",{organization_id:"store",request_key:"slot-fixture-key"});expect(mocks.protectedPost).not.toHaveBeenCalled();
    mocks.protectedGet.mockClear();mocks.readSession.mockResolvedValue({...session,permissions:["availability:manage"]});expect((await GET(new NextRequest(url))).status).toBe(403);expect(mocks.protectedGet).not.toHaveBeenCalled();
  });
  it("retains resource identity in POST and uses scoped GET recovery",async()=>{
    const command={action:"create-resource",organizationId:"store",displayName:"Synthetic bay",kind:"service-bay",skills:["service"]};
    expect((await POST(request(command))).status).toBe(400);expect(mocks.protectedPost).not.toHaveBeenCalled();
    const keyed=request(command);keyed.headers.set("idempotency-key","resource-fixture-key");expect((await POST(keyed)).status).toBe(201);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session,"/v1/franchise/resources",{organization_id:"store",principal_subject:"",display_name:"Synthetic bay",kind:"service-bay",skills:["service"]},"resource-fixture-key");
    mocks.protectedPost.mockClear();mocks.protectedGet.mockResolvedValue({id:"resource"});
    const url="https://portal.example/api/enterprise/franchise/commands?kind=resource&organizationId=store&requestKey=resource-fixture-key";
    expect((await GET(new NextRequest(url))).status).toBe(200);expect(mocks.protectedGet).toHaveBeenCalledWith(session,"/v1/franchise/resources/result",{organization_id:"store",request_key:"resource-fixture-key"});expect(mocks.protectedPost).not.toHaveBeenCalled();
    mocks.protectedGet.mockClear();mocks.readSession.mockResolvedValue({...session,permissions:["availability:manage"]});expect((await GET(new NextRequest(url))).status).toBe(403);expect(mocks.protectedGet).not.toHaveBeenCalled();
  });

  it("maps an initial payment with a stable key but without browser-owned money/provider",async()=>{
    const operator={...session,permissions:["payment:create"]};mocks.readSession.mockResolvedValue(operator);
    const command={action:"request-order-payment",organizationId:"store",orderId:"order",requestKey:"11111111-1111-4111-8111-111111111111"};
    expect((await POST(request(command))).status).toBe(200);
    expect(mocks.protectedPost).toHaveBeenCalledWith(operator,"/v1/commerce/orders/order/payment-request",{organization_id:"store",request_key:command.requestKey});
    vi.clearAllMocks();mocks.readSession.mockResolvedValue(operator);
    for(const bad of [{...command,provider:"stripe"},{...command,amount:1},{...command,requestKey:"invalid"},{...command,actor:"owner"}])expect((await POST(request(bad))).status).toBe(400);
    expect((await POST(request(command,"https://other.invalid"))).status).toBe(403);
    expect((await POST(request({...command,organizationId:"other"}))).status).toBe(403);
    mocks.readSession.mockResolvedValue({...operator,permissions:["admin:read"]});expect((await POST(request(command))).status).toBe(403);
    expect(mocks.protectedPost).not.toHaveBeenCalled();
  });

  it("maps an allocation to the existing stock owner without server-owned fields", async()=>{
    const operator={...session,permissions:["inventory:allocate"]};mocks.readSession.mockResolvedValue(operator);
    const command={action:"allocate-order-stock",organizationId:"store",orderId:"order",lineId:"line",stockUnitId:"unit",orderVersion:3,stockVersion:2};
    expect((await POST(request(command))).status).toBe(200);
    expect(mocks.protectedPost).toHaveBeenCalledWith(operator,"/v1/commerce/orders/order/allocations",{organization_id:"store",line_id:"line",stock_unit_id:"unit",order_version:3,stock_version:2});
    vi.clearAllMocks();
    for(const invalid of [{...command,state:"reserved"},{...command,orderVersion:0},{...command,stockVersion:1.5},{...command,stockUnitId:""}]) expect((await POST(request(invalid))).status).toBe(400);
    expect(mocks.protectedPost).not.toHaveBeenCalled();
  });
  it("does not admit an allocation for a read-only or unrelated organization session",async()=>{
    const command={action:"allocate-order-stock",organizationId:"store",orderId:"order",lineId:"line",stockUnitId:"unit",orderVersion:1,stockVersion:1};
    mocks.readSession.mockResolvedValue({...session,permissions:["admin:read"]});
    expect((await POST(request(command))).status).toBe(403);
    mocks.readSession.mockResolvedValue({...session,organizations:["other"],permissions:["inventory:allocate"]});
    expect((await POST(request(command))).status).toBe(403);
    expect(mocks.protectedPost).not.toHaveBeenCalled();
  });

  it("uses the configured public origin behind an internal proxy", async () => {
    const input=new NextRequest("http://localhost:4173/api/enterprise/franchise/commands",{method:"POST",headers:{origin:"https://portal.example","content-type":"application/json",host:"localhost:4173"},body:JSON.stringify({action:"assign-appointment-resource",organizationId:"store",appointmentId:"a",resourceId:"r",version:1})});
    expect((await POST(input)).status).toBe(200);
  });

  it("does not trust forwarded host or an absent configured public origin", async () => {
    const input=request({action:"assign-lead"},"https://evil.example");
    input.headers.set("x-forwarded-host","evil.example");
    expect((await POST(input)).status).toBe(403);
    vi.stubEnv("APP_BASE_URL","");
    expect((await POST(request({action:"assign-lead"}))).status).toBe(403);
    expect(mocks.readSession).not.toHaveBeenCalled();
  });

  it("rejects cross-origin before session or backend", async () => {
    const response = await POST(request({ action: "assign-lead" }, "https://evil.example"));
    expect(response.status).toBe(403);
    expect(mocks.readSession).not.toHaveBeenCalled();
  });

  it("maps a scoped assignment to the existing Go command", async () => {
    const response = await POST(request({ action: "assign-lead", organizationId: "store", leadId: "lead-1", assignedSubject: "seller", version: 1 }));
    expect(response.status).toBe(200);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/franchise/leads/lead-1/assign", { organization_id: "store", assigned_subject: "seller", version: 1 });
  });

  it("rejects organization and permission escalation", async () => {
    expect((await POST(request({ action: "assign-lead", organizationId: "other", leadId: "lead-1", assignedSubject: "seller", version: 1 }))).status).toBe(403);
    mocks.readSession.mockResolvedValue({ ...session, permissions: [] });
    expect((await POST(request({ action: "assign-lead", organizationId: "store", leadId: "lead-1", assignedSubject: "seller", version: 1 }))).status).toBe(403);
    expect(mocks.protectedPost).not.toHaveBeenCalled();
  });

  it("maps a quote without accepting client price or currency", async () => {
    const req=request({ action: "create-quote", organizationId: "store", leadId: "lead-1", variantId: "variant", priceBookId: "retail", validUntil: "2026-09-01T12:00:00.000Z" });
    req.headers.set("idempotency-key","quote-request-key-0001");
    const response = await POST(req);
    expect(response.status).toBe(201);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/franchise/quotes", { organization_id: "store", lead_id: "lead-1", variant_id: "variant", price_book_id: "retail", valid_until: "2026-09-01T12:00:00.000Z" }, "quote-request-key-0001");
  });

  it("maps bounded appointment capacity to the existing Go owner", async () => {
    const slotRequest = request({ action: "create-appointment-slot", organizationId: "store", kind: "service", startsAt: "2026-09-01T12:00:00.000Z", endsAt: "2026-09-01T13:00:00.000Z", capacity: 2 });
slotRequest.headers.set("idempotency-key","slot-positive-fixture-key");const response=await POST(slotRequest);
    expect(response.status).toBe(201);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/franchise/appointment-slots", { organization_id: "store", kind: "service", starts_at: "2026-09-01T12:00:00.000Z", ends_at: "2026-09-01T13:00:00.000Z", capacity: 2 },"slot-positive-fixture-key");
  });

  it("maps availability and resource commands without accepting server-owned fields", async () => {
    const availabilityRequest=request({ action: "create-availability", organizationId: "store", resourceId: "technician", entryType: "unavailable", reasonCode: "annual-leave", startsAt: "2026-09-01T12:00:00.000Z", endsAt: "2026-09-01T13:00:00.000Z" });
    availabilityRequest.headers.set("idempotency-key","availability-request-0001");
    const availability=await POST(availabilityRequest);
    expect(availability.status).toBe(201);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/franchise/availability", { organization_id: "store", resource_id: "technician", entry_type: "unavailable", reason_code: "annual-leave", starts_at: "2026-09-01T12:00:00.000Z", ends_at: "2026-09-01T13:00:00.000Z" }, "availability-request-0001");
    vi.clearAllMocks();
    const injected = await POST(request({ action: "create-availability", organizationId: "store", entryType: "working", startsAt: "2026-09-01T12:00:00.000Z", endsAt: "2026-09-01T13:00:00.000Z", state: "active", version: 99 }));
    expect(injected.status).toBe(400);
    expect(mocks.protectedPost).not.toHaveBeenCalled();
  });

  it("maps audited appointment transition and customer-owned cancellation", async () => {
    const transition = await POST(request({ action: "transition-appointment", organizationId: "store", appointmentId: "appointment-1", current: "confirmed", target: "no-show", version: 3, reasonCode: "customer-absent" }));
    expect(transition.status).toBe(200);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/franchise/appointments/appointment-1/transitions", { organization_id: "store", current: "confirmed", target: "no-show", version: 3, reason_code: "customer-absent" });
    vi.clearAllMocks();
    const cancellation = await POST(request({ action: "cancel-customer-appointment", organizationId: "store", appointmentId: "appointment-2", version: 1, reasonCode: "customer-request" }));
    expect(cancellation.status).toBe(200);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/customer/appointments/appointment-2/cancel", { organization_id: "store", version: 1, reason_code: "customer-request" });
  });

  it("maps customer quote acceptance without client price, currency or evidence digest", async () => {
    const response = await POST(request({ action: "accept-quote", organizationId: "store", quoteId: "quote-1", version: 1 }));
    expect(response.status).toBe(200);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/customer/quotes/quote-1/accept", { organization_id: "store", version: 1 });
    vi.clearAllMocks();
    const injected = await POST(request({ action: "accept-quote", organizationId: "store", quoteId: "quote-1", version: 1, totalMinorUnits: 1, currency: "USD" }));
    expect(injected.status).toBe(400);
    expect(mocks.protectedPost).not.toHaveBeenCalled();
  });

  it("maps handover acceptance without accepting a browser-authored evidence digest", async () => {
    const response = await POST(request({ action: "accept-handover", organizationId: "store", handoverId: "handover-1", version: 2, confirmedReceived: true, serialNumber: "SERIAL-1", checklistId: "standard-delivery", checklistVersion: 1 }));
    expect(response.status).toBe(200);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/customer/handovers/handover-1/accept", { organization_id: "store", version: 2, confirmed_received: true, serial_number: "SERIAL-1", checklist_id: "standard-delivery", checklist_version: 1 });
    vi.clearAllMocks();
    const injected = await POST(request({ action: "accept-handover", organizationId: "store", handoverId: "handover-1", version: 2, confirmedReceived: true, serialNumber: "SERIAL-1", checklistId: "standard-delivery", checklistVersion: 1, evidenceSha256: "a".repeat(64) }));
    expect(injected.status).toBe(400);
    expect(mocks.protectedPost).not.toHaveBeenCalled();
  });

  it("publishes and completes an exact delivery checklist version", async () => {
    const items = [{ id: "serial-observed", prompt: "Verificar serie", response_type: "serial", required: true }];
    const published = await POST(request({ action: "publish-delivery-checklist", organizationId: "store", checklistId: "standard-delivery", version: 1, title: "Entrega estándar", items }));
    expect(published.status).toBe(201);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/franchise/delivery-checklists", { organization_id: "store", checklist_id: "standard-delivery", version: 1, title: "Entrega estándar", items });
    vi.clearAllMocks();
    mocks.readSession.mockResolvedValue(session);
    mocks.protectedPost.mockResolvedValue({ version: 2 });
    const responses = [{ item_id: "serial-observed", response_text: "SERIAL-1" }];
    const completed = await POST(request({ action: "complete-delivery-checklist", organizationId: "store", handoverId: "handover-1", version: 1, checklistId: "standard-delivery", checklistVersion: 1, responses }));
    expect(completed.status).toBe(200);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/franchise/handovers/handover-1/complete-checklist", { organization_id: "store", version: 1, checklist_id: "standard-delivery", checklist_version: 1, responses });
  });

  it("maps customer rejection and operator resolution without browser authority fields", async () => {
    const rejected = await POST(request({ action: "reject-handover", organizationId: "store", handoverId: "handover-1", version: 2, reasonCode: "visible-damage", details: "Rayón visible" }));
    expect(rejected.status).toBe(200);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/customer/handovers/handover-1/reject", { organization_id: "store", version: 2, reason_code: "visible-damage", details: "Rayón visible" });
    vi.clearAllMocks();
    mocks.readSession.mockResolvedValue(session);
    mocks.protectedPost.mockResolvedValue({ exception: { state: "resolved" } });
    const resolved = await POST(request({ action: "resolve-delivery-exception", organizationId: "store", exceptionId: "exception-1", version: 1, resolutionAction: "correct-and-represent", notes: "Corregir preparación" }));
    expect(resolved.status).toBe(200);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/franchise/delivery-exceptions/exception-1/resolve", { organization_id: "store", version: 1, action: "correct-and-represent", notes: "Corregir preparación" });
    vi.clearAllMocks();
    const injected = await POST(request({ action: "resolve-delivery-exception", organizationId: "store", exceptionId: "exception-1", version: 1, resolutionAction: "refund-now", notes: "invalid", stockState: "available" }));
    expect(injected.status).toBe(400);
    expect(mocks.protectedPost).not.toHaveBeenCalled();
  });

  it("maps physical return receipt and disposition without browser-owned effects", async () => {
    const received = await POST(request({ action: "receive-return", organizationId: "store", authorizationId: "authorization-1", serialNumber: "SERIAL-1", conditionCode: "damaged", notes: "Daño confirmado" }));
    expect(received.status).toBe(200);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/franchise/return-authorizations/authorization-1/receive", { organization_id: "store", serial_number: "SERIAL-1", condition_code: "damaged", notes: "Daño confirmado", evidence_sha256: expect.stringMatching(/^[0-9a-f]{64}$/) });
    vi.clearAllMocks();
    mocks.readSession.mockResolvedValue(session);
    mocks.protectedPost.mockResolvedValue({ customer_remedy: "refund" });
    const decided = await POST(request({ action: "decide-return", organizationId: "store", receiptId: "receipt-1", inventoryAction: "quarantine", notes: "Separar y solicitar efectos" }));
    expect(decided.status).toBe(200);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/franchise/return-receipts/receipt-1/decide", { organization_id: "store", inventory_action: "quarantine", notes: "Separar y solicitar efectos" });
    vi.clearAllMocks();
    const injected = await POST(request({ action: "decide-return", organizationId: "store", receiptId: "receipt-1", inventoryAction: "restock", notes: "invalid", refundNow: true, paymentId: "browser-owned" }));
    expect(injected.status).toBe(400);
    expect(mocks.protectedPost).not.toHaveBeenCalled();
  });
});

describe("quote recovery boundary",()=>{
 beforeEach(()=>{vi.clearAllMocks();vi.stubEnv("APP_BASE_URL","https://portal.example");mocks.readSession.mockResolvedValue(session)});
 afterEach(()=>vi.unstubAllEnvs());
 it("rejects absent and malformed quote keys before the backend",async()=>{
  for(const key of [null,"short","a".repeat(129),"has invalid spaces"]){
   const req=request({action:"create-quote",organizationId:"store",leadId:"lead",variantId:"variant",priceBookId:"book",validUntil:"2026-09-20T12:00:00Z"});if(key!==null)req.headers.set("idempotency-key",key);
   expect((await POST(req)).status).toBe(400);
  }expect(mocks.protectedPost).not.toHaveBeenCalled();
 });
 it("uses scoped GET without posting and conserves unavailable results",async()=>{
  const req=new NextRequest("https://portal.example/api/enterprise/franchise/commands?organizationId=store&leadId=lead&requestKey=quote-request-key-0001");
  mocks.protectedGet.mockResolvedValue({id:"q"});const response=await GET(req);expect(response.status).toBe(200);expect(response.headers.get("cache-control")).toBe("no-store");
  expect(mocks.protectedGet).toHaveBeenCalledWith(session,"/v1/franchise/quotes/result",{organization_id:"store",lead_id:"lead",request_key:"quote-request-key-0001"});
  mocks.protectedGet.mockRejectedValue(new Error("synthetic private detail"));const failed=await GET(req);expect(failed.status).toBe(502);expect(await failed.text()).not.toContain("private");
  mocks.readSession.mockResolvedValue({...session,permissions:["lead:read"]});expect((await GET(req)).status).toBe(403);
  mocks.readSession.mockResolvedValue(null);expect((await GET(req)).status).toBe(401);expect(mocks.protectedPost).not.toHaveBeenCalled();
 });
});

describe("availability creation recovery BFF",()=>{
 beforeEach(()=>{vi.clearAllMocks();vi.stubEnv("APP_BASE_URL","https://portal.example");mocks.readSession.mockResolvedValue(session)});afterEach(()=>vi.unstubAllEnvs());
 it("reads an exact availability reference with manage permission and no write",async()=>{
  const req=new NextRequest("https://portal.example/api/enterprise/franchise/commands?kind=availability&organizationId=store&requestKey=availability-request-0001");
  mocks.protectedGet.mockResolvedValue({id:"entry"});expect((await GET(req)).status).toBe(200);
  expect(mocks.protectedGet).toHaveBeenCalledWith(session,"/v1/franchise/availability/result",{organization_id:"store",request_key:"availability-request-0001"});
  mocks.readSession.mockResolvedValue({...session,permissions:["availability:read"]});expect((await GET(req)).status).toBe(403);expect(mocks.protectedPost).not.toHaveBeenCalled();
 });
 it("rejects missing creation identity before forwarding",async()=>{
  const req=request({action:"create-availability",organizationId:"store",entryType:"working",startsAt:"2026-09-20T12:00:00Z",endsAt:"2026-09-20T13:00:00Z"});
  expect((await POST(req)).status).toBe(400);expect(mocks.protectedPost).not.toHaveBeenCalled();
 });
});


describe("checklist completion recovery boundary",()=>{
 beforeEach(()=>{vi.clearAllMocks();mocks.readSession.mockResolvedValue({...session,permissions:["handover:manage"]})});
 it("reads exact completion without posting and rejects scope or extra fields",async()=>{
  const url="https://portal.example/api/enterprise/franchise/commands?kind=checklist-completion&organizationId=store&handoverId=handover";
  mocks.protectedGet.mockResolvedValue({handover_id:"handover",state:"presented"});const result=await GET(new NextRequest(url));expect(result.status).toBe(200);expect(result.headers.get("cache-control")).toBe("no-store");
  expect(mocks.protectedGet).toHaveBeenCalledWith({...session,permissions:["handover:manage"]},"/v1/franchise/handovers/handover/checklist-result",{organization_id:"store"});
  expect((await GET(new NextRequest(url+"&version=1"))).status).toBe(400);expect((await GET(new NextRequest(url.replace("store","other")))).status).toBe(403);
  mocks.protectedGet.mockRejectedValue(new Error("private synthetic detail"));const failed=await GET(new NextRequest(url));expect(failed.status).toBe(502);expect(await failed.text()).not.toContain("private");
  mocks.readSession.mockResolvedValue({...session,permissions:["admin:read"]});expect((await GET(new NextRequest(url))).status).toBe(403);expect(mocks.protectedPost).not.toHaveBeenCalled();
 });
});


describe("return recovery boundary",()=>{
 beforeEach(()=>{vi.clearAllMocks();vi.stubEnv("APP_BASE_URL","https://portal.example")});afterEach(()=>vi.unstubAllEnvs());
 it("uses an exact authorization with management scope, no writes and redacted failures",async()=>{
  const caller={...session,permissions:["handover:manage"]};mocks.readSession.mockResolvedValue(caller);mocks.protectedGet.mockResolvedValue({authorization_id:"authorization"});
  const url="https://portal.example/api/enterprise/franchise/commands?kind=return&organizationId=store&authorizationId=authorization";
  const result=await GET(new NextRequest(url));expect(result.status).toBe(200);expect(result.headers.get("cache-control")).toBe("no-store");expect(mocks.protectedGet).toHaveBeenCalledWith(caller,"/v1/franchise/returns/result",{organization_id:"store",authorization_id:"authorization"});
  expect((await GET(new NextRequest(url+"&receiptId=receipt"))).status).toBe(400);expect((await GET(new NextRequest(url.replace("&authorizationId=authorization","")))).status).toBe(400);expect((await GET(new NextRequest(url.replace("organizationId=store","organizationId=other")))).status).toBe(403);
  mocks.protectedGet.mockRejectedValue(new Error("synthetic private failure"));const failed=await GET(new NextRequest(url));expect(failed.status).toBe(502);expect(await failed.text()).not.toContain("private");
  mocks.readSession.mockResolvedValue({...session,permissions:["admin:read"]});expect((await GET(new NextRequest(url))).status).toBe(403);mocks.readSession.mockResolvedValue(null);expect((await GET(new NextRequest(url))).status).toBe(401);expect(mocks.protectedPost).not.toHaveBeenCalled();
 });
});
