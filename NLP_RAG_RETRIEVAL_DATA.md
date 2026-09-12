# NLP, RAG, retrieval y datos — manual operativo

> **Estado:** *Natural Language Processing Specialization* y *Retrieval Augmented Generation (RAG)* completos y auditados de punta a punta.
> **Fuentes primarias actuales:** Younes Bensouda Mourri, Łukasz Kaiser, Zain Hasan y una conversación identificada de Andrew Ng; DeepLearning.AI.
> **Cobertura inventariada:** NLP —4 cursos, 14 semanas, 412 unidades, 176 videos, 38 ejemplos de código y 20 asignaciones—; RAG —5 módulos, 77 unidades, 49 videos, 9 ejemplos de código y 10 evaluaciones certificables—.
> **Propósito:** convertir NLP clásico y neuronal en decisiones, algoritmos, pruebas y sistemas de representación/recuperación que complementen Transformers, RAG y agentes sin duplicar los otros manuales.

## Cómo usar este documento

1. Consultar el mapa curricular para ubicar el problema.
2. Cargar sólo la sección del método requerido.
3. Para redes/Transformers, usar además `DEEP_LEARNING_ANDREW_NG.md` o `PYTORCH_TRANSFORMERS_GENERATIVE_MODELS.md`; aquí se conservará el delta de NLP y retrieval.
4. Para despliegue, drift y observabilidad, usar `ML_PRODUCTION_LLMOPS_EVALUATION.md`.
5. Ninguna sección marcada pendiente debe tratarse como instrucción ya verificada.

## Procedencia y límites

- **[DLAI]** paráfrasis fiel de Younes Bensouda Mourri, Łukasz Kaiser o contenido oficial del programa; no es transcripción literal.
- **[NG]** reservado para intervenciones directamente identificadas de Andrew Ng.
- **[IMPL]** traducción a especificación, código, prueba, métrica o proceso profesional.
- **[EXT]** práctica externa al programa, identificada expresamente.
- Las asignaciones se inventarían pero no se reproducen ni resuelven. Los ejemplos de código se convierten en habilidades y contratos.
- La especialización es histórica y pedagógica: APIs, TensorFlow/Trax, modelos, benchmarks y recomendaciones se modernizarán sin atribuir las extensiones al curso.

## Inventario maestro verificado

### Natural Language Processing Specialization

| Curso | Semanas | Videos | Código visible | Lecturas visibles | Unidades totales | Asignaciones |
|---|---:|---:|---:|---:|---:|---:|
| 1. Classification and Vector Spaces | 4 | 47 | 9 | 39 | 113 | 5 |
| 2. Probabilistic Models | 4 | 57 | 12 | 49 | 133 | 6 |
| 3. Sequence Models | 3 | 33 | 8 | 30 | 83 | 6 |
| 4. Attention Models | 3 | 39 | 9 | 19 | 83 | 3 |
| **Total** | **14** | **176** | **38** | **137** | **412** | **20** |

Los 351 elementos tipados como video/código/lectura más 61 notas, quizzes, entrevistas, divisores curriculares y asignaciones suplementarias forman las 412 unidades. La diferencia respecto del primer inventario se debe a dos secciones colapsadas: *N-grams vs. Sequence Models* (Curso 3/Semana 1) y *Hugging Face* (Curso 4/Semana 3). Ambas se abrieron en modo lectura y sus unidades quedaron contabilizadas. La cuenta permaneció en 0 %; `Free/Pro` se usó sólo para comprobar el inventario, no para alterar el plan ni marcar lecciones.

### Curso 1 — Classification and Vector Spaces

| Semana | Tema | Videos | Código | Lecturas |
|---:|---|---:|---:|---:|
| 1 | Sentiment Analysis with Logistic Regression | 14 | 3 | 12 |
| 2 | Sentiment Analysis with Naïve Bayes | 13 | 1 | 11 |
| 3 | Vector Space Models | 10 | 3 | 9 |
| 4 | Machine Translation and Document Search | 10 | 2 | 7 |

**Ruta de extracción:** representación sparse → baselines lineales/probabilísticos → error analysis → embeddings/analogías → alignment de espacios → LSH/ANN y búsqueda documental.

### Curso 2 — Probabilistic Models

| Semana | Tema | Videos | Código | Lecturas |
|---:|---|---:|---:|---:|
| 1 | Autocorrect | 11 | 2 | 9 |
| 2 | POS Tagging and Hidden Markov Models | 13 | 2 | 11 |
| 3 | Autocomplete and Language Models | 11 | 3 | 9 |
| 4 | Word Embeddings with Neural Networks | 22 | 5 | 20 |

**Ruta de extracción:** edit distance/dynamic programming → HMM/Viterbi → n-grams/smoothing/perplexity → CBOW/skip-gram y evaluación de embeddings.

### Curso 3 — Sequence Models

| Semana | Tema | Videos | Código | Lecturas |
|---:|---|---:|---:|---:|
| 1 | RNNs for Language Modeling | 15 | 4 | 14 |
| 2 | LSTMs and Named Entity Recognition | 8 | 1 | 8 |
| 3 | Siamese Networks | 10 | 3 | 8 |

**Ruta de extracción:** dense sentiment/RNN language modeling → GRU/LSTM y NER → metric learning, triplet loss y duplicate-question retrieval.

### Curso 4 — Attention Models

| Semana | Tema | Videos | Código | Lecturas |
|---:|---|---:|---:|---:|
| 1 | Neural Machine Translation | 14 | 3 | 3 |
| 2 | Text Summarization | 10 | 3 | 5 |
| 3 | Question Answering | 15 | 3 | 11 |

**Ruta de extracción:** encoder–decoder/attention → Transformer causal y summarization → transfer learning, T5/BERT, QA y Hugging Face.

### Retrieval Augmented Generation (RAG)

Curso intermedio de DeepLearning.AI, instructor Zain Hasan, duración publicada de 26 h 03 min. La página oficial declara 49 videolecciones, 9 ejemplos de código y 10 evaluaciones. La inspección del syllabus añade 9 lecturas y confirma **77 unidades únicas**; evaluaciones inventariadas, no abiertas ni reproducidas.

| Módulo | Tema | Videos | Código | Lecturas | Evaluaciones | Total |
|---:|---|---:|---:|---:|---:|---:|
| 1 | RAG Overview | 8 | 2 | 2 | 2 | 14 |
| 2 | Information Retrieval and Search Foundations | 10 | 2 | 1 | 2 | 15 |
| 3 | Information Retrieval with Vector Databases | 9 | 2 | 1 | 2 | 14 |
| 4 | LLMs and Text Generation | 11 | 2 | 1 | 2 | 16 |
| 5 | RAG Systems in Production | 11 | 1 | 4 | 2 | 18 |
| **Total** |  | **49** | **9** | **9** | **10** | **77** |

**Ruta de extracción:** arquitectura y casos → lexical/semantic/hybrid retrieval y métricas → ANN/vector DB/chunking/query parsing/reranking → prompting/generación/agentic RAG → evaluación, trazas, costo, latencia, seguridad y multimodalidad.

## Fronteras contra duplicación

| Tema | Autoridad principal | Aquí se incorporará |
|---|---|---|
| regresión logística, redes, RNN/LSTM y optimización | ML/Math y Deep Learning | features textuales, supuestos lingüísticos y fallos |
| atención/Transformers, fine-tuning y serving | PyTorch/Transformers | traducción, resumen, QA y evaluación específica |
| scoping, datasets, drift y producción | ML Production | métricas/slices propios de NLP y retrieval |
| agentes y tool use | Agentic AI | recuperación, representaciones, corpus y grounding |
| RAG/índices/vector databases | este documento | pipeline completo, evaluación y datos |

# Parte I — Classification and Vector Spaces

## 1. Semana 1 — Sentiment Analysis with Logistic Regression

### El baseline que hay que construir antes de una red

**[DLAI]** Logistic regression es rápida, interpretable y fácil de entrenar, por lo que ofrece un baseline fuerte para sentimiento binario. El workflow enseñado es: procesar texto, extraer features, entrenar minimizando costo y evaluar sobre datos no usados en entrenamiento.

**[IMPL]** Un baseline simple prueba cinco cosas antes de pagar un modelo complejo:

1. si la etiqueta es aprendible desde el texto disponible;
2. si el split y la métrica representan producción;
3. qué tokens/patrones aportan señal o leakage;
4. qué errores requieren contexto/orden/semántica;
5. cuánto lift real debe justificar el siguiente modelo.

Nunca saltar directo a Transformer sin guardar este resultado, latencia, tamaño, costo y errores. En routing, moderación ligera o alto throughput, el baseline puede ser el modelo final o un filtro/cascade.

### Representación sparse y compresión por clase

**[DLAI]** Una representación binaria bag-of-words tiene dimensión `|V|` y muchos ceros. El curso propone un diccionario `freq[(word,class)]` y comprime cada texto a tres features: bias, suma de frecuencias positivas y suma de frecuencias negativas.

```text
freq[(w,c)] = número de apariciones de w en documentos train de clase c

x(text) = [
  1,
  Σ_{w ∈ preprocess(text)} freq[(w,1)],
  Σ_{w ∈ preprocess(text)} freq[(w,0)]
]

z = θᵀx
p(y=1|x) = σ(z) = 1/(1+e^-z)
```

**[IMPL]** Construir vocabulario/frecuencias **sólo con train** y aplicar el transform congelado a dev/test. Si se cuentan labels de todo el corpus, se filtra la respuesta. Decidir y testear si se suman ocurrencias o palabras únicas: el video alterna lenguaje que puede confundirse; son features distintas.

La compresión a tres dimensiones acelera y regulariza, pero pierde identidad, orden, negación, intensidad y contexto. Dos textos con sumas iguales se vuelven indistinguibles. Comparar con:

- count/binary bag-of-words sparse;
- TF-IDF y n-gramas de palabras/caracteres;
- log-count ratios o Naïve Bayes–SVM;
- embeddings/modelos neuronales sólo si el error residual lo exige.

Normalizar conteos (`log1p`, longitud, TF-IDF) si textos largos dominan. Inspeccionar contribuciones de tokens y ablar cada familia de features.

### Preprocesamiento es una hipótesis, no limpieza universal

**[DLAI]** El ejemplo elimina stopwords, puntuación, URLs y handles; convierte a minúscula y aplica stemming para reducir vocabulario. El propio curso advierte que puntuación puede ser informativa según la tarea.

**[IMPL] Matriz de decisión:**

| Transformación | Puede ayudar | Puede destruir |
|---|---|---|
| lowercase | sparsity/variantes | siglas, entidades, énfasis |
| stopword removal | modelos de conteo | negación, modalidad, estilo |
| punctuation removal | ruido de formato | emoticons, intensidad, límites |
| stemming | variantes morfológicas | legibilidad y distinciones léxicas |
| borrar URL/handle | privacidad/sparsity | fuente, comunidad, reputación |
| normalizar emoji | cobertura | matiz si el mapping es pobre |

Versionar `preprocess(text)`, aplicarlo idénticamente online/offline y medir ablations. No asumir que handles/URLs “no añaden valor”; pueden ser señal legítima o leakage/proxy indebido. Para Transformers subword, stemming y stopword removal suelen ser contraproducentes: evaluar el pipeline propio del checkpoint.

Pruebas: Unicode, emojis, hashtags, contracciones, `not good`, URLs, handles, HTML, texto vacío, repetición, idioma mixto y caracteres adversariales. No mutar el texto original requerido para auditoría.

### Entrenamiento y numerics

**[DLAI]** Logistic regression optimiza binary cross-entropy mediante gradient descent; el costo penaliza con fuerza predicciones confiadas y equivocadas.

```text
J(θ) = -(1/m) Σ_i [y_i log p_i + (1-y_i) log(1-p_i)]
∇J(θ) = (1/m) Xᵀ(p-y)
θ ← θ - α∇J(θ)
```

**[IMPL]** Calcular la loss desde logits con una primitiva estable (`logaddexp`/`BCEWithLogitsLoss`), no `log(sigmoid(z))` ingenuo. Estandarizar features de gran escala o ajustar optimizer/LR; añadir regularización si corresponde. Unit tests:

- gradient check por diferencias finitas;
- loss finita para logits extremos;
- loss decrece en dataset diminuto y separable;
- permutar labels destruye rendimiento;
- serializar/cargar reproduce probabilidades;
- inferencia batch y ejemplo individual coinciden.

La transcripción atribuye `4,92` a la salida de sigmoide; es imposible porque `σ(z) ∈ (0,1)`. Interpretarlo como **logit** `z=4,92`, cuya probabilidad sí es cercana a 1.

### Evaluación que representa el costo real

**[DLAI]** El curso obtiene clases con un threshold, usualmente `0,5`, y calcula accuracy en un conjunto separado.

**[IMPL]** `0,5` sólo corresponde a costos/prior/calibración particulares. Elegir threshold en dev según matriz de costo, capacidad humana o restricción de precision/recall; congelarlo antes de test. Reportar:

- confusion matrix y soporte por clase/slice;
- precision, recall, F1 y PR-AUC para desbalance;
- ROC-AUC cuando sea pertinente;
- log loss/Brier y reliability plot si se consumen probabilidades;
- cobertura/abstención, latencia, memoria y throughput;
- intervalos de confianza y comparación pareada con baseline.

Split por usuario, hilo/producto, fuente y tiempo para impedir duplicados/paráfrasis entre conjuntos. Accuracy de 80 % no significa que cada request tenga “80 % de probabilidad de funcionar”; es una frecuencia estimada bajo ese dataset/protocolo.

### Error analysis de sentimiento

Clasificar una muestra de fallos por: negación/scope; sarcasmo/ironía; intensidad; emoji/puntuación; ambigüedad; múltiples aspectos; contexto conversacional; OOV/morfología; idioma/dialecto; label noise; spam; cambio temporal; entidad/fuente. Para cada categoría registrar frecuencia, severidad, ejemplo, causa y cambio mínimo propuesto.

No agregar complejidad por un error anecdótico. Priorizar `frecuencia × costo`, corregir primero labels/split/preprocesamiento y volver a ejecutar toda regresión. Guardar falsos positivos y negativos críticos como test cases permanentes.

### Contrato de implementación de sentimiento

```text
Entradas: texto UTF-8 + modelo/preprocess/frequency revisions
Artefactos: preprocess config, vocabulary/freq train-only, θ, threshold, schema
Tests: leakage, transforms, numerics, parity batch/online, slices, adversariales
Métricas: task + calibration + operación
Salida: {label, probability, model_version}; explicación sólo como evidencia de features
Fallback: abstener/escalar cuando confidence y cobertura estén fuera de contrato
```

**Cobertura de la semana:** 14 videos/transcripciones contrastados, 12 lecturas y 3 ejemplos de código inventariados; práctica/quiz/asignación no reproducidos.

## 2. Semana 2 — Sentiment Analysis with Naïve Bayes

### Por qué conservar este baseline

**[DLAI]** Naïve Bayes resuelve el mismo problema de sentimiento con conteos y probabilidades. Es un baseline rápido, interpretable y transferible a clasificación de texto, spam, atribución de autor, recuperación y desambiguación. Se llama *naïve* porque supone independencia condicional entre features/palabras, algo rara vez verdadero en lenguaje.

**[IMPL]** Su valor profesional no es competir universalmente con redes: establece en minutos si la señal léxica basta, deja una puntuación descomponible por token y crea un control de regresión barato. Mantenerlo como:

- referencia de calidad, costo, memoria y latencia;
- detector/filtro temprano en un cascade;
- diagnóstico de priors, vocabulario y drift léxico;
- fallback reproducible cuando el modelo principal no está disponible.

Si un modelo complejo no mejora slices importantes frente a este baseline, la complejidad no está justificada.

### Modelo multinomial binario

Para clases `c ∈ {0,1}`, vocabulario `V`, conteo de tokens `N_c` y suavizado `α > 0`:

```text
P(c) = documentos_train_de_clase_c / documentos_train

P(w|c) = (count_train(w,c) + α) / (N_c + α|V|)

λ(w) = log P(w|1) - log P(w|0)
log_prior = log P(c=1) - log P(c=0)

score(text) = log_prior + Σ_w count(w,text) · λ(w)
ŷ = 1[score > τ]
```

**[DLAI]** El curso usa Laplace/add-one (`α=1`), suma las razones en espacio logarítmico para evitar underflow y usa `τ=0` en el ejemplo binario. Una palabra fuera del vocabulario no aporta al score. Entrenar aquí significa contar y estimar; no hay gradient descent.

**[IMPL]** `score` es una **log-odds bajo los supuestos del modelo**, no una probabilidad calibrada. `τ=0` sólo aplica a la regla/costos implícitos; seleccionar otro threshold en dev cuando el costo de errores, la prevalencia o la capacidad operativa lo requieran. Si se necesita `P(y|x)`, medir calibración y calibrar sobre datos separados.

La suma debe incluir repeticiones en Multinomial NB; Bernoulli NB usa presencia/ausencia y otro likelihood. Elegir la variante con una ablation, no mezclarlas accidentalmente. El vocabulario, los conteos, `N_c`, el prior y cualquier selección de features se estiman **sólo con train**.

### Suavizado, priors y estabilidad

**[DLAI]** Sin suavizado, una sola probabilidad cero anula un producto completo. Add-one agrega uno a cada conteo y `|V|` al denominador. El prior importa cuando las clases están desbalanceadas.

**[IMPL]** Generalizar add-one a `α` y elegirlo en dev. Add-one es una opción pedagógica, no una constante universal. Comprobar:

- cada distribución `Σ_{w∈V} P(w|c)` suma aproximadamente uno;
- no existen logs de cero, `NaN` o infinitos;
- `α`, tokenización y vocabulario están versionados junto al modelo;
- el prior representa el entorno objetivo o se corrige explícitamente ante *prior shift*;
- los OOV, texto vacío y tokens sólo presentes en una clase tienen comportamiento definido.

Usar `log P(w|1) - log P(w|0)` directamente. No multiplicar cientos de probabilidades en coma flotante. Para scores extremos, evitar convertir innecesariamente con `exp`; las comparaciones funcionan en log-space.

### Pipeline reproducible

```text
fit(train):
  validar labels y split
  ajustar/versionar preprocess con train
  construir V y count(w,c) con train
  calcular log_prior y λ(w) con α
  elegir α/variantes/features usando dev
  congelar artefactos y threshold

predict(text):
  tokens = preprocess_versionado(text)
  score = log_prior + suma de λ[token] por ocurrencia conocida
  devolver score, decisión, versión y cobertura léxica

evaluate(test congelado):
  métricas globales + slices + calibración + operación
  comparar pareado contra regresión logística y baseline trivial
```

Persistir schema, vocabulario/hash, tokenizer, conteos o lambdas, prior, `α`, threshold y corpus lineage. Verificar paridad batch/online, determinismo y compatibilidad hacia atrás. Para alto throughput, convertir tokens a IDs, almacenar lambdas en un array contiguo y sumar sin crear diccionarios/copias por request; medir antes de optimizar.

### Supuestos y límites que deben quedar explícitos

**[DLAI]** El curso muestra que palabras correlacionadas violan independencia; el orden cambia significado; una distribución balanceada artificialmente puede no representar el flujo real. Los modelos de conteo pierden negación, contexto, sarcasmo, ironía y eufemismos.

**[IMPL] Correcciones necesarias:**

- El desbalance de train no obliga a balancear para “ser correcto”: se debe estimar/tratar el prior y evaluar sobre una distribución que corresponda al uso o a un protocolo de costos explícito.
- Naïve Bayes sí tiene decisiones ajustables: `α`, variante multinomial/Bernoulli/complement, preprocesamiento, vocabulario, n-gramas, selección de features y threshold.
- Un OOV con contribución cero significa “sin evidencia aprendida”, no necesariamente neutralidad semántica.
- Una alta magnitud de `λ(w)` indica asociación en ese corpus; no causalidad, sentimiento intrínseco ni explicación humana completa.
- La independencia falsa puede duplicar evidencia correlacionada y producir scores sobreconfiados.

Antes de abandonar el baseline, probar bigramas/char n-grams para negación, ortografía y morfología; Complement NB para ciertas colecciones desbalanceadas; límites de vocabulario/min-frequency; y features que preserven emoji/puntuación. Validar cada cambio contra latencia, memoria y slices, no sólo accuracy agregada.

### Error analysis: aprender de lo que el algoritmo rompe

**[DLAI]** La instrucción práctica central es inspeccionar **cómo quedó realmente el texto procesado**. Eliminar puntuación puede borrar un emoticono decisivo; quitar stopwords puede transformar una frase negativa en una bolsa de palabras positivas. El orden y el alcance de `not`, además de sarcasmo/ironía/eufemismo, explican fallos que el supuesto bag-of-words no representa.

**[IMPL] Bucle obligatorio:**

1. Muestrear falsos positivos, falsos negativos, baja cobertura y scores extremos; incluir errores de alto costo.
2. Guardar texto original, tokens resultantes, contribución `count(w,x)·λ(w)`, prior y decisión.
3. Etiquetar causa dominante: label/ambigüedad, split/leakage, transformación destructiva, negación/orden, OOV, dominio/idioma, correlación repetida, sarcasmo/pragmática o prior shift.
4. Cuantificar `frecuencia × severidad`; no diseñar por una anécdota.
5. Proponer el cambio mínimo dirigido a la categoría: corregir datos, preservar token, añadir n-grama, recalibrar prior/threshold o escalar de familia de modelo.
6. Convertir ejemplos representativos en tests de regresión; reejecutar métricas globales y slices para detectar trade-offs.

No editar el conjunto de test final durante este ciclo. Usar train/dev y una cola de errores fechada; abrir el test congelado sólo al cerrar una decisión. Si la causa es dependencia/contexto, registrar que el límite es estructural: aumentar conteos no lo arregla necesariamente.

### Contrato de pruebas

- `P(w|c)` normaliza por clase y permanece finita con conteos cero.
- Un token repetido suma repetidamente sólo en la variante multinomial.
- OOV/texto vacío devuelven el prior y una cobertura explícita.
- Corpus balanceado produce `log_prior≈0`; corpus desbalanceado coincide con el cálculo manual.
- Un ejemplo diminuto calculado a mano coincide con el score implementado.
- Alterar test no cambia vocabulario, lambdas ni threshold.
- Permutar labels destruye la señal; duplicar documentos entre splits dispara el control de leakage.
- Unicode, emojis, contracciones, negación, URLs, handles y mezcla de idiomas tienen casos de regresión.
- Se reportan confusion matrix, precision/recall/F1, PR-AUC según necesidad, costo y slices; accuracy sola no basta.

**Criterio de salida:** baseline reproducible, comparado con regresión logística bajo el mismo split/preprocesamiento, con errores priorizados y una decisión documentada: conservar, usar en cascade o reemplazar por una representación contextual.

**Cobertura de la semana:** 13 videos/transcripciones contrastados, 11 lecturas y 1 ejemplo de código inventariados; práctica/quiz/asignación no reproducidos.

## 3. Semana 3 — Vector Space Models

### De texto a geometría útil

**[DLAI]** Todo sistema NLP necesita una representación numérica. Los modelos vectoriales explotan la idea distribucional: las palabras que aparecen en contextos parecidos tienden a adquirir representaciones parecidas. El curso construye espacios *word-by-word* y *word-by-document*, compara vectores con distancia euclídea/coseno, usa aritmética para relaciones y aplica PCA para visualizar.

**[IMPL]** Un vector no “contiene el significado” de forma completa. Codifica regularidades del corpus, objetivo, contexto y preprocesamiento que lo produjeron. La similitud es válida sólo dentro de ese contrato. Antes de usar embeddings en búsqueda, clustering o features, fijar:

```text
unidad representada: token | palabra | frase | documento | usuario/item
contexto: ventana, documento, sesión o secuencia
objetivo: conteo, predicción, contraste o tarea supervisada
normalización: ninguna | L1/L2 | TF-IDF/PPMI | whitening
métrica: dot product | cosine | L2
corpus/idioma/fecha/dominio + versión
```

### Matrices de coocurrencia

**[DLAI]** En *word-by-word*, `C[i,j]` cuenta cuántas veces `w_j` aparece dentro de una ventana `k` alrededor de `w_i`. La fila representa la palabra. En *word-by-document*, cada dimensión cuenta apariciones en documentos o categorías; pueden compararse también las columnas para representar colecciones.

```text
C_word[i,j] = Σ ocurrencias de contexto(w_j dentro de ±k de w_i)
X_doc[d,j]  = count(w_j, documento_d)
```

**[IMPL]** `k`, dirección izquierda/derecha, distancia ponderada, límites de oración, subsampling, min-count y tratamiento de documentos son hiperparámetros semánticos. Versionarlos. Conteos crudos favorecen palabras frecuentes y documentos largos; comparar al menos:

- binario/count y normalización por longitud;
- TF-IDF para documentos;
- PMI positivo (`PPMI`) o log-count para coocurrencias;
- n-gramas de palabras/caracteres;
- reducción por SVD/embeddings aprendidos cuando la matriz sparse sea demasiado grande.

Construir estadísticas sólo con train cuando la representación interviene en una evaluación supervisada. Para un índice no supervisado, documentar corpus y fecha de corte; no permitir que documentos futuros o privados contaminen una prueba histórica.

### Distancias: elegir la geometría que coincide con el índice

```text
L2(x,y) = sqrt(Σ_i (x_i-y_i)^2)
cos(x,y) = (x·y) / (||x||₂ ||y||₂)
```

**[DLAI]** L2 mide la línea recta; en conteos puede quedar dominada por el tamaño del corpus/documento. Coseno compara dirección y reduce ese efecto. En vectores de conteos no negativos, el ejemplo queda entre `0` y `1`.

**[IMPL] Correcciones y decisiones:**

- Coseno en vectores generales vive en `[-1,1]`, no siempre `[0,1]`.
- Coseno es indefinido para norma cero: filtrar o definir fallback explícito.
- En vectores L2-normalizados, ordenar por mayor coseno equivale a ordenar por menor L2 al cuadrado (`||x-y||²=2-2cos`). Esto permite alinear métrica e índice.
- Coseno elimina magnitud; si la magnitud contiene confianza, popularidad o cantidad de evidencia, puede perder señal.
- En alta dimensión aparecen concentración/hubness/an-isotropy; medir distribución de scores, vecinos dominantes y calidad por slice.
- No llamar “similares” a dos elementos sin especificar representación, métrica y umbral.

Tests: simetría, identidad, invariancia de escala positiva para coseno, comportamiento ante cero/NaN, equivalencia coseno–L2 tras normalizar y resultados manuales en 2–3 dimensiones. Usar tolerancias numéricas y la misma precisión/normalización offline y online.

### Aritmética, analogías y nearest neighbors

**[DLAI]** Una relación puede aproximarse por un desplazamiento: `capital_A - país_A + país_B`; luego se busca el vector más cercano para recuperar `capital_B`. Vecinos de una palabra pueden revelar términos relacionados.

**[IMPL]** Las analogías son una sonda del espacio, no una garantía algebraica ni una prueba de comprensión. Son sensibles a frecuencia, polisemia, corpus, tokenización, sesgo y métrica. Evaluar con múltiples relaciones y negativos difíciles; excluir los términos de la consulta al buscar vecinos. Reportar `top-k`, MRR/Recall@k cuando corresponda, no sólo un ejemplo exitoso.

Para producción:

1. congelar encoder/normalización;
2. generar vectores y metadatos con lineage;
3. usar búsqueda exacta como *oracle* en una muestra;
4. comparar el índice ANN contra el oracle (`Recall@k`) y contra relevancia humana/de negocio;
5. medir p50/p95/p99, QPS, memoria, build time y actualización;
6. versionar juntos embeddings, distancia e índice.

Un ANN puede recuperar con fidelidad los vecinos del embedding y aun así devolver resultados irrelevantes: **calidad del índice ≠ calidad semántica**.

### PCA: compresión lineal y visualización, no veredicto semántico

Para una matriz `X` con filas como ejemplos:

```text
μ = mean(X, axis=0)
Xc = X - μ
Xc = U S Vᵀ                    # SVD estable
Z_k = Xc V_k
explained_variance_ratio_k = Σ_{i≤k} S_i² / Σ_i S_i²
```

**[DLAI]** Los ejes principales son direcciones no correlacionadas ordenadas por varianza; proyectar en las primeras componentes retiene la mayor varianza posible entre subespacios lineales de esa dimensión. Dos o tres componentes permiten graficar relaciones del embedding.

**[IMPL] Correcciones:**

- Centrar es obligatorio para PCA estándar; escalar cada feature es una decisión distinta y puede cambiar por completo el resultado.
- PCA maximiza **varianza**, no “información” o semántica en sentido general.
- Una proyección 2D destruye distancias y vecindarios; un cluster aparente no valida el embedding.
- Ajustar PCA sólo en train/referencia y transformar dev/test con la misma media/componentes.
- Signos de componentes pueden invertirse sin cambiar la solución; tests y plots no deben depender del signo.
- Para matrices muy grandes usar SVD truncada/incremental, conservando un oracle pequeño.

Validar `explained_variance_ratio`, error de reconstrucción y estabilidad de vecinos/relaciones relevantes antes y después de reducir. Para una visualización, mostrar proporción explicada, número de puntos, criterio de selección y advertir que es una vista proyectada.

### Error analysis de representaciones

Muestrear consultas con vecinos correctos e incorrectos y clasificar: token/OOV; polisemia; frecuencia; dominio/idioma; documento largo; duplicados; sesgo social; hubness; mala métrica/normalización; embedding desactualizado; pérdida por reducción; relevancia que exige metadatos o contexto.

Para cada fallo, comparar:

```text
texto crudo → tokens → vector/norma → top-k exacto → top-k ANN
            → filtro/reranker → decisión final
```

Así se localiza si el problema nace en datos/encoder, geometría, aproximación del índice, filtros o ranking. Cambiar sólo la etapa responsable y convertir el caso en benchmark de regresión.

### Contrato para Codex

```text
Entrada: corpus/consulta + versión de preprocess/encoder
Artefactos: vocabulario o encoder, normalización, PCA opcional, índice, metadata schema
Invariantes: dimensión/métrica compatibles; vectores finitos; norma definida; lineage completo
Calidad: baseline léxico + exact-k oracle + relevancia + slices
Operación: latencia percentiles, throughput, memoria, frescura, costo de rebuild/update
Salida: ids, scores comparables sólo dentro de versión, metadata y trazas de etapa
Fallback: búsqueda léxica/híbrida o abstención si embedding/index no cumple contrato
```

**Cobertura de la semana:** 10 videos/transcripciones contrastados, 9 lecturas y 3 ejemplos de código inventariados; práctica/quiz/asignación no reproducidos.

## 4. Semana 4 — Machine Translation and Document Search

### Un mismo patrón: transformar, recuperar, verificar

**[DLAI]** La semana conecta dos tareas mediante nearest neighbors. Para traducción de palabras, aprende una matriz que lleva embeddings ingleses al espacio francés y busca el vecino destino más cercano. Para documentos, suma embeddings de palabras y recupera textos similares. Locality-sensitive hashing (LSH) reduce el espacio de candidatos para acelerar k-NN.

**[IMPL]** El patrón general reutilizable es:

```text
objeto → representación versionada → transformación opcional
       → generación aproximada de candidatos → score exacto/reranking
       → filtros/política → respuesta + trazabilidad
```

Separar evaluación de cada etapa. Una mala respuesta puede provenir del encoder/alineación, candidate recall, métrica, reranker, filtros o corpus; “el vector search falló” no es un diagnóstico.

### Alineación lineal entre espacios

Con pares de traducción alineados por fila, `X∈R^{m×d_src}`, `Y∈R^{m×d_tgt}`:

```text
R* = argmin_R (1/m) ||XR - Y||²_F
∇_R = (2/m) Xᵀ(XR-Y)
R ← R - η∇_R

consulta destino = x_src R
traducción = argmax_{y∈V_tgt} cosine(x_src R, y)
```

**[DLAI]** El curso entrena `R` con gradient descent sobre el Frobenius norm cuadrado y luego usa k-NN para palabras no incluidas en el diccionario de entrenamiento.

**[IMPL]** Antes de optimizar, verificar dimensiones, orden de multiplicación y correspondencia exacta de filas. Separar lexicón train/dev/test por pares y, cuando se mida generalización léxica, por lema/familia para evitar variantes casi duplicadas.

Comparar el entrenamiento iterativo con una solución least-squares/Procrustes. Si se exige una rotación ortogonal, resolver `R=UVᵀ` desde la SVD de `XᵀY`; una matriz libre puede escalar/deformar además de rotar. Centrado, normalización, anisotropía y espacios creados con objetivos/corpus incompatibles afectan la alineación.

Este es un traductor **word-level**, no un sistema general de traducción: no resuelve polisemia, morfología, orden, contexto, frases, entidades ni generación. Evaluar precision@1/top-k por frecuencia, categoría e idioma, con diccionario congelado y revisiones humanas de ambigüedad.

### k-NN exacto como oracle

**[DLAI]** Tras transformar un vector, rara vez coincide exactamente con un vector destino; se necesita buscar vecinos. Recorrer todo el vocabulario es costoso, por lo que el curso introduce buckets y LSH.

**[IMPL]** Mantener un buscador exacto vectorizado como referencia:

```text
build: normalizar base si la métrica es cosine
query: normalizar q; scores = Base @ q
result: top-k por score, excluyendo ids prohibidos/duplicados
```

Complejidad aproximada por query: `O(Nd)` y memoria `O(Nd)`. El oracle sirve para datasets pequeños, tests y medición del ANN; no confundir vecino geométrico exacto con relevancia correcta.

### Hash tables y random-hyperplane LSH

**[DLAI]** Un hash común asigna determinísticamente a buckets pero no preserva cercanía. Random-hyperplane LSH usa el signo de `v·p_i` para cada plano normal `p_i` y empaqueta bits en un bucket:

```text
b_i(v) = 1[v·p_i ≥ 0]
h(v) = Σ_i 2^i b_i(v)
```

Vectores de dirección similar tienen mayor probabilidad de compartir firma. Múltiples conjuntos independientes de planos/tablas elevan la probabilidad de generar buenos candidatos; se unen buckets y se rerankea con la métrica exacta.

**[IMPL] Correcciones:**

- `p·v` conserva el signo del lado del hiperplano, pero sólo es la longitud de la proyección escalar si `p` es unitario (en general se divide por `||p||`).
- Más planos producen buckets más finos: menos candidatos y colisiones útiles; más tablas/probes aumentan cobertura y costo. “Más regiones = más accuracy y más lentitud” no es una ley aislada: depende de ambos parámetros y la distribución.
- ANN sacrifica normalmente **recall/fidelidad de vecinos** por eficiencia; no asumir que baja la métrica estadística llamada precision.
- Definir conducta para `v·p=0`, semillas, duplicados y candidatos insuficientes.

Parámetros mínimos a versionar: dimensión, normalización/métrica, semilla/planos, bits por tabla, número de tablas, probes/buckets vecinos, `k`, límite de candidatos y política de updates.

### Evaluar ANN sin autoengañarse

```text
candidate_recall@k = |ANN_k(q) ∩ Exact_k(q)| / k
task_recall@k      = relevantes_recuperados / relevantes_totales
precision@k        = relevantes_en_top_k / k
MRR, nDCG@k        = calidad del orden cuando aplica
```

Medir además p50/p95/p99, QPS, memoria, tiempo de build, tamaño, frescura y costo de insert/delete/rebuild. Barrer parámetros y conservar la frontera Pareto calidad–latencia–memoria. El benchmark debe incluir distribución real de queries, cold/warm cache, concurrencia y filtros; separar tiempo de embedding, búsqueda, reranking y red.

Tests de índice:

- exact search coincide con cálculo manual pequeño;
- serializar/cargar preserva vecinos;
- misma semilla reproduce firmas;
- probes/tablas adicionales no pierden candidatos por bug;
- ningún id borrado/filtrado reaparece;
- actualización no mezcla dimensiones/versiones;
- ANN alcanza el recall mínimo contra oracle bajo cada slice.

### Representación y búsqueda documental

**[DLAI]** El ejemplo suma embeddings de las palabras conocidas del documento y aplica k-NN; el curso señala que la estructura “texto a vector y vecinos semánticos” reaparece en NLP moderno.

```text
d = Σ_{w∈doc∩V} e_w
```

**[IMPL]** La suma pierde orden, negación y contexto, ignora OOV y crece con longitud. Comparar:

- promedio/normalización por longitud;
- ponderación TF-IDF/SIF y eliminación de componente dominante;
- n-gramas y recuperación lexical BM25;
- encoder de frases/documentos;
- búsqueda híbrida lexical+dense y reranker.

No elegir por una demo de paráfrasis. Crear queries y juicios de relevancia representativos, con negativos léxicamente parecidos y semánticamente distintos, entidades exactas, negación, fechas, código/identificadores, idiomas y consultas sin respuesta. Evitar duplicados/paráfrasis entre index/train/test.

### Error analysis de retrieval

Para cada fallo registrar:

```text
query + intención/idioma/slice
documentos relevantes conocidos
top-k exacto del embedding
candidatos ANN y candidate_recall
top-k tras filtros/reranker
versión de corpus/encoder/index
latencias por etapa
```

Taxonomía: consulta ambigua; corpus sin respuesta; chunking/metadatos; lexical exact-match; embedding/alineación; OOV/multilingüe; ANN miss; filtro incorrecto; reranker; documento obsoleto; duplicado; política de abstención. Priorizar `frecuencia × severidad`, corregir la etapa causal y agregar el caso al benchmark.

Para RAG, retrieval recall es una condición necesaria pero no suficiente: evaluar por separado contexto recuperado, uso fiel del contexto, calidad de respuesta, citas y abstención. Nunca “arreglar” ausencia de evidencia aumentando creatividad del generador.

### Contrato de implementación de traducción y búsqueda

```text
Artefactos: encoder/alignment R, preprocess, corpus/chunks, vectors, ANN, metadata, reranker
Compatibilidad: dimensión + métrica + normalización + versiones deben coincidir
Oracles: búsqueda exacta y conjunto de relevancia humano
Gates: Recall@k, task metrics, slices, p99/QPS/memoria/frescura
Observabilidad: latencia/errores por etapa, cobertura, zero-result, drift de score/norma
Actualización: estrategia idempotente, tombstones/deletes, rebuild y rollback probado
Seguridad: autorización antes y después de recuperar; no filtrar datos entre tenants
Salida: ids, scores, versiones, filtros y evidencia; abstención cuando no hay soporte
```

**Cobertura de la semana:** 10 videos/transcripciones contrastados, 7 lecturas y 2 ejemplos de código inventariados; práctica/quiz/asignación no reproducidos.

### Cierre del Curso 1

El curso construye una progresión coherente: baseline discriminativo → baseline generativo/probabilístico → geometría de representaciones → recuperación aproximada. La transferencia profesional no consiste en usar siempre esos modelos, sino en conservar sus oracles, diagnósticos y contratos al escalar hacia embeddings contextuales, índices modernos y RAG.

**Gate alcanzado:** 47 videos contrastados de punta a punta; 39 lecturas y 9 ejemplos de código inventariados; 5 evaluaciones certificables contabilizadas pero no abiertas ni reproducidas. El contenido operativo fue corregido donde una simplificación pedagógica no debía convertirse en regla de producción.

# Parte II — Probabilistic Models

## 5. Curso 2 · Semana 1 — Autocorrect and Minimum Edit Distance

### Alcance: error ortográfico no es error contextual

**[DLAI]** El baseline detecta una palabra que no pertenece al vocabulario, genera cadenas a una o más ediciones, filtra palabras conocidas y elige el candidato de mayor probabilidad unigram. No corrige una palabra válida pero inadecuada en contexto —por ejemplo, un homófono o término real distinto—; eso requiere contexto.

```text
input → ¿está en vocabulario? → candidatos por edits → ∩ vocabulario
      → prior de palabra → ranking → reemplazo/abstención
```

**[IMPL]** Definir producto antes del algoritmo: corrección automática, sugerencias top-k, subrayado o búsqueda tolerante a errores tienen costos distintos. En campos de código, símbolos, nombres, medicamentos, direcciones o mercados, una corrección silenciosa puede ser peor que no corregir. Conservar siempre el texto original y ofrecer abstención.

### Generación y ranking de candidatos

**[DLAI]** Las operaciones enseñadas para generar candidatos son insertar, borrar, reemplazar e intercambiar caracteres adyacentes. El corpus proporciona:

```text
P(w) = count(w) / Σ_v count(v)
ĉ = argmax_{c ∈ candidates(x)∩V} P(c)
```

**[IMPL]** Esto es un prior de popularidad, no una probabilidad contextual ni una probabilidad de que `c` haya originado el typo. Un ranking más completo separa:

```text
score(c|x,context) = log P(c|context) + log P(x|c) + features/policy
```

`P(x|c)` puede ponderar confusiones de teclado, fonética, idioma y hábitos; `P(c|context)` evita elegir siempre la palabra globalmente frecuente. Mantener cada componente auditable.

La generación exhaustiva explota con longitud, alfabeto y radio. Aproximadamente, un solo paso produce `n` borrados, `n-1` transposiciones, `n·|Σ|` reemplazos y `(n+1)|Σ|` inserciones antes de deduplicar. Evitar materializar edits de radio 2–3 sin límites. Opciones:

- generar radio 1, filtrar vocabulario y ampliar sólo si hace falta;
- deduplicar con set e impedir reemplazar por el mismo carácter;
- trie/DAWG o autómata de Levenshtein para buscar vocabulario sin enumerar todo;
- SymSpell/BK-tree según distribución, distancia y updates;
- candidate cap, deadline y fallback por request.

Tokenización/normalización forman parte del modelo: Unicode NFC/NFKC, mayúsculas, acentos, apóstrofes, guiones, emoji, scripts y grapheme clusters. No operar ingenuamente sobre bytes ni asumir que un carácter visible equivale a un code point.

### Minimum edit distance por programación dinámica

Para source `s[1..m]`, target `t[1..n]` y costos configurables:

```text
D[0,0] = 0
D[i,0] = D[i-1,0] + delete_cost(s_i)
D[0,j] = D[0,j-1] + insert_cost(t_j)

D[i,j] = min(
  D[i-1,j]   + delete_cost(s_i),
  D[i,j-1]   + insert_cost(t_j),
  D[i-1,j-1] + (0 if s_i=t_j else replace_cost(s_i,t_j))
)
```

**[DLAI]** La tabla reutiliza soluciones de prefijos; la esquina inferior derecha contiene el costo mínimo. Un backpointer por celda reconstruye la secuencia de operaciones. La configuración pedagógica usa insertar/borrar `1` y reemplazar `2`.

**[IMPL] Corrección terminológica:** Levenshtein estándar usa normalmente costo unitario para inserción, borrado y sustitución, sin transposición elemental. La configuración `1/1/2` es una **edit distance ponderada** compatible con interpretar sustitución como delete+insert. Si se permite intercambio adyacente, especificar Damerau–Levenshtein/Optimal String Alignment y sus diferencias. Nunca llamar a variantes distintas con el mismo nombre en API/tests.

Complejidad del DP completo: tiempo `O(mn)`, memoria `O(mn)`. Si sólo se necesita la distancia, conservar dos filas da `O(min(m,n))` memoria. Si se requiere el camino/alineación, guardar backpointers o usar reconstrucción divide-and-conquer. Aplicar banding/Ukkonen cuando sólo interesan distancias `≤k`; cortar por deadline/longitud para inputs adversariales.

### Pruebas del algoritmo

- vacío↔vacío, vacío↔texto y simetría cuando los costos sean simétricos;
- identidad `D(x,x)=0` y resultados manuales de inserción/borrado/sustitución;
- caso donde sustitución `2` empata delete+insert;
- transposición demuestra la diferencia entre Levenshtein y Damerau;
- backtrace reproduce exactamente el target y su costo suma `D[m,n]`;
- implementación de dos filas coincide con matriz completa;
- Unicode/graphemes, texto muy largo y límite `k` no rompen memoria/latencia;
- fuzz/property tests comparan una implementación optimizada contra oracle simple.

Si los costos dependen de dirección (`P(x|c)`), la distancia puede dejar de ser simétrica; documentarlo y no usar estructuras que requieran una métrica matemática sin comprobar sus invariantes.

### Evaluación y error analysis

Crear un dataset de triples `(contexto, texto observado, intención/corrección)` a partir de logs consentidos o errores sintéticos calibrados. Separar por usuario/documento/tiempo para no memorizar hábitos. Reportar:

- detection precision/recall: ¿se decidió intervenir correctamente?;
- correction accuracy@1 y recall@k/MRR;
- keystroke savings o aceptación/rechazo;
- tasa de corrección dañina y abstención;
- latencia p50/p95/p99, candidatos generados y memoria;
- slices por idioma, longitud, distancia, palabra rara, entidad y dispositivo.

Inspeccionar fallos en etapas: falso OOV por vocabulario/normalización; candidato correcto no generado; podado por vocabulario; prior/contexto lo rankea mal; nombre propio/dominio; confusión de teclado no modelada; error real-word; idioma equivocado; overcorrection. Medir `frecuencia × costo`, arreglar la etapa responsable y convertir cada categoría recurrente en test.

### Contrato operativo de autocorrect

```text
Entradas: texto original + locale/contexto + versiones de vocabulario/modelo
Salida: original, candidates[{term,score,edit_path}], acción y razón/abstención
Artefactos: normalizer, vocab/frequencies, confusion costs, LM/ranker, thresholds
Límites: max length/radius/candidates/time; sin corrección silenciosa fuera de política
Monitoreo: OOV, aceptación, daño, no-result, drift léxico, latencia y cache hit
Privacidad: minimizar/retener con política; no aprender secretos ni mezclar tenants
Rollback: diccionario/modelo versionado y despliegue gradual
```

**Cobertura de la semana:** 11 videos/transcripciones contrastados, 9 lecturas y 2 ejemplos de código inventariados; quiz/asignación no abiertos ni reproducidos.

## 6. Curso 2 · Semana 2 — POS Tagging, HMM and Viterbi

### Qué modela y qué no

**[DLAI]** POS tagging asigna una categoría gramatical a cada token. La misma palabra puede ser nombre, verbo, adverbio, etc. según contexto. El curso usa estados POS, probabilidades de transición entre tags y probabilidades de emisión de palabras; Viterbi recupera la secuencia de estados más probable.

**[IMPL]** POS puede aportar features a parsing, búsqueda, NER, coreference o speech, pero no resuelve esas tareas por sí solo. El tagset, tokenización y guía de anotación definen la verdad operacional; mapear Penn Treebank/Universal Dependencies o cambiar granularidad requiere una evaluación nueva.

### HMM supervisado de primer orden

Para tokens observados `w_1:T` y tags ocultos durante inferencia `t_1:T`:

```text
P(t_1:T,w_1:T) = P(t_1|START) · Π_i P(w_i|t_i)
                 · Π_{i=2..T} P(t_i|t_{i-1}) · P(END|t_T)

A[a,b] = P(t_i=b | t_{i-1}=a)
B[b,w] = P(w_i=w | t_i=b)
```

**[DLAI]** Se estiman `A` contando pares de tags y `B` contando pares tag–palabra; se agregan tokens de inicio y smoothing para evitar ceros. Cada fila es una distribución normalizada.

**[IMPL]** Los supuestos fuertes deben quedar registrados:

- **Markov:** el tag actual depende sólo del tag anterior;
- **emisión:** la palabra depende sólo de su tag actual;
- **estacionariedad:** las mismas matrices gobiernan toda posición/dominio;
- train anotado representa el uso futuro.

Los tags son observables en el corpus supervisado de entrenamiento y ocultos para la oración nueva. La emisión es `P(palabra|tag)`, no “probabilidad de ir de palabra a tag”; `P(tag|palabra)` exigiría Bayes/normalización diferente.

Con add-`α`:

```text
A[a,b] = (count(a,b)+α_A) / (count(a,*)+α_A|TagsNext|)
B[t,w] = (count(t,w)+α_B) / (count(t,*)+α_B|V|)
```

Elegir `α` en dev y definir si START/END se suavizan. No prohibir o permitir inicios/puntuación por accidente. Verificar que cada fila sume uno y que la convención de orientación de matrices sea única en código/documentación.

### OOV y rare words

Suavizar sobre un vocabulario fijo no da una política útil para cualquier string desconocido. Antes de estimar conteos, reemplazar palabras raras por clases `UNK` reproducibles, por ejemplo:

```text
UNK, UNK_NUM, UNK_CAP, UNK_ALLCAPS, UNK_SUFFIX_ING,
UNK_SUFFIX_ED, UNK_PUNCT, UNK_URL, UNK_EMOJI
```

La misma función debe operar online. Validar categorías mutuamente excluyentes/priorizadas y conservar cobertura. Alternativas: features morfológicas, char/subword models o tagger contextual si OOV/idioma mixto domina el error.

### Viterbi en log-space

Definir `δ_i(s)` como el mejor log-score de una ruta que termina en estado `s` en posición `i`, y `ψ_i(s)` como su predecesor:

```text
δ_1(s) = log A[START,s] + log B[s,w_1]

δ_i(s) = log B[s,w_i] + max_p (δ_{i-1}(p) + log A[p,s])
ψ_i(s) = argmax_p (δ_{i-1}(p) + log A[p,s])

t_T* = argmax_s (δ_T(s) + log A[s,END])
t_{i-1}* = ψ_i(t_i*)
```

**[DLAI]** Inicialización, forward y backward/backtrace llenan una matriz de scores y otra de predecesores. Se recomiendan logs para evitar underflow.

**[IMPL] Correcciones:** Viterbi no elige greedily el tag local más probable; conserva el mejor prefijo para **cada estado** y optimiza la secuencia conjunta. Incluir transición END si el modelo la entrenó. Fijar tie-breaking determinista y usar `-∞` para rutas imposibles; nunca `log(0)` por descuido.

Para `T` tokens y `K` tags: tiempo `O(TK²)`, memoria `O(TK)` para backtrace. Vectorizar el máximo sobre estados previos, usar matrices sparse/beam sólo si el benchmark lo exige y comparar siempre contra implementación oracle pequeña.

### Invariantes y pruebas

- `A`/`B` no contienen `NaN`; filas normalizan dentro de tolerancia.
- Conteos de START/END coinciden con número de oraciones válidas.
- Un corpus diminuto produce matrices y mejor ruta calculables a mano.
- Viterbi coincide con enumeración exhaustiva para secuencias/estados pequeños.
- Multiplicar probabilidades y sumar logs eligen la misma ruta en casos no extremos.
- Backtrace tiene exactamente `T` tags y reproduce el score terminal.
- Permutar orden de tags sólo permuta índices, no predicciones semánticas.
- OOV categories son idénticas en train/serving.
- Oraciones vacías, un token, muy largas y puntuación tienen comportamiento definido.

### Evaluación y aprendizaje de errores

Baselines mínimos: tag global mayoritario; tag más frecuente por palabra con fallback; regla/morfología simple. Reportar token accuracy, macro/per-tag precision–recall–F1, sentence exact match, confusion matrix y rendimiento por OOV, palabra ambigua, longitud, dominio, mayúsculas y puntuación. Accuracy agregada puede ocultar tags raros críticos.

Para cada error guardar token/oración, gold/predicción, `UNK`, emisiones candidatas, transición dominante y camino alrededor. Clasificar:

- anotación/tagset ambiguo;
- tokenización/normalización;
- OOV/morfología;
- emisión equivocada por frecuencia;
- dependencia mayor que orden 1;
- contexto largo/semántica;
- dominio/idioma drift;
- bug de START/END, smoothing, índices o backtrace.

Corregir datos/pipeline antes de aumentar modelo. Si la mayoría de fallos exige contexto bilateral/largo, el límite es estructural y justifica CRF/biLSTM/Transformer; conservar HMM como baseline y oracle de pipeline.

### Contrato operativo de POS/HMM/Viterbi

```text
Entrada: tokens + tokenizer/tagset/model versions
Artefactos: tag index, vocab/UNK mapper, A/B log-probs, START/END policy
Salida: tags, sequence log-score, OOV mask, model version
Tests: normalización, oracle exhaustivo, log numerics, backtrace, parity batch/online
Métricas: token/per-tag/sentence + OOV/ambiguity/domain + p99/throughput
Fallback: tagger léxico/reglas o abstención para schema/idioma no soportado
```

**Cobertura de la semana:** 13 videos/transcripciones contrastados, 11 lecturas y 2 ejemplos de código inventariados; práctica/quiz/asignación no abiertos ni reproducidos.

## 7. Curso 2 · Semana 3 — Autocomplete and N-gram Language Models

### Modelo y aproximación

**[DLAI]** Un language model asigna probabilidades a secuencias y predice el próximo token. La regla de la cadena es exacta, pero un n-gram aproxima el historial completo usando sólo los `n-1` tokens anteriores.

```text
P(w_1:T) = Π_i P(w_i | w_1:i-1)                         # chain rule

P(w_i | w_1:i-1) ≈ P(w_i | w_{i-n+1:i-1})              # Markov n-gram

P_MLE(w|h) = count(h,w) / count(h)
```

**[IMPL]** `n` controla contexto frente a sparsity/memoria. Ajustarlo en dev y comparar unigram/bigram/trigram bajo idéntico vocabulario, splits y normalización. Conteos grandes deben almacenarse sparse (`history → {next: count}`), no como una matriz densa `|V|^n`.

### Fronteras de oración y corpus

**[DLAI]** Añadir `n-1` tokens START y un token END permite estimar los primeros tokens, modelar longitud/terminación y conservar denominadores coherentes.

**[IMPL]** No concatenar documentos como si la última palabra de uno precediera a la primera del siguiente. Definir si párrafo, turno, query o línea es una secuencia. START no cuenta como predicción; END normalmente sí al calcular likelihood/perplexity. Tokenización, casing, puntuación y límites son parte del artefacto.

Dividir por documento/usuario/sesión/tiempo antes de crear ventanas. Un split aleatorio de n-gramas filtra frases vecinas/duplicadas y produce perplexity irreal. Deduplicar de forma consciente, conservar dominio y versionar fecha/corpus.

### OOV y smoothing son problemas diferentes

**[DLAI]** `UNK` representa palabras fuera del vocabulario; smoothing asigna masa a n-gramas no observados aunque sus palabras sean conocidas. Comparar perplexity sólo con el mismo vocabulario porque muchos `UNK` pueden reducirla artificialmente.

**[IMPL]** Construir vocabulario **sólo con train**, reemplazar rare words en train y aplicar el mismo mapping a dev/test/serving. Reportar OOV rate y `UNK` rate junto a perplexity. Un modelo que predice `UNK` con facilidad no es más útil.

Add-`k`:

```text
P_k(w|h) = (count(h,w)+k) / (count(h)+k|V_next|)
```

Add-one es pedagógico y suele sobreasignar masa a eventos no vistos en vocabularios grandes. Ajustar `k` en dev y comparar con interpolación/backoff; para n-gramas serios, considerar Katz/Good–Turing o interpolated Kneser–Ney. *Stupid backoff* es un score práctico de ranking, no una distribución normalizada; no usarlo donde se requieran probabilidades calibradas/perplexity válida.

En interpolación:

```text
P(w|h_2) = λ_3 P_MLE(w|h_2) + λ_2 P(w|h_1) + λ_1 P(w)
λ_i ≥ 0,  Σ_i λ_i = 1
```

Ajustar lambdas sin tocar test. En backoff, descontar masa de n-gramas vistos antes de reasignarla a órdenes menores cuando se necesita una distribución válida.

### Perplexity correctamente definida

Para `M` tokens evaluados y probabilidades en log natural:

```text
NLL/token = -(1/M) Σ_i log P(w_i|h_i)
PP = exp(NLL/token)

# con log base 2:
cross_entropy_bits = -(1/M) Σ_i log₂ P(w_i|h_i)
PP = 2^(cross_entropy_bits)
```

**[DLAI]** Menor perplexity significa mayor probabilidad media asignada al conjunto bajo ese modelo.

**[IMPL] Correcciones:**

- La media negativa de logs es cross-entropy/log-perplexity; el signo negativo no puede omitirse.
- No existe rango universal “bueno” como `20–60`: PP depende de idioma, corpus, tokenización, vocabulario, OOV, longitud y unidad (word/char/subword).
- Comparar sólo sobre exactamente el mismo test y pipeline.
- Menor PP no garantiza texto más humano, mejor autocomplete, veracidad, seguridad ni métrica downstream.
- Calcular en log-space; una probabilidad cero vuelve PP infinita y revela smoothing/OOV roto.

Unit tests: secuencia diminuta calculada a mano; filas normalizan; START/END correctos; perplexity de corpus repetible; modelo que asigna mayor likelihood obtiene menor PP; cambiar base de log conserva PP tras conversión.

### Autocomplete y generación

**[DLAI]** Para sugerir, tomar el último historial `n-1`, obtener la distribución de próximos tokens y devolver el más probable o muestrear para generar hasta END.

**[IMPL]** Producto real:

```text
history → normalize/tokenize → lookup/backoff → candidate scores
        → filtros de política/contexto → top-k → UI/aceptación
```

No optimizar únicamente PP. Medir top-k accuracy/recall, MRR, keystroke savings, aceptación, tiempo hasta completar, false-trigger/annoyance, diversidad cuando sea necesaria, seguridad/PII y p95/p99. El candidato correcto puede no estar en top-1 pero ahorrar más escritura.

Para sampling, explicitar greedy/temperature/top-k y semilla; imponer longitud máxima y detectar loops. Un n-gram reproduce fragmentos frecuentes: auditar memorization/PII/profanidad y aplicar filtros antes de mostrar, no sólo al entrenar.

### Eficiencia

- Precomputar top candidatos por history frecuente; cache acotada por versión/locale.
- Usar IDs enteros, counts compactos y lookup sin crear substrings/copias innecesarias.
- Backoff temprano para histories ausentes; deadline y fallback unigram.
- Medir tamaño de tablas, cold start, cache hit, QPS y actualizaciones.
- Mantener construcción offline separada del hot path; publicar snapshot inmutable/atómico y rollback.

Para contextos personales, aislar usuario/tenant, minimizar retención y combinar modelo global con adaptación local controlada. No registrar texto sensible por defecto.

### Error analysis de autocomplete

Por error guardar history original/tokenizado, nivel n-gram usado, OOV, candidatos/scores, elección/aceptación y latencia. Etiquetar: split/leakage; tokenización; corpus ausente/desactualizado; OOV; history sparse; smoothing/backoff; entidad/idioma; intención contextual larga; unsafe/PII; ranking/UI.

Separar:

1. **candidate miss**: la continuación deseada no fue generada;
2. **ranking miss**: estaba pero quedó abajo;
3. **policy miss**: fue filtrada o debió filtrarse;
4. **interaction miss**: buena sugerencia, mal momento/formato.

Priorizar por volumen y daño. Si la señal requiere contexto más largo/semántica, escalar a neural LM/Transformer sólo después de demostrar el límite y conservar n-gram como fallback/benchmark.

### Contrato operativo de autocomplete

```text
Artefactos: tokenizer, vocab/UNK, n-gram counts, smoothing/backoff, top-k cache
Compatibilidad: n + vocab + boundaries + normalization + corpus version
Calidad: PP comparable + top-k/MRR/acceptance + OOV + slices + seguridad
Operación: p99, QPS, memoria, cache, frescura, atomic publish/rollback
Salida: candidates con score/nivel de backoff/version; nunca fingir calibración
Fallback: orden inferior/unigram o ninguna sugerencia
```

**Cobertura de la semana:** 11 videos/transcripciones contrastados, 9 lecturas y 3 ejemplos de código inventariados; práctica/quiz/asignación no abiertos ni reproducidos.

## 8. Curso 2 · Semana 4 — Word Embeddings with Neural Networks

### La representación se aprende de un objetivo

**[DLAI]** One-hot preserva identidad sin imponer orden, pero es sparse, grande y no expresa semejanza. Un embedding denso se aprende de un corpus mediante una tarea autosupervisada; CBOW predice la palabra central desde palabras circundantes. Corpus, contexto, objetivo y dimensión determinan lo que el vector captura.

**[IMPL]** No describir un embedding como “el significado” de una palabra. Es una solución no única a un objetivo estadístico, condicionada por:

```text
corpus/dominio/fecha/idioma
tokenizer + vocabulario + rare/OOV policy
ventana/dirección/distancia
objetivo/negativos/loss/optimizer
dimensión/regularización/seed
```

Dos ejecuciones pueden preservar vecindarios pero rotar/reflejar ejes; no testear coordenadas absolutas. Vectores estáticos dan una representación por token y mezclan polisemia; representaciones contextuales cambian con la oración. *Polisemia* significa varios sentidos de una palabra, no “palabras con significados similares”.

### Elegir método antes de entrenar desde cero

**[DLAI]** Se presentan Word2Vec (CBOW/skip-gram), GloVe, FastText y modelos contextuales como ELMo/BERT/GPT. FastText usa char n-grams y puede formar vectores para muchas palabras no vistas.

**[IMPL] Matriz de decisión:**

| Necesidad | Primera opción |
|---|---|
| baseline pequeño/interpretable | TF-IDF/PPMI/SVD o CBOW/skip-gram |
| morfología/OOV | FastText/subword |
| polisemia y contexto | encoder contextual |
| búsqueda de frases/documentos | encoder entrenado/evaluado para retrieval |
| dominio con pocos datos | pretrained + evaluación/adaptación |
| privacidad/edge/latencia extrema | vocab/modelo compacto propio si el lift lo justifica |

Promediar palabras para frases no es una propiedad exclusiva de FastText y pierde orden/contexto. Fine-tuning de un Transformer tampoco garantiza “embeddings de alta calidad”: elegir pooling, objetivo contrastivo y benchmark de la tarea.

### Datos CBOW

Con half-window `c`, para cada posición válida:

```text
context_i = [w_{i-c:i}, w_{i+1:i+c+1}]
target_i  = w_i
x_i       = average(one_hot(context_i))
y_i       = one_hot(target_i)
```

**[DLAI]** Una ventana deslizante produce ejemplos; un generator con `yield` evita materializar todo. El promedio es bag-of-words: cuenta repeticiones pero pierde orden.

**[IMPL]** No cruzar límites de oración/documento salvo decisión explícita. `c` cambia el tipo de relación: ventanas cortas suelen enfatizar sintaxis; más amplias, tópico/semántica. Probar distancia ponderada/dirección si la tarea lo requiere. Subsamplear tokens extremadamente frecuentes y controlar rare words; medir distribución de targets.

La recomendación pedagógica de lowercase/eliminar símbolos no es universal. Puntuación, casing, números, fórmulas, moneda, hashtags, emojis y markup pueden ser señal o datos críticos. Decidir por ablation y dominio; versionar Unicode/tokenizer. En código, finanzas o biomedicina, borrar símbolos puede destruir el significado.

### Arquitectura y shapes

Para vocabulario `V`, dimensión `D` y batch de columnas `M`:

```text
X:  [V,M]          Y:  [V,M]
W1: [D,V]          b1: [D,1]
W2: [V,D]          b2: [V,1]

Z1 = W1X + b1
H  = ReLU(Z1)                         # [D,M]
Z2 = W2H + b2                         # logits [V,M]
Ŷ  = softmax(Z2, axis=vocab)          # [V,M]
J  = -(1/M) Σ_examples Σ_vocab Y·log Ŷ
```

**[IMPL]** Softmax regression/multinomial logistic regression sí soporta más de dos clases; la frontera no es “logistic regression sólo tiene dos salidas”. Implementar cross-entropy desde logits con `logsumexp` o primitiva fusionada:

```text
log_softmax(z)_i = z_i - logsumexp(z)
```

No calcular `exp(z)` sin restar el máximo. Comprobar eje de softmax, broadcasting de bias y promedio de loss. Para una implementación moderna usar IDs+embedding lookup; no construir one-hot denso `V×M`.

El full softmax cuesta aproximadamente `O(V·D)` por ejemplo y matrices `O(V·D)`. Para vocabularios grandes comparar negative sampling, sampled/adaptive softmax o hierarchical softmax; declarar que cambian objetivo/calibración. Batch, mixed precision y vectorización se eligen por profiling, no por intuición.

### Entrenamiento y verificación

**[DLAI]** Backpropagation calcula gradientes de `W1,W2,b1,b2`; gradient descent actualiza parámetros. Los embeddings pueden extraerse como columnas de `W1`, filas de `W2` o su promedio compatible.

**[IMPL]** `W1` (input/context) y `W2ᵀ` (output/target) aprenden roles distintos; promediar es una opción a validar, no una identidad teórica. Versionar cuál se publica.

Tests/gates:

- shapes y ejes en single/batch; probabilities suman uno;
- loss/logits finitos en valores extremos;
- gradient check en red diminuta;
- overfit de microcorpus y descenso de loss;
- label/token shuffle destruye vecinos esperados;
- batching/shuffle no cruza documentos ni pierde ejemplos;
- lookup con IDs coincide con one-hot oracle;
- checkpoint/reload conserva scores y vocab mapping;
- seeds, corpus hash y config permiten reproducir métricas, no coordenadas exactas.

Monitorear loss train/dev, norma/varianza, tokens por segundo, memoria, grad norm, frecuencia y vecinos de probes congelados. Dead ReLU o embeddings colapsados requieren revisar activación/LR/inicialización/datos.

### Evaluación intrínseca, extrínseca y sesgo

**[DLAI]** Intrínseca examina analogías, similitud, clustering o visualización; extrínseca usa el embedding dentro de la tarea final. Las analogías admiten múltiples respuestas y no son perfectas.

**[IMPL] Corrección:** la **evaluación extrínseca/downstream** es la prueba decisiva de utilidad para el producto, aunque mezcla calidad del embedding con el resto del sistema. La intrínseca es más rápida y diagnóstica, pero puede no correlacionar con el objetivo. Una proyección visual es exploración, no métrica.

Protocolo:

1. congelar modelo downstream/seed/presupuesto;
2. sustituir sólo el embedding;
3. evaluar pares/intervalos y slices OOV/frecuencia/dominio/idioma;
4. medir calidad, latencia, memoria y cobertura;
5. usar intrinsic probes para explicar regresiones, no para anular el resultado downstream.

Auditar asociaciones estereotipadas, nombres/dialectos, términos sensibles, corpus tóxico/privado y memorization. Vecindad refleja asociación del corpus, no verdad, causalidad ni recomendación. La herramienta de similitud debe incluir mitigación/evaluación de daño correspondiente a su uso.

### Error analysis de embeddings

Clasificar fallos: tokenización/OOV; rare/frequent imbalance; ventana inadecuada; corpus/domain drift; polisemia estática; objetivo no alineado; optimization/collapse; métrica/normalización; sesgo; evaluación contaminada; downstream bottleneck.

Para cada probe guardar vecinos exactos con scores, frecuencia/norma, ejemplos de contexto, versión y resultado downstream. Comparar `W1`, `W2ᵀ`, promedio y baselines bajo el mismo protocolo. Si el problema exige contexto, no “arreglar” el vector estático con más dimensiones sin evidencia.

### Contrato operativo de word embeddings

```text
Entrada: token/id + tokenizer/vocab/model versions
Artefactos: corpus lineage, vocab, config, W1/W2 o embedding publicado, normalization
Calidad: intrinsic probes + downstream A/B + bias/OOV/domain slices
Operación: lookup/encode p99, memory, throughput, size, compatibility
Actualización: retrain/adapt, alignment/migration, index rebuild y rollback
Salida: vector sin interpretación causal; score sólo con métrica/version compatibles
Fallback: pretrained/lexical/subword o abstención para token/domain no soportado
```

**Cobertura de la semana:** 22 videos/transcripciones contrastados, 20 lecturas y 5 ejemplos de código inventariados; práctica/quiz/asignación no abiertos ni reproducidos.

### Cierre del Curso 2

La progresión es coherente: edit distance/autocorrect → HMM/Viterbi → n-gram language models → embeddings autosupervisados. El hilo transversal es estimar estructura desde conteos o predicción, manejar lo no visto y usar dynamic programming/error analysis antes de escalar.

**Gate alcanzado:** 57 videos contrastados de punta a punta; 49 lecturas y 12 ejemplos de código inventariados; 6 evaluaciones certificables contabilizadas pero no abiertas ni reproducidas. Las reglas pedagógicas dependientes de corpus, vocabulario, smoothing o arquitectura quedaron convertidas en opciones verificables, no defaults universales.

# Parte III — Sequence Models

## 9. Curso 3 · Semana 1 — Deep Sentiment and Recurrent Language Models

### Dos baselines neuronales distintos

La semana contiene dos bloques que no deben confundirse:

1. **embedding + mean + dense** para sentimiento: aprende features densas, pero sigue siendo bag-of-embeddings y pierde orden;
2. **RNN/GRU** para secuencias/language modeling: mantiene un estado recurrente y comparte parámetros a través del tiempo.

**[DLAI]** Las redes aprenden representaciones adecuadas para la tarea y pueden superar conteos simples. **[IMPL]** “No requieren feature engineering” es una simplificación: tokenizer, corpus, padding, arquitectura, objetivo, contexto y labels siguen siendo diseño de representación.

### Sentimiento con mean pooling

```text
ids = tokenize(text)
E = embedding(ids)                         # [T,D]
v = masked_mean(E, valid_tokens)           # [D]
logits = MLP(v)
loss = cross_entropy(logits, label)
```

No promediar PAD: usar mask y denominador de tokens válidos; definir texto vacío. Mean pooling borra orden/negación y puede diluir tokens raros en textos largos. Comparar contra regresión logística/Naïve Bayes y contra pooling max/attention/RNN. Si `not good` y `good not` quedan iguales, es límite estructural, no falta de epochs.

Split por usuario/hilo/tiempo, padding dinámico por batch, packed/ragged sequences cuando compense y threshold/calibración igual que en Semana 1 del Curso 1.

### RNN vanilla

Para input `x_t`, estado `h_t` y logits `z_t`:

```text
h_t = φ(W_x x_t + W_h h_{t-1} + b_h)
z_t = W_y h_t + b_y
p_t = softmax(z_t)
```

**[DLAI]** Los mismos parámetros se reutilizan en cada posición, permitiendo longitudes variables y un resumen del prefijo. Se presentan many-to-one, one-to-many y many-to-many, además de encoder–decoder.

**[IMPL]** El estado recibe información de todo el prefijo en teoría, pero la capacidad efectiva de conservarla es limitada. Una RNN no garantiza dependencia larga ni siempre usa menos memoria/latencia que un n-gram; comparar footprint, throughput y calidad bajo el hardware real.

Para language modeling autoregresivo:

```text
L = - Σ_t mask_t · log p_t[target_{t+1}] / Σ_t mask_t
```

No promediar sobre padding. Shift input/target una posición, incluir política BOS/EOS y calcular perplexity con el mismo tokenizer. En entrenamiento usar secuencias truncadas/batches con estado reseteado o detach explícito; no filtrar estado entre documentos/usuarios.

### Backpropagation through time y estabilidad

La semana muestra forward/loss, pero el contrato profesional incluye BPTT. Productos repetidos de Jacobianos causan gradients que desaparecen o explotan. Instrumentar grad norm, activation/state norms y longitud efectiva.

- gradient clipping mitiga explosión, no vanishing;
- inicialización/normalización/residual/truncated BPTT cambian dinámica;
- clipping debe aplicarse después de unscale en mixed precision;
- ocultar `NaN` con clipping no sustituye investigar LR/datos/numerics.

Tests: overfit secuencia mínima; causal shift correcto; PAD no cambia loss; permutar tiempo empeora tarea orden-dependiente; estado reseteado entre ejemplos; step-by-step coincide con secuencia batched; logits/loss finitos para longitud límite.

### GRU

**[DLAI]** Reset/relevance y update gates controlan cuánto pasado usa el candidato y cuánto estado conserva la unidad, mitigando pérdida de información frente a RNN vanilla.

```text
r_t = σ(W_r x_t + U_r h_{t-1}+b_r)
u_t = σ(W_u x_t + U_u h_{t-1}+b_u)
h~_t = tanh(W_h x_t + U_h(r_t⊙h_{t-1})+b_h)
h_t = u_t⊙h_{t-1} + (1-u_t)⊙h~_t       # convención; algunas APIs invierten u
```

Documentar convención de gate/orden de pesos al portar checkpoints. GRU **mitiga**, no elimina, dependencias largas; cuesta más por step que vanilla. Comparar RNN/GRU/LSTM por task metric, tokens/s, memoria, p99 y estado serializado.

### Deep y bidirectional RNN

**[DLAI]** Stacking aumenta profundidad; bidirectional concatena estados que recorren izquierda→derecha y derecha→izquierda.

**[IMPL]** Una BiRNN usa tokens futuros: apropiada para clasificación/NER offline con secuencia completa, inválida para generación causal o streaming sin lookahead. Declarar presupuesto de lookahead. Más capas no garantizan dependencias “más abstractas”; requieren evidencia, regularización, residual/normalization y profiling.

`scan` expresa la recurrencia para que el framework compile/vectorice donde pueda, pero las dependencias temporales siguen limitando paralelismo intra-secuencia. No prometer GPU-parallel total; medir kernels, padding waste y batch shapes.

### Generación y evaluación

Para char/word LM, separar teacher-forced loss de generación libre. Decoding greedy/sampling/temperature altera calidad sin cambiar pesos. Definir seed, max length, EOS, repetition guard y filtros.

Reportar perplexity comparable, next-token accuracy sólo como diagnóstico, muestras cegadas, diversidad/repetición, memorization y métricas de la aplicación. Un generador de caracteres pedagógico no demuestra summarization o autocomplete productivo sin datos/evals propios.

### Error analysis de deep sentiment

Trazar texto→tokens/mask→estados/normas→logits→decoding. Categorías: tokenizer/OOV; padding/state leakage; dependencia larga; gradiente/numerics; exposición train–inference; corpus/domain; repetición/no EOS; bidirectional leakage; etiqueta/benchmark.

Comparar cada categoría contra baseline n-gram/mean. Si RNN gana sólo por memorizar duplicados, corregir split. Convertir errores de secuencia, límites y state reset en tests permanentes.

### Contrato operativo de deep sentiment

```text
Artefactos: tokenizer/vocab, embedding, recurrent weights, state/init, decoding config
Shapes: batch/time convention + masks + hidden layers/directions documentados
Gates: quality/perplexity + slices + stability + tokens/s/memory/p99
Estado: reset/detach/tenant isolation explícitos; streaming schema versionado
Salida: logits/probabilidades no asumidas calibradas; generación con policy/version
Fallback: n-gram/mean model o no-output ante estado/schema inválido
```

**Cobertura de la semana:** 15 videos/transcripciones contrastados, 14 lecturas y 4 ejemplos de código inventariados; prácticas/quizzes/asignaciones no abiertos ni reproducidos. El conteo incluye la sección inicialmente colapsada *N-grams vs. Sequence Models*.

## 10. Curso 3 · Semana 2 — LSTMs and Named Entity Recognition

### LSTM: ruta aditiva con compuertas

**[DLAI]** LSTM agrega cell state y forget/input/output gates para aprender qué conservar, incorporar y exponer, mitigando vanishing gradients.

```text
f_t = σ(W_f[x_t;h_{t-1}] + b_f)
i_t = σ(W_i[x_t;h_{t-1}] + b_i)
g_t = tanh(W_g[x_t;h_{t-1}] + b_g)
c_t = f_t⊙c_{t-1} + i_t⊙g_t
o_t = σ(W_o[x_t;h_{t-1}] + b_o)
h_t = o_t⊙tanh(c_t)
```

**[IMPL] Correcciones:** LSTM mejora el flujo de gradiente, no elimina vanishing/exploding ni garantiza memoria larga. Gates sigmoid no tienen por qué quedar “casi 0 o 1”. Mantener gradient clipping —preferentemente por norma global con umbral medido—, monitoreo de states/gates y tests de longitud.

Variantes cambian ecuaciones (peepholes, coupled gates, projection, sin `tanh(c)`); checkpoints sólo son portables si la convención, orden de gates/shapes y biases coinciden. Comparar LSTM con GRU/Transformer según calidad, paralelismo, latencia y estado streaming.

### NER es extracción de spans, no sólo clasificación de tokens

**[DLAI]** NER localiza/clasifica personas, lugares, organizaciones, fechas y otras entidades. El pipeline pedagógico usa token IDs, padding, LSTM, dense y log-softmax por posición.

```text
tokens → ids/mask → embeddings → (Bi)LSTM → logits[tag] por token
loss = masked token cross-entropy
```

**[IMPL]** Definir ontology y esquema BIO/BIOES/BILOU. `B-` inicia entidad, `I-` continúa y `O` queda fuera; una etiqueta `B/O` aislada no representa spans multi-token robustamente. Validar transiciones y, si hace falta coherencia global, CRF/decoding restringido.

Con subwords, alinear offsets del texto original a tokens: etiquetar primer subtoken y enmascarar el resto, propagar esquema de forma definida o usar word-level aggregation. Nunca perder offsets; son necesarios para devolver spans auditables.

### Dataset y batching

- Split por documento/fuente/tiempo; evitar frases duplicadas y la misma entidad memorizada en todos los splits.
- Preservar casing/puntuación si ayudan a entidades; no importar preprocessing de sentimiento.
- Padding de inputs y labels con valores distintos; mask en loss y métricas.
- Bucketing por longitud reduce padding; batch como potencia de 2 es una heurística de hardware, no garantía. Medir throughput/memoria.
- Mapear OOV/subwords y etiquetas con schemas versionados; validar longitudes/offsets antes del batch.
- Revisar calidad de anotación y agreement; las fronteras y tipos pueden ser ambiguos.

### Loss y numerics

Usar `cross_entropy(logits, labels, ignore_index=PAD_LABEL)` o log-softmax + NLL estable; no aplicar softmax antes de cross-entropy. Normalizar por tokens válidos o por secuencia según objetivo y documentarlo. Class weights/focal loss pueden ayudar a tipos raros, pero deben evaluarse por span y calibración.

Tests:

- PAD no altera loss/metric;
- offsets reconstruyen exactamente substring original;
- BIO transitions válidas o reparadas de forma determinista;
- oración/entidad multi-token manual produce spans correctos;
- batch individual/batched coincide en eval;
- longitudes límite, Unicode y truncation no cortan entidades silenciosamente;
- label permutation rompe desempeño;
- checkpoint conserva tag/id mappings.

### Evaluación correcta

**[DLAI]** Se enseña token accuracy con mask. **[IMPL]** `O` suele dominar, por lo que accuracy puede ser alta aunque no se extraiga nada. La métrica primaria es exact span-level precision/recall/F1 con tipo y frontera correctos; añadir:

- relaxed/partial match sólo como diagnóstico separado;
- per-entity-type F1 y soporte;
- entity-level confusion, boundary-only errors;
- OOV/unseen-entity, longitud, idioma, fuente y tiempo;
- latencia, throughput, memoria y truncation/coverage.

Evaluar generalización con entidades no vistas, no sólo nuevas frases con nombres memorizados. Definir manejo de nested/overlapping entities; BIO plano no puede representar todas.

### Error analysis y uso responsable

Taxonomía: frontera; tipo; entidad perdida/espuria; tokenización/subword; casing/puntuación; OOV; label noise; long context/truncation; nested entity; domain/temporal drift; invalid tag transition. Revisar ejemplos por volumen×costo y congelarlos como regresión.

El curso menciona NER+sentimiento de noticias para trading. **[IMPL]** Es sólo una idea de aplicación, no evidencia de ventaja financiera. NER/sentimiento de noticias pertenece al path analítico asíncrono, nunca al hot path atómico del order book; antes de cualquier acción exige event-time alignment, deduplicación, causal backtest, costos/slippage, controles de riesgo y aprobación. Una entidad detectada no prueba impacto ni dirección del mercado.

### Contrato operativo de LSTM/NER

```text
Entrada: texto + locale + tokenizer/ontology/model versions
Salida: entities[{type,start,end,text,score}] + offsets/version
Artefactos: tokenizer, tag map, model, decoding constraints, thresholds
Gates: exact span F1/per-type + unseen/OOV/slices + p99/throughput
Privacidad: redactar/autorizar entidades sensibles; tenant isolation
Fallback: reglas/diccionario o revisión; abstener ante truncation/schema inválido
```

**Cobertura de la semana:** 8 videos/transcripciones contrastados, 8 lecturas y 1 ejemplo de código inventariados; práctica/quiz/asignación no abiertos ni reproducidos.

## 11. Curso 3 · Semana 3 — Siamese Networks and Metric Learning

### Encoder compartido y espacio comparable

**[DLAI]** Dos inputs pasan por la **misma** subred/pesos para producir vectores comparables; cosine y un threshold deciden duplicado/no duplicado. El ejemplo usa embedding+LSTM, pero el patrón admite otros encoders.

```text
u = normalize(f_θ(q1))
v = normalize(f_θ(q2))
s(q1,q2) = u·v                         # cosine si norma L2=1
ŷ = 1[s > τ]
```

**[IMPL]** Implementar un solo encoder invocado dos veces, no dos copias que puedan divergir. En eval desactivar dropout/estado mutable; padding/mask y pooling deben ser idénticos. Normalizar con epsilon y definir zero vector. La simetría `s(a,b)=s(b,a)` es esperable para duplicados, pero no para relaciones dirigidas como entailment.

### Triplet loss

Con anchor `a`, positive `p`, negative `n` y margen `α`:

```text
L_triplet = max(0, s(a,n) - s(a,p) + α)
```

El objetivo exige que el positivo supere al negativo al menos por `α`; un negativo no tiene que alcanzar cosine `-1`. Elegir `α` en dev y monitorear proporción de triplets activos. Loss cero masivo puede significar tarea resuelta **o** negativos demasiado fáciles.

Si se usa distancia, invertir signos/inecuaciones consistentemente. Unit test manual debe comprobar que acercar `a,p` o alejar `a,n` no aumenta loss. L2-normalization cambia geometría y gradientes; versionarla.

### In-batch negatives y hard-negative mining

**[DLAI]** Organizar dos batches de pares positivos alineados hace que la diagonal de `S=UVᵀ` sean positivos y off-diagonals candidatos negativos. La pérdida modificada combina mean negative y closest/hard negative.

**[IMPL]** La condición “ninguna off-diagonal es duplicada” debe validarse por grupos/labels; en datos reales aparecen falsos negativos, paráfrasis no anotadas y preguntas del mismo cluster. Un false negative duro empuja semántica correcta a separarse.

Pipeline seguro:

1. deduplicar/group split por intención/cluster;
2. construir batches sin positivos cruzados conocidos;
3. enmascarar diagonal, ids iguales y positivos múltiples;
4. elegir random/semi-hard/hard negatives;
5. inspeccionar manualmente los más duros;
6. refrescar mining con versiones de encoder controladas.

Hardest negative puro puede amplificar label noise/outliers y desestabilizar. Comparar semi-hard, top-k promedio y memory bank; registrar tasa de falsos negativos y active-triplet rate. El término mean+closest del curso es una estrategia concreta, no una loss universalmente óptima; contrastar contrastive loss/InfoNCE/multiple-negatives ranking loss según batch y objetivo.

### One-shot y serving de retrieval

**[DLAI]** Aprender una función de similitud permite comparar una clase/identidad nueva sin reentrenar un softmax de `K+1` clases.

**[IMPL]** “One-shot” describe enrolamiento/reconocimiento con pocas referencias, no que el encoder se haya entrenado con un solo ejemplo total. Para duplicados a escala:

```text
offline: encode corpus → vector index
online: encode query → ANN top-k → optional cross-encoder/rules → threshold/action
```

Precomputar embeddings, medir exact-vs-ANN recall, actualizar/tombstone ids y versionar encoder+índice. El bi-encoder aporta velocidad; un cross-encoder puede rerankear pares difíciles con más precisión. Añadir lexical/hybrid retrieval para ids, nombres y coincidencias exactas.

### Threshold y evaluación

`τ` se selecciona en dev según costos/capacidad, no por apariencia del cosine. Reportar ROC/PR, precision/recall/F1 al threshold operativo, false-positive/false-negative costs, calibration si el score se convierte en probabilidad, y Recall@k/MRR para candidate retrieval.

Splits por cluster/pregunta/usuario/tiempo: ninguna paráfrasis del mismo grupo debe cruzar train/test. Evaluar hard negatives: alta superposición léxica/no duplicado, baja superposición/duplicado, negación, entidades/números, idioma y preguntas ambiguas.

Tests:

- ramas comparten exactamente parámetros/gradientes;
- swap de inputs conserva score cuando la tarea es simétrica;
- diagonal positiva y masks correctas en `S`;
- ningún positive conocido se usa como negative;
- loss manual/gradientes finitos; zero-vector seguro;
- batch size/candidate composition no cambia inferencia;
- serializar encoder/índice conserva top-k/threshold contract.

### Error analysis y seguridad

Categorías: tokenizer/OOV; pooling/truncation; false negative/label noise; entity/number mismatch; negación/dirección; domain/language; threshold/calibration; ANN miss; corpus ausente; duplicate cluster leakage. Localizar si el relevante falta en candidate set o el score/reranker lo ordena mal.

El ejemplo de autenticar firmas es sólo analogía de metric learning. Un sistema biométrico real exige threat model, anti-spoof/liveness, calidad de captura, ataques adversariales, fairness, privacidad, tasas FAR/FRR y revisión regulatoria; cosine+threshold no basta.

### Contrato operativo de redes siamesas

```text
Artefactos: tokenizer, shared encoder, normalization, ANN, threshold, reranker
Datos: pair/cluster lineage, negative-mining policy, false-negative audit
Calidad: pair PR/F1 + retrieval Recall@k/MRR + hard slices + calibration
Operación: encode/search/rerank p99, QPS, memory, freshness, rebuild/rollback
Salida: candidate ids/scores/version; score no es probabilidad sin calibración
Fallback: lexical/hybrid/review o crear nueva pregunta si no hay evidencia
```

**Cobertura de la semana:** 10 videos/transcripciones contrastados, 8 lecturas y 3 ejemplos de código inventariados; práctica/quiz/asignación no abiertos ni reproducidos.

### Cierre del Curso 3

El curso escala desde dense baselines a recurrence, gating y metric learning. La lección transversal es elegir qué información debe persistir —estado temporal o geometría compartida— y probar explícitamente cuándo se pierde.

**Gate alcanzado:** 33 videos contrastados de punta a punta; 30 lecturas y 8 ejemplos de código inventariados; 6 evaluaciones certificables contabilizadas pero no abiertas ni reproducidas. El inventario anterior de 23 videos fue corregido al descubrir la sección colapsada de la Semana 1; no queda tratada como material suplementario invisible.

# Parte IV — Attention Models

## 12. Curso 4 · Semana 1 — Neural Machine Translation with Attention

### Seq2seq y el cuello de botella

**[DLAI]** Un encoder recurrente convierte source variable en estado; un decoder autoregresivo genera target. Usar sólo el último estado fuerza toda la fuente a una memoria fija y degrada secuencias largas. Atención permite que cada paso decoder combine estados del encoder según relevancia.

```text
e_{i,j} = score(s_{i-1}, h_j)
α_{i,:} = softmax(e_{i,:} + source_mask)
c_i = Σ_j α_{i,j} h_j
p(y_i|y_<i,x) = decoder(s_{i-1}, y_{i-1}, c_i)
```

**[IMPL]** Atención alivia el bottleneck, pero **sí retiene** encoder states y tiene costo/memoria aproximados `O(T_target·T_source)` para cross-attention denso. No describirla como solución de memoria gratuita. Attention weights son variables internas útiles para diagnóstico, no explicaciones causales garantizadas ni alineamientos lingüísticos perfectos.

### Query, key, value y scaled dot-product

```text
Attention(Q,K,V,M) = softmax(QKᵀ/√d_k + M) V
```

**[DLAI]** Query busca, keys se comparan, values se agregan; masks excluyen padding. **[IMPL]** La división por `√d_k` controla la escala/varianza de dot products para evitar softmax saturado; no es una “constante de regularización”. Q/K/V suelen surgir de proyecciones aprendidas y no tienen que ser iguales.

Tests: shapes batch/head/time; mask deja probabilidad cero en PAD; cada fila válida suma uno; todo-masked row tiene fallback definido sin `NaN`; resultado manual diminuto; causal/source masks no se invierten; mixed precision coincide dentro de tolerancia.

### Datos NMT

- Pares source–target alineados, deduplicados y con licencias/PII controlados.
- Split por documento/origen/tiempo; evitar traducciones paralelas o near-duplicates cruzadas.
- Tokenizers/vocabularios por idioma o compartidos conscientemente; conservar idioma y offsets.
- BOS/EOS/PAD distintos, masks y target shift correctos.
- Bucketing por longitudes reduce padding; truncation nunca debe cortar silenciosamente información crítica.
- Evaluar rare words, nombres, números, formato, dialectos, code-switching y dominio.

One-hot es sólo explicación pedagógica: implementar embedding lookups. Pretrained no implica adecuación; medir data/locale/domain coverage.

### Teacher forcing y exposición

Durante train:

```text
decoder_input = [BOS] + target[:-1]
target_label  = target
loss = masked_cross_entropy(logits, target_label)
```

**[DLAI]** Alimentar gold previous tokens acelera/estabiliza aprendizaje temprano. **[IMPL]** En inferencia el modelo consume sus propias salidas: aparece exposure bias. Usar outputs progresivamente se llama comúnmente *scheduled sampling*; no es sinónimo general de curriculum learning y puede introducir un objetivo inconsistente. Comparar teacher forcing estándar, label smoothing, data/sequence-level training o fine-tuning de preferencias sólo con evals finales; no adoptar una técnica por nombre.

### Decoding

```text
greedy: y_t = argmax log p(y_t|prefix,x)
sampling: p_T = softmax(logits/T)
beam: conservar B prefijos por log-score acumulado + length/coverage policy
```

**[IMPL] Correcciones:** temperatura no está restringida a `[0,1]`; `T>0`, `T→0` aproxima greedy, `T=1` conserva distribución y `T>1` la aplana. Beam `B=1` equivale a greedy, pero beam mayor no garantiza mejor calidad. Usar log-probs, hypotheses finalizadas, EOS, max length, length penalty y batching/cache medidos. Evitar promedio de probabilidad ingenuo que favorezca patologías.

Comparar decoding manteniendo pesos/dataset fijos. Medir calidad, p99, memoria, tokens/s, repetición, longitud y fallos de terminación. Para traducción factual, sampling creativo rara vez debe ser default.

### BLEU, ROUGE y evaluación

**[DLAI]** Se presentan overlap n-gram, clipping, BLEU/ROUGE y sus límites semánticos/estructurales.

**[IMPL] Corrección completa:** BLEU estándar combina modified n-gram precisions mediante media geométrica y **brevity penalty**, normalmente a nivel corpus; no es sólo unigram precision. Su nombre es *BiLingual Evaluation Understudy*. ROUGE-N se usa históricamente en summarization y puede reportar recall/precision/F1 según implementación. Combinar “BLEU como precision” y “ROUGE como recall” en un F1 inventado no es una métrica BLEU/ROUGE estándar.

Fijar implementación, tokenización, case, smoothing, n-order y número de referencias; no comparar scores de protocolos distintos. Añadir:

- exactitud de nombres, números, negación, terminología y formato;
- chrF u otra métrica de caracteres para morfología;
- métrica semántica aprendida sólo validada contra humanos;
- evaluación humana ciega de adecuación, fluidez, errores críticos y preferencias;
- slices por longitud/idioma/dominio y significancia pareada/bootstrap.

Ninguna métrica de overlap prueba fidelidad, ausencia de toxicidad o seguridad.

### Minimum Bayes Risk

**[DLAI]** Generar múltiples candidates, compararlos y elegir el consenso por mayor similitud/promedio es una aproximación práctica. **[IMPL]** El nombre correcto es **Minimum Bayes Risk**, no “Base/Bias Risk”. Formalmente elige la hipótesis con menor pérdida esperada bajo una distribución posterior; el consenso entre samples con ROUGE es una aproximación dependiente de sample y utility.

```text
ŷ = argmin_y E_{y'~p(.|x)}[loss(y,y')]
```

Más samples aumentan compute; controlar seed, diversidad, duplicates y estimator variance. MBR no garantiza contextualidad/veracidad si modelo, samples o utility comparten el mismo sesgo.

### Error analysis NMT

Taxonomía: omisión; adición/hallucination; mistranslation; nombre/número; negación; agreement/morphology; word order; terminology; repetition/EOS; source truncation; tokenization; domain/dialect; metric blind spot. Separar model error de decoding error comparando gold-prefix/teacher-forced, greedy, beam y MBR.

Guardar source/reference/candidate, token versions, log-score, attention/masks, decoding config y metric breakdown. Priorizar errores críticos aunque sean raros; convertir nombres/números/negación/formato en suites deterministas.

### Contrato operativo de NMT con atención

```text
Artefactos: source/target tokenizers, model, masks, decoding + terminology config
Calidad: corpus metrics + human critical-error eval + language/domain slices
Operación: encode/decode p99, tokens/s, beam memory, max lengths, batching
Seguridad: PII, toxicidad, prompt/data injection si traduce contenido no confiable
Salida: translation + model/decoding versions; score no equivale a confianza
Fallback: glossary/rule/human review/abstención para dominio crítico
```

**Cobertura de la semana:** 14 videos/transcripciones contrastados, 3 lecturas y 3 ejemplos de código inventariados; práctica/quiz/asignación no abiertos ni reproducidos.

## 13. Curso 4 · Semana 2 — Transformer Text Summarization

> Para implementación interna, KV cache, arquitecturas modernas, fine-tuning y serving, la autoridad es `PYTORCH_TRANSFORMERS_GENERATIVE_MODELS.md`. Aquí se conserva el delta de NLP/summarization y las correcciones del curso.

### Qué cambia frente a RNN

**[DLAI]** Self-attention elimina la recurrencia del cómputo de entrenamiento, habilita procesamiento paralelo de posiciones y acceso directo entre tokens. Positional encodings aportan orden; residuals, LayerNorm y position-wise FFN completan cada bloque.

**[IMPL]** Transformer no es universalmente “más rápido y con menos memoria”. Self-attention densa cuesta `O(T²)` en scores/memoria y `O(T²D)` aproximadamente en compute; RNN cuesta lineal en longitud pero es secuencial. El ganador depende de longitud, batch, hardware, kernel, precision, KV cache y régimen train/decode. Important information no “debe” perderse siempre en RNN; es una limitación empírica/optimization, no ley absoluta.

### Bloque mínimo y masks

```text
X = token_embedding(ids) + positional_encoding(positions)
H = residual_norm(X, MHA(X,X,X, padding_or_causal_mask))
Y = residual_norm(H, FFN(H))
```

En multi-head:

```text
head_i = Attention(QW_i^Q, KW_i^K, VW_i^V)
MHA = concat(head_1..head_h) W^O
```

**[IMPL] Correcciones:**

- causal mask bloquea posiciones **futuras/a la derecha**, no información “hacia la izquierda”;
- decoder-only Transformer es una familia arquitectónica; no es automáticamente GPT-2, que fija tokenizer, tamaños, training data/objective y configuración concretos;
- FFN comparte pesos entre posiciones dentro de una capa, no necesariamente entre capas;
- pre-norm y post-norm son variantes distintas; no mezclar orden al portar checkpoints;
- `d_model`, heads y profundidad son hiperparámetros con restricciones divisibles/compute, no reglas como “hasta 10k”.

Tests: padding/causal masks combinados; cambiar token futuro no altera logits pasados; posición cambia representación; heads concat/shapes; residual dims; FFN por posición; reference implementation; gradientes finitos; attention all-masked segura.

### Summarization decoder-only

**[DLAI]** Concatenar artículo, separador/EOS y resumen; entrenar next-token prediction, enmascarando loss del artículo para que el objetivo principal recaiga en summary.

```text
input  = [BOS] + article + [SEP] + summary[:-1]
labels = article_shifted + [SEP] + summary
loss_mask = 0 sobre prompt/article/PAD; 1 sobre summary incl. EOS
```

**[IMPL]** Prevenir off-by-one: label en posición `t` corresponde al próximo token. El summary no puede aparecer en el contexto de tokens cuya loss se calcula causalmente de manera incorrecta. Tests de máscara deben cambiar article weights sin modificar qué targets summary se evalúan.

Ponderar article LM loss puede regularizar con poco dato, pero también desvía capacidad y hace el objetivo dependiente del dominio; tratar `λ_article` como ablation, no recomendación universal. Comparar pretrained fine-tuning antes de training from scratch.

### Long documents y producto

Definir primero summary contract: extractive/abstractive; longitud/compresión; audiencia; hechos obligatorios; citas; idioma; formato; qué omitir. Para textos mayores al context window:

- truncation sólo con medición de información perdida;
- chunk/map-reduce puede perder relaciones globales y propagar errores;
- retrieve-then-summarize necesita eval de retrieval separada;
- hierarchical/long-context model requiere profiling y benchmark propio.

No confundir ventana aceptada con contexto efectivamente usado. Instrumentar posición de hechos y recall por sección.

### Decoding y fidelidad

Usar greedy/beam de baja entropía como baseline; sampling sólo si diversidad es requisito y con fact-checking. Definir EOS, min/max length, length/repetition penalty y no-repeat n-gram como parámetros de decoding, no arreglos de factualidad.

Un resumen fluido puede inventar, fusionar entidades, negar mal u omitir información crítica. Separar:

- **relevance/coverage:** incluye lo importante;
- **faithfulness/attribution:** cada afirmación está soportada;
- **coherence/fluency:** se lee bien;
- **compression/style:** cumple longitud/formato;
- **safety/privacy:** no revela ni distorsiona contenido sensible.

ROUGE/BLEU sólo miden overlap parcial. Añadir QA/fact extraction consistency con validación humana, entity/number checks, citas/entailment calibrados y revisión ciega. Métricas automáticas aprendidas también fallan; no son árbitro único.

### Error analysis de summarization

Para cada artículo guardar reference/candidate, prompt/template, truncation/chunks/retrieval, decoding, citas y scores descompuestos. Categorías: hallucinación; atribución errónea; entidad/número/fecha; negación; omisión; redundancia; incoherencia; longitud/formato; toxicidad/PII; truncation/retrieval; benchmark/reference defect.

Localizar si el hecho no llegó al contexto, llegó pero se ignoró o fue contradicho al generar. La intervención cambia: retrieval/chunking, data/objective, decoding, verifier o abstención. Convertir hechos críticos en evals y requerir revisión humana para dominios de alto riesgo.

### Contrato operativo de summarization

```text
Entrada: documento(s) + summary contract + model/tokenizer/prompt versions
Artefactos: preprocess/chunker/retriever opcional, model, decoding, verifier
Calidad: coverage + faithfulness + entity/number + human + slices
Operación: context utilization, p99, tokens/s, memory/KV, cost, truncation rate
Salida: summary + evidence/citations cuando se exijan + versions
Fallback: extractive template, partial result o abstención/revisión
```

**Cobertura de la semana:** 10 videos inspeccionados/contrastados, 5 lecturas y 3 ejemplos de código inventariados; práctica/quiz/asignación no abiertos ni reproducidos.

## 14. Curso 4 · Semana 3 — Transfer Learning and Question Answering

### Primero definir qué significa “responder”

**[DLAI]** Se distinguen QA con contexto —extraer/generar desde un pasaje— y closed-book, que responde desde parámetros. BERT se aplica a extractive QA; T5 formula tareas como text-to-text. Hugging Face ejemplifica uso y fine-tuning de checkpoints.

**[IMPL]** Contratos distintos:

| Tipo | Evidencia permitida | Salida | Riesgo central |
|---|---|---|---|
| extractive QA | contexto dado | span exacto | respuesta ausente/mal span |
| generative grounded QA | contexto recuperado/dado | texto + citas | hallucination/attribution |
| closed-book QA | parámetros | texto | conocimiento obsoleto/no verificable |
| RAG | corpus + retrieval | texto + evidencia | candidate miss + generation error |

Para software profesional, closed-book no debe ser default cuando se exige trazabilidad o frescura. Implementar no-answer/abstención y separar retrieval de reader/generator.

### Transfer learning: selección y adaptación

**[DLAI]** Pretraining autosupervisado permite usar features congeladas o fine-tunear parámetros en downstream, reduciendo tiempo/datos y mejorando resultados potencialmente.

**[IMPL]** Transfer puede ser negativo por dominio, idioma, tokenizer, licencia, sesgos o objetivo. Más datos/modelo no garantiza mejor desempeño: importan calidad, cobertura, duplicados y compute. Decision ladder:

1. baseline/checkpoint task-specific sin entrenamiento;
2. linear head/frozen encoder;
3. PEFT/partial unfreeze;
4. full fine-tune;
5. continued pretraining de dominio;
6. scratch sólo con datos/compute/razón demostrada.

Evaluar cada escalón con seeds, dev/test congelados, costo/latencia/memoria y error analysis. Model size/reputación/leaderboard no sustituyen benchmark local.

### BERT: objetivo y fine-tuning

**[DLAI]** BERT encoder bidireccional suma token, position y segment embeddings. El pretraining original selecciona 15% de posiciones para MLM (80% `[MASK]`, 10% random, 10% unchanged) y combina Next Sentence Prediction. Fine-tuning adapta todos los parámetros y una head.

**[IMPL] Correcciones:**

- BERT original usa **Masked Language Modeling**, no “multi-mask language model”.
- NSP es parte de BERT original, no requisito universal; modelos posteriores lo modifican/eliminan.
- ELMo combina language models forward/backward contextuales; no equivale simplemente a “predecir el centro con BiLSTM”.
- BERT encoder no es por sí solo un generador/summarizer autoregresivo; la head/arquitectura debe coincidir con la tarea.
- El vector `[CLS]` no es universalmente el mejor sentence embedding sin training/pooling adecuado.

Para extractive QA, concatenar `[CLS] question [SEP] context [SEP]`, producir logits start/end y escoger span válido:

```text
score(i,j) = start_logit[i] + end_logit[j]
subject to context_mask, i≤j, length≤L_max
```

Con contextos largos, sliding windows+stride; mapear offsets al texto original, combinar candidatos entre chunks y deduplicar. Incluir score `[CLS]`/head explícita para no-answer y ajustar threshold en dev. No permitir spans en question/PAD/special tokens.

### T5 y multitask text-to-text

**[DLAI]** T5 usa encoder–decoder, span corruption con sentinel tokens y prefixes textuales para unificar clasificación, QA, traducción y summarization. Mezcla tareas con estrategias proportional/equal/temperature y estudia fine-tuning/adapters.

**[IMPL]** Prefix sólo funciona porque el modelo fue entrenado con esa convención; texto arbitrario no crea una capacidad. T5 span corruption no es idéntico al MLM token individual de BERT. Encoder–decoder no es universalmente superior; elegir según task/latency.

En multitask, controlar sampling por tarea, tamaño/gradientes, interferencia negativa y métricas separadas. Un checkpoint distinto por tarea al reportar no demuestra un único artefacto óptimo simultáneo. Adapters/PEFT reducen parámetros entrenables, no necesariamente memoria de activaciones o latencia base.

### Tokenización subword y alineación

SentencePiece/BPE permiten vocabulario acotado y composición OOV, pero no eliminan unknown/coverage ni significado. Versionar normalización, special tokens y vocab con checkpoint. Para QA:

- `offset_mapping` obligatorio para spans;
- distinguir question/context/special tokens;
- truncation side y stride testeados;
- Unicode/whitespace/casing reproducen substrings;
- fast/slow tokenizer no se sustituyen sin parity test.

Cambiar tokenizer invalida embedding rows, offsets, dataset features e índice.

### Hugging Face como dependencia, no magia

**[DLAI]** Model/dataset hubs, model cards, tokenizers, pipelines y training abstractions aceleran prototipos/fine-tuning. **[IMPL]** APIs, defaults y cifras del curso son históricas; consultar documentación primaria actual al implementar.

Contrato de supply chain:

- fijar model/dataset/library revision y hashes; no depender de `latest`;
- revisar licencia, provenance, dataset/model cards, intended use y limits;
- evitar remote custom code o aislar/revisar explícitamente; nunca habilitar confianza ciega;
- usar formatos seguros de weights cuando estén disponibles; escanear artefactos;
- cache/download en build controlado, no hot path; soportar offline/rollback;
- registrar tokenizer/config/preprocess/postprocess además de weights;
- benchmark del checkpoint exacto; pipeline default sólo smoke test.

Una línea de código puede ocultar truncation, labels, thresholds, device, batching y postprocessing. Expandir pipeline antes de producción y escribir parity tests.

### GLUE/benchmarks

**[DLAI]** GLUE agrupa tareas NLU y habilita comparación model-agnostic. **[IMPL]** Un benchmark público mide datasets/protocolos específicos, puede saturarse o contaminarse y no prueba QA/RAG/producto propio. No elegir modelo por score agregado. Reportar dataset por dataset, varianza, compute y benchmark local con slices/errores.

### Evaluación QA

Extractive:

- Exact Match normalizado y token F1;
- no-answer precision/recall y threshold curve;
- answerable/unanswerable, long context, answer position/length, adversarial distractors;
- entity/number/date/negation y multilingual/domain slices.

Grounded/generative:

- retrieval Recall@k antes de answer quality;
- correctness, faithfulness/citation entailment, completeness y abstention;
- exact structured checks cuando sea posible;
- evaluación humana ciega y costo de error.

No usar LLM judge sin calibración humana, agreement y adversarial audit. Latencia se separa: tokenize/embed, retrieve, reader/generate, verify.

### Error analysis de question answering

Árbol causal:

```text
¿hay respuesta válida en corpus/contexto?
 ├─ no → corpus/coverage/abstención
 └─ sí → ¿llegó al chunk/candidate set?
          ├─ no → chunking/retrieval/filter
          └─ sí → ¿reader/generator la seleccionó y atribuyó?
                   ├─ no → model/span/decoding
                   └─ sí pero salida mala → postprocess/citation/policy
```

Registrar query, gold/evidence, candidates, offsets, start/end logits o generation trace, thresholds, versions y latencias. Categorías: ambiguity; unanswerable; multi-hop; coreference; long context; tokenizer/span; retrieval miss; distractor; stale/contradictory source; hallucination; citation; domain/language; label defect.

### Contrato operativo de question answering

```text
Entrada: question + authorized context/corpus + versions/tenant
Artefactos: tokenizer, reader/generator, retriever/index, thresholds, verifier
Seguridad: ACL pre/post retrieval, injection isolation, PII, source allowlist
Calidad: retrieval + answer + faithfulness + no-answer + slices + human
Operación: stage p99/QPS/memory/cost/freshness, caching and rollback
Salida: answer, evidence/offsets/citations, confidence semantics, versions
Fallback: no-answer, search results o revisión; nunca fabricar evidencia
```

**Cobertura de la semana:** 15 videos/transcripciones contrastados, 11 lecturas y 3 ejemplos de código inventariados; práctica/quiz/asignación no abiertos ni reproducidos. Incluye la sección inicialmente colapsada de Hugging Face, revisada sin pulsar `Mark As Complete`.

### Cierre del Curso 4 y de la especialización

El curso conecta atención recurrente, Transformer y transfer learning con tres tareas: NMT, summarization y QA. La habilidad senior es descomponer representación, contexto, decoding, evidencia, métrica y operación; no tratar “Transformer” como solución indivisible.

**Gate del Curso 4:** 39 videos contrastados, 19 lecturas y 9 ejemplos de código inventariados; 3 evaluaciones certificables contabilizadas pero no abiertas ni reproducidas.

**Gate de Natural Language Processing Specialization:** 4 cursos, 14 semanas, 412 unidades, **176 videos contrastados**, 137 lecturas y 38 ejemplos de código inventariados, 20 evaluaciones certificables excluidas. La cuenta/sesión permaneció en 0 % y no se modificó el plan.

## Estado de incorporación

- [x] Catálogo, instructores, duración, cursos y semanas verificados.
- [x] 412 unidades y 20 asignaciones contabilizadas sin abrir evaluaciones.
- [x] Curso 1 contrastado e integrado.
  - [x] Semana 1 — Logistic Regression.
  - [x] Semana 2 — Naïve Bayes.
  - [x] Semana 3 — Vector Space Models.
  - [x] Semana 4 — Machine Translation and Document Search.
- [x] Curso 2 contrastado e integrado.
  - [x] Semana 1 — Autocorrect and Minimum Edit Distance.
  - [x] Semana 2 — POS Tagging and Hidden Markov Models.
  - [x] Semana 3 — Autocomplete and Language Models.
  - [x] Semana 4 — Word Embeddings with Neural Networks.
- [x] Curso 3 contrastado e integrado.
  - [x] Semana 1 — Deep Sentiment and Recurrent Language Models.
  - [x] Semana 2 — LSTMs and Named Entity Recognition.
  - [x] Semana 3 — Siamese Networks.
- [x] Curso 4 contrastado e integrado.
  - [x] Semana 1 — Neural Machine Translation.
  - [x] Semana 2 — Text Summarization.
  - [x] Semana 3 — Question Answering.
- [x] Auditoría transversal, compresión y contrato final para Codex.

## 15. Auditoría transversal y contrato final para Codex

### Qué documento gobierna cada decisión

| Decisión | Autoridad primaria | Complemento cuando haga falta |
|---|---|---|
| representación léxica, edit distance, n-gramas, HMM, embeddings, búsqueda y NLP aplicado | este documento | matemática/ML clásico para derivaciones |
| optimización, regularización, RNN/GRU/LSTM y atención como fundamentos | `DEEP_LEARNING_ANDREW_NG.md` | este documento para tareas NLP |
| implementación, fine-tuning, profiling y serving de Transformers | `PYTORCH_TRANSFORMERS_GENERATIVE_MODELS.md` | este documento para tokenización, QA y métricas de lenguaje |
| scoping, datos, error analysis, despliegue, drift y monitoreo | `ML_PRODUCTION_LLMOPS_EVALUATION.md` | este documento para slices y fallos NLP |
| agentes, tools, memoria y evaluación de trayectorias | `AGENTIC_AI_SOFTWARE_ENGINEERING_CODEX.md` | RAG sólo como capa de conocimiento |
| RAG, embeddings de recuperación y vector stores | Parte V de este documento | ML Production para operación y seguridad; knowledge graphs como extensión futura |

No duplicar teoría entre archivos. Cuando dos documentos se solapan, usar la autoridad para decidir y el complemento para implementar u operar.

### Pipeline invariante

```text
scope/impacto
→ contrato de entrada, salida, idioma, evidencia y latencia
→ split por entidad/tiempo antes de aprender vocabulario o features
→ normalización/tokenizer versionados y adecuados a la tarea
→ baseline trivial + baseline clásico fuerte
→ modelo mínimo que supere el baseline en slices relevantes
→ evaluación offline descompuesta
→ error analysis causal con artefactos reproducibles
→ gate de seguridad, costo, rendimiento y rollback
→ shadow/canary
→ monitoreo de datos, calidad, abstención y negocio
→ nuevos errores etiquetados regresan al dataset de evaluación
```

### Tensiones resueltas

| Regla aparentemente universal | Resolución senior |
|---|---|
| quitar puntuación, stopwords o casing | sólo si una ablación demuestra que no son señal; negación, entidades, código y sentimiento suelen necesitarlos |
| accuracy decide el modelo | usar la métrica que represente el costo de error: F1/PR, span EM, retrieval, ranking, calibración, latencia y slices |
| full softmax es el algoritmo de embeddings | es útil pedagógicamente; vocabularios grandes suelen requerir sampled objectives, batching y medición real |
| Transformer siempre es más rápido que RNN | paraleliza training, pero latencia/memoria dependen de longitud, hardware, batching, KV cache y decoding autoregresivo |
| attention explica por qué decidió el modelo | los pesos ayudan a inspeccionar, pero no constituyen una explicación causal ni evidencia suficiente |
| perplexity, BLEU o ROUGE prueban calidad | son proxies; añadir evaluación de tarea, factualidad, evidencia, slices y revisión humana calibrada |
| un checkpoint pretrained ya está listo | verificar dominio, idioma, licencia, tokenizer, provenance, seguridad, costo y benchmark local |
| más parámetros/datos siempre mejora | calidad, cobertura, deduplicación, objetivo y restricciones operativas pueden dominar el tamaño |
| el modelo debe procesar cada evento del order book | mantener IA fuera del hot path atómico determinista; consumir snapshots/features asíncronos salvo evidencia que justifique otra arquitectura |

### Protocolo de implementación que un agente debe ejecutar

1. Escribir primero el contrato: tarea, población, fuente, salida, SLA, costo de falsos positivos/negativos y política de abstención.
2. Congelar splits y construir un baseline auditable; prohibido elegir arquitectura por moda.
3. Versionar datos, normalización, tokenizer/vocabulario, modelo, thresholds y postproceso como un solo artefacto lógico.
4. Agregar pruebas de shapes, masks, padding, offsets, leakage, determinismo, casos Unicode y ejemplos límite antes de entrenar.
5. Medir loss y métrica de producto por slices; guardar seeds, configuración, checkpoint y entorno.
6. Revisar errores individuales y agruparlos por causa, no sólo por etiqueta final.
7. Cambiar una familia causal por iteración: datos, representación, objetivo, arquitectura, decoding, retrieval o threshold.
8. Convertir cada fallo crítico corregido en una prueba o caso de evaluación permanente.
9. Antes de desplegar, medir p50/p95/p99, throughput, memoria, costo, cold start, degradación y rollback.
10. En producción, registrar versiones y señales suficientes para reproducir el fallo sin guardar texto sensible innecesario.

### Árbol de diagnóstico mínimo

```text
¿la especificación y la etiqueta son correctas?
 ├─ no → corregir contrato/datos/evaluación
 └─ sí → ¿la señal necesaria existe en la entrada autorizada?
          ├─ no → mejorar datos/contexto o abstenerse
          └─ sí → ¿preproceso/tokenizer/retrieval la preservan?
                   ├─ no → reparar representación/pipeline
                   └─ sí → ¿el modelo aprende y generaliza?
                            ├─ no aprende → bugs, capacidad, loss, optimización
                            ├─ memoriza → leakage, regularización, más datos útiles
                            └─ aprende offline → shift, threshold, serving o métrica de producto
```

Para RAG/QA se inserta una pregunta adicional: **¿la evidencia correcta entró en el candidate set?** Si no, corregir corpus/chunking/retrieval antes del generador. Para secuencias, inspeccionar además masks, estados, longitudes, exposure bias y decoding.

### Perfil mínimo de carga para ahorrar tokens

| Encargo a Codex | Secciones de este archivo | Complemento |
|---|---|---|
| clasificador de texto | 3–5, 9 y 15 | ML Production al desplegar |
| autocomplete/language model recurrente | 7, 9 y 15 | Deep Learning para optimización |
| embeddings o búsqueda semántica | 5, 8 y 15 | futura parte RAG para producción |
| traducción | 5, 12 y 15 | PyTorch/Transformers para implementación moderna |
| summarization | 13 y 15 | PyTorch/Transformers + ML Production |
| QA | 14 y 15 | parte RAG cuando exista corpus recuperado |
| NER/tagging | 6, 10 y 15 | ML Production para drift/slices |

El agente debe citar sección, declarar supuestos y separar `[DLAI]` de `[IMPL]`. Si una decisión depende de versión actual de una librería, modelo, licencia o servicio, debe consultar su documentación primaria vigente: este manual fija principios y contratos, no APIs eternas.

### Gate final

- Inventario completo: **4 cursos, 14 semanas y 412 unidades**.
- Evidencia contrastada: **176 videos**, 137 lecturas y 38 ejemplos de código inventariados.
- Exclusión preservada: 20 evaluaciones certificables contabilizadas, no abiertas ni reproducidas.
- Procedencia separada; no se presentan paráfrasis como transcripción literal.
- Baselines, métricas, error analysis, pruebas, operación, seguridad y rollback están conectados a decisiones de implementación.
- La especialización NLP queda **completa y auditada**. La ampliación RAG/retrieval avanzada es la siguiente fuente, no un hueco oculto de este programa.

# Parte V — Retrieval Augmented Generation

## 16. Módulo 1 — RAG Overview

### Modelo mental y frontera del sistema

**[DLAI]** RAG combina dos procesos: recuperar información útil desde una base de conocimiento y proporcionarla en el prompt para que el LLM genere una respuesta apoyada en ese contexto. Resuelve el caso en que la información relevante es privada, reciente o demasiado especializada para asumir que quedó aprendida durante pretraining.

```text
pregunta
→ retriever(query, corpus/index)
→ documentos candidatos + scores + metadata
→ construcción del contexto/prompt aumentado
→ LLM
→ respuesta sustentada y, cuando se exija, citas
```

La interfaz puede parecer idéntica a usar un LLM directamente; internamente agrega recuperación, ensamblado de contexto y latencia. El knowledge base no es sinónimo de vector database: es la colección autorizada; el índice y el método de búsqueda son decisiones de implementación.

**[IMPL] Corrección crítica:** RAG aumenta la probabilidad de precisión y grounding, pero no lo garantiza. Un documento recuperado puede ser irrelevante, falso, obsoleto, contradictorio, mal autorizado o contener instrucciones hostiles; el modelo puede ignorarlo o citarlo incorrectamente. “Encontró texto” y “respondió correctamente” son gates distintos.

### Cuándo usarlo

**[DLAI]** Casos mostrados:

- asistencia sobre un repositorio: clases, funciones, convenciones y archivos propios;
- soporte al cliente o empleados con productos, inventario, políticas y documentación empresarial;
- medicina y derecho con material especializado/privado, donde la precisión es crítica;
- búsqueda web con resumen de resultados recientes;
- asistentes personales apoyados en correo, mensajes, calendario o documentos.

**[IMPL]** Usar RAG si el valor depende de conocimiento que cambia, debe citarse, está sujeto a permisos o no conviene incorporar en weights. No usarlo automáticamente para tareas puramente transformativas, conocimiento pequeño/estable que cabe en configuración determinista, ni para operaciones exactas que requieren SQL/API/cálculo: esas fuentes deben consultarse como tools y sus resultados pueden luego formar contexto.

### Aportes identificados de Andrew Ng

**[NG]** Andrew presenta RAG como una de las clases de aplicaciones LLM más construidas. Destaca que sus prácticas evolucionan junto con mejores modelos, context windows mayores y extracción agentiva de PDF/slides; RAG también puede ser sólo un componente intermedio de un workflow agentivo mayor.

**[NG]** En agentic RAG, el agente recibe herramientas y decide qué fuente consultar, cómo formular la búsqueda y si la evidencia recuperada basta o exige otra ronda. Esto sustituye parte de las reglas rígidas escritas por una persona y permite revisar una ruta que falló.

**[IMPL]** Flexibilidad no elimina control: limitar pasos, fuentes, tiempo y costo; registrar cada query/candidato; validar ACL antes de exponer contenido; exigir condición de parada y fallback. El agente puede recuperar iterativamente, pero no autootorgarse permisos ni convertir ausencia de evidencia en una respuesta inventada.

### LLM: capacidad útil y límite

**[DLAI]** El curso describe el LLM como generador autoregresivo de tokens: estima una distribución del siguiente token, selecciona y repite condicionándose por lo generado. Fue entrenado para producir texto probable, no para verificar verdad; por eso un texto plausible puede ser incorrecto. El contexto recuperado puede orientar la generación aun cuando no formó parte del entrenamiento.

La cantidad de contexto está limitada por compute, costo y context window. Más texto no equivale a más señal: una selección amplia aumenta recall, pero también ruido, latencia y riesgo de distraer al generador.

**[IMPL]** Versionar modelo, tokenizer, prompt y sampling. Para respuestas factuales iniciar con decoding de baja variabilidad y no usar temperatura como reparación de factualidad. Medir utilización del contexto y atribución, no sólo fluidez.

### Retriever: contrato inicial

**[DLAI]** El retriever mantiene una base de documentos, crea un índice, interpreta la consulta, calcula relevancia, ordena candidatos y devuelve una cantidad seleccionada. Debe recuperar lo pertinente y excluir lo irrelevante. Devolver todo preserva recall pero destruye precisión/context budget; devolver sólo el primer resultado puede omitir evidencia valiosa.

El campo de information retrieval antecede a los LLM. Una vector database es común para búsqueda por similitud a escala, pero no es obligatoria: búsqueda lexical, filtros, SQL, índices invertidos y combinaciones híbridas también son válidos.

**Contrato mínimo:**

```text
retrieve(
  query,
  tenant/identity,
  filters,
  top_k,
  index_version,
  as_of_time
) -> [
  {document_id, chunk_id, content_or_reference, score, metadata, source_version}
]
```

El score sólo tiene significado dentro del método/index/version; no presentarlo como probabilidad sin calibración. La autorización filtra candidatos antes de entregar contenido, y se vuelve a validar al construir la respuesta.

### Esqueleto de implementación

```python
def answer(request, identity, versions):
    query = normalize_without_losing_signal(request.question)
    candidates = retrieve(
        query=query,
        identity=identity,
        filters=request.filters,
        top_k=versions.top_k,
        index_version=versions.index,
    )
    evidence = select_within_budget(candidates, versions.context_policy)
    if not evidence_is_sufficient(evidence, request.answer_contract):
        return abstain_with_search_results(evidence)

    prompt = render_grounded_prompt(request, evidence, versions.prompt)
    draft = generate(prompt, versions.model, versions.decoding)
    checked = verify_attribution_and_policy(draft, evidence)
    return checked if checked.accepted else checked.fallback
```

Este esquema es deliberadamente independiente de proveedor. La implementación concreta debe conservar IDs/offsets, impedir que el texto recuperado reescriba instrucciones del sistema, soportar timeout/cancelación y emitir una traza por etapa.

### Evaluación y aprendizaje desde errores

Separar al menos cuatro niveles:

| Nivel | Pregunta | Señales iniciales |
|---|---|---|
| corpus | ¿existía una fuente correcta y autorizada? | coverage, freshness, duplicados, ACL |
| retrieval | ¿entró en los candidatos y quedó bien ordenada? | Recall@k, MRR/nDCG, precision/context relevance |
| context | ¿se preservó la evidencia al chunkear/ensamblar? | truncation, lost-in-middle, token budget, contradictions |
| generation | ¿la respuesta es correcta y atribuida? | correctness, faithfulness, citation support, abstention |

```text
respuesta mala
→ ¿gold/fuente válida existe?
  → no: corpus o contrato
  → sí: ¿aparece en top-k?
    → no: query/index/filter/ranking
    → sí: ¿entra al prompt sin truncarse?
      → no: chunking/context policy
      → sí: ¿el modelo la usa y cita correctamente?
        → no: prompt/model/decoding/verifier
        → sí: postproceso, UI o criterio de evaluación
```

Guardar por caso: pregunta, identidad/tenant anonimizado, corpus/index/model/prompt versions, candidatos con scores, evidencia final, tokens, respuesta, citas, decisión de abstención, latencia/costo por etapa y etiqueta causal. Los fallos corregidos se convierten en evals de regresión; el dataset de evaluación no se modifica para maquillar la métrica.

### Gate del Módulo 1

- 8 videolecciones/transcripciones contrastadas, incluida la conversación con Andrew Ng.
- 2 ejemplos de código y 2 lecturas inventariados; se extrajo el contrato técnico sin reproducir notebooks.
- 2 evaluaciones certificables contabilizadas, no abiertas ni reproducidas.
- Arquitectura, casos de uso, límites, seguridad, implementación y error analysis quedaron conectados.

## 17. Módulo 2 — Information Retrieval and Search Foundations

### Arquitectura híbrida

**[DLAI]** Un retriever moderno combina normalmente búsqueda lexical, búsqueda semántica y filtros por metadata. Lexical preserva coincidencias exactas; semántica recupera significado aunque cambien las palabras; metadata aplica condiciones rígidas. Las listas se filtran y fusionan antes de devolver top-k.

```text
query ─┬→ sparse/keyword search ─┐
       └→ dense/semantic search ─┼→ ACL + metadata filters → rank fusion → top-k
                                 └──────────────────────────────────────────────
```

**[IMPL]** ACL no es una mejora de ranking: es una restricción de seguridad. Debe aplicarse en la consulta o antes de exponer contenido, no sólo después de que el modelo lo recibió. Metadata necesita schema, tipos, timezone, null semantics e índices versionados; condiciones mal formadas deben fallar cerradas.

### Búsqueda lexical: TF-IDF como explicación, BM25 como baseline

**[DLAI]** Keyword search representa query/documentos como vectores sparse sobre vocabulario. TF pondera frecuencia local; IDF aumenta el peso de términos raros en el corpus; la matriz/inverted index permite pasar rápidamente de término a documentos. TF-IDF es un baseline fundacional.

Una forma reproducible, entre varias convenciones, es:

```text
tfidf(t,d) = tf(t,d) · log((N + 1)/(df(t) + 1))
score(q,d) = Σ[t ∈ q] tfidf(t,d)
```

**[DLAI]** BM25 añade saturación de term frequency y normalización ajustable por longitud. **[IMPL]** La fórmula y defaults varían por motor; fijar implementación y parámetros:

```text
BM25(q,d) = Σ IDF(t) · tf(t,d)(k1+1)
                         ─────────────────────────────
                         tf(t,d)+k1(1-b+b·|d|/avgdl)
```

- `k1` controla la saturación por repeticiones;
- `b` controla cuánto normaliza longitud;
- analyzer, Unicode, casing, stemming, synonyms y stopwords forman parte del índice;
- IDs, SKU, nombres, errores textuales y terminología técnica suelen favorecer lexical.

No retirar términos o puntuación por hábito. Validar con ablaciones y queries reales. BM25 es el baseline fuerte que dense/hybrid debe superar, no un artefacto anticuado que pueda omitirse.

### Búsqueda semántica y contrato de embeddings

**[DLAI]** Un embedding model transforma texto en vectores densos donde pares semánticamente positivos se entrenan para quedar cercanos y negativos para separarse. La consulta se codifica y se ordenan documentos por cosine similarity, dot product u otra distancia compatible.

```text
cos(q,d) = (q·d) / (||q|| ||d||)
```

**Invariantes de implementación:**

- query y documentos usan el modelo/modo/prefix correcto y compatible;
- no comparar vectores de modelos distintos, aunque coincida su dimensión;
- versionar checkpoint, tokenizer, pooling, normalización y dimensión;
- cambiar cualquiera de ellos exige nuevo índice o migración dual con parity gate;
- la escala de cosine/dot/L2 y el sentido mayor/menor no son intercambiables;
- benchmark local debe incluir idioma, dominio, queries cortas/largas, negación y entidades.

Semantic search resuelve sinónimos y paráfrasis, pero puede perder identificadores exactos, confundir conceptos cercanos o heredar huecos del contraste de entrenamiento. Vector proximity es una señal aprendida de similitud, no verdad ni autorización.

### Fusión por rangos

**[DLAI]** Reciprocal Rank Fusion combina listas sin necesitar que sus scores tengan la misma escala:

```text
RRF(d) = Σ[s ∈ searchers] w_s / (K + rank_s(d))
```

`K` reduce la dominancia de los primeros puestos; `w_s`/`beta` controla la contribución de cada buscador. El curso propone mayor peso semántico como punto inicial, no como óptimo universal.

**[IMPL]** Conservar ausencia de un documento en una lista, empates, deduplicación por identidad estable y trazabilidad de cada contribución. Comparar RRF con score normalization/reranker usando el mismo candidate budget. Tunar BM25, embedding, filtros, pesos, candidate-k y final-k de forma controlada; cambiar todo a la vez impide aprender del error.

### Métricas con ground truth

Para query `q`, relevantes `Rel(q)` y top-k `R_k(q)`:

```text
Precision@k = |R_k ∩ Rel| / k
Recall@k    = |R_k ∩ Rel| / |Rel|
RR(q)       = 1 / rank(primer relevante)
MRR         = mean_q RR(q)
```

Average Precision acumula `Precision@i` sólo donde el elemento `i` es relevante; MAP promedia AP entre queries. Declarar el denominador exacto —todos los relevantes, `min(k, |Rel|)` u otra convención— porque librerías difieren. Añadir nDCG cuando existan grados de relevancia; el curso no lo necesita para explicar los fundamentos.

Interpretación:

- `Recall@k`: ¿la evidencia necesaria entró al candidate set?;
- `Precision@k`: ¿cuánto presupuesto se desperdicia en ruido?;
- MAP/nDCG: ¿los relevantes están bien ordenados en general?;
- MRR: ¿aparece pronto al menos un relevante?;
- ninguna mide por sí sola si la respuesta generada es correcta.

El dataset de evaluación necesita query, corpus/index version, qrels y reglas de juicio. Incluir head/tail, no-result, multi-intent, temporal, acrónimos, IDs, lenguaje natural, idiomas, permisos y adversarial distractors. Separar dev de test; los clicks de producción son señales sesgadas, no relevancia perfecta.

### Ruta de ejecución eficiente

```python
def hybrid_retrieve(query, identity, cfg):
    filt = authorized_filter(identity, cfg.metadata_schema)
    sparse_future = search_sparse(query, filt, cfg.sparse_candidates)
    dense_future = search_dense(embed_query(query), filt, cfg.dense_candidates)
    sparse, dense = join_with_deadline(sparse_future, dense_future)
    fused = reciprocal_rank_fusion(sparse, dense, cfg.weights, cfg.rrf_k)
    return deduplicate(fused)[: cfg.final_k]
```

Precomputar índices y document embeddings; ejecutar ramas independientes en paralelo; empujar filtros al storage cuando preserve recall autorizado; limitar candidates y deadline por etapa; cachear sólo con tenant, filtros, query normalization e index version en la key. Medir p50/p95/p99, QPS, memory, cache hit, timeout y partial-result rate junto a calidad.

### Error analysis del retriever

Clasificar cada fallo antes de cambiar el modelo:

| Síntoma | Inspección | Intervención probable |
|---|---|---|
| exact ID/producto ausente | analyzer y sparse rank | preservar token, field boost, BM25 |
| paráfrasis ausente | dense rank/model/domain | embedding o query expansion |
| relevante filtrado | ACL/metadata/query plan | corregir datos/regla; jamás relajar seguridad por score |
| relevante en una rama, perdido al fusionar | ranks/pesos/candidate-k | fusión o budgets |
| relevantes presentes pero muy abajo | hard negatives/duplicados | reranking, mejores labels |
| métrica buena, respuesta mala | contexto y generador | no “mejorar retrieval” a ciegas |
| latencia alta | trace sparse/embed/dense/fusion | índice, parallelism, cache o budget |

La unidad de aprendizaje es el caso reproducible: query + qrels + candidatos de cada rama + filtros + ranks/scores + versiones + latencias. Toda corrección importante se convierte en una regresión.

### Gate del Módulo 2

- 10 videolecciones/transcripciones contrastadas.
- 2 ejemplos de código y 1 lectura inventariados; habilidades transformadas en contratos y pseudocódigo independiente de proveedor.
- 2 evaluaciones certificables contabilizadas, no abiertas ni reproducidas.
- TF-IDF, BM25, embeddings, metadata, hybrid/RRF, métricas, eficiencia y error analysis quedaron reconciliados con NLP y Producción.

## 18. Módulo 3 — Information Retrieval with Vector Databases

### Exactitud, escala y ANN

**[DLAI]** Exact k-nearest neighbors compara la query con todos los vectores y crece linealmente con el corpus. Approximate nearest neighbors reduce drásticamente el trabajo mediante índices precomputados, aceptando que puede omitir el vecino exacto. HNSW organiza un proximity graph por capas: navega saltos gruesos arriba y refina abajo.

**[IMPL] Precisiones necesarias:**

- la construcción HNSW real suele ser incremental; no asumir que siempre calcula todas las distancias por pares;
- “aproximadamente logarítmico” describe comportamiento esperado, no una garantía independiente de dimensión/distribución/parámetros;
- ANN recall se mide contra exact kNN sobre una muestra, separado de relevancia humana;
- los knobs de construcción y búsqueda cambian build time, RAM, disk, latency y recall; registrar valores, no sólo el nombre del índice;
- filtros selectivos pueden alterar rutas y recall; probar el patrón de filtros real.

```text
ANN gate = Recall_vs_exact@k + relevance@k + p99 + QPS + memory + build/update cost
```

No adoptar ANN por anticipación: un índice exacto, motor relational con vector extension o búsqueda brute-force puede ser suficiente para corpus pequeño. Elegir con benchmark del tamaño, churn, filtros y SLA reales.

#### Práctica pública Meta/Microsoft — el índice es una frontera de recursos

**[CODE]** Faiss, commit `7059eaf7da7eddda62e71367e684d4bdedd7f94f` (MIT), materializa una familia de índices exactos, IVF, quantization, graph y CPU/GPU. Su propia estructura de benchmark barre puntos de operación y separa tiempo de búsqueda, calidad, bytes por vector, entrenamiento y carga. **[CODE]** DiskANN, commit `860cf47bc11b6c2b938818773831a9374553d976` (MIT), expone proveedores de almacenamiento para memoria, disco o key-value store, además de updates y filtros. Juntos convierten “usar una vector DB” en decisiones medibles:

```text
vector/metric contract
→ exact oracle y qrels de tarea
→ storage tier + index family + compression
→ build/train/add/update/delete contract
→ search knobs + filters + rerank
→ recall/quality + latency/QPS + bytes/cost + freshness
```

Reglas de adopción:

- coseno exige normalización compatible cuando el índice opera por inner product; misma dimensión no implica mismo espacio;
- quantization reduce bytes/movimiento a cambio de error: medir `Recall_vs_exact@k` y relevancia, no sólo tamaño;
- IVF requiere datos de entrenamiento representativos y un `nprobe`/equivalente de búsqueda; HNSW/graph tiene otros costos de memoria, build y update;
- un índice que no cabe en RAM cambia el problema a cache, I/O, layout, prefetch/beam y tails; “SSD” no es un parámetro suficiente;
- filtros y ACL interactúan con candidate generation: aplicar un filtro después puede dejar menos de `k`; introducirlo durante traversal puede cambiar recall. Probar selectividades reales y autorización negativa;
- updates, deletes y compaction/rebuild deben preservar identidad, freshness y rollback; medir recall después de streams largos, no sólo tras build limpio;
- CPU/GPU requiere contabilizar residencia y transferencias. Un kernel rápido con queries/results movidos por request puede perder end-to-end;
- comparar branches/configuraciones con repeticiones A/A para estimar ruido antes de aceptar una diferencia A/B. El repositorio DiskANN publica tolerancias separadas para QPS y latencias media/p95, señal de que un número aislado no es evidencia suficiente.

**Gate reproducible:** fijar commit/library, hardware/NUMA/storage, dataset y split de train/base/query, dimensión/distribución, métrica/normalización, exact ground truth, threads/batch/concurrency, cold/warm cache, index parameters y lifecycle de updates. Reportar frontera Pareto y guardar índice/config/digests. Repetir con queries y filtros del producto; los benchmarks SIFT/Deep/Wikipedia validan mecánica, no relevancia local.

### Vector database como sistema de estado

**[DLAI]** El flujo común es crear colección/schema, cargar objetos, generar sparse/dense representations, construir índice ANN y ejecutar búsquedas vectoriales, BM25, híbridas y filtradas. El curso usa Weaviate como implementación, pero los conceptos son transferibles.

**Contrato de ingestión:**

```text
parse source → normalize → chunk → assign stable IDs/offsets
→ attach ACL/tenant/metadata/source version
→ embed with pinned model
→ upsert sparse+dense+payload
→ verify counts/checksums/sample queries
→ publish index version atomically
```

Registrar y reintentar errores por lote; no publicar un índice parcialmente cargado como sano. Diseñar deletes/tombstones, updates, deduplicación, backfill, re-embedding, snapshot/restore y rollback. Dual-write/read sólo con reconciliación medible. Un vector store no reemplaza transacciones/relaciones de la fuente de verdad.

La afirmación histórica de que bases relacionales rinden mal en vector search no es universal: capacidades cambian. Comparar candidatos actuales con documentación primaria y benchmarks propios; evitar acoplar el contrato de dominio a una API particular.

### Chunking como parte del modelo

**[DLAI]** Dividir documentos evita exceder el límite del embedder, representa temas con mayor precisión y reduce contexto enviado al LLM. Chunks grandes diluyen señal/consumen ventana; demasiado pequeños pierden contexto. Fixed-size con overlap y recursive splitting son baselines; el curso ofrece 500 caracteres con 50–100 de overlap como punto inicial ilustrativo.

**[IMPL]** No convertir ese número en regla. Tunar en tokens y por tipo documental. Preferir boundaries naturales:

- texto: títulos, secciones, párrafos, oraciones;
- código: módulo/clase/función/símbolo y dependencias;
- tablas: headers + filas/grupos preservando unidades;
- PDF/slides: jerarquía visual, página/slide y orden de lectura;
- conversaciones: turnos y ventanas con participantes/tiempo.

Cada chunk conserva `document_id`, `chunk_id`, offsets/página, jerarquía, metadata/ACL, checksum y vecinos. Overlap duplica almacenamiento y puede inflar ranking/citas; deduplicar por fuente y unir contexto contiguo después de recuperar.

Semantic chunking corta al detectar cambio de significado. LLM-based chunking/contextualization puede mejorar relevancia offline, pero añade costo, nondeterminismo, riesgo de alterar significado y complejidad de auditoría. Guardar texto original y generated context por separado, versionar prompt/model y comprobar que el contexto añadido no inventa hechos. Probar en muestra antes de reindexar el corpus.

### Query parsing sin perder intención

**[DLAI]** Las preguntas conversacionales no siempre son buenas búsquedas. Un LLM puede reescribirlas, aclarar ambigüedad, quitar ruido y añadir terminología/sinónimos. NER extrae personas, fechas u otras entidades para search/filtering. HyDE genera un documento hipotético, lo embebe y busca documentos parecidos.

**[IMPL] Guardrails:**

- conservar query original, rewrite, modelo/prompt y razón del cambio;
- no permitir que el rewriter decida permisos ni convierta conjeturas en filtros obligatorios;
- proteger entidades, números, negación, temporalidad e idioma con parity checks;
- tratar el hypothetical document como query expansion, nunca como evidencia;
- medir recall ganado contra latencia, costo y falsos conceptos introducidos;
- fallback directo si rewrite/NER/HyDE falla o agota deadline.

La reescritura básica es un candidato fuerte, no una obligación universal. Exact lookup puede empeorar si se altera un ID; conversaciones necesitan además resolver referencias usando historial autorizado.

### Bi-encoder, cross-encoder y ColBERT

| Arquitectura | Interacción query-document | Precompute | Ventaja | Costo/riesgo |
|---|---|---|---|---|
| bi-encoder | vectores independientes | documentos | máxima escala y ANN | interacción comprimida en un vector |
| cross-encoder | texto conjunto → score | no para query desconocida | ranking contextual fuerte | costo por cada par; rerank pequeño |
| ColBERT/late interaction | vectores por token + MaxSim | documentos | interacción fina con más escala | almacenamiento e índice mucho mayores |

**[DLAI]** Cross-encoder se usa normalmente después de recuperar un conjunto acotado. ColBERT suma, para cada token de query, su mejor similitud con tokens del documento. **[IMPL]** Un output en `[0,1]` no es probabilidad calibrada salvo que se demuestre; no usarlo como confidence de respuesta.

### Retrieve ancho, rerank estrecho

```text
corpus
→ hybrid ANN/BM25: candidate_k
→ ACL/dedup
→ cross-encoder/late-interaction/LLM reranker
→ final_k dentro del context budget
```

**[DLAI]** El curso muestra over-fetch de decenas y entrega final de pocos documentos; son heurísticas, no constantes. **[IMPL]** Tunar candidate-k/final-k juntos con batch size y deadline. Un reranker sólo puede reordenar lo recuperado: si el relevante no está en candidates, mejorar first stage. Un LLM reranker añade variabilidad/prompt injection y requiere output estructurado, repetibilidad y evaluación contra cross-encoder.

### Operación y error analysis

Instrumentar:

```text
parse/embed/search_sparse/search_dense/fuse/rerank/fetch_context
latency + candidates in/out + filters + versions + timeout/error
```

Slices mínimos: corpus size/churn, filter selectivity, query length/language, rare IDs, long documents, chunk position, duplicates, recent updates y tenant. Diagnóstico:

- relevante fuera de ANN pero presente en exact → parámetros/índice ANN;
- documento correcto, chunk incorrecto → boundaries/metadata/contextualization;
- relevante en candidates pero cae → reranker/truncation/dedup;
- rewrite cambia intención → parser/fallback;
- índice bueno pero stale → ingest/freshness/publish;
- p99 alto → embedding, search breadth, filters o reranker según trace.

### Gate del Módulo 3

- 9 videolecciones/transcripciones contrastadas.
- 2 ejemplos de código y 1 lectura inventariados; Weaviate se mantuvo como ejemplo, no dependencia normativa.
- 2 evaluaciones certificables contabilizadas, no abiertas ni reproducidas.
- ANN, ciclo de índice, chunking, query parsing, arquitecturas y reranking quedaron operacionalizados con métricas y fallos.

## 19. Módulo 4 — LLMs and Text Generation

### Delta frente al manual de Transformers

La arquitectura Transformer, attention, tokenización, pretraining, fine-tuning y serving se gobiernan desde `PYTORCH_TRANSFORMERS_GENERATIVE_MODELS.md`. Aquí importa su consecuencia para RAG:

- el modelo puede relacionar pregunta y evidencia dentro del contexto;
- la generación sigue siendo autoregresiva y no verifica verdad por diseño;
- input, output y reasoning tokens consumen ventana, tiempo y dinero;
- más contexto puede introducir ruido, conflicto y *lost in the middle*;
- prompt/model/chat template/tokenizer forman una unidad versionada.

**[IMPL]** KV cache evita recomputar partes de tokens previos durante decoding, pero no vuelve gratis la atención ni prompts largos. Medir prefill/TTFT y decoding/TPS por separado. La explicación pedagógica del curso no debe convertirse en una fórmula universal de complejidad para todas las arquitecturas/kernels.

### Decoding controlado

**[DLAI]** Greedy elige el token más probable; temperature reescala la distribución; top-k restringe cantidad de candidatos y top-p el conjunto de masa acumulada; repetition penalty/logit bias alteran tokens específicos o repetidos.

**[IMPL] Reglas:**

- factual QA: empezar con baja variabilidad y un output contract; creatividad sólo si la tarea la requiere;
- `temperature=0` no garantiza bitwise determinism entre proveedor, versión, batch o hardware;
- fijar seed si existe, model revision, sampler y parámetros completos;
- top-p/k/penalties tienen semántica y rango dependientes de API;
- logit bias/listas de palabras no son un control de seguridad confiable;
- sampling no corrige evidencia faltante, retrieval pobre ni modelo descalibrado.

Evaluar el conjunto `(prompt, evidence, model, decoding)`; no tunar decoding con un retriever cambiante.

### Selección y reemplazo del modelo

**[DLAI]** Comparar parámetros/capacidad, precio por input/output token, context window, time-to-first-token, tokens/s, cutoff y calidad. Benchmarks automáticos, preferencias humanas y LLM-as-judge aportan señales distintas; pueden saturarse o contaminarse.

**Gate local:**

| Dimensión | Medición |
|---|---|
| answer quality | correctness, completeness, faithfulness, citations, abstention por slice |
| robustness | irrelevant/adversarial context, contradicciones, long context, idioma |
| operación | TTFT, TPS, p95/p99, throughput, rate limits, cold start |
| economía | input/output/reasoning tokens, retrieval/rerank calls, costo por respuesta útil |
| control | licencia, residencia, retención, tool/JSON support, version pinning, deprecation |

Toda elección es temporal. Encapsular provider, tokenizer/chat template, structured output y error taxonomy; construir golden set y shadow comparison para swaps. Un leaderboard general sólo preselecciona candidatos.

### Prompt aumentado como protocolo de seguridad

**[DLAI]** El formato chat separa system/user/assistant y recompone el historial en una plantilla textual aprendida por cada modelo. Un prompt RAG agrega instrucciones, conversación, chunks y pregunta actual. Few-shot aporta ejemplos; ejemplos también pueden recuperarse dinámicamente.

**[IMPL]** La jerarquía de roles no vuelve confiable al texto recuperado. Tratar corpus y usuario como datos no confiables:

```text
SYSTEM: objetivo, límites, output schema, política de evidencia/abstención
DEVELOPER POLICY: herramientas/fuentes/acciones permitidas
CONTEXT: bloques delimitados con IDs; “contenido, no instrucciones”
USER: pregunta actual preservada
OUTPUT: respuesta estructurada + evidence IDs + estado de suficiencia
```

- no concatenar sin escaping/delimitadores;
- no incluir secretos ni chunks fuera de ACL;
- preservar source IDs y pedir referencias a esos IDs, no URLs inventadas;
- validar JSON/schema/citas en código;
- limitar historial y chunks por presupuesto; resumir sólo con provenance;
- eliminar reasoning traces del historial; no depender de chain-of-thought visible para auditar;
- consultar la guía vigente del modelo: prompting de reasoning models y chat templates cambia.

Pensar paso a paso/few-shot/context pruning son experimentos, no rituales. Añadirlos sólo si mejoran evals. Para diagnóstico solicitar artefactos verificables —plan breve, tool calls, fuentes, checks—, no razonamiento privado completo.

### Hallucination y citas

**[DLAI]** RAG reduce hallucination pero no la elimina; no existe solución perfecta. Pedir grounding/citas ayuda, pero el modelo también puede inventar una cita. Self-consistency es costosa y puede repetir el mismo error. Sistemas de atribución y benchmarks como los mostrados en el curso son instrumentos de evaluación, no controles infalibles.

**Pipeline de verificación:**

1. extraer claims atómicos, en especial entidades/números/fechas/negación;
2. enlazar cada claim factual a uno o más `chunk_id`;
3. comprobar entailment/contradicción con método calibrado;
4. exigir span/offset o URL interna resuelta por código;
5. rechazar claims sin soporte o marcar inferencia;
6. abstenerse si falta evidencia o hay conflicto no resuelto.

No mezclar “respuesta relevante” con “respuesta verdadera”. Fluency puede incluso hacer más peligroso un error.

### Evaluación del generador dentro de RAG

**[DLAI]** El generador debe responder la pregunta, usar lo relevante, ignorar ruido y citar. El curso presenta response relevancy, faithfulness y sensibilidad al contexto mediante evaluadores LLM, y complementa con feedback humano/experimentos.

**[IMPL]** RAGAS y otros frameworks cambian; fijar versión y auditar la definición concreta. LLM judge requiere rubric, muestras humanas, agreement, position/style/self-family bias y casos adversariales. No optimizar directamente contra un único juez.

```text
eval instance = question + authorized evidence + optional reference
              + retrieved candidates + final context + response
              + component/system scores + human label when sampled
```

Hacer component isolation:

- evaluar generador con evidence gold para no culparlo por retrieval;
- evaluar retriever sin generación;
- evaluar end-to-end con tráfico representativo;
- ejecutar ablaciones de prompt/model/decoding sobre artefactos congelados.

### Agentic RAG con límites explícitos

**[DLAI]** Patrones: secuencial; condicional/router; iterativo/evaluator loop; paralelo/orchestrator–workers–synthesizer. Modelos pequeños pueden especializarse en routing/evaluación y uno mayor generar la respuesta.

```text
router → retrieve? → evidence evaluator → (re-query | answer)
                                     max_steps/budget ┘
→ citation/verification → accept | abstain
```

**[IMPL]** Cada nodo tiene schema, timeout, permisos, costo, métricas y fallback. Las transiciones críticas usan código/policy cuando sea posible. Evaluar router confusion matrix, evidence-evaluator false accept/reject y loop termination. Paralelismo reduce wall time sólo si respeta rate limits y si synthesis no pierde evidencia. Agentes aumentan superficie de prompt injection y no deben expandir su propia autoridad.

### RAG, fine-tuning o ambos

| Necesidad | Opción inicial |
|---|---|
| conocimiento privado/reciente/citable | RAG |
| estilo, formato, comportamiento o tarea repetida | prompting → PEFT/SFT si no basta |
| router/reranker pequeño especializado | fine-tuning con dataset de tarea |
| conocimiento + comportamiento especializado | RAG + fine-tuning, evaluados por separado |
| operación exacta/estado vivo | API/SQL/tool; no memorizarlo en weights |

**[DLAI]** Fine-tuning adapta mejor cómo responde que inyecta hechos nuevos; puede empeorar otros dominios. **[IMPL]** No es una dicotomía ni una garantía: comparar contra prompting/RAG baseline, controlar forgetting, leakage, licencia, costo y rollback.

### Error analysis de generación con LLM

Categorías: evidence absent; retrieval miss; context ignored; distractor followed; unsupported claim; wrong citation; instruction injection; query misunderstood; history contamination; format/schema; sampling variance; model capability; judge defect; latency/cost. Registrar etapa causal y cambiar un componente por experimento.

### Gate del Módulo 4

- 11 videolecciones/transcripciones contrastadas.
- 2 ejemplos de código y 1 lectura inventariados; teoría Transformer remitida a su autoridad y aquí conservada sólo en su delta RAG.
- 2 evaluaciones certificables contabilizadas, no abiertas ni reproducidas.
- Decoding, selección de modelo, prompting, grounding, evaluación, agentic workflows y RAG/fine-tuning quedaron convertidos en gates implementables.

## 20. Módulo 5 — RAG Systems in Production

### Prototipo y producción son problemas distintos

**[DLAI]** Producción agrega tráfico/concurrencia, costo, latencia, prompts imprevisibles, datos desordenados y multimodales, privacidad/seguridad y consecuencias comerciales reales. El sistema debe anticipar, localizar y reparar fallos, y demostrar que cada cambio mejora el resultado.

**[IMPL]** Antes de producción exigir SLO y error budget por tarea/tenant:

```text
availability + p95/p99 latency + throughput
+ answer quality/abstention + security/privacy
+ freshness + cost per successful outcome
```

Un demo que responde preguntas seleccionadas no prueba coverage, aislamiento, degradación, idempotencia, backpressure, rollback ni capacidad bajo carga.

### Matriz de evaluación y observabilidad

**[DLAI]** Dos ejes organizan evals: alcance de componente frente a end-to-end, y evaluador por código, LLM o humano.

| Alcance | Código/determinista | LLM judge calibrado | Humano |
|---|---|---|---|
| ingest/index | counts, checksum, schema, freshness, duplicate rate | extracción/layout sample | revisión de fuentes críticas |
| retrieval | P/R/MRR/nDCG, ANN recall, latency | relevancia/contradicción | qrels y adjudicación |
| generation | schema, citas resolubles, tokens, latency | correctness/faithfulness/style | riesgo, utilidad y casos ambiguos |
| end-to-end | SLO, error/abstention/cost rate | rubric de respuesta completa | satisfacción y revisión experta |

**Trace mínimo:**

```text
request/correlation ID, tenant/policy version
→ original query → rewrite/router
→ sparse+dense candidates, filters, fusion/rerank
→ selected chunks/offsets/context
→ model/prompt/decoding/tools
→ answer/citations/verifier/abstention
→ latency, tokens, cost, errors per span
```

OpenTelemetry o equivalente para propagación; plataforma LLM para spans/evals; métricas de infraestructura para CPU/GPU/RAM/vector DB. La herramienta del curso es un ejemplo, no una dependencia obligatoria. No almacenar prompts/chunks/razonamientos sensibles por defecto: redacción, hashing, sampling, cifrado, controles de acceso, retención y borrado forman parte del diseño.

### Dataset nacido del tráfico y los errores

**[DLAI]** Guardar tráfico real permite reproducir prompts, inspeccionar cada componente, agrupar temas y medir cambios. Un caso presentado localizó un router que enviaba solicitudes de diagramas a generación de imágenes en vez de gráficos basados en código.

**[IMPL] Flywheel:**

1. muestrear fallos, baja confianza, abstenciones, thumbs down y casos de alto impacto;
2. desidentificar y respetar consentimiento/retención;
3. reconstruir trace y asignar categoría causal;
4. corregir label/qrel si el esperado era erróneo;
5. agregar el caso a dev o test de regresión según protocolo, evitando contaminación;
6. ejecutar componente, end-to-end, slices, seguridad y load tests;
7. shadow/canary con guardrails;
8. promover o rollback según gates predefinidos.

Feedback de usuario está sesgado por quién responde y por estilo; nunca usarlo solo. Clustering ayuda a descubrir slices, pero requiere revisar y nombrar clusters humanos antes de tomar decisiones.

### Presupuesto de costo

```text
cost/request = retrieval + rerank
             + Σ LLM(input_tokens·price_in + output_tokens·price_out)
             + storage/index/compute/observability amortizados
```

**[DLAI]** Palancas: modelo menor o cuantizado, menos input/output tokens, hardware dedicado a escala, tiering de RAM/disk/object storage y partición/multitenancy. **[IMPL]** Optimizar `costo por respuesta aceptada`, no costo bruto: un modelo barato que fuerza reintentos/revisión puede ser más caro.

Model routing necesita eval por ruta; autoscaling compara utilización, queue time y cold start; reserved/dedicated sólo gana con volumen y operación demostrados. Storage tiering debe preservar SLA/freshness y no mover datos por geografía/horario con suposiciones rígidas sobre usuarios.

### Presupuesto de latencia

```text
p99_total ≠ suma ciega de promedios
critical path = auth + parse/embed + search/fuse/rerank
              + context fetch + prefill/TTFT + decode + verify
```

**[DLAI]** Atacar primero llamadas Transformer/LLM, luego rewriter/reranker/router y finalmente storage; considerar modelos menores/cuantizados, routing, cache, binary vectors y sharding. **[IMPL]** Confirmar con trace: retrieval puede dominar bajo filtros difíciles, cold caches o red remota.

- paralelizar ramas independientes y cancelar trabajo ya innecesario;
- streaming mejora latencia percibida, no necesariamente completion time;
- limitar pasos agentivos y aplicar deadline propagation;
- batch mejora throughput pero puede empeorar tail latency;
- semantic cache key incluye tenant/ACL, locale, time/freshness, model/prompt/index version;
- cachear evidencia/planes suele ser más seguro que una respuesta final mutable;
- stale o cross-tenant cache hit es incidente, no optimización.

La frontera de order books permanece: RAG/LLM fuera del hot path atómico; sólo análisis asíncrono sobre snapshots/events ya capturados por infraestructura determinista.

### Cuantización medida, no asumida

**[DLAI]** Reducir precisión de weights o vectores ahorra memoria y puede acelerar. Se presentan int8/int4 para LLM, int8/binary vectors con posible rescoring full precision y embeddings Matryoshka que admiten prefijos dimensionales.

**[IMPL] Correcciones:**

- precision original, quantizer, calibration y hardware varían; no asumir siempre FP16 ni speedup automático;
- medir calidad por tarea/slice, no sólo benchmark general o Recall@k promedio;
- cuantizar corpus y query con transformación compatible;
- binary retrieval suele necesitar candidate expansion + full-precision/reranker;
- Matryoshka ordena capacidad útil mediante su objetivo de entrenamiento, no por una regla genérica de “varianza” aplicable a cualquier embedding;
- versionar quantization config con índice/modelo y conservar rollback.

Gate: memory/disk, build time, QPS, p99, ANN recall, relevance, answer quality y costo total antes/después.

### Threat model RAG

**[DLAI]** Autenticar, aislar información por tenant/rol, controlar proveedores externos y considerar on-prem; chunks pueden cifrarse, mientras índices/vector representations también son activos sensibles y pueden filtrar información.

**[IMPL]** Aislamiento físico por tenant puede ser fuerte, pero no es el único diseño seguro; row-level security/ACL bien implementado puede ser válido. Metadata de personalización no debe sustituir autorización. On-prem reduce una frontera de terceros, pero no garantiza seguridad.

| Amenaza | Controles mínimos |
|---|---|
| acceso cruzado/IDOR | authn, authz en ingest y query, tenant-scoped IDs/index/cache, pruebas negativas |
| prompt injection en documentos | contenido delimitado/no ejecutable, tool allowlist, policy fuera del corpus, output validation |
| poisoning/stale source | provenance, firma/checksum, review, trust tier, freshness/rollback |
| exfiltración por LLM/proveedor | data classification, minimización, DPA/retención, regional/on-prem cuando corresponda |
| fuga por logs/evals/cache | redacción, cifrado, least privilege, TTL, deletion audit |
| vector/model inversion | tratar embeddings como datos sensibles; cifrado/segmentación y acceso mínimo |
| abuso/DoS/costo | rate limit, quotas, max tokens/steps, timeouts, circuit breaker |
| dependencia/artefacto | pinning, SBOM/licencia, revisión de remote code, secrets separados |

Además: incident response, key rotation, backups probados y auditoría de deletes. Nunca confiar en el system prompt como barrera de confidencialidad.

### Multimodal RAG

**[DLAI]** Embeddings multimodales alinean texto/imágenes en un espacio; un vision-language model procesa tokens de texto e imagen. PDF/slides pueden recuperarse como páginas/regiones; técnicas tipo late interaction comparan tokens de query con múltiples patches, a cambio de muchos vectores.

**[IMPL] Pipeline verificable:**

```text
file → malware/type validation → layout/OCR/table/image extraction
→ hierarchical regions + page/bbox/source IDs
→ text and/or multimodal embeddings
→ retrieve/rerank at region and document level
→ send original crop + OCR + provenance to VLM
→ cite page/bbox; validate numbers/tables
```

Una cuadrícula de patches es una técnica emergente del curso, no default universal. Comparar contra OCR/layout-aware baselines. Medir cross-modal Recall@k, page/region localization, OCR/table accuracy, faithfulness y costo visual. Proteger contra instrucciones incrustadas en imágenes/PDF y archivos maliciosos.

### Runbook de producción

1. definir answer/evidence/security/SLO/cost contracts;
2. validar corpus, provenance, ACL y freshness antes de indexar;
3. publicar índice versionado tras tests de integridad y golden queries;
4. desplegar modelo/prompt/config como release indivisible;
5. ejecutar offline gates, load/security tests y shadow;
6. canary por tenant/traffic slice con rollback automático;
7. monitorear SLO + calidad + costo + seguridad por componente;
8. triage por trace y árbol causal, no por intuición;
9. convertir fallos reales en evals desidentificadas;
10. reindexar/reentrenar sólo cuando la intervención supera baseline y presupuesto.

### Gate del Módulo 5 y del curso

- 11 videolecciones/transcripciones contrastadas.
- 1 ejemplo de código y 4 lecturas inventariados; 2 evaluaciones certificables contabilizadas, no abiertas ni reproducidas.
- Producción, observabilidad, custom evals, costo, latencia, cuantización, seguridad y multimodalidad quedaron operacionalizados.
- **Gate del curso:** 5 módulos, 77 unidades, **49 videos contrastados**, 9 ejemplos de código y 9 lecturas inventariados; 10 evaluaciones certificables excluidas.

## 21. Auditoría transversal final NLP + RAG

### Invariantes sin contradicción

| Decisión | Regla consolidada |
|---|---|
| sparse vs dense | BM25 es baseline y señal complementaria; dense se justifica por recall semántico local |
| vector DB vs base tradicional | elegir por escala, filtros, churn, SLA y benchmark; vector store no es fuente transaccional |
| chunking | forma parte del modelo; preservar estructura, offsets, ACL y evaluar retrieval + respuesta |
| retrieval vs generation | evaluar aislados y end-to-end; reparar la etapa causal |
| prompt vs seguridad | prompt orienta conducta; autorización/validación se implementan fuera del LLM |
| RAG vs fine-tuning | conocimiento/citas → RAG; conducta/tarea → tuning; combinar si evals lo justifican |
| long context vs retrieval | context window grande no elimina selección, permisos, freshness, costo ni atribución |
| agentic vs pipeline fijo | usar agencia donde decisiones iterativas agregan calidad; mantener budgets, schemas y gates deterministas |
| calidad vs costo/latencia | Pareto por SLO y costo de error; nunca optimizar un promedio único |
| observabilidad vs privacidad | capturar causalidad con IDs/versiones y minimizar/redactar contenido sensible |

### Contrato que Codex debe producir al construir RAG

```text
1. Problem/answer/evidence/abstention contract
2. Threat model + identity/ACL/tenant/data classification
3. Corpus and ingestion/index lifecycle
4. Baselines: no-RAG, BM25, dense, hybrid
5. Chunk/query/rank/context policies versioned
6. Generator prompt/model/decoding/output schema
7. Component + end-to-end eval datasets and slices
8. Trace/metrics/SLO/cost dashboards
9. Load, failure, injection, privacy and rollback tests
10. Shadow/canary release and error-to-eval feedback loop
```

Si faltan answer contract, qrels/evidence, ACL o rollback, Codex debe detener el despliegue y declarar el hueco; puede construir un prototipo aislado, no llamarlo production-ready.

### Resultado final

- NLP: 4 cursos, 14 semanas, 412 unidades y 176 videos, completo.
- RAG: 5 módulos, 77 unidades y 49 videos, completo.
- Total de este documento: **225 videolecciones contrastadas**, con prácticas visibles traducidas a contratos y 30 evaluaciones certificables excluidas.
- Procedencia `[NG]`, `[DLAI]`, `[IMPL]` y `[EXT]` permanece separada; APIs o benchmarks temporales no se elevan a principios eternos.
- El manual ya permite especificar, implementar, probar, desplegar, observar y mejorar un sistema NLP/RAG aprendiendo de sus errores; no sustituye formación externa de bases de datos internas, redes, sistemas, seguridad o low latency.

### Fuentes industriales verificadas de retrieval

- Faiss/Meta, snapshot auditado: <https://github.com/facebookresearch/faiss/tree/7059eaf7da7eddda62e71367e684d4bdedd7f94f>
- Faiss benchmark suite: <https://github.com/facebookresearch/faiss/blob/7059eaf7da7eddda62e71367e684d4bdedd7f94f/benchs/README.md>
- DiskANN/Microsoft, snapshot auditado: <https://github.com/microsoft/DiskANN/tree/860cf47bc11b6c2b938818773831a9374553d976>
- DiskANN A/A benchmark policy: <https://github.com/microsoft/DiskANN/blob/860cf47bc11b6c2b938818773831a9374553d976/.github/docs/disk-benchmarks-aa.md>

## Estado de incorporación de RAG

- [x] Curso oficial, instructor, duración, módulos y 77 unidades verificados.
- [x] Módulo 1 — RAG Overview.
- [x] Módulo 2 — Information Retrieval and Search Foundations.
- [x] Módulo 3 — Information Retrieval with Vector Databases.
- [x] Módulo 4 — LLMs and Text Generation.
- [x] Módulo 5 — RAG Systems in Production.
- [x] Auditoría transversal RAG y reconciliación con NLP, PyTorch/Transformers, Producción y Agentes.
