# Amazon SP-API External Fulfillment Inventory Pack Plan

Materializa el lock oficial y el adapter Amazon External Fulfillment Inventory v2024-09-11 para FETCH/UPDATE batch de inventario absoluto location-level. No contiene credenciales ni asume que sus location IDs sean Supply Sources internos.

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
      "path": "implementation_packs/PYTHON_AMAZON_SPAPI_EXTERNAL_FULFILLMENT_INVENTORY_ADAPTER.md",
      "packId": "PYTHON-AMAZON-SPAPI-EXTERNAL-FULFILLMENT-INVENTORY-ADAPTER",
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

Antes de ejecutar: cuenta/app, programa External Fulfillment, rol/scope, canal, mapping location+SKU, sequence authority, sandbox, cuota/costo, privacidad y owner deben estar `PROVEN`. FETCH precede UPDATE; approval liga bytes exactos; batch parcial o excepción produce `UNKNOWN_EFFECT` sin retry. El ERP/POS/WMS conserva autoridad y no hace commit automático por un acknowledgement Amazon.
