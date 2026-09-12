# Microsoft Business Central Capability Profile

Fecha de corte: 2026-08-31.

## 1. Decisión que gobierna

Este perfil sólo se aplica cuando `platform_mode: BUSINESS_CENTRAL_PLATFORM`. No convierte Business Central en una dependencia portable de Elite, no autoriza copiar su runtime y no permite atribuir a Microsoft reglas o adapters creados localmente.

Fuente oficial fijada:

```yaml
source: microsoft/BCApps
revision: 2eae56d704a1fd035d104f333602aea7091b7749
tree: f9846fb1254c6311985c6131bb2af11c7f169e1b
commit_signature_verified: true
archive_bytes: 226293141
archive_sha256: f7e984f2a1e9784a351068f42f0a0cefdc8317704b309f2a94fce5f0f91bc5b6
source_license: MIT
admission: CONDITIONAL_PLATFORM
language: AL
release_kind: development snapshot from main
```

Para revisión/generación AL, adquirir además `microsoft-bcquality-1.5` y comenzar por su `skills/entry.md`. Microsoft BCQuality es una base oficial, versionada y consumible por agentes; su release 1.5 fue fijada y sus validadores locales pasaron. Su scope actual es revisión técnica AL, no conocimiento funcional integral de Finance/SCM/Manufacturing/Service, por lo que no reemplaza las suites de producto ni este mapa de gaps.

Microsoft declara en el README del repositorio que allí construye la aplicación Business Central con el mismo código, pipelines y tooling que usa el equipo de producto. La fuente incluye System Application, Business Foundation, Base Application, first-party apps, Test Framework, Performance Toolkit y AI Test Toolkit. Esa afirmación se limita al repositorio oficial; no demuestra que este snapshot particular sea una release estable ni que pueda ejecutarse sin la plataforma Business Central. La release firmada `releases/28.4/StrictMode@cb07eef12935e07dd4258c122d0e7b1920fc1223` fue verificada por separado y omite BaseApp; no sustituye el snapshot integral.

## 2. Superficie de código real inspeccionada

El archive contiene 36.673 archivos `.al`. La Base Application mundial está en `src/Layers/W1/BaseApp` y contiene 8.154 archivos `.al`, incluidos objetos raíz además de las áreas de directorio siguientes:

| Área W1 | Archivos AL | Uso admitido |
|---|---:|---|
| `Finance` | 1.038 | contabilidad, posting y base de ledger dentro de Business Central |
| `Inventory` | 958 | items, variants, inventory y costing |
| `Service` | 637 | service management y base de postventa |
| `Sales` | 591 | cotizaciones, pedidos, ventas y clientes |
| `Manufacturing` | 530 | planificación y ejecución de fabricación |
| `Purchases` | 371 | proveedores, documentos y procesos de compra |
| `CRM` | 370 | contactos, oportunidades e integración CRM |
| `Warehouse` | 355 | depósitos, ubicaciones, movimientos, picking y shipping |
| `Projects` | 338 | proyectos, recursos y seguimiento |
| `Integration` | 332 | fronteras e integración del producto |
| `Bank` | 239 | cuentas y procesos bancarios |
| `FixedAssets` | 213 | activos fijos |
| `Foundation` | 209 | fundamentos compartidos de la Base Application |
| `HumanResources` | 115 | recursos humanos base |
| `Assembly` | 113 | assembly orders y componentes |
| `OtherCapabilities` | 99 | capacidades transversales |
| `CostAccounting` | 82 | contabilidad de costos |
| `Pricing` | 78 | precios y price calculation |
| `eServices` | 72 | servicios electrónicos |
| `Utilities` | 65 | utilidades compartidas |
| `System` | 563 | comportamiento de sistema de la aplicación |
| `CashFlow` | 49 | cash-flow forecasting y soporte de tesorería |
| `RoleCenters` | 34 | experiencias por rol |
| `Entitlements` | 25 | entitlements del producto |
| `Invoicing` | 21 | capacidad de invoicing separada |
| `FinancialMgt` | 1 | superficie residual |
| `Permissions` | 235 | permission sets y controles del producto |
| `Modules` | 287 | módulos compartidos |
| `Removed` | 2 | compatibilidad/remoción; no seleccionar como capability |

El número de archivos demuestra superficie real, no calidad suficiente ni cobertura funcional automática. La admisión final se hace por journeys, objetos AL, pruebas y runtime ejecutado, nunca por conteo.

### First-party apps Microsoft relevantes

El mismo snapshot contiene aplicaciones W1 publicadas por Microsoft que complementan BaseApp:

| App `src/Apps/W1` | Archivos AL | Capability reutilizable dentro de BC |
|---|---:|---|
| `EDocument` | 431 | documentos electrónicos, estados, processing y extensibilidad |
| `EDocumentConnectors` | 249 | conectores de e-documents; cada endpoint/país conserva sus condiciones |
| `PEPPOL` | 53 | intercambio PEPPOL condicionado a red, identidad y jurisdicción |
| `ExpenseAgent` | 445 | procesamiento agentic de gastos dentro del producto |
| `PayablesAgent` | 66 | automatización de payables dentro del producto |
| `BankAccRecWithAI` | 18 | asistencia de conciliación bancaria dentro del producto |
| `Shopify` | 589 | conector first-party Shopify |
| `Subscription Billing` | 434 | contratos/suscripciones/facturación recurrente |
| `Quality Management` | 253 | quality processes e inspección |
| `Subcontracting` | 226 | procesos de subcontratación/fabricación externa |
| `FieldServiceIntegration` | 83 | frontera first-party con field service |
| `APIV1` + `APIV2` | 411 | APIs publicadas por Business Central |
| `External File Storage - Azure Blob Service Connector` | 15 | almacenamiento externo Azure Blob |
| `External File Storage - Azure File Service Connector` | 14 | almacenamiento externo Azure Files |
| `External File Storage - SFTP Connector` | 14 | transferencia/almacenamiento SFTP |
| `External File Storage - SharePoint Connector` | 22 | documentos en SharePoint |
| `External Storage - Document Attachments` | 11 | attachments sobre external storage |
| `Email - Microsoft 365 Connector` | 7 | email vía Microsoft 365 |
| `Email - SMTP Connector` | 24 | email SMTP configurado |
| `PowerBIReports` | 304 | reporting/analytics first-party |

Estas apps no se habilitan todas por defecto. `PROJECT_BUSINESS_CENTRAL_PLAN.md` selecciona sólo las requeridas, sus dependencias y suites. Un conector first-party sigue necesitando cuenta, términos, permisos y contract probes; un agente de gastos no sustituye el expediente documental por campo.

## 3. Mapeo a las capacidades empresariales solicitadas

| Capacidad | Código Microsoft seleccionable | Estado exacto | Brecha que no se puede ocultar |
|---|---|---|---|
| procurement/imports | `Purchases`, `Inventory`, `Warehouse`, `Manufacturing`, `Assembly`, `Subcontracting`, `Quality Management`, `EDocument` | `CONDITIONAL_PLATFORM` | PI/PL específicos, despachante, aduana argentina, nacionalización y carriers requieren modelo/configuración/adapters del proyecto |
| treasury ledger | `Finance`, `Bank`, `CashFlow`, `CostAccounting`, `FixedAssets`, `ExpenseAgent`, `PayablesAgent`, `BankAccRecWithAI` | `CONDITIONAL_PLATFORM` | bancos/fintech, impuestos y facturación de la jurisdicción deben elegirse y probarse; no se observó localización `AR` en Layers |
| product catalog | `Inventory`, `Pricing`, `Sales`, `CRM` | `CONDITIONAL_PLATFORM` | schema comercial, medidas, placas, colores y documentos REVESTEX deben configurarse y validarse |
| inventory costing | `Inventory`, `Warehouse`, `CostAccounting`, `Purchases`, `Manufacturing` | `CONDITIONAL_PLATFORM` | costo importado, monedas, landed cost y concurrencia se aceptan sólo después de journeys/tests del proyecto |
| B2B/B2C commerce | `Sales`, `Pricing`, `Invoicing`, `CRM`, `Shopify`, `Subscription Billing`, `EDocument`, `PEPPOL` | `CONDITIONAL_PLATFORM` | otros marketplaces, checkout/pagos, facturación local y storefront requieren canales/adapters y cuentas reales |
| franchise network | company/location/permissions como fundamentos | `GAP_ON_TOP_OF_PLATFORM` | no se encontró un módulo oficial exacto de territorios, regalías, políticas y aislamiento REVESTEX listo para activar |
| logistics/customs | `Warehouse`, `Purchases`, `Inventory`, shipping de BaseApp, `EDocument`, external storage/SFTP | `CONDITIONAL_PLATFORM` | customs, navieras, tracking externo, containers/pallets y discrepancias requieren jurisdicción y adapters reales |
| customer/aftersales | `Service`, `CRM`, `Projects`, `Sales`, `FieldServiceIntegration` | `CONDITIONAL_PLATFORM` | garantías, instalación, claims, recall y SLA propios deben modelarse/configurarse y probarse |
| factory/production | `Manufacturing`, `Assembly`, `Inventory`, `Warehouse`, `CostAccounting`, `Subcontracting`, `Quality Management` | `CONDITIONAL_PLATFORM` | proveedores/fábricas externos, telemetría, quality gates y planificación específica requieren integración real |

Si una fila `REQUIRED` queda en `GAP_ON_TOP_OF_PLATFORM`, el proyecto no entra en `READY_TO_BUILD` hasta que el usuario autorice explícitamente configuración/extensión `AUTHORED` o cambie esa capacidad a `NONE_WITH_REASON`. No se inventa un módulo Microsoft inexistente.

## 4. Pruebas oficiales disponibles

`src/Layers/W1/Tests` contiene suites reales; entre las inspeccionadas:

| Suite | Archivos AL |
|---|---:|
| `ERM` | 146 |
| `ERM-Finance` | 62 |
| `ERM-Purchase` | 32 |
| `ERM-Sales` | 25 |
| `Bank` | 19 |
| `Cost Accounting` | 12 |
| `Physical Inventory` | 12 |
| `Prepayment` | 10 |
| `Reverse` | 18 |
| `SCM` | 77 |
| `SCM-Assembly` | 37 |
| `SCM-Costing` | 43 |
| `SCM-Manufacturing` | 38 |
| `SCM-Planning` | 26 |
| `SCM-Reservation` | 22 |
| `SCM-Service` | 57 |
| `SCM-Warehouse` | 39 |
| `SCM-Workflow` | 3 |

Estas suites son candidatos obligatorios cuando se seleccionan sus áreas. No se registran como PASS en Elite: el entorno actual no tiene Docker en modo Windows, Business Central runtime/container ni licencia/tenant para compilarlas y ejecutarlas.

El snapshot fija además manifests reproducibles `path + bytes + SHA-256` para las superficies prioritarias: Availability 49 archivos (`cf0b1c27...0189b`), Tracking 108 (`a1d432cb...a0f9f0`), Warehouse 372 (`95fab623...f08e62`) y SCM-Reservation 23 archivos incluyendo `app.json` (`86de5f6d...1c0b2`). Estos manifests prueban identidad del source adquirido, no que las suites hayan pasado sin el runtime.

V162 fija adicionalmente los cuatro objetos exactos que gobiernan la adaptación portable de transfer orders; no se atribuye a Microsoft el adapter Go/PostgreSQL:

| Path en `microsoft/BCApps@2eae56d` | Git blob | Bytes | SHA-256 observado |
|---|---|---:|---|
| `src/Layers/W1/BaseApp/Inventory/Transfer/TransferHeader.Table.al` | `0f61480f5837ed9ac48d380c10c8434f7364a5d7` | 82.670 | `952ad36743a2ed14313eb3c99a52fd1efbdd8891344919a8cbbc0417c3b8f44c` |
| `src/Layers/W1/BaseApp/Inventory/Transfer/TransferLine.Table.al` | `c78a9a3583498c2cb4ca0698ab20b473e93cfedc` | 108.502 | `49b906242b8ae707c94459301c560f489d8a885a229949f446cb40d30817fdd3` |
| `src/Layers/W1/BaseApp/Inventory/Transfer/TransferOrderPostShipment.Codeunit.al` | `211bf8f4574b57a718524d659392ce29f874e888` | 53.015 | `1e9ad92a31e0e8751d1660e6eaadf48defbba903b85a60e9486ae5a06f7af708` |
| `src/Layers/W1/BaseApp/Inventory/Transfer/TransferOrderPostReceipt.Codeunit.al` | `adede7876bfa59682c8693625676b3bd723fb4a7` | 49.811 | `5523af9ce81b5f972919681beaf48c75edb68fa7bfb92449ce557548c8d6e039` |

La documentación Microsoft Learn fijada exige shipment antes de receipt cuando existe ubicación in-transit, hace no disponible la cantidad despachada, impide editar tracking entre ambos postings y modela cada línea como demanda outbound más suministro inbound. Esas invariantes —no el runtime ni el código AL— son la frontera admitida por el pack portable.

## 5. Condiciones de ejecución reales

El `LOCAL_DEV_ENV.md` oficial exige:

- Docker Desktop ejecutando Windows containers;
- módulo PowerShell `BcContainerHelper`;
- `build/scripts/DevEnv/NewDevEnv.ps1` para crear el container Business Central;
- entorno/artefactos Business Central compatibles;
- para las views de localización, permisos de symbolic links o Developer Mode;
- una country/region view elegida y una evaluación explícita de la localización aplicable.

La fuente MIT no vuelve MIT al servicio, runtime, tenant, imágenes, datos, conectores ni licencias comerciales. Antes de empezar, `PROJECT_ACCESS_PROBES.md` debe demostrar estos recursos sin almacenar secretos en Markdown.

## 6. Plan mínimo materializable del proyecto

Cuando el usuario elige esta plataforma, el agente debe crear `PROJECT_BUSINESS_CENTRAL_PLAN.md`:

```yaml
bcapps_revision: 2eae56d704a1fd035d104f333602aea7091b7749
accepted_development_snapshot: false
target_country_view: ""
runtime_version: ""
tenant_or_container_reference: ""
license_reference: ""
selected_areas: []
selected_first_party_apps: []
required_test_suites: []
project_extensions: []
external_adapters: []
localization_gaps: []
access_probe_evidence: []
build_evidence: []
test_evidence: []
rollback_evidence: []
```

`accepted_development_snapshot: false`, un access probe ausente, una localization gap crítica o una suite requerida no ejecutada bloquean materialización del dominio. El agente puede seguir investigando y preparando el plan; no puede simular un PASS.

## 7. Gates para promoción

1. El archive adquirido coincide en bytes y SHA-256 con el source lock.
2. Se preservan `LICENSE`, notices y licencias de terceros por path.
3. El propietario acepta explícitamente el snapshot de desarrollo o se fija una release integral posterior verificada.
4. Runtime, licencia/tenant y Windows container pasan probes reales.
5. Se construye sólo la view/localización elegida y se registra el artefacto exacto.
6. Las suites oficiales de cada área seleccionada pasan en limpio.
7. Las extensiones REVESTEX están marcadas `AUTHORED`, con pruebas propias y sin atribución falsa.
8. Adapters, impuestos, facturación, bancos, aduana y carriers requeridos pasan sandbox/contract/reconciliation tests.
9. Backup/restore, aislamiento, concurrencia, carga, seguridad y rollback pasan en el target.

Hasta entonces: `CONDITIONAL_PLATFORM`, no `REUSABLE_PACK` universal y no “sistema REVESTEX terminado”.
