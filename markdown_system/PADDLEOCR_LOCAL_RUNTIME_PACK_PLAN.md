# PaddleOCR Local Runtime Pack Plan

Este perfil materializa una lane OCR local completa y condicionada: adquisición oficial, seguridad local ClamAV/Magika/YARA-X, runtime PaddleOCR/PaddleX/PaddlePaddle con grafo y modelos fijados, y evaluación estricta de campos. Produce 58 archivos desde cuatro packs. No instala, descarga modelos, habilita clases ni autoriza almacenamiento automáticamente.

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
      "path": "implementation_packs/PADDLEPADDLE_PADDLEOCR_LOCAL_RUNTIME.md",
      "packId": "PADDLEPADDLE-PADDLEOCR-LOCAL-RUNTIME",
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
    }
  ]
}
```

Antes de instalar: aprobar las dos fuentes Paddle por hash, usar CPython 3.12 Windows x86-64, descargar los ocho archivos de modelos sólo desde sus revisiones fijadas y completar formatos, clases, límites, corpus y política. Antes de cada OCR: exigir un receipt ADMITTED de los mismos bytes. Antes de persistir: el gate estricto debe demostrar schema/ground truth/métricas/rechazo/revisión; la lane siempre entrega `automatic_business_storage=false`.
