# Capability Gap Resolution Pack Plan

Perfil portable para una capability `REQUIRED` que no tiene implementación compatible admitida o cuya autoridad quedó obsoleta/rechazada. No agrega producto: obliga a investigar, demostrar procedencia y cerrar con adopción válida o bloqueo explícito.

```json
{
  "compositionVersion": "1.0",
  "secretVariables": {},
  "packs": [
    {
      "path": "implementation_packs/CAPABILITY_GAP_RESOLUTION_GATE.md",
      "packId": "CAPABILITY-GAP-RESOLUTION-GATE",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": ["*"],
      "variables": {}
    }
  ]
}
```

Secuencia: detectar gap → enlazar requirement/owner/autoridades → investigar docs/repos/releases/advisories oficiales actuales → fijar candidatos/licencia/revisión/hash/claim → G0–G8 → validar expediente → actualizar mapa/lock/catálogo/pack o informar bloqueo y trigger de reapertura → ejecutar gates target. `RESEARCH_INCOMPLETE` nunca autoriza código.
