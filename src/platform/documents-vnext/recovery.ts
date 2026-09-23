// AUTHORED recovery metadata. Never persist original bytes, names, field values or decision text.
import {hash,pendingSchema,type PendingDocument} from "./contract";
const prefix=(scope:string)=>`elite-documents-vnext:${hash.parse(scope)}:`;
export function readMarkers(scope:string):PendingDocument[]{
 const key=prefix(scope),markers:PendingDocument[]=[];
 for(let i=0;i<sessionStorage.length;i++){const name=sessionStorage.key(i);if(!name?.startsWith(key))continue;const raw=sessionStorage.getItem(name);if(!raw)continue;const marker=pendingSchema.parse(JSON.parse(raw));if(name!==key+marker.id)throw new Error("Invalid recovery binding");markers.push(marker);if(markers.length>50)throw new Error("Too many unresolved documents")}
 return markers;
}
export function readMarker(scope:string,id:string){return readMarkers(scope).find(marker=>marker.id===id)??null}
export function retainMarker(scope:string,marker:PendingDocument){pendingSchema.parse(marker);const previous=readMarker(scope,marker.id);if(previous&&JSON.stringify(Object.entries(previous).sort())!==JSON.stringify(Object.entries(marker).sort()))throw new Error("Different request pending");sessionStorage.setItem(prefix(scope)+marker.id,JSON.stringify(marker));}
export function clearMarker(scope:string,id:string){sessionStorage.removeItem(prefix(scope)+id)}
