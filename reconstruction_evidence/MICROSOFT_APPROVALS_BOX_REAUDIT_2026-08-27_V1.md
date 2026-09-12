# Microsoft Approvals Box MCP — reauditoría 2026-08-27 V1

## Decisión

`SAMPLE_ONLY / REJECTED_FOR_IMMEDIATE_ADOPTION`.

Es código público oficial de Microsoft y contiene ideas útiles de UI/MCP para colas, detalle, riesgo, comentarios, approve/reject y creación. No es una implementación empresarial admitida de revisión humana: el source exacto no compila con su grafo npm actual, no publica lock ni tests, no autentica la frontera HTTP, adopta un gerente demo cuando no resuelve actor y no demuestra transición concurrente/idempotente ni auditoría durable.

## Identidad oficial fijada

- repositorio: `microsoft/mcp-interactiveUI-samples`;
- componente: `oai-apps-sdk/approvals-box/node`;
- commit firmado/verificado: `8c2cb6eed8d916dd4d8af355c55133be98da95fd`, 2026-08-07;
- tree exacto obtenido por la API Git oficial: `46c63050d2d31ab533568109c643d7148fe8ad17`;
- archive del repositorio: 316.870.113 bytes, SHA-256 `6894c595e969603a8c20c33221dde040304d2a59c2130c6f05aa1ebfcbfde19a`;
- licencia MIT: SHA-256 `275b4dd619de4e16a017b10d0beec72abbbbf14ee8a2fc68f8bdb398e821f623`;
- `README.md`: 4.816 bytes, SHA-256 `fde35b4e5b930a6d597f2619d8f371f2fdf6ed02bbaec9ea925e2278ab3e93ef`;
- `package.json`: 629 bytes, SHA-256 `6dbbad8ff0dd860b9487fc8c428609462831c5bec1f93dcf07436f30cc3f41e1`;
- `src/db.ts`: SHA-256 `04103e29ecbd470b21cbf1869bf2d5adf674bc48cd6bdcf13a56366852b7d3bb`;
- `src/index.ts`: SHA-256 `08e53362efb63e37f20ba3d7bac7dace71b0d359d3f52326ae38372242e64548`;
- `src/mcp-server.ts`: SHA-256 `b7fbe420bda69895cc8e15c19e46854b122325e6e7719d8c03a4f2b35c1b25bd`;
- `src/risk.ts`: SHA-256 `30866f760a9585ed78ac46cba6bdc12fccc0f6f4b365fb1bfa1f70cb68d2491d`;
- `src/types.ts`: SHA-256 `190dae2fd26a9fa1303f6c7a020553ca9d058e6f03118f40074e176baf0b9529`.

El árbol oficial no contiene `package-lock.json`, `npm-shrinkwrap.json`, pnpm/yarn lock ni tests dentro del componente. El README lo llama explícitamente sample, auto-siembra unos 50 registros demo y propone personalizar tipos, seed y scoring. Además, su comando de clonación apunta a `mcp-apps/approvals-box/node`, mientras el path fijado en el árbol es `oai-apps-sdk/approvals-box/node`.

## Reconstrucción local sin atribución falsa

Entorno auditado: Node `v24.14.1`, npm `11.19.0`, Windows x86-64. `npm install --ignore-scripts --no-fund --no-audit` resolvió de forma mutable 144 paquetes y generó localmente un lock de 78.773 bytes/SHA-256 `f93e18667b3f37345d6268b15df698adc205a058ed3e1b70984f1b8258c0a01d`. Ese lock es evidencia de esta corrida, no un artefacto publicado por Microsoft.

Resultados exactos:

```text
npm run build: exit 1 en Windows, 32 errores TypeScript
npx tsc -p tsconfig.json --noEmit: exit 2, 32 errores TypeScript
npm audit: 0 low / 0 moderate / 0 high / 0 critical
OSV Scanner 2.5.1: 162 paquetes inspeccionados, 0 filas afectadas
```

El script oficial es `tsc -p tsconfig.json || true`. En POSIX oculta un compilador fallido; en Windows `true` no existe y el script termina uno. Los errores incluyen tipos incompatibles de `node:sqlite`, `extra.meta` inexistente frente a `_meta`, contratos `structuredContent` incompatibles y un `ApprovalListItem` sin `type`. No se parcheó source ni se trató el JS emitido con errores como build válido.

## Fronteras que bloquean adopción

1. **Identidad no autenticada.** `app.all('/mcp')` no instala middleware de autenticación. `resolveActor` confía en `extra.meta.subject` y, si falta o no coincide, devuelve siempre `mgr_201`/`Samantha Sandy`, usuario demo con capacidad de aprobar. La versión resuelta del SDK ni siquiera expone `meta` con ese nombre. Un request sin identidad verificable no puede convertirse en aprobador.
2. **Datos operativos expuestos.** `GET /stats` carece de auth y devuelve actividad agregada y reciente por actor. La app permite requests sin `Origin` y publica una frontera MCP mediante túnel de desarrollo. No existe política de issuer, audience, tenant, scopes, CSRF/origin vinculante ni rate limit.
3. **Concurrencia e idempotencia ausentes.** approve/reject ejecutan `SELECT`, luego `UPDATE approvals SET status = ... WHERE id = ?`, luego insertan auditoría, sin transacción, version/CAS, idempotency key ni condición `status='pending'`. Dos procesos pueden decidir desde el mismo estado; una caída entre update y audit deja el estado sin evidencia correspondiente. Bulk llama cada decisión individualmente y puede completar parcialmente.
4. **Persistencia no empresarial.** usa `DatabaseSync` experimental y un SQLite local; el schema activa foreign keys pero no declara relaciones `FOREIGN KEY`. No hay multi-tenant/franquicia, backup/restore, replicación, cifrado/KMS, retención inmutable, outbox, reconciliación ni handoff a un sistema de negocio.
5. **Adjuntos no cerrados.** el preview resuelve `storage/<storage_key>` y construye `/attachments/<storage_key>`, pero la frontera no demuestra allowlist de path, scan, MIME real, límites, autorización por aprobación ni una ruta HTTP durable equivalente. Los IDs recibidos desde la UI no demuestran custodia del archivo.
6. **Scoring demo.** el riesgo es una heurística local editable por umbrales y cadenas. No existe policy version, owner, corpus, evaluación, explanation integrity ni aprobación del negocio. No puede decidir autoaprobación ni sustituir revisión.

## Lo que sí puede conservarse como referencia

- definición de widgets de lista, detalle y alta;
- separación conceptual entre cálculo de acciones permitidas y presentación;
- esquema de comentarios/historial y reason obligatorio para rechazo;
- MCP tools para list/search/detail/risk/similar/approve/reject/bulk/create.

Estas ideas no habilitan copiar el componente como auth, workflow, storage o audit productivo. Cualquier adaptación sería `AUTHORED` o `ADAPTED`, con provenance explícita, y tendría que pasar los gates bloqueantes anteriores.

## Resultado

Las condiciones se registran como `UP-FAIL-101` a `UP-FAIL-103`. El componente ya está contenido byte a byte en el archive Microsoft fijado como fuente 80; el lock se amplía con artefactos adicionales exactos para impedir que su path se interprete como una fuente nueva o una promoción. Ningún source profile lo selecciona.

## Fuente oficial

- <https://github.com/microsoft/mcp-interactiveUI-samples/tree/8c2cb6eed8d916dd4d8af355c55133be98da95fd/oai-apps-sdk/approvals-box/node>
