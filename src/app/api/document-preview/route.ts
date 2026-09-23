// AUTHORED isolated viewer. PDF.js is separately pinned VERBATIM distribution.
import {randomBytes} from "node:crypto";
import {readSession,allowed} from "@/platform/auth/session";
import {documentConfig} from "@/platform/documents-vnext/server";
export const runtime="nodejs";
export async function GET(){
 const session=await readSession();if(!session)return new Response("Unauthorized",{status:401});
 let config;try{config=documentConfig()}catch{return new Response("Unavailable",{status:503})}
 if(session.tenantId!==config.tenant||!session.organizations.includes(config.organization)||!["read","write","process","review"].some(p=>allowed(session,"documents:"+p)))return new Response("Forbidden",{status:403});
 const origin=new URL(config.origin).origin, nonce=randomBytes(24).toString("base64");
 const asset=origin+"/api/document-preview/assets/";
 const policy=`sandbox allow-scripts; default-src 'none'; script-src 'nonce-${nonce}' ${asset} blob:; worker-src blob:; connect-src 'none'; img-src data: blob:; style-src 'unsafe-inline'; font-src 'none'; object-src 'none'; base-uri 'none'; form-action 'none'; frame-ancestors ${origin}`;
 const html=`<!doctype html><html lang="es"><meta charset="utf-8"><meta name="referrer" content="no-referrer"><meta name="viewport" content="width=device-width, initial-scale=1"><title>Vista del documento</title><style>body{margin:0;padding:12px;background:#e7e7e7;color:#222;font:14px system-ui}canvas{display:block;max-width:100%;height:auto;margin:0 auto 12px;background:white}p{margin:0 0 12px}label{display:flex;gap:8px;align-items:center;margin-bottom:12px}select{font:inherit;padding:4px}main{overflow:auto}</style><body><label>Zoom <select id="zoom" aria-label="Zoom"><option value="fit">Ajustar</option><option value="1">100%</option><option value="1.5">150%</option><option value="2">200%</option></select></label><p id="status" role="status">Preparando vista…</p><main id="pages"></main><script type="module" nonce="${nonce}">
import{getDocument,GlobalWorkerOptions}from ${JSON.stringify(asset+"pdf-6.3.289.min.mjs")};
const origin=${JSON.stringify(origin)},nonce=location.hash.slice(1),status=document.getElementById('status'),pages=document.getElementById('pages');
GlobalWorkerOptions.workerSrc=${JSON.stringify(asset+"pdf.worker-6.3.289.min.mjs")};
const zoom=document.getElementById("zoom");const applyZoom=()=>{for(const canvas of pages.querySelectorAll("canvas")){canvas.style.maxWidth=zoom.value==="fit"?"100%":"none";canvas.style.width=zoom.value==="fit"?"":Math.ceil(canvas.width*Number(zoom.value))+"px";}};zoom.addEventListener("change",applyZoom);
let started=false,task;
const notify=(kind,extra={})=>parent.postMessage({kind,nonce,...extra},origin);
if(!/^[a-f0-9-]{36}$/.test(nonce))throw Error('viewer binding');
addEventListener('message',async event=>{
 if(event.source!==parent||event.origin!==origin||event.data?.nonce!==nonce||event.data?.kind!=='document-bytes'||started)return;
 started=true;const bytes=event.data.bytes;
 if(!(bytes instanceof ArrayBuffer)||bytes.byteLength<5||bytes.byteLength>2097152){status.textContent='Vista no disponible.';notify('rejected');return;}
 const timer=setTimeout(()=>{task?.destroy();pages.replaceChildren();status.textContent='La vista tardó demasiado. Descargá el original.';notify('rejected');},15000);
 try{
  if(new TextDecoder().decode(new Uint8Array(bytes,0,5))!=='%PDF-')throw Error('PDF required');
  // No viewer scripting manager or annotation layer is imported or instantiated.
  task=getDocument({data:new Uint8Array(bytes),isEvalSupported:false,enableScripting:false,enableXfa:false,useSystemFonts:false,disableFontFace:true,useWasm:false,isOffscreenCanvasSupported:false,isImageDecoderSupported:false,stopAtErrors:true,maxImageSize:16777216,disableAutoFetch:true});
  const pdf=await task.promise;
  if(pdf.numPages<1||pdf.numPages>30)throw Error('page limit');
  let pixels=0;
  for(let i=1;i<=pdf.numPages;i++){
   const page=await pdf.getPage(i),native=page.getViewport({scale:1}),scale=Math.min(1.4,1000/native.width),viewport=page.getViewport({scale});
   const width=Math.ceil(viewport.width),height=Math.ceil(viewport.height);pixels+=width*height;
   if(!Number.isFinite(pixels)||width<1||height<1||width>4096||height>4096||pixels>16777216)throw Error('pixel limit');
   const canvas=document.createElement('canvas');canvas.width=width;canvas.height=height;canvas.setAttribute('role','img');canvas.setAttribute('aria-label','Página '+i);pages.append(canvas);
   await page.render({canvas,canvasContext:canvas.getContext('2d'),viewport,annotationMode:0}).promise;
   page.cleanup();applyZoom();
  }
  status.textContent=pdf.numPages+' página(s)';notify('rendered',{pages:pdf.numPages,pixels});await task.destroy();
 }catch{pages.replaceChildren();status.textContent='No pudimos mostrar este archivo. Descargá el original.';notify('rejected');await task?.destroy().catch(()=>{});}
 finally{clearTimeout(timer);}
});
notify('ready');
</script></body></html>`;
 return new Response(html,{headers:{"content-type":"text/html; charset=utf-8","content-security-policy":policy,"cache-control":"no-store","x-content-type-options":"nosniff","x-frame-options":"SAMEORIGIN","referrer-policy":"no-referrer","permissions-policy":"camera=(), microphone=(), geolocation=(), fullscreen=(), payment=()"}});
}
