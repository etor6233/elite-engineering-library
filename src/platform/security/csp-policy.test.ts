import { describe, expect, it } from "vitest";
import { CspEvaluator } from "csp_evaluator/dist/evaluator.js";
import { CspParser } from "csp_evaluator/dist/parser.js";
import { Severity } from "csp_evaluator/dist/finding.js";
import { buildContentSecurityPolicy } from "./csp-policy";

const nonce = "dGVzdC1ub25jZS0xMjM0NTY=";

describe("strict Content Security Policy", () => {
  it("has no actionable finding in Google's CSP Evaluator", () => {
    const policy = buildContentSecurityPolicy(nonce, false);
    const findings = new CspEvaluator(new CspParser(policy).csp).evaluate();
    expect(findings.filter((finding) => finding.severity < Severity.NONE)).toEqual([]);
    expect(policy).not.toContain("'unsafe-inline'");
    expect(policy).not.toContain("'unsafe-eval'");
  });

  it("keeps unsafe-eval restricted to development", () => {
    expect(buildContentSecurityPolicy(nonce, true)).toContain("'unsafe-eval'");
    expect(buildContentSecurityPolicy(nonce, false)).not.toContain("'unsafe-eval'");
  });

  it("fails closed for malformed or short nonces", () => {
    expect(() => buildContentSecurityPolicy("short", false)).toThrow(/nonce/);
    expect(() => buildContentSecurityPolicy("not-a-valid-base64-token!", false)).toThrow(/nonce/);
  });

  it("proves the previous header was not acceptable", () => {
    const previous = "base-uri 'self'; object-src 'none'; frame-ancestors 'none'; form-action 'self'; upgrade-insecure-requests";
    const findings = new CspEvaluator(new CspParser(previous).csp).evaluate();
    expect(findings.some((finding) => finding.severity === Severity.HIGH && finding.directive === "script-src")).toBe(true);
  });
});
