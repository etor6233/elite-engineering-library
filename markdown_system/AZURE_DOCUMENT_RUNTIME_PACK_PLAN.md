# Azure Document Runtime Pack Plan

Este perfil materializa la lane Azure completa anterior al negocio: adquisición oficial exacta, seguridad local de archivos, lock de artefactos SDK, runtime GA Microsoft Azure Content Understanding, los samples Python oficiales Microsoft `prebuilt-invoice`, análisis de bytes locales con `prebuilt-documentSearch` y same-resource analyzer copy, evaluación estricta, almacenamiento Azure Blob WORM y orquestación durable de punta a punta. Produce 102 archivos desde diez packs. No descarga ni instala automáticamente, no contiene credenciales y mantiene todo efecto bloqueado hasta approvals y gates del proyecto. Se usa sólo cuando `PROJECT_DOCUMENT_INTELLIGENCE_DECISION.md` selecciona Azure para al menos una clase.

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
      "path": "implementation_packs/AZURE_CONTENT_UNDERSTANDING_DOCUMENT_RUNTIME.md",
      "packId": "AZURE-CONTENT-UNDERSTANDING-DOCUMENT-RUNTIME",
      "version": "0.1.2",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/MICROSOFT_AZURE_CONTENT_UNDERSTANDING_OFFICIAL_INVOICE_SAMPLE.md",
      "packId": "MICROSOFT-AZURE-CONTENT-UNDERSTANDING-OFFICIAL-INVOICE-SAMPLE",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/MICROSOFT_AZURE_CONTENT_UNDERSTANDING_OFFICIAL_BINARY_DOCUMENT_SAMPLE.md",
      "packId": "MICROSOFT-AZURE-CONTENT-UNDERSTANDING-OFFICIAL-BINARY-DOCUMENT-SAMPLE",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    },
    {
      "path": "implementation_packs/MICROSOFT_AZURE_CONTENT_UNDERSTANDING_OFFICIAL_ANALYZER_COPY_SAMPLE.md",
      "packId": "MICROSOFT-AZURE-CONTENT-UNDERSTANDING-OFFICIAL-ANALYZER-COPY-SAMPLE",
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
      "path": "implementation_packs/PYTHON_AZURE_BLOB_IMMUTABLE_EVIDENCE_ADAPTER.md",
      "packId": "PYTHON-AZURE-BLOB-IMMUTABLE-EVIDENCE-ADAPTER",
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

Antes de adquirir se completan approvals enlazados a los hashes exactos. Antes de analizar, el gate local debe emitir un receipt `ADMITTED` que coincida con el archivo; el runtime sólo acepta una clase `REQUIRED` y deriva el analyzer desde el perfil. Luego la evaluación conserva `automatic_storage_authorized=false` hasta que corpus/ground truth/política demuestren la promoción. Sólo entonces el adapter Azure recibe approvals de acceso/costo/retención y un receipt de autoridad fresco, exige contenedor privado con version-level WORM y scope no reemplazable, crea con `Locked` atómico y verifica propiedades más SHA-256 descargando el `version_id` exacto. La orquestación enlaza esos receipts, exige revisión hash-bound cuando corresponde, persiste idempotentemente y conserva evidencia final; no crea ni migra infraestructura. Toda clase `BLOCKED_*` permanece sin automatizar.
