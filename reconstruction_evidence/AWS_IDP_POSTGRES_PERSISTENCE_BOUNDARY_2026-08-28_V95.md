# AWS IDP PostgreSQL persistence boundary — evidencia V95

Fecha: 2026-08-28  
Estado: `REBUILD_VERIFIED / CONDITIONED`

## Autoridades y procedencia

La composición usa contratos públicos oficiales de Amazon S3 EventBridge (`object.version-id`, major version 1), AWS SQS/Lambda partial-batch+DLQ+idempotencia, Aurora RDS Data API (`BeginTransaction`, statements parametrizados, commit/rollback), AWS transactional outbox, AWS Powertools 3.34.0 commit `376757…` y el runtime/tests Dapr 1.18.3 ya preservado. No se copió ni se atribuyó a AWS/Dapr/PostgreSQL el ensamblaje local.

El sample AWS `transactional-outbox-pattern` `23e0519…` sigue rechazado por `UP-FAIL-181`; no se reutilizó. El pack boundary contiene nueve archivos `AUTHORED`, dos `ADAPTED`, cero `VERBATIM`.

## Packs

- `AWS-IDP-EVALUATION-DECISION-WORKER` 0.2.0 — SHA-256 `9040e952d08fb1bead9ab9f3b56a243aa96f17b44015bb452cf217b2e6c6417e`; decision schema v2 incorpora la referencia exacta del snapshot sin valores extraídos.
- `AWS-IDP-POSTGRES-PERSISTENCE-BOUNDARY` 0.1.0 — SHA-256 `a3a20a82fd877c27d081a11a73d6c8a2b1898f118efec119cff947395852fda3`; 11 archivos.

| Archivo boundary | Bytes | SHA-256 |
|---|---:|---|
| `001_document_intake_outbox.sql` | 1.675 | `e405148b44db54a1a2b32e0f48a82b0fa9f4f07bec992984b42f41dc6468dc83` |
| `apply_migration.ps1` | 1.021 | `4ca5f56eb0d0d41aa840dc1e48ad91af069c80987ddc65ebec5d746acebd027e` |
| `deploy.ps1` | 3.562 | `881d32595ad2bb22ac9959f540ac8b79cf51ce855f9b47a355d16fbc52319a10` |
| `handler.py` | 21.083 | `ed95db6c86310d2d20078d731ad3050aa20fd796168c20f4ab198ed134871e41` |
| `official-contract-lock.json` | 1.500 | `238fa1efac74b0f720e2674bd26d26f8f5f5368237db7d108910e9002e1ce3ed` |
| `persistence-policy.template.json` | 420 | `f8f19fc52f122729c506b42b3af0cc2934b8f6d1e9ab43b9d35b5bbbd1cbbf3d` |
| `project-profile.template.json` | 1.049 | `ce7c94063d58b0b470af65971f33f09e99c00e87efdb25cdaa74b1c7f3c5c87e` |
| `README.md` | 2.393 | `77b03645b2cb526596d2e0b520cfaeb330130405c33ae6ea1d0a4d47823e64e3` |
| `template.yaml` | 6.508 | `4171df26580af38bdc84daa73797585e6fd80572400853cb8ffd374885df8ac5` |
| `test_handler.py` | 10.459 | `f5c475c50b45cd5a617b0c157371b9a0b95b1f75bad36dc1a7621d3d63e63c67` |
| `verify_contract.ps1` | 2.324 | `76ddf643267fa042ad199738a6e9e60a857e63ee35fdb65f6a708a3b399bbc71` |

## Contrato demostrado offline

- worker v2 round-trip 8/8, 13/13 tests y cfn-lint 1.55.1: PASS;
- boundary round-trip 11/11, 12/12 tests y cfn-lint 1.55.1: PASS;
- positive: decision/policy/snapshot exactos → tenant RLS context → intake idempotente → hash reconciliation → outbox → commit;
- negative decision: acknowledged sin llamada DB; su decision object inmutable conserva la autoridad negativa;
- duplicate: mismo intake/outbox y estado `RECONCILED`; divergencia o fallo outbox: rollback;
- cuenta/región/bucket/prefix/VersionId/KMS/Object Lock/checksum/hash, policy/profile/classes y binding snapshot se prueban fail-closed;
- EventBridge admite extensiones menores y mantiene major 1; clasificaciones duplicadas conservan cardinalidad de secciones;
- policy/project templates comienzan cerrados; no se registran valores de documento, reviewer o secrets.

El perfil focal compone 22 archivos (worker 8 + foundation 3 + boundary 11). El perfil email compone 12 packs/106 archivos. `VERIFY_LIBRARY` además fue endurecido para exigir que todo `*PACK_PLAN.md` descubierto se componga exactamente una vez.

## Memoria de fallos

`LIB-FAIL-1058` a `1063` preservan query Windows, patch atómico, cardinalidad, typo de path, versión foundation y cobertura del verificador. Memoria V95: 1.063 fallos locales + 181 condiciones upstream = 1.244 IDs, cero abiertos locales y cero duplicados.

## Límites honestos

No hubo AWS/Aurora live ni costos autorizados. Continúan condicionados: migration real, role/secret no-owner, RLS cross-tenant, S3 EventBridge, policy Object-Locked aprobada, duplicates/races/load, DLQ/redrive, outbox publisher/consumer inbox, failover/restore/rollback, observabilidad y mapping ERP/CRM específico. La frontera guarda staging revisado y emite `document.intake.accepted`; no inventa asientos, stock, facturas fiscales ni exactly-once.

## Gates globales

El cierre V95 confirmó `VERIFY_LIBRARY_PASS` con 68 packs/641 archivos/419 Markdown/36 perfiles, email 106 y focal 22. `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit` confirmó 68 packs/111 fuentes, `aws_idp_evaluation_decision_workers=1` y `aws_idp_postgres_persistence_boundaries=1`; las 13+12 pruebas ejecutaron dentro del recorrido integral.
