# Project Failure Lessons

## FAIL-20260910-742 — Firefox context shutdown rejects connected regression

V377 / BROWSER_LIFECYCLE / DIAGNOSED. Connected execution passed Chromium desktop/mobile and WebKit; Firefox passed initial/read-failure/restart, then Browser.removeBrowserContext failed during operation teardown: SessionStore _maybeDontRestoreTabs referenced an undefined _windows entry. Preserve screenshot/error-context/trace and the full failed run; all owned PostgreSQL processes stopped. No92phase/whole-regression PASS or canonical SEO promotion. Inspect exact trace and prior lifecycle lessons before a focused correction/retry; do not suppress context-close errors, disable assertions, change browser version or assume the product passed from the teardown label.

## FAIL-20260910-741 — BFF dependency bill retained the vulnerable version label

V377 / FRESHNESS / DIAGNOSED. The V376 package/lock and four consumer plans correctly select Next16.3.4, but the current dependency bill still labels it16.3.2. The historical V117 build record is distinct and must be preserved. Correct only the current bill alongside the next BFF revision; verify the reconstructed package remains the already-qualified16.3.4 and do not rewrite historical evidence.

## FAIL-20260910-740 — metadata probe assumed a no-store framework contract

V377 / TEST_ORACLE / DIAGNOSED. The real sitemap response is public,max-age=0,must-revalidate, exactly the REVALIDATE header in authenticated Next16.3.4 next-metadata-route-loader.js. The provisional probe required no-store without a product requirement for that stronger policy and rejected before testing runtime origin changes. Preserve rejected runs and observed headers. The acceptance condition is current runtime metadata with mandatory revalidation: assert the exact framework header and rebuild-once/run-three-configurations behavior. Do not call this a private/no-store response or assume CDN invalidation/index removal. FAIL739 also had a transient helper indentation rejection before execution; fixed without source changes.

## FAIL-20260910-739 — HTTP probe treated header names as case-sensitive

V377 / TEST_HARNESS / DIAGNOSED. The isolated server returned a successful robots response, but converting HTTPMessage into an ordinary dict made the Content-Type lookup case-sensitive. The probe rejected before claiming any scenario PASS; its process exited and owned server stopped. Preserve first result/log and read headers through the case-insensitive HTTPMessage API. A following shell-quoted Python generator also failed parsing before file writes, repeating FAIL729; use apply_patch for the helper.

## FAIL-20260910-738 — exact optional environment input type

V377 / AUTHORING / DIAGNOSED. The new public-indexing helper passed22tests but TypeScript exactOptionalPropertyTypes correctly rejected its default input: explicit environment values can be undefined while the optional field type excluded explicit undefined. No canonical code promotion. Preserve the failed typecheck log; declare the actual string-or-undefined input and rerun tests/typecheck/build without weakening tsconfig. A follow-up log read also requested build.log before the rejected typecheck allowed build to run; no build result inferred (FAIL719 recurrence).

### FAIL-20260908-536 — ruta de máquina en reproducción portable

- **Estado:** `REGRESSION_PROVEN`
- **Observación:** VERIFY_LIBRARY rechazó ENTRY_LATENCY_AND_INTERRUPTION_V326.md por ruta local de home en el programa de reproducción incluido. El benchmark original sí terminó; no hay PASS del gate final.
- **Corrección:** ejemplo portable con ELITE_LIBRARY_ROOT/ELITE_PWSH_EXECUTABLE y Python de la invocación, conservando bytes/script/receipts originales; no debilitar el selector de distribución.

- **Cierre V326:** ejemplo portable con variables explícitas/sintaxis válida y VERIFY_LIBRARY_PASS161/1453/761/52 en checkpoint83; rechazo inicial conservado.

### FAIL-20260908-535 — cursor conserva el sender ya corregido

- **Estado:** `REGRESSION_PROVEN`
- **Observación:** el cursor79 seguía listando FAIL457/defectos del sender genérico, aunque el ledger registra cierre V312 del scope inventariado.
- **Verificación:** reconstrucción actual del portal20/20 byte-idéntica a V321 y ausencia del sender/mensaje obsoletos; V312 conserva92fases conectadas y límites. Se retira sólo el blocker obsoleto y su ID; T2804 permanece pendiente. Evidencia FIREFOX_LIFECYCLE_AND_CONTINUITY_V325.md.

### FAIL-20260908-534 — resumen de readiness desactualizado

- **Estado:** `REGRESSION_PROVEN`
- **Observación:** la tabla Markdown mantenía C/PENDING aunque el JSON y la respuesta V287 dicen ANSWERED. Ambos resúmenes de blockers conservaban FAIL386 abierto después del cierre focal V292. El bloqueo nativo actual FAIL532 no estaba reflejado en open_dependency_blockers.
- **Corrección:** conciliar sólo hechos con owners existentes, preservar before y ejecutar el validator real. Mantener D–H pendientes, FAIL384/385 y los límites del logger CANDIDATE; no convertir corrección documental en readiness aprobada.

- **Verificación:** validator0.6.1 exit2/BLOCKED; antes41observaciones, ahora42. Único delta: blocker de dependencia antes omitido; ninguna condición retirada del reporte. Before/after y hashes en UTF8_CHECKPOINT_RECOVERY_V324.md.

### FAIL-20260908-533 — lectura de JSON UTF-8 con codepage implícita de Windows

- **Estado:** `REGRESSION_PROVEN`
- **Observación:** checkpoint78 rechazó open.blockers[4] por longitud445. Scripts locales de checkpoint usaron Path.read_text() sin encoding y se ejecutaron sin -X utf8; cada reserialización UTF-8 amplificó el mojibake de textos españoles. El validator canónico detectó el límite; no se emitió EVT0078. También debe inspeccionarse PROJECT_ENGINEERING_CONTRACT.json y cualquier JSON reescrito por esos scripts.
- **Contención:** preservar bytes before, usar -X utf8 y encoding explícito en lecturas/escrituras. Restaurar sólo transformaciones reversibles comprobadas contra sus bytes/origen; no cambiar IDs, hashes, decisiones o historial de eventos. No repetir scripts mutantes ni declarar válido el checkpoint fallido. Añadir regresión de texto no ASCII y devolver la lección a la fuente canónica.

- **Cierre focal V324:**209campos recuperados mediante inversión estricta;77eventos intactos;18scripts corregidos. Kit canónico15archivos,15tests únicos en ambos modos UTF-8. Evidencia portable UTF8_CHECKPOINT_RECOVERY_V324.md. El checkpoint78 se emite sólo después del gate; su intento previo fallido permanece registrado.

### FAIL-20260908-532 — cobertura del lock no acredita el runtime nativo distribuido

- **Estado:** `OPEN`
- **Observación V323:** reflink0.1.19 tiene75paquetes Cargo exactos y0avisos OSV; cuatro addons incluyen rustc f6e511ee (Rust1.82.0), con ocho statements npm/SLSA verificados y nueve negativos. Eso no demuestra un build reproducible o SBOM del binario. Los trees exactos reflink/fastlist no contienen archivo de licencia completo; metadata/README sólo declaran MIT. Fastlist0.3.0 no aporta proyecto/build/CRT exactos ni hashes publicados en sus assets. Los cuatro PE carecen de Authenticode y version-resource, lo cual no se confunde con falta de firma del contenedor pnpm.
- **Contención:** mantener la admisión integral nativa/licencias de redistribución BLOCKED; no llamar limpio al binario por el lock ni inventar licencia/copyright/versión de CRT. No ejecutar nuevos binarios, descargar otro grupo sin perfil o alterar pnpm para ocultar componentes.
- **Reapertura:** evidencia upstream de licencias completas y build/SBOM/stdlib/CRT vinculados a los hashes exactos, o alternativa compatible que cierre sus propios gates. Conservar revisión oficial actual y resultados focales; continuar los owners operativos sin promoción.

### FAIL-20260908-531 — conteos de perfiles dependientes no actualizados

- **Estado:** `REGRESSION_PROVEN`
- **Observación V322:** el preflight74 superó notices y checkpoint pero VERIFY_LIBRARY rechazó AMAZON_SPAPI_SHIPPING_TRACKING_PACK_PLAN: expected48/actual54, porque el owner de adquisición añadió seis archivos. No hay receipt final ni PASS integrado.
- **Corrección:** revisar las30selecciones reales del owner; ajustar sólo los presupuestos exactos de perfiles registrados que seleccionan esos seis archivos, conservando igualdad y conteos estrictos. Registrar before y delta comprobado, volver a checkpoint y ejecutar gate completo.

- **Recurrencia de diagnóstico V322 / FAIL473:** se intentó leer un directorio V305 no inventariado. La lectura falló; no se usó ni modificó PostgreSQL. Resolver futuros paths desde inventario literal o el cursor, sin construir UUIDs.

- **Cierre focal V322:** owner0.4.78 reconstruido37/37,17profiles/7negatives,9lock negatives,65transport assertions;151preflight steps y161/1453/757/52 PASS. Adquisición real3sources y10outputs inmutables; ver GOVERNED_MAINTENANCE_ACQUISITION_V322.md. Las condiciones productivas permanecen abiertas.

### FAIL-20260908-530 — inventario de provenance pendiente en notices

- **Estado:** `REGRESSION_PROVEN`
- **Observación V322:** después de validar checkpoint73, VERIFY_LIBRARY rechazó THIRD_PARTY_NOTICES por conservar el conteo anterior. El conteo real esperado es AUTHORED1209/ADAPTED137/VERBATIM107/TOTAL1453.
- **Corrección:** actualizar el owner de notices con los tres nuevos archivos AUTHORED y los tres sidecars ADAPTED, preservando los términos y bytes de las licencias decodificadas; checkpoint antes de repetir gate. No relabelar como VERBATIM los envelopes ni modificar el verificador.

- **Cierre focal V322:** owner0.4.78 reconstruido37/37,17profiles/7negatives,9lock negatives,65transport assertions;151preflight steps y161/1453/757/52 PASS. Adquisición real3sources y10outputs inmutables; ver GOVERNED_MAINTENANCE_ACQUISITION_V322.md. Las condiciones productivas permanecen abiertas.

### FAIL-20260908-529 — auditoría integral lanzada con checkpoint anterior

- **Estado:** `REGRESSION_PROVEN`
- **Observación V322:** VERIFY_LIBRARY rechazó evidence hash mismatch for failures porque los owners cambiaron después del checkpoint72. La auditoría no se declara PASS.
- **Corrección:** capturar el delta comprobado en un nuevo checkpoint DISCOVERY/BLOCKED antes de repetir la auditoría integral; conservar verify-library.log rojo. No eliminar ni relajar la validación del cursor.

- **Recurrencia V322 / FAIL473 (diagnóstico de rebuild y mapa):** se supusieron COMPOSITION_RECORD.md y CAPABILITY_GAP_RESOLUTION_MASTER_MAP.md. El inventario materializado y el cursor mostraron MATERIALIZATION_RECORD.md y markdown_system/FRANCHISE_GAP_MAP.md. La escritura se detuvo en la aserción del path inexistente; se continuó sólo con el owner pendiente, sin repetir los tres append ya realizados. No se atribuye lectura ni modificación al path inexistente.

- **Cierre focal V322:** owner0.4.78 reconstruido37/37,17profiles/7negatives,9lock negatives,65transport assertions;151preflight steps y161/1453/757/52 PASS. Adquisición real3sources y10outputs inmutables; ver GOVERNED_MAINTENANCE_ACQUISITION_V322.md. Las condiciones productivas permanecen abiertas.

### FAIL-20260908-528 — licencia de artefacto no roundtripea como texto Markdown

- **Estado:** `REGRESSION_PROVEN`
- **Observación V322:** el precheck de integración rechazó la licencia completa de Node porque sus bytes no son LF terminados en newline. El pack canónico no se escribió; el before y el candidato quedan preservados. Normalizarla invalidaría su SHA y el claim VERBATIM.
- **Corrección requerida:** sidecar con codificación reversible declarada, límites de lectura/decodificación y hash de los bytes originales; prueba de licencia con CRLF y sin newline final. No atribuir VERBATIM al envoltorio codificado. Reanudar sólo los pasos no ejecutados de integración.

- **Cierre focal V322:** owner0.4.78 reconstruido37/37,17profiles/7negatives,9lock negatives,65transport assertions;151preflight steps y161/1453/757/52 PASS. Adquisición real3sources y10outputs inmutables; ver GOVERNED_MAINTENANCE_ACQUISITION_V322.md. Las condiciones productivas permanecen abiertas.

### FAIL-20260908-527 — cache de source y oráculo negativo aceptan evidencia insuficiente

- **Estado:** `REGRESSION_PROVEN`
- **Observación V322:** acquire_upstream_sources0.4.77 devuelve PRESENT cuando existen carpeta+receipt con commit/archiveSHA iguales al lock, sin comprobar los archivos. Su test positivo crea una carpeta vacía y ese receipt sintético. Además Expect-Failure lanza dentro del try un error que contiene el texto esperado y luego lo acepta en el catch, de modo que una acción exitosa puede pasar el negativo.
- **Contención/corrección:** reproducir en staging; bloquear PRESENT legado que no puede probar sus bytes, sin borrar ni sobrescribir el source existente. El transporte opaco nuevo debe rehashear artefacto/licencia y bindings antes de PRESENT. Corregir el helper de tests para distinguir ausencia de excepción y comprobar su regresión con una acción exitosa. Mantener FAIL525 abierto hasta admisión y adquisición real del perfil.

- **Recurrencia V322 / FAIL473:** se supusieron el argumento --profile y exit1 del resolver; su interfaz real es --record, y RESEARCH_INCOMPLETE devuelve2 sin receipt. Se preservan ambas salidas, se corrigió el diagnóstico y no se interpretó como ready. HEAD de Sigstore no trae Content-Length; un probe Range0-0 devolvió206/Content-Range con tamaño10329 sin descargar el artefacto completo. No se reejecutó la preparación parcialmente escrita.

- **Cierre focal V322:** owner0.4.78 reconstruido37/37,17profiles/7negatives,9lock negatives,65transport assertions;151preflight steps y161/1453/757/52 PASS. Adquisición real3sources y10outputs inmutables; ver GOVERNED_MAINTENANCE_ACQUISITION_V322.md. Las condiciones productivas permanecen abiertas.

### FAIL-20260908-526 — fixture del preflight conserva el ID pnpm anterior

- **Estado:** `REGRESSION_PROVEN`
- **Observación V321:** al actualizar el floor/ID pnpm11.25.0, el test_toolchain_resolution rechazó Incorrect explicit routing: PnpmExecutable porque esperaba el ID11.19.0. El preflight se detuvo sin receipt final; no se informa PASS por leer la salida parcial.
- **Recuperación:** actualizar el fixture owner para el mismo ID exacto y ejecutar el test de routing y el preflight completo; conservar salida roja. No quitar el test ni relajar el routing explícito.

- **Cierre focal V321:** evidencia y límites en PNPM_BUNDLE_SECURITY_V321.md; owners canónicos actualizados, logs rojos preservados. FAIL525 de perfil permanece OPEN y no se oculta por estos PASS.

### FAIL-20260908-525 — registro de perfiles quedó en el slice histórico

- **Estado:** `REGRESSION_PROVEN`
- **Observación V321:** PROJECT_OFFICIAL_SOURCE_PROFILE_RECORD conserva el sliceV285 sin adquisición. Los expedientesV319–V321 sí contienen adquisiciones locales de tooling, hashes, firmas y SCA, pero no un perfil materializado/apply_source_profile con sus receipts. Esas pruebas no sustituyen el gate normativo de perfil ni corrigen retroactivamente su orden.
- **Contención:** no declarar ACQUIRED, ready_for_project_pack_plan o producción por los receipts manuales; preservar el registro before y añadir la discrepancia. No inventar respuestas/approval histórico ni relajar el contrato.
- **Resolución siguiente:** adaptar el owner de adquisición para un perfil de tooling oficial compatible, materializarlo, responder inputs reversibles desde instrucciones/probes y ejecutar validación/adquisición exacta con receipts nuevos. Mantener biblioteca y producto CONDITIONED/BLOCKED mientras falte ese gate; V321 sólo admite resultados focales de candidate.

- **Cierre focal V322:** owner0.4.78 reconstruido37/37,17profiles/7negatives,9lock negatives,65transport assertions;151preflight steps y161/1453/757/52 PASS. Adquisición real3sources y10outputs inmutables; ver GOVERNED_MAINTENANCE_ACQUISITION_V322.md. Las condiciones productivas permanecen abiertas.

### FAIL-20260908-524 — caché offline incompleto para políticas pnpm11.25

- **Estado:** `REGRESSION_PROVEN`
- **Observación V321:** web frozen/offline,115tests/1skip y build pasan. La instalación aislada Playwright enlaza3 paquetes pero rechaza ERR_PNPM_NO_OFFLINE_META para @playwright/test; faltaba metadata oficial del mirror para verificar políticas del lock. Conservar candidate-browser-install.log rojo, no tratar enlaces presentes como PASS.
- **Recuperación:** precargar metadata mediante instalación frozen online desde registry oficial, manteniendo versiones/integridades y políticas; después repetir frozen/offline y probar verifiers/browser. Documentar requisito del caché en owners canónicos. No desactivar trust/minimumReleaseAge ni mutar lock para eludir el error.

- **Recurrencia de diagnóstico V321 / FAIL473:** se supusieron nombres de packs/planes y ruta de pwsh; las lecturas fueron rechazadas antes de modificar producto. Inventario de archivos, búsqueda sólo de nombres y Get-Command resolvieron owners y ejecutable. Algunas consultas de líneas JSON extensas se truncaron y no se aceptaron como lectura íntegra. PyYAML no está disponible; la corroboración usa sólo claves literales del bloque packages, con recibo de línea, sin afirmar parseo YAML general. Mantener consultas acotadas antes de ejecutar.

- **Cierre focal V321:** evidencia y límites en PNPM_BUNDLE_SECURITY_V321.md; owners canónicos actualizados, logs rojos preservados. FAIL525 de perfil permanece OPEN y no se oculta por estos PASS.

### FAIL-20260908-523 — pnpm11.19.0 compilado incorpora brace-expansion vulnerable

- **Estado:** `REGRESSION_PROVEN`
- **Observación V321:** los891 archivos observados coinciden con el tarball oficial V319. Inventario del código compilado:5586 headers con identidad de store y48 manifests físicos,471 paquetes únicos,0 headers sin resolver. OSV2.5.1 all-packages encuentra brace-expansion5.0.8/GHSA-rgw5-rvv9-x895 (1paquete/1advisory).
- **Contención:** admisión pnpm11.19.0 reabierta; no tratar identidad891/891 ni package.json sin dependencias como SCA limpio de un bundle. ReachabilityUNKNOWN. Conservar baseline-osv.json/log, inventario y headers; consultar advisory/release oficiales antes de evaluar candidato aislado y corroborar identidades con source lock upstream.
- **Continuidad:** staging elite-v321-184b4d5926e44dc5a3cd1ecbb05ff077, T2803. No cambiar global ni manifests canónicos sin candidate gates/rollback; SCA no cubre aún componentes nativos transitivos.

- **Cierre focal V321:** evidencia y límites en PNPM_BUNDLE_SECURITY_V321.md; owners canónicos actualizados, logs rojos preservados. FAIL525 de perfil permanece OPEN y no se oculta por estos PASS.

### FAIL-20260908-522 — nuevo perfil no estaba en las composiciones del verificador raíz

- **Estado:** `REGRESSION_PROVEN`
- **Continuación del mismo fix:** el guard de checkpoint rechazó VERIFY_LIBRARY.ps1 porque no figuraba entre los owners permitidos del delta. No escribió nuevo estado. Se añadió sólo ese owner efectivamente editado, sin wildcard ni reejecutar la mutación de informe/ledger; se conservan todas las comprobaciones de SHA.
- **Observación V320:** checkpoint68 y materialización de todos los packs pasan, pero VERIFY_LIBRARY detecta52 planes y sólo51 compuestos. El nuevo perfil Node1/10 faltó en el registro explícito del verificador.
- **Corrección prevista:** integrar su composición y aserciones de1pack/10archivos en el owner VERIFY_LIBRARY, mantener la igualdad discovered/composed y repetir el gate. No rebajar el conteo ni excluir el perfil nuevo.

- **Corrección V320:** VERIFY_LIBRARY compone el perfil Node y exige1 pack del ID esperado y10 archivos; la comparación52discovered/composed se conserva. Before y log rojo preservados; gate completo se reejecuta sobre este owner.

### FAIL-20260908-521 — registro llamó firma a un binding SHA

- **Estado:** `REGRESSION_PROVEN`
- **Observación/revisión V320:** el registro de freshness decía «lock firmado por su binding» aunque el runner sólo compara SHA256. No había firma digital del lock.
- **Corrección:** conservar freshness-before-binding-wording.md y declarar «vinculado por SHA al runner»; la firma PGP comprobada pertenece al artefacto Node adquirido enV319, no al nuevo gate. Fuente/código intactos, no aceptación por integridad confundida con autenticidad.
- **Prevención/evidencia:** diferenciar digest, firma y trust root por artefacto; source-lock/runner y NODE_OFFICIAL_ADVISORY_GATE_V320.md. La primera aplicación de patch no coincidió y no mutó archivos; se releyó la línea exacta, recurrencia de diagnóstico FAIL473.

### FAIL-20260908-520 — semver del ZIP Node no coincide byte a byte con tarball npm

- **Estado:** `REGRESSION_PROVEN`
- **Observación V320:** la comprobación del tarball oficial semver7.8.5 falla primero en LICENSE respecto de la copia incluida por Node24.20.0. El tarball ya tiene integrity SHA512 válida; no se infiere identidad por versión ni se descarta la licencia discrepante.
- **Contención:** conservar artefacto/metadata y53 hashes before del lock candidato; inventariar todas las diferencias, contrastar fuente oficial y probar el grafo seleccionado antes de cambiar lock. No presentar la copia del ZIP Node como tarball npm VERBATIM.

- **Cierre V320:** NODE_OFFICIAL_ADVISORY_GATE_V320.md conserva before/after, rojos y arreglo canónico;30regresiones/4CLI/rebuild10/10 y SCA CLI2/0. No debilitar bindings, normalizar artefactos upstream para afirmar identidad ni repetir writers parciales.

### FAIL-20260908-519 — adapter candidato no vinculaba el bridge a su source-lock

- **Estado:** `REGRESSION_PROVEN`
- **Observación V320:** el runner verificaba motor oficial y semver pero copiaba bridge.cjs sin digest. Un bridge sustituido por JSON fijo produjo PASS sin ejecutar el motor; prove-bridge-red.py falla con la aserción esperada. Candidato aislado, sin pack ni promoción.
- **Corrección prevista:** incluir SHA del bridge en el lock, validarlo antes de copiar/ejecutar y conservar negativo de sustitución; completar límites de inventario y recibo de perfil antes de G4–G8.

- **Cierre V320:** NODE_OFFICIAL_ADVISORY_GATE_V320.md conserva before/after, rojos y arreglo canónico;30regresiones/4CLI/rebuild10/10 y SCA CLI2/0. No debilitar bindings, normalizar artefactos upstream para afirmar identidad ni repetir writers parciales.

### FAIL-20260908-518 — ledger reusable recibió incidentes locales abiertos

- **Estado:** `REGRESSION_PROVEN`
- **Observación V319:** checkpoint65 valida, pero VERIFY_LIBRARY rechaza4 filas OPEN del ledger canónico (2226/2228/2231/2232). Los incidentes locales de runtime/cobertura permanecen abiertos y no se pueden declarar resueltos para pasar el gate.
- **Contención:** conservar library-final.log; consultar contrato/validador del ledger, separar lección demostrada de cierre del incidente y corregir sólo la representación canónica admitida. No retirar bloqueos ni inventar cierre de npm/monitor/gate reusable.

- **Corrección:** ledger§1 y FAILURE_LEARNING_CONTRACT§5/§8 admiten REJECTED_COMPONENT y UPSTREAM_OPEN. Las4 filas se mueven a condiciones upstream con esos estados; IDs y causas intactos. Los incidentes511/513/516/517 del proyecto siguen OPEN, sin aceptación de riesgo. Verificador estructural sin modificar; repetir gate completo.

### FAIL-20260908-517 — probe oficial Node acepta corpus vacío y carece de deadline interno

- **Estado:** `OPEN`
- **Observación V319:** seis probes de auditoría con motor oficial fijado: baseline vulnerable exit1, candidato exit0, plataforma inválida/JSON malformado exit1; HTTP200 con corpus{} devuelve exit0 para Node24.14.1 vulnerable. Un servidor que no responde requiere kill externo a3s.
- **Contención:** no admitir el CLI upstream sin condiciones como gate fail-closed ni monitor. La consulta puntual V319 conserva hash exacto del corpus194 avisos/schedule y timeout externo30s; sus resultados no autorizan un reusable pack. El harness usa fixtures locales en copias independientes, sin modificar la evidencia original ni ejecutar un matcher propio.
- **Continuidad:** cerrar el expediente G0–G8 con estas condiciones demostradas; cualquier adapter reusable debe preservar motor oficial y probar integridad/freshness/límites y errores antes de promocionarse.

- **Sucesor V320:** gap del proyecto resuelto sólo mediante adapter offline separado,30tests/4CLI/10rebuild y G0–G8 de instancia. El alias OSV rechazado y el CLI raw upstream no se repararon ni quedaron admitidos; su condición permanece explícita.

### FAIL-20260908-516 — distribución Node24.20.0 incluye npm con advisories

- **Estado:** `OPEN`
- **Observación V319:** OSV2.5.1 sobre145 paquetes npm de la distribución firmada devuelve4 paquetes afectados y9 advisories: brace-expansion5.0.7, ip-address10.2.0, tar7.5.19 y undici6.27.0. El candidato completo no está limpio por tener cero avisos Node core.
- **Contención:** promoción bloqueada; preservar bundled-npm-osv.json/log y manifests exactos. Consultar primero advisories/releases oficiales por paquete, identificar consumers y reachability (UNKNOWN mientras no se demuestre), no silenciar findings ni editar la distribución firmada como si fuera upstream.

### FAIL-20260908-515 — conteo de cobertura OSV sin all-packages

- **Estado:** `REGRESSION_PROVEN`
- **Observación V319:** wrapper comparó145 componentes SBOM con4 resultados de OSV sin solicitar all-packages. Los4 eran los paquetes con advisories; la aserción falló y no se emitió recibo limpio.
- **Corrección prevista:** conservar salida inicial, repetir con la interfaz scan source y all-packages ya usada en V319, verificar identidad exacta del inventario y contar findings aparte. No confundir ausencia en la salida filtrada con ausencia de extracción.

- **Cierre V319:** evidencia before/after y prevención devueltas a autoridades; ver NODE_RUNTIME_SECURITY_V319.md. Firma VALIDSIG/byte identity, identity JSON con Node exacto y cobertura OSV all-packages probadas según el caso. No se transfieren estos PASS a seguridad integral.

### FAIL-20260908-514 — salida de pnpm confundida con identidad del proceso

- **Estado:** `REGRESSION_PROVEN`
- **Observación V319:** el wrapper rechazó la identidad antes del install: pnpm exec añade «Already up to date / Done» al stdout de process.execPath. El diagnóstico estructurado confirma Node24.20.0 en el destino aislado; no hubo selección del Node global.
- **Corrección prevista:** obtener identidad desde un archivo JSON escrito por el proceso hijo, verificar exit, versión, ruta y SHA; conservar stdout como log, sin quitar el guard ni filtrar una ruta por suposición.
- **Regresión:** repetir el probe con pnpm11.19.0 y comprobar la identidad antes de test/build; devolver prevención a la autoridad canónica.










- **Cierre V319:** evidencia before/after y prevención devueltas a autoridades; ver NODE_RUNTIME_SECURITY_V319.md. Firma VALIDSIG/byte identity, identity JSON con Node exacto y cobertura OSV all-packages probadas según el caso. No se transfieren estos PASS a seguridad integral.

### FAIL-20260908-513 — alias npm/node de OSV no detecta avisos conocidos del runtime

- **Estado:** `OPEN`
- **Observación V319:** OSV2.5.1 exacto encontró un paquete y0hallazgos para las proyecciones npm/node24.14.1 y24.20.0. El control24.14.1 tiene avisos oficiales posteriores conocidos; no se valida esta equivalencia entre paquete npm y runtime nativo.
- **Contención:** COVERAGE_BLOCKED para esa consulta; no declararla SCA limpio, no alterar JSON original ni agregar ignores. Consultar soporte oficial de identidad Git/runtime y mantener aparte la evaluación de rangos/advisories oficiales Node. No sustituir OSV por un scanner propio.
- **Evidencia:** V319 node-alias-probe.json, dos JSON/logs OSV y snapshot oficial security-wg195 archivos/194 avisos fijado por commit21bf6214b10d3a70250828c05b5e7026f4e27f50.

- **Sucesor V320:** gap del proyecto resuelto sólo mediante adapter offline separado,30tests/4CLI/10rebuild y G0–G8 de instancia. El alias OSV rechazado y el CLI raw upstream no se repararon ni quedaron admitidos; su condición permanece explícita.

### FAIL-20260908-512 — gpgv de Git rechaza ruta Windows del keyring

- **Estado:** `REGRESSION_PROVEN`
- **Observación V319:** gpgv2.4.9 de Git interpreta C:\\... como resource URL inválida y no carga el keyring oficial; NO_PUBKEY no demuestra falta de clave en el archivo. La descarga del runtime quedó detenida.
- **Contención:** conservar gpg-verification.log y el plaintext generado por el intento fallido como no autenticado. Retomar sólo verificación con cwd/rutas relativas compatibles, output nuevo y fingerprint oficial exacto; no repetir registros ni metadata, no saltar la firma.

- **Cierre V319:** evidencia before/after y prevención devueltas a autoridades; ver NODE_RUNTIME_SECURITY_V319.md. Firma VALIDSIG/byte identity, identity JSON con Node exacto y cobertura OSV all-packages probadas según el caso. No se transfieren estos PASS a seguridad integral.

### FAIL-20260908-511 — runtime Node24.14.1 anterior a correcciones oficiales de seguridad

- **Estado:** `OPEN`
- **Observación V319:** process.versions del Node instalado confirma24.14.1, OpenSSL3.5.5, llhttp9.3.0, nghttp2 1.68.0 y undici7.24.4. Las publicaciones oficiales de seguridad18-jun y29-jul2026 son posteriores; la de junio incluye fixes HIGH aplicables a la rama24 y release24.17.0. V317 cubría sólo Go, no Node.
- **Autoridad consultada primero:** https://nodejs.org/en/blog/vulnerability/june-2026-security-releases y https://nodejs.org/en/blog/vulnerability/july-2026-security-releases; release actual consultado https://nodejs.org/en/blog/release/v24.20.0.
- **Contención:** admisión Node reabierta, promotion bloqueada; no llamar limpio al runtime por el scan npm. Conservar versión/componentes/hash y evaluar candidato oficial exacto en destino aislado con gates y rollback antes de actualizar locks/packs. No instalar sobre el Node global ni ejecutar exploits destructivos.
- **Continuidad:** preparar expediente V319 y consultar el contrato de actualización/lock canónico; los PASS funcionales históricos siguen siendo de su snapshot, no certificación de seguridad de Node.

### FAIL-20260908-510 — espera de fin RSC bloquea Chromium y conteo parcial confunde progreso

- **Estado:** `REGRESSION_PROVEN`
- **Observación V317:** response.finished() del refresh RSC no termina en Chromium desktop/mobile; ambos operation-returns agotan60s. Firefox/WebKit completan23 fases. Run synchronized terminó80 fases/2 proyectos PASS, no92.
- **Corrección factual:** los mensajes de progreso asignaron conteos agregados a Chromium antes de recibir el resultado final. Retirada esa afirmación; contar por proyecto y detectar phase failed/Playwright failure durante el run.
- **Contención:** conservar log/traces, no aumentar timeout ni retirar negativos. Investigar ciclo de vida RSC y acotar sincronización al envío autenticado antes de cambiar cookies, o aislar el contexto del lector. Canonical0.1.23 todavía no está validado.

- **Cierre V317:** TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md;24 erratas preservadas,1.26.7 observado,OSV2/0 por versión y control2/41;fix canónico Playwright0.1.23,1/1 rebuild,92fases/4proyectos+agenda4/refund6 PASS. No cierre integral.

### FAIL-20260908-509 — host de reembolsos imprime error con configuración privada

- **Estado:** `REGRESSION_PROVEN`
- **Reproducción V318:** binario refund de V317/Go1.26.7 ejecutado con DSN sintético de puerto inválido, sin conectar. Exit1, pero stderr incluye PRIVATE-SYNTHETIC-DATABASE desde el error del parser propagado a log.Fatal(err).
- **Owner:** GO_OFFICIAL_RETURN_REFUND_WORKER0.1.2; main también imprime errores completos del Step. No se atribuye una fuga real de credenciales; el dato de prueba es sintético.
- **Contención:** salida del host debe conservar severidad/exit y códigos estáticos sin serializar error arbitrario; no cambiar retry, proveedores ni resultados durables. Preparar regresiones de proceso/startup y step antes del fix canónico.
- **Evidencia:** staging elite-v318-993950e79c854c4c9b26938962ffbdd3, startup-before.log/json. No red a provider ni datos reales.

- **Cierre V318:** REFUND_HOST_LOG_PRIVACY_V318.md; binario before divulgó marcador sintético, after conserva exit1 sin causa privada. Dos tests raíz ×3 por árbol,6PG,3/3 rebuild y Go gates PASS. Owner0.1.3; sin retry nuevo ni cierre de monitoring/producción.

### FAIL-20260908-508 — WebKit interrumpe navegación al cambiar a rol lector

- **Estado:** `REGRESSION_PROVEN`
- **Observación V317:** 86 fases completadas; WebKit falla en operation-returns al navegar /franchise después de setSession(READER), con Frame load interrupted. Chromium desktop/mobile y Firefox completaron sus23 fases. Logs TLS del fixture y trace conservados; no se atribuye todavía a Go ni a una vulnerabilidad.
- **Contención:** no declarar92/4 PASS. Analizar trace/config de aislamiento y postcondiciones antes de recuperar WebKit; no reintentar POST ni relajar TLS de un target. Agenda/refund todavía pendientes; resto de gates pasados se conserva.
- **Evidencia:** V317/baseline-browser-recovered.log y quote-browser-artifacts-3737324191 de la composición aislada.

- **Cierre V317:** TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md;24 erratas preservadas,1.26.7 observado,OSV2/0 por versión y control2/41;fix canónico Playwright0.1.23,1/1 rebuild,92fases/4proyectos+agenda4/refund6 PASS. No cierre integral.

### FAIL-20260908-507 — instalación del gate ignoró aislamiento ya documentado

- **Estado:** `REGRESSION_PROVEN`
- **Observación:** pnpm install del subdirectorio reportó exit0/already up to date pero resolvió el workspace padre; los cuatro browsers fallaron antes de operar por CLI de Playwright ausente. No es un fallo de Go1.26.7 ni del producto.
- **Fuente canónica existente:** MICROSOFT_PLAYWRIGHT_BROWSER_GATE ya exige --ignore-workspace en README y comandos; el wrapper V317 omitió el flag. No se modifica ni rebaja el pack correcto.
- **Corrección:** conservar log inicial, instalar exacto offline/frozen con --ignore-workspace, verificar CLI y retomar sólo browser/agenda/refund pendientes. PG33 y web115/build ya pasados no se repiten por este error de instalación.
- **Evidencia:** V317/baseline-browser.log y browser-install-before.log; recovery pendiente.

- **Cierre V317:** TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md;24 erratas preservadas,1.26.7 observado,OSV2/0 por versión y control2/41;fix canónico Playwright0.1.23,1/1 rebuild,92fases/4proyectos+agenda4/refund6 PASS. No cierre integral.

### FAIL-20260908-506 — versión de Go atribuida sin comprobar ejecutable

- **Estado:** `REGRESSION_PROVEN`
- **Observación V317:** current-toolchain del staging V281 contiene go.exe, VERSION y compile 1.26.8; reportes recientes V313/V314/V316 consignan1.26.7. V281 ya documentaba1.26.8 en ese directorio. No se atribuye un cambio externo reciente sin evidencia.
- **Contención:** no promover ni reutilizar claim de versión exacta. Preservar registros, identificar artefactos por hash, distinguir baseline1.26.7 de candidate1.26.8 y reconstruir/verificar el scope afectado con versión observada. No reescribir logs ni licencias ni suponer compatibilidad por el número de patch.
- **Scope:** trazabilidad y revalidación del toolchain/stdlib; no es un fallo demostrado de los tests ni un hallazgo de vulnerabilidad. Informe V317 pendiente; snapshot de los informes antes de corregir.

- **Cierre V317:** TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md;24 erratas preservadas,1.26.7 observado,OSV2/0 por versión y control2/41;fix canónico Playwright0.1.23,1/1 rebuild,92fases/4proyectos+agenda4/refund6 PASS. No cierre integral.

### FAIL-20260908-505 — observador de row lock conserva snapshot estadístico

- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** cuatro rojos reales de outbox; dos probes de lock no veían al waiter al consultar pg_stat_activity repetidamente dentro del mismo tx que retiene el lock. El diagnóstico no prueba todavía el check de tiempo post-lock.
- **Corrección:** observar pg_stat_activity mediante consultas autocommit del pool, snapshots frescos; mantener la transacción que bloquea y esperar expiración usando reloj PostgreSQL. Conservar outbox-red.log y repetir rojo focal, sin recrear DB ni ejecutar nuevamente red.py.
- **Evidencia:** V316/outbox-red-refined.log pendiente. No cambio de política ni relajación del oráculo.

- **Cierre V316:** seis rojos SQL/tres publisher preservados;33 tests raíz por árbol, tres receipts conectados por árbol y4/4 rebuild;Go test/vet/build/verify. Snapshot del observador corregido. OUTBOX_GENERATION_FENCING_V316.md; no evidencia de provider real ni operación integral.

### FAIL-20260908-504 — outbox sin generación ni vigencia de claim

- **Estado:** `REGRESSION_PROVEN`
- **Observación:** MarkPublished/Release sólo cotejan worker ID y estado no publicado, no lease vigente ni attempts. OutboxProcessor publica sin revalidar presupuesto. Un mismo worker reutilizado permite que un intento viejo altere otro nuevo; jobs ya tenían fencing en owner compartido.
- **Reproducción:** DB exclusiva elite_outbox_test_<unique>, fundación0001, sin filas del audit compartido. Casos de expiración, reclaim por mismo ID y espera de row lock que excede el lease, cada uno publish/release. Ledger existente de jobs se conserva como método, no se infiere prueba para outbox.
- **Scope:** fences de estado local y presupuesto del publisher; no exactly-once externo, supervisor/provider live, nuevo esquema ni regla comercial. Test usa OUTBOX_TEST_DATABASE_URL con guard, no reclama eventos de DB genérica.

- **Cierre V316:** seis rojos SQL/tres publisher preservados;33 tests raíz por árbol, tres receipts conectados por árbol y4/4 rebuild;Go test/vet/build/verify. Snapshot del observador corregido. OUTBOX_GENERATION_FENCING_V316.md; no evidencia de provider real ni operación integral.

### FAIL-20260908-503 — evidence mode omite SBOM/provenance sin ruta

- **Estado:** `REGRESSION_PROVEN`
- **Observación:** validate_operational_readiness sólo comprueba archivos si la referencia es nonempty_string; no exige supply_chain.sbom.artifact ni provenance.artifact. Una ruta omitida/nula/inválida puede saltar el modo require-evidence-files. Una ruta NUL puede abortar el resultado estructurado.
- **Reproducción:** expediente sintético válido y siete archivos de presencia, retirar/malformar ambas referencias en modo plan/files; comparar con control válido, archivos ausentes/directorios/escape y NUL. No se presenta un archivo vacío como evidencia semántica productiva.
- **Owner y scope:** SECURE_OPERATIONS_DELIVERY_CORE; presencia/resolución del gate operativo, no hashes/semántica/admisión productiva ni nuevos campos de negocio. Staging V315, baseline/ y readiness-red.log.

- **Resultado y corrección del diagnóstico:** NUL ya fue rechazado como file-not-found en Windows, sin aborto; primer oráculo se corrigió. Rojo refinado28 fallos de referencias y una excepción PermissionError inyectada.12tests/4CLI y2/2 rebuild PASS; evidencia OPERATIONAL_EVIDENCE_GATE_V315.md. No prueba acceso denegado real ni semántica de archivos.

### FAIL-20260908-502 — hosts salen antes del drenaje HTTP

- **Estado:** `REGRESSION_PROVEN`
- **Observación:** cmd/api y cmd/electromobility-api lanzan Shutdown en goroutine y retornan main al recibir ErrServerClosed de ListenAndServe. La documentación del Go1.26.7 fijado exige esperar a que Shutdown termine; salir cierra el proceso/pool con requests activas. El core además descarta errores de Shutdown.
- **Reproducción en curso:** extraer mecánicamente el mismo ciclo a owner HTTP común sin corregir aún su semántica, usar TCP loopback real y mantener un request activo al cancelar. Test exige detener admisión, esperar y conservar respuesta.
- **Scope de fix:** lifecycle HTTP/1 normal con deadline15s existente, no WebSockets/hijacked, señales OS reales ni rollback del target. No regla comercial ni dependencia nueva.
- **Evidencia:** staging V314 baseline/, red.py y shutdown-red.log; autoridad SOFTWARE_BACKEND_API_ENGINEERING11.2 y net/http.Server.Shutdown del toolchain1.26.7.

- **Cierre V314:** rojo demostrado; cinco tests raíz x3 en candidato y rebuild PASS,4/4 identidad, Go test/vet/build/verify. Evidencia HTTP_HOST_SHUTDOWN_V314.md; límites de señal OS/hijacked/handler no cooperativo intactos.

### FAIL-20260908-501 — matcher PowerShell confunde escenario failure con test fallido

- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** 92 receipts y cuatro proyectos PASS, pero el wrapper usa -match 'FAIL|SKIP' case-insensitive y rechaza phase=read-failure/list_failure=PASS. Abortó antes de agenda/refund; no fallo de navegador.
- **Corrección:** usar -cmatch multilínea anclado a tokens reales FAIL o --- FAIL:/SKIP:. Conservar log original y retomar sólo agenda/refund pendientes; no repetir las92 fases exitosas.
- **Evidencia:** browser-final.log, receipt-parser-regression.json y gates posteriores V313. Incluir probes de nombres negativos exitosos y fallos/skips auténticos.

- **Evidencia final V313:** COMPOSITION_DEPENDENCY_SECURITY_V313.md;OSV347 paquetes/0findings,5/5 canónico,Go/PG3/92 fases+agenda4/fuzz y2/2 binarios. Rojos/proyecciones previas preservados. FAIL498 cierra target arbitrario del test; fixture con bypass no acredita constraints de origen y sigue como límite T2802.

### FAIL-20260908-500 — grafo Go completo reabre el candidato mínimo x/text

- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** aunque fuente directa y wheels verificadas dan exit0, la proyección de go list -m all añade módulos seleccionados no compilados por los consumers. OSV detecta x/mod0.37.0 con GO-2026-6179/6180 en ambos grafos. El candidato x/text0.39.0 no cierra el grafo completo.
- **Alcance:** x/mod no aparece entre los9 módulos usados en paquetes/tests de cada consumer; sí figura en la selección transitiva de herramientas del upstream x/text. No se infiere explotación del producto ni se oculta el advisory.
- **Corrección en curso:** consultar advisories y cadena oficial para un candidato compatible con cierre transitivo, registrar before/after; no suprimir findings ni promover por scan directo. Gates de0.39.0 conservan valor para ese candidato, no se transfieren sin revisión al sucesor.
- **Evidencia:** reports/go-module-closure.json, final-combined-scan.json/exit1 y graphs de V313. La ejecución de navegadores sobre0.39.0 ya iniciada se conserva; nueva admisión requiere gates del delta efectivo.

- **Evidencia final V313:** COMPOSITION_DEPENDENCY_SECURITY_V313.md;OSV347 paquetes/0findings,5/5 canónico,Go/PG3/92 fases+agenda4/fuzz y2/2 binarios. Rojos/proyecciones previas preservados. FAIL498 cierra target arbitrario del test; fixture con bypass no acredita constraints de origen y sigue como límite T2802.

### FAIL-20260908-499 — identidad de build web incluye archivo generado por Next

- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** web-identity.py abortó en next-env.d.ts; V312 compilado agrega imports de routes/root-params generados y la composición fresca conserva el template. No hay prueba fallida de interfaz ni cambio de dependencia web.
- **Corrección:** conservar ambos SHA/contenidos y construir de nuevo el portal recompuesto V313 con locks frozen/offline; usar ese artifact para los gates conectados. No se elimina el guard ni se afirma identidad total no demostrada.
- **Evidencia:** web-identity.py original, web-generation-before.json y logs build/test V313. next-env es derivado del framework; no se edita para fingir identidad.

- **Evidencia final V313:** COMPOSITION_DEPENDENCY_SECURITY_V313.md;OSV347 paquetes/0findings,5/5 canónico,Go/PG3/92 fases+agenda4/fuzz y2/2 binarios. Rojos/proyecciones previas preservados. FAIL498 cierra target arbitrario del test; fixture con bypass no acredita constraints de origen y sigue como límite T2802.

### FAIL-20260908-498 — integración de reembolso permite cleanup en cualquier TEST_DATABASE_URL

- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado antes de ejecutar:** test histórico usa tenant fijo y session_replication_role=replica para seed/cleanup, pero aceptaba cualquier URL. Podía borrar filas del mismo tenant en un target ajeno y no constituye prueba de integridad del journey de origen.
- **Corrección:** guard canónico de loopback literal y base descartable elite_refund_test_<unique>, timeout30s y ocho targets negativos, incluido fallback remoto. El gate del parche usa base nueva separada con53 migraciones. Nunca ejecutar este fixture sobre elite_confirmation_v296 ni bases de proyecto.
- **Límite conservado:** seed/cleanup histórico aún omite triggers en la base descartable; su PASS prueba compatibilidad de reembolso, no constraints o aceptación end-to-end del negocio. Mejorar fixture de origen sigue T2802, no se oculta como integración completa.
- **Evidencia:** staging V313; source anterior reconstruible desde checkpoint52, test corregido y receipts de base aislada. No limpieza ejecutada contra target compartido.

- **Evidencia final V313:** COMPOSITION_DEPENDENCY_SECURITY_V313.md;OSV347 paquetes/0findings,5/5 canónico,Go/PG3/92 fases+agenda4/fuzz y2/2 binarios. Rojos/proyecciones previas preservados. FAIL498 cierra target arbitrario del test; fixture con bypass no acredita constraints de origen y sigue como límite T2802.

### FAIL-20260908-497 — probe de regresión exige timeout aunque su guard detecta no progreso

- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** el probe contra x/text0.29.0 falló inmediatamente con iteration failed to advance; su guard finito detectó el defecto antes del timeout. El wrapper esperaba sólo timeout y abortó antes de probar0.39.0.
- **Corrección:** conservar rojo e input oficial CL794100; aceptar como oráculo negativo el fallo explícito de progreso, ejecutar sólo el candidato pendiente. No quitar el guard ni forzar consumo indefinido para obtener un timeout.
- **Evidencia:** staging V313 norm-before.log, norm-probe.py original y norm-resume.py sucesor. Código del producto intacto por este fallo del wrapper.

- **Evidencia final V313:** COMPOSITION_DEPENDENCY_SECURITY_V313.md;OSV347 paquetes/0findings,5/5 canónico,Go/PG3/92 fases+agenda4/fuzz y2/2 binarios. Rojos/proyecciones previas preservados. FAIL498 cierra target arbitrario del test; fixture con bypass no acredita constraints de origen y sigue como límite T2802.

### FAIL-20260908-496 — comprobación de licencia usa nombre corporativo supuesto

- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** authority.py recuperó patch oficial y delta, luego exigió texto Google Inc. en LICENSE. La fuente íntegra dice Google LLC; no es incompatibilidad de licencia ni permiso faltante.
- **Corrección:** leer cláusulas completas BSD-3-Clause, conservar LICENSE y su SHA; comparar entre módulos y registrar el titular exacto. No cambiar la licencia ni ocultar assertion inicial.
- **Evidencia:** authority.py inicial, LICENSE original del módulo autenticado y receipt de revisión sucesor; no publicación ni copiado de patch local.

- **Evidencia final V313:** COMPOSITION_DEPENDENCY_SECURITY_V313.md;OSV347 paquetes/0findings,5/5 canónico,Go/PG3/92 fases+agenda4/fuzz y2/2 binarios. Rojos/proyecciones previas preservados. FAIL498 cierra target arbitrario del test; fixture con bypass no acredita constraints de origen y sigue como límite T2802.

### FAIL-20260908-495 — locks Python por URL producen cero paquetes en OSV

- **Estado:** `REGRESSION_PROVEN`
- **Fallo reproducido:** scan recursivo omite seis requirements*.lock; forzar parser oficial requirements.txt visita seis archivos pero devuelve128/0paquetes. Los locks son wheel URL+SHA256, no pins name==version reconocidos. No equivale a dependencias seguras ni cobertura completa.
- **Corrección en curso:** derivar inventario sólo de wheels fijadas verificando hash y METADATA, conservar correspondencia por consumer/URL/lockSHA y consultar OSV con artefacto soportado. No modificar locks originales ni inventar versión desde nombre ambiguo.
- **Evidencia:** staging V313 reports/python-scan.json, python-scan.log y python-scan-exit.json. Coverage bloqueada hasta packages verificables por los seis consumers.

- **Evidencia final V313:** COMPOSITION_DEPENDENCY_SECURITY_V313.md;OSV347 paquetes/0findings,5/5 canónico,Go/PG3/92 fases+agenda4/fuzz y2/2 binarios. Rojos/proyecciones previas preservados. FAIL498 cierra target arbitrario del test; fixture con bypass no acredita constraints de origen y sigue como límite T2802.

### FAIL-20260908-494 — OSV vigente de composición devuelve hallazgos

- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** OSV oficial2.5.1, SHA/commit verificados, scan puntual V313 sobre67 packs/745 archivos frescos termina exit1. No se interpreta como gate limpio ni readiness de producción.
- **Scope:** dependencies fijas, sin resolución de manifests ni call analysis. Consultas OSV/deps.dev; sin instalaciones, cuentas, payloads ni scheduled task. JSON all-packages y exit preservados.
- **Corrección en curso:** triage por paquete/consumer/advisory oficial, cobertura de extractores y patches candidatos con rollback. No ejecutar osv fix ni actualizar por latest. Python .lock no apareció en el primer inventario; cerrar cobertura explícita antes de declarar scan del conjunto.
- **Evidencia:** staging V313 reports/source-scan.json, source-scan.log, source-scan-exit.json, intake.json y triage.json. Promoción sigue bloqueada.

- **Evidencia final V313:** COMPOSITION_DEPENDENCY_SECURITY_V313.md;OSV347 paquetes/0findings,5/5 canónico,Go/PG3/92 fases+agenda4/fuzz y2/2 binarios. Rojos/proyecciones previas preservados. FAIL498 cierra target arbitrario del test; fixture con bypass no acredita constraints de origen y sigue como límite T2802.

### FAIL-20260908-493 — generador de informe cuenta subtests como tests raíz

- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** report.py abortó antes de escribir el informe: contaba81 PASS (incluye subtests) frente a33 esperados. El log contiene33 tests raíz,81 PASS totales, cero FAIL/SKIP; el gate PostgreSQL fue exitoso.
- **Corrección:** anclar conteo a ^--- PASS: con modo multilínea; no se altera el log ni el número de11 tests por3 repeticiones. Reinicio ya PASS, no se repite por un fallo del reporte.
- **Evidencia:** report-guard-red.log, postgres-rebuilt.log e informe V312 posterior; scripts antes/después retenidos.

- **Evidencia final V312:** RETURN_OPERATIONS_RECOVERY_V312.md;92 fases/4,agenda4,11 PG3,fuzz10s,rebuild10/10 y restart idéntico. Rojo conservado, fuente canónica reparada; sin debilitar gates ni reglas comerciales.

### FAIL-20260908-492 — wrapper de agenda confunde exit0 con cuatro skips

- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** la agenda se invocó sin ELITE_CONFIRMATION_E2E=1; Go terminó exit0 con cuatro SKIP. El wrapper imprimió BROWSER_AGENDA_PASS; no es evidencia de integración. Las92 fases del gate principal sí pasaron.
- **Corrección:** conservar agenda-skipped.log, activar explícitamente el gate descartable y exigir cuatro APPOINTMENT_AGENDA_BROWSER_POSTGRES_PASS y cero SKIP/FAIL además del exit. El generador del informe ya rechazaba ausencia de receipts.
- **Prevención:** siempre validar count de receipts esperados y skips, no solamente el exit del runner.

- **Evidencia final V312:** RETURN_OPERATIONS_RECOVERY_V312.md;92 fases/4,agenda4,11 PG3,fuzz10s,rebuild10/10 y restart idéntico. Rojo conservado, fuente canónica reparada; sin debilitar gates ni reglas comerciales.

### FAIL-20260908-491 — identidad de reconstrucción detecta CRLF del editor Python

- **Estado:** `REGRESSION_PROVEN`
- **Recurrencia V320:** primer compose del nuevo gate aborta en profile.example.json: el writer Python generó CRLF, el compositor exige LF. Mantener rebuild-compose.log y bytes before; normalizar sólo archivos del candidato, recalcular sus bindings y sincronizar el pack mediante el updater canónico. No repetir el writer que ya creó pack/plan ni debilitar el hash.
- **Fallo observado:** canonical.log aborta comparación del spec Playwright después de sincronizar tres packs y recomponer. El sincronizador normaliza LF; el generador Python escribió CRLF en Windows. No se declara identidad ni se reejecuta la mutación parcial.
- **Corrección en curso:** conservar bytes anteriores, normalizar sólo terminadores de los diez archivos cambiados y comparar de nuevo con rebuilt; ejecutar gates finales sobre rebuilt. La comparación hash sigue obligatoria.
- **Evidencia:** staging V312 canonical.log, normalization.json e identity-rebuilt.log; no promoción ni salida de producción.

- **Evidencia final V312:** RETURN_OPERATIONS_RECOVERY_V312.md;92 fases/4,agenda4,11 PG3,fuzz10s,rebuild10/10 y restart idéntico. Rojo conservado, fuente canónica reparada; sin debilitar gates ni reglas comerciales.

- **Cierre V320:** NODE_OFFICIAL_ADVISORY_GATE_V320.md conserva before/after, rojos y arreglo canónico;30regresiones/4CLI/rebuild10/10 y SCA CLI2/0. No debilitar bindings, normalizar artefactos upstream para afirmar identidad ni repetir writers parciales.

### FAIL-20260908-490 — limpieza del sender retira un tipo compartido

- **Estado:** `REGRESSION_PROVEN`
- **Corrección de fecha:** el encabezado inicial decía FAIL-20260909-490; reloj observado 2026-09-08 03:00:37 UTC. Se conserva aquí la referencia anterior; no se cambia el fallo ni su evidencia.
- **Fallo:** V312 web-build-r1.log: TypeScript TS2304 en franchise-command-panel.tsx, símbolo Command aún usado por otro control. Las114 pruebas unitarias pasan pero el build falla; no se declara compilación correcta.
- **Causa:** se retiró el alias junto al sender genérico sin inventariar todas sus referencias.
- **Acción:** restaurar sólo el tipo compartido, conservar la eliminación del sender engañoso y verificar el build. Antes de borrar símbolos, buscar sus usos reales.

- **Evidencia final V312:** RETURN_OPERATIONS_RECOVERY_V312.md;92 fases/4,agenda4,11 PG3,fuzz10s,rebuild10/10 y restart idéntico. Rojo conservado, fuente canónica reparada; sin debilitar gates ni reglas comerciales.

### FAIL-20260908-489 — devoluciones omiten binding vigente de stock/pedido/entrega

- **Estado:** `REGRESSION_PROVEN`
- **Fallo reproducido:** V312 returns-scope-red.log, fixture sintético con stock movido a otra organización: ReturnCases devuelve2 casos; ReceiveReturn y DecideReturn aceptan la relación incoherente y generan1 recibo/1 decisión/4 solicitudes/2 eventos. No hay efectos externos ni evidencia de incidente productivo.
- **Causa observada:** lecturas/comandos se apoyan en autorización/recibo sin validar y bloquear el grafo completo de organización, entrega, pedido, cliente y stock; listado además carga efectos sin acotarlos a sus casos seleccionados.
- **Acción:** validar bindings completos y locks coherentes en los owners existentes, rechazar antes de efecto, acotar lectura de efectos y conservar snapshot consistente. Mantener rojos, probar carrera/rollback/scope y devolver fix al pack canónico antes de cerrar.

- **Evidencia final V312:** RETURN_OPERATIONS_RECOVERY_V312.md;92 fases/4,agenda4,11 PG3,fuzz10s,rebuild10/10 y restart idéntico. Rojo conservado, fuente canónica reparada; sin debilitar gates ni reglas comerciales.

### FAIL-20260908-488 — recepción y decisión de devoluciones sin recuperación exacta

- **Estado:** `REGRESSION_PROVEN`
- **Fallo reproducido:** V312 returns-red.log: recepción200/replay409 y decisión200/replay409 dejan un recibo, una decisión, cuatro solicitudes y dos eventos auditados; ambas consultas de resultado devuelven404.
- **Impacto:** los formularios por caso conservan sender genérico, recarga y mensaje no aplicado tras respuesta incierta; lista limitada100 no es oracle de ausencia.
- **Acción:** consulta exacta por autorización con recibo/decisión/solicitudes y scope/binding comprobados; referencia global de operación recuperable aunque falle la lista, sin reenviar ni atribuir ejecución downstream a solicitudes.

- **Evidencia final V312:** RETURN_OPERATIONS_RECOVERY_V312.md;92 fases/4,agenda4,11 PG3,fuzz10s,rebuild10/10 y restart idéntico. Rojo conservado, fuente canónica reparada; sin debilitar gates ni reglas comerciales.

### FAIL-20260908-487 — edición inline rechazada por parser PowerShell

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** se intentó pasar un reemplazo Python/SQL con comillas escapadas de forma incompatible con PowerShell; ParserError antes de ejecutar Python. Sin modificación del candidato por ese comando.
- **Corrección:** reemplazo literal mediante apply_patch. Usar scripts separados o patch para código multilínea; conservar valores SQL literales sin interpolación de shell. Regresión pendiente del gate V311.

- **Evidencia final V311:** CHECKLIST_COMPLETION_RECOVERY_V311.md;88 fases/4,agenda4,9 tests PG3,rebuild10/10 y restart idéntico. Rojos conservados; sin alterar unicidad/inmutabilidad. FAIL486 recurre a438: fixture con VIN/batería sintéticos, no compatibilidad universal demostrada.

### FAIL-20260908-486 — gate PostgreSQL V311 falla durante verificación

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** postgres-r1.log conserva el primer gate rojo de persistencia V311. No se atribuye PASS a los casos restantes ni se debilitan restricciones.
- **Acción:** localizar la aserción/fixture o defecto por evidencia, corregir en su owner y repetir el gate con los negativos intactos.

- **Evidencia final V311:** CHECKLIST_COMPLETION_RECOVERY_V311.md;88 fases/4,agenda4,9 tests PG3,rebuild10/10 y restart idéntico. Rojos conservados; sin alterar unicidad/inmutabilidad. FAIL486 recurre a438: fixture con VIN/batería sintéticos, no compatibilidad universal demostrada.

### FAIL-20260908-485 — completado de checklist sin recuperación exacta

- **Estado:** `REGRESSION_PROVEN`
- **Fallo reproducido:** V311 completion-red.log: POST200 seguido de replay409, dos respuestas y un evento con actor persisten; GET de recuperación404. La UI mantiene sender genérico y puede afirmar no aplicado después del commit.
- **Impacto:** estado incierto y riesgo de reintento manual sobre entrega ya presentada; familia FAIL457.
- **Acción:** recuperar resultado durable por identidad natural, scope y binding completos, conservar respuestas inmutables/actor y estado vigente; UI consulta sin reenviar ni reemplazar reglas de entrega.

- **Evidencia final V311:** CHECKLIST_COMPLETION_RECOVERY_V311.md;88 fases/4,agenda4,9 tests PG3,rebuild10/10 y restart idéntico. Rojos conservados; sin alterar unicidad/inmutabilidad. FAIL486 recurre a438: fixture con VIN/batería sintéticos, no compatibilidad universal demostrada.

### FAIL-20260908-484 — wrapper de staging V311 detiene la preparación

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** el wrapper informó composition failed antes de instalar dependencias. Log de composición y destino conservados para determinar si falló el compositor o la comprobación de exit code.
- **Acción:** inspeccionar resultado real sin recomponer sobre destino existente ni atribuir PASS al wrapper; corregir la causa antes de continuar.

- **Causa y prueba:** LASTEXITCODE no expresa éxito de un script PowerShell: conserva valor nulo o previo. Dos probes con script exitoso reproducen falso fallo. Usar excepciones/estado PowerShell y reservar LASTEXITCODE para ejecutables nativos. Composición67/745, identidad10/10 e instalaciones offline/frozen PASS en staging V311; logs conservados.

### FAIL-20260908-483 — reconstrucción V310 difiere del árbol probado

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** canonical.ps1 detectó SHA distinto en enterprise-web.spec.mjs después de sincronizar tres packs y cuatro planes. Gate detenido; no se declara reconstrucción idéntica.
- **Acción:** conservar ambos árboles, comparar bytes y corregir fuente canónica o normalización antes de verificar y checkpoint. No repetir el script de mutación completo.

- **Evidencia final V310:** CHECKLIST_PUBLICATION_RECOVERY_V310.md;84 fases/4,agenda4,8 tests PG3,rebuild10/10 y restart idéntico. Rojos preservados; sin alterar unicidad/inmutabilidad ni ocultar fallos.

### FAIL-20260908-482 — nuevo handler de versión omite import strconv

- **Estado:** `REGRESSION_PROVEN`.
- **Fallo:** compilación V310 rechaza undefined strconv en el GET de checklist; no se ejecutó la prueba HTTP en ese intento.
- **Corrección:** declarar la importación estándar utilizada, repetir el gate y devolver el archivo completo al pack; checklist-r1.log conserva el rojo.

- **Evidencia final V310:** CHECKLIST_PUBLICATION_RECOVERY_V310.md;84 fases/4,agenda4,8 tests PG3,rebuild10/10 y restart idéntico. Rojos preservados; sin alterar unicidad/inmutabilidad ni ocultar fallos.

### FAIL-20260908-481 — publicación de checklist omite actor del evento

- **Estado:** `REGRESSION_PROVEN`.
- **Fallo:** POST HTTP publica la versión y conserva su creador en tabla, pero outbox no incluye actor_subject; TestChecklistPublicationHTTPRecovery observa cadena vacía. GET de recuperación inexistente devuelve404 (familia457).
- **Corrección prevista:** añadir actor autenticado al mismo evento atómico y consulta exacta de versión publicada; no cambiar identidad natural, inmutabilidad ni reglas de la checklist. checklist-red.log preservado.

- **Evidencia final V310:** CHECKLIST_PUBLICATION_RECOVERY_V310.md;84 fases/4,agenda4,8 tests PG3,rebuild10/10 y restart idéntico. Rojos preservados; sin alterar unicidad/inmutabilidad ni ocultar fallos.

### FAIL-20260908-480 — replay sintético cambia bytes JSON del BFF

- **Estado:** `REGRESSION_PROVEN`.
- **Fallo:** browser-r2 recupera cupo pero replay directo devuelve409: el fixture serializa capacity antes de las fechas; el BFF lo serializa al final. decodeHashedJSON fija bytes del request.
- **Corrección:** replay exacto con orden del BFF; conservar el caso reordenado como negativo409 y no cambiar hashing ni semántica de idempotencia para acomodar la prueba.

- **Evidencia final V309:** SLOT_CREATION_RECOVERY_V309.md;80 fases/4,agenda4,7 tests PG3,rebuild9/9,plan PASS. Logs rojos y sus causas conservados; no se retiran guards ni negativos.

### FAIL-20260908-479 — fixture de datetime-local usa UTC como hora local

- **Estado:** `REGRESSION_PROVEN`.
- **Fallo:** Chromium V309 envió04:34–05:34Z desde un input rellenado con ISO UTC; la jornada sintética terminaba03:34Z. El guard rechazó409 correctamente.
- **Evidencia:** trace operation-create-slot / recurso694c338b…json y lectura SQL de slot-create-working-30; browser-r1.log conservado.
- **Corrección:** generar texto datetime-local usando el timezone del propio navegador; no ensanchar jornada ni retirar el guard. Reejecutar matriz y devolver harness canónico.

- **Evidencia final V309:** SLOT_CREATION_RECOVERY_V309.md;80 fases/4,agenda4,7 tests PG3,rebuild9/9,plan PASS. Logs rojos y sus causas conservados; no se retiran guards ni negativos.

### FAIL-20260908-478 — fixtures de capacidad conservan POST sin clave

- **Estado:** `REGRESSION_PROVEN`.
- **Fallo:** web-r1 espera201 para alta sin Idempotency-Key ahora requerida;1 FAIL/111 PASS/1 SKIP. El contrato nuevo rechaza400 correctamente.
- **Corrección:** adaptar sólo los positivos heredados con clave y cuarto argumento; conservar negativos sin clave/overposting. Actualizar fake y Subject del unit HTTP, ejecutar suites y reconstruir.

- **Evidencia final V309:** SLOT_CREATION_RECOVERY_V309.md;80 fases/4,agenda4,7 tests PG3,rebuild9/9,plan PASS. Logs rojos y sus causas conservados; no se retiran guards ni negativos.

### FAIL-20260908-477 — fixture de turno omite jornada requerida

- **Estado:** `REGRESSION_PROVEN`.
- **Fallo:** rojo V309 devolvió409/409 y0 cupos porque faltaba availability working requerida por trigger0015; no prueba todavía replay.
- **Corrección:** materializar jornada sintética compatible mediante owner antes del alta; preservar guard de agenda y repetir la reproducción. Log inicial slot-red.log conservado.

- **Evidencia final V309:** SLOT_CREATION_RECOVERY_V309.md;80 fases/4,agenda4,7 tests PG3,rebuild9/9,plan PASS. Logs rojos y sus causas conservados; no se retiran guards ni negativos.

### FAIL-20260908-476 — captura UTF-8 no coincide con stdout Windows del validador

- **Estado:** `REGRESSION_PROVEN`.
- **Fallo:** el contrato ya quedó conciliado con TEST27, pero la captura de nivel evidence falló al decodificar byte0xe1 de stdout; no se toma ese error como PASS ni se vuelve a insertar las pruebas.
- **Corrección:** forzar UTF-8 en el proceso Python hijo y repetir sólo validación; snapshot previo del contrato y logs conservados.

- **Evidencia final V308:** RESOURCE_CREATION_RECOVERY_V308.md;76 fases/4,agenda4,rebuild9/9. Contrato plan PASS y evidence BLOCKED esperado; no elimina pruebas/bloqueos pendientes.

### FAIL-20260908-475 — trazabilidad de assurance detenida en V305

- **Estado:** `REGRESSION_PROVEN`.
- **Fallo:** al conciliar el contrato de ejecución, TEST-24/EVID-15 eran los últimos registros; V306/V307 tenían expedientes y checkpoints pero no entradas propias en implementation_assurance.
- **Corrección:** incorporar pruebas/evidencia exactas V306–V308 y sus límites por dimensión; conservar TESTs bloqueados y no convertir el gate de plan en assurance integral. Validar contrato antes del siguiente checkpoint.

- **Evidencia final V308:** RESOURCE_CREATION_RECOVERY_V308.md;76 fases/4,agenda4,rebuild9/9. Contrato plan PASS y evidence BLOCKED esperado; no elimina pruebas/bloqueos pendientes.

### FAIL-20260908-474 — fixture de recurso no respeta normalización del BFF

- **Estado:** `REGRESSION_PROVEN`.
- **Fallo:** nuevo test esperaba principal_subject undefined para bahía; contrato existente envía cadena vacía. web-test-r1.log:1 FAIL/110 PASS/1 SKIP.
- **Corrección:** exigir exactamente cadena vacía; no cambiar validación ni transporte productivo para acomodar el mock. Reejecutar web/build y gate conectado, devolver fixture canónico.

- **Evidencia final V308:** RESOURCE_CREATION_RECOVERY_V308.md;76 fases/4,agenda4,rebuild9/9. Contrato plan PASS y evidence BLOCKED esperado; no elimina pruebas/bloqueos pendientes.

### FAIL-20260908-473 — diagnóstico repite ruta inferida y exceso de salida

- **Reincidencia V319:** el digest Node copiado al script perdió un carácter y la aserción paró antes de descargar metadata. Se conserva el script fallido y se retoma después del snapshot, con identidad observada y contraste contra SHASUMS oficial; no se repite la creación ni se atribuye mutación del ejecutable sin evidencia.

- **Reincidencia V317 al preparar continuidad:** se supuso workers/processor.go inexistente; se recuperó inventario literal antes de leer. No hubo mutación ni se dedujo ausencia de la capability por el error de path.

- **Reincidencia V317:** ReadAllText sobre log redirigido activo falló por sharing; lectura posterior usa FileStream/FileShare.ReadWrite y no convierte excepción en cero fallos. Una búsqueda demasiado amplia volvió a truncar salida; se recuperaron campos materiales con consultas acotadas antes de editar.


- **Reincidencia diagnóstica V317:** scan sbom no era un subcomando del OSV fijado; la ayuda confirmó scan source con extractor SBOM. verify_pack de fuzz exigía GoExecutable; se leyó su firma y se pasó el ejecutable1.26.7 exacto. Ninguno ejecutó producto ni habilitó un PASS; recibos iniciales conservados.


- **Reincidencia operativa V316:** un apply_patch con contexto abreviado no coincidió y no modificó archivos. Se usó el anchor exacto del loop del checkpoint y se validó resume; no se repitieron scripts de mutación ni gates de producto ya pasados.


- **Reincidencia V315:** se supuso un owner GO_FRANCHISE_RUNTIME_HOST ausente y un glob posicional JOB* falló. Inventario observado ubicó GO_CONNECTED_CONVERSATION_RUNTIME y los owners GO_ENTERPRISE_BACKEND_CORE/GO_RELIABLE_ASYNC_WORKERS. Ninguna ausencia de implementación se infiere del comando inválido.


- **Reincidencia V314:** se supuso canonical-final.ps1 inexistente y búsquedas amplias excedieron salida. Se inventariaron los scripts y se leyó canonical-successor.ps1; diagnóstico posterior acotado por símbolo. No se trató la lectura fallida como evidencia ni se cambió código por ella. Se corrigió además un glob posicional inválido de rg usando --files -g; no mutación ni falso claim de ausencia.


- **Estado:** `RECURRENCE_PROVEN`.
- **Clase:** operación del agente. Se intentó PORTABLE_RUNTIME_CONFIG_V295.md inexistente y lecturas excedieron salida; no implica defecto de producto ni cambio de gate.
- **Corrección:** rg --files localizó CREDENTIAL_PREPARATION_V295.md; usar rangos menores y resúmenes estructurados. Consultada LIB-FAIL-2183. No repetir la lectura completa ni atribuir contenido truncado como revisado.

- **Reincidencia V308:** sufijo de staging V289 y nombre de migración0008 supuestos no existían; se usa inventario literal y0007_appointment_capacity.up.sql. No cambios por esos errores; no se toman búsquedas fallidas como ausencia de capability.

- **Reincidencia V309:** se buscó un directorio contracts inexistente; inventario literal mostró docs/config/ops y owners Go. No atribuir OpenAPI ausente globalmente a esa búsqueda parcial.

- **Reincidencia V311:** lectura especulativa de ruta customer/actions ausente y salidas extensas truncadas. Diagnóstico resuelto desde el handler inventariado; usar inventario y recortes por símbolo, no rutas supuestas ni silenciamiento de errores.

- **Reincidencia V312, diagnóstico:** búsqueda amplia de `setMessage` sumada a lecturas extensas produjo salida truncada. Sin mutación de producto. Se corrigió leyendo sólo la sección de devoluciones y los anchors del harness. No inferir de salidas incompletas.

### FAIL-20260908-472 — selector de tipo de intervalo incluye las opciones

- **Estado:** `REGRESSION_PROVEN`.
- **Clase:** fixture/contrato accesible. Chromium V307 r1 agota 60s en getByLabel Tipo exacto; el snapshot prueba combobox Tipo visible y habilitado.
- **Corrección:** seleccionar el rol combobox y nombre accesible exacto; conservar la operación y todas sus aserciones. Reejecutar y devolver el harness al pack canónico.
- **Evidencia:** staging V307, browser-r1.log y trace operation-create-availability.

- **Evidencia final:** AVAILABILITY_CREATION_RECOVERY_V307.md;72 fases/4 proyectos y harness canónico reconstruido idéntico.

### FAIL-20260907-471 — inventario preflight no incluye evidencia V306

- **Estado:** `REGRESSION_PROVEN`
- **Observado:** VERIFY_LIBRARY llegó al gate de inventario y rechazó el contador anterior de Markdown; esperaba160/1436/739/51 tras añadir el expediente V306.
- **Corrección:** actualizar únicamente el conteo del owner FRANCHISE_PREFLIGHT_GAP al inventario observado; conservar el log rojo y repetir el verificador sin debilitar su control.

- **Evidencia:** checkpoint38/39 y resume PASS; VERIFY_LIBRARY_PASS posterior sobre inventario160/1436/739/51. Detalle V306, logs rojos y cadena previa conservados.

### FAIL-20260907-470 — next_action excede el cursor compacto

- **Estado:** `REGRESSION_PROVEN`
- **Observado V306:** checkpoint38 rechazado antes de append porque next_action excede800 caracteres; event log conserva37 entradas intactas.
- **Corrección:** compactar siguiente acción y enlazar detalle V306; conservar revisión38 pendiente sin saltar secuencia ni alterar eventos. Revalidar antes de guardar.
- **Prevención:** comprobar longitud local antes del writer; no duplicar roadmap en el cursor.

- **Evidencia:** checkpoint38/39 y resume PASS; VERIFY_LIBRARY_PASS posterior sobre inventario160/1436/739/51. Detalle V306, logs rojos y cadena previa conservados.

### FAIL-20260907-469 — emisión de cotización pierde actor autenticado

- **Estado:** `REGRESSION_PROVEN`
- **Clasificación:** CODE / auditoría T2803, mismo recorrido V306.
- **Reproducción:** quote-actor-red.log consulta únicamente fixtures quote-create-*; seis eventos quotation.issued persistidos y cero con actor_subject=writer.
- **Causa:** CreateQuote recibe tenant y clave pero no Subject HTTP; el outbox atómico conserva organización/lead, sin autor de emisión.
- **Corrección prevista:** método auditado fail-closed para HTTP sobre el owner existente, actor tomado del token, atomicidad con cotización/idempotencia/outbox y replay sin sustituir actor original. Métodos internos legacy no heredan el claim.

- **Evidencia final V306:** QUOTE_CREATION_RECOVERY_V306.md;11 archivos recompuestos, Go/web, PostgreSQL3,68 fases/4 proyectos y agenda4. Logs fallidos conservados; no promoción integral.

### FAIL-20260907-468 — edición y tipos del harness V306

- **Estado:** `REGRESSION_PROVEN`
- **Clasificación:** CODE / tooling de prueba, sin promoción.
- **Observado:** Python read_text sin UTF-8 falla al leer TS; ancla no única detiene script parcial; continuación elimina declaración repo del fixture de disponibilidad. Go focused rechaza undefined repo. TypeScript rechaza acceso calls[0][1] de mock inferido sin argumentos aunque Vitest pasa.
- **Corrección:** UTF-8 explícito, continuar sólo desde etapa observada, ancla contextual, restaurar repo y tipar RequestInfo/RequestInit. Conservar logs go-focused/web-build-first; repetir gates afectados.

- **Evidencia final V306:** QUOTE_CREATION_RECOVERY_V306.md;11 archivos recompuestos, Go/web, PostgreSQL3,68 fases/4 proyectos y agenda4. Logs fallidos conservados; no promoción integral.

### FAIL-20260907-466 — cotización del portal omite clave idempotente

- **Estado:** `REGRESSION_PROVEN`
- **Clasificación:** CONTRACT / T2804 / relacionada FAIL457.
- **Observado V306:** formulario y BFF create-quote invocan protectedPost sin Idempotency-Key; CreateQuote exige longitud mínima16. El test BFF mockeado esperaba201 y no atravesaba el contrato HTTP real.
- **Corrección prevista:** conservar identidad antes de enviar, transferir header exacto, recuperar por consulta al registro idempotente scoped y probar efectos PostgreSQL sin segundo POST ante incertidumbre. No cambiar pricing.

- **Evidencia final V306:** QUOTE_CREATION_RECOVERY_V306.md;11 archivos recompuestos, Go/web, PostgreSQL3,68 fases/4 proyectos y agenda4. Logs fallidos conservados; no promoción integral.

### FAIL-20260907-467 — diagnóstico de continuidad presupone rutas

- **Estado:** `REGRESSION_PROVEN`
- **Clasificación:** ENVIRONMENT / diagnóstico; recurrencia FAIL346/LIB-FAIL-2134.
- **Observado V306:** rutas inferidas PORTAL*V305, gate.ps1, MARKDOWN_PROJECT_COMPOSITOR y internal/modules ausentes; glob nativo de rg inválido, Go fuera del PATH y lecturas extensas truncadas.
- **Corrección:** inventario literal y refs del checkpoint; compositor real MARKDOWN_COMPOSITOR_CORE, servicio internal/franchisejourney, runtime local fijado. No modificar PATH global ni asumir que salida truncada fue leída.
- **Evidencia parcial:** validator reconstruido15 archivos, resume37 PASS y composición67/745 a destino ausente; gates del delta pendientes.

- **Evidencia final V306:** QUOTE_CREATION_RECOVERY_V306.md;11 archivos recompuestos, Go/web, PostgreSQL3,68 fases/4 proyectos y agenda4. Logs fallidos conservados; no promoción integral.

### FAIL-20260907-465 — launcher posterga readiness hasta apagar PostgreSQL

- **Estado:** `REGRESSION_PROVEN`
- **Evidencia:** postgres-control.ps1 espera sólo PID pg_ctl, timeout15s; start/pg_isready/stop del mismo cluster PASS en0.80s. Datos conservados; no prueba de recuperación empresarial.
- **Observado V305:** Start-Process -Wait dejó la consola pendiente durante la vida del servidor; tras apagar el cluster, el pg_isready encadenado respondió no response/exit1. La disponibilidad durante gates se comprobó por un probe independiente PASS.
- **Corrección prevista:** esperar sólo al proceso pg_ctl mediante WaitForExit acotado, luego pg_isready inmediato; verificar un ciclo start/stop del mismo cluster sintético sin borrar datos.

### FAIL-20260907-464 — probe público usa la ventana del portal interno

- **Estado:** `REGRESSION_PROVEN`
- **Evidencia:** browser-final-r3.log pasa64 fases/4 proyectos; retiro de2 cupos y rechazo de2 reservas HTTP por proyecto, cero efectos. Sólo se corrige rango del fixture.
- **Observado V305:** browser-final-r2.log recibe400 al consultar cupos con90 días; PublicAppointmentSlots admite31 y exige from futuro. Falla de fixture, no regresión del arreglo PostgreSQL.
- **Corrección prevista:** intervalo futuro de14 días, dentro del contrato público, repetir matriz sin ampliar límites del servicio.

### FAIL-20260907-463 — fixture de reserva usa clave idempotente demasiado corta

- **Estado:** `REGRESSION_PROVEN`
- **Evidencia:** baseline válido confirma FAIL462, focales10 y suite PostgreSQL final3 PASS; constraint intacta y logs rojos preservados.
- **Observado V305:** availability-booking-red.log demuestra publicación indebida en dos casos, pero la clave "key" viola23514 antes de probar reserva/concurrencia. No atribuir esas reservas al producto.
- **Corrección prevista:** clave sintética compatible con el contrato existente y repetir baseline en log separado; no alterar la restricción.

### FAIL-20260907-462 — cupos y reservas no consultan disponibilidad vigente

- **Estado:** `REGRESSION_PROVEN`
- **Evidencia:** publicación/reserva con jornada vigente, lock cooperativo antes de sentencia y Read Committed explícito. Cuatro casos repetidos10, suite PG3,64 fases/4 proyectos y seis archivos reconstruidos PASS. No cubre writes SQL arbitrarios ni toda asignación de recursos.
- **Observado V305:** PublicAppointmentSlots y RequestAppointment sólo consultan el estado/capacidad del slot; CancelAvailability cambia la jornada sin cerrar slots. La creación del slot sí exige disponibilidad vigente.
- **Hipótesis a reproducir:** tras cancelar la jornada sin citas, un slot abierto podría seguir publicándose y aceptar una reserva; también falta exclusión común entre reserva y cambio de disponibilidad.
- **Reproducción:** availability-booking-red-valid.log confirma publicación de1 cupo y reserva/outbox/idempotencia indebidos (3 efectos) tras cancelación o indisponibilidad; dos carreras terminan con2 éxitos/0 conflictos. Datos sólo sintéticos.
- **Siguiente:** regresión PostgreSQL con owners reales antes de afirmar impacto; corregir en el owner de agenda y probar concurrencia, replay y consultas, sin inventar reglas comerciales nuevas.

### FAIL-20260907-461 — fixture de disponibilidad omite import fmt

- **Estado:** `REGRESSION_PROVEN`
- **Evidencia:** compilación/focal y suite Go completa PASS, harness ejecutado en64 fases/4 proyectos; import estándar, sin dependencia nueva.
- **Observado V305:** TestAvailabilityAndCustomerCancellationBindActorAndScope no compila por undefined: fmt al generar IDs sintéticos; no prueba del producto ejecutada.
- **Corrección prevista:** import estándar explícito y repetir build/gate conectado; no cambia dependencias ni esquema.

### FAIL-20260907-460 — fixture de atomicidad usa event_id no UUID

- **Estado:** `REGRESSION_PROVEN`
- **Evidencia:** V304 usa UUID del generador existente; atomicidad3 aislada más3 junto al owner completo, desde fuente canónica PASS. Log fallido preservado.
- **Observado V304:** lead-atomicity.log conserva tres rechazos22P02 al primer evento, antes de probar rollback; no es defecto del producto.
- **Corrección prevista:** generar UUID con el generador existente y reutilizar exactamente el primero para provocar el conflicto único; repetir la prueba sin debilitar el oráculo.

### FAIL-20260907-459 — cambios de lead no conservan actor HTTP en outbox

- **Estado:** `REGRESSION_PROVEN`
- **Evidencia V304:** cuatro eventos/actores autenticados por cada uno de4 proyectos; rollback conjunto del cambio/outbox y pool nuevo, suite PostgreSQL3 PASS. Métodos HTTP auditados fail-closed; los métodos legacy internos no heredan este claim.
- **Observado V304:** handlers AssignLead/TransitionLead descartan Subject; eventos conservan organización/destinatario pero no actor autenticado. No equivale a bypass ni prueba de corrupción.
- **Reproducción local:** primer assignment de operator-lead atraviesa HTTP/PG en el baseline; query read-only posterior registra1 evento lead.assigned y0 con actor_subject. Mismo efecto cuya respuesta fue abortada por el test de UI.
- **Corrección prevista:** métodos auditados compatibles en los mismos owners, actor exclusivamente del principal HTTP, sin aceptar identidad de actor del JSON. Verificar eventos PostgreSQL, permisos, concurrencia y recuperación del mismo journey.

### FAIL-20260907-458 — instalación del gate anidado resolvió workspace padre

- **Estado:** `REGRESSION_PROVEN`
- **Evidencia V303:** --ignore-workspace/offline/frozen-lockfile; CLI comprobado, gate56 fases/4 proyectos más agenda4 PASS desde reconstrucción canónica.
- **Observado V303:** pnpm install desde el directorio Playwright terminó0 pero resolvió el workspace web padre; el gate posterior no encontró @playwright/test/cli.js y no ejecutó navegador.
- **Corrección:** instalar el paquete anidado con --ignore-workspace, lock exacto y offline; verificar CLI presente antes de repetir. No es fallo de la aplicación ni PASS de integración.

### FAIL-20260907-457 — comandos del operador conservan respuesta incierta engañosa

- **Cierre parcial V311:** checklist completada/presentación con consulta exacta, respuestas/actor originales y estado vigente;88 fases/4. Familia OPEN para recibir y decidir retornos.

- **Cierre parcial V310:** publicación de checklist con identidad natural, huella y consulta exacta; actor atómico,84 fases/4 y versiones inmutables. Familia OPEN para completar checklist y recibir/decidir retornos.

- **Cierre parcial V309:** publicación de turnos con clave/lookup/actor y conservación de restricciones;80 fases/4,3 cupos/keys/eventos/actores por proyecto. Familia OPEN para publicar/completar checklist y recibir/decidir retornos.

- **Reproducción válida V309:** jornada sintética activa; POST de turno201 y replay idéntico409 con1 cupo durable. Falta recuperar recibo por clave y actor HTTP en evento; slot-red-valid-fixture.log. No es duplicación de cupos ni bypass de capacidad.

- **Cierre parcial V308:** recursos con clave/consulta scoped/actor y habilidades atómicos;76 fases/4,3 recursos/keys/eventos/actores y6 skills por proyecto. Familia sigue OPEN para capacidad/checklist/recepción/disposición; no repetir recursos/intervalos/cotización ya demostrados.

- **Cierre parcial V307:** alta de intervalos con clave/lookup scoped;72 fases/4,3 intervalos/keys/eventos/actores por proyecto y cancelación no resucitada. Familia OPEN para recursos/capacidad/checklist/recepción/disposición. No repetir el fix demostrado sin delta.

- **Reproducción V308:** dos POST autenticados idénticos de service-bay con la misma Idempotency-Key devuelven201/201, IDs distintos y2 recursos persistidos; resource-red.log. El evento tampoco incluye actor_subject por inspección del owner. Continúa OPEN hasta corrección canónica y gates de recuperación.

- **Reproducción V307:** create-availability por UI/BFF/HTTP/PG responde201 con ID y state active antes de abortar la respuesta. availability-create-red.log y trace conservados; falla el status de recuperación inexistente. Continúa la corrección sobre el owner de disponibilidad, con clave retenida y lookup exacto, sin cambiar horarios ni reglas de agenda.

- **Cierre parcial V306:** cotización conserva clave UI/BFF/HTTP y consulta scoped sin reenviar;68 fases/4,3 cotizaciones/keys/eventos/actores por proyecto. Familia sigue OPEN para creación de intervalos, recursos, checklist y recepción/disposición. No repetir los fixes de cotización sin delta.

- **Estado:** `REGRESSION_PROVEN`
- **Observado en código V302:** FranchiseCommandPanel.send todavía muestra «No se aplicó» y vuelve a habilitar formularios ante un fetch/JSON fallido. Puede ser pérdida posterior al commit, como el defecto ya corregido en recepción cliente V300.
- **Reproducción V303:** resolución correct-and-represent por BFF/HTTP/PG devolvió200/resolved con sucesor prepared antes de abortar la respuesta. Baseline mostró «No se aplicó: Failed to fetch» y Registrar resolución habilitado; la regresión de recuperación falló. resolution-red.log y trace preservados.
- **Scope de corrección en curso:** resolver de discrepancias, tres decisiones existentes, marker de sesión sin notas, fence y GET con estado durable. El sender genérico para otros comandos no queda readmitido por este tramo.
- **Cierre parcial V303 demostrado:** resolución3/sucesor1/autorizaciones2/eventos3 por cada uno de4 proyectos;56 fases y agenda4, Go/web y rebuild4/4 PASS. Mensaje incierto, GET503 y recuperación, JSON vacío, storage sin POST, carrera200/409 y permisos403. La familia permanece OPEN para los otros comandos; no marcar cobertura integral.
- **Siguiente paso ejecutable T2804:** reproducir post-commit mediante comando autorizado de fixture, conservar estado incierto/fence y resolver consulta del resultado antes de reenvío. No afirmar que V302 prueba esa recuperación: sólo corrige permisos y consultas de secciones.
- **Límite:** no se demostró aún un doble efecto ni un bypass HTTP. No usar este sender como garantía de resultado de mutaciones hasta cerrar su regresión.
- **Reproducción V304:** asignación HTTP200/versión2 confirmada antes de abortar; baseline falla al no ofrecer el status incierto en la tarjeta del lead. lead-red.log y trace preservados. Asignación/transiciones se corrigen separadamente de los comandos genéricos restantes.
- **Cierre parcial V304:** asignación y tres transiciones por cada proyecto,60 fases/4 proyectos más agenda4; marker mínimo, GET503 y recuperación, storage sin POST, JSON vacío, carrera200/409 y permisos. Los otros comandos permanecen OPEN; siguiente cancelación de intervalo existente. LEAD_COMMAND_RECOVERY_V304.md.
- **Reproducción V305:** cancelación de intervalo HTTP200/cancelled/versión2 verificada por el test antes de abortar respuesta; baseline falla por status de recuperación ausente. availability-red.log y trace preservados; corregir sólo cancelación, sin readmitir creación de intervalos.
- **Cierre parcial V305:**4 cancelaciones y1 rechazo por cita activa por proyecto;64 fases/4, recuperación GET503, storage sin POST, permisos y auditoría,6 archivos recompuestos PASS. Sigue OPEN para cotización/creación de intervalos y los otros comandos. Ver AVAILABILITY_CANCELLATION_RECOVERY_V305.md.

- **Cierre V312 del scope inventariado:** receive-return/decide-return reemplazan el último sender genérico;92 fases/4, consulta exacta aun sin lista, evidencia/actor/digest y cero reenvío. No equivale a T2804 completo, multi-tab ni ejecución downstream. Cierres parciales anteriores conservan su corte histórico.

### FAIL-20260907-456 — aserción SQL de fixture con escape incorrecto

- **Estado:** `REGRESSION_PROVEN`
- **Evidencia V302:** consulta parametrizada sobre tabla real del owner; matriz52 fases final PASS y cuatro controles de cero efectos. Ambos logs fallidos preservados.
- **Observado V302:** la UI candidata pasó su fase de roles/recuperación, pero el control SQL final usó comillas con barra dentro de un raw string Go y PostgreSQL rechazó42601. No es defecto del backend de producto ni prueba final PASS.
- **Corrección prevista:** valor parametrizado en la consulta de evidencia, sin cambiar el oráculo de cero efectos.
- **Segunda observación:** al corregir el escape, las cuatro fases UI pasaron pero PostgreSQL rechazó42P01 porque la aserción supuso sales.delivery_checklist. Recuperar el nombre desde PublishDeliveryChecklist y repetir sin quitar el control durable; browser-final.log conserva este fallo.

### FAIL-20260907-454 — cluster de prueba reanudado sin su puerto fijado

- **Estado:** `REGRESSION_PROVEN`
- **Evidencia V302:** arranque con host/puerto exactos, pg_isready y suites52 fases más agenda4 PASS.
- **Observado V302:** pg_ctl start terminó0 pero sin -o del puerto51913; el primer gate no pudo conectar y no ejercitó la UI. Log sections-red.log preservado.
- **Corrección:** detener sólo el cluster propio y arrancarlo con host loopback/puerto exactos; verificar pg_isready antes de repetir en otro log.

### FAIL-20260907-455 — secciones del portal no separan permisos ni fallos de consulta

- **Estado:** `REGRESSION_PROVEN`
- **Evidencia V302:** baseline muestra formulario sin permiso; candidato reconstruido pasa52 fases/4 proyectos y agenda4. Fallo503 recuperado por GET, cero POST automático, HTTP403 y cero checklist prohibido. No prueba bypass previo ni recuperación de comandos (FAIL457).
- **Observado en código V302:** el panel se monta por lead:read y renderiza formularios handover/resource/availability sin permisos específicos; Promise.all propaga un fallo de discrepancias a toda la página, incluidos pedidos independientes.
- **Límite:** todavía no demuestra bypass HTTP ni pérdida de datos; reproducir interfaz y preservar autorización backend. Regresión de navegador en preparación.

### FAIL-20260907-453 — histograma oculta overflow y permite métricas corruptas

- **Estado:** `REGRESSION_PROVEN`
- **Observado V301:** cuatro regresiones sobre0.2.0 fallan: valor200 reporta límite p100=100; p0 selecciona bucket vacío; NaN contamina count/sum; count MaxInt64 desborda a negativo.
- **Evidencia:** histogram-red.log, staging elite-v301-33aaeda4aa534c09bb127dd854f5c43c. Sólo candidato aislado, no integrado al perfil.
- **Corrección prevista:** overflow separado, configuración validada, rechazo atómico de entradas/overflow, consulta sin bucket vacío y snapshot coherente. Preservar CANDIDATE; no afirmar observabilidad integral.
- **Evidencia:** V301, fuente0.3.0 reconstruida2/2;24 tests y17 semillas tres veces PASS, tres targets fuzz10s, vet/build PASS. Rechazo CANDIDATE conserva destino ausente. -race no ejecutado: CGO_ENABLED=0, gcc/clang no resueltos. No promoción/integración productiva ni SLO.

### FAIL-20260907-452 — recepción/discrepancia afirma no registrada ante respuesta incierta

- **Estado:** `REGRESSION_PROVEN`
- **Observado V300:** CustomerHandoverActions usa `No se registró` y habilita de nuevo los formularios ante cualquier error de fetch/JSON; también ante una respuesta perdida después de un commit. Es el patrón ya corregido para cotizaciones, ahora encontrado en entrega.
- **Hipótesis a reproducir:** abortar la respuesta después de la aceptación/rechazo real por BFF/HTTP/PG; la interfaz no debe afirmar que no ocurrió ni permitir reenvío antes de consultar.
- **Reproducción:** action-red.log recibió200/accepted desde el backend antes de abortar la respuesta; la UI mostró `No se registró: Failed to fetch`. La regresión esperaba resultado incierto y falló. Fixture local, no entrega productiva.
- **Corrección prevista:** fence síncrono, recibo validado y resultado incierto con recuperación GET, sin aceptación/rechazo automático. No cambia política comercial ni habilita entrega inicial.
- **Evidencia final:** cuatro proyectos browser pierden respuestas200 después de aceptación y rechazo reales por BFF/HTTP/PG; incertidumbre, botones bloqueados, dos POSTs, estados accepted/rejected versión3 y un evento por acción. Recuperación GET sin reenvío. Rebuild7/7, Go y web PASS; V300. No prueba dedicada de timeout ni todos los recibos malformados.

### FAIL-20260907-451 — fixture intenta alterar una excepción inmutable

- **Estado:** `REGRESSION_PROVEN`
- **Observado V300:** al explorar la lectura de excepciones se intentó cambiar su cliente con UPDATE; el trigger existente lo rechazó correctamente con23514. No demuestra fuga ni defecto de ese trigger.
- **Corrección prevista:** conservar la inmutabilidad y probar un registro sintético inconsistente por INSERT con FK válidas en el tenant descartable, al final de la suite para no alterar su cardinalidad previa. No desactivar guards del producto.
- **Evidencia final:** fixture corregido y suite completa PG repetida tres veces PASS; conserva las aserciones de cardinalidad y no desactiva el trigger. V300.

### FAIL-20260907-450 — proyecciones de entrega exponen relaciones incoherentes

- **Estado:** `REGRESSION_PROVEN`
- **Observado V300:** tres lecturas CustomerJourney y dos Operations devuelven datos con acta/pedido/cliente/stock incoherentes; V299 protegía comandos, no estas lecturas. Baseline read-red.log preservado en staging elite-v300-ef0b29864b3441b3ba76a0859b6332d6.
- **Extensión demostrada:** exception-import-red.log reproduce dos lecturas de una excepción con cliente distinto al del acta (incluye nota sintética). Son siete lecturas indebidas en total; no se cuentan como bypass los rechazos correctos del trigger.
- **Corrección prevista:** validar relaciones en la proyección, devolver error sin datos parciales y usar un snapshot read-only común en CustomerJourney. Acotar items de checklist a las actas seleccionadas; comprobar restauración de lectura y clientes no afectados.
- **Evidencia final:** siete negativos de lectura, snapshot con escritor concurrente, consulta acotada y recuperación; suite PG tres veces PASS. HTTP401/403/403/500 sin datos privados y recuperación200. Ocho fases browser corrupt/restored desde PG y regresión completa48 fases PASS. V300; sin resultados parciales.
- **Límite:** no crea entrega inicial ni resuelve reglas financieras; no cierra assurance integral.

### FAIL-20260907-449 — checkpoint sin incrementar revisión del cursor

- **Estado:** `REGRESSION_PROVEN`
- **Observado V299:** checkpoint rechazado con `state revision must advance exactly once from event count`; el candidato conservaba revisión30 aunque se solicitó EVT-0031. No se escribió el evento ni se falsificó un PASS.
- **Corrección:** incrementar explícitamente a31 antes de sellar; conservar los30 eventos anteriores. Validación resume de31 y verificador raíz son los oráculos, no la ejecución del comando rechazado.

### FAIL-20260907-448 — regresión sintética cambia cardinalidad del journey

- **Estado:** `REGRESSION_PROVEN`
- **Observado V299:** los nueve negativos nuevos pasan después del fix, pero agregan nueve actas sintéticas al fixture existente; su aserción final conservaba cuatro y ahora observa trece.
- **Corrección prevista:** contar exactamente las cuatro actas originales y las nueve del conjunto negativo; no reducir la comprobación ni ocultar actas. Conservar regresiones y repetir el journey completo.
- **Evidencia:** aserción trece/nueve y suite completa del owner PASS tres veces desde fuente canónica reconstruida. V299; ningún filtro retira actas del producto.

### FAIL-20260907-447 — handover valida identidad aislada sin coherencia pedido/stock

- **Estado:** `REGRESSION_PROVEN`
- **Observado V299:** nueve negativos PostgreSQL fallan: completar checklist, aceptar y rechazar permiten pedido de otra organización, cliente del pedido distinto o stock de otra organización. FK tenant-scoped no prueba esos vínculos; se escriben estados/outbox con relaciones incoherentes.
- **Evidencia inicial:** staging elite-v299-8b1c2130f91d4abab8a34b3049914528, relational-red.log; fixtures sintéticos del owner, sin datos ni efectos externos.
- **Frontera adicional:** resolution-red.log demuestra que resolver una excepción también crea un sucesor desde un pedido de otra organización. El primer caso altera el fixture; los otros dos observan ese efecto previo, no prueban por separado tres bypass. Se extiende la misma defensa transaccional a ResolveDeliveryException.
- **Corrección prevista:** validar y bloquear relaciones del mismo handover/pedido/stock dentro de la transacción antes de efectos; preservar el circuito de excepciones. No inferir políticas financieras ni autorización de entrega inicial.
- **Regresión demostrada:**12 negativos, ausencia de efectos, dos locks reales,16 aceptaciones concurrentes con un efecto, rollback por outbox y pool nuevo; tres ejecuciones PostgreSQL, Go test/vet/build y36 fases browser PASS. Fuente canónica2/2 reconstruida. Pool nuevo no prueba reinicio de PG/PITR. No promoción integral ni cierre de entrega inicial. V299.

### FAIL-20260907-446 — ruta de pago previa permite eludir exclusión del intento inicial

- **Estado:** `REGRESSION_PROVEN`
- **Observado:** revisión cruzada V298 encuentra que RecordOrderPayment rechaza otra clave para el mismo pedido pero el método previo RecordPaymentIntent no aplica esa exclusión. Ambas rutas HTTP solicitan el total; el frontend no puede ser la única protección.
- **Corrección prevista:** probar llamada por la ruta previa con clave distinta y llevar la exclusión al método transaccional compartido, preservando replay de la clave original. No habilitar pagos parciales ni política de segundo intento por esta corrección.
- **Demostrado:** baseline202 reproducido, regresión409 en ruta previa; mismo key conserva ID. Fuente final reconstruida, Go/PG/vet/build y36 fases navegador PASS. V298.

### FAIL-20260907-445 — WebKit: Enter antes del estado interactivo tras recargar

- **Estado:** `REGRESSION_PROVEN`
- **Observado:** V298 reconstruido pasa tres navegadores; WebKit encuentra status vacío al esperar respuesta perdida. La secuencia focus/Enter tras reload no esperaba botón habilitado; focus no aplica las comprobaciones de click. No se demuestra por este fallo una petición enviada ni un cobro.
- **Hipótesis/corrección:** esperar toBeEnabled y toBeFocused antes de Enter. Mantener teclado, abort post-commit, timeout y oráculos. Repetir sobre fuente canónica y declarar resuelto sólo si el evento/receipt durables se comprueban.
- **Demostrado:** espera de acción explícita en test canónico; WebKit completa request, respuesta perdida post-commit, recuperación y dos intents/eventos.36 fases finales PASS; no se cierra el fallo Firefox423 previo por esta evidencia.

### FAIL-20260907-444 — solicitudes de pago montadas sin selección de proveedor

- **Estado:** `REGRESSION_PROVEN`
- **Observado:** endpoint previo recibe provider_code del cliente y no exige selección del entorno. No prueba ni habilita un proveedor real, pero no expresa correctamente integración diferida.
- **Corrección prevista:** configuración backend explícita/allowlist y ausencia cerrada; UI guiada por el snapshot del backend; comando inicial deriva importe/moneda del pedido y rechaza otro intento inicial. Probar configuración ausente/inválida, permiso, proveedor inyectado, concurrencia y recuperación desde navegador. No se elige proveedor del usuario: stripe sólo en fixture local sin llamadas externas.
- **Demostrado:** main1.9.2 y Commerce0.6.0, siete casos selección/seis rutas ausente-inválido; cuatro fases de desactivación en navegador mantienen pedidos/intents y rechazan nuevos requests. V298 reconstruido12/12; guía ligada al consumer real.

### FAIL-20260907-443 — fixture V297 no respeta tipos/código de tenant y scope JSON

- **Estado:** `REGRESSION_PROVEN`
- **Observado:** fixture inicial mezcla parámetros UUID/text y usa código que comienza con dígito; PostgreSQL lo rechaza antes de probar producto. Reemplazo de org alteraba el nombre organization_id en vez de su valor.
- **Corrección:** cast UUID explícito y tenant_code con prefijo test; sustituir el literal JSON completo. No rebajar DDL ni aceptar 400 como prueba de permiso 403.
- **Demostrado:** fixture canónico reconstruido, seis negativos HTTP y suite PG V297 PASS.

### FAIL-20260907-441 — reintento HTTP de pago pierde el recibo durable

- **Estado:** `REGRESSION_PROVEN`
- **Observado:** CreatePaymentAttempt inserta con clave única pero propaga el conflicto de PostgreSQL; no recupera el identificador ya persistido ni distingue misma intención de payload divergente. El evento tampoco conserva el actor HTTP.
- **Corrección prevista:** recuperación transaccional por tenant/proveedor/clave, igualdad del intent original, autorización de organización, actor derivado del token y un solo outbox. Pruebas HTTP/PG concurrentes, respuesta perdida y reconexión; sin cobro real.
- **Demostrado:** baseline400 reproducido;16 replays y16 creaciones concurrentes, recibo descartado/recuperado, outbox fallido sin escritura parcial y reconexión con estado actual PASS. Handlers con principal sintético, no nueva prueba OIDC/navegador. Evidencia V297 y reconstrucción7/7.

### FAIL-20260907-442 — referencia NULL impide dos pagos aún no enviados

- **Estado:** `REGRESSION_PROVEN`
- **Observado:** migración 0003 declara UNIQUE NULLS NOT DISTINCT sobre tenant/proveedor/provider_reference. Dos pedidos diferentes sin referencia del proveedor colisionan aunque sus claves sean distintas.
- **Corrección prevista:** migración aditiva que preserve unicidad de referencias no nulas y permita múltiples intents sin referencia. Downgrade debe fallar sin borrar datos si el estado ya no admite la restricción anterior. Demostrar ambas direcciones y conflictos reales.
- **Demostrado:** segundo pedido400 antes de0053 y202 después; referencias duplicadas no nulas rechazadas; up/down/up vacío y down incompatible sin pérdida PASS. Base nueva53 migraciones; snapshot tras reinicio idéntico. V297.

### FAIL-20260907-440 — reserva HTTP no conservaba el actor autenticado

- **Estado:** `REGRESSION_PROVEN`
- **Observado:** la reserva canónica guardaba línea/unidad en outbox pero el endpoint no transfería Subject. Permitía demostrar efecto, no quién lo ordenó.
- **Corrección:** AllocateStockAs deriva actor exclusivamente del principal verificado y usa la misma transacción/tabla. El fallback sin soporte auditado rechaza la operación; no acepta actor del JSON ni duplica inventario. Los callers internos previos conservan su contrato y requieren revisión propia de actor.
- **Regresión requerida:** navegador→HTTP→PostgreSQL con outbox actor_subject=operator, fuente canónica y reconstrucción exacta. Sin identidad personal real en fixtures.
- **Demostrado:** 24 fases desde reconstrucción11/11; cada fase operativa exige dos eventos con actor verificado. V296 final PASS, snapshot tras reinicio PG idéntico.

### FAIL-20260907-439 — selector de alerta incluye announcer Next

- **Estado:** `REGRESSION_PROVEN`
- **Observado:** primera fase stock V296 probó dos reservas/outbox durables; fase de lectura fallida encontró dos role=alert (mensaje y announcer Next) y Playwright rechazó ambigüedad.
- **Corrección prevista:** acotar selector al main del producto; no retirar strict mode ni la prueba de error. Repetir tres fases operativas sobre reconstrucción canónica.
- **Demostrado:** cuatro navegadores, tres fases operativas PASS, selector scoped y strict mode conservado, evidencia V296.

### FAIL-20260907-438 — fixture de dos unidades hereda VIN vacío único

- **Estado:** `REGRESSION_PROVEN`
- **Observado:** V296 completó las tres fases de cotización, pero al sembrar dos unidades ambas tenían VIN NULL; PostgreSQL rechazó unique nulls not distinct. La repetición con VIN distinto encontró la misma regla para battery_serial_number. No era un default vacío: el DDL 0003 lo demuestra. Las restricciones no se debilitan.
- **Corrección prevista:** asignar VIN y batería sintéticos distintos por unidad y repetir el recorrido. Corregir también nombre Organizations en el nuevo fake de identidad antes de compilar. La aplicabilidad de esa restricción a unidades sin VIN/batería sigue sin demostrar; no usar el fixture para afirmar compatibilidad universal del catálogo.
- **Demostrado:** fixture reconstruido siembra dos unidades distintas; carrera produce una sola reserva de la unidad disputada. V296 final PASS, no alteración del DDL.

### FAIL-20260907-437 — puerto local no permitido por Windows

- **Estado:** `REGRESSION_PROVEN`
- **Observado:** el cluster sintético V296 no pudo bindear 127.0.0.1:56696 (Permission denied); createdb rechazó la conexión. No se aplicaron migraciones ni hubo efectos externos.
- **Corrección prevista:** seleccionar un puerto comprobado con TcpListener loopback, comprobar exit de pg_ctl y readiness antes de createdb. Conservar initdb/postgres logs bajo staging elite-v296-a380dd37465046c4a166bf68cdc2baec.
- **Demostrado:** puerto permitido, 52 migraciones, gates PG/navegador y reinicio→snapshot idéntico→parada propia PASS. Sin servidor productivo ni datos reales.

### FAIL-20260907-436 — ledger de procedencia no contabiliza nuevo test AUTHORED

- **Estado:** `REGRESSION_PROVEN`
- **Observado:** VERIFY_LIBRARY reconstruye 1434 archivos y exige AUTHORED=1196/ADAPTED=133/VERBATIM=105; el resumen conservaba 1433 y 1195.
- **Corrección:** actualizar sólo recuento actual en THIRD_PARTY_NOTICES; no reclasificar el test como upstream. Repetir verificador y conservar evidencia del rechazo.
- **Recurrencia:** el control siguiente encontró FRANCHISE_PREFLIGHT_GAP aún en 160/1433/727/51 y 742 compuestos. Se actualizó ese owner explícito a 160/1434/728/51 y 743; no cambiar los recuentos históricos de V258/V294.

### FAIL-20260907-435 — ejemplo de guía contiene home del equipo de mantenimiento

- **Estado:** `REGRESSION_PROVEN`
- **Observado:** VERIFY_LIBRARY rechaza home local literal en PROJECT_SECRETS_TEMPLATE; no se cuenta el intento como PASS. La descripción inicial de esta lección repitió el mismo literal y también fue rechazada; ambas referencias se retiraron.
- **Corrección:** ejemplo resuelve biblioteca desde terminal situada en su raíz, sin home del mantenedor. Conserva ruta del proyecto por completar; ningún secreto en comando.
- **Regresión focal:** búsqueda de home literal en todos los Markdown vacía, incluido staging expresado mediante LOCALAPPDATA en el reporte. Verificador general repetido después de checkpoint actualizado.

### FAIL-20260907-434 — inventario omite configuración del módulo refund anidado

- **Estado:** `REGRESSION_PROVEN`
- **Observado:** escaneo de composición encuentra MERCADO_PAGO_ACCESS_TOKEN y REFUND_WORKER_ID fuera del cmd raíz. STRIPE_SECRET_KEY también tiene consumer real de devoluciones, no de cobro principal. Regex YAML aplicada a TypeScript detectó FLOW_COOKIE falsamente como entorno.
- **Corrección:** incluir módulo anidado y distinguir cobros/devoluciones; limitar interpolación de entorno YAML a YAML, conservar lectura manual de Config/SDK/IaC y probar nuevamente. No afirmar exhaustividad por buscar sólo cmd raíz.
- **Regresión:** 69 bindings documentados sobre reconstrucción 67/743; FLOW_COOKIE excluido por ser constante, no entorno. Artifact credential-consumer-inventory.json en staging V295.

### FAIL-20260907-433 — utilidad DPAPI rechaza caso de una clave y ACL sin privilegios

- **Estado:** `REGRESSION_PROVEN`
- **Observado:** primer test sintético falla por Count sobre escalar bajo StrictMode; al corregirlo, Set-Acl con DirectorySecurity nuevo requiere SeSecurityPrivilege no disponible. Ningún secreto real usado.
- **Corrección:** envolver colección en @(); modificar sólo DACL existente, sin solicitar privilegios de auditoría ni elevación; repetir pruebas de escritura cifrada, ACL, negativos e inyección aislada.
- **Regresión:** FileSystemAclExtensions.SetAccessControl aplica DACL sin SeSecurityPrivilege. Salida de hijo drenada sin exponer, no se presupone captura de Start. Suite final 16 checks PASS, sintéticos exclusivamente; sin nueva dependencia.

### FAIL-20260907-430 — guía de secretos no coincide con consumidores

- **Estado:** `REGRESSION_PROVEN`
- **Observado:** PROJECT_SECRETS_TEMPLATE mezcla environment y bóveda, declara alternativas obligatorias y nombres genéricos sin consumer; no distingue Config Go inyectada de variables leídas por ejecutables.
- **Corrección prevista:** inventario desde archivos compuestos, tipo/selección/custodia/probe y límites por integración; ninguna credencial en evidencia.
- **Regresión:** guía reemplazada, 69 bindings cotejados y revisión manual de Config/CLI/IaC/perfiles. Recetas no elegidas/inaccesibles siguen pendientes, no se consideran resueltas por este cierre local.

### FAIL-20260907-431 — arranque no fiscal exige socket ARCA

- **Estado:** `REGRESSION_PROVEN`
- **Observado:** cmd/electromobility-api/main.go exige ARCA_WSFE_SOCKET antes del arranque y monta FiscalModule siempre. Contradice la decisión explícita del usuario de diferir ARCA.
- **Corrección prevista:** selección fiscal explícita, rutas/provider ausentes al diferir, rechazo de configuración incompleta al activar; pruebas sin servicios fiscales.
- **Regresión:** main1.9.1 reconstruido; ocho escenarios de selección/ruta ausente/401 al activar y ruta independiente en mux sintético. No es nuevo journey PostgreSQL/navegador ni activación fiscal live.

### FAIL-20260907-432 — diagnóstico con paths supuestos y salida excesiva

- **Estado:** `REGRESSION_PROVEN`
- **Recurrencia V305:** lectura combinada truncada, path availability.go supuesto y glob Windows inválido; recuperar protocolo por tramos y owner inventariado franchisejourney_availability.go. No usar intentos fallidos como evidencia.
- **Recurrencia V304:** glob Windows inválido, salida de componente combinada truncada y patch inicial rechazado por hunk repetido; corregidos con rutas inventariadas y anclas únicas. No cuentan como prueba del producto ni lectura completa.
- **Recurrencia V302:** lecturas combinadas truncadas y ancla de patch no observada rechazaron el primer intento antes de escribir. Se usaron tramos/anchors comprobados; no se atribuye éxito a ese intento.
- **Recurrencia V303:** salida combinada de contratos excedió el presupuesto y se infirió journey.go inexistente. Se recuperaron las secciones omitidas y el inventario service.go/availability.go antes de modificar el owner; no usar lecturas fallidas como evidencia.
- **Recurrencia V301:** lectura combinada de protocolos e índice produjo truncamiento; se recuperaron los tramos omitidos y las secciones materiales por separado. No usar salida truncada como lectura completa.
- **Recurrencia V301 de path:** se supuso go_native_fuzz_gate/README.md; el manifest canónico identifica go_fuzz_gate y se reconstruye su perfil dedicado. No se interpreta el archivo ausente como runner verificado.
- **Recurrencia V300:** lectura de src/app/error.tsx inexistente. Inventario confirmó la página real customer/handovers y su falta de fallback local; no se atribuyó una boundary inexistente al producto.
- **Recurrencia V299:** glob Windows, nombre salesshipment.go supuesto y lecturas combinadas truncadas. Reencauzado al archivo inventariado sales_shipment.go y lecturas por función; no se presentan como verificaciones exitosas.
- **Observado:** lecturas combinadas truncadas y paths cmd/enterprise, internal/identity, assembly.go inexistentes. No se atribuye éxito a esas lecturas.
- **Corrección:** inventario rg --files y lecturas acotadas de cmd/electromobility-api, auth y config reales; no cambios de producto por el error de diagnóstico.
- **Recurrencias V295:** workdir incorrecto, glob Windows inválido, truncamientos y patches rechazados por owner repetido o hunk sin contexto. No se cuentan como evidencia; reconstrucción y hashes validan sólo los cambios aplicados.
- **Recurrencia V296:** lecturas combinadas truncadas, glob inválido y confusión de raíz con composición durante diagnóstico; corregir con paths inventariados y lecturas acotadas. No son verificaciones ejecutadas ni cambios funcionales.
- **Recurrencia V297:** nombre supuesto de migración, búsqueda de db desde la raíz de biblioteca y salida truncada; sólo se utilizan después los paths observados en la composición. No cuentan como pruebas.
- **Recurrencia V297 posterior:** glob Windows inválido, patch sin delimitador inicial y LASTEXITCODE heredado interpretado erróneamente como fallo de script PowerShell. Patch repetido correctamente; materialización efectiva745 y comparación7/7 verificadas sin recomponer un destino existente.
- **Recurrencia V298:** lecturas combinadas truncadas y paths supuestos src/lib, command-client, config y e2e; utilizar inventario observado src/platform y microsoft_playwright_browser_gate. No se cuentan como gates ni cambian producto.
- **Recurrencia V298 checkpoint:** descripción de blocker excedió300 caracteres; el validator rechazó sellarlo. Se acorta el resumen y se conserva el detalle en owners; no se elimina el bloqueo ni se altera el evento previo.

### FAIL-20260907-429 — restart PostgreSQL retiene stdout del proceso hijo

- **Estado:** `REGRESSION_PROVEN`
- **Observado:** pg_ctl restart sin -l informa server started, pero la captura PowerShell no retorna mientras PostgreSQL hereda el pipe. El servidor sí responde al probe independiente; la comparación posterior aún no se ejecutó.
- **Corrección:** separar probe de persistencia y parada del servidor propio; en futuras ejecuciones usar -l con log explícito. No interpretar server started como evidencia de comparación de datos.
- **Recurrencias:** start sin -o volvió al puerto loopback por defecto; Start-Process -Wait siguió esperando descendientes. Se detuvieron sólo los orquestadores propios identificados y sus clusters, sin tocar servidores ajenos.
- **Regresión:** Start-Process Hidden sin -Wait, Process.WaitForExit(10000) sobre pg_ctl, -l y -o explícitos; start→probe56594→stop PASS. Comparación independiente antes/después del reinicio SHA-256 idéntica. PostgreSQL descartable quedó detenido; no borrar fixtures.

### FAIL-20260907-428 — foreach canalizado sin expresión colectora

- **Estado:** `REGRESSION_PROVEN`
- **Observado:** diagnóstico de hashes PowerShell rechazado por empty pipe element. No modificó archivos ni ejecutó la comparación.
- **Corrección:** asignar la colección producida por foreach antes de ConvertTo-Json; repetir y conservar 5/5 hashes. No confundir error de diagnóstico con fallo del código.

### FAIL-20260907-427 — formatter presupone dígitos presentes en tipo Intl

- **Estado:** `REGRESSION_PROVEN`
- **Observado:** quote-ux-build.log rechaza TS2345 porque maximumFractionDigits es number|undefined. No se publicó ni se suprimió typecheck.
- **Corrección:** validar dígitos enteros antes de BigInt/padStart; si no son verificables, no ofrecer aceptación ni inventar precisión monetaria.

### FAIL-20260907-426 — presentación de cotización expone representación técnica

- **Estado:** `REGRESSION_PROVEN`
- **Observado:** captura móvil V294 muestra unidades menores, issued/accepted, fecha ISO y botón en issued vencida. Backend rechaza vencida, pero UX permite intento inútil.
- **Corrección:** presentación server-side con Intl fijado, moneda explícita sin conversión ni redondeo; estados legibles, zona UTC explícita y guardas de vigencia/importe seguro. Agregar aserciones de navegador, sin alterar transacción.

### FAIL-20260907-425 — inventario preflight no sincronizado con expediente V294

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / CONTRACT.
- **Observado:** VERIFY_LIBRARY valida checkpoint20 pero rechaza inventario preflight: esperaba 160/1433/727/51. No es fallo del journey ni PASS general.
- **Corrección:** actualizar sólo el owner del inventario al recuento real, conservando el rechazo y repetir verificador tras sellar el siguiente checkpoint. No ampliar allowlists.
- **Regresión focal:** comprobación exacta de la cadena normalizada 160/1433/727/51 PASS; verificador general se ejecuta nuevamente tras el checkpoint final, sin atribuir PASS al primer intento.

### FAIL-20260907-424 — instalación anidada resuelta contra workspace padre

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / ENVIRONMENT; V294.
- **Observado:** install del gate devolvió Already up to date, pero el CLI Playwright no existía en el destino reconstruido. quote-rebuilt.log termina MODULE_NOT_FOUND antes de los journeys; no se cuenta como prueba funcional.
- **Causa / corrección:** pnpm detectó el workspace padre; instalar el package anidado con --ignore-workspace --offline --frozen-lockfile --ignore-scripts mantiene el lock 1.62.1 y agrega sus tres paquetes locales. Repetir el gate completo y conservar ambos resultados.
- **Recurrencia de diagnóstico:** consultas relativas a owners de biblioteca se ejecutaron desde el staging del gate y fallaron; usar workdir de biblioteca o paths absolutos para esa lectura, sin atribuir éxito.
- **Evidencia posterior:** reconstrucción completa 12 fases PASS en quote-rebuilt-corrected.log; versiones intactas.

### FAIL-20260907-423 — Firefox falla durante cierre de contexto

- **Estado:** `DIAGNOSED`
- **Severidad / clasificación:** MEDIUM / UPSTREAM_RUNTIME; V294.
- **Observado:** Firefox/restart termina con Browser.removeBrowserContext: _maybeDontRestoreTabs, this._windows[aWindow.__SSi] undefined. Ocurre en teardown después de las aserciones; el runner queda FAIL y no se cuenta como PASS.
- **Evidencia:** quote-four-browsers-green.log (nombre no implica éxito), screenshot/error-context conservados. Investigar runtime oficial y repetir sin retries ocultos ni supresión de errores; no actualizar pin automáticamente.
- **Seguimiento:** repetición focal y reconstrucción4/4 PASS; causa exacta no confirmada, no cerrar por una repetición verde.

### FAIL-20260907-422 — selector de alerta colisiona con announcer de Next

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / TEST_FIXTURE; V294.
- **Observado:** fase read-failure encuentra la alerta real y el route-announcer del framework; strict mode rechaza el selector global en cuatro navegadores. La UI sí muestra el error seguro esperado.
- **Corrección:** seleccionar alerta dentro del landmark main; no desactivar strict mode ni suprimir el announcer accesible. Preservar quote-four-browsers.log.

### FAIL-20260907-420 — cotización persistida se anuncia como no aceptada

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** HIGH / UX_CONSISTENCY; V294.
- **Reproducción:** Playwright deja completar el POST real (accepted con order_id) y pierde sólo la respuesta. La pantalla anuncia «No se aceptó: Failed to fetch». quote-baseline-ux.log conserva la aserción fallida; datos exclusivamente sintéticos.
- **Causa:** la UI interpreta error de transporte como resultado comercial negativo y vuelve a habilitar el comando.
- **Corrección:** estado incierto explícito, bloquear reenvío y recuperar por lectura; mostrar pedido durable, ayuda versionada y alcances de aceptación. No modificar reglas transaccionales.

### FAIL-20260907-421 — sanitización del harness reemplaza flags no secretos

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / TEST_EVIDENCE; V294.
- **Observado:** reemplazar todos los valores ELITE_QUOTE_* sustituye también el flag «1» y corrompe números de líneas en el log rojo.
- **Corrección:** sanitizar sólo JWTs sintéticos y AUTH_SESSION_SECRET; conservar log original y repetir sin ocultar aserciones ni eventos.

### FAIL-20260907-419 — tipo de parámetro SQL ambiguo en fixture

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / TEST_FIXTURE; V294.
- **Observado:** fixture tenant usa $1 como uuid y texto; PostgreSQL rechaza 42P08 antes del recorrido. No es fallo de producto ni evidencia del caso de respuesta perdida.
- **Corrección:** cast UUID explícito en ambos usos; repetir gate conservando quote-baseline.log.
- **Segundo rechazo del fixture:** tenant_code no admite UUID sin prefijo alfabético; usar quote- más UUID sin guiones, conforme al patrón del harness anterior. Conservar quote-baseline-after-fixture.log; no alterar constraints.

### FAIL-20260907-418 — cache de dependencias mediante junction fuera del root

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / ENVIRONMENT; V294.
- **Observado:** Next/Turbopack rechaza node_modules como junction externa; build baseline no pasó. No defecto del recorrido ni autorización para omitir build.
- **Corrección:** reinstalar exactamente el lock desde cache offline en el staging aislado; conservar log rojo y no modificar el árbol V278.

Regresiones V294 anteriores: evidencia QUOTE_ACCEPTANCE_CONNECTED_V294.md; los fallos rojos se conservan y no se cuentan como pases.

### FAIL-20260907-417 — subdirectorio de staging supuesto

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / ENVIRONMENT; V294.
- **Observado:** lectura bajo verified/backend falló: la composición V278 es plana.
- **Corrección:** inventario literal y rg --files ubicaron verified/internal/platform/httpapi/franchisejourney_test.go; lectura correcta sin alterar el baseline. No se reinicia arquitectura ni se atribuye éxito a la lectura fallida.
- **Recurrencia V294:** nombre supuesto de migración no existe; localizar archivo por inventario antes de leer schema. Aplicación real de las 52 migraciones sí terminó PASS.

### FAIL-20260907-416 — hunks en orden inverso

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / ENVIRONMENT; V293.
- **Observado:** patch de sincronización buscaba un hunk anterior después de otro posterior; rechazado sin cambios, comprobado por rg.
- **Corrección:** ordenar hunks por aparición y repetir sobre anchors leídos; no modificar el código de producto ni afirmar que el patch rechazado se aplicó.
- **Recurrencia V294:** hunk vacío posterior a la versión intentó volver al título anterior; rechazado sin aplicar. Corrección de versión con un único hunk y cabecera en orden ascendente, sin cambios parciales.

### FAIL-20260907-415 — nombre supuesto de validador

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / ENVIRONMENT; V292.
- **Observado:** búsqueda de validate_engineering_contract.py produjo archivo inexistente; no ejecutó ni modificó código.
- **Corrección:** usar el inventario materializado, cuyo ejecutable real es validate_project.py; leer sus argumentos antes de invocarlo. No inferir éxito del comando fallido.

### FAIL-20260907-414 — representación Markdown oficial no recuperable

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / AUTHORITY_ACCESS; V292.
- **Observado:** el enlace index.md publicado por OpenTelemetry devolvió Internal Error en el navegador de investigación.
- **Corrección / evidencia:** consulta de la página HTML oficial y su sección Data minimization; contenido recuperado, sin sustituirlo por fuente secundaria ni inventar acceso al Markdown. PRIVACY_POLICY_LOGGER_REPAIR_V292.md conserva URL y alcance.

### FAIL-20260907-413 — referencia simultánea de lectura y reutilización

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / CONTRACT; checkpoint de mantenimiento V291.
- **Observado:** observability-owner quedó en must_read_refs y reuse_without_reload_refs; el checkpoint lo rechazó antes de agregar EVT-0017. No se modifica la cadena previa ni el código del fix.
- **Corrección:** retirar del grupo reutilizable todas las referencias activas obligatorias, recalcular hash del ledger y repetir checkpoint/resume; no debilitar disjunción del contrato.
- **Evidencia posterior:** validación del candidato corregido PASS; el primer VERIFY conserva salida fallida porque se inició sin checkpoint válido. Repetir el verificador general sólo después del checkpoint/resume encadenados; conservar primer log como verify-library.log y nuevo resultado como verify-library-final.log.

### FAIL-20260907-412 — patch generado repite owner

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / ENVIRONMENT; sincronización V291.
- **Observado:** el patch agrupaba cada línea como operación nueva sobre el mismo archivo y apply_patch rechazó operaciones duplicadas, sin cambios parciales.
- **Corrección:** agrupar hunks por archivo dentro de una única operación. No relajar validadores ni cambiar índices a mano sin lectura literal.
- **Evidencia posterior:** cuatro owners sincronizados con una operación por archivo, sin modificar reglas de admisión; V291. Revisión/checkpoint final comprueba hashes de los owners cambiados.

### FAIL-20260907-411 — concatenación ambigua al comprobar hashes

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / ENVIRONMENT; verificación del rebuild V291.
- **Observado:** pasar una ruta entre comillas seguida de $name se interpretó como dos argumentos; Get-FileHash trató el nombre como Algorithm y rechazó la llamada. El materializador ya había producido dos archivos; no indica corrupción.
- **Corrección prevista:** construir ambos paths con Join-Path, comparar hashes y repetir suite sobre rebuild. No modificar los archivos para ajustar el hash esperado.
- **Evidencia posterior:** Join-Path, hashes 2/2 y tests/vet/build/fuzz reconstruidos PASS; COUNTER_MONOTONIC_REPAIR_V291.md.

### FAIL-20260907-410 — compositor incluido dentro de su propio plan

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / CONTRACT; ensayo de tooling V290.
- **Observado:** el plan de tres packs falla con unresolved variable en tools/test-compose-markdown-project.ps1: sus fixtures contienen placeholders literales. No es una variable empresarial faltante ni un defecto de los validators.
- **Hipótesis / corrección:** respetar bootstrap existente: materializar compositor por el extractor y componer sólo los dos validators. No modificar fixtures ni debilitar el rechazo de placeholders del compositor. Probar colisiones, reconstrucción independiente y suites.
- **Evidencia inicial:** staging elite-v290-5fcdc9bf739d417babe3691898563e2b; destino composition-a no debe contener archivos tras rechazo de preflight.
- **Evidencia posterior:** cero archivos tras rechazo; selección corregida produce 23/23 archivos idénticos en destinos independientes, 24 tests execution y 70 readiness PASS, suite compositor PASS; colisión rechazada sin cambiar 24/24 archivos. V290 conserva límites y plan corregido.

### FAIL-20260907-409 — pipeline después de foreach statement

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / ENVIRONMENT; lectura de índices V289.
- **Observado:** PowerShell rechazó sintaxis foreach {...} | ConvertTo-Json antes de ejecutar; sin cambios de archivos.
- **Corrección / evidencia:** asignar la salida del foreach a una variable y serializarla después. Lectura devuelve exactamente cuatro owners, una línea material por owner; patch generado con esas líneas completas.

### FAIL-20260907-406 — suite aislada sin raíz de autoridades

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / ENVIRONMENT; V289, invocación local de tests.
- **Observado:** unittest desde kit en Temp no recibió ELITE_AUTHORITY_ROOT; cinco fallos por authority_docs no encontrados, no por contratos de producto.
- **Hipótesis / corrección:** fijar explícitamente la raíz existente sólo en el proceso de prueba, conservando/restaurando la variable previa. No modificar tests, fixtures ni validadores.
- **Evidencia:** staging elite-v289-ad41218f69b848bba366dd0944fba52e/contract-regression.log conserva cinco fallos; contract-regression-fixed.log demuestra 11 tests PASS con raíz explícita, sin cambiar código. Mantener la raíz al probar el kit fuera del workspace.

### FAIL-20260907-407 — estado de salida equivocado para script PowerShell

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / ENVIRONMENT; V289, invocación del materializador.
- **Observado:** el script informa 15 archivos materializados, pero el caller comprueba LASTEXITCODE (de programas nativos) y lanza Rebuild failed. Ese valor no representa la finalización de un script PowerShell que no invoca exit.
- **Corrección prevista:** usar proceso pwsh hijo con salida propia y destino nuevo; comprobar todos los archivos contra la reconstrucción previa y ejecutar el validador recién reconstruido. No modificar materializador para ocultar el error del caller.
- **Evidencia posterior:** proceso hijo exit 0, dos reconstrucciones 15/15 byte-idénticas y 11 tests del kit nuevo PASS. Logs materialize-independent.log y rebuilt-regression.log en staging V289; no era un fallo del materializador.

### FAIL-20260907-408 — ancla parcial en patch del roadmap

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / ENVIRONMENT; recurrencia LIB-FAIL-2061.
- **Observado:** patch múltiple rechazado porque un prefijo no era una línea completa; se comprobó que el informe nuevo tampoco se creó, sin mutaciones parciales.
- **Corrección:** separar creación del informe y usar una línea literal completa del roadmap; no tratar prefijos como anclas. Informe V289 conserva el fallo y el cambio final.

### FAIL-20260907-405 — documentación histórica Meta no recuperada

- **Estado:** `BLOCKED_EXTERNAL`
- **Severidad / clasificación:** LOW / ENVIRONMENT; V288, investigación documental, no fallo de un servicio del usuario.
- **Observado:** onboarding-business-app-users devuelve HTTP 429 al lector web; la página propuesta Conversations API es rechazada por el lector como no recuperable. No se obtuvo el contrato primario de importación histórica.
- **Límite / reapertura:** no afirmar ventanas, elegibilidad, schemas, permisos ni cobertura total. Reabrir con documentación oficial accesible y canal/cuenta identificados; no sustituir por un blog ni scraping de sesiones personales. La evaluación conceptual y el inventario local siguen disponibles.

### FAIL-20260907-404 — diagnóstico textual demasiado amplio

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / ENVIRONMENT; V288, recurrencia de límite de salida.
- **Observado:** resultados rg incluyeron líneas extensas del lock y truncaron la salida. No se dedujo cobertura ni se editó código desde el segmento truncado.
- **Corrección / evidencia:** búsqueda posterior por nombres con rg -l y lectura de 38 líneas exactas recupera Presidio 2.2.364 y su owner/revisión/licencia; History se inspeccionó directamente. CHAT_HISTORY_LEARNING_AUDIT_V288.md. No equivale a tests de producto.
- **Recurrencia V289:** lectura conjunta de tres manuales excedió el presupuesto de salida; aumentar sólo el límite interno no eliminó el límite externo. Recuperación por rangos de un archivo y emisión de `.output` sin encapsulado JSON. No se modificó código documental ni se afirmó lectura completa del segmento omitido. Los contratos/schema y evidencia usados para el cambio sí fueron leídos sin truncamiento. Prevención: un archivo/rango material por salida, menos de 5.500 tokens.
- **Recurrencia V290:** arranque agrupado volvió a truncar manuales; se recuperaron los rangos materiales de inicio/readiness y el protocolo de composición en salidas menores. No se afirmó relectura completa de los segmentos omitidos. Usar el cursor validado y el índice antes de pedir un manual extenso.
- **Recurrencia V292:** salida de arranque extensa truncada; se reutilizó el checkpoint previamente validado y se leyeron por separado el owner, fallo y contratos materiales. No se atribuye lectura íntegra a los segmentos omitidos; limitar salidas por rango y emitir sólo texto relevante.
- **Recurrencia V293:** búsquedas agregadas excedieron límites; se recuperaron aparte Overview SQL, ayuda versionada y ruta accept-quote. La auditoría se declara estática/focal, no lectura íntegra ni equivalencia completa. Ausencia de coincidencia no demuestra ausencia funcional.

### FAIL-20260907-403 — fixture NEW declara delta de proyecto existente

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / CONTRACT; V287, fixture AUTHORED.
- **Observado:** checkpoint rechaza `NEW project cannot declare an existing-project delta scope`, exit 2. El control funcionó; el fixture asignaba delta_scope_ref para ambos modos.
- **Corrección prevista:** NEW conserva baseline vacío observado y delta_scope_ref null; EXISTING conserva inventario y delta. No debilitar validator.
- **Evidencia:** 55 checks NEW/EXISTING, 28 procesos; repetidos en copia limpia 13/13. AGENT_ENTRY_LIFECYCLE_V287.md.

### FAIL-20260907-402 — ensayo integrado recibe dos Python de PATH

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / ENVIRONMENT; V287, test AUTHORED.
- **Observado:** Get-Command python -CommandType Application devuelve Python314 y WindowsApps; el parámetro string los concatena y no es ejecutable.
- **Corrección prevista:** elegir explícitamente el primer Application observado igual que el resolver canónico; no instalar ni modificar PATH. Repetir flujo NEW/EXISTING completo.
- **Evidencia:** selección única y 55 checks/copia limpia PASS. AGENT_ENTRY_LIFECYCLE_V287.md.

### FAIL-20260907-401 — lectura diagnóstica excede límite de salida

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / ENVIRONMENT; V286, diagnóstico propio.
- **Observado:** lectura conjunta de owners extensos truncada por transporte; JSON.parse rechazó el prefijo Warning. No se aplicó ningún edit desde datos truncados.
- **Corrección:** patches locales de líneas exactas y append de filas; obtener sólo hash o línea seleccionada para owners grandes. Las escrituras siguen apply_patch; no releer el corpus entero ni ampliar el contexto.

### FAIL-20260907-400 — suite conserva exit code de un negativo esperado

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / CODE; V286, test AUTHORED.
- **Observado:** copia limpia muestra 46 checks PASS, pero caller recibe LASTEXITCODE=1 heredado del último child-process RequireAll, cuyo rechazo era esperado; el caller correctamente rechaza continuar.
- **Corrección prevista:** la suite retorna exit 0 explícito sólo después de todas las aserciones y cleanup; repetir copia limpia con caller que exige exit 0.
- **Resultado:** copia limpia y caller exit 0; 46 checks PASS. ISOLATED_TOOLCHAIN_RESOLUTION_V286.md.

### FAIL-20260907-399 — fixture CLI compone incorrectamente argv

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / CODE; V286, test AUTHORED.
- **Observado:** regresión de proceso rechaza `Informational toolchain probe failed`; ejecución directa con siete paths ausentes produce correctamente TOOLS_MISSING/7, exit 0. La concatenación dentro del array de argumentos debe reemplazarse por dos elementos explícitos.
- **Corrección prevista:** interpolar únicamente el nombre de parámetro y mantener path como argumento separado; conservar stdout del fallo en la aserción y repetir child-process tests.
- **Resultado:** argv separado; procesos reales y 46 checks PASS desde copia limpia. ISOLATED_TOOLCHAIN_RESOLUTION_V286.md.

### FAIL-20260907-398 — rutas aisladas sin interfaz uniforme en preflight

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / ENVIRONMENT; T2801/T2808, V286, runner AUTHORED.
- **Observado:** test real rechaza `Missing explicit tool parameter: NodeExecutable`; tampoco existen PnpmExecutable, DockerExecutable o PsqlExecutable. El PostgreSQL local disponible fuera del PATH no puede seleccionarse como Go/dotnet. Python tiene un fallback a python3 incluso si se proporcionó una ruta explícita inválida.
- **Corrección prevista:** parámetros uniformes, selección explícita sin fallback silencioso, test de todos los consumers; probe separado de disponibilidad sin presentar ausencia de audit como PASS integrado.
- **Resultado:** 46 checks PASS, probe real con 7/8 disponibles, sin promoción ni cambios de PATH global. ISOLATED_TOOLCHAIN_RESOLUTION_V286.md.

### FAIL-20260907-397 — glob de ruta pasado directamente a rg en Windows

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / ENVIRONMENT; V286, diagnóstico propio.
- **Observado:** rg rechazó reconstruction_evidence/ARCA* con error 123. El inventario previo sí localizó los cuatro ejecutables PostgreSQL; la búsqueda fallida no cuenta como ausencia de evidencia.
- **Corrección prevista:** buscar en el directorio literal usando -g para seleccionar archivos; no mezclar wildcard de shell con paths del programa.
- **Resultado:** búsqueda corregida localiza los expedientes oficiales Microsoft ARCA V123–V129; recurrencia conservada, sin inventar ausencia. ISOLATED_TOOLCHAIN_RESOLUTION_V286.md.

### FAIL-20260907-396 — cursor de reanudación supera su presupuesto

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / CONTRACT; V285, expediente propio.
- **Observado:** el checkpoint rechazó next_action de más de 800 caracteres; no agregó EVT-0011. El comando posterior de verificación se había iniciado sin comprobar ese exit, por lo que no puede acreditarse como cierre.
- **Corrección prevista:** resumir el cursor, conservar el detalle en el reporte enlazado, comprobar longitud y código de salida antes de ejecutar el siguiente gate. No aumentar el límite para eludir el presupuesto de contexto.
- **Evidencia:** 815 caracteres rechazados; cursor resumido por debajo de 800 pasa validate_execution_state sin eventos como prevalidación semántica. La comprobación completa con cadena se exige al checkpoint posterior; el intento fallido no agregó ni eliminó eventos. V285.

### FAIL-20260907-395 — preflight conserva el conteo anterior de readiness

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** HIGH / CONTRACT; V285, VERIFY_EXECUTABLE_LIBRARY.ps1 AUTHORED.
- **Observado:** preflight real detuvo antes de inventariar toolchains con `project readiness regression count drifted: 70`; el pack canónico 0.6.1 contiene y ejecuta 70 tests pero el runner mantiene un conteo anterior.
- **Corrección prevista:** conciliar la aserción exacta con los 70 tests demostrados, sin quitar la comprobación de drift, y repetir preflight desde cero. Un VERIFY_LIBRARY_PASS estructural no detectó esta aserción del runner ejecutable.
- **Evidencia:** rerun completo llega al inventario, step project-readiness-syntax-and-70-regressions PASS / 2955 ms. Preflight final BLOCKED por cuatro herramientas no detectadas mediante resolución por defecto, no por este conteo. V285 conserva el fallo inicial y SHA del reporte final.

### FAIL-20260907-394 — compositor falla desde un caller con StrictMode

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** HIGH / CONTRACT; V285, T2801, MARKDOWN-COMPOSITOR 0.2.1 AUTHORED.
- **Observado:** composición real del perfil readiness desde un script con StrictMode Latest falla en línea 218: `The property 'Count' cannot be found on this object`; destino aún ausente, cero archivos incorporados.
- **Invariante:** el resultado de materializar un perfil válido no debe depender de la política de variables del script llamador; no desactivar StrictMode para ocultar defectos.
- **Corrección / evidencia:** selección escalar convertida explícitamente a string[] en 0.2.2; suite completa StrictMode con wildcard/empty/null/explicit PASS, 3/3 archivos idénticos al reconstruir, tres entradas reales y 70 tests readiness PASS. V285. Sin retirar StrictMode ni cambiar admisión.

### FAIL-20260907-393 — pipeline PowerShell sobre foreach statement

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** LOW / ENVIRONMENT; V285, comando diagnóstico propio.
- **Observado:** ParserError `An empty pipe element is not allowed` al colocar ConvertTo-Json detrás de un foreach statement; no se ejecutó la materialización anterior porque el parser rechazó todo el comando.
- **Corrección prevista:** capturar el resultado del foreach en variable y serializar esa variable; verificar destino ausente y repetir una vez. No es un defecto de packs ni código upstream.
- **Evidencia:** comando corregido produjo tres identidades observadas y materializó exactamente tres archivos. Una lectura posterior supuso un receipt del bootstrap ausente; se usó inventario real y comparación 3/3, sin fabricar receipt. URLs propuestas de OSV/StrictMode fallaron y se resolvieron desde fuentes oficiales. V285 conserva límites y errores diagnósticos; no son aprobaciones de producto.

### FAIL-20260907-390 — bridge atraviesa directorios redirigidos

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** HIGH / SECURITY; T2801, V284, owner INSTALL_AGENT_BRIDGE.ps1 AUTHORED.
- **Observado:** prueba real con cuatro junctions aisladas escribe fuera del proyecto y retorna éxito; el control sólo inspecciona el archivo final, no sus ancestros.
- **Invariante:** todos los destinos deben estar dentro del árbol autorizado sin enlaces intermedios; rechazo antes de cualquier escritura.
- **Corrección prevista:** validar ancestros antes de leer/planificar y nuevamente al escribir; regresión con destinos propios descartables, nunca datos del usuario. No prometer defensa ante un escritor hostil concurrente sin primitivas de handles/ACL.
- **Evidencia:** V284; cuatro junctions rechazadas antes de escritura, destinos externos intactos; suite 23 checks y copia limpia PASS.

### FAIL-20260907-391 — marcadores invertidos aceptados

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** MEDIUM / CONTRACT; T2801, V284, owner INSTALL_AGENT_BRIDGE.ps1 AUTHORED.
- **Observado:** END anterior a BEGIN cuenta como par válido, regex no reemplaza y el instalador afirma READY sin instalar el bloque.
- **Corrección prevista:** exigir orden además de cardinalidad y rechazar sin cambios; conservar hash original.
- **Evidencia:** V284; rechazo, hash intacto y cero skills nuevas; suite original y copia limpia PASS.

### FAIL-20260907-392 — rollback del bridge falla y deja cambios parciales

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** HIGH / RECOVERY; T2801, V284, owner INSTALL_AGENT_BRIDGE.ps1 AUTHORED.
- **Observado:** FileShare.Read permite planificar pero impide reemplazar el segundo archivo; rollback usa Select-Object -Reverse, parámetro inexistente, y no restaura el primer archivo. Reproducido con AGENTS existente y ausente.
- **Corrección prevista:** iteración inversa por índice y restauración de bytes exactos (incluido BOM), propagando fallo original. Suite inicial V284: 18 aserciones FAIL de 23; no se atribuyen a código upstream.
- **Evidencia:** V284; reemplazo fallido conserva error original y restaura SHA-256/BOM o ausencia, sin temporales residuales; 23 checks PASS y repetidos desde copia limpia. No prueba recuperación de corte eléctrico o filesystem hostil.

### FAIL-20260907-389 — readiness admite journeys CLI pero exige otra interfaz

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** HIGH / CONTRACT; preparación T2801, V283.
- **Owner:** PROJECT-START-READINESS-VALIDATOR 0.6.0, código AUTHORED.
- **Observado:** CONNECTED_JOURNEY_SURFACES y TOTAL_SYSTEM_CAPABILITY_CONTRACT aceptan CLI/AUTOMATION, pero delivery_surfaces exige WEB/MOBILE/DESKTOP/API. Un proyecto sólo CLI no puede cumplir sin declarar una interfaz que no tiene.
- **Invariante:** un gate no debe exigir arquitecturas ficticias; tampoco debe aceptar una etiqueta CLI que omita autorización, efectos, pruebas, ayuda o recovery.
- **Reproducción prevista:** fixture completo del gate con primer journey CLI o AUTOMATION, sin web ficticia; debe pasar sin eliminar ninguna evidencia. Negativos conservan rechazo por ausencia de auth/evidencia, selección ajena, interfaz mixta no clasificada o versiones divergentes.
- **Corrección prevista:** resolver superficies desde los journeys seleccionados, mantener las 48 capacidades y todos los controles de conexión, exigir autorización para la vía operativa explícita y consistencia de interfaces públicas con su capability.
- **Límite:** corregir este rechazo no completa el readiness local ni autoriza producto; no atribuir el defecto o corrección a una empresa externa.
- **Evidencia V283:** dos fallos reproducidos por exigencia de interfaz; 70 tests PASS tras corregir, repetidos desde reconstrucción canónica 0.6.1 con 8/8 archivos idénticos. Se prueban catorce refs ausentes, auth retirada, journey no seleccionado, web oculta, versiones divergentes y status no probado. El gate real de mantenimiento permanece BLOCKED; no se sustituyó por el fixture de tests.

### FAIL-20260907-388 — toolchain temporal incompleto para tests oficiales de slog

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** MEDIUM / ENVIRONMENT; investigación V281, T2809.
- **Comando:** Go 1.26.7, `go test -json -count=1 log/slog testing/slogtest`.
- **Fallo observado:** exit 1 antes de ejecutar tests; faltan `src/internal/testenv` y `src/testing/slogtest` en el toolchain temporal. También faltan VERSION y LICENSE. No demuestra un defecto de Go ni invalida por sí solo los builds anteriores de aplicaciones.
- **Invariante:** una versión declarada por el ejecutable no acredita integridad del árbol de fuentes ni habilita atribuir PASS a suites que no arrancaron.
- **Corrección prevista:** adquirir distribución oficial exacta en staging separado, contrastar SHA-256 publicado y repetir sin modificar el toolchain compartido ni omitir tests.
- **Evidencia:** staging V281, slog-official-tests.jsonl; conserva setup failed de ambos paquetes.
- **Corrección comprobada:** distribuciones oficiales completas 1.26.7 y 1.26.8 verificadas por SHA-256; 278 casos oficiales y 11 probes de alcance PASS en cada runtime, cero FAIL/SKIP. El toolchain compartido no fue modificado. Expediente OFFICIAL_STRUCTURED_LOGGING_AUDIT_V281.md.
- **Prevención:** antes de ejecutar suites internas, comprobar distribución íntegra y dependencias de prueba; no inferirlas desde go version. Este cierre afecta sólo el entorno del probe, no la admisión de observabilidad ni el upgrade de producto.
- **Tooling adicional:** nombre de ejemplo inferido inexistente corregido al inventariar `example_logvaluer_secret_test.go`; consultas web a URL versionada y JSON de descargas no pudieron abrirse. No atribuir al pin la documentación de otra versión. Lecturas agregadas truncadas no cuentan como revisión completa; continuar mediante fragmentos acotados.

### FAIL-20260907-385 — cobertura declarada no prueba composición

- **Estado:** `DIAGNOSED`
- **Severidad / clasificación:** HIGH / CONTRACT; auditoría V280.
- **Observado:** 21 packs de negocio/operación citados en el mapa no están seleccionados por FRANCHISE_COMPLETE_PACK_PLAN. Es brecha de trazabilidad, no prueba de 21 funciones ausentes: se debe contrastar cada una con owners alternativos seleccionados.
- **Invariante afectada:** existencia de un pack no demuestra integración, persistencia ni procedencia empresarial; un PASS histórico focal no cubre entradas que no probó.
- **Evidencia/acción:** CONNECTED_COVERAGE_AUDIT_V280.md conserva los 21 IDs y el roadmap distribuye exactamente las 48 superficies. Queda pendiente demostrar equivalencia o integración por core; no añadirlos al perfil para completar casillas. El contraste está completo; la conexión funcional no.
- **Procedencia observada:** 21/21 AUTHORED y 15/21 con upstream_sources limitado a go.dev. No es código empresarial publicado ni evidencia de calidad equivalente del dominio. La admisión y los claims de cobertura deben reconciliar ese hecho.
- **Recurrencias de tooling:** lectura agregada truncada; selección rg sin coincidencia terminó 1 sin modificar archivos. No se atribuye revisión a la porción truncada.
- **V293:** conciliación por función/owner: 12 solapamientos parciales, 8 sin equivalente demostrado y 1 candidato. FRANCHISE_GAP_MAP retira verdes/resumen global no sustentados y enlaza rutas/SQL/tests/guías actuales; perfil67 intacto. Sigue DIAGNOSED para equivalencia funcional/integración, no son 21 capacidades ausentes. Próxima discontinuidad: portal accept-quote→pedido durable en navegador. CORE_OWNER_RECONCILIATION_V293.md.
- **Patch:** se propuso cabecera abreviada de CAPABILITY_CATALOG; rechazo atómico, ausencia de cambios parciales comprobada y repetición con la cabecera observada.

### FAIL-20260907-386 — redacción parcial presentada como protección general de PII

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** HIGH / SECURITY.
- **Owner:** GO-OBSERVABILITY-CORE 0.1.0, código local AUTHORED; no defecto de Google ni de OpenTelemetry.
- **Reproducción:** email en campo superior, password en metadata anidada y email en Msg llegan intactos; tres pruebas adversariales fallan con datos sintéticos. Los cinco tests históricos pasan.
- **Causa:** Redact sólo visita el primer nivel y patrones concretos; Msg se copia directamente. No hubo observación de filtración productiva.
- **Contención:** admisión CANDIDATE y rechazo de composición probado, incluso con acknowledgeConditions=true; ningún perfil actual selecciona el pack. Código e histórico conservados.
- **Pendiente:** fuente oficial compatible o adaptación admitida, política explícita de datos permitidos, pruebas de superficie completa y montaje; T2809. No declarar este fallo corregido por bloquear admisión.
- **Evidencia:** reconstruction_evidence/CONNECTED_COVERAGE_AUDIT_V280.md; observability-audit.log.
- **Cierre local V292:** tres fallos baseline 0.1.1 reproducidos; 0.2.0 sustituye filtro heurístico por evento/valores cerrados, rechaza antes del sink y demuestra JSON útil mediante handler oficial Go. 16 tests, 11 semillas, fuzz de dos targets, vet/build, rebuild 2/2 y rechazo de composición PASS. Política/glue AUTHORED; configuración revisada e integración/operación pendientes, pack CANDIDATE. Evidencia: reconstruction_evidence/PRIVACY_POLICY_LOGGER_REPAIR_V292.md. No se declara PII universal ni se cierra T2809.

### FAIL-20260907-387 — contador anunciado monotónico acepta decrementos

- **Estado:** `REGRESSION_PROVEN`
- **Severidad / clasificación:** MEDIUM / CODE.
- **Owner:** GO-OBSERVABILITY-CORE 0.1.0.
- **Reproducción:** Inc(5), Inc(-1) devuelve 4; cuarta prueba adversarial falla. Contradice la invariante explícita de contador monotónico.
- **Causa:** Inc suma cualquier delta sin restricción; no se demostró contrato de overflow.
- **Contención:** misma reapertura CANDIDATE del pack, no cambio silencioso del código.
- **Pendiente:** semántica oficial de instrumentos, corrección/reemplazo admitido y regresiones de negativos/overflow; T2809.
- **Evidencia:** reconstruction_evidence/CONNECTED_COVERAGE_AUDIT_V280.md; observability-audit.log.
- **V291:** dos regresiones reproducidas sobre bytes 0.1.0: Inc(-1) reduce 5 a 4 e Inc(1) desde MaxInt64 devuelve MinInt64. Candidate agrega TryInc(value, accepted) con rechazo sin mutación/panic, conserva firma Inc y exige verificar el resultado cuando importe diagnosticar rechazos. Go spec/OpenTelemetry consultados; tests locales pasan. Pendiente retorno canónico/rebuild/fuzz para cierre, sin cerrar privacidad 386 ni admisión del pack.
- **Cierre V291:** 0.1.1 canónico, 2/2 archivos reconstruidos idénticos; 10 tests y siete semillas, vet/build, fuzz 10s candidate y rebuild PASS (1.791.391/1.774.273 ejecuciones). Composición CANDIDATE rechazada, privacidad 386 sigue fallando. No promoción, race/carga/operación aún pendientes. Evidencia: reconstruction_evidence/COUNTER_MONOTONIC_REPAIR_V291.md.

### FAIL-20260907-384 — continuidad de mantenimiento sin checkpoint y distribución incompatible

- **Estado:** `DIAGNOSED`
- **Severidad / clasificación:** HIGH / CONTRACT; auditoría de mantenimiento V279.
- **Observado:** no existen PROJECT_EXECUTION_STATE.json ni PROJECT_EXECUTION_EVENTS.jsonl en la biblioteca. AGENT_SYSTEM_START exige ambos al reanudar; VERIFY_LIBRARY.ps1 y CREATE_PORTABLE_ARCHIVE.ps1 rechazan cualquier archivo raíz fuera de su allowlist, incluidos esos dos nombres.
- **Invariante:** los tests del execution validator no prueban que el trabajo de mantenimiento tenga un checkpoint; el estado local no debe distribuirse como readiness de otro proyecto.
- **Alcance:** continuidad del agente y empaquetado; no demuestra pérdida, mezcla ni corrupción de los bloques V278.
- **Corrección prevista:** reconocer exclusivamente el par de archivos de estado local, validar su naturaleza regular y excluirlo del payload portable; conservar todos los demás rechazos. Crear un checkpoint DISCOVERY/BLOCKED con owners y hashes observados, sin fabricar readiness, Spec Kit ni implementación autorizada.
- **Regresión requerida:** par completo, ausencia válida en distribución, par incompleto, directorio, nombre casi igual y hash/historial alterados; reconstrucción del validator exacto y verificador general.
- **Pendiente separado:** el expediente de readiness/implementation_assurance del mantenimiento no está demostrado. Resolverlo antes de ampliar código de producto; no marcar este hallazgo cerrado por crear JSON.
- **V282 / diagnóstico adicional:** el par de checkpoint ya funciona, pero PROJECT_READINESS_GATE.json, registros PROJECT_* y los paths Spec Kit obligatorios siguen fuera de la allowlist. La corrección debe excluir sólo nombres locales exactos con checkpoint presente, rechazar directorios/reparse/paths inesperados y mantenerlos fuera del payload; no habilitar archivos arbitrarios ni declarar readiness por presencia.
- **Evidencia parcial V282:** exclusión exacta implementada en ambos selectores, 32 checks PASS y 62 regresiones del readiness validator materializado PASS. Primer expediente real produce BLOCKED/102 observaciones; no es fallo inesperado ni 102 defectos nuevos. MAINTENANCE_READINESS_BOUNDARY_V282.md conserva el alcance. Readiness/assurance siguen pendientes: el incidente no se cierra por esta reparación.
- **Recurrencia de tooling V282:** se consultó un nombre inferido PROJECT_READINESS_GATE.md; la ruta observada en el plan es PROJECT_START_READINESS_VALIDATOR.md. Repetición corregida mediante el owner literal del plan; no hubo edición sobre la ruta inexistente.
- **Recurrencia de tooling:** búsqueda con wildcard de path no expandido y nombre de documento inferido; lecturas agregadas truncadas. Usar rg -g sobre directorios literales y lecturas acotadas; no atribuir revisión a la salida omitida.
- **Patch:** un hunk de readiness usó una línea abreviada como contexto literal; apply_patch rechazó atómicamente la operación. Se confirmó que no había modificaciones parciales y se repitió con cabecera exacta.
- **Evidencia parcial V279:** incompatibilidad de distribución corregida con 14 regresiones de funciones reales; 24 tests del validator reconstruido PASS; checkpoint local y VERIFY_LIBRARY_PASS ejecutados. El estado sigue DISCOVERY/BLOCKED por los owners de readiness/assurance ausentes, no por la allowlist ya corregida. No promover el incidente completo a REGRESSION_PROVEN mientras ese pendiente siga abierto.

### FAIL-20260906-383 — estados accesibles ambiguos al conectar notificaciones con agenda

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** round-agenda.log muestra strict mode violation en tres navegadores; Agenda de turnos ahora contiene dos role=status sin nombre específico. Chromium desktop y el BFF conectado sí pasan; no ocultar los tres fallos.
- **Corrección prevista:** nombres accesibles distintos para resultado de agenda y consulta de notificación; seleccionar el estado por rol y nombre, no por índice ni eliminando el caso de respuesta perdida.
- **Regresión requerida:** cuatro navegadores de agenda y cuatro de notificación sobre el artefacto recompilado, con efectos SQL únicos y reconstrucción canónica.
- **Evidencia:** final-agenda.log conserva cuatro PASS con ambos paneles; verified-agenda/whatsapp.log repiten desde Markdown con configuración off/on respectivamente y efectos únicos. 102 web PASS y BFF conectado ejecutado; 22/22 hashes idénticos.

### FAIL-20260906-382 — PostgreSQL de pruebas iniciado en puerto predeterminado

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** V278 createdb a 55959 rechazado: pg_ctl start omitió opciones explícitas y el servidor quedó en loopback 5432. Log y postmaster.opts lo demuestran; ninguna migración ejecutada en ese intento.
- **Corrección:** detener exclusivamente el data directory de pruebas inspeccionado y reiniciar con -h 127.0.0.1 -p 55959; comprobar readiness antes de crear la base descartable.
- **Prevención:** conservar startup completo con puerto y host, no asumir que start reutiliza opciones de una ejecución anterior.
- **Evidencia:** ambas bases aisladas recibieron 52 migraciones; 43 Go/WhatsApp, BFF y ocho configuraciones de navegador verificadas. Capstones: 24 PASS con autoridad real; cuatro owners/22 rutas reconstruidas idénticas. Patch con dos operaciones del mismo archivo fue rechazado sin editar y se repitió agrupando hunks por archivo.
- **Recurrencia de materialización:** SourceRoot con separador final produjo rechazo seguro Source escaped root antes de copiar bloques; normalizar la ruta raíz sin separador final y repetir los cuatro owners. No se relajó la comprobación de contención.
- **Recurrencia de entorno de verificación:** los capstones del execution validator se ejecutaron aislados sin ELITE_AUTHORITY_ROOT y cinco casos rechazaron manuales ausentes. Fijar la raíz real mediante el contrato existente y repetir; no copiar manuales ficticios ni desactivar el control. Las búsquedas siguientes se limitan a rutas literales confirmadas para evitar confundir biblioteca y staging.
- **Recurrencia de quoting PowerShell:** helper Lighthouse pasó comillas con backslash a Node y el hijo terminó MODULE_NOT_FOUND antes del gate. Corregir sólo el ArgumentList del helper; conservar lighthouse-target.log y lighthouse-next.stderr.log originales. Ninguna aplicación ni política del gate se modifica.
- **Tooling recurrente:** inventario intentó web inexistente y glob de ruta Unix; usar src y directorios literales. Lecturas agregadas truncadas se sustituyen por secciones acotadas.
- **Comparación de artefactos:** el guard esperaba tres docs distintos y rechazó next-env.d.ts generado. Se comprobaron las dos importaciones exactas de Next y se mantuvo rechazo para cualquier otra diferencia. El helper tuvo errores de sobrecarga char/string y parseo JSON antes de comprobar errores PowerShell; se repitió con ErrorActionPreference Stop y comparación de líneas exactas. Los hunks deben seguir orden de archivo: actualizar primero cabecera y después tabla de hashes.
- **Recurrencia LIB-FAIL-1119:** pnpm ascendió al workspace padre; browser-first.log falló MODULE_NOT_FOUND del CLI. Instalar gate con --ignore-workspace --frozen-lockfile --offline. Se añade comprobación de CLI/build antes de crear fixtures o lanzar cuatro navegadores.

### FAIL-20260906-381 — código de tenant sintético aleatorio fuera del contrato

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** V277 worker-green.log conserva SQLSTATE 23514 en tenant_tenant_code_check; UUID sin prefijo puede comenzar con número. Recurrencia de LIB-FAIL-2101.
- **Corrección:** prefijo scope- conforme al CHECK real; no modificar constraints ni datos de negocio. Repetir suite y reconstrucción antes de cerrar.
- **Evidencia:** nueve tests worker ejecutados dos veces, 34 Go/WhatsApp reconstruidos sin SKIP, siete cola/processors, 42 Python, vet/build y seis hashes idénticos PASS. V277 conserva también los intentos fallidos.
- **Prevención:** fixtures sintéticos deben cumplir el contrato completo incluso para valores aleatorios. Revisar lecciones anteriores de la misma tabla.
- **Recurrencias de tooling V277:** lectura de paths inferidos y salida agregada truncada; propuesta de versión usó packVersion en vez del campo real version y fue rechazada sin editar. Después un hunk sin packId como contexto cambió Journey 0.10.0 en vez de WhatsApp (LIB-FAIL-2090); compositor rechazó identidad antes de crear destino. Se restauró Journey y fijó WhatsApp con contexto explícito; conservar composition.log y repetir en destino ausente. No confundir estos intentos fallidos con evidencia verde.

### FAIL-20260906-380 — fixture de permisos usa tipo de set incorrecto

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** router-tests.log conserva build failed: map[string]bool no puede asignarse a map[string]struct{}.
- **Causa:** fixture nuevo asumió la representación del Principal sin comprobar su tipo real.
- **Corrección/evidencia:** usar el set real de identity.Principal; no cambiar autorización productiva. Reconstrucción 6/6, 25 Go/PG sin SKIP, 42 Python y vet/build PASS; negativos de permisos/organización mantienen cero observaciones. Expediente V276.

### FAIL-20260906-379 — fixture SQL reutiliza parámetro UUID como text

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** red.log inicial contiene SQLSTATE 42P08; los cuatro casos no llegaron a probar la transición.
- **Causa:** mismo $2 para job_id UUID y queue text.
- **Corrección:** parámetros separados, sin alterar schema. Se conserva red.log y se exige una nueva reproducción con el fallo semántico esperado antes de corregir implementación.
- **Evidencia:** red-semantic.log alcanza los cuatro fallos del claim; suite corregida y reconstrucción pasan sin 42P08. Expediente V275.

### FAIL-20260906-378 — cierre de job sin vigencia ni generación del claim

- **Estado:** `REGRESSION_PROVEN`
- **Clasificación:** CODE / RECOVERY; mantenimiento V275.
- **Fallo observado:** Complete y Fail comparan sólo worker ID; no verifican vencimiento ni attempts. Reutilizar el mismo worker ID permite que una ejecución antigua cierre el claim nuevo. El processor tampoco revalida antes de ejecutar cada elemento del batch.
- **Invariante:** confirmar/reprogramar exclusivamente el intento vigente; no comenzar handlers con un claim ya perdido. Esto no convierte un efecto remoto en exactly-once.
- **Plan de regresión:** vencimiento sin takeover, takeover con worker ID repetido, generación incorrecta, callback posterior a pérdida y propagación del presupuesto al handler; PostgreSQL aislado, sin proveedores reales.
- **Corrección/evidencia:** owner GO_RELIABLE_ASYNC_WORKERS 0.2.0 exige generación, row lock antes del reloj y presupuesto por handler. Cuatro fallos semánticos antes del fix; siete tests/once subcasos nuevos sin SKIP, vet/build y cuatro hashes idénticos desde reconstrucción. Expediente V275. No extiende el claim a efectos remotos ni a outbox.


### FAIL-20260906-377 — guard de scope confunde suite HTTP con gates de navegador

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** la suite HTTP terminó exit 0, pero el guard de mantenimiento rechazó siete SKIP de pruebas de navegador/PostgreSQL no activadas (tres tests y cuatro variantes de agenda). Por ese guard, vet/build no llegaron a ejecutarse en el comando inicial.
- **Causa:** se seleccionó todo el package HTTP sin inventariar sus gates opt-in ajenos a provider integration.
- **Corrección y evidencia:** cuatro tests de dominio provider, HTTP focal y dos tests PostgreSQL del adapter con siete subcasos pasan sin SKIP desde reconstrucción; vet/build globales ejecutados y PASS. Se mantienen los siete SKIP iniciales en roundtrip-transport.log y no se atribuye nuevo PASS de navegador. Expediente V273.

### FAIL-20260906-376 — admisión durable omite estado y scope del replay

- **Estado:** `REGRESSION_PROVEN`
- **Clasificación:** SECURITY / CODE; mantenimiento V273.
- **Fallo reproducido:** TestProviderWebhookConnectionAndReplayScope falla en seis subcasos. Una conexión disabled inserta evento/job; replay disabled y otro proveedor/tipo/payload retornan éxito. Conexión inexistente sólo cae por FK. Resultado durable observado: dos eventos/jobs en lugar de uno.
- **Causa:** AcceptWebhook no reconsulta provider_connection y compara sólo body_sha256_hex cuando ON CONFLICT evita el INSERT.
- **Invariante:** tenant/conexión/proveedor activos en la transacción y replay del mismo tipo/contenido; ningún nuevo job tras revocación.
- **Corrección/evidencia:** bloqueo compartido de conexión activa, comparación exacta de provider/type/hash/payload JSONB; siete subcasos y test anterior pasan desde reconstrucción, incluidos desactivación concurrente y replay semánticamente igual. 2/2 hashes idénticos; dominio, HTTP focal, vet/build PASS. Pack 0.1.2 permanece CONDITIONED. Expediente V273.

### FAIL-20260906-375 — inventario preflight desactualizado tras evidencia V272

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** VERIFY_LIBRARY_FAILED: franchise preflight inventory is stale; expected 160/1413/705/51.
- **Causa:** se añadió un Markdown de evidencia y se actualizó readiness, pero no el snapshot del preflight.
- **Corrección:** sincronizar sólo el inventario vigente a 705 Markdown, mantener contadores históricos y la aserción original del verificador.
- **Regresión focal:** normalizar whitespace y comprobar la cadena exacta exigida por VERIFY_LIBRARY: PREFLIGHT_SNAPSHOT_REGRESSION_PASS. El control general se repite por separado y sólo su salida determina su PASS.

### FAIL-20260906-374 — blancos de patch interpretados como contexto

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** el patch generado de cuatro bloques añadió líneas vacías sin prefijo, por lo que apply_patch buscó contexto inexistente y rechazó la edición canónica.
- **Corrección:** prefijar cada línea agregada y separar metadata/manifiesto del ancla literal única de bloques; metadata 0.7.0 intacta comprobada tras el rechazo.
- **Evidencia:** reconstrucción 67/722, diez hashes iguales, 42 Python, package Go/PG completo, vet/build y 0051 down/up. No se relaja compositor ni se conserva un placeholder.

### FAIL-20260906-373 — recurrencia de ubicación y frontera de pruebas

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** se buscaron owners materializados desde la raíz documental y rutas backend/cli.py no observadas; hubo salidas agregadas truncadas. Un enlace Postman secundario no pudo abrirse; la página oficial principal sí mostró el contrato necesario. En revisión del test se detectó un patch sobre urllib.request.urlopen en lugar del alias whatsapp_cloud.urlopen realmente usado.
- **Corrección:** usar inventario literal del árbol roundtrip, módulo observado whatsapp_cloud.py, lectura acotada y la fuente oficial principal; corregir el patch sobre el alias real antes del rebuild. No afirmar que se revisaron partes truncadas ni que el enlace fallido fue verificado.
- **Prevención:** separar lectura de biblioteca y materializado; comprobar el símbolo importado antes de interceptar una llamada en pruebas.
- **Recurrencia V273:** el enlace developers.facebook.com de componentes no fue accesible mediante búsqueda. Se verificó la página oficial Meta/Postman de statuses; no se usaron agregadores ni se afirmó una política actual de retries no demostrada.
- **Recurrencia V274:** el overview actual de business-messaging devolvió 429 y su variante .md no se pudo abrir. Meta/Postman permaneció accesible; la documentación oficial del SDK Node advierte archivo y sólo se conserva como contexto histórico, nunca como runtime admitido ni prueba actual de retries. El contrato live queda pendiente. Una lectura agregada local también fue truncada; se releyeron las definiciones de fixtures necesarias sin presumir revisión de la parte omitida.
- **Recurrencia V275:** lecturas agregadas truncadas y comandos mezclando biblioteca/árbol materializado; rg recibió un wildcard de path no expandido en Windows y se supuso migrations en lugar del db/migrations inventariado. Se corrigieron cwd y paths literales; no se atribuye revisión a contenido omitido. La URL AWS inglesa redirige a una página sin texto en el extractor: se localizó la edición oficial española y el PDF oficial, sin sustituirlos por agregadores.
- **Recurrencia V276:** un diagnóstico agregado excedió por pocas líneas el presupuesto de salida; no se atribuye revisión a la porción truncada. Se leyeron los schemas concretos en llamadas acotadas antes de preparar fixtures SQL.
- **Recurrencia V277:** se intentó abrir principal.go sin inventario previo; el tipo vive en oidc.go. rg --files confirmó la ruta; no se creó un archivo alternativo. Los fixtures nuevos conservaron el set real map[string]struct{} y sus pruebas pasaron sin modificar autorización. La invocación inicial Go sólo compiló whatsappbridge (no tests to run); sus tests de worker se ejecutaron después explícitamente, sin presentar la compilación como PASS de esos escenarios.
- **Evidencia:** fuente Meta principal verificada, patch sobre whatsapp_cloud.urlopen, 41 tests Python repetidos desde reconstrucción, tres tests focales Go de proceso/contrato, vet/build y 5/5 hashes idénticos. No se presenta revisión de contenido truncado ni del enlace secundario.

### FAIL-20260906-372 — propuesta mecánica de conteos rechazada antes de editar

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** arrays anidados de un único par se aplanaron en PowerShell; la propuesta de reemplazos produjo demasiada salida y JSON.parse rechazó el resultado truncado antes de aplicar cambios.
- **Corrección:** diccionarios old/new explícitos, comparación de líneas completas y aserción de exactamente siete líneas esperadas. Se aplicaron sólo esos siete cambios; ningún reemplazo por carácter se escribió.
- **Prevención:** no usar arrays anidados susceptibles de aplanarse para pares de sustitución; validar cantidad y forma antes de generar apply_patch.
- **Recurrencia V273:** una línea JSON serverless excedió el presupuesto y el guard rechazó la propuesta antes de editarla. Se capturó completa sin volcarla y se sustituyó únicamente el par packId/version exacto, con coincidencia única. Un patch posterior de lecciones fue rechazado por hunks fuera de orden; se releyeron sus anclas y se aplicaron en orden, sin cambios parciales ni debilitamiento de checks.

### FAIL-20260906-371 — límite BFF posterior a descarga y medido en caracteres

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** los clientes público y protegido leen response.text() completo antes de aplicar el límite; body.length no mide bytes UTF-8. Un payload multibyte supera el presupuesto real sin rechazo y la respuesta de media type inválido no se cancela.
- **Corrección:** un solo lector incremental compartido, límite por bytes antes de decodificar, cancelación y errores sin payload; preservar Problem Details, deadline, autorización y ausencia de retry.
- **Regresión:** red 11 FAIL/16 PASS/1 SKIP; green y reconstrucción limpia 70 PASS/1 SKIP, typecheck/build. Tres pruebas HTTP nativas demuestran gzip, cancelación chunked y timeout del cuerpo; 4/4 archivos cambiados y lock idénticos. Expediente V269. El SKIP conectado PostgreSQL no se presenta como PASS.

### FAIL-20260906-370 — ancla parcial de README

- **Estado:** `RECURRENCE_PROVEN`
- **Fallo:** patch buscó una frase que empieza a mitad de la línea real del README.
- **Corrección:** comprobar ausencia de cambios parciales y usar la línea completa observada; no duplicar hunks ni relajar validadores.

### FAIL-20260906-369 — fixture de reasignación con lead nulo

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** TestNotificationStatusScopeAndUnavailable/reassigned viola NOT NULL de crm.appointment.lead_id antes de probar aislamiento.
- **Causa:** el fixture confundió reasignación con ausencia de lead; la inspección del schema y ejecución se agruparon sin aplicar la evidencia antes del test.
- **Corrección:** crear un segundo lead sintético válido y reasignar a su ID; conservar NOT NULL/FK y repetir desde fuente canónica.
- **Regresión:** suite focal staging y full Go desde reconstrucción limpia PASS, incluidos reasignación, SQL read-only, POST→GET y negativos; 5/5 hashes iguales. No se deshabilita ni modifica ninguna constraint.

### FAIL-20260906-368 — recurrencia de scope de lectura

- **Estado:** `RECURRENCE_PROVEN`
- **Fallo:** filtro sobre paths absolutos coincidió con AppData y truncó inventario; se supuso channel.go para outbound y se consultó un manual de biblioteca desde el árbol materializado.
- **Corrección:** inventario relativo, owner observado delivery.go y raíz explícita para manuales; no se crea ningún owner alternativo. Lecturas siguientes acotadas a contratos observados.
- **Prevención:** no mezclar descubrimiento y apertura de una ruta todavía no confirmada en el mismo comando.
- **Recurrencia V268:** lectura agregada de routers excedió el presupuesto de salida. Las siguientes lecturas se separaron por owner y no se tomó la parte truncada como evidencia de revisión completa.
- **Recurrencia V269:** lectura agregada y resultados de planes compactados excedieron salida. No se presume revisión integral del contenido truncado; consultar sólo bloques materiales delimitados.
- **Recurrencia V272:** se usó un wildcard como ruta literal de rg y se intentó abrir un test materializado desde la raíz documental; ambos comandos fallaron sin mutaciones. También hubo lecturas agregadas truncadas. Corrección: inventario relativo con rg --files en work, rutas observadas y workdir explícito; no usar la salida omitida como evidencia. El test real fue localizado en internal/whatsappbridge/appointment_approval_test.go antes de editar.

### FAIL-20260906-367 — UUID usado como tenant_code de fixture

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** la suite de aprobaciones encuentra tenant_tenant_code_check cuando el UUID sintético comienza por un dígito.
- **Causa:** el fixture no respeta el código legible existente, aunque tenant_id sí sea UUID.
- **Corrección:** prefijar el código de fixture con wa-, mantener UUID como identidad y conservar la constraint; repetir suite completa y reconstrucción.
- **Regresión:** suite focal staging y full Go desde 711 archivos reconstruidos PASS con PostgreSQL configurado; 19 negativos, seis invalidaciones y ocho aprobaciones concurrentes. Migraciones 0001–0050 y down/up 0050 en base vacía separada PASS. Siete archivos del delta idénticos por hash; expediente V266.

### FAIL-20260906-366 — path de identidad inferido

- **Estado:** `RECURRENCE_PROVEN`
- **Fallo:** se intentó internal/identity/identity.go, que no existe.
- **Corrección:** inventario rg --files localizó internal/platform/identity/oidc.go; leído antes de elegir autorización. No se creó un segundo owner.
- **Recurrencia en el mismo ciclo:** nombre supuesto de migración 0001 no existe; búsqueda por contenido en db/migrations resuelve el outbox sin inventar filenames. Se mantiene la huella de inventario previo obligatorio.

### FAIL-20260906-365 — notices omite dos archivos AUTHORED nuevos

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** VERIFY_LIBRARY rechaza Current-Provenance-Counts antiguo; espera AUTHORED=1159, ADAPTED=133, VERBATIM=105, TOTAL=1397.
- **Causa:** se sincronizaron composición/hashes pero no el resumen de procedencia del distribuidor.
- **Corrección:** actualizar el owner THIRD_PARTY_NOTICES sin cambiar atribuciones individuales ni licencias. Repetir gate raíz; no omitir el control.
- **Evidencia:** VERIFY_LIBRARY_PASS exit 0, 160 packs/1397 archivos/698 Markdown, después de sincronizar notices; control intacto.

### FAIL-20260906-364 — protocolo stdio cambia LF por CRLF

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** cinco fixtures Go y una prueba Python reciben CRLF en Windows cuando el contrato byte-exacto espera LF.
- **Causa:** escritura por streams de texto con traducción de newline.
- **Corrección prevista:** protocolo y marcador de prueba por bytes binarios, sin debilitar las aserciones de una sola invocación ni redacción.
- **Evidencia:** 22 pruebas Python y bridge Go/PG pasan reconstruidos, 6/6 hashes idénticos; `reconstruction_evidence/WHATSAPP_GO_DURABLE_BRIDGE_V265.md` conserva el primer fallo.

### FAIL-20260906-363 — recurrencia glob como path de Windows

- **Estado:** `RECURRENCE_PROVEN`
- **Fallo:** `rg implementation_packs/GO*` devuelve os error 123.
- **Corrección:** repetir sobre `implementation_packs --glob 'GO*.md'`, sin expansión de shell; búsqueda corregida ejecutada antes de elegir implementación.
- **Prevención:** usar literalmente el directorio inventariado y el filtro nativo. Relacionado: FAIL-20260905-336.

### FAIL-20260906-362 — WhatsApp firmado sin scope de cuenta/teléfono

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** el normalizador acepta HMAC válido pero descarta WABA, phone_number_id y recipient_id del status; permite eventos ajenos al perfil y estados/identidades vacíos.
- **Impacto:** la firma de una app no demuestra el binding al tenant ni autoriza marcar una entrega de otro número.
- **Corrección prevista:** exigir perfil, WABA/teléfono exactos, estructuras e identidades acotadas; conservar scope y destinatario pseudonimizados. No crear otro inbox ni inferir delivered desde accepted.
- **Regresión:** antes del fix 1 FAIL/8 PASS (PermissionError no levantado); después 17/17 PASS en staging y reconstrucción limpia, cinco archivos idénticos también en composición 67/704. Fuente oficial preservada. Evidencia: `reconstruction_evidence/WHATSAPP_SCOPED_WEBHOOK_V264.md`.

**Cierre V263:** FAIL-20260906-357–361 se corrigen y verifican en `reconstruction_evidence/OPERATOR_AGENDA_BROWSER_POSTGRES_V263.md`: cuatro navegadores sobre TLS real loopback, recuperación post-commit, origen configurado, 54 tests web/build/typecheck, Go/PG y round-trip 12/12. Los rechazos previos se conservan como historia.

### FAIL-20260906-361 — hunks con contexto superpuesto en plan

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** el plan web rechazó cambios con contexto anterior/posterior que se superponía entre versiones cercanas. Perfil integral sí quedó actualizado; no se asume actualización de los restantes.
- **Corrección:** un único hunk desde el contenido literal completo de cada plan corto o minificado; preservar todos los campos no seleccionados y reconstruir.


### FAIL-20260906-360 — encabezado README supuesto

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** patch del README materializado rechazado por asumir su título; la operación no aplicó el hunk de overflow incluido.
- **Corrección:** leer título literal, separar archivos y volver a verificar diff y tests; no inferir aplicación parcial.


### FAIL-20260906-359 — sesión segura de fixture WebKit

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** Chromium desktop/mobile y Firefox pasan; WebKit redirige a login y exige OIDC_ISSUER. El fixture debía aportar sesión cifrada, no iniciar login externo.
- **Acción:** inspeccionar presencia de cookie sin imprimir su contenido y transporte loopback. No desactivar Secure ni simular una identidad aceptada por Go.
- **Resolución:** el shim de routing no transportó la cookie segura en WebKit. Cambiar sólo forma de addCookies y headers no lo resolvió. Proxy TLS local real + certificado efímero + excepción de CA únicamente loopback permitieron 4/4; Secure/HttpOnly y verificación RSA siguen activas. No se certifica PKI productiva.


### FAIL-20260906-358 — asignación en navegador no confirmada

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** Chromium abrió la agenda autenticada y eligió recurso, pero no apareció el recurso asignado tras el POST. El gate paró al primer fallo; no equivale a PASS ni a confirmación.
- **Causa observada:** respuesta CROSS_ORIGIN_REJECTED; Origin/Host públicos correctos, pero la comparación usaba request.nextUrl interno. Dos regresiones adicionales reproducen rechazo detrás de proxy y falta de dependencia del origen configurado (2 FAIL/12 PASS).
- **Corrección:** reutilizar applicationBaseUrl/APP_BASE_URL, comparación exacta y fail-closed si falta. No confiar en Host/X-Forwarded-Host. El fixture mantiene cookie Secure y origen HTTPS mediante shim loopback declarado, no certifica TLS real.


### FAIL-20260906-357 — diagnóstico y patch fuera de orden V263

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** lectura agregada truncada; URL inicial Microsoft no disponible (la búsqueda encontró work-with-schedule-board); patch HTTP pidió volver a un encabezado anterior tras modificar rutas y fue rechazado.
- **Corrección:** leer por archivo/rango y fuente oficial resuelta; inspeccionar postcondiciones y aplicar hunks de arriba hacia abajo. No dar por aplicada la transacción rechazada.


### FAIL-20260906-356 — composición detenida por identidad/version de ARCA

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** al reconstruir V262, compositor rechaza MICROSOFT_ARCA_WSFE_GENERATED_CLIENT por diferencia con el perfil. No se atribuye al cambio de confirmación sin revisar el owner.
- **Causa confirmada:** el patch generado incluyó sólo la línea de versión, repetida en el JSON pretty; modificó la primera 0.2.1 (ARCA), no OIDC. Fue un error de edición de esta tarea, no cambio ajeno.
- **Corrección/evidencia:** restaurar ARCA 0.2.1 y actualizar OIDC 0.2.2 con contexto packId. Composición nueva roundtrip-r1 67/704 PASS; nunca debilitar el gate ni usar un hunk de versión sin identidad.


### FAIL-20260906-355 — contrato de errores del cliente protegido

- **Estado:** `REGRESSION_PROVEN`
- **Fallo reproducido:** 8 FAIL/6 PASS: GET y POST convierten 401/403/409 Problem Details en 502; aceptan application/jsonp en respuesta exitosa.
- **Corrección/evidencia:** media type normalizado exacto y problem+json sólo en error: 14/14 unitarios focales, 52 web, typecheck/build y gate BFF/Go/PG con 409 real concurrente PASS. Fuente devuelta al pack TS-OIDC-PORTAL-ADAPTER 0.2.2.


### FAIL-20260906-354 — recurrencia de rutas y salida de diagnóstico

- **Estado:** `RECURRENCE_PROVEN`
- **Fallo:** nombres de tests inexistentes, ledger buscado desde materialización, globs como rutas rg, lectura agregada truncada y serialización de session_id ausente tras un proceso ya finalizado.
- **Corrección:** inventario literal, raíz explícita, directorio con --glob, lecturas acotadas y no persistir undefined. El compositor terminó 67/704; no se repite su ejecución ni se alteraron datos de producto.

**Cierre V261 de FAIL-20260906-348…351:** correcciones devueltas a cuatro packs canónicos; ocho proyectos conectados, PostgreSQL, 40 tests web, typecheck/build, Go test/vet/build y round-trip 6/6 PASS. Evidencia: `reconstruction_evidence/PUBLIC_APPOINTMENT_RECOVERY_V261.md`. Las hipótesis iniciales se conservan como historia.

### FAIL-20260906-353 — orden de hunks en registro

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** patch de cierre rechazado al pedir al editor volver al encabezado después de hunks posteriores.
- **Corrección:** ordenar los cambios de arriba hacia abajo y separar archivos; el rechazo no cambió código ni evidencia ejecutada.

### FAIL-20260906-352 — reconstrucción durante diagnóstico activo

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** el primer build corregido se ejecutó antes de finalizar toda la matriz roja. Esa ejecución no identifica un artefacto inmutable y no se usa como evidencia de cierre.
- **Corrección/evidencia:** esperar a que el proceso finalice, construir y repetir sin escrituras; round-trip 6/6, ocho proyectos conectados y suites PostgreSQL PASS posteriores.

### FAIL-20260906-351 — interacción antes de hidratación

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** con Chromium/Firefox verdes, WebKit conserva status vacío y vuelve a placeholder tras selectOption; no hay POST de turno.
- **Hipótesis/corrección:** hidratación reinicia la selección SSR antes de instalar handlers; deshabilitar controles hasta el efecto de montaje, como documenta Microsoft Playwright. No usar sleeps ni omitir WebKit.

### FAIL-20260906-350 — replay de turno sin señal HTTP

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** Firefox primero y Chromium después de corregir el consumo del receipt reciben 200 sin Idempotency-Replayed. Esperado true, obtenido undefined.
- **Corrección:** la API Go debe emitir el header sólo cuando el repositorio devuelve replayed=true; el BFF ya sabe propagarlo. Mantener prueba conectada y agregar recuperación ante respuesta perdida después del commit.

### FAIL-20260906-349 — interfaz de turno sin receipt validado

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** la UI mostró Turno solicitado sin leer el body; la prueba de navegador esperó response.json hasta 30 s. La causa exacta del timeout de transporte no se presume demostrada; el defecto de validación sí se observa en el código.
- **Corrección:** consumir y validar id/state de la respuesta antes de anunciar éxito, mostrar referencia y conservar clave ante fallo; repetir browser/Go/PG.

### FAIL-20260906-348 — fixture sin disponibilidad de organización

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** TestPublicAppointmentBrowserPostgres rechazado con SQLSTATE 23514: appointment slot is outside organization availability.
- **Causa:** la fixture agregó slots sin la jornada exigida por el contrato existente.
- **Corrección:** preparar disponibilidad sintética acotada antes de los slots; conservar el trigger y repetir el gate. No es fallo del control productivo.

## 1. Control

```yaml
ledger_version: "1.0"
project_id: "elite-engineering-library"
owner: "workspace-owner"
created_at: "2026-09-01"
updated_at: "2026-09-04"
open_high_or_critical: ["FAIL-20260904-009", "FAIL-20260904-024"]
recurring_failures: []
last_clean_rebuild_evidence: "V171 canonical: supply 78/78, backend 29/370, PostgreSQL 0001-0036 plus down/up, full Go/vet/build; VERIFY_LIBRARY_PASS packs=94 materialized_files=1087 markdown_files=528 and VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=94 upstream_sources=121 on 2026-09-01"
```

- failure_id: FAIL-20260901-007
  fingerprint: "v173-connected-shipment-discovery-truncation"
  related_failures: ["LIB-FAIL-1705", "LIB-FAIL-1706", "LIB-FAIL-1707", "LIB-FAIL-1708", "LIB-FAIL-1709", "LIB-FAIL-1710", "LIB-FAIL-1711", "LIB-FAIL-1712", "LIB-FAIL-1713", "LIB-FAIL-1714", "LIB-FAIL-1715", "LIB-FAIL-1716", "LIB-FAIL-1717", "LIB-FAIL-1718", "LIB-FAIL-1719", "LIB-FAIL-1720", "LIB-FAIL-1721", "LIB-FAIL-1722", "LIB-FAIL-1723", "LIB-FAIL-1724"]
  first_seen: "2026-09-01"
  last_seen: "2026-09-01"
  recurrence_count: 5
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: DATA_INTEGRITY
  phase: implementation
  capability: "CONNECTED-CUSTOMER-SHIPMENT"
  scope: "Registered customer-order pick to durable shipment and order progress; carrier delivery and fiscal posting remain separate owners"
  command_redacted: "cross-owner Go and SQL search for shipment, delivery and inventory issue"
  toolchain: ["PowerShell", "ripgrep"]
  error_signature_redacted: "cross-owner output truncated and multiple unobserved conventional paths/globs were rejected"
  violated_invariant: "Discovery must be complete and any physical shipment must conserve physical, UOM composition, reservation and cost state atomically."
  blast_radius: "No production project is certified, but the prior transfer-shipment claim did not prove UOM composition consumption after V171."
  reproduction: "Search every Go and SQL file for shipment, fulfillment, delivery and inventory terms in one command."
  hypotheses_rejected: ["Visible matches form a complete ownership map", "The existing transfer shipment automatically inherited packaging conservation"]
  root_cause: "Owner boundaries/ranges were not separated early and unobserved conventional paths/globs were added."
  canonical_correction: "Enumerate paths, use ripgrep file filters, locate exact symbols per owner and inspect only bounded ranges."
  regression_test: "Transfer/customer shipment conservation, exact/divergent replay, one-winner concurrency, 0038 down/up, PostgreSQL clean, full Go/vet/build and canonical composition."
  clean_rebuild_evidence: "MICROSOFT_BC_CUSTOMER_SALES_SHIPMENT_2026-09-01_V173.md; supply 88/88, backend 29/380, PostgreSQL 0001-0038 plus 24 SQL tests, full Go/vet/build and 0038 round-trip focal PASS."
  dependent_gates_rerun: ["owner-scoped discovery", "supply round-trip", "enterprise composition", "PostgreSQL", "Go", "root", "Audit"]
  upstream_reference: "microsoft/BCApps@2eae56d704a1fd035d104f333602aea7091b7749 exact SalesPost/SalesShipmentHeader/SalesShipmentLine/SalesWhsePostShipment/ERMSalesOrder files."
  residual_risk: "Carrier, delivery/handover, fiscal/payment, serial/specific customer shipment and project production gates remain conditioned."
  owner: "workspace-owner"
  next_action: "Connect the next authoritative owner without conflating carrier, delivery or fiscal completion."

## 2. Entradas

```yaml
- failure_id: FAIL-20260901-001
  fingerprint: "preflight-playwright-offline-tarball-missing"
  related_failures: ["LIB-FAIL-1119", "LIB-FAIL-1120"]
  first_seen: "2026-09-01"
  last_seen: "2026-09-01"
  recurrence_count: 2
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: ENVIRONMENT
  phase: preflight
  capability: "TEST-PLATFORM"
  scope: "Microsoft Playwright browser gate in VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Preflight"
  command_redacted: ".\\VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Preflight"
  toolchain: ["pnpm 11.19.0", "@playwright/test 1.62.1"]
  error_signature_redacted: "ERR_PNPM_NO_OFFLINE_TARBALL for @playwright/test 1.62.1"
  violated_invariant: "A frozen offline install requires every exact locked tarball to exist in the local pnpm store."
  blast_radius: "The executable Preflight cannot finish on this host; structural verification and the already recorded historical browser evidence remain intact."
  reproduction: "Materialize the Microsoft Playwright browser gate and run its frozen pnpm install in offline mode on the current host."
  hypotheses_rejected: ["Library pack metadata or manifest corruption; VERIFY_LIBRARY.ps1 passed independently."]
  root_cause: "The exact locked npm tarball is absent from this host's pnpm store and network acquisition is not authorized in this run."
  canonical_correction: "No source correction is justified. Acquire and verify the exact locked tarball through an authorized dependency lane, then rerun Preflight from a clean destination."
  regression_test: "Frozen offline install of the materialized Playwright gate followed by its runtime contracts."
  clean_rebuild_evidence: "The exact lock installed frozen/offline, four browser runtime tests passed, and VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit exited 0 on 2026-09-01."
  dependent_gates_rerun: ["VERIFY_LIBRARY.ps1: PASS", "Microsoft Playwright browser runtime: 4 tests PASS", "VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Audit: PASS"]
  upstream_reference: "reconstruction_evidence/MICROSOFT_PLAYWRIGHT_BROWSER_GATE_2026-08-29_V103.md"
  residual_risk: "The current host reproduces the offline runtime gate; other hosts still require the exact locked cache inputs and browser binaries."
  owner: "workspace-owner"
  next_action: "Retain the exact lock/integrity and require each new host to pass the frozen install and browser runtime before relying on it."

- failure_id: FAIL-20260901-002
  fingerprint: "v169-receipt-putaway-packaging-maintenance-chain"
  related_failures: ["LIB-FAIL-1625..LIB-FAIL-1656"]
  first_seen: "2026-09-01"
  last_seen: "2026-09-01"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: IMPLEMENTATION
  phase: maintenance
  capability: "SUPPLY-WAREHOUSE-PACKAGING"
  scope: "GO-SUPPLY-FACTORY-INVENTORY-API 0.13.0 receipt-to-put-away packaging preservation, authorized breakbulk, reservation, registration and cancellation"
  command_redacted: "materialize supply and backend; migrate PostgreSQL 0001..0035; run focal, full, vet, build, root and Audit gates"
  toolchain: ["Go 1.26.7", "PostgreSQL 18.6", "pgx 5.10.0"]
  error_signature_redacted: "Detailed unique signatures are retained in LIB-FAIL-1625 through LIB-FAIL-1656."
  violated_invariant: "No V169 claim may be promoted until exact quantity/composition conservation, atomic rollback, migration compatibility, reconstruction equality and all dependent gates are observed."
  blast_radius: "Each failed attempt was invalidated; the canonical library was not promoted from partial output."
  reproduction: "Use the detailed ledger command/context for each ID, or rerun the final clean reconstruction and PostgreSQL corpus from the V169 evidence."
  hypotheses_rejected: ["A structural PASS proves runtime behavior", "Base quantity alone proves packaging conservation", "An omitted or truncated gate may be treated as PASS"]
  root_cause: "The chain exposed fixture, SQL scope, transfer reservation, pack fence, provenance/pin and verification-output defects while integrating the new invariant."
  canonical_correction: "The canonical pack now reserves physical and composition quantities, preserves or explicitly breaks packages, fixes transfer-receipt parity, reconstructs 75/75, composes backend 29/367 and keeps unsupported pick/replenishment packaging fail-closed."
  regression_test: "TestWarehouseReceiptPackagingPreservationBreakbulkAndCancellation plus combined warehouse/UOM/transfer, legacy 0034↔0035, full Go/vet/build, reconstruction diff, root verify and executable Audit."
  clean_rebuild_evidence: "reconstruction_evidence/MICROSOFT_BC_RECEIPT_PUTAWAY_PACKAGING_2026-09-01_V169.md"
  dependent_gates_rerun: ["supply 75/75 path+SHA equality", "backend 29/367 Go/PostgreSQL/vet/build PASS", "VERIFY_LIBRARY.ps1: PASS", "VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Audit: PASS"]
  upstream_reference: "microsoft/BCApps commit 2eae56d704a1fd035d104f333602aea7091b7749 and the fixed Microsoft Learn authorities in the V169 evidence"
  residual_risk: "Packaging-aware picking/replenishment, external WMS, target load/security/recovery/deploy and business acceptance remain unproved."
  owner: "workspace-owner"
  next_action: "Implement and prove the packaging-aware pick/replenishment lane without broadening the V169 claim."

- failure_id: FAIL-20260901-003
  fingerprint: "readiness-review-query-shape-and-powershell-metadata-count"
  related_failures: ["LIB-FAIL-1657", "LIB-FAIL-1658", "LIB-FAIL-1659"]
  first_seen: "2026-09-01"
  last_seen: "2026-09-01"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  phase: review
  capability: "DOCS-OPS"
  scope: "Read-only measurement of current library readiness"
  command_redacted: "bounded reads and metadata classification over implementation_packs"
  toolchain: ["PowerShell"]
  error_signature_redacted: "output truncation; Get-Content Raw and TotalCount cannot be combined"
  violated_invariant: "A readiness percentage may use only complete, observed measurements."
  blast_radius: "No library implementation or pack was modified; invalid and omitted output was excluded."
  reproduction: "Run the rejected broad reads or combine Get-Content -Raw with -TotalCount."
  hypotheses_rejected: ["Partial historical search output is sufficient", "Raw and TotalCount can be combined"]
  root_cause: "The review queries were broader than required and one PowerShell read shape was invalid."
  canonical_correction: "Use complete table ranges and join bounded line arrays before regex classification."
  regression_test: "Corrected counter exits 0 and accounts for exactly 94 pack metadata records."
  clean_rebuild_evidence: "Corrected result: REBUILD_VERIFIED=94; CONDITIONED=92; REUSABLE_PACK=1; CANDIDATE=1."
  dependent_gates_rerun: ["VERIFY_LIBRARY.ps1"]
  upstream_reference: "none"
  residual_risk: "The resulting percentage remains an explicitly labeled estimate, not a formal product acceptance score."
  owner: "workspace-owner"
  next_action: "Use the 48-capability project record for formal project-specific scoring."

- failure_id: FAIL-20260901-004
  fingerprint: "v170-ripgrep-literal-markdown-glob"
  related_failures: ["LIB-FAIL-1660", "LIB-FAIL-1661", "LIB-FAIL-1662", "LIB-FAIL-1663", "LIB-FAIL-1664", "LIB-FAIL-1665", "LIB-FAIL-1666", "LIB-FAIL-1667"]
  first_seen: "2026-09-01"
  last_seen: "2026-09-01"
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  phase: maintenance
  capability: "DOCS-OPS"
  scope: "Discovery of the existing readiness materializer"
  command_redacted: "rg over a literal PowerShell wildcard path"
  toolchain: ["PowerShell", "ripgrep"]
  error_signature_redacted: "Windows error 123 for markdown_system/*.md"
  violated_invariant: "File discovery must not assume shell wildcard expansion."
  blast_radius: "Read-only discovery failed; no canonical artifact changed."
  reproduction: "Pass markdown_system/*.md as a positional path to rg on this shell."
  hypotheses_rejected: ["The materializer is absent"]
  root_cause: "The wildcard was treated as a literal path."
  canonical_correction: "Use rg --files and filter the observed paths."
  regression_test: "The corrected query located materialize_markdown_pack.ps1."
  clean_rebuild_evidence: "PROJECT-START-READINESS-VALIDATOR 0.6.0 rebuilt 8/8 with zero path/SHA differences; 62/62 tests and the complete executable Audit passed."
  dependent_gates_rerun: ["PROJECT-START-READINESS-VALIDATOR tests", "VERIFY_LIBRARY.ps1"]
  upstream_reference: "none"
  residual_risk: "The gate validates evidence structure and connected coverage; each project must still prove semantic truth and every target/live production condition."
  owner: "workspace-owner"
  next_action: "Use the readiness validator in each target project."

- failure_id: FAIL-20260901-005
  fingerprint: "v171-packaging-aware-pick-replenishment-maintenance-chain"
  related_failures: ["LIB-FAIL-1668", "LIB-FAIL-1669", "LIB-FAIL-1670", "LIB-FAIL-1671", "LIB-FAIL-1672", "LIB-FAIL-1673", "LIB-FAIL-1674", "LIB-FAIL-1675", "LIB-FAIL-1676", "LIB-FAIL-1677", "LIB-FAIL-1678", "LIB-FAIL-1679", "LIB-FAIL-1680", "LIB-FAIL-1681", "LIB-FAIL-1682", "LIB-FAIL-1683", "LIB-FAIL-1684", "LIB-FAIL-1685", "LIB-FAIL-1686", "LIB-FAIL-1687", "LIB-FAIL-1688", "LIB-FAIL-1689"]
  first_seen: "2026-09-01"
  last_seen: "2026-09-01"
  recurrence_count: 2
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: IMPLEMENTATION
  phase: maintenance
  capability: "SUPPLY-WAREHOUSE-PACKAGING"
  scope: "Packaging-aware picking and bin replenishment governed by fixed Microsoft BCApps owners"
  command_redacted: "materialize supply, compose backend, migrate PostgreSQL, run focal/full suites and global gates"
  toolchain: ["Go 1.26.7", "PostgreSQL 18.6", "PowerShell 7", "ripgrep"]
  error_signature_redacted: "partial-module probe; mutable-state FK; cleanup/fixture/query contamination; stale aggregate metadata"
  violated_invariant: "No packaging-aware activity may promote unless physical quantity, exact composition, provenance, registration, cancellation and clean reconstruction all agree."
  blast_radius: "All failed attempts stayed in isolated trees/databases or were rejected by canonical gates; no project production state exists."
  reproduction: "Use the invalid 0036 FK, omit item-bin contracts, reuse a contaminated gate database, or leave provenance counts at V170."
  hypotheses_rejected: ["Base quantity reservation alone protects packages", "A historical FK may target a consumable composition row", "An overlay PASS is equivalent to canonical Markdown reconstruction"]
  root_cause: "The first adaptation did not yet model packaging reservations as immutable snapshots independent from mutable composition state, and its gate fixtures/metadata were incomplete."
  canonical_correction: "Resolve target UOM first; require explicit breakbulk; persist From/To factors; reserve base and source composition; consume or release atomically; retain immutable snapshots; rebuild 78/78 and compose/test canonically."
  regression_test: "Exact BOX pick, unauthorized and authorized BOX-to-EA pick/replenishment, competing conversion rejection, cancellation, registration, conservation, 0036 down/up, full PostgreSQL and Go gates."
  clean_rebuild_evidence: "MICROSOFT_BC_PICK_REPLENISHMENT_PACKAGING_2026-09-01_V171.md; root 94/1087/528/backend370 and executable Audit 94/121 PASS."
  dependent_gates_rerun: ["supply round-trip", "ENTERPRISE_BACKEND composition", "PostgreSQL 0001-0036", "go test/vet/build", "VERIFY_LIBRARY.ps1", "VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Audit"]
  upstream_reference: "microsoft/BCApps@2eae56d704a1fd035d104f333602aea7091b7749"
  residual_risk: "Project-specific demand classes, calendars, operator UX, load, recovery, integrations and production acceptance remain conditioned."
  owner: "workspace-owner"
  next_action: "Select the next project-required conditioned capability; do not infer production readiness from this narrow warehouse proof."

- failure_id: FAIL-20260901-006
  fingerprint: "v172-connected-demand-discovery-truncation"
  related_failures: ["LIB-FAIL-1690", "LIB-FAIL-1691", "LIB-FAIL-1692", "LIB-FAIL-1693", "LIB-FAIL-1694", "LIB-FAIL-1695", "LIB-FAIL-1696", "LIB-FAIL-1697", "LIB-FAIL-1698", "LIB-FAIL-1699", "LIB-FAIL-1700", "LIB-FAIL-1701", "LIB-FAIL-1702", "LIB-FAIL-1703", "LIB-FAIL-1704"]
  first_seen: "2026-09-01"
  last_seen: "2026-09-01"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  phase: closure
  capability: "CONNECTED-DEMAND-WAREHOUSE"
  scope: "Customer-order demand connection to warehouse; service and production remain separate owners"
  command_redacted: "bounded-code discovery attempted across four owners"
  toolchain: ["PowerShell", "ripgrep"]
  error_signature_redacted: "14,779-token result truncated"
  violated_invariant: "A connection inventory may use only complete owner-scoped reads."
  blast_radius: "Read-only discovery only; no implementation changed."
  reproduction: "Search commerce, fulfillment, inventory and all HTTP handlers in one command."
  hypotheses_rejected: ["Visible matches form a complete connected-demand inventory"]
  root_cause: "Owner boundaries were not separated before searching."
  canonical_correction: "Inspect commerce order state, service parts state and warehouse demand validation independently, then join only exact owners."
  regression_test: "Exact source identities; binding replay/divergence; draft refusal; placed 3-BOX demand; base outstanding; exact/divergent pick replay; cancellation; registration; one-winner concurrent remainder; 0037 down/up; full PostgreSQL/Go/root/Audit gates."
  clean_rebuild_evidence: "MICROSOFT_BC_SALES_ORDER_WAREHOUSE_DEMAND_2026-09-01_V172.md; supply 83/83, backend 29/375, root 94/1092/529 and executable Audit 94/121 PASS."
  dependent_gates_rerun: ["owner-scoped discovery", "supply round-trip", "ENTERPRISE_BACKEND composition", "PostgreSQL 0001-0037", "go test/vet/build", "0037 down/up focal", "VERIFY_LIBRARY.ps1", "VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Audit"]
  upstream_reference: "microsoft/BCApps@2eae56d704a1fd035d104f333602aea7091b7749 exact SalesLine/source-document/warehouse tests"
  residual_risk: "Shipment/fiscal completion and service/production demand owners remain disconnected until their own exact verticals are admitted."
  owner: "workspace-owner"
  next_action: "Connect the next authoritative journey without weakening this customer-order boundary."

- failure_id: FAIL-20260902-001
  fingerprint: "v174-required-authority-batch-read-truncation"
  related_failures: ["LIB-FAIL-1725"]
  first_seen: "2026-09-02"
  last_seen: "2026-09-02"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  phase: discovery
  capability: "CONNECTED-CARRIER-DELIVERY"
  scope: "Required authority loading before V174 maintenance"
  command_redacted: "parallel Get-Content of eight mandatory authorities"
  toolchain: ["PowerShell"]
  error_signature_redacted: "46,984-token aggregate output truncated"
  violated_invariant: "Every selected mandatory instruction file must be read completely before acting."
  blast_radius: "Read-only discovery only; no V174 implementation changed."
  reproduction: "Read all mandatory authorities into one aggregated tool response."
  hypotheses_rejected: ["Parallel file reads imply complete model-visible contents"]
  root_cause: "The output budget was not sized from file line counts before aggregation."
  canonical_correction: "Read each authority in bounded line ranges, verify contiguous coverage to EOF, then close this failure before implementation."
  regression_test: "Recorded line counts plus contiguous bounded reads for all eight selected authorities."
  clean_rebuild_evidence: "Line/byte inventory plus eight individual complete reads through EOF before V174 implementation."
  dependent_gates_rerun: ["mandatory authority loading", "VERIFY_LIBRARY.ps1"]
  upstream_reference: "None; local protocol only"
  residual_risk: "Later authorities must still be loaded only when the owner inventory makes them material."
  owner: "workspace-owner"
  next_action: "Proceed with owner-scoped V174 discovery and keep every further read bounded."

- failure_id: FAIL-20260902-002
  fingerprint: "v174-cross-owner-inventory-truncation"
  related_failures: ["LIB-FAIL-1726"]
  first_seen: "2026-09-02"
  last_seen: "2026-09-02"
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  phase: discovery
  capability: "CONNECTED-CARRIER-DELIVERY"
  scope: "Shipment, carrier tracking and customer handover owner inventory"
  command_redacted: "cross-pack ripgrep over six implementation packs and plans"
  toolchain: ["ripgrep"]
  error_signature_redacted: "45,536-token output truncated"
  violated_invariant: "An owner inventory may use only complete owner-scoped reads."
  blast_radius: "Read-only discovery; no V174 implementation changed."
  reproduction: "Search generic carrier/tracking/handover terms across all large packs together."
  hypotheses_rejected: ["Visible cross-pack matches are an exhaustive ownership map"]
  root_cause: "The query mixed domain boundaries and broad terms before locating exact symbols."
  canonical_correction: "Inspect each pack metadata, exact service/repository symbols, DDL and focal tests separately; join only verified boundaries."
  regression_test: "Bounded owner-specific symbol inventories with no truncated response."
  clean_rebuild_evidence: "Exact DDL, service/repository/HTTP and provider operation reads completed without truncation."
  dependent_gates_rerun: ["owner inventory", "VERIFY_LIBRARY.ps1"]
  upstream_reference: "Microsoft BCApps and Amazon SP-API exact revisions already locked; no claim yet"
  residual_risk: "Implementation must still prove the cross-owner binding and transition mapping."
  owner: "workspace-owner"
  next_action: "Design the minimal cross-owner binding without merging their responsibilities."

- failure_id: FAIL-20260902-003
  fingerprint: "v174-materialized-backend-symbol-search-truncation"
  related_failures: ["LIB-FAIL-1727"]
  first_seen: "2026-09-02"
  last_seen: "2026-09-02"
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  phase: discovery
  capability: "CONNECTED-CARRIER-DELIVERY"
  scope: "Materialized V173 backend owner inventory"
  command_redacted: "ripgrep over all internal Go and migration SQL"
  toolchain: ["ripgrep"]
  error_signature_redacted: "10,197-token output truncated"
  violated_invariant: "Discovery must be complete at each selected owner boundary."
  blast_radius: "Read-only temporary V173 tree; no canonical or temporary file changed."
  reproduction: "Search all shipment and handover symbols across internal and migrations together."
  hypotheses_rejected: ["A smaller global search is necessarily bounded enough"]
  root_cause: "The query still combined every downstream test and migration after exact paths were known."
  canonical_correction: "Read exact DDL, service, repository, HTTP and focal test ranges separately for sales shipment, logistics shipment and handover."
  regression_test: "Owner-scoped bounded reads complete without truncation."
  clean_rebuild_evidence: "Exact logistics/customer shipment/handover DDL and owner code ranges read completely."
  dependent_gates_rerun: ["owner inventory", "VERIFY_LIBRARY.ps1"]
  upstream_reference: "No new upstream claim yet"
  residual_risk: "The generic logistics owner still lacks durable idempotency/version/provider-event evidence."
  owner: "workspace-owner"
  next_action: "Correct those gaps in the canonical owner with regression tests."

- failure_id: FAIL-20260902-004
  fingerprint: "v174-bcapps-package-test-generic-pattern-truncation"
  related_failures: ["LIB-FAIL-1728"]
  first_seen: "2026-09-02"
  last_seen: "2026-09-02"
  recurrence_count: 3
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  phase: authority-admission
  capability: "CONNECTED-CARRIER-DELIVERY"
  scope: "Exact Microsoft BCApps sales serial/package tests"
  command_redacted: "generic tracking/package/shipment search inside one long AL test file"
  toolchain: ["ripgrep"]
  error_signature_redacted: "15,641-token result truncated"
  violated_invariant: "An official test claim must be backed by complete selected procedures."
  blast_radius: "Read-only fixed Microsoft snapshot; no source was copied or adapted."
  reproduction: "Search every package/tracking occurrence in the 2,400-line test codeunit."
  hypotheses_rejected: ["One-file scope alone guarantees bounded output"]
  root_cause: "The terms occur throughout every package-tracking scenario in the codeunit."
  canonical_correction: "Enumerate procedure signatures, select exact sales serial/package scenarios and read only their contiguous ranges."
  regression_test: "Selected official procedure ranges plus exact hashes and local vertical tests."
  clean_rebuild_evidence: "Complete BCApps procedure ranges for sales serial/package posting and shipping-agent persistence were read from the fixed snapshot."
  dependent_gates_rerun: ["authority inventory", "VERIFY_LIBRARY.ps1"]
  upstream_reference: "microsoft/BCApps fixed snapshot; claim not yet promoted"
  residual_risk: "BCApps does not prove our Go mapping, provider webhook or business acceptance; those remain local gates."
  owner: "workspace-owner"
  next_action: "Bind the narrow Microsoft claims to local tests without copying AL as Go."

- failure_id: FAIL-20260902-005
  fingerprint: "v174-amazon-adapter-cycle-search-truncation"
  related_failures: ["LIB-FAIL-1729"]
  first_seen: "2026-09-02"
  last_seen: "2026-09-02"
  recurrence_count: 4
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  phase: authority-admission
  capability: "CONNECTED-CARRIER-DELIVERY"
  scope: "Amazon Shipping tracking, Easy Ship handover and MCF delivery evidence"
  command_redacted: "combined function and status search over three complete adapter packs"
  toolchain: ["ripgrep"]
  error_signature_redacted: "14,574-token aggregate output truncated"
  violated_invariant: "Provider contracts must be read at the exact operation boundary used by the claim."
  blast_radius: "Read-only canonical Markdown; no adapter or profile changed."
  reproduction: "Search every lifecycle function and status term across the full Shipping V2 cycle."
  hypotheses_rejected: ["Filtering by definitions is enough when one pack embeds multiple large files"]
  root_cause: "The Shipping V2 pack combines seven operations and their complete tests."
  canonical_correction: "Read only query_tracking, Easy Ship reconcile_package and fetch_delivery_evidence plus their focal tests and non-claims."
  regression_test: "Bounded exact-range reads and existing pack tests for those three operations."
  clean_rebuild_evidence: "Exact tracking query, Easy Ship reconciliation and delivery-evidence receipt/test ranges read completely."
  dependent_gates_rerun: ["provider authority inventory", "VERIFY_LIBRARY.ps1"]
  upstream_reference: "Amazon SP-API SDK/models exact revisions already locked; no new claim yet"
  residual_risk: "Provider-specific reports cannot drive internal transitions until a validated mapping exists."
  owner: "workspace-owner"
  next_action: "Implement an explicit provider-event mapping that never auto-accepts customer delivery."

- failure_id: FAIL-20260902-006
  fingerprint: "v174-logistics-null-provider-reference-uniqueness"
  related_failures: ["LIB-FAIL-1730"]
  first_seen: "2026-09-02"
  last_seen: "2026-09-02"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: CODE
  phase: design
  capability: "CONNECTED-CARRIER-DELIVERY"
  scope: "logistics.shipment provider identity"
  command_redacted: "Exact DDL review of logistics.shipment"
  toolchain: ["PostgreSQL 18.6"]
  error_signature_redacted: "UNIQUE NULLS NOT DISTINCT admits only one NULL provider reference per tenant/provider"
  violated_invariant: "Multiple planned shipments may exist before the provider assigns an external reference."
  blast_radius: "Any second shipment for the same provider without a reference conflicts in the current schema."
  reproduction: "Insert two rows with the same tenant/provider and NULL provider_reference."
  hypotheses_rejected: ["A null provider reference should participate as a durable external identity"]
  root_cause: "The initial schema treated a not-yet-assigned provider identity as a unique value."
  canonical_correction: "Drop the nulls-not-distinct constraint and create a unique partial index where provider_reference is not null."
  regression_test: "Two NULL references succeed; repeated non-NULL reference fails; 0039 down/up preserves the intended baseline."
  clean_rebuild_evidence: "Fresh PostgreSQL 18.6 migrations 0001-0039, two-NULL/duplicate-non-NULL SQL test, complete integration, full Go suite and 0039 down/up all PASS."
  dependent_gates_rerun: ["PostgreSQL migrations/tests", "fulfillment integration", "Go suite", "root", "Audit"]
  upstream_reference: "Microsoft BCApps separates Shipping Agent Code from Package Tracking No.; exact fixed paths inspected"
  residual_risk: "Provider reference assignment and event reconciliation still require idempotent evidence."
  owner: "workspace-owner"
  next_action: "Retain the partial unique-index test in every reconstruction."

- failure_id: FAIL-20260902-007
  fingerprint: "v174-connected-transport-integration-patch-output-truncated"
  related_failures: ["LIB-FAIL-1731"]
  first_seen: "2026-09-02"
  last_seen: "2026-09-02"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  phase: implementation
  capability: "CONNECTED-CARRIER-DELIVERY"
  scope: "PostgreSQL integration-test source"
  command_redacted: "Single apply_patch containing the complete integration test"
  toolchain: ["apply_patch", "Go"]
  error_signature_redacted: "Wrapper output exceeded model context; atomic patch result was not observable"
  violated_invariant: "A source edit must have an observable, byte-verifiable outcome before it can participate in a gate."
  blast_radius: "One temporary integration-test file; no canonical pack was updated."
  reproduction: "Apply a 15 KB source file in one patch through the orchestration wrapper."
  hypotheses_rejected: ["No visible patch result means the file was not written"]
  root_cause: "The orchestration response, not the filesystem edit, exceeded the available response context."
  canonical_correction: "Inspect exact destination bytes and boundaries, then use gofmt, compilation and PostgreSQL execution as the integrity proof."
  regression_test: "gofmt plus package compilation and the connected-carrier integration test from a clean database."
  clean_rebuild_evidence: "File observed complete, gofmt/parser/compile PASS and PostgreSQL integration PASS."
  dependent_gates_rerun: ["Go compile", "PostgreSQL integration", "full Go suite", "root", "Audit"]
  upstream_reference: "None; local procedural failure"
  residual_risk: "An internal malformed region remains possible until parser/compiler/test execution succeeds."
  owner: "workspace-owner"
  next_action: "Keep large source additions split and always verify the destination before use."

- failure_id: FAIL-20260902-008
  fingerprint: "v174-connected-transport-empty-bin-minimum"
  related_failures: ["LIB-FAIL-1732"]
  first_seen: "2026-09-02"
  last_seen: "2026-09-02"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TEST
  phase: integration
  capability: "CONNECTED-CARRIER-DELIVERY"
  scope: "Warehouse fixture item-bin policy"
  command_redacted: "Connected-carrier PostgreSQL integration test"
  toolchain: ["Go 1.26.7", "PostgreSQL 18.6"]
  error_signature_redacted: "SQLSTATE 22P02 invalid input syntax for numeric: empty string"
  violated_invariant: "A required numeric policy boundary must be explicit in every integration fixture."
  blast_radius: "Fixture setup only; no receipt, shipment, carrier report or business state was created."
  reproduction: "Call ConfigureItemBin with an empty MinQuantity."
  hypotheses_rejected: ["The repository defaults an omitted minimum to zero"]
  root_cause: "Only MaxQuantity uses nullif; MinQuantity is a mandatory numeric cast."
  canonical_correction: "Set MinQuantity to the exact numeric string zero in every V174 bin policy."
  regression_test: "Rerun the complete connected-carrier integration on the cleaned tenant."
  clean_rebuild_evidence: "Focal journey, full Go suite, SQL suite, vet/build and PostgreSQL 18.6 migration reversal PASS."
  dependent_gates_rerun: ["connected-carrier integration", "full Go suite", "root", "Audit"]
  upstream_reference: "None; local test contract"
  residual_risk: "Later fixture steps remain unexecuted until rerun."
  owner: "workspace-owner"
  next_action: "Preserve explicit numeric minima in future warehouse fixtures."

- failure_id: FAIL-20260902-009
  fingerprint: "v174-new-pack-blocks-missing-operation"
  related_failures: ["LIB-FAIL-1733"]
  first_seen: "2026-09-02"
  last_seen: "2026-09-02"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: PACKAGING
  phase: materialization
  capability: "CONNECTED-CARRIER-DELIVERY"
  scope: "Six new fulfillment pack blocks"
  command_redacted: "Materialize GO_FULFILLMENT_SERVICE_FRANCHISE_API 0.4.0"
  toolchain: ["materialize_markdown_pack.ps1"]
  error_signature_redacted: "Block metadata field is missing: operation"
  violated_invariant: "Every executable pack block must declare the complete canonical metadata contract."
  blast_radius: "Materialization stopped before the first new file; no partial reconstruction is evidence."
  reproduction: "Materialize a block containing hash/source but no operation."
  hypotheses_rejected: ["Path and language metadata substitute for operation/provenance/license"]
  root_cause: "The new block template was written from memory instead of copying an existing canonical block."
  canonical_correction: "Add operation, provenance, license, variables and secrets_allowed to all six blocks."
  regression_test: "Fresh 15-file materialization and byte-for-byte comparison with the validated tree."
  clean_rebuild_evidence: "Fresh materialization 15/15 and byte-for-byte SHA-256 comparison with the validated tree PASS."
  dependent_gates_rerun: ["pack roundtrip", "backend composition", "root", "Audit"]
  upstream_reference: "None; local pack schema"
  residual_risk: "Other pack schema constraints may surface on the next clean run."
  owner: "workspace-owner"
  next_action: "Copy canonical block metadata before adding future files."

- failure_id: FAIL-20260902-010
  fingerprint: "v174-pack-metadata-patch-assumed-hashes"
  related_failures: ["LIB-FAIL-1734"]
  first_seen: "2026-09-02"
  last_seen: "2026-09-02"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  phase: materialization
  capability: "CONNECTED-CARRIER-DELIVERY"
  scope: "Canonical fulfillment pack edit"
  command_redacted: "Multi-file apply_patch anchored to inferred SHA-256 values"
  toolchain: ["apply_patch"]
  error_signature_redacted: "Expected block lines not found"
  violated_invariant: "Canonical edits must anchor to exact observed bytes, especially after a mechanical updater."
  blast_radius: "None; apply_patch rejected every hunk atomically."
  reproduction: "Patch a current block using a stale or inferred hash line."
  hypotheses_rejected: ["The source-tree hash can be safely transcribed without reading the updated pack"]
  root_cause: "The patch was composed before inspecting the updater's exact output."
  canonical_correction: "Read exact blocks and add only missing metadata around stable block identifiers."
  regression_test: "Fresh materialization plus 15-file hash roundtrip."
  clean_rebuild_evidence: "Exact observed block patch, fresh 15-file materialization and hash roundtrip PASS."
  dependent_gates_rerun: ["pack roundtrip", "root", "Audit"]
  upstream_reference: "None; local edit discipline"
  residual_risk: "None from the rejected patch; the original missing metadata remains open."
  owner: "workspace-owner"
  next_action: "Never infer updater output; inspect it before canonical patches."

- failure_id: FAIL-20260904-001
  fingerprint: "post-deepseek-unknown-franchise-accelerator-root-file"
  related_failures: ["LIB-FAIL-1735"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: PACKAGING
  phase: library-reaudit
  capability: "FRANCHISE-ACCELERATOR"
  scope: "Root file allowlist and library gate"
  command_redacted: "VERIFY_LIBRARY.ps1"
  toolchain: ["PowerShell 7"]
  error_signature_redacted: "VERIFY_LIBRARY_FAILED unknown top-level file FRANCHISE_ACCELERATOR.md"
  violated_invariant: "Every root entrypoint must be explicitly governed and the complete library must pass before readiness claims are accepted."
  blast_radius: "The gate stopped before validating the 142-pack/1264-file claim and the new franchise composition."
  reproduction: "Run VERIFY_LIBRARY.ps1 against the current workspace."
  hypotheses_rejected: ["A useful new Markdown root file is automatically admitted by the existing verifier"]
  root_cause: "New root surface and verifier allowlist evolved independently."
  canonical_correction: "Audit ownership and references, then update the canonical allowlist/test only if the root placement is intentional."
  regression_test: "Root verifier, executable Audit and franchise composition from an empty target."
  clean_rebuild_evidence: "Second verifier run passed the root-file inventory and stopped only on this intentionally open row, proving both intended entrypoints are now governed."
  dependent_gates_rerun: ["VERIFY_LIBRARY.ps1", "VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Audit", "FRANCHISE_COMPLETE composition"]
  upstream_reference: "None; local control-plane defect"
  residual_risk: "Additional post-V174 inconsistencies may remain behind the first fail-fast gate."
  owner: "workspace-owner"
  next_action: "Continue the full post-V223 verification; do not infer downstream PASS from this focal correction."

- failure_id: FAIL-20260904-002
  fingerprint: "post-deepseek-materialization-cardinality-overcount"
  related_failures: ["LIB-FAIL-1736"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: EVIDENCE
  phase: library-reaudit
  capability: "LIBRARY-PROVENANCE"
  scope: "142-pack materialization inventory"
  command_redacted: "VERIFY_LIBRARY.ps1"
  toolchain: ["PowerShell 7"]
  error_signature_redacted: "Expected AUTHORED=1052 ADAPTED=105 VERBATIM=105 TOTAL=1262"
  violated_invariant: "Readiness and notices must use the verifier's materialization-block inventory, not document counts."
  blast_radius: "The advertised 1264 count and current provenance notice are invalid until corrected."
  reproduction: "Count canonical block provenance across implementation_packs and compare to README/readiness/notices."
  hypotheses_rejected: ["Root accelerator entrypoints are executable pack blocks"]
  root_cause: "Two root Markdown entrypoints were added to the reconstructed-file total even though they have no materialization blocks."
  canonical_correction: "Use 1262 total with exact provenance 1052/105/105 in every current-state consumer."
  regression_test: "Root verifier and executable Audit."
  clean_rebuild_evidence: "All current-state consumers now use the verifier-calculated 1052/105/105 = 1262 inventory; full rerun remains the next independent gate."
  dependent_gates_rerun: ["VERIFY_LIBRARY.ps1", "VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Audit"]
  upstream_reference: "None; local evidence accounting"
  residual_risk: "Further profile or pack inconsistencies may remain behind this fail-fast gate."
  owner: "workspace-owner"
  next_action: "Rerun root gate and continue fail-fast audit."

- failure_id: FAIL-20260904-003
  fingerprint: "post-deepseek-franchise-plan-duplicate-license-pack"
  related_failures: ["LIB-FAIL-1737"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: COMPOSITION
  phase: library-reaudit
  capability: "FRANCHISE-ACCELERATOR"
  scope: "FRANCHISE_COMPLETE_PACK_PLAN"
  command_redacted: "VERIFY_LIBRARY.ps1"
  toolchain: ["PowerShell 7", "Markdown compositor"]
  error_signature_redacted: "duplicate packId DEPENDENCY-LICENSE-EVIDENCE-CORE"
  violated_invariant: "A composition selects each pack ID exactly once."
  blast_radius: "The advertised 40-pack franchise plan cannot materialize."
  reproduction: "Parse the plan and group packs by packId."
  hypotheses_rejected: ["Backend 29 plus web 6 are disjoint sets"]
  root_cause: "The license-evidence pack belongs to both profiles and was counted/selected twice."
  canonical_correction: "Keep its backend selection, remove the web duplicate and state 39 unique packs."
  regression_test: "Root verifier and clean franchise composition."
  clean_rebuild_evidence: "Compositor 0.2.1 passed its tests; the corrected plan materialized 510 files from 39 unique packs in an empty destination; Go 1.26.7 test/vet/build and web frozen install/typecheck/38 tests/build passed."
  dependent_gates_rerun: ["VERIFY_LIBRARY.ps1", "FRANCHISE_COMPLETE composition", "Go/web build", "Audit"]
  upstream_reference: "None; local composition defect"
  residual_risk: "Project-specific PostgreSQL, browser, provider and production gates remain separate and cannot be inherited from this composition proof."
  owner: "workspace-owner"
  next_action: "Rerun root/Audit gates and retain project-specific runtime gates."

- failure_id: FAIL-20260904-004
  fingerprint: "post-deepseek-franchise-manifest-one-line-patch-anchor"
  related_failures: ["LIB-FAIL-1738"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  phase: library-reaudit
  capability: "FRANCHISE-ACCELERATOR"
  scope: "Single-line JSON manifest edit"
  command_redacted: "Grouped apply_patch across plan and documentation"
  toolchain: ["apply_patch"]
  error_signature_redacted: "Expected serialized JSON fragment not found"
  violated_invariant: "A long serialized line must be transformed from exact observed bytes, not a partial visual excerpt."
  blast_radius: "None; patch was rejected atomically."
  reproduction: "Use a partial line as an apply_patch deletion anchor."
  hypotheses_rejected: ["The visible rg excerpt is a complete patchable line"]
  root_cause: "The plan stores all pack objects on one physical line."
  canonical_correction: "Perform an exact counted mechanical object removal and validate JSON/duplicate count immediately."
  regression_test: "Parse JSON fence, assert 39 unique IDs, then compose cleanly."
  clean_rebuild_evidence: "An exact byte-counted transformation found two identical objects, removed only the second, reparsed the fenced JSON and proved 39 selections with 39 unique pack IDs."
  dependent_gates_rerun: ["manifest parse", "root", "franchise composition"]
  upstream_reference: "None; local edit discipline"
  residual_risk: "The full 39-pack composition remains governed separately by FAIL-20260904-003."
  owner: "workspace-owner"
  next_action: "Keep long serialized manifests subject to counted parse-and-assert edits."

- failure_id: FAIL-20260904-005
  fingerprint: "post-deepseek-serverless-plan-missing-json"
  related_failures: ["LIB-FAIL-1739"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: COMPOSITION
  phase: library-reaudit
  capability: "FRANCHISE-SERVERLESS"
  scope: "FRANCHISE_SERVERLESS_PACK_PLAN"
  command_redacted: "VERIFY_LIBRARY.ps1"
  toolchain: ["PowerShell 7", "Markdown compositor"]
  error_signature_redacted: "JSON plan missing"
  violated_invariant: "A file named *_PACK_PLAN.md must contain a compositor-readable JSON manifest."
  blast_radius: "The advertised serverless lane cannot be materialized as documented."
  reproduction: "Run the root verifier or search the document for its fenced JSON plan."
  hypotheses_rejected: ["A prose table is a sufficient executable composition"]
  root_cause: "The new file was authored as guidance while being wired as an executable pack plan."
  canonical_correction: "Add an exact, deduplicated and conditioned manifest from canonical pack metadata, then compose it in an empty destination."
  regression_test: "Root verifier plus clean serverless composition and applicable build gates."
  clean_rebuild_evidence: "A 45-unique-pack manifest materialized 537 files into an empty destination; Go 1.26.7 test/vet/build and web frozen install/typecheck/38 tests/build passed."
  dependent_gates_rerun: ["VERIFY_LIBRARY.ps1", "serverless composition", "Go/web build", "Audit"]
  upstream_reference: "None; local composition contract"
  residual_risk: "Cloud account, billing, exact deploy adapter and live provider/production gates remain CONDITIONED as stated."
  owner: "workspace-owner"
  next_action: "Rerun root/Audit and retain live cloud/provider gates per project."

- failure_id: FAIL-20260904-006
  fingerprint: "post-deepseek-pack-plan-regression-coverage-drift"
  related_failures: ["LIB-FAIL-1740"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: TEST-COVERAGE
  phase: library-reaudit
  capability: "PACK-COMPOSITION"
  scope: "VERIFY_LIBRARY profile coverage"
  command_redacted: "VERIFY_LIBRARY.ps1"
  toolchain: ["PowerShell 7", "Markdown compositor"]
  error_signature_redacted: "discovered=44 composed=42"
  violated_invariant: "Every pack plan must be composed exactly once by the canonical root gate."
  blast_radius: "Two flagship franchise plans could regress without failing the root gate."
  reproduction: "Compare discovered *_PACK_PLAN.md names with profileResults in the verifier."
  hypotheses_rejected: ["Manual one-off composition is equivalent to permanent regression coverage"]
  root_cause: "The two new franchise plans were omitted from the verifier's explicit profile set."
  canonical_correction: "Add both exact plan names and rerun the full root gate."
  regression_test: "VERIFY_LIBRARY.ps1 must report 44 composed profiles exactly once."
  clean_rebuild_evidence: "Both exact plan names were added to the canonical profile set after independent clean compositions produced 39/510 and 45/537."
  dependent_gates_rerun: ["VERIFY_LIBRARY.ps1", "VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Audit"]
  upstream_reference: "None; local verifier coverage"
  residual_risk: "Future additions remain protected by the discovered-versus-composed equality assertion."
  owner: "workspace-owner"
  next_action: "Rerun root/Audit and keep the equality assertion active."

- failure_id: FAIL-20260904-007
  fingerprint: "omnichannel-audit-broad-search-truncated-and-wrong-paths"
  related_failures: ["LIB-FAIL-1741"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: DISCOVERY
  phase: library-reaudit
  capability: "OMNICHANNEL-LEAD-INGESTION"
  scope: "Pack inventory search"
  command_redacted: "rg across broad lead/channel terms"
  toolchain: ["ripgrep"]
  error_signature_redacted: "truncated output plus four missing paths"
  violated_invariant: "Audit evidence must be complete and use observed canonical paths."
  blast_radius: "The first search cannot support a completeness claim."
  reproduction: "Search broad patterns across large pack bodies before enumerating filenames."
  hypotheses_rejected: ["Adapter filenames omit their implementation language prefix"]
  root_cause: "Owners and exact filenames were not separated before content search."
  canonical_correction: "Inventory exact paths first and inspect each owner with narrow symbols/sections."
  regression_test: "Bounded per-owner reads with no truncation or missing-file diagnostics."
  clean_rebuild_evidence: "Canonical Python adapter names and six Go owner packs were enumerated successfully."
  dependent_gates_rerun: ["per-owner content audit"]
  upstream_reference: "None; local discovery discipline"
  residual_risk: "Cross-owner gaps must still be synthesized after bounded inspection."
  owner: "workspace-owner"
  next_action: "Inspect capture, attribution, messaging and outcome owners separately."

- failure_id: FAIL-20260904-008
  fingerprint: "omnichannel-audit-grouped-pack-read-truncated"
  related_failures: ["LIB-FAIL-1742"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: LOW
  classification: DISCOVERY
  phase: library-reaudit
  capability: "OMNICHANNEL-LEAD-INGESTION"
  scope: "Pack content inspection"
  command_redacted: "Grouped first/tail reads across eleven packs"
  toolchain: ["PowerShell"]
  error_signature_redacted: "output truncated"
  violated_invariant: "Selected instruction/evidence must be read without hidden output."
  blast_radius: "The grouped output is only a router, not completeness evidence."
  reproduction: "Aggregate first and evidence sections across eleven large Markdown packs."
  hypotheses_rejected: ["Header/tail extraction remains small enough when repeated eleven times"]
  root_cause: "The audit batch was still too broad."
  canonical_correction: "Inspect materialized owner files and one pack at a time with exact symbols."
  regression_test: "Every subsequent output stays below the tool limit and contains no truncation warning."
  clean_rebuild_evidence: "The truncated batch was explicitly invalidated before gap conclusions."
  dependent_gates_rerun: ["bounded per-owner audit"]
  upstream_reference: "None; local discovery discipline"
  residual_risk: "Owner-level semantics remain to be checked."
  owner: "workspace-owner"
  next_action: "Use the 537-file materialized tree for bounded code inspection."

- failure_id: FAIL-20260904-009
  fingerprint: "deepseek-app-wiring-non-null-is-not-e2e-connection"
  related_failures: ["LIB-FAIL-1743"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: CRITICAL
  classification: INTEGRATION
  phase: library-reaudit
  capability: "AGENT-RUNTIME-WIRING"
  scope: "GO-APP-WIRING 0.1.0"
  command_redacted: "Materialize and inspect config.go/app.go/app_test.go"
  toolchain: ["Go 1.26.7"]
  error_signature_redacted: "Non-nil components asserted without an executable connected journey"
  violated_invariant: "A capability is complete only when its real end-to-end effect and recovery are connected and tested."
  blast_radius: "The advertised chatbot runtime cannot currently prove automatic lead-to-appointment/sale execution."
  reproduction: "Inspect New: Agent receives router/tools/safety only; LLM, economy, approval and FinOps remain detached fields."
  hypotheses_rejected: ["Object construction proves runtime integration"]
  root_cause: "The test checks pointer presence instead of behavior across component boundaries."
  canonical_correction: "Implement an orchestrator with explicit interfaces and E2E tests spanning ingress, durable event, response policy, domain tool, audit and handoff; then update plans."
  regression_test: "Provider fixture → durable lead → attributed conversation → governed response → fake domain appointment plus failure/replay/handoff cases."
  clean_rebuild_evidence: "V236: app E2E with controlled Responses/domain HTTP plus PostgreSQL 18.6 proved exact replay, one domain effect, full 48-package suite, vet/build and 17/17 Markdown round-trip."
  dependent_gates_rerun: ["new omnichannel core tests", "Go full suite", "PostgreSQL integration", "franchise composition", "root", "Audit"]
  upstream_reference: "Official provider webhook/API contracts plus the repository's connected-capability rule"
  residual_risk: "Other V181-V223 packs may likewise be isolated packages rather than wired runtime capabilities."
  owner: "workspace-owner"
  next_action: "Re-run the connected gates in each selected project with live identity/provider configuration."

- failure_id: FAIL-20260904-010
  fingerprint: "new-binance-rust-raw-log-path-assumed"
  related_failures: ["LIB-FAIL-1744"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: DISCOVERY
  phase: external-project-pattern-audit
  capability: "DURABLE-REALTIME-INGESTION"
  scope: "NEW BINANCE Rust module inventory"
  command_redacted: "rg targeted symbols across assumed Rust paths"
  toolchain: ["ripgrep"]
  error_signature_redacted: "raw_log.rs not found"
  violated_invariant: "Read paths must come from observed inventory."
  blast_radius: "The combined query cannot prove complete module coverage."
  reproduction: "Assume the Python/conceptual raw_log name exists in Rust."
  hypotheses_rejected: ["Rust mirrors the Python module filename"]
  root_cause: "A conceptual component name was treated as a physical path."
  canonical_correction: "Use observed segment_chain, durability_progress, transport_journal, ownership, liveness and hot_redundancy paths only."
  regression_test: "All later NEW BINANCE reads use paths returned by rg --files."
  clean_rebuild_evidence: "Exact Rust paths were already enumerated; no source was modified."
  dependent_gates_rerun: ["bounded NEW BINANCE pattern audit"]
  upstream_reference: "User-owned repository"
  residual_risk: "No code may be copied until ownership/license is recorded."
  owner: "workspace-owner"
  next_action: "Inspect only observed files and extract design invariants, not financial semantics."

- failure_id: FAIL-20260904-011
  fingerprint: "omnichannel-official-research-batch-truncated"
  related_failures: ["LIB-FAIL-1745"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: DISCOVERY
  phase: authority-research
  capability: "OMNICHANNEL-LEAD-INGESTION"
  scope: "AWS, Google Cloud, CloudEvents and Google SRE official architecture sources"
  command_redacted: "Four official-source web searches in one request"
  toolchain: ["web search"]
  error_signature_redacted: "Output exceeded context and was truncated"
  violated_invariant: "Authority evidence must be complete enough to verify each narrow claim."
  blast_radius: "No architecture claim may use the truncated output."
  reproduction: "Batch four broad documentation searches into one long response."
  hypotheses_rejected: ["A truncated search result is sufficient evidence"]
  root_cause: "Too many broad authorities were queried together."
  canonical_correction: "Search and open one official authority per invariant and preserve exact URLs."
  regression_test: "Every claim in the re-audit names an opened primary source."
  clean_rebuild_evidence: "Not applicable; read-only discovery was repeated in bounded calls."
  dependent_gates_rerun: ["authority-link verification", "post-DeepSeek re-audit"]
  upstream_reference: "Official documentation only"
  residual_risk: "Sources may change and remain governed by freshness checks."
  owner: "workspace-owner"
  next_action: "Repeat bounded searches and record exact official documents."

- failure_id: FAIL-20260904-012
  fingerprint: "readiness-root-path-assumed"
  related_failures: ["LIB-FAIL-1746"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  phase: library-maintenance
  capability: "READINESS-EVIDENCE"
  scope: "MARKDOWN_SYSTEM_READINESS path resolution"
  command_redacted: "Get-Content against unobserved root path"
  toolchain: ["PowerShell", "ripgrep"]
  error_signature_redacted: "Path not found"
  violated_invariant: "Filesystem reads use observed paths."
  blast_radius: "No file changed; the first read did not provide readiness evidence."
  reproduction: "Read the filename from repository root without enumerating it."
  hypotheses_rejected: ["Readiness document is top-level"]
  root_cause: "Owner directory was omitted."
  canonical_correction: "Resolve with rg --files and use markdown_system/MARKDOWN_SYSTEM_READINESS.md."
  regression_test: "Exact-path read succeeds."
  clean_rebuild_evidence: "Exact-path read succeeded immediately afterward."
  dependent_gates_rerun: ["readiness review"]
  upstream_reference: "Local filesystem inventory"
  residual_risk: "None beyond future moves, which maps must record."
  owner: "workspace-owner"
  next_action: "Use observed paths only."

- failure_id: FAIL-20260904-013
  fingerprint: "git-status-on-non-git-library"
  related_failures: ["LIB-FAIL-1747"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  phase: library-maintenance
  capability: "NON-GIT-LIBRARY-VERIFICATION"
  scope: "Elite Engineering Library root"
  command_redacted: "git status --short"
  toolchain: ["Git"]
  error_signature_redacted: "Not a git repository"
  violated_invariant: "Git is not a prerequisite for this library."
  blast_radius: "No mutation and no impact on the library gates."
  reproduction: "Run git status in the library root."
  hypotheses_rejected: ["The library gained a Git repository"]
  root_cause: "A generic workspace diagnostic ignored the explicit non-Git design."
  canonical_correction: "Use native manifests, SHA-256, materialization and verification scripts."
  regression_test: "Subsequent maintenance performs no Git-dependent gate."
  clean_rebuild_evidence: "Not applicable; no files were changed by Git."
  dependent_gates_rerun: ["VERIFY_LIBRARY", "VERIFY_EXECUTABLE_LIBRARY"]
  upstream_reference: "AGENTS.md and project initialization contract"
  residual_risk: "None; optional project Git remains separate."
  owner: "workspace-owner"
  next_action: "Do not treat Git as required."

- failure_id: FAIL-20260904-014
  fingerprint: "project-ledger-open-critical-control-drift"
  related_failures: ["LIB-FAIL-1748"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: EVIDENCE
  phase: library-reaudit
  capability: "FAILURE-LEDGER-CONTROL"
  scope: "PROJECT_FAILURE_LESSONS control metadata"
  command_redacted: "Compare control list with detailed OPEN/CRITICAL entries"
  toolchain: ["PowerShell", "ripgrep"]
  error_signature_redacted: "Open critical entry absent from control list"
  violated_invariant: "The control view must expose every open high or critical failure."
  blast_radius: "A reader could miss the runtime wiring blocker."
  reproduction: "Compare open_high_or_critical with FAIL-20260904-009 status/severity."
  hypotheses_rejected: ["Detailed entries alone keep the control view consistent"]
  root_cause: "The detailed failure was added without updating control metadata."
  canonical_correction: "Synchronize the open ID and updated_at in the same change."
  regression_test: "Control list contains FAIL-20260904-009 while and only while it remains OPEN/CRITICAL."
  clean_rebuild_evidence: "Header and detailed entry now agree."
  dependent_gates_rerun: ["failure-ledger validation", "VERIFY_LIBRARY"]
  upstream_reference: "FAILURE_LEARNING_CONTRACT.md"
  residual_risk: "A machine-check should enforce this relation in a later validator revision."
  owner: "workspace-owner"
  next_action: "Keep the blocker open until behavioral E2E evidence exists."

- failure_id: FAIL-20260904-015
  fingerprint: "cross-provider-local-source-search-truncated"
  related_failures: ["LIB-FAIL-1749"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: LOW
  classification: DISCOVERY
  phase: library-reaudit
  capability: "PROVIDER-SOURCE-INVENTORY"
  scope: "Meta, Google, TikTok, CloudEvents, Pub/Sub and SQS local references"
  command_redacted: "Cross-tree rg over six provider/architecture terms"
  toolchain: ["ripgrep"]
  error_signature_redacted: "Output truncated after 20,005 tokens"
  violated_invariant: "Discovery must be complete within each owner before claiming coverage."
  blast_radius: "No absence or completeness claim can use this query."
  reproduction: "Search all markdown_system and implementation_packs for all providers together."
  hypotheses_rejected: ["Visible results form a complete source inventory"]
  root_cause: "Provider inventory and large embedded source packs were not separated."
  canonical_correction: "Use exact observed pack paths and narrow metadata/ranges per owner."
  regression_test: "Subsequent local reads open one exact owner at a time within bounded output."
  clean_rebuild_evidence: "Not applicable; no source changed."
  dependent_gates_rerun: ["bounded provider pack audit", "post-DeepSeek re-audit"]
  upstream_reference: "Local source locks and provider packs"
  residual_risk: "Unknown owners remain until bounded inventory completes."
  owner: "workspace-owner"
  next_action: "Inspect exact Meta, Google and TikTok owners separately."

- failure_id: FAIL-20260904-016
  fingerprint: "tooling-directories-assumed"
  related_failures: ["LIB-FAIL-1750"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: DISCOVERY
  phase: library-maintenance
  capability: "PACK-MATERIALIZATION-TOOLING"
  scope: "Library script inventory"
  command_redacted: "rg --files against assumed scripts and tooling directories"
  toolchain: ["ripgrep"]
  error_signature_redacted: "Directory not found"
  violated_invariant: "Paths must be observed before use."
  blast_radius: "No source changed; tool routing only."
  reproduction: "Assume conventional scripts/ and tooling/ directories exist."
  hypotheses_rejected: ["Materializer is stored in scripts/"]
  root_cause: "Conventional paths were supplied after exact root paths were already visible."
  canonical_correction: "Use the root materializer, markdown_system updater, and materialized compositor paths."
  regression_test: "Later calls use exact paths returned by VERIFY_LIBRARY.ps1 references."
  clean_rebuild_evidence: "Not applicable; read-only lookup."
  dependent_gates_rerun: ["pack materialization"]
  upstream_reference: "VERIFY_LIBRARY.ps1"
  residual_risk: "None."
  owner: "workspace-owner"
  next_action: "Use observed tooling paths only."

- failure_id: FAIL-20260904-017
  fingerprint: "composition-manifest-pack-id-casing-assumed"
  related_failures: ["LIB-FAIL-1751"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  phase: composition-audit
  capability: "PACK-SELECTION-INVENTORY"
  scope: "FRANCHISE_SERVERLESS_PACK_PLAN JSON"
  command_redacted: "Select non-existent pack_id property"
  toolchain: ["PowerShell"]
  error_signature_redacted: "Property pack_id cannot be found"
  violated_invariant: "Consumers use the observed manifest schema."
  blast_radius: "No source changed and no pack selection was inferred from empty output."
  reproduction: "Project plan.packs.pack_id instead of plan.packs.packId."
  hypotheses_rejected: ["Manifest property follows YAML pack_id spelling"]
  root_cause: "JSON casing was assumed rather than inspected."
  canonical_correction: "Inspect the first object and use exact packId casing."
  regression_test: "Projection returns exactly 45 non-empty unique IDs."
  clean_rebuild_evidence: "First object and total count were observed successfully."
  dependent_gates_rerun: ["composition selection audit"]
  upstream_reference: "MARKDOWN_COMPOSITOR_CORE manifest schema"
  residual_risk: "None after exact schema use."
  owner: "workspace-owner"
  next_action: "Repeat selection inventory using packId."

- failure_id: FAIL-20260904-018
  fingerprint: "serverless-composition-module-root-assumed"
  related_failures: ["LIB-FAIL-1752"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  phase: composition-audit
  capability: "FRANCHISE-SERVERLESS-COMPOSITION"
  scope: "Temporary clean materialization"
  command_redacted: "Compose 45 packs then assert backend/go.mod"
  toolchain: ["PowerShell", "markdown compositor 0.2.1"]
  error_signature_redacted: "Backend missing after composition"
  violated_invariant: "Postconditions use the exact plan output layout."
  blast_radius: "537 files were materialized correctly; only the diagnostic assertion failed."
  reproduction: "Assume a backend subdirectory instead of inspecting destination root."
  hypotheses_rejected: ["The serverless plan nests its primary Go module under backend"]
  root_cause: "A previous evidence layout was reused without checking this plan."
  canonical_correction: "Use the observed root go.mod and package.json."
  regression_test: "Destination root contains both files and exact owner directories."
  clean_rebuild_evidence: "Clean composition reported 537 files from 45 packs; root inventory resolved the exact paths."
  dependent_gates_rerun: ["Go full suite", "web suite"]
  upstream_reference: "FRANCHISE_SERVERLESS_PACK_PLAN.md"
  residual_risk: "None for layout; live target gates remain separate."
  owner: "workspace-owner"
  next_action: "Materialize missing packs into this observed root and test."

- failure_id: FAIL-20260904-019
  fingerprint: "pack-materializer-overlay-rejected"
  related_failures: ["LIB-FAIL-1753"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  phase: integration-fixture
  capability: "PACK-COMPOSITION"
  scope: "Four missing runtime packs over serverless fixture"
  command_redacted: "Materialize pack into non-empty destination"
  toolchain: ["materialize_markdown_pack.ps1"]
  error_signature_redacted: "Destination must be empty"
  violated_invariant: "Single-pack materialization is clean-room only."
  blast_radius: "The guard rejected before writing; base fixture is unchanged."
  reproduction: "Pass the 537-file composition as Destination to the pack materializer."
  hypotheses_rejected: ["Single-pack materializer supports overlay"]
  root_cause: "Materializer and compositor responsibilities were conflated."
  canonical_correction: "Use isolated materializations for inspection; use a compositor plan for canonical overlay."
  regression_test: "Every single-pack destination starts empty; final plan composes all selected packs once."
  clean_rebuild_evidence: "Base fixture still has the 45-pack materialization record."
  dependent_gates_rerun: ["isolated materialization", "canonical composition"]
  upstream_reference: "materialize_markdown_pack.ps1"
  residual_risk: "Manual fixture merge must assert zero path collisions."
  owner: "workspace-owner"
  next_action: "Materialize isolated and verify collision set before fixture merge."

- failure_id: FAIL-20260904-020
  fingerprint: "workspace-dependencies-dynamic-alias-unavailable"
  related_failures: ["LIB-FAIL-1754"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  phase: toolchain-discovery
  capability: "WORKSPACE-DEPENDENCIES"
  scope: "Codex bundled runtime lookup"
  command_redacted: "Legacy dynamic workspace dependency call"
  toolchain: ["Codex app"]
  error_signature_redacted: "Tool no longer available through dynamic tools"
  violated_invariant: "Use the current provider returned by the tool surface."
  blast_radius: "No files or runtimes changed."
  reproduction: "Call the deprecated dynamic alias."
  hypotheses_rejected: ["Legacy alias remains callable"]
  root_cause: "Tool provider changed."
  canonical_correction: "Call mcp__codex_app__load_workspace_dependencies."
  regression_test: "Current provider returns bundle version and paths."
  clean_rebuild_evidence: "Current MCP provider succeeded."
  dependent_gates_rerun: ["toolchain preflight"]
  upstream_reference: "Codex app tool response"
  residual_risk: "Go is not part of the returned bundle and remains separately pinned."
  owner: "workspace-owner"
  next_action: "Use current provider only."

- failure_id: FAIL-20260904-021
  fingerprint: "go-archive-extraction-window-exhausted"
  related_failures: ["LIB-FAIL-1755"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  phase: toolchain-preflight
  capability: "GO-RUNTIME"
  scope: "Official Go 1.26.7 temporary archive"
  command_redacted: "Expand archive and verify within 30 seconds"
  toolchain: ["PowerShell", "Go 1.26.7"]
  error_signature_redacted: "No final output before command window ended"
  violated_invariant: "Uncertain tool output cannot be treated as completion."
  blast_radius: "Only a temporary extraction was potentially incomplete."
  reproduction: "Combine large archive extraction and validation under a short command window."
  hypotheses_rejected: ["The silent command proves a complete extraction"]
  root_cause: "Verification was not separated from a slow extraction."
  canonical_correction: "Verify executable existence, exact hash and version in a separate bounded call."
  regression_test: "go.exe hash and go version match the admitted lock."
  clean_rebuild_evidence: "go.exe SHA-256 5463fe58fa999d74420f00ee1b36d31c3da90a57ff204159f159859567ad61fc; go1.26.7 windows/amd64."
  dependent_gates_rerun: ["Go test", "Go vet", "Go build"]
  upstream_reference: "https://go.dev/dl/go1.26.7.windows-amd64.zip"
  residual_risk: "Toolchain stays temporary and is not installed globally."
  owner: "workspace-owner"
  next_action: "Use the verified absolute go.exe path."

- failure_id: FAIL-20260904-022
  fingerprint: "aifoundation-llm-interface-file-assumed"
  related_failures: ["LIB-FAIL-1756"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: LOW
  classification: DISCOVERY
  phase: runtime-wiring-audit
  capability: "LLM-PROVIDER-CONTRACT"
  scope: "internal/aifoundation"
  command_redacted: "Read provider.go then llm.go without observed inventory"
  toolchain: ["PowerShell", "ripgrep"]
  error_signature_redacted: "Path not found"
  violated_invariant: "Read only observed filesystem paths."
  blast_radius: "No source changed; LLM interface evidence remains pending until exact files are read."
  reproduction: "Infer physical filenames from the LLM concept."
  hypotheses_rejected: ["The interface is in provider.go", "The interface is in llm.go"]
  root_cause: "Conceptual naming was substituted for directory inventory."
  canonical_correction: "Enumerate the package first and open contract.go/gateway.go."
  regression_test: "All later reads use paths returned by rg --files."
  clean_rebuild_evidence: "The exact package inventory was captured."
  dependent_gates_rerun: ["bounded interface audit", "runtime E2E"]
  upstream_reference: "Materialized GO-ML-AI-FOUNDATION pack"
  residual_risk: "None after exact owner read."
  owner: "workspace-owner"
  next_action: "Read contract.go and gateway.go only."

- failure_id: FAIL-20260904-023
  fingerprint: "tiktok-one-click-crm-page-timeout"
  related_failures: ["LIB-FAIL-1757"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: UPSTREAM_ACCESS
  phase: authority-research
  capability: "TIKTOK-LEAD-INTEGRATION"
  scope: "TikTok one-click CRM integration help page"
  command_redacted: "Open official TikTok help URL"
  toolchain: ["web fetch"]
  error_signature_redacted: "TimeoutError"
  violated_invariant: "Unavailable content cannot govern an implementation claim."
  blast_radius: "That page is excluded; no payload schema or behavior is inferred from it."
  reproduction: "Fetch the exact public help URL during this audit."
  hypotheses_rejected: ["Search snippets are equivalent to opened authority"]
  root_cause: "Upstream page did not respond within the fetch window."
  canonical_correction: "Use only opened TikTok CRM Webhooks and signal optimization pages; require portal/API schema for adapter implementation."
  regression_test: "Every TikTok claim cites an accessible official page or remains BLOCKED."
  clean_rebuild_evidence: "Two alternative official TikTok pages opened and support only the narrow real-time/postback claims."
  dependent_gates_rerun: ["TikTok portal schema acquisition", "contract tests"]
  upstream_reference: "https://ads.tiktok.com/resources/help/article/about-one-click-crm-integration"
  residual_risk: "Exact Custom API payload and auth contract remain unavailable publicly."
  owner: "workspace-owner"
  next_action: "Do not author a TikTok parser until exact official schema access exists."

- failure_id: FAIL-20260904-024
  fingerprint: "franchise-entrypoint-overstates-connected-readiness"
  related_failures: ["LIB-FAIL-1743", "LIB-FAIL-1758"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: CRITICAL
  classification: EVIDENCE
  phase: library-reaudit
  capability: "FRANCHISE-ACCELERATOR"
  scope: "FRANCHISE_ACCELERATOR.md, START_FRANCHISE.md and complete/serverless manifests"
  command_redacted: "Compare entrypoint claims with exact selected pack IDs and materialized runtime"
  toolchain: ["PowerShell", "markdown compositor", "Go 1.26.7"]
  error_signature_redacted: "Advertised connected chatbot and ads readiness not present in selected composition"
  violated_invariant: "An entrypoint may advertise only connected capabilities selected and behaviorally proven by its plan."
  blast_radius: "A new agent could start from an incomplete composition and falsely report a lead-to-sale chatbot."
  reproduction: "Materialize the 39/45-pack plans and compare selected IDs plus app behavior with the accelerator matrix."
  hypotheses_rejected: ["Pack existence implies plan inclusion", "Reporting adapters imply lead ingestion", "Non-nil app fields imply runtime wiring"]
  root_cause: "Documentation and plans were not synchronously revalidated as new packs were added."
  canonical_correction: "Downgrade entrypoint claims, publish exact re-audit, create a connected plan and prove provider fixture to durable domain outcome before restoration."
  regression_test: "Entrypoint claim verifier plus clean composition includes every required runtime owner and connected E2E tests."
  clean_rebuild_evidence: "V236: 56-pack plan fixes current versions exactly once, accelerator distinguishes Google lead/Meta reconcile/reporting/TikTok, and connected runtime evidence passes."
  dependent_gates_rerun: ["franchise composition", "Go/PostgreSQL/web E2E", "root", "Audit"]
  upstream_reference: "TOTAL_SYSTEM_CAPABILITY_CONTRACT.md and official provider contracts"
  residual_risk: "Live credentials, policy, jurisdiction and target infrastructure remain per-project gates even after library closure."
  owner: "workspace-owner"
  next_action: "Generate the project-specific subset and close live provider/browser/production gates without weakening the claims."

- failure_id: FAIL-20260904-025
  fingerprint: "postgres-owner-and-up-migration-paths-assumed"
  related_failures: ["LIB-FAIL-1759"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: DISCOVERY
  phase: omnichannel-design
  capability: "POSTGRES-LEAD-INGRESS"
  scope: "Materialized PostgreSQL owners and migration 0041"
  command_redacted: "Read assumed repository.go and unsuffixed migration"
  toolchain: ["PowerShell", "ripgrep"]
  error_signature_redacted: "Two paths not found"
  violated_invariant: "New owners follow the observed repository and migration layout."
  blast_radius: "No implementation decision uses the absent files."
  reproduction: "Assume a shared repository.go and migration without .up suffix."
  hypotheses_rejected: ["Postgres has one shared Repository", "Up migration omits .up"]
  root_cause: "Generic naming was used before local inventory."
  canonical_correction: "Add a domain-specific LeadIngress owner and paired .up/.down migration after exact inventory."
  regression_test: "New files follow observed constructor/pool and migration conventions."
  clean_rebuild_evidence: "Exact owner constructors and 0040/0041 names were enumerated."
  dependent_gates_rerun: ["Go compile", "migration up/down", "PostgreSQL integration"]
  upstream_reference: "Materialized PG foundation and analytics pack"
  residual_risk: "The new migration number must be checked against every selected pack before admission."
  owner: "workspace-owner"
  next_action: "Design 0042 only after global migration-number collision check."

- failure_id: FAIL-20260904-026
  fingerprint: "omnichannel-full-go-gates-output-truncated"
  related_failures: ["LIB-FAIL-1760"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: EVIDENCE
  phase: omnichannel-verification
  capability: "GO-OMNICHANNEL-LEAD-INGRESS"
  scope: "Full Go test, vet and build on the composed franchise tree"
  command_redacted: "Run three global Go gates in one captured invocation"
  toolchain: ["Go 1.26.7", "PowerShell"]
  error_signature_redacted: "Output exceeded available model context and was truncated"
  violated_invariant: "A gate is proven only by a complete, attributable result."
  blast_radius: "No full-suite PASS may be claimed for the new lead-ingress code."
  reproduction: "Combine the three verbose global gates under one output budget."
  hypotheses_rejected: ["A truncated result implies the commands completed successfully"]
  root_cause: "Independent verification gates were combined into an oversized output stream."
  canonical_correction: "Run test, vet and build independently with compact success output and bounded failure tails."
  regression_test: "Each gate emits its own explicit PASS or its complete relevant failure tail and exit code."
  clean_rebuild_evidence: "Go 1.26.7: go test ./... -count=1 PASS across 45 packages; go vet ./... PASS; go build ./... PASS, each run separately with complete output."
  dependent_gates_rerun: ["go test ./...", "go vet ./...", "go build ./..."]
  upstream_reference: "Go toolchain 1.26.7 exact locked archive"
  residual_risk: "PostgreSQL migration tests and connected browser journey remain separate gates."
  owner: "workspace-owner"
  next_action: "Preserve independent compact gate execution in the canonical pack evidence."

- failure_id: FAIL-20260904-027
  fingerprint: "windows-ripgrep-literal-sql-glob-recurrence"
  related_failures: ["LIB-FAIL-1704", "LIB-FAIL-1761"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 1
  status: RECURRENCE_PROVEN
  severity: LOW
  classification: DISCOVERY
  phase: omnichannel-verification
  capability: "POSTGRES-LEAD-INGRESS"
  scope: "Read-only search for owner table migrations"
  command_redacted: "Pass wildcard-looking paths directly to rg on Windows"
  toolchain: ["PowerShell", "ripgrep"]
  error_signature_redacted: "Filename, directory name, or volume label syntax is incorrect"
  violated_invariant: "Ripgrep file selection uses -g filters, not wildcard path operands."
  blast_radius: "No SQL claim uses the invalid search output."
  reproduction: "Run rg with db/migrations/*.up.sql as a path operand on Windows."
  hypotheses_rejected: ["PowerShell expands these path arguments for rg"]
  root_cause: "A known Windows glob lesson was not applied to the combined read."
  canonical_correction: "Search the real directory with explicit -g filters."
  regression_test: "The corrected command returns platform.outbox_event in 0001 and org.organization in 0002 with exit 0."
  clean_rebuild_evidence: "Corrected read-only inventory completed with exact observed files."
  dependent_gates_rerun: ["migration dependency order inspection"]
  upstream_reference: "Local failure learning contract"
  residual_risk: "None beyond enforcing the known command pattern."
  owner: "workspace-owner"
  next_action: "Use -g for all subsequent ripgrep file filters."

- failure_id: FAIL-20260904-028
  fingerprint: "go-race-requires-cgo-and-native-compiler"
  related_failures: ["LIB-FAIL-1762"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: BLOCKED_EXTERNAL
  severity: MEDIUM
  classification: TOOLCHAIN
  phase: omnichannel-verification
  capability: "GO-OMNICHANNEL-LEAD-INGRESS"
  scope: "Go race detector for concurrent duplicate ingestion"
  command_redacted: "go test -race ./internal/leadstream"
  toolchain: ["Go 1.26.7 windows/amd64"]
  error_signature_redacted: "-race requires cgo; GCC_NOT_AVAILABLE"
  violated_invariant: "Race evidence requires an available supported native toolchain."
  blast_radius: "Functional concurrency can be tested, but race-detector PASS cannot be claimed on this host."
  reproduction: "Run -race with CGO disabled and no C compiler."
  hypotheses_rejected: ["The locked Go ZIP alone is sufficient for Windows race detection"]
  root_cause: "The audit host lacks the optional native compiler required by Go race instrumentation."
  canonical_correction: "Run the same pack under a zero-additional-cost CI/host image with CGO plus supported compiler; keep promotion conditioned until then."
  regression_test: "go test -race ./internal/leadstream -count=1 exits zero on that environment."
  clean_rebuild_evidence: "Blocked externally; ordinary concurrent regression still runs on this host."
  dependent_gates_rerun: ["focal Go tests", "full Go test", "race-enabled CI"]
  upstream_reference: "Go race detector requirements"
  residual_risk: "A data race not exposed by deterministic functional assertions remains possible."
  owner: "workspace-owner"
  next_action: "Add the race gate to the canonical pack and execute it when the required native toolchain is available."

- failure_id: FAIL-20260904-029
  fingerprint: "temporary-runtime-inventory-output-truncated"
  related_failures: ["LIB-FAIL-1763"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: DISCOVERY
  phase: postgres-toolchain-discovery
  capability: "POSTGRES-LEAD-INGRESS"
  scope: "Temporary directories named elite*"
  command_redacted: "Enumerate every matching directory under the user Temp root"
  toolchain: ["PowerShell"]
  error_signature_redacted: "Output truncated after hundreds of unrelated directories"
  violated_invariant: "Discovery queries must be bounded to the decision being made."
  blast_radius: "No absence or completeness claim is derived from the truncated inventory."
  reproduction: "List every elite* temp directory without a narrow predicate."
  hypotheses_rejected: ["A broad visual listing is an efficient runtime locator"]
  root_cause: "Known exact evidence was not used to narrow runtime discovery first."
  canonical_correction: "Inspect only the observed elite-v169-postgres-runtime candidate and its pg_ctl.exe."
  regression_test: "Exact candidate path returns one pg_ctl.exe with bounded output."
  clean_rebuild_evidence: "The candidate directory was read directly and pg_ctl.exe was found."
  dependent_gates_rerun: ["runtime version/hash", "isolated PostgreSQL cluster"]
  upstream_reference: "Existing V169 reconstruction evidence"
  residual_risk: "Runtime identity must still be rechecked before use."
  owner: "workspace-owner"
  next_action: "Verify exact runtime identity and start a fresh isolated cluster."

- failure_id: FAIL-20260904-030
  fingerprint: "windows-postgres-child-retains-exec-session"
  related_failures: ["LIB-FAIL-1764"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: HARNESS
  phase: postgres-runtime
  capability: "POSTGRES-LEAD-INGRESS"
  scope: "Fresh PostgreSQL 18.6 isolated cluster on loopback port 55444"
  command_redacted: "Start pg_ctl directly inside the unified exec shell and wait"
  toolchain: ["PostgreSQL 18.6 windows/amd64", "PowerShell"]
  error_signature_redacted: "Server reached ready state but exec session did not return; interruption terminated server; subsequent connection refused"
  violated_invariant: "A runtime gate needs a stable independently probeable server and a completed start command."
  blast_radius: "No SQL migration or repository PASS is claimed from that attempt."
  reproduction: "Launch pg_ctl directly through the PTY-backed command wrapper."
  hypotheses_rejected: ["The ready log alone proves the server remains available after wrapper completion"]
  root_cause: "The spawned PostgreSQL process remained coupled to the command session on this host."
  canonical_correction: "Keep hidden Start-Process -Wait as the lifecycle controller, probe independently while it is live, then stop by exact data directory and collect both exit codes."
  regression_test: "Independent psql returns 18.6/checksums on/loopback; controlled stop and lifecycle controller both exit zero; no process remains."
  clean_rebuild_evidence: "PostgreSQL 18.6 served on 127.0.0.1:55444, all gates ran, pg_ctl stop passed, controller returned exit 0 and process inventory was empty."
  dependent_gates_rerun: ["PostgreSQL health", "migrations", "SQL tests", "Go repository integration"]
  upstream_reference: "Existing Windows PostgreSQL reconstruction harness pattern"
  residual_risk: "Server lifecycle and cleanup must remain bound to the exact data directory."
  owner: "workspace-owner"
  next_action: "Preserve the start/probe/stop lifecycle in pack evidence."

- failure_id: FAIL-20260904-031
  fingerprint: "lead-ingress-integration-test-loop-brace-omitted"
  related_failures: ["LIB-FAIL-1765"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: IMPLEMENTATION
  phase: postgres-integration-test-authoring
  capability: "POSTGRES-LEAD-INGRESS"
  scope: "Concurrent error collection block in lead_ingress_integration_test.go"
  command_redacted: "Inspect newly authored integration test before compilation"
  toolchain: ["PowerShell"]
  error_signature_redacted: "Missing closing brace before cardinality calculation"
  violated_invariant: "Every authored test must be syntax-checked before its result is considered."
  blast_radius: "Temporary test only; no canonical pack was updated."
  reproduction: "Read the tail of the initial test draft."
  hypotheses_rejected: ["The nested block was structurally complete"]
  root_cause: "Large single-pass test edit omitted one block terminator."
  canonical_correction: "Insert the exact terminator and run gofmt plus focal compilation immediately."
  regression_test: "gofmt accepts the file and the PostgreSQL integration test compiles/runs."
  clean_rebuild_evidence: "Pending focal runtime gate after the correction."
  dependent_gates_rerun: ["gofmt", "PostgreSQL Go integration"]
  upstream_reference: "Local failure learning contract"
  residual_risk: "Behavioral database assertions remain to run."
  owner: "workspace-owner"
  next_action: "Execute the corrected test against the fresh PostgreSQL cluster."

- failure_id: FAIL-20260904-032
  fingerprint: "append-only-integration-fixture-static-tenant-not-rerunnable"
  related_failures: ["LIB-FAIL-1766"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: TEST_DATA
  phase: postgres-integration
  capability: "POSTGRES-LEAD-INGRESS"
  scope: "Repeated execution of the new Go integration test"
  command_redacted: "Run focal test, then full suite against the same ephemeral database"
  toolchain: ["Go 1.26.7", "PostgreSQL 18.6"]
  error_signature_redacted: "duplicate tenant primary key SQLSTATE 23505"
  violated_invariant: "A focused integration test remains repeatable within its declared database lifecycle."
  blast_radius: "The global suite stopped; no PASS is claimed."
  reproduction: "Execute the same test twice without rebuilding its append-only database."
  hypotheses_rejected: ["A static tenant is safe when append-only evidence intentionally survives the test"]
  root_cause: "The fixture identity did not account for intentional append-only retention."
  canonical_correction: "Generate a unique test tenant and query all assertions within that scope."
  regression_test: "Run the focal test twice consecutively on one database."
  clean_rebuild_evidence: "The corrected focal PostgreSQL test passed twice consecutively on one migrated database and again in the clean full suite."
  dependent_gates_rerun: ["focal twice", "full PostgreSQL-backed Go suite"]
  upstream_reference: "Local append-only evidence contract"
  residual_risk: "Ephemeral databases still require lifecycle cleanup outside table mutation."
  owner: "workspace-owner"
  next_action: "Preserve per-run scoped identity and ephemeral database lifecycle."

- failure_id: FAIL-20260904-033
  fingerprint: "future-lead-event-cross-contaminates-global-outbox-worker-test"
  related_failures: ["LIB-FAIL-1767"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: TEST_DATA
  phase: postgres-integration
  capability: "CONNECTED-OUTBOX"
  scope: "Cross-test interaction between lead ingress and global outbox worker"
  command_redacted: "Run all Go tests with PostgreSQL integrations enabled"
  toolchain: ["Go 1.26.7", "PostgreSQL 18.6"]
  error_signature_redacted: "outbox_event published_at >= occurred_at constraint SQLSTATE 23514"
  violated_invariant: "A new fixture cannot poison another owner's global queue."
  blast_radius: "The full connected persistence suite is not green."
  reproduction: "Insert an available outbox event with occurred_at later than database clock, then run global Outbox.Claim."
  hypotheses_rejected: ["The existing outbox implementation violated its ordering constraint", "The constraint should be weakened"]
  root_cause: "The new deterministic fixture used a future business timestamp and exposed it to a cross-tenant global worker."
  canonical_correction: "Use an unambiguously past fixture time, rebuild the database and rerun the complete integration suite."
  regression_test: "Omnichannel test followed by all PostgreSQL-backed Go tests passes on a clean database."
  clean_rebuild_evidence: "Fresh database: 44 migrations, 30 SQL tests, all 45 Go packages with PostgreSQL enabled and 0044 down/up/test passed."
  dependent_gates_rerun: ["44 migrations", "30 SQL tests", "full Go integration suite"]
  upstream_reference: "platform.outbox_event invariant and global worker semantics"
  residual_risk: "Production ingress must separately validate excessive provider clock skew."
  owner: "workspace-owner"
  next_action: "Keep ingress fixtures in the past and mark their outbox evidence published after assertions."

- failure_id: FAIL-20260904-034
  fingerprint: "apply-patch-duplicate-update-owner"
  related_failures: ["LIB-FAIL-1768"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: HARNESS
  phase: failure-recording
  capability: "FAILURE-MEMORY"
  scope: "Atomic patch for project ledger header and entries"
  command_redacted: "Declare the same file twice in one apply_patch transaction"
  toolchain: ["apply_patch"]
  error_signature_redacted: "multiple operations target PROJECT_FAILURE_LESSONS.md"
  violated_invariant: "One patch transaction uses one update owner per file."
  blast_radius: "No file was modified by the rejected patch."
  reproduction: "Use two Update File blocks for the same absolute path."
  hypotheses_rejected: ["Separate hunks require separate Update File declarations"]
  root_cause: "Patch structure, not target content, was invalid."
  canonical_correction: "Use one Update File block with multiple hunks."
  regression_test: "The corrected atomic patch applies and all new IDs are discoverable exactly once."
  clean_rebuild_evidence: "Corrected patch applied atomically."
  dependent_gates_rerun: ["ledger uniqueness", "readiness count"]
  upstream_reference: "apply_patch contract"
  residual_risk: "None after exact ID/count verification."
  owner: "workspace-owner"
  next_action: "Verify ledger and readiness synchronization."

- failure_id: FAIL-20260904-035
  fingerprint: "single-pack-roundtrip-assumes-manifest-and-emits-false-pass"
  related_failures: ["LIB-FAIL-1769"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: EVIDENCE
  phase: pack-roundtrip
  capability: "GO-OMNICHANNEL-LEAD-INGRESS"
  scope: "Ten-file materialization output verification"
  command_redacted: "Read an assumed materializer manifest and continue after non-terminating PowerShell errors"
  toolchain: ["PowerShell", "materialize_markdown_pack.ps1"]
  error_signature_redacted: "MATERIALIZATION_MANIFEST.json missing; null hash operations; false ROUNDTRIP_HASH_PASS text"
  violated_invariant: "A success marker is valid only if every prerequisite assertion executed and errors terminate the gate."
  blast_radius: "No round-trip hash PASS may be claimed yet; canonical pack status remains conditioned."
  reproduction: "Expect compositor metadata from the single-pack materializer without setting ErrorActionPreference Stop."
  hypotheses_rejected: ["The individual materializer emits a JSON manifest", "A final success string overrides preceding errors"]
  root_cause: "Two tool contracts were conflated and the verifier was not fail-closed."
  canonical_correction: "Use a fixed ten-path list, require each source/destination file, compare SHA-256 under terminating errors and assert exact destination cardinality."
  regression_test: "All ten pairs exist, all hashes match and destination contains exactly ten files."
  clean_rebuild_evidence: "Fresh materialization produced exactly ten files and a fail-closed source/destination SHA-256 comparison matched all ten."
  dependent_gates_rerun: ["pack round-trip", "pack admission"]
  upstream_reference: "PACK_CONTRACT.md and materialize_markdown_pack.ps1"
  residual_risk: "Metadata/provenance catalog gates remain separate."
  owner: "workspace-owner"
  next_action: "Keep direct cardinality and hash comparison for single-pack materialization."

- failure_id: FAIL-20260904-036
  fingerprint: "apply-patch-duplicate-update-owner-recurrence"
  related_failures: ["LIB-FAIL-1768", "LIB-FAIL-1770"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 1
  status: RECURRENCE_PROVEN
  severity: LOW
  classification: HARNESS
  phase: failure-recording
  capability: "FAILURE-MEMORY"
  scope: "Patch used to register the round-trip verifier failure"
  command_redacted: "Repeat two Update File declarations for the same project ledger"
  toolchain: ["apply_patch"]
  error_signature_redacted: "multiple operations target PROJECT_FAILURE_LESSONS.md"
  violated_invariant: "A known patch construction lesson must be applied on the next edit."
  blast_radius: "No files were changed by the rejected transaction."
  reproduction: "Repeat the structure from FAIL-20260904-034."
  hypotheses_rejected: ["The earlier rejection was transient"]
  root_cause: "The repeated patch was assembled from the prior draft without consolidating both hunks."
  canonical_correction: "Use exactly one Update File declaration per target file."
  regression_test: "Corrected patch applies and both IDs appear exactly once in their owner ledgers."
  clean_rebuild_evidence: "This corrected patch applied atomically with one owner declaration per file."
  dependent_gates_rerun: ["ledger uniqueness", "readiness count"]
  upstream_reference: "apply_patch contract"
  residual_risk: "None after exact count verification."
  owner: "workspace-owner"
  next_action: "Do not reuse duplicate-owner patch templates."

- failure_id: FAIL-20260904-037
  fingerprint: "google-webhook-verification-key-persisted-with-byte-exact-raw"
  related_failures: ["LIB-FAIL-1771"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: CRITICAL
  classification: SECURITY
  phase: pack-admission
  capability: "GO-OMNICHANNEL-LEAD-INGRESS"
  scope: "Google webhook raw payload persistence"
  command_redacted: "Review provider schema against immutable raw persistence before catalog promotion"
  toolchain: ["Go 1.26.7", "PostgreSQL 18.6", "official Google webhook contract"]
  error_signature_redacted: "google_key remains inside persisted payload bytea"
  violated_invariant: "Provider credentials are verified but never persisted in cleartext evidence."
  blast_radius: "A database reader or backup exposure could reveal the webhook verification credential."
  reproduction: "Decode valid payload and compare persisted RawEvent.Payload with its google_key field."
  hypotheses_rejected: ["Raw fidelity overrides secret minimization", "Database-at-rest encryption alone removes cleartext exposure to authorized readers"]
  root_cause: "Source evidence and safe durable representation were modeled as one byte sequence."
  canonical_correction: "Split source-body SHA-256 from stored-payload SHA-256; authenticate against source, remove google_key before persistence and test source divergence plus secret absence."
  regression_test: "Persisted bytes and outbox contain no secret; source hash detects changed original; stored hash validates redacted bytes; all Go/PostgreSQL gates rerun."
  clean_rebuild_evidence: "Google/OWASP-governed redaction plus source/stored hashes passed HTTP tests, PostgreSQL 44 migrations/30 SQL tests, 45 Go packages, vet/build, down/up and ten-file round-trip."
  dependent_gates_rerun: ["HTTP tests", "PostgreSQL integration", "44 migrations", "30 SQL tests", "45-package Go suite", "pack round-trip"]
  upstream_reference: "Google Lead Form Webhook implementation; OWASP Logging and Secrets Management guidance"
  residual_risk: "PII still requires approved retention, access control and encryption at the project target."
  owner: "workspace-owner"
  next_action: "Retain live secret rotation and target PII controls as project gates."

- failure_id: FAIL-20260904-038
  fingerprint: "secret-redaction-authority-web-batch-truncated"
  related_failures: ["LIB-FAIL-1772"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: EVIDENCE
  phase: authority-research
  capability: "GO-OMNICHANNEL-LEAD-INGRESS"
  scope: "Official Google/OWASP authority lookup for secret-safe evidence"
  command_redacted: "Search three broad official topics with long response"
  toolchain: ["web search"]
  error_signature_redacted: "Output exceeded model context and was truncated"
  violated_invariant: "Only complete source reads govern security claims."
  blast_radius: "No final citation or implementation claim uses unseen content."
  reproduction: "Request long aggregated results for multiple extensive security pages."
  hypotheses_rejected: ["Search aggregation is equivalent to exact opened pages"]
  root_cause: "Research scope was wider than needed for one invariant."
  canonical_correction: "Open exact official pages separately and retain only their narrow claims."
  regression_test: "Each used security assertion maps to one completely opened official page."
  clean_rebuild_evidence: "Focused complete opens captured Google lines 42-49, 71-116, 188-213 and OWASP lines 354-380."
  dependent_gates_rerun: ["authority citation check", "pack provenance check"]
  upstream_reference: "Official Google Ads and OWASP pages"
  residual_risk: "Web content freshness remains date-bound."
  owner: "workspace-owner"
  next_action: "Use only those narrow dated claims in pack/evidence."

- failure_id: FAIL-20260904-039
  fingerprint: "windows-ripgrep-literal-0044-glob-recurrence"
  related_failures: ["LIB-FAIL-1761", "LIB-FAIL-1773"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 1
  status: RECURRENCE_PROVEN
  severity: LOW
  classification: DISCOVERY
  phase: secret-redaction-refactor
  capability: "GO-OMNICHANNEL-LEAD-INGRESS"
  scope: "Search across Go and migration 0044 sources"
  command_redacted: "Pass 0044* wildcard paths to rg on Windows"
  toolchain: ["PowerShell", "ripgrep"]
  error_signature_redacted: "Filename, directory name, or volume label syntax is incorrect"
  violated_invariant: "Known Windows searches use -g file filters."
  blast_radius: "SQL reference completeness remains unproven by that call."
  reproduction: "Use wildcard path operands rather than -g."
  hypotheses_rejected: ["The wildcard path will be expanded automatically"]
  root_cause: "The search combined exact directories and wildcard operands without normalizing the selection method."
  canonical_correction: "Search directory roots with -g '0044*' and terminating errors."
  regression_test: "A corrected search covers both Go and SQL references with exit zero."
  clean_rebuild_evidence: "Pending corrected reference sweep after implementation changes."
  dependent_gates_rerun: ["reference sweep", "compile", "SQL migration"]
  upstream_reference: "Local failure learning contract"
  residual_risk: "None after complete corrected sweep."
  owner: "workspace-owner"
  next_action: "Use only valid -g filters for the remaining refactor."

- failure_id: FAIL-20260904-040
  fingerprint: "redaction-refactor-empty-patch-hunk"
  related_failures: ["LIB-FAIL-1774"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: HARNESS
  phase: secret-redaction-refactor
  capability: "GO-OMNICHANNEL-LEAD-INGRESS"
  scope: "Multi-file Go and SQL patch"
  command_redacted: "Submit an update marker after an empty hunk"
  toolchain: ["apply_patch"]
  error_signature_redacted: "Unexpected Update File line inside update hunk"
  violated_invariant: "Every patch hunk has valid context or changed lines."
  blast_radius: "Patch was rejected atomically; no source was partially changed."
  reproduction: "Place @@ immediately before the next file marker."
  hypotheses_rejected: ["An empty hunk can delimit files"]
  root_cause: "Patch syntax contained an unnecessary empty hunk."
  canonical_correction: "Use one valid context-bearing hunk per changed region."
  regression_test: "Corrected patch applies and subsequent complete reference sweep sees all intended fields."
  clean_rebuild_evidence: "Corrected patch applied after removing the empty hunk."
  dependent_gates_rerun: ["gofmt", "Go tests", "SQL migration"]
  upstream_reference: "apply_patch contract"
  residual_risk: "Behavioral refactor still requires all gates."
  owner: "workspace-owner"
  next_action: "Execute the corrected refactor and verify every changed owner."

- failure_id: FAIL-20260904-041
  fingerprint: "v224-documentation-invalid-javascript-unicode-escape"
  related_failures: ["LIB-FAIL-1775"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: HARNESS
  phase: v224-documentation
  capability: "GO-OMNICHANNEL-LEAD-INGRESS"
  scope: "Evidence and catalog patch wrapper"
  command_redacted: "Embed a misspelled JavaScript Unicode escape in patch text"
  toolchain: ["functions.exec", "apply_patch"]
  error_signature_redacted: "SyntaxError: Invalid Unicode escape sequence"
  violated_invariant: "Tool wrapper source must parse before patch application."
  blast_radius: "No documentation or catalog file changed."
  reproduction: "Use a non-hex character inside a JavaScript \\u escape."
  hypotheses_rejected: ["The wrapper will pass malformed escapes through literally"]
  root_cause: "A Spanish word was encoded with an invalid manual escape."
  canonical_correction: "Use valid Unicode text and keep only path backslashes escaped."
  regression_test: "Corrected patch applies and the evidence/catalog files contain the intended accented text."
  clean_rebuild_evidence: "Corrected patch applied after wrapper syntax was fixed."
  dependent_gates_rerun: ["evidence link check", "catalog gate"]
  upstream_reference: "JavaScript string literal syntax"
  residual_risk: "None after exact file reads."
  owner: "workspace-owner"
  next_action: "Verify V224 evidence and catalog entries."

- failure_id: FAIL-20260904-042
  fingerprint: "post-deepseek-reaudit-leaks-machine-local-home-path"
  related_failures: ["LIB-FAIL-1776"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: PRIVACY_PORTABILITY
  phase: root-verification
  capability: "LIBRARY-RELEASE"
  scope: "Post-DeepSeek re-audit reference to the user-authorized NEW BINANCE project"
  command_redacted: "Run VERIFY_LIBRARY.ps1"
  toolchain: ["PowerShell", "VERIFY_LIBRARY.ps1"]
  error_signature_redacted: "machine-local home path in release file"
  violated_invariant: "Distributed Markdown cannot expose a machine-local user path."
  blast_radius: "Root verification stopped; no release PASS was claimed."
  reproduction: "Search the re-audit file for the local user profile prefix."
  hypotheses_rejected: ["An authorized read path is appropriate release metadata"]
  root_cause: "Runtime access location was confused with portable provenance identity."
  canonical_correction: "Refer only to NEW BINANCE/Binance and retain its classification as user-owned, not Binance-authored."
  regression_test: "Root path scan finds no machine-local home and verification advances beyond this assertion."
  clean_rebuild_evidence: "The absolute path was removed before the next root gate."
  dependent_gates_rerun: ["VERIFY_LIBRARY", "release path scan"]
  upstream_reference: "Library release portability contract"
  residual_risk: "Other newly added release files must pass the same scan."
  owner: "workspace-owner"
  next_action: "Repeat root verification and address only newly observed failures."
- failure_id: FAIL-20260904-043
  fingerprint: "v224-inventory-patch-context-mismatch"
  related_failures: ["LIB-FAIL-1777"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  phase: inventory-synchronization
  capability: "LIBRARY-RELEASE"
  scope: "V224 canonical inventory consumers"
  command_redacted: "Apply the multi-file inventory patch"
  toolchain: ["apply_patch", "PowerShell"]
  error_signature_redacted: "Failed to find expected readiness line"
  violated_invariant: "A documentation update must use exact observed context and must not infer a partial replacement succeeded."
  blast_radius: "No file changed because the patch was rejected atomically."
  reproduction: "Attempt the former full-line hunk against the current readiness paragraph."
  hypotheses_rejected: ["The readiness sentence ended immediately after the inventory counts"]
  root_cause: "The patch expected only the prefix of a longer physical line."
  canonical_correction: "Record the failure, replace only exact count tokens and verify all consumers after the edit."
  regression_test: "All inventory consumers report the parser-derived V224 totals and the root count assertion advances."
  clean_rebuild_evidence: "Atomic rejection preserved the prior state; corrected edits are followed by exact searches and parser recount."
  dependent_gates_rerun: ["inventory parser", "VERIFY_LIBRARY"]
  upstream_reference: "Local verifier Read-MaterializationProvenance"
  residual_risk: "The two known critical OPEN findings still intentionally block root PASS."
  owner: "workspace-owner"
  next_action: "Synchronize V224 counts and continue fixing the connected runtime."

- failure_id: FAIL-20260904-044
  fingerprint: "v224-ledger-patch-backslash-anchor-recurrence"
  related_failures: ["LIB-FAIL-1778", "FAIL-20260904-043"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  phase: failure-registration
  capability: "LIBRARY-RELEASE"
  scope: "Canonical failure ledgers"
  command_redacted: "Apply a ledger patch anchored on a Windows path"
  toolchain: ["apply_patch"]
  error_signature_redacted: "Failed to find expected LIB-FAIL-1776 row"
  violated_invariant: "Failure records must be inserted with robust exact context."
  blast_radius: "No file changed because the patch was rejected atomically."
  reproduction: "Use an escaped Windows path as an exact patch anchor."
  hypotheses_rejected: ["Rendered and literal backslash representations are interchangeable patch context"]
  root_cause: "A path-sensitive row was selected instead of a stable section heading."
  canonical_correction: "Insert before stable section headings and verify unique IDs."
  regression_test: "Both ledgers contain unique 1777/1778 and project 043/044 entries."
  clean_rebuild_evidence: "Atomic rejection prevented partial state; corrected insertion is structurally anchored."
  dependent_gates_rerun: ["failure ID uniqueness", "VERIFY_LIBRARY"]
  upstream_reference: "Local failure-learning contract"
  residual_risk: "None for the rejected patch after structurally anchored insertion."
  owner: "workspace-owner"
  next_action: "Continue with exact count replacements."
- failure_id: FAIL-20260904-045
  fingerprint: "v224-partial-line-patch-recurrence"
  related_failures: ["LIB-FAIL-1779", "FAIL-20260904-043", "FAIL-20260904-044"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 2
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  phase: inventory-synchronization
  capability: "LIBRARY-RELEASE"
  scope: "V224 inventory lines"
  command_redacted: "Apply substring-style hunks to long Markdown lines"
  toolchain: ["apply_patch"]
  error_signature_redacted: "Failed to find expected partial line"
  violated_invariant: "Patch deletions must match complete physical lines."
  blast_radius: "No file changed because the patch was rejected atomically."
  reproduction: "Use only the beginning of the THIRD_PARTY_NOTICES inventory sentence as a removed line."
  hypotheses_rejected: ["apply_patch treats a removed line as a substring"]
  root_cause: "The hunk was not built from the complete observed line."
  canonical_correction: "Use complete line replacements captured directly from each file."
  regression_test: "Exact search reports only V224 totals across all inventory consumers."
  clean_rebuild_evidence: "Atomic rejection preserved state before the exact-line rewrite."
  dependent_gates_rerun: ["inventory parser", "VERIFY_LIBRARY"]
  upstream_reference: "apply_patch line-oriented contract"
  residual_risk: "None once exact lines and counts are verified."
  owner: "workspace-owner"
  next_action: "Apply exact physical-line replacements."
- failure_id: FAIL-20260904-046
  fingerprint: "openai-responses-doc-fetch-timeout"
  related_failures: ["LIB-FAIL-1780"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: UPSTREAM_AVAILABILITY
  phase: official-source-research
  capability: "AGENT-RUNTIME-WIRING"
  scope: "OpenAI Responses/function-calling authority"
  command_redacted: "Open official Responses create reference"
  toolchain: ["OpenAI Docs", "web"]
  error_signature_redacted: "400 timeout fetching official page"
  violated_invariant: "An unopened search snippet cannot solely govern an implementation contract."
  blast_radius: "No adapter change was made from the failed fetch."
  reproduction: "Open the timed-out official reference under transient conditions."
  hypotheses_rejected: ["Search-result text alone is adequate implementation evidence"]
  root_cause: "Transient official-site fetch timeout."
  canonical_correction: "Open a different official guide that establishes the required API/tool contract."
  regression_test: "Official GPT-4.1 guide is fetched and lines 843-880 establish API-native tool use and evals."
  clean_rebuild_evidence: "No build affected; source selection recovered before implementation."
  dependent_gates_rerun: ["adapter contract review"]
  upstream_reference: "https://developers.openai.com/api/docs/guides/latest-model?model=gpt-4.1"
  residual_risk: "A live provider contract test is still project-conditioned."
  owner: "workspace-owner"
  next_action: "Use only the opened official guide for the narrow claim."

- failure_id: FAIL-20260904-047
  fingerprint: "approval-owner-filename-assumed"
  related_failures: ["LIB-FAIL-1781"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: DISCOVERY
  phase: runtime-wiring-audit
  capability: "HUMAN-APPROVAL"
  scope: "Materialized approval package"
  command_redacted: "Read assumed model.go"
  toolchain: ["PowerShell"]
  error_signature_redacted: "approval/model.go not found"
  violated_invariant: "Paths must come from observed inventory."
  blast_radius: "Only one read failed; no source changed."
  reproduction: "Request the nonexistent model.go path."
  hypotheses_rejected: ["Approval model uses the conventional model.go filename"]
  root_cause: "The package uses approval.go for its model."
  canonical_correction: "Enumerate the directory and read approval.go plus registry.go."
  regression_test: "Both actual owners are read completely before design."
  clean_rebuild_evidence: "Read-only recovery; build state unchanged."
  dependent_gates_rerun: ["bounded owner audit"]
  upstream_reference: "Local materialized source tree"
  residual_risk: "Registry remains in-memory and cannot prove durable production approval."
  owner: "workspace-owner"
  next_action: "Keep durable approval as an explicit connected-runtime condition."

- failure_id: FAIL-20260904-048
  fingerprint: "crm-directory-path-assumed"
  related_failures: ["LIB-FAIL-1782"]
  first_seen: "2026-09-04"
  last_seen: "2026-09-04"
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: DISCOVERY
  phase: lead-to-domain-audit
  capability: "CRM-JOURNEY"
  scope: "Materialized domain owners"
  command_redacted: "Search an assumed internal/crm directory plus observed platform paths"
  toolchain: ["ripgrep"]
  error_signature_redacted: "internal/crm path not found"
  violated_invariant: "Conceptual capability names cannot be assumed to be physical directories."
  blast_radius: "The missing directory warning prevents claiming exhaustive CRM coverage from that query."
  reproduction: "Run rg against internal/crm in the current composition."
  hypotheses_rejected: ["CRM owns one top-level package"]
  root_cause: "CRM behavior is split across electromobility, franchisejourney and platform adapters."
  canonical_correction: "Use exact symbol matches and observed paths; do not infer a single CRM owner."
  regression_test: "Exact searches identify lead capture, consent, appointments and quotations without nonexistent paths."
  clean_rebuild_evidence: "Read-only correction; no build affected."
  dependent_gates_rerun: ["connected journey owner map"]
  upstream_reference: "Local materialized source tree"
  residual_risk: "Provider candidate promotion to the existing consent-aware CRM remains unimplemented."
  owner: "workspace-owner"
  next_action: "Design promotion against observed owners and idempotency contracts."
- failure_id: FAIL-20260904-049
  fingerprint: "v225-rg-literal-glob-recurrence"
  related_failures: ["LIB-FAIL-1783"]
  recurrence_count: 3
  status: REGRESSION_PROVEN
  severity: LOW
  classification: DISCOVERY
  violated_invariant: "Use observed paths or rg -g filters on Windows."
  root_cause: "Literal wildcard paths were reused in a combined query."
  canonical_correction: "Open exact 0001 and 0003 files."
  regression_test: "Schema definitions were read from exact paths."
  clean_rebuild_evidence: "Read-only correction; no build affected."
  residual_risk: "None for the invalid query."

- failure_id: FAIL-20260904-050
  fingerprint: "v225-go-toolchain-sibling-path"
  related_failures: ["LIB-FAIL-1784"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  violated_invariant: "Use the exact observed toolchain path."
  root_cause: "The sibling toolchain was addressed as a child of franchise."
  canonical_correction: "Use elite_reaudit_20260904_v1/toolchain/go/bin/go.exe."
  regression_test: "Go 1.26.7 focal/full tests, vet and build pass."
  clean_rebuild_evidence: "Full materialized composition gates passed."
  residual_risk: "Race detector remains externally blocked by CGO/GCC."

- failure_id: FAIL-20260904-051
  fingerprint: "v225-undefined-postgres-uuid-helper"
  related_failures: ["LIB-FAIL-1785"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: IMPLEMENTATION
  violated_invariant: "Every referenced owner must exist and compile."
  root_cause: "A stableEventID helper was assumed in postgres."
  canonical_correction: "Use the existing tested leadstream.StableUUID owner."
  regression_test: "Full Go test/vet/build pass."
  clean_rebuild_evidence: "46-package suite with PostgreSQL real passed."
  residual_risk: "None for UUID linkage."

- failure_id: FAIL-20260904-052
  fingerprint: "v225-postgres-wrapper-no-return"
  related_failures: ["LIB-FAIL-1786"]
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: LOW
  classification: ENVIRONMENT
  violated_invariant: "Silence or timeout is not readiness evidence."
  root_cause: "The Windows wrapper retained execution while postmaster started."
  canonical_correction: "Inspect PID, pidfile, socket and log independently."
  regression_test: "Postmaster identity and ready log were observed."
  clean_rebuild_evidence: "PostgreSQL gates subsequently passed."
  residual_risk: "Lifecycle wrapper should be replaced before automation reuse."

- failure_id: FAIL-20260904-053
  fingerprint: "v225-postgres-port-assumed"
  related_failures: ["LIB-FAIL-1787"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: ENVIRONMENT
  violated_invariant: "Connect only to the observed server endpoint."
  root_cause: "The prior 55444 port was reused although this postmaster reported 5432."
  canonical_correction: "Read postmaster.pid and use 127.0.0.1:5432."
  regression_test: "Database probe and all V225 PostgreSQL gates pass."
  clean_rebuild_evidence: "0045 up/down/up and 31 SQL tests passed."
  residual_risk: "Project automation must pin and probe its own port."

- failure_id: FAIL-20260904-054
  fingerprint: "v225-duplicate-postmaster-start"
  related_failures: ["LIB-FAIL-1788"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: ENVIRONMENT
  violated_invariant: "Never start a second postmaster on an active data directory."
  root_cause: "The existing wrapper-started PID was not checked first."
  canonical_correction: "Probe active PID/socket/log before start."
  regression_test: "Exactly one postmaster cluster served subsequent tests."
  clean_rebuild_evidence: "No new cluster or data directory was created."
  residual_risk: "None after controlled stop."

- failure_id: FAIL-20260904-055
  fingerprint: "v225-go-test-wrong-workdir"
  related_failures: ["LIB-FAIL-1789"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  violated_invariant: "Go gates run from the module root."
  root_cause: "SQL and Go commands shared the library workdir."
  canonical_correction: "Repeat from the materialized franchise root."
  regression_test: "Focal and full tests pass with the real database URL."
  clean_rebuild_evidence: "Full suite exit 0."
  residual_risk: "None for the invalid invocation."

- failure_id: FAIL-20260904-056
  fingerprint: "v225-materializer-exit-cuts-roundtrip"
  related_failures: ["LIB-FAIL-1790"]
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: EVIDENCE
  violated_invariant: "Materialization alone is not round-trip hash evidence."
  root_cause: "Direct script invocation propagated its exit to the caller."
  canonical_correction: "Run the materializer in a child pwsh and compare every file."
  regression_test: "Cardinality 7 and 7/7 SHA-256 equality."
  clean_rebuild_evidence: "Fresh destination reports ROUNDTRIP_HASH_PASS=7."
  residual_risk: "None for V225 round-trip."

- failure_id: FAIL-20260904-057
  fingerprint: "v225-start-process-space-splitting"
  related_failures: ["LIB-FAIL-1791"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  violated_invariant: "Paths with spaces must remain one argument."
  root_cause: "Start-Process ArgumentList flattened the script path."
  canonical_correction: "Invoke pwsh directly with PowerShell argument binding."
  regression_test: "Fresh materialization and 7/7 hash comparison pass."
  clean_rebuild_evidence: "Child process exit 0; parent continued checks."
  residual_risk: "None for the corrected invocation."

- failure_id: FAIL-20260904-058
  fingerprint: "v225-failure-ledger-local-path-leak"
  related_failures: ["LIB-FAIL-1792", "FAIL-20260904-042"]
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: PRIVACY_PORTABILITY
  violated_invariant: "Distributed evidence cannot contain a machine-local home path."
  root_cause: "The exact quoting error was copied into the canonical ledger."
  canonical_correction: "Describe the truncated path portably and remove the personal prefix."
  regression_test: "Root release path scan advances beyond the new ledger row."
  clean_rebuild_evidence: "VERIFY_LIBRARY blocked before PASS; corrected source is rescanned."
  residual_risk: "New evidence must continue through the same release scan."

- failure_id: FAIL-20260904-059
  fingerprint: "meta-leads-project-source-profile-assumed-at-library-root"
  related_failures: ["LIB-FAIL-1793"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: PROVENANCE_GOVERNANCE
  violated_invariant: "A project-scoped source profile must not be assumed to exist at the library-maintenance root."
  root_cause: "The inspection requested a per-project generated record before resolving whether the current root was a project or the canonical library."
  canonical_correction: "Do not acquire a new upstream group here; reuse only already-admitted exact Meta artifacts and require the generated record in each consuming project."
  regression_test: "No network acquisition occurs and the new pack must reference the existing exact SDK artifact lock."
  clean_rebuild_evidence: "The failed lookup was read-only and occurred before any acquisition."
  residual_risk: "A consuming project remains blocked until its own source profile and live provider access are validated."

- failure_id: FAIL-20260904-060
  fingerprint: "pack-tooling-root-scripts-directory-assumed"
  related_failures: ["LIB-FAIL-1794"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING_DISCOVERY
  violated_invariant: "Physical paths must be observed before targeted access."
  root_cause: "A conventional root scripts directory was assumed without inventory evidence."
  canonical_correction: "Resolve materializer/updater paths from the repository file inventory."
  regression_test: "The next lookup starts from rg --files and opens only observed paths."
  clean_rebuild_evidence: "The failed lookup was read-only."
  residual_risk: "None beyond continued path discipline."

- failure_id: FAIL-20260904-061
  fingerprint: "official-source-profile-registry-filename-assumed"
  related_failures: ["LIB-FAIL-1795"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING_DISCOVERY
  violated_invariant: "Canonical source records must be resolved from observed inventory."
  root_cause: "A conceptual registry name was used as a physical path without prior evidence."
  canonical_correction: "Use the observed master map and exact embedded SDK artifact lock."
  regression_test: "All subsequent source reads target paths returned by rg --files."
  clean_rebuild_evidence: "The failed lookup was read-only; valid file matches remain independently usable."
  residual_risk: "None beyond continued path discipline."

- failure_id: FAIL-20260904-062
  fingerprint: "v226-meta-field-cardinality-incompatible-with-durable-owner"
  related_failures: ["LIB-FAIL-1796"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: CONTRACT_INTEGRATION
  violated_invariant: "A candidate artifact must be accepted without semantic rewriting by its declared downstream owner."
  root_cause: "The provider array cardinality was preserved in the normalized candidate even though the central candidate contract is one value per field."
  canonical_correction: "Keep all values in raw evidence; accept exactly one non-empty value for normalized candidates and reject ambiguity."
  regression_test: "V226 unique-value and multi-value tests plus provider-neutral importer tests and PostgreSQL integration."
  clean_rebuild_evidence: "Detected before canonical downstream importer or production claim."
  residual_risk: "Forms requiring true multi-select semantics need an explicitly versioned mapping contract."

- failure_id: FAIL-20260904-063
  fingerprint: "v224-migration-suffix-assumed-before-prefix-enumeration"
  related_failures: ["LIB-FAIL-1797", "FAIL-20260904-026"]
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING_DISCOVERY
  violated_invariant: "Exact SQL paths must be observed before access."
  root_cause: "The conceptual short name was substituted for the canonical migration suffix."
  canonical_correction: "Enumerate the exact 0044 prefix, then read only returned paths."
  regression_test: "The next SQL inspection uses 0044_omnichannel_lead_ingress.* exactly."
  clean_rebuild_evidence: "Both failed reads were non-mutating."
  residual_risk: "None beyond enforcing observed paths."

- failure_id: FAIL-20260904-064
  fingerprint: "v226-rejected-outcome-lost-tenant-organization-scope"
  related_failures: ["LIB-FAIL-1798"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: CRITICAL
  classification: TENANT_ISOLATION
  violated_invariant: "Rejected evidence must remain tenant/organization scoped even when no normalized candidate exists."
  root_cause: "Scope was nested only under the nullable candidate object."
  canonical_correction: "Add immutable batch-level scope, cross-check normalized candidates, and persist rejected raw using the envelope scope."
  regression_test: "Python rejection fixtures and Go memory/PostgreSQL tests prove scoped rejected evidence."
  clean_rebuild_evidence: "Detected before importer compilation or release claim."
  residual_risk: "A project must still validate its tenant/form routing configuration."

- failure_id: FAIL-20260904-065
  fingerprint: "v227-go-fmt-files-cross-package"
  related_failures: ["LIB-FAIL-1799"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TEST_HARNESS
  violated_invariant: "Formatting, tests, vet and build require independent successful exits."
  root_cause: "Named files from two Go directories were passed to one go fmt invocation and native failure was not made terminating."
  canonical_correction: "Format each package separately and check LASTEXITCODE after every native tool call."
  regression_test: "Two format commands plus focal/full test, vet and build all exit zero."
  clean_rebuild_evidence: "The failed format did not write across packages; no promotion occurred."
  residual_risk: "None after independent rerun."

- failure_id: FAIL-20260904-066
  fingerprint: "v226-powershell-updater-stale-lastexitcode"
  related_failures: ["LIB-FAIL-1800"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TEST_HARNESS
  violated_invariant: "A successful PowerShell script must not be invalidated by stale native-process state."
  root_cause: "LASTEXITCODE was read after a PowerShell script instead of relying on its terminating-error contract."
  canonical_correction: "Use exception propagation for PowerShell scripts and LASTEXITCODE only immediately after native commands."
  regression_test: "Fresh V226 materialization and tests proceed after the successful updater output."
  clean_rebuild_evidence: "Updater reported eight-block source remained valid and three updated blocks were written; no materialization claim was made."
  residual_risk: "None after the fresh round-trip."

- failure_id: FAIL-20260904-067
  fingerprint: "v227-postgres-test-sha-helper-import-and-dead-call"
  related_failures: ["LIB-FAIL-1801"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: IMPLEMENTATION
  violated_invariant: "New integration fixtures must compile from explicit imports without irrelevant computation."
  root_cause: "The compact SHA helper omitted its package import and retained a discarded StableUUID call."
  canonical_correction: "Use crypto/sha256 directly and remove the dead call."
  regression_test: "gofmt, focal PostgreSQL integration, full Go test, vet and build."
  clean_rebuild_evidence: "Detected before the first compiler invocation."
  residual_risk: "None after gates pass."

- failure_id: FAIL-20260904-068
  fingerprint: "v225-canonical-integration-test-not-gofmt-normalized"
  related_failures: ["LIB-FAIL-1802"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: RELEASE_REPRODUCIBILITY
  violated_invariant: "Canonical materialized Go sources must already be gofmt-stable."
  root_cause: "The V225 focal formatting step omitted the PostgreSQL integration test before embedding."
  canonical_correction: "Re-embed the gofmt output and rebuild the exact V225 pack."
  regression_test: "Fresh 7/7 materialization is unchanged by gofmt and passes focal/full Go, vet and build."
  clean_rebuild_evidence: "A byte-level diff proved formatting-only drift before correction."
  residual_risk: "All future Go packs must test gofmt idempotence over every materialized .go file."

- failure_id: FAIL-20260904-069
  fingerprint: "legacy-openai-adapter-two-files-not-gofmt-stable"
  related_failures: ["LIB-FAIL-1803", "FAIL-20260904-068"]
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: RELEASE_REPRODUCIBILITY
  violated_invariant: "Every canonical Go materialization must be gofmt-stable."
  root_cause: "The historical LLM pack embedded compact production files without a full-pack formatting idempotence gate."
  canonical_correction: "Re-embed only adapter.go/config.go after gofmt and reconstruct the exact three-file pack."
  regression_test: "Fresh 3/3 round-trip, gofmt no-op, package/full tests, vet and build."
  clean_rebuild_evidence: "SHA comparison isolated the two files before canonical update."
  residual_risk: "Continue whole-pack gofmt idempotence checks for Go packs touched by this re-audit."

- failure_id: FAIL-20260904-070
  fingerprint: "v228-promotion-ai-map-heading-assumed"
  related_failures: ["LIB-FAIL-1804"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: DOCUMENTATION_MAINTENANCE
  violated_invariant: "Multi-file patches require exact observed anchors for every target."
  root_cause: "An unobserved AI map heading was added to an otherwise valid promotion transaction."
  canonical_correction: "Split promotion edits and inspect exact headings before the AI map insertion."
  regression_test: "All intended V228 evidence/count/map rows are present exactly once."
  clean_rebuild_evidence: "The failed patch was rejected atomically."
  residual_risk: "None after narrow patches and inventory verification."

- failure_id: FAIL-20260904-071
  fingerprint: "appointment-binding-idempotency-header-omitted"
  related_failures: ["LIB-FAIL-1805"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: DATA_INTEGRITY
  violated_invariant: "Every retryable side effect must receive a stable durable idempotency identity."
  root_cause: "The domain adapter and fake test covered path/body but omitted the header consumed by the real appointment endpoint."
  canonical_correction: "Require and propagate a provider-message-derived Idempotency-Key and prove the header at the adapter boundary."
  regression_test: "Same key/same request is one durable appointment; malformed/missing keys fail closed; different payload conflicts."
  clean_rebuild_evidence: "Cross-owner inspection found the omission before connected runtime promotion."
  residual_risk: "The legacy method must not remain reachable from the connected runtime without a stable key."

- failure_id: FAIL-20260904-072
  fingerprint: "quote-creation-no-durable-idempotency"
  related_failures: ["LIB-FAIL-1806"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: CRITICAL
  classification: DATA_INTEGRITY
  violated_invariant: "At-least-once channel delivery may create at most one quotation and one logical outbox event."
  root_cause: "The quote handler/service/repository generate fresh identifiers on every invocation and own no request key/hash."
  canonical_correction: "Implement tenant-scoped durable idempotency through HTTP, service and PostgreSQL transaction before the quote tool is enabled."
  regression_test: "Same key/same hash replays the original result; same key/different hash conflicts; concurrent duplicates persist one quote/event."
  clean_rebuild_evidence: "V229: PostgreSQL 18.6 focal twice plus 47-package suite, vet/build and 33/33 materialization PASS."
  residual_risk: "The connected caller must still supply a stable ingress identity; legacy calls fail closed."

- failure_id: FAIL-20260904-073
  fingerprint: "sql-discovery-windows-literal-glob-recurrence"
  related_failures: ["LIB-FAIL-1807", "LIB-FAIL-1761", "LIB-FAIL-1773", "LIB-FAIL-1797"]
  recurrence_count: 3
  status: RECURRENCE_PROVEN
  severity: LOW
  classification: TOOLING
  violated_invariant: "Windows file discovery uses observed directories and rg filters, never literal wildcard paths."
  root_cause: "A combined owner query reintroduced literal SQL glob arguments."
  canonical_correction: "Search the exact directories with -g filters and invalidate the rejected branches."
  regression_test: "Subsequent SQL discovery returns without path errors."
  clean_rebuild_evidence: "No SQL conclusion used the rejected output."
  residual_risk: "Future combined commands must preserve platform-safe discovery."

- failure_id: FAIL-20260904-074
  fingerprint: "domain-gateway-protected-routes-have-no-service-identity"
  related_failures: ["LIB-FAIL-1808"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: CRITICAL
  classification: SECURITY_INTEGRATION
  violated_invariant: "A connected agent may reach protected domain actions only through an authenticated, authorized and tenant-scoped principal."
  root_cause: "The HTTP gateway has no token source and its fake server did not exercise the backend verifier."
  canonical_correction: "Add an injectable short-lived credential provider, Bearer propagation for protected calls, fail-closed behavior and real verifier tests."
  regression_test: "Protected calls reject absent credentials, carry the current token, never persist/log it and rotate without reconstructing the gateway."
  clean_rebuild_evidence: "V230: HTTP tests prove Bearer presence/rotation, absent token and cross-tenant rejection; 3/3 materialization and full Go/vet/build PASS."
  residual_risk: "Each project must provide and verify its live IdP/token acquisition; the library does not invent credentials."

- failure_id: FAIL-20260904-075
  fingerprint: "domainbind-refactor-future-anchor-assumed"
  related_failures: ["LIB-FAIL-1809"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: IMPLEMENTATION_MAINTENANCE
  violated_invariant: "Patches anchor only on text already observed in the target."
  root_cause: "A new var block was accidentally used as context before it existed."
  canonical_correction: "Split the refactor into narrow patches around existing imports, declarations and functions."
  regression_test: "Each narrow patch applies and the package is gofmt-stable."
  clean_rebuild_evidence: "The failed transaction was rejected atomically."
  residual_risk: "None after narrow application and compilation."

- failure_id: FAIL-20260904-076
  fingerprint: "failure-record-readiness-partial-line-recurrence"
  related_failures: ["LIB-FAIL-1810", "LIB-FAIL-1779"]
  recurrence_count: 1
  status: RECURRENCE_PROVEN
  severity: LOW
  classification: DOCUMENTATION_MAINTENANCE
  violated_invariant: "apply_patch replaces complete observed physical lines."
  root_cause: "The readiness hunk contained only the beginning of the current paragraph."
  canonical_correction: "Use the complete observed readiness line in the corrected transaction."
  regression_test: "Failure IDs and readiness counts are present exactly once."
  clean_rebuild_evidence: "The invalid patch was rejected atomically."
  residual_risk: "None after exact-line update."

- failure_id: FAIL-20260904-077
  fingerprint: "domainbind-protected-conditional-missing-close"
  related_failures: ["LIB-FAIL-1811"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: IMPLEMENTATION_DEFECT
  violated_invariant: "Every source edit must pass formatting/parser gate before behavioral claims."
  root_cause: "The outer protected branch was not closed after the inner authorization error branch."
  canonical_correction: "Add the missing brace and rerun formatting and tests."
  regression_test: "gofmt no-op plus domainbind focal and full Go suites."
  clean_rebuild_evidence: "gofmt stopped at the exact line before any test promotion."
  residual_risk: "None after parser and behavioral gates pass."

- failure_id: FAIL-20260904-078
  fingerprint: "quote-http-test-symbol-name-abbreviated"
  related_failures: ["LIB-FAIL-1812"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TEST_MAINTENANCE
  violated_invariant: "Test patches use exact observed symbols."
  root_cause: "The target test name was shortened in the patch context."
  canonical_correction: "Anchor the new test at the exact existing function name."
  regression_test: "The new HTTP quote tests compile and pass."
  clean_rebuild_evidence: "Rejected atomically before source modification."
  residual_risk: "None after exact patch and test execution."

- failure_id: FAIL-20260904-079
  fingerprint: "quote-concurrency-fixture-downstream-count-stale"
  related_failures: ["LIB-FAIL-1813"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TEST_MAINTENANCE
  violated_invariant: "Every durable fixture addition updates all downstream cardinality assertions."
  root_cause: "The final customer journey assertion still expected the pre-concurrency quote count."
  canonical_correction: "Update only the aggregate count while retaining the accepted-quote invariant."
  regression_test: "Fresh database integration passes and still proves one concurrent resource/event."
  clean_rebuild_evidence: "Failure output enumerated exactly the four expected quotes."
  residual_risk: "None after clean rerun."

- failure_id: FAIL-20260904-080
  fingerprint: "pack-metadata-patch-contained-spurious-line"
  related_failures: ["LIB-FAIL-1814"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: DOCUMENTATION_MAINTENANCE
  violated_invariant: "Promotion patches contain only intended, observed replacements."
  root_cause: "A transient malformed line remained at the end of the patch."
  canonical_correction: "Remove it and apply exact metadata hunks per pack."
  regression_test: "Version, claim, compatibility and verification date occur exactly once."
  clean_rebuild_evidence: "The malformed transaction was rejected atomically."
  residual_risk: "None after exact patch verification."

- failure_id: FAIL-20260904-081
  fingerprint: "version-reference-search-regex-quoting-invalid"
  related_failures: ["LIB-FAIL-1815"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  violated_invariant: "Exact version discovery uses literal queries when regex adds no value."
  root_cause: "Nested quoting truncated an alternation expression."
  canonical_correction: "Run independent fixed-string searches over Markdown."
  regression_test: "All three literal searches execute without parser errors."
  clean_rebuild_evidence: "The invalid command returned only a regex error."
  residual_risk: "None after literal discovery."

- failure_id: FAIL-20260904-082
  fingerprint: "materializer-discovery-literal-glob-recurrence"
  related_failures: ["LIB-FAIL-1816", "LIB-FAIL-1807"]
  recurrence_count: 4
  status: RECURRENCE_PROVEN
  severity: LOW
  classification: TOOLING
  violated_invariant: "Windows discovery never passes wildcard paths directly to rg."
  root_cause: "A redundant follow-up query used a literal glob after the inventory had already found the script."
  canonical_correction: "Use the exact inventory result for all subsequent reads."
  regression_test: "The exact materializer path opens and executes."
  clean_rebuild_evidence: "No artifact decision depended on the rejected query."
  residual_risk: "Avoid combined discovery/read commands with path globs."

- failure_id: FAIL-20260904-083
  fingerprint: "isolated-roundtrip-go-fmt-requires-module"
  related_failures: ["LIB-FAIL-1817"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: RELEASE_REPRODUCIBILITY
  violated_invariant: "Format idempotence on isolated files uses gofmt, not module package resolution."
  root_cause: "Absolute directories were passed to go fmt outside a go.mod tree."
  canonical_correction: "Run gofmt.exe on every materialized .go file and compare SHA-256 before/after."
  regression_test: "All materialized Go files remain byte-identical after gofmt."
  clean_rebuild_evidence: "All three round-trip hash comparisons passed before the format substep failed."
  residual_risk: "None after the corrected no-op check."

- failure_id: FAIL-20260904-084
  fingerprint: "v180-evidence-search-literal-glob-recurrence"
  related_failures: ["LIB-FAIL-1818", "LIB-FAIL-1816"]
  recurrence_count: 5
  status: RECURRENCE_PROVEN
  severity: LOW
  classification: TOOLING
  violated_invariant: "File selection on Windows starts from rg --files or exact paths."
  root_cause: "A wildcard path was appended to an otherwise valid plan query."
  canonical_correction: "Use exact evidence paths returned by inventory."
  regression_test: "Subsequent reads contain no invalid wildcard path."
  clean_rebuild_evidence: "The valid plan output is separable from the rejected evidence branch."
  residual_risk: "Stop composing wildcard paths in rg arguments."

- failure_id: FAIL-20260904-085
  fingerprint: "channel-dispatcher-drops-ingress-command-identity"
  related_failures: ["LIB-FAIL-1819", "LIB-FAIL-1743"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: CRITICAL
  classification: INTEGRATION_INTEGRITY
  violated_invariant: "Channel-to-domain flow preserves a stable provider message identity and occurrence time."
  root_cause: "Responder accepted only tenant and text, discarding the envelope needed for replay-safe commands."
  canonical_correction: "Carry the full validated message through dispatcher and require ingress-owned message ID/time."
  regression_test: "Missing ingress identity fails; responder observes exact envelope; reply preserves channel/tenant/contact/thread."
  clean_rebuild_evidence: "V232: full envelope regression, 47-package suite, 4/4 materialization and gofmt idempotence PASS."
  residual_risk: "Runtime persistence and provider adapters still require their own evidence."

- failure_id: FAIL-20260904-086
  fingerprint: "channels-message-owner-filename-assumed"
  related_failures: ["LIB-FAIL-1820"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: IMPLEMENTATION_MAINTENANCE
  violated_invariant: "Physical owners are enumerated before patching concatenated output."
  root_cause: "The package content was read with a wildcard and the first filename was inferred incorrectly."
  canonical_correction: "Patch the enumerated owner channel.go."
  regression_test: "Package formatting and tests pass after the exact edit."
  clean_rebuild_evidence: "The invalid patch was rejected atomically."
  residual_risk: "None after exact owner selection."

- failure_id: FAIL-20260904-087
  fingerprint: "domain-gateway-static-lead-organization-across-conversations"
  related_failures: ["LIB-FAIL-1821", "LIB-FAIL-1743"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: CRITICAL
  classification: TENANT_DATA_ISOLATION
  violated_invariant: "Every inbound contact resolves its own authorized organization and lead before domain effects."
  root_cause: "Gateway configuration stores one global lead/organization and command methods reuse it."
  canonical_correction: "Introduce validated dynamic scope passed by an authorized contact resolver and prove two contacts cannot cross."
  regression_test: "Requests for two scopes carry their own path/body; mismatched/empty scope fails before HTTP."
  clean_rebuild_evidence: "V233: two-contact HTTP regression, full 47-package suite, 3/3 materialization and gofmt PASS."
  residual_risk: "The connected runtime must inject an authorized contact resolver and call only *For methods."

- failure_id: FAIL-20260904-088
  fingerprint: "domainbind-03-metadata-partial-paragraph"
  related_failures: ["LIB-FAIL-1822", "LIB-FAIL-1810"]
  recurrence_count: 2
  status: REGRESSION_PROVEN
  severity: LOW
  classification: DOCUMENTATION_MAINTENANCE
  violated_invariant: "Markdown paragraph edits use complete observed physical lines."
  root_cause: "Only an interior sentence was supplied as deletion context."
  canonical_correction: "Replace the full observed paragraph."
  regression_test: "Pack version/claim/contract describe dynamic scope exactly once."
  clean_rebuild_evidence: "The invalid transaction was rejected atomically."
  residual_risk: "None after exact-line update."

- failure_id: FAIL-20260904-089
  fingerprint: "channel-reply-has-no-stable-delivery-key"
  related_failures: ["LIB-FAIL-1823", "LIB-FAIL-1743"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: CRITICAL
  classification: DELIVERY_INTEGRITY
  violated_invariant: "Retrying a completed turn must expose one stable outbound delivery identity."
  root_cause: "The dispatcher reconstructed only address/scope/text for the reply."
  canonical_correction: "Require and derive a deterministic DeliveryKey from the immutable inbound identity."
  regression_test: "Repeated dispatch of one inbound produces the same key; different inbound IDs produce different keys."
  clean_rebuild_evidence: "V234: deterministic replay/distinct-ID regression, 47-package suite and 4/4 materialization/gofmt PASS."
  residual_risk: "A provider without idempotent send requires a durable delivery receipt/reconciliation lane."

- failure_id: FAIL-20260904-090
  fingerprint: "delivery-key-test-hardcodes-first-message-id"
  related_failures: ["LIB-FAIL-1824"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TEST_MAINTENANCE
  violated_invariant: "Expanded fixtures update callback assertions for every admitted case."
  root_cause: "The responder asserted the first ID while the test intentionally dispatched a second."
  canonical_correction: "Assert non-empty preserved identity in callback and exact key behavior after dispatch."
  regression_test: "Same inbound key is stable; distinct inbound key differs."
  clean_rebuild_evidence: "Failure output proved the new ID reached the responder unchanged."
  residual_risk: "None after corrected test."

- failure_id: FAIL-20260904-091
  fingerprint: "safety-allowlist-bypasses-injection-and-output-leak-detection"
  related_failures: ["LIB-FAIL-1825", "LIB-FAIL-1743"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: CRITICAL
  classification: SECURITY
  violated_invariant: "A scoped PII exception never disables prompt-injection or secret-leak controls."
  root_cause: "Allowlisted returned before independent blocked-phrase checks."
  canonical_correction: "Run blocked-signal detection first and apply allowlist only to PII/refusal semantics."
  regression_test: "Allowlisted injection and allowlisted system-prompt output both remain blocked."
  clean_rebuild_evidence: "V235: bypass regressions, full 47-package suite, 10/10 materialization and gofmt PASS."
  residual_risk: "Heuristics remain conditioned on a production classifier and PII recognizer."

- failure_id: FAIL-20260904-092
  fingerprint: "aifoundation-historical-eval-test-gofmt-drift"
  related_failures: ["LIB-FAIL-1826", "LIB-FAIL-1803"]
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: RELEASE_REPRODUCIBILITY
  violated_invariant: "Every Go pack is gofmt-stable across all materialized files."
  root_cause: "The historical foundation gate did not check whole-package formatting idempotence."
  canonical_correction: "Re-embed eval_test.go from gofmt output and verify all ten files after reconstruction."
  regression_test: "Fresh 10/10 materialization remains byte-identical after gofmt."
  clean_rebuild_evidence: "Full tests passed after the formatter named the drifted file."
  residual_risk: "Keep whole-pack gofmt checks for every touched historical pack."

- failure_id: FAIL-20260904-093
  fingerprint: "conversation-runtime-go-toolchain-absent-from-shell-path"
  related_failures: ["LIB-FAIL-1827"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: BUILD_ENVIRONMENT
  violated_invariant: "Build evidence must come from the exact resolved toolchain, never from an assumed shell PATH."
  root_cause: "The first runtime gate invoked go/gofmt by global name although this host keeps the fixed toolchain under Temp."
  canonical_correction: "Resolve and invoke the existing go.exe and sibling gofmt.exe by absolute path."
  regression_test: "Run package and full-suite gates using the explicit toolchain and record its version."
  clean_rebuild_evidence: "The failed invocation changed no source and supplies no compilation evidence."
  residual_risk: "A new host must materialize or resolve the pinned toolchain during preflight."

- failure_id: FAIL-20260904-094
  fingerprint: "conversation-runtime-postgres-migration-directory-assumed"
  related_failures: ["LIB-FAIL-1828"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: DISCOVERY
  violated_invariant: "Filesystem owners are enumerated before use."
  root_cause: "A generic migrations/ path was assumed instead of db/migrations/."
  canonical_correction: "Use rg --files and the observed db/migrations owner."
  regression_test: "All 46 observed migrations applied in sorted order."
  clean_rebuild_evidence: "Fresh database migration gate passed."
  residual_risk: "None for this tree."

- failure_id: FAIL-20260904-095
  fingerprint: "conversation-runtime-package-filenames-assumed"
  related_failures: ["LIB-FAIL-1829"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: DISCOVERY
  violated_invariant: "Interfaces are taken from exact owners."
  root_cause: "client.go and budget.go were guessed filenames."
  canonical_correction: "Enumerate packages and read responses_tools.go and ledger.go."
  regression_test: "Runtime compiles against exact interfaces."
  clean_rebuild_evidence: "Full 48-package suite passed."
  residual_risk: "None."

- failure_id: FAIL-20260904-096
  fingerprint: "conversation-runtime-race-gate-cgo-unavailable"
  related_failures: ["LIB-FAIL-1830"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: BUILD_ENVIRONMENT
  violated_invariant: "Unavailable gates are not reported as executed."
  root_cause: "The pinned Windows toolchain lacks CGO for -race."
  canonical_correction: "Retain race as a capable-host condition and run all available concurrency/database gates."
  regression_test: "Ordinary unit, PostgreSQL concurrency, full suite, vet and build pass."
  clean_rebuild_evidence: "V236 records no race PASS claim."
  residual_risk: "Race detector must run in project CI on a CGO-capable host."

- failure_id: FAIL-20260904-097
  fingerprint: "portable-postgres-dropdb-binary-absent"
  related_failures: ["LIB-FAIL-1831"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: BUILD_ENVIRONMENT
  violated_invariant: "Portable tool contents are inspected before invocation."
  root_cause: "dropdb.exe was assumed present."
  canonical_correction: "Use exact psql.exe with DROP DATABASE WITH (FORCE)."
  regression_test: "Fresh database was created and migrated 46/46."
  clean_rebuild_evidence: "PostgreSQL 18.6 gate passed."
  residual_risk: "None."

- failure_id: FAIL-20260904-098
  fingerprint: "domain-scope-patch-context-inexact"
  related_failures: ["LIB-FAIL-1832"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: EDITING
  violated_invariant: "Patches use observed exact context."
  root_cause: "The test invocation was remembered incorrectly."
  canonical_correction: "Re-read and patch the literal noLead.Validate block."
  regression_test: "Domain package and full suite pass."
  clean_rebuild_evidence: "Initial patch was atomic and changed nothing."
  residual_risk: "None."

- failure_id: FAIL-20260904-099
  fingerprint: "app-test-delete-add-same-patch-rejected"
  related_failures: ["LIB-FAIL-1833"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: EDITING
  violated_invariant: "Use supported atomic patch operations."
  root_cause: "Delete and add targeted the same path in one patch."
  canonical_correction: "Separate the operations."
  regression_test: "Final app_test.go compiles and passes E2E."
  clean_rebuild_evidence: "Rejected operation changed nothing."
  residual_risk: "None."

- failure_id: FAIL-20260904-100
  fingerprint: "app-e2e-quote-route-name-obsolete"
  related_failures: ["LIB-FAIL-1834"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: CONTRACT_TEST
  violated_invariant: "Test expectations mirror canonical gateway routes."
  root_cause: "The test used quotations instead of quotes."
  canonical_correction: "Expect /v1/franchise/quotes."
  regression_test: "PostgreSQL app replay E2E passes."
  clean_rebuild_evidence: "Observed request proved bearer and key were already correct."
  residual_risk: "Backend route must be revalidated in each composed release."

- failure_id: FAIL-20260904-101
  fingerprint: "pack-generator-textdecoder-unavailable"
  related_failures: ["LIB-FAIL-1835"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  violated_invariant: "Generator uses only runtime-supported primitives."
  root_cause: "TextDecoder was assumed present in the JS isolate."
  canonical_correction: "Return normalized source strings through PowerShell JSON."
  regression_test: "New pack materializes 7/7 with exact hashes."
  clean_rebuild_evidence: "Failed attempt created no file."
  residual_risk: "None."

- failure_id: FAIL-20260904-102
  fingerprint: "multi-pack-roundtrip-stale-lastexitcode"
  related_failures: ["LIB-FAIL-1836"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  violated_invariant: "A gate reads status owned by the invoked operation."
  root_cause: "LASTEXITCODE was inherited across a PowerShell script call."
  canonical_correction: "Use terminating script exceptions and fresh destinations."
  regression_test: "All four packs materialize and compare 17/17."
  clean_rebuild_evidence: "The first pack was valid despite the false throw."
  residual_risk: "Avoid LASTEXITCODE after pure PowerShell scripts."

- failure_id: FAIL-20260904-103
  fingerprint: "franchise-plan-unsupported-composition-schema-version"
  related_failures: ["LIB-FAIL-1837"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: COMPOSITION
  violated_invariant: "Plans use only the schema version accepted by the verifier."
  root_cause: "Editorial revision was encoded as compositionVersion 1.1."
  canonical_correction: "Keep schema 1.0 and version packs/content independently."
  regression_test: "Root verifier accepts plan and reaches only the expected open-failure gate."
  clean_rebuild_evidence: "The value was corrected before release verification."
  residual_risk: "Schema evolution requires verifier support first."

- failure_id: FAIL-20260904-104
  fingerprint: "v236-pack-required-sections-omitted"
  related_failures: ["LIB-FAIL-1838"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: PACK_CONTRACT
  violated_invariant: "Every pack contains the ten exact non-empty normative sections."
  root_cause: "Regenerated packs stopped at a custom verification heading."
  canonical_correction: "Restore sections 6 through 10 in exact order."
  regression_test: "VERIFY_LIBRARY pack contract scan."
  clean_rebuild_evidence: "Root PASS accepted all four V236 packs."
  residual_risk: "None."

- failure_id: FAIL-20260904-105
  fingerprint: "meta-import-pack-dependency-section-missing"
  related_failures: ["LIB-FAIL-1839"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: PACK_CONTRACT
  violated_invariant: "Dependency bill is explicit."
  root_cause: "V227 omitted required H2 7."
  canonical_correction: "Record Go/PostgreSQL/pgx/Meta artifact and renumber sections."
  regression_test: "Root pack contract scan."
  clean_rebuild_evidence: "VERIFY_LIBRARY_PASS."
  residual_risk: "Meta live access remains a declared condition."

- failure_id: FAIL-20260904-106
  fingerprint: "responses-pack-dependency-section-missing"
  related_failures: ["LIB-FAIL-1840"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: PACK_CONTRACT
  violated_invariant: "Dependency bill is explicit."
  root_cause: "V228 omitted required H2 7."
  canonical_correction: "Record exact Go and remote API dependency and renumber."
  regression_test: "Root pack contract scan."
  clean_rebuild_evidence: "VERIFY_LIBRARY_PASS and executable Audit PASS."
  residual_risk: "Provider terms/model/evals remain target conditions."

- failure_id: FAIL-20260904-107
  fingerprint: "meta-reconcile-reconstruction-heading-shifted"
  related_failures: ["LIB-FAIL-1841"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: PACK_CONTRACT
  violated_invariant: "Reconstruction evidence is exact H2 10."
  root_cause: "Production blockers was introduced as an extra numbered H2."
  canonical_correction: "Make blockers an H3 and evidence H2 10."
  regression_test: "Root structural parser."
  clean_rebuild_evidence: "VERIFY_LIBRARY_PASS."
  residual_risk: "None structural."

- failure_id: FAIL-20260904-108
  fingerprint: "v236-third-party-provenance-count-drift"
  related_failures: ["LIB-FAIL-1842"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: LICENSE_PROVENANCE
  violated_invariant: "Published provenance totals equal parsed materialization blocks."
  root_cause: "V236 blocks were added without synchronizing notices."
  canonical_correction: "Set exact 1074 authored, 124 adapted, 105 verbatim, total 1303."
  regression_test: "Root derived-count assertion."
  clean_rebuild_evidence: "VERIFY_LIBRARY_PASS packs=148 files=1303."
  residual_risk: "Keep derived ledger synchronized on every pack change."

- failure_id: FAIL-20260904-109
  fingerprint: "historical-composition-pack-version-drift"
  related_failures: ["LIB-FAIL-1843"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: COMPOSITION
  violated_invariant: "Every plan pins the exact current pack metadata version."
  root_cause: "Backend and serverless consumers were not updated with promotions."
  canonical_correction: "Synchronize plan entries mechanically from pack metadata."
  regression_test: "All 44 profiles pass root plan validation."
  clean_rebuild_evidence: "VERIFY_LIBRARY_PASS; franchise complete materializes 612 files."
  residual_risk: "Any future promotion must update every consuming plan."

- failure_id: FAIL-20260904-110
  fingerprint: "broad-stale-string-scan-output-truncated"
  related_failures: ["LIB-FAIL-1844"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: EVIDENCE_COLLECTION
  violated_invariant: "An audit conclusion is based only on complete, inspectable command output."
  root_cause: "A repository-wide multi-pattern search mixed current owners with historical snapshots and exceeded the output budget."
  canonical_correction: "Discard the truncated output; enumerate candidate paths with rg -l and inspect each current owner with bounded queries."
  regression_test: "Bounded path-first stale-reference audit plus VERIFY_LIBRARY_PASS."
  clean_rebuild_evidence: "No file was changed from the truncated result; bounded audit is executed before closure."
  residual_risk: "Historical records intentionally retain old counts and must not be rewritten as current state."

- failure_id: FAIL-20260904-111
  fingerprint: "crm-migration-filename-assumed"
  related_failures: ["LIB-FAIL-1845"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: DISCOVERY
  violated_invariant: "File owners are observed before they are referenced."
  root_cause: "The CRM migration filename was inferred from its domain instead of enumerated from the composed tree."
  canonical_correction: "Discard that read, enumerate db/migrations, and open the exact CRM owner."
  regression_test: "Path-first CRM/contact identity discovery."
  clean_rebuild_evidence: "The failed read made no change and no schema decision depends on it."
  residual_risk: "None after exact owner discovery."

- failure_id: FAIL-20260904-112
  fingerprint: "contact-identity-e2e-residual-fixture"
  related_failures: ["LIB-FAIL-1846"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: TEST_ISOLATION
  violated_invariant: "Database integration tests begin from a dependency-complete clean fixture."
  root_cause: "A transient earlier test created 0047 binding rows; the restored historical cleanup did not delete those new dependencies."
  canonical_correction: "Delete 0047 decisions/bindings before lead, organization and tenant, then rerun all affected suites."
  regression_test: "Clean rerun of app, contactidentity and platform/postgres with PostgreSQL 18.6."
  clean_rebuild_evidence: "The failed run stopped on fixture creation before exercising the journey."
  residual_risk: "Fixed IDs remain safe only with dependency-complete cleanup."

- failure_id: FAIL-20260904-113
  fingerprint: "immutable-contact-decision-cleanup-rejected"
  related_failures: ["LIB-FAIL-1847"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: TEST_ISOLATION
  violated_invariant: "Tests must not weaken or fight append-only production evidence."
  root_cause: "The cleanup strategy issued DELETE against an immutable decision ledger and reused fixed tenant IDs."
  canonical_correction: "Use unique tenant UUIDs on each run; retain evidence inside the disposable test database and reset that database only through controlled migration lifecycle."
  regression_test: "Run the PostgreSQL integration suites twice consecutively without duplicate or immutable-ledger errors."
  clean_rebuild_evidence: "The trigger rejected deletion as designed; no production invariant was bypassed."
  residual_risk: "Long-lived shared test databases require periodic whole-database recreation, not row deletion."

- failure_id: FAIL-20260904-114
  fingerprint: "unique-tenant-id-static-tenant-code"
  related_failures: ["LIB-FAIL-1848"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: TEST_ISOLATION
  violated_invariant: "Every unique key of a retained test fixture is unique across reruns."
  root_cause: "The UUID became per-run while tenant_code remained constant."
  canonical_correction: "Derive tenant_code from the generated UUID as well."
  regression_test: "Two consecutive app/contact/PostgreSQL suite runs."
  clean_rebuild_evidence: "No journey executed before the uniqueness violation."
  residual_risk: "None for the observed unique keys."

- failure_id: FAIL-20260904-115
  fingerprint: "public-map-heading-assumed"
  related_failures: ["LIB-FAIL-1849"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: DOCUMENTATION_EDIT
  violated_invariant: "Multi-owner patches use observed anchors for every target."
  root_cause: "The public architecture map heading was guessed."
  canonical_correction: "Apply exact known rows separately, inspect public-map headings, then patch its observed owner."
  regression_test: "Root structural verification after bounded map updates."
  clean_rebuild_evidence: "apply_patch rejected the transaction atomically."
  residual_risk: "None after exact-anchor update."

- failure_id: FAIL-20260904-116
  fingerprint: "v237-third-party-provenance-count-drift"
  related_failures: ["LIB-FAIL-1850"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: LICENSE_PROVENANCE
  violated_invariant: "Published provenance totals equal parsed materialization blocks."
  root_cause: "Eight V237 AUTHORED blocks were added before synchronizing the derived notices marker."
  canonical_correction: "Publish exactly the gate-computed totals: AUTHORED=1082, ADAPTED=124, VERBATIM=105, TOTAL=1311."
  regression_test: "VERIFY_LIBRARY_PASS."
  clean_rebuild_evidence: "The root verifier failed closed at the provenance assertion."
  residual_risk: "Every future pack change must synchronize the derived ledger."

- failure_id: FAIL-20260904-117
  fingerprint: "go-list-outside-materialized-module"
  related_failures: ["LIB-FAIL-1851"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: EVIDENCE_COLLECTION
  violated_invariant: "Toolchain inventory runs from the module root it claims to measure."
  root_cause: "go list was invoked from the Markdown library root, which has no go.mod."
  canonical_correction: "Repeat only go list from the V237 franchise composition root."
  regression_test: "Non-zero exact package inventory plus full suite already executed from that root."
  clean_rebuild_evidence: "The invalid zero count is discarded; file hashes from exact paths remain valid."
  residual_risk: "None after module-root execution."

- failure_id: FAIL-20260904-118
  fingerprint: "v237-evidence-file-increments-markdown-count"
  related_failures: ["LIB-FAIL-1852"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: EVIDENCE_COLLECTION
  violated_invariant: "Published inventory is captured after all evidence files exist."
  root_cause: "The V237 evidence recorded the pre-evidence Markdown cardinality."
  canonical_correction: "Use final root output: markdown_files=656."
  regression_test: "Final VERIFY_LIBRARY_PASS after evidence synchronization."
  clean_rebuild_evidence: "Pack, materialized-file and profile counts remained unchanged."
  residual_risk: "Any new Markdown intentionally increments this non-normative display count."

- failure_id: FAIL-20260904-119
  fingerprint: "v238-unbounded-output-capture-truncation"
  related_failures: ["LIB-FAIL-1853"]
  recurrence_count: 6
  status: REGRESSION_PROVEN
  severity: LOW
  classification: EVIDENCE_COLLECTION
  violated_invariant: "Only complete bounded output may support a verification conclusion."
  root_cause: "The first Audit read and a later broad Markdown search requested more accumulated output than the consumer context could retain."
  canonical_correction: "Discard truncated output, poll the existing process with a small output budget, and use focal file-scoped searches."
  regression_test: "Bounded poll returned exit_code=0 and the complete final VERIFY_EXECUTABLE_LIBRARY_PASS line."
  clean_rebuild_evidence: "Audit closed with packs=150 and upstream_sources=121; the broad search was not used as an exhaustive inventory."
  residual_risk: "Future high-cardinality searches must remain file-scoped or paginated."

- failure_id: FAIL-20260904-120
  fingerprint: "v239-rg-windows-wildcard-and-assumed-template"
  related_failures: ["LIB-FAIL-1854"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  violated_invariant: "Enumerate paths before reading them."
  root_cause: "A wildcard path and an unobserved template filename were passed directly to rg."
  canonical_correction: "Filter rg --files output, then read exact observed paths."
  regression_test: "The exact source-profile template was located and read."
  clean_rebuild_evidence: "No failed output governed V239."
  residual_risk: "None when file inventory precedes reads."

- failure_id: FAIL-20260904-121
  fingerprint: "v239-assumed-meta-import-package-path"
  related_failures: ["LIB-FAIL-1855"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: DISCOVERY
  violated_invariant: "Package paths are observed, never inferred from pack names."
  root_cause: "The Meta importer was guessed as internal/metaleadimport instead of internal/leadstream."
  canonical_correction: "Enumerate internal files and open meta_import.go/test at their real path."
  regression_test: "Both real files were read before V239 design."
  clean_rebuild_evidence: "V239 is compatible with the observed importer API."
  residual_risk: "None."

- failure_id: FAIL-20260904-122
  fingerprint: "v239-tool-discovery-shell-glob"
  related_failures: ["LIB-FAIL-1856"]
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  violated_invariant: "Windows path discovery does not pass shell globs as literal paths."
  root_cause: "The tool search repeated a non-portable *.ps1 path and assumed a tools directory."
  canonical_correction: "Use rg --files and filter filenames."
  regression_test: "update_pack_from_tree.ps1 was found, read and used successfully."
  clean_rebuild_evidence: "Nine blocks round-tripped exactly."
  residual_risk: "None."

- failure_id: FAIL-20260904-123
  fingerprint: "v239-go-fmt-cross-directory-named-files"
  related_failures: ["LIB-FAIL-1857"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  violated_invariant: "Formatting invocation matches Go tool semantics."
  root_cause: "go fmt received named files from multiple directories."
  canonical_correction: "Use the exact isolated gofmt executable over the allowlist."
  regression_test: "gofmt completed and all focal/full Go gates passed."
  clean_rebuild_evidence: "51 packages, vet and build PASS."
  residual_risk: "None."

- failure_id: FAIL-20260904-124
  fingerprint: "v239-meta-batch-fixture-nesting"
  related_failures: ["LIB-FAIL-1858"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: TEST_FIXTURE
  violated_invariant: "A batch fixture preserves the official entry/changes hierarchy."
  root_cause: "A suffix replacement nested the second entry inside changes."
  canonical_correction: "Use an explicit two-entry official-shape fixture."
  regression_test: "Decoder returns two signals; replay/divergence tests pass."
  clean_rebuild_evidence: "Focal suite passed twice."
  residual_risk: "Live batching remains a project probe."

- failure_id: FAIL-20260904-125
  fingerprint: "v239-postgres-stale-port-assumption"
  related_failures: ["LIB-FAIL-1859"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TEST_ENVIRONMENT
  violated_invariant: "Runtime endpoints are observed before use."
  root_cause: "Port 55432 was reused although the active cluster listened on 55445."
  canonical_correction: "Read postmaster.pid/process state and probe the observed port."
  regression_test: "pg_isready 55445 accepted connections."
  clean_rebuild_evidence: "0049 and a clean 49-migration database passed."
  residual_risk: "Ports remain host-local evidence."

- failure_id: FAIL-20260904-126
  fingerprint: "v239-pgctl-start-active-data-directory"
  related_failures: ["LIB-FAIL-1860"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: TEST_ENVIRONMENT
  violated_invariant: "Inspect process ownership before starting a database cluster."
  root_cause: "A start was attempted after checking the wrong port, while PID 99204 already owned the data directory."
  canonical_correction: "Do not mutate the process; connect to its observed port."
  regression_test: "Existing process served all SQL/Go integration gates."
  clean_rebuild_evidence: "No competing postmaster was created."
  residual_risk: "None after endpoint discovery."

- failure_id: FAIL-20260904-127
  fingerprint: "v239-powershell-foreach-empty-pipe"
  related_failures: ["LIB-FAIL-1861"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: EVIDENCE_COLLECTION
  violated_invariant: "Hash inventory commands must parse before execution."
  root_cause: "A foreach statement was piped directly to ConvertTo-Json."
  canonical_correction: "Accumulate objects in an array and pipe the completed array."
  regression_test: "Nine complete hashes and sizes were emitted."
  clean_rebuild_evidence: "Pack round-trip validated the hashes."
  residual_risk: "None."

- failure_id: FAIL-20260904-128
  fingerprint: "v239-pack-required-sections-omitted"
  related_failures: ["LIB-FAIL-1862"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: PACK_CONTRACT
  violated_invariant: "Every implementation pack exposes sections 1-10 exactly once."
  root_cause: "The initial skeleton compressed three required sections."
  canonical_correction: "Add Configuration surface, Dependency bill and Apply order with exact headings."
  regression_test: "Root gate advanced past pack structure."
  clean_rebuild_evidence: "Final root verification passed."
  residual_risk: "None."

- failure_id: FAIL-20260904-129
  fingerprint: "v239-provenance-count-drift"
  related_failures: ["LIB-FAIL-1863"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: PROVENANCE
  violated_invariant: "Notices equal gate-computed materialization provenance."
  root_cause: "The new pack preceded synchronization of the derived V238 marker."
  canonical_correction: "Publish only AUTHORED=1100, ADAPTED=126, VERBATIM=105, TOTAL=1331 from the gate."
  regression_test: "VERIFY_LIBRARY_PASS 151/1331/660/44 and profile 59/640."
  clean_rebuild_evidence: "No source classification was changed to force PASS."
  residual_risk: "Future packs intentionally reopen this derived gate."

- failure_id: FAIL-20260904-130
  fingerprint: "v239-nonfocal-recurrence-counter-patch"
  related_failures: ["LIB-FAIL-1864"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: DOCUMENTATION_MAINTENANCE
  violated_invariant: "A counter edit targets the exact failure fingerprint."
  root_cause: "A context-poor patch changed the first recurrence_count=1 instead of FAIL-20260904-119."
  canonical_correction: "Restore FAIL-20260902-003 to one and update FAIL-20260904-119 to two using unique failure IDs."
  regression_test: "Targeted rg shows the intended counters at both IDs."
  clean_rebuild_evidence: "Historical failure identity and chronology are preserved."
  residual_risk: "None when patches anchor on failure_id."

- failure_id: FAIL-20260904-131
  fingerprint: "v239-postgres-batch-commandtag-syntax"
  related_failures: ["LIB-FAIL-1865"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: IMPLEMENTATION
  violated_invariant: "The hardened organization check must compile before promotion."
  root_cause: "The patch retained an if-style semicolon after converting Exec to a short assignment."
  canonical_correction: "Separate Exec assignment from the error branch and run gofmt."
  regression_test: "Focal leadstream/PostgreSQL suite passes and cross-organization replay is rejected."
  clean_rebuild_evidence: "No failing bytes were promoted to the pack; four blocks were refreshed from corrected sources."
  residual_risk: "Full/Audit gates are reopened after the byte change."

- failure_id: FAIL-20260904-132
  fingerprint: "v239-meta-batch-organization-replay-scope"
  related_failures: ["LIB-FAIL-1866"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: MULTITENANCY
  violated_invariant: "A provider batch cannot cross organization ownership inside a tenant."
  root_cause: "The first conflict path keyed batch by tenant/hash but did not compare its stored organization before processing signals."
  canonical_correction: "Compare organization on batch replay, fail divergent ownership, and validate provider identifier lengths before persistence."
  regression_test: "Same tenant/payload under store-2 is rejected; tenant B remains isolated."
  clean_rebuild_evidence: "Corrected pack 27e60c84..., clean 49 migrations, full 51 packages, root and Audit pass."
  residual_risk: "Real Page-to-organization mapping remains a live project input and probe."

- failure_id: FAIL-20260904-133
  fingerprint: "aidlc-audit-unnecessary-temp-recursive-cleanup"
  related_failures: ["LIB-FAIL-1867"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  violated_invariant: "Read-only source acquisition must not request an unnecessary destructive cleanup operation."
  root_cause: "The first command reused a fixed temp name and combined acquisition with recursive deletion, so host policy rejected it before execution."
  canonical_correction: "Use a fresh UUID temporary directory and perform no deletion before download or expansion."
  regression_test: "The rejected attempt created no files; the next acquisition uses only New-Item, download, hash and expansion inside the fresh path."
  clean_rebuild_evidence: "No library or project content changed during the rejected command."
  residual_risk: "Temporary audit bytes remain disposable and outside the portable archive."

- failure_id: FAIL-20260904-134
  fingerprint: "execution-kit-tests-without-authority-root"
  related_failures: ["LIB-FAIL-1868"]
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TEST_ENVIRONMENT
  violated_invariant: "Historical execution-validator tests must run with their declared authority documents available at the workspace root."
  root_cause: "The first combined test command ran against only the isolated materialized kit, so five historical fixtures correctly rejected missing authority paths."
  canonical_correction: "Retain the 12/12 focal execution-state result; inspect find_workspace; copy its SYSTEMS_ENGINEERING_MASTER_MAP.md root sentinel plus the exact authority files before rerunning."
  regression_test: "The integral suite must pass from the reconstructed project-like root recognized by find_workspace before the pack is promoted."
  clean_rebuild_evidence: "The failing output named only absent authority documents; no source bytes were changed to bypass the checks."
  residual_risk: "Other projects must likewise retain or remap their authority references rather than disabling path validation."

- failure_id: FAIL-20260904-135
  fingerprint: "nested-pwsh-array-parameter-flattening"
  related_failures: ["LIB-FAIL-1869"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  violated_invariant: "The pack updater must receive the exact seven-file allowlist as one PowerShell array."
  root_cause: "A nested pwsh -File boundary flattened the local array into positional arguments."
  canonical_correction: "Invoke update_pack_from_tree.ps1 directly with the call operator in the current PowerShell session."
  regression_test: "Updater reports seven blocks updated and clean materialization matches every source hash."
  clean_rebuild_evidence: "The failed updater invocation terminated before writing any block."
  residual_risk: "Future array-valued script parameters should avoid unnecessary process boundaries."

- failure_id: FAIL-20260904-136
  fingerprint: "multi-file-doc-patch-stale-bootstrap-command-context"
  related_failures: ["LIB-FAIL-1870"]
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: LOW
  classification: DOCUMENTATION_TOOLING
  violated_invariant: "A documentation integration patch must match the current canonical command text before changing multiple entry points."
  root_cause: "The combined patch expected a newer long-form validate_project invocation, while the bootstrap still used its historical positional syntax."
  canonical_correction: "Inspect the exact local excerpts and apply small file-scoped patches against observed context; then run the root gates."
  regression_test: "All intended entry points contain the execution-state startup contract and VERIFY_LIBRARY passes without broken references."
  clean_rebuild_evidence: "Corrected file-scoped patches applied; all six startup surfaces reference the 1.2.0 state contract; a CLI mismatch caught by rg was corrected to the actual positional/--level interface; clean 15-file round-trip and 20/20 tests pass."
  residual_risk: "No open local regression; project-specific live and production gates remain separate."

- failure_id: FAIL-20260904-137
  fingerprint: "tiktok-gap-search-included-unobserved-paths"
  related_failures: ["LIB-FAIL-1871"]
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: LOW
  classification: DISCOVERY_TOOLING
  violated_invariant: "Repository searches must target observed paths and must not treat partial stdout from an rg path error as a complete inventory."
  root_cause: "Two scoped rg commands included a nonexistent directory and a guessed catalog filename."
  canonical_correction: "Enumerate exact paths first or search only confirmed roots/files; repeat the query against implementation_packs and markdown_system without the missing targets."
  regression_test: "The corrected targeted searches and the exact TikTok source acquisition complete without path errors; no claim relies on the partial outputs."
  clean_rebuild_evidence: "Official TikTok archive f809c396 was reacquired at 3515705 bytes/SHA-256 d77e75ba...34cb and inspected independently."
  residual_risk: "TikTok webhook/retrieval remains conditioned until its exact admitted contract and project live probe are implemented."

- failure_id: FAIL-20260904-138
  fingerprint: "portable-archive-omits-current-franchise-and-failure-entrypoints"
  related_failures: ["LIB-FAIL-1872"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: RELEASE_PACKAGING
  violated_invariant: "The portable library must contain every canonical entrypoint and the failure memory required by AGENTS.md."
  root_cause: "CREATE_PORTABLE_ARCHIVE.ps1 kept an older explicit root allowlist and did not include FRANCHISE_ACCELERATOR.md, START_FRANCHISE.md or PROJECT_FAILURE_LESSONS.md."
  canonical_correction: "Add the three current canonical files to the explicit allowlist and generate/verify a fresh V240 ZIP plus sidecar outside the library root."
  regression_test: "Archive creation must pass, the three entries must be present, DISTRIBUTION_SHA256SUMS must match them and the lateral ZIP SHA-256 must verify."
  clean_rebuild_evidence: "The corrected archivator produced 670 entries; the three owners and DISTRIBUTION_SHA256SUMS are present, the internal manifest links them and ZIP SHA-256 ddcd6bda...107a matches its sidecar."
  residual_risk: "A final R1 archive is generated after recording this proof so the distributed failure memory matches the source."

- failure_id: FAIL-20260904-139
  fingerprint: "tiktok-lead-contract-freshness-correction"
  related_failures: ["LIB-FAIL-1873"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: AUTHORITY_FRESHNESS
  violated_invariant: "A current official provider contract must supersede an older broad statement that its payload fields and API version are unknown."
  root_cause: "The re-audit retained the earlier TikTok blocker after the official 2026 API reference exposed the v1.3 lead webhook field table and lead/subscription endpoints."
  canonical_correction: "Record v1.3 request_id/object=1/entry/id/page/campaign/adgroup/ad/create_time/changes/field/value/time and the /lead/get/ plus subscription endpoints; keep implementation blocked because no webhook signature is documented and the current signed official SDK source omits these APIs."
  regression_test: "Current TikTok main remains signed commit f809c396 and exact archive d77e75ba...34cb; targeted source inspection finds no lead/subscription client, so no code is falsely attributed or promoted."
  clean_rebuild_evidence: "Official TikTok docs and help pages were rechecked on 2026-09-04; only the narrow contract claim changes, not the live/admission gate."
  residual_risk: "A provider account/test lead or an immutable official artifact containing the current endpoint contract is still required before a production adapter can be promoted."

- failure_id: FAIL-20260904-140
  fingerprint: "franchise-clean-compose-powershell-stray-quote"
  related_failures: ["LIB-FAIL-1874"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  violated_invariant: "The clean-room composition command must parse before creating a product destination."
  root_cause: "A stray quote followed the observed temporary compositor root assignment."
  canonical_correction: "Reuse the literal observed root without the extra quote and rerun compositor tests before composition."
  regression_test: "The corrected command must pass the compositor suite and materialize exactly the declared 640 files plus its record in a new destination."
  clean_rebuild_evidence: "The corrected literal assignment parsed, the compositor suite printed PASS and the final fresh run materialized 640 files from 59 packs plus its record."
  residual_risk: "None for command parsing; product-runtime gates are tracked separately."

- failure_id: FAIL-20260904-141
  fingerprint: "franchise-compositor-pass-stale-lastexitcode-wrapper"
  related_failures: ["LIB-FAIL-1875", "LIB-FAIL-1836"]
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  violated_invariant: "A successful PowerShell test script must not be converted into a false failure by inherited native-process state."
  root_cause: "The wrapper inspected `$LASTEXITCODE` after a PowerShell script that reports failure through terminating errors and success by normal return; the variable retained an unrelated nonzero native exit code even though the suite printed PASS."
  canonical_correction: "Let terminating errors propagate and use `$?` only for the immediate PowerShell invocation when an explicit check is required; never interpret stale `$LASTEXITCODE` as that script's result."
  regression_test: "From a new empty destination, the compositor suite must print PASS and the next statement must materialize exactly 640 selected files plus MATERIALIZATION_RECORD.md."
  clean_rebuild_evidence: "The wrapper used the immediate PowerShell result, the suite printed PASS and composition continued normally to 59 packs/640 files in a fresh destination."
  residual_risk: "None for PowerShell result handling; product-runtime gates are tracked separately."

- failure_id: FAIL-20260904-142
  fingerprint: "franchise-composition-plan-path-assumed-inside-tool-pack"
  related_failures: ["LIB-FAIL-1876"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: DISCOVERY
  violated_invariant: "The composition runner must use the observed canonical plan path rather than infer a plans subdirectory inside the materialized tool pack."
  root_cause: "The corrected wrapper correctly passed the compositor suite, then addressed `plans/franchise-complete-pack-plan.json` under the temporary tool root even though that tool pack contains only compositor tooling."
  canonical_correction: "Enumerate the library for the exact JSON plan owner, pass that observed literal path, and use another new destination."
  regression_test: "The canonical plan must resolve and materialize its declared 59 packs/640 files plus one record without destination reuse."
  clean_rebuild_evidence: "The observed Markdown owner was passed directly; it resolved its single fenced JSON object and materialized 59 packs/640 files plus MATERIALIZATION_RECORD.md."
  residual_risk: "None for plan resolution; product-runtime gates are tracked separately."

- failure_id: FAIL-20260904-143
  fingerprint: "franchise-plan-discovery-search-expanded-embedded-fixtures"
  related_failures: ["LIB-FAIL-1877", "LIB-FAIL-1844", "LIB-FAIL-1853"]
  recurrence_count: 2
  status: REGRESSION_PROVEN
  severity: LOW
  classification: DISCOVERY
  violated_invariant: "Discovery output must remain complete and bounded; embedded large official fixtures must not enter a broad textual result."
  root_cause: "A multi-root `rg` for generic composition terms matched large embedded pack payloads and exceeded the output budget even though the desired owner had already been identified."
  canonical_correction: "Discard the truncated output as an exhaustive inventory; read only the observed `markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md` header and enumerate tool filenames separately."
  regression_test: "The next reads must be path-exact and line-bounded, with no truncation, before choosing the extraction/composition command."
  clean_rebuild_evidence: "Bounded reads of the exact plan header and `Read-Plan` implementation proved direct Markdown support; no truncated output governed composition."
  residual_risk: "None for this discovery; broad searches remain prohibited by the standing contract."

- failure_id: FAIL-20260904-144
  fingerprint: "franchise-runtime-discovery-foreach-direct-pipeline"
  related_failures: ["LIB-FAIL-1878", "LIB-FAIL-1861"]
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  violated_invariant: "PowerShell discovery commands must parse and must not pipe directly from a `foreach` statement."
  root_cause: "The combined runtime inventory attempted `foreach (...) { ... } | ConvertTo-Json`, which PowerShell parses as an empty pipe element in this statement form."
  canonical_correction: "Assign the foreach output to a typed/ordinary array, then pipe that array to JSON in a separate statement; repeat all reads because the parser prevented the command from running."
  regression_test: "The bounded target inventory and toolchain JSON must both complete without truncation or parser errors."
  clean_rebuild_evidence: "The corrected array form returned a bounded inventory, exact go.mod/package.json and tool paths: Node/pnpm present, Go/PostgreSQL absent from PATH."
  residual_risk: "Exact Go and PostgreSQL binaries must be resolved from the governed workspace runtime before gates."

- failure_id: FAIL-20260904-145
  fingerprint: "workspace-dependency-locator-no-result-within-bounded-wait"
  related_failures: ["LIB-FAIL-1879"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  violated_invariant: "A dependency locator must return actionable runtime paths within a bounded verification interval or yield to a local exact-path fallback."
  root_cause: "The workspace dependency locator produced no result across multiple 30-second waits and was terminated to avoid unbounded token/time cost."
  canonical_correction: "Search existing evidence and the known Codex runtime tree only for exact `go.exe`, `psql.exe` and `pg_isready.exe` paths; validate their versions before use."
  regression_test: "Resolved binaries must exist and report Go 1.26.7 plus the PostgreSQL version used by the library's governed gates."
  clean_rebuild_evidence: "The focal fallback resolved Go `go1.26.7 windows/amd64` with known executable SHA-256 `5463fe58...61fc`, PostgreSQL 18.6, Node 24.14.1 and the installed pnpm path/version."
  residual_risk: "Product gates still must execute; packageManager requires pnpm 11.19.0 rather than the globally installed 11.25.0."

- failure_id: FAIL-20260904-146
  fingerprint: "franchise-complete-plan-missing-go-import-owners"
  related_failures: ["LIB-FAIL-1880"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: CRITICAL
  classification: COMPOSITION
  violated_invariant: "A complete franchise composition must include every internal package imported by its selected commands and application wiring, with exactly one canonical owner."
  root_cause: "The clean 59-pack composition materialized callers of `elite.local/enterprise/internal/agent` and `internal/finops` but did not materialize either package, so `go test ./...` failed at setup across commands and connected runtime packages."
  canonical_correction: "Trace both imports to their exact pack owners, add only the missing owner files/packs without duplicating the connected runtime, update the canonical composition cardinality and rebuild from a new empty destination."
  regression_test: "Clean composition must pass `go test ./...`, `go vet ./...` and `go build ./...`; plan verification must also assert that every local Go import resolves inside the product tree."
  clean_rebuild_evidence: "The canonical plan now selects 61 packs/657 files, materializes both owners, and the clean product passes Go 1.26.7 test/vet/build. `VERIFY_LIBRARY` also passes all 44 profiles with the new local-import closure gate."
  residual_risk: "Database-backed integration, browser, provider/live and production gates remain separate; the local-import defect is closed."

- failure_id: FAIL-20260904-147
  fingerprint: "frontend-verify-nested-pnpm-global-version-drift"
  related_failures: ["LIB-FAIL-1881"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: TOOLCHAIN
  violated_invariant: "Every frontend gate must execute under the packageManager-pinned pnpm 11.19.0, including commands launched by convenience scripts."
  root_cause: "Corepack correctly ran the outer `pnpm@11.19.0 verify`, but that package script invoked the global `pnpm` command (11.25.0), and the version policy failed closed before typecheck/test/build."
  canonical_correction: "Keep the successful frozen/offline install and invoke `typecheck`, `test` and `build` separately through exact `corepack pnpm@11.19.0`, avoiding the nested unpinned convenience command."
  regression_test: "All three commands must report success under the exact pinned package manager; no bypass of the version policy is permitted."
  clean_rebuild_evidence: "Exact Corepack pnpm 11.19.0 then ran typecheck, Vitest 9 files/38 tests and Next.js production build with 19 generated routes; all passed without policy bypass."
  residual_risk: "Browser/runtime provider gates remain separate; package-manager drift is closed."

- failure_id: FAIL-20260904-148
  fingerprint: "postgres-clean-gate-createdb-wrapper-stalled-before-database"
  related_failures: ["LIB-FAIL-1882"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: TOOLING
  violated_invariant: "The isolated PostgreSQL gate must make bounded, observable progress and must not leave a server running after a stalled client wrapper."
  root_cause: "The first clean runner initialized and started PostgreSQL 18.6, then produced no database or client activity and remained blocked around the `createdb.exe` step for multiple bounded waits."
  canonical_correction: "Interrupt the stalled wrapper, stop its exact data directory, then create the database through `psql -d postgres -c CREATE DATABASE` with statement/connect timeouts and execute migrations/tests with per-file progress and bounded timeouts."
  regression_test: "The replacement run must stop its server in `finally`, report exact counts for migrations and tests, and leave no process owning its data directory."
  clean_rebuild_evidence: "A new logfile-based PostgreSQL 18.6 runtime created `elite_franchise`, applied 49/49 migrations, passed 35/35 SQL tests and stopped normally."
  residual_risk: "Go integration suites with TEST_DATABASE_URL remain a separate gate."

- failure_id: FAIL-20260904-149
  fingerprint: "postgres-interrupted-wrapper-stale-postmaster-pid"
  related_failures: ["LIB-FAIL-1883"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: CLEANUP
  violated_invariant: "Cleanup evidence must distinguish a live PostgreSQL process from a stale postmaster.pid after interruption."
  root_cause: "Interrupting the parent wrapper also ended PostgreSQL before the explicit cleanup; `pg_ctl stop` then found a stale PID file and correctly reported that PID 59860 no longer existed."
  canonical_correction: "Verify process ownership and port closure, never reuse this data directory, and make the replacement runner own/stop its server without external interruption."
  regression_test: "The old PID must have no live process/listener; the replacement run's `finally` must complete normally and remove its live postmaster identity."
  clean_rebuild_evidence: "The old owner count was zero and TCP 59817 was closed; the replacement runtime used a different data directory/port and completed its own normal shutdown."
  residual_risk: "None for the interrupted runtime; its stale directory is never reused."

- failure_id: FAIL-20260904-150
  fingerprint: "postgres-pgctl-start-inherited-output-pipe-hang"
  related_failures: ["LIB-FAIL-1884", "LIB-FAIL-1882"]
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: TOOLING
  violated_invariant: "A detached PostgreSQL server must not retain the wrapper's output pipe and prevent the gate from advancing."
  root_cause: "`pg_ctl start` was piped to `Out-Null` without `-l`; on this Windows runtime the ready postgres child inherited the pipe handle, so PowerShell waited even after pg_ctl exited."
  canonical_correction: "Interrupt and stop the exact second data directory; on the replacement use `pg_ctl -l <literal logfile> start` with no pipeline, then verify readiness explicitly before creating the database."
  regression_test: "After the ready probe the wrapper must emit `started`, create the database, report per-stage progress, run migrations/tests and stop cleanly."
  clean_rebuild_evidence: "The replacement used `pg_ctl -l`, emitted every stage, applied 49 migrations, passed 35 SQL tests and stopped cleanly in `finally`."
  residual_risk: "None for pg_ctl handle ownership; future Windows runners must retain the logfile form."

- failure_id: FAIL-20260904-151
  fingerprint: "go-postgres-integration-runner-reintroduced-pgctl-pipeline"
  related_failures: ["LIB-FAIL-1885", "LIB-FAIL-1884"]
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: TOOLING
  violated_invariant: "After proving the Windows pg_ctl no-pipeline rule, no later gate may reintroduce a pipeline around start/stop."
  root_cause: "The integration runner retained `-l` but piped `pg_ctl start` to `Out-Host`; the server reached ready on port 56621 while the wrapper again failed to reach database creation."
  canonical_correction: "Interrupt/quarantine the exact runtime and repeat with a bare pg_ctl invocation—no `Out-Host`, `Out-Null` or pipeline—matching the already proven SQL runner."
  regression_test: "The bare invocation must advance beyond startup, create the database, apply 49 migrations, execute both Go packages with zero skips and stop normally."
  clean_rebuild_evidence: "A replacement PostgreSQL 18.6 runtime used bare pg_ctl invocations, applied all 49 migrations, ran `./internal/platform/postgres` and `./internal/app` with 45 tests, zero skips/failures, and stopped normally."
  residual_risk: "Live load, provider and production gates remain separate; the runner defect is closed."

- failure_id: FAIL-20260904-152
  fingerprint: "franchise-python-tests-run-without-profile-dependency-environment"
  related_failures: ["LIB-FAIL-1886"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: DEPENDENCY
  violated_invariant: "Product adapter tests must execute in a dependency environment built from their exact materialized locks, never from an assumed global Python installation."
  root_cause: "The aggregate runner invoked `google_ads_reporting/test_reporting_query.py` with global Python before creating the adapter's locked environment; import failed on `google.ads.googleads`."
  canonical_correction: "Enumerate exact lock files and supported Python versions, build an isolated temporary environment from the relevant hash-pinned locks, then rerun every Python test; do not install globally or claim the one earlier test as aggregate proof."
  regression_test: "All discovered product Python tests must pass in the isolated locked environment with zero missing imports and a recorded interpreter/dependency identity."
  clean_rebuild_evidence: "Fifteen dependency-free files passed; Google Ads 3, Google Merchant 7, Meta Ads 6, Meta Lead 10 and TikTok 6 tests passed in isolated hash-locked environments with `pip check`. The symlink negative also passed on WSL Linux, eliminating the Windows-only skip."
  residual_risk: "Provider account/live behavior remains separate; the composed Python dependency/test gate is closed."

- failure_id: FAIL-20260904-153
  fingerprint: "google-ads-adapter-sdk-stale-and-transitive-lock-absent"
  related_failures: ["LIB-FAIL-1887"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: AUTHORITY_FRESHNESS
  violated_invariant: "A selected provider adapter must use a current immutable official release and a complete reproducible dependency lock before aggregate execution."
  root_cause: "The franchise-selected Google Ads adapter still fixed official SDK 31.3.0 and only its direct wheel, while official PyPI/GitHub now publish 31.4.0 and the pack explicitly lacked a transitive environment lock."
  canonical_correction: "Admit the official Google 31.4.0 release/tag and PyPI wheel SHA-256 `210a7f6a...94474`, construct/verify a Windows Python 3.14 hash-pinned transitive lock in isolation, update the pack/plan, and rerun adapter plus library gates."
  regression_test: "Fresh materialization must include the lock, install with `--require-hashes`, import the exact SDK version, pass its adapter tests and preserve license/notices and rollback identity."
  clean_rebuild_evidence: "Pack 0.2.0 reconstructs seven files; a fresh venv installed all 24 exact wheels, `pip check`, 3 tests and SDK 31.4.0/API-v25 probe passed; root/Audit passed and the official OSV batch returned 0 vulnerabilities across the combined provider locks."
  residual_risk: "Developer token, OAuth, account hierarchy, quota and live query remain project gates."

- failure_id: FAIL-20260904-154
  fingerprint: "google-ads-lock-tool-discovery-direct-foreach-pipeline"
  related_failures: ["LIB-FAIL-1888", "LIB-FAIL-1878", "LIB-FAIL-1861"]
  recurrence_count: 2
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  violated_invariant: "The already learned PowerShell rule requiring array materialization before piping must apply to every subsequent discovery command."
  root_cause: "The lock-tool inventory again ended a `foreach` statement directly with `| ConvertTo-Json`, producing a parser error before any requested read executed."
  canonical_correction: "Repeat the entire bounded read, assign the two command rows to `$rows`, and serialize `$rows` in a separate statement; use this form for the rest of the task."
  regression_test: "Pack excerpt, existing lock excerpt and tool inventory must all return in one valid bounded command."
  clean_rebuild_evidence: "The corrected array-first discovery returned the requested pack, lock and tool inventory and was used to build the exact 24-wheel environment."
  residual_risk: "None for the discovery syntax; the array-first rule remains mandatory."

- failure_id: FAIL-20260904-155
  fingerprint: "google-ads-lock-pyasn1-hash-transcription-error"
  related_failures: ["LIB-FAIL-1889"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: INTEGRITY
  violated_invariant: "Every dependency URL hash must be copied exactly from the resolver report and verified by a fail-closed installer."
  root_cause: "The manually transferred pyasn1 hash omitted `3d` in the middle: the lock expected `...5233cdcd8d...` while the official wheel and resolver report are `...5233cdcd8d3d...`."
  canonical_correction: "Correct only the pyasn1 hash from the retained pip report and rerun the entire `--require-hashes --no-deps` install in a new virtual environment; never weaken hash enforcement."
  regression_test: "All 24 direct wheel URLs must install under hash enforcement, followed by `pip check`, adapter tests and SDK/API probes."
  clean_rebuild_evidence: "A fresh venv installed all 24 exact URLs with `--require-hashes --no-deps`; `pip check`, 3 adapter tests and the 31.4.0/API-v25 import probe passed."
  residual_risk: "Provider credentials, quota, live account and SCA freshness remain separate gates."

- failure_id: FAIL-20260904-156
  fingerprint: "google-ads-verifier-multi-hunk-context-name-mismatch"
  related_failures: ["LIB-FAIL-1890"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: EDITING
  violated_invariant: "Verifier integration patches must use exact observed variable names and remain atomic on mismatch."
  root_cause: "A multi-hunk patch assumed `amazon-spapi-external-inventory-adapter` while the current network-gate variable uses a different literal destination name."
  canonical_correction: "Retain the atomic rejection, reread the exact network function header, and apply static-audit and network-gate changes as separate narrow patches."
  regression_test: "Both Google Ads roots/locks must appear exactly once in their respective functions and the PowerShell parser plus Audit must pass."
  clean_rebuild_evidence: "Separate narrow patches added the static and network gates exactly once; the PowerShell script parsed and `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit` completed with 151 packs."
  residual_risk: "Network installation remains an explicit AllowNetwork gate; static enforcement is closed."

- failure_id: FAIL-20260904-157
  fingerprint: "composition-tool-pack-owner-path-assumed"
  related_failures: ["LIB-FAIL-1891", "LIB-FAIL-1876"]
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: LOW
  classification: DISCOVERY
  violated_invariant: "Canonical pack-owner paths must be enumerated before targeted reads; a known materialized tool path does not prove an assumed source filename."
  root_cause: "A bounded command correctly found the existing materialized compositor but also queried a guessed `markdown_system/PROJECT_COMPOSITION_TOOL_PACK.md` path that does not exist."
  canonical_correction: "Use the already observed materialized compositor for the clean run and enumerate source pack filenames with `rg --files` before any later source read."
  regression_test: "The next composition must use the exact observed tool path, produce a fresh destination and report 61 packs/659 files."
  clean_rebuild_evidence: "The exact observed compositor materialized the current canonical plan into a fresh UUID destination: 61 packs and 659 product files."
  residual_risk: "Runtime/test gates on the fresh product remain separate."

- failure_id: FAIL-20260904-158
  fingerprint: "powershell-get-help-literalpath-unsupported"
  related_failures: ["LIB-FAIL-1892"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING
  violated_invariant: "Help/discovery commands must use parameters actually supported by the installed PowerShell command."
  root_cause: "`Get-Help` was invoked with `-LiteralPath`, which is not a parameter of this host's command."
  canonical_correction: "Read the already trusted local script's parameter header with `Get-Content -LiteralPath`, or pass its exact name to `Get-Help -Name`; do not infer parameters."
  regression_test: "The compositor parameter block must be observed before invocation and the clean composition must pass."
  clean_rebuild_evidence: "`Get-Content -LiteralPath` exposed the exact three-parameter header; invocation with those parameters produced the clean 61/659 composition."
  residual_risk: "None for parameter discovery."

- failure_id: FAIL-20260904-159
  fingerprint: "franchise-python-production-admission-one-skipped-test"
  related_failures: ["LIB-FAIL-1893"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: TEST_EVIDENCE
  violated_invariant: "The composed product's aggregate test claim requires every discovered test to be classified; an unexplained skip cannot be counted as full evidence."
  root_cause: "The dependency-free Python run passed 15 files, but `test_validate_production_admission.py` reported one skipped test without its reason in default verbosity."
  canonical_correction: "Rerun that exact file verbosely, inspect the declared skip condition and either satisfy the local prerequisite or retain it as an explicit target-only gate."
  regression_test: "Every test in the file must be PASS or have a named, justified target/environment condition; no silent skip may support readiness."
  clean_rebuild_evidence: "Verbose Windows execution identified only host symlink privilege; the same final composed bytes ran under Ubuntu/WSL and all 9/9 tests passed, including symlink rejection."
  residual_risk: "None for the test semantics; Windows hosts without symlink privilege retain the explicit local skip while Linux production evidence is proven."

- failure_id: FAIL-20260904-160
  fingerprint: "python-unittest-package-import-from-parent-without-package"
  related_failures: ["LIB-FAIL-1894"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TEST_TOOLING
  violated_invariant: "A standalone Python test must be invoked from the import root assumed by its source."
  root_cause: "The verbose rerun used dotted-module notation from the product root, but the test imports a sibling module and its directory is not a Python package on that invocation path."
  canonical_correction: "Run unittest from the exact `production_admission_gate` directory with the observed filename."
  regression_test: "Verbose output must enumerate all nine tests and disclose the exact skipped test/reason without import errors."
  clean_rebuild_evidence: "The test was rerun from its exact directory: Windows disclosed the skip reason and Ubuntu/WSL executed all 9/9 including symlink rejection."
  residual_risk: "None for import-root selection."

- failure_id: FAIL-20260904-161
  fingerprint: "franchise-tiktok-adapter-lock-omits-official-sdk-import"
  related_failures: ["LIB-FAIL-1895", "LIB-FAIL-1886"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: DEPENDENCY
  violated_invariant: "Every selected provider adapter must reconstruct all imports used by its tests from its materialized files and exact locks."
  root_cause: "The fresh TikTok adapter venv installed all four entries in `requirements.lock` and passed `pip check`, but its official-SDK contract test could not import `business_api_client`; the lock omits that package/runtime source."
  canonical_correction: "Trace the exact official TikTok SDK artifact already admitted by the pack, determine whether it is materialized outside the adapter directory or absent from the product plan, and add only the exact licensed/locked source required by the adapter."
  regression_test: "A fresh isolated venv/product composition must install/materialize the official SDK and pass all six TikTok reporting tests plus identity/hash probes."
  clean_rebuild_evidence: "The official 1.1.3 PyPI wheel was added to the five-entry lock; a fresh final composition installed it with hash enforcement, passed pip check, 6/6 tests and ReportingApi/version probes."
  residual_risk: "TikTok account/scope/query/live reconciliation remain project gates."

- failure_id: FAIL-20260904-162
  fingerprint: "official-tiktok-sdk-pinned-commit-missing-imported-model"
  related_failures: ["LIB-FAIL-1896", "LIB-FAIL-1895"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: UPSTREAM_CODE_DEFECT
  violated_invariant: "A signed official source snapshot is not reusable merely because its receipt and critical hashes match; its actual import graph must execute."
  root_cause: "After exact source/receipt verification and PYTHONPATH binding, official commit `f809c396...` failed importing its own package: `models/__init__.py` references `dmpcustom_audiencerulecreate_rule_spec_exclusion_rule_set_filter_set`, but that module is absent from the snapshot."
  canonical_correction: "Reopen official-source freshness and inspect official TikTok history/releases for a signed revision that closes the import graph; reject the current snapshot if no official fix exists. Do not author or silently patch a model."
  regression_test: "The replacement official snapshot must pass complete package import, ReportingApi import, exact-source hashes/license, six adapter tests and a fresh lock/source receipt."
  clean_rebuild_evidence: "The broken GitHub snapshot was removed as runtime authority and replaced by TikTok's official PyPI 1.1.3 wheel; clean import, endpoint/header contract and 6/6 adapter tests passed."
  residual_risk: "The rejected source snapshot remains historical evidence and cannot be reintroduced."

- failure_id: FAIL-20260904-163
  fingerprint: "tiktok-adapter-freshness-missed-official-pypi-1-1-3"
  related_failures: ["LIB-FAIL-1897", "LIB-FAIL-1896"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: AUTHORITY_FRESHNESS
  violated_invariant: "A provider admission must inspect all official distribution channels named by the provider before claiming that no official package/release exists."
  root_cause: "The pack said TikTok publishes no official Python wheel, but TikTok's own README recommends `tiktok-business-api-sdk-official` and PyPI currently exposes author `TikTok Pte. Ltd.` version 1.1.3, uploaded 2026-02-27 with wheel SHA-256 `663b4a...f33b7`."
  canonical_correction: "Test the exact official 1.1.3 wheel in isolation, verify package metadata/license/import/API surface, replace the broken GitHub runtime only if all gates pass, and retain commit source solely as a provenance/reference input where still material."
  regression_test: "Hash-required install of the official wheel plus all transitive dependencies must pass `pip check`, package/ReportingApi import, adapter tests and SCA."
  clean_rebuild_evidence: "Exact wheel 1.1.3 installed in two fresh venvs including the final composition; pip check, import, 6 tests, root/Audit and an OSV batch covering 41 unique provider packages with 0 vulnerabilities passed."
  residual_risk: "Upstream freshness must be rechecked on future releases; live credentials remain project-only."

- failure_id: FAIL-20260904-164
  fingerprint: "apply-patch-same-file-delete-add-in-one-patch"
  related_failures: ["LIB-FAIL-1898"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: EDITING
  violated_invariant: "A source-edit patch must use an operation sequence supported by apply_patch and remain atomic on validation failure."
  root_cause: "The TikTok staging patch attempted both Delete File and Add File for the same README path within one patch; apply_patch rejected the full patch."
  canonical_correction: "Use Update File for an existing file and separate the SDK lock rename/content operations from ordinary updates."
  regression_test: "All seven resulting staging files must match intended content, tests must pass, and canonical round-trip hashes must match after synchronization."
  clean_rebuild_evidence: "The edits were split into supported updates/rename operations; eight blocks synchronized, clean materialization 61/659 and both root/Audit gates passed."
  residual_risk: "None for patch operation sequencing."

- failure_id: FAIL-20260904-165
  fingerprint: "global-provider-verifier-omits-meta-lead-reconciliation"
  related_failures: ["LIB-FAIL-1899", "LIB-FAIL-1886"]
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: VERIFICATION_COVERAGE
  violated_invariant: "Every provider adapter selected by the franchise profile must be materialized and runtime-tested by the global provider gate."
  root_cause: "The composed-product manual suite proved Meta Lead reconciliation 10/10, but inspection of `Invoke-ProviderAdapterAudit`/`Invoke-NetworkProviderAdapterGates` found no Meta Lead root, lock or runtime invocation."
  canonical_correction: "Add static lock/artifact checks and an isolated hash-locked runtime test for `PYTHON_META_LEAD_RECONCILIATION_ADAPTER`, then add one official OSV batch over every provider venv in the network gate."
  regression_test: "Global AllowNetwork execution must report Meta Lead 10/10 and zero OSV vulnerabilities, while Audit must enforce its exact 18-wheel/source identity."
  clean_rebuild_evidence: "Final AllowNetwork Audit materialized and exercised Meta Lead under its exact 18-wheel environment and completed with `provider_adapters=15`; the common OSV batch remained fail-closed."
  residual_risk: "Live Meta credentials/webhooks and tenant traffic remain project gates, not library-runtime claims."

- failure_id: FAIL-20260904-166
  fingerprint: "provider-adapter-count-hardcoded-after-meta-lead-gate"
  related_failures: ["LIB-FAIL-1900", "LIB-FAIL-1899"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: INVENTORY
  violated_invariant: "A verifier's reported adapter cardinality must match every adapter it actually materializes."
  root_cause: "After adding Meta Lead static enforcement, Audit still reported 14 because `providerAdapterCount` is hardcoded and was not updated with the new gate."
  canonical_correction: "Set the governed count to 15 and rerun Audit/AllowNetwork; later replace hardcoding with derived cardinality if the adapter registry is centralized."
  regression_test: "Final verifier output must report `provider_adapters=15` and include the Meta Lead runtime PASS."
  clean_rebuild_evidence: "Final AllowNetwork Audit reported `provider_adapters=15` after executing the added Meta Lead gate."
  residual_risk: "The count is coherent with current gates; a future registry refactor should derive it automatically."

- failure_id: FAIL-20260904-167
  fingerprint: "allow-network-audit-launched-with-python314-against-markitdown-py312-lock"
  related_failures: ["LIB-FAIL-1901"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: EXECUTION_ENVIRONMENT
  violated_invariant: "A frozen artifact lock must be verified with its exact governed interpreter and platform."
  root_cause: "The final AllowNetwork audit was launched with CPython 3.14 to satisfy provider adapters, but the MarkItDown lock is intentionally fixed to CPython 3.12 Windows x86-64."
  canonical_correction: "Rerun the unchanged verifier with the installed CPython 3.12.14 x86-64 runtime; preserve all MarkItDown version/platform checks and let provider adapters create their own exact environments."
  regression_test: "AllowNetwork Audit must pass the MarkItDown frozen install/import gate and later report all provider and OSV gates without relaxing the lock."
  clean_rebuild_evidence: "The unchanged AllowNetwork Audit reran with CPython 3.12.14 x86-64, installed all 44 MarkItDown lock distributions, passed `pip check`, tests and all subsequent gates, ending in `VERIFY_EXECUTABLE_LIBRARY_PASS ... provider_adapters=15`."
  residual_risk: "Another host must reacquire the exact governed runtime/artifacts; this host evidence is complete."

- failure_id: FAIL-20260904-168
  fingerprint: "broad-temp-rg-crossed-access-denied-directories"
  related_failures: ["LIB-FAIL-1902"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: DISCOVERY
  violated_invariant: "A discovery result must be complete for its declared search scope before it can support an absence claim."
  root_cause: "A broad `rg --files` over the whole user temp directory crossed four unrelated access-denied directories and returned only a partial inventory."
  canonical_correction: "Discard the absence result and enumerate only observed `elite-*` temporary roots with accessible literal paths before locating .NET and generated ARCA artifacts."
  regression_test: "The scoped enumeration must complete without access errors and every selected executable/receipt must be identity-checked before use."
  clean_rebuild_evidence: "Scoped enumeration completed without access errors, identified dotnet 10.0.400 and eleven exact ARCA generation receipts under accessible `elite-*` roots."
  residual_risk: "Each selected receipt still requires its own content/hash gate; inventory alone is not admission."

- failure_id: FAIL-20260904-169
  fingerprint: "stale-dotnet-v125-candidate-has-no-sdk"
  related_failures: ["LIB-FAIL-1903"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: EXECUTION_ENVIRONMENT
  violated_invariant: "A candidate toolchain path is not usable until its exact version command succeeds."
  root_cause: "An older `elite-dotnet-v125/sdk/dotnet.exe` host executable remained on disk without a discoverable SDK."
  canonical_correction: "Reject that candidate and select only `elite-dotnet-v127/dotnet.exe` after it reports the governed SDK 10.0.400."
  regression_test: "Every ARCA gate must receive the literal v127 executable and independently enforce `dotnet --version == 10.0.400`."
  clean_rebuild_evidence: "The v127 path returned 10.0.400; v125 returned the official no-SDK diagnostic and is excluded."
  residual_risk: "None for selection while the ARCA scripts retain their own exact version checks."

- failure_id: FAIL-20260904-170
  fingerprint: "postgres-v169-runtime-missing-pg-dump-for-restore-drill"
  related_failures: ["LIB-FAIL-1904"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: EXECUTION_ENVIRONMENT
  violated_invariant: "A backup/restore proof requires the exact PostgreSQL server and complete native backup toolchain, including `pg_dump` and `pg_restore`."
  root_cause: "The selected PostgreSQL 18.6 runtime could initialize, start and migrate the database but its bin directory lacks `pg_dump.exe`; the product gate rejected it before creating a backup."
  canonical_correction: "Locate or acquire the official complete PostgreSQL 18.6 Windows distribution, verify every required executable/version, then repeat all 49 migrations plus backup/hash/restore/invariants from a fresh cluster."
  regression_test: "`Test-BackupRestore.ps1` must return PASS, create a nonempty custom dump and manifest, restore into a new database, pass invariants and remove the drill database."
  clean_rebuild_evidence: "A new PostgreSQL 18.6 cluster using the complete official toolchain applied all 49 up migrations, produced a custom dump plus SHA-256 manifest (`1bee6e9a...304e8`), restored it into a separate database, passed platform invariants, removed the drill database and shut down normally."
  residual_risk: "Project-specific PITR, retention, encrypted remote storage and operational RTO/RPO still require the selected production provider; this local logical restore claim is proven."

- failure_id: FAIL-20260904-171
  fingerprint: "global-audit-temp-pack-cleaned-before-target-browser-test"
  related_failures: ["LIB-FAIL-1905"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TEST_ORCHESTRATION
  violated_invariant: "A target test must materialize its canonical pack independently instead of depending on a verifier's disposable workspace."
  root_cause: "The global verifier correctly removed its `elite-executable-verification-*` directory at completion, while a follow-up command assumed the Playwright pack still existed there."
  canonical_correction: "Locate the canonical Playwright pack, materialize it into a new explicit test root and run its governed target mode against the live final product."
  regression_test: "The independently materialized pack must install from its exact lock and pass browser target tests without referring to the deleted audit directory."
  clean_rebuild_evidence: "Fresh canonical materialization installed Microsoft Playwright 1.62.1 with pnpm 11.19.0 and passed runtime 4/4 plus final-product target 8/8."
  residual_risk: "Authenticated/project-provider journeys and real assistive-technology testing remain target-project gates."

- failure_id: FAIL-20260904-172
  fingerprint: "playwright-materializer-wrapper-read-inherited-lastexitcode"
  related_failures: ["LIB-FAIL-1906", "LIB-FAIL-1875"]
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TEST_ORCHESTRATION
  violated_invariant: "A PowerShell script that signals failure by terminating exception must not be judged by an unrelated inherited native exit code."
  root_cause: "The wrapper checked `$LASTEXITCODE` immediately after `materialize_markdown_pack.ps1`; the script succeeded and wrote ten files, but the variable still held an earlier native command result."
  canonical_correction: "Use the script's terminating error semantics and `$?` immediately if an explicit success check is required; never read inherited `$LASTEXITCODE` at this boundary."
  regression_test: "Fresh materialization must proceed to exact pnpm installation and eight target browser tests."
  clean_rebuild_evidence: "A new UUID destination used terminating-error semantics and continued through exact installation and 12 total browser tests."
  residual_risk: "None for this wrapper boundary."

- failure_id: FAIL-20260904-173
  fingerprint: "isolated-corepack-shim-resolved-pnpm-11-23-not-11-19"
  related_failures: ["LIB-FAIL-1907", "LIB-FAIL-1881"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: DEPENDENCY_TOOLCHAIN
  violated_invariant: "The browser gate must install and execute with the exact package manager declared by `packageManager` and its lock evidence."
  root_cause: "A new Corepack shim directory still used Corepack's shared known-good pnpm 11.23.0 because its home/activation state was not isolated and pinned."
  canonical_correction: "Create a dedicated `COREPACK_HOME`, explicitly install/activate pnpm 11.19.0 there, prepend only its shim and verify the version before installation."
  regression_test: "The clean target run must print `PNPM=11.19.0`, frozen offline install, Playwright 1.62.1 runtime 4/4 and product target 8/8."
  clean_rebuild_evidence: "Dedicated Corepack home and shim printed pnpm 11.19.0, then frozen offline install and all Playwright gates passed."
  residual_risk: "Other hosts must repeat exact package-manager activation."

- failure_id: FAIL-20260904-174
  fingerprint: "powershell-home-case-insensitive-variable-used-for-corepack-cache"
  related_failures: ["LIB-FAIL-1908"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: TEST_ORCHESTRATION
  violated_invariant: "Task scripts must never repurpose PowerShell's reserved `$HOME` variable, including case-insensitive aliases such as `$home`."
  root_cause: "The Corepack probe named a local path variable `$home`; PowerShell rejected the assignment and retained the user's real home path, so the intended cache isolation was not demonstrated."
  canonical_correction: "Use a task-specific `$corepackCache` variable and UUID path, set `COREPACK_HOME` only to that verified path, and do not attempt destructive cleanup of ambiguous global state."
  regression_test: "A new process must show isolated `COREPACK_HOME`, select pnpm 11.19.0 through its own shim, and complete the target gate."
  clean_rebuild_evidence: "A fresh process used `$corepackCache` under an explicit UUID path, printed that exact `COREPACK_HOME`, selected pnpm 11.19.0 and passed all gates."
  residual_risk: "Corepack may have refreshed its user-level cache/known-good record; no user files are removed because the exact mutation boundary is ambiguous."

- failure_id: FAIL-20260904-175
  fingerprint: "playwright-install-run-at-destination-parent-not-pack-subdirectory"
  related_failures: ["LIB-FAIL-1909"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TEST_ORCHESTRATION
  violated_invariant: "Commands must run from the exact pack root produced by its declared file manifest."
  root_cause: "The materializer wrote `microsoft_playwright_browser_gate/*` below the destination, while pnpm was invoked at the destination parent."
  canonical_correction: "Resolve and verify the nested pack root containing `package.json`, then run the exact isolated install and verifier there."
  regression_test: "The corrected root must retain pnpm 11.19.0 and complete runtime 4/4 plus product target 8/8."
  clean_rebuild_evidence: "The new run resolved `materialized/microsoft_playwright_browser_gate/package.json`, installed 3 exact packages and passed runtime 4/4 plus product target 8/8."
  residual_risk: "None for pack-root resolution."

- failure_id: FAIL-20260904-176
  fingerprint: "lighthouse-post-pass-summary-filename-assumed"
  related_failures: ["LIB-FAIL-1910"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: EVIDENCE_RETRIEVAL
  violated_invariant: "Evidence filenames must be discovered from actual output or the runner contract, never guessed after a successful gate."
  root_cause: "The Lighthouse target printed PASS for five runs, but the wrapper then requested a nonexistent `quality-summary.json` path."
  canonical_correction: "Enumerate the explicit results directory, locate the actual summary/receipt, validate its schema and five report references, then rerun only if the evidence is incomplete."
  regression_test: "All five raw reports and one governed summary must exist and agree with the PASS scores."
  clean_rebuild_evidence: "`quality-gate-summary.json` and five `lighthouse-run-*.json` files were enumerated; all report Lighthouse 13.4.1, exact target URL, 155 audits/run, zero warnings, performance minimum/median 0.99 and all other category minimums/medians 1."
  residual_risk: "This proves automated local lab quality only, exactly as the summary limits state."

- failure_id: FAIL-20260904-177
  fingerprint: "web-tool-rejected-direct-github-api-tiktok-tree"
  related_failures: ["LIB-FAIL-1911"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: OFFICIAL_SOURCE_ACQUISITION
  violated_invariant: "A retrieval-tool restriction must not be misinterpreted as absence of official source or endpoint support."
  root_cause: "The web safety layer rejected direct opens of GitHub commit/tree API URLs for the official TikTok repository."
  canonical_correction: "Use HTTPS Git against the same official repository, pin the observed commit SHA, inspect the exact tree locally and retain license/commit receipts."
  regression_test: "The Git acquisition must resolve one exact commit and a complete searchable tree without relying on a mobile branch during admission."
  clean_rebuild_evidence: "HTTPS Git resolved and detached exactly `f809c396520df2d7b201a9ccc5378d822b728ed3`; the complete tree was searched and still contains no Lead retrieval/subscription client. The dynamic official portal was then inspected read-only and supplied exact v1.3 contracts."
  residual_risk: "An adapter must remain explicit ADAPTED/AUTHORED glue over official primitives; it cannot be presented as a TikTok-published Lead SDK."
```

- failure_id: FAIL-20260904-178
  fingerprint: "pack-sync-tool-path-assumed-outside-markdown-system"
  related_failures: ["LIB-FAIL-1912", "LIB-FAIL-1856"]
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING_DISCOVERY
  violated_invariant: "Canonical tooling paths must be observed from the workspace inventory before invocation."
  root_cause: "The continuation state named `tools/update_pack_from_tree.ps1`, but the canonical script lives at `markdown_system/update_pack_from_tree.ps1`."
  canonical_correction: "Resolve the tool with `rg --files`, use the observed literal path and retain the failed read only as a discovery failure."
  regression_test: "The canonical updater and its test owner must both resolve from the workspace inventory before any pack synchronization."
  clean_rebuild_evidence: "`rg --files` resolved `markdown_system/update_pack_from_tree.ps1` and `markdown_system/test_update_pack_from_tree.ps1`; no pack bytes were changed by the failed read."
  residual_risk: "None for path resolution; pack synchronization still requires its own round-trip tests."

- failure_id: FAIL-20260904-179
  fingerprint: "tiktok-lead-staging-pycompile-generated-cache"
  related_failures: ["LIB-FAIL-1913"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: BUILD_HYGIENE
  violated_invariant: "A pack staging tree may contain only files declared by its exact source manifest."
  root_cause: "`python -m py_compile` wrote three interpreter-specific `.pyc` files below staging before the manifest inventory was calculated."
  canonical_correction: "Delete only the verified generated `__pycache__` directory, enumerate the eight intended source files and ensure clean materialization never includes bytecode caches."
  regression_test: "The source inventory and materialized inventory must each contain exactly eight declared files and no `__pycache__` or `.pyc`."
  clean_rebuild_evidence: "Canonical pack materialized exactly 8/8 declared files with matching SHA-256 and no cache; its 10 tests passed from that clean destination."
  residual_risk: "None after the manifest and destination inventories agree."

- failure_id: FAIL-20260904-180
  fingerprint: "policy-rejected-recursive-staging-cache-cleanup"
  related_failures: ["LIB-FAIL-1914"]
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: LOW
  classification: SAFE_CLEANUP
  violated_invariant: "Generated staging artifacts must be removed through an action accepted by the host safety policy and confined to exact targets."
  root_cause: "The cleanup used recursive deletion for a directory containing only three already-enumerated bytecode files; the host rejected the whole command before execution."
  canonical_correction: "Build the canonical pack strictly from the eight-file allowlist, verify its clean destination, then move the entire disposable staging directory to a verified UUID temp destination instead of deleting individual bytecode files."
  regression_test: "All four exact paths become absent and the eight-file allowlist remains intact."
  clean_rebuild_evidence: "The complete `.pack_staging` root was moved with one exact PowerShell operation to validated UUID temp `elite-pack-staging-retained-806ec0a139c446ff880cc150330b8aec`; source is absent and canonical destination remained 8/8 clean."
  residual_risk: "None if the explicit paths are absent and no source file changes."

- failure_id: FAIL-20260904-181
  fingerprint: "apply-patch-cannot-delete-binary-pyc"
  related_failures: ["LIB-FAIL-1915"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: SAFE_CLEANUP
  violated_invariant: "A selected cleanup mechanism must support the exact artifact type before its result is assumed."
  root_cause: "The transactional patcher rejected the binary `.pyc` because it requires UTF-8 input; the patch was atomic and did not change logs or files."
  canonical_correction: "Build from the explicit eight-file allowlist, prove a clean destination and move the disposable staging directory to a verified UUID temp destination."
  regression_test: "Canonical pack and materialized inventories contain exactly eight files; no `.pack_staging` remains below the library root."
  clean_rebuild_evidence: "Pack construction used the eight-text allowlist, clean materialization/tests passed and the polluted staging root was moved outside the workspace without deleting user data."
  residual_risk: "None after the workspace staging path is absent."

- failure_id: FAIL-20260904-182
  fingerprint: "tiktok-lead-pack-large-patch-stale-ledger-context"
  related_failures: ["LIB-FAIL-1916", "LIB-FAIL-1849"]
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: LOW
  classification: PATCH_ORCHESTRATION
  violated_invariant: "A multi-owner patch must use freshly observed exact anchors or be split before it can govern canonical files."
  root_cause: "The pack-creation patch also expected a ledger row variant that had never been committed because the preceding binary-delete patch failed atomically."
  canonical_correction: "Read the exact current rows, update failure records separately and add the new pack in a standalone patch."
  regression_test: "The log patch and pack patch each apply independently; the new pack path exists only after the latter succeeds."
  clean_rebuild_evidence: "The rejected patch left the pack absent; current contexts were read and the log correction now applies independently."
  residual_risk: "None for atomicity; pack verification remains pending."

- failure_id: FAIL-20260904-183
  fingerprint: "tiktok-lead-roundtrip-assumed-materialization-record"
  related_failures: ["LIB-FAIL-1917"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TEST_ORCHESTRATION
  violated_invariant: "Round-trip cardinality must follow the actual materializer contract, not a record emitted by a different compositor."
  root_cause: "The wrapper expected eight declared files plus `MATERIALIZATION_RECORD.md`, but `materialize_markdown_pack.ps1` intentionally emits only declared files."
  canonical_correction: "Require exactly the eight declared files, compare each source/output SHA-256 and reject any extra cache or undeclared path."
  regression_test: "Destination count is 8, allowlist equality and eight hash comparisons pass, then tests run from that destination."
  clean_rebuild_evidence: "Corrected comparison proved exact inventory 8/8, SHA-256 8/8, frozen install, pip check and 10/10 tests from the materialized destination."
  residual_risk: "None after corrected comparison and materialized tests pass."

- failure_id: FAIL-20260904-184
  fingerprint: "failure-close-patch-ledger-md-md-typo"
  related_failures: ["LIB-FAIL-1918"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: PATCH_ORCHESTRATION
  violated_invariant: "Canonical owner paths must be copied from the observed inventory and non-empty hunks must be supplied."
  root_cause: "A close-status patch appended `.md` twice to the ledger path and included an empty update hunk, so the patcher rejected it atomically."
  canonical_correction: "Apply the project status changes and ledger status changes in separate patches using the exact observed paths."
  regression_test: "Each intended status/evidence string is read back from its canonical owner and the nonexistent `.md.md` path remains absent."
  clean_rebuild_evidence: "The rejection was atomic; the corrected project-owner patch follows this record and the ledger owner is patched separately."
  residual_risk: "None after read-back."

- failure_id: FAIL-20260904-185
  fingerprint: "tiktok-lead-profile-validator-not-bound-to-client"
  related_failures: ["LIB-FAIL-1919"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: HIGH
  classification: AUTHORIZATION_GATE
  violated_invariant: "A provider adapter may not make a live-capable call unless its required access/terms/delivery/idempotency/reconciliation profile is enforced on the execution path."
  root_cause: "The first adapter exposed `validate_profile()` but `TikTokLeadClient` did not invoke it, allowing callers to bypass the documented gate."
  canonical_correction: "Require the profile in `TikTokLeadClient`, validate it in `__post_init__` and add a test proving blocked construction and successful proven construction."
  regression_test: "No client instance can be created from the template/blocked profile; all existing endpoint tests construct through a fully proven profile."
  clean_rebuild_evidence: "Constructor profile-bound, materialización limpia, 16/16 Python frozen y Audit global con provider_adapters=16 PASS."
  residual_risk: "Live proofs remain project-bound and cannot be manufactured by the library."

- failure_id: FAIL-20260904-186
  fingerprint: "v242-provenance-ledger-pre-pack-counts"
  related_failures: ["LIB-FAIL-1920"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: LOW
  classification: PROVENANCE_LEDGER
  violated_invariant: "Every materialization block must be reflected exactly in the canonical provenance count before root verification can pass."
  root_cause: "The eight-file TikTok Lead pack was admitted before `THIRD_PARTY_NOTICES.md` changed from V241 counts."
  canonical_correction: "Publish only the verifier-computed totals 1116 AUTHORED, 127 ADAPTED, 105 VERBATIM, 1348 total and explain the new 7/1 boundary."
  regression_test: "Root verifier independently recalculates and accepts the exact provenance marker."
  clean_rebuild_evidence: "Root final recalculó 1356 bloques y aceptó el marker post-V245 1123 AUTHORED/128 ADAPTED/105 VERBATIM."
  residual_risk: "None after independent recount; external dependency licenses remain governed by the pack locks."

- failure_id: FAIL-20260904-187
  fingerprint: "v242-franchise-inspection-assumed-migrations-root"
  related_failures: ["LIB-FAIL-1921", "LIB-FAIL-1871"]
  recurrence_count: 1
  status: REGRESSION_PROVEN
  severity: LOW
  classification: TOOLING_DISCOVERY
  violated_invariant: "Inspection commands must enumerate actual composed paths before passing multiple roots to a search."
  root_cause: "The post-composition search included `$productRoot/migrations`, but migrations live under another observed directory; `rg` returned useful partial output and exit 2."
  canonical_correction: "Retain the independently printed 62/667 composition PASS, discard the partial discovery as exhaustive evidence, enumerate the product tree and search only literal existing roots."
  regression_test: "Every searched root passes `Test-Path`; the resulting inventory identifies canonical owners without stderr."
  clean_rebuild_evidence: "Inventario literal confirmó `db/migrations`; composición limpia 63/675 y 49 migraciones PostgreSQL 18.6 PASS."
  residual_risk: "None after bounded search; no source was changed by the read-only command."

- failure_id: FAIL-20260904-188
  fingerprint: "tiktok-postgres-test-assumed-integration-helpers"
  related_failures: ["LIB-FAIL-1922"]
  recurrence_count: 0
  status: REGRESSION_PROVEN
  severity: MEDIUM
  classification: TEST_INTEGRATION
  violated_invariant: "A new integration test must compile against helpers actually owned by the composed package before its pack is admitted."
  root_cause: "`tiktok_import_integration_test.go` referenced `integrationPool` and `createLeadFixtureTenant`, names inferred from another test style but absent in `internal/platform/postgres`."
  canonical_correction: "Inspect the package's actual PostgreSQL fixture helpers, reuse or define uniquely named local helpers, format and rerun full Go test/vet/build."
  regression_test: "The entire `internal/platform/postgres` package compiles with and without DATABASE_URL, then the TikTok test runs against real PostgreSQL without skip."
  clean_rebuild_evidence: "Helpers locales observados; full Go test/vet/build PASS y PostgreSQL 18.6 ejecutó 49 migraciones, test TikTok focal y 38/38 package sin skips."
  residual_risk: "Live TikTok permanece condicionado, no el owner PostgreSQL local."

### FAIL-20260904-189 — formatter Go no estaba en PATH

- `status`: `REGRESSION_PROVEN`
- `severity`: `LOW`
- `fingerprint`: `tiktok-import-gofmt-host-path-missing`
- `classification`: `TOOLCHAIN_ENVIRONMENT`
- `violated_invariant`: "Todo gate debe invocar el toolchain exacto ya descubierto, sin asumir disponibilidad global."
- `root_cause`: "La edición focal invocó `gofmt` por nombre aunque el host sólo conserva el toolchain Go fijado bajo un path temporal conocido."
- `canonical_correction`: "Invocar `gofmt.exe` y los gates con el path literal del toolchain Go 1.25.7 ya verificado."
- `regression_test`: "Formatear, materializar y ejecutar tests/vet/build desde el mismo toolchain exacto."
- `clean_rebuild_evidence`: "`gofmt.exe` del toolchain fijado formateó ambos archivos; materialización limpia, full Go test/vet/build PASS."
- `residual_risk`: "Ningún archivo canónico se modificó por la invocación fallida; el staging conserva el patch sin formatear hasta aplicar la corrección."

### FAIL-20260904-190 — updater de packs requiere allowlist explícita

- `status`: `REGRESSION_PROVEN`
- `severity`: `LOW`
- `fingerprint`: `tiktok-import-pack-updater-files-required`
- `classification`: `TOOLING_CONTRACT`
- `violated_invariant`: "Una herramienta canónica se invoca desde su interfaz observada, incluidos todos sus parámetros obligatorios."
- `root_cause`: "Se reutilizó el updater correcto pero se omitió su parámetro obligatorio `Files`."
- `canonical_correction`: "Leer el header y pasar sólo los dos archivos modificados como allowlist literal."
- `regression_test`: "El pack materializado debe reproducir byte a byte ambos archivos y superar tests/vet/build."
- `clean_rebuild_evidence`: "Updater sincronizó sólo los dos paths; root verificó hashes y la composición ejecutó sus regresiones."
- `residual_risk`: "El updater falló antes de escribir el pack; el formato del staging sí quedó corregido."

### FAIL-20260904-191 — suite PostgreSQL usó dos nombres de URL

- `status`: `REGRESSION_PROVEN`
- `severity`: `MEDIUM`
- `fingerprint`: `postgres-suite-database-url-alias-skip`
- `classification`: `TEST_ENVIRONMENT`
- `violated_invariant`: "Una prueba integral real no puede contener skips ambientales no clasificados."
- `root_cause`: "El nuevo importador consume `TEST_DATABASE_URL`, mientras una suite heredada de provider webhooks aún consume `DATABASE_URL`; el primer runtime sólo definió la primera."
- `canonical_correction`: "Reiniciar el mismo cluster migrado, definir ambas variables con la misma URL aislada y repetir el package completo comprobando ausencia de SKIP."
- `regression_test`: "Todos los tests de `internal/platform/postgres`, incluido provider webhook y TikTok import, deben ejecutarse contra PostgreSQL 18.6."
- `clean_rebuild_evidence`: "Rerun con ambos aliases ejecutó 38/38 pruebas PostgreSQL reales, cero SKIP."
- `residual_risk`: "Ninguno para el test TikTok focal; la afirmación de suite PostgreSQL completa queda bloqueada hasta el rerun sin skip."

### FAIL-20260904-192 — wheel TikTok 1.1.3 expone `__version__` 1.2.1

- `status`: `REGRESSION_PROVEN`
- `severity`: `HIGH`
- `fingerprint`: `tiktok-official-wheel-distribution-module-version-divergence`
- `classification`: `UPSTREAM_PROVENANCE`
- `violated_invariant`: "La identidad de un artefacto oficial debe verificarse sin inferir equivalencia entre metadata de distribución y constantes internas contradictorias."
- `root_cause`: "El wheel oficial hash-pinned publicado como distribución 1.1.3 instaló correctamente, pero `business_api_client.__version__` respondió 1.2.1."
- `canonical_correction`: "Inspeccionar metadata RECORD/METADATA y código exactos del wheel; gobernar por filename, SHA-256 y metadata de distribución o rechazarlo si la contradicción afecta la identidad ejecutable."
- `regression_test`: "El gate debe verificar y registrar explícitamente ambas versiones, sin presentarlas como iguales ni aceptar otra rueda."
- `clean_rebuild_evidence`: "METADATA/direct_url demostraron distribución 1.1.3 y SHA exacto; módulo 1.2.1 quedó fijado en lock y test. Frozen install, pip check, 16/16 y Audit PASS."
- `residual_risk`: "No afecta por sí sola el contrato HTTP probado, pero impide cerrar procedencia/versionado sin explicación verificable."

### FAIL-20260904-193 — patch de párrafos largos usó líneas parciales

- `status`: `REGRESSION_PROVEN`
- `severity`: `LOW`
- `fingerprint`: `startup-doc-long-line-partial-context`
- `classification`: `DOCUMENTATION_TOOLING`
- `violated_invariant`: "Un patch debe anclarse en líneas completas observadas o insertar en un límite estructural estable."
- `root_cause`: "Se intentó reemplazar sólo un fragmento de dos párrafos almacenados como líneas únicas muy extensas."
- `canonical_correction`: "Insertar una nota TikTok Lead como párrafo independiente junto al límite numerado, evitando reescribir líneas gigantes."
- `regression_test`: "Buscar el plan Lead en ambos protocolos y ejecutar el gate raíz."
- `clean_rebuild_evidence`: "Notas autónomas TikTok Lead insertadas en ambos protocolos y root PASS."
- `residual_risk`: "El patch fue atómico y no modificó ninguno de los dos archivos."

### FAIL-20260904-194 — Durable Task dejó un event loop sin cerrar en Audit

- `status`: `REGRESSION_PROVEN`
- `severity`: `MEDIUM`
- `fingerprint`: `durable-document-runtime-unclosed-proactor-event-loop`
- `classification`: `TEST_RUNTIME_RESOURCE_LEAK`
- `violated_invariant`: "Una suite aprobada no debe ocultar recursos asíncronos sin cerrar."
- `root_cause`: "Durante `test_pipeline.py`, CPython 3.12 emitió `ResourceWarning: unclosed event loop` entre instancias del worker Durable Task."
- `canonical_correction`: "La fuente exacta de `durabletask==1.9.0` confirmó que `TaskHubGrpcWorker.start` creaba un event loop sin cerrarlo. `ManagedTaskHubGrpcWorker` conserva la secuencia oficial y añade `shutdown_asyncgens`, `shutdown_default_executor`, `set_event_loop(None)` y `loop.close()`; el bloque se registra como `ADAPTED`."
- `regression_test`: "El verificador exige estáticamente el worker administrado y ejecuta contratos/pipeline con `-W error::ResourceWarning`."
- `clean_rebuild_evidence`: "Materialización canónica 13/13, suite 13/13, pipeline focal 8/8 y compileall PASS sin `ResourceWarning` ni sockets abiertos; Audit global posterior exit 0 con 153 packs, 121 fuentes y 16 adapters."
- `residual_risk`: "El SDK registra `StatusCode.CANCELLED` al detener el stream y luego `Worker shutdown completed`; no es una fuga. Backend productivo, identidad, carga y recovery del target siguen siendo gates separados."

### FAIL-20260904-195 — glob Markdown literal volvió a fallar en Windows

- `status`: `REGRESSION_PROVEN`
- `severity`: `LOW`
- `fingerprint`: `ripgrep-windows-markdown-glob-recurrence`
- `classification`: `DISCOVERY_TOOLING`
- `violated_invariant`: "En Windows, ripgrep recibe un directorio observado y filtros propios, no un wildcard literal de shell."
- `root_cause`: "La búsqueda del owner Durable Task pasó `implementation_packs/*.md` como path literal."
- `canonical_correction`: "Repetir contra `implementation_packs` y filtrar por contenido/extensión con ripgrep."
- `regression_test`: "La búsqueda devuelve el owner sin error y la recurrencia incrementa la lección previa."
- `clean_rebuild_evidence`: "Búsqueda contra directorio con `-g '*.md'` devolvió el owner Durable exacto sin error."
- `residual_risk`: "El segundo comando focal sí consultó un nombre supuesto; ninguna escritura ocurrió."

### FAIL-20260904-196 — segundo wildcard literal durante búsqueda de consumidores

- `status`: `REGRESSION_PROVEN`
- `severity`: `LOW`
- `fingerprint`: `ripgrep-windows-pack-plan-glob-second-recurrence`
- `classification`: `DISCOVERY_TOOLING`
- `violated_invariant`: "Tras registrar una recurrencia, la corrección se aplica a todas las búsquedas siguientes del mismo tramo."
- `root_cause`: "La búsqueda de consumers volvió a pasar `markdown_system/*PACK_PLAN.md` como path literal."
- `canonical_correction`: "Usar roots directorio sin wildcard y filtros `-g` explícitos."
- `regression_test`: "Inventario de consumers completo sin error 123."
- `clean_rebuild_evidence`: "Inventario con directorio + filtro devolvió exactamente cuatro planes consumidores y gobernó el bump."
- `residual_risk`: "Salida parcial descartada; ninguna escritura."

### FAIL-20260904-197 — versión Durable Task asumida en plan de franquicia

- `status`: `REGRESSION_PROVEN`
- `severity`: `LOW`
- `fingerprint`: `durable-pack-version-franchise-consumer-assumed`
- `classification`: `COMPOSITION_DISCOVERY`
- `violated_invariant`: "Actualizar sólo consumers devueltos por el inventario exacto, no inferir que un perfil amplio incluye el pack."
- `root_cause`: "El patch incluyó `FRANCHISE_COMPLETE_PACK_PLAN.md` aunque la búsqueda de consumers sólo había demostrado cuatro perfiles documentales."
- `canonical_correction`: "Quitar el owner supuesto y actualizar únicamente los cuatro planes observados."
- `regression_test`: "Búsqueda posterior muestra cero referencias a 0.1.1 para el pack y los cuatro perfiles componen."
- `clean_rebuild_evidence`: "Sólo los cuatro consumers observados cambiaron a 0.1.2; root compuso los cuatro perfiles."
- `residual_risk`: "El patch atómico fue rechazado; ningún archivo cambió."

### FAIL-20260904-198 — metadata provenance se aplicó al primer bloque coincidente

- `status`: `REGRESSION_PROVEN`
- `severity`: `HIGH`
- `fingerprint`: `durable-provenance-unanchored-block-mutation`
- `classification`: `PROVENANCE_INTEGRITY`
- `violated_invariant`: "Una reclasificación de procedencia debe anclarse al `block_id` exacto."
- `root_cause`: "El hunk genérico cambió `requirements-direct.in` a `ADAPTED` en vez de `documentflow/runtime.py`."
- `canonical_correction`: "Restaurar requirements a `AUTHORED` y cambiar únicamente el bloque runtime usando sus encabezados/block IDs como contexto."
- `regression_test`: "Inventario del pack debe contar exactamente un bloque runtime `ADAPTED`, doce categorías previas preservadas y hashes válidos."
- `clean_rebuild_evidence`: "Requirements volvió a AUTHORED, sólo runtime quedó ADAPTED; root aceptó 1123/128/105."
- `residual_risk`: "La metadata canónica está incorrecta hasta aplicar la corrección inmediata; contenido y hash no cambiaron."

### FAIL-20260904-199 — reemplazo parcial sobre línea extensa de notices

- `status`: `REGRESSION_PROVEN`
- `severity`: `LOW`
- `fingerprint`: `third-party-notices-long-line-partial-replacement`
- `classification`: `EDITING`
- `violated_invariant`: "Las correcciones documentales deben usar anclas completas observadas y conservar atomicidad."
- `root_cause`: "`apply_patch` recibió sólo un fragmento de una línea extensa y exigía que la línea completa coincidiera."
- `canonical_correction`: "Conservar el rechazo atómico y publicar la corrección cronológica como una oración autónoma junto al marcador de procedencia; editar la fila corta de reauditoría por separado."
- `regression_test`: "Buscar la corrección V246 y ejecutar el verificador raíz sobre notices y todos los hashes."
- `clean_rebuild_evidence`: "El intento no modificó archivos; la corrección autónoma y el root final gobiernan el cierre."
- `residual_risk`: "Ninguno después del read-back y root PASS; la frase original queda acompañada por su corrección explícita."

### FAIL-20260904-200 — preflight de franquicia conservaba inventario y brechas obsoletos

- `status`: `REGRESSION_PROVEN`
- `severity`: `MEDIUM`
- `fingerprint`: `franchise-preflight-stale-142-pack-gap-classification`
- `classification`: `AUTHORITY_FRESHNESS`
- `violated_invariant`: "El documento que decide qué falta debe reflejar el inventario ejecutable vigente y no volver a declarar ausente una capacidad incorporada."
- `root_cause`: "`FRANCHISE_PREFLIGHT_GAP.md` no se actualizó después de incorporar los runners/validators k6, ZAP/SCA, Debezium y los gates frontend posteriores."
- `canonical_correction`: "Reclasificar cada fila contra los packs observados: distinguir reusable presente, ejecución target pendiente y código realmente ausente; actualizar conteos sólo con el root vigente."
- `regression_test`: "Read-back sin 142 packs ni afirmaciones falsas de ausencia, seguido por `VERIFY_LIBRARY.ps1`."
- `clean_rebuild_evidence`: "Read-back confirma 153/1.356/667/45, k6 y ZAP/SCA presentes, Debezium presente y las brechas target separadas; `VERIFY_LIBRARY_PASS` posterior conserva 153 packs/1.356 archivos/667 Markdown/45 perfiles."
- `residual_risk`: "OpenAPI/codegen genérico, SAST/fuzz reusable y las ejecuciones target permanecen explícitos; no se presentan como resueltos."

### FAIL-20260905-201 — RESTler 9.3.1 no conserva global.json en la raíz

- `status`: `REGRESSION_PROVEN`
- `severity`: `LOW`
- `fingerprint`: `restler-9.3.1-root-global-json-assumed`
- `classification`: `UPSTREAM_LAYOUT`
- `violated_invariant`: "La auditoría debe observar el layout exacto de la revisión antes de inferir un archivo de toolchain."
- `root_cause`: "El probe reutilizó la expectativa histórica de `global.json` raíz, pero el archive oficial RESTler v9.3.1 no contiene ese path."
- `canonical_correction`: "Conservar hashes ya comprobados, enumerar los proyectos y archivos de configuración reales, derivar el SDK/TFM sólo desde ellos y repetir el probe."
- `regression_test`: "El expediente debe fijar paths existentes y el comando de build debe ejecutarse contra la revisión exacta sin depender de `global.json` supuesto."
- `clean_rebuild_evidence`: "La revisión exacta fue enumerada desde sus paths reales (`src/Restler.sln`, proyectos F# net8.0, Dockerfile .NET 8 y paquetes bloqueados); el Audit V247 pasó sin reutilizar el supuesto `global.json`. Esto cierra el fallo de inspección, no promueve RESTler ni atribuye un build aún no demostrado."
- `residual_risk`: "El toolchain y la reconstrucción RESTler siguen abiertos."

## 3. Resumen operativo

| ID | Huella | Severidad | Estado | Recurrencias | Capability | Invariante | Próximo gate |
|---|---|---|---|---:|---|---|---|
| FAIL-20260901-001 | preflight-playwright-offline-tarball-missing | MEDIUM | REGRESSION_PROVEN | 1 | TEST-PLATFORM | frozen offline inputs completos | repetir por host nuevo |
| FAIL-20260901-002 | v169-receipt-putaway-packaging-maintenance-chain | HIGH | REGRESSION_PROVEN | 0 | SUPPLY-WAREHOUSE-PACKAGING | conservación y gates completos | V170 pick/replenishment |
| FAIL-20260901-003 | readiness-review-query-shape-and-powershell-metadata-count | LOW | REGRESSION_PROVEN | 0 | DOCS-OPS | sólo mediciones completas | scoring por proyecto |
| FAIL-20260901-004 | v170-ripgrep-literal-markdown-glob | LOW | REGRESSION_PROVEN | 1 | DOCS-OPS | paths observados, no wildcard literal | target evidence por proyecto |
| FAIL-20260901-005 | v171-packaging-aware-pick-replenishment-maintenance-chain | HIGH | REGRESSION_PROVEN | 2 | SUPPLY-WAREHOUSE-PACKAGING | conservación física/composición y reconstrucción canónica | capability siguiente del proyecto |
| FAIL-20260901-006 | v172-connected-demand-discovery-truncation | LOW | REGRESSION_PROVEN | 0 | CONNECTED-DEMAND-WAREHOUSE | sólo inventarios completos por owner | implementar V172 |

## 4. Lecciones aplicadas al plan actual

| Lección previa | Decisión preventiva | Archivo/contrato/gate | Evidencia |
|---|---|---|---|
| LIB-FAIL-1119/1120 | Mantener instalación y ejecución Playwright aisladas del workspace padre y con lock exacto. | `VERIFY_EXECUTABLE_LIBRARY.ps1` | La instalación frozen/offline, cuatro tests browser y Audit pasaron sin ascender dependencias al workspace. |
| LIB-FAIL-1625..1656 | Invalidar cada intento parcial y repetir desde Markdown/PostgreSQL aislados hasta cerrar conservación, migración, composición y gates. | pack supply + evidence V169 + verificadores raíz/ejecutable | 75/75, backend 29/367, PostgreSQL 0001–0035, full Go/vet/build, root y Audit PASS. |
| LIB-FAIL-1657..1659 | Medir readiness sólo desde lecturas completas y comandos válidos. | `TOTAL_SYSTEM_CAPABILITY_CONTRACT.md` + metadata de packs | inventario completo de 48 superficies y clasificación exacta de 94 packs. |
| LIB-FAIL-1660..1667 | Enumerar paths observados, conservar cardinalidades y reconstruir readiness antes de Audit. | materializador + readiness tests + Audit | V170 cerró 62 pruebas, reconstrucción 8/8 y Audit. |
| LIB-FAIL-1668..1689 | Mantener arrays tipados en la misma sesión, fixtures dentro de contratos, columnas SQL calificadas, cleanup por dependencia, bases globales nuevas, evidencia histórica desacoplada de estado consumible y todos los consumidores de cardinalidad sincronizados. | supply 0.14.0 + backend + PostgreSQL + evidence V171 | 78/78, backend 29/370, 0001–0036/down/up, full Go/vet/build, root y Audit PASS. |
| LIB-FAIL-1690..1704 | Separar discovery, fijar identidades, usar paths/toolchains exactos, fixtures autoritativos, dependencias, párrafos literales completos y reconstrucción por etapas. | commerce + warehouse + PostgreSQL + BCApps tree lock | V172 cerró 83/83, backend 29/375, PostgreSQL/Go, root 94/1092/529 y Audit 94/121. |

## 5. Condiciones upstream seleccionadas

| Fuente/revisión | Fallo o condición retenida | Impacto en el proyecto | Aislamiento/gate | Estado |
|---|---|---|---|---|
| Microsoft Playwright 1.62.1 | El runtime exacto requiere artefactos npm/browser presentes o adquisición gobernada. | Reproducible en este host; otros hosts no heredan el PASS. | frozen lock + integridad + runtime browser + Audit | REGRESSION_PROVEN |

## 6. Cierre de ciclo

- [x] cada fallo observado tiene ID y huella;
- [x] ningún warning material, skip o test deshabilitado quedó sin clasificación;
- [x] el estado de readiness refleja el bloqueo ambiental real;
- [x] el artefacto exacto está presente y su instalación frozen/offline fue validada por el lock;
- [x] se reconstruyó desde destinos aislados y se repitió Audit completo;
- [x] V169 conserva un expediente separado y deja su siguiente brecha explícita.
- [x] los fallos de la revisión de readiness fueron registrados y su salida inválida no se usó para puntuar.
### FAIL-20260905-202 — nombre incorrecto al verificar el ZIP autocontenido de Kiota

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** admisión Microsoft Kiota 1.35.0 / receipt de artefacto
- **Fallo observado:** `Get-FileHash` buscó `kiota-win-x64.zip`, pero el artefacto descargado conserva el nombre versionado `kiota-win-x64-1.35.0.zip`; RESTler sí fue hasheado y el comando global continuó.
- **Causa:** el probe reutilizó un nombre abreviado distinto del nombre materializado.
- **Corrección inmediata:** resolver el path literal inventariado y verificar tamaño, digest, contenido y versión desde ese nombre exacto.
- **Prevención:** derivar futuras comprobaciones desde el receipt/asset oficial fijado, nunca desde un alias supuesto.
- **Evidencia de cierre requerida:** SHA-256 del ZIP exacto igual al digest oficial del release, contenido esperado y `kiota.exe --version` igual a la release/commit fijados.
### FAIL-20260905-203 — el entorno rechazó el probe combinado de Kiota por borrado recursivo

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** prueba de reproducibilidad Microsoft Kiota 1.35.0
- **Fallo observado:** el ejecutor rechazó antes de iniciar un comando que validaba y luego eliminaba recursivamente dos outputs temporales.
- **Causa:** el probe combinó una prueba no destructiva con limpieza recursiva, activando correctamente la política preventiva del entorno.
- **Corrección inmediata:** crear dos directorios temporales únicos y nunca limpiar ni sobrescribir durante la demostración.
- **Prevención:** los runners de generación deben fallar si el output existe; la limpieza debe ser una operación explícita, separada y recuperable.
- **Evidencia de cierre requerida:** dos outputs nuevos, generación exitosa y manifiestos/hash de árbol idénticos sin borrado previo.
### FAIL-20260905-204 — SHA de OSV transcripto desde una salida truncada

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** SCA del grafo Go generado por Kiota
- **Fallo observado:** el binario OSV tenía el tamaño esperado, pero el probe rechazó el SHA escrito manualmente porque no coincidía con el hash real descargado.
- **Causa:** se completó manualmente un digest que una salida previa había mostrado abreviado.
- **Corrección inmediata:** leer el digest completo desde `OFFICIAL_MAINTENANCE_TOOL_ARTIFACT_LOCK.md` y comparar el artefacto programáticamente; no ejecutar el scanner antes del PASS de identidad.
- **Prevención:** ningún hash se reconstruye desde UI truncada ni memoria; todos los runners consumen el lock canónico.
- **Evidencia de cierre requerida:** tamaño y SHA iguales al lock, `--version` exacto y scan focal exitoso.
### FAIL-20260905-205 — updater de packs no acepta bloques de esqueleto completamente vacíos

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** construcción del pack Microsoft Kiota 1.35.0
- **Fallo observado:** `update_pack_from_tree.ps1` actualizó el primer bloque y luego informó cero coincidencias para `acquire_kiota.ps1`; la materialización posterior no se ejecutó.
- **Causa:** los fences del esqueleto no contenían ninguna línea y el patrón canónico del updater requiere contenido capturable antes del fence de cierre.
- **Corrección inmediata:** usar un marcador temporal no vacío en cada bloque restante y volver a sincronizar todos los archivos desde el árbol de staging.
- **Prevención:** los esqueletos creados para `update_pack_from_tree.ps1` deben contener una línea placeholder explícita.
- **Evidencia de cierre requerida:** siete bloques actualizados, roundtrip materializado y hashes iguales al staging.
### FAIL-20260905-206 — invocación del materializador sin prefijo de path

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** roundtrip del pack Microsoft Kiota 1.35.0
- **Fallo observado:** los siete bloques se actualizaron, pero `materialize_markdown_pack.ps1` no fue resuelto como comando y el inventario posterior tampoco existió.
- **Causa:** PowerShell no busca scripts del directorio actual por nombre desnudo.
- **Corrección inmediata:** invocar `./materialize_markdown_pack.ps1` por path explícito y exigir exit/resultados antes del inventario.
- **Prevención:** todos los runners locales deben usar `Join-Path` o prefijo `./` para scripts del workspace.
- **Evidencia de cierre requerida:** materialización exitosa y siete archivos byte-idénticos al staging.
### FAIL-20260905-207 — recaída al comprobar LASTEXITCODE tras el materializador PowerShell

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** roundtrip del pack Microsoft Kiota 1.35.0
- **Fallo observado:** el materializador informó `Materialized 7 files`, pero el wrapper lanzó `materializer exit=` porque `$LASTEXITCODE` quedó nulo.
- **Causa:** se repitió el patrón ya prohibido de tratar un script PowerShell invocado en sesión como un ejecutable nativo con exit code persistente.
- **Corrección inmediata:** conservar el destino exitoso y validar contenido/hashes en una llamada separada; no repetir materialización.
- **Prevención:** para scripts PowerShell canónicos, usar excepciones/salida esperada o un proceso `pwsh` separado; nunca exigir `$LASTEXITCODE` no nulo después de `& script.ps1`.
- **Evidencia de cierre requerida:** siete archivos materializados byte-idénticos al staging y ejecución del verifier desde ese destino.
### FAIL-20260905-208 — hash de árbol Kiota no portable entre raíces temporales

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** verifier materializado Microsoft Kiota 1.35.0
- **Fallo observado:** dos generaciones dentro del verifier fueron idénticas (`adaefe6e…`), pero no coincidieron con el hash manual previo (`783003c6…`) y el gate falló antes de compilar.
- **Causa en investigación:** al menos un archivo generado parece incorporar contexto de la ubicación del contrato; comparar el árbol completo entre raíces distintas puede capturar metadata no portable.
- **Corrección inmediata:** comparar manifests por archivo contra la ejecución anterior, identificar exactamente el campo variable y definir un hash funcional que nunca excluya código.
- **Prevención:** una prueba de reproducibilidad debe ejecutarse tanto dentro de la misma raíz como entre raíces distintas antes de fijar un hash dorado.
- **Evidencia de cierre requerida:** archivo/campo variable demostrado, código generado byte-idéntico, criterio de normalización explícito y verifier completo PASS.
### FAIL-20260905-209 — patch Kiota mezcló contextos de archivos distintos

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** corrección del hash portable Kiota
- **Fallo observado:** `apply_patch` rechazó el patch completo al buscar líneas de `generate_go_client.ps1` dentro de `source-lock.json`; no aplicó cambios.
- **Causa:** faltaron encabezados `Update File` al pasar de un archivo de staging a otro.
- **Corrección inmediata:** dividir la modificación en patches atómicos para lock, generator, verifier y README.
- **Prevención:** un cambio multiarchivo sólo se envía después de comprobar cada transición de encabezado; preferir atomicidad por archivo en scripts de gates.
- **Evidencia de cierre requerida:** diff focal por archivo, resincronización del pack y verifier PASS.
### FAIL-20260905-210 — encabezado final supuesto en Official Upstream pack

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** actualización documental del lock Kiota 1.35.0
- **Fallo observado:** `apply_patch` rechazó todo porque `## 6. Verification evidence` no existe literalmente en el pack; no hubo cambios parciales.
- **Causa:** se asumió un encabezado inglés sin enumerar los headings reales del archivo.
- **Corrección inmediata:** leer metadata y encabezados exactos, luego aplicar patches focales.
- **Prevención:** en packs extensos, descubrir headings con `rg '^## '` antes de parchear secciones.
- **Evidencia de cierre requerida:** metadata 0.4.77/2026-09-05 y nota gobernante Kiota 1.35.0 visibles tras readback.
### FAIL-20260905-211 — staging Kiota detectado como directorio raíz no permitido

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** `VERIFY_LIBRARY.ps1` posterior al pack Kiota
- **Fallo observado:** la verificación raíz terminó 1 con `unknown top-level directory: .kiota_pack_stage`.
- **Causa:** el árbol de staging usado para sincronizar el pack permaneció en la raíz después del roundtrip exitoso.
- **Corrección inmediata:** eliminar sólo los siete archivos de staging ya sincronizados y su directorio vacío; conservar pack, evidencia y outputs temporales fuera de la biblioteca.
- **Prevención:** crear staging bajo `%TEMP%` o retirarlo inmediatamente después de `update_pack_from_tree` y roundtrip.
- **Evidencia de cierre requerida:** staging ausente y verificación raíz PASS.
### FAIL-20260905-212 — pack Kiota omitió las secciones contractuales 6–10

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** validación estructural global
- **Fallo observado:** `VERIFY_LIBRARY` exigió `## 6. Configuration surface` exactamente una vez fuera de fences y detuvo el pack Kiota.
- **Causa:** el pack nuevo condensó el cierre en `## 6. Verification record` y no siguió la plantilla completa 6–10.
- **Corrección inmediata:** reemplazar esa sección por Configuration surface, Dependency bill, Apply order, Verification y Reconstruction evidence.
- **Prevención:** todo pack nuevo se contrasta contra la lista de secciones obligatorias del verificador antes del primer root audit.
- **Evidencia de cierre requerida:** root verifier reconoce las diez secciones y materializa los siete archivos.
### FAIL-20260905-213 — THIRD_PARTY_NOTICES no incorporó los siete bloques Kiota

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** ledger global de procedencia
- **Fallo observado:** el root verifier calculó `AUTHORED=1130; ADAPTED=128; VERBATIM=105; TOTAL=1363` y rechazó el resumen anterior.
- **Causa:** el pack nuevo pasó estructura/materialización, pero el ledger de notices aún no había sido actualizado.
- **Corrección inmediata:** agregar la fuente Microsoft Kiota 1.35.0 con su licencia/claim y actualizar sólo `Current-Provenance-Counts` al conteo derivado.
- **Prevención:** después de agregar o cambiar bloques de un pack, ejecutar primero el conteo derivado y sincronizar notices antes del Audit global.
- **Evidencia de cierre requerida:** conteo exacto 1130/128/105/1363 y root verifier PASS.
### FAIL-20260905-214 — planes aún fijaban Official Upstream 0.4.76

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** composición/version closure posterior a Kiota 1.35.0
- **Fallo observado:** `VERIFY_LIBRARY` detectó `plan=0.4.76, pack=0.4.77` comenzando por `AMAZON_SPAPI_CATALOG_PACK_PLAN.md`.
- **Causa:** el pack fuente cambió de versión después de actualizar sus dos archivos materializados y los planes consumidores no fueron sincronizados aún.
- **Corrección inmediata:** enumerar cada plan que referencia `OFFICIAL-UPSTREAM-ACQUISITION-CORE` 0.4.76 y cambiar exclusivamente esa versión a 0.4.77.
- **Prevención:** toda actualización de pack dispara búsqueda de `packId` en todos los pack plans antes del root audit.
- **Evidencia de cierre requerida:** cero referencias consumidoras 0.4.76 y composición/version gate PASS.
### FAIL-20260905-215 — nuevo pack plan Kiota descubierto pero no compuesto por la suite

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** cobertura de pack plans
- **Fallo observado:** `VERIFY_LIBRARY` reportó `discovered=46 composed=45`.
- **Causa:** se creó `MICROSOFT_KIOTA_OPENAPI_CLIENT_PACK_PLAN.md` sin agregarlo a la lista explícita de composición del verificador.
- **Corrección inmediata:** incorporar el plan a la matriz de composición y verificar sus siete archivos/resultados.
- **Prevención:** todo nuevo `*_PACK_PLAN.md` debe agregarse en el mismo patch al listado gobernante del verificador.
- **Evidencia de cierre requerida:** `discovered=46 composed=46` y root PASS.
### FAIL-20260905-216 — consulta de rango asumió una sola coincidencia de plan

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** localización de matriz de pack plans
- **Fallo observado:** `Select-String` devolvió varios objetos y `$m.LineNumber-8` falló; la salida de `rg` sí identificó las líneas 350, 501 y 590.
- **Causa:** el probe no redujo explícitamente una colección a una coincidencia.
- **Corrección inmediata:** leer el rango fijo 470–605 indicado por la búsqueda y editar la matriz observada.
- **Prevención:** envolver búsquedas potencialmente múltiples en `Select-Object -First 1` antes de aritmética de líneas.
- **Evidencia de cierre requerida:** plan Kiota agregado y cobertura 46/46.
### FAIL-20260905-217 — Audit con Python por defecto incompatible con lock MarkItDown

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** auditoría ejecutable global con red
- **Fallo observado:** tras pasar Kiota, Playwright, Lighthouse y los gates intermedios, MarkItDown bloqueó porque el Python resuelto por defecto no era CPython 3.12 Windows x86-64.
- **Causa:** el Audit fue lanzado con Go explícito pero sin `-PythonExecutable`, permitiendo seleccionar el runtime global incompatible.
- **Corrección inmediata:** resolver el Python 3.12 configurado del workspace, verificar versión/plataforma y repetir Audit con ambos toolchains explícitos.
- **Prevención:** cualquier Audit reproducible debe fijar todos los toolchains requeridos, no sólo el añadido por la modificación actual.
- **Evidencia de cierre requerida:** MarkItDown lock/runtime PASS y auditoría global final PASS.

### FAIL-20260905-218 — patch de cierre V247 contenía un prefijo espurio

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** cierre de ledgers posterior al Audit V247
- **Fallo observado:** `apply_patch` rechazó atómicamente el diff al encontrar `n+` en vez de una línea añadida válida.
- **Causa:** un carácter de salto quedó serializado dentro de un patch manual extenso.
- **Corrección inmediata:** dividir el cierre entre el ledger de proyecto y el ledger de biblioteca, revisar cada diff y aplicarlos por separado.
- **Prevención:** evitar patches manuales demasiado extensos para tablas de muchas filas; validar que cada línea añadida empiece exactamente con `+`.
- **Evidencia de cierre requerida:** ambos ledgers sin estados abiertos, read-back correcto y verificación raíz PASS.

### Cierre común V247 de FAIL-20260905-202…218

Los fixes fueron reconstruidos desde Markdown y comprobados conjuntamente por `VERIFY_LIBRARY_PASS` (154 packs, 1.363 archivos, 670 Markdown y 46/46 perfiles) y por `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit` con CPython 3.12.14 x64 y Go 1.26.7 explícitos. El Audit adquirió Kiota 1.35.0 y OSV 2.5.1 por identidad exacta, generó dos árboles independientes, validó siete archivos y seis archivos Go byte-idénticos, compiló/tipó el cliente, obtuvo cero findings actuales en 14 paquetes y pasó el runtime MarkItDown bloqueado. Las condiciones productivas y el lane RESTler continúan separados; este cierre sólo prueba las regresiones enumeradas.

### FAIL-20260905-219 — RESTler no encuentra .NET en el host

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** reconstrucción Microsoft RESTler 9.3.1
- **Fallo observado:** `where.exe dotnet` no encontró binario y `dotnet --info` terminó 1.
- **Causa:** el host no tiene un SDK .NET expuesto en `PATH`; no es un fallo del source RESTler.
- **Corrección inmediata:** adquirir un SDK .NET 8 Windows x64 oficial desde metadata Microsoft, fijar tamaño/SHA y ejecutarlo por path temporal explícito.
- **Prevención:** el futuro pack debe adquirir/verificar su toolchain exacto o aceptar un path ya verificado; nunca depender del `PATH` global.
- **Evidencia de cierre requerida:** identidad del SDK, restore locked, build/tests y ejecución local RESTler o condición upstream exacta.

### FAIL-20260905-220 — URL de metadata .NET 8 incompleta

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** adquisición del toolchain RESTler
- **Fallo observado:** Azure Blob respondió `BlobNotFound` porque el path consultado contenía `release-metadata/8./`.
- **Causa:** se omitió el cero de `8.0` al escribir manualmente la URL.
- **Corrección inmediata:** usar el endpoint oficial literal `release-metadata/8.0/releases.json`, exigir campos no nulos y verificar el hash antes de extraer.
- **Prevención:** validar esquema y campos requeridos de toda metadata remota antes de consumirlos.
- **Evidencia de cierre requerida:** metadata 8.0 válida y artefacto cuyo SHA coincide con el campo oficial.

### FAIL-20260905-221 — ordenamiento de SDK .NET no acepta prereleases

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** selección del SDK oficial para RESTler
- **Fallo observado:** convertir versiones preview/RC a `System.Version` produjo errores y dejó la selección sin artefacto.
- **Causa:** se intentó inferir la versión más reciente ordenando todo el historial en vez de respetar `latest-sdk` de la metadata Microsoft.
- **Corrección inmediata:** localizar la entrada cuyo `sdk.version` coincide exactamente con `latest-sdk` y rechazar cardinalidad distinta de uno.
- **Prevención:** los selectors consumen el campo autoritativo explícito; no implementan ordenamiento propio cuando la autoridad ya publica la selección.
- **Evidencia de cierre requerida:** una sola release coincidente y ZIP Windows x64 con URL/hash oficiales.

### FAIL-20260905-222 — selector de nombre del ZIP .NET demasiado estricto

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** selección de artifact .NET 8
- **Fallo observado:** la release `latest-sdk` fue única, pero el filtro que exigía versión dentro del nombre encontró cero ZIP Windows x64.
- **Causa:** el patrón de nombre se escribió sin enumerar primero la forma publicada por la metadata.
- **Corrección inmediata:** listar los archivos `win-x64`, seleccionar por `rid` + extensión + tipo SDK observado y fijar la URL/hash exactos.
- **Prevención:** los filtros de artifacts se prueban contra el inventario publicado antes de exigir cardinalidad.
- **Evidencia de cierre requerida:** inventario observado y selección única verificable.

### FAIL-20260905-223 — RESTler v9.3.1 es tag, no GitHub Release

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** identidad oficial RESTler
- **Fallo observado:** el endpoint `releases/tags/v9.3.1` devolvió 404, aunque el ref del tag existe y resuelve a un commit verificado.
- **Causa:** se usó “release” informalmente sin comprobar el objeto GitHub exacto.
- **Corrección inmediata:** registrar `tag v9.3.1` y commit firmado `eaacd927…`; adquirir el archive del tag y no afirmar que existe un release asset.
- **Prevención:** verificar por separado ref/tag, release object, commit signature y artifact identity.
- **Evidencia de cierre requerida:** lock que diga tag, commit verificado y archive SHA-256 exacto.

### FAIL-20260905-224 — suite Python RESTler sin límite dejó de informar progreso

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** pruebas upstream RESTler 9.3.1
- **Fallo observado:** `pytest -q unit_tests` avanzó nueve casos y permaneció más de 90 segundos sin output; se interrumpió controladamente y no se contó como PASS.
- **Causa:** el scope global upstream no establece timeout por test y el modo quiet oculta el nodeid activo.
- **Corrección inmediata:** recolectar nodeids y repetir en verbose por archivo/caso con límites externos, conservando cualquier slow/fallo exacto.
- **Prevención:** toda suite fuzz/engine se ejecuta con visibilidad y límite por unidad desde el primer intento.
- **Evidencia de cierre requerida:** inventario completo y resultado por nodeid, con cualquier condición lenta explícitamente separada.

### FAIL-20260905-225 — path RESTler mistipeado durante inspección

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** diagnóstico del test RESTler lento
- **Fallo observado:** `rg` recibió `restler-fuzzer-9.c1` y no encontró el archivo; la suite activa no fue afectada.
- **Causa:** versión reescrita manualmente en vez de reutilizar el path `$src` ya resuelto.
- **Corrección inmediata:** derivar el path desde la raíz inventariada y repetir sólo la lectura.
- **Prevención:** conservar paths largos/versionados en variables observadas; no retranscribirlos.
- **Evidencia de cierre requerida:** fuente exacta del nodeid leída y resultado de la suite conservado.

### FAIL-20260905-226 — comentario y presupuesto de test_fuzz RESTler divergen

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** suite Python RESTler 9.3.1
- **Fallo observado:** el docstring afirma tres minutos/100 secuencias, pero el código fija `Fuzz_Time=0.1` horas (seis minutos) y `Num_Sequences=300`; el comentario inline también dice seis minutos.
- **Causa:** documentación local upstream desactualizada respecto de los valores ejecutados.
- **Corrección inmediata:** gobernar por el código exacto, no por el docstring, y registrar duración/300 secuencias del resultado.
- **Prevención:** el runner reusable expone explícitamente presupuesto y distingue suite source-exact de smoke acotado; nunca reduce una corrida y la llama upstream exacta.
- **Evidencia de cierre requerida:** `test_fuzz` exacto termina y el expediente conserva la divergencia sin parchear upstream.

### FAIL-20260905-227 — RESTler 9.3.1 tiene hallazgos en su grafo oficial

- **Fecha:** 2026-09-05
- **Estado:** `UPSTREAM_OPEN`
- **Etapa:** SCA de admisión Microsoft RESTler
- **Fallo observado:** OSV Scanner 2.5.1 recorrió 13 manifests/resultados, 558 ocurrencias de paquetes y devolvió exit 1 con 23 vulnerabilidades reportadas (19 IDs únicos) sobre la revisión exacta.
- **Causa:** los `packages.lock.json` y requirements del tag conservan versiones hoy afectadas; compilar y pasar tests no neutraliza advisories.
- **Corrección inmediata:** rechazar RESTler 9.3.1 para incorporación reusable; conservar evidencia, no actualizar sus dependencias localmente bajo falsa procedencia y evaluar otra fuente/revisión oficial.
- **Prevención:** todo motor de seguridad debe pasar SCA de su propio grafo antes de poder gobernar seguridad ajena.
- **Evidencia de cierre requerida:** una revisión oficial futura con locks limpios o un candidato oficial alternativo que pase identidad, licencia, build, tests, SCA y ejecución focal.

### FAIL-20260905-228 — escritura atómica del runner Go usó else inválido

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** construcción del gate de fuzz nativo Go
- **Fallo observado:** PowerShell rechazó `run_go_fuzz_gate.ps1` con “Try statement is missing its Catch or Finally block”.
- **Causa:** el bloque de escritura temporal terminó con `else` en vez de `catch`.
- **Corrección inmediata:** reemplazar por `catch`, retirar sólo el temporal de receipt y relanzar el verifier en raíz nueva.
- **Prevención:** parsear cada script PowerShell materializable antes de ejecutar su happy path.
- **Evidencia de cierre:** parser y round-trip 4/4; Audit V249 ejecutó baseline, semilla y fuzz real con `GO_NATIVE_FUZZ_GATE_PASS`, y el verifier focal cerró un positivo y dos negativos.

### FAIL-20260905-229 — validador confundió ./... con traversal

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** positivo del gate Go fuzz
- **Fallo observado:** el perfil seguro con package `./...` fue rechazado como unsafe.
- **Causa:** la regla prohibía cualquier `..`, incluyendo los tres puntos del wildcard oficial Go.
- **Corrección inmediata:** permitir sólo `.` o segmentos relativos alfanuméricos seguros y el sufijo literal `/...`; seguir rechazando `..` como segmento.
- **Prevención:** cada regla path-safety debe tener positivos canónicos además de negativos de traversal.
- **Evidencia de cierre:** Audit V249 aceptó `./...`, ejecutó el target y `verify_pack.ps1` rechazó `../x`; `GO_NATIVE_FUZZ_PACK_PASS positives=1 negatives=2`.

### FAIL-20260905-230 — helper update_pack no está en la raíz

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** sincronización del pack Go fuzz
- **Fallo observado:** `Get-Content update_pack_from_tree.ps1` no encontró el archivo.
- **Causa:** se asumió un path raíz no inventariado.
- **Corrección inmediata:** localizar por `rg --files` y usar el path literal observado.
- **Prevención:** resolver helpers por inventario antes de invocarlos.
- **Evidencia de cierre:** se resolvió `markdown_system/update_pack_from_tree.ps1`, se sincronizaron cuatro bloques y `VERIFY_LIBRARY_PASS` confirmó 155 packs/1.367 archivos.

### FAIL-20260905-231 — array -Files no cruzó correctamente pwsh -File

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** sincronización Go fuzz
- **Fallo observado:** el updater trató el segundo filename como argumento posicional.
- **Causa:** la sintaxis `@(...)` fue textual al cruzar otro proceso PowerShell.
- **Corrección inmediata:** invocar el script en la sesión actual con un array tipado real.
- **Prevención:** no serializar arrays PowerShell complejos a `pwsh -File`.
- **Evidencia de cierre:** el updater fue invocado en la sesión actual con array real, informó cuatro bloques y el materializador reprodujo los cuatro hashes.

### FAIL-20260905-232 — fences de contenido usaban tres backticks

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** materialización Go fuzz
- **Fallo observado:** el materializador informó `Missing content fence` en el primer archivo.
- **Causa:** el pack nuevo usó fences Markdown de tres backticks; el contrato materializable exige cuatro para distinguir contenido anidado.
- **Corrección inmediata:** cambiar exclusivamente los cuatro fences de contenido a cuatro backticks.
- **Prevención:** todo pack nuevo se inicia desde el patrón literal del PACK_CONTRACT.
- **Evidencia de cierre:** fences corregidos; root y Audit V249 materializaron 4/4 y el verifier focal pasó round-trip, positivo y negativos.

### FAIL-20260905-233 — consulta de estados abiertos usó regex mal escapada

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** cierre del gate Go fuzz V249
- **Fallo observado:** `rg` rechazó el patrón combinado con `unclosed character class`.
- **Causa:** backticks y clase negada atravesaron el quoting PowerShell/regex sin una representación segura.
- **Corrección inmediata:** usar búsquedas literales independientes para cada estado y cerrar esta lección sólo después del inventario limpio.
- **Prevención:** para tokens de estado conocidos, preferir `rg -F` sobre expresiones combinadas innecesarias.
- **Evidencia de cierre:** búsquedas `rg -F` independientes enumeraron únicamente 228–233/1962–1967 antes del cierre; la repetición posterior no devuelve estados locales abiertos.

### FAIL-20260905-234 — primer dotnet.exe inventariado no contenía SDK

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** reauditoría Microsoft DevSkim posterior a V249
- **Fallo observado:** `elite-dotnet-v125/sdk/dotnet.exe --version` informó que no encontraba SDKs; el segundo runtime inventariado devolvió 10.0.400.
- **Causa:** se trató la presencia del host ejecutable como prueba de que su carpeta conservaba un SDK completo.
- **Corrección inmediata:** seleccionar exclusivamente `elite-dotnet-v127/dotnet.exe` cuya versión 10.0.400 fue ejecutada, y no contar el primer path como toolchain disponible.
- **Prevención:** toda selección de runtime exige ejecutar su comando de identidad antes de restore/build.
- **Evidencia de cierre:** SDK 10.0.400 identificado y usado para restore/build/tests; V250 rechaza el candidato por su grafo, no por el primer path inválido.

### FAIL-20260905-235 — DevSkim main requiere feed Microsoft no autenticado

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** reconstrucción SAST Microsoft DevSkim commit `a452aa5…`
- **Fallo observado:** `dotnet test` no pudo restaurar porque `nuget.config` dirige a `PublicRegistriesFeed` de Azure DevOps y el endpoint respondió 401.
- **Causa:** el source oficial fijado depende de una configuración de feed que no es públicamente anónima desde este entorno; no es un fallo de compilación demostrado.
- **Corrección inmediata:** repetir sin editar upstream, pasando explícitamente el source público NuGet.org; conservar la divergencia respecto del feed oficial.
- **Prevención:** inspeccionar `nuget.config` y autenticar/acordar sources antes del restore; nunca interpretar 401 como vulnerabilidad ni PASS.
- **Evidencia de cierre:** NuGet.org explícito restauró el source sin editarlo; la divergencia del feed queda en V250 y el grafo produjo NU1902.

### FAIL-20260905-236 — dotnet test no acepta --source en esta invocación

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** restore público DevSkim main
- **Fallo observado:** MSBuild recibió `--source` como switch desconocido al usarlo directamente con `dotnet test`.
- **Causa:** se aplicó una opción del comando `dotnet restore` al wrapper `dotnet test` en vez de separar restore y test.
- **Corrección inmediata:** ejecutar `dotnet restore --source ...` y luego `dotnet test --no-restore` sobre el mismo target/framework.
- **Prevención:** separar opciones de adquisición NuGet de las opciones de VSTest.
- **Evidencia de cierre:** restore separado con NuGet.org y test `--no-restore` ejecutados; V250 conserva commands y resultado.

### FAIL-20260905-237 — apphost de tests DevSkim no resolvió runtime portable

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** pruebas DevSkim main con SDK 10.0.400 portable
- **Fallo observado:** restore y build net10.0 terminaron, pero el apphost de Microsoft Testing Platform no encontró el runtime porque `DOTNET_ROOT_X64`/`DOTNET_ROOT` no estaban definidos.
- **Causa:** el SDK está extraído en un path temporal no registrado en Windows y el testhost se lanzó como apphost independiente.
- **Corrección inmediata:** repetir `--no-build --no-restore` con `DOTNET_ROOT` y `DOTNET_ROOT_X64` apuntando al SDK observado.
- **Prevención:** runtimes portables .NET deben exponer explícitamente su raíz a procesos hijos.
- **Evidencia de cierre:** repetición con `DOTNET_ROOT` y `DOTNET_ROOT_X64` explícitos terminó exit 0; la vulnerabilidad se registra por separado en FAIL-238/UP-FAIL-204.

### FAIL-20260905-238 — DevSkim actual sigue incorporando SharpCompress 0.40.0 afectado

- **Fecha:** 2026-09-05
- **Estado:** `UPSTREAM_OPEN`
- **Etapa:** SCA de Microsoft DevSkim current main
- **Fallo observado:** restore del commit firmado `a452aa5…` emitió NU1902 por `SharpCompress 0.40.0`/`GHSA-6c8g-7p36-r338`; `.deps.json` confirma la versión en library/CLI/tests.
- **Causa:** las dependencias transitivas de Application Inspector 1.9.50 del source actual todavía resuelven la versión afectada.
- **Corrección inmediata:** mantener DevSkim `REJECTED_COMPONENT`; no hacer override/backport local bajo procedencia Microsoft.
- **Prevención:** todo analizador de seguridad pasa SCA de su propio runtime y archive-handling antes de gobernar código ajeno.
- **Evidencia de cierre requerida:** release oficial futura con grafo corregido, tests y regresiones de archives hostiles; expediente V250.

### FAIL-20260905-239 — Ruta asumida para PACK_CONTRACT

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** auditoría del flujo de resolución de capabilities ausentes
- **Fallo observado:** la lectura inicial buscó `implementation_packs/PACK_CONTRACT.md`, path inexistente.
- **Causa:** se infirió la ubicación por cercanía semántica en vez de inventariar el repositorio.
- **Corrección inmediata:** localizar el archivo con `rg --files` y leer `markdown_system/PACK_CONTRACT.md`.
- **Prevención:** resolver siempre los documentos canónicos por inventario antes de abrirlos; no derivar paths por memoria.
- **Evidencia de cierre:** `rg --files` devolvió una única ubicación canónica y la investigación continuó sobre ese archivo.

### FAIL-20260905-240 — Sincronización y cleanup combinados rechazados

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Recurrencia:** 3
- **Etapa:** construcción del capability-gap pack
- **Fallo observado:** el host rechazó antes de ejecutar un comando que combinaba updater, materialización, tests y borrado de un temporal calculado.
- **Causa:** una operación destructiva sobre variable computada quedó encadenada a otras shells/acciones y no cumplió el control de paths de Windows.
- **Corrección inmediata:** separar sincronización, materialización/test y cleanup; usar un destino absoluto explícito y verificar que permanezca bajo `%TEMP%` antes de borrarlo.
- **Prevención:** no combinar creación/prueba con borrado recursivo en una misma invocación; validar el target final en la misma shell nativa.
- **Evidencia de cierre:** updater ejecutado solo; tests fuente PASS; el target fue inventariado por path absoluto y los archivos/directorios propios se retiraron con APIs .NET nativas en PowerShell, sin cross-shell ni path calculado.

### FAIL-20260905-241 — apply_patch no elimina bytecode binario

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** cleanup del staging capability-gap
- **Fallo observado:** `apply_patch` rechazó dos `.pyc` porque no son UTF-8.
- **Causa:** el editor de parches es textual y el bytecode compilado es binario.
- **Corrección inmediata:** después de verificar paths absolutos dentro del staging, eliminar los dos archivos exactos mediante `[IO.File]::Delete` y retirar sólo sus directorios vacíos mediante `[IO.Directory]::Delete` en la misma shell PowerShell.
- **Prevención:** compilar materializaciones en `%TEMP%`, no dentro de un staging textual del workspace; reservar `apply_patch` para source UTF-8.
- **Evidencia de cierre:** `.pack_staging` ya no existe; no se usó shell cruzada ni borrado amplio.

### FAIL-20260905-242 — Loader legacy de dependencias no disponible

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** resolución de toolchains para Audit V251
- **Fallo observado:** `codex_app__load_workspace_dependencies` respondió que la vía dinámica dejó de estar disponible y que debe usarse el servidor MCP de Codex.
- **Causa:** se invocó el alias legacy aunque el entorno expone la variante MCP vigente.
- **Corrección inmediata:** llamar `mcp__codex_app__load_workspace_dependencies` y usar únicamente paths observados.
- **Prevención:** preferir la herramienta MCP explícitamente publicada por el contexto actual.
- **Evidencia de cierre:** el servidor MCP devolvió CPython 3.12.14 y demás paths del bundle; Go 1.26.7 fue observado en el cache temporal existente y el Audit V251 terminó exit 0 con ambos paths exactos.

### FAIL-20260905-243 — path del workspace retranscrito incorrectamente

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Recurrencia:** 3
- **Etapa:** inspección de fixtures y hashes DevSkim V252
- **Fallo observado:** dos comandos read-only y un patch no iniciaron porque la raíz fue retranscrita como `Publicaneers`, `Public EliteBEL` y `Public EliteXana`.
- **Causa:** se escribió manualmente un path ya conocido en vez de reutilizar la raíz exacta del workspace.
- **Corrección inmediata:** repetir ambos comandos con la raíz exacta observada `Public Elite Codes`.
- **Prevención:** conservar y reutilizar el workdir observado; no abreviar ni completar nombres de carpetas por memoria.
- **Evidencia de cierre:** los probes y el patch se repitieron desde la raíz exacta y terminaron correctamente.

### FAIL-20260905-244 — auditoría NuGet heredó feed no autenticado

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** SCA de la adaptación DevSkim
- **Fallo observado:** `dotnet list package --vulnerable --source NuGet.org` informó cero findings pero emitió `NU1900` al consultar también el feed Azure definido por upstream.
- **Causa:** `--source` limita la búsqueda del comando pero el restore implícito todavía consumió `nuget.config` del repositorio.
- **Corrección inmediata:** crear una configuración con `<clear/>` y NuGet.org único, ejecutar restore explícito y después auditoría `--no-restore --configfile`.
- **Prevención:** ningún resultado SCA con warning de adquisición cuenta como limpio; sources y restore deben ser cerrados y observables.
- **Evidencia de cierre:** tres proyectos restaurados sin warnings y tres reportes JSON de vulnerabilidades con cero findings.

### FAIL-20260905-245 — interpolación PowerShell inválida antes de dos puntos

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** parser del runner DevSkim V252
- **Fallo observado:** `$name:` fue interpretado como referencia de variable inválida.
- **Causa:** faltaron llaves al interpolar una variable inmediatamente antes de `:`.
- **Corrección inmediata:** reemplazar por `${name}:` y repetir parser/negativos.
- **Prevención:** delimitar con `${}` toda variable seguida por un carácter válido en nombres scoped.
- **Evidencia de cierre:** sintaxis de ambos scripts y dos negativos fail-closed PASS.

### FAIL-20260905-246 — bytes del archive GitHub cambiaron sin cambiar el código

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** adquisición exacta DevSkim V252
- **Fallo observado:** el mismo URL por commit produjo ZIPs de 935.159 y 949.833 bytes con SHA-256 distintos.
- **Causa:** GitHub documenta que puede regenerar source archives con distinto layout de compresión aunque los archivos del commit permanezcan iguales.
- **Corrección inmediata:** fijar commit/tree/firma y hash canónico ordenado de path+tamaño+SHA-256 por archivo; conservar sólo como observación el hash de transporte.
- **Prevención:** reservar digest de contenedor como identidad para release assets inmutables; para source archives generados verificar contenido canónico y extracción segura.
- **Evidencia de cierre:** ambos ZIPs produjeron 268 archivos y el mismo árbol canónico `4914ee46…`; el runner final consulta commit firmado y rechaza drift de contenido.

### FAIL-20260905-247 — nombre raíz variable del archive GitHub

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** inventario ZIP DevSkim V252
- **Fallo observado:** el archive regenerado usó `DevSkim-<commit>/` en lugar de `microsoft-DevSkim-a452aa5/`.
- **Causa:** se fijó un detalle de empaquetado que no forma parte del contenido del commit.
- **Corrección inmediata:** exigir exactamente una raíz segura, derivarla del archive y eliminarla antes de calcular el árbol canónico.
- **Prevención:** validar cardinalidad y seguridad de raíz, no su spelling, cuando la autoridad no lo garantiza.
- **Evidencia de cierre:** ambos root layouts convergen al mismo árbol; runtime V252 PASS.

### FAIL-20260905-248 — hash adaptado dependía del escritor exacto

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** patch reproducible DevSkim V252
- **Fallo observado:** los hashes de csproj adaptados con `apply_patch` no coincidieron con los producidos por `ReadAllText/Replace/WriteAllText` del runner.
- **Causa:** el primer experimento y el runner no compartían el mismo tratamiento byte-exacto de terminadores de línea.
- **Corrección inmediata:** derivar y fijar los hashes desde el algoritmo exacto que ejecuta el runner sobre los bytes oficiales.
- **Prevención:** la identidad de un patch reproducible se calcula con su productor final, no con un editor alternativo.
- **Evidencia de cierre:** ambos csproj finales coinciden con hashes fijados y el build/test/SCA completo pasa.

### FAIL-20260905-249 — JSON sin `frameworks` bajo StrictMode

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** parsing del reporte NuGet sin findings
- **Fallo observado:** el reporte JSON de cero vulnerabilidades omitió `frameworks` y StrictMode rechazó el acceso directo.
- **Causa:** el parser trató un campo opcional como obligatorio precisamente en el caso vacío.
- **Corrección inmediata:** consultar `PSObject.Properties`, aceptar ausencia como cero y hacer lo mismo para listas top-level/transitivas.
- **Prevención:** todo parser de output externo prueba explícitamente esquemas vacíos, positivos y campos opcionales.
- **Evidencia de cierre:** tres reportes vacíos se procesan como cero y la suite runtime completa pasa.

### FAIL-20260905-250 — `LiteralPath` no expande wildcard

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** publicación atómica del CLI adaptado
- **Fallo observado:** `Copy-Item -LiteralPath publish\*` buscó literalmente un archivo llamado `*`.
- **Causa:** se mezcló intención de glob con una API deliberadamente literal.
- **Corrección inmediata:** enumerar hijos exactos de `publish` y copiar cada `FullName` con `-LiteralPath`.
- **Prevención:** no usar wildcard con parámetros literales; inventariar antes de copiar.
- **Evidencia de cierre:** runtime, licencias, lock, provenance y manifest se publican; verifier completo PASS.

### FAIL-20260905-251 — patch multiarchivo perdió el contexto de destino

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** endurecimiento de identidad canónica DevSkim
- **Fallo observado:** `apply_patch` buscó líneas del runner dentro de `source-lock.json` y rechazó el patch completo.
- **Causa:** faltó separar explícitamente el segundo `Update File` al combinar cambios de JSON y PowerShell.
- **Corrección inmediata:** dividir lock/verifier y runner en patches focales con contexto observado.
- **Prevención:** patches multiarchivo largos se separan por formato o responsabilidad.
- **Evidencia de cierre:** cambios aplicados por archivo, read-back implícito en parser y runtime PASS.

### FAIL-20260905-252 — array no cruza `pwsh -File` como lista tipada

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Recurrencia:** 2
- **Etapa:** sincronización del pack DevSkim
- **Fallo observado:** el updater recibió argumentos posteriores como posicionales y luego una cadena única con comas.
- **Causa:** se volvió a intentar transportar `string[]` mediante parsing CLI de `pwsh -File`.
- **Corrección inmediata:** invocar el script en la sesión PowerShell actual con `&` y un array real `@(...)`.
- **Prevención:** reutilizar la lección LIB-FAIL-1965 para todos los parámetros array de scripts locales.
- **Evidencia de cierre:** cuatro bloques sincronizados, materializados y verificados byte-idénticos.

### FAIL-20260905-253 — pack V252 omitió secciones normativas 6–10

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** verificación raíz del pack DevSkim
- **Fallo observado:** `VERIFY_LIBRARY` rechazó la ausencia de `## 6. Configuration surface`.
- **Causa:** el primer documento condensó verificación y límites en secciones 6/7 no compatibles con el contrato común de packs.
- **Corrección inmediata:** añadir exactamente Configuration surface, Dependency bill, Apply order, Verification y Reconstruction evidence en orden.
- **Prevención:** todo pack nuevo se inicia con las diez secciones literales de `PACK_CONTRACT` antes de sincronizar materialization blocks.
- **Evidencia de cierre:** estructura 1–10 completa y verificación raíz posterior PASS.

### FAIL-20260905-254 — se repitió el loader legacy de dependencias

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Recurrencia:** 2
- **Etapa:** resolución de Python para Audit V252
- **Fallo observado:** se invocó otra vez el alias dinámico retirado y devolvió la misma indicación de usar MCP.
- **Causa:** la elección de herramienta no reutilizó la lección FAIL-242 ya registrada.
- **Corrección inmediata:** ejecutar directamente `mcp__codex_app__load_workspace_dependencies` y consumir su Python exacto.
- **Prevención:** cuando un alias declara retiro permanente, eliminarlo de la ruta operativa y usar sólo el proveedor MCP vigente.
- **Evidencia de cierre:** MCP devolvió el bundle 26.904.11930 y Python ejecutable; Audit V252 usa ese path.

### FAIL-20260905-255 — Docker no está instalado en el host

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** discovery del lane SAST gratuito oficial
- **Fallo observado:** `docker version` no pudo resolver el ejecutable; WSL2 sí está configurado con Ubuntu.
- **Causa:** se probó disponibilidad antes de seleccionar la variante de ejecución y este host no tiene Docker CLI/runtime.
- **Corrección inmediata:** descartar cualquier afirmación de prueba del contenedor GitLab/ASH en este host y priorizar el binario Windows firmado de Opengrep más orquestación ASH.
- **Prevención:** todo pack con variante container debe probar toolchain explícito y mantenerla `BLOCKED` cuando no exista, sin bloquear un lane local independiente.
- **Evidencia de cierre:** la ausencia quedó clasificada como condición del host; no se descargó ni ejecutó ninguna imagen.

### FAIL-20260905-256 — ASH 3.7.1 contiene un advisory en el lock de desarrollo

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** SCA del candidato oficial AWS ASH
- **Fallo observado:** Google OSV-Scanner 2.5.1 encontró `GHSA-xvg9-69gf-fjrf` en `mkdocs-material==9.7.6` dentro del `uv.lock` oficial.
- **Causa:** la release AWS fija una versión de documentación afectada; el grafo runtime separado no contiene el paquete.
- **Corrección inmediata:** conservar el source exacto sin promoverlo y separar explícitamente grafo runtime de grafo dev; una adaptación futura debe actualizar el lock y repetir toda la suite/SCA.
- **Prevención:** escanear siempre locks completos y runtime exportado por separado; cero hallazgos runtime no borra una vulnerabilidad del toolchain.
- **Evidencia de cierre:** lock completo 169 paquetes/1 advisory; export runtime 83 paquetes/cero hallazgos, ambos mediante OSV 2.5.1 exacto.

### FAIL-20260905-257 — el comando de desarrollo ASH no instala el extra CDK requerido por pruebas

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** suite source-exact AWS ASH 3.7.1
- **Fallo observado:** `uv sync --locked --all-groups` seguido de la suite predeterminada produjo 26 fallos: seis requieren privilegio symlink en Windows y veinte pertenecen al wrapper CDK sin el extra opcional instalado.
- **Causa:** el procedimiento documentado `uv sync` no instala `[project.optional-dependencies].cdk`, aunque esas pruebas se recolectan por defecto; Windows sin Developer Mode tampoco puede crear los symlinks de seguridad.
- **Corrección inmediata:** instalar `--all-extras` para evaluar CDK y clasificar los seis tests symlink como evidencia POSIX/privilegio aún no demostrada, no como PASS.
- **Prevención:** verificar extras efectivamente instalados y prerequisitos OS antes de interpretar una suite agregada.
- **Evidencia de cierre:** extra CDK exacto instalado desde el mismo `uv.lock`; los focales avanzaron hasta revelar la interferencia paralela separada en FAIL-258.

### FAIL-20260905-258 — tests CDK de ASH interfieren bajo xdist

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** repetición focal AWS ASH 3.7.1 con extras
- **Fallo observado:** 13/53 tests del wrapper CDK fallaron con `-n auto`, pero exactamente los mismos 53/53 pasaron con `-n 0 --no-cov`.
- **Causa:** los doubles focales mutan estado/entorno global sensible a ejecución paralela; el runner predeterminado usa 24 workers.
- **Corrección inmediata:** no presentar la suite predeterminada como verde; conservar el PASS secuencial únicamente como diagnóstico y mantener el componente condicionado.
- **Prevención:** ejecutar tanto el comando upstream exacto como una repetición secuencial; una repetición alternativa explica el fallo pero no reemplaza el gate oficial.
- **Evidencia de cierre:** default agregado 4.556 PASS/156 skip/26 fail; focal CDK secuencial 53/53 PASS.

### FAIL-20260905-259 — pipeline vacío en inventario de reglas GitLab

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** inventario del paquete oficial GitLab SAST rules 2.9.3
- **Fallo observado:** PowerShell rechazó el comando antes de ejecutarlo por un pipe colocado inmediatamente después del cierre del `foreach`.
- **Causa:** se intentó canalizar directamente una sentencia `foreach` sin capturar primero su salida.
- **Corrección inmediata:** asignar los objetos producidos a una colección y formatearla en una sentencia separada.
- **Prevención:** en probes PowerShell con bucles, separar producción y presentación de resultados.
- **Evidencia de cierre:** inventario repetido con colección explícita y hashes/rule counts obtenidos.

### FAIL-20260905-260 — validación serial de reglas excedió la ventana de ejecución

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** compatibilidad OpenGrep 1.29.0 / reglas GitLab abiertas 2.9.3
- **Fallo observado:** ejecutar cada configuración en un proceso independiente consumió la ventana de 30 segundos sin devolver un resultado concluyente.
- **Causa:** el probe pagó repetidamente el arranque y parsing del motor para cada YAML.
- **Corrección inmediata:** validar las configuraciones en lotes pequeños o en una única invocación con múltiples `--config`, manteniendo conteo de errores.
- **Prevención:** no usar un proceso SAST por archivo cuando el CLI admite composición de reglas.
- **Evidencia de cierre:** lote compuesto posterior terminó dentro del límite con errores explícitamente contados.

### FAIL-20260905-261 — wrapper interpretó mal el exit code de ASH

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** integración AWS ASH 3.7.1 + OpenGrep 1.29.0
- **Fallo observado:** ASH detectó dos hallazgos y terminó correctamente con código `2`, pero el wrapper experimental sólo admitía `0` o `1` y lanzó una excepción falsa.
- **Causa:** se aplicó el convenio del CLI OpenGrep al orquestador ASH; ASH reserva `2` para findings accionables cuando `fail_on_findings: true`.
- **Corrección inmediata:** modelar `0=CLEAN`, `2=FINDINGS`; cualquier otro código es error, y validar además el SARIF/estado del scanner.
- **Prevención:** fijar y probar por herramienta la semántica exacta de salida, sin generalizar entre CLIs.
- **Evidencia de cierre:** caso hostil esperado acepta `2` con dos findings; corpus limpio separado exige `0`.

### FAIL-20260905-262 — nombre inferido del pack DevSkim

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** reutilización de patrón de pack SAST
- **Fallo observado:** la lectura apuntó a `MICROSOFT_DEVSKIM_ADAPTED_LINT_GATE.md`, inexistente; el inventario de la misma llamada mostraba el nombre real `MICROSOFT_DEVSKIM_ADAPTED_SAST_GATE.md`.
- **Causa:** se retranscribió un nombre semánticamente parecido en vez de consumir la ruta enumerada.
- **Corrección inmediata:** usar la ruta exacta entregada por `rg --files`.
- **Prevención:** toda lectura posterior a discovery toma el path observado sin reconstruirlo manualmente.
- **Evidencia de cierre:** pack SAST exacto leído y usado como patrón de composición.

### FAIL-20260905-263 — se repitió el pipe inválido después de `foreach`

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Recurrencia:** 2
- **Etapa:** hashes de licencias del lane SAST
- **Fallo observado:** un segundo probe volvió a canalizar directamente el cierre de `foreach` y PowerShell lo rechazó antes de ejecutar.
- **Causa:** la corrección de FAIL-259 no se propagó al bloque final del comando compuesto.
- **Corrección inmediata:** capturar también las filas de licencias en una variable antes de serializar.
- **Prevención:** no permitir la forma textual `} |` después de un `foreach`; usar siempre `$rows = foreach (...)`.
- **Evidencia de cierre:** identidades y licencias repetidas con colecciones explícitas.

### FAIL-20260905-264 — tercera recurrencia del pipe post-`foreach`

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Recurrencia:** 3
- **Etapa:** distribución de severidades GitLab
- **Fallo observado:** el probe volvió a usar `} | Format-Table` tras un `foreach` y fue rechazado antes de ejecutar.
- **Causa:** se siguió componiendo el loop y su presentación en una sola línea pese a la prevención ya registrada.
- **Corrección inmediata:** imponer una variable `$rows` como estructura del comando, no sólo como recomendación.
- **Prevención:** los probes restantes de este tramo usarán exclusivamente `$rows = foreach (...) { ... }; $rows | ...`.
- **Evidencia de cierre:** distribución de severidades obtenida con esa estructura literal.

### FAIL-20260905-265 — caché confundido con verificación totalmente offline

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** contrato de instalación GitLab/OpenGrep SAST
- **Fallo observado:** el README inicial afirmaba que un caché completo permitía reconstruir sin red, aunque el procedimiento oficial `cosign verify-blob` con certificado/firma puede consultar material de confianza/transparencia Sigstore.
- **Causa:** se equiparó evitar descargas de artefactos con ausencia total de egress durante verificación criptográfica.
- **Corrección inmediata:** exigir `AllowNetwork` para toda instalación; conservar el caché sólo como optimización de bytes/latencia y ejecutar APIs de autoridad siempre.
- **Prevención:** una variante se denomina offline únicamente después de probar todos los resolvers de confianza con red bloqueada y material de confianza fijado.
- **Evidencia de cierre:** contrato, runner y README ya no prometen reconstrucción offline; receipt exige `authority_apis_checked: true`.

### FAIL-20260905-266 — static y runtime SAST combinados ocultaron el cierre

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** repetición posterior al endurecimiento de red
- **Fallo observado:** una sola invocación ejecutó verifier estático y runtime y llegó al límite de 30 segundos después del PASS limpio, sin devolver la línea final del gate.
- **Causa:** se sumaron dos ejecuciones completas en una ventana pensada para una sola.
- **Corrección inmediata:** mantener el PASS estático ya observado y repetir únicamente el runtime en una llamada independiente.
- **Prevención:** gates de aproximadamente 25–30 segundos se ejecutan solos; no se anteponen probes ya demostrados.
- **Evidencia de cierre:** runtime aislado posterior devuelve la línea final `GITLAB_OPENGREP_SAST_PACK PASS`.

### FAIL-20260905-267 — ruta inferida del updater pese al inventario

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Recurrencia:** 2 del patrón de path inferido
- **Etapa:** sincronización del pack GitLab/OpenGrep
- **Fallo observado:** `rg --files` devolvió `markdown_system/update_pack_from_tree.ps1`, pero la lectura de la misma invocación apuntó a `implementation_packs/update_pack_from_tree.ps1`.
- **Causa:** el comando se escribió con una ubicación supuesta antes de reutilizar el path observado.
- **Corrección inmediata:** abrir y ejecutar literalmente `markdown_system/update_pack_from_tree.ps1`.
- **Prevención:** separar discovery y consumo cuando el path aún no está fijado; no incluir una segunda ruta manual en la llamada de discovery.
- **Evidencia de cierre:** helper exacto leído y usado para sincronizar cinco bloques.

### FAIL-20260905-268 — SourceRoot no conservó el prefijo del manifest

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** sincronización del pack GitLab/OpenGrep
- **Fallo observado:** el updater encontró cero bloques para `source-lock.json` porque el pack declara `gitlab_opengrep_signed_sast/source-lock.json`.
- **Causa:** el staging se pasó desde dentro del directorio destino y perdió el prefijo relativo que forma parte del manifest exacto.
- **Corrección inmediata:** crear un source tree padre que contenga `gitlab_opengrep_signed_sast/` y sincronizar las cinco rutas completas.
- **Prevención:** el `SourceRoot` del updater debe reproducir literalmente las rutas `FILE` del pack.
- **Evidencia de cierre:** cinco bloques sincronizados y round-trip byte-idéntico.

### FAIL-20260905-269 — `$LASTEXITCODE` heredado tras updater PowerShell

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Recurrencia:** patrón ya documentado
- **Etapa:** sincronización del pack GitLab/OpenGrep
- **Fallo observado:** el updater imprimió `Updated 5 blocks`, pero el wrapper lanzó `updater failed` al evaluar `$LASTEXITCODE` heredado de un ejecutable nativo previo.
- **Causa:** un script PowerShell invocado en la sesión actual no define contractualmente `$LASTEXITCODE`.
- **Corrección inmediata:** aceptar la ausencia de excepción como éxito y verificar el resultado mediante materialización/hashes en una llamada separada.
- **Prevención:** no inspeccionar `$LASTEXITCODE` tras `& <script.ps1>` en la misma sesión; úselo sólo para procesos nativos o `pwsh` hijo.
- **Evidencia de cierre:** cinco bloques presentes y materialización byte-idéntica posterior.

### FAIL-20260905-270 — tercera recurrencia del loader dinámico retirado

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Recurrencia:** 3
- **Etapa:** resolución de toolchains para Audit V253
- **Fallo observado:** se invocó otra vez el alias dinámico `codex_app__load_workspace_dependencies`, que respondió que debe usarse el servidor MCP.
- **Causa:** no se aplicó la exclusión operativa ya fijada por FAIL-242/254.
- **Corrección inmediata:** usar directamente `mcp__codex_app__load_workspace_dependencies`.
- **Prevención:** considerar el alias dinámico eliminado de manera permanente en esta sesión y no volver a seleccionarlo.
- **Evidencia de cierre:** MCP vigente devuelve CPython/bundle y el Audit V253 usa ese path.

### FAIL-20260905-271 — limpieza temporal rechazada por la política del ejecutor

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** higiene posterior al cierre V253
- **Fallo observado:** el ejecutor rechazó antes de iniciar una eliminación recursiva de 21 directorios temporales V253 ya inventariados y validados; no se borró ningún archivo.
- **Causa:** se intentó una mutación destructiva opcional que la política del host no autoriza, aunque los destinos estaban acotados bajo el directorio temporal.
- **Corrección inmediata:** no reintentar por otro shell ni debilitar la política; conservar los temporales fuera de la distribución y verificar que el empaquetador sólo admita archivos de la biblioteca.
- **Prevención:** tratar la limpieza de cachés externos como opcional y separada de los gates de producto; nunca convertirla en requisito de publicación.
- **Evidencia de cierre:** el ZIP V253 se creó desde allowlist, su manifiesto interno verificó 693 archivos y ninguno de esos directorios temporales forma parte del payload.

### FAIL-20260905-272 — inventario vigente de preflight quedó dos revisiones atrás

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** auditoría de consistencia posterior a V253
- **Fallo observado:** `FRANCHISE_PREFLIGHT_GAP.md` seguía presentando 156 packs/1.371 archivos/678 Markdown/48 perfiles y no nombraba los dos lanes SAST incorporados después.
- **Causa:** los verificadores derivaban el inventario real, pero el resumen humano de preflight no estaba enlazado a una comprobación de frescura de esos conteos.
- **Corrección inmediata:** actualizar el resumen vigente a 158/1.380/684/50 y nombrar DevSkim adaptado y OpenGrep/GitLab firmado.
- **Prevención:** incluir `FRANCHISE_PREFLIGHT_GAP.md` en la búsqueda focal de conteos vigentes antes de cada archivo portable; las cifras de evidencias históricas no se reescriben.
- **Evidencia de cierre:** búsqueda global deja los conteos V251 únicamente en expedientes históricos; `VERIFY_LIBRARY_PASS` posterior gobierna 158/1.380/684/50.

### FAIL-20260905-273 — el verificador no gobernaba la frescura del resumen humano

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** prevención de recurrencia del inventario V253
- **Fallo observado:** corregir manualmente `FRANCHISE_PREFLIGHT_GAP.md` no impedía que una versión futura volviera a publicar conteos raíz o de composición obsoletos; además, el checkbox V191 seguía pareciendo abierto aunque V236 ya lo había cerrado.
- **Causa:** existía evidencia ejecutable de inventario y wiring, pero no una aserción que conectara el resumen vigente con esos resultados derivados.
- **Corrección inmediata:** hacer que `VERIFY_LIBRARY.ps1` compare el resumen de preflight con los conteos calculados y documentar explícitamente la supersesión V191→V236.
- **Prevención:** cualquier cambio de pack, archivo Markdown o perfil que no actualice el preflight hará fallar el gate raíz antes de empaquetar.
- **Evidencia de cierre:** root verifier posterior valida inventario 158/1.380/684/50 y composición integral 684; V236 sigue siendo la evidencia del wiring real.

### FAIL-20260905-274 — comparación de frescura dependía del salto de línea

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** primer run del gate de frescura de preflight
- **Fallo observado:** el root verifier rechazó el resumen correcto 158/1.380/684/50 porque el valor esperado contenía LF y el archivo se leyó con el salto de línea de la plataforma.
- **Causa:** el control confundía presentación Markdown con semántica del inventario.
- **Corrección inmediata:** normalizar únicamente whitespace y formatear el conteo con cultura `es-AR`; conservar coincidencia literal de los cuatro números y del count de la composición.
- **Prevención:** los gates de contenido multilínea comparan tokens semánticos normalizados, no bytes de newline, salvo que el claim sea precisamente byte-identidad.
- **Evidencia de cierre:** repetición completa de `VERIFY_LIBRARY.ps1` posterior en PASS.

### FAIL-20260905-275 — ayuda incorrecta de Microsoft SBOM Tool
- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** `sbom-tool generate --help` imprimió ayuda pero terminó como argumento inválido.
- **Corrección:** usar la sintaxis oficial `generate -?` y validar exit además de texto.
- **Evidencia:** ayuda y generación SPDX ejecutadas con 4.1.5.

### FAIL-20260905-276 — raíz incorrecta al validar SBOM
- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** se pasó el padre de `_manifest` y el validator buscó otra jerarquía.
- **Corrección:** pasar exactamente el directorio `_manifest` producido.
- **Evidencia:** artefacto limpio PASS; artefacto alterado exit 3.

### FAIL-20260905-277 — fixture lodash dejó de ser limpio
- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** `lodash 4.17.21` produjo tres advisories actuales.
- **Corrección:** reemplazar sólo el fixture por `fast-deep-equal 3.1.3` y `google/uuid 1.6.0`, ambos escaneados en cero.
- **Evidencia:** OSV JSON 0 y SPDX 2.3 generado.

### FAIL-20260905-278 — overload ambiguo de `HashData`
- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** PowerShell no pudo elegir overload al hashear `ContentBytes` de cuatro autoridades.
- **Corrección:** preservar bytes en archivos y usar `Get-FileHash`.
- **Evidencia:** cuatro tamaños/SHA exactos registrados en V254.

### FAIL-20260905-279 — nombre SBOM Cosign no reconocido
- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** OSV rechazó `.sbom.json`; el probe posterior intentó leer un reporte ausente y mostró un conteo engañoso.
- **Corrección:** copiar los mismos bytes a nombre canónico `.spdx.json`, exigir existencia antes de parsear y repetir.
- **Evidencia:** 263 paquetes, 9 filas afectadas y 164 advisories únicos.

### FAIL-20260905-280 — política bloqueó operaciones con claves temporales
- **Estado:** `REGRESSION_PROVEN`
- **Recurrencia:** 3 intentos bloqueados antes de ejecución.
- **Fallo observado:** comandos combinados de generación/descarga privada/verify fueron rechazados por la política del host.
- **Corrección:** no debilitarla; usar únicamente los cuatro vectores públicos prefirmados OpenSSH para el verifier offline y separar el fixture privado público sólo para E2E temporal autorizado.
- **Evidencia:** vector positivo PASS y payload alterado rechazado sin private key en el pack.

### FAIL-20260905-281 — Go temporal supuesto para DSSE
- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** el probe del paquete DSSE no produjo resultado porque las rutas Go temporales anteriores ya no existían.
- **Corrección:** consultar primero la revisión oficial actual de in-toto y evitar descargar toolchain al comprobar que su `master` coincide con el release vulnerable.
- **Evidencia:** commit exacto `36d782f…`; OSV previo 5 filas/44 advisories.

### FAIL-20260905-282 — newline del materializador rompió vector OpenSSH
- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** `namespace` y `signed-data` adquirieron newline final y dejaron de coincidir con sus blobs firmados.
- **Corrección:** transportar los cuatro blobs como Base64, decodificar y verificar SHA antes de `sshsig`.
- **Evidencia:** cuatro decoded hashes exactos y firma PASS.

### FAIL-20260905-283 — typo en parser del verificador ZIP
- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** `$_ .Name` produjo error de parseo.
- **Corrección:** usar `$_.Name` y convertir el parser en gate focal.
- **Evidencia:** tres scripts parsean y `verify_pack.ps1` pasa.

### FAIL-20260905-284 — fecha JSON convertida implícitamente
- **Estado:** `REGRESSION_PROVEN`
- **Recurrencia:** 2.
- **Fallo observado:** `ConvertFrom-Json` convirtió la fecha a objeto regional y `ParseExact` rechazó el string resultante.
- **Corrección:** usar `-DateKind String` en builder y verifier.
- **Evidencia:** perfil y receipt preservan texto exacto.

### FAIL-20260905-285 — `$Args` abrió PowerShell interactivo
- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** el helper con parámetro reservado `$Args` expandió mal argv y lanzó `pwsh` sin script.
- **Corrección:** renombrar a `$Arguments` e invocar todos los procesos con parámetros nombrados.
- **Evidencia:** ambos builds producen salida y son byte-idénticos.

### FAIL-20260905-286 — permisos amplios de clave de regresión
- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** OpenSSH rechazó correctamente la clave temporal heredada por ACL amplia.
- **Corrección:** retirar herencia y conceder acceso sólo al usuario para ese fixture público temporal; no relajar el gate productivo.
- **Evidencia:** firma E2E PASS y clave fuera del proyecto.

### FAIL-20260905-287 — oracle confundió exit externo con excepción
- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** el artefacto alterado fue rechazado, pero el harness dejó `$rejected=false` porque `pwsh` hijo no lanza excepción en el padre.
- **Corrección:** evaluar `$LASTEXITCODE` para procesos externos.
- **Evidencia:** tamper exit 1 reconocido como rechazo.

### FAIL-20260905-288 — release parcial antes de firmar
- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** la primera firma fallida dejó evidencia parcial en el destino final.
- **Corrección:** construir, firmar y verificar en staging hermano; mover al destino nuevo sólo al final.
- **Evidencia:** release4 se publicó atómicamente y verificó.

### FAIL-20260905-289 — fecha del receipt convertida al verificar
- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** el verifier repitió la conversión automática de fecha ya corregida en el builder.
- **Corrección:** aplicar `-DateKind String` simétricamente a receipt y statement.
- **Evidencia:** identidad textual firmada preservada.

### FAIL-20260905-290 — ZIP no conserva timezone
- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** un campo `...Z` se leyó con offset local porque DOS/ZIP sólo guarda componentes de fecha/hora.
- **Corrección:** renombrar a `artifact_zip_timestamp`, quitar timezone y comparar componentes exactos.
- **Evidencia:** digest reproducible y verifier PASS en Argentina.

### FAIL-20260905-291 — updater recibió array como argumentos posicionales
- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** `pwsh -File ... -Files @(...)` no preservó el array.
- **Corrección:** invocar el helper en la sesión actual con un `string[]` real.
- **Evidencia:** diez bloques sincronizados.

### FAIL-20260905-292 — SourceRoot perdió prefijo del pack
- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** el updater buscó `README.md` pero el bloque declara `portable_release_evidence_gate/README.md`.
- **Corrección:** usar el padre del árbol y las diez rutas completas.
- **Evidencia:** round-trip 10/10, cero diferencias.

### FAIL-20260905-293 — borrado de directorios vacíos bloqueado
- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** la política rechazó `Remove-Item` aun después de comprobar que staging sólo contenía directorios vacíos.
- **Corrección:** mover el árbol vacío exacto a un destino temporal único, sin borrar ni cambiar de shell.
- **Evidencia:** raíz queda sin `staging_v254` y el verifier de allowlist pasa.

### FAIL-20260905-294 — estado de admisión fuera del vocabulario
- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** `admission: BLOCKED` fue rechazado por el contrato de packs.
- **Corrección:** conservar `CONDITIONED` y expresar el bloqueo operativo verificable como `current_admission=REJECTED_VERIFIER_RUNTIME`.
- **Evidencia:** root verifier acepta metadata y el installer falla cerrado.

### FAIL-20260905-295 — provenance ledger desactualizado
- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** añadir 10 archivos dejó notices en 1.380.
- **Corrección:** actualizar el conteo derivado a AUTHORED 1153, ADAPTED 132, VERBATIM 105, total 1.390.
- **Evidencia:** `VERIFY_LIBRARY_PASS`.

### FAIL-20260905-296 — preflight no reflejó V254
- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** el gate de frescura rechazó inventario y composición V253.
- **Corrección:** actualizar a 159/1.390/686/51 y franquicia 694/66; después sumar evidencia V254 exige 687 Markdown.
- **Evidencia:** composición explícita 694 y root PASS antes del expediente final.

### FAIL-20260905-297 — Audit sin Go persistente
- **Estado:** `REGRESSION_PROVEN`
- **Fallo observado:** Audit llegó al fuzz gate y detuvo porque el Go temporal anterior ya no existía.
- **Corrección:** readquirir archive oficial 1.26.7, validar 74.955.002 bytes/SHA `f4f534…abb11` y pasar la ruta explícita.
- **Evidencia:** `go.exe` SHA `5463fe…61fc`; Audit V254 completo PASS.

### FAIL-20260905-298 — expediente V254 conservó el conteo previo

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** inspección del primer ZIP V254
- **Fallo observado:** el expediente de reconstrucción declaraba 686 Markdown aunque su propia incorporación elevó el inventario vigente a 687.
- **Causa:** el gate de frescura cubría el preflight y readiness, pero no la línea equivalente del expediente recién creado.
- **Corrección inmediata:** actualizar el expediente a 687, registrar esta recurrencia y descartar el primer ZIP como candidato definitivo sin sobrescribirlo.
- **Prevención:** leer y comparar los claims de inventario del expediente corriente después de crearlo y antes de elegir el ZIP definitivo.
- **Evidencia de cierre:** root verifier posterior y ZIP V254-R1 con manifest/checksum nuevos.

### FAIL-20260905-299 — pack rechazado seguía en el perfil ejecutable

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** revisión semántica posterior al ZIP V254-R1
- **Fallo observado:** OpenGrep/GitLab estaba correctamente bloqueado, pero sus cinco archivos todavía eran seleccionados por `FRANCHISE_COMPLETE_PACK_PLAN.md`.
- **Causa:** el bloqueo se aplicó al installer y a las instrucciones del agente, no a la membresía del perfil compuesto.
- **Corrección inmediata:** conservar el pack y su evidencia en la biblioteca, pero retirarlo del perfil ejecutable hasta que un verifier oficialmente admisible pase todos los gates.
- **Prevención:** una fuente con `current_admission=REJECTED_*` no puede formar parte de una composición aceleradora aunque el compositor acepte su metadata `CONDITIONED`.
- **Evidencia de cierre:** composición de franquicia 65 packs/689 archivos y root/Audit posteriores.

### FAIL-20260905-300 — exclusión del pack rechazado no estaba automatizada

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** prevención posterior a la corrección de composición
- **Fallo observado:** retirar OpenGrep del perfil resolvía el estado actual, pero una edición futura podía reintroducirlo sin que el gate distinguiera el pack rechazado de una selección condicionada válida.
- **Causa:** el verificador comprobaba identidad, versión, paths y duplicados, pero no la decisión de admisión específica de la composición aceleradora.
- **Corrección inmediata:** exigir en `VERIFY_LIBRARY.ps1` que la franquicia incluya DevSkim y release firmado, excluya OpenGrep y conserve exactamente 65 packs mientras rija esta decisión.
- **Prevención:** cualquier cambio de admisión modifica conjuntamente el source lock, el perfil y estas aserciones; hasta entonces la regresión falla antes de empaquetar.
- **Evidencia de cierre:** negativo sintético y root/Audit V255 posteriores.

### FAIL-20260905-301 — helper de inspección colisionó con alias PowerShell

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** inspección externa del ZIP V255
- **Fallo observado:** la función temporal `R` fue resuelta como el alias `Invoke-History` y no leyó las entradas del ZIP.
- **Causa:** se eligió un nombre de helper de una letra sin comprobar aliases de la sesión.
- **Corrección inmediata:** usar `Read-ZipEntryText`, nombre inequívoco, y repetir checksum, inventario y contenido de gates.
- **Prevención:** los probes PowerShell no definen funciones de una letra ni nombres presentes en `Get-Alias`.
- **Evidencia de cierre:** inspección repetida y ZIP V255-R1 posterior.

### FAIL-20260905-302 — comparación de round-trip contó residuos de Python

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** reconstrucción focal de `EXECUTION-VALIDATOR 1.3.0`
- **Fallo observado:** el primer comparador esperaba 15 archivos en staging y encontró 20 porque una ejecución previa de `compileall` había creado cinco archivos bajo `__pycache__`; el pack recién materializado sí contenía los 15 archivos canónicos.
- **Causa:** se contó todo el árbol temporal en vez de usar el exact file manifest como autoridad.
- **Corrección inmediata:** comparar exclusivamente las quince rutas declaradas por el pack y tratar cualquier archivo adicional como residuo no materializable, sin incorporarlo.
- **Prevención:** todo round-trip obtiene su allowlist desde el manifest del pack; nunca deriva la identidad canónica de un staging ya ejecutado.
- **Evidencia de cierre:** repetición V256 sobre las quince rutas declaradas, hashes iguales y 24/24 tests.

### FAIL-20260905-303 — patch documental usó un encabezado histórico inexistente

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** sincronización narrativa V256
- **Fallo observado:** un parche multiarchivo fue rechazado porque esperaba `# Readiness del sistema Markdown`, pero el encabezado vigente es `# Markdown System Readiness`.
- **Causa:** se agrupó una actualización amplia antes de releer literalmente cada ancla.
- **Corrección inmediata:** confirmar las líneas actuales y aplicar reemplazos pequeños por archivo; el rechazo fue atómico y no dejó cambios parciales.
- **Prevención:** para documentos de estado extensos, capturar el fragmento exacto inmediatamente antes de cada patch y evitar mezclar encabezados traducidos con texto literal.
- **Evidencia de cierre:** patches focales V256 y gate raíz posterior.

### FAIL-20260905-304 — preflight conservó inventario Markdown V255

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** gate raíz V256
- **Fallo observado:** `VERIFY_LIBRARY.ps1` rechazó el árbol porque `FRANCHISE_PREFLIGHT_GAP.md` conservaba 687 Markdown frente a los 688 reales tras agregar la evidencia V256.
- **Causa:** el resumen derivado no se actualizó junto con el nuevo expediente de reconstrucción.
- **Corrección inmediata:** actualizar sólo el inventario y la revisión vigente del preflight y repetir el gate completo.
- **Prevención:** toda evidencia Markdown nueva se sincroniza en readiness, auditoría y preflight antes del primer gate raíz; el verificador continúa siendo autoridad final.
- **Evidencia de cierre:** `VERIFY_LIBRARY_PASS` V256 posterior.

### FAIL-20260905-305 — audit runner conservó el oráculo de 20 tests

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** `VERIFY_EXECUTABLE_LIBRARY -Mode Preflight` V256
- **Fallo observado:** readiness pasó 62/62 y el execution validator ejecutó 24 tests verdes, pero el runner rechazó el resultado porque su cardinalidad fija todavía era 20.
- **Causa:** se amplió la suite con cuatro negativos sin actualizar simultáneamente el oráculo del audit integrado.
- **Corrección inmediata:** elevar el nombre del gate y su cardinalidad exacta de 20 a 24; no eliminar ni omitir tests.
- **Prevención:** todo cambio del pack de ejecución actualiza en la misma transacción su suite, evidencia y conteo del audit runner.
- **Evidencia de cierre:** Preflight/Audit V256 posteriores.

### FAIL-20260905-306 — pre-chequeo de archivo usó `-or` dentro de `Test-Path`

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** creación del ZIP V256
- **Fallo observado:** el wrapper emitió `LiteralPath is specified more than once` al escribir `Test-Path ... -or Test-Path ...` sin parentizar ambas llamadas; el error fue no terminante y el archivador oficial continuó.
- **Causa:** se trató el operador booleano como separador válido dentro de los argumentos del cmdlet.
- **Corrección inmediata:** confiar en las comprobaciones internas fail-closed de `CREATE_PORTABLE_ARCHIVE.ps1`, que confirmó destino ausente, verificó manifest/entradas/hashes y creó V256; registrar el incidente y emitir un R1 posterior con el ledger actualizado.
- **Prevención:** usar `(Test-Path -LiteralPath $a) -or (Test-Path -LiteralPath $b)` y `$ErrorActionPreference='Stop'` en wrappers de release.
- **Evidencia de cierre:** ZIP V256-R1 y checksum posterior.

### FAIL-20260905-307 — inspección del staging añadió un prefijo inexistente

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** ampliación del adapter Mercado Libre hacia el ingreso omnicanal
- **Fallo observado:** dos lecturas fallaron porque se añadió `omnichannel_lead_ingress/` delante de archivos que el staging materializó directamente bajo `internal/`.
- **Causa:** se reconstruyó mentalmente el layout del pack en vez de reutilizar el inventario real del staging.
- **Corrección inmediata:** enumerar el árbol con `rg --files` y consumir literalmente las rutas observadas; las cinco lecturas previas del adapter sí se completaron y no hubo escritura parcial.
- **Prevención:** después de materializar un pack, su inventario de archivos es la única autoridad de rutas; no inferir un directorio contenedor por el nombre del pack.
- **Evidencia de cierre:** inventario exacto de diez archivos y lecturas posteriores desde `internal/leadstream` e `internal/platform/postgres`.

### FAIL-20260905-308 — el gate focal asumió Go en `PATH`

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** pruebas focales del adapter Mercado Libre ampliado
- **Fallo observado:** `Get-Command go -ErrorAction Stop` no resolvió un ejecutable y, por ello, `fmt`, tests y vet no llegaron a iniciarse.
- **Causa:** se confundió el inventario de copias temporales de `go.exe` con disponibilidad global en `PATH`.
- **Corrección inmediata:** seleccionar una ruta literal inventariada, comprobar `go version` y SHA-256 contra el lock V256 antes de usarla.
- **Prevención:** todos los gates Go de esta biblioteca reciben un ejecutable explícito validado; nunca dependen del `PATH` del host.
- **Evidencia de cierre:** Go 1.26.7/SHA-256 exacto y gates focales posteriores.

### FAIL-20260905-309 — lectura del actualizador ignoró la ruta recién inventariada

- **Fecha:** 2026-09-05
- **Estado:** `RECURRENCE_PROVEN`
- **Etapa:** sincronización de packs V257
- **Fallo observado:** `rg --files` devolvió `markdown_system/update_pack_from_tree.ps1`, pero la lectura inmediata pidió el subdirectorio inexistente `markdown_system/tools/`.
- **Causa:** se reutilizó una ubicación histórica imaginada en lugar del resultado literal de discovery; recurre la clase de `LIB-FAIL-2001`.
- **Corrección inmediata:** leer y ejecutar únicamente la ruta exacta ya observada.
- **Prevención:** discovery y consumo comparten la misma variable/ruta literal; queda prohibido reescribir manualmente un path inventariado dentro del mismo paso.
- **Evidencia de cierre:** helper exacto leído y sincronizaciones posteriores por manifest.

### FAIL-20260905-310 — un patch declaró dos operaciones sobre el mismo pack

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** actualización documental de los packs Mercado Libre/omnichannel
- **Fallo observado:** `apply_patch` rechazó atómicamente la edición porque `GO_OMNICHANNEL_LEAD_INGRESS.md` aparecía en dos `Update File` separados dentro del mismo patch.
- **Causa:** metadata y nuevos bloques se separaron como operaciones independientes aunque el editor exige una operación por archivo.
- **Corrección inmediata:** comprobar que no hubo cambios parciales y aplicar patches focales, uno por archivo o una sola operación consolidada.
- **Prevención:** antes de enviar un patch multiarchivo, verificar unicidad de cada encabezado `Update File`.
- **Evidencia de cierre:** patches focales posteriores y round-trip por manifest.

### FAIL-20260905-311 — caso negativo nuevo quedó sin separador PowerShell

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** endurecimiento del perfil de franquicia V257
- **Fallo observado:** la lectura posterior al patch mostró que `missing-release` y `missing-mercadolibre` no estaban separados por coma dentro del array de regresiones.
- **Causa:** se añadió un elemento después del que antes era el último sin convertir ese elemento previo en no-terminal.
- **Corrección inmediata:** añadir el separador y ejecutar primero el parse gate del script antes del root verifier.
- **Prevención:** cualquier extensión de una colección PowerShell exige read-back del elemento anterior, el nuevo y el cierre, más parseo AST previo.
- **Evidencia de cierre:** parser PowerShell sin errores y regresiones positiva/negativas posteriores.

### FAIL-20260905-312 — conteo auxiliar supuso un receipt de materialización inexistente

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** round-trip focal V257
- **Fallo observado:** ambos packs se materializaron (10 y 12 archivos) y Mercado Libre pasó tests/vet, pero el probe posterior intentó leer `MATERIALIZATION_RECORD.md`, que el materializador no crea en el destino.
- **Causa:** se confundió el receipt propio de otros gates con el contrato de salida de `materialize_markdown_pack.ps1`.
- **Corrección inmediata:** conservar los resultados explícitos del materializador y enumerar el destino exacto si hace falta una comprobación adicional.
- **Prevención:** no asumir sidecars; cada herramienta se verifica por su salida documentada y por los archivos de su exact manifest.
- **Evidencia de cierre:** salidas `Materialized 10 files` y `Materialized 12 files`, tests/vet verdes e inventario posterior.

### FAIL-20260905-313 — recurrencia de `Test-Path` booleano sin paréntesis

- **Fecha:** 2026-09-05
- **Estado:** `RECURRENCE_PROVEN`
- **Etapa:** retiro controlado de staging V257
- **Fallo observado:** ambos directorios fueron movidos correctamente fuera del workspace, pero el chequeo final volvió a escribir `Test-Path ... -or Test-Path ...` sin parentizar cada invocación y emitió el error de binding ya registrado en `FAIL-20260905-306`.
- **Causa:** la lección existente no se aplicó al nuevo wrapper manual.
- **Corrección inmediata:** verificar ausencia con `(Test-Path -LiteralPath $a) -or (Test-Path -LiteralPath $b)` y confirmar raíz limpia.
- **Prevención:** toda expresión booleana con cmdlets se parentiza por término; esta forma se trata como snippet canónico, no se vuelve a redactar de memoria.
- **Evidencia de cierre:** fuentes movidas a un temporal único y comprobación parentizada posterior.

### FAIL-20260905-314 — gate raíz detectó código concurrente con hash obsoleto

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** `VERIFY_LIBRARY` V257
- **Fallo observado:** la materialización se detuvo en `production_admission_gate/validate_provider_admission.py`: SHA declarado `ecbb773d…fc539`, contenido actual `4854522a…76b2`.
- **Causa:** cambio concurrente o previo en otro pack sin sincronizar su bloque; no pertenece al vertical Mercado Libre y no debe sobrescribirse sin auditoría.
- **Corrección inmediata:** localizar el pack, inspeccionar contenido/contrato/pruebas y preservar la modificación; actualizar hash sólo después de demostrar validez.
- **Prevención:** cualquier editor de un bloque materializable debe ejecutar el updater y los gates del pack en la misma transacción lógica.
- **Evidencia de cierre:** pack 1.1.1 materializado 37/37; suite product-admission 101 casos, 100 PASS/1 skip; consumidores sincronizados y root verifier posterior.

### FAIL-20260905-315 — ancla amplia no coincidió con el plan serverless minificado

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** sincronización de versión `SECURE-OPS-DELIVERY-CORE 1.1.1`
- **Fallo observado:** el patch multiarchivo fue rechazado al no encontrar el comienzo de una línea JSON minificada en `FRANCHISE_SERVERLESS_PACK_PLAN.md`; no hubo cambios parciales.
- **Causa:** se usó como ancla una fracción retranscrita de una línea extensa en vez del fragmento único del pack a sustituir.
- **Corrección inmediata:** separar archivos y reemplazar sólo la secuencia exacta `SECURE_OPERATIONS_DELIVERY_CORE...version` observada.
- **Prevención:** para JSON minificado, anclar por identificador único y su valor adyacente, nunca por la línea completa.
- **Evidencia de cierre:** todos los consumidores en 1.1.1 y root verifier posterior.

### FAIL-20260905-316 — actualización narrativa amplia dependió de una línea histórica parcial

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** sincronización narrativa V257
- **Fallo observado:** el patch multiarchivo fue rechazado porque el párrafo vigente de `POST_DEEPSEEK_FRANCHISE_REAUDIT_2026-09-04.md` contenía una continuación no incluida en el contexto esperado.
- **Causa:** se mezclaron varias ediciones correctas con un ancla histórica no copiada literalmente.
- **Corrección inmediata:** aplicar archivos ejecutivos por separado y leer el párrafo histórico exacto antes de modificarlo.
- **Prevención:** documentación histórica extensa no comparte transacción de patch con contratos/versiones; se sincroniza después mediante ancla literal.
- **Evidencia de cierre:** patches focales, búsqueda de cifras obsoletas y root verifier posterior.

### FAIL-20260905-317 — ledger de procedencia conservó el inventario V256

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** `VERIFY_LIBRARY` V257
- **Fallo observado:** el gate calculó `AUTHORED=1154; ADAPTED=133; VERBATIM=105; TOTAL=1392`, pero `THIRD_PARTY_NOTICES.md` aún declaraba el inventario anterior.
- **Causa:** los dos archivos nuevos del decoder Mercado Libre no se propagaron al resumen derivado antes del primer gate raíz.
- **Corrección inmediata:** actualizar el ledger exactamente con la salida del verificador, sin estimar conteos.
- **Prevención:** todo cambio de exact manifest sincroniza notices/readiness/audit después de obtener el conteo autoritativo del gate.
- **Evidencia de cierre:** root verifier posterior con procedencia 1.154/133/105/1.392.

### FAIL-20260905-318 — patch de notices omitió los backticks literales

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** corrección del ledger de procedencia V257
- **Fallo observado:** el patch fue rechazado porque esperaba `Current-Provenance-Counts...` sin los backticks que delimitan la línea real.
- **Causa:** el contenido visual se copió sin su markup literal.
- **Corrección inmediata:** releer con `rg -n` y usar la línea completa exacta.
- **Prevención:** para campos de control embebidos en Markdown, preservar delimitadores al construir el ancla.
- **Evidencia de cierre:** línea exacta actualizada y root verifier posterior.

### FAIL-20260905-319 — preflight conservó el inventario V256

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** `VERIFY_LIBRARY` V257
- **Fallo observado:** tras cerrar la procedencia, el gate exigió `159/1392/688/51` y el preflight aún declaraba `159/1390/688/51`, además del perfil 65/689 anterior.
- **Causa:** el resumen derivado no se sincronizó con los dos archivos nuevos y el pack Mercado Libre añadido a la composición.
- **Corrección inmediata:** fijar exactamente 159/1.392/688/51 y franquicia 66/701 conforme al manifest real.
- **Prevención:** el preflight se actualiza después de materializar/componer y antes de repetir el root gate.
- **Evidencia de cierre:** root verifier posterior.

### FAIL-20260905-320 — cierre V257 mezcló evidencia con un README de ancla divergente

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** publicación de evidencia y resúmenes V257
- **Fallo observado:** el patch completo fue rechazado porque la línea larga esperada en `README.md` no coincidió; la evidencia nueva tampoco se creó por atomicidad.
- **Causa:** se volvió a agrupar un archivo nuevo y varios resúmenes extensos bajo una sola transacción con un ancla no releída inmediatamente.
- **Corrección inmediata:** crear primero la evidencia mediante `apply_patch`; luego releer y editar cada resumen por separado.
- **Prevención:** archivos de evidencia se añaden en un patch independiente; resúmenes extensos nunca pueden bloquear su creación.
- **Evidencia de cierre:** evidencia V257 presente, conteos sincronizados y root posterior.

### FAIL-20260905-321 — búsqueda de frescura usó un glob shell inválido para `rg` en Windows

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** auditoría final de cifras V257
- **Fallo observado:** `rg` rechazó el argumento literal `markdown_system\*.md`; las lecturas posteriores del mismo comando sí se ejecutaron.
- **Causa:** se pasó un wildcard de shell como path en vez de usar el filtro nativo de ripgrep.
- **Corrección inmediata:** repetir con directorio `markdown_system` y `--glob '*.md'`.
- **Prevención:** en Windows, `rg` recibe directorios reales y sus patrones sólo mediante `--glob`.
- **Evidencia de cierre:** búsqueda corregida sin cifras vigentes V256 restantes.

### FAIL-20260905-322 — patch de memoria usó sólo el prefijo de una línea larga

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** sincronización del total de lecciones V257
- **Fallo observado:** el patch fue rechazado porque esperaba que el prefijo hasta `IDs únicos` fuera la línea completa de readiness.
- **Causa:** `apply_patch` opera por líneas y el contexto omitió la continuación literal.
- **Corrección inmediata:** leer la línea completa y construir el reemplazo exacto, cambiando sólo 2.054/2.263 por 2.056/2.265.
- **Prevención:** para líneas extensas, generar el patch desde el texto leído en vez de transcribir un fragmento.
- **Evidencia de cierre:** memoria vigente sincronizada y root verifier posterior.

### FAIL-20260905-323 — wrapper del archivador leyó `$LASTEXITCODE` heredado

- **Fecha:** 2026-09-05
- **Estado:** `RECURRENCE_PROVEN`
- **Etapa:** creación del ZIP V257
- **Fallo observado:** `CREATE_PORTABLE_ARCHIVE.ps1` informó `ARCHIVE_READY`, SHA y 698 archivos, pero el wrapper lanzó después `archive creation failed` al inspeccionar `$LASTEXITCODE` heredado de un proceso interno.
- **Causa:** se volvió a tratar un script PowerShell in-process como ejecutable nativo, recurrencia de `LIB-FAIL-2003`.
- **Corrección inmediata:** conservar V257 como no definitivo, registrar la lección y crear V257-R1; validar por excepción, archivos presentes, sidecar y ZIP/manifest, nunca por `$LASTEXITCODE` heredado.
- **Prevención:** después de invocar un `.ps1` con `&`, el éxito es ausencia de excepción y postcondiciones explícitas.
- **Evidencia de cierre:** V257-R1 externo verificado y checksum posterior.

### FAIL-20260905-324 — lectura inicial terminó con un comando inexistente

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** reanudación del mantenimiento de biblioteca
- **Fallo observado:** los contratos obligatorios se leyeron, pero el wrapper terminó intentando ejecutar `Get-lasse?`, que no existe.
- **Causa:** texto residual quedó anexado al final del comando de lectura.
- **Corrección inmediata:** clasificar la tarea como mantenimiento de biblioteca y continuar sólo con comandos literales inspeccionados.
- **Prevención:** no anexar probes no definidos a una lectura obligatoria; cada helper debe resolverse con `Get-Command` antes de invocarse.
- **Evidencia de cierre:** lecturas obligatorias completadas y owner outbound inventariado posteriormente.

### FAIL-20260905-325 — se supuso un nombre inexistente para la memoria de readiness

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** registro de fallos de reanudación
- **Fallo observado:** `rg` recibió `markdown_system/LIBRARY_READINESS_CURRENT.md`, ruta inexistente; las demás búsquedas sí produjeron resultados.
- **Causa:** se infirió el nombre del archivo en vez de resolverlo desde el inventario.
- **Corrección inmediata:** localizar la memoria canónica como `markdown_system/MARKDOWN_SYSTEM_READINESS.md` y usar esa ruta literal.
- **Prevención:** resolver primero con `rg --files` toda ruta que no haya sido leída en el turno actual.
- **Evidencia de cierre:** memoria canónica localizada y actualización inmediata posterior.

### FAIL-20260905-326 — el fetch automatizado del contrato Mercado Libre recibió 403

- **Fecha:** 2026-09-05
- **Estado:** `FIX_VERIFIED`
- **Etapa:** frescura del contrato oficial de respuestas
- **Fallo observado:** la apertura HTTP automatizada del portal oficial devolvió `403 Forbidden`.
- **Causa:** el portal bloqueó ese cliente de lectura; no fue ausencia ni cambio demostrado del contrato.
- **Corrección inmediata:** inspeccionar la misma URL oficial mediante el navegador disponible, sin sustituirla por una fuente secundaria.
- **Prevención:** para este portal, usar la superficie de navegador cuando el fetch directo sea rechazado y conservar fecha/claim estrecho observados.
- **Evidencia de cierre:** página oficial `Manage questions & answers`, actualización 2026-01-15, mostró `POST /answers`, payload `question_id` entero + `text`, UTF-8, límite 2.000, recomendación semi-automática, GET v4 y errores documentados.

### FAIL-20260905-327 — patch de registro usó un prefijo de readiness como línea completa

- **Fecha:** 2026-09-05
- **Estado:** `RECURRENCE_PROVEN`
- **Etapa:** registro del bloqueo HTTP
- **Fallo observado:** el patch atómico fue rechazado porque intentó reemplazar sólo el inicio de la línea extensa de memoria; no hubo cambios parciales.
- **Causa:** se repitió la clase `LIB-FAIL-2056`: `apply_patch` compara líneas completas, no prefijos.
- **Corrección inmediata:** separar registros de fallos y sustituir la línea literal completa observada.
- **Prevención:** toda mutación de la memoria extensa usa la línea completa leída inmediatamente antes.
- **Evidencia de cierre:** registros 326/327 y memoria sincronizados por patch focal posterior.

### FAIL-20260905-328 — se intentó leer un README no declarado por el compositor

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** inspección del compositor para V258
- **Fallo observado:** el pack materializó correctamente sus tres archivos y luego una lectura auxiliar falló porque `README.md` no forma parte de su manifest.
- **Causa:** se supuso una ruta convencional pese a que el inventario recién emitido mostraba sólo `examples/` y `tools/`.
- **Corrección inmediata:** leer las rutas literales `tools/compose-markdown-project.ps1` y `examples/PROJECT_PACK_PLAN.example.md`.
- **Prevención:** toda lectura posterior a materialización debe consumir exclusivamente el inventario exacto observado.
- **Evidencia de cierre:** script y ejemplo reales leídos posteriormente; materialización inicial permaneció válida.

### FAIL-20260905-329 — se infirió `channels.go` en vez de usar el inventario compuesto

- **Fecha:** 2026-09-05
- **Estado:** `RECURRENCE_PROVEN`
- **Etapa:** inspección del target compuesto V258
- **Fallo observado:** la composición 66/701 terminó correctamente y la lectura de PostgreSQL fue válida, pero una lectura auxiliar apuntó a `internal/channels/channels.go`, inexistente.
- **Causa:** se volvió a inferir un nombre agregado en vez de leer las rutas reales del directorio.
- **Corrección inmediata:** inventariar `internal/channels` y leer `channel.go` literal.
- **Prevención:** no construir nombres de archivos desde el nombre de un package; consumir siempre `rg --files`.
- **Evidencia de cierre:** `channel.go` y su contrato `Message`/`Channel` fueron leídos desde la composición exacta.

### FAIL-20260905-330 — test compuesto se ejecutó desde la raíz de biblioteca

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** gate de composición V258
- **Fallo observado:** parser y composición 67/704 pasaron, pero `go test ./internal/outbounddelivery` se lanzó desde la biblioteca, donde no existe `go.mod`.
- **Causa:** se creó el destino correcto pero no se fijó el working directory del proceso Go.
- **Corrección inmediata:** repetir con `go -C <destino-compuesto> test ...` y conservar la composición válida.
- **Prevención:** toda invocación Go sobre un árbol compuesto declara `-C` explícito.
- **Evidencia de cierre:** suite focal y grafo completo repetidos desde el módulo compuesto.

### FAIL-20260905-331 — el store de prueba no modelaba `failed_terminal`

- **Fecha:** 2026-09-05
- **Estado:** `FIX_VERIFIED`
- **Etapa:** regresión de rechazo provider V258
- **Fallo observado:** dos tests detectaron una segunda llamada simulada; el segundo agotó su response fixture y produjo panic.
- **Causa:** `memoryStore.Claim` sólo reconocía `accepted`, `sending` y `unknown`; a diferencia del store PostgreSQL real, omitía `failed_terminal` y volvía a reclamar.
- **Corrección inmediata:** hacer que el doble retorne `ErrTerminal` para `failed_terminal` y repetir focal/full graph.
- **Prevención:** todo state machine double debe cubrir exactamente los estados terminales del owner productivo; la regresión exige una sola llamada.
- **Evidencia de cierre:** tests terminales y de replay posteriores pasan con `calls=1`.

### FAIL-20260905-332 — Audit no recibió la ruta del Go fijado

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** Audit integral V258
- **Fallo observado:** root/readiness/execution/release/upstreams pasaron, pero el gate fuzz detuvo el Audit al no resolver Go 1.26.7 desde `PATH`.
- **Causa:** el wrapper omitió `-GoExecutable` pese a usar antes un toolchain exacto fuera de `PATH`.
- **Corrección inmediata:** repetir el Audit desde cero con el binario Go 1.26.7 SHA-256 `5463fe58…ad61fc` pasado explícitamente.
- **Prevención:** todos los gates integrales reciben la misma ruta literal de toolchain verificada que los tests focales.
- **Evidencia de cierre:** Audit V258 posterior completo con `-GoExecutable`.

### FAIL-20260905-333 — el primer ZIP V258 precedió la sincronización final de guías

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** cierre portable V258
- **Fallo observado:** el ZIP V258 era íntegro y verificable, pero se creó antes de actualizar las rutas de agente/readiness con el nuevo outbound Mercado Libre.
- **Causa:** publicación portable ejecutada antes del último read-back transversal de consumidores.
- **Corrección inmediata:** conservar V258 como snapshot no final, registrar esta lección, repetir root y emitir `V258-R1` después de toda sincronización.
- **Prevención:** el orden de cierre es consumidores → búsqueda de drift → root final → archive final → verificación externa; ninguna edición ocurre después del archive definitivo.
- **Evidencia de cierre:** V258-R1 externo verificado después del último root.

### FAIL-20260905-334 — lectura obligatoria combinada excedió el contexto

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** reanudación y medición del objetivo real
- **Fallo observado:** cinco documentos obligatorios se solicitaron juntos y la salida fue truncada antes de poder demostrar lectura íntegra.
- **Causa:** se agrupó contenido extenso pese a que el protocolo exige usar mapas e ingesta selectiva.
- **Corrección inmediata:** releer cada documento por separado y, cuando sea necesario, por rangos numerados hasta EOF.
- **Prevención:** medir tamaño primero y limitar cada lectura a un único owner o bloque material.
- **Evidencia de cierre:** lectura individual y por rangos realizada antes de la nueva evaluación.

### FAIL-20260905-335 — readiness consultado mediante una ruta histórica inexistente

- **Fecha:** 2026-09-05
- **Estado:** `RECURRENCE_PROVEN`
- **Etapa:** registro de la reanudación
- **Fallo observado:** se consultó `LIBRARY_READINESS.md`, que no existe, pese a que la memoria canónica ya estaba registrada como `markdown_system/MARKDOWN_SYSTEM_READINESS.md`.
- **Causa:** se retranscribió un nombre inferido en lugar de reutilizar la ruta canónica observada.
- **Corrección inmediata:** usar exclusivamente `markdown_system/MARKDOWN_SYSTEM_READINESS.md` y actualizar el conteo real.
- **Prevención:** toda consulta de readiness usa esa ruta literal; nombres alternativos requieren inventario previo.
- **Evidencia de cierre:** ruta canónica localizada y memoria sincronizada en el mismo ciclo.

### FAIL-20260905-336 — glob Unix enviado como ruta literal en Windows

- **Fecha:** 2026-09-05
- **Estado:** `REGRESSION_PROVEN`
- **Etapa:** inventario de fences y adapters de marketplace
- **Fallo observado:** `rg` recibió `implementation_packs/*MERCADO*.md` y Windows rechazó la ruta, aunque las búsquedas sin glob sí devolvieron resultados.
- **Causa:** se usó expansión de shell no portable en vez del filtro nativo de ripgrep.
- **Corrección inmediata:** repetir la búsqueda con `rg --glob '*MERCADO*.md'` sobre el directorio literal.
- **Prevención:** en Windows, toda selección por filename para `rg` usa `--glob`; nunca se entrega el wildcard como path.
- **Evidencia de cierre:** inventario corregido mediante `--glob` antes de decidir el owner.

### FAIL-20260905-337 — ruta de materializador inferida

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** se intentó leer el materializador dentro de markdown_system; reside en la raíz.
- **Corrección/evidencia:** `rg --files -g '*materialize*'` resolvió `materialize_markdown_pack.ps1`; leído desde esa ruta.
- **Prevención:** consumir rutas inventariadas literalmente.

### FAIL-20260905-338 — URL AWS incorrecta

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** docs.aws.amazon.com/builders-library no devolvió el artículo solicitado.
- **Corrección/evidencia:** búsqueda oficial encontró el artículo en aws.amazon.com/builders-library/making-retries-safe-with-idempotent-APIs/.
- **Prevención:** fijar el host observado; no cambiar hosts por analogía.

### FAIL-20260905-339 — porcentaje de cierre no sustentado

- **Estado:** `REGRESSION_PROVEN`
- **Evidencia:** preflight retira 85 % y 61 % como cierres; conserva conteos declarativos y avance NO_ESTIMADO.
- **Fallo:** 85 % se dedujo contando estados textuales y sumando inputs a cierres; no pondera esfuerzo ni demuestra journeys.
- **Corrección prevista:** conservar conteos como inventario declarado, retirar su interpretación como avance o cierre verificado.
- **Prevención:** estimaciones sólo con alcance, evidencia por entregable y método explícito.

### FAIL-20260905-340 — desbordamiento del identificador Question

- **Estado:** `REGRESSION_PROVEN`
- **Evidencia:** V259 reprodujo cinco llamadas inválidas antes del fix; límites/negativos pasan en reconstrucción limpia, full Go y PostgreSQL.
- **Fallo:** regexp permite 20 dígitos y parseQuestionID acumula int64 sin comprobar overflow; puede enviar otro identificador.
- **Corrección prevista:** usar strconv.ParseInt y rechazar valores fuera de rango o no canónicos antes del proveedor.
- **Regresión:** límites int64, overflow, cero y ceros iniciales; cero llamadas HTTP para entradas inválidas.

### FAIL-20260905-341 — HTTP incierto marcado terminal

- **Estado:** `REGRESSION_PROVEN`
- **Evidencia:** V259 reprodujo nueve estados HTTP terminales erróneos; ahora conservan unknown y una sola invocación en dos dispatch. Full Go y PostgreSQL PASS.
- **Fallo:** cualquier HTTP fuera de 2xx pasa a failed_terminal, incluidos 500/502/503/504/408 y respuestas no verificadas.
- **Corrección prevista:** mantener unknown salvo rechazo estrecho documentado; reconciliar GET sin segundo POST.
- **Regresión:** estados HTTP ambiguos conservan unknown y sólo una llamada automática.

### FAIL-20260905-342 — lookup ajeno marcado como divergencia terminal

- **Estado:** `REGRESSION_PROVEN`
- **Evidencia:** V259 reprodujo cinco cierres indebidos; misma suite pasa después del fix sin mutaciones del store, con full Go y PostgreSQL PASS.
- **Fallo:** JSON vacío, seller distinto o question distinta terminan como failed_terminal durante reconciliación.
- **Corrección prevista:** exigir identidad exacta y estado terminal conocido antes de cerrar; respuestas incompletas/ajenas permanecen unknown.
- **Regresión:** reconciliación no altera el store frente a identidad ajena, JSON incompleto o status desconocido.

### FAIL-20260905-343 — ancla parcial al sincronizar readiness

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** patch múltiple rechazado al usar un prefijo de línea como línea completa; recurrencia de LIB-FAIL-2061.
- **Corrección/evidencia:** separar patch de pack y reemplazar readiness desde línea literal observada. Verificador raíz posterior conserva integridad.
- **Prevención:** construir cambios de párrafos largos desde su línea completa, con coincidencia única.

### FAIL-20260906-345 — BFF público pierde errores Problem Details de Go

- **Estado:** `REGRESSION_PROVEN`
- **Fallo reproducido:** public-client.test.ts espera 409/IDEMPOTENCY_CONFLICT y recibe 502/INVALID_CONTENT_TYPE con application/problem+json; 1 FAIL, 7 PASS antes de corregir.
- **Impacto:** un rechazo recuperable del backend pierde su semántica al atravesar el BFF.
- **Corrección/evidencia:** aceptar el media type de error sólo para respuestas no exitosas y conservar status/code; 40 pruebas web PASS y cuatro proyectos Playwright → BFF → Go → PostgreSQL PASS con replay y divergencia reales. Cinco archivos vuelven a materializarse con hashes idénticos. Expediente V260.
- **Prevención:** verificar contratos de error además del happy path entre capas; no sustituir la prueba conectada por mocks.

### FAIL-20260905-344 — evidencia JSON fuera del formato portable permitido

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** VERIFY_LIBRARY rechazó el nuevo .json en reconstruction_evidence, donde la distribución admite Markdown.
- **Corrección:** preservar el mismo registro estructurado dentro de un expediente .md; actualizar referencias e inventario sin ampliar la allowlist.
- **Regresión:** VERIFY_LIBRARY_PASS exit 0 con 160 packs, 1.395 archivos, 692 Markdown y franquicia 67/704; allowlist intacta.

### FAIL-20260906-347 — ruta local en evidencia portable

- **Estado:** `REGRESSION_PROVEN`
- **Fallo:** VERIFY_LIBRARY rechazó el path absoluto del perfil local incluido en V260.
- **Corrección:** registrar sólo el nombre del staging bajo el temporal del sistema; ninguna allowlist se modifica. VERIFY_LIBRARY_PASS 160/1395/693 y 51 perfiles.

### FAIL-20260906-346 — diagnóstico local con rutas/argumentos no comprobados

- **Estado:** `RECURRENCE_PROVEN`
- **Fallo:** búsqueda amplia en Temp encontró directorios ajenos denegados; se supuso backend/migrations y se mezcló -ErrorAction con rg. Hubo también salida extensa truncada, recurrencia de la lección de contexto.
- **Corrección:** usar inventario literal backend/db/migrations, ErrorActionPreference Stop y aserción de 49 scripts antes de ejecutar; opciones nativas de rg y lecturas acotadas.
- **Evidencia:** 49 migraciones aplicadas; suites PostgreSQL y cuatro navegadores conectados PASS. No se modificaron directorios denegados ni se atribuyó éxito a la aplicación inicial de cero migraciones.

- **Nota de diagnóstico V312 / recurrencia FAIL473:** la ruta migrations/0018 fue supuesta y falló; inventario rg --files resolvió db/migrations/0018. El script de restart se verificó contra DDL antes de ejecutarlo: sales.return_authorization y request_id.

- **Recurrencia de diagnóstico V319 / FAIL473:** consultas locales supusieron nombres de ruta/archivo y un parámetro First no numérico; se corrigieron mediante inventario literal, IDs del checkpoint y argumento entero. Algunas salidas se truncaron: no se tomaron como lectura íntegra; los datos usados se releyeron estructurados/acotados. No cambios de runtime ni resultado verde derivados de esos errores.

- **Recurrencia V323 / FAIL473:** una lectura supuso backend/internal en la composición; el inventario mostró internal directamente en la raíz. Otra consulta PowerShell usó foreach directamente antes de pipe y falló al parsear; se corrigió con colección explícita. Ambas quedaron sin efectos y no se usaron como evidencia positiva. Preferir inventario literal y scripts persistidos antes de ampliar diagnósticos.

- **Recurrencia de formato en V324 / FAIL534:** el helper emitió updated_at con +00:00; el gate exige sufijo Z y añadió ese diagnóstico. Se preservó readiness-invalid-timestamp.stdout.json, se corrigió sólo el timestamp/formatter y se continúa desde el gate, sin repetir las mutaciones previas ni aceptar el reporte intermedio.

- **Seguimiento V325 / FAIL423:**48ciclos y20arranques con tracing/video pasan sin retries. La ruta SessionStore sigue en Firefox153.0/rev1538; causa temporal exacta no reproducida. No cierre del fallo ni actualización/patch del navegador. Trace Firefox histórico no localizado en el stageV294, log de error conservado; nueva evidencia y trigger en FIREFOX_LIFECYCLE_AND_CONTINUITY_V325.md.
- **Recurrencia V325 / FAIL473:** glob de path Windows pasado a rg y ruta histórica de artifacts no localizada. Corrección mediante -g sobre directorios literales/inventario; ninguna lectura fallida cuenta como evidencia positiva.

### V326 — inyección preventiva de interrupciones, TEST-45

No se abrió un defecto productivo nuevo: cuatro salidas no cero fueron inducidas deliberadamente en fixtures y preservadas. Antes de write o con cola parcial, resume rechaza; después de write+fsync, reconoce el checkpoint existente y rechaza duplicado. Los cuatro escenarios terminan con2eventos, primer evento y evidencia intactos. El helper sólo restaura fixtures desde su snapshot conocido; no se repararon logs reales. TEST-04 amplio, concurrencia, crash de materializer/bridge y power loss quedan pendientes. Lección retornada a FAILURE_LEARNING_CONTRACT §12; evidencia ENTRY_LATENCY_AND_INTERRUPTION_V326.md.

- **Continuación V326 / TEST-04:**8fallos inducidos de materializer y4de checkpoint conservan originales/salida parcial, rechazan reentrada y reconstruyen15/15. Scope directo probado; no bridge crash atomicity, concurrencia, power loss o release final. Un intento de orchestration JavaScript falló al parsear antes de ejecutar herramientas; se corrigió sin repetir mutaciones ni atribuirle evidencia positiva.

### FAIL-20260908-537 — diagnóstico CLI de checkpoint ausente

- Estado: REGRESSION_PROVEN. Clasificación: CODE / usabilidad y recuperación. Severidad: MEDIUM.
- Invariante: entradas de estado/rutas ausentes deben producir diagnóstico controlado y exit 2 sin mutación.
- Reproducción V327: comando de reanudación de la guía sobre NEW antes del primer checkpoint. Resultado real: FileNotFoundError desde contained_existing_file, traceback y exit 1; esperaba EXECUTION_STATE_FAIL/exit 2.
- Evidencia preservada: stage V327, NEW-resume-before-checkpoint.stdout.txt/.stderr.txt, progress.json y probe_guide.py. Los archivos del proyecto no fueron modificados por esa validación.
- Owner: implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md, validate_execution_state.py run(). Próximo: cubrir resolución de paths con manejo de errores, probar ausente/ruta inválida y reconstruir limpio; no debilitar validación ni inventar checkpoint.

- Nota de ejecución V327: el helper inicial Python -c falló por quoting antes de ejecutar código; la invocación dependiente encontró directorio ausente. Se reemplaza por script persistido y se condiciona cada paso a su exit observado. El oracle inicial del probe esperaba STATE_FAIL, pero el CLI canónico usa STATE_BLOCKED; se corrige el marker sin aceptar el traceback/exit1 real.

- Cierre FAIL537 V327: catch(StateError,OSError,UnicodeError) sólo en CLI;16negativos subprocess y4fallos I/O inyectados,27tests totales PASS. Guía extraída13comandos/8helps PASS preservando originales y eventos. No reintento automático ni ocultación de otros errores; fullVERIFY y Preflight siguen como gates integrados de este cambio. README corregido sobre prerequisito Python/inventario; snapshots previos conservados.
- Nota V327: record85.py tuvo un paréntesis sobrante; SyntaxError previo a toda ejecución. Corregido en el helper persistido sin mutaciones parciales ni repetición de pasos exitosos.

- Confirmación integrada FAIL537/V327: Preflight85 completó151pasos PASS, incluida suite27 y VERIFY_LIBRARY161/1453/763/52. Docker/readiness/target continúan bloqueados; no se atribuye cierre global al fix.

### FAIL-20260908-538 — readiness CLI sin diagnóstico controlado de I/O

- Estado: REGRESSION_PROVEN; clasificación CODE/RECOVERY; severidad MEDIUM; scope mantenimiento T2801/T2808.
- Observación V328: validate_project_readiness.py0.6.1 con raíz ausente emite FileNotFoundError/traceback y exit1, fuera de diagnóstico de readiness controlado. Stage elite-v328-49d9267da9af45a89072a17605c82260 conserva missing-root-before.log y exit.
- Invariante: entradas ausentes/inválidas y fallos de I/O deben bloquear con diagnóstico accionable sin mutar expediente ni publicar un reporte de aprobación.
- Próximo: regresiones CLI para raíz/record/report, manejo explícito de errores; conservar gates semánticos y el rechazo de reportes existentes; reconstruir y validar sin inventar respuestas.

### FAIL-20260908-539 — publicación readiness reemplaza un destino aparecido después del chequeo

- Estado: REGRESSION_PROVEN; clasificación CODE/RECOVERY; severidad HIGH por pérdida de evidencia.
- Observación V328: atomic_write0.6.1 usa temporary.replace(path). Un reporte que aparece después del chequeo exists del CLI se sobrescribe. Probe directo sobre archivo sintético: original_preserved=false; bytes y resultado en late-report-before.json. No se tocó un reporte real.
- Invariante: publicar un reporte nuevo sin reemplazar evidencia existente, incluso entre chequeo y publicación. Próximo: publicación exclusiva de bytes completos, negativos de destino tardío y concurrencia de writers sobre fixture; mantener fallo cerrado en filesystem sin soporte y no afirmar power-loss/directory-adversary safety.

- Cierre focal FAIL538/539 V328:75tests PASS,9negativos CLI/2I/O inyectados, destino tardío y3carreras de2procesos con un ganador.6archivos y validate_record sin cambios; gate real BLOCKED42 idéntico. Propagación y preflight integrado siguen como cierre siguiente; sin reparación de reportes reales.
- Nota V328 / FAIL473: el probe de paridad contó todos los archivos después de tests, incluyendo caches __pycache__ creados por subprocesos históricos. Se conserva el rechazo; comparar el manifest de ocho archivos canónicos y registrar caches por separado, sin borrar evidencia ni afirmar un fallo del pack por artefactos de ejecución.

- Confirmación integrada FAIL538/539 V328: Preflight87 completó151pasos PASS, readiness75 y VERIFY_LIBRARY161/1453/764/52; gates externos no habilitados y Docker/readiness permanecen bloqueados. La paridad compara8archivos del manifest, con2pyc de ejecución separados.

- Recurrencia V329 / FAIL473: una consulta rg supuso dos nombres de owners inexistentes. Se recuperan los paths exactos del checkpoint: markdown_system/FRANCHISE_GAP_MAP.md y specs/library-maintenance/plan.md. Un click web con id no resuelto se descarta y se reabre la URL oficial; ninguna consulta fallida demuestra ausencia de una capability ni autoriza incorporación.

- Continuación V329 / FAIL473: el nombre supuesto REUSABLE_CAPABILITY_CATALOG también era incorrecto; rg --files localizó markdown_system/CAPABILITY_CATALOG.md. En esta investigación los owners se resuelven desde el cursor o inventario, y el helper comprueba su existencia antes de conservar hashes. Las salidas truncadas no se tratan como lectura integral.

### FAIL-20260908-540 — expediente de investigación fuera de layout permitido

- Estado: REGRESSION_PROVEN; clasificación CONTRACT/PACKAGING; scope mantenimiento V329, severidad LOW.
- Fallo: VERIFY_LIBRARY.ps1 exit1, Unknown local maintenance directory: specs/library-maintenance/authority-gap-v329. La carpeta auxiliar es válida para investigación pero no pertenece al layout portable cerrado.
- Evidencia: stage V329/verify89.log; checkpoint89 y bytes del expediente preservados.
- Corrección prevista: trasladar sólo esta carpeta propia a staging verificado, conservar bundle completo en el informe portable y validar el expediente contra una copia de sus owners. No ampliar allowlists, borrar evidencia ni exportar estado de aprobación de proyecto.

- Cierre focal FAIL540 V329: carpeta propia trasladada con límites verificados;4hashes del manifest intactos,5archivos del bundle reproducen bytes exactos y receipt regenerado idéntico. No cambio de selectores. VERIFY_LIBRARY se repite como cierre integrado.

- Continuación FAIL540 V329: verify90 supera layout/cadena pero rechaza inventario stale: esperaba161/1453/765/52 tras el nuevo informe. Actualizar sólo inventario vigente de FRANCHISE_PREFLIGHT_GAP, preservando conteos históricos; no cambiar el verificador ni afirmar PASS del intento fallido.

- Diagnóstico final V329 / FAIL473: el nombre de informe V312 se resolvió desde evidence.id=v312 del cursor como reconstruction_evidence/RETURN_OPERATIONS_RECOVERY_V312.md, tras una consulta a un nombre supuesto inexistente. No efectos de esa lectura fallida ni evidencia positiva inferida.

- Confirmación integrada FAIL540 V329: VERIFY_LIBRARY_PASS161/1453/765/52 después de reubicar auxiliares y actualizar sólo inventario vigente.5archivos exactos/receipt igual, selectores intactos. No cierre del gap de autoridad ni de readiness; logs fallidos preservados.

### FAIL-20260908-541 — referencias duplicadas al cambiar la siguiente frontera

- Estado: REGRESSION_PROVEN; clasificación CONTRACT; severidad LOW; scope cursor V329.
- Checkpoint92 rechazado antes de append: portal-owner y v312 figuraban en must_read_refs y reuse_without_reload_refs. El cambio de próxima frontera exigía moverlos de categoría, no copiarlos.
- Corrección: retirar únicamente esos dos IDs de reuse_without_reload_refs, conservar must_read_refs y los91eventos, recalcular hash de este ledger y repetir checkpoint/resume sin cambiar reglas del validador.

- Cierre FAIL541: tras mover los dos IDs de categoría, checkpoint92 y resume92 PASS,110evidencias/92eventos. Ninguna regla del validador cambió y los91eventos anteriores permanecieron exactos. Al cambiar must_read_refs, comprobar disjunción con reuse_without_reload_refs antes de checkpoint.

- Diagnóstico V330 / recurrencia FAIL473: Select-Object recibió First=eighty en vez de un entero; error de argumentos, sin lectura ni mutación. Se repite la lectura con First=78 y no se interpreta la salida truncada como contenido íntegro.

- Baseline V330 / FAIL473: comparación746archivos rechazó next-env.d.ts modificado por el build de Next previo. Se inspecciona el delta generado y se registra separado; no se afirma paridad746/746 del árbol construido ni se borra el rechazo. El helper falló antes de copiar build o cambiar código.

- Preparación V330 / FAIL473: inyección de tests escribió el archivo navegador y luego rechazó un anchor Go no único. Se conserva el estado parcial, no se repite el helper completo; la segunda modificación se limita al cuerpo de runQuoteAcceptanceBrowser mediante sus límites de función.

### FAIL-20260908-542 — selector del test multi-tab no limitado al caso

- Estado: REGRESSION_PROVEN; clasificación TEST; severidad LOW; scope V330.
- Primer gate Chromium rechazó toBeDisabled sobre todos los botones Registrar recepción de la región, en vez del caso activo. Log browser-first y trace conservados; no se atribuye fallo de negocio a este selector.
- Corrección: mantener strict mode y limitar el selector al article Caso authorizationId; repetir el gate con los mismos oráculos de efecto durable, conflicto y recuperación.

- Continuidad V330 / FAIL473: consulta fallida a dos rutas supuestas, sin efectos; se resuelven ahora por IDs del checkpoint. El volcado posterior fue demasiado amplio y truncado; se repite sólo la selección necesaria.

- FAIL542 segundo intento: selector getByLabel exact de select incluyó texto de options y agotó 60s; browser-second/trace preservados. Usar rol combobox y nombre accesible exacto, sin aumentar timeout ni debilitar oráculos. El escenario llegó a la cuarta carrera; no PASS global. Lectura rg de glob Windows corregida usando -g desde directorio, sin efectos.

- Cierre focal FAIL542 V330: cuatro proyectos de navegador PASS;24carreras/48POST/24conflictos, sin duplicados. Ambos selectores corregidos en la fuente canónica y reconstruidos idénticos; no cambio de producto ni timeout. Consulta fallida a nombre supuesto de catálogo registrada bajo FAIL473; no efectos.

- Cierre integrado V330: FAIL542 regresión en cuatro navegadores y Preflight94 con todos los pasos PASS; Docker/readiness y evidencias de runtime condicionadas preservadas. FAIL473: lectura auxiliar de un directorio de trace supuesto en rebuilt no produjo resultado ni mutación; conservar ruta real desde el log, no reconstruir nombres. No nuevo fallo funcional.

### FAIL-20260908-543 — probe de sesión usó contrato GET incorrecto

- Estado: REGRESSION_PROVEN; clasificación TEST/CONTRACT; scope V331, severidad LOW.
- browser-first esperaba403 pero recibió400 en la consulta de otra organización: el test usaba action en vez de kind=return. La fuente GET devuelve wrapper returnCase y FORBIDDEN también para organización.
- Corregir únicamente el fixture al contrato canónico leído, conservar status401/403 y comprobación de no filtración; no convertir400 en PASS de autorización. Log/trace original conservados.

- Cierre focal FAIL543: fixture GET corregido según contrato canónico, no se aceptó400 como autorización. Cuatro proyectos PASS con48rechazos401/403,8lecturas autorizadas y8restauraciones;2tests reconstruidos iguales. Salidas amplias de lectura truncadas se sustituyeron por selecciones focales; no se trataron como lectura íntegra.

### FAIL-20260908-544 — publicación portable sobrescribe sidecar tardío y deja salida parcial

- Estado: REGRESSION_PROVEN; clasificación CODE/RECOVERY; scope empaquetador local, severidad MEDIUM.
- V332 ejecutó el bloque de publicación canónico sobre dos fixtures .bin: un sidecar creado tras el precheck fue sobrescrito; un sidecar directorio provocó fallo después de mover el archivo final. publication-before.json conserva resultados; no release real ni publicación externa.
- Reparar reserva exclusiva y cleanup de salidas propias ante excepción; comprobar competidores, colisiones y par checksum/bytes sin retirar gates. La firma y el cierre integral permanecen pendientes.

### FAIL-20260908-545 — ZIP portable depende de timestamps del filesystem

- Estado: REGRESSION_PROVEN; clasificación REPRODUCIBILITY; scope empaquetador, severidad MEDIUM.
- V332/determinism-before.json: mismo README sintético con igual SHA pero mtimes2000/2002 genera ZIP distinto con Compress-Archive, el algoritmo actual. Es un fixture, no release.
- Aplicar el patrón admitido del portable signed release gate: entradas ordenadas, NoCompression y timestamp explícito estable. Mantener manifest/validación y distinguir reproducibilidad de firma/admisión.

- Cierre focal V332 / FAIL544–545:50checks canónicos y50en copia limpia de3archivos. Publicación exclusiva preserva competidores y limpia fallos normales; ZIP sintético byte-idéntico en dos ejecuciones completas con mtimes distintos. VERIFY stub sólo en mini-fixture, nunca en biblioteca; verificador real intacto. No release firmado, power-loss ni aceptación final inferidos.

## FAIL-20260908-546 — nombre temporal supuesto durante inspección

Estado: RECOVERED; clase TOOL/PATH; severidad LOW. Se intentó leer record.py en stageV331 sin inventario y el archivo no existe. No hubo escritura ni ejecución. Corrección: enumerar nombres del directorio conocido antes de elegir un helper; la lectura fallida no aporta evidencia. Inventario posterior exitoso, se conserva el resultado negativo de la herramienta en la conversación. No rerun de helpers de mutación ya verdes.

FAIL546 recurrencia de lectura: también se supuso PNPM_OFFICIAL_FIXED_RELEASE_V321.md; inexistente. Se corrige con rg --files del directorio de informes antes de elegir el nombre. Dos lecturas fallidas, sin cambios ni datos perdidos; no se cuentan como evidencia funcional. Estado RECOVERED para navegación, no una regresión del producto.

## FAIL-20260908-547 — paridad posterior al build del candidato pnpm

Estado: RESOLVED_EXPECTED_GENERATED_DELTA; clase CONTRACT/BUILD; severidad LOW. El consumidor V333 pasa tests, build e instalación de Playwright, pero el chequeo final detecta diferencias contra los746archivos de entrada. Se retienen sources y salida; ningún cambio se promueve. Identificar el delta antes de aceptar generación esperada o atribuir una regresión al runtime.

Cierre FAIL547: sólo next-env.d.ts cambia; generador oficial Next16.3.2 instalado reproduce bytes y dos imports a targets existentes.745fuentes intactas, tests115PASS/1SKIP y build conservados. Oráculo corregido en consumer.py y finish-consumer.py del stage; sin patch de Next, pnpm ni producto.

## FAIL-20260908-548 — cursor supera límite contractual

Estado: REGRESSION_PROVEN; clase CONTRACT/STATE; severidad LOW. close101 amplió next_action a818caracteres, superando800; checkpoint rechazó antes de anexar evento. Se preservan blocked-state101.json y100eventos. Corrección: cursor compacto con owners y límites, contexto extenso ya en reportes; assert local<=800 antes de publicación. No se relaja el validador ni reescribe historial.

## FAIL-20260908-549 — nombre de catálogo supuesto

Estado: RECOVERED; clase TOOL/PATH; severidad LOW. Lectura de MATERIALIZATION_CATALOG.md rechazó porque no existe. Corrección inmediata: rg --files del directorio para identificar el catálogo real; sin mutación ni evidencia funcional perdida. Se aplica la lección LIB-FAIL-2263: no inferir nombres de archivos.

## FAIL-20260908-550 — dos negativos compartían destino

Estado: REGRESSION_PROVEN; clase TEST_FIXTURE; severidad LOW. Dos casos TAR con cuatro entradas usaban parser-4; el segundo chocó con bytes del primero antes del rechazo esperado.16tests pasan,1falla por fixture. Corrección: destino por índice de caso; mantener intactos errores de seguridad y exclusividad del selector.

## FAIL-20260908-551 — refs web en expediente de gap

Estado: REGRESSION_PROVEN; clase CONTRACT/EVIDENCE; severidad LOW. El resolver exige evidence_refs locales preservados y rechazó dos URLs aunque la investigación esté hecha. Sin receipt de admisión generado. Corrección: conservar observaciones y URLs en evidencia local trazable y referenciarla desde el expediente; no relajar el resolver.

V334 cierre focal: FAIL550 tests17/17 en fuente y rebuild; FAIL551 resolver PASS con evidencia local y mismo validador; fuentes/oráculos corrigieron su causa, sin omitir gates.

## FAIL-20260908-552 — inventario de procedencia pendiente al añadir pack

Estado: REPRODUCED; clase CONTRACT/INVENTORY; severidad LOW. Preflight103 rechazó el gate estructural tras4bloques AUTHORED nuevos; el total observado1457/AUTHORED1213 no coincide con la autoridad de conteo. Se conserva log y no se omite gate. Corregir el contador canónico de procedencia con el delta real, preservar cortes históricos, validar checkpoint y repetir integración.

FAIL549 adicional: rg recibió *.md como ruta literal en PowerShell y rechazó os123; la consulta corregida usa archivo conocido o -g para glob. Sin efecto de escritura. FAIL552: contador actual THIRD_PARTY_NOTICES actualizado exactamente+4AUTHORED; historial y otras procedencias intactos. Estado CORRECTED_PENDING_INTEGRATION, gate completo se repite tras checkpoint104.

## FAIL-20260908-553 — perfil nuevo fuera de la matriz explícita

Estado: REPRODUCED; clase CONTRACT/COMPOSITION; severidad LOW. Preflight104 detecta53planes y compone52: el nuevo perfil requiere una fila en la matriz explícita de VERIFY_LIBRARY. Añadir su composición1pack/4files y conservar el rechazo de planes no probados; no reducir discovered ni saltar la validación.

FAIL553 correction: canonical explicit composition matrix includes the new profile; exhaustive discovered/composed assertion unchanged. CORRECTED_PENDING_INTEGRATION; checkpoint105 retry.

## FAIL-20260908-554 — consulta asumió contenedor de tests

Estado: RECOVERED; clase TOOL/SCHEMA; severidad LOW. La consulta de lectura trató tests como objeto con items, pero el contrato usa una lista. TypeError sin escrituras de producto ni cambio de test. Recuperación: inspeccionar el tipo real y filtrar directamente la lista. No cambia ningún gate ni evidencia funcional.
V334 FAIL549 recurrence: read-only app/page.tsx guess failed; recovered by rg --files for actual src/app layout. No mutation; same path-inventory lesson applies.

V334 verified closure: FAIL552 and FAIL553 REGRESSION_PROVEN. Preflight105 composes all53profiles, checks162packs/1457files/772Markdown and completes153/153steps. Current provenance1213/137/107/1457 and exhaustive composition assertion both pass. Failed103/104logs preserved; neither gate weakened. FAIL554 recovered by querying the actual tests list; no test status modified.

## FAIL-20260908-555 — receta offline rechazada por pnpm real

Estado: REPRODUCED; clase CONTRACT/CONSUMER; severidad LOW. Primer consumidor aislado enterprise-web devuelve exit1 desde la invocación generada. Se conservan receta y log de V335; sin promoción ni cambio de consumer canónico. Leer diagnóstico, corregir la fuente de routing y agregar regresión antes de continuar.

## FAIL-20260908-556 — sustitución amplia en harness de reintento

Estado: RECOVERED; clase TOOL/HARNESS; severidad LOW. Renombrar consumer- globalmente cambió la referencia consumer-parity.json por un nombre inexistente. Fallo antes de cualquier copia/instalación. Corrección exacta de la referencia; mantener sólo destinos y logs nuevos. No cambia el gate canónico.

## FAIL-20260908-557 — cache de metadata ausente en entorno aislado

Estado: REPRODUCED; clase CONTRACT/OFFLINE; severidad LOW. pnpm acepta los flags corregidos, copia75paquetes desde store y rechaza ERR_PNPM_NO_OFFLINE_META para @emnapi/runtime: el store de contenido no sustituye cache de metadata. El entorno limpio detecta dependencia ambiental no declarada. Sin descarga ni promoción; fijar cache explícita separada, preservar logs y añadir rechazo cuando no existe.
FAIL557 recurrence r3: explicit pnpm-cache path resolves @emnapi but @next/env metadata-full is absent; failed supply-chain policy check. Retain fail-closed behavior and investigate actual baseline configuration/cache before changing any policy.
FAIL557 diagnosis: ignoring the known enterprise workspace dropped its hash-bound minimumReleaseAgeExclude policy and invalidated the upstream verification cache. Preserve the exact known workspace for enterprise; --ignore-workspace remains only on standalone Playwright/Lighthouse. No --trust-lockfile, policy relaxation, or fabricated cache record.

## FAIL-20260908-558 — lectura auxiliar sin destino válido

Estado: RECOVERED; clase TOOL/PATH; severidad LOW. Una consulta auxiliar a ruta inexistente terminó exit1, sin escritura ni pérdida de evidencia. Se elimina esa consulta; los archivos requeridos se identificaron con rg --files. No altera gates ni producto.

V335 focal closure: FAIL555/557 REGRESSION_PROVEN by final rebuilt43tests and3real frozen/offline installs. Unknown CLI options corrected; metadata cache explicit; enterprise workspace policy retained without trust-lockfile. Original failed logs preserved. FAIL556/558 read-only/harness recoveries do not count as product tests.

## FAIL-20260908-559 — glob de búsqueda pasado como ruta

Estado: RECOVERED; clase TOOL/PATH; severidad LOW. rg recibió un glob literal en PowerShell y rechazó os123; corregir usando -g. Sin escritura de producto ni pérdida de evidencia; se conserva lección LIB-FAIL-2266.

## FAIL-20260908-560 — refresh aislado de metadata

Estado: REPRODUCED; clase TOOL/PROBE; severidad LOW. Primer refresh lockfile-only en cache privada devuelve exit1; comandos/log/resultados V336 preservados. Lock/workspace/package intactos y node_modules ausente. Leer diagnóstico antes de reintentar; no promoción ni cambio de política.

## FAIL-20260908-561 — metadata refrescada insuficiente para revalidación offline

Estado: REPRODUCED; clase DEPENDENCY/CACHE; severidad MEDIUM. V336 cold-policy retira resultados históricos en copia privada y obtiene ERR_PNPM_NO_OFFLINE_META para @next/swc-linux-arm64-musl pese al refresh online previo. Tres instalaciones con políticas recién verificadas sí pasaron. No equiparar ese PASS con capacidad de revalidación offline completa. Logs/inputs preservados; investigar formato full/filtered y reparar el protocolo antes de claim.

## FAIL-20260908-562 — nombre de campo inferido en auditor auxiliar

Estado: REPRODUCED; clase TOOL/SCHEMA; severidad LOW. El auditor V336 asumió lockfileHash y obtuvo KeyError al leer registros locales. No escribió qualification ni alteró cache; inspeccionar claves reales, corregir el auditor y repetir sus asserts. Los logs de pnpm permanecen intactos.
V336 FAIL562 seguimiento: corregido campo real lockfile.hash; assertion posterior detecta layout store diferente tras probes adicionales. Preservar observación inicial de index.db y revisar todos los archivos actuales antes de clasificar contenido.

V336 closure: FAIL560 RECOVERED by3empty-cache online checks without unsupported flag. FAIL561 REGRESSION_PROVEN:17official full JSON bodies preserved,3cold offline policy replays and3final fresh offline installs pass. FAIL562 RECOVERED: actual nested lockfile.hash and hidden projects snapshots verified against exact source hashes; cache-qualification auditor PASS. No failed history removed.

## FAIL-20260908-563 — verifier estructural tras integración V336

Estado: REPRODUCED; clase GATE; severidad MEDIUM. VERIFY_LIBRARY devuelve exit1 tras documentación de cache; log íntegro V336/verify-library.log preservado. No checkpoint ni cierre global; leer diagnóstico y corregir owner antes de continuar.

V336 FAIL563 RECOVERED: current evidence hashes recorded in checkpoint109, resume109 PASS, then VERIFY_LIBRARY r2 PASS162/1459/774/53. The gate correctly rejected stale checkpoint order; no validator bypass or weakened assertion. Closure110 preserves109events.

## FAIL-20260908-564 — suite acquisition candidato V337

Estado: REPRODUCED; clase GATE; severidad MEDIUM. test_upstream_acquisition exit1 al ampliar roster de npm-lifecycle; log V337 preservado. No adquisición ni promoción realizadas. Diagnosticar y corregir fuente canónica candidata antes de continuar.

## FAIL-20260908-565 — suite perfiles candidato V337

Estado: REPRODUCED; clase GATE; severidad MEDIUM. test_source_profiles exit1 sobre candidato V337; suite opaque sí pasa. No perfil adquirido todavía; preservar log, corregir fuente y repetir ambos gates.

## FAIL-20260908-566 — nombre inferido de template de perfil

Estado: RECOVERED; clase TOOL/PATH; severidad LOW. Lectura de template con nombre supuesto falló; resolver con rg --files -g. Sin adquisición ni cambios de producto.

## FAIL-20260908-567 — resumen de tests de perfiles desactualizado

Estado: REPRODUCED; clase EVIDENCE/COUNT; severidad LOW. Nueva suite ejecuta18perfiles pero imprime valid=17 heredado. Corregir resumen desde la matriz real, sin alterar las aserciones ni presentar el texto anterior como18.

## FAIL-20260908-568 — serializador auxiliar asume fence único

Estado: REPRODUCED; clase TOOL/PACK; severidad LOW. build-pack V337 detectó29de31bloques por asumir cuatro backticks; el pack no se escribió y materialización dependiente falló por ausencia. Preservado original; inspeccionar delimitadores reales, corregir serializador y ejecutar dependencias sólo tras éxito. No adquisición ni fuente canónica sustituida.

V337 focal closure: FAIL564/565 REGRESSION_PROVEN by actual four-digit1100fixture,78opaque tests source/rebuild,18profiles and real SHA512 acquisition. FAIL567 RECOVERED via dynamic executed profile counter18. FAIL566 lookup recovered; FAIL568 actual cause was blank-line variation, corrected with all31baseline payloads checked and33file byte-identical rebuild. Original failed logs preserved.

## FAIL-20260908-569 — referencias de versión con formato compacto omitidas

Estado: REPRODUCED; clase COMPOSITION/GATE; severidad MEDIUM. Preflight111 falla estructuralmente: AWS_SECURE_EMAIL_ATTACHMENT_PROCESSING_PACK_PLAN conserva0.4.78 en JSON compacto. La selección por espacio literal actualizó29planes pero omitió otros formatos. Preservar log/estado111; enumerar todos los planes por JSON/regex de whitespace, corregir referencias exactas y repetir gate tras nuevo checkpoint.

V337 FAIL569 corrective source update: one compact JSON reference fixed; parsed audit confirms30composition plans use0.4.79. No version-field whitespace assumption. Full Preflight rerun follows checkpoint112.

## FAIL-20260908-570 — segundo gate integrado V337

Estado: REPRODUCED; clase GATE; severidad MEDIUM. Preflight112 vuelve a exit1 en library-structural-and-materialization. Log completo preservado; diagnosticar owner exacto y corregir antes de declarar cierre.

V337 FAIL570 diagnosis/fix: six static Amazon composition counts need+2 for the two explicit new core files; all30upstream-core consumers select star. Compared against53verified V336profile counts. Ten profile introductions were already stale and are corrected with before/verified/after receipt. Assertions remain explicit; full replay after checkpoint113.

V337 closure114: FAIL569 and FAIL570 -> REGRESSION_PROVEN. Preflight113 completed154/154steps PASS; all53actual composition counts equal expected. Compact JSON reference and six count assertions/ten introductions remain corrected in canonical owners. Failed111/112evidence retained. Overall availability BLOCKED only Docker; no runtime admission. See PNPM_LIFECYCLE_LICENSE_V337.md integrated closure.

## FAIL-20260908-571 — root package omitted from observed union

V338 research probe counted only external child manifests and bundle markers, so pnpm11.25.0 was incorrectly marked NOT_FOUND_BY_THIS_PROBE despite root package.json and LICENSE. Status REPRODUCED; class EVIDENCE/INVENTORY; severity LOW. Preserve initial audit and script before correction; include root manifest independently, retain external21 count, and independently reconcile selected identities before use. No runtime or production promotion occurred.

## FAIL-20260908-572 — guessed inventory filename

V338 read-only search referenced nonexistent EXECUTABLE_LIBRARY_INVENTORY.md. RECOVERED by resolving inventory evidence path from validated execution state and using actual file inventory. No artifact mutation resulted from failed lookup. Prior LIB-FAIL-2282 already requires discovery before reads; apply it consistently.

V338 FAIL571 REGRESSION_PROVEN: independent path-splitting+22manifest reconciliation gives468observed identities and includes pnpm root;7marker boundaries PASS,442original file hashes unchanged. Initial erroneous467union preserved. FAIL572 recovered using README.md evidence owner. No full-runtime admission.

## FAIL-20260908-573 — V338 structural gate115 failure

REPRODUCED; class GATE; severity MEDIUM. VERIFY_LIBRARY returned exit1 before structural closure. Complete verify-library115.log retained; diagnose exact gate and correct canonical source or invocation before replay. No PASS or production claim.

## FAIL-20260908-574 — implementation directory count includes README

REPRODUCED; class EVIDENCE/INVENTORY; severity LOW. All-file Markdown baseline parity probe assumed162entries; directory also contains its README index. Preserve entire baseline including index; distinguish implementation pack definitions from directory Markdown count and require both parity and162definitions.

V338 FAIL573 diagnosis: closure114 appended duplicate canonical IDs after Preflight113, and V338 repeated this pattern. Resume hash validation alone did not test ledger uniqueness. Keep original rows and successor content, remove canonical-token formatting only from three successor references, then execute the exact VERIFY_LIBRARY regex. Prior lesson065/531/620 already describes this rule. Regression uniqueness passes before checkpoint116; full structural replay required. FAIL574 corrected:163Markdown directory files are162packs+README, all163hashes unchanged.

V338 closure117: FAIL573 REGRESSION_PROVEN by full VERIFY_LIBRARY116 PASS and exact ledger uniqueness check after closure. FAIL574163directoryMarkdown=162packs+index independently verified with allhashes unchanged. Other V338 notice evidence remains research/draft, not runtime admission. No new canonical failure-ID declaration is appended for a successor.

## FAIL-20260908-575 — Yarn PnP declared gitHead does not resolve

V340 source-provenance probe; EXTERNAL_SOURCE / MEDIUM; REPRODUCED. Exact npm4.1.7 metadata still declares5761b03feb2146da8ce8cefeba6482f3cb6edb79, but official GitHub tree and package-directory queries return404. Both responses retained in V340 research receipts. This does not prove the source never existed. Official annotated release tag resolves to a different commit4fe4d4bf45a13dca90181c5b7fee61376aa21794 with version4.1.7; record both, do not replace npm gitHead or assert byte equivalence. Reopen on registry-to-release provenance or governed archive/build comparison. Continue independent Undici and tagged source licence research with admission blocked.

V340 FAIL575 follow-up: direct git/commits probe also404, so failure is not resolved by substituting tree/commit endpoints. Official tag/source provides separate licence evidence only. Preserve declared and tagged identities; no closure of published-byte binding. Existing ledger lessons1572/1691 informed direct-commit check.

## FAIL-20260908-576 — optional parser path assumed from older TypeScript layout

V341 TOOLING/LOW; REPRODUCED. Existing TypeScript7.0.2 junction has no lib/typescript.js. Probe failed read-only; no source or payload mutated. Inventory actual parser files before selection; never assume older package layout or install a replacement silently. Evidence: %TEMP%/elite-v341-171c0a35821f4cf6aacf2ab09cd9d6c1/parser-probe.txt. Existing filename-discovery lessons consulted.

V341 FAIL576 RECOVERED: actual file inventory located Next16.3.2 compiled Acorn8.14.0; hash/licence/version recorded. Existing parser executed five source comparisons plus eight actual-source negative mutations successfully; no replacement TypeScript install or payload change. Initial missing-file probe remains preserved.

## FAIL-20260908-577 — V341 structural gate rejects local Windows home path

GATE/DOCUMENT_PORTABILITY; LOW; REPRODUCED. VERIFY_LIBRARY122 exit1 rejects the backslash Windows home path introduced in the new FAIL576 evidence reference. Preserve failed verify-library122.log and before-fix122 snapshot in V341 staging; normalize only that new reference to temp-relative portable evidence notation, then recheckpoint and replay the affected structural gate. No code or release-payload claim changes.

V341 FAIL577 REGRESSION_PROVEN: VERIFY_LIBRARY123 PASS after replacing the new home reference with temp-relative notation. Failed122 log and before-state retained. No validator weakening, code change or runtime promotion. FAIL576 remains recovered using inventoried existing parser; Yarn FAIL575 remains unresolved.

## FAIL-20260908-578 — V342 discovery recurrence

TOOLING/LOW; RECOVERED. guessed V301 filename absent; inventory resolved HISTOGRAM_BOUNDARY_REPAIR_V301.md. Read-only failure preserved in %TEMP%/elite-v342-9b34063622f44f5782a8e63857306dc1/discovery-errors.txt. No missing-source or exhaustive-search claim inferred from failed read. Prior filename/glob lessons apply.

## FAIL-20260908-579 — V342 discovery recurrence

TOOLING/LOW; RECOVERED. literal rg wildcard path rejected on Windows; use directory plus native -g filter. Read-only failure preserved in %TEMP%/elite-v342-9b34063622f44f5782a8e63857306dc1/discovery-errors.txt. No missing-source or exhaustive-search claim inferred from failed read. Prior filename/glob lessons apply.

## FAIL-20260908-580 — V342 guessed catalogue path

TOOLING/LOW; RECOVERED. A read-only search assumed a nonexistent implementation catalogue name; no absence inferred. Use exact pack owner and rg file inventory for dependent references. Prior filename lessons apply; original missing path was markdown_system/IMPLEMENTATION_PACK_CATALOG.md.

## FAIL-20260908-581 — observability reports success for incomplete sink write

CODE/DELIVERY; MEDIUM; REPRODUCED. Exact candidate0.3.0 plus boundary regressions returns nil for zero/short/negative/oversized byte counts paired with nil error. Four red cases retained in %TEMP%/elite-v342-9b34063622f44f5782a8e63857306dc1/delivery-red.log. Writer violates the io.Writer contract; this is defensive validation missing from the local explicit-delivery claim, not a Go vulnerability. Canonical owner GO_OBSERVABILITY_CORE; contain as CANDIDATE and add bounded count validation before JSONHandler. No retry, raw-error disclosure or runtime promotion.

## FAIL-20260908-582 — V342 mixed shell/tool options during discovery

TOOLING/LOW; RECOVERED. rg received PowerShell ErrorAction as a native option while searching a guessed archive-script name. Its option parse rejected the call; no result used. File inventory located CREATE_PORTABLE_ARCHIVE.ps1, then its exact source was read. Use tool-native arguments only.

V342 FAIL581 REGRESSION_PROVEN locally: four baseline false-success cases now return ErrDeliveryUnknown with one writer call; six boundary modes and complete-write control pass. Two canonical files rebuilt identically;28tests x3,17seed cases,real HTTP/file reference and3finite fuzz targets pass. Candidate admission and actual runtime operation remain blocked.

## FAIL-20260908-583 — manual ledger count off by one in integration helper

EVIDENCE/LOW; RECOVERED. Helper expected2508before recording this failure; actual2507IDs were all unique (2502baseline plus5new). Preserve integration-count-failure.txt and original helper in V342 staging. Resume partial integration without repeating its one-shot writes; derive count from canonical token set and verify baseline plus newly introduced IDs. No test result or code changed.

## FAIL-20260908-584 — guessed fuzz runner extension

TOOLING/LOW; RECOVERED. V343 read-only Get-Content and rg targeted unobserved run_fuzz_gate.py, which does not exist. The four-file directory inventory identifies run_go_fuzz_gate.ps1. No test or gate result was inferred from the failed lookup. Recurrent filename-discovery error: inspect the exact manifest/directory before calling any entry point; use the observed PowerShell runner. Evidence: V343 discovery-error.txt.

## FAIL-20260908-585 — update helper assumed uniform code fence

TOOLING/LOW; RECOVERING. V343 update_pack.py matched the Go/PROVENANCE blocks but failed to locate README because its regex assumed exactly four backticks. Canonical pack/profile writes had not started; candidate documentation edits are preserved. Inspect actual Markdown fence delimiters, preserve failed helper and recover using matching opening/closing fence lengths. Do not rerun one-shot documentation insertion blindly. No product test failed or passed because of this helper.
V343 FAIL585 diagnosis corrected after exact source read: the delimiters are four backticks in both cases; the cause is required double-newline separators in the helper, while README/PROVENANCE use one newline. Recovery accepts one or more newlines, preserves matching fences and detects the already-inserted exact section. Initial fence-length hypothesis is rejected.

V343 FAIL585 RECOVERED: corrected spacing matcher, fresh67/746composition and four-block hash parity; no duplicate documentation insertion. Computer interruption is recorded separately as incomplete final gates, not a fabricated code failure or PASS.

## FAIL-20260908-586 — export-data reads fail after machine interruption

ENVIRONMENT/BUILD; OPEN. Recovery go vet ./... and the rebuilt WhatsApp test invocation report cannot find export data (bufio: buffer full) for existing packages including returneffects/aifoundation, before target tests complete. Preserve rebuilt-vet-recovered.log and rebuilt-whatsapp.jsonl; concurrent fuzz baseline is being inspected. No product source change or test weakening. Hypothesis: interrupted shared Go build cache contains unusable artifacts. Verify the pinned toolchain and retry the exact canonical tree with a new stage-owned GOCACHE; preserve the old cache and failed logs. Claim recovery only from complete receipts.

V343 FAIL586 RECOVERED: same pinned Go executable and canonical tree pass full tests/vet/build under a fresh stage-owned GOCACHE; old cache and failed logs preserved. Final rebuilt14host tests x3,finite fuzz and full baseline/vet/build PASS after adding retained post-commit report-loss regression. No product fix or assertion weakening was needed for the environment failure.

## FAIL-20260908-587 — ledger row status not reconciled after proven recovery

EVIDENCE/LOW; RECOVERING. VERIFY_LIBRARY128 rejected one OPEN local row; Preflight128 correctly stopped at its structural dependency. FAIL586 recovery was documented as a successor but its canonical row remained OPEN. Preserve both failed logs and the before ledger. Update only the actual row statuses after linking demonstrated recovery, retain the original failure narrative and successors, checkpoint the updated hashes, then rerun both gates. No gate or regex is weakened.

## FAIL-20260908-588 — redirected PostgreSQL launcher left PowerShell waiting

ENVIRONMENT/HARNESS; RECOVERING. Controlled restart stopped and started the owned PostgreSQL successfully (ready on loopback55959); pg_ctl/psql had no live process, but the wrapper PowerShell remained waiting and had not created restart-after.jsonl. Its exact unique command marker identified task-owned PID14152. Preserve logs; stop only that wrapper, keep the owned database running, capture the after snapshot independently, then stop PostgreSQL through a process invocation that does not inherit a PowerShell output pipe. No product data or source changes. Do not infer snapshot equality until measured.

V343 FAIL587 RECOVERED: corrected canonical status passes embedded VERIFY_LIBRARY129 and154Preflight steps; no verifier change. FAIL588 RECOVERED: exact owned waiting wrapper stopped, independent after snapshot matches before SHA for412synthetic rows across4tables, PostgreSQL stopped via hidden direct process redirection. Scope is controlled local restart; not production DR or crash-proof power-loss certification.

## FAIL-20260908-589 — refund host discards log delivery failures

CODE/OPERABILITY; REPRODUCED. V344 exact0.1.3 main.go plus a behavioral probe observes no report error for six sink modes: explicit error, full-length/error, zero, short, negative and oversized count. The old void logStepError calls log.Print and run continues without observing write failure. Reflection permits the same probe across old void/new error-returning internal API; it does not alter baseline source. delivery-red.log and baseline-pack.md are preserved. Contain by retaining CONDITIONED; return a static report error, validate one-write counts and make the actual run loop stop on reporting failure. No refund policy, provider, retry, amount, currency or durable-state change is authorized by this fix.

## FAIL-20260908-590 — manual test tally included the fuzz target

EVIDENCE/LOW; RECOVERING. integrate131.py expected18ordinary tests but the JSON log contains17ordinary tests plus one fuzz target. All17passed three times (51executions); the integration helper failed before writing report/owners/state. Preserve failed helper; derive totals directly from Test/Fuzz JSON events, keep seed/skip counts separate and correct draft prose before integrating. No code or test failed and no assertion is weakened.
V344 FAIL590 RECOVERED: observed JSON records prove17ordinary tests x3=51executions, one separate fuzz target and14distinct seed cases. Draft/helper totals corrected before integration; no source/test change. observed-test-counts.json preserves per-test counts.

V344 FAIL589 REGRESSION_PROVEN: all six exact-baseline failure modes now return static report error; actual run loop stops before another claim/token and shared main exit adapter exits1. Full module17tests x3,real closed file,subprocess/binary,14seeds,finite fuzz and15/15rebuild PASS. No provider/business/DB/lock changes; target supervision and stalled writers remain separate.

V344 closure132:589REGRESSION_PROVEN and590RECOVERED;17ordinary tests/51executions and14seeds kept separate. All154Preflight checks/structural131 PASS,three-file consumer delta and743unchanged files. No source/test weakening or new critical local failure.

## FAIL-20260908-591 — native rg glob passed as a Windows path

TOOLING/LOW; RECOVERED. V345 two read-only searches passed wildcard paths directly to rg, producing OS error123. Successful results from other explicit paths remain usable; missing glob results are not absence evidence. Recurrence of filename-discovery misuse: use rg --files and rg -g filters on observed directories. No canonical code or gate changed. Exact rejected patterns: markdown_system/CAPABILITY_GAP* and implementation_packs/OFFICIAL_TOOL*.

## FAIL-20260908-592 — operational tool runner bounds output only after collection and exposes timeout argv

CODE/OPERABILITY-PRIVACY; REPRODUCED. V345 exact1.1.2 runner plus two behavioral probes: a child exceeded the patched1024-byte cap yet continued to a later synthetic marker; timeout traceback disclosed PRIVATE_MARKER_345 supplied as argv. subprocess.run(capture_output=True) accumulates before the length check and TimeoutExpired is chained. Preserve baseline-pack.md, baseline/ and red.log. Correct the existing owner to bounded nonblocking reads, finite child wait/cleanup and static errors with suppressed private causes. No generic service-supervisor admission, process-tree kill or production effect is implied.

## FAIL-20260908-593 — composition delta included generated receipt

EVIDENCE/LOW; RECOVERING. V345 parity helper proved37canonical files equal, then classified MATERIALIZATION_RECORD.md as product because it guessed metadata exclusions by suffix. The comparison correctly found that receipt differs between rebuilds; assertion stopped before refund-reference execution or its receipt. Preserve original helper/output. Inspect actual materialization record and use its explicit output inventory for the746product paths, retaining generated receipts separately. No product test failed; no source change or assertion weakening.

V345 FAIL593 RECOVERED: actual MATERIALIZATION_RECORD.md file table supplied the746unique product paths and exact output hashes; previous/current path sets match, three product files changed and743remain identical. Generated receipt differences remain preserved separately. All37owner files equal candidate/rebuild/consumer; unchanged refund0.1.4 binary executed by the new runner returns1 with fixed REFUND_HOST_FAILED and no provider/database configuration. No source or test changes for this evidence recovery.

V345 FAIL592 REGRESSION_PROVEN: two exact-baseline red cases now pass;14runner tests,109other-module/runner aggregate PASS plus1existing symlink skip,12readiness PASS,37file parity and actual refund subprocess. Finite child correction does not admit service supervision.

## FAIL-20260908-594 — native qualification exposes create-before-assign crash window

CODE/RELIABILITY; REPRODUCED. Isolated V345 Job Object prototype creates a child suspended, then assigns it to the job. A controlled owner os._exit between those steps closes the job yet leaves the suspended child alive. No user code ran; the exact owned process handle was terminated/reaped by the controller. launch-gap-red.json preserves evidence. This is NOT admitted product code. Consulted official UpdateProcThreadAttribute/STARTUPINFOEXW: qualify job membership at CreateProcess via JOB_LIST, then repeat abrupt pre-resume owner death and all prior tree/breakaway cases. Do not promote the initial prototype on its prior twelve PASS experiments.

V345 FAIL594 REGRESSION_PROVEN in the isolated candidate: JOB_LIST supplied within CreateProcess removes the demonstrated post-create assignment window. Five scenarios x3 pass, then repeat x3 after exact canonical reconstruction; includes pre-resume owner death and no child user effect. Initial prototype remains rejected, no product supervisor admission.

## FAIL-20260908-595 — native I/O candidate reports cleanup failure on terminated trees

CODE/RELIABILITY; REGRESSION_PROVEN. Initial observation was REPRODUCED. V346 isolated job+pipe capture passes successful capture paths but five negative/termination cases return CLEANUP_FAILED instead of their expected reason. Preserve candidate-tests.log/json and initial candidate source. No product owner/profile has changed. Diagnose the cleanup/accounting/handle lifecycle with direct OS observations; keep tree_empty false until verified, do not weaken assertions or promote the candidate.

## FAIL-20260908-596 — peak job counter assumed bounded by enforced memory limit

ASSUMPTION/RESOURCE; RECOVERED_WITH_LIMITATION. Initial observation was REPRODUCED. V346 real64MiB job limit refused a128MiB allocation (MemoryError branch and exit0), but PeakJobMemoryUsed read141565952bytes, exceeding the configured limit. The test incorrectly equated that peak counter with a verified committed-memory ceiling. Preserve resource-tests.log/json and source. Inspect documented/actual allocation semantics and verify denied native commit plus permitted small commit; retain the counter honestly, do not claim bounded RSS or hide the observed value.

## FAIL-20260908-597 — owner-death I/O fixture did not complete its ready handshake

HARNESS/RELIABILITY; RECOVERED. Initial observation was REPRODUCED. V346 resource suite includes an error in the redirected-pipe owner-death fixture, alongside the separate peak-counter assertion. Preserve complete failed log; inspect the exact error and add task-owned stderr capture before diagnosing. No product adoption, no inferred tree cleanup PASS and no unrelated process mutation.

V346 closure of current defects:595REGRESSION_PROVEN by shared-budget wait and real/injected transitions;596RECOVERED_WITH_LIMITATION via native denied/granted allocation controls, observed peak retained and RSS/peak-ceiling claim rejected;597RECOVERED by hidden descendant readiness before shared-pipe saturation and actual parent-death handle checks. Two current files reconstructed exactly,23tests and15base lifecycle observations pass. No product supervisor admission.

## FAIL-20260908-598 — CRT permission error lacks Win32 error field in identity-lock test

HARNESS/EVIDENCE; RECOVERED. Initially REPRODUCED. V347 governed-launch suite: 17 tests pass and the lock-denial test fails because Python open(r+b) raises a permission OSError with winerror=None, contrary to its asserted [5,32]. The operation was refused, but native denial identity is not yet demonstrated by that assertion. Preserve launch-tests.log/json and pre-fix test source. Add a direct CreateFileW write probe and inspect GetLastError; distinguish CRT errno from native winerror without accepting a successful mutation. Product owner/profiles unchanged.


V347 FAIL598 RECOVERED: direct CreateFileW write request is refused with GetLastError32; CRT errno EACCES and native winerror remain distinct. All write/delete/file-and-directory-rename attempts must still fail until pin closure; close permits rename. Rebuilt 20 launch tests and 23 prior capture/resource tests PASS, zero skips. Original assertion/log/source preserved; no product admission.

## FAIL-20260908-599 — wait-only event handle is not a hostile-worker security boundary

ASSUMPTION/SECURITY; RECOVERED_WITH_LIMITATION. Initially REPRODUCED. V348 direct SetEvent on the inherited SYNCHRONIZE-only handle is denied5, but a same-user child can DuplicateHandle requesting EVENT_MODIFY_STATE and signal that duplicate. event-access-probe.json preserves direct0/error5, duplicate1 and signal1. The earlier wording that the child can only observe is too broad. Official DuplicateHandle documents that some object types permit greater access. Keep the protocol limited to trusted cooperative workers, preserve the exact direct-handle right test and characterize the amplification; do not present it as an untrusted-worker sandbox or acknowledgement of drained business state. Rename the result to EXITED_DURING_GRACE, which describes observed process lifecycle only. Product packs unchanged.


## FAIL-20260908-600 — descendant capture passes but immediate fixture directory removal is busy

HARNESS/PROCESS-LIFECYCLE; RECOVERED. Initially REPRODUCED. V348 retained 23-case capture suite passed the silent-descendant timeout/tree_empty assertions, then TemporaryDirectory teardown failed with WinError32 removing tmpltrefwfb (intermediate WinError5). Other22 cases pass. Preserve capture-regression.log/json and original test source. Job accounting empty and root signaled do not establish immediate descendant filesystem teardown. Strengthen this fixture with an exact descendant process handle acquired while owned and alive, require its termination signal, and use a bounded explicit owned-directory release check; never turn a live descendant or deadline expiry into PASS. No arbitrary PID termination or unbounded retry.


V348 closure of current findings:599 RECOVERED_WITH_LIMITATION, direct wait-only handle rights do not prevent same-user amplification; protocol requires trusted workers and reports only EXITED_DURING_GRACE.600 RECOVERED by acquiring/verifying the live owned descendant handle, requiring its signal and finite directory-release check. Rebuilt19+20+23tests PASS/0skips. Neither correction certifies hostile-worker isolation, domain drain or immediate filesystem quiescence for every child.

## FAIL-20260909-601 — recovery search used an unexpanded native wildcard and wrong contract extension

TOOLING/LOW; RECOVERED. V350 read-only rg targeted PROJECT_ENGINEERING_CONTRACT.md (actual owner is JSON) and markdown_system/*ROADMAP* as a native path. rg reported missing path/invalid syntax; no mutation or product failure. Recovered by reading the known root REUSABLE_CODE_READINESS_ROADMAP.md and PROJECT_EXECUTION_STATE.json; subsequent searches must use existing roots with -g filters. Original tool output is retained in the task history; no successful owner read is inferred from the failed search.

V350 FAIL601 recurrence, RECOVERED: a README read guessed go_native_fuzz_gate although the returned inventory named go_fuzz_gate. Materialization and gofmt succeeded; only the read failed. Corrected using the observed path. A following multi-file apply_patch used an incomplete ledger line and was rejected atomically; an explicit read confirmed neither intended edit was applied. The complete exact line is used for recovery. Retain tool errors and never infer success from preceding commands.

## FAIL-20260909-602 — Go fixture flags rejected before ready handshake

HARNESS/INTEGRATION; RECOVERED. V350 first lifecycle subprocess cases fail their ready handshake. The fixture dispatch used --fixture without terminating Go test flag parsing. Preserve worker-lifecycle.log/json, original Python source and process receipts before correcting the explicit argument boundary. Native Go tests, prior host/processor tests, vet and fixture build already passed; they do not substitute for the failed end-to-end launch. No product source changed, no real provider effects. Adding the explicit -- boundary yields ten integrated tests PASS; exact canonical rebuilt confirmation remains required.

## FAIL-20260909-604 — open-pipe descendant fixture repeats immediate directory-release assumption

HARNESS/LIFECYCLE; RECOVERED. V350 unchanged retained capture suite passes its timeout/tree assertions, then test_timeout_cleans_child_and_descendant_with_open_pipe teardown fails WinError32/5 for owned tmp8xbja8xi.22other tests PASS. This is the same resource-release boundary learned in FAIL600, now in the open-pipe sibling fixture. Preserve rebuilt-capture.log/json and source. Apply the existing exact-owned-descendant handle/termination proof plus finite directory-release observation to this fixture; do not hide teardown, delete unrelated paths, weaken the timeout assertion or claim universal quiescence from job accounting.

## FAIL-20260909-603 — already-signaled event initially returns an uncancelled worker context

CODE/LIFECYCLE; REGRESSION_PROVEN. Review-derived regression TestNativeStopPreSignaledPreventsFirstRealHostClaim fails against the initial V350 adapter: cancellation is only performed by a newly scheduled goroutine, so nativeStopContext returns a live context despite an already-signaled event. Preserve presignaled-red.log and initial adapter. Add synchronous zero-time initial wait before returning the context, retain joined cleanup and finite subsequent waits; prove the actual host does not make its first claim. This is isolated qualification code, not a product regression or permission to alter refund policies.

V350 closure of local findings:601 recovered from observed paths/exact patch context;602 recovered with explicit Go argument boundary and TestMain fixture dispatch;603 regression-proven by synchronous event read and actual no-first-claim test;604 recovered by extending exact owned-descendant termination and finite release observation to the open-pipe fixture. Fresh canonical10native tests x3 +10integrated +90supervisor tests PASS. Original failed sources/logs/receipts remain. Product owner files unchanged; PostgreSQL/provider/target admission remains open.

## FAIL-20260909-605 — structural gate rejects stale current Markdown inventory

EVIDENCE/INTEGRATION; RECOVERED. Preflight143 correctly stops at VERIFY_LIBRARY because FRANCHISE_PREFLIGHT_GAP still advertises787Markdown after the new V350 report raises the inventory to788; expected162/1461/788/53. Integration143 appended the new evidence but omitted the authoritative header update. Preserve preflight143.log (no successful Preflight receipt exists). Correct that current header, record a successor checkpoint with changed evidence hashes and rerun the full gate. No test or pack claim is waived; original dated inventories remain historical.

V350 FAIL605 RECOVERED: canonical inventory corrected788; successor checkpoint144 retained the failed143 log, and full Preflight144 passed all executed gates. No code gate bypass or history rewrite.

## FAIL-20260909-606 — next-action cursor exceeds validator size limit

EVIDENCE/CHECKPOINT; RECOVERED. Full Preflight144 passed154checks. close145 prepared state145, but checkpoint creation correctly rejected next_action longer than800characters after a PostgreSQL fixture caution was appended. No EVT-0145 was written. Preserve prepared state and passed gap receipt145; reduce the cursor to actionable routing, put the detailed caution in the canonical V350 report, refresh hashes and validate before retry. No source/test or admission claim changes; do not bypass the compact-cursor contract.

V351 FAIL607 further recurrence: an unnecessary final read guessed0001_enterprise_platform.up.sql and failed after the meaningful inspected source succeeded. Keep the failed command visible. Migration filenames must come from the observed migration inventory, not from inferred naming conventions.

## FAIL-20260909-611 — advisory renderer rejects an absolute output reference

TOOLING/READINESS; RECOVERED. V351 requested roundD renderer output via an absolute external filename; its canonical-relative-path guard rejected before writing. Preserve the tool rejection and existing root advisory. Re-render with the observed CLI contract in an isolated copy of the readiness record, compare with the historical prompt and do not mark the pending user answer as supplied.

## FAIL-20260909-610 — full database snapshot sent through a bounded process-metadata encoder

HARNESS/EVIDENCE; RECOVERED. V351 third run passes11cases, completes controlled PostgreSQL restart and compares before/after snapshots equal, then fails constructing the digest receipt: result_store.encode correctly rejects the combined database snapshots above its16KiB metadata limit. Preserve postgres-third.log/json and source. Use a separate canonical JSON snapshot encoding for evidence digests; never enlarge or bypass the operational result-store limit. Full canonical rerun must still prove restart and cleanup.

## FAIL-20260909-609 — controlled PostgreSQL shutdown fails before restart proof

HARNESS/OPERATION; RECOVERED. V351 all11native PostgreSQL cases passed, but the first pg_ctl fast shutdown exceeded its10second wait. Preserve postgres-second.json/log and owned server/stop logs. Initial wording implied cleanup retry also failed; detailed evidence corrects that: server checkpoint took14.885seconds flushing12841files, cleanup completed and owned_postgres_stopped=true. No restart PASS was produced. Increase only this test fixture's explicit budget to40seconds (command45), preserve timeout failure, and use the configured postgres role in pg_isready. Do not mutate other instances, weaken cleanup or claim a production shutdown SLO.

## FAIL-20260909-608 — fresh owned PostgreSQL exits during startup

ENVIRONMENT/INTEGRATION; RECOVERED. V351 first owned PostgreSQL18.6 fixture initialized its fresh cluster but the postmaster exited before readiness. No seed or worker ran. postgres-first.log/json and owned postgres-0.log preserve the error; inspect the server output before changing initialization/launch. No reuse of an existing cluster, weakened readiness check or inferred database integration PASS.

## FAIL-20260909-607 — readiness-report recovery path guessed instead of resolving its owner

TOOLING/LOW; RECOVERED. V351 read-only rg used nonexistent PROJECT_READINESS_GATE_REPORT.json. The file inventory establishes PROJECT_READINESS_REPORT.json. Resume145 was validated before modifying this ledger; subsequent reads use observed paths. This repeats the path-discovery lesson601; no product source or evidence status was inferred from the failed read. A second dependent read was incorrectly launched in parallel with inventory discovery and guessed0040_return_refund.test.sql; actual returned path is0020_return_refund_execution.test.sql. The missing-file error is retained; discovery-dependent reads must be sequential.

## FAIL-20260909-612 — parity helper assumed a different materialization table

HARNESS/EVIDENCE; RECOVERED. V351 preserve351.py wrote the canonical qualification report and reconstructed the consumer, then rejected its own guessed MATERIALIZATION_RECORD table regex before producing parity evidence. Preserve the helper/output and existing reconstruction. Inspect the actual table and compare all recorded product paths; do not repeat one-shot materialization or infer parity from composition success.

V351 closure of local findings608–612: exact rebuilt11 PostgreSQL cases, controlled restart with11identical database snapshots and owned server cleanup PASS; initdb timezone retained, finite40s shutdown budget and independent evidence JSON encoding. RoundD rendered via canonical relative path is byte-identical; no user answer inferred. Actual recorded table proves746product files equal. Original failed logs/sources/receipts remain, with no widened operational metadata or production claim.

## FAIL-20260909-613 — installation acceptance deferred behind final release publication

PROCESS/ACCEPTANCE; RECOVERED. V352 review finds TEST06 explicitly accepts an exact candidate, while guide0.1.1 and prior V327/V332 narratives require a final release before TEST06/08. Combined with final release waiting for material gates, that interpretation postpones acceptance indefinitely. Preserve original sources and checkpoint147; qualify a complete frozen candidate locally before publication. TEST06/08 results must identify its archive SHA and consumed tooling projection; TEST07/final publication requires the same payload or renewed acceptance after a material change. No signer/SCA/notice/provider/production gate is removed. Remaining test statuses stay pending until actual candidate evidence passes.

## FAIL-20260909-614 — intermediate read requested an output published only at command completion

TOOLING/LOW; RECOVERED. V352 progress probe requested extracted-structural.stdout while its subprocess was still running; the harness publishes stdout/stderr only after completion. The read returned exit1 despite silent missing-file handling. No verifier failure or PASS is inferred. Use the always-present accept-candidate.log and process completion before reading per-command files; retain the tool output. Candidate and tests remain unchanged.

## FAIL-20260909-615 — synthetic acceptance fixture omitted its declared blocker

HARNESS/CONTRACT; RECOVERED. V352 exact archive verification and extracted full structural validator PASS, then fixture creation fails: status BLOCKED requires a blocker, failure or blocked task. The AUTHORED harness changed ACTIVE to BLOCKED without declaring the synthetic intake blocker. Preserve accept-candidate.log/source and partial fixture; add the actual fixture-only blocker, reconstruct in a fresh destination and rerun acceptance. Do not weaken the execution validator or invent production readiness.

## FAIL-20260909-616 — locked replacement reports Win32 access denied rather than sharing violation

HARNESS/OS-BOUNDARY; RECOVERED. V352 canonical candidate verification and actual guide13commands/8help checks PASS. Real installer reaches File.Move at line93 with the test-owned Skill held FileShare.Read; replacement is rejected with native HRESULT0x80070005, while the fixture required only0x80070020. Preserve canonical-acceptance.log and source. Accept the observed access-denied variant only together with exact rollback, preserved events and successful identical installation after releasing the lock; no error alone proves recovery. Candidate bytes and installer stay unchanged.

V352 recovery613/615/616: complete candidate148 independently extracted and verified; actual guide13commands/8help, locked replacement NEW/EXISTING with native access-denied5, exact rollback, same operation succeeding after release, unchanged events and duplicate-append refusal, then55entry checks PASS. Synthetic BLOCKED fixture now declares incomplete intake. Acceptance ordering correction is demonstrated on exact archive, not an unconditional release grant. Original failed sources/logs remain.

## FAIL-20260909-617 — acceptance-delta helper counted directory documents as packs

HARNESS/INVENTORY; RECOVERED. V352 delta audit already confirms48unchanged test IDs/actions/oracles, unchanged requirements and exactly TEST06/08 planned→passed. Its final assertion equated every implementation_packs/ document with one of162materializable packs and failed. Preserve helper/task output; inspect the directory inventory, distinguish index documents from pack metadata, verify every file hash and count actual pack_id owners. No acceptance criterion or product change is inferred from this helper error. Full Preflight149 already passed independently.

V352 FAIL617 recovered: directory contains162pack owners plus README.md; all163files remain byte-identical.48test IDs/actions/oracles and requirements unchanged; only TEST06/08 statuses advance after exact candidate evidence. This metadata-only verification follows passed Preflight149 and does not change executable sources.

## FAIL-20260909-618 — FX numeric validation accepts non-finite inputs and unrepresentable products

CODE/CORRECTNESS; REGRESSION_PROVEN. V353 Go1.26.7 baseline5tests PASS, but three adversarial tests reproduce6failing subcases against exact GO-FX-CORE0.1.0: NaN/+Inf rates replace a valid stored rate, NaN/+Inf amounts return success, finite-input multiplication overflows to+Inf or underflows tozero with nil error. Negative infinity controls are already rejected. Preserve fx-baseline.log/fx-red.log and original two files. Reject non-finite/nonpositive input before mutation and unrepresentable positive result; preserve prior valid rate and ordinary binary64 conversion behavior. This is an AUTHORED in-memory core outside the67pack profile, not a production exchange-rate or ledger admission. Retain precision/rounding/provider/persistence gaps in TEST02.

## FAIL-20260909-619 — stale native exit-code check stopped the integration script early

TOOLING/CONTROL-FLOW; RECOVERED. V353 update_pack_from_tree.ps1 successfully updated the two canonical blocks, then the surrounding script incorrectly checked LASTEXITCODE after an in-process PowerShell script. With no native exit code set, the guard exited before metadata update and reconstruction. Preserve the tool output; inspect actual stage/version and run the remaining dependent steps explicitly. PowerShell scripts use terminating errors, native subprocesses use LASTEXITCODE; no later step is inferred from exit0 of the wrapper.

V353 FAIL618 regression proven: canonical GO-FX-CORE0.1.1 reconstructs2/2files; all9ordinary tests x3,14fuzz seeds, vet/build and10s/13,047,867fuzz executions PASS. Exact rational oracle checks binary64 rounding independently of production helper. Invalid rates preserve prior rate; NaN/infinities and overflow/zero-underflow return typed errors. Original baseline5PASS and6failing regression subcases remain in stage and report. Financial policy/provider/persistence and21core admission remain blocked. FAIL619 recovered by explicit sequential native subprocess execution; no wrapper exit is treated as proof of omitted work.

## FAIL-20260909-620 — qualification attempted two pack materializations into one destination

HARNESS/COMPOSITION; RECOVERED. V354 first giftcards pack materialized, then loyalty materializer correctly refused nonempty destination. No source or existing destination was overwritten. Preserve original script/log and partial red stage. Materialize each pack into a distinct absent staging directory, then copy the four exact source files into a separate qualification-only module; verify source hashes. This is a test mount, not an admitted product composition.

## FAIL-20260909-621 — credit core identity/boundary audit

CODE/CORRECTNESS; REGRESSION_PROVEN. V354 Loyalty0.1.0 admits int64 credit overflow, concatenated account-key alias across distinct tenant/customer tuples and globally shared entry IDs. Static source inspection violates existing nonnegative/tenant-scoped claims. Reproduce original baseline and focused boundary/identity tests before canonical fix; keep in-memory/financial-policy and integration gates conditioned.

## FAIL-20260909-622 — credit core identity/boundary audit

CODE/CORRECTNESS; REGRESSION_PROVEN. V354 Giftcards0.1.0 globally shares redeem IDs across tenants and accepts blank IDs; String reads map length without the mutex used by writers. Issue rejects duplicate codes, contrary to top-up/idempotent-issue wording. Reproduce tenant/blank-ID failures; correct canonical validation and identity scope, synchronize String, narrow documentation; do not invent durable checkout or financial rules.

## FAIL-20260909-623 — fuzz coordinator deadline at end of giftcards budget

ENVIRONMENT/QUALIFICATION; RECOVERED_WITH_LIMITATION. V354 canonical19ordinary tests x3/vet/build pass. First fuzz target loyalty passes5,529,823executions; giftcards ends after7,480,937executions with context deadline exceeded and no saved crasher/testdata. Whole two-target gate failed and emitted no success receipt. Preserve fuzz.log; do not call these failed-run executions a PASS. Pinned Go source shows cancellation paths and worker parallelism; cause not proven. Hypothesis: coordinator shutdown under24worker scheduling. Qualify same unchanged sources/oracles/budget with explicit GOMAXPROCS=4, GOFLAGS=-parallel=4, at most two runs with distinct receipts; retain original limitation, no compiler/runtime patch.

V354 recovery621/622: original11tests PASS then5focused tests FAIL. Corrected canonical4/4sources reconstruct exactly;19ordinary tests x3, vet/build and two complete2target fuzz gates at4workers/10s per target PASS. Tenant-scoped operation identity remains unique across customer/card within each tenant; tuple keys avoid delimiter alias; int64 overflow/blank ID rejects before mutation. String locking is static repair exercised concurrently, not race-detector proof. FAIL623: qualified only with explicit4workers; original24worker deadline remains unexplained, not hidden or claimed fixed upstream. No persisted crash input found. Original sources/logs remain.

## FAIL-20260909-624 — reminder scheduler has no tenant-scoped read/ack API

CODE/ISOLATION; REGRESSION_PROVEN. V355 static audit: Schedule keys globally by ID, Due() returns pending records for every tenant and MarkSent(id) cannot identify the calling tenant. This contradicts the canonical tenant-scoped claim. String also reads its map without the writer mutex; static race-risk finding only. Reproduce baseline and cross-tenant probes, then introduce explicit tenant parameters with a breaking0.2.0 API and tuple identity. Preserve old source/tests, document migration and reject unscoped callers at compile time. Due is not a reservation, MarkSent is a local flag, and neither demonstrates exactly-once delivery. No provider effects or new business rules authorized.

V355 FAIL624:5baseline tests PASS,3isolation tests FAIL on0.1.0; canonical0.2.0 requires Due(tenant), MarkSent(tenant,id) and keys by tuple.11tests x3,vet/build,2/2rebuild and10s/4worker fuzz 2112907executions PASS. Five old assertion bodies preserved with explicit tenant arguments only; legacy unscoped callers fail compilation as required. String synchronized, exercised concurrently but no -race claim. Waitlist0.1.1 changes metadata only, removes unsafe legacy compatibility; no integration admitted. Due can repeat before local MarkSent, so exactly-once send promise withdrawn.

## FAIL-20260909-625 — promotions int64 arithmetic boundary

CODE/CORRECTNESS; REGRESSION_PROVEN. V356 static source: percent multiplication can overflow int64 before dividing, returning negative or incorrect discount while consuming a use. Preserve0.1.0 sources, baseline and red probes; use exact integer oracle. Repair arithmetic without new tax/discount/rounding policy, preserve typed error precedence and no-mutation rejections. AUTHORED in-memory core outside franchise; no financial admission.

## FAIL-20260909-626 — payroll int64 arithmetic boundary

CODE/CORRECTNESS; REGRESSION_PROVEN. V356 static source: percent multiplication can overflow int64; repeated subtraction can wrap a mathematically negative net back positive and store success. Preserve0.1.0 sources, baseline and red probes; use exact integer oracle. Repair arithmetic without new tax/discount/rounding policy, preserve typed error precedence and no-mutation rejections. AUTHORED in-memory core outside franchise; no financial admission.

## FAIL-20260909-627 — wrong historical payroll revision in new note

DOCUMENTATION/REFERENCE; RECOVERED. V356 new metadata note typed V209 without confirming the payroll evidence reference; the existing canonical footer identifies V216. Preserved both notes before correction and qualification script. Corrected new note to V195/V216; old evidence untouched. No materialization blocks, tests or claims of newly executed historical compositions changed. Derive references from observed canonical footer rather than memory.

V356 FAIL625/626 regression proven:11baseline tests PASS;3new groups fail (10percentage subcases and wrapped-negative-net case). Canonical4/4files,18tests x3,vet/build and37seeds/4614438fuzz executions PASS. Quotient/remainder avoids percentage overflow while preserving truncation; payroll tracks exceeded gross without subtracting through int64 floor and still validates later rules to preserve ErrInvalidRun precedence. Invalid runs remain unstored and reusable; coupon error does not consume a use. FAIL627 reference corrected V216, preserved erroneous note; no historical suite falsely rerun.

## FAIL-20260909-628 — referral index aliases tenant/referee tuples

CODE/ISOLATION; REGRESSION_PROVEN. V357 five original tests pass; three public-API regressions fail on exact0.1.0. Concatenated tenant+NUL+referee aliases distinct tuples. It rejects a distinct creation, reports an absent referee present, and can Convert/Reward another local referral when a shared code exists in the requested tenant. No auth bypass in an installed product demonstrated. Preserve original/red sources and logs; replace ambiguous index keys with typed tuples, preserve API, duplicate precedence and local lifecycle. No reward issuance or business eligibility invented.

V357 FAIL628: claves struct tenant/id eliminan alias sin rechazar nuevos caracteres ni normalizar identidades.5tests originales intactos;3regresiones fallaban antes,10tests x3 y modelo lineal independiente/12semillas/1709469fuzz PASS tras rebuild2/2. Concurrencia64goroutines por Create/Convert/Reward demuestra una transición exitosa por fase; no detector -race ni entrega/antifraude productivos.

## FAIL-20260909-629 — reviews identity boundary

CODE/ISOLATION; REGRESSION_PROVEN. V358 Three-field concatenated identity aliases tenant/subject/author. Distinct Submit rejects as duplicate and foreign Approve/Reject can mutate another record. String reads map length without writer mutex; static race-risk finding, no race detector claim. Baseline11tests PASS;4public-API regressions FAIL across both cores. Preserve original/red source and logs; correct canonical tuple keys, preserve API and lifecycle, qualify independent model and rebuild. No auth, real publication, booking or provider effects inferred.

## FAIL-20260909-630 — waitlist identity boundary

CODE/ISOLATION; REGRESSION_PROVEN. V358 Join forbids NUL subjects, but read/transition APIs do not enforce that constraint. Crafted tenant/subject lookup aliases a valid record in another tenant; State leaks it and Notify/Seat/Cancel mutate it. Position may remove the original pending entry as a consequence. Baseline11tests PASS;4public-API regressions FAIL across both cores. Preserve original/red source and logs; correct canonical tuple keys, preserve API and lifecycle, qualify independent model and rebuild. No auth, real publication, booking or provider effects inferred.

## FAIL-20260909-631 — vet detects ignored String result in new probes

TEST-HARNESS/STATIC; RECOVERED. V358 canonical rebuild4/4 and19tests x3 pass, then go vet rejects five ignored String return values in new concurrency/fuzz probes. The intended effect is to exercise synchronized debug reads; preserve rejected test source, fix explicit discard, update canonical test blocks, reconstruct again and rerun gates. No production contract or vet check disabled.

V358 FAIL629/630:11baseline tests intactos,4regresiones de API fallaban antes;19tests x3,4fuentes reconstruidas y modelo lineal de moderación/FIFO con24semillas/4216603ejecuciones fuzz PASS. Tuplas aíslan todas las lecturas/transiciones; reviews String toma el mutex. FAIL631 recuperado con descarte explícito de5resultados String intencionales; tests rechazados y vet anterior preservados, nueva reconstrucción/vet PASS. No gate deshabilitado ni prueba -race inferida.

## FAIL-20260909-632 — onboarding identity boundary

CODE/ISOLATION; REGRESSION_PROVEN. V359 Concatenated tenant/franchisee aliases distinct identities: Start rejects a different checklist, Progress leaks counts and Complete changes the other checklist. String reads map length without writer mutex. Complete scans steps O(N), contradicting metadata O(1); stored order does not enforce prerequisite completion. Original9tests PASS,4regression groups FAIL. Preserve original/red sources and logs. Correct canonical tuple keys without new identity restrictions or business rules; retain error/transition semantics and qualify an independent model. No installed auth bypass or real franchise activation/refund inferred.

## FAIL-20260909-633 — warranty identity boundary

CODE/ISOLATION; REGRESSION_PROVEN. V359 Open validates ID syntax, but State and five transition methods accept lookups whose concatenated identity aliases a valid claim in another tenant. All transitions and both Close branches reproduce a foreign state change. Original9tests PASS,4regression groups FAIL. Preserve original/red sources and logs. Correct canonical tuple keys without new identity restrictions or business rules; retain error/transition semantics and qualify an independent model. No installed auth bypass or real franchise activation/refund inferred.

V359 FAIL632/633: 9 tests históricos intactos, 4 grupos de regresión fallaban antes; 4 fuentes reconstruidas, 17 tests x3, vet/build y modelos independientes con 30 semillas/4760595 ejecuciones fuzz PASS. Claves estructuradas preservan inputs/API/errores. Onboarding String sincronizado; Complete documentado O(N), sin prerequisites nuevos. Garantías conserva nota antes de lookup y ambas ramas Close; diagrama corregido para no implicar Resolve desde rejected. Sin -race ni efecto empresarial real.

## FAIL-20260909-634 — wildcard passed as literal rg path

TOOLING/SEARCH; RECOVERED. V360 rg received markdown_system/*PACK_PLAN.md as an unexpanded Windows path and reported OS123; partial implementation_packs output does not prove plan coverage. Preserved search-path-error.txt. Repeat using -g *PACK_PLAN.md against markdown_system; exit1 with no stderr means no matches, not a filesystem error. No product change from this recovery.

## FAIL-20260909-635 — surveys identity/ownership boundary

CODE/ISOLATION; REGRESSION_PROVEN. V360 Concatenated tenant/survey/customer aliases distinct responses and rejects the second with ErrDuplicate. Responses reads compare stored fields and were not observed leaking across tenants. String reads map length without writer mutex.11original tests PASS,5public API regression groups FAIL across both cores. Preserve baseline/red sources and logs; fix canonical keys and ownership. Marketing requires explicit tenant API migration, no unsafe global shim. No actual recipients/providers/business effects.

## FAIL-20260909-636 — marketing identity/ownership boundary

CODE/ISOLATION; REGRESSION_PROVEN. V360 Global campaign IDs, Due(now), Send(id,recipient) and Remaining(id) cannot isolate callers by tenant. Due returns shared recipient slice allowing external mutation of the store. Concatenated campaign/recipient aliases acknowledgements. String lacks writer mutex. Send only records local flags; no provider send or exactly-once guarantee exists.11original tests PASS,5public API regression groups FAIL across both cores. Preserve baseline/red sources and logs; fix canonical keys and ownership. Marketing requires explicit tenant API migration, no unsafe global shim. No actual recipients/providers/business effects.

V360 FAIL635/636: 11 baseline tests PASS y5grupos red; surveys0.1.1/marketing0.2.0 reconstruyen4fuentes,20tests x3,vet/build y24semillas/3864824fuzz PASS. Marketing exige tenant en Due/Send/Remaining, usa identidad estructurada de campaña/acuse y copia Recipients. Callers antiguos fallan tres firmas; sin shim global. String sincronizado en ambos, sin -race. Aserciones históricas preservadas con migración de argumentos/map index del test. FAIL634 búsqueda corregida sin asumir que la salida parcial fuera exhaustiva.

## FAIL-20260909-637 — recurrence of unexpanded rg wildcard

TOOLING/SEARCH; RECOVERED. V361 repeated FAIL634 by passing reconstruction_evidence/NATIVE* as a literal Windows path. Partial output not used to infer absence. Preserved error summary and corrected to rg -g on the directory; read exact V351 and verifier lines. Never pass shell-style wildcards as rg positional paths in PowerShell.

## FAIL-20260909-638 — POS numerical/ownership boundary

CODE/CORRECTNESS; REGRESSION_PROVEN. V361 Unchecked quantity/product/sum and tender subtraction overflow can accept impossible totals or a negative tender; completed Lines alias the caller; concatenated tenant/id rejects another identity. Total/ChangeDue cannot signal invalid arithmetic with current int64-only signatures. String lacks writer mutex. Baseline originals pass,5regression groups fail across both cores. Preserve original and red evidence. Fix canonical arithmetic/ownership and reject NaN, keeping existing economic rules and no payment/alert effects.

## FAIL-20260909-639 — SLO numerical/ownership boundary

CODE/CORRECTNESS; REGRESSION_PROVEN. V361 NewBudget admits NaN because ordered comparisons do not reject it. This violates target in (0,1) and can suppress meaningful burn/error-budget output. No counter-overflow or complete temporal-window implementation claim is made by this focused test. Baseline originals pass,5regression groups fail across both cores. Preserve original and red evidence. Fix canonical arithmetic/ownership and reject NaN, keeping existing economic rules and no payment/alert effects.

## FAIL-20260909-640 — newline retained in historical stage pointer

TOOLING/PATH; RECOVERED. V361 read elite-v351-current.txt without stripping its trailing newline, producing Errno22 before receipt access. Preserve pointer-path-error.txt; strip pointer whitespace and then verify observed receipt SHA. No missing or corrupt PostgreSQL evidence inferred from this path error.

V361 FAIL638/639:11original tests PASS;5red groups reproduce numeric/identity/ownership errors. POS0.2.0 returns checked (int64,error) totals/change, detects product/sum overflow, compares tender before subtracting, clones Lines and keys by tuple. SLO0.1.1 rejects NaN.4canonical files,20tests x3,vet/build,2legacy signatures rejected and34seeds/3381083fuzz PASS. Existing invalid-line precedence, positive-sale requirement and SLO finite-target calculations retained. Docker scope clarified from verifier plus SHA-verified historical11case/11database V351 receipt; no container/service admission.

## FAIL-20260909-641 — V362 domain boundary

CODE/CORRECTNESS; REGRESSION_PROVEN. Dashboard YYYY-MM accepts months00/13/99, violating valid-period invariant. Command: pinned Go1.26.7 test -json -count=1 ./... in <LOCALAPPDATA>/Temp/elite-v362-afda5e60566943998f33a6559f376d05/red. Baseline8 tests PASS;5red groups. Owner: codex maintenance. Preserve before169/original/red sources/logs. Canonical correction and clean rebuild required. No external data, auth grant or publication.

## FAIL-20260909-642 — V362 domain boundary

CODE/CORRECTNESS; REGRESSION_PROVEN. Help-center concatenated tenant/id aliases distinct tuples: duplicate rejection and foreign Update/Publish/Archive succeed. String also reads map length without writer mutex; race detector not run. Command: pinned Go1.26.7 test -json -count=1 ./... in <LOCALAPPDATA>/Temp/elite-v362-afda5e60566943998f33a6559f376d05/red. Baseline8 tests PASS;5red groups. Owner: codex maintenance. Preserve before169/original/red sources/logs. Canonical correction and clean rebuild required. No external data, auth grant or publication.

## FAIL-20260909-643 — V362 domain boundary

CODE/CORRECTNESS; REGRESSION_PROVEN. Help-center version increment wraps MaxInt64 to MinInt64 in Update and transitions. Boundary is injected through private map in synthetic test, not reached by billions of public API calls. Command: pinned Go1.26.7 test -json -count=1 ./... in <LOCALAPPDATA>/Temp/elite-v362-afda5e60566943998f33a6559f376d05/red. Baseline8 tests PASS;5red groups. Owner: codex maintenance. Preserve before169/original/red sources/logs. Canonical correction and clean rebuild required. No external data, auth grant or publication.

V362 FAIL641/642/643: canonical dashboards/helpcenter0.1.1 rebuilt4sources;8historical assertions retained,5red groups,19tests x3,vet/build,34seeds/4266789fuzz PASS. Months01-12 checked after shape; tenant/id tuple; String mutex; MaxInt64 version rejects unchanged with ErrVersionExhausted. Failure precedence and archived Update retained. Exact numerical and nested-map history oracles; maximum version fixture is internal synthetic state. No external auth/search/source integration or -race claim.

## FAIL-20260909-644 — external authority overclaim

CONTRACT/FRESHNESS; REGRESSION_PROVEN. V363 local SEO metadata/comment attributes60/160 limits to Google; official title-link/snippet pages read2026-09-09 state no fixed maximum. i18n universal n==1 claims CLDR although official plural-rules are language-specific. Preserve before171 and authority-review.json. Correct attribution without changing local business/text policy; no upstream code acquired.

## FAIL-20260909-645 — V363 baseline count assertion

TOOLING/EVIDENCE; RECOVERED. red363.py assumed14 baseline tests, actual Go JSON events prove15 (4/6/5). Assertion stopped evidence aggregation after baselinePASS and expected red failure. Preserve original script/logs and count-error.txt. recover_red363.py reads exact events and requires15; no tests removed, rerun or rewritten as green.

## FAIL-20260909-646 — V363 local boundary

CODE/CORRECTNESS; REGRESSION_PROVEN. i18n SetPlural accepts blank one/other strings and replaces previously valid forms; String lacks reader lock. Pinned Go test -json -count=1 ./... in <LOCALAPPDATA>/Temp/elite-v363-d569bfbd627d436085851aa034ab1c90/red:15baseline tests PASS,8red groups. Preserve before171 and red logs; fix canonical owner then reconstruct and test. No real account, message, URL fetch or user data.

## FAIL-20260909-647 — V363 local boundary

CODE/CORRECTNESS; REGRESSION_PROVEN. SEO prefix-based tenant lookup leaks another tenant URL; hostless URL accepted; LocalBusiness validates only nonblank URL; robots CR/LF/NUL can introduce extra directives. Pinned Go test -json -count=1 ./... in <LOCALAPPDATA>/Temp/elite-v363-d569bfbd627d436085851aa034ab1c90/red:15baseline tests PASS,8red groups. Preserve before171 and red logs; fix canonical owner then reconstruct and test. No real account, message, URL fetch or user data.

## FAIL-20260909-648 — V363 local boundary

CODE/CORRECTNESS; REGRESSION_PROVEN. Social global post ID rejects another tenant same ID. Due/MarkPosted/MarkFailed cannot express caller tenant; Due exposes all tenants. String unsynchronized. No provider dispatch exists. Pinned Go test -json -count=1 ./... in <LOCALAPPDATA>/Temp/elite-v363-d569bfbd627d436085851aa034ab1c90/red:15baseline tests PASS,8red groups. Preserve before171 and red logs; fix canonical owner then reconstruct and test. No real account, message, URL fetch or user data.

## FAIL-20260909-649 — guessed historical evidence filename

TOOLING/SEARCH; RECOVERED. V363 guessed CORE_OWNER_CAPABILITY_AUDIT_V293.md and rg returned OS2. Exact discovery found CORE_OWNER_RECONCILIATION_V293.md. Preserve search-error.txt; discover filenames before reading historical evidence. No missing-owner or blocked-evidence inference.

V363 FAIL644/646/647/648: i18n0.1.1,SEO0.1.1,social0.2.0 canonical6sources;15original test bodies retained with tenant-call migration only,8red groups,32tests x3,vet/build,3legacy signatures rejected,62seeds/4518771fuzz PASS. Official docs correct local60/160 and binary-rule attribution; policies retained as AUTHORED. Observability0.3.1 unchanged,28tests/vet/build isolated PASS,still CANDIDATE. FAIL645 baseline counting and649 guessed path recovered with original evidence preserved.

## FAIL-20260909-650 — numeric core ID omitted by census regex

TOOLING/EVIDENCE; RECOVERED. V363 integration parser only allowed letters/dashes, excluding GO-I18N-CORE and failing21-ID assertion before report/checkpoint. Preserve census-error.txt and initial script. Resume after completed mutations with alphanumeric ID grammar and explicit21unique check. Do not rerun initial appends or fabricate missing core.

## FAIL-20260909-651 — generated continuation newline escaping

TOOLING/SYNTAX; RECOVERED. V363 recovery generator inserted literal newlines in single-quoted Python text. SyntaxError before any execution. Preserve continue171-syntax-error.py.txt and continuation-syntax-error.txt; replace with escaped newlines and ast.parse whole script before run. No checkpoint/file effects from failed parse.

## FAIL-20260909-652 — guessed source profile index path

TOOLING/SEARCH; RECOVERED. V364 guessed OFFICIAL_SOURCE_PROFILES.md, rg returned OS2. Discover actual filenames before exact reads; used existing PROJECT_OFFICIAL_SOURCE_PROFILE_RECORD and PNPM V333-V341 reports. Preserve search-error.txt in V364 stage. No source absence or approval inferred.

## FAIL-20260909-653 — guessed V337 receipt filename

TOOLING/PATH; RECOVERED. A read requested qualification-verification.json in V337 without first discovering its actual filenames. PowerShell returned path-not-found. No qualification result was inferred; V337 canonical report and known source receipts remain the authority. V364 now inventories before reads.

## FAIL-20260909-654 — official research endpoints unavailable through web

RESEARCH/ACCESS; ACCESS_BLOCKED. V364 web probes for the relocated semver-utils fixed commit, its npm version metadata and Yarn pnp4.1.7 attestation returned non-retryable safe-open rejection. These are access observations, not proof of absent sources. No alternate transport bypass or admission inferred. Continue read-only qualification of already captured primary bytes; reopen the network branch when the endpoint becomes accessible through the authorized tool. FAIL575 remains unresolved.

## FAIL-20260909-655 — manual compressed-byte count rejected

EVIDENCE/REPORT; RECOVERED. V364 report generator assumed12373 compressed bytes; qualification receipts independently measure12375. Assertion stopped before checkpoint173 or contract promotion. Preserved generator and pre-fix report in V364 stage. Canonical report now uses measured12375; decoded72463 and all correspondence proofs unchanged. Generate quantitative prose from measured receipts rather than hand-derived padding assumptions.

## FAIL-20260909-656 — structural inventory cursor stale

STRUCTURAL_GATE; REGRESSION_PROVEN. VERIFY_LIBRARY173 rejected the current FRANCHISE_PREFLIGHT_GAP header at801Markdown after V364 added the802nd. The failed log is preserved in V364 stage. Update only the current inventory to162/1461/802/53 and versionV364; historical evidence remains. No verifier change. Checkpoint174 precedes the required replay.

V364 FAIL656 replay174: VERIFY_LIBRARY_PASS162/1461/802/53. Failed173 remains intact; no verifier weakening.

## FAIL-20260909-657 — materializer incorrectly used as overlay

TOOLING/MATERIALIZATION; RECOVERED. Recurrence of FAIL620/LIB-FAIL-2336 and LIB-FAIL-1753: V365 attempted the second pack into isolated21 after dashboards. The empty-destination guard rejected FX before writes. Preserve failed log and first output. Correction: one new absent directory per pack, then a separate analysis-only module with collision checks and byte equality. No admission/compositor guard bypass.

## FAIL-20260909-658 — compatibility metadata references unresolved owners

CONTRACT/PROVENANCE; REGRESSION_PROVEN. V365 reads all21current cores and all canonical pack IDs/versions: six claims do not resolve to the declared range. FX/giftcards/promotions refer to old commerce ranges; observability/payroll use nonexistent GO-ENTERPRISE-BACKEND-CORE; SLO refers to observability0.1.x although current0.3.1 is CANDIDATE. Preserve exact before refs/SHA in compatibility-invalid-before.json. These helpers import only stdlib and have no existing adapter to those owners. Correct six metadata declarations to isolated test scope, not a guessed latest version; preserve historical wording and all code hashes. Full domain equivalence/integration FAIL385 remains DIAGNOSED.

## FAIL-20260909-659 — assumed materialization receipt filename

TOOLING/EVIDENCE; RECOVERED. V365 assumed MATERIALIZATION_RECORD.md after a successful isolated FX rebuild. Inventory and the canonical materializer tail show only manifest source files and a stdout count are produced. Preserved failure script and successful FX log; use its log SHA plus explicit reconstructed-file digests, not a fabricated receipt. Continue remaining five absent destinations and verify the already rebuilt FX bytes without overwriting.

## FAIL-20260909-660 — guessed ML authority filename

TOOLING/ROUTING; RECOVERED. V366 searched nonexistent ML_PRODUCTION_MLOPS.md. The already observed master map names ML_PRODUCTION_LLMOPS_EVALUATION.md; read its provenance/lineage and deployment gates plus the existing AI security training boundary. No absent ML capability inferred from the path error. Use exact paths from the authority map.

## FAIL-20260909-661 — consumer training requirement confused with library intake

CONTRACT/COMMUNICATION; SCOPE_CORRECTION_VERIFIED. The assistant asked for real customer chat history to unblock library completion, while V288 framed learning mainly as memory/knowledge/evals. The user clarified NEW has no history, EXISTING has history, and the purpose is training a model with elite code and explicit use. Existing ownership/safety controls stand. New canonical PROJECT_HISTORY_MODEL_TRAINING_CONTRACT separates library, per-project original/derived datasets, explicit offline training, evaluation and runtime promotion. REQ05/TEST09 are corrected to that purpose; no status promoted. Real corpus/channel/model/cost/authorization are consumer inputs, not demanded now as library data. No training, data acquisition, upload or new executable pipeline occurred. Historical V288/289 remains; before177 records exact old contract.
## FAIL-20260909-662 — directory passed to file reader

TOOLING; RECOVERED. V367 appended Get-Content for the existing V337 stage directory to a metadata read. PowerShell returned exit1; -ErrorAction SilentlyContinue hid that final reader message. The preceding metadata and transport searches completed. No product/source mutation occurred. Correction: remove that unrelated directory read; inspect only the already observed metadata-index.json. Preserve the failed command/exit in V367/tooling-failure662.json. Never suppress inspection errors or treat a directory as a receipt file.

## FAIL-20260909-663 — training capability admission remains incomplete

CAPABILITY_GAP; RESEARCH_INCOMPLETE. V367 materialized the existing gate and ran the real REQ05 record. It rejects with exit2 and emits no admission receipt: official sources are discovered, but artifact acquisition, full dependencies/licenses and executable training/promotion gates are absent. Preserve reconstruction_evidence/training_gap_v367/record.json and gap-investigation.log. Do not mark CONDITIONED or fabricate source hashes from research notes. Continue governed acquisition/admission of the narrow TRL/optional PEFT lane; consumer private corpus is not this blocker. No HUMAN_DECISION_REQUIRED or source exhaustion inferred.
## FAIL-20260909-664 — research JSON outside the portable release allowlist

STRUCTURAL_GATE; REGRESSION_PROVEN. V367 checkpoint180/resume/plan passed, but VERIFY_LIBRARY180 rejected reconstruction_evidence/training_gap_v367/discovered-roster.json as an unknown release file. The library distributes approved Markdown authorities, not arbitrary new raw research payloads. Preserve failed log and all JSON bytes before correction. Embed the hash-bound research records in the canonical V367 Markdown report and keep their materialization in the isolated stage; remove only the newly introduced raw research directory from the release tree after a verified snapshot. Do not expand the release allowlist or weaken its verifier. Repeat structural gate after checkpoint181.
## FAIL-20260909-665 — guessed portable authority filename

TOOLING/ROUTING; RECOVERED. During FAIL664 diagnosis, rg included nonexistent markdown_system/PORTABLE_RELEASE_CONTRACT.md. Its error is preserved in the V367 repair record; no absence of a release authority is inferred. The already-read VERIFY_LIBRARY.ps1 lines12-17 and227 explicitly define reconstruction_evidence as Markdown-only. Use this observed executable owner and the existing maps, not invented filenames. No verifier changes.
## FAIL-20260909-666 — text embedding normalized CRLF evidence

EVIDENCE/SERIALIZATION; REGRESSION_PROVEN. V367 repair181 stopped before checkpoint181 because read_text normalized the preserved gap log's CRLF bytes to LF inside a Markdown fence. Hash comparison correctly rejected the mismatch. All originals remain in before181 and research-root. Preserve original repair181.py. Encode the log as explicit base64 with its original SHA-256, retain JSON as UTF-8, then verify all eight decoded records byte-for-byte before replaying the structural gate. Do not replace the original receipt hash with a normalized-text hash.
V367 FAIL666, segundo intento: la aserción identificó metadata-index.json, también escrito originalmente con CRLF por Python en Windows. La corrección limitada al log era insuficiente. Se conserva repair181-recovered.py y el reporte anterior. Aplicar envelope base64 a los ocho registros, verificar bytes originales de todos y recién después registrar checkpoint; no suponer LF por extensión .json.
## FAIL-20260909-667 — inspection scope flag omitted metadata copies

EVIDENCE/SCOPE; REGRESSION_PROVEN. V369 bounded TAR inspection did not extract an archive tree or execute package code, but did copy selected manifests/VERSION into the isolated stage. Its extracted_to_filesystem=false field was too broad. Preserve the original source-correspondence.json and inspection script; replace that field with archive_tree_extracted=false plus exact selected_metadata_files_copied paths/hashes, and verify the copied bytes against the archive inventory. Transport receipts remain accurate for their opaque acquisition step. Do not erase or relabel the copied evidence as if no bytes were written.

V369 FAIL663 partial progress: exact TRL/PEFT sdists acquired under validated profile;559regular files inspected,544match source Git blobs,15generated metadata files reviewed. Generic training admission remains RESEARCH_INCOMPLETE: full dependency/runtime/security/fixtures/promotion evidence absent. Do not repeat completed source acquisition or ask for consumer corpus to prepare the library.

## FAIL-20260909-668 — inventory header formatting assumption

EVIDENCE/INTEGRATION; REGRESSION_PROVEN. V369 integrate185.py stopped at the current preflight header assertion before checkpoint185: expected ungrouped 1461 and a single line, but the canonical text contains 1.461 followed by a newline. Earlier mutations are preserved; state remains184. Preserve original script and observed header; repair only this exact header and resume the unexecuted tail. Do not rerun already applied composition/version mutations or weaken any gate.

## FAIL-20260909-669 — unverified previous preflight receipt path

TOOLING/ROUTING; RECOVERED. V369 inspection attempted to read a prior-stage preflight183.json without first locating it. PowerShell reported path absent; the subsequent source reads succeeded, so the overall shell exit0 did not validate that read. No old preflight result is used. Read the canonical Write-Evidence schema and wait for the actual new preflight185.json instead; no acquisition or gate mutation.

V369 FAIL668 regression: checkpoint185/resume/plan and Preflight185 structural/audit checks passed;availability remains BLOCKED; all53composition counts equal predicted values,162packs/1464files/808Markdown verified. Only unexecuted integration tail resumed.

## FAIL-20260909-670 — current roadmap table retained an older inventory

EVIDENCE/FRESHNESS; REGRESSION_PROVEN. V370 found the roadmap current Estado real table still lists1461materializable files (historical Preflight149) while its V369 current header and Preflight185 prove1464. before187 preserves the exact original. Update only this current table to1464 with Preflight185 attribution, retain historical evidence and verify the full structural gate; no pack or completed-control status changes.

## FAIL-20260909-671 — nested script delimiter collision

TOOLING/SERIALIZATION; REGRESSION_PROVEN. V370 outer Python script generator used triple-single-quoted text containing the same delimiter in the generated Markdown strings. Parsing failed before creating integrate187.py or applying any integration mutation; the subsequent invocation correctly reported the script absent. Preserve the command transcript and exact errors. Write the script directly with a literal PowerShell here-string, syntax-check it before execution, and keep original research receipts intact.

V370 FAIL663: subpaso de resolución metadata cerrado59/103;59METADATA hash-verified,resolver pip independiente+2negativos,OSV59versiones0hallazgos. Siguen bytes/licencias/nativos,firmas cuando apliquen,runtime y pipeline. TEST09 no cierra por consulta de nombres/versiones.

## FAIL-20260909-672 — machine-local path in portable evidence prose

STRUCTURAL_GATE/PORTABILITY; REGRESSION_PROVEN. VERIFY_LIBRARY187 rejected TRAINING_DEPENDENCY_GRAPH_V370.md because its prose included the absolute machine-local stage path [machine-local home]. Original report and failed log preserved. Replace only the prose stage reference with %TEMP% plus the observed stage basename; hash-bound raw research records remain byte-exact inside their explicit base64 envelopes. Do not weaken the portable-path verifier. New checkpoint188 and structural replay required.

V370 FAIL672 recurrence: replay188 correctly found the same absolute home copied into the diagnostic lesson itself. Original lesson and failed188 preserved. Replace the diagnostic path with a neutral placeholder; scan every Markdown with the same home-path pattern before checkpoint189. This is a scope correction to the repair, not a relaxed verifier.

V370 regressions670/671/672: roadmap current count1464 verified; directly written integration script syntax-check/checkpoints passed; original report/lesson and failed187/188 retained. All826Markdown home-path scan has0hits; full VERIFY_LIBRARY189 PASS with162/1464/809/53. No source/verifier weakening.

## FAIL-20260909-673 — next action exceeded compact cursor limit

CHECKPOINT/CONTRACT; REGRESSION_PROVEN. Checkpoint190 rejected next_action longer than800characters before appending EVT-0190. Preserve failed state and checkpoint190.log. The detailed roadmap belongs to canonical owners; reduce the cursor to the immediate governed acquisition step and reference V370 for remaining work. Do not change the validator limit or repeat completed integration.

V370 FAIL673: checkpoint190 retry, resume190 and plan190 PASS with concise cursor; no event appended by failed attempt. Closure191 retains this evidence.

## FAIL-20260909-674 — wheel RECORD root selection

INSPECTION/ASSUMPTION; REGRESSION_PROVEN. V371 wheel inspection stopped because its suffix-based RECORD selector found a count other than1. Preserve original script and partial inventory/metadata copies. Identify root versus nested vendored dist-info directories before correcting the selector; do not weaken any RECORD hash or ignore an unexpected package root.

## FAIL-20260909-675 — embedded Rust component advisories in candidate wheels

DEPENDENCY/SECURITY; DIAGNOSED. V371 inspected actual wheel SBOMs and found6OSVrecords over embedded crates: memmap2 0.9.10 and pyo3 0.28.3 in safetensors0.8.0, paste1.0.15 in tokenizers0.23.2. Some records may be aliases or maintenance advisories; do not conflate records with unique vulnerabilities. Existing59PyPI query had no results but did not cover these embedded identities. Preserve SBOMs/request/response, keep wheels quarantined/uninstalled and inspect official advisories/fixed releases before a candidate replacement or declared source rebuild. No ignore or admission inferred.

V371 FAIL674 regression: root dist-info RECORD selected by path depth2, preserving nested vendored metadata. All59wheels inspected;23795RECORD hashes verified over23854files. FAIL675 remains open:6records are4distinct advisories (3security/unsoundness plus1unmaintained notice). Official safetensors source candidate pins fixes pyo3 0.29.2/memmap2 0.9.11;45registry crates query0findings, no source build or wheel replacement yet.

## FAIL-20260909-676 — assumed upstream lock path

RESEARCH/ASSUMPTION; REGRESSION_PROVEN. Isolated V371 read-only probe assumed tokenizers root Cargo.lock and received404. First recovery queried the exact tree but asserted another unobserved path. Scripts and returned bytes preserved in V371; local PROJECT_FAILURE_LESSONS opened immediately while Preflight192 snapshot ran. Complete exact tree contains4Cargo.toml and0Cargo.lock. Corrected probe reads observed manifest only and makes no resolved-graph claim. Root record now includes the isolated lesson; no source/lock/production mutation arose from failed probes.

## FAIL-20260910-677 — wildcard passed as a literal rg path

TOOLING/SEARCH; RECOVERED. Read-only V372 discovery passed a wildcard as a positional Windows path and received error123. No absent-file or core-admission conclusion is inferred. Discover exact filenames with rg --files then read those paths; this repeats the existing path-discovery lesson. No source mutation.

V372 FAIL677 recurrence: another positional wildcard was used during plan discovery; error123 preserved. Corrected search uses directory plus -g filter. No source selection or absence inferred.

## FAIL-20260910-678 — SBOM purl version assumption

INSPECTION/ASSUMPTION; DIAGNOSED. Official Collector SPDX includes Go purls without an @version. Initial query preparation raised ValueError before sending a batch. Preserve original SBOM/script; classify missing registry versions as unresolved identities and inspect versionInfo rather than inventing versions or omitting them from coverage claims.


## FAIL-20260910-679 — Collector contrib binary graph is not clean

DEPENDENCY/SECURITY; DIAGNOSED. Exact official0.160.0 Windows contrib SPDX contains1048versioned Go identities plus obi UNKNOWN; OSV returns7records. The distribution is not installed/executed/admitted. Preserve all requests/responses and consult official advisories. Evaluate a minimal official Collector distribution and required components instead of silently ignoring unused integrations.


## FAIL-20260910-680 — contrib archive is incompatible with Windows acquisition

PROFILE/COMPATIBILITY; DIAGNOSED. Candidate profile/lock suites reject the contrib source archive because it contains a symlink and lacks the required non-Windows restriction/kind. Transport suites55/98 pass. Preserve candidate40 and failed logs. Do not weaken the existing platform guard or flatten links: exclude this whole-repository archive from the Windows profile. Selected contrib Go modules require their own official module identity/checksum/license evidence; no whole-repository admission inferred.

V372 FAIL677 recurrence retained: read-only runtime-evidence search again passed a literal wildcard path and returned123. No result was relied on; subsequent searches use existing exact paths or directory with -g.
V372 FAIL679 extends to official OCB builder graph:32queries return GO-2026-5024 for x/sys0.41.0; build guard stopped before producing/running the builder. Preserve graph/request/response and inspect the official fixed version before a declared dependency-only candidate delta.

## FAIL-20260910-681 — reporter TLS fixture did not reach the sink

TEST/IDENTITY; DIAGNOSED. Initial new reporter tests receive TLS bad-certificate before handlers, so positive/response-boundary assertions fail and the cancellation fixture waits for an event that cannot arrive. Preserve original test/source and log. Inspect actual certificate/httptest behavior, bound test waits, and repair the fixture without disabling verification or treating TLS failure as a passed sink test.

V372 FAIL681 cause: the synthetic leaf was signed with its own key instead of the CA key. Independent x509.Verify reproduces the signature failure. Preserve original test and diagnostic. Use CA signer, add finite handler-entry wait; terminate/reap only the identified owned hung test process. TLS verification remains enabled.

## FAIL-20260910-682 — minimal Collector configuration scheme

RUNTIME/CONFIGURATION; DIAGNOSED. The first real Collector validate command rejected DefaultScheme because the generated distribution selected only fileprovider while its default resolver scheme remained env. No service was launched. Preserve original generated source/config/binary and failed runtime receipt. Set the official builder default_uri_scheme to file, regenerate into an absent directory, retain the explicit security dependency patch, rebuild and replay validate; do not add environment/network providers or bypass validation.

V372 FAIL677 recurrence: read-only contract discovery used an unobserved name and positional wildcard. Corrected by rg --files and exact PROJECT_ENGINEERING_CONTRACT.json; no absence claim.

## FAIL-20260910-683 — real TLS client interoperability failed

INTEGRATION/IDENTITY; DIAGNOSED. Collector validate and promtool config passed after682, but the Python reference client's authenticated scrape repeatedly failed with TLS bad record MAC. Preserve full Collector log and run receipt; no healthy-sink claim. Isolate the exact Python/OpenSSL versus Go client handshake and certificate chain before changing TLS policy. Keep authentication and certificate verification enabled. Readiness probes must also accept an HTTP200 empty metrics body.
V372 FAIL683 diagnosis: the exact client reports Missing Authority Key Identifier under strict OpenSSL certificate verification. Add standard SKI/AKI and leaf digitalSignature usage to ephemeral fixture certificates. Preserve strict verification and TLS1.3; this is a certificate-fixture defect, not an observed Go/OpenSSL interoperability vulnerability. Preserve failed script and direct client diagnostic.

## FAIL-20260910-684 — host reference environment rejected below launch profile

INTEGRATION/CONTRACT; DIAGNOSED. Authenticated Collector and Prometheus endpoints started, but the native host launch receipt is REJECTED: the v3 profile accepted its four explicit reference runtime variables while native_capture still only accepted four OS variables. No worker ran. Preserve result receipt and source. Extend the low-level boundary only behind an explicit reference_environment flag with the exact reference key set; retain the default four-key behavior and reject arbitrary provider/PATH/stop-handle injection. Add regressions at both boundaries.

## FAIL-20260910-685 — replay test supplied a different invalid profile

TEST/CONTRACT; DIAGNOSED. The actual main produced all three signals and Prometheus scraped its counter, but the replay assertion supplied {} instead of the original digest-bound launch profile and therefore correctly did not return ALREADY_RESERVED. Preserve the failing run. Replay must use the exact original bytes; malformed/digest-mismatched profiles remain rejected and cannot be counted as replay acceptance.

## FAIL-20260910-686 — upstream Prometheus Windows expected fixtures are stale

UPSTREAM_TEST/PORTABILITY; DIAGNOSED. Official Collector component suites and selected builder tests pass; Prometheus rules pass but config has two failed equality tests. Actual diff shows Windows path separators and missing OTLP defaults in the Windows expected fixture. Preserve original suite and source before diagnosis. Compare non-Windows fixture and source behavior, then declare a test-only portability adaptation with unchanged assertions; no runtime fix or full upstream-green claim inferred. A follow-up guessed config_unix_test.go path failed; actual inventory identifies config_default_test.go.

## FAIL-20260910-687 — supervisor suite invocation omitted fixture arguments

TEST_RUNNER/CONTRACT; DIAGNOSED. Capture ran22checks then its real-binary check raised IndexError; result-store suite failed before tests because the runner omitted their required executable/fixture arguments. Launch20 and shutdown19 pass. Preserve original invocation/logs; use the exact executable for capture and hash-bound fixture JSON for result-store. No runtime defect or skipped PASS inferred.
V372 FAIL687 extends to fuzz gate self-test invocation: mandatory GoExecutable omitted, rejected before execution. Preserve fuzz-self.log; corrected wrapper supplies the exact pinned executable. Capture invocation also required explicit SHA-256; capture23/result28/launch20/shutdown19 now all pass.
V372 review correction: the integrated main must preserve native watcher protocol failure after the loop returns on cancellation. Candidate now returns nativeStopFailure(ctx), with a specific Windows regression distinguishing protocol failure from cooperative stop. This prevents a clean process result from hiding watcher failure; final rebuild/integration will use the corrected executable.

## FAIL-20260910-688 — portable reference lock generator helper scope

TOOLING/GENERATION; DIAGNOSED. The portable reference source was written, but its local lock generator raised NameError because sha was only present inside the generated source string. Preserve original generator; define the helper in the generator itself and replay candidate generation. No lock, admission or completed-control claim arose from the failed invocation.

## FAIL-20260910-689 — refund README canonical newline hash

MATERIALIZATION/INTEGRITY; DIAGNOSED. Fresh refund reconstruction rejected README SHA-256 because the candidate README was written with Windows CRLF while the Markdown block materializes canonical LF. Preserve rejected candidate and log context. Normalize candidate UTF-8/LF before hashing and reconstruct into a new destination. Do not change the materializer hash check. Reference runtime13files already reconstructs correctly.

## FAIL-20260910-690 — inventory parser assumed no blank metadata lines

INVENTORY/ASSUMPTION; DIAGNOSED. Canonical integration updated the three packs and composition versions, then its auxiliary provenance assertion undercounted because it required a YAML fence immediately after each FILE heading. Existing valid packs include intervening blank lines. Preserve the partially applied integration and original parser; count the first provenance metadata of each exact FILE section, matching the canonical materializer. A raw global provenance grep also includes one source example and is not the block census. Do not relax expected inventory or rerun already-applied version mutations.

## FAIL-20260910-691 — metric timestamp variable shadowed process starter

RUNTIME/NAME_BINDING; DIAGNOSED. Fresh reconstructed reference suites96tests pass, but the strengthened integration fails before launching Collector because a local metric variable named start shadows the start() helper throughout main. Preserve failed canonical pack and reconstructed run. Rename only the metric timestamp binding to series_started, reconstruct from Markdown again and replay integration. No process/runtime PASS is inferred from the failed invocation.

## FAIL-20260910-692 — PostgreSQL startup lost a required clean-environment setting

RUNTIME/ENVIRONMENT; DIAGNOSED. With the new closed service environment, Collector/Prometheus mutual-TLS checks passed but initdb exited without output before creating the reference database. Preserve exact environment policy, failed log and stopped-process receipt. Diagnose the actual Windows exit code and required OS/runtime variable; do not restore the full user environment or infer an absent tool from this failure.
V372 FAIL692 isolated bootstrap: minimal environment reproduces exit0xC0000005; adding only COMSPEC pointing to the Windows System32 cmd.exe succeeds. A controlled OS-only PATH also succeeds, but is not selected. Canonical runner supplies fixed COMSPEC and keeps user PATH, proxies, OTEL attributes and provider variables absent. This is a demonstrated reference precondition; no claim of an upstream patch.

## FAIL-20260910-693 — upstream component tests modified module-cache tree

SUPPLY_CHAIN/INTEGRITY; DIAGNOSED. Final go mod verify rejects the fileexporter0.160.0 module directory after upstream component tests. Preserve modified bytes and compare every entry with the checksum-bound Go module ZIP before attributing the difference. No final source-integrity PASS or binary promotion is inferred. Quarantine the changed tree, restore by exact module extraction, verify all graphs and bind any rebuild to the tested artifact.

V372 final recovery evidence: FAIL680/681/682/683/684/685/686/687/688/689/690/691/692/693 REGRESSION_PROVEN for their narrow corrected cases. Final canonical runner and actual host complete all integrated assertions; Go tests101PASS/1explicit opt-in skip, final policy6/0skip and native suites90PASS. FAIL693 exact module ZIP/tree comparison finds only one added testdata/log.json, all35original files unchanged. That owned test output is preserved outside the module cache; four go mod verify graphs pass and Collector rebuild is byte-identical. Full-release security rejection FAIL679 is retained: selected artifacts use corrected locks with separately evidenced NOT_PRESENT package findings, not a claim of zero module advisories. FAIL677 search wildcard recurred during closure read-only inventory; corrected using exact observed roadmap path, no inference from failed search. Historical diagnoses/logs remain unchanged.

## FAIL-20260910-694 — reference pack omitted mandatory contract sections

STRUCTURE/CONTRACT; DIAGNOSED. Full Preflight194 rejects the new reference pack because mandatory section 6 Configuration surface is absent outside its code fences. Original pack and failed Preflight preserved in V372 staging. Correct the canonical document structure without weakening the structural verifier; rerun affected reconstruction/checkpoint and complete Preflight. Runtime PASS does not substitute for this gate.


V372 FAIL694 REGRESSION_PROVEN: canonical ten-section pack contract restored without changing any of13materialized files. Complete Preflight195 passes all executed steps and all54composition counts. Both failed Preflight194 and before-pack preserved in V372. TEST05 alone closes after actual runtime and canonical verification.

## FAIL-20260910-695 — release tag spelling assumed during source research

RESEARCH/ASSUMPTION; RECOVERED. Web open for safetensors v0.9.0rc0 returned404. Preserve query/result; PyPI version syntax does not establish Git tag spelling. Exact official releases API was already acquired successfully; use its returned tag and commit before reads. No archive/runtime execution or source promotion occurred.


## FAIL-20260910-696 — prerelease wheel profile rejected before acquisition

ACQUISITION/CONTRACT; DIAGNOSED. Reconstructed source-profile validation rejects the new safetensors candidate; exact log retained in V373. No candidate wheel downloaded or run. Dependent inspection also failed because the prior script never emitted it; future dependent shell steps must check exit status. Read the validator reason and correct candidate metadata only if consistent with the canonical transport contract.


V373 FAIL696 diagnosed two independent constraints: profile lock must sit beside its canonical acquirer, then version parser supported only final/post releases. Exact numeric rcN inspection is now a declared0.4.83 extension with19new regression checks (74total), preserving all binding/no-admission semantics and prior stable-opaque behavior. Fresh40file rebuild matches;98opaque/21profile/7negative/9upstream suites pass. Dependent commands now gate on subprocess exit. No prerelease runtime promotion follows.

## FAIL-20260910-697 — candidate wheel installation exceeded its finite budget

ENVIRONMENT/PROCESS; DIAGNOSED. Offline58wheel pip install exceeded240seconds after partial installation; wrapper killed its direct pip process, so inspect and stop only any observed owned descendants before reuse. Preserve partial venv and original log; do not infer completion from files. Next qualification uses a fresh destination, disables unnecessary bulk bytecode compilation and supervises the process tree with a longer finite budget. No package imports, model job or control closure claimed.

## FAIL-20260910-698 — Windows search used positional wildcards again

TOOLING/SEARCH; RECOVERED. Read-only rg passed PROJECT* and markdown_system/POST_DEEPSEEK* as literal paths and exited1 with error123. No file inference or code change follows. Use rg with the root and -g filters or exact enumerated paths. This repeats historical FAIL677; preserve recurrence rather than hide it.

V373 FAIL697 REGRESSION_PROVEN for installation: fresh runtime2 installed58 hash-locked wheels offline with bytecode compilation disabled, native Job supervision,104738816 peak committed bytes, exit0 and empty owned process tree; pip check exit0. Partial runtime1 remains rejected and preserved. Actual model qualification is separate.
## FAIL-20260910-699 — native read rejection leaked through governance API

INTEGRATION/ERROR_CONTRACT; DIAGNOSED.19-test candidate suite reports two errors for concurrent/partial imports: native pinned-read correctly rejects an absent manifest but exposes governed_launch.Rejected instead of the pipeline Rejected contract. No unsafe import or overwrite occurred. Preserve candidate and failed test evidence; translate only the native rejection at the read adapter, then replay the same crash/concurrency tests. Do not weaken the missing-manifest gate.
## FAIL-20260910-700 — private-storage ACL verifier rejects newly private directory

SECURITY/ABI; DIAGNOSED. The strengthened ACL gate rejects freshly created0700 reference directories before importing data. Preserve rejected suite; inspect actual owner/DACL and native parser layout before changing policy. Do not disable the ACL gate or broaden allowed principals. Prior19test PASS precedes this new gate and is not final.
## FAIL-20260910-701 — catalog path assumed instead of enumerated

TOOLING/READ; RECOVERED. Read-only Get-Content requested absent IMPLEMENTATION_PACKS_CATALOG.md. No inventory or admission inferred. Enumerate catalog paths with rg --files and read the observed canonical file. Existing runtime E2E log independently records PASS; final reconstruction remains pending.
## FAIL-20260910-702 — copied wheel-lock text retained CRLF

MATERIALIZATION/INTEGRITY; DIAGNOSED. Pack generator rejected a copied candidate text file with CRLF before publishing the new canonical pack. Preserve those candidate bytes; normalize only text-file transport newlines to canonical LF and regenerate hashes. Wheel artifact hashes and bytes are unchanged; do not weaken the materializer.
## FAIL-20260910-703 — inventory mixed local maintenance and distributable files

INVENTORY/ASSUMPTION; DIAGNOSED. Auxiliary rglob/FILE regex counted835Markdown and1507payloads, including local maintenance files and one nested example. This is not the release selector used by VERIFY_LIBRARY. Preserve that auxiliary result; invoke the canonical read-only manifest/provenance/release selector functions and correct current inventory/notices before checkpoint/Preflight. No control percentage changes.
## FAIL-20260910-704 — canonical inventory helper omitted strict UTF-8 binding

TOOLING/INTEGRATION; RECOVERED. Extracted canonical provenance reader requires strictUtf8 from VERIFY_LIBRARY; the wrapper omitted it and failed before producing an inventory. Add the exact UTF8Encoding(false,true) binding, preserve the wrapper and rerun without changing the canonical parser.

V373 FAIL699/700/702/703/704 REGRESSION_PROVEN for the corrected adapter, actual private OWNER RIGHTS DACL, LF reconstruction, canonical inventory selectors and strict UTF8 dependency. Canonical census164packs/1506files/818Markdown/55profiles;1259AUTHORED/140ADAPTED/107VERBATIM. Original diagnosis and failed candidates preserved. SDK gap license-evidence reference is corrected from instructions to LICENSE.md by a separately validated successor receipt, preserving before/after.
## FAIL-20260910-705 — full verifier did not compose the newly discovered training plan

STRUCTURE/COMPOSITION; DIAGNOSED. Complete Preflight198 reconstructs packs and validates checkpoint198, then rejects discovered55/composed54 plans. New training plan passed its isolated real composition but was absent from VERIFY_LIBRARY's explicit composition table. Preserve failed Preflight log; add its exact21file/two-pack entry to the canonical table without weakening full-coverage assertion. No TEST09 closure before a new full PASS. The final JSON receipt is absent because the gate failed, not a zero-step success.

V373 FAIL705 REGRESSION_PROVEN: complete Preflight199 executes160steps PASS and composes all55plans; the exact21file table addition fixes coverage. Overall availability remains BLOCKED only Docker. Final audit then strengthens exact decision retention and non-null baseline rollback:27policy tests plus fresh actual SDK/runtime pass; current-artifact full Preflight200 remains required before TEST09 closure.

V373 final closure: Preflight200 repeats160executed steps PASS on the final27test/audit-complete artifact, preserving all55composition coverage checks. FAIL695-705 corrections are preserved in their stated narrow scopes; acquisition/SFT/lineage/rollback never admit full native security or private consumer data. TEST09 alone closes after all required library evidence. Remaining TEST02/03/07 are not hidden or relabeled.

## FAIL-20260910-706 — contract path assumed during continuation

TOOLING/READ; RECOVERED. Read-only rg included absent markdown_system/LIBRARY_COMPLETION_CONTRACT.md (os error 2); other explicit inputs returned matches but the combined search is not a complete contract read. No state or coverage inferred from the missing file. Discover the actual normative owner from checkpoint evidence paths and read it directly; preserve this failed path instead of retrying guessed names.

## FAIL-20260910-707 — source comments overstate SRE and aggregation integration

CONTRACT/ATTRIBUTION; DIAGNOSED. V374 inspection finds the SLO source header and MultiWindowAlert comment claiming Google SRE multi-window alerting despite no window lifecycle/ingestion, and dashboard header saying aggregation from domain stores despite a caller-supplied Input only. Narrow current pack claims already disclaim these integrations. Preserve source bytes and wording; align comments with the actual cumulative/formula and pure-input contracts, retain official comparison proofs and reconstruct canonical sources. Do not change runtime behavior or declare full SRE/domain equivalence.

## FAIL-20260910-708 — comparison intake used an unobserved telemetry alias

CONTRACT/ASSUMPTION; RECOVERED. Draft V374 comparison intake reported unresolved GO_RUNTIME_TELEMETRY_ADAPTER; no such canonical pack exists. No audit report or admission used this alias. V372 identifies GO_OFFICIAL_RETURN_REFUND_WORKER as the reporter owner. Resolve the observed path, preserve the draft and change the intake to reject any unresolved owner rather than continue with a partial owner set.

## FAIL-20260910-709 — oversized report orchestration did not parse

TOOLING/SERIALIZATION; RECOVERED. The JavaScript orchestration for the V374 report failed to parse before any nested shell call or file write (Invalid or unexpected token). No report was emitted and no verification result is inferred. Split report construction into bounded scripts and use raw literal transport with explicit PowerShell here-string boundaries; retain failure identity and verify the emitted report before checkpointing.

V374 FAIL709 recurrence_count=2: a second orchestration parse failed before execution because embedded Markdown triple-backticks terminated the JavaScript raw literal. The first smaller raw-literal script was valid, but this delimiter boundary was still unguarded. Generate Markdown fences from character96 inside the transported script; no nested raw-literal delimiter may occur in its source. No canonical verifier was changed by either rejected invocation.

V374 FAIL707 REGRESSION_PROVEN:42fuentes canónicas reconstruidas,228tests sin skip/fallo, las seis comparaciones nuevas y13negativos del verificador PASS; sólo dos comentarios de implementación cambian, no API/comportamiento. FAIL706/708/709 recuperados con resolución de owners completa, reporte emitido y hashes/nombres de tests vinculados. FAIL385 de integración empresarial sigue DIAGNOSED; no se borra para cerrar una auditoría más estrecha.

V374 final: fullPreflight202160steps PASS, including the canonical21source bindings and13negative regressions. Final-source23target fuzz PASS supplements the preserved baseline. FAIL706-709 corrections verified; FAIL709 recurrence2 retained. TEST02 closes only its unchanged audit oracle. FAIL385 remains DIAGNOSED for broader functional integration, with no deletion or fabricated equivalence; TEST03/07 remain blocked.

## FAIL-20260910-710 — audit completion was promoted before its integration blocker closed

CONTRACT/ACCEPTANCE; DIAGNOSED. Checkpoint203 changed TEST02 to passed after valid
source/reference suites but while FAIL385 still explicitly required functional
equivalence/integration. Preserving the literal oracle did not satisfy that agreed
broader blocker. Final self-review caught the premature promotion before delivery
of a completion answer. Preserve checkpoint203 and its reports; restore TEST02
to blocked through successor204, retain the useful tests/source comparisons and
continue the actual owner integration work. Audit completeness is not control
completeness; no bookkeeping-only promotion is an acceptable shortcut.

## FAIL-20260910-711 — low-traffic reproduction receipt not yet classified

TEST/HARNESS; DIAGNOSED. Promtool returned1 on the new low-traffic rule suite, but the wrapper assertion expecting a particular diagnostic phrase failed. The nonzero exit is not yet proof of the suspected rule defect. Preserve exact log and fixtures; inspect whether failure is rule semantics, fixture/schema or output-format assumption before changing any canonical rule.

V374 FAIL711 RECOVERED: promtool used got:[] without the wrapper's assumed space. The actual log contains exactly the two expected missing alerts (low traffic and reset at low traffic); no fixture/schema error. Match structural whitespace and preserve the original wrapper failure.

## FAIL-20260910-712 — error-rate denominator suppresses low-traffic outages

CODE/OPERABILITY; DIAGNOSED. Official-source Prometheus3.14.0 promtool, exact admitted SHA2d8428b4..., evaluates the canonical SECURE-OPS-DELIVERY-CORE1.1.3 PlatformHighErrorRate rule. With one failed request per four minutes and no successes, the true error fraction is1 but clamp_min(total_rate,1) produces about0.00417, so the1%/10minute rule never fires. The same defect persists across a counter reset. Both expected alerts are absent; high-traffic and healthy controls behave normally. Preserve canonical baseline and red log. Divide by actual total rate, retaining zero/no-series non-firing semantics, the1% threshold and10minute hold; add permanent source-engine regressions and reconstruct all affected consumers. This is a real rule defect, not a substitute for closing TEST02 as a whole.

## FAIL-20260910-713 — mutable test oracle emptied a recovery expectation

TEST/HARNESS; DIAGNOSED. Candidate rule suite correctly fired the recovery scenario at20m, while the fixture incorrectly expected no alert: the helper reread index0 after inserting a pending-state assertion there. No canonical promotion occurred. Preserve failed fixture/log; use an independent immutable expected-alert value and rerun the strengthened cases before canonical publication.

V374 FAIL706 recurrence2: a later read again guessed an absent SECURE_OPERATIONS_DELIVERY_PACK_PLAN path. No write or admission followed; actual JSON plan discovery found exactly ENTERPRISE_BACKEND, FRANCHISE_COMPLETE and FRANCHISE_SERVERLESS consumers. Use parsed pack IDs rather than naming conventions.

## FAIL-20260910-714 — candidate newline transport differs from canonical reconstruction

MATERIALIZATION/INTEGRITY; DIAGNOSED. Candidate rule/fixture were written using platform-default CRLF while canonical materialization correctly emits LF. Raw candidate-versus-rebuild byte assertion failed after38files reconstructed. Preserve candidate bytes; verify the difference is exclusively CRLF, then execute the exact LF reconstructed files and bind their SHA. Do not label the earlier CRLF hash as canonical or weaken block verification.

## FAIL-20260910-715 — nested PowerShell here-string collided with script transport

TOOLING/SERIALIZATION; RECOVERED. A Python generator containing a PowerShell here-string end marker terminated its outer shell transport. PowerShell rejected parsing before file writes. Deliver the script with apply_patch instead of nesting delimiters; verify emitted source and parse the resulting PowerShell before execution.

V374 successor205: FAIL712 REGRESSION_PROVEN on canonical SECUREOPS1.1.4: official-source engine evaluates9scenarios/13assertions PASS, low-traffic outage/reset failures reproduced before fix. FAIL713 recovered with independent expected alert; FAIL714 proves CRLF-only candidate delta and runs canonical LF; FAIL715 emitted PowerShell parses and focal wrapper passes5rejections plus honest absent-engine behavior. Full Preflight remains pending; FAIL385 and TEST02 remain open.

## FAIL-20260910-716 — source-engine report exposed machine-local paths

RELEASE/HYGIENE; DIAGNOSED. FullPreflight205 rejected the newly appended V374 report before structural completion because verbatim receipts/logs contained the local Windows user profile path. No PASS or release emitted. Preserve original raw receipts outside the distributable tree and their hashes; publish explicitly normalized display paths, then rerun the unchanged release-hygiene gate. Runtime rule proof and exact executable digest remain unchanged.

## FAIL-20260910-717 — receipt parser matched alertname as name

FAIL-20260910-717 TOOLING/TEST_RECEIPT: unanchored name regex also captured alertname, rejecting a correct three-scenario baseline failure. The original log remains intact. Anchor the line label before classification; do not rerun or alter the oracle. Root ledger pending completion of active Preflight to preserve checkpoint stability.

RECOVERED: anchored line-label parser confirms exactly the three expected failing baseline scenarios. The source-engine log is unchanged. Canonical37to38file delta proves only the denominator expression changed;36existing files identical and one new fixture. No rule oracle or status modified for this harness repair.

V374 successor207: Preflight206162executed steps PASS, all55compositions verified and promtool identity/rule/fixture hashes retained. FAIL716 REGRESSION_PROVEN by the unchanged release-hygiene selector; FAIL712 rule fix and canonical regression validated in complete library context. Docker availability remains separate BLOCKED; TEST02/03/07 and FAIL385 remain open.

## FAIL-20260910-718 — canonical gap-map header still directed readers to V305

FRESHNESS/CONTINUITY; RECOVERED. Despite appended V306–V374 evidence, FRANCHISE_GAP_MAP began with Estado vigenteV305, a next create-quote action and a67/743current count. These stale lead-ins could trigger duplicate work and hide the V372 runtime closure. Preserve before hash, mark dated entries historical, set the current45/48and67/755header and credit exact alternate-owner evidence without promoting21core equivalence. This is a cursor correction, not a new control closure.

## FAIL-20260910-719 — continuation guessed a source-lock path

TOOLING/READ; RECOVERED. Read-only rg used absent markdown_system/OFFICIAL_UPSTREAM_SOURCE_LOCK.md. No source or admission was inferred. Reuse the observed materialized acquisition lock and enumerate filenames before reading; this repeats FAIL706 and does not justify another source search.

V375 FAIL719 recurrence2: read-only request guessed internal/semconv/httpconv.go in the exact module; that filename is absent. No semantic result inferred. All following source reads use the ZIP/path inventory rather than prior-version naming.

## FAIL-20260910-720 — copied Go cache retained read-only manifest flags

TOOLING/ENVIRONMENT; RECOVERED. Source ZIP parity succeeded before candidate preparation; copytree retained Go-cache read-only flags, so changing the staging-only go.mod was rejected before tests. Original module/cache unchanged. Make only staged go.mod/go.sum writable, record dependency-only delta and run official suites. Future test copies use copyfile semantics or explicit staged manifest attributes.

V375 FAIL719 recurrence3: read-only request guessed a foundation migration filename. It was absent; no migration skipped or inferred. The integration harness will consume the exact sorted migration set already bound in the V372 runtime lock, avoiding name assumptions.

## FAIL-20260910-721 — offline graph audit lacked two metadata records

ENVIRONMENT/DEPENDENCY; DIAGNOSED. Reference tidy/test/vet/build passed using existing cached source artifacts. Complete go list -m -json all failed because the offline mirrors lacked .info for protobuf1.5.0 and reflect2 1.0.2. No full-graph scan is claimed. Retrieve only those exact official proxy metadata documents, preserve receipts, then rerun the unchanged graph command; no dependency version update or new code artifact acquisition.

## FAIL-20260910-722 — HTTP reference initdb failed with empty output

ENVIRONMENT/INTEGRATION; DIAGNOSED. Exact admitted PostgreSQL initdb returned nonzero before launching servers, with empty captured output. The rejected fresh reference directory/result remain; no business requests or alert proof occurred. Compare the launch environment with V372 and record exact exit codes before retrying a new owned cluster. Do not weaken database ownership or migration gates.

## FAIL-20260910-723 — standalone exporter suite retained monorepo replacements

ENVIRONMENT/UPSTREAM_TEST_LAYOUT; DIAGNOSED. Unchanged official exporter tests could not set up because its root go.mod contains six relative SDK replacements valid only in the upstream monorepo. Runtime consumption ignores dependency replacements and the actual reference build/test had passed. Preserve failed suite/mod; remove only those six staging replacements so tests resolve the exact same locked SDK artifacts as the consumer. No source/test assertion is changed, and no PASS is inferred from setup failures.

## FAIL-20260910-724 — parity helper assumed the probe copied an unused API command

TEST/ASSEMBLY; RECOVERED. Current canonical backend reconstruction produced18files; parity helper found cmd/api/main.go absent from the probe consumer, which intentionally copied internal/db/go manifests for the separate reference command. No mismatch in executed source was shown. Complete the consumer with the unchanged canonical command after proving its bytes equal V372, then compare all18files. The HTTP evidence continues to claim the real handler/domain/repository under a synthetic reference main, not execution of the production OIDC main.

## FAIL-20260910-725 — two-minute cadence did not guarantee sub-one-percent legacy ratio

TEST/NUMERICAL_ORACLE; DIAGNOSED. Actual canonical alert fired after599.883observed seconds pending on seven real HTTP500 responses, but the added legacy instantaneous-ratio assertion failed: three increments in the current five-minute rate window plus extrapolation produced0.010033243480063946. The average1/120 rate was not a bound on each sliding window. Preserve full758.823second failed run and actual firing prefix; recovery did not run. Use the original V374 four-minute sparse cadence (at most two increments per five-minute window), retain the unchanged1%/10minute/30second rule and all assertions, extend only the finite harness budget within the existing20minute host lifetime, then repeat complete failure/recovery. No test threshold relaxation or source mutation.
## FAIL-20260910-726 — reference rebuild digest differs

BUILD/PROVENANCE; DIAGNOSED. Eight candidate files reconstructed identically and tests/build passed, but the binary built over the minimal backend differs from the binary under runtime qualification. Preserve both artifacts. Compare complete imported package inputs and build metadata before accepting reproducibility; an eighteen-file backend subset does not establish parity with a composed consumer. No byte-identical receipt or G8 PASS emitted.

V375: FAIL721/722/723/725/726 recovered with exact metadata, COMSPEC environment, staging-only test replacements, actual four-minute sparse fault/alert/recovery and byte-identical complete-parent reconstruction. Original failures preserved; see HTTP_SLI_ALERT_INTEGRATION_V375.md. FAIL719 recurred in further guessed read-only compositor/catalog paths; no semantic evidence was inferred from absent paths.
## FAIL-20260910-727 — verifier insertion anchor matched two loops

TOOLING/EDIT; RECOVERED. The canonical update helper refused an ambiguous foreach anchor before changing verifier/owner files. Qualify the anchor with the loop body; preserve the failed invocation and existing validation logic. Runtime and source evidence are unchanged.
## FAIL-20260910-728 — new pack omitted the required apply-order heading

PACK/STRUCTURE; DIAGNOSED. Full Preflight210 rejects the new HTTP reference pack because its section8 is named Verification, not the canonical Apply order. Materializer and runtime tests do not enforce the whole documentation contract. Preserve the failed full-gate log; align all required sections with the canonical pack template without changing materialization blocks or weakening the structural validator, then repeat final checks.
## FAIL-20260910-729 — nested PowerShell expanded a status variable

TOOLING/TRANSPORT; RECOVERED. A nested shell wrapper expanded its exit-status expression before the inner parser, so validation never launched. Replace nested shell strings with subprocess argument arrays and direct log handles. No validation result inferred from this parse rejection.
## FAIL-20260910-730 — structural rerun preceded successor checkpoint

CONTINUITY/ORDER; RECOVERED. The unchanged structural gate accepted the repaired pack sections, then rejected the stale failures-owner hash because the ledger changed after checkpoint209. The checkpoint guard worked. Record the coherent successor before any verifier that validates local owners; do not disable hash verification or treat the interrupted suite as PASS.
## FAIL-20260910-731 — provenance summary omitted the eight authored reference files

PACK/INVENTORY; RECOVERED. Full Preflight211 correctly rejected the stale THIRD_PARTY_NOTICES current provenance total. Canonical census proves 1268 AUTHORED,140 ADAPTED,107 VERBATIM,total1515. Update the existing derived summary, preserve the failed log and repeat the unchanged gate after checkpoint. No provenance reclassification.

## FAIL-20260910-732 — pinned Next.js is affected by official Windows RCE advisory

UPSTREAM/SECURITY; DIAGNOSED. Current composed web pins Next.js16.3.2. The official maintainer advisory GHSA-p293-qw3h-jr36 / CVE-2026-75604 lists16.0 through less than16.3.3 as affected on Windows; no known workaround. Official advisories were inspected before selecting any successor. Block web promotion and TEST03/07, preserve baseline/locks and current HTTP reference's Go-only runtime scope. Qualify an official patched release with exact artifacts, graph/license/SCA, build and connected browser regressions before canonical update. Related AVIF advisory GHSA-2xp9-vwfh-vxw4 also requires range/closure review. Source: https://github.com/vercel/next.js/security/advisories/GHSA-p293-qw3h-jr36 .
## FAIL-20260910-733 — broad receipt count included unrelated profile fields

TEST_RECEIPT; RECOVERED. Final report helper counted every profile= substring rather than only VERIFY_LIBRARY result lines and rejected the valid log before writing the report. This repeats FAIL717 parser anchoring. Use the exact beginning/end-anchored profile result shape; preserve the original Preflight212 log. Checkpoint212 recorded the observed runtime/164 PASS prefix; successor213 adds the verified report after parser repair. No test or source rerun is inferred.

V375 checkpoint212: fullPreflight212164executed steps PASS,56profiles composed;165/1515/823/56 and1268AUTHORED/140ADAPTED/107VERBATIM. HTTP reference508requests,four actual failures,canonical ten-minute alert and recovery,one durable order/idempotency/outbox,exact rebuilt executable;380official+2local tests PASS. Only Docker unavailable in that receipt.45/48 unchanged; TEST02/03/07 remain blocked. Next action: FAIL732 official Next/Sharp patched-release qualification, preserving TEST02 owner integration work.
## FAIL-20260910-734 — registry HEAD omitted Content-Length

UPSTREAM_METADATA/TOOLING; RECOVERED. Exact npm metadata was preserved, but artifact-size discovery assumed every HEAD response contains Content-Length. The registry omitted it; no archive was installed or executed. Treat the header as optional, retain provenance metadata and probe a bounded one-byte HTTP range for total size before constructing an exact acquisition lock.
## FAIL-20260910-735 — published env package omits a standalone license file

UPSTREAM_PACKAGING/LICENSE; DIAGNOSED. Next16.3.4 archive passed exact integrity/path inventory (8527 files); the separate official @next/env artifact has no file named LICENSE/NOTICE, so the artifact-local license assertion failed. Preserve inspected bytes. Bind its declared MIT license to the already preserved exact verified Next repository root license and require that companion notice for redistribution; never invent an artifact-local file or equate an npm label alone with license evidence. Repeat bounded inventory with this explicit distinction and retain native license gates.
## FAIL-20260910-736 — inspection budget was below the declared SWC artifact size

TOOLING/BUDGET; RECOVERED. The fixed100MB per-file inspection cap rejected the official SWC native file of106190336bytes; exact registry metadata declares106190853unpacked bytes across three files and the signed archive already matches its digest. Preserve the rejected partial stage. Use the exact declared unpacked size as the per-file ceiling, still bounded by300MB total, while retaining path/type/count/hash checks. This is a declared resource-budget correction, not a security finding waiver.

V376 FAIL719 recurrence: a read guessed an absent web-pack filename and rg received unexpanded wildcard paths. No file changed or fact inferred; actual owner resolved from packId and glob inventory. Prepared changes use exact38/8/46file manifests and parsed30plan dependencies.

V376: Next16.3.4/Sharp0.35.4 replaces the affected Next16.3.2/Sharp0.35.3 baseline in TS-GO-API-WEB-BRIDGE0.5.5. Five exact public artifacts, ten verified provenance attestations, 8589 installed files identical to inspected archives,152 declared npm versions/0OSV findings versus3baseline advisories,115web tests plus1explicit integration skip, typecheck/build,92connected browser phases and4agenda projects with53exact migrations PASS. Multitab/session gates enabled on all four browser projects. Native AVIF/PNG/WebP smoke confirms libheif1.23.2. This is scoped advisory remediation; full vendored/native security and redistribution obligations remain blocked. Acquisition core0.4.84/46files,193IDs/22profiles and all four transport suites PASS;30dependent source plans updated. HTTP reference0.1.1 records the changed web-parent plan; Go/SQL/runtime evidence is reused only with byte parity.45/48 unchanged. See reconstruction_evidence/NEXT_SHARP_SECURITY_REPAIR_V376.md. Final canonical composition, HTTP binary parity and full Preflight remain pending in this entry.

FAIL732 scoped regression proven by the official fixed version ranges and exact candidate graph/runtime evidence; whole native/vendor admission remains blocked. FAIL735 next/env/SWC companion MIT root terms are bound explicitly; complete Sharp/libvips redistribution license/source obligations remain open. FAIL734 and736 inspection helpers recovered without weakening identity/path/type/hash checks.

## FAIL-20260910-737 — generated Next declaration differed from fresh source

FAIL-20260910-737 / LIB-FAIL-2453: final source parity stopped at next-env.d.ts. The connected candidate passed Next build, which inserted two generated type imports, while the fresh canonical materialization preserves the pre-build file. The failed strict assertion is preserved; no full parent parity receipt exists yet. Inspect the exact signed Next generator, bind the two-line transformation explicitly, and distinguish 754 identical source files from the generated declaration. No source or runtime predicate is waived. Root ledger append waits for the already-running Preflight214 to finish, preserving the validated root checkpoint during that read-only gate.


REGRESSION_PROVEN: exact signed Next writeAppTypeDeclarations generator adds only routes.d.ts and root-params.d.ts imports. All754other parent files match the connected candidate byte-for-byte. Fresh68/763composition differs from V375 only in package.json,pnpm-lock.yaml and the reference source-lock. HTTP tests/vet/build PASS and executable SHA4643654625c71aaf5b45367bd3a198381257588d073ae0349fc90c1bc1b5ab85 matches the real V375 run. No oracle weakened; rejected helper retained. Preflight214 completed before this root ledger update.

V376 FAIL733 recurrence: anchored final receipt parser used the unobserved field files= instead of actual implementation_files=. It rejected before root writes. Corrected to the exact observed line shape;56unique compositions classified from the unchanged successful Preflight214 log. No suite rerun or altered result.

V376 closure: Preflight214 passed all164executed steps and all56compositions; only Docker unavailable. Canonical165packs/1521files/824Markdown/56plans,1269AUTHORED/145ADAPTED/107VERBATIM; ordinary67/755 and HTTP68/763. Exact Next/Sharp package-advisory repair, signed artifacts and92connected phases/4agendas complete. Current45/48; TEST02/03/07 and native/vendor license/security gates remain blocked. Final canonical consumer has754parent source files byte-identical to the connected candidate plus the two proven Next-generated type imports. HTTP executable matches V375 byte-for-byte, so its unchanged15minute runtime evidence is reused without repetition. Next work returns to TEST02 owner integration, including the public-web SEO/i18n boundary; do not repeat completed browser creation/recovery work or the fixed dependency qualification.

V377 integrates actual Next canonical/robots/sitemap/WebSite JSON-LD under the existing BFF owner0.5.6/43files:5new AUTHORED files and3page changes; no dependency update. Public indexing is opt-in, limited to the real home/model-list pages and public render assets; private pages default noindex.137web tests/1preserved explicit integration skip, strict typecheck/build and3actual HTTP configurations PASS. The same build proves disabled/enabled/changed origin,Host poisoning isolation,mandatory cache revalidation,private noindex and JSON-LD HTML escaping/CSP nonce. Connected qualification: three browser projects passed initially, Firefox failed teardown with known FAIL423/742; a separate exact Firefox23phase rerun plus4agenda projects passed. These results cover4projects/92distinct phases across two runs, not a flawless92phase initial run or a Firefox lifecycle fix. Known upstream incident stays open.45/48 unchanged. See reconstruction_evidence/PUBLIC_WEB_METADATA_INTEGRATION_V377.md; final canonical composition and Preflight pending.

FAIL738/739/740 recovered in the executed final sources/probes; failed logs/stages preserved. FAIL741 corrected in current dependency bill and current inventory lead-ins only. FAIL742/423 remains open despite focused Firefox rerun PASS. No product assertions suppressed. The empty prepare directories were reused only after proving no files had been written by the rejected heading assertion.

V377 cierre217: integración pública canonical/robots/sitemap/JSON-LD en BFF0.5.6/43archivos comprobada con137tests,3configuraciones HTTP reales y composición final67/760. Regresión conectada:69fases de tres perfiles PASS en el primer run y23Firefox+4agendas PASS en una repetición aislada; primer fallo de teardown conservado y FAIL423/742 abierto. Preflight216:164pasos ejecutados PASS,56planes compuestos;165packs/1526files/825Markdown/56planes,1274AUTHORED/145ADAPTED/107VERBATIM. Docker sigue ausente. Binario HTTP idéntico a V375; fuentes Go/SQL/regla/dependencias sin cambios. Los48criterios/status son idénticos:45/48, TEST02/03/07 bloqueados. Esta frontera SEO de referencia está integrada; faltan i18n/localización y otras equivalencias funcionales, seguridad/admisión nativa y release integral. No contar el replay Firefox como reparación upstream ni repetir creación/recuperación ya probadas.

V378 recurrence FAIL719: two premature reads guessed connected/appointment-form.tsx and browser_gate/tests/enterprise-web.spec.ts. Both absent, no mutation. Enumerated actual locations/appointment-form.tsx and materialized manifests before subsequent use. Inspection filter also used uppercase PASS although contract stores passed; corrected to exact blocked status, no acceptance change.

### FAIL-20260910-744 — local UI state shadowed translation function

- Status: DIAGNOSED. V378 candidate only; canonical source unchanged.
- Strict typecheck rejected message(...) because an existing string state named message shadowed the imported function in both public forms. Unit catalog tests had passed; build was not executed.
- Correction: rename presentation state to statusMessage; preserve explicit localized invalid-receipt feedback via guarded returns and use a fixed localized message for unknown transport exceptions. Do not change TypeScript strictness, request payload or retry key.
- Evidence: V378/typecheck.log; run corrected candidate tests/typecheck/build and actual connected public flow before promotion.

### FAIL-20260910-743 — public locale and appointment time ignored configuration

- Status: REGRESSION_PROVEN for focused HTTP boundary; connected qualification pending.
- Baseline V377 actual Next build with en-US/UTC labels HTML en-US but renders Spanish forms and12:00 for15:00Z, because appointment formatting hardcodes Buenos Aires.
- Candidate resolves the existing trusted configuration, supplies complete es/en catalogs and explicit market timeZone. Same HTTP probe verifies English labels,15:00UTC and correct lang; preserves all raw HTML and exact before/after.
- No backend locale parsing, financial policy, source impersonation or product approval is added.

### FAIL-20260910-745 — cross-engine Intl spacing in test oracle

- Status: DIAGNOSED. First es-AR connected run: Chromium desktop/mobile and Firefox pass; WebKit differs only by U+00A0 inside p. m. Node's expected string contains U+0020.
- Browser trace/video/screenshot preserved before rerun. This is a string-oracle portability defect, not a wrong appointment instant. The first run does not claim final SQL verification.
- Correct only U+00A0/U+202F to spaces in the comparison; preserve date digits, punctuation and full text, and require submitted UTC instant equals the actual Go/PostgreSQL receipt. Baseline en-US/UTC red-green independently rejects the old three-hour offset.
- Re-execute all four projects against fresh owned PostgreSQL for es-AR/en-US/fr-FR fallback; keep original failed run and zero Playwright retries.

### FAIL-20260910-746 — Windows long path during exact installed-byte comparison

DIAGNOSED. pnpm nested dependency path exceeds legacy Win32 path length. Candidate comparison failed on a deep Next bundle-analyzer manifest; no missing-file waiver. Retry the same hash/size inventory using an extended absolute Windows path and require all8589files. Failed command retained as artifact-parity378.py.

V378 FAIL719 recurrence: guessed nonexistent compose_implementation_packs.ps1; corrected by reading VERIFY_LIBRARY and locating the existing generated tools/compose-markdown-project.ps1 before invoking.

V378 integra el recorrido público home→modelos→consulta→ubicaciones→solicitud de turno en español/inglés, con fallback español explícito y zona horaria del mercado.12/12recorridos navegador/configuración PASS,12leads y12turnos durables sin duplicados/sobrecupo; pérdida de respuesta y replay probados.164tests/1skip explícito conservado, typecheck/build PASS. El probe HTTP reproduce y corrige labels españoles bajo en-US y offset fijo de3horas bajo UTC. Cuatro archivos nuevos y ocho existentes modificados bajo owners BFF/portales/browser, sin nuevas dependencias ni reglas de dominio.45/48 permanece; no cerrar TEST02 por una sola frontera.

V378 follow-through: FAIL743/744/745/746 REGRESSION_PROVEN in final candidate and exact canonical reconstruction; failed source/logs/artifacts retained. Original Firefox lifecycle incident remains open. No broader control promotion.

V378 cierre219: recorrido público localizado es/en/fallback y zona horaria configurada,12/12casos navegador/configuración,12leads y12turnos durables sin duplicados/sobrecupo;164tests y typecheck/build PASS. Composición final67/764 y HTTP68/772; Preflight218164pasos ejecutados/56planes PASS, Docker ausente. Sin nuevas dependencias;8589archivos autenticados idénticos y ejecutable HTTP idéntico a V375.45/48, TEST02/03/07 siguen bloqueados. La brecha pública de idiomas queda demostrada; portales privados y demás integraciones mantienen su alcance. Se conservan los fallos iniciales y FAIL423/742 abierto.

V379 conecta búsqueda/lectura de ayuda operativa por permiso y versión exacta: cotizaciones y estado WhatsApp reutilizan sus textos1.0.0.186tests/1skip explícito heredado, typecheck/build y4navegadores PASS sin reintentos; sesión vencida/adulterada, cambio de permiso, versión ausente, lecturas tardías y recuperación503/malformed comprobados. Sin operaciones ni nuevas dependencias. Sólo cierra esta frontera de referencia;45/48 y TEST02/03/07 bloqueados, CMS/capacitación integral y target pendientes.

V379 continuity: no failed qualification command in this boundary. Preserve the inherited one-test skip and upstream incidents; the new browser negatives are expected controls, not incidents. Reuse the exact admitted session/toolchains and test the changed read boundary without starting unrelated databases/providers.

V379 cierre221: búsqueda/lectura de ayuda de cotizaciones y WhatsApp por permiso y versión exacta, compartiendo textos inline;186tests, typecheck/build y4/4navegadores PASS sin reintentos. Composición final67/774 y HTTP68/782; Preflight220164pasos ejecutados/56planes PASS, Docker ausente.8589archivos autenticados y ejecutable HTTP idénticos a evidencia previa. Esta frontera de lectura queda cerrada en referencia;45/48 permanece con TEST02/03/07 bloqueados. CMS, capacitación integral, otras guías y target conservan sus pendientes. No cambios de dependencia, reglas comerciales ni operaciones.

## FAIL-20260910-747 — V380 inline guide render parity

Status: DIAGNOSED_PENDING. Classification: contract/oracle. The first V380 candidate test run exited1:13inline guide digest checks, one mismatch;213other tests passed and the existing connected skip remains. Preserve tests.log and candidate before changing the oracle or content. Invariant: extracting versioned operational guidance must preserve its rendered text/version. Compare the exact failing original JSX and actual rendered HTML, correct the responsible source or serialization oracle, rerun the focused regression and dependent gates, reconstruct canonical blocks before promotion. tests.log SHA-256 8694de1958c91b981dbb729c78e4d70bbf3025e3367ffa5ac672569b3a58d2b0

FAIL747 diagnosis: the sole mismatch is JSX className="card" versus rendered HTML class="card". Replacing only that syntax in the captured original reproduces the actual digest exactly; all paragraphs, version and class value remain byte-identical. Correct only the captured HTML oracle; no production source change. original/red evidence retained; focused and dependent reruns pending.

## FAIL-20260910-748 — V380 strict typecheck

Status: OPEN. Candidate full tests214pass plus13focal render checks pass, but strict typecheck exited1; build did not execute. Preserve typecheck-fixed.log SHA-256 4712112eadcb344a9061437f70bac7d48b22febdc510caa2eccbf3cc9b8f3033. Diagnose the exact type error and fix source without weakening compiler settings, then rerun affected gates.

FAIL719 recurrence in V380: attempted read of build-fixed.log in the same shell batch as phase status without first checking that build had executed. File absent because typecheck failed. No data mutation. Read structured phase results before phase-dependent artifacts; original shell exit1 preserved in turn history.

FAIL748 diagnosis: the new render test supplied className:undefined despite exactOptionalPropertyTypes. Omit the optional prop when absent, preserving the production component contract and strict compiler settings. This is a test-call fix; rerun the13render oracles, full214unit tests, typecheck and build.

V380 completa el índice de las15/15guías versionadas existentes:13extracciones exactas, permisos por pantalla/sección y enlaces fijados.214tests,13regresiones de HTML inline, typecheck/build y8casos en4navegadores PASS;15perfiles ×15guías ×4navegadores=900decisiones de acceso, más60navegaciones exactas. Tres archivos de operación conservan todos sus handlers idénticos fuera de la ayuda. Sin nuevas dependencias ni reglas.45/48 permanece: no equivale a CMS, guías para tareas que antes carecían de ellas, capacitación integral ni target.

V380 follow-through: FAIL747 and748 REGRESSION_PROVEN in focused/full gates and exact canonical reconstruction; original failed logs/oracle retained. FAIL719 recurrence corrected by sequential phase receipt inspection. FAIL423/742 remains open. No global control promotion.

V380 cierre223: índice15/15de las guías versionadas existentes, permisos por pantalla/sección y enlaces exactos;13extracciones preservan texto/versión/HTML y handlers operativos intactos.214tests, typecheck/build y8casos en4navegadores PASS:900decisiones de acceso y60enlaces UI, sin reintentos. Composición67/777 y HTTP68/785; Preflight222164pasos ejecutados/56planes PASS, Docker ausente. FAIL747/748 corregidos con evidencia;45/48 mantiene TEST02/03/07 bloqueados. El inventario de guías existentes queda cerrado; CMS, tareas sin guía previa, capacitación integral, seguridad/release y target permanecen separados.

## FAIL-20260910-749 — multi-role candidate contract contradiction

Status: OPEN. TS-MULTIROLE-ONBOARDING0.1.0, outside the active composition, promises common authenticated guide access but visibleSections treats section.permission=* as requiring wildcard in the session. Its own tests expect the opposite. Role guide strings also suggest high-value approval flows and other operations without selecting them by current permission/implemented evidence. Preserve the exact original pack; qualify an isolated candidate, replace unsupported instructions with links to current authorized journeys/guides and prove that choosing owner/admin labels grants no access. No admission by compilation.

FAIL719 recurrence V381: rg was given a Windows path containing *PACK_PLAN.md instead of a directory plus -g; it returned os error123. Reissued the query with the documented directory/filter form. No source mutation from the failed read.

## FAIL-20260910-750 — browser generator escape warning

Status: REGRESSION_PROVEN for generator compilation. Python warned about invalid slash/question-mark escapes in embedded JavaScript regex syntax. Preserve the original generator, double only unescaped slash/question-mark backslashes, and compile with SyntaxWarning treated as error: PASS. Generated product/browser files are unchanged during the active test. No compiler setting was relaxed.

FAIL719 recurrence V381: a browser log was requested before proxy-build had completed; absent file, no mutation. Enumerated the running work directory and consumed phase output before subsequent dependent reads.

## FAIL-20260910-751 — first workspace browser qualification

Status: OPEN. First enabled-workspace run reports failures in the new workspace scenario on Chromium desktop/mobile while the existing help-recovery scenario passes. Preserve log, screenshots, traces and test source; inspect the exact failed assertion and fixture/application state before repair. Do not retry the unrelated green help lane or relax an authorization oracle without a diagnosed cause.

FAIL751 diagnosis: the browser reached the login URL, but the harness supplied plaintext where Next client navigation expects its own response format. The test now records and aborts the actual local login request, checking its exact return_to, instead of claiming an artificial login page rendered. No real IdP contacted. First four help-recovery cases passed; repeat only the four changed workspace cases with a fresh owned Next/proxy instance. All initial traces/videos/screenshots preserved in browser-first-artifacts. The unused customer backend-success assertion is replaced by the actual pre-backend admin denial for a customer session; help navigation remains the exercised successful downstream read.

## FAIL-20260910-752 — redirected request interception assumption

Status: DIAGNOSED. Second workspace browser run failed 4/4 at loginRequests=0. Trace proves actual dashboard HTTP 307 and exact Location /api/auth/login?return_to=/dashboard; followed request reaches local login and fails closed because OIDC_ISSUER is absent. Playwright route interception did not catch the redirect follow-up. Preserve second artifacts. Correct oracle to request the real protected route with maxRedirects:0 and require HTTP 307 plus exact Location, using the same browser context cookies. No real IdP, no synthetic success page, no production code change. Browser matrix remains required independently.

## FAIL-20260911-753 — manifest newline materialization

Status: DIAGNOSED. FAIL-20260911-753: canonical browser package.json digest mismatched. pnpm installation left19CRLF newlines; the Markdown materializer canonicalizes LF. Preserve original package bytes and pending pack, prove JSON object equivalence, normalize only that manifest and recompute block/pack digests. No script/dependency/value change and no test relaxation.

V381 integra el pack de panel por rol que estaba fuera del perfil y fallaba2de5tests: navegación compartida con BFF, permisos efectivos y cuatro vistas sin concesión por etiqueta.22tests de rol,236web/1skip previo, typecheck/build y8casos en4navegadores PASS;15perfiles ×5páginas ×4navegadores=300estados de página, además de flag desactivado, sesión inválida,404 y acceso admin denegado. Sin reglas financieras nuevas, operaciones, datos privados ni entrenamiento.45/48 permanece; cierre acotado del panel, no capacitación o integración integral.

V381 follow-through: FAIL749/751/752/753 REGRESSION_PROVEN by original red evidence, corrected real HTTP/browser oracles, JSON identity and canonical rebuild; FAIL750 retains its warning-free compile evidence. FAIL719 also recurred on a guessed nonexistent catalog path; use rg --files before named-file reads. No global promotion and FAIL423/742 remains open.

V381 cierre225: panel y cuatro vistas por rol integrados con navegación única BFF y permisos efectivos; pack antes fuera del perfil y2tests fallidos ahora reparado.22tests de rol,236web/1skip previo, typecheck/build,8casos en4navegadores y300estados de página PASS sin reintentos. Composición68/785 y HTTP69/793,4packs reconstruidos,8589artefactos idénticos y binario HTTP sin cambios. Preflight224:164pasos ejecutados/56planes PASS; Docker ausente. Panel cerrado en referencia, opt-in, sin conceder permisos por etiqueta ni inventar reglas/entrenamiento.45/48 mantiene TEST02/03/07 pendientes.

## FAIL-20260911-754 — administrative page combines different read grants

Status: OPEN. AdminPage requires admin:read but unconditionally fetches /v1/franchise/leads; the actual Go route requires lead:read. The current Promise.all rejects the entire admin page for an otherwise valid admin-only identity. Reproduce against the existing page before changing the UI; preserve both HTTP permission contracts, condition only the optional lead section, and prove real browser behavior.

## FAIL-20260911-755 — test suffix outside configured discovery

Status: DIAGNOSED. First V382 baseline exited1 because Vitest includes src/**/*.test.ts, while the new non-JSX test was named page.test.tsx. No test ran; this is not the product regression. Rename only the new test to page.test.ts, preserve the initial log, require actual9tests and the diagnosed two failures before repair. Do not broaden discovery or claim exit1 alone proves a defect.

## FAIL-20260911-756 — raw minor units rendered as currency amount

Status: DIAGNOSED. Visual review of the real V382 admin browser screenshot shows USD100 for a database total_minor_units=100; expected display is USD1,00. The same raw interpolation exists in customer orders and quote summaries. Reuse the exact BigInt/Intl presentation already implemented in the dedicated quote page; share its canonical owner, preserve validity behavior, and prove numeric labels plus unchanged stored values. No pricing/conversion/rounding-policy change.

V382 corrige dos defectos reales de lectura: admin:read ya no consulta leads sin lead:read y los portales dejan de mostrar unidades menores como importes mayores. Se comparte exactamente el formatter de cotizaciones, sin alterar precio/estado/aceptación.24regresiones focales,260tests web/1skip previo, typecheck/build y4navegadores contra Go/JWKS/PostgreSQL PASS;5negativos HTTP, consultas indebidas0, escrituras0 y snapshot de6tablas idéntico.45/48 permanece, sin promoción global.

V382 follow-through: FAIL754/755/756 REGRESSION_PROVEN. Preserve the no-tests run,2failed/7passed permission baseline, monetary red and both actual browser stages. Exact quote extraction and canonical parity recorded. No ACL/business-policy change and FAIL423/742 remains open.

V382 cierre227: admin:read funciona sin consultar leads no autorizados; admin/customer muestran importes correctos mediante el formatter exacto de cotizaciones compartido.260tests web/1skip previo,24focales,typecheck/build y4navegadores con Go/JWKS/PostgreSQL PASS;5negativos HTTP,0consultas indebidas,0escrituras y snapshot6tablas idéntico.6packs reconstruidos,composición68/790 y HTTP69/798;8589artefactos y binario HTTP idénticos. Preflight226164pasos ejecutados/56planes PASS; Docker ausente. FAIL754/755/756 corregidos con originales preservados.45/48 conserva TEST02/03/07 pendientes; no cerrar integración,seguridad o release completos por estos fixes.

V383 recurrence FAIL719: attempted nonexistent PROJECT_ENGINEERING_CONTRACT.md; actual owner is the existing PROJECT_ENGINEERING_CONTRACT.json. Read checkpoint refs before constructing an owner filename; no source mutation occurred.

## FAIL-20260911-757 — private read error has no safe recovery entry

Status: DIAGNOSED. V383 connected baseline against unchanged V382 production build preserves backend403 as a generic Next500 but lacks the expected recovery alert/link. Raw browser artifacts retained. Add fixed-text private-route error boundaries, no error payload display or automatic retry; prove explicit GET reconsultation and unchanged API permissions/domain snapshots.

## FAIL-20260911-758 — Next route announcer shares alert role

Status: DIAGNOSED. Four-browser first fixed run finds both the actual recovery paragraph and Next route announcer for global getByRole(alert). Preserve original artifacts; scope all three alert assertions to the existing main landmark, retain exact text/denial/recovery/count checks and all four browsers. No product behavior change.

V383 completa recuperación de lecturas admin/customer/factory: error sin detalles internos, rechazo persistente sin permisos y nueva consulta GET tras corregir sesión; ninguna repetición de operaciones.263tests web/1skip previo,27focales,typecheck/build y4navegadores Go/JWKS/PostgreSQL PASS;5negativos HTTP,0escrituras y snapshot6tablas idéntico. Cinco archivos AUTHORED nuevos; producción Go/SQL/dependencias intactas.45/48 mantiene TEST02/03/07 pendientes.

V383 follow-through: FAIL757/758 REGRESSION_PROVEN with original stages retained, exact canonical reconstruction and four browser projects. No upstream Firefox incident closed.

## FAIL-20260911-759 — inventory substitution corrupted baseline path

Status: DIAGNOSED. Generated parity helper globally replaced798→803 inside the previous stage UUID as well as the file count. Composition and795parent-file parity passed, but delta comparison treated all files as added against an absent path. Preserve failed helper/log. Resolve baseline from tools.json, assert its materialization record exists, and rerun every hash/delta assertion against the already reconstructed consumer. No product changes or gate weakening.

V383 FAIL759 REGRESSION_PROVEN: composition803and parent795byte checks had passed; a global count substitution corrupted the baseline UUID. Original helper/log preserved, baseline now resolved from recorded tools.json with required materialization record. All final added/modified/byte/hash assertions pass and the rebuilt HTTP executable is identical to V375. No source, permission or acceptance relaxation.

V383 cierre229: recuperación admin/customer/factory integrada, error fijo sin datos internos y GET explícito; bearer insuficiente permanece denegado y sesión válida recupera datos.263tests/1skip previo,27focales,typecheck/build y4navegadores Go/JWKS/PostgreSQL PASS;0escrituras/snapshot6tablas idéntico.5packs reconstruidos,68/795 y HTTP69/803,8589artefactos/binario HTTP idénticos. Preflight228164pasos/56planes PASS, Docker ausente.45/48; TEST02/03/07 pendientes, sin cambiar oráculos.

V383 continuation FAIL735/532: exact Next MIT companion captured and byte-matched to the verified repository and three signed packages. Sharp native29notice rows/28versions reconciled; libnsgif owner identified as13vendored blobs at verified libvips426af3f44246fce9cfa8dd51a353aa4dfd48c553, four document/blob matches. No fabricated NetSurf version, binary correspondence, distribution or full native/SCA admission. Raw receipts in V383 notice-analysis; production blocks retained.

V384 recurrence FAIL719: attempted node_modules/sharp/package.json although Sharp is only a transitive pnpm package. No source mutation. Resolve every package root from authenticated inventory as already implemented in notice-analysis383.py.

V384 FAIL719 recurrence: looked for the new web profile in historical source-core; inventory identified source-rebuilt and source-candidate as its actual owners. Existing canonical46file source baseline is now explicit and every candidate block/hash will be compared.

## FAIL-20260911-760 — source-tooling profile mistaken for runtime parent

Status: DIAGNOSED. V384 pending canonical helper assumed the acquisition core was selected by FRANCHISE_COMPLETE_PACK_PLAN and attempted to update its missing pending plan. Actual plan excludes this tooling owner. No canonical files were published by that failed preparation. Preserve original pending files/helper, derive affected profiles only from their actual selections, and leave franchise/HTTP source lock/version/file counts unchanged. Only30source plans and6explicit count assertions require+2; verify all compositions before closure.

V384 demuestra correspondencia exacta del insumo nativo Sharp: libvips-42.dll18614784bytes, versions.json y avisos coinciden con @img/sharp-libvips-win32-x641.3.3 firmado. Perfil previo/2attestations/tamper, cuarentena sin ejecución y4suites del transporte PASS; core0.4.85/48files/194IDs/23perfiles. DLL C++ y build MXE/notices/SCA completos pendientes. Runtime y perfil68/795 intactos;45/48 sin promoción.

FAIL760 corrected from actual source profile selections; original unpublished preparation retained. Runtime parent unchanged. Preserve moving-notice input and C++/MXE/SCA/source obligations; no substituted or executed native runtime.

V384 FAIL760 REGRESSION_PROVEN for exact source/runtime scope:48canonical source files match the tested candidate, and803fresh runtime consumer files equal V383. Original incorrect preparation remains preserved. FAIL719 also recurred with guessed canonical helper and Windows wildcard-path lookup; current paths were resolved from directory inventory before use.

V384 cierre231: input nativo Sharp1.3.3 firmado y adquirido con perfil; libvips-42.dll18614784bytes, versiones y avisos byte-idénticos. Core0.4.85/48files/194IDs/23perfiles y4suites PASS;30planes actualizados. Preflight230164pasos/56composiciones PASS, Docker ausente. Runtime68/795 y69/803 idéntico aV383.45/48; integración completa, seguridad nativa y release pendientes.

V385 FAIL719 recurrence: an unenumerated nested scripts/acquire_opaque_artifact.ps1 path was absent. Enumerate source-final before reads; canonical transport is elite_sources/opaque_artifact_transport.ps1. No artifact or runtime changed.

## FAIL-20260911-761 — NATIVE_BUILD_INPUT_PROVENANCE

Status: DIAGNOSED. Classification: DEPENDENCY/LICENSE. V385: exact ZIP/DLL/versions/LICENSE correspondence PASS, but release ZIP has no complete linked-component source set. Pinned MXE recipe uses mutable base container and llvm-mingw-20260605 branch; source build enables GLib fallback downloads and librsvg Cargo dependencies. Observed branch commit d973945bb92c7783d5afa41bb2b8d2e1a04eaba3 is fixed for inspection, not proven historical compiler input. 28 component versions and SHA256 source checksums now match recipes. Preserve component graph/source/notice/SCA and historical toolchain evidence as TEST03 blockers; no inferred source-complete admission. Evidence: V385 recipe-receipts.json, native-origin-correspondence.json and component-recipe-map.json. Next: inspect exact source archives and vendored graphs, recover build attestations or qualify an explicitly adapted rebuilt native artifact with target gates.

## FAIL-20260911-762 — UPSTREAM_SOURCE_LOCATION_DRIFT

Status: REGRESSION_PROVEN for source-location metadata only. V385 pinned base src/libxml2.mk URL uses directory2.14 for2.15.3 and returns404. Official GNOME2.15 index exposes the same source, HEAD200/3152452bytes and official checksum78262a6e7ac170d6528ebfe2efccdf220191a5af6a6cd61ea4a9a9a5042c7a07 equal to recipe. No source archive downloaded or different version substituted. Preserve original URL/error; canonical V385 evidence and continuity point to corrected official2.15URL. Browser could not open news/hash endpoints; direct bounded official HTTPS supplied their bytes. Metadata recorded outside root while Preflight232 ran, then flushed after all164stepsPASS.

## FAIL-20260911-763 — LIBXML2_NATIVE_SECURITY_REVIEW

Status: DIAGNOSED; admission remains BLOCKED. Official libxml22.15.4news published September lists security fixes; installed authenticated Sharp0.35.4/MXE8.18.6 versions declare2.15.3. Base and release recipe patches only adjust pkg-config/DllMain, not those security functions. Official latest package/release metadata still Sharp0.35.4, nativeinput1.3.3, MXE/libvips8.18.6; no newer fixed packaged native release found. Do not infer exact exploit reachability or backport absence beyond inspected recipes. Need exact official fix/source graph, confirmed binary/source reachability and qualified adapted rebuild or equivalent containment. TEST03/07 stayblocked; no runtime modified. GitLab encoded project API returned404 for both tags/compare; official GNOME release news/hash bytes retained, alternate official source discovery continues. Evidence native-security-discovery.json, xml2-security-official-review.json, both fixed recipe patches and source-location-discovery/libxml2-2.15.4.news.

V385 cierre233: origen ZIP-DLL-version-license exacto,488files inspeccionados y28versiones/SHA256 de fuentes fijados; core0.4.86/50files/195IDs/24perfiles,4suites PASS. Preflight232164pasos/56composiciones PASS, Docker ausente. Runtime68/795 y69/803 idéntico aV384.45/48. FAIL763: libxml22.15.4 corrige seguridad frente al2.15.3deSharp; fix oficial fijado, sin release Sharp/MXE posterior observado ni reparación binaria demostrada.

FAIL763 official fix c94eb021→96498992,39commits/47files captured through verified GNOME mirror. No repaired binary/runtime admission. Next work follows fixed source graph and controlled candidate rebuild/containment. FAIL762 original404 and exact checksum correction retained.

V386 FAIL719 recurrence: recursive search across unrelated Temp directories hit protected paths and was stopped. Subsequent reads use only inventoried V385 paths. No source/runtime mutation. WSLUbuntu explicitly inventoried: GCC15.2.0, autoconf/automake/libtoolize available; Docker absence does not imply no native build environment.

## FAIL-20260911-764 — INLINE_PYTHON_SHELL_QUOTING

Status: DIAGNOSED. A python -c helper-generation command had unclosed quoting; no helper or source was written. Correction: use apply_patch to create the explicit acquisition helper, avoiding nested shell/Python quoting. Regression: execute that helper with profile/hash checks and retain source validation/acquisition receipts.

## FAIL-20260911-765 — LIBXML2_FIXED_BUILD_TOOLING

Status: DIAGNOSED. Fixed2.15.4 configure requires pkg-config; old2.15.3 configures without it. Initial fixed configure exit1 retained under V386/2.15.4-configure.log and owned WSL config.log. Do not stub or bypass the dependency check. Locate an existing exact pkgconf/pkg-config or admit an official bounded tool package, record provenance, then rebuild both versions from clean copies.

V386 FAIL719 recurrence: initial make check-build listed removed target testAutomata. Source Makefile.am now supplies actual check_PROGRAMS; no test removed or gate skipped. Failed output preserved; clean rebuild will invoke upstream make check directly.

## FAIL-20260911-766 — LIBXML2_STATIC_ONLY_MODULE_TEST

Status: DIAGNOSED. On recovery, both retained rebuild1 make check logs exit2 at testModule: Failed to open module. The qualification configured --disable-shared while the unchanged upstream test opens .libs/testdso.so. Source testModule.c and Makefile.am define this dependency. Preserve both failures; clean rebuild2 enables shared and static artifacts with the same ASan/UBSan flags and executes the full upstream make check, including module loading. No disabled test, runtime replacement or Windows qualification inferred.

V386 recovery: checkpoint233 state and all233events verified; interrupted local source/dependency/failure records identified by before/after SHA. FAIL719 recurred reading guessed MATERIALIZE_IMPLEMENTATION_PACK.ps1; actual inventoried owner is materialize_markdown_pack.ps1. Recovered both failed module checks as FAIL766; no product changed.

V386: usuario informa restriccion de contenido de ciberseguridad con acceso Daybreak y pide registrar/bypassear en orden de trabajo, no el control. Investigacion libxml2 suspendida ACCESS_BLOCKED; TEST03/07 siguen BLOCKED. Ambos make check completos ASan/UBSan Linux PASS; ninguna DLL Windows reparada o promovida. Siguiente TEST02: integracion funcional pendiente de FAIL385, sin repetir auditoria21yahecha.45/48.

FAIL764/765/766 REGRESSION_PROVEN for acquisition/tooling/module-suite scope only. FAIL763 remains ACCESS_BLOCKED; user reported platform restriction and directed proceeding to independent TEST02. Evidence reconstruction_evidence/NATIVE_QUALIFICATION_ACCESS_BLOCK_V386.md.

## FAIL-20260911-768 — CANDIDATE_DEPENDENCY_COPY_LONG_PATH

Status: DIAGNOSED. shutil.copytree of installed pnpm tree failed on Windows paths over260characters, leaving a partial candidate; subsequent fix refused the missing baseline.json before mutation. Preserve partial tree and traceback. Recreate from803canonical files and install the unchanged lock offline through pnpm instead of copying installed links. Validate artifact8589hashes before use. FAIL719 recurrence on guessed prepare382.py is retained; source paths now come from tools.json/manifest.

## FAIL-20260911-767 — PRIVATE_PORTAL_PAGINATION_MISSING

Status: DIAGNOSED. Customer orders/service cases and factory units request limit25 and ignore backend next_cursor, hiding later records. Two baseline rendering tests fail against unchanged canonical pages. Repair navigation using existing after query and session scope, retain fixed page size and independent cursors; prove all records against disposable PostgreSQL and no writes. No business policy or dependency change. V387 pagination-red.log preserved. FAIL719 also recurred on an irrelevant guessed i18n/public.ts path; use actual manifest inventory.

## FAIL-20260911-769 — PAGINATION_FIXTURE_UNIQUE_IDENTITIES

Status: DIAGNOSED. Connected seed fails before browser at production_unit_tenant_id_vin_key; schema intentionally uses UNIQUE NULLS NOT DISTINCT for VIN and battery. The extra26synthetic units omitted both identifiers. Preserve failed database/log; assign each unit explicit unique synthetic VIN/battery and rebuild a fresh database. No constraint or assertion removed, no product policy changed.

V387: paginacion real en pedidos/casos cliente y unidades fabrica;3listas/81registros por navegador,4navegadores PASS,0escrituras/snapshot6tablas intacto.272tests/1skip heredado,typecheck/build;6packs reconstruidos y805files exactos,franquicia68/797. TEST02 sigue pendiente integral;TEST03libxml2 ACCESS_BLOCKED por indicacion del usuario,TEST07dependiente.45/48.
FAIL767/768/769 REGRESSION_PROVEN for the bounded pagination/environment/fixture claims in reconstruction_evidence/PRIVATE_PORTAL_PAGINATION_V387.md. Native security scope remains ACCESS_BLOCKED; no test status promotion.

## FAIL-20260911-770 — ADMIN_PAGINATION_MISSING

Status: DIAGNOSED. FAIL-20260911-770 ADMIN_PAGINATION_MISSING: admin orders/cases/leads ignore next_cursor; one baseline test fails with granted scoped fixture. Preserve independent lead permission and existing page size. Pending merge to root after Preflight236 stable snapshot.

## FAIL-20260911-771 — PREFLIGHT_INVENTORY_OWNER_NOT_UPDATED

Status: DIAGNOSED. Preflight236 rejected stale FRANCHISE_PREFLIGHT_GAP inventory after two new evidence documents/two source files. Expected165packs/1561files/835Markdown/56plans and797franchise files. Source reconstruction and checkpoint passed before this assertion. Preserve failed log; update that canonical summary without weakening the exact verifier, include it in all subsequent inventory updates and rerun full Preflight after the adjacent admin paging delta.

V388: seis listados cliente/fabrica/admin paginados con APIs existentes;163registros de lista por navegador,4navegadores PASS,0escrituras/snapshot6tablas intacto.276tests/1skip heredado,typecheck/build;6packs reconstruidos,805files exactos,68/797integral.45/48;TEST02integral pendiente,TEST03libxml2 ACCESS_BLOCKED,TEST07dependiente.
FAIL770 REGRESSION_PROVEN in reconstruction_evidence/ADMIN_PORTAL_PAGINATION_V388.md; FAIL771 full Preflight rerun pending. V386 native research remains ACCESS_BLOCKED; no promotion.

V388 cierre239: seis listados cliente/fabrica/admin conectados al cursor real;163registros de lista por navegador,4navegadores PASS,0escrituras/snapshot6tablas idéntico.276tests/1skip previo,typecheck/build;6packs y805files reconstruidos,HTTPbinario idéntico. Preflight238164pasos/56perfiles PASS; Docker ausente.165packs/1561files/836Markdown; integral68/797.45/48,TEST02integral pendiente,TEST03 ACCESS_BLOCKED por restriccion reportada/usuario,TEST07dependiente.
FAIL771 REGRESSION_PROVEN; original236 failure preserved, full238 verifies canonical inventory and both pagination repairs.
`nFAIL719 recurrence 2026-09-11, resume239: guessed src/platform/config/tenant.ts absent. Read-only lookup failed; inventory config directory before selecting its owner. No product change; original tool output retained.
FAIL719 same-cycle recurrence: guessed src/app/(public)/page.tsx absent; use rg --files output exclusively for subsequent reads.

## FAIL-20260911-772 — PRIVATE_APPOINTMENT_IMPLICIT_TIME_ZONE

Status: DIAGNOSED. Customer account renders appointment using server default; customer appointments client repeats host-default formatting. Both ignore admitted business locale/market timeZone already used in public booking. Baseline regression tests fail with explicit date-boundary/DST expectations; receipts preserved in V389 staging. Repair presentation using existing configuration and formatter, without changing stored UTC timestamps, cancellation authorization or business policy.
`n## FAIL-20260911-773 — TIME_ORACLE_HOUR_CYCLE_ASSUMPTION`n`nStatus: DIAGNOSED. V389 first green run has 14/15PASS: es-AR Intl formatter on pinned runtime emits 10:30 p. m. for 22:30; test incorrectly assumed a 24-hour locale default. Preserve failure log; correct literal oracle to equivalent 12-hour representation without changing formatter, timestamp or zone.

## FAIL-20260911-774 — TIME_ORACLE_ISO_LEXICAL_ASSUMPTION

Status: DIAGNOSED. Four first connected browsers correctly render two time nodes but test selects literal Z representation. Preserved traces show 2035-01-01T22:30:00-03:00 equals the seeded 2035-01-02T01:30:00Z. Product preserves original API string as intended; compare exact epoch instants and independently assert expected local date/hour/zone, without normalizing or changing product/backend. First failed database/trace/receipts preserved and owned processes stopped.

## FAIL-20260911-775 — WEBKIT_LOOPBACK_PREFETCH_ACCESS_ERRORS

Status: DIAGNOSED / fixture limitation retained. Three browsers pass all assertions; WebKit also preserves both times but newly added all-pageerror assertion sees four access-control errors from loopback Next prefetch URLs (/ , /dashboard, /help, /models with _rsc). This is not demonstrated to be a hydration failure, nor is its transport root cause repaired. Preserve exact failures/trace. Restrict the new assertion to the date-rendering claim: record every exception, allow only those four exact same-origin prefetch paths with the observed suffix, reject every other exception and React hydration error. Full clean-WebKit transport certification remains unclaimed; no production security or permission changes.

FAIL775 recurrence: NY first run has3PASS/1WebKit failure on customer quotes/handovers prefetch, with correct date assertions. Do not keep broadening an exception list. Candidate hypothesis: forced full navigations/reload cancel background prefetched requests in self-signed WebKit fixture. Remove the exception classifier entirely; restore zero page/hydration errors, explicitly await network quiescence before beginning the date segment and before its forced reload. Preserve all previous failures/classifier candidate; rerun both configurations.

V389 checkpoint241: horas de turnos cliente coherentes entre cuenta/gestion con locale/zona del negocio, labels SSR preservados al hidratar.282tests/1skip previo,typecheck/build;2configuraciones x4navegadores,2instantes con limite de dia/DST,snapshot7tablas intacto y0escrituras.6packs reconstruidos/805files,HTTPbinario idéntico;68/797.45/48; TEST02integral pendiente, V386native diferido,TEST07dependiente. Preflight241 pendiente.
FAIL772/773/774 REGRESSION_PROVEN in reconstruction_evidence/CUSTOMER_APPOINTMENT_TIMEZONE_V389.md. No native/security or whole control promotion.

V389 cierre242: turnos cliente muestran fecha/hora/zona del negocio coherentes en cuenta/gestion/recarga, sin depender del dispositivo.282tests/1skip previo,8runs navegador(2config x4),0errores en segmento de fechas,0escrituras/snapshot7tablas idéntico.6packs/805files reconstruidos,8589artefactos yHTTPbinario idénticos. Preflight241164pasos/56composiciones PASS; Docker ausente.165packs/1561files/837Markdown;68/797.45/48; TEST02integral pendiente, V386native diferido yTEST07dependiente.
FAIL772–774 repaired; FAIL775 qualified fixture sequencing only, not an upstream browser fix. Temporary error classifier rejected and absent from canonical code.

## FAIL-20260911-776 — CUSTOMER_CANCELLATION_AMBIGUOUS_RESULT

Status: DIAGNOSED. Current catch reports No se canceló and enables another action after losing a response, even if backend committed. Any HTTP200 JSON is accepted without binding result ID/state/version. Isolated browser regression will deliver actual cancellation then drop its response; preserve it before repairing. Use existing PostgreSQL version/CAS, audit and outbox owner; no new business policy.
FAIL719 recurrence V390: guessed quote-acceptance.spec.mjs absent; available test inventory shows enterprise-web.spec.mjs. Read-only path failure preserved; use only discovered filenames.
`n## FAIL-20260911-777 — CUSTOMER_CANCELLATION_DOUBLE_SUBMIT`n`nStatus: DIAGNOSED. V390 first baseline fails in all four browsers: same-event double click sends two POSTs. Backend CAS rejects second effect, but client lacks synchronous in-flight guard. Original source/logs preserved. A second isolated diagnostic uses one click only to reach the separate lost-response defect; final gate restores double click and its exact one-request requirement.

## FAIL-20260911-778 — CANCELLATION_TEST_RELOAD_RACE

Status: DIAGNOSED. First corrected connected run passes Firefox/WebKit; two Chromium projects read DOM while automatic success reload starts (Execution context was destroyed). Baseline single-click diagnostic had same fixture race before reaching lost-response check on those projects. Register DOMContentLoaded wait before each successful action and await that actual new document, then network quiescence, before inspecting. Preserve all business/result/concurrency oracles and original logs; no retries added.

V390 checkpoint244: cancelacion cliente conectada y recuperable: receipt ligado aID/org/estado/version,guard sincrono,GETexplicito tras resultado incierto.294tests/1skip previo,typecheck/build;4navegadores,16cancelaciones/16audits/16outbox con28POSTincluidos rechazos;8regresiones lectura/fechas sin escrituras.6packs/805files reconstruidos,8589artefactos yHTTPbinario idénticos;68/797.45/48;V386diferido,TEST02integral yTEST07pendientes. Preflight244 pendiente.
FAIL776/777/778 REGRESSION_PROVEN in reconstruction_evidence/CUSTOMER_CANCELLATION_RECOVERY_V390.md. No native/security or whole control promotion.

V390 cierre245: recorrido cliente cuenta→turnos→cancelacion→recuperacion GET cerrado en referencia.294tests/1skip previo,4runs cancelacion con16cambios/audits/outbox unicos y28POSTincluidos rechazos;8runs lectura/fechas siguen sin escrituras.6packs/805files reconstruidos,8589artefactos yHTTPbinario idénticos. Preflight244164pasos/56composiciones PASS; Docker ausente.165packs/1561files/838Markdown;68/797.45/48; TEST02integral pendiente, V386native diferido yTEST07dependiente.
FAIL776/777/778 proven with final canonical source and unchanged full verifier. No production/security or21core-equivalence claim.

FAIL719 recurrence V391: rg received an unexpanded Windows wildcard path, and readiness log references E-H were mistaken for existing advisory files. Both read-only failures are preserved in the tool transcript; no completion inferred. Resolve source inventory before reading and generate future advisories only through the canonical renderer. No product mutation.

FAIL719 V391 follow-through: evidence parsing assumed acceptance_tests instead of inventoried tests, and a map basename was used without its markdown_system directory. Read-only errors, no verdict inferred; keys and exact owner paths now inventoried. Original tool transcript retained.

### FAIL-20260911-779 — Windows locale in readiness reconciliation

V391 / ASSUMPTION / REGRESSION_PROVEN. Stage script omitted explicit UTF-8 when reading Spanish canonical Markdown. Exact-row assertion rejected the text before writing; dependent E invocation also rejected unchanged PENDING D. Both source files remain byte-identical to the snapshots. Preserve failed script/transcript; use explicit UTF-8 on every read and inspect each stage exit before executing dependent work. No answers or approval inferred.

V391 FAIL779 follow-through: readiness CLI correctly rejected --report outside project root; no report was fabricated or replaced. Qualification now captures stdout via subprocess to external evidence, with no --report argument, and validates exit2 plus JSON BLOCKED. Preserve rejected invocation in transcript.

V391 FAIL779: explicit UTF-8 edits complete D/E/F sequentially; validator confirms exactly42→41→40→39 observations, only their three round-status errors removed. All other input fields and48test definitions/statuses byte-semantically unchanged. Rejected outside-root report probe recovered using captured CLI output. FAIL719 wildcard recurrence during staging inventory retained; corrected subsequent reads use explicit paths.

V391 FAIL716 recurrence: release hygiene rejected a machine-local home path in the new report. Original report and structural.log preserved in stage; portable report now records only Temp/stage basename, same receipts and claims. No selector bypass; rerun unchanged structural gate.

V391 FAIL563 recurrence: structural verification ran before checkpointing the changed local evidence, correctly rejecting stale checkpoint hashes. Preserve structural-corrected.log; append successor246 with all real evidence hashes, then rerun without disabling checkpoint validation. No source failure or lost events.

V391 FAIL770 recurrence: inventory summary counted three local advisory files as portable Markdown. Canonical VERIFY_LIBRARY selector reports165/1561/839/56, since local advisories are excluded. Original expected842 and failed structural246.log preserved. Correct summary to839, keep selectors unchanged, checkpoint before final structural validation.

V391-V392: FAIL716/563/770 de edición y checkpoint recuperados: verificador estructural247 y ambos empaquetadores PASS sobre fuente congelada; aceptación desde ZIP PASS, originales conservados. No defecto de código/release cerrado por estos fixes administrativos.

V393 FAIL719 recurrence: source search included absent next.config.mjs; actual manifest supplies next.config.ts. No config/source change or coverage inferred from the rejected filename. One guessed staging lookup produced no data and is not evidence; subsequent staging locations come only from recorded pointers and inventories.

V393 FAIL770 recurrence: candidate source inventory counted MATERIALIZATION_RECORD.md as an implementation file, rejecting806 vs805 before edits. Preserve tool transcript; explicitly inventory805source files plus the composition receipt separately, which remains intact. No dependency or source modification in the rejected preparation.

V393 FAIL768 recurrence: remaining-install parity rejected a 270-character Windows path because ordinary Path.is_file returned false while extended-length Path.is_file returned true for the same existing file. Original artifacts393.py and failed transcript preserved. Apply extended-length paths to both comparison sides; retain every file/hash assertion. No missing-package or integrity failure inferred.

V393 FAIL768 follow-through: generated backslash literal caused a syntax rejection before comparison. Rejected script preserved; construct the Windows prefix from character codes, keeping comparison assertions intact.

V393 FAIL719 recurrence: a guessed completion-contract basename was absent. The real inventoried owner is PROJECT_ENGINEERING_CONTRACT.json. Read-only error retained; no evidence inferred.

V393 FAIL719 follow-through: guessed roadmap basename rejected; the specification resolves its canonical owner to root REUSABLE_CODE_READINESS_ROADMAP.md. No state inferred from absent path. FAIL768 parity recovery now proves12555remaining installed files byte-identical, including all long paths and8534signed artifact files.

V393 lock comparison rejected: handwritten section extraction included the new top-level ignoredOptionalDependencies setting in the last zod package stanza. Exact diff shows zod integrity unchanged. Preserve rejected parser; bound package section at both following top-level keys and retain exact equality for all122surviving stanzas. Not a dependency-byte failure.

V393 FAIL719 follow-through: guessed roadmap basename rejected; the specification resolves its canonical owner to root REUSABLE_CODE_READINESS_ROADMAP.md. No state inferred from absent path. FAIL768 parity recovery now proves12555remaining installed files byte-identical, including all long paths and8534signed artifact files.

V393 FAIL719 recurrence: two guessed acquisition/profile basenames were absent. File inventory resolves OFFICIAL_UPSTREAM_ACQUISITION_CORE.md and the existing root profile record; only the separately authorized native inspection profile selects the old Sharp artifact. Active composition plans do not select it. No source/profile mutation inferred from failed paths.

V393 FAIL780 / CANONICAL_NEWLINE_MISMATCH: first pack rebuild rejected next.config.ts because candidate writes used Windows CRLF while canonical reconstruction emits LF. Root packs unchanged; rejected candidate bytes, pending packs and logs retained. Normalize only the five changed text files to explicit UTF-8/LF, verify identical decoded text, regenerate hashes, and rerun affected web/build/connected gates before promotion.

V393 FAIL781 / LONG_PATH_ENUMERATION: fresh-install parity counted11fewer Next declaration files than the candidate, although the independent extended-path signed-artifact check found all8534files. Ordinary os.walk silently skipped directories beyond Windows path limits; this is an incomplete inventory, not missing package content. Preserve both inventories/scripts; enumerate from an extended-length root and compare all files without an exclusion or count adjustment. Root application bytes and existing runtime results remain unchanged.

V393 FAIL781 REGRESSION_PROVEN: extended-root enumeration gives12571files in both candidate and fresh canonical install, complete maps identical and every hash matches the earlier qualified installation. Ordinary-path counts12555/12544 undercounted16/27files respectively; prior artifacts retained. Signed8534file identity was independently correct throughout. No package content or runtime code changed by this verifier repair.

V393: perfil BFF0.5.16 omite Sharp opcional y desactiva optimización runtime; lock152→122sin versiones nuevas,122stanzas y12571archivos restantes exactos,8534artefactos firmados exactos;294tests/1skip previo,92fases/4agendas/4lecturas/4cancelaciones PASS,OSV122/0.4packs reconstruidos/805files,68/797; HTTPbinario idéntico. No Sharp fix ni reapertura V386; native/tooling restantes yTEST02/03/07 abiertos,45/48. Ver reconstruction_evidence/OPTIONAL_IMAGE_DEPENDENCY_CONTAINMENT_V393.md.

V393 FAIL770 recurrence: Preflight249 correctly rejects the stale franchise summary840Markdown after the new report raises the exact portable count to841. Original full log/process receipt retained; no final Preflight JSON exists after this structural failure, and its attempted read is not evidence. Update only the inventory owner to165/1561/841/56, preserve all selectors/oracles, checkpoint250 and rerun the complete gate.

V393 cierre251: cuatro packs incorporados/805files reconstruidos; lock122versiones/0avisos OSV,12571archivos y8534firmados exactos en candidato e instalación limpia.294tests/1skip heredado,92fases/4agendas/4lecturas/4cancelaciones PASS. Preflight250:164pasos PASS/56perfiles; disponibilidad global BLOCKED sólo Docker. Inventario165/1561/841/56,integral68/797.45/48; TEST02integral,TEST03tooling/native restante yTEST07siguen pendientes. Sharp no se instala en este perfil; V386investigación diferida.

V393 tooling note: first finalizer patch had unprefixed multiline patch content and was rejected before creating the helper or changing files. Correct patch construction prefixes every added line; rejection retained in tool transcript. No product or gate failure inferred.

V394 FAIL782 / STALE_PINNED_CONSUMER_POLICY: canonical pnpm planner0.2.0 rejects current enterprise-web0.5.16 and Playwright0.1.37 with consumer input changed: package.json; Lighthouse remains valid. Exact before/after hashes and original code retained in V394 baseline-routing.json. Update only the reviewed consumer identities, preserve strict rejection of drift, reconstruct and test current real recipes before closing. No runtime or redistribution admission.

V394 FAIL782 REGRESSION_PROVEN: two stale current-consumer pins repaired in canonical selector0.3.0; three fresh canonical consumers prepare/verify, seven input mutations and three notice mutations rejected. Preflight now compares the planner policy to freshly composed current consumers, so future manifest drift cannot pass solely on synthetic unit fixtures. No fallback to arbitrary inputs.

V394 tooling note: first JavaScript helper construction rejected a nested Markdown backtick before any command or file write. Rebuilt the helper using chr(96) for fence matching; preserved rejection in the tool transcript. Canonical reconstruction and51tests then passed.

V394 cierre253: FAIL782 reparado con selector0.3.0/6files;51tests PASS,3recetas actuales,5declaraciones BlueOak exactas con aviso local y10mutaciones reales rechazadas. Preflight252:164pasos PASS/56perfiles, nuevo control de compatibilidad de3consumers; disponibilidad BLOCKED sólo Docker.165/1561/842/56;45/48sin promoción. Payload442 y22notices intactos; pnpm no ejecutado, runtime/redistribución pendientes. Ver reconstruction_evidence/PNPM_CURRENT_ROUTING_NOTICES_V394.md.

V395 FAIL719 recurrence: guessed pnpm.cjs was absent; source inventory and the existing V341 harness identify dist/pnpm.mjs. Rejected read retained in transcript, no evidence inferred. All subsequent paths come from inventory/receipts.

V395 FAIL783 / UNDELIVERED_VENDOR_PERMISSION_NOTICE: V338 identified QRCode MIT attribution inside qrcode-terminal0.12.0, but all3current V394 recipes still omit the vendor notice. Baseline and exact10bundle regions retained. Candidate adds the exact author header plus complete MIT permission text, outside the immutable payload, with digest verification; no Apache-to-MIT blanket relabeling or full source/build equivalence.

V395 FAIL784 / UNSUPPORTED_BLOCK_PROVENANCE: materializer rejected MIXED in the pending plan_install.py block before canonical promotion; executable block schema accepts AUTHORED/ADAPTED/VERBATIM only. Preserve failed pack/log. Use ADAPTED with explicit AUTHORED control logic plus VERBATIM notice origin, exact source and AND MIT license. Do not weaken materializer or relabel embedded third-party text as authored. Candidate source bytes unchanged; reconstruct into new destination and rerun59tests.

V395 FAIL783 REGRESSION_PROVEN:3recetas contienen atribución y permiso completos,12mutaciones reales rechazadas,59tests y6/6paridad de reconstrucción. FAIL784 REGRESSION_PROVEN: metadata ADAPTED soportada con fuente mixta explícita, licencia AND MIT y notices preservados;no cambio del materializer ni del código probado.

V395 FAIL785 / STALE_PROVENANCE_LEDGER: Preflight255 correctly rejected THIRD_PARTY_NOTICES after the notice-bearing planner changed AUTHORED→ADAPTED. Failed full log/process retained; no completed Preflight JSON exists. Update canonical counts1307/147/107→1306/148/107, same1561total, and document exact third-party header/permission origins. No verifier weakening or source-code change. Checkpoint256 and rerun complete Preflight.

V395 cierre257: aviso QRCode con copyright/permiso completos entregado y ligado a3recetas;59tests/12negativos PASS,6fuentes reconstruidas;BlueOak ypayload442 intactos. Preflight256164pasos/56perfiles PASS;Docker ausente.165/1561/843/56;procedencia1306/148/107. FAIL783/784/785 reparados;45/48sin promoción. semver-utils fuente gitHead en GitHub404, investigación pendiente sin bypass ni licencia inventada. Ver reconstruction_evidence/QRCODE_VENDOR_NOTICE_DELIVERY_V395.md.

V395 FAIL786 / OVERBROAD_NONEXECUTION_CLAIM: V394 report incorrectly described the entire Preflight252 as not executing pnpm; preserved logs show its existing offline frozen Playwright/Lighthouse installs. Only the three new recipes were never executed. Original report retained; corrected canonical wording distinguishes planner qualification from existing general runtime gates. No hidden execution, new authorization or changed runtime-admission verdict inferred; no code change. Repeat final structural/state validation after this documentary correction, not unchanged runtime suites.

V396 FAIL719 recurrence: guessed official_source_lock.json was absent; the materialized test identifies upstream-source-lock.json. No source or result inferred from the failed read. Subsequent paths use the materialized inventory. Rejection preserved in tool transcript.

V396 FAIL719 second recurrence: candidate README.md was assumed, but no such file is in its exact manifest. Use pack prose for documentation and inventory before reads; no missing-file evidence inferred.

V396 FAIL787 / HELPER_FENCE_WHITESPACE: candidate reconstruction helper stopped at48/50 source blocks because two canonical FILE sections contain blank lines. No canonical promotion or acquisition occurred. Preserve failure and diagnose exact missing paths; accept normative whitespace in the helper while retaining all50-block assertion, compositor hash validation and51-file byte parity.

V396 FAIL787 follow-up: broad helper string replacement also inserted literal plus signs into two generated fence lines; compositor rejected the isolated pack before materialization. Preserve failed helper, pack and log; constrain correction to helper regex only and rebuild into fresh destination. No canonical parser or gate changed.

V396 FAIL788 / SEMVER_LICENSE_TEXT_UNRESOLVED: prior APACHEv2 registry declaration did not provide original permission text. Governed exact4193-byte npm quarantine now reveals original LICENSE1839bytes SHA25680b98c1b20edfc51abd2c802ee7a1d3c5561151d108a24f19c207b6935beaed8, explicitly MIT OR Apache-2.0 with full MIT permission. Preserve declaration unchanged and add original notice delivery; expired signing key and whole pnpm admission remain conditioned. Failed guessed V395 rebuilt path resolved by exact inventory; no evidence inferred from absence.

V396 FAIL789 / UNDERESTIMATED_PROFILE_FANOUT: integration helper expected only initialization to select the acquisition core; exact JSON selection scan found30core plans plus the pnpm plan. All matched version pins were updated, then assertion stopped before inventory/evidence integration. Preserve old31plans and helper. Audit each selected plan against its before snapshot, changing only reviewed pack version, and continue all56composition gates rather than narrowing or omitting actual consumers.

V396 FAIL787 REGRESSION_PROVEN: canonical compositor reconstructs all51files byte-identically after helper-only whitespace correction; four core suites pass. FAIL788 REGRESSION_PROVEN for original license discovery and delivery only: exact npm artifact/6members, original dual-license file delivered to3real consumers,12mutations rejected,67planner tests; expired key and broader admission remain open.

V396 FAIL790 / HELPER_LABEL_REPLACEMENT: runner copied Preflight256 with an overbroad numeric replacement, corrupting parameter PromtoolSha256 to PromtoolSha258. PowerShell rejected argument binding before any gate executed. Preserve runner/log/process; rebuild from original with replacement scoped to preflight256 label only. No tool identity, digest, parameter or gate changed. Checkpoint259 before retry.

V396 FAIL791 / NONPORTABLE_EVIDENCE_PATH: structural release gate rejected machine-local home path embedded in V396 report. Preserve before report and failed Preflight259 logs; replace portable report locator with Temp/stage and pointer, preserving raw evidence in place. No exclusion or weakening of release validator. Retry complete Preflight260 after checkpoint.

V396 FAIL792 / STALE_COMPOSITION_COUNT_ORACLES: Preflight260 structural gate correctly rejected Shipping profile expected73/actual74 after new AUTHORED source-profile file. Audit all actual core-selected plans and derive exact before/after file sets; update only fixed counts affected by that one new file. Preserve strict equality, all56profiles and unchanged test/admission oracles. No suppressions or lower thresholds. Failed logs retained; retry261 after checkpoint.

V396 cierre262: licencia original semver-utils1.1.4 localizada y entregada bajo opción MIT completa;3recetas/12negativos/67tests PASS. Core0.4.88/51files/196sources/25profiles,140checks; planner0.5.0/6files. Preflight261164pasos/56perfiles PASS;Docker ausente.165/1562/844/56;procedencia1307/148/107.45/48sin promoción, Daybreak/V386 diferido. Fuente/licencia semver resuelta en alcance local; firma registry vencida y admisión pnpm restante condicionadas. Ver reconstruction_evidence/SEMVER_ORIGINAL_LICENSE_DELIVERY_V396.md.

V397 FAIL793 / VERIFIED_NOTICE_COLLECTION_NOT_DELIVERED: V340 retained47 exact texts but current recipes deliver only the three later supplements. Revalidate the47copy manifest and151253bytes; create a portable, hash-bound collection and deliver to all3consumers while preserving22original/25research scopes and open source/relinking/runtime conditions. No package clearance inferred from text presence.

V397 FAIL793 REGRESSION_PROVEN:47textos/151253bytes entregados en3recetas (141copias),12negativos reales,82tests y7fuentes reconstruidas PASS. Scope de evidencia retenida explícito y no promoción de paquete. La invocación inicial con cwd todavía ausente se rechazó antes de ejecutar; orden corregido y verificado.

V397 cierre264: conjunto47textos/151253bytes entregado en3recetas,141copias exactas/12negativos/82tests PASS. Preflight263164pasos/56perfiles PASS;165/1563/845/56 y1307/149/107. Docker ausente,45/48sin promoción. README/roadmap vigentes sincronizados y cortes previos preservados (FAIL794). Licencias/source/relinking/publicación restante y Daybreak diferido siguen explícitos. Ver reconstruction_evidence/PNPM_RETAINED_NOTICE_DELIVERY_V397.md.

V398 FAIL719 recurrence: follow-up guessed qualification.json in V341 evidence; read failed and no result inferred. The inspected script explicitly names ast-correspondence.json; this existing receipt supplies parser identity. A wildcard directory path in rg also yielded no usable inventory; exact parent inventory resolved the stage. Use observed paths only. Prior V397 closure264 remains verified.

V398 FAIL795 / NEXT_PATH_SOURCE_NOT_DELIVERED: original MPL text existed but recipes did not supply the fixed source/main manifest or transformation record. Official metadata gitHead and three Git blobs now match;4module statements structurally compare under explicit bundler adapters and10mutations reject. Deliver original317byte source,1010byte manifest and15978byte license with qualifications, not an invented runtime/build equivalence or blanket MPL clearance.

V398 FAIL795 REGRESSION_PROVEN_LOCAL_SOURCE_DELIVERY:3originales next-path entregados en3recetas,9copias/12negativos/90tests/8fuentes reconstruidas PASS.4sentencias AST y10mutaciones con adaptadores explícitos; no equivalencia runtime, tarball ni build pnpm. Fuente/licencia entregadas, resto de obligaciones conservado.

V398 cierre266: next-path fuente/manifiesto/MPL entregados,9copias exactas/12negativos/90tests/8files PASS. Preflight265164pasos/56perfiles PASS;165/1564/846/56 y1307/150/107. Docker ausente,45/48 sin promoción. Entradas vigentes sincronizadas con historia conservada. Daybreak diferido; integración funcional,source/relinking/publicación restante yrelease pendientes. Ver reconstruction_evidence/NEXT_PATH_MPL_SOURCE_DELIVERY_V398.md.

### FAIL-20260911-796 — continuidad de cierre desactualizada y verificación redundante

CONTRACT/OPERATION; DIAGNOSED. La assurance vigente seguía citando
GO_OBSERVABILITY_CORE0.3.0 y ausencia de integración general, aunque el candidato
actual es0.3.3 y TEST05/V372 demuestra un runtime distinto, no su adopción.
V393 ya prueba ausencia de Sharp en el perfil web; Daybreak/V386 conserva su
investigación diferida, no prueba una DLL instalada en ese perfil.
La repetición de suites completas por cambios pequeños consumió recursos sin
cerrar TEST02/03/07. El usuario exige agrupar por bloqueo material y no contar
pruebas repetidas como avance. Se preservan antes/después y42hashes de fuentes
V374: son idénticos, por lo tanto no repetir sus suites sin nuevo delta.
No reclasificar candidatos ni retirar las integraciones funcionales faltantes.
Recurrencia FAIL719: GO_ENTERPRISE_POSTGRES_MIGRATIONS.md era un nombre supuesto
inexistente; lectura rechazada, no se infirió ausencia de SQL ni se modificó código.

V399 FAIL796 CORRECTION_EVIDENCED: before/after y42hashes conservados; descripciones vigentes corregidas sin alterar48tests/requisitos/estados. No se confunde esta corrección con completar integración o seguridad.

### FAIL-20260911-797 — fixture de encuesta reutiliza un parámetro con tipos incompatibles

TEST/DIAGNOSED. El primer ensayo PostgreSQL del candidato de encuestas falla antes
de las operaciones HTTP: SQLSTATE42P08, UUID/text deducidos para el mismo $1.
Las54migraciones pasaron; servidor propio detenido, datos sintéticos conservados.
Separar tenant_id y tenant_code en parámetros distintos. No modificar el schema
ni atribuir un PASS al ensayo fallido. Evidencia: stageV399/feedback-r1/connected.log.

V399 FAIL798 / SURVEY_WINDOW_LOCK_RACE — REPRODUCED in staged candidate.
A request waiting on the definition row lock was accepted201 after closes_at
because accepting was computed before the lock wait. feedback-r2 preserves the
failing HTTP/OIDC/PostgreSQL test and rejected source; the owned server stopped.
Fix: hold definition/customer locks, check DB clock in INSERT itself, recover
only the original matching receipt before retain_until. Service and repository
share validation. No canonical executable pack or control promoted.

V400 FAIL797/FAIL798 CORRECTION_EVIDENCED: fixture typed parameters and database clock-after-lock correction pass feedback-final1 connected/restart/browser/retention CLI batch. 4 browser profiles,8 durable browser responses,8 writes; prior rejected r1/r2 logs remain. No claim of full TEST02 closure.
V400 FAIL719 RECURRENCE / FAIL799: guessed roadmap paths, PowerShell rg glob arguments and materializer PackPath parameter caused tool-only failures. No product mutation from these errors. Corrective rule: inventory exact paths and read parameter declarations before invocation; preserve tool errors, do not rerun unaffected product suites.

V400 FAIL800 CORRECTION_EVIDENCED: canonical publication stopped on generated next-env.d.ts imports. Only the two exact .next/type imports are removed for baseline comparison; generated and original SHA retained in next-env-generated-delta.json. Product sources are not exempted. reconstruction_evidence/CONNECTED_CUSTOMER_SURVEYS_V400.md

V400 FAIL801: final preflight268 rejected at resume because VERIFY_LIBRARY profile count805→828 was corrected after checkpoint268. No downstream preflight tests were credited. Product828/828 reconstruction remained valid. Freeze all canonical files before refreshing checkpoint; rerun final batch only after269 validates. Preserve preflight268 log/process receipt.

V400 FAIL802 CORRECTION_EVIDENCED: survey0054 lacked explicit transaction unlike prior migrations. Controlled duplicate-function failure left2tables; BEGIN/COMMIT correction leaves0 and succeeds in up/down/up. Stage feedback-atomic1 retains red/green logs and stopped owned server. Canonical Go survey0.1.1 and web compatibility0.1.1; no re-execution of identical browser/Go code.

V400 cierre270: FAIL801 corregido, Preflight269 verifica164pasos y checkpoint congelado; FAIL802 corregido por regresión SQL y reconstrucción final825archivos sin delta/2SQL/1lock. No se repite la batería general por cambiar sólo las transacciones de migración y su versión.45/48 conserva requisitos/estados.

V401 FAIL803: nested source scan preparation joined Noble against the BlueOak-only identity table and stopped before scanning. Corrected to hash-locked V339 source-license-cases.noble, exact source manifests and Git blobs. No scanner PASS or promotion was attributed to the failed preparation.


V401 / FAIL804 — 2026-09-11: adm-zip0.6.0 bundled in exact pnpm11.25.0 is affected by GHSA-vwc7-r8mq-g2x9 / CVE-2026-76845. Pinned OSV2.5.1 identifies one affected package among476 identified components; this is not a complete recursive native SBOM. The historical V321473/0 observation remains historical and cannot authorize current promotion. pnpm11.26.0 also contains adm-zip0.6.0; rejected as a fix without executing it. pnpm12.4.1 is an isolated inspection candidate with native optional executables, not a clean22-package runtime. Current pnpm promotion and new use are frozen; prior payload retained. Advisory lists no patched release. Reachability UNKNOWN pending call-site examination; no exploit or destructive probe executed. Response lane NO_UPSTREAM_FIX pending official patch review. Daybreak/V386 remains deferred. Evidence: Temp/elite-v401-20260911-closure/nested-scan-result.json and candidate-acquisition-receipts.json. Authority: https://github.com/advisories/GHSA-vwc7-r8mq-g2x9 . Reopen on demonstrated removal/containment or admissible fixed replacement; no risk acceptance.

V401 FAIL805: all three offline installs and command dispatches succeeded; package-byte comparison then failed on a Windows path longer than260characters. No installation is repeated or claimed failed. Compare the existing outputs through Windows extended-length paths and preserve the original comparison error.

V401 cierre: FAIL803 preparación corregida y scan real476completado. FAIL804 CONTAINMENT_PROVEN_IN_ISOLATED_CANDIDATE; original y admisión general permanecen BLOCKED. FAIL805 lectura larga recuperada mediante rutas extendidas, comparación opcional detenida tras4minutos sin PASS de equivalencia; instalaciones previas conservadas. FAIL806 inventario regex omitía espacios y contaba README como pack: sustituido por función canónica payload-aware,167/1588/1331/150/107. No suites repetidas. reconstruction_evidence/PNPM_ZIP_CONTAINMENT_V401.md

V401 FAIL807: the three explicit frozen offline installs succeeded without downloads, but the later pnpm exec Next probe omitted the matching store/offline settings. pnpm automatically recreated node_modules and downloaded71packages before reporting Next16.3.4. This invalidates zero-network for that dispatch, not the initial install receipts. Prior logs preserved. Fix the diagnostic invocation with explicit offline/store/cache/no-script/no-hook settings; no global runtime promotion.

V401 reconciliación final: comparación12571/12571idéntica completada antes de la solicitud de detenerla; no se detuvo proceso. FAIL807 separa las3instalaciones iniciales offline de un exec Next que descargó71paquetes por configuración omitida. Exec corregido con store/offline explícitos PASS sin descargas. Historial/log anterior retenido; no prueba global de ausencia de red. reconstruction_evidence/PNPM_ZIP_CONTAINMENT_V401.md

V401 / FAIL719 tooling recurrence: a readiness --report path outside project root was correctly rejected before output. Retried via JSON stdout, preserved previous root report and recorded the actual BLOCKED result; no readiness promotion.

V402 / FAIL719 recurrence: guessed input paths, a PowerShell literal glob and a delete/add patch targeting the same file caused tool-only errors. The patch was rejected before mutation; START_FRANCHISE was then changed by one Update operation. Exact existing paths/CLI declarations are the basis for subsequent calls. No failed read or patch is credited as a product check.

V402 / FAIL808 OPEN: Accounting.CreateJournal forwards negative columns and overflowing int64 totals to its repository. Reproduced with the exact prior source by provenance_scope in business-red.log: two lines with opposite negative columns pass as5/5; MaxInt64+MaxInt64+3 wraps to1/1 on each side. No PostgreSQL persistence claim. Candidate correction uses traced BC balance semantics and representation guards; integration/green evidence pending.

V402 / FAIL808 CORRECTED: exact BC derivation, negative/overflow caller regressions, real Commerce/accounting PostgreSQL and finite native fuzz passed;23/23outputs reconstructed,9changed/newfiles merged by hash into payment-connected-reference. See reconstruction_evidence/BC_EXACT_AMOUNT_ADAPTATION_V402.md. V402 / FAIL809: payment candidate RED reproduced partial refund promoted to full refund and old capture observation regressing terminal financial state. Corrected with distinct partial hold, monotonic refund observation and terminal-state preservation; semantic-green2 and admission PostgreSQL runs PASS. Harness semantic-green stopped before initdb because0056 increased migration inventory; it is not a product regression or PASS. Initial first-run module downloads are retained, later runs explicitly disable downloads.

V402 / FAIL810 TOOLING: staging inventory parser rejected the valid go.mod language fence in the Kiota pack before composition. The original snapshot is preserved; local packaging helper now accepts the same non-backtick language suffix as the canonical compositor. No product test failed or PASS was inferred. New snapshot uses a separate destination.
V402 / FAIL811: bounded native fuzz reproduced an invalid UTF-8 provider URL admitted by CustomerCheckout.Valid; JSON changed its bytes. The candidate now requires utf8.ValidString and keeps seed#7. Red and green receipts are preserved under checkout-policy-fuzz-agent; the exact three-file delta is integrated in the external reference. Canonical publication and frozen composition verification remain pending.
V402 / FAIL719 recurrence: a literal Windows rg glob and an empty apply_patch hunk were rejected. Reads/patch were corrected without claiming verification credit; no product effect from failed calls.

V402 / FAIL812 HARNESS: frozen combined probe gave connected/handover tests a database name that their disposable-target guards correctly reject. Two payment persistence/customer authorization tests passed; remaining tests failed before entering product behavior. Preserve connected-frozen-pg-final1 with stopped PostgreSQL. Correct runner uses elite_payment_connected_ and a separate elite_handover_ template clone; product bytes remain unchanged.

V402 / FAIL811 CORRECTED_CANONICAL: UTF-8 URL admission fix and permanent fuzz seed published with exact915file reconstruction and new graph gates. FAIL809 financial observation corrections also returned to canonical source. FAIL810 staging-only parser/candidate union errors corrected without product changes. FAIL813 harness selection: initial unit regex matched names ending in keywords and omitted intended checks; subsequent explicit required-test verification found25PASS plus one unrelated older PG test SKIP, which receives no PASS credit.16connected PG tests are separately no-skip. Histories retained.

V402 / FAIL815 TOOLING: default Python cp1252 decoded selected UTF-8 publication metadata; source payload packs used explicit UTF-8 and all915outputs remained exact. Restored plans and policy evidence from original UTF-8 snapshots, reversed the known complete-text transform in notices/verifier, and restored prior contract/state entries from exact backups. Rebind metadata receipts; use -X utf8 and explicit encodings in subsequent tooling.

V402 / FAIL814 METADATA: checkpoint275 rejected an evidence producer label longer than160characters before appending an event. Artifact/probe results were unaffected. Retain rejected state/log; shorten only producer metadata, refresh exact evidence hashes and retry the same unrecorded event.

V402 / FAIL816 TOOLING: handover packaging asserted an incomplete three-pack owner inventory after the verified12file overlay. Actual protected-client.ts owner is TYPESCRIPT_OIDC_PORTAL_ADAPTER, and franchise/page.tsx belongs to TYPESCRIPT_FRANCHISE_JOURNEY_PORTALS. No canonical pack was written. Correct the owner set to the discovered four packs; resume only after matching every already-integrated file hash, preserving the original pre-overlay backup. Product tests remain unchanged.

V402 / FAIL817 CORRECTED_CANONICAL: Next build rejected page export signatures with optional/union PageProps; required overload preserves framework contract and old direct calls. Browser then reproduced React418 from client locale/time formatting; server formatting removed nondeterminism. Preserve original failures and final65focused tests/build/browser receipts in HANDOVER_BROWSER_V402.md. FAIL816 packaging owner inventory corrected before canonical publication.

V402 / FAIL818 COMPOSITION_OPEN: building the independent ENTERPRISE_BACKEND profile with the new optional WhatsApp activation hook exposed pre-existing missing paymentbridge imports internal/channels and internal/outbounddelivery. Full franchise selected those owners; the previous six-profile receipt proved exact reconstruction, not compilation of every independent profile. Preserve whatsapp-backend-compat-build.log. Resolve the actual selected-pack dependency closure, then rebuild this affected backend. Do not treat the absent packages as new external modules or run go get.

V402 / FAIL818 CORRECTED_CANONICAL: ENTERPRISE_BACKEND selects the actual internal checkout dependency closure,39packs/531files. go build ./... PASS offline; separate optional-host absence test PASS. Previous franchise77/943 unchanged. See BACKEND_DEPENDENCY_CLOSURE_V402.md; no external module/version added.

V402 / FAIL819 PACKAGING_CORRECTED_IN_STAGING: WhatsApp delivery contained15CRLF/mixed payloads; the pack LF guard rejected them before canonical publication. Original source hashes and bytes remain preserved. Normalize only line endings, compare Python AST and installed Next/SWC emitted JavaScript, and verify SQL quoted literals are unchanged. A diagnostic first assumed the removed pre-7 TypeScript compiler API; inspection of the actual7.0.2package corrected the tool choice. No install or old PG/browser suite repeated. The authored host JSON template also required explicit newline='\n'. Evidence: Temp/elite-v402-library-infra/whatsapp-payload-lf-rebind.json and preserved first diagnostic log. Always write materializable UTF-8 payloads with explicit LF.

V402 / FAIL820 CORRECTED_IN_STAGING: the connected WhatsApp negative fixture reproduced a future-effective contact reaching model/domain execution. The canonical-contact candidate adds effective_at <= statement_timestamp() to the existing scoped resolver. Focused PostgreSQL regression now produces a durable handoff with zero model/domain/provider calls. The original RED and final proof remain in whatsapp-connected-final-receipt.json; do not count formatting-only rebinding as another runtime test.

V402 / FAIL821 COMPOSITION_OPEN: static selected-import audit of the new independent profiles found a serverless payment dependency on the absent outbound fence, and a backend fence test importing leadstream outside that profile. This is not a compiler execution or test failure. Select the existing fence/contact digest files for serverless; keep the omnichannel integration test in the full franchise profile where its owner exists. The backend payment integration tests still exercise its selected payment/fence path. Rebuild and compile affected profiles including test packages; do not add an external dependency or infer test PASS from compilation.

V402 / FAIL819,FAIL820,FAIL821 CORRECTED_CANONICAL: normalized payloads preserve explicit before/after semantics; future-contact resolver correction and connected/browser receipts published. Backend40/554 and serverless57/779 dependency selections reconstruct and compile including test packages; zero test execution claimed from compilation. Full franchise82/1032 exact. See COMMUNICATIONS_RUNTIME_V402.md.

V402 / FAIL822 HARNESS_ONLY: the added FX migration downgrade probe guessed entry_no instead of reading accounting0011, whose actual key is entry_id. The read-only snapshot query failed before any downgrade or product mutation; owned PostgreSQL stopped. Original fx-journal-downgrade/result.json and history-before.log retained. Correction uses the observed schema column in a new probe revision; product source and the already passing connected journal test remain unchanged. The same pre-read inventory rule also applies to locating the external execution kit; a missing guessed workspace directory is not a product failure.

V402 / FAIL822 CORRECTED_HARNESS: schema-derived entry_id snapshot query used in a new preserved probe; populated downgrade refused without history change, empty down/up PASS, owned PostgreSQL stopped. No product defect or repeated payment/FX domain test credited. See FX_JOURNAL_CONNECTION_V402.md.

V402 / FAIL823 CANDIDATE_POLICY_BOUNDARY: new stored-value tender glue changed the amount accepted by initial handover while retaining provider-only algorithm revision2. A new real SDK/PG regression reproduced an unselected stored-value handover effect. Candidate only, never published. RED: gift-tender-policy-red/connected.log and result.json; server stopped. Fix requires a separately hash-bound revision3/FULL_ORDER_WITH_STORED_VALUE option and rejects adjusted provider due under legacy profiles; do not re-label revision2 or its historical evidence. Full zero-provider branch and final admission remain pending.

V402 / FAIL823 implementation note: a guarded edit script stopped at the gofmt-normalized test helper signature after applying the candidate product/profile edits. No test ran in that attempt; the dependent runner file had not been written. Resume only the unexecuted tail using the observed signature, never replay applied substitutions. This is tooling preparation, distinct from the reproduced policy defect.

V402 / FAIL823 CORRECTED_IN_STAGING: exact revision3/FULL_ORDER_WITH_STORED_VALUE selection enforced; prior provider-only profile rejected with zero preparation effects. Partial SDK and full local funding through commercial receipt PASS. Candidate not yet composed; HTTP/host/UI, TTL final wait and source packaging remain. See reconstruction_evidence/STORED_VALUE_TENDER_V402.md.

V402 / FAIL824 CANDIDATE_COMMIT_EXPIRY: real PostgreSQL regression reproduced a stored-value operation committing after its approval lease expired while blocked on the outbox. RED gift-expiry-red/connected.log: seven operations/events instead of six; owned server stopped. Add a deferred commit constraint for proposals and applied operations, using database clock after all waits. Keep RED and verify expired transaction has no decision, account/entry, operation or event effect. Candidate only.

V402 / FAIL825 CANDIDATE_COMPILE_CORRECTED: optional stored-value host hook referenced logger instead of the observed slog owner. Initial gift-http-compile.log preserved; no runtime test was claimed. Two call sites corrected to existing slog. This was host glue compilation, not a source or dependency defect.

V402 / FAIL824 CORRECTED_IN_STAGING: deferred SQL constraints use clock_timestamp at COMMIT. gift-expiry-green/result.json PASS: expired outbox-wait operation leaves zero operation/decision/event and pending approval; expired proposal insert also rolls back. Existing gift/loyalty core remains six operations/events. Source profile minimum60seconds unchanged; only a typed test fixture uses shorter bound leases.

V402 / FAIL826 CANDIDATE_TYPECHECK: funding sum-type refinement parameter omitted explicit undefined under exactOptionalPropertyTypes. gift-web-first/typecheck.log retained. Widened only the callback optional-field type; runtime XOR rule unchanged. Repeat typecheck, then focused altered transport tests.

V402 / FAIL827 BROWSER_HARNESS: first stored-value browser run reached the correct scoped order but exact getByLabel did not resolve the nested select. Retained error-context shows the accessible combobox Operación and full form; no proposed operation ran. Use its observed role/name and a bounded12second action timeout. Product build unchanged; preserve gift-browser-pg-first and browser trace. Fixture extraction compile error for op was corrected before runtime and is recorded in gift-browser-compile[-corrected].log.

V402 / FAIL828 CANDIDATE_SOURCE_CLOSURE: two harmless exception fixtures proved the launcher imported an unmanifested package initializer and unchecked cached bytecode after verifying source hashes. gift-source-loader-red.log retains both real failures; changed locked source was already rejected. Load only verified in-memory source bytes, create a package with no filesystem search path, and use Python -I -S -B. No change to derived commercial arithmetic; candidate only.

V402 / FAIL829 EXACT_SOURCE_PACKAGING: official Odoo LICENSE is 43529 bytes with no final LF (SHA256 abc09dad5f84a76e1b0279237053cae16c03228ab27d8d467677054c2bd17eeb). The existing materializer/compositor always append LF, so exact VERBATIM packaging stopped before publishing either new stored-value pack. Seven staged owner revisions are preserved and canonical product packs remain unchanged. Add explicit optional final_newline=false with legacy default, reject malformed/duplicate flags, and test exact source bytes and updater round trips before promotion.

V402 / FAIL828 CORRECTED_IN_STAGING: verified in-memory package source only, -I -S -B; initializer/bytecode/hash-tamper regressions, IPC and host activation PASS. Source commercial engine unchanged. See STORED_VALUE_SOURCE_CLOSURE_V402.md.

V402 / FAIL829 CORRECTED_CANONICAL: explicit optional final_newline flag, legacy default retained; exact official43529byte license SHA preserved,42 reconstruction/negative/updater cases and legacy suites PASS. Root bootstrap/updater/contract and MARKDOWN-COMPOSITOR0.3.0 published from exact reconstruction. Product packaging resumes from seven staged owners.

V402 / FAIL823–828 CORRECTED_CANONICAL: handover profile boundary, deferred commit expiry, host/TS/browser harness fixes and verified source closure are included in reconstructed stored-value release. RED receipts retained; no runtime suite repeated solely for publication. See STORED_VALUE_CONNECTED_RELEASE_V402.md.

V402 / FAIL830 DISTRIBUTION_REGISTRATION: full structural verifier stopped before product checks because the new exact-EOF regression script was absent from the exact markdown_system script allowlists. gift-tender-library-verify.log preserves the rejection. Register this reviewed canonical test by exact filename in both verification and portable archive selectors; retain rejection of unknown scripts. This is distribution integration, not a stored-value runtime failure.

V402 / FAIL830 registration delta: exact EOF test filename added to verifier/archive allowlists; next run passed that boundary. Full verifier then exposed independent stale release controls (existing COMMUNICATIONS_RUNTIME_V402.json not classified; profile still hard-codes77 despite current84). Preserve gift-tender-library-verify-registered.log. Do not rerun broad suites until distribution inventory/profile controls are reconciled under T2810; five current exact compositions remain valid independent evidence.

V402 / FAIL831 STALE_RELEASE_BOUNDARY: structural release selector predates V402 JSON receipts and current profile84/1112. Requires explicit separation of distributable reconstruction evidence from local readiness/checkpoint receipts, symmetric verifier/archive selectors and current profile identity guards. Never allow arbitrary JSON or omit unknown files to force PASS; preserve historical evidence and negative selector tests. Owner T2810, no failure of the verified gift/loyalty runtime inferred.

V402 / FAIL832 PREPARATION_SCRIPT_CORRECTED: warranty owner extraction assumed IssueBulk had a successor function and then confused SQL alias r. with a Go receiver. Preserve before-owner snapshots and partial/complete logs; resume only the unexecuted tail, mask raw SQL for lexical checks, and stop Python/shell pipelines on first nonzero. Final SQL literals/order preserved and test package compilation passed; no runtime credit for compilation. Connected terms test now passes independently. See WARRANTY_CONNECTED_TERMS_V402.md.

V402 / FAIL833 PREPARATION_SCRIPT: exact aligned Go field search failed after the verified terms-source snapshot and before product edits. Snapshot retained; resume only the unexecuted tail using the observed JSON tag and precompute all edit assertions before writing. No runtime or source defect inferred.

V402 / FAIL834 FIXTURE_CONTRACT: first warranty J4 run applied67migrations and reached appointment resource assignment, then rejected version0. The fixture called the repository without the domain-prepared appointment version. Use the existing NewServiceWithProfile with the same exact policy and real clock for slot/request preparation; preserve warranty-connected-claim-first. Product writer is unchanged; do not loosen CAS.

V402 / FAIL835 WARRANTY_RESERVATION_MAPPING: connected J4 reached the shared human approval, which correctly rejected the CAS on a plan reservation carrying zero version. The existing writer persists reservation/version1 and returns the domain-prepared value. Warranty composition omitted Status/Version in its typed input. Supply the existing prepared reservation contract (reservation,1), retain original writer SQL/CAS and preserve warranty-connected-claim-domain-service. No relaxed approval or stock checks.

V402 / FAIL836 FIXTURE_DECIMAL_PRESENTATION: J4 approval and real FIFO work/rollback passed, but the fixture expected cost text740.0000 while the existing exact decimal formatter returns740. Both are the same numeric value; preserve the source representation and compare exact rational cost, never floating point or altered source formatting. Retain warranty-connected-claim-reservation-binding.

V402 / FAIL833–836 CORRECTED_IN_STAGING: prepared edits applied once; appointment uses domain service/profile, reservation supplies persisted initial state/version, and FIFO assertion compares exact numeric value without changing source format. Full connected J4 now PASS with4cases,1issue/2FIFOapplications, work rollback and released rejected/cancelled holds. Evidence WARRANTY_CONNECTED_CLAIM_V402.md; canonical publication remains pending.

V402 / FAIL837 CHECKPOINT_METADATA: checkpoint287 rejected the new evidence producer string over160characters before appending its event. Product J4 PASS unchanged. Preserve checkpoint287.log, shorten only producer metadata, refresh hashes and resume checkpoint/resume/plan at287 without replaying publication or incrementing the cursor again.

V402 / FAIL838 FIXTURE_HTTP_ACCOUNTING: all4claims passed through real HTTP and two lost responses recovered by GET. Final fixture counter omitted the intentionally denied store-to-factory reconciliation POST, observing3 instead of2. Retain warranty-connected-http-connected-first; account separately for denied attempt, committed request and explicit replay. No product retry or authorization change.

V402 / FAIL839 EXPIRY_FIXTURE_SETUP: the focused lease fixture supplied ends_at without a slot, violating the existing appointment_slot_time_check before any lease wait. Preserve warranty-connected-expiry-first. Represent the explicit historical unslotted appointment with both slot_id/ends_at null; full real scheduling remains independently proven by the HTTP J4 test. No product constraint changed.

V402 / FAIL837–839 CORRECTED_IN_STAGING: checkpoint287 repaired without duplicate event, HTTP counter includes denied attempt, historical unslotted appointment obeys original nullability. HTTP/host/downgrade and real60second commit-expiry gates PASS. No product authorization/cost/schema rule loosened. Evidence WARRANTY_INTERFACE_AND_PORTABILITY_V402.md.

V402 / FAIL832–839 CORRECTED_CANONICAL in WARRANTY_CONNECTED_RELEASE_V402.md: extraction/fixture/prepared reservation/amount/HTTP counter/expiry fixture and checkpoint metadata corrections included. RED history retained. Publication did not rerun unchanged runtime suites. Command generation encountered a tool syntax rejection and automatic long-command review rejection before filesystem mutation; separate apply_patch file creation and bounded execution succeeded, no product/security gate bypass.

V402 / FAIL840 SERIAL_OPTIONAL_IDENTIFIERS OPEN: actual PostgreSQL test rejects FRAME-2 without VIN at production_unit_tenant_id_vin_key. The public domain contract marks VIN/battery optional, but four NULLS NOT DISTINCT constraints allow only one missing value per tenant. Retain serial-supply-optional-identifiers-red/result.json and connected.log; fix only optional identifier null semantics using PostgreSQL18 documented NULLS DISTINCT, preserve mandatory serial and all present-value uniqueness. No synthetic VIN workaround.

V402 / FAIL840 CORRECTED_CANONICAL: inventory0.17.1/migration71 and exact regression published. Final runtime uses original canonical Operations; proposed J2 transaction helpers remain staging. Missing optional identifiers allowed, known identifiers/serial unique; rejected effects leak no events, downgrade preserves data. SERIAL_OPTIONAL_IDENTIFIERS_V402.md/json.

V402 / FAIL841 SERIAL_SUPPLY_HTTP_TEST_IMPORT_CYCLE OPEN: serial-supply-http-first/connected.log records a compile failure before runtime. Package postgres's new transport fixture imports httpapi, whose existing social adapter imports postgres. Move the two connected J2 fixtures to postgres_test and use exported repository constructors, matching warranty's external-package integration pattern. Do not change production dependencies to accommodate a test. Retain the failed receipt; require actual HTTP and host execution before closure.

V402 / FAIL473 diagnostic recurrence: a read-only rg invocation mistakenly included PowerShell ErrorAction and a guessed migration filename. It executed no mutation; corrected by querying actual migrations and reading0042. Reviewer is approval.decision.reviewer, not an invented request.decided_by; candidate guard corrected before the connected test. Final J2 first runtime PASS. Use filename inventory and actual column definitions before edits.

V402 / FAIL841 CORRECTED_CANONICAL: J2 integration fixtures use postgres_test and exported constructors. Real HTTP/host and final transport/downgrade gates PASS; failed compile retained. No production import changes. SERIAL_SUPPLY_CONNECTED_RELEASE_V402.md/json.

V402 / FAIL842 CATALOG_PREPARER_ASSERTION CONTAINED: extend-catalog-release-owners.py stopped before emitting migration73 because its preparer expected seven occurrences of serial_quality in the immutable migration72 approval header; the actual canonical header has six. The five candidate owner edits already completed. Preserve their before snapshots and resume only the schema/receipt tail with the observed count; do not replay prior mutations or rewrite migration72. No database runtime claim from this preparation.

V402 / FAIL843 CATALOG_PACKAGE_NAME_COLLISION OPEN: catalog-release-backend-first/connected.log fails compilation because contract.go's package variable hex conflicts with media.go's encoding/hex import. Migration73 applied successfully; no test executed. Rename the validation expression digestPattern, keep the standard library encoding/hex, and rerun the connected delta with a corrected J3/media target selector. Retain the first failed receipt; no claim from an unexecuted test.

V402 / FAIL843 CORRECTED_CANDIDATE: package compiles; media structural test PASS. The connected runtime reached its second distinct review and exposed FAIL844, so no J3 PASS yet.

V402 / FAIL844 CATALOG_OUTBOX_IDENTITY OPEN: catalog-release-backend-package-fix/connected.log shows the technical review colliding with the legal review at the existing outbox aggregate/type/version key. Both commands incorrectly reused draft resource/version1 with event type catalog-release.review. Bind the event aggregate to the immutable command identity/version1; retain draft and publication generation in its payload. Do not weaken the outbox unique constraint or invent increasing resource versions.

V402 / FAIL844 CORRECTED_CANDIDATE: independent review command events now persist; two versions publish and concurrency passes. The first rollback exposed the distinct source lifecycle mismatch FAIL845.

V402 / FAIL845 PRICE_BOOK_REACTIVATION_CONTRACT OPEN: catalog-release-backend-command-events/connected.log reaches rollback, then the retained Commerce owner rejects a second price-book.activated event at immutable version2 for the same book. Its lifecycle is one-time creation/activation. Preserve that source contract: each catalog publication must clone the exact approved snapshot into a new price book through the original creation writer, then activate it once. Bind the effective book to the publication and assert header/entries equal the captured snapshot. Old published books and quote references remain immutable; no event-key weakening or rewriting source versions.

V402 / FAIL845 CORRECTED_CANDIDATE: backend-effective-books and http-first execute five publications with fresh exact effective books, immutable source snapshots and rollback/expiry PASS. Original CreatePriceBook/ActivatePriceBook SQL and event versions retained. Canonical publication still awaits remaining J3 gates.

V402 / FAIL846 PUBLIC_CATALOG_CACHE_HEADER OPEN: catalog-cache-policy-red.log shows public200 overwritten to Cache-Control:no-store by the existing generic writeJSON helper. Conditional ETag values already change correctly, but the declared public revalidation policy was not emitted. Add a dedicated public catalog writer preserving max-age=0,must-revalidate, leave private/default writeJSON unchanged, and prove both200/304 headers. This was conservative caching, not an exposure or stale-price incident.

V402 / FAIL847 CATALOG_BUILD_RECIPE: schema/SEO tests PASS, but the selected old standard build helper repeated Turbopack rejection of the external node_modules junction. This limitation and the admitted build --webpack recipe were already documented by the handover/factory browser closure. Preserve catalog-next-first-build.log; use that existing recipe with all type checks active. Do not copy/re-admit dependencies or widen the filesystem root. Read the actual successful build receipt before deriving a new harness.

V402 / FAIL848 HOST_FIXTURE_KEY_LENGTH: host-feed-downgrade rejected activation because the manually counted synthetic HMAC string was33bytes, while the real contract correctly requires32. Preserve host.log. Build the fixture as bytes.Repeat with an explicit32 count; do not weaken key validation. The Next --webpack build passed independently (FAIL847 corrected candidate).

V402 / FAIL849 CATALOG_COMPOSITION_SELECTION CORRECTED_STAGED: candidate pack creation completed, but exact full-output assertion found that the storefront pack had not been selected. The selector inferred TYPESCRIPT-GO-API-WEB-BRIDGE from a filename; actual canonical pack_id is TS-GO-API-WEB-BRIDGE. Read observed metadata, preserve before plan hashes, and resume only composition generation. No canonical mutation or product runtime defect; source tests remain valid.

V402 / FAIL842 CONTAINED; FAIL843–849 CORRECTED_CANONICAL in CATALOG_CONNECTED_RELEASE_V402.md/json. Exact source price lifecycle/cache policies preserved, transport/test build recipe and32byte fixture corrected, actual pack_id selected. Retained RED history; five exact profiles and all new connected claims PASS.

V402 / FAIL850 CHECKPOINT_TASK_SYNC CORRECTED: T2802 was checked in the roadmap with canonical runtime evidence, but checkpoint294 still classified it in_progress. Validator refused before appending EVT-0294. Align task_progress (T2802 completed, T2804 in_progress), refresh evidence hashes and resume the same revision/event once. Product release and T2802 proof already valid; never undo a proven task solely to satisfy a stale cursor.

V402 / FAIL851 TRAINING_BROWSER_LOCATOR: browser reached submitted immutable answers and the distinct reviewer, then its getByLabel exact-text selector timed out on the select. The captured accessibility tree exposes a combobox named Resultado. Use that observed role/name and a12second per-action deadline; do not alter UI/domain/security rules or rerun the already-passed physical guard test. Preserve training-current-browser-first/connected.log and error-context.md.

V402 / FAIL851 TRAINING_BROWSER_SELECTOR CORRECTED_CANONICAL: observed accessible combobox role replaces wrapping-label lookup; final connected browser PASS3.51s, four committed effects and GET-only loss recovery. Source/domain guards unchanged by selector fix. RED and successful receipts retained in TRAINING_CONNECTED_RELEASE_V402.md/json.

V402 / FAIL852 CONTROL_RECEIPT_COMPOSITION_DRIFT: preparation refused stale T2802.inputs[0] plan SHA after training publication. Training release and execution294 remain intact; no event295 appended. Preserve the294 receipt and mint295 correspondence against exact1228 files,1193unchanged prior files and9declared shared-owner deltas with focused proof; replace the active control pointer. Do not blindly rehash historical proof or replay runtime journeys.

V402 / FAIL852 CORRECTED_CANONICAL: fresh295source-correspondence receipt selected; preparation READY_FOR_LIBRARY_WORK and checkpoint/resume/plan295 PASS. Historical294receipt retained. Release remains honestly BLOCKED for other controls.

V402 / FAIL853 CATALOG_SOURCE_FIXTURE_UUID: first focused test failed before product writes because fixture INSERT used the same untyped parameter as UUID and text concatenation (42P08). Preserve RED; cast through UUID explicitly in the synthetic tenant code expression. Product source not changed by this correction; resume only this new source regression.

V402 / FAIL854 CATALOG_NEXT_TEST_INIT:14focused tests and webpack compilation passed, but strict TypeScript rejected the test-only broad RequestInit cast because Next disallows a nullable signal. Preserve RED. Infer the concrete synthetic request options object instead of widening to platform RequestInit; no request assertions, product behavior or compiler strictness changed.

V402 / FAIL854 follow-up: concrete object also needs duplex literal preservation (as const). Second RED retained; run standalone tsc to expose all remaining type diagnostics before rebuilding again. No product or test assertion changed.

V402 / FAIL855 CATALOG_REVIEW_LOCATOR: browser completed source creation/lost-response recovery, PNG, draft and legal review, then exact getByLabel failed on a populated wrapped textarea. Snapshot proves textbox name Motivo de la decisión still present. Use observed textbox role/name for repeated review/publication fields, extending FAIL851 learning beyond selects. Product source and assertions unchanged; preserve RED.

V402 / FAIL853–855 CORRECTED_CANONICAL: synthetic UUID typing, exact NextRequest duplex type, and observed textbox role/name corrections retained with REDs. All new catalog role effects, host/down/goldens/fuzz and five exact profiles PASS. No product or assertion weakening. CATALOG_ROLE_AUTHORING_RELEASE_V402.md/json.

V402 / FAIL856 SUPPLY_GOLDEN_TRAILING_SEPARATOR: new cross-language fixtures placed U+2028/U+2029 at text end; both Go strings.TrimSpace and JS trim reject that domain input. Keep the RED; move separators into the interior with a suffix and recalculate only fixture expected bytes. Do not weaken product validation. BFF four tests and three other contract tests already PASS; repeat corrected contract fixtures and first build only.

V402 / FAIL857 SUPPLY_FORM_JSX: first new page build rejected a missing handler-closing brace in the registration form. Preserve RED and repair syntax; run standalone typecheck before another bundle. No runtime assertions or compiler strictness waived.

V402 / FAIL858 SUPPLY_NAV_TRANSLATION: visual inspection of the successful23-write browser exposed raw nav.supply. Add the new key to both existing Spanish/English dictionaries; retain the browser proof and verify translation/type closure without replaying its transactions.

V402 / FAIL856–858 CORRECTED_CANONICAL: domain-valid Unicode goldens, JSX handler closure and nav dictionary keys retained with REDs. Source atomic/browser/golden/type/fuzz proof and four exact profiles PASS. SUPPLY_ROLE_RELEASE_V402.md/json; no assertion/guard weakening.

V402 / FAIL859 WARRANTY_BROWSER_LEDGER_SIGN: both browser phases passed all18 role writes and both lost-response recoveries. The final SQL assertion incorrectly expected positive inventory issue cost; the retained FIFO ledger writes -740.0000 for an outbound issue while the service receipt displays positive740.0000. Preserve original failed wrapper and passed browser logs. Correct only the fixture assertion, then verify the already committed owned database read-only; do not replay the successful browser journey solely for an expected-sign typo.

V402 / FAIL860 PACKAGE_HELPER_NEWLINE CORRECTED: candidate packs/plans and exact inventory were already emitted, but generating rebuild helpers passed literal backslash-n to Python newline. Resume only helper generation with actual LF; do not replay the completed pack mutator. No product/source change.

V402 / FAIL859–860 CORRECTED_CANONICAL: signed inventory-ledger fixture expectation corrected and existing DB verified read-only; newline helper tail resumed only. Both actual browser phases, goldens/types and four exact profiles retained. WARRANTY_ROLE_RELEASE_V402.md/json.

V402 / FAIL861 NETWORK_EXTRACTION_TEMPLATE CORRECTED_CANDIDATE: the helper generator globally replaced SQL in an explanatory comment, duplicating a body outside a declaration. gofmt rejected before any execution. Regenerate only this candidate file from the retained exact pre-extraction source with direct interpolation, not global word replacement. Original SQL and contracts preserved; no canonical code affected.

V402 / FAIL862 NETWORK_ADAPTER_TYPE_QUALIFICATION CORRECTED_CANDIDATE: interface-derived transaction adapter retained unqualified fulfillment domain type names after moving packages. First compile rejected them. Qualify only observed domain types; preserve the explicit deny methods and original writer signatures. No test was executed and no schema/database mutated.

V402 / FAIL863 NETWORK_FIXTURE_ORGANIZATION_CODE CORRECTED_CANDIDATE: PG rejected uppercase ROOT against existing organization_organization_code_check before any network effect. The command transport now reflects that existing lowercase/hyphen constraint and the fixture uses valid codes. The overlap fixture uses different starts_on so it reaches territory interval exclusion rather than the earlier exact-key constraint. No database guard weakened; first RED preserved.

V402 / FAIL861–863 CORRECTED_CANONICAL: exact extraction template/type qualification and retained lowercase-code constraints fixed. Original SQL SHA correspondence, atomic/browser/host/golden/fuzz and five exact profiles PASS. NETWORK_ROLE_RELEASE_V402.md/json; no guard weakening.

V402 / FAIL864 HELP_CMS_UNICODE_SEARCH: first connected CMS run proved create/update/publish and late-outbox rollback but failed the accented uppercase search under the deliberately C-locale PostgreSQL fixture. Preserve help-cms-pg-first/connected.log. Bind both search operands explicitly to PostgreSQL18 built-in pg_unicode_fast, require UTF8/collation at host activation, and retain the accented regression. Do not weaken its assertion or change the cluster default. No canonical product published.

V402 / FAIL865 TRAINING_EXPORT_MODULE_ROOT: the exporter compiled the new21-guide source but its absent output was outside the materialized dependency tree, so the generated notification module could not resolve the existing pinned zod. Preserve help-release-export.log. Resume only the export/profile tail with an absent output under training_content, then retain its receipt outside. No dependency installation, source repetition or validator weakening.

V402 / FAIL866 HELP_CMS_FIXTURE_INDEX:11focused web tests and webpack compilation passed; strict TypeScript rejected unchecked indices into imported JSON golden rows in test-only code. Retain the RED build. Add an explicit required-row accessor, keep all expected values/assertions, and run standalone types before another build.

V402 / FAIL867 HELP_CURRICULUM_LIMIT: the materialized profile gate rejected seven admin lessons against the existing four-lesson limit. Preserve help-cms-corrected-http-profile.log. Split the intended material into two explicit admin courses, keeping every required guide, original runtime limits and human-review/no-grants policy. This is curriculum composition correction, not a relaxed guard.

V402 / FAIL864–867 CORRECTED_CANONICAL: explicit Unicode search, consumer-root export, checked fixture indices and bounded curriculum split. REDs retained; CMS/browser/training/host/goldens/fuzz and five exact profiles PASS. HELP_CMS_RELEASE_V402.md/json.

V402 / FAIL868 OVERVIEW_CONVERTED_LEAD: new source-grounded KPI regression proves the existing Overview counts converted as open (actual2, expected1 for new/converted/lost). The current domain and database use converted, not won. Preserve role-kpi-leads-baseline/connected.log; change only the terminal-state literal in the canonical-owner candidate, then validate the same assertion. No lead lifecycle change.

V402 / FAIL869 KPI_PATCH_MATCH_COUNT: the scoped UI/host script finished those files, then stopped before the generic BFF body edit because it assumed two identical session blocks. The GET uses different formatting. Resume only the unexecuted tail with the observed single POST block; preserve prior source edits. No gate weakened.

V402 / FAIL870 KPI_FIXTURE_TENANT_CODE: the first new metric scope/precision test PASS; the following source-owner test stopped on the retained globally unique tenant_code. Replace only the copied fixed fixture code with its generated tenant UUID. Preserve first RED and rerun only the previously unexecuted source-owner case.

V402 / FAIL871 KPI_REVERSAL_PRECONDITION: new metric proof passed pending/approved issue/replay checks, then the existing reverse owner correctly rejected a confirmed source order. Mark the synthetic source order cancelled before reading its new version, as required by the already admitted reverse contract. Keep the rejection and do not relax the business guard; only this source-owner test resumes.

V402 / FAIL872 KPI_RETENTION_FIXTURE: the source-owner test passed issue/reversal, program isolation and NPS0/threshold/two submissions, then PostgreSQL correctly rejected rewriting immutable survey retention dates. Seed a separate already-expired definition for the retention read. Preserve immutable policy and the RED; never disable its trigger.

V402 / FAIL873 KPI_SUPPLY_FIXTURE_BINDING: source-owner KPI test now PASS. The new factory join fixture was rejected because direct serial plan insertion lacked the original draft/lines/same-transaction creation evidence. Replace that setup with existing CreatePurchaseOrder→BindPlan→submit/confirm/start/register writers. No trigger bypass or fabricated durable step.

V402 / FAIL874 GENERIC_BFF_BODY_AUTH_BOUND: targeted original-source test RED (invalid-body400 before401/403), corrected nine-test batch PASS including streaming65537bytes, stalled-body deadline/cancel and invalidUTF8. Source baseline/restoration hashes in role-kpi-body-regression.json. Historical FAIL457 sender/recovery closure V312/V325 stays closed; this distinct body-boundary defect was inaccurately grouped under that old identifier in the cursor. New code only authorizes before bounded read and preserves action-specific checks.

V402 / FAIL868–874 CORRECTED_CANONICAL: current lead terminal value, scoped fixture/source writer corrections and bounded/authenticated BFF. All named tests PASS, RED history retained, exact five profiles. ROLE_METRICS_RELEASE_V402.md/json. Historical FAIL457 remains closed; it is not this distinct FAIL874 body defect.

V402 / FAIL875 INVENTORY_COMPILER_API: the private-locale inventory assumed the old TypeScript package JS entrypoint/createSourceFile API. The fixed TypeScript7.0.2 package exposes version and new unstable APIs instead. Read actual exports and use the already selected Next16.3.4 compiled Babel parser only for inventory tooling; no new product dependency. Candidate1343 copied once; do not replay its materialization.

V402 / FAIL876 PRIVATE_LOCALE_INVENTORY_BOUNDARY: candidate translation helper rejected an intentional leading separator before writing the batch; preserve nonempty display whitespace. Semantic inventory review found one event-handler notice incorrectly excluded as a JSX attribute. Restrict attribute exclusion to its direct value, retain original inventory and add p1193 in a separate complete inventory. Catalog candidate only; UI coverage is not proven by translated row count.

V402 / FAIL877 PRIVATE_LOCALE_ATOMIC_CODEMOD: pre-write parse rejected a locale binding collision with CMS content language, then a grouped constant declaration. No source was written by either rejected plan. Alias the private display locale, move only message-bearing declarators into their actual component consumers, and validate every planned file before the atomic publication loop. Final40-file direct wiring parsed; semantic/render/wire tests remain pending.

V402 / FAIL878 PRIVATE_LOCALE_WIRE_REASON: semantic review found two existing WhatsApp audit reasons translated inside a decision body. Restore both original wire literals; they are audit data, not UI messages. Forty-file AST audit now preserves all90 JSON serialization/recovery-storage argument expressions, with zero display calls in wire fields. Also include the two platform TSX owners omitted from app/components inventory; split paging presentation from server cursor utilities. Browser proof remains required.

V402 / FAIL879 PRIVATE_LOCALE_GUIDE_TEST_TYPE: five runtime display tests PASS; strict types found an unchecked version derived by split in the pre-existing notification guide type. The new test explicitly requires the version string before constructing the display input. Do not weaken noUncheckedIndexedAccess or change the immutable guide source.

V402 / FAIL880 PRIVATE_LOCALE_BROWSER_LOCALE: first connected browser test stopped before any business write: Chromium document navigation retained its context en-US even after setExtraHTTPHeaders(es); trace distinguishes document en-US from subsequent fetch es. Use two explicit browser contexts with real es-AR/en-US locale options. Keep the original rejection; do not change product negotiation to satisfy a misconfigured test. Fix test wire assertions from the actual CMS form contract before execution.

V402 / FAIL881 PRIVATE_LOCALE_BROWSER_TEXT_SCOPE: browser passed Spanish create, lost-response uncertainty, preference POST, English reload and GET-only recovery of the exact reference. A text locator then matched both the article paragraph and editor textarea. Scope the assertion to the article role; retain the RED and do not change rendered content or the lang assertion.

V402 / FAIL875–881 CORRECTED_CANONICAL: pinned parser API, complete literal inventory, atomic binding-aware codemod, immutable audit wire reasons, checked guide version and actual browser locale/accessible selectors. Final1article/3versions/3events with GET-only recovery; source/display/types/build/fourexactprofiles PASS. Search input selector also corrected within FAIL881. PRIVATE_LOCALE_RELEASE_V402.md/json.

V402 / FAIL882 MARKETPLACE_OWNER_INVENTORY: Owner lookup incorrectly assumed JSON manifest fields and guessed inventory/migration paths; the copied 1363-file candidate was intact. Resume verified files without recopying; use actual FILE headers and filename inventory. Official HTTPS acquisition succeeded for all six URLs after web extraction returned403; preserve transport discrepancy, not a provider-access blocker.

V402 / FAIL883 MARKETPLACE_MESSAGE_SCHEMA: First new unit compilation rejected a guessed Metadata field in channels.Message. The admitted fence hashes Text; remove the redundant metadata field and retain the immutable approval payload as the authoritative body. Evidence: marketplace-first-unit.log; compile and targeted contract tests must pass before promotion.

V402 / FAIL884 MARKETPLACE_FIXTURE_RESPONSE_CONTROLLER: The first PG/HTTP fixture stopped before provider effects because its response recorder hid the real connection required by the new five-second read deadline. Add Unwrap to the test recorder, preserving the actual ResponseController and fail-closed product deadline. Keep marketplace-connected-pg-first logs; source publication/migrations succeeded, provider path remains unproven until rerun.

V402 / FAIL885 MARKETPLACE_SOURCE_SQL_OWNER: The new source join referenced p.snapshot_sha256 even though the admitted publication owner stores the hash and canonical snapshot on its draft. Correct the binding to d.snapshot_sha256/d.snapshot_canonical, retaining the current-publication view and original catalog transaction advisory lock. No provider effects occurred; retain marketplace-connected-pg-response-controller logs. This is a new SQL glue defect, not an upstream or credential failure.

V402 / FAIL886 MARKETPLACE_FIXTURE_ROLE_SCOPE: The negative organization fixture used the root wildcard permission, which intentionally bypasses organization restrictions in the original identity owner. Replace the fixture with the exact catalog/marketplace permission set before testing wrong-organization denial. Preserve source-owner RED logs; do not weaken or rewrite the existing root authorization semantics.

V402 / FAIL887 MARKETPLACE_JSONB_BODY_CANONICAL: PG persisted approval and claimed the send, but the HTTP mapper compared canonical body bytes with JSONB output whitespace, so it failed closed as unknown before the provider mutation. Canonicalize the stored body with the existing approval serializer before comparison and sending; retain exact semantic approval and wire hashes, reject duplicate keys. Add a JSONB whitespace regression and cross-User-Product path tampering test. Scoped-role RED retained; no uncertainty was reset to permit a write.

V402 / FAIL888 MARKETPLACE_FUZZ_SOURCE_BOUNDARY: The isolated fuzz closure accidentally included dependency files named fuzz_test.go in addition to the intended target; the existing stored-value fixture then failed before the new fuzz campaign because its data was intentionally not copied. Rebuild a separate v2 closure allowing only internal/marketplacebridge/fuzz_test.go. Keep original module manifests/runtime sources, preserve the RED log and do not repeat unrelated business tests.

V402 / FAIL889 MARKETPLACE_PACK_TEXT_EOL: The pack builder rejected two new JSON texts emitted with Windows default CRLF. Normalize only those candidate text line endings to UTF-8 LF, record before/after hashes, and resume the still-unpublished package after its two existing-owner drafts. Official HTML snapshots and proven Go/PG behavior are unchanged. Always set newline explicitly for composed text artifacts.

V402 / FAIL890 MARKETPLACE_PROFILE_COMPILE_SELECTION: The reduced-profile compile command included internal/approval even though those profiles select only the application slot delta and omit that package. Their three selected packages compiled; the aggregate command correctly failed for the extra path. Derive the corrected package list from materialized directories and preserve first RED logs. No tests or unrelated journey suites are rerun; only the corrected compile aggregate is required.

V402 / FAIL882–890 CORRECTED_CANONICAL: actual owner/field/SQL/role/profile selection, response-controller test wrapper, JSONB canonical wire, isolated fuzz target and UTF8LF pack text. Narrow connected mutation, host/downgrade/fuzz/unit/vet/build/exact four profiles PASS. All RED retained; no duplicate provider attempt authorized.

V402 / FAIL891 MARKETPLACE_CHECKPOINT_FILE_WRITE: Checkpoint303 preparation encountered OSError22 opening the existing preparation receipt for write. The2658byte receipt remained valid, drive health/space and normal Archive attribute were checked; cause is not established. State remains302 with EVID148 already appended and two303 correspondence receipts minted. Resume only the remaining gate/report/checkpoint tail using atomic replacement; do not rerun publication or append duplicate evidence.

V402 / FAIL892 MARKETPLACE_MULTIPART_BOUNDARY: The first initial-publication unit run rejected every media upload before provider I/O: the deterministic multipart boundary used a14character prefix plus64hex characters, exceeding Go MIME boundary maximum70. Reduce only the glue prefix to six characters. Preserve marketplace-initial-first-unit.log; existing mutation regressions passed. No provider effect or canonical product promotion occurred. The corrected delta must pass before incorporation.

V402304 / FAIL892 CORRECTED_CANONICAL: deterministic multipart prefix reduced to fit70characters. Seven contracts, real API/PG/concurrency/recovery, host/downgrade/fuzz/build/vet and exact two profiles PASS. Original RED preserved in MARKETPLACE_INITIAL_RELEASE_V402.json.

V402304 / FAIL893 MARKETPLACE_PATCH_DELIMITERS CORRECTED: Three orchestration patch submissions were rejected by JavaScript parsing before apply_patch executed because literal Markdown fences were embedded in a template literal. No partial file edit occurred. Corrected submissions use plain text or chr(96) for generated fences. This recurrence is tool-input glue failure, not provider/account or product failure; preserve completed mutations and do not rerun them.

V402 / FAIL894 MARKETPLACE_CONTENT_TERMINAL_FIXTURE: The first connected CONTENT proof correctly rejected a multi-item User Product before any content PUT, but the new test expected a guessed delivery state failed instead of the original fence state failed_terminal. Correct only the fixture expectation; preserve marketplace-content-pg-first evidence, including accepted media upload and zero content writes. Do not weaken or rename the admitted fence state.

V402305 / FAIL894 CORRECTED_CANONICAL: new fixture now expects original failed_terminal state. Source/API/PG CONTENT, scope rejection before PUT, GET recovery, host/downgrade/fuzz/build/vet/two exact profiles PASS. No product fence change. RED preserved in MARKETPLACE_CONTENT_RELEASE_V402.json.

V402 / FAIL895 MERCHANT_SDK_QUERY_FIXTURE: The initial actual-SDK fixture rejected GET before any insert because it compared the whole URL against a path without query. The pinned generated REST transport always adds $alt=json;enum-encoding=int. Validate the path and that exact documented/generated transport query separately. Preserve merchant-sdk-worker-first.json; no product write occurred. The SDK/source integrity loader worked and is not replaced by a hand-written HTTP surrogate.

V402 / FAIL896 MERCHANT_SCA_IDENTITY_COVERAGE: Targeted OSV2.5.1 scan exited0, but the independent exact20-package identity comparison failed. Preserve merchant-connected-sca/scan.json and result.json; do not call the scan complete until every installed pinned package is accounted for. Inspect actual parser output before correcting the receipt projection or scanner input. No dependency/runtime version or product code changed.

V402 / FAIL896 CLOSED: OSV returned20 unique fixed PyPI identities with zero vulnerabilities;22 raw rows include two normalized-name aliases. The receipt comparator had treated PEP440-equivalent3.19/3.19.0 and3.0/3.0.0 as different versions. Exact wheel/metadata pins remain unchanged. Compare names canonically and versions with official pip-vendored packaging.Version; preserve the original scan and FAIL receipt. coverage-result.json proves20/20 without another network scan.

V402 / FAIL897 UTF8_RECEIPT_ENVELOPE CLOSED: the isolated Python receipt closer lacked -X utf8 and stopped on the first canonical Markdown read after writing its successful SCA receipt. No canonical history changed. Explicit UTF8 reads and idempotent receipt verification resume the tail; original network scan is not repeated.

V402 / FAIL898 / FAIL898 SCHEDULE_STATUS_FIXTURE_SECRET: initial connected test passed source/manual approval, no early dispatch and12concurrent callers->1actual adapter POST, then rejected the synthetic signed status. The test signed with app-secret while the reused appSecretFixture returns test-app-secret. Correct the fixture signer to that existing synthetic constant; no provider verifier/authentication product code changes. Preserve scheduled-communications-pg-first result/log; remaining cancellation/reschedule/recovery/optout cases were not yet reached.

V402 / FAIL899 / FAIL899 SCHEDULE_STATUS_OWNER_LABEL: after correcting the synthetic signature, the original signed observer inserted one status and the send/job were accepted/completed. The new test expected delivered, while the admitted status projection deliberately returns observed_delivered (existing router/history tests prove that contract). Correct only the test expectation; preserve scheduled-communications-pg-signature RED. No product status label or delivery guarantee is changed. Future fixture assertions must be read from the actual owner contract, not guessed.

V402 / FAIL900 / FAIL900 SCHEDULE_TEMPLATE_RECEIPT_RECOVERY: the connected test preserved a real fixture template send/receipt then lost stdout; the shared Python recovery verifier only supported conversational text. Add the exact existing approved-template payload builder to that verifier; retain the text branch and every source/request/recipient/profile/response/time binding. No recovery network call or resend. Preserve scheduled-communications-pg-observed RED, then verify both message kinds and tampering before the connected rerun.

V402 / FAIL901 / FAIL901 SCHEDULE_HOST_BASELINE_COUNT: the delta host runner retained 79 from its Merchant predecessor although the scheduled composition has 80 up migrations. The assertion stopped before PostgreSQL started. Correct the observed exact migration count; preserve scheduled-communications-host-first and this traceback. No migration or expected application behavior weakened.

V402 / FAIL902 / FAIL902 SCHEDULE_HOST_FIXTURE_SCOPE: the original constructor-only fixture used a nonexistent tenant/connection. On running its status loop against PostgreSQL, the original owner correctly reported RECONCILE_REQUIRED for that absent connection. Bind the host fixture to the exact already-created scheduled fixture tenant/org/connection, retaining all original scope validation; no provider traffic or connected journey is repeated. Preserve scheduled-communications-host-count RED.

V402 / FAIL903 / FAIL903 SCHEDULE_FUZZ_SOURCE_CLOSURE: baseline compilation of the isolated fuzz closure failed because the Go-only copier omitted the actual business-policy go:embed JSON and replaced nested payment module manifests. No fuzz inputs or product failure occurred. Include those exact source dependencies with hashes before the finite run; preserve scheduled-communications-fuzz.log RED. Isolated closures must include embedded assets and local replacement manifests, not only *.go.

V402307: FAIL898–903 RESOLVED_LOCAL for this delta: fixture HMAC/status label, product template receipt recovery900, host count/scope, fuzz embedded/module closure. RED traces retained; connected/host/boundary/fuzz/rebuild PASS. No replay of closed provider/domain journeys. SCHEDULED_COMMUNICATIONS_RELEASE_V402.md/json.

V402 / FAIL904 / FAIL904 CAMPAIGN_FIXTURE_QUOTE_NAME: the newly composed connected proof reused quote for both its existing Python-string quoting helper and a domain quotation value. Compilation stopped before any test body ran;81 migrations themselves passed. Rename the domain fixture value, retain campaign-connected-pg-first RED, and resume that exact already-migrated owned database rather than repeat migrations. No product change.

V402 / FAIL905 / FAIL905 CAMPAIGN_RESUME_EXPECTED_BASELINE: the copied host-resume harness required a prior whole-test PASS; this fixture intentionally had a compile-only failure after successful migrations. It stopped before PostgreSQL startup. Narrow the resume condition to the exact observed FAIL/connected, successful migration/final-stop receipts and build-failed log with no test execution; retain all runtime/path/migration ownership checks. No product or fixture expectation changed.

V402 / FAIL906 / FAIL906 CAMPAIGN_FIXTURE_LEAD_TRANSITION: the campaign source-drift fixture attempted forbidden new->qualified and the original CRM trigger rejected it. Earlier campaign create/replay/manual review/two-step dispatch/unknown predecessor/recovery/stop/marketing opt-out cases had passed. Use the actual original TransitionLeadAs new->contacted path with version1 and outbox evidence, then prepare the next campaign from contacted. No CRM rule or product campaign guard weakened; preserve campaign-connected-pg-guard RED.

V402 / FAIL907 / FAIL907 CAMPAIGN_PARTIAL_DECISION_RECOVERY: partial preparation/review correctly remained visible, but retrying the batch treated the already-decided first request as DECISION_UNCONFIRMED because the original HumanApprovals owner deliberately returns ErrNotPending. Add a read-only exact recovery branch for that error, bound to the same request hash, actor, decision, reason and unique immutable decision. Never reset/reapprove or synthesize a reviewer. Preserve campaign-connected-pg-partial RED and rerun only the partial-progress proof.

V402 / FAIL908 / FAIL908 CAMPAIGN_PUBLICATION_BUILDER_QUOTING: a Python generator used the same triple-quote delimiter for its enclosing script and Markdown payload. It failed parsing before producing or publishing any canonical file. Correct the enclosing delimiter and inspect the generated count/source-link assertions before publication. The candidate code and all passing runtime/rebuild receipts are unchanged.

V402308: FAIL904–908 RESOLVED_LOCAL: fixture quoting name, exact failed-build resume guard, original new->contacted transition, and actual partial decision retry fix907. RED traces retained; campaign/partial/replay negatives/host/appointment compatibility/fuzz/rebuild PASS. CAMPAIGN_CONNECTED_RELEASE_V402.md/json.

V402 / FAIL909 / FAIL909 CHECKPOINT_TASK_CHECKBOX_CORRESPONDENCE: T2805 was marked completed in state308 and its receipt, but the roadmap task checkbox still said pending. The checkpoint correctly refused to append event308. Correct the task row to the actually proven local scope, refresh hashes and run only checkpoint/resume/plan on the existing revision308; do not replay publication, proofs or the revision-increment helper. Original rejected log retained.

V402 / FAIL910 / IDENTITY_REBASE_MODULE_LINK: exact1451-output copy and26-file rebase completed; os.symlink failed with Windows privilege1314 before tests or runner creation. Replaced only the absent node_modules link with a native PowerShell junction to the already verified contained consumer. No installation, elevation or source replay. Future Windows candidate runtime links use Junction. RESOLVED_LOCAL.

V402 / FAIL911 / IDENTITY_TYPECHECK_ENTRYPOINT: host12boundaries/activation/shutdown and18TypeScript regressions passed. The typecheck launcher guessed @typescript/native-preview/bin/tsgo.js; this composition exposes its pinned TypeScript7 native wrapper at node_modules/typescript/bin/tsc. Node refused the absent path before compiling any source. Correct to the observed existing entrypoint and resume only types/vet/build; preserve original types.log. No toolchain download/version change.

V402 / FAIL912 / IDENTITY_CANDIDATE_RELEASE_REVIEW: saved0062 down migration unconditionally removed live session/revocation evidence; Go profile accepted a bridge URL path the TypeScript owner rejected; worker accepted a non-JSON media-type prefix. Before promotion, add populated rollback refusal, align exact root bridge URL validation and parse media type. Targeted host/profile and empty/populated rollback proofs required; prior lifecycle PASS remains valid for unchanged runtime SQL/SDK flow.

V402 / FAIL913 / PORTAL_STANDALONE_TEST_OWNERSHIP: the new protocol test read a config file owned only by the full Go lifecycle pack. Full composition passed, but standalone web would lack that fixture. Inline the exact nonsecret example in the test so the existing TypeScript pack remains independently testable. Four profiles contain this test and need updated exact rebuilds; the backend-only receipt is reusable because its plan, packs and577outputs are byte-identical.

FAIL910 follow-up: nested pwsh -Command did not bind positional args for the standalone junction setup; no tests executed. Used direct PowerShell New-Item -Literal known path/target pattern already proven. Keep shell setup native; do not treat -Command suffixes as script parameters. RESOLVED_LOCAL.

V402309: FAIL910–913 RESOLVED_LOCAL: native junction/CLI entrypoint; populated rollback refusal/Go URL and media checks; standalone test owns its exact fixture. Connected/host/types/profiles/fuzz/rebuild evidence retained. IDENTITY_PORTAL_RELEASE_V402.md/json.

V402 / FAIL914 / J5_EXPLICIT_BOOTSTRAP_PERMISSION: network root creation and its UI offered bootstrap only through the global * permission, while the admitted portal profile explicitly rejects *. Preserve compatibility for existing owners but allow the narrow network:bootstrap capability together with network:admin; provider remains role authority. Prove root creation with no wildcard, scoped ongoing administration and denial after role withdrawal using signed OIDC tokens and original HTTP/PG owner.

V402 / FAIL915 / J5_FIXTURE_OWNER_IMPORT: separating the passing cross-owner fixture into its actual portal/network owners retained an unused time import in the portal test. Compilation refused before test execution; remove only that unused import and retain owner-closure RED. No product change or protocol test rerun required; the signed network fixture itself is requalified after owning its issuer helper.

V402 / FAIL916 / J5_FIXTURE_TENANT_CODE: the new signed-issuer fixture resumed the existing owned database but reused the first test tenant_code, globally unique in the original schema. It failed before any domain command. Derive the synthetic code from the new fixture UUID; retain earlier tenant/evidence and owner-fixed RED. Each repeatable fixture must isolate both primary IDs and natural unique keys.

V402310: FAIL914–916 RESOLVED_LOCAL: specific bootstrap permission; separated fixture ownership/unused import; unique synthetic natural tenant key. All prior receipts retained. IDENTITY_J5_RELEASE_V402.md/json.

V402 / FAIL917 / PNPM_REGISTRY_HEAD_SIZE: all467read-only HEAD requests returned200 without Content-Length; no artifact body was acquired. The header-size assumption was wrong. Preserve license-header-receipts.json and probe one actual Range0-0 contract before batching it. Do not repeat HEAD requests or treat missing size as zero/known.

V402 / FAIL918 / PNPM_INSPECTION_HELPER_CONTRACT: local lock generator passed a Path to blocks(text), rejected before any acquisition or output. Read the existing helper signature; now pass UTF-8 text and compare the returned raw bytes. RESOLVED_LOCAL.

V402 / FAIL919 / PNPM_WASM_FLAG_DISCOVERY: Node24.20 rejects guessed --no-expose-wasm before running JavaScript. No pnpm invocation or changes. Inspect the actual V8/Node options before selecting a built-in local execution restriction; do not claim the unsupported flag works.

V402 / FAIL920 / PNPM_PACKAGING_REGEX_REPLACEMENT: JSON-escaped Spanish claim passed as regex replacement triggered bad escape before any canonical publication. Candidate pack/source files already exist; use a callable replacement and resume only the metadata/plan/result tail. Do not rerun prefix or duplicate README. RESOLVED_LOCAL.

V402311: FAIL917–920 RESOLVED_LOCAL: boundedRange sizing, actual helper text/raw signature, supported --jitless/--no-addons flags, callable regex replacement with tail-only resume. Receipts retained. PNPM_LOCAL_RUNTIME_V402.md/json.

V402 / FAIL921 / COMPLETE_GO_GRAPH_OFFLINE_MISS: three current module graphs resolved from cache, but standalone official_payment_webhooks needs uncached exact module metadata. GOPROXY=off refused before SCA. Preserve successful graphs and negative log; identify only missing metadata using go list -m -e, acquire fixed official metadata with checksum/profile evidence, then resume this graph and the scan tail. Do not rerun domain tests or relax version pins.

V402 / FAIL922 / GO_GRAPH_SECURITY_FLOOR_REGRESSION: current GO_ENTERPRISE_BACKEND_CORE0.4.5 go.mod no longer preserves admitted x/mod0.40 floor fromV313; full graph selects vulnerable0.37 (GO-2026-6179/6180). Official current Go advisories confirm fixed0.40. Restore the canonical owner floor using already admitted/cache-verified module; preserve graph and demonstrate compiled package/binary correspondence so no unchanged domain journeys rerun.

V402 / FAIL923 / GO_METADATA_ERROR_INTERPRETATION: GoMod in the -e record names an existing checksum-valid file, not the missing metadata. A file-only GOPROXY diagnosed missing .info for exactly two versions; selecting fewer JSON fields does not avoid it. One unnecessary .mod response matched existing bytes and no cache file was replaced. An unquoted PowerShell comma argument was rejected before Go ran; quote native arguments. Retain diagnostics and resolve only the two .info files. RESOLVED_LOCAL after the following verified acquisition.

V402 / FAIL924 / DEVSKIM_FAILURE_DIAGNOSTIC_MASK: the admitted runner failed a dotnet invocation but its strict-mode error formatter indexed [-20..-1] on fewer than20lines, hiding the actual cause. Candidate changes only formatting to Select-Object -Last20; preserve initial log and record actual command/error on the next bounded attempt. Return this correction to the canonical DevSkim owner before continuing.

V402 / FAIL925 / DEVSKIM_SDK_INCOMPLETE: retained dotnet10.0.400 passes --version but its shared10.0.11 contains only32files and lacks System.Security.Cryptography.dll; restore fails before source compilation. Preserve original failed logs. Reacquire the same official SDK ZIP with published SHA512 to an absent isolated destination; do not repair or overwrite old evidence, change SDK version or conflate host-executable hash with full SDK integrity.

V402 / FAIL926 / SARIF_UNICODE_PROJECTION: Python splitlines counted Unicode paragraph separators in golden JSON as newlines, so draft displayed lines differed from SARIF LF coordinates. Final per-occurrence review uses exact UTF16 charOffset/charLength and source SHA plus independently computed LF line; golden digest spans inspected. No source change or analyzer rerun; initial draft retained. RESOLVED_LOCAL.

V402312: FAIL921–926 RESOLVED_LOCAL:2metadata.info fijados; floor0.40; diagnóstico Go correcto; DevSkim short-error; mismoSDK completo5577files; ubicación SARIF UTF16. RED y receipts conservados. COMPOSITION_SECURITY_RELEASE_V402.md/json.

V402 / FAIL927 / DOCUMENT_MIGRATION_COUNT: initial document harness assumed the highest migration number85 equals the selected count. The current composition intentionally omits historical63/64; there are83 selected up migrations. Failure occurred before PostgreSQL start or migration execution. Derive the exact baseline set from the1481file inventory and add only0085; never infer count from filenames. RESOLVED_LOCAL.

V402 / FAIL928 / DOCUMENT_HARNESS_CACHE: inherited historical probe environment selected an older isolated GOMODCACHE, missing the already verified AWS modules. Setup passed83migrations and empty down/up; Go test stopped before any document test. Pin the known verified <USERPROFILE>/go/pkg/mod cache explicitly, keep GOPROXY=off, retain migration receipts and resume the existing owned database without rebuilding it. RESOLVED_LOCAL.

V402 / FAIL929 / DOCUMENT_TIDY_FLOOR_SUM: go mod tidy removed unused security-floor checksum entries; restoring the x/mod0.40 require line alone allowed compilation but full readonly module-graph resolution correctly refused. Merge the exact prior verified go.sum entries with new AWS entries, reject any conflicting hash, retain the floor comment and resolve the complete graph before publication. No module version or compiled source change. RESOLVED_LOCAL.

V402 / FAIL930 / DOCUMENT_NOTICE_MULTIPLE_OWNERS: final wording was updated in only the first of two mutually selected THIRD_PARTY_NOTICES owners. Four exact-profile verifiers caught a one-file mismatch; code and other files matched. Update both declared canonical owners from the same candidate text, retain failed receipts and rebuild only the affected profiles. Never select the first owner when alternative profiles own the same target path.

V402313: FAIL927–930 RESOLVED_LOCAL: inventario real83migraciones, caché de módulos fijado, floor+checksums preservados, ambos owners de notices actualizados. RED retenidos. DOCUMENT_REFERENCE_RELEASE_V402.md/json.

V402 / FAIL931 / DOCUMENT_CHECKPOINT_CURSOR_LENGTH: checkpoint313 rejected next_action over800characters before event insertion. The state revision was already advanced by the helper. Shorten the cursor, refresh receipts, and execute only checkpoint/resume/plan tail with event313; never rerun the revision-incrementing prefix. RESOLVED_LOCAL.

FAIL-20260913-932 OPEN — T2807 review: conversation store completes/retries without attempt generation or live lease; retention does not filter history/replay. Responses continuation omits instructions, which official API does not inherit. Reproduce against current owner, fix canonically and retain regression evidence before local promotion. Historical training already has admitted HISTORY-MODEL-TRAINING-PIPELINE0.1.0 (LIB-R10); reuse, do not invent another trainer.

V402314: FAIL-20260913-932 RESOLVED_LOCAL_DELTA: generation/lease/retention/continuation RED-to-GREEN. T2807 sigue abierto. FAIL-20260913-933 RESOLVED_LOCAL: actualización de firma en test duplicó generación; error de compilación retenido pg-lease-green/lease-red.log, reparación puntual y test PASS. AI_RUNTIME_GOVERNANCE_V402.md/json.

FAIL-20260913-934 OPEN — T2807 connected audit: domain gateway accepts quote quantity but drops it from the single-vehicle quote contract; order-status customer endpoint uses the shared service subject instead of the resolved contact. Complete reference composition with existing domain owners, contact-bound status and explicit single-unit guard. Existing EvalSuite also accepts NaN threshold/nil assertion and lacks cancelable per-case/required-case reporting; adapt existing owner, no model-quality claim from fixtures.

V402315: FAIL-20260913-934 RESOLVED_LOCAL: cantidad, consulta por contacto y evals con caso requerido/cancelación/NaN/nil/identidad probados. FAIL-20260913-935 RESOLVED_LOCAL: usar recibo histórico audit-final y code_sha, conservar manifest previo. FAIL-20260913-936 RESOLVED_LOCAL: guard de log existente bloqueó launcher contra staging anterior; ruta corregida, ningún gate falso. AI_CONNECTED_REFERENCE_RELEASE_V402.md/json.

FAIL-20260913-937 RESOLVED_LOCAL: textual checkpoint revision replacement altered the Python314 runtime path. Missing executable rejected before validation/event; correct fixed path and resume only the tail. Do not use global numeric replacement for runtime-bearing scripts.

FAIL-20260913-938 OPEN — T2808: Dockerfile tags are mutable and COPY misses selected local Go modules; no native local build/deploy/rollback entrypoint binds this composition. Add AUTHORED orchestration around existing pinned Go/Next/pnpm and Windows reference launcher; exclude unqualified container/cloud route from local acceptance. Also restore stale gap-map counts and preserve historical cuts.

FAIL-20260913-939 OPEN — First T2808 Go and Next builds PASS; standalone copy rejected nested pnpm paths at MAX_PATH. Preserve build-a result/logs, fix only delivery projection with extended local paths, confined link flattening and finite budgets. Do not rerun successful compilation to diagnose a copy-only failure.

FAIL-20260913-940 OPEN — Local operator contention test found lock-file read before byte-range lock and leaked handle on failure; lock first, close on contention. FAIL939 follow-up: Next standalone on Windows retains junctions into the exact installed consumer; projection correctly rejected an undeclared root. Declare only that same qualified consumer node_modules as an additional read root, flatten with budgets, reject every other escape.

FAIL-20260913-941 OPEN — Local rollout/failed-candidate rollback/controller recovery and order replay actually succeeded; fixture assertion used ID instead of JSON id. Shutdown correctly emptied every process tree but wrapper misclassified Next native cleanup: its pinned start-server.js intentionally exits 143 on SIGTERM, and synthetic emit must pass the signal name. Correct bridge to that actual contract; accept 143 only after requested graceful web stop, require Go0 and wrapper0. Reuse unchanged compiled artifacts for this lifecycle delta.

FAIL-20260913-942 OPEN — Two independent full Go/Next builds from exact source both run, but complete artifact comparison rejects differing Next preview keys and absolute RSC resource identities. Go binary and unchanged dependency/source/notices bytes match. Pin public LOCAL fixture preview cache via the exact Next16.3.4 cache contract and require an explicit absent, stable build workspace; archive the owned completed workspace before the second independent build. Preserve both failed comparisons. Never call different chunks byte-identical or patch application bundles to hide drift.

V402316: FAIL938 RESOLVED_LOCAL_DELIVERY; FAIL939 RESOLVED_LOCAL con proyección de enlaces declarada y paths extendidos; FAIL940 lock antes de leer; FAIL941 JSON id y cierre Next143; FAIL942 dos builds exactos con workspace estable/preview fixture/traces. Resultados RED intactos. LOCAL_REFERENCE_DELIVERY_V402.md/json.

FAIL-20260913-943 OPEN — T2809 scope audit: current production-shaped API host lacks the existing closed-attribute OTel/Prometheus middleware; the qualified instrument() lives in a separate synthetic identity executable. Reuse the exact middleware in an importable owner, bind optional loopback metrics to the real host identity, and connect local alert/retention/load evidence. Do not promote the synthetic principal or claim production operations from reference proofs.

FAIL-20260913-943 update: first offline tidy could not find already-admitted HTTP telemetry modules in the root cache. Resolve only through the retained V375 file proxy plus the current local cache; no floating lookup or source change.

FAIL-20260913-943 update: first offline tidy could not find already-admitted HTTP telemetry modules in the root cache. Resolve only through the retained V375 file proxy plus the current local cache; no floating lookup or source change.

FAIL-20260913-944 RESOLVED_LOCAL_EVIDENCE — Middleware correspondence initially omitted gofmt normalization; retain failed assertion and compare exact pinned gofmt output after explicit package/name changes. Staging helper inspect.py also shadowed Python stdlib during traceback formatting; renamed inspect_observability.py. No product algorithm changed and no failed runtime claim was promoted.

V402317: FAIL-20260913-943 RESOLVED_LOCAL: middleware importable en host OIDC, activación loopback explícita, caché offline histórica recuperada y piso x/mod conservado, alerta real/repair/replay/WAL/native stop demostrados. Candidatura target separada, nunca producción inferida. LOCAL_REFERENCE_OPERATIONS_V402.md/json.

FAIL-20260913-945 OPEN — T2801 residual source audit: GO_CUSTOMER_SURVEY_API still calculates NPS in AUTHORED business code, matching an explicit V402 provenance rejection. Preserve the old owner and find/admit a narrow official source implementation; do not relabel by citing Bain formula or passing tests. T2809 closed317 is unaffected.
FAIL-20260913-946 DIAGNOSED — isolated NPS regression harness used subprocess capture_output for pg_ctl start. PostgreSQL became ready on the correct owned loopback54507 cluster, pg_ctl exited, but inherited capture handles kept Python waiting; no test/database creation happened. Cancelled only the verified task runner16400 and stopped only the owned cluster. Preserve test319-first.py/nps-postgres.log. Use file-backed stdout/stderr for child service launch and cleanup based on actual postmaster identity, then run the still-unexecuted regression. No product defect or PASS inferred.

V402319: FAIL-20260913-945 RESOLVED_LOCAL with actual PostHog derivation invoked by existing PG summary; FAIL-20260913-946 RESOLVED_LOCAL file-backed subprocess logs, real84migration/OIDC/HTTP/concurrency/restart PASS, all owned processes stopped. NPS_SOURCE_ADAPTATION_V402.md/json.

V402320: disposiciones por scope en LIBRARY_ASSURANCE_REVIEW_V402.json.384 CONTAINED_LOCAL con release/FAIL831 pendiente;385 RESOLVED_LOCAL_SCOPE sin equivalencia general de21cores;532/804 CONTAINED_LOCAL_USE sólo receta restringida311 y excluida de distribución.945/946 RESOLVED_LOCAL319. Historial intacto; ningún fallo nativo se declara reparado.
FAIL-20260913-947 DIAGNOSED — first full public-source verifier321 stopped on official relative marketplace endpoint /users/123456789. The machine-home detector applies case-insensitivity to POSIX paths, confusing lower-case API users with the case-sensitive macOS home directory name. Product URL and fixture stay intact. Scope case-insensitivity only to Windows home paths; preserve POSIX home checks and add positive/negative path regressions on the actual helper before resuming. verify321.log/JSON preserve the rejection; no global PASS.
FAIL-20260913-948 DIAGNOSED — corrected verifier now rejects a real machine-local home locator in public historical CONNECTED_DELTA_ASSURANCE_V402.md. This is evidence portability, not a product/source failure. Preserve original evidence bytes externally; inventory all public textual occurrences before another full run, replace only rendered local locators with explicit portable aliases and record before/after. Do not weaken the path gate or change embedded source hashes.
FAIL-20260913-949 DIAGNOSED — full verifier321c reached structural pack validation and rejected GO_BUSINESS_POLICY_PROFILE: heading Apply order and rollback does not equal the required Apply order contract. Compositor outputs previously passed because payloads are intact; that did not validate the outer normative heading. Preserve321c rejection; run the actual structural contract over all205packs once and correct only malformed outer headings/metadata, keeping exact payload hashes and substantive rollback text.

FAIL-20260913-950 DIAGNOSED — verifier321d reached core-claim audit and rejected stale GO-HELP-CENTER-CORE SHA. Inspect all21 rows and retain V374 history; reconcile current pack metadata/tests with actual later evidence, without rerunning unchanged domain suites or inferring enterprise source.

FAIL-20260913-951 DIAGNOSED —321e completed all56composition profiles but final inventory summary still counted954Markdown instead of955 after the successor audit was added. Preserve full log. Recalculate public Markdown inventory after assembling all publication evidence, before running the gate; no product code or profile changed.

V402321: FAIL831/947/948/949/950/951 RESOLVED_LOCAL_PORTABLE_SELECTION. Exact43selector+50distribution cases,205outercontracts/58unchanged repaired-pack payloads,21current core bindings/13negative cases and full VERIFY_LIBRARY56profiles PASS. Three corrected rendered public locators retain original bytes/SHAs externally. Final current archive/durable product/signature is separate T2810 work; no historical native or production failure is promoted.

FAIL-20260913-952 DIAGNOSED — release builder catalogued71installed top-level npm packages but omitted nested Next runtime legal scope. Actual standalone contains129package manifests, including unversioned bundles and15named roots without direct license files. Inventory exact runtime manifests, retain original bundled legal texts and fixed-source supplements, explicitly exclude unshipped build-only tools, and reject unknown/mismatched/missing runtime notices before final release. No prior production or redistribution acceptance inferred.

V402322: FAIL-20260913-952 RESOLVED_LOCAL_NOTICE_INFRASTRUCTURE. Actual129runtime manifests bind original/declared notice sources;155unique retained texts,12negative publication tests,4exactsource profiles and8triaged lint findings. Final current runtime build/SCA/signature remains T2810. No historical316 binary is labeled as322/ARCA.

FAIL-20260913-953 DIAGNOSED — final signed-release integration review found the legacy Windows gate replacing two literal backslashes instead of the single path separator; its verifier also trusts an allowed-signers file embedded inside the candidate itself. Before real signing, correct portable path normalization, require an independently supplied external signer policy and fixed namespace, recheck input hashes around both builds, and bind SCA to the declared dependency manifest set. Preserve0.1.0 source/history; qualify with meaningful path/source/trust regressions, then use real final artifact. No private key or live account requested.

FAIL-20260913-954 DIAGNOSED — signed-gate regression attempted a backslash ZIP name through Python ZipInfo construction, which normalized Windows separators to forward slash. The gate correctly accepted the resulting safe a/b.txt archive; the negative fixture was wrong. Preserve signed-contract323 original logs/ZIP. Set the fixture entry filename after construction and assert the persisted ZIP name equals the intended adverse path before expecting rejection. No product gate weakening or false security finding.

FAIL954 fixture refinement: ZipInfo also normalizes its display filename while reading. The second archive contains the intended original name; compare orig_filename (and bind that original name in the synthetic manifest), not the normalized filename API. Preserve signed-contract323b; still no gate or source claim changed.

V402323: FAIL953 RESOLVED_LOCAL_GATE_CONTRACT (external signer trust, single-separator ZIP, source/profile stability, exact scan refs and failed evidence retention). FAIL954 RESOLVED_FIXTURE_INPUT with persisted orig_filename assertion; original erroneous tests retained. Actual final product/SCA/signature remains pending, no synthetic receipt promoted.

FAIL-20260913-955 DIAGNOSED — Windows OpenSSH key generation installs an explicit BUILTIN Administrators ACE on the new private key after the protected directory ACL was set. The post-generation assertion correctly refused publication. Retain the existing key and failed creation script; restrict only this task-owned key directory/files to the current user and SYSTEM, verify the resulting protected ACL and public-key correspondence, then resume the public receipt. Never print, copy, overwrite or regenerate the private key to hide this failure.
FAIL955 repair refinement: applying a new FileSecurity object asks Set-Acl for SeSecurityPrivilege on this host. No key bytes changed. Use the existing Get-Acl descriptor and modify only its access rules; retain the original pre-repair ACL metadata and first repair source.
V402324 preparation: FAIL955 RESOLVED_LOCAL_ACL. Native icacls modifies only the access DACL without requesting audit privilege; current user and SYSTEM are the sole explicit FullControl principals on directory, private/public key and allowed_signers. Public/private correspondence verified, original key retained. No private key printed or copied.
FAIL-20260913-956 DIAGNOSED — ARCA selected generator Offline only guarded WSDL acquisition; tool/project restore could still contact NuGet and its audit claimed a current online result. Generated client output is also a required external step and the worker lock predates the existing Cryptography.Xml floor. Complete deterministic offline generation/explicit SCA separation, reconcile exact existing-version locks and wire generated clients into reference build. Preserve old owner/source; never fabricate a NuGet audit PASS or request user fiscal credentials.
FAIL-20260913-957 DIAGNOSED — cached dotnet-svcutil8.0.0 nupkg SHA512 does not equal its .nupkg.metadata contentHash. Generator has not executed. Preserve archive/metadata, inspect both cached originals and official fixed package identity/signature before using tooling; do not overwrite cache or relax the pin.
FAIL957 RESOLVED_TOOL_CORRESPONDENCE: original svcutil nupkg exactly matches V123 official SHA256 b889780662e42d0e132e91c1821d1dcf810728a3775af226a1c8cca0ea6a2d8a and17915141bytes; all113 payload files match the original archive. Cache metadata hashes are retained as differing observations, not treated as authority or rewritten. No new version/source acquisition or substituted payload.
FAIL-20260913-958 DIAGNOSED — ARCA test preparation copied the omitted credential tests but not their same-owner integration project; restore warned and compilation correctly rejected the missing reference. No runtime test ran. Complete the selected owner dependency closure and validate all ProjectReference paths before resuming only the failed build/test and untouched projects.
FAIL958 refinement: the integration-test ProjectReference forwards an empty GeneratedClientRoot as a global MSBuild property, overriding the child default. Set the same relative default on the caller as well; preserve both failed builds and reuse successful restore receipts before resuming compilation.
FAIL-20260913-959 RESOLVED_FIXTURE_TYPE: connected ARCA fixture used a long sequence for the generated FERecuperaLastCbteResponse CbteNro int field. Compiler rejected the fixture before execution. Use checked int conversion at the exact generated boundary and resume only fixture build/untouched Go compile; successful locked restores retained.
FAIL-20260913-960 DIAGNOSED — first actual HTTP/PG/Go/UDS/.NET fiscal composition cannot reach sequence allocation: invoice stays queued after LastAuthorized failure, despite separate suites passing. Hypothesis: Go encodes absent tax-line slices as null while .NET validators require non-null collections. Preserve failed connected run/database; compare null versus empty-array payloads on the actual UDS host before changing the canonical wire owner. No completed fiscal claim.
FAIL-20260913-961 RESOLVED_DIAGNOSTIC_SCOPE: Python3.14 on this Windows build does not expose AF_UNIX; the auxiliary probe failed before sending any request and stopped its owned .NET fixture. Use the already proven Go/.NET UDS path and a direct wire regression, not another runtime or repeated unsuccessful Python probe.
FAIL-20260913-962 DIAGNOSED — corrected fiscal wire reaches LastAuthorized and the deliberate lost authorization response, leaving reconcile_required with voucher1. The subsequent real consult reaches PostgreSQL but FinishInvoice returns fiscal conflict. Preserve r2 database/logs and inspect the exact database constraint/transition before changing code; null-array correction already progressed the intended boundary.
FAIL962 cause proven: PostgreSQL invoice_provider_codes_check rejects JSON null. IPC copied an empty provider_codes array into a nil Go slice; FinishInvoice marshaled it directly. Reuse the existing fiscal-parameter owner normalization []string{} at persistence, preserve arrays in IPC responses, and add nil-code real PG coverage. No constraint weakened, no fabricated provider event code.
FAIL-20260913-963 OPEN — explicit ARCA generator-graph SCA found GHSA-g4vj-cjjj-v7hg in bundled NuGet.Packaging6.12.1 and NuGet.Protocol6.12.1 inside official svcutil8.0.0. Current fiscal application locks have no reported finding. Stop promoting/reusing that generator runtime; inspect the official advisory and fixed compatible package graph, then admit a declared adaptation or official patched source. Preserve successful fixtures/generated hashes and original tool; no Daybreak/libxml2 investigation.
FAIL-20260913-964 DIAGNOSED — NuGet patch quarantine verified fixed catalog length/SHA512 but the package identity/license-text assertion rejected before execution. Preserve all downloaded archives; inspect their exact nuspec and embedded notice set, then complete official commit-fixed license text if absent. Resume metadata inspection only, not acquisition.

FAIL-20260913-965 DIAGNOSED — OSV2.5.1 explicit lockfile auto-detection rejects the nonstandard svcutil-packages.lock.json basename before querying vulnerabilities. Preserve dotnet-fixed-osv.log. Store the graph projection at fiscal/svcutil/packages.lock.json so the actual admitted NuGet extractor is selected; no exclusion of the generator graph or inferred SCA PASS.

FAIL-20260913-966 DIAGNOSED — first adapted svcutil execution passes inventory/feed gates but fails loading NuGet.Common6.12.1.1 while the fixed composition contains6.12.5.1. Preserve generation-adapted-first and original binaries. Inspect the exact loader/bootstrap boundary before qualification; no fallback to vulnerable NuGet and no generated-client PASS.

V402 ARCA957 refinement: actual SDK verification exit0 proves metadata.contentHash equals the verified NuGet signed-package Content hash. The .nupkg.sha512 sidecar differs from that signed Content hash and remains a separate recorded observation. Raw SHA512 mismatch alone is not evidence of archive corruption; exact original SHA256 and113payload correspondence stay valid. Five current fiscal runtime packages also match their lock contentHash to verified signed Content hash. Original6.12.1 generator dependencies remain rejected963 independently.
V402964 RESOLVED_FIXED_LEGAL: nine fixed patch archives declare Apache-2.0 without embedded license text; same immutable NuGet.Client commit LICENSE and official Apache text retained, nine author/repository signatures PASS.965 RESOLVED_EXTRACTOR_CONTRACT: conventional svcutil/packages.lock.json gives actual12-manifest OSV2.5.1 exit0/zero findings.966 actual generation compatibility remains open.

FAIL-20260913-967 DIAGNOSED — isolated portable build reaches exact svcutil assembly but PowerShell native invocation returns no LASTEXITCODE for pinned dotnet --version under the builder's minimal environment. Preserve build-qualification; no compiler/generation ran in this attempt. Inspect Windows executable recognition in that environment, then resume only the unexecuted generation/publish steps.

FAIL-20260913-968 RESOLVED_BUILD_LAYOUT_GUARD — after explicit Windows PATHEXT/COMSPEC, dotnet starts correctly. The generation overlap gate then rejects the intended selected-project arca/fiscal/generated subtree. Admit only that exact absent subtree while retaining all other source/cache/tool/feed overlap checks; retain build-qualification-r2. Prior isolated generations already passed; actual integrated build remains to prove this layout.

FAIL967 refinement: r3 confirms native invocation and approved generated subtree work; NuGet ConfigurationDefaults initialization then requires Windows machine directory locators absent from the minimal environment. Supply only ProgramFiles, ProgramFiles(x86), ProgramData, SystemDrive; isolate APPDATA/LOCALAPPDATA under the owned build. Package sources remain the exact local feed, all locks enforced. Preserve r3; no inherited credential environment.

FAIL-20260913-969 DIAGNOSED — canonical pack assembly rejects CRLF serialization introduced by staging Python write_text for AUTHORED JSON/PowerShell glue. Preserve before bytes externally; normalize only AUTHORED candidate text to the canonical LF contract, keep all VERBATIM legal files byte-exact, and rebuild from the pack. No business algorithm changed; final artifact builds will consume the normalized exact sources.

FAIL-20260913-970 DIAGNOSED — exact four-profile reconstruction rejects stale GeneratedClientRoot materializer variables on the two newly selected credential integration project blocks. Their payloads already use the tested MSBuild property/default and contain no materializer placeholder. Remove only those two obsolete metadata declarations, preserve the failed metadata/receipts and reconstruct again; do not invent a project path or ask the user.

V402325: FAIL956/958/960/962/963/966/967/968/969/970 RESOLVED_LOCAL_INFRA with ARCA_CONNECTED_INFRA_V402. Actual portable CMS/UDS/PG fixture, null-array regressions, fixed signed NuGet cluster, offline fresh-cache generation, complete selected credential closure and exact4profile rebuilds. Failed logs/metadata/archives preserved.957 signed content/hash observations retained; no cache rewrite.1539 plaintext-URL finding resolved as verified UDS message syntax; no TCP provider transport. Global release remains T2810.

FAIL-20260913-971 RESOLVED_LOCAL_WRITE: Python text-mode write of PROJECT_ENGINEERING_CONTRACT returned EINVAL before changing its original bytes. Preserved backup, prepared exact successor externally and applied native File.WriteAllBytes; SHA256430cd51c7a89367a907c61c29147feb2ac1694da6ea1a30aaa40c534b6e57fb8 verified. Resume only unpublished cursor/readiness tail; no repeated pack/source publication.

FAIL-20260913-972 DIAGNOSED — actual gate ParseExact with DateTimeStyles.None interprets the fixed ZIP timestamp in Windows UTC-03, yielding315543600; builder correctly requires315532800. No full build attempted. Parse the unzoned release timestamp explicitly as UTC and test the actual gate assignment against the UTC epoch; preserve before bytes.

FAIL-20260913-972 RESOLVED_LOCAL: explicit UTC in canonical gate0.2.1, actual AST epoch regression and exact1653file reconstruction PASS. No failed full build was repeated; first actual final build starts with this fix.

FAIL-20260913-973 DIAGNOSED — current START_FRANCHISE Preflight catches stale expected profile count in VERIFY_LIBRARY: signed gate expects10 but selects11 after323. Preserve preflight log, derive every affected expected count from actual selected canonical blocks and reconcile the verifier inventory once. No product code or dependency change; running final builds remain valid.

FAIL973 counter correction: all56profile selections derived, only signed gate count changed10→11. Public Markdown962. Product source326 unchanged. Final Preflight rerun still required.

FAIL-20260913-974 DIAGNOSED — first actual full reference build completed but publish_reference_build rejected artifact/receipt correspondence. Preserve durable build-run-r326/build-1 and failed signed staging55c208a3325d4d6398bb414fd86b04ba; inspect exact error and source/receipt mapping before any rebuild. No signed release admitted.

FAIL974 cause: Next16.3.4 adds two generated route-type imports to next-env.d.ts. Builder copied that mutated workspace declaration into corresponding source; final publication correctly rejected one mismatch among1653. Preserve exact original source and separately retain only this precisely checked framework-derived declaration; reject every other source mutation.
FAIL-20260913-975 DIAGNOSED — preflight327 ran while failure974 was appended to canonical ledger, invalidating its checkpoint hash. Stop overlapping canonical mutations with full verification. Preserve log; final preflight/archive runs only after publication/checkpoint is quiescent. No product failure or gate bypass.

V402328 FAIL974 corrected in canonical CI0.1.7: exact original source, separately retained narrow Next derivation;16fixtures and1653actual source projection PASS. New current full builds/signature pending. FAIL973 count fixed, final verification pending; FAIL975 audit race prevented by scheduling final verification after publication.

FAIL-20260913-976 RESOLVED_CURSOR_ONLY — checkpoint328 rejected next_action exceeding800characters before writing its event. Shorten the cursor, retain detailed evidence/paths in the328 report and resume the same revision/event once. Product/build unchanged.

FAIL-20260913-977 OPEN — two actual builds328 and both deterministic ZIPs match; exact OSV2.5.1 source scan reports454 vulnerability records and prevents signing. Preserve both builds and staging8dceea47fa4a4254940759ad01bec2da. Inspect package identities, dependency selection and actual compiled runtime correspondence; no finding excluded or credential-only downgrade without evidence.

FAIL977 cause proven: official OSV report has454package records and zero vulnerabilities. PowerShell @($pkg.vulnerabilities) materialized one $null for every missing property, fabricating454empty-string findings. Both gate and verifier now enumerate nullable collections directly and reject empty vulnerability IDs. Actual454package report passes both new helpers; synthetic positive and malformed records remain detected. No dependency finding suppressed, source version changed or Daybreak branch opened.

V402329 FAIL977 RESOLVED_LOCAL_PARSER: actual OSV JSON454packages/0vulnerabilities; both actual helpers and17contract assertions PASS. Previous454claim corrected with preserved before/after evidence. New source329 includes the fixed release tools; current builds/signature still pending.

FAIL-20260913-978 OPEN — actual329signature/extraction/SCA pass, but SPDX reports only1root package while the all-packages OSV report contains454records. Signed content integrity is valid; dependency inventory completeness is not. Preserve329release and verification receipt, qualify OSV all-packages SPDX and enforce report-to-SBOM coverage in both builder/verifier before final acceptance. No vulnerability claim or production promotion.

FAIL-20260913-979 OPEN: extracted signed329 artifact native activation failed in web_local_stop junction creation before health. Fresh84migrations/replay passed; inspect missing link parents after ZIP extraction, preserve failed logs, return fix to canonical launcher. No production or final acceptance.

FAIL-20260913-980 CONTAINED_FIXTURE_HARNESS: signed329 runtime probe took219.814s; a synthetic120second OIDC token expired before post-rollback replay. Actual activation, order, manual rollback, redeploy and failed-candidate rollback succeeded. Authentication correctly rejected expired token401. Resume from committed pointer/existing database using newlyissued fixture tokens; do not weaken expiry or rerun completed migrations/rollback phases.

V402330 FAIL978/979 RESOLVED_LOCAL: full437SPDX records cover454OSV rows/335identities; ZIP-omitted parents recreated by launcher, actual extracted329API/web/order/rollback/recovery/native stop PASS. Source330full release remains pending.

FAIL-20260913-981 RESOLVED_METADATA:19CI notice blocks preserved real licenses/provenance but used a generic local-glue source label. Exact source receipt mapping already exists in SOURCE_NOTICES.json. Labels now bind that lock and exact payload SHA;CI0.1.9. Every materialized byte unchanged; current330two builds remain valid. No source re-acquisition or reputation claim.

V402 final local acceptance:384/385/532/804 contained or resolved in LIBRARY_FAIL_CONTAINMENT_V402.json;405historical provider import excluded,423/742Firefox teardown not claimed fixed,735actualnotice resolved.978/979SPDX/ZIP corrected;980fixture expiry correctly rejected and resumed;981metadata corrected. No open code blocker for selected LIBRARY_INFRASTRUCTURE. Final Preflight/library archive qualification remains the last operational check before COMPLETE.

FAIL-20260913-982 RESOLVED_CURSOR_ONLY: checkpoint332 rejected two evidence IDs sharing the same report path before writing its event. Business-acceptance evidence now references its existing exact T2802 final receipt; system-validation retains the A-G report. Product/source and readiness PASS unchanged. Resume same332event once, no business suite rerun.

FAIL-20260913-983 RESOLVED_SCOPE_RECORD: final evidence validator requires every active benchmark passed. BENCH01 is an unmeasured historical claim of reduced agent work/tokens, outside the user LIBRARY_INFRASTRUCTURE acceptance; it is not marked PASS. Full original and explicit NONE_WITH_REASON decision preserved in specs/library-maintenance/tasks.md; removed only from active benchmark list. BENCH02 and finite local load evidence retained;48required controls unchanged.

FAIL-20260913-984 RESOLVED_PORTABLE_DOC: final Preflight correctly rejected a personal home path in the public A-G guide. Public guide now uses relative reference and LOCALAPPDATA public-policy locator; the absolute local guide lives outside canonical distribution in START_REFERENCE_V402.md. No source, product, signature, dependency or test delta. Actual shared path checker scans1243public files before retry; do not weaken selector.

FAIL-20260913-985 DIAGNOSED: final Preflight334 passed structural/composition checks then rejected the HTTP metrics pack because its root wrapper retained8 while the admitted pack contains11 files. Original failed log and148.453second receipt preserved. Audit all mappable static materialization counts and the separate history plan before correcting the wrapper; product source330, signature and runtime unchanged.

FAIL-20260913-985 RESOLVED_WRAPPER:13mappable static assertions and separate history21 selection audited; only HTTP8→11 drift. Exact count check retained, no gate relaxed. Canonical root wrapper corrected; final frozen Preflight remains required.

FAIL-20260913-986 RESOLVED_TASK_METADATA: final checkpoint336 correctly refused before event336 because roadmap T2801/T2810 checkboxes remained unchecked although actual local readiness, signed release and portable335 NEW/EXISTING had passed. Synchronize those two checkboxes with existing proof; preserve qualified335archive and failed checkpoint. Rebuild only the library distribution after public task metadata changed; carry forward all unchanged executable checks and exact source330 builds. No domain, security, ARCA or operations rerun.

FAIL-20260913-987 RESOLVED_CURSOR_FORMAT: same336 checkpoint refused a bare64hex release digest; schema requires sha256: prefix. Actual archive hash unchanged. Correct the cursor representation and finalizer template; no new artifact or test claim. Failed logs preserved; no event appended until validation succeeds.
