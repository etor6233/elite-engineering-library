# Official Field Evidence Code Reaudit — 2026-08-26 V1

## 1. Pregunta y criterio

Se buscó código público oficial de compañías líderes que pudiera comparar extracción documental contra ground truth, detectar campos faltantes o inventados, resolver listas uno-a-uno y sostener una decisión de almacenamiento sin presentar glue local como código de la empresa. El corte exigió identidad inmutable, licencia reusable, ejecución reproducible y casos adversariales; una descripción comercial o una suite verde no bastaron.

## 2. Microsoft Content Processing no supera el corte

Fuente fijada: `microsoft/content-processing-solution-accelerator` 2.1.2, commit `b47cec48475cfedd7109debf2f407ef0c5530311`; el `main` firmado `659eaa1f503dd08b1e1aea1c72eab11c7c191d00` del 25 de agosto se volvió a inspeccionar y conserva la misma frontera material.

- `confidence.py`, 7.775 bytes, SHA-256 `219cf029581e6a702edccfd40186a0b1206efd6447dfce628d7ca7c0e4fbd7f0`: `get_confidence_values` no incorpora confianza cero al promedio; `merge_confidence_values` conserva el valor de la primera extracción cuando la segunda difiere y sólo reduce la confianza.
- `gap_executor.py`, 8.334 bytes, SHA-256 `a621249c9ecf69caffc63fa55bed4d8c290c25d94601ec90d26990c9ad6251d0`: entrega reglas a un agente LLM y persiste su respuesta; no es un adjudicador determinista universal.
- `fnol_gap_rules.dsl.yaml`, 8.003 bytes, SHA-256 `515ca131161ef9a4478d3bb08c7d95aed99cdfd5156cbddaee6fa5064a7b53ee`: las reglas públicas son un ejemplo de seguros/FNOL, no reglas oficiales para facturas, proformas, packing lists, órdenes, aduana o tesorería del proyecto.

Decisión: las piezas quedan como `REJECTED_COMPONENT` para admisión/merge de campos. No se copiaron ni se corrigieron localmente. El acelerador general permanece condicionado por sus otras pruebas y por servicios Azure.

## 3. AWS OCR Evaluation Workbench no supera el corte de storage

Fuente fijada: `aws-samples/ocr-with-aws-ai-services`, commit firmado `5341910c3aee2da872f2966caed0fe0246b8a55c`, archive 4.127.057 bytes/SHA-256 `0e419c085a1957a37de791050a72e004b3d54b91f392ccf0311a4a5be1d2ac8f`, MIT-0. El archivo oficial `shared/evaluator.py` mide 14.137 bytes y tiene SHA-256 `412815b5be3d7f5819552f7aeec6be085240511e338d8fd274fdc0e868f0f4e3`.

Casos ejecutados sobre ese archivo exacto, cargado sin modificar:

1. expected `{"a":"x"}` contra actual `{"a":"x","b":"hallucinated"}` produjo 100%; el recorrido sólo parte del ground truth.
2. dos objetos esperados idénticos contra un único objeto extraído produjeron 100%; cada esperado puede volver a elegir el mismo mejor candidato.

Decisión: conserva valor como workbench condicionado, pero es `REJECTED_COMPONENT` como puerta de almacenamiento. El fallo inicial de import por un `gradio` ausente y la repetición aislada están registrados como `LIB-FAIL-137`; los resultados lógicos son `UP-FAIL-051`.

## 4. AWS Labs Stickler 0.6.0: candidato oficial superior, todavía condicionado

Autoridad y supply chain:

- repositorio: `awslabs/stickler`;
- release/tag: `v0.6.0`;
- commit firmado/verificado: `174ca9d3476c1ea2d0a36628c73d595604e0398d`;
- archive GitHub: 1.356.931 bytes, SHA-256 `f24ae4650b3ffc46d23e5b447a0e72a8f75a20e98df1b4ee0ca31490759707a3`;
- licencia Apache-2.0: 10.142 bytes, SHA-256 `09e8a9bcec8067104652c168685ab0931e7868f9c8284b66f5ae6edae5f1130b`;
- NOTICE: 67 bytes, SHA-256 `d4290ed64c2edd0fce1d84e3f9dfb2881240fe534def76b8cd29ed6af683e287`;
- `uv.lock`: 677.578 bytes, SHA-256 `325f27921c235b3a3620d3523e615fa3f63525a09deb0ed2fad96da901884dd6`;
- wheel PyPI atestado: 257.072 bytes, SHA-256 `d906c7d3f6538ac6e8582d1b71e4f74388a26073b2f1606956e12812efa9e4ae`;
- sdist PyPI atestado: 219.729 bytes, SHA-256 `405291e374604811dcac60fbd701afd2026c3b7e1576feb37605572699bf3e98`;
- la attestation PyPI enlaza ambos artefactos al commit AWS Labs anterior mediante Trusted Publishing;
- export runtime frozen: 25 dependencias más `stickler-eval`; consulta OSV batch actual: 26 paquetes, 0 findings conocidos.

Reproducción exacta: `uv 0.11.30` oficial verificado por SHA-256, Python 3.12.13 descargado por uv y `uv sync --frozen` sobre el checkout corto. Resultado Windows: 1.543 casos, 1.535 PASS, 2 skip y 6 fallos exclusivos del uso upstream de `signal.SIGALRM`, que Windows no implementa. No se editaron tests ni source.

Fortalezas verificadas del código oficial:

- `ExactComparator` y thresholds explícitos;
- matching húngaro uno-a-uno para listas;
- campos extras de raíz y objetos anidados directos contabilizados como `FA/FP`;
- filas faltantes contabilizadas como `FN`;
- valores críticos distintos contabilizados como `FD/FP`;
- `validate_instance_against_schema` usa `jsonschema.Draft7Validator` y bloquea propiedades adicionales cuando el schema las cierra.

Defecto adversarial reproducido:

- `comparison_engine.py`, 25.501 bytes, SHA-256 `f5fdd28ca6ccb14ffa83a43086350985f7cb26516db5094245e194b2321e59fc`, no cuenta un campo extra dentro de un objeto perteneciente a `List[StructuredModel]`; el probe dejó `FA=0`, `FP=0`, score/F1/accuracy 1.0.
- el mismo caso fue rechazado correctamente por la función oficial `validate_instance_against_schema` cuando cada objeto del JSON Schema declaró `additionalProperties: false`.
- `all_fields_matched` permaneció `true` aun con una falsa alarma de raíz y también con una fila faltante; no es un indicador de admisión.
- el `main` firmado `349aa8403f9bb80f59bd7ba86596229b3fbcf4fb`, candidato 0.7.0 del 25 de agosto, conserva la misma condición defectuosa; no reemplaza el release.

Decisión: `awslabs-stickler-eval-0.6.0` entra al source lock como `PINNED_CANDIDATE_CONDITIONED`, no como storage gate. Una composición segura debe validar primero ambos documentos con un JSON Schema cerrado recursivamente y luego exigir simultáneamente `FA=FD=FN=FP=0`; score, F1 o `all_fields_matched` por sí solos están prohibidos. El schema y el ground truth siguen siendo inputs aprobados del proyecto, no datos que la biblioteca pueda inventar.

## 5. Resultado y próximos gates

El lock canónico sube de 68 a 69 fuentes. Los ocho perfiles pasan contra el lock; AWS selecciona 9 y document intelligence 30. Las cinco negativas de lock/acquisition, cinco negativas de perfiles y una positiva offline siguen verdes.

Esto agrega código real de AWS Labs disponible para adquisición inmediata, con licencia, commit, artefactos y lock exactos. No autoriza afirmar extracción perfecta en documentos nuevos: primero hacen falta schemas cerrados por clase, corpus representativo, ground truth aprobado, evaluación completa, política de ambigüedad/revisión y evidencia del proveedor. Los defectos quedan retenidos como `UP-FAIL-050` a `UP-FAIL-053`; los fallos de auditoría y sus regresiones como `LIB-FAIL-137` a `LIB-FAIL-144`.

Fuentes oficiales gobernantes: <https://github.com/awslabs/stickler/tree/174ca9d3476c1ea2d0a36628c73d595604e0398d>, <https://pypi.org/project/stickler-eval/0.6.0/>, <https://github.com/aws-samples/ocr-with-aws-ai-services/tree/5341910c3aee2da872f2966caed0fe0246b8a55c> y <https://github.com/microsoft/content-processing-solution-accelerator/tree/659eaa1f503dd08b1e1aea1c72eab11c7c191d00>.
