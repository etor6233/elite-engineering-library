# IBM Document Extraction Toolkit Reaudit — 2026-08-27 V1

## Resultado

Se auditó `IBM/document-extraction-toolkit`, código público oficial de IBM que reúne una interfaz React, servidor Express, worker Python, PostgreSQL/PostgREST, MinIO/S3, Sqitch, Docker Compose y Terraform alrededor de extracción de PDF con watsonx. El repositorio cubre una parte arquitectónica relevante, pero no es una implementación de revisión humana segura ni reutilizable de inmediato.

El propio README lo define como `Prototype UI`, declara que usuarios/roles son placeholders que no funcionan y exige ajustar manualmente el worker desplegado por una limitación Terraform. La evaluación exacta además falla tests/build y contiene un grafo de dependencias severamente desactualizado. La fuente sólo puede registrarse como `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`; no se embebe, compone ni despliega.

## Identidad exacta

- repositorio: `IBM/document-extraction-toolkit`;
- rama oficial: `main`; repositorio no archivado;
- no existen releases publicados al corte;
- commit: `d9144ec82194f173c5a73cc29fc08b757f6eced1`, fechado 2025-02-04;
- tree: `3f393d2e62a0e21e5417eb1f7b277e230018516d`;
- commit unsigned;
- archive por commit: 2.174.373 bytes, SHA-256 `53dc40c6f5fcdbc729cc45abee487ff2476bcccd79464ca3a042097228c1f73d`;
- licencia Apache-2.0: SHA-256 `c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4`;
- README: SHA-256 `15d76d5f7c5f4f9ef7e9a1432cd899f19e9152343992ff3ce33973340818e230`;
- requirements Python: SHA-256 `083771fb42a0bc6ddacae3f1fa1a41a9f12a5512d2eafe5ccfee8329954f0504`;
- lock webclient: SHA-256 `33c4fbbb96e6a26d6ecee8bf51557a6606120098913c6d2dd109d795cf65645a`;
- lock cliente: SHA-256 `ff5879b2fe9f7ff06590b984398a12fc95eafb4aafb3aeb05fcb927ebee7ad60`;
- lock servidor: SHA-256 `0756c8f3457471166bc7b4f46fb62239d097c4d1708a8c41fad8bdc3cff15287`.

## Contratos publicados

El `package.json` raíz y el del servidor definen su test como `Error: no test specified`. El cliente contiene dos archivos detectados por Jest, pero no una suite funcional de negocio. El README exige WatsonX API key/project, Docker, Postgres, MinIO y pasos manuales; no demuestra autorización, aislamiento, decisión/reasignación de reviewer, conflictos, trazabilidad inmutable ni rollback.

No se usaron credenciales, Docker, IBM Cloud ni servicios con costo.

## Ejecución exacta del frontend

El cliente se copió a una raíz Windows corta preservando el SHA-256 del lock y se instaló con el comando upstream `npm ci --legacy-peer-deps`:

- 1.731 paquetes instalados;
- npm informó 76 vulnerabilidades: 15 bajas, 17 moderadas, 39 altas y 5 críticas;
- npm advirtió nueve install scripts no admitidos automáticamente;
- Jest: 2 suites fallidas, 0 tests ejecutados; Babel/Jest no transforma el JSX/import del source;
- build: FAIL por import no resoluble `utils/authProviderRefresh.js` desde `src/components/Admin`.

No se parcheó Babel, el import, los locks ni el código IBM para fabricar un PASS.

## Dependencias y seguridad

`npm audit --package-lock-only` sobre los tres locks de aplicación devolvió:

| lock | dependencias | low | moderate | high | critical | total |
|---|---:|---:|---:|---:|---:|---:|
| `webclient/package-lock.json` | 97 | 0 | 0 | 4 | 1 | 5 |
| `webclient/client/package-lock.json` | 1.765 | 15 | 18 | 39 | 5 | 77 |
| `webclient/server/package-lock.json` | 426 | 8 | 8 | 12 | 5 | 33 |

OSV-Scanner oficial 2.5.1 examinó 113 dependencias Python y 1.976 entradas npm repartidas entre los locks reconocidos. Después de excluir explícitamente arrays nulos/vacíos, el resultado correcto es 130 filas afectadas y 108 combinaciones paquete-versión distintas. Los hallazgos alcanzan librerías de procesamiento de documentos y fronteras web como Pillow, PyPDF, aiohttp, cryptography, LangChain, Express, Multer, DOMPurify, pdfjs, React routing y Webpack. El scan de lock no demuestra alcanzabilidad, pero impide una admisión sin correcciones oficiales y revalidación completa.

El primer agregador local contó erróneamente paquetes no afectados; quedó registrado como `LIB-FAIL-338` y la condición nula explícita produjo estos conteos reproducibles.

## Decisión

Las condiciones upstream están en `UP-FAIL-095` a `UP-FAIL-097`. El componente queda `REJECTED_COMPONENT`: su valor es histórico/arquitectónico para evitar repetir la búsqueda, no como base de implementación.

Una revisión posterior exige release/commit firmado, locks oficiales actualizados, tests reales offline verdes, build reproducible, usuarios/roles efectivos, workflow de revisión con trazabilidad, seguridad de archivos, aislamiento, idempotencia, browser/a11y/performance, observabilidad sin PII, backup/restore y despliegue/rollback verificables. Ninguna de esas capacidades se infiere por pertenecer el repositorio a IBM.

Fuente oficial gobernante: <https://github.com/IBM/document-extraction-toolkit/tree/d9144ec82194f173c5a73cc29fc08b757f6eced1>.
