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
