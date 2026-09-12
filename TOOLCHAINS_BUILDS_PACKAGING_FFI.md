# Toolchains, Builds, Packaging & FFI — manual operativo

> **Estado:** núcleo general v1 auditado el 2026-08-21; fuentes, integridad documental y cruces materiales cerrados. Artifact/release/evidence contract ejecutable disponible en `ENGINEERING_EXECUTION_PLAYBOOK.md`; cada target real debe generar su evidencia.
> **Autoridad:** Reproducible Builds; Bazel; CMake; LLVM/GNU platform docs; Cargo/Rust; PyPA/CPython; Go; npm; SemVer; SLSA/Sigstore como frontera de supply chain.
> **Propósito:** convertir fuente versionada en artefactos reproducibles, verificables, instalables y compatibles, e integrar lenguajes sin violar ABI, ownership, errores ni concurrencia.

## Cómo usar este manual

```text
source + manifest + lock + toolchain + target + config
→ action graph hermético
→ compile/codegen/link/package
→ tests + ABI/API checks
→ SBOM/provenance/signature
→ publish immutable
→ install/upgrade/rollback
→ reproduce y verificar
```

“Compila en mi máquina” significa que el ambiente local contiene inputs no declarados. La solución no es documentar rituales manuales, sino hacer explícito y verificable el grafo de inputs y acciones.

## Procedencia

- `[SPEC]`: estándar, ABI o documentación oficial versionada.
- `[ACADEMIC]`: material de curso/texto identificado.
- `[CODE]`: repositorio/implementación oficial.
- `[PROD]`: práctica operativa atribuida.
- `[MEASURED]`: resultado reproducible en entorno declarado.
- `[IMPL]`: traducción a manifest, build, test, release o runbook.
- `[OPEN]`: decisión dependiente de plataforma/proyecto no fijada.

## 1. Contrato del artefacto

Antes del build:

```text
producto y entrypoints:
source revision y dirty policy:
targets: OS/arch/libc/runtime/accelerator:
API/ABI mínima soportada:
toolchain/sysroot/SDK:
dependency policy y lock:
features/profiles:
outputs y layout:
debug symbols/source maps:
tests y quality gates:
reproducibility target:
signing/provenance/SBOM:
install/upgrade/rollback:
support/EOL:
```

Una build es una función de todos esos inputs, no sólo del commit.

## 2. Fases y artefactos

```text
preprocess/code generation
→ compile/transpile
→ assemble/object/module
→ link/bundle
→ strip/symbol split
→ package/image/archive
→ attest/sign
→ publish
```

Por fase registrar comando efectivo, working directory, environment allowlist, inputs, outputs, exit status y tool version. No publicar outputs parciales.

### 2.1 Host, build y target

- **host:** donde corre la herramienta;
- **build machine:** término frecuente para la máquina que construye;
- **target:** donde correrá el artefacto.

En cross-compilation, code generators/build scripts suelen correr en host mientras objetos generados apuntan a target. Confundir ambos produce probes falsos y binaries no ejecutables.

### 2.2 Target contract

```text
triple/architecture
endianness y pointer width
ISA baseline/extensions
OS deployment target
libc/CRT y C++ runtime
calling convention/ABI
dynamic loader
SDK/sysroot
GPU/driver/runtime si aplica
```

No compilar accidentalmente para la CPU del builder si se distribuirá a hardware más antiguo.

## 3. Hermeticidad, determinismo y reproducibilidad

### 3.1 Conceptos distintos

| Propiedad | Significado |
|---|---|
| hermetic | acciones sólo leen inputs declarados/aislados |
| deterministic | mismos inputs producen mismo output |
| reproducible | terceros pueden recrear entorno y obtener output verificable |
| repeatable | mismo equipo/proceso repite resultado; garantía más débil |

**[SPEC]** Bazel define hermeticidad como aislamiento del host y versionado explícito de tools/dependencies. Esto habilita cache y ejecución paralela/remota; no elimina por sí solo timestamps, randomness o nondeterminism del tool.

### 3.2 Fuentes de variación

- reloj/timezone/locale;
- orden de filesystem/map/concurrency;
- paths absolutos/build user/hostname;
- random seeds/build IDs;
- archive metadata/permissions;
- tool/SDK/system libraries;
- network/latest tags;
- generated code y volatile VCS info;
- uninitialized bytes;
- compression/link order;
- environment variables.

### 3.3 Protocolo reproducible

1. clean checkout por digest;
2. toolchain y dependencies por digest/version;
3. network disabled o inputs pre-fetched verificados;
4. environment allowlist;
5. `SOURCE_DATE_EPOCH` donde aplique;
6. stable input/output order;
7. normalize archive metadata/paths;
8. build en dos entornos independientes;
9. comparar hashes y usar diffoscope/equivalente si diverge;
10. registrar recipe y evidencia.

Reproducible Builds enumera control de timestamps, locale, timezone, paths, ordering, randomness y environment como fuentes típicas.

## 4. Build graph correcto

Cada acción declara:

```text
inputs directos y transitivos
tool executable + config
outputs exclusivos
target/profile/features
environment permitido
resources/network policy
```

### 4.1 Incremental correctness

Probar:

- clean build = incremental build;
- null build ejecuta cero acciones innecesarias;
- tocar un input invalida consumidores correctos;
- borrar output lo reconstruye;
- cambiar feature/flag/toolchain invalida cache;
- dos targets no escriben el mismo output;
- generated files no modifican source tree.

Una cache hit incorrecta es corrupción, no optimización.

### 4.2 Remote cache/execution

La cache key incluye todos los inputs semánticos. Aislar credenciales y secretos; no almacenarlos en action outputs/logs. Validar trust, tenant isolation, digest y eviction. Comparar local/remote outputs.

### 4.3 Práctica pública Meta/Buck2 — incrementalidad sin perder identidad

**[CODE]** Buck2, commit `a3b535bb144294cbda6b04ae3b58f38c3f503fa2` (Apache-2.0/MIT), publica una arquitectura que separa parsing, configuración, análisis, action graph, ejecución y materialización. Cada acción obtiene un digest de su comando e inputs; un resultado remoto puede conservarse por digest y materializarse sólo cuando otro paso local o el usuario lo necesita. La lección no es copiar Buck2, sino separar tres identidades que suelen mezclarse:

```text
identidad de análisis = regla + configuración + dependencias semánticas
identidad de acción   = comando + tools + environment permitido + inputs
identidad de artefacto = bytes + metadata que el contrato considera relevante
```

El daemon y el motor incremental reutilizan computaciones entre invocaciones; ejecución local, remota o en carrera decide dónde obtener el resultado, no qué significa. Dependencias descubiertas por contenido y acciones incrementales siguen necesitando límites declarados: si un depfile, un estado previo o un output dinámico omite una dependencia, el sistema puede producir un cache hit coherente con una clave incorrecta.

**Puerta operacional:**

1. comparar resultado clean, daemon/incremental, cache local, cache remota y remote execution por digest o equivalencia explícita;
2. mutar individualmente source, tool, flag, environment permitido, configuración de target y dependencia dinámica, verificando la invalidación exacta;
3. ejecutar una acción en sandbox sin filesystem/network implícito y auditar accesos no declarados;
4. probar cold/warm/null build, critical path, bytes CAS, latencia de materialización, hit útil y hit rechazado;
5. si existe estado incremental, escribirlo de forma transaccional: invalidarlo antes de mutar el output y publicar el nuevo estado sólo después de completar la actualización;
6. aislar daemon/cache por proyecto, configuración y tenant; probar restart, cache corrupta/evictada y cambio de versión del build system;
7. no aceptar local/remote racing si las dos ejecuciones pueden producir bytes o efectos externos distintos.

**[IMPL]** La aceleración correcta conserva semántica. Un build más rápido con invalidación incompleta, outputs no herméticos o materialización ambigua es corrupción diferida, no optimización.

### 4.4 Código Bazel — una cache hit es una afirmación verificable

**[PROD]** Código upstream fijado en `4ecfa6bc4d7c792772dedd8f630e2d60cb018c06` (Apache-2.0). `ActionCacheChecker` compara una identidad construida con action key, inputs y metadata, environment permitido, execution salt, permisos y outputs; un output ausente/no confiable invalida la entrada. La ruta remota distingue miss esperado, fallo de cache/capabilities y resultado reutilizable; si no puede reconstruir con fidelidad un fallo previo, vuelve a ejecutar. Los tests corrompen el índice de cache y exigen build exitoso con warning/rebuild, mientras el sandbox hermético intenta accesos absolutos, dependencias omitidas, symlink escape y modificación de inputs.

Gates adicionales al protocolo Buck2:

- enumerar formalmente qué bytes/config/environment/permisos/toolchain forman la action key y probar cada dimensión mediante mutación;
- tratar cache corrupta, blob ausente, digest mismatch, partial download/upload y eviction como miss/rebuild seguro o fallo explícito, nunca como output válido parcial;
- ejecutar acciones no confiables sin host filesystem/network implícito y detectar modificación de inputs, outputs fuera del árbol y symlink escape;
- no cachear fallo/resultado si no puede reproducirse con el mismo diagnóstico y contract; separar deduplicación in-flight de persistencia de resultados;
- medir `hit útil`: hit rate sin bytes, latencia, materialización y correctness no decide si cache/remoto acelera el critical path.

## 5. Manifests, locks y resolución

### 5.1 Roles

| Artefacto | Función |
|---|---|
| manifest | intención/rangos/features/metadata |
| lockfile | resolución concreta y, según ecosistema, integridad/source |
| registry/index | discovery/distribución |
| package | contenido publicado |
| SBOM | inventario del artefacto resultante |
| provenance | cómo/dónde/con qué inputs se produjo |

Un lockfile no fija OS packages, compiler, SDK, build scripts remotos o mutable URLs salvo que el sistema los incluya explícitamente.

### 5.2 Dependency policy

- direct/transitive y runtime/build/dev;
- source/registry allowlist;
- exact integrity/digest;
- license/security policy;
- duplicate versions y feature unification;
- minimum supported runtime/toolchain;
- update cadence y rollback;
- vendoring/mirror/offline plan;
- yanked/deleted/compromised package response.

### 5.3 Resolver tests

```text
fresh resolve/install
locked/offline install
minimum-supported set donde ecosistema lo permita
latest-compatible set
multiple OS/arch/runtime
conflicting/optional features
registry outage/mirror
dependency removal/yank
```

## 6. Versionado y compatibilidad

**[SPEC]** SemVer 2.0.0 requiere declarar una API pública; major rompe compatibilidad, minor agrega compatible y patch corrige compatible. Una versión publicada no se modifica.

La API pública incluye más que firmas:

- wire/file/schema formats;
- CLI flags/output/exit codes;
- configuration/env;
- import paths y symbols;
- error types/messages si consumidores dependen;
- performance/resource guarantees documentadas;
- plugin/extension interfaces;
- observable order/side effects.

SemVer comunica intención, no prueba compatibilidad. Automatizar API/ABI diff y consumer tests.

### 6.1 Deprecación

```text
announce + replacement
telemetry de uso
compatibility window
migration tooling/docs
warning sin romper automation
removal en release correcta
support y rollback
```

## 7. C/C++ con CMake

CMake 4.4.2 era la documentación observada. Modelar targets, no listas globales de flags.

### 7.1 Target contract

- sources y generated sources;
- `PUBLIC/PRIVATE/INTERFACE` includes, definitions y libraries;
- compile features/standard;
- warnings/sanitizers por target/profile;
- PIC/visibility;
- install/export/package config;
- transitive usage requirements.

Evitar `include_directories()`/flags globales que filtran a targets no relacionados.

### 7.2 Presets y toolchains

Versionar presets no secretos para configure/build/test. Toolchain file fija compiler, sysroot, target y search policy. No usar runtime probes (`try_run`) sobre target no ejecutable sin emulator/result predefinido; `try_compile` prueba compilación, no capacidad runtime.

### 7.3 Libraries

| Tipo | Consecuencia |
|---|---|
| static | código copiado; licencia/size/update |
| shared | loader/ABI/dependency discovery |
| object | objetos reutilizados sin interfaz de link completa |
| interface | usage requirements, sin objeto |

No mezclar runtimes/allocators incompatibles en boundaries sin contrato.

## 8. Compile y link

### 8.1 Translation units y ODR

Headers, macros, compile flags y generated config forman input. Diferencias pueden violar One Definition Rule aunque el linker no falle. Evitar ABI types condicionados por macros distintos entre library/consumer.

### 8.2 Link

- symbol definition/reference y visibility;
- static archive order/group semantics;
- dead stripping/LTO;
- weak/interposable symbols;
- version scripts/export lists;
- duplicate runtimes;
- RPATH/RUNPATH/install names;
- delay/load-time resolution;
- debug symbol mapping.

El orden de argumentos puede afectar link; Cargo documenta que incluso el orden de instrucciones de `build.rs` puede propagarse al linker.

### 8.3 Loader

Testear en ambiente limpio, no sólo builder. Verificar dependencia transitiva, search path, architecture, code signing, relocatability y error al faltar/ser incompatible una library. Nunca depender de current directory o PATH no controlado para cargar código.

## 9. ABI

API es contrato fuente; ABI es contrato binario:

- calling convention;
- name mangling/linkage;
- object layout/alignment/vtables;
- size/representation de tipos;
- exception/unwind;
- runtime/allocator;
- symbol versioning/visibility.

**[SPEC]** Itanium C++ ABI documenta layout, vtables, mangling, RTTI y exception ABI para plataformas que lo adoptan. Windows/MSVC y otras plataformas tienen ABIs distintas: fijar toolchain/target reales.

### 9.1 Estabilizar boundary

Preferir C ABI estrecha:

- opaque handles;
- fixed-width scalars;
- pointer + length;
- explicit status/error;
- create/destroy del mismo lado;
- version/size field en structs extensibles;
- no STL/Rust/Python objects cruzando;
- no exceptions/panics cruzando por defecto.

No asumir layout de `bool`, enum, bitfield, long, wchar o packed struct entre compilers/languages.

### 9.2 ABI gate

- exported symbol list;
- size/alignment/offset assertions;
- calling convention;
- baseline libc/OS;
- old consumer + new library;
- new consumer + old library cuando prometido;
- debug/release/sanitizer variants no mezcladas;
- architecture matrix.

## 10. FFI contract

### 10.1 Obligaciones

```text
ownership: quién libera y con qué allocator
lifetime: cuánto vive pointer/view/callback context
mutability/aliasing
thread affinity y thread safety
encoding y nullability
errors/exceptions/panics
cancellation
callbacks y shutdown
blocking/reentrancy
layout/alignment
version/feature negotiation
```

### 10.2 Memory

Regla segura: quien asigna libera, salvo allocator API explícita compartida. Para buffers:

```c
struct buffer_v1 {
  uint32_t struct_size;
  const uint8_t* data;
  size_t len;
  void* owner;
  void (*release)(void* owner);
};
```

El ejemplo expresa un patrón, no una ABI universal. Documentar si es borrowed/owned y cuándo invocar `release`.

### 10.3 Errors

No usar global `last_error` sin thread contract. Retornar status estable y permitir obtener detalle owned/copy. Separar programmer error, invalid input, transient y fatal. Ningún error language-specific cruza sin ABI soportada.

### 10.4 Callbacks

- context pointer estable;
- unregister/drain antes de destruir;
- ninguna callback después de cierre confirmado;
- thread de invocación documentado;
- reentrancy/deadlock;
- callback lenta/backpressure;
- panic/exception capturada en boundary.

El Rustonomicon subraya que callbacks asíncronas no deben apuntar a objetos Rust destruidos y que unwind sólo cruza ABIs que lo permiten; de otro modo puede abortar o ser undefined behavior.

## 11. Rust/C/C++

### 11.1 Rust export/import

- `extern "C"`/ABI correcto;
- `#[repr(C)]` para layouts acordados;
- opaque pointer para tipos internos;
- no `String`, `Vec`, trait object o Rust enum sin representación contractual;
- `catch_unwind`/abort policy;
- `Send/Sync` y callback threads;
- pinning/self-reference;
- allocator/free API.

`unsafe` wrapper demuestra preconditions antes de exponer API segura.

### 11.2 Cargo

Cargo workspaces comparten lock y output; features son aditivas según resolución y pueden cambiar código/ABI. `build.rs` ejecuta código durante build: declarar `rerun-if-*`, escribir sólo en `OUT_DIR`, distinguir host/target y auditar network/environment.

Profiles, target features, panic strategy y LTO forman parte del artifact identity.

## 12. CPython extensions

**[SPEC]** Una extensión CPython es shared library compatible que exporta init function. Multi-phase initialization es el mecanismo recomendado; modules deben tratar multiple instances/subinterpreters según capacidad.

### 12.1 Contract

- Python versions/ABIs y platform tags;
- limited API/Stable ABI o build por versión;
- refcount new/borrowed/stolen;
- GIL/free-threaded/subinterpreter assumptions;
- buffer protocol lifetime;
- exception mapping;
- release GIL sólo sin tocar Python objects;
- module state per interpreter;
- wheel matrix/repair/audit.

Stable ABI reduce rebuilds al limitar API; no garantiza que código propio o libraries nativas sean compatibles.

### 12.2 Wheels

Wheel tag expresa interpreter/ABI/platform. Probar instalación en imagen mínima por tag, import, symbols, bundled libraries y license. No renombrar wheel para fingir compatibilidad.

## 13. Python packaging

PyPA especifica `pyproject.toml` con `[build-system]`, `[project]` y `[tool]`. Separar:

- library distribution: ranges compatibles y wheel/sdist;
- application environment: resolución bloqueada;
- build requirements: aislados;
- optional/development groups.

`pylock.toml` (PEP 751) es formato estandarizado para instalaciones reproducibles observado en 2026; verificar soporte real del tool elegido.

Gate:

- build sdist y wheel en isolation;
- instalar wheel, no source tree;
- test sdist→wheel;
- metadata/license/readme;
- package data/imports;
- no secretos/local paths;
- hashes/index policy;
- TestPyPI/staging antes de publish;
- versión immutable.

## 14. JavaScript/npm

`package.json` expresa manifest; `package-lock.json` describe árbol concreto para reproducibilidad del install con npm compatible.

Revisar:

- `npm ci` en CI, lock sin mutación;
- `engines`, package manager y Node version;
- `exports`/`imports`, ESM/CJS;
- browser/node conditional exports;
- lifecycle scripts como ejecución de código;
- bundled/native dependencies;
- files publicados con dry run;
- peer/optional dependencies;
- source maps y minification determinista;
- registry/integrity/provenance.

Lockfile no protege contra script malicioso de una dependency ya admitida.

## 15. Go modules

Go define module como conjunto de packages versionados/distribuidos juntos, con path y `go.mod`; `go.sum` verifica contenido observado. Major `v2+` usa suffix de module path salvo excepciones históricas.

Gate:

- `go mod tidy` sin diff;
- module proxy/private policy;
- checksum DB/offline/vendor según threat model;
- `go version -m`/build info;
- CGO cross-toolchain si aplica;
- reproducible flags/paths;
- compatibility por import path;
- test con minimum Go version declarada.

CGO reintroduce C compiler, libc, loader y ABI al artifact contract.

## 16. Code generation

Generated code tiene source of truth, generator version y command. Elegir:

- commit generated output si consumidores/build no tienen generator o review del diff aporta;
- generar en build si hermeticidad y tiempo están controlados.

Siempre:

- deterministic;
- no timestamps/banner volátil;
- stable ordering;
- schema compatibility;
- regenerate check en CI;
- generator input/output ownership;
- safe parser/templates;
- license/provenance.

No editar output manualmente.

## 17. Monorepo, multirepo y workspaces

Decidir por atomic changes, ownership, build graph, release cadence y access control.

| Monorepo | Multirepo |
|---|---|
| refactoring/visibility/atomicidad | isolation/autonomy/access separados |
| exige scalable build/ownership | exige contract/version coordination |

No copiar herramientas Google sin su escala ni resolver coordination. No crear repos por microservice automáticamente. Workspaces no implican una única release.

## 18. CI pipeline de build

```text
checkout immutable
→ verify source/signature
→ resolve locked inputs
→ build hermetic
→ unit/static/sanitizers
→ integration/ABI/package tests
→ reproduce/compare según riesgo
→ SBOM/provenance/sign
→ publish staging
→ promote same digest
```

No rebuild entre staging y production. Promover mismo digest con policy/evidence.

### 18.1 Matrix

- OS/arch/libc;
- compiler/toolchain supported range;
- debug/release;
- feature combinations importantes;
- runtime/interpreter versions;
- static/shared;
- sanitizer/coverage;
- cross-compile + target smoke test.

Evitar combinatoria completa mediante pairwise/risk, pero probar cada artifact publicado.

## 19. Supply chain

Para cada release:

- commit/tag y reviewer policy;
- builder identity y isolated runner;
- dependency integrity;
- SBOM por artifact;
- provenance con inputs/recipe;
- signature/attestation;
- vulnerability/license policy;
- secret-free logs/artifacts;
- registry immutability;
- revocation/incident procedure.

SBOM sin artifact digest o provenance no identifica qué se desplegó. Firma prueba key/identity, no calidad; verificar policy, transparency/freshness y trust root.

## 20. Artifact types

### 20.1 Archives/installers

Normalizar paths, permissions, uid/gid, timestamps y order. Defenderse de path traversal/symlinks al extraer. Installer declara privileges, locations, services, config preservation, upgrade y uninstall.

### 20.2 Containers

- base por digest;
- multi-stage y minimal runtime;
- non-root/read-only cuando viable;
- no compilers/secrets en final;
- architecture manifest correcto;
- entrypoint/signal/health;
- SBOM/provenance;
- same digest promotion.

Container no corrige binary ABI incompatible con kernel/CPU/driver.

### 20.3 Plugins

Plugin API/ABI, discovery, version negotiation, isolation, permissions, lifecycle/unload y failure containment. Código plugin es código ejecutable; firmarlo no vuelve confiable su comportamiento.

## 21. Release engineering

### 21.1 Release manifest

```text
version + commit
artifact digests/platforms
toolchain/build recipe
dependencies/SBOM
provenance/signatures
compatibility/migrations
known issues
install/upgrade/rollback
support window
```

### 21.2 Rollback

Binary rollback puede fallar si cambió schema/data/config/protocol. Probar mixed version y downgrade o declarar forward-only con recovery. Guardar artifacts y symbols dentro de retention operacional.

### 21.3 Hotfix

Derivar de revision conocida, cambiar mínimo, ejecutar gates, nueva versión immutable, documentar divergence/merge-back. Nunca reemplazar bytes bajo la misma versión/tag.

## 22. Testing del paquete real

1. build limpio;
2. instalar artifact publicado/staging;
3. ejecutar desde directorio ajeno;
4. usuario sin privilegios;
5. máquina sin source/dev dependencies;
6. offline/runtime network policy;
7. smoke/API/CLI;
8. upgrade/downgrade/uninstall;
9. dependency/library missing/corrupt;
10. locale/timezone/path con espacios/Unicode;
11. verify digest/signature/SBOM;
12. reproduce desde recipe.

Los tests del source tree no validan package metadata, loader paths ni archivos omitidos.

## 23. Debuggability

Conservar mapping release→symbols/source/build ID. Separar symbols si tamaño/privacidad exige, pero no perderlos. Deterministic paths/source maps, frame pointers según estrategia, crash dumps redacted y exact binary/container digest.

Optimized/LTO builds pueden alterar debugging; reproducir con artifact real antes de suponer equivalencia debug.

## 24. Performance del build

Medir:

- clean/incremental/null build;
- critical path/action count;
- cache hit local/remote;
- download/upload bytes;
- CPU/memory/I/O;
- codegen/link/package time;
- flaky/retry rate;
- developer wait p50/p95.

Optimizar graph/dependencies/hermeticity antes de más hardware. Unity builds/PCH/LTO mejoran ciertos tiempos y pueden afectar isolation, memory o incremental behavior: medir.

## 25. Anti-patrones

- curl latest durante build;
- version/tag mutable;
- PATH/system SDK implícito;
- generated output con timestamps;
- lockfile ignorado en CI;
- publish desde laptop;
- build scripts con network/secrets;
- global compiler flags;
- C++ objects/exceptions cruzando ABI no fijada;
- liberar memoria con allocator del otro runtime;
- callback después de destroy;
- wheel/container probado sólo en builder;
- rebuild para producción;
- SemVer sin API declarada;
- SBOM sin digest;
- cache compartida no autenticada.

## 26. Gate profesional

- [ ] Host/target/toolchain/sysroot fijados.
- [ ] Manifest/lock/source integrity.
- [ ] Action inputs/outputs herméticos.
- [ ] Clean/incremental/null equivalence.
- [ ] Reproducibility/variance policy.
- [ ] API/ABI/format compatibility.
- [ ] FFI ownership/lifetime/error/thread contract.
- [ ] Package real instalado y testeado.
- [ ] OS/arch/runtime/feature matrix.
- [ ] SBOM/provenance/signature.
- [ ] Immutable publish/promote same digest.
- [ ] Upgrade/mixed-version/rollback.
- [ ] Symbols/debug evidence y support/EOL.
- [ ] Build performance/cost measured.

## 27. Contrato para Codex

```text
Producto/entrypoints:
Lenguajes y repos/workspaces:
Host/targets/toolchains:
Dependencies/registries:
API/ABI/FFI:
Artifacts y distribución:
CI/security/reproducibility:

Entregar:
1) artifact y target contract;
2) build graph/manifests/locks;
3) hermetic/reproducible inputs;
4) compile/link/package/install;
5) FFI ownership/errors/threads con wrappers seguros;
6) compatibility y platform matrix;
7) tests del artifact real;
8) SBOM/provenance/sign/publish;
9) upgrade/rollback/debugging;
10) benchmark, riesgos y runbook.

No usar latest ni dependencias del host implícitas.
No cruzar exceptions/panics/allocators sin contrato.
No publicar una versión mutable.
No promover un rebuild distinto.
No confundir lock, SBOM, provenance y firma.
```

## 28. Suites de aceptación

| Suite | Inyección | Aceptación |
|---|---|---|
| `hermetic_clean` | host sin tools/libs extras | mismo build desde inputs declarados |
| `reproduce` | dos paths/timezones/users/builders | hashes iguales o variance explicada |
| `incremental` | cambios por tipo de input | invalidación exacta; null build vacío |
| `resolver` | offline/yank/conflict/mirror | locked install verificable o fallo seguro |
| `abi_matrix` | old/new lib/consumer y platforms | compatibilidad prometida, break detectado |
| `ffi_lifecycle` | null/error/callback-after-close/panic/thread | sin UAF/leak/unwind UB/deadlock |
| `package_install` | máquina mínima/path Unicode/no source | install/import/run/uninstall correctos |
| `release_chain` | tamper artifact/SBOM/attestation | verification bloquea y preserva evidencia |
| `rollback` | binary+config+schema mixed | camino soportado o forward recovery probado |

## 29. Capstones

### A. Biblioteca C++ + Rust + Python

C ABI opaca, CMake/Cargo, Python wheel. Probar alloc/free, errors, callbacks, threads, stable symbols, Linux/Windows y reproducible artifacts.

### B. CLI multi-plataforma

Go o Rust con packages/installers/container, signed release, SBOM/provenance, update/rollback y symbols.

### C. Monorepo polyglot

Bazel o graph equivalente para C++/Python/TS, codegen schema, remote cache, target matrix y hermeticity audit.

## 30. Inventario de autoridades

| Fuente | Cobertura |
|---|---|
| Reproducible Builds | variance, environment, deterministic output y verification |
| Bazel | hermetic action graph, sandbox/cache/remote execution |
| CMake 4.4.2 docs | targets, generators, presets, libraries, tests, install/export/dependencies |
| Cargo/Rustonomicon | workspaces, resolution/features/profiles/build scripts, FFI/unwind |
| PyPA/CPython 3.14 docs | pyproject/pylock/build/publish, extensions/C API/Stable ABI |
| Go modules | module paths, versions, resolution, proxy/checksum |
| npm v11 docs | package-lock y install tree |
| SemVer 2.0.0 | public API y version semantics |
| Itanium C++ ABI | layout, mangling, vtables, RTTI/unwind en targets aplicables |

Las versiones observadas no sustituyen la matrix del proyecto. Estándares/implementaciones pueden evolucionar; pin y verificar.

## 31. Fuentes oficiales verificadas

- Reproducible Builds docs: <https://reproducible-builds.org/docs/>
- Bazel hermeticity: <https://bazel.build/basics/hermeticity>
- Bazel fijado `4ecfa6bc4d7c792772dedd8f630e2d60cb018c06`: [ActionCacheChecker](https://github.com/bazelbuild/bazel/blob/4ecfa6bc4d7c792772dedd8f630e2d60cb018c06/src/main/java/com/google/devtools/build/lib/actions/ActionCacheChecker.java), [RemoteSpawnCache](https://github.com/bazelbuild/bazel/blob/4ecfa6bc4d7c792772dedd8f630e2d60cb018c06/src/main/java/com/google/devtools/build/lib/remote/RemoteSpawnCache.java), [corrupted-cache test](https://github.com/bazelbuild/bazel/blob/4ecfa6bc4d7c792772dedd8f630e2d60cb018c06/src/test/java/com/google/devtools/build/lib/buildtool/CorruptedActionCacheTest.java) y [hermetic sandbox tests](https://github.com/bazelbuild/bazel/blob/4ecfa6bc4d7c792772dedd8f630e2d60cb018c06/src/test/shell/bazel/bazel_hermetic_sandboxing_test.sh).
- CMake tutorial: <https://cmake.org/cmake/help/latest/guide/tutorial/index.html>
- Cargo reference: <https://doc.rust-lang.org/cargo/reference/>
- Cargo build scripts: <https://doc.rust-lang.org/cargo/reference/build-scripts.html>
- Rustonomicon FFI: <https://doc.rust-lang.org/nomicon/ffi.html>
- PyPA guide/specifications: <https://packaging.python.org/>
- `pyproject.toml`: <https://packaging.python.org/en/latest/specifications/pyproject-toml/>
- `pylock.toml`: <https://packaging.python.org/en/latest/specifications/pylock-toml/>
- CPython extension modules: <https://docs.python.org/3/c-api/extension-modules.html>
- Go Modules Reference: <https://go.dev/ref/mod>
- npm package lock: <https://docs.npmjs.com/cli/v11/configuring-npm/package-lock-json/>
- SemVer 2.0.0: <https://semver.org/spec/v2.0.0.html>
- Itanium C++ ABI: <https://itanium-cxx-abi.github.io/cxx-abi/abi.html>
- Buck2, snapshot auditado: <https://github.com/facebook/buck2/tree/a3b535bb144294cbda6b04ae3b58f38c3f503fa2>
- Buck2 architecture: <https://github.com/facebook/buck2/blob/a3b535bb144294cbda6b04ae3b58f38c3f503fa2/docs/concepts/architecture.md>
- Buck2 incremental actions: <https://github.com/facebook/buck2/blob/a3b535bb144294cbda6b04ae3b58f38c3f503fa2/docs/rule_authors/incremental_actions.md>

## 32. Cruces iniciales

| Manual | Frontera |
|---|---|
| `SOFTWARE_ARCHITECTURE_SYSTEM_DESIGN.md` | artifact boundaries, build/buy, compatibility y evolución |
| `COMPUTER_SYSTEMS_PERFORMANCE_LOW_LATENCY.md` | compiler/ISA/link/runtime/performance |
| `SOFTWARE_BACKEND_API_ENGINEERING.md` | package/service contracts, CI/release |
| `GPU_ACCELERATED_COMPUTING.md` | CUDA/driver/arch/toolkit kernels y extensions |
| `SECURITY_SRE_CLOUD_INFRASTRUCTURE.md` | SLSA/Sigstore/SBOM, CI IAM y deployment |
| `FRONTEND_PRODUCT_ENGINEERING_UX.md` | TS/npm/bundling/source maps/browser targets |
| `DATA_ENGINEERING_ANALYTICS.md` | DAG/jobs/codegen/connectors y reproducibilidad |
| `AI_ENGINEERING_MASTER_MAP.md` | native ML extensions, model/runtime packaging |
| `NATIVE_MOBILE_DESKTOP_ENGINEERING.md` | SDK/OS/arch matrix, native modules, firma, store/package y updater |

Cruces materiales cerrados el 2026-08-21: arquitectura, sistemas, backend, GPU, seguridad, frontend, datos, IA y aplicaciones nativas incorporan sus obligaciones. Este manual gobierna artifact identity, build/package/ABI/FFI y release chain; cada receptor conserva semántica de dominio. Para nativo, deployment target, SDK, entitlements/capabilities, signing identity, nested binaries, store metadata y update path pertenecen al artifact contract; probar el package instalado en targets reales, no sólo el build de desarrollo. No se añadieron enlaces ceremoniales donde no cambia una decisión.

## 33. Puerta de cierre

- [x] Artifact/target/build graph contract.
- [x] Hermeticity/reproducibility/cache.
- [x] Manifests/locks/versioning.
- [x] CMake/compile/link/loader/ABI.
- [x] FFI/Rust/CPython.
- [x] Python/npm/Go packaging.
- [x] CI/supply chain/release/testing/debugging.
- [x] Codex/suites/capstones/fuentes.
- [x] Auditoría Markdown/referencias: fences, encabezados, caracteres y referencias internas.
- [x] Cruces recíprocos materiales y arbitraje.
- [x] Contrato ejecutable de artifact/release/evidence y capstones disponible en `ENGINEERING_EXECUTION_PLAYBOOK.md`; builds de targets reales no se presumen ejecutados.

El conocimiento general v1 queda auditado. No declarar suites ejecutadas sin artefactos y resultados reales.
