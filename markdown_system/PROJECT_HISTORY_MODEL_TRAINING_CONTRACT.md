# Historial y entrenamiento de modelos en proyectos consumidores

Estado: contrato de alcance y admisión, AUTHORED. No es un pack ejecutable ni
evidencia de un importador o entrenamiento terminado. Clarificación del usuario
2026-09-09; sustituye la interpretación del requisito histórico V288 sin borrar
su evidencia. Owners existentes: T2805/T2807, con T2803/T2806; no crea otro roadmap.

## 1. Propósito y separación de la biblioteca

El requisito es preparar una capacidad para entrenar o adaptar un modelo del
proyecto consumidor a partir de conversaciones humanas pertinentes. No debe
sustituirse silenciosamente por RAG, memoria de chats o evaluaciones: son usos
distintos. La técnica concreta —incluido fine-tuning como candidato—, el modelo,
el objetivo medible y el presupuesto se deciden en ese proyecto. No se presupone
entrenamiento desde cero ni se elige proveedor por reputación.

La biblioteca distribuye instrucciones, contratos, código admitido y pruebas
reutilizables. No recibe conversaciones de clientes, datasets privados, pesos o
adapters entrenados con ellos. Leer estos Markdown no entrena al agente lector
ni activa entrenamiento, recolección o envío de datos. Ningún dato de un negocio
se convierte automáticamente en material de otro proyecto o en aprendizaje global.

| Situación | Comportamiento requerido |
|---|---|
| Proyecto NEW, sin historial | Inicia sin importación ni entrenamiento histórico. La ausencia de chats previos no bloquea el bootstrap ni exige inventarlos. La selección de un modelo base y su uso tienen sus propios gates. |
| Proyecto EXISTING con historial | Inventariar la fuente y su exportación/acceso oficial dentro de ese proyecto. Tener conversaciones no demuestra calidad, cobertura ni permiso de usarlas para entrenamiento. |
| NEW que acumula conversaciones después | Sólo genera un candidato de dataset si se decide activar esa finalidad y sus controles; no aprende continuamente de cada mensaje por defecto. |
| Mantenimiento de esta biblioteca | Puede investigar fuentes y preparar/verificar código y fixtures admitidos sin pedir el historial de una empresa hipotética. La validación de calidad con chats reales pertenece al proyecto que los utiliza. |

## 2. Dónde se usan los datos y dónde vive cada resultado

```text
Fuente del proyecto EXISTING
  → importación aislada sin efectos
  → originales privados con identidad, hash, permisos y retención
  → derivados revisados y versionados para la finalidad de entrenamiento
  → particiones de entrenamiento / validación / evaluación independiente
  → job de entrenamiento explícito del modelo seleccionado
  → modelo o adapter candidato + registro de dataset/código/configuración
  → evaluación contra baseline + revisión de seguridad, costo y calidad
  → promoción explícita al runtime del mismo proyecto + rollback
```

Los originales, derivados, particiones, checkpoints y modelos resultantes viven
en almacenamiento privado elegido y gobernado por el proyecto consumidor,
fuera del árbol distribuible de la biblioteca. No hay un directorio de clientes
compartido por defecto. El lugar físico —equipo propio o servicio seleccionado—,
responsable, acceso, cifrado, retención, borrado y egress deben quedar registrados
antes de mover datos. Nada se sube a un proveedor de entrenamiento por inferencia.

Los originales se conservan byte a byte dentro de la política autorizada; corregir,
redactar, desidentificar, etiquetar y deduplicar crea derivados con lineage.
No se sobrescribe el original ni se interpreta raw como retención ilimitada.
Un identificador hash no acredita anonimato. Si la fuente contiene adjuntos o
exports documentales, aplicar también OFFICIAL_DOCUMENT_INTELLIGENCE_PROFILE.md,
OFFICIAL_DOCUMENT_SDK_ARTIFACT_LOCK.md y OFFICIAL_DOCUMENT_FIXTURE_CATALOG.md,
con evaluación por formato/campo: un fixture público no valida el corpus privado.

## 3. Qué significa aprender y qué no

- **Entrenamiento/adaptación:** un proceso explícito modifica pesos o adapters y
  produce un candidato versionado. Guardar chats no demuestra que ocurrió.
- **RAG/conocimiento recuperable:** suministra contexto en inferencia; no cumple
  por sí solo el requisito de entrenamiento. Su uso requiere decisión separada.
- **Memoria de conversación:** contexto limitado al contacto/hilo/negocio y sus
  permisos. No implica uso como ejemplo de entrenamiento.
- **Evaluación:** mide candidatos en casos reservados; esos casos no se usan
  también para entrenar ni para proporcionar las respuestas esperadas al modelo.

Las conversaciones humanas se revisan: pueden contener errores, secretos, datos
de terceros, instrucciones maliciosas, respuestas obsoletas o decisiones no
aprobadas. El objetivo del dataset, labels, criterios de exclusión, particiones
y métricas se fijan antes del job. Evitar contaminación por conversación,
contacto, incidentes duplicados y tiempo según el dominio; registrar qué se probó.
No entrenar con credenciales, ni aceptar todo mensaje histórico como instrucción
del sistema. Precios, stock, permisos y estado transaccional vigentes conservan
sus owners; el modelo no sustituye las reglas ni inventa efectos del negocio.

## 4. Código admitido y puertas de ejecución

Antes de implementar, reutilizar un pack compatible admitido para el claim
exacto. Si falta, ejecutar CAPABILITY_GAP_RESOLUTION_GATE: investigación de
fuentes oficiales, identidad/revisión fijada, licencia por componente, SCA,
contratos, reconstrucción y pruebas G0–G8. RESEARCH_INCOMPLETE obliga a continuar
la investigación; no permite escribir un motor improvisado y llamarlo élite.
Un sample, tutorial o README de una empresa no admite una pipeline completa.
El glue local imprescindible conserva procedencia AUTHORED/ADAPTED, cambios y
pruebas explícitos; nunca se atribuye a la empresa del SDK ni reemplaza la
selección de una implementación admitida.

Cada proyecto registra antes de ejecutar: propósito del modelo; fuente y permiso
de entrenamiento; alcance por tenant; dataset/splits y hashes; fuente y licencia
del código/modelo base; técnica y configuración; lugar de ejecución/egress;
presupuesto finito; baseline/umbrales; responsable de revisión; políticas de
retención/borrado; versión candidata y rollback. Las decisiones ausentes bloquean
el job afectado, no todo el desarrollo independiente del proyecto.

La importación nunca llama Handle ni reencola históricos como mensajes live;
no envía respuestas ni invoca herramientas de ventas, cobro o reservas. El job
de entrenamiento es offline y separado del runtime. La evaluación comprueba
calidad por caso/slice, fugas entre tenants y de información sensible, abuso y
poisoning pertinentes, latencia/costo y regresiones contra baseline. Un job
exitoso o una métrica promedio mayor no autoriza desplegar el candidato.
La promoción requiere evidencia, aceptación del responsable y rollback probado.

## 5. Aceptación de biblioteca frente a aceptación de proyecto

La biblioteca debe demostrar mediante código admitido y fixtures adecuados:
NEW sin corpus no dispara importación/entrenamiento; EXISTING requiere un perfil
completo antes de tocar datos; aislamiento, lineage, train/eval separados,
rechazo de configuración incompleta, cero efectos de backfill y promoción cerrada.
V373 demuestra esos controles reutilizables en HISTORY-MODEL-TRAINING-PIPELINE0.1.0,
con27pruebas de política, entrenamiento SFT real con fixtures, evaluación
independiente,13negativos, timeout real, auditoría inmutable y rollback local
probado. Ver HISTORY_MODEL_TRAINING_PACK_PLAN.md y
reconstruction_evidence/HISTORY_MODEL_TRAINING_CONTROL_V373.md desde la raíz.
El perfil CPU/Windows es opt-in y conserva sus condiciones; este PASS de biblioteca
no admite otros modelos/targets ni hereda autoridad, seguridad o calidad del consumer.

El proyecto consumidor debe demostrar además contrato del canal y cobertura del
historial real, autorización de uso, calidad del dataset, ejecución reproducible,
resultado de entrenamiento, evaluación independiente, costo y aceptación del
modelo en su target. Esta evidencia no se hereda de fixtures ni de la biblioteca.
No pedir datos reales del proyecto consumidor como requisito para terminar el
contrato o la investigación de biblioteca; tampoco prometer calidad sin probarlos.

Autoridades internas: ML_PRODUCTION_LLMOPS_EVALUATION.md (provenance, lineage,
datasets/evaluación/promoción), AI_SECURITY_GOVERNANCE_PRIVACY.md (entrenamiento
independiente y autorizado, aislamiento y leakage), PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD.md
y ENGINEERING_EXECUTION_PLAYBOOK.md. Este contrato organiza sus límites; no
impone un modelo, SDK, proveedor o técnica al consumer. El pack SFT opcional
es una implementación compatible demostrada; otra selección exige su propia
admisión y evidencia conforme a este mismo contrato.
