import { describe, expect, it } from "vitest";
import { parseSubmission, summarySchema } from "./contract";
describe("survey submission grammar", () => {
 it("preserves zero and explicit consent", () => {
  expect(parseSubmission('{"score":0,"consent":true,"consent_version":"v1"}')).toEqual({score:0,consent:true,consent_version:"v1"});
 });
 it.each([
  '{}', '[]', '{"score":0,"consent":true}', '{"score":null,"consent":true,"consent_version":"v1"}',
  '{"score":1e0,"consent":true,"consent_version":"v1"}', '{"score":1.5,"consent":true,"consent_version":"v1"}',
  '{"score":11,"consent":true,"consent_version":"v1"}', '{"score":0,"consent":false,"consent_version":"v1"}',
  '{"score":0,"score":9,"consent":true,"consent_version":"v1"}',
  '{"score":0,"scor\\u0065":9,"consent":true,"consent_version":"v1"}',
  '{"score":0,"consent":true,"consent_version":"v1","customer_id":"other"}',
  '{"score":0,"consent":true,"consent_version":"v1",}',
  '{"score":0,"consent":true,"consent_version":"v1"}{}',
  '{"__proto__":{},"score":0,"consent":true,"consent_version":"v1"}',
  ' '.repeat(2049)
 ])("rejects ambiguous or invalid body: %s", text => { expect(() => parseSubmission(text)).toThrow(); });
 it("does not show a score below the configured threshold", () => {
  expect(summarySchema.safeParse({responses:1,available:false,nps:null}).success).toBe(true);
  expect(summarySchema.safeParse({responses:1,available:false,nps:100}).success).toBe(false);
  expect(summarySchema.safeParse({responses:0,available:true,nps:0}).success).toBe(false);
 });
});
