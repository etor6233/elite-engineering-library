import {beforeEach,expect,test,vi} from "vitest";
import {renderToStaticMarkup} from "react-dom/server";
const mocks=vi.hoisted(()=>({session:vi.fn(),get:vi.fn(),locale:vi.fn()}));
vi.mock("@/platform/auth/session",()=>({readSession:mocks.session,allowed:(s:{permissions:string[]},p:string)=>s.permissions.includes(p)}));
vi.mock("@/platform/backend/protected-client",()=>({protectedGet:mocks.get}));
vi.mock("next/navigation",()=>({redirect:(p:string)=>{throw new Error("REDIRECT "+p)}}));
vi.mock("@/platform/i18n/load-public-locale",()=>({loadPublicLocale:mocks.locale}));
import AppointmentsPage from "./appointments/page";
import CustomerPage from "./page";
import FactoryPage from "../factory/page";
beforeEach(()=>{vi.resetAllMocks();mocks.locale.mockResolvedValue({language:"es",locale:"es-AR",timeZone:"America/Argentina/Buenos_Aires"});mocks.session.mockResolvedValue({subject:"customer",tenantId:"tenant",organizations:["org-a"],permissions:["customer:self","factory:read"],accessToken:"test-only"});mocks.get.mockImplementation(async(_s:unknown,p:string)=>p.endsWith("journey")?{appointments:[],quotes:[],handovers:[]}:{items:[],next_cursor:p.endsWith("orders")?"order-025":p.endsWith("service-cases")?"case-025":"unit-025"});});
test("customer exposes each backend next cursor independently",async()=>{
 const html=renderToStaticMarkup(await CustomerPage());
 expect(html).toContain('href="/customer?orders_after=order-025"');expect(html).toContain('href="/customer?cases_after=case-025"');
});
test("factory exposes backend next cursor",async()=>{expect(renderToStaticMarkup(await FactoryPage())).toContain('href="/factory?units_after=unit-025"');});

test("customer forwards opaque cursors separately and preserves the other list",async()=>{
 const html=renderToStaticMarkup(await CustomerPage({searchParams:Promise.resolve({orders_after:"old order/&",cases_after:"old-case",organization_id:"ignored",limit:"1000"})}));
 const calls=mocks.get.mock.calls;expect(calls.find(c=>c[1].endsWith("orders"))?.[2]).toEqual({organization_id:"org-a",limit:"25",after:"old order/&"});
 expect(calls.find(c=>c[1].endsWith("service-cases"))?.[2]).toEqual({organization_id:"org-a",limit:"25",after:"old-case"});
 expect(html).toContain('/customer?orders_after=order-025&amp;cases_after=old-case');expect(html).toContain('/customer?orders_after=old+order%2F%26&amp;cases_after=case-025');
 expect(html).toContain('href="/customer?cases_after=old-case"');expect(html).toContain('href="/customer?orders_after=old+order%2F%26"');
});
test("last and empty pages have no next link and retain first-page recovery",async()=>{
 mocks.get.mockImplementation(async(_s:unknown,p:string)=>p.endsWith("journey")?{appointments:[],quotes:[],handovers:[]}:{items:[]});
 const html=renderToStaticMarkup(await CustomerPage({searchParams:Promise.resolve({orders_after:"last"})}));expect(html).not.toContain('Ver siguientes');expect(html).toContain('href="/customer"');expect(html).toContain('No hay registros en esta página.');
 const factory=renderToStaticMarkup(await FactoryPage({searchParams:Promise.resolve({units_after:"last"})}));expect(factory).not.toContain('Ver siguientes');expect(factory).toContain('href="/factory"');
});
test("factory reads the session organization and fixed limit",async()=>{await FactoryPage({searchParams:Promise.resolve({units_after:"unit-025",organization_id:"ignored",limit:"1000"})});expect(mocks.get).toHaveBeenCalledWith(expect.anything(),"/v1/factory/units",{organization_id:"org-a",limit:"25",after:"unit-025"});});
test.each([{value:["duplicate","cursor"]},{value:"x".repeat(513)}])("ambiguous or excessive cursors stop before reads",async({value})=>{await expect(CustomerPage({searchParams:Promise.resolve({orders_after:value})})).rejects.toThrow("Invalid page cursor");expect(mocks.get).not.toHaveBeenCalled();});
test("missing session still redirects before querying",async()=>{mocks.session.mockResolvedValue(null);await expect(CustomerPage({searchParams:Promise.resolve({orders_after:"cursor"})})).rejects.toThrow("REDIRECT");expect(mocks.get).not.toHaveBeenCalled();});
test("existing permission checks run before paging",async()=>{mocks.session.mockResolvedValue({organizations:["org-a"],permissions:[]});const html=renderToStaticMarkup(await FactoryPage({searchParams:Promise.resolve({units_after:"cursor"})}));expect(html).toContain("Acceso denegado");expect(mocks.get).not.toHaveBeenCalled();});

const appointmentDates = [
 {locale:"es-AR",language:"es",timeZone:"America/Argentina/Buenos_Aires",date:"2035-01-02T01:30:00Z",dateLabel:"1 ene",hour:"10:30 p. m."},
 {locale:"en-US",language:"en",timeZone:"America/New_York",date:"2035-01-02T01:30:00Z",dateLabel:"Jan 1",hour:"8:30 PM"},
 {locale:"en-US",language:"en",timeZone:"America/New_York",date:"2035-07-02T01:30:00Z",dateLabel:"Jul 1",hour:"9:30 PM"},
 {locale:"en-US",language:"en",timeZone:"Asia/Tokyo",date:"2035-01-02T01:30:00Z",dateLabel:"Jan 2",hour:"10:30 AM"},
];
test.each(appointmentDates)("both customer appointment views use configured $timeZone / $date",async(row)=>{
 mocks.locale.mockResolvedValue(row);
 mocks.get.mockImplementation(async(_s:unknown,p:string)=>p.endsWith("journey")?{appointments:[{id:"appointment",kind:"consultation",starts_at:row.date,state:"requested",version:1}],quotes:[],handovers:[]}:{items:[]});
 for(const view of [CustomerPage,AppointmentsPage]){
  const html=renderToStaticMarkup(await view()).replace(/\u00a0|\u202f/g," ");
  expect(html).toContain('dateTime="'+row.date+'"');expect(html).toContain(row.dateLabel);expect(html).toContain(row.hour);expect(html).toContain(row.timeZone);
 }
});
test("unverifiable appointment time fails instead of showing Invalid Date",async()=>{
 mocks.get.mockImplementation(async(_s:unknown,p:string)=>p.endsWith("journey")?{appointments:[{id:"appointment",kind:"consultation",starts_at:"not-a-date",state:"requested",version:1}],quotes:[],handovers:[]}:{items:[]});
 for(const view of [CustomerPage,AppointmentsPage]) await expect(view()).rejects.toThrow("Invalid appointment date");
});
test("appointment page denies access before loading time configuration or journey",async()=>{
 mocks.session.mockResolvedValue({organizations:["org-a"],permissions:[]});expect(renderToStaticMarkup(await AppointmentsPage())).toContain("Acceso denegado");expect(mocks.get).not.toHaveBeenCalled();expect(mocks.locale).not.toHaveBeenCalled();
});

import {createElement} from "react";
import {isCancellationReceipt,CustomerAppointmentActions} from "@/components/customer-appointment-actions";
const cancellationReceipt={id:"appointment",organization_id:"org-a",state:"cancelled",version:2};
test.each([
 {name:"exact receipt",value:cancellationReceipt,valid:true},
 {name:"null",value:null,valid:false},
 {name:"array",value:[],valid:false},
 {name:"different appointment",value:{...cancellationReceipt,id:"other"},valid:false},
 {name:"different organization",value:{...cancellationReceipt,organization_id:"org-b"},valid:false},
 {name:"unconfirmed state",value:{...cancellationReceipt,state:"requested"},valid:false},
 {name:"stale version",value:{...cancellationReceipt,version:1},valid:false},
 {name:"future version",value:{...cancellationReceipt,version:3},valid:false},
 {name:"string version",value:{...cancellationReceipt,version:"2"},valid:false},
 {name:"nonfinite version",value:{...cancellationReceipt,version:NaN},valid:false},
])("cancellation response binding: $name",({value,valid})=>{expect(isCancellationReceipt(value,"org-a",{id:"appointment",version:1})).toBe(valid)});
test("customer can reach appointment management from their account",async()=>{expect(renderToStaticMarkup(await CustomerPage())).toContain('href="/customer/appointments"');expect(renderToStaticMarkup(await CustomerPage())).toContain('Gestionar mis turnos')});
test("empty appointment view provides an explicit state",()=>{expect(renderToStaticMarkup(createElement(CustomerAppointmentActions,{organizationId:"org-a",appointments:[]}))).toContain("No hay turnos disponibles")});
