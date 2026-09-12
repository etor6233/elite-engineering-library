# Human Review Leader Code Reaudit — 2026-08-28 V1

## Resultado gobernante

La biblioteca no encontró un repositorio público oficial de una empresa líder que, en una sola implementación reutilizable, demuestre simultáneamente identidad humana autenticada, autorización, motivo obligatorio, corrección por campo, asignación con expiración, conflicto optimista, aprobación ligada a la solicitud, evidencia durable y separación de funciones. No se inventó una integración entre muestras incompatibles y no se promovió ningún sample como sistema productivo completo.

`OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.46` fija 95 archives oficiales. Cinco nuevas fuentes AWS/Microsoft quedan disponibles por commit y SHA-256 tras aprobación hash-linked: una como candidato condicionado y cuatro como referencias rechazadas para adopción inmediata. AWS GenAI IDP v0.6.5 ya fijado pasa de candidato genérico a condicionado explícito con evidencia focal de revisión.

## Matriz de capacidades observadas en código exacto

| Fuente oficial exacta | Identidad/RBAC | Motivo de decisión | Corrección | Expiración | Conflicto/CAS | Solicitud ligada | Evidencia | Pruebas focales | Admisión |
|---|---|---|---|---|---|---|---|---|---|
| AWS Accelerated IDP v0.6.5 `1b5fd744…` | sí: Cognito + Admin/Reviewer/Annotator | no en `completeSectionReview` | sí, server-side por campo | no para claim | sí al reclamar mediante condición DynamoDB | parcial: owner de sección | `_editHistory`, reviewedBy/email/at | backend 37 PASS; UI 291 PASS | `PINNED_CANDIDATE_CONDITIONED` |
| SAP invoice validation `990f4d865…` | roles accounting/validator | no en decisión final | sí para deducciones/retenciones/posiciones | no | no | parcial por current validator | usuario/estado y versiones de correcciones | sin lock/release integral | `SAMPLE_ONLY` condicionado |
| AWS HITL patterns `e7188fdf…` | no en la decisión remota | request reason, no decision reason | no | heartbeat 300 s | no | no: aprobación persiste por sesión | historial Step Functions parcial | scripts manuales | rechazado |
| AWS timecards `3c1d9564…` | no | no | no | no | no | no | sólo timestamp `review_completed` | no se encontraron tests | rechazado |
| AWS insurance claims `c3dc49c9…` | Cognito, MFA, grupos, `actionBy` | notas/reason opcionales | allowlist de campos | no | no `ConditionExpression` | claim ID/estado | timestamps y `actionBy` | sin suite focal de decisión | rechazado como aplicación |
| AWS Nitro multi-approver `774991f7…` | registro de aprobadores dentro del enclave | no es workflow documental | no | nonce+TTL | replay local por instancia | sí: firma del digest canónico exacto | gates y firma | 67/67 PASS | `PINNED_CANDIDATE_CONDITIONED` |
| Microsoft Azure agents escalation `a2dc4248…` | no demostrada; email fijo | no | no | no demostrada | no | conversation ID parcial | respuesta Logic App/Service Bus | sin tests focales | rechazado |

## AWS Accelerated IDP v0.6.5

- Repositorio: `aws-solutions-library-samples/accelerated-intelligent-document-processing-on-aws`.
- Release/commit: `v0.6.5`, `1b5fd74454e593de233342a02ee911af8ee38359`.
- Archive: 55.014.587 bytes, SHA-256 `62c84cdf5bed3a2029285c9a398d9532863953333dd09b115150af1cb7e0f57e`.
- Licencia: MIT-0; tag ligero y commit sin firma.
- Código focal: `src/lambda/complete_section_review/index.py`, su test, `lib/idp_common_pkg` review-curve feedback y schema GraphQL.
- El backend oficial pasó 6 tests Lambda más 31 tests de feedback. El harness Python de 28 paquetes quedó sin vulnerabilidades conocidas en OSV Scanner 2.5.1; eso no limpia el acelerador completo.
- El lock UI instaló 1.284 paquetes, typecheck pasó, la ejecución serial y el recibo JSON compacto confirmaron 28 archivos y 291/291 tests. El entrypoint JS exacto de Vite construyó 6.475 módulos.
- Condiciones retenidas: no hay argumento de motivo en la mutación final; claim sin TTL; no hay test UI focal de claim/complete/release; `npm audit` informa nanoid high y react-router/react-router-dom moderate; el script oficial de build no es portable a Windows; packaging `idp_common` falla en Windows; dos errores basedpyright y benchmark v0.6.5 no reproducible reconocido por PR oficial #669 abierto.

## Nuevas identidades del lock

| ID | Commit | Archive bytes | SHA-256 | Licencia |
|---|---|---:|---|---|
| `aws-human-in-loop-patterns-e7188fd` | `e7188fdf108b66d1052580355db6d8e54cd647d2` | 261.338 | `7e23340a5fde83ee5df5374567ffd2fb7a23ab029dad71bf79ed839c79dce852` | MIT-0 |
| `aws-timecards-bedrock-3c1d956` | `3c1d95646d7e180c039b4bf0c14a7e65b10bfe24` | 1.657.527 | `ac3f66def82446f12175a278ccf5d13a9a82f0b6d3b66e8a3e51fa46d54b462a` | MIT |
| `aws-agentic-insurance-claims-c3dc49c` | `c3dc49c93f678300b384610718db549019a3edbb` | 375.029 | `e63fccd2d3497efc5dbb8948be0d8c0292b1bc5037d4e80b7546a1ce9a9036f4` | MIT-0 |
| `aws-nitro-multi-approver-774991f` | `774991f74de933f37abee85f08e77b59d2157a65` | 3.963.432 | `267e9681a19fff42f8fbe1d5278ff7de567f6d77aa0f4104629532faf5160da8` | MIT-0 |
| `azure-agents-escalation-a2dc424` | `a2dc42483ecf72d87c6152f0da02a9a3f4fc9120` | 227.045 | `7681c198f056b8db49f858ed8fd0ccb68eccece609dd4c66eb96123c214cd501` | MIT |

## Condiciones que el agente debe aplicar

1. Adquirir sólo después de completar los inputs del perfil y producir una aprobación enlazada a los hashes exactos.
2. No copiar de una referencia rechazada a un proyecto como si fuera foundation aprobada.
3. No combinar estas fuentes afirmando autoría upstream: cualquier glue o adaptación futura será código del proyecto y necesitará pruebas propias.
4. Para revisión documental productiva, exigir antes de persistir: identidad autenticada, rol, reason obligatorio, correcciones con evidencia por campo, lease/expiry, versión/CAS, replay ligado a request, recibo durable, dual control cuando corresponda y regresiones de conflicto/timeout/reintento.
5. Mantener proveedor, costos, credenciales y despliegue fuera de la biblioteca; son decisiones y autorizaciones del proyecto.

## Recibos locales

- Source lock: 95 fuentes, IDs únicos, JSON válido.
- Perfiles: 14 válidos; AWS selecciona 14 y document-intelligence 40.
- Materialización: 23 archivos desde el Markdown canónico; los tres archivos modificados reconstruyen con SHA idéntico al árbol fuente.
- Suite upstream: `UPSTREAM_ACQUISITION_TEST_PASS negatives=9`.
- AWS Nitro: 67/67 tests oficiales PASS en venv CPython 3.12 aislado.
- No se creó Git, no se usaron credenciales cloud, no se desplegaron recursos y no se efectuaron llamadas a proveedores documentales.
