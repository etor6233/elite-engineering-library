# Portable Signed Release Evidence — 2026-09-05 V254

## Resultado

`PORTABLE-SIGNED-RELEASE-EVIDENCE-GATE 0.1.0` materializa 10/10 archivos: dos builds, ZIP determinista, source/artifact manifests, SCA JSON, SPDX 2.3, in-toto Statement v1 con SLSA provenance v1, firma Ed25519 OpenSSH, verificación independiente y receipt hash-linked. No requiere Git, GitHub Actions, cuenta cloud ni servicio pago. El proyecto aporta su build, inputs, dependency manifests, notices, allowed signers y una clave protegida fuera del árbol.

## Autoridades exactas

- SLSA spec `v1.2`: commit verificado `19e4e2f005f871270c4f555fc47afecfb37f3efe`, tree `851eab11f4234c9ef75163bad6918daa69bd15dc`; build provenance SHA-256 `182ab5cf823c374874136371466ba43f99f0a9ba87a8046757d26a2ed74854a0` y CUE `0ddd00372622d1c04f71d1fb2666579af089eb733c9c5d838d158489b2ea6145`.
- in-toto attestation `v1.2.0`: commit verificado `df02077bf97218a8860a5c534eff1f1381f56984`; Statement SHA-256 `cbe684a18b812b8b613d9202eb43b2ea24477f91a2ad6ca5be935185a455ebea`, envelope `c02c65880ccc117bbdecc6c916d5108541c1a864de67f8f9082efbaf15265d11`.
- OpenSSH portable `V_10_2_P1`: tag y commit verificados `d01efaa1c9ed84fd9011201dbc3c7cb0a82bcee3`, tree `c91beec59da56218afced35b14087d362b5bba8e`. Cuatro blobs públicos `sshsig` se transportan por Base64 y reconstruyen hashes exactos; ninguna clave privada entra al pack.
- Microsoft Windows OpenSSH observado: `OpenSSH_for_Windows_9.5p2`, file version `9.5.6.1`, `ssh-keygen.exe` SHA-256 `44c6809b7bbc917f1310ba92857f983e2788e9b0015aa7896fa0362eddb6338b`, Authenticode válido de Microsoft Windows.
- Google OSV-Scanner 2.5.1: commit `c84fa4568f2526d0333e9a914ea8a0a5f74ad68b`, Windows executable 58.970.112 bytes/SHA-256 `25e42f5ef6711fd8c0fb45390972205891dd44c6bd02ac93f0f63e8e98d9bfb6`.

## Reconstrucción funcional

Un fixture Go sin dependencias se construyó dos veces desde un script hash-locked. Ambos ZIP dieron SHA-256 `c6a018ae452644877c6f4c37b1acb36f753a05c2d110cef3ff89dac0da199fe3`. OSV reportó cero vulnerabilidades y generó SPDX 2.3. La provenance se firmó con la clave privada de regresión pública de OpenSSH, protegida por ACL y ubicada fuera del proyecto; `verify_release.ps1` pasó. Una copia con un byte alterado fue rechazada por `artifact hash mismatch`. La publicación se cambió a staging hermano + move final para no dejar un release parcial ante fallo.

## Candidatos rechazados

- Microsoft SBOM Tool `v4.1.5` funcionó y detectó tamper, pero su SBOM oficial produjo 12 filas afectadas/18 advisories; `main` exacto `4091b7bcce1640c4db42b5fad63d7d1b7bc0e4cf` todavía mostró `NuGet.Packaging` y `NuGet.Protocol` 7.3.0 afectados. No se incorpora.
- Sigstore Cosign `v3.1.3` se auto-verificó contra su bundle oficial y probó attestation/tamper, pero su SPDX oficial produjo 9 filas afectadas/164 advisories. No se incorpora. Esta evidencia reabre y bloquea el verifier del pack OpenGrep/GitLab V253.
- in-toto-golang `v0.11.0`, que también es el `master` actual exacto, produjo 5 filas afectadas/44 advisories. No se incorpora.

## Verificación integral

- Round-trip Markdown: 10 archivos, cero diferencias.
- `verify_pack.ps1`: vector oficial positivo, payload alterado rechazado y Authenticode válido.
- `VERIFY_LIBRARY_PASS`: 159 packs, 1.390 archivos, 687 Markdown, 51 perfiles; franquicia 689 archivos desde 65 packs.
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: Go oficial 1.26.7 explícito, 121 fuentes, 16 adapters, `portable_signed_release_evidence_gates=1`; OpenGrep runtime informado como `BLOCKED`, no ejecutado.

## Límites

El pack prueba el candidate release local y su evidencia. No certifica CDN/WAF, IdP, provider accounts, PostgreSQL/recovery, carga, seguridad ofensiva, despliegue, rollback live ni aceptación empresarial. Cada proyecto debe demostrar esas condiciones con el artefacto exacto firmado.
