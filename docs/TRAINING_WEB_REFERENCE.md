# Portal de capacitación — misma release

Seleccionar TS-GO-API-WEB-BRIDGE, TS-OIDC-PORTAL-ADAPTER,
TS-FRANCHISE-JOURNEY-PORTALS y este pack. El perfil completo incorpora el módulo
Go de capacitación. Configurar features.training_portal=true explícitamente.
La navegación sólo muestra /guide/training a training:learn o training:review;
las comprobaciones del servidor siguen siendo obligatorias.

El BFF valida sesión/permisos antes de leer el cuerpo. Conserva el límite temporal
del boundedCommandBody existente y selecciona un tamaño acotado para respuestas.
Sólo permite acciones tipadas y rutas Go fijas. No toma URLs arbitrarias del cliente.
Los formularios de lectura, respuesta y evaluación escriben en los owners durables.
Ante pérdida de respuesta, consultar la referencia guardada; no repetir el POST.
sessionStorage guarda sólo identidad del intento/curso/hash, nunca respuestas.
Las guías son públicas. Una evaluación permanece sujeta a revisión humana.

Con las dependencias fijadas, reconstruir el bundle de ayuda a un destino ausente:
node tools/export-training-content.mjs src/platform/help/content.ts <ruta-absoluta-output.json>
Comparar content_sha256 y source_sha256 con el perfil Go antes de activarlo.
Un cambio de guías produce una nueva revisión y evidencia de contenido.
Para compilar con node_modules externo enlazado: next build --webpack.
La receta existente evita la restricción de rutas del builder Turbopack.

El fixture config/training.browser.fixture.json sólo activa la UI en pruebas;
no configura IdP ni concede permisos. La prueba real usa sesiones JWE y tokens
RS256/JWKS sintéticos, alumno/revisor/lector/otro alumno/otra organización.
Un resultado de capacitación no concede permisos ni acredita aptitud laboral.
