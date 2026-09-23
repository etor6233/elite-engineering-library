import { listModels } from "@/platform/backend/public-client";
import { errorResponse } from "@/platform/http/problem";
export const dynamic = "force-dynamic";
export async function GET() { try { return Response.json({ models: await listModels() }, { headers: { "cache-control": "public, max-age=30, stale-while-revalidate=120" } }); } catch (error) { return errorResponse(error); } }
