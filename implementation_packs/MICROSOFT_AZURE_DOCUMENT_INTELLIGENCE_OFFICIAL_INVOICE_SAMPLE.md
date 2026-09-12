# Microsoft Azure Document Intelligence Official Invoice Sample

## 1. Metadata

```yaml
pack_id: "MICROSOFT-AZURE-DOCUMENT-INTELLIGENCE-OFFICIAL-INVOICE-SAMPLE"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa sin modificar el sample Python oficial de Microsoft Azure Document Intelligence 1.0.2 para analizar una invoice desde bytes, su licencia MIT, un lock del fixture oficial fijado por commit y un adquiridor fail-closed con aprobación, hashes, staging, receipt y regresiones offline."
stacks: ["PowerShell 7", "Python 3.8+", "azure-ai-documentintelligence 1.0.2"]
compatible_with: ["OFFICIAL-DOCUMENT-SDK-ARTIFACT-CORE 0.3.x", "SECURE-LOCAL-FILE-INGESTION-GATE 0.1.x", "DOCUMENT-EVIDENCE-ENVELOPE-CORE 0.1.x"]
incompatible_with: ["persistencia automática", "fixture no fijado", "secretos embebidos", "modelo distinto de prebuilt-invoice sin nueva evidencia"]
license_expression: "LicenseRef-Workspace-Owner AND MIT"
upstream_sources: ["https://github.com/Azure/azure-sdk-for-python/tree/8555d14532a9688b751d8408d822d1dd5feb47f6/sdk/documentintelligence/azure-ai-documentintelligence/samples", "https://pypi.org/project/azure-ai-documentintelligence/1.0.2/"]
verified_at: "2026-08-28"
```

## 2. Applicability

Use este pack cuando el blueprint haya elegido Azure Document Intelligence y necesite partir del sample de invoice desde bytes publicado por Microsoft, sin reescribir su llamada al SDK. Requiere adquirir el SDK exacto por el artifact core, una cuenta/endpoint/key aportados por el proyecto para la prueba live y aceptación explícita de MIT para el fixture oficial.

No use el sample como autorización para persistir resultados, como garantía de exactitud de campos ni como cobertura de proformas, packing lists u otras clases. Es una referencia ejecutable y verificable de integración para una invoice; el proyecto debe añadir corpus, schema, evaluación, evidencia por campo, reconciliación y revisión según sus contratos documentales.

## 3. Architecture contract

El sample y `upstream/LICENSE.txt` son bytes `VERBATIM` del repositorio Microsoft fijado al commit firmado `8555d14532a9688b751d8408d822d1dd5feb47f6`. `fixture-lock.json` identifica el JPEG oficial por URL commit-pinned, tamaño y SHA-256. El adquiridor exige lock y approval con propiedades exactas, IDs únicos, aceptación MIT, selección idéntica, host `raw.githubusercontent.com`, ruta segura, destino inexistente, staging y verificación de bytes/hash antes del movimiento final.

`test_official_invoice_sample.py` ejecuta la función real del sample con módulos del SDK sustituidos sólo en el test y prueba que abre el fixture, llama `prebuilt-invoice`, usa locale `en-US` y el endpoint esperado. `test_fixture_acquisition.ps1` prueba validación, adquisición offline, receipt, hash y rechazos. Ningún archivo descargado se incorpora al pack; el fixture conserva licencia Microsoft MIT.

Presupuesto: validación y tests offline menores a 10 segundos; descarga limitada a 184.686 bytes más overhead HTTP. Rollback: eliminar únicamente el destino nuevo después de preservar receipt/evidencia. Una actualización exige nuevo commit/tag, hashes, diff upstream, pruebas y versión del pack.

## 4. Exact file manifest

```text
CREATE azure_document_intelligence_official_invoice/README.md
CREATE azure_document_intelligence_official_invoice/upstream/LICENSE.txt
CREATE azure_document_intelligence_official_invoice/upstream/sample_analyze_invoices_from_bytes_source.py
CREATE azure_document_intelligence_official_invoice/fixture-lock.json
CREATE azure_document_intelligence_official_invoice/fixture-approval.template.json
CREATE azure_document_intelligence_official_invoice/acquire_official_fixture.ps1
CREATE azure_document_intelligence_official_invoice/test_official_invoice_sample.py
CREATE azure_document_intelligence_official_invoice/test_fixture_acquisition.ps1
```

## 5. Materialization blocks

### FILE: `azure_document_intelligence_official_invoice/README.md`
```yaml
block_id: "AZURE-DI-OFFICIAL-INVOICE:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration guide constrained by Microsoft upstream"
license: "LicenseRef-Workspace-Owner"
sha256: "f0e21349e4a939cf85ac07f415da269b6d9310fc5d76fe509f92c68c63a0dc7d"
variables: []
secrets_allowed: false
```
````markdown
# Microsoft Azure Document Intelligence official invoice sample

This directory materializes Microsoft's exact Azure AI Document Intelligence 1.0.2 invoice-from-bytes sample. The Python sample and MIT license are byte-verbatim from the signed release commit 8555d14532a9688b751d8408d822d1dd5feb47f6. The acquisition wrapper, lock, approval template and tests are Elite-authored packaging; they are not Microsoft product code.

Use order:

1. Validate fixture-lock.json for azure-documentintelligence-1.0.2-sample-invoice-jpg.
2. Copy fixture-approval.template.json and bind it to the exact SHA-256 of fixture-lock.json.
3. Acquire the exact Microsoft fixture to a new azure_di_official_invoice_sample/upstream directory, from the official raw URL or a verified offline cache.
4. Install the exact azure-ai-documentintelligence 1.0.2 wheel through OFFICIAL-DOCUMENT-SDK-ARTIFACT-CORE. Bootstrap source builds only with a separately verified pip 26.2 or a newer scanned version.
5. Set DOCUMENTINTELLIGENCE_ENDPOINT and DOCUMENTINTELLIGENCE_API_KEY through the project's secret mechanism, then run upstream/sample_analyze_invoices_from_bytes_source.py.

The sample calls prebuilt-invoice and prints fields/confidences. It does not perform malware admission, validate a business schema, compare ground truth, authorize persistence, write durable evidence or reconcile provider effects. It is an executable Microsoft integration sample only after the user supplies an Azure resource and accepts possible service cost. Route production documents through the Elite security, routing, evaluation, review, persistence and evidence gates instead of treating printed fields as trusted business data.
````

### FILE: `azure_document_intelligence_official_invoice/upstream/LICENSE.txt`
```yaml
block_id: "AZURE-DI-OFFICIAL-INVOICE:microsoft-license:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/Azure/azure-sdk-for-python/8555d14532a9688b751d8408d822d1dd5feb47f6/LICENSE"
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

### FILE: `azure_document_intelligence_official_invoice/upstream/sample_analyze_invoices_from_bytes_source.py`
```yaml
block_id: "AZURE-DI-OFFICIAL-INVOICE:microsoft-sample:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/Azure/azure-sdk-for-python/8555d14532a9688b751d8408d822d1dd5feb47f6/sdk/documentintelligence/azure-ai-documentintelligence/samples/sample_analyze_invoices_from_bytes_source.py"
license: "MIT"
sha256: "ebc32b0fc534e625b15e636c6b55376f01d28a83f35217dd141202073297a293"
variables: []
secrets_allowed: true
```
````python
# coding: utf-8

# -------------------------------------------------------------------------
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License. See License.txt in the project root for
# license information.
# --------------------------------------------------------------------------

"""
FILE: sample_analyze_invoices_from_bytes_source.py

DESCRIPTION:
    This sample demonstrates how to analyze invoices from bytes source.

    See fields found on a invoice here:
    https://aka.ms/azsdk/documentintelligence/invoicefieldschema

USAGE:
    python sample_analyze_invoices_from_bytes_source.py

    Set the environment variables with your own values before running the sample:
    1) DOCUMENTINTELLIGENCE_ENDPOINT - the endpoint to your Document Intelligence resource.
    2) DOCUMENTINTELLIGENCE_API_KEY - your Document Intelligence API key.
"""

import os


def analyze_invoice():
    path_to_sample_documents = os.path.abspath(
        os.path.join(
            os.path.abspath(__file__),
            "..",
            "./sample_forms/forms/sample_invoice.jpg",
        )
    )

    from azure.core.credentials import AzureKeyCredential
    from azure.ai.documentintelligence import DocumentIntelligenceClient
    from azure.ai.documentintelligence.models import AnalyzeResult, AnalyzeDocumentRequest

    endpoint = os.environ["DOCUMENTINTELLIGENCE_ENDPOINT"]
    key = os.environ["DOCUMENTINTELLIGENCE_API_KEY"]

    document_intelligence_client = DocumentIntelligenceClient(endpoint=endpoint, credential=AzureKeyCredential(key))
    with open(path_to_sample_documents, "rb") as f:
        poller = document_intelligence_client.begin_analyze_document(
            "prebuilt-invoice", AnalyzeDocumentRequest(bytes_source=f.read()), locale="en-US"
        )
    invoices: AnalyzeResult = poller.result()

    if invoices.documents:
        for idx, invoice in enumerate(invoices.documents):
            print(f"--------Analyzing invoice #{idx + 1}--------")
            if invoice.fields:
                vendor_name = invoice.fields.get("VendorName")
                if vendor_name:
                    print(f"Vendor Name: {vendor_name.get('content')} has confidence: {vendor_name.get('confidence')}")
                vendor_address = invoice.fields.get("VendorAddress")
                if vendor_address:
                    print(
                        f"Vendor Address: {vendor_address.get('content')} has confidence: {vendor_address.get('confidence')}"
                    )
                vendor_address_recipient = invoice.fields.get("VendorAddressRecipient")
                if vendor_address_recipient:
                    print(
                        f"Vendor Address Recipient: {vendor_address_recipient.get('content')} has confidence: {vendor_address_recipient.get('confidence')}"
                    )
                customer_name = invoice.fields.get("CustomerName")
                if customer_name:
                    print(
                        f"Customer Name: {customer_name.get('content')} has confidence: {customer_name.get('confidence')}"
                    )
                customer_id = invoice.fields.get("CustomerId")
                if customer_id:
                    print(f"Customer Id: {customer_id.get('content')} has confidence: {customer_id.get('confidence')}")
                customer_address = invoice.fields.get("CustomerAddress")
                if customer_address:
                    print(
                        f"Customer Address: {customer_address.get('content')} has confidence: {customer_address.get('confidence')}"
                    )
                customer_address_recipient = invoice.fields.get("CustomerAddressRecipient")
                if customer_address_recipient:
                    print(
                        f"Customer Address Recipient: {customer_address_recipient.get('content')} has confidence: {customer_address_recipient.get('confidence')}"
                    )
                invoice_id = invoice.fields.get("InvoiceId")
                if invoice_id:
                    print(f"Invoice Id: {invoice_id.get('content')} has confidence: {invoice_id.get('confidence')}")
                invoice_date = invoice.fields.get("InvoiceDate")
                if invoice_date:
                    print(
                        f"Invoice Date: {invoice_date.get('content')} has confidence: {invoice_date.get('confidence')}"
                    )
                invoice_total = invoice.fields.get("InvoiceTotal")
                if invoice_total:
                    print(
                        f"Invoice Total: {invoice_total.get('content')} has confidence: {invoice_total.get('confidence')}"
                    )
                due_date = invoice.fields.get("DueDate")
                if due_date:
                    print(f"Due Date: {due_date.get('content')} has confidence: {due_date.get('confidence')}")
                purchase_order = invoice.fields.get("PurchaseOrder")
                if purchase_order:
                    print(
                        f"Purchase Order: {purchase_order.get('content')} has confidence: {purchase_order.get('confidence')}"
                    )
                billing_address = invoice.fields.get("BillingAddress")
                if billing_address:
                    print(
                        f"Billing Address: {billing_address.get('content')} has confidence: {billing_address.get('confidence')}"
                    )
                billing_address_recipient = invoice.fields.get("BillingAddressRecipient")
                if billing_address_recipient:
                    print(
                        f"Billing Address Recipient: {billing_address_recipient.get('content')} has confidence: {billing_address_recipient.get('confidence')}"
                    )
                shipping_address = invoice.fields.get("ShippingAddress")
                if shipping_address:
                    print(
                        f"Shipping Address: {shipping_address.get('content')} has confidence: {shipping_address.get('confidence')}"
                    )
                shipping_address_recipient = invoice.fields.get("ShippingAddressRecipient")
                if shipping_address_recipient:
                    print(
                        f"Shipping Address Recipient: {shipping_address_recipient.get('content')} has confidence: {shipping_address_recipient.get('confidence')}"
                    )
                print("Invoice items:")
                items = invoice.fields.get("Items")
                if items:
                    for idx, item in enumerate(items.get("valueArray")):
                        print(f"...Item #{idx + 1}")
                        item_description = item.get("valueObject").get("Description")
                        if item_description:
                            print(
                                f"......Description: {item_description.get('content')} has confidence: {item_description.get('confidence')}"
                            )
                        item_quantity = item.get("valueObject").get("Quantity")
                        if item_quantity:
                            print(
                                f"......Quantity: {item_quantity.get('content')} has confidence: {item_quantity.get('confidence')}"
                            )
                        unit = item.get("valueObject").get("Unit")
                        if unit:
                            print(f"......Unit: {unit.get('content')} has confidence: {unit.get('confidence')}")
                        unit_price = item.get("valueObject").get("UnitPrice")
                        if unit_price:
                            unit_price_code = (
                                unit_price.get("valueCurrency").get("currencyCode")
                                if unit_price.get("valueCurrency").get("currencyCode")
                                else ""
                            )
                            print(
                                f"......Unit Price: {unit_price.get('content')}{unit_price_code} has confidence: {unit_price.get('confidence')}"
                            )
                        product_code = item.get("valueObject").get("ProductCode")
                        if product_code:
                            print(
                                f"......Product Code: {product_code.get('content')} has confidence: {product_code.get('confidence')}"
                            )
                        item_date = item.get("valueObject").get("Date")
                        if item_date:
                            print(
                                f"......Date: {item_date.get('content')} has confidence: {item_date.get('confidence')}"
                            )
                        tax = item.get("valueObject").get("Tax")
                        if tax:
                            print(f"......Tax: {tax.get('content')} has confidence: {tax.get('confidence')}")
                        amount = item.get("valueObject").get("Amount")
                        if amount:
                            print(f"......Amount: {amount.get('content')} has confidence: {amount.get('confidence')}")
                subtotal = invoice.fields.get("SubTotal")
                if subtotal:
                    print(f"Subtotal: {subtotal.get('content')} has confidence: {subtotal.get('confidence')}")
                total_tax = invoice.fields.get("TotalTax")
                if total_tax:
                    print(f"Total Tax: {total_tax.get('content')} has confidence: {total_tax.get('confidence')}")
                previous_unpaid_balance = invoice.fields.get("PreviousUnpaidBalance")
                if previous_unpaid_balance:
                    print(
                        f"Previous Unpaid Balance: {previous_unpaid_balance.get('content')} has confidence: {previous_unpaid_balance.get('confidence')}"
                    )
                amount_due = invoice.fields.get("AmountDue")
                if amount_due:
                    print(f"Amount Due: {amount_due.get('content')} has confidence: {amount_due.get('confidence')}")
                service_start_date = invoice.fields.get("ServiceStartDate")
                if service_start_date:
                    print(
                        f"Service Start Date: {service_start_date.get('content')} has confidence: {service_start_date.get('confidence')}"
                    )
                service_end_date = invoice.fields.get("ServiceEndDate")
                if service_end_date:
                    print(
                        f"Service End Date: {service_end_date.get('content')} has confidence: {service_end_date.get('confidence')}"
                    )
                service_address = invoice.fields.get("ServiceAddress")
                if service_address:
                    print(
                        f"Service Address: {service_address.get('content')} has confidence: {service_address.get('confidence')}"
                    )
                service_address_recipient = invoice.fields.get("ServiceAddressRecipient")
                if service_address_recipient:
                    print(
                        f"Service Address Recipient: {service_address_recipient.get('content')} has confidence: {service_address_recipient.get('confidence')}"
                    )
                remittance_address = invoice.fields.get("RemittanceAddress")
                if remittance_address:
                    print(
                        f"Remittance Address: {remittance_address.get('content')} has confidence: {remittance_address.get('confidence')}"
                    )
                remittance_address_recipient = invoice.fields.get("RemittanceAddressRecipient")
                if remittance_address_recipient:
                    print(
                        f"Remittance Address Recipient: {remittance_address_recipient.get('content')} has confidence: {remittance_address_recipient.get('confidence')}"
                    )


if __name__ == "__main__":
    from azure.core.exceptions import HttpResponseError
    from dotenv import find_dotenv, load_dotenv

    try:
        load_dotenv(find_dotenv())
        analyze_invoice()
    except HttpResponseError as error:
        # Examples of how to check an HttpResponseError
        # Check by error code:
        if error.error is not None:
            if error.error.code == "InvalidImage":
                print(f"Received an invalid image error: {error.error}")
            if error.error.code == "InvalidRequest":
                print(f"Received an invalid request error: {error.error}")
            # Raise the error again after printing it
            raise
        # If the inner error is None and then it is possible to check the message to get more information:
        if "Invalid request".casefold() in error.message.casefold():
            print(f"Uh-oh! Seems there was an invalid request: {error}")
        # Raise the error again
        raise
````

### FILE: `azure_document_intelligence_official_invoice/fixture-lock.json`
```yaml
block_id: "AZURE-DI-OFFICIAL-INVOICE:fixture-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "registry derived from Microsoft commit-pinned fixture"
license: "LicenseRef-Workspace-Owner"
sha256: "ac4893f44b79dd50012ee6fc48fa0347a47bf81b0e8cb96408fb524c4dfa8c83"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-azure-documentintelligence-official-fixture-lock/v1",
  "verified_at": "2026-08-28",
  "source": {
    "owner": "Microsoft",
    "repository": "Azure/azure-sdk-for-python",
    "release": "azure-ai-documentintelligence_1.0.2",
    "commit": "8555d14532a9688b751d8408d822d1dd5feb47f6",
    "commit_verified": true,
    "license_expression": "MIT",
    "license_sha256": "7c77a44a8acd9b41fdc209864a8016b3d430b5d0e09309818d5b7444336df744"
  },
  "fixtures": [
    {
      "id": "azure-documentintelligence-1.0.2-sample-invoice-jpg",
      "url": "https://raw.githubusercontent.com/Azure/azure-sdk-for-python/8555d14532a9688b751d8408d822d1dd5feb47f6/sdk/documentintelligence/azure-ai-documentintelligence/samples/sample_forms/forms/sample_invoice.jpg",
      "filename": "sample_invoice.jpg",
      "destination_path": "sample_forms/forms/sample_invoice.jpg",
      "bytes": 184686,
      "sha256": "489f0c63b6a05e0ae0fcfbf29121299fc2ccacab35d5cd40a5d63d42764078cb"
    }
  ]
}
````

### FILE: `azure_document_intelligence_official_invoice/fixture-approval.template.json`
```yaml
block_id: "AZURE-DI-OFFICIAL-INVOICE:approval-template:v1"
operation: CREATE
provenance: AUTHORED
source: "local approval template"
license: "LicenseRef-Workspace-Owner"
sha256: "656e285e538feaf05963689e65631c8d0d4a45b1da72eb74aad6a281139e8e33"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-azure-documentintelligence-official-fixture-approval/v1",
  "lock_sha256": "REPLACE_WITH_FIXTURE_LOCK_SHA256",
  "approved_fixture_ids": [],
  "approved_by": "",
  "approved_at": "YYYY-MM-DD",
  "accept_microsoft_mit_license": false
}
````

### FILE: `azure_document_intelligence_official_invoice/acquire_official_fixture.ps1`
```yaml
block_id: "AZURE-DI-OFFICIAL-INVOICE:fixture-acquirer:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed fixture acquisition adapter"
license: "LicenseRef-Workspace-Owner"
sha256: "ad172fe50b9b63bef692d4346586066d4d57ea9c43bd4f02534f7e26308d24a6"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0

[CmdletBinding()]
param(
  [string]$LockPath = (Join-Path $PSScriptRoot 'fixture-lock.json'),
  [Parameter(Mandatory = $true)]
  [string[]]$FixtureId,
  [string]$DestinationRoot,
  [string]$ApprovalPath,
  [string]$ImportDirectory,
  [switch]$AllowNetwork,
  [switch]$ValidateOnly
)

$ErrorActionPreference = 'Stop'
$utf8 = [Text.UTF8Encoding]::new($false)

function Fail([string]$Message) { throw "AZURE_DI_OFFICIAL_FIXTURE_FAILED: $Message" }

function Assert-ExactProperties($Object, [string[]]$Expected, [string]$Label) {
  $actual = @($Object.PSObject.Properties.Name | Sort-Object)
  $wanted = @($Expected | Sort-Object)
  if (($actual -join '|') -cne ($wanted -join '|')) { Fail "$Label properties are not exact" }
}

function Assert-SafeRelativePath([string]$Value, [string]$Label) {
  $normalized = $Value.Replace('\', '/')
  $parts = $normalized.Split('/', [StringSplitOptions]::None)
  if ([string]::IsNullOrWhiteSpace($Value) -or [IO.Path]::IsPathRooted($Value) -or $parts -contains '' -or $parts -contains '.' -or $parts -contains '..') {
    Fail "$Label is unsafe"
  }
  foreach ($part in $parts) {
    if ($part.EndsWith('.') -or $part.EndsWith(' ')) { Fail "$Label is non-canonical" }
  }
  $normalized
}

if (-not (Test-Path -LiteralPath $LockPath -PathType Leaf)) { Fail 'lock not found' }
$lockBytes = [IO.File]::ReadAllBytes((Resolve-Path -LiteralPath $LockPath).Path)
$lockHash = [Convert]::ToHexString([Security.Cryptography.SHA256]::HashData($lockBytes)).ToLowerInvariant()
$lock = [Text.Encoding]::UTF8.GetString($lockBytes) | ConvertFrom-Json -Depth 20
Assert-ExactProperties $lock @('schema','verified_at','source','fixtures') 'lock'
if ($lock.schema -cne 'elite-azure-documentintelligence-official-fixture-lock/v1') { Fail 'unsupported lock schema' }
if ($lock.verified_at -cnotmatch '^[0-9]{4}-[0-9]{2}-[0-9]{2}$') { Fail 'invalid verified_at' }
Assert-ExactProperties $lock.source @('owner','repository','release','commit','commit_verified','license_expression','license_sha256') 'source'
if ($lock.source.owner -cne 'Microsoft' -or $lock.source.repository -cne 'Azure/azure-sdk-for-python' -or $lock.source.release -cne 'azure-ai-documentintelligence_1.0.2') { Fail 'unexpected source identity' }
if ($lock.source.commit -cnotmatch '^[0-9a-f]{40}$' -or $lock.source.commit_verified -cne $true) { Fail 'unverified source commit' }
if ($lock.source.license_expression -cne 'MIT' -or $lock.source.license_sha256 -cnotmatch '^[0-9a-f]{64}$') { Fail 'invalid license identity' }

$fixtureIds = @($lock.fixtures | ForEach-Object { $_.id })
if ($fixtureIds.Count -eq 0 -or ($fixtureIds | Sort-Object -Unique).Count -ne $fixtureIds.Count) { Fail 'fixture IDs must be non-empty and unique' }
$selectedIds = @($FixtureId)
if ($selectedIds.Count -eq 0 -or ($selectedIds | Sort-Object -Unique).Count -ne $selectedIds.Count) { Fail 'selected fixture IDs must be non-empty and unique' }
$selected = @()
foreach ($id in $selectedIds) {
  $fixture = @($lock.fixtures | Where-Object id -ceq $id)
  if ($fixture.Count -ne 1) { Fail "unknown fixture id: $id" }
  Assert-ExactProperties $fixture[0] @('id','url','filename','destination_path','bytes','sha256') "fixture $id"
  $uri = [Uri]$fixture[0].url
  if ($uri.Scheme -cne 'https' -or $uri.Host -cne 'raw.githubusercontent.com') { Fail "untrusted fixture host: $id" }
  if (-not $uri.AbsolutePath.Contains("/$($lock.source.commit)/", [StringComparison]::Ordinal)) { Fail "fixture URL is not commit pinned: $id" }
  if ($fixture[0].filename -cne [IO.Path]::GetFileName($uri.AbsolutePath)) { Fail "fixture filename mismatch: $id" }
  $null = Assert-SafeRelativePath $fixture[0].destination_path "fixture destination $id"
  if ([long]$fixture[0].bytes -le 0 -or $fixture[0].sha256 -cnotmatch '^[0-9a-f]{64}$') { Fail "invalid fixture integrity: $id" }
  $selected += $fixture[0]
}

if ($ValidateOnly) {
  Write-Output "AZURE_DI_OFFICIAL_FIXTURE_LOCK_VALID fixtures=$($lock.fixtures.Count) selected=$($selected.Count)"
  exit 0
}
if ([string]::IsNullOrWhiteSpace($DestinationRoot)) { Fail 'DestinationRoot is required' }
if (Test-Path -LiteralPath $DestinationRoot) { Fail 'destination already exists' }
if (-not $AllowNetwork -and [string]::IsNullOrWhiteSpace($ImportDirectory)) { Fail 'offline import or explicit network is required' }
if ($ImportDirectory -and -not (Test-Path -LiteralPath $ImportDirectory -PathType Container)) { Fail 'import directory not found' }
if (-not $ApprovalPath -or -not (Test-Path -LiteralPath $ApprovalPath -PathType Leaf)) { Fail 'approval not found' }
$approval = Get-Content -Raw -LiteralPath $ApprovalPath | ConvertFrom-Json -Depth 10
Assert-ExactProperties $approval @('schema','lock_sha256','approved_fixture_ids','approved_by','approved_at','accept_microsoft_mit_license') 'approval'
if ($approval.schema -cne 'elite-azure-documentintelligence-official-fixture-approval/v1') { Fail 'unsupported approval schema' }
if ($approval.lock_sha256 -cne $lockHash) { Fail 'approval lock hash mismatch' }
if ([string]::IsNullOrWhiteSpace($approval.approved_by) -or $approval.approved_at -cnotmatch '^[0-9]{4}-[0-9]{2}-[0-9]{2}$') { Fail 'approval identity/date missing' }
if ($approval.accept_microsoft_mit_license -cne $true) { Fail 'Microsoft MIT license not accepted' }
$approvedIds = @($approval.approved_fixture_ids)
$approvedSet = (($approvedIds | Sort-Object) -join '|')
$selectedSet = (($selectedIds | Sort-Object) -join '|')
if ($approvedSet -cne $selectedSet) { Fail 'approval fixture set mismatch' }

$destination = [IO.Path]::GetFullPath($DestinationRoot)
$parent = Split-Path -Parent $destination
if (-not (Test-Path -LiteralPath $parent -PathType Container)) { Fail 'destination parent not found' }
$staging = Join-Path $parent ('.azure-di-fixture-stage-' + [guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Path $staging | Out-Null
try {
  $receiptItems = @()
  foreach ($fixture in $selected) {
    $download = Join-Path $staging $fixture.filename
    if ($ImportDirectory) {
      $cached = Join-Path ([IO.Path]::GetFullPath($ImportDirectory)) $fixture.filename
      if (-not (Test-Path -LiteralPath $cached -PathType Leaf)) { Fail "cached fixture missing: $($fixture.id)" }
      Copy-Item -LiteralPath $cached -Destination $download
    } else {
      Invoke-WebRequest -UseBasicParsing -Uri $fixture.url -OutFile $download
    }
    $item = Get-Item -LiteralPath $download
    $hash = (Get-FileHash -LiteralPath $download -Algorithm SHA256).Hash.ToLowerInvariant()
    if ($item.Length -ne [long]$fixture.bytes -or $hash -cne $fixture.sha256) { Fail "fixture bytes/hash mismatch: $($fixture.id)" }
    $relative = Assert-SafeRelativePath $fixture.destination_path "fixture destination $($fixture.id)"
    $target = [IO.Path]::GetFullPath((Join-Path $staging $relative))
    if (-not $target.StartsWith($staging + [IO.Path]::DirectorySeparatorChar, [StringComparison]::OrdinalIgnoreCase)) { Fail 'fixture escaped staging' }
    $targetParent = Split-Path -Parent $target
    New-Item -ItemType Directory -Path $targetParent -Force | Out-Null
    Move-Item -LiteralPath $download -Destination $target
    $receiptItems += [ordered]@{ id=$fixture.id; destination_path=$relative; bytes=[long]$fixture.bytes; sha256=$fixture.sha256; source=$fixture.url }
  }
  $receipt = [ordered]@{ schema='elite-azure-documentintelligence-official-fixture-receipt/v1'; acquired_at=[DateTimeOffset]::UtcNow.ToString('o'); lock_sha256=$lockHash; approved_by=$approval.approved_by; fixtures=$receiptItems }
  [IO.File]::WriteAllText((Join-Path $staging 'FIXTURE_RECEIPT.json'), ($receipt | ConvertTo-Json -Depth 10) + [Environment]::NewLine, $utf8)
  Move-Item -LiteralPath $staging -Destination $destination
  Write-Output "AZURE_DI_OFFICIAL_FIXTURE_ACQUISITION_PASS fixtures=$($selected.Count)"
} catch {
  if (Test-Path -LiteralPath $staging) { Remove-Item -LiteralPath $staging -Recurse -Force }
  throw
}
````

### FILE: `azure_document_intelligence_official_invoice/test_official_invoice_sample.py`
```yaml
block_id: "AZURE-DI-OFFICIAL-INVOICE:sample-regression:v1"
operation: CREATE
provenance: AUTHORED
source: "local offline regression over verbatim Microsoft sample"
license: "LicenseRef-Workspace-Owner"
sha256: "1c3f29cbbde23300a4b4a97b212986917a64c8669f1f7770e5bbc05946aa7c1b"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import hashlib
import importlib.util
import json
import os
from pathlib import Path
import shutil
import sys
import tempfile
import types
import unittest


ROOT = Path(__file__).resolve().parent
SAMPLE = ROOT / "upstream" / "sample_analyze_invoices_from_bytes_source.py"
LICENSE = ROOT / "upstream" / "LICENSE.txt"
LOCK = ROOT / "fixture-lock.json"
EXPECTED_SAMPLE_SHA256 = "ebc32b0fc534e625b15e636c6b55376f01d28a83f35217dd141202073297a293"
EXPECTED_LICENSE_SHA256 = "7c77a44a8acd9b41fdc209864a8016b3d430b5d0e09309818d5b7444336df744"
EXPECTED_FIXTURE_SHA256 = "489f0c63b6a05e0ae0fcfbf29121299fc2ccacab35d5cd40a5d63d42764078cb"


class OfficialInvoiceSampleTests(unittest.TestCase):
    def test_verbatim_hashes_and_fixture_lock(self) -> None:
        self.assertEqual(hashlib.sha256(SAMPLE.read_bytes()).hexdigest(), EXPECTED_SAMPLE_SHA256)
        self.assertEqual(hashlib.sha256(LICENSE.read_bytes()).hexdigest(), EXPECTED_LICENSE_SHA256)
        compile(SAMPLE.read_bytes(), str(SAMPLE), "exec")
        lock = json.loads(LOCK.read_text(encoding="utf-8"))
        self.assertEqual(lock["source"]["commit"], "8555d14532a9688b751d8408d822d1dd5feb47f6")
        self.assertTrue(lock["source"]["commit_verified"])
        self.assertEqual(lock["fixtures"][0]["sha256"], EXPECTED_FIXTURE_SHA256)
        self.assertEqual(lock["fixtures"][0]["destination_path"], "sample_forms/forms/sample_invoice.jpg")

    def test_official_sample_executes_against_exact_fixture_contract(self) -> None:
        with tempfile.TemporaryDirectory() as temp:
            sample_root = Path(temp)
            copied = sample_root / SAMPLE.name
            shutil.copyfile(SAMPLE, copied)
            fixture = sample_root / "sample_forms" / "forms" / "sample_invoice.jpg"
            fixture.parent.mkdir(parents=True)
            fixture_bytes = b"fixture-contract-probe"
            fixture.write_bytes(fixture_bytes)
            observed: dict[str, object] = {}

            class AzureKeyCredential:
                def __init__(self, key: str):
                    observed["key"] = key

            class AnalyzeDocumentRequest:
                def __init__(self, *, bytes_source: bytes):
                    observed["bytes_source"] = bytes_source

            class Poller:
                def result(self):
                    return types.SimpleNamespace(documents=[])

            class DocumentIntelligenceClient:
                def __init__(self, *, endpoint: str, credential: object):
                    observed["endpoint"] = endpoint
                    observed["credential"] = credential

                def begin_analyze_document(self, model_id: str, request: object, *, locale: str):
                    observed["model_id"] = model_id
                    observed["request"] = request
                    observed["locale"] = locale
                    return Poller()

            modules = {
                "azure": types.ModuleType("azure"),
                "azure.core": types.ModuleType("azure.core"),
                "azure.core.credentials": types.ModuleType("azure.core.credentials"),
                "azure.ai": types.ModuleType("azure.ai"),
                "azure.ai.documentintelligence": types.ModuleType("azure.ai.documentintelligence"),
                "azure.ai.documentintelligence.models": types.ModuleType("azure.ai.documentintelligence.models"),
            }
            modules["azure.core.credentials"].AzureKeyCredential = AzureKeyCredential
            modules["azure.ai.documentintelligence"].DocumentIntelligenceClient = DocumentIntelligenceClient
            modules["azure.ai.documentintelligence.models"].AnalyzeDocumentRequest = AnalyzeDocumentRequest
            modules["azure.ai.documentintelligence.models"].AnalyzeResult = object
            previous = {name: sys.modules.get(name) for name in modules}
            sys.modules.update(modules)
            old_endpoint = os.environ.get("DOCUMENTINTELLIGENCE_ENDPOINT")
            old_key = os.environ.get("DOCUMENTINTELLIGENCE_API_KEY")
            os.environ["DOCUMENTINTELLIGENCE_ENDPOINT"] = "https://example.invalid"
            os.environ["DOCUMENTINTELLIGENCE_API_KEY"] = "test-only"
            try:
                spec = importlib.util.spec_from_file_location("official_invoice_sample", copied)
                assert spec and spec.loader
                module = importlib.util.module_from_spec(spec)
                spec.loader.exec_module(module)
                module.analyze_invoice()
            finally:
                if old_endpoint is None:
                    os.environ.pop("DOCUMENTINTELLIGENCE_ENDPOINT", None)
                else:
                    os.environ["DOCUMENTINTELLIGENCE_ENDPOINT"] = old_endpoint
                if old_key is None:
                    os.environ.pop("DOCUMENTINTELLIGENCE_API_KEY", None)
                else:
                    os.environ["DOCUMENTINTELLIGENCE_API_KEY"] = old_key
                for name, prior in previous.items():
                    if prior is None:
                        sys.modules.pop(name, None)
                    else:
                        sys.modules[name] = prior

            self.assertEqual(observed["bytes_source"], fixture_bytes)
            self.assertEqual(observed["model_id"], "prebuilt-invoice")
            self.assertEqual(observed["locale"], "en-US")
            self.assertEqual(observed["endpoint"], "https://example.invalid")


if __name__ == "__main__":
    unittest.main()
````

### FILE: `azure_document_intelligence_official_invoice/test_fixture_acquisition.ps1`
```yaml
block_id: "AZURE-DI-OFFICIAL-INVOICE:acquisition-regression:v1"
operation: CREATE
provenance: AUTHORED
source: "local acquisition regression"
license: "LicenseRef-Workspace-Owner"
sha256: "18461d7e24537d5cd52cff8ede867823b77853eef1f658a0c3f1e080a89ad0d6"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
$ErrorActionPreference = 'Stop'
$runner = Join-Path $PSScriptRoot 'acquire_official_fixture.ps1'
$realLock = Join-Path $PSScriptRoot 'fixture-lock.json'

function Expect-Failure([scriptblock]$Action, [string]$Label) {
  $failed = $false
  try { & $Action } catch { $failed = $true }
  if (-not $failed) { throw "Expected failure: $Label" }
}

& $runner -LockPath $realLock -FixtureId 'azure-documentintelligence-1.0.2-sample-invoice-jpg' -ValidateOnly
Expect-Failure { & $runner -LockPath $realLock -FixtureId 'missing' -ValidateOnly } 'unknown fixture'

$temp = Join-Path ([IO.Path]::GetTempPath()) ('elite-azure-di-fixture-test-' + [guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Path $temp | Out-Null
try {
  $cache = Join-Path $temp 'cache'
  New-Item -ItemType Directory -Path $cache | Out-Null
  $fixtureName = 'fixture.bin'
  $fixturePath = Join-Path $cache $fixtureName
  [IO.File]::WriteAllBytes($fixturePath, [Text.Encoding]::UTF8.GetBytes('fixture-regression'))
  $bytes = (Get-Item -LiteralPath $fixturePath).Length
  $sha = (Get-FileHash -LiteralPath $fixturePath -Algorithm SHA256).Hash.ToLowerInvariant()
  $lockObject = [ordered]@{
    schema='elite-azure-documentintelligence-official-fixture-lock/v1'
    verified_at='2026-08-28'
    source=[ordered]@{owner='Microsoft';repository='Azure/azure-sdk-for-python';release='azure-ai-documentintelligence_1.0.2';commit='8555d14532a9688b751d8408d822d1dd5feb47f6';commit_verified=$true;license_expression='MIT';license_sha256='7c77a44a8acd9b41fdc209864a8016b3d430b5d0e09309818d5b7444336df744'}
    fixtures=@([ordered]@{id='fixture-test';url="https://raw.githubusercontent.com/Azure/azure-sdk-for-python/8555d14532a9688b751d8408d822d1dd5feb47f6/test/$fixtureName";filename=$fixtureName;destination_path='sample_forms/forms/fixture.bin';bytes=$bytes;sha256=$sha})
  }
  $fixtureLock = Join-Path $temp 'lock.json'
  [IO.File]::WriteAllText($fixtureLock, ($lockObject | ConvertTo-Json -Depth 10) + [Environment]::NewLine, [Text.UTF8Encoding]::new($false))
  $lockHash = (Get-FileHash -LiteralPath $fixtureLock -Algorithm SHA256).Hash.ToLowerInvariant()
  $approvalObject = [ordered]@{schema='elite-azure-documentintelligence-official-fixture-approval/v1';lock_sha256=$lockHash;approved_fixture_ids=@('fixture-test');approved_by='regression';approved_at='2026-08-28';accept_microsoft_mit_license=$true}
  $approval = Join-Path $temp 'approval.json'
  [IO.File]::WriteAllText($approval, ($approvalObject | ConvertTo-Json -Depth 10) + [Environment]::NewLine, [Text.UTF8Encoding]::new($false))

  $destination = Join-Path $temp 'output'
  & $runner -LockPath $fixtureLock -FixtureId 'fixture-test' -DestinationRoot $destination -ApprovalPath $approval -ImportDirectory $cache
  $installed = Join-Path $destination 'sample_forms/forms/fixture.bin'
  if ((Get-FileHash -LiteralPath $installed -Algorithm SHA256).Hash.ToLowerInvariant() -cne $sha) { throw 'installed fixture hash mismatch' }
  if (-not (Test-Path -LiteralPath (Join-Path $destination 'FIXTURE_RECEIPT.json'))) { throw 'fixture receipt missing' }
  Expect-Failure { & $runner -LockPath $fixtureLock -FixtureId 'fixture-test' -DestinationRoot $destination -ApprovalPath $approval -ImportDirectory $cache } 'occupied destination'

  $tamperedCache = Join-Path $temp 'tampered-cache'
  New-Item -ItemType Directory -Path $tamperedCache | Out-Null
  [IO.File]::WriteAllText((Join-Path $tamperedCache $fixtureName), 'tampered', [Text.UTF8Encoding]::new($false))
  $tamperedDestination = Join-Path $temp 'tampered-output'
  Expect-Failure { & $runner -LockPath $fixtureLock -FixtureId 'fixture-test' -DestinationRoot $tamperedDestination -ApprovalPath $approval -ImportDirectory $tamperedCache } 'tampered cache'
  if (Test-Path -LiteralPath $tamperedDestination) { throw 'tampered destination survived' }
  Write-Output 'AZURE_DI_OFFICIAL_FIXTURE_TEST_PASS positives=2 negatives=3'
} finally {
  if (Test-Path -LiteralPath $temp) { Remove-Item -LiteralPath $temp -Recurse -Force }
}
````

## 6. Configuration surface

- `AZURE_DOCUMENT_INTELLIGENCE_ENDPOINT` y `AZURE_DOCUMENT_INTELLIGENCE_KEY`: requeridos sólo por la ejecución live del sample oficial; son secretos/aportes del proyecto y nunca se guardan en lock, approval ni Markdown.
- `fixture-approval.json`: copia completada de la plantilla con identidad, fecha, aceptación MIT y SHA-256 del lock exacto.
- `-FixtureId`, `-CacheRoot`, `-AllowNetwork`, `-ValidateOnly`: selección y frontera de adquisición; red permanece deshabilitada por defecto.
- El modelo y locale permanecen exactamente como upstream (`prebuilt-invoice`, `en-US`); cualquier cambio exige un adapter/version/evidencia nueva.

## 7. Dependency bill

| Dependencia | Versión/identidad | Licencia | Uso | Condición |
|---|---|---|---|---|
| Microsoft Azure SDK sample | commit `8555d14532a9688b751d8408d822d1dd5feb47f6`, tag `azure-ai-documentintelligence_1.0.2` | MIT | código de integración invoice desde bytes | bytes y SHA-256 exactos |
| azure-ai-documentintelligence | `1.0.2` | MIT | runtime live del sample | adquirir por artifact core y lock transitivo del target |
| sample_invoice.jpg | mismo commit Microsoft, 184.686 bytes | MIT | fixture oficial | aprobación, cache/network explícita, SHA-256 |
| PowerShell | `7+` | MIT | adquisición/regresión | local o CI autorizado |
| Python | `3.8+` | PSF-2.0 | sample/regresión | entorno aislado del proyecto |

## 8. Apply order

1. Complete el blueprint y seleccione Azure Document Intelligence; no use este pack por defecto si el proveedor no está decidido.
2. Materialice/adquiera `OFFICIAL-DOCUMENT-SDK-ARTIFACT-CORE` 0.2.x y construya un entorno Python con lock transitivo aprobado.
3. Materialice este pack y ejecute ambos tests offline.
4. Copie la plantilla de aprobación, fije `lock_sha256`, complete identidad/fecha y acepte MIT.
5. Valide el lock; adquiera el fixture desde cache o con red explícita y preserve receipt.
6. Sólo con cuenta/secretos aportados por el proyecto ejecute el sample live en un entorno no productivo.
7. Integre su salida mediante los contratos de evidencia/evaluación/revisión; no persista automáticamente desde el sample.

## 9. Verification

```powershell
python .\azure_document_intelligence_official_invoice\test_official_invoice_sample.py
pwsh -NoProfile -File .\azure_document_intelligence_official_invoice\test_fixture_acquisition.ps1
pwsh -NoProfile -File .\azure_document_intelligence_official_invoice\acquire_official_fixture.ps1 -LockPath .\azure_document_intelligence_official_invoice\fixture-lock.json -ApprovalPath .\azure_document_intelligence_official_invoice\fixture-approval.json -FixtureId azure-di-1.0.2-sample-invoice-jpg -DestinationRoot <new-directory> -CacheRoot <approved-cache>
```

PASS exige SHA-256 exacto de sample/licencia/fixture lock, ejecución offline de la función upstream, dos positivos, tres negativos, receipt y hash del fixture. Una prueba live requiere además cuenta Azure, endpoint, key, cuota/costo aceptados y evidencia separada; su ausencia no se representa como PASS productivo.

## 10. Reconstruction evidence

- Repositorio oficial Microsoft `Azure/azure-sdk-for-python`, tag `azure-ai-documentintelligence_1.0.2`, commit `8555d14532a9688b751d8408d822d1dd5feb47f6`, firma GitHub verificada y fecha 2025-03-26.
- Sample sync 13.084 bytes SHA-256 `ebc32b0fc534e625b15e636c6b55376f01d28a83f35217dd141202073297a293`; coincide con raw GitHub y el sdist oficial PyPI 1.0.2.
- Licencia MIT 1.074 bytes SHA-256 `7c77a44a8acd9b41fdc209864a8016b3d430b5d0e09309818d5b7444336df744`.
- Fixture oficial `sample_invoice.jpg` 184.686 bytes SHA-256 `489f0c63b6a05e0ae0fcfbf29121299fc2ccacab35d5cd40a5d63d42764078cb`.
- Regresiones offline reconstruidas desde este Markdown; la ejecución live y exactitud de negocio permanecen condiciones del proyecto.
