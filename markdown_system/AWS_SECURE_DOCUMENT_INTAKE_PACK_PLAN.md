# AWS Secure Document Intake Pack Plan

Este perfil materializa una lane condicionada de entrada portal/API: lock/adquiridor oficial, sesión S3 de carga directa a cuarentena y gate local Magika/ClamAV/YARA-X. Produce 52 archivos desde tres packs. No crea cuenta/bucket/IAM/KMS/CORS, no descarga binarios automáticamente, no abre un endpoint anónimo y no autoriza almacenamiento empresarial.

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
      "path": "implementation_packs/GO_AWS_SECURE_QUARANTINE_INTAKE.md",
      "packId": "GO-AWS-SECURE-QUARANTINE-INTAKE",
      "version": "0.1.0",
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

Orden obligatorio: readiness del proyecto; cuenta/IAM/costo y bucket/KMS/versioning/Block Public Access/policy/CORS probados; endpoint autenticado que guarda la sesión en store transaccional; presigned PUT con todos los headers; finalize por version ID/ETag; descarga aislada de esa versión; gate local de seguridad sobre los mismos bytes; recién entonces routing/clasificación/evaluación. Email, SFTP, scanner/mobile, eventos/reconciliación y ERP/marketplace siguen fuera de este perfil y no se suponen resueltos.
