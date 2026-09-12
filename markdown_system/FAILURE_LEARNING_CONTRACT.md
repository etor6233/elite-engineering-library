# Failure Learning Contract

## 1. Propósito

Cada fallo observado es evidencia de una frontera débil, un supuesto incorrecto o una condición no demostrada. No se borra cuando el siguiente intento pasa: se registra, se diagnostica, se corrige en la fuente canónica y se convierte en una prueba que impida repetirlo.

Este contrato es `AUTHORED` por Elite Engineering Library. No atribuye su metodología ni sus entradas a las empresas cuyos proyectos públicos se auditan.

## 2. Artefactos obligatorios

Todo proyecto que use la biblioteca mantiene desde discovery:

```text
PROJECT_FAILURE_LESSONS.md
```

Se crea desde `PROJECT_FAILURE_LESSONS_TEMPLATE.md`. Para mantenimiento de esta biblioteca, el índice compartido es `LIBRARY_FAILURE_LEARNING_LEDGER.md`; la evidencia técnica detallada permanece en `reconstruction_evidence/`.

Un PASS posterior no autoriza eliminar una entrada. Se puede redacted o archivar payload sensible, pero debe conservarse la huella, causa, resolución y evidencia no sensible.

## 3. Momento de captura

El agente abre o actualiza la entrada **en el mismo ciclo** en que observa cualquiera de estas condiciones:

- comando, build, test, migración, probe, despliegue, restore o rollback no cero;
- excepción, crash, timeout, deadlock, race, corrupción o resultado distinto del contrato;
- test omitido, deshabilitado, deselected o no ejecutado sin justificación ya demostrada;
- dependencia, licencia, notice, runtime, credencial, permiso, cuota o acceso ausente/incompatible;
- warning de seguridad, deprecación o resolución no determinista material;
- supuesto de negocio, datos, entorno o API demostrado falso;
- resultado “verde” que ocultó una aserción débil, fixture inválido o prueba que no ejercitó la ruta declarada;
- recurrencia de una huella previamente corregida.

No se registran como fallos los negativos esperados que prueban correctamente un control, salvo que el control acepte lo que debía rechazar.

## 4. Identidad y esquema mínimo

ID: `FAIL-YYYYMMDD-NNN`. La huella se calcula conceptualmente sobre `{fase, herramienta+versión, comando normalizado, clase de error, invariante}`, nunca sobre secretos o payloads privados.

El ledger de mantenimiento usa los namespaces monotónicos `LIB-FAIL-NNN` y `UP-FAIL-NNN`. `NNN` expresa un mínimo inicial de tres dígitos, no un máximo: al superar 999 continúa con cuatro o más dígitos sin reiniciar ni reutilizar identidades. El verificador debe reconocer `[0-9]{3,}`, exigir unicidad global y rechazar toda fila local abierta.

Cada entrada contiene:

```yaml
failure_id: FAIL-YYYYMMDD-001
fingerprint: "estable, sin secretos"
first_seen: "YYYY-MM-DDTHH:MM:SSZ"
last_seen: "YYYY-MM-DDTHH:MM:SSZ"
recurrence_count: 1
status: OPEN
severity: LOW|MEDIUM|HIGH|CRITICAL
classification: CONTRACT|CODE|DATA|DEPENDENCY|ENVIRONMENT|OPERATION|ASSUMPTION|LICENSE|SECURITY|PERFORMANCE|RECOVERY
phase: discovery|materialization|build|test|migration|integration|deploy|operation|restore|rollback|audit
capability: ""
scope: ""
command_redacted: ""
toolchain: []
error_signature_redacted: ""
violated_invariant: ""
blast_radius: ""
reproduction: ""
hypotheses_rejected: []
root_cause: ""
canonical_correction: ""
regression_test: ""
clean_rebuild_evidence: ""
dependent_gates_rerun: []
upstream_reference: ""
residual_risk: ""
owner: ""
next_action: ""
```

## 5. Estados y transiciones

```text
OPEN
→ DIAGNOSED
→ FIXED
→ REGRESSION_PROVEN
```

También se admiten:

- `BLOCKED_EXTERNAL`: la causa o verificación requiere acceso, cuenta, plataforma, datos o decisión externa;
- `UPSTREAM_OPEN`: existe en la revisión oficial fijada y no se parcheó silenciosamente;
- `REJECTED_COMPONENT`: licencia, seguridad o incompatibilidad impide usar el componente;
- `ACCEPTED_RISK`: sólo con owner, motivo, vencimiento, compensating controls y evidencia; nunca para corrupción, secreto expuesto o autorización rota.

`FIXED` sin regresión no cierra la lección. Sólo `REGRESSION_PROVEN` permite considerarla cerrada, y exige reconstrucción limpia más los gates dependientes. Una recurrencia reabre la entrada, incrementa `recurrence_count` y se considera señal de que el control preventivo era insuficiente.

## 6. Deduplicación y causalidad

- Misma huella y misma causa: actualizar la entrada e incrementar recurrencia.
- Síntoma igual con causa distinta: crear otro ID y enlazar `related_failures`.
- Varios errores independientes en un comando: una entrada por causa material, no un único renglón ambiguo.
- Un retry transitorio conserva cada intento y su presupuesto; no se reintenta indefinidamente.
- Un test omitido se mantiene `OPEN` o `BLOCKED_EXTERNAL` hasta demostrar por qué no pertenece a la lane seleccionada.

## 7. Bucle obligatorio

Aplicar `AGENT_ERROR_RECOVERY_PROTOCOL.md` y, además:

```text
capturar entrada
→ consultar huellas/lecciones anteriores
→ reproducir con mínimo alcance
→ corregir fuente canónica
→ añadir regresión
→ reconstruir desde vacío
→ repetir gates dependientes
→ actualizar estado/evidencia
→ promover la lección reusable cuando corresponda
```

El agente consulta `PROJECT_FAILURE_LESSONS.md` antes de elegir una herramienta, dependencia o patrón que tenga una huella relacionada. No repite una aproximación ya rechazada sin nueva evidencia explícita.

## 8. Promoción de una lección

Una lección puede pasar del proyecto al ledger de biblioteca sólo si:

1. la causa es generalizable y no una regla privada del negocio;
2. no contiene secretos, PII, documentos, endpoints privados ni telemetría irrestricta;
3. identifica la invariante que habría prevenido el fallo;
4. enlaza regresión y reconstrucción limpia, o conserva honestamente `UPSTREAM_OPEN`/`BLOCKED_EXTERNAL`;
5. no transforma código propio en supuesto código de una empresa líder.

La promoción puede fortalecer un contrato, pack, validator, fixture o gate. No debilita una aserción para obtener verde.

## 9. Stop conditions

Ante corrupción o estado durable ambiguo, secreto/PII expuesto, entorno productivo equivocado, autorización quebrada o acción destructiva no autorizada: detener mutaciones, preservar evidencia segura y seguir incident/recovery. Registrar el fallo no concede autoridad para continuar una acción riesgosa.

## 10. Readiness

`READY_TO_BUILD` exige:

- ledger creado y owner asignado;
- cero fallos `CRITICAL`/`HIGH` en `OPEN` o `DIAGNOSED` que afecten el primer vertical slice;
- toda condición upstream elegida enlazada al source lock y a un gate del proyecto;
- toda lección cerrada material para la lane seleccionada consultada e incorporada al plan;
- ningún fallo, warning material o skip ocultado para mejorar una métrica.

La existencia de fallos históricos cerrados no bloquea: acelera. Lo que bloquea es ignorarlos, no comprenderlos o carecer de evidencia para la capacidad que se pretende usar.

## 11. Integridad de texto durable en Windows

Leer y escribir JSON/JSONL y Markdown de control con UTF-8 explícito; ejecutar los scripts Python auxiliares con `-X utf8`. No depender de la codepage local ni de que el contenido actual sea ASCII. Antes de reserializar un cursor, comprobar textos no ASCII y sus límites. Si aparece mojibake, detener mutaciones, preservar bytes y eventos, y restaurar sólo mediante origen verificable o una transformación estrictamente reversible por campo; no normalizar artefactos oficiales ni reescribir eventos. Separar recuperación de texto y nuevos cambios de estado. Regresión ejecutable y evidencia: `reconstruction_evidence/UTF8_CHECKPOINT_RECOVERY_V324.md` (15tests, UTF-8 activado/desactivado).

## 12. Respuesta incierta de checkpoint

Una excepción o caída del proceso no prueba que el evento esté ausente. Primero validar los bytes y la cadena observados: un evento completo coincidente puede existir aunque el caller no haya recibido éxito; el append duplicado debe rechazarse sin mutación. Si hay cola parcial, detener mutaciones y preservar todos los bytes dañados antes de considerar recuperación desde un snapshot exacto demostrado; jamás adivinar registros faltantes. V326 prueba cuatro fronteras con archivos sintéticos y procesos hijos, no autoriza reparación automática de logs reales ni acredita cortes eléctricos o writers concurrentes. Evidencia: `reconstruction_evidence/ENTRY_LATENCY_AND_INTERRUPTION_V326.md`.

Un materializador interrumpido puede dejar un destino parcial: conservarlo, rechazar reentrada y reconstruir en un destino nuevo comprobado antes de seleccionarlo. Presencia de algunos archivos no demuestra finalización: contrastar manifest completo y bytes. V326 ejecuta8fronteras de archivo completo/parcial y error/terminación, preservando original/BOM y salida parcial. No atribuir atomicidad global del árbol ni seguridad de writers concurrentes.
