import{createElement}from"react";
import{renderToStaticMarkup}from"react-dom/server";
import{expect,it}from"vitest";
import{CustomerHandoverActions}from"@/components/customer-handover-actions";
it("renders the exact server-formatted handover labels regardless of client timezone",()=>{
 const html=renderToStaticMarkup(createElement(CustomerHandoverActions,{organizationId:"store",exceptions:[],handovers:[{id:"handover",order_id:"order",stock_unit_id:"stock",state:"accepted",version:3,checklist_items:[{id:"serial",ordinal:1,prompt:"Serial",response_type:"serial",required:true}],checklist_id:"checklist",checklist_version:1,checklist_completed_at:"2026-09-11T12:00:00Z",customer_accepted_at:"2026-09-11T12:05:00Z",checklistCompletedLabel:"11/09/2026 09:00 (America/Argentina/Buenos_Aires)",acceptedLabel:"11/09/2026 09:05 (America/Argentina/Buenos_Aires)"}]}));
 expect(html).toContain("11/09/2026 09:00 (America/Argentina/Buenos_Aires)");expect(html).toContain("11/09/2026 09:05 (America/Argentina/Buenos_Aires)");
});
