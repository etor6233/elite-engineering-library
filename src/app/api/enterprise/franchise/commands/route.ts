import {isQuoteAcceptanceReceipt} from "@/platform/payments/quote-receipt";
import {validateGuidedSelection} from "@/platform/experience/options";
import { NextResponse, type NextRequest } from "next/server";
import { z } from "zod";
import { allowed, readSession } from "@/platform/auth/session";
import { applicationBaseUrl } from "@/platform/auth/oidc-client";
import { boundedCommandBody } from "@/platform/backend/bounded-command";
import { BackendProblem } from "@/platform/backend/public-client";
import { protectedGet, protectedPost } from "@/platform/backend/protected-client";

const identifier = z.string().min(1).max(128);
const organization = z.string().min(1).max(128);
const version = z.number().int().positive();
const reasonCode = z.string().regex(/^[a-z][a-z0-9]*(-[a-z0-9]+)*$/);
const checklistItem = z.object({ id: identifier, prompt: z.string().trim().min(1).max(500), response_type: z.enum(["confirmation", "text", "serial", "evidence"]), required: z.boolean() }).strict();
const checklistResponse = z.object({ item_id: identifier, response_text: z.string().trim().min(1).max(2048), evidence_sha256: z.string().regex(/^[0-9a-f]{64}$/).optional() }).strict();

const commandSchema = z.discriminatedUnion("action", [
  z.object({action:z.literal("request-order-payment"),organizationId:organization,orderId:identifier,requestKey:z.uuid()}).strict(),
  z.object({action:z.literal("allocate-order-stock"),organizationId:organization,orderId:identifier,lineId:identifier,stockUnitId:identifier,orderVersion:version,stockVersion:version}).strict(),
  z.object({ action: z.literal("assign-lead"), organizationId: organization, leadId: identifier, assignedSubject: z.string().min(1).max(256), version }).strict(),
  z.object({ action: z.literal("transition-lead"), organizationId: organization, leadId: identifier, current: z.enum(["new", "contacted", "qualified"]), target: z.enum(["contacted", "qualified", "converted", "lost"]), version }).strict(),
  z.object({ action: z.literal("create-quote"), organizationId: organization, leadId: identifier, variantId: identifier, priceBookId: identifier, validUntil: z.iso.datetime({ offset: true }) }).strict(),
  z.object({ action: z.literal("create-appointment-slot"), organizationId: organization, kind: z.enum(["consultation", "test-drive", "delivery", "service"]), startsAt: z.iso.datetime({ offset: true }), endsAt: z.iso.datetime({ offset: true }), capacity: z.number().int().min(1).max(100) }).strict(),
  z.object({ action: z.literal("create-resource"), organizationId: organization, principalSubject: z.string().max(256).optional(), displayName: z.string().min(1).max(256), kind: z.enum(["employee", "contractor", "service-bay", "vehicle"]), skills: z.array(z.enum(["consultation", "test-drive", "delivery", "service"])).min(1).max(16) }).strict(),
  z.object({ action: z.literal("create-availability"), organizationId: organization, resourceId: identifier.optional(), entryType: z.enum(["working", "unavailable"]), reasonCode: reasonCode.optional(), startsAt: z.iso.datetime({ offset: true }), endsAt: z.iso.datetime({ offset: true }) }).strict(),
  z.object({ action: z.literal("cancel-availability"), organizationId: organization, availabilityId: identifier, version, reasonCode }).strict(),
  z.object({ action: z.literal("assign-appointment-resource"), organizationId: organization, appointmentId: identifier, resourceId: identifier, version }).strict(),
  z.object({ action: z.literal("transition-appointment"), organizationId: organization, appointmentId: identifier, current: z.enum(["requested", "confirmed"]), target: z.enum(["confirmed", "completed", "cancelled", "no-show"]), version, reasonCode: reasonCode.optional() }).strict(),
  z.object({ action: z.literal("cancel-customer-appointment"), organizationId: organization, appointmentId: identifier, version, reasonCode }).strict(),
  z.object({ action: z.literal("accept-quote"), organizationId: organization, quoteId: identifier, version }).strict(),
  z.object({ action: z.literal("publish-delivery-checklist"), organizationId: organization, checklistId: identifier, version, title: z.string().trim().min(1).max(160), items: z.array(checklistItem).min(1).max(64) }).strict(),
  z.object({ action: z.literal("complete-delivery-checklist"), organizationId: organization, handoverId: identifier, version, checklistId: identifier, checklistVersion: version, responses: z.array(checklistResponse).min(1).max(64) }).strict(),
  z.object({ action: z.literal("accept-handover"), organizationId: organization, handoverId: identifier, version, confirmedReceived: z.literal(true), serialNumber: z.string().min(1).max(128), checklistId: identifier, checklistVersion: version }).strict(),
  z.object({ action: z.literal("reject-handover"), organizationId: organization, handoverId: identifier, version, reasonCode, details: z.string().trim().min(1).max(1000) }).strict(),
  z.object({ action: z.literal("resolve-delivery-exception"), organizationId: organization, exceptionId: identifier, version, resolutionAction: z.enum(["correct-and-represent", "return", "exchange"]), notes: z.string().trim().min(1).max(1000) }).strict(),
  z.object({ action: z.literal("receive-return"), organizationId: organization, authorizationId: identifier, serialNumber: z.string().min(1).max(128), conditionCode: z.enum(["sealed", "opened", "damaged", "incomplete"]), notes: z.string().trim().min(1).max(1000) }).strict(),
  z.object({ action: z.literal("decide-return"), organizationId: organization, receiptId: identifier, inventoryAction: z.enum(["quarantine", "restock", "repair", "scrap"]), notes: z.string().trim().min(1).max(1000) }).strict()
]);

function sameOrigin(request: NextRequest) {
  const supplied = request.headers.get("origin");
  if (!supplied) return false;
  try { return supplied === applicationBaseUrl().origin; } catch { return false; }
}

export async function POST(request: NextRequest) {
  if (!sameOrigin(request)) return NextResponse.json({ code: "CROSS_ORIGIN_REJECTED" }, { status: 403 });
  const session = await readSession();
  if (!session) return NextResponse.json({ code: "UNAUTHENTICATED" }, { status: 401 });
  if (!["appointment:manage", "availability:manage", "customer:self", "handover:manage", "inventory:allocate", "lead:assign", "lead:update", "payment:create", "quote:write", "resource:manage"].some(permission => allowed(session,permission))) return NextResponse.json({ code: "FORBIDDEN" }, { status: 403 });
  if (request.headers.get("content-type") !== "application/json") return NextResponse.json({ code: "UNSUPPORTED_MEDIA_TYPE" }, { status: 415 });
  const length = Number(request.headers.get("content-length") ?? "0");
  if (!Number.isSafeInteger(length) || length < 0 || length > 65_536) return NextResponse.json({ code: "COMMAND_TOO_LARGE" }, { status: 413 });
  let text: string;
  try { text = await boundedCommandBody(request,65_536); } catch(e) { return NextResponse.json({code:e instanceof Error && e.message==="BODY_TOO_LARGE"?"COMMAND_TOO_LARGE":"INVALID_BODY"},{status:e instanceof Error && e.message==="BODY_TOO_LARGE"?413:400}); }
  let json: unknown;
  try { json = JSON.parse(text); } catch { return NextResponse.json({ code: "INVALID_JSON" }, { status: 400 }); }
  const parsed = commandSchema.safeParse(json);
  if (!parsed.success) return NextResponse.json({ code: "INVALID_COMMAND" }, { status: 400 });
  const command = parsed.data;
  if (!session.organizations.includes(command.organizationId)) return NextResponse.json({ code: "ORGANIZATION_FORBIDDEN" }, { status: 403 });

  let permission: string;
  let path: string;
  let body: Record<string, unknown>;
  switch (command.action) {
    case "request-order-payment":
      permission="payment:create";
      path=`/v1/commerce/orders/${encodeURIComponent(command.orderId)}/payment-request`;
      body={organization_id:command.organizationId,request_key:command.requestKey};
      break;
    case "allocate-order-stock":
      permission="inventory:allocate";
      path=`/v1/commerce/orders/${encodeURIComponent(command.orderId)}/allocations`;
      body={organization_id:command.organizationId,line_id:command.lineId,stock_unit_id:command.stockUnitId,order_version:command.orderVersion,stock_version:command.stockVersion};
      break;
    case "assign-lead":
      permission = "lead:assign";
      path = `/v1/franchise/leads/${encodeURIComponent(command.leadId)}/assign`;
      body = { organization_id: command.organizationId, assigned_subject: command.assignedSubject, version: command.version };
      break;
    case "transition-lead":
      permission = "lead:update";
      path = `/v1/franchise/leads/${encodeURIComponent(command.leadId)}/transitions`;
      body = { organization_id: command.organizationId, current: command.current, target: command.target, version: command.version };
      break;
    case "create-quote":
      permission = "quote:write";
      path = "/v1/franchise/quotes";
      body = { organization_id: command.organizationId, lead_id: command.leadId, variant_id: command.variantId, price_book_id: command.priceBookId, valid_until: command.validUntil };
      break;
    case "create-appointment-slot":
      permission = "appointment:manage";
      path = "/v1/franchise/appointment-slots";
      body = { organization_id: command.organizationId, kind: command.kind, starts_at: command.startsAt, ends_at: command.endsAt, capacity: command.capacity };
      break;
    case "create-resource":
      permission = "resource:manage";
      path = "/v1/franchise/resources";
      body = { organization_id: command.organizationId, principal_subject: command.principalSubject ?? "", display_name: command.displayName, kind: command.kind, skills: command.skills };
      break;
    case "create-availability":
      permission = "availability:manage";
      path = "/v1/franchise/availability";
      body = { organization_id: command.organizationId, resource_id: command.resourceId ?? "", entry_type: command.entryType, reason_code: command.reasonCode ?? "", starts_at: command.startsAt, ends_at: command.endsAt };
      break;
    case "cancel-availability":
      permission = "availability:manage";
      path = `/v1/franchise/availability/${encodeURIComponent(command.availabilityId)}/cancel`;
      body = { organization_id: command.organizationId, version: command.version, reason_code: command.reasonCode };
      break;
    case "assign-appointment-resource":
      permission = "appointment:manage";
      path = `/v1/franchise/appointments/${encodeURIComponent(command.appointmentId)}/resources`;
      body = { organization_id: command.organizationId, resource_id: command.resourceId, version: command.version };
      break;
    case "transition-appointment":
      permission = "appointment:manage";
      path = `/v1/franchise/appointments/${encodeURIComponent(command.appointmentId)}/transitions`;
      body = { organization_id: command.organizationId, current: command.current, target: command.target, version: command.version, reason_code: command.reasonCode ?? "" };
      break;
    case "cancel-customer-appointment":
      permission = "customer:self";
      path = `/v1/customer/appointments/${encodeURIComponent(command.appointmentId)}/cancel`;
      body = { organization_id: command.organizationId, version: command.version, reason_code: command.reasonCode };
      break;
    case "accept-quote":
      permission = "customer:self";
      path = `/v1/customer/quotes/${encodeURIComponent(command.quoteId)}/accept`;
      body = { organization_id: command.organizationId, version: command.version };
      break;
    case "publish-delivery-checklist":
      permission = "handover:manage";
      path = "/v1/franchise/delivery-checklists";
      body = { organization_id: command.organizationId, checklist_id: command.checklistId, version: command.version, title: command.title, items: command.items };
      break;
    case "complete-delivery-checklist":
      permission = "handover:manage";
      path = `/v1/franchise/handovers/${encodeURIComponent(command.handoverId)}/complete-checklist`;
      body = { organization_id: command.organizationId, version: command.version, checklist_id: command.checklistId, checklist_version: command.checklistVersion, responses: command.responses };
      break;
    case "accept-handover":
      permission = "customer:self";
      path = `/v1/customer/handovers/${encodeURIComponent(command.handoverId)}/accept`;
      body = { organization_id: command.organizationId, version: command.version, confirmed_received: command.confirmedReceived, serial_number: command.serialNumber, checklist_id: command.checklistId, checklist_version: command.checklistVersion };
      break;
    case "reject-handover":
      permission = "customer:self";
      path = `/v1/customer/handovers/${encodeURIComponent(command.handoverId)}/reject`;
      body = { organization_id: command.organizationId, version: command.version, reason_code: command.reasonCode, details: command.details };
      break;
    case "resolve-delivery-exception":
      permission = "handover:manage";
      path = `/v1/franchise/delivery-exceptions/${encodeURIComponent(command.exceptionId)}/resolve`;
      body = { organization_id: command.organizationId, version: command.version, action: command.resolutionAction, notes: command.notes };
      break;
    case "receive-return": {
      permission = "handover:manage";
      path = `/v1/franchise/return-authorizations/${encodeURIComponent(command.authorizationId)}/receive`;
      const digest = await crypto.subtle.digest("SHA-256", new TextEncoder().encode(JSON.stringify(command)));
      const evidence = Array.from(new Uint8Array(digest), (byte) => byte.toString(16).padStart(2, "0")).join("");
      body = { organization_id: command.organizationId, serial_number: command.serialNumber, condition_code: command.conditionCode, notes: command.notes, evidence_sha256: evidence };
      break;
    }
    case "decide-return":
      permission = "handover:manage";
      path = `/v1/franchise/return-receipts/${encodeURIComponent(command.receiptId)}/decide`;
      body = { organization_id: command.organizationId, inventory_action: command.inventoryAction, notes: command.notes };
      break;
  }
  if (!allowed(session, permission)) return NextResponse.json({ code: "FORBIDDEN" }, { status: 403 });
  try { if(!await validateGuidedSelection(command,session))return NextResponse.json({code:"SELECTION_NOT_AUTHORIZED",operation_effect:"NOT_ATTEMPTED"},{status:403}); }
  catch { return NextResponse.json({code:"SELECTION_CONFIGURATION_UNAVAILABLE",operation_effect:"NOT_ATTEMPTED"},{status:503}); }

  const requestKey = request.headers.get("idempotency-key");
  if (["create-quote","create-availability","create-resource","create-appointment-slot"].includes(command.action) && (!requestKey || !/^[A-Za-z0-9_-]{16,128}$/.test(requestKey))) return NextResponse.json({code:"INVALID_IDEMPOTENCY_KEY"},{status:400});
  try {
    const value = ["create-quote","create-availability","create-resource","create-appointment-slot"].includes(command.action) ? await protectedPost<unknown>(session, path, body, requestKey!) : await protectedPost<unknown>(session, path, body);
    if (command.action === "accept-quote" && !isQuoteAcceptanceReceipt(value, command.organizationId, {id:command.quoteId,version:command.version})) return NextResponse.json({code:"QUOTE_ACCEPTANCE_UNCONFIRMED"},{status:502,headers:{"cache-control":"no-store"}});
    return NextResponse.json(value, { status: ["create-quote", "create-appointment-slot", "create-resource", "create-availability", "publish-delivery-checklist"].includes(command.action) ? 201 : 200 });
  } catch (error) {
    if (error instanceof BackendProblem) return NextResponse.json({ code: error.code }, { status: error.status });
    return NextResponse.json({ code: "UPSTREAM_UNAVAILABLE" }, { status: 502 });
  }
}

export async function GET(request: NextRequest) {
  const session = await readSession();
  if (!session) return NextResponse.json({code:"UNAUTHENTICATED"},{status:401});
  const parameters=Object.fromEntries(new URL(request.url).searchParams);
  if(parameters.kind==="return"){
    const parsed=z.object({kind:z.literal("return"),organizationId:organization,authorizationId:identifier}).strict().safeParse(parameters);
    if(!parsed.success)return NextResponse.json({code:"INVALID_QUERY"},{status:400});
    const q=parsed.data;if(!session.organizations.includes(q.organizationId)||!allowed(session,"handover:manage"))return NextResponse.json({code:"FORBIDDEN"},{status:403});
    try{
      const value=await protectedGet<unknown>(session,"/v1/franchise/returns/result",{organization_id:q.organizationId,authorization_id:q.authorizationId});
      return NextResponse.json({returnCase:value},{headers:{"cache-control":"no-store"}});
    }catch(error){
      if(error instanceof BackendProblem)return NextResponse.json({code:error.code},{status:error.status,headers:{"cache-control":"no-store"}});
      return NextResponse.json({code:"UPSTREAM_UNAVAILABLE"},{status:502,headers:{"cache-control":"no-store"}});
    }
  }
  if(parameters.kind==="checklist-completion"){
    const parsed=z.object({kind:z.literal("checklist-completion"),organizationId:organization,handoverId:identifier}).strict().safeParse(parameters);
    if(!parsed.success)return NextResponse.json({code:"INVALID_QUERY"},{status:400});
    const q=parsed.data;if(!session.organizations.includes(q.organizationId)||!allowed(session,"handover:manage"))return NextResponse.json({code:"FORBIDDEN"},{status:403});
    try{
      const value=await protectedGet<unknown>(session,`/v1/franchise/handovers/${encodeURIComponent(q.handoverId)}/checklist-result`,{organization_id:q.organizationId});
      return NextResponse.json({completion:value},{headers:{"cache-control":"no-store"}});
    }catch(error){
      if(error instanceof BackendProblem)return NextResponse.json({code:error.code},{status:error.status,headers:{"cache-control":"no-store"}});
      return NextResponse.json({code:"UPSTREAM_UNAVAILABLE"},{status:502,headers:{"cache-control":"no-store"}});
    }
  }
  if(parameters.kind==="checklist"){
    const parsed=z.object({kind:z.literal("checklist"),organizationId:organization,checklistId:z.string().max(64).regex(/^[a-z][a-z0-9]*(-[a-z0-9]+)*$/),version:z.string().regex(/^[1-9][0-9]{0,15}$/).transform(Number).refine(Number.isSafeInteger)}).strict().safeParse(parameters);
    if(!parsed.success)return NextResponse.json({code:"INVALID_QUERY"},{status:400});
    const q=parsed.data;if(!session.organizations.includes(q.organizationId)||!allowed(session,"handover:manage"))return NextResponse.json({code:"FORBIDDEN"},{status:403});
    try{
      const value=await protectedGet<unknown>(session,"/v1/franchise/delivery-checklists/result",{organization_id:q.organizationId,checklist_id:q.checklistId,version:String(q.version)});
      return NextResponse.json({checklist:value},{headers:{"cache-control":"no-store"}});
    }catch(error){
      if(error instanceof BackendProblem)return NextResponse.json({code:error.code},{status:error.status,headers:{"cache-control":"no-store"}});
      return NextResponse.json({code:"UPSTREAM_UNAVAILABLE"},{status:502,headers:{"cache-control":"no-store"}});
    }
  }
  const query = z.object({organizationId:organization,leadId:identifier.optional(),kind:z.enum(["quote","availability","resource","slot"]).default("quote"),requestKey:z.string().regex(/^[A-Za-z0-9_-]{16,128}$/)}).strict().refine(v=>v.kind==="quote"?!!v.leadId:v.leadId===undefined).safeParse(Object.fromEntries(new URL(request.url).searchParams));
  if (!query.success) return NextResponse.json({code:"INVALID_QUERY"},{status:400});
  const {organizationId,leadId,requestKey,kind}=query.data;
  if (!session.organizations.includes(organizationId) || !allowed(session,kind==="quote"?"quote:write":kind==="availability"?"availability:manage":kind==="resource"?"resource:manage":"appointment:manage")) return NextResponse.json({code:"FORBIDDEN"},{status:403});
  try {
    const value=await protectedGet<unknown>(session,kind==="quote"?"/v1/franchise/quotes/result":kind==="availability"?"/v1/franchise/availability/result":kind==="resource"?"/v1/franchise/resources/result":"/v1/franchise/appointment-slots/result",{organization_id:organizationId,...(kind==="quote"?{lead_id:leadId}:{}),request_key:requestKey});
    return NextResponse.json({request_key:requestKey,[kind]:value},{headers:{"cache-control":"no-store"}});
  } catch(error) {
    if(error instanceof BackendProblem)return NextResponse.json({code:error.code},{status:error.status,headers:{"cache-control":"no-store"}});
    return NextResponse.json({code:"UPSTREAM_UNAVAILABLE"},{status:502,headers:{"cache-control":"no-store"}});
  }
}
