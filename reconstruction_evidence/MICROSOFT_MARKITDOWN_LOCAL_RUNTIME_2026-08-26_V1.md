# Microsoft MarkItDown Local Runtime — 2026-08-26 V1

## 1. Claim probado

Se materializó y ejecutó un runtime local restringido sobre el código oficial Microsoft MarkItDown 0.1.7. Convierte PDF, Word, Excel, PowerPoint, Outlook MSG y formatos textuales mediante la API oficial `convert_local()`, sin habilitar URL, plugins, archives, audio/video, YouTube, LLM ni almacenamiento automático de datos de negocio.

El pack nuevo es `MICROSOFT-MARKITDOWN-LOCAL-RUNTIME 0.1.0`, seis archivos. El engine/converters son Microsoft; wrapper, profile, lock y tests son `AUTHORED` y no se atribuyen a Microsoft.

## 2. Identidad oficial

```text
repository: microsoft/markitdown
release: v0.1.7
commit: fd239d5d2be43d9b68329730206b9312c7d5a388
commit verification: verified=true, reason=valid
commit date: 2026-07-29T18:16:03Z
source archive bytes: 4,021,911
source archive SHA-256: 13b3d78a0ba3807df85208f662e89a2155b0137b1d3bc32f9635d03f329e1c70
source LICENSE: MIT
source LICENSE SHA-256: c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383
PyPI wheel: markitdown-0.1.7-py3-none-any.whl
wheel bytes: 71,093
wheel SHA-256: 4eca912c87c6aa6897284a7f4bf6769a23bccf8544530f5d8b175fbe3797c916
wheel upload: 2026-07-29T18:20:30.226447Z
wheel Trusted Publishing: no
```

La URL del wheel, tamaño y digest se obtuvieron del JSON oficial de PyPI y el archivo descargado coincidió. No se infirió una URL. El source ya estaba fijado por el source lock de 68 upstreams.

## 3. Selección de superficie

Microsoft publica extras separados. Se seleccionaron exclusivamente:

```text
pdf, docx, pptx, xlsx, xls, outlook + core
```

Se excluyeron deliberadamente:

- `az-content-understanding`: MarkItDown 0.1.7 pide `azure-ai-contentunderstanding>=1.2.0b1`, mientras la biblioteca usa el SDK GA 1.1.0 por una lane separada y verificada; no se mezcló beta con GA;
- audio/YouTube/media: requieren surfaces y toolchains ajenos al claim;
- archives: la ingesta segura exige una lane de extracción aislada específica;
- URL y plugins: Microsoft advierte que MarkItDown accede recursos con los privilegios del proceso;
- LLM: conversión local no requiere inferencia ni costo cloud.

El grafo resuelto para CPython 3.12 Windows x86-64 contiene 44 wheels y 82.672.112 bytes. Cada fila fija `name==version` y SHA-256; `pip install --require-hashes --only-binary` puede reconstruirlo sin source builds. Los license expressions/classifiers y license files de los wheels se inspeccionaron; la aprobación legal del artefacto final permanece como gate del proyecto.

## 4. Seguridad y vulnerabilidades

Una consulta batch oficial OSV de los 44 pares PyPI nombre/versión devolvió cero findings conocidos el 2026-08-26. La primera interpretación del JSON produjo 42 filas nulas por semántica `@($null)` de PowerShell; se registró `LIB-FAIL-130`, se exigió ID no vacío y se repitió la consulta correctamente.

El wrapper aplica:

- `convert_local()` solamente;
- archivo regular local y no symlink;
- allowlist de extensiones;
- bloqueo de archives, media y executables;
- 50 MiB de input y 100 MiB de output por default;
- plugins deshabilitados y ausencia de argumentos URL/secret/LLM;
- staging nuevo y rename final;
- hash SHA-256 de input/output y receipt sin ruta original;
- `business_storage_authorized=false` inmutable.

El proceso de despliegue todavía debe negar egress, aplicar límites/timeout y ejecutarse después de type/antimalware/quarantine. MarkItDown convierte estructura a Markdown; no prueba campos comerciales.

## 5. Reconstrucción ejecutada

```text
wheelhouse exact: 44 wheels / 82,672,112 bytes
offline install into empty CPython 3.12 venv: PASS
pip check: PASS — No broken requirements found
markitdown import/version: PASS — 0.1.7
environment verifier: PASS — distributions=44
wrapper unit regressions: 5 PASS
pack materialization: 6 files, source tree hash-equivalent PASS
OSV: 44 packages / 0 findings known
```

El entorno verificó además `olefile==0.47` y `xlrd==2.0.2` para MSG/XLS. No se usó una instalación global ni se descargaron dependencias durante los tests offline.

## 6. Fixtures oficiales convertidos

Todos pertenecen a `packages/markitdown/tests/test_files` del commit exacto:

| Fixture | Input bytes | Input SHA-256 | Output bytes | Output SHA-256 |
|---|---:|---|---:|---|
| `test.pdf` | 92.971 | `77c4014cf15ea5e56663e9b8093ddc2a89a9224c15d6a9eb2aa8af2b48e1b237` | 5.200 | `911930696aa4cea4a0fd27f6f64eeeb6589c699397470ef07214fe1218b86fe1` |
| `test.docx` | 135.824 | `ee1974633f3b1e8bb54201abce779a755cb1cd5259641a08cf7b1b9b36dec42b` | 4.703 | `dd23168c57aaff38ea627ca03a6ebb368d543daf1f772d5e4ee575b27efc24e2` |
| `test.xlsx` | 11.562 | `a867be6ece38a4224bcfe34312e5c296ed520a76a431e0454603ceb171fbd4af` | 808 | `346eee409feebbca0bf6bf144f655c851045eb77b16fec4d873d06e8c7df8caf` |
| `test.pptx` | 277.515 | `f0b9e5252aec3730c91f276a7c3b8f4a9893a7b66540f6c6af160a49855ca507` | 2.048 | `c56bfb7203a996f609f14bb04659cc7e45ba39990c1af4dbbc3df23193090fcd` |
| `test.xls` | 27.648 | `17a94b6514e8998f4dc25bc77265b6f62982c18614ba4401a87fa01f90f53f1d` | 808 | `346eee409feebbca0bf6bf144f655c851045eb77b16fec4d873d06e8c7df8caf` |
| `test_outlook_msg.msg` | 13.312 | `028d84ffe67e1865009669d13d4c12682943b32eccf7f84a8da1899db63b0131` | 173 | `1fdb3db829ce3de0cb6c9a9702ea8c14dfd1c6447ede0395eef9fc3b97deb982` |

El wrapper materializado se ejecutó además sobre `test.pdf` y produjo exactamente el mismo output SHA-256 `911930…`, con receipt atómico y `business_storage_authorized=false`.

## 7. Fallos propios convertidos en lecciones

- `LIB-FAIL-128`: `uv` fue asumido sin estar en PATH;
- `LIB-FAIL-129`: se intentó una URL PyPI inferida y el objeto no existía;
- `LIB-FAIL-130`: parser OSV contó filas nulas;
- `LIB-FAIL-131`: el primer intervalo dejó un venv parcial y `pip check` solo no detectaba ausencia de MarkItDown;
- `LIB-FAIL-132`: la sesión larga perdió su ID y dejó padre/hijo bloqueados por output sin consumidor.

La repetición obtuvo la URL desde JSON oficial, conservó session ID, usó salida acotada, esperó exit 0 y exigió simultáneamente distribución presente, versión, 44 packages, imports, `pip check`, tests y fixtures.

## 8. Admisión

```text
source/wheel identity: PASS
Windows CPython 3.12 frozen graph: PASS
local PDF/Office/MSG conversion: PASS
wrapper atomic/security regressions: PASS
known-vulnerability query: 0 findings at dated query
untrusted-input process isolation: CONDITIONED_BY_TARGET
Linux/macOS lock: NOT_CREATED_OR_CLAIMED
semantic field accuracy: NOT_PROVEN_BY_CONVERSION
automatic business storage: FORBIDDEN
admission: REBUILD_VERIFIED / CONDITIONED
```
