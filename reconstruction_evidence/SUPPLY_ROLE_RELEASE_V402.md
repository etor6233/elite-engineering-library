# V402 / execution297 — suministro conectado por rol

PROVEN_LOCAL para el recorrido T2804/J2: alta de orden/cantidades desde formulario,
confirmación y producción, series/calidad/rechazo/reemplazo, despacho, recepción
parcial en cuarentena, rechazo/reinspección y stock disponible.95packs/1264archivos
exactos;191packs/2030blocks en biblioteca. No cierre de todoT2804 ni READYglobal.

El alta compone Operations.CreatePurchaseOrder con BindPlan en una transacción,
conservando SQL y reglas originales, planned receipt y outbox existentes. Cuatro
solicitudes concurrentes crean una orden y tres replays. Variante inválida y
fallo tardío del outbox revierten cinco tablas, incluida la orden inicial.
La recuperación liga actor/orden/comando/hash; int64 viaja como texto exacto.

La prueba de navegador real Next/BFF/Go/PG usa JWE/RS256/JWKS y sólo siembra
maestros tenant/org/supplier/catalog. Sus23POST dejan1orden recibida,3series
(1rechazada y reemplazada),2stocks disponibles,6solicitudes/6decisiones de calidad.
Se pierden respuestas de alta y recepción; tras recarga se recupera porGET,
sin unPOST adicional. Rechazo de calidad y reinspección usan evidencia nueva;
la separación humana sigue gobernada por el owner y el UI muestra sus límites.
Duración del test conectado9.44s; escritorio/390px sin desborde inspeccionados.

4pruebas BFF comprueban permisos antes del body, límite32KiB/cancelación,
separación de vistas/org y recuperación de actor.6checks de contrato y4goldens
Go/TS verifican Unicode/HTML/cantidades/int64/versiones. Next--webpack y tipos
PASS; fuzz finito2s/3semillas/99596ejecuciones. No migración ni dependencia nueva.
Los guards del host J2 y la máquina de estados no cambiaron.

FAIL856 corrigió goldens con separadores Unicode al final, inválidos portrim
en ambos lenguajes. FAIL857 corrigió un brace de JSX, con tsc antes del rebuild.
FAIL858 fue un nav.supply sin traducción detectado en revisión visual: dos
diccionarios corregidos y compilados, sin repetir la transacción del navegador.
Todos los RED y PASS están fijados en SUPPLY_ROLE_RELEASE_V402.json.

Cuatro perfiles con delta se reconstruyen exactos:95/1264franquicia,12/200web,
61/833serverless y96/1272HTTP. Backend sin delta conserva su proof anterior.
Los perfiles menores compilan el cierre de imports/tipos.16archivos nuevos y
7revisados son AUTHORED glue declarado; ninguna atribución empresarial nueva.
Packs nuevos GO-CONNECTED-SUPPLY-CREATION0.1.0(6files) y
TS-SERIAL-SUPPLY-PORTAL0.1.0(10files); tres owners revisados.

Guía del destino: docs/SUPPLY_ROLE_REFERENCE.md. Ayuda inline en la misma release.
Siguen warranty/J5/CMS/guías/KPIs/i18n privado bajoT2804, luego el orden acordado.
TEST02/03/07 siguen abiertos; no se usa el antiguo45/48 como porcentaje.
