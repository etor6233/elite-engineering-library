import {readSession} from "@/platform/auth/session";
import {createCaptureGateway} from "@/platform/capture-vnext/gateway";
import {captureConfig,captureTransport} from "@/platform/capture-vnext/server";
export const runtime="nodejs";
export const POST=createCaptureGateway({session:readSession,config:captureConfig,transport:captureTransport});
