# V353 — auditoría de aceptación y corrección de límites FX

Fecha:2026-09-09. Clasificación: mantenimiento/revisión correctiva de biblioteca.
Checkpoint de entrada150 validado; estado posterior151 y cierre152 separados.
Stage observado: `<LOCALAPPDATA>/Temp/elite-v353-403d584c08784164bea1cbf2ea2d6566`. Evidencia previa: `<LOCALAPPDATA>/Temp/elite-v352-32e0c00909af4e31a62ae41f95568773`.

## Resultado y alcance

Se verificaron31filas SHA-256 y18recibos de subprocess de V352. Los48IDs,
acciones, oráculos y requisitos originales permanecen iguales. El harness
canónico V352 se reconstruyó y ejecutó desde cero contra el mismo ZIP completo:
`0590588c4ab139e86e581b13c1db561f9342462646105c109f190aeefc755085`. Verifica804fuentes+manifest=805entradas, hashes contra
manifest y source148 congelado, verificador real extraído,13comandos de guía,
8helps, fallo real de reemplazo y rollback exacto NEW/EXISTING, retry tras liberar
el lock, cadena de eventos preservada, append duplicado rechazado y55checks
de lifecycle PASS. No se encontró discrepancia en esa evidencia auditada.

TEST06/08 conservan PASS sólo para ese candidato interno148. No es release final,
no incluye estudio UX humano independiente ni prueba de corte eléctrico. Los9
owners de tooling/guía consumidos permanecen iguales. El delta actual SÍ cambia
código FX: el ZIP anterior contiene0.1.0; no incluye esta corrección0.1.1. TEST07
debe construir, comparar y revalidar el payload final antes de publicación.

## Defecto real y cambio canónico

GO-FX-CORE0.1.0 pasaba sus5tests históricos. Tres grupos adversariales adicionales
reprodujeron6fallos: tasas NaN/+Inf reemplazaban una tasa válida; montos NaN/+Inf
devolvían éxito; inputs finitos podían devolver+Inf por overflow o cero por
underflow sin error. Dos controles de infinito negativo ya eran rechazados.
FAIL-20260909-618 / LIB-FAIL-2334 conserva el antes y registra la regresión.

GO-FX-CORE0.1.1 valida inputs positivos finitos antes de mutar, conserva la tasa
anterior ante rechazo y añade ErrInvalidResult para resultado no positivo o no
finito. Mantiene subnormales positivos representables y el redondeo binary64.
Son2archivos AUTHORED/LicenseRef-Workspace-Owner; no se copió código financiero
externo ni se introdujeron dependencias, tasas o políticas monetarias.

Se corrigió además el claim ISO4217: el código comprueba tres letras ASCII
mayúsculas, no pertenencia al catálogo ISO. El pack sigue SUPPORTED_REFERENCE,
REBUILD_VERIFIED / CONDITIONED y fuera de todos los pack plans. Los otros161
packs no cambiaron. Franquicia permanece67packs/746archivos. No se demuestra
compatibilidad con commerce0.6.1; compatible_with0.1.x permanece como dato legado.

## Verificación ejecutada

Reconstrucción desde los bloques canónicos:2/2fuentes byte-idénticas al candidato.
9tests ordinarios x3 y14semillas PASS; go vet ./... y go build ./... PASS.
GO_NATIVE_FUZZ_GATE con presupuesto10s:13.047.867ejecuciones PASS. El oráculo
interpreta bits IEEE754, multiplica con math/big.Rat exacto y convierte el racional
a binary64 para contrastar bits; no llama al helper de validación de producción.
Comprueba clasificación de errores, preservación de tasas y ausencia de mutación.
El gate de fuzz también pasó su selftest:1positivo y2negativos. go.mod existe
sólo en el montaje local de calificación; no es tercer archivo del pack.

Toolchain existente Go1.26.7 Windows/amd64, Python3.14.4, PowerShell7.6.5;
GOTOOLCHAIN=local, GOWORK/GOPROXY/GOSUMDB=off, CGO_ENABLED=0. Sin upgrade.
La semántica se contrastó con autoridades oficiales [tipos flotantes Go](https://go.dev/ref/spec#Floating_point_types)
y [math](https://pkg.go.dev/math#IsNaN). La versión actual mostrada por pkg.go.dev
no fue adoptada. No se afirma identidad numérica entre arquitecturas ni política
de redondeo monetario, proveedor/frescura de tasas, ledger durable o autorización.

## Estado de entrega

43/48controles pasan; TEST02/03/05/07/09 siguen bloqueados.5/48=10,4167% de los
controles, no del trabajo total. TEST02 recibe evidencia parcial de corrección,
no admisión completa de los21cores. Roadmap1completo/10abiertos, BENCH01pendiente.
No se inventa historial/corpus autorizado, servicio instalado ni release firmado.
Preflight151 pendiente en este punto; su resultado se registra al cierre.

## Recibos preservados

| Archivo relativo al stage V353 | SHA-256 |
|---|---|
| before151/implementation_packs/GO_FX_CORE.md | `5086a042d522bd756a12a315412a763e9943ade3f88ed408205532ce03731de8` |
| prior-evidence-audit.json | `0814acb8ee1673803d0d411a3a5ab3d2ccc5a638fc9918fd7ff87394aa0239df` |
| audit353.py | `962183fff0e42549538874a8579874c4fcd97f049c691205fb335455474c05ef` |
| fresh-acceptance.log | `aace1066fd11fe00c78402ab374261aa279adc1399c67cf54530dafa0e301fb9` |
| fresh-acceptance/acceptance.json | `d6e3d1c416a62799d2fa27ac24273df29e3af0dc73bd67619502358dbd006498` |
| fresh-acceptance/archive-verification.json | `efddef232819fd41409cf1420a785891568b314182ce1f434cf558833b6aefd5` |
| fresh-acceptance/guide-run/results.json | `98108aad5b1b92dab5c2c096ce2153fe86b25d4868d4f451c50d1878fa2594c0` |
| fx-baseline.log | `4152531cc2d56f0b83798d18abe99ae532ad64075f1a87cfd006904f9b1902a7` |
| fx-red.log | `ca615ed15c0f2096a8365565f9f2a9396568514869cc84504d90f57802a3ff8d` |
| red353.py | `0219178ecef6b99b50ec52bb7886f198908f3f70b9b7436e98e3817309805044` |
| fx-green.log | `6b8ad60cfde605045f2d7d81228a0064d1168eebc24b43097806975e2bfddcbd` |
| fx-vet.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| fx-rebuilt.log | `fe7bd2cf70da9df55fc4485b32dd883e9b2c70f6f42a235ae3c96cca2a079359` |
| rebuilt-tests.log | `c4993b38bc080376846cf4c90a8ba49d3e70c40b2fbdef685f0e849e1269c755` |
| rebuilt-vet.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| rebuilt-build.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| fuzz-profile.json | `8641079402c6687649478a4b59226c8d766391d7343a79862019cb046a6a9391` |
| fuzz-receipt.json | `d84e42bf97b3cb376f1765ad1628a11f98fbe7454d0ce207647ced374a0053c7` |
| fuzz.log | `9d7f9952fa9612f90d032be37f26580bf635a34d2c733ee85fad1cd420ea93dd` |
| fuzz-self.log | `1540f2abdd5365ec40d3e10e7caf332e2869a1504c599615e4e67d94526d299c` |
| source-delta151.json | `fb6e1b27b1067312ce8a6fe684f9c954b8a198260150ebe74e1f5197f92a4940` |
| ledger151.json | `960bdd71d5a192ceed55841cc8df1f083e032c5dc422985e3073f47c719188f9` |

### Reproducción de la regresión contra0.1.0

Conservar los5tests originales, añadir este archivo y ejecutar go test ./...
con el toolchain fijado. Los6subcasos descritos deben fallar antes y pasar después.

````go
package fx
import("testing";"math";"errors")
func TestNonFiniteRateDoesNotReplacePrior(t *testing.T){
 for name,v:=range map[string]float64{"nan":math.NaN(),"positive-infinity":math.Inf(1),"negative-infinity":math.Inf(-1)} {t.Run(name,func(t *testing.T){
  s:=NewStore();if e:=s.SetRate("t","USD","ARS",2);e!=nil{t.Fatal(e)}
  if e:=s.SetRate("t","USD","ARS",v);!errors.Is(e,ErrInvalidRate){t.Errorf("invalid rate accepted: %v",e)}
  if got,ok:=s.Rate("t","USD","ARS");!ok||got!=2 {t.Errorf("invalid update replaced valid rate: %v",got)}
 })}
}
func TestNonFiniteAmountRejected(t *testing.T){
 for name,v:=range map[string]float64{"nan":math.NaN(),"positive-infinity":math.Inf(1),"negative-infinity":math.Inf(-1)} {t.Run(name,func(t *testing.T){
  s:=NewStore();_ = s.SetRate("t","USD","ARS",2);got,e:=s.Convert("t","USD","ARS",v)
  if !errors.Is(e,ErrInvalidAmount)||got!=0 {t.Errorf("invalid amount returned result=%v error=%v",got,e)}
 })}
}
func TestNonRepresentableProductRejected(t *testing.T){
 for name,p:=range map[string][2]float64{"overflow":{math.MaxFloat64,2},"underflow-to-zero":{math.SmallestNonzeroFloat64,0.25}} {t.Run(name,func(t *testing.T){
  s:=NewStore();if e:=s.SetRate("t","USD","ARS",p[0]);e!=nil {t.Fatal(e)}
  got,e:=s.Convert("t","USD","ARS",p[1]);if e==nil||got!=0 {t.Errorf("nonrepresentable product returned result=%v error=%v",got,e)}
 })}
}

````

## Cierre152 — resultado verificado

Preflight151 finalizó: 154 controles ejecutados PASS,162packs,
1461archivos materializables,791Markdown y53perfiles. Docker permanece ausente;
las comprobaciones omitidas/condicionadas del target no se cuentan como ejecutadas.
Los48tests originales conservan definición y estado; requisitos idénticos.
43passed/5blocked/0planned. Se incorporó evidencia EVID-63, sin relajar criterios.

Se conserva el antes de FAIL618 y su regresión demostrada, más FAIL619 recuperado
por ejecución explícita de las operaciones omitidas. Pack FX0.1.1 sigue condicionado.
Los9owners consumidos por aceptación148 siguen iguales; el payload completo ya
cambió por FX y evidencia. No se hereda esa aceptación como release actual.
Checkpoint152 conserva hashes actualizados y cadena append-only; su recibo de
validación se escribe después de este reporte para evitar dependencia circular.

| Recibo de cierre relativo al stage | SHA-256 |
|---|---|
| preflight151.json | `70eac8f114855601df3f48baf950cf134d302b06ae4443bc6c6df27750f87017` |
| preflight151.log | `8ffcf65a167ea47b7752240666ffc1a3d9b3d0ad675cf243fd7fdc23348dc0ff` |
| closure151.json | `78664d4f5192c30e772a4a21879bfc58da3771c61a0aaf4cdbe4afce458f0a9c` |
| resume151.log | `a4efd416932c7737118f88456dd06555941073cdede10ff656704128f11ded72` |
| plan151.log | `c912fb6cdf9a56e15785aad245756d19de1ae3116be72c57e00780542982c232` |
