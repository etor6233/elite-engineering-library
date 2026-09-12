# AWS IDP PostgreSQL Persistence Pack Plan

Perfil focal de 5 packs/56 archivos para reconstruir decisión v2, foundation PostgreSQL, frontera transaccional intake+outbox, runtime Debezium PostgreSQL fijado por digest OCI y consumer inbox/proyección basado en el sample oficial Red Hat con mapping Google-origin CEL hash-bound integrado en su transacción. No incluye infraestructura upstream de extracción ni Kafka/broker; compóngalo con el perfil documental y la plataforma de mensajería elegidos. Exige aprobación explícita y evidencia target antes de red o datos reales.

```json
{
  "compositionVersion": "1.0",
  "secretVariables": {},
  "packs": [
    {"path":"implementation_packs/AWS_IDP_EVALUATION_DECISION_WORKER.md","packId":"AWS-IDP-EVALUATION-DECISION-WORKER","version":"0.2.0","acknowledgeConditions":true,"files":["*"],"variables":{}},
    {"path":"implementation_packs/POSTGRES_TRANSACTIONAL_FOUNDATION.md","packId":"PG-TX-FOUNDATION","version":"0.1.0","acknowledgeConditions":true,"files":["*"],"variables":{}},
    {"path":"implementation_packs/AWS_IDP_POSTGRES_PERSISTENCE_BOUNDARY.md","packId":"AWS-IDP-POSTGRES-PERSISTENCE-BOUNDARY","version":"0.2.0","acknowledgeConditions":true,"files":["*"],"variables":{}},
    {"path":"implementation_packs/DEBEZIUM_POSTGRES_OUTBOX_RUNTIME.md","packId":"DEBEZIUM-POSTGRES-OUTBOX-RUNTIME","version":"0.1.0","acknowledgeConditions":true,"files":["*"],"variables":{}},
    {"path":"implementation_packs/DEBEZIUM_POSTGRES_INBOX_CONSUMER.md","packId":"DEBEZIUM-POSTGRES-INBOX-CONSUMER","version":"0.2.0","acknowledgeConditions":true,"files":["*"],"variables":{}}
  ]
}
```

Orden: materializar; ejecutar 13+12+7+18 pruebas, cfn-lint y SCA consumer integrado 170/0; fijar schemas/mapping/clase/corpus/aprobación; aplicar foundation+migrations; crear roles/secrets no-owner, CDC read-only y consumer insert/select; probar RLS; aprobar/subir policy Object-Locked; habilitar S3 EventBridge; configurar publicación/slot; registrar conector y consumer exactos en DEV; probar duplicates/cross-class/mapping-change/races/DLQ/outbox outage/restart/redelivery/order/restore/rollback. Los mappings ERP/CRM son configuración versionada del consumer, no autoridad automática: cada negocio debe probar su clase/schema/corpus.
