# AWS Lambda Durable Execution Component

## 1. Metadata

```yaml
pack_id: "AWS-LAMBDA-DURABLE-EXECUTION-COMPONENT"
pack_version: "0.1.0"
status:
  authority: ELITE_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa byte-verbatim ejemplos oficiales AWS Durable Execution SDK Python 1.7.0 para retry, replay logging y pasos at-most-once, con wheel/runtime hash-locked y contratos offline."
stacks: ["Python 3.14", "AWS Lambda Durable Functions", "AWS Durable Execution SDK Python 1.7.0"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "AWS-POWERTOOLS-IDEMPOTENT-SQS-BATCH-COMPONENT 0.1.x"]
incompatible_with: ["event-source mapping tratado como inicio idempotente", "testing wheel 1.2.1 atribuido al árbol exacto del tag", "efecto empresarial no idempotente", "producción sin cuenta, región, IAM, replay y reconciliación"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0"
upstream_sources: ["https://github.com/aws/aws-durable-execution-sdk-python/tree/075b65aacb80de8bb1507e3df2e53cc90cb3b874"]
verified_at: "2026-08-28"
```

## 2. Applicability

Use este pack cuando un proyecto en AWS necesite código oficial actual para checkpoints, retry durable, espera/reanudación y logging consciente de replay. La selección entrega ejemplos y tests publicados por AWS, no una reimplementación local del motor.

No lo use como worker S3/SQS completo ni como garantía exactly-once. AWS documenta que un event-source mapping no hace idempotente el inicio de una durable execution; el proyecto debe usar un dispatcher con execution name estable o idempotencia Powertools, y cada efecto empresarial debe tolerar repetición.

## 3. Architecture contract

Los ocho archivos bajo `upstream/` son bytes `VERBATIM` del tag compuesto `sdk-v1.7.0,otel-v0.3.0`, commit verificado `075b65aacb80de8bb1507e3df2e53cc90cb3b874`, archive 889.895 bytes/SHA-256 `d3e42423…02d11`, licencia Apache-2.0 y NOTICE preservado. Los cuatro archivos restantes son packaging `AUTHORED` claramente separado.

El core wheel 1.7.0 queda fijado a 109.716 bytes/SHA-256 `1482e1f4…72a0`. El lock runtime de ocho wheels es configuración Elite, no un lock publicado por AWS. El testing wheel 1.2.1 no se incorpora: su árbol difiere del tag y no exporta `OperationPaginatorState`; la reproducción de tests usa source exacto o espera un release alineado.

Límites: no aporta parser de S3, version/eTag/sequencer, owner fencing, definición del efecto empresarial, dispatcher, DLQ/redrive automático, reconciliación de recursos externos, IAM/IaC ni cuenta cloud probada. `AT_MOST_ONCE_PER_RETRY` evita reejecución de un step interrumpido dentro de su semántica publicada; no transforma un efecto externo arbitrario en exactly-once. Rollback conserva historiales/receipts y retira sólo la nueva composición autorizada.

## 4. Exact file manifest

```text
CREATE aws_lambda_durable_execution/README.md
CREATE aws_lambda_durable_execution/source-lock.json
CREATE aws_lambda_durable_execution/requirements-runtime.lock
CREATE aws_lambda_durable_execution/test_contracts.py
CREATE aws_lambda_durable_execution/upstream/LICENSE
CREATE aws_lambda_durable_execution/upstream/NOTICE
CREATE aws_lambda_durable_execution/upstream/step_with_retry.py
CREATE aws_lambda_durable_execution/upstream/test_step_with_retry.py
CREATE aws_lambda_durable_execution/upstream/step_semantics_at_most_once.py
CREATE aws_lambda_durable_execution/upstream/test_step_semantics_at_most_once.py
CREATE aws_lambda_durable_execution/upstream/replay_logging.py
CREATE aws_lambda_durable_execution/upstream/test_replay_logging.py
```

## 5. Materialization blocks

### FILE: `aws_lambda_durable_execution/README.md`
```yaml
block_id: "AWS-LAMBDA-DURABLE-EXECUTION:readme-md:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite packaging constrained by exact AWS upstream"
license: "LicenseRef-Workspace-Owner"
sha256: "e13313f788aee36d7da4f28d7f7f57c61c61b0ca1283bd68fa67943fce64265f"
variables: []
secrets_allowed: false
```
````markdown
# AWS Lambda Durable Execution component

This directory packages byte-verbatim AWS Durable Execution SDK for Python examples from the verified composite tag `sdk-v1.7.0,otel-v0.3.0`, commit `075b65aacb80de8bb1507e3df2e53cc90cb3b874`.

The upstream examples demonstrate retryable durable steps, `AT_MOST_ONCE_PER_RETRY` step semantics, checkpoint replay and replay-aware logging. They are building blocks, not a complete S3/SQS worker. AWS documents that event-source mappings can start duplicate durable executions; the project must use a dispatcher with a stable execution name or Powertools idempotency, and its business steps must remain idempotent.

Install the runtime with `python -m pip install --require-hashes -r requirements-runtime.lock`. The lock is Elite-authored packaging around AWS's exact core wheel; it is not an AWS-published requirements file. Preserve Apache-2.0 and NOTICE.

Do not install the published testing wheel 1.2.1 as evidence for this tag. The tag source and that wheel differ: the source test suite imports `OperationPaginatorState`, which the wheel lacks. Acquire the exact commit source for upstream test reproduction or wait for a testing release aligned to the selected tree.

Run `python -m unittest -v test_contracts.py` for the offline packaging gate. Production still requires an AWS account/region where Durable Functions is available, execution-name idempotency, IAM, limits, S3 version/eTag/sequencer ordering, business-effect idempotency, DLQ/redrive, reconciliation, observability, load, recovery and rollback.
````

### FILE: `aws_lambda_durable_execution/source-lock.json`
```yaml
block_id: "AWS-LAMBDA-DURABLE-EXECUTION:source-lock-json:v1"
operation: CREATE
provenance: AUTHORED
source: "exact inventory derived from AWS release and PyPI"
license: "LicenseRef-Workspace-Owner"
sha256: "803f6f42b56751ffe407686daacdc2ecae7e458ce5eab00a927f3226429ee605"
variables: []
secrets_allowed: false
```
````json
{
  "source_id": "aws-lambda-durable-execution-python-1.7.0",
  "owner": "Amazon Web Services",
  "repository": "aws/aws-durable-execution-sdk-python",
  "release": "sdk-v1.7.0,otel-v0.3.0",
  "commit": "075b65aacb80de8bb1507e3df2e53cc90cb3b874",
  "commit_signature_verified": true,
  "archive_bytes": 889895,
  "archive_sha256": "d3e42423e61e258191fdd078a29d1fd58cfc55380b25a357f879b73665a02d11",
  "core_wheel": {
    "version": "1.7.0",
    "bytes": 109716,
    "sha256": "1482e1f439e36deb1fcf9a086499b3788c3facf3ac755e50376a739081cc72a0"
  },
  "core_sdist": {
    "version": "1.7.0",
    "bytes": 240798,
    "sha256": "6ce773cfba243df98ec35f375a30b98cc4779a9f921d8c6d1180773811bf57f0"
  },
  "testing_wheel_condition": {
    "version": "1.2.1",
    "bytes": 101233,
    "sha256": "c6d3de645e55aeb277bdb685a3ebd016fd155b395b984f05cf8dc7fb6100d8f2",
    "classification": "INCOMPATIBLE_WITH_SELECTED_TAG_TEST_TREE",
    "source_execution_py_sha256": "abc2785722246de1c386a91af2d3e9d9a0f9692f8f03d91234a1f2688b697ee9",
    "wheel_execution_py_sha256": "4269ee0b876fe36a5b3d8ddb69b4f75613a775405e29c27ababbcca3a82b8917"
  },
  "license_expression": "Apache-2.0",
  "classification": "PINNED_COMPONENT_CONDITIONED",
  "files": [
    {"path":"upstream/LICENSE","source_path":"LICENSE","sha256":"09e8a9bcec8067104652c168685ab0931e7868f9c8284b66f5ae6edae5f1130b"},
    {"path":"upstream/NOTICE","source_path":"NOTICE","sha256":"d4290ed64c2edd0fce1d84e3f9dfb2881240fe534def76b8cd29ed6af683e287"},
    {"path":"upstream/step_with_retry.py","source_path":"packages/aws-durable-execution-sdk-python-examples/src/step/step_with_retry.py","sha256":"a809441f9174be84b54df17cec8af5028e3fb3c9e09cb3c0dde1116788e48987"},
    {"path":"upstream/test_step_with_retry.py","source_path":"packages/aws-durable-execution-sdk-python-examples/test/step/test_step_with_retry.py","sha256":"b3c46f64fc7619a89adbd302efca384412823ab7cc8e4a33ef3e9d33184b57eb"},
    {"path":"upstream/step_semantics_at_most_once.py","source_path":"packages/aws-durable-execution-sdk-python-examples/src/step/step_semantics_at_most_once.py","sha256":"071a0aeae4e4eeb294da389a7aa28bf71213fb513428de018e6b710c146040d5"},
    {"path":"upstream/test_step_semantics_at_most_once.py","source_path":"packages/aws-durable-execution-sdk-python-examples/test/step/test_step_semantics_at_most_once.py","sha256":"ae9513b146347c0c7cbb189bdd12531eecd5e3aa8669c4f17e9983997c0388cf"},
    {"path":"upstream/replay_logging.py","source_path":"packages/aws-durable-execution-sdk-python-examples/src/logger_example/replay_logging.py","sha256":"4da5c9529859104b91c5210578704f91505e7e0bf653ce412983a3371e02c06a"},
    {"path":"upstream/test_replay_logging.py","source_path":"packages/aws-durable-execution-sdk-python-examples/test/logger_example/test_replay_logging.py","sha256":"5f47faa486407aaa72ef289c8a0202f13ad4409406b3ca1a2496862d5cd39a6e"}
  ]
}
````

### FILE: `aws_lambda_durable_execution/requirements-runtime.lock`
```yaml
block_id: "AWS-LAMBDA-DURABLE-EXECUTION:requirements-runtime-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "hash-locked project runtime resolution around AWS core wheel"
license: "LicenseRef-Workspace-Owner"
sha256: "e38a6df20f235238a215b196f50e215dce9c4d82a265d60854193815659f8ffb"
variables: []
secrets_allowed: false
```
````text
aws-durable-execution-sdk-python==1.7.0 --hash=sha256:1482e1f439e36deb1fcf9a086499b3788c3facf3ac755e50376a739081cc72a0
boto3==1.43.82 --hash=sha256:65319e8bba6e30f74a6e2727a5688725222da2ab71c6069bc484ce5bfd101c73
botocore==1.43.82 --hash=sha256:97b3e89061decc91e7745d726dae595cbe2053894611c49ea91d6fdcb2ecc36b
jmespath==1.1.0 --hash=sha256:a5663118de4908c91729bea0acadca56526eb2698e83de10cd116ae0f4e97c64
python-dateutil==2.9.0.post0 --hash=sha256:a8b2bc7bffae282281c8140a97d3aa9c14da0b136dfe83f850eea9a5f7470427
s3transfer==0.19.2 --hash=sha256:d8168eccca828cbb2cd573675333f3bddd254313a9c42494b84c76b539e8ba25
six==1.17.0 --hash=sha256:4721f391ed90541fddacab5acf947aa0d3dc7d27b2e1e8eda2be8970586c3274
urllib3==2.7.0 --hash=sha256:9fb4c81ebbb1ce9531cce37674bbc6f1360472bc18ca9a553ede278ef7276897
````

### FILE: `aws_lambda_durable_execution/test_contracts.py`
```yaml
block_id: "AWS-LAMBDA-DURABLE-EXECUTION:test-contracts-py:v1"
operation: CREATE
provenance: AUTHORED
source: "offline packaging contracts"
license: "LicenseRef-Workspace-Owner"
sha256: "4829b590402f1f4e569ffe7bf33e064f0eee6c3d1b4d28613f277a97bef5a7ca"
variables: []
secrets_allowed: false
```
````python
import hashlib
import json
import re
import unittest
from pathlib import Path


ROOT = Path(__file__).resolve().parent


class DurableComponentContracts(unittest.TestCase):
    def test_verbatim_files_match_source_lock(self):
        lock = json.loads((ROOT / "source-lock.json").read_text(encoding="utf-8"))
        self.assertEqual(lock["commit"], "075b65aacb80de8bb1507e3df2e53cc90cb3b874")
        self.assertEqual(lock["core_wheel"]["sha256"], "1482e1f439e36deb1fcf9a086499b3788c3facf3ac755e50376a739081cc72a0")
        for item in lock["files"]:
            payload = (ROOT / item["path"]).read_bytes()
            self.assertEqual(hashlib.sha256(payload).hexdigest(), item["sha256"], item["path"])

    def test_runtime_lock_is_exact_and_hash_required(self):
        lines = [line for line in (ROOT / "requirements-runtime.lock").read_text(encoding="utf-8").splitlines() if line]
        self.assertEqual(len(lines), 8)
        self.assertTrue(all(re.fullmatch(r"[a-z0-9-]+==[^ ]+ --hash=sha256:[0-9a-f]{64}", line) for line in lines))
        self.assertTrue(lines[0].startswith("aws-durable-execution-sdk-python==1.7.0 "))

    def test_retry_and_at_most_once_examples_remain_exact_contracts(self):
        retry = (ROOT / "upstream/step_with_retry.py").read_text(encoding="utf-8")
        at_most_once = (ROOT / "upstream/step_semantics_at_most_once.py").read_text(encoding="utf-8")
        self.assertIn("RetryStrategyConfig", retry)
        self.assertIn("max_attempts=3", retry)
        self.assertIn("retryable_error_types=[RuntimeError]", retry)
        self.assertIn("StepSemantics.AT_MOST_ONCE_PER_RETRY", at_most_once)

    def test_replay_example_and_upstream_tests_are_present(self):
        replay = (ROOT / "upstream/replay_logging.py").read_text(encoding="utf-8")
        replay_test = (ROOT / "upstream/test_replay_logging.py").read_text(encoding="utf-8")
        self.assertIn("context.wait", replay)
        self.assertIn("context.is_replaying()", replay)
        self.assertIn("messages.count(\"Workflow started (before wait)\") == 1", replay_test)
        self.assertIn("messages.count(\"Workflow completed\") == 1", replay_test)


if __name__ == "__main__":
    unittest.main()
````

### FILE: `aws_lambda_durable_execution/upstream/LICENSE`
```yaml
block_id: "AWS-LAMBDA-DURABLE-EXECUTION:upstream-license:v1"
operation: CREATE
provenance: VERBATIM
source: "aws/aws-durable-execution-sdk-python 075b65a LICENSE"
license: "Apache-2.0"
sha256: "09e8a9bcec8067104652c168685ab0931e7868f9c8284b66f5ae6edae5f1130b"
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
````

### FILE: `aws_lambda_durable_execution/upstream/NOTICE`
```yaml
block_id: "AWS-LAMBDA-DURABLE-EXECUTION:upstream-notice:v1"
operation: CREATE
provenance: VERBATIM
source: "aws/aws-durable-execution-sdk-python 075b65a NOTICE"
license: "Apache-2.0"
sha256: "d4290ed64c2edd0fce1d84e3f9dfb2881240fe534def76b8cd29ed6af683e287"
variables: []
secrets_allowed: false
```
````text
Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
````

### FILE: `aws_lambda_durable_execution/upstream/step_with_retry.py`
```yaml
block_id: "AWS-LAMBDA-DURABLE-EXECUTION:upstream-step-with-retry-py:v1"
operation: CREATE
provenance: VERBATIM
source: "AWS Durable Execution Python example"
license: "Apache-2.0"
sha256: "a809441f9174be84b54df17cec8af5028e3fb3c9e09cb3c0dde1116788e48987"
variables: []
secrets_allowed: false
```
````python
from itertools import count
from typing import Any

from aws_durable_execution_sdk_python.config import StepConfig
from aws_durable_execution_sdk_python.context import (
    DurableContext,
    StepContext,
    durable_step,
)
from aws_durable_execution_sdk_python.execution import durable_execution
from aws_durable_execution_sdk_python.retries import (
    RetryStrategyConfig,
    create_retry_strategy,
)


# Counter for deterministic behavior across retries
_attempts = count(1)  # starts from 1


@durable_step
def unreliable_operation(
    _step_context: StepContext,
) -> str:
    # Use counter for deterministic behavior
    # Will fail on first attempt, succeed on second
    attempt = next(_attempts)
    if attempt < 2:
        msg = f"Attempt {attempt} failed"
        raise RuntimeError(msg)
    return "Operation succeeded"


@durable_execution
def handler(_event: Any, context: DurableContext) -> str:
    retry_config = RetryStrategyConfig(
        max_attempts=3,
        retryable_error_types=[RuntimeError],
    )

    result: str = context.step(
        unreliable_operation(),
        config=StepConfig(create_retry_strategy(retry_config)),
    )

    return result
````

### FILE: `aws_lambda_durable_execution/upstream/test_step_with_retry.py`
```yaml
block_id: "AWS-LAMBDA-DURABLE-EXECUTION:upstream-test-step-with-retry-py:v1"
operation: CREATE
provenance: VERBATIM
source: "AWS Durable Execution Python example test"
license: "Apache-2.0"
sha256: "b3c46f64fc7619a89adbd302efca384412823ab7cc8e4a33ef3e9d33184b57eb"
variables: []
secrets_allowed: false
```
````python
"""Tests for step_with_retry example."""

import pytest
from aws_durable_execution_sdk_python.execution import InvocationStatus
from aws_durable_execution_sdk_python.lambda_service import OperationType
from src.step import step_with_retry
from test.conftest import deserialize_operation_payload


@pytest.mark.example
@pytest.mark.durable_execution(
    handler=step_with_retry.handler,
    lambda_function_name="step with retry",
)
def test_step_with_retry(durable_runner):
    """Test step with retry configuration.

    With counter-based deterministic behavior:
    - Attempt 1: counter = 1 < 2 → raises RuntimeError ❌
    - Attempt 2: counter = 2 >= 2 → succeeds ✓

    The function deterministically fails once then succeeds on the second attempt.
    """
    with durable_runner:
        result = durable_runner.run(input="test", timeout=30)

    # With counter-based deterministic behavior, succeeds on attempt 2
    assert result.status is InvocationStatus.SUCCEEDED
    assert deserialize_operation_payload(result.result) == "Operation succeeded"

    # Verify step operation exists with retry details
    step_ops = [
        op for op in result.operations if op.operation_type == OperationType.STEP
    ]
    assert len(step_ops) == 1

    # The step should have succeeded on attempt 2 (after 1 failure)
    # Attempt numbering: 1 (initial attempt), 2 (first retry)
    step_op = step_ops[0]
    assert step_op.attempt == 2  # Succeeded on first retry (1-indexed: 2=first retry)
````

### FILE: `aws_lambda_durable_execution/upstream/step_semantics_at_most_once.py`
```yaml
block_id: "AWS-LAMBDA-DURABLE-EXECUTION:upstream-step-semantics-at-most-once-py:v1"
operation: CREATE
provenance: VERBATIM
source: "AWS Durable Execution Python example"
license: "Apache-2.0"
sha256: "071a0aeae4e4eeb294da389a7aa28bf71213fb513428de018e6b710c146040d5"
variables: []
secrets_allowed: false
```
````python
from typing import Any

from aws_durable_execution_sdk_python.config import StepConfig, StepSemantics
from aws_durable_execution_sdk_python.context import DurableContext
from aws_durable_execution_sdk_python.execution import durable_execution


@durable_execution
def handler(_event: Any, context: DurableContext) -> str:
    # Step with AT_MOST_ONCE_PER_RETRY semantics
    config = StepConfig(step_semantics=StepSemantics.AT_MOST_ONCE_PER_RETRY)

    result = context.step(
        lambda _: "AT_MOST_ONCE_PER_RETRY semantics",
        name="at_most_once_step",
        config=config,
    )
    return f"Result: {result}"
````

### FILE: `aws_lambda_durable_execution/upstream/test_step_semantics_at_most_once.py`
```yaml
block_id: "AWS-LAMBDA-DURABLE-EXECUTION:upstream-test-step-semantics-at-most-once-py:v1"
operation: CREATE
provenance: VERBATIM
source: "AWS Durable Execution Python example test"
license: "Apache-2.0"
sha256: "ae9513b146347c0c7cbb189bdd12531eecd5e3aa8669c4f17e9983997c0388cf"
variables: []
secrets_allowed: false
```
````python
"""Tests for step_semantics_at_most_once example."""

import pytest
from aws_durable_execution_sdk_python.execution import InvocationStatus
from aws_durable_execution_sdk_python.lambda_service import OperationType

from src.step import step_semantics_at_most_once
from test.conftest import deserialize_operation_payload


@pytest.mark.example
@pytest.mark.durable_execution(
    handler=step_semantics_at_most_once.handler,
    lambda_function_name="step semantics at most once",
)
def test_step_semantics_at_most_once(durable_runner):
    """Test step with at-most-once semantics."""
    with durable_runner:
        result = durable_runner.run(input="test", timeout=10)

    assert result.status is InvocationStatus.SUCCEEDED
    assert (
        deserialize_operation_payload(result.result)
        == "Result: AT_MOST_ONCE_PER_RETRY semantics"
    )

    # Verify step operation exists with correct name
    step_ops = [
        op for op in result.operations if op.operation_type == OperationType.STEP
    ]
    assert len(step_ops) == 1
    assert step_ops[0].name == "at_most_once_step"
````

### FILE: `aws_lambda_durable_execution/upstream/replay_logging.py`
```yaml
block_id: "AWS-LAMBDA-DURABLE-EXECUTION:upstream-replay-logging-py:v1"
operation: CREATE
provenance: VERBATIM
source: "AWS Durable Execution Python example"
license: "Apache-2.0"
sha256: "4da5c9529859104b91c5210578704f91505e7e0bf653ce412983a3371e02c06a"
variables: []
secrets_allowed: false
```
````python
"""Example demonstrating replay-aware logging across a wait boundary."""

from typing import Any

from aws_durable_execution_sdk_python.config import Duration
from aws_durable_execution_sdk_python.context import (
    DurableContext,
    StepContext,
    durable_step,
    durable_with_child_context,
)
from aws_durable_execution_sdk_python.execution import durable_execution


@durable_step
def prepare(step_context: StepContext, item: str) -> str:
    """A step that runs before the wait.

    Its log is emitted on the first invocation. On replay this step is not
    re-executed (it returns its checkpointed result), so this log does not
    repeat.
    """
    step_context.logger.info("Preparing item", extra={"item": item})
    return f"prepared:{item}"


@durable_step
def finalize(step_context: StepContext, prepared: str) -> str:
    """A step that runs after the wait (new work on the replay invocation)."""
    step_context.logger.info("Finalizing item", extra={"prepared": prepared})
    return f"done:{prepared}"


@durable_with_child_context
def audit(child_ctx: DurableContext, prepared: str) -> str:
    """Child context with its own logger and its own replay status."""
    child_ctx.logger.info(
        "Auditing in child context (before child wait)",
        extra={"prepared": prepared, "child_is_replaying": child_ctx.is_replaying()},
    )

    # The child's own replay boundary.
    child_ctx.wait(duration=Duration.from_seconds(5), name="audit_cooldown")

    # After the child's wait: emitted as new work on the child's replay.
    child_ctx.logger.info(
        "Resumed in child context (after child wait)",
        extra={"child_is_replaying": child_ctx.is_replaying()},
    )

    return child_ctx.step(lambda _: f"audited:{prepared}", name="record_audit")


@durable_execution
def handler(event: Any, context: DurableContext) -> dict[str, Any]:
    """Handler demonstrating replay-aware logging across a wait."""
    item: str = event.get("item", "widget") if isinstance(event, dict) else "widget"

    # --- Before the wait ---
    # On the replay invocation these lines are de-duplicated by the replay-aware
    # logger because the context is still replaying when it reaches them.
    context.logger.info(
        "Workflow started (before wait)",
        extra={"item": item, "is_replaying": context.is_replaying()},
    )

    prepared: str = context.step(prepare(item), name="prepare")

    context.logger.info(
        "Prepared, about to wait",
        extra={"prepared": prepared, "is_replaying": context.is_replaying()},
    )

    # --- The replay boundary ---
    # The wait suspends the execution. When it resumes, the handler replays from
    # the top; everything above is de-duplicated, and everything below is new.
    context.wait(duration=Duration.from_seconds(5), name="cooldown")

    # --- After the wait ---
    # These logs are emitted on the replay invocation because the context has
    # crossed its replay boundary and is no longer replaying.
    context.logger.info(
        "Resumed after wait",
        extra={"is_replaying": context.is_replaying()},
    )

    audited: str = context.run_in_child_context(audit(prepared), name="audit")

    result: str = context.step(finalize(audited), name="finalize")

    context.logger.info("Workflow completed", extra={"result": result})

    return {"result": result, "item": item}
````

### FILE: `aws_lambda_durable_execution/upstream/test_replay_logging.py`
```yaml
block_id: "AWS-LAMBDA-DURABLE-EXECUTION:upstream-test-replay-logging-py:v1"
operation: CREATE
provenance: VERBATIM
source: "AWS Durable Execution Python example test"
license: "Apache-2.0"
sha256: "5f47faa486407aaa72ef289c8a0202f13ad4409406b3ca1a2496862d5cd39a6e"
variables: []
secrets_allowed: false
```
````python
"""Tests for the replay_logging example.

Most assertions here verify the workflow runs end-to-end across the
wait/replay boundary and produces the expected operations and result. One test
additionally asserts on the replay-aware logger: messages emitted before the
wait are de-duplicated on replay, so each appears exactly once despite the
handler replaying from the top after the wait resumes.
"""

import logging

import pytest

from aws_durable_execution_sdk_python.execution import InvocationStatus
from aws_durable_execution_sdk_python.lambda_service import OperationType
from src.logger_example import replay_logging
from test.conftest import deserialize_operation_payload


@pytest.mark.example
@pytest.mark.durable_execution(
    handler=replay_logging.handler,
    lambda_function_name="Replay Logging",
)
def test_replay_logging(durable_runner):
    """Test the replay-aware logging example runs across the wait boundary."""
    with durable_runner:
        result = durable_runner.run(input={"item": "widget"}, timeout=30)

    assert result.status is InvocationStatus.SUCCEEDED
    assert deserialize_operation_payload(result.result) == {
        "result": "done:audited:prepared:widget",
        "item": "widget",
    }

    # Two wait operations force suspend/replay cycles: one in the parent context
    # and one inside the child (audit) context. This exercises per-context replay
    # status in different contexts.
    wait_ops = [
        op for op in result.operations if op.operation_type == OperationType.WAIT
    ]
    assert len(wait_ops) >= 1

    # Steps before (prepare) and after (finalize) the wait both ran. The child
    # context's record_audit step is nested inside the CONTEXT operation.
    step_ops = [
        op for op in result.operations if op.operation_type == OperationType.STEP
    ]
    assert len(step_ops) >= 2

    # The audit child context produces a CONTEXT operation.
    context_ops = [
        op for op in result.operations if op.operation_type.value == "CONTEXT"
    ]
    assert len(context_ops) >= 1


@pytest.mark.example
@pytest.mark.durable_execution(
    handler=replay_logging.handler,
    lambda_function_name="Replay Logging",
)
def test_replay_logging_dedupes_logs_across_wait(durable_runner, caplog):
    """Verify the replay-aware logger de-duplicates messages across the wait.

    The handler logs before the wait, suspends, then replays from the top when
    the wait resumes. The replay-aware logger suppresses logs while the context
    is replaying, so a message emitted before the wait must appear exactly once
    even though the code that produces it runs again on replay. Messages emitted
    after the wait are new work and also appear exactly once.

    This only holds in local mode, where every invocation runs in-process and is
    captured by a single ``caplog``; cloud mode spreads invocations across
    separate Lambda executions, so the test is skipped there.
    """
    if durable_runner.mode != "local":
        pytest.skip("Log capture is only available in local (in-process) mode")

    with caplog.at_level(logging.INFO):
        with durable_runner:
            result = durable_runner.run(input={"item": "widget"}, timeout=30)

    assert result.status is InvocationStatus.SUCCEEDED

    messages = [record.getMessage() for record in caplog.records]

    # Emitted before the wait: the handler replays past this line, but the
    # replay-aware logger suppresses the duplicate, so it appears exactly once.
    assert messages.count("Workflow started (before wait)") == 1
    assert messages.count("Prepared, about to wait") == 1

    # Emitted after the wait as new work: also exactly once.
    assert messages.count("Resumed after wait") == 1
    assert messages.count("Workflow completed") == 1
````

## 6. Configuration surface

| Variable/decisión | Tipo | Default seguro | Validación | Secret | Efecto |
|---|---|---|---|---|---|
| execution name | identidad estable | ninguno | único por cuenta/región y hash/payload compatible | no | inicio idempotente |
| step semantics | enum AWS | at-least-once publicado | seleccionar at-most-once sólo con análisis de interrupción | no | replay/retry |
| retry strategy | intentos/backoff/errores | ninguno empresarial | errores transitorios explícitos, límites y timeout | no | disponibilidad/costo |
| durable retention/timeout | días/segundos | ninguno | región, cuotas, datos y RTO/RPO aprobados | no | operación/costo |
| business effect key | identidad tenant-bound | ninguno | durable, estable y validado contra payload | no | consistencia |
| AWS account/region/IAM | referencias | ninguno | servicio disponible y least privilege probado | posiblemente | despliegue |

## 7. Dependency bill

| Paquete/herramienta | Pin exacto | Uso | Licencia | Runtime/build | Autoridad |
|---|---|---|---|---|---|
| `aws-durable-execution-sdk-python` | 1.7.0, wheel SHA-256 `1482e1f4…72a0` | checkpoint/replay/retry | Apache-2.0 | runtime | AWS/PyPI |
| boto3 graph | 1.43.82/1.43.82 + seis transitivas exactas | Lambda API client | Apache-2.0/MIT | runtime | PyPI lock Elite |
| Python | 3.14 target audit | runtime Lambda | PSF | runtime | AWS supported runtime |
| testing SDK source | commit `075b65a…` only for reproduction | local tests | Apache-2.0 | test | AWS GitHub |
| testing wheel 1.2.1 | rechazado para este tag | no usar como evidencia | Apache-2.0 | test | condición upstream 154 |

## 8. Apply order

1. Materializar en destino vacío y ejecutar cuatro contratos offline.
2. Instalar `requirements-runtime.lock` con `--require-hashes` y preservar LICENSE/NOTICE.
3. Seleccionar dispatcher/execution name, identidad empresarial y política de replay en el blueprint.
4. Componer con Powertools batch/idempotency sólo mediante código de proyecto atribuido, nunca como archivo AWS inexistente.
5. Ejecutar upstream source-exact, unit/integration, duplicados, interrupción, timeout, DLQ/redrive/reconciliation, IAM, SCA, carga y target-account E2E.
6. Rollback: detener nuevos inicios, preservar histories/receipts y revertir la composición/IaC autorizada.

## 9. Verification

```text
python -m unittest -v aws_lambda_durable_execution/test_contracts.py
expected: 4 tests, OK
```

Evidencia upstream: Windows source-exact 2.791 PASS/3 fallos de portabilidad/2 skips; Linux Python 3.14 source-exact 115/115 para los 112 focales más los tres nodeids Windows. OSV-Scanner 2.5.1 consultó el lock de ocho paquetes el 2026-08-28 y devolvió cero vulnerabilidades conocidas; esto no sustituye cloud E2E, un scan fresco del grafo final ni pruebas del efecto del proyecto.

## 10. Reconstruction evidence

- identidad: tag compuesto y commit AWS verificado `075b65a…`;
- archive: 889.895 bytes/SHA-256 `d3e42423…02d11`;
- artifacts: core wheel/sdist PyPI exactos; ocho archivos AWS byte/hash locked;
- ejecución: 112/112 focal Windows; integral Windows 2.791/3/2; Linux 115/115;
- SCA fechada: lock SHA-256 `e38a6df2…f8ffb`, ocho paquetes, cero vulnerabilidades conocidas; reporte 1.553 bytes/SHA-256 `70ec56bd…b99b`;
- condición: testing wheel 1.2.1 difiere del source del tag y queda excluido;
- divergencias: cero en archivos verbatim; README/lock/requirements/contratos son Elite `AUTHORED`;
- fecha: 2026-08-28; agente: Codex.
