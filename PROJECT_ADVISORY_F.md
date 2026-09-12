# Asesoramiento de proyecto — ronda F

Catalog-Version: 1.0.0
Catalog-SHA256: d5c336bf2c493e334ecc253aac5ed3d9e7c90d0ebd09fb78120da0c28bff6015
Topic-ID: F-UX-CHANNELS

Este prompt fue generado desde el catálogo ejecutable. El agente debe explicarlo en lenguaje simple, presentar un grupo manejable de preguntas y registrar la respuesta en los artefactos de readiness. No debe completar respuestas por inferencia ni pedir secretos en chat.

## Experiencia completa para público, clientes y empresa

UX completa significa que cada actor entiende qué ocurre, puede corregir errores y finaliza su objetivo en web pública, cliente, franquicia, administración, fábrica o proveedor. Incluye responsive, accesibilidad, rendimiento, contenido, contacto, estados vacíos/error/espera, ayuda contextual, capacitación y soporte vinculados a la misma versión operativa.

## Por qué se necesita

Una API correcta no garantiza que visitantes conviertan, clientes confíen ni operadores trabajen sin atajos manuales. Una pantalla, capacitación o respuesta de soporte desconectada de permisos, efectos durables y versión tampoco cierra una capacidad. Los canales y journeys deben validarse de extremo a extremo con usuarios y dispositivos reales.

## Preguntas que el agente debe ayudarte a responder

- ¿Qué pantallas y canales necesita cada actor para completar cada journey P1?
- ¿Qué branding, contenido, SEO, analytics, consentimiento, accesibilidad, dispositivos y navegadores se aceptan?
- ¿Qué presupuestos de rendimiento y pruebas de error, espera, recuperación y entrega final gobiernan la aceptación?
- ¿Qué ayuda contextual, práctica, capacitación, evaluación, soporte y reentrenamiento necesita cada rol, y cómo se ligan a permisos y versión?

## Información que puedes proporcionar

- mapa de pantallas
- contenido y assets autorizados
- matriz de dispositivos/navegadores
- criterios de accesibilidad
- performance budgets
- owners editoriales
- matriz rol→journey→ayuda→capacitación→soporte
- versiones de contenido operativo

## Ejemplo ilustrativo

EJEMPLO, NO RESPUESTA ASUMIDA: visitante móvil filtra modelos y contacta; cliente ve turno y entrega; franquicia gestiona disponibilidad con ayuda contextual y práctica versionada; admin audita; soporte recibe contexto autorizado y runbook de la misma release; Chrome/Edge/Safari actuales, teclado y lector de pantalla; p75 web vitals objetivo y textos de error aprobados.

## Cómo se comprobará

- E2E por actor y efecto durable
- navegadores y responsive
- accesibilidad
- métricas de rendimiento
- ayuda/capacitación/soporte iguales a la release
- reentrenamiento ante cambio de journey
- aceptación del owner de negocio

## Cuándo puede responderse NO_APLICA

Un canal puede ser NO_APLICA si ningún journey P1 lo usa. Estados de error y recuperación no son opcionales para una pantalla incluida; ayuda, capacitación o soporte sólo pueden omitirse con razón aprobada cuando ninguna persona opera ni recibe el efecto.

## Autoridades oficiales

- Google — Use curated, practical and adaptable launch questions, shared infrastructure, explicit dependencies, failure modes and staged rollout. — https://sre.google/sre-book/reliable-product-launches/
- NASA — Separate requirements, verification, validation, integration and stakeholder acceptance evidence. — https://www.nasa.gov/reference/system-engineering-handbook-appendix/

## Regla de avance

La ronda sólo se cierra con respuesta o evidencia real. Si falta una decisión o acceso material, el agente explica exactamente qué falta, conserva BLOCKED/AWAITING_USER y continúa con trabajo read-only que no dependa de inventarlo.
