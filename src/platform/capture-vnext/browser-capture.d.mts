export class CaptureError extends Error {code:string;constructor(code:string)}
export const limits:Readonly<{fileBytes:number;pixels:number;dimension:number;textBytes:number}>;
export function frameFromFile(file:File,options?:{signal?:AbortSignal}):Promise<Uint8Array>;
export function frameFromVideo(video:HTMLVideoElement):Uint8Array;
export function frameFromRGBA(bytes:Uint8ClampedArray,width:number,height:number):Uint8Array;
export function startCamera(video:HTMLVideoElement,options?:{signal?:AbortSignal}):Promise<{stream:MediaStream;stop:()=>void}>;
export function readerText(value:string):string;
export function captureText(value:string):string;
export function selectedCandidate(result:unknown,index:number):string;
