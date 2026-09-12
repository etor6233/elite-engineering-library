# AWS Amazon SQS best-practices CDK sample — 2026-08-28 V1

## Decision

`aws-samples/amazon-sqs-best-practices-cdk` is an official, public, MIT-0 AWS sample, but it is **`SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`**. It is useful as an exact reference for `S3 -> Lambda -> SQS -> Lambda -> DynamoDB`, a source queue with DLQ, queue-age alarm and least-privilege intent. It is not the complete object-event worker required by this library and none of its application code is promoted as Elite reusable code.

## Exact upstream identity

- repository: `https://github.com/aws-samples/amazon-sqs-best-practices-cdk`;
- default branch: `main`; repository active/not archived at audit time;
- head commit: `54143c7efd610b2345964951403bbde628f54ba0`, tree `944b437b9d94519353d6a6ed222361914e210f42`, GitHub verification `valid`;
- no tags and no releases;
- ZIP: 58,128 bytes; SHA-256 `d6dea055180dc4951c1fba5dca58834fdd3a3089f72ca9cf0518c87fce3e9c9b`;
- expanded source: 21 files / 77,365 bytes;
- MIT-0 `LICENSE`: 947 bytes; SHA-256 `7cb750713252efd1d578837ba8785b61319109906f0c10b87adab7cf4badfc42`; no NOTICE/third-party notice file;
- the four application/dependency files were last changed by the unsigned initial commit `0777c2b1b6b054be5a3c250a4b99616a2dd2ea17` dated 2023-06-03. The signed head only updates README and is dated 2023-06-14; a recent repository activity timestamp is not represented as recent application code.

Critical file locks:

| Path | Bytes | SHA-256 |
|---|---:|---|
| `README.md` | 3,940 | `3c72d9d1aea8772a8f010fe0065f08c12b60e365982fd484ea1ee031e2b69c13` |
| `requirements.txt` | 47 | `9be4442e04303763407bbecb7e0db7de58226e3ceccc2245101df3a7410ca310` |
| `app.py` | 930 | `3124261fd6a70b9d0a323c6cfe27f3a3b5baa27cf993cfb8af1e431268d0e906` |
| `sqs_blog/sqs_blog_stack.py` | 6,674 | `01e841e2af0e5101bb4cd547ddf388047b1b3954eb9fa5b37fbc1702f24c98c8` |
| `sqs_blog/lambda/CSVProcessingToSQSFunction.py` | 1,544 | `eb01e9dc775dba79b13ce59a5fe96b849662f4aedbddeee289ba733a6b9ec622` |
| `sqs_blog/lambda/SQSToDynamoDBFunction.py` | 1,033 | `21ee326e8d968828c10854a31a31f7e4c2d0c3ec84f1c491654c678b9179979d` |
| `sqs_blog/sample_file.csv` | 715 | `87dd10632549c97f974bb93730f36e06ab2dff0f8e9bc1670ddbfab67f24482d` |

## Reconstruction and dependency evidence

- five Python files compiled in memory with CPython 3.14.4;
- upstream supplies 0 test files, 0 workflows and 0 dependency locks; four `.pyc` files are committed;
- `requirements.txt` fixes only `aws-cdk-lib==2.76.0` and leaves `constructs>=10,<11`; an isolated CPython 3.12.14 install resolved 14 distributions and `pip check` passed, but that is a dated mutable resolution, not an upstream lock;
- Google OSV-Scanner 2.5.1 exact binary `25e42f5e…bfb6` scanned 15 package rows and closed 1. One package/version is affected: `aws-cdk-lib 2.76.0`, `GHSA-464c-974j-9xm6`, score 3.3/LOW, fixed in 2.253.0. Report: 64,580 bytes; SHA-256 `184a13abb74d283f60b3f84584c51fecee9ea37048957307d96e33161c5327ef`;
- direct `app.py` execution with explicit temporary `CDK_OUTDIR` synthesized 17 CloudFormation resources. Template: 15,993 bytes; SHA-256 `683a67c3188369615ee7c4588f32ce2c5ab6cea93695d1b39ae8585b0a02ba15`;
- Node 24 emitted the `DEP0169 url.parse()` deprecation from the old dependency graph. No account, credentials, deployment or billable AWS call was used.

The synthesis proves reconstructibility only. It generated two Python 3.8 application Lambdas plus a Python 3.9 notification Lambda, one SQS event-source mapping with batch size 10 and no `FunctionResponseTypes`, an unencrypted-by-template standard queue with DLQ/maxReceiveCount 5, a bucket without versioning or explicit encryption, and a DynamoDB table without PITR configuration.

## Executed semantic audit

An audit-only nine-assertion harness imported the exact two handlers with fake AWS clients. It proved:

- processing the same stable SQS message twice performs two DynamoDB writes with different random UUID partition keys;
- if the second record fails, the first DynamoDB write already occurred and the entire handler raises; there is no partial-batch response;
- `send_message_batch` returning an item in `Failed` is ignored and the CSV handler returns normally;
- an S3 notification containing two records reads only `Records[0]`.

The source also has no object version/eTag/sequencer idempotency key, conditional business write, owner/lease/fencing, payload schema, tenant binding, quarantine/malware boundary, DLQ consumer/redrive/replay/reconciliation, immutable audit receipt, load/recovery tests or target-account E2E.

## Comparison with current AWS authority

AWS Prescriptive Guidance currently requires a DLQ to avoid snowball failure, `ReportBatchItemFailures`, per-job idempotency and code-level metrics for SQS partial-batch processing. It recommends Powertools Batch Processing and Idempotency. This sample configures only the DLQ/age alarm portion and does not implement the other controls: `https://docs.aws.amazon.com/prescriptive-guidance/latest/lambda-event-filtering-partial-batch-responses-for-sqs/best-practices-partial-batch-responses.html`.

## Admission boundary

Do not copy, compose or deploy this sample as the library's worker. Keep it as a pinned rejected reference so an agent cannot mistake the AWS organization name or README “best practices” claim for production evidence. A future official revision must at minimum provide current runtimes, a frozen clean graph, tests, checked `SendMessageBatch` partial results, `ReportBatchItemFailures`, stable idempotency and ordering keys, conditional business effects, DLQ replay/reconciliation, secure storage/IAM and live recovery/load evidence before re-audit.
