# Microsoft Azure Content Understanding Official Analyzer Copy Sample

## 1. Metadata

```yaml
pack_id: "MICROSOFT-AZURE-CONTENT-UNDERSTANDING-OFFICIAL-ANALYZER-COPY-SAMPLE"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa byte a byte el sample Python y el test oficial Microsoft Azure AI Content Understanding 1.1.0 para crear, inspeccionar, copiar, actualizar, verificar y eliminar un analyzer dentro del mismo recurso, junto con licencia, source lock y tres contratos offline."
stacks: ["Python 3.9+", "azure-ai-contentunderstanding 1.1.0", "Azure Content Understanding API 2025-11-01"]
compatible_with: ["OFFICIAL-DOCUMENT-SDK-ARTIFACT-CORE 0.3.x", "AZURE-CONTENT-UNDERSTANDING-DOCUMENT-RUNTIME 0.1.x"]
incompatible_with: ["ejecución automática", "promoción automática a producción", "copy cross-resource", "dependencias auxiliares sin lock"]
license_expression: "LicenseRef-Workspace-Owner AND MIT"
upstream_sources: ["https://github.com/Azure/azure-sdk-for-python/tree/129a5cbb06fba44e11ffc92f6bcdca14a2f5b6eb/sdk/contentunderstanding/azure-ai-contentunderstanding", "https://pypi.org/project/azure-ai-contentunderstanding/1.1.0/", "https://learn.microsoft.com/azure/ai-services/content-understanding/concepts/prebuilt-analyzers"]
verified_at: "2026-08-28"
```

## 2. Applicability

Use este pack cuando un proyecto Azure Content Understanding deba partir del lifecycle de copy publicado por Microsoft para fijar un analyzer dentro del mismo recurso. El código upstream queda disponible inmediatamente y los contratos offline prueban identidad, sintaxis y llamadas exactas.

No lo ejecute automáticamente: crea, copia, actualiza y elimina recursos externos. Exige endpoint/identidad/región/modelos/costo, autoridad explícita sobre cada mutación, IDs target, receipt de versión y reconciliación posterior. El sample usa `prebuilt-document`; no demuestra `prebuilt-procurement`, `prebuilt-purchaseOrder`, proforma, packing list, BOL, aduana, exactitud ni promoción productiva.

## 3. Architecture contract

Los tres archivos bajo `upstream/` son bytes `VERBATIM` Microsoft MIT del commit firmado `129a5cbb…`; sample y test son byte-idénticos al sdist 1.1.0. README, lock y contrato offline son `AUTHORED` de packaging y nunca se atribuyen a Microsoft.

El sample ejecuta create → get → copy → get → update → get → delete. El test oficial usa Azure SDK test-proxy/fixtures, no es standalone. El sample absorbe excepciones de cleanup: cualquier ejecución live requiere inventario y reconciliación independientes. Rollback no es “ejecutar delete a ciegas”; debe comparar el receipt/ID target y conservar evidencia antes de eliminar.

## 4. Exact file manifest

```text
CREATE azure_content_understanding_official_analyzer_copy/README.md
CREATE azure_content_understanding_official_analyzer_copy/upstream/LICENSE.txt
CREATE azure_content_understanding_official_analyzer_copy/upstream/sample_copy_analyzer.py
CREATE azure_content_understanding_official_analyzer_copy/upstream/test_sample_copy_analyzer.py
CREATE azure_content_understanding_official_analyzer_copy/source-lock.json
CREATE azure_content_understanding_official_analyzer_copy/test_official_copy_sample.py
```

## 5. Materialization blocks

### FILE: `azure_content_understanding_official_analyzer_copy/README.md`
```yaml
block_id: "AZURE-CU-OFFICIAL-ANALYZER-COPY:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local guide constrained by exact Microsoft source"
license: "LicenseRef-Workspace-Owner"
sha256: "231bebc292b343c1fc8e10e33db68ae95b672a1c645bc898baae4fb9a3ee34dd"
variables: []
secrets_allowed: false
```
````markdown
# Microsoft Azure Content Understanding official analyzer copy sample

This directory preserves the exact Microsoft Python SDK 1.1.0 sample and its exact official recorded/live test for same-resource analyzer copy. It demonstrates create → inspect → copy → inspect → update → verify → delete with the official `ContentUnderstandingClient`.

The upstream sample creates external analyzer resources and deletes them. Do not run it automatically. Live use requires a Microsoft Foundry resource, authorized identity, supported region/models, cost approval, unique analyzer IDs, explicit create/update/delete authority and post-run reconciliation. Cleanup exceptions are intentionally swallowed by the sample and therefore require an independent inventory/reconciliation step.

The sample creates a custom analyzer based on `prebuilt-document`; it does not call `prebuilt-procurement`, `prebuilt-purchaseOrder`, proforma, packing-list, bill-of-lading or customs analyzers. The code is evidence for the copy lifecycle only. It does not prove extraction accuracy, a production promotion decision, cross-resource copy or rollback.

`upstream/` is byte-verbatim Microsoft MIT code. `source-lock.json`, this README and `test_official_copy_sample.py` are Elite packaging/verification and are not attributed to Microsoft.
````

### FILE: `azure_content_understanding_official_analyzer_copy/upstream/LICENSE.txt`
```yaml
block_id: "AZURE-CU-OFFICIAL-ANALYZER-COPY:microsoft-license:v1"
operation: CREATE
provenance: VERBATIM
source: "Azure/azure-sdk-for-python@129a5cbb repository root LICENSE"
license: "MIT"
sha256: "7c77a44a8acd9b41fdc209864a8016b3d430b5d0e09309818d5b7444336df744"
variables: []
secrets_allowed: false
```
````text
Copyright (c) Microsoft Corporation.

MIT License

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED *AS IS*, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
````

### FILE: `azure_content_understanding_official_analyzer_copy/upstream/sample_copy_analyzer.py`
```yaml
block_id: "AZURE-CU-OFFICIAL-ANALYZER-COPY:microsoft-sample:v1"
operation: CREATE
provenance: VERBATIM
source: "Azure/azure-sdk-for-python@129a5cbb SDK sample"
license: "MIT"
sha256: "4b8aa70c941c9b0b8fcd54ddac588ef0d1a6dbd4cefe5e2c21d11fd8e44bdb71"
variables: []
secrets_allowed: false
```
````python
# pylint: disable=line-too-long,useless-suppression
# coding=utf-8
# --------------------------------------------------------------------------
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License. See License.txt in the project root for license information.
# --------------------------------------------------------------------------
"""
FILE: sample_copy_analyzer.py

DESCRIPTION:
    This sample demonstrates how to copy an analyzer from source to target within the same
    Microsoft Foundry resource using the begin_copy_analyzer API. This is useful for
    creating copies of analyzers for testing, staging, or production deployment.

    About copying analyzers
    The begin_copy_analyzer API allows you to copy an analyzer within the same Azure resource:
    - Same-resource copy: Copies an analyzer from one ID to another within the same resource
    - Exact copy: The target analyzer is an exact copy of the source analyzer

    Note: For cross-resource copying (copying between different Azure resources or subscriptions),
    use the grant_copy_auth sample instead.

USAGE:
    python sample_copy_analyzer.py

    Set the environment variables with your own values before running the sample:
    1) CONTENTUNDERSTANDING_ENDPOINT - the endpoint to your Content Understanding resource.
    2) CONTENTUNDERSTANDING_KEY - your Content Understanding API key (optional if using DefaultAzureCredential).
"""

import os
import time

from dotenv import load_dotenv
from azure.ai.contentunderstanding import ContentUnderstandingClient
from azure.ai.contentunderstanding.models import (
    ContentAnalyzer,
    ContentAnalyzerConfig,
    ContentFieldSchema,
    ContentFieldDefinition,
    ContentFieldType,
    GenerationMethod,
)
from azure.core.credentials import AzureKeyCredential
from azure.identity import DefaultAzureCredential

load_dotenv()


def main() -> None:
    endpoint = os.environ["CONTENTUNDERSTANDING_ENDPOINT"]
    key = os.getenv("CONTENTUNDERSTANDING_KEY")
    credential = AzureKeyCredential(key) if key else DefaultAzureCredential()

    client = ContentUnderstandingClient(endpoint=endpoint, credential=credential)

    base_id = f"my_analyzer_{int(time.time())}"
    source_analyzer_id = f"{base_id}_source"
    target_analyzer_id = f"{base_id}_target"

    # Step 1: Create the source analyzer
    print(f"Creating source analyzer '{source_analyzer_id}'...")

    analyzer = ContentAnalyzer(
        base_analyzer_id="prebuilt-document",
        description="Source analyzer for copying",
        config=ContentAnalyzerConfig(
            enable_formula=False,
            enable_layout=True,
            enable_ocr=True,
            estimate_field_source_and_confidence=True,
            return_details=True,
        ),
        field_schema=ContentFieldSchema(
            name="company_schema",
            description="Schema for extracting company information",
            fields={
                "company_name": ContentFieldDefinition(
                    type=ContentFieldType.STRING,
                    method=GenerationMethod.EXTRACT,
                    description="Name of the company",
                ),
                "total_amount": ContentFieldDefinition(
                    type=ContentFieldType.NUMBER,
                    method=GenerationMethod.EXTRACT,
                    description="Total amount on the document",
                ),
            },
        ),
        models={"completion": "gpt-4.1"},
        tags={"modelType": "in_development"},
    )
    poller = client.begin_create_analyzer(
        analyzer_id=source_analyzer_id,
        resource=analyzer,
    )
    poller.result()
    print(f"Source analyzer '{source_analyzer_id}' created successfully!")

    # Get the source analyzer to see its description and tags before copying
    source_analyzer_info = client.get_analyzer(analyzer_id=source_analyzer_id)
    print(f"Source analyzer description: {source_analyzer_info.description}")
    if source_analyzer_info.tags:
        print(
            f"Source analyzer tags: {', '.join(f'{k}={v}' for k, v in source_analyzer_info.tags.items())}"
        )

    # [START copy_analyzer]
    print(
        f"\nCopying analyzer from '{source_analyzer_id}' to '{target_analyzer_id}'..."
    )

    poller = client.begin_copy_analyzer(
        analyzer_id=target_analyzer_id,
        source_analyzer_id=source_analyzer_id,
    )
    poller.result()

    print("Analyzer copied successfully!")
    # [END copy_analyzer]

    # [START update_and_verify_analyzer]
    # Get the target analyzer first to get its BaseAnalyzerId
    print(f"\nGetting target analyzer '{target_analyzer_id}'...")
    target_analyzer = client.get_analyzer(analyzer_id=target_analyzer_id)

    # Update the target analyzer with a production tag
    updated_analyzer = ContentAnalyzer(
        base_analyzer_id=target_analyzer.base_analyzer_id,
        tags={"modelType": "model_in_production"},
    )

    print("Updating target analyzer with production tag...")
    client.update_analyzer(analyzer_id=target_analyzer_id, resource=updated_analyzer)

    # Verify the update
    updated_target = client.get_analyzer(analyzer_id=target_analyzer_id)
    print(f"Updated target analyzer description: {updated_target.description}")
    if updated_target.tags:
        print(
            f"Updated target analyzer tag: {updated_target.tags.get('modelType', 'N/A')}"
        )
    # [END update_and_verify_analyzer]

    # [START delete_copied_analyzers]
    print("\nCleaning up analyzers...")

    try:
        client.delete_analyzer(analyzer_id=source_analyzer_id)
        print(f"  Source analyzer '{source_analyzer_id}' deleted successfully.")
    except Exception:
        pass  # Ignore cleanup errors

    try:
        client.delete_analyzer(analyzer_id=target_analyzer_id)
        print(f"  Target analyzer '{target_analyzer_id}' deleted successfully.")
    except Exception:
        pass  # Ignore cleanup errors
    # [END delete_copied_analyzers]


if __name__ == "__main__":
    main()
````

### FILE: `azure_content_understanding_official_analyzer_copy/upstream/test_sample_copy_analyzer.py`
```yaml
block_id: "AZURE-CU-OFFICIAL-ANALYZER-COPY:microsoft-test:v1"
operation: CREATE
provenance: VERBATIM
source: "Azure/azure-sdk-for-python@129a5cbb SDK recorded/live test"
license: "MIT"
sha256: "227d59a71c0e1ecaa67a9eb406654c79bc6a1294334a1d1aca6eceeb58c18f8e"
variables: []
secrets_allowed: false
```
````python
# coding: utf-8

# -------------------------------------------------------------------------
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License. See License.txt in the project root for
# license information.
# --------------------------------------------------------------------------

"""
TEST FILE: test_sample_copy_analyzer.py

DESCRIPTION:
    These tests validate the sample_copy_analyzer.py sample code.
    This sample demonstrates how to copy an analyzer from source to target within the same
    Microsoft Foundry resource using the begin_copy_analyzer API.

USAGE:
    pytest test_sample_copy_analyzer.py
"""

import uuid
import pytest
from typing import Dict
from azure.core.exceptions import ResourceNotFoundError
from devtools_testutils import recorded_by_proxy
from testpreparer import ContentUnderstandingPreparer, ContentUnderstandingClientTestBase
from azure.ai.contentunderstanding.models import (
    ContentAnalyzer,
    ContentAnalyzerConfig,
    ContentFieldSchema,
    ContentFieldDefinition,
    ContentFieldType,
    GenerationMethod,
)


class TestSampleCopyAnalyzer(ContentUnderstandingClientTestBase):
    """Tests for sample_copy_analyzer.py"""

    @ContentUnderstandingPreparer()
    @recorded_by_proxy
    def test_sample_copy_analyzer(self, contentunderstanding_endpoint: str, **kwargs) -> Dict[str, str]:
        """Test copying an analyzer (within same resource or across resources).

        This test validates:
        1. Creating a source analyzer with complex configuration
        2. Initiating a copy operation
        3. Verifying the copy completed successfully
        4. Validating the target analyzer has the same configuration

        14_CopyAnalyzer.CopyAnalyzerAsync()

        Note: This test requires copy API support. If not available, test will be skipped.
        """
        # Get variables from test proxy (recorded values in playback, empty dict in recording)
        variables = kwargs.pop("variables", {})

        try:
            client = self.create_client(endpoint=contentunderstanding_endpoint)

            # Generate unique analyzer IDs for this test
            # Use variables from recording if available (playback mode), otherwise generate new ones (record mode)
            default_source_id = f"test_analyzer_source_{uuid.uuid4().hex[:16]}"
            default_target_id = f"test_analyzer_target_{uuid.uuid4().hex[:16]}"
            source_analyzer_id = variables.setdefault("copySourceAnalyzerId", default_source_id)
            target_analyzer_id = variables.setdefault("copyTargetAnalyzerId", default_target_id)

            print(f"[INFO] Source analyzer ID: {source_analyzer_id}")
            print(f"[INFO] Target analyzer ID: {target_analyzer_id}")

            assert source_analyzer_id is not None, "Source analyzer ID should not be null"
            assert len(source_analyzer_id) > 0, "Source analyzer ID should not be empty"
            assert target_analyzer_id is not None, "Target analyzer ID should not be null"
            assert len(target_analyzer_id) > 0, "Target analyzer ID should not be empty"
            assert source_analyzer_id != target_analyzer_id, "Source and target IDs should be different"
            print("[PASS] Analyzer IDs verified")

            # Step 1: Create the source analyzer with complex configuration
            source_config = ContentAnalyzerConfig(
                enable_formula=False,
                enable_layout=True,
                enable_ocr=True,
                estimate_field_source_and_confidence=True,
                return_details=True,
            )

            # Verify source config
            assert source_config is not None, "Source config should not be null"
            assert source_config.enable_formula is False, "EnableFormula should be false"
            assert source_config.enable_layout is True, "EnableLayout should be true"
            assert source_config.enable_ocr is True, "EnableOcr should be true"
            assert (
                source_config.estimate_field_source_and_confidence is True
            ), "EstimateFieldSourceAndConfidence should be true"
            assert source_config.return_details is True, "ReturnDetails should be true"
            print("[PASS] Source config verified")

            # Create field schema
            source_field_schema = ContentFieldSchema(
                name="company_schema",
                description="Schema for extracting company information",
                fields={
                    "company_name": ContentFieldDefinition(
                        type=ContentFieldType.STRING, method=GenerationMethod.EXTRACT, description="Name of the company"
                    ),
                    "total_amount": ContentFieldDefinition(
                        type=ContentFieldType.NUMBER,
                        method=GenerationMethod.EXTRACT,
                        description="Total amount on the document",
                    ),
                },
            )

            # Verify field schema
            assert source_field_schema is not None, "Source field schema should not be null"
            assert source_field_schema.name == "company_schema", "Field schema name should match"
            assert (
                source_field_schema.description == "Schema for extracting company information"
            ), "Field schema description should match"
            assert len(source_field_schema.fields) == 2, "Should have 2 fields"
            print(f"[PASS] Source field schema verified: {source_field_schema.name}")

            # Verify individual fields
            assert "company_name" in source_field_schema.fields, "Should contain company_name field"
            company_name_field = source_field_schema.fields["company_name"]
            assert company_name_field.type == ContentFieldType.STRING, "company_name should be String type"
            assert company_name_field.method == GenerationMethod.EXTRACT, "company_name should use Extract method"
            print("[PASS] company_name field verified")

            assert "total_amount" in source_field_schema.fields, "Should contain total_amount field"
            total_amount_field = source_field_schema.fields["total_amount"]
            assert total_amount_field.type == ContentFieldType.NUMBER, "total_amount should be Number type"
            assert total_amount_field.method == GenerationMethod.EXTRACT, "total_amount should use Extract method"
            print("[PASS] total_amount field verified")

            # Create source analyzer
            source_analyzer = ContentAnalyzer(
                base_analyzer_id="prebuilt-document",
                description="Source analyzer for copying",
                config=source_config,
                field_schema=source_field_schema,
                models={"completion": "gpt-4.1"},
                tags={"modelType": "in_development"},
            )

            # Create the source analyzer
            create_poller = client.begin_create_analyzer(
                analyzer_id=source_analyzer_id, resource=source_analyzer, allow_replace=True
            )
            source_result = create_poller.result()
            print(f"[PASS] Source analyzer '{source_analyzer_id}' created successfully")

            # Step 2: Copy the analyzer
            # Note: Copy API may require authorization token for cross-resource copying
            # For same-resource copying, no authorization is needed
            print(f"\n[INFO] Attempting to copy analyzer from '{source_analyzer_id}' to '{target_analyzer_id}'")

            # Copy analyzer
            # begin_copy_analyzer requires:
            # - analyzer_id: target analyzer ID
            # - source_analyzer_id: source analyzer ID (as keyword arg)
            copy_poller = client.begin_copy_analyzer(
                analyzer_id=target_analyzer_id, source_analyzer_id=source_analyzer_id
            )
            copy_result = copy_poller.result()
            print(f"[PASS] Analyzer copied successfully to '{target_analyzer_id}'")

            print("\n[SUCCESS] All test_sample_copy_analyzer assertions passed")
            print("[INFO] Copy analyzer functionality demonstrated")

            # Return variables to be recorded for playback mode
            return variables

        except Exception as e:
            # Let all exceptions fail - don't skip
            raise
        finally:
            # Clean up: delete test analyzers
            try:
                if "source_analyzer_id" in locals() and "client" in locals():
                    client.delete_analyzer(analyzer_id=source_analyzer_id)  # type: ignore
                    print(f"\n[INFO] Source analyzer deleted: {source_analyzer_id}")  # type: ignore
            except Exception as cleanup_error:
                print(f"\n[WARN] Could not delete source analyzer: {str(cleanup_error)[:100]}")

            try:
                if "target_analyzer_id" in locals() and "client" in locals():
                    # Only try to delete if copy succeeded
                    if "copy_result" in locals():
                        client.delete_analyzer(analyzer_id=target_analyzer_id)  # type: ignore
                        print(f"[INFO] Target analyzer deleted: {target_analyzer_id}")  # type: ignore
            except Exception as cleanup_error:
                print(f"[WARN] Could not delete target analyzer: {str(cleanup_error)[:100]}")
````

### FILE: `azure_content_understanding_official_analyzer_copy/source-lock.json`
```yaml
block_id: "AZURE-CU-OFFICIAL-ANALYZER-COPY:source-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "local exact provenance and conditions lock"
license: "LicenseRef-Workspace-Owner"
sha256: "e52869fec1b1e29d5801c109aa6546149907e24da674716ab7bd750f2945fffc"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-azure-content-understanding-official-analyzer-copy/v1",
  "owner": "Microsoft Corporation",
  "repository": "Azure/azure-sdk-for-python",
  "release": "azure-ai-contentunderstanding_1.1.0",
  "commit": "129a5cbb06fba44e11ffc92f6bcdca14a2f5b6eb",
  "commit_signature_verified": true,
  "package": "azure-ai-contentunderstanding",
  "version": "1.1.0",
  "api_version": "2025-11-01",
  "sdist": {
    "url": "https://files.pythonhosted.org/packages/b6/5a/6dbcb8278c8d3779fc03992d8a9f87a9ce5280437d76e72f8121f9a90d46/azure_ai_contentunderstanding-1.1.0.tar.gz",
    "bytes": 230330,
    "sha256": "00f39c7cf2ba50c8d4586315f4295a45fad0881daaff07d13d33aea2bafc7ec8"
  },
  "files": [
    {
      "path": "upstream/sample_copy_analyzer.py",
      "source_path": "sdk/contentunderstanding/azure-ai-contentunderstanding/samples/sample_copy_analyzer.py",
      "bytes": 6071,
      "sha256": "4b8aa70c941c9b0b8fcd54ddac588ef0d1a6dbd4cefe5e2c21d11fd8e44bdb71",
      "provenance": "VERBATIM",
      "license": "MIT",
      "commit_sdist_byte_identical": true
    },
    {
      "path": "upstream/test_sample_copy_analyzer.py",
      "source_path": "sdk/contentunderstanding/azure-ai-contentunderstanding/tests/samples/test_sample_copy_analyzer.py",
      "bytes": 9468,
      "sha256": "227d59a71c0e1ecaa67a9eb406654c79bc6a1294334a1d1aca6eceeb58c18f8e",
      "provenance": "VERBATIM",
      "license": "MIT",
      "commit_sdist_byte_identical": true,
      "execution": "recorded/live Azure SDK test harness; not standalone"
    },
    {
      "path": "upstream/LICENSE.txt",
      "source_path": "repository-root/LICENSE",
      "bytes": 1074,
      "sha256": "7c77a44a8acd9b41fdc209864a8016b3d430b5d0e09309818d5b7444336df744",
      "provenance": "VERBATIM",
      "license": "MIT"
    }
  ],
  "official_claim": {
    "operation": "same-resource exact analyzer copy",
    "source_base_analyzer_id": "prebuilt-document",
    "automatic_execution": false,
    "automatic_production_promotion": false
  },
  "conditions": [
    "Microsoft Foundry resource and supported region",
    "endpoint and key or DefaultAzureCredential supplied outside files",
    "models and possible provider cost approved",
    "explicit authority for analyzer create, copy, update and delete",
    "target-specific immutable analyzer identity and version receipt",
    "independent inventory and reconciliation because cleanup errors are swallowed",
    "project corpus and ground truth evaluated before production use"
  ]
}
````

### FILE: `azure_content_understanding_official_analyzer_copy/test_official_copy_sample.py`
```yaml
block_id: "AZURE-CU-OFFICIAL-ANALYZER-COPY:offline-contract:v1"
operation: CREATE
provenance: AUTHORED
source: "local offline identity and non-claim regression"
license: "LicenseRef-Workspace-Owner"
sha256: "89c46e57d5982d1374d6f300bfc0bbde70f36a12f4452fdac8ac2e2e77656349"
variables: []
secrets_allowed: false
```
````python
import ast
import hashlib
import json
import unittest
from pathlib import Path


ROOT = Path(__file__).resolve().parent
SAMPLE = ROOT / "upstream" / "sample_copy_analyzer.py"
UPSTREAM_TEST = ROOT / "upstream" / "test_sample_copy_analyzer.py"
LICENSE = ROOT / "upstream" / "LICENSE.txt"
LOCK = ROOT / "source-lock.json"


def sha256(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


class OfficialAnalyzerCopyContract(unittest.TestCase):
    def test_exact_official_bytes_and_lock(self) -> None:
        lock = json.loads(LOCK.read_text(encoding="utf-8"))
        expected = {item["path"]: item for item in lock["files"]}
        self.assertEqual(lock["commit"], "129a5cbb06fba44e11ffc92f6bcdca14a2f5b6eb")
        self.assertEqual(set(expected), {
            "upstream/sample_copy_analyzer.py",
            "upstream/test_sample_copy_analyzer.py",
            "upstream/LICENSE.txt",
        })
        for relative, item in expected.items():
            path = ROOT / relative
            self.assertEqual(path.stat().st_size, item["bytes"])
            self.assertEqual(sha256(path), item["sha256"])
            self.assertEqual(item["provenance"], "VERBATIM")

    def test_sample_preserves_official_lifecycle_and_non_claims(self) -> None:
        source = SAMPLE.read_text(encoding="utf-8")
        tree = ast.parse(source)
        calls = {
            node.func.attr
            for node in ast.walk(tree)
            if isinstance(node, ast.Call) and isinstance(node.func, ast.Attribute)
        }
        self.assertTrue({
            "begin_create_analyzer",
            "begin_copy_analyzer",
            "get_analyzer",
            "update_analyzer",
            "delete_analyzer",
        }.issubset(calls))
        self.assertIn('base_analyzer_id="prebuilt-document"', source)
        self.assertIn("estimate_field_source_and_confidence=True", source)
        self.assertNotIn("prebuilt-procurement", source)
        self.assertNotIn("prebuilt-purchaseOrder", source)

    def test_official_test_is_preserved_but_not_misrepresented_as_standalone(self) -> None:
        source = UPSTREAM_TEST.read_text(encoding="utf-8")
        tree = ast.parse(source)
        lock = json.loads(LOCK.read_text(encoding="utf-8"))
        test_record = next(item for item in lock["files"] if item["path"].endswith("test_sample_copy_analyzer.py"))
        self.assertIn("@recorded_by_proxy", source)
        self.assertIn("ContentUnderstandingPreparer", source)
        self.assertIn("finally:", source)
        self.assertNotIn("pytest.skip", source)
        self.assertEqual(test_record["execution"], "recorded/live Azure SDK test harness; not standalone")
        self.assertGreater(len(list(ast.walk(tree))), 100)


if __name__ == "__main__":
    unittest.main()
````

## 6. Configuration surface

| Input | Offline | Live rule |
|---|---:|---|
| source lock | yes | immutable exact commit/sdist identity |
| `CONTENTUNDERSTANDING_ENDPOINT` | no | external secret-free endpoint reference |
| `CONTENTUNDERSTANDING_KEY` / identity | no | external secret store; prefer approved identity |
| analyzer source/target IDs | no | unique, project-owned, receipt-bound |
| models/region/cost | no | explicit project authority before mutation |

## 7. Dependency bill

| Dependency | Pin/status | License | Use |
|---|---|---|---|
| `azure-ai-contentunderstanding` | 1.1.0 wheel/sdist exactos en artifact core | MIT | official client/models/LRO |
| `azure-identity` | project lock required | MIT | `DefaultAzureCredential` |
| `python-dotenv` | project lock required | BSD-3-Clause | official sample environment loading |
| Azure SDK test-proxy/devtools | monorepo harness only | MIT | official test, not standalone |

No se inventan pins auxiliares: el proyecto debe resolverlos conforme al dependency update contract y probar el grafo exacto.

## 8. Apply order

1. Materializar este pack y ejecutar sólo sus contratos offline.
2. Materializar/adquirir el artifact core 1.1.0 exacto.
3. Cerrar cuenta, identidad, modelos, costo, IDs, authority y rollback/reconciliation.
4. Ejecutar el test oficial únicamente dentro del harness Azure soportado o una reproducción live autorizada.
5. Registrar source/target analyzer identity, operación, resultado e inventario posterior antes de usar el target.

## 9. Verification

```powershell
python -m unittest -v azure_content_understanding_official_analyzer_copy/test_official_copy_sample.py
python -m py_compile azure_content_understanding_official_analyzer_copy/upstream/sample_copy_analyzer.py azure_content_understanding_official_analyzer_copy/upstream/test_sample_copy_analyzer.py azure_content_understanding_official_analyzer_copy/test_official_copy_sample.py
```

Debe ejecutar 3/3 contratos offline, conservar hashes Microsoft exactos y no realizar llamadas externas. La prueba live no forma parte del gate offline.

## 10. Reconstruction evidence

Estado: `REBUILD_VERIFIED / CONDITIONED`. Commit y sdist 1.1.0 son byte-idénticos para sample/test; ambos parsean. La promoción sigue bloqueada hasta resource/identity/model/cost, mutación explícita, test-proxy o sandbox live, receipt, inventario/reconciliación, corpus y rollback target.
