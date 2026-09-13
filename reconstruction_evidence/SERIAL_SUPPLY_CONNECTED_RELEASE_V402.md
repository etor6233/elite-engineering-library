# V402 — J2 supply conectado e incorporado

PROVEN_LOCAL para backend/API/host de referencia LIBRARY_INFRASTRUCTURE.
GO-CONNECTED-SERIAL-SUPPLY0.1.0,18archivos nuevos; cinco archivos de tres owners
actualizados: supply0.18.0, shared approval0.5.0 y app1.14.0.
Biblioteca183packs/1932blocks (AUTHORED1632/ADAPTED159/VERBATIM141).
Franquicia87/1166archivos reconstruidos exactos. No cierre global ni producción.

Demanda/PO con cantidades→confirmación proveedor/fábrica→producción por serie→
calidad distinta→ASN/manifiesto dividido→recepción parcial→cuarentena→revisión→
stock disponible. Unidades planeadas3, registradas4, una rechazada/reemplazada,
dos ASN, tres recibidas/revisadas/disponibles,30versiones. El total90000ARS fixture
no se recalcula. Stock/FIFO/ATP y transiciones siguen sus owners existentes.

POST/GET HTTP reales y PostgreSQL18.6 con69migraciones: ocho comandos concurrentes
producen1efecto y7replays; falla de outbox revierte snapshots de15tablas. Dos
respuestas ya confirmadas se pierden y recuperan por GET actor/command/hash sin
otro POST aceptado. Cambio de actor/payload, organización/rol incorrectos,
sobreproducción, manifiesto ajeno y saltos de aprobación/stock/PO son rechazados.
13límites de transporte incluyen duplicados JSON/case, actor, ruta, scope,
32768bytes y versión9007199254740993 preservada como string.

Host optativo exige hash exacto de política fixture y15guards activos; apagado
no lee política. Down0072 poblado rechaza conservando1plan/30pasos/4factory/3stock;
base vacía down/up PASS. Fuzz nativo finito2s,7semillas,130427ejecuciones PASS;
UTF-8 válido preserva identidad serial al serializar.0dependencias añadidas.
El fixture de identidad HTTP asigna principals: no afirma criptografía JWT/IdP.

Cuatro perfiles exactos:87/1166,40/570,57/796,88/1174. Backend/serverless compilan
la dependencia de los owners modificados;0tests runtime repetidos allí.
FAIL841 fue sólo un ciclo de imports de pruebas, resuelto con postgres_test
como garantía. RED conservado; implementación productiva no cambió por ese error.

G0: Existing admitted Operations, Approval and BC-derived Inventory owners selected; no new upstream acquisition.

G1: All23delta files are AUTHORED transactional/interface/configuration/test glue; existing source/notices unchanged.

G2: Required J2 blueprint bound to fixed explicit serialized-unit reference policy, no invented pricing or legal rules.

G3: One purchase/stock/approval owner; existing source transitions reused transactionally. No second ledger.

G4: PG18.6/69migrations: quantities, rejection/replacement, split ASN, receiving/quarantine and distinct review through version30.

G5: Scoped HTTP, malformed input and exact int64, two lost-response GET recoveries, generic API/SQL bypass denied.

G6: 15table rollback snapshots, eight concurrent commands1new7replay, restarted repository,130427finite fuzz executions.

G7: Hash/15enabled guards before activation; populated down refuses intact, empty down/up passes; PG stopped.

G8: Four exact reconstructions87/1166,40/570,57/796,88/1174; two narrower profiles compile only, no repeated runtime suite.

Hashes por archivo, packs y receipts: SERIAL_SUPPLY_CONNECTED_RELEASE_V402.json.
Guía materializada: docs/SERIAL_SUPPLY_REFERENCE.md.
T2802 conserva J3 REQUIRED; T2804 conserva el frontend por rol. No se usa este
PASS para sustituir supply-chain/identidad de composición, ops ni release firmado.
