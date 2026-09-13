# V402335 — distribution verifier count repair

AUTHORED library verifier glue; no upstream attribution. FAIL985 corrected to the actual admitted pack file count. Thirteen mappable static counts plus the separate21-file history plan audited; only HTTP metrics drifted. Original failed Preflight preserved. No source330 rebuild, dependency or runtime change.

```json
{
  "state": "PASS_STATIC_COUNT_AUDIT",
  "scope": "ROOT_DISTRIBUTION_VERIFIER_ONLY",
  "failed_preflight_receipt_sha256": "b4722bb8d10191b8949c172dbab35e5f0d9ef65901a8090d448e7d8999cfe583",
  "before_script_sha256": "be9a1319bed5c545e3fee0b13c7e475b58c79e7116b0d87b83739b3e5e5685f4",
  "after_script_sha256": "c2ecab51844657f1bedb2befa47b62cd9cec89bfb0dabe5aa8beb21bca5e36a6",
  "audit_before": {
    "state": "AUDITED_ALL_MAPPABLE_STATIC_MATERIALIZATION_COUNTS",
    "checks": [
      {
        "line": 626,
        "variable": "httpMetricsRoot",
        "pack": "GO_HTTP_METRICS_REFERENCE.md",
        "prefix": "",
        "recursive": true,
        "expected": 8,
        "actual": 11,
        "matches": false,
        "pack_sha256": "00aad47aec6767fcc2e125822e5a5dc261edcd4f8014e24910fabd3b0aa7cded"
      },
      {
        "line": 1148,
        "variable": "officialPgDurableDirectory",
        "pack": "MICROSOFT_PG_DURABLE_HUMAN_HANDOFF.md",
        "prefix": "microsoft_pg_durable_human_handoff/upstream/",
        "recursive": true,
        "expected": 26,
        "actual": 26,
        "matches": true,
        "pack_sha256": "9cae15cd99c021da0579a837b7063cdb67d680189c0a2ea8dcdec18eb1f3df53"
      },
      {
        "line": 1178,
        "variable": "officialAvmSftpDirectory",
        "pack": "MICROSOFT_AVM_SECURE_SFTP_INTAKE.md",
        "prefix": "azure_avm_secure_sftp_intake/",
        "recursive": true,
        "expected": 8,
        "actual": 8,
        "matches": true,
        "pack_sha256": "63ab14f8a368c9a68240c74f67f0901b132ec20789f8400f29ba08904dae2b1b"
      },
      {
        "line": 1202,
        "variable": "secureEmailMimeDirectory",
        "pack": "SECURE_EMAIL_MIME_QUARANTINE_CORE.md",
        "prefix": "secure_email_mime_quarantine_core/",
        "recursive": false,
        "expected": 5,
        "actual": 5,
        "matches": true,
        "pack_sha256": "4904e7cf8db690df96ec0aaa3ecfc16ce7575a6c33ffcbbf90ba175779a9a335"
      },
      {
        "line": 1216,
        "variable": "awsSesReceiverDirectory",
        "pack": "AWS_SES_IMMUTABLE_EMAIL_RECEIVER.md",
        "prefix": "aws_ses_immutable_email_receiver/",
        "recursive": false,
        "expected": 9,
        "actual": 9,
        "matches": true,
        "pack_sha256": "239a8194b8ed8b04ce0c87d76a8dea938b9df57e7373691cb930105e900dfe81"
      },
      {
        "line": 1234,
        "variable": "awsGuardDutyReleaseDirectory",
        "pack": "AWS_GUARDDUTY_IMMUTABLE_RELEASE_GATE.md",
        "prefix": "aws_guardduty_immutable_release_gate/",
        "recursive": false,
        "expected": 9,
        "actual": 9,
        "matches": true,
        "pack_sha256": "538e288ec6838f7345cb5c2f0e394fdb3a05292f3f1f7e1907f9ce44eaf03f68"
      },
      {
        "line": 1248,
        "variable": "awsGuardDutyIdpDispatchDirectory",
        "pack": "AWS_GUARDDUTY_MAGIKA_IDP_DISPATCH_GATE.md",
        "prefix": "aws_guardduty_magika_idp_dispatch/",
        "recursive": false,
        "expected": 11,
        "actual": 11,
        "matches": true,
        "pack_sha256": "d192efe903a02d0c5db5d55dfbc7329cbd6142213f87d2755033c7a99ef45012"
      },
      {
        "line": 1268,
        "variable": "awsIdpHandoffDirectory",
        "pack": "AWS_IDP_IMMUTABLE_EVALUATION_HANDOFF.md",
        "prefix": "aws_idp_immutable_evaluation_handoff/",
        "recursive": false,
        "expected": 8,
        "actual": 8,
        "matches": true,
        "pack_sha256": "68d12314d2b5853eca1898988f88bc1c9a3e31c49d6304bf86c0da6d1da074af"
      },
      {
        "line": 1282,
        "variable": "awsIdpDecisionDirectory",
        "pack": "AWS_IDP_EVALUATION_DECISION_WORKER.md",
        "prefix": "aws_idp_evaluation_decision_worker/",
        "recursive": false,
        "expected": 8,
        "actual": 8,
        "matches": true,
        "pack_sha256": "9040e952d08fb1bead9ab9f3b56a243aa96f17b44015bb452cf217b2e6c6417e"
      },
      {
        "line": 1296,
        "variable": "awsIdpPersistenceDirectory",
        "pack": "AWS_IDP_POSTGRES_PERSISTENCE_BOUNDARY.md",
        "prefix": "aws_idp_postgres_persistence_boundary/",
        "recursive": false,
        "expected": 11,
        "actual": 11,
        "matches": true,
        "pack_sha256": "e2c7be0ff5de0aba7f7ef17af4e57dd50935c7ac8bdd821ab7a76d55c14d1bab"
      },
      {
        "line": 1311,
        "variable": "debeziumPostgresOutboxDirectory",
        "pack": "DEBEZIUM_POSTGRES_OUTBOX_RUNTIME.md",
        "prefix": "debezium_postgres_outbox_runtime/",
        "recursive": false,
        "expected": 10,
        "actual": 10,
        "matches": true,
        "pack_sha256": "68f287784dd24d377a0ec5f727f51c05ff0d81a9995a28ba718d477b69255347"
      },
      {
        "line": 1326,
        "variable": "debeziumPostgresInboxDirectory",
        "pack": "DEBEZIUM_POSTGRES_INBOX_CONSUMER.md",
        "prefix": "debezium_postgres_inbox_consumer/",
        "recursive": true,
        "expected": 24,
        "actual": 24,
        "matches": true,
        "pack_sha256": "f044c0b63ee3dd5624c6cb8de742710284dfe8e178e3897b06bcd52007efff05"
      },
      {
        "line": 1353,
        "variable": "googleCelMappingDirectory",
        "pack": "GOOGLE_CEL_DOCUMENT_MAPPING.md",
        "prefix": "google_cel_document_mapping/",
        "recursive": true,
        "expected": 7,
        "actual": 7,
        "matches": true,
        "pack_sha256": "9f1e4a49f7e1351c181db5fb2705dd50f02603b5802d04e071c03b443195f935"
      }
    ],
    "drift": [
      {
        "line": 626,
        "variable": "httpMetricsRoot",
        "pack": "GO_HTTP_METRICS_REFERENCE.md",
        "prefix": "",
        "recursive": true,
        "expected": 8,
        "actual": 11,
        "matches": false,
        "pack_sha256": "00aad47aec6767fcc2e125822e5a5dc261edcd4f8014e24910fabd3b0aa7cded"
      }
    ]
  },
  "separate_history_plan": {
    "files": 21,
    "sha256": "743fb9db0f15096b1209dabc90faa01f0a4797caa74c378da2790c8ffa28852b"
  },
  "fix": "HTTP metrics count8 to11; exact equality retained; neutral step label. Product selected payload unchanged.",
  "full_preflight": "Pending next frozen distribution run; this receipt does not claim full Preflight PASS."
}
```
