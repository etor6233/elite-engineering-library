import {expect,it} from "vitest";
import {operationMarkerSchema,retainOperation,currentReleaseSchema} from "./contracts";
const key="a03c7850-13ad-4b38-81d0-120000000001";
it("retains one key and rejects corrupt/cross-order recovery without overwriting",()=>{
 const saved=new Map<string,string>(),storage={getItem:(k:string)=>saved.get(k)??null,setItem:(k:string,v:string)=>{saved.set(k,v)}};
 const old=retainOperation(storage,"scope",{action:"prepare",orderId:"order",requestKey:key});
 expect(retainOperation(storage,"scope",{action:"prepare",orderId:"order",requestKey:"a03c7850-13ad-4b38-81d0-120000000002"})).toEqual(old);
 expect(()=>retainOperation(storage,"scope",{action:"prepare",orderId:"other",requestKey:key})).toThrow();expect(JSON.parse(saved.get("scope")!)).toEqual(old);
 saved.set("scope","bad-json");expect(()=>retainOperation(storage,"scope",old)).toThrow();expect(saved.get("scope")).toBe("bad-json");
});
it("requires a handover reference for release recovery even if current context is unavailable",()=>{expect(operationMarkerSchema.safeParse({action:"release",orderId:"order",requestKey:key}).success).toBe(false);expect(operationMarkerSchema.safeParse({action:"release",orderId:"order",handoverId:"handover",requestKey:key}).success).toBe(true)});
it("storage failure prevents creation of an unsafe operation reference",()=>{expect(()=>retainOperation({getItem:()=>null,setItem:()=>{throw new Error("disabled")}},"scope",{action:"prepare",orderId:"order",requestKey:key})).toThrow()});
it("a historical receipt cannot manufacture a current positive result",()=>{expect(currentReleaseSchema.safeParse({receipt:{id:"receipt"},current:true}).success).toBe(false)});
