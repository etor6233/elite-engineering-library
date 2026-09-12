# Especificación del cierre de biblioteca

Versión: 0.1.2. Estado: DISCOVERY. No autoriza expandir producto.
Owner decisor: usuario solicitante; ejecución técnica: agente dentro de la autoridad concedida.

## Resultado y baseline

Convertir los pendientes del roadmap existente en capacidades conectadas, reproducibles,
con procedencia y condiciones claras, utilizables tanto en NEW como en EXISTING.
Baseline V282: 160 packs, 1433 archivos materializables, 67/742 en el perfil integral.
Esto es inventario, no equivalencia con 160 capacidades terminadas.

## Requisitos y aceptación

| ID | Requisito del usuario | Owner del trabajo | Aceptación / evidencia y pendiente |
|---|---|---|---|
| LIB-R01 | iniciar o ampliar sin perder contexto ni cambios existentes | T2801 | V287 prueba entrada CLI NEW/EXISTING con bridge, composición, readiness BLOCKED, checkpoint, colisión, tamper y recuperación; 55 checks. Falta authority-map automático, assurance y expansión del producto, no confundir este ensayo con esos claims |
| LIB-R02 | código y conocimiento con procedencia real y calidad demostrada | T2801/T2803 | clasificación por archivo, admisión por claim, licencias y assurance; V280 detecta 21 cores cuya equivalencia sigue abierta |
| LIB-R03 | sistema conectado de punta a punta | T2802/T2804/T2805 | journey con interfaz/rol/auth/dominio/efecto/auditoría/respuesta/soporte/recovery; V278 es evidencia parcial, no cierra todo negocio |
| LIB-R04 | ingesta documental completa y confiable | T2806 | lanes y corpus por clase/campo, revisión y persistencia autorizadas; no asumir precisión universal |
| LIB-R05 | IA y comunicación útiles con costos controlados | T2805/T2807 | identidad, consentimiento, herramientas autorizadas, presupuesto, reconciliación y evals; cuentas/modelos reales separados |
| LIB-R06 | fallos registrados y correcciones reutilizables | T2801/T2809 | ledger, reproducción, fix canónico, regresión, checkpoint; fallos 384–387 siguen abiertos |
| LIB-R07 | dependencias y seguridad actualizadas sin gasto externo automático | T2803/T2808/T2809 | candidato exacto, advisory/SCA, pruebas, rollback y monitoring runtime; tests de Go V281 no certifican el grafo de producto |
| LIB-R08 | ensamblaje rápido y portable, sin backend TypeScript impuesto ni Git obligatorio | T2801/T2808 | composición desde Markdown y contexto por índice; medir ensayo representativo, no prometer latencia o tokens sin medición |
| LIB-R09 | archivo final reutilizable sin heredar aprobaciones | T2810 | artefacto exacto reconstruido, manifest/notices/receipts y exclusión local; V282 prueba selectores, no un ZIP final nuevo |

## Escenarios Given / When / Then

LIB-R10, aclarado por el usuario en V366, owners T2805/T2807 con T2803/T2806:
preparar código admitido para entrenar/adaptar un modelo del proyecto consumidor
con historial humano pertinente. NEW inicia sin historial; EXISTING inventaría su
fuente dentro del proyecto. Originales/derivados privados, entrenamiento explícito,
evaluación independiente, promoción y runtime son etapas separadas; memoria/RAG
no sustituyen el requisito. No se piden chats reales para preparar la biblioteca.
Modelo/técnica/ubicación/presupuesto/permisos se fijan por proyecto antes del job.
Contrato normativo: markdown_system/PROJECT_HISTORY_MODEL_TRAINING_CONTRACT.md.
V288 conserva la interpretación anterior y su investigación; V366 la corrige sin
afirmar un pipeline implementado, calidad demostrada ni autorización de entrenar.

Clarificación explícita V289: conservar los bytes originales tal como los entregue
la fuente autorizada, antes de interpretar. Normalización, deduplicación, redacción
y anotaciones pertenecen a derivados trazables, no reemplazan el original. Un hash
no demuestra comprensión semántica ni que el proveedor haya entregado todo su
historial. Desconocidos/ambigüedades se registran, no se fuerzan. Retención, borrado,
acceso y coste requieren política del proyecto; raw no significa conservar PII
para siempre ni enviarla íntegra a un modelo. TEST-09 del contrato local expresa
la aceptación pendiente, no una importación ya implementada.

- Dada una biblioteca válida y destino nuevo, al seleccionar packs compatibles, se materializan sólo sus archivos exactos y se conserva receipt. Admisión CANDIDATE, colisión o hash incorrecto rechazan antes de incorporar código.
- Dado un proyecto existente, al reanudar, se valida baseline/delta/historial y se limita el cambio al alcance autorizado. No se sobrescriben archivos ajenos ni se hereda readiness de otra tarea.
- Dado un fallo de un componente, la prueba negativa se conserva; sólo una corrección canónica reconstruida y sus gates dependientes permiten cerrarlo.
- Dado un servicio sin cuenta o un documento sin corpus, el agente conserva el bloqueo correspondiente y continúa sólo el trabajo independiente autorizado. No emula una llamada live como prueba real.

## Exclusiones del mantenimiento actual

No desplegar una franquicia, iniciar campañas, enviar mensajes reales, cobrar, emitir documentos fiscales,
crear recursos cloud ni declarar leyes/configuraciones de un negocio todavía no definido.
Estas exclusiones de efectos reales no eliminan sus capacidades del roadmap de biblioteca.

## Cierre pendiente

Cada fila anterior necesita evidencia integrada y versionada, no sólo un path existente.
El detalle de tareas vive exclusivamente en REUSABLE_CODE_READINESS_ROADMAP.md §10.
Antes de expansión: readiness semántico y implementation_assurance reales.

## Significado verificable de biblioteca cerrada — aclaración V368

La pregunta del usuario exige precisar qué significa que todo funcione al cerrar.
Aplicar ENGINEERING_EXECUTION_PLAYBOOK.md §§22–23 al alcance comprometido:

- Cada requisito LIB-R01–LIB-R10 debe tener comportamiento implementado y
  evidencia de ejecución reproducible en las composiciones/entornos declarados.
  Un documento, archivo presente, compilación o fixture aislado no reemplaza
  las pruebas de integración que requiera su claim.
- Los recorridos reutilizables prometidos se prueban conectados, incluyendo
  errores, seguridad, compatibilidad, límites de recursos y recuperación según
  riesgo. NEW y EXISTING conservan sus oráculos propios.
- Cero bloqueos pendientes que contradigan ese alcance; los gates de evidencia,
  implementation_assurance y checkpoint completo deben pasar, junto con el
  artefacto final exacto. No reducir requisitos ni convertir bloqueos en
  OPTIONAL/NONE_WITH_REASON para fabricar el cierre.
- El catálogo puede conservar candidatos/rechazos como investigación identificada;
  eso no los convierte en capacidades listas ni autoriza prometerlas. Una
  capability comprometida sin implementación admitida sigue bloqueando el cierre.
- Quedan para cada consumer sus reglas, datos, cuentas, secretos, infraestructura,
  volúmenes, requisitos regulatorios y aceptación del negocio. La biblioteca
  debe aportar interfaces y controles ejecutados para esas variaciones; no puede
  esconder un defecto genérico como futura configuración del proyecto.

Funcionar significa cumplir los contratos probados en versiones y condiciones
declaradas. No garantiza ausencia de errores desconocidos ni compatibilidad con
cualquier combinación o cambio futuro. El cierre de la biblioteca no es una
certificación de producción heredable por otro proyecto. Esta aclaración no
cambia oráculos ni estados de los48controles y no cierra ningún pendiente.


## Decisión sucesora V402 — alcance de infraestructura local

La instrucción explícita del usuario en reconstruction_evidence/LIBRARY_INFRA_SCOPE_V402.md
rige el cierre desde checkpoint272. READY significa biblioteca/infra lista para
usar con referencia materializada fuera de la fuente, contratos y pruebas locales
completos. Producción live permanece separada. No solicitar secretos; ARCA se
completa penúltimo como infraestructura, Daybreak último queda diferido y no impide
el cierre infra. No confundir credenciales pendientes con código incompleto.
AUTHORED queda limitado a glue inevitable: cualquier algoritmo de negocio local
sin source ADAPTED/VERBATIM/DEPENDENCY_PIN admisible requiere resolución real.
Los estados/evidencias previos se conservan; no se marca ningún pendiente PASS por
registrar el nuevo alcance. Orden exacto y nueve criterios en el expediente V402.
