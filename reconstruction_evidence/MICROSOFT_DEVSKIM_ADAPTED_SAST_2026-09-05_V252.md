# Microsoft DevSkim adapted SAST baseline — reconstrucción V252

## Decisión

`MICROSOFT-DEVSKIM-ADAPTED-SAST-GATE 0.1.0` queda `ELITE_REFERENCE / REBUILD_VERIFIED / CONDITIONED` para el claim estrecho de security linting local sin costo sobre Go y TypeScript. No es código Microsoft sin modificar, SAST interprocedural completo ni prueba de que un sistema sea seguro.

El upstream oficial continúa `REJECTED_COMPONENT / UPSTREAM_OPEN`: Microsoft DevSkim `a452aa5…` resuelve SharpCompress 0.40.0 afectado. V252 crea una lane distinta y honesta: agrega referencias directas 0.48.0 bajo clasificación `ADAPTED`, conserva source/licencias/notices y ejecuta nuevamente toda la evidencia material.

## Autoridades y estado actual

- Microsoft DevSkim oficial: <https://github.com/microsoft/DevSkim>.
- GitHub CodeQL CLI: <https://docs.github.com/en/code-security/concepts/code-scanning/codeql/codeql-cli>. Para repositorios privados de organizaciones Team exige además GitHub Code Security; no se presume entitlement.
- GitHub source archive stability: <https://docs.github.com/en/repositories/working-with-files/using-files/downloading-source-code-archives>. GitHub garantiza contenido por commit, no bytes idénticos del contenedor comprimido.
- Microsoft commit API observado el 2026-09-05: `a452aa506f80928b611c4c7bfe306b8115821aae`, tree `f6be283e73f106afdb1345c4123539a7a2f9a136`, `verified=true`, fecha `2026-08-16T15:47:06Z`; release actual `v1.0.90`, publicada `2026-07-17T18:58:08Z`.
- SharpCompress 0.48.0 exacto: nupkg 9.262.981 bytes, SHA-256 `d8c5da8a76d325eb81c1103a78953e025513f22ade36b5b11d8342324146f0b7`, MIT, repository commit `6e59c7d7bbf8c19a8a92c3c382599906684bb93d`.

## Identidad reproducible

Dos descargas del mismo archive por commit produjeron transportes distintos —935.159/949.833 bytes y SHA-256 `2e75…`/`c3b7…`— pero exactamente 268 archivos y el mismo árbol canónico `4914ee46a553ee3c80b7e4849cbbcc349f1deab99951373229f177c88ad43ecb`. El runner consulta commit/tree/firma, exige una raíz única segura, hashea cada `path|bytes|sha256`, rechaza traversal y registra el digest de transporte sólo como receipt.

La adaptación cambia exclusivamente los proyectos library/CLI para resolver SharpCompress 0.48.0. Hashes finales reproducidos por el mismo productor:

- library csproj: `a5439c44f6d5dafae269185dfd090816c43f37c08f81d9283432f20b913b6016`;
- CLI csproj: `9ffdc50a3fb44861d5cd6987ce805f3f4ef250b610ffe5339c47c3892e71a085`.

## Evidencia ejecutada

Con .NET SDK oficial 10.0.400 Windows x64:

```text
restore NuGet.org exclusivo: PASS, cero warnings NU
build Release net10.0: PASS
Microsoft.DevSkim.Tests: 300 passed / 0 failed / 0 skipped
dotnet package vulnerability reports: 3 proyectos / 0 findings
resolved SharpCompress runtime: 0.48.0 exacto
TypeScript hostile fixture: DS189424 exacto
Go hostile fixture: DS112852 exacto
clean TypeScript fixture: 0 findings
static parser/contract: PASS
fail-closed negatives: 2/2 PASS
materialization: 4/4 PASS
```

La auditoría integral de la biblioteca también ejecutó esta lane desde cero y terminó con código 0:

```text
VERIFY_EXECUTABLE_LIBRARY_PASS
mode=Audit
packs=157
upstream_sources=121
microsoft_devskim_adapted_sast_gates=1
MICROSOFT_DEVSKIM_ADAPTED_SAST_PASS tests=300 vulnerabilities=0 typescript=DS189424 go=DS112852 clean=0
MICROSOFT_DEVSKIM_ADAPTED_PACK_PASS static=1 negatives=2 runtime=1
```

El publisher conserva runtime, `MICROSOFT_DEVSKIM_LICENSE.txt`, notice, licencia SharpCompress, `SOURCE_LOCK.json`, `PROVENANCE.json`, manifest por archivo y receipt atómico no sobrescribible.

## Hashes de biblioteca

- pack: `254b0c50971ca1509fc2f1111ff38b6a8be1d12f1d211f2b3c8e8ee45a2e4192`;
- plan: `f5f4f4919b6508a162971dd08552cd746ef3bc138d51cd4a803025b84d6fc5b7`;
- `source-lock.json`: `c730eec041250a68bfe61aad33630233d473d9706050f48793e5d6ba00291b6d`;
- runner: `1bea4f07ec53cb07e4e41f8c19df9ed96352e3365b42f1c19c1ff20cb0a88e35`;
- verifier: `e62d3646dcf9a6fed4af336987f549ce25ae583e3ba39b728205b07dcc2a9ea2`;
- README: `1a9791b25bb7db555c84b2e874da8b8e3dac1dd4eefb75a8ce83362de463a84b`.

## Límites y reapertura

DevSkim encuentra patrones de reglas; no demuestra dataflow/interproceduralidad universal. Cada proyecto debe fijar scope, exclusiones, severidades, baseline, suppressions con vencimiento, triage y regresiones, y ejecutar además fuzzing, DAST, threat model, revisión y seguridad ofensiva. Reabrir ante release/commit Microsoft nuevo, advisory NuGet, cambio de SDK/support, cambio de reglas o claim más amplio.
