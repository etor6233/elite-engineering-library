import { supplyPlan } from "@/platform/supply/contract";
export type AuthorizedCodeUnit = {
    id: string;
    serial_number: string;
    product: string;
    organization: string;
    selectable: boolean;
};
export type CodeMatch = {
    state: "INVALID" | "PARTIAL" | "NONE" | "AMBIGUOUS" | "UNAVAILABLE";
} | {
    state: "MATCH";
    unit: AuthorizedCodeUnit;
};
// AUTHORED literal lookup, using the admitted supply owner's exact serial validator.
// Input is never normalized, executed, interpreted as a URL, or used to expand authorization.
export function validLiteralUnitCode(code: string): boolean {
    return supplyPlan.shape.units.element.shape.serial_number.safeParse(code).success;
}
export function matchAuthorizedUnitCode(code: string, units: AuthorizedCodeUnit[], complete: boolean): CodeMatch {
    if (!validLiteralUnitCode(code))
        return { state: "INVALID" };
    if (!complete)
        return { state: "PARTIAL" };
    const matches = units.filter(unit => unit.serial_number === code);
    if (!matches.length)
        return { state: "NONE" };
    if (matches.length !== 1)
        return { state: "AMBIGUOUS" };
    if (!matches[0]!.selectable)
        return { state: "UNAVAILABLE" };
    return { state: "MATCH", unit: matches[0]! };
}
