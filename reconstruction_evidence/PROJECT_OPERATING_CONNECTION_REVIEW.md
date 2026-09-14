# Conexión operativa de proyecto — DONE / PASS

Alcance: protocolo, instalación complementaria y continuidad de una sesión de agente mediante archivos. Expediente separado del mantenimiento V402. No se abrió una franquicia real ni se autorizó producto o producción.

## Resultado observado

18 controles PASS. Se reconstruyeron los seis archivos del pack PROJECT-OPERATING-CONNECTION 0.1.0 desde su Markdown. El complemento invocó el bridge original, conservó instrucciones ajenas y una segunda instalación produjo los mismos bytes. Las referencias ausentes/alteradas y las escrituras concurrentes/obsoletas fueron rechazadas.

El harness procesó tres etiquetas sintéticas, provocó un error real de JSON, registró el fallo y restauró el fixture. Una sesión de agente sin el historial de esta conversación descubrió T103 desde los archivos, verificó el resultado independientemente y registró el handoff. El execution validator original aceptó la revisión 5 y cinco eventos. El proyecto ficticio permanece READINESS, release NONE; la conexión DONE no se disfraza como proyecto COMPLETE.

| Control observable | Resultado |
|---|---|
| entry_routes_to_protocol | PASS |
| foreign_instructions_preserved | PASS |
| repeat_install_byte_identical | PASS |
| missing_entry_detected | PASS |
| altered_authority_reference_detected | PASS |
| unmanaged_collision_preserved | PASS |
| material_progress_failure_recovery | PASS |
| stale_revision_rejected | PASS |
| concurrent_writer_rejected | PASS |
| real_execution_validator_pass | PASS |
| fresh_agent_file_only_resume | PASS |
| no_product_admission | PASS |
| six_pack_files_reconstructed_exactly | PASS |
| protected_337_files_unchanged | PASS |
| immutable_zips_unchanged | PASS |
| signed_release_unchanged | PASS |
| official_method_snapshots_retained | PASS |
| durable_receipt_hashes_valid | PASS |

## Entorno y adapters

Ensayo de tooling local Windows, PowerShell 7 y Python 3.14, más una sesión de agente con herramientas y contexto de conversación vacío. Es evidencia de esa sesión, no de descubrimiento automático en una aplicación de terceros.

| Plataforma nativa | Ejecución de su adapter |
|---|---|
| Codex | NOT_TESTED |
| Claude | NOT_TESTED |
| Grok | NOT_TESTED |
| Otros hosts | NOT_TESTED |

Las plantillas declaran estos límites. Grok cloud necesita rutas/accesos efectivos en su entorno. No se instalaron conectores, Skills remotas, cuentas ni automatizaciones.

## Evidencia verificable

En la referencia durable, `qualification/operating-connection-v1/CONNECTION_DONE.json` enumera comprobaciones, hashes y límites. `smoke-evidence/smoke-receipt.json` conserva comandos, exit codes, artefactos y SHA; los archivos asociados se retuvieron en la misma carpeta para verificarla con `verify_smoke_receipt.py --root <smoke-evidence>`. Sus rutas históricas de ejecución apuntan al ensayo descartable; la copia durable es evidencia archivada, no un proyecto reubicado para operar.

SHA del recibo del ensayo: `bbe88317be266cc798789eae6cd02b9f57bac6aaf10bee4ff023c0126164bf42`.
La sesión nueva dejó `evidence/agent-handoff.json` y sus comprobaciones en el proyecto ficticio. Los controles del harness se ejecutaron una vez; no se reabrieron suites de dominio ni las de mantenimiento V402.

## Fuentes oficiales fijadas

Snapshots consultados el 2026-09-13, conservados junto al expediente local. Son fuentes metodológicas; el glue es AUTHORED y no se atribuye a esos proveedores.

| Fuente | URL | SHA256 del snapshot de consulta |
|---|---|---|
| dora-small-batches | [Fuente](https://dora.dev/capabilities/working-in-small-batches/) | `6b43fddc441be5fbf6f5c1e523e130b4e803d10c4d25725d7c7ea4835638af6e` |
| dora-deployment | [Fuente](https://dora.dev/capabilities/deployment-automation/) | `f1a2de151c0d15800d271c473d8bbc3ee488864aecda8290c258974908a1b198` |
| microsoft-environments | [Fuente](https://learn.microsoft.com/en-us/azure/well-architected/mission-critical/mission-critical-deployment-testing) | `74d6001b61a5f4cb0b348b947ba75a7eaa3750e7c94f5481099fef9346b2c78b` |
| microsoft-safe-deployments | [Fuente](https://learn.microsoft.com/en-us/azure/well-architected/operational-excellence/safe-deployments) | `bf6b75f609cb4b07440280bb00cf64cd5cdf942ae7afb03484e8b0e8f7b8648d` |
| grok-skills | [Fuente](https://docs.x.ai/grok-bot/skills-routines-and-automations) | `96bf4f084ab5a229c0c2bbd2a42fd40c708d18f97b154eac513a4fd330a22187` |

Reglas y relación fuente→acción→prueba: [FRANCHISE_PROJECT_OPERATING_PROTOCOL](../markdown_system/FRANCHISE_PROJECT_OPERATING_PROTOCOL.md). Código reconstruible y condiciones: [PROJECT_OPERATING_CONNECTION](../implementation_packs/PROJECT_OPERATING_CONNECTION.md). Arranque completo: [START_FRANCHISE](../START_FRANCHISE.md).

## Preservación y límites

Los 300 archivos protegidos del checkpoint, ejecución 337/eventos y los 10 archivos del release firmado coinciden con sus SHA previos. Ambos ZIPs permanecen exactos: biblioteca `cd24757744d01dbc37b666495d659f76824a8a680015eee469187e975542591f`; producto `dba7978d04d08dc2eedb112524708e220ba18cccf86c3a2142f8b1556ed1bf76`. El release original sigue READY_FOR_LIBRARY_USE con errors=[]; no se reemitió.

El árbol vivo añade protocolo/pack/evidencia y actualiza la ruta de arranque y el inventario a 206 packs / 2366 archivos / 977 Markdown. El ZIP original sigue siendo la instantánea canónica del READY V402, sin esta extensión. Los scripts nuevos se materializan fuera de la biblioteca para respetar su selector inmutable. El lock es local y requiere un único integrador entre hosts. Los hashes detectan drift; no son una firma de confianza. Un corte entre escrituras de estado y eventos se detecta al retomar y requiere recuperación preservando evidencia.

DONE de conexión no autoriza el siguiente proyecto. Esperar la instrucción explícita del usuario para comenzar una franquicia real.
