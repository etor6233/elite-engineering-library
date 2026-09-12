# V260 — captura pública conectada a PostgreSQL

Fecha: 2026-09-06. Clasificación: mantenimiento de biblioteca. Admisión: `CONDITIONED`; no promoción global ni certificación productiva.

## Resultado ejecutado

`navegador → formulario → BFF Next → API Go → PostgreSQL → receipt visible`

- Microsoft Playwright 1.62.1 fijado: Chromium desktop/mobile, Firefox desktop y WebKit desktop, **4/4 PASS** en 9,2 s de runner (9,71 s de harness).
- Catálogo leído desde Go/PostgreSQL, formulario real enviado y referencia visible con enlace a reservar. Ninguna respuesta de catálogo o lead está simulada. WebKit conserva el shim de transporte loopback declarado por el gate histórico; no prueba HTTPS productivo.
- Replay idéntico devuelve el mismo ID y header de replay; payload divergente conserva 409/IDEMPOTENCY_CONFLICT. Sin consentimiento o sin clave se rechaza con 400.
- Consulta SQL posterior: exactamente cuatro leads, cuatro consentimientos, cuatro eventos `lead.captured` y cuatro receipts completados. No se duplicaron efectos al repetir cada envío.
- Base exclusiva `elite_browser_v260`, loopback 127.0.0.1:55959, PostgreSQL 18.6, **49 migraciones** aplicadas con ON_ERROR_STOP y aserción de cantidad.
- El harness usa un tenant aleatorio y limpia sólo sus filas. No entrega TEST_DATABASE_URL/DATABASE_URL al runner. La autenticación del harness rechaza todos los intentos; no existe un bypass productivo ni se afirma OIDC probado.

## Fallo reproducido y corregido

El BFF admitía únicamente `application/json`, mientras Go publica sus errores como `application/problem+json`. Antes: la regresión esperaba 409/IDEMPOTENCY_CONFLICT y obtuvo 502/INVALID_CONTENT_TYPE; **1 FAIL y 7 PASS**. Después se admite el media type exacto de Problem Details sólo para HTTP no exitoso, preservando status/code; un problem document sobre HTTP 200 sigue rechazado. Registro: LIB-FAIL-2079 / FAIL-20260906-345.

## Verificación complementaria

- Web: frozen offline install (sin nuevas dependencias), **40/40 tests**, typecheck y build de producción PASS.
- Go 1.26.7: `go test ./... -count=1`, `go vet ./...`, `go build ./...` PASS. La ejecución general sin variables no demuestra DB/browser; esas suites se ejecutaron por separado.
- PostgreSQL: `go test ./internal/platform/postgres -count=1` PASS con URL de la base descartable y código reconstruido desde los Markdown.
- Playwright histórico: **4 runtime + 8 home/CSP PASS**. En ese comando los cuatro tests conectados se omiten explícitamente; su evidencia es la ejecución opt-in distinta descrita arriba.
- Composición posterior desde los packs canónicos: **67 packs / 704 archivos**. Comparación SHA-256 de los cinco archivos modificados contra los ejecutados: **5/5 idénticos**.
- `VERIFY_LIBRARY.ps1`: **PASS**, 160 packs / 1.395 archivos / 693 Markdown / 51 perfiles. La ruta local de staging fue rechazada en la primera pasada y retirada del expediente portable sin cambiar allowlists (LIB-FAIL-2081).
- Limpieza comprobada: cero tenants `browser-*` residuales; PostgreSQL temporal detenido correctamente al finalizar.

## Archivos y reproducción

Owners actualizados sin crear módulos de dominio duplicados:

1. `TS-GO-API-WEB-BRIDGE` 0.5.1: public-client y dos regresiones de media type.
2. `GO-ELECTROMOBILITY-PUBLIC-CRM-API` 0.2.2: test HTTP existente ampliado con `TestPublicLeadBrowserPostgres`.
3. `MICROSOFT-PLAYWRIGHT-BROWSER-GATE` 0.1.3: prueba conectada dentro de enterprise-web.spec y README con prerequisitos/comando/limpieza/no-claims.

Comando desde un backend compuesto, con `TEST_DATABASE_URL` apuntando exclusivamente a una base descartable loopback `elite_browser_*`, `ELITE_WEB_ROOT` absoluto al web construido y dependencias exactas instaladas:

```powershell
$env:ELITE_PUBLIC_LEAD_E2E = '1'
try {
  go test ./internal/platform/httpapi -run '^TestPublicLeadBrowserPostgres$' -v -count=1 -timeout=5m
  if ($LASTEXITCODE -ne 0) { throw 'Connected gate failed' }
} finally {
  Remove-Item Env:ELITE_PUBLIC_LEAD_E2E -ErrorAction SilentlyContinue
}
```

Staging observado: directorio `elite-v260-f7ac67345f984056852366b4c1fe09a2` bajo el temporal del sistema (`web`, `backend`, `roundtrip`). Los archivos canónicos son los bloques de los packs, no este staging temporal. No se distribuyen rutas del perfil local de usuario.

## Autoridad y procedencia

Hashes de los owners y perfil al cierre V260 (snapshots históricos, no locks eternos):

```json
[
  {
    "path": "implementation_packs/TYPESCRIPT_GO_API_WEB_BRIDGE.md",
    "sha256": "d7100178e3b8419dcbdb968294232a123af9e4e6eb13ef015c0a795d5695b5df"
  },
  {
    "path": "implementation_packs/GO_ELECTROMOBILITY_PUBLIC_CRM_API.md",
    "sha256": "8a0ef4ff87be00be3570143ec7a932e20a121e6cc84fd4f45b210ef01e864c5b"
  },
  {
    "path": "implementation_packs/MICROSOFT_PLAYWRIGHT_BROWSER_GATE.md",
    "sha256": "769825704f1a68438dbd3f8413c3bf43ff9be9af604d8d03619fa4614d46c74e"
  },
  {
    "path": "markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md",
    "sha256": "7e5ea3db4cb400c9f90e77d0651630b8e900a759da0f95e8ce9885289049bffb"
  }
]
```

Consulta oficial 2026-09-06:

- [Microsoft Playwright: buenas prácticas](https://playwright.dev/docs/best-practices): pruebas de comportamiento visible, aislamiento, datos controlados y locators semánticos.
- [Microsoft Playwright: API testing](https://playwright.dev/docs/api-testing): combinación de acciones de navegador con pre/postcondiciones API.
- [Microsoft: manejo de errores HTTP y Problem Details](https://learn.microsoft.com/en-us/aspnet/core/fundamentals/error-handling-api?view=aspnetcore-10.0): formatos de error HTTP; no se adopta ASP.NET ni su sample como implementación Go.

Runtime Microsoft: commit `26a9e470a7b3c7822084b09fb7f13902c5f37b51`, @playwright/test 1.62.1, Apache-2.0; lock, licencia, graph y browsers fijados permanecen sin cambios. El source-contract y runtime verifier locales pasan. No se afirma una nueva auditoría global de vulnerabilidades en V260.

La corrección, el harness, las fixtures y las aserciones son `AUTHORED`; **no son código copiado de Microsoft**. Se apoyan en el runtime oficial admitido y verifican los contratos de los componentes existentes.

## Límites que siguen abiertos

Esto cierra sólo el primer tramo del journey de referencia. No demuestra respuesta de un operador/proveedor, reserva confirmada, cotización/pedido, onboarding de roles, recuperación completa, política legal de consentimiento, antiabuse distribuido, accesibilidad con tecnología asistiva, carga, seguridad ofensiva, IdP/edge ni producción. No hay un porcentaje global sustentado ni un nuevo ZIP final certificado. El roadmap canónico conserva estas brechas.
