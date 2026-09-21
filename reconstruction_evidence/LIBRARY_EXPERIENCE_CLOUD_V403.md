# V403 — experiencia, funciones y composición de biblioteca

Revisión de código: UI0.3.0 + funciones0.1.0 + cloud0.1.0; cada snapshot se distingue por los SHA completos de este expediente. Alcance: mantenimiento/infraestructura y fixtures. No se abrió una franquicia real.

## Tres veredictos

**A) Preparada para iniciar la construcción de una franquicia: SÍ**, como biblioteca y composición local reutilizable. El proyecto debe ejecutar su intake/gates, elegir su alcance completo y aportar las decisiones del negocio. La experiencia visual nueva es candidata y requiere aprobación explícita; no representa todas las pantallas de un producto terminado.

**B) Ejecución cloud demostrada: NO.** Se ejecutaron Windows y Linux/WSL locales, simuladores y reconstrucción. Eso no demuestra una cuenta, CI hospedado, despliegue, OIDC/IAP real, facturación, carga ni recuperación cloud.

**C) Producción autorizada: NO.** Se conserva production_authorized=false y el gate de producto DISCOVERY/BLOCK.

## Qué se puede reutilizar y cómo

[START_FRANCHISE](../START_FRANCHISE.md) → complemento del bridge original → [protocolo único](../markdown_system/FRANCHISE_PROJECT_OPERATING_PROTOCOL.md) → owners del consumidor → [perfil completo120](../markdown_system/FRANCHISE_COMPLETE_EXTENSION_PACK_PLAN_V403.md) → compositor original → overlay UI → overlay IAM → verificación y próximo checkpoint. Biblioteca y proyecto se fijan por separado; sus registros no se duplican.

El perfil conserva exactamente los116packs del V402 y agrega4extensiones. El destino siguiente es un **ensayo de mantenimiento**, no una franquicia real: `<WORKSHOP>\Elite Library Extension V403\integrated-reference-final`. La verificación de bytes cubrió 1825 archivos después de 58 destinos de overlay disjuntos. [evidence/integrated-composition.json](<WORKSHOP>/Elite Library Extension V403/evidence/integrated-composition.json).

La UI pública y operativa comparte tokens, componentes, navegación oscura configurable por marca y menú móvil. Los contactos comienzan por un resumen; las acciones y controles se abren según la tarea. Los selectores mantienen sus validaciones, autorización y recuperación. IDs/versiones quedan internos donde existe una resolución autorizada; no se inventa información comercial ausente.

El lector de serie probado es teclado/manual: preserva caracteres y ceros iniciales, consulta datos autorizados y separa captura de confirmación. Enter no confirma la entrega. Cámara/decoder y hardware físico no están demostrados. La UI documental agrega recepción, lectura, corrección de4campos, revisión por otra identidad y decisión sobre el hash exacto; JPEG/PDF2MiB/una página/supplier-invoice. Los tests nuevos no acreditan OCR infalible ni amplían el corpus del owner.

Las11funciones/22tareas son contratos, con permisos propuestos o procedentes del owner identificado. No son11roles de autenticación ni11operaciones de negocio implementadas. El binder vincula6owners existentes, resetea claims y no concede permisos. El catálogo técnico opt-in sigue mostrando una muestra identificada; no es un tablero del estado real del consumidor. Galaxy permanece pendiente para sus materiales concretos, sin bloquear los métodos oficiales ya referenciados.

## Verificaciones y límites

| Criterio | Ejecución y resultado | Límite |
|---|---|---|
| Composición120 | PASS; 1789 archivos del compositor original + overlays verificados | No compila ni despliega por sí misma |
| UI nueva | Materialización65/65; overlay54/54; build; 27 criterios de navegador cubiertos | 26PASS iniciales + corrección focal del fallo de foco; histórico retenido |
| Documentos | 9/9 pruebas de navegador, real Next/BFF, backend simulado | No acredita proveedor, base de datos ni corpus real |
| Accesibilidad automática | 10 vistas/anchos, 0 violaciones de axe4.13.0 | No certifica WCAG completo, lector de pantalla o aceptación humana |
| Lighthouse13.4.1 | 5 ejecuciones: medianas rendimiento98, accesibilidad100, buenas prácticas100, SEO63 | Gate original completo **FAIL**, app local deliberadamente noindex; no se relajó |
| Contratos de funciones | 33tests; dos reconstrucciones; 11funciones/22tareas verificadas | Scaffold, no ejecución comercial |
| Cloud local | Python: Windows114PASS+2SKIP, Linux116PASS; Node8/8 por plataforma; certificados Linux7PASS | SKIPWindows de symlink ejecutados en Linux; proveedor real NOT_RUN |
| FinOps | Consulta acotada, diario con CAS, intención durable, Job/Scheduler y dispatch simulados | Publicar no acredita entrega; alerta no es restricción; no hay techo absoluto de gasto |
| Continuidad | Agente nuevo sin conversación, entry→protocolo→estado/eventos→9/9evidencia | PASS tras reparación documentada de ruta; adapters nativos NOT_TESTED |
| Visual | 6 comparaciones automáticas sin diferencias | **PENDING_USER_APPROVAL**; UI0.1/UI0.2 siguen rechazadas |
| V402/337 | 317 comprobaciones PASS | ZIPs y fuentes protegidas inmutables |

La repetibilidad visual no sustituye aprobación del usuario. [evidence/VISUAL_CANDIDATES_V403.md](<WORKSHOP>/Elite Library Extension V403/evidence/VISUAL_CANDIDATES_V403.md>). Zoom real del navegador, otros motores, lector de pantalla, teléfono/lector físico y tareas con personas conservan NOT_RUN/NOT_TESTED; no se infiere su PASS de cambiar el viewport.

La integración detectó un pin obsoleto del engine compartido UI/IAM. Se preservó el rechazo antes de escritura y se readmitió el engine actual con10/10casos de compatibilidad; sólo cambió el manifiesto IAM y los otros62payloads cloud conservaron sus bytes. No se inventó un diff contra el engine histórico cuyos bytes no estaban disponibles. La reconstrucción integrada usa el SHA nuevo indicado abajo.

## Procedencia por claim

| Claim | Fuente y método | SHA256 de pack | Condición |
|---|---|---|
| Contratos de función | AUTHORED glue; métodos oficiales referenciados;9packs históricos enlazados; [BUSINESS_FUNCTION_OPERATING_V403.md](../implementation_packs/BUSINESS_FUNCTION_OPERATING_V403.md) | `dd7203436fd4ee9c6fa4dbfd47db0ca6a9696a7aff43af8541589ec68209319b` | Vincular owners y demostrar cada operación del consumidor |
| Sistema visual | AUTHORED; React/Next admitidos existentes; fuentes del sistema e iconos locales declarados; [TS_DESIGN_SYSTEM_V403.md](../implementation_packs/TS_DESIGN_SYSTEM_V403.md) | `cc9d089574aa4564782845b1ed7c5a3b455fea1e8bca30d3a44b10447bdb24d1` | Aprobación visual y aceptación de target |
| Interacciones/BFF/documentos | AUTHORED/ADAPTED sobre owners existentes Go/API y lector acotado; sin nueva dependencia; [TYPESCRIPT_FRANCHISE_EXPERIENCE_V403.md](../implementation_packs/TYPESCRIPT_FRANCHISE_EXPERIENCE_V403.md) | `0cc1269a61d44ecd55d660371c99568ab5550ebeb60f8e1bd30e0e0f787ad724` | Corpus/proveedor/dispositivos/métodos no ejecutados conservan su gate |
| Cloud/identidad/FinOps | AUTHORED orquestación; comandos/docs oficiales y herramientas/imágenes fijadas candidatas; [FRANCHISE_CLOUD_EXECUTION_V403.md](../implementation_packs/FRANCHISE_CLOUD_EXECUTION_V403.md) | `ae998d28ee6fe19253fda4902269ac69356fd75baba3effa9539fced99a7bb5f` | G0–G8 de target, builds OCI y ejecución externa aún requeridos |

La observación de [xAI Console](https://console.x.ai/) orientó jerarquía y navegación; no se copió ni se atribuyó su código o recursos. Las reglas de foco se contrastaron con [WAI-ARIA Dialog](https://www.w3.org/WAI/ARIA/apg/patterns/dialog-modal/). React Aria/Spectrum, Material Symbols y decoder de cámara se evaluaron como candidatos: no se incorporaron ni admitieron por reputación. No se atribuye código local a xAI, Tesla, SpaceX, Google, Microsoft o Adobe.

Los locks/package/notices V402 siguen vinculados a sus bytes originales. Las extensiones declaran LicenseRef-Workspace-Owner para glue; no introducen fuentes/iconos externos. Las licencias/SCA/runtime de herramientas e imágenes nuevas siguen en el expediente cloud/ADMISSION.md: fijar un digest no es admitirlo.

## Cloud: composición completa, promoción aún no ejecutada

Se definen50pasos/8manifiestos y5recetas OCI: web, API/background, migraciones, ARCA y controlador FinOps; datos persistentes, identidad de aplicación y workload separadas, herramientas, CI y recuperación. El proceso API con background necesita CPU siempre/min1 hasta un drenado explícito. No se promete desarrollo siempre gratuito ni suspensión sin costo de datos/backups/logs/APIs.

El build invoca los owners originales; el empaquetado consume exclusivamente sus salidas verificadas y la promoción reutiliza ese digest. El artefacto Windows firmado V402 no se renombra ni se presenta como Linux. Full Linux product build/OCI/SCA/runtime, CI hospedado, IAP+OIDC, migraciones/restore, shutdown/grace period y aceptación siguen NOT_RUN. Esas son condiciones técnicas explícitas, no sólo credenciales.

FinOps distingue presupuesto interno, facturación observada, alerta y restricción efectiva. Incluye total compartido por cuenta/período, consulta limitada, persistencia con CAS e intención previa a dispatch; requiere export real completo, límite aprobado de consulta, identidades y comprobación de entrega/dedupe. Suspender requiere plan de drenado/aprobación separado y conserva datos.

## Git, artefactos y preservación

Conexión original preparada en commit local `d2729ab31adcb82906413e0f5b9f9b228c851d0d` (5archivos), sin push. La extensión agrega fuentes versionadas; no se reescribe historia pública. Commit fuente probado: `b58b30d7fa3269ad1bc43552270c3c832db8fde9`.

Recibos Git120 independientes incorporados: 2. Los hashes y métodos exactos están en el índice de evidencia; una reconstrucción Git local no equivale a disponibilidad del commit en GitHub ni a ejecución hospedada.

El árbol post-pulido/ampliación no equivale a los bytes del ZIP de cierre; **el ZIP sigue siendo la instantánea canónica del READY original**. El checkpoint337 conserva COMPLETE y library release READY_FOR_LIBRARY_USE/errors=0. Los48controles originales no se reabrieron.

| Artefacto inmutable | SHA256 |
|---|---|
| elite-library-v402-source330-meta331-distribution336.zip | `cd24757744d01dbc37b666495d659f76824a8a680015eee469187e975542591f` |
| signed-release-v402-r330/artifact.zip | `dba7978d04d08dc2eedb112524708e220ba18cccf86c3a2142f8b1556ed1bf76` |

## Condiciones externas — cuando el usuario decida

- [ ] Aprobar o corregir las6candidatas visuales; completar aceptación humana y métodos de accesibilidad/dispositivos pendientes.
- [ ] Autorizar push del commit elegido, proyecto cloud y presupuesto; preparar identidad/WIF/IAP/OIDC y referencias de secretos en el proyecto consumidor.
- [ ] Cerrar admisión de herramientas/imágenes/dependencias target, build/OCI, CI hospedado, deploy/promoción/rollback, restore/carga y seguridad del target.
- [ ] Conectar proveedor de pago, WhatsApp y demás cuentas seleccionadas; comprobar callbacks/reconciliación real.
- [ ] ARCA: la infraestructura/fixtures continúa disponible; cuentas/certificados y homologación real cuando el usuario lo autorice. No se pidieron datos.
- [ ] Daybreak/libxml2: ACCESS_BLOCKED diferido, expediente/trigger existente. No se investigó ni se usó para bloquear este delta.

## Índice verificable de evidencia

Los paths Desktop son la ubicación del ensayo de este equipo; no son defaults del producto. El repo lleva fuentes reconstruibles y este índice; los recibos externos contienen sus hashes, método y salidas. Un SHA demuestra correspondencia de bytes, no veracidad editorial ni certificación.

| Archivo del ensayo | SHA256 |
|---|---|
| [evidence/integrated-composition.json](<WORKSHOP>/Elite Library Extension V403/evidence/integrated-composition.json) | `aa5a5889d66f266d79b77afd9a832bfabcf81a2efe95aed7d8d2881f41f37e3a` |
| [ui/evidence/UI_FINAL_V403_0.3.0.json](<WORKSHOP>/Elite Library Extension V403/ui/evidence/UI_FINAL_V403_0.3.0.json>) | `a036a4a40dd3cb5ea1ff35a2b535d4cf2c1d899558c3d23d10bebffb6176d6d1` |
| [ui/evidence/browser-aggregate-final-0.3.0.json](<WORKSHOP>/Elite Library Extension V403/ui/evidence/browser-aggregate-final-0.3.0.json>) | `52d3eb0fe56285ab7d0ca70e2598320f9908a0a00a1bee9f0a1ee61af1f286fa` |
| [evidence/documents-browser-final.json](<WORKSHOP>/Elite Library Extension V403/evidence/documents-browser-final.json>) | `01eea4f9f63408a3b726978c844c5ee3775ddd8006bbd451f4f534bda95a2150` |
| [evidence/audit-ui-quality/axe-v03-final.json](<WORKSHOP>/Elite Library Extension V403/evidence/audit-ui-quality/axe-v03-final.json>) | `33161cc7c17c8a734522587b23ffa21ceaa8df1736b97c432e55cfc17a20e13a` |
| [evidence/audit-ui-quality/lighthouse-ui03-summary.json](<WORKSHOP>/Elite Library Extension V403/evidence/audit-ui-quality/lighthouse-ui03-summary.json>) | `38f8f06aecb89cf14606213479f871e7ef3e1f2f9f43ad1eed41483eeaafbfea` |
| [roles/evidence/validation-final.json](<WORKSHOP>/Elite Library Extension V403/roles/evidence/validation-final.json>) | `e969d3959bf5557a0cbaba618a9ec0b917306932a44f6fc2b43263e53c4c22b0` |
| [roles/evidence/standard-reconstruction-result.json](<WORKSHOP>/Elite Library Extension V403/roles/evidence/standard-reconstruction-result.json>) | `b052640482329e12667ee61d49dc9e421db28f643cdb4dd8ecec10cd38d13692` |
| [cloud/evidence/cloud-extension-final.json](<WORKSHOP>/Elite Library Extension V403/cloud/evidence/cloud-extension-final.json>) | `bc5e5cdce1625278fe72ab38b1f8fb9b1e34626e243e182c404c851fb30afdaa` |
| [cloud/evidence/qualification-windows-final-03.json](<WORKSHOP>/Elite Library Extension V403/cloud/evidence/qualification-windows-final-03.json>) | `72ac467c83793fc337f5e7e3d85cf7ad5a2d52bbe8b50aee6628c99ffdd1fd1a` |
| [cloud/evidence/qualification-linux-final-03.json](<WORKSHOP>/Elite Library Extension V403/cloud/evidence/qualification-linux-final-03.json>) | `4388545f73fc605a258189a305b4a6b8f623fb991daf48febbc5cba1ef3fe058` |
| [evidence/fresh-file-only-resume-v403.json](<WORKSHOP>/Elite Library Extension V403/evidence/fresh-file-only-resume-v403.json>) | `cde02f4f35f691f684ff16055ecfa99b87370d776e802725e0daefaf22fb2820` |
| [evidence/visual-candidates-v403.json](<WORKSHOP>/Elite Library Extension V403/evidence/visual-candidates-v403.json>) | `580e1c7c0ae0063af041b1a62a4e384ed2d49a5805e887aed1fd3d98adeb9897` |
| [evidence/preservation-final.json](<WORKSHOP>/Elite Library Extension V403/evidence/preservation-final.json>) | `5c5e7b0435e09512f6e42779e5b64b7cdfbacd9be96a09e54989ca161aa920f1` |
| [evidence/connection-impact-v403.json](<WORKSHOP>/Elite Library Extension V403/evidence/connection-impact-v403.json>) | `2c6eaef9d5a25ced9cad051b88c3773bd11b89ed4beaffeda34191a7c7f4477a` |
| [evidence/extension-source-commit.json](<WORKSHOP>/Elite Library Extension V403/evidence/extension-source-commit.json>) | `5811e91803cd4b99c8488ce96f9dc0ed6115295cddb7c0253d0e707c5c1fbdec` |
| [cloud/evidence/cloud-extension-iam-delta.json](<WORKSHOP>/Elite Library Extension V403/cloud/evidence/cloud-extension-iam-delta.json>) | `a1ddfa7b23cb68c783b57c1ca44ad1407e6c948392bc9e9f631baddfc60923f3` |
| [cloud/evidence/overlay-compat-01/receipt.json](<WORKSHOP>/Elite Library Extension V403/cloud/evidence/overlay-compat-01/receipt.json>) | `f5685a1e7b5f0d8d214d02b7272a8e178236502027b898680dfb3a7958ff6f6d` |
| [cloud/evidence/consumer120-windows-01/qualification-receipt.json](<WORKSHOP>/Elite Library Extension V403/cloud/evidence/consumer120-windows-01/qualification-receipt.json>) | `d3ba6fed04e7f7ad3e1b16c284e44e5caffe9fada5e2a95d3b0d7832cb0e4d14` |
| [cloud/evidence/consumer120-linux-01.json](<WORKSHOP>/Elite Library Extension V403/cloud/evidence/consumer120-linux-01.json>) | `9b3984f73b28fabda72b1d51740cc422438ad39766bda875ac5916be204d2eb1` |
