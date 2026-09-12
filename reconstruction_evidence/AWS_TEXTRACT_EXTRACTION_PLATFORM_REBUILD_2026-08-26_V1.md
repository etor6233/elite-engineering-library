# AWS Textract Document Data Extraction Platform — rebuild 2026-08-26 V1

## 1. Fuente exacta y criterio

Se reaudito el source publico oficial `aws-samples/aws-textract-document-data-extraction-platform` sin corregirlo ni relabelar glue local como codigo AWS.

| Propiedad | Valor verificado |
|---|---|
| commit | `c7fed5353abcae184fbd705d5cf6a8512af0e9c4` |
| archive | 1.124.136 bytes; SHA-256 `52755b28d8afd2c8cee23ecba76b11c364edd3ecad066b2e4a51eb2f77219bb1` |
| licencia | Apache-2.0 |
| archivos source | 281 |
| `pnpm-lock.yaml` | SHA-256 `240755c9c46f7d8e41985f50d141c7f2f93d15bbae461c595f21cd655b147e22`; lockfile 6.0 |
| engines | Node `>=16`; pnpm `>=8 <9` |

El archive exacto sigue siendo la cabeza oficial segun la comprobacion de frescura V1. Esto prueba identidad, no readiness.

## 2. Toolchains aislados

No se cambiaron runtimes globales ni se uso Git. La repeticion valida uso:

- Node 18.20.8 y su `SHASUMS256.txt` oficial; tar Linux SHA-256 `5467ee62d6af1411d46b6a10e3fb5cacc92734dbcef465fea14e7b90993001c9`;
- pnpm 8.15.9, dentro del rango declarado;
- Python 3.11.9 portable adquirido mediante Astral `uv` 0.11.30;
- Poetry 1.8.3, version que genero el lock principal de `packages/lib`;
- Amazon Corretto 17.0.20.10.1, descargado desde el endpoint estable AWS y verificado contra su sidecar oficial: Windows ZIP SHA-256 `af002c5f7dd3d09ad4a1f1643266e28fb368077de50560534a2aa84e376ddfde`; Linux tar SHA-256 `74ff458657da91ca222681993e3c6b9a8e3629ca8e61c0d8cd90527280da9aa5`;
- Google OSV-Scanner 2.5.1, commit `c84fa4568f2526d0333e9a914ea8a0a5f74ad68b`, binario Linux verificado contra `osv-scanner_SHA256SUMS`.

Los temporales no contenian credenciales AWS. No se desplego infraestructura ni se llamo Textract.

## 3. Instalacion: termina, pero no es frozen

`pnpm install --frozen-lockfile` termino cero con Node 18/pnpm 8 tanto en Windows desde un path de 35 caracteres como en WSL2/Linux. El lock pnpm permanecio byte-identico, pero el `postinstall` oficial ejecuta el target Python `install`, no `install:ci`; ese target contiene `poetry update` en tres proyectos.

Locks originales:

| Lock | SHA-256 original |
|---|---|
| runtime Python generado | `2af7b66f09672aabeec43beb4457239ce25558d8e574f90fa2256a838f6a55b9` |
| handlers Python | `5d732ef4eda13c8ddb5eb7f8d5e52f6c7b4a64bdac2e97c04a522f07e120f603` |
| libreria | `8aca277ba6b1a0e8fb32742eb1adb6149d48b6edd2731369d5b06b604ba95490` |

Despues de la instalacion Windows pasaron a `956f97e2...`, `0a9d589b...` y `bf32c424...`; la repeticion Linux resolvio otro conjunto y termino en `36d5ecbf...`, `87141f77...` y `67fe5e4d...`. Por tanto, el comando documentado no reconstruye el grafo Python fijado: consulta indices actuales y reescribe los locks aun cuando la invocacion superior declara `--frozen-lockfile`.

## 4. Build y tests oficiales

### Windows

La instalacion solo fue posible desde un path corto; paths mas largos fallaron por el limite de 260 caracteres en MSBuild/virtualenv. Luego `pnpm build` termino 1: los tasks usan `$AWS_PDK_VERSION`, que `cmd.exe` no expande, y `npx` rechazo el package spec resultante. No se altero el task para hacerlo pasar.

### Linux/WSL2

El mismo commit supero el primer codegen, pero el build integral termino 1. Los subproyectos Python ejecutan `pnpm dlx projen build` sin version. En la fecha de corte `dlx` resolvio Projen 0.103.5, mientras el lock oficial contiene Projen 0.82.8. La version nueva interpreto el task generado `cp -f ...` como `cp: unsupported flag: -f`; seis proyectos posteriores no se ejecutaron. El build, ademas, instalo dinamicamente 1.352 paquetes bajo `.pdk`, fuera del lock raiz, reporto 34 subdependencias deprecadas y dos conflictos de peers.

`pnpm test` tambien termino 1 por la misma frontera de `dlx`: se observaron un test web `placeholder` PASS, un test Python de libreria PASS y cero tests colectados en handlers; luego Projen 0.103.5 rechazo construcciones shell `if` y `|| ([ ... ])` guardadas por el source. Dos proyectos finales no se ejecutaron. No existe suite integral verde en esta reconstruccion.

## 5. Seguridad del grafo exacto

Sobre el `pnpm-lock.yaml` original, `pnpm audit --prod --json` detecto 1.485 dependencias productivas y 153 ocurrencias: 18 low, 65 moderate, 65 high y 5 critical, agrupadas en 147 advisories.

OSV-Scanner 2.5.1 recorrio los cuatro locks originales: 2.334 packages npm, 52+33+15 packages PyPI. Resultado fechado:

- 91 registros de package afectados;
- 75 registros npm y 16 PyPI;
- 218 IDs de vulnerabilidad unicos: 171 npm y 47 PyPI.

Entre los findings hay criticidad para `sha.js`, `form-data`, `fast-xml-parser`, `shell-quote` y `tar`; el conteo no sustituye triage de reachability, pero impide promocion inmediata.

## 6. Decision

La plataforma aporta codigo oficial real para UI, APIs, correccion humana, metricas, CDK y Textract. No esta lista para copiar a un sistema actual. Queda `SAMPLE_ONLY / REJECTED_FOR_IMMEDIATE_ADOPTION` hasta que una revision oficial:

1. elimine `poetry update` del camino frozen y publique locks reproducibles;
2. fije Projen/codegen y todo grafo dinamico;
3. pase build y suite integral desde cero en una plataforma declarada;
4. elimine findings criticos/altos o publique triage y fixes verificables;
5. demuestre despliegue, seguridad, corpus empresarial, recovery y costo del target.

No se incorpora un parche Elite para ocultar estos defectos. Las condiciones quedan retenidas en `UP-FAIL-054` a `UP-FAIL-058`; los errores de la auditoria y sus regresiones en `LIB-FAIL-171` a `LIB-FAIL-181`.

Fuentes oficiales: <https://github.com/aws-samples/aws-textract-document-data-extraction-platform/tree/c7fed5353abcae184fbd705d5cf6a8512af0e9c4>, <https://docs.aws.amazon.com/corretto/latest/corretto-17-ug/downloads-list.html>, <https://nodejs.org/download/release/v18.20.8/> y <https://github.com/google/osv-scanner/releases/tag/v2.5.1>.
