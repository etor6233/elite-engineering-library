# Microsoft pg_durable Human Handoff — Reconstruction Evidence V1

- verified_at: `2026-08-28`
- upstream: `https://github.com/microsoft/pg_durable`
- release: `v0.2.6`, published `2026-08-24T01:08:16Z`
- annotated tag: `67e75c7395c3689ac5790d2e4502f9f69c61ea37` (`unsigned`)
- exact commit: `6799781da58c7f51b9046445759362d2301610dd` (`verified=true`, GitHub reason `valid`)
- source archive: `pg_durable-0.2.6.tar.gz`, 814259 bytes
- archive SHA-256: `54b7dd75d7361844bec01ce1645da13b0a3c583e3e25a838249b25c1aeac8e60`
- official `SHA256SUMS` SHA-256: `f70d2dc4289f914d3a5632d51bae48858128cdde12966afed555266728ed62dc`
- pack SHA-256: `9cae15cd99c021da0579a837b7063cdb67d680189c0a2ea8dcdec18eb1f3df53`
- environment: Windows + PowerShell `7.6.4`; WSL2 Linux `6.18.33.1`; Bash `5.3.9`; Python `3.14.4`
- local Rust/PostgreSQL/Docker: unavailable; no global toolchain was installed and no cloud resource/account was touched

## Official release artifacts

| Artifact | Bytes | SHA-256 / API digest |
|---|---:|---|
| `pg-durable-postgresql-17_0.2.6-1_amd64.deb` | 2916792 | `a2bab86756d33f765b8b55b2a261552946b1641b65ac156d2da330c78fe6514c` |
| `pg-durable-postgresql-18_0.2.6-1_amd64.deb` | 2916624 | `b3bc653ed63be24602a27cf37aa5eb400f20360aace84e7fe5ce1aad18f2f5b0` |
| `pg_durable-0.2.6.tar.bz2` | 645893 | `34a3e02c624f57359ac859f1b3089b59ff673a1f24e376272b2e3b9d571e4aad` |
| `pg_durable-0.2.6.tar.gz` | 814259 | `54b7dd75d7361844bec01ce1645da13b0a3c583e3e25a838249b25c1aeac8e60` |
| `SHA256SUMS` | 399 | `f70d2dc4289f914d3a5632d51bae48858128cdde12966afed555266728ed62dc` |

El archive descargado coincidió con la línea publicada en `SHA256SUMS`. `LICENSE.txt`, `NOTICE`, `Cargo.lock` y `docs/SCENARIOS.md` estaban presentes.

## Official commit checks

La API oficial de GitHub devolvió 18 check-runs completados con `conclusion=success` para el commit exacto:

1. Build PG17 amd64 package
2. Build PG18 amd64 package
3. Build source archives
4. Clippy & Tests (PG17)
5. Clippy & Tests (PG18)
6. Dependabot
7. Docker Build & E2E Tests
8. Example Smoke Checks
9. Format Check
10. macOS Source Install (PG17)
11. Prepare Matrix
12. Publish PG17 image
13. Publish PG18 image
14. Source Install Checks
15. Upload release assets
16. Validate package inputs
17. Validate PG17 amd64 package
18. Validate PG18 amd64 package

Primary check endpoint: `https://api.github.com/repos/microsoft/pg_durable/commits/6799781da58c7f51b9046445759362d2301610dd/check-runs`.

## Local reconstruction

1. Los cuatro `examples/*/scripts/smoke_check.sh` del archive oficial terminaron correctamente; contador afirmado: `4`.
2. El pack se materializó desde Markdown en un destino vacío.
3. Se compararon contra el release los 26 archivos manifestados: `26/26` SHA-256 exactos.
4. Sobre el árbol reconstruido, `examples/invoice-approval/scripts/smoke_check.sh` produjo exactamente un marcador `Invoice approval example smoke checks passed` y exit `0`.
5. `py_compile` generó un cache temporal; se retiró el archivo exacto y el directorio vacío, y el árbol volvió a 26 fuentes.

Hashes reconstruidos:

```text
e3efac6d1903b32eff538716231c6f61c2fb489864f784ff5db6882b5b7bc16c  examples/invoice-approval/.gitignore
ae1ca693e1e6277c07f4ea62cb2e7886a3a61740d354a4f2d26a81e023eac207  examples/invoice-approval/function-app/classify_invoice/__init__.py
c19b9f527c7619a1609907ac713253895b9bd234932111391a59586d6067c13f  examples/invoice-approval/function-app/classify_invoice/function.json
6807c805e01f645e9249b28fb65e8221873a1251802b4c112746b77f202cd2e2  examples/invoice-approval/function-app/host.json
119e1c75628b69ee8a7a519230a9b18e1d7ccc16540a3c33c355b98f4d3e8e13  examples/invoice-approval/function-app/requirements.txt
f261dc2fda8b4d318aee81f92bdbe79f0a7039cc08d51905f49db181fc54ac7c  examples/invoice-approval/README.md
0e1b7918ad07e054afa5758ed63f4af66647327d84c54aa7fdf1a69123da54d1  examples/invoice-approval/scripts/cleanup_azure.sh
439d80b70a66133e653587254b531d1e63cc0d19f1e62153e57cc6d6de2fd449  examples/invoice-approval/scripts/configure_pg.sh
2702f212dcfc74915c59109a6c7cdcc7c052c543a81e71945b5ad16161421cb5  examples/invoice-approval/scripts/create_function_app.sh
f484895eaca057a9188e92bf7691f0cad7af4d0397c94f20ca143c4f6c555212  examples/invoice-approval/scripts/deploy_function.sh
771d2498d84db90df188525ef50fa549e1aee84f817f60d8f1fbee5181dae731  examples/invoice-approval/scripts/feed_invoices.sh
cbf8a30bc7bf25983efcaae16431d4819232c4880be9b7824a2824bddf43048e  examples/invoice-approval/scripts/smoke_check.sh
5a2eed5b252b324918291e97e077a7bc9ca389c83dcb5908e85e1c848bffc8b3  examples/invoice-approval/sql/01_schema.sql
340d14a266caf6bd4b9d2929e17e5cd53a5b002037af1c7efc194d5aef4925a5  examples/invoice-approval/sql/02_set_vars.sql
c4cc1e17cc4080a02880920117f80fff9ab0cbe674de5079fd99fa0d85a38f24  examples/invoice-approval/sql/03_seed_data.sql
bdb413ae553be796185f8bc5c2af6c02651fb79bab6417d00807b7888df246a6  examples/invoice-approval/sql/04_explain.sql
11bf11a2d8ecb6e22dbf24be31ffb636b72407232c010bfbd04cd06452f2505d  examples/invoice-approval/sql/05_start_workflow.sql
a3e702a973a1ea6e17e7ec5407f160779641eb2a06e54f7f3df972af48b069d3  examples/invoice-approval/sql/06_monitor.sql
eefb437da3d82c84f01b571bec6b222524c36f7ab60d1e62ef55177a0db3fa30  examples/invoice-approval/sql/07_approve.sql
e6be3350212dad8bdb8a80559e3a2f2ced4d88fc671ee362c3524804d6f13a50  examples/invoice-approval/sql/08_explain_live.sql
c5cee9e30eb1639cc75a687ef6d8c0a5f92be41d8f871b6cc58672e388e50d4e  examples/invoice-approval/sql/09_verify.sql
3fe19a9cffb334a477fa99cb11d20b8c3bcf334c19bae9685cb6ee64239f4c0f  examples/invoice-approval/sql/10_cancel.sql
45e4e2defc048a9f4f9a5ee97705f3ef31a750b1ab9fd3335e0ddd4506b6d6b8  examples/invoice-approval/sql/11_reset.sql
4b9dd8bc9b5267d9ee86dc3d91072e5b352ad04ce298d4c1c0ec1804d06ac254  examples/invoice-approval/without-df.md
aab465651982b811f0dcd1f6403eb8e89282b75d3bf6ac685ada7e67023a1009  LICENSE.txt
a5b3fe6b28cdfefbbdb2bd0fd60d5329144aa531c8ab92f9794acf4f4ca06675  NOTICE
```

## Non-claims and retained condition

- `live_smoke_check.sh` no termina en LF. El formato de bloques Markdown no puede preservarlo VERBATIM antes del fence; no se materializa y permanece disponible únicamente en el archive oficial hash-bound.
- El ejemplo prueba sintaxis, JSON, Python y el patrón signal/timeout; no autentica el string `approver` enviado en el payload.
- El runtime sí implementa ownership/RLS: el test oficial rechaza que Bob señale una instancia de Alice. Eso no equivale a un modelo completo de revisores humanos.
- Faltan reason obligatorio, correcciones por campo, versión/conflict detection, separación de funciones, corpus documental, backup/restore, rollback y piloto productivo.
- `azure-functions>=1.20.0` es un rango abierto upstream y requiere un lock aprobado por el proyecto.
- Microsoft declara pg_durable Preview y la imagen GHCR sólo para evaluación/aprendizaje. El pack permanece `CONDITIONED`; `UP-FAIL-127` conserva esta frontera.
