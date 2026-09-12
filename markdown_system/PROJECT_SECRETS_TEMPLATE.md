# Preparar cuentas, configuración y secretos — guía operativa única

Revisión V295, 2026-09-07. Owner: T2801/T2803, con T2805–T2809 del roadmap vigente.
Guía AUTHORED contrastada con los consumidores de la composición; no es código
de Microsoft/Google/Meta ni certificación de esos proveedores. Aquí sólo se
registran identificadores, referencias y estados: **nunca valores secretos**.

## 1. Qué preparar y qué no contratar

No hay cloud, IdP, proveedor de cobros ni proveedor documental elegidos para un
proyecto productivo. Seleccionar un pack no selecciona una cuenta tuya ni autoriza
gastos. ARCA: **DIFERIDA_POR_USUARIO**, sin pedir claves ni trámites fiscales ahora.

| Orden | Preparación del usuario | Trabajo previo del agente |
|---|---|---|
| Ahora | Responsable, correo de recuperación, administrador suplente, MFA y códigos de rescate custodiados fuera del chat | Detectar NEW/EXISTING; preservar cuentas y operación existentes |
| Primer recorrido local | Proyecto, organizaciones y roles | PostgreSQL/migraciones, IdP compatible, URLs y claims exactos; cotización→pedido no requiere Ads/IA/cobro externo/ARCA |
| Antes de acceso público | Qué dominio/cuentas ya tenés y quién los administra | Elegir un target, región, residencia, presupuesto, gestor de secretos y despliegue; demostrar recuperación |
| Según selección | Canales/proveedores que usa el negocio | Admitir adapter y completar receta específica; no contratar todas las alternativas |
| Después de dominio/autorización | Autorizar sólo cuentas propias y finalidades acordadas | Callbacks reales, TLS, permisos mínimos, verificación y contrato |
| Diferible | Ads/reporting, push/SMS, SFTP, nuevas clases documentales, ARCA | Dejar proceso/ruta/worker deshabilitado y probar recorridos independientes |

Diccionario: **cuenta/trámite** = acceso administrativo del titular; **configuración**
= dato no secreto; **secreto** = contraseña/token; **clave interna** = bytes aleatorios
generados por herramienta; **identidad temporal** = permiso renovable sin clave
permanente. Issuer identifica emisor de tokens; audience identifica su API destinataria;
client ID identifica la aplicación; callback es la dirección de vuelta del navegador;
scope limita permisos; webhook es una notificación entrante. Un ID puede ser
confidencial aunque no sea contraseña.

## 2. Registro sin secretos y custodia

Una fila por acceso seleccionado. Estados: DEFERRED, DECISION_REQUIRED,
ACCESS_REQUIRED, STORED_NOT_VALIDATED, LOCAL_PROVEN, LIVE_PROVEN, ADAPTER_BLOCKED.
Completar campos no cambia un estado a PROVEN.

| Integración/entorno | Tipo | Cuenta/ID no secreto | Consumer exacto | Referencia protegida | Owner/recuperación | Caduca/rota | Estado/evidencia |
|---|---|---|---|---|---|---|---|
| ARCA | Diferida | No solicitar | ARCA_ENABLED=false; workers ausentes | No requerida ahora | Titular | Activación posterior | DEFERRED |
| Sólo selección real | Cuenta/config/secreto/clave/temporal | Sin tokens | Archivo + campo | Ruta/ID del gestor, no valor | Responsable | Fecha/condición | Pendiente |

### Local Windows: procedimiento disponible

Usar `markdown_system/project_local_secrets.ps1` desde la biblioteca. Set pregunta
con entrada oculta; Generate crea 32 bytes con CSPRNG .NET; Check no revela valores
ni hashes; Run entrega los nombres elegidos sólo al hijo. DPAPI cifra los archivos
en `%LOCALAPPDATA%/EliteEngineeringSecrets/<ProjectId>/dev/<NOMBRE>.clixml`,
fuera del proyecto, con DACL para usuario actual y SYSTEM. Un ProjectId diferente
por instalación. No es una bóveda productiva ni convierte .env en almacenamiento cifrado.

En terminal local sin compartir pantalla, transcripción ni grabación, el agente
sustituye las rutas y vos introducís el valor únicamente en el prompt oculto:

```powershell
# Abrir primero la terminal en la carpeta raíz de esta biblioteca.
$library = (Get-Location).Path
$project = 'C:\ruta\del\proyecto'
$tool = Join-Path $library 'markdown_system/project_local_secrets.ps1'
& $tool -Mode Generate -ProjectId 'mi-franquicia' -ProjectRoot $project -Names AUTH_SESSION_SECRET
& $tool -Mode Set -ProjectId 'mi-franquicia' -ProjectRoot $project -Names DATABASE_URL
& $tool -Mode Set -ProjectId 'mi-franquicia' -ProjectRoot $project -Names OIDC_CLIENT_SECRET
& $tool -Mode Check -ProjectId 'mi-franquicia' -ProjectRoot $project -Names DATABASE_URL,OIDC_CLIENT_SECRET,AUTH_SESSION_SECRET
```

AVAILABLE_NOT_PROVIDER_VALIDATED significa sólo descifrado local de un valor no
vacío, no validez/permisos/saldo. No mostrar el archivo ni descifrarlo en el chat.
Run requiere Executable absoluto revisado, Arguments sin secretos, Names limitados
al servicio y PublicConfigFile JSON sólo con campos públicos admitidos. API:
DATABASE_URL; web: AUTH_SESSION_SECRET y OIDC_CLIENT_SECRET. Ejemplo público API:
ARCA_ENABLED=false, OIDC_ISSUER, OIDC_AUDIENCE; no contraseñas en ese JSON.
El helper no hereda secretos/proxies arbitrarios y descarta stdout/stderr del hijo
para no filtrarlos. Health probes y observabilidad redactada van por separado:
no es supervisor ni solución de observabilidad. El secreto existe en memoria del
proceso autorizado; DPAPI no protege de malware/administradores/dumps.

Microsoft limita este cifrado al **mismo usuario y equipo Windows**. CLIXML no es
backup portable. Custodiar recuperación en gestor aprobado o reemitir en proveedor;
MFA y códigos de rescate separados. El helper rechaza sobrescrituras: rotar exige
detener consumidor, preservar versión anterior cifrada fuera del nombre activo,
cargar nueva con nombre exacto, probar y revocar anterior. El agente prepara rutas
exactas/backup; no borrar claves de identidad como si fueran tokens reemplazables.
[Microsoft Export-Clixml](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.utility/export-clixml?view=powershell-7.5).

### Despliegue: elegir una alternativa

| Target | Custodia/entrega | Datos pendientes y acceso mínimo |
|---|---|---|
| AWS | Secrets Manager + rol del workload/credenciales temporales | Account/region, ARN/versión, rol; lectura del secreto y KMS sólo si corresponde |
| Google Cloud | Secret Manager + identidad del workload; ADC donde el SDK lo soporta | Project/region, recurso/versión/identidad; Secret Accessor sobre secretos necesarios |
| Azure | Key Vault + identidad administrada; RBAC de datos separado de administración | Tenant/subscription/resource group/vault/versión; Key Vault Secrets User según operación |
| Otro/on-premise | Gestor y adapter de entrega por decidir | Autenticación, permisos, rotación, recuperación y probe; no reemplazar por .env |

La API lee entorno, **no resuelve vault://, kms:// ni ARN por sí sola**. El despliegue
resuelve/injecta sin volcar valores en manifests/logs. KMS no es automáticamente
una bóveda. DEV/PROD, migración/runtime/backup/deploy tienen accesos separados.
El target no elegido impide dar nombres reales/políticas IAM completas; se completan
aquí, no se inventan.
[AWS](https://docs.aws.amazon.com/secretsmanager/latest/userguide/best-practices.html),
[Google](https://docs.cloud.google.com/secret-manager/docs/best-practices),
[Azure](https://learn.microsoft.com/en-us/azure/key-vault/general/rbac-guide).

## 3. Primer recorrido: PostgreSQL, identidad y portal

Rutas relativas al proyecto materializado. Estos nombres sí tienen consumidores:

| Consumer | Nombre exacto/tipo | Obtención, comprobación y ausencia |
|---|---|---|
| cmd/electromobility-api/main.go | DATABASE_URL: secreto si contiene contraseña | Admin PostgreSQL entrega host/puerto/base/usuario/CA y contraseña protegida; probar conexión/migraciones completas con identidad limitada. No superusuario runtime ni sslmode=disable fuera del local aislado |
| Mismo | OIDC_ISSUER, OIDC_AUDIENCE: config | Discovery IdP y audiencia del access token de nuestra API; ausencia bloquea arranque |
| Mismo | HTTP_ADDRESS: config, default :8080 | Bind local/red privada; HTTPS por edge elegido |
| Mismo | ARCA_ENABLED: config, default desactivado | false ahora; true requiere ARCA_WSFE_SOCKET absoluto |
| cmd/electromobility-api/main.go → paymentRequestProvider → CommerceModule.PaymentProvider | PAYMENT_REQUEST_PROVIDER: configuración, no secreto. Ausente/vacío desactivado; sólo stripe o mercadopago | Selección del usuario antes de habilitar. V298 sólo registra intención local; no ejecuta cobro ni elige cuenta. Valores distintos impiden iniciar sin mostrarse en logs |
| src/platform/auth/oidc-client.ts | APP_BASE_URL, OIDC_ISSUER, OIDC_CLIENT_ID: config; OIDC_CLIENT_SECRET: secreto | Aplicación web confidencial, Authorization Code + PKCE; callback new URL('/api/auth/callback', APP_BASE_URL), en la raíz de ese origen; HTTP sólo localhost en desarrollo |
| src/platform/auth/session.ts | AUTH_SESSION_SECRET: clave interna | Generar 32 bytes/base64; no NEXT_PUBLIC. Cookie cifrada server-side; rotar invalida sesiones/flows, no hay keyring con solapamiento |
| src/platform/backend/public-client.ts y protected-client.ts | ENTERPRISE_API_BASE_URL, ENTERPRISE_TENANT_CODE, ENTERPRISE_ORGANIZATION_CODE: config | URL/códigos existentes, no IDs deducidos; probar público/cliente/ajeno/operador |
| src/platform/config/load.ts | BUSINESS_CONFIG_FILE: config | Nombre dentro de config/; business.example.json es ejemplo, no expediente aprobado |
| deploy/compose.yaml | POSTGRES_PASSWORD: secreto; OIDC_ISSUER/OIDC_AUDIENCE/API_PORT: config | Alternativa Docker local. PGHOST/PGPORT/PGDATABASE/PGUSER y POSTGRES_DB/USER en Compose; PGPASSWORD deriva del secreto. No imprimir compose config interpolado |

Workers seleccionados: return-accounting-worker usa DATABASE_URL y
RETURN_ACCOUNTING_WORKER_ID; return-exchange-worker usa RETURN_EXCHANGE_WORKER_ID;
return-effect-worker usa RETURN_EFFECT_WORKER_ID (todos con DATABASE_URL).
IDs de worker son configuración, no tokens. Los diferidos return-fiscal-worker,
arca-parameter-worker y arca-fiscal-worker consumen respectivamente
RETURN_FISCAL_WORKER_ID, FISCAL_PARAMETER_WORKER_ID y FISCAL_WORKER_ID; los dos últimos
también ARCA_WSFE_SOCKET. No arrancarlos ahora. Importadores meta-lead-import y
tiktok-lead-import usan DATABASE_URL y rutas de perfil/batch/receipt, no un token
permanente implícito. TEST_DATABASE_URL, ELITE_* de tests, NODE_ENV y SystemRoot
son controles del entorno de prueba/plataforma, no credenciales comerciales.
Los IDs y datos de fixtures no se trasladan al negocio real.

El módulo anidado **return_refund_worker/cmd/return-refund-worker/main.go** sí
consume STRIPE_SECRET_KEY y MERCADO_PAGO_ACCESS_TOKEN (con ese guion bajo exacto),
además de DATABASE_URL y REFUND_WORKER_ID. Esto no convierte esos nombres en loader
de cobros del main: habilitan providers de **devolución** en ese worker. No arrancar
el worker hasta elegir proveedor, registrar autorización de devoluciones, probar
reconciliación y revisar su logging de errores. No inyectar ambos por defecto.

Una vez elegido el IdP: identificar tenant/realm y administrador; comprobar
discovery/RS256/claims; registrar aplicación web confidencial y API; configurar
callback exacto; custodiar client secret; asignar usuarios por rol/organización;
probar login/expiración/usuario ajeno/logout. **El agente debe completar aquí el
enlace oficial y menú exacto del IdP elegido antes de pedirte navegar.** No están
demostrados indistintamente Auth0/Okta/Cognito/Entra por escribir OIDC.

Contrato `internal/platform/identity/oidc.go`: issuer/audience/RS256, `sub`,
`tenant_id`, `permissions[]`, `organization_ids[]`. Roles en business.json no
crean claims. El frontend requiere esos claims en el ID token y la API en el access
token: configurar y probar ambos. Login pide `openid profile email`, no parámetro audience específico
de proveedor: si se necesita, falta adaptación admitida. Logout en
`src/app/api/auth/logout/route.ts` borra sesión local, **no hace logout/revocación
global IdP**. No inventar callback de logout. Errores: callback mismatch→URL exacta;
401→issuer/audience/RS256/caducidad; 403→permiso/organización; sesión ausente→
HTTPS/cookie/secreto. No imprimir JWT para diagnosticar.

## 4. Claves internas y runtime conversacional

| Consumer | Input real | Límite |
|---|---|---|
| internal/platform/postgres/contact_identity.go NewContactIdentityStore(pool,hmacKey) | Clave independiente >=32 bytes; HMAC tenant/canal en internal/contactidentity/identity.go | No env canónica para el constructor. Rotar cambia digests: requiere migración/dual-read/recuperación, no reemplazo ciego |
| internal/app/config.go app.Config | LLMAPIKey/LLMModel/LLMBaseURL, DomainBaseURL/DomainTokenProvider, TenantID/TenantCode, ServiceKinds/ProductVariants/PriceBookID, ConversationStore/ContactResolver/ChannelRegistry/Conversation/LLMTokenBudget/ApprovalPolicy | LLM_API_KEY y otros MissingRequired son etiquetas de diagnóstico, **no loader os.Getenv**; objetos no se crean rellenando .env |
| internal/domainbind/config.go AccessTokenProvider | Fuente de token corto, adquisición/renovación del target | No JWT permanente; OrganizationID/LeadID legado ambos o ninguno, no un lead fijo para todo el negocio |
| internal/llmopenai/config.go y responses_tools.go | Config.APIKey/BaseURL/Model/Timeout | Responses/Chat Completions según llamada. Grok/u otro endpoint no compatibles por semejanza; modelo/tools/PII/evals/costo pendientes |
| internal/providerintegration/config.go | PROVIDER_WEBHOOK_CONNECTIONS_JSON: tenant_id/connection_id/provider_code/secret_env; secreto nombrado >=32 caracteres | Vacío = registro vacío. HMAC genérico no es firma Meta/TikTok/Stripe |
| portable_release_evidence_gate/build_release.ps1 | -SigningKey: ruta privada externa; profile: release_identity/allowed_signers_ref | OpenSSH Ed25519 según gate existente; passphrase/custodia/recuperación separadas. Pública distribuible; privada nunca ZIP |

Si se elige OpenAI: titular selecciona proyecto API DEV separado de PROD; identidad
de servicio/clave con scopes mínimos para llamadas usadas; registrar project ID y
referencia, inyectar Config.APIKey. Nunca clave admin de organización en runtime.
La API distingue alertas y `spend_limit` duro: verificar disponibilidad/enforcement
real, rate limits y presupuestos locales. Una alerta no corta consumo. Tests
sintéticos primero; inferencia/tools remotos pueden costar y requieren autorización.
Pérdida→revocar/reemitir; 401 identidad/proyecto, 403 scope/modelo, 429 cuota/límite:
reintentos acotados, no eludir cambiando cuenta.
[OpenAI proyectos/identidades/gasto](https://developers.openai.com/api/reference/typescript/resources/admin/subresources/organization/subresources/projects).

## 5. Canales, leads y publicidad

Copiar el profile del directorio indicado al expediente; flags sólo PROVEN con
evidencia. No ejecutar CLI live sólo para comprobar que hay un valor.

| Función/consumer | Config y secretos exactos | No habilita |
|---|---|---|
| whatsapp_cloud/whatsapp_cloud.py + provider-profile.template.json | business_account_id/phone_number_id/graph_api_version; WHATSAPP_ACCESS_TOKEN, META_APP_SECRET, WHATSAPP_VERIFY_TOKEN | Messenger/Instagram, Ads ni exportación universal de chats |
| cmd/meta-lead-webhook/main.go | DATABASE_URL, META_APP_SECRET, META_VERIFY_TOKEN; TENANT_ID/ORGANIZATION_ID/LISTEN_ADDR; /webhooks/meta-leads | Sólo señal durable; falta retrieval/promoción autorizada/consentimiento |
| meta_lead_reconciliation/fetch_leads.py + profile | META_APP_ID (ID), META_APP_SECRET/META_ACCESS_TOKEN; form_id/tenant_id/organization_id/allowlists | No permiso de envío ni entrenamiento automático con todo |
| meta_ads_reporting/run_insights.py + profile | META_APP_ID/META_APP_SECRET/META_ACCESS_TOKEN; ad_account_id/approved_fields | Reporting no publica anuncios ni captura leads |
| google_ads_reporting/run_reporting_query.py + profile | customer_id/login_customer_id/query_allowlist; auth SDK abajo | No crea campañas ni flujo de formularios |
| tiktok_ads_reporting/run_reporting.py + profile | TIKTOK_ACCESS_TOKEN; advertiser_id/approved_dimensions/approved_metrics | No publicación/mensajería |
| tiktok_lead_adapter/run_retrieval.py + profile | TIKTOK_ACCESS_TOKEN; TIKTOK_APP_ID/TIKTOK_APP_SECRET son referencias del perfil, no prueba de OAuth implementado | No certifica autenticación webhook ni live |

### WhatsApp y Meta

Titular confirma portfolio, app, WABA (cuenta WhatsApp Business), número y
administrador. Teléfono visible, phone number ID, WABA ID y app ID son diferentes.
Seguir alta/test de colección oficial Meta con recursos de prueba y destinatario
consentido; no migrar número operativo sin evaluar convivencia. Envío:
whatsapp_business_messaging; metadatos/plantillas/cuenta:
whatsapp_business_management cuando corresponda. No business_management por defecto.
Registrar duración real del token y activos; token temporal de test no es operación
permanente. Verify token generado por nosotros, app secret emitido por Meta.
WHATSAPP_VERIFY_TOKEN no es META_VERIFY_TOKEN del receptor Go de leads.
Webhook HTTPS/challenge/firma antes de suscribir; plantillas y consentimiento aparte.
Tarifa según mensaje/mercado: enviar no es probe sin efectos. Offline primero;
real consentido después. Permiso erróneo→app/token/activo, no admin global;
firma errónea→bytes crudos/secreto, no desactivar firma.
[Colección oficial Meta](https://www.postman.com/meta/whatsapp-business-platform/documentation/wlk6lh4/whatsapp-cloud-api).

Ads/leads: confirmar cuenta publicitaria, Page/form y acceso específico; profile
exige ads_read para reporting y aprobación de acceso a leads para retrieval.
Lectura pública de páginas Meta quedó bloqueada esta revisión: **menús y conjunto
completo de permisos de leads no certificados**. Completar receta desde documentación
oficial accesible en cuenta elegida antes de pedir autorización. No reemplazar por
permisos WhatsApp ni cambiar su gate a PROVEN.

### Google Ads y Merchant

Ads: entrar al [API Center](https://ads.google.com/aw/apicenter) de cuenta
administradora real, no anunciante ni manager de test; reutilizar token empresarial
si corresponde. Completar API Access con contacto atendido y presencia real; términos
sólo por titular. Registrar nivel concedido (Explorer/Test/otro), no asumir live.
OAuth del usuario con acceso a cuenta va aparte. Elegir **uno**: ADC mediante
GOOGLE_ADS_USE_APPLICATION_DEFAULT_CREDENTIALS; o GOOGLE_ADS_CLIENT_ID (config),
GOOGLE_ADS_CLIENT_SECRET y GOOGLE_ADS_REFRESH_TOKEN (secretos). Ambos necesitan
GOOGLE_ADS_DEVELOPER_TOKEN, GOOGLE_ADS_USE_PROTO_PLUS y, si aplica,
GOOGLE_ADS_LOGIN_CUSTOMER_ID. No heredar GOOGLE_ADS_CONFIGURATION_FILE_PATH: SDK
priorizaría su YAML. Revocar OAuth/controlar refresh; sin secretos en receipts.
Perfil/query offline primero, consulta acotada a test sólo autorizada.
[Alta](https://developers.google.com/google-ads/api/docs/api-policy/developer-token),
[variables oficiales SDK](https://developers.google.com/google-ads/api/docs/client-libs/python/configuration).

Merchant: cuenta Merchant y proyecto cloud diferentes. Administrador habilita
Merchant API y acceso de identidad en Merchant. Registrar merchant_account_id,
cloud_project_reference, data_source_name, adc_identity_reference en
`google_merchant_product_sync/provider-profile.template.json`.
`sync_product.py` usa ADC, no GOOGLE_MERCHANT_* inventadas. OAuth o service account
autorizada en cuenta; `gcloud auth application-default login` no es
`gcloud auth login`, y su archivo contiene credenciales. Preferir temporal,
revocar acceso cloud y Merchant ante pérdida. El CLI **inserta producto**: no usar
como probe inocuo. Dominio/comercio/políticas/aprobación de producto siguen gates;
no habilita campañas.
[Acceso Merchant](https://developers.google.com/merchant/api/guides/authorization/access-your-account),
[ADC local](https://docs.cloud.google.com/docs/authentication/set-up-adc-local-dev-environment).

### TikTok

Confirmar Business Center/advertiser, app API for Business y autorización por
reporting o leads. No intercambiar OAuth de creator por advertiser. Perfiles
exigen permisos/retención/cuotas/consentimiento. La página oficial de
[crear suscripción](https://business-api.tiktok.com/portal/docs/create-a-subscription/v1.3)
no devolvió contenido legible: **pasos exactos y autenticación entrante siguen
ACCESS_REQUIRED**. No suscribir ni exponer receiver por tener token. Reapertura:
doc oficial accesible de app elegida + contrato archivado + firma/retrieval/replay
probados. Mientras no ejecutar scheduler/reporting/retrieval; portal independiente continúa.

## 6. Cobros, marketplaces y canales adicionales

V298 permite preparar la solicitud **local** después de una selección explícita:
PAYMENT_REQUEST_PROVIDER se guarda con configuración no secreta del proceso Go,
separada por entorno; no es una API key ni requiere abrir cuenta para el probe
sintético. Dejar ausente desactiva requests/transiciones HTTP de pago y la UI,
sin desactivar la lectura de pedidos ni la reserva de stock. El frontend obtiene
la selección del GET autorizado; no tiene otra variable ni lista independiente.
No completar este valor por el usuario: stripe se usó sólo en los tests locales.
Al activarlo en un proyecto, verificar permiso payment:create, organización,
misma clave, response recovery y monto derivado del pedido. La pantalla guarda
la clave UUID en sessionStorage, no tokens/secretos; si no puede conservarla no
envía el comando. No borrar storage para forzar otro intento. No crear un consumer
de payment.requested hasta completar cuenta/sandbox/adapter/webhook/reconciliación
y autorización de efectos. Una intención local no habilita dinero ni entrega.

| Capacidad | Preparación disponible y brecha |
|---|---|
| Mercado Pago o Stripe | Elegir por entidad/país y cobro propio/plataforma. GO_OFFICIAL_PAYMENT_WEBHOOK_ADAPTERS ofrece constructores/verificadores con secreto como argumento, no loader STRIPE_SECRET_KEY/MERCADOPAGO_ACCESS_TOKEN. Wiring Commerce/cuenta/reconciliación pendiente; payment created no es cobro |
| Mercado Libre | mercadolibre_marketplace/provider-profile.template.json: application_id/seller_id/site_id; MERCADOLIBRE_ACCESS_TOKEN. Autorizar vendedor, permisos por operación/refresh/revocación. Perfil reconoce ausencia de preproducción global aislada. Questions no cubre publicar/vender/logística completa |
| Amazon SP-API, fuera de 67 base | PYTHON_AMAZON_SPAPI_CATALOG_ADAPTER: amazon_spapi_catalog/provider-profile.template.json; AMAZON_SP_API_CLIENT_ID, AMAZON_SP_API_CLIENT_SECRET, AMAZON_SP_API_REFRESH_TOKEN; región/marketplace/seller. Catálogo no habilita shipping/fulfillment/programas restringidos |
| Email/recepción, selección adicional | GO_AWS_ENTERPRISE_STORAGE_EMAIL_ADAPTERS: AWS default credential chain/región, no EMAIL_API_KEY. AWS_SES_IMMUTABLE_EMAIL_RECEIVER y SECURE_EMAIL_MIME_QUARANTINE_CORE conservan intake/cuarentena separados |
| SMS | Fuera de perfil integral: no pedir Twilio/SMS_ACCOUNT_SID universal. Seleccionar/admitir adapter, cuenta/país/remitente/consentimiento antes de receta |
| Push | GO_FIREBASE_CLOUD_MESSAGING_ADAPTER fuera del perfil: Firebase project/ADC/FCM y FIREBASE_TEST_DEVICE_TOKEN del profile. SendDryRun llama a Google; no prueba entrega/cliente/opt-out |
| Messenger/Instagram/publicación Ads/chats históricos | No se cierran con Meta token/reporting. Contrato de canal/exportación, permisos y adapter admitido; sin envíos retroactivos ni entrenamiento automático |

Mercado Pago: Developers → **Tus integraciones**, seleccionar app y credenciales
de prueba/producción. Public Key no es Access Token; Client ID/Secret para
autorización según integración. Access Token/Client Secret en bóveda; firma webhook
independiente y URL del receptor implementado. No compartir por chat. Cuentas/métodos
test y firmas primero; cobro/devolución/live requieren autorización. Comisiones y
activación según cuenta; pérdida→reemitir/revocar/reconciliar antes de repetir
cobro incierto.
[Credenciales oficiales](https://www.mercadopago.com.ar/developers/en/docs/credentials).

Stripe: Dashboard **API keys** en sandbox, clave restringida con permisos de
operaciones elegidas y secreto de webhook separado. Publicable puede ir al frontend;
restringida/secreta nunca. No live/Connect sin titular. Rotar desde API keys, no
mezclar objetos sandbox/live. Offline antes de sandbox; no certifica pago real.
[Claves oficiales](https://docs.stripe.com/keys).

Amazon: registrar app en Solution Provider Portal según propia/terceros y roles;
seller autoriza. No claves AWS raíz ni acceso PII restringido para catálogo.
Después consentimiento OAuth, refresh protegido/endpoints/región. Sandbox no es
operación real. [Registro](https://developer-docs.amazon.com/sp-api/docs/registering-your-application).
Mercado Libre: autenticación pública devolvió error esta revisión; receta de consola
pendiente de verificación, sin inventar menú/caducidad. Owner
`MERCADOLIBRE_MARKETPLACE_PACK_PLAN.md` conserva su contrato.

SES si elegido: identidad dominio/remitente en región; publicar valores DKIM/SPF/
MAIL FROM/MX emitidos para esa configuración, sin sobrescribir correo existente.
DMARC, recepción, rebotes/quejas y cuarentena conectados. Sandbox limita
destinatarios/cuotas; salida exige titular. SES → Account dashboard → Get set up →
Request production access: uso real, sitio/contacto, sólo al autorizar. Simulador/envío
no necesariamente gratuito.
[AWS SES](https://docs.aws.amazon.com/ses/latest/dg/request-production-access.html),
[Firebase Admin](https://firebase.google.com/docs/admin/setup).

## 7. Documentos, infraestructura y operación

OCR no es un único secreto. Elegir proveedor por clase/dataset autorizado/residencia.
Los expedientes documentales obligatorios siguen vigentes. Los 67 packs base **no
materializan todo el pipeline OCR**: incorporar plan elegido, no declararlo presente.

| Alternativa (no simultáneamente obligatoria) | Campos y condiciones |
|---|---|
| Google Document AI | GOOGLE_DOCUMENT_AI_RUNTIME: processor_version_resource exacto projects/.../locations/.../processors/.../processorVersions/..., schema_version por clase, ADC. IAM para operación/recurso/almacenamiento, región/API/billing/versión/corpus antes de process_document |
| Azure Content Understanding | AZURE_CONTENT_UNDERSTANDING_DOCUMENT_RUNTIME: --endpoint-env default CONTENTUNDERSTANDING_ENDPOINT, --key-env default CONTENTUNDERSTANDING_KEY. Recurso Foundry/analyzer admitido. Consumer usa AzureKeyCredential, no Managed Identity sin adaptación probada |
| AWS Textract | GO_AWS_TEXTRACT_DOCUMENT_RUNTIME: default AWS credential chain/identidad temporal, región, operación y S3/cifrado según modo; no STORAGE_KEY universal |
| Local | Plan local + recursos/corpus/evals; sin credencial cloud no significa precisión perfecta ni hardware gratuito |
| SFTP | MICROSOFT_AVM_SECURE_SFTP_INTAKE si elegido: endpoint/usuario/directorio/permisos/clave pública/fingerprint host. Privada protegida. Activar SFTP/storage tiene efectos y costo, no ejecutar para inventario |

[Azure recurso/endpoint](https://learn.microsoft.com/en-us/azure/ai-services/content-understanding/how-to/create-multi-service-resource),
[AWS Textract IAM](https://docs.aws.amazon.com/textract/latest/dg/security-iam.html).
El agente completa consola/permisos de opción elegida antes de solicitar acceso.
Procesar documentos consume servicio/trata datos: fixture autorizado después del
gate local, no factura privada improvisada. 401/403→identidad/roles; recurso ausente→
región/ID; campos ambiguos→corpus/revisión, no persistencia forzada.

Dominio/edge: titular/registrador, zone/account ID, DNS actual, hostnames,
origen/puertos/certificados/renovación. Si Cloudflare elegido: **My Profile → API
Tokens → Create Token**, recursos/permisos de zona/acción, no Global API Key.
Leer/verificar antes de autorizar cambios DNS/WAF; plan/reglas disponibles probados.
No CLOUD_PROJECT_ID universal.
[Cloudflare tokens](https://developers.cloudflare.com/fundamentals/api/get-started/create-token/).

`ops/otel/collector.yaml` consume OTEL_EXPORTER_OTLP_ENDPOINT (config) y
OTEL_EXPORTER_AUTHORIZATION (secreto exporter), no token usuario ni permiso para
registrar payloads. Backend/retención/redacción/alertas/receptores pendientes del
target. `ops/readiness/operational-readiness.example.json` y
`production_admission_gate/` son evidencia, no instalaciones. Backup: identidad
separada, destino cifrado, retención, WAL/PITR cuando corresponda, restore probado;
conservar acceso a claves históricas, no regenerarlas para recuperar. CI/registry/
deploy: identidad temporal, permisos push/deploy separados, firma/MFA/rollback;
sin workflows pagos ni GitHub obligatorio por esta guía.

## 8. Diferir y reactivar

ARCA false/ausente en composition root1.9.1 omite FiscalModule y socket. No desplegar
arca-parameter-worker, arca-fiscal-worker, return-fiscal-worker ni WSAA/WSFE sidecars.
Conservar schema/dependencias/expedientes; emisión/nota de crédito fiscal bloqueadas
aunque pedido/reserva funcione. Activación futura: titular + procedimiento fiscal
owner + credenciales + parámetros/POS + pruebas/reconciliación/autorización, recién
entonces ARCA_ENABLED=true.

`config/business.example.json`: enabled=false y credentialRefEnv lógicos
MERCADO_PAGO_CREDENTIAL_REF, AMAZON_SP_API_CREDENTIAL_REF, MERCADO_LIBRE_CREDENTIAL_REF,
GOOGLE_ADS_CREDENTIAL_REF, META_ADS_CREDENTIAL_REF. **No son interruptor global Go.**
No tokens dummy. No montar canal/receiver, no arrancar CLI/scheduler/worker,
no suscribir webhooks y no habilitar egress. Registro HMAC vacío/perfiles bloqueados
rechazan uso, pero no sustituyen despliegue explícito.

Reactivar en mismo owner: seleccionar/admitir → cuenta/autorización → custodia/
permisos → config exacta → migraciones necesarias → negativos/duplicados/caducidad
→ probe externo autorizado → journey/reconciliación → activación/monitor/rollback.
Sin cobros se puede cotizar/pedir, no confirmar pago; sin canal no notificación;
sin fiscal no emisión; sin OCR no autoingesta. Credencial válida no termina adapter.

## 9. Checklist y fallos

- [ ] NEW/EXISTING, opción por necesidad, DEV/PROD y titular confirmados.
- [ ] Inventario seleccionado/diferido/bloqueado con owner; no todas las alternativas obligatorias.
- [ ] MFA/recuperación/suplente; ningún valor real en chat/Markdown/capturas/ZIP/historial.
- [ ] Consumer, referencia, permisos y rotación/revocación; Config no presentado como loader.
- [ ] Ausente/incompleto/inválido/permisos rechazan sin filtrar; independientes probados.
- [ ] Costos aceptados antes de efectos; alertas no confundidas con límites.
- [ ] Callback/roles/cookies coinciden; logout local no presentado como global.
- [ ] Acceso no marca adapter/corpus/reconciliación/journey terminado.
- [ ] ARCA diferida con workers ausentes y dependencias visibles.
- [ ] Evidencia local/live separada; producción sólo con gates del target.

Errores: registrar código redactado, consumer/revisión/entorno en
PROJECT_FAILURE_LESSONS; fuente oficial/advisory, corrección canónica, regresión y
checkpoint. No desactivar firma/TLS/permisos para obtener PASS.
En otro proyecto reconciliar **sus** packs/config/IaC/cuentas, no aplicar a ciegas.
Evidencia: `reconstruction_evidence/CREDENTIAL_PREPARATION_V295.md`.
Recetas del target/IdP sin elegir y accesos oficiales bloqueados de Meta leads/
TikTok/Mercado Libre siguen pendientes: no ocultarlos como “sólo credenciales”.
Completarlos aquí y en owners vigentes; no nuevo roadmap.
