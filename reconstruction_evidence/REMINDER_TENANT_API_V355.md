# V355 — API de recordatorios con tenant explícito

2026-09-09. Mantenimiento correctivo/revisión, TEST02/03 y T2802/T2805/T2809.
Entrada154 validada; integración155 y cierre156. Stage `C:/Users/NL/AppData/Local/Temp/elite-v355-9b316eaf126843548066c004205b8276`.
No nuevo producto, cuenta, provider, datos reales, política comercial o envío.

## Defecto reproducido y reparación

GO-REMINDERS-CORE0.1.0 declaraba tenant-scoped, pero Schedule indexaba únicamente
por ID, Due() devolvía pendientes de todos los tenants y MarkSent(id) no podía
representar el tenant del caller.5tests originales PASS (la metadata decía4);
3tests nuevos FAIL: mismo ID entre tenants, lectura de otro tenant y ack ajeno.
Los probes de0.1.0 son adapters de calificación que evidencian la ausencia de
scope; no son implementación de compatibilidad ni entran al código productivo.
String leía len(map) sin mutex: hallazgo estático, sin ejecución de detector races.

0.2.0 introduce un cambio de API deliberado: Due(tenant) y MarkSent(tenant,id),
sin overload global ni shim. La clave interna separa tenant/ID. El mismo ID puede
existir en tenants distintos, los duplicados dentro del tenant siguen rechazados,
la lectura retorna copias del tenant exacto y el ack no modifica otro tenant.
String usa el mismo mutex de los writers; su contador global es diagnóstico
interno, no una vista tenant autorizada. Pasar un string no autentica al actor.

## Migración y semántica observable

| Antes0.1.x | Ahora0.2.0 | Condición |
|---|---|---|
| Due() | Due(tenant) | tenant validado desde identidad/contexto del caller |
| MarkSent(id) | MarkSent(tenant,id) | conservar identidad del registro y autorización |
| ID global | (tenant,ID) | mismo tenant conserva deduplicación |

Las cinco funciones de test originales conservan exactamente sus aserciones;
se adaptaron únicamente los argumentos de esas llamadas y gofmt. Un montaje
separado con los tests/callers originales rechaza ambas firmas al compilar:
ese negativo es esperado y demuestra que no se acepta una llamada sin scope.
No se presenta como compatibilidad binaria/fuente con0.1.x. No se hallaron
callers/imports de reminders en otros packs; sí una referencia de metadata de
waitlist. GO-WAITLIST-CORE0.1.1 retira esa compatibilidad; sus2fuentes no cambiaron
y no se implementó ni admitió un adapter waitlist→reminders nuevo.

Due sigue siendo lectura: dos consultas pueden devolver el mismo pendiente.
MarkSent cambia un flag local de manera idempotente; no envía ni valida receipt.
Schedule conserva el Sent aportado por el caller. Ninguno demuestra entrega
exactamente una vez. Se retira esa promesa y se remite a los owners PostgreSQL,
worker y fence outbound de V293 para persistencia y efectos. No hay cancelación,
reagenda, lease/claim, delivery receipt, integración de canales ni recovery durable.

## Evidencia ejecutada

Reconstrucción2/2fuentes byte-idénticas al candidato.11tests ordinarios x3PASS,
vet/build PASS. Se comprueban frontera at==now, rechazo sin consumir identidad,
copias sin alias, ack idempotente, independencia de tenants y64goroutines con
schedule/ack/due/String. No se ejecutó -race ni una suite de carga de target.
GO_NATIVE_FUZZ_GATE:6semillas,10s,4workers, 2112907ejecuciones PASS.
El oráculo usa una lista y busca campos exactos, distinto del map de producción;
tras cada una de hasta32operaciones comprueba todas las vistas tenant, errores,
tiempos y flags. Se usa reloj fijo/avances controlados, sin depender de red ni
reloj real para la secuencia. Fuzz es acotado, no prueba exhaustiva ni DAST.

Toolchains existentes: Go1.26.7 Windows/amd64, Python3.14.4, PowerShell7.6.5;
GOTOOLCHAIN=local, GOWORK/GOPROXY/GOSUMDB=off, CGO_ENABLED=0,
GOMAXPROCS=4 y GOFLAGS=-parallel=4 para calificación. Se conserva la limitación
V35424workers; no se repite ni se declara reparado ese deadline del coordinador.
Método: SOFTWARE_BACKEND_API_ENGINEERING, ENGINEERING_EXECUTION_PLAYBOOK,
PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD y GO_NATIVE_FUZZ_GATE_PACK_PLAN.
Código AUTHORED/LicenseRef-Workspace-Owner; no fuente externa adquirida/adoptada.

## Assurance y disposición

| Dimensión | Demostrado / pendiente |
|---|---|
| Corrección | red/green, reloj controlado, modelo y reconstrucción |
| Integración | migración compile negativa/positiva; sin caller de canales ni waitlist integrado |
| Seguridad/privacidad | reads/acks scoped; autorización y seguridad integrada siguen pendientes |
| Resiliencia | rechazo/replay local; sin restart durable o provider ambiguo |
| Rendimiento | concurrencia acotada; sin benchmark del target |
| Operación | sólo registro en memoria; no servicio/retención/alertas |
| Recovery | fuente anterior preservada; rollback0.1 reabre aislamiento inseguro |
| Release | identidad de fuentes; sin ZIP final/SCA/firma nuevos |
| Usabilidad | migración API documentada y obligada; sin interfaz de recordatorios |

Ambos owners conservan CONDITIONED/SUPPORTED_REFERENCE/REBUILD_VERIFIED;
ninguno aparece en pack plans.160otros packs intactos, franquicia67/746 igual.
TEST02/03 reciben evidencia parcial,43/48controles sin cambios de criterio/estado.
No admisión completa de21cores ni porcentaje de esfuerzo. TEST05/07/09 y demás
owners conservan bloqueos. ZIP148 no contiene este cambio de API ni los fixes
V353/V354; TEST06/08 permanecen acotados a ese candidato y TEST07 exige nuevo
payload/gates. Preflight155 pendiente antes del cierre; sin publicación.

## Recibos

| Archivo relativo al stage | SHA-256 |
|---|---|
| baseline154.json | `0234a8958b07e10924396dbd51c38b77e96a659696ba16694d8b24a8437c3f3e` |
| before155/implementation_packs/GO_REMINDERS_CORE.md | `efe16e3d13cfa7d109d2783d665f33b9e61d40c861a13d186beaef07e70dc0cb` |
| before155/implementation_packs/GO_WAITLIST_CORE.md | `b8b4b5768f78181bfa9a4a91d6900b5d9b1e3ec9ce713f5a58d87d34c333d6a5` |
| baseline-tests.log | `21d97e3804589d7c8e1d34ad12af53e1b97ec3c5ad2c2d3731b5f6d35a8f7a1c` |
| red-tests.log | `d3e714ce1fb1e0e9d054869d68966cd97e14ea9f6545e8233dcc4972fe76d1e5` |
| red-summary.json | `87e9cacaa76a33d5540e5b8e56181b695a987ac20703d2737d5ad8a591bb1eca` |
| red355.py | `6481175a3996cc7b881649ef96ce96cf370eeccfa42b5b239f74522335c7b3fb` |
| fix355.py | `45e2ff230925339fa9f734fb932aa961afb26aa4a5cd0cca6dc795ad36d64995` |
| rebuild355.py | `bea9075e04e311b317cb5e58ccbdb093af0cf8fe041b70f5337ae1d86ff8cb77` |
| green-tests.log | `7b32a040a3c5207a74595d5ac53ab5f5fc2c7199cbcf2b00b05aa29a48801cd3` |
| legacy-api-rejection.log | `80f8e7985b088d1466a588fe05fd664190084d8a7a611fe1c68563818b4dc399` |
| adapted-baseline-tests.go.txt | `ef72e724834ebfe9434487c1c20d1a9e7ce5c2b4beec009e13315311689486f0` |
| rebuilt-tests.log | `8c69b0a09c3689060230a36d16f6fa8777416ac427544a2f5d619a4ec64fcf40` |
| rebuilt-vet.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| rebuilt-build.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| fuzz-profile.json | `497cde11ecef6e63322930650602149610368ae432a20a13552fa5b4044b32e0` |
| fuzz-receipt.json | `8a278d906add0acb188f68cef731bb6585cc68b5590fe86267f77f6ebb840f69` |
| fuzz.log | `19fad84018144cddfa6dd47e9481b9acc2c831d96fb16770f83340e2b8855985` |
| fuzz-self.log | `08ae8b9a55816b44b64fe3a9aec6c663ec42397e928c23a3bd1f8d4698192a4c` |
| qualification155.json | `64525e2fd357a39b435443ec654bdd85611bed9525726be1a028829f176f5c62` |
| ledger155.json | `01e359ae3194bd05fdb00fd2cc40606296c942ee6478385c6352688a4a3528bc` |

### Probes originales de aislamiento (sólo calificación)

````go
package reminders
import("errors";"testing";"time")
// Qualification-only adapters expose the old absence of scope while allowing
// these assertions to run unchanged against the corrected explicit API.
func dueFor(s *Scheduler, tenant string) []Reminder {
 if scoped,ok:=any(s).(interface{Due(string) []Reminder});ok{return scoped.Due(tenant)}
 return any(s).(interface{Due() []Reminder}).Due()
}
func markFor(s *Scheduler, tenant,id string) error {
 if scoped,ok:=any(s).(interface{MarkSent(string,string) error});ok{return scoped.MarkSent(tenant,id)}
 return any(s).(interface{MarkSent(string) error}).MarkSent(id)
}
func reminderFixture(t *testing.T)(*Scheduler,*time.Time){
 t.Helper();now:=time.Unix(1000000,0).UTC();s:=NewScheduler();s.clock=func()time.Time{return now};return s,&now
}
func TestSameReminderIDIndependentTenants(t *testing.T){
 s,now:=reminderFixture(t)
 for _,tenant:=range []string{"a","b"}{if e:=s.Schedule(Reminder{TenantID:tenant,ID:"shared",AppointmentID:"appt",At:now.Add(time.Minute),Channel:ChannelEmail});e!=nil{t.Fatalf("tenant %s blocked by another: %v",tenant,e)}}
}
func TestDueNeverReturnsOtherTenant(t *testing.T){
 s,now:=reminderFixture(t)
 for _,tenant:=range []string{"a","b"}{if e:=s.Schedule(Reminder{TenantID:tenant,ID:tenant,AppointmentID:"appt-"+tenant,At:now.Add(time.Minute),Channel:ChannelEmail});e!=nil{t.Fatal(e)}}
 *now=now.Add(time.Minute)
 for _,tenant:=range []string{"a","b","unknown",""}{for _,r:=range dueFor(s,tenant){if r.TenantID!=tenant{t.Fatalf("requested %q got tenant %q",tenant,r.TenantID)}}}
}
func TestOtherTenantCannotMarkReminder(t *testing.T){
 s,now:=reminderFixture(t)
 if e:=s.Schedule(Reminder{TenantID:"owner",ID:"private",AppointmentID:"appt",At:now.Add(time.Minute),Channel:ChannelPush});e!=nil{t.Fatal(e)}
 if e:=markFor(s,"other","private");!errors.Is(e,ErrNotFound){t.Fatalf("foreign mark accepted: %v",e)}
 *now=now.Add(time.Minute);if got:=dueFor(s,"owner");len(got)!=1{t.Fatal("foreign mark changed owner state")}
}

````

## Cierre156 — verificación de biblioteca

Preflight155: 154controles ejecutados PASS,162packs,
1461archivos materializables,793Markdown,53perfiles. Docker continúa ausente;
gates live/target omitidos o condicionados no se cuentan como ejecutados.
Los48tests originales y requisitos conservan criterios/estado:43PASS/5BLOCKED.
EVID-65 agrega prueba parcial de corrección/aislamiento, sin rebajar TEST02/03.

Reminders0.2.0:2fuentes reconstruidas,11tests x3,vet/build y2112907
ejecuciones fuzz PASS. Las cinco funciones originales conservaron aserciones;
sólo se añadieron argumentos tenant. Los callers originales sin tenant fallan
al compilar deliberadamente. Waitlist0.1.1 retira metadata de compatibilidad y
sus dos fuentes siguen intactas. No se adquirió código externo ni hubo envíos.

La migración y sus límites quedan canónicos. Ambos cores permanecen CONDITIONED
y fuera del perfil. TEST02/03/05/07/09 siguen pendientes; roadmap1completo/10abiertos.
El ZIP148 no representa este payload; ningún release nuevo se firma o entrega.
Checkpoint156 actualiza hashes y preserva cadena append-only; su validación se
ejecuta después de este informe y se conserva aparte sin referencia circular.

| Recibo de cierre relativo al stage | SHA-256 |
|---|---|
| preflight155.json | `3c83b3c80bc6ad613259b4a2e02ca0ef57060ef0627584b563f93de85311d4b5` |
| preflight155.log | `19d6a701cf279af03302b744727d815fae0baa55f19039680e6c63be4302438b` |
| closure155.json | `a9c74ee5f6c02fdfdc5de084b67bf044ac41c394ecdb2f25e3c7ba70995575f5` |
| resume155.log | `14b1b25d7af1440a5bb7ab185748a749a29fa2a25ef0f5c3251b9bea4c8cd7f2` |
| plan155.log | `c912fb6cdf9a56e15785aad245756d19de1ae3116be72c57e00780542982c232` |
