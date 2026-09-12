# Asesoramiento de proyecto — ronda D

Catalog-Version: 1.0.0
Catalog-SHA256: d5c336bf2c493e334ecc253aac5ed3d9e7c90d0ebd09fb78120da0c28bff6015
Topic-ID: D-DOCUMENT-EVIDENCE

Este prompt fue generado desde el catálogo ejecutable. El agente debe explicarlo en lenguaje simple, presentar un grupo manejable de preguntas y registrar la respuesta en los artefactos de readiness. No debe completar respuestas por inferencia ni pedir secretos en chat.

## Archivos, extracción y almacenamiento preciso

Cada clase y variante identifica un layout real. Corpus es el conjunto representativo de archivos; ground truth son valores revisados; schema define tipos y campos; mapping transforma esos campos a los nombres exactos del sistema consumidor. Ningún proveedor convierte un archivo desconocido en dato correcto sin estas evidencias.

## Por qué se necesita

Persistir una extracción incorrecta contamina compras, stock, pagos y fiscalidad. La precisión debe demostrarse por clase, variante y campo antes de automatizar.

## Preguntas que el agente debe ayudarte a responder

- ¿Qué canales, clases, variantes, emisores, idiomas, layouts, volúmenes y límites de archivo existen?
- ¿Dónde están originales, checksums, corpus autorizado, ground truth, schemas y reglas cruzadas por clase?
- ¿Qué errores son tolerables por campo y quién revisa, aprueba, rechaza o corrige ambigüedades?
- Si habrá persistencia automática, ¿cuál es el modo y el mapping exacto al consumer con hashes y evidencia?

## Información que puedes proporcionar

- inventario clase/variante
- archivos autorizados
- ground truth
- schemas
- mapping
- umbrales por campo
- owner de aprobación
- retención y seguridad

## Ejemplo ilustrativo

EJEMPLO, NO RESPUESTA ASUMIDA: invoice/supplier-a-v1 y packing-list/factory-b-es-v2; 200 documentos anonimizados por variante; campos y unidades en schemas versionados; número, moneda y total requieren exactitud aprobada; discrepancias van a revisión; persistencia sólo tras evaluación y mapping hash-locked.

## Cómo se comprobará

- bytes y SHA-256 del original
- evaluación separada por clase/variante/campo
- schema y mapping reproducibles
- revisión humana donde la política lo exige

## Cuándo puede responderse NO_APLICA

Persistencia automática puede ser NO_APLICA y conservar extracción para revisión. Corpus o schema no son NO_APLICA si el sistema almacenará datos extraídos de esa clase.

## Autoridades oficiales

- Google — Document processing capabilities are processor- and document-type-specific; project data, schemas and evaluation remain project inputs. — https://cloud.google.com/document-ai/docs/processors-list
- Amazon Web Services — Document services have explicit format, size, page, operation and service limits that must be selected and tested. — https://docs.aws.amazon.com/textract/latest/dg/limits-document.html
- NASA — Separate requirements, verification, validation, integration and stakeholder acceptance evidence. — https://www.nasa.gov/reference/system-engineering-handbook-appendix/

## Regla de avance

La ronda sólo se cierra con respuesta o evidencia real. Si falta una decisión o acceso material, el agente explica exactamente qué falta, conserva BLOCKED/AWAITING_USER y continúa con trabajo read-only que no dependa de inventarlo.
