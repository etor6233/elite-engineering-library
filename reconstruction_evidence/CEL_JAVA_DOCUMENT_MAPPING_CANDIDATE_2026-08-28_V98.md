# CEL-Java document mapping pack — evidencia V98

## Autoridad exacta

- proyecto oficial CEL, originado y publicado por Google y trasladado a `cel-expr/cel-java`: release `v0.13.0`, commit `c57d9d02259795c761e82f8870303b93a7a672c8`, 2026-05-15, Apache-2.0;
- tag lightweight y commit consultado por API: commit unsigned. La identidad no se presenta como firma verificada;
- archive exacto: 1.892.506 bytes, SHA-256 `088d90e7e178fe09fd664d5dcb5c60a7d946fbe9d27cadb7400f9868e69b2cc8`;
- Maven Central `dev.cel:cel:0.13.0`: JAR 499.701 bytes/SHA-256 `a3bc22d78affe749f81b17cce2131a3807275436157ba4d6fe3f10d8f541ea94`; POM 4.259 bytes/SHA-256 `3a6eccb1ce1aececfec7e90f56073e87f1bc14491b6d00104bf84beeedaed0f6`.

## Pruebas ejecutadas

- prototipo Java 21/Maven 3.9.16 con CEL 0.13.0: mapping cerrado a claves aprobadas, identidad/version/hash enlazados al resultado, rechazo de hash/output extra/input ausente/multilínea, macros deshabilitadas y expresión mayor al límite; 6/6 PASS;
- compilador usa type-check previo, `Map<string,dyn>` de entrada/salida, sin funciones host, macros vacías, límite de 4.096 code points y profundidad de parse 32;
- CycloneDX 2.9.1 generó SBOM 1.6 con 18 componentes; Google OSV Scanner 2.5.1 exacto terminó con cero findings conocidos;
- Bazelisk oficial 1.29.0, SHA-256 `092a8738d5b41aae7a85c42cc961b1034e3389aba43ffc20c0fabda7b43e095b`, adquirió Bazel 9.2.0 firmado por Bazel Developer;
- once suites upstream focales de common/compiler/runtime/bundle ejecutaron 2.836 pruebas: 2.831 PASS y cinco FAIL, cero skips. Diez suites pasaron; los cinco fallos pertenecen a `ErrorsTest` en Windows por semántica LF frente a `System.lineSeparator()`;
- el comando upstream total `bazelisk test ...` no fue admitido: analiza targets Android y este host carece de Android SDK.

## Decisión

`REBUILD_VERIFIED / CONDITIONED` para `GOOGLE-CEL-DOCUMENT-MAPPING` 0.1.0. No se copió source Google: los siete archivos materializables son integración `AUTHORED` y dependen del artefacto oficial fijado. El pack Markdown mide 24.284 bytes y su SHA-256 es `9f1e4a49f7e1351c181db5fb2705dd50f02603b5802d04e071c03b443195f935`.

La reconstrucción desde Markdown produjo 7/7 archivos byte-exactos. Ejecutada desde un directorio ajeno, compiló con Java 21/Maven 3.9.16, pasó 7/7 pruebas propias, regeneró CycloneDX 1.6 con 18 componentes y OSV Scanner 2.5.1 exacto reportó cero findings. `VERIFY_LIBRARY_PASS` cerró 71 packs/680 archivos/425 Markdown y ambos perfiles actualizados: email 15/145 y persistence focal 6/61. `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit` cerró con `google_cel_document_mappings=1`.

El template sigue bloqueado: ninguna expresión, schema ni corpus viene aprobado. Sólo admite expresiones single-line, macros vacías, input/output JSON acotado, claves de salida exactas e identidad/version/SHA del mapping. Antes de persistir, el proyecto debe fijar y probar schemas de entrada/salida, corpus, determinismo, límites de carga, rollback y unión transaccional con el consumer. Los cinco fallos upstream Windows permanecen visibles; el glue no se atribuye a Google.

Memoria vigente: 1.091 fallos locales + 185 condiciones upstream = 1.276 IDs únicos.
