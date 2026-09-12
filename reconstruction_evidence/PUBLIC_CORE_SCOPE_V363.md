# V363 — idiomas, SEO y agenda social; conciliación de los21cores

2026-09-09. Mantenimiento correctivo/revisión; entrada170 validada mediante
execution kit1.3.1 recién reconstruido y todos los hashes del cursor comprobados.
Integración171/cierre172. Stage `C:/Users/NL/AppData/Local/Temp/elite-v363-d569bfbd627d436085851aa034ab1c90`.
Sin nueva plataforma, producto READY_TO_BUILD, gasto, mensajes o publicación.

## Defectos demostrados y cambios

15tests originales PASS (i18n4,SEO6,social5),8grupos red sobre fuentes previas:
plural vacío reemplaza valores; URLs de otro tenant filtradas por prefijo;
CR/LF/NUL inyectados en robots; LocalBusiness admite URL relativa/javascript;
sitemap admite hostname vacío; ID social global rechaza otro tenant; Due mezcla
tenants y las tres firmas sociales no permiten expresar ámbito de tenant.
Los logs previos se conservan; no se usan mocks como prueba de publicación real.

i18n0.1.1 rechaza formas vacías/whitespace antes de mutar y sincroniza String.
Mantiene firmas, fallback exacto→base→default→key y n==1one/elseother, sin interpolar.
SetPlural conserva ErrInvalidKey para key/locale inválidos y devuelve
ErrInvalidLocale para forma vacía, consistente con el texto vacío de Set.
NewCatalog no valida defaultLocale; el caller debe configurarlo correctamente.
Gramática de locale local no prueba existencia ISO/BCP47 ni soporte lingüístico.

SEO0.1.1 usa tupla tenant/URL y comparación exacta de tenant, sin prefijos.
Mantiene deduplicación literal/orden de inserción y devuelve copias. Sitemap y
LocalBusiness exigen HTTP(S)/hostname no vacío; no DNS/fetch/propiedad de dominio.
Robots rechaza CR/LF/NUL sin mutación. No implementa parser RFC completo, XML,
canonicalización URL semántica, autorización ni certificación de resultados Google.
JSON-LD conserva encoding/json y escape HTML; roundtrip con texto </script> probado.

Social0.2.0 exige **Due(tenant), MarkPosted(tenant,id), MarkFailed(tenant,id,reason)**.
Identidad almacenada por struct tenant/id; duplicado sólo dentro de ese ámbito.
String sincronizado. Tres llamadas legacy fallan compilación por argumento faltante;
no shim global. Due no reclama/reserva trabajo. Marcar posted/failed sólo cambia
estado local, no invoca adapter ni demuestra delivery. Conserva transiciones antes
del vencimiento, reason vacío y validación de schedule futuro; no agrega política
editorial/de envío. Reloj inyectado sólo en tests privados, no API pública de clock.

## Atribuciones corregidas con fuente oficial

Los límites locales60/160 no se presentan más como requerimientos de Google:
su guía indica que la longitud del título/meta-description no tiene ese máximo
fijo y que la presentación puede truncarse según ancho. Se conservan las reglas
existentes como política AUTHORED de code points tras TrimSpace, no graphemes.
Fuentes leídas2026-09-09: [títulos](https://developers.google.com/search/docs/appearance/title-link)
y [descripciones](https://developers.google.com/search/docs/appearance/snippet).

CLDR define reglas/categorías según idioma; n==1 universal no implementa ese motor.
Se retira equivalencia de metadata/comentarios conservando API y comportamiento
binario local: [Unicode CLDR plural rules](https://cldr.unicode.org/index/cldr-spec/plural-rules).
No se adquirió/copió código o datos CLDR, no nueva licencia ni motor adoptado.
authority-review.json es paráfrasis del análisis con URLs/fecha; su SHA identifica
ese registro local, no un hash inventado de bytes upstream descargados.

## Evidencia técnica

**6fuentes actualizadas**, reconstruidas byte a byte desde Markdown a destinos
nuevos; **32tests ×3**, vet/build PASS. Las15funciones históricas preservan sus
aserciones; sólo se agregan argumentos tenant a las llamadas sociales.
Tests concurrentes:48writers/lectores del catálogo;48adds de sitemap sobre8URLs;
48tenants con mismo post ID y transición independiente. Sin -race.

| Target | Semillas | Presupuesto | Workers | Ejecuciones PASS |
|---|---:|---|---:|---:|
| FuzzBinaryPluralFallback |30|10s|4|909489|
| FuzzTenantURLsAndRobotLines |16|10s|4|1473199|
| FuzzScopedScheduleHistory |16|10s|4|2136083|

Total **4518771**,62semillas. Fuzz gate reconstruido/self-test/receipt nuevo PASS.
Oráculos: fallback esperado por decisiones explícitas; map tenant→lista URL con
dedup lineal y salida robots exacta; map tenant→id→post con historia/transiciones.
No usan el helper de clave de producción como oráculo. Límites: labels64bytes,
texto128bytes,64operaciones social/SEO64adds; no carga o exhaustividad.

Observability0.3.1 se reconstruyó aparte, **2fuentes sin cambios,28tests/vet/build
PASS**. Se preserva CANDIDATE y sus reparaciones V291/V292/V301/V342; no se compone
ni se declara nuevo exporter, privacidad universal, servicio o retención.
No se ejecutó fuzz extendido de observability en esta corrida.

Go1.26.7/Python3.14.4/PowerShell7.6.5 existentes, sin dependencia nueva. Go SHA:
5463fe58fa999d74420f00ee1b36d31c3da90a57ff204159f159859567ad61fc.
GOTOOLCHAIN=local,GOWORK/GOPROXY/GOSUMDB=off,CGO_ENABLED=0,GOMAXPROCS=4,
GOFLAGS=-parallel=4. Cache recuperada V343. Sin SAST/DAST/SCA integrada nueva.

| Assurance | Prueba y condición |
|---|---|
| Corrección |8red,15tests preservados,32tests x3,3modelos fuzz |
| Integración | social rompe3firmas; no provider/CLDR/frontend conectados |
| Seguridad/privacidad | separación sintética y rechazo de líneas; falta auth/SCA del target |
| Resiliencia | rechazo sin mutación; datos sólo en memoria |
| Rendimiento | scans acotados en pruebas, sin budget/benchmark productivo |
| Operación | no collector, canal real, entrega o indexación demostrados |
| Recovery | before171 y recibos; rollback reabre defectos, no restore durable |
| Supply chain | pin/SHA y stdlib; no nuevo SBOM/firma/release |
| Usabilidad | migración explícita/atribuciones corregidas; sin aceptación de UI |

## Estado conciliado de los21cores y bloqueos

Con estas correcciones,20cores tienen reparaciones acotadas recientes y el core
de observabilidad conserva reparaciones previas y estadoCANDIDATE. Esto termina
la lista de tres cores aún sin esa revisión reciente, **no la admisión integral**.
El censo siguiente fija los21owners/versiones/claims actuales. Las suites de los
otros17no se reejecutaron aquí; se conserva su evidencia previa.159otros packs
byte-intactos; perfil integral67/746sin cambios. Ninguno de los tres se selecciona.

V293 distingue12solapamientos parciales,8sin equivalente demostrado y1candidato;
esta reparación no convierte esas disposiciones en equivalencia con owners durables.
No duplicar dominios ni incorporar estos stores para completar una casilla.
TEST02 exige la comparación/condiciones/equivalencia aplicable y permanece abierto;
TEST03 requiere grafo exacto e integración de seguridad, no sólo tests de claves.
TEST05 necesita servicio/identidad/configuración/collector/retención/alertas reales;
V351 PostgreSQL nativo y estos28tests no son esa aceptación. TEST07 requiere payload
final/builds/SCA/SBOM/firma más gates materiales; ZIP148 no hereda nuevos bytes.
TEST09 aún no dispone de corpus histórico autorizado/contrato del canal; no se
reemplaza por mensajes sintéticos ni se activan efectos de backfill.

43PASS/5BLOCKED;48definiciones/requisitos conservados,EVID-73 parcial TEST02/03.
1macro-tarea completa/10abiertas. Sin porcentaje global de esfuerzo inventado.
Próximo trabajo: evaluar admisión/equivalencia de esta matriz contra owners
seleccionados y gates conectados ya existentes; no reiniciar la cuenta de fixes.
Docker no impidió esta ruta y no se instaló ni se promovió como requisito universal.
Preflight171 pendiente en integración; cierre se agrega al final.

| Core | Revisión | Admisión vigente | Owner SHA-256 |
|---|---|---|---|
| GO-DASHBOARDS-CORE | 0.1.1 | CONDITIONED | `83291a92995ea0ea6bb27624c83557ec11605fa889476b6ca33c2fec752950e2` |
| GO-FX-CORE | 0.1.1 | CONDITIONED | `b4ba89a8d4072330a10767285bebef0b5ac95c7723100deb4047064c6c434b36` |
| GO-GIFT-CARDS-CORE | 0.1.1 | CONDITIONED | `61f5ca1c91e69884d034009de6e26614f86a6c0267ca874036cae8a5b46662ae` |
| GO-HELP-CENTER-CORE | 0.1.1 | CONDITIONED | `ab297d100f2c5d6973321ade442e731f247e4ceafcfaedb0e78c84341dcbb2ed` |
| GO-I18N-CORE | 0.1.1 | CONDITIONED | `5a9bbcc0e4bf1704a02e4a35a26ac01cf9fb02af40a6ccfcefcc66cfd4c36d9f` |
| GO-LOYALTY-CORE | 0.1.1 | CONDITIONED | `69f889e55ecd7e8bf474934f466dead005b542eb205c72042eb8332fffd1a04f` |
| GO-MARKETING-CORE | 0.2.0 | CONDITIONED | `24cc8a426971aee612d3b9ae2b423302985c5c44073a74365dc63c996181e7f0` |
| GO-OBSERVABILITY-CORE | 0.3.1 | CANDIDATE | `bc24509c8295121f42f614ece09e22b153f7ff3bc9c654ee608d4920f6638a4c` |
| GO-ONBOARDING-CORE | 0.1.1 | CONDITIONED | `3dc86976cc05c3d633ab421742c83532c8191953b9323b324943a885145c1ff7` |
| GO-PAYROLL-CORE | 0.1.1 | CONDITIONED | `c8ff4da4a669af30027956b47073c6ff291e10766dfc36d5678ba19831afde54` |
| GO-POS-CORE | 0.2.0 | CONDITIONED | `23d0137864abd2ae04433e01f61c1303dd16cd5f90e8208c2c793ef64bb518e1` |
| GO-PROMOTIONS-CORE | 0.1.1 | CONDITIONED | `d7c8f92191d041557209023139df6964bffbc4c1787b4deb4eb7dcb58388cd74` |
| GO-REFERRALS-CORE | 0.1.1 | CONDITIONED | `cd11aacdcd7c46234849d77d7e5149bbee9618dd1b0696a83cdaf06f9a2870b6` |
| GO-REMINDERS-CORE | 0.2.0 | CONDITIONED | `16532c41087dc151b5c8b7fbf2bf3d2a276b9d6a18a8c5653aa1cdf5d104c363` |
| GO-REVIEWS-CORE | 0.1.1 | CONDITIONED | `153cfc34f94696e58aecf845570b86f53678bee0fef8e069b8f850d00f8c7531` |
| GO-SEO-CORE | 0.1.1 | CONDITIONED | `b4d240929bde659d3c0abdc7964dd1dd1b2e4f1e6bc66ddc8920cb702fa57b34` |
| GO-SLO-CORE | 0.1.1 | CONDITIONED | `bda35a02080bc531b8fee44553564995f80dcd6ca490c49a271a82ebea1561d0` |
| GO-SOCIAL-POSTING-CORE | 0.2.0 | CONDITIONED | `10beb7a78cc47e917569c2b77eb2d24185b319639721735f76b1fbe22ddcdcb1` |
| GO-SURVEYS-CORE | 0.1.1 | CONDITIONED | `3f450c7a8336d7962a7bb0ce329624a24a889a71803a876b9498b3c404d84145` |
| GO-WAITLIST-CORE | 0.1.2 | CONDITIONED | `2a423e6c7a83589f8d50f62528dc099e53c87e52d2b9eae36f149e95d96a6abf` |
| GO-WARRANTY-CLAIMS-CORE | 0.1.1 | CONDITIONED | `2a93d2b1c2e68b924ee5e50d087c7e61f9d355a06d04c49e3936d2949e42dc31` |

## Recibos preservados

FAIL645 conserva error del contador supuesto14 frente a15tests observados;
FAIL649 conserva nombre V293 adivinado antes de descubrir la ruta exacta.
Ambos recuperados sin alterar tests ni inferir ausencia de evidencia.
FAIL650 conserva parser de IDs que excluía I18N por sus dígitos; censo reanudado
con gramática alfanumérica sin repetir mutaciones anteriores ni checkpoints.

| Archivo relativo al stage | SHA-256 |
|---|---|
| baseline170.json | `c768e7ea287207e8a9d581ebb6c6847be208b630f29443307aaa2c142c10f0ae` |
| resume170.log | `41252903b452debfcc181647dc28bf99f0add48bea11347e1acd368a140c84f6` |
| plan170.log | `2fd65a7ad27c8b968c64f8079abf8610894d6e7732b0b0f514dba962cf93b565` |
| authority-review.json | `f4ccc1fc0ea9e55a8d94f542c26424833370c6eb80dea6201e56f3b1f150b78b` |
| census-error.txt | `63c8f67b0a39180ca344a31c5300a908062250ddcb61b944aee762b9875ea66b` |
| continuation-syntax-error.txt | `099ddd317c3403ef6d13d9cff22cddbe1f11a916620253a159e2c22cb5fe2b0e` |
| continue171-syntax-error.py.txt | `caa1c6d34e9367a2a68ea2b5fc3b5e7c966adc807ddac239527c4fa3cc2afec3` |
| integrate171.py | `68989a75f95ef248d5bf1c89db40bcfc26dc98092be2fbf52b554bba3c9d2ceb` |
| count-error.txt | `0ce5b1fc9a50048ef7aa86888648be5dc9c3add7053eaf8ac9e9a36a810be42b` |
| search-error.txt | `a8708f9d7f1a328952d9e44a3ebadae98a2891ee79b20c1c5790a5d5f69b5f12` |
| baseline-tests.log | `d043fad8e2ee162063df6f9eee85cce5ab2205e935b99817dff91855ce63bc02` |
| red-tests.log | `490f929f62d84aae1e471853da563cb5aeb9dbb3149a41b6ab2e1e8b01bfd94c` |
| red-summary.json | `c45412893ea82f6ac21b868550bee31d4e51e59b1743e6bfc2fc8f53af30c181` |
| red363.py | `eb8563a8d434ef5eadbba8228080ea0f2be81b6fb7cf997c064ad77f59f16ab6` |
| recover_red363.py | `fdffcdce69c251711814468641e19a266d5596673730ffc9e69ba0b5af013142` |
| fix363.py | `fc5a8b5ce9b674350a009124bd08d94ab3d07d66b6e309fc3349df483556a46f` |
| canonical363.py | `b8b4911e6a393408ffbed83dbddfcb66c55cf39797fbff26147dd4e957de2eae` |
| green-tests.log | `a81589630af8d067855d0c801a71a91bbe44691dbd0a270ea6183cbe608fb087` |
| rebuilt-tests.log | `9bd0a2a8fe2fa6316fd895dab7954163a7c5840f0914ece896c90f3681f8e425` |
| rebuilt-vet.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| rebuilt-build.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| legacy-calls-rejected.log | `aec5c6e4bd1bb7c33d30e18ce7dc153564d5e8c968c3b9e15f4d0a41b2286a66` |
| fuzz-profile.json | `5f56c381c0e88eab31ad20561ced50c1e2cba5463f5b4275021246561e757f3d` |
| fuzz-receipt.json | `460056f99083de318844500fc9fc34d214c608f4cb14cb392e99c921686bf25b` |
| fuzz.log | `02a623ba7b4040fd27483c5cd3e1affef36991111e641e8ad5a7d5f49a369f14` |
| fuzz-self.log | `759f402f9210a517a16cd7fd8592902ebe73cb975f4ab066ffd1d90d89c4d861` |
| observability-tests.log | `6f2ee1a82c0c45a5511aae1a4aba12f450080a4934cea00dbe10859a07600c91` |
| observability-vet.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| observability-build.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| qualification171.json | `f122722c5871c5ef4dcd4629e73b5e3331a9e2498f49742cf7e33db6b8a91bcc` |
| core-census171.json | `635c6f7b8e8110449decf00f5a14a6403aded10ea3faea0051c21957f0d9dbfd` |
| ledger171.json | `507e427d6623b60624cfed12461ce5368b33fcbce92e4ed78170715e043164ca` |

## Probes originales i18n

````go
package i18n
import("testing")
func TestPluralRejectsBlankWithoutReplacing(t *testing.T){
 for _,forms:=range [][2]string{{"","many"},{"one",""},{" \t","many"},{"one","\n"}}{c:=NewCatalog("es");if e:=c.SetPlural("items","es","one","many");e!=nil{t.Fatal(e)};if e:=c.SetPlural("items","es",forms[0],forms[1]);e==nil{t.Error("blank form admitted")};if c.Tn("items","es",1)!="one"||c.Tn("items","es",2)!="many"{t.Error("rejected plural changed catalog")}}
}

````

## Probes originales seo

````go
package seo
import("testing";"errors")
func TestSitemapExactTenant(t *testing.T){
 s:=NewSitemap();if e:=s.Add("a\x00b","https://example.test/path");e!=nil{t.Fatal(e)};if got:=s.URLs("a");len(got)!=0{t.Fatal("prefix tenant leaked URLs",got)}
}
func TestRobotsRejectLineInjection(t *testing.T){
 for _,path:=range []string{"/ok\nDisallow: /","/ok\rUser-agent: secret","/ok\x00hidden"}{for _,allow:=range []bool{true,false}{r:=NewRobots();before:=r.Render("t");var e error;if allow{e=r.Allow("t",path)}else{e=r.Disallow("t",path)};if !errors.Is(e,ErrInvalidURL)||r.Render("t")!=before{t.Error("injected rule accepted",path,e)}}}
}
func TestLocalBusinessAbsoluteURL(t *testing.T){
 for _,raw:=range []string{"relative/path","javascript:alert(1)","https://:80/path"}{s,e:=(LocalBusiness{Name:"A",URL:raw}).JSONLD();if !errors.Is(e,ErrInvalidBiz)||s!=""{t.Error("invalid business URL accepted",raw,s,e)}}
}
func TestSitemapRejectsHostlessURL(t *testing.T){
 s:=NewSitemap();if e:=s.Add("t","https://:80/path");!errors.Is(e,ErrInvalidURL){t.Fatal("missing hostname accepted",e)}
}

````

## Probes originales social

````go
package social
import("testing";"reflect";"time")
func TestSamePostIDAcrossTenants(t *testing.T){
 s:=NewStore();a,b:=post(),post();b.TenantID="other";if e:=s.Schedule(a);e!=nil{t.Fatal(e)};if e:=s.Schedule(b);e!=nil{t.Fatal("other tenant cannot use its own ID",e)}
}
func TestDueHasNoTenantBoundary(t *testing.T){
 s:=NewStore();now:=time.Unix(100,0);s.clock=func()time.Time{return now};a,b:=post(),post();a.ScheduledAt=now.Add(time.Second);b.ScheduledAt=a.ScheduledAt;b.TenantID="other";b.ID="p2";if e:=s.Schedule(a);e!=nil{t.Fatal(e)};if e:=s.Schedule(b);e!=nil{t.Fatal(e)};now=now.Add(time.Second)
 for _,p:=range s.Due(){if p.TenantID!="t"{t.Error("due view cannot select requesting tenant",p.TenantID)}}
}
func TestMutationAndReadRequireTenantArgument(t *testing.T){
 typ:=reflect.TypeOf(NewStore());for name,n:=range map[string]int{"Due":2,"MarkPosted":3,"MarkFailed":4}{m,ok:=typ.MethodByName(name);if !ok||m.Type.NumIn()!=n{t.Error("API cannot express tenant scope",name)}}
}

````

## Cierre172 — controles generales y continuidad

Preflight171: **154 controles ejecutados PASS**,162packs,
1461archivos materializables,801Markdown y53perfiles. El inventario general
conserva Docker no disponible; no impidió esta ruta ni cambió el verifier.
6fuentes actualizadas exactas,32tests x3,15tests históricos preservados,
3firmas legacy rechazadas,vet/build,62semillas y4518771fuzz PASS.
Observability2fuentes sin cambios,28tests/vet/build PASS,CANDIDATE intacto.
Se verificaron 32SHA de recibos y los21SHA de owners del censo.

FAIL644/646/647/648 con corrección demostrada; FAIL645/649/650/651 conservan
fallos de conteo, ruta y generación del informe y su recuperación, sin ocultarlos
ni repetir checkpoints.48definiciones/requisitos/estados intactos:43PASS/5BLOCKED.
EVID-73 parcial;21cores inventariados no equivalen a21admitidos. Cierra la tanda
local de tres núcleos pendientes; la próxima fase es admisión/equivalencia e
integración del scope existente, no nuevas copias de dominio ni volver al conteo.
Sin corpus autorizado nuevo, servicio admitido, release final o aceptación
transferida del ZIP148. Checkpoint172/eventos append-only conserva continuidad.

| Recibo de cierre relativo al stage | SHA-256 |
|---|---|
| preflight171.json | `d6ea1d8f2ad7b891fe8ac25657415f987c1d7ff2dbe14d07f4fa63b76e47db21` |
| preflight171.log | `92d9cfa46f2acfb61c23635306a2d3c08990946cbb76a5d4564cd5f800d19c8b` |
| closure171.json | `7dad4f1e9559f13d9dc46ebe792a77f652319cacf7f0c53f73776fa679e1ff08` |
| resume171.log | `4bf1d2001bffc68c860594f0fa483723d95d8f44fd2d64abadc564030b0e113c` |
| plan171.log | `2fd65a7ad27c8b968c64f8079abf8610894d6e7732b0b0f514dba962cf93b565` |
