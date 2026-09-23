// AUTHORED bounded reader separate from the Node-only BFF gateway.
export async function boundedDocumentBytes(response:Response,limit:number){
 if(!response.body)throw Error("Missing original");const reader=response.body.getReader();const chunks:Uint8Array[]=[];let size=0;
 try{for(;;){const part=await reader.read();if(part.done)break;size+=part.value.length;if(size>limit)throw Error("Original too large");chunks.push(part.value)}
 const bytes=new Uint8Array(size);let offset=0;for(const chunk of chunks){bytes.set(chunk,offset);offset+=chunk.length}return bytes;
 }finally{void reader.cancel().catch(()=>{});reader.releaseLock()}
}
