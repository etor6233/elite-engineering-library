# Red de franquicia por rol — infraestructura de referencia

/network es opt-in features.network_portal; NETWORK_ROLE_ENABLED=true activa
el host después de verificar migración0077,3triggers,PK e índice del receipt.
No lee cuentas/secretos al estar deshabilitado. UI requiere network:admin o
franchise:write; cada comando conserva el permiso original y el scope exacto.
Sólo el permiso global existente permite crear organizaciones raíz. No se
hereda autoridad del padre para operar una organización hija. El alta no
otorga permisos; el proveedor de identidad gobierna admins/sesiones/accesos.

Fulfillment.Service conserva validaciones/estados y cuatro writers PostgreSQL
conservan exactamente SQL/guards. Se extrae sólo su frontera transaccional.
El adapter permite esas cuatro operaciones y rechaza otras. Organización e ID
de acuerdo se preasignan para recuperar; IDs de eventos independientes aleatorios.
El receipt inmutable liga comando/actor/scope/hash y entidad/version/estado en
la misma transacción del writer y outbox. No segundo ledger de negocio.
GET histórico no sustituye lectura actual antes de una nueva transición.

La estructura activa permite acuerdos y sucursales, no producción/despliegue.
Términos, territorio y fechas son inputs explícitos; no se inventa contrato,
jurisdicción, regalía o criterio de habilitación comercial. Se conservan
no-superposición de territorio, jerarquía y bloqueos para cerrar padres con
hijos/acuerdos activos. Estado cerrado no afirma revocar sesiones del IdP.

BFF autoriza antes del body y limita32KiB/tiempo. El navegador persiste sólo
identidad del comando/hash/resultado esperado bajo tenant+subject, nunca nombres,
términos ni secretos. Una respuesta incierta se recupera GET, sin POST automático.
Si aún no puede confirmarse, conserva referencia para revisión autorizada.
La consulta de un resultado root admite scope vacío/omitido porque protectedGet
omite valores vacíos; el receipt sólo vuelve al actor con autoridad root original.
Importes no se usan; versiones int64 viajan como texto exacto.

TestNetworkRoleAtomic:1nuevo+3replay concurrentes, cuatro tablas revierten en
fallo tardío del receipt. Actor/hash/tenant/org/stale y cierre prematuro rechazados;
territorio superpuesto impedido;12receipts/12outbox finales, histórico/current
separados. TestNetworkRoleBrowser usa Next/BFF/Go/PG/Chromium con JWE/RS256/JWKS;
siembra sólo tenant.13POST crean/activan raíz, franquicia, acuerdo y sucursal,
suspenden/reactivan sucursal y cierran acuerdo/sucursal/franquicia en orden.
Respuestas de alta root y cierre sucursal perdidas y recuperadas tras reload
con0POST adicional. Scope limitado no hereda permisos. Capturas desktop/390px.
TestNetworkRoleHostGuards rechaza ausencia de cada uno de5guards; downgrade
vacío/reapply PASS y downgrade poblado rechaza perder evidencia. Fuzz2s.
10pruebas web y4goldens Go/TS; build/types exactos, sin nueva dependencia.

Pruebas Go requieren PAYMENT_CONNECTED_DB_URL owned/loopback con74migraciones.
Browser añade ELITE_NETWORK_ROLE_BROWSER=1, ELITE_WEB_ROOT/ELITE_NODE_BIN absolutos,
Next compilado--webpack y issuer sintético existente del catálogo. No IdP live.
El companion portal posee los goldens Go/TS; composición completa requerida
para pruebas conectadas, perfil web aislado sólo imports/tipos.

J5: frontend de estructura/acuerdo por rol y ResourceCreate/progreso/evaluación
humana existentes son piezas conectables. Identidad/admin/revocación/rotación,
composición final y readiness de despliegue siguen sus controles T2803/8/1;
este delta no los certifica. Nuevas guías/CMS/KPIs/i18n privados siguen T2804.
AUTHORED sólo glue declarado; ninguna atribución a empresas famosas.
