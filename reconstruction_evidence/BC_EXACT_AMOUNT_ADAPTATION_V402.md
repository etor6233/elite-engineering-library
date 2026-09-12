# V402 — adaptación exacta BCApps de importes y balance

Estado: PASS del adapter estrecho de biblioteca. No cierra TEST02/03/07, no declara infraestructura integral ni producción y no reclasifica el resto del negocio AUTHORED. Se preserva la auditoría V374 por identidad21/21packs y42/42bloques, sin repetirla.

## Cambio verificado

GO-BC-EXACT-AMOUNT-ADAPTER0.1.0 aporta6archivos:2ADAPTED (aritmética y expresiónSQL de reversión),3AUTHORED (tests/documentación de integración) y1VERBATIM (MIT Microsoft). Commerce0.6.2 yAccounting0.1.1 llaman al adapter; sus bloques permanecen AUTHORED y sus admisiones CONDITIONED. Sólo el adapter matemático delimitado queda REUSABLE_PACK; no es Business Central ni código Go escrito/aprobado por Microsoft.

La composición canónica de referencia71/828 coincidió con los828SHA de V400 antes del overlay. Cambiaron3fuentes de owners y se agregaron6archivos. La corrección se limita a importe de línea, suma/balance y traducción de signo de reversión; no incorpora descuento comercial, FX, fiscalidad, nuevos estados ni proveedor de pago.

## Defecto y evidencia

FAIL808: el CreateJournal original admitía pares débito10/crédito-5 y débito-5/crédito10, retornando5/5; admitía también totales MaxInt64+MaxInt64+3 por lado como1/1 por overflow. Dos regresiones RED sobre fuente canónica exacta demuestran que esos valores llegaban al repositorio fake. No se afirma que PostgreSQL los persistiera: sus constraints son otra frontera. El servicio corregido los rechaza antes de persistencia.

- Unitarios de negocio y adapter PASS:5tests principales,8vectores de importe,200balances más casos negativos/representación.
- Gate nativo Go materializado:7semillas,573594ejecuciones,3segundos,2workers PASS. Oracle de overflow por división independiente del big.Int usado por la traducción; sin inputs fallidos.
- PostgreSQL18.6 real en loopback descartable:54migraciones, TestAccountingPostingConcurrencyReversalAndClose y TestCommercePriceOrderAllocationPaymentFlow PASS sin skips. Cobertura incluye posting concurrente, reversión/close y precio→pedido→stock→request de pago.
- Vet y build del módulo Go compuesto completo PASS. PostgreSQL detenido al finalizar.
- Materialización de los3packs:23/23bloques SHA exactos; adapter nuevo6/6. No cambian go.mod/go.sum, ningún pin, migración o runtime.

## Derivación y deltas que no se ocultan

| Fuente fijada BCApps | Función/línea | Destino | Delta explícito |
|---|---|---|---|
| Sales/Document/SalesLine.Table.al | UpdateAmounts:5883 | internal/bcamounts/amounts.go:LineAmount, llamado desde Commerce.AddOrderLine | Cantidad integral, unidadmenor y precisión1: Round es identidad. No VAT/FX; Commerce pasa descuento0. |
| Finance/GeneralLedger/Posting/GenJnlPostBatch.Codeunit.al | ProcessBalanceOfLines:431;CheckBalance:546 | JournalTotals yCheckBalance; CreateJournal/PostJournal | AL Decimal→intermedio exacto big.Int→int64 comprobado; columnas no negativas y turnover positivo son restricciones locales existentes. |
| Finance/GeneralLedger/Reversal/GenJnlPostReverse.Codeunit.al | ReverseGLEntry:209;negación:225 | internal/platform/postgres/bc_accounting_sql.go, llamado dentro de ReverseJournal | Negar debit-credit y normalizar positivo preserva representación local. No equivale a turnover de columnas negativas Correction de BC ni habilita reversión de reversión. |

No hay fuente/admisión nueva para handover inicial, liberación comercial, campañas, payroll, giftcards, loyalty o waitlist. Cotización/replay y todo negocio no sustituido siguen declarados como pendientes bajo criterio9. Una referencia oficial de algoritmo no atribuye el resto del archivo al proveedor.

## G0–G8 del claim estrecho

- **G0 PASS:** Three source functions were selected from fourteen inspected AL files; official fixed commit, Git-blob and SHA256 comparison; no supplier Go attribution.
- **G1 PASS:** Exact Microsoft MIT license hash preserved verbatim; two translated files ADAPTED, local tests/docs and caller glue remain AUTHORED; no new dependency license.
- **G2 PASS:** SalesLine.UpdateAmounts, GenJnlPostBatch balance and GenJnlPostReverse signed negation mapped to exact destination functions. Integral minor units, positive columns and one currency/journal are explicit deltas.
- **G3 PASS:** Existing Commerce and Accounting callers invoke the extracted implementation. No parallel owner, schema, monetary policy or external side effect.
- **G4 PASS:** FAIL808 red/green; exact-amount/200-balance fixtures; fresh PostgreSQL54 migrations and connected pricing/order/stock/payment-request plus concurrent accounting posting/reversal/close; vet/build.
- **G5 PASS:** No new module dependency, network, dynamic execution, secret or PII. Existing tenant/org/period transaction guards unchanged; overflow and negative-sided representation rejected. Exact Go1.26.8 inherited source/runtime identity; broader TEST03 remains separate.
- **G6 PASS:** Bounded int64 operands with exact intermediates, linear journal accumulation and existing SQL transaction batch;7-seed native fuzz573594 executions/3s/2workers; concurrent posting and stock operation tests. No target-load or production-SLO claim.
- **G7 PASS:** Stateless adapter, no migration/durable format change; existing posting/reversal/close and recovery contracts retained, same API and SQL representation. Target operation and all other business owners remain outside this admission.
- **G8 PASS:** Three candidate packs reconstructed23/23blocks; new adapter6/6 match tested code, compatible callers wired and MIT notice included. Source mapping, local evidence, remaining constraints and hash-bound gate record published.

## Identidad de fuentes inspeccionadas

Commit Microsoft fijado `2eae56d704a1fd035d104f333602aea7091b7749`;14archivos individuales raw, cada uno cotejado con Gitblob de laAPI oficial para ese ref. Se seleccionan únicamente3fuentes de implementación anteriores; las restantes informan comparaciones y límites. No se ejecutó AL ni se adquirió runtime/plataforma BC.

| Fuente | SHA-256 |
|---|---|
| [PriceCalculationBufferMgt.Codeunit.al](https://raw.githubusercontent.com/microsoft/BCApps/2eae56d704a1fd035d104f333602aea7091b7749/src/Layers/W1/BaseApp/Pricing/Calculation/PriceCalculationBufferMgt.Codeunit.al) | `6bd6c18dbdea9d0af11f132c4c5e50b975c8c2053838c0fe43fc668f1998297b` |
| [PriceCalculationV16.Codeunit.al](https://raw.githubusercontent.com/microsoft/BCApps/2eae56d704a1fd035d104f333602aea7091b7749/src/Layers/W1/BaseApp/Pricing/Calculation/PriceCalculationV16.Codeunit.al) | `72b28468771f5413bc3a1aeb1ed2b53886a3b8655e17fefc752c99a54e054b44` |
| [PriceListLine.Table.al](https://raw.githubusercontent.com/microsoft/BCApps/2eae56d704a1fd035d104f333602aea7091b7749/src/Layers/W1/BaseApp/Pricing/PriceList/PriceListLine.Table.al) | `8d0b4349ab5e02a0cd445374bc727582bc08ea31e840dcf667f70616c1737e64` |
| [SalesQuotetoOrder.Codeunit.al](https://raw.githubusercontent.com/microsoft/BCApps/2eae56d704a1fd035d104f333602aea7091b7749/src/Layers/W1/BaseApp/Sales/Document/SalesQuotetoOrder.Codeunit.al) | `5ba30dbcb3d5bb9fb250df6691ad7fc461e1e023665d1572652b3561a736270d` |
| [GenJnlPostLine.Codeunit.al](https://raw.githubusercontent.com/microsoft/BCApps/2eae56d704a1fd035d104f333602aea7091b7749/src/Layers/W1/BaseApp/Finance/GeneralLedger/Posting/GenJnlPostLine.Codeunit.al) | `fdecc9a5e52831552e9addfe4c7e66d1b53fa3425c95b2d49d606478a24d54f3` |
| [GenJnlPostReverse.Codeunit.al](https://raw.githubusercontent.com/microsoft/BCApps/2eae56d704a1fd035d104f333602aea7091b7749/src/Layers/W1/BaseApp/Finance/GeneralLedger/Reversal/GenJnlPostReverse.Codeunit.al) | `97d9268c6be04de015c50595e94b935710cd62610f987a08d66c9aab6790cba5` |
| [ReversalEntry.Table.al](https://raw.githubusercontent.com/microsoft/BCApps/2eae56d704a1fd035d104f333602aea7091b7749/src/Layers/W1/BaseApp/Finance/GeneralLedger/Reversal/ReversalEntry.Table.al) | `6cbaae248437832efd3e9b9e683640d436abd06d1aae3e38c1d052db2fcca75b` |
| [PriceListLineUT.Codeunit.al](https://raw.githubusercontent.com/microsoft/BCApps/2eae56d704a1fd035d104f333602aea7091b7749/src/Layers/W1/Tests/ERM/PriceListLineUT.Codeunit.al) | `937be20d961932c8ac45f29c921467bfb4d55d6151737bec96ea56c7d23672ba` |
| [TestPriceCalculationV16.Codeunit.al](https://raw.githubusercontent.com/microsoft/BCApps/2eae56d704a1fd035d104f333602aea7091b7749/src/Layers/W1/Tests/ERM/TestPriceCalculationV16.Codeunit.al) | `4b199f235bb527c8853157187bfa5fa35c72a7816ced0833f0ff3bc043c211bb` |
| [GenJnlPostBatch.Codeunit.al](https://raw.githubusercontent.com/microsoft/BCApps/2eae56d704a1fd035d104f333602aea7091b7749/src/Layers/W1/BaseApp/Finance/GeneralLedger/Posting/GenJnlPostBatch.Codeunit.al) | `be8be19c361c4610f0c3894eb97cfe0554ab54c1de0aaa491a8a3903fee4b150` |
| [ERMReverseGLEntries.Codeunit.al](https://raw.githubusercontent.com/microsoft/BCApps/2eae56d704a1fd035d104f333602aea7091b7749/src/Layers/W1/Tests/Reverse/ERMReverseGLEntries.Codeunit.al) | `6c6690de941bbe88b84879792270b44057fb696da1c868bd15e6dbabc932180d` |
| [ERMReverseGLEntriesII.Codeunit.al](https://raw.githubusercontent.com/microsoft/BCApps/2eae56d704a1fd035d104f333602aea7091b7749/src/Layers/W1/Tests/Reverse/ERMReverseGLEntriesII.Codeunit.al) | `6cbfe2c3ec2a8e1cc52983883d5c0d5c8dd424498f26d58a843b00311e782752` |
| [ERMJournalPosting.Codeunit.al](https://raw.githubusercontent.com/microsoft/BCApps/2eae56d704a1fd035d104f333602aea7091b7749/src/Layers/W1/Tests/General%20Journal/ERMJournalPosting.Codeunit.al) | `1c465cfe8342e7c28b28e139dba0ff58ab91b27fe8c3227a98c87f2945ec44af` |
| [SalesLine.Table.al](https://raw.githubusercontent.com/microsoft/BCApps/2eae56d704a1fd035d104f333602aea7091b7749/src/Layers/W1/BaseApp/Sales/Document/SalesLine.Table.al) | `ca65615dc05bc555a878ba2635c955897930332ef171e313a00103074fbc7278` |

MIT raíz Microsoft: `c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383`, preservada literal. Go1.26.8 go.exe:`21761eceb9302062c9623fb699f332c8c7fe000f15f70efe8da01a2cfbbc16b9`. Fuentes/límites de método Go se reutilizan de V374; no cambió identidad.

## Receipts locales y fuentes publicadas

Stage fuera de distribución: `Temp/elite-v402-library-infra`. Los hashes ligan bytes concretos; los logs no se presentan como portable runtime.

| Archivo | SHA-256 |
|---|---|
| `business-before-receipts.json` | `5330a43dd4d52745704af61419c0ae86189e08091ea5aa1a9fed9e16334ab0e8` |
| `core-hash-reuse.json` | `ea4e21bda226b61b32851d768d50fe5f8e052109e6338d95f81d7c17c3074158` |
| `business-red.log` | `e02b11a4b90b0f8c01290707d09c4fda4cea59ac49d77c9b7bf720bb874da320` |
| `business-green-unit.log` | `d596acefffc3257b8c5de8b114c3bc1fa72e12a771d48c83b66b60bcaef68444` |
| `business-pg-1/result.json` | `e85795459e554fc67146f89d79d1efd84cf7535677cfa5ef101e87b83f365561` |
| `business-pg-1/connected.log` | `e9b79200fb0f45e1e88543cc6de0a4f326845a221429b2469d11fd371c754b72` |
| `business-pg-1/vet.log` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `business-pg-1/build.log` | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `business-fuzz-profile.json` | `8582dded09a53b419a9a7e17fccbe6aeaf92d51ff9412c29fc67a00c59bd680b` |
| `business-fuzz-receipt.json` | `e6e9bb35544ff131cb67a087cac085af250528813eed3bca1cc947e0e1d42251` |
| `business-fuzz.log` | `d0ea0486b5e4844baff90076813fb03d1f91ebdf4d8abef4bf66538fc8fe87af` |
| `business-roundtrip-receipt.json` | `5bc220d8f46d5d954607f06e67cabfa80dd8d0edd9eb89f79501b1f21ebda699` |
| `official-source-inspection/license-receipt.json` | `8ebbbdd29cb0731148146d279bb29d3b4c105e3dc6aeae7e1c276fd6e44f68c8` |
| `implementation_packs/GO_BC_EXACT_AMOUNT_ADAPTER.md` | `d08b9f5c676286b460340a493902d43024a8e4eac8e83ae41c65d8d02b5cbbd2` |
| `implementation_packs/GO_COMMERCE_PRICING_PAYMENT_API.md` | `1f9b7213e2b73718bc42727a1989fb3bdc2310847d14c351335970301d2f847a` |
| `implementation_packs/GO_ENTERPRISE_ACCOUNTING_LEDGER_API.md` | `f3841975ac2acc9432a667b24559f05538901f18531c280e2dc82e76510eb9f4` |

## Límites de adopción

El futuro proyecto confirma agrupación/currency/políticas e identidad/operación de sus owners, ejecuta gates afectados y conserva esta licencia. El adapter no toma decisiones comerciales, no conoce cuentas, no habilita integración externa y no certifica balances financieros. La política del usuario permite infraestructura local y difiere cuentas; esas cuentas no subsanan los demás dominios AUTHORED pendientes.
