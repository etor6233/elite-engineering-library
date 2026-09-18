# V359 — aislamiento del checklist y los reclamos

2026-09-09. Mantenimiento correctivo y revisión de TEST02/03. Checkpoint 162
validado antes de editar; integración 163 y cierre 164. Stage `<LOCALAPPDATA>/Temp/elite-v359-cb0d32db936c4a41bce814e26edcadf5`.

## Resultado observado

GO-ONBOARDING-CORE 0.1.0 y GO-WARRANTY-CLAIMS-CORE 0.1.0 superaban sus 9 tests
originales. Cuatro grupos nuevos fallan contra esas fuentes conservadas:

- Onboarding: (tenant="a\0b", franchisee="c") y ("a", "b\0c") colisionan en
  la clave concatenada. Start rechaza un checklist distinto; Progress lee los
  contadores ajenos y Complete modifica ese checklist. Done también expone el
  resultado de esa finalización indebida.
- Garantías: Open admite tenant con NUL y valida ID/subject mediante regex.
  La consulta alternativa con ID que contiene NUL no podría crear un reclamo,
  pero State y Review/Approve/Reject/Resolve/Close resolvían la clave ajena.
  Se reprodujo modificación en los cinco métodos y ambas ramas de Close.

Ambos 0.1.1 usan igualdad estructural de tenant/identidad. No hay normalización
ni nuevo filtro de caracteres: los dos campos se comparan completos y separados.
Se conserva API, validación, unicidad, estados y orden de errores. String de
onboarding usa el mutex de escritura; prueba concurrente local, sin -race.

## Semántica preservada

Onboarding copia y deduplica Steps con igualdad exacta, conserva orden de alta
y permite completar cualquier paso conocido. No impone prerequisites: hacerlo
sería una regla nueva. Complete repetido devuelve nil y no incrementa dos veces
el progreso. Start con identidad existente sigue ErrDuplicate. La metadata O(1)
de Complete se corrigió a O(N), porque recorre Steps; Progress también es O(N).
No alta real de franquiciado, capacitación evaluada, aprobación o persistencia.

Garantías fuerza StateOpen al abrir; requiere Review antes de Approve/Reject.
Resolve sólo desde approved; Close desde resolved o rejected. El diagrama previo
se acota a esas dos ramas, sin agregar una transición. Reject/Resolve validan nota
antes de identidad/estado y mantienen ErrNoteRequired prioritario. Una nota en
blanco no cambia estado; un reclamo cerrado mantiene su identidad. No eligibilidad,
término legal, devolución, reemplazo o reembolso ejecutado ni inferido.

## Evidencia ejecutada

4/4 fuentes se reconstruyen desde destinos nuevos y coinciden byte por byte.
17 tests (9 onboarding, 8 warranty), cada uno repetido 3 veces, vet/build PASS.
Los 9 cuerpos históricos preservan sus aserciones tras gofmt. Las regresiones
originales y sus fallos se guardan, sin cambiar el resultado anterior.

Las pruebas concurrentes usan 64 goroutines por fase. Onboarding permite un único
Start y todas las repeticiones válidas de Complete, con progreso final 2/2.
Garantías permite una mutación exitosa por fase Open/Review/Approve/Resolve/Close.
State/Progress/String se ejercitan durante concurrencia; sin detector de races,
benchmark de carga, proceso durable o proveedor real.

Los oráculos de fuzz son listas lineales, con comparaciones directas de campos y
validación independiente; no llaman key, validate ni transition de producción.
El checklist usa slices de flags en lugar del map de completados. Cada operación
comprueba Progress/Done de todas las identidades candidatas. Garantías aplica
su tabla de estados y prioridad de nota y compara cada State visible. No existe
getter público de Note: esta suite no atribuye una prueba de retención durable.

| Target | Semillas | Presupuesto | Workers | Ejecuciones PASS |
|---|---:|---|---:|---:|
| FuzzChecklistProgressModel | 12 | 10s | 4 | 2415852 |
| FuzzClaimLifecycleModel | 18 | 10s | 4 | 2344743 |

Total **4760595** ejecuciones, sin exhaustividad. Inputs sintéticos Unicode/NUL,
blancos, duplicados, desconocidos, errores y repeticiones; strings hasta 64 bytes,
32 pasos de checklist o 40 de reclamo por caso. Go 1.26.7 Windows/amd64 existente,
Python 3.14.4, PowerShell 7.6.5. GOTOOLCHAIN=local, GOWORK/GOPROXY/GOSUMDB=off,
CGO_ENABLED=0, GOMAXPROCS=4, GOFLAGS=-parallel=4. Sin nuevas dependencias.
La limitación del fuzz a 24 workers de V354 no se declara reparada.

## Implementation assurance y límites

| Dimensión | Evidencia / condición pendiente |
|---|---|
| Corrección | red/green, modelo independiente, fuentes exactas, tests/vet/build |
| Integración | API estable, sin activación/servicio de devoluciones conectado |
| Seguridad/privacidad | aislamiento sintético; caller debe autorizar tenant/actor; sin SCA integrada |
| Resiliencia | rechazos sin efecto local, concurrencia por mutex; sin restart durable |
| Rendimiento | claims O(1) esperado; checklist O(N); sin benchmark target |
| Operación | memoria de proceso, sin supervisor, alertas ni retención |
| Recovery | original preservado; rollback reabre fallos, no restore de datos |
| Supply chain | stdlib existente, cuatro SHA; sin nuevo SBOM/firma/release |
| Usabilidad | reglas y errores explícitos; sin UX ni aceptación comercial |

AUTHORED / LicenseRef-Workspace-Owner, SUPPORTED_REFERENCE,
REBUILD_VERIFIED / CONDITIONED. Compatibilidad histórica no implica integración
actual. V206/V214 siguen como historia; sus composiciones no se reejecutaron.
160 otros packs intactos, estos dos fuera de pack plans; franquicia 67/746 intacta.
V293 conserva owners de organizaciones/identidad/servicio y gaps de evaluación,
eligibilidad/estados específicos/recuperación; esta corrección no duplica dominio.

EVID-69 sólo agrega evidencia parcial TEST02/03. Criterios y estados de los 48
controles intactos: 43 PASS/5 BLOCKED, 1 macro-tarea completa/10 abiertas.
10,4167% de controles pendientes no mide esfuerzo global. TEST02/03/05/07/09
permanecen abiertos. ZIP148 carece de estas correcciones y las de V353–358;
su aceptación no se hereda a un payload diferente. Preflight163 pendiente.

## Recibos preservados

| Archivo relativo al stage | SHA-256 |
|---|---|
| baseline162.json | `911d71b16c2a30199832412bf730ffe9cc4aeeab5924507661f9234c73c401cf` |
| resume162.log | `0f7c00f0174f459a724cf0f877b17615ef268524624bf1958ad8ea33146c123a` |
| before163/implementation_packs/GO_ONBOARDING_CORE.md | `0348f51991a57109678e608451fb5e2cd39c12b19c2eef6fbdb30a9b5315aca5` |
| before163/implementation_packs/GO_WARRANTY_CLAIMS_CORE.md | `5633258019a207a5a6317e90123b77c20f1970e52331d8f6a0a7bfa368d2b723` |
| baseline-tests.log | `6565f3b803575914e8c39071070e3f38fa6dd49fcb9f109b5fa5299e1c257d68` |
| red-tests.log | `0e7524f2c6489a62547c3e2379a55ea5bdc64e3310cd4c404dc6c296347e2834` |
| red-summary.json | `82acd6fdaeca1d3ab14e1d65d6fffacea9e46c1fa5707c7e31e86ad061b9d5d6` |
| red359.py | `070aedac91270ac1f1aefccbdbbecdb5d587ec0c8a8118efb4d56875d7366c0d` |
| fix359.py | `78d52c25b1ff70e4901b78632564160ec8578675a8c8c4880138019eb37a4315` |
| green-tests.log | `069e457f9c152f8adcc67d0b987df15f6753ac539bad2c86c30c331d5518c433` |
| rebuilt-tests.log | `0e8cbb458f87ddb8015653e3661f8f0efe319a7607c7be3ccb1eb5bfe55eb0ce` |
| rebuilt-vet.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| rebuilt-build.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| fuzz-profile.json | `cce8332a15d72c3a7e19e8b5bf27d846072e21d3119540238a7550beef2d6830` |
| fuzz-receipt.json | `3e07af22fe8e4917ef3744a1ed617590cd35232af3ba054383b0321010131c4e` |
| fuzz.log | `7e967009f5e733cb24c08297bbdda0d14d114b58e333b774c97c755cd2996c28` |
| fuzz-self.log | `8971a6578399bbcdb9bfc6db8026b884ec505d830b4d90913a6bef675d561e4b` |
| qualification163.json | `3e53c573c5170d403cdb14efa9ebac1c96c8ee3bee9994b604c33fd9f7bf1fc0` |
| ledger163.json | `8e73aa35d259b26da17c905ec5b6d24c217039dd79b54334b2b9ccfae037ff49` |

## Probes originales onboarding

````go
package onboarding
import("testing";"errors")
func TestDistinctOnboardingTuplesDoNotCollide(t *testing.T){
 s:=NewStore();if e:=s.Start("a\x00b","c",[]string{"step"});e!=nil{t.Fatal(e)}
 if e:=s.Start("a","b\x00c",[]string{"other"});e!=nil{t.Fatalf("distinct onboarding rejected: %v",e)}
}
func TestAliasedOnboardingCannotReadOrComplete(t *testing.T){
 s:=NewStore();if e:=s.Start("a\x00b","c",[]string{"step"});e!=nil{t.Fatal(e)}
 if n,total,ok:=s.Progress("a","b\x00c");ok||n!=0||total!=0{t.Errorf("foreign progress: %d %d %v",n,total,ok)}
 if e:=s.Complete("a","b\x00c","step");!errors.Is(e,ErrNotFound){t.Errorf("foreign completion: %v",e)}
 if s.Done("a","b\x00c"){t.Error("foreign Done")}
 if n,total,ok:=s.Progress("a\x00b","c");!ok||n!=0||total!=1{t.Errorf("original progress changed: %d %d %v",n,total,ok)}
}

````

## Probes originales warranty

````go
package warranty
import("testing";"errors")
func TestAliasedClaimCannotReadForeignState(t *testing.T){
 s:=NewStore();if e:=s.Open(Claim{TenantID:"a\x00b",ID:"c",SubjectID:"item"});e!=nil{t.Fatal(e)}
 if st,ok:=s.State("a","b\x00c");ok||st!=""{t.Fatalf("foreign state: %q %v",st,ok)}
}
func TestAliasedClaimCannotTransition(t *testing.T){
 for _,action:=range []string{"review","approve","reject","resolve","close-resolved","close-rejected"}{t.Run(action,func(t *testing.T){
  s:=NewStore();if e:=s.Open(Claim{TenantID:"a\x00b",ID:"c",SubjectID:"item"});e!=nil{t.Fatal(e)};want:=StateOpen
  if action!="review"{if e:=s.Review("a\x00b","c");e!=nil{t.Fatal(e)};want=StateInReview}
  if action=="resolve"||action=="close-resolved"{if e:=s.Approve("a\x00b","c");e!=nil{t.Fatal(e)};want=StateApproved}
  if action=="close-resolved"{if e:=s.Resolve("a\x00b","c","resolved");e!=nil{t.Fatal(e)};want=StateResolved}
  if action=="close-rejected"{if e:=s.Reject("a\x00b","c","rejected");e!=nil{t.Fatal(e)};want=StateRejected}
  var e error;switch action{case "review":e=s.Review("a","b\x00c");case "approve":e=s.Approve("a","b\x00c");case "reject":e=s.Reject("a","b\x00c","reason");case "resolve":e=s.Resolve("a","b\x00c","resolution");default:e=s.Close("a","b\x00c")}
  if !errors.Is(e,ErrNotFound){t.Errorf("foreign transition got %v",e)}
  if st,ok:=s.State("a\x00b","c");!ok||st!=want{t.Errorf("original state changed: %q want %q",st,want)}
 })}
}

````

## Cierre 164 — verificación del delta completo

Preflight163: **154 controles ejecutados PASS**, 162 packs,
1461 archivos materializables, 797 Markdown y 53 perfiles. Docker ausente;
checks live omitidos no se cuentan como ejecutados. Ambos cores 0.1.1 mantienen
4 fuentes exactas, 17 tests x3, vet/build, 30 semillas y 4760595 ejecuciones
fuzz PASS. Los 9 tests originales conservan sus aserciones. FAIL632/633 tienen
regresión demostrada; el antes permanece en fuentes, pruebas y logs.

Se verificaron 19 SHA de los recibos listados antes de este cierre y
los 4 SHA de fuentes reconstruidas. Los 48 controles/requisitos mantienen sus
definiciones y estados: 43 PASS/5 BLOCKED. EVID-69 parcial para TEST02/03, sin
admisión integral ni integración productiva de estos cores. No se reabre una
decisión empresarial ni se habilita una activación/reembolso. Perfil 67/746
intacto; ZIP148 no recibe aceptación sobre un payload diferente.

Checkpoint164 conserva eventos append-only y refresca owners por hash. Sus
recibos se generan después del reporte para evitar una referencia circular.

| Recibo de cierre relativo al stage | SHA-256 |
|---|---|
| preflight163.json | `b1132c02052cb84d66439264ff03824ed446111cda3d5aadf96a8f446ff88538` |
| preflight163.log | `83d1f5aa800cc4d51d5ded2bda01dab37a5de9e38ff5f5bdaaa497b2307441e4` |
| closure163.json | `8452bb6414c4360dbacc85059ba7bdc2078a0127530c3a249fe8d75f4a0f82b3` |
| resume163.log | `bc3b2a57ddf591edc2113c6d821e93f01352b32bcab308f49242058124f59847` |
| plan163.log | `2fd65a7ad27c8b968c64f8079abf8610894d6e7732b0b0f514dba962cf93b565` |
