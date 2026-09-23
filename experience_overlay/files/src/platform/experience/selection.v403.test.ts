import { beforeEach, describe, it, expect, vi } from "vitest";
import { NextRequest } from "next/server";
const m = vi.hoisted(() => ({ read: vi.fn(), session: vi.fn(), get: vi.fn() }));
vi.mock("node:fs/promises", () => ({ readFile: m.read }));
vi.mock("@/platform/auth/session", () => ({ readSession: m.session, allowed: (s: {
        permissions: string[];
    }, p: string) => s.permissions.includes(p) }));
vi.mock("@/platform/backend/protected-client", () => ({ protectedGet: m.get }));
import { GET } from "@/app/api/enterprise/experience-options/route";
import { GET as checklistsGET } from "@/app/api/enterprise/experience-checklists/route";
import { permittedOptions, validateGuidedSelection, groups as allGroups, organizationLabel, type OptionGroup } from "./options";
import { checklistResponses } from "./checklists";
const s = { subject: "operator", tenantId: "tenant", organizations: ["store"], permissions: ["quote:write", "lead:assign", "availability:manage", "handover:manage"], accessToken: "fixture" };
const groups = { variants: [{ value: "variant", label: "Producto Uno", allowed_subjects: ["operator"] }, { value: "private", label: "Restringido", allowed_subjects: ["other"] }], priceBooks: [{ value: "retail", label: "Minorista", allowed_subjects: ["operator"] }], assignees: [], resources: [], availabilityReasons: [], appointmentReasons: [], evidence: [] };
const checklist = { id: "delivery", organization_id: "store", version: 1, state: "published" as const, title: "Entrega", items: [{ id: "serial", prompt: "Serie", response_type: "serial" as const, required: true }, { id: "condition", prompt: "Condición", response_type: "confirmation" as const, required: true }] };
beforeEach(() => { vi.clearAllMocks(); m.session.mockResolvedValue(s); m.read.mockResolvedValue(JSON.stringify({ schema_version: "1", tenants: { tenant: { organizations: { store: groups } } } })); });
describe("server selection boundary", () => {
    it("returns permitted named choices only", async () => { expect(await permittedOptions(s, "store", "variants")).toEqual([{ value: "variant", label: "Producto Uno" }]); const r = await GET(new NextRequest("https://example.invalid/api?organizationId=store&group=variants")); expect(r.status).toBe(200); expect(await r.json()).toEqual({ options: [{ value: "variant", label: "Producto Uno" }] }); });
    it("rejects cross organization without reading config", async () => { m.session.mockResolvedValue({ ...s, organizations: ["foreign"] }); expect((await GET(new NextRequest("https://x.invalid/?organizationId=store&group=variants"))).status).toBe(403); expect(m.read).not.toHaveBeenCalled(); });
    it("isolates same organization names by tenant", async () => expect(await permittedOptions({ ...s, tenantId: "different" }, "store", "variants")).toEqual([]));
    it("checks both quote references again at POST boundary", async () => { expect(await validateGuidedSelection({ action: "create-quote", organizationId: "store", variantId: "variant", priceBookId: "retail" }, s)).toBe(true); expect(await validateGuidedSelection({ action: "create-quote", organizationId: "store", variantId: "private", priceBookId: "retail" }, s)).toBe(false); });
    it("missing config fails explicitly rather than admitting arbitrary input", async () => { m.read.mockRejectedValue(new Error("missing")); expect((await GET(new NextRequest("https://x.invalid/?organizationId=store&group=variants"))).status).toBe(503); });
    it("unknown groups and duplicate queries are rejected", async () => {
        for (const query of ["organizationId=store&group=__proto__", "organizationId=store&group=variants&group=variants"])
            expect((await GET(new NextRequest("https://x.invalid/?" + query))).status).toBe(400);
    });
    it("absent subject and grant reject options", async () => { m.session.mockResolvedValue(null); expect((await GET(new NextRequest("https://x.invalid/?organizationId=store&group=variants"))).status).toBe(401); m.session.mockResolvedValue({ ...s, permissions: [] }); expect((await GET(new NextRequest("https://x.invalid/?organizationId=store&group=variants"))).status).toBe(403); });
});
describe("guided published checklist", () => {
    beforeEach(() => { m.read.mockResolvedValue(JSON.stringify({ schema_version: "1", tenants: { tenant: { organizations: { store: [{ id: "delivery", version: 1 }] } } } })); m.get.mockResolvedValue(checklist); });
    it("reads configured version using current principal", async () => { expect((await checklistsGET(new NextRequest("https://x.invalid/?organizationId=store"))).status).toBe(200); expect(m.get).toHaveBeenCalledWith(s, "/v1/franchise/delivery-checklists/result", { organization_id: "store", checklist_id: "delivery", version: "1" }); });
    it("rejects upstream organization or revision mismatches", async () => {
        for (const bad of [{ organization_id: "foreign" }, { version: 2 }, { id: "different" }, { state: "draft" }]) {
            m.get.mockResolvedValue({ ...checklist, ...bad });
            expect((await checklistsGET(new NextRequest("https://x.invalid/?organizationId=store"))).status).toBe(502);
        }
    });
    it("keeps response IDs internal and requires confirmations", () => { const data = new FormData(); data.set("answer:serial", "SERIAL-123"); expect(() => checklistResponses(checklist, data)).toThrow(); data.set("answer:condition", "confirmed"); expect(checklistResponses(checklist, data)).toEqual([{ item_id: "serial", response_text: "SERIAL-123" }, { item_id: "condition", response_text: "confirmed" }]); });
});
import { definitiveSelectionRejection } from "./rejection";
it("does not release a fence for generic errors or unverifiable result", () => { expect(definitiveSelectionRejection({ status: 403 }, { code: "FORBIDDEN" })).toBe(false); expect(definitiveSelectionRejection({ status: 503 }, { code: "SELECTION_CONFIGURATION_UNAVAILABLE" })).toBe(false); expect(definitiveSelectionRejection({ status: 403 }, { code: "SELECTION_NOT_AUTHORIZED", operation_effect: "NOT_ATTEMPTED" })).toBe(true); expect(definitiveSelectionRejection({ status: 503 }, { code: "SELECTION_CONFIGURATION_UNAVAILABLE", operation_effect: "NOT_ATTEMPTED" })).toBe(true); });
describe("every guided reference group boundary", () => {
    for (const [group, permissions] of Object.entries(allGroups)) {
        it(group + " isolates identity and requires one actual product grant", async () => {
            const rows = [{ value: "allowed-ref", label: "Nombre comercial aprobado", allowed_subjects: ["operator"] }, { value: "private-ref", label: "No divulgar", allowed_subjects: ["other"] }];
            m.read.mockResolvedValue(JSON.stringify({ schema_version: "1", tenants: { tenant: { organizations: { store: { [group]: rows } } } } }));
            for (const permission of permissions.split("|")) {
                const session = { ...s, permissions: [permission] };
                expect(await permittedOptions(session, "store", group as OptionGroup)).toEqual([{ value: "allowed-ref", label: "Nombre comercial aprobado" }]);
                expect(await permittedOptions({ ...session, subject: "other-subject" }, "store", group as OptionGroup)).toEqual([]);
            }
            await expect(permittedOptions({ ...s, permissions: [] }, "store", group as OptionGroup)).rejects.toThrow("FORBIDDEN");
        });
    }
});
it("does not reveal another tenant or organization branch label", async () => { m.read.mockResolvedValue(JSON.stringify({ schema_version: "1", tenants: { tenant: { organization_labels: { store: "Sucursal Centro" }, organizations: {} } } })); expect(await organizationLabel(s, "store")).toBe("Sucursal Centro"); expect(await organizationLabel({ ...s, tenantId: "other" }, "store")).toBeNull(); await expect(organizationLabel({ ...s, organizations: ["other"] }, "store")).rejects.toThrow("FORBIDDEN"); });
