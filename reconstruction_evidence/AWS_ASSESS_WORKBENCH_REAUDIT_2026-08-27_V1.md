# AWS Assess Workbench Reaudit — 2026-08-27 V1

## Resultado

Se auditó el código público oficial `aws-samples/sample-assess-workbench`, una muestra de AWS para evaluación documental con agentes, frontend de revisión, infraestructura y flujos de trabajo. El commit exacto aporta material arquitectónico relevante, pero no supera los gates de adopción inmediata: su frontend verificable funciona localmente, mientras la suite Python falla y ambos locks conservan vulnerabilidades conocidas.

La fuente entra en `OFFICIAL-UPSTREAM-ACQUISITION-CORE` únicamente como `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`. El agente puede localizarla y adquirir el archive exacto para investigación aprobada; ningún perfil puede componerla, desplegarla o presentarla como código listo.

## Identidad exacta

- repositorio: `aws-samples/sample-assess-workbench`;
- rama oficial: `main`; no existen releases publicados al corte;
- commit: `79a57b557fce417bb2bbbcc0794a8f7a6ce9245f`;
- fecha del commit: 2026-08-18;
- tree: `aa8403a00b981a8a8a9cd5f8ea3249f803a0c204`;
- el commit no figura firmado;
- archive por commit: 1.665.545 bytes, SHA-256 `cab872e5919eadffac962527f9347fc098badbf875343aa179462c31467a966a`;
- licencia MIT-0: SHA-256 `5d0ed32d916046277a19e21df9cabc50e88de80e62bfff7f4df38b89800355ef`;
- README: SHA-256 `8d104962327b09bb843ec5ab9c11fd7ed083e026ed238632acdc83308fac4592`;
- `pyproject.toml`: SHA-256 `38b30f3848de1b41d54a8f60c6af5588b4f70b4fa0efec46efe513766e2f6ac0`;
- `uv.lock`: SHA-256 `27c9fc24d1494e1e7eabfa6ac3cc2061daebb624a1849f34bf7020eef7418611`;
- frontend `package.json`: SHA-256 `2a2ba02c85d77d6e52b20f405cd732a62710f60ed11d2a7dd3ae1652d84a573d`;
- frontend lock: SHA-256 `967d6c829f8b16823d39e74507bf58069f0c9b52461974df36a4db083dfcb10a`.

## Ejecución sin modificar source

### Python

La instalación congelada mediante `uv` y Python 3.14 incorporó 47 paquetes. Ruff y `compileall` sobre `api`, `agents` y `tests` pasaron. La colección inicial en Windows falló en dos módulos porque los tests abren YAML sin declarar UTF-8. Al repetir con UTF-8 de proceso para separar portabilidad de comportamiento, pytest ejecutó 482 casos: 459 PASS, 17 FAIL y 6 skip.

Los 17 fallos no se ocultaron: uno contradice la semántica esperada de `_is_not_found`; diez crean un cliente Bedrock durante importación y fallan sin región pese a pertenecer a la categoría unitaria; seis construyen invocaciones Bash con paths Windows sin conversión o quoting correcto. Esto impide afirmar que la suite sea offline, hermética o multiplataforma.

### Frontend

Desde una raíz Windows corta y con el lock exacto:

- `npm ci`: 342 paquetes instalados;
- lint: PASS;
- Vitest: 47/47 PASS;
- Knip sobre el checkout limpio: PASS;
- build: PASS, 77 módulos; bundle principal 308,10 kB, gzip 95,02 kB;
- `npm audit`: 5 vulnerabilidades —4 altas y 1 moderada—;
- auditoría productiva: 1 moderada en DOMPurify.

La primera ejecución de Knip posterior al build señaló cuatro chunks de `dist` como no usados. Se registró como `LIB-FAIL-337`, se eliminó sólo ese output temporal y la repetición limpia pasó; no se atribuyó el falso positivo al upstream.

## Dependencias y seguridad

OSV-Scanner oficial 2.5.1, SHA-256 `25e42f5e16d3e2884ace3f40ca19e0f5c391a8aac9afde6a997045f2969bfb6`, examinó los locks exactos: 139 paquetes Python y 437 paquetes npm. Reportó 43 registros de vulnerabilidad en 9 paquetes:

- npm: `brace-expansion` 1.1.15, `dompurify` 3.4.11, `js-yaml` 4.2.0, `nanoid` 3.3.13 y `postcss` 8.5.15;
- PyPI: `bedrock-agentcore` 1.15.0, `cryptography` 49.0.0, `mcp` 1.28.0 y `pillow` 12.2.0.

El scan de lock no demuestra alcanzabilidad. Tampoco permite ignorar los hallazgos: no se ejecutaron fixes automáticos ni se alteraron versiones para fabricar un PASS.

## Límites de la prueba

Terraform y Task CLI no estaban instalados, por lo que no se validaron infraestructura ni tareas declaradas. No se usaron credenciales, no se desplegó AWS y no se realizó ninguna llamada con costo. Tampoco se demostraron aislamiento, autorización productiva, navegador/accesibilidad, carga, backup/restore, rollback ni observabilidad sin PII.

## Fallos y decisión

Las condiciones upstream quedan en `UP-FAIL-092` a `UP-FAIL-094`; los errores locales y sus correcciones están en `LIB-FAIL-335` a `LIB-FAIL-337`. Todos los fallos locales quedaron `REGRESSION_PROVEN`; la condición del componente permanece `REJECTED_COMPONENT`.

Assess Workbench contiene código oficial útil para estudiar agentes, evaluación y una interfaz documental. No es código instantáneamente reutilizable en el estado publicado. Una futura admisión requiere snapshot oficial inmutable —preferentemente release y firma—, locks corregidos, suites unitarias offline y multiplataforma totalmente verdes, infraestructura verificable, scan con análisis de alcanzabilidad y los gates productivos del target.

Fuente oficial gobernante: <https://github.com/aws-samples/sample-assess-workbench/tree/79a57b557fce417bb2bbbcc0794a8f7a6ce9245f>.
