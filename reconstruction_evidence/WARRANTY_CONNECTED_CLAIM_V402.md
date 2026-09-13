# V402 — reclamo de garantía conectado al producto de referencia

Checkpoint287. El backend J1/J4 queda probado localmente en staging: cotización
con términos/consentimiento, pedido, SDK pago fixture, callback/reconciliación,
entrega y recibo comercial, activación, cita real, diagnóstico, aprobación humana,
repuestos, calidad, aceptación del cliente y conciliación de fábrica. T2802 no
se cierra aún: faltan API/host/UI, gates estrechos pendientes y publicación.

Un test PostgreSQL conectado sin skips cubre cuatro casos: uno closed y tres
cancelled por rechazo, exclusión de cobertura o retiro antes de consumir stock.
La agenda usa su servicio y repositorio con el mismo perfil materializado fixture
(adcabf25387287f12b08600d910bee45991b09c2c27099610e2d3a6fcce751a3),
leadtime0 y slots1segundo explícitos. Request, asignación, confirmación y completion
se ejecutan realmente con reloj real. Calendario/localización son inputs fixture.
El perfil operativo posterior no cambia los términos de garantía ya vendidos.

Los writers de stock/caso se componen en Serializable. Dos capas FIFO de5unidades
a100/120 producen una salida de7, dos aplicaciones y costo exacto740; quedan
3unidades y0reservadas. Una falla después de los writes revierte stock, costos,
caso y recibo. La aprobación conserva separación de sujetos; el endpoint genérico
de caso no puede saltar el recibo del comando de garantía. Calidad rechazada no
permite aceptación; su corrección conserva evidencia. Otro cliente/organización
no acepta/concilia. Replays no duplican decisiones, consumo ni conciliación.

La conciliación es INTERNAL_WARRANTY_SERVICE_ACKNOWLEDGEMENT, con costo de
inventario y referencias reales; nunca reembolso pagado ni documento fiscal.
La cancelación técnica previa al consumo requiere permiso propio, libera la
reserva existente y conserva la decisión humana. La reserva pendiente tiene
lease explícito; aprobada pasa a hold durable de reparación y puede cancelarse
antes del consumo. El test de expiración al commit aún está pendiente.

Un test focal de registry verifica que warranty_repair no se autoaprueba aunque
la policy general tenga umbral; los dos tests del perfil/configuración siguen
PASS y vet final PASS. Las extracciones de release/decision conservan SQL/orden.
No se acredita ejecución de AL ni una arquitectura empresarial por reputación.

FAIL833 preparación, FAIL834 contrato del fixture de cita, FAIL835 mapeo de versión
de reserva y FAIL836 formato textual decimal están corregidos en staging con RED
retenidos. Se reutilizaron owners y se repitió únicamente el recorrido afectado.
El último marcador heredado del log de términos todavía dice claim_pending; el
marcador J4 nuevo y los cuatro casos muestran lo ejecutado. No reescribir el log.

Inventario/hashes y delta congelado: WARRANTY_CONNECTED_CLAIM_V402.json. Producto
canónico84/1112 sin cambio; garantía todavía no suma un pack completo. PG67
migraciones y proceso detenido. Siguiente: transporte/host, expiración, downgrade,
publicación exacta y reconciliación del gap map; luego T2804.
