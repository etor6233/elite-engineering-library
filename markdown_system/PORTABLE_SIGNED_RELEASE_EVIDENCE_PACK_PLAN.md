# Portable Signed Release Evidence Pack Plan

Perfil de release local para proyectos Windows que necesitan un artefacto reproducible, SBOM, SCA, provenance y firma verificable sin Git, GitHub Actions ni servicios pagos. No genera ni almacena claves; el proyecto aporta una identidad Ed25519 protegida y una política pública explícita.

```json
{
  "compositionVersion": "1.0",
  "secretVariables": {},
  "packs": [
    {
      "path": "implementation_packs/PORTABLE_SIGNED_RELEASE_EVIDENCE_GATE.md",
      "packId": "PORTABLE-SIGNED-RELEASE-EVIDENCE-GATE",
      "version": "0.1.0",
      "acknowledgeConditions": true,
      "files": ["*"],
      "variables": {}
    }
  ]
}
```

Secuencia: materializar 10/10 → verificar vector OpenSSH/Authenticode → adquirir OSV exacto desde el lock existente → completar inputs, notices, build y signer policy → dos builds → SCA 0 → SPDX 2.3 → in-toto/SLSA → firma/verify → verificación independiente → preservar por digest. Deployment, recuperación, carga, seguridad ofensiva y aceptación del target siguen siendo gates separados.
