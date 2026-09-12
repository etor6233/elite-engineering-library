# V288 — chats históricos para atención realtime de bajo costo

2026-09-07. Requisito explícito del usuario: incorporar esta capacidad a proyectos
existentes aunque no tengan módulo de chat, aprovechando conversaciones atendidas
por personas. Investigación de mantenimiento, **no importador implementado ni
capacidad REUSABLE_PACK**. Amplía T2805/T2807 con T2803/T2806 como dependencias;
no crea otro CRM, broker, roadmap o runtime conversacional.

## 1. Lo que existe y lo que no está demostrado

| Owner inspeccionado | Evidencia de alcance | Lo que NO demuestra |
|---|---|---|
| GO_CONNECTED_CONVERSATION_RUNTIME | History consulta communication.conversation_turn, state=completed, tenant/canal/contacto/hilo y límite 0–20; Handle puede ejecutar tools y responder | recuperar conversaciones humanas anteriores a la instalación o importar un archivo histórico |
| PYTHON_META_WHATSAPP_CLOUD_ADAPTER 0.12.0 | webhook messages, envío aprobado, estados, inbox/job y lectura de notificaciones | sincronización histórica completa de WhatsApp Business del teléfono, Instagram o CRM |
| GO_PG_CONTACT_CHANNEL_IDENTITY | asociación de identidad externa y lead, scope y revocación | identificar automáticamente a una persona sólo por nombres/texto o reutilizar IDs entre negocios |
| OFFICIAL_UPSTREAM_ACQUISITION_CORE | Presidio 2.2.364, commit 779dbd286d5ef4d1fbe2514275fb1bce358f2417, MIT, PINNED_CANDIDATE | redacción perfecta, integración histórica ni SCA vigente de un pipeline elegido |
| AI_SECURITY_GOVERNANCE_PRIVACY / AI_ENGINEERING_MASTER_MAP | privacidad/ACL, poisoning, separación de evals y entrenamiento | código ejecutable de importación/corpus/retrieval ya admitido |

Por tanto: **hay piezas reutilizables condicionadas y conocimiento; no está
demostrada la cadena histórica humana → conocimiento aprobado → respuesta realtime**.
Una lista vacía de resultados de búsqueda no se usa como prueba de ausencia universal.

## 2. Arquitectura de trabajo propuesta, AUTHORED

```text
fuente histórica autorizada                  nuevos mensajes autorizados
          ↓                                            ↓
importación sin efectos                     runtime/inbox/identidad existentes
          ↓                                            ↓
normalizar + deduplicar + privacidad → contexto privado del contacto
          ↓
selección/revisión → conocimiento aprobado del negocio + casos de evaluación
                                                        ↓
                                      respuesta con datos actuales y tools autorizadas
                                                        ↓
                                     resultado verificado → revisión → nueva versión
```

Tres usos diferentes, nunca un único almacén indiscriminado:

1. **Memoria privada de la conversación:** sólo para ese contacto, negocio y propósito;
   acceso y retención propios. No se publica en la biblioteca portable.
2. **Conocimiento del negocio:** FAQ, objeciones, estilo y procedimientos revisados,
   desidentificados cuando corresponda. Una respuesta humana histórica puede estar
   equivocada u obsoleta: no se promueve automáticamente como verdad.
3. **Evaluaciones:** conversaciones/casos reservados para comparar versiones, con
   separación por hilo/contacto y tiempo para evitar contaminación y fuga de respuestas.
   No usar a la vez los mismos casos como ejemplos recuperables y prueba independiente.

**Invariante crítica propuesta para el importador:** cargar historia nunca llama
Handle ni reencola sus mensajes como live, nunca envía respuestas, cobra, reserva
o invoca tools del negocio. Replay de prueba sin efectos separado de delivery real.
Este requisito está registrado; todavía falta implementarlo y probarlo.

Integración sin duplicación: conservar IDs oficiales scoped por negocio/cuenta/canal,
dirección/autor humano o bot, tiempo del mensaje y de importación, fuente y versión,
cursor/rango/paginación, recuentos, gaps, duplicados, hashes/conflictos y mapeo autorizado
al owner de contacto. No deduplicar sólo por texto o timestamp ni fusionar identidades
por nombres. Una colisión de ID con payload distinto se reconcilia, no se sobrescribe.
Empalme historia/live exige watermark y prueba de solapamientos/huecos/reinicio.
No inventar atribución de anuncio: sin identificador o enlace autorizado queda desconocida.

## 3. Menor costo: hipótesis que se debe medir

- Empezar con procesamiento local/incremental del historial autorizado, una sola vez
  por revisión y sin reenviar el corpus entero por respuesta. Local también consume
  cómputo/almacenamiento; no significa costo total cero.
- Evaluar primero reglas/respuestas aprobadas y recuperación acotada sobre infraestructura
  existente. No exigir embeddings, nueva base vectorial o GPU sin demostrar beneficio.
- Contexto realtime: datos vigentes, breve historial autorizado y conocimiento relevante;
  precios, stock y disponibilidad proceden del sistema actual, no de chats antiguos.
- Seleccionar modelo y volumen por calidad/costo/latencia medidos; no por marca o mínimo
  precio por token. Incluir canal/BSP, almacenamiento, transcripción, operación y revisión.
- Batch es candidato para clasificación/evaluación offline, nunca para responder un chat
  urgente. No ejecutar llamadas pagas sin presupuesto y autoridad.
- Fine-tuning no es requisito inicial ni consecuencia automática de importar chats.
  Promoción versionada y reversible; errores nuevos alimentan casos de regresión revisados,
  no entrenamiento online sin control.

Métricas requeridas por piloto: cobertura/gaps de importación, duplicados/conflictos,
fuga de PII/tenant, factualidad y cumplimiento, selección/argumentos de tools, derivación
humana, latencia y costo por conversación resuelta. Cita/reserva/venta sólo cuentan con
su evento empresarial confirmado; no inferir causalidad de conversión desde una charla.

## 4. Fuentes oficiales verificadas y cambios de vigencia

- [OpenAI — Evaluation best practices](https://developers.openai.com/api/docs/guides/evaluation-best-practices): admite datos históricos, producción y revisión experta para construir evaluaciones; no vuelve correcto todo chat ni garantiza aprendizaje.
- [OpenAI — Model optimization](https://developers.openai.com/api/docs/guides/model-optimization): evaluar primero y mejorar contexto/instrucciones; fine-tuning es una opción según caso, no paso obligatorio.
- [OpenAI — Batch](https://developers.openai.com/api/docs/guides/batch): documenta descuento del 50% frente a síncrono y ventana de 24 h para trabajos compatibles sin respuesta inmediata. No es una cotización de costo total, universal a cualquier proveedor ni autorización de gasto.
- [OpenAI — Deprecations](https://developers.openai.com/api/docs/deprecations): Evals pasa a read-only el 2026-10-31 y API/dashboard tienen cierre anunciado 2026-11-30; self-serve fine-tuning restringe organizaciones nuevas desde 2026-05-07. No basar una nueva integración en disponibilidad asumida de esas superficies. La metodología de evaluación no obliga a usar la API Evals.
- [Presidio — origen Microsoft](https://microsoft.github.io/presidio/getting_started/) redirige a [mantenedor actual](https://presidio.dataprivacystack.org/getting_started/). La biblioteca ya conserva Data Privacy Stack como owner, no se reetiqueta como publicación Microsoft actual. Revalidar modelos/idioma, licencia/SCA y falsos negativos antes de cualquier dato real.

No se copió código externo ni se adoptó una versión nueva en V288. OpenAI Docs
influyó en separar evals/contexto/fine-tuning y en no elegir servicios deprecados.

## 5. Acceso histórico: investigación sin afirmaciones inventadas

Se intentaron las páginas oficiales Meta onboarding-business-app-users y Conversations
API; el lector devolvió 429/no recuperable. FAIL-405 conserva esta limitación. No se
verificó aún la elegibilidad de coexistencia, ventana histórica, campos, permisos,
licencia de exportación o disponibilidad por cuenta. Ninguna cifra de retención ni
capacidad de obtener todos los chats se admite desde snippets/blogs.

Pregunta enviada al usuario: dónde viven hoy los chats (WhatsApp Business teléfono,
API/CRM, Instagram, Messenger, correo u otro). No se solicitaron chats privados ni
credenciales. Después: cuenta propietaria/acceso oficial, rango, volumen, idiomas,
adjuntos/audio, permisos de uso y retención/borrado, fuentes de resultados y presupuesto.
No asumir que poder leer un chat autoriza entrenamiento o compartirlo entre empresas.

Antes de programar el importador: resolver canal y fuente oficial mediante el expediente
CAPABILITY_GAP_RESOLUTION_GATE vigente; admitir código compatible con pin/licencia/tests,
verificar privacidad y extracción de exports bajo los contratos documentales cuando
apliquen. Mantener acquisition, corpus, revisión, retrieval, evals y delivery con owners
existentes; probar explícitamente cero envíos/efectos durante backfill antes de promover.

## 6. Estado de cierre

Investigación parcial y requisito conectado al roadmap. **No hay nuevo pack ejecutable**,
no se importaron datos reales, no se entrenó un modelo ni se midieron ahorros/ventas.
Readiness global sigue BLOCKED/43; inventario de código y perfil 67/742 sin cambios.
Los fallos previos 384–387 y los demás cierres no desaparecen por este requisito.
