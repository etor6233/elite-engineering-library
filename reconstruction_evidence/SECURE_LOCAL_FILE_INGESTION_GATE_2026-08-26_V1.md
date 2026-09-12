# Secure Local File Ingestion Gate — 2026-08-26 V1

## Resultado

Se materializó y ejecutó una puerta local única sobre tres motores públicos oficiales fijados. La integración contiene seis archivos `AUTHORED`; no copia código de Google, Cisco Talos ni VirusTotal y no presenta la coordinación como si fuera publicada por esas empresas. Los motores ejecutados fueron sus binarios release exactos.

La capacidad queda `REBUILD_VERIFIED / CONDITIONED`: se puede reconstruir y ejecutar inmediatamente después de que el proyecto complete/apruebe política, acepte licencias, actualice bases, aporte ruleset revisado y configure aislamiento. `ADMITTED` significa ausencia de matches en esos detectores/configuración/momento, no seguridad universal ni exactitud documental.

## Identidades upstream ejecutadas

| Motor | Identidad | Artefacto Windows x86-64 | Hash del ejecutable observado | Contrato observado |
|---|---|---|---|---|
| Google Magika | CLI 1.1.0, commit `5e2f437fb7b7452368c8c1fa9354858f5487a5c4` | ZIP 10.076.032 bytes, SHA-256 `d4de347c53e2d25f9663780c93fdf26d156284c64d6a5ec4bcdb829f9994f6f3` | `3631dab2f57ec42b6646ce141397671132707cd7e7633f4511fe47171efe69eb` | `magika 1.1.0 standard_v3_3`; `--json` entrega status/output label/MIME/score |
| Cisco Talos ClamAV | 1.5.4, commit `fa59fca15872bb8a914ba4c68188bcc8a502cbdf` | ZIP 225.774.978 bytes, SHA-256 `0d9e0228b2674137ea1a2853566c98a0278ad52ab2582c3d6dbd75373848c395` | clamscan `f368912f61ddea8b302acf9114891ec57f8a45b0de77f076b1d19c960bc6f7cf`; freshclam `4032fdd45184333d7c963eec87664a8339c6f31dc5ae94bb5008fe5e88feaf1f` | exit 0 sin detección, 1 detección, 2 error; `official-db-only`, edad y límites disponibles |
| VirusTotal/Google YARA-X | 1.20.0, commit `60ad06971467029e77967e59d580cbbe85a1474d` | ZIP 8.640.562 bytes, SHA-256 `b1e2840bac593aea353d2b2b341f5a862c9d61c0c406d9abbbad9e1fa35163a1` | `d3e648651f3eaa5e833faeb394fc7dc7dfb583a31dd5fc572ab6cae7ac41f3fd` | `scan --compiled-rules --output-format=json`, timeout/matches/threads/no-mmap; match y no-match devuelven exit 0 y se distinguen por JSON |

Licencias upstream: Magika Apache-2.0; ClamAV GPL-2.0-only; YARA-X BSD-3-Clause. Las identidades completas, LICENSE hashes, source archives y locks permanecen en `OFFICIAL-UPSTREAM-ACQUISITION-CORE` y en las evidencias V1 previas de file security.

Fuentes oficiales consultadas: https://github.com/google/magika/tree/5e2f437fb7b7452368c8c1fa9354858f5487a5c4, https://docs.clamav.net/manual/Usage/Scanning.html, https://docs.clamav.net/manual/Usage/SignatureManagement.html, https://github.com/Cisco-Talos/clamav/blob/fa59fca15872bb8a914ba4c68188bcc8a502cbdf/docs/man/clamscan.1.in y https://virustotal.github.io/yara-x/docs/cli/.

## Bases ClamAV reales

`freshclam` actualizó y probó las bases oficiales el 2026-08-26:

- daily 28104, `daily.cvd` 23.428.985 bytes, SHA-256 `1a04a59df8aebcc6d3a9340846479e732579fc936db92f46867061a033aee717`;
- main 63, `main.cvd` 89.072.577 bytes, SHA-256 `0b2182d229f46981ec8f535382222f7c9dfdd656b250ad47988b910a8d302365`;
- bytecode 339, `bytecode.cvd` 281.702 bytes, SHA-256 `6d4aa01f219e988060fc419f495d07f27e0cdf1a2cccc065971da922c76f7ffb`.

Los tres self-tests de base pasaron. Después, el preparador emitido ejecutó `clamscan --database ... --official-db-only=yes --fail-if-cvd-older-than=2` sobre un probe benigno y sólo entonces escribió receipt/inventario hash. Esto evita depender de `freshclam` exit 0: el issue oficial https://github.com/Cisco-Talos/clamav/issues/965 documenta un caso donde puede devolver cero bajo cool-down sin producir una descarga utilizable (`UP-FAIL-049`).

No se atribuye una detección EICAR a ClamAV: Microsoft Defender interceptó ese fixture en la auditoría anterior (`LIB-FAIL-087`). Este ciclo demuestra carga/frescura/base/probe benigno, mientras la detección EICAR sigue siendo gate obligatorio en CI aislada del proyecto.

## Ruleset YARA-X real

El binario exacto compiló `foo.yar` del commit oficial como fixture de contrato y produjo un `.yarc` de SHA-256 `6923c4ef5d8829c1bb8dc8eca5807aab58aee839efb17a8d55f433135b0c832a`. El scan JSON real de un Markdown autorizado produjo cero matches. La biblioteca no entrega este fixture como ruleset de seguridad: cada proyecto debe aportar reglas con procedencia, revisión, hash, falsos positivos, canary y rollback.

## Recorrido integral real

Sobre un Markdown local autorizado:

```text
ClamAV official-db-only + max-age + limits: exit 0
Magika: markdown / text/markdown / score 0.9929999709129333
YARA-X compiled rules JSON: 0 matches
SECURE_LOCAL_FILE_GATE decision=ADMITTED reason=ALL_SELECTED_GATES_PASSED
```

El receipt contiene SHA-256/bytes del input, SHA de política, hashes/versiones de los tres binarios, versión de bases, resultado de tipo, hash de ruleset y reglas coincidentes; omite paths y output crudo. Conserva `business_storage_authorized=false`.

## Verificación reconstruida

- seis bloques materializados con hashes exactos;
- `compileall` PASS;
- once tests unitarios PASS;
- negativos: política pendiente, binario alterado, malware/detección, error/base vieja, tipo no aprobado, match YARA, overwrite y FreshClam no probado;
- preparador ClamAV real PASS;
- compilador YARA-X real PASS;
- recorrido integral real PASS;
- artefactos temporales no forman parte de la distribución.

Fallos propios incorporados: `LIB-FAIL-133` a `LIB-FAIL-136`. Condición upstream incorporada: `UP-FAIL-049`.

## Condiciones que siguen siendo del proyecto

El usuario/agente debe completar canales, tenants, tipos, límites, cuarentena/clean/retención, aislamiento, firma máxima, outage/HA, ruleset, falsos positivos, EICAR/hostiles/polyglots/bombs, carga, métricas, recall e incident owner. Luego debe seleccionar conversor y analyzer/processor oficial, schema/corpus/ground truth y aceptación por campo. Ninguna puerta genérica puede prometer cero errores para todos los documentos.
