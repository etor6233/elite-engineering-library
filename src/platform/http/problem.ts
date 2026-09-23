import { NextResponse } from "next/server";
import { ZodError } from "zod";
import { BackendProblem } from "@/platform/backend/public-client";

export function problem(status: number, title: string, detail: string, code: string): NextResponse {
  return NextResponse.json({ type: "about:blank", title, status, detail, code }, {
    status,
    headers: { "cache-control": "no-store", "content-type": "application/problem+json" }
  });
}

export function errorResponse(error: unknown): NextResponse {
  if (error instanceof ZodError) return problem(400, "Solicitud inválida", "Revise los datos enviados.", "VALIDATION_FAILED");
  if (error instanceof BackendProblem) {
    const status = error.status >= 400 && error.status < 500 ? error.status : 502;
    return problem(status, "Backend no disponible", "No pudimos completar la operación.", error.code);
  }
  return problem(500, "Error interno", "La operación no pudo completarse.", "INTERNAL_ERROR");
}
