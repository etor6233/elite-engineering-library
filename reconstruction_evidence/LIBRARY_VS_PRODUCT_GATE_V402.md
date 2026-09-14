# V402 — biblioteca lista y gate de producto separado

## Estado posterior a la publicación completa

El commit `b0ef1ad6230dd8d12e4a4aaa457f173c0c8d9e9f` publicó los packs, scripts y evidencias V402 que estaban pendientes. Esta aclaración reemplaza las frases del historial inferior que describían un repo sólo documental. La revisión de transporte y la reparación de bytes se registran en [GIT_PUBLICATION_INTEGRITY_REPAIR_V402](GIT_PUBLICATION_INTEGRITY_REPAIR_V402.md).

**READY_FOR_LIBRARY_USE / LIBRARY_INFRASTRUCTURE / LOCAL_FIXTURES**, ejecución 337 COMPLETE. Perfil **116 packs / 1.653 archivos**; biblioteca **205 packs / 2.360 archivos materializables**. La fuente normativa del inventario es el JSON del [plan integral](../markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md), no cifras antiguas en párrafos históricos del README.

El README y gap map canónicos están ligados por SHA al cierre 337: se conservan/restauran exactamente, sin modificar el evento ni fingir una nueva certificación de producto. El bloque superior del gap map expresa el cierre; toda la sección «Registro histórico conservado», incluidas sus tablas y próximos pasos, es **HISTÓRICO — NO VIGENTE — no usar para estado actual**.

**árbol post-pulido ≠ bytes del ZIP de cierre; el ZIP sigue siendo la instantánea canónica del READY original**. Los packs de producto no cambian. El árbol Git añade aclaraciones de publicación; no se regeneran los dos ZIPs.

### Dos ZIPs: cuál usar para qué

| Artefacto inmutable | Uso | SHA-256 |
|---|---|---|
| `elite-library-v402-source330-meta331-distribution336.zip` | Extraer la biblioteca y materializar el perfil aceptado en otro destino. | `cd24757744d01dbc37b666495d659f76824a8a680015eee469187e975542591f` |
| `signed-release-v402-r330/artifact.zip` | Verificar y usar el producto de referencia local con su firma y política externa. | `dba7978d04d08dc2eedb112524708e220ba18cccf86c3a2142f8b1556ed1bf76` |

Los ZIPs están fuera del repo, en `Desktop/Elite Franchise Reference V402`. Dentro de esa carpeta, `START_REFERENCE_V402.md` y `qualification/FINAL_LIBRARY_READY_V402.json` guían el uso. El gate de producto sigue DISCOVERY/BLOCK; production_authorized=false; pnpm general/Docker mantienen el diagnóstico genérico BLOCKED. ARCA es PROVEN_LOCAL en infraestructura/fixtures, conexión CONDITIONED credenciales. Daybreak/libxml2 continúa ACCESS_BLOCKED diferido, sin investigación.

FAIL384 sigue **CONTAINED_LIBRARY_USE**, no cerrado universalmente. Fuentes actuales: [assurance JSON](LIBRARY_ASSURANCE_REVIEW_V402.json), [contención final](LIBRARY_FAIL_CONTAINMENT_V402.json) y [failure lessons](../PROJECT_FAILURE_LESSONS.md). El [assurance Markdown 320](LIBRARY_ASSURANCE_REVIEW_V402.md) conserva su revisión histórica anterior al release final.

## Historial del pulido documental anterior

> **HISTÓRICO — NO VIGENTE — no usar para estado actual del repo.** Las frases siguientes sobre cambios de ingeniería aún no publicados describen exclusivamente el commit 583e76c previo a b0ef1ad.

<details>
<summary>Expediente original del pulido583e76c, conservado</summary>

# V402 — biblioteca lista y gate de producto separado

Expediente de higiene documental posterior al cierre. No crea código, admisión ni un release nuevo. Observación local: ejecución337 COMPLETE y release READY_FOR_LIBRARY_USE sin errores. El commit público base es `21064b36cdbca720bbc9c02791f7d7b45d683a9f`; este pulido no publica los cambios de ingeniería previos que permanecen en la carpeta canónica.

## Gates que coexisten

| Registro / alcance | Estado observado del cierre local | Significado |
|---|---|---|
| `PROJECT_LIBRARY_READINESS_GATE.json` y su release report | READY_FOR_LIBRARY_USE / LIBRARY_INFRASTRUCTURE / LOCAL_FIXTURES | Uso de la composición y los artefactos locales exactos; T2801–T2810 y ARCA_INFRA PROVEN_LOCAL. |
| `PROJECT_EXECUTION_STATE.json` | Revisión337, fase/status COMPLETE; release LOCAL VERIFIED | Cursor del mantenimiento terminado. No se modifica ni se añade un evento en este pulido. |
| `PROJECT_READINESS_GATE.json` | `project.status=DISCOVERY`, `platform_mode=BLOCK` | Sigue siendo el gate PROJECT/target. Su ID histórico de mantenimiento no lo convierte en el gate de release de biblioteca. No se sube a READY_TO_BUILD. |

`production_authorized=false`. Los registros de producto y sus rondas/accesos no se completan con evidencia local de biblioteca. El Preflight general conserva BLOCKED: pnpm general rechazado y Docker ausente. El instalador restringido y PostgreSQL nativo tienen pruebas propias en la referencia V402; ese hecho no admite pnpm general ni instala Docker.

## FAIL-20260907-384: contenido, no borrado

La entrada sigue listada en `PROJECT_READINESS_GATE.json`. La disposición final de biblioteca es **CONTAINED_LIBRARY_USE**, no un cierre universal: selectores, preservación NEW/EXISTING y fuentes/producto exactos están probados. La aceptación final del archivo distributivo vive en el expediente local excluido `specs/library-maintenance/tasks.md`. No se hereda aprobación productiva.

`LIBRARY_ASSURANCE_REVIEW_V402.md` conserva la revisión histórica320, cuando el release todavía estaba pendiente. El JSON de assurance posterior y la contención final deben leerse con sus propios SHA; no se altera el Markdown antiguo para hacerlo parecer actual.

Fuentes verificadas en la carpeta canónica V402 / instantánea de cierre (estos hashes no identifican necesariamente los archivos históricos del checkout público):

| Fuente | SHA-256 |
|---|---|
| `reconstruction_evidence/LIBRARY_ASSURANCE_REVIEW_V402.json` | `dd304d3d91ba2f14cb9d02e1f48f07aeb119c35830ac78e4aa029b256a10554c` |
| `reconstruction_evidence/LIBRARY_FAIL_CONTAINMENT_V402.json` | `6eb5c808bc9ca31b6aa27101f092d24471d8731371ade585020f8458d4005fe6` |
| `PROJECT_FAILURE_LESSONS.md` | `ac4cecaec7b0e71aa546e77aebbf06d3053cb8186065a94f8af092957e107b9e` |
| `reconstruction_evidence/LIBRARY_READINESS_RELEASE_REPORT_V402.json` | `2c9a2840c9c8c3efd431cd2bbaf27b756d7e66adceca890b5542a93c83932702` |
| `qualification/FINAL_LIBRARY_READY_V402.json` (referencia durable local) | `01cc08637802fc01af8edca825eae0fc3d8b99e2dc9bb508528c957b1e41e5c5` |

La guía local `START_REFERENCE_V402.md` enlaza directamente a assurance, contención y failure lessons. Esos enlaces de disco no se publican como supuestas descargas de GitHub.

## Historia y artefactos inmutables

**HISTÓRICO — NO VIGENTE — no usar para estado actual** aplica a todo el [registro histórico del gap map](../markdown_system/FRANCHISE_GAP_MAP.md#registro-histórico-conservado), incluidas las tablas inferiores y cualquier «siguiente acción» original. El gap map canónico protegido por337 permanece intacto; el banner público es una capa documental separada.

**árbol post-pulido ≠ bytes del ZIP de cierre; el ZIP sigue siendo la instantánea canónica del READY original**.

La [tabla de los dos ZIPs](../README.md#dos-zips-cuál-usar-para-qué) es la referencia de nombre, uso y SHA. No se regeneran ni reempaquetan. Los archivos locales están en `Desktop/Elite Franchise Reference V402`, con `START_REFERENCE_V402.md` y `qualification/FINAL_LIBRARY_READY_V402.json`. El repositorio público mantiene su snapshot de código anterior más estas aclaraciones: este commit documental no certifica ese checkout como producto V402 ni como una nueva distribución verificada.

## Invariantes del pulido y próximos pasos del usuario

Los documentos canónicos `README.md`, `START_FRANCHISE.md`, `markdown_system/FRANCHISE_GAP_MAP.md`, los gates/receipts y los 300 paths protegidos auditados no cambian. Las ediciones públicas se realizan en un checkout documental aislado; se conserva el árbol canónico anterior y sus 91 modificaciones/357 archivos nuevos previos fuera del commit. La guía local editable se respalda antes del pulido.

Bridge → [START_FRANCHISE](../START_FRANCHISE.md) → verify sobre la instantánea elegida. Checklist vacío para cuando el usuario conecte sus cuentas: pago, WhatsApp/Page/ML, IdP/cloud, IA y ARCA. ARCA permanece PROVEN_LOCAL en infraestructura/fixtures y CONDITIONED credenciales para conectarse. Daybreak/libxml2 conserva expediente y trigger de reapertura, sin investigación. Reglas fiscales/país, corpus y aceptación operativa/productiva siguen perteneciendo al target. No se afirma BENCH01, equivalencia payroll/POS/waitlist ni una franquicia en menos de una semana.

</details>
