// Public parser code only; no original document, token or tenant state is served.
import {readFile} from "node:fs/promises";
import {createHash} from "node:crypto";
import path from "node:path";
import lock from "@/platform/documents-vnext/pdfjs-lock.json";
export const runtime="nodejs";
const files:Record<string,string>={"pdf-6.3.289.min.mjs":"pdf.min.mjs","pdf.worker-6.3.289.min.mjs":"pdf.worker.min.mjs"};
export async function GET(_request:Request,{params}:{params:Promise<{asset:string}>}){
 const name=files[(await params).asset];if(!name)return new Response("Not found",{status:404});
 const bytes=await readFile(path.join(process.cwd(),"third_party/pdfjs",name));
 if(createHash("sha256").update(bytes).digest("hex")!==lock.files[name as keyof typeof lock.files])return new Response("Unavailable",{status:503});
 return new Response(bytes,{headers:{"content-type":"text/javascript; charset=utf-8","cache-control":"public, max-age=31536000, immutable","x-content-type-options":"nosniff","access-control-allow-origin":"*","cross-origin-resource-policy":"cross-origin"}});
}
