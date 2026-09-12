# Transparency.dev Tessera POSIX Evidence Log — 2026-08-27 V1

## Resultado

`TRANSPARENCY-DEV-TESSERA-POSIX-EVIDENCE-LOG` V0.1.0 quedó `REBUILD_VERIFIED / CONDITIONED`. Materializa 11 archivos: cinco archivos oficiales Apache-2.0 byte-verbatim y seis archivos de integración `AUTHORED`. El plan `TESSERA_POSIX_EVIDENCE_LOG_PACK_PLAN.md` compone el pack junto con la adquisición oficial V0.4.29 y produce 30 archivos desde Markdown.

Esto aporta un log local verificable para receipts de evidencia documental. No aporta extracción, OCR, exactitud de campos, storage WORM ni garantía productiva. El log sólo acepta hashes y referencias opacas; prohíbe bytes de documentos, valores extraídos, PII, credenciales y secretos.

## Autoridad y atribución

- Repositorio oficial: `https://github.com/transparency-dev/tessera`.
- Commit firmado: `a8f33c56b80be808712c0537cb1b71f2fd5846ea`, tree `71d51fdd8d119e5dd31394d99859c073c9544881`, 2026-08-26T18:05:57Z.
- Archive exacto: `https://github.com/transparency-dev/tessera/archive/a8f33c56b80be808712c0537cb1b71f2fd5846ea.zip`, 2.215.098 bytes, SHA-256 `33bb2394a0d6ce0ed34b3dadbf6fdac6059afada6611f694c88be1b5fbdd779c`.
- El `AUTHORS` oficial enumera Anthropic PBC, Google LLC e Internet Security Research Group; SHA-256 `7ae698cb6e29d0df7d8f485d60d8993aa1de88e04e5b85b86afda76b415c8d58`.
- Google Trillian está en maintenance mode y recomienda Tessera para operadores nuevos. La biblioteca conserva la atribución Transparency.dev/Tessera authors; no lo presenta como producto exclusivo de Google, Anthropic ni ISRG.
- Licencia Apache-2.0, `LICENSE` SHA-256 `c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4`.
- Latest stable inspeccionado: V1.0.4, commit `6bca8e8d5e23c9941f2b8a08f512b373f7131730`, archive 2.087.515 bytes/SHA-256 `6eecee02ace6750e39140339503f8b46070094e8c535f48b592cce62ea66cc23`. El release no incluye binarios y su commit no tiene firma verificable; el pack no lo sustituye silenciosamente por el head.

## Identidad del módulo y archivos empaquetados

- Pseudo-version: `v1.0.3-0.20260826180557-a8f33c56b80b`.
- Go module sum: `h1:IqUpmTdA0rRtIPRszAnMwNl+gyCZXZz0WCrqkplQAHQ=`; Go mod sum `h1:cH4u91YS85J/lmRxL/zcLEKYNg9lzLvYHQ0qq5OURN0=`.
- `go.mod` SHA-256 `c4a71771d54d1724e7729e22135cccdc23f4689eabd6854168402ad2c1f783c8`; `go.sum` `70305ce806819ad63d4ae1c3ed4cf7f958e12fe83da9009b925dcd0875ba6cc2`.
- `cmd/examples/posix-oneshot/main.go` SHA-256 `8f8e53c6c26b917d49fe131c5baa27b77588e40ebf1c0008116015d56ed63c77`.
- `cmd/examples/fsck/main.go` SHA-256 `48ed7524f671c8b0f62746a367ee637c321047d5ea0a4d8112ab20b9a36fa91a`.
- `storage/posix/README.md` SHA-256 `c50218a0df7de7c5e92cbd373a6721bc905967c950aeb67e0f2c98a1d84358e8`.
- `SECURITY.md` SHA-256 `8b8bd2544980bfe2a1b76d7a8a658ece6bd5e5d071d536c899523e89cddf3bff`.

## Reconstrucción Linux exacta

Entorno temporal WSL2/Ubuntu sobre filesystem Linux, Go oficial 1.26.7 Linux amd64 (`go1.26.7.linux-amd64.tar.gz`, SHA-256 `ffb5f8de10c62550dfddab66b36b57030721e0a44a3218e9e1181d7b59f121ca`). Para las pruebas oficiales de fault injection se extrajeron sin instalación global `strace 6.19+ds-0ubuntu5` (deb SHA-256 `aa530bb652e46beaf92d5c916403182a5ac6b2ba79a94f18e5b9da45d1799ee5`) y `libunwind8 1.8.3-0ubuntu1` (deb SHA-256 `7a9978fadf0940f45500ced0fba219cb5954f322a31068b7b32ad04f3ddc5c3f`).

Resultados sobre el commit firmado:

- `go mod verify`: PASS.
- `go test ./...`: PASS, incluida I/O POSIX e inyección de fallos con I/O y `SIGKILL`.
- `go vet ./...`: PASS.
- Builds deterministas: `posix-oneshot` SHA-256 `d36f44eaca595fb4109b1c211c37d82f46b83326bdb91e4d796d682149cc0c15`; `fsck` SHA-256 `5e75e44a93bc5ad0d56c80f2f95f4655fc19ca61f4fed383c005e036563cf591`.
- E2E oficial: diez entradas integradas; root `421c41a01e8d1f4d7b70f5fb57e358f37b098d7e60af6071b1462a3d8947b891`; un bundle truncado fue rechazado con exit 1.
- E2E del runner local: create, append y verify PASS; root de dos receipts `77fdfd4bd7a154b6a89300a4f8d95293afde046f9aefa273cad4e6b917262c3c`; invocación sin argumentos rechazada con exit 64.

## Vulnerabilidades y condición

`govulncheck` V1.7.0, commit `617f44b718537dccdea1915395650e0529e3b72e`, binario SHA-256 `8f7bc456defc10509f3821cd2d9271e730fe5f7a56251a95f22b2e0c43b9db80`, reportó cero vulnerabilidades alcanzables y cero en paquetes importados para source y binario POSIX. Eso no equivale a grafo limpio.

El build list del commit contiene 321 módulos versionados. El batch OSV devolvió cinco registros en tres módulos requeridos pero no llamados por el CLI: `GO-2026-5841` en `github.com/klauspost/compress v1.18.0`, `GO-2026-5320` en `github.com/yuin/goldmark v1.4.13`, `GO-2026-5932` en `golang.org/x/crypto v0.54.0` y `GO-2026-6179`/`GO-2026-6180` en `golang.org/x/mod v0.39.0`. Cualquier revisión futura debe repetir inventario OSV y análisis de alcanzabilidad; no se permite editar pins localmente y conservar esta evidencia.

## Gates del pack y biblioteca

- Materialización individual: 11/11 archivos y SHA-256 PASS.
- `test_contracts.py`: 4/4 PASS — identidad/condiciones, bytes oficiales, scripts y template con efectos deshabilitados.
- `bash -n` de `build_and_verify_official.sh` y `run_verified_log.sh`: PASS.
- `build_and_verify_official.sh` ejecutado contra el source exacto: `TESSERA_OFFICIAL_BUILD_PASS`.
- `run_verified_log.sh` create/append/verify y negativo noargs: PASS.
- Adquisición V0.4.29: 76 fuentes, diez perfiles, seis negativos de lock/acquirer y cinco negativos de perfiles PASS; `tamper-evident-local-log` selecciona exactamente un source.
- `VERIFY_LIBRARY_PASS`: 46 packs, 413 archivos materializables, 285 Markdown y 23 perfiles; Tessera compone 30 archivos.
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: 46 packs, 76 fuentes, 5 artefactos SDK documentales, 1 artefacto Business Central, 9 adapters, 1 orquestador documental y 1 evidence log.

## Restricciones retenidas

El commit empaquetado es un snapshot firmado más nuevo que V1.0.4, no un release estable. Sólo se admite Linux amd64 sobre filesystem POSIX real; Windows, NTFS y rutas Windows montadas en WSL quedan fuera. El template habilita cero efectos. Producción exige source receipt aprobado, claves existentes con custodia/rotación, checkpoints publicados independientemente o witnesses, backup/restore, acceso, retención, observabilidad, carga y responsables de incidentes. Un log local detecta manipulación; no impide eliminación. Storage regulatorio inmutable debe resolverse con la lane WORM aprobada del proyecto.

## Fallos preservados

Los fallos locales `LIB-FAIL-307` a `LIB-FAIL-321` y upstream `UP-FAIL-083` a `UP-FAIL-085` permanecen en `markdown_system/LIBRARY_FAILURE_LEARNING_LEDGER.md`. Incluyen ruta oficial mal supuesta, fronteras PowerShell/Bash, prerrequisitos Linux incompletos, invocaciones/rutas erróneas, pérdida de newlines del updater y restricciones reales del release/grafo. Un PASS posterior no los elimina; sus regresiones evitan recurrencia.
