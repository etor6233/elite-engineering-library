# Amazon SP-API Easy Ship handover — execution inventory V153

Fecha: 2026-08-31. Resultado focal: `REBUILD_VERIFIED / CONDITIONED`.

## Autoridad

Amazon SDK Python 1.11.1, commit `8e792ae345a8d334ccdbdd03181f05f040e6a4fc`, modelos `8e429486005c4ebdce5099e48cc48515a65359bb`, Apache-2.0; wheel oficial SHA-256 `a88e4059c9fbfd454bfef67b82288a926f73cbdc608a8831f8383565736c70b0`.

Se inspeccionaron las firmas oficiales `EasyShipApi.list_handover_slots`, `create_scheduled_package` y `get_scheduled_package`, junto con `ListHandoverSlotsRequest`, `CreateScheduledPackageRequest`, `Dimensions`, `Weight`, `TimeSlot`, `PackageDetails`, `Item`, `ScheduledPackageId`, `TrackingDetails` y `PackageStatus`. Referencias: https://developer-docs.amazon.com/sp-api/reference/listhandoverslots, https://developer-docs.amazon.com/sp-api/reference/createscheduledpackage y https://developer-docs.amazon.com/sp-api/reference/getscheduledpackage.

Los nueve archivos materializables son `AUTHORED`; importan el SDK oficial exacto y no se atribuyen a Amazon. La API fijada no expone idempotency key para schedule ni una operación de cancelación inequívoca; el pack no las inventa.

## Inventario

- Pack 0.1.0: 49.795 bytes/SHA-256 `265482e01cb03fe61541b795b8178e541519411500e21dfb43f92d866d267209`; 9/9 archivos.
- Perfil: `AMAZON_SPAPI_EASYSHIP_HANDOVER_PACK_PLAN.md`, 1.243 bytes/SHA-256 `2aeb0a1298c5278f5c3903aebbddfa462f377bdf32e3405cb725729ecf2b4204`; composición esperada 2/34.
- Provenance global resultante: 838 `AUTHORED`, 37 `ADAPTED`, 105 `VERBATIM`; 980 bloques.

## Contrato demostrado

`slots` construye modelos oficiales de dimensiones/peso, exige marketplace aprobado, valida 1..100 slots únicos/futuros y conserva request/response/receipt hash-bound. Una aprobación separada selecciona exactamente un slot devuelto y 1..500 items dentro de límites oficiales.

`schedule` construye `CreateScheduledPackageRequest`, persiste el intento antes del POST y valida order/package/slot/status/tracking. Excepción o respuesta inválida produce `SCHEDULE_UNKNOWN_EFFECT`, desautoriza retry y exige reconciliación. `reconcile` consulta el paquete oficial: `PickedUp` se conserva como `provider_reports_handover=true`, pero `automatic_internal_handover_completion=false` permanece fijo.

## Gates focales

El pack materializó 9/9 con SHA exacto. Desde el árbol autor y desde el Markdown materializado pasaron `compileall`, `pip check` y 9/9 pruebas: firmas/modelos oficiales, slot inventory, duplicado/expiración, schedule exacto, slot inventado/tamper, timeout ambiguo, respuesta inválida, reconciliación `PickedUp` y no-overwrite.

## Condiciones abiertas

Cuenta, rol/scope, support table exacta del marketplace/operación, sandbox, rate headers, cuota/costo, privacidad/retención, concurrencia, reconciliación live, alarmas y aceptación siguen bloqueados. Easy Ship no es carrier universal. Cancelación, firma/foto, prueba empresarial de recepción, cross-location y target productivo quedan fuera de V153.
