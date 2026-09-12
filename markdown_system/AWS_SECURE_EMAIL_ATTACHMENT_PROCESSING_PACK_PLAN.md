# AWS Secure Email Attachment Processing Pack Plan

Este perfil compone 14 packs/157 archivos: adquisición oficial, receptor SES/S3/SQS/Lambda/DynamoDB materializable, extracción MIME fail-closed, liberación GuardDuty por versión exacta, identificación con Google Magika, dispatch exact-version al AWS GenAI IDP Accelerator v0.6.5, handoff postprocessing/review inmutable, componente Powertools oficial, evaluación offline estricta, worker de decisión idempotente v2, foundation PostgreSQL, frontera Aurora Data API de intake+outbox, runtime Debezium PostgreSQL fijado por digests OCI y consumer inbox/proyección Red Hat adaptado con mapping CEL-Java integrado dentro de la misma transacción. Los workers y la infraestructura local son adaptaciones declaradas sobre contratos/samples oficiales fijados; no se atribuye a AWS, Google, PostgreSQL, Dapr, Red Hat ni Debezium el hardening local. Toda infraestructura nace sólo después de completar profiles, probar cuenta/identity/MX/KMS/aislamiento/configuración/corpus/HITL/database/Kafka Connect y presentar aprobaciones de efectos/costo; el PASS offline no autoriza AWS/Aurora/Kafka ni mappings ERP finales.

```json
{
  "compositionVersion": "1.0",
  "secretVariables": {},
  "packs": [
    {
      "path": "implementation_packs/OFFICIAL_UPSTREAM_ACQUISITION_CORE.md",
      "packId": "OFFICIAL-UPSTREAM-ACQUISITION-CORE",
      "version": "0.4.88",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/AWS_SES_IMMUTABLE_EMAIL_RECEIVER.md",
      "packId": "AWS-SES-IMMUTABLE-EMAIL-RECEIVER",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/SECURE_EMAIL_MIME_QUARANTINE_CORE.md",
      "packId": "SECURE-EMAIL-MIME-QUARANTINE-CORE",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/AWS_GUARDDUTY_IMMUTABLE_RELEASE_GATE.md",
      "packId": "AWS-GUARDDUTY-IMMUTABLE-RELEASE-GATE",
      "version": "0.1.1",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/AWS_GUARDDUTY_MAGIKA_IDP_DISPATCH_GATE.md",
      "packId": "AWS-GUARDDUTY-MAGIKA-IDP-DISPATCH-GATE",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/AWS_POWERTOOLS_IDEMPOTENT_SQS_BATCH_COMPONENT.md",
      "packId": "AWS-POWERTOOLS-IDEMPOTENT-SQS-BATCH-COMPONENT",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/AWS_IDP_IMMUTABLE_EVALUATION_HANDOFF.md",
      "packId": "AWS-IDP-IMMUTABLE-EVALUATION-HANDOFF",
      "version": "0.1.2",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/STRICT_DOCUMENT_FIELD_EVALUATION_GATE.md",
      "packId": "STRICT-DOCUMENT-FIELD-EVALUATION-GATE",
      "version": "0.1.1",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/AWS_IDP_EVALUATION_DECISION_WORKER.md",
      "packId": "AWS-IDP-EVALUATION-DECISION-WORKER",
      "version": "0.2.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/POSTGRES_TRANSACTIONAL_FOUNDATION.md",
      "packId": "PG-TX-FOUNDATION",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/AWS_IDP_POSTGRES_PERSISTENCE_BOUNDARY.md",
      "packId": "AWS-IDP-POSTGRES-PERSISTENCE-BOUNDARY",
      "version": "0.2.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/DEBEZIUM_POSTGRES_OUTBOX_RUNTIME.md",
      "packId": "DEBEZIUM-POSTGRES-OUTBOX-RUNTIME",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/DEBEZIUM_POSTGRES_INBOX_CONSUMER.md",
      "packId": "DEBEZIUM-POSTGRES-INBOX-CONSUMER",
      "version": "0.2.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/SECURE_LOCAL_FILE_INGESTION_GATE.md",
      "packId": "SECURE-LOCAL-FILE-INGESTION-GATE",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    }
  ]
}
```

Orden: readiness/proveedor/profiles; recepción SES; raw Object Lock; MIME/cuarentena; GuardDuty exact-version; Magika y dispatch versionado al IDP; defensa local; clasificación/extracción/HITL; hook `onError: fail`; evaluación/policy; worker decision v2; foundation; boundary intake+outbox append-only; conector Debezium; consumer que lee el payload revisado, exige clase documental exacta, aplica mapping CEL ID/version/hash/schema/corpus aprobados, rechaza cross-class y transacciona inbox+proyección antes del ack. Gate consumer: 18/18 PostgreSQL 18.6, concurrencia 4×8 y SCA 170/0. DNS/SES/AWS/IDP/Aurora/Kafka/costos/canaries/SNS/DLQ-redrive/corpus/HITL/RLS/race/load/restart/redelivery/order/restore/rollback y aprobación de mappings ERP/CRM finales siguen condicionados.
