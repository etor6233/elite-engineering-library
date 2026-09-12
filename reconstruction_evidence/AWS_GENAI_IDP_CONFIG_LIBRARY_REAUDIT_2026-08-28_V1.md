# AWS GenAI IDP configuration library re-audit — V1

Fecha: 2026-08-28  
Clasificación: `OFFICIAL_SOURCE_VERIFIED / PINNED_CANDIDATE_CONDITIONED / NOT_COMPLETE_FOR_LOGISTICS`

## Identidad

- repositorio oficial: `aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws`;
- release: `v0.6.5`;
- commit: `1b5fd74454e593de233342a02ee911af8ee38359`;
- licencia: MIT-0 + NOTICE;
- ZIP: 55.014.587 bytes, SHA-256 `62c84cdf5bed3a2029285c9a398d9532863953333dd09b115150af1cb7e0f57e`;
- el source ya pertenece al lock global; esta revisión no altera ni reemplaza ese lock.

## Ejecución exacta

Se extrajo el archive verificado y se ejecutó el test oficial `config_library/test_config_library.py` sin modificar source. El entorno aislado usó CPython 3.12 y las siete distribuciones exactas fijadas por `lib/idp_common_pkg/uv.lock`: PyYAML 6.0.2, pytest 9.0.3, colorama 0.4.6, iniconfig 2.3.0, packaging 26.0, pluggy 1.6.0 y Pygments 2.20.0. Cada wheel descargado tuvo un SHA-256 presente literalmente en el lock y la instalación offline terminó `pip check` limpio.

Resultado oficial: **114 PASS en 2,09 s**. La suite valida parse YAML, estructura de configs, tipos de campos y cobertura regex de policy classes. No llama AWS ni evalúa exactitud de extracción sobre el corpus empresarial del propietario.

## Cobertura observada

Los trece presets `config_library/unified/*/config.yaml` contienen 50 IDs de clase únicos. Incluyen, entre otros:

- invoice con 15 propiedades en RVL-CDIP;
- bank statement, W2, payslip, bank check, driver's license y homeowners insurance;
- delivery note dentro del benchmark OCR;
- email, form, letter, questionnaire, specification, budget y handwritten;
- diez clases healthcare y cinco clases prior-authorization con siete policy classes.

El release conserva configuraciones Level 1/2 o tests con muestras acotadas. `rvl-cdip` y `docsplit` declaran Level 2; el propio README advierte variación ante tipos especializados o layouts no representados. Ningún resultado se eleva a validación productiva universal.

## Brecha determinante

La búsqueda exacta sobre todos los YAML del release no encontró clases/schemas para:

- packing list;
- proforma invoice;
- bill of lading;
- commercial invoice.

`purchase order` aparece sólo como un atributo/referencia dentro de configs de invoice, no como una clase logística independiente. Por tanto, el catálogo no se materializa como cobertura completa ni habilita persistencia automática para esos documentos. AWS BDA permite crear blueprints custom y ofrece catálogo en el servicio, pero un blueprint generado por prompt sigue requiriendo schema aprobado, corpus, ground truth, evaluación, revisión y acceso real antes de convertirse en autoridad del sistema.

## Decisión

No se copió ni reetiquetó código en un nuevo implementation pack durante esta revisión. El source permanece `PINNED_CANDIDATE_CONDITIONED`: es código oficial real y probado en su scope, pero incorporarlo como “todos los documentos empresariales” sería falso. La condición upstream 159 obliga al agente a pedir las clases/campos/corpus faltantes o a encontrar una revisión oficial con esos artefactos y repetir todos los gates.

## Recuperación local

Las lecciones 862–865 conservan cuatro fallos de harness: globs literales Windows, stderr oculto, venv sin PyYAML y path de manifest supuesto. Ninguno se atribuyó a AWS. Las repeticiones usaron listas/path reales, errores visibles y dependencias exactas del lock.
