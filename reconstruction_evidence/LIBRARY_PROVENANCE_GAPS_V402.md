# V402 — procedencia y camino de sustitución; auditoría de sólo lectura

Fecha: 2026-09-11. Alcance: mantenimiento de biblioteca, T2802/TEST02 y criterio 9 del nuevo HECHO. No se modifica canónico, no se ejecuta producto, no se cambia admisión ni se reduce el alcance. ARCA y Daybreak no se investigan en esta tarea.

## Conclusión que cambia la ejecución

El perfil de referencia selecciona **70 packs / 820 archivos: 717 AUTHORED, 99 ADAPTED y 4 VERBATIM**. Es un conteo de bloques materializables seleccionados, no de dependencias transitivas ni porcentaje de calidad. Diez archivos WSAA del pack completo no se seleccionan; contar todos los bloques de los setenta packs produciría incorrectamente 830.

El criterio nuevo «AUTHORED sólo glue inevitable» no está satisfecho. Hay negocio sustantivo local en precios, conversión de cotización, contabilidad, regalías y encuestas. Una etiqueta CONDITIONED por cuentas no elimina esa incompatibilidad. Tests o manuales oficiales tampoco convierten una implementación previa en fuente derivada.

Se pueden preservar los owners y contratos. El camino es sustituir decisiones/algoritmos por derivaciones verificadas, separar el transporte/persistencia/identidad propio y conservar deltas semánticos explícitos; no crear dominios paralelos ni reetiquetar archivos sin cambios/evidencia. La licencia/admisión de la plataforma BC completa no se hereda por una adaptación estrecha.

## Cinco contradicciones determinantes

| Frontera | Evidencia actual | Por qué excede glue | Camino mínimo real |
|---|---|---|---|
| Precio y colocación de pedido | GO_COMMERCE_PRICING_PAYMENT_API.md:21/84/547/625/674; 8 AUTHORED | Selección de price book, precio de línea, total permitido y transición de orden son reglas propias. Go/PG/pgx no aportan un motor comercial. | Crear expediente de derivación estrecha BCApps para selección de precio y cálculo de importe; extraer algoritmo probado y delegar desde owner Commerce. Mantener contrato de libro explícito, exactitud monetaria y políticas ajenas rechazadas. Conectar SDK oficial de pago ya existente; el SDK no sustituye precio/pedido. |
| Contabilidad y reversión | GO_ENTERPRISE_ACCOUNTING_LEDGER_API.md:21/47/274/391/440/490/986; 9 AUTHORED | Balance, posting, reversal, período y ledger son un dominio; los ocho hashes BCApps citados son autoridad documental, sin derivación ejecutable declarada. | Derivación trazada GenJnlPostBatch/GenJnlPostLine/GenJnlPostReverse y tests oficiales. Separar representación de importes, transacción PG y scopes locales. No portar impuestos, redondeo fiscal, FX ni cierre estatutario. |
| Cotización → pedido y entrega | GO_FRANCHISE_CUSTOMER_JOURNEY_API.md:171/495/1030/2385; 33 AUTHORED | Aceptación, expiración, conversión y handover son estado de negocio; no basta conectarlos por HTTP. | Traducir estrechamente SalesQuotetoOrder manteniendo autorización del cliente y aceptación inmutable del contrato actual como delta propio explícito. Hay que investigar source para handover inicial/liberación; la fuente Sales shipment ya admitida termina antes de aceptación comercial. |
| Los 21 cores y encuestas V400 | CORE_CLAIM_ADMISSION_V374.md:25/59/122; V400:3/46; GO_CUSTOMER_SURVEY_API.md:59/258/334/390 | Los 21 son helpers propios, 0 equivalencias empresariales; el owner nuevo de encuesta es negocio AUTHORED, incluso con fórmula Bain y pruebas reales. | Reutilizar auditoría sin repetir tests. Para cada función REQUIRED: admitir fuente real y sustituir; para no requerida conservar exclusión explícita ya autorizada, sin cambiar REQUIRED para fabricar PASS. No rellenar waitlist/gift/loyalty/payroll/reviews/social con otro core AUTHORED. |
| Conteos/documentación y promoción | FRANCHISE_COMPLETE_PACK_PLAN.md:3 dice68 pero JSON selecciona70; texto de migraciones termina0048/0049 aunque hay0054; V374 cierre203 se retracta204 | Un perfil materializable o verificador estructural PASS no demuestra el nuevo criterio de procedencia ni la integración comercial. | Actualizar notas y alcance por esta revisión, promover sólo gates cuya evidencia exacta exista. TEST02 puede separar live diferido por instrucción nueva, pero no ignorar source/glue ni huecos locales reales. |

## Fuentes reales ya fijadas: camino Commerce / Accounting

Se leyeron individualmente 13 archivos AL oficiales del **mismo commit ya fijado** `microsoft/BCApps@2eae56d704a1fd035d104f333602aea7091b7749`, tree histórico `f9846fb1254c6311985c6131bb2af11c7f169e1b`. Cada respuesta raw se validó contra el Git blob devuelto por la API de contenidos para ese ref; se preservan bytes y receipts en `official-source-inspection/`. No se adquirió un grupo/plataforma, no se ejecutó AL y no se declara PASS de suites Microsoft. El snapshot histórico tiene licencia raíz MIT SHA `c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383`; confirmar obligations por los archivos derivados en su gate.

| Candidato de derivación | Fuente y funciones concretas | Prueba de origen aprovechable | Límite y trabajo imprescindible |
|---|---|---|---|
| COMMERCE_PRICE_SELECTION_BC | Pricing/Calculation/PriceCalculationV16.Codeunit.al:244 PickBestLine,265 IsDegradedLine,279 IsImprovedLine,293 IsBetterLine,307 IsBetterPrice; PriceCalculationBufferMgt.Codeunit.al:258 SetFiltersOnPriceListLine y278 VerifySelectedLine | TestPriceCalculationV16:1278 T060,1304 T061,1330 T062–T067 para orden/mejor precio/currency/variant;2434 T171 mínimo;2661 T176 fecha;2699 T177 UOM;2747 T178 moneda | Un único libro/moneda/variante no es motor BC completo. Portar función real con licencia y source map, tests derivados y casos negativos; adaptar contrato actual sin activar descuentos, FX, UOM/pricing nuevas silenciosamente. Fechas BC son Date con ending inclusivo; biblioteca usa timestamps con valid_until exclusivo. Ese delta exige boundary mapping explícito y pruebas; no llamar idénticos. |
| COMMERCE_AMOUNT_BC | PriceCalculationV16:307 calcula precio por factor de descuento antes de comparar; PriceCalculationBufferMgt:109 precision,126 RoundPrice,131 IsInMinQty | T060–T067 y escenarios de precisión/UOM que el futuro gate seleccione; inspección de cada oráculo antes de portar | El owner actual limita int64; AL usa Decimal. Mantener enteros/racional exacto, overflow antes de multiplicar y política explícita para redondeo. No ampliar descuentos si no seleccionados. El cálculo de total de pedido requiere además SalesLine, identificado históricamente V172 con SHA `ca65615dc05bc555a878ba2635c955897930332ef171e313a00103074fbc7278`, no leído completo en este informe. |
| QUOTE_TO_ORDER_BC | Sales/Document/SalesQuotetoOrder.Codeunit.al:36 tipo Quote;41 restricciones,43–49 cliente bloqueado;58 líneas bloqueadas;62 cabecera,64 transferencia de líneas;130 CreateSalesHeader; TransferQuoteToOrderLines posterior | Fuentes ERM-Sales históricas V172/V173 están fijadas; deben seleccionar específicamente escenarios quote→order y extraer fixture/assertions (pendiente, no PASS) | BC elimina/archiva quote tras convertir; actual biblioteca conserva aceptación durable y replay. Declarar adaptación de persistencia/consentimiento, no borrar cotizaciones para imitar upstream. El source no es autoridad de handover/liberación financiera. |
| ACCOUNTING_BALANCE_POST_BC | Finance/GeneralLedger/Posting/GenJnlPostBatch.Codeunit.al:373 acumuladores,431–437 suma del balance,477/546 CheckBalance; GenJnlPostLine.Codeunit.al es byte-idéntico al source ya fijado en V121 (`fdecc9...`) | ERMJournalPosting:30/49/69 cantidad conservada/negativos;390 journal sin ForceBalance. Complementar con source-test específico de balance y períodos antes de atribuir ese claim. | Traducir unidad moneda contable a signo exacto, comprobar sumas/rechazo, conservación de cuenta/período y una sola transición durable. PostgreSQL CAS/transacción/outbox sigue glue propio. No cambiar chart, currency, fiscalidad o política de cierre. |
| ACCOUNTING_REVERSAL_BC | Finance/GeneralLedger/Reversal/GenJnlPostReverse.Codeunit.al:209 ReverseGLEntry,217 doble reversión,225 signo invertido,246 vínculos de origen; ReversalEntry.Table.al implementa checks previos | ERMReverseGLEntries:58 out-of-balance;98 blocked account;234/247 y279/291 fuera de período. ERMReverseGLEntriesII conservado para selección de escenarios | **Diferencia material:** BC crea débitos/créditos negativos y flags Correction; biblioteca intercambia columnas positivas. Mapear y probar el signed balance `debit-credit` y vínculos; no afirmar igualdad del turnover bruto. BC soporta reversión de reversión; biblioteca la rechaza. Conservar rechazo como delta explícito y no fingir equivalencia general. |

Para cada fila: materializar CAPABILITY_GAP_RESOLUTION_GATE con el claim acotado y estos refs; G0 identidad/source; G1 licencia; G2 mapeo de semántica y deltas; G3 ownership/interfaces; G4 traducción y tests de fuente; G5 tenant/permisos y rechazo; G6 precisión/concurrencia/cancelación acotada; G7 rollback/recovery; G8 pack/rebuild/notices/evidencia y consumers exactos. **Este informe entrega candidatos y bytes; no cambia RESEARCH_INCOMPLETE a USE_REUSABLE_PACK.** No basta insertar una función derivada nunca llamada ni reclasificar todo el archivo mixto.

## Reutilización que sí ahorra trabajo

- V374: comparación de hashes de 21 packs y 42 bloques exactos **21/21 y42/42**. Se reutilizan 228 tests y decisiones históricas dentro de su límite; no se reejecutaron suites/fuzz. `core-hash-reuse.json` contiene cada comprobación.
- Supply0.16.0 ya contiene63 ADAPTED derivados de BCApps MIT para reserva, UOM, bins, coste/transfer y customer shipment; Fulfillment0.4.0 tiene5 ADAPTED de BCApps/Amazon. Reutilizar estos owners; no construir otro stock.
- GO_OFFICIAL_PAYMENT_WEBHOOK_ADAPTERS tiene SDKs Stripe/MP exactos, autenticación/verifier y create PaymentIntent/Payment; hoy va en perfil separado, no entre los70 del integral. La normalización/invocación SDK es glue admisible, pero reconciliación/booking y wiring al owner Commerce deben completarse y demostrarse. No usar un cambio de etiqueta para sustituir ese trabajo.
- GO_APP_WIRING:45/149 es ejemplo de AUTHORED inevitable (config e inyección de owners), no de dominio empresarial. HTTP/BFF, maps de datos, harnesses, fixture generation y scripts pueden ser glue cuando sólo componen componentes declarados y no deciden precio, crédito, reglas laborales o movimientos financieros.
- Los hashes y diffs previos de V400 se preservan. Su calidad funcional no queda negada por el nuevo criterio, pero su procedencia no lo cumple automáticamente.

## Qué no puede cerrarse por reencuadre local

Credenciales y aceptación productiva se separan sin problema por instrucción del usuario. No ocurre lo mismo con source no admitido, política irreversible no definida o código local ausente: no son sólo «espera cuenta». BCApps plataforma completa sigue CONDITIONAL_PLATFORM, runtime/licencia no demostrados; la demo Google microservices-demo es SAMPLE_ONLY. Ninguno puede sustituir silenciosamente el backend Go entero. Las funciones como nómina necesitan país/política además de cuenta; infraestructura genérica debe aceptar contratos configurados y rechazar ausencia, nunca seleccionar legislación.

## Reproducción de esta auditoría

`write-provenance-audit.py` lee JSON del plan y metadata de bloques seleccionados; no materializa ni ejecuta producto. `provenance-selected-inventory.json` contiene paths, versiones, hashes de cada pack y de cada bloque. Los receipts de fuentes contienen URL fija, SHA-256 y Git blob. No usa un resultado de búsqueda ni reputación como aprobación.

## Identidad de entradas actuales

| Archivo | SHA-256 |
|---|---|
| `markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md` | `8bb5eaac51e36e0a6ce8811415f0c61f4e1a744d76db6f40ada7b3cb0f91ccb2` |
| `markdown_system/FRANCHISE_GAP_MAP.md` | `f4b09fa380b0f1f4f780e1618d5ce1a510a55489340a8c1277ad6af75777cc1a` |
| `reconstruction_evidence/CORE_OWNER_RECONCILIATION_V293.md` | `b892a4f63b1decdba3d47727e6c7c55cea59c5e400e4b8ecb1236e58baeb1f05` |
| `reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md` | `26eca8f36241bd3bc2197b93128315fdb03ca495764201ca5a7638cca04be271` |
| `reconstruction_evidence/CONNECTED_CUSTOMER_SURVEYS_V400.md` | `29970739d1d1e8d77a22c12619677f4fd0dcb491a6149b2c07ddcb394c4f2ca9` |
| `implementation_packs/GO_COMMERCE_PRICING_PAYMENT_API.md` | `0e7b07bba6c4046ee369f52f2a2964990b54b4510645e29a8e40259edc4c7ed3` |
| `implementation_packs/GO_ENTERPRISE_ACCOUNTING_LEDGER_API.md` | `d6a788b518904c56b6c2c0a35c745827964f9d3caafa9567719784dec36b9daa` |
| `implementation_packs/GO_FRANCHISE_CUSTOMER_JOURNEY_API.md` | `9f64509e4ca6d658c16800b6923cc55af947f4eb2bfe40c75a234c508b458fea` |
| `implementation_packs/GO_CUSTOMER_SURVEY_API.md` | `bbc64649680cd2d8d50bd48c597ffdb3741b21ee59f9c1cac83ed3ef54bcb891` |
| `implementation_packs/GO_SUPPLY_FACTORY_INVENTORY_API.md` | `0aa8312e7d5026ee7177849530fcfa594c9642b82913c7ad1586060d3a176a5d` |
| `implementation_packs/GO_OFFICIAL_PAYMENT_WEBHOOK_ADAPTERS.md` | `8faafb118d1bf92445c6a4f0464f3e9678608c6c88a68f1fca0e69a1b09ac99a` |

## Fuentes oficiales leídas — identidad por archivo

| Archivo | SHA-256 |
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
