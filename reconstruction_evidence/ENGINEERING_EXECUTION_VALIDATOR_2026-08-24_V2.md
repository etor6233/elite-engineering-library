# Engineering Execution Validator — Reconstruction Evidence V2

```yaml
evidence_id: "EXEC-VALIDATOR-20260824-V2"
pack_id: "EXECUTION-VALIDATOR"
pack_version: "1.1.0"
pack_sha256: "10f02173e2b2ddbee564df18d50617378be88a7c21dd995f0087b87bf6520967"
materialized_file_count: 8
ordered_file_hash_aggregate_sha256: "77a75af4e1558dc74d799460f55f38aab0c711818553a46b1a27c0fa975729d2"
toolchain: "Python 3.14.4 standard library"
result: "PASS_CONDITIONED"
verified_at: "2026-08-24"
```

V2 añadió un authority root explícito por `--authority-root` o `ELITE_AUTHORITY_ROOT`. Esto permite materializar el kit dentro de un proyecto distinto sin copiar 70 manuales ni desactivar la comprobación de autoridad.

```text
materialización de 8 archivos       PASS
test_validate_project.py            PASS (7 tests)
project.example.json --level plan   PASS
authority root inexistente/inválido FAIL cerrado por diseño
```

El pack sigue condicionado a una revisión accesible de la biblioteca. El contenido generado en `__pycache__` por Python no forma parte de los ocho archivos ni del agregado.

