# AWS Transactional Outbox Pattern — exact-source re-audit V1

## Scope and verdict

Date: `2026-08-27` (`America/Buenos_Aires`).

Question: can the official AWS sample close `UP-FAIL-109` and be copied immediately as the durable audit/outbox implementation for authenticated human decisions?

Verdict: `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`.

The [AWS Prescriptive Guidance pattern](https://docs.aws.amazon.com/prescriptive-guidance/latest/cloud-design-patterns/transactional-outbox.html) is an official architectural authority for avoiding dual writes. The exact [AWS sample repository](https://github.com/aws-samples/transactional-outbox-pattern) is retained as licensed source and negative evidence, but its code is not admitted as an executable production baseline. No local patch is attributed to AWS.

## Exact upstream identity

| Field | Evidence |
|---|---|
| repository | `aws-samples/transactional-outbox-pattern` |
| state | active, not archived or disabled; default branch `main` |
| commit | verified signed `23e0519c7a8048c6db4acd8f3d7de534a9d5deac` |
| tree | `4d41530ad300cd46aeea43171e3e65bd19e5e881` |
| commit date | `2024-02-19T07:00:40Z` |
| tags/releases | GitHub API raw responses `[]` / `[]` on `2026-08-27`; SHA-256 of each trimmed payload `4f53cda18c2baa0c0354bb5f9a3ecbe5ed12ab4d8e11ba873c2f11161202b945` |
| commit ZIP | `1,250,782` bytes; SHA-256 `662baa72971fba86a558bb927c1a3fbef6cd5d64d72aa96fbf04c377b6698fc6` |
| archive inventory | 74 files; `1,542,500` uncompressed bytes |
| license | MIT-0; `LICENSE` 1,030 bytes/SHA-256 `0ae7c59c8ace792b1d14e816ed148c144d1901c3de1bd9b2e63d024ed4edfe39` |
| README | 6,243 bytes/SHA-256 `3c65667da4bf8ce7799a2b20c8c0ef74e045ea32bfe2668b1c34cfc373a311b2` |

The source is now entry `aws-transactional-outbox-pattern-23e0519` in the 83-source exact lock. Eleven audited artifacts are independently hash-bound inside the archive.

## Official code inspected

| Artifact | SHA-256 | Result |
|---|---|---|
| relational `FlightController.java` | `0cb9df3b446fe33cc736b708b1525f463fb40bf4a565833fe07cfaebc0990e79` | business row and outbox row share one Spring transaction |
| relational `QueueService.java` | `83d34c805c46506dfab5bd4ca1e2de4de6bccf47a165eb488398ae242b609560` | polls ordered rows, calls SQS batch, then deletes all selected rows |
| relational test | `447ad3276d834cd08fbe526d421ae56d9c95f9b195c7046114cfcf52b2eb6a5f` | only `contextLoads`; no outbox behavior assertion |
| CDC `DynamoDBMessageConverter.java` | `6665a4b884f06e2eef68e70342e66412e187518442581a0212528ffc83b6a0d7` | reconstructs fields but does not assign the DynamoDB item ID |
| CDC `QueueService.java` | `c352aacc20ed140cc1ea2be216ed1a786cf598488e84c26aa6ab79a8fdf06ec3` | dereferences missing `flight.id`; logs/returns on serialization `IOException`; imports SDK-v1 exception beside an SDK-v2 SQS client |
| CDC test | `447ad3276d834cd08fbe526d421ae56d9c95f9b195c7046114cfcf52b2eb6a5f` | only `contextLoads`; no CDC/redelivery assertion |
| CDK `baseStack.ts` | `b8ecd0f084ec41089238646465b20ff3a34519c3da21d1961b0b724216914e43` | FIFO/DLQ present; destructive log removal and no demonstrated restore |
| CDK `cdcStack.ts` | `09c56da0f201542458fd20e2d1ba549608f43e5e1030d9bc6a46c94539ccb7d5` | DynamoDB/Kinesis wiring; `kinesis:*`, `dynamodb:*`, wildcard resource and deletion protection disabled |
| CDK test | `1c8f43aa2981a09024f0d760eb62becc2c401f166ef366944dee3ff64b843bb0` | Jest function runs, but all CDK assertions/imports are commented |
| npm lock | `c412150b5380fa76b8df4f545eee2ba1d577976d5d2461a658c023e747d5152d` | reproducible but vulnerable |
| relational POM | `dcaed6e3477b642b427a0f633908b27fce180db1f43743acb09776048701609e` | Spring Boot 3.1.2 / Java 17 graph; vulnerable |

## Rebuild and tests

Tooling was isolated under a temporary audit directory; nothing was installed globally and no AWS account or billable service was used.

- Amazon Corretto `17.0.20.1` (`17.0.20.10.1`) from the official latest Windows JDK ZIP: 188,084,668 bytes, SHA-256 `af002c5f7dd3d09ad4a1f1643266e28fb368077de50560534a2aa84e376ddfde`.
- Apache Maven `3.9.16`: official ZIP 9,395,475 bytes; official SHA-512 matched `ed41650d42485cfc243fad22158caf9cbb5dc408ce7a09ddb94dd42a019de929ca43065bfa450612cf12bf78b5cafa3884b96c090de326ff590448c933454af3`.
- Node `24.14.1`, npm `11.19.0`, OSV Scanner `2.5.1`.

Results:

1. Relational `mvn test`: nine main sources and one test compiled; `1` test errored because `contextLoads` attempted DDL against the deliberately unreachable external PostgreSQL endpoint. The repository has no hermetic database fixture and no functional outbox test.
2. Relational `mvn clean package -DskipTests`: `BUILD SUCCESS`; this proves compilation/packaging only.
3. CDC `mvn clean package -DskipTests`: eight main sources and one test compiled; `BUILD SUCCESS`; this proves compilation/packaging only.
4. CDC `mvn test` with non-secret dummy credentials and metadata disabled: `1` error before cloud access because both `ApacheSdkHttpService` and `UrlConnectionSdkHttpService` are present and no implementation is selected.
5. Infra `npm ci --ignore-scripts --no-audit --no-fund`: 316 packages installed from the exact lock.
6. Infra `npm run build`: TypeScript PASS.
7. Infra Jest: `1/1` reported PASS, but the source contains no active infrastructure assertion; it is not counted as behavioral evidence.

No PostgreSQL, DynamoDB, Kinesis, SQS, ECS, WAF, restore, replay, race or live deployment gate ran.

## Security and dependency evidence

OSV Scanner `2.5.1` results from exact upstream manifests/lock:

| Lane | Packages observed | Finding rows | Unique advisories | Affected package names |
|---|---:|---:|---:|---:|
| CDK `package-lock.json` | 321 | 28 | 28 | 13 |
| relational Maven | transitive graph resolved from POM | 120 | 117 | 28 |
| CDC Maven | transitive graph resolved from POM | 143 | 139 | 40 |

The npm findings include six advisories on `aws-cdk-lib 2.90.0`. The Java graphs include obsolete Spring Boot/Tomcat/Netty/Jackson and related versions. Counts are retained as audit evidence; they are not a reachability verdict and do not authorize automated upgrades.

## Why it does not close durable audit

AWS's pattern correctly states that the domain change and outbox event must be atomic, consumers must be idempotent because duplicate delivery can occur, ordering matters and rollback must not publish an event. The exact sample does not demonstrate the guarantees needed for the library's approval audit gap:

- partial `SendMessageBatch` failures are not inspected before deleting the full relational batch;
- concurrent polling/lease/CAS and post-send database rollback behavior are not tested;
- CDC conversion omits the identifier subsequently required for FIFO deduplication;
- serialization failure is logged without a failed message contract;
- the SDK HTTP implementation is ambiguous in the resolved CDC graph;
- no real retry, replay, poison-message, DLQ redrive, reconciliation or immutable receipt test exists;
- the CDK uses broad IAM and destructive defaults unsuitable as an enterprise baseline;
- all three dependency lanes have known advisories.

Therefore `UP-FAIL-109` remains open. The official architecture can govern a future implementation, but the library will not invent or falsely attribute the missing hardened code to AWS.

## Admission and follow-up gate

The source may be acquired only for research through the exact lock; it is absent from all automatic source profiles. Any future candidate must originate from a maintained official source or be clearly classified local `AUTHORED`/`ADAPTED`, and must prove at minimum:

```text
atomic domain+outbox write
→ claim/lease with concurrent workers
→ per-entry broker acknowledgement
→ idempotent consumer and ordering contract
→ retry/DLQ/redrive/replay
→ crash after send / before mark
→ crash after commit / before dispatch
→ reconciliation and immutable receipt
→ least privilege, encryption, retention, backup/restore
→ clean dependency admission
```

This evidence adds no production-ready outbox code. It prevents an official but insufficient sample from being copied as “elite”.
