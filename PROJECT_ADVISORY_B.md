# Asesoramiento de proyecto — ronda B

Catalog-Version: 1.0.0
Catalog-SHA256: d5c336bf2c493e334ecc253aac5ed3d9e7c90d0ebd09fb78120da0c28bff6015
Topic-ID: B-ACTORS-JOURNEYS

Este prompt fue generado desde el catálogo ejecutable. El agente debe explicarlo en lenguaje simple, presentar un grupo manejable de preguntas y registrar la respuesta en los artefactos de readiness. No debe completar respuestas por inferencia ni pedir secretos en chat.

## Personas, permisos y recorridos completos

Describe qué puede hacer cada visitante, cliente, empleado, franquicia y proveedor de principio a fin. Licencia laboral o ausencia significa indisponibilidad de una persona; cancelación revierte una reserva; no-show registra que alguien no asistió; entrega completa incluye preparación, identidad, serie física, aceptación y evidencia.

## Por qué se necesita

Los roles y estados evitan exposición entre organizaciones, doble reserva, acciones sin autoridad y experiencias donde el usuario queda atrapado a mitad del recorrido.

## Preguntas que el agente debe ayudarte a responder

- ¿Qué actores existen y qué organización, objeto y acción puede operar cada uno?
- ¿Cuáles son los journeys P1, incluidos error, cancelación, no-show, reversa y recuperación?
- ¿Cómo se administran horarios, ausencias/licencias laborales, turnos, recursos y entrega final?

## Información que puedes proporcionar

- actores
- matriz de permisos
- journeys Given/When/Then
- reglas de cancelación y no-show
- reglas de ausencia
- evidencia de entrega

## Ejemplo ilustrativo

EJEMPLO, NO RESPUESTA ASUMIDA: el visitante consulta modelos sin login; el cliente reserva test-drive; la franquicia confirma recurso; una ausencia bloquea ese recurso; el cliente cancela antes del límite; no-show requiere actor y motivo; la entrega compara cliente, organización y número de serie antes de aceptar.

## Cómo se comprobará

- journeys positivos y negativos
- pruebas de autorización y aislamiento
- transiciones e identidad auditables

## Cuándo puede responderse NO_APLICA

Sólo NO_APLICA si el actor o journey no existe en el alcance aprobado; no se permite omitir cancelación o recuperación de una operación que sí existe.

## Autoridades oficiales

- Microsoft — Employee absence is explicit operational data with causes, periods and units; it is distinct from software licensing. — https://learn.microsoft.com/en-us/dynamics365/business-central/hr-how-manage-absence
- Google — Use curated, practical and adaptable launch questions, shared infrastructure, explicit dependencies, failure modes and staged rollout. — https://sre.google/sre-book/reliable-product-launches/
- NASA — Separate requirements, verification, validation, integration and stakeholder acceptance evidence. — https://www.nasa.gov/reference/system-engineering-handbook-appendix/

## Regla de avance

La ronda sólo se cierra con respuesta o evidencia real. Si falta una decisión o acceso material, el agente explica exactamente qué falta, conserva BLOCKED/AWAITING_USER y continúa con trabajo read-only que no dependa de inventarlo.
