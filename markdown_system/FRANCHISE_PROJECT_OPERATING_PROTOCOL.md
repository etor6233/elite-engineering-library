# Protocolo operativo de proyecto — conexión 1.0.0

Alcance: coordinar el uso de la biblioteca en un proyecto consumidor y conservar continuidad sin memoria de chat. Composición metodológica y glue `AUTHORED`; no código atribuido a Google, Microsoft, xAI ni otro proveedor. No altera el cierre V402/337 ni autoriza construir o desplegar un producto por el solo hecho de instalar el bridge.

Entrada: [START_FRANCHISE](../START_FRANCHISE.md). Instalador complementario y verificador: [PROJECT_OPERATING_CONNECTION](../implementation_packs/PROJECT_OPERATING_CONNECTION.md). En el proyecto, `AGENTS.md` y `CLAUDE.md` remiten a `PROJECT_AGENT_ENTRY.md`, que verifica `.elite/operating-connection.json` y conduce a este protocolo. Otro agente recibe explícitamente esa misma entrada mediante su plantilla. Las reglas del usuario y las políticas de la plataforma siguen aplicando.

## 1. Dónde vive cada cosa

| Owner existente | Biblioteca | Proyecto consumidor |
|---|---|---|
| Método y procedencia | Contratos, mapas, packs, versiones, fuentes oficiales, lecciones reutilizables | `PROJECT_AUTHORITY_MAP.md`, `PROJECT_EXTERNAL_SOURCE_LOCK.md`: selección y condiciones aplicables |
| Alcance y decisiones | Templates y contratos de blueprint/readiness | `PROJECT_BLUEPRINT.md`, readiness y rondas A–H, decisiones y aceptación de ese negocio |
| Implementación | Packs compatibles; ningún estado de cliente | Código, configuración sin secretos, `PROJECT_PACK_PLAN.md`, spec/plan/tasks existentes |
| Continuidad | Execution validator y sus templates | `PROJECT_EXECUTION_STATE.json` y `PROJECT_EXECUTION_EVENTS.jsonl` |
| Fallos y mantenimiento | Protocolos de recuperación y ledger reusable | `PROJECT_FAILURE_LESSONS.md`, dependency/freshness/monitoring records y evidencias |
| Conexión | Este protocolo y su pack | Entrada y receipt de conexión; sólo rutas, versiones y hashes, nunca un segundo backlog |

La biblioteca se consulta como autoridad de sólo lectura durante el trabajo del producto. Una corrección reusable se propone en un expediente de biblioteca separado, con reproducción y procedencia, sin datos del cliente. El agente consumidor no cambia un pack canónico ni un cierre protegido por su cuenta. No copiar el estado 337, su event log ni sus aprobaciones al proyecto.

## 2. Secuencia obligatoria

Cada paso tiene una salida observable. Al retomar se ejecutan 1–3 y luego el primer paso pendiente; no se repite todo el inicio si su evidencia sigue siendo válida. Los gates normativos conservan sus exigencias.

| Paso | Lectura y acción mínima | Salida / continuación |
|---|---|---|
| 1. Situarse | Abrir la raíz real del proyecto y leer `PROJECT_AGENT_ENTRY.md`. Ejecutar el verificador de conexión; contrastar autoridad y revisión fijadas. | PASS de integridad o discrepancia exacta. Una ruta rota no autoriza elegir otra biblioteca silenciosamente. |
| 2. Detectar NEW/EXISTING | Inspeccionar archivos antes de escribir. NEW registra vacío/scaffold observado; EXISTING registra inventario y delta autorizado. El scaffolding del bridge cuenta como scaffold. | Baseline propio; nunca heredar el baseline del mantenimiento de biblioteca. |
| 3. Recuperar contexto | Si existe el par estado/eventos, validar cadena, revisión y hashes con el execution validator. Leer resumen, `context.must_read_refs`, fallos abiertos y `next_action`. Resolver referencias antes de actuar. | Próxima acción sustentada por archivos. Reutilizar `reuse_without_reload_refs` sólo mientras sus condiciones sigan válidas. Un par incompleto se recupera, no se recrea ocultando historia. |
| 4. Preparar el proyecto | Leer [readiness](PROJECT_START_READINESS_GATE.md), [autonomía](AGENT_AUTONOMY_CONTRACT.md) y [blueprint](PROJECT_BLUEPRINT_CONTRACT.md). Materializar sus templates y el execution kit. Generar cada ronda pendiente con el renderer y registrar respuestas reales. | Gates aplicables al target; sólo implementar producto después de READY_TO_BUILD. Un ensayo de conexión no emite ese permiso. |
| 5. Planificar todo el alcance | Clasificar las 48 superficies de [TOTAL_SYSTEM_CAPABILITY_CONTRACT](TOTAL_SYSTEM_CAPABILITY_CONTRACT.md), conectar journeys y criterios de aceptación completos. Usar spec/plan/tasks como backlog único. | REQUIRED, OPTIONAL, NONE_WITH_REASON o BLOCKED justificados. No trasladar exclusiones de una referencia a otro negocio. El alcance final completo se conserva durante la ejecución incremental. |
| 6. Elegir autoridad y componentes | Consultar mapas por la tarea seleccionada; registrar fuentes, packs, revisión, SHA, condiciones y archivos afectados en authority map / pack plan existentes. | Componente compatible admitido. Si falta, seguir [CAPABILITY_GAP_RESOLUTION_PACK_PLAN](CAPABILITY_GAP_RESOLUTION_PACK_PLAN.md), sin inventar equivalencia ni procedencia. |
| 7. Preparar entrega y entornos | Identificar OS/arquitectura, datos, identidades y runtime del target; definir CI y staging temprano junto al plan. Registrar diferencias con local y cómo probarlas. | Build y despliegue reproducibles para ese target. La referencia Windows local no acredita Linux ni cloud. |
| 8. Ejecutar un incremento conectado | Elegir una tarea con dependencias resueltas y aceptación definida. Aplicar packs y cambios mínimos; integrar con el mismo producto y sus journeys. | Artefacto observable, cambios identificados y comprobación pertinente. Avanzar autónomamente dentro del alcance autorizado. |
| 9. Verificar lo afectado | Ejecutar contratos y pruebas exigidos por el delta, riesgo y entorno. Capturar comando, revisión, salida, exit code y hashes. Ampliar sólo por impacto o exigencia normativa. | Evidencia que prueba ese claim; ningún PASS de formato sustituye una ejecución. |
| 10. Recuperar o continuar | Ante un fallo real, registrar firma, causa, alcance, corrección, prueba y condición de reapertura en [failure lessons](FAILURE_LEARNING_CONTRACT.md). Consultar lecciones pertinentes. | Fallo contenido/corregido con evidencia; repetir sólo verificaciones afectadas. Si depende de un humano, dejarlo visible y continuar trabajo independiente autorizado. |
| 11. Checkpoint material | Un único escritor concilia tasks, referencias de evidencia y estado; incrementa revisión y añade evento mediante el writer admitido. Verifica el par final. | Próxima acción exacta y contexto mínimo. No checkpoint por cada lectura/comando ni segundo registro de progreso. |
| 12. Retomar, transferir o cerrar | Siguiente agente vuelve al paso 1. Para cerrar, comprobar todos los REQUIRED y el gate del entorno solicitado. | Cierre del alcance demostrado, o siguiente acción con dueño/dependencia. Nunca declarar producción por READY de biblioteca. |

## 3. Bootstrap y continuidad ejecutables

El instalador no inventa negocio, readiness ni estado. Si no hay execution kit, materializar desde la biblioteca fijada `implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md` en un destino ausente, recomendado `.elite/execution`; sus archivos quedan en `.elite/execution/engineering_execution_kit`. Copiar `execution_state.template.json` como candidato y completar sólo hechos observados. Crear los owners del proyecto desde sus templates, en particular el failure ledger requerido desde el primer checkpoint. Ver [ENGINEERING_EXECUTION_PLAYBOOK](../ENGINEERING_EXECUTION_PLAYBOOK.md).

Validación de reanudación, desde la raíz del proyecto (Python 3.14+):

```powershell
python -B .elite/execution/engineering_execution_kit/validate_execution_state.py PROJECT_EXECUTION_STATE.json --project-root . --events PROJECT_EXECUTION_EVENTS.jsonl --level resume
```

Para checkpoint, el complemento ofrece una exclusión de escritor y comprobación de revisión alrededor del writer original, sin sustituir su schema o cadena. Preparar un candidato separado dentro del proyecto, con evidencia y `revision = anterior + 1`, y ejecutar desde la raíz:

```powershell
# Sustituir revisión e identificador por los reales; 0 sólo para el primer estado.
python -B .elite/connection-tools/operating_connection/connection.py checkpoint --project-root . --candidate .elite/state-candidate.json --expected-revision 0 --event-id CONNECTION-001 --event-type BASELINE --actor project-agent --summary 'Baseline observado y referencias fijadas'
```

El lock `.elite/operating-write.lock` impide escritores simultáneos que usan este comando. Un lock existente exige identificar al dueño; no borrarlo automáticamente por edad. No coordina equipos en máquinas diferentes: allí asignar un solo integrador y conservar revisión base/commit y entrega por tarea en el plan existente. Ningún proceso puede prometer atomicidad de dos archivos ante corte de energía: un par interrumpido se detecta al retomar y se recupera preservando ambos bytes, nunca se borra el event log para obtener PASS.

`next_action` debe incluir tarea, archivo/owner a consultar, acción concreta, comprobación de éxito y dependencia si existe. Ejemplo: «T120: leer spec de conciliación, corregir el mapping del adapter seleccionado, ejecutar su contrato C03; registrar receipt y actualizar tasks; requiere fixture F02». Un «continuar» sin objetivo ni comprobación no sirve.

## 4. Velocidad y control del trabajo

- Mantener un router corto. Al cambiar de sesión validar la conexión y el cursor; abrir sólo referencias necesarias para la siguiente decisión. La lectura inicial obligatoria de autoridades sigue vigente.
- Reutilizar evidencia cuando artefactos, dependencias, configuración, entorno, alcance y vigencia sigan siendo los probados. Cambio material, expiración, advisory oficial o contradicción invalida la parte afectada. Un hash idéntico por sí solo no prueba vigencia de una cuenta o servicio externo.
- Registrar en la tarea por qué un control se ejecuta, se reutiliza o debe renovarse. No agregar una suite repetida para aparentar progreso. No omitir controles requeridos por el contrato.
- Tomar decisiones reversibles autorizadas sin pedir confirmación ceremonial. Las reglas comerciales, gasto, datos sensibles y acciones externas conservan sus límites reales.
- Si una tarea se bloquea, declarar causa, trabajo dependiente y trigger de reapertura; elegir otra tarea independiente del mismo plan. No marcar la bloqueada como terminada para desbloquear el tablero.
- Un escritor del estado/eventos y un integrador por revisión. Otros agentes entregan artefactos y evidencia sobre tareas/paths asignados; no editan el mismo estado simultáneamente.
- Conservar el mapa de alcance completo; demostrar avances conectados en cada incremento. No prometer un plazo menor a una semana sin medir el alcance y sus dependencias.

## 5. Criterios separados de entorno

| Alcance | Evidencia que permite avanzar | Lo que no hereda |
|---|---|---|
| Biblioteca | Packs, locks, reconstrucción y verificaciones de su referencia | Reglas y aceptación de otro negocio |
| Desarrollo local | Gate PROJECT aplicable, ejecución del incremento y fixtures identificados | Credenciales, callbacks o aceptación live |
| CI | Build del target, pruebas pertinentes, procedencia y artefacto versionado | Aprobación de usuarios y operación productiva |
| Staging | Entorno representativo, integraciones autorizadas, permisos, migraciones, observabilidad y recovery ensayados | Producción automática por tener una URL funcional |
| Producción | Aceptación del alcance completo solicitado, identidades/datos reales, seguridad, operación, carga, recuperación y despliegue aprobados | Ningún PASS de biblioteca/local sustituye esos gates |

Preparar staging durante el desarrollo y promover el artefacto verificado para el target, manteniendo configuración y secretos fuera del paquete. Una dependencia externa pendiente debe quedar explícita en el gate correspondiente; no clasificarla artificialmente como opcional para implementar el producto. Se puede avanzar en probes, planificación y trabajo independiente permitido por el estado actual.

## 6. Agentes y adopción por plataforma

La entrada neutral es `PROJECT_AGENT_ENTRY.md`; el complemento añade esa ruta a lo instalado por el bridge original y entrega `.elite/platform-adapter.template.md`. Completar placeholders en el proyecto para su plataforma/versión, rutas realmente accesibles, rol y permisos. Los templates no instalan una Skill en Grok ni prueban que un host concreto descubra instrucciones.

Una prueba de harness demuestra archivos y validadores. Una sesión de agente sin conversación previa demuestra el recorrido observado en ese entorno. Codex/Claude/Grok/otros sólo reciben PASS de ejecución cuando fueron ejecutados en su plataforma identificada; en caso contrario `NOT_TESTED`. La copia local o cloud debe usar la revisión fijada; nunca cambiar de biblioteca por memoria del bot. Ventas/SDR/Marketing/Support requieren procedimientos de función y permisos del producto; este protocolo de conexión no los certifica.

## 7. Fuentes oficiales y aplicación comprobable

Consultadas el 2026-09-13. Los snapshots y SHA de consulta se conservan en el expediente de conexión; son fuentes metodológicas, no dependencias ejecutables ni autoría del glue.

| Fuente primaria | Regla aplicada aquí | Comprobación de la conexión |
|---|---|---|
| [DORA / Google Cloud: Working in small batches](https://dora.dev/capabilities/working-in-small-batches/) | Incrementos verificables y feedback breve, con tareas y alcance visibles | Tarea acotada con salida y próxima acción, sin repetir todo el inicio |
| [DORA: Deployment automation](https://dora.dev/capabilities/deployment-automation/) | Paquete identificable, configuración separada y proceso repetible | Protocolo distingue fuente, entorno y claim; el ensayo no dice haber desplegado |
| [Microsoft Well-Architected: deployment/testing](https://learn.microsoft.com/en-us/azure/well-architected/mission-critical/mission-critical-deployment-testing) | Preparar entornos temprano y validar condiciones representativas | Matriz local/CI/staging/producción; ninguna promoción heredada |
| [Microsoft Well-Architected: safe deployments](https://learn.microsoft.com/en-us/azure/well-architected/operational-excellence/safe-deployments) | Detectar fallos y recuperar antes de promover | Fallo provocado, evidencia y recuperación; no aprobación por texto |
| [Grok: skills and routines](https://docs.x.ai/grok-bot/skills-routines-and-automations) | Entradas, secuencia, validación, límites; probar antes de automatizar | Plantilla explícita y NOT_TESTED sin ejecución de plataforma |

## 8. Conexión DONE

El expediente debe enumerar cada control con PASS/FAIL, comandos/exit codes, archivos y SHA, entorno y límites: arranque conectado; instalación repetible; preservación de instrucciones ajenas; tarea real acotada; fallo y recuperación; sesión nueva sin chat; detección de referencia inválida; exclusión de escritor y rechazo de revisión obsoleta; evidencia verificable; V402 intacto. Una comprobación pendiente impide declarar DONE de la conexión. El PASS de la conexión no abre una franquicia real: requiere una instrucción posterior del usuario.
