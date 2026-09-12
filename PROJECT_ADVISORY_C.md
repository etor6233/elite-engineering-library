# Asesoramiento de proyecto — ronda C

Catalog-Version: 1.0.0
Catalog-SHA256: d5c336bf2c493e334ecc253aac5ed3d9e7c90d0ebd09fb78120da0c28bff6015
Topic-ID: C-ENTERPRISE-FISCAL

Este prompt fue generado desde el catálogo ejecutable. El agente debe explicarlo en lenguaje simple, presentar un grupo manejable de preguntas y registrar la respuesta en los artefactos de readiness. No debe completar respuestas por inferencia ni pedir secretos en chat.

## Modelo empresarial, stock, dinero y fiscalidad

Define productos, variantes, series, stock, compras, ventas, pagos, franquicias y contabilidad con sus estados y dueños. Regla fiscal aprobada es una decisión validada por responsable contable; certificado es la identidad criptográfica de la empresa; POS es el punto/canal de venta; homologación ARCA es el ambiente de prueba oficial antes de producción.

## Por qué se necesita

Estos conceptos generan compromisos de inventario, dinero, impuestos y auditoría. Inventarlos puede producir stock imposible, documentos fiscales inválidos o conciliaciones irreparables.

## Preguntas que el agente debe ayudarte a responder

- ¿Cuáles son las entidades, estados, invariantes, fuentes de verdad y reglas configurables del negocio?
- ¿Cómo fluyen compra, producción, importación, stock, precio, venta, cobro, devolución, regalía y postventa?
- ¿Qué autoridad contable aprueba reglas fiscales, puntos de venta, certificados, homologación y paso a producción ARCA?

## Información que puedes proporcionar

- glosario y estados
- owners por agregado
- reglas de stock y dinero
- política de franquicia
- contador/autoridad fiscal
- referencias seguras a certificado y ambientes

## Ejemplo ilustrativo

EJEMPLO, NO RESPUESTA ASUMIDA: serie única por vehículo; reserva no descuenta stock hasta confirmación; devolución genera reversa; contador externo aprueba IVA y comprobantes; certificado se entrega por secret store; primero homologación ARCA y luego producción con aprobación separada.

## Cómo se comprobará

- invariantes y concurrencia
- conciliación y reversa
- evidencia del responsable contable
- probe separado de homologación y producción

## Cuándo puede responderse NO_APLICA

Fiscalidad o una etapa comercial sólo puede ser NO_APLICA si el primer slice no emite, cobra, registra ni promete esa operación y existe trigger explícito para reabrirla.

## Autoridades oficiales

- Microsoft — Posting sales documents creates durable operational and accounting consequences and requires explicit business decisions. — https://learn.microsoft.com/en-us/dynamics365/business-central/ui-post-sales
- ARCA — Argentine electronic invoicing requires the applicable official service contract, authentication and homologation/production authority. — https://www.arca.gob.ar/ws/documentacion/manuales/manual-desarrollador-ARCA-COMPG.pdf
- NASA — Separate requirements, verification, validation, integration and stakeholder acceptance evidence. — https://www.nasa.gov/reference/system-engineering-handbook-appendix/

## Regla de avance

La ronda sólo se cierra con respuesta o evidencia real. Si falta una decisión o acceso material, el agente explica exactamente qué falta, conserva BLOCKED/AWAITING_USER y continúa con trabajo read-only que no dependa de inventarlo.
