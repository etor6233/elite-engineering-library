import { NextResponse } from "next/server";
import { readSession } from "@/platform/auth/session";
import { loadBusinessConfig } from "@/platform/config/load";
import { availableGuides, parseHelpQuery, selectGuides } from "@/platform/help/catalog";
import { helpResponseSchema } from "@/platform/help/contract";

const reply = (value: unknown, status = 200) => NextResponse.json(value, { status, headers: { "Cache-Control": "no-store", "X-Content-Type-Options": "nosniff" } });
export async function GET(request: Request) {
  if (["cross-site", "none"].includes(request.headers.get("sec-fetch-site") ?? "")) return reply({ code: "CROSS_SITE_REJECTED" }, 403);
  const query = parseHelpQuery(request.url);
  if (!query) return reply({ code: "INVALID_QUERY" }, 400);
  try {
    const session = await readSession();
    if (!session) return reply({ code: "UNAUTHENTICATED" }, 401);
    const config = await loadBusinessConfig();
    const guides = availableGuides(session, config.features.whatsapp_status_history === true);
    if (!guides.length) return reply({ code: "FORBIDDEN" }, 403);
    const articles = selectGuides(guides, query);
    // An unavailable version and an inapplicable article share the same response.
    if (query.article !== null && !articles.length) return reply({ code: "GUIDE_NOT_AVAILABLE" }, 404);
    return reply(helpResponseSchema.parse({ schema: "journey-help/1", articles }));
  } catch { return reply({ code: "HELP_UNAVAILABLE" }, 503); }
}
