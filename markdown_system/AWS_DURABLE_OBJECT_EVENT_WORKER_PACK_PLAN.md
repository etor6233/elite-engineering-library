# AWS Durable Object Event Worker Pack Plan

Este perfil compone tres autoridades separadas: adquisición oficial, idempotencia/partial batch SQS de AWS Powertools y checkpoint/retry/replay de AWS Lambda Durable Executions. Entrega código tangible, pero no inventa ni atribuye a AWS un worker integrado que AWS no publicó.

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
      "path": "implementation_packs/AWS_LAMBDA_DURABLE_EXECUTION_COMPONENT.md",
      "packId": "AWS-LAMBDA-DURABLE-EXECUTION-COMPONENT",
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

Antes de desplegar, el proyecto debe aportar dispatcher/execution name, identidad tenant-bound, S3 version/eTag/sequencer, efecto idempotente, owner/fencing si existe reclaim, IaC/IAM, DLQ/redrive, reconciliación, carga y evidencia en la cuenta destino. El perfil acelera esos componentes; no declara exactly-once ni autorización productiva.
