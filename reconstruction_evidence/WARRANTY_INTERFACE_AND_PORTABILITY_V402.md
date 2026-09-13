# V402 — API, activación y portabilidad de garantía

Checkpoint288. Los cambios de garantía están completos en staging para el claim
local de backend/API y listos para incorporación canónica; el frontend por rol
continúa en T2804. No se cambia todavía84/1112 ni se declara T2802 terminado.

PASS:10casos HTTP de JSON/tamaño/identidad de ruta/scope y representación int64
exacta. El mismo J4 ejecuta cuatro casos por HTTP real, con permisos/cliente/fábrica
verificados en handlers y repository. Dos conexiones se cierran después del commit;
GET recupera recibos ligados a actor/payload, sin POST oculto. El contador incluye
el POST de fábrica deliberadamente rechazado. El verifier de este fixture entrega
principales explícitos: no acredita criptografía JWT ni navegador.

El host desactivado no lee perfiles. Activado verifica archivo acotado regular,
hash, tenant/org, esquema y guards antes de montar rutas; dos tests de host con
PostgreSQL PASS sin skips. El downgrade de tablas pobladas falla conservando4
reclamos; base vacía baja/sube0069/0070 y recupera sus7tablas. Ambos PG terminaron.

Lease: un test PostgreSQL focal espera el mínimo60segundos real dentro de outbox.
Una decisión tomada antes del vencimiento falla al COMMIT posterior:0decisiones,
0steps/eventos/consumos filtrados, reserva original/version1 y aprobación pending
conservadas. El rechazo posterior libera la reserva. Este test usa handover/cita
históricos explícitos como inputs; no repite ni acredita su creación completa.

El CLI usa el mismo loader Go del host, exige hash de origen/destino ausente y
valida overrides antes de escribir. Dos destinos contienen profile/activation
byte-idénticos, con ruta relativa y hashes; rechaza overwrite, hash o scope inválido.
Template sin secretos y guía operativa incluidos en staging.

GO-NATIVE-FUZZ-GATE exacto:2segundos/cuatro workers/103752ejecuciones PASS sobre
parser/calendario, copiados desde6archivos Go exactos con MIT. El scaffold conserva
module/versión Go y sólo esa closure sin imports externos; no representa toda la
franquicia ni reemplaza SCA/SAST. Los dos tests/seeds de baseline exigidos por el
gate no suman nueva cobertura. Vet de host/API/CLI PASS.

FAIL837 metadata de checkpoint,838 contador de fixture HTTP y839 fixture histórico
slot/ends_at corregidos con RED preservado. No se relajaron guards ni se cambió
el costeo para lograr el PASS. Nuevos archivos/funciones son AUTHORED glue; el
predicado BC ADAPTED conserva su fuente/licencia y admisión estrecha anterior.

Manifiesto/hash/delta congelado: WARRANTY_INTERFACE_AND_PORTABILITY_V402.json.
Siguiente acción única: publicar pack de fechas y pack de garantía conectada,
actualizar sus owners/planes/notices/locks y reconstruir exactamente. Reutilizar
estos resultados por identidad de fuente; no repetir suites sin delta. Después
reconciliar las filas restantes del blueprint para cerrar T2802 y seguir T2804.
