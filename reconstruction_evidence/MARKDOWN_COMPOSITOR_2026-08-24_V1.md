# Reconstruction Evidence — Markdown Compositor V1

```yaml
evidence_id: "MD-COMPOSITOR-20260824-V1"
pack_id: "MARKDOWN-COMPOSITOR"
pack_version: "0.1.0"
pack_sha256: "c286eeed8e028b0f188fc88ddcc9452c06523e1af3378eb5681cf24a91a34985"
materialized_file_count: 3
ordered_file_hash_aggregate_sha256: "014eeb69fa2ada79c1d8fd2092c287ef5064542ea1b87d25d88e81f65044af9a"
toolchain: "PowerShell 7.6.4"
result: "PASS"
```

El target temporal comenzó vacío y recibió exclusivamente los tres archivos declarados por `MARKDOWN_COMPOSITOR_CORE.md`; cada SHA-256 fue comprobado antes de escribir.

Pruebas aprobadas:

- composición satisfactoria con sustitución no secreta y `MATERIALIZATION_RECORD.md`;
- rechazo de pack condicionado sin aceptación explícita;
- rechazo de un bloque cuyo contenido no coincide con el SHA-256;
- rechazo de una colisión divergente antes de escribir el target.

Alcance: demuestra el compositor y su fail-closed básico. No demuestra todavía una composición empresarial de dos packs reales, porque PostgreSQL/PBC siguen candidatos y Go todavía es incompatible con la foundation SQL. Esa condición permanece abierta.

Tres directorios de reconstrucción bajo `%TEMP%` no pudieron eliminarse: la política de ejecución bloqueó tanto `Remove-Item -Recurse` como la eliminación descendente, aun después de verificar paths absolutos contenidos en el temporal del sistema. No están en el workspace y no contienen dependencias instaladas.
