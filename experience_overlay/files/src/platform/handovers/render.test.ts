import {createElement}from"react";
import{renderToStaticMarkup}from"react-dom/server";
import{expect,it,vi}from"vitest";
vi.mock("next/navigation",()=>({useRouter:()=>({refresh:vi.fn()})}));
import{HandoverOperationsPanel}from"@/components/handover-operations-panel";
import{ChecklistCompletionPanel}from"@/components/franchise-command-panel";
it("exposes unavailable data distinctly and does not authorize a mutation before context",()=>{
 const absent=renderToStaticMarkup(createElement(HandoverOperationsPanel,{orders:null,organization:"store",scope:"a".repeat(64)}));expect(absent).toContain("lista de pedidos no está disponible");expect(absent).not.toContain("Registrar cierre comercial");
 const present=renderToStaticMarkup(createElement(HandoverOperationsPanel,{orders:[{id:"order"}],organization:"store",scope:"a".repeat(64)}));expect(present).toContain("Consultar estado de entrega");expect(present).not.toContain("Registrar cierre comercial");expect(present).not.toContain("Preparar entrega");
});
it("keeps presentation disabled until published fields load; never asks operator for IDs or versions",()=>{
 const html=renderToStaticMarkup(createElement(ChecklistCompletionPanel,{organization:"store",scope:"a".repeat(64),initialHandover:{id:"handover-from-server",version:7}}));
 expect(html).not.toContain('name="handoverId"');expect(html).not.toContain('name="handoverVersion"');expect(html).toContain("Cargando información");expect(html).toContain("disabled=\"\"");expect(html).toContain("Confirmar revisión");
});
