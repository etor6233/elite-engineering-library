# Mercado Libre Marketplace Pack Plan

Materializa la adaptación Go declarada contra los contratos HTTP oficiales vigentes de Mercado Libre. No usa el SDK oficial archivado, no contiene tokens y no implementa publicación ni respuesta: lectura, questions v4, reconciliación de preguntas sin responder, validador previo y fetch jobs de notificaciones.

```json
{
  "compositionVersion": "1.0",
  "secretVariables": {},
  "packs": [
    {
      "path": "implementation_packs/GO_MERCADOLIBRE_MARKETPLACE_ADAPTER.md",
      "packId": "GO-MERCADOLIBRE-MARKETPLACE-ADAPTER",
      "version": "0.2.0",
      "acknowledgeConditions": true,
      "files": ["*"],
      "variables": {}
    }
  ]
}
```

Antes de usar: aplicación/owner legal, seller grant, token store, scopes, site/seller/application IDs, términos y políticas, inexistencia de preproducción, origen/topics de callbacks —incluido `questions`—, cuotas/costo y reconciliación deben estar `PROVEN`. `VALIDATE_ITEM` exige aprobación separada; no publica. Las respuestas v4 se normalizan con `GO-OMNICHANNEL-LEAD-INGRESS 0.2.x`; contacto y respuesta permanecen separados.
