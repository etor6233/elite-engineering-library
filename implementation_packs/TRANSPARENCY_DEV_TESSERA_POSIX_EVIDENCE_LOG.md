# Transparency.dev Tessera POSIX Evidence Log

## 1. Metadata

```yaml
pack_id: "TRANSPARENCY-DEV-TESSERA-POSIX-EVIDENCE-LOG"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa el código oficial de Transparency.dev Tessera para construir y verificar un log local append-only con checkpoints y pruebas Merkle, más glue fail-closed que sólo admite recibos documentales canónicos sobre Linux/POSIX."
stacks: ["Go 1.26.7", "Transparency.dev Tessera signed commit a8f33c56", "Linux amd64", "POSIX filesystem"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "MICROSOFT-DURABLE-DOCUMENT-ORCHESTRATION 0.1.x", "PYTHON-AZURE-BLOB-IMMUTABLE-EVIDENCE-ADAPTER 0.1.x"]
incompatible_with: ["Windows or NTFS storage", "document bytes or extracted values as log entries", "unverified source", "unwitnessed production claim", "automatic key generation", "silent dependency pin changes"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0"
upstream_sources: ["https://github.com/transparency-dev/tessera/tree/a8f33c56b80be808712c0537cb1b71f2fd5846ea", "https://github.com/google/trillian", "https://c2sp.org/tlog-tiles"]
verified_at: "2026-08-27"
```

## 2. Applicability

Use cuando el proyecto necesite detectar modificación del historial de evidencias documentales sin depender de un servicio pago. El log recibe únicamente recibos canónicos con hashes y referencias opacas después de seguridad, extracción, evaluación, revisión y persistencia. No almacena documentos, campos, PII ni secretos.

El upstream `AUTHORS` enumera Google LLC, Anthropic PBC e Internet Security Research Group. Google Trillian recomienda Tessera para logs nuevos, pero el repositorio se atribuye a Transparency.dev/Tessera authors y no se renombra como producto Google. La implementación POSIX exige Linux y semántica de filesystem real; Windows/NTFS y montajes Windows de WSL quedan rechazados.

## 3. Architecture contract

El source fijado es el commit firmado `a8f33c56b80be808712c0537cb1b71f2fd5846ea`, posterior al release estable `v1.0.4`. Los CLI oficiales `posix-oneshot` y `fsck`, el diseño POSIX, `LICENSE` y `AUTHORS` se conservan byte a byte. El resto es glue local declarado AUTHORED.

Build y ejecución fallan cerrados ante identidad, hash, toolchain, plataforma, filesystem, binarios, claves, recibos, suite, vet, análisis alcanzable o `fsck` inválidos. El template no habilita nada. El log local es tamper-evident, no deletion-proof: producción exige publicar checkpoints independientemente o usar witnesses, controlar acceso, backup/restore, retención, monitoreo y responsables.

El grafo exacto del commit contiene cinco registros OSV en tres módulos requeridos pero no llamados por el CLI POSIX; `govulncheck 1.7.0` encontró cero vulnerabilidades alcanzables y cero en paquetes importados. Ese resultado condicionado no autoriza ignorar futuras revisiones de la base ni cambiar pins localmente.

## 4. Exact file manifest

```text
CREATE tessera_posix/source-lock.json
CREATE tessera_posix/deployment-profile.template.json
CREATE tessera_posix/build_and_verify_official.sh
CREATE tessera_posix/run_verified_log.sh
CREATE tessera_posix/test_contracts.py
CREATE tessera_posix/README.md
CREATE tessera_posix/upstream/LICENSE
CREATE tessera_posix/upstream/AUTHORS
CREATE tessera_posix/upstream/POSIX_README.md
CREATE tessera_posix/upstream/posix-oneshot/main.go
CREATE tessera_posix/upstream/fsck/main.go
```

## 5. Materialization blocks

### FILE: `tessera_posix/source-lock.json`
```yaml
block_id: "TESSERA-POSIX:source-lock-json:v1"
operation: CREATE
provenance: AUTHORED
source: "local exact identity and evidence lock around official Tessera"
license: "LicenseRef-Workspace-Owner"
sha256: "4c09d8600ff78f32e9cbb6e90a78b363973f2c1b820e87181bf65578ce436a37"
variables: []
secrets_allowed: false
```
````text
{
  "schema": "elite-tessera-posix-source-lock/v1",
  "verified_at": "2026-08-27",
  "source": {
    "repository": "transparency-dev/tessera",
    "commit": "a8f33c56b80be808712c0537cb1b71f2fd5846ea",
    "commit_signature": "verified",
    "commit_signature_reason": "valid",
    "commit_time": "2026-08-26T18:05:57Z",
    "tree": "71d51fdd8d119e5dd31394d99859c073c9544881",
    "archive_url": "https://github.com/transparency-dev/tessera/archive/a8f33c56b80be808712c0537cb1b71f2fd5846ea.zip",
    "archive_bytes": 2215098,
    "archive_sha256": "33bb2394a0d6ce0ed34b3dadbf6fdac6059afada6611f694c88be1b5fbdd779c",
    "module_version": "v1.0.3-0.20260826180557-a8f33c56b80b",
    "module_sum": "h1:IqUpmTdA0rRtIPRszAnMwNl+gyCZXZz0WCrqkplQAHQ=",
    "module_go_mod_sum": "h1:cH4u91YS85J/lmRxL/zcLEKYNg9lzLvYHQ0qq5OURN0=",
    "go_mod_sha256": "c4a71771d54d1724e7729e22135cccdc23f4689eabd6854168402ad2c1f783c8",
    "go_sum_sha256": "70305ce806819ad63d4ae1c3ed4cf7f958e12fe83da9009b925dcd0875ba6cc2",
    "license": "Apache-2.0",
    "license_sha256": "c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4",
    "authors_sha256": "7ae698cb6e29d0df7d8f485d60d8993aa1de88e04e5b85b86afda76b415c8d58",
    "security_policy_sha256": "8b8bd2544980bfe2a1b76d7a8a658ece6bd5e5d071d536c899523e89cddf3bff",
    "posix_cli_sha256": "8f8e53c6c26b917d49fe131c5baa27b77588e40ebf1c0008116015d56ed63c77",
    "fsck_cli_sha256": "48ed7524f671c8b0f62746a367ee637c321047d5ea0a4d8112ab20b9a36fa91a",
    "posix_design_sha256": "c50218a0df7de7c5e92cbd373a6721bc905967c950aeb67e0f2c98a1d84358e8"
  },
  "target": {
    "os": "linux",
    "architecture": "amd64",
    "go": "1.26.7",
    "govulncheck": "1.7.0",
    "govulncheck_binary_sha256": "8f7bc456defc10509f3821cd2d9271e730fe5f7a56251a95f22b2e0c43b9db80",
    "posix_oneshot_binary_sha256": "d36f44eaca595fb4109b1c211c37d82f46b83326bdb91e4d796d682149cc0c15",
    "fsck_binary_sha256": "5e75e44a93bc5ad0d56c80f2f95f4655fc19ca61f4fed383c005e036563cf591"
  },
  "evidence": {
    "source_suite": "PASS",
    "go_vet": "PASS",
    "fault_injection": "PASS",
    "end_to_end_entries": 10,
    "end_to_end_root_hex": "421c41a01e8d1f4d7b70f5fb57e358f37b098d7e60af6071b1462a3d8947b891",
    "tampered_bundle_rejected": true,
    "module_count": 321,
    "osv_records": 5,
    "govulncheck_reachable": 0,
    "govulncheck_imported_packages": 0,
    "govulncheck_affected_modules_not_called": 3
  },
  "conditions": [
    "This exact signed commit is newer than the latest stable v1.0.4 release and is not itself a release.",
    "The source graph retains five OSV records in three modules that govulncheck reports as not called by the POSIX CLI.",
    "The POSIX driver is not compatible with Windows or NTFS and must run on a filesystem with the semantics documented by upstream.",
    "A local log is tamper-evident, not deletion-proof; independent checkpoint publication, witnessing, backup and access control remain project gates.",
    "The log may contain only approved canonical evidence receipts, never document bytes, extracted values, credentials or personal data."
  ]
}
````

### FILE: `tessera_posix/deployment-profile.template.json`
```yaml
block_id: "TESSERA-POSIX:deployment-profile-template-json:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed configuration surface"
license: "LicenseRef-Workspace-Owner"
sha256: "b5011eef89cbe6993a0621769e683aa8510e42ce93d98f3b7467d565b5c01a3b"
variables: []
secrets_allowed: false
```
````text
{
  "schema": "elite-tessera-posix-deployment-profile/v1",
  "enabled": false,
  "production_authorized": false,
  "source_receipt_path": "",
  "linux_distribution": "",
  "filesystem_type": "",
  "storage_directory": "",
  "private_key_file": "",
  "public_key_file": "",
  "approved_receipt_glob": "",
  "checkpoint_publication_target": "",
  "witness_policy_file": "",
  "backup_target": "",
  "retention_policy_reference": "",
  "incident_owner": "",
  "recovery_owner": ""
}
````

### FILE: `tessera_posix/build_and_verify_official.sh`
```yaml
block_id: "TESSERA-POSIX:build-and-verify-official-sh:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed build gate invoking exact upstream source"
license: "LicenseRef-Workspace-Owner"
sha256: "748766368791a9a34b0254fe5c869cbec373a8855d9f92340b10d54eac2584a3"
variables: []
secrets_allowed: false
```
````text
#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 3 ]]; then
  echo "usage: $0 <exact-tessera-source-root> <empty-output-directory> <govulncheck-binary>" >&2
  exit 64
fi

source_root=$(readlink -f -- "$1")
output_root=$(readlink -m -- "$2")
govulncheck_bin=$(readlink -f -- "$3")

[[ $(uname -s) == Linux ]] || { echo "Linux is required" >&2; exit 65; }
[[ $(uname -m) == x86_64 ]] || { echo "linux/amd64 is required" >&2; exit 65; }
[[ -d "$source_root" ]] || { echo "source root is missing" >&2; exit 66; }
[[ -x "$govulncheck_bin" ]] || { echo "govulncheck is missing" >&2; exit 66; }
if [[ -e "$output_root" ]] && [[ -n $(find "$output_root" -mindepth 1 -maxdepth 1 -print -quit) ]]; then
  echo "output directory must be absent or empty" >&2
  exit 73
fi
mkdir -p -- "$output_root"

source_fs=$(stat -f -c %T -- "$source_root")
case "$source_fs" in
  ext2/ext3|xfs|zfs|ceph|tmpfs|overlayfs) ;;
  *) echo "source must be on an admitted Linux filesystem, got $source_fs" >&2; exit 65 ;;
esac

[[ $(go version) == "go version go1.26.7 linux/amd64" ]] || {
  echo "exact Go 1.26.7 linux/amd64 is required" >&2
  exit 69
}
"$govulncheck_bin" -version | grep -F 'govulncheck@v1.7.0' >/dev/null || {
  echo "exact govulncheck 1.7.0 is required" >&2
  exit 69
}

cd "$source_root"
sha256sum -c <<'HASHES'
c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4  LICENSE
7ae698cb6e29d0df7d8f485d60d8993aa1de88e04e5b85b86afda76b415c8d58  AUTHORS
8b8bd2544980bfe2a1b76d7a8a658ece6bd5e5d071d536c899523e89cddf3bff  SECURITY.md
c4a71771d54d1724e7729e22135cccdc23f4689eabd6854168402ad2c1f783c8  go.mod
70305ce806819ad63d4ae1c3ed4cf7f958e12fe83da9009b925dcd0875ba6cc2  go.sum
8f8e53c6c26b917d49fe131c5baa27b77588e40ebf1c0008116015d56ed63c77  cmd/examples/posix-oneshot/main.go
48ed7524f671c8b0f62746a367ee637c321047d5ea0a4d8112ab20b9a36fa91a  cmd/fsck/main.go
c50218a0df7de7c5e92cbd373a6721bc905967c950aeb67e0f2c98a1d84358e8  storage/posix/README.md
HASHES

go mod download
go mod verify
go test ./...
go vet ./...
go build -trimpath -o "$output_root/posix-oneshot" ./cmd/examples/posix-oneshot
go build -trimpath -o "$output_root/fsck" ./cmd/fsck

echo 'd36f44eaca595fb4109b1c211c37d82f46b83326bdb91e4d796d682149cc0c15  posix-oneshot' > "$output_root/SHA256SUMS"
echo '5e75e44a93bc5ad0d56c80f2f95f4655fc19ca61f4fed383c005e036563cf591  fsck' >> "$output_root/SHA256SUMS"
(cd "$output_root" && sha256sum -c SHA256SUMS)

"$govulncheck_bin" ./cmd/examples/posix-oneshot > "$output_root/govulncheck-source.txt"
"$govulncheck_bin" -mode=binary "$output_root/posix-oneshot" > "$output_root/govulncheck-binary.txt"
grep -F 'Your code is affected by 0 vulnerabilities.' "$output_root/govulncheck-source.txt" >/dev/null
grep -F '0 vulnerabilities in packages you import' "$output_root/govulncheck-source.txt" >/dev/null
grep -F 'Your code is affected by 0 vulnerabilities.' "$output_root/govulncheck-binary.txt" >/dev/null
grep -F '0 vulnerabilities in packages you import' "$output_root/govulncheck-binary.txt" >/dev/null

echo "TESSERA_OFFICIAL_BUILD_PASS source=$source_root output=$output_root"
````

### FILE: `tessera_posix/run_verified_log.sh`
```yaml
block_id: "TESSERA-POSIX:run-verified-log-sh:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed runtime gate around official Tessera CLIs"
license: "LicenseRef-Workspace-Owner"
sha256: "2bcf94a45d45d67a2e0e2323362aae19b8e18c1ce309ff71eaf2e9b248c2e33b"
variables: []
secrets_allowed: false
```
````text
#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 6 ]]; then
  echo "usage: $0 <create|append|verify> <tools-directory> <storage-directory> <private-key-file-or-dash> <public-key-file> <approved-receipt-glob-or-dash>" >&2
  exit 64
fi

mode=$1
tools_root=$(readlink -f -- "$2")
storage_root=$(readlink -m -- "$3")
private_key_arg=$4
public_key_file=$(readlink -f -- "$5")
receipt_glob=$6

[[ $(uname -s) == Linux ]] || { echo "Linux is required" >&2; exit 65; }
[[ -x "$tools_root/posix-oneshot" && -x "$tools_root/fsck" ]] || { echo "verified tools are missing" >&2; exit 66; }
[[ -f "$tools_root/SHA256SUMS" ]] || { echo "tool checksum manifest is missing" >&2; exit 66; }
(cd "$tools_root" && sha256sum -c SHA256SUMS)
[[ -f "$public_key_file" && ! -L "$public_key_file" ]] || { echo "public key file is invalid" >&2; exit 66; }

case "$mode" in
  create)
    [[ ! -e "$storage_root" ]] || { echo "create requires an absent storage directory" >&2; exit 73; }
    ;;
  append|verify)
    [[ -f "$storage_root/checkpoint" ]] || { echo "$mode requires an existing checkpoint" >&2; exit 66; }
    ;;
  *) echo "unsupported mode" >&2; exit 64 ;;
esac

if [[ "$mode" != verify ]]; then
  [[ "$private_key_arg" != '-' ]] || { echo "write modes require a private key file" >&2; exit 64; }
  private_key_file=$(readlink -f -- "$private_key_arg")
  [[ -f "$private_key_file" && ! -L "$private_key_file" ]] || { echo "private key file is invalid" >&2; exit 66; }
  private_mode=$(stat -c %a -- "$private_key_file")
  [[ "$private_mode" == 600 || "$private_mode" == 400 ]] || { echo "private key permissions must be 600 or 400" >&2; exit 77; }
  [[ "$receipt_glob" != '-' && "$receipt_glob" == *.json* ]] || { echo "only an approved JSON receipt glob is allowed" >&2; exit 64; }
  mapfile -t receipt_files < <(compgen -G "$receipt_glob" || true)
  [[ ${#receipt_files[@]} -gt 0 ]] || { echo "receipt glob matched no files" >&2; exit 66; }
  for receipt in "${receipt_files[@]}"; do
    [[ -f "$receipt" && ! -L "$receipt" ]] || { echo "receipt must be a regular non-symlink file: $receipt" >&2; exit 66; }
    [[ $(stat -c %s -- "$receipt") -le 65535 ]] || { echo "receipt exceeds Tessera entry limit: $receipt" >&2; exit 65; }
    grep -F '"document_sha256"' "$receipt" >/dev/null || { echo "receipt lacks document_sha256: $receipt" >&2; exit 65; }
    grep -F '"evidence_sha256"' "$receipt" >/dev/null || { echo "receipt lacks evidence_sha256: $receipt" >&2; exit 65; }
  done
  "$tools_root/posix-oneshot" --storage_dir="$storage_root" --entries="$receipt_glob" --private_key="$private_key_file"
fi

storage_fs=$(stat -f -c %T -- "$storage_root")
case "$storage_fs" in
  ext2/ext3|xfs|zfs|ceph|tmpfs|overlayfs) ;;
  *) echo "storage filesystem is not admitted: $storage_fs" >&2; exit 65 ;;
esac

"$tools_root/fsck" --storage_url="file://$storage_root/" --public_key="$public_key_file" --ui=false
checkpoint_sha256=$(sha256sum "$storage_root/checkpoint" | cut -d ' ' -f 1)
echo "TESSERA_LOG_VERIFIED mode=$mode checkpoint_sha256=$checkpoint_sha256 storage=$storage_root"
````

### FILE: `tessera_posix/test_contracts.py`
```yaml
block_id: "TESSERA-POSIX:test-contracts-py:v1"
operation: CREATE
provenance: AUTHORED
source: "local static regression contracts"
license: "LicenseRef-Workspace-Owner"
sha256: "073cd6364a421fb6354a9f2f97546b1ccfeb16b99822a957dfbe2859bf183065"
variables: []
secrets_allowed: false
```
````text
from __future__ import annotations

import hashlib
import json
import pathlib
import unittest


ROOT = pathlib.Path(__file__).resolve().parent


def sha256(path: pathlib.Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


class TesseraPackContracts(unittest.TestCase):
    def test_template_enables_nothing(self) -> None:
        profile = json.loads((ROOT / "deployment-profile.template.json").read_text(encoding="utf-8"))
        self.assertIs(profile["enabled"], False)
        self.assertIs(profile["production_authorized"], False)
        for key, value in profile.items():
            if key not in {"schema", "enabled", "production_authorized"}:
                self.assertEqual(value, "")

    def test_official_bytes_match_lock(self) -> None:
        lock = json.loads((ROOT / "source-lock.json").read_text(encoding="utf-8"))["source"]
        self.assertEqual(sha256(ROOT / "upstream" / "LICENSE"), lock["license_sha256"])
        self.assertEqual(sha256(ROOT / "upstream" / "AUTHORS"), lock["authors_sha256"])
        self.assertEqual(sha256(ROOT / "upstream" / "POSIX_README.md"), lock["posix_design_sha256"])
        self.assertEqual(sha256(ROOT / "upstream" / "posix-oneshot" / "main.go"), lock["posix_cli_sha256"])
        self.assertEqual(sha256(ROOT / "upstream" / "fsck" / "main.go"), lock["fsck_cli_sha256"])

    def test_identity_and_conditions_are_fail_closed(self) -> None:
        lock = json.loads((ROOT / "source-lock.json").read_text(encoding="utf-8"))
        self.assertEqual(lock["source"]["commit"], "a8f33c56b80be808712c0537cb1b71f2fd5846ea")
        self.assertEqual(lock["source"]["commit_signature"], "verified")
        self.assertEqual(lock["evidence"]["govulncheck_reachable"], 0)
        self.assertEqual(lock["evidence"]["govulncheck_imported_packages"], 0)
        self.assertGreater(lock["evidence"]["osv_records"], 0)
        self.assertTrue(any("not itself a release" in item for item in lock["conditions"]))
        self.assertTrue(any("never document bytes" in item for item in lock["conditions"]))

    def test_scripts_preserve_required_gates(self) -> None:
        build = (ROOT / "build_and_verify_official.sh").read_text(encoding="utf-8")
        run = (ROOT / "run_verified_log.sh").read_text(encoding="utf-8")
        for token in ("go mod verify", "go test ./...", "go vet ./...", "govulncheck-source.txt", "sha256sum -c"):
            self.assertIn(token, build)
        for token in ("create requires an absent storage directory", "private key permissions", "document_sha256", "evidence_sha256", "fsck"):
            self.assertIn(token, run)


if __name__ == "__main__":
    unittest.main()
````

### FILE: `tessera_posix/README.md`
```yaml
block_id: "TESSERA-POSIX:readme-md:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration and admission runbook"
license: "LicenseRef-Workspace-Owner"
sha256: "5ac3b72f7a9c58d6d28a004a2dab21565f1545f233dabace60485c24a4d107c1"
variables: []
secrets_allowed: false
```
````text
# Tessera POSIX transparency log lane

This lane materializes exact public code from `transparency-dev/tessera` at signed commit `a8f33c56b80be808712c0537cb1b71f2fd5846ea`. The upstream `AUTHORS` file names Google LLC, Anthropic PBC and Internet Security Research Group. Google Trillian recommends Tessera for new transparency logs; this pack does not relabel the repository as a Google product.

The packaged `posix-oneshot`, `fsck`, POSIX design, license and authors files are byte-for-byte upstream evidence. `build_and_verify_official.sh`, `run_verified_log.sh`, the profile, lock and tests are workspace-authored integration glue and are never vendor code.

## What it proves

On Linux/amd64 with Go 1.26.7, the exact signed commit passed `go mod verify`, its complete Go suite including POSIX I/O and SIGKILL fault injection, `go vet`, deterministic builds, an append of ten entries, cryptographic `fsck`, and rejection of a truncated entry bundle. `govulncheck 1.7.0` found zero reachable or imported-package vulnerabilities in the POSIX CLI. The graph still contains five OSV records in three required but uncalled modules; a new upstream revision must repeat every gate.

This is tamper-evident storage, not deletion-proof storage. Production still requires independent checkpoint publication or witnesses, access control, backup, restore, retention and incident ownership. Log only canonical evidence receipts containing hashes and opaque references. Never log document bytes, extracted field values, personal data, secrets or credentials.

## Apply

1. Acquire source ID `transparency-dev-tessera-a8f33c5` through the official acquisition pack after the user completes its Linux, filesystem, checkpoint, witness, backup and retention inputs.
2. Use a Linux filesystem with upstream-supported POSIX semantics. Windows/NTFS and WSL-mounted Windows paths are not admitted storage targets.
3. Install or acquire exact Go 1.26.7 and official `govulncheck` 1.7.0. Run:

   `./build_and_verify_official.sh <source-root> <empty-tools-dir> <govulncheck>`

4. Complete `deployment-profile.template.json`; do not change `enabled` or `production_authorized` until the project supplies all paths and gates.
5. For an absent log directory, run `run_verified_log.sh create ...`; for an existing verified log, use `append`; use `verify` for read-only integrity checking. The write modes require a private-key file with mode 600 or 400 and JSON evidence receipts with `document_sha256` plus `evidence_sha256`.

The scripts never create keys, choose retention, authorize production, publish checkpoints, configure a witness or repair a failed log automatically. Those decisions remain explicit project configuration and evidence.
````

### FILE: `tessera_posix/upstream/LICENSE`
```yaml
block_id: "TESSERA-POSIX:upstream-license:v1"
operation: CREATE
provenance: VERBATIM
source: "https://github.com/transparency-dev/tessera/blob/a8f33c56b80be808712c0537cb1b71f2fd5846ea/LICENSE"
license: "Apache-2.0"
sha256: "c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4"
variables: []
secrets_allowed: false
```
````text
                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "[]"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright [yyyy] [name of copyright owner]

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
````

### FILE: `tessera_posix/upstream/AUTHORS`
```yaml
block_id: "TESSERA-POSIX:upstream-authors:v1"
operation: CREATE
provenance: VERBATIM
source: "https://github.com/transparency-dev/tessera/blob/a8f33c56b80be808712c0537cb1b71f2fd5846ea/AUTHORS"
license: "Apache-2.0"
sha256: "7ae698cb6e29d0df7d8f485d60d8993aa1de88e04e5b85b86afda76b415c8d58"
variables: []
secrets_allowed: false
```
````text
# This is the official list of benchmark authors for copyright purposes.
# This file is distinct from the CONTRIBUTORS files.
# See the latter for an explanation.
#
# Names should be added to this file as:
#       Name or Organization <email address>
# The email address is not required for organizations.
#
# Please keep the list sorted.

Anthropic PBC
Google LLC
Internet Security Research Group
````

### FILE: `tessera_posix/upstream/POSIX_README.md`
```yaml
block_id: "TESSERA-POSIX:upstream-posix-readme:v1"
operation: CREATE
provenance: VERBATIM
source: "https://github.com/transparency-dev/tessera/blob/a8f33c56b80be808712c0537cb1b71f2fd5846ea/storage/posix/README.md"
license: "Apache-2.0"
sha256: "c50218a0df7de7c5e92cbd373a6721bc905967c950aeb67e0f2c98a1d84358e8"
variables: []
secrets_allowed: false
```
````text
# Tessera on POSIX Filesystems

This document describes the storage implementation for running Tessera on a POSIX-compliant filesystem.

## Overview

POSIX provides for a small number of atomic operations on compliant filesystems.

This design leverages those to safely maintain a Merkle tree log on disk, in a format
which can be exposed directly via a read-only endpoint to clients of the log (for example,
using `nginx` or similar).

In contrast with some of the other storage backends, sequencing and integration of entries into
the tree is synchronous.

The implementation uses a `.state/` directory to coordinate operation.
This directory does _not_ need to be visible to log clients, but it does not contain sensitive
data and so it isn't a problem if it is made visible.

## Life of a Leaf

In the description below, when we talk about writing to files - either appending or creating new ones,
the _actual_ process used always follows the following pattern:
1. Create a temporary file on the same filesystem as the target location
1. If we're appending data, copy the contents of the prefix location into the temporary file
1. Write any new/additional data into the temporary file
1. Close the temporary file
1. Rename the temporary file into the target location.

The final step in the dance above is atomic according to the POSIX spec, so in performing this sequence
of actions we can avoid corrupt or partially written files being part of the tree.

1. Leaves are submitted by the binary built using Tessera via a call the storage's `Add` func.
1. The storage library batches these entries up in memory, and, after a configurable period of time has elapsed
   or the batch reaches a configurable size threshold, the batch is sequenced and appended to the tree:
   1. An advisory lock is taken on `.state/treeState.lock` file.
   1. Flushed entries are assigned contiguous sequence numbers, and written out into entry bundle files.
   1. Integrate newly added leaves into Merkle tree, and write tiles out as files.
   1. Update `./state/treeState` file with the new size & root hash.
1. Asynchronously, at an interval determined by the `WithCheckpointInterval` option, the `checkpoint` file
will be updated:
   1. An advisory lock is taken on `.state/publish.lock`
   1. If the last-modified date of the `checkpoint` file is older than the checkpoint update interval,
      a new checkpoint which commits to the latest tree state is produced and written to the `checkpoint`
      file.

## Filesystems

This implementation has been somewhat tested on local `ext4` and `ZFS` filesystems, and on a distributed
[CephFS](https://docs.ceph.com/en/reef/cephfs/) instance on GCP, in all cases with multiple
personality binaries attempting to add new entries concurrently.

Other POSIX compliant filesystems such as `XFS` _should_ work, but filesystems which do not offer strong
POSIX compliance (e.g. `s3fs` or `NFS`) are unlikely to result in long term happiness.

If in doubt, tools like https://github.com/saidsay-so/pjdfstest may help in determining whether a given
filesystem is suitable.

## Memory Management (GOMEMLIMIT)

The POSIX persistent antispam feature is powered by BadgerDB, which can be memory-intensive under heavy write loads.
This is particularly acute during Badger's periodic compaction windows.

To prevent Out-Of-Memory (OOM) crashes, it is **strongly** recommended to configure the `GOMEMLIMIT` environment variable (e.g., `GOMEMLIMIT=2GiB`). This restricts the Go runtime's memory usage aggressively and helps ensure process stability.

````

### FILE: `tessera_posix/upstream/posix-oneshot/main.go`
```yaml
block_id: "TESSERA-POSIX:upstream-posix-oneshot-main-go:v1"
operation: CREATE
provenance: VERBATIM
source: "https://github.com/transparency-dev/tessera/blob/a8f33c56b80be808712c0537cb1b71f2fd5846ea/cmd/examples/posix-oneshot/main.go"
license: "Apache-2.0"
sha256: "8f8e53c6c26b917d49fe131c5baa27b77588e40ebf1c0008116015d56ed63c77"
variables: []
secrets_allowed: false
```
````text
// Copyright 2024 The Tessera authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// posix-oneshot is a command line tool for adding entries to a local
// tlog-tiles log stored on a posix filesystem.
// The command takes a list of new entries to add to the log, and exits
// when they are successfully integrated.
// See the README in this package for more detailed usage instructions.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/mod/sumdb/note"

	"log/slog"

	"github.com/transparency-dev/tessera"
	"github.com/transparency-dev/tessera/storage/posix"
)

var (
	storageDir        = flag.String("storage_dir", "", "Root directory to store log data.")
	entries           = flag.String("entries", "", "File path glob of entries to add to the log.")
	privKeyFile       = flag.String("private_key", "", "Location of private key file. If unset, uses the contents of the LOG_PRIVATE_KEY environment variable.")
	witnessPolicyFile = flag.String("witness_policy_file", "", "(Optional) Path to the file containing the witness policy in the format describe at https://git.glasklar.is/sigsum/core/sigsum-go/-/blob/main/doc/policy.md")
	witnessTimeout    = flag.Duration("witness_timeout", tessera.DefaultWitnessTimeout, "Maximum time to wait for witness responses.")
	witnessFailOpen   = flag.Bool("witness_fail_open", false, "Still publish a checkpoint even if witness policy could not be met")
	slogLevel         = flag.Int("slog_level", 0, "The cut-off threshold for structured logging. Default is 0 (INFO). See https://pkg.go.dev/log/slog#Level for other levels.")
)

// entryInfo binds the actual bytes to be added as a leaf with a
// user-recognisable name for the source of those bytes.
// The name is only used below in order to inform the user of the
// sequence numbers assigned to the data from the provided input files.
type entryInfo struct {
	name string
	f    tessera.IndexFuture
}

func main() {
	flag.Parse()
	ctx := context.Background()
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.Level(*slogLevel)})))

	slog.DebugContext(ctx, "Initialising driver")

	// Gather the info needed for reading/writing checkpoints
	s := getSignerOrDie()
	// Construct a new Tessera POSIX log storage, anchored at the correct directory, and initialising it if requested.
	// The options provide the checkpoint signer & verifier, and batch options.
	// In this case, we want to create a single batch containing all of the leaves being added in order to
	// add all of these leaves without creating any intermediate checkpoints.
	driver, err := posix.New(
		ctx,
		posix.Config{
			Path: *storageDir,
		},
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to construct storage", slog.Any("error", err))
		os.Exit(1)
	}

	slog.DebugContext(ctx, "Reading entries")
	// Evaluate the glob provided by the --entries flag to determine the files containing leaves
	filesToAdd := readEntriesOrDie()
	batchSize := uint(len(filesToAdd))
	if batchSize == 0 {
		// batchSize can't be zero
		batchSize = 1
	}

	slog.DebugContext(ctx, "Configuring options")
	opts := tessera.NewAppendOptions().
		WithCheckpointSigner(s).
		// Hint to Tessera the number of entries we're about to add via the batchSize parameter below,
		// this will cause the batch to flush as soon as we've called Add on the final entry.
		WithBatching(batchSize, 100*time.Millisecond).
		// We're unlikely to ever wait this long to publish a checkpoint because of the batchSize hint
		// passed in to the option above, but we set this interval low primarily such that if the user re-runs this
		// tool to add further entries to the log, they don't have to wait for the previous checkpoint to
		// become old enough to be overwritten.
		WithCheckpointInterval(100 * time.Millisecond)

	if *witnessPolicyFile != "" {
		f, err := os.ReadFile(*witnessPolicyFile)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to read witness policy file", slog.String("witnesspolicyfile", *witnessPolicyFile), slog.Any("error", err))
			os.Exit(1)
		}
		wg, err := tessera.NewWitnessGroupFromPolicy(f)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to create witness group from policy", slog.Any("error", err))
			os.Exit(1)
		}

		wOpts := &tessera.WitnessOptions{
			FailOpen: *witnessFailOpen,
			Timeout:  *witnessTimeout,
		}
		opts.WithWitnesses(wg, wOpts)
	}

	slog.DebugContext(ctx, "Creating appender")
	appender, shutdown, r, err := tessera.NewAppender(ctx, driver, opts)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create new appender", slog.Any("error", err))
		os.Exit(1)
	}

	slog.DebugContext(ctx, "Creating awaiter")
	// We don't want to exit until our entries have been integrated into the tree, so we'll use Tessera's
	// PublicationAwaiter to help with that.
	await := tessera.NewPublicationAwaiter(ctx, r.ReadCheckpoint, 100*time.Millisecond)

	slog.DebugContext(ctx, "Adding entries")
	// Add each of the leaves in order, and store the futures in a slice
	// that we will check once all leaves are sent to storage.
	indexFutures := make([]entryInfo, 0, len(filesToAdd))
	for _, fp := range filesToAdd {
		b, err := os.ReadFile(fp)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to read entry file", slog.String("fp", fp), slog.Any("error", err))
			os.Exit(1)
		}

		f := appender.Add(ctx, tessera.NewEntry(b))
		indexFutures = append(indexFutures, entryInfo{name: fp, f: f})
	}

	slog.DebugContext(ctx, "Awaiting entries")
	// Two options to ensure all work is done:
	// 1) Check each of the futures to ensure that the leaves are sequenced.
	for _, entry := range indexFutures {
		seq, _, err := await.Await(ctx, entry.f)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to sequence", slog.String("name", entry.name), slog.Any("error", err))
			os.Exit(1)
		}
		slog.InfoContext(ctx, "Integrated entry", slog.Uint64("index", seq.Index), slog.String("name", entry.name))
	}
	slog.DebugContext(ctx, "Futures resolved")
	slog.DebugContext(ctx, "Shutting down")

	// 2) shutdown the appender
	if err := shutdown(ctx); err != nil {
		slog.ErrorContext(ctx, "Failed to shut down cleanly", slog.Any("error", err))
		os.Exit(1)
	}
	slog.DebugContext(ctx, "Finished")
}

// Read log private key from file or environment variable
func getSignerOrDie() note.Signer {
	var privKey string
	var err error
	if len(*privKeyFile) > 0 {
		privKey, err = getKeyFile(*privKeyFile)
		if err != nil {
			slog.ErrorContext(context.Background(), "Unable to get private key", slog.Any("error", err))
			os.Exit(1)
		}
	} else {
		privKey = os.Getenv("LOG_PRIVATE_KEY")
		if len(privKey) == 0 {
			slog.ErrorContext(context.Background(), "Supply private key file path using --private_key or set LOG_PRIVATE_KEY environment variable")
			os.Exit(1)
		}
	}
	s, err := note.NewSigner(privKey)
	if err != nil {
		slog.ErrorContext(context.Background(), "Failed to instantiate signer", slog.Any("error", err))
		os.Exit(1)
	}
	return s
}

func getKeyFile(path string) (string, error) {
	k, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read key file: %w", err)
	}
	return string(k), nil
}

func readEntriesOrDie() []string {
	toAdd, err := filepath.Glob(*entries)
	if err != nil {
		slog.ErrorContext(context.Background(), "Failed to glob entries", slog.String("entries", *entries), slog.Any("error", err))
		os.Exit(1)
	}
	slog.DebugContext(context.Background(), "toAdd", slog.Any("files", toAdd))
	return toAdd
}
````

### FILE: `tessera_posix/upstream/fsck/main.go`
```yaml
block_id: "TESSERA-POSIX:upstream-fsck-main-go:v1"
operation: CREATE
provenance: VERBATIM
source: "https://github.com/transparency-dev/tessera/blob/a8f33c56b80be808712c0537cb1b71f2fd5846ea/cmd/fsck/main.go"
license: "Apache-2.0"
sha256: "48ed7524f671c8b0f62746a367ee637c321047d5ea0a4d8112ab20b9a36fa91a"
variables: []
secrets_allowed: false
```
````text
// Copyright 2025 The Tessera authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// fsck is a command-line tool for checking the integrity of a tlog-tiles based log.
package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"time"

	"log/slog"

	f_note "github.com/transparency-dev/formats/note"
	"github.com/transparency-dev/merkle/rfc6962"
	"github.com/transparency-dev/tessera/api"
	"github.com/transparency-dev/tessera/client"
	"github.com/transparency-dev/tessera/cmd/fsck/internal/tui"
	"github.com/transparency-dev/tessera/fsck"
	"golang.org/x/mod/sumdb/note"
	"golang.org/x/time/rate"
)

var (
	storageURL  = flag.String("storage_url", "", "Base tlog-tiles URL")
	bearerToken = flag.String("bearer_token", "", "The bearer token for authorizing HTTP requests to the storage URL, if needed")
	N           = flag.Uint("N", 1, "The number of workers to use when fetching/comparing resources")
	origin      = flag.String("origin", "", "Origin of the log to check, if unset, will use the name of the provided public key")
	pubKey      = flag.String("public_key", "", "Path to a file containing the log's public key")
	qps         = flag.Float64("qps", 0, "Max QPS to send to the target log. Set to zero for unlimited")
	ui          = flag.Bool("ui", true, "Set to true to use a TUI to display progress, or false for logging")
	slogLevel   = flag.Int("slog_level", 0, "The cut-off threshold for structured logging. Default is 0 (INFO). See https://pkg.go.dev/log/slog#Level for other levels.")
)

func main() {
	flag.Parse()
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.Level(*slogLevel)})))
	ctx, cancel := context.WithCancel(context.Background())
	logURL, err := url.Parse(*storageURL)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid --storage_url", slog.String("param", *storageURL), slog.Any("error", err))
		os.Exit(1)
	}
	var src fsck.Fetcher

	if logURL.Scheme == "file" {
		src = &client.FileFetcher{
			Root: logURL.Path,
		}
	} else {
		httpSrc, err := client.NewHTTPFetcher(logURL, nil)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to create HTTP fetcher", slog.Any("error", err))
			os.Exit(1)
		}
		if *bearerToken != "" {
			httpSrc.SetAuthorizationHeader(fmt.Sprintf("Bearer %s", *bearerToken))
		}
		src = httpSrc
	}
	if *qps > 0 {
		src = &rateLimitedSrc{
			rl:       rate.NewLimiter(rate.Limit(*qps), 10),
			delegate: src,
		}
	}
	v := verifierFromFlags(ctx)
	if *origin == "" {
		*origin = v.Name()
	}
	f := fsck.New(*origin, v, src, defaultMerkleLeafHasher, fsck.Opts{N: *N})

	// checkResult is a channel to receive the result of the fsck check.
	// It is used to signal the main goroutine when the check is complete.
	checkResult := make(chan error, 1)
	go func() {
		err := f.Check(ctx)

		// Only log the error if we're using the TUI, otherwise it will be double logged in headless mode.
		if err != nil && *ui {
			slog.ErrorContext(ctx, "fsck failed", slog.Any("error", err))
		}
		slog.DebugContext(ctx, "Completed ranges", slog.Any("status", f.Status()))

		checkResult <- err
		close(checkResult)
		cancel() // We're done, so we can exit.
	}()

	if *ui {
		if err := tui.RunApp(ctx, f); err != nil {
			slog.ErrorContext(ctx, "App exited", slog.Any("error", err))
		}
		// User may have exited the UI, cancel the context to signal to the async fsck process too.
		cancel()
	} else {
		for {
			select {
			case err := <-checkResult:
				if err != nil {
					slog.ErrorContext(ctx, "fsck failed", slog.Any("error", err))
					os.Exit(1)
				}
				return
			case <-time.After(time.Second):
				slog.DebugContext(ctx, "Ranges", slog.Any("status", f.Status()))
			}
		}
	}
}

// defaultMerkleLeafHasher parses a C2SP tlog-tile bundle and returns the Merkle leaf hashes of each entry it contains.
func defaultMerkleLeafHasher(bundle []byte) ([][]byte, error) {
	eb := &api.EntryBundle{}
	if err := eb.UnmarshalText(bundle); err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}
	r := make([][]byte, 0, len(eb.Entries))
	for _, e := range eb.Entries {
		h := rfc6962.DefaultHasher.HashLeaf(e)
		r = append(r, h[:])
	}
	return r, nil
}

func verifierFromFlags(ctx context.Context) note.Verifier {
	if *pubKey == "" {
		slog.ErrorContext(ctx, "Must provide the --public_key flag")
		os.Exit(1)
	}
	b, err := os.ReadFile(*pubKey)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to read verifier from", slog.String("pubkey", *pubKey), slog.Any("error", err))
		os.Exit(1)
	}
	v, err := f_note.NewVerifier(string(b))
	if err != nil {
		slog.ErrorContext(ctx, "Invalid verifier in", slog.String("pubkey", *pubKey), slog.Any("error", err))
		os.Exit(1)
	}
	return v
}

type rateLimitedSrc struct {
	rl       *rate.Limiter
	delegate fsck.Fetcher
}

func (r *rateLimitedSrc) ReadCheckpoint(ctx context.Context) ([]byte, error) {
	if err := r.rl.Wait(ctx); err != nil {
		return nil, err
	}
	return r.delegate.ReadCheckpoint(ctx)
}

func (r *rateLimitedSrc) ReadTile(ctx context.Context, l, i uint64, p uint8) ([]byte, error) {
	if err := r.rl.Wait(ctx); err != nil {
		return nil, err
	}
	return r.delegate.ReadTile(ctx, l, i, p)
}

func (r *rateLimitedSrc) ReadEntryBundle(ctx context.Context, i uint64, p uint8) ([]byte, error) {
	if err := r.rl.Wait(ctx); err != nil {
		return nil, err
	}
	return r.delegate.ReadEntryBundle(ctx, i, p)
}
````

## 6. Configuration surface

Obligatorio: Linux distribution, exact source receipt, filesystem type, log path, private/public key paths, approved receipt glob, checkpoint publication or witness policy, backup, retention, incident owner and recovery owner. The distributed profile has `enabled=false` and `production_authorized=false`; every value is empty.

The private key never enters Markdown, profile or log history. Write modes require an external regular file with mode 600 or 400. Receipts must be regular non-symlink JSON files no larger than the upstream limit and include `document_sha256` plus `evidence_sha256`.

## 7. Dependency bill

- Go 1.26.7 linux/amd64, acquired from the official Go distribution and verified by its published SHA-256.
- `transparency-dev/tessera` exact signed commit and module pseudo-version recorded in `source-lock.json`.
- Google `govulncheck` 1.7.0 for source and binary reachability; its vulnerability DB is time-varying, so every run preserves output and date.
- `strace` is required only by the upstream Linux POSIX fault-injection tests; its distro package and runtime closure remain target evidence.
- No paid cloud service, provider credential, background daemon, Docker or Git repository is required by the POSIX one-shot lane.

## 8. Apply order

1. Complete project initialization and source-profile approval.
2. Acquire the exact Tessera source archive through `OFFICIAL-UPSTREAM-ACQUISITION-CORE` and verify its receipt.
3. Place source and log storage on admitted Linux filesystems; never use an NTFS-mounted source as storage.
4. Acquire exact Go and `govulncheck`, then run `build_and_verify_official.sh` into an absent/empty tools directory.
5. Complete the deployment profile, external keys, checkpoint publication/witness and backup policy.
6. Run `run_verified_log.sh create`, `append` or read-only `verify`; retain output and checkpoint hash.
7. Attach the log receipt to the durable document evidence chain. Do not authorize production until recovery, witness/checkpoint, access, monitoring and load gates pass.

## 9. Verification

- Static: materialization hashes, official byte hashes, zero-enabled profile, attribution and required fail-closed script tokens.
- Upstream: exact archive/license/manifest identity; complete Linux source suite including I/O and SIGKILL fault tests; `go vet`; deterministic builds.
- Security: OSV build-list inventory plus official `govulncheck` on source and binary; any reachable/imported result fails.
- Runtime: create and append approved receipts, `fsck` each checkpoint, verify-only mode, absent-output and key-permission negatives.
- Tamper: mutation of any tile or bundle must make official `fsck` non-zero; never auto-repair or overwrite evidence.

## 10. Reconstruction evidence

Verified on 2026-08-27 in WSL2/Linux amd64 with Go 1.26.7. The signed commit archive (2,215,098 bytes; SHA-256 `33bb2394…d779c`) passed the complete suite, fault injection, vet and two deterministic builds. The official CLIs created and verified a ten-entry log, then rejected a truncated bundle. The packaged runner independently created, appended and reverified two canonical test receipts; its usage negative returned 64.

The current commit is a signed main snapshot, not a stable release. Five OSV records remain in three required but uncalled modules; official `govulncheck` reports zero reachable and zero imported-package vulnerabilities for the POSIX CLI. The pack therefore remains `CONDITIONED`, with exact evidence rather than a false production-ready claim.
