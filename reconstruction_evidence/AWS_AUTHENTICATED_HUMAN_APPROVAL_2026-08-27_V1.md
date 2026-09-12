# AWS Authenticated Human Approval — Evidence V1

## Scope

Read-only reauditoría y adquisición local, sin Git, del repositorio oficial `aws-samples/sample-autonomous-cloud-coding-agents` para determinar si su implementación Cognito/DynamoDB de aprobación humana puede incorporarse como fuente exacta. No se usaron credenciales AWS, no se desplegó infraestructura y no se efectuaron llamadas facturables.

## Official identity

- repository: `https://github.com/aws-samples/sample-autonomous-cloud-coding-agents`
- exact commit: `1b17c28dfbd134c4cc498f3d17d905934bc94df5`
- exact tree: `3a155e9d68427ad3f57ec183649e2b2031667fc8`
- commit verification: GitHub API `verified=true`, `reason=valid`
- commit date: `2026-08-26T16:05:55Z`
- repository state at audit: active, default branch `main`, no published tag or release
- exact archive URL: `https://github.com/aws-samples/sample-autonomous-cloud-coding-agents/archive/1b17c28dfbd134c4cc498f3d17d905934bc94df5.zip`
- archive bytes: `9024752`
- archive SHA-256: `736f9aa2862ed1fae4c83e215246dea0ac01983b1cb5fe3822ee9ec19a26ee69`
- license: `MIT-0`
- `LICENSE` SHA-256: `07802d663781bf3dc08e827ca3b2bf0f9fdf7aad5032198b97837731ec3974d0`
- `README.md` SHA-256: `36353aeb264d21e28adcfa73c58aae9c28f131124f73652ff51055869b62e6cc`
- `yarn.lock` SHA-256: `b26e9d36df5bd032bbe892ed1a9231b3eae5b1d96e4a91c17cfe143ee3c948e3`
- `cdk/package.json` SHA-256: `8d90b4cf4f40e16d9ea5ad3dc26c3b35d1beff9deee85ce2b7d513d49dddcc3c`

AWS declara en el README que el ejemplo es experimental/educativo y no está destinado al uso directo en producción. La licencia se conserva `AS IS`, sin garantía.

## Exact approval artifacts

| Path | SHA-256 |
|---|---|
| `cdk/src/handlers/approve-task.ts` | `f7cdade22d370b0e55acda18084033acaa2002647b5b61c82f8c04c148360f76` |
| `cdk/src/handlers/deny-task.ts` | `5858792e5a6bf87482e3c6209dab5d1e78960350079172f6a4a3ddeebf90e248` |
| `cdk/src/handlers/shared/approval-scope.ts` | `233af4d93db6abf221983b17eb2ac088ceb49d8a9d02384682ffde8d4fe8ef51` |
| `cdk/src/handlers/shared/gateway.ts` | `5d5c4b5c60019e1ea0a1c741d2dc1e8ff35f19c7450d27cdd62818dc0be5ff29` |
| `cdk/src/constructs/task-approvals-table.ts` | `4d1f2f278a8697fcb781272540f8ce662f051c9040793314c9f1a8b7818c94b3` |
| `cdk/src/constructs/task-api.ts` | `c315ffcc2df28993b8265713eed3557827ac739206707f24aba03bd00193902e` |
| `cdk/test/handlers/approve-task.test.ts` | `739298144e33dd5c2c6d39cd4963e777bae2035bb0bbf8b99d53984bef52806a` |
| `cdk/test/handlers/deny-task.test.ts` | `15fb6a80ab47d0c2b41b69759947d906d3e070e36817b8ca0ebca22346043ec7` |
| `cdk/test/handlers/approval-concurrency.test.ts` | `fb90e49145dcfee938828d9bfc1856777fc6722669c2bc9358339a33d5321977` |
| `cdk/test/handlers/shared/approval-scope.test.ts` | `04325efdc97bcda9fb505194678ca95e36cb6133320ddd9f5d6c2ca73b2c90d8` |
| `cdk/test/constructs/task-approvals-table.test.ts` | `4763de527555c20c6e1e4f3a8237c74f4ad7ecd5fe7dd26da744ce4ed66917c7` |
| `cdk/test/constructs/task-api.test.ts` | `a37a9dffe558bd322fc7a2d2d4fa3829d4d4a9a520d5533696ca8eb876930e93` |
| `cdk/test/stacks/agent.test.ts` | `a6db5da3fa75cadea031785163142b9280af0aae02316f78e2c32b8b69dcbe04` |
| `docs/design/CEDAR_HITL_GATES.md` | `75f46eb69bcc41b40ee3cd5b0a1bbadb4c3a6cfddc68cff6cec9a5b87b805b9a` |

## Verified executable behavior

Toolchain usado: Node `v24.14.1`, Corepack `0.34.6`, Yarn `1.22.22` y OSV Scanner `2.5.1`.

1. `yarn install --frozen-lockfile --ignore-scripts --non-interactive`: PASS; el lock no cambió.
2. `yarn workspace @abca/cdk run compile`: PASS.
3. Cinco suites focalizadas de handler/scope/concurrencia/tabla: `5 passed`, `65 passed`, exit `0`.
4. `task-api.test.ts` + `agent.test.ts`: `2 passed`, `158 passed`, exit `0`.
5. Suite CDK completa offline: `200 passed`, `4142 passed`, `1 snapshot passed`, exit `0`.
6. OSV Scanner sobre el `yarn.lock` exacto: `0` filas afectadas, `0` paquetes afectados, `0` advisories.
7. El compilador y tests emitieron deprecaciones de CDK para `Version.addAlias` y `ClusterProps.containerInsights`; no rompieron el snapshot, pero se mantienen como condición de actualización.

Los tests de carrera usan mocks de AWS SDK y fijan que un writer gana y el segundo recibe rechazo condicional. Esto prueba el contrato y el comando emitido, no una carrera distribuida real contra DynamoDB.

## Security contract actually wired

- API Gateway conecta `POST /tasks/{task_id}/approve` y `POST /tasks/{task_id}/deny` a `CognitoUserPoolsAuthorizer` con `AuthorizationType.COGNITO`.
- El handler extrae `requestContext.authorizer.claims.sub`; no acepta actor del body.
- La transacción de decisión exige que la fila exista, esté `PENDING` y tenga `user_id = caller`; la tarea debe estar `AWAITING_APPROVAL` y enlazada al mismo `request_id`.
- El scope pasa por parser fail-closed; existen rate limit por usuario, 404 anti-enumeración y rechazo de tarea fuera de estado.
- Approve y deny comparten el invariante first-writer-wins.

## Residual gap that prevents promotion

El evento `approval_decision_recorded` no forma parte de la transacción de decisión. El handler lo escribe después y, si falla, registra warning pero conserva respuesta exitosa. Por ello TaskEvents no constituye por sí solo un ledger completo. Una adopción empresarial necesita transacción/outbox durable o reconciliación demostrada y receipts inmutables.

Tampoco se demostraron live Cognito/DynamoDB, tenant/franchise/resource authority del negocio, backup/PITR/restore, deployment, costos, observabilidad real, incident response, rollback ni cleanup. Esos gates requieren configuración y autorización del proyecto.

## Library admission

- source ID: `aws-autonomous-coding-agents-1b17c28`
- classification: `PINNED_CANDIDATE_CONDITIONED`
- source profile: `authenticated-human-approval-aws`
- profile output after acquisition-only approval: `BLOCKED_PENDING_USER_INPUTS_AND_LIVE_GATES`
- exact acquisition regression: archive, license, lock, manifest and `14` additional artifacts verified; source/profile receipts emitted.
- no existing executable profile selects this source implicitly.

This evidence admits exact acquisition and informed reuse. It does not call the sample production-ready and does not attribute any future Elite adapter or audit closure to AWS.
