# V354 — aislamiento y límites de gift cards y fidelización

Fecha:2026-09-09. Mantenimiento correctivo/revisión; continuación de TEST02,
T2802 y auditoría V293. Entrada152 validada. No nuevo producto, owner financiero,
capability, dependencia ni autoridad comercial asumida.
Stage: `C:/Users/NL/AppData/Local/Temp/elite-v354-fab8f124d0124c95b4906f8350cac023`. Checkpoint153 registra integración; cierre154 posterior.

## Resultado observado

Los11tests originales pasaban. Cinco tests nuevos fallaron contra0.1.0:
CreditOverflowAtomicAndRetry, AccountTupleIsolation, EntryIdentityScopedByTenant,
RedeemIdentityScopedByTenant y BlankRedeemIdentityRejected (tres subcasos).
Loyalty aceptaba MaxInt64+1, devolviendo saldo negativo y consumiendo ID/historia.
Las cuentas (tenant=a\0b,cliente=c) y (tenant=a,cliente=b\0c) compartían saldo
por concatenación. Los IDs globales en ambos cores permitían que la actividad
de un tenant bloqueara el mismo ID en otro. Giftcards aceptaba IDs vacíos/blancos.
String leía el map sin mutex; ese hallazgo fue estático, no detector de races.

Se corrigieron GO-LOYALTY-CORE y GO-GIFT-CARDS-CORE a0.1.1. AccountID y los IDs de
operación usan campos separados. La identidad de una operación es (tenant,ID),
conservando unicidad entre cuentas/tarjetas del MISMO tenant y entre Earn/Burn.
No se relaja deduplicación: el replay sigue ErrDuplicate/ErrSpent sin otro efecto.
Overflow retorna ErrOverflow antes de consumir ID, balance o historia; un intento
rechazado puede repetirse al existir capacidad. Giftcards valida el ID antes del
débito y sincroniza String. Se preserva API exportada salvo nuevo ErrOverflow.

Se corrigió el texto que prometía top-up/Issue idempotente: Issue rechaza código
duplicado y no recarga. History es copia en memoria y cuesta O(total movimientos),
no un journal durable ni un mecanismo de restore. No se añadieron reglas de
puntos, moneda, reversión, impuestos, expiry, checkout ni permisos de usuario.

## Pruebas y límites

Fuentes canónicas reconstruidas4/4byte-idénticas.19tests ordinarios (loyalty10,
giftcards9) x3PASS; vet/build PASS. Pruebas concurrentes con64goroutines verifican
crédito duplicado único, débito limitado al saldo, historial consistente y copia
aislada; String se ejercita durante emisión concurrente. No se ejecutó -race.

Cada fuzz target tiene6semillas y secuencias finitas de hasta32operaciones loyalty
o48giftcards. El oráculo loyalty usa big.Int y mapas anidados, independientes de
la representación de claves y cálculo de capacidad de producción; contrasta
errores, saldo e historial de todas las cuentas tras cada operación. Giftcards
contrasta saldo/errores/IDs por tenant contra su modelo. Dos corridas completas
con las mismas fuentes/invariantes y presupuesto10s por target:

| Corrida | Loyalty ejecuciones | Giftcards ejecuciones | Workers por target |
|---|---:|---:|---:|
| 1 | 1220787 | 1673312 |4|
| 2 | 1268674 | 1676582 |4|

El primer gate con24workers FALLÓ: loyalty pasó5.529.823ejecuciones; giftcards
terminó7.480.937con context deadline exceeded al finalizar10s. No se guardó
crasher/testdata ni se emitió receipt PASS. FAIL623 conserva este límite. El
origen exacto del deadline no está demostrado; la hipótesis de presión del
coordinador motivó la concurrencia explícita, no un parche al runtime. Sólo las
dos corridas4workers se admiten como calificación completa. No son exhaustivas.

Toolchain fijado Go1.26.7 Windows/amd64, Python3.14.4, PowerShell7.6.5;
GOTOOLCHAIN=local, GOWORK/GOPROXY/GOSUMDB=off, CGO_ENABLED=0; para corridas finales,
GOMAXPROCS=4 y GOFLAGS=-parallel=4. Ninguna versión nueva fue adoptada.
Autoridad: [Go: overflow entero](https://go.dev/ref/spec#Integer_overflow) y
[Go: claves de maps](https://go.dev/ref/spec#Map_types). La especificación actual
consultada muestra1.27; se usa para semántica estable y el runtime probado sigue
1.26.7. Estos métodos son AUTHORED, no código financiero de Go ni de otra empresa.

## Assurance y disposición

| Dimensión | Evidencia acotada / pendiente |
|---|---|
| Corrección | red/green, bordes, modelo independiente, reconstrucción |
| Integración | no montaje con checkout/ledger empresarial; compatibilidad histórica no probada para owners actuales |
| Seguridad/privacidad | aislamiento de claves local; faltan autorización, threat model y SCA integrada |
| Resiliencia | rechazos conservan estado; no crash/restart durable |
| Rendimiento | sin benchmark de target; concurrencia acotada de prueba no es carga productiva |
| Operación | stores en memoria, sin servicio, retención ni monitoring productivo |
| Recovery | fuente anterior preservada; no API restore y rollback reabriría defectos |
| Supply chain/release | hashes de4fuentes; sin nuevo release firmado ni promoción del ZIP148 |
| Usabilidad | errores/clases de API probados; no UX/ayuda de checkout |

Ambos packs permanecen AUTHORED/LicenseRef-Workspace-Owner,
SUPPORTED_REFERENCE / REBUILD_VERIFIED / CONDITIONED.160otros packs intactos;
los dos siguen fuera de todos los pack plans y franquicia67/746 no cambia.
TEST02 recibe evidencia parcial;43/48checks conservan estado/criterios,5bloqueados.
No porcentaje de esfuerzo ni cierre de los21cores. TEST03/05/07/09 y los owners
de V293 conservan sus exigencias. El ZIP148 no contiene estos fixes ni FX0.1.1;
TEST06/08 conservan alcance histórico exacto, TEST07 exige nuevo payload/gates.
Preflight153 pendiente antes del cierre. No publicación ni efectos externos.

## Recibos y reproducción

El materializador rechazó el segundo pack en carpeta no vacía (FAIL620).
Se preservó ese stage; se reconstruyó cada pack separado y se copiaron sólo sus
cuatro fuentes exactas a un módulo de calificación. go.mod no es archivo del pack.
La fuente de pruebas/fuzz corregida está en los bloques canónicos. Los probes
originales abajo se ejecutan contra las revisiones0.1.0 conservadas.

| Archivo relativo al stage | SHA-256 |
|---|---|
| baseline152.json | `88697d83cc45cdecd221bab5e2bbbdb8aa39d87fcde6cad7e6e54687fb711a6f` |
| before153/implementation_packs/GO_LOYALTY_CORE.md | `b2f3e07c4b127be207351b82d77be2a977c708b0aaed9aa200ed2103699cbcc8` |
| before153/implementation_packs/GO_GIFT_CARDS_CORE.md | `0ce8666f978d918c076652ec11a99fb1ba2e0477ac44ab4ce3c96b8b0490ad2d` |
| baseline-tests.log | `6831983281438fc51a2c6c36a66eb5a47ae4bdd8c10be21b5a2057e896658d62` |
| red-tests.log | `59aced5d7184ab28e72685a2a066bda835eecde6065322744aac80b3eadafa51` |
| red-summary.json | `5e2a0af5659f2112e04c55032e555a4ff03b010a25c74f0b56562913f16d686b` |
| red354.py | `d2d6ba0902711c39544f58b68ff073e4ed69fb50fb93219775ad440f3418a685` |
| red354-recovered.py | `63ac204b82e8e643b3806baa5cfaa25418f61e40cc347e0ecc3a65890bca37fa` |
| fix354.py | `2982d80eef66b3fd95eeaf60226161d3d36fa3d10178aef313c4a615c12d3dbc` |
| rebuild354.py | `4138255b50fbc291e41486197ac0e5732d150273f3be596e4956aa3cb57b222c` |
| bounded-fuzz.py | `ed1ff54f2e72ec78c4a1ce2c0ce060abc9b768449d32f773552a741398a1e048` |
| green-tests.log | `e026e010941f6c0c651d1ed757c64a9191c5c0ce82a40d042d8ab696bb70a5f6` |
| rebuilt-tests.log | `ba216829707879c9507adb13974d359666869ec14a96503cce771ee87db82938` |
| rebuilt-vet.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| rebuilt-build.log | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| fuzz.log | `68b26572c1c79a6b4894800523875ea61c7994fedd7d6cd9aa8859253aa6e8cd` |
| fuzz-profile.json | `baffa758f5116ba38e8e7b9cf168908c578700cc6ce49c5df8adade5f9e1c0fd` |
| fuzz-self.log | `e11c09145ce74279fd0abf050882c7dcf444edd4e6e95c2dfe6f234f08bf1148` |
| fuzz-bounded-1.log | `22b26c198aa66f9b32a65ee0425f14e3d80a807e0618c582304c2c0a6277aa1b` |
| fuzz-bounded-1.json | `5693956fb6b85a6c1e46b26856cbbb596ae4e6fe2af901c019f1360b0637e292` |
| fuzz-bounded-2.log | `2714c7aaf6bd9beccdcbdafe74ffe9d3d1cf1241aee2064fe338802884b1b5c4` |
| fuzz-bounded-2.json | `5693956fb6b85a6c1e46b26856cbbb596ae4e6fe2af901c019f1360b0637e292` |
| qualification153.json | `0f39a995611301f5ac8d9384a296a10fa66d28b0ba6c8bf48774e733050c9485` |
| ledger153.json | `2524d0b3a4f55ec952351db9f91b9013a39da3ff842b317ab8f352527eb23339` |

### Probes originales loyalty

````go
package loyalty
import("testing";"math";"errors")
func TestCreditOverflowAtomicAndRetry(t *testing.T) {
 l:=NewLedger();if e:=l.Earn("t","c","max",math.MaxInt64,"initial");e!=nil{t.Fatal(e)}
 if e:=l.Earn("t","c","next",1,"overflow");e==nil{t.Fatal("overflow accepted")}
 if l.Balance("t","c")!=math.MaxInt64 || len(l.History("t","c"))!=1{t.Fatal("rejection mutated balance/history")}
 if e:=l.Burn("t","c","space",1,"release");e!=nil{t.Fatal(e)}
 if e:=l.Earn("t","c","next",1,"retry");e!=nil{t.Fatalf("failed attempt consumed ID: %v",e)}
 if e:=l.Earn("t","c","next",1,"replay");!errors.Is(e,ErrDuplicate){t.Fatalf("duplicate accepted: %v",e)}
}
func TestAccountTupleIsolation(t *testing.T) {
 l:=NewLedger();if e:=l.Earn("a\x00b","c","first",20,"initial");e!=nil{t.Fatal(e)}
 if l.Balance("a","b\x00c")!=0 {t.Fatal("distinct account tuples alias")}
 if e:=l.Burn("a","b\x00c","steal",1,"other");!errors.Is(e,ErrInsufficient){t.Fatalf("other account debit: %v",e)}
 if l.Balance("a\x00b","c")!=20 || len(l.History("a","b\x00c"))!=0 {t.Fatal("cross-account mutation")}
}
func TestEntryIdentityScopedByTenant(t *testing.T) {
 l:=NewLedger();for _,tenant:=range []string{"a","b"}{
 if e:=l.Earn(tenant,"c","shared-earn",10,"initial");e!=nil{t.Fatalf("tenant %s cannot earn: %v",tenant,e)}
 if e:=l.Burn(tenant,"c","shared-burn",1,"burn");e!=nil{t.Fatal(e)}
 }
 if e:=l.Earn("a","other","shared-earn",10,"reuse");!errors.Is(e,ErrDuplicate){t.Fatalf("same-tenant duplicate: %v",e)}
}

````

### Probes originales giftcards

````go
package giftcards
import("testing";"errors")
func TestRedeemIdentityScopedByTenant(t *testing.T){
 s:=NewStore();for _,tenant:=range []string{"a","b"}{
 if e:=s.Issue(tenant,"CARD",100);e!=nil{t.Fatal(e)}
 if e:=s.Redeem(tenant,"CARD","shared",10);e!=nil{t.Fatalf("tenant %s cannot redeem: %v",tenant,e)}
 }
 if e:=s.Issue("a","OTHER",100);e!=nil{t.Fatal(e)}
 if e:=s.Redeem("a","OTHER","shared",10);!errors.Is(e,ErrSpent){t.Fatalf("same-tenant duplicate: %v",e)}
 if b,_:=s.Balance("a","OTHER");b!=100{t.Fatal("duplicate debited other card")}
}
func TestBlankRedeemIdentityRejected(t *testing.T){for _,id:=range []string{""," ","\t\n"}{t.Run(id,func(t *testing.T){
 s:=NewStore();if e:=s.Issue("t","CARD",100);e!=nil{t.Fatal(e)}
 if e:=s.Redeem("t","CARD",id,1);!errors.Is(e,ErrInvalidCard){t.Fatalf("blank identity accepted: %v",e)}
 if b,_:=s.Balance("t","CARD");b!=100{t.Fatal("invalid identity mutated balance")}
})}}

````

## Cierre154 — biblioteca y checkpoint verificados

Preflight153 completó 154controles ejecutados PASS:162packs,
1461archivos materializables,792Markdown y53perfiles. Docker permanece ausente;
gates live/target omitidos o condicionados no se cuentan como ejecutados.
Ambos packs0.1.1 conservan4fuentes reconstruidas y19tests x3PASS.
Las dos corridas completas4workers suman5839355ejecuciones entre los2targets.
El primer deadline24workers sigue preservado y no se suma como gate aprobado.

Se verificó que los48tests y requisitos originales mantienen sus definiciones
y estados:43PASS/5BLOCKED. Sólo se adjuntó EVID-64 a TEST02/03 como evidencia
parcial, sin cambiar el denominador ni eliminar condiciones. TEST02admisión,
TEST03seguridad integrada, TEST05operación, TEST07release y TEST09historial
siguen pendientes;1macro-tarea completa/10abiertas. No se entrega release final.

Checkpoint154 actualiza hashes, vincula los dos owners corregidos y conserva
eventos append-only. Su validación se ejecuta después de escribir este informe;
el recibo posterior queda en el stage para no crear referencias circulares.

| Recibo de cierre relativo al stage | SHA-256 |
|---|---|
| preflight153.json | `d39be9d358d404341366d27a61f10328561a9a35b3284f3e2b81244b8b0b32a5` |
| preflight153.log | `31df5a41f28cef660a38bca059c4d07eb287e58f1eed99eb58887b8ba159c290` |
| closure153.json | `84a3502ac9b7df32cc2f868020e6484b5548786d9df1fe9e684f074176d0d747` |
| resume153.log | `ec050a5e4132bd9fceb9636b467d952d8d5340da6e2057c05f05cd27950f77c0` |
| plan153.log | `c912fb6cdf9a56e15785aad245756d19de1ae3116be72c57e00780542982c232` |
