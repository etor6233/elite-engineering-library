# Amazon SP-API Easy Ship Handover Pack Plan

Materializa el lock oficial y el adapter Easy Ship v2022-03-23 para slots, programación y reconciliación. No contiene credenciales, no habilita llamadas y no convierte estados del proveedor en aceptación interna automática.

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
      "path": "implementation_packs/PYTHON_AMAZON_SPAPI_EASYSHIP_HANDOVER_ADAPTER.md",
      "packId": "PYTHON-AMAZON-SPAPI-EASYSHIP-HANDOVER-ADAPTER",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    }
  ]
}
```

Antes de ejecutar: registro, aplicación, rol/scope Easy Ship, marketplace y operación soportados, sandbox, cuota/costo, datos y owners deben estar `PROVEN`. Schedule exige slot inventariado y aprobación exacta; efecto ambiguo se reconcilia, no se reintenta. `PickedUp` sólo es reporte Amazon. Cancelación, firma/foto y carrier universal quedan fuera.
