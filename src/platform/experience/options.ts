import "server-only";
import { readFile } from "node:fs/promises";
import { join } from "node:path";
import { z } from "zod";
import { allowed, type PortalSession } from "@/platform/auth/session";
import {protectedGet} from "@/platform/backend/protected-client";
// AUTHORED config/authorization glue. This selection catalogue grants no business effect.
// Its deployment-reviewed references are checked again by original Go command owners.
export const groups = { variants: "quote:write", priceBooks: "quote:write", assignees: "lead:assign", resources: "availability:manage", availabilityReasons: "availability:manage", appointmentReasons: "appointment:manage", evidence: "handover:manage", "catalogModels": "catalog:draft|catalog:publish", "catalogMedia": "catalog:draft|catalog:publish", "catalogBooks": "catalog:draft|catalog:publish", "catalogDrafts": "catalog:draft|catalog:publish", "customerReasons": "customer:self", "resourcePrincipals": "resource:manage", "helpOrganizations": "help:write|help:publish", "networkOrganizations": "network:admin|franchise:write|network:read|franchise:read", "networkTerritories": "franchise:write", "networkTerms": "franchise:write", "networkEntities": "network:read|franchise:read|network:admin|franchise:write", "metricReferences": "surveys:read|stored_value:read", "storedOrders": "stored_value:read", "storedAccounts": "stored_value:request", "storedOperations": "stored_value:request", "supplySuppliers": "supply:plan", "supplyFactories": "supply:plan", "supplyVariants": "supply:plan", "supplyOrders": "supply:read|supply:factory-read|supply:plan|supply:factory", "warrantyQuotes": "warranty:offer|warranty:self", "warrantyHandovers": "warranty:activate|warranty:request|warranty:self", "warrantyAppointments": "warranty:request|warranty:self", "warrantyClaims": "warranty:read|warranty:factory-read|warranty:self", "warrantyFaults": "warranty:diagnose", "warrantyParts": "warranty:plan", "warrantyBins": "warranty:plan", "warrantyLots": "warranty:plan", "warrantyReceipts": "warranty:plan" } as const;
export type OptionGroup = keyof typeof groups;
export function groupAllowed(session: PortalSession, group: OptionGroup) { return groups[group].split("|").some(permission => allowed(session, permission)); }
const option = z.object({ value: z.string().min(1).max(256), label: z.string().trim().min(1).max(160), allowed_subjects: z.array(z.string().min(1).max(256)).min(1).max(100) }).strict();
export const optionsSchema = z.object({ schema_version: z.literal("1"), tenants: z.record(z.string(), z.object({ organization_labels: z.record(z.string(), z.string().trim().min(1).max(160)).optional(), organizations: z.record(z.string(), z.partialRecord(z.enum(Object.keys(groups) as [
            OptionGroup,
            ...OptionGroup[]
        ]), z.array(option).max(100))) })) }).strict();
export async function loadOptions() { return optionsSchema.parse(JSON.parse(await readFile(join(process.cwd(), "config", "experience-options.json"), "utf8"))); }
export async function permittedOptions(session: PortalSession, organization: string, group: OptionGroup) {
    if (!session.organizations.includes(organization) || !groupAllowed(session, group))
        throw new Error("FORBIDDEN");
    if(group==="networkOrganizations"||group==="networkEntities"){
        const result=await protectedGet(session,"/v1/franchise/network/options/"+(group==="networkOrganizations"?"organizations":"entities"),{});
        return z.array(z.object({value:z.string().min(1).max(128),label:z.string().min(1).max(260)}).strict()).max(100).parse(result);
    }
    const config = await loadOptions(), rows = config.tenants[session.tenantId]?.organizations[organization]?.[group] ?? [];
    return rows.filter(row => row.allowed_subjects.includes(session.subject)).map(({ value, label }) => ({ value, label }));
}
export async function validateGuidedSelection(command: Record<string, unknown>, session: PortalSession) {
    const needs: Partial<Record<string, [
        OptionGroup,
        string
    ][]>> = { "create-quote": [["variants", "variantId"], ["priceBooks", "priceBookId"]], "assign-lead": [["assignees", "assignedSubject"]], "create-availability": [["resources", "resourceId"], ["availabilityReasons", "reasonCode"]], "transition-appointment": [["appointmentReasons", "reasonCode"]] };
    for (const [group, key] of needs[String(command.action)] ?? []) {
        const value = command[key];
        if ((value === undefined || value === "") && ["resourceId", "reasonCode"].includes(key))
            continue;
        const options = await permittedOptions(session, String(command.organizationId), group);
        if (!options.some(option => option.value === value))
            return false;
    }
    if (command.action === "complete-delivery-checklist" && Array.isArray(command.responses)) {
        const references = command.responses.filter((r: {
            evidence_sha256?: string;
        }) => Boolean(r.evidence_sha256));
        if (references.length) {
            const allowed = await permittedOptions(session, String(command.organizationId), "evidence");
            if (references.some((r: {
                evidence_sha256: string;
            }) => !allowed.some(x => x.value === r.evidence_sha256)))
                return false;
        }
    }
    return true;
}
export async function organizationLabel(session: PortalSession, organization: string) {
    if (!session.organizations.includes(organization))
        throw new Error("FORBIDDEN");
    try {
        return (await loadOptions()).tenants[session.tenantId]?.organization_labels?.[organization] ?? null;
    }
    catch {
        return null;
    }
}
