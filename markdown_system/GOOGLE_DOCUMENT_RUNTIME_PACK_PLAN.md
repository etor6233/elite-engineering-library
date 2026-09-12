# Google Document Runtime Pack Plan

Este perfil materializa la lane Google completa anterior al negocio: adquisición oficial exacta, seguridad local de archivos, lock de artefactos SDK, código Google Custom Document Extractor con schema variable por solicitud, lifecycle oficial de processor/dataset/import/train/evaluate/deploy/default/undeploy, batch GCS y response handling OCR/form/table/entity/split/layout/custom, runtime Google Cloud Document AI + Storage, evaluación estricta, adapter GCS con creación única/CRC32C/KMS/retención verificada y orquestación durable de punta a punta. Produce 119 archivos desde nueve packs. No descarga ni instala automáticamente, no contiene credenciales y mantiene todo efecto bloqueado hasta approvals y gates del proyecto. Se usa sólo cuando `PROJECT_DOCUMENT_INTELLIGENCE_DECISION.md` selecciona Google para al menos una clase.

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
      "path": "implementation_packs/OFFICIAL_DOCUMENT_SDK_ARTIFACT_CORE.md",
      "packId": "OFFICIAL-DOCUMENT-SDK-ARTIFACT-CORE",
      "version": "0.3.1",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GOOGLE_DOCUMENT_AI_RUNTIME.md",
      "packId": "GOOGLE-DOCUMENT-AI-RUNTIME",
      "version": "0.1.3",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GOOGLE_CLOUD_DOCUMENT_AI_OFFICIAL_PROCESS_SAMPLE.md",
      "packId": "GOOGLE-CLOUD-DOCUMENT-AI-OFFICIAL-PROCESS-SAMPLE",
      "version": "0.2.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/GOOGLE_DOCUMENT_AI_OFFICIAL_LIFECYCLE.md",
      "packId": "GOOGLE-DOCUMENT-AI-OFFICIAL-LIFECYCLE",
      "version": "0.2.0",
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
      "path": "implementation_packs/GO_GOOGLE_CLOUD_STORAGE_ADAPTER.md",
      "packId": "GO-GOOGLE-CLOUD-STORAGE-ADAPTER",
      "version": "0.1.0",
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

Antes de adquirir se completan approvals enlazados a los hashes exactos. Antes de analizar, el gate local debe emitir un receipt `ADMITTED` que coincida con el archivo; el runtime sólo acepta una clase `REQUIRED` y deriva el processor version desde el perfil. Luego la evaluación conserva `automatic_storage_authorized=false` hasta que corpus/ground truth/política del proyecto demuestren la promoción. Sólo entonces el adapter GCS puede recibir approvals de acceso/costo/retención y verificar owner, bucket, generación, checksum, KMS y retención; la orquestación enlaza esos receipts, exige revisión hash-bound cuando corresponde, persiste idempotentemente y conserva evidencia final. No crea ni reconfigura infraestructura. Toda clase `CONFIGURATION_REQUIRED` o `BLOCKED_*` permanece sin automatizar. Document AI y GCS tienen cargos/cuotas: el gate de inicio debe obtener aprobación de costo antes de una llamada real.
