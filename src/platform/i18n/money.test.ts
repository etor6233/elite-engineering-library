import { expect, test } from "vitest";
import { minorAmountPresentation } from "./money";
test.each([
 {minor:100,currency:"USD",label:"USD 1,00"},
 {minor:0,currency:"USD",label:"USD 0,00"},
 {minor:123,currency:"JPY",label:"JPY 123"},
 {minor:1234,currency:"KWD",label:"KWD 1,234"},
 {minor:9007199254740991,currency:"USD",label:"USD 90.071.992.547.409,91"},
 {minor:100,currency:"CLP",label:"CLP 100"},
])("exact existing quote representation: $currency $minor",({minor,currency,label})=>{const result=minorAmountPresentation(minor,currency);expect(result.amountLabel.replace(/\u00a0/g," ")).toBe(label);expect(result.amountValid).toBe(true);});
test.each([
 {minor:-1,currency:"USD"},{minor:1.5,currency:"USD"},{minor:NaN,currency:"USD"},
 {minor:Infinity,currency:"USD"},{minor:9007199254740992,currency:"USD"},
 {minor:100,currency:"unknown"},{minor:100,currency:"usd"},
])("invalid display cannot imply a verified amount: $currency $minor",({minor,currency})=>{expect(minorAmountPresentation(minor,currency)).toEqual({amountLabel:"Importe no verificable; consultá al soporte.",amountValid:false});});
