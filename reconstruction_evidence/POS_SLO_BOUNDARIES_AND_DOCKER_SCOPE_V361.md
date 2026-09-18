# V361 — POS, objetivos SLO y alcance real de Docker

2026-09-09. Mantenimiento correctivo/revisión TEST02/03. Entrada166 validada;
integración167 y cierre168. Stage `<LOCALAPPDATA>/Temp/elite-v361-3b233cd29c2b47eb9e0cd793b1b5ed81`.

## Docker: disponibilidad de herramienta y avance son cosas distintas

El probe actual no encuentra docker en PATH ni el CLI en la ruta estándar de
Docker Desktop. No afirma búsqueda exhaustiva de toda instalación posible.
Get-ToolchainReport de VERIFY_EXECUTABLE_LIBRARY.ps1 usa --version, no prueba
daemon ni contenedor. Preflight etiqueta BLOCKED si cualquier toolchain de su
inventario amplio falta, aun cuando sus pasos ejecutados pasan; ese inventario
incluye rutas de contenedores PostgreSQL y opcionales Business Central.

La ruta PostgreSQL nativa ya está calificada en V351: **11 casos y 11 snapshots
de bases idénticos tras reinicio controlado**. Se verificó ahora el SHA del receipt
postgres-rebuilt.json y sus campos pass/cases/restart/owned_postgres_stopped;
no se reejecutó el servidor ni se presenta ese PASS histórico como una corrida nueva.
SHA: 364369b5da63ea5e4bb072c82c23f38947503ed319cee0f22afc0fe545948ef3.

Por tanto, Docker ausente no explica los cinco controles globales aún abiertos
ni impide continuar con estos cores o esa ruta nativa. No se instaló/adoptó Docker,
se cambiaron criterios ni se declaró equivalente toda prueba de contenedores.
El target aún necesita identidad/configuración/collector/retención/alertas y
aceptación, además de admisión/security/release/histórico de los otros controles.
La mención genérica anterior se acota aquí; el historial no se reescribe.

## Errores reproducidos

Los5tests originales de POS y6de SLO pasaban. El footer POS decía4/4 pero las
fuentes observadas contienen5tests; se preserva el texto histórico y se usa el
conteo real en esta evidencia. Cinco grupos nuevos fallan en las fuentes previas:

- POS multiplica qty/precio y suma sin comprobar overflow.3*MaxInt64 wrappea a
  MaxInt64-2 y dos MaxInt64 más4 wrappean a2, permitiendo Complete indebido.
- Tender=MinInt64 y total=1 produce vuelto positivo por overflow de resta y se
  registra una venta con tender negativo.
- Dos tuplas tenant/id distintas colisionan al concatenarlas con NUL.
- Las Lines guardadas comparten slice con el caller; mutarlo después de Complete
  modifica el registro interno. La regresión inspecciona storage interno: no existe
  getter público y no se inventa una API de consulta.
- NewBudget(NaN) retorna budget con nil error, contradiciendo target en(0,1).

POS0.2.0 valida todas las líneas, detecta producto con price>MaxInt64/qty y suma
con subtotal>MaxInt64-total. Todo operando validado es no negativo y qty>0; esas
comparaciones evitan intermediarios desbordados. Compara tender<total antes de
restar, por lo que toda resta aceptada está entre0yMaxInt64. Copia Lines al guardar,
usa clave struct tenant/id y sincroniza String. No precio/regla fiscal nueva.

SLO0.1.1 rechaza NaN en NewBudget. Inf y límites0/1 ya eran rechazados y continúan
igual. No se modifica Record, thresholds ni semántica de alertas por esta corrección.

## Migración y reglas conservadas

Total(sale) y ChangeDue(sale) pasan de int64 a **(int64,error)**. ErrOverflow indica
producto/suma no representable; ErrInvalidSale indica líneas inválidas/vacías,
y ErrShortTender indica pago inferior al total. No saturación ni sentinel ambiguo.
Un caller con asignación int64 antigua falla compilación en ambas firmas.
No se encontraron otros packs que importen/referencien POS; montajes externos
deben migrar y gestionar el error. Se retira compatibilidad no integrada.

Total puede devolver0para líneas válidas de precio0; Complete conserva exigencia
de total>0. Identidad de venta se valida en Complete, no en los helpers numéricos.
Se validan todas las líneas antes del cálculo para conservar ErrInvalidSale frente
a un overflow anterior. Rechazos no consumen ID; duplicado válido sigue ErrDuplicate,
sin replay exitoso. El caller conserva responsabilidad por no mutar inputs durante
la llamada; se prueba propiedad de copia después del retorno, sin -race.

Budget SLO es contador en memoria, sin timestamps/reset/expiración. MultiWindowAlert
compara budgets dados por el caller, no construye ventanas ni configura alertas.
No validación de thresholds ni admisión de contadores ilimitados: el target debe
acotar duración/número de eventos antes del límite int64. No se declara reparado
overflow de contadores; el modelo aquí ejecuta hasta256eventos por caso.

## Verificación y assurance

4/4fuentes reconstruidas byte a byte,20tests (POS11,SLO9) x3,vet/build PASS.
Los11tests originales preservan aserciones; POS migra las lecturas numéricas por
helpers de test que además comprueban nil error.64goroutines verifican unicidad de
ventas por tenant y SLO cuenta6400eventos/3200errores con lecturas concurrentes.
Dos firmas legacy fallan compilación como negativo esperado, no se rebaja el gate.

El fuzz POS calcula productos/suma/vuelto con math/big.Int, clasifica error antes
de almacenamiento y prueba reutilizar ID tras rechazo. SLO compara disponibilidad,
burn rate y remaining con math/big.Rat desde el target binary64 exacto, tolerancia
relativa1e-12; invalidez se expresa como negación de target>0&&target<1 para incluir
NaN. No usa helpers de producción como oráculo.

| Target | Semillas | Presupuesto | Workers | Ejecuciones PASS |
|---|---:|---|---:|---:|
| FuzzSaleExactArithmetic |14|10s|4|3144808|
| FuzzFiniteBudgetModel |20|10s|4|236275|

Total **3381083**, no exhaustividad ni benchmark de throughput. El contador de
ejecuciones del target SLO quedó en236275desde el primer reporte de3s hasta cierre;
la corrida terminó PASS, sin inferir una causa de esa meseta ni una garantía de
rendimiento. Cotas POS9líneas, SLO256eventos y inputs sintéticos. Go1.26.7 Windows,
Python3.14.4 y PowerShell7.6.5 existentes; sin dependencia nueva. GOTOOLCHAIN=local,
GOWORK/GOPROXY/GOSUMDB=off,CGO_ENABLED=0,GOMAXPROCS=4,GOFLAGS=-parallel=4.

| Dimensión | Evidencia / condición pendiente |
|---|---|
| Corrección | red/green, oráculos exactos,4fuentes,tests/vet/build |
| Integración | POS rompe2firmas; sin checkout/ledger o alertas conectadas |
| Seguridad/privacidad | copia/tuplas sintéticas; sin auth ni SCA integrada |
| Resiliencia | rechazo sin registro; sin persistencia/reinicio de estos cores |
| Rendimiento | POS O(N) por líneas; SLO O(1) Record; sin benchmark target |
| Operación | sin collector/ventanas/retención/alerta/servicio admitidos |
| Recovery | before167 preservado; rollback reabre defectos, no restore durable |
| Supply chain | stdlib y SHA; sin nuevo SBOM/firma/release/runtime |
| Usabilidad | migración explícita; sin UI ni aceptación comercial |

AUTHORED/LicenseRef-Workspace-Owner, SUPPORTED_REFERENCE, REBUILD_VERIFIED /
CONDITIONED.160otros packs intactos; ambos fuera de perfiles. Franquicia67/746
sin cambio. V202/V220 históricos no reejecutados; V293 conserva owners de precio,
pedidos/pagos y operación.15cores con fixes acotados no equivalen a admitir los21.
EVID-71 parcial TEST02/03,43PASS/5BLOCKED con48definiciones/requisitos sin cambio;
1macro-tarea completa/10abiertas. Sin porcentaje de esfuerzo ni release final.
ZIP148 no hereda aceptación al payload cambiado. Preflight167 pendiente.

FAIL637 conserva recurrencia de wildcard literal rg; FAIL640 conserva newline
en pointer histórico que provocó Errno22 antes de leer el receipt. Ambos corregidos
sin inferir pérdida/corrupción de evidencia. FAIL638/639 tienen regresión demostrada.

## Recibos preservados

| Archivo relativo al stage | SHA-256 |
|---|---|
| baseline166.json | `65972fd5cbc894ee6319e919a6a71d8e6673aa6bbaed5b6e344acecadaa3805e` |
| resume166.log | `2c6fcb86449d42594bc77f3a0ed820c0057b135607350216d017b60f77f53780` |
| before167/implementation_packs/GO_POS_CORE.md | `1b470a7c7a44389a34b3f4879c3492249e79ca3c681845d8175a177bb486f017` |
| before167/implementation_packs/GO_SLO_CORE.md | `a1ae0b2d1dd61981a137b93c8ac6ef843d5813e9ace7a27f915dd58c3cfe8a9d` |
| search-path-error.txt | `6b9e5608abaa645e74d2036ded0993fd010bafb48292a45bd49b7a56684a67da` |
| pointer-path-error.txt | `4a211bc1f57b646206356377123822c834207e4006cae029ad024cb919c67b8b` |
| docker-probe.json | `c878b1423f108a52f7f49df917c22f62aaf3cb21a41d8ab5135a076227cb4eec` |
| docker-scope.json | `4f5009cb5e2867f5f66898509dc869f33ef35598bd927a761ec1bdc2d963ce90` |
| baseline-tests.log | `c6cc6a511ae430d2b8a0c2b63c10733dc042c5587ddd7b9ad1967c4f007e7755` |
| red-tests.log | `a35362b442ea459688e4f00ac9bc8d39058f2e445c26c4d4007b9a3f96b7bc11` |
| red-summary.json | `59af46c2bc196a72fefb460daf47e832043178498ce1a7bcc7934824b501bca7` |
| red361.py | `e31f39ac3e4389d9f5b72462b641b47800b2e4a5789f1754170297a2586546d3` |
| fix361.py | `445c8050c64a0b5e36a73f559695520598f6b8df774ffce293cbbca662dd6fe5` |
| green-tests.log | `863e6e0b07138f3720c0d2acd65f8a93886e51195982653c72555e8ac5a63aac` |
| rebuilt-tests.log | `f879d296ece8c217af84c233f902e9dde1d37059b5bacaa59577a8202b5297dd` |
| rebuilt-vet.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| rebuilt-build.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| legacy-callers-rejected.log | `551fb54af0539e4417dc067dbe4156ec8a28fe9db540347dfaec3d0681c04e41` |
| fuzz-profile.json | `3f2ff7618b5439da86ffdbf6db769356f37d1d5f39d5f85cf698be409969cf3d` |
| fuzz-receipt.json | `00bf492629efb170a3178ccf40e805beadb065542ca5db323a2d5c5401cdfacb` |
| fuzz.log | `00a7db207e08aeeb630961d4d6c527484fc37aedeac7619c1cf83f1f75aff8c7` |
| fuzz-self.log | `c4126070ceedb74dd3036920fe24444ff3b526132c416945e1adaea8bbfa3f90` |
| qualification167.json | `0f2f3c0da1d9b5e3b4506c2bc6097d5116bc60f490bbb74345fe38b81deb5710` |
| ledger167.json | `fdeaedf99ccc1868e4f368401a90a27044508419781aa9419590ae3a1291f8d9` |

## Probes originales pos

````go
package pos
import("testing";"math")
func TestOverflowingSaleRejectedBeforeRecord(t *testing.T){
 for _,lines:=range [][]LineItem{
  {{ProductID:"p",Quantity:3,UnitPriceMinorUnits:math.MaxInt64}},
  {{ProductID:"p",Quantity:1,UnitPriceMinorUnits:math.MaxInt64},{ProductID:"q",Quantity:1,UnitPriceMinorUnits:math.MaxInt64},{ProductID:"r",Quantity:1,UnitPriceMinorUnits:4}},
 }{s:=NewStore();r:=sale();r.Lines=lines;r.TenderMinorUnits=math.MaxInt64;if e:=s.Complete(r);e==nil{t.Error("unrepresentable total accepted")};if len(s.sales)!=0{t.Error("overflow stored a sale")}}
}
func TestNegativeTenderCannotWrapToPositiveChange(t *testing.T){
 s:=NewStore();r:=sale();r.Lines=[]LineItem{{ProductID:"p",Quantity:1,UnitPriceMinorUnits:1}};r.TenderMinorUnits=math.MinInt64
 if e:=s.Complete(r);e==nil{t.Error("negative tender accepted via wrap")};if len(s.sales)!=0{t.Error("invalid tender stored")}
}
func TestDistinctSaleIdentityTuples(t *testing.T){
 s:=NewStore();a,b:=sale(),sale();a.TenantID="a\x00b";a.ID="c";b.TenantID="a";b.ID="b\x00c";if e:=s.Complete(a);e!=nil{t.Fatal(e)};if e:=s.Complete(b);e!=nil{t.Fatal("distinct sale rejected",e)}
}
func TestCompletedSaleOwnsItsLines(t *testing.T){
 s:=NewStore();r:=sale();if e:=s.Complete(r);e!=nil{t.Fatal(e)};r.Lines[0].UnitPriceMinorUnits=9999
 // Internal stored-value inspection; this API does not expose a sale getter.
 for _,stored:=range s.sales{if stored.Lines[0].UnitPriceMinorUnits!=500{t.Fatal("caller mutated completed record")}}
}

````

## Probes originales slo

````go
package slo
import("testing";"math";"errors")
func TestNaNTargetRejected(t *testing.T){
 b,e:=NewBudget(math.NaN());if b!=nil||!errors.Is(e,ErrInvalidTarget){t.Fatalf("invalid NaN target admitted: budget=%v err=%v",b,e)}
}

````

## Cierre168 — delta numérico y alcance de herramientas

Preflight167: **154 controles ejecutados PASS**,162packs,
1461archivos materializables,799Markdown y53perfiles. Docker sigue no disponible
para rutas de contenedores; no bloquea estas pruebas ni la ruta nativa V351.
El verifier y sus criterios se conservaron byte-idénticos. No se transforma
disponibilidad de herramienta en admisión de servicio, daemon o target.

POS0.2.0/SLO0.1.1 mantienen4fuentes exactas,20tests x3,vet/build,34semillas y
3381083ejecuciones fuzz PASS.11aserciones originales conservadas, con comprobación
adicional de error en helpers POS. Dos firmas legacy de resultado único fallan
compilación. Se verificaron 24SHA de recibos y4SHA de fuentes, además
del receipt histórico PostgreSQL11casos/11bases con reinicio controlado.

FAIL638/639 con regresión demostrada; FAIL637/640 conservan errores de búsqueda/
pointer y recuperación.48controles/requisitos mantienen estado y definición:
43PASS/5BLOCKED. EVID-71 sólo agrega prueba parcial a TEST02/03.15cores con
correcciones acotadas no equivalen a admitir los21; perfil67/746 y límites de
negocio/operación intactos. Sin release final ni aceptación transferida de ZIP148.
Checkpoint168 conserva eventos append-only; receipts se generan después del
informe para evitar dependencia circular.

| Recibo de cierre relativo al stage | SHA-256 |
|---|---|
| preflight167.json | `5a21b1d81d33926f43d76b744acc7f948a5ac60e02406b0388ec66ddbffa106d` |
| preflight167.log | `dc94cc07971b73884c94a4398f90907d869179d1d751d9fd8849b935e252716c` |
| closure167.json | `c6b7467596befa60d1fa84a93a3b49afdf2ccd590d0aecad20b51adacd7a4d1e` |
| resume167.log | `d413de4b95c9562f04205278aee11c42ceda6c86f332604ed9e563173c43e806` |
| plan167.log | `2fd65a7ad27c8b968c64f8079abf8610894d6e7732b0b0f514dba962cf93b565` |
