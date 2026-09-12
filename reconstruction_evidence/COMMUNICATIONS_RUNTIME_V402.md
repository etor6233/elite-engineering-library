# V402 — infraestructura conectada de comunicaciones

Revisión canónica:82packs/1032archivos únicos en franquicia;178packs/1797bloques en biblioteca. Reconstrucción exacta de los seis planes afectados y compilación de producto y paquetes de tests en franquicia/backend/serverless. El comando de compilación ejecuta cero tests; no recibe crédito de prueba funcional. Este incremento no declara READY global ni producción.

WhatsApp: ingreso autenticado a la bandeja y trabajo durables existentes, resolución actual de contacto por tenant/organización, ejecución del runtime seleccionado y propuesta inmutable. La aprobación humana exige un sujeto diferente, permiso y payload exacto; el envío usa el fence existente. Estados y recuperación reutilizan la misma observación durable. Un contacto con effective_at futuro ahora produce handoff sin modelo, dominio ni proveedor. Seis escenarios PostgreSQL nuevos quedaron verificados en ejecuciones focales, sin repetir suites anteriores.

Navegador real Chromium→BFF→Go con OIDC/JWE→PostgreSQL: tres propuestas, dos aprobaciones/envíos al fixture, un rechazo, dos estados entregados y recuperación del segundo envío con cero POST adicionales. Navegación y envío respetan permisos; organización ajena no ve propuestas. Capturas desktop/390px revisadas, sin errores ni desbordamiento. Siete tests TS previos y build Next se conservan con sus receipts; la mejora visual sólo repitió el navegador afectado.

Host: activa de forma explícita un único app/runtime, módulos de aprobación/agenda, webhook y worker. Desactivado no lee secretos; base sin pack opcional rechaza su activación. Perfiles/procesos/CA fijados por hash, secretos por archivos externos, informe TLS1.3 con autenticación mutua y parada ante pérdida del reporter. Constructor con Python real y perfil bloqueado, transporte TLS y configuración se verificaron. Constructor no ejecuta SQL ni llama proveedor/modelo; esa evidencia es distinta de los escenarios conectados. La retención/supervisión del collector y persistencia del presupuesto/modelo siguen en T2809/T2807.

Identidad de servicio: broker con OAuth client_credentials y verificadores oficiales ya fijados, renovación serializada, secreto reread por grant, RS256/JWKS y claims exactos. Prueba de host con12llamadores concurrentes al adaptador dominio/worker obtiene un solo grant real de fixture; rechaza perfil de otra organización, dos fuentes de identidad y loopback sin opt-in de referencia. El broker agrega38subcasos, TLS/rotación de claves y fuzz finito75749ejecuciones/3s. El ciclo refresh/logout del usuario de portal es otra implementación en curso; el broker no lo sustituye.

Publicación de Page: comando independiente cmd/social-publishing con solicitud/revisión/agenda/worker, SDK oficial Meta26.0.1, POST y GET de observación, DELETE aprobado aparte. PostgreSQL:POST4/GET7/DELETE1, aislamiento, expiración bajo lock, replay/recovery, presupuesto de reintentos, recibo perdido y downgrade que preserva evidencia. POST ambiguo sin ID y DELETE sin confirmación mantienen UNKNOWN; no se reenvían ni se inventa éxito. Las32salidas de sus tres packs incluyen el owner humano compartido una sola vez. La integración aporta26archivos nuevos frente al baseline; el texto de conteos anteriores en la evidencia de staging conserva su momento histórico.

Procedencia: SDKs oficiales DEPENDENCY_PIN con sus revisiones/archivos y licencias fijadas; adaptación Python declarada desde su fuente exacta; host/BFF/SQL/recibos son AUTHORED glue. No se atribuye lógica local a Meta, Microsoft u otra empresa. Se normalizaron15payloads CRLF a LF con hashes antes/después, AST Python idéntico, JavaScript emitido por Next/SWC idéntico y literales SQL preservados. Los manifests y logs originales permanecen intactos; el rebind de formato no se cuenta como nuevos tests.

Los perfiles independientes ahora seleccionan el fence y digest de contacto necesarios para pagos. El test exclusivo de omnicanal permanece en el perfil completo que selecciona leadstream; las pruebas de pagos no fueron quitadas. Backend40/554 y serverless57/779 compilan sus paquetes de producto y tests. Web8/162, HTTP83/1040 y WhatsAppsource2/107 se reconstruyen exactos; este último es un conjunto de fuentes/adaptadores, no otro host autónomo.

Condiciones restantes de este incremento: publicación incorporada; source/security/SCA de toda la revisión, T2807 y operaciones finales pendientes. Las cuentas/grants/profiles del usuario siguen externos. No se etiquetan esas tareas locales como falta de credenciales. ARCA y Daybreak no fueron investigados en este incremento.

## Recibos exactos

Stage: C:/Users/NL/AppData/Local/Temp/elite-v402-library-infra.

| Archivo | SHA-256 |
|---|---|
| `whatsapp-connected-final-manifest.json` | `db136d8d32e4c8c4e65ba444149fc17176459db45efb81c421f0167961702171` |
| `whatsapp-connected-final-receipt.json` | `4b787388c416430e259e9dfc253c31bf9ab4dee7161e08b40368555cb68c9069` |
| `whatsapp-browser-final-manifest.json` | `61cb936cd1925392b3c531d9c083cdcec98d7f03e3255bb4463b4648aa264c89` |
| `whatsapp-browser-final-receipt.json` | `b371e4a4a99fdb88e525cb1bd075b89f5d9dc537d5af1ef86b7308f15b0b149f` |
| `WHATSAPP_CONNECTED_ASSURANCE_V402.candidate.md` | `7210ce9254e5502720a1d4d51a2d334d23a00190676215e5d19f154aa1eb2cc7` |
| `whatsapp-payload-lf-rebind.json` | `9db2fb4f515f0985fd4466dd0d4b41b65d0f189af4ebafd4a414b16d64328128` |
| `connected-communications-host-gates.json` | `3a89fae8b67a1f097dcd552d7144319d078a95c145da9fb221b6048b54f080fd` |
| `whatsapp-host-broker-connected.log` | `77d7d67d47725cbf7ccac4ca75547d2c0801ef578e2e6b18241bb6a25e833f61` |
| `oidc-service-manifest.json` | `b40a1efedef643bc29d34dfa057e37f3ba98bc89dc5714c592ac473ecc4606c4` |
| `social-connected-manifest.json` | `964893452f0e74d619a07a6e4132db55364fbddb9739492060d92f17ef5c0b51` |
| `social-connected-evidence.md` | `b652784ea84c6022b4d8ad059f894d04c0676aa1a49615ec06a3ad539d307bfe` |
| `communications-package-result.json` | `bb74f702ae877e17c81b3c0cfe7492bc646ab421812129721c6ddb847671a5ec` |
| `communications-composition-stage.json` | `d355da38570168bc9f7cfa5ce19b63fd6838dbdf48ab7af7a8dcfc56c0f46329` |
| `communications-reconstruction.json` | `9c08438a62bfcfd32df138557cdfbaca77b0329274d16e1f3ea382d8926671c5` |
