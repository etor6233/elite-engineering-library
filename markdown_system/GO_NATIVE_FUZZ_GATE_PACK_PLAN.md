# Go Native Fuzz Gate Pack Plan

Perfil opt-in para proyectos Go que ya poseen invariantes reales expresadas como targets `FuzzXxx`. No inventa targets ni sustituye DAST, autorización o carga del target.

```json
{
  "compositionVersion": "1.0",
  "secretVariables": {},
  "packs": [
    {
      "path": "implementation_packs/GO_NATIVE_FUZZ_GATE.md",
      "packId": "GO-NATIVE-FUZZ-GATE",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": ["*"],
      "variables": {}
    }
  ]
}
```

Secuencia: materializar → verificar runner → identificar invariantes/semillas reales → habilitar perfil → baseline → corpus → fuzz acotado → conservar receipt/corpus de fallo → corregir y promover la semilla a regresión → repetir gates afectados.
