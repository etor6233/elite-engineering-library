# V402 — auditoría local T2802/TEST02, T2804 y T2805

Fecha: 2026-09-11. Mantenimiento de biblioteca; revisión read-only. No se ejecutaron suites, proveedores, npm/pnpm, ni se modificó código canónico. Los hashes del anexo identifican el código observado; un hash de test existente no se presenta como una nueva ejecución PASS.

## Resultado que habilita el próximo trabajo

El siguiente vertical útil es **cotización aceptada → pedido/reserva → solicitud de pago → SDK oficial en servidor fixture → webhook autenticado → reconciliación por GET → preparación inicial de handover → checklist/aceptación → gate de liberación**, en el mismo destino/revisión. Reusar Commerce, FranchiseJourney, ProviderIntegration, outbox/jobs y el módulo oficial de pagos existente. No crear otra biblioteca ni otro ledger.

La referencia V400 tiene 828 archivos reconstruidos (71 packs, incluye el pack HTTP de observabilidad adicional a franquicia). El plan franquicia real contiene 70 entradas; su descripción inicial todavía dice 68, un drift documental a corregir en cierre. El módulo oficial de pagos **no está seleccionado** en ese plan.

## Matriz de aceptación local

| Superficie | Implementado y evidencia reusable | Código / wiring faltante | Prueba conectada faltante y cambio mínimo |
|---|---|---|---|
| Cotización → pedido → reserva | FranchiseJourney/Commerce + rutas BFF y Postgres reales; TestQuoteAcceptanceBrowserPostgres y V298 enlazan pedidos creados por cliente, reserva e intent. V297 verifica replay, 16 contendientes, outbox atómico, scopes y recuperación. | No hace falta duplicar cotización/pedido ni request UI. | Extender ese harness con el pago real al fixture y handover nacido por API, conservando sus IDs y la misma DB. Reejecutar sólo delta y negativos afectados. |
| Solicitud al proveedor | GO-OFFICIAL-PAYMENT-WEBHOOK-ADAPTERS 0.2.0 ya contiene StripeIntentClient.CreateIntent y MercadoPagoPaymentClient.CreatePayment, serializados por SDKs oficiales; sus tests usan transportes fixture. | El perfil principal no contiene official_payment_webhooks; cmd/electromobility-api sólo selecciona string PAYMENT_REQUEST_PROVIDER. Ningún consumer payment.requested llama al SDK. Resultado ProviderReference/ClientSecret no queda integrado. | Componer módulo aislado y módulo raíz mediante go.mod/replace o workspace gobernado; bridge con intención durable, claim/fence, IDs/importe/moneda de DB y secreto por referencia. Test de POST al fixture con cuerpo/idempotency exactos y pérdida de respuesta antes/después de efecto. |
| Callback de pagos | VerifyStripe y VerifyMercadoPago delegan firma/tolerancia/version a SDK oficial. ProviderIntegration.AcceptWebhook ya hace conexión durable activa, body-hash replay check e inbox+job atómicos. | No hay ruta Stripe/MP montada que convierta Event a Receipt, ni processor de provider.webhook.received para actualizar el pago. El genérico HMACSHA256 NO sustituye la firma del proveedor. | Reusar verifier → Receipt con tenant/conexión fijados en servidor → inbox → handler de jobs. Casos firma errónea, replay exacto, evento divergente, importe/moneda/account/order mismatch, desactivación concurrente y duplicado tras commit. |
| Reconciliación pagos | ReconciliationService y PG mantienen mappings y discrepancias matched/missing/amount/state. Commerce.TransitionPayment persiste CAS y outbox atómicos. | El módulo de SDK no tiene Retrieve/Get; reconciliación es registro, no observación oficial + comparación + transición. Callback exitoso no se traduce en evidencia monetaria por sí solo. | Añadir adaptación fina de GET oficial y binding comprobado al pago; nunca inventar captured ni corregir dinero por manual string. Simular entrega fuera de orden, pérdida POST, GET recovery y replay con una sola transición/evento. |
| Handover inicial | Tabla sales.delivery_handover; checklist publicado/completado; aceptación con serial y sujeto/hash server-side; rechazo/excepción/devolución. V136 y tests de recovery existentes. | No hay Create/PrepareHandover para pedido normal. Único INSERT de producto hallado es returnexchange.go para reemplazo. Los harness actuales insertan prepared por SQL fixture. No hay gate comercial demostrable inicial. | Añadir comando al owner FranchiseJourney/Commerce usando pedido/line/stock/customer derivados y comprobados bajo lock, identidad idempotente y outbox. No inferir crédito/entrega/regla fiscal: gate contractual externo con fixture de referencia explícito y fail-closed fuera de configuración admitida. |
| Frontend por rol (T2804) | Role-dashboard y visibility tests; admin/factory/customer views; comandos reales de cotización, agenda, pedido/stock/request pago, checklist/aceptación; ayuda versionada15/15. | Handover inicial requiere comando BFF/form/error/recovery. El inventario de links por rol no equivale a alta personas, factory ops, CMS o training/evaluación completas. | Priorizar la acción handover sobre vistas existentes y sus permisos, misma versión de help/training/support; cerrar explícitamente cada claim del mapa restante, sin contar nombres de módulos como equivalencia. |
| WhatsApp envío/status (T2805) | Python adapter oficial fijado, Sender con hash de proceso/adapter y aprobación; fence PostgreSQL; appointment HTTP real; raw callback firmado con inbox/job; StatusWorker/Host y reconciliación con recepción perdida. Tests concretos en anexo. | cmd/electromobility-api no monta receiver/host; no hay cmd de conversación/WhatsApp en referencia. App.NewProduction acepta assembly pero el test utiliza LLM/domain fixture y sender stub. | Crear composition host/config que instancie los owners existentes y los deje desactivados sin perfil/credenciales. Test de referencia con transporte Meta fixture, raw inbound→contact identity→runtime→respuesta→fence→status, más shutdown/restart. |
| WhatsApp conversación libre | Ingress guarda whatsapp.raw_webhook.received.v1; app/runtime/contact resolver existen. | Falta transformación durable de mensajes entrantes autenticados a channels.Message y dispatch en host; Sender existente sólo acepta JSON de template (recipient/template_name/language_code/body_parameters), no cualquier texto del LLM. | No conectar texto libre a sender template. Cerrar modo de respuesta admitido/policy, adaptar el envío oficial ya fijado al tipo permitido, enlazar approval/budget/handoff y probar negativa fuera de sesión/template permitido. |
| Mercado Libre mutations | GO-MERCADOLIBRE-QUESTION-OUTBOUND ya aporta POST aprobado y reconciliación GET-only mediante fence compartido; tests oficiales/adaptación fijada existentes. | No asumir todas las mutations marketplace implementadas. Frontera de preguntas es concreta; lifecycle publicaciones/órdenes requiere selección o exclusión explícita. | Reusar pregunta/respuesta en host si requerida; prohibir extender automáticamente Ads writes/marketplace por la palabra integración. |

## Callpaths exactos para pago

1. `POST /v1/commerce/orders/{id}/payment-request` → httpapi.commerceAPI.requestOrderPayment → commerce.Service.RequestOrderPayment → postgres.Commerce.RecordOrderPayment / RecordPaymentIntent. Produce payment.payment_attempt(state=created), outbox `payment.requested`, actor tomado del principal. Clave es tenant/provider/key; inicial es un intento del total, no parcial/segundo intento.
2. Nuevo bridge debe consumir ese evento sin volver a crear intención; carga estado durable y llama `officialpayments.StripeIntentClient.CreateIntent(ctx, StripeIntentRequest{AmountMinor, Currency, OrderID, IdempotencyKey})`. SDK Stripe-go86.3.0 es dependency pin oficial commit a2df585a800a97fe8ec4ebf551b4449bdb3d90a1; MercadoPago1.14.0 commit f910ee53fbb6819e435eaf3d0f800cb1fe74ae09. Revalidar source/license/SCA exactos antes de ejecutarlos: esta auditoría no renueva su admisión.
3. StripeResult devuelve Provider, ID, Status, ClientSecret. ClientSecret no entra en log/evidencia; el consumer cliente necesitaría endpoint autenticado/alcance si se elige flujo de confirmación. El intent por sí mismo no cobra; hace falta cerrar confirmación/checkout sandbox del flujo elegido. MercadoPago CreatePayment recibe token/método/email/notificationURL y requiere evitar incluir PII en outbox/evidencias.
4. Handler callback específico llama `VerifyStripe(payload, signatureHeader, signingSecret, tolerance)` o `VerifyMercadoPago(payload, signatureHeader, requestID, dataID, signingSecret, tolerance)`. Normalizar la diferencia actual ProviderCode `mercadopago` (Commerce) vs `mercado_pago` (adapter) en el único bridge y testearla, no cambiar IDs históricos silenciosamente.
5. `postgres.ProviderIntegration.AcceptWebhook` bloquea por SHARE conexión durable active y registra inbox+job. Nuevos jobs no pueden ejecutar transición desde strings sin vincular payment/provider-reference/account/currency/amount/order. GET oficial debe comparar esos datos.
6. El estado Commerce actual permite created→pending|authorized|failed; pending→authorized|failed; authorized→captured|failed; captured→refunded|disputed. Evento provider succeeded puede llegar sin authorized previo; requiere mapeo y transacción con evidencia que trate saltos/reordenamiento según contrato, sin sintetizar dos eventos financieros ficticios.
7. Handover normal no existe todavía. La función lockDeliveryScope explícitamente advierte que es integridad, NO proof of payment/release/delivery. Reusar sus checks, pero no confundirlos con la nueva autorización de preparación/liberación.

## Fallos de interpretación que evitar

- GO_OFFICIAL_PAYMENT_WEBHOOK_ADAPTERS claim0.2.0 incluye creación, pero Applicability conserva texto anterior que la rechaza: corregir contradicción junto al cambio, no inventar ausencia del cliente.
- SDK oficial + tests de transporte no equivalen a producto pago conectado. Debe demostrarse el consumer, callback y reconciliación.
- CONDITIONED sólo por credenciales es incorrecto hoy para pagos o montaje WhatsApp: faltan los bridges/hosts descritos.
- Hash de test source permite reutilizar alcance histórico sólo si también coinciden sus entradas/runtime; no es comprobante de ejecución nueva.
- Toda composición/glue nueva es AUTHORED inevitable o ADAPTED sólo cuando existe source derivado exacto; nunca atribuir negocio local al SDK oficial.
- No se inspeccionaron ARCA ni Daybreak. Permanecen en el orden separado del usuario.

## Primer incremento de código después de T2801

Prioridad: integrar el módulo oficial de pagos con consumer+GET+callback en destino limpio, y un test SDK→HTTP→PG sin UI. Después reutilizar TestQuoteAcceptanceBrowserPostgres para terminar la misma jornada y creación handover; luego BFF/UX. Esto reduce el riesgo de construir UI sobre un cobro inexistente y no repite pruebas de listados/roles ya cerradas.

## Anexo de archivos exactos observados

El archivo hermano `integration-audit-hashes.json` contiene rutas absolutas, bytes y SHA-256 de cada fuente/test/receipt documental. Los registros de V297/V298/V136 son históricos; las pruebas no se ejecutaron otra vez por esta auditoría.

| Scope | Archivo | SHA-256 |
|---|---|---|
| reference | `internal/commerce/service.go` | `32100fdc384284acef4db78d034121f86e5d290f48e88360245e6f9eab4cb983` |
| reference | `internal/platform/postgres/commerce.go` | `f50f64b51242bc71a5d92ccf87a7b2186440d865dc207ccecd8481ad7c00be4a` |
| reference | `internal/platform/httpapi/commerce.go` | `a72d15067ab28bf00ce7805bad2bc883ff739e674a5e8122becff09efea3018c` |
| reference | `internal/platform/postgres/commerce_integration_test.go` | `91e8c31cbf133f81044bcad51e7f382c5f21fff18303d99a1571635585663945` |
| reference | `internal/platform/httpapi/commerce_test.go` | `8c62000fb647c78217fe60b45ac92974e0006f937a6c450848c92d126919fcb4` |
| reference | `internal/platform/httpapi/franchisejourney_test.go` | `83b304013de77fbc41097704a4e59967c1b7b84a48825facba198c563cc39cef` |
| reference | `internal/platform/postgres/franchisejourney_integration_test.go` | `af93b84ec8f18f2b7afad402f399f23fde3865d8ab22d874f42a8da74767afe4` |
| reference | `internal/platform/httpapi/franchisejourney.go` | `2439a78ae334f0a62376963c8a627e52a4627c39e7926b3927fda361230964fe` |
| reference | `internal/platform/postgres/franchisejourney.go` | `7fdee6385a6fd7ed9be4991251d4cd03f97af61ca0092b89706618be80794e8e` |
| reference | `internal/platform/postgres/returnexchange.go` | `f6843aba7f8e1d63ce13ee5be8c8fc98cba99bf335397758df39e3c6fe4edc6f` |
| reference | `internal/providerintegration/service.go` | `2639ceb1e734d8615886c714051a9a8049c6dc5d49b072505a7571d856825be1` |
| reference | `internal/providerintegration/reconciliation.go` | `77f7edb91a51bcbb9f48331303a2e5972eba65d7bacd6270252e3d172166a93d` |
| reference | `internal/platform/postgres/providerintegration.go` | `82bfc01c9615f9c2b96115c1dcf7f261c75036ce67fd819d314f7832ad0a41bd` |
| reference | `internal/platform/postgres/providerintegration_integration_test.go` | `cb81f9228899be5b7f83b7d24c05b09a07dbf3a9562ae7ec252818a48d28d997` |
| reference | `internal/platform/workers/processors.go` | `7bff6315aae3f6c0a2b04253860c9056649ee6527208c5191b9f84a560c62314` |
| reference | `internal/platform/postgres/outbox.go` | `ac68c30fdc76b1b9d9fb2a3c6014a9d805604ffa3bf6625ee8b939b399d72328` |
| reference | `cmd/electromobility-api/main.go` | `65f4ab673ddb7f8a7ca95d22325ce4b13ad30302b28f898d6b9f167f7d1e625f` |
| reference | `internal/app/app.go` | `73396b4dbab245c66b8250a185ee2b39d58812241f10470b55203bc51b5def67` |
| reference | `internal/app/production.go` | `48c2628528b8f5f467d5827af88df5290aa2d564a25b71d5a4c4176b6b0acf4e` |
| reference | `internal/app/durable_delivery_e2e_test.go` | `cfd51aaf676c28e2624f15a35b4a234ff5cc496d97938023e2e65ffe4098cc34` |
| reference | `internal/app/contact_identity_e2e_test.go` | `52754b240ac2c89bc15b1b98154649fc7f466b3e2ade4611fb5c4525006aaadb` |
| reference | `internal/whatsappbridge/sender.go` | `094db800ac72af136449260d607b216f07741c9ed19715f0b219eb7cde3f3735` |
| reference | `internal/whatsappbridge/webhook_receiver.go` | `349ef32e52bcacdd216029cff466b976781ecf9ae83277a7b7fb016245a08724` |
| reference | `internal/whatsappbridge/webhook_receiver_test.go` | `c15b1a8352d9b75d8e35670d502b3e625ce945b4864759389e281bbfbb0e6c40` |
| reference | `internal/whatsappbridge/appointment_notification_test.go` | `32ae8990b26158e727bc990fca3109183c5d4d98c32dda33dcd99a5172283a51` |
| reference | `internal/whatsappbridge/status_host_test.go` | `aa1b0082cd2789d63180eff1374078a9f84ed117760e58f780910cd919e1b04d` |
| reference | `internal/whatsappbridge/status_worker.go` | `a94051c187923cb19a10807369c28f78dea7c3e676d2da5bd2659ec85e7b5ffa` |
| reference | `internal/outbounddelivery/mercadolibre_question.go` | `5c30e9937bad08db395d98243198740e17f186a1da63c64ad49065e35587c6f0` |
| reference | `internal/outbounddelivery/mercadolibre_question_test.go` | `3243d098d830c574f11f1495c6c107939ebe3b2ad2dcea0414778a845e05bb28` |
| reference | `src/components/role-dashboard.tsx` | `ad78a3223dda60dc0d8d58de291324b6f6e7c698479bbac717be77222a8511b0` |
| reference | `src/platform/roles/role-visibility.test.ts` | `c4e085922de8de753148bfd6a99127952032a967469b7a9170fad78f8d177f14` |
| reference | `src/platform/help/coverage.test.ts` | `3d5bd9aa236545599cca8b050017ebb0cd7125fa0107233cbdc40660d5f5a02d` |
| reference | `src/components/customer-quote-actions.tsx` | `452c0ebb920b5729022c7f9edd1a540dfc3b22c663dfccfe165290ea75f38611` |
| reference | `src/components/customer-handover-actions.tsx` | `c1f264d9a83d58f7c34976e8479b107ef8b1d472357123ebd97c9b87f380e5af` |
| reference | `src/app/factory/page.tsx` | `c2aa3463d897d42d66d89d38d4e1ebf9c2db50bf1a32efd4458a997964d3e679` |
| reference | `src/app/admin/page.tsx` | `cfa32d36822be1ead4378f0b3136e2181297374faba2b391dcebaf3f9458e3c2` |
| reference | `src/app/customer/handovers/page.tsx` | `836b56e83748ec6e0063f7c8bb2111801ae20b5f72cca6fcce65d9b9c585046d` |
| reference | `whatsapp_cloud/official-source.lock.json` | `044f5b73195a44afca1fb1be1ceb1a7ec9d5c15092dfe03bb491cb17b52453da` |
| reference | `whatsapp_cloud/whatsapp_cloud.py` | `6310abef1497b31139aed605a8a3d6239ee8cfe23a12fe643cf22a681d673a28` |
| canonical | `implementation_packs/GO_OFFICIAL_PAYMENT_WEBHOOK_ADAPTERS.md` | `8faafb118d1bf92445c6a4f0464f3e9678608c6c88a68f1fca0e69a1b09ac99a` |
| canonical | `markdown_system/PAYMENT_WEBHOOK_ADAPTERS_PACK_PLAN.md` | `fa4c35f697e21d9e08239db38d8050d1626d85d35e17f689c229b97b55dcb017` |
| canonical | `markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md` | `8bb5eaac51e36e0a6ce8811415f0c61f4e1a744d76db6f40ada7b3cb0f91ccb2` |
| canonical | `reconstruction_evidence/PAYMENT_REQUEST_PORTAL_V298.md` | `1096ee8996625a6fcb4b9cd704fb69613882d453157d7fe0b7e125f10f1d1bdf` |
| canonical | `reconstruction_evidence/PAYMENT_INTENT_RECOVERY_V297.md` | `fbeec2007681d5dbe38f028c6d2fba581100d70c747a0f7f76bbc670b1ac4778` |
| canonical | `reconstruction_evidence/FRANCHISE_HANDOVER_SERVER_EVIDENCE_2026-08-30_V136.md` | `ca9e632679b73f90f052102587cd97aec3571bb7eab731d197966ffa8c1648d2` |
| canonical | `reconstruction_evidence/WHATSAPP_HTTP_DURABLE_INGRESS_V274.md` | `175ace4fd6d7b5591b66c1f30cf96161441fb64670022067ad0c0ae5d6b84ad0` |
| canonical | `reconstruction_evidence/WHATSAPP_HOST_OPERATOR_HISTORY_V278.md` | `3754b95d1d5f6a57bbb895717002d0ce1f25135f0667350cfbd3b6eec38643a7` |
| canonical | `reconstruction_evidence/GO_APP_WIRING_2026-09-04_V231.md` | `b4f93e8cd4f5dd5c1ede7d045797caca0f22cf2c0aff879eba9f5781fca20dd7` |
