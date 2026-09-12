# V358 — aislamiento de moderación y lista de espera

2026-09-09. Mantenimiento correctivo/revisión, TEST02/03. Entrada160 validada;
integración161 y cierre162. Stage `C:/Users/NL/AppData/Local/Temp/elite-v358-3ddd304098ba458bb81f3a1218088722`. Continúa V357 en el mismo
turno: tres cores corregidos en total, sin nuevo dominio o política de negocio.

## Antes y después observados

GO-REVIEWS-CORE0.1.0 y GO-WAITLIST-CORE0.1.1 pasaban11tests originales.
Cuatro regresiones públicas fallan:

- Reviews concatena tenant,subject,author con NUL: dos tuplas distintas colisionan
  tanto en la frontera tenant/subject como subject/author. Submit rechaza otra
  reseña como duplicada. Approve/Reject con la tupla inexistente modifica una
  reseña ajena; Approve la incluye en Approved y en el promedio del tenant original.
- Waitlist valida el subject al Join, pero State/Notify/Seat/Cancel aceptan un
  lookup arbitrario. Join("a\0b","c") es válido; lookup("a","b\0c") resuelve
  ese registro ajeno. State lo lee y las tres transiciones lo modifican, quitándolo
  del orden pending. No se demostró colisión entre dos Join válidos ni explotación
  de autenticación en un producto instalado: el defecto está en la resolución local.

Reviews0.1.1 usa clave struct tenant/subject/author; waitlist0.1.2 usa struct
tenant/subject tanto en map como en slice de orden. Igualdad de tuplas conserva
fronteras para todos los bytes y también protege lookups no admisibles al Join.
No se cambian firmas, conjuntos de entrada válida, validación ni errores.
Reviews String adquiere el mutex de escritura; el riesgo de race se reparó por
inspección y se ejercitó concurrentemente, sin ejecución de -race.

## Semántica preservada y assurance

Reviews fuerza pending en Submit; una reseña por autor/sujeto/tenant, rating1..5.
Approve y Reject sólo desde pending; repetición falla. Approved devuelve copias
de valor y su orden no está garantizado. AverageRating sólo usa aprobadas y es0
sin ellas. No compra verificada, moderador autorizado, antiabuso o publicación web.

Waitlist ordena por Join exitoso bajo mutex, no por reloj del cliente. Position
cuenta sólo pending dentro del tenant. Notify es un estado local. Seat y Cancel
admiten pending o notified; no se añade obligación nueva de Notify previo a Seat.
Entradas terminales conservan su identidad y no pueden volver a Join; no política
de reingreso, capacidad, oferta/expiración, reserva ni envío implementada.

| Dimensión | Prueba / límite |
|---|---|
| Corrección |11tests históricos intactos,4red,19tests x3,4fuentes exactas,vet/build |
| Integración | API estable; fuera de perfil, sin publicación/booking/adapter de envío |
| Seguridad/privacidad | aislamiento sintético; caller autoriza IDs/roles, sin SCA integrada |
| Resiliencia | rechazo sin efecto, una transición terminal concurrente; sin restart durable |
| Rendimiento | índices esperados O(1), lecturas agregadas/Position O(N); sin benchmark target |
| Operación | memoria de proceso, sin supervisor/alerta/retención |
| Recovery | fuentes previas preservadas; rollback reabre fallos, sin restore durable |
| Supply chain | stdlib existente, sin dependencia nueva/SBOM/release/firma |
| Usabilidad | límites explícitos; no UX ni aceptación comercial |

## Verificación independiente

Tests ordinarios: reviews9,waitlist10, cada uno3veces. Los11cuerpos originales
permanecen tras gofmt. Cuatro fuentes finales se reconstruyen desde destinos
nuevos y comparan byte por byte.64goroutines por fase prueban unicidad de Submit/
Join y una sola transición terminal (Approve/Reject o Seat/Cancel), con lecturas
simultáneas. Un test confirma que mutar Approved no modifica el store.

Los oráculos de fuzz son listas lineales con comparación directa de campos;
no llaman key/transition/validadores de producción. Reviews compara todas las
reseñas aprobadas y promedio por tenant/sujeto después de cada paso. Waitlist
calcula rango pending desde su lista, y compara State/Position de cada candidato.
Incluyen Unicode/NUL, campos vacíos/espacios, rating fuera de rango, duplicados,
lookups inválidos, tenants ajenos y transiciones repetidas. Cotas120bytes de
operaciones (24pasos reviews/40waitlist), strings64bytes. No exhaustividad.

| Target | Semillas | Presupuesto | Workers | Ejecuciones PASS |
|---|---:|---|---:|---:|
| FuzzReviewModerationModel |12|10s|4|2388166|
| FuzzWaitlistFIFOModel |12|10s|4|1828437|

Total4216603. Go1.26.7 Windows/amd64, Python3.14.4, PowerShell7.6.5 existentes;
GOTOOLCHAIN=local, GOWORK/GOPROXY/GOSUMDB=off,CGO_ENABLED=0,GOMAXPROCS=4,
GOFLAGS=-parallel=4. No nueva adquisición. Limitación24workers V354 permanece.

FAIL631: primer rebuild/tests pasaron pero vet rechazó5String results ignorados
en probes nuevos. Conservados fixed/rebuilt y rebuilt-vet.log. Descarte explícito
documenta intención de ejercitar lectura; no se silencia vet ni cambia aserción.
corrected/qualified reconstruyen otra vez desde canonical; vet/build/fuzz PASS.

## Alcance de admisión

Ambos AUTHORED/LicenseRef-Workspace-Owner, SUPPORTED_REFERENCE,
REBUILD_VERIFIED / CONDITIONED. Suites V197/V213 históricas no reejecutadas.
Compatibilidad antigua reviews/backend no acredita integración; waitlist conserva
retiro de compatibilidad reminders de V355. V293 mantiene gaps de compra verificada/
moderación durable/UX y cola durable/oferta/expiración/notificación.160otros packs
intactos; ninguno de los dos aparece en pack plans; franquicia67/746 sin cambio.

EVID-68 parcial para TEST02/03, sin cambiar48definiciones/requisitos ni estados.
43PASS/5BLOCKED (10,4167% de controles, no esfuerzo), roadmap1completo/10abiertos.
TEST02/03/05/07/09 siguen pendientes. ZIP148 no incluye estos fixes ni V353–357;
no transferir aceptación al payload final. Preflight161 pendiente al integrar.

## Recibos preservados

| Archivo relativo al stage | SHA-256 |
|---|---|
| baseline160.json | `53de04097b599ec6aa6f8f91d234384b5e4399f43f34cf74e82abfc93d0bb090` |
| resume160.log | `6e7e509236c9e848ff098260d33059393fbaeaf79402aef8a5079cf3ab28d1ea` |
| before161/implementation_packs/GO_REVIEWS_CORE.md | `80456619b92b15ee479b57fe760c81d2cb5d5d280be3ae9e7ccf6fa599c16b6e` |
| before161/implementation_packs/GO_WAITLIST_CORE.md | `195e187b35027ece6212e02dbb75f52af2daf2b628e9cf7d0c1cc7f1a5c4ea43` |
| baseline-tests.log | `03b23ea9df5a9817c175b9f37c726dfa639bb398f019bfb2c36944a309ddd9e2` |
| red-tests.log | `7afae05d61126c3edb53dbf374fc0afe5782aa66f9fa9b301573b869ee559909` |
| red-summary.json | `c7bde8a6300d4b27ed645e2cae885191f3c821560ae27153c8e44f872016ff0d` |
| red358.py | `0b5b0a5392935aed46757d2ab22f7ac866a0088161ecd9590d35add1ae6b1422` |
| fix358.py | `26ef6eb106d98d369ddce8c4b01b505c2c4118f93884f65e113f37637066ba6b` |
| recover358.py | `520c040e73eab774e6a08a9d92443201bc50a4ccdf7d2061b348d1bc873cb9d0` |
| green-tests.log | `89607035f0205dedab640a496b7209256226881739b39242546b01c4710366ca` |
| rebuilt-tests.log | `7715b8dce6858d3ea78232bb8094e3402001cc663082bf9cc0f788ab5fc2905f` |
| rebuilt-vet.log | `59f211813fcf7f341eecad77cf814f980453c7f3dad358256592b443b1dda17a` |
| rebuilt/internal/reviews/reviews_test.go | `109f25889f60b973a3731c67b6c515a03385e7b1dbd5f4c4902b3dd7b4cb9d6b` |
| rebuilt/internal/waitlist/waitlist_test.go | `98254f4ae8df46effdce688eb0bf931a04d572367db689c526c943d87db520b9` |
| qualified-tests.log | `ef6813e6614801cd2b076449ed70e8d799e16d877960825a6ebb453b15004bac` |
| qualified-vet.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| qualified-build.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| fuzz-profile.json | `9c5faec72d147e1626a85d567bbc128f9e27ba6981672e6647951880db1b70e4` |
| fuzz-receipt.json | `0f5d49508083b06c1e357bab44271cb48f70538bedf83e005bde07f05aaf6400` |
| fuzz.log | `f6c10e0de753f6273c9c991b1713d6f675443d44f4a596fc9d91eb45be124540` |
| fuzz-self.log | `3adf9799b746a5b46bedadde27179efee394de52585a408828cd4c2ff13df937` |
| qualification161.json | `1e7e262ba73af71849a245b4a7432e662ec5557a23201b7ef6cc344d1b4c597e` |
| ledger161.json | `4641d8ad18d02864d445aa3b9225de8499b71abb75cbad6ffcb9e37836e246a6` |

## Probes originales reviews

````go
package reviews
import("testing";"errors")
func TestDistinctModerationTuplesDoNotCollide(t *testing.T){
 for _,p:=range [][2]Review{
  {{TenantID:"a\x00b",SubjectID:"c",AuthorID:"d",Rating:5},{TenantID:"a",SubjectID:"b\x00c",AuthorID:"d",Rating:3}},
  {{TenantID:"a",SubjectID:"b\x00c",AuthorID:"d",Rating:5},{TenantID:"a",SubjectID:"b",AuthorID:"c\x00d",Rating:3}},
 }{s:=NewStore();if e:=s.Submit(p[0]);e!=nil{t.Fatal(e)};if e:=s.Submit(p[1]);e!=nil{t.Errorf("distinct review rejected: %v",e)}}
}
func TestAliasedModerationCannotPublishOrReject(t *testing.T){
 for _,approve:=range []bool{true,false}{s:=NewStore();if e:=s.Submit(Review{TenantID:"a\x00b",SubjectID:"c",AuthorID:"d",Rating:5});e!=nil{t.Fatal(e)}
 var e error;if approve{e=s.Approve("a","b\x00c","d")}else{e=s.Reject("a","b\x00c","d")};if !errors.Is(e,ErrNotFound){t.Errorf("foreign moderation got %v",e)}
 if len(s.Approved("a\x00b","c"))!=0||s.AverageRating("a\x00b","c")!=0{t.Error("foreign moderation published review")}
 if e:=s.Approve("a\x00b","c","d");e!=nil{t.Errorf("original pending state changed: %v",e)}
 }
}

````

## Probes originales waitlist

````go
package waitlist
import("testing";"errors")
func TestAliasedWaitlistLookupDoesNotReadForeignState(t *testing.T){
 s:=NewStore();if e:=s.Join("a\x00b","c");e!=nil{t.Fatal(e)}
 // This subject cannot Join, but read/transition APIs used to accept its alias.
 if st,ok:=s.State("a","b\x00c");ok||st!=""{t.Fatalf("foreign state %q %v",st,ok)}
}
func TestAliasedWaitlistCannotNotifySeatOrCancel(t *testing.T){
 for _,action:=range []string{"notify","seat","cancel"}{s:=NewStore();if e:=s.Join("a\x00b","c");e!=nil{t.Fatal(e)};var e error
 switch action{case "notify":e=s.Notify("a","b\x00c");case "seat":e=s.Seat("a","b\x00c");case "cancel":e=s.Cancel("a","b\x00c")}
 if !errors.Is(e,ErrNotFound){t.Errorf("%s foreign effect: %v",action,e)}
 if st,ok:=s.State("a\x00b","c");!ok||st!=StatePending{t.Errorf("%s changed foreign state: %q %v",action,st,ok)}
 if pos,ok:=s.Position("a\x00b","c");!ok||pos!=1{t.Errorf("%s removed foreign pending position",action)}
 }
}

````

## Cierre162 — delta completo y continuidad

Preflight161: 154controles ejecutados PASS,162packs/1461archivos/
796Markdown/53perfiles. Docker ausente; checks live omitidos no cuentan como
ejecutados. Reviews0.1.1/waitlist0.1.2:4fuentes reconstruidas,19tests x3,
vet/build,24semillas y4216603ejecuciones fuzz PASS. FAIL629/630 demostrados;
FAIL631 recuperado sin silenciar vet y conservando el intento fallido.

Acumulado de este turno V357–358: tres cores corregidos,6fuentes exactas,
29tests ordinarios x3,36semillas y5926072ejecuciones fuzz PASS. Los16tests
originales preservan aserciones.7grupos de regresión fallaban antes; no se
traslada este progreso a un porcentaje de esfuerzo global o admisión completa.

Los48tests y requisitos mantienen definición/estado:43PASS/5BLOCKED.
EVID-68 sólo agrega prueba parcial TEST02/03. AUTHORED/CONDITIONED y fuera del
perfil67/746; TEST02/03/05/07/09 pendientes,1macro-tarea completa/10abiertas.
Se verificaron los SHA de recibos ya listados en ambos informes antes de agregar
este cierre. Sin release nuevo o aceptación trasladada de ZIP148. Checkpoint162
actualiza hashes y eventos append-only; receipts finales se generan después del
reporte para evitar referencias circulares.

| Recibo de cierre relativo al stage | SHA-256 |
|---|---|
| preflight161.json | `1b7448c0e38e01ea6f92a4df41076a6f89c4ae21dcb3824933470a1104141d95` |
| preflight161.log | `d44352e3a188cc401bce1114d80a172b2bb413201e274b230872802f66147b99` |
| closure161.json | `8ad3de22e80160a543e6bdcbaa2021a32765cf76c6a471fbde08d8df0a542390` |
| resume161.log | `5b88beca5af96c3aff3677e2bf06978be7c7a9f8b60352722658b8236930ec9c` |
| plan161.log | `2fd65a7ad27c8b968c64f8079abf8610894d6e7732b0b0f514dba962cf93b565` |
