# V360 — encuestas y alcance de campañas por tenant

2026-09-09. Mantenimiento correctivo/revisión TEST02/03. Checkpoint164 validado
antes de editar; integración165 y cierre166. Stage `C:/Users/NL/AppData/Local/Temp/elite-v360-f944579a754644c3961c5c12b1f568bb`.

## Defectos demostrados y correcciones

11 tests originales pasaban en surveys/marketing0.1.0. Cinco grupos nuevos fallan:

- Encuestas: tuplas distintas tenant/survey/customer colisionan al concatenarlas
  con NUL y rechazan la segunda respuesta; se prueban ambas fronteras. Responses
  compara campos completos y no se observó fuga de lecturas entre tenants.
- Campañas: ID global impide usar el mismo ID en otro tenant. Due(now) carece de
  argumento de tenant y devuelve campañas de todos. Send(id,recipient) y
  Remaining(id) tampoco permiten expresar el scope del caller.
- Due devuelve una copia superficial cuyo slice Recipients comparte backing
  array. Mutarlo cambia destinatarios internos, permite reconocer el inyectado
  y elimina el original. Create ya copiaba/deduplicaba su input; se conserva.
- La concatenación campaign/recipient hace que un acuse de ("a\0b","c")
  cuente también para ("a","b\0c"), falseando Remaining y bloqueando su acuse.

Surveys0.1.1 cambia el índice por tupla tenant/survey/customer. Marketing0.2.0
usa tuplas tenant/campaign y tenant/campaign/recipient, exige tenant en lecturas
y acuses, y copia Recipients en Due. Ambos String usan el mutex de escritura.
Concurrencia local ejercitada; no se ejecutó -race ni se afirma bypass de auth
en un producto instalado. Las credenciales y efectos externos no se tocaron.

## Migración y contrato conservado

| API antigua | API requerida desde marketing0.2.0 |
|---|---|
| Due(now) | Due(tenant, now) |
| Send(id, recipient) | Send(tenant, id, recipient) |
| Remaining(id) | Remaining(tenant, id) |

No hay shim global. El tenant debe provenir del contexto autorizado del caller.
Un consumer que usa las tres firmas anteriores produce tres errores de compilación
por argumentos insuficientes; el negativo se conserva y no se cuenta como fallo
de producto. Ningún otro implementation pack referencia/importa este core en la
búsqueda registrada. No se admite por esa ausencia un montaje externo desconocido.
Se retiró la compatibilidad declarada con dispatch no conectado.

Send registra sólo un acuse local, no envía SMS/email/push ni valida ScheduledAt.
El duplicado retorna ErrAlreadySent, no replay exitoso. Due es snapshot repetible,
sin reserva; incluye ScheduledAt==now y no garantiza orden. Remaining desconocido
es0, sin distinguirlo de completado. Create conserva Sent dado por el caller;
Sent=true no fabrica acuses por destinatario. Estas reglas existentes se preservan,
no se inventa política de campaña o consentimiento/supresión. Metadata Send O(1)
se corrige a O(R), pues busca destinatario y recuenta; Due copia destinatarios.

Encuestas conserva score0..10, promotores9..10, detractores0..6 y pasivos7..8;
denominador incluye todas las respuestas, vacío devuelve0. Igualdad exacta de IDs,
sin normalización ni restricciones nuevas. No compra, consentimiento, invitación,
representatividad estadística ni retención durable demostrados.

## Pruebas y oráculos

4/4 fuentes reconstruidas desde destinos nuevos, byte-idénticas.20tests ordinarios
(surveys8,marketing12) x3; vet/build PASS. Los5tests de surveys mantienen cuerpos;
los6de marketing sólo migran argumentos tenant y el índice privado del test,
sin rebajar aserciones. Cinco grupos red preservados contra las fuentes antiguas.

64goroutines prueban unicidad de respuestas/campañas por tenant y un acuse por
destinatario/tenant, con lecturas concurrentes. Snapshots de encuestas y campañas
no mutan storage. Se prueban límites de score, horario inclusivo, input copy,
acuse ajeno, repetición, Sent importado y Remaining. Sin prueba -race ni carga.

Fuzz usa modelos de listas, igualdad directa de campos y flags por destinatario;
no emplea keys ni helpers de producción. Encuestas compara el conjunto completo
de Responses y NPS con math/big.Rat (tolerancia1e-12 por división/multiplicación
binary64 existente). Marketing compara Due/Remaining después de cada operación,
valida/rechaza entradas, deduplica con búsquedas lineales y muta cada snapshot
devuelto para detectar alias en la siguiente consulta. No función de envío real.

| Target | Semillas | Presupuesto | Workers | Ejecuciones PASS |
|---|---:|---|---:|---:|
| FuzzSurveyResponseModel |12|10s|4|1491142|
| FuzzCampaignAcknowledgementModel |12|10s|4|2373682|

Total **3864824**, sin exhaustividad. Strings64bytes, hasta24respuestas o25pasos
por caso. Unicode/NUL, blancos, duplicados, horarios, flags y errores sintéticos.
Go1.26.7 Windows/amd64, Python3.14.4 y PowerShell7.6.5 existentes; GOTOOLCHAIN=local,
GOWORK/GOPROXY/GOSUMDB=off,CGO_ENABLED=0,GOMAXPROCS=4,GOFLAGS=-parallel=4.
math/big es stdlib de pruebas, sin dependencia nueva. Límite24workers V354 intacto.

## Implementation assurance y admisión

| Dimensión | Evidencia / condición pendiente |
|---|---|
| Corrección | red/green, oráculos independientes,4fuentes exactas,tests/vet/build |
| Integración | marketing rompe3firmas; sin dispatch/consumer productivo integrado |
| Seguridad/privacidad | tuplas/snapshots sintéticos; sin auth, consentimiento ni SCA integrada |
| Resiliencia | rechazo y acuse local atómico; sin reserva/receipt/restart durable |
| Rendimiento | NPS O(N); Send/Remaining O(R); Due O(C+copias); sin benchmark target |
| Operación | memoria de proceso; sin supervisión, alertas o retención |
| Recovery | before165 preservado; rollback reabre defectos, no restore durable |
| Supply chain | stdlib existente y SHA; sin SBOM/firma/release nuevo |
| Usabilidad | migración documentada; sin UX ni aceptación de negocio |

Ambos AUTHORED/LicenseRef-Workspace-Owner, SUPPORTED_REFERENCE,
REBUILD_VERIFIED/CONDITIONED y fuera de pack plans.160otros packs intactos,
franquicia67/746 sin cambio. V200/V201 históricos no reejecutados. V293 conserva
owners de lead/ads/notificación; segmentación/consentimiento/conversión y encuesta
durable siguen pendientes. Trece cores tienen correcciones acotadas desde V353;
esto no cierra la admisión de los21cores ni habilita efectos.

EVID-70 parcial para TEST02/03.48definiciones y requisitos intactos,43PASS/5BLOCKED,
1macro-tarea completa/10abiertas.10,4167% de controles pendientes no mide esfuerzo.
TEST02/03/05/07/09 abiertos; ZIP148 no recibe aceptación para este payload cambiado.
Preflight165 pendiente. FAIL634: rg recibió wildcard literal de Windows, OS123;
se conservó el error y repitió con -g/directorio, sin usar la salida parcial como
prueba de ausencia. FAIL635/636 tienen fuentes/red/green; no se ocultan fallos.

## Recibos preservados

| Archivo relativo al stage | SHA-256 |
|---|---|
| baseline164.json | `1ca451b99fbde6a9c9d907a059951be310bde8f172b695e4a178e48a9a72611d` |
| resume164.log | `14a2bef4ecca8a0e30bc1c666fefd45050d8db73d1f3e88a8968143ee5193115` |
| before165/implementation_packs/GO_SURVEYS_CORE.md | `e5786433b98612a5ff24b86336d610fe8818ce98dc9d84bb877bfd9a31565b9e` |
| before165/implementation_packs/GO_MARKETING_CORE.md | `a9122f7daeed08463f6ab5035dab94150dede1ab50798c000f682bfc29c350d8` |
| search-path-error.txt | `bb0d0352d0752f56e514817f6aa5ba010534869ccaf439c1e732dc0b8ca91d6a` |
| plan-search.log | `3281800fe3e0d199dec368f2c9770b2d7d1a68550de5772b6dbb327089b38241` |
| baseline-tests.log | `01a816cc6c77001895579b45e76ac0ec246da16dd12e554e81a285336779ab00` |
| red-tests.log | `68953ca73075e1b4c7a7fe22c7f5903cf63563869464bc53f04052143165ba5f` |
| red-summary.json | `6476951b0d0a88eb7ed47ceef453898e7470239b586cd6d2833ac703a26889c2` |
| red360.py | `2973c1be9ee430482f1d07dcbf9253b8d9858e886c3929dfd5f1f178bd401883` |
| fix360.py | `64def32ea79c067c5364757c2afe592ab99d66187c73caf58038b3eace207e7e` |
| green-tests.log | `7c9ceb7a35e457d3e05b5ac5ac038fd55a0aa76d47665136167b9be98672c972` |
| rebuilt-tests.log | `495961f82ffc664680e93b53ece7113100c75f18a4dcf5b9caf98153ce35c003` |
| rebuilt-vet.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| rebuilt-build.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| legacy-callers-rejected.log | `3ce3d966a4657cbb92d67e2a2112def2a64fcd4ac9bde417406499b21ad5f131` |
| legacy/internal/marketing/legacy_usage_test.go | `5d677b98ce6df59bc52b8f910d4c0d47a7d867d2457b197f09d9821117d5b82c` |
| fuzz-profile.json | `65753460771de6a26a6fb1aa9749799add605a5c6fb02505699d95a6bbe795a5` |
| fuzz-receipt.json | `13962029d11ad3572180fcc4c4cd057c3fadde8c93da6731b64d0f94e4054e11` |
| fuzz.log | `9a130e2b3e9fae34d666e7f2981823b50aa20476f1151b8171180ff886dcf733` |
| fuzz-self.log | `86a60bb47b987539ea759b155c5908768dbb33f023d2829bb35c53806075ed93` |
| qualification165.json | `8e4f420203550976ec330d440681b311996c042f28cde04c1bcee6146499bc76` |
| ledger165.json | `3e286446deb21c1111128b028624a7ed79caa1bc019821dfad326153e5c5eed9` |

## Probes originales surveys

````go
package surveys
import("testing")
func TestDistinctSurveyTuplesDoNotCollide(t *testing.T){
 for _,pair:=range [][2]Response{
  {{TenantID:"a\x00b",SurveyID:"c",CustomerID:"d",Score:10},{TenantID:"a",SurveyID:"b\x00c",CustomerID:"d",Score:0}},
  {{TenantID:"a",SurveyID:"b\x00c",CustomerID:"d",Score:10},{TenantID:"a",SurveyID:"b",CustomerID:"c\x00d",Score:0}},
 }{s:=NewStore();if e:=s.Submit(pair[0]);e!=nil{t.Fatal(e)};if e:=s.Submit(pair[1]);e!=nil{t.Errorf("distinct response rejected: %v",e)};if s.NPS(pair[0].TenantID,pair[0].SurveyID)!=100||s.NPS(pair[1].TenantID,pair[1].SurveyID)!=-100{t.Error("independent survey aggregate missing")}}
}

````

## Probes originales marketing

````go
package marketing
import("testing";"errors";"time")
func TestCampaignIDsIndependentPerTenant(t *testing.T){
 s:=NewStore();a,b:=camp(),camp();b.TenantID="other";if e:=s.Create(a);e!=nil{t.Fatal(e)};if e:=s.Create(b);e!=nil{t.Fatalf("another tenant blocked: %v",e)}
}
func TestDueCanRestrictCallerTenant(t *testing.T){
 s:=NewStore();a,b:=camp(),camp();b.TenantID="other";b.ID="OTHER";for _,c:=range []Campaign{a,b}{if e:=s.Create(c);e!=nil{t.Fatal(e)}}
 // The legacy API has no tenant argument with which to enforce this boundary.
 due:=s.Due(time.Now());if len(due)!=1||due[0].TenantID!="t"{t.Fatalf("unscoped due returned %d campaigns",len(due))}
}
func TestDueRecipientsDoNotAliasStore(t *testing.T){
 s:=NewStore();if e:=s.Create(camp());e!=nil{t.Fatal(e)};out:=s.Due(time.Now());out[0].Recipients[0]="injected"
 if e:=s.Send("C1","injected");!errors.Is(e,ErrNotFound){t.Fatalf("snapshot changed real recipients: %v",e)}
 if e:=s.Send("C1","a");e!=nil{t.Fatalf("original recipient removed: %v",e)}
}
func TestDeliveryIdentityBoundary(t *testing.T){
 s:=NewStore();a,b:=camp(),camp();a.ID="a\x00b";a.Recipients=[]string{"c"};b.ID="a";b.Recipients=[]string{"b\x00c"};for _,c:=range []Campaign{a,b}{if e:=s.Create(c);e!=nil{t.Fatal(e)}}
 if e:=s.Send("a\x00b","c");e!=nil{t.Fatal(e)}
 if n:=s.Remaining("a");n!=1{t.Errorf("other campaign falsely acknowledged: %d",n)}
 if e:=s.Send("a","b\x00c");e!=nil{t.Errorf("distinct acknowledgement rejected: %v",e)}
}

````

## Cierre166 — delta y migración verificados

Preflight165: **154 controles ejecutados PASS**,162packs,
1461archivos materializables,798Markdown y53perfiles. Docker ausente; checks
live omitidos no cuentan como ejecutados. Surveys0.1.1/marketing0.2.0 mantienen
4fuentes exactas,20tests x3,vet/build,24semillas y3864824ejecuciones fuzz PASS.
Las11aserciones históricas se conservan con migración de argumentos tenant e
índice privado en marketing. Tres firmas sin tenant fallan compilación.

Se verificaron 23SHA de recibos y los4SHA de las fuentes reconstruidas.
FAIL635/636 tienen regresión demostrada; FAIL634 conserva búsqueda fallida y
corrección. Snapshot defensivo y scope explícito no prueban dispatch, consentimiento,
supresión, receipt de proveedor o entrega exactamente una vez. No efectos externos.

Los48controles y requisitos conservan definición/estado:43PASS/5BLOCKED.
EVID-70 parcial TEST02/03;13cores con fixes acotados no equivalen a admisión
de21cores. Perfil67/746 intacto; sin release nuevo ni aceptación transferida del
ZIP148. Checkpoint166 refresca hashes y conserva eventos append-only; sus
receipts se generan después del informe para evitar referencias circulares.

| Recibo de cierre relativo al stage | SHA-256 |
|---|---|
| preflight165.json | `e5af17b1751edd834193c0a614637f29e160cbee0bbb5cc23781796ceb2c7286` |
| preflight165.log | `bb19973611e02f33c9fb179c49de2cb5fbe413b352e451b4af787cd67edbad08` |
| closure165.json | `62d87bb3f905a050f484419214a978f30408a204211d6a890a6e0ef35d9b367a` |
| resume165.log | `70f1a2b9077ca97820a8fe2e148f5c73b80cfe746c563ccc8c32a64dc1b19c62` |
| plan165.log | `2fd65a7ad27c8b968c64f8079abf8610894d6e7732b0b0f514dba962cf93b565` |
