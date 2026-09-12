# V293 — conciliación de 21 cores y próximo cierre conectado

Fecha: 2026-09-07. Mantenimiento/discovery. Checkpoint de entrada: revisión18,
EVT-0018, resume PASS. Scope: T2801/T2802/T2804/T2805/T2809 existentes.
No expansión de producto, nuevo plan, stack, módulo, dependencia, cuenta o costo.
No se modificó código ejecutable. Esto resuelve selección/prioridad, no implementa
las funciones faltantes ni certifica una franquicia.

## Resultado de la revisión

Los 21 cores siguen fuera del perfil integral. La revisión cruza sus claims y
firmas con rutas, persistencia y tests/guías de los owners seleccionados:

- 12 presentan solapamiento parcial con owners existentes: no duplicar dominios.
- 8 no tienen equivalente funcional demostrado en lo revisado; conservar el gap.
- 1, observabilidad, permanece candidato aislado, con reparación focal V292.
- Cero equivalencias completas aprobadas por esta auditoría documental/estática.

“Sin equivalente demostrado” no significa búsqueda exhaustiva de todos los
archivos ni imposibilidad técnica. Se inspeccionó la selección de 67 packs,
manifests/claims de los 21 y los puntos de implementación listados abajo.
La ausencia de un nombre en rg no prueba por sí sola ausencia de una función.
No son 21 nuevas capacidades REQUIRED ni nuevas decisiones de negocio: respetar
el blueprint aprobado y los owners de las 48 superficies.

| Core fuera del perfil | Disposición | Reutilización comprobada en código/evidencia histórica | Brecha que NO cubre | Owner roadmap |
|---|---|---|---|---|
| GO-DASHBOARDS-CORE | Solapamiento parcial | Overview: pedidos/leads/stock/casos/envíos desde PostgreSQL; portales admin/customer/factory. | KPIs de puntos/reseñas/NPS y experiencia por rol no equivalentes. | T2802/T2804 |
| GO-FX-CORE | Sin equivalente demostrado | Monedas y unidades menores existentes; no prueba conversión. | Tasas versionadas, fuente, redondeo, vigencia y asientos FX durables. | T2802 |
| GO-GIFT-CARDS-CORE | Sin equivalente demostrado | Pagos/pedidos no son saldo de gift card. | Emisión/canje/reversión durable, permisos, checkout y conciliación. | T2802 |
| GO-HELP-CENTER-CORE | Solapamiento parcial | Ayuda versionada whatsapp-status-view/1.0.0 conectada a agenda y recuperación (V278). | No equivale a CMS/KB general, capacitación evaluada/progreso por rol o soporte integral. | T2804 |
| GO-I18N-CORE | Solapamiento parcial | defaultLocale → html lang; textos de interfaz inspeccionados en español. | Catálogo/fallback/plurales y journeys localizados, no sólo atributo lang. | T2804 |
| GO-LOYALTY-CORE | Sin equivalente demostrado | Pedidos/pagos no acreditan ledger de puntos. | Earn/burn/reversión durable conectado a compra y reglas aprobadas. | T2802 |
| GO-MARKETING-CORE | Solapamiento parcial | Ingreso de leads, reporting Ads y notificaciones con owners separados ya presentes. | Segmentación/drip, supresión/consentimiento, campaña→conversión y proveedor real. | T2805/T2807 |
| GO-OBSERVABILITY-CORE | Candidato aislado | Políticas/gates/alertas declaradas; logger candidato reparado V292 fuera del perfil. | Host/exportación/retención/alertas, histograma y carga; no insertar core para cerrar casilla. | T2809 |
| GO-ONBOARDING-CORE | Solapamiento parcial | Organizaciones/acuerdos con rutas y persistencia; sesión/permisos son otro owner. | Alta de persona, activación, práctica, evaluación/progreso/reentrenamiento por rol. | T2803/T2804 |
| GO-PAYROLL-CORE | Sin equivalente demostrado | Ledger empresarial no calcula ni liquida nómina. | Política laboral/fiscal, liquidación durable, recibos y aprobación; no sólo configuración de deducción. | T2802 |
| GO-POS-CORE | Solapamiento parcial | Pedidos, precio servidor, intentos de pago y fiscalidad tienen owners durables. | Caja/turnos/tender/vuelto/dispositivos/offline cuando requeridos y UX de mostrador. | T2802/T2804 |
| GO-PROMOTIONS-CORE | Solapamiento parcial | Price books activos y líneas server-side; reutilizar ese cálculo. | Cupones, límites de uso concurrente, acumulación, expiración y reversión ligados a pedido. | T2802 |
| GO-REFERRALS-CORE | Solapamiento parcial | Origen del lead/canal e identidad CRM reutilizables. | Código de referido, elegibilidad/antifraude y recompensa durable por conversión. | T2802/T2805 |
| GO-REMINDERS-CORE | Solapamiento parcial | Jobs PostgreSQL, worker/status notificaciones y fence ya existen; V275–V278. | Programación de cada recordatorio, cancelación/reagenda y envío permitido; no confundir aviso de estado con scheduler completo. | T2805/T2809 |
| GO-REVIEWS-CORE | Sin equivalente demostrado | Identidad de cliente no demuestra reseña moderada/publicada. | Compra verificada, moderación durable, publicación/frontend y protección antiabuso. | T2802/T2804 |
| GO-SEO-CORE | Solapamiento parcial | Web pública y title/description existentes; gate Lighthouse acotado. | Sitemap/robots/JSON-LD/canonical y publicación/indexación verificadas; no importar core en memoria. | T2804 |
| GO-SLO-CORE | Solapamiento parcial | SLI/objetivo/ventana/owner/runbook y reglas de alerta ya declarados. | Conectar señales y alertas al journey, inyectar fallo y demostrar respuesta; Budget en memoria no lo prueba. | T2809 |
| GO-SOCIAL-POSTING-CORE | Sin equivalente demostrado | Reporting read-only NO es publicación social. | Adapter oficial de escritura admitido, aprobación, cuotas, publicación/confirmación/reconciliación y revocación. | T2805 |
| GO-SURVEYS-CORE | Sin equivalente demostrado | Queries/reporting no son captura de respuestas NPS. | Consentimiento, invitación/respuesta durable, autorización y UX/reporting. | T2802/T2804 |
| GO-WAITLIST-CORE | Sin equivalente demostrado | Agenda/capacidad no demuestra FIFO ni promoción de lista de espera. | Cola durable, liberación/concurrencia, oferta/expiración y notificación sin sobrecupo. | T2802 |
| GO-WARRANTY-CLAIMS-CORE | Solapamiento parcial | Casos de servicio PostgreSQL, lecturas scope cliente/org y devoluciones ya existentes. | Elegibilidad/términos/estados específicos de garantía y recorrido completo de reclamo con recuperación. | T2802/T2804 |

## Owners exactos para evitar implementación paralela

- **GO-DASHBOARDS-CORE**: `implementation_packs/GO_ENTERPRISE_QUERY_API.md`, `implementation_packs/TYPESCRIPT_GO_API_WEB_BRIDGE.md`.
- **GO-FX-CORE**: `implementation_packs/GO_COMMERCE_PRICING_PAYMENT_API.md`, `implementation_packs/GO_ENTERPRISE_ACCOUNTING_LEDGER_API.md`.
- **GO-GIFT-CARDS-CORE**: `implementation_packs/GO_COMMERCE_PRICING_PAYMENT_API.md`.
- **GO-HELP-CENTER-CORE**: `implementation_packs/TYPESCRIPT_FRANCHISE_JOURNEY_PORTALS.md`, `implementation_packs/PYTHON_META_WHATSAPP_CLOUD_ADAPTER.md`.
- **GO-I18N-CORE**: `implementation_packs/TYPESCRIPT_GO_API_WEB_BRIDGE.md`.
- **GO-LOYALTY-CORE**: `implementation_packs/GO_COMMERCE_PRICING_PAYMENT_API.md`.
- **GO-MARKETING-CORE**: `implementation_packs/GO_OMNICHANNEL_LEAD_INGRESS.md`, `implementation_packs/PYTHON_META_ADS_REPORTING_ADAPTER.md`, `implementation_packs/PYTHON_META_WHATSAPP_CLOUD_ADAPTER.md`.
- **GO-OBSERVABILITY-CORE**: `implementation_packs/SECURE_OPERATIONS_DELIVERY_CORE.md`.
- **GO-ONBOARDING-CORE**: `implementation_packs/GO_FULFILLMENT_SERVICE_FRANCHISE_API.md`, `implementation_packs/TYPESCRIPT_OIDC_PORTAL_ADAPTER.md`.
- **GO-PAYROLL-CORE**: `implementation_packs/GO_ENTERPRISE_ACCOUNTING_LEDGER_API.md`.
- **GO-POS-CORE**: `implementation_packs/GO_COMMERCE_PRICING_PAYMENT_API.md`, `implementation_packs/GO_ARCA_FISCAL_ISSUANCE_API.md`.
- **GO-PROMOTIONS-CORE**: `implementation_packs/GO_COMMERCE_PRICING_PAYMENT_API.md`.
- **GO-REFERRALS-CORE**: `implementation_packs/GO_OMNICHANNEL_LEAD_INGRESS.md`, `implementation_packs/GO_ELECTROMOBILITY_PUBLIC_CRM_API.md`.
- **GO-REMINDERS-CORE**: `implementation_packs/GO_RELIABLE_ASYNC_WORKERS.md`, `implementation_packs/PYTHON_META_WHATSAPP_CLOUD_ADAPTER.md`, `implementation_packs/GO_PG_OUTBOUND_DELIVERY_FENCE.md`.
- **GO-REVIEWS-CORE**: `implementation_packs/GO_ELECTROMOBILITY_PUBLIC_CRM_API.md`.
- **GO-SEO-CORE**: `implementation_packs/TYPESCRIPT_GO_API_WEB_BRIDGE.md`, `implementation_packs/GOOGLE_LIGHTHOUSE_WEB_QUALITY_GATE.md`.
- **GO-SLO-CORE**: `implementation_packs/SECURE_OPERATIONS_DELIVERY_CORE.md`.
- **GO-SOCIAL-POSTING-CORE**: `implementation_packs/PYTHON_META_ADS_REPORTING_ADAPTER.md`, `implementation_packs/PYTHON_TIKTOK_ADS_REPORTING_ADAPTER.md`.
- **GO-SURVEYS-CORE**: `implementation_packs/GO_ENTERPRISE_QUERY_API.md`.
- **GO-WAITLIST-CORE**: `implementation_packs/GO_FRANCHISE_CUSTOMER_JOURNEY_API.md`.
- **GO-WARRANTY-CLAIMS-CORE**: `implementation_packs/GO_FULFILLMENT_SERVICE_FRANCHISE_API.md`, `implementation_packs/GO_ENTERPRISE_QUERY_API.md`, `implementation_packs/GO_FRANCHISE_CUSTOMER_JOURNEY_API.md`.

Anclas revisadas:

- GO_ENTERPRISE_QUERY_API: EnterpriseQuery.Overview consulta cinco agregados
  tenant/organización PostgreSQL; TestEnterpriseQueriesEnforceOrganizationCustomerAndCursor
  y rutas admin/customer/factory. No KPI de puntos/reseñas/NPS.
- GO_COMMERCE_PRICING_PAYMENT_API: TestCommercePriceOrderAllocationPaymentFlow,
  price books y rutas de líneas/place/allocations/payments. Price book no es cupón.
- GO_FULFILLMENT_SERVICE_FRANCHISE_API: OpenServiceCase/TransitionServiceCase,
  repository PostgreSQL y red/acuerdos; no se readmite el Store de warranty paralelo.
- GO_FRANCHISE_CUSTOMER_JOURNEY_API: AcceptQuote, sales.quotation_acceptance,
  transacción PostgreSQL y TestFranchiseJourneyPersistenceIsolationAndReplay.
- TYPESCRIPT_FRANCHISE_JOURNEY_PORTALS: customer-quote-actions.tsx,
  customer/quotes/page.tsx, command accept-quote→customer:self→ruta Go; ayuda
  docs/whatsapp-status-operations.md vinculada a whatsapp-status-view/1.0.0.
- TYPESCRIPT_GO_API_WEB_BRIDGE: layout con metadata y defaultLocale, textos
  españoles; no equivalencia de i18n completo ni sitemap/robots/JSON-LD.
- SECURE_OPERATIONS_DELIVERY_CORE: contratos SLO y reglas
  PlatformHighErrorRate/PlatformOutboxBacklog/PlatformBackupStale. Declarar regla
  no demuestra que hay un exporter, destino o supervisor funcionando.
- MICROSOFT_PLAYWRIGHT_BROWSER_GATE: tests actuales de lead, cita, confirmación
  en agenda e historial/ayuda WhatsApp. No contiene un test navegador de
  accept-quote→pedido durable. Esa discontinuidad no exige crear otro POS.

## Evidencia histórica reutilizada y alcance

V260 demuestra navegador→BFF→Go→PostgreSQL para lead, cuatro proyectos browser,
replay/divergencia y efectos durables. Sus versiones son históricas: no se
presenta como reejecución actual ni como prueba de login OIDC.
V278 demuestra consulta de estados→sesión→RS256/JWKS→scope→PostgreSQL→ayuda con
fixture sintético y recuperación. No demuestra envío live, capacitación de toda
la empresa o importación de chats históricos. V292 cierra el defecto local del
logger; no host/exportación ni observabilidad del journey completo.

No se repitieron esos tests porque no cambió código en V293. Se preservan los
receipts previos y se reejecutarán los gates afectados por el siguiente delta.

## Próximo camino de cierre, dentro del roadmap

Primero T2802 + T2804: **cotización presentada → cliente autorizado la acepta
desde el portal → pedido durable único → respuesta comprensible → recuperación**.
Razón: ya existen frontend/BFF/identidad/dominio/PostgreSQL y sus propietarios,
pero falta demostrarlos como un recorrido común de navegador. Conecta el lead
y la agenda previamente probados con commerce; habilita después stock/pago/entrega.
No agregar GO-POS-CORE ni otro ledger para conseguirlo.

Criterios de la prueba pendiente (no declarados PASS aquí):

1. Tenant/cotización/cliente sintéticos en PostgreSQL descartable; sesión validada
   por el harness OIDC/JWKS existente, no permisos ficticios en producto.
2. Navegación real al portal; aceptar sin introducir IDs/precios/identidad del
   cliente manualmente. Comprobar permiso, pertenencia, versión y expiración.
3. UI→BFF→Go→transacción: mismo pedido en retry, un registro de aceptación y
   eventos/auditoría esperados; doble clic y dos solicitudes concurrentes.
4. Perder respuesta tras commit, recargar/reiniciar proceso y recuperar el estado
   sin pedido extra ni mensaje de éxito falso.
5. Cliente/organización incorrectos, versión obsoleta y sesión ausente: rechazo
   seguro y error comprensible; no ampliar permisos para pasar.
6. Ayuda/soporte contextual de esa versión, estados carga/vacío/error, teclado,
   responsive y navegadores representativos. No sustituir por Lighthouse del home.
7. Emitir evidencia del mismo árbol; devolver fixes al owner, reconstruir y
   repetir sólo gates afectados. No cambiar reglas comerciales sin autorización.

T2809 sigue requerido y NO se cancela: el histograma genérico aislado no desplaza
indefinidamente el cierre comercial. Instrumentar el recorrido con el owner
admitido que corresponda; no incorporar el logger CANDIDATE ni fingir exporter.
T2801/assurance/readiness sigue bloqueado y restringe promoción/expansión de
producto; probes y tests correctivos aislados no otorgan READY_TO_BUILD.
Para gaps REQUIRED sin implementación admitida, ejecutar CAPABILITY_GAP_RESOLUTION_GATE,
sin reemplazar el expediente por este mapa ni crear una nueva solución al azar.

Ensayo de la meta “menos de una semana”: T2810 conserva prueba desde carpeta
vacía, tiempos técnicos separados de decisiones humanas/esperas externas y
escenarios de proyecto EXISTING. No hay duración ni porcentaje sustentados aún.

## Fuentes de método oficiales, consultadas 2026-09-07

- [Google SRE — Reliable Product Launches](https://sre.google/sre-book/reliable-product-launches/):
  convergencia sobre infraestructura común y controles concretos/proporcionales.
  Sustenta evitar duplicación; no valida nuestros módulos ni promete plazos.
- [Microsoft Playwright — Best practices](https://playwright.dev/docs/best-practices):
  comportamiento visible, aislamiento y datos controlados.
- [Microsoft Playwright — API testing](https://playwright.dev/docs/api-testing):
  combinar experiencia de navegador con comprobación de efectos vía API.

Esta conciliación y sus decisiones son AUTHORED. No se adquirió código nuevo,
no se cambió pin/licencia, no se equipara documentación oficial con aprobación
de código local y no se promociona ningún pack.

## Verificación de identidad

Se comprobó: 21 IDs únicos, 21 no seleccionados, perfil67 y todas las referencias
alternativas seleccionadas. Los hashes siguientes fijan exactamente qué se
revisó, no acreditan equivalencia funcional. Staging:
Temp/elite-v293-6cf286442a9f4fb5911499420ec7820a. El mapa anterior se conserva allí
y en la evidencia V280; se corrige su resumen “élite/verificado” no sustentado.

| Archivo | SHA-256 |
|---|---|
| implementation_packs/GO_ARCA_FISCAL_ISSUANCE_API.md | 2d8dd3e13e5ae38986cb7232cf345adb70778f94e968d83b3155bd8f21d2a83d |
| implementation_packs/GO_COMMERCE_PRICING_PAYMENT_API.md | 2a8732e4083ca474ac75f2666437142d28f5d72a176494afa0ed173b612ae405 |
| implementation_packs/GO_DASHBOARDS_CORE.md | 6ef647ef1e3084a717a6f00521663a64c96ffa70caa37cccad94cffd1f418085 |
| implementation_packs/GO_ELECTROMOBILITY_PUBLIC_CRM_API.md | 69b0f715eab64ef831f510505ceb4b59d2a36444c7e2c557524eb44d424e6557 |
| implementation_packs/GO_ENTERPRISE_ACCOUNTING_LEDGER_API.md | d6a788b518904c56b6c2c0a35c745827964f9d3caafa9567719784dec36b9daa |
| implementation_packs/GO_ENTERPRISE_QUERY_API.md | 84601b5f6b3917f4f540ffb90d9d569b6cf3e9a8e0fa4955a736c2c4e9a38f70 |
| implementation_packs/GO_FRANCHISE_CUSTOMER_JOURNEY_API.md | f74b8e0758c051bccaa8ea7c4b7dd8e1c3af11dda5ee54f728e3e2edae6b856d |
| implementation_packs/GO_FULFILLMENT_SERVICE_FRANCHISE_API.md | b647b19aa840b1f6723fa9e349425f8d042084579da6cda4cace91c28bb021a6 |
| implementation_packs/GO_FX_CORE.md | 5086a042d522bd756a12a315412a763e9943ade3f88ed408205532ce03731de8 |
| implementation_packs/GO_GIFT_CARDS_CORE.md | 0ce8666f978d918c076652ec11a99fb1ba2e0477ac44ab4ce3c96b8b0490ad2d |
| implementation_packs/GO_HELP_CENTER_CORE.md | 20fc7f9ab1ea1e41102c7fc16bd1a352b95df2b561194d1ff9abbab343b11216 |
| implementation_packs/GO_I18N_CORE.md | 6ca21bafc4a82977a169fb2b60240f2a46e5509a3cb8eef1924293f810e563eb |
| implementation_packs/GO_LOYALTY_CORE.md | b2f3e07c4b127be207351b82d77be2a977c708b0aaed9aa200ed2103699cbcc8 |
| implementation_packs/GO_MARKETING_CORE.md | a9122f7daeed08463f6ab5035dab94150dede1ab50798c000f682bfc29c350d8 |
| implementation_packs/GO_OBSERVABILITY_CORE.md | 82bf34505c3fbdbdb9359f659f40b3f8bc32089a51e14495dd1f9260508d8c03 |
| implementation_packs/GO_OMNICHANNEL_LEAD_INGRESS.md | 5ff3325a3d1e2e33f79f80034b92fd771e97c104416e96d67e902b64c8ebc242 |
| implementation_packs/GO_ONBOARDING_CORE.md | 0348f51991a57109678e608451fb5e2cd39c12b19c2eef6fbdb30a9b5315aca5 |
| implementation_packs/GO_PAYROLL_CORE.md | b23d1ca807194768b2bf095bde037e186338c07dc528fe83993b8d7d30e3804e |
| implementation_packs/GO_PG_OUTBOUND_DELIVERY_FENCE.md | a3084eec92ce1a918c9240f815999724daf44d010fa53f45173e85173367b061 |
| implementation_packs/GO_POS_CORE.md | 1b470a7c7a44389a34b3f4879c3492249e79ca3c681845d8175a177bb486f017 |
| implementation_packs/GO_PROMOTIONS_CORE.md | ef721b1c53b55fbc39ee251a6e15184350ca01bf0dfef105aa2bc043eea91780 |
| implementation_packs/GO_REFERRALS_CORE.md | 4397e3aeb0236f4ab36f83dd910e016cdda24fef905580c14999d39d7049f28e |
| implementation_packs/GO_RELIABLE_ASYNC_WORKERS.md | 45b29df774959f6cf2b745bbb9a55a5829341bb3363619243106d14dcf565a08 |
| implementation_packs/GO_REMINDERS_CORE.md | efe16e3d13cfa7d109d2783d665f33b9e61d40c861a13d186beaef07e70dc0cb |
| implementation_packs/GO_REVIEWS_CORE.md | 80456619b92b15ee479b57fe760c81d2cb5d5d280be3ae9e7ccf6fa599c16b6e |
| implementation_packs/GO_SEO_CORE.md | 9d6d0c2cd1633d47c8562a289f1a578ed2a3e4b95c936835ada502024b35ec49 |
| implementation_packs/GO_SLO_CORE.md | a1ae0b2d1dd61981a137b93c8ac6ef843d5813e9ace7a27f915dd58c3cfe8a9d |
| implementation_packs/GO_SOCIAL_POSTING_CORE.md | 04791bba47e6173fb72098c8f690f8ffa62b7c2535ecd677da40b77e71e15330 |
| implementation_packs/GO_SURVEYS_CORE.md | e5786433b98612a5ff24b86336d610fe8818ce98dc9d84bb877bfd9a31565b9e |
| implementation_packs/GO_WAITLIST_CORE.md | b8b4b5768f78181bfa9a4a91d6900b5d9b1e3ec9ce713f5a58d87d34c333d6a5 |
| implementation_packs/GO_WARRANTY_CLAIMS_CORE.md | 5633258019a207a5a6317e90123b77c20f1970e52331d8f6a0a7bfa368d2b723 |
| implementation_packs/GOOGLE_LIGHTHOUSE_WEB_QUALITY_GATE.md | c805b9e01c9ffc0dc79396281a9df2e1b134c77662c9414e4251d4cf52b6d532 |
| implementation_packs/MICROSOFT_PLAYWRIGHT_BROWSER_GATE.md | 5e7c84835c0aa31b82d0d0567e5187839dbf01c00c2aedd9d29176fa2e1a8f46 |
| implementation_packs/PYTHON_META_ADS_REPORTING_ADAPTER.md | 2ed2981302f3c7c6caf69eaead0205cd9cf2e54b6799b719664241cb6339e53f |
| implementation_packs/PYTHON_META_WHATSAPP_CLOUD_ADAPTER.md | 1ea7bd10ca1cf2e989071b25356b160a87a4211ffb02fc0dfa6dae1697691101 |
| implementation_packs/PYTHON_TIKTOK_ADS_REPORTING_ADAPTER.md | 8804c81cfe58e25a062830f5f75f8f664b6bb1b097565f5d3201cbf01056b787 |
| implementation_packs/SECURE_OPERATIONS_DELIVERY_CORE.md | f9fbd06b522a7f4ea9248fe2fee29b3a23a1c8aef5f1d2c78bd2ed33dcf6c8a5 |
| implementation_packs/TYPESCRIPT_FRANCHISE_JOURNEY_PORTALS.md | 27c328d631c8bf9d67520aff73842919bb48703084e803b80b77d7fa30e3ab1b |
| implementation_packs/TYPESCRIPT_GO_API_WEB_BRIDGE.md | 8f23d48b1fd4f463f4785e76d86bace3c8d14b20f5f1e33ac49e5c0d718186cc |
| implementation_packs/TYPESCRIPT_OIDC_PORTAL_ADAPTER.md | 5e39c0353ae7414ff56ca9d4de1d81562cea4bd6d9dd9513a5e67eed1e5524e3 |
| markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md | 5b359cda20724e38edaea486ea0948c433cb9548dc3ec22d3b516091308adaca |
