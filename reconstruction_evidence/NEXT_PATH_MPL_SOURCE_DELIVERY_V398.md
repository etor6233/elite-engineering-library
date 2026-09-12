# V398 — fuente y licencia MPL de next-path entregadas y verificadas

2026-09-11. Mantenimiento de biblioteca desde checkpoint264 validado.
Cierre local: fuente original next-path1.0.0, manifiesto y licencia completa
entregados por las tres recetas pnpm actuales.9copias byte-idénticas,
12negativos reales y90tests PASS;8archivos reconstruidos desde Markdown.
El contrato conserva45/48 y TEST02/03/07 bloqueados.

## Autoridad fijada y alcance de la fuente

Commit oficial ce2c1386836339bc473b76c9888947899ead8d56:
https://github.com/sholladay/next-path/tree/ce2c1386836339bc473b76c9888947899ead8d56
Se adquirieron index.js, package.json y LICENSE desde raw.githubusercontent.com
con revisión fija, el árbol Git oficial y el registro npm next-path/1.0.0.
Cinco respuestas HTTP200 sin redirects, cuerpos con límites explícitos.
El árbol no está truncado y cada blob Git calculado coincide con el árbol.
El gitHead npm coincide con el commit; nombre/versión/licencia coinciden.
El manifiesto declara index.js como main y único elemento files.

| Original | Bytes | SHA256 |
|---|---:|---|
| index.js |317|ef96862a543004139a10f8f52f0a45b4c972d0c4aa26e85396fd7aa2575d290e|
| package.json |1010|07d01d2d44846a76138c89985a17eda58358db997b5a25cdc7e87163aadcfdf4|
| LICENSE |15978|786ef75c24eb986a2ffe31eb878b51a16826af1f5c9bd5cbdb5ad9a4223b5ab7|

No tarball npm adquirido, instalación ni ejecución de next-path.
No se afirma firma Git ni identidad byte a byte del tarball publicado.
La licencia completa coincide con la evidencia MPL anterior V336.

## Correspondencia estructural delimitada

Acorn8.14.0 previamente fijado, SHA256
758cead0e9764f94320f938ac169fb95c7eeea30dc675a6e5b1133fae239fa19
El receipt conserva la versión/hash realmente verificados del parser.
Se analiza el módulo completo en dist/pnpm.mjs, bytes7619794..7620528
(extremo final exclusivo),734bytes:
d8b58764b8a66de031a13ad8b76cccbb936e49f69c7df6d6bda306ced71fab83.
BundleSHA256:
ddc64218bc85fb88d28b5def06eb01fafb39cb67f7a9465ce736456658cc7f11.

Las4sentencias coinciden estructuralmente sólo después de:
- desempaquetar un wrapper __commonJS con sus parámetros observados;
- convertir sólo las dos declaraciones superiores nombradas var a const;
- mapear path246/path, nextPath2/nextPath, from5/from, diff2/diff,
  next2/next, module2/module y __require/require;
- quitar posiciones AST y Literal.raw, preservando semántica literal,
  propiedades, operadores, orden de argumentos, directiva y export.

__require(path) a require(path) es un adaptador condicional del loader:
la resolución real del host no está demostrada. No se afirma equivalencia
runtime ni build reproducible de todo pnpm. Diez controles mutados fallan:
builtin, orden relative, propiedad sep, comparador, inicio substring,
orden join, valor de retorno, export, directiva strict y binding superior.
El receipt detalla cada diferencia y mantiene false los cuatro claims de
ejecución/equivalencia/build/tarball. AST normalizado:
14befa9a0127d30d1d160aec651ab7ed32c7a6a669537c370a45f7cba13da2d7.

## Implementación canónica

Pack0.7.0/8files y recipev6. El envelope ADAPTED
next-path-source-evidence.json conserva base64 de los tres originales,
commit, URL, Git blob, SHA256, bytes y la comparación limitada.
Licencia del envelope: LicenseRef-Workspace-Owner AND MPL-2.0;
no cambia la autoría ni los términos de la fuente original.

NEXT-PATH-MPL-SOURCE.md entrega los tres textos completos y explica
los cambios de empaquetado y la limitación del loader. La fuente cubierta
se ofrece bajo MPL2.0 sin añadir restricciones a esos derechos; las demás
licencias de pnpm conservan sus ámbitos.
Límites: evidencia65536bytes, original32768, documento65536.
El planner verifica pin, schema, inventario exacto, scope, región del bundle,
base64 estricto, longitudes y hashes. verify exige el archivo, fija el receipt
y regenera el contenido esperado; una falsificación conjunta tampoco pasa.
La publicación mantiene destino ausente y limpieza ante fallos de escritura.

90tests:17selection+73routing. Ocho regresiones nuevas cubren entrega íntegra,
ausencia, cuerpo fuente mutado, licencia completa reemplazada por etiqueta,
receipt/documento rehasheados, región bundle alterada, pin de evidencia
alterado y fallo de escritura con limpieza. No nuevas dependencias runtime.

## Recorrido real y preservación

Composición fresca enterprise web6packs/126files.
Tres recetas: enterprise-web, Playwright, Lighthouse. Extracción independiente
por encabezado y longitud de los tres originales en cada receta:9copias
completas coinciden con SHA256 y bytes de origen.
Doce negativos reales: mutación de fuente, omisión MPL completa, documento
ausente y rehash conjunto de documento/receipt, en cada consumidor.
Restaurar los originales devuelve PASS.
Las tres recetas nuevas no ejecutan pnpm; executed/runtime_admitted/
redistribution_admitted permanecen false.

442archivos payload y su receipt intactos. Colección47textos V397 y suplementos
BlueOak/QRCode/semver byte-idénticos en los tres consumidores.
Selector, política,17tests selection y catálogo47 permanecen intactos.

## Fallos y límites abiertos

FAIL795: texto MPL conocido pero fuente original y transformación documentada
no llegaban a los consumidores. Cerrado en alcance de entrega local verificable.
FAIL719 recurrente: nombre de receipt supuesto y glob de directorio no soportado;
se inventariaron las rutas reales, se leyó ast-correspondence.json y no se
atribuyó resultado a lecturas fallidas. Se conservan fallos y before/after.

No cierra todo el bundle pnpm, OpenPGP LGPL/source/relinking, condiciones
Artistic, publicación/procedencia restante Yarn/pnpm ni admisión de runtime.
Daybreak/libxml2/V386 permanece diferido por instrucción del usuario.
TEST02 exige aún integración funcional material; TEST03 seguridad/source/
nativos y TEST07 release mantienen sus blockers. Los oráculos no se modifican.

Rollback: snapshots before conservan pack0.6.0, planes, contratos y mapas;
se conservan recetas V397 y payload original. No sustitución de runtime.

## Receipts

Stage: Temp/elite-v398-cb7a3831938478; pointer elite-v398-current.txt.
Evidencia local preservada; enlaces de autoridad y bloques canónicos portables.

| Evidencia relativa al stage | SHA256 |
|---|---|
| source-receipts.json | f7056a62c0d14240e14b73c6b6b84f2fe8263197ecc37ba817bd1497066ba932 |
| fixed-source-verification.json | e98a9d38122be69f30b8a7fe00f38dfa76aec60c8c891c9fb17ce3508e46a3ac |
| module-correspondence.json | 072d18dcf41463b15e258385178fba9ae456183a8d4fe0efc20244af76b5f73f |
| source-parity.json | 7885f240912abc41dae6ca7eead37635fc3647387da4ad2551bc4b244793c274 |
| rebuild.log | f461c5d66fc379dfc6a628e6cd4ba5daa9fdf697b8a26ec0b8b9a883063e01c2 |
| rebuilt-tests.log | ae8fc6df1c2ec1d88987fe704e192bccb5f74a0f65e136d69aafb5f5d5cd79f0 |
| current-web.log | 37bd03340746ee9fca065f9ecd44f299b2cb44403439f3e2fb4510b4cfc6a567 |
| real-delivery.json | f53273c64ccc9e4809409bd0ea7a3896653d355cb8f0dab79ddbc69e6aa4da21 |
| projection-before-after.json | 4b56bb3ee5a5bfa16b7ac0832b67866ff5d5348e7357be81e1e3d200eca9f27b |

Preflight265 pendiente tras checkpoint; aún no se atribuye resultado global a esta revisión.

## Cierre comprobado — checkpoint266

Preflight265:164pasos ejecutados PASS y56composiciones verificadas.
Inventario165packs/1564archivos/846Markdown;1307AUTHORED/150ADAPTED/107VERBATIM.
Disponibilidad BLOCKED sólo Docker ausente; no equivale a admisión productiva.
Franquicia68/797, HTTP69/805 y enterpriseweb6/126 no cambian.
El Preflight general ejecuta instalaciones offline congeladas ya existentes;
las tres recetas nuevas verificadas no ejecutan pnpm.

Las8fuentes canónicas coinciden con la reconstrucción probada. Selector,
política,17tests selection y catálogo47 permanecen idénticos. Los48oráculos,
acciones y estados del contrato siguen iguales:45passed; TEST02/03/07blocked.
Se cierra entrega local de fuente/licencia next-path, con sus límites de
correspondencia explícitos, y se preserva el resto de los blockers.

| Evidencia final relativa al stage | SHA256 |
|---|---|
| preflight265.json | 642a94b75041df226afc6fec16491957bfc895d83f5fdc82ff3d364661e9795a |
| preflight265-process.json | 155f5484850364caa55fd7dbefb6d741ce89f76fed7339dd9098a692d819fbaf |
| preflight265.log | f4ef86235efc5faea11db17ad3710430dcc33c4a6f987cd4204516ed4df47679 |
| delta-audit.json | e8075e40b4366140aced994269423ee8738a9e7fb9ed99229dfeb913213507a3 |
