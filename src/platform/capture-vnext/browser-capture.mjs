// AUTHORED browser glue using standard File/ImageBitmap/Canvas/MediaStream APIs.
// No vendor assets, decoder, QR authenticity or business permission is implied.
export const limits = Object.freeze({fileBytes:2*1024*1024,pixels:4*1024*1024,dimension:4096,textBytes:1024});
export class CaptureError extends Error { constructor(code){super(code);this.name='CaptureError';this.code=code;} }
function dimensions(width,height){if(!Number.isInteger(width)||!Number.isInteger(height)||width<=20||height<=20||width>limits.dimension||height>limits.dimension||width*height>limits.pixels)throw new CaptureError('IMAGE_DIMENSIONS');}
// Header checks bound decoded dimensions before asking the browser image codec.
function imageDimensions(bytes){
 const view=new DataView(bytes.buffer,bytes.byteOffset,bytes.byteLength);
 if(bytes.length>=24&&bytes.slice(0,8).every((v,i)=>v===[137,80,78,71,13,10,26,10][i]))return [view.getUint32(16),view.getUint32(20),'image/png'];
 if(bytes.length>4&&bytes[0]===255&&bytes[1]===216){
  let offset=2;
  while(offset+4<=bytes.length){
   if(bytes[offset++]!==255)throw new CaptureError('INVALID_IMAGE_HEADER');while(bytes[offset]===255)offset++;
   const marker=bytes[offset++];if(marker===218||marker===217)break;
   if(marker===1||(marker>=208&&marker<=215))continue;
   const length=view.getUint16(offset);if(length<2||offset+length>bytes.length)throw new CaptureError('INVALID_IMAGE_HEADER');
   if([192,193,194,195,197,198,199,201,202,203,205,206,207].includes(marker)){if(length<8)throw new CaptureError('INVALID_IMAGE_HEADER');return [view.getUint16(offset+5),view.getUint16(offset+3),'image/jpeg'];}offset+=length;
  }
 }
 throw new CaptureError('IMAGE_FORMAT');
}
export function frameFromRGBA(rgba,width,height){
 dimensions(width,height);if(!(rgba instanceof Uint8ClampedArray)||rgba.length!==width*height*4)throw new CaptureError('FRAME_FORMAT');
 const frame=new Uint8Array(8+width*height);frame.set([69,71,89,49]);const view=new DataView(frame.buffer);view.setUint16(4,width);view.setUint16(6,height);
 for(let i=0,j=8;i<rgba.length;i+=4,j++){const a=rgba[i+3];const luma=(77*rgba[i]+150*rgba[i+1]+29*rgba[i+2]+128)>>8;frame[j]=Math.round((luma*a+255*(255-a))/255);}return frame;
}
function raster(source,width,height){
 dimensions(width,height);const canvas=document.createElement('canvas');canvas.width=width;canvas.height=height;
 try{const context=canvas.getContext('2d',{willReadFrequently:true});if(!context)throw new CaptureError('CANVAS_UNAVAILABLE');context.drawImage(source,0,0);return frameFromRGBA(context.getImageData(0,0,width,height).data,width,height);}finally{canvas.width=0;canvas.height=0;}
}
export async function frameFromFile(file,{signal}={}){
 if(!(file instanceof Blob)||file.size===0||file.size>limits.fileBytes)throw new CaptureError('IMAGE_SIZE');signal?.throwIfAborted();
 const bytes=new Uint8Array(await file.arrayBuffer());const[width,height,mime]=imageDimensions(bytes);dimensions(width,height);
 if(file.type&&file.type!==mime)throw new CaptureError('IMAGE_MIME_MISMATCH');signal?.throwIfAborted();
 const bitmap=await createImageBitmap(new Blob([bytes],{type:mime}));try{signal?.throwIfAborted();dimensions(bitmap.width,bitmap.height);return raster(bitmap,bitmap.width,bitmap.height);}finally{bitmap.close();}
}
export function frameFromVideo(video){if(!(video instanceof HTMLVideoElement)||video.readyState<2)throw new CaptureError('CAMERA_NOT_READY');return raster(video,video.videoWidth,video.videoHeight);}
export async function startCamera(video,{signal}={}){
 if(!navigator.mediaDevices?.getUserMedia)throw new CaptureError('CAMERA_UNAVAILABLE');signal?.throwIfAborted();
 const stream=await navigator.mediaDevices.getUserMedia({audio:false,video:{facingMode:{ideal:'environment'},width:{ideal:1280,max:2048},height:{ideal:720,max:2048}}});
 const stop=()=>{for(const track of stream.getTracks())track.stop();if(video.srcObject===stream)video.srcObject=null;signal?.removeEventListener('abort',stop);};
 try{if(signal?.aborted){stop();signal.throwIfAborted();}signal?.addEventListener('abort',stop,{once:true});video.srcObject=stream;video.muted=true;video.playsInline=true;await video.play();return {stream,stop};}catch(error){stop();throw error;}
}
export function captureText(value){if(typeof value!=='string'||!value||new TextEncoder().encode(value).length>limits.textBytes||/[\u0000-\u001f\u007f]/.test(value))throw new CaptureError('TEXT_FORMAT');return value;}
// Bind to a focused dedicated input; never install global key interception.
export function readerText(inputValue){return captureText(inputValue);}
// A MULTIPLE result requires user selection. Every selected text is still
// untrusted and must be resolved through the authenticated server session.
export function selectedCandidate(result,index){if(!result||!['DECODED','MULTIPLE'].includes(result.status)||!Array.isArray(result.candidates)||!Number.isInteger(index)||!result.candidates[index])throw new CaptureError('SELECTION_REQUIRED');return captureText(result.candidates[index].text);}
