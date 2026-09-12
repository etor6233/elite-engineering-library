# Durable Document Pipeline Pack Plan

Este perfil materializa el plano de control documental vendor-neutral: adquisición oficial exacta, decisión de routing por clase, seguridad local de archivos, evaluación estricta y orquestación durable Microsoft. Produce 68 archivos desde cinco packs. No selecciona proveedor, no descarga fuentes, no contiene credenciales y el template habilita cero efectos. Sirve para completar la configuración y demostrar el flujo antes de añadir exactamente una lane Azure, Google, AWS o local aprobada.

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
      "path": "implementation_packs/DOCUMENT_PIPELINE_ROUTING_GATE.md",
      "packId": "DOCUMENT-PIPELINE-ROUTING-GATE",
      "version": "0.3.0",
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

El agente completa primero las 21 clases documentales base y toda clase adicional del proyecto, con schemas/campos críticos, owners, MIME/tamaños, corpus y ground truth; luego el receipt de routing selecciona una sola lane por clase. La orquestación conserva sólo referencias opacas y hashes, retiene el original mediante el adapter aprobado, extrae, evalúa, espera revisión humana cuando corresponde, persiste con clave idempotente y retiene evidencia final. El backend in-memory oficial Microsoft es únicamente de prueba; producción exige DTS/host aprobado, identidad, TLS, task hub, storage, load, recovery, versionado y pruebas live. Ningún fallo ambiguo autoriza retry ciego ni almacenamiento comercial.
