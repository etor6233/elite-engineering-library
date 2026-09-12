# Asesoramiento de proyecto — ronda E

Catalog-Version: 1.0.0
Catalog-SHA256: d5c336bf2c493e334ecc253aac5ed3d9e7c90d0ebd09fb78120da0c28bff6015
Topic-ID: E-INTEGRATION-ACCESS

Este prompt fue generado desde el catálogo ejecutable. El agente debe explicarlo en lenguaje simple, presentar un grupo manejable de preguntas y registrar la respuesta en los artefactos de readiness. No debe completar respuestas por inferencia ni pedir secretos en chat.

## Proveedores, cuentas y accesos demostrables

Una integración no está lista por conocer su nombre. Requiere producto y versión exactos, cuenta owner, sandbox, identidad, scopes, cuotas, costos, webhooks, reconciliación y evidencia de un probe seguro. Los secretos nunca se pegan en chat ni Markdown.

## Por qué se necesita

El código puede compilar y aun fallar por cuenta, país, permisos, contrato, límite o ambiente incorrectos. El probe evita simular acceso inexistente.

## Preguntas que el agente debe ayudarte a responder

- ¿Qué integraciones son REQUIRED, OPTIONAL o NO_APLICA y para qué journey?
- ¿Quién posee cada cuenta y cómo entregará una referencia segura de identidad sin revelar secretos?
- ¿Qué sandbox, scopes, cuotas, costos, webhooks, fixtures, reconciliation y soporte se demostrarán?

## Información que puedes proporcionar

- proveedor/producto/país
- owner de cuenta
- ambiente
- identity reference
- scopes
- cuotas y costos
- probe y evidencia
- fallback/exit plan

## Ejemplo ilustrativo

EJEMPLO, NO RESPUESTA ASUMIDA: Mercado Pago sandbox, owner Finanzas, identidad mp-sandbox-franchise en secret store, scopes mínimos, webhook firmado, tope de costo aprobado, fixture de pago y conciliación diaria; producción queda bloqueada hasta probe separado.

## Cómo se comprobará

- probe read-only o sandbox autorizado
- cuenta y ambiente observados
- request ID/schema redactados
- fallo cerrado ante ambiente equivocado

## Cuándo puede responderse NO_APLICA

Sólo si ningún journey aprobado depende del proveedor; debe registrarse el motivo y la condición futura que reabre la integración.

## Autoridades oficiales

- Google — Use curated, practical and adaptable launch questions, shared infrastructure, explicit dependencies, failure modes and staged rollout. — https://sre.google/sre-book/reliable-product-launches/
- OpenAI — Keep agent instructions lean, state each instruction once, expose relevant tools and validate changes on representative work. — https://developers.openai.com/api/docs/guides/latest-model

## Regla de avance

La ronda sólo se cierra con respuesta o evidencia real. Si falta una decisión o acceso material, el agente explica exactamente qué falta, conserva BLOCKED/AWAITING_USER y continúa con trabajo read-only que no dependa de inventarlo.
