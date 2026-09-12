# Algorithms, Data Structures & Problem Solving — manual operativo

> **Estado:** núcleo general v1 auditado el 2026-08-21; fuentes, cobertura curricular, estructura interna y cruces recíprocos cerrados. Contrato ejecutable y capstones de planificación disponibles en `ENGINEERING_EXECUTION_PLAYBOOK.md`; la evidencia de una implementación concreta sigue perteneciendo al proyecto.
> **Autoridad académica:** MIT 6.006 (Spring 2020 y contrato de entrega Fall 2011); MIT 6.046J (Spring 2015); Princeton *Algorithms, 4th Edition*.
> **Propósito:** convertir problemas reales en modelos, algoritmos correctos, estructuras adecuadas, implementaciones robustas, tests adversariales y evidencia de costo.
> **Frontera:** enseña a decidir y demostrar *qué trabajo debe hacerse*. El costo físico, concurrencia, storage, red y GPU se remiten a sus manuales especializados.

## Cómo usar este manual

No pedir a una persona o agente solamente “implementá el algoritmo”. Entregar primero el contrato del problema, después exigir una solución verificable y recién entonces optimizar.

```text
especificación
→ modelo y restricciones
→ baseline correcto
→ estructura/paradigma
→ prueba de corrección
→ complejidad bajo un costo declarado
→ implementación defensiva
→ tests diferenciales y adversariales
→ medición representativa
→ decisión documentada
```

Toda solución material debe poder responder:

1. ¿Qué entrada acepta y qué salida promete?
2. ¿Qué supuestos hacen posible el algoritmo?
3. ¿Por qué termina y por qué el resultado es correcto?
4. ¿Qué recurso domina el costo y bajo qué modelo?
5. ¿Cuál es el peor caso, el caso esperado y el costo amortizado pertinente?
6. ¿Qué entradas rompen la implementación, no sólo la idea?
7. ¿Qué alternativa se descartó y por qué?
8. ¿La medición confirma que el modelo representa este sistema?

## Procedencia

- `[SPEC]`: comportamiento normativo o contrato oficial versionado.
- `[ACADEMIC]`: síntesis de curso, texto o paper identificado; no es transcripción literal.
- `[CODE]`: implementación oficial usada para contrastar una idea.
- `[PROD]`: evidencia operacional publicada y atribuida.
- `[MEASURED]`: experimento reproducible dentro de un entorno declarado.
- `[IMPL]`: traducción propia a diseño, código, prueba o runbook.
- `[OPEN]`: afirmación o extensión todavía no auditada.

## 1. Contrato de un problema computacional

### 1.1 Especificar antes de elegir

**[ACADEMIC]** MIT 6.006 presenta algoritmos como modelado matemático de problemas computacionales y enfatiza la relación entre algoritmo, programa y medición de rendimiento. La estructura de datos no se elige por costumbre: se deriva de las operaciones que el problema necesita.

**[IMPL]** Completar esta ficha:

```text
objetivo observable:
entradas, tipos y rangos:
salida y criterio de igualdad:
precondiciones:
postcondiciones:
mutabilidad permitida:
duplicados, orden y estabilidad:
online/batch; exacto/aproximado:
tamaño n y demás parámetros relevantes:
distribución normal y entradas adversariales:
presupuesto de tiempo, memoria, I/O y latencia tail:
concurrencia y consistencia:
errores permitidos y política de fallo:
```

La medida `n` rara vez basta. Un grafo requiere al menos `|V|` y `|E|`; strings requieren cantidad y longitud total; un algoritmo numérico puede depender de magnitud o precisión; una consulta externa depende de bytes y accesos, no sólo de registros.

### 1.2 Dominio, representación e interfaz

Separar tres decisiones:

| Capa | Pregunta |
|---|---|
| problema | ¿qué relación matemática debe cumplirse? |
| representación | ¿cómo se codifican estados, claves, aristas o secuencias? |
| implementación | ¿qué tipos, ownership, memoria y API se usan? |

Una representación puede hacer fácil o imposible una operación. Una matriz de adyacencia da consulta de arista constante pero ocupa espacio cuadrático; una lista de adyacencia ocupa `Θ(V + E)` y favorece recorrer vecinos. La decisión depende de densidad y operaciones, no de una regla universal.

### 1.3 Baseline y oráculo

La primera implementación útil suele ser la más simple que sea obviamente correcta, aunque no escale. Sirve como:

- especificación ejecutable;
- oráculo diferencial para tamaños pequeños;
- generador de ejemplos;
- estimación del beneficio máximo;
- defensa contra una optimización que cambia la semántica.

No desplegar el baseline fuera de su rango seguro; conservarlo en tests cuando su costo lo permita.

## 2. Contrato de entrega de una solución

**[ACADEMIC]** El syllabus oficial de MIT 6.006 Fall 2011 exige, al solicitar un algoritmo: descripción textual o pseudocódigo, ejemplo trabajado o diagrama, justificación de corrección y análisis temporal —y espacial cuando corresponda—. También califica corrección, calidad e implementación, incluidas conductas deducibles del enunciado aunque no aparezcan en tests públicos.

**[IMPL]** Para trabajo profesional, ampliar ese contrato:

```text
1. Reformulación precisa y supuestos.
2. Ejemplo mínimo, ejemplo límite y contraejemplo a la solución ingenua.
3. Algoritmo independiente del lenguaje.
4. Invariante o argumento de corrección y terminación.
5. Complejidad con parámetros y modelo de costo explícitos.
6. Elección y representación de estructuras de datos.
7. Implementación con errores, overflow, límites y ownership tratados.
8. Tests unitarios, de propiedades, diferenciales y adversariales.
9. Benchmark sólo si la decisión depende del rendimiento real.
10. Alternativas, umbral de cambio y riesgos residuales.
```

Una respuesta que sólo contiene código está incompleta. Una prueba sin implementación puede ignorar costos reales. Un benchmark sin corrección sólo mide qué tan rápido se produce una respuesta posiblemente equivocada.

## 3. Modelos de costo y crecimiento

### 3.1 Qué significan los límites

| Notación | Significado operativo |
|---|---|
| `O(f(n))` | cota superior asintótica |
| `Ω(f(n))` | cota inferior asintótica |
| `Θ(f(n))` | cota superior e inferior ajustadas |
| `o(f(n))` | crecimiento estrictamente menor |
| `ω(f(n))` | crecimiento estrictamente mayor |

No escribir `O(n)` cuando se conoce `Θ(n)`. No ocultar parámetros: `O(V + E)` es más informativo que “lineal”. No comparar órdenes asintóticos sin considerar tamaños, constantes, memoria, vectorización y distribución de entrada.

### 3.2 Casos distintos

- **Peor caso:** garantía sobre toda entrada válida de tamaño dado.
- **Promedio:** expectativa bajo una distribución explícita de entradas.
- **Esperado:** puede provenir de aleatoriedad interna aun con entrada fija.
- **Amortizado:** costo de una secuencia completa, no probabilidad ni promedio empírico.
- **Output-sensitive:** depende también del tamaño de la salida.
- **Pseudopolinomial:** polinomial en el valor numérico, no en los bits que lo representan.
- **Parametrizado:** separa un parámetro estructural `k` del tamaño total.

### 3.3 Declarar la operación básica

**[ACADEMIC]** Princeton formula el costo a partir de frecuencia de ejecución y costo de operaciones, y exige un modelo explícito. Su método experimental observa, formula una hipótesis falsable, predice y vuelve a medir de forma reproducible.

Modelos habituales:

| Modelo | Operación o recurso dominante |
|---|---|
| RAM | acceso/palabra y operación elemental constantes |
| comparación | cantidad de comparaciones |
| word-RAM | operaciones sobre palabras de `w` bits |
| I/O externo | transferencias de bloques entre niveles |
| cache-aware | conoce tamaño de caché y bloque |
| cache-oblivious | algoritmo no fija esos tamaños, análisis sí los modela |
| paralelo | trabajo total, profundidad/span, comunicación |
| streaming | memoria, pasadas, actualización y error |

Big-O bajo RAM no predice por sí solo latencia real. Para la máquina concreta, enlazar con `COMPUTER_SYSTEMS_PERFORMANCE_LOW_LATENCY.md`.

### 3.4 Recurrencias y sumas

Para divide-and-conquer, expresar y resolver el trabajo:

```text
T(n) = a T(n/b) + f(n)
```

Usar árbol de recursión, sustitución o Master theorem sólo cuando sus condiciones aplican. Para loops no rectangulares, contar iteraciones mediante sumas. Para operaciones sobre secuencias, usar agregación, contabilidad o función potencial.

### 3.5 Límites inferiores

Una cota inferior depende del modelo. Sorting por comparaciones necesita `Ω(n log n)` comparaciones en el peor caso, pero claves enteras con rango aprovechable admiten otras técnicas. Antes de declarar “óptimo”, indicar:

- problema exacto;
- modelo computacional;
- recursos contados;
- tipo de garantía;
- reducción o argumento que establece el límite.

## 4. Pruebas de corrección

### 4.1 Obligaciones mínimas

```text
seguridad: nunca produce un resultado inválido
progreso: avanza hacia una condición terminal
terminación: no continúa indefinidamente
completitud: encuentra una solución cuando existe
```

Para algoritmos aproximados o aleatorizados, añadir cota de calidad y probabilidad de fallo.

### 4.2 Patrones de prueba

| Patrón | Uso típico |
|---|---|
| invariante de loop | estado correcto antes/durante/después de iterar |
| inducción | recursión, tamaño o estructura |
| intercambio | demostrar que una elección greedy puede formar parte de un óptimo |
| cut/cycle | árboles de expansión mínima |
| optimal substructure | programación dinámica |
| contradicción | descartar una solución mejor o un estado imposible |
| reducción | transferir solución o dificultad entre problemas |
| argumento de carga | costo amortizado agregado |
| función potencial | energía contable de una secuencia |

### 4.3 Plantilla de invariante

```text
Inicialización: vale antes de la primera iteración.
Mantenimiento: si vale antes de una iteración, vale después.
Terminación: al finalizar, invariante + condición de salida implican la postcondición.
```

### 4.4 Errores comunes de razonamiento

- demostrar sólo un ejemplo;
- asumir la propiedad que se intenta probar;
- probar optimalidad sin probar factibilidad;
- omitir terminación o el caso vacío;
- usar inducción sin hipótesis suficientemente fuerte;
- confundir correlación experimental con prueba;
- inferir que greedy funciona porque una elección local “parece mejor”;
- tratar tests exitosos como demostración universal.

## 5. Estructuras secuenciales

| Estructura | Fortalezas | Riesgos/limitaciones |
|---|---|---|
| array fijo | contigüidad, acceso indexado, mínimo overhead | tamaño fijo; insertar en medio mueve elementos |
| vector dinámico | acceso `O(1)`, append amortizado `O(1)` | realocación; invalida referencias según lenguaje |
| lista enlazada | inserción local con nodo conocido | pobre localidad, overhead por nodo, búsqueda lineal |
| stack | LIFO, parsing, DFS, undo | overflow lógico; recursión puede consumir stack real |
| queue | FIFO, BFS, scheduling | crecimiento sin límite y backpressure |
| deque | extremos eficientes | layout y constantes dependen de implementación |
| ring buffer | memoria acotada, streaming y colas | wraparound, estado full/empty, concurrencia |

**[IMPL]** Preferir almacenamiento contiguo cuando las operaciones lo permiten. No elegir lista enlazada sólo porque insertar sea asintóticamente constante: primero hay que encontrar la posición y después pagar allocation y locality.

## 6. Diccionarios, conjuntos y orden

### 6.1 Hash table

Promete lookup/inserción/borrado esperados `O(1)` bajo supuestos de hashing y control de carga; no garantiza orden. Revisar:

- función de hash y dominio de claves;
- colisiones y estrategia de resolución;
- load factor, resize y costo amortizado;
- adversario que controla claves;
- estabilidad/reproducibilidad del seed;
- igualdad consistente con hash;
- memoria y locality.

No usar hash table si se necesitan rangos, predecessor/successor, orden determinista o garantía logarítmica adversarial.

#### 6.1.1 Práctica pública Google/Meta — hash tables orientadas a localidad

**[CODE]** Abseil Swiss Tables, commit `a690167a55b9e210b7884d311759977c0f5669c9` (Apache-2.0), y Folly F14, commit `8dc549c54875146647295331282db573c586831e` (Apache-2.0), muestran una decisión industrial reutilizable: separar metadatos compactos de ocupación/fingerprint de los valores y comparar varios candidatos por grupo antes de tocar claves completas. Esto reduce accesos de memoria inútiles; no elimina colisiones ni convierte el caso esperado en garantía adversarial.

No son la misma implementación ni deben tratarse como reemplazos ciegos:

- `flat`/value storage favorece localidad y evita una indirección, pero una mutación o `rehash` puede mover elementos e invalidar referencias;
- node storage paga allocations/indirecciones para conservar estabilidad de punteros o referencias según el contrato concreto;
- F14 además publica variantes vector/fast con otros compromisos de iteración, memoria y movimiento;
- lookup heterogéneo evita construir una clave temporal sólo si `Hash` y `Eq` son transparentes, consistentes y aceptan ambos tipos;
- ni el orden de iteración ni su reproducibilidad deben formar parte de la lógica; Abseil usa sal por tabla y Folly introduce aleatorización de depuración para exponer dependencias accidentales;
- `reserve` puede amortizar crecimiento, pero no justifica retener capacidad ilimitada ni asumir que `erase` la devuelve.

**Puerta de selección y benchmark.** Antes de migrar `unordered_map` o diseñar una tabla propia, registrar:

1. contrato de iteradores, punteros/referencias, `erase`, exceptions, allocator y frontera ABI/DSO;
2. tamaño y alineación de clave/valor, número de elementos y bytes residentes por entrada;
3. proporción real de hit/miss/insert/erase/iterate, distribución y skew de claves, churn y picos de `rehash`;
4. throughput y `p50/p95/p99/p99.9`, no sólo nanosegundos promedio; cycles, branches, cache/TLB misses y allocations;
5. hashes normales, débiles y claves controladas por adversario; equivalencia `Hash`/`Eq` mediante property tests;
6. comportamiento con tabla vacía, pequeña, cerca de su carga máxima, después de borrar masivamente y bajo memoria limitada;
7. sanitizers/hardening para use-after-rehash y referencias inválidas; carrera de datos aparte, porque estos contenedores no añaden sincronización al contrato.

**[IMPL]** Elegir la representación por evidencia del workload. Si se requiere dirección estable, considerar node storage o valores indirectos; si domina lookup y caben valores inline, medir flat/value; si se itera intensamente, medir por separado tras churn. Una mejora de microbenchmark no autoriza cambiar semántica, lifetime ni seguridad del servicio.

### 6.2 Árbol de búsqueda balanceado

Mantiene claves ordenadas y operaciones principales `O(log n)` en el peor caso cuando conserva su invariante de balance. AVL favorece balance más estricto; red-black reduce ciertas rotaciones; la biblioteca del lenguaje puede ocultar otra implementación. Lo importante es el contrato, no el nombre.

Operaciones adicionales mediante augmentación:

- rank/select con tamaño de subárbol;
- interval queries con máximo de extremo;
- sumas/rangos con agregados;
- min/max, predecessor y successor.

Cada campo agregado debe actualizarse en todas las mutaciones y rotaciones.

### 6.3 Trie, radix y búsquedas por prefijo

Útiles cuando el costo debe depender de longitud de clave y se requieren prefijos. Pueden consumir mucha memoria con alfabetos grandes. Aplicar compresión de caminos o representaciones compactas cuando corresponda. Para búsqueda textual avanzada, considerar suffix array/tree/automata según la consulta y el presupuesto.

### 6.4 Estructuras probabilísticas

Un Bloom filter puede responder “definitivamente ausente” o “posiblemente presente”; admite falsos positivos, no falsos negativos bajo su contrato básico sin borrado defectuoso. Nunca tratarlo como fuente de verdad. Declarar capacidad, tasa de falso positivo, función de hash y estrategia de rebuild.

## 7. Priority queues, heaps y selección

Un heap binario ofrece mínimo/máximo en `O(1)`, inserción y extracción en `O(log n)`, y construcción bottom-up en `O(n)`. No ofrece búsqueda arbitraria eficiente ni iteración ordenada.

Usos:

- top-`k` con heap acotado;
- scheduling;
- Dijkstra/Prim;
- merge de `k` secuencias;
- simulación por eventos.

Para decrease-key, definir cómo se localiza el elemento: handle/índice externo, duplicados perezosos o una estructura distinta. La API concreta modifica costo y corrección.

Selección:

| Objetivo | Opción |
|---|---|
| mínimo/máximo único | scan `Θ(n)` |
| `k` pequeño online | heap de tamaño `k` |
| k-ésimo esperado in-place | quickselect esperado `Θ(n)` |
| garantía lineal | median-of-medians, mayor constante |
| datos estáticos y muchas consultas | ordenar una vez |

## 8. Sorting y búsqueda

### 8.1 Decisión de sorting

| Algoritmo/familia | Tiempo | Espacio/propiedad | Usarlo cuando |
|---|---|---|---|
| insertion sort | `O(n²)`, cercano a lineal si casi ordenado | in-place, estable | bloques pequeños o casi ordenados |
| mergesort | `Θ(n log n)` | buffer típico; estable | estabilidad, linked/external merge |
| heapsort | `Θ(n log n)` | in-place; no estable | memoria estricta y peor caso |
| quicksort introspectivo | esperado `Θ(n log n)` | in-place; pivote crítico | uso general con defensa de peor caso |
| counting/radix | depende de rango/dígitos | no comparison-based | claves discretas explotables |
| external merge | I/O por bloques | archivos/runs | datos mayores que memoria |

Antes de ordenar, preguntar si basta selección, partición, bucket, índice o mantener orden incremental.

Propiedades que deben declararse:

- estabilidad;
- in-place real y memoria auxiliar;
- comparator total, consistente y sin efectos laterales;
- tratamiento de NaN, null, Unicode y locale;
- duplicados;
- peor caso y defensa ante adversario;
- costo de mover elementos grandes.

### 8.2 Binary search sin errores de borde

Definir intervalo y mantenerlo siempre. Para un rango semiabierto `[lo, hi)`:

```text
invariante: la respuesta, si existe, permanece dentro de [lo, hi)
while lo < hi:
    mid = lo + (hi - lo) // 2
    descartar una mitad sin violar el invariante
```

No limitar binary search a arrays: sirve sobre cualquier predicado monótono, siempre que evaluar el predicado y acceder al punto tengan costos declarados. Diferenciar exact match, lower bound, upper bound, first true y last false.

## 9. Union–find y conectividad incremental

Disjoint-set union mantiene una partición bajo `find` y `union`. Con union-by-rank/size y path compression, una secuencia de operaciones tiene costo amortizado casi constante, formalmente ligado a la inversa de Ackermann.

Aplicaciones:

- conectividad incremental;
- Kruskal;
- componentes en procesamiento offline;
- detección de ciclos no dirigidos.

No resuelve borrado dinámico general. Verificar inicialización, índices, unión repetida y si el representante observable forma parte del contrato.

## 10. Modelar con grafos

### 10.1 Preguntas previas

```text
dirigido o no dirigido
simple, multigrafo o self-loops
ponderado o no
pesos negativos
estático, incremental o dinámico
disperso o denso
un origen, varios o todos los pares
necesidad de camino, distancia, conectividad, orden o flujo
```

### 10.2 Recorridos

| Problema | Herramienta | Garantía principal |
|---|---|---|
| alcance y distancia sin pesos | BFS | menor cantidad de aristas |
| exploración/estructura | DFS | tiempos, forest y clasificación |
| DAG | topological sort | orden si y sólo si no hay ciclo dirigido |
| componentes no dirigidas | BFS/DFS/DSU | pertenencia a componente |
| componentes fuertemente conectadas | Kosaraju/Tarjan | partición dirigida |

BFS y DFS son `Θ(V + E)` con adjacency lists si se visitan vértices y aristas una vez. Marcar visitado en el momento correcto evita duplicación explosiva.

### 10.3 Caminos mínimos

| Condición | Algoritmo |
|---|---|
| aristas sin peso / peso unitario | BFS |
| pesos 0/1 | 0–1 BFS |
| pesos no negativos | Dijkstra |
| pesos negativos, detectar ciclos negativos alcanzables | Bellman–Ford |
| DAG ponderado | relajación en orden topológico |
| todos los pares, denso | Floyd–Warshall según escala |
| todos los pares, disperso y negativos sin ciclo negativo | Johnson |

La relajación correcta no habilita Dijkstra con pesos negativos. Predecessor y distancia necesitan estados distintos para “inalcanzable”, infinito y overflow.

### 10.4 Minimum spanning tree

Kruskal ordena aristas y usa DSU; Prim expande desde un árbol mediante prioridad. Ambos dependen de la propiedad de corte. Un MST minimiza suma total de aristas de un árbol, no distancias desde una raíz.

### 10.5 Flujo, corte y matching

Modelar capacidades, conservación y residual graph. Un augmenting path modifica aristas directas y reversas. Max-flow/min-cut conecta el valor del flujo máximo con la capacidad del corte mínimo bajo el modelo. Matching bipartito puede reducirse a flujo, pero algoritmos especializados pueden ser mejores.

Antes de usar la reducción, probar que toda solución factible de un lado corresponde a una solución del otro y que el objetivo se conserva.

## 11. Diseño por paradigmas

### 11.1 Divide and conquer

```text
dividir en subproblemas
resolver recursivamente
combinar
```

Probar que los subproblemas son menores, que cubren los casos y que la combinación preserva la solución. Medir si recursion, allocation y layout destruyen el beneficio teórico. Ejemplos: mergesort, FFT, selección y geometría.

### 11.2 Greedy

Una elección local necesita una propiedad global. Ruta de prueba:

1. caracterizar soluciones factibles;
2. proponer elección;
3. tomar una solución óptima arbitraria;
4. intercambiar sin empeorar el objetivo;
5. reducir al subproblema restante.

Si el intercambio falla, buscar DP, flujo, matroid structure o aproximación. Interval scheduling por earliest finish es un patrón; no generalizarlo a cualquier scheduling.

### 11.3 Dynamic programming

DP requiere subproblemas superpuestos y estructura óptima. Diseñar en este orden:

```text
estado = información mínima que determina el futuro
transición = decisiones y estados previos
base = instancias mínimas
orden = dependencias antes que consumidores
respuesta = estado o agregación final
```

Después:

- probar que toda solución induce una transición;
- probar que cada transición produce una solución válida;
- contar estados × transiciones por estado;
- reconstruir solución si no basta el valor;
- comprimir memoria sólo si no elimina datos necesarios;
- distinguir valor numérico de longitud en bits para evitar pseudopolinomialidad oculta.

### 11.4 Backtracking y branch-and-bound

Usar cuando se necesita búsqueda exacta y la explosión combinatoria es inevitable en el peor caso. Definir:

- orden de decisiones;
- poda por inviabilidad;
- bound optimista para poda por objetivo;
- memoization/symmetry breaking;
- límite de tiempo y mejor incumbent;
- respuesta parcial o “desconocido” distinta de “no existe”.

### 11.5 Randomization

Separar:

- **Las Vegas:** resultado siempre correcto; tiempo aleatorio.
- **Monte Carlo:** tiempo acotado; probabilidad de error controlada.

Declarar fuente de aleatoriedad, seed reproducible, adversario, probabilidad de fallo por operación y acumulada, y estrategia ante peor caso.

### 11.6 Incremental improvement

Parte de una solución factible y la mejora mediante una estructura residual o local. Max flow y matching son ejemplos. La ausencia de una mejora local sólo implica optimalidad si existe el teorema correspondiente.

## 12. Strings y secuencias

### 12.1 La unidad de texto importa

Bytes, code points, grapheme clusters y collation no son equivalentes. Un algoritmo correcto sobre ASCII puede cortar Unicode o comparar texto de forma inesperada. Normalización, locale y case folding forman parte del contrato.

### 12.2 Herramientas

| Necesidad | Familia |
|---|---|
| prefijos/diccionario | trie/radix |
| patrón único lineal | KMP/Z según interfaz |
| muchos patrones | Aho–Corasick |
| búsquedas repetidas sobre corpus estático | suffix array/tree/FM-index según escala |
| similitud/alineación | edit distance/LCS por DP |
| patrones regulares | autómata/engine con semántica declarada |
| compresión | Huffman/LZW/context models según datos y formato |

Evitar regex con backtracking no acotado sobre entrada adversarial. Diferenciar matching exacto, tokenización, ranking y semantic retrieval; los dos últimos se conectan con `NLP_RAG_RETRIEVAL_DATA.md`.

## 13. Optimización, reducciones y límites de tratabilidad

### 13.1 Linear programming

Formular variables, objetivo y restricciones. Verificar unidades, dominios e integridad. Una relajación lineal puede dar bound, no necesariamente solución entera. Solver status, tolerancias numéricas y gap forman parte del resultado.

### 13.2 Reducciones

Dos usos distintos:

- **algorítmico:** transformar una instancia a otra que sabemos resolver;
- **complejidad:** demostrar que resolver un problema permitiría resolver otro difícil.

Siempre incluir función de transformación, equivalencia de respuestas y costo de transformar/recuperar.

### 13.3 P, NP y NP-completeness

No interpretar “NP” como “no polinomial”. Para declarar NP-completo se necesita pertenencia a NP y reducción desde un problema NP-hard conocido en la dirección correcta. La clasificación asintótica no decide automáticamente qué instancias reales son resolubles.

### 13.4 Cuando exactitud total no escala

| Estrategia | Garantía que debe declararse |
|---|---|
| aproximación | ratio aditivo/multiplicativo |
| FPT | `f(k) · poly(n)` y parámetro `k` |
| heurística | ninguna universal; medir calidad/tiempo |
| anytime | calidad del incumbent y bound al interrumpir |
| relajación | bound y procedimiento de redondeo |
| muestreo/sketch | error y probabilidad |

Nunca presentar heurística como exacta ni un buen promedio como garantía de peor caso.

## 14. Datos mayores que memoria y locality

### 14.1 I/O y cache

Cuando mover datos cuesta más que operar sobre ellos, contar transferencias. Un algoritmo con más instrucciones y accesos secuenciales puede superar otro con mejor Big-O RAM pero locality deficiente.

Técnicas:

- blocking/tiling;
- external merge y runs;
- batching;
- layouts contiguos y structure-of-arrays según acceso;
- B-trees/LSM para persistencia;
- cache-oblivious divide-and-conquer cuando su análisis aplica.

MIT 6.046 dedica sus dos últimas clases a medianas, matrices, búsqueda y sorting cache-oblivious. El detalle físico y la medición pertenecen a `COMPUTER_SYSTEMS_PERFORMANCE_LOW_LATENCY.md`; persistencia a `DATABASE_STORAGE_INTERNALS.md`.

### 14.2 Streaming y online

Definir:

- una o varias pasadas;
- memoria acotada;
- orden y lateness;
- exactitud o error;
- mergeability de estados;
- ventanas y expiración;
- adversario y skew.

Reservoir sampling, sketches y heavy hitters son útiles sólo con sus garantías. Para event time, checkpoints y entrega, usar `NETWORKING_DISTRIBUTED_STREAMING.md`.

## 15. Ingeniería de implementación

### 15.1 Tipos y aritmética

- comprobar overflow/underflow antes de que ocurra;
- elegir signed/unsigned conscientemente;
- no usar floating point como igualdad exacta sin contrato;
- distinguir sentinel válido de ausencia;
- representar infinito sin que una suma haga wrap;
- validar conversiones, tamaños e índices;
- documentar precisión y tolerancias.

### 15.2 Recursión, memoria y ownership

- profundidad máxima y riesgo de stack overflow;
- allocations por elemento o transición;
- aliasing y mutación durante iteración;
- invalidación de iteradores/referencias;
- lifetime de keys/values almacenados;
- copias invisibles en slices, strings o closures;
- cleanup ante salida temprana y excepciones.

### 15.3 Determinismo y comparadores

Un comparator debe respetar el orden que exige la biblioteca. Evitar restar enteros para comparar si puede overflow. Si el output necesita reproducibilidad, definir tie-breaker total y no depender del orden accidental de un hash map.

### 15.4 Concurrencia cambia el problema

Una estructura correcta secuencial no es automáticamente thread-safe. Declarar:

- linearization point;
- atomicidad de operaciones compuestas;
- orden de memoria;
- reclamación segura de nodos;
- ABA, starvation y fairness;
- snapshot/iteration semantics.

El diseño lock-free y la memoria atómica se tratan en `COMPUTER_SYSTEMS_PERFORMANCE_LOW_LATENCY.md`.

### 15.5 Seguridad adversarial

- hash flooding;
- regex denial of service;
- nesting o recursion bombs;
- decompression bombs;
- inputs que disparan peor caso de quicksort;
- tamaños que desbordan multiplicaciones de allocation;
- grafos/strings malformados;
- consumo no acotado por output.

Todo input no confiable necesita límites de trabajo, memoria, profundidad y salida.

## 16. Estrategia de tests

### 16.1 Pirámide algorítmica

1. ejemplos mínimos y casos vacíos;
2. bordes de rango y duplicados;
3. invariantes internos en debug/test;
4. property-based testing;
5. comparación con oráculo lento;
6. metamorphic testing;
7. fuzzing y entradas adversariales;
8. stress de tamaño y memoria;
9. benchmark reproducible separado de correctness.

### 16.2 Propiedades útiles

| Dominio | Propiedad |
|---|---|
| sorting | permutación de entrada + orden total |
| search | índice válido y valor correcto; ausencia real |
| heap | raíz extrema + invariante padre/hijo |
| DSU | equivalencia reflexiva/simétrica/transitiva |
| shortest path | path válido; desigualdades; comparación con Bellman–Ford pequeño |
| MST | árbol, conecta, `V-1` aristas; cut property/oráculo pequeño |
| flow | capacidad, conservación, valor = corte certificado |
| DP | equivalencia con enumeración pequeña |
| parser/string | round-trip o semántica preservada |

### 16.3 Generación adversarial

Incluir:

- vacío, uno y dos elementos;
- todos iguales, todos distintos;
- ascendente, descendente, casi ordenado;
- máximos/mínimos representables;
- claves con colisiones;
- grafos desconectados, ciclos, self-loops y multiaristas;
- pesos cero, negativos y suma cercana al límite;
- strings vacíos, Unicode combinado y secuencias inválidas según API;
- inputs que maximizan profundidad, estados o output.

## 17. Benchmark y selección real

**[IMPL]** El análisis descarta diseños que no escalan; el benchmark decide entre implementaciones plausibles en el entorno objetivo.

Registrar:

```text
distribución y tamaño de entradas
hardware/runtime/compilador/flags
warm-up y número de repeticiones
tiempo, memoria, allocations y bytes movidos
p50/p95/p99 si hay variabilidad
correctness check del output
datos crudos y seed
```

Comparar crossover points, no un único tamaño. Mantener baseline y candidate intercalados. No usar microbenchmark para inferir throughput end-to-end sin confirmar colas, I/O y concurrencia.

## 18. Matriz de decisión rápida

| Señal del problema | Primera familia a evaluar | Pregunta de descarte |
|---|---|---|
| acceso por clave | hash/tree/trie | ¿se necesita orden, rango o adversarial guarantee? |
| extremo repetido | heap | ¿también se busca/borrar arbitrariamente? |
| conectividad incremental | DSU | ¿hay borrados? |
| dependencia por niveles | BFS/toposort | ¿hay pesos o ciclos? |
| costo acumulado en red | shortest path | ¿hay pesos negativos? |
| conectar con costo total | MST | ¿se confunde con rutas desde origen? |
| asignación/capacidad | flow/matching/LP | ¿integralidad y escala? |
| subproblemas repetidos | DP | ¿estado suficiente y acotado? |
| elección irreversible | greedy | ¿existe intercambio/cut property? |
| corpus estático, muchas búsquedas | índice/trie/suffix | ¿memoria y updates? |
| datos fuera de RAM | external/cache-aware | ¿cuál es el tamaño de bloque/working set? |
| exacto intratable | FPT/aprox/heurística | ¿qué garantía necesita negocio? |

## 19. Flujo de resolución para persona o Codex

### Fase A — entender

1. Reescribir el problema sin solución propuesta.
2. Enumerar casos y contradicciones del enunciado.
3. Identificar parámetros, operaciones y límites.
4. Preguntar qué significa igualdad, orden, error y timeout.

### Fase B — modelar

1. Elegir dominio: secuencia, set, grafo, intervalo, flujo, string, estado.
2. Crear baseline y ejemplos.
3. Establecer lower bound obvio o volumen mínimo de lectura/output.
4. Elegir estructura por operaciones.

### Fase C — demostrar

1. Escribir pseudocódigo sin detalles accidentales.
2. Fijar invariante y terminación.
3. Analizar tiempo/espacio con parámetros.
4. Buscar contraejemplo adversarial.

### Fase D — implementar

1. Usar tipos que representen el dominio.
2. Separar política de mecanismo.
3. Tratar overflow, errores y límites.
4. Añadir assertions de invariantes donde sean asequibles.

### Fase E — verificar

1. Tests deterministas.
2. Propiedades y diferencial contra baseline.
3. Fuzz/adversarial.
4. Sanitizers/static analysis pertinentes.
5. Benchmark sólo después de correctness.

### Fase F — entregar

1. Complejidad y supuestos visibles.
2. Casos límite documentados.
3. Resultados de tests y medición.
4. Umbral que obliga a cambiar de diseño.
5. Riesgos residuales y rollback si integra un sistema.

## 20. Contrato de prompt para Codex

```text
Problema y criterio de éxito:
Entradas, salidas y pre/postcondiciones:
Tamaños y distribuciones esperadas/adversariales:
Presupuesto de tiempo, memoria, I/O y latencia:
Mutabilidad, estabilidad, determinismo y concurrencia:
Lenguaje/runtime y estructuras permitidas:

Entregar:
1) supuestos y casos límite;
2) baseline/oráculo pequeño;
3) alternativas y trade-offs;
4) algoritmo elegido con pseudocódigo;
5) prueba de corrección y terminación;
6) tiempo/espacio bajo modelo declarado;
7) implementación sin pseudocódigo oculto;
8) tests unitarios, properties, diferencial y adversarial;
9) benchmark reproducible si afecta la decisión;
10) riesgos, límites y criterio de cambio.

No afirmar optimalidad sin límite inferior aplicable.
No inferir rendimiento físico únicamente desde Big-O.
No usar una estructura cuyo contrato no cubra las operaciones requeridas.
```

## 21. Señales de una revisión senior

Rechazar o devolver para corrección cuando:

- el algoritmo no corresponde al problema especificado;
- faltan precondiciones o semántica de error;
- el proof sketch no cubre inicialización/mantenimiento/terminación;
- se usa promedio donde se necesita garantía de tail;
- se oculta otro parámetro dominante;
- una estructura ofrece peor operación crítica que una alternativa simple;
- se desconoce la semántica de comparator/hash/equality;
- no existe oráculo o propiedad independiente;
- benchmark y producción usan distribuciones distintas;
- la optimización degrada claridad sin beneficio medido;
- “no encontrado”, overflow y timeout comparten el mismo sentinel;
- el estado puede crecer sin límite.

## 22. Capstones de validación

### A. Índice mutable

Implementar lookup, inserción, borrado, range query y top-`k`. Comparar hash + heap, árbol augmentado y sorted vector bajo tres distribuciones. Exigir invariantes, property tests, crossover de lectura/escritura y memoria por elemento.

### B. Motor de rutas

Soportar grafos dirigidos, desconectados y distintos regímenes de peso. Seleccionar BFS/DAG/Dijkstra/Bellman–Ford según contrato, reconstruir caminos y generar certificados/test diferencial en grafos pequeños.

### C. Scheduler con restricciones

Comenzar con interval scheduling, demostrar cuándo greedy es correcto y construir un contraejemplo cuando se agregan pesos. Resolver la variante ponderada con DP y medir reconstrucción/memoria.

### D. Procesador de stream acotado

Top-`k`, deduplicación aproximada y ventanas. Declarar error, memoria, orden, timestamps y comportamiento ante skew. Conectar checkpoint/entrega con el manual de streaming.

### E. Order book local

Aplicar mapas ordenados, heaps sólo donde sean semánticamente correctos, colas FIFO por price level y secuencias. Verificar invariantes de precio-tiempo, snapshot/delta y gaps con `MARKET_MICROSTRUCTURE_EXCHANGE_SYSTEMS.md`.

### F. Corpus de búsqueda

Comparar trie/radix, hash e índice invertido para exact/prefix/token lookup; medir Unicode, memoria, updates y queries. La recuperación semántica se remite a `NLP_RAG_RETRIEVAL_DATA.md`.

## 23. Auditoría curricular de fuentes

### 23.1 MIT 6.006 Spring 2020 — 20/20 clases cubiertas

| # | Tema oficial | Receptor principal |
|---:|---|---|
| 1 | Introduction | §§1–4 |
| 2 | Data Structures | §§5–7 |
| 3 | Sorting | §8 |
| 4 | Hashing | §6 |
| 5 | Linear Sorting | §8 |
| 6 | Binary Trees 1 | §6 |
| 7 | Binary Trees 2: AVL | §6 |
| 8 | Binary Heaps | §7 |
| 9 | Breadth-First Search | §10 |
| 10 | Depth-First Search | §10 |
| 11 | Weighted Shortest Paths | §10 |
| 12 | Bellman–Ford | §10 |
| 13 | Dijkstra | §10 |
| 14 | All-pairs Shortest Paths/Johnson | §10 |
| 15 | Dynamic Programming 1 | §11 |
| 16 | DP 2: LCS, LIS, Coins | §§11–12 |
| 17 | Dynamic Programming 3 | §11 |
| 18 | DP 4: Rods, Subset Sum, Pseudopolynomial | §§11, 13 |
| 19 | Complexity | §13 |
| 20 | Review | §§18–21 |

### 23.2 MIT 6.046J Spring 2015 — 24/24 clases cubiertas

| Rango | Temas oficiales | Receptor |
|---|---|---|
| 1 | overview, interval scheduling | §§1, 11 |
| 2–4 | divide-and-conquer: hull, median, FFT, van Emde Boas | §§7, 11; geometría/FFT como extensión |
| 5 | amortized analysis | §3 |
| 6–8 | randomization, quicksort, skip lists, universal/perfect hashing | §§6, 8, 11 |
| 9 | augmentation/range trees | §6 |
| 10–11 | advanced DP, all-pairs shortest paths | §§10–11 |
| 12 | greedy/MST | §§10–11 |
| 13–14 | max-flow/min-cut, matching | §10 |
| 15 | LP, reductions, simplex | §13 |
| 16–18 | P/NP, approximation, fixed-parameter | §13 |
| 19–20 | synchronous/asynchronous distributed algorithms | frontera con manual distribuido |
| 21–22 | hash functions, encryption | frontera con manual de seguridad |
| 23–24 | cache-oblivious median/matrix/search/sort | §14 y manual de sistemas |

### 23.3 Princeton *Algorithms, 4th Edition* — 6/6 áreas cubiertas

| Área oficial | Receptor |
|---|---|
| Fundamentals: model, abstraction, stacks/queues, analysis, union–find | §§1–5, 9 |
| Sorting: elementary, merge, quick, priority queues, applications | §§7–8 |
| Searching: symbol tables, BST, balanced trees, hashing | §6 |
| Graphs: undirected/directed, MST, shortest paths | §10 |
| Strings: sorts, tries, substring, regex, compression | §12 |
| Context: simulation, B-trees, suffix arrays, maxflow, reductions, intractability | §§7, 10, 12–14 y manuales de storage |

La cobertura significa que cada tema tiene receptor conceptual; no reproduce clases, soluciones, código ni texto protegido. Geometría computacional, FFT, van Emde Boas, criptografía y algoritmos distribuidos conservan receptor o frontera explícita y se profundizan sólo cuando el proyecto los requiera.

## 24. Cruces con la biblioteca

| Manual | Relación |
|---|---|
| `COMPUTER_SYSTEMS_PERFORMANCE_LOW_LATENCY.md` | costo físico, cache, layout, SIMD, concurrencia y benchmarks |
| `SOFTWARE_BACKEND_API_ENGINEERING.md` | contratos, validación, errores, testing y diseño de componentes |
| `NETWORKING_DISTRIBUTED_STREAMING.md` | grafos/colas, orden, ventanas, consenso y entrega |
| `DATABASE_STORAGE_INTERNALS.md` | B-trees/LSM, joins, query planning, external algorithms |
| `GPU_ACCELERATED_COMPUTING.md` | trabajo/span, tiling, parallel primitives y costo de transferencia |
| `SECURITY_SRE_CLOUD_INFRASTRUCTURE.md` | adversarios, resource exhaustion, límites y operación |
| `MARKET_MICROSTRUCTURE_EXCHANGE_SYSTEMS.md` | price levels, queues, secuencias, matching y métricas |
| `FRONTEND_PRODUCT_ENGINEERING_UX.md` | state/data structures, rendering, búsqueda, workers y budgets |
| `AI_ENGINEERING_MASTER_MAP.md` | optimización, muestreo, grafos, retrieval, entrenamiento e inferencia |

Cruces recíprocos cerrados el 2026-08-21: cada manual contiene ahora su obligación algorítmica y conserva autoridad sobre su semántica de dominio. Este manual decide modelado, corrección y costo abstracto; sistemas decide costo físico, y cada dominio decide qué resultado es válido.

Cruce con `SOFTWARE_ARCHITECTURE_SYSTEM_DESIGN.md` cerrado el 2026-08-21: este manual prueba el algoritmo y su costo; arquitectura decide boundaries, ownership, quality scenarios, distribución y evolución. Un diseño local óptimo puede ser sistémicamente incorrecto si rompe consistency, operability o change isolation; una capa arquitectónica tampoco justifica complejidad algorítmica evitable.

Cruce con `DATA_ENGINEERING_ANALYTICS.md` cerrado el 2026-08-21: algoritmos decide exactitud/error y costo de joins, aggregations, sketches, sampling, graph/stream/external processing; datos decide source contract, event time, lineage, publication, quality y reproducibilidad. Un sketch o aproximación porta bound/probabilidad y no puede presentarse como dato exacto; un pipeline distribuido no cambia la prueba de la transformación.

## 25. Fuentes oficiales verificadas

### MIT

- MIT 6.006 Spring 2020, curso: <https://ocw.mit.edu/courses/6-006-introduction-to-algorithms-spring-2020/>
- MIT 6.006 Spring 2020, 20 lecture notes: <https://ocw.mit.edu/courses/6-006-introduction-to-algorithms-spring-2020/pages/lecture-notes/>
- MIT 6.006 Fall 2011, syllabus y contrato de respuesta: <https://www.ocw.mit.edu/courses/6-006-introduction-to-algorithms-fall-2011/pages/syllabus/>
- MIT 6.046J Spring 2015, curso: <https://ocw.mit.edu/courses/6-046j-design-and-analysis-of-algorithms-spring-2015/>
- MIT 6.046J Spring 2015, 24 lecture notes: <https://ocw.mit.edu/courses/6-046j-design-and-analysis-of-algorithms-spring-2015/pages/lecture-notes/>

### Princeton

- Sedgewick y Wayne, *Algorithms, 4th Edition*, sitio oficial: <https://algs4.cs.princeton.edu/home/>
- Analysis of Algorithms y método experimental: <https://algs4.cs.princeton.edu/14analysis/>

### Implementaciones industriales verificadas

- Abseil C++ Containers Guide: <https://abseil.io/docs/cpp/guides/container>
- Abseil Swiss Tables, snapshot auditado: <https://github.com/abseil/abseil-cpp/tree/a690167a55b9e210b7884d311759977c0f5669c9>
- Folly F14 design, snapshot auditado: <https://github.com/facebook/folly/blob/8dc549c54875146647295331282db573c586831e/folly/container/F14.md>
- Folly, snapshot auditado: <https://github.com/facebook/folly/tree/8dc549c54875146647295331282db573c586831e>

## 26. Puerta de cierre del manual

- [x] Autoridades y alcance definidos.
- [x] Contrato de problema y de solución.
- [x] Corrección, complejidad y modelos de costo.
- [x] Estructuras fundamentales y avanzadas de uso general.
- [x] Sorting, searching, grafos, strings y paradigmas.
- [x] Complejidad, aproximación, FPT, external/streaming.
- [x] Implementación, seguridad, tests y benchmark.
- [x] Inventario 20/20 + 24/24 + 6/6 con receptor.
- [x] Auditoría de consistencia interna automatizada: fences, encabezados y referencias.
- [x] Cruces recíprocos con los manuales afectados.
- [x] Contrato y capstones convertidos en planes ejecutables mediante `ENGINEERING_EXECUTION_PLAYBOOK.md`; ninguna implementación se marca `passed` sin evidencia de proyecto.

El conocimiento general v1 queda auditado. No declarar los capstones `ejecutables` hasta que sus repositorios, tests y mediciones existan realmente.
