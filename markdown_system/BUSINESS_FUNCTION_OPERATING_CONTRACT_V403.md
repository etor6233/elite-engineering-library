# Business Function Operating Contract — V403

Alcance: LIBRARY_INFRASTRUCTURE / FUNCTION_CONTRACT_SCAFFOLD. Este nuevo bloque preserva V402/checkpoint337. La implementación canónica única es [BUSINESS_FUNCTION_OPERATING_V403.md](../implementation_packs/BUSINESS_FUNCTION_OPERATING_V403.md), compatible con el materializador existente. Esta guía no duplica código ni crea otro formato de pack.

## Función, tarea y permiso

Grok Bot 101 es LEARNING_TRACK; los otros diez tracks son BUSINESS_FUNCTION. Ninguno concede acceso. La matriz completa incorpora en el JSON del pack inputs/outputs, aceptación por tarea, fuentes, permiso explícito, sistemas, alcance, owner y reanudación.

| Función | Tareas delimitadas | Requisito de permiso (sin grant) | Vista del contrato | Referencia existente |
|---|---|---|---|---|
| Grok Bot 101 | Delimitar un ejercicio de aprendizaje / Revisar la salida propuesta por un agente | learning:read, learning:review | #/functions/grok-bot-101 | DOCUMENTED_ONLY |
| Engineering | Delimitar cambio y aceptación / Verificar el cambio | engineering:plan, engineering:verify | #/functions/engineering | implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md |
| Product Managers | Enmarcar problema y resultado / Revisar feedback y decisión | product:plan, product:review | #/functions/product-managers | implementation_packs/GO_CUSTOMER_SURVEY_API.md |
| Founders | Registrar hipótesis empresarial / Revisar decisión de alcance | business:plan, business:review | #/functions/founders | DOCUMENTED_ONLY |
| Sales Engineering | Evaluar encaje técnico / Preparar una demostración acotada | sales-engineering:assess, sales-engineering:demo | #/functions/sales-engineering | implementation_packs/GO_FRANCHISE_CUSTOMER_JOURNEY_API.md |
| Sales | Revisar oportunidad / Preparar propuesta | sales:read, sales:request | #/functions/sales | implementation_packs/GO_FRANCHISE_CUSTOMER_JOURNEY_API.md |
| SDRs | Revisar procedencia del lead / Preparar handoff comercial | lead:read, lead:request | #/functions/sdrs | implementation_packs/GO_LEAD_CANDIDATE_PROMOTION.md |
| Customer Support | Clasificar consulta / Preparar respuesta o escalación | support:read, support:request | #/functions/customer-support | implementation_packs/GO_CONNECTED_CONVERSATION_RUNTIME.md |
| Marketing Operations | Preparar operación de campaña / Revisar estado y recuperación | marketing:request, marketing:read | #/functions/marketing-operations | implementation_packs/GO_CONNECTED_WHATSAPP_CAMPAIGNS.md |
| Post-Sales | Revisar entrega y onboarding / Preparar seguimiento o escalación | post-sales:read, post-sales:request | #/functions/post-sales | implementation_packs/GO_FRANCHISE_CUSTOMER_JOURNEY_API.md |
| Marketing | Preparar brief de campaña / Revisar aprendizaje de campaña | marketing:request, marketing:read | #/functions/marketing | implementation_packs/GO_CONNECTED_WHATSAPP_CAMPAIGNS.md |

Los nombres marketing:request y marketing:read se observaron en el owner existente de campañas. Los otros identificadores de permisos son requisitos propuestos para el consumidor, no nombres atribuidos a su producto. Cada asignación efectiva requiere principal autenticado, scope de organización y objeto, enforcement y evidencia. Founders no implica administrador.

La consola consulta contratos y referencias: acciones de envío, publicación, gasto o CRM están deshabilitadas sólo allí. Un producto consumidor puede ejecutar las operaciones que autorice y pruebe explícitamente. Los links a packs históricos conservan los claims de sus pruebas locales; no implementan Salesforce, HubSpot, Grok ni un negocio universal.

## Estados y cierre

SPEC_READY describe el contrato; OFFICIAL_METHOD_REFERENCED permite referencia documental estrecha; WAITING_GALAXY_ADMISSION se aplica exclusivamente a material futuro del evento. Las ocho fuentes actuales no esperan al evento. Nada atribuye currículo, métodos o capacidades todavía no publicados.

El JSON reusable conserva NOT_TESTED sin autoatestación. El PASS ejecutado es un receipt externo ligado al SHA exacto del contrato, fuentes y validator. Un PROVEN solicitado dentro del contrato exige receipt existente PASS con hash, función y claim scope. Producto/target siguen NOT_TESTED. La vista debe hacer visibles ambos niveles.

DONE sólo cierra la revisión y controles locales declarados, con cobertura frontend independiente y cero defectos conocidos abiertos en ese alcance según el owner existente. CONDITIONED/NOT_TESTED no se cuentan como PASS de una garantía requerida; un check de links no prueba usabilidad, negocio, proveedor, seguridad del target ni aceptación humana.

## Owners y reanudación

El fixture se identifica MAINTENANCE_FIXTURE y usa los owners reales del sandbox V403; nunca se copian como progreso de un consumidor. El binder recibe un mapa explícito de los seis owners existentes, resuelve el tasks_ref real, rechaza identidades y rutas de mantenimiento, y resetea claims/evidencias de verificación. Sólo escribe un contrato nuevo y no toca state/events.

La fuente de progreso sigue siendo tasks; state/events es cursor y cadena; PROJECT_FAILURE_LESSONS conserva incidentes. Antes de reanudar se usa ENGINEERING_EXECUTION_VALIDATOR1.3.1 del consumidor. No existe segundo backlog ni autorización de resume por OWNERS_BOUND_ONLY.

## Método y procedencia

Referencias verificadas: [Microsoft RBAC](https://learn.microsoft.com/en-us/azure/role-based-access-control/overview), [DORA user-centric focus](https://dora.dev/capabilities/user-centric-focus/), [DORA small batches](https://dora.dev/capabilities/working-in-small-batches/), [Microsoft prospect to quote](https://learn.microsoft.com/en-us/dynamics365/guidance/business-processes/prospect-to-quote-overview), [Microsoft case to resolution](https://learn.microsoft.com/en-us/dynamics365/guidance/business-processes/case-to-resolution-introduction), [Microsoft service to deliver](https://learn.microsoft.com/en-us/dynamics365/guidance/business-processes/service-to-cash-introduction), [HubSpot campaigns](https://knowledge.hubspot.com/campaigns/create-campaigns) y [xAI Grok overview](https://docs.x.ai/grok/overview). Consultas, claims y límites en method-sources.v403.json.

Son método documental y paráfrasis local; no código oficial, licencia de software ni ELITE_REFERENCE automática. La URL Microsoft service-to-cash conserva un nombre histórico; su texto explica service-to-deliver. Citar Dynamics no selecciona Business Central. Los hashes web son huellas históricas de observación sin cuerpos retenidos; los hashes de packs locales sí son comparables con sus bytes existentes.

Autoridades internas materiales: AGENT_SYSTEM_START; PROJECT_START_READINESS_GATE §2.1; ENGINEERING_EXECUTION_PLAYBOOK; FRONTEND_PRODUCT_ENGINEERING_UX §§1–2; SECURITY_SRE_CLOUD_INFRASTRUCTURE; PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD; AGENT_ERROR_RECOVERY_PROTOCOL y FAILURE_LEARNING_CONTRACT. Todos los scripts nuevos de este bloque son AUTHORED glue y no crean lógica comercial.
