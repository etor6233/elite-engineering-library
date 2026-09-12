# Amazon SP-API Supply Sources Execution Inventory — V155

## Resultado focal

`REBUILD_VERIFIED / CONDITIONED` para el adapter local Amazon Supply Sources v2020-07-01.

- Pack 0.1.0: 58.879 bytes; SHA-256 `b0fc2ce3b031d9ec26927a667e26467b58452a72aba6f92f847865bbb078e68f`.
- Perfil: 2 packs / 34 archivos; 1.336 bytes; SHA-256 `82bee18cf3f7ad060b57ba7861c64409fe630ab4946b5fa6f095412020443eae`.
- Materialización: 9/9 archivos con igualdad SHA-256 entre árbol autor y reconstruido.
- Runtime: wheel oficial fijado `amzn-sp-api==1.11.1`; `pip check` PASS.
- Pruebas: 14/14 PASS en árbol autor y 14/14 PASS desde Markdown.
- Biblioteca: `VERIFY_LIBRARY_PASS`; 92 packs / 997 archivos materializables / 507 Markdown / 40 perfiles de composición.
- Auditoría ejecutable: `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`; 92 packs / 121 fuentes oficiales / 12 adapters de proveedor.

## Autoridad y procedencia

Amazon documenta Supply Sources v2020-07-01 para configurar información y capacidades de supply sources como tiendas y depósitos, y publicó Multi-Location Inventory para aportar ubicación, inventario y capacidad real al promise system. El adapter usa los métodos oficiales `get_supply_sources`, `get_supply_source`, `create_supply_source`, `update_supply_source`, `update_supply_source_status` y `archive_supply_source`, junto con sus modelos oficiales.

- SDK: `amzn/selling-partner-api-sdk` commit `8e792ae345a8d334ccdbdd03181f05f040e6a4fc`.
- Modelos: `amzn/selling-partner-api-models` commit `8e429486005c4ebdce5099e48cc48515a65359bb`.
- Wheel: 3.557.108 bytes; SHA-256 `a88e4059c9fbfd454bfef67b82288a926f73cbdc608a8831f8383565736c70b0`.
- Licencia upstream: Apache-2.0.
- Los nueve archivos son integración/configuración `AUTHORED`; no se presentan falsamente como producto escrito por Amazon.

## Contratos demostrados

1. Perfil distribuido bloqueado: registro, app, rol/scope, marketplace/programa, dynamic sandbox, cuota/costo, tratamiento de ubicación/contacto y owner son obligatorios.
2. Inventario paginado sigue tokens únicos, obtiene cada source completo y rechaza ciclo, exceso, ID/código duplicado o divergencia summary/detail.
3. Receipts persisten ID/código/alias/status y hashes; dirección y contacto no se persisten.
4. CREATE/UPDATE/STATUS construyen `CreateSupplySourceRequest`, `UpdateSupplySourceRequest`, modelos anidados y `UpdateSupplySourceStatusRequest` oficiales; campos inventados fallan antes del efecto.
5. Cada mutación se liga a bytes exactos de request/inventario y snapshot actual; un inventory drift bloquea antes del provider.
6. `MUTATION_ATTEMPT.json` existe antes del efecto. Excepción o reconciliación divergente crea `UNKNOWN_EFFECT` sin error sensible y sin retry automático.
7. Cada éxito se reconcilia por `get_supply_source`; CREATE verifica código, UPDATE campos solicitados, STATUS enum oficial y ARCHIVE status `Archived`.
8. Archive exige primero `Inactive`; output existente nunca se sobrescribe.
9. `automatic_internal_location_activation=false`: Amazon no sustituye la autoridad del dominio interno.

## Límites honestos

No se ejecutó cuenta, dynamic sandbox o marketplace real. No demuestra location-level inventory, promise accuracy, transporte, sincronización con sucursales, carga, seguridad o producción. Supply Sources es una frontera Amazon concreta, no una red multi-location universal. El proyecto conserva gates live para CDN/WAF, IdP, proveedores, PostgreSQL/recovery, carga, seguridad ofensiva, deploy/rollback y aceptación empresarial.

Los gates que requieren red, credenciales, cuenta o target real permanecieron `SKIPPED` por falta de autorización/contexto; no fueron simulados ni contados como evidencia productiva.
