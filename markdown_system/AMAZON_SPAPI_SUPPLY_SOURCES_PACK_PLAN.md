# Amazon SP-API Supply Sources Pack Plan

Materializa el lock oficial y el adapter Amazon Supply Sources v2020-07-01 para inventariar y gobernar tiendas/depósitos. No contiene credenciales, no habilita llamadas, no persiste dirección/contacto y no convierte status Amazon en activación interna.

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
      "path": "implementation_packs/PYTHON_AMAZON_SPAPI_SUPPLY_SOURCES_ADAPTER.md",
      "packId": "PYTHON-AMAZON-SPAPI-SUPPLY-SOURCES-ADAPTER",
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

Antes de ejecutar: developer/app, rol/scope, marketplace/programa, dynamic sandbox, cuota/costo, tratamiento de ubicación/contacto y owner deben estar `PROVEN`. Cada mutación se liga a inventario y snapshot exactos, se intenta una vez y se reconcilia; `UNKNOWN_EFFECT` bloquea retry. Archive exige `Inactive`. El modelo interno de franquicias conserva autoridad sobre activación, inventario, promesas, transporte y aceptación.
