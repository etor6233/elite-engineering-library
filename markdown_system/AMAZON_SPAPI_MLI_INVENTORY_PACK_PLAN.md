# Amazon SP-API Multi-Location Inventory Pack Plan

Materializa el lock oficial, el inventario Supply Sources y el adapter Listings Items para snapshot/reemplazo de inventario por ubicación Amazon MLI. No contiene credenciales, no habilita llamadas ni convierte el receipt provider en commit automático del inventario interno.

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
    },
    {
      "path": "implementation_packs/PYTHON_AMAZON_SPAPI_MLI_INVENTORY_ADAPTER.md",
      "packId": "PYTHON-AMAZON-SPAPI-MLI-INVENTORY-ADAPTER",
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

Antes de ejecutar: developer/app, rol Product Listing, inscripción MLI, marketplace, sources activos, dynamic sandbox, cuota/costo, privacidad y owner deben estar `PROVEN`. El reemplazo usa `PRODUCT` y `/attributes/fulfillment_availability`, exige approval y snapshots exactos, se intenta una vez y se reconcilia. `UNKNOWN_EFFECT` bloquea retry; el sistema interno mantiene autoridad hasta consumir un receipt reconciliado.
