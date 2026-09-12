# V297 — recibo durable de solicitud de pago y referencias pendientes

> **Errata de trazabilidad V317 (2026-09-08), prevalece sobre las menciones históricas de Go1.26.7 debajo.** Se retira la atribución histórica exacta Go1.26.7 no vinculada a un receipt de identidad por ejecución. No se sustituye por1.26.8 sin prueba del artefacto concreto. La ruta current-toolchain de V281 contiene1.26.8; official-toolchain contiene1.26.7. El cuerpo original queda conservado como historia. V317 revalida únicamente la composición actual y el scope allí enumerado con1.26.7 observado, no todos los snapshots previos. Ver [TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md](TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md).

Fecha: 2026-09-07. Continuación del checkpoint28, T2802/T2804 del roadmap vigente.
V295 y ARCA diferida se conservan. Cierre local acotado, no cobro ni promoción.

## Implementación y procedencia

GO-COMMERCE-PRICING-PAYMENT-API 0.4.0 → 0.5.0. Cinco archivos Go existentes
cambiados y dos migraciones aditivas 0053, siete hashes de fuente iguales a su
reconstrucción. Ningún SDK, dependencia, cuenta o proveedor nuevo. AUTHORED;
no se presenta como código copiado/certificado por una empresa externa.

- RecordPaymentIntent reutiliza el owner, tabla y outbox existentes. Bloquea
  el pedido autorizado, coteja importe/moneda/estado server-side para crear y
  devuelve el registro durable en reintentos de la misma intención.
- Tenant/proveedor/Idempotency-Key delimitan la clave. Pedido, organización,
  moneda e importe deben corresponder; divergencias son409, sin divulgar recibos
  ajenos. Una transición posterior no impide recuperar el ID/estado actual.
- CreatePaymentAttemptAs toma el actor del principal HTTP. No acepta actor del
  cuerpo. Un repositorio sin el contrato nuevo falla cerrado. El caller interno
  legado mantiene error ante duplicado y no adquiere auditoría de actor por magia.
- Un solo payment.requested por intento; intento y outbox atómicos. La respuesta
  no se almacena en caché. Nunca se emite una petición a un proveedor desde el test.
- Migración0053 deja intacta0003: NULL no representa un identificador compartido
  del proveedor. Mantiene unicidad de referencias no nulas y de clave idempotente.
  Downgrade atómico falla si hay múltiples NULL; no borra ni inventa referencias.

Fuentes oficiales consultadas el mismo día, claims estrechos:
[AWS Builders’ Library](https://aws.amazon.com/builders-library/making-retries-safe-with-idempotent-APIs/)
para identificación explícita de la intención y tratamiento del reintento;
[PostgreSQL INSERT](https://www.postgresql.org/docs/18/sql-insert.html) para
ON CONFLICT y [Unique Constraints](https://www.postgresql.org/docs/18/ddl-constraints.html#DDL-CONSTRAINTS-UNIQUE-CONSTRAINTS)
para la semántica NULL. Son fundamentos de la implementación propia, no su
aprobación. Se conservan runtimes y licencias del lock existente.

## Pruebas reproducibles

Staging: `%LOCALAPPDATA%/Temp/elite-v297-c9a2214ca6f44941b48387c663247cf3`.
Go1.26.7/pgx5.10.0/PostgreSQL18.6, loopback, datos sintéticos exclusivamente.
El perfil integral recompone67 packs/745 archivos sin colisiones.
Base nueva elite_confirmation_v297 con53 migraciones desde cero.

1. Baseline real: segundo POST con misma intención devuelve400 en vez de202.
2. Tras reparar replay, segundo pedido sin provider_reference devuelve400;
   se corrige0053, no se simula una referencia para eludir la restricción.
3. Gate reconstruido TestPaymentHTTPDurableReplayPostgres: primer recibo y replay
   iguales;16 reintentos concurrentes;16 creaciones concurrentes para otro pedido
   producen un mismo ID y un evento; importe/moneda/pedido/tenant/organización/
   permiso/anónimo/clave corta rechazados según su contrato.
4. Fallo inducido de unicidad del outbox después de INSERT del intento: cero
   registros parciales; siguiente intento con la misma clave funciona. Actor
   payment-operator comprobado en el evento, no aportado por JSON.
5. Cliente descarta primer recibo y lo recupera. No es una desconexión de red
   inyectada. Cambia estado durable a pending y pedido a cancelled en fixture;
   un pool nuevo recupera mismoID/version2/estado actual, sin nuevo intento.
6. Referencia no nula duplicada rechazada. Up/down/up sobre esquema vacío pasa;
   down sobre datos incompatibles falla y conserva la restricción vigente.
7. Suites focales commerce, HTTP y PostgreSQL pasan; suite Go raíz completa,
   vet y build pasan. Sin TEST_DATABASE_URL, tests PG opt-in se omiten: por eso
   se ejecutó además el gate PostgreSQL explícito contra la reconstrucción.
8. PostgreSQL propio detenido/reiniciado: snapshot ordenado de pedidos, pagos y
   outbox idéntico; SHA256
   `17c9f41f8c69ad26ea24004fcee55051a6cb21df06fc5688314a6d8db4a82f96`.
   No equivale a crash recovery/PITR. Cluster propio queda detenido.

Para repetir: componer FRANCHISE_COMPLETE_PACK_PLAN en carpeta ausente, aplicar
las53 migraciones sobre DB sintética loopback cuyo nombre comience
elite_confirmation_, establecer TEST_DATABASE_URL sólo en ese entorno y ejecutar:

```text
go test ./internal/platform/httpapi ./internal/commerce ./internal/platform/postgres -run 'TestPayment|TestCommerce|TestOperationRead' -count=1 -v
go test ./...
go vet ./...
go build ./...
```

Los handlers usan un principal verificado sintético inyectado; esta regresión
no prueba JWT/JWKS, navegador, BFF ni transporte HTTP real. Se conserva evidencia
OIDC/frontend de V296 sin afirmar que cubre la nueva escritura de pago.

| Receipt local | SHA256 |
|---|---|
| baseline-replay.log | 4f252ac56ee5be7c8a70d26997117e5419ee1c96e58b27bc296ecac614ab26cf |
| baseline-reference.log | 9e2a8e7e0191d9fcaaf29587c63c12eea6251f5639a2c4dc90bcf495d7ff29a0 |
| rebuilt-payment.log | 73b6785800211ffd5fce0e055fed5ce650eb6b63916ff3789c82cce492758e94 |
| go-all.log | d1481801b8df225d72ffae7d52ba137113ba9d521e86b3610301ca79fc4ab186 |
| downgrade-blocked.log | 87ef57532019c2d8ad0ea4cc9c81796394de213b5391397a8bc51bb49a9f1cc5 |
| restart.log | ce149400e6429969e11aeeed4cb9805838cc547a97668c61214b77b395e6169b |

## Límites y siguiente paso del mismo roadmap

No se ha completado el recorrido de pago: falta UI/BFF con clave estable,
selección explícita y allowlist de proveedor, adapter de solicitud/cobro, callbacks
autenticados y reconciliación, sandbox y política comercial de liberación.
No hay garantía de un solo cobro entre claves distintas ni autorización de
procesar payment.requested externamente. No arrancar consumers con estos fixtures.
No hay todavía creación inicial de entrega del pedido ordinario conectada;
los flows de reemplazo/devolución no la sustituyen. Capacitación, operación y
gates target continúan según owners; T2802/T2804 no pasan a completos.

Próximo paso ejecutable: continuar el cierre de pago UI/BFF/HTTP sobre este recibo
con configuración explícita desactivada por defecto y validación de selección,
sin activar proveedor no elegido. La decisión mínima previa al adapter live es
Mercado Pago o Stripe, cuenta titular y ambiente sandbox. La regla de entrega
requiere decisión comercial; no deducir financiación, crédito o autorización
de entrega de este PASS. ARCA sigue diferida; no solicitar sus credenciales.

FAIL441/442 reparados; FAIL443 fixture corregido. El primer wrapper de composición
interpretó LASTEXITCODE heredado como fallo pese a materializar745 archivos:
se comprobaron salida efectiva y siete hashes, no se recompuso inútilmente.
Los errores de diagnóstico/patch se registran como recurrencias de FAIL432.
