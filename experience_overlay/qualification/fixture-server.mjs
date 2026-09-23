import {evidence} from "./paths.mjs";
// AUTHORED loopback fixture transport. Not a replacement for Go/PG owner acceptance.
import {handle as handleDocuments} from './document-fixture.mjs';
import http from 'node:http';
import https from 'node:https';
import {readFileSync,writeFileSync} from 'node:fs';
import {resolve} from 'node:path';
const root=resolve(import.meta.dirname),writes=[];let completion=null;
let unitMode='complete';
const unitPlan={purchase_order_id:'supplyOrders-v403',destination_organization_id:'store',factory_organization_id:'factory',supplier_id:'supplier',state:'shipped',purchase_version:'1',version:'1',currency:'ARS',total_minor_units:'120000',policy_code:'strict-serial-reference/v1',demand_reference:'fixture-order',lines:[{id:'line',variant_id:'variant-v403',quantity:2}],units:[{id:'unit-1',line_id:'line',state:'shipped',serial_number:'001aB',stock_state:'in-transit'}]};
const fixed='2026-09-14T02:00:00Z';
const checklist={id:'standard-delivery',organization_id:'store',version:1,state:'published',title:'Entrega verificada',items:[{id:'serial-observed',prompt:'Número de serie observado',response_type:'serial',required:true},{id:'asset-condition',prompt:'Revisé el estado del activo',response_type:'confirmation',required:true},{id:'inspection-proof',prompt:'Expediente de inspección',response_type:'evidence',required:true}]};
function send(res,status,value){res.writeHead(status,{'content-type':'application/json'});res.end(JSON.stringify(value));}
const backend=http.createServer(async(req,res)=>{
 if(await handleDocuments(req,res))return;
 const url=new URL(req.url,'http://127.0.0.1:4314');
 if(url.pathname==='/fixture/reset' && req.method==='POST'){writes.length=0;completion=null;return send(res,200,{reset:true});}
 if(url.pathname==='/fixture/unit-mode'&&req.method==='POST'){unitMode=url.searchParams.get('mode')??'complete';return send(res,200,{mode:unitMode});}
 if(url.pathname==='/fixture/stats')return send(res,200,{writes,completion});
 if(req.method==='GET'&&url.pathname==='/v1/public/elite-mobility/models')return send(res,200,{models:[{id:'m',code:'urban',displayName:'Urban',vehicleClass:'bicycle',specification:{}}]});
 if(req.headers.authorization!=='Bearer fixture-access-token')return send(res,403,{code:'FORBIDDEN'});
 if(req.method==='GET'){
  if(url.searchParams.get('organization_id')!=='store')return send(res,403,{code:'FORBIDDEN'});
  if(url.pathname==='/v1/franchise/supply/orders/supplyOrders-v403')return send(res,200,{...unitPlan,...(unitMode==='partial'?{next_unit_id:'unit-1'}:{}),units:unitMode==='ambiguous'?[...unitPlan.units,{...unitPlan.units[0],id:'unit-2'}]:unitPlan.units});
  if(url.pathname==='/v1/franchise/delivery-checklists/result')return send(res,200,checklist);
  if(url.pathname.endsWith('/checklist-result'))return completion?send(res,200,completion):send(res,404,{code:'NOT_FOUND'});
  if(url.pathname==='/v1/commerce/orders')return send(res,200,{orders:[],stock:[],truncated:false});
  if(url.pathname==='/v1/franchise/leads')return send(res,200,{items:[{id:'lead-v403',organization_id:'store',state:'new',source_code:'website',created_at:fixed,version:1}]});
  if(url.pathname==='/v1/franchise/agenda')return send(res,200,{appointments:[],resources:[],truncated:false});
  return send(res,200,{items:[]});
 }
 if(req.method==='POST'){
  const chunks=[];let length=0;for await(const chunk of req){length+=chunk.length;if(length>65536)return send(res,413,{code:'TOO_LARGE'});chunks.push(chunk);}
  let body;try{body=JSON.parse(Buffer.concat(chunks).toString());}catch{return send(res,400,{code:'INVALID_JSON'});}
  if(body.organization_id!=='store')return send(res,403,{code:'FORBIDDEN'});
  writes.push({path:url.pathname,body});writeFileSync(resolve(evidence,'fixture-transactions.json'),JSON.stringify({method:'AUTHORED_HTTP_BACKEND_FIXTURE_WITH_REAL_NEXT_BFF',writes},null,2));
  if(url.pathname.endsWith('/complete-checklist')){
   if(completion||body.version!==1)return send(res,409,{code:'CONFLICT'});
   if(body.checklist_id!==checklist.id||body.checklist_version!==1||body.responses?.length!==3||body.responses[1]?.response_text!=='confirmed')return send(res,400,{code:'INVALID_CHECKLIST'});
   completion={handover_id:'handover-v403',organization_id:'store',state:'presented',version:2,checklist_id:checklist.id,checklist_version:1,completed_at:fixed,actor_subject:'operator',responses:body.responses};
   return send(res,200,{id:'handover-v403',organization_id:'store',state:'presented',version:2,checklist_id:checklist.id,checklist_version:1,checklist_completed_at:fixed});
  }
  if(url.pathname==='/v1/franchise/quotes')return send(res,200,{id:'quote-v403',organization_id:'store',lead_id:body.lead_id,variant_id:body.variant_id,price_book_id:body.price_book_id,valid_until:body.valid_until,state:'issued',currency:'ARS',total_minor_units:120000,version:1});
  if(url.pathname==='/v1/franchise/availability')return send(res,200,{id:'availability-v403',organization_id:'store',resource_id:body.resource_id,entry_type:body.entry_type,reason_code:body.reason_code,starts_at:body.starts_at,ends_at:body.ends_at,state:'active',version:1});
  if(url.pathname.endsWith('/assign'))return send(res,200,{id:'lead-v403',organization_id:'store',state:'new',assigned_subject:body.assigned_subject,version:2});
  return send(res,404,{code:'NOT_FOUND'});
 }
 return send(res,405,{code:'METHOD_NOT_ALLOWED'});
});
backend.listen(4314,'127.0.0.1');
https.createServer({key:readFileSync(resolve(evidence,'fixture.key')),cert:readFileSync(resolve(evidence,'fixture.crt'))},(req,res)=>{
 const outgoing=http.request({hostname:'127.0.0.1',port:4311,path:req.url,method:req.method,headers:{...req.headers,host:'127.0.0.1:4313','x-forwarded-host':'127.0.0.1:4313','x-forwarded-proto':'https'}},up=>{res.writeHead(up.statusCode,up.headers);up.pipe(res);});outgoing.on('error',()=>{res.writeHead(502);res.end('Fixture proxy unavailable');});req.pipe(outgoing);
}).listen(4313,'127.0.0.1');
console.log('Loopback HTTPS fixture 4313 and backend 4314 ready');
