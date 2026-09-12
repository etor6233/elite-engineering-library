# Google Merchant Product Sync Pack Plan

Materializa el source lock y el adapter oficial `productInputs.insert` + `products.get`. No contiene credenciales y la plantilla bloquea toda escritura externa.

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
      "path": "implementation_packs/PYTHON_GOOGLE_MERCHANT_PRODUCT_SYNC_ADAPTER.md",
      "packId": "PYTHON-GOOGLE-MERCHANT-PRODUCT-SYNC-ADAPTER",
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

Antes de ejecutar: proyecto/billing/API, ADC, Merchant account, data source API primary, test account, offers, product spec, escritura externa, quota/cost, refresh, retention y reconciliation deben estar `PROVEN`. `automatic_local_business_write=false` y una inserción nunca se considera aprobación.
