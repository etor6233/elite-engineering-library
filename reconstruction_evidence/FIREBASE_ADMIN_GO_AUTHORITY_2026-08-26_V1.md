# Firebase Admin Go Authority — 2026-08-26 V1

Official active repository: `firebase/firebase-admin-go`. Latest observed release on 2026-08-26: `v4.21.0`, published 2026-07-08, commit `eebb06f2a643fbb59b1cb262874a943584475128`. GitHub reported the commit `verified=false`, reason `unsigned`; no signing claim is made.

The exact source archive measured 322,131 bytes with SHA-256 `8d7bbc5621ddf19ddf0ebd74df867690cd3ee1bed7499eabf973eba668284b5e`. Root Apache-2.0 LICENSE: 11,357 bytes, SHA-256 `c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4`. Critical files observed:

- `go.mod`: 2,996 bytes, SHA-256 `a4baa157…db1`;
- `go.sum`: 15,679 bytes, SHA-256 `9cd89415…925`;
- `messaging/messaging.go`: 37,835 bytes, SHA-256 `b80b2ab5…0d9`;
- `messaging/messaging_batch.go`: 17,826 bytes, SHA-256 `b7896b56…0d9`.

The official `fcmClient` exposes `Send` and `SendDryRun`; the latter sets `ValidateOnly: true`. This is a valid high-value authority candidate for Firebase Cloud Messaging. It is not yet an Elite implementation pack: exact adapter module/transitive graph, credential/project approval, token/topic policy, dry-run/real-send separation, receipts, quota/cost, consent and provider reconciliation still require implementation and clean gates.

Failure learning: `LIB-FAIL-070` records that the first critical-file query guessed nonexistent `messaging/example_test.go` and stopped on HTTP 404; the official Git tree was enumerated before selecting actual paths. `UP-FAIL-027` retains the unsigned release commit. No temp source or binary was retained.
