# V402 — seguimiento de fábrica desde navegador

El perfil materializa77packs/943archivos exactos fuera de la carpeta canónica. Este delta conecta el seguimiento operativo de una unidad existente con su owner; no cierra todos los recorridos de supply/fábrica/stock ni todo T2804.

La pantalla selecciona una unidad autorizada y llama al owner existente para planned→assembly. Una lectura exacta por tenant/organización/unidad usa la proyección actual de EnterpriseQuery. La escritura conserva operations.Service, SQL CAS y outbox anteriores sin modificaciones. En respuesta perdida, la UI bloquea un nuevo envío y consulta el estado. La prueba introduce un avance competidor hasta quality: la recuperación informa el estado observado, sin afirmar un recibo idempotente ni atribuir ese avance al primer actor.

PASS:18tests BFF, Next producción Webpack con TypeScript, Go compilación/vet, Chromium→BFF→Go→PostgreSQL. El comando Go de compilación ejecutó cero tests y no recibe crédito de pruebas. El navegador ejecutó un POST de transición y verificó los efectos durables reales; otra transición ocurrió por el owner existente. Se rechazaron rol de lectura, organización/tenant ajenos y estado obsoleto. Viewport390×844 sin overflow ni errores de página. PostgreSQL detenido. No efectos nuevos de stock, dinero o calidad productiva.

Procedencia:10archivos AUTHORED de proyección/UI/transporte/pruebas; GO_ENTERPRISE_QUERY_API0.2.0/11files y TYPESCRIPT_GO_API_WEB_BRIDGE0.6.0/64files. El grafo operativo anterior sigue siendo local y no se promociona por este wiring. No se agregó dependencia.

Corrección del inventario: sí hay alta acotada de recurso employee/contractor vinculado a principal_subject mediante CreateServiceResource, POST /franchise/resources y ResourceCreatePanel con recibo durable. Esto no equivale a provisioning de IdP ni alta general de personas. Capacitación/evaluación y onboarding persistente continúan abiertos; el checklist en memoria GO_ONBOARDING_CORE no los cubre.

## Receipts exactos

Stage `<LOCALAPPDATA>/Temp/elite-v402-library-infra`.

| Receipt | SHA-256 |
|---|---|
| `factory-browser-manifest.json` | `635261dc26a767e4e1e3f1492bc10ecfd7324baa24709e7b3ea53ce18c283063` |
| `factory-browser-package-result.json` | `2243e980c7b74053560a026f7a3afb9325a8486862461f9f931d52e9a52bde29` |
| `factory-browser-composition-stage.json` | `fbf8982c1e858543788b0a19b67a14bb55bc95691e208d569184816f2cb2fbbb` |
| `factory-browser-reconstruction.json` | `1dbb4e25c52430323729f8491125e15f297f22a672ac5839be2efed9ffe819bf` |
| `factory-browser-pg-final/result.json` | `759f464b5643bdb5e9c13bd465548c6707fb43bcc518bb213c5e94188a687f90` |
| `factory-role-inventory.md` | `43d3e26ab037ccdce06af217e1d8dbf07f48ac0fad6cb11fd0c1bb95b41a4ce2` |
