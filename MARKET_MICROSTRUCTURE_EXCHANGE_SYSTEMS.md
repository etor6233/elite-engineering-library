# Microestructura de mercados y sistemas de exchange — manual operativo

> **Estado:** núcleo general v1 auditado y cruces recíprocos cerrados el 2026-08-21, incluido frontend. Triangulación: Stanford MATH 237 y papers fundacionales/empíricos; Nasdaq TotalView-ITCH/OUCH 5.0; Eurex T7; CME MDP; FIX Trading Community; Binance Spot, Coinbase Exchange/International y Kraken Spot; Uniswap v2/v3; SEC/CFTC. Codecs, rulebooks, regulación y contratos dependientes de venue/producto se fijan por proyecto en §19.

> **Alcance:** comprender, reconstruir, medir y operar sistemas de market data/order lifecycle. No promete predecir precios, rentabilidad, best execution ni cumplimiento regulatorio; esos resultados requieren venue, jurisdicción, cuenta, fees, latencia y evidencia propios.

---

## 0. Cómo usar este manual

Orden obligatorio para un humano o agente:

1. fijar instrumento, venue, sesión, reglas y versión de protocolo;
2. definir qué estado es observable: trades, L1, L2, L3, órdenes propias;
3. implementar decoder + secuencia + snapshot/delta + invariantes;
4. demostrar replay determinista y resync ante gap/corrupción;
5. recién después calcular métricas o conectar order entry;
6. si existe capacidad de actuar, aplicar risk gate independiente y probar kill/recovery.

Etiquetas:

- `[SPEC]`: regla o protocolo oficial;
- `[ACADEMIC]`: modelo/resultado académico con supuestos;
- `[PROD]`: práctica documentada de operación;
- `[MEASURED]`: resultado reproducido en el entorno objetivo;
- `[IMPL]`: decisión de implementación derivada;
- `[OPEN]`: falta venue/versión/evidencia; no tratar como verdad cerrada.

Regla central:

```text
Un feed es una secuencia parcial bajo un contrato.
Un libro local es estado derivado, no la verdad del matching engine.
Una métrica describe una ventana observada, no una predicción universal.
```

## 1. Contrato de venue antes de código

```yaml
venue:
  name: ""
  market_and_jurisdiction: ""
  rulebook_version: ""
instrument:
  symbol_and_ids: []
  base_quote_settlement: []
  tick_lot_notional_filters: {}
  session_calendar_timezone: ""
  corporate_or_contract_lifecycle: ""
matching:
  continuous_or_auction: []
  price_priority: ""
  same_price_allocation: ""
  amend_priority_rules: ""
  hidden_reserve_peg_stp: {}
market_data:
  protocol_schema_version: ""
  l1_l2_l3: ""
  transport_and_channels: []
  sequence_scope: ""
  snapshot_recovery_checksum: {}
order_entry:
  protocol_version: ""
  client_and_venue_ids: []
  ack_reject_fill_cancel_semantics: {}
  rate_position_credit_controls: {}
  cancel_on_disconnect_kill: {}
time:
  venue_event_clock: ""
  receive_capture_process_publish_clocks: []
open_questions: []
```

**[SPEC]** No copiar semántica entre venues. Price-time, pro-rata, auction allocation, amend-in-place, self-trade prevention, trigger orders, price bands, snapshot limits y sequence scopes pueden diferir incluso entre productos de un mismo exchange.

## 2. Primitivas microestructurales

### 2.1 Libro, cotización y ejecución

- **bid/ask:** mejor precio comprador/vendedor visible bajo el feed;
- **spread:** `ask - bid`; puede expresarse en ticks, moneda o relativo al mid;
- **mid:** `(bid + ask) / 2`; referencia, no precio necesariamente ejecutable;
- **depth:** cantidad agregada por nivel o por orden según L2/L3;
- **trade:** ejecución reportada; no equivale a nueva orden ni siempre revela todo el match;
- **aggressor/taker:** lado que cruza liquidez según reglas/datos; inferirlo sólo desde precio puede fallar;
- **maker/resting:** interés que ya descansaba, salvo mecanismos/auctions especiales;
- **queue priority:** posición económica bajo la regla de asignación; no se obtiene fielmente de L2 agregado.

### 2.2 L1, L2 y L3

| Vista | Contiene | Permite | No permite por sí sola |
|---|---|---|---|
| L1 | best bid/ask y tamaños | spread, top imbalance, top depletion | forma completa, orden individual |
| L2 | cantidad por precio | depth curve, walk/slippage visible | prioridad individual, hidden/reserve |
| L3 | eventos/órdenes identificables | reconstruir colas visibles y lifecycle | identidades ocultas, intención, totalidad del venue |

Un snapshot truncado no vuelve conocido el resto del libro. Binance Spot advierte que su snapshot REST llega hasta 5.000 niveles por lado: niveles más profundos permanecen desconocidos hasta que cambian.

### 2.3 Órdenes y restricciones

- **market:** prioriza inmediatez, no precio total; puede recorrer múltiples niveles;
- **limit:** restringe peor precio, no garantiza fill;
- **IOC/FOK/post-only:** cambian resting, partial fill y cancel semantics;
- **stop/trigger/peg:** no asumir que reside o se publica antes de activarse;
- **iceberg/reserve/hidden:** displayed quantity no es total disponible;
- **amend/replace:** puede mantener o perder prioridad según campo/regla;
- **STP:** evita o transforma self-match según modo; puede expirar, cancelar o decrementar lados distintos.

“Market”, “limit” o “cancel” son nombres insuficientes: registrar time-in-force, price protection, trigger source, rounding, fees, priority y estados terminales exactos.

### 2.4 Por qué puede existir el spread

Separar mecanismos que pueden coexistir:

- **adverse selection:** una ejecución puede revelar que la contraparte poseía mejor información;
- **inventory risk:** acumular posición expone al dealer mientras espera descargarla;
- **order processing/capital:** infraestructura, clearing, balance y riesgo operativo no son gratuitos;
- **market power, tick y reglas:** competencia, discretización y diseño del venue también delimitan el spread.

**[ACADEMIC]** En Glosten–Milgrom, traders con información superior bastan para producir un spread positivo aun con un especialista risk-neutral y de beneficio esperado cero. Es un mecanismo teórico de adverse selection, no una descomposición universal de cada spread observado.

**[ACADEMIC]** Kyle modela subastas secuenciales entre un insider estratégico, noise traders y market makers competitivos. Su `lambda` relaciona order flow agregado con ajuste de precio y representa profundidad/liquidez dentro de ese equilibrio. No tratarla como constante física ni estimarla sin declarar ventana, normalización, régimen y supuestos.

**[ACADEMIC]** Avellaneda–Stoikov muestra cómo inventory risk puede desplazar el reservation price y volver asimétricas las cotizaciones óptimas bajo utilidad CARA, precio difusivo y llegadas de órdenes parametrizadas. Sirve como lente/control baseline; en producción debe incorporar tick, fees, queue position, fill uncertainty, jumps, latency, inventory hard limits y calibración walk-forward.

### 2.5 Modelos de colas y flujo

- **Markovian LOB:** Cont–de Larrard representa las colas best bid/ask y sus depleciones para estudiar duración y dirección del siguiente movimiento bajo hipótesis Markovianas.
- **Queue-reactive:** las intensidades de add/cancel/market orders dependen del estado visible del libro; útil para simulación y execution/TCA cuando se valida por instrumento y régimen.
- **Hawkes/autoexcitation:** captura clustering y excitación cruzada entre eventos, pero estabilidad, identificabilidad y causalidad no vienen garantizadas por un buen fit.

Contrato de uso de cualquier modelo:

```text
def model_card(model):
    return {
        "question": "qué magnitud intenta explicar/simular",
        "state_observed": "L1/L2/L3, trades, clocks y cobertura",
        "assumptions": "matching, stationarity, arrivals, hidden liquidity",
        "calibration_window": "instrumento, venue, fechas y régimen",
        "falsification": "hechos que lo invalidan",
        "out_of_sample": "distribuciones, colas, paths e invariantes",
        "operational_limit": "no sustituye estado, riesgo ni venue spec",
    }
```

Un modelo generativo que reproduce spread medio pero falla en queue depletion, inter-arrivals, cancel bursts o resiliency no reproduce el mecanismo. El survey de Gould et al. documenta precisamente que muchos modelos LOB no se parecen suficientemente a los libros reales; ajuste visual o una métrica aislada no basta.

## 3. Matching: mecanismo, no folklore

**[SPEC]** Eurex T7 mantiene price priority y ofrece distintas asignaciones dentro del mejor precio: time, pro-rata y time-pro-rata. CME Globex también usa algoritmos configurados por producto. Por tanto, “todos los libros son FIFO” es falso.

### 3.1 Continuous matching

Contrato mínimo:

1. elegibilidad por estado de mercado, instrumento, order type y restricciones;
2. mejor precio primero;
3. asignación dentro del precio según regla del producto;
4. execution price según resting/incoming/auction rule;
5. residual: resting, cancelado o rechazado según TIF;
6. eventos de execution/order state emitidos en el orden definido.

En FIFO, cantidad anterior al order propio sólo aproxima queue-ahead si se posee L3 correcto y no hay hidden/reserve/priority-changing events. En pro-rata, llegada más temprana no determina por sí sola allocation; size, mínimo, top order o residual rules pueden importar.

### 3.2 Auctions, halts y transiciones

Opening/closing/reopening/volatility auctions agregan interés y determinan un precio bajo objetivos/tie-breakers del venue. No aplicar el algoritmo continuous. El feed debe modelar:

```text
PREOPEN → AUCTION/CROSS → CONTINUOUS
CONTINUOUS → HALTED/PAUSED → QUOTE_ONLY/AUCTION → RESUMED
SESSION_END → CLOSED
```

Imbalance/indicative price no es fill garantizado. Mensajes administrativos, directory, trading action, price bands y session events forman parte del estado; descartarlos por “no ser trades” corrompe interpretación.

### 3.3 La prioridad se versiona por producto y operación

Ejemplos oficiales que impiden abstraer de más:

| Venue/producto | Regla verificada | Consecuencia de implementación |
|---|---|---|
| Coinbase Exchange Spot | continuous first-come/price-time; ejecución al precio de la resting order | arrival al matching engine, no timestamp cliente, fija prioridad; STP del taker decide el cruce propio |
| Binance Spot | cancel-replace pierde time priority; `amend keep priority` sólo reduce cantidad y depende de `amendAllowed` | no mapear ambos endpoints a un `replace()` genérico |
| Eurex T7 | price priority y asignación time/pro-rata/time-pro-rata según producto | cargar regla por instrumento/session, no por nombre del venue |
| CME Globex | algoritmos y porcentajes pueden cambiar por product group | capturar el atributo/configuración vigente y fechar el rulebook |
| Nasdaq TotalView/OUCH | referencia/prioridad puede cambiar por replace o update del sistema | consumir mensajes de priority update; no estimar cola por ID original |

`venue + market + instrument + session + effective_date + protocol_version` identifica el contrato. Un cambio de algoritmo es un schema/config migration: invalida simuladores, estimadores de fill y baselines aunque el wire format no cambie.

## 4. Dos state machines separadas

### 4.1 Estado de market data

```text
DISCONNECTED
  → CONNECTING
  → BUFFERING/SYNCING
  → LIVE
  → STALE | GAPPED | CHECKSUM_FAILED
  → RESYNCING
  → LIVE
```

Sólo `LIVE` publica derivados como actuales. `STALE/GAPPED` conserva evidencia, pero no se presenta como libro válido. Cada transición registra `venue/channel/instrument/epoch/last_sequence/reason`.

### 4.2 Estado de orden propia

```text
INTENT_CREATED → RISK_ACCEPTED → SENT
SENT → ACCEPTED | REJECTED | UNKNOWN
ACCEPTED → LIVE | PARTIALLY_FILLED | FILLED | CANCELED | EXPIRED
cancel_requested ≠ canceled
replace_requested ≠ replaced
```

`UNKNOWN` es un estado real ante timeout/disconnect. Nunca reenviar ciegamente: reconciliar por client ID, venue ID, execution reports y query. Nasdaq OUCH permite repetir inbound benignamente bajo su referencia day-unique/creciente; eso es semántica del protocolo, no licencia para asumir idempotencia en cualquier API.

## 5. Reconstrucción correcta del libro local

### 5.1 Patrón snapshot + delta

**[SPEC]** Binance Spot documenta el patrón:

1. abrir stream y bufferizar deltas;
2. obtener snapshot;
3. repetir snapshot si termina antes del primer delta buffered;
4. descartar deltas ya incluidos;
5. exigir que el primer delta puente el `lastUpdateId`;
6. instalar snapshot y aplicar buffered + live en secuencia;
7. si aparece un gap, descartar el libro y reiniciar.

La regla de aplicación L2 es **set quantity**, no sumar: cantidad cero elimina el nivel. Implementar con enteros fixed-point tras validar tick/lot; nunca `float` binario como key de precio.

### 5.2 Variantes reales

- **Coinbase Exchange L2:** entrega snapshot completo del canal y updates; L3/full requiere subscribir, encolar, pedir snapshot REST, descartar secuencias cubiertas y reproducir el resto. `received` no significa resting; modificar por él produce un libro falso.
- **Kraken Spot WS v2:** snapshot/update incluyen CRC32 de top 10; procesar múltiples cambios del mismo precio en orden y verificar checksum según formato oficial.
- **Nasdaq ITCH 5.0:** `Add` crea la referencia visible; `Executed` reduce cantidad y puede repetirse acumulativamente; `Cancel` es reducción parcial; `Delete` elimina el remanente; `Replace` retira la referencia original y crea otra con precio/cantidad nuevos. No reinterpretar `Trade (Non-Cross)` como add: existe para matches de interés no displayable. Los mensajes administrativos, directory, trading action y system event también mutan la interpretación.
- **Prioridad ITCH:** una nueva order reference tras replace representa una orden nueva para el handler; no conservar la posición de cola anterior salvo que la especificación/rulebook aplicable lo garantice expresamente.
- **Codificación ITCH:** enteros big-endian unsigned salvo nota, ASCII padded, precios fixed-point y timestamp como nanosegundos desde medianoche. `stock locate` se asigna diariamente: sirve como índice intradía, no como identidad persistente.
- **Transporte ITCH:** el payload es una serie sequenced, pero la entrega/recuperación pertenece a SoupBinTCP, MoldUDP64 y GLIMPSE. El decoder no puede atribuir al payload una garantía que sólo existe en el protocolo inferior.
- **CME MDP 3.0:** market data de futuros/opciones usa SBE sobre arquitectura multicast dual-feed. A/B son copias redundantes, no dos mercados para mezclar; channel sequencing, incremental refresh, recovery/snapshot y security definition forman un contrato conjunto.
- **Eurex T7 EOBI:** incremental y snapshot pueden viajar por canales multicast distintos sincronizados mediante campos definidos por la versión. La recuperación debe seguir el manual EOBI vigente; ETI es el plano privado de trading, no sustituto del feed público.

“WebSocket usa TCP” no prueba consistencia de aplicación: reconnect, resubscribe, snapshot race, servidor distinto, mensaje parseado parcialmente o buffer overflow pueden abrir huecos.

### 5.3 Invariantes de libro

```text
all quantities >= 0
all prices and quantities satisfy venue grids
best_bid < best_ask during uncrossed continuous state
sequence/checksum contract holds
zero quantity removes a level
same event replayed under its contract does not double effect
snapshot published atomically belongs to one epoch/sequence
no derived metric crosses an invalid interval silently
```

Crossed/locked book puede ser evento válido de auction, consolidación, latency o feed; la invariante se condiciona al estado/venue, no se “arregla” ordenando precios.

### 5.4 Arquitectura single-writer

```text
socket/NIC capture
  → framing + decoder
  → channel sequencer
  → single writer per book/shard
  → atomic immutable view/version
  → analytics/readers
  → raw event log + checkpoints
```

Single writer elimina locks del mutation path, no backpressure. Todo ring/queue tiene slots, bytes, age y overflow policy; para market state, drop silencioso está prohibido: gap→invalidar→resync. Los readers reciben `{epoch, sequence, state, view}` coherente, no referencias mutables.

### 5.5 Order entry: OUCH y FIX no son la misma sesión

**[SPEC] OUCH 5.0:** host→client debe llegar secuenciado mediante un protocolo inferior; client→host no está secuenciado por OUCH y sus mensajes pueden repetirse benignamente. Esta propiedad depende de `UserRefNum`: es day/account-scoped, único y estrictamente creciente para enter/replace. Un número menor al esperado puede considerarse duplicado e ignorarse; tras restart se consulta el siguiente disponible. Replace posee semántica propia —incluida responsabilidad acumulada sobre la cadena— y puede ser ignorado, cancelar la original, rechazarse o producir `Replaced`; sólo el response confirma el resultado.

**[SPEC] FIX Session Layer:** `MsgSeqNum(34)` ordena mensajes de sesión y aplicación en un único espacio; `NextNumIn/NextNumOut` debe persistir a través de conexiones. Un gap genera `ResendRequest(35=2)` y los mensajes nuevos se retienen hasta cerrar el hueco. Retransmission conserva el sequence original y usa `PossDupFlag(43)=Y`; `SequenceReset(35=4, GapFillFlag=Y)` salta mensajes que no se retransmiten. Resetear la sesión no es una reparación inocua: puede ocultar estado de aplicación y causar overfill.

FIX distingue dos operaciones que el código no debe fusionar:

```text
session retransmission:
  same MsgSeqNum + PossDup=Y -> recuperar delivery; deduplicar por sesión

application resend after missing business ACK:
  new MsgSeqNum + PossResend=Y/ID persistente según rules of engagement
  -> la aplicación, no la sesión, resuelve duplicidad e idempotencia
```

Por tanto, `TCP connected`, `FIX logged on`, `sequence synchronized`, `order acknowledged` y `book reconciled` son hechos distintos. El rules-of-engagement bilateral fija qué application messages se retransmiten, ventanas, IDs, resets y recuperación; la biblioteca genérica no puede decidirlo.

## 6. Tiempo y causalidad operacional

Mantener al menos:

```text
venue_event_time
gateway_or_feed_send_time  # si existe
kernel_or_hw_receive_time
userspace_receive_time
decode_apply_time
publish_time
consumer_observe_time
```

Un timestamp con nanosegundos de resolución no implica exactitud nanosegundo. Documentar clock domain, sync/PTP/NTP, hardware/software capture, monotonicidad, leap/session handling y error bound. Nunca ordenar venues distintos sólo por timestamps sin incertidumbre; cada venue tiene secuencia propia y no existe un “orden global verdadero” derivado automáticamente.

Medir latencia por distribución y fase: wire/capture, queue, decode, sequence, apply, publish. Correlacionar `sequence/epoch`, no sólo wall clock.

## 7. Métricas para comprender el libro

Para `bid=b`, `ask=a`, tamaños top `q_b,q_a`:

```text
spread_ticks = (a - b) / tick
mid = (a + b) / 2
top_imbalance = (q_b - q_a) / (q_b + q_a)
microprice = (a*q_b + b*q_a) / (q_b + q_a)
```

Validar denominadores y estado uncrossed/live. Microprice y imbalance resumen top-of-book; hidden liquidity, queue position, event intensity y regime pueden invalidar su interpretación.

Medir además:

- depth acumulada por distancia en ticks/bps y costo de walk para notional fijo;
- add/cancel/execute rates por lado/nivel, queue depletion y refill;
- spread duration, book age, update intensity y burstiness;
- order-flow imbalance (OFI) que combina cambios de precio/tamaño en best bid/ask;
- realized volatility/returns en event time y clock time;
- trade size, inter-arrival, signed flow con método de clasificación declarado;
- resiliency: tiempo/proceso de refill tras shock;
- slippage/implementation shortfall sólo contra un arrival benchmark y fees reales.

Para trade price `P_t`, mid inmediatamente anterior `M_t`, mid futuro comparable `M_{t+Δ}` y lado agresor `D_t∈{+1 buy,-1 sell}`:

```text
effective_spread = 2 * D_t * (P_t - M_t)
realized_spread_Δ = 2 * D_t * (P_t - M_{t+Δ})
price_impact_Δ    = 2 * D_t * (M_{t+Δ} - M_t)
effective_spread  = realized_spread_Δ + price_impact_Δ
```

Declarar cómo se obtuvo `D_t`, qué quote era observable antes del trade, qué clock y qué horizonte usa `Δ`. Trades fuera de secuencia, midpoint stale, auctions, crossed markets y clasificación por tick rule pueden sesgar la descomposición. Reportar distribución y tamaño ponderado además de media; no confundir bid–ask bounce mecánico con información.

**[ACADEMIC]** Cont–Kukanov–Stoikov hallaron para su muestra/ventanas una relación aproximadamente lineal entre cambio de precio y OFI, con pendiente inversa a depth. Es evidencia condicional y explicativa; no autoriza asumir causalidad estable, transferencia a cripto, ausencia de costs ni capacidad predictiva out-of-sample.

### 7.1 Disciplina estadística

- separar evento contemporáneo de predictor disponible antes de la decisión;
- train/validation/test por tiempo, con purging/embargo si labels se solapan;
- no mezclar sesiones, ticks, fees o regimes sin normalización/strata;
- evitar leakage de snapshot futuro, late corrections o consolidated feed;
- reportar coverage, missing/gap intervals y survivorship/listing bias;
- comparar contra baseline simple y costos; intervalos/confianza, no sólo promedio;
- drift y falsificación: una relación debe poder dejar de cumplirse.

### 7.2 Análisis sistemático de errores

Cada discrepancia entre replay, observación y modelo entra en un ledger; no se “limpia” antes de clasificarla:

| Clase | Ejemplos | Respuesta |
|---|---|---|
| adquisición | gap, truncation, stale snapshot, packet loss | invalidar intervalo, recuperar y medir cobertura |
| decoder/schema | endian, precision, campo nuevo, framing | golden bytes, version gate y fail closed |
| state machine | evento legal no contemplado, terminal conflict | reducir a secuencia mínima y añadir invariant/test |
| sincronización | snapshot race, duplicate, replay boundary | reproducir con epoch/sequence y corregir recovery |
| semántica | STP/amend/auction/rulebook mal entendido | volver a spec vigente, no ajustar el dato al modelo |
| reloj | drift, reorder aparente, clock step | preservar orden de secuencia y marcar incertidumbre |
| modelo | residual por régimen, hidden liquidity, misspecification | slice, falsificar, recalibrar o retirar el modelo |
| operación | overload, logging stall, credential/rate limit | capacity test, runbook y control independiente |

Priorizar por `frecuencia × severidad × alcance × detectabilidad`, revisar muestras reales y convertir cada clase material en test, métrica o guardrail. Un menor error agregado no compensa intervalos silenciosamente inválidos: primero data/state correctness; luego representatividad; finalmente ajuste del modelo.

## 8. Liquidez no es una cifra

Separar:

- **tightness:** spread/costo inmediato;
- **depth:** cantidad visible a distintos precios;
- **immediacy:** tiempo para ejecutar/reponer;
- **resiliency:** recuperación tras flujo/shock;
- **adverse selection:** movimiento posterior condicionado al fill;
- **capacity:** tamaño antes de impacto/material delay;
- **reliability:** probabilidad de que dato/venue/control esté disponible.

Displayed depth puede cancelarse antes de alcanzarla, ocultar reserve o competir con otra demanda. Un libro “grande” no garantiza ejecución; estimar fill requiere queue/matching/latency y sigue siendo probabilístico.

## 9. Multi-venue y consolidación

Antes de comparar:

- canonicalizar instrumento, contract multiplier, base/quote, expiry y settlement;
- convertir tick, lot, notional, currency y fee/rebate tier;
- registrar trading state, session y staleness por venue;
- mantener feed sequence y clock uncertainty separados;
- distinguir displayed executable price de index/mark/last/indicative;
- incorporar transfer/custody/capital/credit/withdrawal constraints cuando sean relevantes.

NBBO o “best global crypto price” es una construcción con routing, latencia y reglas; no fusionar niveles como si pertenecieran a un único matching engine. Un crossed consolidated view puede expresar desalineación temporal, fees o fragmentación, no arbitraje ejecutable.

### 9.1 Spot, futuros y perpetuals son estados económicos distintos

Para derivados, el libro CLOB es sólo una parte. El contrato debe añadir:

```yaml
derivative:
  type: future | option | perpetual
  contract_multiplier: ""
  inverse_or_linear: ""
  expiry_settlement_delivery: {}
  index_composition_and_quality: {}
  mark_price_formula_and_source: {}
  funding_formula_interval_cap: {}
  initial_maintenance_margin: {}
  liquidation_and_insurance_adl: {}
  position_open_interest_limits: {}
```

- **last/trade price** describe executions; **index** agrega referencias externas según metodología; **mark** alimenta riesgo/margen; **settlement** cristaliza obligaciones; no son intercambiables.
- **funding** es un cash flow periódico long↔short definido por venue, no fee de trading ni retorno garantizado; distinguir predicted del final aplicado.
- **liquidation** es flujo forzado por margin engine. Puede impactar el libro, pero una liquidation stream puede ser agregada/truncada y no revela necesariamente cada orden interna.
- **open interest** cuenta exposición contractual abierta bajo reglas del venue; volume cuenta transferencias ejecutadas. No inferir uno del otro.
- Cambios de margin tiers, index constituents, price bands, funding caps o maintenance rules son cambios de contrato aunque el símbolo permanezca.

Coinbase International, por ejemplo, publica bid/ask/trade junto con index, mark, settlement, limits y predicted/final funding como entry types separados. Esa separación debe sobrevivir desde decoder hasta almacenamiento y dashboards.

### 9.2 CLOB, RFQ y AMM/DEX no comparten la misma microestructura

| Mecanismo | Estado/prioridad | Precio/ejecución | Riesgos adicionales |
|---|---|---|---|
| CLOB | órdenes y colas bajo matching rule | cruza resting liquidity/auction | hidden, queue, latency, halts |
| RFQ | quotes dirigidas con expiry/firmness | aceptación bilateral o venue workflow | information leakage, last look, credit |
| AMM/CFMM | reservas, curva, ticks/ranges y estado on-chain | función/invariant + fee + route | gas, block ordering/MEV, revert, token/oracle/contract |

**[SPEC/CODE]** Uniswap v2 impone el invariant constant-product después de fees; v3 concentra liquidez en rangos y atraviesa ticks. No convertir reserves virtuales o posiciones LP en órdenes FIFO, ni comparar “depth” sin incluir fee, price impact, gas, route, slippage guard y estado del bloque. Una transacción atómica puede revertir completa: `submitted`, `included`, `finalized` y `economically settled` son estados diferentes.

Esta biblioteca sólo fija la frontera. Implementar AMM/DEX exige versión de contratos, chain/finality, MEV model, oracle, token behavior, addresses/deployments y auditorías de seguridad específicas.

## 10. Order entry y controles independientes

### 10.1 Gate pre-trade

Antes del encoder/socket:

```text
authorized strategy/account/session
instrument and trading-state allowed
price/qty/tick/lot/notional valid
fat-finger and price collar
max order/message rate and open orders
position/exposure/credit/capital limits
duplicate/client-ID rule
STP and restricted-list rule
```

**[SPEC]** SEC Rule 15c3-5 exige a broker-dealers alcanzados controles pre-trade de exposición/órdenes erróneas, requisitos regulatorios, acceso autorizado y reporting post-trade. No se aplica automáticamente a toda persona/cripto/jurisdicción, pero muestra una separación profesional: el strategy process no controla ni puede saltar su propio risk gateway.

Ejemplos de semántica que el adapter debe preservar:

- Binance `-1006`/`-1007` declara execution/send status `unknown`; reconciliar por client/order ID y user stream antes de decidir un resend.
- Coinbase puede publicar `received` en el feed antes de que vuelva la respuesta REST con server order ID; correlacionar mediante client ID y aceptar ordering races documentadas.
- STP no significa siempre “rechazo”: según venue/modo puede decrementar cantidad, cancelar resting, aggressing o ambos y generar eventos específicos.
- Un amend exitoso puede conservar prioridad sólo para campos/endpoint autorizados; todo otro cancel-replace debe asumir prioridad nueva.

### 10.2 Durante y después

- drop-copy/execution stream independiente del strategy path;
- positions/exposure desde executions confirmadas, no desde intención;
- cancel-on-disconnect y mass cancel probados por scope;
- kill switch local/gateway/venue con estado observable y access dual cuando proceda;
- reject/partial/late fill/duplicate/out-of-order reconciliados;
- limits fail closed para nuevas exposiciones; risk-reducing actions conservadas bajo policy;
- clock/session/cert/key expiry y rate-limit exhaustion probados.

Kill no significa flat: puede quedar fill in-flight, posición, orden en otro session/venue o cancel rechazado. El cierre verifica open orders + executions + positions contra fuente autorizada.

### 10.3 Integridad de mercado y límites de inferencia

Un observador puede generar **señales para revisión**, no sentencias sobre intención:

- cancel-to-add/fill ratios condicionados por nivel, lifetime, spread y régimen;
- layering-like shape, quote stuffing/bursts, marking-window concentration;
- self-match/wash indicators cuando existe beneficial-owner/account linkage autorizado;
- opposite-side execution tras retirar depth, cross-venue/cross-instrument sequence;
- concentración por participante sólo si el feed realmente contiene attribution.

**[SPEC]** La guía CFTC sobre spoofing hace material la intención de cancelar antes de ejecución y exige contexto, patrón, fills y demás circunstancias; una cancelación legítima o modificación de buena fe no se vuelve ilícita por una ratio alta. L2 público normalmente no identifica actor ni beneficial ownership, y hidden/reserve/other venues alteran el denominador. Por tanto:

```text
heuristic alert -> evidence bundle -> human/compliance investigation
                -> applicable rule + jurisdiction + intent/context analysis
```

Conservar raw lineage, clocks, rulebook, parámetros/versiones del detector y explicación de por qué alertó. Medir precision/recall sobre casos adjudicados cuando existan, revisar disparate impact y permitir reproducir/falsificar. Nunca automatizar acusaciones públicas o sanciones irreversibles desde una anomalía estadística.

## 11. Matching engine y exchange interno

Si se construye un simulador o venue:

```text
gateway/risk
  → deterministic sequencer
  → matching shard single-writer
  → executions + order events
  → market-data encoder
  → durable audit/drop copy/clearing
```

Definir explícitamente:

- total order por instrument/shard y qué timestamp asigna prioridad;
- integer price/qty y overflow/rounding;
- order IDs, duplicate/replay y session epoch;
- matching/auction/STP/amend/hidden rules;
- commit point: cuándo ack/fill se vuelve authoritative y durable;
- replication/failover/fencing y recovery desde log+checkpoint;
- deterministic replay que produce mismos events/sequence;
- fairness: gateway queues, throttles, partitions y observability iguales bajo policy.

No replicar el matching state con una base general y asumir que una transacción crea orden global. El manual distribuido define consensus/fencing; el de datos define WAL/recovery; éste define la state machine económica y la secuencia autoritativa.

## 12. Failure model de market data

| Falla | Estado correcto | Acción |
|---|---|---|
| disconnect limpio | `DISCONNECTED/STALE` | reconnect, nueva epoch y resync según spec |
| sequence gap | `GAPPED` | dejar de publicar; snapshot/retransmit/replay |
| duplicate/old | contrato dependiente | ignorar sólo si la secuencia demuestra inclusión |
| checksum mismatch | `CHECKSUM_FAILED` | preservar raw evidence y resync |
| parser/schema desconocido | `INVALID` | no adivinar; alert/version gate |
| snapshot demasiado viejo | `SYNCING` | repetir sin instalar estado imposible |
| local queue overflow | `GAPPED` | invalidar; nunca continuar “best effort” silencioso |
| clock step/drift | data válida, tiempo incierto | marcar clock quality; no reordenar secuencia |
| disk/logger lento | policy acotada | no bloquear feed indefinidamente; degradar analytics o fail capture explícito |

## 13. Suites ejecutables

| Suite | Casos | Gate |
|---|---|---|
| `feed_codec` | cada message type, boundary, malformed, schema/version | parse exacto/bounded o rechazo |
| `book_rebuild` | snapshot race, buffered deltas, add/set/delete, same-price repeats | estado coincide con golden/reference |
| `book_gap` | drop/duplicate/reorder/reconnect/checksum/overflow | nunca publica `LIVE` incorrecto; resync converge |
| `session_state` | open/auction/halt/resume/close/directory change | reglas sólo en estado elegible |
| `order_lifecycle` | ack/reject/partial/fill/cancel/replace/expire/unknown | estado legal e idempotencia/reconcile |
| `matching` | price levels, FIFO/pro-rata, TIF, STP, amend, auction | executions exactas según rulebook fijado |
| `risk_gateway` | fat finger, credit/position/rate, stale data, kill/disconnect | orden bloqueada antes de venue; cancel/closure verificadas |
| `replay` | raw log + checkpoint, crash en cada boundary | mismo book/events/sequence o divergencia detectada |
| `latency_capacity` | burst/soak, gap storm, resync, slow consumer, disk slow | goodput, p99.9 y bounded queues dentro de SLO |
| `metrics` | hand-worked books, invalid intervals, missing data, regimes | fórmula/ventana/availability correctas |

Property tests útiles:

```text
replay(prefix + suffix) == replay(checkpoint(prefix), suffix)
every LIVE version has one continuous validated lineage
no order has two incompatible terminal states
filled_qty + leaves_qty == accepted_qty after venue adjustments
book state never depends on reader scheduling
```

## 14. Benchmark de baja latencia

Reportar:

- hardware/NIC/OS/runtime/NUMA/affinity/IRQ y mitigations;
- feed/protocol/schema, symbols, message mix, burst model y bytes;
- cold/warm, duration, loss/gap/resync y background logging;
- offered rate, decoded/applied/published **goodput** y drops;
- p50/p95/p99/p99.9/max por fase y queue age;
- allocations, cache misses, branch misses, context switches, GC si existe;
- correctness hash/checksum junto al resultado.

Messages/s sin validar libro puede premiar un decoder que descarta trabajo. Un median-only benchmark oculta stalls capaces de abrir gaps. Kernel bypass, busy polling o FPGA sólo se justifican después de localizar el costo end-to-end y mantener recovery/observability.

## 15. Contrato para Codex

Entrada mínima:

```yaml
task: "capture | book | analytics | order gateway | matching simulation"
venue_instrument_environment: ""
official_specs_and_versions: []
observable_data_level: "L1 | L2 | L3 | own orders"
correctness_invariants: []
sequence_snapshot_recovery: {}
latency_capacity_slo: {}
risk_and_permissions: "read-only by default"
deliverables: []
```

El agente debe:

1. citar la spec vigente y marcar `[OPEN]` si falta;
2. escribir state machines antes de handlers;
3. usar fixed-point/checked arithmetic y límites;
4. separar raw event, normalized event, authoritative state y metric;
5. implementar gap/resync antes de optimizar;
6. entregar golden replay, fault tests y benchmark con correctness;
7. mantener order entry deshabilitado salvo autorización explícita y entorno declarado;
8. no presentar una correlación como estrategia, causalidad o rentabilidad.

## 16. Ruta práctica profesional

1. **Decoder offline:** ITCH o un feed cripto capturado; golden messages, endian/fixed-point y malformed corpus.
2. **L2 local:** Binance/Coinbase/Kraken read-only; snapshot/delta/checksum, epochs y deterministic replay.
3. **Observador:** spread/depth/OFI/event rates con missing-data masks y dashboards de feed health.
4. **Stress:** bursts, reconnects, gaps, disk slow, clock drift; medir goodput/p99.9 y resync time.
5. **L3/matching simulator:** price-time primero, luego algoritmo alternativo; order lifecycle y auctions.
6. **Gateway simulado/testnet:** pre-trade risk, IDs, ambiguous outcome, drop copy, mass cancel/kill; sin capital real.
7. **Multi-venue:** normalización y staleness, nunca inventar orden global.

## 17. Auditoría transversal cerrada

| Manual | Autoridad de frontera | Obligación aquí |
|---|---|---|
| `COMPUTER_SYSTEMS_PERFORMANCE_LOW_LATENCY.md` | CPU/cache/NUMA/atomics/profiling | single writer y queues conservan correctness; medir path real |
| `SOFTWARE_BACKEND_API_ENGINEERING.md` | API/idempotencia/deadline/lifecycle | outcome `UNKNOWN`, auth/rate/reconcile y no retry ciego |
| `NETWORKING_DISTRIBUTED_STREAMING.md` | transport/sequence/backpressure/consensus | feed gap y retransmit no se confunden con book semantics |
| `DATABASE_STORAGE_INTERNALS.md` | WAL/checkpoint/replay/durability | raw log y matching commit declaran ACK/durable point |
| `GPU_ACCELERATED_COMPUTING.md` | batch/device/async completion | analytics GPU no bloquea book ni publica output tardío/inválido |
| `SECURITY_SRE_CLOUD_INFRASTRUCTURE.md` | IAM/provenance/SLO/incident/recovery | read-only default; risk/kill independientes; secrets y audit protegidos |
| `FRONTEND_PRODUCT_ENGINEERING_UX.md` | navegador, interacción, accesibilidad y estado visible | UI porta venue/instrument/epoch/sequence/freshness; muestra gaps/resync y nunca habilita acción fuera del risk gate |

Cruces recíprocos cerrados el 2026-08-21. Reglas de arbitraje:

- este manual decide si book/order/matching/metric es semánticamente correcto;
- sistemas decide el costo físico; red, delivery; backend, adapter/lifecycle; datos, durability/replay; GPU, device/batch; seguridad/SRE, permisos y operación;
- ninguna capa puede “reparar” silenciosamente un gap, inventar prioridad, cambiar estado terminal o habilitar acción fuera del risk gate;
- cada output derivado porta `venue/instrument/contract_version/epoch/sequence/state/freshness/model_version` aplicables.

Cruce frontend cerrado el 2026-08-21. Coalescing/virtualización puede reducir renders, no borrar gaps, reordenar eventos auditables ni mostrar un mid viejo como vivo. Confirmaciones, hotkeys y optimistic state no reemplazan pre-trade risk, ACK/reject/drop-copy o reconciliación de outcome `UNKNOWN`; keyboard/focus y alertas accesibles forman parte del control operacional.

Cruce con `ALGORITHMS_DATA_STRUCTURES_PROBLEM_SOLVING.md` cerrado el 2026-08-21. Price levels requieren orden/predecessor, FIFO requiere cola estable y sequence/gap requiere state machine; un heap por sí solo no representa cancelación arbitraria ni prioridad precio-tiempo. Elegir array/radix/tree/hash/queue por rango de ticks, profundidad, updates, locality y venue semantics. Oráculo lento, replay determinista, invariantes y secuencias adversariales preceden al benchmark; ninguna mejora asintótica puede inventar liquidez, prioridad o continuidad.

Cruce con `SOFTWARE_ARCHITECTURE_SYSTEM_DESIGN.md` cerrado el 2026-08-21. Arquitectura separa feed/hot state, raw durable log, analytics, order entry/risk, API y UI por quality/failure drivers; este manual decide book, matching, order lifecycle y venue semantics. El runtime view porta epoch/sequence/freshness/ACK/durable point; ninguna boundary autoriza reconstrucción silenciosa, retry ciego o IA en el hot path sin evidencia.

Cruce con `DATA_ENGINEERING_ANALYTICS.md` cerrado el 2026-08-21. Preservar capas raw frame → decoded event → validated book → derived metric → aggregate, con venue/instrument/contract/epoch/sequence/event/receive time. Un gap invalida derivados hasta resync; backfill/replay no inventa observabilidad pasada. Analytics no controla matching state y ninguna imputación, forward-fill o aggregate se presenta como evento del venue.

## 18. Matriz de trazabilidad v1

| Área | Autoridad | Estado |
|---|---|---|
| teoría LOB/spread/impact/inventory/colas | Stanford MATH 237; Glosten–Milgrom; Kyle; Cont; Avellaneda–Stoikov; Gould; queue-reactive | mecanismos y límites incorporados; validación empírica por proyecto |
| matching/allocation | Eurex T7, CME y rulebooks de venue | principio incorporado; matriz por producto condicionada |
| market data L3 | Nasdaq TotalView-ITCH 5.0 | lifecycle, binary types, daily identity y transport boundary incorporados |
| order entry/session | Nasdaq OUCH 5.0 + FIX Session Layer | lifecycle, duplicate/replay, sequence recovery y app idempotency incorporados; SBE/FIXP queda por proyecto |
| crypto L2/L3 + order semantics | Binance, Coinbase, Kraken oficiales | snapshot/delta/gap/checksum y ejemplos de priority/STP/unknown outcome incorporados; adapter completo se fija por proyecto |
| risk/regulación | SEC 15c3-5 + CFTC | gate técnico incorporado; jurisdicción condicionada |
| market integrity | CFTC spoofing guidance + SEC enforcement context | heurística separada de determinación legal/intención; evidence workflow incorporado |

## 19. Extensiones deliberadamente condicionadas

- Las notas no públicas/restantes de Stanford MATH 237 sólo se incorporan si se obtiene acceso legítimo y agregan una obligación operacional no cubierta; ningún modelo se convierte en receta de trading.
- Completar codecs message-by-message y transports SoupBinTCP/MoldUDP64/GLIMPSE sólo al seleccionar ese feed; FIXP/SBE requiere versión y venue concretos.
- Implementar CME MDP 3.0/iLink o Eurex T7 ETI/EOBI message-by-message sólo con versión/release, entitlement y producto seleccionados; aquí quedan arquitectura y recovery boundaries.
- Expandir derivatives/liquidation/funding o AMM/DEX sólo contra contrato, jurisdicción y protocolo elegidos; aquí queda su separación arquitectónica.
- El adapter productivo de Binance/Coinbase/Kraken debe fijar el rulebook vigente completo: matching, priority, fees, STP, limits, halts y order states; el núcleo no congela configuración cambiante.
- Clearing/settlement/custody, derivados/liquidación/funding completos y compliance por jurisdicción requieren módulos específicos cuando el proyecto los declare.

## 20. Fuentes primarias auditadas

### Academia

- Stanford MATH 237, George Papanicolaou: <https://math.stanford.edu/~papanico/Math237/CourseInfo.html>.
- Stanford MATH 237, plan público: <https://math.stanford.edu/~papanico/Math237/LecturePlan.html>. El índice confirma LOB/impact/control estocástico, Hawkes/feedback y statistical arbitrage; el enlace de materiales devolvió `401` al auditarlo el 2026-08-21, por lo que no se afirma cobertura literal/completa de esas notas.
- Cont, Kukanov y Stoikov, *The Price Impact of Order Book Events*: <https://arxiv.org/abs/1011.6402>.
- Cont y de Larrard, *Price dynamics in a Markovian limit order market*: <https://arxiv.org/abs/1104.4596>.
- Avellaneda y Stoikov, *High-frequency trading in a limit order book*: <https://math.nyu.edu/inmemoriam/avellaneda/HighFrequencyTrading.pdf>.
- Glosten y Milgrom, *Bid, Ask and Transaction Prices in a Specialist Market with Heterogeneously Informed Traders*: <https://business.columbia.edu/faculty/research/bid-ask-and-transaction-prices-specialist-market-heterogeneously-informed-traders>.
- Kyle, *Continuous Auctions and Insider Trading*: <https://people.stern.nyu.edu/lpederse/courses/LAP/papers/Information%2CFundamental/Kyle85.pdf>.
- Gould et al., *Limit order books*: <https://arxiv.org/abs/1012.0349>.
- Huang, Lehalle y Rosenbaum, *Simulating and analyzing order book data: The queue-reactive model*: <https://arxiv.org/abs/1312.0563>.
- Wu et al., *Queue-reactive Hawkes models for the order flow*: <https://arxiv.org/abs/1901.08938>.

### Venues y protocolos

- Nasdaq TotalView-ITCH 5.0: <https://classic.nasdaqtrader.com/content/technicalsupport/specifications/dataproducts/NQTVITCHSpecification.pdf>.
- Nasdaq OUCH 5.0: <https://www.nasdaqtrader.com/content/technicalsupport/specifications/TradingProducts/OUCH5.0.pdf>.
- Eurex T7 matching principles: <https://www.eurex.com/ex-en/trade/order-book-trading/matching-principles>.
- Eurex T7 Release 14.0 EOBI manual: <https://www.eurex.com/resource/blob/4597908/bf8b02f4f4d1f220e54aa1a93f0482df/data/T7_R.14.0_%20EOBI_Manual_Version_1.pdf>.
- CME Group market data platform/MDP 3.0: <https://www.cmegroup.com/market-data/distributor/market-data-platform.html>.
- FIX Trading Community: <https://www.fixtrading.org/standards/>.
- FIX Trading Community, *FIX Session Layer*: <https://www.fixtrading.org/standards/fix-session-layer-online/>.
- FIX Trading Community, session conformance test cases: <https://www.fixtrading.org/standards/fix-session-testcases-online/>.
- Binance Spot official API docs: <https://github.com/binance/binance-spot-api-docs>.
- Binance Spot, REST/order errors y amend keep priority: <https://github.com/binance/binance-spot-api-docs/blob/master/rest-api.md>, <https://github.com/binance/binance-spot-api-docs/blob/master/errors.md>, <https://github.com/binance/binance-spot-api-docs/blob/master/faqs/order_amend_keep_priority.md>.
- Coinbase Exchange WebSocket channels: <https://docs.cdp.coinbase.com/exchange/websocket-feed/channels>.
- Coinbase Exchange matching engine: <https://docs.cdp.coinbase.com/exchange/concepts/matching-engine>.
- Kraken Spot WebSocket v2 book: <https://docs.kraken.com/api/docs/websocket-v2/book/>.
- Coinbase International FIX market data: <https://docs.cdp.coinbase.com/international-exchange/fix-api/market-data>.

### AMM/DEX

- Uniswap v2 whitepaper: <https://docs.uniswap.org/whitepaper.pdf>.
- Uniswap v3 whitepaper: <https://app.uniswap.org/whitepaper-v3.pdf>.
- Uniswap developer protocol overview: <https://developers.uniswap.org/docs/get-started/concepts/how-uniswap-works>.

### Riesgo y regulación

- SEC Rule 15c3-5, Market Access: <https://www.sec.gov/rules-regulations/2011/06/risk-management-controls-brokers-or-dealers-market-access>.
- CFTC, Electronic Trading Risk Principles: <https://www.cftc.gov/LawRegulation/FederalRegister/finalrules/2020-27622.html>.
- CFTC, interpretive guidance on disruptive practices/spoofing: <https://www.cftc.gov/LawRegulation/FederalRegister/FinalRules/2013-12365.html>.
- SEC, market-structure enforcement and data/code evidence: <https://www.sec.gov/newsroom/speeches-statements/ceresney-speech-sifma-ny-regional-seminar>.

## Regla final

```text
First recover the exact state.
Then measure the mechanism.
Only then form a hypothesis.
Never let a hypothesis bypass sequence, risk or evidence.
```
