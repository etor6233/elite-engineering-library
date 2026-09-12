# V356 — aritmética de puntos básicos y neto sin overflow

2026-09-09. Mantenimiento correctivo/revisión TEST02, T2802. Entrada156
validada; integración157 y cierre158 posteriores. Stage `C:/Users/NL/AppData/Local/Temp/elite-v356-558d2ce3a1ec4275aa11a16eafd3715c`.
No reglas comerciales/laborales/fiscales nuevas ni cuentas, documentos o efectos.

## Defectos reproducidos

Los11tests originales (promotions6/payroll5) pasaban.3grupos nuevos fallaron:
porcentaje de promociones, porcentaje de nómina y neto negativo que wrappea
a positivo. Diez subcasos porcentuales fallan con MaxInt64 y2/1000/5000/9999/
10000bps;1bps es control que ya pasaba. Con gross=0 y dos deducciones fijas
MaxInt64, payroll0.1.0 devolvía2 como net y guardaba éxito, aunque el resultado
matemático era negativo. Se conserva el antes en FAIL625/626 y logs red.

Ambos cores0.1.1 corrigen esos defectos. Ninguna firma pública cambia. Los11
tests previos conservan sus cuerpos/aserciones; no se rebajó el contrato.

## Razón de corrección y semántica conservada

Para amount>=0 y0<=bps<=10000, escribir amount=10000*q+r con0<=r<10000 permite:

`floor(amount*bps/10000) = q*bps + floor(r*bps/10000)`.

q*bps<=amount, r*bps<=99.990.000 y el resultado final<=amount<=MaxInt64.
No intermedio necesita representar amount*bps, que era el origen del overflow.
Promotions exige bps>0; payroll admite0bps. Se preserva truncado entero positivo,
porcentaje sobre gross y redondeo individual por deducción. No se elige una
política monetaria nueva ni se afirma que la existente sea legalmente correcta.

Payroll evita restar cuando la deducción excede net disponible y retiene un flag.
Como todas las deducciones válidas son no negativas, ninguna regla posterior
puede recuperar ese neto. Sigue validando reglas posteriores: ErrInvalidRun
mantiene prioridad sobre ErrNegativeNet. Sólo guarda después de cálculo válido;
un fallo no consume ID. La detección de duplicado permanece en el orden previo.

Promotions conserva expiración inclusiva at==ExpiresAt, mínimo, cuota y truncado;
Redeem exitoso consume un uso incluso si el descuento es cero. No existe ID por
pedido, replay seguro ni reversión. Coupon válido se copia al store; tests no
promueven modificación externa de sus maps privados. Payroll conserva rechazo
ErrDuplicate por tenant/run; no es replay que devuelve resultado previo.

## Gates ejecutados

Reconstrucción4/4fuentes byte-idénticas.18tests (9por core) x3,vet/build PASS.
64goroutines prueban límite de7cupones con importe máximo y una nómina exitosa
por tenant con ID compartido, sin resultados corruptos. Rechazos no consumen
ID/cuota. No se ejecutó -race ni una suite de carga del target.

Oráculos de fuzz con math/big.Int calculan producto/suma completos y división
exacta entera; no replican la descomposición ni el flag de producción. Payroll
valida todas las reglas antes de su suma exacta y compara prioridad/error, net,
storage y reutilización tras rechazo. Promotions compara descuento y cuota.

| Target | Semillas | Presupuesto | Workers | Ejecuciones PASS |
|---|---:|---|---:|---:|
| FuzzPercentageCoupon |13|10s|4|2443829|
| FuzzPayrollExactDeductions |24|10s|4|2170609|

Total4614438ejecuciones acotadas, no exhaustividad. Go1.26.7 Windows/amd64,
Python3.14.4, PowerShell7.6.5 existentes. GOTOOLCHAIN=local, GOWORK/GOPROXY/
GOSUMDB=off,CGO_ENABLED=0,GOMAXPROCS=4,GOFLAGS=-parallel=4. La limitación V354
con24workers permanece; no se declara reparado el runtime. Ninguna dependencia
nueva: math/big sólo aparece en tests y pertenece a la stdlib ya fijada.

## Assurance y límites

| Dimensión | Evidencia / condición pendiente |
|---|---|
| Corrección | red/green, prueba algebraica, modelos exactos y rebuild |
| Integración | firmas estables; sin checkout/ledger/recibos empresariales montados |
| Seguridad/privacidad | tests tenant/cuota locales; no auth ni SCA integrada |
| Resiliencia | rechazos sin efecto local; sin restart durable |
| Rendimiento | O(1) descuento/O(D) nómina; sin benchmark de producción |
| Operación | memoria de proceso; no supervisor/monitoring/retención |
| Recovery | fuente anterior preservada; rollback reabre errores, sin restore durable |
| Supply chain/release |4fuentes con SHA; no release/SCA/firma nuevos |
| Usabilidad | errores y límites documentados; sin UX/aceptación de negocio |

Ambos packs AUTHORED/LicenseRef-Workspace-Owner, SUPPORTED_REFERENCE,
REBUILD_VERIFIED / CONDITIONED.160otros packs intactos; ninguno de estos dos
aparece en pack plans y franquicia67/746 no cambia. Fuentes Go/manuales gobiernan
el método numérico, no una arquitectura financiera prestada ni una política fiscal.
Compatibilidad histórica no equivale a integración con owners actuales de V293.
TEST02 recibe evidencia parcial;43/48checks mantienen criterios/estado,5bloqueados.
No porcentaje global de esfuerzo. TEST03/05/07/09 y admisión integral pendientes.
ZIP148 carece de estos fixes y los anteriores; no se hereda su aceptación a un
payload final distinto. Preflight157 pendiente antes del cierre.

FAIL627: una nueva nota escribió V209 como referencia histórica de payroll;
el footer canónico identifica V216. Se preservó la nota errónea y se corrigió
a V195/V216 sin alterar fuentes/test ni afirmar reejecución de aquellas suites.

## Recibos

| Archivo relativo al stage | SHA-256 |
|---|---|
| baseline156.json | `85a18669974382acdc0ebdc07d9c3e7dc5e324f61492673bb41413fb4b765614` |
| before157/implementation_packs/GO_PROMOTIONS_CORE.md | `ef721b1c53b55fbc39ee251a6e15184350ca01bf0dfef105aa2bc043eea91780` |
| before157/implementation_packs/GO_PAYROLL_CORE.md | `b23d1ca807194768b2bf095bde037e186338c07dc528fe83993b8d7d30e3804e` |
| baseline-tests.log | `83fe9ea5aed41c89d39297b5b754c24ce9b1a7deeeb023fbb0430e5b24e3cfd5` |
| red-tests.log | `ad685929e1c19f2d61b49be00e70245bef0e7e4a23167a0eecb5bf6c1a6376d1` |
| red-summary.json | `e6637cde20d9b87220ace60424499ac422d6f3823fb79c69a789ecb60fe840b3` |
| work356.py | `40bf8ed241bf017806aa2a84f52f3e5831152fda6dea97d834edc9a258512347` |
| qualify356.py | `3a5919d09d8139ceca74b85b8456e34a1db9cb726b533b6c5edaf34da7df3349` |
| GO_PROMOTIONS_CORE.md.before-reference-fix.txt | `c9aaba9f28586deaa604d5a85198c1e22de0716734f3513036412631040c4a39` |
| GO_PAYROLL_CORE.md.before-reference-fix.txt | `2aa87130ad1b03cc4c3525c54d7ce39d4eddc6572a345bd5533b21f33a3bd878` |
| green-tests.log | `84d331dbff5c2cffb0067f32a734eef354564f0fb67725697e7a6073f61c9922` |
| rebuilt-tests.log | `fe0ef04993d9484ef55557ef61fbf7fe9b46392f385b3438ebb9d6ad2a456302` |
| rebuilt-vet.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| rebuilt-build.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| fuzz-profile.json | `4a38eb5891d49ca33edbb2e0011c5f59b91489949aab9a9fac9ffab86c6d92c8` |
| fuzz-receipt.json | `164147a15a0ed2b8225da69b650e18c5d1043e58ed9935eeef491a68d530c7a4` |
| fuzz.log | `afbc9a1012bad10f5db27748c652982ca0563eaa20fe79d0775be1ee2d4bd160` |
| fuzz-self.log | `2f795cf18c17044de7e23b6f85a03fae84fb39783aca0a4d6905d56449bd1b70` |
| qualification157.json | `d5454cbbc9ec6688678fd6bb74958ec7b90e2149bce67724ab4c1f654715e1d8` |
| ledger157.json | `177a1bc455a626ff4e32f9adce65cd1d30a08f295692861c31a8a307b3f33ed3` |

### Probes originales promotions

````go
package promotions
import("testing";"math";"math/big";"fmt";"time")
func TestFullRangePercentageDiscount(t *testing.T){
 for _,bps:=range []int64{1,2,1000,5000,9999,10000}{t.Run(fmt.Sprint(bps),func(t *testing.T){
  s:=NewStore();now:=time.Unix(1000000,0);s.clock=func()time.Time{return now}
  c:=Coupon{TenantID:"t",Code:"C",Kind:KindPercent,Value:bps,MaxUses:1,ExpiresAt:now.Add(time.Hour),Active:true}
  if e:=s.Create(c);e!=nil{t.Fatal(e)};exact:=new(big.Int).Mul(big.NewInt(math.MaxInt64),big.NewInt(bps));exact.Quo(exact,big.NewInt(10000))
  got,e:=s.Redeem("t","C",math.MaxInt64);if e!=nil||got!=exact.Int64(){t.Fatalf("discount=%d want=%d error=%v",got,exact.Int64(),e)}
 })}
}

````

### Probes originales payroll

````go
package payroll
import("testing";"math";"math/big";"fmt";"errors")
func TestFullRangePercentageNet(t *testing.T){
 for _,bps:=range []int64{1,2,1000,5000,9999,10000}{t.Run(fmt.Sprint(bps),func(t *testing.T){
  s:=NewStore();r:=Run{TenantID:"t",ID:"r",EmployeeID:"e",GrossMinorUnits:math.MaxInt64,Deductions:[]Deduction{{Kind:KindPercent,BasisPoints:bps}}}
  deduction:=new(big.Int).Mul(big.NewInt(math.MaxInt64),big.NewInt(bps));deduction.Quo(deduction,big.NewInt(10000));want:=new(big.Int).Sub(big.NewInt(math.MaxInt64),deduction)
  got,e:=s.Execute(r);if e!=nil||got!=want.Int64(){t.Fatalf("net=%d want=%d error=%v",got,want.Int64(),e)}
 })}
}
func TestNegativeNetCannotWrapToStoredSuccess(t *testing.T){
 s:=NewStore();r:=Run{TenantID:"t",ID:"r",EmployeeID:"e",GrossMinorUnits:0,Deductions:[]Deduction{{Kind:KindFixed,AmountMinorUnits:math.MaxInt64},{Kind:KindFixed,AmountMinorUnits:math.MaxInt64}}}
 if got,e:=s.Execute(r);!errors.Is(e,ErrNegativeNet)||got!=0{t.Fatalf("negative mathematical net accepted: %d %v",got,e)}
 if _,ok:=s.Net("t","r");ok{t.Fatal("rejected run stored")};r.Deductions=nil
 if got,e:=s.Execute(r);e!=nil||got!=0{t.Fatalf("failure consumed identity: %d %v",got,e)}
}

````

## Cierre158 — verificación completa del delta

Preflight157 completó 154controles ejecutados PASS,162packs,
1461archivos materializables,794Markdown y53perfiles. Docker permanece ausente;
checks omitidos/condicionados del target no se cuentan como ejecutados.
Promotions/payroll0.1.1 conservan4fuentes reconstruidas,18tests x3,vet/build,
37semillas y4614438ejecuciones fuzz PASS.11tests originales mantienen aserciones.
FAIL625/626 tienen regresión demostrada; FAIL627 conserva antes y referencia
corregida, sin modificar código ni reescribir evidencia histórica.

Los48tests y requisitos originales conservan definición/estado:43PASS/5BLOCKED.
EVID-66 sólo suma evidencia parcial a TEST02. Ambos packs siguen CONDITIONED y
fuera del perfil; no política fiscal/comercial ni compatibilidad actual admitida.
TEST02/03/05/07/09 siguen abiertos,1macro-tarea completa/10abiertas. No release
final ni transferencia silenciosa de aceptación ZIP148 al payload modificado.
Checkpoint158 actualiza hashes y mantiene eventos append-only; su recibo de
validación se genera después de este informe, fuera del reporte hash-linked.

| Recibo de cierre relativo al stage | SHA-256 |
|---|---|
| preflight157.json | `dd2c3c40a42bd71c765bbb0571a07eba512fbc43d1a588161a21700169b847b5` |
| preflight157.log | `a6ee5b8d2fcb35f7c2d283e1792b189b4573fd6e57291c7af648c59f6d6dc0ad` |
| closure157.json | `7548fb6b964e065fc33ec51dbce2c4aba63d477d2e097792bc21f4d635a834c1` |
| resume157.log | `ab2384ffb18e13ce71c1897d0f9aa04a0f749c8d76f5f4b93e38d855958ae293` |
| plan157.log | `c912fb6cdf9a56e15785aad245756d19de1ae3116be72c57e00780542982c232` |
