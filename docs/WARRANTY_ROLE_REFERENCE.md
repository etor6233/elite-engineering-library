# Garantía por rol — infraestructura de referencia

/warranty es opt-in features.warranty_portal. Conserva los owners de términos
vendidos, cobertura, reclamos, revisión humana, FIFO, calidad y conciliación.
El nombre de rol no concede permisos. Se exige warranty:read, warranty:self o
warranty:factory-read según vista, más el permiso concreto para cada escritura.
La consulta de cotización exige warranty:offer y entrega la versión actual,
importe exacto como texto y vigencia; BindOffer sigue bloqueando y validando
en PostgreSQL. RoleContext lee los pasos inmutables autorizados dentro de la
misma transacción repeatable-read de Claim, verificando hash y case_id.

El operador revisa términos configurados y los vincula a la cotización.
El cliente lee esos términos y registra recepción explícita. La entrega
aceptada activa la garantía vendida. Los reclamos requieren entrega y atención.
Diagnóstico→plan con repuestos/labor→revisión independiente→trabajo y emisión
FIFO→calidad distinta del técnico→aceptación del cliente→conciliación de fábrica.
El UI presenta las referencias verificadas del paso anterior; no exige copiar
payloads/hashes internos. Exclusión bloquea aprobación, revisión fallida exige
corrección, autoaprobación/autocalidad permanecen rechazadas por el owner.
Cancelación previa al trabajo y rechazo por exclusión usan las rutas originales.
La conciliación es un acuse interno, no un pago ni documento fiscal.

El BFF limita32KiB/tiempo, autentica y autoriza antes del body, fija rutas y
valida recibos contra tenant/org/actor/referencia/versión/hash. Command recovery
guarda sólo referencia y hash, aislados por tenant/org/subject/vista. Guarda
antes de POST y recupera por GET tras una respuesta incierta incluso con reload;
no reenvía la escritura. Ni archivos ni motivos se conservan en ese marcador.
La evidencia se representa por su huella; el portal no afirma almacenar o
clasificar documentos. Una referencia aún incierta permanece para revisión.

TestWarrantyRoleBrowser usa Next/BFF/Go/PG/Chromium y sesiones JWE/RS256/JWKS
sintéticas. Habilitar ELITE_WARRANTY_ROLE_BROWSER=1, ELITE_WEB_ROOT,
ELITE_NODE_BIN y PAYMENT_CONNECTED_DB_URL aislada elite_payment_connected_*.
HANDOVER_PROFILE_PYTHON fija el builder admitido. La primera fase realiza
términos/consentimiento en UI antes de que el fixture continúe los owners
quote/order/SDK/callback/reconciliación/checklist/handover/release ya admitidos.
Los owners de atención y stock preparan tres turnos y dos capas FIFO. La
segunda fase ejecuta activación y todos los pasos de tres casos desde forms.
18POST,15pasos,1caso cerrado/2cancelados,3unidades restantes, salida de inventario
-740.0000 y costo positivo740.0000 presentado en el servicio. Dos respuestas
perdidas recuperadas con GET y ceroPOST adicional. Calidad fallida/corrección,
exclusión, cancelación, foreign y otro cliente incluidos. Desktop/390px.

12goldens Go/TS verifican los comandos tipados, Unicode/interiores/HTML,
int64 fuera del rango seguro de JavaScript y arrays;18tests web, build/types.
No nuevo algoritmo, migración, runtime o dependencia. El fuzz de dominio y
guards de host/migración de WARRANTY_CONNECTED_RELEASE_V402 permanecen.
El nuevo fixture de browser comparte el issuer sintético del pack de catálogo
y los goldens son compartidos con el portal. Composición completa obligatoria
para la prueba conectada; perfil web aislado sólo prueba imports/tipos.

FAIL859 fue una expectativa positiva en el ledger de salida, después de que
ambas fases browser pasaron. Se conserva el RED y se corrige la expectativa;
la base ya guardada se verificó read-only sin repetir escrituras.
Ayuda inline1.0.0 en esta release; integración CMS/capacitación e i18n privados
continúan en T2804. AUTHORED sólo proyección/transporte/UI/fixtures; fuentes
BC/Go/PG y sus atribuciones exactas conservadas, sin atribución empresarial nueva.
No cierre de todo T2804, infraestructura global ni producción.
