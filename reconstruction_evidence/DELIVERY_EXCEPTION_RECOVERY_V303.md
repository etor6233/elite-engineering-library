# V303 — resolución de discrepancias y recuperación verificable

> **Errata de trazabilidad V317 (2026-09-08), prevalece sobre las menciones históricas de Go1.26.7 debajo.** Se retira la atribución histórica exacta Go1.26.7 no vinculada a un receipt de identidad por ejecución. No se sustituye por1.26.8 sin prueba del artefacto concreto. La ruta current-toolchain de V281 contiene1.26.8; official-toolchain contiene1.26.7. El cuerpo original queda conservado como historia. V317 revalida únicamente la composición actual y el scope allí enumerado con1.26.7 observado, no todos los snapshots previos. Ver [TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md](TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md).

2026-09-07. Checkpoint34/V302 → corrección T2803/T2804 del roadmap existente.
V295/preparación de credenciales y ARCA diferida conservadas. Sin cuentas,
compras, secretos ni efectos externos. No reemplaza el creador inicial pendiente.

## Lo implementado

TS-FRANCHISE-JOURNEY-PORTALS0.14.1 conserva el owner de página/componentes y
agrega DeliveryExceptionResolution dentro del mismo archivo. La resolución ya
existente deja de usar el sender genérico defectuoso. Tres decisiones originales:
correct-and-represent, return y exchange; no nueva política comercial.

- Controles cerrados antes de hidratación; useRef excluye otro envío mientras
  el resultado se resuelve. Datos de formulario se conservan ante error.
- Antes del POST guarda únicamente versión y decisión en sessionStorage, bajo
  hash de tenant/subject/organización y el ID opaco de excepción. Sin notas,
  credenciales ni tokens. Si falla esa escritura no se envía el comando.
- Timeout10s; recibo debe coincidir con excepción, organización, acta, versión,
  estado, decisión y vínculo de sucesor/autorización antes de decir registrado.
  JSON vacío, error de transporte o respuesta no demostrada quedan inciertos.
- «Consultar resolución» sólo recarga por GET. El marker sobrevive a esa recarga;
  consulta fallida/ausencia/versiones no demostradas no lo eliminan. Un estado
  resolved posterior recuperado permite mostrar la decisión real, sin atribuir
  a esta solicitud una acción de otro operador. Una decisión distinta se señala.
- Ayuda delivery-resolution-view/1.0.0 explica incertidumbre, práctica sintética,
  soporte autorizado y límites. No es capacitación evaluada ni soporte integral.

GO-FRANCHISE-CUSTOMER-JOURNEY-API0.10.8 y
MICROSOFT-PLAYWRIGHT-BROWSER-GATE0.1.13 amplían únicamente el harness y sus tests.
El dominio, SQL de producto, autorización y migraciones no cambian. Se conserva
ResolveDeliveryException: CAS/versiones, transacción, resolución y outbox únicos.
No hay módulo, dependencia o archivo de producto nuevo:4 archivos modificados.
Perfil integral67 packs/745 archivos sin paths duplicados; planes backend/web/
integral/serverless sincronizados. Los archivos vuelven al Markdown canónico.

## Procedencia y fuentes

Todo el delta es AUTHORED / LicenseRef-Workspace-Owner, no código extraído de AWS
o Meta. Los frameworks/SDK ya fijados conservan locks/licencias y condiciones.
Consultados el2026-09-07:

- [AWS Builders' Library — retries e idempotencia](https://aws.amazon.com/builders-library/making-retries-safe-with-idempotent-APIs/): una respuesta perdida no demuestra ausencia de efecto; no reintentar ciegamente. Aquí no se incorpora un SDK AWS ni se inventa un contrato idempotente nuevo.
- [React — refs](https://react.dev/learn/referencing-values-with-refs): referencia mutable entre renders; el ref no sustituye la exclusión transaccional en PostgreSQL.

FRONTEND_PRODUCT_ENGINEERING_UX sección2.2 y SOFTWARE_BACKEND_API_ENGINEERING
sección8 gobiernan visibilidad, recuperación y semántica de reintentos. Citar
estas autoridades no implica que aprueben este código local.

## Evidencia reproducible

Staging: `%LOCALAPPDATA%/Temp/elite-v303-d4bbf1e6dc86493694cd2f8ed861e68f`.
baseline: V302 construido, con nuevo test temporal; source original restaurado
después de reproducir. candidate compila la web; rebuilt se compone desde los
packs actualizados a destino ausente y sus4 archivos coinciden por SHA256.

- resolution-red.log: POST real devolvió200/resolved y sucesor prepared antes
  de abortar la respuesta. UI V302 mostró «No se aplicó: Failed to fetch» y
  Registrar resolución habilitado; faltó el status de recuperación esperado.
  Trace/error-context preservados, no se atribuye el rojo a bypass de permisos.
- browser-final.log:56 fases (14 por proyecto) PASS en Chromium desktop/mobile,
  Firefox desktop y WebKit desktop;168.807s de paquete. Web/BFF/Go/OIDC sintético/
  PostgreSQL18.6 reales locales; conserva los recorridos V294–V302.
- Cada fase nueva resuelve3 excepciones: correct-and-represent pierde respuesta
  después del commit; return recibe JSON200 vacío después del commit; exchange
  primero prueba storage fallido con cero POST y luego dos requests concurrentes
  de prueba contra el backend (200/409), con recibo válido para la UI.
- Tras la primera resolución se inyecta503 en la consulta. No se muestran notas
  internas, marker permanece y otro GET recupera. Cada acción tiene un único
  POST de navegador; los duplicados adicionales son inyección explícita del test,
  no retries automáticos del producto. Replay409, organización ajena403 y token
  sin permiso403 se prueban para cada decisión. GET/recarga no reenvían comandos.
- SQL exige por tenant3 resoluciones y3 eventos; correct-and-represent tiene
  exactamente1 sucesor prepared y return/exchange1 autorización cada una.
  Actor de las resoluciones=handover. Persistencia y recarga demostradas, no
  reinicio de PostgreSQL/PITR. Dos actas preparadas adicionales son fixtures;
  se completan/rechazan por repositorio real. No creador inicial de producto.
- Agenda: cuatro proyectos PASS,13.060s. Web105 PASS/1 SKIP en12 archivos y build
  PASS; Go test ./..., vet y build ./... sobre rebuilt PASS. Suites opt-in reales
  acreditadas por los logs conectados, no por tests omitidos del comando genérico.
- node24.14.1, pnpm11.19.0, Go1.26.7, Next16.3.2, Playwright1.62.1 fijados.
  Instalaciones offline/frozen-lockfile; gate anidado con --ignore-workspace.

Reproducción: componer FRANCHISE_COMPLETE_PACK_PLAN en destino ausente, instalar
la web con lock y el gate anidado con `pnpm install --ignore-workspace --offline
--frozen-lockfile`; ejecutar tests/build web. PostgreSQL descartable loopback con
nombre elite_confirmation_* y53 migraciones, nunca base de negocio. Definir
TEST_DATABASE_URL, ELITE_WEB_ROOT absoluto, ELITE_QUOTE_E2E=1, ELITE_ORDER_E2E=1,
ELITE_PAYMENT_E2E=1, ELITE_DELIVERY_READ_E2E=1, ELITE_DELIVERY_ACTION_E2E=1,
ELITE_OPERATOR_SECTIONS_E2E=1, ELITE_RESOLUTION_RECOVERY_E2E=1 y dejar vacío
ELITE_QUOTE_PROJECT; ejecutar `go test ./internal/platform/httpapi
-run '^TestQuoteAcceptanceBrowserPostgres$' -count=1 -v -timeout 15m`.
Para agenda, ELITE_CONFIRMATION_E2E=1 y TestAppointmentAgendaBrowserPostgres.
El runner gate.ps1 y logs preservados registran las invocaciones exactas locales.
Cluster propio detenido al terminar; datos/evidencia conservados.

## Fallos y límites

FAIL458: primera instalación del paquete anidado resolvió el workspace padre y
el CLI Playwright quedó ausente. resolution-candidate.log permanece FAIL, no
ejercitó el producto. Instalación aislada y gate posterior completo cierran esa
causa. FAIL457 sigue OPEN como familia: V303 corrige resolución de discrepancias;
asignación/estado de lead, emisión de cotización, recursos, disponibilidad,
checklist, recepción/disposición y otros comandos del sender genérico necesitan
su recuperación específica. No readmitir el sender completo por este PASS.

No prueba de todas las variantes de marker inválido, pérdida de storage después
de enviar, cierres de navegador, multi-tab ni todos los recibos malformados.
sessionStorage no es bóveda ni registro durable empresarial: servidor/CAS son
la autoridad. Consulta sin el registro requerido conserva incertidumbre y pide
revisión; no se interpreta un resultado vacío como permiso para reenviar.
Prepared/authorized no significan entrega física, dinero devuelto ni cambio
finalizado. No workers externos, fiscalidad, carga, nuevo SCA, deploy, rollback,
PITR, aceptación productiva, cierre T2804 integral ni promoción/ZIP.

| Archivo/log | SHA256 |
|---|---|
| src/app/franchise/page.tsx | f6e1121a6bbe33d98f6798aed24bb1f2e915d522cde9128eca01f090c7b0ef54 |
| src/components/franchise-command-panel.tsx | 65eda66afae243f8abf0dc9c25f5033620320db3c434401ce411dde87f896b3e |
| internal/platform/httpapi/franchisejourney_test.go | 159e47bad6ce2ef2de0e3840fa9509f866d48b4046d749864f58fbb6dfff076a |
| microsoft_playwright_browser_gate/tests/enterprise-web.spec.mjs | 0a8b5723a68baeaf43c79fadda75c6d2c5febde22c7d020d952334e4b70549f4 |
| browser-final.log | c260d9e04869fa3ae1c29a3fe28648466f03588dcb47680e393bf231cdda5884 |
| resolution-red.log | 8868b1733b4b7358e11adc60fa44ac59143dbd795d5c07c55b7ab4070e1e5be5 |
| resolution-candidate.log (FAIL tooling) | b7d415592abe64eead9af005b7fe550ae7e55c3252a65babb425743e05cee36b |
| resolution-candidate-r1.log | 5cd7d6d827bccb80af38bcb1f9cb7294bf0d6ab61218d0351e7be4e30f133c84 |
| agenda-final.log | 2118e0af425fbd0f21bb720a19478f24094fc1c9e05cdc51e9a5d61c648971ea |
| web-test.log | b09eaf38e79f38483671b430f5710328189ae6eee078ac6ff6ad2f4990004fa6 |
| web-build.log | b456e461962d7b6f23f2d84ad2ecdb946eab40600c6dfdc4c175a834ee2f89d0 |
| go-test.log | d7795c623f7f5246060adad6ce29c48c244ad0e02d72ee593cadbea1cd7a1d43 |
| go-vet.log / go-build.log (vacíos, exit0) | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
