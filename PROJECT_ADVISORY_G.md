# Asesoramiento de proyecto — ronda G

Catalog-Version: 1.0.0
Catalog-SHA256: d5c336bf2c493e334ecc253aac5ed3d9e7c90d0ebd09fb78120da0c28bff6015
Topic-ID: G-SECURITY-PRIVACY

Este prompt fue generado desde el catálogo ejecutable. El agente debe explicarlo en lenguaje simple, presentar un grupo manejable de preguntas y registrar la respuesta en los artefactos de readiness. No debe completar respuestas por inferencia ni pedir secretos en chat.

## Datos, aislamiento, seguridad y cumplimiento

Clasifica cada dato, propósito, acceso y retención. Aislamiento significa que una franquicia no puede leer ni modificar otra. IdP autentica identidades; autorización decide acciones; KMS/PKI protege claves y certificados; threat model enumera abusos y defensas.

## Por qué se necesita

Sin límites de datos e identidad, cualquier módulo puede filtrar información, mezclar franquicias o usar privilegios excesivos aunque sus funciones de negocio pasen tests.

## Preguntas que el agente debe ayudarte a responder

- ¿Qué datos son personales, financieros, secretos o restringidos y con qué propósito, base, retención y borrado?
- ¿Qué IdP, MFA, sesiones, roles, políticas por organización y break-glass se usarán?
- ¿Qué amenazas, fraude, abuso, logs, auditoría, pruebas ofensivas y acceptance owner deben demostrarse?

## Información que puedes proporcionar

- clasificación por campo
- política de tenant
- IdP y roles
- secret/KMS references
- retención y derechos
- threat model
- plan de pruebas

## Ejemplo ilustrativo

EJEMPLO, NO RESPUESTA ASUMIDA: email CONFIDENTIAL por atención comercial; retención aprobada; OIDC con MFA para administración; policy por organization_id; certificados sólo por referencia KMS; pruebas negativas entre dos franquicias y revisión ofensiva antes de piloto.

## Cómo se comprobará

- tests de aislamiento
- tokens y scopes reales redactados
- secret scanning
- audit logs sin PII innecesaria
- security acceptance

## Cuándo puede responderse NO_APLICA

Una categoría de dato puede ser NO_APLICA si se demuestra que no se recoge ni deriva. Autorización y aislamiento no son NO_APLICA cuando existen usuarios u organizaciones múltiples.

## Autoridades oficiales

- Google — Use curated, practical and adaptable launch questions, shared infrastructure, explicit dependencies, failure modes and staged rollout. — https://sre.google/sre-book/reliable-product-launches/
- NASA — Separate requirements, verification, validation, integration and stakeholder acceptance evidence. — https://www.nasa.gov/reference/system-engineering-handbook-appendix/

## Regla de avance

La ronda sólo se cierra con respuesta o evidencia real. Si falta una decisión o acceso material, el agente explica exactamente qué falta, conserva BLOCKED/AWAITING_USER y continúa con trabajo read-only que no dependa de inventarlo.
