import { NextResponse } from "next/server";
import { readSession } from "@/platform/auth/session";

export async function GET() {
  const session = await readSession();
  if (!session) return NextResponse.json({ authenticated: false }, { status: 401, headers: { "cache-control": "no-store" } });
  return NextResponse.json({ authenticated: true, subject: session.subject, tenant_id: session.tenantId, permissions: session.permissions, organization_ids: session.organizations }, { headers: { "cache-control": "no-store" } });
}
