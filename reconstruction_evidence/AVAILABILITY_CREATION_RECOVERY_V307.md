# V307 — creación de intervalos con recuperación durable

> **Errata de trazabilidad V317 (2026-09-08), prevalece sobre las menciones históricas de Go1.26.7 debajo.** Se retira la atribución histórica exacta Go1.26.7 no vinculada a un receipt de identidad por ejecución. No se sustituye por1.26.8 sin prueba del artefacto concreto. La ruta current-toolchain de V281 contiene1.26.8; official-toolchain contiene1.26.7. El cuerpo original queda conservado como historia. V317 revalida únicamente la composición actual y el scope allí enumerado con1.26.7 observado, no todos los snapshots previos. Ver [TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md](TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md).

Mantenimiento correctivo de T2802/T2803/T2804 y FAIL457, desde checkpoint40
validado antes del delta. Fecha 2026-09-08 UTC. Sin nueva dependencia, esquema,
módulo, regla de agenda, dato privado ni efecto externo. Perfil 67 packs/745
archivos. No certifica producción ni completa el roadmap integral.

## Fallo reproducido y corrección canónica

El alta original generaba ID nuevo en cada POST. Después de commit201 y respuesta
perdida, el sender genérico afirmaba no aplicado y volvía a habilitar envío.
La lista por rango/limit1000 no identifica una operación y excluye intervalos
lejanos. El rojo conectado conservado creó un intervalo activo antes del aborto.

GO-FRANCHISE-CUSTOMER-JOURNEY-API 0.10.12 incorpora CreateAvailabilityOnce y
GET /v1/franchise/availability/result. Usa el idempotency_record existente,
scope franchise-availability y hash del request, en transacción con intervalo,
restricciones aplicables y outbox. La ruta requiere clave16..128 y
availability:manage; sujeto autenticado original en evento. Se conserva
ReadCommitted y lock organizacional antes del snapshot de agenda. El método
legacy interno queda fuera del claim de idempotencia de la ruta HTTP.

La consulta une clave, tipo, estado completado, código201 y referencia del recibo
con tenant/organización/intervalo. No hay búsqueda heurística por fechas. Replay
devuelve el estado actual; un intervalo cancelado no resucita. HTTP201 creación,
200 replay; el BFF conserva su envoltorio201 existente para respuestas de POST.
La retención expires_at24h es la del owner; no se promete recuperación tras purga.

TYPESCRIPT-FRANCHISE-JOURNEY-PORTALS 0.14.5 conserva sólo requestKey en sessionStorage
scoped por hash tenant/sujeto/organización. Hidratación cerrada, fence sincrónico,
storage fallido impide POST, recibo contrastado y respuesta incierta queda bloqueada.
GET explícito recupera el intervalo sin POST. Sólo tras consulta positiva y sin
POST activo se habilita Preparar otro intervalo, que retira la referencia
confirmada y limpia el formulario, sin crear ni cancelar nada. Ayuda versionada
availability-create-view/1.0.0. No persistir payload, fechas, motivos ni secretos.
No garantía multi-tab ni autorización de reintentos comerciales automáticos.

MICROSOFT-PLAYWRIGHT-BROWSER-GATE 0.1.17 integra la nueva fase en los cuatro
proyectos. Las tres correcciones son AUTHORED; no código atribuido a upstream.
Fuentes/métodos fijados y restricciones anteriores siguen vigentes; no upgrade.

## Verificación observada

| Gate | Resultado y límite |
|---|---|
| Web | 110 PASS, 1 SKIP conectado explícito; Next build PASS |
| Go reconstruido | go test ./..., go vet ./..., go build ./... PASS; integraciones ausentes no contadas como conexión |
| PostgreSQL reconstruido | Cinco tests de persistencia/actor/cancelación/booking/identidad, tres repeticiones PASS |
| Concurrencia de alta | 8 goroutines:1 creación y7 replay;1 intervalo/clave/evento; hash distinto rechaza |
| Atomicidad y lectura | evento duplicado revierte intervalo/clave; identidad/type/code/body inconsistentes rechazan; tenant/org aislados; autor original preservado; nueva pool recupera |
| Navegadores | 72 fases, Chromium desktop/mobile, Firefox desktop, WebKit desktop PASS |
| Nueva fase | Por proyecto3 intervalos/keys/eventos/actores;1 cancelado; pérdida de respuesta, JSON inválido, concurrente/replay, GET503, storage fallido0POST, permisos, clave ausente, actor forjado y fecha180d |
| Reconstrucción | 9/9 archivos byte-idénticos; composición67/745 |
| Reinicio PostgreSQL | snapshot de intervalos/keys/eventos idéntico antes/después; cluster sintético aislado, no PITR ni restore productivo |

FAIL472: primer Chromium r1 agotó60s en getByLabel Tipo exacto porque la etiqueta
envolvía opciones. Snapshot probó combobox habilitado. getByRole combobox con
nombre exacto corrigió la prueba; se conservaron acciones y aserciones.
FAIL473: reincidencia diagnóstica de ruta supuesta/salida excesiva, corregida con
inventario y lecturas acotadas; no se toma salida truncada como evidencia leída.
FAIL457 sigue OPEN para recursos, capacidad, checklist y recepción/disposición;
no repetir alta de intervalos ni cotizaciones como defectos todavía presentes.

## Reproducción y hashes

Staging conservado: `%LOCALAPPDATA%/Temp/elite-v307-667484e0c80441bcad349fdb85cea497`.
Go1.26.7, Node24.14.1, pnpm11.19.0, Python3.14.4, PostgreSQL18.6 fijados previamente.
gate.ps1 -Mode Go|Postgres|Browser -Root rebuilt; Browser -Project '' usa cuatro.
Web corre en candidate; sus archivos son byte-idénticos a reconstrucción.
Baseline rojo, scripts, logs y traces se preservan fuera de la distribución.
Preflight exit0 informó BLOCKED por Docker/.NET no resueltos en ese comando;
no se convirtió en PASS de acceso ni produjo instalación nueva.

| Archivo reconstruido | SHA256 |
|---|---|
| `internal/franchisejourney/availability.go` | `cb7de6b97ba5074c4a275241004ff8988e90c116f39b4d604d3b1185903f7e6b` |
| `internal/platform/httpapi/franchisejourney.go` | `5c3a52816ba1c007dfe0fc0928fed4ed981b6c2140a28190e0859ba64e312778` |
| `internal/platform/httpapi/franchisejourney_test.go` | `f6acd02fd8c75c40894475cfa27b6cc8c458e4fdd9515b0fe4bca78e077e3c25` |
| `internal/platform/postgres/franchisejourney_availability.go` | `8b28584844fb74c6697d04f20d809e06ea121a4748be03f317886df962bd0e50` |
| `internal/platform/postgres/franchisejourney_integration_test.go` | `d9325e226a5f22206a6b9e05ce4ac8e2d7eeeae378c3593ab2a3341ab7611f3a` |
| `src/components/franchise-command-panel.tsx` | `7849d8ebc2325e0ba01c1353748e14099fcbdf6454d6ad7f8193d32eae14f681` |
| `src/app/api/enterprise/franchise/commands/route.ts` | `c1247fd9fe2caef4cd06391311f4290b59c0726f23ca9aeb57ea80119ce0fdfc` |
| `src/app/api/enterprise/franchise/commands/route.test.ts` | `ae4d8dd4bb90ca60a9456ddee81ae7670f7171ef8a6ca5758c2552344eefc953` |
| `microsoft_playwright_browser_gate/tests/enterprise-web.spec.mjs` | `4f02bf53dc6d6e614f59dff4711c1fd6018575563b9a932c8194c909cf171f04` |

| Log | SHA256 |
|---|---|
| `availability-create-red.log` | `04c21b940cb10a7a4d3a98dc4698d629764ad4fd12616c3b95e845be4e8d4356` |
| `browser-r1.log` | `94706d57da2a36f11e15c4b1f97e11f52802b38f47ae05e776b25d0881889216` |
| `browser-r2.log` | `5123961aa9854cbb5f33ac01c8e50b2c5296492c555fe46f679cf64eb5c7c1c8` |
| `browser-final.log` | `41d3b398acaa3a665c782cbe0e24e844055c418b924e8e05074485588f1e5abe` |
| `postgres-final.log` | `eeb7e244309f3c977c12b959e74db2f80530af466a088a3f8ad731405a7bc258` |
| `go-test-final.log` | `674bc1f79a4f6bb3267f6a3eec4b7479b268a47a484ea804d64db316bb9a66b5` |
| `go-vet-final.log` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `go-build-final.log` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `web-test-r1.log` | `5f8359cdcdd612ce99944baa4587f50ecc755781e689f4a8db129c3e28bc7a06` |
| `web-build-r1.log` | `f63eae47b0ba76cc59696e12902a7cc2067985949639dd4e3d9afb2b7c0dbc62` |
| `compose-rebuilt.log` | `ba7d298931d4aff757b7c23fc901b3f2e8b9a89df446ae58d327504081fb873b` |
| `restart.log` | `1518ed002a569ba6ec9321fcf891736c4f3e43678634e40107979080421f1c4c` |
| `preflight.log` | `a8001e500610d6fa7f31980425980574ec2d4432788b501193dc9dc8f9bde912` |

## Verificación de biblioteca

VERIFY_LIBRARY_PASS

Verificador ejecutado con checkpoint41; cierre documental en checkpoint42 y resume.
SHA256 library-final.log: `616c77b977f8bcba7745f59ba0b17f2acc0ec673904d9e5c245e1c0fd8b9ae6a`.
