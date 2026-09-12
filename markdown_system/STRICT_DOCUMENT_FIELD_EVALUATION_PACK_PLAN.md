# Strict Document Field Evaluation Pack Plan

Este perfil materializa únicamente la adquisición oficial exacta y el gate de evaluación estructurada. No descarga sin aprobación, no inventa corpus/schema/ground truth, no ejecuta un extractor y no autoriza almacenamiento automático.

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

Antes de adquirir: completar el perfil `elite_sources/source-profiles/strict-field-evaluation.json`, su aprobación hash-linked y todos los blockers. Antes de evaluar: aprobar clases, schema cerrado, corpus representativo, ground truth independiente, privacidad/retención y owner de rechazo/revisión. Un PASS sólo prueba esos casos y conserva `automatic_storage_authorized=false`.
