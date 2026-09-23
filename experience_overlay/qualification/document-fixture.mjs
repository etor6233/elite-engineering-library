// AUTHORED loopback HTTP fixture for the real V403 Next BFF/browser journey.
// No Go/PG execution, document scanner, OCR provider or OCR ground truth claim.
import {createHash} from 'node:crypto';

export const profile = 'a'.repeat(64);
const sha = value => createHash('sha256').update(value).digest('hex');
const segment = (marker, payload) => Buffer.concat([Buffer.from([255, marker, (payload.length + 2) >> 8, (payload.length + 2) & 255]), payload]);
// AUTHORED one-pixel grayscale baseline JPEG: quantization=1; DC=0 and EOB.
// The test's four suggested fields are independent synthetic values, not decoded from this pixel.
export const syntheticJpeg = Buffer.concat([
  Buffer.from([255,216]),
  segment(0xdb, Buffer.from([0,...Array(64).fill(1)])),
  segment(0xc0, Buffer.from([8,0,1,0,1,1,1,0x11,0])),
  segment(0xc4, Buffer.from([0,1,...Array(15).fill(0),0,0x10,1,...Array(15).fill(0),0])),
  segment(0xda, Buffer.from([1,1,0,0,63,0])),
  Buffer.from([0x3f,255,217]),
]);
const identities = {
  'Bearer fixture-document-writer': {subject:'document-writer', tenant:'tenant', organization:'store', write:true, process:true, review:false},
  'Bearer fixture-document-reviewer': {subject:'document-reviewer', tenant:'tenant', organization:'store', write:false, process:false, review:true},
};
const records = new Map(), effects = [], requests = [];
const send = (res,status,value,headers={}) => {res.writeHead(status,{'content-type':'application/json','cache-control':'no-store','x-content-type-options':'nosniff',...headers});res.end(Buffer.isBuffer(value)?value:JSON.stringify(value));return true;};
const no = (res,status,code) => send(res,status,{code});
const exact = (value,keys) => value && typeof value==='object' && !Array.isArray(value) && Object.keys(value).sort().join('|')===[...keys].sort().join('|');
const text = (value,max) => typeof value==='string' && value && value.trim()===value && !value.includes('\0') && value.isWellFormed() && Buffer.byteLength(value)<=max;
const validFields = v => exact(v,['invoice_number','vendor','total','currency']) && text(v.invoice_number,128) && text(v.vendor,512) && text(v.total,64) && /^[A-Z]{3}$/.test(v.currency);
const view = r => structuredClone(r.view);

async function body(req,limit) {
  const chunks=[];let size=0;
  for await (const chunk of req) {size+=chunk.length;if(size>limit)throw new Error('TOO_LARGE');chunks.push(chunk);}
  return Buffer.concat(chunks);
}

export async function handle(req,res) {
  const url = new URL(req.url,'http://127.0.0.1:4314');
  if (!url.pathname.startsWith('/fixture/documents/') && !url.pathname.startsWith('/v1/documents/')) return false;
  if (!['127.0.0.1','::1','::ffff:127.0.0.1'].includes(req.socket.remoteAddress)) return no(res,403,'LOOPBACK_ONLY');
  if (url.pathname==='/fixture/documents/reset' && req.method==='POST') {records.clear();effects.length=0;requests.length=0;return send(res,200,{reset:true,method:'AUTHORED_DOCUMENT_HTTP_FIXTURE'});}
  if (url.pathname==='/fixture/documents/stats' && req.method==='GET') return send(res,200,{method:'AUTHORED_DOCUMENT_HTTP_FIXTURE',provider_executed:false,database_executed:false,ocr_truth:false,synthetic_jpeg_sha256:sha(syntheticJpeg),effects,requests,documents:[...records.values()].map(view)});
  if (url.pathname.startsWith('/fixture/')) return no(res,404,'NOT_FOUND');
  const actor=identities[req.headers.authorization];
  if (!actor) return no(res,401,'UNAUTHENTICATED');
  const route = /^\/v1\/documents\/([a-f0-9]{8}-[a-f0-9]{4}-[1-8][a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12})(?:\/(original|process|review|decision|evidence\/(?:security|provider|analysis)))?$/.exec(url.pathname);
  if (!route || url.search) return no(res,400,'INVALID_QUERY');
  const [,id,action='read']=route;
  requests.push({method:req.method,path:url.pathname,id,actor:actor.subject});
  let record=records.get(id);
  if (req.method==='PUT' && action==='original') {
    if (!actor.write) return no(res,403,'FORBIDDEN');
    let bytes;try{bytes=await body(req,2*1024*1024);}catch{return no(res,413,'TOO_LARGE');}
    const name=req.headers['x-document-name'];
    if(req.headers['content-type']!=='application/octet-stream' || req.headers['x-document-profile-sha256']!==profile || !text(name,128) || !/\.jpe?g$/i.test(name) || !bytes.equals(syntheticJpeg) || req.headers['x-document-sha256']!==sha(bytes))return no(res,400,'INVALID_CONTRACT');
    if(record) return record.view.uploader===actor.subject && record.view.name===name && record.view.original_sha256===sha(bytes)?send(res,201,view(record)):no(res,409,'CONSULT_RECORDED_STATE');
    record={tenant:actor.tenant,organization:actor.organization,original:bytes,view:{document_id:id,name,uploader:actor.subject,original_sha256:sha(bytes),profile_sha256:profile,mode:'FIXTURE',state:'QUARANTINED'},evidence:{}};
    records.set(id,record);effects.push({id,action:'receive',actor:actor.subject});return send(res,201,view(record));
  }
  if (!record || record.tenant!==actor.tenant || record.organization!==actor.organization) return no(res,404,'NOT_FOUND');
  if(req.method==='GET') {
    if(action==='read')return send(res,200,view(record));
    if(action==='original')return send(res,200,record.original,{'content-type':'application/octet-stream','X-Document-SHA256':record.view.original_sha256});
    if(action.startsWith('evidence/')) {const bytes=record.evidence[action.slice(9)];return bytes?send(res,200,bytes,{'X-Document-Evidence-SHA256':sha(bytes)}):no(res,404,'NOT_FOUND');}
    return no(res,405,'METHOD_NOT_ALLOWED');
  }
  if(req.method!=='POST' || req.headers['content-type']!=='application/json')return no(res,405,'METHOD_NOT_ALLOWED');
  let payload;try{payload=JSON.parse((await body(req,32768)).toString('utf8'));}catch{return no(res,400,'INVALID_CONTRACT');}
  if(action==='process') {
    if(!actor.process)return no(res,403,'FORBIDDEN');
    if(!exact(payload,[]))return no(res,400,'INVALID_CONTRACT');
    if(record.view.state!=='QUARANTINED')return no(res,409,'CONSULT_RECORDED_STATE');
    for(const part of ['security','provider','analysis'])record.evidence[part]=Buffer.from(JSON.stringify({method:'AUTHORED_DOCUMENT_HTTP_FIXTURE',part,mode:'FIXTURE',document_id:id,original_sha256:record.view.original_sha256,profile_sha256:profile,provider_executed:false,ocr_truth:false}));
    record.view.evidence_sha256=sha(Buffer.concat(Object.values(record.evidence)));
    record.view.suggested={invoice_number:'SYNTHETIC-001',vendor:'Synthetic supplier',total:'123.45',currency:'USD'};
    record.view.state='REVIEW_REQUIRED';effects.push({id,action,actor:actor.subject});return send(res,200,view(record));
  }
  if(action==='review') {
    if(!actor.write || record.view.uploader!==actor.subject)return no(res,403,'FORBIDDEN');
    if(!exact(payload,['evidence_sha256','fields']) || payload.evidence_sha256!==record.view.evidence_sha256 || !validFields(payload.fields))return no(res,400,'INVALID_CONTRACT');
    if(record.view.state!=='REVIEW_REQUIRED')return no(res,409,'CONSULT_RECORDED_STATE');
    record.view.proposal={schema:'document-review/v1',document_id:id,organization_id:actor.organization,original_sha256:record.view.original_sha256,profile_sha256:profile,evidence_sha256:record.view.evidence_sha256,mode:'FIXTURE',fields:structuredClone(payload.fields)};
    record.view.payload_sha256=sha(JSON.stringify(record.view.proposal));record.view.state='REVIEW_PENDING';effects.push({id,action,actor:actor.subject});return send(res,200,view(record));
  }
  if(action==='decision') {
    if(!actor.review)return no(res,403,'FORBIDDEN');
    if(actor.subject===record.view.uploader)return no(res,409,'CONSULT_RECORDED_STATE');
    if(!exact(payload,['payload_sha256','approved','reason']) || payload.payload_sha256!==record.view.payload_sha256 || typeof payload.approved!=='boolean' || !text(payload.reason,2048))return no(res,400,'INVALID_CONTRACT');
    if(record.view.state!=='REVIEW_PENDING')return no(res,409,'CONSULT_RECORDED_STATE');
    record.view.state=payload.approved?'PERSISTED':'REJECTED';record.view.reviewer=actor.subject;record.view.reason=payload.reason;effects.push({id,action,actor:actor.subject,approved:payload.approved});return send(res,200,view(record));
  }
  return no(res,405,'METHOD_NOT_ALLOWED');
}
