# V331 — recuperación de devoluciones y cambio de sesión

Mantenimiento/verificación del flujo existente, 2026-09-08. PASS focal en
Chromium desktop/mobile, Firefox y WebKit. Dos archivos de pruebas AUTHORED:
GO-FRANCHISE-CUSTOMER-JOURNEY-API0.10.19 y PLAYWRIGHT-BROWSER-GATE0.1.26.
No cambia código de producto, permisos, pins, dependencias ni reglas de negocio.

## Oráculo conectado

Se extiende TEST48 con evidencia EVID41, sin sumar otra fila al denominador.
En recepción y decisión de un mismo caso, después del commit con respuesta
perdida, la pestaña conserva una referencia pendiente y cambia su cookie:

| Transición de fixture | Resultado exigido por acción |
|---|---|
| Lector admin:read | UI aún abierta consulta403; POST explícito403; navegación oculta panel |
| Cookie ausente | GET401 y POST401, sin resultado de devolución |
| Sesión selecciona otra organización | GET403/FORBIDDEN y POST403/ORGANIZATION_FORBIDDEN; referencia anterior no se carga |
| Segundo operador handover:manage | GET200 permitido en organización compartida, actor original intacto; no hereda intención pendiente |
| Restaurar operador original | referencia exacta reaparece, se consulta y continúa sin reenviar |

Por navegador:8transiciones,12rechazos,2lecturas autorizadas y2restauraciones.
Total4proyectos:32transiciones,48rechazos (24GET/24POST),8lecturas y8restauraciones.
Además se repiten las24carreras/48POST de UI de V330:24éxitos/24conflictos,
8respuestas perdidas tras commit y4copias opener. Los24POST denegados por API
son adicionales a los48POST observados en las páginas; no se confunden sus
contadores. PostgreSQL conserva12recibos/12decisiones/44requests/24eventos y
24actores comprobados entre4tenants sintéticos separados. Requests siguen requested.

Se observa el getItem real de sessionStorage después de navegar, sin cambiar su
resultado: la clave corresponde al subject/organización actuales y no se lee la
clave pendiente original bajo el otro scope. En el segundo operador se exige
además un formulario habilitado: un botón deshabilitado antes de hidratar no
basta. El resto de V330 verifica el snapshot durable exacto tras recuperar.
Un segundo operador con permiso de gestión puede leer la evidencia común;
no se inventa una restricción de propiedad por autor que el dominio no establece.

## Fronteras

Cookie cifrada y JWT/JWKS sintéticos reales del harness loopback; no login,
logout ni revocación de un IdP externo. Cookie ausente no prueba expiración.
El caso de organización seleccionada usa un token API limitado a store; prueba
rechazo de membership del BFF, no migración real entre tenants. No elimina
información ya entregada al navegador antes del cambio de sesión. No demuestra
seguridad ofensiva, aislamiento multi-tenant completo, efectos downstream,
aceptación humana, deploy/restore/rollback ni cierre de TEST03/TEST08/T2804.
FAIL423 histórico de Firefox permanece abierto para su incidente particular.

## Reconstrucción y ejecución

67packs/746archivos;744sin cambios y2tests reconstruidos idénticos al candidate.
Build Next previo con misma fuente de producto y dependencias ya admitidas;
procedencia/paridad del build heredada de V330/V321, sin nuevo build ni install.
Node24.20.0/Go1.26.7 fijados por SHA256 en runner; PostgreSQL18.6 loopback.
go test ./internal/platform/httpapi -count=1 PASS en rebuilt; sus suites opt-in
se verifican por separado en los cuatro logs conectados. Sin retry automático.

Reproducir en composición nueva con build de idéntica fuente, dependencias
admitidas, TEST_DATABASE_URL sintético, ELITE_WEB_ROOT absoluto y proyecto
de navegador seleccionado en ELITE_QUOTE_PROJECT. Habilitar con valor1:
ELITE_QUOTE_E2E, ELITE_ORDER_E2E, ELITE_DELIVERY_READ_E2E,
ELITE_RETURN_RECOVERY_E2E, ELITE_RETURN_MULTITAB_E2E y ELITE_RETURN_SESSION_E2E.

```powershell
go test ./internal/platform/httpapi -run '^TestQuoteAcceptanceBrowserPostgres$' -count=1 -v -timeout 10m
```

Stage local elite-v331-96f811b1fece48d1970f1521310fd9cb bajo el temporal del sistema.
Antes, logs y artefactos fallidos se conservan. Sólo Markdown canónico e informe
son portables. No ZIP ni promoción productiva.

## Fallo y corrección

FAIL543: primer probe directo GET usaba action en vez de kind=return; recibió400,
que no sirve como PASS de autorización. Contrato leído: wrapper returnCase y
FORBIDDEN para GET de otra organización. Corrección sólo del test; browser-second
PASS y browser-third agrega observación real de lectura/hidratación. Rebuild y
otros3navegadores PASS con el mismo oráculo. No se modificó el backend para
acomodar el fixture. El primer log/trace permanece disponible.

## Porcentaje solicitado

El registro conserva48tests:41passed,5blocked y2planned. Pendientes7/48;
porcentaje exacto175/12%=14,583333…%, redondeado14,58%. Describe filas de
verificación, sin ponderar complejidad, esfuerzo o cobertura. V331 amplía TEST48
y no mejora artificialmente ese cociente añadiendo un test ya aprobado.
Los7pendientes son TEST02,03,05,06,07,08,09. Roadmap:1macrofrente cerrado y10
abiertos con avances parciales; readiness conserva42observaciones. No existe
una ponderación/calibración que permita calcular porcentaje global exacto de
trabajo restante; no se transforma14,58% en estimación del objetivo completo.

## Evidencia SHA256

| Archivo del stage | SHA256 |
|---|---|
| browser-first.log | `9992faaf62eb4a02503c45a35c7c4a58f81dd4a42233acfcec90553865444070` |
| browser-second.log | `68179155b080c3b612236e79b32ab84a98496f64e61a602020bb375132392efa` |
| browser-third.log | `b592d8e869cedcd4688e7188ba305a66e203550f55c9312e602ef8bfbd1de6b1` |
| browser-mobile.log | `a6e1cb83d29f097467e4fe5f283dc2f02a590653cbd58a8acd30a9e04b0a3f56` |
| browser-firefox.log | `e231f2217be4797f8610f209ee3346cc2c3c876433bda529750836f7fc3f89c2` |
| browser-webkit.log | `5b88bc3a8f8429acdd09a2cab42d147c3781ca3177dd3d447fdcae7513156d88` |
| go-http.log | `0d09723ff04c73b6a0c2f4e2152428c34f8ba5eb77b37d48baca5732769ef3f8` |
| compose.log | `e898323913b79ee835d54404265b0c17656ba599a79ff9d4244db96c09ed75f9` |
| rebuild.log | `b1a6f01a264f38d06d689234e7228a1be40c0021ef002e1ef255a90f5f096120` |
| baseline.json | `29232749b93de8af25938ff68d533575b79bb66ae652409431f31199fc04473a` |
| parity.json | `2b25937490285964fc308b4c78555c2decee29d8260567f37512d073070a1afa` |
| run.ps1 | `1c79c26d9ab6fc2f8abfd6e98a1aa53bedff27f495713cb37311a84aa88ae083` |

| Test reconstruido | SHA256 |
|---|---|
| microsoft_playwright_browser_gate/tests/enterprise-web.spec.mjs | `e15ec87989ce7d72b1062e88fd36570ada3f6f028b91b11a12bc6d782a7da966` |
| internal/platform/httpapi/franchisejourney_test.go | `83b304013de77fbc41097704a4e59967c1b7b84a48825facba198c563cc39cef` |

Cierre integrado pendiente tras checkpoint96; no claim global PASS.

## Cierre integrado

Preflight96:151pasos PASS; VERIFY_LIBRARY_PASS161packs/1453archivos/767Markdown/
52perfiles. Disponibilidad general BLOCKED sólo por Docker ausente; los runtime
SKIPPED/BLOCKED de cada pack no se convierten en ejecutados. Plan del contrato
PASS; evidence integral BLOCKED y TEST48/EVID41 permanecen acotados. Checkpoint96
y resume96 PASS con112evidencias. Checkpoint97 conserva96eventos y sólo registra
este cierre, sin cambio de código después de los gates.

TEST06 se inspeccionó: exige release candidato y restauración NEW/EXISTING desde
el artefacto distribuido. Por tanto no se cierra con estas pruebas de sesión ni
con una copia temporal; TEST06/07/08 requieren su propio cierre del artefacto.
No se genera release/ZIP para eludir los gates pendientes.7/48continúa siendo
14,583333…% de filas pendientes, sin porcentaje global de esfuerzo calculable.

| Cierre del stage | SHA256 |
|---|---|
| preflight96.json | `3076ec2f47083b4397758e2c359f685f78a957f1f00d105e04a726687caa441c` |
| preflight96.log | `3eece28694419ebaf8860480e60d26efa5fec5b678148b3f173d0b76ca18559b` |
| contract-plan.log | `2fd65a7ad27c8b968c64f8079abf8610894d6e7732b0b0f514dba962cf93b565` |
| contract-evidence.log | `551b5754a1c2ab245e45e30ba82964b6ed1abafe87680604a6d1cd60b2a890a1` |
| checkpoint96.log | `0043eaf49b5b507380eddebc233f9c0864d8d04fea6fa90e6aeb4ce5de564acf` |
