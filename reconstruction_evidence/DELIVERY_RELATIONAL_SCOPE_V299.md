# V299 — coherencia relacional del circuito de entrega

> **Errata de trazabilidad V317 (2026-09-08), prevalece sobre las menciones históricas de Go1.26.7 debajo.** Se retira la atribución histórica exacta Go1.26.7 no vinculada a un receipt de identidad por ejecución. No se sustituye por1.26.8 sin prueba del artefacto concreto. La ruta current-toolchain de V281 contiene1.26.8; official-toolchain contiene1.26.7. El cuerpo original queda conservado como historia. V317 revalida únicamente la composición actual y el scope allí enumerado con1.26.7 observado, no todos los snapshots previos. Ver [TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md](TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md).

2026-09-07. Continuación del checkpoint30, T2802/T2803/T2804. Se conserva
V295–V298, selección de credenciales pendiente y ARCA diferida. Corrección local
de un defecto encontrado al retomar entrega; no cierre de la entrega inicial.

## Qué cambió realmente

Journey0.10.5 modifica dos archivos existentes, sin dependencias ni migraciones:
`internal/platform/postgres/franchisejourney.go` y su integration_test.
La composición conserva67 packs/745 archivos, sin módulos duplicados.
El cambio es AUTHORED / CONDITIONED; no código copiado de Microsoft ni una
readmisión productiva de todo el pack por reputación o compilación.

Las FK existentes relacionan por tenant, pero no prueban que acta, pedido y stock
pertenezcan a la misma organización ni que el cliente del pedido sea el del acta.
Se reprodujeron nueve acciones indebidas: completar checklist, aceptar y rechazar
con tres clases de inconsistencia. Se persistían estado y outbox. Un segundo
experimento demostró que resolver una excepción también creaba un sucesor desde
un pedido de otra organización; los otros dos casos de ese experimento vieron
el efecto previo, no son tres bypass independientes demostrados.

Ahora las cuatro operaciones validan los vínculos dentro de su transacción,
bloquean el acta para escritura y pedido/stock para lectura estable. La defensa
no exige ni inventa pago total, crédito o financiación. Ante incoherencia devuelve
ErrConflict antes de respuestas, aceptación, excepción, resolución y outbox.
Se conservan la versión optimista y la atomicidad del evento de auditoría.

Fundamento de locking verificado en la documentación oficial de
[PostgreSQL18](https://www.postgresql.org/docs/18/explicit-locking.html).
Microsoft distingue la [liberación de documentos comerciales](https://learn.microsoft.com/en-gb/dynamics365/business-central/release-reopen-documents)
de su preparación/procesamiento posterior; esa referencia no decide la regla
financiera de esta franquicia ni convierte este cambio Go en código Microsoft.
No se adquirieron ni actualizaron upstreams.

## Evidencia reproducible

Staging: `%LOCALAPPDATA%/Temp/elite-v299-8b1c2130f91d4abab8a34b3049914528`.
Baseline V298 recompuesto; corrección devuelta al Markdown, recompuesta otra vez
en `rebuilt`. Los dos archivos son byte-idénticos al candidato probado.
Sólo esos archivos y el registro generado de materialización difieren del
baseline. Se verificaron70 archivos de web/gate/locks iguales al runtime V298
reutilizado; no se repitió build web ni se contaron sus105 tests como nuevos.

- Go1.26.7: test ./..., vet ./... y build ./... PASS. Los opt-in sin entorno
  no se presentan como ejecutados por esta suite raíz.
- PostgreSQL18.6/schema53 existente: TestFranchiseJourneyPersistenceIsolationAndReplay
  repetido tres veces desde la reconstrucción.12 negativos por ejecución:
  checklist/aceptar/rechazar/resolver contra organización de pedido, cliente
  de pedido y organización de stock incoherentes. Cero efectos indebidos.
- 16 aceptaciones concurrentes por ejecución: un éxito, quince conflictos y un
  evento durable. Fallo real de inserción de outbox conserva presented/version2
  y permite recuperar. Un pool nuevo observa accepted/version3 y un evento.
- Dos intentos de mutar relaciones bajo el lock devuelven55P03 con lock_timeout
  de50ms; se revierten ambos. Prueba de bloqueo, no benchmark ni SLO.
- Regresión de navegador:36 fases en Chromium desktop/mobile, Firefox y WebKit,
 106.12s observados. Revalida cotización/stock/solicitud local de pago de V298
  sobre el backend reconstruido, BFF y issuer RS256/JWKS sintético. No es una
  nueva prueba de interfaz de entrega, IdP real, cobro o proveedor externo.
- Base sintética loopback solamente; el cluster propio se detuvo al terminar.
  Pool nuevo no equivale a restauración, PITR, crash recovery o producción.

| Receipt local | SHA256 |
|---|---|
| relational-red.log | 13f558af75d0ef38700b2a274b3ba968eef558751afedfc653ad7ab2c17d20dd |
| resolution-red.log | 0f1f5c4db8fd6d4c6bcddc3700cba47b29327706090f37eefbffa0f5ed58f812 |
| rebuilt-pg.log | 6a5bae159a46946fa96443cf15ba8d82bbcbe1998224db3ed233eefc9c26202d |
| go-test.log | 64e0722478ec3ed2397b13030307fa8d988442d1906a4d1958b69446dd64acd7 |
| browser.log | 14cde9035f7952cc0036f3bc78271e8455e9f08161d07208ce8166717f76efd7 |

Reproducir desde el perfil integral vigente y schema53 en base sintética aislada:

```text
go test ./internal/platform/postgres -run '^TestFranchiseJourneyPersistenceIsolationAndReplay$' -count=3 -v
go test ./...
go vet ./...
go build ./...
go test ./internal/platform/httpapi -run '^TestQuoteAcceptanceBrowserPostgres$' -count=1 -timeout 22m -v
```

El primer comando necesita TEST_DATABASE_URL. El último además requiere web
construido/lock exacto en ELITE_WEB_ROOT y ELITE_QUOTE_E2E=1, ELITE_ORDER_E2E=1,
ELITE_PAYMENT_E2E=1. El harness rechaza bases que no sean loopback descartables
elite_confirmation_*. No pasar credenciales reales a este entorno.

FAIL447 corregido; FAIL448 era cardinalidad del fixture (cuatro actas originales
más nueve negativas), no se elimina su aserción: trece exactas y nueve marcadas.
FAIL432 conserva recurrencia de diagnóstico. FAIL384/385/405/423 no se cierran
por un PASS no relacionado; readiness/assurance integral permanecen bloqueados.

## Próximo tramo y decisión pendiente

Owner correcto: Journey posee acta/checklist/excepción; Commerce pedido/pago/stock;
Fulfillment transporte. No existe aún comando inicial conectado para el pedido
ordinario: los inserts actuales de producto crean sucesores o reemplazos.
Todavía faltan creación inicial y frontend, autorización de liberación, cobro
verificado/callback/reconciliación y cierre coherente de pedido/stock/entrega.
No resolver esas ausencias insertando un acta sintética ni copiando devolución.

Se preguntó qué condición comercial autoriza entregar. Su respuesta debe
registrarse en los owners de configuración/negocio existentes; no se asume
anticipo, pago total o aprobación manual. Hasta definirla no habilitar liberación.
El agente puede continuar con verificaciones locales independientes, pero el
checkpoint no permite declarar la franquicia lista ni producir el ZIP final.
