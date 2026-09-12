# Official Maintenance Tool Artifact Lock

Fecha de corte: 2026-08-25. Este lock fija herramientas oficiales que apoyan mantenimiento; no autoriza cambios automáticos ni convierte un scan en garantía de seguridad.

## Google OSV-Scanner 2.5.1

```yaml
tool_id: google-osv-scanner
version: 2.5.1
release: https://github.com/google/osv-scanner/releases/tag/v2.5.1
source_commit: c84fa4568f2526d0333e9a914ea8a0a5f74ad68b
license_expression: Apache-2.0
source:
  url: https://github.com/google/osv-scanner/archive/c84fa4568f2526d0333e9a914ea8a0a5f74ad68b.zip
  bytes: 13461636
  sha256: 688ceb4ab62f38b5fd8fb4c49a9220e7589d49b88eee422aa3061ab6e339995a
  license_sha256: cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30
  go_sum_sha256: 660d54260464f85ae1800b9ee17ca8dbd440270e7cea3a98d2fd417e4a50b331
windows_amd64:
  url: https://github.com/google/osv-scanner/releases/download/v2.5.1/osv-scanner_windows_amd64.exe
  bytes: 58970112
  sha256: 25e42f5ef6711fd8c0fb45390972205891dd44c6bd02ac93f0f63e8e98d9bfb6
checksums:
  url: https://github.com/google/osv-scanner/releases/download/v2.5.1/osv-scanner_SHA256SUMS
  bytes: 554
  sha256: 40120ab7d670199f6663d76f9806b6d11ab3c402df808e5cc4e13c5cb0643356
provenance:
  url: https://github.com/google/osv-scanner/releases/download/v2.5.1/multiple.intoto.jsonl
  bytes: 24443
  sha256: 664273601759eb6540f35bfdbfa6ff083026c09e2e9bb7ad3cded795967dd462
```

Verificación local del binario:

```text
osv-scanner version: 2.5.1
osv-scalibr version: 0.5.2
commit: c84fa4568f2526d0333e9a914ea8a0a5f74ad68b
built at: 2026-08-17T03:44:26Z
```

El SHA-256 del binario coincide con `osv-scanner_SHA256SUMS`. La presencia de `multiple.intoto.jsonl` queda fijada por hash; verificar criptográficamente su statement/trust root con SLSA verifier compatible antes de una promoción productiva, no sólo su integridad de bytes.

## Política de uso

1. Descargar por URL exacta a un directorio nuevo sin secretos.
2. Verificar bytes y SHA-256 contra este lock y contra el checksum oficial descargado.
3. Verificar provenance/signature cuando el target disponga del verifier y trust policy.
4. Ejecutar `scan`, no `fix`, sobre manifests/locks/SBOM del artefacto exacto.
5. Guardar JSON, fecha, versión de scanner y hash del input/output.
6. Tratar cada finding como candidato de triage: exposición, reachability, impacto, exploit maturity, fixed version, compatibilidad y rollback.
7. `fix` permanece prohibido por defecto porque puede ejecutar package managers/scripts o modificar manifests. Sólo se habilita en workspace desechable, sin credenciales, y vuelve a todo `DEPENDENCY_UPDATE_CONTRACT.md`.
8. Ignorar un finding exige ID, razón, evidence, owner y vencimiento; un ignore permanente no es válido.
9. Un resultado vacío sólo significa “sin vulnerabilidades conocidas por la base consultada en ese momento para los packages detectados”. No prueba ausencia de vulnerabilidades, malware o componentes no detectados.

## Plataformas adicionales

La release oficial publica binarios Darwin/Linux/Windows amd64/arm64. Antes de usar otra plataforma, agregar sus bytes/SHA desde `osv-scanner_SHA256SUMS`, verificar provenance y ejecutar `--version`; no inferir sus digests desde el binario Windows.

### V319 — Node24.20.0 Windows x64, candidato local acotado

Fecha2026-09-08. No es runtime productivo promovido ni se cambia el global.
Fuente https://nodejs.org/dist/v24.20.0/win-x64/node.exe;
commit71b8b174857e25106d39b61a9e6f30d927da8b01,93381448 bytes,
SHA2565c976096e04e5c2c1f091938926234cc9fbebfe9787ddd149351b3b0ecc707b5.
LICENSE íntegra de la distribución: SHA256
ed34dd8e3f0a78dbaf00d0444ce8e285b015b765379c2e17880455f70370f8e9;
MIT más los notices/licencias de componentes embebidos descritos allí.
Verificar SHASUMS256.txt.asc con fingerprint
5BE8A3F6C8A5C01D106C0AD820B1A390B168D356 y keyring oficial nodejs/release-keys
5b7f55f4a7e35d1176d27a6b81b0c3c3b794216b,
SHA2562733c3959ac2843194801dd7a17162439a996a1471e007624c27f2dda30206a6.
La selección local AUTHORED contiene sólo el node.exe oficial VERBATIM y esa
licencia; pnpm11.19.0 es herramienta separada y debe verificar su propio grafo.
El ZIP6cac9ffbca8f6a47091e4b5c772e0606049c3871cb67d900c0cedde630e545ba
queda rechazado como distribución completa por9 advisories de npm; no corregido.

Probe oficial evaluado: nodejs/is-my-node-vulnerable1.6.1,
commitc37a56bad56e34fe5223ddd3cb223cc4158136ae, MIT, adaptación de2 URLs
de datos a commits fijos. EngineSHA256
9ea3eb3ca804e918ef53968299090cf62b93a77f458f41ac40be253d4570c408;
semver7.8.5 ISC desde la distribución firmada. Sólo CLI, sin Actions entrypoint.
G0–G8 CONDITIONED, implementation_ready=false: no hay pack reusable admitido.
No ejecutar automáticamente el CLI sin los controles del expediente.
Evidencia, before/after y límites: reconstruction_evidence/NODE_RUNTIME_SECURITY_V319.md.


### V320 — adapter offline Node0.1.0 y semver byte-exacto

El pack NODE_OFFICIAL_RUNTIME_ADVISORY_GATE.md contiene el source-lock completo
del motor1.6.1/c37a56bad56e34fe5223ddd3cb223cc4158136ae,3 archivos de fuente/licencia,
bridgeAUTHORED y snapshots fijados. Sólo getJson cambia a lectura local verificada.
Tarball semver7.8.5: https://registry.npmjs.org/semver/-/semver-7.8.5.tgz,
SHA256d85045d4300d7d57c891336b95df532e73f34c22ffcd222452b6d08b9d127d5d,
commit6e05b7637396ac66522cff8731f07cfe0ef49a29,ISC;
LICENSEsha2564ec3d4c66cd87f5c8d8ad911b10f99bf27cb00cdfcff82621956e379186b016b.
No reutilizar como byte-exacta la copia de Node:53 archivos difieren por formato
(52 sóloCRLF/LF,README whitespace). El runner verifica53 hashes del tarball sin
normalizar. Nodeexe+LICENSE conservan pins V319. No nueva versión npm/pnpm
admitida, sourceprofile de producto ni instalación global. Evidence V320.

## V321 — pnpm11.25.0 candidate selection

Supersedes11.19.0 for three active consumers after bundle471/1 finding. Official registry tarballSHA256:33dd0748f27e7916c4f1c8b6943461983e3453b06bbda6312a6280130b4881e5;455files/473packages/0SCA,MIT root plus473license declarations, exact published notices. Commit6d90c71efdffbc909b499490b64c66badc720327 verified by SLSA identity policy; npm ECDSA/publish also verified. Only directentry+V319 standaloneNode24.20.0 was tested. Entrybootstrap hash is identical across versions; require whole published artifact/bundle hash, not entry hash alone. Retain6unchangednative binaries as a separate coverage condition. FAIL525 profile acquisition pending, no globalupdate. See PNPM_BUNDLE_SECURITY_V321.md.

## V322 — adquisición de mantenimiento comprobada

Core0.4.78 / initialization37files; narrow G0–G8 USE_REUSABLE_PACK con owner global CONDITIONED. Perfil materializado maintenance-runtime-artifacts, lock124 y3sources adquiridos por HTTPS con approval actual del agente bajo autorización previa; no aprobación retrospectiva. Node24.20.0 win-x64, pnpm11.25.0 y Sigstore4.1.1: bytes/licencias exactos y10outputs inmutables tras PRESENT. FAIL525/527/528/529/530/531 cerrados focalmente, TEST41/EVID32. Biblioteca161packs/1453files/757Markdown/52profiles y151preflight steps PASS; preflight global BLOCKED por3tools no resueltos y readiness/target separados. Ver GOVERNED_MAINTENANCE_ACQUISITION_V322.md y PROJECT_OFFICIAL_SOURCE_PROFILE_RECORD.md. Continuar6binarios nativos de pnpm y T2809; no scheduler, ZIP, install global, ejecución de artefactos adquiridos ni producción.

## V323 — cobertura nativa explícita

Reflink0.1.19 commit2223f06f:75Cargo identities exactas/0avisos OSV2.5.1;4addons contienen rustc f6e511ee→1.82.0.8statements npm/SLSA verificados/9negativos;18RustSec std y4avisos oficiales std revisados sin afectar1.82.0. No tarball nativo nuevo, build, ejecución ni SBOM del binario. Fastlist0.3.0 CRT/build no fijados; full license text ausente en ambos trees inspeccionados. FAIL532/LIB-FAIL2248 BLOCKED_EXTERNAL para admisión integral nativa y redistribución; no vulnerabilidad inventada. TEST42/EVID33, PNPM_NATIVE_COVERAGE_V323.md. Continuar T2809 por owners y conservar target/readiness pendientes.
