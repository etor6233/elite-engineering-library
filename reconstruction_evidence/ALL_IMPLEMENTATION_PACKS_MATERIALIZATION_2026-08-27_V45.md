# All Implementation Packs Materialization — 2026-08-27 V45

## Alcance del snapshot

Este snapshot sucede a V44 sin reescribirlo. Conserva 46 packs, 413 archivos materializables y 23 perfiles; amplía el lock a 81 fuentes oficiales con Microsoft Agents Human Oversight como referencia exacta rechazada, sin seleccionarla ni atribuirle autenticación que no implementa.

Estado canónico al cierre:

- 46 implementation packs;
- 413 bloques `CREATE`: 401 `AUTHORED`, 3 `ADAPTED`, 9 `VERBATIM`;
- 23 perfiles de composición y 10 source profiles;
- 81 upstream source archives fijados por identidad, bytes, SHA-256 y licencia;
- 5 artefactos SDK documentales y 1 artefacto Business Central;
- 9 provider adapters, 1 orquestador documental y 1 log de evidencia;
- 356 fallos locales y 107 condiciones upstream preservados al corte;
- 301 Markdown y 8 scripts PowerShell canónicos después de incorporar este expediente.

## Fuente 81

`microsoft/agents-humanoversight` quedó fijado en el commit Microsoft firmado `9eca38d78f7d26027a925ffbf99088d74e7ebcdd`:

- ZIP por commit: 2.301.334 bytes;
- SHA-256: `b82c5b221bba29207f158315ebfafea7f51bce46bcf4f0043e18527cf6fd6b00`;
- MIT SHA-256: `c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383`;
- README, requirements, decorador, servicio, dos Bicep y contrato live fijados por path/hash;
- commit verificado por GitHub: `valid`;
- repositorio archivado por Microsoft el 2026-08-24.

El código oficial aporta email Office 365 con opciones, responder, correlation ID, timeout y log Table Storage. No aporta una frontera productiva: el trigger HTTP no implementa Entra; el README manda tratar la URL como secreto y considerar agregar Azure AD. El caller elige emails/correlation ID y no se demuestran resource authority, nonce/replay, CAS, idempotencia ni transacción entre decisión y efecto.

## Ejecución upstream

En CPython 3.12.13 aislado, una resolución audit-only produjo 89 paquetes. La suite oficial terminó:

```text
28 passed
2 skipped
1 warning
```

Los dos tests live requieren Logic App/email/interacción y no se ejecutaron. OSV Scanner 2.5.1 sobre el freeze instalado produjo 16 advisory IDs en cuatro paquetes. El repositorio queda `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`; ningún perfil lo selecciona.

## Pack de adquisición

`OFFICIAL-UPSTREAM-ACQUISITION-CORE` sube a `0.4.36`. Los 18 planes consumidores fueron alineados. Desde materialización limpia:

```text
Materialized 19 files
UPSTREAM_ACQUISITION_TEST_PASS negatives=7
UPSTREAM_LOCK_VALID sources=81 selected=81
SOURCE_PROFILE_TEST_PASS valid=10 negatives=5 positives=1
```

Las selecciones se mantienen 6, 1, 10, 7, 16, 35, 1, 25, 12 y 1; source 81 no aparece en ninguna.

## Auditoría global

Antes de incorporar este snapshot:

```text
VERIFY_LIBRARY_PASS
packs=46 materialized_files=413 markdown_files=300

VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=46 upstream_sources=81 document_sdk_artifacts=5 business_central_artifacts=1 provider_adapters=9 document_orchestrators=1 evidence_logs=1
```

El gate final posterior debe producir:

```text
VERIFY_LIBRARY_PASS
packs=46 materialized_files=413 markdown_files=301
```

## Lecciones preservadas

`LIB-FAIL-349` a `LIB-FAIL-356` retienen errores reales de sintaxis PowerShell, formato de lock OSV, paths/parametrización supuestos, contexto de patch, formato de archive, selección de root y array de updater. `UP-FAIL-105` a `UP-FAIL-107` conservan archivo/auth, supply chain y falta de evidencia live/consistencia del candidato.

## Estado honesto

La brecha de revisión humana autenticada/autorizada de punta a punta permanece abierta. V45 evita que un agente adopte por error un repositorio oficial cuyo título promete human oversight, pero no inventa el adapter faltante ni lo atribuye a Microsoft.

No existe autorización para Git, ZIP final, credenciales, cuentas, despliegue ni servicios pagos.

## Fuentes oficiales

- <https://github.com/microsoft/agents-humanoversight>
- <https://github.com/microsoft/agents-humanoversight/tree/9eca38d78f7d26027a925ffbf99088d74e7ebcdd>
