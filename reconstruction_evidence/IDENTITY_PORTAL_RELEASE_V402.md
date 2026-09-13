# V402 / checkpoint309 — identidad de portal y sesiones durables

PROVEN_LOCAL para este claim de T2803. Perfil completo109packs/1477archivos;
biblioteca204packs/2241bloques. T2803 integral sigue abierto por administración
J5 y admisión source/SCA. No READY global ni producción.

Código reutilizado: openid-client6.8.5/jose6.2.10/zod4.4.3 y los owners Go
OIDC/service-token/pgx ya fijados.26archivos nuevos AUTHORED y6revisados;
1445salidas anteriores intactas. No nueva dependencia ni autoría empresarial.
Go application1.22.0 conserva todos los módulos posteriores a la rama guardada;
OIDC portal0.3.0 conserva su modo compatible y añade lifecycle explícito;
GO-OIDC-PORTAL-SESSION0.1.0 aporta persistencia y mantenimiento.

La prueba conectada usa el SDK original, issuer sintético firmado, Go HTTP y
PostgreSQL:7casos de login/PKCE/refresh/logout/recuperación/retención y16writers
CAS→1owner;12rechazos de ámbito. El navegador recibe un handle JWE aleatorio;
la bóveda liga sessionID/profileSHA y PostgreSQL ciphertextSHA. El reloj se
consulta después del lock; una respuesta de refresh incierta no autoriza repetir
el grant. Logout queda durable antes de revocar en el proveedor. Ciphertexts
intercambiados entre dos sesiones del mismo perfil se rechazan.

Mantenimiento sin cookies: hasta2revocaciones y100filas de retención porbatch,
con conteo durable de revocaciones no confirmadas. Host probado con grant real
del SDK Go, JWKS firmado, un poll autenticado y shutdown conjunto.13fronteras
de respuesta; tipos/perfil/protocolo/sesión18tests PASS también en web separado.
Rollback con evidencia poblada se rechaza; base vacía down/up PASS. Perfil Go,
vet/build y fuzz2s PASS. No equivalencia a DAST/SAST o revisión de acceso IdP.

Cinco perfiles afectados reconstruidos exactos. Se reusó sólo el receipt de
backend577files tras demostrar plan/pack/output idénticos; la corrección del
fixture TypeScript sólo requirió reconstruir los otros cuatro perfiles. FAIL910
setup local,911entrypoint,912rollback/URL/media y913fixtureowner cerrados local.

Cada usuario sigue gobernado por su IdP: bootstrap admin, roles, MFA y baja
requieren el contrato J5 separado. Estos adaptadores no inventan un principal
ni administran passwords. La aceptación de tokens ya emitidos queda acotada
por su expiración; no se afirma revocación inmediata de bearer en todos los IdP.
Credenciales y configuración live serán del usuario; no se solicitaron ahora.
