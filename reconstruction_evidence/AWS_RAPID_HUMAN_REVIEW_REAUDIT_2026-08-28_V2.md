# AWS RAPID Human Review Reaudit — 2026-08-28 V2

## Decisión

Se revalidó el `main` oficial y firmado de `aws-samples/review-and-assessment-powered-by-intelligent-documentation` (RAPID) sin modificar su source ni sus locks. El commit aporta código real de AWS para ingesta, checklist, revisión y override humano, pero permanece `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`: compila y parte de sus tests pasan, aunque sus propios lints fallan, los cuatro grafos conservan vulnerabilidades y el contrato de decisión humana no satisface identidad, concurrencia ni aislamiento empresarial.

No se agregó el `main` como una segunda autoridad adquirible: no es release y no sustituye el release estable `v1.25.1` ya fijado en `OFFICIAL-UPSTREAM-ACQUISITION-CORE`. Esta evidencia permite que el agente descarte el candidato actual sin repetir una búsqueda incompleta.

## Identidad oficial

- repositorio: <https://github.com/aws-samples/review-and-assessment-powered-by-intelligent-documentation>;
- commit: `adb2e9ddeee09875ee6e1a7e40ab1cd1eaff5fc7`;
- firma GitHub: `verified=true`, razón `valid`;
- fecha del commit: `2026-08-19T02:52:18Z`;
- tree: `73f373e76704d5ed68e9ef4a6fbab207c6ecee3a`;
- comparación con el release estable: `ahead`, 22 commits, 97 archivos modificados;
- archive: 7.733.507 bytes, SHA-256 `11310a2875d37d5d2b102af6340227f6b721a3a3e6a1b2127fcb7b10cba7cfd9`;
- licencia MIT-0: SHA-256 `7cb750713252efd1d578837ba8785b61319109906f0c10b87adab7cf4badfc42`;
- backend lock: `e5138ca75a5bf3642ca00ca7a2384649de61c3d830eea8102e7751f4d4c15ced`;
- frontend lock: `ce86a2247474c5d3a300654601155d89cc600b7237d0f0205afdb3f6da0abd65`;
- CDK lock: `ba4db436208d8ecb33b507b0237e77ecec3644b23ac701e7fe46033ff8e55838`;
- review processor lock: `bf2cc88a276283242e6ca5848cbfec1829ecf6a79f09ae48c786a6ad7f6143fd`.

## Ejecución exacta

| Área | Evidencia positiva | Condición que bloquea promoción |
|---|---|---|
| backend | npm 10.9.2 instaló 538 paquetes; Prisma generate y TypeScript build PASS; 5 suites/32 tests PASS | lint oficial: 6 errores y 34 warnings; audit completo 14 y productivo 9 —6 altas— |
| frontend | instalación PASS; TypeScript+Vite build PASS; 1.013 artefactos; bundle principal 1.199.682 bytes/SHA-256 `c61d64ef320c32f2dcef445c4e5ae4bbcf21fb069813021dd1e477f9c55be188` | el script lint usa `--ext` incompatible con su propia flat config; audit completo 14 y productivo 7 —4 altas— |
| CDK | instalación/preinstall PASS; TypeScript build y 2 suites/12 tests PASS | audit completo 6 y productivo 2 —ambas altas— |
| review-item-processor | lock frozen sincronizado con Astral uv 0.12.7; 11 tests PASS, 2 integration tests SKIP explícito | uv audit reportó 128 advisories, todos con fix publicado, sobre 14 paquetes; los skips requieren S3 desplegado y Bedrock live |

La primera instalación backend con npm 11.19.0 omitió ambos binarios opcionales esbuild de Windows presentes en el lock. Microsoft Defender informó cero detecciones relacionadas. Una extracción idéntica con npm 10.9.2 instaló los binarios SHA-256 `ec02ee9b14ab332416fedd10614dfb80eed5304d94f67745067c011934a8c3c3` y `44ce6728d54c891b1c5a6d7dbfb1a0f13419884cca0b090f1fbcf0dcd8bee0e9`; por eso el fallo local no se atribuye a AWS.

Para Python se adquirió la release oficial e inmutable Astral uv 0.12.7 desde GitHub. El ZIP Windows x86-64 coincidió con el checksum publicado y con el digest de la API: `bf1518af459a3915511a11fdc6e2f43ef9a2afa138b9d498eeb9642fe9d85218`; `uv.exe` tiene SHA-256 `a98374aee19f113b7f2527f748306b59fc3244fb548e838ec8a9893685861f96`. La documentación oficial admite binarios directos de GitHub: <https://docs.astral.sh/uv/getting-started/installation/>.

## Contrato humano observado

RAPID sí implementa:

- autenticación previa y autorización owner/admin;
- override de resultado con `userOverride=true` y `userComment`;
- `updatedAt`, transacción para actualización en cascada y referencias/evidencia del resultado.

RAPID no implementa en esa frontera:

- identidad del revisor persistida en `ReviewResult`;
- `expectedVersion`, `If-Match`, row version o conflicto optimista;
- tenant/organización/franquicia para aislamiento empresarial;
- lease/claim con expiración;
- dual control;
- tests del handler/use case de override, del comentario o del conflicto;
- vinculación comprobada entre el `jobId` de la ruta y el `resultId`: el use case recibe `reviewJobId`, pero busca autorización usando `current.reviewJobId` y no compara ambos IDs.

Una búsqueda exacta sobre schema y feature review devolvió cero hits para reviewer identity, CAS/version conflict, tenant/organization/franchise y tests de override. Por tanto, `userComment` no se eleva falsamente a identidad humana, evidencia completa ni control de concurrencia.

## Seguridad y límites de ejecución

Los audits productivos no se arreglaron con `npm audit fix`; los 128 advisories Python tampoco se ignoraron ni se resolvieron alterando `uv.lock`. No se usaron credenciales, no se llamó Bedrock/S3 y no se desplegó AWS. Las dos pruebas live quedan correctamente contabilizadas como `SKIP`, no como PASS.

La condición gobernante es `UP-FAIL-138`. El agente sólo puede usar el release estable fijado para investigación aprobada por hash; no puede copiar el `main`, desplegar RAPID ni afirmar garantía de AWS. Una release posterior debe repetir identidad, licencia, locks, lint, builds, tests offline/live autorizados, SCA, identidad humana, reason, corrections, lease, CAS, request binding, aislamiento, evidencia, recovery y rollback antes de reconsiderar promoción.

