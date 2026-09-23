# Portable signed release evidence gate

This gate produces a release candidate without Git, a hosted CI service, or paid signing infrastructure. It executes one hash-locked project build twice, requires byte-identical artifacts, generates a source manifest, scans the source with the exact Google OSV-Scanner 2.5.1, emits SPDX 2.3, creates an in-toto Statement v1 with the SLSA provenance v1 predicate, signs that statement through OpenSSH `sshsig`, immediately verifies the signature against an explicit allowed-signers file, and writes a hash-linked receipt.

The project supplies the business code, an existing release signing key, the public allowed-signers policy, an explicit input manifest, dependency manifests, notices, and the build script. The private key is never copied into the release directory or receipt. A pass proves only the recorded local candidate; it does not prove cloud deployment, runtime security, legal acceptance, production data recovery, or business acceptance.

Run `verify_pack.ps1` first. Then copy and complete `release-profile.template.json` and run:

```powershell
pwsh ./build_release.ps1 -ProjectRoot <project> -Profile <profile.json> -ReleaseRoot <new-directory> -OsvScanner <exact-osv.exe> -SigningKey <existing-key>
pwsh ./verify_release.ps1 -ReleaseRoot <new-directory> -TrustedAllowedSigners <independently-trusted-public-policy>
```

The release directory is create-only. Altering the artifact, manifest, provenance, signature, signer policy, SBOM, SCA report, or receipt makes verification fail. `ssh-keygen` must be an Authenticode-valid Microsoft Windows component when `require_microsoft_authenticode` is enabled. OSV must match the locked 2.5.1 executable hash.

Authorities are pinned in the pack metadata: SLSA v1.2 provenance, in-toto attestation v1.2.0 and OpenSSH portable `V_10_2_P1`. The four files under `fixtures/` are Base64 transport wrappers that reconstruct the exact public OpenSSH regression bytes; no private fixture key is included.

V402323 contract: supply an independently trusted public allowed-signers file
outside the received release directory. A copy carried by the artifact is never
its own trust anchor. The verifier checks its exact hash, identity and the fixed
elite-release-v1 namespace. This public policy contains no private key.

The builder scans exactly dependency_manifest_refs with the pinned OSV executable;
the caller must enumerate the complete selected dependency manifest set, including
nested Go/Node/.NET owners when selected. Source input SHA/size/reparse checks run
before and after both builds and again after scanning. Windows ZIP paths use
single-separator normalization and the verifier rejects backslashes/unsafe parts.
These checks do not turn a local candidate into a production release.

Failed release staging is retained with an explicit path for diagnosis; it is never
published as PASS. Successful staging is cleaned only after exclusive publication.
The consumed profile hash is frozen and checked again before signing.
