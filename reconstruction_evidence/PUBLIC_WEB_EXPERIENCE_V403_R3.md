# FASE2 pública V403 — DONE técnico, aprobación visual pendiente

**P1–P6 implementado y verificado en el alcance local/fixtures autorizado.** Revisión `V403-PUBLIC-WEB-0.1.0-r3`, build r4. `phase2_authorized=true`. La aprobación de las 12 candidatas públicas sigue en **PENDING_USER_APPROVAL**; no se hereda del admin ni de las referencias Brave.

## DONE observable

| Criterio solicitado | Resultado | Evidencia |
|---|---|---|
| Composición pública navegable en destino nuevo | PASS | [reference](<<WORKSHOP>/Elite Library Extension V403/public-web/reference>): portada → catálogo → detalle → consulta confirmada/incierta → ubicación/turno; Next real y backend fixture |
| Admin sin regresión demostrada | PASS acotado | 9 pruebas privadas, 17 comparaciones con cero píxeles distintos; mismo layout y fuentes privadas. Reutilización por hashes del ensayo previo, no reejecución ficticia |
| Candidatas públicas desktop/móvil listadas y hasheadas | PASS de entrega | [Galería](<<WORKSHOP>/Elite Library Extension V403/public-web/PUBLIC_VISUAL_CANDIDATES.md>),12 candidatas; aceptación pendiente |
| Patrón Brave → regla nuestra → evidencia | PASS documental trazable | [15 patrones](<<WORKSHOP>/Elite Library Extension V403/public-web/docs/BRAVE_PATTERN_IMPLEMENTATION.md>) y anexo BRAVE-OBS-20260914-01 preservado |
| Autorización en el recibo de implementación | PASS | [PHASE2_IMPLEMENTATION_RECEIPT.json](<<WORKSHOP>/Elite Library Extension V403/public-web/evidence/PHASE2_IMPLEMENTATION_RECEIPT.json>) |

## Qué cambió y qué se probó

Nuevo overlay TS-PUBLIC-WEB-EXPERIENCE-V403: 22 targets públicos bajo una única autoridad visual V403. Portada centrada en producto, navegación con pictogramas y etiquetas, catálogo/detalle con precios y condiciones del owner, consulta mínima y siguiente paso contextual. La bicicleta se muestra completa en móvil; la portada sirve un JPEG de 250.602 bytes frente al PNG original de 1.923.305 bytes, que se conserva con SHA y procedencia. No se copió media, código, CSS, fuentes ni identidad de los sitios observados.

**17 verificaciones públicas PASS**: 15 pruebas del recorrido/estados/orígenes/404/teclado/reflow/contraste y 2 ejecuciones de cobertura (es-AR/en-US), con escenarios cero/uno de texto largo/varios modelos y tablet 768. **9 verificaciones privadas PASS** conservadas por identidad de las fuentes y el límite del cambio. No se suman capturas como si fueran pruebas de negocio.

La primera corrida pública encontró dos fallos reales: Origin extranjero aceptado por los POST de consultas/turnos y respuesta 200 para modelo inexistente bajo streaming. Se corrigieron sólo sus copias en el overlay nuevo: origen configurado obligatorio antes de efectos y Suspense del listado separado del detalle. El reensayo rechaza orígenes faltantes/null/ajenos, conserva consentimiento/idempotencia y devuelve 404. Fallos de preparación y fixtures se conservaron, corrigieron y reensayaron; no se etiquetaron como falta de cuenta.

Contraste de hero medido sobre el fondo raster real: mínimo 9,01:1 desktop y 9,09:1 móvil en los rectángulos examinados; los estados/formularios tienen comprobación de superficies sólidas. Rendimiento: observación local sin throttling, portada desktop LCP 200 ms y CLS 0 en esa navegación; **no constituye PASS de red móvil, percentiles o INP de campo**. Los métodos y límites exactos están en los receipts.

## Reconstrucción, notices y revisión

- [Pack nuevo](<<LIBRARY_ROOT>/implementation_packs/TYPESCRIPT_PUBLIC_WEB_EXPERIENCE_V403.md>), SHA256 `8ae8f4d714c074b2fd349d24f80cd5566c20456ed8b7e84f3b13029fce207324`.
- [Plan complementario T2804](<<LIBRARY_ROOT>/markdown_system/PUBLIC_WEB_EXPERIENCE_PACK_PLAN_V403.md>), SHA256 `9b0983448fc266dad02e305aa721ad46b466d385069e1b41ceb71efb26a3e6c5`.
- Manifest SHA256 `441fde82852f5d56d21ae6c6a6edf412343db731b409b8df5f766d3f94ad2ae8`; 22 targets y 226 fuentes protegidas. Dependencias/lock y tokens originales preservados; nuevo código de unión/presentación declarado AUTHORED, sin atribución a empresas por reputación.
- Materializador estándar reconstruye exactamente el pack y su fixture codificada. La aplicación fresca R2 probó 248/248 archivos idénticos al destino ejecutado. R3 conserva todos los bytes de producto/applier y amplía sólo el harness/cobertura; reutiliza ese apply/build, con recibo de identidad explícito.
- [Notices](<<WORKSHOP>/Elite Library Extension V403/public-web/pack-output-r3/bundle/THIRD_PARTY_NOTICES.md>), [provenance](<<WORKSHOP>/Elite Library Extension V403/public-web/pack-output-r3/bundle/PROVENANCE.json>) y receipts de exportación conservan origen y hashes de media. El normalizador Go de media no se ejecutó en esta fase; el catálogo usa un fixture PNG acotado, servido por hash.

## V402 intacto y límites

317/317 controles de preservación PASS. Execution 337 sigue COMPLETE; release original READY_FOR_LIBRARY_USE, errors=[]; conexión 18/18 y bridge original intactos. Ambos Get-FileHash coinciden:

- distribution336: `cd24757744d01dbc37b666495d659f76824a8a680015eee469187e975542591f`.
- signed-release artifact.zip: `dba7978d04d08dc2eedb112524708e220ba18cccf86c3a2142f8b1556ed1bf76`.

Árbol con ampliación pública ≠ bytes de ZIP/cierre original. Los ZIP siguen siendo las instantáneas canónicas del READY V402. Sólo se añaden pack, plan y este expediente; no hay push ni franquicia real.

Lector de pantalla, teléfono físico, zoom real, pruebas con personas, campo/INP y cloud: **NOT_RUN**. Aprobación visual pública: **PENDING_USER_APPROVAL**. Producción: **NO autorizada**. La recuperación incierta está probada dentro del mismo formulario montado; no se promete continuidad tras recarga. Los efectos del backend son fixtures en memoria; los datos comerciales/sedes/media reales pertenecen al futuro proyecto. No se afirma ausencia universal de bugs ni franquicia productiva terminada.

**Próxima acción: control de usuario/Grok sobre las 12 candidatas identificadas por hash. El agente se detiene aquí.**
