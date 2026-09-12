# Document Pipeline Routing Pack Plan

Este perfil materializa el source lock/adquiridor oficial y el gate determinista de decisión documental. Produce 44 archivos desde dos packs, no descarga fuentes, no llama providers, no contiene credenciales y no autoriza storage. Su receipt selecciona perfiles posteriores sólo desde respuestas y evidencia explícitas del proyecto.

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
    }
  ]
}
```

El template habilita cero rutas: el agente pregunta y registra las 21 clases base y cualquier clase adicional del proyecto, junto con owners, schemas/fields/MIME/límites, corpus, acceso, seguridad, evaluación y candidatos. Sólo después de un receipt V3 válido materializa los planes Azure, Google, AWS o MarkItDown exactos indicados; el receipt separa IDs base/adicionales y conserva cada ruta clase→provider→role→plan→provider profile. El gate no transforma un sample oficial condicionado en código productivo y conserva `automatic_storage_authorized=false`.
