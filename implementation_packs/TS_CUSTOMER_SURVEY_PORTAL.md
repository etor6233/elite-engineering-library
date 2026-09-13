# TypeScript Customer Survey Portal

## 1. Metadata

```yaml
pack_id: "TS-CUSTOMER-SURVEY-PORTAL"
pack_version: "0.1.2"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Connected customer survey form, strict BFF, GET-only recovery and separately authorized NPS view"
stacks: ["Next.js 16.3.4", "React 19.2.8", "TypeScript 7.0.2"]
compatible_with: ["GO-CUSTOMER-SURVEY-API 0.1.1", "TS-GO-API-WEB-BRIDGE 0.5.16", "MICROSOFT-PLAYWRIGHT-BROWSER-GATE 0.1.37"]
incompatible_with: ["automatic campaign or legal-consent certification", "unscoped response access", "production admission inferred from synthetic fixtures"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://www.netpromotersystem.com/about/measuring-your-net-promoter-score/", "https://www.postgresql.org/docs/18/explicit-locking.html", "https://go.dev/doc/security/fuzz/"]
verified_at: "2026-09-11"
```

## 2. Applicability

Opt-in connected survey reference. Project-owned definition, notice, recommendation question,
retention, threshold and authorized distribution of invitation links are required.
The existing GO-SURVEYS-CORE remains an unchanged isolated reference. No upstream
enterprise implementation is claimed. See docs/customer-surveys.md from the Go pack.

## 3. Architecture contract

One answer per tenant/organization/survey/customer. OIDC identity and active CRM
membership govern writes; customer:self and surveys:read are separate grants.
PostgreSQL owns commit time, conflict/replay and expiry. UI performs no automatic
POST retry and recovers by GET. No raw errors, credentials or foreign identifiers
are returned. Count is visible to an authorized aggregate reader below the NPS
threshold. Synthetic tests do not supply production grants or business policy.

## 4. Exact file manifest

```text
CREATE src/platform/surveys/contract.ts
CREATE src/platform/surveys/contract.test.ts
CREATE src/app/api/enterprise/surveys/route.ts
CREATE src/app/api/enterprise/surveys/route.test.ts
CREATE src/components/customer-survey.tsx
CREATE src/app/customer/surveys/page.tsx
CREATE src/app/admin/surveys/page.tsx
CREATE microsoft_playwright_browser_gate/tests/customer-survey.spec.mjs
```

## 5. Materialization blocks

### FILE: `src/platform/surveys/contract.ts`
```yaml
block_id: "TS-CUSTOMER-SURVEY-PORTAL:src/platform/surveys/contract.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local; existing admitted reference contracts, no external business source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "eba644ae21718d6033387ba27f96e9c4fda5c37ebcddd299f252a5e414351d2f"
variables: []
secrets_allowed: false
```
````typescript
import { z } from "zod";
export const surveyID = z.string().regex(/^[A-Za-z0-9][A-Za-z0-9_-]{0,127}$/);
const version = z.string().regex(/^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$/);
export const definitionSchema = z.object({ id: surveyID, prompt: z.string().min(1).max(1000), consent_version: version, consent_notice: z.string().min(1).max(4000), accepting: z.boolean() }).strict();
export const answerSchema = z.object({ score: z.number().int().min(0).max(10), consent_version: version, received_at: z.iso.datetime({ offset: true }) }).strict();
export const resultSchema = z.object({ answer: answerSchema, replay: z.boolean() }).strict();
export const submissionSchema = z.object({ score: z.number().int().min(0).max(10), consent: z.literal(true), consent_version: version }).strict();
export const summarySchema = z.object({ responses: z.number().int().min(0).max(Number.MAX_SAFE_INTEGER), available: z.boolean(), nps: z.number().min(-100).max(100).nullable() }).strict().refine(v => v.available ? v.responses > 0 && v.nps !== null : v.nps === null);
export type SurveyDefinition = z.infer<typeof definitionSchema>;
export type SurveyAnswer = z.infer<typeof answerSchema>;
export type SurveySubmission = z.infer<typeof submissionSchema>;

// Flat, bounded request grammar preserves duplicate-key rejection through the BFF.
export function parseSubmission(text: string): SurveySubmission {
 if (text.length > 2048) throw new Error("invalid survey body");
 let at = 0;
 const ws = () => { while (at < text.length && /[ \t\r\n]/.test(text[at]!)) at++; };
 const expect = (c: string) => { ws(); if (text[at++] !== c) throw new Error("invalid survey body"); };
 const string = () => {
  ws(); const token = /^"(?:[^"\\\u0000-\u001f]|\\(?:["\\/bfnrt]|u[0-9a-fA-F]{4}))*"/.exec(text.slice(at));
  if (!token) throw new Error("invalid survey body"); at += token[0].length; return JSON.parse(token[0]) as string;
 };
 const values: Record<string, unknown> = Object.create(null);
 expect("{"); ws();
 while (text[at] !== "}") {
  const key = string(); if (!["score", "consent", "consent_version"].includes(key) || Object.hasOwn(values, key)) throw new Error("invalid survey body");
  expect(":"); ws();
  if (key === "consent_version") values[key] = string();
  else {
   const token = (key === "score" ? /^-?(?:0|[1-9][0-9]*)/ : /^(?:true|false)/).exec(text.slice(at));
   if (!token) throw new Error("invalid survey body"); at += token[0].length; values[key] = JSON.parse(token[0]) as unknown;
  }
  ws(); if (text[at] === "}") break;
  expect(","); ws(); if (text[at] === "}") throw new Error("invalid survey body");
 }
 expect("}"); ws(); if (at !== text.length) throw new Error("invalid survey body");
 return submissionSchema.parse(values);
}

````

### FILE: `src/platform/surveys/contract.test.ts`
```yaml
block_id: "TS-CUSTOMER-SURVEY-PORTAL:src/platform/surveys/contract.test.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local; existing admitted reference contracts, no external business source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "cdc6631a3cf220477bcedaff9326beda591f37e2a1a5e54d7334c02f6ce6c1d9"
variables: []
secrets_allowed: false
```
````typescript
import { describe, expect, it } from "vitest";
import { parseSubmission, summarySchema } from "./contract";
describe("survey submission grammar", () => {
 it("preserves zero and explicit consent", () => {
  expect(parseSubmission('{"score":0,"consent":true,"consent_version":"v1"}')).toEqual({score:0,consent:true,consent_version:"v1"});
 });
 it.each([
  '{}', '[]', '{"score":0,"consent":true}', '{"score":null,"consent":true,"consent_version":"v1"}',
  '{"score":1e0,"consent":true,"consent_version":"v1"}', '{"score":1.5,"consent":true,"consent_version":"v1"}',
  '{"score":11,"consent":true,"consent_version":"v1"}', '{"score":0,"consent":false,"consent_version":"v1"}',
  '{"score":0,"score":9,"consent":true,"consent_version":"v1"}',
  '{"score":0,"scor\\u0065":9,"consent":true,"consent_version":"v1"}',
  '{"score":0,"consent":true,"consent_version":"v1","customer_id":"other"}',
  '{"score":0,"consent":true,"consent_version":"v1",}',
  '{"score":0,"consent":true,"consent_version":"v1"}{}',
  '{"__proto__":{},"score":0,"consent":true,"consent_version":"v1"}',
  ' '.repeat(2049)
 ])("rejects ambiguous or invalid body: %s", text => { expect(() => parseSubmission(text)).toThrow(); });
 it("does not show a score below the configured threshold", () => {
  expect(summarySchema.safeParse({responses:1,available:false,nps:null}).success).toBe(true);
  expect(summarySchema.safeParse({responses:1,available:false,nps:100}).success).toBe(false);
  expect(summarySchema.safeParse({responses:0,available:true,nps:0}).success).toBe(false);
 });
});
````

### FILE: `src/app/api/enterprise/surveys/route.ts`
```yaml
block_id: "TS-CUSTOMER-SURVEY-PORTAL:src/app/api/enterprise/surveys/route.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local; existing admitted reference contracts, no external business source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "3070da7d8a67af77bc3d8050edc12b16f28a2a0fee3a76edf67db0938afc9941"
variables: []
secrets_allowed: false
```
````typescript
import { NextResponse } from "next/server";
import { allowed, readSession } from "@/platform/auth/session";
import { applicationBaseUrl } from "@/platform/auth/oidc-client";
import { loadBusinessConfig } from "@/platform/config/load";
import { protectedGet, protectedPost } from "@/platform/backend/protected-client";
import { BackendProblem } from "@/platform/backend/public-client";
import { answerSchema, definitionSchema, resultSchema, summarySchema, surveyID, parseSubmission } from "@/platform/surveys/contract";

const reply = (value: unknown, status = 200) => NextResponse.json(value, { status, headers: { "Cache-Control": "no-store", "X-Content-Type-Options": "nosniff" } });
async function scope(request: Request, writing: boolean) {
 try { if ((await loadBusinessConfig()).features.customer_surveys !== true) return reply({ code: "CAPABILITY_NOT_ENABLED" }, 404); }
 catch { return reply({ code: "SURVEY_UNAVAILABLE" }, 503); }
 if (["cross-site", "none"].includes(request.headers.get("sec-fetch-site") ?? "")) return reply({ code: "CROSS_SITE_REJECTED" }, 403);
 if (writing) {
  try { if (request.headers.get("origin") !== applicationBaseUrl().origin) return reply({ code: "CROSS_ORIGIN_REJECTED" }, 403); }
  catch { return reply({ code: "CROSS_ORIGIN_REJECTED" }, 403); }
 }
 const session = await readSession(); if (!session) return reply({ code: "UNAUTHENTICATED" }, 401);
 const query = new URL(request.url).searchParams;
 const organization = query.get("organizationId"), id = query.get("surveyId"), view = query.get("view");
 const keys = writing ? ["organizationId", "surveyId"] : ["organizationId", "surveyId", "view"];
 if (request.url.length > 2048 || [...query.keys()].length !== keys.length || keys.some(key => query.getAll(key).length !== 1) ||
     !surveyID.safeParse(organization).success || !surveyID.safeParse(id).success || (!writing && !["definition", "response", "summary"].includes(view ?? ""))) return reply({ code: "INVALID_QUERY" }, 400);
 if (!allowed(session, view === "summary" ? "surveys:read" : "customer:self")) return reply({ code: "FORBIDDEN" }, 403);
 if (!session.organizations.includes(organization!) && !session.organizations.includes("*")) return reply({ code: "ORGANIZATION_FORBIDDEN" }, 403);
 return { session, organization: organization!, id: id!, view };
}
function failure(error: unknown) {
 if (error instanceof BackendProblem && [400, 401, 403, 404, 409, 503].includes(error.status) &&
     ["INVALID_SURVEY_REQUEST", "SURVEY_RESPONSE_CONFLICT", "SURVEY_NOT_AVAILABLE", "SURVEY_UNAVAILABLE", "UNAUTHENTICATED", "FORBIDDEN", "ORGANIZATION_FORBIDDEN"].includes(error.code)) return reply({ code: error.code }, error.status);
 return reply({ code: "SURVEY_UNAVAILABLE" }, 503);
}
export async function GET(request: Request) {
 const value = await scope(request, false); if (value instanceof Response) return value;
 const { session, organization, id, view } = value;
 const path = view === "summary" ? "/v1/admin/surveys/" + encodeURIComponent(id) + "/summary" :
  "/v1/customer/surveys/" + encodeURIComponent(id) + (view === "response" ? "/response" : "");
 try {
  const raw = await protectedGet<unknown>(session, path, { organization_id: organization });
  const parsed = (view === "summary" ? summarySchema : view === "response" ? answerSchema : definitionSchema).safeParse(raw);
  if (!parsed.success || (view === "definition" && (!("id" in parsed.data) || parsed.data.id !== id))) return reply({ code: "INVALID_SURVEY_RESPONSE" }, 502);
  return reply(parsed.data);
 } catch (error) { return failure(error); }
}
async function boundedBody(request: Request) {
 if (!request.body) throw new Error("invalid survey body");
 const reader = request.body.getReader(); const chunks: Uint8Array[] = []; let bytes = 0; let timer: ReturnType<typeof setTimeout> | undefined;
 const deadline = new Promise<never>((_, reject) => { timer = setTimeout(() => reject(new Error("body deadline")), 2500); });
 try {
  while (true) {
   const part = await Promise.race([reader.read(), deadline]); if (part.done) break;
   bytes += part.value.byteLength; if (bytes > 2048) throw new Error("body limit"); chunks.push(part.value);
  }
  const data = new Uint8Array(bytes); let offset = 0; for (const chunk of chunks) { data.set(chunk, offset); offset += chunk.byteLength; }
  return new TextDecoder("utf-8", { fatal: true }).decode(data);
 } finally { clearTimeout(timer); void reader.cancel().catch(() => {}); reader.releaseLock(); }
}
export async function POST(request: Request) {
 const value = await scope(request, true); if (value instanceof Response) return value;
 const type = (request.headers.get("content-type") ?? "").split(";", 1)[0]?.trim().toLowerCase();
 if (type !== "application/json" || !["", "identity"].includes(request.headers.get("content-encoding") ?? "")) return reply({ code: "INVALID_CONTENT_TYPE" }, 415);
 let input;
 try { input = parseSubmission(await boundedBody(request)); } catch { return reply({ code: "INVALID_SURVEY_REQUEST" }, 400); }
 const { session, organization, id } = value;
 try {
  const raw = await protectedPost<unknown>(session, "/v1/customer/surveys/" + encodeURIComponent(id) + "/response?organization_id=" + encodeURIComponent(organization), input);
  const parsed = resultSchema.safeParse(raw);
  if (!parsed.success || parsed.data.answer.score !== input.score || parsed.data.answer.consent_version !== input.consent_version) return reply({ code: "INVALID_SURVEY_RESPONSE" }, 502);
  return reply(parsed.data, parsed.data.replay ? 200 : 201);
 } catch (error) { return failure(error); }
}

````

### FILE: `src/app/api/enterprise/surveys/route.test.ts`
```yaml
block_id: "TS-CUSTOMER-SURVEY-PORTAL:src/app/api/enterprise/surveys/route.test.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "local; existing admitted reference contracts, no external business source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "f28c869694c0672d515ee78d34d819d3e7ec884c5e193f3b7d58b994f1e42db5"
variables: []
secrets_allowed: false
```
````typescript
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
````

### FILE: `src/components/customer-survey.tsx`
```yaml
block_id: "TS-CUSTOMER-SURVEY-PORTAL:src/components/customer-survey.tsx:v1"
operation: CREATE
provenance: AUTHORED
source: "local; existing admitted reference contracts, no external business source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "bb90357422c9207258ca04119e22973f37ba9ae5d97e466f2f50113fc973a7e7"
variables: []
secrets_allowed: false
```
````tsx
"use client";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

import { useRef, useState, type FormEvent } from "react";
import { answerSchema, resultSchema, type SurveyAnswer, type SurveyDefinition, type SurveySubmission } from "@/platform/surveys/contract";

export function CustomerSurvey({ organization, definition, initialAnswer }: { organization: string; definition: SurveyDefinition; initialAnswer: SurveyAnswer | null }) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

 const [answer, setAnswer] = useState(initialAnswer); const [score, setScore] = useState(""); const [consent, setConsent] = useState(false);
 const [busy, setBusy] = useState(false); const [uncertain, setUncertain] = useState(false); const [message, setMessage] = useState("");
 const pending = useRef(false); const attempted = useRef<SurveySubmission | null>(null);
 const url = "/api/enterprise/surveys?organizationId=" + encodeURIComponent(organization) + "&surveyId=" + encodeURIComponent(definition.id);
 async function submit(event: FormEvent) {
  event.preventDefault(); if (pending.current || answer || uncertain || !consent || score === "" || !definition.accepting) return;
  const input: SurveySubmission = { score: Number(score), consent: true, consent_version: definition.consent_version };
  attempted.current = input; pending.current = true; setBusy(true); setMessage("");
  try {
   const response = await fetch(url, { method: "POST", headers: { "content-type": "application/json" }, body: JSON.stringify(input), signal: AbortSignal.timeout(7000) });
   if (!response.ok) throw new Error("confirmation unavailable");
   const value = resultSchema.parse(await response.json());
   if (value.answer.score !== input.score || value.answer.consent_version !== input.consent_version) throw new Error("receipt mismatch");
   setAnswer(value.answer); setMessage(t("p0281"));
  } catch { setUncertain(true); setMessage(t("p0282")); }
  finally { pending.current = false; setBusy(false); }
 }
 async function recover() {
  if (pending.current) return; pending.current = true; setBusy(true);
  try {
   const response = await fetch(url + "&view=response", { cache: "no-store", signal: AbortSignal.timeout(7000) });
   if (response.status === 404) { setMessage(t("p0283")); return; }
   if (!response.ok) throw new Error("read unavailable");
   const value = answerSchema.parse(await response.json());
   if (value.consent_version !== definition.consent_version) throw new Error("version mismatch");
   setAnswer(value); setUncertain(false); setMessage(t("p0284"));
  } catch { setMessage(t("p0285")); }
  finally { pending.current = false; setBusy(false); }
 }
 return <section lang="es" aria-labelledby="survey-question">
  <h1 id="survey-question" className="pageTitle">{definition.prompt}</h1>
  {answer ? <div className="notice"><p>{t("p0286")} <strong data-testid="survey-stored-score">{answer.score}</strong> {t("p0287")}</p><p>{t("p0288")}</p></div> :
   <form onSubmit={submit} aria-describedby="survey-result">
    <label>{t("p0289")}<select name="score" required value={score} disabled={busy || uncertain || !definition.accepting} onChange={event => setScore(event.target.value)}>
     <option value="">{t("p0290")}</option>{Array.from({ length: 11 }, (_, n) => <option key={n} value={n}>{n}</option>)}
    </select></label>
    <p>{definition.consent_notice}</p>
    <label><input name="consent" type="checkbox" checked={consent} required disabled={busy || uncertain || !definition.accepting} onChange={event => setConsent(event.target.checked)} />{t("p0291")}</label>
    <button type="submit" disabled={busy || uncertain || !definition.accepting || score === "" || !consent}>{busy ? t("p0292") : t("p0293")}</button>
    {!definition.accepting ? <p>{t("p0294")}</p> : null}
   </form>}
  {uncertain && !answer ? <button type="button" onClick={recover} disabled={busy}>{t("p0295")}</button> : null}
  <p id="survey-result" role="status" aria-live="polite">{message}</p>
  <details><summary>{t("p0296")}</summary><p>{t("p0297")}</p><p>{t("p0298")}</p><p>{t("p0031")}</p></details>
 </section>;
}

````

### FILE: `src/app/customer/surveys/page.tsx`
```yaml
block_id: "TS-CUSTOMER-SURVEY-PORTAL:src/app/customer/surveys/page.tsx:v1"
operation: CREATE
provenance: AUTHORED
source: "local; existing admitted reference contracts, no external business source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "fa13f9dd93ef86526a17e4405e6f2de88de82f3a237a66cc61db478135fe087b"
variables: []
secrets_allowed: false
```
````tsx
import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import Link from "next/link";
import type { Route } from "next";
import { redirect, notFound } from "next/navigation";
import { readSession, allowed } from "@/platform/auth/session";
import { loadBusinessConfig } from "@/platform/config/load";
import { protectedGet } from "@/platform/backend/protected-client";
import { BackendProblem } from "@/platform/backend/public-client";
import { definitionSchema, answerSchema, surveyID, type SurveyAnswer } from "@/platform/surveys/contract";
import { CustomerSurvey } from "@/components/customer-survey";

export default async function CustomerSurveyPage({ searchParams }: { searchParams: Promise<Record<string, string | string[] | undefined>> }) {
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

 if ((await loadBusinessConfig()).features.customer_surveys !== true) notFound();
 const query = await searchParams; const id = query.surveyId, org = query.organizationId;
 if (Object.keys(query).length !== 2 || !surveyID.safeParse(id).success || !surveyID.safeParse(org).success) return <section lang="es"><h1 className="pageTitle">{t("p0063")}</h1><p>{t("p0064")}</p></section>;
 const survey = id as string, organization = org as string;
 const href = "/customer/surveys?surveyId=" + encodeURIComponent(survey) + "&organizationId=" + encodeURIComponent(organization);
 const session = await readSession(); if (!session) redirect(("/api/auth/login?return_to=" + encodeURIComponent(href)) as Route);
 if (!allowed(session, "customer:self") || (!session.organizations.includes(organization) && !session.organizations.includes("*"))) return <section lang="es"><h1 className="pageTitle">{t("p0006")}</h1><p>{t("p0065")}</p></section>;
 try {
  const path = "/v1/customer/surveys/" + encodeURIComponent(survey);
  const definition = definitionSchema.parse(await protectedGet<unknown>(session, path, { organization_id: organization }));
  if (definition.id !== survey) throw new Error("survey mismatch");
  let answer: SurveyAnswer | null = null;
  try { answer = answerSchema.parse(await protectedGet<unknown>(session, path + "/response", { organization_id: organization })); }
  catch (error) { if (!(error instanceof BackendProblem && error.status === 404 && error.code === "SURVEY_NOT_AVAILABLE")) throw error; }
  return <CustomerSurvey organization={organization} definition={definition} initialAnswer={answer} />;
 } catch {
  return <section lang="es"><h1 className="pageTitle">{t("p0066")}</h1><p>{t("p0067")}</p><Link href={href as Route}>{t("p0033")}</Link></section>;
 }
}

````

### FILE: `src/app/admin/surveys/page.tsx`
```yaml
block_id: "TS-CUSTOMER-SURVEY-PORTAL:src/app/admin/surveys/page.tsx:v1"
operation: CREATE
provenance: AUTHORED
source: "local; existing admitted reference contracts, no external business source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "c6af14d2747d191f8fdff4eeaf2b84611ccd2e061a1092685032ef7f478b4fad"
variables: []
secrets_allowed: false
```
````tsx
import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import Link from "next/link";
import type { Route } from "next";
import { notFound, redirect } from "next/navigation";
import { loadBusinessConfig } from "@/platform/config/load";
import { readSession, allowed } from "@/platform/auth/session";
import { protectedGet } from "@/platform/backend/protected-client";
import { summarySchema, surveyID } from "@/platform/surveys/contract";
export default async function SurveySummaryPage({ searchParams }: { searchParams: Promise<Record<string, string | string[] | undefined>> }) {
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

 if ((await loadBusinessConfig()).features.customer_surveys !== true) notFound();
 const query = await searchParams;
 if (Object.keys(query).length !== 2 || !surveyID.safeParse(query.surveyId).success || !surveyID.safeParse(query.organizationId).success) return <section lang="es"><h1 className="pageTitle">{t("p0021")}</h1><p>{t("p0022")}</p></section>;
 const id = query.surveyId as string, organization = query.organizationId as string;
 const href = "/admin/surveys?surveyId=" + encodeURIComponent(id) + "&organizationId=" + encodeURIComponent(organization);
 const session = await readSession(); if (!session) redirect(("/api/auth/login?return_to=" + encodeURIComponent(href)) as Route);
 if (!allowed(session, "surveys:read") || (!session.organizations.includes(organization) && !session.organizations.includes("*"))) return <section lang="es"><h1 className="pageTitle">{t("p0006")}</h1><p>{t("p0023")}</p></section>;
 try {
  const value = summarySchema.parse(await protectedGet<unknown>(session, "/v1/admin/surveys/" + encodeURIComponent(id) + "/summary", { organization_id: organization }));
  return <section lang="es"><h1 className="pageTitle">{t("p0024")}</h1><p>{t("p0025")} {value.responses}</p>{value.available ? <p>{t("p0026")} <strong data-testid="survey-nps">{value.nps!.toFixed(2)}</strong></p> : <p>{t("p0027")}</p>}<p>{t("p0028")}</p><details><summary>{t("p0029")}</summary><p>{t("p0030")}</p><p>{t("p0031")}</p></details></section>;
 } catch { return <section lang="es"><h1 className="pageTitle">{t("p0032")}</h1><Link href={href as Route}>{t("p0033")}</Link></section>; }
}

````

### FILE: `microsoft_playwright_browser_gate/tests/customer-survey.spec.mjs`
```yaml
block_id: "TS-CUSTOMER-SURVEY-PORTAL:microsoft_playwright_browser_gate/tests/customer-survey.spec.mjs:v1"
operation: CREATE
provenance: AUTHORED
source: "local; existing admitted reference contracts, no external business source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "e8c2b5a70290b4d7a5035990e6055cda34e8ee24905915f6faa82d9576f94476"
variables: []
secrets_allowed: false
```
````javascript
import { test, expect } from '@playwright/test';
import { createHash } from 'node:crypto';
import { createRequire } from 'node:module';
import { resolve } from 'node:path';
import { pathToFileURL } from 'node:url';

test('customer survey persists and recovers without repeated writes',async({page,context},info)=>{
 test.skip(process.env.ELITE_FEEDBACK_E2E!=='1','Requires disposable survey API and PostgreSQL');
 const identities=JSON.parse(process.env.ELITE_FEEDBACK_IDENTITIES);
 const require=createRequire(resolve(process.env.ELITE_WEB_ROOT,'package.json'));
 const {EncryptJWT}=await import(pathToFileURL(require.resolve('jose')).href);
 const id=info.project.name;
 async function identity(name) {
  const jwt=await new EncryptJWT(identities[name]).setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('300s').encrypt(createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest());
  await context.clearCookies();
  await context.addCookies([{name:'__Host-elite_session',value:jwt,url:process.env.ELITE_BASE_URL,secure:true,httpOnly:true,sameSite:'Lax'}]);
 }
 const customer='/customer/surveys?organizationId=org-a&surveyId='+id;
 const admin='/admin/surveys?organizationId=org-a&surveyId='+id;
 const api='/api/enterprise/surveys?organizationId=org-a&surveyId='+id;
 let writes=0;page.on('request',r=>{if(r.method()==='POST'&&r.url().includes('/api/enterprise/surveys'))writes++;});
 await identity('admin');await page.goto(admin);await expect(page.getByText('El resultado todavía no está disponible',{exact:false})).toBeVisible();
 await identity('alice');await page.goto(customer);
 await expect(page.getByRole('heading',{name:'Synthetic survey <img src=x onerror=alert(1)>'})).toBeVisible();
 await expect(page.locator('section img')).toHaveCount(0);
 await expect(page.getByRole('button',{name:'Enviar respuesta'})).toBeDisabled();
 await page.getByRole('combobox',{name:'Puntuación'}).selectOption('10');
 await expect(page.getByRole('button',{name:'Enviar respuesta'})).toBeDisabled();
 await page.getByRole('checkbox').check();
 // Two synchronous submissions before React can paint a new disabled state.
 await page.locator('form').evaluate(form=>{form.requestSubmit();form.requestSubmit();});
 await expect(page.getByTestId('survey-stored-score')).toHaveText('10');expect(writes).toBe(1);
 await page.reload();await expect(page.getByTestId('survey-stored-score')).toHaveText('10');expect(writes).toBe(1);
 expect((await context.request.get(api+'&view=summary')).status()).toBe(403);
 await page.goto(admin);await expect(page.getByRole('heading',{name:'Acceso denegado'})).toBeVisible();
 await identity('admin');await page.goto(admin);await expect(page.getByText('El resultado todavía no está disponible',{exact:false})).toBeVisible();
 await identity('bob');await page.goto(customer);
 await page.getByRole('combobox',{name:'Puntuación'}).selectOption('0');await page.getByRole('checkbox').check();
 // Complete the real upstream POST, then discard only its response at the browser.
 await page.route('**/api/enterprise/surveys?*',async route=>{
  if(route.request().method()!=='POST'){await route.continue();return;}
  const response=await route.fetch();expect(response.status()).toBe(201);await route.abort('failed');
 });
 await page.getByRole('button',{name:'Enviar respuesta'}).click();
 await expect(page.getByRole('button',{name:'Consultar respuesta guardada'})).toBeVisible();expect(writes).toBe(2);
 await expect(page.getByRole('button',{name:'Enviar respuesta'})).toBeDisabled();
 await page.getByRole('button',{name:'Consultar respuesta guardada'}).click();
 await expect(page.getByTestId('survey-stored-score')).toHaveText('0');expect(writes).toBe(2);
 await page.unrouteAll();await page.reload();await expect(page.getByTestId('survey-stored-score')).toHaveText('0');expect(writes).toBe(2);
 await identity('admin');await page.goto(admin);await expect(page.getByTestId('survey-nps')).toHaveText('0.00');await expect(page.getByText('Respuestas: 2',{exact:true})).toBeVisible();
 expect((await context.request.get(api.replace('org-a','org-b')+'&view=summary')).status()).toBe(403);
 await page.screenshot({path:info.outputPath('survey-results.png'),fullPage:true});
});
````

## 6. Configuration surface

CUSTOMER_SURVEYS_ENABLED=1 explicitly mounts API routes; features.customer_surveys=true
exposes BFF/pages. Both default to disabled. Existing OIDC/session/API-origin config
is reused. Definitions are inactive by default. Terms/dates/thresholds are immutable;
new policy versions require a new survey ID. DATABASE_URL is an operator secret.

## 7. Dependency bill

No dependency, version or external source added. Exact unchanged go.mod/go.sum and
pnpm lock/workspace/package.json come from the selected reference composition.
Go1.26.8, PostgreSQL18.6, Node24.20.0 and frozen pnpm11.25.0 tool identity is recorded
in the V400 evidence; existing dependency/provenance/redistribution conditions remain.

## 8. Apply order

Compose the compatible Go application/query/customer-journey foundations and, for
UI, the matching web/session/BFF and Playwright gate. Apply SQL0001–0054. Configure
project-owned definitions and permissions; enable only after target review.
Disable both flags to roll back exposure. A down migration destroys responses and
must not be used as a data-preserving production rollback. Use the documented
bounded retention command under the project retention/backup policy.

## 9. Verification

Focal HTTP/OIDC/PostgreSQL negatives,24 concurrent submissions, lost response and
actual database restart; production Next browser execution in4projects;8browser
writes produce8responses. Empty migration down/up and bounded retention CLI are
checked. UI typecheck/build and38focal tests PASS. Fuzz and reconstruction receipts
are bound in the V400 record. Readiness, full-domain coverage, runtime monitoring,
load/security/admission and signed final release remain independent gates.

## 10. Reconstruction evidence

See reconstruction_evidence/CONNECTED_CUSTOMER_SURVEYS_V400.md from the library
root. Exact file manifest and hashes bind the rebuilt sources to tested bytes.
This is CONDITIONED for this narrow reference, not 48/48 or target approval.

V400 final0.1.1: source bytes unchanged; compatibility binds atomic Go survey migration0.1.1. Existing browser evidence remains exact for all8web sources.

V402 composed delta: T2804 private es/en display, per-user/tenant preference and browser negotiation;1198messages,21exact guides,5hash-bound curricula; original commands/content/policies preserved. PRIVATE_LOCALE_RELEASE_V402.md/json. AUTHORED glue; no new dependency or corporate attribution.
