# V308 — recursos con identidad y recuperación verificable

> **Errata de trazabilidad V317 (2026-09-08), prevalece sobre las menciones históricas de Go1.26.7 debajo.** Se retira la atribución histórica exacta Go1.26.7 no vinculada a un receipt de identidad por ejecución. No se sustituye por1.26.8 sin prueba del artefacto concreto. La ruta current-toolchain de V281 contiene1.26.8; official-toolchain contiene1.26.7. El cuerpo original queda conservado como historia. V317 revalida únicamente la composición actual y el scope allí enumerado con1.26.7 observado, no todos los snapshots previos. Ver [TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md](TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md).

Mantenimiento correctivo T2802/T2803/T2804, FAIL457, desde checkpoint42 validado.
Fecha 2026-09-08 UTC. La pauta V294 continúa después de V307 en el mismo owner.
No módulo, dependencia, migración, regla comercial ni efecto externo nuevo.
Perfil67/745; el roadmap integral sigue abierto.

## Rojo conservado y corrección

Dos POST HTTP autenticados idénticos de service-bay con la misma Idempotency-Key
devolvían201/201, IDs distintos y2 recursos durables. El evento creado omitía
actor_subject. El formulario usaba el sender genérico con falso no aplicado
ante respuesta incierta. resource-red.log conserva la reproducción PostgreSQL.

GO-FRANCHISE-CUSTOMER-JOURNEY-API0.10.13: CreateServiceResourceOnce mantiene
el idempotency_record existente, scope franchise-resource, hash del request,
recurso, habilidades y outbox en una transacción. Conserva la validación de
tipo/identidad/habilidades, unicidad del principal y organización activa.
Actor procede del principal HTTP, nunca del payload. HTTP201 crea,200 replay.
El método legacy interno conserva su contrato; no se atribuye a él la nueva
garantía de actor o idempotencia del endpoint.

GET /v1/franchise/resources/result exige resource:manage y scope tenant/org;
une registro completado, tipo, código201, recibo y recurso exactos. Habilidades
se recuperan ordenadas; devuelve estado vigente, no reactiva un recurso inactivo.
Sin búsqueda heurística por nombre ni nuevo listado como sustituto de identidad.
expires_at24h conserva política del owner; no se promete recuperación tras purga.

TYPESCRIPT-FRANCHISE-JOURNEY-PORTALS0.14.6 propaga la clave y consulta el owner.
ResourceCreation persiste sólo requestKey en namespace de sesión por hash
tenant/sujeto/organización. Hidrata cerrado y conserva fence sincrónico, timeout,
contraste del recibo y bloqueo frente a respuesta incierta. GET recupera sin POST.
Preparar otro recurso exige consulta positiva y ningún POST/GET activo; limpia
la referencia confirmada sin crear ni modificar el recurso anterior. Habilidades
y tipos se validan antes del envío. Ayuda resource-create-view/1.0.0.
No se almacenan nombre, principal, habilidades, payload ni credenciales en marker.
No garantía multi-tab ni readmisión de los otros comandos genéricos por este fix.

MICROSOFT-PLAYWRIGHT-BROWSER-GATE0.1.18 conecta la fase de recursos a las18 previas.
Todo el delta es AUTHORED, con métodos/locks previamente fijados; no código de
Microsoft/Google atribuido ni nueva adquisición de upstream. Condiciones de
seguridad, carga, restore, target y producción permanecen independientes.

## Gates observados

| Gate | Resultado y límite |
|---|---|
| HTTP rojo→verde | mismo request/key pasa de2 recursos a1;201/200, mismo ID y actor resource-manager |
| Web |111 PASS/1 SKIP conectado explícito, Next build PASS |
| Go | test/vet/build ./... PASS en reconstrucción; skips sin DB no se cuentan como conexión |
| PostgreSQL |6 tests de persistencia/actor/agenda/identidad,3 repeticiones PASS |
| Carrera de recursos |8 goroutines→1 creación/7 replay, hash distinto409, aislamiento tenant/org/key |
| Atomicidad | fallo outbox revierte recurso, habilidades y clave;1 recurso/clave/evento y2 habilidades |
| Lectura | recibos inconsistentes rechazados; autor original preservado, recurso inactivo recuperado sin reactivar, nueva pool recupera |
| Navegadores |76 fases/4 proyectos: Chromium desktop/mobile, Firefox desktop, WebKit desktop PASS |
| Nueva fase | Por proyecto3 recursos/keys/eventos/actores y6 habilidades; bahía, empleado y vehículo; respuesta perdida/JSON inválido/concurrencia, storage0POST, GET503, replay/divergencia, actor forjado, clave ausente, permisos y recuperación tras reload |
| Agenda existente |4 proyectos PASS; alta HTTP de técnico, habilidades y citas mantienen sus constraints |
| Reconstrucción |9/9 archivos byte-idénticos, perfil67/745 |
| Reinicio | snapshot de recursos/habilidades/keys/eventos idéntico tras reinicio PostgreSQL sintético; no restore/PITR productivo |

FAIL474: el nuevo mock esperaba principal_subject undefined para bahía; el BFF
normaliza a cadena vacía. Se corrigió la expectativa exacta; no se cambió el
contrato productivo para acomodarla. Web-r1 fallido se conserva.
FAIL457 sigue OPEN para capacidad/checklist/recepción/disposición. No repetir
recursos ni intervalos como fallos aún presentes; no cierre integral ni ZIP.
Ronda D canal/formato de conversaciones preguntada, sin respuesta inventada.
Decisiones V295/comerciales y ARCA diferida se conservan.

## Reproducción

Staging `%LOCALAPPDATA%/Temp/elite-v308-8dd5ad55124c4b05ad4abb4cb63e786a`.
Go1.26.7, Node24.14.1, pnpm11.19.0, PostgreSQL18.6, Python3.14.4 ya fijados.
gate.ps1 -Mode Go|Postgres|Browser -Root rebuilt; Browser -Project '' selecciona4.
Web usa candidate compilado, archivos byte-idénticos al reconstruido. Logs,
scripts, snapshots y traces sintéticos preservados fuera de la distribución.

| Archivo | SHA256 |
|---|---|
| `internal/franchisejourney/service.go` | `f8ecf6d25ebd7e39d5c48f3119647271c12fe3bdb42404fd915ca2017cb35438` |
| `internal/platform/httpapi/franchisejourney.go` | `8348486da2a1f6ff8c0b4abe1074ad234fd42e611e62ac7555fc299e76e5f45c` |
| `internal/platform/httpapi/franchisejourney_test.go` | `e3f669f89bb1af4063476040aec7f04f53486d321e96d6cc8b3d09341d03168e` |
| `internal/platform/postgres/franchisejourney.go` | `a191226a008760b4ad6ec020f769d7b594f0a9249de6b998cf33a2c5e0778260` |
| `internal/platform/postgres/franchisejourney_integration_test.go` | `81ed83f53405c6a1ad6bdaa1846dd6a391db6213f1a246d73d97e61a966b3623` |
| `src/components/franchise-command-panel.tsx` | `f6701ca8207de5ebd03b39e4644e6672dc51c7cdb4e80502b3428416ebdc45bf` |
| `src/app/api/enterprise/franchise/commands/route.ts` | `641c7c575bfbc6533e7bf24815fa0ac1287d5441692f89ef9d84044384b36b1a` |
| `src/app/api/enterprise/franchise/commands/route.test.ts` | `f7009cd814f827997b71f63b53d21873e05ffd1fbdee3dee7a45257377f8c93c` |
| `microsoft_playwright_browser_gate/tests/enterprise-web.spec.mjs` | `f503947178020b17e372946a8407576dd8d2ff94ac22f9eed3d17148accc1299` |

| Log | SHA256 |
|---|---|
| `resource-red.log` | `b988c43776ce830353628b74f186b4fec9a16c7df1ed4ebd56e23b875adce3ac` |
| `resource-r1.log` | `71ed2adc6f310037efe3befadf1d6c33f0e92a2b2156f20b359d55ae9185b28e` |
| `resource-pg-r1.log` | `394614ab9a2e3e45ff901a094ca332e611661d82068adf5f18d7c628ff0ee04d` |
| `browser-r1.log` | `7f6f8117a7b271f76cb4e20b83a98b48f8a181f28313d7b07dc40cdd2727f14b` |
| `browser-final.log` | `3daf3bd8bc546612eb73ba064b21f490bf8db9fd191632649d790c8a4e874c46` |
| `postgres-final.log` | `75e122c1faa6c2ed9ff03b4db6684d62b125013a5b097ba8cf63865394d6ed86` |
| `go-test-final.log` | `9b67c2971bfb00f82dfaff95dab9c03e91302509a6b5fe4acd5fb49e93ad2c79` |
| `go-vet-final.log` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `go-build-final.log` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `web-test-r1.log` | `68c3659916e0289e1e2d3fb665c8beaaf4b8cec106fecd3d6a0e1b4c68880b15` |
| `web-test-r2.log` | `84686bf46784185227585cbd1f1d2845358e0e903a7e603162cfb438df1da6af` |
| `web-build-r2.log` | `7ffb4e972f072f1c20448a470622212bf5f75e4bd8aec89461999a866913363c` |
| `agenda-r1.log` | `2d354380c37962e499183cbbea4b0135c572a313f2429d4a0ee0d8416902c6f2` |
| `compose-rebuilt.log` | `ffc33b47077b9dc7c8e30c805a2263f0dec32193e647cef82af016821f559af7` |
| `restart.log` | `f47c0b665eea179db9d2f60834d4232ddb0d389283d3f88a9daf20267a19df8f` |

## Conciliación de assurance

TEST25–27/EVID16–18 enlazan V306–V308 con REQ03/04 y las dimensiones aplicables.
Plan PASS; nivel evidence retorna1 por los TESTs/benchmark pendientes, sin retirarlos.
Contrato previo SHA256: `a59946cf4db7a2a43dd9aa3b891eaed84928ba2b772c78035f94e51dc80717b6`.
FAIL476: captura con encoding Windows falló tras escribir contrato; se forzó UTF-8
y se repitió sólo validación, sin duplicar entradas. Logs plan/evidence conservados.

## Verificación final de biblioteca

VERIFY_LIBRARY_PASS

Verificador ejecutado con checkpoint43; cierre documental en checkpoint44 y resume.
SHA256 library-final.log: `06bbbfee184a3f02ed577d57da774ae328e8a0c85ca6df8d128a27b522c543c9`.
