# Amazon SP-API Shipping Tracking Pack Plan

Materializa el lock oficial y el adapter Amazon Shipping V2 para ciclo, collection form y claim aprobados. No contiene credenciales, no habilita llamadas y no completa handover/entregas.

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
      "path": "implementation_packs/PYTHON_AMAZON_SPAPI_SHIPPING_TRACKING_ADAPTER.md",
      "packId": "PYTHON-AMAZON-SPAPI-SHIPPING-TRACKING-ADAPTER",
      "version": "0.6.0",
      "acknowledgeConditions": true,
      "files": [
        "*"
      ],
      "variables": {}
    }
  ]
}
```

Antes de ejecutar: developer/application registration, Shipping role/scope, región/business/carrier, sandbox, secretos, costo, datos y reconciliación deben estar `PROVEN`. Ciclo/collection/claim requieren approvals separados y artefactos hash-bound; claims además insured value y proof hosts. Form/claim no autorizan handover ni entrega.
