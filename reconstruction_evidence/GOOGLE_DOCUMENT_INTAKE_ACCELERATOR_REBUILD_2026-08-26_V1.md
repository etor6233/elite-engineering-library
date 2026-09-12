# Google Document Intake Accelerator — reconstrucción 2026-08-26 V1

## 1. Alcance e identidad exacta

La auditoría comenzó el 2026-08-26 y cerró el 2026-08-27. Se inspeccionó y ejecutó código público oficial sin modificarlo, sin Git, sin credenciales Google Cloud y sin presentar una corrección Elite como código Google.

| Propiedad | Resultado verificado |
|---|---|
| repositorio | `GoogleCloudPlatform/document-intake-accelerator` |
| commit | `956e3cc7900338d1146d82c5ee2bd4b2551f70bb`; cabeza exacta de `main` al cierre |
| firma | GitHub API `verified=true`, reason `valid` |
| fecha del commit | 2022-11-08 |
| licencia | Apache-2.0; `LICENSE` SHA-256 `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30` |
| archive | 43.108.717 bytes; SHA-256 `38204b7365ed2df6b00633c63943d588fcb085d876e5decd7cd96d99ca2309bb` |
| árbol | 596 archivos; 48.959.485 bytes expandidos |
| lock observado | `microservices/adp_ui/package-lock.json`; SHA-256 `722a9cee3d481c92f51fa3ef7d1335b8babc3bbca0545c82670e6ee71716b703` |

Que el commit siga siendo oficial y firmado prueba identidad; no prueba que un snapshot de 2022 sea seguro o reproducible en 2026.

## 2. Toolchains oficiales aislados

No se alteró ninguna instalación global. Para reproducir los runtimes declarados por upstream se fijaron temporalmente:

- `uv 0.11.30`, archive Linux verificado con su checksum oficial;
- CPython `3.7.17`, source oficial Python.org SHA-256 `7911051ed0422fd54b8f59ffc030f7cf2ae30e0f61bda191800bb040dce4f9d2`; firma GPG válida del fingerprint `0D96DF4D4110E5C43FBFB17F2D347EA6AA65421D`;
- OpenSSL `1.1.1w`, SHA-256 `cf3098950cb4d853ad95c0841f1f9c6d3dc102dccfcacd521d93925208b76ac8`;
- libffi `3.4.4`, SHA-256 `d66c56ad259a82cf2a9dfc408b32bf5da52371500b84745f7fb8b645712df676`;
- SQLite `3.53.4`, SHA3-256 `454e45f61c6bd75b7420e7190732dea03ce6639c63ada47bbc592f67fc340338`;
- Node `17.9.1`/npm `8.11.0`, tar oficial SHA-256 `2e5fcad237d934d1bced978b5a53a7586fe83aa3c20de5f4791601789dcb4f5c` y `SHASUMS256.txt` PASS;
- Google OSV-Scanner `2.5.1`, binario/checksums/version oficiales y commit `c84fa4568f2526d0333e9a914ea8a0a5f74ad68b`;
- Firebase CLI `15.28.1`, Firestore Emulator `1.22.0` JAR SHA-256 `9b6498b7f62714d67f48f59b3818883cd682dbcd46b9f59511de81c97bb5166c` y OpenJDK `21.0.11`.

El CPython temporal fue reconstruido hasta que `ssl`, `zlib`, `ctypes` y `sqlite3` importaron correctamente. Los fallos de bootstrap y sus correcciones se conservaron en `LIB-FAIL-183` a `LIB-FAIL-198`; no se atribuyeron al upstream.

## 3. Reproducibilidad Python

Los workflows oficiales seleccionan Python 3.7, versión EOL desde 2023. El árbol contiene 23 `requirements.txt`, 16 hashes de contenido únicos, 264 entradas no comentadas y 51 entradas sin pin exacto. No hay lock integral con hashes.

`utils/fake_data_generation/requirements.txt` conserva marcadores de conflicto `<<<<<<<`, `=======` y `>>>>>>>`. No se eligió ninguna rama del conflicto porque eso habría fabricado una revisión Google inexistente.

El resolver actual instaló 67 paquetes para `common`, pero ese grafo móvil no se promovió como reconstrucción histórica. Con CPython 3.7.17 exacto, Firestore Emulator local y `PROJECT_ID=GCLOUD_PROJECT=demo-elite-audit`, la colección de `common` todavía abortó: `common/db_client.py` construye `bigquery.Client()` durante import y solicita Application Default Credentials. Se observó un item y dos errores de collection; no se proporcionaron credenciales reales ni se declaró la suite verde.

Los 163 archivos Python analizados previamente continúan con AST válido. Sintaxis no sustituye dependencias reproducibles, colección o tests.

## 4. UI exacta

El Dockerfile usa `node:17-alpine` y `nginx:alpine` móviles y ejecuta `npm install`. Bajo Node 17.9.1/npm 8.11.0, `npm ci` sí terminó e instaló 1.635 paquetes sin cambiar el lock. También emitió incompatibilidades/deprecaciones; `react-bs-datatable@3.3.1` declara Node 16.

Los gates oficiales no pasan:

- Jest aborta antes de ejecutar tests por `GlobalWorkerOptions` indefinido en `DocumentReview`: una suite failed, cero tests ejecutados;
- `CI=true npm run build` termina no-cero por errores/warnings ESLint, incluidos dos `target=_blank` sin la protección exigida.

No se cambió configuración para convertir esos resultados en PASS.

## 5. Dependencias vulnerables

`npm audit --omit=dev --json` sobre el lock exacto reportó 115 paquetes vulnerables: 18 low, 20 moderate, 68 high y 9 critical. Diez son dependencias directas y cinco no tienen fix disponible dentro del grafo auditado; el reporte referencia 155 advisories únicos.

Google OSV-Scanner 2.5.1 analizó el package lock y los 23 requirements:

- 27 fuentes de resultado;
- 189 registros de package afectados;
- 111 combinaciones únicas package-versión;
- 337 IDs primarios y severidad máxima observada 9.8;
- npm: 78 registros, 64 packages únicos y 155 IDs;
- PyPI: 111 registros, 27 packages únicos y 182 IDs.

Los conteos no sustituyen reachability, pero los críticos/altos, el runtime EOL y la ausencia de locks reproducibles impiden adopción inmediata.

## 6. Exactitud de extracción y autoaprobación

El acelerador contiene arquitectura real de clasificación, extracción, validación, matching, autoapproval y revisión humana. Su decisión automática no demuestra el estándar solicitado para almacenar datos empresariales:

- `AUTO_APPROVAL_MAPPING` fija thresholds por ejemplo de 0.2, 0.3 y 0.4;
- `get_autoapproval_status` devuelve aprobación cuando la primera combinación agregada supera sus límites;
- un `application_form` puede aprobar únicamente por extraction score;
- el valor nombrado `Avg Matching Score` es `sum(matched)`, y el matching textual usa `fuzz.token_sort_ratio` ponderado;
- faltantes reciben cero/`None`, pero no existe política visible de campo crítico o schema cerrado en esta decisión;
- no se encontró una suite dedicada que pruebe autoapproval con cero false additions, false deletions, false negatives y false positives.

Por eso este algoritmo queda `REJECTED_COMPONENT` como puerta de almacenamiento. Conservar una muestra oficial no autoriza inventarle exactitud que su código y tests no demuestran.

## 7. Seguridad del ingreso y consistencia

Firebase Auth aparece en el frontend, pero la búsqueda de los microservicios no encontró verificación servidor de ID token, `HTTPBearer`, OAuth/JWT o dependency de autorización. El upload FastAPI:

- habilita CORS `*`;
- acepta archivos por `filename.lower().endswith(".pdf")`;
- usa el nombre aportado dentro del object key GCS;
- no demuestra límite de bytes/rate, tipo por contenido, antimalware, cuarentena ni normalización de filename/tenant.

El flujo tampoco demuestra garantías cross-store: crea estado antes de subir, publica a Pub/Sub sin idempotency key visible, extracción escribe BigQuery y luego actualiza otro servicio por HTTP, y handlers amplios capturan `Exception` después de efectos parciales. La búsqueda no encontró contrato de inbox/outbox, dedupe, transacción, DLQ/replay o compensación. No se reutiliza como libro empresarial hasta probar duplicados y fallos intermedios.

## 8. Infraestructura y despliegue

El README instruye desactivar la enforcement de `constraints/compute.requireOsLogin` y eliminar `constraints/compute.vmExternalIpAccess`, o permitir IPs externas. Terraform asigna grupos de roles admin/owner a workloads, backup y CI, incluidos BigQuery admin, Datastore owner, Document AI admin, Firebase admin, Storage admin, IAM security admin, Network admin y Secret Manager admin.

No existe provider lock Terraform; hay rangos abiertos, Kubernetes 1.22.12, Actions por tags mayores, imágenes base móviles y deploy de `latest`. El pipeline usa una service-account key secret. Esa infraestructura queda `REJECTED_COMPONENT` como baseline empresarial.

## 9. Decisión y condiciones de reapertura

Clasificación final: `SAMPLE_ONLY / REJECTED_FOR_IMMEDIATE_ADOPTION`.

Se preservan commit, licencia, archive, hash y arquitectura para estudio. No se copia UI, API, autoapproval, dependencias ni IaC a un proyecto. La admisión sólo puede reabrirse ante una revisión oficial que:

1. use runtimes soportados y locks integrales inmutables;
2. pase instalación, build, suites offline y cloud separadas desde cero;
3. cierre o justifique con reachability los findings críticos/altos;
4. pruebe authn/authz servidor, secure file ingress, aislamiento y least privilege;
5. demuestre schema cerrado, evidencia por campo y cero FA/FD/FN/FP sobre corpus autorizado;
6. demuestre idempotencia, dedupe, DLQ/replay, recuperación, carga, canary y rollback.

Las condiciones se retienen en `UP-FAIL-059` a `UP-FAIL-068`. El estado global permanece `NOT_READY_UNDER_EXPANDED_USER_STANDARD`.

Fuentes oficiales: <https://github.com/GoogleCloudPlatform/document-intake-accelerator/tree/956e3cc7900338d1146d82c5ee2bd4b2551f70bb>, <https://devguide.python.org/versions/>, <https://www.python.org/downloads/release/python-3717/>, <https://nodejs.org/download/release/v17.9.1/> y <https://github.com/google/osv-scanner/releases/tag/v2.5.1>.
