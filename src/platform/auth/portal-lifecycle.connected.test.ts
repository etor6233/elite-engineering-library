import {describe,it,expect,vi,afterEach} from "vitest";
import {GET as login} from "@/app/api/auth/login/route";
import {GET as callback} from "@/app/api/auth/callback/route";
import {POST as logout} from "@/app/api/auth/logout/route";
import {SESSION_COOKIE,FLOW_COOKIE,openFlow} from "./session";
import {portalProfile} from "./portal-profile";
import {readPortalSession,createPortalSession} from "./portal-lifecycle";
import {portalServiceAccessToken} from "./portal-protocol";
import {POST as maintenance} from "@/app/api/internal/oidc/session-maintenance/route";

const enabled=!!process.env.PORTAL_FIXTURE_CONTROL_URL;
const controls=process.env.PORTAL_FIXTURE_CONTROL_URL??"http://127.0.0.1:9";
async function control(value:unknown){const response=await fetch(controls+"/fixture/control",{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify(value)});expect(response.status).toBe(200)}
async function stats(){return await(await fetch(controls+"/fixture/stats")).json() as Record<string,number>}
async function begin(){const response=await login(new Request("http://127.0.0.1:4567/api/auth/login?return_to=/customer"));expect(response.status).toBe(303);const flow=response.cookies.get(FLOW_COOKIE)?.value;expect(flow).toBeTruthy();const authorize=await fetch(response.headers.get("location")!,{redirect:"manual"});expect(authorize.status).toBe(303);return{flow:flow!,url:authorize.headers.get("location")!}}
async function finish(flow:{flow:string;url:string}){return callback(new Request(flow.url,{headers:{cookie:`${FLOW_COOKIE}=${flow.flow}`}}))}
async function loginCookie(){const response=await finish(await begin());if(response.status!==303){const diagnostic=await begin();try{await createPortalSession(portalProfile(),await openFlow(diagnostic.flow),new URL(diagnostic.url))}catch(error){const e=error as {name?:string;code?:string;message?:string;cause?:{code?:string;message?:string}};throw new Error(JSON.stringify({name:e.name,code:e.code,message:/^OIDC_[A-Z_]+$/.test(e.message??"")?e.message:"SDK_FAILURE",causeCode:e.cause?.code,causeMessage:e.cause?.message?.replace(/[A-Za-z0-9_-]{60,}/g,"REDACTED"),stats:await stats()}))}}expect(response.status).toBe(303);const cookie=response.cookies.get(SESSION_COOKIE)?.value;expect(cookie).toBeTruthy();expect(cookie).not.toContain("refresh");return cookie!}
afterEach(()=>vi.useRealTimers());
describe.skipIf(!enabled)("official OIDC SDK → BFF route → authenticated Go bridge → PostgreSQL",()=>{
 it("performs code/PKCE login, one concurrent refresh, lost commit recovery and durable logout",async()=>{
  const p=portalProfile(),cookie=await loginCookie();const session=await readPortalSession(p,cookie);expect(session?.subject).toBe("person");expect(session?.permissions).toEqual(["customer:read"]);
  const before=await stats();await control({expire:true,mode:"commit_lost"});const sessions=await Promise.all(Array.from({length:20},()=>readPortalSession(p,cookie)));expect(sessions.every(s=>s?.subject==="person")).toBe(true);expect((await stats()).refresh_grants).toBe(before.refresh_grants!+1);
  const denied=await logout(new Request(p.post_logout_url+"api/auth/logout",{method:"POST",headers:{origin:"https://foreign.invalid",cookie:`${SESSION_COOKIE}=${cookie}`}}));expect(denied.status).toBe(403);expect(await readPortalSession(p,cookie)).not.toBeNull();
  const done=await logout(new Request(p.post_logout_url+"api/auth/logout",{method:"POST",headers:{origin:new URL(p.post_logout_url).origin,cookie:`${SESSION_COOKIE}=${cookie}`}}));expect(done.status).toBe(303);expect(done.cookies.get(SESSION_COOKIE)?.value).toBe("");expect(await readPortalSession(p,cookie)).toBeNull();expect((await stats()).revocations).toBeGreaterThanOrEqual(2);
 });
 it("never reuses consumed refresh after an ambiguous provider response",async()=>{
  const p=portalProfile(),cookie=await loginCookie();await control({expire:true,mode:"grant_lost"});const before=await stats();expect(await readPortalSession(p,cookie)).toBeNull();expect(await readPortalSession(p,cookie)).toBeNull();const after=await stats();expect(after.refresh_grants).toBe(before.refresh_grants!+1);expect(after.reauth_required).toBeGreaterThan(0);
 });
 it("invalidates locally before failed remote revocation and recovers explicitly",async()=>{
  const p=portalProfile(),cookie=await loginCookie();const request=()=>new Request(p.post_logout_url+"api/auth/logout",{method:"POST",headers:{origin:new URL(p.post_logout_url).origin,cookie:`${SESSION_COOKIE}=${cookie}`}});
  await control({mode:"revoke_unavailable"});const first=await logout(request());expect(first.status).toBe(503);expect(await first.json()).toMatchObject({local_session_revoked:true});expect(await readPortalSession(p,cookie)).toBeNull();await control({mode:""});expect((await logout(request())).status).toBe(303);
 });
 it("rejects bad state, replayed code, issuer/audience/expiry/nonce/tenant and excess permissions",async()=>{
  const original=await begin();const wrong=new URL(original.url);wrong.searchParams.set("state","foreign");expect((await finish({...original,url:wrong.href})).status).toBe(400);
  const response=await finish(original);expect(response.status).toBe(303);expect((await finish(original)).status).toBe(400);
  for(const mode of ["wrong_issuer","wrong_audience","expired","wrong_nonce","foreign_tenant","extra_permission"]){await control({mode});expect((await finish(await begin())).status).toBe(400)}await control({mode:""});
 });
 it("rejects unsigned/tampered handles and accepts SDK JWKS rotation after its cache interval",async()=>{
  const p=portalProfile(),cookie=await loginCookie();expect(await readPortalSession(p,cookie.slice(0,-4)+"bad!")).toBeNull();const before=await stats();await control({expire:true,rotate:true});vi.useFakeTimers({toFake:["Date"]});vi.setSystemTime(Date.now()+61000);expect((await readPortalSession(p,cookie))?.subject).toBe("person");expect((await stats()).jwks).toBeGreaterThan(before.jwks!);
 });
 it("rejects a valid ciphertext swapped between two sessions of the same profile",async()=>{
  const p=portalProfile(),first=await loginCookie(),second=await loginCookie();expect(await readPortalSession(p,first)).not.toBeNull();await control({swap:true});expect(await readPortalSession(p,first)).toBeNull();expect(await readPortalSession(p,second)).toBeNull();await control({swap:true});expect(await readPortalSession(p,first)).not.toBeNull();
 });
 it("recovers provider revocation without cookies and sweeps bounded retention with durable unconfirmed counts",async()=>{
  const p=portalProfile();await loginCookie();await control({expire_absolute:true});const request=(bearer:string)=>new Request(p.post_logout_url+"api/internal/oidc/session-maintenance",{method:"POST",headers:{authorization:bearer}});
  expect((await maintenance(request("Bearer invalid"))).status).toBe(401);
  const bearer="Bearer "+await portalServiceAccessToken(p);let confirmed=0;
  for(let i=0;i<8;i++){const response=await maintenance(request(bearer));expect(response.status).toBe(200);const value=await response.json() as {claimed:number;confirmed:number};expect(value.claimed).toBeLessThanOrEqual(2);confirmed+=value.confirmed;if(!value.claimed)break}expect(confirmed).toBeGreaterThan(0);
  await loginCookie();await control({retention:true});const response=await maintenance(request(bearer));expect(response.status).toBe(200);expect(await response.json()).toMatchObject({unconfirmed_purged:1});
 });
});
