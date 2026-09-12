# Portable Business Configuration — Reconstruction Evidence V1

```yaml
evidence_id: "PBC-20260824-V1"
pack_id: "PBC-CORE"
pack_version: "0.1.0"
pack_sha256: "6b1a94abbe5b13b7ef82336a0cd42a14d12c526255058828b1b1297d64f77034"
materialized_file_count: 7
ordered_file_hash_aggregate_sha256: "3d678c17a7b8cb3d1a57a7416e49b44b15c7d325c1963afd073f8392ae6154be"
toolchains:
  - "Python 3.14.4"
  - "Go 1.26.7 windows/amd64"
  - "jsonschema 4.26.0 (external verification environment)"
result: "PASS_CONDITIONED"
verified_at: "2026-08-24"
```

## Gates ejecutados

```text
materialize_markdown_pack.ps1                 PASS (7 hashes)
Draft202012Validator.check_schema             PASS
Draft202012Validator + FORMAT_CHECKER example PASS
python -m unittest -v                         PASS (3 tests)
gofmt -d                                      PASS (sin diff)
go test ./...                                 PASS
```

Los tests compartidos admiten el ejemplo válido y rechazan jerarquías cíclicas o inexistentes, capabilities desactivadas con permisos activos y estados de transición huérfanos. Ambos runtimes aplican códigos de error estables para referencias, unicidad, permisos, workflows y custom fields.

## Dependencia de verificación

`jsonschema 4.26.0` se instaló sólo en un directorio temporal para comprobar metaschema, Draft 2020-12 y formatos. Su wheel oficial declara MIT y SHA-256 `d489f15263b8d200f8387e64b4c3a75f06629559fb73deb8fdfb525f2dab50ce`; no se incorporó silenciosamente al pack ni al runtime de producto.

## Límite del claim

La evidencia prueba reconstrucción, schema/example y reglas semánticas comunes. Faltan un lockfile estructural portable elegido por proyecto, corpus diferencial/property exhaustivo, benchmark de 1 MiB y extensiones legales/fiscales/dominio. No demuestra que una configuración particular sea legal o segura para producción.

