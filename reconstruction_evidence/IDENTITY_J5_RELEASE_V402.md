# V402 / checkpoint310 — J5 identidad y administración acotada

PROVEN_LOCAL para J5 en el perfil OIDC externo seleccionado. T2803 integral
continúa con source/SCA. Biblioteca204packs/2245bloques; perfil109/1481.
No READY global, vendor admin API ni producción.

El owner original de red y su frontend aceptan network:bootstrap junto con
network:admin para la organización inicial. No se necesita ni se agrega *.
Después del alta, los ámbitos de organización siguen siendo explícitos;
ninguna nueva organización otorga automáticamente acceso a su creador.
El IdP gobierna el alta de cuentas, MFA, asignación de roles y su retirada.
El contrato completo de esa frontera está en docs/J5_IDP_ACCESS_CONTRACT.md.

Prueba firmada RS256/JWKS→HTTP→PostgreSQL del owner original:2organizaciones,
3recibos inmutables con subject/hash/ámbito y3eventos outbox. Rechazados:
bootstrap sin admin, admin sin bootstrap, ámbito ajeno, permiso retirado,
token expirado/alterado, tenant ajeno y recuperación por otro actor.
El segundo issuer fixture queda dentro del owner de red; no añade dependencia
oculta con el pack de sesiones. El corpus/natural-key de cada tenant de prueba
queda aislado y la evidencia anterior se conserva.

Tres casos SDK originales→portal/Go/PG: refresh aplica permissions vacíos sin
inventar otro rol; rechaza cambiar subject; un usuario deshabilitado queda en
reauth_required y no puede terminar otro authorization code. Se enlaza el ciclo
completo PKCE/CAS/logout/revocación probado en309; no se lo repite sin delta.
Bearer ya emitido conserva validez hasta su expiración firmada. No se afirma
revocación inmediata universal ni implementación de un IdP/password stack.

4nuevos AUTHORED y3revisados;1474salidas anteriores intactas. Owners:
GO_NETWORK_ROLE_COMPOSITION0.2.0, TYPESCRIPT_NETWORK_ROLE_PORTAL0.2.0,
GO_OIDC_PORTAL_SESSION0.1.1 y TYPESCRIPT_OIDC_PORTAL_ADAPTER0.3.1.
Tipos/vet y4perfiles exactos PASS. Wire/canonicalización/SQL previo sin cambios;
no se inventa una corrida nueva de fuzzing. FAIL914/915/916 cerrados local.
Credenciales/configuración del IdP quedan para el usuario; source/SCA integral
sigue como trabajo de T2803, nunca oculto como falta de credenciales.
