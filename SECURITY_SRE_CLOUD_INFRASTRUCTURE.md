# Seguridad, SRE, cloud e infraestructura — manual operativo

> **Estado:** núcleo general v1 auditado y cruces recíprocos cerrados el 2026-08-21. Corpus verificado: Stanford CS155 Spring 2026; Dan Boneh y Victor Shoup; NIST CSF 2.0, SSDF 1.1, SP 800-207/207A, SP 800-63-4 y SP 800-61r3; OWASP Top 10:2025, ASVS 5.0.0 y API Security 2023; Google SRE y *Building Secure and Reliable Systems*; SLSA 1.2, in-toto 1.0 y Sigstore; OpenTelemetry, Prometheus y Kubernetes 1.36; fuentes completas en §27. Incluye cruce recíproco con frontend.
> **Propósito:** diseñar, construir, desplegar y operar sistemas cuya seguridad, fiabilidad, costo y capacidad sean propiedades verificables, no configuraciones asumidas.
> **Alcance:** aplicación, identidad, datos, host, red, supply chain, cloud, contenedores, Kubernetes, observabilidad, incidentes y recuperación. Cumplimiento regulatorio específico requiere jurisdicción y asesoramiento especializado.

## 0. Cómo usar este manual

Antes de programar o provisionar, completar §1 y §2. Después usar:

- §3–§8 para arquitectura y desarrollo seguro;
- §9–§12 para cloud, infraestructura como código y Kubernetes;
- §13–§19 para SRE, telemetría, entrega, incidentes y recuperación;
- §20–§24 para gates, instrucciones a Codex y suites ejecutables;
- §25 para arbitrar contradicciones con los demás manuales;
- §27 para verificar procedencia, versión y límites.

Etiquetas:

- **[SPEC]** estándar, RFC o documentación oficial versionada;
- **[ACADEMIC]** curso o texto académico;
- **[PROD]** experiencia publicada por operadores expertos;
- **[CODE]** comportamiento comprobable en implementación oficial;
- **[MEASURED]** resultado reproducible en entorno declarado;
- **[IMPL]** traducción operativa propia de esta biblioteca;
- **[OPEN]** decisión todavía no resuelta; nunca tratarla como garantía.

Reglas de precedencia:

1. La especificación/version del componente decide semántica.
2. El threat model decide qué ataques importan.
3. El contrato del producto decide impacto y prioridades.
4. La evidencia operativa decide si el control funciona y es sostenible.
5. Un framework de cumplimiento organiza evidencia; no prueba seguridad.

## 1. Contrato del sistema y del riesgo

No existe “seguro”, “altamente disponible” ni “cloud-ready” sin contexto. Registrar:

```yaml
system:
  owners: []
  users_and_tenants: []
  critical_user_journeys: []
  environments: [dev, staging, prod]
  data_classes: []
  trust_boundaries: []
  dependencies_and_providers: []
  regions_and_failure_domains: []
security:
  assets: []
  adversaries: []
  entry_points: []
  attacker_capabilities: []
  unacceptable_outcomes: []
  required_controls_and_evidence: []
reliability:
  slis: []
  slos: []
  error_budget_policy: ""
  capacity_and_quotas: []
  rto: ""
  rpo: ""
operations:
  deploy_and_rollback: ""
  oncall_and_escalation: []
  backup_and_restore: ""
  retention_and_deletion: ""
open_questions: []
```

### 1.1 Activos e impactos

**[IMPL]** Un activo no es sólo un archivo. Incluye:

- datos, claves, credenciales, modelos, artefactos y propiedad intelectual;
- identidad y capacidad de actuar: cuenta, role, token, signing authority;
- integridad de código, configuración, pipeline, telemetría y decisiones;
- disponibilidad, cuota, presupuesto cloud y capacidad de recuperación;
- confianza de usuario, seguridad física y obligaciones legales.

Para cada resultado adverso, estimar impacto y reversibilidad en confidencialidad, integridad, disponibilidad, privacidad, dinero y personas. No multiplicar números arbitrarios para fingir precisión; usar categorías con criterios observables y propietario de aceptación.

### 1.2 Adversario y failure model

**[ACADEMIC]** Stanford CS155 organiza ataques y defensas sobre sistema, web, criptografía/HTTPS, red, DoS, privacidad, supply chain y sistemas de IA. La pregunta transferible es siempre qué capacidad posee el atacante y qué supuesto rompe.

Separar:

- error humano, bug y fallo aleatorio;
- abuso por usuario autenticado;
- credencial robada o sesión secuestrada;
- dependency/build/provider comprometido;
- workload, nodo, cuenta cloud o administrador comprometido;
- insider y colusión;
- adversario de red pasivo o activo;
- agotamiento económico o de recursos;
- desastre regional y pérdida del control plane.

**[IMPL]** No usar MITRE ATT&CK como lista universal. Mapear sólo técnicas plausibles a activos, superficies, señales, prevención y respuesta. Un control sin ataque o fallo declarado puede ser costo sin reducción de riesgo.

### 1.3 TCB y raíces de confianza

**[ACADEMIC]** El *Trusted Computing Base* (TCB) es el conjunto mínimo de componentes, claves y decisiones cuya corrupción puede violar una invariante. Inventariarlo incluye CPU/firmware, boot chain, kernel/hypervisor, identity/policy plane, builder, signing keys, recovery tooling y operadores privilegiados según el sistema. “Confiable” significa que **debe** comportarse correctamente, no que se haya demostrado seguro.

Minimizar el TCB, fijar sus versiones/owners, endurecer sus interfaces, actualizarlo y probar cómo se reconstruye desde raíces verificadas. Una enclave puede excluir OS/RAM/hypervisor de cierta confidencialidad, pero incorpora hardware, firmware, attestation chain y código interno al contrato; no convierte la aplicación en correcta ni elimina side channels.

## 2. Invariantes de seguridad y fiabilidad

**[PROD]** *Building Secure and Reliable Systems* une ambas disciplinas mediante invariantes, entendibilidad, least privilege, resilience, investigación, recuperación y aprendizaje.

Ejemplos que deben convertirse en tests/policies:

```text
ninguna identidad accede fuera de tenant, recurso y acción autorizados
ningún artefacto llega a prod sin provenance e identidad de build aceptadas
ningún secreto aparece en repo, imagen, plan, estado, log o trace
ningún ACK durable se emite antes del punto de persistencia declarado
ningún resultado tardío o de dependencia fallida se publica
ningún restore se considera válido sin prueba funcional e integridad
ningún deploy puede consumir error budget sin rollback observable
ningún control crítico falla abierto salvo decisión explícita y probada
```

### 2.1 Principios

- **Least privilege:** mínima identidad, acción, recurso, contexto y duración.
- **Deny by default:** ausencia de regla no concede.
- **Complete mediation:** autorizar cada operación sensible sobre el objeto real.
- **Separation of duties:** crear, aprobar, firmar y desplegar no dependen de una única credencial cuando el riesgo lo exige.
- **Defense in depth:** capas independientes contra el mismo resultado adverso.
- **Blast radius:** cuentas/proyectos, tenants, regiones, claves y pools son failure/security domains deliberados.
- **Secure by default:** el camino fácil no exige al usuario descubrir cómo protegerse.
- **Simplicidad y entendibilidad:** una arquitectura que nadie puede explicar tampoco puede auditarse ni recuperarse con confianza.
- **Assume compromise:** detectar, contener, rotar, reconstruir y demostrar recuperación.
- **Evidence over configuration:** declarar un control no prueba que se aplique a la ruta real.

Eliminar complejidad accidental es una medida de fiabilidad y seguridad: dead code, flags permanentes, APIs “misc”, rutas operativas raras y múltiples formas de hacer lo mismo amplían estados que nadie prueba. Mantener APIs mínimas, módulos con una responsabilidad y releases pequeños/atribuibles; version control permite recuperar código eliminado sin dejarlo latente en producción.

### 2.2 Fail-open, fail-closed y degradación

**[IMPL]** Elegir por operación, no por servicio completo:

| Operación | Fallo de control | Default defendible |
|---|---|---|
| lectura pública cacheable | authorizer auxiliar caído | servir sólo contenido inequívocamente público |
| dato privado/escritura | identidad/policy indeterminada | denegar y preservar evidencia mínima |
| pago/orden/acción irreversible | dedupe o ledger no disponible | detener, reconciliar y no adivinar |
| telemetría no crítica | collector caído | buffer acotado/drop conocido; nunca bloquear indefinidamente |
| control de emergencia | sistema normal inaccesible | breakglass temporal, fuerte, auditado y revisado |

Degradación segura especifica qué capacidad permanece, por cuánto tiempo, con qué límite y cómo vuelve al estado normal.

### 2.3 Superficie administrativa mínima

**[PROD]** Google BSaRS recomienda APIs funcionales pequeñas en vez de acceso interactivo amplio a producción. Un endpoint administrativo expresa una acción acotada, tipada, autorizable, idempotente cuando proceda y auditable con `actor/action/target/reason/result/artifact`. Shell/root/console no son la API normal: quedan como breakglass raro. Cada uso de breakglass dispara alarma, expiración, revisión contextual y una pregunta de diseño: qué primitiva segura faltó.

Para acciones masivas/destructivas: selector vacío significa cero targets, nunca “todos”; límite absoluto y por failure domain; preview/dry-run; confirmación ligada al target set/hash; rate de ejecución; cancel/kill switch; y reconciliación idempotente. Probar repetición, target que cambia entre preview/apply, partial failure y credencial con scope excesivo.

## 3. Gobernanza, inventario y clasificación

**[SPEC]** NIST CSF 2.0 organiza resultados en `Govern`, `Identify`, `Protect`, `Detect`, `Respond` y `Recover`; no prescribe una implementación única.

### 3.1 Inventario vivo

Mantener, con owner y última observación:

- servicios, repositorios, dominios, APIs y endpoints;
- cuentas/proyectos/subscriptions, regiones, clusters y namespaces;
- identidades humanas/workload, roles, keys, certificates y trust roots;
- datasets, stores, backups, replicas y rutas de egress;
- dependencias directas/transitivas, images, builders y registries;
- vendors, contratos, subprocesadores, quotas y fechas EOL;
- SLIs/SLOs, runbooks, dashboards, alerts y on-call.

Un inventario manual sin reconciliación con DNS, cloud APIs, repositorios, registry, IAM y runtime envejece silenciosamente.

### 3.2 Datos y lifecycle

Por campo/dataset declarar:

```text
owner → propósito → base/autorización → clasificación
→ lugares y copias → readers/writers → cifrado/keys
→ retención → backup → exportación → borrado verificable
```

- recopilar lo mínimo y no reutilizar fuera del propósito sin decisión;
- separar datos de prod de test y evitar PII/secrets en fixtures;
- registrar residency y transferencias sólo cuando sean requisitos reales;
- cifrado no reemplaza autorización, minimización ni eliminación;
- backup conserva también datos que se pretendían borrar: definir propagación y excepciones.

### 3.3 Privacidad, metadata y anonimato

**[ACADEMIC]** Privacidad incluye control sobre datos y metadata de actividad. Separar cuatro objetivos: confidencialidad de contenido, anonimato respecto de un conjunto/adversario, unlinkability entre identidad y acciones, y unobservability del uso. Un sistema puede satisfacer uno y filtrar los demás.

TLS protege el contenido y autentica el peer bajo su contrato, pero deja observables —según capa y observador— endpoints, DNS, IP, volumen, timing y patrones. VPN/Tor/proxy cambian quién observa y en quién se confía; no eliminan metadata, compromiso de extremos, traffic correlation, fingerprinting ni el tramo de salida sin cifrado end-to-end. Registrar adversario, vantage points, anonymity set, retention y leakage residual antes de afirmar “anónimo” o “privado”.

En producto: minimizar eventos e identificadores, separar propósito, acotar retención, evitar terceros innecesarios, tratar scripts/SDKs como recipients y probar export/delete. Cookies particionadas o bloqueadas no eliminan fingerprinting; no prometer privacidad desde un solo control del navegador.

## 4. Identidad, autenticación y autorización

### 4.1 Identidades distintas

**[SPEC]** NIST SP 800-63-4 separa identity proofing, autenticación y federación con niveles de assurance; la revisión 4 incorpora passkeys sincronizables, fraude y métricas continuas.

No mezclar:

- persona, cuenta, sesión y dispositivo;
- workload, instancia, service account y deployment;
- autenticación (“quién demuestra ser”) con autorización (“qué puede hacer ahora”);
- identidad del issuer con identidad/audience del receptor.

### 4.2 Credenciales humanas

- preferir autenticadores resistentes a phishing para roles privilegiados;
- password manager, rate limits y detección de credential stuffing si persisten passwords;
- recovery de cuenta no puede ser más débil que login;
- sesiones con expiración, revocación, rotación y reautenticación proporcional al riesgo;
- proteger enrollment, cambio de factor, cambio de email y soporte como operaciones privilegiadas;
- breakglass: cuentas mínimas, offline cuando proceda, alertado, probado y con revisión posterior.

### 4.3 Workload identity

**[SPEC]** NIST 800-207/207A rechaza confianza implícita por ubicación de red. SPIFFE modela identidad verificable de workload y SPIRE la emite tras attestation de nodo y proceso.

- preferir identidad corta y renovable sobre secretos estáticos desplegados;
- audience, issuer, subject, trust domain y uso deben validarse;
- un service account compartido destruye atribución y amplía blast radius;
- mTLS autentica peers, pero todavía hace falta autorizar acción/recurso;
- token bearer robado puede reproducirse; usar PoP/mTLS cuando el riesgo lo exige;
- la robustez de workload identity depende del aislamiento: un nodo comprometido puede cambiar el problema.

### 4.4 Policy model

**[IMPL]** La decisión mínima es:

```text
allow = policy(subject, action, resource, tenant, environment, context, time)
```

Probar permitidos y denegados, pertenencia de objeto, campos, funciones administrativas, cambio de roles y revocación. RBAC simple sirve para roles estables; ABAC/ReBAC añade expresividad y también complejidad, ciclos, stale relations y necesidad de explicación.

## 5. Criptografía aplicada y gestión de secretos

### 5.1 No diseñar primitivas propias

**[ACADEMIC]** Boneh/Shoup distingue confidencialidad, integridad, autenticación, semantic security, chosen-ciphertext security, authenticated encryption, signatures y key exchange. Cada construcción posee un objetivo y supuestos; “cifrado” no implica integridad ni autenticidad.

Usar bibliotecas/protocolos revisados y configuración vigente. No inventar cipher, padding, nonce scheme, KDF, firma, password hashing ni serialización criptográfica.

### 5.2 Contrato criptográfico

```yaml
purpose: "encrypt | authenticate | sign | derive | establish"
data_and_metadata: ""
algorithm_and_library: ""
key_owner_and_boundary: "KMS/HSM/process"
nonce_iv_generation: ""
associated_data: "tenant/type/version/context"
rotation_and_overlap: ""
revocation_and_compromise: ""
crypto_agility_and_migration: ""
```

- usar AEAD para datos cuando se necesitan confidencialidad e integridad;
- nonce reuse puede ser catastrófico según esquema: generación forma parte del diseño;
- asociar tenant, versión, tipo o contexto como AAD cuando evita ciphertext swapping;
- AEAD no impide replay: el protocolo añade sequence/nonce uniqueness, freshness, dedupe o state según operación;
- derivar claves separadas por propósito y contexto mediante KDF; no reutilizar una key entre encryption, MAC, tenants, protocolos o ambientes;
- autenticar antes de actuar sobre plaintext y devolver fallos indistinguibles; comparación de tags/secrets evita short-circuit observable;
- signatures autentican artefacto e identidad bajo una trust policy; no prueban que el contenido sea benigno;
- firmar bytes/canonical form inequívocos más domain separation; verificar algoritmo, key identity, context, freshness y anti-replay antes del efecto;
- password hashing usa KDF específica con salt y parámetros medidos; nunca cifrado reversible por comodidad.

Deterministic/searchable/format-preserving/disk encryption tienen objetivos y leakage distintos: igualdad, frequency, access pattern, sector relation o tamaño pueden permanecer visibles. No reemplazarlos por “AES” genérico ni usarlos para habilitar búsqueda sin un leakage/threat contract explícito.

### 5.3 TLS

**[SPEC]** RFC 9325 es BCP para despliegues TLS/DTLS y RFC 9852 exige TLS 1.3 para protocolos nuevos. Para sistemas existentes, verificar BCP/errata y compatibilidad real; no congelar una lista de ciphers en este manual.

- validar hostname/identity y cadena completa; encryption sin peer authentication admite MITM;
- rotar certificados antes de expirar y probar reloj, renovación y revocación;
- 0-RTT puede reproducirse: sólo en operaciones/protocolos cuyo contrato lo permita;
- TLS termina en un punto concreto; mapear plaintext posterior y confianza entre proxies;
- compression, logs y tamaños pueden filtrar metadatos aunque el canal esté cifrado.

### 5.4 Keys y secrets

**[SPEC]** NIST SP 800-57 gobierna lifecycle de material criptográfico. **[IMPL]** Por secreto/key registrar creación, storage, acceso, uso, rotación, overlap, revocación, backup y destrucción.

- secrets fuera de repo, image, command line, crash dump, logs, plan y state;
- KMS/HSM reduce exposición de key material, no corrige IAM excesivo;
- envelope encryption separa data keys y KEK, y permite rotar por capas;
- caches de secrets tienen TTL, zeroization best-effort y conducta ante provider caído;
- rotación debe aceptar old+new durante ventana explícita y retirar old con evidencia;
- post-quantum: inventariar cryptographic dependencies y crypto-agility; adoptar FIPS 203/204/205 sólo mediante protocolos/bibliotecas interoperables aprobados, no ensamblajes caseros.

**[SPEC]** A 2026-08-21 NIST pide iniciar migración a ML-KEM/ML-DSA/SLH-DSA y su timeline IR 8547 apunta a deprecar/eliminar algoritmos vulnerables a quantum de estándares NIST hacia 2035, antes para alto riesgo. Esto no autoriza cambiar un protocolo unilateralmente: primero inventario, data lifetime/harvest-now-decrypt-later, interoperabilidad, performance, tamaños, rollback y validación del módulo.

## 6. Threat modeling y diseño seguro

### 6.1 Flujo mínimo

1. Dibujar procesos, stores, actores, flujos, providers y trust boundaries.
2. Enumerar activos y operaciones de alto impacto.
3. Escribir adversarios/capacidades y abuse cases concretos.
4. Mapear prevención, detección, respuesta y recovery.
5. Identificar control owner, señal, test y riesgo residual.
6. Revisar ante nueva frontera, privilegio, dato, dependency o modo de despliegue.

STRIDE ayuda a preguntar; attack trees ayudan a descomponer objetivos; ATT&CK conecta comportamiento observado. Ninguno sustituye conocimiento del dominio ni prueba exhaustividad.

### 6.2 Tabla de decisión

| Amenaza | Precondición | Control preventivo | Señal | Respuesta | Prueba |
|---|---|---|---|---|---|
| cross-tenant read | ID controlado por cliente | authz por objeto/tenant | denied/allowed audit | revoke + scope | generative access matrix |
| SSRF a metadata | URL/redirect/DNS controlado | destination policy + egress | blocked destination | isolate credential | redirect/rebind/IPv6 corpus |
| artifact replacement | registry/build comprometido | digest+provenance policy | verification failure | halt rollout | tampered artifact |
| secret exfiltration | log/debug/build access | short-lived identity+redaction | canary/egress/audit | rotate+rebuild | seeded secret drill |
| regional loss | shared failure domain | tested alternate path | SLI + provider signal | failover | game day |

## 7. Seguridad de aplicación y API

**[SPEC]** OWASP Top 10:2025 es awareness; ASVS 5.0.0 ofrece requisitos verificables versionados. El manual backend sigue siendo autoridad sobre HTTP/API; aquí se añade threat model y assurance.

### 7.1 Controles por frontera

- parsear con límites de bytes, profundidad, elementos, tiempo y expansión;
- validar sintaxis y semántica, pero autorizar después sobre el recurso resuelto;
- usar queries/APIs parametrizadas; escaping depende del sink/contexto;
- output encoding es contextual: HTML, attribute, URL, JS y CSS no son intercambiables;
- cookies de sesión: `Secure`, `HttpOnly`, `SameSite` y scope mínimo según flujo;
- CORS es política del navegador, no control de acceso del servidor;
- CSRF importa cuando credenciales se adjuntan automáticamente;
- `GET` y demás métodos seguros no producen side effects; CSP reduce fuentes/ejecución permitidas y SRI liga un recurso tercero a bytes esperados, pero ninguno corrige autorización o lógica defectuosa;
- cada script, iframe, SDK y tag de analítica de tercero cruza una frontera de datos/capacidades; reducirlos, aislarlos y fijar integridad/origin según el caso;
- SSRF exige resolver redirects, DNS rebinding, schemes, IPv4/IPv6, link-local/metadata y egress;
- upload: tamaño, tipo real, nombre, descompresión, parser isolation, storage y serving domain;
- errores no filtran stack/secrets, pero preservan correlation y causa interna segura.

### 7.2 Consumo hostil y costo

Autenticación no vuelve benigno al cliente. Presupuestar por tenant/actor:

```text
request bytes + parse CPU + DB rows/locks + fanout
+ downstream calls + response bytes + logs/traces + dinero externo
```

Rate limit sin quotas por operación/costo puede proteger requests baratos y dejar abierto el flujo caro. Usar admission, concurrency limits, byte/work budgets, timeouts, pagination, circuit breaking y límites de cardinalidad.

### 7.3 Assurance

- requisitos ASVS citados como `v5.0.0-x.y.z` cuando se adopten;
- SAST/SCA/DAST/fuzzing detectan clases diferentes y producen falsos positivos/negativos;
- pentest valida un alcance y momento; no reemplaza SDLC;
- bug bounty amplía observación, pero necesita safe harbor, scope y respuesta;
- findings tienen reachability, exploitability, asset exposure, impacto, owner y SLA; CVSS solo no decide.

### 7.4 Gate ASVS 5.0.0 por aplicabilidad

**[SPEC]** El artefacto JSON oficial versionado contiene 345 requisitos en 17 capítulos. No se marca “ASVS compliant” por ejecutar un scanner: el proyecto fija alcance/nivel, registra `aplica/no aplica + razón`, cita IDs `v5.0.0-x.y.z` y conserva evidencia reproducible.

| Capítulo (requisitos) | Pregunta/gate de implementación |
|---|---|
| V1 Encoding/Sanitization (30) | ¿decode canónico una vez y encode/sanitize según el intérprete final; queries/commands sin concatenación? |
| V2 Validation/Business Logic (13) | ¿reglas combinadas, secuencia, límites, atomicidad y locking evitan saltos/doble consumo? |
| V3 Web Frontend (31) | ¿browser boundary, headers, origins, CSP, cookies, DOM y terceros están controlados? |
| V4 API/Web Service (16) | ¿método/media/schema/size/URL y consumo upstream se validan y acotan? |
| V5 File Handling (13) | ¿nombre/path/type/content/archive/storage/serving se tratan como hostiles? |
| V6 Authentication (47) | ¿enrollment, factores, recovery, throttling, federation y eventos sensibles prueban assurance elegido? |
| V7 Session (19) | ¿tokens son impredecibles, rotan, expiran, se invalidan y resisten fixation/replay/abuse? |
| V8 Authorization (13) | ¿trusted service autoriza cada operación/objeto/tenant con cambios efectivos y tests negativos? |
| V9 Self-contained Tokens (7) | ¿algoritmo/key/claims/audience/tiempo/replay y revocación compensan su estado autocontenido? |
| V10 OAuth/OIDC (36) | ¿roles, flow, redirect, PKCE/state/nonce, issuer, clients y tokens siguen el threat model del protocolo? |
| V11 Cryptography (24) | ¿inventario, librería validada, agility, randomness, AEAD/MAC/hash/key lifecycle y failure mode son correctos? |
| V12 Secure Communication (12) | ¿TLS/certificados/configuración y canales internos/externos autentican el peer esperado? |
| V13 Configuration (21) | ¿defaults, debug, headers, dependencies, admin, backup y deploy config están endurecidos y auditados? |
| V14 Data Protection (13) | ¿clasificación, minimización, retention/deletion, cache/client/log y privacidad siguen el lifecycle? |
| V15 Secure Coding/Architecture (21) | ¿threat model, components, trust boundaries, unsafe primitives y memory/concurrency errors tienen controles? |
| V16 Logging/Error (17) | ¿security events están inventariados, estructurados, protegidos/separados y errores fallan sin filtrar datos? |
| V17 WebRTC (12) | si aplica: ¿TURN, DTLS-SRTP, fingerprint, signaling/media malformed/flood y recording resisten abuso? |

La suite §23 es el mínimo transversal; el proyecto añade casos para cada requisito aplicable. Los capítulos con 0 requisitos aplicables también dejan razón y reviewer, porque omitir por accidente no equivale a `N/A`.

## 8. Host, runtime y red

### 8.1 Host/runtime

- reducir packages, daemons, capabilities y writable paths;
- ejecutar no-root y separar identidad del proceso de la del host;
- mantener kernel/runtime, reiniciar cuando el patch lo requiere y verificar versión efectiva;
- ASLR, NX, CFI, sanitizers y memory-safe languages reducen clases distintas; ninguna elimina lógica insegura;
- sandbox/seccomp/AppArmor/SELinux limitan syscalls/objetos según policy probada;
- debug endpoint, profiler, socket runtime y metadata service son superficies privilegiadas;
- proteger time sync: certificados, expiraciones, logs y protocolos dependen del reloj.

Un reference monitor útil intercepta **toda** operación protegida inmediatamente antes de ejecutarla, decide con policy inequívoca y resiste manipulación. Si existe una ruta alternativa —debug socket, filesystem, node proxy, cache, sidecar o TOCTOU entre check/use— no hay complete mediation. En sandboxing, reducir syscalls/capabilities, aplicar `no_new_privs`, seccomp/MAC y resource limits, y elegir una frontera más fuerte —proceso, container, VM o hardware— según adversario; recursos compartidos y covert/side channels impiden prometer aislamiento absoluto.

### 8.2 Red

- segmentar por flujo necesario, no por diagrama ornamental;
- ingress y egress tienen policy, DNS behavior, proxy chain y ownership;
- firewall permite paquetes; la aplicación todavía autentica y autoriza;
- DDoS se trata en edge, red, aplicación, dependencia y costo económico;
- registrar source original sólo mediante proxies confiables; headers reenviados son input hostil;
- network encryption no oculta volumen, timing, endpoints ni todos los metadatos.

### 8.3 Zero trust sin eslogan

**[SPEC]** SP 800-207 protege recursos y evalúa subject/device/workload, no concede por “red interna”. Esto no significa consultar un control remoto en cada instrucción ni comprar un service mesh.

La implementación debe definir identity issuance, policy decision/enforcement, freshness, revocation, failure mode, logging y bootstrap. Si el policy plane cae, la conducta por operación ya está decidida en §2.2.

### 8.4 Hardware, attestation y microarquitectura

**[ACADEMIC]** El estado arquitectónico puede revertir una ejecución especulativa mientras caches, TLB, branch predictors, consumo u otros efectos permanecen observables. Por eso bounds checks, separación de procesos o una enclave no bastan contra Spectre y familias relacionadas cuando el adversario comparte recursos y puede medirlos.

- attestation demuestra una medición/identidad bajo una chain y policy concretas; no demuestra ausencia de bugs, input seguro, configuración correcta ni que el resultado sea correcto;
- verificar measurement, signer/TCB version, freshness, revocation, claims y binding de la sesión/clave antes de entregar secretos;
- aplicar microcode, kernel, hypervisor y compiler mitigations vigentes; aislar cores/tenants o evitar co-residencia sólo cuando el threat model lo justifique;
- una barrera o mitigación localizada debe cubrir todos los gadgets relevantes y medirse: las defensas incompletas fallan y las globales pueden degradar severamente el rendimiento;
- confidential computing exige pruebas propias de attestation, patching, side-channel assumptions, observabilidad, backup/restore y respuesta ante revocación.

## 9. Secure SDLC y supply chain

### 9.1 SSDF operativo

**[SPEC]** NIST SSDF 1.1 agrupa preparar organización, proteger software, producir software bien asegurado y responder a vulnerabilidades. Convertirlo en pipeline:

```text
threat model + requirements
→ protected source + review
→ hermetic/isolated build where justified
→ tests + scanners + fuzzing
→ immutable artifact + SBOM + provenance + signature
→ policy verification at promotion/deploy
→ runtime inventory + vulnerability response
→ rebuild, revoke and learn
```

#### Memory safety, fuzzing y sanitizers

Preferir lenguajes memory-safe para componentes nuevos expuestos a input hostil; runtimes y `unsafe` conservan TCB y bugs de lógica. Cuando C/C++ sea necesario, combinar safe subsets/APIs, hardening del compilador, ASLR/NX/CFI y aislamiento. Estas capas reducen explotabilidad; no prueban ausencia de vulnerabilidades.

Flujo continuo de fuzzing: construir un harness pequeño sobre el parser/estado real → seed corpus representativo → mutation/generation y coverage guidance → sanitizer build → minimizar/deduplicar/reproducir → corregir → regression test → conservar corpus. Ejecutar también el corpus que alcanzó plateau bajo ASan/UBSan/MSan/TSan según clase; sanitizers aumentan cobertura diagnóstica, pero no son una defensa productiva ni una prueba de completitud. Sandbox, timeout, memoria/disco y crash artifacts quedan acotados porque el target procesa datos hostiles.

### 9.2 Source y CI/CD

- MFA resistente a phishing y mínimo admin sobre org/repo;
- protected branches/tags, review y status checks; controlar bypass;
- pin actions/plugins/dependencies por digest/revisión cuando sea posible;
- forks/PR no confiables no reciben secrets ni runner persistente privilegiado;
- separar build de signing/deploy; credenciales efímeras y audience restringida;
- runner limpio, red/egress acotada, logs protegidos y artefactos inmutables;
- cambio de pipeline, policy o dependency bot tiene el mismo rigor que código productivo.

### 9.3 SBOM, provenance y firma

**[SPEC]** SLSA 1.2 separa source y build tracks. Build L1 crea provenance; L2 la firma desde plataforma hosted; L3 endurece aislamiento y protege signing material frente al build controlado por usuario.

No confundir:

- **SBOM:** qué componentes declara el artefacto;
- **provenance:** quién/cómo/con qué inputs lo produjo;
- **signature:** integridad + identidad bajo trust policy;
- **vulnerability scan:** coincidencias conocidas en un momento;
- **reproducible build:** dos builds coinciden bajo condiciones declaradas.

Separar también:

- **hermetic:** todos los inputs —source, dependencies, toolchain y ambiente relevante— están resueltos antes del build;
- **reproducible:** mismos inputs/comandos producen output bit a bit idéntico;
- **verifiable:** el consumidor puede confiar en el vínculo entre artefacto, inputs y builder bajo su threat model.

El admission gate de provenance realiza tres verificaciones independientes: autenticidad/integridad de la attestation, digest que la liga exactamente al artefacto y cumplimiento de policy para ese ambiente. Firma válida no demuestra por sí sola source permitido, review, tests ni ausencia de vulnerabilidades. Preferir provenance inequívoca propagada con el artefacto y una policy aplicable inequívoca; builders controlados por usuarios no acceden al signing material.

**[SPEC]** Sigstore keyless liga una ephemeral key a identidad OIDC mediante certificado corto y registra evidencia en transparency logs. Verificación debe fijar issuer, subject/identity, artifact digest, inclusion/bundle y policy; “signed” sin identidad esperada no basta. Transparency exige monitoreo para detectar emisión indebida.

### 9.4 Dependency policy

Registrar versión/digest, maintainer/source, licencia, transitive tree, release cadence, EOL, privileges y alternativa. Priorizar remediation por exposición, reachability, impacto, exploit maturity y CISA KEV, no sólo por score. Un update puede ser tanto corrección como nuevo riesgo: canary, rollback y test.

## 10. Arquitectura cloud

### 10.1 Cloud no transfiere responsabilidad del workload

Los frameworks Well-Architected de AWS, Google Cloud y Azure convergen en seguridad, fiabilidad, operación, performance y costo; AWS/GCP añaden sostenibilidad explícita. Son lentes de revisión, no certificaciones automáticas.

Por cada servicio administrado declarar:

| Capa | Proveedor | Nuestro equipo | Evidencia |
|---|---|---|---|
| hardware/control plane | alcance exacto del SLA y patching | configuración, dependencias y fallback | docs + status test |
| data plane | disponibilidad/cifrado ofrecidos | schema, authz, backup, restore, quotas | drill + metrics |
| IAM/keys | primitives y logs | policy, identities, rotation, review | access tests |
| red | fabric/edge básico | topology, ingress/egress, DNS, TLS | flow tests |
| incidentes | provider response | detección propia, escalation y continuidad | runbook/game day |

### 10.2 Landing zone y blast radius

- organización/root protegido y fuera de uso diario;
- cuentas/proyectos/subscriptions por ambiente, riesgo y ownership;
- identidad federada central, roles temporales y separación billing/security/platform/workload;
- regiones/zones elegidas por SLO, data residency, latency, service availability y costo;
- centralizar audit logs con write protection y acceso separado;
- policies/guardrails impiden clases de configuración, pero tienen rollout y breakglass;
- tags/labels obligatorios: owner, environment, data class, cost center, expiry;
- quotas se inventarían, monitorean y prueban antes de necesitarse.

### 10.3 Selección de compute

| Opción | Favorece | Costos/riesgos que permanecen |
|---|---|---|
| serverless/FaaS | escala/eventos, menor host operation | cold start, limits, concurrency, lock-in y costos por invocación |
| managed container | packaging y control moderado | image/runtime/IAM/network/observability propios |
| Kubernetes | plataforma multi-workload y policies comunes | gran control plane humano/técnico; no justificarlo por moda |
| VM | control de OS/runtime | patch, image, process supervision y capacity |
| bare metal | hardware/latency/control | provisioning, failure, security y ciclo completo |

Elegir la abstracción más alta que satisfaga semántica, SLO, compliance, hardware y portabilidad. Kubernetes no es requisito de un backend sólido.

### 10.4 Datos, regiones y DR

- multi-zone no equivale a multi-region; replica no equivale a backup;
- cada servicio tiene consistency, durability, backup y restore semantics propios;
- cross-region aumenta latency, egress, conflicto, key/data residency y complejidad;
- active/active exige resolver ownership, routing, split brain, replication lag y reconciliación;
- active/passive exige capacidad reservada/provisionable, datos recuperables y promoción probada;
- provider SLA no es SLO end-to-end y sus exclusiones deben modelarse.

### 10.5 FinOps como corrección operativa

Costo es una salida del workload:

```text
cost/request or cost/event
= compute + memory + storage-capacity + storage-ops
 + network-egress + managed requests + observability + redundancy
```

Presupuestos y anomaly alerts no sustituyen hard limits/admission. Medir unit economics, idle, commitment utilization, egress y costo de recovery. Spot/preemptible requiere checkpoint/idempotencia; no es capacidad garantizada.

## 11. Infraestructura como código

### 11.1 Contrato

```text
versioned desired configuration
→ validate/lint/security policy
→ speculative plan
→ reviewed immutable plan
→ apply with scoped identity
→ verify runtime and drift
→ preserve state/audit/rollback path
```

- pin CLI, providers, modules y lockfiles;
- módulos pequeños con inputs tipados, defaults seguros e invariantes;
- policy-as-code valida propiedades; revisar también excepciones/bypass;
- cambios destructivos, reemplazos, IAM y egress requieren atención explícita;
- `target`/manual console changes son excepción documentada y reconciliada;
- rollback de infraestructura puede no revertir datos: expand/contract y restore aparte.

### 11.2 State y plan son sensibles

**[SPEC]** Terraform/OpenTofu state mapea configuración a recursos y puede contener secretos; marcar `sensitive` puede ocultar UI sin eliminar valor del state. Plan guardado puede contener configuración y valores completos.

- backend remoto con access control, encryption, locking/versioning y audit;
- state/plan nunca en VCS ni artifact público;
- separar states por blast radius/ownership, no un monolito privilegiado;
- backup y restore del state probados; proteger contra replay/rollback malicioso;
- minimizar `remote_state` amplio: publicar outputs mínimos o consultar API;
- identidad de plan puede leer; identidad de apply modifica; ambas mínimas y temporales.

### 11.3 Drift y reconciliación

Detectar drift periódicamente, clasificar autorizado/incidente y reconciliar hacia código o importar con review. Auto-remediation sólo para invariantes seguros y libres de bucles; medir acciones, fallos y blast radius.

## 12. Contenedores y Kubernetes

### 12.1 Versiones y API lifecycle

**[SPEC]** Kubernetes 1.36 es la serie activa verificada al 2026-08-21; el patch target debe consultarse al desplegar. Se soportan las tres ramas minor recientes y rige version-skew/upgrade order oficial. Alpha/beta/GA expresan estabilidad, no idoneidad automática.

### 12.2 Image y runtime

- minimal base por digest, usuario no-root, read-only rootfs cuando sea posible;
- eliminar shell/package manager si no cumplen función; conservar debug mediante mecanismo separado y autorizado;
- no hornear secrets; evitar environment si pueden filtrarse en dumps/diagnóstico;
- requests/limits, ephemeral storage y PID limits según medición;
- drop capabilities, `allowPrivilegeEscalation: false`, seccomp y MAC policy;
- firmar/provenance y verificar antes de admission; scan no reemplaza verificación;
- registry retention, immutability, replication y recovery definidos.

### 12.3 Control plane e IAM

Kubernetes security oficial enfatiza TLS/control de API, encryption at rest, Pod Security Standards, NetworkPolicy, admission y audit.

- acceso humano federado; reservar `system:masters` sólo para bootstrap/break-glass, con custodia y auditoría;
- RBAC por verbs/resources/namespaces; evitar wildcard, secret read y impersonate;
- tratar `create`/`patch` sobre Pods y recursos que los crean como capacidad de ejecutar código: admission debe impedir host access, identidad o privilegios superiores al caller;
- auditar como escalación: `bind`, `escalate`, `impersonate`, CSR approval, `serviceaccounts/token`, `nodes/proxy`, PV/hostPath, webhook configuration y mutación de labels de Namespace que gobiernan PSS/NetworkPolicy;
- `get nodes/proxy` no es read-only: alcanza APIs privilegiadas de kubelet y puede omitir admission/audit del API server;
- service account por workload, projected short-lived tokens y audience;
- desactivar automount cuando no se necesita API;
- bloquear metadata de cloud desde workloads salvo necesidad explícita y usar identidad de workload en vez de credenciales de nodo;
- proteger kubeconfig: un `exec` plugin puede ejecutar código local;
- proteger root/intermediate CA, expiración y rotación; etcd usa mTLS y trust root separado cuando la arquitectura lo permite;
- API server, kubelet y etcd no son públicos; autenticar/autorizar kubelet explícitamente y revisar bindings de `system:unauthenticated`;
- admission webhook/policy tiene availability, timeout, failure policy, version skew y emergency bypass explícitos;
- evaluar `NodeRestriction` y `DenyServiceExternalIPs`; cambios de admission son privileged code con review, canary y rollback;
- etcd contiene estado crítico y secrets; cifrado, network isolation, backup y restore.

### 12.4 Aislamiento y multi-tenancy

Namespace no es una frontera completa. Combinar según riesgo:

- RBAC + quotas/LimitRange;
- Pod Security `restricted` o baseline justificado;
- default-deny ingress/egress con DNS/control-plane paths explícitos;
- node pools/runtime classes/sandbox para cargas de distinta confianza;
- storage class, PV, snapshots y backup aislados;
- admission y image policy;
- clusters/cuentas separados cuando compartir kernel/control plane no satisface el threat model.

NetworkPolicy depende del plugin y no cubre por sí sola todos los paths (hostNetwork, node, control plane, DNS, metadata). Probar tráfico permitido y prohibido.

Restringir `LoadBalancer` y `externalIPs`; un Service o route puede ampliar exposición fuera del namespace. Verificar que el CNI realmente implemente ingress/egress policy y que el deny-all seleccione todos los Pods, incluidos los creados después.

### 12.5 Scheduling y disponibilidad

- requests son input del scheduler; limits actúan en runtime y pueden causar throttling/OOM;
- el overcommit de memoria (`request < limit`) expone al nodo a presión/OOM; igualdad request-limit obtiene QoS más predecible pero no es regla universal: elegir por workload y verificar bajo presión;
- CPU limit puede producir throttling aun con CPU disponible; adoptarlo, omitirlo o separarlo del request exige medición de latencia, fairness y blast radius;
- probes tienen semánticas distintas: startup protege arranque, readiness enruta, liveness reinicia;
- una liveness mal diseñada convierte lentitud de dependencia en restart storm;
- topology spread/anti-affinity se alinean con zones/failure domains reales;
- PDB limita disrupciones voluntarias, no garantiza capacidad ante fallo;
- autoscaling usa señales retrasadas y tiene máximos, stabilization y dependency quotas;
- graceful termination coordina endpoint removal, drain, deadline y trabajo en curso.

### 12.6 Operación del cluster

- control plane/etcd en failure domains suficientes para el SLO;
- no autoscalar etcd como servicio stateless; quorum y latency de storage gobiernan;
- upgrade: patch actual → control plane → componentes → nodes, respetando skew y APIs removidas;
- validar CRDs/webhooks/controllers frente a nuevos fields/versions;
- backup de manifests no reemplaza snapshot/backup de stores de aplicación;
- restore etcd detiene API servers y se ensaya en entorno aislado;
- medir API saturation, workqueues, admission latency, scheduler, node pressure, DNS, CNI, CSI y image pulls.

### 12.7 Secrets y gate de cluster

Secret usa base64, no cifrado, y por defecto puede quedar sin cifrar en etcd. Activar encryption at rest/KMS según threat model; `get`, `list` y `watch` revelan contenido. Crear workloads también puede exfiltrar Secrets, ConfigMaps, PVs y ServiceAccounts montables del namespace. Restringir cada mount al container que lo necesita, preferir credenciales cortas, evitar manifests en VCS y auditar lecturas masivas.

Antes de producción, demostrar en runtime: versión/skew soportados; CA/cert rotation; API/kubelet/etcd no públicos; RBAC/escalation matrix; PSS enforce; seccomp/MAC; NetworkPolicy positiva/negativa; metadata/egress; quotas/object-count; Secret encryption/mount; image digest/signature; admission failure mode; audit protection; node isolation; upgrade, drain y etcd/application restore. El checklist upstream es baseline mutable, no certificación.

### 12.8 Código Kubernetes — reconciliar desde verdad observable

**[PROD]** Código upstream fijado en `e81f39c0e03ce8ed8e2660c9147b391edd9e262b` (Apache-2.0). La workqueue mantiene conjuntos separados de keys `dirty` y `processing`: eventos repetidos se comprimen, pero una key marcada mientras se procesa vuelve a encolarse al llamar `Done`. El Deployment controller espera sincronización inicial de informers, procesa una key sin concurrencia consigo misma, hace `DeepCopy` antes de mutar objetos cacheados, admite tombstones de deletes perdidos y aplica retry rate-limited con límite. Para una adopción peligrosa revalida deletion mediante lectura fresca: el cache es una vista útil, no autoridad universal.

Patrón de controller/operador transferible:

```text
evento/resync → enqueue key → leer desired + observed state
→ calcular delta idempotente → aplicar con precondition/version
→ publicar status/condition → Forget o retry con backoff
```

- La notificación es un hint para reconciliar, no una orden que deba ejecutarse una vez; puede duplicarse, comprimirse, llegar tarde o faltar hasta el relist.
- El sync converge desde estado actual y soporta restart entre side effect y status update. No guardar corrección únicamente en memoria del worker.
- Conflicto/NotFound y objeto que cambió durante el trabajo son transiciones normales; requeue/leer otra vez, no sobrescribir ciegamente.
- Un finalizer es un protocolo de borrado durable: marcar progreso, ejecutar cleanup reintentable, comprobar postcondición y sólo entonces retirarlo. Debe tener timeout/escalation y runbook para finalizers huérfanos; quitarlo manualmente acepta el recurso residual.
- Gate: tests con add/update/delete duplicado, tombstone stale, evento perdido + resync, conflicto, crash antes/después del side effect, API throttling, retry exhaustion, shutdown/drain y dos workers intentando recursos relacionados. Observar queue depth, work duration, retries, oldest item y reconcile errors por razón.

## 13. SRE: contrato de fiabilidad

### 13.1 SLI, SLO y SLA

**[PROD]** Google SRE usa SLO y error budget para negociar riesgo de forma observable.

- **SLI:** medición de una experiencia o propiedad relevante.
- **SLO:** target del SLI en una ventana.
- **SLA:** compromiso contractual y consecuencias; puede diferir del SLO interno.

Para eventos totales `N` y buenos `G`:

```text
SLI = G / N
error_budget_fraction = 1 - SLO
bad_events_allowed = N × (1 - SLO)
```

Disponibilidad temporal sirve si el usuario percibe tiempo; event-based suele describir mejor requests. Definir qué cuenta, qué se excluye, source of truth, window, late data y reset.

### 13.2 SLIs por workload

Separar **SLI specification** —resultado que importa al usuario, independiente del sensor— de **SLI implementation** —fuente y cálculo concretos—. Logs de servidor, load balancer, black-box y cliente observan fronteras distintas; documentar coverage, bias, costo y failure mode del medidor. Empezar con no más de unas pocas dimensiones críticas y refinar mediante evidencia evita convertir cada métrica en un objetivo.

| Workload | Good event/estado |
|---|---|
| API | respuesta semánticamente correcta dentro del umbral |
| stream | evento íntegro procesado antes de deadline, sin gap indebido |
| batch | resultado correcto y completo antes del horario |
| storage | operación con durability/consistency declaradas dentro de latency |
| ML inference | output válido dentro de SLO y quality/safety gate |
| control plane | cambio converge correctamente dentro del tiempo |

No contar `200` si el body es error o stale fuera de contrato. No excluir failures internos para mejorar el ratio. Separar availability, latency, freshness, correctness y durability cuando un solo número los oculta.

### 13.3 Error budget policy

Definir ventana, burn thresholds y acciones:

```text
burn_rate = observed_bad_fraction / (1 - SLO)
budget_consumed_over_window = burn_rate × observation_window / SLO_window
```

`burn_rate = 1` consumiría exactamente el budget a lo largo de toda la ventana; `10` lo consume diez veces más rápido. Alertar con pares de ventanas —corta+larga— para detectar burn rápido sin paginar por un pico aislado, y otro par para burn lento persistente. Los factores, ventanas y umbrales se calculan desde el SLO y el tiempo humano de respuesta del servicio; no se copian como constantes universales.

En servicios de poco tráfico, un solo evento puede aparentar burn extremo sin señal estadística suficiente. Decidir por impacto: synthetic/black-box traffic que no oculte fallos reales, agregación sólo de operaciones con failure mode común, fallback/retry seguro o SLO/window renegociados. Si una única operación es irreversible y valiosa, investigar cada fallo puede ser correcto, pero la alerta llega después del daño y exige prevención adicional.

- budget sano: cambios normales con gates;
- burn elevado: canary más lento, congelar cambios riesgosos, corregir top causes;
- agotado: detener feature risk salvo emergencia, invertir en fiabilidad;
- incidente externo/exclusión: sólo según política previa y transparente.

Error budget no concede permiso para violar seguridad, integridad o obligaciones duras.

### 13.4 Toil y automatización

Toil es manual, repetitivo, automatizable, reactivo y crece con servicio. Automatizar después de entender invariant, failure modes, authorization, rate y rollback. La automatización multiplica tanto corrección como error; usar dry-run, límites, canary, approvals por riesgo y kill switch.

### 13.5 On-call sostenible

Una rotación profesional define cobertura, carga máxima aceptable de páginas, secondary/escalation, handoff, permisos y compensación. Cada página exige acción urgente, dashboard y runbook; si sólo informa, es ticket o métrica. Evitar dependencia de una persona mediante shadowing, ejercicios y acceso break-glass probado. Revisar periódicamente páginas por turno, wake-ups, tiempo de diagnóstico/mitigación, falsos positivos, salud del equipo y aprendizaje; una guardia crónicamente saturada es un defecto del sistema.

**[PROD]** Como referencia explícitamente contextual, Google publica un objetivo máximo de dos incidentes por turno y al menos 50% de tiempo de proyecto para corregir causas estructurales. No son constantes universales: sirven como alarma para dimensionar la rotación y proteger salud, seguimiento y trabajo de ingeniería.

La overload del equipo también se mide: tiempo en pages/tickets/toil, carga posterior al turno, backlog/edad, context switches, horas, descansos, capacidad de proyecto y percepción de seguridad psicológica. Triage conjunto puede eliminar, devolver, diferir, documentar/self-service o automatizar trabajo; ocultar alertas sólo es temporal, con expiry y ticket. Cuando urgencias impiden corregir sus fuentes, el servicio ya está en fallo operativo.

## 14. Capacidad, overload y resilience

### 14.1 Capacity model

```text
offered_load
→ admission
→ queues/pools
→ CPU/memory/network/storage/GPU/dependencies
→ goodput within SLO
```

Medir headroom por recurso y failure domain, no sólo promedio global. Inventariar quotas administrativas y físicas, provisioning lead time y capacidad de failover. Una región secundaria sin datos/capacidad/cuota no es DR.

### 14.2 Overload

- bounded queues y admission antes de trabajo caro;
- fairness por tenant/prioridad; reservar capacidad crítica;
- load shedding temprano, explícito y observable;
- deadlines propagados y cancellation real;
- retries con budget, idempotencia, backoff+jitter y máximo;
- circuit breaking limita dependencia fallida, pero puede sincronizar y oscilar;
- caches/stale fallback sólo bajo contrato de freshness/authz;
- degradation ladder probada, no improvisada durante incidente.

**[ACADEMIC]** Diseñar contra asimetría: antes de autenticar/validar al cliente, el servidor no crea estado ni realiza trabajo mucho mayor que el del solicitante. Medir `server CPU/bytes/state/money ÷ attacker CPU/bytes`; reducir amplification, validar reachability antes de reservar estado cuando el protocolo lo permita y mover admission delante de parsing, crypto, fanout o storage caros. DDoS puede agotar red, conexiones, threads, memoria, dependencia, logs o presupuesto aunque el promedio de CPU parezca sano.

Un retry budget es global al request tree, no “N por cada capa”: con cuatro intentos en tres niveles una acción puede producir `4³ = 64` intentos al backend. Un solo nivel coordina retries; clasifica errores permanentes/retriables/overload, respeta deadline y corta trabajo cancelado. Separar **process health** —scheduler: el proceso vive— de **service readiness/health** —load balancer: puede servir esta clase— para que overload no retire capacidad en masa y realimente la cascada.

Balanceo, autoscaling, admission y shedding forman un solo sistema de control. Instrumentar sus intersecciones y probar feedback loops: si el shed mantiene CPU plana, un balanceador que interpreta menor CPU/request como eficiencia puede enviar aún más carga al dominio saturado. Escalar antes del umbral de shed cuando sea viable, mantener capacidad mínima por failure domain y propagar deadlines/cancelación para no trabajar en requests ya inútiles.

### 14.3 Failure domains y blast radius

Particionar por tenant, shard, cell, zone, region o riesgo. Cada compartment tiene límites de recursos, rollout y dependencia; evitar global locks/configs/credentials que vuelvan ficticio el aislamiento.

### 14.4 Resilience validation

Game days y fault injection prueban hipótesis concretas: node loss, dependency timeout, DNS failure, disk full, quota, expired cert, revoked credential, corrupt artifact, region loss. Abort conditions y observabilidad se validan antes; chaos sin hipótesis sólo crea ruido.

## 15. Observabilidad y alertas

### 15.1 Señales

**[SPEC]** OpenTelemetry define APIs/SDK/protocol por señal con estabilidad distinta; verificar implementación concreta. Context propagation correlaciona, pero baggage puede exponer/propagar datos y no es un almacén arbitrario.

- **metrics:** agregación y alertas; controlar units, temporality y cardinalidad;
- **logs:** eventos discretos y audit; estructurados, redacted, con retention/access;
- **traces:** causalidad distribuida muestreada; no sustituyen métricas completas;
- **profiles:** costo por stack; overhead y sensibilidad explícitos;
- **audit:** quién hizo qué sobre qué, resultado y policy; integridad/acceso reforzados.

### 15.2 Instrumentación mínima

**[PROD]** Prometheus recomienda para online serving count, errors, latency e in-progress; para pipelines, items/edad por etapa; para batch, último éxito y duración.

Registrar además queues, retries, shed, dependency calls, resource saturation, deploy/version y state transitions. No usar user ID, email, token, URL arbitraria o error text como label: cada combinación crea una serie.

### 15.3 Distribuciones

Promedios ocultan tails. Histograms agregables permiten percentiles con error según resolución/buckets; summaries precomputadas normalmente no se agregan entre instancias. Elegir buckets alrededor de SLOs y validar costo/cardinalidad.

### 15.4 Alertas accionables

Una página necesita:

- impacto o riesgo inmediato para usuario/invariante;
- urgencia que exige acción humana;
- señal suficientemente específica y sostenible;
- runbook, dashboard, owner y escalation;
- dedupe/inhibition y test del canal.

Alertar síntomas/SLO burn antes que cada causa posible. Causas alimentan diagnóstico o tickets. Detectar también silencio: exporter/collector/heartbeat caído.

### 15.5 Telemetría como sistema hostil/costoso

- redaction antes de exportar y allowlist de attributes;
- sampling por decisión consciente, preservando errores/slow traces según política;
- buffers bounded y conducta ante backend caído;
- acceso y retención por sensibilidad; audit de consultas;
- medir CPU, allocation, network, storage, cardinalidad y drops de la propia telemetría;
- logs no deben convertirse en canal de injection ni ejecutar contenido al visualizarse.

### 15.6 Interfaz de investigación

Read-only/debug access también expone datos y capacidad. Preferir consultas y endpoints pequeños con redaction, límites, identidad temporal y audit sobre shell o dumps completos. El cuaderno de investigación separa `observado`, `esperado`, `hipótesis`, `prueba`, `resultado` y `descartado`; conserva resultados negativos para evitar repetir trabajo. Conocer el baseline no autoriza normalizar desviaciones crónicas. Reproducir fuera de producción e aislar el mínimo caso cuando sea posible; si no, usar canary/fault injection con blast radius y stop condition.

## 16. Entrega y cambios seguros

### 16.1 Artifact promotion

Construir una vez y promover el mismo digest; configuración/secret/environment se enlazan aparte con versión. El gate verifica tests, provenance, firma, policy, schema compatibility, threat model delta, rollback y observabilidad.

Tests de release se ejecutan sobre el artefacto y combinación de configuración que realmente se promueven, incluidos cherry-picks/mixed versions; registrar revision, toolchain, dependencies, test results y change set. Configuración cambia conducta con la misma capacidad de causar outage que un binary: versionarla, revisarla, probar su parser/schema/compatibilidad y comparar **desired vs deployed state**. Breakglass acelera la aplicación, no apaga tests: estos continúan con prioridad y retroanotan el cambio urgente.

### 16.2 Rollout

| Estrategia | Ventaja | Riesgo/gate |
|---|---|---|
| rolling | simple/capacidad eficiente | mixed versions y rollback real |
| canary | limita blast radius | muestra representativa + automated analysis |
| blue/green | cutover/rollback rápido | doble capacidad, state/data compatibility |
| feature flag | separa deploy/release | flags son config privilegiada, deuda y combinaciones |

Ordenar cambios de schema/protocol con expand→migrate/backfill→switch→contract. Rollback de binary no revierte mensajes, side effects ni datos escritos bajo schema nuevo.

### 16.3 Gate de canary

Comparar baseline/candidate con tráfico, tenant mix y tiempo suficientes:

- SLI/error budget burn, p99 y goodput;
- errors por clase, saturation, queues y retries;
- seguridad: authz denials/anomalías, secrets, policy failures;
- costo/unit y capacity;
- correctness/domain/ML quality;
- ausencia de evidencia es fallo del gate, no aprobación.

La muestra debe ser representativa y el cambio atribuible: misma clase de hosts/zones, carga y routing; evitar que varios cambios o un selection bias expliquen la diferencia. Artefactos pequeños y autocontenidos abaratan rollback y diagnóstico. Definir antes tamaño, duración mínima/máxima, error-budget permitido, hipótesis, thresholds, abort/rollback y conducta si el análisis falla; aumentar exposición por etapas sólo con evidencia fresca.

### 16.4 Launch readiness

Un launch no es un deploy ordinario con más tráfico. Antes de abrir audiencia/región/partner, un reviewer independiente verifica arquitectura y request flow, dependencies/owners, estimación y mezcla de carga, cuotas/headroom por failure domain, SLO/alerts/on-call, privacidad/abuse, data migration, terceros, rollback/kill switches y comunicación. Modelar launch spike y fanout —una page view puede originar muchas requests—; empezar con audiencia/región/rate acotados, observar el sistema completo y ampliar por etapas. La checklist es reproducible pero proporcional: excepciones tienen riesgo, owner y expiry; un check marcado sin evidencia no es gate.

## 17. Incident response

**[SPEC]** NIST SP 800-61r3 integra incident response en CSF 2.0: preparación y reducción de impacto atraviesan `Govern/Identify/Protect`, mientras `Detect/Respond/Recover` manejan el evento y aprendizaje.

### 17.1 Roles y estado

- incident commander coordina y decide prioridades;
- operations/subject leads investigan y mitigan;
- communications mantiene updates internos/externos;
- scribe registra timeline, hipótesis, acciones y evidencia;
- security/legal/privacy/business se incorporan por impacto, sin bloquear contención urgente preautorizada.

Declarar temprano. El incident commander coordina, no monopoliza diagnóstico; operations mitiga y communications mantiene una interpretación común. Handoff transmite estado, impacto, evidencia, decisiones, riesgos, tareas y “qué haría el turno saliente a continuación”. **[PROD]** BSaRS recomienda no superar 12 horas continuas en crisis; es un techo contextual, no una meta de explotación humana.

```text
DETECTED → TRIAGED → CONTAINING → MITIGATING
→ RECOVERING → MONITORING → CLOSED → LEARNING
```

### 17.2 Protocolo

1. Confirmar señal, scope, usuarios/invariantes y severidad.
2. Abrir canal/documento, roles, reloj común y cadencia.
3. Registrar por separado hechos, expectativas, hipótesis, pruebas y resultados descartados.
4. Preservar evidencia proporcional sin impedir seguridad humana/sistema.
5. Contener y mitigar impacto antes de esperar una root cause perfecta.
6. Erradicar causa y caminos de persistencia; rotar/rebuild desde trusted source.
7. Recuperar por etapas con verification y observación reforzada.
8. Notificar según contrato/ley y mantener hechos separados de hipótesis.
9. Cerrar sólo con criterios definidos y owner de acciones.

### 17.3 Seguridad durante incidente

- no pegar secrets/PII en chat/ticket;
- commands destructivos requieren target verificado y segundo control según riesgo;
- breakglass y forensics access expiran y quedan auditados;
- snapshot forense no se vuelve backup confiable;
- no confiar en logs/control plane potencialmente comprometidos sin corroboración;
- rotar credenciales siguiendo dependencias para evitar outage secundario.

Ante adversario activo, asumir que puede observar y reaccionar. Investigación, coordinación, build, identity o restore tooling potencialmente comprometidos requieren un canal/entorno limpio y aislado. Coordinar ejection con investigación: esperar más evidencia puede reducir reinfección, pero daño activo a personas, datos o servicio puede exigir contención inmediata. La decisión registra incertidumbre, costo de demora y posibles respuestas del atacante.

### 17.4 Postmortem que aprende

Timeline factual; impacto; detección; respuesta; contributing conditions; qué funcionó/falló; acciones con owner, deadline y verificación. Evitar “error humano” como causa final: preguntar qué diseño permitió que una acción normal tuviera ese blast radius.

Las acciones fuertes cambian código, control, test, alerta, arquitectura o simulacro. “Tener más cuidado” no cierra recurrencia.

Clasificar cada acción como prevención, detección, mitigación, reparación o investigación; asignar prioridad, owner único, tracking y criterio verificable. “Mejorar X” no es una acción. Revisar cierre e impacto transversal: el valor del postmortem está en el cambio ejecutado y compartido, no en publicar el documento.

### 17.5 Ejercicios que cambian el sistema

Tabletop no es una conversación pasiva: escenario creíble, artefactos/señales realistas, decision points, roles que ejecutan escalation/comunicación y acciones con owner. Luego fault injection o game day valida los caminos que necesitan producción, con scope y abort conditions. Probar también dependencias de respuesta —identity, chat, telefonía, status, logs y acceso— porque el canal normal puede fallar durante el mismo incidente.

## 18. Backup, restore y disaster recovery

### 18.1 Objetivos

- **RPO:** pérdida temporal máxima de datos aceptable.
- **RTO:** tiempo máximo para restaurar servicio/capacidad declarados.
- **MTPD/impact tolerance:** límite de interrupción del negocio cuando aplique.

Un backup existe sólo si es íntegro, accesible durante el desastre, compatible y restaurable dentro de objetivos.

### 18.2 Diseño

- inventario de datos/config/keys/dependencies necesario para restore;
- copias con failure/security domain distinto y protección contra borrado/ransomware;
- cifrado con keys recuperables sin depender del sistema perdido;
- retention/versioning e inmutabilidad según amenaza;
- consistency/cut elegido para bases, objects, queues y metadata;
- restore isolated, malware/corruption check, schema migration y functional verification;
- reconciliar side effects/eventos ocurridos después del punto restaurado.

Para integridad de datos aplicar capas distintas: undelete/soft-delete contra error de usuario/desarrollador → backups definidos desde escenarios de **recovery** → validators fuera de banda sobre invariantes devastadoras. Soft-delete tiene ventana y purge verificable compatibles con privacidad; no es permiso para retener indefinidamente. Validación a escala puede ser tiered —checks críticos frecuentes y auditoría más rigurosa menos frecuente— con rate limits, logs drill-down, alertas y playbook. Réplicas, backups y validators que comparten bug, credencial o failure domain no son defensas independientes.

Un backup puede conservar persistencia, configuración maliciosa o corrupción anterior a la detección. Estimar timeline de compromiso, poner copias dudosas en cuarentena, restaurar desde un punto justificadamente limpio y buscar variantes en assets que comparten image, librería, credencial o control plane. Cuando no se puede confiar en limpieza in-place, reconstruir desde source/image/config conocidos y rotar credenciales en un orden que no vuelva a entregar acceso al entorno viejo.

Codificar **estado pretendido** y observar **estado desplegado** para código, configuración, identidad/policy, firmware cuando aplique, datos y estado en memoria. Detectar drift y reparar/reconstruir mediante el camino rutinario; la ruta “sólo para emergencias” se pudre si no se ejercita. Incluir images/configs fallback o inactivas, cuotas, capacidad y dependencias de arranque: encender componentes sin recuperar sus invariantes no es recuperación.

### 18.3 Drill

Medir desde declaración hasta servicio verificado: localizar backup, conseguir acceso, restaurar infra/data/keys, validar invariantes, cambiar routing y observar. Registrar RPO real, RTO real, pasos manuales, dependencias ocultas y capacidad.

## 19. Vulnerability y exposure management

Pipeline:

```text
inventory → discover → normalize/dedupe → enrich
→ exposure/reachability/KEV/impact → decide
→ patch/mitigate/accept/remove → verify → learn
```

- cubrir internet assets, cloud configs, hosts, images, packages, repos, SaaS e identidades;
- mantener version/runtime real, no sólo lockfile;
- CISA KEV prueba explotación observada y eleva prioridad; no contiene todas las amenazas;
- exceptions tienen scope, compensating controls, owner y expiry;
- patch incluye build, test, deploy y verification; descargar no remedia;
- medir time-to-triage/remediate por risk y reincidencia, no premiar cerrar falsos positivos;
- EOL es una deuda con fecha y migration plan.

## 20. Seguridad de sistemas de IA

El manual de IA define calidad/modelo; éste define superficies operativas:

- prompt, retrieved content, tool output y model artifact son inputs no confiables;
- separar instrucciones de datos por arquitectura y capability enforcement, no sólo delimitadores;
- tool calls pasan schema validation, authz por usuario/recurso, budgets y confirmación para side effects críticos;
- secrets nunca en prompt/system prompt si el modelo o logs pueden exponerlos;
- RAG preserva ACL/tenant/freshness y protege ingestion contra poisoning;
- model/checkpoint/adapter/tokenizer son artefactos ejecutables o data parsable: provenance, trust y sandbox;
- training data, evals y traces tienen consentimiento, licencia, privacidad y retention;
- rate/cost/loop limits, cancellation y human escalation para agentes;
- evaluar prompt injection, data exfiltration, excessive agency, unsafe output, denial-of-wallet y supply chain.

Un LLM no es policy engine, secret store, transaction coordinator ni detector de intrusión confiable por defecto.

**[ACADEMIC]** Stanford CS155 muestra que prompt injection es una confusión entre datos e instrucciones y puede llegar por texto, RAG, email, web, base64 o imagen. Escapar/delimitar el input o preguntar a un segundo LLM si hay una inyección no establece una frontera fiable. El patrón dual —modelo privilegiado que no consume contenido hostil y modelo en cuarentena sin acceso a datos/capacidades adicionales— reduce alcance, pero tampoco resuelve todo si al modelo aislado se le entrega información excesiva.

Contrato defendible para un agente:

1. un planner propone un flujo tipado; no concede permisos;
2. un intérprete determinista y pequeño valida capabilities, identity/tenant, destinatario, schema, secuencia y presupuesto;
3. el componente que lee contenido no confiable recibe sólo datos mínimos y handles opacos, nunca secrets o tools ambientales;
4. lectura, escritura, comunicación externa y acción irreversible son capacidades distintas, short-lived y deny-by-default;
5. cada side effect se revalida justo antes de commit; confirmación humana liga intención, diff/targets y expiración;
6. evals adversariales cubren contenido activo/pasivo, multimodal/encoded, tool output, RAG poisoning, exfiltración y combinaciones de varios turnos.

La seguridad del modelo también requiere threat model de training/inference. Pesos ocultos no impiden por sí solos adversarial examples transferibles ni model extraction mediante queries. Prompts/respuestas pueden revelar datos memorizados y APIs con logits/probabilidades amplían extracción: minimizar outputs, acceso, rate y retention, y medir leakage. Para verificar training, attestation prueba que cierto código/dato medido ejecutó bajo un TCB; rerun/auditoría necesita controlar o registrar fuentes de no determinismo —incluida aritmética/orden GPU— y aun así no prueba que los datos, objetivos o código sean benignos.

## 21. Gates profesionales

### 21.1 Antes de implementar

- [ ] owners, users, data, trust boundaries y critical journeys declarados;
- [ ] threat/failure model y riesgos residuales aceptados por owner;
- [ ] SLI/SLO, RTO/RPO, quotas y budget definidos;
- [ ] identity/authz y secret/key lifecycle diseñados;
- [ ] dependency/provider responsibility y exit/recovery comprendidos;
- [ ] rollout, rollback, telemetry y incident path diseñados.

### 21.2 Antes de producción

- [ ] tests de autorización negativa y tenant isolation;
- [ ] input/resource/egress limits y overload probados;
- [ ] artifact digest, SBOM/provenance/signature policy verificados;
- [ ] IaC plan revisado, state protegido y drift observable;
- [ ] least privilege y breakglass probados;
- [ ] dashboards/alerts/runbooks/on-call activos;
- [ ] load, failure, rollback y restore drills dentro de objetivo;
- [ ] logs/traces sin secretos/PII no autorizados;
- [ ] canary compara correctness, SLO, seguridad y costo;
- [ ] known risks/exceptions con owner y expiry.

### 21.3 Operación continua

- [ ] inventory y ownership reconciliados;
- [ ] access/key/cert review y rotación;
- [ ] vulnerability/KEV/EOL y dependency updates;
- [ ] error budget/capacity/quotas/cost;
- [ ] backup restore y disaster exercise;
- [ ] incident/postmortem actions verificadas;
- [ ] threat model y ASVS/control baseline revalidados tras cambios.

## 22. Contrato para un agente Codex

Entrada mínima:

```yaml
task: ""
system_and_environment: ""
assets_and_data_classes: []
users_tenants_and_trust_boundaries: []
threat_and_failure_model: []
identity_and_authorization: ""
slo_rto_rpo_and_capacity: ""
cloud_runtime_versions: ""
deployment_and_rollback: ""
evidence_available: []
constraints_and_open_questions: []
```

Codex debe devolver:

1. supuestos y `[OPEN]` que bloquean garantías;
2. invariantes y abuse/failure cases;
3. diseño mínimo y trust boundaries;
4. código/config/IaC con versiones y defaults seguros;
5. pruebas positivas, negativas, fuzz, concurrency/failure y restore según riesgo;
6. identity, secrets, supply chain y data lifecycle;
7. SLIs, telemetry, alerts y cardinality/cost budgets;
8. rollout, rollback, incident y recovery;
9. evidencia que acepta/rechaza el cambio y riesgo residual.

Prohibido para el agente:

- inventar compliance o garantías del proveedor;
- afirmar “zero trust”, “encrypted”, “HA”, “exactly-once” o “secure” sin frontera y prueba;
- registrar/exponer credenciales para depurar;
- desactivar authz, TLS, validation, audit, durability o policy para mejorar benchmark;
- usar wildcard IAM, cluster-admin, privileged container o public network sin justificación explícita;
- ejecutar cambio destructivo sin target, backup/rollback y autorización;
- recomendar crypto propia o algoritmo por memoria sin documentación vigente;
- cerrar un incidente porque el síntoma desapareció sin comprobar invariantes.

## 23. Suites de aceptación ejecutables

| Suite | Casos mínimos | Gate |
|---|---|---|
| `security_authz` | matriz subject/action/resource/tenant; revoked/stale/admin | ninguna operación no autorizada |
| `security_input` | límites, encoding, injection, SSRF, upload, parser fuzz | bounded y sin boundary bypass |
| `host_memory_isolation` | sanitizer builds, corpus regression, exploit mitigations, syscall/capability policy, TOCTOU/escape | bug reproducible; frontera y blast radius demostrados |
| `security_identity` | expiry, audience, issuer, rotation, recovery, breakglass | fail mode y audit correctos |
| `security_secrets` | scanning, log/trace/crash/state, rotate/revoke | cero exposición; recovery probado |
| `supply_chain` | tamper source/build/artifact/provenance/identity | promotion/admission rechaza |
| `cloud_iam_network` | public exposure, wildcard, egress, metadata, cross-account | policies efectivas en runtime |
| `k8s_isolation` | RBAC, SA, PSS, NetworkPolicy, quota, node compromise assumptions | tenant/workload boundary demostrada |
| `reliability_load` | ramp, burst, soak, dependency slow/down, retry storm | goodput/SLO y bounded resources |
| `reliability_change` | canary, mixed version, schema, rollback | no invariant ni budget violado |
| `incident` | credential/artifact/tenant/region scenarios | roles, contain, evidence y comms |
| `backup_restore` | deleted/corrupt/ransomware/region loss | RPO/RTO e invariantes cumplidos |
| `observability` | collector down, high cardinality, redaction, missing telemetry | señal útil, bounded y privada |
| `privacy_metadata` | content vs endpoint/timing/volume, identifiers, third parties, export/delete | leakage residual y lifecycle demostrados |
| `ai_agent_boundary` | injection directa/indirecta/multimodal, poisoned RAG/tool output, exfil, side effects, loops/costo | contenido no concede capability ni cruza tenant/recipient |

## 24. Ruta práctica profesional

1. **Servicio local:** API con authz por objeto, secrets externos, ASVS subset, fuzz y audit.
2. **Supply chain:** build aislado, SBOM, SLSA provenance, firma/verificación y tamper tests.
3. **IaC:** red/IAM/compute/data con plan review, remote state, policy, drift y teardown seguro.
4. **Kubernetes:** workload no-root, SA mínima, PSS, NetworkPolicy, quotas, rollout y restore.
5. **SRE:** SLIs/SLO/error budget, load shedding, dashboard, burn alerts y runbooks.
6. **Incidente:** tabletop + game day de credencial comprometida, dependency outage y región perdida.
7. **Capstone:** servicio multi-tenant con canary, supply-chain gate, DR medido y postmortem.

Cada práctica entrega threat model, decisiones, repo/config, tests, evidencias, benchmark/costo, runbook y reflexión sobre fallos encontrados. No se copian soluciones de cursos.

## 25. Auditoría transversal — reglas de frontera

| Manual relacionado | Conserva autoridad sobre | Obligación de este manual |
|---|---|---|
| `COMPUTER_SYSTEMS_PERFORMANCE_LOW_LATENCY.md` | OS, CPU, memory model, NUMA, profiling y tails | hardening/telemetry no rompe correctness ni SLO; acceso debug es mínimo y auditado |
| `SOFTWARE_BACKEND_API_ENGINEERING.md` | contrato HTTP/API, lifecycle, idempotencia y tests | authn/z, abuso, supply chain, secrets y operación se integran sin cambiar semántica |
| `NETWORKING_DISTRIBUTED_STREAMING.md` | protocolo, partial failure, consensus, delivery y backpressure | TLS/identity/policy no inventan delivery; partitions y DDoS entran al failure model |
| `DATABASE_STORAGE_INTERNALS.md` | isolation, durability, WAL, backup/restore y query semantics | cifrado/IAM/cloud conservan recovery; réplica no se llama backup |
| `GPU_ACCELERATED_COMPUTING.md` | device/stream/allocator/collective correctness | drivers/artifacts/tenants aislados; GPU failure y OOM tienen containment/rollback |
| `MARKET_MICROSTRUCTURE_EXCHANGE_SYSTEMS.md` | book/matching/order lifecycle, venue rules y risk semantics | feed read-only por defecto; trading keys segregadas; pre-trade/kill/drop-copy/audit no son saltables por strategy/agent |
| manuales de IA | modelo, datos, eval, error analysis y calidad | proteger artifacts/data/tools y operar dentro de authz, SLO, costo y incident policy |
| `FRONTEND_PRODUCT_ENGINEERING_UX.md` | browser lifecycle, semantic UI, accessibility, state visible y UX | threat model cubre DOM XSS, CSRF/CORS, storage/cache/SW, third parties, privacy y supply chain sin delegar authz al cliente |

Regla: seguridad puede rechazar una optimización insegura; no puede redefinir silenciosamente el contrato de datos/red/modelo. Fiabilidad puede exigir redundancia; no puede copiar datos o ampliar privilegios sin threat/privacy review.

Cruce recíproco completado el 2026-08-21: los manuales de cada frontera material contienen ahora sus obligaciones de seguridad/SRE/cloud, por lo que la regla existe en ambos lados y no sólo en este documento. Para frontend, CSP/Trusted Types son defensa en profundidad; server-side authz, output encoding/sanitization, session/CSRF, dependency provenance y privacy lifecycle conservan autoridad. Para mercados, secrets/roles separan market-data, order-entry, drop-copy y admin; para aplicaciones nativas, sandbox/permissions/secure storage/IPC/signing/update forman el límite local. Incident response preserva raw evidence, bloquea exposición nueva y valida estado antes de declarar cierre.

Cruce con `ALGORITHMS_DATA_STRUCTURES_PROBLEM_SOLVING.md` cerrado el 2026-08-21: complejidad y estructura forman parte del threat model. Hash flooding, regex/backtracking, nesting, decompression bombs, graph/path explosion, output no acotado y peores casos de sorting/parsing requieren límites de bytes, nodos, profundidad, tiempo y memoria antes de ejecutar. Aleatorización, aproximación o Bloom filters declaran probabilidad/error y nunca sustituyen autorización, fuente de verdad, SLO ni fallback seguro.

Cruce con `SOFTWARE_ARCHITECTURE_SYSTEM_DESIGN.md` cerrado el 2026-08-21: threat/trust boundaries, SLO, failure domains, RTO/RPO, identity y operación son drivers iniciales, no anexos. Este manual valida controles, plataforma y respuesta; arquitectura conserva escenarios, vistas y ADR. Multi-region, managed service, zero trust o Kubernetes no prueban disponibilidad/seguridad: capacity, policy efectiva, restore, failover, rollback y incident paths deben ejecutarse.

Cruce con `DATA_ENGINEERING_ANALYTICS.md` cerrado el 2026-08-21: clasificación, purpose, IAM por dataset/column/tenant, encryption/key lifecycle, retention/delete, egress y audit viajan por raw, staging, mart, cache, notebook y backup. Lineage localiza posibles copias, no prueba borrado ni autorización. Data SLO/cost/quality se operan con bounded cardinality, owners, incident/backfill/restore y supply-chain provenance del job/DAG.

Cruce con `TOOLCHAINS_BUILDS_PACKAGING_FFI.md` cerrado el 2026-08-21: CI/build scripts/package lifecycle son código privilegiado; usar identity efímera, inputs por digest, aislamiento/network policy y secret-free outputs. SBOM inventaría, provenance enlaza recipe/inputs, signature identifica signer: ninguno sustituye los otros ni vulnerability/runtime policy. Promover mismo digest, verificar en admission y conservar revocation/rollback/rebuild reproducible.

Cruce con `NATIVE_MOBILE_DESKTOP_ENGINEERING.md` cerrado el 2026-08-21: sandbox, permissions/entitlements, secure storage, deep links/files, WebView/IPC, publisher identity y updater forman trust boundaries. El cliente no guarda secretos universales ni toma decisiones autoritativas; backend reautoriza. SRE segmenta crash/hang/startup/update por build/OS/device con minimización de datos, staged rollout y kill switch; stores lentos y clients antiguos obligan compatibility/forward recovery.

## 26. Matriz de trazabilidad inicial

| Capacidad | Autoridad | Salida verificable | Estado |
|---|---|---|---|
| threat/system security | Stanford CS155 + Google BSaRS | threat model, invariants, abuse/failure tests | 17 decks públicos/1.038 páginas auditados; incorporado |
| criptografía | Boneh/Shoup + NIST/IETF | crypto contract, key lifecycle, protocol tests | 7 semanas públicas inventariadas; contrato auditado |
| governance/risk | NIST CSF 2.0 | current/target profile y owners | incorporado |
| identity/zero trust | NIST 800-63-4, 800-207/A + SPIFFE | assurance, workload identity, policy tests | incorporado |
| AppSec | OWASP Top 10:2025, ASVS 5.0.0, API 2023 | versioned requirements + security suites | 345 requisitos/17 capítulos inventariados; gate incorporado |
| secure SDLC | SSDF 1.1 + SLSA 1.2 + Sigstore | provenance/signature/promotion policy | incorporado |
| cloud | AWS/GCP/Azure Well-Architected | responsibility, landing zone, cost/DR review | incorporado; provider deltas condicionados |
| Kubernetes | upstream 1.36 | manifests/policies/upgrade/restore tests | security checklist, RBAC, Secrets y PSS auditados |
| SRE | Google SRE/Workbook/BSaRS | SLO, budget, load, incident/postmortem | auditoría selectiva de alta densidad cerrada; incorporado |
| observability | OTel + Prometheus | schema, cardinality, alerts, redaction tests | incorporado |
| incident response | NIST 800-61r3 + Google | playbook, drill, recovery y actions | incorporado |

### 26.1 Inventario Stanford CS155 Spring 2026

El syllabus público enumera 19 clases. Se inspeccionaron los 17 decks PDF publicados —1.038 páginas— y se transfirió su semántica operativa; las clases invitadas 13 y 19 no publican slides. Recordings privadas en Canvas, evaluaciones y materiales no enlazados no se consideran inspeccionados.

| Clases | Núcleo público | Incorporación operativa |
|---|---|---|
| 1 | overview, threat model, TCB y principios | §§1–2 y §6 |
| 2–7 | control hijacking, defensa en profundidad, fuzzing/sanitizers, OS, reference monitor, sandboxing, enclaves/Spectre | §§6, 8–9, 21 y 23 |
| 8–10 | web model, ataques y defensas | §7 y suite `web_api_abuse` |
| 11–12 | criptografía aplicada y HTTPS | §5 y suites crypto/transport |
| 13 | software supply chain | sin deck público; contrato cubierto por SSDF/SLSA/Sigstore y lectura indicada en syllabus |
| 14–17 | protocolos, Internet, DoS/asimetría, VPN/zero trust, metadata, anonimato/Tor y censura | §§3, 5, 8, 14 y 15 |
| 18 | poisoning/adversarial examples, prompt injection/agentes, privacidad/extracción y training verificable | §20 |
| 19 | invitado | sin material público verificable; no se atribuye contenido |

### 26.2 Inventario OWASP ASVS 5.0.0

ASVS contiene 345 requisitos en 17 capítulos según el JSON inglés del tag `v5.0.0`. El manual conserva el mapa y deriva controles por riesgo; no replica el catálogo normativo. Cuando un proyecto adopta un requisito, debe fijar versión e identificador exacto.

| Capítulos | Tema | Salida en este manual |
|---|---|---|
| V1–V2 | encoding/sanitization; validation/business logic | parser limits, validación semántica y abuse tests |
| V3–V5 | frontend web; API/web service; file handling | browser/API boundary, CORS/CSRF, uploads y parsers hostiles |
| V6–V10 | authentication, session, authorization, tokens, OAuth/OIDC | identity contract, authz por objeto y suites de sesión/token |
| V11–V12 | cryptography; secure communication | crypto/transport contract y key/cert lifecycle |
| V13–V14 | configuration; data protection | hardened defaults, secrets, minimización y lifecycle |
| V15–V17 | secure architecture/coding; logging/error; WebRTC | invariants, audit/redaction/error handling; WebRTC condicionado al producto |

### 26.3 Inventario Google de fiabilidad y seguridad

| Corpus público | Cobertura verificable | Uso aquí | Estado |
|---|---:|---|---|
| *Site Reliability Engineering* | 34 capítulos | SLO, toil, monitoring, automation, release, incidents, overload, consensus, data integrity y launch | auditoría selectiva: caps. 8–9, 17–18, 22 y 26–27 |
| *The Site Reliability Workbook* | 21 capítulos | implementación SLO, burn alerts, on-call, postmortems, load, canary y engagement | auditoría selectiva: caps. 2, 5, 8–11 y 16–17 |
| *Building Secure and Reliable Systems* | 21 capítulos | adversarios, least privilege, invariants, resilience, recovery, secure SDLC, deployment, investigations y crisis | auditoría selectiva: caps. 5, 8–9 y 14–18 |

La auditoría selectiva privilegia capítulos con reglas transferibles y los cruza con los manuales especializados; no afirma extracción línea por línea de 76 capítulos totales. Se incorporaron: release/config real, simplicidad, tests de producción, capacity intent, retry cascades, recovery/data validation, launch readiness, SLO/burn alerts, on-call/incidents/postmortems, load/canary y diseño/operación bajo adversario.

### 26.4 Inventario del curso público de Dan Boneh

| Semana | Material público | Semántica transferida |
|---|---|---|
| 1 | overview, probability, stream ciphers, PRG y semantic security | objetivo/supuesto antes de elegir construcción; keystream/nonce no se reutiliza |
| 2 | block ciphers, PRP/PRF, CPA y modes | primitive ≠ scheme; IV/nonce/mode son parte de seguridad |
| 3 | MAC, collision resistance, HMAC y timing | integridad separada, domain/key separation y verificación constante |
| 4 | authenticated encryption, CCA, padding/non-atomic attacks, KDF, deterministic/tweakable/FPE | AEAD, authenticate-before-use y leakage contract |
| 5 | trusted parties, Diffie–Hellman y number theory | key establishment necesita autenticación y parámetros validados |
| 6 | public-key encryption, RSA/PKCS#1, ElGamal y attacks | usar schemes/protocols revisados, no textbook primitives |
| 7 | digital y hash-based signatures | signature contract, identity/policy, canonicalization y migration |

El curso enseña fundamentos y ataques históricos para comprender por qué fallan combinaciones; las decisiones de despliegue actuales se fijan con NIST/IETF y bibliotecas/protocolos vigentes. Videos enlazados públicamente se inventariaron; evaluaciones no se usaron.

## 27. Fuentes primarias inventariadas

### 27.1 Academia y textos expertos

- Stanford CS155 Spring 2026: <https://cs155.stanford.edu/>, <https://cs155.stanford.edu/syllabus.html>. El syllabus contiene 19 clases; se inspeccionaron los 17 decks PDF públicos (1.038 páginas). Las clases 13 y 19 no enlazan slides; recordings están en Canvas y no se consideran vistas.
- Dan Boneh, curso público: <https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/>.
- Boneh y Shoup, *A Graduate Course in Applied Cryptography* v0.6: <https://crypto.stanford.edu/~dabo/cryptobook/>.

### 27.2 NIST, IETF, CISA y MITRE

- NIST CSF 2.0: <https://www.nist.gov/cyberframework>.
- NIST SSDF 1.1, SP 800-218: <https://csrc.nist.gov/pubs/sp/800/218/final>.
- Zero Trust SP 800-207 y cloud-native 207A: <https://csrc.nist.gov/pubs/sp/800/207/final>, <https://csrc.nist.gov/pubs/sp/800/207/a/final>.
- Digital Identity SP 800-63-4 suite, final July 2025: <https://pages.nist.gov/800-63-4/>.
- Incident response SP 800-61r3, final April 2025: <https://csrc.nist.gov/pubs/sp/800/61/r3/final>.
- Key management SP 800-57: <https://csrc.nist.gov/projects/key-management/key-management-guidelines>.
- PQC FIPS 203/204/205: <https://csrc.nist.gov/Projects/Post-Quantum-Cryptography>.
- TLS/DTLS BCP RFC 9325 y protocolos nuevos RFC 9852: <https://www.rfc-editor.org/rfc/rfc9325.html>, <https://www.rfc-editor.org/rfc/rfc9852.html>.
- CISA Secure by Design y KEV: <https://www.cisa.gov/securebydesign>, <https://www.cisa.gov/known-exploited-vulnerabilities-catalog>.
- MITRE ATT&CK Enterprise: <https://attack.mitre.org/>.

### 27.3 OWASP

- OWASP Top 10:2025: <https://owasp.org/Top10/>.
- ASVS 5.0.0: <https://owasp.org/www-project-application-security-verification-standard/>.
- API Security Top 10 2023: <https://owasp.org/www-project-api-security/>.
- Cheat Sheet Series: <https://cheatsheetseries.owasp.org/>.
- SAMM: <https://owasp.org/www-project-samm/>.

### 27.4 Seguridad, supply chain e identidad de workload

- *Building Secure and Reliable Systems*: <https://google.github.io/building-secure-and-reliable-systems/>.
- SLSA 1.2: <https://slsa.dev/spec/v1.2/>.
- in-toto stable 1.0: <https://in-toto.io/docs/specs/>.
- Sigstore security/sign/verify: <https://docs.sigstore.dev/about/security/>, <https://docs.sigstore.dev/cosign/signing/overview/>, <https://docs.sigstore.dev/cosign/verifying/verify/>.
- CycloneDX 1.7: <https://cyclonedx.org/specification/overview/>.
- SPIFFE/SPIRE: <https://spiffe.io/docs/latest/spiffe/concepts/>, <https://spiffe.io/docs/latest/spire-about/spire-concepts/>.

### 27.5 SRE, observabilidad y cloud

- Google SRE books: <https://sre.google/books/>.
- SRE book: <https://sre.google/sre-book/table-of-contents/>.
- SRE Workbook: <https://sre.google/workbook/table-of-contents/>.
- OpenTelemetry specification: <https://opentelemetry.io/docs/specs/otel/>, <https://opentelemetry.io/docs/specs/status/>.
- Prometheus practices: <https://prometheus.io/docs/practices/instrumentation/>, <https://prometheus.io/docs/practices/histograms/>, <https://prometheus.io/docs/practices/naming/>.
- AWS Well-Architected: <https://docs.aws.amazon.com/wellarchitected/latest/framework/welcome.html>.
- Google Cloud Well-Architected: <https://docs.cloud.google.com/architecture/framework>.
- Azure Well-Architected: <https://learn.microsoft.com/azure/well-architected/>.

### 27.6 Kubernetes, containers e IaC

- Kubernetes 1.36 release series: <https://kubernetes.io/releases/1.36/>; confirmar patch vigente al aplicar.
- Kubernetes fijado `e81f39c0e03ce8ed8e2660c9147b391edd9e262b`: [workqueue](https://github.com/kubernetes/kubernetes/blob/e81f39c0e03ce8ed8e2660c9147b391edd9e262b/staging/src/k8s.io/client-go/util/workqueue/queue.go), [Deployment controller](https://github.com/kubernetes/kubernetes/blob/e81f39c0e03ce8ed8e2660c9147b391edd9e262b/pkg/controller/deployment/deployment_controller.go) y [CRD finalizer](https://github.com/kubernetes/kubernetes/blob/e81f39c0e03ce8ed8e2660c9147b391edd9e262b/staging/src/k8s.io/apiextensions-apiserver/pkg/controller/finalizer/crd_finalizer.go).
- Security overview/checklist: <https://kubernetes.io/docs/concepts/security/>, <https://kubernetes.io/docs/concepts/security/security-checklist/>.
- Version skew: <https://kubernetes.io/releases/version-skew-policy/>.
- Multiple zones y etcd recovery: <https://kubernetes.io/docs/setup/best-practices/multiple-zones/>, <https://kubernetes.io/docs/tasks/administer-cluster/configure-upgrade-etcd/>.
- Terraform state/plan: <https://developer.hashicorp.com/terraform/language/state>, <https://developer.hashicorp.com/terraform/cli/commands/plan>.
- OpenTofu state security: <https://opentofu.org/docs/language/state/sensitive-data/>.

## 28. Extensiones deliberadamente condicionadas

- Reauditar capítulos Google SRE/Workbook/BSaRS adicionales sólo cuando un proyecto revele una laguna no cubierta; el gate ASVS 5.0.0 ya está inventariado, pero cada proyecto debe seleccionar requisitos aplicables.
- Reauditar Stanford CS155/Boneh sólo ante una edición/material público nuevo; no acceder ni reproducir soluciones evaluadas.
- Añadir controles concretos AWS/GCP/Azure sólo cuando un proyecto declare proveedor, servicios, regiones y versión.
- Los cruces recíprocos de seguridad/SRE/cloud con todas las fronteras materiales actuales están cerrados; reabrirlos sólo ante una frontera o versión material nueva.
- Compliance específico —PCI DSS, HIPAA, SOC 2, ISO 27001, GDPR, normativa argentina u otra— requiere alcance, jurisdicción y evidencia propia; no se infiere desde este manual.
- Hardware roots, confidential computing, formal verification, offensive security avanzada y forensics profundo requieren fuentes especializadas si el proyecto los necesita.

## Regla final

```text
Threat model defines relevance.
Invariants define correctness.
Least privilege limits blast radius.
SLO and error budget govern reliability risk.
Provenance and identity establish what may run.
Telemetry must detect real failure without leaking the system.
Restore and incident drills decide whether recovery exists.
```
