import { NextResponse, type NextRequest } from "next/server";
import { readSession } from "@/platform/auth/session";
import { groups, groupAllowed, permittedOptions, type OptionGroup } from "@/platform/experience/options";
export async function GET(request: NextRequest) {
    const session = await readSession();
    if (!session)
        return NextResponse.json({ code: "UNAUTHENTICATED" }, { status: 401 });
    const q = new URL(request.url).searchParams, organization = q.get("organizationId"), group = q.get("group");
    if (!organization || organization.length > 128 || !group || !Object.hasOwn(groups, group) || [...q.keys()].sort().join(",") !== "group,organizationId")
        return NextResponse.json({ code: "INVALID_QUERY" }, { status: 400 });
    if (!session.organizations.includes(organization) || !groupAllowed(session, group as OptionGroup))
        return NextResponse.json({ code: "FORBIDDEN" }, { status: 403 });
    try {
        return NextResponse.json({ options: await permittedOptions(session, organization, group as OptionGroup) }, { headers: { "cache-control": "no-store" } });
    }
    catch {
        return NextResponse.json({ code: "SELECTION_CONFIGURATION_UNAVAILABLE" }, { status: 503, headers: { "cache-control": "no-store" } });
    }
}
