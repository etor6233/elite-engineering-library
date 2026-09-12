# Microsoft BC + Amazon connected carrier delivery — V174

## Resultado estrecho

V174 conecta un `sales.customer_shipment` ya posteado con un transporte durable del cliente, normaliza reportes provider mediante un schema exacto, conserva replay/evidencia/versiones y marca el fulfillment del pedido como `delivered` sólo cuando todos sus shipments de cliente poseen transporte entregado.

No crea ni acepta `sales.delivery_handover`. No afirma cuenta carrier, webhook, polling, custodia, firma, checklist ni aceptación del cliente. Esos límites continúan fail-closed.

## Autoridad oficial exacta

Microsoft BCApps, licencia MIT, commit `2eae56d704a1fd035d104f333602aea7091b7749`, tree `f9846fb1254c6311985c6131bb2af11c7f169e1b`:

| Path oficial | Bytes | SHA-256 | Claim estrecho |
|---|---:|---|---|
| `src/Layers/W1/BaseApp/Foundation/Shipping/ShippingAgent.Table.al` | 3,982 | `1f18ed96830b0bb8af68f22012661f886f48b49403b354eb5f923fa37aa54d73` | agente y tracking URL explícitos |
| `src/Layers/W1/BaseApp/Foundation/Shipping/ShippingAgentServices.Table.al` | 2,473 | `60d79a7c71e3d32e03a9e275160f7680810f3c8e4849a3e51cf160899462c9ee` | servicio separado del agente |
| `src/Layers/W1/BaseApp/Sales/History/SalesShipmentHeader.Table.al` | 54,368 | `9ea6772c4e470ee65cbb033a8e99ca6154607851ec3d8c0704f4c25deaa1dd41` | posted shipment conserva agente/servicio/tracking |
| `src/Layers/W1/Tests/SMB/O365ShippingAgent.Codeunit.al` | 23,865 | `d99105e91b43b234cd5b18d0f90f78e1ec758b85317a2b79c6877038e4caeda2` | persistencia order→posted shipment |
| `src/Layers/W1/Tests/SCM-Reservation/SCMPackageTrackingSales.Codeunit.al` | 130,076 | `96ed3a37ace2a44f14d835d6a75676f264900d41750544ebc48cbabc7a5df1bb` | package/lot/serial y cantidad del partial shipment |

Amazon se consume por los tres packs ya admitidos, derivados de SDK oficial `8e792ae345a8d334ccdbdd03181f05f040e6a4fc`, models `8e429486005c4ebdce5099e48cc48515a65359bb` y referencias oficiales:

- Shipping Tracking `0.6.x`: evidencia hash-linked/redactada, sin completar entrega automáticamente;
- Easy Ship Handover `0.1.x`: statuses exactos y `automatic_internal_handover_completion=false`;
- Fulfillment Delivery Evidence `0.1.x`: documento provider sin URL firmada/identidad y sin aceptación comercial.

El Go/SQL local es `ADAPTED`; el expediente de derivación es `AUTHORED`. Nada se presenta como Go de Microsoft/Amazon ni como carrier universal.

## Implementación incorporada

- fulfillment pack `0.4.0`: 15/15 archivos reconstruibles;
- seis bloques nuevos: repository de transporte, integración, migraciones `0039` up/down, test SQL y expediente de derivación;
- cuatro owners existentes actualizados: service/contract, unit tests, HTTP routes y HTTP tests;
- backend compuesto: 29 packs / 386 archivos;
- migración `0039` reemplaza la unicidad `NULLS NOT DISTINCT` defectuosa por índice parcial para referencias no nulas;
- `logistics.shipment_provider_report` es inmutable y único por shipment/version;
- endpoints autorizados crean transporte conectado y aceptan reportes provider normalizados;
- un `Delivered` provider puede cerrar fulfillment técnico del pedido, pero deja exactamente cero handovers aceptados.

## Pruebas ejecutadas

Toolchains fijados:

- Go 1.26.7 desde el runtime local verificado;
- PostgreSQL 18.6 x86-64 con checksums, cluster aislado y base nueva en puerto 55440.

Resultados:

- preflight de biblioteca completó sus gates offline y declaró explícitamente toolchains/access target ausentes en PATH;
- migraciones limpias `0001`–`0039`: PASS;
- 25 archivos de pruebas SQL: PASS;
- `0039` down → ausencia de schema nuevo → up → test focal: PASS;
- pack fulfillment materializado 15/15 y comparación SHA-256 byte a byte: PASS;
- compositor: suite negativa/positiva PASS;
- composición `ENTERPRISE_BACKEND`: 29 packs / 386 archivos;
- full `go test ./... -count=1` sobre árbol de trabajo y sobre composición desde Markdown: PASS;
- `go vet ./...` y `go build ./...` en ambos: PASS;
- pruebas domain/HTTP rechazan schema/status/hash inválido y aceptación automática inventada;
- `VERIFY_LIBRARY_PASS`: 94 packs / 1.103 archivos materializables / 531 Markdown / backend 386;
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: 94 packs / 121 fuentes oficiales y todos los componentes ejecutables seleccionados PASS. Los gates live sin cuenta, acceso, gasto o target permanecen explícitamente condicionados/SKIPPED.

Journey PostgreSQL focal:

1. stock FIFO real → receipt → put-away → binding comercial → pick registrado → customer shipment;
2. creación de transporte conectada y replay exacto sin duplicación; replay divergente rechazado;
3. `PickedUp` aceptado;
4. `AtOriginFC` y `AtDestinationFC` concurrentes: exactamente un ganador de versión;
5. `OutForDelivery` y luego `Delivered`;
6. replay exacto de delivery converge y evidencia divergente se rechaza;
7. order fulfillment final `delivered`, cuatro reportes provider y cero `sales.delivery_handover`.

## Fallos conservados y corregidos

`LIB-FAIL-1725`–`LIB-FAIL-1734` registran lecturas truncadas, el defecto real de unicidad nula, resultado de patch no observable, mínimo numérico omitido, metadata incompleta del pack y un parche rechazado por hashes inferidos. Los fallos materiales/procedurales quedan `REGRESSION_PROVEN` únicamente tras PostgreSQL limpio, integración, compilación y roundtrip 15/15.

## Condiciones residuales

- cuenta/región/sandbox/costo/cuota del carrier elegido;
- autenticación de webhook o polling/reconciliación, orden fuera de secuencia, outage y dead-letter;
- dirección, PII, retención, custodia, foto/firma y política de evidencia;
- customer shipment serial/específico para vehículos identificables;
- creación del handover, checklist versionada, identidad/autorización y aceptación/rechazo del cliente;
- carga/resiliencia, backup/restore, seguridad ofensiva, observabilidad, cloud/edge, rollout/rollback y aceptación empresarial target.

V174 mejora una capacidad ejecutable y conectada; no declara lista para producción una franquicia completa.
