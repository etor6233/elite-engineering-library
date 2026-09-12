# V357 — identidad y estados del registro de referidos

2026-09-09. Mantenimiento correctivo/revisión, TEST02/T2802. Checkpoint158
validado antes de editar; integración159 y cierre160. Stage `C:/Users/NL/AppData/Local/Temp/elite-v357-32715539145f4f2f8293c1004dafa388`.

## Defecto y corrección

GO-REFERRALS-CORE0.1.0 superaba sus5tests originales. Tres regresiones de API
pública reproducen alias de (tenant,referee): ("a\0b","c") y ("a","b\0c")
comparten la cadena concatenada. El segundo Create se rechaza indebidamente.
Si sólo existe el primero, State del segundo devuelve estado vacío con true y
Convert/Reward retornan ErrBadTransition en vez de ErrNotFound. Si además existe
en tenant a otro referido legítimo con el mismo código, las consultas aliased
leen y Convert/Reward modifican ese otro registro. Son efectos locales reales
reproducidos; no se infiere explotación de autenticación de un producto instalado.

0.1.1 usa identity{tenant,id} como clave en ambos índices. La igualdad compara
los campos completos por separado, por lo que bytes NUL en cualquier campo no
pueden desplazar la frontera de la tupla. No se bloquean caracteres previamente
válidos, no se normalizan IDs ni se cambian firmas, errores o precedencia.
Create mantiene validación, rechazo de auto-referido, prioridad de código
duplicado, unicidad de referee por tenant y State inicial pending. Convert y
Reward conservan transiciones estrictas; repetición retorna error, no replay.

## Pruebas nuevas y conservadas

Los5tests originales mantienen cuerpos/aserciones tras gofmt. Tres regresiones
fallan contra0.1.0 y pasan contra0.1.1. Dos tests adicionales verifican rechazos
sin consumo de identidad, State impuesto por Create y prioridad de errores,
y64goroutines por fase con exactamente un éxito Create/Convert/Reward.
Reconstrucción desde destino nuevo:2/2fuentes byte-idénticas;10tests x3,
go vet y go build PASS. No se ejecutó -race ni carga productiva.

FuzzReferralSequenceModel usa una lista lineal de registros con comparación
directa de campos, validación de códigos ASCII propia y máquina de estados
independiente; no llama ck ni validate de producción. Después de cada operación
compara todas las identidades candidatas, incluso errores y lecturas sin efecto.
Datos sintéticos: duplicados, campos vacíos, espacios, NUL, Unicode, código
inválido, tenant ajeno y transiciones repetidas.12semillas,10s,4workers,
**1709469ejecuciones PASS**. Cotas: strings64bytes, operaciones160bytes/32pasos;
no exhaustividad ni representación de todo negocio.

Go1.26.7 exacto existente, Windows/amd64, Python3.14.4 y PowerShell7.6.5.
GOTOOLCHAIN=local; GOWORK/GOPROXY/GOSUMDB=off;CGO_ENABLED=0;
GOMAXPROCS=4/GOFLAGS=-parallel=4. Sin nuevas dependencias ni adquisición pública.
La limitación del fuzz con24workers de V354 no se declara resuelta.

## Implementation assurance y admisión

| Dimensión | Demostrado / límite |
|---|---|
| Corrección | baseline, red/green, modelo independiente,2SHA canónicos |
| Integración | API sin cambio; sin CRM/loyalty/checkout conectado |
| Seguridad/privacidad | aislamiento de claves sintético; caller debe autorizar tenant/persona |
| Resiliencia | rechazo sin cambio y transición local atómica; sin durabilidad/reinicio |
| Rendimiento | dos índices hash, coste esperado O(1) operaciones, memoria O(N); sin benchmark target |
| Operación | sólo proceso en memoria; sin monitor, alertas o retención productiva |
| Recovery | original preservado; rollback reabre FAIL628, no restore durable |
| Supply chain | sin dependencias nuevas; no SCA integrada, SBOM o release nuevos |
| Usabilidad | límites y errores documentados; sin aceptación de flujo comercial |

Convert no prueba conversión comercial y Reward sólo cambia State; no emite
recompensa, saldo ni cupón. No elegibilidad, política, antifraude o permisos nuevos.
Procedencia AUTHORED / LicenseRef-Workspace-Owner, SUPPORTED_REFERENCE /
REBUILD_VERIFIED / CONDITIONED. La compatibilidad0.1.x declarada con loyalty
es histórica, no integración aprobada. Suite V196 histórica preservada.

V293 asigna referidos al origen lead/CRM y mantiene la brecha de elegibilidad
y recompensa durable.161otros packs intactos; este core no figura en pack plans.
Franquicia67packs/746archivos sin cambio. Evidencia parcial TEST02 y aislamiento
TEST03; ninguno se cierra.43/48controles PASS,5BLOCKED (10,4167% de controles,
no de esfuerzo). Roadmap1completo/10abiertos. TEST02/03/05/07/09 pendientes.
ZIP148 no contiene estas correcciones ni las seis anteriores; su aceptación
no se transfiere al payload final modificado. Preflight159 pendiente al integrar.

## Recibos y fuentes conservadas

| Archivo relativo al stage | SHA-256 |
|---|---|
| baseline158.json | `55a0d746f6fe74068d49ac86c5b8d9772788230c9cc2361a54073f0f37429754` |
| before159/implementation_packs/GO_REFERRALS_CORE.md | `4397e3aeb0236f4ab36f83dd910e016cdda24fef905580c14999d39d7049f28e` |
| resume158.log | `41738146aec2a0219bf598165fec5e572dc2f0fd28107f1bcca9ea1a5d0c2cda` |
| baseline-tests.log | `df5e242d28cb631c4ba2a9290912385f25975fa338a9cba99c8004e32b4d9c69` |
| red-tests.log | `7fa3d731ba760a90970a042281a11bf2bd00fe2e22e786059d5b79e71b9f8dc2` |
| red-summary.json | `59e59dbf8e00474f711c9b687c073248350b4cb0c209f2e0c72e91121edca41a` |
| red357.py | `2bab343ea3e16f14af33225b39e1b7f9c5c08dba06d8613a0fab82a68cdba26b` |
| fix357.py | `94bb788fd200f3ef8cc93c3857d6e11e247fbdfc561139ddf05d198ac1da9281` |
| green-tests.log | `26ce252b907289c53224f0ed4662125a7d833ad691775a7aeb5b4099eee6dadc` |
| rebuilt-tests.log | `f640a7c6550fa00c17895c0e004c360d565842d23b65a5419cf5bea8490c5fa7` |
| rebuilt-vet.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| rebuilt-build.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| fuzz-profile.json | `7a82e73d46aaab4ccdef2e14044aaf2958614ee41bc438a59e7521ed2d496dc3` |
| fuzz-receipt.json | `620b4f24afa0ec3e67107441769956661ba2e961f10c2f12e2f580d8623365a7` |
| fuzz.log | `6f8377e533ba02141581e882438164879ee96dac55a05e7080f9a9b8614d20c7` |
| fuzz-self.log | `8f3d3a6d5dd03cacd1390b74aa7421ddd6da9667ba898c7a7f936dd42583538c` |
| qualification159.json | `cb5168a71994823684b013baa5d0b210e4491a2605886b71f2b6440c7280af24` |
| ledger159.json | `ac6c8a3f75ac1c5f26aba33b29efe712d21c4434811b1812fdb9fbcff9846354` |

## Regresiones reproducidas

````go
package referrals
import("errors";"testing")
func identityReferral(tenant,code,referee string)Referral{return Referral{TenantID:tenant,Code:code,ReferrerID:"source",RefereeID:referee}}
func TestDistinctReferralTuplesDoNotCollide(t *testing.T){
 s:=NewStore();if e:=s.Create(identityReferral("a\x00b","FIRST","c"));e!=nil{t.Fatal(e)}
 if e:=s.Create(identityReferral("a","SECOND","b\x00c"));e!=nil{t.Fatalf("distinct identity rejected: %v",e)}
}
func TestAbsentAliasedRefereeIsNotFound(t *testing.T){
 s:=NewStore();if e:=s.Create(identityReferral("a\x00b","SHARED","c"));e!=nil{t.Fatal(e)}
 if st,ok:=s.State("a","b\x00c");ok||st!=""{t.Errorf("absent referee reported present: %q %v",st,ok)}
 if e:=s.Convert("a","b\x00c");!errors.Is(e,ErrNotFound){t.Errorf("convert got %v",e)}
 if e:=s.Reward("a","b\x00c");!errors.Is(e,ErrNotFound){t.Errorf("reward got %v",e)}
}
func TestAliasedLookupCannotTransitionAnotherReferee(t *testing.T){
 s:=NewStore();for _,r:=range []Referral{identityReferral("a\x00b","SHARED","c"),identityReferral("a","SHARED","legitimate")}{if e:=s.Create(r);e!=nil{t.Fatal(e)}}
 if st,ok:=s.State("a","b\x00c");ok||st!=""{t.Errorf("read unrelated referral: %q %v",st,ok)}
 if e:=s.Convert("a","b\x00c");!errors.Is(e,ErrNotFound){t.Errorf("converted unrelated referral: %v",e)}
 if e:=s.Reward("a","b\x00c");!errors.Is(e,ErrNotFound){t.Errorf("rewarded unrelated referral: %v",e)}
 for _,p:=range [][2]string{{"a\x00b","c"},{"a","legitimate"}}{if st,ok:=s.State(p[0],p[1]);!ok||st!=StatePending{t.Errorf("other identity changed: %q %v",st,ok)}}
}

````

## Cierre160 — verificación del delta completo

Preflight159: 154controles ejecutados PASS,162packs/1461archivos/
795Markdown/53perfiles. Docker ausente; checks live omitidos no se cuentan como
ejecutados. Referrals0.1.1 conserva2fuentes reconstruidas,10tests x3,vet/build,
12semillas y1709469ejecuciones fuzz PASS. FAIL628 REGRESSION_PROVEN.

Los48tests y requisitos conservan definición/estado:43PASS/5BLOCKED.
EVID-67 sólo agrega evidencia parcial TEST02/03; no cierre de admisión integral.
AUTHORED/CONDITIONED, perfil67/746 intacto; sin emisión de recompensa ni release.
Checkpoint160 mantiene hash/eventos append-only; su receipt se genera después
de este informe para evitar una dependencia circular de hashes.

| Recibo de cierre relativo al stage | SHA-256 |
|---|---|
| preflight159.json | `344044634848cd8fc61611b7bc1a727da8a4db8c20395de024a12cd7548a3fbe` |
| preflight159.log | `3336c4a79e418f7b1136397e2540137185d8a7fd55c9f97261e5f95a1b8349b6` |
| closure159.json | `3c14982d1a1b728aa31aaa533f5b105b9ffca80a38391c779260f9f73e94cd76` |
| resume159.log | `87b0a6a21fe86ca275e0a3a6e35d5cfd4e4133309a775adf302c968c46a8a01a` |
| plan159.log | `2fd65a7ad27c8b968c64f8079abf8610894d6e7732b0b0f514dba962cf93b565` |
