# AWS Powertools Idempotent SQS Batch Component

## 1. Metadata

```yaml
pack_id: "AWS-POWERTOOLS-IDEMPOTENT-SQS-BATCH-COMPONENT"
pack_version: "0.1.0"
status:
  authority: ELITE_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa byte-verbatim el ejemplo oficial AWS Powertools Python 3.34.0 que combina idempotencia DynamoDB por registro con respuesta parcial SQS, más los templates oficiales separados de tabla/IAM y SQS cifrada/DLQ/ReportBatchItemFailures, un lock exacto y cuatro contratos offline."
stacks: ["Python 3.12", "AWS Lambda Powertools Python 3.34.0", "AWS SAM", "Amazon SQS", "Amazon DynamoDB"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "AWS-SECURE-QUARANTINE-INTAKE 0.1.x"]
incompatible_with: ["adopción sin clave idempotente de negocio", "efecto empresarial no condicional", "templates fusionados atribuidos falsamente a AWS", "producción sin DLQ/replay/reconciliation"]
license_expression: "LicenseRef-Workspace-Owner AND MIT-0"
upstream_sources: ["https://github.com/aws-powertools/powertools-lambda-python/tree/376757161b002f0c2f5d19d5cdf9def5b8c45704"]
verified_at: "2026-08-28"
```

## 2. Applicability

Use este pack cuando un proyecto AWS necesite un punto de partida oficial para procesar lotes SQS con respuesta parcial e idempotencia DynamoDB por registro. Es adecuado como componente reusable condicionado después de que el blueprint fije la identidad empresarial de la operación, el efecto durable, la política de expiración y la topología AWS.

No lo use como worker completo de eventos S3, como replay/reconciliation, ni como autorización productiva. El handler oficial devuelve el body como efecto demostrativo y usa `messageId`; el proyecto debe probar que esa identidad representa la operación de negocio o elegir otra clave conforme a la documentación oficial.

## 3. Architecture contract

Los cinco archivos bajo `upstream/` son bytes `VERBATIM` de AWS Powertools for Lambda Python 3.34.0, commit `376757161b002f0c2f5d19d5cdf9def5b8c45704`, archivo SHA-256 `b7b8b059…c4742`, licencia MIT-0. AWS publica por separado el ejemplo handler, el template de idempotencia y el template batch SQS; el pack preserva esa frontera y no inventa una plantilla AWS combinada.

El handler registra el tiempo restante de Lambda, aplica `idempotent_function` a cada `SQSRecord` y devuelve `process_partial_response`. Un template aporta tabla DynamoDB con TTL, pay-per-request y cuatro acciones item-scoped; el otro aporta SQS managed SSE, DLQ, redrive, visibility timeout seis veces el timeout y `ReportBatchItemFailures`.

Fallos/límites: no hay semántica S3 version/eTag/sequencer, owner fencing después de expiración/reclaim, efecto empresarial real, replay, reconciliation ni prueba en la cuenta destino. Performance/costo dependen de batch, concurrencia, DynamoDB y payload elegidos. Rollback: retirar sólo la composición del proyecto y preservar receipts/evidencia. Toda actualización reabre admisión, hashes, SCA, tests y target E2E.

## 4. Exact file manifest

```text
CREATE aws_powertools_idempotent_sqs_batch/README.md
CREATE aws_powertools_idempotent_sqs_batch/source-lock.json
CREATE aws_powertools_idempotent_sqs_batch/test_contracts.py
CREATE aws_powertools_idempotent_sqs_batch/upstream/LICENSE
CREATE aws_powertools_idempotent_sqs_batch/upstream/integrate_idempotency_with_batch_processor.py
CREATE aws_powertools_idempotent_sqs_batch/upstream/integrate_idempotency_with_batch_processor_payload.json
CREATE aws_powertools_idempotent_sqs_batch/upstream/idempotency_sam.yaml
CREATE aws_powertools_idempotent_sqs_batch/upstream/sqs_batch_processing.yaml
```

## 5. Materialization blocks

### FILE: `aws_powertools_idempotent_sqs_batch/README.md`
```yaml
block_id: "AWS-POWERTOOLS-IDEMPOTENT-SQS-BATCH:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local packaging constrained by exact AWS upstream"
license: "LicenseRef-Workspace-Owner"
sha256: "e138590602d2be3a36ab8e3160259b36aad10c7163b28c5b705c1b6870f46499"
variables: []
secrets_allowed: false
```
````markdown
# AWS Powertools idempotent SQS batch component

This directory contains byte-verbatim files from AWS Powertools for Lambda Python 3.34.0 at commit `376757161b002f0c2f5d19d5cdf9def5b8c45704`. The upstream handler is AWS's documented composition of per-record DynamoDB idempotency and SQS partial-batch processing. The two upstream SAM templates separately provide the idempotency table/IAM boundary and the SQS queue/DLQ/`ReportBatchItemFailures` boundary.

The examples remain separate because AWS publishes them separately. Do not claim that AWS published a single merged deployment template. A project must select its business idempotency key, replace the demonstration `record_handler` effect, and create a project-owned composition with tests. The upstream example's `messageId` is suitable only when producer semantics prove that SQS message identity equals business-operation identity.

Before deployment, acquire the exact Powertools 3.34.0 wheel through `OFFICIAL-UPSTREAM-ACQUISITION-CORE`, preserve MIT-0, configure `IDEMPOTENCY_TABLE`, enable `ReportBatchItemFailures`, and prove DynamoDB, retry, DLQ, replay, reconciliation, IAM, encryption, observability and target-account E2E gates. These files do not provide S3 sequencer ordering, owner fencing after record expiry, a business effect, replay automation or reconciliation.

Run `python -m unittest -v test_contracts.py` for the offline byte and contract gate. A PASS proves exact packaging and the narrow published controls; it does not authorize production.
````

### FILE: `aws_powertools_idempotent_sqs_batch/source-lock.json`
```yaml
block_id: "AWS-POWERTOOLS-IDEMPOTENT-SQS-BATCH:lock:v1"
operation: CREATE
provenance: AUTHORED
source: "exact inventory derived from AWS Powertools Python v3.34.0"
license: "LicenseRef-Workspace-Owner"
sha256: "f89ea094fd7eea7d9ceb39eb32ab1d6f6f15d78af50b35caa291e81978427cc1"
variables: []
secrets_allowed: false
```
````json
{
  "source_id": "aws-powertools-python-3.34.0",
  "owner": "Amazon Web Services",
  "repository": "aws-powertools/powertools-lambda-python",
  "release": "v3.34.0",
  "commit": "376757161b002f0c2f5d19d5cdf9def5b8c45704",
  "archive_sha256": "b7b8b0599b039275b68138837a4644061f84fbedba624c2f434d4d9e7b5c4742",
  "wheel_sha256": "ab1354c58085ccecf92e9c548b1336ad5bae72c7ae23e61d2aa0ff15520d679c",
  "license_expression": "MIT-0",
  "classification": "PINNED_COMPONENT_CONDITIONED",
  "files": [
    {
      "path": "upstream/LICENSE",
      "source_path": "LICENSE",
      "sha256": "bcc1d2d9e0609b03d420c3dbecb438a7b89bfc53d54176ab56301ab1f96b4f0a"
    },
    {
      "path": "upstream/integrate_idempotency_with_batch_processor.py",
      "source_path": "examples/idempotency/src/integrate_idempotency_with_batch_processor.py",
      "sha256": "d364d7b6ac242568e7c25bf4292da1576b26c0129ecf431fe4912b9191a2df29"
    },
    {
      "path": "upstream/integrate_idempotency_with_batch_processor_payload.json",
      "source_path": "examples/idempotency/src/integrate_idempotency_with_batch_processor_payload.json",
      "sha256": "2a5bab34e2865b8a889033f01642fce1f6841fe092bf84ebb10ae168cf3e12f0"
    },
    {
      "path": "upstream/idempotency_sam.yaml",
      "source_path": "examples/idempotency/templates/sam.yaml",
      "sha256": "27d5a2db8673bd0e1c45799168b68dbde4f3e2e7d025c5c3c7148eaf332292a5"
    },
    {
      "path": "upstream/sqs_batch_processing.yaml",
      "source_path": "examples/batch_processing/sam/sqs_batch_processing.yaml",
      "sha256": "df8341badb7c81de8d7818846071d48554f3dbdb279b8043b00eb5316941ba0e"
    }
  ],
  "proven_controls": [
    "per-record idempotent_function composition",
    "Lambda remaining-time registration",
    "SQS partial-batch response generation",
    "DynamoDB TTL and least-privilege item actions",
    "SQS managed encryption, DLQ, redrive count and visibility-timeout ratio",
    "ReportBatchItemFailures event-source configuration"
  ],
  "production_blockers": [
    "business idempotency key and effect are project-specific",
    "published templates are separate and require project-owned composition",
    "S3 version, eTag and sequencer ordering are not implemented",
    "expired-record owner fencing is not implemented",
    "DLQ replay and reconciliation are not implemented",
    "target AWS account, load and recovery evidence are absent"
  ]
}
````

### FILE: `aws_powertools_idempotent_sqs_batch/test_contracts.py`
```yaml
block_id: "AWS-POWERTOOLS-IDEMPOTENT-SQS-BATCH:tests:v1"
operation: CREATE
provenance: AUTHORED
source: "offline integrity and narrow contract tests over exact AWS bytes"
license: "LicenseRef-Workspace-Owner"
sha256: "11f34a47d4e3aebe374ea7cac5b84c73bf546ae618ddd48d9abfee337b832519"
variables: []
secrets_allowed: false
```
````python
import ast
import hashlib
import json
import pathlib
import unittest


ROOT = pathlib.Path(__file__).resolve().parent


class AwsPowertoolsIdempotentSqsBatchContracts(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.lock = json.loads((ROOT / "source-lock.json").read_text(encoding="utf-8"))

    def test_source_identity_and_verbatim_hashes(self):
        self.assertEqual(self.lock["owner"], "Amazon Web Services")
        self.assertEqual(self.lock["release"], "v3.34.0")
        self.assertEqual(self.lock["commit"], "376757161b002f0c2f5d19d5cdf9def5b8c45704")
        self.assertEqual(self.lock["license_expression"], "MIT-0")
        self.assertEqual(len(self.lock["files"]), 5)
        for item in self.lock["files"]:
            payload = (ROOT / item["path"]).read_bytes()
            self.assertEqual(hashlib.sha256(payload).hexdigest(), item["sha256"], item["path"])

    def test_handler_is_valid_python_and_composes_both_official_utilities(self):
        text = (ROOT / "upstream/integrate_idempotency_with_batch_processor.py").read_text(encoding="utf-8")
        ast.parse(text)
        for required in (
            "BatchProcessor(event_type=EventType.SQS)",
            'IdempotencyConfig(event_key_jmespath="messageId")',
            '@idempotent_function(data_keyword_argument="record"',
            "config.register_lambda_context(context)",
            "process_partial_response(",
        ):
            self.assertIn(required, text)

    def test_official_templates_preserve_idempotency_and_partial_batch_controls(self):
        idempotency = (ROOT / "upstream/idempotency_sam.yaml").read_text(encoding="utf-8")
        batch = (ROOT / "upstream/sqs_batch_processing.yaml").read_text(encoding="utf-8")
        for required in (
            "TimeToLiveSpecification:",
            "BillingMode: PAY_PER_REQUEST",
            "dynamodb:PutItem",
            "dynamodb:GetItem",
            "dynamodb:UpdateItem",
            "dynamodb:DeleteItem",
            "IDEMPOTENCY_TABLE: !Ref IdempotencyTable",
        ):
            self.assertIn(required, idempotency)
        for required in (
            "FunctionResponseTypes:",
            "- ReportBatchItemFailures",
            "SqsManagedSseEnabled: true",
            "VisibilityTimeout: 30 # Fn timeout * 6",
            "deadLetterTargetArn: !GetAtt SampleDLQ.Arn",
            "maxReceiveCount: 2",
        ):
            self.assertIn(required, batch)

    def test_pack_does_not_claim_missing_business_or_recovery_layers(self):
        readme = (ROOT / "README.md").read_text(encoding="utf-8")
        self.assertIn("AWS publishes them separately", readme)
        self.assertIn("replace the demonstration `record_handler` effect", readme)
        self.assertIn("do not provide S3 sequencer ordering", readme)
        blockers = set(self.lock["production_blockers"])
        self.assertIn("DLQ replay and reconciliation are not implemented", blockers)
        self.assertIn("business idempotency key and effect are project-specific", blockers)


if __name__ == "__main__":
    unittest.main()
````

### FILE: `aws_powertools_idempotent_sqs_batch/upstream/LICENSE`
```yaml
block_id: "AWS-POWERTOOLS-IDEMPOTENT-SQS-BATCH:aws-license:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/aws-powertools/powertools-lambda-python/376757161b002f0c2f5d19d5cdf9def5b8c45704/LICENSE"
license: "MIT-0"
sha256: "bcc1d2d9e0609b03d420c3dbecb438a7b89bfc53d54176ab56301ab1f96b4f0a"
variables: []
secrets_allowed: false
```
````text
MIT No Attribution

Copyright 2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
````

### FILE: `aws_powertools_idempotent_sqs_batch/upstream/integrate_idempotency_with_batch_processor.py`
```yaml
block_id: "AWS-POWERTOOLS-IDEMPOTENT-SQS-BATCH:aws-handler:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/aws-powertools/powertools-lambda-python/376757161b002f0c2f5d19d5cdf9def5b8c45704/examples/idempotency/src/integrate_idempotency_with_batch_processor.py"
license: "MIT-0"
sha256: "d364d7b6ac242568e7c25bf4292da1576b26c0129ecf431fe4912b9191a2df29"
variables: []
secrets_allowed: false
```
````python
import os
from typing import Any, Dict

from aws_lambda_powertools.utilities.batch import BatchProcessor, EventType, process_partial_response
from aws_lambda_powertools.utilities.data_classes.sqs_event import SQSRecord
from aws_lambda_powertools.utilities.idempotency import (
    DynamoDBPersistenceLayer,
    IdempotencyConfig,
    idempotent_function,
)
from aws_lambda_powertools.utilities.typing import LambdaContext

processor = BatchProcessor(event_type=EventType.SQS)

table = os.getenv("IDEMPOTENCY_TABLE", "")
dynamodb = DynamoDBPersistenceLayer(table_name=table)
config = IdempotencyConfig(event_key_jmespath="messageId")


@idempotent_function(data_keyword_argument="record", config=config, persistence_store=dynamodb)
def record_handler(record: SQSRecord):
    return {"message": record.body}


def lambda_handler(event: Dict[str, Any], context: LambdaContext):
    config.register_lambda_context(context)  # see Lambda timeouts section

    return process_partial_response(
        event=event,
        context=context,
        processor=processor,
        record_handler=record_handler,
    )
````

### FILE: `aws_powertools_idempotent_sqs_batch/upstream/integrate_idempotency_with_batch_processor_payload.json`
```yaml
block_id: "AWS-POWERTOOLS-IDEMPOTENT-SQS-BATCH:aws-payload:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/aws-powertools/powertools-lambda-python/376757161b002f0c2f5d19d5cdf9def5b8c45704/examples/idempotency/src/integrate_idempotency_with_batch_processor_payload.json"
license: "MIT-0"
sha256: "2a5bab34e2865b8a889033f01642fce1f6841fe092bf84ebb10ae168cf3e12f0"
variables: []
secrets_allowed: false
```
````json
{
  "Records": [
    {
      "messageId": "059f36b4-87a3-44ab-83d2-661975830a7d",
      "receiptHandle": "AQEBwJnKyrHigUMZj6rYigCgxlaS3SLy0a...",
      "body": "Test message.",
      "attributes": {
        "ApproximateReceiveCount": "1",
        "SentTimestamp": "1545082649183",
        "SenderId": "replace-to-pass-gitleak",
        "ApproximateFirstReceiveTimestamp": "1545082649185"
      },
      "messageAttributes": {
        "testAttr": {
          "stringValue": "100",
          "binaryValue": "base64Str",
          "dataType": "Number"
        }
      },
      "md5OfBody": "e4e68fb7bd0e697a0ae8f1bb342846b3",
      "eventSource": "aws:sqs",
      "eventSourceARN": "arn:aws:sqs:us-east-2:123456789012:my-queue",
      "awsRegion": "us-east-2"
    }
  ]
}
````

### FILE: `aws_powertools_idempotent_sqs_batch/upstream/idempotency_sam.yaml`
```yaml
block_id: "AWS-POWERTOOLS-IDEMPOTENT-SQS-BATCH:aws-idempotency-sam:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/aws-powertools/powertools-lambda-python/376757161b002f0c2f5d19d5cdf9def5b8c45704/examples/idempotency/templates/sam.yaml"
license: "MIT-0"
sha256: "27d5a2db8673bd0e1c45799168b68dbde4f3e2e7d025c5c3c7148eaf332292a5"
variables: []
secrets_allowed: false
```
````yaml
Transform: AWS::Serverless-2016-10-31
Resources:
  IdempotencyTable:
    Type: AWS::DynamoDB::Table
    Properties:
      AttributeDefinitions:
        - AttributeName: id
          AttributeType: S
      KeySchema:
        - AttributeName: id
          KeyType: HASH
      TimeToLiveSpecification:
        AttributeName: expiration
        Enabled: true
      BillingMode: PAY_PER_REQUEST

  HelloWorldFunction:
    Type: AWS::Serverless::Function
    Properties:
      Runtime: python3.12
      Handler: app.py
      Policies:
        - Statement:
            - Sid: AllowDynamodbReadWrite
              Effect: Allow
              Action:
                - dynamodb:PutItem
                - dynamodb:GetItem
                - dynamodb:UpdateItem
                - dynamodb:DeleteItem
              Resource: !GetAtt IdempotencyTable.Arn
      Environment:
        Variables:
          IDEMPOTENCY_TABLE: !Ref IdempotencyTable
````

### FILE: `aws_powertools_idempotent_sqs_batch/upstream/sqs_batch_processing.yaml`
```yaml
block_id: "AWS-POWERTOOLS-IDEMPOTENT-SQS-BATCH:aws-batch-sam:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/aws-powertools/powertools-lambda-python/376757161b002f0c2f5d19d5cdf9def5b8c45704/examples/batch_processing/sam/sqs_batch_processing.yaml"
license: "MIT-0"
sha256: "df8341badb7c81de8d7818846071d48554f3dbdb279b8043b00eb5316941ba0e"
variables: []
secrets_allowed: false
```
````yaml
AWSTemplateFormatVersion: "2010-09-09"
Transform: AWS::Serverless-2016-10-31
Description: partial batch response sample

Globals:
  Function:
    Timeout: 5
    MemorySize: 256
    Runtime: python3.12
    Tracing: Active
    Environment:
      Variables:
        POWERTOOLS_LOG_LEVEL: INFO
        POWERTOOLS_SERVICE_NAME: hello

Resources:
  HelloWorldFunction:
    Type: AWS::Serverless::Function
    Properties:
      Handler: app.lambda_handler
      CodeUri: hello_world
      Policies:
        - SQSPollerPolicy:
            QueueName: !GetAtt SampleQueue.QueueName
      Events:
        Batch:
          Type: SQS
          Properties:
            Queue: !GetAtt SampleQueue.Arn
            FunctionResponseTypes:
              - ReportBatchItemFailures

  SampleDLQ:
    Type: AWS::SQS::Queue

  SampleQueue:
    Type: AWS::SQS::Queue
    Properties:
      VisibilityTimeout: 30 # Fn timeout * 6
      SqsManagedSseEnabled: true
      RedrivePolicy:
        maxReceiveCount: 2
        deadLetterTargetArn: !GetAtt SampleDLQ.Arn
````

## 6. Configuration surface

| Variable/decision | Type | Safe default | Validation | Secret | Mutability/effect |
|---|---|---|---|---|---|
| `IDEMPOTENCY_TABLE` | DynamoDB table name/ref | none; required | non-empty and exact deployed table | no | deployment; persistence boundary |
| business idempotency key | schema/JMESPath | none | stable, tenant-bound, immutable operation identity | no | code/schema migration |
| Lambda timeout | seconds | official example 5 | queue visibility at least 6x and workload proven | no | deployment |
| batch/concurrency | integers | provider/project decision | quotas, downstream capacity and load evidence | no | deployment/runtime |
| max receive count | integer | official example 2 | incident/retry budget | no | deployment |
| business effect | project handler | none | idempotent/conditional effect plus tests | possibly | code |

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| `aws-lambda-powertools` | `3.34.0`, wheel SHA-256 `ab1354…d679c` | idempotency + batch | MIT-0 | runtime | PyPI/AWS release SLSA |
| Python | `3.12` in upstream SAM | Lambda runtime | PSF | runtime | AWS template |
| AWS SAM | project-pinned compatible version | synth/deploy templates | Apache-2.0 | build/deploy | AWS |
| DynamoDB/SQS/Lambda | account/region selected by blueprint | managed runtime | service terms | runtime | AWS |

## 8. Apply order

1. Materializar este pack en destino vacío y ejecutar sus cuatro contratos offline.
2. Materializar/adquirir `OFFICIAL-UPSTREAM-ACQUISITION-CORE` con la fuente `aws-powertools-python-3.34.0` y verificar wheel/SLSA/licencia.
3. Fijar identidad de operación, tenant, efecto y expiración en el blueprint/ADR.
4. Crear una composición de proyecto —declarada como propia— que integre ambos templates sin modificar los archivos verbatim.
5. Ejecutar synth/lint/SCA/unit/integration y luego target-account E2E, fallo parcial, duplicado, timeout, DLQ, replay y reconciliación.
6. Rollback: retirar la composición nueva y sus recursos según el plan del proyecto; no borrar tablas/colas con estado sin preservación y autorización.

## 9. Verification

```text
python -m unittest -v aws_powertools_idempotent_sqs_batch/test_contracts.py
expected: 4 tests, OK
```

El proyecto añade `sam validate`, `sam build`, tests con DynamoDB/SQS reales o emulador admitido, duplicados concurrentes, fallo de un registro dentro del batch, expiración/reclaim, DLQ/redrive/reconciliation, IAM negativos, SCA/SBOM/notices, carga y recuperación. Ninguno se hereda del gate offline.

## 10. Reconstruction evidence

- entorno limpio: temporal V75, Windows 11, PowerShell 7, CPython 3.14.4 para contratos stdlib;
- fuente: archive exacto 9.650.888 bytes/SHA-256 `b7b8b059…c4742`;
- cinco archivos AWS verificaron byte/hash contra el lock;
- `python -m unittest -v test_contracts.py`: 4/4 PASS;
- materializador y gates globales deben repetirse antes de promover el snapshot V75;
- divergencias: ninguna en bytes AWS; README/lock/tests son packaging propio claramente atribuido;
- fecha: 2026-08-28; agente: Codex.
