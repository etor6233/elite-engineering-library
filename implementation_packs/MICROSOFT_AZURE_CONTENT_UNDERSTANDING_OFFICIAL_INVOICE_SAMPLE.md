# Microsoft Azure Content Understanding Official Invoice Sample

## 1. Metadata

```yaml
pack_id: "MICROSOFT-AZURE-CONTENT-UNDERSTANDING-OFFICIAL-INVOICE-SAMPLE"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa byte a byte el sample Python síncrono oficial Microsoft Azure AI Content Understanding 1.1.0 que llama prebuilt-invoice y expone campos, confidence, sources, spans, line items y usage para invoices, utility bills, sales orders y purchase orders, junto con licencia, source lock y tres contratos offline."
stacks: ["Python 3.9+", "azure-ai-contentunderstanding 1.1.0", "Azure Content Understanding API 2025-11-01"]
compatible_with: ["OFFICIAL-DOCUMENT-SDK-ARTIFACT-CORE 0.3.x", "AZURE-CONTENT-UNDERSTANDING-DOCUMENT-RUNTIME 0.1.x", "STRICT-DOCUMENT-FIELD-EVALUATION-GATE 0.1.x"]
incompatible_with: ["persistencia automática", "URL móvil como input empresarial", "dependencias auxiliares sin lock", "clases no documentadas por prebuilt-invoice"]
license_expression: "LicenseRef-Workspace-Owner AND MIT"
upstream_sources: ["https://github.com/Azure/azure-sdk-for-python/tree/129a5cbb06fba44e11ffc92f6bcdca14a2f5b6eb/sdk/contentunderstanding/azure-ai-contentunderstanding", "https://pypi.org/project/azure-ai-contentunderstanding/1.1.0/", "https://learn.microsoft.com/azure/ai-services/content-understanding/concepts/prebuilt-analyzers"]
verified_at: "2026-08-28"
```

## 2. Applicability

Use este pack cuando el proyecto haya elegido Microsoft Azure Content Understanding y necesite partir del sample oficial estable de `prebuilt-invoice` sin reescribir su llamada, object model ni extracción demostrativa. Es inmediatamente materializable y verificable offline. La ejecución live requiere recurso Foundry, región, modelos, identidad, costo aceptado y un lock target para las dependencias auxiliares del sample.

No use este pack para proformas, packing lists, bills of lading, aduana o layouts particulares: Microsoft no los atribuye a este analyzer. No lo use para persistencia automática ni como prueba de precisión. `prebuilt-procurement` y `prebuilt-purchaseOrder` existen en la documentación vigente, pero este source exacto llama únicamente `prebuilt-invoice`; adoptar otros IDs exige su propia evidencia.

## 3. Architecture contract

`upstream/sample_analyze_invoice.py` y `upstream/LICENSE.txt` son bytes `VERBATIM` del commit firmado `129a5cbb…` y del sdist PyPI 1.1.0. El sample/commit y el sample/sdist son byte-idénticos. Los otros tres archivos son packaging `AUTHORED` y nunca se atribuyen a Microsoft.

El sample crea `ContentUnderstandingClient`, llama `begin_analyze(analyzer_id="prebuilt-invoice")`, espera el LRO y presenta `CustomerName`, `InvoiceDate`, `TotalAmount`, `LineItems`, confidence, source, spans y usage. No recibe documentos del proyecto de forma segura: la demo usa una URL bajo `main`. El flujo empresarial debe usar la lane local segura y sólo reutilizar los contratos del SDK/analyzer.

Fallos: ausencia de endpoint produce error; credencial, red, cuota, modelo o servicio fallidos permanecen errores del SDK; contenido vacío termina sin persistir. El sample sólo imprime y no implementa retries de negocio, idempotencia, revisión, reconciliación ni storage. Rollback elimina sólo el directorio recién materializado. Presupuesto offline menor a cinco segundos; performance/costo live dependen del corpus y proveedor.

## 4. Exact file manifest

```text
CREATE azure_content_understanding_official_invoice/README.md
CREATE azure_content_understanding_official_invoice/upstream/LICENSE.txt
CREATE azure_content_understanding_official_invoice/upstream/sample_analyze_invoice.py
CREATE azure_content_understanding_official_invoice/source-lock.json
CREATE azure_content_understanding_official_invoice/test_official_sample.py
```

## 5. Materialization blocks

### FILE: `azure_content_understanding_official_invoice/README.md`
```yaml
block_id: "AZURE-CU-OFFICIAL-INVOICE:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local guide constrained by exact Microsoft source"
license: "LicenseRef-Workspace-Owner"
sha256: "50f20ece578433df61dd6e9bf1d911b05aa92a62100d08fb761df7634e3d763a"
variables: []
secrets_allowed: false
```
````markdown
# Microsoft Azure Content Understanding official invoice sample

This directory contains Microsoft's byte-verbatim `azure-ai-contentunderstanding` 1.1.0 synchronous invoice sample and MIT license. The source commit, PyPI sdist and sample SHA-256 are fixed in `source-lock.json`. `test_official_sample.py` is Elite-authored packaging verification; it is not Microsoft product code.

The official sample calls `prebuilt-invoice`, which Microsoft documents for invoices, utility bills, sales orders and purchase orders. It prints structured fields, confidence, source, spans and usage. It does not authorize persistence, define a project's business schema or cover proformas, packing lists, bills of lading, customs or supplier-specific variants.

Offline verification:

```powershell
python .\test_official_sample.py
```

Live execution requires the exact official wheel plus its project-specific locked graph, `python-dotenv`, `azure-identity`, a Microsoft Foundry endpoint, an API key or approved identity, configured model deployments, region/quota/cost authority and service access. The upstream demo reads an asset from a moving `main` URL. Do not send project documents through that URL. For real files use the secure local ingestion and Azure document runtime profile, which submits admitted local bytes and preserves evidence without modifying this upstream sample.

No output from this sample is trusted business data until the project's corpus, ground truth, field evaluation, review, reconciliation, retention and storage authorization gates pass.
````

### FILE: `azure_content_understanding_official_invoice/upstream/LICENSE.txt`
```yaml
block_id: "AZURE-CU-OFFICIAL-INVOICE:microsoft-license:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/Azure/azure-sdk-for-python/129a5cbb06fba44e11ffc92f6bcdca14a2f5b6eb/LICENSE"
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

### FILE: `azure_content_understanding_official_invoice/upstream/sample_analyze_invoice.py`
```yaml
block_id: "AZURE-CU-OFFICIAL-INVOICE:microsoft-sample:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/Azure/azure-sdk-for-python/129a5cbb06fba44e11ffc92f6bcdca14a2f5b6eb/sdk/contentunderstanding/azure-ai-contentunderstanding/samples/sample_analyze_invoice.py"
license: "MIT"
sha256: "0cb9d7b0e183cd1c3f677c79a9d8d5dd9a1ed30bf397f84b753eedd7819c53fd"
variables: []
secrets_allowed: true
```
````python
# pylint: disable=line-too-long,useless-suppression
# mypy: disable-error-code="attr-defined"
# coding=utf-8
# --------------------------------------------------------------------------
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License. See License.txt in the project root for license information.
# --------------------------------------------------------------------------
"""
FILE: sample_analyze_invoice.py

DESCRIPTION:
    This sample demonstrates how to analyze an invoice from a URL using the prebuilt-invoice analyzer
    and extract structured fields from the result.

    ## About analyzing invoices

    Content Understanding provides a rich set of prebuilt analyzers that are ready to use without any
    configuration. These analyzers are powered by knowledge bases of thousands of real-world document
    examples, enabling them to understand document structure and adapt to variations in format and
    content.

    Prebuilt analyzers are ideal for:
    - Content ingestion in search and retrieval-augmented generation (RAG) workflows
    - Intelligent document processing (IDP) to extract structured data from common document types
    - Agentic flows as tools for extracting structured representations from input files

    ### The prebuilt-invoice analyzer

    The prebuilt-invoice analyzer is a domain-specific analyzer optimized for processing invoices,
    utility bills, sales orders, and purchase orders. It automatically extracts structured fields
    including:

    - Customer/Vendor information: Name, address, contact details
    - Invoice metadata: Invoice number, date, due date, purchase order number
    - Line items: Description, quantity, unit price, total for each item
    - Financial totals: Subtotal, tax amount, shipping charges, total amount
    - Payment information: Payment terms, payment method, remittance address

    The analyzer works out of the box with various invoice formats and requires no configuration.
    It's part of the financial documents category of prebuilt analyzers, which also includes:
    - prebuilt-receipt - Sales receipts from retail and dining establishments
    - prebuilt-creditCard - Credit card statements
    - prebuilt-bankStatement.us - US bank statements
    - prebuilt-check.us - US bank checks
    - prebuilt-creditMemo - Credit memos and refund documents

USAGE:
    python sample_analyze_invoice.py

    Set the environment variables with your own values before running the sample:
    1) CONTENTUNDERSTANDING_ENDPOINT - the endpoint to your Content Understanding resource.
    2) CONTENTUNDERSTANDING_KEY - your Content Understanding API key (optional if using DefaultAzureCredential).

    Before using prebuilt analyzers, you MUST configure model deployments for your Microsoft Foundry
    resource. See sample_update_defaults.py for setup instructions.
"""

import os
from typing import cast

from dotenv import load_dotenv
from azure.ai.contentunderstanding import ContentUnderstandingClient
from azure.ai.contentunderstanding.models import (
    AnalysisInput,
    AnalysisResult,
    DocumentContent,
    ContentField,
    ArrayField,
    ObjectField,
)
from azure.core.credentials import AzureKeyCredential
from azure.identity import DefaultAzureCredential

load_dotenv()


def main() -> None:
    endpoint = os.environ["CONTENTUNDERSTANDING_ENDPOINT"]
    key = os.getenv("CONTENTUNDERSTANDING_KEY")
    credential = AzureKeyCredential(key) if key else DefaultAzureCredential()

    client = ContentUnderstandingClient(endpoint=endpoint, credential=credential)

    # [START analyze_invoice]
    # You can replace this URL with your own invoice file URL
    invoice_url = "https://raw.githubusercontent.com/Azure-Samples/azure-ai-content-understanding-assets/main/document/invoice.pdf"

    print("Analyzing invoice with prebuilt-invoice analyzer...")
    print(f"  URL: {invoice_url}\n")

    poller = client.begin_analyze(
        analyzer_id="prebuilt-invoice",
        inputs=[AnalysisInput(url=invoice_url)],
    )
    result: AnalysisResult = poller.result()
    # [END analyze_invoice]

    # [START extract_invoice_fields]
    if not result.contents or len(result.contents) == 0:
        print("No content found in the analysis result.")
        return

    # Get the document content (invoices are documents)
    document_content = cast(DocumentContent, result.contents[0])

    # Print document unit information
    # The unit indicates the measurement system used for coordinates in the source field
    print(f"Document unit: {document_content.unit or 'unknown'}")
    print(
        f"Pages: {document_content.start_page_number} to {document_content.end_page_number}"
    )

    # Print page dimensions if available
    if document_content.pages and len(document_content.pages) > 0:
        page = document_content.pages[0]
        unit = document_content.unit or "units"
        print(f"Page dimensions: {page.width} x {page.height} {unit}")
    print()

    if not document_content.fields:
        print("No fields found in the analysis result.")
        return

    # Extract simple string fields
    customer_name_field = document_content.fields.get("CustomerName")
    print(
        f"Customer Name: {customer_name_field.value or '(None)' if customer_name_field else '(None)'}"
    )
    if customer_name_field:
        print(
            f"  Confidence: {customer_name_field.confidence:.2f}"
            if customer_name_field.confidence
            else "  Confidence: N/A"
        )
        print(f"  Source: {customer_name_field.source or 'N/A'}")
        if customer_name_field.spans and len(customer_name_field.spans) > 0:
            span = customer_name_field.spans[0]
            print(f"  Position in markdown: offset={span.offset}, length={span.length}")

    # Extract simple date field
    invoice_date_field = document_content.fields.get("InvoiceDate")
    print(
        f"Invoice Date: {invoice_date_field.value or '(None)' if invoice_date_field else '(None)'}"
    )
    if invoice_date_field:
        print(
            f"  Confidence: {invoice_date_field.confidence:.2f}"
            if invoice_date_field.confidence
            else "  Confidence: N/A"
        )
        print(f"  Source: {invoice_date_field.source or 'N/A'}")
        if invoice_date_field.spans and len(invoice_date_field.spans) > 0:
            span = invoice_date_field.spans[0]
            print(f"  Position in markdown: offset={span.offset}, length={span.length}")

    # Extract object fields (nested structures)
    total_amount_field = document_content.fields.get("TotalAmount")
    if isinstance(total_amount_field, ObjectField) and total_amount_field.value:
        amount_field = total_amount_field.value.get("Amount")
        currency_field = total_amount_field.value.get("CurrencyCode")
        amount = amount_field.value if amount_field else None
        # Use currency value if present, otherwise default to ""
        currency = (
            currency_field.value if currency_field and currency_field.value else ""
        )
        if isinstance(amount, (int, float)):
            print(f"\nTotal: {currency}{amount:.2f}")
        else:
            print(f"\nTotal: {currency}{amount or '(None)'}")
        print(
            f"  Amount Confidence: {amount_field.confidence:.2f}"
            if amount_field and amount_field.confidence
            else "  Amount Confidence: N/A"
        )
        print(
            f"  Source for Amount: {amount_field.source or 'N/A'}"
            if amount_field
            else "  Source: N/A"
        )

    # Extract array fields (collections like line items)
    line_items_field = document_content.fields.get("LineItems")
    if isinstance(line_items_field, ArrayField) and line_items_field.value:
        print(f"\nLine Items ({len(line_items_field.value)}):")
        for i, item in enumerate(line_items_field.value, 1):
            if isinstance(item, ObjectField) and item.value:
                description_field = item.value.get("Description")
                quantity_field = item.value.get("Quantity")
                description = (
                    description_field.value
                    if description_field and description_field.value
                    else "N/A"
                )
                quantity = (
                    quantity_field.value
                    if quantity_field and quantity_field.value
                    else "N/A"
                )
                print(f"  Item {i}: {description}")
                print(f"    Quantity: {quantity}")
                print(
                    f"    Quantity Confidence: {quantity_field.confidence:.2f}"
                    if quantity_field and quantity_field.confidence
                    else "    Quantity Confidence: N/A"
                )
    # [END extract_invoice_fields]

    # [START get_usage]
    # Access usage details from the poller (available after result() completes).
    # Usage reports resource consumption for billing estimation:
    #
    # - document_pages_standard/basic/minimal: Pages processed at each extraction tier.
    #   Standard = layout + OCR (scanned docs), Basic = OCR only, Minimal = digital formats
    #   (DOCX, XLSX, HTML, TXT) that need no OCR. Charged per 1,000 pages.
    #
    # - contextualization_tokens: Fixed-rate tokens charged by Content Understanding for
    #   preparing context, generating confidence scores, source grounding, and formatting
    #   output. Typically 1,000 tokens per page. Charged separately from LLM tokens.
    #
    # - tokens: Dict of "{model}-input" / "{model}-output" token counts consumed by your
    #   Foundry model deployment (e.g. "gpt-4.1-input", "gpt-4.1-output"). These are
    #   billed on your Foundry deployment, not on Content Understanding.
    #
    # For full pricing details, see:
    # https://learn.microsoft.com/azure/ai-services/content-understanding/pricing-explainer
    usage = poller.usage
    if usage:
        print("\nUsage Details:")
        if usage.document_pages_standard is not None:
            print(f"  Document pages (standard): {usage.document_pages_standard}")
        if usage.contextualization_tokens is not None:
            print(f"  Contextualization tokens: {usage.contextualization_tokens}")
        if usage.tokens:
            print("  Model tokens:")
            for model, count in usage.tokens.items():
                print(f"    {model}: {count}")
    # [END get_usage]


if __name__ == "__main__":
    main()
````

### FILE: `azure_content_understanding_official_invoice/source-lock.json`
```yaml
block_id: "AZURE-CU-OFFICIAL-INVOICE:source-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "registry derived from Microsoft signed release, PyPI artifacts and official docs"
license: "LicenseRef-Workspace-Owner"
sha256: "8ec2f600bbc996d126b19668c0c3f28701032b451b9f545ace50805be277814f"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-azure-content-understanding-official-sample/v1",
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
  "wheel": {
    "url": "https://files.pythonhosted.org/packages/67/d6/6f60ac39b73a0b8a2bc86cd6895e26eccc20d9b7bde86814fe6a36964dff/azure_ai_contentunderstanding-1.1.0-py3-none-any.whl",
    "bytes": 101987,
    "sha256": "d1d6bdeffe02f5c8cc5309f0e120a1338f51af04638605f9617b7b501e2d1dd6"
  },
  "files": [
    {
      "path": "upstream/sample_analyze_invoice.py",
      "source_path": "sdk/contentunderstanding/azure-ai-contentunderstanding/samples/sample_analyze_invoice.py",
      "bytes": 10539,
      "sha256": "0cb9d7b0e183cd1c3f677c79a9d8d5dd9a1ed30bf397f84b753eedd7819c53fd",
      "provenance": "VERBATIM",
      "license": "MIT",
      "commit_sdist_byte_identical": true
    },
    {
      "path": "upstream/LICENSE.txt",
      "source_path": "repository-root/LICENSE",
      "bytes": 1074,
      "sha256": "7c77a44a8acd9b41fdc209864a8016b3d430b5d0e09309818d5b7444336df744",
      "provenance": "VERBATIM",
      "license": "MIT",
      "sdist_variant_bytes": 1073,
      "sdist_variant_sha256": "fd532481d828e13a0b13ccb598e02338a3617740675a862ee6bdc1541b68e93d",
      "difference": "the sdist license omits only the repository file's final LF; this pack uses repository commit bytes"
    }
  ],
  "official_claim": {
    "analyzer_id": "prebuilt-invoice",
    "document_types": ["invoice", "utility bill", "sales order", "purchase order"],
    "automatic_storage": false
  },
  "conditions": [
    "Microsoft Foundry resource and supported region",
    "model deployments configured before prebuilt analyzers",
    "endpoint and key or DefaultAzureCredential supplied outside files",
    "possible provider and model token cost accepted",
    "moving demo asset URL not used for project documents",
    "project corpus and ground truth evaluated per class and field",
    "review and persistence authority remain external gates"
  ]
}
````

### FILE: `azure_content_understanding_official_invoice/test_official_sample.py`
```yaml
block_id: "AZURE-CU-OFFICIAL-INVOICE:offline-contract:v1"
operation: CREATE
provenance: AUTHORED
source: "local offline contract over exact Microsoft sample"
license: "LicenseRef-Workspace-Owner"
sha256: "d893bd91658b68f71b4b723edaecec5cfebab3e7abb2f72e93853cfa740242b8"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import ast
import hashlib
import importlib.util
import json
import os
from pathlib import Path
import sys
import types
import unittest
from unittest.mock import patch


ROOT = Path(__file__).resolve().parent
SAMPLE = ROOT / "upstream" / "sample_analyze_invoice.py"
LICENSE = ROOT / "upstream" / "LICENSE.txt"


def sha256(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


class AnalysisInput:
    def __init__(self, *, url: str):
        self.url = url


class AnalysisResult:
    pass


class DocumentContent:
    pass


class ContentField:
    pass


class ArrayField:
    pass


class ObjectField:
    pass


class AzureKeyCredential:
    def __init__(self, key: str):
        self.key = key


class DefaultAzureCredential:
    pass


class FakePoller:
    usage = None

    def result(self):
        return types.SimpleNamespace(contents=[])


class FakeClient:
    instances: list["FakeClient"] = []

    def __init__(self, *, endpoint, credential):
        self.endpoint = endpoint
        self.credential = credential
        self.calls = []
        self.__class__.instances.append(self)

    def begin_analyze(self, **kwargs):
        self.calls.append(kwargs)
        return FakePoller()


def fake_modules() -> dict[str, types.ModuleType]:
    dotenv = types.ModuleType("dotenv")
    dotenv.load_dotenv = lambda: None
    azure = types.ModuleType("azure")
    ai = types.ModuleType("azure.ai")
    contentunderstanding = types.ModuleType("azure.ai.contentunderstanding")
    models = types.ModuleType("azure.ai.contentunderstanding.models")
    core = types.ModuleType("azure.core")
    credentials = types.ModuleType("azure.core.credentials")
    identity = types.ModuleType("azure.identity")
    contentunderstanding.ContentUnderstandingClient = FakeClient
    for value in (AnalysisInput, AnalysisResult, DocumentContent, ContentField, ArrayField, ObjectField):
        setattr(models, value.__name__, value)
    credentials.AzureKeyCredential = AzureKeyCredential
    identity.DefaultAzureCredential = DefaultAzureCredential
    return {
        "dotenv": dotenv,
        "azure": azure,
        "azure.ai": ai,
        "azure.ai.contentunderstanding": contentunderstanding,
        "azure.ai.contentunderstanding.models": models,
        "azure.core": core,
        "azure.core.credentials": credentials,
        "azure.identity": identity,
    }


class OfficialSampleContract(unittest.TestCase):
    def setUp(self):
        FakeClient.instances.clear()
        self.lock = json.loads((ROOT / "source-lock.json").read_text(encoding="utf-8"))

    def test_exact_official_bytes_and_identity(self):
        files = {item["path"]: item for item in self.lock["files"]}
        self.assertEqual(SAMPLE.stat().st_size, files["upstream/sample_analyze_invoice.py"]["bytes"])
        self.assertEqual(sha256(SAMPLE), files["upstream/sample_analyze_invoice.py"]["sha256"])
        self.assertEqual(LICENSE.stat().st_size, files["upstream/LICENSE.txt"]["bytes"])
        self.assertEqual(sha256(LICENSE), files["upstream/LICENSE.txt"]["sha256"])
        self.assertEqual(self.lock["commit"], "129a5cbb06fba44e11ffc92f6bcdca14a2f5b6eb")
        self.assertTrue(files["upstream/sample_analyze_invoice.py"]["commit_sdist_byte_identical"])

    def test_source_parses_and_keeps_official_analyzer(self):
        source = SAMPLE.read_text(encoding="utf-8")
        tree = ast.parse(source, filename=str(SAMPLE))
        analyzer_values = [
            keyword.value.value
            for node in ast.walk(tree)
            if isinstance(node, ast.Call)
            for keyword in node.keywords
            if keyword.arg == "analyzer_id" and isinstance(keyword.value, ast.Constant)
        ]
        self.assertEqual(analyzer_values, ["prebuilt-invoice"])
        self.assertIn("Azure-Samples/azure-ai-content-understanding-assets/main/document/invoice.pdf", source)
        self.assertNotIn("business_storage", source)

    def test_official_main_calls_sdk_without_persistence(self):
        module_name = "official_microsoft_sample_analyze_invoice"
        with patch.dict(sys.modules, fake_modules(), clear=False), patch.dict(
            os.environ,
            {"CONTENTUNDERSTANDING_ENDPOINT": "https://example.services.ai.azure.com", "CONTENTUNDERSTANDING_KEY": "test-only"},
            clear=False,
        ):
            spec = importlib.util.spec_from_file_location(module_name, SAMPLE)
            self.assertIsNotNone(spec)
            self.assertIsNotNone(spec.loader)
            module = importlib.util.module_from_spec(spec)
            spec.loader.exec_module(module)
            module.main()
        client = FakeClient.instances[-1]
        self.assertEqual(client.endpoint, "https://example.services.ai.azure.com")
        self.assertEqual(len(client.calls), 1)
        call = client.calls[0]
        self.assertEqual(call["analyzer_id"], "prebuilt-invoice")
        self.assertEqual(len(call["inputs"]), 1)
        self.assertTrue(call["inputs"][0].url.endswith("/document/invoice.pdf"))


if __name__ == "__main__":
    unittest.main()
````

## 6. Configuration surface

| Variable | Tipo | Default | Validación | Secreto | Efecto |
|---|---|---|---|---|---|
| `CONTENTUNDERSTANDING_ENDPOINT` | URL HTTPS | ninguno | recurso Foundry autorizado y región elegida | no | endpoint del cliente |
| `CONTENTUNDERSTANDING_KEY` | string | ausente | secret manager; si falta el sample usa DefaultAzureCredential | sí | autenticación API key |
| analyzer | literal | `prebuilt-invoice` | fijado por source; otro ID exige nueva evidencia | no | contrato de extracción |
| input demo | URL | asset `main` Microsoft | sólo demo; prohibido para documentos del proyecto | no | input de ejemplo |

No se guardan credenciales en Markdown, `.env`, receipts o logs. El target productivo fija API version, analyzer copy/version, modelos, región y schema en la decisión documental.

## 7. Dependency bill

| Dependencia | Pin/identidad | Licencia | Uso | Condición |
|---|---|---|---|---|
| `azure-ai-contentunderstanding` | 1.1.0 wheel SHA `d1d6bdef…d1dd6` | MIT | cliente/modelos/LRO | adquirir por artifact core |
| source distribution | 1.1.0, 230.330 bytes, SHA `00f39c7c…ec8` | MIT | source/sample/test | exacto |
| `python-dotenv` | upstream no fija versión | BSD-3-Clause | carga `.env` del sample | resolver/pinnear/escanear en target; no inferido aquí |
| `azure-identity` | upstream no fija versión | MIT | DefaultAzureCredential | resolver/pinnear/escanear en target; no inferido aquí |
| Python | 3.9+ | PSF-2.0 | runtime | versión elegida por proyecto |

El `dev_requirements.txt` oficial del commit tiene 128 bytes/SHA-256 `c00259b7…215`, pero enumera estas dependencias sin versiones. Por eso el pack no promete live reproducible hasta que el proyecto produzca el lock exacto de su plataforma.

## 8. Apply order

1. Complete readiness, provider/costo/región/modelos y decisión documental.
2. Materialice este pack en destino vacío y ejecute los tres contratos offline.
3. Adquiera `azure-ai-contentunderstanding` 1.1.0 mediante el artifact core.
4. Resuelva en aislamiento las dependencias auxiliares para OS/Python, fije hashes, licencias y SCA; no instale rangos abiertos.
5. Para smoke live autorizado, configure secretos fuera del árbol y use sólo el asset demo.
6. Para archivos reales, componga la lane Azure segura local, evalúe corpus/ground truth y conserve storage=false hasta que review/persist/reconciliation pasen.

En workspace existente, abortar ante cualquier ruta ocupada. Una actualización requiere nuevo release/commit/artifact/hash, diff, locks, tests y versión del pack.

## 9. Verification

```powershell
pwsh -NoProfile -File .\materialize_markdown_pack.ps1 -PackFile .\implementation_packs\MICROSOFT_AZURE_CONTENT_UNDERSTANDING_OFFICIAL_INVOICE_SAMPLE.md -Destination <empty>
python <empty>\azure_content_understanding_official_invoice\test_official_sample.py
python -m compileall -q <empty>\azure_content_understanding_official_invoice
```

Éxito offline esperado: cinco archivos con SHA exacto, `Ran 3 tests ... OK`, analyzer literal `prebuilt-invoice`, un único call al fake SDK y cero storage. El test no simula calidad del proveedor. PASS live requiere recurso, modelos, credencial, costo y evidencia separada; PASS de persistencia requiere corpus/ground truth/review/reconciliation.

## 10. Reconstruction evidence

- source commit firmado `129a5cbb…`; sdist 230.330 bytes/SHA `00f39c7c…ec8`; wheel 101.987 bytes/SHA `d1d6bdef…d1dd6`;
- sample commit/sdist byte-idéntico: 10.539 bytes/SHA `0cb9d7b0…c53fd`;
- árbol temporal limpio: tres contratos PASS y compileall PASS en 2026-08-28;
- no se realizó llamada Azure ni se atribuye precisión productiva;
- evidencia de cierre: `reconstruction_evidence/MICROSOFT_AZURE_CU_OFFICIAL_INVOICE_SAMPLE_2026-08-28_V1.md`.
