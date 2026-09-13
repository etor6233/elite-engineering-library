# V402 — identificadores opcionales de unidades, corrección canónica

FAIL840 CORRECTED_CANONICAL. La segunda unidad sin VIN fallaba en
production_unit_tenant_id_vin_key: el dominio admite VIN/batería opcionales, pero
NULLS NOT DISTINCT trataba dos ausencias como duplicado. El rojo se conserva;
no se fabricó un VIN para hacer pasar el fixture.

Migración0071 cambia sólo cuatro restricciones de VIN/batería a NULLS DISTINCT.
Preserva filas, unicidad de valores informados y serial obligatorio; no reescribe
0003 ni cambia APIs. Es glue AUTHORED, sin atribución a otra empresa.
Semántica normativa verificada en PostgreSQL18:
https://www.postgresql.org/docs/18/ddl-constraints.html#DDL-CONSTRAINTS-UNIQUE-CONSTRAINTS

PASS local: dos unidades sin VIN/batería y una con valores informados; tres
unidades de fábrica y tres de stock persistidas. Duplicados VIN, batería y serial
rechazados en ambos owners; siete eventos exactos, ninguno filtrado por fallos.
Down poblado rechazado, tres/tres filas y cuatro índices nuevos intactos; base
vacía down/up PASS. PG18.6 propio con68migraciones, detenido, cero SKIP.
Primer verde usó la extracción candidata de transacciones J2 con SQL inalterado.
El verde final usa1148archivos reconstruidos y operations.go canónico original:
la extracción J2 y su contrato candidato no se publican por este arreglo.

GO-SUPPLY-FACTORY-INVENTORY-API0.17.1 incorpora sólo dos migraciones y la regresión.
Cuatro perfiles exactos:86/1148,40/570,57/796,87/1156. Biblioteca182/1914bloques,
AUTHORED1614/ADAPTED159/VERBATIM141. No nueva dependencia o fuente externa copiada.

G0/G1: revisión de owner y provenance previas conservadas;0dependencias añadidas.
G2/G3: contrato Go opcional y semántica PG18, un owner de stock, migración aditiva.
G4/G5: reproducción roja y verde; valores ausentes/presentes y rechazo sin eventos.
G6/G7: down poblado fail-closed, empty down/up y datos preservados.
G8: tres archivos exactos en cuatro composiciones; runtime final sobre la referencia.
No se acredita J2 cantidad/ASN ni producción ni cierre global de TEST02 con esto.

Receipts/hashes/delta: SERIAL_OPTIONAL_IDENTIFIERS_V402.json.
