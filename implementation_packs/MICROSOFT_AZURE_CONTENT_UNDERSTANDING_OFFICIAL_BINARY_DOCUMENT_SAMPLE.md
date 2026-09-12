# Microsoft Azure Content Understanding Official Binary Document Sample

## 1. Metadata

```yaml
pack_id: "MICROSOFT-AZURE-CONTENT-UNDERSTANDING-OFFICIAL-BINARY-DOCUMENT-SAMPLE"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa byte a byte el sample Python binario y su test oficial Microsoft para prebuilt-documentSearch, fijados simultáneamente al commit Azure SDK y al sdist 1.1.0, junto con licencia, source lock y contratos offline."
stacks: ["Python 3.9+", "azure-ai-contentunderstanding 1.1.0", "Azure Content Understanding API 2025-11-01"]
compatible_with: ["OFFICIAL-DOCUMENT-SDK-ARTIFACT-CORE 0.3.x", "AZURE-CONTENT-UNDERSTANDING-DOCUMENT-RUNTIME 0.1.x", "SECURE-LOCAL-FILE-INGESTION-FOUNDATION 0.1.x"]
incompatible_with: ["persistencia automática", "exactitud de campos sin corpus", "schemas universales de negocio", "ejecución live sin autoridad Azure"]
license_expression: "LicenseRef-Workspace-Owner AND MIT"
upstream_sources: ["https://github.com/Azure/azure-sdk-for-python/tree/129a5cbb06fba44e11ffc92f6bcdca14a2f5b6eb/sdk/contentunderstanding/azure-ai-contentunderstanding", "https://pypi.org/project/azure-ai-contentunderstanding/1.1.0/", "https://learn.microsoft.com/azure/ai-services/content-understanding/concepts/prebuilt-analyzers"]
verified_at: "2026-08-28"
```

## 2. Applicability

Use este pack cuando un proyecto haya seleccionado Microsoft Azure Content Understanding y necesite el código oficial que envía bytes locales de PDF, conserva layout/tablas/figuras/markdown y demuestra rangos de páginas mediante `prebuilt-documentSearch`. Es materializable y verificable offline; el test Microsoft queda incluido como evidencia ejecutable bajo su propio harness.

No lo use como extractor de campos de proforma, packing list, bill of lading, aduana o purchase order. `prebuilt-documentSearch` preserva contenido documental; no sustituye el analyzer de negocio. La selección de `prebuilt-procurement`, `prebuilt-purchaseOrder` o un analyzer custom requiere configuración oficial, corpus, ground truth y gates propios.

## 3. Architecture contract

Los tres archivos bajo `upstream/` son bytes `VERBATIM` de Microsoft. Muestra y test coinciden byte a byte entre el commit firmado `129a5cbb…` y el sdist PyPI 1.1.0. Los otros tres archivos son packaging `AUTHORED` y nunca se atribuyen a Microsoft.

El sample lee bytes locales, llama `begin_analyze_binary`, espera el LRO, preserva `AnalysisResult` y expone markdown, páginas y tablas. También demuestra rangos `3-` y `1-3,5,9-`. No implementa intake hostil, antivirus, idempotencia, persistencia, revisión ni reconciliación: esas fronteras pertenecen a los packs de ingestión segura y evaluación estricta.

Fallos de archivo, credencial, endpoint, modelo, cuota, red o servicio deben permanecer errores. No se ejecuta live automáticamente. Rollback elimina sólo el directorio recién materializado. La verificación offline debe terminar sin credenciales ni red de servicio.

## 4. Exact file manifest

```text
CREATE azure_content_understanding_official_binary/README.md
CREATE azure_content_understanding_official_binary/upstream/LICENSE.txt
CREATE azure_content_understanding_official_binary/upstream/sample_analyze_binary.py
CREATE azure_content_understanding_official_binary/upstream/test_sample_analyze_binary.py
CREATE azure_content_understanding_official_binary/source-lock.json
CREATE azure_content_understanding_official_binary/test_official_sample.py
```

## 5. Materialization blocks

### FILE: `azure_content_understanding_official_binary/README.md`
```yaml
block_id: "AZURE-CU-OFFICIAL-BINARY:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local guide constrained by exact Microsoft source"
license: "LicenseRef-Workspace-Owner"
sha256: "eab757a9f1ffe0c01e131b5a6e29f1c5fba520a41598860ed491a21d6870a068"
variables: []
secrets_allowed: false
```
````text
# Microsoft Azure Content Understanding official binary document sample

This directory contains Microsoft's byte-verbatim synchronous binary-document sample, its official SDK test and the MIT license from `azure-ai-contentunderstanding` 1.1.0. The exact commit, sdist and file hashes are fixed in `source-lock.json`. `test_official_sample.py` is Elite-authored packaging verification; it is not Microsoft product code.

The official sample submits local PDF bytes to `prebuilt-documentSearch`, preserves markdown/layout/tables/figures, and demonstrates page ranges. It does not extract a procurement schema and does not authorize persistence. Choose `prebuilt-procurement`, `prebuilt-purchaseOrder`, another documented analyzer or a custom analyzer only through the project's document decision record.

Offline verification:

```powershell
python .\test_official_sample.py
```

Live use requires the exact Microsoft SDK, a frozen and audited project dependency graph, an approved Microsoft Foundry resource, identity/secret authority, configured model deployments, region/quota/cost approval and admitted local documents. Compose secure intake, malware handling, evidence, ground truth, field thresholds, review, reconciliation and storage authorization before business persistence.

The included Microsoft test uses the Azure SDK test-proxy/live harness. It is preserved as official evidence but is not presented as a standalone offline test.
````

### FILE: `azure_content_understanding_official_binary/upstream/LICENSE.txt`
```yaml
block_id: "AZURE-CU-OFFICIAL-BINARY:microsoft-license:v1"
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

### FILE: `azure_content_understanding_official_binary/upstream/sample_analyze_binary.py`
```yaml
block_id: "AZURE-CU-OFFICIAL-BINARY:microsoft-sample:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/Azure/azure-sdk-for-python/129a5cbb06fba44e11ffc92f6bcdca14a2f5b6eb/sdk/contentunderstanding/azure-ai-contentunderstanding/samples/sample_analyze_binary.py"
license: "MIT"
sha256: "dce0d3c2684bb0d5bca01f1015e601d4f15407b8272c72b773e2cfb822588893"
variables: []
secrets_allowed: true
```
````text
# pylint: disable=line-too-long,useless-suppression
# coding=utf-8
# --------------------------------------------------------------------------
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License. See License.txt in the project root for license information.
# --------------------------------------------------------------------------
"""
FILE: sample_analyze_binary.py

DESCRIPTION:
    This sample demonstrates how to analyze a PDF file from disk using the prebuilt-documentSearch
    analyzer.

    ## About analyzing documents from binary data

    One of the key values of Content Understanding is taking a content file and extracting the content
    for you in one call. The service returns an AnalysisResult that contains an array of AnalysisContent
    items in AnalysisResult.contents. This sample starts with a document file, so each item is a
    DocumentContent (a subtype of AnalysisContent) that exposes markdown plus detailed structure such
    as pages, tables, figures, and paragraphs.

    This sample focuses on document analysis. For prebuilt RAG analyzers covering images, audio, and
    video, see sample_analyze_url.py.

    ## Prebuilt analyzers

    Content Understanding provides prebuilt RAG analyzers (the prebuilt-*Search analyzers, such as
    prebuilt-documentSearch) that return markdown and a one-paragraph Summary for each content item,
    making them useful for retrieval-augmented generation (RAG) and other downstream applications:

    - prebuilt-documentSearch - Extracts content from documents (PDF, images, Office documents) with
      layout preservation, table detection, figure analysis, and structured markdown output.
      Optimized for RAG scenarios.
    - prebuilt-audioSearch - Transcribes audio content with speaker diarization, timing information,
      and conversation summaries. Supports multilingual transcription.
    - prebuilt-videoSearch - Analyzes video content with visual frame extraction, audio transcription,
      and structured summaries. Provides temporal alignment of visual and audio content.
    - prebuilt-imageSearch - Analyzes standalone images and returns a one-paragraph Summary of the
      image content. For images that contain text (including hand-written text), use
      prebuilt-documentSearch.

    This sample uses prebuilt-documentSearch to extract structured content from PDF documents.

USAGE:
    python sample_analyze_binary.py

    Set the environment variables with your own values before running the sample:
    1) CONTENTUNDERSTANDING_ENDPOINT - the endpoint to your Content Understanding resource.
    2) CONTENTUNDERSTANDING_KEY - your Content Understanding API key (optional if using DefaultAzureCredential).

    See sample_update_defaults.py for model deployment setup guidance.
"""

import os

from dotenv import load_dotenv
from azure.ai.contentunderstanding import ContentUnderstandingClient
from azure.ai.contentunderstanding.models import (
    AnalysisResult,
    DocumentContent,
)
from azure.core.credentials import AzureKeyCredential
from azure.identity import DefaultAzureCredential

load_dotenv()


def main() -> None:
    endpoint = os.environ["CONTENTUNDERSTANDING_ENDPOINT"]
    key = os.getenv("CONTENTUNDERSTANDING_KEY")
    credential = AzureKeyCredential(key) if key else DefaultAzureCredential()

    client = ContentUnderstandingClient(endpoint=endpoint, credential=credential)

    # [START analyze_document_from_binary]
    # Replace with the path to your local document file.
    file_path = "sample_files/sample_invoice.pdf"

    with open(file_path, "rb") as f:
        file_bytes = f.read()

    print(f"Analyzing {file_path} with prebuilt-documentSearch...")
    poller = client.begin_analyze_binary(
        analyzer_id="prebuilt-documentSearch",
        binary_input=file_bytes,
    )
    result: AnalysisResult = poller.result()
    # [END analyze_document_from_binary]

    # [START analyze_binary_with_content_range]
    # Use a multi-page document for content range demonstrations.
    multi_page_path = "sample_files/mixed_financial_invoices.pdf"
    with open(multi_page_path, "rb") as f:
        multi_page_bytes = f.read()

    # Analyze only pages 3 onward.
    print("\nAnalyzing pages 3 onward with content range '3-'...")
    range_poller = client.begin_analyze_binary(
        analyzer_id="prebuilt-documentSearch",
        binary_input=multi_page_bytes,
        content_range="3-",
    )
    range_result: AnalysisResult = range_poller.result()

    if isinstance(range_result.contents[0], DocumentContent):
        range_doc = range_result.contents[0]
        print(f"Content range analysis returned pages {range_doc.start_page_number} - {range_doc.end_page_number}")
    # [END analyze_binary_with_content_range]

    # [START analyze_binary_with_combined_content_range]
    # Analyze pages 1-3, page 5, and pages 9 onward.
    print("\nAnalyzing combined pages (1-3, 5, 9-) with content range '1-3,5,9-'...")
    combine_range_poller = client.begin_analyze_binary(
        analyzer_id="prebuilt-documentSearch",
        binary_input=multi_page_bytes,
        content_range="1-3,5,9-",
    )
    combine_range_result: AnalysisResult = combine_range_poller.result()

    if isinstance(combine_range_result.contents[0], DocumentContent):
        combine_doc = combine_range_result.contents[0]
        print(f"Combined content range analysis returned pages {combine_doc.start_page_number} - {combine_doc.end_page_number}")
    # [END analyze_binary_with_combined_content_range]

    # [START extract_markdown]
    print("\nMarkdown Content:")
    print("=" * 50)

    # A PDF file has only one content element even if it contains multiple pages
    content = result.contents[0]
    print(content.markdown)

    print("=" * 50)
    # [END extract_markdown]

    # [START access_document_properties]
    # Check if this is document content to access document-specific properties
    if isinstance(content, DocumentContent):
        print(f"\nDocument type: {content.mime_type or '(unknown)'}")
        print(f"Start page: {content.start_page_number}")
        print(f"End page: {content.end_page_number}")

        # Check for pages
        if content.pages and len(content.pages) > 0:
            print(f"\nNumber of pages: {len(content.pages)}")
            for page in content.pages:
                unit = content.unit or "units"
                print(f"  Page {page.page_number}: {page.width} x {page.height} {unit}")

        # Check for tables
        if content.tables and len(content.tables) > 0:
            print(f"\nNumber of tables: {len(content.tables)}")
            table_counter = 1
            for table in content.tables:
                print(
                    f"  Table {table_counter}: {table.row_count} rows x {table.column_count} columns"
                )
                table_counter += 1
    # [END access_document_properties]


if __name__ == "__main__":
    main()
````

### FILE: `azure_content_understanding_official_binary/upstream/test_sample_analyze_binary.py`
```yaml
block_id: "AZURE-CU-OFFICIAL-BINARY:microsoft-test:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/Azure/azure-sdk-for-python/129a5cbb06fba44e11ffc92f6bcdca14a2f5b6eb/sdk/contentunderstanding/azure-ai-contentunderstanding/tests/samples/test_sample_analyze_binary.py"
license: "MIT"
sha256: "88f8cb38d39a351639c835f246da5c443d5c50cc068a363032ab43ee229b276f"
variables: []
secrets_allowed: false
```
````text
# pylint: disable=line-too-long,useless-suppression
# coding: utf-8

# -------------------------------------------------------------------------
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License. See License.txt in the project root for
# license information.
# --------------------------------------------------------------------------

"""
TEST FILE: test_sample_analyze_binary.py

DESCRIPTION:
    These tests validate the sample_analyze_binary.py sample code.
    
    This sample demonstrates how to analyze a PDF file from disk using the `prebuilt-documentSearch`
    analyzer. The service returns an AnalysisResult that contains an array of AnalysisContent items
    in AnalysisResult.contents. For documents, each item is a DocumentContent that exposes markdown
    plus detailed structure such as pages, tables, figures, and paragraphs.
    
    The prebuilt-documentSearch analyzer transforms unstructured documents into structured, machine-
    readable data optimized for RAG scenarios. It extracts rich GitHub Flavored Markdown that preserves
    document structure and can include: structured text, tables (in HTML format), charts and diagrams,
    mathematical formulas, hyperlinks, barcodes, annotations, and page metadata.
    
    Content Understanding supports many document types including PDF, Word, Excel, PowerPoint, images
    (including scanned image files with hand-written text), and more.

USAGE:
    pytest test_sample_analyze_binary.py
"""

import os
import pytest
from devtools_testutils import recorded_by_proxy
from testpreparer import ContentUnderstandingPreparer, ContentUnderstandingClientTestBase
from azure.ai.contentunderstanding.models import DocumentContent


class TestSampleAnalyzeBinary(ContentUnderstandingClientTestBase):
    """Tests for sample_analyze_binary.py"""

    @ContentUnderstandingPreparer()
    @recorded_by_proxy
    def test_sample_analyze_binary(self, contentunderstanding_endpoint: str) -> None:
        """Test analyzing a document from binary data.

        This test validates:
        1. File loading and binary data creation
        2. Document analysis using begin_analyze_binary
        3. Markdown content extraction
        4. Document properties (MIME type, pages, tables)

        """
        client = self.create_client(endpoint=contentunderstanding_endpoint)

        # Read the sample file
        # Use test_data directory from parent tests folder
        tests_dir = os.path.dirname(os.path.dirname(__file__))
        file_path = os.path.join(tests_dir, "test_data", "sample_invoice.pdf")

        # Assertion: Verify file exists
        assert os.path.exists(file_path), f"Sample file not found at {file_path}"
        print(f"[PASS] Sample file exists: {file_path}")

        with open(file_path, "rb") as f:
            file_bytes = f.read()

        # Assertion: Verify file is not empty
        assert len(file_bytes) > 0, "File should not be empty"
        print(f"[PASS] File loaded: {len(file_bytes)} bytes")

        # Assertion: Verify binary data
        assert file_bytes is not None, "Binary data should not be null"
        print("[PASS] Binary data created successfully")

        # Analyze the document
        poller = client.begin_analyze_binary(analyzer_id="prebuilt-documentSearch", binary_input=file_bytes)

        result = poller.result()

        # Assertion: Verify analysis operation completed
        assert poller is not None, "Analysis operation should not be null"
        assert poller.done(), "Operation should be completed"
        print("[PASS] Analysis operation completed successfully")

        # Assertion: Verify result
        assert result is not None, "Analysis result should not be null"
        assert hasattr(result, "contents"), "Result should have contents attribute"
        assert result.contents is not None, "Result contents should not be null"
        print(f"[PASS] Analysis result contains {len(result.contents)} content(s)")

        # Test markdown extraction
        self._test_markdown_extraction(result)

        # Test document properties access
        self._test_document_properties(result)

        print("\n[SUCCESS] All test_sample_analyze_binary assertions passed")

    def _test_markdown_extraction(self, result):
        """Test markdown content extraction."""
        # Assertion: Verify contents structure
        assert result.contents is not None, "Result should contain contents"
        assert len(result.contents) > 0, "Result should have at least one content"
        assert len(result.contents) == 1, "PDF file should have exactly one content element"

        content = result.contents[0]
        assert content is not None, "Content should not be null"

        # Assertion: Verify markdown content
        markdown = getattr(content, "markdown", None)
        if markdown:
            assert isinstance(markdown, str), "Markdown should be a string"
            assert len(markdown) > 0, "Markdown content should not be empty"
            assert markdown.strip(), "Markdown content should not be just whitespace"
            print(f"[PASS] Markdown content extracted successfully ({len(markdown)} characters)")
        else:
            print("[WARN]  No markdown content available")

    def _test_document_properties(self, result):
        """Test document property access."""
        content = result.contents[0]
        assert content is not None, "Content should not be null for document properties validation"

        # Check if this is DocumentContent
        content_type = type(content).__name__
        print(f"[INFO] Content type: {content_type}")

        # Validate this is document content (should have document-specific properties)
        is_document_content = hasattr(content, "mime_type") and hasattr(content, "start_page_number")
        if not is_document_content:
            print(f"[WARN] Expected DocumentContent but got {content_type}, skipping document-specific validations")
            return

        # Validate MIME type
        mime_type = getattr(content, "mime_type", None)
        if mime_type:
            assert isinstance(mime_type, str), "MIME type should be a string"
            assert mime_type.strip(), "MIME type should not be empty"
            assert mime_type == "application/pdf", f"MIME type should be application/pdf, but was {mime_type}"
            print(f"[PASS] MIME type verified: {mime_type}")

        # Validate page numbers
        start_page = getattr(content, "start_page_number", None)
        if start_page is not None:
            assert start_page >= 1, f"Start page should be >= 1, but was {start_page}"

            end_page = getattr(content, "end_page_number", None)
            if end_page is not None:
                assert end_page >= start_page, f"End page {end_page} should be >= start page {start_page}"
                total_pages = end_page - start_page + 1
                assert total_pages > 0, f"Total pages should be positive, but was {total_pages}"
                print(f"[PASS] Page range verified: {start_page} to {end_page} ({total_pages} pages)")

                # Validate pages collection
                pages = getattr(content, "pages", None)
                if pages and len(pages) > 0:
                    assert len(pages) > 0, "Pages collection should not be empty when not null"
                    assert (
                        len(pages) == total_pages
                    ), f"Pages collection count {len(pages)} should match calculated total pages {total_pages}"
                    print(f"[PASS] Pages collection verified: {len(pages)} pages")

                    # Validate individual pages
                    self._validate_pages(pages, start_page, end_page, content)
                else:
                    print("[WARN] No pages collection available in document content")

        # Validate tables collection
        tables = getattr(content, "tables", None)
        if tables and len(tables) > 0:
            self._validate_tables(tables)
        else:
            print("No tables found in document content")

        # Final validation message
        print("[PASS] All document properties validated successfully")

    def _validate_pages(self, pages, start_page, end_page, content=None):
        """Validate pages collection details."""
        page_numbers = set()
        unit = getattr(content, "unit", None) if content else None
        unit_str = str(unit) if unit else "units"

        for page in pages:
            assert page is not None, "Page object should not be null"
            assert hasattr(page, "page_number"), "Page should have page_number attribute"
            assert page.page_number >= 1, f"Page number should be >= 1, but was {page.page_number}"
            assert (
                start_page <= page.page_number <= end_page
            ), f"Page number {page.page_number} should be within document range [{start_page}, {end_page}]"

            assert (
                hasattr(page, "width") and page.width > 0
            ), f"Page {page.page_number} width should be > 0, but was {page.width}"
            assert (
                hasattr(page, "height") and page.height > 0
            ), f"Page {page.page_number} height should be > 0, but was {page.height}"

            # Ensure page numbers are unique
            assert page.page_number not in page_numbers, f"Page number {page.page_number} appears multiple times"
            page_numbers.add(page.page_number)

            # Print page details with unit
            print(f"  Page {page.page_number}: {page.width} x {page.height} {unit_str}")

        print(f"[PASS] All {len(pages)} pages validated successfully")

    def _validate_tables(self, tables):
        """Validate tables collection details."""
        assert len(tables) > 0, "Tables collection should not be empty when not null"
        print(f"[PASS] Tables collection verified: {len(tables)} tables")

        for i, table in enumerate(tables, 1):
            assert table is not None, f"Table {i} should not be null"
            assert hasattr(table, "row_count"), f"Table {i} should have row_count attribute"
            assert hasattr(table, "column_count"), f"Table {i} should have column_count attribute"
            assert table.row_count > 0, f"Table {i} should have at least 1 row, but had {table.row_count}"
            assert table.column_count > 0, f"Table {i} should have at least 1 column, but had {table.column_count}"

            # Validate table cells if available
            if hasattr(table, "cells") and table.cells:
                assert len(table.cells) > 0, f"Table {i} cells collection should not be empty when not null"

                for cell in table.cells:
                    assert cell is not None, "Table cell should not be null"
                    assert hasattr(cell, "row_index"), "Cell should have row_index"
                    assert hasattr(cell, "column_index"), "Cell should have column_index"
                    assert (
                        0 <= cell.row_index < table.row_count
                    ), f"Cell row index {cell.row_index} should be within table row count {table.row_count}"
                    assert (
                        0 <= cell.column_index < table.column_count
                    ), f"Cell column index {cell.column_index} should be within table column count {table.column_count}"

                    if hasattr(cell, "row_span"):
                        assert cell.row_span >= 1, f"Cell row span should be >= 1, but was {cell.row_span}"
                    if hasattr(cell, "column_span"):
                        assert cell.column_span >= 1, f"Cell column span should be >= 1, but was {cell.column_span}"

                print(
                    f"[PASS] Table {i} validated: {table.row_count} rows x {table.column_count} columns ({len(table.cells)} cells)"
                )
            else:
                print(f"[PASS] Table {i} validated: {table.row_count} rows x {table.column_count} columns")

    @ContentUnderstandingPreparer()
    @recorded_by_proxy
    def test_sample_analyze_binary_with_content_range(self, contentunderstanding_endpoint: str) -> None:
        """Test analyzing a document from binary data with content range strings.

        This test validates:
        1. "2" — single page
        2. "1-3" — page range
        3. "1,3-4" — combined page ranges
        4. "3-" — analyze pages 3 onward
        5. "1-3,5,9-" — analyze disjoint page ranges

        01_AnalyzeBinary.AnalyzeBinaryWithPageContentRangesAsync()
        """
        client = self.create_client(endpoint=contentunderstanding_endpoint)

        # Read the sample file (use multi-page document for content range testing)
        tests_dir = os.path.dirname(os.path.dirname(__file__))
        file_path = os.path.join(tests_dir, "test_data", "mixed_financial_invoices.pdf")
        assert os.path.exists(file_path), (
            f"Required multi-page test file not found: {file_path}. "
            "This test requires 'mixed_financial_invoices.pdf' to validate content range scenarios."
        )

        with open(file_path, "rb") as f:
            file_bytes = f.read()

        # Full analysis for comparison
        full_poller = client.begin_analyze_binary(
            analyzer_id="prebuilt-documentSearch", binary_input=file_bytes
        )
        full_result = full_poller.result()
        assert full_result.contents is not None
        full_doc = full_result.contents[0]
        assert isinstance(full_doc, DocumentContent)
        full_page_count = len(full_doc.pages) if full_doc.pages else 0
        print(f"[PASS] Full document: {full_page_count} pages, {len(full_doc.markdown or '')} chars")

        # "2" — single page
        print("\nAnalyzing page 2 only with content range '2'...")
        page2_poller = client.begin_analyze_binary(
            analyzer_id="prebuilt-documentSearch",
            binary_input=file_bytes,
            content_range="2",
        )
        page2_result = page2_poller.result()
        assert page2_result.contents is not None
        page2_doc = page2_result.contents[0]
        assert isinstance(page2_doc, DocumentContent)
        page2_page_count = len(page2_doc.pages) if page2_doc.pages else 0
        assert page2_page_count == 1, f"'2' should return exactly 1 page, got {page2_page_count}"
        assert page2_doc.start_page_number == 2, f"'2' should start at page 2, got {page2_doc.start_page_number}"
        assert page2_doc.end_page_number == 2, f"'2' should end at page 2, got {page2_doc.end_page_number}"
        assert page2_doc.pages[0].page_number == 2, (
            f"'2' page[0].page_number should be 2, got {page2_doc.pages[0].page_number}"
        )
        print(f"[PASS] '2': {page2_page_count} page, page number: {page2_doc.pages[0].page_number}")

        # "1-3" — page range
        print("\nAnalyzing pages 1-3 with content range '1-3'...")
        pages13_poller = client.begin_analyze_binary(
            analyzer_id="prebuilt-documentSearch",
            binary_input=file_bytes,
            content_range="1-3",
        )
        pages13_result = pages13_poller.result()
        assert pages13_result.contents is not None
        pages13_doc = pages13_result.contents[0]
        assert isinstance(pages13_doc, DocumentContent)
        pages13_page_count = len(pages13_doc.pages) if pages13_doc.pages else 0
        assert pages13_page_count == 3, f"'1-3' should return exactly 3 pages, got {pages13_page_count}"
        assert pages13_doc.start_page_number == 1, f"'1-3' should start at page 1, got {pages13_doc.start_page_number}"
        assert pages13_doc.end_page_number == 3, f"'1-3' should end at page 3, got {pages13_doc.end_page_number}"
        actual_pages13 = sorted([p.page_number for p in pages13_doc.pages])
        assert actual_pages13 == [1, 2, 3], (
            f"'1-3' page numbers should be [1, 2, 3], got {actual_pages13}"
        )
        print(f"[PASS] '1-3': {pages13_page_count} pages, page numbers: {actual_pages13}")

        # "1,3-4" — combined page ranges
        print("\nAnalyzing combined pages (1, 3-4) with content range '1,3-4'...")
        combine2_poller = client.begin_analyze_binary(
            analyzer_id="prebuilt-documentSearch",
            binary_input=file_bytes,
            content_range="1,3-4",
        )
        combine2_result = combine2_poller.result()
        assert combine2_result.contents is not None
        combine2_doc = combine2_result.contents[0]
        assert isinstance(combine2_doc, DocumentContent)
        combine2_page_count = len(combine2_doc.pages) if combine2_doc.pages else 0
        assert combine2_page_count == 3, (
            f"'1,3-4' should return exactly 3 pages, got {combine2_page_count}"
        )
        assert combine2_doc.start_page_number == 1, f"'1,3-4' should start at page 1, got {combine2_doc.start_page_number}"
        assert combine2_doc.end_page_number == 4, f"'1,3-4' should end at page 4, got {combine2_doc.end_page_number}"
        actual_combine2_pages = sorted([p.page_number for p in combine2_doc.pages])
        assert actual_combine2_pages == [1, 3, 4], (
            f"'1,3-4' page numbers should be [1, 3, 4], got {actual_combine2_pages}"
        )
        print(f"[PASS] '1,3-4': {combine2_page_count} pages, page numbers: {actual_combine2_pages}")

        # "3-" — pages 3 onward
        print("\nAnalyzing pages 3 onward with content range '3-'...")
        range_poller = client.begin_analyze_binary(
            analyzer_id="prebuilt-documentSearch",
            binary_input=file_bytes,
            content_range="3-",
        )
        range_result = range_poller.result()
        assert range_result.contents is not None
        range_doc = range_result.contents[0]
        assert isinstance(range_doc, DocumentContent)
        expected_range_page_count = full_page_count - 2
        range_page_count = len(range_doc.pages) if range_doc.pages else 0
        assert range_page_count == expected_range_page_count, (
            f"'3-' should return exactly {expected_range_page_count} pages, got {range_page_count}"
        )
        assert range_doc.start_page_number == 3, f"'3-' should start at page 3, got {range_doc.start_page_number}"
        assert range_doc.end_page_number == full_doc.end_page_number, (
            f"'3-' should end at page {full_doc.end_page_number}, got {range_doc.end_page_number}"
        )
        expected_range_pages = list(range(3, full_doc.end_page_number + 1))
        actual_range_pages = sorted([p.page_number for p in range_doc.pages])
        assert actual_range_pages == expected_range_pages, (
            f"'3-' page numbers should be {expected_range_pages}, got {actual_range_pages}"
        )
        print(f"[PASS] '3-': {range_page_count} pages (pages {range_doc.start_page_number}-{range_doc.end_page_number})")

        # "1-3,5,9-" — combined ranges
        print("\nAnalyzing combined pages (1-3, 5, 9-) with content range '1-3,5,9-'...")
        combine_poller = client.begin_analyze_binary(
            analyzer_id="prebuilt-documentSearch",
            binary_input=file_bytes,
            content_range="1-3,5,9-",
        )
        combine_result = combine_poller.result()
        assert combine_result.contents is not None
        combine_doc = combine_result.contents[0]
        assert isinstance(combine_doc, DocumentContent)
        # Expected pages: 1, 2, 3, 5, 9, 10, ..., N => count = 3 + 1 + (N - 8) = N - 4
        expected_combine_page_count = full_page_count - 4
        combine_page_count = len(combine_doc.pages) if combine_doc.pages else 0
        assert combine_page_count == expected_combine_page_count, (
            f"'1-3,5,9-' should return exactly {expected_combine_page_count} pages, got {combine_page_count}"
        )
        expected_combine_pages = [1, 2, 3, 5] + list(range(9, full_doc.end_page_number + 1))
        actual_combine_pages = sorted([p.page_number for p in combine_doc.pages])
        assert actual_combine_pages == expected_combine_pages, (
            f"'1-3,5,9-' page numbers should be {expected_combine_pages}, got {actual_combine_pages}"
        )
        assert combine_doc.start_page_number == 1, (
            f"'1-3,5,9-' should start at page 1, got {combine_doc.start_page_number}"
        )
        assert combine_doc.end_page_number == full_doc.end_page_number, (
            f"'1-3,5,9-' should end at page {full_doc.end_page_number}, got {combine_doc.end_page_number}"
        )
        print(f"[PASS] '1-3,5,9-': {combine_page_count} pages, page numbers: {actual_combine_pages}")

        print("\n[SUCCESS] All content range binary test assertions passed")
````

### FILE: `azure_content_understanding_official_binary/source-lock.json`
```yaml
block_id: "AZURE-CU-OFFICIAL-BINARY:source-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "local acquisition receipt from official Microsoft artifacts"
license: "LicenseRef-Workspace-Owner"
sha256: "aef052805b9879d937a787ba2438497294a1e3f869c323334fa26676827cf927"
variables: []
secrets_allowed: false
```
````text
{
  "schema_version": 1,
  "source": "Azure/azure-sdk-for-python",
  "revision": "129a5cbb06fba44e11ffc92f6bcdca14a2f5b6eb",
  "package": "azure-ai-contentunderstanding",
  "package_version": "1.1.0",
  "sdist": {
    "filename": "azure_ai_contentunderstanding-1.1.0.tar.gz",
    "bytes": 230330,
    "sha256": "00f39c7cf2ba50c8d4586315f4295a45fad0881daaff07d13d33aea2bafc7ec8"
  },
  "files": [
    {
      "path": "upstream/LICENSE.txt",
      "provenance": "VERBATIM",
      "license": "MIT",
      "sha256": "7c77a44a8acd9b41fdc209864a8016b3d430b5d0e09309818d5b7444336df744"
    },
    {
      "path": "upstream/sample_analyze_binary.py",
      "provenance": "VERBATIM",
      "license": "MIT",
      "bytes": 6979,
      "sha256": "dce0d3c2684bb0d5bca01f1015e601d4f15407b8272c72b773e2cfb822588893",
      "commit_equals_sdist": true
    },
    {
      "path": "upstream/test_sample_analyze_binary.py",
      "provenance": "VERBATIM",
      "license": "MIT",
      "bytes": 20749,
      "sha256": "88f8cb38d39a351639c835f246da5c443d5c50cc068a363032ab43ee229b276f",
      "commit_equals_sdist": true
    }
  ],
  "admission": "CONDITIONED",
  "conditions": [
    "prebuilt-documentSearch preserves document content but is not a business-field procurement schema",
    "live execution requires approved Azure access, identity, region, quota, model deployments and cost",
    "project dependencies must be frozen and audited",
    "business persistence requires corpus, ground truth, field evidence, review and reconciliation"
  ]
}
````

### FILE: `azure_content_understanding_official_binary/test_official_sample.py`
```yaml
block_id: "AZURE-CU-OFFICIAL-BINARY:offline-contracts:v1"
operation: CREATE
provenance: AUTHORED
source: "local offline packaging contracts; not Microsoft product code"
license: "LicenseRef-Workspace-Owner"
sha256: "28ed9f9713294e11ef6cf07fc7dda83121df1573782ef7fd1c23088c1a39a6e2"
variables: []
secrets_allowed: false
```
````text
from __future__ import annotations

import ast
import hashlib
import json
from pathlib import Path


ROOT = Path(__file__).resolve().parent
LOCK_PATH = ROOT / "source-lock.json"


def sha256(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def require(condition: bool, message: str) -> None:
    if not condition:
        raise AssertionError(message)


lock = json.loads(LOCK_PATH.read_text(encoding="utf-8"))
require(lock["source"] == "Azure/azure-sdk-for-python", "unexpected source")
require(
    lock["revision"] == "129a5cbb06fba44e11ffc92f6bcdca14a2f5b6eb",
    "unexpected revision",
)
require(lock["package_version"] == "1.1.0", "unexpected package version")
require(
    lock["sdist"]["sha256"]
    == "00f39c7cf2ba50c8d4586315f4295a45fad0881daaff07d13d33aea2bafc7ec8",
    "unexpected sdist hash",
)
print("contract_source_lock PASS")

for entry in lock["files"]:
    path = ROOT / entry["path"]
    require(path.is_file(), f"missing official file: {entry['path']}")
    require(sha256(path) == entry["sha256"], f"hash mismatch: {entry['path']}")
print("contract_official_hashes PASS")

sample_path = ROOT / "upstream" / "sample_analyze_binary.py"
official_test_path = ROOT / "upstream" / "test_sample_analyze_binary.py"
sample_tree = ast.parse(sample_path.read_text(encoding="utf-8"))
official_test_tree = ast.parse(official_test_path.read_text(encoding="utf-8"))
require(isinstance(sample_tree, ast.Module), "sample AST missing")
require(isinstance(official_test_tree, ast.Module), "official test AST missing")
print("contract_python_ast PASS")

calls = [
    node
    for node in ast.walk(sample_tree)
    if isinstance(node, ast.Call)
    and isinstance(node.func, ast.Attribute)
    and node.func.attr == "begin_analyze_binary"
]
require(len(calls) == 3, f"expected 3 binary analyze calls, found {len(calls)}")
for call in calls:
    keywords = {item.arg: item.value for item in call.keywords if item.arg}
    analyzer = keywords.get("analyzer_id")
    require(
        isinstance(analyzer, ast.Constant)
        and analyzer.value == "prebuilt-documentSearch",
        "unexpected analyzer id",
    )
    require("binary_input" in keywords, "binary_input keyword missing")
ranges = sorted(
    keyword.value.value
    for call in calls
    for keyword in call.keywords
    if keyword.arg == "content_range"
    and isinstance(keyword.value, ast.Constant)
    and isinstance(keyword.value.value, str)
)
require(ranges == ["1-3,5,9-", "3-"], f"unexpected content ranges: {ranges}")
print("contract_binary_analysis PASS")

official_test_text = official_test_path.read_text(encoding="utf-8")
require("class TestSampleAnalyzeBinary" in official_test_text, "official class missing")
require(official_test_text.count("@recorded_by_proxy") == 2, "proxy test count changed")
require(
    official_test_text.count('analyzer_id="prebuilt-documentSearch"') >= 7,
    "official analyzer assertions changed",
)
print("contract_official_test_surface PASS")

print("OFFICIAL_AZURE_BINARY_SAMPLE_TEST_PASS contracts=5")
````

## 6. Configuration surface

| Variable | Required live | Secret | Authority |
|---|---:|---:|---|
| `CONTENTUNDERSTANDING_ENDPOINT` | yes | no | Microsoft Foundry resource receipt |
| `CONTENTUNDERSTANDING_KEY` | conditional | yes | approved secret store; omit for `DefaultAzureCredential` |
| local input path | yes | no | secure local ingestion receipt |
| analyzer selection | yes | no | official analyzer inventory plus project decision record |

El sample exact conserva paths demostrativos relativos. El agente debe suministrar archivos admitidos en el árbol del proyecto sin editar el upstream y registrar analyzer/API/version/region/costo en el record del proyecto.

## 7. Dependency bill

| Dependency | Pin | Purpose | Condition |
|---|---|---|---|
| `azure-ai-contentunderstanding` | `1.1.0` | client/LRO/models | artefacto oficial ya fijado por hash |
| `azure-core` | graph oficial de 1.1.0 | credential/pipeline | congelar en proyecto |
| `azure-identity` | proyecto | Entra ID | congelar y auditar sólo si se usa |
| `python-dotenv` | proyecto | demo env loading | reemplazable por config admitida; congelar si se conserva |
| Microsoft test proxy | SDK repo | test oficial | no se presenta como suite standalone |

## 8. Apply order

1. Materializar en destino vacío y ejecutar contratos offline.
2. Componer intake local seguro, routing y evaluación estricta.
3. Elegir analyzer oficial y fijar su definición/copia cuando Microsoft lo recomiende.
4. Congelar el grafo target, auditar vulnerabilidades/licencias y registrar acceso/costo.
5. Ejecutar sólo sobre corpus no productivo; medir fields/evidence y revisión.
6. Autorizar persistencia únicamente después de los gates del proyecto.

## 9. Verification

```powershell
pwsh -NoProfile -File .\materialize_markdown_pack.ps1 -PackFile .\implementation_packs\MICROSOFT_AZURE_CONTENT_UNDERSTANDING_OFFICIAL_BINARY_DOCUMENT_SAMPLE.md -Destination <empty-dir>
python <empty-dir>\azure_content_understanding_official_binary\test_official_sample.py
python -m py_compile <empty-dir>\azure_content_understanding_official_binary\upstream\sample_analyze_binary.py <empty-dir>\azure_content_understanding_official_binary\upstream\test_sample_analyze_binary.py
```

Esperado: materialización 6/6, hashes exactos, contratos offline completos y sintaxis Python limpia. El test Microsoft requiere su harness oficial y acceso/recording; incluirlo no equivale a ejecutarlo standalone ni a una prueba live.

## 10. Reconstruction evidence

- commit Microsoft: `129a5cbb06fba44e11ffc92f6bcdca14a2f5b6eb`;
- sdist PyPI 1.1.0: 230.330 bytes, SHA-256 `00f39c7cf2ba50c8d4586315f4295a45fad0881daaff07d13d33aea2bafc7ec8`;
- sample: 6.979 bytes, SHA-256 `dce0d3c2684bb0d5bca01f1015e601d4f15407b8272c72b773e2cfb822588893`, commit=sdist;
- test oficial: 20.749 bytes, SHA-256 `88f8cb38d39a351639c835f246da5c443d5c50cc068a363032ab43ee229b276f`, commit=sdist;
- licencia MIT Microsoft: SHA-256 `7c77a44a8acd9b41fdc209864a8016b3d430b5d0e09309818d5b7444336df744`.
