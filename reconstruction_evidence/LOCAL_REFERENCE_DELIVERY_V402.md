# V402316 — T2808 entrega local de referencia

PROVEN_LOCAL para LIBRARY_INFRASTRUCTURE. No producción live.
Perfil115packs/1560archivos, biblioteca205packs/2294blocks:
1990AUTHORED159ADAPTED145VERBATIM. Tres archivos previos cambian (host,
Next config y template OCI);1527previos exactos. Doce bloques nuevos son glue,
tests y documentación AUTHORED;18archivos existentes pnpm/supervisor se
seleccionan con sus hashes y condiciones ya demostradas. Sin dependencia nueva.

Dos builds con instalación/caché de compilación/fuentes nuevas producen
4066archivos byte-idénticos. SHA del inventario del artefacto:
5ec38349f25c63e932d6cc22a4c4a22d5b07eb63e45e1559c1d74138d82a48ad.
La receta usa una ruta estable de build ausente en cada arranque: Next16.3.4
emplea rutas absolutas como identidad RSC. Se archiva el workspace anterior,
se fijan datos públicos de preview/Server Actions sólo para fixtures y se
normalizan localizadores generados/traces NFT declarados ADAPTED. No se
reescriben bundles de aplicación ni se promete igualdad en cualquier host/ruta.
La entrega se arrancó con el workspace de compilación ausente.

PostgreSQL real:84migraciones desde base vacía, ledger de hash/ordinal,
transacción por paso, lock de sesión; replay sin cambios y drift rechazado.
API+Next+OIDC fixture reales, rollback manual entre binarios, candidata fallida
con restauración de la anterior, recovery del controlador y replay de un único
pedido durable. Go adopta el evento de apagado ya existente en el worker;
Next conserva su cierre SIGTERM143, wrapper0 y Job Object con árbol vacío.

Cinco perfiles de fuentes exactos. Lint del delta14archivos/25avisos revisados,
0nuevos blockers/0sin triar. Sólo se reescaneó el builder tras ajustar la receta;
el resto conserva hashes. No se reinvestigó Daybreak/libxml2 ni se repitió SCA
del grafo sin cambios. Next16.3.5 observado: fixes publicados para superficies
no seleccionadas; se conserva16.3.4 y el expediente build-authority.json.

El artefacto conserva fuentes correspondientes, notices de owners, Node LICENSE
exacto y catálogos de71paquetes npm/39módulos Go. Seis paquetes npm sin texto
legal directo quedan explícitos en el JSON para resolver su licencia de
upstream/parent en T2810 antes de redistribuir el ejecutable. TEST07/firma y
destino durable siguen pendientes; no se presenta este artefacto como release
firmado ni como READY global. Contenedores/cloud permanecen fuera del target
local seleccionado, con condiciones de target y ejecución, no sólo credenciales.

FAIL938–942 contenidos/corregidos en el alcance demostrado. Se retienen el
copiado MAX_PATH rechazado, la proyección de enlaces rechazada, el test de lock,
el primer cierre SIGTERM mal interpretado y los dos builds no idénticos.
No se alteran sus resultados fallidos. Ver JSON para receipts y hashes.
Siguiente: T2809 operación local; ARCA penúltimo y dossier nativo último.
