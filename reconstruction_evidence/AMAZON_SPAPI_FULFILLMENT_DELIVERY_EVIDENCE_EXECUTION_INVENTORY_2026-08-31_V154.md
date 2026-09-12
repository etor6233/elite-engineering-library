# Amazon SP-API Fulfillment Delivery Evidence Execution Inventory — V154

## Resultado

`REBUILD_VERIFIED / CONDITIONED` para el adapter local de evidencia de entrega Amazon Fulfillment Outbound v2020-07-01.

- Pack: `PYTHON-AMAZON-SPAPI-FULFILLMENT-DELIVERY-EVIDENCE-ADAPTER` 0.1.0.
- Pack Markdown: 45.889 bytes; SHA-256 `8b6bec2902717b2ab21229131d8e140667020a2adfaac6a776680c6b6b43ff9e`.
- Perfil: 2 packs / 33 archivos; 1.394 bytes; SHA-256 `b54ac32a117e49cb859322807bea09ba069921f95678dbd076477a2b8bb47dad`.
- Materialización focal: 8/8 archivos; igualdad SHA-256 completa entre árbol autor y reconstruido.
- Runtime exacto: `amzn-sp-api==1.11.1`; `pip check` sin dependencias rotas.
- Pruebas: 12/12 PASS tanto sobre árbol autor como sobre árbol reconstruido.
- Gate estructural global: `VERIFY_LIBRARY_PASS`, 91 packs / 988 archivos / 504 Markdown / 39 perfiles; el perfil focal reconstruye 33 archivos.
- Gate ejecutable global: `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`, 91 packs / 121 fuentes / 11 adapters de proveedor; los gates live no autorizados permanecieron explícitamente `SKIPPED`, no se contabilizaron como producción.

## Autoridad fijada

Amazon publicó que `getFulfillmentOrder` puede devolver `deliveryInformation.deliveryDocumentList` con foto o firma para un shipment Multi-Channel Fulfillment entregado, y `dropOffLocation` para indicar dónde fue entregado. El adapter importa la API y modelos exactos del wheel oficial fijado y comprueba en ejecución la firma `get_fulfillment_order(self, seller_fulfillment_order_id, ...)` y los modelos `DeliveryDocument`, `DeliveryInformation` y `FulfillmentShipmentPackage`.

- SDK oficial: `amzn/selling-partner-api-sdk` commit `8e792ae345a8d334ccdbdd03181f05f040e6a4fc`.
- Modelos oficiales: `amzn/selling-partner-api-models` commit `8e429486005c4ebdce5099e48cc48515a65359bb`.
- Wheel PyPI: `amzn_sp_api-1.11.1-py3-none-any.whl`, 3.557.108 bytes, SHA-256 `a88e4059c9fbfd454bfef67b82288a926f73cbdc608a8831f8383565736c70b0`.
- Release note: `https://developer-docs.amazon.com/sp-api/docs/sp-api-release-notes?ld=SDESSOADirect`.
- Referencia: `https://developer-docs.amazon.com/sp-api/docs/fulfillment-outbound-api-v2020-07-01-reference`.
- Licencia upstream: Apache-2.0. Los ocho archivos materializables son glue/configuración `AUTHORED` y no se atribuyen a Amazon.

## Contratos demostrados

1. El perfil distribuido permanece bloqueado y exige registro, aplicación, rol/scope, MCF/marketplace, sandbox, cuota/costo, privacidad, retención y owner.
2. El request queda ligado por SHA-256 a un receipt empresarial real, pedido, marketplace y paquetes esperados.
3. Sólo se selecciona exactamente un documento requerido por tipo y paquete, bajo status terminal y host HTTPS exactos aprobados.
4. No se inventa un enum oficial: los tipos documentales deben observarse y aprobarse en el target.
5. Descarga sin redirect; tamaño, content type y magic bytes se validan para JPEG, PNG o PDF.
6. Las URLs firmadas no se persisten ni aparecen en el error de red probado.
7. Identidad del destinatario, atributos drop-off y tracking no se persisten; sólo queda tipo categórico y presencia de atributos.
8. Binarios con nombre SHA-256 y receipt/respuesta filtrada se publican mediante staging; directorio preexistente o aparecido concurrentemente nunca se sobrescribe.
9. Contenido duplicado, documento faltante, host/status no aprobado, orden/marketplace incorrectos, tamper, magic inválido u oversize fallan cerrados sin output parcial.
10. `automatic_business_delivery_acceptance=false` es inmutable: documento Amazon no prueba inspección física, aceptación contractual ni cierre empresarial.

## Límites honestos

No se ejecutó una cuenta Amazon real ni se afirma disponibilidad, exactitud del proveedor, soporte de marketplace, retención legal, aceptación empresarial o producción. El proyecto debe aportar credenciales por ambiente, listas observadas, sandbox y live reconciliation, además de CDN/WAF, IdP, PostgreSQL recovery, carga, seguridad ofensiva, deploy/rollback y aceptación empresarial reales antes de afirmar producción.

## Fallos retenidos

`LIB-FAIL-1475` a `LIB-FAIL-1480` conservan el constructor oficial incompleto del primer fixture, la firma de materializador supuesta, los fences generados incorrectamente, claims de router obsoletos, el patch documental atómico rechazado y el detector final inicialmente incapaz de distinguir snapshots históricos. Los PASS posteriores no borran esas lecciones.
