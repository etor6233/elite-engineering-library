# V332 — publicación y reproducibilidad de copia portable

Mantenimiento correctivo T2808/T2810, 2026-09-08. Dos defectos reproducidos y
corregidos en CREATE_PORTABLE_ARCHIVE.ps1; no nuevo runtime ni código público
adquirido. El artefacto final de la biblioteca sigue pendiente.

## Defectos observados

FAIL544: se ejecutó el bloque de publicación anterior sobre fixtures .bin. Un
checksum creado después del precheck fue sobrescrito por WriteAllText. Cuando
el destino del checksum era un directorio, el script falló después de mover el
archivo final y dejó salida parcial. publication-before.json conserva ambos casos.

FAIL545: dos README sintéticos con SHA idéntico, mtimes2000/2002 y el algoritmo
Compress-Archive anterior produjeron ZIP diferentes. El manifest de bytes no
capturaba esa dependencia de fechas; determinism-before.json conserva hashes.

## Correcciones canónicas

Publish-VerifiedArchive mantiene el candidato abierto sin compartir escritura,
calcula su hash, reserva checksum y archivo con CreateNew/FileShare.None, copia
y fuerza flush de ambos. Las colisiones conservan el destino existente y las
excepciones normales retiran únicamente las salidas propias. El candidato no se
mueve ni se borra por esta función; el staging de la invocación conserva su cleanup.
Tres parejas de publicadores con contenido diferente producen un ganador y un
rechazo, con checksum correspondiente a los bytes ganadores.

Write-DeterministicArchive aplica el patrón AUTHORED ya existente del pack
PORTABLE-SIGNED-RELEASE-EVIDENCE-GATE0.1.0: NoCompression, orden ordinal, timestamp
UTC explícito y atributos externos fijos. Manifest ordenado ordinalmente.
SourceDateEpoch predeterminado946684800 (2000-01-01), parametrizable sólo con un
segundo par representable por ZIP (1980–2107). NoCompression aumenta el tamaño;
no se afirma reproducibilidad entre runtimes/OS no probados.

README explica el parámetro, tamaño, verificación del par y límite de firma.
Los dos scripts modificados son AUTHORED; el patrón del pack es autoridad local,
no código atribuido a un upstream nuevo. VERIFY_LIBRARY.ps1 permanece idéntico.

## Evidencia

50checks del script canónico, luego50sobre copia limpia de3archivos idénticos:
32selectores de distribución,8publicación/rollback,3parejas concurrentes,
1reproducibilidad con orden/metadata/payload y6rechazos ZIP. Las dos inyecciones
de fallo sustituyen sólo la llamada de copia o de escritura del checksum por un
throw; las reservas y cleanup son el cuerpo canónico. Candidatos/originales
quedan intactos. Se conserva la fuente previa y cada fixture.

Dos ejecuciones completas del empaquetador sobre un mini-payload sintético
generan ZIP byte-idénticos pese a mtimes distintos, manifests internos válidos
y checksum lateral correcto. La dependencia VERIFY_LIBRARY fue sustituida por
un stub rotulado para aislar ese fixture: esto NO es validación del payload
completo de biblioteca ni un release. El verificador real se ejecuta aparte en
el cierre integrado. No se reemplazó el verificador en la fuente canónica.

## Límites de cierre

Proceso/OS interrumpido puede dejar par parcial. Consumidores deben verificar
checksum y manifest; no se afirma una transacción atómica de dos archivos ante
power loss, seguridad de directorio hostil ni firma/autorización de release.
Los ensayos .bin/ZIP contienen fixtures sintéticos; no se distribuye una copia
completa ni se crea una identidad de firma. TEST06/07/08 continúan requiriendo
el artefacto final y sus gates; TEST04 recibe esta evidencia focal adicional.
Se mantienen48tests/7pendientes, readiness42,10macrofrentes y FAIL532 nativo.
Consulta de historial crudo enviada al usuario: respuesta aún no inferida.

## Reproducción y hashes

Ejecutar markdown_system/test_library_maintenance_state.ps1 con PowerShell7.6.5.
Fixture retenido en el temporal impreso por el comando. Sólo extrae funciones y
declaraciones reales mediante AST; las pruebas del publisher no firman releases.
Stage de esta observación: elite-v332-a64048e690cc44e8aebbb56c17729445, temporal.

| Evidencia del stage | SHA256 |
|---|---|
| publication-before.json | `54291c5a4f0e5bac70f339e22f1ac8d0340ce88adc6225678dc7681bb827c2d0` |
| determinism-before.json | `820c5363d6fcecd355028d5dac58c7e8a029ce419eccd44f5ac403bdc649b1f0` |
| publication-after.log | `48504ab5a59b75e812d294cc86901f4b9f84332d41ca78455e86dcc0d7aa43e4` |
| archive-after.log | `a7a5105eba12a1c7f8bbd34ab0d38ea9f52a80ef85b3576df3e5c598154c658c` |
| archive-rebuilt.log | `56f12a0596d1d1ac1a3fecdd0b113529d0259929c9c72fdadd1b2a90827d80aa` |
| main-1.log | `2682dbe8c893f5431c1f74fdce1fe5e4225725d0dcd1bfbc5c419c67079320bb` |
| main-2.log | `3cfbb8d94e2af079014a989adc0c352789edf55fc18a27bbcd9963a8956f706f` |
| synthetic-main-1.zip | `db1815e8fe2e211c6f3a6e2f9d9874f211484bc3264162b18cf3a145cb3728e6` |
| synthetic-main-2.zip | `db1815e8fe2e211c6f3a6e2f9d9874f211484bc3264162b18cf3a145cb3728e6` |

| Fuente reconstruida | SHA256 |
|---|---|
| CREATE_PORTABLE_ARCHIVE.ps1 | `12d4973627a5771193bb4a2d77f390f186422944097f5ebbc27792dac48f9b0c` |
| VERIFY_LIBRARY.ps1 | `f73a6c09279e195d645d443c4d15d910f0aee89023b80592980fc25856152cd5` |
| markdown_system/test_library_maintenance_state.ps1 | `9750688f7c2697e881c615603321dcc28cb57c438a1cea5711d70b66773f985b` |

Checkpoint98 prepara el gate integrado; no claim de cierre global.

## Cierre integrado

Preflight98 ejecutó151pasos PASS. VERIFY_LIBRARY_PASS161packs/1453archivos/
768Markdown/52perfiles. Disponibilidad general BLOCKED sólo por Docker ausente;
las lanes marcadas SKIPPED o BLOCKED no pasan por el resultado estructural.
Contrato plan PASS; evidence integral sigue bloqueado por sus gates pendientes.
Checkpoint98/resume98 PASS con113evidencias y98eventos. Checkpoint99 conserva
ese historial y registra este cierre sin cambios de código posteriores.

FAIL544/545 quedan REGRESSION_PROVEN para publicación ante excepciones normales
y reproducibilidad del fixture. TEST06/07/08, admisión nativa FAIL532, readiness
y alcance global no se cierran por estas pruebas. No se creó archivo completo
de biblioteca ni firma.7de48filas de verificación siguen pendientes; no hay una
medida exacta de porcentaje global restante.

| Cierre del stage | SHA256 |
|---|---|
| preflight98.json | `f59e1e242ce87df29daa406d7114929083d5954269cf63b9f36b985844c7c49d` |
| preflight98.log | `9b3f1946d1a69aee02b0a49cb2d2091e583e620cab993d251ac021f7e95b3be6` |
| contract-plan.log | `2fd65a7ad27c8b968c64f8079abf8610894d6e7732b0b0f514dba962cf93b565` |
| contract-evidence.log | `4e5a66765adc406a0929c60445312e6b4bc8f234fb3f8ee34568a7808a5d7243` |
| checkpoint98.log | `837054bc9e630142ea08d903dcece74ac8a1c527728b6dc002ec091b9f43b30b` |
