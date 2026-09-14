# V402 — reparación de integridad tras publicación Git

Alcance: revisión del commit público b0ef1ad6230dd8d12e4a4aaa457f173c0c8d9e9f (448 archivos), sin reabrir T2801–T2810, ARCA o Daybreak y sin modificar los ZIPs inmutables. No se promueve producción.

## Observación antes de reparar

Los dos ZIPs y su firma permanecían intactos. El checkout tenía 39 diferencias frente a los 300 paths protegidos del cierre 337:36 por conversión uniforme de saltos de línea; README y gap map con contenido fusionado distinto; un registro de fuentes con saltos mixtos normalizados. Al comparar el payload completo del ZIP, 88 archivos diferían sólo en saltos y 3 documentos tenían deltas de contenido. Los packs conservaban su contenido; el problema era integridad/transporte y documentación.

La lectura completa de los objetos Git distinguió el problema local del publicado: 28 blobs del payload distributivo diferían del ZIP, 25 sólo por saltos de línea y 3 documentos por contenido. Además, el registro canónico de fuentes externas, fuera de ese payload, perdió sus saltos mixtos. Algunos archivos estaban correctos en disco gracias a core.autocrlf, pero tenían otro SHA en Git: la comprobación inicial del índice detectó esa diferencia antes del commit. Se recuperaron los bytes exactos en el índice y se exige un checkout independiente. El detalle queda en all-payload-git-blob-audit.json e index-transport-findings.json dentro del expediente local de revisión.

## Corrección aplicada

- Restaurar 104 archivos locales desde distribution336 o mediante recuperación que reproduce exactamente el SHA preexistente. Muchos ya tenían el blob Git correcto; restaurar el checkout no implica 104 cambios de código.
- Recuperar el registro mixto desde el respaldo 319 y la secuencia histórica de append, aceptándolo sólo al reproducir SHA `0bea06bf54700179447fb394892799cc0e5e0c023be3fd025b55f5a1da632722`.
- Conservar README SHA `f17bc0cfc40216cd6c98c2d70450daa5edbde007bdfb67c9bc95bcad4b2220ab` y gap map SHA `c72117518385d3a11fb54d1a5ef38e106a65a046481250abcee4d423e859f171`: el checkpoint 337 los referencia. Cifras vigentes y banners se aclaran en [el expediente no sellado](LIBRARY_VS_PRODUCT_GATE_V402.md), sin alterar el cierre.
- Declarar `* -text` en `.gitattributes` para impedir conversiones de bytes en index/checkout. Git sigue mostrando diffs de texto; no se desactiva ningún gate. [Semántica oficial de Git](https://git-scm.com/docs/gitattributes).
- Actualizar START_FRANCHISE y el expediente post-publicación para que no sigan diciendo que el código V402 está pendiente de subir. Inventario actual: 205 packs / 2360 archivos materializables; perfil 116 / 1653. Los conteos históricos del README no gobiernan el plan JSON.

**árbol post-pulido ≠ bytes del ZIP de cierre; el ZIP sigue siendo la instantánea canónica del READY original**. La metadata documental posterior y este expediente son un delta de publicación; fuente 330 / metadatos 331 de producto y ejecución 337 no cambian.

## Comprobaciones exigidas antes del push

Los 300 SHA protegidos deben coincidir con el baseline; state 337 COMPLETE y cadena de337 eventos válidos; release READY_FOR_LIBRARY_USE sin errores. Los 116 packs seleccionados y el compositor deben coincidir con la instantánea; la composición debe producir los 1653 SHA del inventario aceptado. Un checkout Git nuevo, incluso con core.autocrlf=true, debe conservar esos bytes. Ambos ZIPs deben mantener sus SHA y ningún archivo de signed-release puede cambiar. Las evidencias de ejecución de esta revisión quedan fuera del payload, en `qualification/grok-publish-review` de la referencia durable; este texto define las comprobaciones y no sustituye sus resultados.

No se repiten builds de producto, pruebas de negocio ni suites live por esta reparación de transporte. Las modificaciones de b0ef1ad permanecen en la historia; la corrección se publica en un commit sucesor con push normal, nunca force.
