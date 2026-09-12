# Go FinOps Core — V190

## Resultado estrecho

V190 materializa `GO-FINOPS-CORE 0.1.0`, el gobierno de costo de nube (`COST-FINOPS`): todo recurso etiquetado e inventariado (nada se escapa de la facturación), audit que detecta untagged/drift/leak, budget por tenant fail-closed y teardown que contabiliza cada recurso.

Gobernado por el contrato `COST-FINOPS` (48 superficies) y el corpus SRE (`SECURITY_SRE_CLOUD_INFRASTRUCTURE.md`). Código `AUTHORED`; no copia ningún producto cloud.

## Autoridad y procedencia

- Go stdlib (BSD-3-Clause) — core `AUTHORED`.
- PostgreSQL 18.6 (PostgreSQL License) — migración `0043` `AUTHORED`.

## Archivos materializados (8)

| Archivo | SHA-256 |
|---|---|
| internal/finops/resource.go | dd138e2cadee55b2e5ff9571cd97e2f6046f0376788fc8d7c58c3480416ec181 |
| internal/finops/inventory.go | 97b98838e50cc52f2a006e6c98cf35c13576ecb6181b4ce6477f638c2338ce82 |
| internal/finops/budget.go | e56d8e03939c99c28483e77577261f18dc2be878567941238026ad1c864eef60 |
| internal/finops/teardown.go | e80a4310601fd2090fba63cc9a4074ead5509caa437e36f3a5fb8d5eff10e116 |
| internal/finops/finops_test.go | fcbd6c3a64d146ab7e7aa9c189f916a023f4d82d8e23399c2a88f100d2caa5d2 |
| db/migrations/0043_finops.up.sql | 5736c7a328b34a450322fd7dfa9424779f97cdd8c0fd448cc7e347927b2c0f37 |
| db/migrations/0043_finops.down.sql | b5c6c45a83180321dc40a23e9f131e1f5d43bb0fd6b2d4cd39e8b173fd1cd572 |
| db/tests/0043_finops.test.sql | a0530e2237db6c06171413df4d5ea885c5a28d064e1ab8b84a385303df19bd4f |

SHA-256 del pack: `9c5a700c8a11cc59f88cf0496ce81217ad1e19ca8758e61488e33cc8798f84cb`.

## Toolchain fijado

- Go 1.26.7 windows/amd64 (runtime local verificado).
- PostgreSQL 18.6 x86-64 (cluster aislado, puerto 55444).

## Verificación ejecutada

- `go test ./internal/finops/ -count=1`: **7/7 PASS** (tags obligatorios, inventario dedup/costo, rechazo de untagged, audit detecta escapes y leak, budget fail-closed, teardown limpio vs sucio).
- `go test ./... -count=1` (13 paquetes): PASS.
- `go vet ./...`: exit 0.
- Round-trip de hashes: 8/8 bloques reproducen byte a byte; `go test` sobre el árbol materializado PASS.
- PostgreSQL 18.6 real (initdb → up → test → down): `ON_ERROR_STOP=1` exit 0; aserciones DO: untagged rechazado (`check_violation`), over-cap rechazado; down deja `finops_ns=true`.

## Cobertura de invariantes probada

1. Tags `tenant`/`environment`/`owner` obligatorios (nada sin atribuir).
2. Inventario único por (tenant, provider, id); duplicados rechazados.
3. Budget por tenant fail-closed (`spent <= cap`).
4. Audit detecta untagged, drift (desconocido) y leak (faltante); teardown exige audit limpio.

## Condiciones residuales

- Scan real del provider (inventario vs realidad cloud): adapter del proyecto.
- Perfil serverless (Cloud Run) por-uso con scale-to-zero: pendiente.
- Gate pnpm del frontend: entorno.

V190 cierra la superficie `COST-FINOPS`; el scan contra el provider real es runtime del proyecto.
