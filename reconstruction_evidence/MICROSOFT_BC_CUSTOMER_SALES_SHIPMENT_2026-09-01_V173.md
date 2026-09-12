# Microsoft BC customer sales shipment — V173

## Resultado estrecho

V173 conecta `customer-order` placed/confirmed/paid/allocated → pick registrado → customer shipment inmutable → consumo de balance físico, composición UOM, reserva y FIFO cost → `shipped_quantity` por línea → fulfillment parcial/completo → outbox. El request exacto se reproduce sin duplicar, el divergente falla y dos requests distintos sobre el mismo pick producen un ganador.

No promueve carrier, entrega/handover, factura/ARCA, pago/crédito, stock serial, costo específico, service/production demand ni producción del sistema completo.

## Autoridad oficial exacta

- Repositorio: `microsoft/BCApps`.
- Commit: `2eae56d704a1fd035d104f333602aea7091b7749`.
- Tree: `f9846fb1254c6311985c6131bb2af11c7f169e1b`.
- Licencia raíz: MIT.

| Path oficial | Bytes | SHA-256 |
|---|---:|---|
| `src/Layers/W1/BaseApp/Sales/Posting/SalesPost.Codeunit.al` | 780,050 | `f0f3a8a4c914e2257a8186abe0c80a56b060b05a694d2e3cfabdfbf711103341` |
| `src/Layers/W1/BaseApp/Sales/History/SalesShipmentHeader.Table.al` | 54,368 | `9ea6772c4e470ee65cbb033a8e99ca6154607851ec3d8c0704f4c25deaa1dd41` |
| `src/Layers/W1/BaseApp/Sales/History/SalesShipmentLine.Table.al` | 78,739 | `607cc7697cd019fbecf82bfa7f0c8b287b5a444186184b2c4f44a0ee3eb1550d` |
| `src/Layers/W1/BaseApp/Sales/SalesWhsePostShipment.Codeunit.al` | 46,608 | `84121bd38555f3a7cba59dd97f985ae41af2da349359a60217bcab14ee7c3747` |
| `src/Layers/W1/Tests/ERM-Sales/ERMSalesOrder.Codeunit.al` | 376,840 | `817bc8bf7f3efbd46b6e0e095a79f4e7fb3dffe501fc86400ce44ee04c62829b` |

Los archivos oficiales crean encabezado/línea de shipment posteado, inicializan la línea desde Sales Line, trasladan cantidad y base, actualizan cantidad enviada y prueban su correspondencia. El Go/SQL local es `ADAPTED`; no es código Go de Microsoft ni una copia AL.

## Cambio canónico

- `GO-SUPPLY-FACTORY-INVENTORY-API` 0.16.0: 88/88 archivos.
- Cinco archivos nuevos: adapter PostgreSQL, migraciones 0038 up/down, test SQL y derivación.
- Ocho owners existentes actualizados: contrato/servicio, HTTP, packaging, transfer shipment y pruebas.
- `sales.customer_shipment`, `sales.customer_shipment_line` e `inventory.customer_shipment_allocation` son evidencia inmutable.
- `sales.customer_order_line.shipped_quantity` y `sales.customer_order.fulfillment_state` no sobrescriben el estado comercial.
- `consumeRegisteredPickPackaging` convierte composición no-base al delta base exacto que consume el trigger físico; el mismo owner corrige transfer shipment.
- FIFO bulk es el claim admitido; específico y serial fallan cerrados.

## Pruebas ejecutadas

Toolchains:

- Go 1.26.7: `go.exe` SHA-256 `5463fe58fa999d74420f00ee1b36d31c3da90a57ff204159f159859567ad61fc`.
- PostgreSQL 18.6, data checksums activos, cluster aislado en puerto 55439.

Resultados:

- updater: tests de nested fence, replacement literal, trailing LF, hash, atomic failure y traversal PASS;
- pack materializado 88/88; trece archivos modificados/nuevos byte-equivalentes al árbol verificado;
- composición `ENTERPRISE_BACKEND`: 29 packs / 380 archivos;
- PostgreSQL limpio: migraciones 0001–0038 PASS;
- 24 pruebas SQL PASS;
- 0038 down/up + test estructural PASS;
- full `go test ./... -count=1` PASS;
- `go vet ./...` PASS;
- `go build ./...` PASS;
- focal customer shipment y transfer packaging después de round-trip PASS;
- `VERIFY_LIBRARY_PASS`: 94 packs / 1.097 archivos materializables / 530 Markdown / backend 380 PASS;
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: 94 packs / 121 fuentes oficiales, incluyendo gates offline de navegador, calidad web, persistencia, documentos, integraciones y supply chain, PASS. Los gates live que requieren red, cuentas, gasto o infraestructura target permanecieron explícitamente `SKIPPED`/condicionados; no se reinterpretan como evidencia productiva.

Journey customer shipment demostrado:

- 3 BOX = 36 EA recibidas con costo unitario base 10;
- primer pick/shipment: 1 BOX = 12 EA, costo 120, fulfillment `partially-shipped`;
- replay exacto conserva shipment/line/version y fecha divergente falla;
- segundo pick: 2 BOX = 24 EA;
- dos posts concurrentes producen un ganador;
- estado final: `shipped_quantity=3`, fulfillment `shipped`, order version 4, dos shipments, costo posteado 360, dos reservas consumidas, físico/composición/costo restante cero.

Regresión transfer shipment:

- recibe tres BOX reales de dos EA;
- primer pick de una EA exige breakbulk explícito;
- segundo pick conserva una BOX intacta;
- shipment convierte la composición BOX al delta base y el trigger consume físico/base sin divergencia;
- transfer parcial, tránsito, recibos, lote y costo conservan las suites anteriores.

## Fallos y correcciones retenidos

`LIB-FAIL-1705`–`LIB-FAIL-1722` preservan búsquedas truncadas, paths/globs supuestos, brecha UOM real, FK/tipo SQL, identidad sobrecargada, compilación, normalización base, cleanups, fixtures de packaging, autocorrección de una lectura falsa y promoción separada del updater. Ningún intento parcial se presenta como PASS.

La brecha material `LIB-FAIL-1709` queda cerrada sólo porque la regresión no-base, el customer shipment, la base limpia, el round-trip y la composición canónica pasan juntos.

## Condiciones residuales

- La política pago/crédito previa a shipment se decide por proyecto; BC no autoriza inventarla como universal.
- El shipment logístico `logistics.shipment` sigue siendo carrier tracking, no documento de ventas.
- `sales.delivery_handover` sigue siendo entrega/aceptación, no posting de almacén.
- ARCA/facturación y pago continúan en owners separados.
- Vehículos serializados y costo específico requieren vínculo exacto de activo/source receipt.
- Carga, backup/restore, seguridad ofensiva, observabilidad, cloud/edge, canary/rollback y aceptación empresarial permanecen gates del proyecto real.

V173 no declara lista para producción una franquicia completa.
