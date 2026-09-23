import { notFound, redirect } from "next/navigation";
import type { Route } from "next";
import { allowed, readSession } from "@/platform/auth/session";
import { loadPrivateI18n } from "@/platform/i18n/load-private-i18n";
import { organizationLabel } from "@/platform/experience/options";
import { documentId } from "@/platform/documents-v403/contract";
import { documentConfig } from "@/platform/documents-v403/server";
import { DocumentWorkspaceV403 } from "@/components/document-workspace-v403";

// AUTHORED opt-in reference composition. The same BFF/owner gates every effect.
export default async function DocumentPage({ searchParams }: {
  searchParams: Promise<{ document?: string | string[] }>;
}) {
  if (process.env.ELITE_EXPERIENCE_CATALOGUE !== "1") notFound();
  const session = await readSession();
  if (!session) redirect("/api/auth/login?return_to=/experience/documents" as Route);
  const { locale } = await loadPrivateI18n();
  const en = locale.language === "en";
  let config;
  try { config = documentConfig(); } catch {
    return <><h1>{en ? "Invoices" : "Facturas"}</h1><p>{en ? "This workspace is not enabled yet." : "Este espacio aún no está habilitado."}</p></>;
  }
  if (session.tenantId !== config.tenant || !session.organizations.includes(config.organization) || !["documents:write", "documents:process", "documents:review"].some(permission => allowed(session, permission))) {
    return <><h1>{en ? "Invoices" : "Facturas"}</h1><p>{en ? "Your account does not have access to this workspace." : "Tu cuenta no tiene acceso a este espacio."}</p></>;
  }
  const requested = (await searchParams).document;
  if (requested !== undefined && (typeof requested !== "string" || !documentId.safeParse(requested).success)) {
    return <><h1>{en ? "Invoices" : "Facturas"}</h1><p>{en ? "This review link is invalid. Request a new link." : "El enlace de revisión no es válido. Pedí un nuevo enlace."}</p></>;
  }
  const branch = await organizationLabel(session, config.organization);
  return <div className="experienceDocument">
    <p className="branchContext">{branch ?? (en ? "Branch: " : "Sucursal: ") + config.organization}</p>
    <h1>{en ? "Invoices" : "Facturas"}</h1>
    <DocumentWorkspaceV403 language={locale.language} {...(requested ? { initialDocumentId: requested } : {})}/>
  </div>;
}
