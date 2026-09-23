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
