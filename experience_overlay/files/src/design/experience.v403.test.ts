import { describe, it, expect } from "vitest";
import { readFileSync } from "node:fs";
import { contrast } from "./contrast.v403";
import tokens from "./tokens.v403.json";
const pairs = [["ink", "panel"], ["ink", "surface"], ["muted", "panel"], ["muted", "surface"], ["muted", "muted-bg"], ["on-accent", "accent"], ["on-accent", "accent-strong"], ["accent", "panel"], ["accent", "surface"], ["success", "success-bg"], ["warning", "warning-bg"], ["danger", "danger-bg"], ["ink", "muted-bg"], ["ink", "ink-soft"], ["warning", "panel"]];
describe("V403 actual text pairs", () => {
    for (const [brand, delta] of Object.entries(tokens.brands)) {
        const t = { ...tokens.default, ...delta } as Record<string, string>;
        for (const [a, b] of pairs)
            it(`${brand}: ${a}/${b} >=4.5`, () => expect(contrast(t[a!]!, t[b!]!)).toBeGreaterThanOrEqual(4.5));
    }
    it("retains the exact historical defect as a regression oracle", () => { expect(contrast("#dc2626", "#fee2e2")).toBeCloseTo(3.9534, 3); expect(contrast(tokens.default.danger, tokens.default["danger-bg"])).toBeGreaterThan(4.5); });
    it("the generated page stylesheet contains the selected tokens", () => {
        const css = readFileSync("src/app/globals.css", "utf8");
        for (const [k, v] of Object.entries(tokens.default))
            expect(css).toContain(`--${k}: ${v};`);
        expect(css).not.toContain("#DC2626");
    });
});
// WCAG non-text controls/focus use 3:1, independently from the 4.5:1 text rule.
describe("V403 selected control and focus boundaries", () => { for (const [brand, delta] of Object.entries(tokens.brands)) {
    const values = { ...tokens.default, ...delta } as Record<string, string>;
    for (const token of ["line", "focus"])
        for (const background of ["surface", "panel"])
            it(`${brand} ${token}/${background} >=3`, () => expect(contrast(values[token]!, values[background]!)).toBeGreaterThanOrEqual(3));
} });
