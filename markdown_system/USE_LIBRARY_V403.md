# Usar la biblioteca sin perder el hilo

**Un punto de entrada para personas y agentes.** La biblioteca es una base de ingeniería; cada proyecto tiene negocio, datos y aceptación propios. El [recibo de la revisión](../reconstruction_evidence/UNIFIED_REFERENCE_V403_R4.md) separa reconstrucción, ejecución local y pendientes externos. Ningún inventario equivale a producción.

## Empezar o retomar

1. **Ubicar el trabajo.** En un proyecto existente, abrir `PROJECT_AGENT_ENTRY.md` y su estado/eventos: revisión, alcance y próxima acción exacta. En uno nuevo, seguir [START_FRANCHISE](../START_FRANCHISE.md) y el [protocolo operativo único](FRANCHISE_PROJECT_OPERATING_PROTOCOL.md), mediante el bridge. No copiar el estado de mantenimiento de esta biblioteca como estado del negocio.
2. **Seleccionar todo el alcance.** El blueprint clasifica las [48 superficies](TOTAL_SYSTEM_CAPABILITY_CONTRACT.md) y los recorridos completos. REQUIRED necesita owner admitido y criterio observable; OPTIONAL no se presenta como terminado. El plan completo se ejecuta en incrementos, no en módulos sin conectar.
3. **Buscar antes de leer.** Una fila del [índice por claim](PACK_PER_CLAIM_INDEX.md), una autoridad y un pack. XERJ sólo si una consulta devuelve pasaje + archivo/línea; si falla, `rg -n -F -- "término" <owner>`. Leer alrededor del hit. No cargar el corpus, reindexar por cada tarea ni prometer ahorro40×: no está medido.
4. **Construir y comprobar.** Elegir la fuente y versión compatibles, ejecutar el incremento conectado y sus pruebas afectadas. Reutilizar pruebas previas únicamente si código, entradas, dependencias y condiciones siguen iguales. Un hash prueba identidad, no funcionamiento.
5. **Dejar continuidad.** Registrar avances materiales, fallos y recuperación en los registros existentes: blueprint/authority map/pack plan, estado/eventos, failure lessons y evidencias. Una próxima acción concreta. Un integrador escribe estado y manifiesto; colaboradores entregan cambios aislados.

Grok, Cursor, Codex y Claude pueden seguir la misma entrada. Sus adapters son configuración; ejecución de un vendor se acredita por un recibo, no por nombrarlo. Para ver el avance diario, leer el resumen/next_action del estado y abrir la referencia de la misma revisión; no mantener otro tablero como segunda verdad.

## Encontrar capacidades

| Necesidad | Autoridad de entrada | Qué debe demostrar el consumidor |
|---|---|---|
| Arquitectura y capacidades | [Sistemas](../SYSTEMS_ENGINEERING_MASTER_MAP.md), [catálogo](CAPABILITY_CATALOG.md) | Selección coherente, sin duplicar owners |
| Público, operación y permisos | [Frontend y UX](../FRONTEND_PRODUCT_ENGINEERING_UX.md), [perfil integrado](UNIFIED_REFERENCE_PACK_PLAN_V403_R4.md) | UI→owner→efecto durable; PC y viewport móvil; aceptación humana aparte |
| Franquicias, compras, stock, venta, entrega y postventa | [Acelerador](../FRANCHISE_ACCELERATOR.md), [índice por claim](PACK_PER_CLAIM_INDEX.md) | Quote/order/payment/handover/fiscal son efectos distintos; B2B/B2C no elimina reglas por marca |
| Documentos y captura | [Perfil documental](OFFICIAL_DOCUMENT_INTELLIGENCE_PROFILE.md) | 21 clases, recepción/revisión/commit/recuperación; QR/checksum no es firma; hardware y OCR real se califican aparte |
| Funciones comerciales/Galaxy | [Selección V403](FRANCHISE_COMPLETE_EXTENSION_PACK_PLAN_V403.md) | 11 funciones con tarea/permiso/owner/evidencia; no 11 roles de auth ni playbooks de evento inventados |
| IA, formación, contenido y red | [Mapa IA](../AI_ENGINEERING_MASTER_MAP.md), [índice por claim](PACK_PER_CLAIM_INDEX.md) | Presupuesto, permisos, herramientas, evals y recuperación, con proveedor fixture/live identificado |
| Nube, seguridad, operación y mantenimiento | [Sistemas](../SYSTEMS_ENGINEERING_MASTER_MAP.md), [dependencias](DEPENDENCY_UPDATE_CONTRACT.md) | Artefacto probado, CI/entorno/rollback, backup/restore, alertas; gasto interno ≠ factura o techo universal |
| Fuentes y licencias | [Admisión](../PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD.md), [notices sucesor](THIRD_PARTY_NOTICES_V403.md) | Versión/SHA/licencia/gates; AUTHORED no se atribuye a empresas |

Móvil nativo, desktop nativo, IoT y GPU son selecciones que necesitan admisión propia; no son sinónimos de una web responsive. POS/payroll y otros helpers no se convierten en aplicaciones completas por listar sus nombres. Este índice no excluye capacidades del proyecto futuro.

## Reconstruir y verificar

[START_V403_LOCAL](START_V403_LOCAL.md) contiene los comandos exactos para materializar el carrier, reconstruir las fuentes, preparar herramientas fijadas, instalar el bridge y ejecutar la calificación fixture. Elegir destinos nuevos. El calificador se detiene al terminar: no es un servidor permanente ni un proyecto real.

Verificación de la biblioteca ampliada desde su raíz:

```powershell
pwsh -File ./markdown_system/VERIFY_LIBRARY_V403.ps1 -LibraryRoot .
```

El verificador original permanece intacto y su inventario es histórico. El sucesor admite de forma estricta el skill `.agents`, comprueba todos los checks históricos y añade el pack/plan/notices/inventario de esta revisión. El comando no sustituye el recibo runtime ni autoriza producción.

## Qué conservar

Biblioteca versionada y fuentes/locks/notices; paquete materializado; evidencia; datos y configuración de cada proyecto. **No borrar reference/workspace/backup como si fueran cachés.** Bases, uploads y recibos irremplazables requieren copia y procedimiento de recuperación. Sólo salidas declaradas regenerables pueden limpiarse.

El árbol post-ampliación no coincide con los bytes del ZIP de cierre V402: ese ZIP sigue siendo la instantánea canónica original. Live, cloud alojada, teléfono físico, lector de pantalla y visual humana conservan estados propios. Daybreak/libxml2 diferido; no se pide ninguna credencial aquí. REVESTEX y X no se inician con esta entrega.
