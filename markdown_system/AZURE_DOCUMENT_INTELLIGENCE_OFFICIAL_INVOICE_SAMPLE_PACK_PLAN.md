# Azure Document Intelligence Official Invoice Sample Pack Plan

Este perfil mínimo materializa catorce archivos desde dos packs: los doce artefactos oficiales documentales fijados por hash y el sample de invoice desde bytes publicado por Microsoft junto con su fixture oficial adquirible. No instala dependencias, no descarga sin aprobación, no contiene credenciales y no autoriza persistencia ni exactitud de negocio. Se usa sólo cuando el blueprint selecciona Azure Document Intelligence 1.0.2 y necesita demostrar primero la integración oficial `prebuilt-invoice`.

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
      "path": "implementation_packs/MICROSOFT_AZURE_DOCUMENT_INTELLIGENCE_OFFICIAL_INVOICE_SAMPLE.md",
      "packId": "MICROSOFT-AZURE-DOCUMENT-INTELLIGENCE-OFFICIAL-INVOICE-SAMPLE",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": ["*"],
      "variables": {}
    }
  ]
}
```

Después de componer, el agente ejecuta primero los tests offline. Para adquirir el JPEG oficial copia la plantilla, enlaza el SHA-256 del lock, registra identidad/fecha/aceptación MIT y usa cache aprobada o red explícita. La ejecución live requiere endpoint, key y costo/cuota autorizados por el usuario. Los resultados del sample entran luego en los gates documentales de schema, evidencia, corpus y revisión; nunca se escriben automáticamente en el sistema empresarial por este perfil.
