import {readSession} from "@/platform/auth/session";
import {createDocumentGateway} from "@/platform/documents-vnext/gateway";
import {documentConfig,documentTransport} from "@/platform/documents-vnext/server";
export const runtime="nodejs";
const handle=createDocumentGateway({session:readSession,config:documentConfig,transport:documentTransport});
export const GET=handle;export const PUT=handle;export const POST=handle;
