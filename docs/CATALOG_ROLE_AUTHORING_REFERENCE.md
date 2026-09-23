# Catálogo por rol — referencia de infraestructura

La página /admin/catalog permite crear un modelo con1–16variantes, precios
explícitos y vigencia UTC, cargar PNG, congelar una versión, revisar sus cuatro
etapas y publicarla. El formulario de versión compone un modelo con su libro;
el contrato backend admite hasta32modelos y128variantes en un snapshot. No se
afirma un editor masivo multi-modelo ni edición de borradores ya congelados.
Para cambiar una versión publicada se crea y revisa otro snapshot; restaurar
una versión aprobada crea una publicación nueva con un libro efectivo nuevo.

Seleccionar los packs de autoría y catálogo conectado junto con la composición
Go/PG/Next/BFF/identidad existente. La referencia completa incorpora todo.
features.catalog_editor=true muestra el portal; la sesión debe tener catalog:read.
catalog:draft crea fuentes/PNG/snapshot. catalog:review:legal, :technical, :media y
:publication permiten las revisiones correspondientes. catalog:publish publica;
el creador del snapshot no puede revisarlo ni publicarlo. Cada permiso sigue
limitado por tenant/organización real en Go; el nombre de rol no autoriza.

Aplicar migration0076 después de0075. El host del catálogo exige ese contrato
además de sus guards anteriores y perfil SHA fijado. La migración down rechaza
historial source; en base vacía el down/up funciona. Los modelos/variantes y
libros iniciales permanecen draft, modelos no públicos y homologación unknown.
La revisión humana no inventa acreditación fiscal ni homologación oficial.

El escritor de modelo conserva SQL/evento originales de Electromobility;
sus services y los de Commerce se componen en la misma transacción mediante
repositories acotados. La variante usa la tabla existente; release_command
conserva actor/identidad/hash/referencias y outbox. No segundo ledger.
La transacción completa se revierte si falla el outbox; cuatro solicitudes
simultáneas de la misma identidad producen un alta y tres replays.

El navegador conserva sólo referencias/hashes de recuperación y comprobantes
de fuente/imagen, ligados a tenant/org/subject. No guarda razones de revisión,
archivos de evidencia ni PNG en sessionStorage. Se verifica actor/command/hash/
snapshot antes de resolver una respuesta perdida; la recuperación sólo consulta.
Las huellas de source/review/publish/media se comparan entre Go y TypeScript,
incluidos caracteres no ASCII, escapes HTML y valores int64 como texto.
Este contrato estrecho no es un serializador JSON/JCS genérico.

Los JSON están limitados a32KiB y el PNG a1MiB, con plazo de lectura y cancelación
del reader. El helper textual conserva sus presupuestos y comportamiento UTF-8.
El BFF verifica sesión/permiso antes de consumir el cuerpo; una publicación
vigente se consulta con autorización privada, no desde otro tenant público.

Ejecutar TestCatalogSourceAuthoringAtomic con PAYMENT_CONNECTED_DB_URL aislada,
TestCatalogRoleCommandGolden y la prueba de navegador
TestCatalogAuthoringRoleBrowser con ELITE_CATALOG_ROLE_BROWSER=1,
ELITE_WEB_ROOT/ELITE_NODE_BIN absolutos y Next compilado --webpack.
La fixture sólo crea tenant/org; las16escrituras salen del navegador:
2modelos/2variantes/2snapshots/8reviews/3publicaciones, con rollback a la primera.
Los dos replays por pérdida de respuesta son GET y no agregan un POST.
Las5price_books finales son2de fuente y3efectivas, sin reactivar una antigua.
config/catalog.authoring.fixture.json y el goldens JSON son públicos sintéticos.

La revisión visual verificó escritorio y390px sin desborde. Los códigos de
estado/tributación/homologación reflejan sus owners; la localización privada
completa continúa bajo T2804. Tampoco se cierra supply/warranty/J5/CMS/KPIs ni
el mapping de marketplace T2805 por este recorrido.
Procedencia: AUTHORED glue; ninguna dependencia o atribución corporativa nueva.
