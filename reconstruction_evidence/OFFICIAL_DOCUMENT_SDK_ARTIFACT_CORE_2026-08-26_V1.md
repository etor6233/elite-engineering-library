# Official Document SDK Artifact Core — 2026-08-26 V1

## Claim probado

El pack `OFFICIAL-DOCUMENT-SDK-ARTIFACT-CORE` 0.1.0 reconstruye un lock de cinco artefactos oficiales exactos de Microsoft Azure, Google Cloud y OpenAI y un adquiridor fail-closed. El código de adquisición/probe es `AUTHORED`; los wheels/sdist descargados son artefactos oficiales externos y no se atribuyen a Elite.

## Reconstrucción y pruebas

```text
Materialized 6 files into ...\elite-document-sdk-pack-verify-20260826
DOCUMENT_SDK_ACQUISITION_TEST_PASS negatives=4 offline_verified=1
DOCUMENT_SDK_LOCK_VALID artifacts=5 selected=1
```

La suite rechazó ID desconocido, ID duplicado, adquisición sin approval y cache alterada. El positivo offline verificó bytes/hash, creó receipt y no dependió de red. El materializador comprobó manifest↔FILE y SHA-256 de los seis archivos.

El perfil previo también quedó ejecutado por el verificador global: `PROJECT_INITIALIZATION_PACK_PLAN.md implementation_files=18`. La primera composición descubrió que dos parámetros runtime estaban declarados como variables de sustitución; se corrigió el pack upstream a 0.4.2 y la repetición completa pasó (`LIB-FAIL-033`).

## Artefacto oficial real verificado

Se importó desde cache el source distribution oficial Microsoft `azure_ai_contentunderstanding-1.1.0.tar.gz` usando el runner reconstruido:

```text
DOCUMENT_SDK_ACQUISITION_PASS artifacts=1
bytes=230330
sha256=00f39c7cf2ba50c8d4586315f4295a45fad0881daaff07d13d33aea2bafc7ec8
source=verified-offline-cache
```

El approval utilizado declara explícitamente que es un fixture de verificación de Codex, no aprobación productiva del propietario. Ninguna cuenta, endpoint, secreto ni llamada facturable fue utilizada.

## Límites preservados

- identidad del artefacto no equivale a precisión documental;
- el grafo transitivo se fija por Python/OS del proyecto antes de instalar;
- cuenta, región, processor/analyzer/model, cuotas, costos y corpus siguen bloqueados hasta decisión y acceso del usuario;
- imports/probes no sustituyen contract tests ni evaluación por campo;
- no se afirma Trusted Publishing para uploads que PyPI no atestiguó.

Fecha/revisor: 2026-08-26 / Codex.
