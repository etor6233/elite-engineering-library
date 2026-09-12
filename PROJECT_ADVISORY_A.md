# Asesoramiento de proyecto — ronda A

Catalog-Version: 1.0.0
Catalog-SHA256: d5c336bf2c493e334ecc253aac5ed3d9e7c90d0ebd09fb78120da0c28bff6015
Topic-ID: A-OUTCOME-AUTHORITY

Este prompt fue generado desde el catálogo ejecutable. El agente debe explicarlo en lenguaje simple, presentar un grupo manejable de preguntas y registrar la respuesta en los artefactos de readiness. No debe completar respuestas por inferencia ni pedir secretos en chat.

## Resultado, responsables y límites reales

Define qué negocio se construye, quién puede decidir, dónde operará, qué fecha y presupuesto existen y qué queda fuera. El owner es la persona que puede aceptar o rechazar una decisión; no es un nombre decorativo.

## Por qué se necesita

Sin resultado medible, jurisdicción y autoridad, el agente no puede distinguir una necesidad real de una suposición ni elegir prioridades, proveedores o evidencia de aceptación.

## Preguntas que el agente debe ayudarte a responder

- ¿Cuál es el resultado comercial u operativo medible del primer lanzamiento?
- ¿Quién decide producto, operación, seguridad y gasto, y qué decisiones requieren tu aprobación?
- ¿Qué países, idiomas, monedas, fecha, presupuesto y exclusiones gobiernan el primer alcance?

## Información que puedes proporcionar

- nombre del proyecto
- owners y responsabilidades
- resultado medible
- jurisdicciones
- fecha y presupuesto
- fuera de alcance

## Ejemplo ilustrativo

EJEMPLO, NO RESPUESTA ASUMIDA: Argentina; owner comercial Ana, owner técnico Luis; primer objetivo: publicar catálogo y convertir leads medidos; presupuesto mensual con tope aprobado; Brasil y pagos quedan fuera de P1.

## Cómo se comprobará

- confirmación del owner
- constitution y spec coherentes
- restricciones y métricas escritas

## Cuándo puede responderse NO_APLICA

Una subpregunta puede marcarse NO_APLICA sólo si no afecta el primer slice y se registra el motivo y el trigger que obligaría a reabrirla.

## Autoridades oficiales

- OpenAI — Keep agent instructions lean, state each instruction once, expose relevant tools and validate changes on representative work. — https://developers.openai.com/api/docs/guides/latest-model
- Google — Use curated, practical and adaptable launch questions, shared infrastructure, explicit dependencies, failure modes and staged rollout. — https://sre.google/sre-book/reliable-product-launches/
- NASA — Separate requirements, verification, validation, integration and stakeholder acceptance evidence. — https://www.nasa.gov/reference/system-engineering-handbook-appendix/

## Regla de avance

La ronda sólo se cierra con respuesta o evidencia real. Si falta una decisión o acceso material, el agente explica exactamente qué falta, conserva BLOCKED/AWAITING_USER y continúa con trabajo read-only que no dependa de inventarlo.
