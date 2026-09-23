import { readSession } from "@/platform/auth/session";
import { createDocumentGateway } from "@/platform/documents-v403/gateway";
import { documentConfig, documentTransport } from "@/platform/documents-v403/server";
export const runtime = "nodejs";
const handle = createDocumentGateway({ session: readSession, config: documentConfig, transport: documentTransport });
export const GET = handle;
export const POST = handle;
export const PUT = handle;
