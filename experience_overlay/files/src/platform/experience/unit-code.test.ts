import { it, expect } from "vitest";
import { matchAuthorizedUnitCode, type AuthorizedCodeUnit } from "./unit-code";
const unit: AuthorizedCodeUnit = { id: "unit", serial_number: "001aB", product: "Model", organization: "Branch", selectable: true };
it("keeps leading zeros and case literal", () => { expect(matchAuthorizedUnitCode("001aB", [unit], true)).toEqual({ state: "MATCH", unit }); for (const code of ["1aB", "001ab", "001AB"])
    expect(matchAuthorizedUnitCode(code, [unit], true).state).toBe("NONE"); });
it.each(["", " x", "x ", "x\n", "\t", "é".repeat(65), "x".repeat(129), "\ud800"])("rejects invalid literal %j", code => expect(matchAuthorizedUnitCode(code, [unit], true).state).toBe("INVALID"));
it("never resolves partial, ambiguous or unavailable data", () => { expect(matchAuthorizedUnitCode("001aB", [unit], false).state).toBe("PARTIAL"); expect(matchAuthorizedUnitCode("001aB", [unit, { ...unit, id: "other" }], true).state).toBe("AMBIGUOUS"); expect(matchAuthorizedUnitCode("001aB", [{ ...unit, selectable: false }], true).state).toBe("UNAVAILABLE"); });
it("preserves allowed 128 UTF-8 bytes without coercion", () => { const unicode = { ...unit, serial_number: "é".repeat(64) }; expect(matchAuthorizedUnitCode(unicode.serial_number, [unicode], true)).toEqual({ state: "MATCH", unit: unicode }); });

import { checklistResponses, type PublishedChecklist } from "./checklists";
it("serial checklist binding preserves the literal code and rejects normalization",()=>{const checklist:PublishedChecklist={id:"delivery",organization_id:"store",version:1,state:"published",title:"Delivery",items:[{id:"serial",prompt:"Serial",response_type:"serial",required:true}]};const data=new FormData();data.set("answer:serial","001aB");expect(checklistResponses(checklist,data)).toEqual([{item_id:"serial",response_text:"001aB"}]);data.set("answer:serial"," 001aB");expect(()=>checklistResponses(checklist,data)).toThrow("Invalid literal serial");});
