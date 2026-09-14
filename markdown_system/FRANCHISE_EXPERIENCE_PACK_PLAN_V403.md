# Experiencia reutilizable V403 — owner T2804

Revisión UI: **V403-UI-0.3.0**. Extensión de mantenimiento, local/fixtures. V402/337, sus archivos protegidos y ZIPs permanecen intactos. No es una franquicia real ni autorización de producción.

## Selección e instalación

La selección completa vive en [FRANCHISE_COMPLETE_EXTENSION_PACK_PLAN_V403.md](FRANCHISE_COMPLETE_EXTENSION_PACK_PLAN_V403.md). Incorpora juntos `TS-DESIGN-SYSTEM-V403 0.3.0` y `TS-FRANCHISE-EXPERIENCE-V403 0.3.0`; no seleccionar además el runtime visual anterior. El único CSS/tokens activo se aplica al frontend existente mediante el overlay con hashes previos exactos. Los contratos de negocio siguen en sus owners.

1. Materializar el perfil completo en destino separado; nunca dentro de la biblioteca ni de V402.
2. Ejecutar `python experience_overlay/apply_overlay.py --target . --report qualification/ui-overlay.json` y repetir con `--verify`.
3. Instalar dependencias con el lock existente, compilar TypeScript y construir `next build --webpack`. No se agregan versiones nuevas de React/Next ni runtimes visuales candidatos.
4. Configurar referencias revisadas en `config/experience-options.json` y listas publicadas en `config/experience-checklists.json`. Ambos se entregan vacíos: **BUSINESS_CONFIGURATION_REQUIRED**, no «falta credencial» ni defecto de código encubierto. Cada tenant mantiene nombres de sucursales en `organization_labels`; la sesión determina qué organización se puede consultar.
5. Para ensayo local, ejecutar el harness incorporado: `python experience_overlay/qualification/run_qualification.py --target . --evidence <directorio-local> --candidate`. Requiere Node/Python/OpenSSL, las dependencias y navegador fijados del pack Playwright, y puertos loopback4311/4313/4314 libres. Genera material efímero fuera de fuente; nunca incluir ese directorio completo en una release.
6. Sin `--candidate`, compara las imágenes existentes. Una coincidencia automática **no aprueba** un baseline. `--grep`/`--grep-invert` permiten repetir solamente el criterio afectado.

La instalación valida todos los archivos antes de mutar, usa lock exclusivo, reemplazo atómico por archivo y rollback de errores ordinarios. No ofrece transacción atómica de toda la carpeta ante un corte de energía. Una diferencia del consumidor exige revisión de compatibilidad, nunca sobrescritura silenciosa.

## Dos familias, una autoridad visual

- Web pública: los owners `/models`, captación y navegación pública conservan sus contratos. El ensayo reutiliza el modelo sintético Urban de `src/platform/backend/public-client.test.ts`; no agrega imágenes comerciales ni propiedades inventadas.
- Aplicación: `/experience` es una entrada breve; `/experience/operations` ejercita la entrega con fixtures opt-in. `/franchise` aplica el mismo patrón de tarea enfocada a pedidos, entregas, agenda y gestión. Los panes visitados permanecen montados para conservar campos y fences; los no visitados se cargan al abrirlos. La navegación usa sidebar en escritorio y diálogo modal nativo en móvil. Las áreas vienen de los permisos del owner; el menú conserva el foco y los paneles visitados. El selector contextual muestra una tarea a la vez. Los contactos empiezan plegados y sólo uno se abre a la vez; cada formulario conserva su fence. Los fallos accionables se muestran donde ocurren y la ayuda queda plegada.
- Referencia de construcción: `/experience/reference` conserva el catálogo técnico, estados, paletas y once funciones. No aparece como navegación principal de la aplicación. Las funciones son especificaciones de mantenimiento y no permisos, operaciones probadas ni once roles. El consumidor vincula sus propios owners mediante el binder del [contrato C](BUSINESS_FUNCTION_OPERATING_CONTRACT_V403.md); nunca copia estado, evidencia o permisos del mantenimiento como si fueran suyos.

Los campos de referencia se seleccionan por nombre; IDs, versiones y huellas siguen internos. Se conservan los números de serie/VIN realmente observados y las entradas legítimas de configuración. El inventario inicial de135 controles está evaluado en el expediente;27 controles adicionales requeridos se conectaron al selector compartido, más los recorridos guiados de cotización, asignación, disponibilidad y entrega. El GET del selector filtra tenant/organización/sujeto/permiso. En los cuatro recorridos citados, el POST repite el catálogo revisado; en las otras superficies la autoridad del efecto sigue siendo el BFF/owner de dominio original. El catálogo UX no concede permisos de negocio.

Los textos nuevos de operación preservan es/en con el provider existente. No se traducen automáticamente datos comerciales, nombres aprobados ni prompts de listas publicadas. El consumidor debe aportar esos contenidos según los idiomas de su blueprint.

## Procedencia y admisión

| Fuente | Revisión | Uso y estado |
|---|---|---|
| React/Next/TypeScript/Playwright del perfil V402 | package.json + pnpm-lock.yaml inalterados | DEPENDENCY_PIN heredado, compatibilidad local verificada |
| Glue de componentes, selección, composición y pruebas | Manifiesto SHA por archivo V403-UI-0.3.0 | AUTHORED declarado; sin atribución a empresas externas |
| [Adobe React Spectrum/Aria](https://github.com/adobe/react-spectrum/tree/4693fcc844a341107e7dd26fe456edb1e269e44c) |4693fcc844a341107e7dd26fe456edb1e269e44c| Evaluado, licencia observada; **NOT_SELECTED / NOT_ADMITTED**. G3–G7 no ejecutados; no incorpora runtime |
| [Google Material Symbols](https://github.com/google/material-design-icons/tree/40a7a292a79d9394157e1ea24f83d52d5e17c556)|40a7a292a79d9394157e1ea24f83d52d5e17c556| Evaluado, licencia observada; **NOT_SELECTED / NOT_ADMITTED**. No iconos adquiridos |
| [W3C contraste mínimo](https://www.w3.org/WAI/WCAG22/Understanding/contrast-minimum.html)| Método citado, no código incorporado | Cálculo sRGB AUTHORED; el danger anterior3,953:1 queda como regresión histórica |
| [Playwright comparaciones](https://playwright.dev/docs/test-snapshots)| Runtime fijado del perfil | Método real de comparación; aprobación humana separada |

El drawer sigue el patrón oficial W3C de [diálogo modal](https://www.w3.org/WAI/ARIA/apg/patterns/dialog-modal/): contenido externo inerte, foco dentro, Escape y retorno al disparador. El mínimo AA de [target size](https://www.w3.org/WAI/WCAG22/Understanding/target-size-minimum.html) es24CSSpx con excepciones;44px para acciones móviles principales es una decisión UX de este ejemplo.

No se copió código, marca, imágenes ni assets de xAI, Tesla, SpaceX, DogeOS o Grok Build. Los candidatos no se admiten por estar nombrados ni por licencia o reputación.

## Verificación y aceptación

El recibo UI enumera SHA exactos del manifiesto, herramientas, pruebas, capturas y límites. Evidencia durable de esta ejecución: `C:/Users/NL/Desktop/Elite Library Extension V403/ui/evidence/`. Los recibos y fuentes portables se incluyen selectivamente; las claves efímeras y trazas con cookies de fixture se excluyen.

La candidata0.1 fue **REJECTED_BY_USER** y sus imágenes están conservadas con hashes en `rejected-visual-20260914`. La candidata0.2 también fue REJECTED_BY_USER, conservada en rejected-visual-0.2.0. La candidata0.3 oscura requiere aprobación explícita del usuario. El contraste automatizado mide los elementos/renderings realmente ensayados; no certifica todos los futuros temas y contenidos.

| Criterio | Estado permitido |
|---|---|
| Contratos, selección, composición y runtime local | PASS sólo con recibo de esta revisión |
| Comparación de imágenes candidatas | PASS automatizado separado de aceptación |
| Baselines visuales de usuario | PENDING_USER_APPROVAL |
| Lector de pantalla, teléfono físico, zoom real y prueba con personas | NOT_RUN; viewport equivalente no sustituye dispositivo o zoom |
| Operaciones de negocio de los once temas | NOT_TESTED por el catálogo; leer cada claim C |
| Cloud/producto real y producción | NOT_RUN aquí / no autorizada |

Los13 digests HTML de ayuda V379 ya fallaban sobre el renderer V402 exacto por la evolución i18n. Se conserva el expediente; el test adaptado compara la salida por defecto V403 byte a byte contra la copia fijada del renderer V402 y mantiene los enlaces/versiones y el contenido de capacitación. No se cambió ningún SHA publicado del cierre.

## Captura literal de unidades

AuthorizedUnitCodePicker admite lector en modo teclado y escritura/pegado manual en un campo con foco. Enter sólo busca dentro de un pedido completo ya consultado por el BFF autorizado; no usa listener global, no interpreta URLs y conserva ceros, mayúsculas y los128bytes UTF-8 del contrato supply. Vista parcial, cero coincidencias, ambigüedad o estado no seleccionable no eligen unidad. Elegir esta unidad cambia una selección local sin enviar un POST; repetir el código no aumenta cantidades. La operación existente mantiene autorización, versión, evidencia y reconciliación en su owner. La búsqueda no es un catálogo universal ni demuestra hardware físico.

Cámara/decoder queda NOT_ADMITTED; no existe botón de cámara sin función. La fuente de datos del pedido sólo aporta algunas referencias técnicas: no se inventan nombres de productos.

La observación de serie en la entrega también intercepta Enter del lector: conserva el texto literal y exige confirmar la revisión por separado. No inventa una búsqueda de inventario donde el owner de handover no la expone.

## Facturas y evidencia de esta revisión

El nuevo `DocumentWorkspaceV403` conecta el owner existente de documentos: recibir original, procesar, revisar cuatro campos y decidir con otra persona autorizada; la decisión del owner realiza el commit. Requiere el transporte compartido de D y la configuración de tenant, organización, perfil y modo. El alcance exacto y sus condiciones viven en `docs/DOCUMENT_EXPERIENCE_V403.md`, incluido en el overlay. No introduce OCR, reglas fiscales, bandeja universal ni búsqueda documental.

Manifiesto final: `e3bd6c00bc1c2b909370519b44065f301022eda471eaef68a12de5ffa8199303`.

| Comprobación | Resultado ejecutado |
|---|---|
| Materializador original | 65/65 archivos idénticos a los bloques |
| Overlay independiente | 54/54; instalación, repetición sin escrituras, conflicto sin cambios y rollback de fallo ordinario |
| Compilación | PASS; snapshot de 332 archivos compilados |
| Unitarios afectados | Dos ejecuciones: 73 PASS y 26 PASS; conjuntos parcialmente solapados, no sumarlos como pruebas únicas |
| Navegador autónomo | 27 criterios cubiertos: 26 PASS iniciales, un fallo de foco corregido; se repitieron sólo foco y comparación visual |
| Comparaciones visuales | 6 candidatas coinciden; aprobación humana PENDING_USER_APPROVAL |
| Documentos | 9 recorridos de navegador incluidos en los 27; backend HTTP simulado, BFF Next real |

El fallo de Tab no se borró: Chromium conservaba rectángulos de controles dentro de `details` cerrados; se excluyeron del orden tabulable y se verificó el ciclo completo. Sólo `application-chrome.tsx` cambió entre aquella ejecución y la corrección; el recibo conserva los hashes antes/después y el FAIL inicial. No se presenta como una única ejecución final de 27/27.

Evidencia: [recibo UI final](<C:/Users/NL/Desktop/Elite Library Extension V403/ui/evidence/UI_FINAL_V403_0.3.0.json>), [criterios de navegador](<C:/Users/NL/Desktop/Elite Library Extension V403/ui/evidence/browser-aggregate-final-0.3.0.json>), [materialización estándar](<C:/Users/NL/Desktop/Elite Library Extension V403/ui/evidence/standard-materialization-final-0.3.0.json>) y [snapshot compilado](<C:/Users/NL/Desktop/Elite Library Extension V403/ui/evidence/compiled-snapshot-final-0.3.0.json>). Las candidatas están en `ui/qualification/tests/visual-candidates/v03-*.png` del destino de extensión.

«Sin fallos abiertos» se limita al código, revisión, métodos y criterios enumerados en ese recibo. No equivale a aceptación humana, ejecución cloud, lector físico, cámara, lector de pantalla, nuevos idiomas, nuevos temas ni producción.
