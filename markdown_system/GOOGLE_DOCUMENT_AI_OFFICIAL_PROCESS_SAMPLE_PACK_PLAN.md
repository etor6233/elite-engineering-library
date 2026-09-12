# Google Document AI Official Process Sample Pack Plan

Este perfil mínimo materializa diecinueve archivos desde dos packs: los artefactos SDK documentales fijados, el sample standalone Google `process_document`, el sample oficial Custom Document Extractor con `schema_override` por solicitud, sus tests live exactos, licencias, adquisición de invoice/packing-list y validación del output oficial. No instala, no descarga sin approval, no contiene ADC/proyecto/processor y mantiene `automatic_storage_authorized=false`.

```json
{
  "compositionVersion": "1.0",
  "secretVariables": {},
  "packs": [
    {
      "path": "implementation_packs/OFFICIAL_DOCUMENT_SDK_ARTIFACT_CORE.md",
      "packId": "OFFICIAL-DOCUMENT-SDK-ARTIFACT-CORE",
      "version": "0.3.1",
      "acknowledgeConditions": true,
      "files": ["*"],
      "variables": {}
    },
    {
      "path": "implementation_packs/GOOGLE_CLOUD_DOCUMENT_AI_OFFICIAL_PROCESS_SAMPLE.md",
      "packId": "GOOGLE-CLOUD-DOCUMENT-AI-OFFICIAL-PROCESS-SAMPLE",
      "version": "0.2.0",
      "acknowledgeConditions": true,
      "files": ["*"],
      "variables": {}
    }
  ]
}
```

Después de componer, el agente ejecuta las regresiones offline —incluido el contrato CDE exacto contra SDK 3.15.0—, completa approval, adquiere los cuatro fixtures exactos y valida el JSON packing-list. Los tests Google live sólo se ejecutan con GCP/ADC/processors/costo autorizados. La ruta oficial CDE permite configurar clases/campos variables; las confidencias oficiales se conservan como evidencia y no se convierten por sí solas en datos almacenables.
