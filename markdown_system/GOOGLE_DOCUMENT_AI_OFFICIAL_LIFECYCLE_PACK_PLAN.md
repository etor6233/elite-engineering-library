# Google Document AI Official Lifecycle Pack Plan

Este perfil materializa 40 archivos desde tres packs: artifacts oficiales `google-cloud-documentai==3.15.0` + `google-cloud-storage==3.13.1`, samples Google de procesamiento/schema variable y el lifecycle oficial para processor, dataset/schema, importación train/test/unassigned, training, evaluation, deploy, default, migración, batch GCS, response handling OCR/form/table/entity/split/layout/custom y undeploy. No instala, no contiene credenciales y no ejecuta efectos cloud.

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
    },
    {
      "path": "implementation_packs/GOOGLE_DOCUMENT_AI_OFFICIAL_LIFECYCLE.md",
      "packId": "GOOGLE-DOCUMENT-AI-OFFICIAL-LIFECYCLE",
      "version": "0.2.0",
      "acknowledgeConditions": true,
      "files": ["*"],
      "variables": {}
    }
  ]
}
```

Después de componer, el agente ejecuta los verificadores offline y los cinco tests Google offline-safe. Antes de live completa GCP/ADC/IAM/storage/CMEK/cuota/costo, configura una copia del script documental oficial y exige corpus, ground truth y umbrales target. El test Google set-default completo y todos los efectos cloud permanecen bloqueados hasta autorización explícita.
