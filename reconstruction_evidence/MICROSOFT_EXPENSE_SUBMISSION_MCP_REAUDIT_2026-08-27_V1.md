# Microsoft Expense Submission MCP Reaudit — 2026-08-27 V1

## Resultado

Se auditó el ejemplo oficial `microsoft/mcp-interactiveUI-samples`, componente `mcp-apps/expense-submission/node`. El código demuestra MCP Apps, widgets, Entra SSO/OBO, Microsoft Graph, Azure Table Storage/Azurite y archivos de Email/OneDrive/SharePoint en un flujo de gastos. Su build aporta patrones oficiales modernos, pero el componente no implementa un handoff empresarial durable y no es reutilizable de inmediato.

La fuente sólo puede registrarse como `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`. No se embebe, compone ni despliega.

## Identidad exacta

- repositorio: `microsoft/mcp-interactiveUI-samples`;
- rama oficial: `main`; repositorio no archivado;
- commit firmado/verificado: `8c2cb6eed8d916dd4d8af355c55133be98da95fd`, fechado 2026-08-07;
- tree: `46c63050d2d31ab533568109c643d7148fe8ad17`;
- archive completo por commit: 316.870.113 bytes, SHA-256 `6894c595e969603a8c20c33221dde040304d2a59c2130c6f05aa1ebfcbfde19a`;
- licencia MIT: SHA-256 `275b4dd619de4e16a017b10d0beec72abbbbf14ee8a2fc68f8bdb398e821f623`;
- README del componente: SHA-256 `3d693620d4ae9c3d2035acede7ac79d1c1c622444a3b62226dbb713d86493436`;
- `package.json`: SHA-256 `6d2bf97ae73cc4b7a5a452a913a29a472948c00d443e2ebb312f27ecefba4afa`;
- `tsconfig.json`: SHA-256 `95d1cd23f7263091c7618994ceb35455932d39887596d0fffd4432006103a671`;
- `src/server.ts`: SHA-256 `94c25e2211eada50fd2502d716cb1bece2768a95c815174d2981f151ee3fa660`;
- `src/db.ts`: SHA-256 `d941be07e7022add79b30cc2dc829e64e260ba48d5ef999f4af019c7439585a4`;
- `appPackage/mcp-tools.json`: SHA-256 `2d1fb34a085e3bf819c4ab03820bc84844943b326bc07bb6d3f99aa4d0a4ea7a`.

El componente no publica `package-lock.json`, `pnpm-lock.yaml` ni otra resolución cerrada. Tampoco define scripts de test o lint.

## Ejecución local

Desde una copia corta que conservó el hash del manifest se ejecutó el comando publicado `npm install`, con scripts de instalación deshabilitados para la auditoría:

- 519 paquetes instalados en esa resolución temporal;
- `npm run build`: PASS;
- `npm run build:widgets`: PASS, tres widgets HTML generados;
- `npm audit` de la resolución temporal: 4 vulnerabilidades moderadas, 0 altas/críticas;
- no existen tests oficiales del componente que puedan ejecutarse.

Este PASS no convierte la resolución local en lock oficial. Una repetición futura puede resolver versiones distintas porque el manifest usa rangos abiertos.

## Bloqueos de corrección y seguridad

El source exacto conserva fronteras explícitamente no productivas:

- `DEV_API_KEY` fijo `mock_mcp_api_key` y TODO para eliminarlo;
- cache OBO en memoria y TODO para sustituirlo por cache distribuido;
- claims JWT se decodifican y registran antes de verificarse; no hay gate demostrado de issuer/tenant;
- se registran `file_url`, URLs finales y metadatos de download; los SAS URLs no deben entrar en telemetry;
- descargas y base URL dependen de URLs/forwarded host suministrados sin demostrar allowlist, tamaño, tipo real, malware scan ni límites;
- los recibos quedan en filesystem local temporal, no en evidencia inmutable ni storage autorizado;
- el único draft usa partition/row constantes `drafts/current`, sin user, tenant ni franquicia;
- `submit_expense_report` genera `RPT-${Date.now()}`, no verifica idempotency key ni concurrencia, no persiste el reporte y no llama a un manager/ERP;
- después de devolver una estructura de éxito, elimina el draft global y los recibos descargados;
- no existen tests de auth, OBO, aislamiento, repetición, timeout, concurrencia, Graph, rollback, browser o accesibilidad.

Por lo tanto el texto “deterministic updates to the system of record” describe el escenario del sample, no una garantía demostrada para adoptar este componente.

## Límites

No se usaron credenciales, Microsoft 365 Copilot, Entra, Graph, Azure, Azurite, túneles ni servicios pagos. El README exige licencia de Copilot y cuenta de desarrollo para la experiencia completa. La auditoría no autoriza un túnel anónimo ni permisos amplios de archivos.

## Decisión

Las condiciones quedan en `UP-FAIL-098` a `UP-FAIL-100`; el error local de búsqueda y su regresión están en `LIB-FAIL-339`. El componente queda `REJECTED_COMPONENT` para handoff empresarial. Puede servir como referencia exacta del protocolo/UI y de `Files.SelectedOperations.Selected`, pero ningún agente puede copiarlo como auth, storage, submit, revisión o ingesta productiva.

Una revisión futura exige lock oficial, tests offline, persistencia idempotente y aislada, transición/approval durable, validación segura de archivos/URLs, logs sin tokens/PII/SAS, issuer/audience/tenant/client estrictos, cache/secret custody, contratos Graph, browser/a11y/load, backup/restore y rollback.

Fuente oficial gobernante: <https://github.com/microsoft/mcp-interactiveUI-samples/tree/8c2cb6eed8d916dd4d8af355c55133be98da95fd/mcp-apps/expense-submission/node>.
