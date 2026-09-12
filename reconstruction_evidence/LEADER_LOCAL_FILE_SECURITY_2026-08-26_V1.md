# Leader Local File Security — 2026-08-26 V1

## 1. Alcance y regla de verdad

Esta evidencia amplía la ingesta de archivos con código público exacto de Cisco Talos, VirusTotal/Google y Google. No afirma que una firma, regla o parser vuelva infalible un pipeline. Cada fuente conserva licencia, límites, plataforma, resultado observado y condición operativa. No se llamó a Azure, AWS, GCP ni otro servicio pago.

También se revisó el último commit oficial de `microsoft/content-processing-solution-accelerator` para comprobar si corregía los fallos de la release 2.1.2. No los corrigió; por eso no reemplaza el source estable fijado.

## 2. Cisco Talos ClamAV 1.5.4

- repositorio: `Cisco-Talos/clamav`;
- release: `clamav-1.5.4`, publicada 2026-08-07;
- commit: `fa59fca15872bb8a914ba4c68188bcc8a502cbdf`;
- commit ZIP: 14.935.925 bytes, SHA-256 `31e15ccf9cfd241f446209707b30f0c13d2209efff12541ed937115b4d75e958`;
- release source tar: 63.702.396 bytes, SHA-256 `1af1117a228f1b5bc7fa91a0dabc37848a99e7d25188e9be8043332ce721dfd3`;
- source signature: 801 bytes, SHA-256 `8acf5816c3c9dc8b02a2e11f93d90d37b25845a0481e16873454df4a30a375c0`;
- Windows x64 portable ZIP: 225.774.978 bytes, SHA-256 `0d9e0228b2674137ea1a2853566c98a0278ad52ab2582c3d6dbd75373848c395`;
- Windows signature: 801 bytes, SHA-256 `42189dd43040cd5fe4ba71a963f8d7a402b6cd2284099b653f42ca73f4055f50`;
- licencia principal: `GPL-2.0-only`; `COPYING.txt` SHA-256 `0c4fd2fa9733fc9122503797648710851e4ee6d9e4969dd33fcbd8c63cd2f584`;
- Rust lock: `Cargo.lock` SHA-256 `4d1d83e531d18e793f348db5a4baa6ec8c5d14345a285eb6959fcd4af12e1e28`.

GnuPG 2.4.8 en WSL2 verificó source tar y ZIP Windows con `Good signature` usando la clave incluida por Cisco Talos. Fingerprint observado: `5BAD CA26 65EF 59DC F8A2 3D8B 707F 0DB4 8083 6771`. El binario devolvió `ClamAV 1.5.4`.

`freshclam` descargó bases oficiales y ejecutó sus tests internos: `daily` 28104, `main` 63 y `bytecode` 339; los tres database tests pasaron. Un Markdown benigno devolvió `OK`. El fixture EICAR oficial fue interceptado por Microsoft Defender antes de que ClamAV pudiera abrirlo; ese resultado queda como `LIB-FAIL-087` y no cuenta como detección ClamAV.

Admisión: `PINNED_CANDIDATE_CONDITIONED`. Para proyecto requiere firmas auténticas y frescas, límite de bytes/tiempo/recursión/expansión, aislamiento, fail-closed, cuarentena, release/recall, HA, métricas y política GPL. La base de firmas es mutable y no se congela como capacidad futura.

## 3. VirusTotal YARA-X 1.20.0

- repositorio: `VirusTotal/yara-x`;
- release: `v1.20.0`, publicada 2026-08-24;
- commit: `60ad06971467029e77967e59d580cbbe85a1474d`;
- source ZIP: 58.622.679 bytes, SHA-256 `fe49ec00e393b89996481f8d13acd1aed388da00756389cda0953f87f7269992`;
- Windows x64 CLI ZIP: 8.640.562 bytes, SHA-256 `b1e2840bac593aea353d2b2b341f5a862c9d61c0c406d9abbbad9e1fa35163a1`;
- licencia: `BSD-3-Clause`, SHA-256 `fdf05444c9178e662fa28810d94a1fa6ec32d7be798241c98094213317265880`;
- `Cargo.lock`: SHA-256 `2b58dd867d95c8854bc150fd8065bb443f368783e010adbbd5a6587d71f0039d`.

El binario devolvió `yara-x-cli 1.20.0`. Compiló la regla oficial `true.yar` y la ejecutó contra un archivo benigno; el match esperado y ambos exit code cero fueron observados. El README oficial afirma que VirusTotal ejecuta YARA-X en producción sobre miles de millones de archivos y decenas de miles de reglas.

Admisión: `PINNED_CANDIDATE`. Es un motor de reglas, no antimalware ni detector universal. Cada ruleset necesita procedencia, revisión, versión/hash, shadow/canary, falsos positivos, rollback e incident owner.

## 4. Google safearchive

- repositorio: `google/safearchive`;
- commit: `f7ce9d7b6f9c6ecd72d0b0f16216b046e55e44dc`, 2024-10-25;
- source ZIP: 44.740 bytes, SHA-256 `438a4cce18b31123dc29c92443f416dd79fa412393fefbdef060445d84997e22`;
- licencia: `Apache-2.0`, SHA-256 `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30`;
- `go.sum`: SHA-256 `6ba1f1ec5cac9d186a7701f560308f919d1d8ddccb49cbf2aaf6281a48268c74`.

El README declara explícitamente que no es un producto Google con soporte. Su propósito declarado es reemplazar `archive/tar` y `archive/zip` con defensas contra traversal y symlinks.

Con Go oficial 1.26.7 para Linux, archive 66.890.901 bytes y SHA-256 `ffb5f8de10c62550dfddab66b36b57030721e0a44a3218e9e1181d7b59f121ca`, sus paquetes `sanitizer`, `tar` y `zip` pasaron en WSL2. En Windows con Go oficial 1.26.7 fallaron 7 tests tar y 8 zip por semántica de paths/handles.

Admisión: `SAMPLE_ONLY_CONDITIONED_LINUX`. Queda en un perfil separado. En cualquier otro target, la extracción de archives permanece prohibida hasta admitir y probar otra implementación.

## 5. Revisión del main actual de Microsoft

El commit exacto `659eaa1f503dd08b1e1aea1c72eab11c7c191d00`, fechado 2026-08-25, produjo ZIP de 18.361.248 bytes y SHA-256 `5f94e4b5e052ed736abd759e500fee42857f0b21e0d7228cc82ccbffd2026565`. Sus tres locks Python y el lock pnpm cambiaron respecto de 2.1.2, pero los fallos materiales continúan:

- Workflow: 268 PASS, 3 warnings;
- Processor: 2 errores de colección; al excluirlos, 207 PASS y 2 fallos;
- API: 1 error de colección; al excluirlo, 198 PASS y 15 fallos;
- Web: instalación frozen con `pnpm@10.28.2` y build PASS, pero 12 suites Jest fallan al transformar TypeScript/JSX y ejecutan 0 tests; el proyecto exige Node `>=22.22.0 <23`, mientras el audit Windows disponible era Node 24.

Decisión: no agregar ni promover el main. La release 2.1.2 exacta permanece `SAMPLE_ONLY_CONDITIONED`. Fallos: `UP-FAIL-038` a `UP-FAIL-040`.

## 6. Integración en Elite

`OFFICIAL-UPSTREAM-ACQUISITION-CORE` 0.4.12 fija 66 fuentes y materializa 14 archivos. Se añadieron:

- `secure-file-ingestion-leaders`: seis sources local-first, portables, con preguntas y blockers explícitos;
- `secure-archive-linux`: sólo safearchive y sólo Linux;
- perfil documental base: 27 sources, ahora incluye ClamAV/YARA-X sin arrastrar safearchive en Windows.

Cinco perfiles válidos y cinco negativos pasan. El agente debe completar approval enlazado a hashes antes de adquirir. El expediente `PROJECT_DOCUMENT_INTELLIGENCE_DECISION.md` ahora registra canales, original inmutable, lifecycle, tipo, malware, YARA, archives, PII, provenance, retries, DLQ, replay, reconciliation, restore y evidencia.

Esto entrega código upstream real y una admisión fail-closed. Todavía no entrega un único pipeline productivo conectado: faltan target/project, corpus, identities, storage, schemas, límites, rulesets, fixtures hostiles, carga, aislamiento y operación demostrados.

Fuentes oficiales: https://github.com/Cisco-Talos/clamav/tree/fa59fca15872bb8a914ba4c68188bcc8a502cbdf, https://github.com/VirusTotal/yara-x/tree/60ad06971467029e77967e59d580cbbe85a1474d, https://github.com/google/safearchive/tree/f7ce9d7b6f9c6ecd72d0b0f16216b046e55e44dc y https://github.com/microsoft/content-processing-solution-accelerator/tree/659eaa1f503dd08b1e1aea1c72eab11c7c191d00.
