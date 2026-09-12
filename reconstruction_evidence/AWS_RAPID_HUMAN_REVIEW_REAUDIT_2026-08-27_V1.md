# AWS RAPID Human Review Reaudit — 2026-08-27 V1

## Resultado

Se auditó código público oficial de Amazon Web Services que implementa una experiencia completa de revisión documental humana: `aws-samples/review-and-assessment-powered-by-intelligent-documentation` (RAPID). La fuente aporta React, backend TypeScript/Prisma, workflow de revisión, CDK y un processor Python/Bedrock. No se copió al producto ni se presentó como listo: tanto el último release estable como el `main` firmado fallan los gates de admisión inmediata.

El source exacto entra en `OFFICIAL-UPSTREAM-ACQUISITION-CORE` sólo como `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`. Esto permite que el agente lo encuentre y adquiera por hash para investigación sin volver a buscarlo, pero impide componerlo o desplegarlo automáticamente.

## Identidad del release estable

- repositorio: `aws-samples/review-and-assessment-powered-by-intelligent-documentation`;
- release: `v1.25.1`, publicado 2026-04-07;
- commit: `68fa528c040f4184316a7b8fd214a8113656f9f1`;
- tree: `ce1f3ce4d82599d71688a21039c80b9f73ba19ca`;
- el tag anotado no está firmado; el commit subyacente figura firmado/verificado por GitHub;
- archive por commit: 7.693.308 bytes, SHA-256 `8130cace698dea37430734fe2d1ac45eff7a5c1c63c963281188f0696ba24dcf`;
- licencia MIT-0: SHA-256 `7cb750713252efd1d578837ba8785b61319109906f0c10b87adab7cf4badfc42`;
- README: SHA-256 `13cb1f194658eb5002894dbdd9b40a7c09aa85b58cfb2795ea1d25beb513dba5`;
- backend lock: `36b38266672af968b344a3ee112e064cb5f093904b89548abd4d97ae8f8f6444`;
- frontend lock: `d1a6eeece8e9480ae181cc19341e3af886b558ac8b9ad59923e83fac8bb11cc7`;
- CDK lock: `4a55f204c8cf8bb286c3825735e68c1874944d312d346af33aa86f3030277ef4`;
- review processor `uv.lock`: `a81a2fce71d85fd69acf6e777c400297db7a14393cbe2f705e83afb74c93c3aa`.

## Ejecución del release sin modificar source

| Componente | Evidencia positiva | Rechazo material |
|---|---|---|
| frontend | instalación y build de 2.492 módulos PASS; formato PASS | 25 vulnerabilidades, 17 altas; lint no inicia por `eslint --ext` incompatible con flat config; bundle principal 1.197,50 kB |
| backend | Prisma generate, TypeScript build y 32 tests/5 archivos PASS; formato PASS | 26 vulnerabilidades, incluida 1 crítica y 17 altas; lint 6 errores/34 warnings |
| CDK/Lambda | TypeScript build y 12 Jest tests/2 suites PASS | CDK 8 vulnerabilidades, 5 altas; Lambda invoke-agent 2, 1 alta; el preinstall vuelve a incorporar el backend vulnerable |
| review-item-processor | `pyproject.toml` y `uv.lock` exactos presentes | tests declarados locales llaman Bedrock; otro test usa CloudFormation/S3, un MCP `@latest` y llamadas de modelo, y captura excepciones; no es suite offline fail-closed |

No se usaron credenciales, no se desplegó AWS y no se ejecutó ninguna llamada con costo.

## Reauditoría del `main` firmado

El `main` oficial vigente al corte es el commit firmado `adb2e9ddeee09875ee6e1a7e40ab1cd1eaff5fc7`, del 2026-08-19, tree `73f373e76704d5ed68e9ef4a6fbab207c6ecee3a`. Su archive exacto mide 7.733.507 bytes y tiene SHA-256 `11310a2875d37d5d2b102af6340227f6b721a3a3e6a1b2127fcb7b10cba7cfd9`. No es release y no hereda admisión.

- frontend: instalación/build/formato PASS; 14 vulnerabilidades —11 altas—, lint todavía incompatible y bundle principal 1.199,68 kB;
- backend, desde raíz Windows corta con lock autenticado: instalación de 538 paquetes, Prisma, build, formato y 32 tests PASS; lint 6 errores/34 warnings; audit productivo 9 vulnerabilidades, 6 altas;
- por lo tanto `main` mejora algunos conteos, pero no corrige la frontera de seguridad/calidad y no reemplaza `v1.25.1`.

## Fallos y aprendizaje preservados

Las condiciones upstream quedan en `UP-FAIL-086` a `UP-FAIL-091`. Los fallos locales de investigación y sus regresiones están en `LIB-FAIL-322` a `LIB-FAIL-334`. En particular, la repetición desde una raíz corta demostró que los errores `EBUSY`/`ENOENT` iniciales eran de la ejecución Windows/ruta y no se usaron para acusar al upstream.

## Decisión para Elite

RAPID contiene arquitectura y código oficial valiosos para checklist, revisión asistida, decisión humana y trazabilidad. No es código de uso instantáneo seguro en su estado publicado. El agente puede adquirir el release exacto sólo mediante aprobación enlazada a hash y para evaluación; no puede seleccionarlo en un perfil productivo, ejecutar fixes automáticos, desplegarlo ni atribuirle una garantía AWS.

La búsqueda de una revisión humana ejecutable continúa. Para admisión se exige como mínimo: release oficial inmutable, licencia exacta, locks cerrados, cero vulnerabilidades bloqueantes alcanzables, lint/build/tests offline verdes, identidad/autorización, aislamiento, razones y conflictos de revisión, timeout/escalamiento, persistencia idempotente, browser/a11y/performance, telemetry sin PII, backup/restore y rollback. Cuenta, servicios y costos AWS requieren autorización posterior del usuario.

Fuentes oficiales gobernantes: <https://github.com/aws-samples/review-and-assessment-powered-by-intelligent-documentation/tree/68fa528c040f4184316a7b8fd214a8113656f9f1> y <https://github.com/aws-samples/review-and-assessment-powered-by-intelligent-documentation/tree/adb2e9ddeee09875ee6e1a7e40ab1cd1eaff5fc7>.
