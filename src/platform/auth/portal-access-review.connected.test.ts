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
describe.skipIf(!enabled)("IdP access review through official SDK and durable portal",()=>{
 it("applies an explicit permission withdrawal on refresh without granting a replacement role",async()=>{const p=portalProfile(),cookie=await loginCookie();expect((await readPortalSession(p,cookie))?.permissions).toEqual(["customer:read"]);await control({expire:true,mode:"permission_withdrawn"});const changed=await readPortalSession(p,cookie);expect(changed?.subject).toBe("person");expect(changed?.permissions).toEqual([]);await control({mode:""})});
 it("rejects a changed subject during refresh and never links it to the original session",async()=>{const p=portalProfile(),cookie=await loginCookie();await control({expire:true,mode:"subject_changed"});expect(await readPortalSession(p,cookie)).toBeNull();expect(await readPortalSession(p,cookie)).toBeNull();await control({mode:""})});
 it("requires reauthentication after IdP deprovisioning and rejects a fresh authorization code",async()=>{const p=portalProfile(),cookie=await loginCookie();await control({expire:true,mode:"account_disabled"});expect(await readPortalSession(p,cookie)).toBeNull();expect(await readPortalSession(p,cookie)).toBeNull();expect((await finish(await begin())).status).toBe(400);await control({mode:""})});
});
