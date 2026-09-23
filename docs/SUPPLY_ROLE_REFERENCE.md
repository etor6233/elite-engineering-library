# Suministro por rol — referencia de infraestructura

/supply conecta compras, fábrica, recepción y calidad con el owner J2 existente.
features.supply_portal es opt-in; requiere supply:read o supply:factory-read.
Las escrituras requieren el permiso explícito supply:plan, :factory, :receive,
:release o :inspect correspondiente, con tenant/organización y versión en Go.
El nombre de rol no concede permisos. La referencia del pedido se comparte
sólo con personas autorizadas. Maestros de proveedor/organizaciones/variantes
activas son precondiciones; el formulario crea la orden y plan de cantidades.

El alta usa Operations.CreatePurchaseOrder y BindPlan en una transacción,
SQL y reglas originales conservadas. ID de orden preasignado para recuperación,
eventos con UUID independiente. El lock transaccional por tenant/orden serializa
replay; una colisión del hash sólo serializa, no une IDs ni concede acceso.
El receipt planned existente liga actor/orden/comando/hash; no segundo ledger.
El importe y moneda son explícitos, sin política financiera nueva. Int64 se
transporta como texto. Admite32líneas y1000unidades; páginas de100series.

El formulario registra series y milestones, decisión humana distinta de la
solicitud, rechazo y reemplazo, ASN, recepción parcial en cuarentena, rechazo,
reinspección con evidencia nueva y liberación. El view proyecta el último
estado de revisión/requester para deshabilitar autoaprobación; Go y sus
constraints siguen siendo la autoridad. No cambia la máquina de estados.

El navegador guarda únicamente orden/comando/tipo/versión/hash en sessionStorage,
aislado por tenant/org/subject/vista. Un resultado incierto bloquea nuevas
escrituras hasta recuperar un receipt concordante mediante GET. Los archivos
de evidencia y motivos no se guardan en el marcador; se registra la huella,
no se afirma almacenar o clasificar documentos. BFF limita32KiB y tiempo,
autoriza antes de leer body y usa rutas backend fijas. Una respuesta409 incierta
se conserva para revisión autorizada; nunca se transforma en reintento automático.

La ayuda inline de esta versión explica el recorrido y su recuperación.
Nuevo nav.supply tiene texto español/inglés. La localización de todo el portal,
el CMS/general help y J5 siguen dentro de T2804; este delta no cierra todoT2804.

Pruebas nuevas: TestSerialSupplyCreateAtomic demuestra1nuevo+3replay concurrentes,
rollback de5tablas por variante inválida y fallo tardío de outbox, actor/hash/org,
lectura durable e int64 exacto. TestSupplyRoleCommandGoldens compara4comandos
Go/TS con Unicode/HTML/int64/orden de arrays. FuzzSupplyCreateBounds usa un
presupuesto finito de2s sobre source closure aislado. No nuevo runtime/migración.

TestSupplyRoleBrowser con ELITE_SUPPLY_ROLE_BROWSER=1, ELITE_WEB_ROOT y
ELITE_NODE_BIN absolutos, Next compilado--webpack y PAYMENT_CONNECTED_DB_URL
aislada usa Next/BFF/Go/PG reales y JWE/RS256/JWKS sintéticos. Sólo maestros
se siembran: sus23POST crean la orden/plan,3series(1rechazada),2stocks disponibles,
6solicitudes y6decisiones de calidad. Pérdidas de respuesta de alta y recepción
se recuperan porGET tras recargar, sinPOST adicional. Incluye permisos, foreign
org, autoaprobación deshabilitada, escritorio y390px sin desborde.

Seleccionar GO-CONNECTED-SERIAL-SUPPLY junto con su nuevo companion de alta
y TS-SERIAL-SUPPLY-PORTAL. El test Go de goldens usa el fixture del pack web;
el browser reutiliza el issuer sintético de GO-CONNECTED-CATALOG-AUTHORING,
seleccionado en la composición completa. El perfil web aislado prueba tipos,
no conectividad Go. Evidencias y reconstrucción en SUPPLY_ROLE_RELEASE_V402.

AUTHORED sólo composición/transporte/formularios/fixtures; sin nuevo upstream,
sin atribución corporativa. No aceptación productiva ni certificación de calidad.
ARCA penúltimo; Daybreak/libxml2 último diferido sin investigación.
