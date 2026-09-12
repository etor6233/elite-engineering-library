# Human Review Leader Code Reaudit — V2

Date: 2026-08-28  
Scope: official AWS document-extraction human review code, exact-source admission only  
Governing result: `SAMPLE_ONLY_REJECTED_FOR_IMMEDIATE_ADOPTION`

## Exact official source

- Owner/repository: Amazon Web Services, `aws-samples/sample-scalable-intelligent-document-processing-with-amazon-bedrock-data-automation`.
- Exact commit: `480f4eb83639c337ec75826c08614e4f7b098700`.
- GitHub commit verification observed: verified, reason `valid`.
- Release: none; the identity is the exact commit, not `main`.
- License: MIT-0, `LICENSE` SHA-256 `ef47b4ae2a1a8d38ef87b38dc5927967e7b472ffa3f7d125e89dbc8353a1e7af`.
- Commit archive: 18,418,459 bytes; SHA-256 `886b163829c95cd45833aea016a914209c4a84e24563f2dfe91da74902ff4986`.
- Official repository: https://github.com/aws-samples/sample-scalable-intelligent-document-processing-with-amazon-bedrock-data-automation
- Official A2I output contract: https://docs.aws.amazon.com/sagemaker/latest/dg/a2i-output-data.html

## What the exact code really provides

| Capability | Exact code evidence | Result |
|---|---|---|
| Per-field confidence and geometry | `multipagepdfbda_confidence/lambda_function.py` builds labels with value, confidence, page, bounding box and vertices | present |
| Editable human correction | `Custom-Template` renders each extracted value in an editable input and submits the edited name/value map | present |
| Private workforce procedure | README requires a private A2I team and verified worker invitation | documented/manual |
| Request/page binding | `multipagepdfbda_analyzepdf/lambda_function.py` derives a page-specific human-loop ID and binds the page S3 object | present, convention-based |
| Durable corrected answer | completion Lambda writes the selected answer map to S3; wrap-up emits separate BDA and human-response objects | present |
| Multi-page completion | DynamoDB tracks page completion and callback tokens before Step Functions success | present |
| Reviewer identity retained in business evidence | AWS A2I output supplies `workerId`, identity-provider `sub`, acceptance/submission times and duration, but the sample selects only `humanAnswers[0].answerContent` | absent |
| Mandatory decision/correction reason | no reason input or server-side validation exists | absent |
| Assignment lease/expiry | no claim expiry or stale-assignment rejection exists in application code | absent |
| Optimistic version/CAS | completion performs an unconditional DynamoDB update; no expected version or conditional decision write exists | absent |
| Dual control/separation of duties | no second independent reviewer or quorum decision exists | absent |
| Automated tests | repository contains no test/spec files | absent |

## Executed verification

- All 12 Python files parsed with CPython: `PYTHON_PARSE files=12 errors=0`.
- Exact archive, license, package manifest and nine focal artifacts were hashed and inserted into `upstream-source-lock.json`.
- Acquisition suite: `UPSTREAM_ACQUISITION_TEST_PASS negatives=9`.
- Source-profile suite: `SOURCE_PROFILE_TEST_PASS valid=14 negatives=5 positives=1`.
- AWS document profile: `SOURCE_PROFILE_VALID ... selected=15`.
- Document-intelligence profile: `SOURCE_PROFILE_VALID ... selected=41`.
- No dependency install, AWS deployment, paid provider call or corpus accuracy run was authorized or performed.

## Production blockers proven from the source

- The root manifest pins only part of the graph; `constructs` and `cdk-nag` are ranges, while the Lambda layer leaves `pytz`, `requests` and `requests_aws4auth` unpinned.
- The CDK stack uses `RemovalPolicy.DESTROY` for evidence-bearing resources.
- CDK Nag suppressions explicitly describe the implementation as a POC and defer least privilege, S3 bucket policy/access logging and KMS hardening to the customer.
- The A2I workflow and private workforce are configured manually in the console, so a one-command reproducible deployment is not supplied.
- S3 persistence of edited values does not make those values an authenticated, reasoned, conflict-safe or accuracy-proven business record.

## Admission decision

The source is retained because its exact editable-field UI, confidence/geometry mapping, A2I page binding and S3 consolidation are useful official AWS mechanics. It is rejected for immediate composition or deployment. An agent must not copy it into a production project or claim complete document-review assurance until an official revision retains reviewer identity/timestamps, requires a reason, enforces lease/CAS and dual control where policy requires it, publishes locked clean dependencies and automated tests, and the target project proves corpus accuracy, security, recovery, load, cost and rollback.

