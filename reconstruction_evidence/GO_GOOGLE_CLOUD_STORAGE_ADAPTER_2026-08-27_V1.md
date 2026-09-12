# Go Google Cloud Storage Adapter — Reconstruction Evidence V1

Date: 2026-08-27. Scope: exact public SDK identity, local materialization, fail-closed behavior and offline executable gates. No Google Cloud account, credential, paid API or production bucket was contacted.

## Official authority and license

- module: `cloud.google.com/go/storage v1.65.1`;
- official tag: `storage/v1.65.1`;
- commit: `3ab7d1390bbbb40cf197a545b941f9df760a3269`, GitHub verification `valid/verified=true`;
- Go proxy publication: `2026-08-26T17:25:54Z`;
- module ZIP: 998,644 bytes, SHA-256 `0030613d00242a7b2cd0133e86e67c12c066925e88909f185d56706c0bdb6d64`;
- license: Apache-2.0, `LICENSE` SHA-256 `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30`;
- toolchain: official Go 1.26.7 Windows amd64 archive, 74,955,002 bytes, SHA-256 `f4f534a486e4bc3387fa18f08208f2f854b7aaea8a08f2a2d829a914a05abb11`.

The official SDK source establishes the employed primitives: `Conditions.DoesNotExist`, caller-supplied CRC32C with `SendCRC32C`, `KMSKeyName`, `ObjectRetention`, exact `Generation` handles, bucket project number and `ObjectRetentionMode`. Google documents that object retention requires a bucket with the feature enabled, that `Locked` is irreversible, and that existing buckets cannot be silently enabled by this adapter. The local `gcsstorage` package is explicitly `AUTHORED` integration glue; it is not presented as Google-authored application code.

## Materialized artifact

Pack SHA-256: `981bcaaa671a58ed44bd2a7c64bdadb418eca4478dfc9c7457c78b84972ebf34`.

| Path | Bytes | SHA-256 |
|---|---:|---|
| `gcsstorage/crc_test.go` | 148 | `5088535a35d1ef378a202e123617dba278bfc7bcb873e4b651beb26c8b762e4d` |
| `gcsstorage/google_backend.go` | 2,547 | `067c9668555f1b40488400026e1c98ee5f92997cc3e042cade4e2115db5e6242` |
| `gcsstorage/storage_test.go` | 7,898 | `d30537bbc7adab99d5c1f1feff3268d92ddc0689d8b925c7b82b1102850df946` |
| `gcsstorage/storage.go` | 9,175 | `4925eb9ae4400adb161964e068d836a47afe188a9bba51f4b34c1edb4af7e0a6` |
| `go.mod` | 2,751 | `3a0d72638310d1ef66969dd6525df7b447815dc14b0c24833e66a7092582e09a` |
| `go.sum` | 12,153 | `374651b97a68537d2b0e43bd0747aec1974eec33ecb110b457378812b50dec70` |
| `official-artifact-lock.json` | 592 | `2c33e9aa7420835a5427fb055324e2000ed0a5fee50f798ef9550e4d0f3f010a` |
| `provider-profile.template.json` | 766 | `3d3bca488f2fb5febcd67ffde5667139c857f7dac7dd91e30b7d606c54a87aba` |
| `README.md` | 1,568 | `88ddf649dfd4b069543a93948599754587fa182a7fcf4ed494d92ef753283979` |

The root materializer recreated all nine files from Markdown and verified every declared SHA-256. The Google document composition profile then produced 6 packs / 50 implementation files and a valid materialization record.

## Executable behavior

The reconstructed module ran with `GOPROXY=off`:

```text
go mod verify  -> PASS, all modules verified
go test ./...  -> PASS, 8 primary tests / 20 passing test and subtest events
go vet ./...   -> PASS
go build ./... -> PASS
```

Covered behavior includes empty/oversize/unsafe names, owner mismatch, retention-disabled bucket, missing approvals, invalid KMS/mode/date, irreversible lock approval, create failure, exact-generation inspection, every post-write mismatch, no raw identifiers in receipts and partial-effect evidence. The provider mapping compiles against the exact SDK and uses `DoesNotExist`, CRC32C, Cloud KMS and retention before the first write.

## Dependency security evidence

An OSV query batch over the exact 198-module build list first found four reports in three selected modules: `GO-2026-5932` (`x/crypto` OpenPGP), `GO-2026-6179`/`GO-2026-6180` (`x/mod`) and `GO-2026-5970` (`x/text`). This result is retained as `UP-FAIL-079`; it was not hidden by claiming that the newest SDK graph was automatically safe.

`x/text` was elevated to the official fixed version `v0.39.0`, followed by tidy, module verify, tests, vet and build. The remaining build-list reports are not runtime calls: `govulncheck v1.7.0` at commit `617f44b718537dccdea1915395650e0529e3b72e` scanned 47 runtime modules plus Go 1.26.7 and reported:

```text
reachable vulnerabilities: 0
vulnerabilities in imported packages: 0
module-only findings not called: 1 (GO-2026-5932 in x/crypto/openpgp)
```

`x/mod` was not among the 47 runtime modules; the equivalent `cmd/go` flaws are fixed before Go 1.26.6, and this reconstruction uses 1.26.7. A future dependency change must repeat both module inventory and symbol reachability; the module-only advisory remains a visible condition.

## Result and limits

Result: `REBUILD_VERIFIED / CONDITIONED`. The code is immediately materializable and locally executable, but a project cannot enable it until real Google identity/IAM, expected project, bucket retention configuration, KMS grant, cost/quota, audit, legal retention, security admission, strict field evaluation, concurrency/load and recovery gates pass. A `Locked` object cannot be used as an experiment or undone by rollback. A post-write verification failure returns the created generation in a partial receipt and requires reconciliation rather than blind retry.
