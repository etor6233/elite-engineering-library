# Authority Freshness and Self-Correction Contract

## 1. Propósito

El agente debe poder detectar que su propia información, un Markdown, un lock o una decisión quedaron obsoletos. No corrige por memoria ni sobreescribe historia: compara contra fuentes oficiales actuales, preserva el snapshot anterior, crea una decisión sucesora y repite los gates afectados.

“Autocorrección” significa revisión trazable de conocimiento y código; no significa que el agente pueda inventar hechos, modificar reglas del negocio o autoaprobar autoridad.

## 2. Artefacto obligatorio

Todo proyecto mantiene:

```text
PROJECT_AUTHORITY_FRESHNESS_RECORD.md
```

desde `PROJECT_AUTHORITY_FRESHNESS_TEMPLATE.md` para toda autoridad material: leyes/contratos provistos, API/provider docs, runtimes, arquitecturas, source locks, modelos/datasets, seguridad, licencias y operación.

## 3. Señales de obsolescencia

- fecha/versión/EOL excedida;
- enlace oficial movido, release superseded o branch móvil;
- contradicción entre docs, SDK, schema, runtime y comportamiento observado;
- advisory, errata, deprecación o default nuevo;
- proveedor cambia API version, scope, webhook, terms, región o quota;
- tests/fixtures ya no reproducen el claim;
- el agent descubre que una respuesta previa fue inferida, incompleta o falsa;
- conteos/hashes/planes no coinciden con el snapshot material;
- decisión de negocio/jurisdicción modifica la autoridad aplicable.

## 4. Jerarquía de evidencia

Para hechos externos actuales, priorizar:

1. contrato, regulación o dato autorizado por su owner competente;
2. release, registry, advisory, schema y documentación oficial vigente;
3. código oficial en revisión inmutable y comportamiento reproducido;
4. evidencia local del proyecto;
5. manuales/mapas de Elite como router y síntesis;
6. memoria del agente sólo como hipótesis a verificar.

Una fuente más nueva no siempre reemplaza una versión contractual fijada. Registrar applicability, effective date y target environment.

## 5. Bucle de autocorrección

```text
detectar contradicción/staleness
→ registrar FAILURE o FRESHNESS event
→ congelar el claim afectado
→ localizar fuente oficial primaria
→ fijar versión/fecha/hash y applicability
→ comparar old/new y blast radius
→ clasificar factual|contract|code|dependency|business decision
→ crear superseding record; no borrar snapshot
→ actualizar mapas/contracts/locks/packs canónicos
→ regenerar hashes/SBOM/notices si aplica
→ ejecutar tests/gates dependientes
→ actualizar readiness y comunicar corrección explícita
```

Si la fuente no puede verificarse, el claim pasa a `STALE_BLOCKED` o `CONFLICT_UNRESOLVED`; no se elige la versión más conveniente.

## 6. Corrección de una respuesta previa del agente

El agente debe decir claramente:

- qué afirmó;
- qué evidencia nueva demuestra el problema;
- cuál es la afirmación corregida y su alcance;
- qué archivos/decisiones/código podrían haberse afectado;
- qué gates se repitieron y qué sigue bloqueado.

La evidencia estructurada conserva al menos `previous_statement`, `contradictory_evidence`, `corrected_statement`, `blast_radius`, `files_or_decisions_changed`, `gates_rerun` y `related_failure_ids`.

No ocultar la corrección en una edición silenciosa. Registrar el evento en `PROJECT_FAILURE_LESSONS.md` cuando el error pudo afectar diseño, código, seguridad, datos, costo u operación.

## 7. Estados

- `CURRENT_VERIFIED`: fuente oficial y applicability verificadas dentro de la ventana;
- `CURRENT_CONDITIONED`: vigente, pero requiere target/access/test;
- `REVIEW_DUE`: ventana vencida; no implica falsedad;
- `STALE_BLOCKED`: fuente/versión ya no gobierna y falta reemplazo probado;
- `CONFLICT_UNRESOLVED`: autoridades aplicables contradicen;
- `SUPERSEDED`: reemplazado por un record nuevo, conservado históricamente;
- `NOT_APPLICABLE_WITH_REASON`: fuera del proyecto con owner/razón.

## 8. Ventanas de revisión

El owner y riesgo determinan la ventana. Hechos temporales de seguridad, pricing, API, legislación, releases y soporte se verifican al usarlos aunque el registro diga vigente. Un manual con fecha de corte sirve como snapshot, no como prueba eterna.

## 9. Límites de autoridad

El agente puede corregir automáticamente errores técnicos demostrables dentro del scope y fuentes canónicas. Debe pedir decisión cuando la corrección cambia jurisdicción, términos comerciales, presupuesto, tratamiento de datos, riesgo aceptado, proveedor, plataforma o regla empresarial. Registrar el blocker no autoriza inventar la decisión.

## 10. Definition of done

### Identidad del runtime ejecutado — V317

Antes de gates dependientes del runtime, registrar versión observada, SHA-256 del
ejecutable seleccionado, versión del compilador cuando corresponda, raíz efectiva
y política de selección automática. Para Go: `go version`, `go tool compile -V`,
`go env GOROOT GOTOOLCHAIN`, ejecutable explícito y `GOTOOLCHAIN=local`; contrastar
el árbol con el archivo oficial fijado cuando se afirme integridad completa.
Después del build, conservar metadata del binario (`go version -m`) y su digest.
El mínimo del manifest y nombres como current-toolchain no identifican el runtime
ejecutado. Si no se registró, retirar la atribución exacta no demostrada y conservar
el resultado funcional con su límite; revalidar el head actual no certifica
retroactivamente todos los snapshots anteriores. Preservar bytes before/after y
actualizar refs/hashes sin reescribir logs históricos. Caso y evidencia:
`reconstruction_evidence/TOOLCHAIN_IDENTITY_AND_RUNTIME_SCA_V317.md`.

- cada claim material tiene source/version/date/applicability/owner;
- ninguna fuente móvil gobierna un artefacto reproducible;
- contradicciones y correcciones conservan before/after;
- locks, mapas, packs, templates y evidencia coinciden;
- gates dependientes pasan o el claim queda honestamente bloqueado;
- la siguiente fecha/trigger de revisión está definida.

### V319 — identidad de hijos y firma en Windows

No interpretar stdout completo de un package manager como process.execPath:
puede incluir progreso. Hacer que el hijo emita identidad estructurada a un
archivo nuevo, comprobar exit, versión, ruta y SHA; conservar stdout como log.
Si hay selección automática de package manager, observar también el entrypoint
real y su versión; para el gate fijar ese entrypoint y el Node explícitos.

Con gpgv de Git/MSYS, una ruta C:\ puede interpretarse como resource URL.
Ante ese fallo, conservar log y cualquier plaintext como NO_AUTENTICADO;
reintentar sólo la verificación con cwd y rutas relativas comprobadas, output
ausente, exit0 y VALIDSIG del fingerprint oficial. No inferir ausencia de clave
desde NO_PUBKEY si no se abrió el keyring. Comparar el plaintext autenticado
con los checksums oficiales antes de descargar/ejecutar el candidato.
Caso: reconstruction_evidence/NODE_RUNTIME_SECURITY_V319.md.

