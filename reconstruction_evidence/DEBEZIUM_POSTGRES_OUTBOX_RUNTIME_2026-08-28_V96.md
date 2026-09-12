# Debezium PostgreSQL outbox runtime — evidencia V96

Fecha: 2026-08-28  
Estado: `REBUILD_VERIFIED / CONDITIONED`

## Autoridad oficial

- Red Hat/Debezium source `v3.6.1.Final`, commit `63371830ceff75437f95c775b91eeeaf55a056c1`, archive SHA-256 `3cc98e…`; 48 tests focales Event Router ya demostrados.
- Quay oficial `quay.io/debezium/connect:3.6.1.Final`: OCI index `sha256:76db18f20116557b2e550d5844d2bbebc927cd9c215ec2271c4ab713cc74f386`.
- linux/amd64 manifest `sha256:1a39c97202ef3da294f9775a935d113e7a2ec7c0ab150026507682d014e99990`, config `sha256:f70c6880d091a2533bb7012b8a498af23720e4d6e7df732f080d5234673faaef`.
- linux/arm64 manifest `sha256:b988179c42d05339aa37c6eaab3439481a6ad11b3531760f3228e4e7725f9e96`, config `sha256:0e2152faa0211397c52af862f27099da60428dd81cff5a9b6605a60b943b8270`.

La configuración local es `ADAPTED`, no código Debezium verbatim. El pack contiene 8 `AUTHORED`, 2 `ADAPTED`, 0 `VERBATIM`.

## Packs y hashes

- `AWS-IDP-POSTGRES-PERSISTENCE-BOUNDARY` 0.2.0 — SHA-256 `e2c7be0ff5de0aba7f7ef17af4e57dd50935c7ac8bdd821ab7a76d55c14d1bab`; outbox append-only Debezium-shaped, RLS tenant writer y NOLOGIN CDC read-only.
- `DEBEZIUM-POSTGRES-OUTBOX-RUNTIME` 0.1.0 — SHA-256 `68f287784dd24d377a0ec5f727f51c05ff0d81a9995a28ba718d477b69255347`.

| Archivo runtime | Bytes | SHA-256 |
|---|---:|---|
| `connector.template.json` | 1.921 | `dc7032f377c2b1239a471dfb163cd2a8edc88a2ca21fcc2e1318d1a07df20dd6` |
| `consumer-contract.json` | 725 | `40defc2adc0354febfbdedeac8b9d1ecaa9cb6ad7a9eec0ec86364776be9b8d4` |
| `postgres_preflight.sql` | 654 | `bd33872fda044f9ddd7ac8698b21374aba486dbf892e31161b04269bdb03ca6b` |
| `project-profile.template.json` | 946 | `4530db3e031b6ca9eab9c541c4062061d0c8d5bc04bd6ddc85053a77b1634f6e` |
| `README.md` | 2.097 | `ea5b4e9dc3ed15acce69abd6aa4485d8918cbdab6dddb2d6538431f6be189d3b` |
| `register_connector.ps1` | 3.256 | `8334f10623a82aef5243d676f6bfee17b6c6eb8e0e976c6605da81e8f73db222` |
| `runtime-lock.json` | 1.311 | `85694132ffe5452b2d846c067c56be46b76f759cd4f98c67c335a52f1e46c139` |
| `test_validate_config.py` | 3.446 | `6500b6e8fcdb2bf0bcb1c539107e3f99c7b634f4a55d5121f6d62a79d2d2d5b8` |
| `validate_config.py` | 7.754 | `54e12314265da2447aa4954595f259905d71179e236ebc398ba4433a9c7299cc` |
| `verify_contract.ps1` | 2.759 | `171c9640f2766d561864ff90483f4f9e8db911dda2e43ebb17af8a75a9cecddd` |

## Gates

- boundary round-trip 11/11, 12/12 tests y cfn-lint 1.55.1: PASS;
- Debezium runtime round-trip 10/10, compile/PowerShell parse/preflight y 7/7 tests: PASS;
- lock tamper, platform/digest, tokens, route, heartbeat, TLS/publication/table/slot/invalid-op y duplicate JSON: fail-closed;
- secrets quedan como file-provider refs; registration exige HTTPS, auth por environment y update explícito ante divergencia;
- consumer contract fija at-least-once, header id+tenant+schema+type, inbox/hash y efecto de negocio atómico.

Memoria V96: 1.067 fallos locales + 181 condiciones upstream = 1.248 IDs, cero abiertos locales y cero duplicados.

## Límites

No se hizo pull de ~1 GB, ni launch, PostgreSQL logical replication, Kafka Connect/broker o registration live. Faltan publicación/slot/CDC membership, TLS/auth, initial snapshot, duplicates/restart/ordering/load, outage/recovery, schema evolution, observabilidad, rollback, costo y consumer inbox E2E. No exactly-once ni audit ledger.

El cierre V96 confirmó `VERIFY_LIBRARY_PASS` sobre 69 packs/651 archivos/421 Markdown/36 perfiles; email 13/116 y persistence focal 4/32. `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit` confirmó 69 packs/111 fuentes y `debezium_postgres_outbox_runtimes=1`.
