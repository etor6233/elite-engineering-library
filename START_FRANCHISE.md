# START_FRANCHISE — arrancar desde la copia elegida

## Qué está listo / Qué NO está listo

La **referencia local V402** está `READY_FOR_LIBRARY_USE`, ejecución337 COMPLETE, en `LIBRARY_INFRASTRUCTURE / LOCAL_FIXTURES`. El gate de producto `PROJECT_READINESS_GATE.json` conserva `DISCOVERY` y `platform_mode=BLOCK`; no se lo promueve ni se transfieren las aprobaciones del mantenimiento al target.

Para reutilizar exactamente el perfil aceptado, extraer el ZIP de biblioteca distribution336 de la [tabla única de artefactos](README.md#dos-zips-cuál-usar-para-qué). El ZIP firmado de producto sirve para verificar y usar la referencia local. Los comandos siguientes se ejecutan desde la raíz de la copia elegida; un checkout público anterior con este pulido documental no equivale a la instantánea V402.

El Preflight genérico permanece BLOCKED por pnpm general rechazado y Docker ausente. La referencia V402 se probó con el instalador restringido y PostgreSQL nativo admitidos. Live y aceptación del target siguen pendientes; ARCA está completa en infraestructura/fixtures y espera credenciales para conectarse. Daybreak/libxml2 sigue diferido. El checklist de cuentas queda vacío: no se pide ningún secreto en este pulido.

En la instalación local, abrir `Desktop/Elite Franchise Reference V402/START_REFERENCE_V402.md` y `qualification/FINAL_LIBRARY_READY_V402.json` dentro de esa carpeta. La [separación de gates](reconstruction_evidence/LIBRARY_VS_PRODUCT_GATE_V402.md) explica FAIL384 y las fuentes. Ningún paso siguiente obliga a reabrir el mantenimiento337 ya completo; aplica al proyecto consumidor y a la revisión que elija.

## Paso 0 — Separar biblioteca y producto

Conservá la biblioteca en una ubicación estable y elegí un destino de producto separado. No hace falta Git. Ejemplo PowerShell 7 para una carpeta nueva:

```powershell
$libraryPath = (Get-Location).ProviderPath # Ejecutar desde la raíz de la copia elegida.
$projectPath = Join-Path (Split-Path -Parent $libraryPath) 'Mi Franquicia'
New-Item -ItemType Directory -Path $projectPath -ErrorAction Stop
```

## Paso 1 — Instalar el bridge (Codex + Claude Code)

```powershell
& "$libraryPath\INSTALL_AGENT_BRIDGE.ps1" -LibraryRoot $libraryPath -ProjectRoot $projectPath -Agent Both
```

Esto crea la Skill progresiva y el import; no instala deps, no pide secretos. Si el proyecto ya existe, omití `New-Item` y usá el flujo EXISTING; el bridge conserva contenido ajeno y rechaza archivos administrados incompatibles.

## Paso 2 — Abrir el agente en la raíz exacta del proyecto

Abrí Codex o Claude Code **en la raíz del proyecto** (obligatorio sin Git) y
decile:

> "Usá START_FRANCHISE.md y FRANCHISE_ACCELERATOR.md. Validá el resultado de biblioteca, reutilizá el perfil y sus recibos exactos y materializá el producto en staging nuevo. Conservá los fixtures para verificarlo antes de configurar mis cuentas."

## Paso 3 — El agente ejecuta en orden

1. **Preflight** (`VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Preflight`) → declara
   toolchains ausentes (no las finge).
2. **Rondas A–H** (`PROJECT-START-READINESS-VALIDATOR`) → te pregunta negocio,
   journeys, supuestos; registra `PROJECT_READINESS_RECORD.md`.
3. **Estado reanudable y aseguramiento** (`EXECUTION-VALIDATOR 1.3.1`) → completa `implementation_assurance`, crea y valida
   `PROJECT_EXECUTION_STATE.json` y `PROJECT_EXECUTION_EVENTS.jsonl`; si el sistema ya
   existe exige inventario/delta y en cada retorno carga sólo el contexto necesario.
4. **Referencia local** → materializa y verifica con fixtures, sin cuentas del usuario. El modo `LIBRARY_INFRASTRUCTURE` se usa para la biblioteca; el readiness `PROJECT` sigue aplicando al target y no hereda autorización productiva.
5. **Composición** (`markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md`) → su bloque JSON fija el inventario y las versiones; el recibo del compositor debe coincidir con la revisión aceptada. No usar cifras históricas como inventario ni llamar productivo al resultado local.
6. **Migraciones y gates de esa revisión** → `go test ./...`, `go vet`, `go build`, PostgreSQL real y frontend con runtime admitido por digest. El pnpm publicado 11.25.0 está rechazado según V401; en la referencia V402, la receta local restringida ya fue probada en sus builds aceptados. Este commit documental no actualiza el runtime del snapshot público. No sustituirlo por una versión de PATH ni usar una instalación con red para ignorar el gate.
7. **Vertical slice** end-to-end (captación → lead/contacto → catálogo → pedido/pago o agenda → respuesta/handoff → auditoría/recovery), checkpoint con hashes/evidencia y después expansión por vertical.

Antes de implementar, leer `markdown_system/POST_DEEPSEEK_FRANCHISE_REAUDIT_2026-09-04.md` y su sucesor de alcance V402. Para materializar desde el plan exacto, ambos destinos del ejemplo siguiente deben estar ausentes:

```powershell
& "$libraryPath\VERIFY_EXECUTABLE_LIBRARY.ps1" -Mode Preflight
& "$libraryPath\materialize_markdown_pack.ps1" -PackFile "$libraryPath\implementation_packs\MARKDOWN_COMPOSITOR_CORE.md" -Destination "$projectPath\tooling-compositor"
& "$projectPath\tooling-compositor\tools\compose-markdown-project.ps1" -PlanFile "$libraryPath\markdown_system\FRANCHISE_COMPLETE_PACK_PLAN.md" -LibraryRoot $libraryPath -Destination "$projectPath\reference-staging"
```

Preflight comprueba la biblioteca y declara toolchains ausentes. El agente debe contrastar además `MATERIALIZATION_RECORD.md` y todos los archivos con el inventario aceptado, ejecutar el verify local y conservar los recibos de esa misma revisión. No mover el staging ni sobrescribir un producto existente antes de ese contraste.

## Paso 4 — Cuando quieras conectar tus cuentas

Usá los campos vacíos de `markdown_system/PROJECT_SECRETS_TEMPLATE.md` y guardá los valores en el gestor de secretos del destino; al chat y a los recibos sólo van referencias lógicas. Los adapters deben estar completos y probados localmente antes: código ausente no es una credencial pendiente. La infraestructura ARCA se prepara sin CUIT/certificados; su activación corresponde a este paso. Daybreak/libxml2 mantiene expediente y trigger de reapertura separados.

Para Cloud Run, el owner es `markdown_system/FRANCHISE_SERVERLESS_PACK_PLAN.md` con `GO-FINOPS-CORE`. El proyecto conserva sus verificaciones de reglas/país/corpus, permisos, callbacks, reconciliación, migraciones, rollout/rollback y aceptación del target antes de habilitar producción.

## Lo que NO es instantáneo (honesto)

- Responder las rondas A–H (negocio, reglas, supuestos).
- Obtener credenciales reales (cada consola tiene su trámite).

El agente te lleva de la mano en cada uno; no lo reemplaza, lo acelera.

## Partes condicionales

Sólo pueden diferirse si el blueprint las clasifica `OPTIONAL` o `NONE_WITH_REASON`. Si el negocio las clasifica `REQUIRED`, son gates y no se ocultan:

- **Fiscal / facturación** (ARCA u otra jurisdicción).
- **Homologación fiscal / notas de crédito / QR/PDF**.
- **Marketplaces** (Amazon/Google Merchant) y **ads** — solo si el canal se usa.
- **SMS** — solo si se elige ese canal.

El orden concreto lo determina el blueprint y las dependencias del journey; el agente debe entregar primero un vertical end-to-end verificable y luego ampliar sin duplicar owners.
