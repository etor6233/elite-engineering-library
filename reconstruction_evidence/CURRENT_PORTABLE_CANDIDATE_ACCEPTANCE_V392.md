# V392 — copia portable actual comprobada

Mantenimiento de biblioteca, 2026-09-11. La copia del checkpoint247 se empaquetó
dos veces desde la misma fuente congelada mediante CREATE_PORTABLE_ARCHIVE.
Los ZIP son byte-idénticos; cada entrada se verificó contra el manifest y la
fuente. Después se probó la guía y el instalador extraídos de ese mismo ZIP.

## Resultado de uso y recuperación

- NEW y EXISTING: guía real,13comandos y8consultas de ayuda; ciclo de entrada55checks.
- Fallo parcial real de escritura provocado con un lock de Windows: rollback
  exacto, archivos previos y eventos intactos en ambos tipos de proyecto.
- Recuperación con instalador real, reanudación válida y siguiente checkpoint;
  el intento duplicado de append se rechaza. Los fixtures permanecen BLOCKED.
- Exclusión comprobada de expedientes locales, formularios/respuestas de esta
  tarea y aprobaciones; no se hereda READY_TO_BUILD.
- Verificador estructural ejecutado desde el archivo extraído; fuentes
  inmutables al terminar. No se modificó código de producto ni toolchains.

El harness AUTHORED V352 conserva sus aserciones: sólo cambian la identidad del
candidato148→247 y el rótulo V352→V392. Su original se verificó contra el SHA
publicado y se conservó junto al diff de procedencia. Se ejecutaron los comandos
reales; no se reemplazaron sus efectos con mocks. No se adquirieron dependencias.

## Alcance de la copia

El archivo es `library-candidate247-not-for-release.zip`, interno y sin firma.
Contiene las fuentes de la biblioteca en checkpoint247:165packs,
1561archivos materializables,839Markdown distribuibles y56perfiles.
El perfil de franquicia sigue68packs/797archivos; la referenciaHTTP69/805.
No contiene runtimes o credenciales de un proyecto consumidor. La materialización
y admisión de cada pack siguen sus condiciones, que este ZIP no elimina.

El gate de release firmado se materializó10/10 como owner disponible, sin
invocarlo como una aprobación: siguen pendientes seguridad/SCA/admisión nativa,
integración integral y release firmado final. Dos empaquetados idénticos no son
dos compilaciones del producto ni un PASS de TEST07. Tampoco prueban pérdida de
energía, instalación concurrente, UX humana independiente o producción.

TEST06/08 se revalidan sólo sobre este SHA concreto; sus criterios y los48status
se conservan. TEST02,TEST03,TEST07 siguen BLOCKED:45/48. La parte restringida en
V386 continúa diferida por instrucción del usuario. Ningún bypass o aceptación
de riesgo se infiere. Todavía no corresponde declarar la biblioteca100%cerrada.

La raíz incorpora después este reporte y su checkpoint; no es byte-idéntica al
snapshot247. Cualquier futuro archivo modificado requiere contraste y aceptación
propia. No reemplazar los receipts de este candidato con los de otro payload.

## Evidencia reproducible

Stage portable: `Temp/elite-v392-298c139d14c14c6e994d425d08fd3f9e`.

```json
{
  "status": "PASS",
  "archive_sha256": "4124b1300df0699c6a007e8412f2cbd66bd6200a6eb3c4cc015fabcdc02464b0",
  "manifest_sha256": "c6d2af7ba716e83660240a8c3c4c91e5767ba4e33ce788969fae8ef4ba26da94",
  "payload_files": 853,
  "guide_commands": 13,
  "help_checks": 8,
  "lifecycle_checks": 55,
  "sources_unchanged": true,
  "test06": "PASS_ON_EXACT_CANDIDATE",
  "test08": "PASS_ON_EXACT_CANDIDATE",
  "test07": "BLOCKED",
  "source_checkpoint": 247,
  "archive_entries": 854,
  "deterministic_archives": 2,
  "source_date_epoch": 946684800,
  "signed": false,
  "published": false,
  "build_elapsed_seconds": [
    89.605,
    85.653
  ],
  "cli_processes": 18,
  "root_code_changed": false,
  "final_release_approved": false
}
```

| Artefacto del stage | SHA-256 |
|---|---|
| accept_candidate392.py | `94b29d3a2afcd307c4d16c2baca87d15ffca89a4387cb850e057c1c420381363` |
| accept_candidate352-original.py | `4fb0f0fd07742d2a006fd28b76982f4d310ca08fd86a00165bda8a5b7961cf4a` |
| guide_journey_v327.py | `5e7b2c008a1235ffa274ff65ef068cbfdc8b8b66e4e68e809e01b12bf77227c1` |
| harness-provenance.json | `e0e3c54eede4f36d2dbcd06d0da0173be7d883342b0eed2c7e68f262edb87ae0` |
| baseline247.json | `b84a566392263c6dd263568ef78620fcd871ac7dfe48e7d58668f7b59ef2c680` |
| reproducibility.json | `2c822debf69feb9c60d01b09af00847aed75c1c26adce0206ee0f16b54f49737` |
| archive-verification.json | `9cba8fd00b0dfaf959406196dabf1d923e36ff2f54d877a9d3829f2136336c84` |
| acceptance.json | `d08d7df56cd8dccacd82b71eac597c725faaeddc05ca398c526297169297e6a7` |
| closure392.json | `c6f30404e0c0dd62f9c62fed6ffbe9af7576a33f1eb792cdde5aba7342dd1b1e` |
| acceptance-run.log | `84f56aad9d91468444f251c6fad6d945bc61605011bfbd928a6738aa4b0324d1` |
| library-candidate247-not-for-release.zip | `4124b1300df0699c6a007e8412f2cbd66bd6200a6eb3c4cc015fabcdc02464b0` |
| library-candidate247-not-for-release.zip.sha256 | `ff9ad03cb5c8a15a2590532112c291e3221ad28f7fe8bbe8c779f5063a26c228` |
