# V402 WhatsApp connected host contract

Version: V402. This supplements the existing whatsapp_cloud/provider-profile.template.json and docs/whatsapp-status-operations.md. It is infrastructure configuration, not an account proof. Values remain absent until supplied/configured by the user; no secrets are embedded in profile artifacts.

## Exact composition

1. Load a bounded, exact-SHA provider profile from a private configured path. Python validate_profile remains normative. Keep all template/readiness/source fields, exact source_commit de70ee908a67026e642aaee3703d20464e2a9466, Graph version/WABA/phone and explicit approved_templates. Template default values are blocked; fixture profile() in whatsapp_cloud/test_whatsapp_cloud.py is synthetic only. No automatic PROVEN defaults.
2. Resolve the same server-side HMAC key for postgres.NewContactIdentityStore, postgres.NewOutboundDeliveryStore(pool,key,time.Minute), NewPostgresAppointmentApprovals and NewStatusRouter. Scope is fixed tenant+organization+connection; require durable integration.provider_connection with provider_code meta-whatsapp and that organization, active.
3. base = whatsappbridge.NewPostgresAppointmentApprovals(pool,key,purpose,policy). Purpose/policy are explicit configured consent contract values, not invented business defaults. replies = NewPostgresReplyApprovals(base,tenant,organization,connection,profileJSON).
4. Process holds PythonExecutable/PythonSHA256, AdapterDirectory/AdapterSHA256 (whatsapp_cloud.py), EvidenceDirectory. ReconcilerSHA256 hashes status_reconciliation.py independently. Hashes must match the materialized revision; directories private. Existing Process launches -I -B and no inherited secrets/proxies/PYTHONPATH.
5. Sender{TenantID,Profile,Tokens,Process}. The token source is runtime-only; other sources implement WhatsAppAppSecret(ctx) and WhatsAppVerifyToken(ctx). NewReplyModule(replies,sender,outboundStore) forces the same approval resolver and registered durable channel. Register(mux,real OIDC verifier). Exports VerifiedInboxChannel{} for the receiver side of an outbounddelivery.Channel or NewAppointmentNotificationModule; it refuses direct receive/send because inbox+job owns ingress.
6. NewWebhookReceiver(WebhookReceiverConfig{TenantID,ConnectionID,RetentionApprovalSHA256,Profile,Process,Secrets,Verification,Store:postgres.NewProviderIntegration(pool),MaxConcurrent}) mounts raw signed /messages ingress. The exact existing callback path is host-selected. Ingress returns receipt only, not a conversation response.
7. Construct the existing app.New(Config) once. OrganizationID and LeadID legacy defaults BOTH EMPTY: contact resolver provides scope. ContactResolver=whatsappbridge.ScopedContactResolver{TenantID,OrganizationID,Resolver:postgresContactStore}. Existing ConversationStore, LLM API/model/base, domain token source/base, variants/services/pricebook, finite token budget, approved retention and handoff text remain explicit configuration. Register a durable WhatsApp channel, but do not execute Dispatcher on webhook replies. The admitted Runtime.Handle is the only callable conversation entry here.
8. StatusObserver{TenantID,Profile,Process,ReconcilerSHA256,Secrets,Approvals:base}; router=NewStatusRouter(observer,connection,retentionSHA,key,8); router.EnableConversation(ConversationRoute{OrganizationID,Runtime:app.Conversation,Proposals:replies}). NewStatusWorker(router,workerID,2*time.Minute,retry). Run(ctx,StatusPrincipalSource,StatusReporter,interval). One existing provider-events queue processes mixed messages/statuses. No parallel status-only consumer.
9. Service identity is revalidated for each poll using actual OIDC/service auth; requires appointment:manage (existing status owner), whatsapp:process and configured organization. A provider webhook never creates a principal. StatusReporter must use the existing bounded exclusive authenticated connection. Failure to authenticate or report stops/holds processing; no invented subject/permissions.

## Human HTTP contract

Bearer-only same verifier/tenant/org, no cookie/body identity. All responses no-store. Human review permission whatsapp:approve; sends/recovery additionally whatsapp:send. Service requester differs from human reviewer through the common approval owner.

- GET /v1/franchise/whatsapp/replies returns at most 50 exact pending/completed request projections, newest first.
- GET /v1/franchise/whatsapp/replies/{request_id} returns request_id,state,payload_sha256,context (exact message text/recipient/service-window/scope).
- POST /.../{request_id}/decision body {payload_sha256,approved,reason}. Approve/reject immutable; no message text/contact fields accepted. Approval itself sends nothing.
- POST /.../{request_id}/send body {payload_sha256}. Reads current bindings/consent/window and uses the existing fence. Accepted does not mean delivered. Retry retains exact identity and never creates a second POST after unknown.
- POST /.../{request_id}/recover body {payload_sha256}. Reads bounded existing local SEND_RECEIPT/provider-response through os.Root and the hash-locked Python evidence validator. Never performs HTTP. With a proven receipt, calls existing ReconcileAccepted; without it returns reconciliation required. A bare webhook ID or user assertion cannot manufacture an acceptance.

Unsupported inbound media and unknown statuses remain retained for explicit review; text with unknown contact produces the existing durable handoff and never reaches the model/tool/provider. An unresolved send remains unknown; do not provide a resend button. A new response requires a new genuine user message and a new reviewed proposal; changing text/key to force retry is prohibited.

## Empty activation inputs

Profile file/hash; fixed tenant/org/connection; purpose/policy/retention approval hash; private evidence directory; Python/source hashes from manifest; WhatsApp access-token/app-secret/verify-token sources; contact-HMAC source; OIDC verifier/service-token source; domain endpoint/token; LLM endpoint/model/token/budget; service/product/pricebook configuration; authenticated status reporter connection. These are configuration/account activation, not proof of deployment, provider reachability or live delivery.

Go integration tests use ELITE_WHATSAPP_CONNECTED_DATABASE_URL for a disposable loopback database named elite_whatsapp_connected_* and explicit ELITE_WHATSAPP_PYTHON. Test fixtures never contact Meta or the model provider.
