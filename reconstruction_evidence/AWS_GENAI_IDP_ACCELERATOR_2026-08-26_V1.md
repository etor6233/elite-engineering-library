# AWS GenAI IDP Accelerator 0.6.5 — evidencia de admisión condicionada

Fecha de verificación: `2026-08-26`.

## 1. Resultado honesto

Se incorporó como `PINNED_CANDIDATE`, no como `REUSABLE_PACK`, el código oficial de Amazon Web Services [`aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws`](https://github.com/aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws). Es el candidato oficial más completo observado para el tramo documental: OCR, clasificación, separación, extracción estructurada, confianza por campo, reglas, evaluación contra ground truth, revisión/corrección, UI, CLI/SDK, almacenamiento, monitoreo y despliegue serverless/headless.

El código ejecuta y sus suites offline son extensas, pero la revisión exacta **no está lista para producción**: sus locks contienen vulnerabilidades actuales y el typecheck integral conserva errores. Elite permite adquirirla inmediatamente por hash y la bloquea antes de adopción productiva. No se copió ni se reatribuyó código AWS dentro de los packs.

## 2. Identidad y licencia fijadas

- release oficial: [`v0.6.5`](https://github.com/aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws/releases/tag/v0.6.5), publicada `2026-08-24`;
- commit: [`1b5fd74454e593de233342a02ee911af8ee38359`](https://github.com/aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws/tree/1b5fd74454e593de233342a02ee911af8ee38359);
- source ZIP por commit: `55.014.587` bytes, SHA-256 `62c84cdf5bed3a2029285c9a398d9532863953333dd09b115150af1cb7e0f57e`;
- `LICENSE`: `925` bytes, SHA-256 `a10cf1f6e78998b710fe868bda5b5f19120a7f36de10f06899356edc1ec0cd07`, `MIT-0`;
- `NOTICE`: `1.048` bytes, SHA-256 `717b00f175303f3ae92ab079c37e5210373ef84a15b55ef4ceef066249c0e4ba`;
- `VERSION`: `0.6.5`, SHA-256 `34bf52562bae401de106933a7565c9d3a5c8dc83c04b0b29492dd3f6f3983b7a`;
- `lib/idp_common_pkg/uv.lock`: `598.612` bytes, SHA-256 `9bf4beec0d82d2094e9bdc0357512ab1b4be154da32448e5cfcb07952382af40`;
- `src/ui/package-lock.json`: `724.612` bytes, SHA-256 `2fcaca03a304d9ee4e325d8214210d9bbc67856a6502a54dbca96bd529ce23a6`.

El árbol GitHub fue no truncado: `2.733` entradas/`2.241` archivos, `91` locks/requirements y `460` rutas de test. El ZIP pasó inspección de traversal antes de extraerse. La rama `main` consultada el mismo día (`720f3052a57c22ae0d1e50886965c065ee86d60b`) conserva byte-idéntico el lock UI vulnerable; por eso no sustituye la release.

## 3. Código oficial materialmente útil

La documentación oficial describe extracción estructurada y confianza, [despliegue headless](https://github.com/aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws/blob/1b5fd74454e593de233342a02ee911af8ee38359/docs/headless-deployment.md), [Discovery de clases/esquemas](https://github.com/aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws/blob/1b5fd74454e593de233342a02ee911af8ee38359/docs/discovery.md), CLI/evaluación y [referencias reales](https://github.com/aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws/blob/1b5fd74454e593de233342a02ee911af8ee38359/docs/references.md). Esto supera un wrapper OCR: contiene ciclo de configuración versionada, procesamiento, evaluación, ground truth y revisión.

No demuestra automáticamente facturas, proformas, packing lists ni documentos de un proyecto concreto. Esas clases, esquemas, corpus y thresholds se exigen como inputs/gates del perfil; el nombre de AWS no sustituye evidencia por campo.

## 4. Reconstrucción y suites ejecutadas

Entorno aislado WSL2, sin instalar toolchains globales ni convertir Elite en repositorio Git:

- Python oficial descargado por Astral uv: `3.13.15`;
- uv oficial `0.12.6`, artefacto Linux SHA-256 `8681d8921e7d520fb368991dcf5f9c1905b80f5bf2a265a0ed085c8d8e342477`;
- Node `22.22.3`, npm `10.9.8`;
- primera parte resuelta en una sola invocación, como ordena el Makefile para evitar dependency confusion;
- cinco paquetes `idp-*` verificados como editables desde el checkout local; `uv pip check`: compatible.

Resultados offline:

- `idp_common`: `3.257 passed`, `57 skipped`, cobertura total `59%`;
- CLI `145 passed`;
- SDK `303 passed`, `3 skipped`, `13 deselected`;
- Feature SDK `119 passed`, `21 skipped`;
- main stack extensions `202 passed`;
- feature API `2 passed`;
- seller entitlement/security/fuzz `136 passed`;
- capacity `33`, circuit breaker `33`, chat document `13`, chat stream `6`, config library `114`, SDLC/IAM `240`: todos PASS;
- UI: `1.286` paquetes desde lock, `28` archivos/`291` tests PASS; ESLint, TypeScript y Vite build PASS;
- total medido: `4.894` tests PASS, `81` skipped y `13` deselected;
- Ruff: `All checks passed`; formato: `635 files already formatted`;
- particiones IAM/servicios: PASS; `875` archivos sin scans DynamoDB filtrados no paginados; buildspec PASS; GraphQL codegen byte-estable.

El guard SDLC necesita `git grep`. El archive por commit no incluye `.git`; se creó sólo un índice temporal equivalente a `actions/checkout`, se repitió el target completo y se destruye con el laboratorio. Elite y el proyecto no requieren Git por este hecho.

## 5. Bloqueos observados, no ocultados

### Seguridad de dependencias

Google OSV-Scanner oficial `2.5.1` (`c84fa456…`) fue descargado con checksum oficial; binario Linux `57.725.090` bytes, SHA-256 `f9f25499a2c8cc367b3af45df2ea7eeca7fbccceab9c35079968f4b3652194be`.

Scan de los dos locks: `1.502` paquetes observados; exit `1`; ocho paquetes con findings y `37` ocurrencias incluyendo aliases:

- Python: `bedrock-agentcore 1.7.0` (fix `1.18.1`, pero upstream limita `<1.9`), `click 8.3.1` (fix `8.3.3`), `Pillow 12.2.0` (múltiples CVE, fix `12.3.0`, upstream la fija exactamente), `setuptools 82.0.1` (fix `83.0.0`);
- UI: `nanoid 3.3.17` (fix `3.3.18`), `react-router`/`react-router-dom 6.30.4` (fix seguro exige rama 7), `xlsx 0.20.2` con dos advisories y sin fix publicado en el canal registrado por OSV.

`npm audit --omit=dev` confirmó producción: `1 high`, `2 moderate`, `0 critical` (`nanoid`, React Router/DOM). No se ejecutó `npm audit fix`.

Dependabot abrió y cerró sin merge los PR oficiales [#554](https://github.com/aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws/pull/554) (Pillow) y [#556](https://github.com/aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws/pull/556) (React Router). Dos jobs Dependabot del commit fallan; sus logs requieren permisos sobre el repositorio ajeno. No se infirió una causa privada.

### Calidad residual

- basedpyright `1.32.1` sobre todo el árbol: `2 errors`, `44 warnings`; el workflow oficial ejecuta sólo `typecheck-pr`;
- warnings de `datetime.utcnow`, ciclos/imports/tipos y tests React sin `act(...)` permanecen visibles;
- Vite produce un chunk principal de `3.448,57 kB` y advierte tamaño superior a `3.000 kB`;
- tests offline no equivalen a AWS live, aislamiento tenant, corpus real, carga, restore, canary, seguridad ofensiva ni control de costos.

## 6. Integración en Elite

`OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.15` fija ahora `67` fuentes y materializa el nuevo source ID `aws-genai-idp-accelerator-0.6.5`. El lock valida `67/67`, la adquisición conserva cinco negativos PASS y los ocho perfiles conservan cinco negativos más una adquisición positiva offline PASS. Esa nueva regresión corrigió `LIB-FAIL-110`: un script PowerShell ya no consulta `$LASTEXITCODE` después de invocar otro script PowerShell.

- `aws-secure-document-pipeline`: `8` fuentes, incluye el acelerador;
- `document-intelligence-leaders`: `28` fuentes;
- ambos perfiles transportan los findings OSV/typecheck como `production_blockers` obligatorios;
- la adquisición continúa exigiendo respuestas completas, approval enlazado por hashes, plataforma/cuentas/costos y destino ausente.

La adquisición real de auditoría descargó/verificó el source exacto y, después de la regresión `LIB-FAIL-110`, emitió `SOURCE_PROFILE_ACQUIRED` con `production_status=BLOCKED_PENDING_USER_INPUTS_AND_LIVE_GATES`; receipt de fuente SHA-256 `5b3d1af9f7c33925f0f5b00c6f5741e486b1d0e3bb00c799749099b08d8b54c7`. El perfil de auditoría prohíbe adopción/despliegue y se elimina junto con el laboratorio.

El source puede descargarse de inmediato después de la aprobación. Su dependencia o despliegue permanece bloqueado hasta una revisión oficial corregida, rescan cero aplicable, suites completas, prueba sandbox/live y corpus del proyecto. Eso es código oficial disponible sin una afirmación falsa de producción.
