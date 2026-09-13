# V402 — catálogo J3 conectado, candidato antes de publicación

PROVEN_LOCAL para los claims ejecutados abajo; J3/T2802 y readiness global no
se cierran todavía. Canónico87packs/1166archivos. Código candidato separado y
delta congelado exacto en CATALOG_CONNECTED_PROGRESS_V402.json.

Tres snapshots inmutables, cuatro reviews humanos compartidos por snapshot,
cinco publicaciones con rollback. Cada publicación crea un libro efectivo nuevo
desde los importes/vigencia aprobados y activa una sola vez usando los escritores
Commerce existentes: conserva eventos created1/activated2 y referencias históricas.
La primera implementación intentaba reactivar el mismo libro: FAIL845 rojo
conservado y corregido, sin debilitar la unicidad ni inventar versión upstream.

PASS en PG18.6 y Go1.26.8: cambio del modelo fuente no altera el snapshot;
storefront existente/precio/búsqueda coinciden; ocho llamados concurrentes
producen1efecto/7replays; fallo de outbox revierte13tablas; rollback vencido
rechazado; repositorio reabierto recupera receipt. HTTP real: respuesta perdida
de publicación recuperada por GET sin otro POST aceptado, tres ETags obsoletos
invalidados y storefront/feed con bytes idénticos. PNG estructural decodificado/
re-encode por stdlib Go fijada, límites de bytes/píxeles; no antivirus ni libxml2.

Feed de referencia conectado al OutboundDelivery existente: cuatro publicaciones,
rechazo terminal, aceptación con respuesta perdida, bloqueo de nuevo POST y
reconciliación GET aceptada. Intención inmutable y fence/eventos en PostgreSQL;
source/payload/profile SHA enlazados. Es contrato explícito de receptor fixture,
no equivale a Google/Meta/ML live. Su mapping a SDKs pertenece a T2805.

Host básico12guards/hash/profile PASS y down73 poblado rechazado con3snapshots/
5publicaciones/12reviews intactos; vacío down/up PASS. El harness reusó la base
propia detenida sin repetir el journey. Carga de perfil fijado PASS. FAIL846:
writeJSON privado pisaba la caché pública a no-store; writer dedicado probado
rojo/verde para200/304. Las APIs privadas mantienen no-store.

Pendiente antes de incorporar J3: activar host con feed, down74 y orden73/74,
fuzz finito, ruta de detalle/canonical SEO en Next existente, empaquetado exacto,
notices/locks y conciliación del gap map. La URL propuesta /models/{code} aún
no existe en Next; no se presenta como página indexable hasta materializarla.

Todo código nuevo es glue AUTHORED; owners admitidos y dependencias exactas
conservan su provenance. No producción, cuentas live, ARCA ni Daybreak.
