import { permittedOptions } from "@/platform/experience/options";
import { NextResponse, type NextRequest } from "next/server";
import { readFile } from "node:fs/promises";
import { join } from "node:path";
import { allowed, readSession } from "@/platform/auth/session";
import { protectedGet } from "@/platform/backend/protected-client";
import { BackendProblem } from "@/platform/backend/public-client";
import { checklistSchema, checklistSelectionSchema } from "@/platform/experience/checklists";
// Only configured references are enumerated; each record is fetched under the current principal.
export async function GET(request: NextRequest) {
    const session = await readSession();
    if (!session)
        return NextResponse.json({ code: "UNAUTHENTICATED" }, { status: 401 });
    const query = new URL(request.url).searchParams, organization = query.get("organizationId");
    if (!organization || organization.length > 128 || [...query.keys()].join(",") !== "organizationId")
        return NextResponse.json({ code: "INVALID_QUERY" }, { status: 400 });
    if (!session.organizations.includes(organization) || !allowed(session, "handover:manage"))
        return NextResponse.json({ code: "FORBIDDEN" }, { status: 403 });
    try {
        const config = checklistSelectionSchema.parse(JSON.parse(await readFile(join(process.cwd(), "config", "experience-checklists.json"), "utf8")));
        const references = config.tenants[session.tenantId]?.organizations[organization] ?? [], checklists = [];
        for (const reference of references) {
            const value = checklistSchema.parse(await protectedGet<unknown>(session, "/v1/franchise/delivery-checklists/result", { organization_id: organization, checklist_id: reference.id, version: String(reference.version) }));
            if (value.organization_id !== organization || value.id !== reference.id || value.version !== reference.version)
                throw new Error("scope mismatch");
            checklists.push(value);
        }
        return NextResponse.json({ checklists, evidenceOptions: checklists.some(c => c.items.some(item => item.response_type === "evidence")) ? await permittedOptions(session, organization, "evidence") : [] }, { headers: { "cache-control": "no-store" } });
    }
    catch (error) {
        if (error instanceof BackendProblem)
            return NextResponse.json({ code: error.code }, { status: error.status, headers: { "cache-control": "no-store" } });
        return NextResponse.json({ code: "CHECKLIST_CONFIGURATION_OR_UPSTREAM_UNAVAILABLE" }, { status: 502, headers: { "cache-control": "no-store" } });
    }
}
