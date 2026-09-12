# Amazon SP-API Fulfillment Delivery Evidence Pack Plan

Materializa el lock oficial y el adapter Fulfillment Outbound v2020-07-01 para recuperar evidencia de entrega de paquetes MCF. No contiene credenciales, no habilita llamadas, no persiste URLs firmadas ni identidad y no convierte evidencia provider en aceptación empresarial.

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
      "path": "implementation_packs/PYTHON_AMAZON_SPAPI_FULFILLMENT_DELIVERY_EVIDENCE_ADAPTER.md",
      "packId": "PYTHON-AMAZON-SPAPI-FULFILLMENT-DELIVERY-EVIDENCE-ADAPTER",
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

Antes de ejecutar: registro, aplicación, rol/scope Amazon Fulfillment, soporte MCF/marketplace, sandbox, cuota/costo, privacidad, retención, owner y listas observadas de status/tipos/hosts deben estar `PROVEN`. El request debe estar ligado al receipt empresarial y a paquetes exactos. Foto, firma o drop-off sólo son evidencia reportada por Amazon; aceptación, inspección y cierre del negocio son decisiones separadas.
