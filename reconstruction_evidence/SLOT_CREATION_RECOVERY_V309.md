# V309 — recuperación de publicación de capacidad

> **Errata de trazabilidad V317 (2026-09-08), prevalece sobre las menciones históricas de Go1.26.7 debajo.** Se retira la atribución histórica exacta Go1.26.7 no vinculada a un receipt de identidad por ejecución. No se sustituye por1.26.8 sin prueba del artefacto concreto. La ruta current-toolchain de V281 contiene1.26.8; official-toolchain contiene1.26.7. El cuerpo original queda conservado como historia. V317 revalida únicamente la composición actual y el scope allí enumerado con1.26.7 observado, no todos los snapshots previos. Ver [TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md](TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md).

Mantenimiento T2802/T2803/T2804, FAIL457, desde checkpoint44 validado.
Fecha2026-09-08 UTC. Continuación del roadmap existente, sin nuevo módulo,
dependencia, migración, regla comercial, proveedor ni efecto externo.
Perfil67/745. Sin cierre integral ni promoción/ZIP.

## Fallo y corrección canónica

Con jornada sintética válida, el POST inicial de turno devolvía201 y el replay
idéntico409; había1 cupo durable. La restricción de unicidad impedía duplicación,
pero faltaban recibo por clave y actor HTTP en el evento. El sender podía afirmar
no aplicado tras perder la respuesta. El listado público excluye completos,
cerrados, fuera de rango o sin disponibilidad; no sirve como prueba de no creación.

Journey0.10.14 agrega CreateAppointmentSlotOnce y GET protegido por appointment:manage
en /v1/franchise/appointment-slots/result. Clave/hash, cupo y outbox con actor
autenticado se confirman juntos usando idempotency_record, scope franchise-slot.
HTTP201 crea y200 replay. Lookup une tenant/org/clave/completed/type/code201/recibo
con el turno exacto y cuenta reservas requested/confirmed en el mismo statement.
Recupera el estado vigente, incluyendo completo/cerrado, sin reabrir ni reservar.
Legacy interno conserva contrato y queda fuera del nuevo claim de auditoría HTTP.
Se mantienen jornada explícita, no superposición, límites de capacidad1..100,
duración≤8h y validación temporal existente. POST conserva esa validación temporal;
GET de recuperación no exige que el turno siga siendo futuro. No promesa tras
purga de la clave (expires_at24h del owner) ni cambio de retención.

Portal0.14.7 añade referencia mínima scoped, fence sincrónico, storage0POST al
fallar, timeout y recibo contrastado. GET recupera sin POST. Preparar otro turno
exige consulta positiva y ausencia de POST/GET activo; no reserva, cancela ni
duplica el previo. Muestra estado/cupos ocupados y ayuda slot-create-view/1.0.0.
No payload/fechas/claves de acceso en sessionStorage ni garantía multi-tab.
Browser gate0.1.19 conecta el nuevo recorrido a las19 fases previas.
Los cambios son AUTHORED y condicionados; no atribuidos a upstream ni admitidos
como sistema productivo por estos gates locales.

## Verificación

| Gate | Resultado y alcance |
|---|---|
| HTTP rojo→verde |201/409 con1 cupo pasa a201/200 mismo ID y actor slot-manager |
| Web |112 PASS/1 SKIP conectado explícito; build PASS |
| Go reconstruido | test/vet/build ./... PASS; skips de integración sin DB no contados como conexión |
| PostgreSQL |7 tests de persistencia/agenda/identidad,3 repeticiones PASS |
| Carrera/atomicidad |8 solicitudes→1 creación/7 replay; hash distinto rechaza; fallo outbox revierte cupo/clave; actor original preservado |
| Recuperación |aislamiento tenant/org/key, recibos inconsistentes rechazados, nueva pool; dos reservas llenan cupo, desaparece de lista pública pero lookup devuelve2/2; cerrado permanece cerrado |
| Browser reconstruido |80 fases/4: Chromium desktop/mobile, Firefox desktop, WebKit desktop PASS |
| Nueva fase |3 cupos/keys/eventos/actores por proyecto; service/delivery/test-drive, pérdida/JSON inválido/concurrencia, storage0POST, GET503, permisos, clave ausente, actor forjado, divergencia y fecha200d |
| Agenda |4 proyectos existentes PASS con sus permisos/asignación/capacidad |
| Reconstrucción |9/9 archivos byte-idénticos; composición67/745 |
| Reinicio |snapshot de cupos/reservas/keys/eventos idéntico tras restart del cluster sintético; no restore/PITR del target |

FAIL477: primera fixture no tenía jornada y devolvió409/409 con0 cupos; se
conserva ese log, luego se incorporó jornada por el owner y se reprodujo el fallo
real. FAIL478: positivo heredado sin clave actualizado junto con fake/actor,
sin quitar negativos. FAIL479: ISO UTC rellenado como hora local produjo turno
fuera de jornada; trace/SQL probaron el desfase. Se genera datetime-local en
timezone del navegador; no se ensancha jornada ni retira guard para pasar.

FAIL457 sigue OPEN para publicación/completado de checklist y recepción/disposición.
FAIL480: replay sintético reordenaba JSON frente al BFF; se conserva409 como
negativo y se verifica200 con bytes/orden exactos del contrato existente.
Ronda D canal/formato de historial preguntada sin respuesta; V295, liberación
comercial y ARCA diferida intactos. Pruebas locales no sustituyen seguridad
integral, carga, providers, observabilidad/restore operativo ni aceptación.

## Reproducción e integridad

Staging `%LOCALAPPDATA%/Temp/elite-v309-9f18cfdc01cc4eb6819dbd1674425f0e`.
Runtimes fijados previos: Go1.26.7, Node24.14.1, pnpm11.19.0, Python3.14.4,
PostgreSQL18.6. gate.ps1 -Mode Go|Postgres|Browser -Root rebuilt; Browser -Project
'' ejecuta4. Web candidate compilada con archivos idénticos. Logs/traces/snapshots
y scripts sintéticos fuera de distribución; no datos privados ni efectos live.

| Archivo | SHA256 |
|---|---|
| `microsoft_playwright_browser_gate/tests/enterprise-web.spec.mjs` | `27d3bdb6cdcc89f8ac890b007a4802aec2d88c6b0642637b9657700253e2370a` |
| `src/components/franchise-command-panel.tsx` | `57fd49ea17a65ef6e25dd7bcb9e2daf71da0eed89d6e677e5aa597d174e1d148` |
| `src/app/api/enterprise/franchise/commands/route.ts` | `de6f65ef21bcaa54b4277238b456432514e1fbcaff495e4ff4eb8203ca0c4c5a` |
| `src/app/api/enterprise/franchise/commands/route.test.ts` | `a1119eec284951d072cb0964867431c94160985337ff872049a1c7b77eff4e79` |
| `internal/franchisejourney/service.go` | `1cccce80455eefda5e37938e3cd244361e53100053634c0cbe0827f8849c303e` |
| `internal/platform/httpapi/franchisejourney.go` | `e5d588a32f64247cd67fb27b7f20ff79b5147fda1e8e531917e4dd9aa7d079c6` |
| `internal/platform/httpapi/franchisejourney_test.go` | `c0347179a1da5c9c3290edb851b3cf7bee8b2eacab04b97137edd89420b8f1c2` |
| `internal/platform/postgres/franchisejourney.go` | `52f09a4f9220a21e45cae4288011c929394487cc28fbaafaf89b428dffe697b9` |
| `internal/platform/postgres/franchisejourney_integration_test.go` | `95c23381c507a358e0bfe655cc4d378b92e2f64b0ea231025fbe7bc64e00bbc9` |

| Log | SHA256 |
|---|---|
| `slot-red.log` | `f5b9722427482ce141b5eea49ef51ad4aa18ae9ad13bbe035393977f4295af70` |
| `slot-red-valid-fixture.log` | `9ac0d5fa5b68be42eab165504d03b2ba4afe32bb3f2f532ff13dd93900368070` |
| `slot-r1.log` | `a5b8ceeff48d9157c5258a16135f0bafd4fe0651849722fe6582a2855baed4c9` |
| `slot-pg-r1.log` | `4834f5a87697c88e3e8b57d6fd5e1854b2d4e3abd758e294a321a10ca497103e` |
| `browser-r1.log` | `fa3de700ec0910db3c847245b30d7e98cfdc72c0ee5d0a1919f469a2b0d59f44` |
| `browser-r2.log` | `d412cd6c604ba0cb4a5f9b2fad3a301381c3a004423d959ccb3d4b4aa7e1d7ff` |
| `browser-r3.log` | `78c00e3ff10aa423d34b48e20b28b7a2d58413bb6974b0584c0b70ef7630a8b3` |
| `browser-final.log` | `4ba40edb9546da25d538c91c073423b669a805260dc89eea13766953942e302a` |
| `postgres-final.log` | `3f4cf26ce0686f81d9214fef3e92f1ca5da9db817e735e58faebe8e8d927cc22` |
| `go-test-final.log` | `01a64921a551ef88c6d5834e2cb188d91019eabe88f11911e17f1b7801a9478a` |
| `go-vet-final.log` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `go-build-final.log` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `web-test-r1.log` | `1d71175a88e44ba259273944b28ab8fe1c0fcf995a532169b7ce92cdd2f0727a` |
| `web-test-r2.log` | `059986c428e0dcc8a77e6cf98244424747d47ea2e26c9d16f08a7ef79414f4a2` |
| `web-build-r2.log` | `0af48bfcafa937b31fc2e9c5f05cd98446a2ddd337ede103b1fc94579b5288ba` |
| `agenda-final.log` | `c2efcb508ec5806ae2006fc806f59d432548fb9568650579202f1f77042fd742` |
| `compose-rebuilt.log` | `78b33cfafb0b8a9e5108a93040a49883e3c63918702861d97962b25cd1e02c46` |
| `restart.log` | `f11fe1d7dc84d68a62b69c1ac76b91d5b774fe9603f8e1eece5f172f279d977e` |

## Verificación final de biblioteca

VERIFY_LIBRARY_PASS

Verificador con checkpoint45; cierre documental en checkpoint46 y resume.
SHA256 library-final.log: `578c3ed0aadb55c657dc5542978b12aa1208d3ebc90de7c36ed060b76f84e6bf`.
