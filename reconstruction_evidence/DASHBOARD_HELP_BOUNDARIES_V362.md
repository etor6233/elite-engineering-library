# V362 — períodos de dashboard y aislamiento del centro de ayuda

2026-09-09. Mantenimiento correctivo/revisión TEST02/03. Entrada168 validada
con kit ENGINEERING_EXECUTION_VALIDATOR1.3.1 recién reconstruido y todos los
owner hashes del cursor verificados. Integración169/cierre170.
Stage `<LOCALAPPDATA>/Temp/elite-v362-afda5e60566943998f33a6559f376d05`. No producto nuevo ni aprobación READY_TO_BUILD.

## Observación, corrección y compatibilidad

Los8tests originales pasaban (dashboards3,helpcenter5). Cinco grupos añadidos
fallaban sobre fuentes anteriores, con salida no cero conservada en red-tests.log:

1. Build acepta2026-00,2026-13,2026-99 aunque exige período válido.
2. Dos artículos con tuplas tenant/id diferentes colisionan: (a\0b,c) y(a,b\0c).
3. Update con la segunda tupla puede editar el artículo de la primera.
4. Publish/Archive con esa tupla extranjera cambian el estado del otro artículo.
5. Update/Publish/Archive incrementan MaxInt64 a MinInt64 y mutan el artículo.
   Este extremo se inyecta en el map privado sintético; Create siempre inicia1.
   No se afirma haberlo alcanzado por un historial real o restore externo.

GO-DASHBOARDS-CORE0.1.1 primero valida forma y luego rango lexicográfico01-12;
la evaluación cortocircuitada evita slice fuera de rango. Conserva cuatro dígitos
de año, incluso0000 antes permitido: no nueva regla fiscal/calendario empresarial.
Conteos, ratings/NPS y fórmulas siguen iguales. Input es caller-owned; el snapshot
devuelve valores agregados sin retener slices. Un tenant escrito en la salida
no demuestra pertenencia, autorización, moneda o período real de los datos.

GO-HELP-CENTER-CORE0.1.1 usa struct tenant/id, sin delimitador ni cambio de firmas.
String toma el mutex del store. Antes de mutar estado/contenido/fecha, el límite
MaxInt64 devuelve ErrVersionExhausted; MaxInt64-1 todavía incrementa con éxito.
Conserva precedencia lookup→versión→contenido/estado→agotamiento; no reset/saturación.
Update continúa permitido en draft/published/archived, sin republicar al editar.
Published y Search filtran el tenant almacenado y sólo statepublished.
No se cambia orden de resultados ni se promete índice/full-text.

La metadata anterior decía que Search delegaba en GO-SEARCH-CORE; el código
observado usa strings.Contains, sin adapter. Se corrige el claim y compatibilidad
en su owner sin atribuir una integración inexistente. Publicar aquí sólo cambia
un estado en memoria. No se envió ni publicó contenido a un servicio.

## Verificación reproducible

4fuentes canonicales actualizadas mediante update_pack_from_tree.ps1 y después
materializadas a destinos nuevos. Las4son byte-idénticas al árbol calificado.
**19tests ordinarios ×3** (7dashboards/12helpcenter), go vet y go build PASS.
Las8funciones de tests históricas mantienen exactamente el cuerpo normalizado
por gofmt; no aserción retirada ni shim de compatibilidad.
64goroutines disputan una misma versión: un solo Update gana; lecturas Published,
Search y String concurrentes. Sin -race ni afirmación de ausencia exhaustiva de races.
Los tests también fijan shape/calendario, agregados negativos, scores inválidos,
valores máximos representables, copia de lectura y límites/precedencia de errores.

| Target | Semillas | Presupuesto | Workers | Ejecuciones PASS |
|---|---:|---|---:|---:|
| FuzzCalendarAndKPIModel |18|10s|4|2825096|
| FuzzArticleHistoryModel |16|10s|4|1441693|

Total **4266789**. Gate nativo reconstruido con self-test PASS y receipt nuevo;
oráculo KPI usa tabla explícita de meses y math/big.Rat, tolerancia absoluta1e-12,
con hasta256scores por lista. Oráculo help usa map anidado tenant→id, historial
de hasta64operaciones y4identidades que incluyen la colisión. Compara errores,
versiones, estado/contenido, cardinalidad y vistas Published/Search; no llama al
helper key de producción para decidir el resultado. Un op sintético puede fijar
MaxInt64 en el store privado y en el modelo; no representa API/restauración real.

Go1.26.7 Windows, Python3.14.4 y PowerShell7.6.5 existentes; ningún runtime ni
dependencia adquiridos. Go ejecutable SHA256:
5463fe58fa999d74420f00ee1b36d31c3da90a57ff204159f159859567ad61fc.
GOTOOLCHAIN=local,GOWORK/GOPROXY/GOSUMDB=off,CGO_ENABLED=0,GOMAXPROCS=4,
GOFLAGS=-parallel=4; cache recuperada V343. Fuzz finito no prueba exhaustividad
ni es benchmark de throughput/SAST/DAST/carga.

## Assurance y condiciones

| Dimensión | Evidencia / límite |
|---|---|
| Corrección |5grupos red/green, oráculos diferentes,8tests históricos,19tests x3 |
| Integración | firmas conservadas; sin adapter KPI/search/backend admitido |
| Seguridad/privacidad | tuplas sintéticas y negativos cross-tenant; falta auth del actor/SCA integrada |
| Resiliencia | rechazos sin mutación y versión acotada; sólo memoria |
| Rendimiento | KPI O(N), Search scan de artículos/texto; sin budget target probado |
| Operación | sin collector, persistencia, UI o publicación externa |
| Recovery | before169 preservado; rollback reabre fallos, no restore durable |
| Supply chain | stdlib/pin/SHA y rebuilt; sin nuevo SBOM/release/firma |
| Usabilidad | errores explícitos/metadata corregida; sin aceptación humana de UI |

La media de ratings usa acumulador int: el target debe acotar tamaño del corpus
para mantener suma representable. Esta revisión no declara slices ilimitados ni
arregla ese límite por un fuzz de256scores. Help no tiene entrada de snapshot;
no admite versiones internas negativas/corruptas ni datos externos por inferencia.
Autorización, límites de entrada/retención, durable store y buscador siguen pendientes.

AUTHORED/LicenseRef-Workspace-Owner, SUPPORTED_REFERENCE, REBUILD_VERIFIED /
CONDITIONED. Ambos fuera de perfiles;160otros packs byte-intactos respecto a168.
Perfil integral67packs/746archivos sin cambio. V207/V203 son evidencia histórica,
no corridas nuevas.17cores con correcciones acotadas no equivalen a21admitidos.
Quedan i18n/observability/SEO/social-posting en la cobertura de claim del conjunto;
observability conserva su trabajo previo. TEST02/03 parciales,EVID-72.
43/48controles y requisitos no se rebajan;5bloqueados y1macro-tarea completa/10abiertas.
TEST03seguridad,TEST05servicio/collector,TEST07release exacto yTEST09historial
autorizado siguen abiertos. No porcentaje de esfuerzo total ni release aprobado.
Docker no bloquea estas pruebas; no se repitió PostgreSQL ni se instaló Docker.
ZIP148 no hereda esta nueva revisión. Preflight169 pendiente en esta integración.

## Recibos preservados

| Archivo relativo al stage | SHA-256 |
|---|---|
| baseline168.json | `3bb590686351455f7fc43429d983a15e4c3a141cfd0676bb24e2544614c59f22` |
| resume168.log | `bc226841f2fc34fec4c900e1243002463724adef038e470f54ba8eb2d32dff92` |
| plan168.log | `2fd65a7ad27c8b968c64f8079abf8610894d6e7732b0b0f514dba962cf93b565` |
| before169/implementation_packs/GO_DASHBOARDS_CORE.md | `6ef647ef1e3084a717a6f00521663a64c96ffa70caa37cccad94cffd1f418085` |
| before169/implementation_packs/GO_HELP_CENTER_CORE.md | `20fc7f9ab1ea1e41102c7fc16bd1a352b95df2b561194d1ff9abbab343b11216` |
| baseline-tests.log | `3bf1a356113535545d290f72e87be088da2ff703162df3c2ef153714506a17b8` |
| red-tests.log | `cf83425a733a35087114ccbbc871658587c4d97cf21171f76291b3620460b38e` |
| red-summary.json | `ed7f08734865b3b3b8f8878681604ee29d5796d79fe7caf3f7c3153e2c709909` |
| red362.py | `793750c8a3902f252cc66ad04ac6d13efc0dc22b3852b3c5e7d085de15e69d54` |
| fix362.py | `c86e26714fa00618aa28fc20557605a70f1939627deee76822b3c15e2a0198b5` |
| canonical362.py | `501d418687cb97863ba01ea8605234163bc340ca6cfb5b8671bf4c90c80cb356` |
| green-tests.log | `ddedf15dd2ba09176a7f22e57e8ed01ce8f978003a582022d28d221f6be37882` |
| rebuilt-tests.log | `2bca8ff636b23eb7ae7e819932a7ae3b47743bc6dc315fa520068501b1e0e54d` |
| rebuilt-vet.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| rebuilt-build.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| fuzz-profile.json | `ce53dc9f380eeccfd3d16b6e79c85193b9df0edaa7aa86ed02d04ed15a62ecf7` |
| fuzz-receipt.json | `ba64acb229136f5223ce2ea83a4bb6f9f05eb632164821d703649e106affedae` |
| fuzz.log | `b0a024b8564e9561214cce5707520bf62ffccfcda13c0973dbee8f949358f12c` |
| fuzz-self.log | `e1b73662ee7228f8f25b5aace2226323c52677929f62ad2669d34b689c60b9f3` |
| qualification169.json | `27d15d63ad85f8b5a3c39636ff0d5293a9d207fbcce28d5472258cb9b82467df` |
| ledger169.json | `a3f3835eabe476f7adf19094119a568e75e907843f3d10fdecb5addd4fdfa1b3` |

## Probes originales dashboards

````go
package dashboards
import("testing";"errors")
func TestInvalidCalendarMonth(t *testing.T){
 for _,period:=range []string{"2026-00","2026-13","2026-99"}{s,e:=Build("t",period,Input{});if !errors.Is(e,ErrInvalidInput)||s!=(Snapshot{}){t.Errorf("invalid month admitted: %q %+v %v",period,s,e)}}
}

````

## Probes originales helpcenter

````go
package helpcenter
import("testing";"errors";"math")
func TestDistinctArticleTuples(t *testing.T){
 s:=NewStore();a,b:=art(),art();a.TenantID="a\x00b";a.ID="c";b.TenantID="a";b.ID="b\x00c"
 if e:=s.Create(a);e!=nil{t.Fatal(e)};if e:=s.Create(b);e!=nil{t.Fatal("distinct article rejected",e)}
}
func TestForeignArticleUpdateRejected(t *testing.T){
 s:=NewStore();a:=art();a.TenantID="a\x00b";a.ID="c";if e:=s.Create(a);e!=nil{t.Fatal(e)}
 if e:=s.Update("a","b\x00c","foreign","foreign",1);!errors.Is(e,ErrNotFound){t.Fatal("foreign edit admitted",e)}
 if e:=s.Publish(a.TenantID,a.ID,1);e!=nil{t.Fatal("foreign edit consumed version",e)}
 if rows:=s.Published(a.TenantID,"");len(rows)!=1||rows[0].Title!=a.Title{t.Fatal("owner article changed",rows)}
}
func TestForeignArticleTransitionsRejected(t *testing.T){
 for _,archive:=range []bool{false,true}{t.Run(map[bool]string{false:"publish",true:"archive"}[archive],func(t *testing.T){
  s:=NewStore();a:=art();a.TenantID="a\x00b";a.ID="c";if e:=s.Create(a);e!=nil{t.Fatal(e)}
  var e error;if archive{if e=s.Publish(a.TenantID,a.ID,1);e!=nil{t.Fatal(e)};e=s.Archive("a","b\x00c",2)}else{e=s.Publish("a","b\x00c",1)}
  if !errors.Is(e,ErrNotFound){t.Fatal("foreign transition admitted",e)}
  for _,row:=range s.articles{want:=StateDraft;if archive{want=StatePublished};if row.State!=want{t.Fatal("owner state changed")}}
 })}
}
func TestVersionExhaustionRejected(t *testing.T){
 for _,op:=range []string{"update","publish","archive"}{t.Run(op,func(t *testing.T){
  s:=NewStore();if e:=s.Create(art());e!=nil{t.Fatal(e)}
  // Synthetic internal boundary: no public API initializes a MaxInt64 version.
  var before Article;for k,a:=range s.articles{a.Version=math.MaxInt64;if op=="archive"{a.State=StatePublished};s.articles[k]=a;before=a}
  var e error;switch op{case "update":e=s.Update("t","a1","new","body",math.MaxInt64);case "publish":e=s.Publish("t","a1",math.MaxInt64);case "archive":e=s.Archive("t","a1",math.MaxInt64)}
  if e==nil{t.Error("exhausted version accepted")};for _,a:=range s.articles{if a!=before{t.Error("failed operation changed article")}}
 })}
}

````

## Cierre170 — reconstrucción y verificación general

Preflight169: **154 controles ejecutados PASS**,162packs,
1461archivos materializables,800Markdown y53perfiles. El inventario general
conserva Docker no disponible; no impidió estas pruebas y no se alteró el verifier.
Dashboards/helpcenter0.1.1:4fuentes byte-idénticas,19tests x3,8tests históricos
intactos,vet/build,34semillas y4266789ejecuciones fuzz PASS.
Se comprobaron 21SHA de recibos y4SHA de fuentes antes del cierre.

FAIL641/642/643 con regresión demostrada;48definiciones/requisitos/estados intactos:
43PASS/5BLOCKED. EVID-72 agrega evidencia parcial a TEST02/03. Perfil67/746sin cambio;
17cores corregidos estrechamente no son21admitidos. Sin producto READY_TO_BUILD,
auth integrada, adapter de búsqueda, fuente KPI conectada o release final.
Checkpoint170/eventos append-only conserva continuidad sin trasladar aceptación
del ZIP148 a los nuevos bytes. Recibos de estado se generan después del informe.

| Recibo de cierre relativo al stage | SHA-256 |
|---|---|
| preflight169.json | `18fdb221461c41c74623a189582f6941b8076473020c8cd8d48da218347a2aec` |
| preflight169.log | `a16135680a5ff7e37f546db63861471cfddd8a0ba42555dafbacfe57d8586deb` |
| closure169.json | `180d93a1faff5544caffc6228548febe73c18c4ff565b1dd92b67a982af8ccd2` |
| resume169.log | `65c0482545323ed1917e16131c660e66d98c01dd0c5509d5dc6e9e531926765d` |
| plan169.log | `2fd65a7ad27c8b968c64f8079abf8610894d6e7732b0b0f514dba962cf93b565` |
