# Microsoft MarkItDown Local Runtime Pack Plan

Este perfil materializa la lane local completa anterior al negocio: source lock/adquiridor oficial, seguridad local ClamAV/Magika/YARA-X, runtime Microsoft MarkItDown fijado para CPython 3.12 Windows x86-64 y evaluación estricta de campos. Produce 57 archivos desde cuatro packs. No descarga ni instala automáticamente, no habilita formatos/URLs/plugins/archives/LLM por defecto y no autoriza almacenar datos de negocio.

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
      "path": "implementation_packs/MICROSOFT_MARKITDOWN_LOCAL_RUNTIME.md",
      "packId": "MICROSOFT-MARKITDOWN-LOCAL-RUNTIME",
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
    }
  ]
}
```

Antes de instalar: aprobar el source profile Microsoft/local, configurar explícitamente cada formato/MIME/límite, ejecutar el gate de seguridad y confirmar CPython 3.12 Windows x86-64 en un venv vacío. El template MarkItDown habilita cero formatos y el runtime no acepta overrides de política por CLI. Antes de usar output: tratar Markdown como input no confiable, seleccionar analyzer/processor semántico y demostrar fields contra ground truth mediante el gate estricto; la conversión sola nunca habilita `automatic_business_storage`.
