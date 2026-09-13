# TypeScript Payment Checkout Portal

## 1. Metadata

```yaml
pack_id: "TS-PAYMENT-CHECKOUT-PORTAL"
pack_version: "0.1.1"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Authenticated customer checkout URL recovery and provider-hosted redirect through the existing BFF; never treats a return URL as payment evidence."
stacks: ["Go 1.26.8", "PostgreSQL 18.6", "Node.js 24.20.0 where frontend selected"]
compatible_with: ["GO-COMMERCE-PRICING-PAYMENT-API0.6.4", "GO-ELECTROMOBILITY-APPLICATION1.10.0", "GO-OFFICIAL-PAYMENT-WEBHOOK-ADAPTERS0.3.2", "TS-GO-API-WEB-BRIDGE0.5.17 where selected"]
incompatible_with: ["unbound tenant/provider/account", "production certification inferred from fixtures", "implicit source or business-policy attribution"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://github.com/stripe/stripe-go/tree/a2df585a800a97fe8ec4ebf551b4449bdb3d90a1", "https://github.com/mercadopago/sdk-go/tree/f910ee53fbb6819e435eaf3d0f800cb1fe74ae09", "https://docs.stripe.com/api/checkout/sessions/object"]
verified_at: "2026-09-11"
```

## 2. Applicability

Use in the current composed franchise reference with the exact declared owner revisions. Missing live credentials leave the provider lane disabled; no credentials are embedded or requested. The local fixture receipts prove the declared contract only. Select the supported materialized profile explicitly; incompatible business options are rejected, not silently invented.

## 3. Architecture contract

Existing Commerce, provider inbox, durable jobs, outbound delivery fence and PostgreSQL remain the state owners. Provider callbacks are signed resource hints; account/mode/order/attempt/currency/amount facts come from the fixed SDK GET. Send acceptance, request binding and pending state share one transaction. Unknown effects are never blindly retried. Handover uses a versioned exact-byte profile and current observations under locks; tenant/org/provider/account/connection/mode are bound. No new financial ledger, statutory decision or implicit physical shipment. Profile and customer routes use existing authorization and server-side identity. Network calls are bounded and outside database locks; queue batches are finite. Roll back the caller pack set before removing its new adapters; do not delete audit rows or provider effects.

## 4. Exact file manifest

```text
CREATE src/app/api/enterprise/checkout/route.test.ts
CREATE src/app/api/enterprise/checkout/route.ts
CREATE src/components/customer-checkout-actions.tsx
CREATE src/platform/payments/checkout.ts
```

## 5. Materialization blocks

### FILE: `src/app/api/enterprise/checkout/route.test.ts`

```yaml
block_id: "TS-PAYMENT-CHECKOUT-PORTAL:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "1d9535a7f844b93ddec8c9de31f89907b7822981a68f33a8d2f5c70eccd1a7be"
variables: []
secrets_allowed: false
```

````typescript
import {beforeEach,expect,it,vi} from "vitest";
const mock=vi.hoisted(()=>({session:vi.fn(),get:vi.fn()}));
vi.mock("@/platform/auth/session",()=>({readSession:mock.session,allowed:(s:{permissions:string[]},p:string)=>s.permissions.includes(p)||s.permissions.includes("*")}));
vi.mock("@/platform/backend/protected-client",()=>({protectedGet:mock.get}));
import {GET} from "./route";
const url="https://portal.example.test/api/enterprise/checkout?organizationId=store&orderId=order";
const checkout=()=>({order_id:"order",provider_code:"stripe",url:"https://checkout.stripe.com/c/pay/cs_test_fixture",expires_at:new Date(Date.now()+3600000).toISOString()});
beforeEach(()=>{vi.clearAllMocks();mock.session.mockResolvedValue({subject:"alice",tenantId:"tenant",organizations:["store"],permissions:["customer:self"],accessToken:"server-only-token"});mock.get.mockResolvedValue(checkout());});
it("uses the session's backend token and returns a validated provider URL",async()=>{const value=checkout();mock.get.mockResolvedValue(value);const r=await GET(new Request(url));expect(r.status).toBe(200);expect(await r.json()).toEqual(value);expect(mock.get).toHaveBeenCalledWith(expect.objectContaining({subject:"alice"}),"/v1/customer/orders/order/checkout",{organization_id:"store"});expect(r.headers.get("cache-control")).toBe("no-store");});
it.each(["organizationId=store&orderId=order&customerId=bob","organizationId=store&organizationId=store&orderId=order","organizationId=store&orderId=../other","organizationId=other&orderId=order"])("rejects injected/foreign scope %s",async query=>{const response=await GET(new Request(url.split("?")[0]+"?"+query));expect([400,403]).toContain(response.status);expect(mock.get).not.toHaveBeenCalled();});
it("requires authentication",async()=>{mock.session.mockResolvedValue(null);expect((await GET(new Request(url))).status).toBe(401);expect(mock.get).not.toHaveBeenCalled();});
it("requires customer permission",async()=>{mock.session.mockResolvedValue({permissions:["admin:read"]});expect((await GET(new Request(url))).status).toBe(403);expect(mock.get).not.toHaveBeenCalled();});
it("rejects cross-site reads",async()=>{expect((await GET(new Request(url,{headers:{"sec-fetch-site":"cross-site"}}))).status).toBe(403);expect(mock.get).not.toHaveBeenCalled();});
it.each(["https://checkout.stripe.com.evil.test/session","javascript:alert(1)","https://user:secret@checkout.stripe.com/session"])("rejects a backend open redirect %s",async target=>{mock.get.mockResolvedValue({...checkout(),url:target});expect((await GET(new Request(url))).status).toBe(502);});
it("rejects another order and expired response",async()=>{mock.get.mockResolvedValue({...checkout(),order_id:"other"});expect((await GET(new Request(url))).status).toBe(502);mock.get.mockResolvedValue({...checkout(),expires_at:"2000-01-01T00:00:00Z"});expect((await GET(new Request(url))).status).toBe(502);});
it("never retries or exposes backend failure details",async()=>{mock.get.mockRejectedValue(new Error("PRIVATE_PROVIDER_DETAIL"));const r=await GET(new Request(url));expect(r.status).toBe(503);expect(await r.text()).not.toContain("PRIVATE_PROVIDER_DETAIL");expect(mock.get).toHaveBeenCalledTimes(1);});
````

### FILE: `src/app/api/enterprise/checkout/route.ts`

```yaml
block_id: "TS-PAYMENT-CHECKOUT-PORTAL:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "0dd98ce24feb982692b6dec6696c638caded8435d2c0bc3ec1b94fc94603ed6a"
variables: []
secrets_allowed: false
```

````typescript
import { NextResponse } from "next/server";
import { allowed, readSession } from "@/platform/auth/session";
import { protectedGet } from "@/platform/backend/protected-client";
import { BackendProblem } from "@/platform/backend/public-client";
import { checkoutID, checkoutSchema } from "@/platform/payments/checkout";

const reply=(value:unknown,status=200)=>NextResponse.json(value,{status,headers:{"Cache-Control":"no-store","Referrer-Policy":"no-referrer","X-Content-Type-Options":"nosniff"}});
export async function GET(request: Request) {
 if (["cross-site","none"].includes(request.headers.get("sec-fetch-site")??"")) return reply({code:"CROSS_SITE_REJECTED"},403);
 const session=await readSession();if(!session)return reply({code:"UNAUTHENTICATED"},401);
 if(!allowed(session,"customer:self"))return reply({code:"FORBIDDEN"},403);
 const query=new URL(request.url).searchParams;
 const order=query.get("orderId"),organization=query.get("organizationId");
 if(request.url.length>2048||[...query.keys()].length!==2||query.getAll("orderId").length!==1||query.getAll("organizationId").length!==1||!checkoutID.safeParse(order).success||!checkoutID.safeParse(organization).success)return reply({code:"INVALID_CHECKOUT_QUERY"},400);
 if(!session.organizations.includes(organization!)&&!session.permissions.includes("*"))return reply({code:"ORGANIZATION_FORBIDDEN"},403);
 try {
  const raw=await protectedGet<unknown>(session,"/v1/customer/orders/"+encodeURIComponent(order!)+"/checkout",{organization_id:organization!});
  const parsed=checkoutSchema.safeParse(raw);
  if(!parsed.success||parsed.data.order_id!==order)return reply({code:"CHECKOUT_UNAVAILABLE"},502);
  return reply(parsed.data);
 } catch(error) {
  if(error instanceof BackendProblem&&[401,403,404,503].includes(error.status))return reply({code:error.status===404?"CHECKOUT_NOT_AVAILABLE":"CHECKOUT_UNAVAILABLE"},error.status);
  return reply({code:"CHECKOUT_UNAVAILABLE"},503);
 }
}
````

### FILE: `src/components/customer-checkout-actions.tsx`

```yaml
block_id: "TS-PAYMENT-CHECKOUT-PORTAL:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "8e8eed64e7b5d6fd7e1ed865efdac9f300f2bda1eb2111b72dc28dfc325944b9"
variables: []
secrets_allowed: false
```

````tsx
"use client";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

import { useState } from "react";
import { checkoutSchema } from "@/platform/payments/checkout";

export function CustomerCheckoutActions({orderId,organizationId}:{orderId:string;organizationId:string}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

 const [pending,setPending]=useState(false),[message,setMessage]=useState("");
 async function openCheckout(){
  setPending(true);setMessage("");
  try{
   const query=new URLSearchParams({orderId,organizationId});
   const response=await fetch("/api/enterprise/checkout?"+query,{cache:"no-store",redirect:"error",credentials:"same-origin",headers:{accept:"application/json"},signal:AbortSignal.timeout(7000)});
   if(response.status===404){setMessage(t("p0237"));return;}
   if(!response.ok)throw new Error("checkout unavailable");
   const value=checkoutSchema.parse(await response.json());
   if(value.order_id!==orderId)throw new Error("checkout does not match order");
   window.location.assign(value.url);
  }catch{setMessage(t("p0238"));}
  finally{setPending(false);}
 }
 return <div><button type="button" disabled={pending} onClick={openCheckout}>{pending?t("p0239"):t("p0240")}</button><p role="status" aria-live="polite">{message}</p></div>;
}
````

### FILE: `src/platform/payments/checkout.ts`

```yaml
block_id: "TS-PAYMENT-CHECKOUT-PORTAL:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "512077333c00084c0fd256fc177a6c9c48f28c2b71df222cc0a5af3b5af7448f"
variables: []
secrets_allowed: false
```

````typescript
// AUTHORED redirect contract; all payment entry remains on the provider UI.
import { z } from "zod";
export const checkoutID = z.string().min(1).max(200).regex(/^[A-Za-z0-9_-]+$/);
export const checkoutSchema = z.object({
 order_id: checkoutID,
 provider_code: z.enum(["stripe", "mercadopago"]),
 url: z.string().min(1).max(8192),
 expires_at: z.iso.datetime({ offset: true }),
}).strict().superRefine((value,ctx) => {
 let target: URL;
 try { target = new URL(value.url); } catch { ctx.addIssue({ code:"custom", message:"invalid checkout URL" }); return; }
 const hosts = value.provider_code === "stripe" ? ["checkout.stripe.com"] : ["www.mercadopago.com.ar", "sandbox.mercadopago.com.ar"];
 if (target.protocol !== "https:" || target.username || target.password || target.hash || (target.port && target.port !== "443") || !hosts.includes(target.hostname) || /[\r\n]/.test(value.url) || Date.parse(value.expires_at) <= Date.now()) ctx.addIssue({code:"custom",message:"checkout is not available"});
});
````

## 6. Configuration surface

Exact future environment and profile fields are listed in the materialized docs. Defaults disable PAYMENT_CHECKOUT_ENABLED and HANDOVER_ENABLED. Enabling requires complete scope and validation before traffic. Provider secrets, webhook secret and outbound HMAC key are external runtime inputs. Hosted return URLs must be explicit trusted HTTPS configuration. Profile ID/revision/SHA and supported algorithm/options bind the same deployment scope. No secret value is present in the pack.

## 7. Dependency bill

Existing exact TS-GO-API-WEB-BRIDGE lock for Node24.20.0/Next16.3.4/React/TypeScript/Vitest. No new package or provider SDK in the browser; backend owns provider secrets and effects.

## 8. Apply order

Compose with TS-GO-API-WEB-BRIDGE and the compatible remote Go payment runtime. Preserve existing OIDC/BFF configuration and lock; no database or financial migration in this portal pack.

## 9. Verification

Run the affected Go unit/HTTP/host tests, bounded native fuzz where the selected input surface supports it, and the explicit loopback PostgreSQL connected-payment/profile/read tests with skips rejected. Run vet and build on the composed revision. Frontend selection additionally uses pinned Node/Next/Vitest/TypeScript from the admitted recipe, no unreviewed package-manager invocation. Verify wrong scope, provider mode, account, amount, URL, hash, generation, claim, duplicate callback, response loss and immutable terminal-state behavior. Compare a clean canonical reconstruction byte-for-byte before promotion. Tests of fixture contracts do not authorize live accounts or production.

## 10. Reconstruction evidence

V402: the complete77pack/915file reference was reconstructed into an absent external destination; every output matched the frozen source inventory. PostgreSQL57migrations and16 connected tests passed without skips, plus vet/build and the explicitly enumerated unit/host cases. Finite fuzz and frontend contract receipts are separately bound to identical files. See reconstruction_evidence/CONNECTED_PAYMENT_HANDOVER_V402.md and CONNECTED_DELTA_ASSURANCE_V402.md. Current SAST/SCA, whole-library operations/performance/portable release and remaining capabilities retain their own gates; no live production or complete TEST02 claim.

Canonical V402 integration: selected by the current profile with exact dependencies and caller overlays. Metadata promotion records byte reconstruction, not closure of every admission/release gate. Payload provenance is unchanged.

V402 composed delta: T2804 private es/en display, per-user/tenant preference and browser negotiation;1198messages,21exact guides,5hash-bound curricula; original commands/content/policies preserved. PRIVATE_LOCALE_RELEASE_V402.md/json. AUTHORED glue; no new dependency or corporate attribution.
