# Firebase Cloud Messaging Adapter — Reconstruction Evidence V1

Date: 2026-08-26  
Result: `REBUILD_VERIFIED / CONDITIONED`  
Pack: `GO-FIREBASE-CLOUD-MESSAGING-ADAPTER 0.1.0`  
Pack SHA-256: `a2dac3d16838f21e0d690f913172932ae3454a273b7e428a7b6d929a2ac2d374`

## Authority fixed before implementation

- Official repository: `firebase/firebase-admin-go`.
- Official release/module: `v4.21.0`.
- Exact commit: `eebb06f2a643fbb59b1cb262874a943584475128`.
- GitHub archive: 322,131 bytes; SHA-256 `8d7bbc5621ddf19ddf0ebd74df867690cd3ee1bed7499eabf973eba668284b5e`.
- Apache-2.0 LICENSE: 11,357 bytes; SHA-256 `c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4`.
- Official upstream `go.mod`: SHA-256 `a4baa15771d86ad3f97e729437a5971d168f138440fe742dad595c9d728a1db1`.
- Official upstream `go.sum`: SHA-256 `9cd89415c44406855db8406b673cef99ca8964d630bd77f052e1071aac33e925`.
- Official `messaging/messaging.go`: SHA-256 `b80b2ab57062a50dce0cf3060ae8857e7b9145a1661895a146277705c30210d9`.
- Official `messaging/messaging_batch.go`: SHA-256 `b7896b56b98159814324ddf95497a3af93adc92633405dc762bd292f5fc3b0d9`.
- GitHub reports the release commit unsigned (`verified=false`, reason `unsigned`); this remains `UP-FAIL-027`, not relabeled as signed.

The adapter code is `AUTHORED`; it imports and invokes the exact official module. No Firebase source was copied or falsely attributed.

## Clean toolchain and reconstruction

- Temporary isolated roots under `%LOCALAPPDATA%/Temp/firebase-admin-go-audit-20260826`.
- Official Go `1.26.7 windows/amd64` ZIP: 74,955,002 bytes; SHA-256 `f4f534a486e4bc3387fa18f08208f2f854b7aaea8a08f2a2d829a914a05abb11`.
- No global Go installation and no project credentials were used.
- Root materializer produced 9 files.
- Every reconstructed file matched the authored source tree byte-for-byte by SHA-256.

| Reconstructed path | Bytes | SHA-256 |
|---|---:|---|
| `firebase_push/cmd/firebase-push/main.go` | 1,445 | `90202057ecf0522beb731adf60cdb0556be507a7f02fc40098558309ad6866cd` |
| `firebase_push/firebase_push_test.go` | 5,791 | `b6c143462a7cc724942d9aff81587ba4498c680dd9d4cd91635efb1d4e0fbf81` |
| `firebase_push/firebase_push.go` | 10,172 | `4a6c7c65e5c272e65a37f37f85a633cc9311a323448e294730fec60c1cf34ffe` |
| `firebase_push/go.mod` | 3,076 | `9543e144c5c213983232f03756ef73da52a40e21811112ae7f28496b8d2d0c2c` |
| `firebase_push/go.sum` | 15,850 | `ad8189a073109a42337ef084ed5e2d7762c27e14e9cddc13d809ec9104a8516a` |
| `firebase_push/message.template.json` | 69 | `18be24fdff42998690cae145990c35422bc3ae3bf3e9883a4710ef604e575d8c` |
| `firebase_push/provider-profile.template.json` | 924 | `fe862ef2788ed1cb368159e546dcba585c3a98246656f0c525d2f6311e8922fc` |
| `firebase_push/README.md` | 1,258 | `3ed5935c2a8c2e79e1b186c0fe7a4f64411ffa55c1941d871f88be6d466eb68d` |
| `firebase_push/sdk-source.lock.json` | 841 | `0b2dabefdbcf1e044881ebd24303223c553cc205b29ee9ad420fe942ca3a1f50` |

## Commands and results

Executed first on the authored isolated tree and again on the Markdown reconstruction:

```text
go fmt ./...                                      PASS
go mod tidy                                      PASS
go mod verify                                    PASS: all modules verified
go test ./...                                    PASS: 7 tests
go vet ./...                                     PASS
go build ./cmd/firebase-push                     PASS
go list -m ... firebase.google.com/go/v4         PASS: v4.21.0
materialize_markdown_pack.ps1                    PASS: 9 files
source vs reconstruction SHA-256                 PASS: 9/9 identical
VERIFY_LIBRARY.ps1                               PASS: 36 packs / 325 files
VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Audit        PASS: 57 sources / 6 provider adapters
VERIFY_EXECUTABLE_LIBRARY.ps1 -Mode Foundation   PASS: Firebase module/tests/vet + full foundation
```

The seven tests prove official `messaging.Message` construction, dry-run selection, separate real-send approval, allowlisted/reserved data handling, token/TTL bounds, unknown/trailing JSON rejection, provider failure atomicity, missing message-ID rejection and receipt redaction.

## Failures retained as learning

- `LIB-FAIL-070`: an assumed upstream test path returned 404; the official Git tree became the path authority.
- `LIB-FAIL-071`: a wrong local module import caused an attempted external resolution; the exact module root is now used.
- `LIB-FAIL-072`: Firebase 4.21.0 requires `*time.Duration` for Android TTL; the official type governs and the message is inspected by test.
- `LIB-FAIL-073`: direct temporary binary deletion was policy-blocked; an exact-root validation script removed only the build artifact.
- `UP-FAIL-027`: the upstream release commit remains unsigned.

## Conditions and non-claims

No Firebase project, ADC identity, FCM API, registration token, device, sandbox or production environment was contacted. No notification was sent. Therefore this evidence does not prove delivery, foreground/background behavior, APNs/WebPush settings, end-user consent/opt-out, quotas, billing, load, retention or operational reconciliation. The project-start gate must collect and prove those inputs; `DRY_RUN` is the first permitted provider operation, and `SEND` requires a separate approval.
