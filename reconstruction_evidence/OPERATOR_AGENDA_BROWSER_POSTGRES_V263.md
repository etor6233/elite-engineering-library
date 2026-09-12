# V263 — agenda operativa conectada a confirmación durable

Fecha: 2026-09-06. Mantenimiento de biblioteca. Estado: REBUILD_VERIFIED / CONDITIONED, sin promoción global.

## Resultado tangible

`operador → agenda por fecha → recurso por nombre → asignar → confirmar → PostgreSQL/audit/outbox → estado Confirmado → recargar y verificar`

Ya no se copian IDs ni versiones para asignar o transicionar un turno. La agenda usa datos del servidor, muestra recursos y estados y valida el receipt antes de anunciar éxito. Una respuesta incierta bloquea más escrituras hasta consultar el estado guardado. Los formularios manuales anteriores de esas dos acciones se eliminaron; no se duplicó el dominio.

Se incorporó GET /v1/franchise/agenda con permiso appointment:manage y aislamiento tenant/organización. Rango RFC3339 máximo 31 días, inicio inclusivo/fin exclusivo; UI por día UTC explícito. Snapshot read-only repeatable-read, hasta 200 turnos y 200 recursos activos; overflow marca truncated y bloquea acciones en una lista incompleta. No se exponen subjects de recursos. Elegibilidad final sigue en los owners de asignación/disponibilidad existentes.

## Evidencia ejecutada

- **4/4 Microsoft Playwright**: Chromium desktop/mobile, Firefox y WebKit. Sesión JWE sintética Secure/HttpOnly, tokens RSA verificados por discovery/JWKS local, Next construido, API Go y PostgreSQL reales.
- Chromium desktop demuestra confirmación normal. Los otros tres pierden la respuesta **después del commit real**; no muestran éxito no comprobado ni vuelven a enviar confirmación. Actualizar agenda obtiene el estado durable. Después de recargar sigue Confirmado; sin overflow horizontal.
- Por tenant: un turno confirmed/version 3, un recurso asignado, una transición con actor operator y un evento appointment.confirmed. Se prueban lecturas de tenant/organización ajenos y límite de recursos con 201 fixtures adicionales.
- Gate BFF anterior de dos confirmaciones concurrentes: una aceptación y un 409, lectura sólo del cliente dueño, PASS separado.
- Siete casos HTTP adicionales: permiso, autenticación, organización, fechas ausentes, rango excesivo/invertido y respuesta no-store.
- Web: **54 tests PASS**, un test conectado omitido expresamente sin opt-in; typecheck/build PASS. El omitido pasó por su runner conectado.
- Go test ./..., vet ./..., build ./... PASS. PostgreSQL 18.6, 49 migraciones y suite de persistencia completa PASS.
- Fuente canónica reconstruida: **67 packs / 704 archivos**, **12/12 archivos modificados idénticos por SHA-256**. Gates conectados repetidos desde el backend reconstruido y cliente web de bytes comparados.

## Defectos encontrados y corregidos

La primera asignación en navegador recibió 403/CROSS_ORIGIN_REJECTED: se comparaba Origin con la URL interna de Next. Dos regresiones reprodujeron el problema. Se reutiliza applicationBaseUrl/APP_BASE_URL como origen público exacto; ausencia/configuración inválida rechazan, Host/X-Forwarded-Host no conceden confianza.

El primer shim HTTPS→HTTP de navegador no transportó la cookie segura en WebKit. Se sustituyó por un **proxy TLS local real** con certificado efímero autofirmado. La excepción de CA se limita al fixture loopback explícito; las advertencias de certificado de las herramientas quedan en el log. No es evidencia de certificados confiables, TLS desplegado ni OIDC login. No se desactivaron Secure, HttpOnly, firma de tokens, permisos ni triggers.

Historial completo: FAIL-20260906-357–361 y LIB-FAIL-2091–2095, incluyendo errores de edición y diagnósticos fallidos. No se contaron ejecuciones fallidas como PASS.

## Canonicalidad y procedencia

- GO-FRANCHISE-CUSTOMER-JOURNEY-API **0.10.0**: cinco archivos existentes.
- TS-FRANCHISE-JOURNEY-PORTALS **0.10.0**: cuatro archivos existentes.
- MICROSOFT-PLAYWRIGHT-BROWSER-GATE **0.1.5**: tres archivos existentes.

Cero dependencias, migraciones o archivos materializables nuevos. Read model, UI, fixture y correcciones son **AUTHORED**, no código copiado de Microsoft/Meta. Runtime Microsoft Playwright 1.62.1 y sus browsers/lock/licencia no cambiaron; se reutilizan los componentes previamente admitidos, sin afirmar SCA global nuevo.

Autoridades consultadas el 2026-09-06:

- [Microsoft Field Service — Schedule board](https://learn.microsoft.com/en-us/dynamics365/field-service/work-with-schedule-board): vista de reservas/recursos y selección por fecha. Patrón funcional, no código Dynamics trasladado a Go.
- [Meta React — Choosing the State Structure](https://react.dev/learn/choosing-the-state-structure): evitar estados redundantes; aquí el servidor conserva el estado de negocio y React controla interacción/resultado incierto.
- [Microsoft Playwright — Best Practices](https://playwright.dev/docs/best-practices): comportamiento visible, datos aislados y assertions sobre la experiencia.
- [Next.js — Self-hosting](https://nextjs.org/docs/app/guides/self-hosting) y [Server Actions security](https://nextjs.org/blog/security-nextjs-server-components-actions): fronteras de proxy/origen. Nuestra ruta es un Route Handler, no recibe automáticamente las protecciones de Server Actions; mantiene su chequeo explícito.

## Reproducir y límites pendientes

El README materializado del gate incluye el comando exacto. Instalar locks frozen, construir web, aplicar migraciones a base descartable elite_confirmation_* en 127.0.0.1, definir TEST_DATABASE_URL, ELITE_WEB_ROOT absoluto y ELITE_CONFIRMATION_E2E=1. Ejecutar:

```text
go test ./internal/platform/httpapi -run '^TestAppointment(AgendaBrowserPostgres|ConfirmationBFFPostgres|AgendaAuthorizationAndRange)$' -v -count=1 -timeout=3m
```

El puerto local 4173 debe estar libre. Las fixtures inmutables se conservan para inspección hasta descartar exclusivamente la base de ensayo; no limpiar desactivando triggers. Staging: elite-v263-04f2eba189414d878259ce3cc879c98d bajo el temporal del sistema; logs connected-final.log y roundtrip-connected.log. Captura confirmed-agenda.png producida por Playwright e inspeccionada visualmente durante esta tarea.

Pendiente: login/IdP real, calendario empresarial con zona horaria elegida, selección cómoda entre organizaciones, paginación más allá del límite, capacitación/ayuda completa, accesibilidad asistiva, carga/ofensiva/PKI/deploy y notificación externa. Completar/cancelar/no-show conservan API/controles, pero esta nueva prueba browser sólo acredita asignación/confirmación/recovery. Un outbox no acredita mensaje entregado. La biblioteca completa y el proyecto productivo no se declaran terminados; no se generó ZIP final ni porcentaje global.

## Verificador y limpieza

VERIFY_LIBRARY_PASS: 160 packs, 1.395 archivos materializables, 696 Markdown, 51 perfiles; franquicia 67/704. Memoria sincronizada: 2.095 lecciones locales + 209 condiciones upstream = 2.304 IDs, sin ocultar condiciones upstream abiertas.

La base de confirmación creada exclusivamente en esta tarea conservaba 19 tenants sintéticos de intentos y repeticiones, 16 turnos confirmados, 16 audits y 16 eventos. Los tres intentos no confirmados proceden de gates fallidos, no se presentan como PASS. Tras comprobar ausencia de tenants ajenos se descartó elite_confirmation_v263. También se descartó elite_browser_v263, creada aquí para las suites de persistencia, con sus nueve tenants de fixture restantes. Son datos sintéticos reconstruibles, no datos de proyecto recuperables. PostgreSQL temporal y los servidores Next/proxy de ensayo quedaron detenidos.
