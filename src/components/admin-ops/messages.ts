// AUTHORED presentation labels. Domain enum values are never translated on the wire.
const messages = {
  deliveries: ["Entregas", "Deliveries"], orders: ["Pedidos", "Orders"], contacts: ["Contactos", "Contacts"],
  search: ["Buscar", "Search"], allStates: ["Todos los estados", "All statuses"], status: ["Estado", "Status"],
  noMatches: ["No hay coincidencias", "No matches"], clearFilters: ["Limpiar filtros", "Clear filters"],
  empty: ["No hay registros todavía", "No records yet"], choose: ["Elegí un registro para continuar", "Choose a record to continue"],
  chooseDelivery: ["Elegí una entrega", "Choose a delivery"], chooseOrder: ["Elegí un pedido", "Choose an order"],
  chooseContact: ["Elegí un contacto", "Choose a contact"], back: ["Volver a la lista", "Back to list"],
  reference: ["Referencia", "Reference"], loading: ["Cargando entrega…", "Loading delivery…"],
  unavailable: ["No pudimos cargar la entrega", "We could not load the delivery"],
  forbidden: ["No tenés permiso para ver esta entrega", "You do not have permission to view this delivery"],
  retry: ["Volver a consultar", "Check again"], noDeliveries: ["No hay pedidos para preparar", "No orders to prepare"],
  unavailableList: ["No pudimos cargar tus pedidos", "We could not load your orders"],
  readyStep: ["Preparación", "Preparation"], reviewStep: ["Revisión", "Review"], customerStep: ["Aceptación", "Acceptance"], releaseStep: ["Liberación", "Release"],
  pendingPrepare: ["Listo para preparar", "Ready to prepare"], noPermissionPrepare: ["La preparación no está habilitada para este pedido.", "Preparation is not available for this order."],
  pendingAcceptance: ["La revisión fue presentada. Esperamos la aceptación del cliente.", "The review was submitted. Waiting for the customer's acceptance."],
  refresh: ["Actualizar", "Refresh"], serialHelp: ["Escribí la serie o usá un lector conectado. Revisá el valor antes de confirmar.", "Type the serial or use a connected scanner. Review the value before confirming."],
  checklist: ["Revisar entrega", "Review delivery"], guide: ["Ayuda", "Help"], optional: ["Opcional", "Optional"],
  referenceOrder: ["Pedido", "Order"], noStatus: ["Sin estado informado", "Status not supplied"],
  recordCount: ["registros", "records"], record: ["registro", "record"], selected: ["Seleccionado", "Selected"],
  confirmed: ["Confirmado", "Confirmed"], placed: ["Recibido", "Placed"], delivered: ["Entregado", "Delivered"], cancelled: ["Cancelado", "Cancelled"],
  prepared: ["En preparación", "Preparing"], presented: ["Espera aceptación", "Awaiting acceptance"], accepted: ["Aceptado", "Accepted"], rejected: ["Requiere revisión", "Needs review"],
  new: ["Nuevo", "New"], contacted: ["Contactado", "Contacted"], qualified: ["Calificado", "Qualified"], converted: ["Convertido", "Converted"], lost: ["No concretado", "Not converted"], closed: ["Cerrado", "Closed"],
  created: ["Creado", "Created"], pending: ["Pendiente", "Pending"], authorized: ["Autorizado", "Authorized"], captured: ["Cobrado", "Captured"], failed: ["Fallido", "Failed"], refunded: ["Reintegrado", "Refunded"], disputed: ["En disputa", "Disputed"],
  receipt: ["Ver comprobante", "View receipt"], safety: ["Antes de confirmar", "Before confirming"],
  oneSelection: ["El detalle se abre al seleccionar un registro.", "Select a record to open its details."],
} as const;
export type AdminMessage = keyof typeof messages;
export function adminText(key: AdminMessage, language: string): string { return messages[key][language === "en" ? 1 : 0]; }
export function adminState(value: string | undefined, language: string): string {
  if (!value) return adminText("noStatus", language);
  return value in messages ? adminText(value as AdminMessage, language) : value;
}
