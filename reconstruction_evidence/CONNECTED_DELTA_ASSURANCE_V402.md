# V402 — evidencia candidata de admisión y assurance del delta conectado

**Estado: CANDIDATE_ASSURANCE_NOT_FINAL_ADMISSION.** Unidad de auditoría: pago alojado → notificación autenticada → observación reconciliada → handover inicial → checklist/aceptación → recibo comercial, más la consulta de checkout en portal. Alcance de biblioteca local, fixtures HTTP oficiales y PostgreSQL real. No acredita producción live, cierre completo de T2802, TEST02 completo, 48/48 ni READY global.

Esta tarea reunió evidencia existente y calculó digests; **ejecutó cero pruebas** y no modificó canónico, estado ni composición. Los manifests son históricos. El inventario adjunto fotografía 82 archivos del delta y verifica que coinciden con la referencia congelada. La reconstrucción de root demuestra 915/915 archivos exactos; sus nuevos receipts se conservan separados de los históricos.

Inventario legible por máquina: `CONNECTED_DELTA_ASSURANCE_V402.candidate.json`, SHA-256 `4afb4ab9bacf7bc21ff5334b709fbb1a7630bf0131559bd9425088ca791082cc`. Contiene cada archivo/rol, digest observado, referencias de manifests, fuentes normativas y evidencias. Los 82 archivos inventariados tienen `frozen_revision_verified=true`: coinciden con el archivo reconstruido y el digest de `connected-frozen-reconstruction.json`. Esto acredita identidad de bytes; el alcance de cada prueba sigue siendo el de su receipt.

## Claim y fronteras

El SDK crea/retrieves la operación del proveedor; el callback sólo aporta una señal autenticada. El worker verifica scope, identidad de cuenta/conexión, modo, pedido/intento, moneda, importe y generación antes de proyectar estado. La preparación usa observación no held y vigente, pedido/cliente y unidad reservada del mismo scope. El perfil materializado fija hash, revisión y opciones admitidas: rev1 elegibilidad de lectura; rev2 `COMMIT_COMMERCIAL_RELEASE_RECEIPT`. El recibo es histórico e inmutable; su vigencia actual se comprueba otra vez. No hace una expedición física, un asiento fiscal, un movimiento bancario adicional ni inventa reglas de crédito/precio.

La UI/BFF obtiene la URL de checkout de un lector autorizado por tenant/organización/cliente. No convierte `status=complete` de una sesión en pago cobrado. La API manual de transiciones no es autoridad financiera cuando `ProviderObservedPayments=true`; sólo conserva cancelación `created→failed` sin referencia, sujeta al CAS del repositorio antes del claim de envío.

## Procedencia y fuentes fijadas

| Unidad | Tipo | Fuente/revisión exacta | Claim admitido y límite |
|---|---|---|---|
| Stripe Go SDK | DEPENDENCY_PIN, MIT | `github.com/stripe/stripe-go/v86@v86.3.0`; commit `a2df585a800a97fe8ec4ebf551b4449bdb3d90a1`; go.sum `h1:BKtYc3NtRa4EGzKAmp4jvl5q7kk2rwMZ+llF18N5vHI=` | APIs oficiales checkout/payment/account y verificación de notificaciones. No autoría del workflow local. |
| Mercado Pago Go SDK | DEPENDENCY_PIN, MIT | `github.com/mercadopago/sdk-go@v1.14.0`; commit `f910ee53fbb6819e435eaf3d0f800cb1fe74ae09`; go.sum `h1:3PYp9GPa+iysx2lcaKbpBEkXgEw4IIpbw3T6Jhl0IzI=` | APIs oficiales preference/payment/user, Search/Get y verificación. No firma QR ni aceptación live inferida. |
| SDK source/license receipts | Fuente oficial inspeccionada, no código de negocio copiado | EV01; licencia Stripe SHA `09db8d93f1ec60eb9f2be32ce2e10e440b35e209ea9fa505e84e788af938670a`; MP SHA `a227666680a292d0d96c78eb832076508d3de5d50366e3757686ed4ae8839597` | Los paths descargados coinciden con el cache de módulos fijado. SDK MIT y notices deben mantenerse en release. |
| Recuperación MP | SDK fijado + AUTHORED glue | EV18; snapshot oficial Search SHA `ccc976806b0313970517bfa36b767817097e2c9c86b6d82297bd0d3b70a0731d` | `external_reference` + site MLA, cardinalidad única y GET de identidad/monto/cuenta/modo/expiry. Search cubre los últimos 90 días; cero/ambigüedad queda unknown y no genera segundo POST. |
| Código local de composición, configuración y pruebas | AUTHORED inevitable glue, LicenseRef-Workspace-Owner donde el pack lo declara | Tabla de archivos e inventario JSON | Conecta owners y contratos existentes. No se atribuye a Stripe, Mercado Pago, Microsoft, Google ni otra empresa. |
| Owners de negocio reutilizados | Provenance original separada | `commercial-release-manifest.json/reused_owner_sources` y manifests canónicos existentes | No se reclasifica el owner entero como AUTHORED/ADAPTED por este delta. Supply/fulfillment inspeccionados no prueban expedición física en este tramo. |

La necesidad del glue es concreta: ningún SDK de proveedor conoce los IDs de pedido, tenant, reserva, inbox, outbox, checklist o perfil de esta composición. Los mappings, fences, contratos y proyecciones unen esas fronteras. Un algoritmo comercial nuevo o una política arbitraria no entraría en esta justificación: requiere su propia admisión. La revisión de todos los demás cores de la biblioteca queda fuera de esta evidencia.

## Matriz G0–G8

| Gate normativo | Evidencia local ya disponible | Límite / pendiente de la revisión congelada |
|---|---|---|
| G0 Identidad y procedencia | Commits/module checksums y receipts oficiales; hashes por archivo en manifests y este inventario. EV01, EV35, EV36 | **Inputs y bytes del delta verificados.** EV39 demuestra 915 archivos exactos; este candidato contrastó sus 82 archivos contra esos outputs. El manifest preliminar sigue siendo histórico y no reemplaza el receipt final. |
| G1 Licencia y reutilización | SDKs MIT, textos y digests oficiales; glue propietario declarado; no nueva dependencia en rev2/recuperación. EV01, EV03 | **Obligaciones identificadas.** Root debe generar notices/THIRD_PARTY_NOTICES/SBOM del artefacto final. No se afirma que esa propagación esté terminada mediante receipts de SDK aislado. |
| G2 Alineación con autoridades | Scope V402 del usuario + norma de admisión + ejecución; contratos explícitos de referencia y perfil; observación financiera y owners conservados. | **Mapeo candidato.** Invariantes concretos abajo. No citar fama como autoridad; las condiciones locales faltantes no se renombran credenciales. |
| G3 Arquitectura observable | DTOs, store/fence, inbox/job, callback, SDK GET y loader; migraciones 0055/0056/0057, docstrings y docs. EV34, EV35 | **Arquitectura local trazable y reconstruida.** Wiring/loader/host del snapshot congelado tienen PASS identificados en EV46. No certifica host productivo ni red externa. |
| G4 Correctitud | SDK HTTP fixtures, PostgreSQL, concurrencia/replay/rollback, MP sin payment metadata, Stripe URL null, recovery único, invariantes y 2 fuzz targets. EV16, EV22, EV19, EV27, EV28 | **G4 del delta conectado local PROVEN en composición congelada:** EV42/43 = 16 tests PG PASS, cero SKIP, vet/build exit0 y PG parado. EV45/46 = 25 tests unit/HTTP/host PASS enumerados, incluido guard financiero; existe un SKIP histórico explícito, por lo que el gate agregado cero-SKIP de units es FAIL y ese test no recibe crédito. No cierra toda T2802. |
| G5 Seguridad y privacidad | Callback firmado, raw body, account/mode/scope, generation/hold, autorización por objeto, URLhost + UTF-8, SDK logger sin payload, OSV de 3 runtime modules sin findings en snapshot SDK. EV07, EV03, EV02 | **Parcial local.** Falta el resultado final SCA/SAST/threat model/retención aplicable sobre composición congelada. Fuzz no sustituye SAST/DAST; no se hereda seguridad productiva ni se excusa con Daybreak. |
| G6 Rendimiento y resiliencia | 12 solicitudes comerciales concurrentes→1 recibo/11 replays; handover concurrente; retries/holds/bounded Search; cancelación/expiry/unknown con recovery sin segundo POST. EV20, EV19, EV15 | **Resiliencia focal probada, performance no cerrada.** No hay benchmark completo offered-load/goodput/p95/p99/high-water de este delta. El throughput de fuzz no es rendimiento del producto. T2809 carga local sigue siendo owner. |
| G7 Operación y evolución | Recuperación de respuesta perdida, rechazo de generation vieja, runbooks/docs, rollback transaccional y lecciones conservadas. EV20, EV36 | **Parcial local.** Root debe ligar supervisor/alertas/retención/restore/rollback y mantenimiento de inputs al artefacto final. Unknown es estado operativo explícito con retry bounded; no se convierte en accepted por timeout. |
| G8 Conversión a pack | SDK recovery 0.3.2 reconstruido exacto: 13 archivos, 12 intactos; perfiles parametrizables; candidatos de packs y manifest. EV31, EV30 | **Reconstrucción independiente exacta disponible (EV39); G8 final pendiente.** Verify global, notices, portable signed release y activación por bridge/START deben corresponder a esta revisión. No se concede REUSABLE_PACK/READY global aquí. |

## Assurance de implementación

Riesgo material **al menos high**: una asociación incorrecta podría aceptar un pedido impago o revelar un checkout de otro cliente. El fixture reduce exposición de ejecución, no la exigencia del owner reutilizable. No se altera el tier global que mantiene root.

| Dimensión obligatoria | Evidencia / oráculo | Estado de este candidato |
|---|---|---|
| functional_correctness | Misma orden, intento, moneda/importe, cuenta/modo; reservas; acceptance/checklist; partial refund no total, terminal no captured | PASS en PG congelado y casos unitarios enumerados; no toda la biblioteca |
| integration_contracts | SDK real contra HTTP fixtures oficiales + PostgreSQL real, no insert fixture captured/prepared en E2E conectado | PASS Stripe/MP en PG congelado EV42/43 |
| security_privacy | Object authz, firma, byte locks, provider identity, rechazo UTF-8, logs sin payload | Parcial: vincular SCA/SAST/security final; la UI no acredita privacidad completa |
| resilience_failure | Hold/generation, concurrencia, rollback atómico, timeout/POST perdido, ambiguity unknown | Probes PASS focales, no failover externo |
| performance_efficiency | Fences y límites observables; 3s/2 workers por fuzz target, inputs acotados | Carga de producto y percentiles no demostrados aquí |
| operability_observability | Estados unknown/hold, correlación por IDs, docs y logger controlado | Supervisor/alertas/retención de composición pendientes del owner operativo |
| recovery_rollback | MP Search/Get + restart entre SaveCheckout/fence; transacciones y migraciones de prueba | Recovery focal PASS; backup/DR del target no heredado |
| release_supply_chain | Go1.26.8 exacto, SDK pin/SCA, manifests y rebuild de SDK | Release final y firmas requieren evidencias separadas del build |
| usability_accessibility | BFF tests + typegen/typecheck y UI con estados de espera/redirección | Función BFF local PASS; navegador/a11y de la revisión final no demostrado por este receipt |

No se propone `implementation_assurance=proven` global en este documento. Artefacto y release requieren evidence IDs distintos en el manifest de ejecución; root preserva blockers reales hasta su resolución.

## Pruebas decisivas y lecciones

1. **Parcial/refund y terminal: RED sí observado.** `EV10, EV11` falló por convertir un refund parcial en total y permitir terminal→captured. `EV12, EV13` conserva los mismos casos en PASS. La proyección mantiene el owner de saldo/reembolso y no rehabilita terminales. El alias `mercado_pago`/`mercadopago` tiene RED/green propio y preserva identidades (EV25, EV38); no se reescribe provider en datos persistidos.
2. **Autoridad financiera manual cerrada y probada en alcance local.** `commerce.go` devuelve 403 `PROVIDER_OBSERVATION_REQUIRED` para authorize/capture/refund/dispute o fail de un request despachado cuando el modo observado está activo. `TestProviderObservationOwnsFinancialTransitions` define seis denegaciones y una cancelación previa permitida. EV46 contiene `TestProviderObservationOwnsFinancialTransitions` y sus siete casos PASS sobre el snapshot congelado. No se atribuye RED histórico a esa guarda. El primer log auxiliar EV44 no ejecutó HTTP/host; la prueba procede del run enumerado EV46, cuyo único SKIP externo se registra aparte. El root financiero es la observación reconciliada, no un permiso genérico `payment:write`.
3. **UTF-8 de URL: RED sí observado.** `EV26` reprodujo URL admitida cuyo byte inválido cambia al pasar por JSON. Fix `utf8.ValidString(c.URL)`; `EV27` PASS con 257.695 ejecuciones y 10 semillas. Prehash `cb1cccaac8449ec1c9f7295f6cac049da710ea526c894d251db2ab58dbfe8af7` → `1238275631751549d8a44c6b7f390a856644c35c80ac76723cbc05227bdb60b7`. Semilla #7 permanente; no evasión de test.
4. **Profile rev1/rev2 y bindings.** `EV28` PASS, 99.625 ejecuciones, 16 semillas: duplicados, hash/mode/scope y perfil comercial. El roundtrip JSON de un contrato no puede fabricar su binding privado; activación o campos públicos alterados no mantienen admisión. Perfil sin política o algoritmo incompatible no es READY.
5. **Expiry comercial: no afirmar RED inexistente.** `TestCommercialReleaseRejectsExpiryDuringOutboxWait` figura PASS en `EV21`. Prueba que expirar mientras espera outbox no deja un recibo válido; también están callback-hold/became-stale, rollback y lost-response. Esto es evidencia positiva de un negativo de contrato, no prueba histórica de que la implementación anterior fallaba.
6. **Recibo histórico versus autorización actual.** `EV23` demuestra quote→order→allocate→hosted checkout→callback firmado→job→GET→prepare→checklist/accept→commit con un recibo/un evento/un POST. Callback posterior o refund deja `current=false`; conservar la historia no autoriza entregar con observación nueva held.
7. **MP recuperación finita.** `EV19` prueba identidad única, crash tras guardar identidad y antes de cerrar fence, ausente/ambigua, metadata/amount/account/mode/expiry incorrectos. Nunca reenvía POST. La señal financiera puede quedar procesada antes de que la identidad de checkout cierre su retry; son receipts distintos. Search de 90 días no prueba ausencia de efectos antiguos.
8. **Runner scope.** `EV24` preserva FAIL de runner porque un regex seleccionó un test SDK sin su entorno y produjo SKIP. Se enumeraron cinco tests del owner y se corrió el SDK con su entorno propio; no se relajó cero-SKIP. El fallo de compilación del harness de fuzz por omitir el owner local `businesspolicy` también queda en su lesson; se copió dependencia exacta sin sus suites.


9. **Freeze y fallos de selección de harness, separados de defectos de producto.** EV40/41 conservan `connected-frozen-pg-final1`: varias pruebas se detuvieron en `requires disposable loopback handover database` por el nombre/guard de base de datos (FAIL812). EV42/43 reejecutó el alcance correcto y dejó 16 PG PASS sin SKIP; no es un RED de lógica comercial. EV44 sólo ejecutó un test SDK y corpus fuzz: HTTPAPI/host registraron `[no tests to run]` por el filtro anclado. EV45/46 corrigió la selección y obtuvo 25 PASS explícitos, pero incluyó además `TestPaymentHTTPDurableReplayPostgres`, SKIP por TEST_DATABASE_URL ausente; su gate agregado cero-SKIP fue rechazado (FAIL813). Se acredita cada PASS enumerado, se excluye el test omitido y se conserva el fallo del harness. No se repitió PostgreSQL para esa corrección de unidades.

Referencia congelada: `Temp/elite-v402-library-infra/connected-frozen-reference`. Pruebas PG/SDK: EV42/43; unidades manual guard/host: EV45/46. Ninguno acredita producción live.

## Pendientes exactos para incorporar a canónico

- Reconstrucción y PG/Go congelados ya enlazados (EV39, EV42/43, EV45/46). Conservar la revisión al refrescar BFF/build y nuevos deltas de UI; cualquier cambio material posterior exige su evidencia propia.
- Mantener los 25 PASS unitarios enumerados y el único SKIP histórico separados: `TestPaymentHTTPDurableReplayPostgres` sin TEST_DATABASE_URL no recibe crédito. No repetir los 25 para borrar el log; no declarar su corrida agregada como cero-SKIP. El PG conectado sí tiene 16 PASS/cero-SKIP.
- Terminar G5–G8 locales aplicables (SCA/SAST admitido, carga/operación/recovery, notices y release portable); no etiquetarlos como falta de credenciales.
- Conservar el checklist externo para futura activación: referencias de cuentas/secretos y modo/country-product/provider compatibles, aceptación de términos y endpoints. No se solicitan secretos ahora. Un resultado sandbox simulado no acredita operación live.
- Reconciliar estos evidence IDs con `implementation_assurance` y mapas de la revisión final sin sobrescribir REDs ni manifests históricos. No actualizar 45/48, T2802 o READY por este documento.

## Archivos y roles del delta

Cada hash completo y manifest histórico está en el JSON adjunto. `snapshot` coincide con la referencia congelada y EV39, **no** prueba retroactiva de que un receipt anterior ejecutó esos bytes. Los deltas sobre owners previos conservan la provenance del baseline.

| Archivo | Rol / procedencia del delta | SHA-256 snapshot |
|---|---|---|
| `.env.example` | Configuración, locks o documentación; AUTHORED metadata glue, licencias de terceros conservadas | `26b31cd5fc1b71a195daceee5c07b4facbb6d484bcc8e81e2ba8964267770e7e` |
| `cmd/electromobility-api/handover.go` | Montaje de runtime/handlers y validación de configuración; AUTHORED wiring glue | `fbd676682d373518ac9a2bf9825957a866e3cd1f96379b7c6add2cbea4b30a22` |
| `cmd/electromobility-api/handover_test.go` | Oráculos, fixtures y pruebas de frontera; AUTHORED verification glue | `ec95742ca6b10f946562c1da44370253f1465896d1c3f83d54b396ea608813f4` |
| `cmd/electromobility-api/main.go` | Montaje de runtime/handlers y validación de configuración; AUTHORED wiring glue | `0b726aab72cc2f2e033c658b699e5cf70b381e1499379980eead6f15f2b0235a` |
| `cmd/electromobility-api/payment.go` | Montaje de runtime/handlers y validación de configuración; AUTHORED wiring glue | `5e58c03b1def8995e9b8ccbb61c8c1b60d5527099b47831a61deb05fd3a7b35c` |
| `cmd/electromobility-api/payment_test.go` | Oráculos, fixtures y pruebas de frontera; AUTHORED verification glue | `a739eb14741686ebbec43e3b750fc8eb4a8229463ae321c7d30e3a4e32fda974` |
| `db/migrations/0055_payment_provider_observation.down.sql` | Persistencia, constraints e idempotencia de proyecciones; AUTHORED schema glue | `8184300a02e5bed5d1c804bc2b93b834eb2c34d7ee0307519c36dc6ae30e2735` |
| `db/migrations/0055_payment_provider_observation.up.sql` | Persistencia, constraints e idempotencia de proyecciones; AUTHORED schema glue | `8cdf10bc792d38d8750f81344d2a42070d75a091744d65f55756018cd1377b8c` |
| `db/migrations/0056_initial_handover.down.sql` | Persistencia, constraints e idempotencia de proyecciones; AUTHORED schema glue | `eb5f6d41bc17ea421a69077eb0f4cfef47fee3942b32c7785d98bd90a9687a6b` |
| `db/migrations/0056_initial_handover.up.sql` | Persistencia, constraints e idempotencia de proyecciones; AUTHORED schema glue | `500631f966f393a39ad8078fc9e7b5ce63af69b9c992124bb31ba15a8cdb8962` |
| `db/migrations/0057_commercial_release.down.sql` | Persistencia, constraints e idempotencia de proyecciones; AUTHORED schema glue | `2e3d61b2935649744d90c397d5fa908875ffbf1f3e25d1486563132211729fe4` |
| `db/migrations/0057_commercial_release.up.sql` | Persistencia, constraints e idempotencia de proyecciones; AUTHORED schema glue | `f258bb182d8f61ea636a29bc52ba0dc4ab4e3aa66043a9927bb3ce2154597baa` |
| `docs/commercial-release.md` | Configuración, locks o documentación; AUTHORED metadata glue, licencias de terceros conservadas | `9ebaa9b8a2f618dcb09e0771aa61eebce6f14636474b6b66d504c22e40f43239` |
| `docs/handover-host-activation.md` | Configuración, locks o documentación; AUTHORED metadata glue, licencias de terceros conservadas | `d452ce00db30dde395d6a7d5d2707876aff05facb679c48a8bb5bd942c454c66` |
| `docs/handover-profile.md` | Configuración, locks o documentación; AUTHORED metadata glue, licencias de terceros conservadas | `e730cb3c39b23516e13229641628d39d9391b028d28be2738ddd1ec8a56319e9` |
| `docs/initial-handover-reference.md` | Configuración, locks o documentación; AUTHORED metadata glue, licencias de terceros conservadas | `f6ae20dd9b368ff87cd150f285af78bb0a17d52fcc1d3eb9be3ac54a4df3c621` |
| `docs/payment-callback-reference.md` | Configuración, locks o documentación; AUTHORED metadata glue, licencias de terceros conservadas | `5e18ffa19bc05e45ef427eedb86954a479784dba5689c9c8796287e89b10489f` |
| `docs/payment-hosted-checkout.md` | Configuración, locks o documentación; AUTHORED metadata glue, licencias de terceros conservadas | `cc6df9d0f86d76f4d4ee3eb308a75a35d321c6966910fc40905b9e6039de9e11` |
| `go.mod` | Configuración, locks o documentación; AUTHORED metadata glue, licencias de terceros conservadas | `45f6843699be85cf99ecbecac2eb41b2dcfedd6483cf38cd61e95db80c6b2e66` |
| `go.sum` | Configuración, locks o documentación; AUTHORED metadata glue, licencias de terceros conservadas | `72eaa4679c14f54cfd95bd7cd9243c4dd836c12a4629a6f9fa461b5cd9be2d81` |
| `internal/franchisejourney/commercial_release.go` | Contratos, activación hash-locked y orquestación de owner existente; AUTHORED composition glue | `1971f0d99a7e68ec0d41171ae77967b71c7153b07225845a1d7fee7d565c5607` |
| `internal/franchisejourney/commercial_release_test.go` | Oráculos, fixtures y pruebas de frontera; AUTHORED verification glue | `2c8cac2695aa38bbb7199223a7b0e0fc97497d40624833d119d19ffcf4e7dabc` |
| `internal/franchisejourney/handover_preparation.go` | Contratos, activación hash-locked y orquestación de owner existente; AUTHORED composition glue | `4fa99592c9877938cac35f98235075c24361deb6f9b59051636f07502cee37b5` |
| `internal/franchisejourney/handover_preparation_test.go` | Oráculos, fixtures y pruebas de frontera; AUTHORED verification glue | `ab1c83aa726d99cf6f1ce4a4fd57a91ddfb84491a45bb785dda1535d374b5fc9` |
| `internal/franchisejourney/handover_profile.go` | Contratos, activación hash-locked y orquestación de owner existente; AUTHORED composition glue | `6c34bf0f19d54a115184821a6a68082828fb1a6dfad2fbe48df9e7d2e74c27ab` |
| `internal/franchisejourney/handover_profile_fuzz_test.go` | Oráculos, fixtures y pruebas de frontera; AUTHORED verification glue | `6053a13ae143e627169020944aad358ba2430456591e964e7f2f60498eb3bc97` |
| `internal/franchisejourney/handover_profile_test.go` | Oráculos, fixtures y pruebas de frontera; AUTHORED verification glue | `75ea6db0c29b095fb5daa810eb7595835f05e46a4fba46bc1c6b01cda91de0ce` |
| `internal/paymentbridge/checkout.go` | DTO, selección de SDK, worker/fence y validación del transporte; AUTHORED integration glue | `4da77f0226c3c018d5553d98f35ee464bc0dc2c7ec64be49931d523816ae30ea` |
| `internal/paymentbridge/customer_checkout.go` | DTO, selección de SDK, worker/fence y validación del transporte; AUTHORED integration glue | `1238275631751549d8a44c6b7f390a856644c35c80ac76723cbc05227bdb60b7` |
| `internal/paymentbridge/customer_checkout_fuzz_test.go` | Oráculos, fixtures y pruebas de frontera; AUTHORED verification glue | `183585b941bc6e2e59740320fe716bedf1b1359552f27e34f78393f2c17e968c` |
| `internal/paymentbridge/sdk_checkout_recovery.go` | DTO, selección de SDK, worker/fence y validación del transporte; AUTHORED integration glue | `f2f3566eedeef4481ef41744d64bd1d21d83652a420285d987f24a8078d3c89f` |
| `internal/paymentbridge/sdk_driver.go` | DTO, selección de SDK, worker/fence y validación del transporte; AUTHORED integration glue | `fc26f4bdae3576b5cafd7e90f6d101a8ae4ae7582a696110a7bfe79993e02b17` |
| `internal/paymentbridge/sdk_driver_test.go` | Oráculos, fixtures y pruebas de frontera; AUTHORED verification glue | `39aee9de29bcb33db266a00ab5502237867ff53ed4333e0c0cd7f553ab451144` |
| `internal/paymentbridge/webhook.go` | DTO, selección de SDK, worker/fence y validación del transporte; AUTHORED integration glue | `2dc4f478994076dab3a47e35c44af7b05c4be2b013cf3081de9ac822e1c6d001` |
| `internal/paymentbridge/webhook_test.go` | Oráculos, fixtures y pruebas de frontera; AUTHORED verification glue | `5f57fd0684d1f2fe933eaeff636e2d41d9fb9578f8daebc9a6cd9a6878f9edb0` |
| `internal/platform/httpapi/commerce.go` | Binding HTTP, autorización y cierre de bypass manual; AUTHORED delta glue | `5ba7a61bd62b41d03514b392062d02cdad514c7d017940f140add1779c430809` |
| `internal/platform/httpapi/commercial_release.go` | Binding HTTP, autorización y cierre de bypass manual; AUTHORED delta glue | `fa821602e9794279a35969ca0e6960196b78027809ca3a5fd7e5c99ff6127dd0` |
| `internal/platform/httpapi/commercial_release_test.go` | Oráculos, fixtures y pruebas de frontera; AUTHORED verification glue | `c538bf9c6f9be6bbca4340663417ab6247d7c8a6bb7b014f4fead0a2834d56a8` |
| `internal/platform/httpapi/handover_preparation.go` | Binding HTTP, autorización y cierre de bypass manual; AUTHORED delta glue | `0f34a467f979ac331f68e71fb60261a90504b805113bbdd01919c0ef2db9843f` |
| `internal/platform/httpapi/handover_preparation_test.go` | Oráculos, fixtures y pruebas de frontera; AUTHORED verification glue | `461e9e3294c21972289ec59b0a74ddcb5cfb0419af8b278c2b20164282108435` |
| `internal/platform/httpapi/payment_checkout.go` | Binding HTTP, autorización y cierre de bypass manual; AUTHORED delta glue | `9f4f08f9f594f9cd9ca8b7206a9e29e8b94321aab7b387b7dd55f4e193e5d2bd` |
| `internal/platform/httpapi/payment_checkout_test.go` | Oráculos, fixtures y pruebas de frontera; AUTHORED verification glue | `74bb67bcbb279fe1c6b43a3cc36f9f68f7f2073545e2c82c5f72e81b9a415655` |
| `internal/platform/httpapi/payment_observation_authority_test.go` | Oráculos, fixtures y pruebas de frontera; AUTHORED verification glue | `51ccf68a05177c18017bfff0ef37f2cf6d16d96cfddd1c83c32b9dff0b62d34a` |
| `internal/platform/postgres/commercial_release.go` | Composición transaccional de owners existentes/observaciones; AUTHORED delta glue, provenance del baseline conservada | `a084b79f6b4874197c54e9adba884ef5a4a9bb37fe03a22fbbbc46acdabaf0ef` |
| `internal/platform/postgres/commercial_release_connected_integration_test.go` | Oráculos, fixtures y pruebas de frontera; AUTHORED verification glue | `c29a4a094941c2571025d8c9ef03aec25f8385fcf70fed7e196429f8b5c73f83` |
| `internal/platform/postgres/commercial_release_integration_test.go` | Oráculos, fixtures y pruebas de frontera; AUTHORED verification glue | `b1a5de13ac5a2fb8e50520a12a5e9bee5b869e63dd3f05386ab6abcf2366b0b9` |
| `internal/platform/postgres/handover_preparation.go` | Composición transaccional de owners existentes/observaciones; AUTHORED delta glue, provenance del baseline conservada | `9da7d25451ff5fbad7c2270cf3ae71b331c4a3fc88b1852442b0b4785b7c27b1` |
| `internal/platform/postgres/handover_preparation_integration_test.go` | Oráculos, fixtures y pruebas de frontera; AUTHORED verification glue | `f6e5a4c8571d66a01b14fe61189782e87a4d484b1bd535f48ac2495d99efe509` |
| `internal/platform/postgres/outbound_delivery.go` | Composición transaccional de owners existentes/observaciones; AUTHORED delta glue, provenance del baseline conservada | `e7161276b124f7bb1a5c0666df59d33b9c0663beb8a8b22d65d47b20ce26a257` |
| `internal/platform/postgres/payment_callback_checkout_recovery.go` | Composición transaccional de owners existentes/observaciones; AUTHORED delta glue, provenance del baseline conservada | `365239d0421552a769c242bf6e4c73c4c12213020560b518fdc80653ed7797c3` |
| `internal/platform/postgres/payment_callback_processor.go` | Composición transaccional de owners existentes/observaciones; AUTHORED delta glue, provenance del baseline conservada | `e7fe432c2d5c492a817b39cf0828f699955e37cb3d15acf35e4f5987188b884b` |
| `internal/platform/postgres/payment_checkout.go` | Composición transaccional de owners existentes/observaciones; AUTHORED delta glue, provenance del baseline conservada | `d6aeaa87260bdedcda029c87118c0e8802f978e74a5e7134bcaac34e43dcc274` |
| `internal/platform/postgres/payment_checkout_integration_test.go` | Oráculos, fixtures y pruebas de frontera; AUTHORED verification glue | `7a104fbe2d840a8031bffd1ae724935fccae9aa5077ff685f62b9738fda71765` |
| `internal/platform/postgres/payment_checkout_read.go` | Composición transaccional de owners existentes/observaciones; AUTHORED delta glue, provenance del baseline conservada | `c05f4b63f75e16bbb9180e937e82b4ce88a4434c92c71bfa2fecace8b7264ec8` |
| `internal/platform/postgres/payment_checkout_read_integration_test.go` | Oráculos, fixtures y pruebas de frontera; AUTHORED verification glue | `44bf43ec7dd3d97d755994965d3c5ad43018f82abc021207d66501e9630ffe34` |
| `internal/platform/postgres/payment_checkout_recovery_integration_test.go` | Oráculos, fixtures y pruebas de frontera; AUTHORED verification glue | `6874b569386e8a7192c8bb3e05f8e39dafdbb23250bd12cb78a628906c757329` |
| `internal/platform/postgres/payment_connected_integration_test.go` | Oráculos, fixtures y pruebas de frontera; AUTHORED verification glue | `3e91c0e323360965b6bb284a6df95c4c8d40d355a060a29f655570419a7aed04` |
| `internal/platform/postgres/payment_dispatch_fence.go` | Composición transaccional de owners existentes/observaciones; AUTHORED delta glue, provenance del baseline conservada | `401c49f162c825e0b7f203a485eaed32388ab8b8f7d75c4214dab7530a4c565a` |
| `internal/platform/postgres/providerintegration.go` | Composición transaccional de owners existentes/observaciones; AUTHORED delta glue, provenance del baseline conservada | `77fd3573d0cf246ec068af19757a8cb58c148d977e874088f31681cb104b8dbd` |
| `official_payment_webhooks/README.md` | Configuración, locks o documentación; AUTHORED metadata glue, licencias de terceros conservadas | `63d63590cf39c65efa80488bf8aaf2ffb71ad2ba46fa47d1ebaf8e12a846bbaf` |
| `official_payment_webhooks/go.mod` | Configuración, locks o documentación; AUTHORED metadata glue, licencias de terceros conservadas | `fc6106ab8e8e4ab2dde71a6210e5d4b31cc674e62088947e29eb8e025387b141` |
| `official_payment_webhooks/go.sum` | Configuración, locks o documentación; AUTHORED metadata glue, licencias de terceros conservadas | `6f564169198e386498d553a0c9a1e7bce16b619d4e88d640429ad1f3ff74dc06` |
| `official_payment_webhooks/officialpayments/checkout.go` | Adaptación de DTO/llamadas SDK y verificación de identidad; AUTHORED glue, SDK separado DEPENDENCY_PIN | `1533b728ca09a06016aa8e0c7b2828cb9ea630b9b3f1bab6834f2d70999b5338` |
| `official_payment_webhooks/officialpayments/checkout_test.go` | Oráculos, fixtures y pruebas de frontera; AUTHORED verification glue | `754bef6d11f6d66574cf5555b95c111ed35cc70d69cbe44242797dfb7f21efa3` |
| `official_payment_webhooks/officialpayments/client.go` | Adaptación de DTO/llamadas SDK y verificación de identidad; AUTHORED glue, SDK separado DEPENDENCY_PIN | `9afae1052d458ab227a748b1a9f79c59e2dba1d91607bbd9cb17ab4132324c41` |
| `official_payment_webhooks/officialpayments/client_test.go` | Oráculos, fixtures y pruebas de frontera; AUTHORED verification glue | `6f0f8596f193a91e8063ef88ce2f5ebd54d5bbce10645145961fa38b1872c2d5` |
| `official_payment_webhooks/officialpayments/notification.go` | Adaptación de DTO/llamadas SDK y verificación de identidad; AUTHORED glue, SDK separado DEPENDENCY_PIN | `11ad6b9441b06c474ff2e5e7862c6fd118347b24d5993db8311b2ec631752da9` |
| `official_payment_webhooks/officialpayments/preference_recovery.go` | Adaptación de DTO/llamadas SDK y verificación de identidad; AUTHORED glue, SDK separado DEPENDENCY_PIN | `fd3dad0004905048125f0017a3e0efd0733aa782628c624303f0f1a12941d6a8` |
| `official_payment_webhooks/officialpayments/retrieve.go` | Adaptación de DTO/llamadas SDK y verificación de identidad; AUTHORED glue, SDK separado DEPENDENCY_PIN | `5148b8b7d9195cb300bbdba3b4b6f16f7b161958770b2a395d214a12cae22b26` |
| `official_payment_webhooks/officialpayments/retrieve_test.go` | Oráculos, fixtures y pruebas de frontera; AUTHORED verification glue | `0a6b926e852f8a5480f953dbf8c79df0622c70ab2b20f6405ed5f341c3da6913` |
| `official_payment_webhooks/officialpayments/verifier.go` | Adaptación de DTO/llamadas SDK y verificación de identidad; AUTHORED glue, SDK separado DEPENDENCY_PIN | `64ba006387b43979648306f4b0acb0e4de3db7ea152f33f69ea0becd6ed74eac` |
| `official_payment_webhooks/officialpayments/verifier_test.go` | Oráculos, fixtures y pruebas de frontera; AUTHORED verification glue | `e924989d2efbe1a3c2be3ecbbbc72a5bf4367e5e5b63e4a73902a6afea1b1939` |
| `return_refund_worker/docs/COMMERCE_PROVIDER_ALIAS.md` | Configuración, locks o documentación; AUTHORED metadata glue, licencias de terceros conservadas | `21d6e4003b6c911acdf690288114d1737e5af419296f15fa769792969dbe3e94` |
| `return_refund_worker/internal/refundworker/postgres.go` | Delta de compatibilidad de alias de proveedor; AUTHORED glue, baseline del owner conservado | `06a0915c86c29d30516ad05a4ab283a1f4f38624224e7a7736e5737230648ae3` |
| `return_refund_worker/internal/refundworker/postgres_integration_test.go` | Oráculos, fixtures y pruebas de frontera; AUTHORED verification glue | `f20113caa74aadd0a3fb28b94ecaf45c6cc49656dd2e7376bafd054b694e5d13` |
| `return_refund_worker/internal/refundworker/provider_alias_integration_test.go` | Oráculos, fixtures y pruebas de frontera; AUTHORED verification glue | `4a3a9d2b9ef20ac08006926b3ed48a5d7009171f14dcc700f8e46cf906626fc9` |
| `src/app/api/enterprise/checkout/route.test.ts` | Oráculos, fixtures y pruebas de frontera; AUTHORED verification glue | `1d9535a7f844b93ddec8c9de31f89907b7822981a68f33a8d2f5c70eccd1a7be` |
| `src/app/api/enterprise/checkout/route.ts` | UI/BFF hacia operaciones y permisos existentes; AUTHORED UI binding glue | `0dd98ce24feb982692b6dec6696c638caded8435d2c0bc3ec1b94fc94603ed6a` |
| `src/app/customer/page.tsx` | UI/BFF hacia operaciones y permisos existentes; AUTHORED UI binding glue | `277bcdfaeb46cf0482fabbc83be1d6b0b7330667583c611753bac2411c0f6cee` |
| `src/components/customer-checkout-actions.tsx` | UI/BFF hacia operaciones y permisos existentes; AUTHORED UI binding glue | `153b13fa7a9296e8d56014763c1a5cec44951b15da5f8e6cec79703e74388568` |
| `src/platform/payments/checkout.ts` | UI/BFF hacia operaciones y permisos existentes; AUTHORED UI binding glue | `512077333c00084c0fd256fc177a6c9c48f28c2b71df222cc0a5af3b5af7448f` |
| `tools/materialize_handover_profile.py` | Materialización determinista de configuración admitida; AUTHORED tooling glue | `a24c419c5e85c59d1e05aa172d6c477be773bd3f57b99e1eebb2aa5541241df1` |

## Índice de evidencia existente

| ID | Artefacto en staging | SHA-256 |
|---|---|---|
| EV01 | `payment-source-receipts-complete.json` | `ea8a210a7d96ff61e9b36ade2764df791f9acceb6e2d04bcf2e348651f8505c7` |
| EV02 | `payment-sdk-final-result.json` | `48fe4802c7eeb9522284ae46890d8ecd023429e7908b3cff2ba73453afc29773` |
| EV03 | `payment-logger-result.json` | `a948997537d479a60368d270eaee7c95e91c8ccea9ecd019fca3600fb6715360` |
| EV04 | `payment-sdk-driver-result.json` | `4035ed65c1a7dc152a3db42ec52cf4a39214749cc3881f3d21a1dcb125c5c3d8` |
| EV05 | `payment-sdk-osv.json` | `01b8c26c57808538d47521d42f2ff3d93f78cae489d480f97c08a60c98b9d200` |
| EV06 | `payment-bridge-test-result.json` | `a355bb0d045d9ee467280a0ab2ad892ec19211ebfe618a12022106bebe2c407d` |
| EV07 | `payment-wiring-go-result.json` | `ed744ced2b072bc73f2f7b1946415aa4bcdf6391004ec569024cd64b7f58d2f8` |
| EV08 | `payment-wiring-web-result.json` | `ac0124e49bc3200a9025ea0d5757b3de8a0b97a6834ba54abc840d48cac0f5da` |
| EV09 | `payment-reader-pg-read1/result.json` | `b2474da1ba7dedadf634f59fd2a055f32e858ee155a12bb2ee84dfb853c3518d` |
| EV10 | `payment-pg-semantic-red/result.json` | `9b1f8ce82786c3af2d6ec1c5ba841bbcb6006521f582df409af95c4b8bb8c8bf` |
| EV11 | `payment-pg-semantic-red/connected.log` | `af7f8aa5ba9464a8b8c9c143e0aa3b004cfcf762bd9f02a71dd712d4a40538f1` |
| EV12 | `payment-pg-semantic-green2/result.json` | `3d6ff283e5a57dcb950f21e0d66162ded79f3ecbd03c0f34210c2234a0dbd290` |
| EV13 | `payment-pg-semantic-green2/connected.log` | `4b3dd3d582187f96c8141186b58f2ebd6759f066dd4c7c818a35ef95c9a205a9` |
| EV14 | `payment-pg-admission/result.json` | `3cf7b733ce3d879bd2533d6b591df792efe668fa81a83703f37d4debbeea319d` |
| EV15 | `payment-pg-deadline-filter/result.json` | `afdc75db2e22800aa047ef983d30472222a4b881ca1d58a8b240917e1421b1b6` |
| EV16 | `handover-agent/connected-run3-mp/result.json` | `c45fd6f1d67c38084230ce6351b5142d60b001214de50e5cb9f47d005a433dd1` |
| EV17 | `handover-agent/connected-run5-profile-binding/result.json` | `ebf3a4c7edfb33da94add3422d7d918098193bcae8825b88a4f523bc64e48bac` |
| EV18 | `mp-recovery-agent/source-lock.json` | `3e59288a17ab6e8e2fd2b5ddced1fe961e8f940cb3301f8f2d33785a1516814a` |
| EV19 | `mp-recovery-agent/run1/result.json` | `612622a4b50bf2c0bebec9acb9ebc0cfeaf60d8afa577265497741e6bde196fe` |
| EV20 | `commercial-release-pg-final-scoped/result.json` | `e7302dd778d8e14b6bd70d8ab35479f80aaf906f5e597bbb8ebe71fc574e97a9` |
| EV21 | `commercial-release-pg-final-scoped/connected.log` | `9110e4581ae9a1b965ddca49743392dbacbf29af898a69c6cc5aa3764bc3392b` |
| EV22 | `commercial-release-official-pg-final/result.json` | `75b21466a48f989fda29776d3efeed43c081cc05b84d1aa2b8149e4c70c3d1e7` |
| EV23 | `commercial-release-official-pg-final/connected.log` | `8fb80b9cc011bd3d090946f0337b85b518ae1b90b0a49d878ed9e6ce11a6b790` |
| EV24 | `commercial-release-failure-lesson.md` | `c32b182ec3cca8df53ba3cd43667388934f2bca420eb0d2f0d5c4f4b078a2e05` |
| EV25 | `refund-alias-pg-2/result.json` | `c402425d90e1cee64ade3cbdd999545bc98eb77e8d91cad8bc14dd5ecc5388ea` |
| EV26 | `checkout-policy-fuzz-agent/checkout-run1/execution-receipt.json` | `f8e8e91884407beb0e955378e9d9ec29ab7f2cecd58641d48c8c43fd685522f6` |
| EV27 | `checkout-policy-fuzz-agent/checkout-run2/gate-receipt.json` | `44542de3fa14e02f0a5ef00b00e5873a7206043e1996feccd02c5c05dfc05c40` |
| EV28 | `checkout-policy-fuzz-agent/profile-run2/gate-receipt.json` | `f2e3900ac113b1fb258b3639dea292e2acdb508df4d4bc368c21fd037c53622a` |
| EV29 | `checkout-policy-fuzz-agent/PROJECT_FAILURE_LESSONS.md` | `b93c3653a4af31d912415888135655a57ea18f0df947300ee404ff0899ed4860` |
| EV30 | `connected-pack-candidate-manifest.json` | `789d3142b5bd05574b288d02de5da22a8226af84bb15ddf94a363e0a440c0914` |
| EV31 | `sdk-recovery-candidate-manifest.json` | `f0822cf4f10d815e08e7513a0d9879cef4dad64a54d96f9b8c885a810e5730d6` |
| EV32 | `payment-wiring-manifest.json` | `11402de4c2d32ebb46f51030068cd8c5797a203571a10f05782cf0df52ba9dd7` |
| EV33 | `handover-host-manifest.json` | `1967e4d145f87213a145b55f3164ddeb1b40ca498110d734ddfa3410e92f4168` |
| EV34 | `handover-agent/payment-connected-result.json` | `5cf582f839d142466836b1a29b52a30b9260041068c9fc73f65325d32e86005b` |
| EV35 | `commercial-release-manifest.json` | `9650840fc2a38ef7cdb9944b840cac1ad2d1c6e8bb0475d77b946e69effc44ea` |
| EV36 | `mp-recovery-agent/result.json` | `ddb721a8eb92afd03d969e2031ab6b3d88108e005c5f26050a0d7266c5070c54` |
| EV37 | `checkout-policy-fuzz-agent/result.json` | `cf5a0f39cac6be3546f19f9e1dc624b350a9236578ffc98536657bb12fb7a2ce` |
| EV38 | `refund-alias-manifest.json` | `90bfb623d516c22b34edc4d2137811e43ac8a62e96894fea2b9f2154b04de52e` |
| EV39 | `connected-frozen-reconstruction.json` | `9a8664d562f805d55fd4a076e8c17895e8cd223e1639c5aef437227c7383fbf3` |
| EV40 | `connected-frozen-pg-final1/result.json` | `176a0a354aea0e4d613d63b22b0736d99098850bd59fc55d95ef6aa19d71b06f` |
| EV41 | `connected-frozen-pg-final1/connected.log` | `a4f691fd4fd048b8cac1c7b903ba9797fe2a04c70a5ec0b72ddbf6529c29e2ae` |
| EV42 | `connected-frozen-pg-final2/result.json` | `b47bde8727143ea90eeb18f67c694e0ee860abb208f4fe7e227704004ce0b422` |
| EV43 | `connected-frozen-pg-final2/connected.log` | `8d7533d73bd87d6c15bd1a658abff3e8c758e0f8fa018f8bd92cac29b2fd8970` |
| EV44 | `connected-frozen-pg-final2/affected-unit-host.log` | `004fd337dcee3a8decdd4ba256a8c35442944dcc20783fb0f92673166ac55ca5` |
| EV45 | `connected-frozen-unit-final/result.json` | `4bf013eb754b219f6ca5b91bb88f984f6cfd2fc9e53553d8cf57f5caa14ea4e8` |
| EV46 | `connected-frozen-unit-final/run.log` | `1731129e053e1b7de8c0bf2092d27fb1074ddeba73d9ae982ec21bbda8334678` |

Fuentes oficiales originales y digests de snapshots figuran en receipts. Esta tarea no realizó adquisiciones ni pruebas nuevas. G5–G8 y assurance global permanecen abiertos hasta sus evidencias propias; no son faltantes de credenciales.

V402321 portability correction: Temp/ denotes the originating system temporary directory; only rendered machine-local locators changed. Original rendering SHA-256 734da8551196755b5934e3b2dba28764bcd4ffe69b0879e6794cacca6655e43f is retained outside the distributed source. Receipt/source hashes and historical outcomes above are unchanged.
