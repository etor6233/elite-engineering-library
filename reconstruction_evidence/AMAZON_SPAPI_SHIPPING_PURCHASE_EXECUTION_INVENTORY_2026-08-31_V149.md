# Amazon SP-API Shipping purchase — execution inventory V149

Fecha: 2026-08-31. Resultado: `REBUILD_VERIFIED / CONDITIONED`.

## Autoridad exacta

- Amazon `selling-partner-api-sdk`, release Python 1.11.1, commit `8e792ae345a8d334ccdbdd03181f05f040e6a4fc`, Apache-2.0.
- Amazon `selling-partner-api-models`, commit `8e429486005c4ebdce5099e48cc48515a65359bb`, Apache-2.0.
- Wheel oficial PyPI `amzn_sp_api-1.11.1-py3-none-any.whl`, SHA-256 `a88e4059c9fbfd454bfef67b82288a926f73cbdc608a8831f8383565736c70b0`.
- Contrato generado inspeccionado: `ShippingApi.purchase_shipment(body, x_amzn_idempotency_key, x_amzn_shipping_business_id)`; Amazon documenta compra dentro de diez minutos de crear la tarifa y uso de la clave para reconocer reintentos.
- Referencias: https://github.com/amzn/selling-partner-api-sdk/tree/8e792ae345a8d334ccdbdd03181f05f040e6a4fc, https://github.com/amzn/selling-partner-api-models/tree/8e429486005c4ebdce5099e48cc48515a65359bb y https://developer-docs.amazon.com/sp-api/reference/purchaseshipment.

Los tres archivos nuevos son `AUTHORED`: integración local gobernada por el SDK/modelo oficial. No se presentan como código Amazon.

## Artefacto reconstruido

- Pack `PYTHON-AMAZON-SPAPI-SHIPPING-TRACKING-ADAPTER` 0.3.0: 67.277 bytes, SHA-256 `2dbe917aeba4696aaa3c4022056c924298beaf19e989edda426e3d714032a5a2`.
- Materialización aislada: 13/13 archivos; nuevos `purchase-selection.template.json` 464 bytes/SHA-256 `e62b3b1f7b83a6663069cfd5e7f34dbf57a2c30da5de79ec696b0ac400872dac`, `run_shipping_purchase.py` 15.690/`e348f13a4090d78fac295522f838867f8abcf68838ef758d319b1ecb1a98b177` y `test_shipping_purchase.py` 7.560/`a4ed2d36792b8b7c7a28f2d200d2ce147e10a2810039db8a03cd559723722ac1`.
- Perfil Amazon Shipping: 2 packs/38 payload files más `MATERIALIZATION_RECORD.md`; record 11.642 bytes/SHA-256 `c059801caa82688553d24f451f9db77c5b6e67e98e530c106e494e4de567eb65`.

## Contrato demostrado

Purchase acepta sólo el `RATES_RECEIPT.json` y `provider-response.json` V148 cuyos hashes coinciden, dentro de la ventana oficial de diez minutos. El `requestToken`, `rateId`, carrier, cargo y opciones documentales proceden de la respuesta cotizada; la selección humana liga el receipt exacto, owner, timestamp y máximo por moneda. Se rechazan additional inputs aún no resueltos.

Antes del POST se persiste un attempt con idempotency key UUID determinista. Éxito exige `shipmentId` y `packageDocumentDetails`, preserva respuesta y receipt; cualquier excepción desde el límite de llamada genera `PURCHASE_UNKNOWN_EFFECT.json`, `automatic_retry_authorized=false` y reconciliación obligatoria. Nunca completa la entrega.

## Gates ejecutados

- 16/16 pruebas tracking + rates + purchase PASS desde el árbol autor y desde Markdown materializado.
- Construcción de `PurchaseShipmentRequest`, `RequestedDocumentSpecification` y `DocumentSize` oficiales PASS.
- Negativos: receipt/response alterados, expiración, carrier, cargo, additional inputs, documentos no ofertados, output existente y timeout ambiguo.
- `compileall` y `pip check` PASS.

## Condiciones que permanecen cerradas

No se ejecutó una cuenta Amazon real. Developer/app/Shipping scope, shipping business, carrier/región, sandbox, costo, privacidad, rate-limit observado y reconciliación live siguen obligatorios. Cancelación, recuperación post-compra de documentos, handover físico, carrier universal y producción no quedan demostrados por V149.
