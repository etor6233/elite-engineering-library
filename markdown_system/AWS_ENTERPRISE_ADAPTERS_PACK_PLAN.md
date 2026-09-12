# AWS Enterprise Adapters Pack Plan

Este perfil materializa el lock/adquiridor oficial y los adapters Amazon S3/SES v2. No contiene ni solicita valores de credenciales en el plan, no crea infraestructura y no habilita storage/email hasta completar el profile fail-closed.

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
      "path": "implementation_packs/GO_AWS_ENTERPRISE_STORAGE_EMAIL_ADAPTERS.md",
      "packId": "GO-AWS-ENTERPRISE-STORAGE-EMAIL-ADAPTERS",
      "version": "0.2.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    }
  ]
}
```

El agente pide región/IAM/costo, expected bucket owner, bucket/versioning/Object Lock/lifecycle/KMS, owner/modo/fecha de retención y SES identity/access/configuration set/events/suppression/outbox antes de una llamada. El profile conserva decisiones bloqueadas y ambos `cost_approved=false` hasta evidencia del usuario/proveedor. El adapter consulta Object Lock, escribe con retención y verifica la versión exacta; nunca habilita Object Lock ni legal hold porque esas mutaciones requieren autoridad separada.
