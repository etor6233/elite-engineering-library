# Asesoramiento de proyecto — ronda H

Catalog-Version: 1.0.0
Catalog-SHA256: d5c336bf2c493e334ecc253aac5ed3d9e7c90d0ebd09fb78120da0c28bff6015
Topic-ID: H-OPERATIONS-PRODUCTION

Este prompt fue generado desde el catálogo ejecutable. El agente debe explicarlo en lenguaje simple, presentar un grupo manejable de preguntas y registrar la respuesta en los artefactos de readiness. No debe completar respuestas por inferencia ni pedir secretos en chat.

## Plataforma, operación, recuperación y producción

Distingue local, CI, DEV, STAGE y PROD. SLI/SLO mide el servicio; RPO/RTO define pérdida y tiempo de recuperación; PITR restaura a un instante; canary expone gradualmente; rollback revierte. CDN/WAF, IdP, PostgreSQL, proveedores y observabilidad deben existir realmente en el target elegido.

## Por qué se necesita

Compilar y probar localmente no demuestra capacidad, seguridad, recuperación ni operación productiva. La aceptación necesita el artefacto y entorno reales con owners y evidencia.

## Preguntas que el agente debe ayudarte a responder

- ¿Qué plataforma, regiones, red, DNS/TLS, CDN/WAF, identidades, IaC y entornos se usarán?
- ¿Cuáles son SLO, carga, capacidad, costo, alertas, on-call, runbooks y respuesta a incidentes?
- ¿Cómo se probarán backup/PITR/restore, RPO/RTO, canary, rollback, seguridad y piloto?
- ¿Qué evidencia y qué autoridad separada habilitan DEV, piloto y producción?

## Información que puedes proporcionar

- target y accesos
- topología e IaC
- SLO/carga/costo
- on-call/runbooks
- backup/restore
- rollout/rollback
- acceptance owners

## Ejemplo ilustrativo

EJEMPLO, NO RESPUESTA ASUMIDA: DEV en región aprobada con IaC; PostgreSQL gestionado; CDN/WAF e IdP definidos; SLO y carga medidos; restore en entorno aislado; canary al 5%; rollback automático; owner técnico habilita DEV y comité negocio-seguridad autoriza piloto/PROD con evidencia separada.

## Cómo se comprobará

- deploy del artefacto exacto
- prueba de carga y seguridad
- restore/PITR observado
- canary y rollback
- aceptación empresarial del piloto

## Cuándo puede responderse NO_APLICA

Una tecnología puede ser NO_APLICA por arquitectura demostrada; producción nunca se considera cubierta por evidencia sólo local o DEV.

## Autoridades oficiales

- Google — Use curated, practical and adaptable launch questions, shared infrastructure, explicit dependencies, failure modes and staged rollout. — https://sre.google/sre-book/reliable-product-launches/
- NASA — Separate requirements, verification, validation, integration and stakeholder acceptance evidence. — https://www.nasa.gov/reference/system-engineering-handbook-appendix/
- OpenAI — Keep agent instructions lean, state each instruction once, expose relevant tools and validate changes on representative work. — https://developers.openai.com/api/docs/guides/latest-model

## Regla de avance

La ronda sólo se cierra con respuesta o evidencia real. Si falta una decisión o acceso material, el agente explica exactamente qué falta, conserva BLOCKED/AWAITING_USER y continúa con trabajo read-only que no dependa de inventarlo.
