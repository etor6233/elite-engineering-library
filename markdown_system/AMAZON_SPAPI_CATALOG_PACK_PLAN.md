# Amazon SP-API Catalog Pack Plan

Materializa el lock oficial y el adapter read-only Amazon Catalog Items. No contiene credenciales ni habilita llamadas por sí solo.

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
      "path": "implementation_packs/PYTHON_AMAZON_SPAPI_CATALOG_ADAPTER.md",
      "packId": "PYTHON-AMAZON-SPAPI-CATALOG-ADAPTER",
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

Antes de ejecutar: developer/application registration, roles/data access, región/marketplace, ASINs/datasets, sandbox contract, OAuth secret references, quota/cost y reconciliation deben estar `PROVEN`. La plantilla sigue `BLOCKED_ACCESS_AND_SCOPE_APPROVAL_REQUIRED` y `automatic_business_write=false`.
