# TypeScript Franchise Journey Portals

V312: corrección AUTHORED de recepción y decisión de devoluciones.
Consulta exacta por autorización, scope de grafo bloqueado en escritura,
lectura atómica, actor y evidencia inmutables. UI conserva huella y consulta
sin reenviar aun con lista inaccesible. Refund4/exchange3 solicitudes.
Evidencia reconstruction_evidence/RETURN_OPERATIONS_RECOVERY_V312.md.
Sin efectos downstream, cambio comercial ni promoción integral.

V311: corrección AUTHORED de recuperación de checklist completada por
identidad natural de entrega, actor y respuestas inmutables. Consulta scoped
preserva estado vigente accepted/rejected/presented; no reenvía ni reabre.
Evidencia reconstruction_evidence/CHECKLIST_COMPLETION_RECOVERY_V311.md.
Sin cambio de reglas comerciales, upstreams ni promoción integral.

V310: corrección AUTHORED de publicación de checklist con consulta exacta
por identidad natural org/id/versión, actor atómico y recuperación sin POST.
Evidencia reconstruction_evidence/CHECKLIST_PUBLICATION_RECOVERY_V310.md.
No altera inmutabilidad, reglas de entrega, upstreams ni promoción integral.

V309: corrección AUTHORED de publicación/consulta de turnos con identidad
retenida, Idempotency-Key y recuperación scoped sin reenvío. Pruebas y límites
en reconstruction_evidence/SLOT_CREATION_RECOVERY_V309.md. No cambio de reglas de agenda,
upstreams ni promoción integral; FAIL457 conserva los otros comandos pendientes.

V308: corrección AUTHORED de alta/consulta de recursos con identidad
retenida, Idempotency-Key y recuperación scoped sin reenvío. Pruebas y límites
en reconstruction_evidence/RESOURCE_CREATION_RECOVERY_V308.md. No cambio de reglas de agenda,
upstreams ni promoción integral; FAIL457 conserva los otros comandos pendientes.

V307: corrección AUTHORED de alta/consulta de disponibilidad con identidad
retenida, Idempotency-Key y recuperación scoped sin reenvío. Pruebas y límites
en reconstruction_evidence/AVAILABILITY_CREATION_RECOVERY_V307.md. No cambio de reglas de agenda,
upstreams ni promoción integral; FAIL457 conserva los otros comandos pendientes.

V306: corrección AUTHORED de emisión/consulta de cotización con identidad
retenida, Idempotency-Key y recuperación scoped sin reenvío. Pruebas y límites
en reconstruction_evidence/QUOTE_CREATION_RECOVERY_V306.md. No cambio de pricing,
upstreams ni promoción integral; FAIL457 conserva los otros comandos pendientes.

V303 / 0.14.1 corrige resolución de discrepancias sin cambiar su dominio:
fence síncrono, controles tras hidratación, marker sessionStorage bajo hash de
tenant/subject/organización y excepción (sólo versión/decisión, sin notas),
timeout10s y recibo contrastado antes de anunciar registrado. GET vuelve a leer
el estado; ausencia/fallo no borra el marker ni permite repetir. Resolver con
otra decisión se comunica sin atribuirla al envío anterior. Ayuda versionada
delivery-resolution-view/1.0.0. No almacena secretos ni confirma efectos downstream.
AUTHORED; AWS idempotencia y React refs son método, no código copiado ni aval.
Evidencia: reconstruction_evidence/DELIVERY_EXCEPTION_RECOVERY_V303.md.
FAIL457 queda parcialmente corregido: los otros comandos del sender genérico
siguen pendientes de su propio circuito; no readmisión integral por este tramo.

V302 / 0.14.0 separa secciones/formularios por los permisos ya consumidos por Go.
FranchiseCommandPanel recibe permissions desde la sesión; omisión equivale a
cero permisos, no acceso general. Se permite entrar con resource:manage o
availability:read/manage sin exigir permisos comerciales ajenos. Consultas
independientes retornan null ante fallo y muestran alerta/GET/ayuda versionada,
sin llamar vacío a indisponible ni derribar pedidos por fallos de discrepancias.
No es nueva autorización: BFF/HTTP/tenant/organización siguen siendo obligatorios.
Evidencia: reconstruction_evidence/OPERATOR_SECTIONS_RECOVERY_V302.md.
Límite pendiente FAIL457: el sender genérico del operador aún afirma «No se
aplicó» ante resultado incierto; no usarlo como garantía de mutaciones hasta
cerrar su propio circuito de consulta/fence. V300 sólo corrigió el sender cliente.

V300 / 0.13.1 añade recuperación de lectura en Mis entregas: error comprensible,
cero acciones ante lectura fallida, reconsulta GET y ayuda/práctica/soporte
handover-read-view/1.0.0. Es AUTHORED, no código atribuido a Microsoft.
Recepción/discrepancia: fence síncrono, timeout10s, recibo contrastado y resultado
incierto sin reenvío habilitado; recuperación mediante GET, no rechazo supuesto.
Evidencia: reconstruction_evidence/DELIVERY_READ_RECOVERY_V300.md.

V296 / 0.12.0 conecta la reserva operativa a Commerce HTTP existente: selección
de serie desde snapshot scoped, autorización BFF y Go, versiones internas,
fence de reenvío y recuperación GET ante respuesta perdida. Muestra pago/entrega
guardados sin inventar cobros ni preparación. Ayuda order-operations-view/1.0.0.
Todo el delta es AUTHORED; no se atribuye a las empresas de sus frameworks.
Evidencia: reconstruction_evidence/ORDER_STOCK_CONNECTED_V296.md.

V294 / 0.11.1 corrige aceptación de cotizaciones: error de transporte implica
resultado incierto, no rechazo comercial. Timeout 10s, reenvío bloqueado hasta
lectura, referencia del pedido aceptado y fallback seguro si falla la consulta.
Ayuda única quote-acceptance-view/1.0.0 compartida entre resultado y error,
incluye práctica en fixture y escalamiento autorizado; no es LMS completo.
No cambia precio, política comercial ni persistencia; todo este delta es AUTHORED.
Prueba conectada y condiciones: reconstruction_evidence/QUOTE_ACCEPTANCE_CONNECTED_V294.md.

## 1. Metadata

V263 connects an operator agenda to GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.10.0. Appointment/resource selections come from the scoped server snapshot, versions remain internal, receipts are validated and an ambiguous response blocks further writes until reload. Confirmation, assignment and other existing appointment transitions have one agenda UI owner; their previous manual-ID forms are removed. The command route compares Origin exactly against the existing configured applicationBaseUrl, fails closed without it and does not trust forwarded headers. UTC-day and 200-row limits are explicit; role-specific simplification is proven only for the appointment-only operator. Local implementation is AUTHORED, not Microsoft/Meta source. Evidence: `reconstruction_evidence/OPERATOR_AGENDA_BROWSER_POSTGRES_V263.md`.

```yaml
pack_id: "TS-FRANCHISE-JOURNEY-PORTALS"
pack_version: "0.21.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Añade al BFF canónico ubicación/agenda/operación comercial y entrega con checklist, excepción, recepción física y disposición; sólo solicita efectos a los owners Go y nunca afirma stock, refund, cambio, accounting o fiscal completados."
stacks: ["Node.js 24.20.0", "TypeScript 7.0.2", "Next.js 16.3.4", "React 19.2.8", "Zod 4.4.3"]
compatible_with: ["TS-GO-API-WEB-BRIDGE 0.5.15", "TS-OIDC-PORTAL-ADAPTER 0.2.4", "GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.10.19", "GO-COMMERCE-PRICING-PAYMENT-API 0.6.1 for operation snapshot and audited allocation", "PYTHON-META-WHATSAPP-CLOUD-ADAPTER 0.13.0 for opt-in notification history"]
incompatible_with: ["browser-direct backend tokens", "alternate frontend owner for the same routes"]
license_expression: "LicenseRef-Workspace-Owner AND MIT dependencies"
upstream_sources: ["https://github.com/microsoft/BCApps/tree/31a860b527f0dc72c7a44a255d7e7d403cfa4789", "https://learn.microsoft.com/en-us/dynamics365/business-central/hr-how-manage-absence", "https://learn.microsoft.com/en-us/dynamics365/business-central/production-how-to-create-work-center-calendars", "https://learn.microsoft.com/en-us/dynamics365/business-central/ui-post-sales", "https://learn.microsoft.com/en-us/dynamics365/field-service/inspections-overview", "https://github.com/GoogleCloudPlatform/microservices-demo/tree/5b3a712ab85ccb8f6f7cd5b720d36ba9a8d041eb", "https://nextjs.org/docs", "https://sre.google/sre-book/reliable-product-launches/"]
verified_at: "2026-09-08"
```

Los treinta y un bloques son `AUTHORED`. Reutilizan el cliente público y el cliente OIDC server-only de los packs propietarios; no crean otro fetch layer, sesión, dominio, esquema ni backend. V278 añade consulta explícita read-only de WhatsApp en la agenda, scope estricto, recuperación sin POST y ayuda versionada. No acredita mensajes live ni capacitación real de personas.

## 2. Applicability

V298 conecta solicitud local de pago en la vista de pedidos ya existente.
Permiso payment:create y selección leída del backend; BFF strict/same-origin no
acepta importe/proveedor/actor. Clave UUID conservada en sessionStorage bajo
scope hash de tenant/subject/organización y pedido; si storage falla no envía.
Respuesta incierta exige GET y conserva clave; request existente no se duplica.
Ayuda/práctica/soporte order-operations-view/1.1.0. No confirma dinero ni entrega.
No dependencia nueva; AUTHORED, con React/event-handlers y AWS/idempotencia como
fundamentos, no atribución del código. Evidencia PAYMENT_REQUEST_PORTAL_V298.md.

V261 corrige el formulario de turnos: controles deshabilitados hasta hidratación, receipt ID/state validado y referencia visible antes de anunciar éxito, clave conservada ante fallo y botón bloqueado tras solicitud aceptada. No afirma confirmación operativa ni recuperación tras reiniciar el navegador. Verificación conectada con Microsoft Playwright y Go/PostgreSQL: `reconstruction_evidence/PUBLIC_APPOINTMENT_RECOVERY_V261.md`; cambios `AUTHORED` bajo método público documentado, no código atribuido a Microsoft.

Use only with the canonical web BFF and Go journey API. Reject when the selected project has no public locations/appointments or when another frontend owns those exact routes. The portal assumes server-side OIDC and binds franchise data to the first admitted organization in the session.

## 3. Architecture contract

The browser posts only to the same-origin BFF. Public commands require an exact Origin match before session access, Zod rejects unknown fields, organization membership and action permission are checked again, and the shared OIDC client keeps bearer tokens server-only. Idempotency is forwarded unchanged for public appointments and provider/backend errors are bounded. Public locations and future non-full slots are read from explicitly published backend rows; the browser cannot invent a datetime and sends the exact server slot kind/start. The franchise page publishes bounded capacity only with `appointment:manage`. The lead form in the bridge returns an opaque lead reference and links directly to the appointment journey; no PII is placed in the URL.

Availability reads run server-side only when the token has `availability:read`. Every mutation crosses one strict discriminated union, rejects unknown/server-owned fields, checks organization membership and maps to one existing Go owner. The UI provides explicit organization/resource intervals, resource skills, assignment and reason-bound transitions; it does not calculate recurrence, entitlement or labor policy. Customer cancellation uses the verified session and never accepts a customer identifier.

The franchise UI constructs a bounded checklist definition without a parallel backend owner, publishes an immutable explicit version, and completes a prepared handover by exact item IDs. The BFF admits only the closed response schema and `handover:manage`; server identity owns the operator and all durable evidence. The customer sees the published controls and completion timestamp, while acceptance remains unavailable unless checklist ID, checklist version and completion are present. Its command sends those exact values back to Go and still cannot supply its own evidence digest.

For a presented handover, the customer may submit either acceptance or a reason-coded discrepancy, never both after the server changes version/state. The BFF does not accept customer subject, exception ID, hash, stock state, refund or accounting fields. Franchise users read scoped exceptions server-side and choose one closed resolution action; Go creates the successor or return authorization. The UI describes authorization accurately and never claims that money, inventory, exact-cost reversal or fiscal posting already occurred.

Franchise users read scoped return cases server-side. Receipt accepts only the durable serial, a closed physical-condition value and bounded notes; the BFF hashes the canonical admitted command instead of accepting browser evidence. Disposition accepts only a closed inventory action. The browser cannot choose refund versus exchange, request IDs, owner contexts, state, payment, accounting or fiscal fields. It displays each downstream request as `requested`, never as completed.

The pack does not duplicate CRM, appointment, availability, price, order, identity or organization ownership. The command UI never accepts price/currency: the Go service resolves both from the active price book. Project edge must still reprove origin/host behavior, rate limiting, bot defense, accessibility with users, real IdP roles, backend failure UX and browser/edge performance.

## 4. Exact file manifest

```text
CREATE src/platform/notifications/status-contract.ts
CREATE src/platform/notifications/status-contract.test.ts
CREATE src/app/api/enterprise/franchise/notifications/route.ts
CREATE src/app/api/enterprise/franchise/notifications/route.test.ts
CREATE src/components/appointment-notification-status.tsx
CREATE docs/whatsapp-status-operations.md
CREATE src/app/api/enterprise/appointments/route.ts
CREATE src/app/api/enterprise/appointments/route.test.ts
CREATE src/app/api/enterprise/franchise/commands/route.ts
CREATE src/app/api/enterprise/franchise/commands/route.test.ts
CREATE src/app/locations/appointment-form.tsx
CREATE src/app/locations/page.tsx
CREATE src/app/franchise/page.tsx
CREATE src/components/franchise-command-panel.tsx
CREATE src/components/customer-quote-actions.tsx
CREATE src/app/customer/quotes/page.tsx
CREATE src/components/customer-appointment-actions.tsx
CREATE src/app/customer/appointments/page.tsx
CREATE src/components/customer-handover-actions.tsx
CREATE src/app/customer/handovers/page.tsx
CREATE src/platform/help/content.ts
CREATE src/platform/help/contract.ts
CREATE src/platform/help/catalog.ts
CREATE src/app/api/enterprise/help/route.ts
CREATE src/app/api/enterprise/help/route.test.ts
CREATE src/components/journey-help-panel.tsx
CREATE src/app/help/page.tsx
CREATE docs/journey-help.md
CREATE src/components/operational-guide.tsx
CREATE src/platform/help/coverage.test.ts
CREATE docs/journey-help-coverage.md
CREATE docs/handover-browser-gate.md
CREATE docs/handover-operator-flow.md
CREATE microsoft_playwright_browser_gate/tests/handover-connected.spec.mjs
CREATE src/app/api/enterprise/handovers/bounded-body.test.ts
CREATE src/app/api/enterprise/handovers/route.test.ts
CREATE src/app/api/enterprise/handovers/route.ts
CREATE src/components/handover-operations-panel.tsx
CREATE src/platform/backend/handover-recovery.test.ts
CREATE src/platform/handovers/contracts.test.ts
CREATE src/platform/handovers/contracts.ts
CREATE src/platform/handovers/customer-dates.test.ts
CREATE src/platform/handovers/render.test.ts
CREATE src/app/franchise/whatsapp/page.tsx
CREATE src/components/whatsapp-reply-review.tsx
CREATE src/platform/backend/bounded-command.ts
CREATE src/app/api/enterprise/franchise/commands/body-bound.test.ts
CREATE src/app/api/enterprise/locale/route.test.ts
CREATE src/app/api/enterprise/locale/route.ts
CREATE src/platform/help/query.ts
CREATE src/platform/i18n/private-guide-display.ts
```

## 5. Materialization blocks

### FILE: `src/platform/notifications/status-contract.ts`

```yaml
block_id: "TS-FRANCHISE-JOURNEY-PORTALS:v278:1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "501130c6f091a1e007d4c6af585cb8acb0ea2fab95f2eadc786329db41ca3fd7"
variables: []
secrets_allowed: false
```

````typescript
import { z } from "zod";

// AUTHORED narrow projection of the existing Go reader, not a Meta SDK.
export const STATUS_VIEW_VERSION = "whatsapp-status-view/1.0.0";
const identifier = z.string().min(1).max(128);
const instant = z.iso.datetime({ offset: true });
const statusSchema = z.object({
  delivery_key: z.string().min(1).max(256),
  fence_state: z.enum(["not_started", "sending", "unknown", "accepted", "failed_terminal"]),
  delivery_status: z.enum(["not_observed_by_this_reader", "ambiguous_latest_timestamp", "observed_sent", "observed_delivered", "observed_read", "observed_failed", "observed_deleted"]),
  provider_timestamp: z.number().int().positive().max(999_999_999_999).optional(),
  provider_event_count: z.number().int().nonnegative().max(Number.MAX_SAFE_INTEGER),
  approval_expires_at: instant, approval_expired: z.boolean(), accepted_at: instant.optional(), updated_at: instant.optional(),
  observed_at: instant, reconciliation_required: z.boolean(),
}).strict().superRefine((value, context) => {
  if (value.provider_event_count === 0 ? value.provider_timestamp !== undefined || value.delivery_status !== "not_observed_by_this_reader" : value.provider_timestamp === undefined || value.delivery_status === "not_observed_by_this_reader" || value.fence_state !== "accepted") context.addIssue({ code: "custom", message: "inconsistent observations" });
  if (value.fence_state === "accepted" && !value.accepted_at) context.addIssue({ code: "custom", message: "missing acceptance timestamp" });
  if (value.fence_state === "unknown" && !value.reconciliation_required) context.addIssue({ code: "custom", message: "unknown outcome requires reconciliation" });
});
export const notificationHistorySchema = z.object({
  organization_id: identifier, appointment_id: identifier,
  items: z.array(z.object({ confirmation_event_id: z.uuid(), status: statusSchema }).strict()).max(20),
}).strict().superRefine((value, context) => {
  if (new Set(value.items.map(item => item.confirmation_event_id)).size !== value.items.length) context.addIssue({ code: "custom", message: "duplicate confirmation" });
});
export type NotificationHistory = z.infer<typeof notificationHistorySchema>;
export type NotificationStatus = NotificationHistory["items"][number]["status"];

export function notificationExplanation(status: NotificationStatus): { title: string; guidance: string } {
  if (status.reconciliation_required) return { title: "Resultado por conciliar", guidance: "No reenvíes el mensaje. Soporte debe revisar el intento y su evidencia antes de autorizar otra acción." };
  if (status.delivery_status === "ambiguous_latest_timestamp") return { title: "Estados contradictorios", guidance: "Hay avisos distintos con la misma fecha del proveedor. Conservá la referencia y pedí revisión; no elijas uno por orden de llegada." };
  const observed: Record<string, string> = { observed_sent: "El proveedor informó envío", observed_delivered: "El proveedor informó entrega", observed_read: "El proveedor informó lectura", observed_failed: "El proveedor informó un fallo", observed_deleted: "El proveedor informó eliminación" };
  if (observed[status.delivery_status]) return { title: observed[status.delivery_status]!, guidance: "Este aviso no acredita una venta ni la aceptación del turno por el cliente. No vuelve a enviar mensajes." };
  switch (status.fence_state) {
    case "accepted": return { title: "Aceptado por el proveedor; entrega sin verificar", guidance: "El envío fue aceptado, pero todavía no hay un aviso verificado de entrega. Actualizá la consulta; no reenvíes por falta de confirmación." };
    case "sending": return { title: "Envío en proceso", guidance: "Esperá y actualizá la consulta. No ejecutes un segundo envío mientras el intento siga abierto." };
    case "failed_terminal": return { title: "Intento cerrado con fallo", guidance: "Revisá la evidencia con soporte. Esta pantalla no habilita reintentos ni cambia el historial." };
    default: return { title: "Sin intento de envío registrado", guidance: "Una aprobación no implica un envío. Esta consulta es sólo de lectura; no inicia mensajes ni captura consentimiento." };
  }
}
````

### FILE: `src/platform/notifications/status-contract.test.ts`

```yaml
block_id: "TS-FRANCHISE-JOURNEY-PORTALS:v278:2"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "1c6b7eacbb277fe3c66e635852c9ffbd6c0b65a98e8892e0028d145339464f62"
variables: []
secrets_allowed: false
```

````typescript
import { describe, expect, it } from "vitest";
import { readFileSync } from "node:fs";
import { notificationExplanation, notificationHistorySchema, STATUS_VIEW_VERSION, type NotificationStatus } from "./status-contract";
const status: NotificationStatus = { delivery_key: "test-delivery-reference", fence_state: "accepted", delivery_status: "not_observed_by_this_reader", provider_event_count: 0, approval_expires_at: "2026-09-07T01:00:00Z", approval_expired: true, accepted_at: "2026-09-06T01:00:00Z", observed_at: "2026-09-07T02:00:00Z", reconciliation_required: false };
const item = { confirmation_event_id: "11111111-1111-4111-8111-111111111111", status };
const history = { organization_id: "store", appointment_id: "appointment", items: [item] };
describe("notification projection and operational explanation", () => {
  it("ships the same operational guide version as the UI", () => { const guide=readFileSync(new URL("../../../docs/whatsapp-status-operations.md", import.meta.url),"utf8"); expect(guide).toContain(`Guía: ${STATUS_VIEW_VERSION}`); expect(guide).toContain("no marca automáticamente personas capacitadas"); });
  it("does not equate accepted with delivery or business acceptance", () => { expect(notificationHistorySchema.safeParse(history).success).toBe(true); expect(notificationExplanation(status).title).toContain("entrega sin verificar"); expect(STATUS_VIEW_VERSION).toBe("whatsapp-status-view/1.0.0"); });
  it.each(["observed_sent", "observed_delivered", "observed_read", "observed_failed", "observed_deleted"] as const)("labels verified %s as a provider report", delivery_status => { const value = { ...status, delivery_status, provider_event_count: 1, provider_timestamp: 1603086314 }; expect(notificationHistorySchema.safeParse({ ...history, items: [{ ...item, status: value }] }).success).toBe(true); expect(notificationExplanation(value).title).toContain("proveedor informó"); expect(notificationExplanation(value).guidance).toContain("no acredita una venta"); });
  it("keeps ambiguous and unknown states actionable without resend", () => { expect(notificationExplanation({ ...status, delivery_status: "ambiguous_latest_timestamp" }).title).toContain("contradictorios"); expect(notificationExplanation({ ...status, fence_state: "unknown", reconciliation_required: true }).guidance).toContain("No reenvíes"); });
  it.each([{ ...status, provider_timestamp: 1 }, { ...status, provider_event_count: 1 }, { ...status, fence_state: "unknown" }, { ...status, accepted_at: undefined }, { ...status, delivery_status: "unrecognized" }, { ...status, raw_message: "secret" }])("rejects contradictory/unknown fields", invalid => { expect(notificationHistorySchema.safeParse({ ...history, items: [{ ...item, status: invalid }] }).success).toBe(false); });
  it("rejects duplicate or oversized history", () => { expect(notificationHistorySchema.safeParse({ ...history, items: [item, item] }).success).toBe(false); expect(notificationHistorySchema.safeParse({ ...history, items: Array(21).fill(item) }).success).toBe(false); });
});
````

### FILE: `src/app/api/enterprise/franchise/notifications/route.ts`

```yaml
block_id: "TS-FRANCHISE-JOURNEY-PORTALS:v278:3"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "ebfe5ab0ca8a45c3601666e2248e894bf1acd7bbf2a7a61e0c8cf813be438455"
variables: []
secrets_allowed: false
```

````typescript
import { NextResponse } from "next/server";
import { allowed, readSession } from "@/platform/auth/session";
import { protectedGet } from "@/platform/backend/protected-client";
import { BackendProblem } from "@/platform/backend/public-client";
import { notificationHistorySchema } from "@/platform/notifications/status-contract";
import { loadBusinessConfig } from "@/platform/config/load";

const reply = (value: unknown, status = 200) => NextResponse.json(value, { status, headers: { "Cache-Control": "no-store", "X-Content-Type-Options": "nosniff" } });
const safeID = (value: string | null): value is string => value !== null && /^[a-zA-Z0-9][a-zA-Z0-9_-]{0,127}$/.test(value);

// Read-only same-site BFF. Token/tenant are resolved server-side, never returned.
export async function GET(request: Request) {
  try {
    if ((await loadBusinessConfig()).features.whatsapp_status_history !== true) return reply({ code: "CAPABILITY_NOT_ENABLED" }, 404);
  } catch { return reply({ code: "NOTIFICATION_READ_UNAVAILABLE" }, 503); }
  if (["cross-site", "none"].includes(request.headers.get("sec-fetch-site") ?? "")) return reply({ code: "CROSS_SITE_REJECTED" }, 403);
  const session = await readSession();
  if (!session) return reply({ code: "UNAUTHENTICATED" }, 401);
  if (!allowed(session, "appointment:manage")) return reply({ code: "FORBIDDEN" }, 403);
  const query = new URL(request.url).searchParams;
  const organization = query.get("organizationId"), appointment = query.get("appointmentId");
  if (request.url.length > 2048 || [...query.keys()].length !== 2 || query.getAll("organizationId").length !== 1 || query.getAll("appointmentId").length !== 1 || !safeID(organization) || !safeID(appointment)) return reply({ code: "INVALID_QUERY" }, 400);
  if (!session.organizations.includes(organization) && !session.organizations.includes("*")) return reply({ code: "ORGANIZATION_FORBIDDEN" }, 403);
  try {
    const raw = await protectedGet<unknown>(session, `/v1/franchise/appointments/${encodeURIComponent(appointment)}/whatsapp-confirmations`, { organization_id: organization });
    const parsed = notificationHistorySchema.safeParse(raw);
    if (!parsed.success || parsed.data.organization_id !== organization || parsed.data.appointment_id !== appointment) return reply({ code: "INVALID_NOTIFICATION_RESPONSE" }, 502);
    return reply(parsed.data);
  } catch (error) {
    const permitted = new Set(["NOTIFICATION_NOT_FOUND", "NOTIFICATION_EVIDENCE_MISMATCH", "NOTIFICATION_HISTORY_LIMIT", "NOTIFICATION_READ_UNAVAILABLE", "UNAUTHENTICATED", "FORBIDDEN", "ORGANIZATION_FORBIDDEN"]);
    if (error instanceof BackendProblem && permitted.has(error.code) && [401, 403, 404, 409, 503].includes(error.status)) return reply({ code: error.code }, error.status);
    return reply({ code: "NOTIFICATION_READ_UNAVAILABLE" }, 503);
  }
}
````

### FILE: `src/app/api/enterprise/franchise/notifications/route.test.ts`

```yaml
block_id: "TS-FRANCHISE-JOURNEY-PORTALS:v278:4"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "b1cac0a7bc3605e88dab51b33d0f044495de55cd5de52079cdba1b7afdf25238"
variables: []
secrets_allowed: false
```

````typescript
import { beforeEach, describe, expect, it, vi } from "vitest";
import { GET } from "./route";
import { BackendProblem } from "@/platform/backend/public-client";
const mocks = vi.hoisted(() => ({ readSession: vi.fn(), protectedGet: vi.fn(), loadBusinessConfig: vi.fn() }));
vi.mock("@/platform/config/load", () => ({ loadBusinessConfig: mocks.loadBusinessConfig }));
vi.mock("@/platform/auth/session", () => ({ readSession: mocks.readSession, allowed: (s: { permissions: string[] }, p: string) => s.permissions.includes(p) }));
vi.mock("@/platform/backend/protected-client", () => ({ protectedGet: mocks.protectedGet }));
const session = { subject: "operator", tenantId: "tenant", permissions: ["appointment:manage"], organizations: ["store"], accessToken: "server-secret" };
const empty = { organization_id: "store", appointment_id: "appointment", items: [] };
const request = (query = "organizationId=store&appointmentId=appointment", site = "same-origin") => new Request(`https://portal.example/api/enterprise/franchise/notifications?${query}`, { headers: { "sec-fetch-site": site } });
beforeEach(() => { vi.clearAllMocks(); mocks.readSession.mockResolvedValue(session); mocks.protectedGet.mockResolvedValue(empty); mocks.loadBusinessConfig.mockResolvedValue({ features: { whatsapp_status_history: true } }); });
describe("notification history BFF", () => {
  it.each([{}, { whatsapp_status_history: false }])("does not activate an unconfigured capability", async features => {
    mocks.loadBusinessConfig.mockResolvedValue({ features });
    const response = await GET(request()); expect(response.status).toBe(404);
    expect(await response.json()).toEqual({ code: "CAPABILITY_NOT_ENABLED" });
    expect(mocks.readSession).not.toHaveBeenCalled(); expect(mocks.protectedGet).not.toHaveBeenCalled();
  });
  it("fails closed when configuration is unavailable", async () => {
    mocks.loadBusinessConfig.mockRejectedValue(new Error("private config path"));
    const response = await GET(request()); expect(response.status).toBe(503);
    expect(await response.text()).not.toContain("private"); expect(mocks.protectedGet).not.toHaveBeenCalled();
  });
  it("forwards only authenticated scope and returns no token", async () => {
    const response = await GET(request()); expect(response.status).toBe(200); expect(await response.json()).toEqual(empty); expect(response.headers.get("cache-control")).toBe("no-store");
    expect(mocks.protectedGet).toHaveBeenCalledWith(session, "/v1/franchise/appointments/appointment/whatsapp-confirmations", { organization_id: "store" });
  });
  it("rejects missing session, permissions and foreign organization", async () => {
    mocks.readSession.mockResolvedValue(null); expect((await GET(request())).status).toBe(401);
    mocks.readSession.mockResolvedValue({ ...session, permissions: [] }); expect((await GET(request())).status).toBe(403);
    mocks.readSession.mockResolvedValue(session); expect((await GET(request("organizationId=other&appointmentId=appointment"))).status).toBe(403);
    expect(mocks.protectedGet).not.toHaveBeenCalled();
  });
  it("rejects cross-site before resolving credentials", async () => {
    expect((await GET(request(undefined, "cross-site"))).status).toBe(403); expect(mocks.readSession).not.toHaveBeenCalled();
  });
  it.each(["organizationId=store&appointmentId=a&tenantId=other", "organizationId=store&appointmentId=a&appointmentId=b", "organizationId=store", "organizationId=store&appointmentId=..", "organizationId=store&appointmentId=a%2Fb", "organizationId=store&appointmentId=" + "a".repeat(129)])("rejects ambiguous/unbounded query %s", async query => { expect((await GET(request(query))).status).toBe(400); expect(mocks.protectedGet).not.toHaveBeenCalled(); });
  it.each([{ ...empty, organization_id: "other" }, { ...empty, appointment_id: "other" }, { ...empty, recipient: "secret" }, { ...empty, items: [{}] }])("rejects foreign or malformed backend projection", async raw => { mocks.protectedGet.mockResolvedValue(raw); const response = await GET(request()); expect(response.status).toBe(502); expect(await response.text()).not.toContain("secret"); });
  it("preserves approved backend errors without arbitrary details", async () => {
    mocks.protectedGet.mockRejectedValue(new BackendProblem(409, "NOTIFICATION_HISTORY_LIMIT")); const response = await GET(request()); expect(response.status).toBe(409); expect(await response.json()).toEqual({ code: "NOTIFICATION_HISTORY_LIMIT" });
    mocks.protectedGet.mockRejectedValue(new Error("private payload")); expect(await (await GET(request())).text()).not.toContain("private");
  });
});
````

### FILE: `src/components/appointment-notification-status.tsx`

```yaml
block_id: "TS-FRANCHISE-JOURNEY-PORTALS:v278:5"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "989f9eb5efaf787c1aa08309c2127fb75868feebf2b45ae38a68deab0693203a"
variables: []
secrets_allowed: false
```

````typescript
"use client";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

import { NOTIFICATION_GUIDE } from "@/platform/help/content";

import { useEffect, useRef, useState } from "react";
import { readBackendResponse } from "@/platform/backend/public-client";
import { notificationExplanation, notificationHistorySchema, STATUS_VIEW_VERSION, type NotificationHistory } from "@/platform/notifications/status-contract";



// No automatic polling/N+1 page load or write action. Abort stale reads when
// scope changes/unmounts, following React's effect cleanup contract.
export function AppointmentNotificationStatus({ organization, appointment }: { organization: string; appointment: string }) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();
 const displayTime=(value:string|number)=>new Date(value).toLocaleString(privateLocale.locale,{timeZone:privateLocale.timeZone,dateStyle:"short",timeStyle:"medium"})+" ("+privateLocale.timeZone+")";

  const [ready, setReady] = useState(false), [busy, setBusy] = useState(false);
  const [history, setHistory] = useState<NotificationHistory | null>(null), [error, setError] = useState("");
  const request = useRef<AbortController | null>(null);
  useEffect(() => {
    setReady(true); setHistory(null); setError(""); setBusy(false);
    return () => { request.current?.abort(); request.current = null; };
  }, [organization, appointment]);
  async function refresh() {
    if (request.current) return;
    const controller = new AbortController(); request.current = controller;
    setBusy(true); setError(""); setHistory(null);
    try {
      const query = new URLSearchParams({ organizationId: organization, appointmentId: appointment });
      const response = await fetch(`/api/enterprise/franchise/notifications?${query}`, { cache: "no-store", redirect: "error", signal: AbortSignal.any([controller.signal, AbortSignal.timeout(7000)]) });
      const value = notificationHistorySchema.parse(await readBackendResponse<unknown>(response));
      if (value.organization_id !== organization || value.appointment_id !== appointment) throw new Error("scope mismatch");
      if (!controller.signal.aborted) setHistory(value);
    } catch {
      if (!controller.signal.aborted) setError(t("p0126"));
    } finally { if (request.current === controller) { request.current = null; setBusy(false); } }
  }
  const visible = history?.organization_id === organization && history.appointment_id === appointment ? history : null;
  return <section aria-label={t("p0127")} className="notificationStatus">
    <h4>{t("p0128")}</h4>
    <button type="button" className="button" disabled={!ready || busy} onClick={() => void refresh()}>{busy ? t("p0129") : t("p0130")}</button>
    <p role="status" aria-label={t("p0131")} aria-live="polite">{error || (busy ? t("p0132") : visible ? t("p0133") : t("p0134"))}</p>
    {visible?.items.length === 0 ? <p>{t("p0135")}</p> : null}
    {visible?.items.map(item => {
      const explanation = notificationExplanation(item.status);
      return <article key={item.confirmation_event_id}>
        <h5>{controlled(explanation.title)}</h5><p>{controlled(explanation.guidance)}</p>
        <p>{t("p0136")} <time dateTime={item.status.observed_at}>{displayTime(item.status.observed_at)}</time></p>
        {item.status.provider_timestamp ? <p>{t("p0137")} <time dateTime={new Date(item.status.provider_timestamp * 1000).toISOString()}>{displayTime(item.status.provider_timestamp * 1000)}</time></p> : null}
        {item.status.approval_expired ? <p>{t("p0138")}</p> : null}
        <details><summary>{t("p0139")}</summary>
          <p>{controlled(NOTIFICATION_GUIDE.paragraphs[0])}</p>
          <dl><dt>{t("p0140")}</dt><dd>{appointment}</dd><dt>{t("p0141")}</dt><dd>{item.confirmation_event_id}</dd><dt>{t("p0142")}</dt><dd>{STATUS_VIEW_VERSION}</dd></dl>
          <p>{controlled(NOTIFICATION_GUIDE.paragraphs[1])}</p>
        <a href={`/help?article=${NOTIFICATION_GUIDE.id}&version=${NOTIFICATION_GUIDE.version}`}>{t("p0143")}</a>
        </details>
      </article>;
    })}
  </section>;
}
````

### FILE: `docs/whatsapp-status-operations.md`

```yaml
block_id: "TS-FRANCHISE-JOURNEY-PORTALS:v278:6"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "5619c6d4d1791107592dccd99334d3758039d7e7200b5672fa2299036ece62a8"
variables: []
secrets_allowed: false
```

````markdown
# Ayuda operativa — estado de notificaciones

Guía: whatsapp-status-view/1.0.0

Procedencia: AUTHORED. Aplica sólo al panel de lectura de notificaciones de un
turno y al worker de statuses de WhatsApp de esta release. No es una metodología
publicada literalmente por Meta/Google ni capacitación de toda la empresa.

## Operador con appointment:manage en la organización actual

Activación por el responsable técnico: usar el owner de configuración existente,
features.whatsapp_status_history=true, únicamente después de montar y verificar
el módulo Go de notificaciones en ENTERPRISE_API_BASE_URL. Ausente o false oculta
el panel y el BFF devuelve CAPABILITY_NOT_ENABLED sin consultar credenciales ni
backend. El flag no concede permisos, consentimiento ni aprobación de envío.
La configuración es server-side y se aplica al reiniciar la release; no crear
un segundo registro de funcionalidades ni activar servicios pagos con este flag.

1. Abrí Franquicia, elegí la fecha UTC y el turno por su referencia visible.
2. Pulsá Consultar estado de WhatsApp. No hay envío ni sondeo automático.
3. Separá la fecha de consulta de la fecha del aviso del proveedor: una consulta
   reciente no vuelve reciente un aviso antiguo.
4. Aceptado por el proveedor no significa entregado. Entrega/lectura informadas
   tampoco significan venta o aceptación del turno por parte del cliente.
5. Si falla la lectura, consultá de nuevo. Si persiste, abrí la ayuda y comunicá
   la referencia únicamente a soporte autorizado. Nunca pegues tokens, teléfonos
   ni contenido del mensaje en tickets o capturas no autorizados.

## Recuperación y soporte

- Resultado incierto, intento en proceso o estado contradictorio: no reenvíes.
  Soporte revisa el intento, las observaciones y sus fechas, permisos/organización,
  originales retenidos y receipts anclados. Preserva evidencia y el ledger.
- Sin aprobaciones: no hay una aprobación visible para ese turno; no demuestra
  ausencia de todo mensaje externo ni autoriza crear consentimiento.
- Límite de historial: la vista no devuelve una lista parcial. Soporte debe
  resolver una lectura acotada/paginada admitida; no borrar historia para entrar
  en el límite ni ignorar el error.
- Un fallo no autoriza ampliar permisos. 401 requiere revisar sesión; 403,
  autorización vigente; 404, scope/referencia; 409, evidencia o límite; 503,
  disponibilidad. El panel puede mostrar un error seguro genérico.
- El worker conserva auditoría y terminales existentes. TERMINAL_REVIEW requiere
  revisión; RECONCILE_REQUIRED sin FailureRecorded no acredita registro durable.
  Una caída del reporter detiene Run y debe alertarse mediante el supervisor.
  No reiniciar generations/attempts ni marcar inbox procesado manualmente.

## Práctica de esta versión, sin datos reales

En el fixture aislado de navegador: consultar por teclado, perder una respuesta
GET y recuperarla, revisar entrega informada, desplegar ayuda, rechazar respuesta
de otra organización y verificar cero POST. El gate SQL conserva un intento de
envío y una observación. El ejercicio no demuestra una cuenta Meta real.

El operador debe poder explicar: (a) diferencia entre aceptación y entrega;
(b) por qué no se reenvía un resultado incierto; (c) qué referencia compartir y
qué datos nunca adjuntar; (d) diferencia entre fecha de consulta y del aviso.
El responsable del proyecto registra evaluación real por rol/versión; esta guía
no marca automáticamente personas capacitadas ni autorizaciones cumplidas.

## Actualización y rollback

Si cambia contrato, significado de estado, permisos, recovery o interacción,
actualizar juntos status-contract.ts, componente, BFF, Go, esta guía y tests.
Cambiar la versión de guía y reentrenar los roles afectados antes de promover.
Pausar el host ante inconsistencia; conservar jobs, approvals, audit, observaciones
y fences. Reponer artefacto compatible verificado, no reactivar SQL inseguro.

Autoridades de método consultadas 2026-09-06: [Google SRE, gestión de incidentes](https://sre.google/sre-book/managing-incidents/)
(coordinación/registro), [Google SRE, monitoreo](https://sre.google/sre-book/monitoring-distributed-systems/)
(señales accionables) y [Microsoft Playwright, buenas prácticas](https://playwright.dev/docs/best-practices)
(comportamiento visible y aislamiento). La traducción a este flujo es local.
````


### FILE: `src/app/api/enterprise/franchise/commands/route.ts`

```yaml
block_id: "TS-FRANCHISE-JOURNEY:file:06"
operation: CREATE
provenance: AUTHORED
source: "local same-origin command adapter governed by pack references"
license: "LicenseRef-Workspace-Owner"
sha256: "2a9e1ffb99754e25a86aecf13e855cf9078e29366c5a7886e630688720e55259"
variables: []
secrets_allowed: false
```

````ts
import { NextResponse, type NextRequest } from "next/server";
import { z } from "zod";
import { allowed, readSession } from "@/platform/auth/session";
import { applicationBaseUrl } from "@/platform/auth/oidc-client";
import { boundedCommandBody } from "@/platform/backend/bounded-command";
import { BackendProblem } from "@/platform/backend/public-client";
import { protectedGet, protectedPost } from "@/platform/backend/protected-client";

const identifier = z.string().min(1).max(128);
const organization = z.string().min(1).max(128);
const version = z.number().int().positive();
const reasonCode = z.string().regex(/^[a-z][a-z0-9]*(-[a-z0-9]+)*$/);
const checklistItem = z.object({ id: identifier, prompt: z.string().trim().min(1).max(500), response_type: z.enum(["confirmation", "text", "serial", "evidence"]), required: z.boolean() }).strict();
const checklistResponse = z.object({ item_id: identifier, response_text: z.string().trim().min(1).max(2048), evidence_sha256: z.string().regex(/^[0-9a-f]{64}$/).optional() }).strict();

const commandSchema = z.discriminatedUnion("action", [
  z.object({action:z.literal("request-order-payment"),organizationId:organization,orderId:identifier,requestKey:z.uuid()}).strict(),
  z.object({action:z.literal("allocate-order-stock"),organizationId:organization,orderId:identifier,lineId:identifier,stockUnitId:identifier,orderVersion:version,stockVersion:version}).strict(),
  z.object({ action: z.literal("assign-lead"), organizationId: organization, leadId: identifier, assignedSubject: z.string().min(1).max(256), version }).strict(),
  z.object({ action: z.literal("transition-lead"), organizationId: organization, leadId: identifier, current: z.enum(["new", "contacted", "qualified"]), target: z.enum(["contacted", "qualified", "converted", "lost"]), version }).strict(),
  z.object({ action: z.literal("create-quote"), organizationId: organization, leadId: identifier, variantId: identifier, priceBookId: identifier, validUntil: z.iso.datetime({ offset: true }) }).strict(),
  z.object({ action: z.literal("create-appointment-slot"), organizationId: organization, kind: z.enum(["consultation", "test-drive", "delivery", "service"]), startsAt: z.iso.datetime({ offset: true }), endsAt: z.iso.datetime({ offset: true }), capacity: z.number().int().min(1).max(100) }).strict(),
  z.object({ action: z.literal("create-resource"), organizationId: organization, principalSubject: z.string().max(256).optional(), displayName: z.string().min(1).max(256), kind: z.enum(["employee", "contractor", "service-bay", "vehicle"]), skills: z.array(z.enum(["consultation", "test-drive", "delivery", "service"])).min(1).max(16) }).strict(),
  z.object({ action: z.literal("create-availability"), organizationId: organization, resourceId: identifier.optional(), entryType: z.enum(["working", "unavailable"]), reasonCode: reasonCode.optional(), startsAt: z.iso.datetime({ offset: true }), endsAt: z.iso.datetime({ offset: true }) }).strict(),
  z.object({ action: z.literal("cancel-availability"), organizationId: organization, availabilityId: identifier, version, reasonCode }).strict(),
  z.object({ action: z.literal("assign-appointment-resource"), organizationId: organization, appointmentId: identifier, resourceId: identifier, version }).strict(),
  z.object({ action: z.literal("transition-appointment"), organizationId: organization, appointmentId: identifier, current: z.enum(["requested", "confirmed"]), target: z.enum(["confirmed", "completed", "cancelled", "no-show"]), version, reasonCode: reasonCode.optional() }).strict(),
  z.object({ action: z.literal("cancel-customer-appointment"), organizationId: organization, appointmentId: identifier, version, reasonCode }).strict(),
  z.object({ action: z.literal("accept-quote"), organizationId: organization, quoteId: identifier, version }).strict(),
  z.object({ action: z.literal("publish-delivery-checklist"), organizationId: organization, checklistId: identifier, version, title: z.string().trim().min(1).max(160), items: z.array(checklistItem).min(1).max(64) }).strict(),
  z.object({ action: z.literal("complete-delivery-checklist"), organizationId: organization, handoverId: identifier, version, checklistId: identifier, checklistVersion: version, responses: z.array(checklistResponse).min(1).max(64) }).strict(),
  z.object({ action: z.literal("accept-handover"), organizationId: organization, handoverId: identifier, version, confirmedReceived: z.literal(true), serialNumber: z.string().min(1).max(128), checklistId: identifier, checklistVersion: version }).strict(),
  z.object({ action: z.literal("reject-handover"), organizationId: organization, handoverId: identifier, version, reasonCode, details: z.string().trim().min(1).max(1000) }).strict(),
  z.object({ action: z.literal("resolve-delivery-exception"), organizationId: organization, exceptionId: identifier, version, resolutionAction: z.enum(["correct-and-represent", "return", "exchange"]), notes: z.string().trim().min(1).max(1000) }).strict(),
  z.object({ action: z.literal("receive-return"), organizationId: organization, authorizationId: identifier, serialNumber: z.string().min(1).max(128), conditionCode: z.enum(["sealed", "opened", "damaged", "incomplete"]), notes: z.string().trim().min(1).max(1000) }).strict(),
  z.object({ action: z.literal("decide-return"), organizationId: organization, receiptId: identifier, inventoryAction: z.enum(["quarantine", "restock", "repair", "scrap"]), notes: z.string().trim().min(1).max(1000) }).strict()
]);

function sameOrigin(request: NextRequest) {
  const supplied = request.headers.get("origin");
  if (!supplied) return false;
  try { return supplied === applicationBaseUrl().origin; } catch { return false; }
}

export async function POST(request: NextRequest) {
  if (!sameOrigin(request)) return NextResponse.json({ code: "CROSS_ORIGIN_REJECTED" }, { status: 403 });
  const session = await readSession();
  if (!session) return NextResponse.json({ code: "UNAUTHENTICATED" }, { status: 401 });
  if (!["appointment:manage", "availability:manage", "customer:self", "handover:manage", "inventory:allocate", "lead:assign", "lead:update", "payment:create", "quote:write", "resource:manage"].some(permission => allowed(session,permission))) return NextResponse.json({ code: "FORBIDDEN" }, { status: 403 });
  if (request.headers.get("content-type") !== "application/json") return NextResponse.json({ code: "UNSUPPORTED_MEDIA_TYPE" }, { status: 415 });
  const length = Number(request.headers.get("content-length") ?? "0");
  if (!Number.isSafeInteger(length) || length < 0 || length > 65_536) return NextResponse.json({ code: "COMMAND_TOO_LARGE" }, { status: 413 });
  let text: string;
  try { text = await boundedCommandBody(request,65_536); } catch(e) { return NextResponse.json({code:e instanceof Error && e.message==="BODY_TOO_LARGE"?"COMMAND_TOO_LARGE":"INVALID_BODY"},{status:e instanceof Error && e.message==="BODY_TOO_LARGE"?413:400}); }
  let json: unknown;
  try { json = JSON.parse(text); } catch { return NextResponse.json({ code: "INVALID_JSON" }, { status: 400 }); }
  const parsed = commandSchema.safeParse(json);
  if (!parsed.success) return NextResponse.json({ code: "INVALID_COMMAND" }, { status: 400 });
  const command = parsed.data;
  if (!session.organizations.includes(command.organizationId)) return NextResponse.json({ code: "ORGANIZATION_FORBIDDEN" }, { status: 403 });

  let permission: string;
  let path: string;
  let body: Record<string, unknown>;
  switch (command.action) {
    case "request-order-payment":
      permission="payment:create";
      path=`/v1/commerce/orders/${encodeURIComponent(command.orderId)}/payment-request`;
      body={organization_id:command.organizationId,request_key:command.requestKey};
      break;
    case "allocate-order-stock":
      permission="inventory:allocate";
      path=`/v1/commerce/orders/${encodeURIComponent(command.orderId)}/allocations`;
      body={organization_id:command.organizationId,line_id:command.lineId,stock_unit_id:command.stockUnitId,order_version:command.orderVersion,stock_version:command.stockVersion};
      break;
    case "assign-lead":
      permission = "lead:assign";
      path = `/v1/franchise/leads/${encodeURIComponent(command.leadId)}/assign`;
      body = { organization_id: command.organizationId, assigned_subject: command.assignedSubject, version: command.version };
      break;
    case "transition-lead":
      permission = "lead:update";
      path = `/v1/franchise/leads/${encodeURIComponent(command.leadId)}/transitions`;
      body = { organization_id: command.organizationId, current: command.current, target: command.target, version: command.version };
      break;
    case "create-quote":
      permission = "quote:write";
      path = "/v1/franchise/quotes";
      body = { organization_id: command.organizationId, lead_id: command.leadId, variant_id: command.variantId, price_book_id: command.priceBookId, valid_until: command.validUntil };
      break;
    case "create-appointment-slot":
      permission = "appointment:manage";
      path = "/v1/franchise/appointment-slots";
      body = { organization_id: command.organizationId, kind: command.kind, starts_at: command.startsAt, ends_at: command.endsAt, capacity: command.capacity };
      break;
    case "create-resource":
      permission = "resource:manage";
      path = "/v1/franchise/resources";
      body = { organization_id: command.organizationId, principal_subject: command.principalSubject ?? "", display_name: command.displayName, kind: command.kind, skills: command.skills };
      break;
    case "create-availability":
      permission = "availability:manage";
      path = "/v1/franchise/availability";
      body = { organization_id: command.organizationId, resource_id: command.resourceId ?? "", entry_type: command.entryType, reason_code: command.reasonCode ?? "", starts_at: command.startsAt, ends_at: command.endsAt };
      break;
    case "cancel-availability":
      permission = "availability:manage";
      path = `/v1/franchise/availability/${encodeURIComponent(command.availabilityId)}/cancel`;
      body = { organization_id: command.organizationId, version: command.version, reason_code: command.reasonCode };
      break;
    case "assign-appointment-resource":
      permission = "appointment:manage";
      path = `/v1/franchise/appointments/${encodeURIComponent(command.appointmentId)}/resources`;
      body = { organization_id: command.organizationId, resource_id: command.resourceId, version: command.version };
      break;
    case "transition-appointment":
      permission = "appointment:manage";
      path = `/v1/franchise/appointments/${encodeURIComponent(command.appointmentId)}/transitions`;
      body = { organization_id: command.organizationId, current: command.current, target: command.target, version: command.version, reason_code: command.reasonCode ?? "" };
      break;
    case "cancel-customer-appointment":
      permission = "customer:self";
      path = `/v1/customer/appointments/${encodeURIComponent(command.appointmentId)}/cancel`;
      body = { organization_id: command.organizationId, version: command.version, reason_code: command.reasonCode };
      break;
    case "accept-quote":
      permission = "customer:self";
      path = `/v1/customer/quotes/${encodeURIComponent(command.quoteId)}/accept`;
      body = { organization_id: command.organizationId, version: command.version };
      break;
    case "publish-delivery-checklist":
      permission = "handover:manage";
      path = "/v1/franchise/delivery-checklists";
      body = { organization_id: command.organizationId, checklist_id: command.checklistId, version: command.version, title: command.title, items: command.items };
      break;
    case "complete-delivery-checklist":
      permission = "handover:manage";
      path = `/v1/franchise/handovers/${encodeURIComponent(command.handoverId)}/complete-checklist`;
      body = { organization_id: command.organizationId, version: command.version, checklist_id: command.checklistId, checklist_version: command.checklistVersion, responses: command.responses };
      break;
    case "accept-handover":
      permission = "customer:self";
      path = `/v1/customer/handovers/${encodeURIComponent(command.handoverId)}/accept`;
      body = { organization_id: command.organizationId, version: command.version, confirmed_received: command.confirmedReceived, serial_number: command.serialNumber, checklist_id: command.checklistId, checklist_version: command.checklistVersion };
      break;
    case "reject-handover":
      permission = "customer:self";
      path = `/v1/customer/handovers/${encodeURIComponent(command.handoverId)}/reject`;
      body = { organization_id: command.organizationId, version: command.version, reason_code: command.reasonCode, details: command.details };
      break;
    case "resolve-delivery-exception":
      permission = "handover:manage";
      path = `/v1/franchise/delivery-exceptions/${encodeURIComponent(command.exceptionId)}/resolve`;
      body = { organization_id: command.organizationId, version: command.version, action: command.resolutionAction, notes: command.notes };
      break;
    case "receive-return": {
      permission = "handover:manage";
      path = `/v1/franchise/return-authorizations/${encodeURIComponent(command.authorizationId)}/receive`;
      const digest = await crypto.subtle.digest("SHA-256", new TextEncoder().encode(JSON.stringify(command)));
      const evidence = Array.from(new Uint8Array(digest), (byte) => byte.toString(16).padStart(2, "0")).join("");
      body = { organization_id: command.organizationId, serial_number: command.serialNumber, condition_code: command.conditionCode, notes: command.notes, evidence_sha256: evidence };
      break;
    }
    case "decide-return":
      permission = "handover:manage";
      path = `/v1/franchise/return-receipts/${encodeURIComponent(command.receiptId)}/decide`;
      body = { organization_id: command.organizationId, inventory_action: command.inventoryAction, notes: command.notes };
      break;
  }
  if (!allowed(session, permission)) return NextResponse.json({ code: "FORBIDDEN" }, { status: 403 });
  const requestKey = request.headers.get("idempotency-key");
  if (["create-quote","create-availability","create-resource","create-appointment-slot"].includes(command.action) && (!requestKey || !/^[A-Za-z0-9_-]{16,128}$/.test(requestKey))) return NextResponse.json({code:"INVALID_IDEMPOTENCY_KEY"},{status:400});
  try {
    const value = ["create-quote","create-availability","create-resource","create-appointment-slot"].includes(command.action) ? await protectedPost<unknown>(session, path, body, requestKey!) : await protectedPost<unknown>(session, path, body);
    return NextResponse.json(value, { status: ["create-quote", "create-appointment-slot", "create-resource", "create-availability", "publish-delivery-checklist"].includes(command.action) ? 201 : 200 });
  } catch (error) {
    if (error instanceof BackendProblem) return NextResponse.json({ code: error.code }, { status: error.status });
    return NextResponse.json({ code: "UPSTREAM_UNAVAILABLE" }, { status: 502 });
  }
}

export async function GET(request: NextRequest) {
  const session = await readSession();
  if (!session) return NextResponse.json({code:"UNAUTHENTICATED"},{status:401});
  const parameters=Object.fromEntries(new URL(request.url).searchParams);
  if(parameters.kind==="return"){
    const parsed=z.object({kind:z.literal("return"),organizationId:organization,authorizationId:identifier}).strict().safeParse(parameters);
    if(!parsed.success)return NextResponse.json({code:"INVALID_QUERY"},{status:400});
    const q=parsed.data;if(!session.organizations.includes(q.organizationId)||!allowed(session,"handover:manage"))return NextResponse.json({code:"FORBIDDEN"},{status:403});
    try{
      const value=await protectedGet<unknown>(session,"/v1/franchise/returns/result",{organization_id:q.organizationId,authorization_id:q.authorizationId});
      return NextResponse.json({returnCase:value},{headers:{"cache-control":"no-store"}});
    }catch(error){
      if(error instanceof BackendProblem)return NextResponse.json({code:error.code},{status:error.status,headers:{"cache-control":"no-store"}});
      return NextResponse.json({code:"UPSTREAM_UNAVAILABLE"},{status:502,headers:{"cache-control":"no-store"}});
    }
  }
  if(parameters.kind==="checklist-completion"){
    const parsed=z.object({kind:z.literal("checklist-completion"),organizationId:organization,handoverId:identifier}).strict().safeParse(parameters);
    if(!parsed.success)return NextResponse.json({code:"INVALID_QUERY"},{status:400});
    const q=parsed.data;if(!session.organizations.includes(q.organizationId)||!allowed(session,"handover:manage"))return NextResponse.json({code:"FORBIDDEN"},{status:403});
    try{
      const value=await protectedGet<unknown>(session,`/v1/franchise/handovers/${encodeURIComponent(q.handoverId)}/checklist-result`,{organization_id:q.organizationId});
      return NextResponse.json({completion:value},{headers:{"cache-control":"no-store"}});
    }catch(error){
      if(error instanceof BackendProblem)return NextResponse.json({code:error.code},{status:error.status,headers:{"cache-control":"no-store"}});
      return NextResponse.json({code:"UPSTREAM_UNAVAILABLE"},{status:502,headers:{"cache-control":"no-store"}});
    }
  }
  if(parameters.kind==="checklist"){
    const parsed=z.object({kind:z.literal("checklist"),organizationId:organization,checklistId:z.string().max(64).regex(/^[a-z][a-z0-9]*(-[a-z0-9]+)*$/),version:z.string().regex(/^[1-9][0-9]{0,15}$/).transform(Number).refine(Number.isSafeInteger)}).strict().safeParse(parameters);
    if(!parsed.success)return NextResponse.json({code:"INVALID_QUERY"},{status:400});
    const q=parsed.data;if(!session.organizations.includes(q.organizationId)||!allowed(session,"handover:manage"))return NextResponse.json({code:"FORBIDDEN"},{status:403});
    try{
      const value=await protectedGet<unknown>(session,"/v1/franchise/delivery-checklists/result",{organization_id:q.organizationId,checklist_id:q.checklistId,version:String(q.version)});
      return NextResponse.json({checklist:value},{headers:{"cache-control":"no-store"}});
    }catch(error){
      if(error instanceof BackendProblem)return NextResponse.json({code:error.code},{status:error.status,headers:{"cache-control":"no-store"}});
      return NextResponse.json({code:"UPSTREAM_UNAVAILABLE"},{status:502,headers:{"cache-control":"no-store"}});
    }
  }
  const query = z.object({organizationId:organization,leadId:identifier.optional(),kind:z.enum(["quote","availability","resource","slot"]).default("quote"),requestKey:z.string().regex(/^[A-Za-z0-9_-]{16,128}$/)}).strict().refine(v=>v.kind==="quote"?!!v.leadId:v.leadId===undefined).safeParse(Object.fromEntries(new URL(request.url).searchParams));
  if (!query.success) return NextResponse.json({code:"INVALID_QUERY"},{status:400});
  const {organizationId,leadId,requestKey,kind}=query.data;
  if (!session.organizations.includes(organizationId) || !allowed(session,kind==="quote"?"quote:write":kind==="availability"?"availability:manage":kind==="resource"?"resource:manage":"appointment:manage")) return NextResponse.json({code:"FORBIDDEN"},{status:403});
  try {
    const value=await protectedGet<unknown>(session,kind==="quote"?"/v1/franchise/quotes/result":kind==="availability"?"/v1/franchise/availability/result":kind==="resource"?"/v1/franchise/resources/result":"/v1/franchise/appointment-slots/result",{organization_id:organizationId,...(kind==="quote"?{lead_id:leadId}:{}),request_key:requestKey});
    return NextResponse.json({request_key:requestKey,[kind]:value},{headers:{"cache-control":"no-store"}});
  } catch(error) {
    if(error instanceof BackendProblem)return NextResponse.json({code:error.code},{status:error.status,headers:{"cache-control":"no-store"}});
    return NextResponse.json({code:"UPSTREAM_UNAVAILABLE"},{status:502,headers:{"cache-control":"no-store"}});
  }
}
````

### FILE: `src/app/api/enterprise/franchise/commands/route.test.ts`

```yaml
block_id: "TS-FRANCHISE-JOURNEY:file:07"
operation: CREATE
provenance: AUTHORED
source: "local security and mapping regressions"
license: "LicenseRef-Workspace-Owner"
sha256: "5e190d894b1e396b6d99e86e321b2216cc2cf64b24c21ea7b1f160b5fe0f086e"
variables: []
secrets_allowed: false
```

````ts
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { NextRequest } from "next/server";

const mocks = vi.hoisted(() => ({ readSession: vi.fn(), protectedPost: vi.fn(), protectedGet: vi.fn() }));
vi.mock("@/platform/auth/session", () => ({ readSession: mocks.readSession, allowed: (session: { permissions: string[] }, permission: string) => session.permissions.includes(permission) }));
vi.mock("@/platform/backend/protected-client", () => ({ protectedPost: mocks.protectedPost, protectedGet:mocks.protectedGet }));

import { GET, POST } from "./route";

const session = { subject: "operator", tenantId: "tenant", organizations: ["store"], permissions: ["lead:assign", "lead:update", "quote:write", "appointment:manage", "resource:manage", "availability:manage", "handover:manage", "customer:self"], accessToken: "token" };

function request(body: unknown, origin = "https://portal.example") {
  return new NextRequest("https://portal.example/api/enterprise/franchise/commands", { method: "POST", headers: { origin, "content-type": "application/json" }, body: JSON.stringify(body) });
}

describe("franchise command BFF", () => {
  beforeEach(() => { vi.clearAllMocks(); vi.stubEnv("APP_BASE_URL","https://portal.example"); mocks.readSession.mockResolvedValue(session); mocks.protectedPost.mockResolvedValue({ version: 2 }); });
  afterEach(() => vi.unstubAllEnvs());

  it("recovers a published checklist by its natural scoped version without POST",async()=>{
    const url="https://portal.example/api/enterprise/franchise/commands?kind=checklist&organizationId=store&checklistId=fixture-checklist&version=1";mocks.protectedGet.mockResolvedValue({id:"fixture-checklist",version:1});
    expect((await GET(new NextRequest(url))).status).toBe(200);expect(mocks.protectedGet).toHaveBeenCalledWith(session,"/v1/franchise/delivery-checklists/result",{organization_id:"store",checklist_id:"fixture-checklist",version:"1"});expect(mocks.protectedPost).not.toHaveBeenCalled();
    mocks.protectedGet.mockClear();for(const bad of [url+"&requestKey=unknown-key",url.replace("version=1","version=0"),url.replace("version=1","version=9007199254740993")])expect((await GET(new NextRequest(bad))).status).toBe(400);expect(mocks.protectedGet).not.toHaveBeenCalled();
    mocks.readSession.mockResolvedValue({...session,permissions:["resource:manage"]});expect((await GET(new NextRequest(url))).status).toBe(403);expect(mocks.protectedGet).not.toHaveBeenCalled();
  });

  it("retains slot identity in POST and uses scoped GET recovery",async()=>{
    const command={action:"create-appointment-slot",organizationId:"store",kind:"service",startsAt:"2030-01-01T10:00:00Z",endsAt:"2030-01-01T11:00:00Z",capacity:2};
    expect((await POST(request(command))).status).toBe(400);expect(mocks.protectedPost).not.toHaveBeenCalled();
    const keyed=request(command);keyed.headers.set("idempotency-key","slot-fixture-key");expect((await POST(keyed)).status).toBe(201);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session,"/v1/franchise/appointment-slots",{organization_id:"store",kind:"service",starts_at:command.startsAt,ends_at:command.endsAt,capacity:2},"slot-fixture-key");
    mocks.protectedPost.mockClear();mocks.protectedGet.mockResolvedValue({id:"resource"});
    const url="https://portal.example/api/enterprise/franchise/commands?kind=slot&organizationId=store&requestKey=slot-fixture-key";
    expect((await GET(new NextRequest(url))).status).toBe(200);expect(mocks.protectedGet).toHaveBeenCalledWith(session,"/v1/franchise/appointment-slots/result",{organization_id:"store",request_key:"slot-fixture-key"});expect(mocks.protectedPost).not.toHaveBeenCalled();
    mocks.protectedGet.mockClear();mocks.readSession.mockResolvedValue({...session,permissions:["availability:manage"]});expect((await GET(new NextRequest(url))).status).toBe(403);expect(mocks.protectedGet).not.toHaveBeenCalled();
  });
  it("retains resource identity in POST and uses scoped GET recovery",async()=>{
    const command={action:"create-resource",organizationId:"store",displayName:"Synthetic bay",kind:"service-bay",skills:["service"]};
    expect((await POST(request(command))).status).toBe(400);expect(mocks.protectedPost).not.toHaveBeenCalled();
    const keyed=request(command);keyed.headers.set("idempotency-key","resource-fixture-key");expect((await POST(keyed)).status).toBe(201);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session,"/v1/franchise/resources",{organization_id:"store",principal_subject:"",display_name:"Synthetic bay",kind:"service-bay",skills:["service"]},"resource-fixture-key");
    mocks.protectedPost.mockClear();mocks.protectedGet.mockResolvedValue({id:"resource"});
    const url="https://portal.example/api/enterprise/franchise/commands?kind=resource&organizationId=store&requestKey=resource-fixture-key";
    expect((await GET(new NextRequest(url))).status).toBe(200);expect(mocks.protectedGet).toHaveBeenCalledWith(session,"/v1/franchise/resources/result",{organization_id:"store",request_key:"resource-fixture-key"});expect(mocks.protectedPost).not.toHaveBeenCalled();
    mocks.protectedGet.mockClear();mocks.readSession.mockResolvedValue({...session,permissions:["availability:manage"]});expect((await GET(new NextRequest(url))).status).toBe(403);expect(mocks.protectedGet).not.toHaveBeenCalled();
  });

  it("maps an initial payment with a stable key but without browser-owned money/provider",async()=>{
    const operator={...session,permissions:["payment:create"]};mocks.readSession.mockResolvedValue(operator);
    const command={action:"request-order-payment",organizationId:"store",orderId:"order",requestKey:"11111111-1111-4111-8111-111111111111"};
    expect((await POST(request(command))).status).toBe(200);
    expect(mocks.protectedPost).toHaveBeenCalledWith(operator,"/v1/commerce/orders/order/payment-request",{organization_id:"store",request_key:command.requestKey});
    vi.clearAllMocks();mocks.readSession.mockResolvedValue(operator);
    for(const bad of [{...command,provider:"stripe"},{...command,amount:1},{...command,requestKey:"invalid"},{...command,actor:"owner"}])expect((await POST(request(bad))).status).toBe(400);
    expect((await POST(request(command,"https://other.invalid"))).status).toBe(403);
    expect((await POST(request({...command,organizationId:"other"}))).status).toBe(403);
    mocks.readSession.mockResolvedValue({...operator,permissions:["admin:read"]});expect((await POST(request(command))).status).toBe(403);
    expect(mocks.protectedPost).not.toHaveBeenCalled();
  });

  it("maps an allocation to the existing stock owner without server-owned fields", async()=>{
    const operator={...session,permissions:["inventory:allocate"]};mocks.readSession.mockResolvedValue(operator);
    const command={action:"allocate-order-stock",organizationId:"store",orderId:"order",lineId:"line",stockUnitId:"unit",orderVersion:3,stockVersion:2};
    expect((await POST(request(command))).status).toBe(200);
    expect(mocks.protectedPost).toHaveBeenCalledWith(operator,"/v1/commerce/orders/order/allocations",{organization_id:"store",line_id:"line",stock_unit_id:"unit",order_version:3,stock_version:2});
    vi.clearAllMocks();
    for(const invalid of [{...command,state:"reserved"},{...command,orderVersion:0},{...command,stockVersion:1.5},{...command,stockUnitId:""}]) expect((await POST(request(invalid))).status).toBe(400);
    expect(mocks.protectedPost).not.toHaveBeenCalled();
  });
  it("does not admit an allocation for a read-only or unrelated organization session",async()=>{
    const command={action:"allocate-order-stock",organizationId:"store",orderId:"order",lineId:"line",stockUnitId:"unit",orderVersion:1,stockVersion:1};
    mocks.readSession.mockResolvedValue({...session,permissions:["admin:read"]});
    expect((await POST(request(command))).status).toBe(403);
    mocks.readSession.mockResolvedValue({...session,organizations:["other"],permissions:["inventory:allocate"]});
    expect((await POST(request(command))).status).toBe(403);
    expect(mocks.protectedPost).not.toHaveBeenCalled();
  });

  it("uses the configured public origin behind an internal proxy", async () => {
    const input=new NextRequest("http://localhost:4173/api/enterprise/franchise/commands",{method:"POST",headers:{origin:"https://portal.example","content-type":"application/json",host:"localhost:4173"},body:JSON.stringify({action:"assign-appointment-resource",organizationId:"store",appointmentId:"a",resourceId:"r",version:1})});
    expect((await POST(input)).status).toBe(200);
  });

  it("does not trust forwarded host or an absent configured public origin", async () => {
    const input=request({action:"assign-lead"},"https://evil.example");
    input.headers.set("x-forwarded-host","evil.example");
    expect((await POST(input)).status).toBe(403);
    vi.stubEnv("APP_BASE_URL","");
    expect((await POST(request({action:"assign-lead"}))).status).toBe(403);
    expect(mocks.readSession).not.toHaveBeenCalled();
  });

  it("rejects cross-origin before session or backend", async () => {
    const response = await POST(request({ action: "assign-lead" }, "https://evil.example"));
    expect(response.status).toBe(403);
    expect(mocks.readSession).not.toHaveBeenCalled();
  });

  it("maps a scoped assignment to the existing Go command", async () => {
    const response = await POST(request({ action: "assign-lead", organizationId: "store", leadId: "lead-1", assignedSubject: "seller", version: 1 }));
    expect(response.status).toBe(200);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/franchise/leads/lead-1/assign", { organization_id: "store", assigned_subject: "seller", version: 1 });
  });

  it("rejects organization and permission escalation", async () => {
    expect((await POST(request({ action: "assign-lead", organizationId: "other", leadId: "lead-1", assignedSubject: "seller", version: 1 }))).status).toBe(403);
    mocks.readSession.mockResolvedValue({ ...session, permissions: [] });
    expect((await POST(request({ action: "assign-lead", organizationId: "store", leadId: "lead-1", assignedSubject: "seller", version: 1 }))).status).toBe(403);
    expect(mocks.protectedPost).not.toHaveBeenCalled();
  });

  it("maps a quote without accepting client price or currency", async () => {
    const req=request({ action: "create-quote", organizationId: "store", leadId: "lead-1", variantId: "variant", priceBookId: "retail", validUntil: "2026-09-01T12:00:00.000Z" });
    req.headers.set("idempotency-key","quote-request-key-0001");
    const response = await POST(req);
    expect(response.status).toBe(201);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/franchise/quotes", { organization_id: "store", lead_id: "lead-1", variant_id: "variant", price_book_id: "retail", valid_until: "2026-09-01T12:00:00.000Z" }, "quote-request-key-0001");
  });

  it("maps bounded appointment capacity to the existing Go owner", async () => {
    const slotRequest = request({ action: "create-appointment-slot", organizationId: "store", kind: "service", startsAt: "2026-09-01T12:00:00.000Z", endsAt: "2026-09-01T13:00:00.000Z", capacity: 2 });
slotRequest.headers.set("idempotency-key","slot-positive-fixture-key");const response=await POST(slotRequest);
    expect(response.status).toBe(201);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/franchise/appointment-slots", { organization_id: "store", kind: "service", starts_at: "2026-09-01T12:00:00.000Z", ends_at: "2026-09-01T13:00:00.000Z", capacity: 2 },"slot-positive-fixture-key");
  });

  it("maps availability and resource commands without accepting server-owned fields", async () => {
    const availabilityRequest=request({ action: "create-availability", organizationId: "store", resourceId: "technician", entryType: "unavailable", reasonCode: "annual-leave", startsAt: "2026-09-01T12:00:00.000Z", endsAt: "2026-09-01T13:00:00.000Z" });
    availabilityRequest.headers.set("idempotency-key","availability-request-0001");
    const availability=await POST(availabilityRequest);
    expect(availability.status).toBe(201);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/franchise/availability", { organization_id: "store", resource_id: "technician", entry_type: "unavailable", reason_code: "annual-leave", starts_at: "2026-09-01T12:00:00.000Z", ends_at: "2026-09-01T13:00:00.000Z" }, "availability-request-0001");
    vi.clearAllMocks();
    const injected = await POST(request({ action: "create-availability", organizationId: "store", entryType: "working", startsAt: "2026-09-01T12:00:00.000Z", endsAt: "2026-09-01T13:00:00.000Z", state: "active", version: 99 }));
    expect(injected.status).toBe(400);
    expect(mocks.protectedPost).not.toHaveBeenCalled();
  });

  it("maps audited appointment transition and customer-owned cancellation", async () => {
    const transition = await POST(request({ action: "transition-appointment", organizationId: "store", appointmentId: "appointment-1", current: "confirmed", target: "no-show", version: 3, reasonCode: "customer-absent" }));
    expect(transition.status).toBe(200);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/franchise/appointments/appointment-1/transitions", { organization_id: "store", current: "confirmed", target: "no-show", version: 3, reason_code: "customer-absent" });
    vi.clearAllMocks();
    const cancellation = await POST(request({ action: "cancel-customer-appointment", organizationId: "store", appointmentId: "appointment-2", version: 1, reasonCode: "customer-request" }));
    expect(cancellation.status).toBe(200);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/customer/appointments/appointment-2/cancel", { organization_id: "store", version: 1, reason_code: "customer-request" });
  });

  it("maps customer quote acceptance without client price, currency or evidence digest", async () => {
    const response = await POST(request({ action: "accept-quote", organizationId: "store", quoteId: "quote-1", version: 1 }));
    expect(response.status).toBe(200);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/customer/quotes/quote-1/accept", { organization_id: "store", version: 1 });
    vi.clearAllMocks();
    const injected = await POST(request({ action: "accept-quote", organizationId: "store", quoteId: "quote-1", version: 1, totalMinorUnits: 1, currency: "USD" }));
    expect(injected.status).toBe(400);
    expect(mocks.protectedPost).not.toHaveBeenCalled();
  });

  it("maps handover acceptance without accepting a browser-authored evidence digest", async () => {
    const response = await POST(request({ action: "accept-handover", organizationId: "store", handoverId: "handover-1", version: 2, confirmedReceived: true, serialNumber: "SERIAL-1", checklistId: "standard-delivery", checklistVersion: 1 }));
    expect(response.status).toBe(200);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/customer/handovers/handover-1/accept", { organization_id: "store", version: 2, confirmed_received: true, serial_number: "SERIAL-1", checklist_id: "standard-delivery", checklist_version: 1 });
    vi.clearAllMocks();
    const injected = await POST(request({ action: "accept-handover", organizationId: "store", handoverId: "handover-1", version: 2, confirmedReceived: true, serialNumber: "SERIAL-1", checklistId: "standard-delivery", checklistVersion: 1, evidenceSha256: "a".repeat(64) }));
    expect(injected.status).toBe(400);
    expect(mocks.protectedPost).not.toHaveBeenCalled();
  });

  it("publishes and completes an exact delivery checklist version", async () => {
    const items = [{ id: "serial-observed", prompt: "Verificar serie", response_type: "serial", required: true }];
    const published = await POST(request({ action: "publish-delivery-checklist", organizationId: "store", checklistId: "standard-delivery", version: 1, title: "Entrega estándar", items }));
    expect(published.status).toBe(201);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/franchise/delivery-checklists", { organization_id: "store", checklist_id: "standard-delivery", version: 1, title: "Entrega estándar", items });
    vi.clearAllMocks();
    mocks.readSession.mockResolvedValue(session);
    mocks.protectedPost.mockResolvedValue({ version: 2 });
    const responses = [{ item_id: "serial-observed", response_text: "SERIAL-1" }];
    const completed = await POST(request({ action: "complete-delivery-checklist", organizationId: "store", handoverId: "handover-1", version: 1, checklistId: "standard-delivery", checklistVersion: 1, responses }));
    expect(completed.status).toBe(200);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/franchise/handovers/handover-1/complete-checklist", { organization_id: "store", version: 1, checklist_id: "standard-delivery", checklist_version: 1, responses });
  });

  it("maps customer rejection and operator resolution without browser authority fields", async () => {
    const rejected = await POST(request({ action: "reject-handover", organizationId: "store", handoverId: "handover-1", version: 2, reasonCode: "visible-damage", details: "Rayón visible" }));
    expect(rejected.status).toBe(200);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/customer/handovers/handover-1/reject", { organization_id: "store", version: 2, reason_code: "visible-damage", details: "Rayón visible" });
    vi.clearAllMocks();
    mocks.readSession.mockResolvedValue(session);
    mocks.protectedPost.mockResolvedValue({ exception: { state: "resolved" } });
    const resolved = await POST(request({ action: "resolve-delivery-exception", organizationId: "store", exceptionId: "exception-1", version: 1, resolutionAction: "correct-and-represent", notes: "Corregir preparación" }));
    expect(resolved.status).toBe(200);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/franchise/delivery-exceptions/exception-1/resolve", { organization_id: "store", version: 1, action: "correct-and-represent", notes: "Corregir preparación" });
    vi.clearAllMocks();
    const injected = await POST(request({ action: "resolve-delivery-exception", organizationId: "store", exceptionId: "exception-1", version: 1, resolutionAction: "refund-now", notes: "invalid", stockState: "available" }));
    expect(injected.status).toBe(400);
    expect(mocks.protectedPost).not.toHaveBeenCalled();
  });

  it("maps physical return receipt and disposition without browser-owned effects", async () => {
    const received = await POST(request({ action: "receive-return", organizationId: "store", authorizationId: "authorization-1", serialNumber: "SERIAL-1", conditionCode: "damaged", notes: "Daño confirmado" }));
    expect(received.status).toBe(200);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/franchise/return-authorizations/authorization-1/receive", { organization_id: "store", serial_number: "SERIAL-1", condition_code: "damaged", notes: "Daño confirmado", evidence_sha256: expect.stringMatching(/^[0-9a-f]{64}$/) });
    vi.clearAllMocks();
    mocks.readSession.mockResolvedValue(session);
    mocks.protectedPost.mockResolvedValue({ customer_remedy: "refund" });
    const decided = await POST(request({ action: "decide-return", organizationId: "store", receiptId: "receipt-1", inventoryAction: "quarantine", notes: "Separar y solicitar efectos" }));
    expect(decided.status).toBe(200);
    expect(mocks.protectedPost).toHaveBeenCalledWith(session, "/v1/franchise/return-receipts/receipt-1/decide", { organization_id: "store", inventory_action: "quarantine", notes: "Separar y solicitar efectos" });
    vi.clearAllMocks();
    const injected = await POST(request({ action: "decide-return", organizationId: "store", receiptId: "receipt-1", inventoryAction: "restock", notes: "invalid", refundNow: true, paymentId: "browser-owned" }));
    expect(injected.status).toBe(400);
    expect(mocks.protectedPost).not.toHaveBeenCalled();
  });
});

describe("quote recovery boundary",()=>{
 beforeEach(()=>{vi.clearAllMocks();vi.stubEnv("APP_BASE_URL","https://portal.example");mocks.readSession.mockResolvedValue(session)});
 afterEach(()=>vi.unstubAllEnvs());
 it("rejects absent and malformed quote keys before the backend",async()=>{
  for(const key of [null,"short","a".repeat(129),"has invalid spaces"]){
   const req=request({action:"create-quote",organizationId:"store",leadId:"lead",variantId:"variant",priceBookId:"book",validUntil:"2026-09-20T12:00:00Z"});if(key!==null)req.headers.set("idempotency-key",key);
   expect((await POST(req)).status).toBe(400);
  }expect(mocks.protectedPost).not.toHaveBeenCalled();
 });
 it("uses scoped GET without posting and conserves unavailable results",async()=>{
  const req=new NextRequest("https://portal.example/api/enterprise/franchise/commands?organizationId=store&leadId=lead&requestKey=quote-request-key-0001");
  mocks.protectedGet.mockResolvedValue({id:"q"});const response=await GET(req);expect(response.status).toBe(200);expect(response.headers.get("cache-control")).toBe("no-store");
  expect(mocks.protectedGet).toHaveBeenCalledWith(session,"/v1/franchise/quotes/result",{organization_id:"store",lead_id:"lead",request_key:"quote-request-key-0001"});
  mocks.protectedGet.mockRejectedValue(new Error("synthetic private detail"));const failed=await GET(req);expect(failed.status).toBe(502);expect(await failed.text()).not.toContain("private");
  mocks.readSession.mockResolvedValue({...session,permissions:["lead:read"]});expect((await GET(req)).status).toBe(403);
  mocks.readSession.mockResolvedValue(null);expect((await GET(req)).status).toBe(401);expect(mocks.protectedPost).not.toHaveBeenCalled();
 });
});

describe("availability creation recovery BFF",()=>{
 beforeEach(()=>{vi.clearAllMocks();vi.stubEnv("APP_BASE_URL","https://portal.example");mocks.readSession.mockResolvedValue(session)});afterEach(()=>vi.unstubAllEnvs());
 it("reads an exact availability reference with manage permission and no write",async()=>{
  const req=new NextRequest("https://portal.example/api/enterprise/franchise/commands?kind=availability&organizationId=store&requestKey=availability-request-0001");
  mocks.protectedGet.mockResolvedValue({id:"entry"});expect((await GET(req)).status).toBe(200);
  expect(mocks.protectedGet).toHaveBeenCalledWith(session,"/v1/franchise/availability/result",{organization_id:"store",request_key:"availability-request-0001"});
  mocks.readSession.mockResolvedValue({...session,permissions:["availability:read"]});expect((await GET(req)).status).toBe(403);expect(mocks.protectedPost).not.toHaveBeenCalled();
 });
 it("rejects missing creation identity before forwarding",async()=>{
  const req=request({action:"create-availability",organizationId:"store",entryType:"working",startsAt:"2026-09-20T12:00:00Z",endsAt:"2026-09-20T13:00:00Z"});
  expect((await POST(req)).status).toBe(400);expect(mocks.protectedPost).not.toHaveBeenCalled();
 });
});


describe("checklist completion recovery boundary",()=>{
 beforeEach(()=>{vi.clearAllMocks();mocks.readSession.mockResolvedValue({...session,permissions:["handover:manage"]})});
 it("reads exact completion without posting and rejects scope or extra fields",async()=>{
  const url="https://portal.example/api/enterprise/franchise/commands?kind=checklist-completion&organizationId=store&handoverId=handover";
  mocks.protectedGet.mockResolvedValue({handover_id:"handover",state:"presented"});const result=await GET(new NextRequest(url));expect(result.status).toBe(200);expect(result.headers.get("cache-control")).toBe("no-store");
  expect(mocks.protectedGet).toHaveBeenCalledWith({...session,permissions:["handover:manage"]},"/v1/franchise/handovers/handover/checklist-result",{organization_id:"store"});
  expect((await GET(new NextRequest(url+"&version=1"))).status).toBe(400);expect((await GET(new NextRequest(url.replace("store","other")))).status).toBe(403);
  mocks.protectedGet.mockRejectedValue(new Error("private synthetic detail"));const failed=await GET(new NextRequest(url));expect(failed.status).toBe(502);expect(await failed.text()).not.toContain("private");
  mocks.readSession.mockResolvedValue({...session,permissions:["admin:read"]});expect((await GET(new NextRequest(url))).status).toBe(403);expect(mocks.protectedPost).not.toHaveBeenCalled();
 });
});


describe("return recovery boundary",()=>{
 beforeEach(()=>{vi.clearAllMocks();vi.stubEnv("APP_BASE_URL","https://portal.example")});afterEach(()=>vi.unstubAllEnvs());
 it("uses an exact authorization with management scope, no writes and redacted failures",async()=>{
  const caller={...session,permissions:["handover:manage"]};mocks.readSession.mockResolvedValue(caller);mocks.protectedGet.mockResolvedValue({authorization_id:"authorization"});
  const url="https://portal.example/api/enterprise/franchise/commands?kind=return&organizationId=store&authorizationId=authorization";
  const result=await GET(new NextRequest(url));expect(result.status).toBe(200);expect(result.headers.get("cache-control")).toBe("no-store");expect(mocks.protectedGet).toHaveBeenCalledWith(caller,"/v1/franchise/returns/result",{organization_id:"store",authorization_id:"authorization"});
  expect((await GET(new NextRequest(url+"&receiptId=receipt"))).status).toBe(400);expect((await GET(new NextRequest(url.replace("&authorizationId=authorization","")))).status).toBe(400);expect((await GET(new NextRequest(url.replace("organizationId=store","organizationId=other")))).status).toBe(403);
  mocks.protectedGet.mockRejectedValue(new Error("synthetic private failure"));const failed=await GET(new NextRequest(url));expect(failed.status).toBe(502);expect(await failed.text()).not.toContain("private");
  mocks.readSession.mockResolvedValue({...session,permissions:["admin:read"]});expect((await GET(new NextRequest(url))).status).toBe(403);mocks.readSession.mockResolvedValue(null);expect((await GET(new NextRequest(url))).status).toBe(401);expect(mocks.protectedPost).not.toHaveBeenCalled();
 });
});
````

### FILE: `src/components/franchise-command-panel.tsx`

```yaml
block_id: "TS-FRANCHISE-JOURNEY:file:08"
operation: CREATE
provenance: AUTHORED
source: "local operator command UI over the canonical BFF"
license: "LicenseRef-Workspace-Owner"
sha256: "408e4ef0d8f40fd51b81d6b837cf6f3cf3e7f80cb5aeb819ca47330fd28bae48"
variables: []
secrets_allowed: false
```

````tsx
"use client";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

import { OperationalGuide } from "@/components/operational-guide";
import { OPERATIONAL_GUIDES } from "@/platform/help/content";


import { useEffect, useRef, useState, useTransition, type FormEvent } from "react";
import { useRouter } from "next/navigation";
import { AppointmentNotificationStatus } from "./appointment-notification-status";

type Lead = { id: string; organization_id: string; state: string; source_code: string; assigned_subject?: string; version: number };
type Availability = { id: string; resource_id?: string; entry_type: "working" | "unavailable"; reason_code?: string; starts_at: string; ends_at: string; state: string; version: number };
type DeliveryException = { id: string; handover_id: string; customer_subject: string; reason_code: string; details: string; state: string; version: number; resolution_action?: string };
type ReturnEffect = { id: string; effect_kind: "inventory" | "refund" | "exchange" | "accounting" | "fiscal"; owner_context: string; state: "requested" };
type ReturnCase = { authorization_id: string; order_id: string; stock_unit_id: string; customer_subject: string; authorized_action: "return" | "exchange"; receipt?: { id: string; received_serial_number: string; condition_code: string; received_at: string }; disposition?: { id: string; inventory_action: string; customer_remedy: string; effect_requests: ReturnEffect[] } };
type Command = Record<string, unknown> & { action: string };
type ChecklistDraftItem = { key: number; id: string; prompt: string; responseType: "confirmation" | "text" | "serial"; required: boolean };

const transitions: Record<string, string[]> = { new: ["contacted", "lost"], contacted: ["qualified", "lost"], qualified: ["converted", "lost"] };

export type OrderOperations = {
  payment_provider?: string;
  orders: { id: string; state: string; version: number; lines: { id: string; variant_id: string; name: string; quantity: number; stock_id: string }[]; payments: { id: string; state: string }[]; handovers: { id: string; state: string }[] }[];
  stock: { id: string; variant_id: string; serial: string; version: number }[];
  truncated: boolean;
};

function OrderPaymentRequest({order,organization,provider,scope,disabled}:{order:OrderOperations["orders"][number];organization:string;provider:string;scope:string;disabled:boolean}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

  const [message,setMessage]=useState("");const [locked,setLocked]=useState(false);const inFlight=useRef(false);
  async function requestPayment(){
    if(disabled||locked||inFlight.current)return;
    inFlight.current=true;setLocked(true);setMessage(t("p0308"));
    let key:string;
    try{
      const storageKey=`elite-payment-request:${scope}:${order.id}`;
      key=window.sessionStorage.getItem(storageKey)??crypto.randomUUID();
      if(!/^[a-f0-9-]{36}$/i.test(key))throw new Error("invalid retained key");
      window.sessionStorage.setItem(storageKey,key);
    }catch{setMessage(t("p0309"));return}
    try{
      const response=await fetch("/api/enterprise/franchise/commands",{method:"POST",headers:{"content-type":"application/json"},signal:AbortSignal.timeout(10000),body:JSON.stringify({action:"request-order-payment",organizationId:organization,orderId:order.id,requestKey:key})});
      const value=await response.json() as {id?:unknown;order_id?:unknown;organization_id?:unknown;provider_code?:unknown;state?:unknown;version?:unknown};
      if(response.ok&&typeof value.id==="string"&&value.id.length>0&&value.order_id===order.id&&value.organization_id===organization&&value.provider_code===provider&&typeof value.state==="string"&&["created","pending","authorized","captured","failed","refunded","disputed"].includes(value.state)&&typeof value.version==="number"&&Number.isSafeInteger(value.version)&&value.version>0){
        setMessage(t("p0310", {request_id:(value.id)}));
      }else if(response.status===409){setMessage(t("p0311"))}
      else if(response.status===401||response.status===403){setMessage(t("p0312"))}
      else throw new Error("unconfirmed");
    }catch{setMessage(t("p0313"))}
  }
  return <div><button className="button" type="button" disabled={disabled||locked} onClick={()=>void requestPayment()}>{t("p0314")}</button><p>{t("p0315")} {provider}{t("p0316")}</p><p role="status" aria-label={`Resultado de pago ${order.id}`} aria-live="polite">{message}</p></div>;
}

export function OrderOperationsPanel({ snapshot, organization, canAllocate, canRequestPayment=false, paymentScope="" }: { snapshot: OrderOperations | null; organization: string; canAllocate: boolean;canRequestPayment?:boolean;paymentScope?:string }) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

  const [ready,setReady]=useState(false);
  const [message,setMessage]=useState("");
  const [locked,setLocked]=useState(false);
  const inFlight=useRef(false);
  useEffect(()=>setReady(true),[]);
  async function allocate(event: FormEvent<HTMLFormElement>, order: OrderOperations["orders"][number], line: OrderOperations["orders"][number]["lines"][number]) {
    event.preventDefault();
    if(inFlight.current || locked || !snapshot || snapshot.truncated || !canAllocate) return;
    const stock=snapshot.stock.find(item=>item.id===new FormData(event.currentTarget).get("stock") && item.variant_id===line.variant_id);
    if(!stock) {setMessage(t("p0317"));return;}
    inFlight.current=true;setLocked(true);setMessage(t("p0318"));
    try {
      const response=await fetch("/api/enterprise/franchise/commands",{method:"POST",headers:{"content-type":"application/json"},signal:AbortSignal.timeout(10000),body:JSON.stringify({action:"allocate-order-stock",organizationId:organization,orderId:order.id,lineId:line.id,stockUnitId:stock.id,orderVersion:order.version,stockVersion:stock.version})});
      const result=await response.json() as {status?:unknown};
      if(response.ok && result.status==="accepted") setMessage(t("p0319"));
      else if(response.status===409) setMessage(t("p0320"));
      else if(response.status===401 || response.status===403) setMessage(t("p0321"));
      else throw new Error("unconfirmed");
    } catch {setMessage(t("p0322"));}
    // Remain fenced after every outcome. Recovery is a fresh server-side GET,
    // never an automatic replay of a possibly committed command.
  }
  const state=(value:string)=>({placed:t("p0323"),confirmed:t("p0324"),delivered:t("p0325"),cancelled:t("p0326"),created:t("p0327"),pending:t("p0328"),authorized:t("p0329"),captured:t("p0330"),failed:t("p0331"),refunded:t("p0332"),disputed:t("p0333"),prepared:t("p0334"),presented:t("p0335"),accepted:t("p0336"),rejected:t("p0337")}[value] ?? t("p0338"));
  return <section className="card" aria-labelledby="order-operations-title">
    <h2 id="order-operations-title">{t("p0339")}</h2>
    <p>{t("p0340")}</p>
    <button className="button" type="button" disabled={!ready} onClick={()=>window.location.reload()}>{t("p0341")}</button>
    <p role="status" aria-label={t("p0342")} aria-live="polite">{message}</p>
    {!snapshot ? <p role="alert">{t("p0343")}</p> : <>
      {snapshot.truncated ? <p role="alert">{t("p0344")}</p> : null}
      {snapshot.orders.length===0 ? <p>{t("p0345")}</p> : snapshot.orders.map(order=><article className="card" key={order.id} aria-label={`Pedido ${order.id}`}>
        <h3>{t("p0251")} {order.id}</h3><p>{state(order.state)}</p>
        {order.lines.map(line=><div key={line.id}><h4>{line.name} {t("p0346")} {line.quantity}</h4>
          {line.stock_id ? <p>{t("p0347")} <span>{line.stock_id}</span></p> : <><p>{t("p0348")}</p>
          {canAllocate && ["placed","confirmed"].includes(order.state) && line.quantity===1 ? <form onSubmit={event=>void allocate(event,order,line)}>
            <label>{t("p0349")} {line.name}<select name="stock" required defaultValue="" disabled={!ready || locked || snapshot.truncated}><option value="" disabled>{t("p0350")}</option>{snapshot.stock.filter(item=>item.variant_id===line.variant_id).map(item=><option key={item.id} value={item.id}>{item.serial}</option>)}</select></label>
            {!snapshot.stock.some(item=>item.variant_id===line.variant_id) ? <p>{t("p0351")}</p> : null}
            <button className="button" disabled={!ready || locked || snapshot.truncated || !snapshot.stock.some(item=>item.variant_id===line.variant_id)}>{snapshot.stock.some(item=>item.variant_id===line.variant_id) ? t("p0352") : t("p0353")}</button>
          </form> : <p>{t("p0354")}</p>}</>}
        </div>)}
        <p>{t("p0355")} {order.payments.length ? order.payments.map(item=>`${state(item.state)} (${item.id})`).join("; ") : t("p0356")}</p>
        {!snapshot.payment_provider ? <p>{t("p0357")}</p> : canRequestPayment&&paymentScope&&order.payments.length===0&&["placed","confirmed"].includes(order.state) ? <OrderPaymentRequest order={order} organization={organization} provider={snapshot.payment_provider} scope={paymentScope} disabled={!ready||locked||snapshot.truncated}/> : null}
        <p>{t("p0358")} {order.handovers.length ? order.handovers.map(item=>`${state(item.state)} (${item.id})`).join("; ") : t("p0359")}</p>
      </article>)}
    </>}
    <OperationalGuide guide={OPERATIONAL_GUIDES["order-operations-view"]}/>
  </section>;
}

export type AppointmentAgenda = {
  appointments: { id: string; organization_id: string; lead_id: string; kind: string; starts_at: string; ends_at: string; state: string; version: number; resource_id?: string }[];
  resources: { id: string; display_name: string; skills: string[] }[];
  truncated: boolean;
};

export function AppointmentAgendaPanel({ agenda, organization, notificationStatusEnabled = false }: { agenda: AppointmentAgenda; organization: string; notificationStatusEnabled?: boolean }) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

  const router = useRouter();
  const [ready, setReady] = useState(false);
  const [busy, setBusy] = useState(false);
  const [refreshing, startRefresh] = useTransition();
  const [uncertain, setUncertain] = useState(false);
  const [message, setMessage] = useState("");
  useEffect(() => setReady(true), []);
  const disabled = !ready || busy || refreshing || uncertain || agenda.truncated;
  async function apply(appointment: AppointmentAgenda["appointments"][number], command: Command, expectedState: string) {
    setBusy(true); setMessage("");
    try {
      const response = await fetch("/api/enterprise/franchise/commands", { method: "POST", headers: { "content-type": "application/json" }, body: JSON.stringify({ ...command, organizationId: organization, appointmentId: appointment.id, version: appointment.version }) });
      const receipt = await response.json() as { id?: unknown; organization_id?: unknown; state?: unknown; version?: unknown; resource_id?: unknown };
      if (!response.ok) throw new Error("command not confirmed");
      if (receipt.id !== appointment.id || receipt.organization_id !== organization || receipt.state !== expectedState || receipt.version !== appointment.version + 1 || (command.action === "assign-appointment-resource" && receipt.resource_id !== command.resourceId)) throw new Error("receipt differs");
      setMessage(expectedState === "confirmed" ? t("p0360") : t("p0361"));
      startRefresh(() => router.refresh());
    } catch {
      setUncertain(true);
      setMessage(t("p0362"));
    } finally { setBusy(false); }
  }
  const names: Record<string, string> = { requested: t("p0363"), confirmed: t("p0324"), completed: t("p0364"), cancelled: t("p0326"), "no-show": t("p0365") };
  const kinds: Record<string,string> = { consultation: t("p0366"), "test-drive": t("p0367"), delivery: t("p0368"), service: t("p0018") };
  return <section className="card" aria-labelledby="appointment-agenda-title">
    <h2 id="appointment-agenda-title">{t("p0369")}</h2>
    <p>{t("p0370")}</p>
    <button className="button" type="button" disabled={!ready || busy || refreshing} onClick={() => window.location.reload()}>{t("p0371")}</button>
    <p role="status" aria-label={t("p0372")} aria-live="polite">{message}</p>
    {agenda.truncated ? <p role="alert">{t("p0373")}</p> : null}
    {!agenda.appointments.length ? <p>{t("p0374")}</p> : <div className="grid">{agenda.appointments.map((appointment) => <article className="card" key={appointment.id} aria-label={`Turno ${appointment.kind} ${appointment.starts_at}`}>
      <h3>{kinds[appointment.kind] ?? appointment.kind} {t("p0016")} <time dateTime={appointment.starts_at}>{new Date(appointment.starts_at).toISOString().slice(0,16).replace("T"," ")} {t("p0375")}</time></h3>
      <p>{t("p0247")} {names[appointment.state] ?? appointment.state}</p>
      <p>{t("p0167")} {appointment.id}</p>
      <p>{t("p0376")} {agenda.resources.find((item) => item.id === appointment.resource_id)?.display_name ?? (appointment.resource_id ? t("p0377") : t("p0378"))}</p>
      {appointment.state === "requested" && !appointment.resource_id ? <form onSubmit={(event) => { event.preventDefault(); const data = new FormData(event.currentTarget); void apply(appointment, { action: "assign-appointment-resource", resourceId: String(data.get("resourceId") ?? "") }, "requested"); }}>
        <label>{t("p0379")}<select name="resourceId" required defaultValue="" disabled={disabled}><option value="" disabled>{t("p0380")}</option>{agenda.resources.filter((item) => item.skills.includes(appointment.kind)).map((item) => <option key={item.id} value={item.id}>{item.display_name}</option>)}</select></label>
        <button className="button" disabled={disabled}>{t("p0381")}</button>
      </form> : null}
      {appointment.state === "requested" && appointment.resource_id ? <button className="button" disabled={disabled} onClick={() => void apply(appointment, { action: "transition-appointment", current: appointment.state, target: "confirmed" }, "confirmed")}>{t("p0382")}</button> : null}
      {["requested","confirmed"].includes(appointment.state) ? <form onSubmit={(event) => { event.preventDefault(); const data=new FormData(event.currentTarget); const target=String(data.get("target")); const reasonCode=String(data.get("reasonCode") ?? "").trim(); void apply(appointment,{action:"transition-appointment",current:appointment.state,target,...(reasonCode?{reasonCode}:{})},target); }}>
        <label>{t("p0383")}<select name="target" disabled={disabled}><option value="cancelled">{t("p0384")}</option>{appointment.state==="confirmed"?<><option value="completed">{t("p0385")}</option><option value="no-show">{t("p0386")}</option></>:null}</select></label>
        <label>{t("p0387")}<input name="reasonCode" disabled={disabled} pattern="[a-z][a-z0-9]*(-[a-z0-9]+)*" /></label><button className="button" disabled={disabled}>{t("p0388")}</button>
      </form> : null}
      {notificationStatusEnabled ? <AppointmentNotificationStatus key={`${organization}:${appointment.id}`} organization={organization} appointment={appointment.id}/> : null}
    </article>)}</div>}
  </section>;
}

function AvailabilityCancellation({item,organization,scope,canCancel}:{item:Availability;organization:string;scope:string;canCancel:boolean}) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

  const [ready,setReady]=useState(false),[locked,setLocked]=useState(true),[message,setMessage]=useState("");
  const inFlight=useRef(false);
  const key=`elite-availability-cancel:${scope}:${item.id}`;
  useEffect(()=>{
    try {
      if(!/^[a-f0-9]{64}$/.test(scope)||!Number.isSafeInteger(item.version)||item.version<1)throw new Error("invalid scope/version");
      const raw=window.sessionStorage.getItem(key);
      if(raw!==null) {
        const marker=JSON.parse(raw) as {version?:unknown}|null;
        if(!marker||typeof marker!=="object"||Object.keys(marker).join(",")!=="version"||!Number.isSafeInteger(marker.version)||Number(marker.version)<1)throw new Error("invalid marker");
        if(item.state==="cancelled"&&item.version>Number(marker.version)) {
          window.sessionStorage.removeItem(key);
          setMessage(t("p0389"));
        } else setMessage(t("p0390"));
        setLocked(true);
      } else setLocked(item.state!=="active");
      setReady(true);
    } catch {setReady(true);setLocked(true);setMessage(t("p0391"))}
  },[key,scope,item.version,item.state]);
  async function cancel() {
    if(!ready||locked||inFlight.current||!canCancel||item.state!=="active")return;
    inFlight.current=true;setLocked(true);
    try {if(window.sessionStorage.getItem(key)!==null)throw new Error("unresolved");window.sessionStorage.setItem(key,JSON.stringify({version:item.version}))}
    catch {setMessage(t("p0392"));return}
    setMessage(t("p0393"));
    try {
      const response=await fetch("/api/enterprise/franchise/commands",{method:"POST",headers:{"content-type":"application/json"},signal:AbortSignal.timeout(10000),body:JSON.stringify({action:"cancel-availability",organizationId:organization,availabilityId:item.id,version:item.version,reasonCode:"schedule-correction"})});
      const receipt=await response.json() as (Partial<Availability>&{organization_id?:unknown})|null;
      if(response.ok&&receipt?.id===item.id&&receipt.organization_id===organization&&receipt.state==="cancelled"&&receipt.version===item.version+1&&receipt.entry_type===item.entry_type&&receipt.resource_id===item.resource_id&&receipt.starts_at===item.starts_at&&receipt.ends_at===item.ends_at) {
        setMessage(t("p0394"));
      } else if(response.status===401||response.status===403)setMessage(t("p0395"));
      else if(response.status===409)setMessage(t("p0396"));
      else throw new Error("unconfirmed");
    } catch {setMessage(t("p0397"))}
  }
  return <>
    <p>{item.state==="active"?t("p0398"):item.state==="cancelled"?t("p0326"):t("p0338")}</p>
    {canCancel?<button type="button" disabled={!ready||locked} onClick={()=>void cancel()}>{t("p0399")}</button>:null}
    <p role="status" aria-label={t("p0400")} aria-live="polite">{message}</p>
    <button type="button" disabled={!ready} onClick={()=>window.location.reload()}>{t("p0401")}</button>
    <OperationalGuide guide={OPERATIONAL_GUIDES["availability-cancel-view"]}/>
  </>;
}

function LeadActions({lead,scope,defaultAssignee,canAssign,canTransition}:{lead:Lead;scope:string;defaultAssignee:string;canAssign:boolean;canTransition:boolean}) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

  const [ready,setReady]=useState(false),[locked,setLocked]=useState(true),[message,setMessage]=useState("");
  const inFlight=useRef(false);
  const key=`elite-lead-command:${scope}:${lead.id}`;
  const terminal=lead.state==="converted"||lead.state==="lost";
  useEffect(()=>{
    try {
      if(!/^[a-f0-9]{64}$/.test(scope)||!Number.isSafeInteger(lead.version)||lead.version<1)throw new Error("invalid scope/version");
      const raw=window.sessionStorage.getItem(key);
      if(raw!==null) {
        const marker=JSON.parse(raw) as {version?:unknown;action?:unknown};
        if(!marker||typeof marker!=="object"||Object.keys(marker).sort().join(",")!=="action,version"||!Number.isSafeInteger(marker.version)||Number(marker.version)<1||!["assign-lead","transition-lead"].includes(String(marker.action)))throw new Error("invalid marker");
        if(lead.version>Number(marker.version)) {
          window.sessionStorage.removeItem(key);inFlight.current=false;setLocked(terminal);
          setMessage(t("p0402"));
        } else {setLocked(true);setMessage(t("p0403"))}
      } else {setLocked(terminal)}
      setReady(true);
    } catch {setReady(true);setLocked(true);setMessage(t("p0391"))}
  },[key,scope,lead.version,terminal]);
  async function submit(event:FormEvent<HTMLFormElement>,action:"assign-lead"|"transition-lead") {
    event.preventDefault();if(!ready||locked||inFlight.current||terminal)return;
    if(action==="assign-lead"?!canAssign:!canTransition)return;
    const data=new FormData(event.currentTarget),value=String(data.get(action==="assign-lead"?"assignedSubject":"target")??"").trim();
    if(action==="assign-lead"?(!value||value.length>256):!transitions[lead.state]?.includes(value))return;
    inFlight.current=true;setLocked(true);
    try {if(window.sessionStorage.getItem(key)!==null)throw new Error("unresolved");window.sessionStorage.setItem(key,JSON.stringify({version:lead.version,action}))}
    catch {setMessage(t("p0392"));return}
    setMessage(t("p0404"));
    const command={action,organizationId:lead.organization_id,leadId:lead.id,version:lead.version,...(action==="assign-lead"?{assignedSubject:value}:{current:lead.state,target:value})};
    try {
      const response=await fetch("/api/enterprise/franchise/commands",{method:"POST",headers:{"content-type":"application/json"},signal:AbortSignal.timeout(10000),body:JSON.stringify(command)});
      const receipt=await response.json() as Partial<Lead>|null;
      if(response.ok&&receipt?.id===lead.id&&receipt.organization_id===lead.organization_id&&receipt.version===lead.version+1&&(action==="assign-lead"?receipt.assigned_subject===value&&receipt.state===lead.state:receipt.state===value)) {
        setMessage(t("p0405"));
      } else if(response.status===401||response.status===403) {setMessage(t("p0406"))}
      else if(response.status===409) {setMessage(t("p0407"))}
      else throw new Error("unconfirmed");
    } catch {setMessage(t("p0408"))}
  }
  return <>
    {canAssign?<form onSubmit={event=>void submit(event,"assign-lead")}><label>{t("p0409")}<input name="assignedSubject" defaultValue={lead.assigned_subject||defaultAssignee} required maxLength={256} disabled={!ready||locked}/></label><button className="primaryAction" disabled={!ready||locked}>{t("p0410")}</button></form>:null}
    {canTransition&&transitions[lead.state]?.length?<form onSubmit={event=>void submit(event,"transition-lead")}><label>{t("p0411")}<select name="target" disabled={!ready||locked}>{transitions[lead.state]!.map(target=><option key={target}>{target}</option>)}</select></label><button disabled={!ready||locked}>{t("p0412")}</button></form>:null}
    <p role="status" aria-label={t("p0413")} aria-live="polite">{message}</p>
    <button type="button" disabled={!ready} onClick={()=>window.location.reload()}>{t("p0414")}</button>
    <OperationalGuide guide={OPERATIONAL_GUIDES["lead-command-view"]}/>
  </>;
}

function DeliveryExceptionResolution({item,organization,scope}:{item:DeliveryException;organization:string;scope:string}) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

  // This owner resolves delivery exceptions, not CRM commands.
  const [ready,setReady]=useState(false);
  const [locked,setLocked]=useState(true);
  const [message,setMessage]=useState("");
  const inFlight=useRef(false);
  const key=`elite-delivery-resolution:${scope}:${item.id}`;
  useEffect(()=>{
    try {
      if(!/^[a-f0-9]{64}$/.test(scope))throw new Error("missing scope");
      const raw=window.sessionStorage.getItem(key);
      if(raw!==null) {
        const marker=JSON.parse(raw) as {version?:unknown;action?:unknown};
        if(!marker||typeof marker!=="object"||Object.keys(marker).sort().join(",")!=="action,version"||!Number.isSafeInteger(marker.version)||Number(marker.version)<1||!["correct-and-represent","return","exchange"].includes(String(marker.action)))throw new Error("invalid marker");
        if(item.state==="resolved"&&Number.isSafeInteger(item.version)&&item.version>Number(marker.version)&&["correct-and-represent","return","exchange"].includes(item.resolution_action??"")) {
          setMessage(item.resolution_action===marker.action?t("p0415"):t("p0416"));
          window.sessionStorage.removeItem(key);
        } else {setMessage(t("p0417"))}
        setLocked(true);
      } else {setLocked(item.state!=="open")}
      setReady(true);
    } catch {setLocked(true);setReady(true);setMessage(t("p0391"))}
  },[key,scope,item.state,item.version,item.resolution_action]);
  async function submit(event:FormEvent<HTMLFormElement>) {
    event.preventDefault();if(!ready||locked||inFlight.current||item.state!=="open")return;
    const data=new FormData(event.currentTarget), action=String(data.get("resolutionAction")??""), notes=String(data.get("notes")??"");
    if(!["correct-and-represent","return","exchange"].includes(action)||!notes.trim()||notes.length>1000)return;
    inFlight.current=true;setLocked(true);
    try {
      if(window.sessionStorage.getItem(key)!==null)throw new Error("unresolved");
      window.sessionStorage.setItem(key,JSON.stringify({version:item.version,action}));
    } catch {setMessage(t("p0392"));return}
    setMessage(t("p0418"));
    try {
      const response=await fetch("/api/enterprise/franchise/commands",{method:"POST",headers:{"content-type":"application/json"},signal:AbortSignal.timeout(10000),body:JSON.stringify({action:"resolve-delivery-exception",organizationId:organization,exceptionId:item.id,version:item.version,resolutionAction:action,notes})});
      const value=await response.json() as {exception?:Record<string,unknown>;successor_handover?:Record<string,unknown>;return_authorization_id?:unknown;disposition?:unknown};
      const e=value?.exception;
      const linked=action==="correct-and-represent" ? typeof e?.successor_handover_id==="string"&&e.successor_handover_id!==""&&value.successor_handover?.id===e.successor_handover_id&&value.successor_handover?.supersedes_handover_id===item.handover_id&&value.successor_handover?.state==="prepared" : typeof e?.return_authorization_id==="string"&&e.return_authorization_id!==""&&value.return_authorization_id===e.return_authorization_id&&value.disposition===action;
      if(response.ok&&e?.id===item.id&&e.organization_id===organization&&e.handover_id===item.handover_id&&e.state==="resolved"&&e.version===item.version+1&&e.resolution_action===action&&linked) {
        setMessage(t("p0419"));
      } else if(response.status===401||response.status===403) {setMessage(t("p0420"))}
      else if(response.status===409) {setMessage(t("p0421"))}
      else throw new Error("unconfirmed");
    } catch {setMessage(t("p0422"))}
  }
  return <>
    {item.state==="open"?<form onSubmit={event=>void submit(event)}><label>{t("p0423")}<select name="resolutionAction" required disabled={!ready||locked}><option value="correct-and-represent">{t("p0424")}</option><option value="return">{t("p0425")}</option><option value="exchange">{t("p0426")}</option></select></label><label>{t("p0427")}<textarea name="notes" required maxLength={1000} disabled={!ready||locked}/></label><button disabled={!ready||locked}>{t("p0428")}</button></form>:<p>{t("p0429")} {item.resolution_action}</p>}
    <p role="status" aria-label={t("p0430")} aria-live="polite">{message}</p>
    <button type="button" disabled={!ready} onClick={()=>window.location.reload()}>{t("p0431")}</button>
    <OperationalGuide guide={OPERATIONAL_GUIDES["delivery-resolution-view"]}/>
  </>;
}

export function FranchiseCommandPanel({ permissions = [], leads, availability, exceptions, returnCases, defaultAssignee, organization, commandScope="" }: { permissions?: readonly string[]; leads: Lead[]; availability: Availability[]; exceptions: DeliveryException[] | null; returnCases: ReturnCase[] | null; defaultAssignee: string; organization: string;commandScope?:string }) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

  const can=(permission:string)=>permissions.includes("*")||permissions.includes(permission);
  return <>
    {(can("availability:read")||can("availability:manage")) ? <section className="card"><h2>{t("p0432")}</h2><p>{t("p0433")}</p>{can("availability:manage") ? <AvailabilityCreation organization={organization} scope={commandScope}/> : null}<div className="grid">{availability.map((item) => <article className="card" key={item.id}><h3>{item.resource_id || t("p0434")}</h3><p>{item.entry_type} {t("p0016")} {new Date(item.starts_at).toLocaleString(privateLocale.locale, {timeZone:privateLocale.timeZone})} {t("p0435")} {new Date(item.ends_at).toLocaleString(privateLocale.locale, {timeZone:privateLocale.timeZone})}</p>{item.reason_code ? <p>{t("p0436")} {item.reason_code}</p> : null}<AvailabilityCancellation item={item} organization={organization} scope={commandScope} canCancel={can("availability:manage")}/></article>)}</div></section> : null}
    {can("resource:manage") ? <section className="card"><h2>{t("p0437")}</h2><ResourceCreation organization={organization} scope={commandScope}/></section> : null}
    {can("appointment:manage") ? <section className="card"><h2>{t("p0438")}</h2><SlotCreation organization={organization} scope={commandScope}/></section> : null}
    {can("handover:manage") ? <section className="card"><h2>{t("p0439")}</h2><p>{t("p0440")} <code>confirmed</code>{t("p0441")}</p><ChecklistPublication key={commandScope} organization={organization} scope={commandScope}/><ChecklistCompletionPanel key={commandScope} organization={organization} scope={commandScope}/></section> : null}
    {can("handover:manage") && exceptions!==null ? <section className="card"><h2>{t("p0246")}</h2>{exceptions.length ? <div className="grid">{exceptions.map((item) => <article className="card" key={item.id}><h3>{item.reason_code}</h3><p>{item.details}</p><p>{t("p0368")} {item.handover_id} {t("p0442")} {item.customer_subject} {t("p0443")} {item.state}</p><DeliveryExceptionResolution item={item} organization={organization} scope={commandScope}/></article>)}</div> : <p>{t("p0444")}</p>}</section> : null}
    {can("handover:manage") ? <ReturnOperationsPanel key={commandScope} organization={organization} scope={commandScope} cases={returnCases}/> : null}
    <div className="grid">{leads.map((lead) => <article className="card" key={lead.id}>
      <h2>{lead.id}</h2><p>{lead.state} {t("p0016")} {lead.source_code}</p><p>{lead.assigned_subject ? `Responsable: ${lead.assigned_subject}` : t("p0378")}</p>
      <LeadActions lead={lead} scope={commandScope} defaultAssignee={defaultAssignee} canAssign={can("lead:assign")} canTransition={can("lead:update")}/>
      {can("quote:write") ? <QuoteCreation lead={lead} scope={commandScope}/> : null}
    </article>)}</div>
  </>;
}

function QuoteCreation({lead,scope}:{lead:Lead;scope:string}) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

 const [ready,setReady]=useState(false),[locked,setLocked]=useState(true),[message,setMessage]=useState("");
 const [result,setResult]=useState<{id:string;state:string;currency:string;total_minor_units:number}|null>(null);
 const sending=useRef(false),reading=useRef(false),operation=useRef<string|null>(null);
 const storageKey=`elite-quote-create:${scope}:${lead.id}`;
 useEffect(()=>{
   try {
     if(!/^[a-f0-9]{64}$/.test(scope))throw new Error("scope");
     const raw=sessionStorage.getItem(storageKey);
     if(raw!==null){const marker=JSON.parse(raw);if(!marker||Object.keys(marker).join(",")!=="requestKey"||!/^quote-[a-f0-9-]{36}$/.test(marker.requestKey))throw new Error("marker");operation.current=marker.requestKey;setLocked(true);setMessage(t("p0445"));}
     else {operation.current=null;setLocked(false)}
     setReady(true);
   }catch{setReady(true);setLocked(true);setMessage(t("p0446"))}
 },[scope,storageKey]);
 function valid(value:unknown): value is {id:string;organization_id:string;lead_id:string;variant_id:string;price_book_id:string;valid_until:string;state:string;currency:string;total_minor_units:number;version:number} {
   if(!value||typeof value!=="object")return false;
   const v=value as Record<string,unknown>;
   return typeof v.id==="string"&&v.id.length>0&&v.id.length<=128&&v.organization_id===lead.organization_id&&v.lead_id===lead.id&&typeof v.variant_id==="string"&&typeof v.price_book_id==="string"&&typeof v.valid_until==="string"&&Number.isFinite(Date.parse(v.valid_until))&&["issued","accepted","expired","cancelled"].includes(String(v.state))&&typeof v.currency==="string"&&/^[A-Z]{3}$/.test(v.currency)&&Number.isSafeInteger(v.total_minor_units)&&Number(v.total_minor_units)>=0&&Number.isSafeInteger(v.version)&&Number(v.version)>0;
 }
 async function submit(event:FormEvent<HTMLFormElement>){
   event.preventDefault();if(!ready||locked||sending.current)return;
   const data=new FormData(event.currentTarget),variantId=String(data.get("variantId")??""),priceBookId=String(data.get("priceBookId")??""),until=new Date(String(data.get("validUntil")??""));
   if(!variantId||!priceBookId||!Number.isFinite(until.getTime())||until<=new Date())return;
   sending.current=true;setLocked(true);
   try{if(sessionStorage.getItem(storageKey)!==null)throw new Error("unresolved");const requestKey="quote-"+crypto.randomUUID();sessionStorage.setItem(storageKey,JSON.stringify({requestKey}));operation.current=requestKey;}
   catch{setMessage(t("p0447"));return}
   setMessage(t("p0448"));
   try{
     const response=await fetch("/api/enterprise/franchise/commands",{method:"POST",headers:{"content-type":"application/json","idempotency-key":operation.current!},signal:AbortSignal.timeout(10000),body:JSON.stringify({action:"create-quote",organizationId:lead.organization_id,leadId:lead.id,variantId,priceBookId,validUntil:until.toISOString()})});
     const receipt=await response.json();
     if(!response.ok||!valid(receipt)||receipt.variant_id!==variantId||receipt.price_book_id!==priceBookId||new Date(receipt.valid_until).getTime()!==until.getTime()||receipt.state!=="issued"||receipt.version!==1)throw new Error("unconfirmed");
     setResult(receipt);setMessage(t("p0449"));
   }catch{setMessage(t("p0450"))}
 }
 async function consult(){
   if(!ready||!operation.current||reading.current)return;reading.current=true;
   try{
     const query=new URLSearchParams({organizationId:lead.organization_id,leadId:lead.id,requestKey:operation.current});
     const response=await fetch("/api/enterprise/franchise/commands?"+query,{cache:"no-store",signal:AbortSignal.timeout(10000)});
     const receipt=await response.json();
     if(!response.ok||receipt?.request_key!==operation.current||!valid(receipt.quote))throw new Error("unconfirmed");
     setResult(receipt.quote);setMessage(t("p0451"));
   }catch{setMessage(t("p0452"))}
   finally{reading.current=false}
 }
 return <div>
  <form onSubmit={event=>void submit(event)}>
   <label>{t("p0182")}<input name="variantId" required maxLength={128} disabled={!ready||locked}/></label>
   <label>{t("p0453")}<input name="priceBookId" required maxLength={128} disabled={!ready||locked}/></label>
   <label>{t("p0454")}<input name="validUntil" type="datetime-local" required disabled={!ready||locked}/></label>
   <button disabled={!ready||locked}>{t("p0455")}</button>
  </form>
  <p role="status" aria-label={t("p0456")} aria-live="polite">{message}</p>
  {result?<p>{t("p0457")} {result.id} {t("p0016")} {result.state} {t("p0016")} {result.currency}</p>:null}
  <button type="button" disabled={!ready||!operation.current} onClick={()=>void consult()}>{t("p0458")}</button>
  <OperationalGuide guide={OPERATIONAL_GUIDES["quote-create-view"]}/>
 </div>;
}

function AvailabilityCreation({organization,scope}:{organization:string;scope:string}) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

 const [ready,setReady]=useState(false),[locked,setLocked]=useState(true),[message,setMessage]=useState("");
 const [confirmed,setConfirmed]=useState(false),[posting,setPosting]=useState(false),[result,setResult]=useState<Availability|null>(null);
 const fence=useRef(false),reading=useRef(false),operation=useRef<string|null>(null),form=useRef<HTMLFormElement>(null);
 const storageKey=`elite-availability-create:${scope}`;
 useEffect(()=>{
  try{
   if(!/^[a-f0-9]{64}$/.test(scope))throw new Error("scope");
   const raw=sessionStorage.getItem(storageKey);
   if(raw!==null){const marker=JSON.parse(raw);if(!marker||Object.keys(marker).join(",")!=="requestKey"||!/^availability-[a-f0-9-]{36}$/.test(marker.requestKey))throw new Error("marker");operation.current=marker.requestKey;setLocked(true);setMessage(t("p0459"));}
   else{operation.current=null;setLocked(false)}setReady(true);
  }catch{setReady(true);setLocked(true);setMessage(t("p0446"))}
 },[scope,storageKey]);
 function valid(value:unknown):value is Availability & {organization_id:string}{
  if(!value||typeof value!=="object")return false;const v=value as Record<string,unknown>;
  return typeof v.id==="string"&&v.id.length>0&&v.id.length<=128&&v.organization_id===organization&&["working","unavailable"].includes(String(v.entry_type))&&["active","cancelled"].includes(String(v.state))&&Number.isSafeInteger(v.version)&&Number(v.version)>0&&typeof v.starts_at==="string"&&typeof v.ends_at==="string"&&Number.isFinite(Date.parse(v.starts_at))&&Date.parse(v.ends_at)>Date.parse(v.starts_at)&&(v.resource_id===undefined||typeof v.resource_id==="string")&&(v.reason_code===undefined||typeof v.reason_code==="string");
 }
 async function submit(event:FormEvent<HTMLFormElement>){
  event.preventDefault();if(!ready||locked||fence.current)return;
  const data=new FormData(event.currentTarget),from=new Date(String(data.get("startsAt")??"")),to=new Date(String(data.get("endsAt")??""));
  if(!Number.isFinite(from.getTime())||!Number.isFinite(to.getTime())||to<=from){setMessage(t("p0460"));return}
  const resourceId=String(data.get("resourceId")??"").trim(),entryType=String(data.get("entryType")??"working"),reasonCode=String(data.get("reasonCode")??"").trim();
  fence.current=true;setLocked(true);setConfirmed(false);
  try{if(sessionStorage.getItem(storageKey)!==null)throw new Error("unresolved");const key="availability-"+crypto.randomUUID();sessionStorage.setItem(storageKey,JSON.stringify({requestKey:key}));operation.current=key;}
  catch{setMessage(t("p0447"));return}
  const key=operation.current;setPosting(true);setMessage(t("p0461"));
  try{
   const response=await fetch("/api/enterprise/franchise/commands",{method:"POST",headers:{"content-type":"application/json","idempotency-key":key!},signal:AbortSignal.timeout(10000),body:JSON.stringify({action:"create-availability",organizationId:organization,...(resourceId?{resourceId}:{}),entryType,...(reasonCode?{reasonCode}:{}),startsAt:from.toISOString(),endsAt:to.toISOString()})});
   const receipt=await response.json();if(operation.current!==key)return;
   if(!response.ok||!valid(receipt)||(receipt.resource_id??"")!==resourceId||receipt.entry_type!==entryType||(receipt.reason_code??"")!==reasonCode||Date.parse(receipt.starts_at)!==from.getTime()||Date.parse(receipt.ends_at)!==to.getTime()||receipt.state!=="active"||receipt.version!==1)throw new Error("unconfirmed");
   setResult(receipt);setMessage(t("p0462"));
  }catch{if(operation.current===key)setMessage(t("p0463"))}
  finally{setPosting(false)}
 }
 async function consult(){
  if(!ready||!operation.current||reading.current)return;reading.current=true;const key=operation.current;
  try{
   const query=new URLSearchParams({kind:"availability",organizationId:organization,requestKey:key});
   const response=await fetch("/api/enterprise/franchise/commands?"+query,{cache:"no-store",signal:AbortSignal.timeout(10000)}),receipt=await response.json();
   if(operation.current!==key)return;
   if(!response.ok||receipt?.request_key!==key||!valid(receipt.availability))throw new Error("unconfirmed");
   setResult(receipt.availability);setConfirmed(true);setMessage(t("p0464"));
  }catch{if(operation.current===key){setConfirmed(false);setMessage(t("p0465"))}}
  finally{reading.current=false}
 }
 function prepare(){
  if(!confirmed||posting||reading.current)return;
  try{sessionStorage.removeItem(storageKey);operation.current=null;fence.current=false;setLocked(false);setConfirmed(false);setResult(null);form.current?.reset();setMessage(t("p0466"))}
  catch{setMessage(t("p0467"))}
 }
 return <div>
  <form ref={form} onSubmit={event=>void submit(event)}>
   <label>{t("p0468")}<input name="resourceId" maxLength={128} placeholder={t("p0469")} disabled={!ready||locked}/></label>
   <label>{t("p0470")}<select name="entryType" defaultValue="working" disabled={!ready||locked}><option value="working">{t("p0471")}</option><option value="unavailable">{t("p0472")}</option></select></label>
   <label>{t("p0473")}<input name="reasonCode" pattern="[a-z][a-z0-9]*(-[a-z0-9]+)*" placeholder={t("p0474")} disabled={!ready||locked}/></label>
   <label>{t("p0475")}<input name="startsAt" type="datetime-local" required disabled={!ready||locked}/></label><label>{t("p0476")}<input name="endsAt" type="datetime-local" required disabled={!ready||locked}/></label>
   <button disabled={!ready||locked}>{t("p0477")}</button>
  </form>
  <p role="status" aria-label={t("p0478")} aria-live="polite">{message}</p>
  {result?<p>{t("p0479")} {result.id} {t("p0016")} {result.state} {t("p0016")} {new Date(result.starts_at).toLocaleString(privateLocale.locale, {timeZone:privateLocale.timeZone})} {t("p0435")} {new Date(result.ends_at).toLocaleString(privateLocale.locale, {timeZone:privateLocale.timeZone})}</p>:null}
  <button type="button" disabled={!ready||!operation.current} onClick={()=>void consult()}>{t("p0480")}</button>
  <button type="button" disabled={!confirmed||posting} onClick={prepare}>{t("p0481")}</button>
  <OperationalGuide guide={OPERATIONAL_GUIDES["availability-create-view"]}/>
 </div>;
}

type ResourceReceipt={id:string;organization_id:string;display_name:string;principal_subject?:string;kind:string;status:string;version:number;skills:string[]};
function ResourceCreation({organization,scope}:{organization:string;scope:string}) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

 const [ready,setReady]=useState(false),[locked,setLocked]=useState(true),[message,setMessage]=useState("");
 const [confirmed,setConfirmed]=useState(false),[posting,setPosting]=useState(false),[result,setResult]=useState<ResourceReceipt|null>(null);
 const fence=useRef(false),reading=useRef(false),operation=useRef<string|null>(null),form=useRef<HTMLFormElement>(null);
 const storageKey=`elite-resource-create:${scope}`;
 useEffect(()=>{
  try{
   if(!/^[a-f0-9]{64}$/.test(scope))throw new Error("scope");
   const raw=sessionStorage.getItem(storageKey);
   if(raw!==null){const marker=JSON.parse(raw);if(!marker||Object.keys(marker).join(",")!=="requestKey"||!/^resource-[a-f0-9-]{36}$/.test(marker.requestKey))throw new Error("marker");operation.current=marker.requestKey;setLocked(true);setMessage(t("p0482"));}
   else{operation.current=null;setLocked(false)}setReady(true);
  }catch{setReady(true);setLocked(true);setMessage(t("p0446"))}
 },[scope,storageKey]);
 function valid(value:unknown):value is ResourceReceipt {
  if(!value||typeof value!=="object")return false;const v=value as Record<string,unknown>;
  return typeof v.id==="string"&&v.id.length>0&&v.id.length<=128&&v.organization_id===organization&&typeof v.display_name==="string"&&v.display_name.length>0&&v.display_name.length<=160&&["employee","contractor","service-bay","vehicle"].includes(String(v.kind))&&["active","inactive"].includes(String(v.status))&&Number.isSafeInteger(v.version)&&Number(v.version)>0&&(v.principal_subject===undefined||typeof v.principal_subject==="string")&&Array.isArray(v.skills)&&v.skills.length>0&&v.skills.length<=4&&new Set(v.skills).size===v.skills.length&&v.skills.every(skill=>typeof skill==="string"&&["consultation","test-drive","delivery","service"].includes(skill));
 }
 async function submit(event:FormEvent<HTMLFormElement>){
  event.preventDefault();if(!ready||locked||fence.current)return;
  const data=new FormData(event.currentTarget),displayName=String(data.get("displayName")??"").trim(),principalSubject=String(data.get("principalSubject")??"").trim(),kind=String(data.get("kind")??"employee"),skills=data.getAll("skills").map(String);
  if(!displayName||!skills.length||((kind==="employee"||kind==="contractor")!==!!principalSubject)){setMessage(t("p0483"));return}
  fence.current=true;setLocked(true);setConfirmed(false);
  try{if(sessionStorage.getItem(storageKey)!==null)throw new Error("unresolved");const key="resource-"+crypto.randomUUID();sessionStorage.setItem(storageKey,JSON.stringify({requestKey:key}));operation.current=key;}
  catch{setMessage(t("p0447"));return}
  const key=operation.current;setPosting(true);setMessage(t("p0484"));
  try{
   const response=await fetch("/api/enterprise/franchise/commands",{method:"POST",headers:{"content-type":"application/json","idempotency-key":key!},signal:AbortSignal.timeout(10000),body:JSON.stringify({action:"create-resource",organizationId:organization,displayName,...(principalSubject?{principalSubject}:{}),kind,skills})});
   const receipt=await response.json();if(operation.current!==key)return;
   if(!response.ok||!valid(receipt)||receipt.display_name!==displayName||(receipt.principal_subject??"")!==principalSubject||receipt.kind!==kind||JSON.stringify([...receipt.skills].sort())!==JSON.stringify([...skills].sort())||receipt.status!=="active"||receipt.version!==1)throw new Error("unconfirmed");
   setResult(receipt);setMessage(t("p0485"));
  }catch{if(operation.current===key)setMessage(t("p0486"))}
  finally{setPosting(false)}
 }
 async function consult(){
  if(!ready||!operation.current||reading.current)return;reading.current=true;const key=operation.current;
  try{
   const query=new URLSearchParams({kind:"resource",organizationId:organization,requestKey:key});
   const response=await fetch("/api/enterprise/franchise/commands?"+query,{cache:"no-store",signal:AbortSignal.timeout(10000)}),receipt=await response.json();
   if(operation.current!==key)return;
   if(!response.ok||receipt?.request_key!==key||!valid(receipt.resource))throw new Error("unconfirmed");
   setResult(receipt.resource);setConfirmed(true);setMessage(t("p0487"));
  }catch{if(operation.current===key){setConfirmed(false);setMessage(t("p0488"))}}
  finally{reading.current=false}
 }
 function prepare(){
  if(!confirmed||posting||reading.current)return;
  try{sessionStorage.removeItem(storageKey);operation.current=null;fence.current=false;setLocked(false);setConfirmed(false);setResult(null);form.current?.reset();setMessage(t("p0489"))}
  catch{setMessage(t("p0467"))}
 }
 return <div>
  <form ref={form} onSubmit={event=>void submit(event)}>
   <label>{t("p0490")}<input name="displayName" required maxLength={160} disabled={!ready||locked}/></label>
   <label>{t("p0491")}<input name="principalSubject" maxLength={256} disabled={!ready||locked}/></label>
   <label>{t("p0470")}<select name="kind" defaultValue="employee" disabled={!ready||locked}><option value="employee">{t("p0492")}</option><option value="contractor">{t("p0493")}</option><option value="service-bay">{t("p0494")}</option><option value="vehicle">{t("p0495")}</option></select></label>
   <fieldset disabled={!ready||locked}><legend>{t("p0496")}</legend>{["consultation","test-drive","delivery","service"].map(skill=><label key={skill}><input name="skills" type="checkbox" value={skill}/>{skill}</label>)}</fieldset>
   <button disabled={!ready||locked}>{t("p0497")}</button>
  </form>
  <p role="status" aria-label={t("p0498")} aria-live="polite">{message}</p>
  {result?<p>{t("p0499")} {result.id} {t("p0016")} {result.display_name} {t("p0016")} {result.kind} {t("p0016")} {result.status}</p>:null}
  <button type="button" disabled={!ready||!operation.current} onClick={()=>void consult()}>{t("p0500")}</button>
  <button type="button" disabled={!confirmed||posting} onClick={prepare}>{t("p0501")}</button>
  <OperationalGuide guide={OPERATIONAL_GUIDES["resource-create-view"]}/>
 </div>;
}

type SlotReceipt={id:string;organization_id:string;kind:string;starts_at:string;ends_at:string;capacity:number;booked:number;state:string;version:number};
function SlotCreation({organization,scope}:{organization:string;scope:string}) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

 const [ready,setReady]=useState(false),[locked,setLocked]=useState(true),[message,setMessage]=useState("");
 const [confirmed,setConfirmed]=useState(false),[posting,setPosting]=useState(false),[result,setResult]=useState<SlotReceipt|null>(null);
 const fence=useRef(false),reading=useRef(false),operation=useRef<string|null>(null),form=useRef<HTMLFormElement>(null);
 const storageKey=`elite-slot-create:${scope}`;
 useEffect(()=>{
  try{
   if(!/^[a-f0-9]{64}$/.test(scope))throw new Error("scope");
   const raw=sessionStorage.getItem(storageKey);
   if(raw!==null){const marker=JSON.parse(raw);if(!marker||Object.keys(marker).join(",")!=="requestKey"||!/^slot-[a-f0-9-]{36}$/.test(marker.requestKey))throw new Error("marker");operation.current=marker.requestKey;setLocked(true);setMessage(t("p0502"));}
   else{operation.current=null;setLocked(false)}setReady(true);
  }catch{setReady(true);setLocked(true);setMessage(t("p0446"))}
 },[scope,storageKey]);
 function valid(value:unknown):value is SlotReceipt {
  if(!value||typeof value!=="object")return false;const v=value as Record<string,unknown>;
  return typeof v.id==="string"&&v.id.length>0&&v.id.length<=128&&v.organization_id===organization&&["consultation","test-drive","delivery","service"].includes(String(v.kind))&&["open","closed"].includes(String(v.state))&&Number.isSafeInteger(v.version)&&Number(v.version)>0&&typeof v.starts_at==="string"&&typeof v.ends_at==="string"&&Number.isFinite(Date.parse(v.starts_at))&&Date.parse(v.ends_at)>Date.parse(v.starts_at)&&Number.isSafeInteger(v.capacity)&&Number(v.capacity)>=1&&Number(v.capacity)<=100&&Number.isSafeInteger(v.booked)&&Number(v.booked)>=0&&Number(v.booked)<=Number(v.capacity);
 }
 async function submit(event:FormEvent<HTMLFormElement>){
  event.preventDefault();if(!ready||locked||fence.current)return;
  const data=new FormData(event.currentTarget),from=new Date(String(data.get("startsAt")??"")),to=new Date(String(data.get("endsAt")??""));
  if(!Number.isFinite(from.getTime())||!Number.isFinite(to.getTime())||to<=from){setMessage(t("p0503"));return}
  const kind=String(data.get("kind")??"test-drive"),capacity=Number(data.get("capacity"));
  if(!Number.isInteger(capacity)||capacity<1||capacity>100){setMessage(t("p0504"));return}
  fence.current=true;setLocked(true);setConfirmed(false);
  try{if(sessionStorage.getItem(storageKey)!==null)throw new Error("unresolved");const key="slot-"+crypto.randomUUID();sessionStorage.setItem(storageKey,JSON.stringify({requestKey:key}));operation.current=key;}
  catch{setMessage(t("p0447"));return}
  const key=operation.current;setPosting(true);setMessage(t("p0505"));
  try{
   const response=await fetch("/api/enterprise/franchise/commands",{method:"POST",headers:{"content-type":"application/json","idempotency-key":key!},signal:AbortSignal.timeout(10000),body:JSON.stringify({action:"create-appointment-slot",organizationId:organization,kind,capacity,startsAt:from.toISOString(),endsAt:to.toISOString()})});
   const receipt=await response.json();if(operation.current!==key)return;
   if(!response.ok||!valid(receipt)||receipt.kind!==kind||receipt.capacity!==capacity||Date.parse(receipt.starts_at)!==from.getTime()||Date.parse(receipt.ends_at)!==to.getTime()||receipt.state!=="open"||receipt.version!==1)throw new Error("unconfirmed");
   setResult(receipt);setMessage(t("p0506"));
  }catch{if(operation.current===key)setMessage(t("p0507"))}
  finally{setPosting(false)}
 }
 async function consult(){
  if(!ready||!operation.current||reading.current)return;reading.current=true;const key=operation.current;
  try{
   const query=new URLSearchParams({kind:"slot",organizationId:organization,requestKey:key});
   const response=await fetch("/api/enterprise/franchise/commands?"+query,{cache:"no-store",signal:AbortSignal.timeout(10000)}),receipt=await response.json();
   if(operation.current!==key)return;
   if(!response.ok||receipt?.request_key!==key||!valid(receipt.slot))throw new Error("unconfirmed");
   setResult(receipt.slot);setConfirmed(true);setMessage(t("p0508"));
  }catch{if(operation.current===key){setConfirmed(false);setMessage(t("p0509"))}}
  finally{reading.current=false}
 }
 function prepare(){
  if(!confirmed||posting||reading.current)return;
  try{sessionStorage.removeItem(storageKey);operation.current=null;fence.current=false;setLocked(false);setConfirmed(false);setResult(null);form.current?.reset();setMessage(t("p0510"))}
  catch{setMessage(t("p0467"))}
 }
 return <div>
  <form ref={form} onSubmit={event=>void submit(event)}>
   <label>{t("p0470")}<select name="kind" defaultValue="test-drive" disabled={!ready||locked}><option value="test-drive">{t("p0367")}</option><option value="consultation">{t("p0366")}</option><option value="delivery">{t("p0368")}</option><option value="service">{t("p0018")}</option></select></label>
   <label>{t("p0475")}<input name="startsAt" type="datetime-local" required disabled={!ready||locked}/></label><label>{t("p0476")}<input name="endsAt" type="datetime-local" required disabled={!ready||locked}/></label>
   <label>{t("p0511")}<input name="capacity" type="number" min={1} max={100} defaultValue={1} required disabled={!ready||locked}/></label>
   <button disabled={!ready||locked}>{t("p0512")}</button>
  </form>
  <p role="status" aria-label={t("p0513")} aria-live="polite">{message}</p>
  {result?<p>{t("p0140")} {result.id} {t("p0016")} {result.state} {t("p0016")} {result.booked}{t("p0280")}{result.capacity} {t("p0514")} {new Date(result.starts_at).toLocaleString(privateLocale.locale, {timeZone:privateLocale.timeZone})} {t("p0435")} {new Date(result.ends_at).toLocaleString(privateLocale.locale, {timeZone:privateLocale.timeZone})}</p>:null}
  <button type="button" disabled={!ready||!operation.current} onClick={()=>void consult()}>{t("p0515")}</button>
  <button type="button" disabled={!confirmed||posting} onClick={prepare}>{t("p0516")}</button>
  <OperationalGuide guide={OPERATIONAL_GUIDES["slot-create-view"]}/>
 </div>;
}



type PublishedChecklist={id:string;organization_id:string;version:number;title:string;state:string;items:{id:string;ordinal:number;prompt:string;response_type:string;required:boolean}[]};
type ChecklistMarker={checklistId:string;version:number;digest:string};
function ChecklistPublication({organization,scope}:{organization:string;scope:string}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

 const [ready,setReady]=useState(false),[locked,setLocked]=useState(true),[message,setMessage]=useState("");
 const [posting,setPosting]=useState(false),[querying,setQuerying]=useState(false),[confirmed,setConfirmed]=useState(false),[result,setResult]=useState<PublishedChecklist|null>(null);
 const [checklistItems,setChecklistItems]=useState<ChecklistDraftItem[]>([{key:1,id:"",prompt:"",responseType:"confirmation",required:true}]);
 const form=useRef<HTMLFormElement>(null),fence=useRef(false),reading=useRef(false),operation=useRef<ChecklistMarker|null>(null);
 const storageKey=`elite-checklist-publish:${scope}`;
 const code=(v:unknown):v is string=>typeof v==="string"&&v.length<=64&&/^[a-z][a-z0-9]*(-[a-z0-9]+)*$/.test(v);
 const textWithin=(v:unknown,max:number):v is string=>typeof v==="string"&&v.trim().length>0&&new TextEncoder().encode(v).length<=max;
 function valid(v:unknown):v is PublishedChecklist{
  if(!v||typeof v!=="object")return false;const value=v as PublishedChecklist;
  return code(value.id)&&value.organization_id===organization&&Number.isSafeInteger(value.version)&&value.version>0&&value.state==="published"&&textWithin(value.title,160)&&Array.isArray(value.items)&&value.items.length>0&&value.items.length<=64&&new Set(value.items.map(i=>i?.id)).size===value.items.length&&value.items.every((i,index)=>i&&code(i.id)&&i.ordinal===index+1&&textWithin(i.prompt,500)&&["confirmation","text","serial","evidence"].includes(i.response_type)&&typeof i.required==="boolean");
 }
 async function digest(value:PublishedChecklist){
  const bytes=new TextEncoder().encode(JSON.stringify({id:value.id,organization_id:value.organization_id,version:value.version,title:value.title,items:value.items.map(i=>({id:i.id,prompt:i.prompt,response_type:i.response_type,required:i.required}))}));
  return Array.from(new Uint8Array(await crypto.subtle.digest("SHA-256",bytes)),b=>b.toString(16).padStart(2,"0")).join("");
 }
 useEffect(()=>{
  try{
   if(!/^[a-f0-9]{64}$/.test(scope))throw new Error("scope");
   const raw=sessionStorage.getItem(storageKey);
   if(raw!==null){const v=JSON.parse(raw);if(!v||Object.keys(v).sort().join(",")!=="checklistId,digest,version"||!code(v.checklistId)||!Number.isSafeInteger(v.version)||v.version<1||!/^[a-f0-9]{64}$/.test(v.digest))throw new Error("marker");operation.current=v;setLocked(true);setMessage(t("p0517"));}
   else{operation.current=null;setLocked(false)}setReady(true);
  }catch{setReady(true);setLocked(true);setMessage(t("p0518"))}
 },[scope,storageKey]);
 async function submit(event:FormEvent<HTMLFormElement>){
  event.preventDefault();if(!ready||locked||fence.current)return;
  const data=new FormData(event.currentTarget);const value:PublishedChecklist={id:String(data.get("checklistId")??""),organization_id:organization,version:Number(data.get("checklistVersion")),title:String(data.get("title")??"").trim(),state:"published",items:checklistItems.map((i,index)=>({id:i.id,ordinal:index+1,prompt:i.prompt.trim(),response_type:i.responseType,required:i.required}))};
  if(!valid(value)){setMessage(t("p0519"));return}
  fence.current=true;setLocked(true);setConfirmed(false);setPosting(true);
  let marker:ChecklistMarker;
  try{if(sessionStorage.getItem(storageKey)!==null)throw new Error("unresolved");marker={checklistId:value.id,version:value.version,digest:await digest(value)};sessionStorage.setItem(storageKey,JSON.stringify(marker));operation.current=marker;}
  catch{setPosting(false);setMessage(t("p0520"));return}
  setMessage(t("p0521"));
  try{
   const response=await fetch("/api/enterprise/franchise/commands",{method:"POST",headers:{"content-type":"application/json"},signal:AbortSignal.timeout(10000),body:JSON.stringify({action:"publish-delivery-checklist",organizationId:organization,checklistId:value.id,version:value.version,title:value.title,items:value.items.map(i=>({id:i.id,prompt:i.prompt,response_type:i.response_type,required:i.required}))})});
   const receipt=await response.json();if(operation.current!==marker)return;
   if(!response.ok||!valid(receipt)||await digest(receipt)!==marker.digest)throw new Error("unconfirmed");
   setResult(receipt);setMessage(t("p0522"));
  }catch{if(operation.current===marker)setMessage(t("p0523"))}
  finally{setPosting(false)}
 }
 async function consult(){
  const marker=operation.current;if(!ready||!marker||reading.current)return;reading.current=true;setQuerying(true);
  try{
   const q=new URLSearchParams({kind:"checklist",organizationId:organization,checklistId:marker.checklistId,version:String(marker.version)});
   const response=await fetch("/api/enterprise/franchise/commands?"+q,{cache:"no-store",signal:AbortSignal.timeout(10000)}),body=await response.json();
   if(operation.current!==marker)return;
   if(!response.ok||!valid(body?.checklist)||body.checklist.id!==marker.checklistId||body.checklist.version!==marker.version||await digest(body.checklist)!==marker.digest)throw new Error("unconfirmed");
   setResult(body.checklist);setConfirmed(true);setMessage(t("p0524"));
  }catch{if(operation.current===marker){setConfirmed(false);setMessage(t("p0525"))}}
  finally{reading.current=false;setQuerying(false)}
 }
 function prepare(){
  if(!confirmed||posting||reading.current)return;
  try{sessionStorage.removeItem(storageKey);operation.current=null;fence.current=false;setLocked(false);setConfirmed(false);setResult(null);form.current?.reset();setChecklistItems([{key:1,id:"",prompt:"",responseType:"confirmation",required:true}]);setMessage(t("p0526"))}
  catch{setMessage(t("p0527"))}
 }
 return <div><h3>{t("p0528")}</h3><form ref={form} onSubmit={event=>void submit(event)}><fieldset disabled={!ready||locked}>
 
      <label>{t("p0529")}<input name="checklistId" required maxLength={64} pattern="[a-z][a-z0-9]*(-[a-z0-9]+)*" /></label><label>{t("p0530")}<input name="checklistVersion" type="number" min={1} required /></label><label>{t("p0531")}<input name="title" required maxLength={160} /></label>
      {checklistItems.map((item, index) => <fieldset key={item.key}><legend>{t("p0532")} {index + 1}</legend><label>{t("p0533")}<input required value={item.id} maxLength={64} pattern="[a-z][a-z0-9]*(-[a-z0-9]+)*" onChange={(event) => setChecklistItems((items) => items.map((candidate) => candidate.key === item.key ? { ...candidate, id: event.target.value } : candidate))}/></label><label>{t("p0534")}<input required value={item.prompt} maxLength={500} onChange={(event) => setChecklistItems((items) => items.map((candidate) => candidate.key === item.key ? { ...candidate, prompt: event.target.value } : candidate))}/></label><label>{t("p0535")}<select value={item.responseType} onChange={(event) => setChecklistItems((items) => items.map((candidate) => candidate.key === item.key ? { ...candidate, responseType: event.target.value as ChecklistDraftItem["responseType"] } : candidate))}><option value="confirmation">{t("p0141")}</option><option value="text">{t("p0536")}</option><option value="serial">{t("p0537")}</option></select></label><label><input type="checkbox" checked={item.required} onChange={(event) => setChecklistItems((items) => items.map((candidate) => candidate.key === item.key ? { ...candidate, required: event.target.checked } : candidate))}/> {t("p0538")}</label>{checklistItems.length > 1 ? <button type="button" onClick={() => setChecklistItems((items) => items.filter((candidate) => candidate.key !== item.key))}>{t("p0539")}</button> : null}</fieldset>)}
      <button type="button" disabled={checklistItems.length >= 64} onClick={() => setChecklistItems((items) => [...items, { key: Math.max(...items.map((item) => item.key)) + 1, id: "", prompt: "", responseType: "confirmation", required: true }])}>{t("p0540")}</button><button disabled={!ready||locked}>{t("p0541")}</button>
    
 </fieldset></form>
 <p role="status" aria-label={t("p0542")} aria-live="polite">{message}</p>
 {result&&result.organization_id===organization?<div><p>{result.id} {t("p0543")} {result.version} {t("p0016")} {result.title} {t("p0016")} {result.state}</p><ol>{result.items.map(item=><li key={item.id}>{item.id}{t("p0213")} {item.prompt} {t("p0016")} {item.response_type} {t("p0016")} {item.required?t("p0538"):t("p0544")}</li>)}</ol></div>:null}
 <button type="button" disabled={!ready||!operation.current||querying} onClick={()=>void consult()}>{t("p0545")}</button>
 <button type="button" disabled={!confirmed||posting||querying} onClick={prepare}>{t("p0546")}</button>
 <OperationalGuide guide={OPERATIONAL_GUIDES["checklist-publication-view"]}/>
 </div>;
}


type CompletionReference={handoverId:string;version:number;checklistId:string;checklistVersion:number;digest:string};
type CompletionResponse={item_id:string;response_text:string;evidence_sha256?:string};
type CompletionResult={handover_id:string;organization_id:string;state:string;version:number;checklist_id:string;checklist_version:number;completed_at:string;actor_subject:string;responses:CompletionResponse[]};
export function ChecklistCompletionPanel({organization,scope,initialHandover}:{organization:string;scope:string;initialHandover?:{id:string;version:number}}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

 const [ready,setReady]=useState(false),[locked,setLocked]=useState(true),[message,setMessage]=useState("");
 const [posting,setPosting]=useState(false),[querying,setQuerying]=useState(false),[confirmed,setConfirmed]=useState(false),[result,setResult]=useState<CompletionResult|null>(null);
 const form=useRef<HTMLFormElement>(null),fence=useRef(false),reading=useRef(false),operation=useRef<CompletionReference|null>(null);
 const storageKey=`elite-checklist-complete:${scope}${initialHandover ? ":"+initialHandover.id : ""}`;
 const id=(v:unknown):v is string=>typeof v==="string"&&/^[A-Za-z0-9_-]{1,128}$/.test(v);
 const code=(v:unknown):v is string=>typeof v==="string"&&v.length<=64&&/^[a-z][a-z0-9]*(-[a-z0-9]+)*$/.test(v);
 const version=(v:unknown):v is number=>typeof v==="number"&&Number.isSafeInteger(v)&&v>=1&&v<Number.MAX_SAFE_INTEGER;
 const bounded=(v:unknown,max:number):v is string=>typeof v==="string"&&v.trim().length>0&&new TextEncoder().encode(v).length<=max;
 function validResponses(value:unknown):value is CompletionResponse[]{return Array.isArray(value)&&value.length>=1&&value.length<=64&&new Set(value.map(r=>r?.item_id)).size===value.length&&value.every(r=>r&&code(r.item_id)&&bounded(r.response_text,2048)&&(r.evidence_sha256===undefined||r.evidence_sha256===""||typeof r.evidence_sha256==="string"&&/^[a-f0-9]{64}$/.test(r.evidence_sha256)))}
 function valid(v:unknown):v is CompletionResult{if(!v||typeof v!=="object")return false;const value=v as CompletionResult;return id(value.handover_id)&&value.organization_id===organization&&version(value.version)&&value.version>=2&&["presented","accepted","rejected"].includes(value.state)&&code(value.checklist_id)&&version(value.checklist_version)&&typeof value.completed_at==="string"&&Number.isFinite(Date.parse(value.completed_at))&&bounded(value.actor_subject,255)&&validResponses(value.responses)}
 async function digest(value:Pick<CompletionResult,"handover_id"|"organization_id"|"checklist_id"|"checklist_version"|"responses">){
  const responses=[...value.responses].sort((a,b)=>a.item_id<b.item_id?-1:a.item_id>b.item_id?1:0).map(r=>({item_id:r.item_id,response_text:r.response_text,evidence_sha256:r.evidence_sha256??""}));
  const bytes=new TextEncoder().encode(JSON.stringify({handover_id:value.handover_id,organization_id:value.organization_id,checklist_id:value.checklist_id,checklist_version:value.checklist_version,responses}));
  return Array.from(new Uint8Array(await crypto.subtle.digest("SHA-256",bytes)),b=>b.toString(16).padStart(2,"0")).join("");
 }
 useEffect(()=>{try{
  if(!/^[a-f0-9]{64}$/.test(scope))throw new Error("scope");const raw=sessionStorage.getItem(storageKey);
  if(raw!==null){const v=JSON.parse(raw);if(!v||Object.keys(v).sort().join(",")!=="checklistId,checklistVersion,digest,handoverId,version"||!id(v.handoverId)||!version(v.version)||!code(v.checklistId)||!version(v.checklistVersion)||!/^[a-f0-9]{64}$/.test(v.digest))throw new Error("reference");operation.current=v;setLocked(true);setMessage(t("p0547"))}
  else{operation.current=null;setLocked(false)}setReady(true);
 }catch{setReady(true);setLocked(true);setMessage(t("p0518"))}},[scope,storageKey]);
 async function submit(event:FormEvent<HTMLFormElement>){
  event.preventDefault();if(!ready||locked||fence.current)return;const data=new FormData(event.currentTarget);
  const lines=String(data.get("responses")??"").split(/\r?\n/).map(line=>line.trim()).filter(Boolean);
  const responses=lines.map(line=>{const separator=line.indexOf("=");return {item_id:separator>0?line.slice(0,separator).trim():"",response_text:separator>0?line.slice(separator+1).trim():""}});
  const value={handover_id:String(data.get("handoverId")??""),organization_id:organization,checklist_id:String(data.get("checklistId")??""),checklist_version:Number(data.get("checklistVersion")),responses};const currentVersion=Number(data.get("handoverVersion"));
  if(!id(value.handover_id)||!code(value.checklist_id)||!version(value.checklist_version)||!version(currentVersion)||!validResponses(responses)){setMessage(t("p0548"));return}
  fence.current=true;setLocked(true);setPosting(true);setConfirmed(false);let reference:CompletionReference;
  try{if(sessionStorage.getItem(storageKey)!==null)throw new Error("unresolved");reference={handoverId:value.handover_id,version:currentVersion,checklistId:value.checklist_id,checklistVersion:value.checklist_version,digest:await digest(value)};sessionStorage.setItem(storageKey,JSON.stringify(reference));operation.current=reference}
  catch{setPosting(false);setMessage(t("p0520"));return}
  setMessage(t("p0549"));
  try{
   const response=await fetch("/api/enterprise/franchise/commands",{method:"POST",headers:{"content-type":"application/json"},signal:AbortSignal.timeout(10000),body:JSON.stringify({action:"complete-delivery-checklist",organizationId:organization,handoverId:value.handover_id,version:currentVersion,checklistId:value.checklist_id,checklistVersion:value.checklist_version,responses})});
   const receipt=await response.json();if(operation.current!==reference)return;
   if(!response.ok||receipt?.id!==value.handover_id||receipt.organization_id!==organization||receipt.state!=="presented"||receipt.version!==currentVersion+1||receipt.checklist_id!==value.checklist_id||receipt.checklist_version!==value.checklist_version||typeof receipt.checklist_completed_at!=="string"||!Number.isFinite(Date.parse(receipt.checklist_completed_at)))throw new Error("unconfirmed");
   setMessage(t("p0550"));
  }catch{if(operation.current===reference)setMessage(t("p0551"))}
  finally{setPosting(false)}
 }
 async function consult(){
  const reference=operation.current;if(!ready||!reference||reading.current)return;reading.current=true;setQuerying(true);
  try{
   const response=await fetch("/api/enterprise/franchise/commands?"+new URLSearchParams({kind:"checklist-completion",organizationId:organization,handoverId:reference.handoverId}),{cache:"no-store",signal:AbortSignal.timeout(10000)}),body=await response.json();if(operation.current!==reference)return;
   if(!response.ok||!valid(body?.completion)||body.completion.handover_id!==reference.handoverId||body.completion.version<reference.version+1||body.completion.checklist_id!==reference.checklistId||body.completion.checklist_version!==reference.checklistVersion||await digest(body.completion)!==reference.digest)throw new Error("unconfirmed");
   setResult(body.completion);setConfirmed(true);setMessage(t("p0552"));
  }catch{if(operation.current===reference){setConfirmed(false);setMessage(t("p0553"))}}
  finally{reading.current=false;setQuerying(false)}
 }
 function prepare(){if(!confirmed||posting||reading.current)return;try{sessionStorage.removeItem(storageKey);operation.current=null;fence.current=false;setLocked(false);setConfirmed(false);setResult(null);form.current?.reset();setMessage(t("p0554"))}catch{setMessage(t("p0527"))}}
 return <div role="region" aria-label={t("p0555")}><form ref={form} onSubmit={event=>void submit(event)}><fieldset disabled={!ready||locked}><h3>{t("p0556")}</h3><label>{t("p0368")}<input name="handoverId" required maxLength={128} defaultValue={initialHandover?.id} readOnly={Boolean(initialHandover)}/></label><label>{t("p0557")}<input name="handoverVersion" type="number" min={1} required defaultValue={initialHandover?.version} readOnly={Boolean(initialHandover)}/></label><label>{t("p0529")}<input name="checklistId" required maxLength={64} pattern="[a-z][a-z0-9]*(-[a-z0-9]+)*"/></label><label>{t("p0558")}<input name="checklistVersion" type="number" min={1} required/></label><label>{t("p0559")}<textarea name="responses" required placeholder={"serial-observed=SERIAL-123\nasset-condition=confirmed"}/></label><button disabled={!ready||locked}>{t("p0561")}</button></fieldset></form>
 <p role="status" aria-label={t("p0562")} aria-live="polite">{message}</p>
 {result&&result.organization_id===organization?<div><p>{result.handover_id} {t("p0563")} {result.state} {t("p0543")} {result.version}</p><p>{t("p0564")} {result.checklist_id} {t("p0543")} {result.checklist_version} {t("p0565")} {result.actor_subject} {t("p0016")} {result.completed_at}</p><ol>{result.responses.map(r=><li key={r.item_id}>{r.item_id}{t("p0213")} {r.response_text}{r.evidence_sha256?` · evidencia ${r.evidence_sha256}`:""}</li>)}</ol></div>:null}
 <button type="button" disabled={!ready||!operation.current||querying} onClick={()=>void consult()}>{t("p0566")}</button>
 <button type="button" disabled={!confirmed||posting||querying} onClick={prepare}>{t("p0567")}</button>
 <OperationalGuide guide={OPERATIONAL_GUIDES["checklist-completion-view"]}/>
 </div>;
}


type ReturnReference={action:"receive-return"|"decide-return";authorizationId:string;receiptId:string;digest:string};
type ReturnCommand={action:"receive-return";organizationId:string;authorizationId:string;serialNumber:string;conditionCode:string;notes:string}|{action:"decide-return";organizationId:string;receiptId:string;inventoryAction:string;notes:string};
type ReceivedReturn={id:string;authorization_id:string;organization_id:string;order_id:string;stock_unit_id:string;customer_subject:string;received_serial_number:string;condition_code:string;notes:string;evidence_sha256:string;received_by_subject:string;received_at:string};
type DecidedReturn={id:string;receipt_id:string;inventory_action:string;customer_remedy:string;notes:string;decided_by_subject:string;decided_at:string;effect_requests:(ReturnEffect&{idempotency_key:string;requested_at:string})[]};
type RecoveredReturnCase=Omit<ReturnCase,"receipt"|"disposition">&{organization_id:string;authorized_at:string;receipt?:ReceivedReturn;disposition?:DecidedReturn};
function ReturnOperationsPanel({organization,scope,cases}:{organization:string;scope:string;cases:ReturnCase[]|null}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

 const router=useRouter();const [caseList,setCaseList]=useState<ReturnCase[]|null>(cases),[ready,setReady]=useState(false),[locked,setLocked]=useState(true),[message,setMessage]=useState("");
 const [posting,setPosting]=useState(false),[querying,setQuerying]=useState(false),[confirmed,setConfirmed]=useState(false),[result,setResult]=useState<RecoveredReturnCase|null>(null);
 const fence=useRef(false),reading=useRef(false),operation=useRef<ReturnReference|null>(null);const storageKey=`elite-return-operation:${scope}`;
 const text=(v:unknown,max:number):v is string=>typeof v==="string"&&new TextEncoder().encode(v).length>=1&&new TextEncoder().encode(v).length<=max;
 const bounded=(v:unknown,max:number):v is string=>text(v,max)&&v.trim().length>0;
 const date=(v:unknown):v is string=>typeof v==="string"&&Number.isFinite(Date.parse(v));
 const hash=(v:unknown):v is string=>typeof v==="string"&&/^[a-f0-9]{64}$/.test(v);
 const effects:Record<string,string>={inventory:t("p0568"),refund:t("p0569"),exchange:t("p0570"),accounting:t("p0571"),fiscal:t("p0572")};
 function validCase(value:unknown):value is RecoveredReturnCase{
  if(!value||typeof value!=="object")return false;const v=value as RecoveredReturnCase;
  if(!text(v.authorization_id,128)||v.organization_id!==organization||!bounded(v.order_id,128)||!bounded(v.stock_unit_id,128)||!bounded(v.customer_subject,255)||!["return","exchange"].includes(v.authorized_action)||!date(v.authorized_at))return false;
  if(v.receipt===undefined)return v.disposition===undefined;const r=v.receipt;
  if(!r||!text(r.id,128)||r.authorization_id!==v.authorization_id||r.organization_id!==organization||r.order_id!==v.order_id||r.stock_unit_id!==v.stock_unit_id||r.customer_subject!==v.customer_subject||!text(r.received_serial_number,128)||!["sealed","opened","damaged","incomplete"].includes(r.condition_code)||!bounded(r.notes,1000)||!hash(r.evidence_sha256)||!bounded(r.received_by_subject,255)||!date(r.received_at))return false;
  if(v.disposition===undefined)return true;const d=v.disposition,remedy=v.authorized_action==="return"?"refund":"exchange";
  if(!d||!text(d.id,128)||d.receipt_id!==r.id||!["quarantine","restock","repair","scrap"].includes(d.inventory_action)||d.customer_remedy!==remedy||!bounded(d.notes,1000)||!bounded(d.decided_by_subject,255)||!date(d.decided_at)||!Array.isArray(d.effect_requests))return false;
  const owners:Record<string,string>={inventory:"inventory",accounting:"accounting",[remedy]:remedy==="refund"?"payment":"fulfillment"};if(remedy==="refund")owners.fiscal="fiscal";
  return d.effect_requests.length===Object.keys(owners).length&&new Set(d.effect_requests.map(e=>e?.id)).size===d.effect_requests.length&&new Set(d.effect_requests.map(e=>e?.effect_kind)).size===d.effect_requests.length&&d.effect_requests.every(e=>e&&text(e.id,128)&&Object.hasOwn(owners,e.effect_kind)&&owners[e.effect_kind]===e.owner_context&&e.state==="requested"&&text(e.idempotency_key,128)&&e.idempotency_key.length>=16&&date(e.requested_at));
 }
 async function digest(command:ReturnCommand){const bytes=new TextEncoder().encode(JSON.stringify(command));return Array.from(new Uint8Array(await crypto.subtle.digest("SHA-256",bytes)),b=>b.toString(16).padStart(2,"0")).join("")}
 async function matches(value:RecoveredReturnCase,reference:ReturnReference){
  if(value.authorization_id!==reference.authorizationId||!value.receipt)return false;
  if(reference.action==="receive-return"){
   const r=value.receipt;const command:ReturnCommand={action:"receive-return",organizationId:organization,authorizationId:reference.authorizationId,serialNumber:r.received_serial_number,conditionCode:r.condition_code,notes:r.notes};
   return r.evidence_sha256===reference.digest&&await digest(command)===reference.digest;
  }
  if(value.receipt.id!==reference.receiptId||!value.disposition)return false;const d=value.disposition;
  return await digest({action:"decide-return",organizationId:organization,receiptId:reference.receiptId,inventoryAction:d.inventory_action,notes:d.notes})===reference.digest;
 }
 useEffect(()=>setCaseList(cases),[cases]);
 useEffect(()=>{try{
  if(!/^[a-f0-9]{64}$/.test(scope))throw new Error("scope");const raw=sessionStorage.getItem(storageKey);
  if(raw!==null){if(raw.length>2048)throw new Error("reference");const v=JSON.parse(raw);if(!v||Object.keys(v).sort().join(",")!=="action,authorizationId,digest,receiptId"||!["receive-return","decide-return"].includes(v.action)||!text(v.authorizationId,128)||!hash(v.digest)||(v.action==="receive-return"?v.receiptId!=="":!text(v.receiptId,128)))throw new Error("reference");operation.current=v;setLocked(true);setMessage(t("p0573"))}
  else{operation.current=null;setLocked(false)}setReady(true);
 }catch{setReady(true);setLocked(true);setMessage(t("p0518"))}},[scope,storageKey]);
 async function submit(event:FormEvent<HTMLFormElement>,item:ReturnCase,action:ReturnReference["action"]){
  event.preventDefault();if(!ready||locked||fence.current)return;const data=new FormData(event.currentTarget);
  const command:ReturnCommand=action==="receive-return"?{action:"receive-return",organizationId:organization,authorizationId:item.authorization_id,serialNumber:String(data.get("serialNumber")??""),conditionCode:String(data.get("conditionCode")??""),notes:String(data.get("notes")??"").trim()}:{action:"decide-return",organizationId:organization,receiptId:item.receipt?.id??"",inventoryAction:String(data.get("inventoryAction")??""),notes:String(data.get("notes")??"").trim()};
  if(!text(item.authorization_id,128)||!bounded(command.notes,1000)||(command.action==="receive-return"?(!text(command.serialNumber,128)||!["sealed","opened","damaged","incomplete"].includes(command.conditionCode)):(!text(command.receiptId,128)||!["quarantine","restock","repair","scrap"].includes(command.inventoryAction)))){setMessage(t("p0574"));return}
  fence.current=true;setLocked(true);setPosting(true);setConfirmed(false);setResult(null);let reference:ReturnReference;
  try{if(sessionStorage.getItem(storageKey)!==null)throw new Error("unresolved");reference={action,authorizationId:item.authorization_id,receiptId:command.action==="decide-return"?command.receiptId:"",digest:await digest(command)};sessionStorage.setItem(storageKey,JSON.stringify(reference));operation.current=reference}
  catch{setPosting(false);setMessage(t("p0520"));return}
  setMessage(action==="receive-return"?t("p0575"):t("p0576"));
  try{
   const response=await fetch("/api/enterprise/franchise/commands",{method:"POST",headers:{"content-type":"application/json"},signal:AbortSignal.timeout(10000),body:JSON.stringify(command)}),body=await response.json();if(operation.current!==reference)return;
   if(!response.ok||!text(body?.id,128)||(action==="receive-return"?(body.authorization_id!==reference.authorizationId||body.organization_id!==organization||body.evidence_sha256!==reference.digest):body.receipt_id!==reference.receiptId))throw new Error("unconfirmed");
   setMessage(t("p0577"));
  }catch{if(operation.current===reference)setMessage(t("p0578"))}
  finally{setPosting(false)}
 }
 async function consult(){
  const reference=operation.current;if(!ready||!reference||reading.current)return;reading.current=true;setQuerying(true);
  try{
   const response=await fetch("/api/enterprise/franchise/commands?"+new URLSearchParams({kind:"return",organizationId:organization,authorizationId:reference.authorizationId}),{cache:"no-store",signal:AbortSignal.timeout(10000)}),body=await response.json();if(operation.current!==reference)return;
   if(!response.ok||!validCase(body?.returnCase)||!await matches(body.returnCase,reference))throw new Error("unconfirmed");
   setResult(body.returnCase);setConfirmed(true);setMessage(t("p0579"));
  }catch{if(operation.current===reference){setConfirmed(false);setMessage(t("p0580"))}}
  finally{reading.current=false;setQuerying(false)}
 }
 function prepare(){
  if(!confirmed||!result||posting||reading.current)return;
  try{sessionStorage.removeItem(storageKey);const recovered=result;setCaseList(previous=>previous?.some(item=>item.authorization_id===recovered.authorization_id)?previous.map(item=>item.authorization_id===recovered.authorization_id?recovered:item):[recovered,...(previous??[])]);operation.current=null;fence.current=false;setLocked(false);setConfirmed(false);setResult(null);setMessage(t("p0581"))}
  catch{setMessage(t("p0527"))}
 }
 return <section className="card" aria-label={t("p0582")}><h2>{t("p0583")}</h2><p>{t("p0584")}</p>
 <p role="status" aria-label={t("p0585")} aria-live="polite">{message}</p>
 <button type="button" disabled={!ready||!operation.current||querying} onClick={()=>void consult()}>{t("p0586")}</button>
 <button type="button" disabled={!confirmed||posting||querying} onClick={prepare}>{t("p0587")}</button>
 <button type="button" disabled={posting||querying} onClick={()=>router.refresh()}>{t("p0588")}</button>
 {result?<div role="region" aria-label={t("p0589")}><h3>{t("p0590")} {result.authorization_id}</h3>{result.receipt?<><p>{t("p0591")} {result.receipt.id} {t("p0592")} {result.receipt.received_by_subject} {t("p0016")} {result.receipt.received_at}</p><p>{t("p0593")} {result.receipt.received_serial_number} {t("p0016")} {result.receipt.condition_code}</p><p>{result.receipt.notes}</p></>:null}{result.disposition?<><p>{t("p0594")} {result.disposition.id} {t("p0595")} {result.disposition.decided_by_subject} {t("p0016")} {result.disposition.decided_at}</p><p>{result.disposition.inventory_action} {t("p0016")} {result.disposition.customer_remedy} {t("p0016")} {result.disposition.notes}</p><ul>{result.disposition.effect_requests.map(e=><li key={e.id}>{effects[e.effect_kind]}{t("p0596")} {e.id}</li>)}</ul></>:null}</div>:null}
 {cases===null?<p role="alert">{t("p0597")}</p>:null}
 {caseList?.length?<div className="grid">{caseList.map(item=><article className="card" key={item.authorization_id} aria-label={`Caso ${item.authorization_id}`}><h3>{item.authorized_action==="return"?t("p0598"):t("p0570")} {t("p0016")} {item.authorization_id}</h3><p>{t("p0251")} {item.order_id} {t("p0599")} {item.stock_unit_id}</p>
 {!item.receipt?<form onSubmit={event=>void submit(event,item,"receive-return")}><fieldset disabled={!ready||locked}><label>{t("p0600")}<input name="serialNumber" required maxLength={128}/></label><label>{t("p0601")}<select name="conditionCode" required><option value="sealed">{t("p0602")}</option><option value="opened">{t("p0603")}</option><option value="damaged">{t("p0604")}</option><option value="incomplete">{t("p0605")}</option></select></label><label>{t("p0606")}<textarea name="notes" required maxLength={1000}/></label><button disabled={!ready||locked}>{t("p0261")}</button></fieldset></form>:!item.disposition?<form onSubmit={event=>void submit(event,item,"decide-return")}><fieldset disabled={!ready||locked}><p>{t("p0607")} {item.receipt.received_serial_number} {t("p0016")} {item.receipt.condition_code}</p><label>{t("p0608")}<select name="inventoryAction" required><option value="quarantine">{t("p0609")}</option><option value="restock">{t("p0610")}</option><option value="repair">{t("p0611")}</option><option value="scrap">{t("p0612")}</option></select></label><label>{t("p0613")}<textarea name="notes" required maxLength={1000}/></label><button disabled={!ready||locked}>{t("p0614")}</button></fieldset></form>:<><p>{t("p0591")} {item.receipt.id} {t("p0615")} {item.disposition.id}</p><p>{item.disposition.inventory_action} {t("p0016")} {item.disposition.customer_remedy}</p><ul>{item.disposition.effect_requests.map(e=><li key={e.id}>{effects[e.effect_kind]}{t("p0596")} {e.id}</li>)}</ul></>}
 </article>)}</div>:cases!==null?<p>{t("p0616")}</p>:null}
 <OperationalGuide guide={OPERATIONAL_GUIDES["return-operations-view"]}/>
 </section>;
}
````

### FILE: `src/components/customer-quote-actions.tsx`

```yaml
block_id: "TS-FRANCHISE-JOURNEY:file:09"
operation: CREATE
provenance: AUTHORED
source: "local customer command UI governed by pack references"
license: "LicenseRef-Workspace-Owner"
sha256: "842e6b7d2f3f7bcf9b23cf74a6b148fdd89f002355413a7d0adab4308ba43a2c"
variables: []
secrets_allowed: false
```

````tsx
"use client";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

import { QUOTE_GUIDE } from "@/platform/help/content";

import { useState } from "react";

type Quote = { id: string; currency: string; total_minor_units: number; valid_until: string; state: string; version: number; order_id?: string; amountLabel: string; validUntilLabel: string; statusLabel: string; acceptanceAvailable: boolean };

export function CustomerQuoteActions({ organizationId, quotes }: { organizationId: string; quotes: Quote[] }) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

  const [pending, setPending] = useState("");
  const [message, setMessage] = useState("");
  const [uncertain, setUncertain] = useState(false);

  async function accept(quote: Quote) {
    if (pending || uncertain || !quote.acceptanceAvailable) return;
    if (Date.parse(quote.valid_until) <= Date.now()) {
      setUncertain(true);
      setMessage(t("p0269"));
      return;
    }
    setPending(quote.id);
    setMessage(t("p0270"));
    const controller = new AbortController();
    const timer = setTimeout(() => controller.abort(), 10000);
    try {
      const response = await fetch("/api/enterprise/franchise/commands", {
        method: "POST",
        headers: { "content-type": "application/json" },
        signal: controller.signal,
        body: JSON.stringify({ action: "accept-quote", organizationId, quoteId: quote.id, version: quote.version })
      });
      const result: unknown = await response.json();
      if (!response.ok || typeof result !== "object" || result === null ||
          !("state" in result) || result.state !== "accepted" ||
          !("order_id" in result) || typeof result.order_id !== "string" || result.order_id.length === 0) {
        throw new Error("RESULT_NOT_CONFIRMED");
      }
      setMessage(t("p0271"));
      window.location.reload();
    } catch {
      setUncertain(true);
      setMessage(t("p0272"));
      // Keep acceptance disabled. GET recovery, never an automatic POST retry.
    } finally {
      clearTimeout(timer);
    }
  }

  return <>
    {message && <p className="status" role="status" aria-live="polite">{message}</p>}
    {uncertain && <button className="button" onClick={() => window.location.reload()}>{t("p0273")}</button>}
    {quotes.length === 0 && <p className="notice">{t("p0274")}</p>}
    <div className="grid">{quotes.map((quote) => <article className="card" key={quote.id}>
      <h2>{quote.id}</h2>
      <p>{t("p0275")} {quote.amountLabel}</p>
      <p>{t("p0247")} {quote.statusLabel}</p>
      <p>{t("p0276")} {quote.validUntilLabel}</p>
      {quote.state === "accepted" && quote.order_id ? <p role="note">{t("p0277")} <strong>{quote.order_id}</strong></p> : null}
      {quote.acceptanceAvailable ? <button className="button" disabled={pending !== "" || uncertain} onClick={() => void accept(quote)}>{t("p0278")}</button> : null}
    </article>)}</div>
    <QuoteAcceptanceHelp/>
  </>;
}

export function QuoteAcceptanceHelp() {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

  return <details className="card">
    <summary>{t("p0279")}</summary>
    <p>{t("p0142")} {QUOTE_GUIDE.id}{t("p0280")}{QUOTE_GUIDE.version}</p>
    {QUOTE_GUIDE.paragraphs.map(text => <p key={text}>{controlled(text)}</p>)}
    <a href={`/help?article=${QUOTE_GUIDE.id}&version=${QUOTE_GUIDE.version}`}>{t("p0143")}</a>
  </details>;
}
````

### FILE: `src/app/customer/quotes/page.tsx`

```yaml
block_id: "TS-FRANCHISE-JOURNEY:file:10"
operation: CREATE
provenance: AUTHORED
source: "local customer quote page over canonical session and protected client"
license: "LicenseRef-Workspace-Owner"
sha256: "18d3b3c0bc7c2bc08376157884057c83b6f16f3c47cb8d7a5fd6a5a5d25286b0"
variables: []
secrets_allowed: false
```

````tsx
import {privateTranslator} from "@/platform/i18n/private-catalog";
import type {PrivateLocale} from "@/platform/i18n/private-locale";
import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import { minorAmountPresentation } from "@/platform/i18n/money";
import { redirect } from "next/navigation";
import type { Route } from "next";
import { allowed, readSession } from "@/platform/auth/session";
import { protectedGet, type CustomerJourney } from "@/platform/backend/protected-client";
import { CustomerQuoteActions, QuoteAcceptanceHelp } from "@/components/customer-quote-actions";
import Link from "next/link";

// AUTHORED presentation only. ISO currency digits from the pinned runtime's Intl;
// no pricing, exchange, tax or rounding rule is introduced. Compute on server
// once so browser ICU/timezone differences do not cause hydration mismatches.
function quotePresentation(quote: CustomerJourney["quotes"][number], asOf: number, locale:PrivateLocale) {
 const t=privateTranslator(locale.language);

  const { amountLabel, amountValid } = minorAmountPresentation(quote.total_minor_units, quote.currency, locale.locale);
  const until = Date.parse(quote.valid_until);
  const dateValid = Number.isFinite(until);
  const expired = dateValid && until <= asOf;
  const states: Record<string, string> = { issued: expired ? t("p0053") : t("p0054"), accepted: t("p0055"), cancelled: t("p0056"), draft: t("p0057"), expired: t("p0053") };
  return { ...quote, amountLabel, statusLabel: states[quote.state] ?? t("p0058"),
    validUntilLabel: dateValid ? new Intl.DateTimeFormat(locale.locale, { dateStyle: "medium", timeStyle: "short", timeZone: locale.timeZone }).format(until) + " (" + locale.timeZone + ")" : t("p0059"),
    acceptanceAvailable: quote.state === "issued" && amountValid && dateValid && !expired };
}

export default async function CustomerQuotesPage() {
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

  const session = await readSession();
  if (!session) redirect("/api/auth/login?return_to=/customer" as Route);
  if (!allowed(session, "customer:self")) return <><h1 className="pageTitle">{t("p0006")}</h1><div className="notice">{t("p0007")} <code>customer:self</code>{t("p0008")}</div></>;
  const organization = session.organizations[0]!;
  let journey: CustomerJourney;
  try {
    journey = await protectedGet<CustomerJourney>(session, "/v1/customer/journey", { organization_id: organization });
    if (!Array.isArray(journey.quotes)) throw new Error("INVALID_JOURNEY");
  } catch {
    return <>
      <h1 className="pageTitle">{t("p0052")}</h1>
      <p role="alert">{t("p0060")}</p>
      <p><a className="button" href="/customer/quotes">{t("p0033")}</a></p>
      <QuoteAcceptanceHelp/>
    </>;
  }
  return <>
    <div className="eyebrow">{t("p0034")} {organization}</div>
    <h1 className="pageTitle">{t("p0052")}</h1>
    <p className="lede">{t("p0061")}</p>
    <p><Link href={"/customer/appointments" as Route}>{t("p0062")}</Link></p>
    <CustomerQuoteActions organizationId={organization} quotes={journey.quotes.map(quote => quotePresentation(quote, Date.now(), privateLocale))}/>
  </>;
}
````

### FILE: `src/app/api/enterprise/appointments/route.ts`

```yaml
block_id: "TS-FRANCHISE-JOURNEY:file:01"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by the upstream references in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "78e6effcc7de35e9ae3fcb00599666609c4429baa0fe95a75664765eef00f120"
variables: []
secrets_allowed: false
```

````ts
import { z } from "zod";
import { requestAppointment } from "@/platform/backend/public-client";
import { errorResponse, problem } from "@/platform/http/problem";

const schema = z.object({
  leadId: z.string().min(1).max(128),
  modelId: z.string().min(1).max(128).optional(),
  kind: z.enum(["consultation", "test-drive", "delivery", "service"]),
  startsAt: z.iso.datetime({ offset: true })
}).strict();

export async function POST(request: Request) {
  const key = request.headers.get("idempotency-key");
  if (!key || key.length < 16 || key.length > 128) return problem(400, "Falta idempotencia", "Envíe una clave estable para este intento.", "IDEMPOTENCY_KEY_REQUIRED");
  try {
    const input = schema.parse(await request.json());
    const result = await requestAppointment({ leadId: input.leadId, kind: input.kind, startsAt: input.startsAt, ...(input.modelId ? { modelId: input.modelId } : {}) }, key);
    return Response.json(result.value, { status: result.status, headers: { "cache-control": "no-store", ...(result.replayed ? { "idempotency-replayed": "true" } : {}) } });
  } catch (error) {
    return errorResponse(error);
  }
}
````

### FILE: `src/app/api/enterprise/appointments/route.test.ts`

```yaml
block_id: "TS-FRANCHISE-JOURNEY:file:02"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by the upstream references in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "f3a01f762c855aafdaef85a386374f94e944cf8fdfe9640fa85b5da5b7c8b53c"
variables: []
secrets_allowed: false
```

````ts
import { beforeEach, describe, expect, it, vi } from "vitest";
import { POST } from "./route";

const { requestAppointment } = vi.hoisted(() => ({ requestAppointment: vi.fn() }));
vi.mock("@/platform/backend/public-client", () => ({ requestAppointment }));

beforeEach(() => requestAppointment.mockReset());

describe("appointment BFF", () => {
  it("rejects requests without an idempotency key", async () => {
    const response = await POST(new Request("https://web.example.test/api/enterprise/appointments", { method: "POST", headers: { "content-type": "application/json" }, body: "{}" }));
    expect(response.status).toBe(400);
    expect(requestAppointment).not.toHaveBeenCalled();
  });

  it("validates and forwards an exact appointment", async () => {
    requestAppointment.mockResolvedValue({ value: { id: "appointment", state: "requested" }, status: 202, replayed: false });
    const response = await POST(new Request("https://web.example.test/api/enterprise/appointments", { method: "POST", headers: { "content-type": "application/json", "idempotency-key": "appointment-key-0001" }, body: JSON.stringify({ leadId: "lead", modelId: "model", kind: "test-drive", startsAt: "2026-09-01T15:00:00Z" }) }));
    expect(response.status).toBe(202);
    expect(requestAppointment).toHaveBeenCalledWith({ leadId: "lead", modelId: "model", kind: "test-drive", startsAt: "2026-09-01T15:00:00Z" }, "appointment-key-0001");
  });
});
````

### FILE: `src/app/locations/appointment-form.tsx`

```yaml
block_id: "TS-FRANCHISE-JOURNEY:file:03"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by the upstream references in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "c8e97aace140951615d2668faec98cd9a54f3577e15b202fce17d4ff1e31de77"
variables: []
secrets_allowed: false
```

````tsx
"use client";

import { publicAppointmentTime, publicCount, publicMessage as message, type PublicLocale } from "@/platform/i18n/public-catalog";
import { FormEvent, useEffect, useRef, useState } from "react";
import type { AppointmentSlot } from "@/platform/backend/public-client";

export function AppointmentForm({ leadId, modelId, slots, locale }: { leadId: string; modelId?: string; slots: AppointmentSlot[]; locale: PublicLocale }) {
  const key = useRef(crypto.randomUUID());
  const [pending, setPending] = useState(false);
  const [statusMessage, setMessage] = useState("");
  const [receiptId, setReceiptId] = useState("");
  const [ready, setReady] = useState(false);
  useEffect(() => { setReady(true); }, []);
  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault(); setPending(true); setMessage("");
    const form = new FormData(event.currentTarget);
    const slot = slots.find((item) => item.id === String(form.get("slotId") ?? ""));
    if (!slot) { setMessage(message(locale.locale, "appointment.select")); setPending(false); return; }
    try {
      const response = await fetch("/api/enterprise/appointments", { method: "POST", headers: { "content-type": "application/json", "idempotency-key": key.current }, body: JSON.stringify({ leadId, ...(modelId ? { modelId } : {}), kind: slot.kind, startsAt: slot.starts_at }) });
      if (!response.ok) { setMessage(message(locale.locale, "appointment.failed")); return; }
      const receipt = await response.json() as { id?: unknown; state?: unknown };
      if (typeof receipt.id !== "string" || !/^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(receipt.id) || receipt.state !== "requested") { setMessage(message(locale.locale, "appointment.invalid")); return; }
      setReceiptId(receipt.id);
      setMessage(`${message(locale.locale, "appointment.received")} ${receipt.id}. ${message(locale.locale, "appointment.confirmation")}`);
    } catch { setMessage(message(locale.locale, "appointment.failed")); }
    finally { setPending(false); }
  }
  return <form lang={locale.locale} onSubmit={submit} aria-describedby="appointment-result">
    <label>{message(locale.locale, "appointment.slot")}<select name="slotId" required defaultValue="" disabled={!ready || pending || receiptId !== ""}><option value="" disabled>{message(locale.locale, "appointment.option")}</option>{slots.map((slot) => <option value={slot.id} key={slot.id}>{message(locale.locale, "kind." + slot.kind)} · {publicAppointmentTime(locale, slot.starts_at)} · {publicCount(locale.locale, "appointment", slot.capacity - slot.booked)}</option>)}</select></label>
    {slots.length === 0 ? <p className="notice">{message(locale.locale, "appointment.empty")}</p> : null}
    <button type="submit" disabled={!ready || pending || slots.length === 0 || receiptId !== ""}>{pending ? message(locale.locale, "appointment.sending") : message(locale.locale, "appointment.submit")}</button>
    <p id="appointment-result" role="status" aria-live="polite">{statusMessage}</p>
  </form>;
}
````

### FILE: `src/app/locations/page.tsx`

```yaml
block_id: "TS-FRANCHISE-JOURNEY:file:04"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by the upstream references in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "294f6ea2f070a400e1e017683030a4b8d595bce6d4c0ff7b34023113bc5f8c5a"
variables: []
secrets_allowed: false
```

````tsx
import { publicMessage as message } from "@/platform/i18n/public-catalog";
import { loadPublicLocale } from "@/platform/i18n/load-public-locale";
import { listAppointmentSlots, listLocations, type AppointmentKind } from "@/platform/backend/public-client";
import { AppointmentForm } from "./appointment-form";

export default async function LocationsPage({ searchParams }: { searchParams: Promise<{ lead_id?: string; model_id?: string }> }) {
  const locale = await loadPublicLocale();
  const from = new Date(Date.now() + 60_000).toISOString();
  const to = new Date(Date.now() + 30 * 24 * 60 * 60_000).toISOString();
  const kinds: AppointmentKind[] = ["consultation", "test-drive", "delivery", "service"];
  const [locations, query, slotGroups] = await Promise.all([listLocations(), searchParams, Promise.all(kinds.map((kind) => listAppointmentSlots(kind, from, to)))]);
  const slots = slotGroups.flat().sort((a, b) => a.starts_at.localeCompare(b.starts_at));
  return <section lang={locale.locale}>
    <div className="eyebrow">{message(locale.locale, "locations.eyebrow")}</div><h1 className="pageTitle">{message(locale.locale, "locations.title")}</h1>
    <div className="grid">{locations.map((location) => <article className="card" key={location.organization_id}><h2>{location.name}</h2><p>{location.city}, {location.region} · {location.country}</p>{location.contact_phone ? <p><a href={`tel:${location.contact_phone}`}>{location.contact_phone}</a></p> : null}{location.contact_email ? <p><a href={`mailto:${location.contact_email}`}>{location.contact_email}</a></p> : null}</article>)}</div>
    {query.lead_id ? <section><h2>{message(locale.locale, "locations.book")}</h2><AppointmentForm leadId={query.lead_id} slots={slots} locale={locale} {...(query.model_id ? { modelId: query.model_id } : {})} /></section> : <div className="notice">{message(locale.locale, "locations.start")}</div>}
  </section>;
}
````

### FILE: `src/app/franchise/page.tsx`

```yaml
block_id: "TS-FRANCHISE-JOURNEY:file:05"
operation: CREATE
provenance: AUTHORED
source: "local implementation governed by the upstream references in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "215190e2750f48da05dceab9bc116146b1f5700d461cdce0e7d32905d9f7ecfd"
variables: []
secrets_allowed: false
```

````tsx
import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import { HandoverOperationsPanel } from "@/components/handover-operations-panel";
import { OperationalGuide } from "@/components/operational-guide";
import { OPERATIONAL_GUIDES } from "@/platform/help/content";
import { redirect } from "next/navigation";
import { createHash } from "node:crypto";
import type { Route } from "next";
import { allowed, readSession } from "@/platform/auth/session";
import { protectedGet, type LeadSummary, type Page } from "@/platform/backend/protected-client";
import { AppointmentAgendaPanel, FranchiseCommandPanel, OrderOperationsPanel, type OrderOperations, type AppointmentAgenda } from "@/components/franchise-command-panel";
import { loadBusinessConfig } from "@/platform/config/load";

type Availability = { id: string; organization_id: string; resource_id?: string; entry_type: "working" | "unavailable"; reason_code?: string; starts_at: string; ends_at: string; state: string; version: number };
type DeliveryException = { id: string; handover_id: string; customer_subject: string; reason_code: string; details: string; state: string; version: number; resolution_action?: string };
type ReturnEffect = { id: string; effect_kind: "inventory" | "refund" | "exchange" | "accounting" | "fiscal"; owner_context: string; state: "requested" };
type ReturnCase = { authorization_id: string; order_id: string; stock_unit_id: string; customer_subject: string; authorized_action: "return" | "exchange"; receipt?: { id: string; received_serial_number: string; condition_code: string; received_at: string }; disposition?: { id: string; inventory_action: string; customer_remedy: string; effect_requests: ReturnEffect[] } };

export default async function FranchisePage({ searchParams }: { searchParams: Promise<{ day?: string }> }) {
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

  const session = await readSession();
  if (!session) redirect("/api/auth/login?return_to=/franchise" as Route);
  const canOperate=["inventory:allocate","payment:create","handover:manage","admin:read"].some(permission=>allowed(session,permission));
  const canPanel=["lead:read","handover:manage","resource:manage","availability:read","availability:manage","appointment:manage"].some(permission=>allowed(session,permission));
  const canReviewWhatsApp=allowed(session,"whatsapp:approve");
  const canStoredValue=allowed(session,"stored_value:read");
  if (!canOperate && !canPanel && !canReviewWhatsApp && !canStoredValue) return <><h1 className="pageTitle">{t("p0006")}</h1><div className="notice">{t("p0074")}</div></>;
  const organization = session.organizations[0]!;
  const notificationStatusEnabled = (await loadBusinessConfig()).features.whatsapp_status_history === true;
  const from = new Date().toISOString();
  const to = new Date(Date.now() + 90 * 24 * 60 * 60_000).toISOString();
  const requestedDay=(await searchParams).day ?? new Date().toISOString().slice(0,10);
  const parsedDay=new Date(requestedDay+"T00:00:00Z");
  if (!/^\d{4}-\d{2}-\d{2}$/.test(requestedDay) || !Number.isFinite(parsedDay.getTime()) || parsedDay.toISOString().slice(0,10)!==requestedDay) return <><h1>{t("p0075")}</h1><a href="/franchise">{t("p0076")}</a></>;
  const agendaFrom=parsedDay.toISOString(), agendaTo=new Date(parsedDay.getTime()+24*60*60_000).toISOString();
  // Each owner may fail independently; null means unavailable, never an empty result.
  const [operationSnapshot,leads,availability,exceptions,returnCases,agenda]=await Promise.all([
    canOperate ? protectedGet<OrderOperations>(session,"/v1/commerce/orders",{organization_id:organization}).catch(()=>null) : Promise.resolve(null),
    allowed(session,"lead:read") ? protectedGet<Page<LeadSummary>>(session,"/v1/franchise/leads",{organization_id:organization,limit:"50"}).catch(()=>null) : Promise.resolve({items:[]}),
    allowed(session,"availability:read") ? protectedGet<{items:Availability[]}>(session,"/v1/franchise/availability",{organization_id:organization,from,to}).catch(()=>null) : Promise.resolve({items:[]}),
    allowed(session,"handover:manage") ? protectedGet<{items:DeliveryException[]}>(session,"/v1/franchise/delivery-exceptions",{organization_id:organization}).catch(()=>null) : Promise.resolve({items:[]}),
    allowed(session,"handover:manage") ? protectedGet<{items:ReturnCase[]}>(session,"/v1/franchise/returns",{organization_id:organization}).catch(()=>null) : Promise.resolve({items:[]}),
    allowed(session,"appointment:manage") ? protectedGet<AppointmentAgenda>(session,"/v1/franchise/agenda",{organization_id:organization,from:agendaFrom,to:agendaTo}).catch(()=>null) : Promise.resolve(null)
  ]);
  const failed=[
    leads===null ? "leads" : null,
    availability===null ? "disponibilidad" : null,
    exceptions===null ? t("p0077") : null,
    returnCases===null ? "devoluciones" : null,
    allowed(session,"appointment:manage") && agenda===null ? "agenda" : null
  ].filter((item):item is string=>item!==null);
  return <>
    <div className="eyebrow">{t("p0078")} {organization}</div>
    <h1 className="pageTitle">{t("p0079")}</h1>
    <p className="lede">{t("p0080")}</p>
    {canStoredValue ? <p><a className="button" href="/franchise/stored-value">{t("p0081")}</a></p> : null}
    {canReviewWhatsApp ? <p><a className="button" href="/franchise/whatsapp">{t("p0082")}</a></p> : null}
    {canOperate ? <OrderOperationsPanel snapshot={operationSnapshot} organization={organization} canAllocate={allowed(session,"inventory:allocate")} canRequestPayment={allowed(session,"payment:create")} paymentScope={createHash("sha256").update(JSON.stringify([session.tenantId,session.subject,organization])).digest("hex")}/> : null}
    {allowed(session,"handover:manage") ? <HandoverOperationsPanel orders={operationSnapshot && !operationSnapshot.truncated ? operationSnapshot.orders.map(o=>({id:o.id})) : null} organization={organization} scope={createHash("sha256").update(JSON.stringify([session.tenantId,session.subject,organization])).digest("hex")}/> : null}
    {failed.length>0 ? <section role="alert">{failed.map(name=><p key={name}>{t("p0083")} {name}{t("p0084")}</p>)}<a href={("/franchise?day="+requestedDay) as Route}>{t("p0085")}</a></section> : null}
    {allowed(session,"appointment:manage") ? <form action="/franchise" method="get"><label>{t("p0086")}<input type="date" name="day" defaultValue={requestedDay} required /></label><button className="button">{t("p0087")}</button></form> : null}
    {agenda ? <AppointmentAgendaPanel agenda={agenda} organization={organization} notificationStatusEnabled={notificationStatusEnabled}/> : null}
    {canPanel ? <FranchiseCommandPanel commandScope={createHash("sha256").update(JSON.stringify([session.tenantId,session.subject,organization])).digest("hex")} permissions={session.permissions} leads={leads?.items??[]} availability={availability?.items??[]} exceptions={exceptions?.items??null} returnCases={returnCases?.items??null} defaultAssignee={session.subject} organization={organization}/> : null}
    <OperationalGuide guide={OPERATIONAL_GUIDES["operation-sections-view"]} className="card"/>
  </>;
}
````

### FILE: `src/components/customer-appointment-actions.tsx`

```yaml
block_id: "TS-FRANCHISE-JOURNEY:file:11"
operation: CREATE
provenance: AUTHORED
source: "local customer self-service UI governed by Microsoft calendar/absence and Next.js references pinned in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "53c1b1036c933abcca77d84919783fb85dd118e605533b64b03d2d2927a1b0e1"
variables: []
secrets_allowed: false
```

````tsx
"use client";
import {usePrivateI18n} from "@/platform/i18n/private-provider";


import { useRef, useState } from "react";

type Appointment = { id: string; kind: string; starts_at: string; startsAtLabel: string; state: string; version: number };

// The response must identify this exact mutation, not merely be an HTTP 200.
export function isCancellationReceipt(value: unknown, organizationId: string, appointment: { id: string; version: number }): boolean {
  if (typeof value !== "object" || value === null || Array.isArray(value)) return false;
  const receipt = value as Record<string, unknown>;
  return receipt.id === appointment.id && receipt.organization_id === organizationId &&
    receipt.state === "cancelled" && Number.isSafeInteger(receipt.version) &&
    Number.isSafeInteger(appointment.version + 1) && receipt.version === appointment.version + 1;
}

export function CustomerAppointmentActions({ organizationId, appointments }: { organizationId: string; appointments: Appointment[] }) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

  const [pending, setPending] = useState("");
  const [message, setMessage] = useState("");
  const [uncertain, setUncertain] = useState(false);
  const inFlight = useRef(false);

  async function cancel(appointment: Appointment) {
    if (inFlight.current) return;
    inFlight.current = true;
    setPending(appointment.id);
    setMessage(t("p0231"));
    const controller = new AbortController();
    const timer = setTimeout(() => controller.abort(), 10000);
    try {
      const response = await fetch("/api/enterprise/franchise/commands", { method: "POST", headers: { "content-type": "application/json" }, signal: controller.signal, body: JSON.stringify({ action: "cancel-customer-appointment", organizationId, appointmentId: appointment.id, version: appointment.version, reasonCode: "customer-request" }) });
      const result: unknown = await response.json();
      if (!response.ok || !isCancellationReceipt(result, organizationId, appointment)) throw new Error("RESULT_NOT_CONFIRMED");
      setMessage(t("p0232"));
      window.location.reload();
    } catch {
      setUncertain(true);
      setMessage(t("p0233"));
      // Keep commands disabled; recover with an explicit GET, never an automatic POST retry.
    } finally {
      clearTimeout(timer);
    }
  }

  return <>
    {message && <p className="status" role="status" aria-live="polite">{message}</p>}
    {uncertain && <p><a href="/customer/appointments">{t("p0234")}</a></p>}
    {appointments.length === 0 && <p className="notice">{t("p0235")}</p>}
    <div className="grid">{appointments.map((appointment) => <article className="card" key={appointment.id}>
      <h2>{appointment.kind}</h2><p>{appointment.state} {t("p0016")} <time dateTime={appointment.starts_at}>{appointment.startsAtLabel}</time></p>
      {["requested", "confirmed"].includes(appointment.state) && new Date(appointment.starts_at).getTime() > Date.now() ? <button disabled={pending !== ""} onClick={() => void cancel(appointment)}>{t("p0236")}</button> : null}
    </article>)}</div>
  </>;
}
````

### FILE: `src/app/customer/appointments/page.tsx`

```yaml
block_id: "TS-FRANCHISE-JOURNEY:file:12"
operation: CREATE
provenance: AUTHORED
source: "local server-rendered customer appointment portal governed by Next.js and journey references pinned in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "bb1c7c2853757f29fca7236163ea0a4bd312c29aed6a5cbdfd6e7dc923af3527"
variables: []
secrets_allowed: false
```

````tsx
import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import { publicAppointmentTime } from "@/platform/i18n/public-catalog";
import Link from "next/link";
import { redirect } from "next/navigation";
import type { Route } from "next";
import { allowed, readSession } from "@/platform/auth/session";
import { protectedGet, type CustomerJourney } from "@/platform/backend/protected-client";
import { CustomerAppointmentActions } from "@/components/customer-appointment-actions";

export default async function CustomerAppointmentsPage() {
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

  const session = await readSession();
  if (!session) redirect("/api/auth/login?return_to=/customer/appointments" as Route);
  if (!allowed(session, "customer:self")) return <><h1 className="pageTitle">{t("p0006")}</h1><div className="notice">{t("p0007")} <code>customer:self</code>{t("p0008")}</div></>;
  const organization = session.organizations[0]!;
  const journey = await protectedGet<CustomerJourney>(session, "/v1/customer/journey", { organization_id: organization });
  const locale = privateLocale;
  const appointments = journey.appointments.map((item) => ({ ...item, startsAtLabel: `${publicAppointmentTime(locale, item.starts_at)} (${locale.timeZone})` }));
  return <>
    <div className="eyebrow">{t("p0034")} {organization}</div><h1 className="pageTitle">{t("p0035")}</h1>
    <p className="lede">{t("p0036")}</p>
    <p><Link href={"/customer/quotes" as Route}>{t("p0037")}</Link></p>
    <p><Link href={"/customer/handovers" as Route}>{t("p0038")}</Link></p>
    <CustomerAppointmentActions organizationId={organization} appointments={appointments}/>
  </>;
}
````

### FILE: `src/components/customer-handover-actions.tsx`

```yaml
block_id: "TS-FRANCHISE-JOURNEY:file:13"
operation: CREATE
provenance: AUTHORED
source: "local authenticated handover UI governed by Microsoft posted-shipment and inspection references pinned in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "10402e5757d5dc7e63ffd34e57360791210902375576088cc9aff50f4e52baa3"
variables: []
secrets_allowed: false
```

````tsx
"use client";
import {usePrivateI18n} from "@/platform/i18n/private-provider";


import { useRef, useState, type FormEvent } from "react";

type ChecklistItem = { id: string; ordinal: number; prompt: string; response_type: string; required: boolean };
type Handover = { id: string; order_id: string; stock_unit_id: string; state: string; version: number; customer_accepted_at?: string; checklist_id?: string; checklist_version?: number; checklist_title?: string; checklist_items: ChecklistItem[]; checklist_completed_at?: string; checklistCompletedLabel?: string; acceptedLabel?: string };
type DeliveryException = { id: string; handover_id: string; reason_code: string; details: string; state: string; version: number; resolution_action?: string; successor_handover_id?: string; return_authorization_id?: string };

export function CustomerHandoverActions({ organizationId, handovers, exceptions }: { organizationId: string; handovers: Handover[]; exceptions: DeliveryException[] }) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

  const [pending, setPending] = useState("");
  const [message, setMessage] = useState("");
  const inFlight = useRef(false);

  async function send(body: Record<string, unknown>) {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), 10_000);
    try {
      const response = await fetch("/api/enterprise/franchise/commands", { method: "POST", headers: { "content-type": "application/json" }, body: JSON.stringify(body), signal: controller.signal });
      const result: unknown = await response.json();
      if (!response.ok || !result || typeof result !== "object" || Array.isArray(result)) throw new Error("UNVERIFIED_RECEIPT");
      return result as Record<string, unknown>;
    } finally { clearTimeout(timeout); }
  }

  async function accept(event: FormEvent<HTMLFormElement>, handover: Handover) {
    event.preventDefault();
    if (inFlight.current) return;
    inFlight.current = true;
    const data = new FormData(event.currentTarget);
    setPending(handover.id);
    setMessage("");
    try {
      const result = await send({ action: "accept-handover", organizationId, handoverId: handover.id, version: handover.version, confirmedReceived: data.get("confirmedReceived") === "yes", serialNumber: String(data.get("serialNumber") ?? ""), checklistId: handover.checklist_id, checklistVersion: handover.checklist_version });
      if (result.id !== handover.id || result.state !== "accepted" || result.version !== handover.version + 1) throw new Error("UNVERIFIED_RECEIPT");
      setMessage(t("p0241"));
      window.location.reload();
    } catch {
      setMessage(t("p0242"));
    }
  }

  async function reject(event: FormEvent<HTMLFormElement>, handover: Handover) {
    event.preventDefault();
    if (inFlight.current) return;
    inFlight.current = true;
    const data = new FormData(event.currentTarget);
    setPending(handover.id);
    setMessage("");
    try {
      const result = await send({ action: "reject-handover", organizationId, handoverId: handover.id, version: handover.version, reasonCode: String(data.get("reasonCode") ?? ""), details: String(data.get("details") ?? "") });
      if (typeof result.id !== "string" || !result.id || result.handover_id !== handover.id || result.state !== "open" || result.version !== 1) throw new Error("UNVERIFIED_RECEIPT");
      setMessage(t("p0243"));
      window.location.reload();
    } catch {
      setMessage(t("p0244"));
    }
  }

  return <>{message && <><p className="status" role="status">{message}</p><p><a href="/customer/handovers">{t("p0245")}</a></p></>}{exceptions.length ? <section className="card"><h2>{t("p0246")}</h2>{exceptions.map((item) => <article key={item.id}><h3>{item.reason_code}</h3><p>{item.details}</p><p>{t("p0247")} {item.state}{item.resolution_action ? t("p0248", {resolution_action:(item.resolution_action)}) : ""}</p>{item.successor_handover_id ? <p>{t("p0249")} {item.successor_handover_id}</p> : null}{item.return_authorization_id ? <p>{t("p0250")} {item.return_authorization_id}</p> : null}</article>)}</section> : null}<div className="grid">{handovers.map((handover) => {
    const checklistReady = Boolean(handover.checklist_id && handover.checklist_version && handover.checklist_completed_at && handover.checklist_items.length);
    return <article className="card" key={handover.id}><h2>{t("p0251")} {handover.order_id}</h2><p>{t("p0247")} {handover.state}</p>{handover.checklist_id ? <section aria-label={t("p0252")}><h3>{handover.checklist_title || t("p0253")} {t("p0254")}{handover.checklist_version}</h3><ol>{handover.checklist_items.map((item) => <li key={item.id}>{item.prompt}{item.required ? " · obligatorio" : ""}</li>)}</ol><p>{handover.checklist_completed_at ? t("p0255", {completed_at:(handover.checklistCompletedLabel ?? t("p0256"))}) : t("p0257")}</p></section> : <p>{t("p0258")}</p>}{handover.state === "presented" && checklistReady ? <><form onSubmit={(event) => void accept(event, handover)}><label>{t("p0259")}<input name="serialNumber" required maxLength={128}/></label><label><input name="confirmedReceived" type="checkbox" value="yes" required/> {t("p0260")}</label><button disabled={pending !== ""}>{t("p0261")}</button></form><form onSubmit={(event) => void reject(event, handover)}><h3>{t("p0262")}</h3><label>{t("p0263")}<input name="reasonCode" required maxLength={64} pattern="[a-z][a-z0-9]*(-[a-z0-9]+)*" placeholder={t("p0264")}/></label><label>{t("p0265")}<textarea name="details" required maxLength={1000}/></label><button disabled={pending !== ""}>{t("p0266")}</button></form></> : <p>{handover.customer_accepted_at ? t("p0267", {accepted_at:(handover.acceptedLabel ?? t("p0256"))}) : t("p0268")}</p>}</article>;
  })}</div></>;
}
````

### FILE: `src/app/customer/handovers/page.tsx`

```yaml
block_id: "TS-FRANCHISE-JOURNEY:file:14"
operation: CREATE
provenance: AUTHORED
source: "local server-rendered delivery portal governed by Microsoft posted-shipment and inspection references pinned in metadata"
license: "LicenseRef-Workspace-Owner"
sha256: "9ddfb74a89069a82e952190f8e86d4d6c97e22e4a7564ade3466593e66cb9775"
variables: []
secrets_allowed: false
```

````tsx
import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import { publicAppointmentTime } from "@/platform/i18n/public-catalog";
import { OperationalGuide } from "@/components/operational-guide";
import { OPERATIONAL_GUIDES } from "@/platform/help/content";
import Link from "next/link";
import { redirect } from "next/navigation";
import type { Route } from "next";
import { allowed, readSession } from "@/platform/auth/session";
import { protectedGet, type CustomerJourney } from "@/platform/backend/protected-client";
import { CustomerHandoverActions } from "@/components/customer-handover-actions";

type DeliveryHandover = CustomerJourney["handovers"][number] & { checklist_id?: string; checklist_version?: number; checklist_title?: string; checklist_items: { id: string; ordinal: number; prompt: string; response_type: string; required: boolean }[]; checklist_completed_at?: string };
type DeliveryException = { id: string; handover_id: string; reason_code: string; details: string; state: string; version: number; resolution_action?: string; successor_handover_id?: string; return_authorization_id?: string };
type DeliveryJourney = Omit<CustomerJourney, "handovers"> & { handovers: DeliveryHandover[]; delivery_exceptions: DeliveryException[] };

export default async function CustomerHandoversPage() {
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

  const session = await readSession();
  if (!session) redirect("/api/auth/login?return_to=/customer/handovers" as Route);
  if (!allowed(session, "customer:self")) return <><h1 className="pageTitle">{t("p0006")}</h1><div className="notice">{t("p0007")} <code>customer:self</code>{t("p0008")}</div></>;
  const organization = session.organizations[0]!;
  const journey = await protectedGet<DeliveryJourney>(session, "/v1/customer/journey", { organization_id: organization }).catch(() => null);
  const locale = privateLocale;
  const handovers = journey?.handovers.map(item => ({ ...item, ...(item.checklist_completed_at ? { checklistCompletedLabel: `${publicAppointmentTime(locale, item.checklist_completed_at)} (${locale.timeZone})` } : {}), ...(item.customer_accepted_at ? { acceptedLabel: `${publicAppointmentTime(locale, item.customer_accepted_at)} (${locale.timeZone})` } : {}) })) ?? [];
  return <><div className="eyebrow">{t("p0034")} {organization}</div><h1 className="pageTitle">{t("p0039")}</h1><p className="lede">{t("p0040")}</p><p><Link href={"/customer/appointments" as Route}>{t("p0041")}</Link></p>
    {journey ? <>{journey.handovers.length === 0 ? <p>{t("p0042")}</p> : null}<CustomerHandoverActions organizationId={organization} handovers={handovers} exceptions={journey.delivery_exceptions}/></> : <><div className="notice" role="alert">{t("p0043")}</div><p><a href="/customer/handovers">{t("p0044")}</a></p></>}
    <OperationalGuide guide={OPERATIONAL_GUIDES["handover-read-view"]}/>
  </>;
}
````

### FILE: `src/platform/help/content.ts`
```yaml
block_id: "TS-FRANCHISE-JOURNEY-PORTALS:HELP:src-platform-help-content.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "Local read-only integration of existing public operational guides and admitted session/runtime; see VERSIONED_JOURNEY_HELP_V379.md"
license: "LicenseRef-Workspace-Owner"
sha256: "1c1ebc33ca2e578725f9d0ca4f51588f6f480e81ab9e8cbbdfc05b963395a581"
variables: []
secrets_allowed: false
```
````typescript
// Release-bound PUBLIC operational guidance, shared with inline help. Not tenant documents.
import { STATUS_VIEW_VERSION } from "@/platform/notifications/status-contract";
export const QUOTE_GUIDE = {
  id: "quote-acceptance-view", version: "1.0.0", title: "Aceptar una cotización",
  paragraphs: ["Revisá el importe, la moneda y la vigencia. Aceptar solicita crear un pedido con esos datos del servidor.", "Aceptar una cotización no confirma el pago, el stock ni la entrega.", "Si se corta la conexión, usá Actualizar estado o volvé a abrir Mis cotizaciones. No repitas la aceptación hasta comprobar el estado.", "Un pedido se confirma aquí cuando la lectura muestra su referencia. Si el problema continúa, contactá al soporte de la empresa con la referencia de cotización y organización, sólo por un canal autorizado; no compartas contraseñas ni tokens.", "Práctica de esta versión en el entorno de capacitación: reconocer un resultado incierto, recuperar la lectura y ubicar el pedido sin reenviar. No practiques creando pedidos en producción."]
} as const;
export const NOTIFICATION_GUIDE = {
  id: "whatsapp-status-view", version: STATUS_VIEW_VERSION.split("/")[1], title: "Consultar el estado de WhatsApp",
  paragraphs: ["Compartí esta referencia sólo con soporte autorizado de tu organización. No adjuntes teléfonos, tokens ni el contenido del mensaje.", "Si el resultado es incierto o contradictorio: preservá el intento, consultá su evidencia y no reenvíes ni borres el historial. La entrega informada no confirma una venta ni la aceptación del turno."]
} as const;

// Extracted verbatim from the existing versioned inline guidance, V380.
export const OPERATIONAL_GUIDES = {
  "availability-cancel-view": {
    "id": "availability-cancel-view",
    "version": "1.0.0",
    "title": "Cancelar un intervalo",
    "summary": "Ayuda para cancelar un intervalo",
    "paragraphs": [
      "Cancelar este registro no cancela ni reprograma citas de clientes. El servidor comprueba los permisos, la versión y las restricciones existentes. Revisá recurso, fechas y tipo antes de actuar.",
      "Si se pierde la respuesta, Consultar intervalo sólo lee el servidor. Una cancelación consultada no demuestra quién la solicitó; la auditoría autorizada conserva al actor. La referencia de esta pestaña guarda sólo la versión, no motivos, horarios ni datos personales.",
      "Si la consulta falla, el intervalo no aparece en la ventana consultada o permanece activo, conservamos el bloqueo: solicitá al soporte autorizado revisar la referencia, hora y estado. No borres la referencia ni abras otra pestaña para forzar un reintento.",
      "Práctica en pruebas: cancelar un intervalo sin citas, perder la respuesta y consultar Cancelado sin reenviar. No practicar con la agenda real. Reactivar o crear otro intervalo es una operación distinta."
    ],
    "inlineLead": false
  },
  "availability-create-view": {
    "id": "availability-create-view",
    "version": "1.0.0",
    "title": "Registrar un intervalo",
    "summary": "Ayuda para registrar un intervalo",
    "paragraphs": [
      "Conservamos sólo la referencia de operación. Una respuesta perdida no demuestra rechazo: consultá antes de volver a actuar. La consulta muestra el registro aunque quede fuera del rango de la agenda. Preparar otro intervalo no reenvía ni cancela el anterior y sólo se habilita tras recuperarlo. Ante referencia ilegible o consulta fallida, pedí revisión autorizada sin borrar datos. Practicá con datos sintéticos el envío, pérdida de respuesta, consulta y preparación explícita de otro intervalo."
    ],
    "inlineLead": false
  },
  "checklist-completion-view": {
    "id": "checklist-completion-view",
    "version": "1.0.0",
    "title": "Presentar una entrega",
    "summary": "Ayuda para presentar una entrega",
    "paragraphs": [
      "La presentación conserva respuestas inmutables. Guardamos sólo identidades, versiones y una huella; no guardamos las respuestas ni evidencias en la referencia del navegador. Una respuesta perdida puede ocultar un cambio correcto: consultá sin reenviarlo. El estado actual puede avanzar a aceptado o rechazado; consultar o preparar otro formulario no lo cambia. Si el contenido no coincide, pedí revisión. Practicá con datos sintéticos antes de operar."
    ],
    "inlineLead": false
  },
  "checklist-publication-view": {
    "id": "checklist-publication-view",
    "version": "1.0.0",
    "title": "Publicar una checklist",
    "summary": "Ayuda para publicar una checklist",
    "paragraphs": [
      "Una versión publicada es inmutable. Conservamos sólo su ID, versión y huella de contenido; no guardamos el título ni los textos en la referencia del navegador. Una respuesta perdida puede ocultar una publicación correcta: consultá sin volver a enviarla. Si no coincide el contenido o no se puede recuperar, pedí revisión sin borrar la referencia. Preparar otra versión no cambia ni vuelve a publicar la anterior. Practicá este recorrido con datos sintéticos antes de usarlo en operación."
    ],
    "inlineLead": false
  },
  "delivery-resolution-view": {
    "id": "delivery-resolution-view",
    "version": "1.0.0",
    "title": "Resolver una discrepancia",
    "summary": "Ayuda para resolver una discrepancia",
    "paragraphs": [
      "Una respuesta perdida no demuestra rechazo. Conservamos en esta pestaña sólo la versión y la decisión pendiente, sin notas ni credenciales. Consultar vuelve a leer el servidor sin reenviar la acción.",
      "Si la consulta falla, el registro no aparece o hay otra decisión, pedí revisión al responsable autorizado. No borres la referencia ni abras otra pestaña para forzar un reintento. Informá la referencia de la entrega y la hora, sin compartir notas privadas ni tokens.",
      "Práctica: en el entorno de prueba, recuperá una resolución tras perder la respuesta y comprobá su estado antes del siguiente paso. Preparada o autorizada no significa entrega física, reembolso ni cambio finalizados."
    ],
    "inlineLead": false
  },
  "handover-read-view": {
    "id": "handover-read-view",
    "version": "1.0.0",
    "title": "Consultar entregas",
    "summary": "Ayuda para consultar entregas",
    "paragraphs": [
      "Si la información no está disponible, volvé a consultar antes de aceptar o rechazar. No crees otra solicitud para recuperar una lectura.",
      "Si el problema continúa, contactá al soporte autorizado e indicá la acción y la hora. No envíes contraseñas, tokens ni información de otras personas.",
      "Práctica en pruebas: interrumpí la consulta, verificá que no haya acciones de entrega disponibles y restablecé la conexión. La próxima consulta debe recuperar el estado sin enviar una aceptación automática."
    ],
    "inlineLead": true
  },
  "lead-command-view": {
    "id": "lead-command-view",
    "version": "1.0.0",
    "title": "Actualizar un lead",
    "summary": "Ayuda para actualizar un lead",
    "paragraphs": [
      "Una respuesta perdida no demuestra que el cambio haya sido rechazado. Consultar lee el servidor sin reenviar comandos. Una versión posterior permite revisar el estado real; no prueba que haya sido modificado exclusivamente por tu solicitud.",
      "Conservamos sólo versión y tipo de acción en esta pestaña, no responsables ni datos de contacto. Si falta el lead, falla su consulta o la versión no avanzó, pedí revisión al soporte autorizado; no borres la referencia ni abras otra pestaña para forzar un reintento.",
      "Práctica: recuperá un cambio en el entorno de prueba, comprobá responsable y estado y continuá con la versión consultada. No asumas una venta o un cliente convertido por una asignación."
    ],
    "inlineLead": false
  },
  "operation-sections-view": {
    "id": "operation-sections-view",
    "version": "1.0.0",
    "title": "Recuperar una sección",
    "summary": "Ayuda para recuperar una sección",
    "paragraphs": [
      "Una consulta fallida no significa que no existan registros. Volvé a consultar antes de decidir sobre esa sección; las secciones independientes conservan sus propios permisos y estados.",
      "La consulta no reenvía acciones. Si el problema persiste, informá al soporte autorizado la sección, la hora y la operación esperada; no compartas tokens, contraseñas ni datos personales en capturas.",
      "Práctica: ante una consulta fallida, identificá la sección, recuperala mediante consulta y verificá su estado antes de actuar."
    ],
    "inlineLead": false
  },
  "order-operations-view": {
    "id": "order-operations-view",
    "version": "1.1.0",
    "title": "Continuar un pedido",
    "summary": "Ayuda para continuar un pedido",
    "paragraphs": [
      "Elegí una unidad por su serie y reservá una sola vez. El servidor vuelve a comprobar organización, variante, disponibilidad y versiones.",
      "Si se pierde la respuesta o alguien actuó primero, usá Actualizar pedidos. No crees otro pedido ni cambies identificadores para forzar la operación.",
      "Pago: registrar una solicitud guarda una intención del total; no confirma dinero recibido ni libera entrega. Una solicitud existente se consulta, no se duplica. Si falla, no borres la clave del navegador: consultá el estado y escalá al responsable antes de intentar otro pago.",
      "Práctica segura: en pruebas, dos operadores solicitan el pago del mismo pedido; sólo debe existir una intención inicial. Simulá respuesta perdida y recuperá el registro con Actualizar pedidos. No practiques con pedidos reales.",
      "Para soporte, comunicá la referencia del pedido o solicitud, acción y hora al responsable autorizado. No adjuntes tokens ni datos de clientes. Cobro, conciliación y entrega mantienen sus gates separados."
    ],
    "inlineLead": true
  },
  "quote-create-view": {
    "id": "quote-create-view",
    "version": "1.0.0",
    "title": "Emitir una cotización",
    "summary": "Ayuda para emitir una cotización",
    "paragraphs": [
      "El servidor determina precio y moneda. La referencia conserva sólo un identificador de operación. Una respuesta perdida no permite emitir otra: consultá el resultado. Si no aparece, conservá la referencia y pedí revisión al responsable autorizado. Practicá emisión, respuesta perdida y consulta con datos sintéticos; no borres la referencia para repetir una cotización."
    ],
    "inlineLead": false
  },
  "resource-create-view": {
    "id": "resource-create-view",
    "version": "1.0.0",
    "title": "Registrar un recurso",
    "summary": "Ayuda para registrar un recurso",
    "paragraphs": [
      "Conservamos sólo la referencia de operación. Una respuesta perdida no demuestra rechazo: consultá antes de volver a actuar. La consulta identifica el registro por la referencia de esta operación, sin buscar por nombre. Preparar otro recurso no reenvía ni cancela el anterior y sólo se habilita tras recuperarlo. Ante referencia ilegible o consulta fallida, pedí revisión autorizada sin borrar datos. Practicá con datos sintéticos el envío, pérdida de respuesta, consulta y preparación explícita de otro recurso."
    ],
    "inlineLead": false
  },
  "return-operations-view": {
    "id": "return-operations-view",
    "version": "1.0.0",
    "title": "Recibir y decidir devoluciones",
    "summary": "Ayuda para recibir y decidir devoluciones",
    "paragraphs": [
      "Recepción, decisión y solicitudes conservan evidencia inmutable. La referencia del navegador guarda sólo operación, IDs y huella; no guarda serie, notas ni tokens. Una respuesta perdida puede ocultar un cambio correcto. Consultá el caso exacto aunque la lista falle o no lo muestre. Si la evidencia no coincide, pedí revisión. Continuar operando actualiza el caso consultado y habilita una acción nueva; no reenvía la anterior. Una devolución con reembolso crea cuatro solicitudes; un cambio crea tres y no solicita emisión fiscal. Practicá con datos sintéticos antes de operar."
    ],
    "inlineLead": false
  },
  "slot-create-view": {
    "id": "slot-create-view",
    "version": "1.0.0",
    "title": "Registrar un turno",
    "summary": "Ayuda para registrar un turno",
    "paragraphs": [
      "Conservamos sólo la referencia de operación. Una respuesta perdida no demuestra rechazo: consultá antes de volver a actuar. La consulta muestra el turno por su referencia, aunque esté cerrado, completo o fuera del listado público. Publicar no confirma una reserva; siguen vigentes la jornada y las restricciones de agenda. Preparar otro turno no reenvía ni cancela el anterior y sólo se habilita tras recuperarlo. Ante referencia ilegible o consulta fallida, pedí revisión autorizada sin borrar datos. Practicá con datos sintéticos el envío, pérdida de respuesta, consulta y preparación explícita de otro turno."
    ],
    "inlineLead": false
  }
} as const;

// Same-release operational guidance. AUTHORED integration text; no corporate attribution.
export const RELEASE_GUIDES = {
  "supply-role-view": {
    "id": "supply-role-view",
    "version": "1.0.0",
    "title": "Guía de suministro",
    "summary": "Guía de suministro · versión 1.0.0",
    "paragraphs": [
      "Compras registra proveedor, destino, fábrica, importe y cantidades. La fábrica confirma la orden, registra cada serie y pide revisión de calidad. Otra persona autorizada decide su liberación o rechazo con evidencia.",
      "El despacho identifica un envío y sus series. Recepción registra sólo lo recibido; queda en cuarentena hasta la decisión de otra persona autorizada. Las cantidades y versiones se comprueban en cada operación.",
      "Si se pierde una respuesta, consultá el resultado pendiente. Esa consulta no vuelve a enviar. No borres referencias para forzar un reintento. En capacitación, practicá con datos sintéticos el recorrido y la recuperación."
    ],
    "inlineLead": false
  },
  "warranty-role-view": {
    "id": "warranty-role-view",
    "version": "1.0.0",
    "title": "Guía de garantía",
    "summary": "Guía de garantía · versión 1.0.0",
    "paragraphs": [
      "Leé los términos vinculados a la cotización antes de acusar su recepción. La garantía activa conserva esos términos vendidos. Un reclamo necesita entrega y atención registradas.",
      "Diagnóstico, plan y aprobación preceden al trabajo. La persona revisora es diferente de quien propone; calidad es diferente de quien ejecuta. El cliente acepta y la fábrica concilia el servicio. La conciliación es un acuse interno, no un pago.",
      "Si se pierde una respuesta, consultá el resultado pendiente: no vuelve a enviar. Conservá su referencia y pedí revisión autorizada si sigue incierto. Practicá con datos sintéticos en capacitación."
    ],
    "inlineLead": false
  },
  "network-role-view": {
    "id": "network-role-view",
    "version": "1.0.0",
    "title": "Guía de red y acuerdos",
    "summary": "Guía de red y acuerdos · versión 1.0.0",
    "paragraphs": [
      "Registrá la organización bajo un padre autorizado y usá la referencia devuelta para consultar su estado. Crear una organización no concede accesos a personas. Los nuevos accesos se asignan mediante el proveedor de identidad autorizado.",
      "La estructura activa permite preparar acuerdos y sucursales. No equivale a autorización de despliegue o producción. Territorio, versión de términos y fechas deben proceder del acuerdo revisado; no se generan condiciones comerciales aquí.",
      "Antes de operar, completá la capacitación de tu rol. La evaluación no concede permisos. Si una respuesta se pierde, recuperá su resultado y luego consultá el estado actual. No repitas una escritura incierta.",
      "Para cerrar una organización, revisá primero sus sucursales y acuerdos activos. No se eliminan los registros ni se revocan sesiones desde este formulario."
    ],
    "inlineLead": false
  },
  "help-cms-view": {
    "id": "help-cms-view",
    "version": "1.0.0",
    "title": "Guía de contenidos",
    "summary": "Guía de contenidos · versión 1.0.0",
    "paragraphs": [
      "El contenido pertenece a una organización y un idioma. Guardá un borrador, revisalo y publicalo con el permiso correspondiente. Archivar retira el artículo de la consulta de lectores y conserva su historial.",
      "Las versiones publicadas no se editan. Para reemplazarlas, prepará un artículo nuevo y archivá el anterior después de revisar la sustitución. El texto se muestra literalmente; no ejecuta HTML.",
      "Publicar aquí no modifica cursos aprobados ni concede permisos. La incorporación a capacitación requiere revisión y una nueva versión del curso."
    ],
    "inlineLead": false
  },
  "catalog-role-view": {
    "id": "catalog-role-view",
    "version": "1.0.0",
    "title": "Guía de catálogo y publicación",
    "summary": "Guía de catálogo · versión 1.0.0",
    "inlineLead": false,
    "paragraphs": [
      "Creá modelo, variantes y precios con importes y vigencias explícitos. Conservá las referencias devueltas, incorporá la imagen PNG y prepará una versión del catálogo.",
      "La versión requiere las revisiones legal, técnica, de imagen y de publicación. Cada decisión refiere al contenido exacto y su evidencia. Publicar comprueba la generación vigente; una revisión desactualizada no autoriza contenido nuevo.",
      "Consultar un resultado pendiente recupera la operación sin repetirla. Antes de publicar otra versión o volver a una anterior, consultá la publicación vigente. La publicación del catálogo no confirma stock, cobro ni habilitación comercial. Practicá únicamente con datos sintéticos."
    ]
  },
  "training-role-view": {
    "id": "training-role-view",
    "version": "1.0.0",
    "title": "Guía de capacitación y evaluación",
    "summary": "Guía de capacitación · versión 1.0.0",
    "inlineLead": false,
    "paragraphs": [
      "Elegí la guía del rol e iniciá una práctica. Registrá la lectura de cada lección y presentá respuestas para revisión humana. Elegir una guía no cambia los permisos de tu sesión.",
      "Una persona autorizada diferente revisa las respuestas contra la versión presentada. Una evaluación favorable conserva evidencia, pero no asigna accesos ni reemplaza la aceptación del responsable de la operación.",
      "Ante una respuesta incierta, consultá el intento guardado. Un curso nuevo no cambia la evidencia de intentos anteriores; prepará otra práctica para la versión activa. Los artículos privados del CMS requieren incorporación explícita a un curso revisado."
    ]
  }
} as const;

export const ALL_GUIDES = [QUOTE_GUIDE, NOTIFICATION_GUIDE, ...Object.values(OPERATIONAL_GUIDES), ...Object.values(RELEASE_GUIDES)] as const;
export type GuideID = typeof ALL_GUIDES[number]["id"];
````

### FILE: `src/platform/help/contract.ts`
```yaml
block_id: "TS-FRANCHISE-JOURNEY-PORTALS:HELP:src-platform-help-contract.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "Local read-only integration of existing public operational guides and admitted session/runtime; see VERSIONED_JOURNEY_HELP_V379.md"
license: "LicenseRef-Workspace-Owner"
sha256: "e452582a735b7cce3a01e025fcf3be33fc7e8f4106783e3f4a450880f28528e2"
variables: []
secrets_allowed: false
```
````typescript
import { z } from "zod";
import { ALL_GUIDES } from "./content";
export const helpResponseSchema = z.object({
  schema: z.literal("journey-help/1"),
  articles: z.array(z.object({
    id: z.enum(ALL_GUIDES.map(guide => guide.id)),
    version: z.string().regex(/^\d+\.\d+\.\d+$/).max(32),
    title: z.string().min(1).max(160),
    paragraphs: z.array(z.string().min(1).max(2000)).min(1).max(20)
  }).strict()).max(ALL_GUIDES.length)
}).strict();
export type HelpResponse = z.infer<typeof helpResponseSchema>;
````

### FILE: `src/platform/help/catalog.ts`
```yaml
block_id: "TS-FRANCHISE-JOURNEY-PORTALS:HELP:src-platform-help-catalog.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "Local read-only integration of existing public operational guides and admitted session/runtime; see VERSIONED_JOURNEY_HELP_V379.md"
license: "LicenseRef-Workspace-Owner"
sha256: "e35adc9a1d42ed9345c0e016a357a2f2b9fefd752a3beedadb33282ec15639f4"
variables: []
secrets_allowed: false
```
````typescript
import{parseHelpQuery}from"./query";
export{parseHelpQuery}from"./query";
import "server-only";
import { allowed, type PortalSession } from "@/platform/auth/session";
import { ALL_GUIDES, type GuideID } from "./content";

// Match the existing page/section permission predicates. Reading is not command authorization.
const operationPermissions = ["inventory:allocate", "payment:create", "handover:manage", "admin:read"];
const sectionPermissions = [...operationPermissions, "lead:read", "resource:manage", "availability:read", "availability:manage", "appointment:manage"];
export function availableGuides(session: PortalSession, notificationsEnabled: boolean) {
  const can = (permission: string) => allowed(session, permission);
  const policy: Record<GuideID, boolean> = {
    "catalog-role-view": can("catalog:read"),
    "supply-role-view": can("supply:read") || can("supply:factory-read"),
    "warranty-role-view": can("warranty:read") || can("warranty:self") || can("warranty:factory-read"),
    "network-role-view": can("network:admin") || can("franchise:write"),
    "help-cms-view": can("help:read") || can("help:write") || can("help:publish"),
    "training-role-view": can("training:learn") || can("training:review"),

    "quote-acceptance-view": can("customer:self"),
    "handover-read-view": can("customer:self"),
    "whatsapp-status-view": notificationsEnabled && can("appointment:manage"),
    "order-operations-view": operationPermissions.some(can),
    "operation-sections-view": sectionPermissions.some(can),
    "lead-command-view": can("lead:read"),
    "quote-create-view": can("lead:read") && can("quote:write"),
    "availability-cancel-view": can("availability:read"),
    "availability-create-view": can("availability:manage"),
    "resource-create-view": can("resource:manage"),
    "slot-create-view": can("appointment:manage"),
    "delivery-resolution-view": can("handover:manage"),
    "checklist-publication-view": can("handover:manage"),
    "checklist-completion-view": can("handover:manage"),
    "return-operations-view": can("handover:manage")
  };
  // Project only the response contract, excluding inline-rendering metadata.
  return ALL_GUIDES.filter(g => policy[g.id]).map(({ id, version, title, paragraphs }) => ({ id, version, title, paragraphs }));
}
export function selectGuides(guides: ReturnType<typeof availableGuides>, query: NonNullable<ReturnType<typeof parseHelpQuery>>) {
  if (query.article !== null) return guides.filter(g => g.id === query.article && g.version === query.version);
  return guides.filter(g => [g.title, ...g.paragraphs].join("\n").normalize("NFC").toLocaleLowerCase("es").includes(query.q));
}
````

### FILE: `src/app/api/enterprise/help/route.ts`
```yaml
block_id: "TS-FRANCHISE-JOURNEY-PORTALS:HELP:src-app-api-enterprise-help-route.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "Local read-only integration of existing public operational guides and admitted session/runtime; see VERSIONED_JOURNEY_HELP_V379.md"
license: "LicenseRef-Workspace-Owner"
sha256: "da176d09f165554cc7bba601dbf95bb65ca6d896d123876053616cfb9455cc29"
variables: []
secrets_allowed: false
```
````typescript
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
````

### FILE: `src/app/api/enterprise/help/route.test.ts`
```yaml
block_id: "TS-FRANCHISE-JOURNEY-PORTALS:HELP:src-app-api-enterprise-help-route.test.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "Local read-only integration of existing public operational guides and admitted session/runtime; see VERSIONED_JOURNEY_HELP_V379.md"
license: "LicenseRef-Workspace-Owner"
sha256: "232505c91b6ac4c950f82bace4a10e6ddb1e3ea1d30b2f9c09d2f836e9b321b3"
variables: []
secrets_allowed: false
```
````typescript
import { beforeEach, expect, test, vi } from "vitest";
const mocks = vi.hoisted(() => ({ readSession: vi.fn(), loadBusinessConfig: vi.fn() }));
vi.mock("@/platform/auth/session", () => ({ readSession: mocks.readSession, allowed: (s: {permissions: string[]}, p: string) => s.permissions.includes(p) || s.permissions.includes("*") }));
vi.mock("@/platform/config/load", () => ({ loadBusinessConfig: mocks.loadBusinessConfig }));
import { GET } from "./route";
import { helpResponseSchema } from "@/platform/help/contract";
import { QUOTE_GUIDE, NOTIFICATION_GUIDE, OPERATIONAL_GUIDES, ALL_GUIDES } from "@/platform/help/content";
import { readFileSync } from "node:fs";
const session = (permissions: string[]) => ({ subject: "private-subject", tenantId: "private-tenant", permissions, organizations: ["private-org"], accessToken: "never-return-token" });
const get = (query = "", headers = {}) => GET(new Request("https://example.test/api/enterprise/help" + query, { headers }));
beforeEach(() => { vi.resetAllMocks(); mocks.readSession.mockResolvedValue(session(["customer:self"])); mocks.loadBusinessConfig.mockResolvedValue({ features: { whatsapp_status_history: true } }); });
test("customer receives only applicable versioned public help without identity", async () => {
 const r = await get(); expect(r.status).toBe(200); expect(r.headers.get("cache-control")).toBe("no-store"); expect(r.headers.get("x-content-type-options")).toBe("nosniff");
 const data = helpResponseSchema.parse(await r.json()); expect(data.articles.map(g=>g.id)).toEqual(["quote-acceptance-view","handover-read-view"]); expect(JSON.stringify(data)).not.toMatch(/private-|never-return-token/);
});
test("operator receives only enabled notification guidance", async () => { mocks.readSession.mockResolvedValue(session(["appointment:manage"])); expect((await (await get()).json()).articles.map((g:{id:string})=>g.id).sort()).toEqual(["operation-sections-view","slot-create-view","whatsapp-status-view"]); mocks.loadBusinessConfig.mockResolvedValue({ features: {} }); expect((await (await get()).json()).articles.map((g:{id:string})=>g.id).sort()).toEqual(["operation-sections-view","slot-create-view"]); });
test("wildcard still obeys feature selection", async () => { mocks.readSession.mockResolvedValue(session(["*"])); expect((await (await get()).json()).articles).toHaveLength(21); mocks.loadBusinessConfig.mockResolvedValue({ features: {} }); expect((await (await get()).json()).articles.map((g:{id:string})=>g.id).sort()).toEqual(ALL_GUIDES.filter(g=>g.id!=="whatsapp-status-view").map(g=>g.id).sort()); });
test("no permission is denied and rechecked on subsequent reads", async () => { expect((await get()).status).toBe(200); mocks.readSession.mockResolvedValue(session(["factory:read"])); const r=await get(); expect(r.status).toBe(403); expect(await r.json()).toEqual({ code: "FORBIDDEN" }); });
test("missing session returns401 without loading configuration", async () => { mocks.readSession.mockResolvedValue(null); expect((await get()).status).toBe(401); expect(mocks.loadBusinessConfig).not.toHaveBeenCalled(); });
test.each(["cross-site", "none"])("rejects %s before credentials", async value => { expect((await get("", { "sec-fetch-site": value })).status).toBe(403); expect(mocks.readSession).not.toHaveBeenCalled(); });
test("exact version succeeds and unavailable version never silently upgrades", async () => { expect((await get("?article=quote-acceptance-view&version=1.0.0")).status).toBe(200); for(const q of ["?article=quote-acceptance-view&version=9.0.0", "?article=whatsapp-status-view&version=1.0.0", "?article=unknown&version=1.0.0"]) { const r=await get(q); expect(r.status).toBe(404); expect(await r.json()).toEqual({code:"GUIDE_NOT_AVAILABLE"}); } });
test("literal search is case/NFC insensitive and does not interpret regex or markup", async () => { for(const q of ["COTIZACIÓN", "cotizacio\u0301n"]) expect((await (await get("?q="+encodeURIComponent(q))).json()).articles).toHaveLength(1); for(const q of [".*", "<script>alert(1)</script>", "unmatched"]) expect((await (await get("?q="+encodeURIComponent(q))).json()).articles).toHaveLength(0); });
test.each(["?q=a&q=b", "?tenant=other", "?role=*", "?article=quote-acceptance-view", "?version=1.0.0", "?article=../escape&version=1.0.0", "?article=quote-acceptance-view&version=latest", "?article=quote-acceptance-view&version=1.0.0&q=", "?q=%00", "?q="+"x".repeat(161), "?article="+"x".repeat(2050)])("invalid query rejects before credentials: %s", async query => { expect((await get(query)).status).toBe(400); expect(mocks.readSession).not.toHaveBeenCalled(); });
test("unexpected session/config failures are sanitized and not cached", async () => { mocks.readSession.mockRejectedValue(new Error("private-secret")); let r=await get(); expect(r.status).toBe(503); expect(await r.text()).not.toContain("private-secret"); mocks.readSession.mockResolvedValue(session(["*"])); mocks.loadBusinessConfig.mockRejectedValue(new Error("config-secret")); r=await get(); expect(r.status).toBe(503); expect(r.headers.get("cache-control")).toBe("no-store"); expect(await r.json()).toEqual({code:"HELP_UNAVAILABLE"}); });
test("inline guides share the source and retain their safety instructions", () => { const quote=readFileSync("src/components/customer-quote-actions.tsx","utf8"), notification=readFileSync("src/components/appointment-notification-status.tsx","utf8"); expect(quote).toContain("QUOTE_GUIDE.paragraphs.map"); expect(notification).toContain("NOTIFICATION_GUIDE.paragraphs[1]"); expect(QUOTE_GUIDE.paragraphs.join(" ")).toContain("no confirma el pago, el stock ni la entrega"); expect(NOTIFICATION_GUIDE.paragraphs.join(" ")).toContain("no reenvíes ni borres el historial"); });

const roleMatrix = [
  {
    "role": "none",
    "permissions": [],
    "expected": []
  },
  {
    "role": "factory",
    "permissions": [
      "factory:read"
    ],
    "expected": []
  },
  {
    "role": "customer",
    "permissions": [
      "customer:self"
    ],
    "expected": [
      "handover-read-view",
      "quote-acceptance-view"
    ]
  },
  {
    "role": "lead",
    "permissions": [
      "lead:read"
    ],
    "expected": [
      "lead-command-view",
      "operation-sections-view"
    ]
  },
  {
    "role": "quote-without-lead-read",
    "permissions": [
      "quote:write"
    ],
    "expected": []
  },
  {
    "role": "quote",
    "permissions": [
      "lead:read",
      "quote:write"
    ],
    "expected": [
      "lead-command-view",
      "operation-sections-view",
      "quote-create-view"
    ]
  },
  {
    "role": "availability-read",
    "permissions": [
      "availability:read"
    ],
    "expected": [
      "availability-cancel-view",
      "operation-sections-view"
    ]
  },
  {
    "role": "availability-manage",
    "permissions": [
      "availability:manage"
    ],
    "expected": [
      "availability-create-view",
      "operation-sections-view"
    ]
  },
  {
    "role": "resource",
    "permissions": [
      "resource:manage"
    ],
    "expected": [
      "operation-sections-view",
      "resource-create-view"
    ]
  },
  {
    "role": "appointment",
    "permissions": [
      "appointment:manage"
    ],
    "expected": [
      "operation-sections-view",
      "slot-create-view",
      "whatsapp-status-view"
    ]
  },
  {
    "role": "inventory",
    "permissions": [
      "inventory:allocate"
    ],
    "expected": [
      "operation-sections-view",
      "order-operations-view"
    ]
  },
  {
    "role": "payment",
    "permissions": [
      "payment:create"
    ],
    "expected": [
      "operation-sections-view",
      "order-operations-view"
    ]
  },
  {
    "role": "handover",
    "permissions": [
      "handover:manage"
    ],
    "expected": [
      "checklist-completion-view",
      "checklist-publication-view",
      "delivery-resolution-view",
      "operation-sections-view",
      "order-operations-view",
      "return-operations-view"
    ]
  },
  {
    "role": "admin",
    "permissions": [
      "admin:read"
    ],
    "expected": [
      "operation-sections-view",
      "order-operations-view"
    ]
  },
  {
    "role": "wildcard",
    "permissions": [
      "*"
    ],
    "expected": [
      "availability-cancel-view",
      "availability-create-view",
      "catalog-role-view",
      "checklist-completion-view",
      "checklist-publication-view",
      "delivery-resolution-view",
      "handover-read-view",
      "help-cms-view",
      "lead-command-view",
      "network-role-view",
      "operation-sections-view",
      "order-operations-view",
      "quote-acceptance-view",
      "quote-create-view",
      "resource-create-view",
      "return-operations-view",
      "slot-create-view",
      "supply-role-view",
      "training-role-view",
      "warranty-role-view",
      "whatsapp-status-view"
]
  }
] as const;
test.each(roleMatrix)("complete guide access for $role", async row => {
 mocks.readSession.mockResolvedValue(session([...row.permissions]));
 const r=await get();expect(r.status).toBe(row.expected.length?200:403);
 if(row.expected.length){const data=helpResponseSchema.parse(await r.json());expect(data.articles.map(g=>g.id).sort()).toEqual([...row.expected]);}
 for(const guide of ALL_GUIDES){const exact=await get(`?article=${guide.id}&version=${guide.version}`);const included=(row.expected as readonly string[]).includes(guide.id);expect(exact.status).toBe(included?200:row.expected.length?404:403);}
});
````

### FILE: `src/components/journey-help-panel.tsx`
```yaml
block_id: "TS-FRANCHISE-JOURNEY-PORTALS:HELP:src-components-journey-help-panel.tsx:v1"
operation: CREATE
provenance: AUTHORED
source: "Local read-only integration of existing public operational guides and admitted session/runtime; see VERSIONED_JOURNEY_HELP_V379.md"
license: "LicenseRef-Workspace-Owner"
sha256: "ebc16333b017e210ca14593e674c6a481c04d517a4d05ac36857a8a7062b16c5"
variables: []
secrets_allowed: false
```
````typescript
"use client";
import {parseHelpQuery} from "@/platform/help/query";
import {guideDisplay} from "@/platform/i18n/private-guide-display";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

import { useCallback, useEffect, useRef, useState } from "react";
import { helpResponseSchema, type HelpResponse } from "@/platform/help/contract";

export function JourneyHelpPanel() {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

  const [query, setQuery] = useState("");
  const [result, setResult] = useState<HelpResponse | null>(null);
  const [status, setStatus] = useState(t("p0690"));
  const [busy, setBusy] = useState(false);
  const request = useRef<AbortController | null>(null);
  const currentQuery = useRef("");
  const clear = useCallback(() => { request.current?.abort(); request.current = null; setResult(null); setBusy(false); }, []);
  const load = useCallback(async (params: string) => {
    clear(); currentQuery.current = params;
    const controller = new AbortController(); request.current = controller;
    setBusy(true); setStatus(t("p0690"));
    try {
      const validated=parseHelpQuery(new URL("/help?"+params,window.location.origin).href);if(!validated)throw new Error("INVALID_HELP_QUERY");
      const parameters=new URLSearchParams(params),search=parameters.get("q")??"";parameters.delete("q");
      const response = await fetch(`/api/enterprise/help${parameters.size ? "?" + parameters.toString() : ""}`, { cache: "no-store", redirect: "error", signal: AbortSignal.any([controller.signal, AbortSignal.timeout(7000)]) });
      if (!response.ok) {
        const text = response.status === 401 ? t("p0691") : response.status === 403 ? t("p0692") : response.status === 404 ? t("p0693") : t("p0694");
        if (request.current === controller) setStatus(text);
        return;
      }
      const data = helpResponseSchema.parse(await response.json());
      data.articles=data.articles.filter(article=>{const view=guideDisplay(article,privateLocale.language);return [view.title,...view.paragraphs].join(" ").normalize("NFC").toLocaleLowerCase(privateLocale.locale).includes(search.trim().normalize("NFC").toLocaleLowerCase(privateLocale.locale))});
      if (request.current === controller && !controller.signal.aborted) { setResult(data); setStatus(data.articles.length ? t("p0695") : t("p0696")); }
    } catch { if (request.current === controller && !controller.signal.aborted) setStatus(t("p0694")); }
    finally { if (request.current === controller) { request.current = null; setBusy(false); } }
  }, [clear,privateLocale]);
  useEffect(() => {
    const initial = window.location.search.slice(1);
    setQuery(new URLSearchParams(initial).get("q")?.slice(0, 160) ?? "");
    void load(initial);
    const visible = () => { if (document.visibilityState === "hidden") { clear(); setStatus(t("p0697")); } else void load(currentQuery.current); };
    const focus = () => { if (document.visibilityState === "visible") void load(currentQuery.current); };
    document.addEventListener("visibilitychange", visible); window.addEventListener("focus", focus);
    return () => { request.current?.abort(); request.current = null; document.removeEventListener("visibilitychange", visible); window.removeEventListener("focus", focus); };
  }, [clear, load]);
  return <section aria-label={t("p0698")}>
    <h1>{t("p0699")}</h1>
    <p>{t("p0700")}</p>
    <form onSubmit={event => { event.preventDefault(); void load(new URLSearchParams({ q: query }).toString()); }}>
      <label>{t("p0701")}<input type="search" maxLength={160} value={query} onChange={event => { clear(); setQuery(event.target.value); currentQuery.current = new URLSearchParams({ q: event.target.value }).toString(); setStatus(t("p1193")); }}/></label>
      <button className="button" type="submit">{t("p0702")}</button>
    </form>
    <button className="button" type="button" disabled={busy} onClick={() => void load(currentQuery.current)}>{t("p0703")}</button>
    <p role="status" aria-live="polite">{status}</p>
    <div aria-busy={busy}>{result?.articles.map(source => {const article=guideDisplay(source,privateLocale.language);return <article lang={article.displayLanguage} className="card" key={article.id + "/" + article.version}>
      <h2>{article.title}</h2><p>{t("p0142")} {article.id}{t("p0280")}{article.version}</p>
      {article.paragraphs.map(text => <p key={text}>{text}</p>)}
      <a href={`/help?article=${encodeURIComponent(article.id)}&version=${encodeURIComponent(article.version)}`}>{t("p0704")}</a>
    </article>})}</div>
  </section>;
}
````

### FILE: `src/app/help/page.tsx`
```yaml
block_id: "TS-FRANCHISE-JOURNEY-PORTALS:HELP:src-app-help-page.tsx:v1"
operation: CREATE
provenance: AUTHORED
source: "Local read-only integration of existing public operational guides and admitted session/runtime; see VERSIONED_JOURNEY_HELP_V379.md"
license: "LicenseRef-Workspace-Owner"
sha256: "3a1644e0563b3741397b137a69dda461f51e3436abe5a5bcbd51221025fa9227"
variables: []
secrets_allowed: false
```
````typescript
import { JourneyHelpPanel } from "@/components/journey-help-panel";
export default function HelpPage() { return <JourneyHelpPanel/>; }
````

### FILE: `docs/journey-help.md`
```yaml
block_id: "TS-FRANCHISE-JOURNEY-PORTALS:HELP:docs-journey-help.md:v1"
operation: CREATE
provenance: AUTHORED
source: "Local read-only integration of existing public operational guides and admitted session/runtime; see VERSIONED_JOURNEY_HELP_V379.md"
license: "LicenseRef-Workspace-Owner"
sha256: "3022186c5e442ffebfd760d4dab5749e2376f80f856836d9ddeaf340f6c0f2f1"
variables: []
secrets_allowed: false
```
````markdown
# Ayuda operativa de los recorridos — V380

PUBLIC, AUTHORED, release-bound Spanish guidance. The same text/version powers the existing 15 existing versioned inline guides and /help. Exact id/version links never silently select another revision. GET /api/enterprise/help revalidates the encrypted BFF session each time, selects the existing page/section permissions documented in journey-help-coverage.md, and requires whatsapp_status_history only for notification guidance. It returns no tenant, organization, token, transaction identifier or personal data. These generic texts already ship in browser bundles; permission filtering selects applicable guidance and is not confidential-document access control.

Search is bounded literal NFC/case-insensitive substring matching over15fixed guides; no external query, regex, database or model. Unknown/repeated parameters, malformed version pairs and oversized queries fail. Unavailable article/version both return404. Responses use no-store. The UI aborts superseded reads, clears previous results before search/retry and when hidden, revalidates on focus, rejects malformed responses and offers explicit GET recovery. It does not infer external identity revocation before session expiry or a server refresh.

Only GET is implemented. No operation, message, acceptance, training enrollment or completion is triggered. No CMS/editor, draft/publication/archive workflow, historical version store, tenant documents, translated private help, progress tracking or general GO-HELP-CENTER-CORE equivalence is claimed. Changed instructions require a new source version plus reviewed tests and release evidence. The quote/notification versions remain1.0.0 because the exact prior text is reused. Future incompatible guide content must not reuse those versions.

Verification: unit/API authorization, feature, version, query and error tests; real production Next/JWE HTTPS loopback browser cases across Chromium desktop/mobile, Firefox and WebKit. Browser synthetic sessions do not prove an external IdP deployment. No database or provider is needed for these15static guides. Target business guidance, identity, support channels and operational acceptance remain consumer obligations.
````

### FILE: `src/components/operational-guide.tsx`
```yaml
block_id: "TS-FRANCHISE-JOURNEY-PORTALS:COVERAGE:src-components-operational-guide.tsx:v1"
operation: CREATE
provenance: AUTHORED
source: "Exact existing operational guide extraction; see VERSIONED_HELP_ROLE_COVERAGE_V380.md"
license: "LicenseRef-Workspace-Owner"
sha256: "37ffa11dbb90c30691d84c4862b79fc95e15d2d469475f6201d461b459f6b206"
variables: []
secrets_allowed: false
```
````typescript
"use client";
import {guideDisplay} from "@/platform/i18n/private-guide-display";
import {usePrivateI18n} from "@/platform/i18n/private-provider";
// Pure shared rendering: no effects, session state, operation or training progress.
type Guide = { id: string; version: string; title: string; summary: string; paragraphs: readonly string[]; inlineLead: boolean };
export function OperationalGuide({ guide, className }: { guide: Guide; className?: string }) {
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

  const display=guideDisplay(guide,privateLocale.language);
  return <details lang={display.displayLanguage} className={className}>
    <summary>{display.summary}</summary>
    <p>{guide.id}{t("p0280")}{guide.version}{guide.inlineLead ? " · " + display.paragraphs[0] : ""}</p>
    {display.paragraphs.slice(guide.inlineLead ? 1 : 0).map(text => <p key={text}>{text}</p>)}
    <a href={`/help?article=${encodeURIComponent(guide.id)}&version=${encodeURIComponent(guide.version)}`}>{t("p0143")}</a>
  </details>;
}
````

### FILE: `src/platform/help/coverage.test.ts`
```yaml
block_id: "TS-FRANCHISE-JOURNEY-PORTALS:COVERAGE:src-platform-help-coverage.test.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "Exact existing operational guide extraction; see VERSIONED_HELP_ROLE_COVERAGE_V380.md"
license: "LicenseRef-Workspace-Owner"
sha256: "3d5bd9aa236545599cca8b050017ebb0cd7125fa0107233cbdc40660d5f5a02d"
variables: []
secrets_allowed: false
```
````typescript
import { expect, test } from "vitest";
import { createElement } from "react";
import { renderToStaticMarkup } from "react-dom/server";
import { createHash } from "node:crypto";
import { OperationalGuide } from "@/components/operational-guide";
import { OPERATIONAL_GUIDES } from "./content";

// These digests are the exact V379 inline HTML, captured before extraction.
const original = [
  {
    "id": "return-operations-view",
    "sha256": "890335861532191b4fb532c11a0c1c07998eb388853cf767b35094411f0b7119",
    "className": null
  },
  {
    "id": "checklist-completion-view",
    "sha256": "2eb9cc5986181d4016397cef06f8373ce6dc9820543f8ecdb6e7a2aa681a9cb5",
    "className": null
  },
  {
    "id": "checklist-publication-view",
    "sha256": "a7487655d312417f4514358491fe0ef3d8b2972cb1113daa3cfaacf0af3f45f6",
    "className": null
  },
  {
    "id": "slot-create-view",
    "sha256": "db23a5340f56125be6842fa6d885995ad9bca2fae0b5c02191aa692711e43265",
    "className": null
  },
  {
    "id": "resource-create-view",
    "sha256": "2adbb54d47ad45c863bbf139fa81d39eb4780ba203a40e712dc00418252f4618",
    "className": null
  },
  {
    "id": "availability-create-view",
    "sha256": "3ea153fa5cc8cc8a9e36bad06ca38fa3581c62b7be72247fc826b0ade24041b8",
    "className": null
  },
  {
    "id": "quote-create-view",
    "sha256": "2361437eaac0a930954e0852779dfc1ac0cc60d4f7ad336a9fa48fcce7cbf93d",
    "className": null
  },
  {
    "id": "delivery-resolution-view",
    "sha256": "28126f161b2c40ab03836c14edefcd2530d9107bc2aa02d7b6f768bfcc85600a",
    "className": null
  },
  {
    "id": "lead-command-view",
    "sha256": "571367bf03b0d3ff229743a2ffde37c044edf4aabdd8bd05c9f927a74eec1dc8",
    "className": null
  },
  {
    "id": "availability-cancel-view",
    "sha256": "7b18e926cbcd8ab34b9a06fefdd5fa99e9162f964a3c22a9434d499e4bcbbd24",
    "className": null
  },
  {
    "id": "order-operations-view",
    "sha256": "d8954474a5b5f6d2f983ec4e755b5496f8617fea16cddc1dada7b27ce060b401",
    "className": null
  },
  {
    "id": "operation-sections-view",
    "sha256": "29fd0212dffe93893892c469c07a6a2a557b6a938d0d7fa2ce9569ddd241932b",
    "className": "card"
  },
  {
    "id": "handover-read-view",
    "sha256": "bcd57c885618d24d15922ee66ba495b2943c2b7aad506c5c7884b14d8a12737e",
    "className": null
  }
] as const;
test.each(original)("inline $id preserves all original paragraphs and version", row => {
 const html=renderToStaticMarkup(createElement(OperationalGuide,{guide:OPERATIONAL_GUIDES[row.id],...(row.className?{className:row.className}:{})}));
 const expectedLink=`<a href="/help?article=${row.id}&amp;version=${OPERATIONAL_GUIDES[row.id].version}">Abrir esta guía en Ayuda</a>`;
 expect(html).toContain(expectedLink);
 expect(createHash("sha256").update(html.replace(expectedLink,"")).digest("hex")).toBe(row.sha256);
});
````

### FILE: `docs/journey-help-coverage.md`
```yaml
block_id: "TS-FRANCHISE-JOURNEY-PORTALS:COVERAGE:docs-journey-help-coverage.md:v1"
operation: CREATE
provenance: AUTHORED
source: "Exact existing operational guide extraction; see VERSIONED_HELP_ROLE_COVERAGE_V380.md"
license: "LicenseRef-Workspace-Owner"
sha256: "ca9cacb671b871290458f19c7bdd39f17902575abdcbf0fcf35bf6b9a890dc7e"
variables: []
secrets_allowed: false
```
````markdown
# Versioned help coverage — V380

Scope:15/15versioned operational guides already present in the V379 selected web sources.
Thirteen were extracted without changing any instruction/version; two were already shared.
This is coverage of existing guidance, not a claim that every possible user task has a guide.
In particular factory:read alone has no existing guide; it receives403, as do empty permissions
and quote:write without lead:read. No factory policy, CMS or training curriculum is invented.

GET resolves the existing JWE session and applies the page/section predicates below. Help is
generic PUBLIC content, not a command grant or tenant-confidential document store. Lead readers
may read lead guidance without acquiring assign/update permissions. Availability readers may
read cancellation guidance without acquiring manage; manage alone can read creation guidance.
Quote creation guidance also needs lead:read because the actual page loads its leads under that
permission. The notification feature flag gates only the WhatsApp article, not other agenda help.

| Session permissions | Available article IDs |
|---|---|
| (none) | (none;403) |
| factory:read | (none;403) |
| customer:self | handover-read-view, quote-acceptance-view |
| lead:read | lead-command-view, operation-sections-view |
| quote:write | (none;403) |
| lead:read, quote:write | lead-command-view, operation-sections-view, quote-create-view |
| availability:read | availability-cancel-view, operation-sections-view |
| availability:manage | availability-create-view, operation-sections-view |
| resource:manage | operation-sections-view, resource-create-view |
| appointment:manage | operation-sections-view, slot-create-view, whatsapp-status-view |
| inventory:allocate | operation-sections-view, order-operations-view |
| payment:create | operation-sections-view, order-operations-view |
| handover:manage | checklist-completion-view, checklist-publication-view, delivery-resolution-view, operation-sections-view, order-operations-view, return-operations-view |
| admin:read | operation-sections-view, order-operations-view |
| * | availability-cancel-view, availability-create-view, checklist-completion-view, checklist-publication-view, delivery-resolution-view, handover-read-view, lead-command-view, operation-sections-view, order-operations-view, quote-acceptance-view, quote-create-view, resource-create-view, return-operations-view, slot-create-view, whatsapp-status-view |

## Exact source extraction

| Article / version | Previous inline source | Original JSX SHA-256 |
|---|---|---|
| availability-cancel-view/1.0.0 | src/components/franchise-command-panel.tsx | 7b18e926cbcd8ab34b9a06fefdd5fa99e9162f964a3c22a9434d499e4bcbbd24 |
| availability-create-view/1.0.0 | src/components/franchise-command-panel.tsx | 3ea153fa5cc8cc8a9e36bad06ca38fa3581c62b7be72247fc826b0ade24041b8 |
| checklist-completion-view/1.0.0 | src/components/franchise-command-panel.tsx | 2eb9cc5986181d4016397cef06f8373ce6dc9820543f8ecdb6e7a2aa681a9cb5 |
| checklist-publication-view/1.0.0 | src/components/franchise-command-panel.tsx | a7487655d312417f4514358491fe0ef3d8b2972cb1113daa3cfaacf0af3f45f6 |
| delivery-resolution-view/1.0.0 | src/components/franchise-command-panel.tsx | 28126f161b2c40ab03836c14edefcd2530d9107bc2aa02d7b6f768bfcc85600a |
| handover-read-view/1.0.0 | src/app/customer/handovers/page.tsx | bcd57c885618d24d15922ee66ba495b2943c2b7aad506c5c7884b14d8a12737e |
| lead-command-view/1.0.0 | src/components/franchise-command-panel.tsx | 571367bf03b0d3ff229743a2ffde37c044edf4aabdd8bd05c9f927a74eec1dc8 |
| operation-sections-view/1.0.0 | src/app/franchise/page.tsx | f4ee283d9423bef7179d8535271d4ba970331409bd5165ddd232e64aceeadf5c |
| order-operations-view/1.1.0 | src/components/franchise-command-panel.tsx | d8954474a5b5f6d2f983ec4e755b5496f8617fea16cddc1dada7b27ce060b401 |
| quote-create-view/1.0.0 | src/components/franchise-command-panel.tsx | 2361437eaac0a930954e0852779dfc1ac0cc60d4f7ad336a9fa48fcce7cbf93d |
| resource-create-view/1.0.0 | src/components/franchise-command-panel.tsx | 2adbb54d47ad45c863bbf139fa81d39eb4780ba203a40e712dc00418252f4618 |
| return-operations-view/1.0.0 | src/components/franchise-command-panel.tsx | 890335861532191b4fb532c11a0c1c07998eb388853cf767b35094411f0b7119 |
| slot-create-view/1.0.0 | src/components/franchise-command-panel.tsx | db23a5340f56125be6842fa6d885995ad9bca2fae0b5c02191aa692711e43265 |

The13render regression oracles compare the full original inline HTML, including summary,
paragraph order, version, punctuation and CSS class, after excluding only the newly added help
link. JSX className becomes HTML class in the one classed details block; its original source
digest above is deliberately retained as history. A separate frozen HTML digest governs that
render test. Semantic guide changes require a new guide version and an explicit before/after
review; do not refresh regression digests merely to conceal an instruction change.

The real HTTPS Next/JWE browser suite checks15permission profiles times15exact article links
in each of four browsers (900direct read decisions), plus all15UI version links in each browser
(60navigations), recovery and absence of browser writes. These figures count checks, not global
controls. Original business handlers outside the extracted help blocks remain byte-identical.
No external IdP, database, provider, confidential data or training activity is used by this gate.
````

### FILE: `docs/handover-browser-gate.md`

```yaml
block_id: "HANDOVER-DELTA-TYPESCRIPT_FRANCHISE_JOURNEY_PORTALS:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "eafa5a9f3494f72402820968e993f0997fd0a09bedea0c28c3a261c7c07b0ffe"
variables: []
secrets_allowed: false
```

````markdown
# Connected handover browser gate

This gate proves the admitted handover operator/customer workflow through an actual Chromium page, Next BFF, Go API and PostgreSQL. It uses only generated local identity credentials and the existing official Stripe SDK HTTP fixture. It does not perform hosted IdP login, live payment, physical shipment or a fiscal effect.

## Required composition

Materialize the connected payment, initial handover, commercial release, handover context/UI, identity session, Next and Playwright packs together. The Go test reuses `newConnectedRun` from `payment_connected_integration_test.go`; its database guard requires an isolated loopback database whose name starts with `elite_payment_connected_`. Apply the composition's migrations before running the test. Do not point the fixture at business data.

Use the admitted Node, Go, Playwright and Chromium versions recorded in the accompanying runtime/evidence locks. Existing dependency caches are sufficient; this gate has no dependency installer. Build the frontend with its standard Next CLI and all TypeScript checks enabled. For a node_modules junction outside the workspace, the checked configuration uses `next build --webpack`; Turbopack rejects that junction. This is a builder choice, not a type-check exception.

Set `ELITE_HANDOVER_BROWSER=1`, `ELITE_WEB_ROOT` to the materialized frontend root, `ELITE_NODE_BIN` to the absolute admitted Node executable, and `PAYMENT_CONNECTED_DB_URL` to the isolated fixture database. The reference run also uses `BUSINESS_CONFIG_FILE=business.example.json`, `NEXT_TELEMETRY_DISABLED=1`, `GOTOOLCHAIN=local`, `GOPROXY=off` and `GOSUMDB=off`.

Run the admitted Go executable directly:

```text
go test -mod=readonly ./internal/platform/postgres -run ^TestHandoverBrowserPostgres$ -count=1 -v -timeout=4m
```

The Go test starts the already-built Next server through the explicit Node executable, creates a loopback HTTPS fixture and invokes the existing Playwright CLI directly for `tests/handover-connected.spec.mjs`, project `chromium-desktop`. It supplies the base URL so Playwright's package-manager webServer path is not used. An omitted fixture flag/database causes an explicit skip; a release gate must reject skipped execution.

## Assertions and effects

- Real encrypted role sessions and RS256/JWKS bearer verification authorize operator/customer identities; foreign organization and unrelated customer attempts fail.
- The payment fixture executes official SDK calls and signed callback reconciliation before preparation. The policy is an explicitly activated, hash-bound materialized profile with exact tenant, organization and provider/account/connection bindings.
- The browser prepares a handover, loses the response after commit, recovers it by GET, completes the existing checklist and obtains the real customer's acceptance.
- The browser records a commercial release receipt, loses the response after commit, then recovers the immutable receipt after a signed refund callback has already placed the payment on hold.
- Current authorization is false both before and after the official GET refund reconciliation. Recovery never converts historical receipt existence into current authorization.
- Double-click preparation produces one POST. PostgreSQL contains one preparation, one acceptance, one commercial receipt and one release event. There are two customer acceptance HTTP requests by design: one forbidden stranger attempt followed by one successful customer request.
- The final run requires zero React page errors and no horizontal overflow at a 390 by 844 mobile viewport. This is one Chromium project with a mobile viewport check, not a four-browser matrix.
- The commercial receipt creates no stock or money journal mutation and does not claim physical dispatch.

The local self-signed HTTPS fixture is allowed only by the existing explicit loopback Playwright configuration. Production TLS verification is not changed. Fixture control endpoints exist only in Go test code and use a generated per-run token.

## Source classification and delta

All new composition code, tests and frontend glue are AUTHORED. The timestamp fix reuses the admitted server locale/timezone formatter; the BFF bounded reader reuses the existing surveys reader with a 4096-byte limit. New PageProps overloads satisfy Next's generated type contract while preserving zero-argument direct-call tests. No upstream business logic, dependency or license is added or attributed to a vendor by these changes.

The evidence manifest records source prehashes, final source hashes, build/test/browser receipts and retained failure logs. Runtime directories, node_modules, .next, videos and screenshots are evidence or generated artifacts; they are not product source files.
````

### FILE: `docs/handover-operator-flow.md`

```yaml
block_id: "HANDOVER-DELTA-TYPESCRIPT_FRANCHISE_JOURNEY_PORTALS:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "105da679f23d81aca759441855a2a24c68f86f6c9bb30113310d1daebc00d52f"
variables: []
secrets_allowed: false
```

````markdown
# Operator handover flow

The existing `/franchise` page now includes an order-bound handover panel for `handover:manage`. It reads an authoritative scoped context, prepares the initial handover, embeds the existing checklist completion owner with the server's handover ID/version, and records or recovers the commercial receipt after the customer's acceptance. The existing customer route `/customer/handovers` and `CustomerHandoverActions` continue to own acceptance, observed serial and discrepancy reporting. No operator can manufacture customer acceptance from this new panel.

## Host binding

Mount `httpapi.HandoverContextModule{Service: service}` beside `InitialHandoverModule` and `CommercialReleaseModule`, with the same admitted `*franchisejourney.HandoverPreparationService`. If using the existing helper, its `InitialHandoverModule.Service` concrete value implements `HandoverContextService`. No additional runtime credential, policy switch, SDK, database or dependency is introduced. The active profile remains hash-bound to tenant, organization, provider/account/connection and mode.

The single new Go route is `GET /v1/franchise/orders/{id}/handover-context?organization_id=...`. Verified claims and `handover:manage` membership authorize the read. SQL derives the allocated line, captured attempt and observed evidence, then invokes `lockInitialHandoverScope`. A hold, foreign scope, profile mismatch or incompatible order is not presented as ready. The result includes the current initial handover/version, if one exists. This read never replaces the write's own locked revalidation.

## BFF and browser behavior

`/api/enterprise/handovers` accepts only the supported actions `prepare` and `release`, an order/organization selection and a persisted request key. It obtains payment bindings server-to-server; browser-entered hashes, payment attempts, customer IDs, amounts and stock IDs are rejected. The context returned to the browser strips the internal observation digest, attempt ID and any unknown backend fields.

GET operations are `context`, `prepare-result`, `release-result` and `release-current`. Recovery uses the existing protected backend client with an optional validated `Idempotency-Key` header. The customer/operator access token stays in the server. Responses disable caching. Authentication, organization, duplicate/unknown query parameters, path identifiers and response bindings are checked before returning data.

The browser stores only an action, order ID, request UUID and (for release) handover ID in session storage scoped by the server-derived tenant/actor/organization digest. A storage failure blocks the mutation. An uncertain outcome retains its key and exposes GET recovery, without automatic POST replay or a new key. Keeping the handover ID allows historical recovery even when a later payment callback makes the current context unavailable. A receipt response contains no current authorization flag; current validity is an explicit separate query and is labeled with its evaluation time and expiry.

The existing checklist publisher and recovery controls remain available. The same checklist component supports optional read-only handover ID/version prefill for this per-order flow; its original independent workflow is preserved. No new checklist algorithm, pricing, shipment, tax or financial posting logic is implemented.

## Local evidence and provenance

All deltas are `AUTHORED` transport, persistence/read, presentation and test glue over the admitted owners. The exact existing Go, Next, React, Zod, TypeScript and Vitest dependency locks remain unchanged. No third-party code was copied, and no new dependency, license or notice was added.

Focal checks cover HTTP authorization, real PostgreSQL scope/profile/hold checks, BFF command derivation and recovery, cross-scope and forged fields, retained keys/storage failure, historical versus current receipts, existing checklist/customer acceptance BFF contracts, server rendering, protected-client recovery headers, type generation and TypeScript compilation. Provider reconciliation and the durable commercial write retain their separate already completed SDK/PostgreSQL evidence. No live account or production deployment is claimed.
````

### FILE: `microsoft_playwright_browser_gate/tests/handover-connected.spec.mjs`

```yaml
block_id: "HANDOVER-DELTA-TYPESCRIPT_FRANCHISE_JOURNEY_PORTALS:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "e5d8281a373df5acbe2bac45e687261990538e1a6979a94dee1ceb6641b905e3"
variables: []
secrets_allowed: false
```

````text
import {test,expect} from '@playwright/test';
import {createHash} from 'node:crypto';
import {createRequire} from 'node:module';
import {resolve} from 'node:path';
import {pathToFileURL} from 'node:url';

test('operator and customer close handover through real BFF API and PostgreSQL',async({page,context},info)=>{
 if(process.env.ELITE_HANDOVER_BROWSER!=='1')throw new Error('explicit local handover fixture required');
 const base=process.env.ELITE_BASE_URL;expect(new URL(base).protocol).toBe('https:');expect(new URL(base).hostname).toBe('127.0.0.1');
 const identities=JSON.parse(process.env.ELITE_HANDOVER_IDENTITIES);
 const require=createRequire(resolve(process.env.ELITE_WEB_ROOT,'package.json'));
 const {EncryptJWT}=await import(pathToFileURL(require.resolve('jose')).href);
 async function identity(name){
  const jwt=await new EncryptJWT(identities[name]).setProtectedHeader({alg:'dir',enc:'A256GCM',typ:'JWT'}).setIssuedAt().setExpirationTime('300s').encrypt(createHash('sha256').update(process.env.AUTH_SESSION_SECRET).digest());
  await context.clearCookies();await context.addCookies([{name:'__Host-elite_session',value:jwt,url:base+'/',secure:true,httpOnly:true,sameSite:'Lax'}]);
 }
 const errors=[];page.on('pageerror',error=>errors.push(error.message));
 const writes=[];page.on('request',request=>{if(request.method()==='POST'&&request.url().includes('/api/enterprise/'))writes.push(request.postDataJSON()?.action)});
 await identity('operator');await page.goto('/franchise');
 const panel=page.getByRole('article',{name:'Entrega del pedido order',exact:true});
 await expect(panel.getByRole('button',{name:'Preparar entrega',exact:true})).toBeEnabled();
 // Discard only the browser response after the actual Go/PG transaction.
 await page.route('**/api/enterprise/handovers',async route=>{
  const action=route.request().postDataJSON()?.action;
  if(route.request().method()!=='POST'||action!=='prepare'){await route.continue();return}
  const response=await route.fetch();expect(response.status()).toBe(201);await route.abort('failed');
 });
 await panel.getByRole('button',{name:'Preparar entrega',exact:true}).evaluate(button=>{button.click();button.click()});
 await expect(panel.getByText('El resultado quedó sin comprobar.',{exact:false})).toBeVisible();
 expect(writes.filter(x=>x==='prepare')).toHaveLength(1);
 await panel.getByRole('button',{name:'Consultar resultado de preparación',exact:true}).click();
 await expect(panel.getByText('Preparación recuperada:',{exact:false})).toBeVisible();
 await page.unrouteAll();await page.reload();
 const presentation=panel.getByRole('region',{name:'Presentación de entrega',exact:true});
 await expect(presentation.getByLabel('Entrega',{exact:true})).toHaveAttribute('readonly','');
 const handoverID=await presentation.getByLabel('Entrega',{exact:true}).inputValue();
 await expect(presentation.getByLabel('Versión actual de entrega')).toHaveValue('1');
 await presentation.getByLabel('ID del checklist',{exact:true}).fill('browser-checklist');
 await presentation.getByLabel('Versión del checklist',{exact:true}).fill('1');
 await presentation.getByLabel('Respuestas, una por línea').fill('serial=SERIAL-SYNTHETIC');
 await presentation.getByRole('button',{name:'Completar y presentar',exact:true}).click();
 await expect(presentation.getByText('Respuesta recibida.',{exact:false})).toBeVisible();
 await presentation.getByRole('button',{name:'Consultar presentación',exact:true}).click();
 await expect(presentation.getByText('Presentación recuperada:',{exact:false})).toBeVisible();
 await identity('stranger');await page.goto('/customer/handovers');
 await expect(page.getByRole('button',{name:'Registrar recepción',exact:true})).toHaveCount(0);
 const unauthorized=await context.request.post(base+'/api/enterprise/franchise/commands',{headers:{origin:base,'content-type':'application/json'},data:{action:'accept-handover',organizationId:'store',handoverId:handoverID,version:2,confirmedReceived:true,serialNumber:'SERIAL-SYNTHETIC',checklistId:'browser-checklist',checklistVersion:1}});
 expect([404,409]).toContain(unauthorized.status());
 await identity('customer');await page.goto('/customer/handovers');
 expect((await context.request.get(base+'/api/enterprise/handovers?kind=context&organizationId=store&orderId=order')).status()).toBe(403);
 await page.getByLabel('Número de serie observado').fill('SERIAL-SYNTHETIC');
 await page.getByRole('checkbox',{name:'Confirmo que recibí el activo identificado y revisé la preparación indicada'}).check();
 await page.getByRole('button',{name:'Registrar recepción',exact:true}).click();
 await expect(page.getByText('Aceptado:',{exact:false})).toBeVisible();
 await expect(page.getByRole('button',{name:'Registrar recepción',exact:true})).toHaveCount(0);
 await identity('operator');await page.goto('/franchise');
 await expect(panel.getByRole('button',{name:'Registrar cierre comercial',exact:true})).toBeEnabled();
 await page.route('**/api/enterprise/handovers',async route=>{
  if(route.request().method()!=='POST'||route.request().postDataJSON()?.action!=='release'){await route.continue();return}
  const response=await route.fetch();expect(response.status()).toBe(201);await route.abort('failed');
 });
 await panel.getByRole('button',{name:'Registrar cierre comercial',exact:true}).click();
 await expect(panel.getByText('El resultado quedó sin comprobar.',{exact:false})).toBeVisible();
 await page.unrouteAll();
 // The signed provider callback holds the observation before official GET.
 const control=await context.request.post(process.env.ELITE_HANDOVER_CONTROL,{headers:{'X-Fixture-Token':process.env.ELITE_HANDOVER_CONTROL_TOKEN},data:{action:'hold'}});expect(control.status()).toBe(200);
 await page.reload();
 await expect(panel.getByRole('button',{name:'Registrar cierre comercial',exact:true})).toHaveCount(0);
 await panel.getByRole('button',{name:'Consultar resultado comercial',exact:true}).click();
 await expect(panel.getByText('Recibo histórico recuperado.',{exact:false})).toBeVisible();
 await panel.getByRole('button',{name:'Consultar vigencia del recibo',exact:true}).click();
 await expect(panel.getByText('El recibo no tiene vigencia comprobada',{exact:false})).toBeVisible();
 const reconciled=await context.request.post(process.env.ELITE_HANDOVER_CONTROL,{headers:{'X-Fixture-Token':process.env.ELITE_HANDOVER_CONTROL_TOKEN},data:{action:'reconcile'}});expect(reconciled.status()).toBe(200);
 await panel.getByRole('button',{name:'Consultar vigencia del recibo',exact:true}).click();
 await expect(panel.getByText('El recibo no tiene vigencia comprobada',{exact:false})).toBeVisible();
 expect(writes.filter(x=>x==='prepare')).toHaveLength(1);expect(writes.filter(x=>x==='release')).toHaveLength(1);expect(writes.filter(x=>x==='complete-delivery-checklist')).toHaveLength(1);expect(writes.filter(x=>x==='accept-handover')).toHaveLength(1);
 await page.setViewportSize({width:390,height:844});await expect(panel.getByRole('button',{name:'Consultar resultado comercial'})).toBeVisible();
 expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth+1)).toBe(true);
 await page.screenshot({path:info.outputPath('handover-recovered-mobile.png'),fullPage:true});
 await identity('foreign-org');expect((await context.request.get(base+'/api/enterprise/handovers?kind=context&organizationId=store&orderId=order')).status()).toBe(403);
 expect(errors).toEqual([]);
 console.log('HANDOVER_BROWSER_PASS prepare=1 checklist=1 customer_accept=1 commercial_release=1 uncertain_response_recovery=true callback_hold_current_false=true role_and_object_scope_denied=true');
});
````

### FILE: `src/app/api/enterprise/handovers/bounded-body.test.ts`

```yaml
block_id: "HANDOVER-DELTA-TYPESCRIPT_FRANCHISE_JOURNEY_PORTALS:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "f1593b49a45932a473ad422d1a72517d156d447f8de4af694450f2a6a4a70af8"
variables: []
secrets_allowed: false
```

````typescript
import{expect,it,vi}from"vitest";
const m=vi.hoisted(()=>({get:vi.fn(),post:vi.fn()}));
vi.mock("@/platform/auth/session",()=>({readSession:async()=>({subject:"operator",tenantId:"tenant",organizations:["store"],permissions:["handover:manage"]}),allowed:()=>true}));
vi.mock("@/platform/auth/oidc-client",()=>({applicationBaseUrl:()=>new URL("https://portal.example.test")}));
vi.mock("@/platform/backend/protected-client",()=>({protectedGet:m.get,protectedPost:m.post}));
import{POST}from"./route";
it("cuts off an oversized chunked UTF-8 request without Content-Length before downstream",async()=>{
 let pulls=0,cancelled=false;
 const body=new ReadableStream<Uint8Array>({pull(controller){pulls++;controller.enqueue(new TextEncoder().encode("é".repeat(1025)))},cancel(){cancelled=true}});
 const r=new Request("https://portal.example.test/api/enterprise/handovers",{method:"POST",headers:{origin:"https://portal.example.test","content-type":"application/json"},body,duplex:"half"}as RequestInit);
 expect(r.headers.get("content-length")).toBeNull();const reply=await POST(r);
 expect(reply.status).toBe(413);expect(await reply.json()).toEqual({code:"COMMAND_TOO_LARGE"});expect(cancelled).toBe(true);expect(pulls).toBeLessThanOrEqual(3);expect(m.get).not.toHaveBeenCalled();expect(m.post).not.toHaveBeenCalled();
});
````

### FILE: `src/app/api/enterprise/handovers/route.test.ts`

```yaml
block_id: "HANDOVER-DELTA-TYPESCRIPT_FRANCHISE_JOURNEY_PORTALS:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "e3e375d6b67aa8212c9caa9d4768dedc6c50fc9a71a11f1d5e89f57ac64b4103"
variables: []
secrets_allowed: false
```

````typescript
import {beforeEach,expect,it,vi}from"vitest";
const m=vi.hoisted(()=>({session:vi.fn(),get:vi.fn(),post:vi.fn()}));
vi.mock("@/platform/auth/session",()=>({readSession:m.session,allowed:(s:{permissions:string[]},p:string)=>s.permissions.includes(p)}));
vi.mock("@/platform/auth/oidc-client",()=>({applicationBaseUrl:()=>new URL("https://portal.example.test")}));
vi.mock("@/platform/backend/protected-client",()=>({protectedGet:m.get,protectedPost:m.post}));
import{GET,POST}from"./route";
const h="a".repeat(64),date="2026-09-11T12:00:00Z",key="a03c7850-13ad-4b38-81d0-120000000001";
const context=()=>({organization_id:"store",order_id:"order",order_line_id:"line",payment_attempt_id:"payment",observation_sha256:h,can_prepare:true,release_effect:"COMMIT_COMMERCIAL_RELEASE_RECEIPT",evaluated_at:date});
const handover=()=>({id:"handover",organization_id:"store",order_id:"order",state:"prepared",version:1});
const prepared=()=>({handover:handover(),order_line_id:"line",reservation_id:"reservation",payment_attempt_id:"payment",observation_sha256:h,contract_id:"profile@1",contract_sha256:h,prepared_by:"operator",prepared_at:date});
const receipt=()=>({id:"receipt",organization_id:"store",handover_id:"handover",order_id:"order",payment_attempt_id:"payment",observation_sha256:h,observation_generation:1,handover_version:3,acceptance_sha256:h,checklist_id:"checklist",checklist_version:1,contract_id:"profile@1",contract_sha256:h,effect:"COMMIT_COMMERCIAL_RELEASE_RECEIPT",released_by:"operator",recorded_at:date,valid_until:"2026-09-11T12:05:00Z"});
const session={subject:"operator",tenantId:"tenant",organizations:["store"],permissions:["handover:manage"],accessToken:"never-browser-token"};
const url="https://portal.example.test/api/enterprise/handovers";
const post=(extra:Record<string,unknown>={})=>new Request(url,{method:"POST",headers:{origin:"https://portal.example.test","content-type":"application/json"},body:JSON.stringify({action:"prepare",organizationId:"store",orderId:"order",requestKey:key,...extra})});
beforeEach(()=>{vi.clearAllMocks();m.session.mockResolvedValue(session);m.get.mockResolvedValue(context());m.post.mockResolvedValue(prepared())});
it("derives preparation bindings server-side and passes the original key",async()=>{const r=await POST(post());expect(r.status).toBe(201);expect(m.post).toHaveBeenCalledWith(session,"/v1/franchise/orders/order/handover",{organization_id:"store",order_line_id:"line",payment_attempt_id:"payment",observation_sha256:h},key);expect(await r.json()).toEqual({handover:handover()})});
it.each(["observationSha256","customerId","amountMinor","stockUnitId","paymentAttemptId"])("rejects browser financial or identity field %s",async field=>{expect((await POST(post({[field]:"forged"}))).status).toBe(400);expect(m.get).not.toHaveBeenCalled();expect(m.post).not.toHaveBeenCalled()});
it.each(["unauthenticated","permission","organization","origin"])("rejects %s before any backend effect",async mode=>{let req=post();if(mode==="unauthenticated")m.session.mockResolvedValue(null);if(mode==="permission")m.session.mockResolvedValue({...session,permissions:[]});if(mode==="organization")req=post({organizationId:"foreign"});if(mode==="origin")req=new Request(req,{headers:{origin:"https://evil.test","content-type":"application/json"}});expect([401,403]).toContain((await POST(req)).status);expect(m.get).not.toHaveBeenCalled();expect(m.post).not.toHaveBeenCalled()});
it("hides provider observation details in the browser context",async()=>{m.get.mockResolvedValue({...context(),client_secret:"not-allowed"});const r=await GET(new Request(url+"?kind=context&organizationId=store&orderId=order"));expect(r.status).toBe(200);const raw=await r.text();expect(raw).not.toContain(h);expect(raw).not.toContain("client_secret");expect(raw).not.toContain("never-browser-token")});
it("recovers preparation via GET with exact idempotency header and no POST",async()=>{m.get.mockResolvedValue(prepared());const r=await GET(new Request(url+"?kind=prepare-result&organizationId=store&orderId=order&requestKey="+key));expect(r.status).toBe(200);expect(m.get).toHaveBeenCalledWith(session,"/v1/franchise/orders/order/handover-result",{organization_id:"store"},key);expect(m.post).not.toHaveBeenCalled()});
it("commits only the accepted handover under explicit commercial profile",async()=>{m.get.mockResolvedValue({...context(),can_prepare:false,handover:{...handover(),state:"accepted",version:3}});m.post.mockResolvedValue(receipt());const r=await POST(post({action:"release"}));expect(r.status).toBe(201);expect(m.post).toHaveBeenCalledWith(session,"/v1/franchise/handovers/handover/commercial-release",{organization_id:"store",observation_sha256:h},key);expect(await r.json()).not.toHaveProperty("current")});
it.each(["prepared","readonly"])("rejects release for %s",async kind=>{m.get.mockResolvedValue({...context(),can_prepare:false,handover:{...handover(),state:kind==="prepared"?"prepared":"accepted"},release_effect:kind==="readonly"?"READ_ONLY_ELIGIBILITY":"COMMIT_COMMERCIAL_RELEASE_RECEIPT"});expect((await POST(post({action:"release"}))).status).toBe(409);expect(m.post).not.toHaveBeenCalled()});
it("recovers historical receipt independently of a held/missing current context",async()=>{m.get.mockResolvedValue(receipt());const r=await GET(new Request(url+"?kind=release-result&organizationId=store&orderId=order&handoverId=handover&requestKey="+key));expect(r.status).toBe(200);expect(await r.json()).toEqual({receipt:receipt()});expect(m.get).toHaveBeenCalledTimes(1);expect(m.post).not.toHaveBeenCalled()});
it("shows current=false after callback without deleting historical receipt",async()=>{m.get.mockResolvedValue({receipt:receipt(),evaluated_at:date,current:false});const r=await GET(new Request(url+"?kind=release-current&organizationId=store&orderId=order&handoverId=handover"));expect(r.status).toBe(200);expect((await r.json()).current).toBe(false)});
it.each(["organizationId=store&organizationId=store","organizationId=store&unexpected=value","organizationId=foreign"])("denies invalid recovery query %s",async org=>{expect([400,403]).toContain((await GET(new Request(url+"?kind=context&orderId=order&"+org))).status);expect(m.get).not.toHaveBeenCalled()});
it("rejects receipt mismatch and never retries or exposes backend errors",async()=>{m.post.mockResolvedValue({...prepared(),handover:{...handover(),order_id:"other"}});expect((await POST(post())).status).toBe(502);expect(m.post).toHaveBeenCalledTimes(1);m.get.mockRejectedValue(new Error("PRIVATE_BACKEND_DETAIL"));const r=await GET(new Request(url+"?kind=context&organizationId=store&orderId=order"));expect(r.status).toBe(503);expect(await r.text()).not.toContain("PRIVATE_BACKEND_DETAIL")});

it("uses exactly one stored-value funding reference for preparation",async()=>{
 m.get.mockResolvedValue({...context(),payment_attempt_id:"",funding_receipt_id:"funding_local"});m.post.mockResolvedValue({...prepared(),payment_attempt_id:"",funding_receipt_id:"funding_local"});
 expect((await POST(post())).status).toBe(201);expect(m.post).toHaveBeenCalledWith(session,"/v1/franchise/orders/order/handover",{organization_id:"store",order_line_id:"line",payment_attempt_id:"",funding_receipt_id:"funding_local",observation_sha256:h},key);
});
it.each(["both","neither"])("rejects %s funding sources before preparation",async mode=>{
 m.get.mockResolvedValue({...context(),payment_attempt_id:mode==="both"?"payment":"",...(mode==="both"?{funding_receipt_id:"funding_local"}:{})});expect((await POST(post())).status).toBe(503);expect(m.post).not.toHaveBeenCalled();
});
````

### FILE: `src/app/api/enterprise/handovers/route.ts`

```yaml
block_id: "HANDOVER-DELTA-TYPESCRIPT_FRANCHISE_JOURNEY_PORTALS:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "25639c5b727f5543d8d35b03950efd344f539d35d5bd435a26e89f8313734955"
variables: []
secrets_allowed: false
```

````typescript
import { boundedCommandBody as boundedBody } from "@/platform/backend/bounded-command";
// AUTHORED scoped BFF; obtains observation bindings server-to-server and forwards
// only supported commands. No browser-entered amount, customer or digest.
import { NextResponse } from "next/server";
import { z } from "zod";
import { allowed,readSession } from "@/platform/auth/session";
import { applicationBaseUrl } from "@/platform/auth/oidc-client";
import { protectedGet,protectedPost } from "@/platform/backend/protected-client";
import { BackendProblem } from "@/platform/backend/public-client";
import { backendHandoverContextSchema,handoverContextSchema,handoverID,preparationSchema,commercialReceiptSchema,currentReleaseSchema } from "@/platform/handovers/contracts";
const reply=(body:unknown,status=200)=>NextResponse.json(body,{status,headers:{"Cache-Control":"no-store","X-Content-Type-Options":"nosniff"}});
const base={organizationId:handoverID,orderId:handoverID};
const command=z.object({...base,action:z.enum(["prepare","release"]),requestKey:z.uuid()}).strict();
const query=z.discriminatedUnion("kind",[
 z.object({...base,kind:z.literal("context")}).strict(),
 z.object({...base,kind:z.literal("prepare-result"),requestKey:z.uuid()}).strict(),
 z.object({...base,kind:z.literal("release-result"),handoverId:handoverID,requestKey:z.uuid()}).strict(),
 z.object({...base,kind:z.literal("release-current"),handoverId:handoverID}).strict(),
]);
function problem(error:unknown){return reply({code:"HANDOVER_UNAVAILABLE"},error instanceof BackendProblem&&[400,401,403,404,409,503].includes(error.status)?error.status:503);}
export async function POST(request:Request){
 try{if(request.headers.get("origin")!==applicationBaseUrl().origin)return reply({code:"CROSS_ORIGIN_REJECTED"},403)}catch{return reply({code:"CROSS_ORIGIN_REJECTED"},403)}
 if(request.headers.get("content-type")!=="application/json")return reply({code:"UNSUPPORTED_MEDIA_TYPE"},415);
 const session=await readSession();if(!session)return reply({code:"UNAUTHENTICATED"},401);if(!allowed(session,"handover:manage"))return reply({code:"FORBIDDEN"},403);
 let body:unknown;try{body=JSON.parse(await boundedBody(request))}catch(error){return error instanceof Error&&error.message==="BODY_TOO_LARGE"?reply({code:"COMMAND_TOO_LARGE"},413):reply({code:"INVALID_COMMAND"},400)}
 const parsed=command.safeParse(body);if(!parsed.success)return reply({code:"INVALID_COMMAND"},400);const c=parsed.data;
 if(!session.organizations.includes(c.organizationId))return reply({code:"FORBIDDEN"},403);
 try{
  const v=backendHandoverContextSchema.parse(await protectedGet<unknown>(session,`/v1/franchise/orders/${encodeURIComponent(c.orderId)}/handover-context`,{organization_id:c.organizationId}));
  if(v.organization_id!==c.organizationId||v.order_id!==c.orderId)return reply({code:"CONTEXT_MISMATCH"},502);
  if(c.action==="prepare"){
   if(!v.can_prepare)return reply({code:"HANDOVER_ALREADY_PREPARED"},409);
   const result=preparationSchema.parse(await protectedPost<unknown>(session,`/v1/franchise/orders/${encodeURIComponent(c.orderId)}/handover`,{organization_id:c.organizationId,order_line_id:v.order_line_id,payment_attempt_id:v.payment_attempt_id,...(v.funding_receipt_id?{funding_receipt_id:v.funding_receipt_id}:{}),observation_sha256:v.observation_sha256},c.requestKey));
   if(result.handover.organization_id!==c.organizationId||result.handover.order_id!==c.orderId||result.order_line_id!==v.order_line_id||result.payment_attempt_id!==v.payment_attempt_id||result.funding_receipt_id!==v.funding_receipt_id||result.observation_sha256!==v.observation_sha256)return reply({code:"RECEIPT_MISMATCH"},502);
   return reply({handover:result.handover},201);
  }
  if(v.release_effect!=="COMMIT_COMMERCIAL_RELEASE_RECEIPT"||v.handover?.state!=="accepted")return reply({code:"RELEASE_NOT_ENABLED"},409);
  const receipt=commercialReceiptSchema.parse(await protectedPost<unknown>(session,`/v1/franchise/handovers/${encodeURIComponent(v.handover.id)}/commercial-release`,{organization_id:c.organizationId,observation_sha256:v.observation_sha256},c.requestKey));
  if(receipt.organization_id!==c.organizationId||receipt.order_id!==c.orderId||receipt.handover_id!==v.handover.id||receipt.observation_sha256!==v.observation_sha256)return reply({code:"RECEIPT_MISMATCH"},502);
  return reply({receipt},201);
 }catch(error){return problem(error)}
}
export async function GET(request:Request){
 if(["cross-site","none"].includes(request.headers.get("sec-fetch-site")??""))return reply({code:"CROSS_ORIGIN_REJECTED"},403);
 const session=await readSession();if(!session)return reply({code:"UNAUTHENTICATED"},401);if(!allowed(session,"handover:manage"))return reply({code:"FORBIDDEN"},403);
 const search=new URL(request.url).searchParams;if(request.url.length>2048||[...search.keys()].some(k=>search.getAll(k).length!==1))return reply({code:"INVALID_QUERY"},400);
 const parsed=query.safeParse(Object.fromEntries(search));if(!parsed.success)return reply({code:"INVALID_QUERY"},400);const q=parsed.data;
 if(!session.organizations.includes(q.organizationId))return reply({code:"FORBIDDEN"},403);
 try{
  if(q.kind==="context"){
   const value=handoverContextSchema.parse(await protectedGet<unknown>(session,`/v1/franchise/orders/${encodeURIComponent(q.orderId)}/handover-context`,{organization_id:q.organizationId}));
   if(value.organization_id!==q.organizationId||value.order_id!==q.orderId)return reply({code:"CONTEXT_MISMATCH"},502);return reply(value);
  }
  if(q.kind==="prepare-result"){
   const value=preparationSchema.parse(await protectedGet<unknown>(session,`/v1/franchise/orders/${encodeURIComponent(q.orderId)}/handover-result`,{organization_id:q.organizationId},q.requestKey));
   if(value.handover.organization_id!==q.organizationId||value.handover.order_id!==q.orderId)return reply({code:"RECEIPT_MISMATCH"},502);return reply({handover:value.handover});
  }
  const raw=await protectedGet<unknown>(session,`/v1/franchise/handovers/${encodeURIComponent(q.handoverId)}/${q.kind==="release-result"?"commercial-release-result":"commercial-release-current"}`,{organization_id:q.organizationId},q.kind==="release-result"?q.requestKey:undefined);
  const value=q.kind==="release-result"?{receipt:commercialReceiptSchema.parse(raw)}:currentReleaseSchema.parse(raw);
  if(value.receipt.organization_id!==q.organizationId||value.receipt.order_id!==q.orderId||value.receipt.handover_id!==q.handoverId)return reply({code:"RECEIPT_MISMATCH"},502);return reply(value);
 }catch(error){return problem(error)}
}
````

### FILE: `src/components/handover-operations-panel.tsx`

```yaml
block_id: "HANDOVER-DELTA-TYPESCRIPT_FRANCHISE_JOURNEY_PORTALS:file7:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "d9790e1c34b165dccab6f498d01eb0f613527df64be3d9a2f772bc25fdd77514"
variables: []
secrets_allowed: false
```

````tsx
"use client";
import {usePrivateI18n} from "@/platform/i18n/private-provider";


// AUTHORED orchestration UI; uses existing checklist and customer acceptance.
import { useEffect,useRef,useState } from "react";
import { ChecklistCompletionPanel } from "@/components/franchise-command-panel";
import { commercialReceiptSchema,currentReleaseSchema,handoverContextSchema,handoverViewSchema,operationMarkerSchema,retainOperation,type CommercialReceipt,type CurrentRelease,type HandoverContext,type OperationMarker } from "@/platform/handovers/contracts";

const endpoint="/api/enterprise/handovers";
export function HandoverOperationsPanel({orders,organization,scope}:{orders:{id:string}[]|null;organization:string;scope:string}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

 return <section className="card" aria-label={t("p0617")}><h2>{t("p0618")}</h2><p>{t("p0619")}</p>
 {orders===null?<p role="alert">{t("p0620")}</p>:orders.length===0?<p>{t("p0621")}</p>:orders.map(order=><HandoverOrder key={`${scope}:${order.id}`} orderId={order.id} organization={organization} scope={scope}/>)}
 </section>;
}

function HandoverOrder({orderId,organization,scope}:{orderId:string;organization:string;scope:string}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();

 const [context,setContext]=useState<HandoverContext|null>(null),[message,setMessage]=useState(""),[ready,setReady]=useState(false),[busy,setBusy]=useState(false),[storageOK,setStorageOK]=useState(false);
 const [markers,setMarkers]=useState<Partial<Record<"prepare"|"release",OperationMarker>>>({}),[receipt,setReceipt]=useState<CommercialReceipt|null>(null),[current,setCurrent]=useState<CurrentRelease|null>(null);
 const operations=useRef<Partial<Record<"prepare"|"release",OperationMarker>>>({}),fence=useRef(false);
 const storageKey=(action:string)=>`elite-handover:${scope}:${orderId}:${action}`;
 async function read(kind:string,extra:Record<string,string>={}){
  const response=await fetch(endpoint+"?"+new URLSearchParams({kind,organizationId:organization,orderId,...extra}),{cache:"no-store",signal:AbortSignal.timeout(10000)});
  if(!response.ok)throw new Error("unavailable");return response.json() as Promise<unknown>;
 }
 async function refresh(){
  setCurrent(null);
  try{const v=handoverContextSchema.parse(await read("context"));if(v.organization_id!==organization||v.order_id!==orderId)throw new Error("scope");setContext(v);setMessage(t("p0622"))}
  catch{setContext(null);setMessage(t("p0623"))}
 }
 useEffect(()=>{
  let mounted=true;
  try{
   if(!/^[a-f0-9]{64}$/.test(scope))throw new Error("scope");const values:Partial<Record<"prepare"|"release",OperationMarker>>={};
   for(const action of ["prepare","release"] as const){const raw=sessionStorage.getItem(storageKey(action));if(raw!==null){if(raw.length>1024)throw new Error("reference");const v=operationMarkerSchema.parse(JSON.parse(raw));if(v.action!==action||v.orderId!==orderId)throw new Error("scope");values[action]=v}}
   operations.current=values;setMarkers(values);setStorageOK(true);
  }catch{setStorageOK(false);setMessage(t("p0624"))}
  setReady(true);
  void read("context").then(raw=>{if(!mounted)return;const v=handoverContextSchema.parse(raw);if(v.order_id!==orderId||v.organization_id!==organization)throw new Error("scope");setContext(v)}).catch(()=>{if(mounted)setContext(null)});
  return()=>{mounted=false};
 // Identity changes remount this keyed component; no background polling.
 // eslint-disable-next-line react-hooks/exhaustive-deps
 },[scope,orderId,organization]);
 async function mutate(action:"prepare"|"release"){
  if(!ready||!storageOK||fence.current||operations.current[action]||!context)return;
  if(action==="prepare"?!context.can_prepare:context.handover?.state!=="accepted"||context.release_effect!=="COMMIT_COMMERCIAL_RELEASE_RECEIPT")return;
  fence.current=true;setBusy(true);setCurrent(null);
  let marker:OperationMarker;
  try{
   const requested:OperationMarker=action==="prepare"?{action,orderId,requestKey:crypto.randomUUID()}:{action,orderId,handoverId:context.handover!.id,requestKey:crypto.randomUUID()};
   if(sessionStorage.getItem(storageKey(action))!==null)throw new Error("existing operation");
   marker=retainOperation(sessionStorage,storageKey(action),requested);operations.current={...operations.current,[action]:marker};setMarkers(operations.current);
  }catch{setStorageOK(false);setBusy(false);fence.current=false;setMessage(t("p0625"));return}
  try{
   const response=await fetch(endpoint,{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify({action,organizationId:organization,orderId,requestKey:marker.requestKey}),signal:AbortSignal.timeout(10000)}),raw=await response.json();
   if(!response.ok)throw new Error("unconfirmed");
   if(action==="prepare"){
    const v=handoverViewSchema.parse(raw.handover);if(v.organization_id!==organization||v.order_id!==orderId)throw new Error("scope");
    setContext({...context,handover:v,can_prepare:false});setMessage(t("p0626"));
   }else{
    const v=commercialReceiptSchema.parse(raw.receipt);if(v.organization_id!==organization||v.order_id!==orderId||marker.action!=="release"||v.handover_id!==marker.handoverId)throw new Error("scope");setReceipt(v);setMessage(t("p0627"));
   }
  }catch{setMessage(t("p0628"))}
  finally{setBusy(false);fence.current=false}
 }
 async function recover(action:"prepare"|"release"){
  const marker=operations.current[action];if(!marker||fence.current)return;fence.current=true;setBusy(true);setCurrent(null);
  try{
   const raw=await read(action+"-result",{requestKey:marker.requestKey,...(marker.action==="release"?{handoverId:marker.handoverId}:{})}) as {handover?:unknown;receipt?:unknown};
   if(action==="prepare"){const h=handoverViewSchema.parse(raw.handover);if(h.order_id!==orderId||h.organization_id!==organization)throw new Error("scope");await refresh();setMessage(t("p0629", {handover_id:(h.id),handover_state:(h.state)}))}
   else{const r=commercialReceiptSchema.parse(raw.receipt);if(r.order_id!==orderId||r.organization_id!==organization||marker.action!=="release"||r.handover_id!==marker.handoverId)throw new Error("scope");setReceipt(r);setMessage(t("p0630"))}
  }catch{setMessage(t("p0631"))}
  finally{fence.current=false;setBusy(false)}
 }
 async function validate(){
  const id=receipt?.handover_id??context?.handover?.id??(operations.current.release?.action==="release"?operations.current.release.handoverId:undefined);if(!id||fence.current)return;
  fence.current=true;setBusy(true);setCurrent(null);
  try{const v=currentReleaseSchema.parse(await read("release-current",{handoverId:id}));if(v.receipt.organization_id!==organization||v.receipt.order_id!==orderId||v.receipt.handover_id!==id)throw new Error("scope");setReceipt(v.receipt);setCurrent(v)}catch{setMessage(t("p0632"))}
  finally{fence.current=false;setBusy(false)}
 }
 return <article className="card" aria-label={`Entrega del pedido ${orderId}`}><h3>{t("p0251")} {orderId}</h3>
 <button className="button" type="button" disabled={!ready||busy} onClick={()=>void refresh()}>{t("p0633")}</button>
 {!context?<p>{t("p0634")}</p>:context.handover?<><p>{t("p0368")} {context.handover.id} {t("p0443")} {context.handover.state}</p>{context.handover.state==="prepared"?<ChecklistCompletionPanel key={context.handover.id} organization={organization} scope={scope} initialHandover={{id:context.handover.id,version:context.handover.version}}/>:context.handover.state==="presented"?<p>{t("p0635")}</p>:null}</>:<button className="button" type="button" disabled={!ready||busy||!storageOK||!context.can_prepare||Boolean(markers.prepare)} onClick={()=>void mutate("prepare")}>{t("p0636")}</button>}
 {markers.prepare?<button className="button" type="button" disabled={!ready||busy} onClick={()=>void recover("prepare")}>{t("p0637")}</button>:null}
 {context?.handover?.state==="accepted"&&context.release_effect==="COMMIT_COMMERCIAL_RELEASE_RECEIPT"?<button className="button" type="button" disabled={!ready||busy||!storageOK||Boolean(markers.release)} onClick={()=>void mutate("release")}>{t("p0638")}</button>:context?.release_effect==="READ_ONLY_ELIGIBILITY"?<p>{t("p0639")}</p>:null}
 {markers.release?<button className="button" type="button" disabled={!ready||busy} onClick={()=>void recover("release")}>{t("p0640")}</button>:null}
 {receipt||context?.handover?.state==="accepted"?<button className="button" type="button" disabled={!ready||busy} onClick={()=>void validate()}>{t("p0641")}</button>:null}
 {receipt?<div><p>{t("p0642")} {receipt.id} {t("p0643")} {receipt.recorded_at}</p><p>{t("p0644")}</p></div>:null}
 {current?<p role="status">{current.current?t("p0645", {evaluated_at:(current.evaluated_at),valid_until:(current.receipt.valid_until)}):t("p0646", {evaluated_at:(current.evaluated_at)})}</p>:receipt?<p>{t("p0647")}</p>:null}
 <p role="status" aria-live="polite">{message}</p>
 </article>;
}
````

### FILE: `src/platform/backend/handover-recovery.test.ts`

```yaml
block_id: "HANDOVER-DELTA-TYPESCRIPT_FRANCHISE_JOURNEY_PORTALS:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "8440d2681fdf89ec1667a0986012280a4eb9f562a20a6308571583cecb82dc4a"
variables: []
secrets_allowed: false
```

````typescript
import{afterEach,expect,it,vi}from"vitest";
import{protectedGet}from"./protected-client";
const original=process.env.ENTERPRISE_API_BASE_URL;
afterEach(()=>{process.env.ENTERPRISE_API_BASE_URL=original;vi.unstubAllGlobals()});
it("forwards the recovery key in a header through the existing bounded client",async()=>{
 process.env.ENTERPRISE_API_BASE_URL="https://api.example.test";
 const fetch=vi.fn(async(_url:unknown,_init:RequestInit)=>new Response("{}",{headers:{"content-type":"application/json"}}));vi.stubGlobal("fetch",fetch);
 const s={subject:"operator",tenantId:"tenant",permissions:["handover:manage"],organizations:["store"],accessToken:"fixture-server-only"};
 await protectedGet(s,"/v1/franchise/orders/order/handover-result",{organization_id:"store"},"a03c7850-13ad-4b38-81d0-120000000001");
 expect(fetch.mock.calls[0]?.[1]).toMatchObject({cache:"no-store",redirect:"error",headers:{"Idempotency-Key":"a03c7850-13ad-4b38-81d0-120000000001",authorization:"Bearer fixture-server-only"}});
 expect(String(fetch.mock.calls[0]?.[0])).not.toContain("a03c7850");
 await expect(protectedGet(s,"/v1/franchise/orders/order/handover-result",{},"injected\r\nvalue")).rejects.toThrow();expect(fetch).toHaveBeenCalledTimes(1);
});
````

### FILE: `src/platform/handovers/contracts.test.ts`

```yaml
block_id: "HANDOVER-DELTA-TYPESCRIPT_FRANCHISE_JOURNEY_PORTALS:file9:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "dd1ea4a20e901d08f376f167c40fe23704977b5117aaa899a95c37a046b5d02a"
variables: []
secrets_allowed: false
```

````typescript
import {expect,it} from "vitest";
import {operationMarkerSchema,retainOperation,currentReleaseSchema} from "./contracts";
const key="a03c7850-13ad-4b38-81d0-120000000001";
it("retains one key and rejects corrupt/cross-order recovery without overwriting",()=>{
 const saved=new Map<string,string>(),storage={getItem:(k:string)=>saved.get(k)??null,setItem:(k:string,v:string)=>{saved.set(k,v)}};
 const old=retainOperation(storage,"scope",{action:"prepare",orderId:"order",requestKey:key});
 expect(retainOperation(storage,"scope",{action:"prepare",orderId:"order",requestKey:"a03c7850-13ad-4b38-81d0-120000000002"})).toEqual(old);
 expect(()=>retainOperation(storage,"scope",{action:"prepare",orderId:"other",requestKey:key})).toThrow();expect(JSON.parse(saved.get("scope")!)).toEqual(old);
 saved.set("scope","bad-json");expect(()=>retainOperation(storage,"scope",old)).toThrow();expect(saved.get("scope")).toBe("bad-json");
});
it("requires a handover reference for release recovery even if current context is unavailable",()=>{expect(operationMarkerSchema.safeParse({action:"release",orderId:"order",requestKey:key}).success).toBe(false);expect(operationMarkerSchema.safeParse({action:"release",orderId:"order",handoverId:"handover",requestKey:key}).success).toBe(true)});
it("storage failure prevents creation of an unsafe operation reference",()=>{expect(()=>retainOperation({getItem:()=>null,setItem:()=>{throw new Error("disabled")}},"scope",{action:"prepare",orderId:"order",requestKey:key})).toThrow()});
it("a historical receipt cannot manufacture a current positive result",()=>{expect(currentReleaseSchema.safeParse({receipt:{id:"receipt"},current:true}).success).toBe(false)});
````

### FILE: `src/platform/handovers/contracts.ts`

```yaml
block_id: "HANDOVER-DELTA-TYPESCRIPT_FRANCHISE_JOURNEY_PORTALS:file10:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "1b5d8a029aa06c4cbe30cf3047f76184cb44bff86a6bff6ae536e3680e93714b"
variables: []
secrets_allowed: false
```

````typescript
// AUTHORED DTO and recovery glue; all business decisions belong to the Go owner.
import { z } from "zod";
export const handoverID=z.string().min(1).max(128).regex(/^[A-Za-z0-9_-]+$/);
const hash=z.string().regex(/^[a-f0-9]{64}$/),version=z.number().int().positive().max(Number.MAX_SAFE_INTEGER);
const date=z.iso.datetime({offset:true});
const fundingFields={payment_attempt_id:z.union([handoverID,z.literal("")]),funding_receipt_id:handoverID.optional()};
const oneFunding=(v:{payment_attempt_id:string;funding_receipt_id?:string|undefined})=>Boolean(v.payment_attempt_id)!==Boolean(v.funding_receipt_id);
export const handoverViewSchema=z.object({id:handoverID,organization_id:handoverID,order_id:handoverID,state:z.enum(["prepared","presented","accepted","rejected"]),version,checklist_id:z.string().max(128).optional(),checklist_version:z.number().int().nonnegative().optional(),checklist_completed_at:date.optional(),customer_accepted_at:date.optional()});
export const handoverContextSchema=z.object({organization_id:handoverID,order_id:handoverID,handover:handoverViewSchema.optional(),can_prepare:z.boolean(),release_effect:z.enum(["READ_ONLY_ELIGIBILITY","COMMIT_COMMERCIAL_RELEASE_RECEIPT"]),evaluated_at:date}).superRefine((v,ctx)=>{
 if(v.can_prepare===Boolean(v.handover)||v.handover&&(v.handover.organization_id!==v.organization_id||v.handover.order_id!==v.order_id))ctx.addIssue({code:"custom",message:"inconsistent context"});
});
export const backendHandoverContextSchema=handoverContextSchema.safeExtend({order_line_id:handoverID,...fundingFields,observation_sha256:hash}).refine(oneFunding);
export const preparationSchema=z.object({handover:handoverViewSchema,order_line_id:handoverID,reservation_id:handoverID,...fundingFields,observation_sha256:hash,contract_id:z.string().min(1).max(128),contract_sha256:hash,prepared_by:z.string().min(1).max(256),prepared_at:date}).refine(oneFunding);
export const commercialReceiptSchema=z.object({id:handoverID,organization_id:handoverID,handover_id:handoverID,order_id:handoverID,...fundingFields,observation_sha256:hash,observation_generation:version,handover_version:version,acceptance_sha256:hash,checklist_id:handoverID,checklist_version:version,contract_id:z.string().min(1).max(128),contract_sha256:hash,effect:z.literal("COMMIT_COMMERCIAL_RELEASE_RECEIPT"),released_by:z.string().min(1).max(256),recorded_at:date,valid_until:date}).refine(oneFunding).refine(v=>Date.parse(v.valid_until)>Date.parse(v.recorded_at));
export const currentReleaseSchema=z.object({receipt:commercialReceiptSchema,evaluated_at:date,current:z.boolean()}).refine(v=>!v.current||Date.parse(v.evaluated_at)<Date.parse(v.receipt.valid_until));
export type HandoverContext=z.infer<typeof handoverContextSchema>;
export type CommercialReceipt=z.infer<typeof commercialReceiptSchema>;
export type CurrentRelease=z.infer<typeof currentReleaseSchema>;
export const operationMarkerSchema=z.discriminatedUnion("action",[
 z.object({action:z.literal("prepare"),orderId:handoverID,requestKey:z.uuid()}).strict(),
 z.object({action:z.literal("release"),orderId:handoverID,handoverId:handoverID,requestKey:z.uuid()}).strict(),
]);
export type OperationMarker=z.infer<typeof operationMarkerSchema>;
export function retainOperation(storage:Pick<Storage,"getItem"|"setItem">,key:string,value:OperationMarker):OperationMarker{
 const parsed=operationMarkerSchema.parse(value),existing=storage.getItem(key);
 if(existing!==null){const old=operationMarkerSchema.parse(JSON.parse(existing));if(old.action!==value.action||old.orderId!==value.orderId||old.action==="release"&&value.action==="release"&&old.handoverId!==value.handoverId)throw new Error("operation scope mismatch");return old;}
 storage.setItem(key,JSON.stringify(parsed));if(storage.getItem(key)!==JSON.stringify(parsed))throw new Error("operation reference not retained");return parsed;
}
````

### FILE: `src/platform/handovers/customer-dates.test.ts`

```yaml
block_id: "HANDOVER-DELTA-TYPESCRIPT_FRANCHISE_JOURNEY_PORTALS:file11:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "e833faa06efcddb89def0ff0ab2c11c885f0d66fbbc42955e34b17179f51f7ac"
variables: []
secrets_allowed: false
```

````typescript
import{createElement}from"react";
import{renderToStaticMarkup}from"react-dom/server";
import{expect,it}from"vitest";
import{CustomerHandoverActions}from"@/components/customer-handover-actions";
it("renders the exact server-formatted handover labels regardless of client timezone",()=>{
 const html=renderToStaticMarkup(createElement(CustomerHandoverActions,{organizationId:"store",exceptions:[],handovers:[{id:"handover",order_id:"order",stock_unit_id:"stock",state:"accepted",version:3,checklist_items:[{id:"serial",ordinal:1,prompt:"Serial",response_type:"serial",required:true}],checklist_id:"checklist",checklist_version:1,checklist_completed_at:"2026-09-11T12:00:00Z",customer_accepted_at:"2026-09-11T12:05:00Z",checklistCompletedLabel:"11/09/2026 09:00 (America/Argentina/Buenos_Aires)",acceptedLabel:"11/09/2026 09:05 (America/Argentina/Buenos_Aires)"}]}));
 expect(html).toContain("11/09/2026 09:00 (America/Argentina/Buenos_Aires)");expect(html).toContain("11/09/2026 09:05 (America/Argentina/Buenos_Aires)");
});
````

### FILE: `src/platform/handovers/render.test.ts`

```yaml
block_id: "HANDOVER-DELTA-TYPESCRIPT_FRANCHISE_JOURNEY_PORTALS:file12:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "e59c00349ee820506510e6e4e739a6e15d0b3b8c7f901b1872e470390701d001"
variables: []
secrets_allowed: false
```

````typescript
import {createElement}from"react";
import{renderToStaticMarkup}from"react-dom/server";
import{expect,it,vi}from"vitest";
vi.mock("next/navigation",()=>({useRouter:()=>({refresh:vi.fn()})}));
import{HandoverOperationsPanel}from"@/components/handover-operations-panel";
import{ChecklistCompletionPanel}from"@/components/franchise-command-panel";
it("exposes unavailable data distinctly and does not authorize a mutation before context",()=>{
 const absent=renderToStaticMarkup(createElement(HandoverOperationsPanel,{orders:null,organization:"store",scope:"a".repeat(64)}));expect(absent).toContain("lista de pedidos no está disponible");expect(absent).not.toContain("Registrar cierre comercial");
 const present=renderToStaticMarkup(createElement(HandoverOperationsPanel,{orders:[{id:"order"}],organization:"store",scope:"a".repeat(64)}));expect(present).toContain("Consultar estado de entrega");expect(present).not.toContain("Registrar cierre comercial");expect(present).not.toContain("Preparar entrega");
});
it("reuses the existing checklist owner with the exact handover/version read from backend",()=>{
 const html=renderToStaticMarkup(createElement(ChecklistCompletionPanel,{organization:"store",scope:"a".repeat(64),initialHandover:{id:"handover-from-server",version:7}}));
 expect(html).toContain('value="handover-from-server"');expect(html).toContain('value="7"');expect((html.match(/readOnly=""/g)??[]).length).toBe(2);expect(html).toContain("Completar y presentar");
});
````

### FILE: `src/app/franchise/whatsapp/page.tsx`

```yaml
block_id: "COMMUNICATIONS-DELTA-TYPESCRIPT_FRANCHISE_JOURNEY_PORTALS:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "4feaa986d7aa743bc27d06d65dd5d5ab493f71c1d86fff49f255b516bfdf5267"
variables: []
secrets_allowed: false
```

````tsx
import {loadPrivateI18n} from "@/platform/i18n/load-private-i18n";
import { redirect } from "next/navigation";
import type { Route } from "next";
import { allowed, readSession } from "@/platform/auth/session";
import { protectedGet } from "@/platform/backend/protected-client";
import { responseText, type ReplyView } from "@/platform/notifications/reply-contract";
import { WhatsAppReplyReview } from "@/components/whatsapp-reply-review";

export default async function WhatsAppReplyPage(){
 const {locale:privateLocale,t,controlled}=await loadPrivateI18n();

 const session=await readSession();if(!session)redirect("/api/auth/login?return_to=/franchise/whatsapp" as Route);
 if(!allowed(session,"whatsapp:approve"))return <><h1 className="pageTitle">{t("p0096")}</h1><p>{t("p0097")}</p></>;
 try{
  const replies=await protectedGet<ReplyView[]>(session,"/v1/franchise/whatsapp/replies",{});
  if(!Array.isArray(replies)||replies.length>50||replies.some(v=>!session.organizations.includes(v.context.organization_id)))throw new Error("INVALID_SCOPE");
  const rows=replies.map(value=>({...value,body:responseText(value),windowLabel:new Intl.DateTimeFormat(privateLocale.locale,{dateStyle:"medium",timeStyle:"short",timeZone:privateLocale.timeZone}).format(new Date(value.context.expires_at))+" ("+privateLocale.timeZone+")"}));
  return <><h1 className="pageTitle">{t("p0096")}</h1><p>{t("p0098")}</p><WhatsAppReplyReview replies={rows} canSend={allowed(session,"whatsapp:send")}/></>;
 }catch{return <><h1 className="pageTitle">{t("p0096")}</h1><p role="alert">{t("p0099")}</p><a className="button" href="/franchise/whatsapp">{t("p0095")}</a></>}
}
````

### FILE: `src/components/whatsapp-reply-review.tsx`

```yaml
block_id: "COMMUNICATIONS-DELTA-TYPESCRIPT_FRANCHISE_JOURNEY_PORTALS:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "58e685649c19ea32aa635a36594bd8e3db2f0fea7aa503fc91cca0ebcf4d6cb5"
variables: []
secrets_allowed: false
```

````tsx
"use client";
import {usePrivateI18n} from "@/platform/i18n/private-provider";

import { useState } from "react";
import type { ReplyView } from "@/platform/notifications/reply-contract";
type Row=ReplyView&{body:string;windowLabel:string};


export function WhatsAppReplyReview({replies,canSend}:{replies:Row[];canSend:boolean}){
 const {locale:privateLocale,t,controlled}=usePrivateI18n();
const sendLabels:Record<string,string>={not_started:t("p1035"),sending:t("p1036"),accepted:t("p1037"),unknown:t("p1038"),failed:t("p1039")};
const deliveryLabels:Record<string,string>={not_observed_by_this_reader:t("p1040"),observed_sent:t("p1041"),observed_delivered:t("p1042"),observed_read:t("p1043"),observed_failed:t("p1044"),ambiguous_latest_timestamp:t("p1045")};

 const [busy,setBusy]=useState<string|null>(null);const [notice,setNotice]=useState("");
 async function act(row:Row,action:"decision"|"send"|"recover",approved?:boolean){
  if(busy)return;setBusy(row.request_id);setNotice("");
  // Existing audit reasons are wire data, not translated display messages.
  const body={action,request_id:row.request_id,payload_sha256:row.payload_sha256,...(action==="decision"?{approved,reason:approved?"Texto y destinatario revisados en el portal":"Respuesta rechazada desde el portal"}:{})};
  try{const r=await fetch("/api/enterprise/franchise/whatsapp/replies",{method:"POST",headers:{"content-type":"application/json"},body:JSON.stringify(body)});if(!r.ok)throw new Error("UNCONFIRMED");window.location.reload()}
  catch{setNotice(t("p1048"))}
  // Stay disabled after uncertainty until the user reloads the authoritative state.
 }
 return <><p><a className="button" href="/franchise/whatsapp">{t("p1049")}</a></p><p role="status" aria-live="polite">{notice}</p>{replies.length===0?<p>{t("p1050")}</p>:replies.map(row=><article className="card" style={{marginTop:16,overflowWrap:"anywhere"}} key={row.request_id}>
  <h2>{t("p1051")} {row.context.message.ExternalID}</h2><p>{t("p0005")} {row.context.organization_id}</p><p>{t("p1052")} {row.windowLabel}</p><p style={{whiteSpace:"pre-wrap",overflowWrap:"anywhere"}}>{row.body}</p>
  <p>{t("p1053")} {row.state==="pending"?t("p0797"):row.state==="approved"?t("p0798"):t("p0337")}{t("p1054")} {row.status ? sendLabels[row.status.fence_state]??t("p1055") : t("p1035")}{t("p1056")} {row.status ? deliveryLabels[row.status.delivery_status]??t("p1055") : t("p1057")}{t("p0008")}</p>
  {row.state==="pending"&&<div style={{display:"flex",gap:12,flexWrap:"wrap"}}><button className="button" disabled={busy!==null} onClick={()=>void act(row,"decision",true)}>{t("p1058")}</button><button className="button" disabled={busy!==null} onClick={()=>void act(row,"decision",false)}>{t("p1059")}</button></div>}
  {row.state==="approved"&&canSend&&row.status?.fence_state==="not_started"&&!row.status.approval_expired&&<button className="button" disabled={busy!==null} onClick={()=>void act(row,"send")}>{t("p1060")}</button>}
  {row.state==="approved"&&canSend&&row.status?.reconciliation_required&&<><p>{t("p1061")}</p><button className="button" disabled={busy!==null} onClick={()=>void act(row,"recover")}>{t("p1062")}</button></>}
 </article>)}</>;
}
````

## 6. Configuration surface

No new variable or secret. It consumes `ENTERPRISE_API_BASE_URL`, `ENTERPRISE_TENANT_CODE`, `ENTERPRISE_ORGANIZATION_CODE` and the OIDC/session variables already validated by the compatible packs.

## 7. Dependency bill

| Package | Pin exacto | Uso | Licencia | Fuente oficial |
|---|---|---|---|---|
| Next.js | 16.3.2 | App Router and same-origin BFF | MIT | https://github.com/vercel/next.js |
| React | 19.2.8 | accessible interactive form | MIT | https://github.com/facebook/react |
| Zod | 4.4.3 | strict request contract | MIT | https://github.com/colinhacks/zod |

All pins are owned by the bridge lockfile; this overlay adds no dependency.

## 8. Apply order

Materialize after bridge 0.5.x, OIDC adapter 0.2.x and Go journey 0.8.x, before browser/Lighthouse gates. Existing systems must reserve `/locations`, `/franchise`, `/customer/quotes`, `/customer/appointments`, `/customer/handovers`, `/api/enterprise/appointments` and `/api/enterprise/franchise/commands`; any collision stops composition. Rollback removes navigation exposure first, then these fourteen files; backend checklist, exception, authorization, receipt, disposition, effect requests and acceptance evidence remain durable.

## 9. Verification

Run frozen offline install, `pnpm test`, strict `tsc --noEmit` and `pnpm build`. Exercise strict schema, backend conflict, cross-origin rejection before session, permission/organization scope, customer rejection, exact serial receipt, closed condition/inventory action and explicit rejection of browser-owned evidence, refund, payment, stock, accounting, fiscal, request ID and completion state. Then run Playwright desktop/mobile and production edge CSP/nonce/CSRF/rate-limit checks selected by the project.

## 10. Reconstruction evidence

On 2026-08-30 the V140 author tree passed strict TypeScript and nine Vitest files/38 tests. The BFF maps return listing, physical receipt and disposition to exact Go endpoints, computes receipt evidence server-side and rejects browser-owned downstream effects. Fresh Markdown reconstruction, frozen-lock production build and full profile gates govern final admission. Microsoft Playwright/Lighthouse remain independent gates; real IdP/backend consumers, approved return policy, provider reconciliation, exact-cost accounting/fiscal posting, legal acceptance and production edge remain conditions.

V378 public locations and appointment forms require the BFF0.5.7 locale helper. Keep canonical UTC appointment instants, consent and idempotency unchanged; resolve configured market timeZone and translate only presentation. Exact es/en/fallback connected qualification: `reconstruction_evidence/PUBLIC_JOURNEY_LOCALIZATION_V378.md`. Private portals keep their existing Spanish language and approval boundaries.

V379 integrates version-pinned search/read of two existing PUBLIC Spanish operational guides under current BFF/JWE permissions. Exact inline text and versions are shared, not copied into a separate CMS. This closes only that reference boundary, not publishing/archive, historical storage, tenant-confidential documents or training progress. See docs/journey-help.md and reconstruction_evidence/VERSIONED_JOURNEY_HELP_V379.md. No new dependency or domain write.

V380 extends the same help owner to all15existing versioned guides;13exact inline extractions, shared pure rendering, precise existing page/section predicates and a bounded15article response. No new domain action or business guidance. See docs/journey-help-coverage.md for the frozen15profile matrix, source before hashes and limits. Existing guide versions unchanged.

V382: exact selected-composition compatibility refreshed. Administrative lead reads are conditional on the existing lead grant; shared quote-format code corrects admin/customer minor-unit labels without repricing. TestAdministrativeReadBrowserPostgres uses real Go/JWKS/PostgreSQL and four browsers with synthetic identities; all production Go/SQL remains unchanged. The role implementation and quote acceptance/date/state logic retain exact source parity. See docs/private-portal-reads.md and reconstruction_evidence/PRIVATE_PORTAL_READS_V382.md. No full business, IdP, security, release or target admission.

V383: private read failures expose a fixed alert and explicit GET reconsultation in admin/customer/factory, without error details or command replay. Production build plus Go/JWKS/PostgreSQL and four browser projects prove persistent denial with the insufficient bearer and recovery after session correction; zero writes and unchanged durable snapshots.263web tests and canonical parity. Existing role/portal/Go business sources unchanged. See reconstruction_evidence/PRIVATE_READ_RECOVERY_V383.md; no full-control or target promotion.

V387: customer/factory cursor traversal connects three existing lists beyond25rows.81durable records per browser, four projects, reload/independent first-page recovery, zero writes and unchanged snapshots verified. Source is AUTHORED, business/API/SQL/dependencies unchanged. Supporting TEST02 progress, no whole-control or native security/release promotion. See reconstruction_evidence/PRIVATE_PORTAL_PAGINATION_V387.md.

V388: admin orders/cases/leads now traverse existing scoped keyset APIs using the shared navigator.276web tests and four browser projects exercise all six lists,163list records, no omissions/duplicates per list, permission removal, reload/first-page recovery and unchanged durable snapshots. All code remains AUTHORED/CONDITIONED; no full TEST02, native security or release promotion. See reconstruction_evidence/ADMIN_PORTAL_PAGINATION_V388.md.

V389: customer appointment times share trusted business locale/zone across SSR and interactive management, preserving original instants and cancellation behavior.282web tests plus eight connected browser/configuration runs, two appointment instants per view and seven-table unchanged snapshots. AUTHORED/CONDITIONED; no whole TEST02, security or release promotion. See reconstruction_evidence/CUSTOMER_APPOINTMENT_TIMEZONE_V389.md.

V390: customer cancellation receipt binding, synchronous single-submit and explicit GET recovery close the existing reference journey.294web tests;4actual cancellation browser projects,16unique durable cancellations/audits/outbox and28API writes;8read/time regression runs remain read-only. AUTHORED/CONDITIONED, no whole TEST02/security/release promotion. See reconstruction_evidence/CUSTOMER_CANCELLATION_RECOVERY_V390.md.

V402 composed delta: Connected handover browser/BFF/Go/PostgreSQL gate; bounded body, stable server date formatting, Next PageProps signatures. AUTHORED integration glue; existing domain and fixed upstreams unchanged.

V402 browser evidence: reconstruction_evidence/HANDOVER_BROWSER_V402.md. Production Webpack build includes TypeScript checking;42 direct-call and23 focused BFF/date tests PASS. Chromium desktop and mobile viewport through real BFF/API/PG cover lost-response recovery and callback invalidation. Hosted IdP/live payment/physical shipment not claimed.

V402 composed delta: Role-scoped navigation and human WhatsApp reply review/send/recovery with readable delivery status. Real browser/BFF/Go/PostgreSQL proof; no new business rule or automatic send.

V402 composed delta: V402 source-backed stored-value integration: exact remaining provider due, explicit payment/funding XOR, shared approval, bounded browser transport and optional host. See STORED_VALUE_OPERATOR_FLOW_V402.md; source/pack admission successor governs final claim. Existing provider-only behavior retained.


### FILE: `src/platform/backend/bounded-command.ts`

```yaml
block_id: "TS-FRANCHISE-JOURNEY-PORTALS:src/platform/backend/bounded-command.ts:v1"
operation: CREATE
provenance: AUTHORED
source: "Local typed source/transaction/transport/UI/recovery glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "bbd8f31396a0262e9416744387d985adf85947d18dbc648af3fa5aa28ada88b7"
variables: []
secrets_allowed: false
```

````typescript
// AUTHORED shared body bound, extracted without semantic change from the handover BFF.
export async function boundedRequestBytes(request: Request, maxBytes: number) {
 if (!Number.isSafeInteger(maxBytes) || maxBytes < 1 || maxBytes > 1048576) throw new Error("invalid body budget");
 if (!request.body) throw new Error("invalid handover body");
 const reader = request.body.getReader(); const chunks: Uint8Array[] = []; let bytes = 0; let timer: ReturnType<typeof setTimeout> | undefined;
 const deadline = new Promise<never>((_, reject) => { timer = setTimeout(() => reject(new Error("body deadline")), 2500); });
 try {
  while (true) {
   const part = await Promise.race([reader.read(), deadline]); if (part.done) break;
   bytes += part.value.byteLength; if (bytes > maxBytes) throw new Error("BODY_TOO_LARGE"); chunks.push(part.value);
  }
  const data = new Uint8Array(bytes); let offset = 0; for (const chunk of chunks) { data.set(chunk, offset); offset += chunk.byteLength; }
  return data;
 } finally { clearTimeout(timer); void reader.cancel().catch(() => {}); reader.releaseLock(); }
}

// Existing text budgets and UTF-8 behavior remain unchanged.
export async function boundedCommandBody(request: Request, maxBytes = 4096) {
 if (!Number.isSafeInteger(maxBytes) || maxBytes < 1 || maxBytes > 65536) throw new Error("invalid body budget");
 return new TextDecoder("utf-8", { fatal: true }).decode(await boundedRequestBytes(request,maxBytes));
}
````

V402 composed delta: T2804 connected training: existing versioned help/audit/shared approvals/outbox/BFF; opt-in host and navigation; bounded body retains original default. No new dependency or automatic grant. TRAINING_CONNECTED_RELEASE_V402.md.

V402 composed delta: T2804 catalog role source/edit/review/publication transport reuses original model/price writers and catalog owner. Bounded PNG and text defaults retained. No new dependency. CATALOG_ROLE_AUTHORING_RELEASE_V402.md.

V402 composed delta: T2804 help CMS and21same-release guides;15oldguides unchanged,5bounded training curricula revision2, no automatic grants. New optional host, shared inline text, actual browser/PG proofs. HELP_CMS_RELEASE_V402.md.

V402 composed delta: T2804 role metrics use existing domain read models with exact strings, organization/customer/factory/program permissions, NPS minimum/retention and no zero on unavailable. FAIL868 converted lead and FAIL457 bounded generic body corrected. ROLE_METRICS_RELEASE_V402.md/json; no new dependency or corporate authorship.

### FILE: `src/app/api/enterprise/franchise/commands/body-bound.test.ts`

```yaml
block_id: "TS-FRANCHISE-JOURNEY-PORTALS-METRIC-DELTA:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "9f0e250e72838a01b2ed676bd95292354ef005a2795527a46c0c87038baa1c0b"
variables: []
secrets_allowed: false
```

````typescript
import{it,expect,vi,afterEach}from"vitest";import{NextRequest}from"next/server";
const m=vi.hoisted(()=>({session:vi.fn(),get:vi.fn(),post:vi.fn()}));
vi.mock("@/platform/auth/session",()=>({readSession:m.session,allowed:(s:{permissions:string[]},p:string)=>s.permissions.includes("*")||s.permissions.includes(p)}));
vi.mock("@/platform/auth/oidc-client",()=>({applicationBaseUrl:()=>new URL("https://portal.example.test")}));
vi.mock("@/platform/backend/protected-client",()=>({protectedGet:m.get,protectedPost:m.post}));
import{POST}from"./route";afterEach(()=>{vi.clearAllMocks();vi.useRealTimers()});
const req=(body:BodyInit)=>new NextRequest("https://portal.example.test/api/enterprise/franchise/commands",{method:"POST",headers:{origin:"https://portal.example.test","content-type":"application/json"},body,duplex:"half"as const});
it("auth precedes body consumption",async()=>{for(const session of[null,{permissions:[],organizations:["org"]}]){m.session.mockResolvedValue(session);const r=req("{}");expect([401,403]).toContain((await POST(r)).status);expect(r.bodyUsed).toBe(false)}expect(m.post).not.toHaveBeenCalled()});
it("bounds byte stream without content length and cancels it",async()=>{m.session.mockResolvedValue({permissions:["customer:self"],organizations:["org"]});let cancelled=false;const r=req(new ReadableStream<Uint8Array>({pull(c){c.enqueue(new Uint8Array(65537))},cancel(){cancelled=true}}));expect((await POST(r)).status).toBe(413);expect(cancelled).toBe(true);expect(m.post).not.toHaveBeenCalled()});
it("cancels stalled body at deadline and rejects malformed UTF8",async()=>{m.session.mockResolvedValue({permissions:["customer:self"],organizations:["org"]});vi.useFakeTimers();let cancelled=false;const pending=POST(req(new ReadableStream<Uint8Array>({cancel(){cancelled=true}})));await vi.advanceTimersByTimeAsync(2600);expect((await pending).status).toBe(400);expect(cancelled).toBe(true);vi.useRealTimers();expect((await POST(req(new Uint8Array([255])))).status).toBe(400);expect(m.post).not.toHaveBeenCalled()});
it("retains action permission checks after parse",async()=>{m.session.mockResolvedValue({permissions:["customer:self"],organizations:["org"]});expect((await POST(req(JSON.stringify({action:"assign-lead",organizationId:"org",leadId:"lead",assignedSubject:"maker",version:1})))).status).toBe(403);expect(m.post).not.toHaveBeenCalled()});
````


V402 composed delta: T2804 private es/en display, per-user/tenant preference and browser negotiation;1198messages,21exact guides,5hash-bound curricula; original commands/content/policies preserved. PRIVATE_LOCALE_RELEASE_V402.md/json. AUTHORED glue; no new dependency or corporate attribution.

### FILE: `src/app/api/enterprise/locale/route.test.ts`

```yaml
block_id: "TS-FRANCHISE-JOURNEY-PORTALS-LOCALE-DELTA:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "971c8180196a9c6a9ad0a11728b29d1b70ed7bf65ee524a82cf000d8f2a47ac9"
variables: []
secrets_allowed: false
```

````typescript
import{it,expect,vi,afterEach}from"vitest";import{NextRequest}from"next/server";import{createHash}from"node:crypto";
const m=vi.hoisted(()=>({session:vi.fn()}));
vi.mock("@/platform/auth/session",()=>({readSession:m.session}));
vi.mock("@/platform/auth/oidc-client",()=>({applicationBaseUrl:()=>new URL("https://portal.example.test")}));
vi.mock("@/platform/i18n/load-private-locale",()=>({privateLocaleCookieName:(s:{tenantId:string;subject:string})=>"elite_locale_"+createHash("sha256").update(JSON.stringify([s.tenantId,s.subject])).digest("hex")}));
import{POST}from"./route";afterEach(()=>vi.clearAllMocks());
const req=(body:BodyInit,origin="https://portal.example.test")=>new NextRequest("https://portal.example.test/api/enterprise/locale",{method:"POST",headers:{origin,"content-type":"application/json"},body,duplex:"half"as const});
it("authenticates before reading and rejects cross-origin preference writes",async()=>{m.session.mockResolvedValue(null);const r=req('{"language":"en"}');expect((await POST(r)).status).toBe(401);expect(r.bodyUsed).toBe(false);expect((await POST(req('{"language":"en"}',"https://other.test"))).status).toBe(403)});
it("allows only one bounded display language and no actor/policy field",async()=>{m.session.mockResolvedValue({tenantId:"t",subject:"user"});for(const body of['{"language":"fr"}','{"language":"en","language":"es"}','{"language":"en","tenantId":"other"}',"x".repeat(65)])expect((await POST(req(body))).status).toBe(400)});
it("preference cookie is private and differs by authenticated subject and tenant",async()=>{const names=[];for(const [tenantId,subject]of[["t","a"],["t","b"],["other","a"]]){m.session.mockResolvedValue({tenantId,subject});const response=await POST(req('{"language":"en"}'));expect(response.status).toBe(200);const cookie=response.headers.get("set-cookie")!;expect(cookie).toContain("HttpOnly");expect(cookie).toContain("SameSite=lax");expect(response.headers.get("cache-control")).toContain("no-store");names.push(cookie.split("=")[0])}expect(new Set(names).size).toBe(3)});
````

### FILE: `src/app/api/enterprise/locale/route.ts`

```yaml
block_id: "TS-FRANCHISE-JOURNEY-PORTALS-LOCALE-DELTA:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "83458d1331bd74794844a2b415670ddec65fb8e5ae406e164d362b70aa9abcf6"
variables: []
secrets_allowed: false
```

````typescript
import{NextResponse,type NextRequest}from"next/server";
import{readSession}from"@/platform/auth/session";
import{applicationBaseUrl}from"@/platform/auth/oidc-client";
import{boundedCommandBody}from"@/platform/backend/bounded-command";
import{privateLocaleCookieName}from"@/platform/i18n/load-private-locale";
const reply=(code:string,status:number)=>NextResponse.json({code},{status,headers:{"cache-control":"private, no-store",vary:"Cookie"}});
export async function POST(r:NextRequest){
 try{if(r.headers.get("origin")!==applicationBaseUrl().origin)return reply("CROSS_ORIGIN_REJECTED",403)}catch{return reply("UNAVAILABLE",503)}
 const session=await readSession();if(!session)return reply("UNAUTHENTICATED",401);
 if(r.nextUrl.searchParams.size)return reply("INVALID_QUERY",400);
 if(r.headers.get("content-type")!=="application/json")return reply("UNSUPPORTED_MEDIA_TYPE",415);
 let text:string;try{text=await boundedCommandBody(r,64)}catch{return reply("INVALID_PREFERENCE",400)}
 // Exact one-field request rejects duplicate keys as well as identity/policy fields.
 const parsed=/^\{\s*"language"\s*:\s*"(es|en)"\s*\}$/.exec(text);if(!parsed)return reply("INVALID_PREFERENCE",400);
 const response=NextResponse.json({language:parsed[1]},{headers:{"cache-control":"private, no-store",vary:"Cookie"}});
 response.cookies.set(privateLocaleCookieName(session),parsed[1]!,{httpOnly:true,secure:process.env.NODE_ENV==="production",sameSite:"lax",path:"/",maxAge:31536000});
 return response;
}
````

### FILE: `src/platform/help/query.ts`

```yaml
block_id: "TS-FRANCHISE-JOURNEY-PORTALS-LOCALE-DELTA:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "45d561ee81c44bc1751db056b7b7585dce90b790c28e03507edecccd26567bac"
variables: []
secrets_allowed: false
```

````typescript
export function parseHelpQuery(url: string) {
  if (url.length > 2048) return null;
  const p = new URL(url).searchParams;
  if ([...p.keys()].some(key => !["q", "article", "version"].includes(key) || p.getAll(key).length !== 1)) return null;
  const q = p.get("q") ?? "", article = p.get("article"), version = p.get("version");
  if (q.length > 160 || /[\u0000-\u001f\u007f]/.test(q)) return null;
  if ((article === null) !== (version === null) || (article !== null && (p.has("q") || !/^[a-z][a-z0-9-]{0,63}$/.test(article) || !/^\d{1,6}\.\d{1,6}\.\d{1,6}$/.test(version!)))) return null;
  return { q: q.trim().normalize("NFC").toLocaleLowerCase("es"), article, version };
}
````

### FILE: `src/platform/i18n/private-guide-display.ts`

```yaml
block_id: "TS-FRANCHISE-JOURNEY-PORTALS-LOCALE-DELTA:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed configuration, persistence, authorization, UI and orchestration glue around explicitly selected owners and fixed official SDKs; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "5348fcc095d238d42a654349919b2f6bb46fd6ad14ee67b6f848784aa7b411a3"
variables: []
secrets_allowed: false
```

````typescript
// AUTHORED display adapter. Immutable source content and training hashes are unchanged.
import{ALL_GUIDES}from"@/platform/help/content";
import{controlledPrivateLabel}from"./private-catalog";
import type{PrivateLanguage}from"./private-locale";
type Guide={id:string;version:string;title:string;paragraphs:readonly string[];summary?:string;inlineLead?:boolean};
export function guideDisplay(source:Guide,language:PrivateLanguage){
 const original=ALL_GUIDES.find(x=>x.id===source.id&&x.version===source.version);
 const known=original!==undefined&&original.title===source.title&&(source.summary===undefined||("summary"in original&&source.summary===original.summary))&&original.paragraphs.length===source.paragraphs.length&&original.paragraphs.every((p,i)=>p===source.paragraphs[i]);
 const translate=(text:string)=>known?controlledPrivateLabel(language,text):text;
 return{id:source.id,version:source.version,title:translate(source.title),paragraphs:source.paragraphs.map(translate),summary:translate(source.summary??source.title),inlineLead:source.inlineLead??false,displayLanguage:known?language:"und"};
}
````

