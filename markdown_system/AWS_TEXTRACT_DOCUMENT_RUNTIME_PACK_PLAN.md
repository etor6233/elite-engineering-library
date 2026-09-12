# AWS Textract Document Runtime Pack Plan

Este perfil materializa la lane AWS completa anterior al negocio: adquisición oficial exacta, seguridad local de archivos, storage S3 con Object Lock verificable, runtime Amazon Textract sobre AWS SDK for Go v2, componente oficial Amazon Textractor Python, evaluación estricta y orquestación durable de punta a punta. Produce 92 archivos desde siete packs. No descarga automáticamente, no contiene credenciales, no habilita Object Lock y no autoriza almacenamiento de campos comerciales. Se usa sólo cuando el receipt de routing selecciona AWS para al menos una clase.

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
      "path": "implementation_packs/SECURE_LOCAL_FILE_INGESTION_GATE.md",
      "packId": "SECURE-LOCAL-FILE-INGESTION-GATE",
      "version": "0.1.0",
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
    },
    {
      "path": "implementation_packs/GO_AWS_TEXTRACT_DOCUMENT_RUNTIME.md",
      "packId": "GO-AWS-TEXTRACT-DOCUMENT-RUNTIME",
      "version": "0.2.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/PYTHON_AWS_TEXTRACTOR_OFFICIAL_COMPONENT.md",
      "packId": "PYTHON-AWS-TEXTRACTOR-OFFICIAL-COMPONENT",
      "version": "0.1.0",
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
      "path": "implementation_packs/MICROSOFT_DURABLE_DOCUMENT_ORCHESTRATION.md",
      "packId": "MICROSOFT-DURABLE-DOCUMENT-ORCHESTRATION",
      "version": "0.1.2",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    }
  ]
}
```

Antes de adquirir se aprueba el source AWS exacto. El usuario debe aportar cuenta/owner, bucket ya configurado, IAM, KMS, modo/fecha de retención, lifecycle y costo; el adapter nunca habilita Object Lock. El gate local debe emitir un receipt `ADMITTED` que coincida con el archivo; cualquier original que se conserve usa checksum, owner esperado, `If-None-Match`, VersionId y retención verificada. El runtime Go deriva región, operación, features, Queries y AdapterId/Version sólo de una clase `REQUIRED`; Textractor aporta el parser/CLI oficial AWS 1.10.0 con wheel/grafo hash-locked, pruebas Queries 2/1 y suite determinista Linux 69/16, sin sustituir el runner seguro ni habilitar `CALL_TEXTRACT` en offline. Luego la evaluación conserva `automatic_storage_authorized=false` hasta que corpus/ground truth/política prueben cero falsos bajo el contrato del proyecto. La orquestación enlaza esos receipts, exige revisión hash-bound cuando corresponde, persiste idempotentemente y conserva evidencia final. Toda clase `CONFIGURATION_REQUIRED` o `BLOCKED_*` permanece sin automatizar. Storage de evidencia no equivale a persistencia de campos comerciales; IAM, presupuesto, cuotas, residencia, reconciliación y sandbox real siguen siendo entradas obligatorias.
