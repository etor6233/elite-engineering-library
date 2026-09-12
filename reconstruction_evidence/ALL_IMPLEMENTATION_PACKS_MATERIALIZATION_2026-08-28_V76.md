# All Implementation Packs Materialization — V76

Fecha: 2026-08-28  
Estado: `VERIFY_LIBRARY_PASS` + `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`

## Resultado gobernante

- 55 implementation packs y 519 archivos materializables desde Markdown;
- 380 archivos Markdown después de incorporar esta evidencia;
- 32 perfiles de composición ejecutados; `AWS_DURABLE_OBJECT_EVENT_WORKER_PACK_PLAN.md` compone 3 packs/43 archivos;
- source lock global permanece en 110 fuentes oficiales/14 perfiles; el nuevo componente conserva además su source lock exacto independiente;
- ledger: 850 fallos locales + 154 condiciones upstream = 1.004 IDs únicos esperados, cero duplicados y cero lecciones locales pendientes.

## Código oficial incorporado

Fuente: `aws/aws-durable-execution-sdk-python`, tag compuesto `sdk-v1.7.0,otel-v0.3.0`, commit firmado/verificado `075b65aacb80de8bb1507e3df2e53cc90cb3b874`, Apache-2.0 + NOTICE. Archive: 889.895 bytes/SHA-256 `d3e42423e61e258191fdd078a29d1fd58cfc55380b25a357f879b73665a02d11`. Core wheel 1.7.0: 109.716 bytes/SHA-256 `1482e1f439e36deb1fcf9a086499b3788c3facf3ac755e50376a739081cc72a0`; sdist SHA-256 `6ce773cfba243df98ec35f375a30b98cc4779a9f921d8c6d1180773811bf57f0`.

`AWS_LAMBDA_DURABLE_EXECUTION_COMPONENT.md` materializa 12 archivos: cuatro de packaging Elite `AUTHORED` y ocho AWS `VERBATIM` —LICENSE, NOTICE y pares source/test de `step_with_retry`, `step_semantics_at_most_once` y `replay_logging`. Los cuatro contratos offline pasan desde una reconstrucción limpia.

El perfil combinado agrega el componente AWS Powertools 3.34.0 ya admitido para idempotencia por registro + partial batch SQS. No existe un archivo integrado atribuido a AWS: cualquier dispatcher, selección de execution name y composición empresarial pertenece al proyecto y debe declararse `AUTHORED`.

## Ejecución upstream

- inventario del commit: 376 archivos, 332 Python, 168 archivos de test, 146 archivos bajo examples y nueve workflows; no publica dependency lock;
- Windows Python 3.14 source-exact focal: 112/112 PASS;
- Windows integral core + testing + examples: 2.791 PASS, 3 FAIL de portabilidad y 2 skips en 257,13 s;
- los tres nodeids Windows son path POSIX supuestamente no escribible y timestamps anteriores a epoch no soportados por Windows;
- WSL2 Linux/Python 3.14, source exacto y wheels runtime: los 112 focales más esos tres nodeids terminaron 115/115 PASS en 29,66 s;
- auditor global: `aws_durable_execution_components=1`, pack 12/12 y contratos 4/4 PASS.
- SCA: OSV-Scanner 2.5.1 escaneó los ocho pins del runtime lock byte-idéntico (SHA-256 `e38a6df20f235238a215b196f50e215dce9c4d82a265d60854193815659f8ffb`) y devolvió cero vulnerabilidades conocidas; reporte JSON 1.553 bytes/SHA-256 `70ec56bd4e623fcdd57873d0923de131c26b7d36a3c2a5f2eeacebe6b832b99b`.

## Condición upstream conservada

El source testing del tag sigue declarando 1.2.1, pero `execution.py` mide 29.720 bytes/SHA-256 `abc2785722246de1c386a91af2d3e9d9a0f9692f8f03d91234a1f2688b697ee9` y define `OperationPaginatorState`. El wheel testing PyPI 1.2.1 mide 101.233 bytes; su archivo homónimo mide 17.581 bytes/SHA-256 `4269ee0b876fe36a5b3d8ddb69b4f75613a775405e29c27ababbcca3a82b8917` y no define ese símbolo. La condición upstream 154 prohíbe combinarlo silenciosamente con los tests del tag. El pack no redistribuye ese wheel.

## Recuperación de fallos

- 838 corrigió paths/globs supuestos de adquisición mediante inventario real;
- 839 resolvió que el tag 1.7.0 es compuesto;
- 840 corrigió una lista PowerShell concatenada como un solo target;
- 841 probó el límite Windows de path 260 y reconstruyó desde raíz corta;
- 842 restauró `--import-mode=importlib` requerido por el monorepo;
- 843 clasificó y reprodujo en Linux los tres fallos Windows;
- 844 y 845 evitaron quoting cruzado y `apt/sudo`, usando wheels puros bajo temporales;
- 846 detectó la deriva real del testing wheel y la promovió como condición upstream;
- 847 corrigió la firma del materializador y demostró 12 archivos/4 contratos.
- 848 separó una invocación SCA bloqueada por política de cualquier resultado del código;
- 849 corrigió la detección del `.lock` mediante un `requirements.txt` de staging con SHA-256 idéntico y demostró ocho paquetes/cero vulnerabilidades conocidas.
- 850 descartó un conteo auxiliar que confundía referencias narrativas con filas y repitió el patrón exacto del verificador oficial.

## Límites vigentes

El componente no cierra exactamente-una-vez ni autoriza producción. Aún se requieren dispatcher con execution name estable, identidad empresarial tenant-bound, S3 version/eTag/sequencer, efecto idempotente, owner/fencing si hay reclaim, cuenta/región/IAM/IaC, DLQ/redrive, reconciliación, observabilidad, SCA del día, carga, recovery/rollback y E2E cloud. Los samples S3 endupe, SQS best-practices y DLQ replay antiguos continúan rechazados; no se usaron para rellenar esos límites.
