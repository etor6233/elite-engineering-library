#requires -Version 7.0

[CmdletBinding()]
param(
  [ValidateSet('Audit','Preflight','Foundation','Toolchain')]
  [string] $Mode = 'Audit',
  [string] $PythonExecutable = '',
  [string] $GoExecutable = '',
  [string] $DotNetExecutable = '',
  [string] $NodeExecutable = '',
  [string] $PnpmExecutable = '',
  [string] $DockerExecutable = '',
  [string] $PsqlExecutable = '',
  [string] $PromtoolExecutable = '',
  [string] $PromtoolSha256 = '',
  [switch] $AllowNetwork,
  [switch] $RequireAll,
  [string] $EvidencePath = ''
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$libraryRoot = [IO.Path]::GetFullPath($PSScriptRoot)
$tempRoot = Join-Path ([IO.Path]::GetTempPath()) ('elite-executable-verification-' + [guid]::NewGuid().ToString('N'))
$utf8 = [Text.UTF8Encoding]::new($false)
$steps = [Collections.Generic.List[object]]::new()
$upstreamSourceCount = 0
$implementationPackCount = 0
$documentArtifactCount = 0
$officialInvoiceSampleCount = 0
$officialAzureCuSampleCount = 0
$officialAzureCuBinarySampleCount = 0
$officialAzureCuCopySampleCount = 0
$officialGoogleProcessSampleCount = 0
$officialGoogleLifecycleCount = 0
$officialDaprOutboxCount = 0
$businessCentralArtifactCount = 0
$providerAdapterCount = 0
$documentOrchestratorCount = 0
$evidenceLogCount = 0
$awsPowertoolsBatchComponentCount = 0
$awsDurableExecutionComponentCount = 0
$awsTextractorOfficialComponentCount = 0
$officialAvmSftpIntakeCount = 0
$secureEmailMimeCoreCount = 0
$awsSesImmutableEmailReceiverCount = 0
$awsGuardDutyImmutableReleaseGateCount = 0
$awsGuardDutyMagikaIdpDispatchGateCount = 0
$awsIdpImmutableEvaluationHandoffCount = 0
$awsIdpEvaluationDecisionWorkerCount = 0
$awsIdpPostgresPersistenceBoundaryCount = 0
$debeziumPostgresOutboxRuntimeCount = 0
$debeziumPostgresInboxConsumerCount = 0
$googleCelDocumentMappingCount = 0
$microsoftPlaywrightBrowserGateCount = 0
$googleLighthouseWebQualityGateCount = 0
$microsoftKiotaOpenApiClientGateCount = 0
$executionStateControlCount = 0
$capabilityGapResolutionGateCount = 0
$microsoftDevSkimAdaptedSastGateCount = 0
$gitlabOpengrepSignedSastGateCount = 0
$portableSignedReleaseEvidenceGateCount = 0
$prometheusRuleArtifact = $null

function Fail([string] $Message) { throw "VERIFY_EXECUTABLE_LIBRARY_FAILED: $Message" }

function Resolve-Tool([string] $Name, [string] $Explicit = '') {
  if ($Explicit) {
    $path = [IO.Path]::GetFullPath($Explicit)
    if (Test-Path -LiteralPath $path -PathType Leaf) { return $path }
    return $null
  }
  $command = Get-Command $Name -ErrorAction SilentlyContinue | Select-Object -First 1
  if ($command) { return $command.Source }
  $null
}

function Get-ToolVersion([string] $Path, [string[]] $Arguments) {
  if (-not $Path) { return '' }
  $output = & $Path @Arguments 2>&1
  if ($LASTEXITCODE -ne 0) { return '' }
  ($output | Out-String).Trim()
}

function Resolve-PythonTool {
  $resolved = Resolve-Tool 'python' $PythonExecutable
  # An explicit selection is authoritative even when unavailable; never substitute it.
  if (-not $resolved -and -not $PythonExecutable) { $resolved = Resolve-Tool 'python3' }
  $resolved
}

function Test-MinimumVersion([string] $Text, [version] $Minimum) {
  $match = [regex]::Match($Text, '(?<version>\d+\.\d+(?:\.\d+)?)')
  if (-not $match.Success) { return $false }
  try { return ([version]$match.Groups['version'].Value -ge $Minimum) } catch { return $false }
}

function Invoke-Checked([string] $Id, [string] $WorkingDirectory, [scriptblock] $Action) {
  $started = [DateTimeOffset]::UtcNow
  Push-Location $WorkingDirectory
  try {
    & $Action
    if ($LASTEXITCODE -ne 0) { Fail "$Id exited $LASTEXITCODE" }
    $steps.Add([ordered]@{ id=$Id; status='PASS'; elapsed_ms=[math]::Round(([DateTimeOffset]::UtcNow-$started).TotalMilliseconds) })
  } finally { Pop-Location }
}

function Stop-OwnedProcessTree([Diagnostics.Process] $Process) {
  if ($null -eq $Process) { return }
  if ($IsWindows) {
    $all = @(Get-CimInstance Win32_Process)
    $children = @{}
    foreach ($item in $all) {
      $parent = [int]$item.ParentProcessId
      if (-not $children.ContainsKey($parent)) { $children[$parent] = [Collections.Generic.List[int]]::new() }
      $children[$parent].Add([int]$item.ProcessId)
    }
    $ordered = [Collections.Generic.List[int]]::new()
    function Add-Descendants([int] $ParentId) {
      if (-not $children.ContainsKey($ParentId)) { return }
      foreach ($childId in @($children[$ParentId])) { Add-Descendants $childId; $ordered.Add($childId) }
    }
    Add-Descendants $Process.Id
    foreach ($id in @($ordered)) { Stop-Process -Id $id -Force -ErrorAction SilentlyContinue }
  }
  if (-not $Process.HasExited) { Stop-Process -Id $Process.Id -Force -ErrorAction SilentlyContinue }
  try { $Process.WaitForExit(5000) | Out-Null } catch {}
}

function Test-HashLockedPythonRequirements([string] $Path, [int] $ExpectedCount) {
  if (-not (Test-Path -LiteralPath $Path -PathType Leaf)) { Fail "requirements lock missing: $Path" }
  $entries = @(
    Get-Content -LiteralPath $Path |
      ForEach-Object { $_.Trim() } |
      Where-Object { $_ -and -not $_.StartsWith('#') }
  )
  if ($entries.Count -ne $ExpectedCount) { Fail "requirements lock entry count mismatch: $Path expected=$ExpectedCount actual=$($entries.Count)" }
  $names = [Collections.Generic.HashSet[string]]::new([StringComparer]::OrdinalIgnoreCase)
  foreach ($entry in $entries) {
    $match = [regex]::Match($entry, '^(?<name>[A-Za-z0-9_.-]+) @ (?<url>https://files\.pythonhosted\.org/[^\s#]+)#sha256=(?<sha>[0-9a-f]{64})$')
    if (-not $match.Success) { Fail "requirements entry is not an exact pythonhosted wheel with SHA-256: $entry" }
    if (-not $names.Add($match.Groups['name'].Value)) { Fail "duplicate distribution in requirements lock: $($match.Groups['name'].Value)" }
  }
}

function Test-PinnedHashPythonRequirements([string] $Path, [int] $ExpectedCount) {
  if (-not (Test-Path -LiteralPath $Path -PathType Leaf)) { Fail "requirements lock missing: $Path" }
  $entries = @(Get-Content -LiteralPath $Path | ForEach-Object { $_.Trim() } | Where-Object { $_ -and -not $_.StartsWith('#') })
  if ($entries.Count -ne $ExpectedCount) { Fail "requirements lock entry count mismatch: $Path expected=$ExpectedCount actual=$($entries.Count)" }
  $names = [Collections.Generic.HashSet[string]]::new([StringComparer]::OrdinalIgnoreCase)
  foreach ($entry in $entries) {
    $match = [regex]::Match($entry, '^(?<name>[A-Za-z0-9_.-]+)==(?<version>[^ ]+) --hash=sha256:(?<sha>[0-9a-f]{64})$')
    if (-not $match.Success) { Fail "requirements entry is not an exact version plus SHA-256: $entry" }
    if (-not $names.Add($match.Groups['name'].Value)) { Fail "duplicate distribution in requirements lock: $($match.Groups['name'].Value)" }
  }
}

function Test-MultilineHashPythonRequirements([string] $Path, [int] $ExpectedCount) {
  if (-not (Test-Path -LiteralPath $Path -PathType Leaf)) { Fail "requirements lock missing: $Path" }
  $lines = @(Get-Content -LiteralPath $Path | Where-Object { $_.Trim() -and -not $_.Trim().StartsWith('#') })
  if ($lines.Count -ne ($ExpectedCount * 2)) { Fail "multiline requirements lock line count mismatch: $Path expected=$($ExpectedCount * 2) actual=$($lines.Count)" }
  $names = [Collections.Generic.HashSet[string]]::new([StringComparer]::OrdinalIgnoreCase)
  for ($index = 0; $index -lt $lines.Count; $index += 2) {
    $package = [regex]::Match($lines[$index], '^(?<name>[A-Za-z0-9_.-]+)==(?<version>[^\s\\]+) \\$')
    $hash = [regex]::Match($lines[$index + 1], '^\s+--hash=sha256:(?<sha>[0-9a-f]{64})$')
    if (-not $package.Success -or -not $hash.Success) { Fail "requirements entry is not an exact version plus SHA-256 pair: $($lines[$index]) / $($lines[$index + 1])" }
    if (-not $names.Add($package.Groups['name'].Value)) { Fail "duplicate distribution in requirements lock: $($package.Groups['name'].Value)" }
  }
}

function Materialize-ProviderAdapter([string] $PackName, [string] $Destination) {
  & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot "implementation_packs/$PackName") -Destination $Destination | Out-Null
  if ($LASTEXITCODE -ne 0) { Fail "provider adapter materialization failed: $PackName exit=$LASTEXITCODE" }
}

function Invoke-ProviderAdapterAudit {
  $amazonRoot = Join-Path $tempRoot 'amazon-spapi-adapter'
  $amazonEasyShipRoot = Join-Path $tempRoot 'amazon-easyship-adapter'
  $amazonFulfillmentDeliveryEvidenceRoot = Join-Path $tempRoot 'amazon-fulfillment-delivery-evidence-adapter'
  $amazonSupplySourcesRoot = Join-Path $tempRoot 'amazon-supply-sources-adapter'
  $amazonMliInventoryRoot = Join-Path $tempRoot 'amazon-mli-inventory-adapter'
  $amazonExternalInventoryRoot = Join-Path $tempRoot 'amazon-external-inventory-adapter'
  $googleAdsRoot = Join-Path $tempRoot 'google-ads-adapter'
  $merchantRoot = Join-Path $tempRoot 'google-merchant-adapter'
  $metaRoot = Join-Path $tempRoot 'meta-ads-adapter'
  $metaLeadRoot = Join-Path $tempRoot 'meta-lead-reconciliation-adapter'
  $tiktokRoot = Join-Path $tempRoot 'tiktok-ads-adapter'
  $tiktokLeadRoot = Join-Path $tempRoot 'tiktok-lead-adapter'
  $whatsAppRoot = Join-Path $tempRoot 'meta-whatsapp-adapter'
  $firebaseRoot = Join-Path $tempRoot 'firebase-push-adapter'
  $mercadoLibreRoot = Join-Path $tempRoot 'mercadolibre-marketplace-adapter'
  $azureBlobRoot = Join-Path $tempRoot 'azure-blob-evidence-adapter'
  Invoke-Checked 'amazon-spapi-adapter-materialization' $libraryRoot {
    Materialize-ProviderAdapter 'PYTHON_AMAZON_SPAPI_CATALOG_ADAPTER.md' $amazonRoot
  }
  Invoke-Checked 'amazon-spapi-hash-lock-contract' $amazonRoot {
    Test-HashLockedPythonRequirements (Join-Path $amazonRoot 'amazon_spapi_catalog/requirements.lock') 7
    $sdk = Get-Content -LiteralPath (Join-Path $amazonRoot 'amazon_spapi_catalog/sdk-artifact.lock.json') -Raw | ConvertFrom-Json -Depth 20
    if ($sdk.provider -ne 'Amazon' -or $sdk.distribution -ne 'amzn-sp-api' -or $sdk.version -ne '1.11.1') { Fail 'Amazon SDK artifact identity mismatch' }
  }
  Invoke-Checked 'amazon-easyship-adapter-materialization' $libraryRoot {
    Materialize-ProviderAdapter 'PYTHON_AMAZON_SPAPI_EASYSHIP_HANDOVER_ADAPTER.md' $amazonEasyShipRoot
  }
  Invoke-Checked 'amazon-easyship-hash-lock-contract' $amazonEasyShipRoot {
    Test-HashLockedPythonRequirements (Join-Path $amazonEasyShipRoot 'amazon_spapi_easyship_handover/requirements.lock') 7
    $sdk = Get-Content -LiteralPath (Join-Path $amazonEasyShipRoot 'amazon_spapi_easyship_handover/sdk-artifact.lock.json') -Raw | ConvertFrom-Json -Depth 20
    if ($sdk.provider -ne 'Amazon' -or $sdk.distribution -ne 'amzn-sp-api' -or $sdk.version -ne '1.11.1') { Fail 'Amazon Easy Ship SDK artifact identity mismatch' }
  }
  Invoke-Checked 'amazon-fulfillment-delivery-evidence-adapter-materialization' $libraryRoot {
    Materialize-ProviderAdapter 'PYTHON_AMAZON_SPAPI_FULFILLMENT_DELIVERY_EVIDENCE_ADAPTER.md' $amazonFulfillmentDeliveryEvidenceRoot
  }
  Invoke-Checked 'amazon-fulfillment-delivery-evidence-hash-lock-contract' $amazonFulfillmentDeliveryEvidenceRoot {
    Test-HashLockedPythonRequirements (Join-Path $amazonFulfillmentDeliveryEvidenceRoot 'amazon_spapi_fulfillment_delivery_evidence/requirements.lock') 7
    $sdk = Get-Content -LiteralPath (Join-Path $amazonFulfillmentDeliveryEvidenceRoot 'amazon_spapi_fulfillment_delivery_evidence/sdk-artifact.lock.json') -Raw | ConvertFrom-Json -Depth 20
    if ($sdk.provider -ne 'Amazon' -or $sdk.distribution -ne 'amzn-sp-api' -or $sdk.version -ne '1.11.1') { Fail 'Amazon Fulfillment delivery evidence SDK artifact identity mismatch' }
  }
  Invoke-Checked 'amazon-supply-sources-adapter-materialization' $libraryRoot {
    Materialize-ProviderAdapter 'PYTHON_AMAZON_SPAPI_SUPPLY_SOURCES_ADAPTER.md' $amazonSupplySourcesRoot
  }
  Invoke-Checked 'amazon-supply-sources-hash-lock-contract' $amazonSupplySourcesRoot {
    Test-HashLockedPythonRequirements (Join-Path $amazonSupplySourcesRoot 'amazon_spapi_supply_sources/requirements.lock') 7
    $sdk = Get-Content -LiteralPath (Join-Path $amazonSupplySourcesRoot 'amazon_spapi_supply_sources/sdk-artifact.lock.json') -Raw | ConvertFrom-Json -Depth 20
    if ($sdk.provider -ne 'Amazon' -or $sdk.distribution -ne 'amzn-sp-api' -or $sdk.version -ne '1.11.1') { Fail 'Amazon Supply Sources SDK artifact identity mismatch' }
  }
  Invoke-Checked 'amazon-mli-inventory-adapter-materialization' $libraryRoot {
    Materialize-ProviderAdapter 'PYTHON_AMAZON_SPAPI_MLI_INVENTORY_ADAPTER.md' $amazonMliInventoryRoot
  }
  Invoke-Checked 'amazon-mli-inventory-hash-lock-contract' $amazonMliInventoryRoot {
    Test-HashLockedPythonRequirements (Join-Path $amazonMliInventoryRoot 'amazon_spapi_mli_inventory/requirements.lock') 7
    $sdk = Get-Content -LiteralPath (Join-Path $amazonMliInventoryRoot 'amazon_spapi_mli_inventory/sdk-artifact.lock.json') -Raw | ConvertFrom-Json -Depth 20
    if ($sdk.provider -ne 'Amazon' -or $sdk.distribution -ne 'amzn-sp-api' -or $sdk.version -ne '1.11.1') { Fail 'Amazon MLI Inventory SDK artifact identity mismatch' }
  }
  Invoke-Checked 'amazon-external-inventory-adapter-materialization' $libraryRoot {
    Materialize-ProviderAdapter 'PYTHON_AMAZON_SPAPI_EXTERNAL_FULFILLMENT_INVENTORY_ADAPTER.md' $amazonExternalInventoryRoot
  }
  Invoke-Checked 'amazon-external-inventory-hash-lock-contract' $amazonExternalInventoryRoot {
    Test-HashLockedPythonRequirements (Join-Path $amazonExternalInventoryRoot 'amazon_spapi_external_inventory/requirements.lock') 7
    $sdk = Get-Content -LiteralPath (Join-Path $amazonExternalInventoryRoot 'amazon_spapi_external_inventory/sdk-artifact.lock.json') -Raw | ConvertFrom-Json -Depth 20
    if ($sdk.provider -ne 'Amazon' -or $sdk.distribution -ne 'amzn-sp-api' -or $sdk.version -ne '1.11.1') { Fail 'Amazon External Inventory SDK artifact identity mismatch' }
    if ($sdk.model_source.external_fulfillment_model_license -ne 'LicenseRef-Amazon-Software-License') { Fail 'Amazon External Inventory model license condition missing' }
  }
  Invoke-Checked 'google-merchant-adapter-materialization' $libraryRoot {
    Materialize-ProviderAdapter 'PYTHON_GOOGLE_MERCHANT_PRODUCT_SYNC_ADAPTER.md' $merchantRoot
  }
  Invoke-Checked 'google-merchant-hash-lock-contract' $merchantRoot {
    Test-HashLockedPythonRequirements (Join-Path $merchantRoot 'google_merchant_product_sync/requirements-windows-py314.lock') 20
    $sdk = Get-Content -LiteralPath (Join-Path $merchantRoot 'google_merchant_product_sync/sdk-artifact.lock.json') -Raw | ConvertFrom-Json -Depth 20
    if ($sdk.provider -ne 'Google' -or $sdk.distribution -ne 'google-shopping-merchant-products' -or $sdk.version -ne '1.8.0') { Fail 'Google Merchant SDK artifact identity mismatch' }
  }
  Invoke-Checked 'google-ads-adapter-materialization' $libraryRoot {
    Materialize-ProviderAdapter 'PYTHON_GOOGLE_ADS_REPORTING_ADAPTER.md' $googleAdsRoot
  }
  Invoke-Checked 'google-ads-hash-lock-contract' $googleAdsRoot {
    Test-HashLockedPythonRequirements (Join-Path $googleAdsRoot 'google_ads_reporting/requirements-windows-py314.lock') 24
    $sdk = Get-Content -LiteralPath (Join-Path $googleAdsRoot 'google_ads_reporting/sdk-artifact.lock.json') -Raw | ConvertFrom-Json -Depth 20
    if ($sdk.provider -ne 'Google' -or $sdk.distribution -ne 'google-ads' -or $sdk.version -ne '31.4.0') { Fail 'Google Ads SDK artifact identity mismatch' }
    if ($sdk.wheel.sha256 -ne '210a7f6a7b5d40ef544090216dceeb90c9d2e7fba51c72329a690c9e5bb94474' -or $sdk.source.commit -ne 'b8eb80ae277920d56fef467617795a7360c492fb' -or $sdk.source.commit_signature_verified -ne $true) { Fail 'Google Ads SDK revision/integrity mismatch' }
  }
  Invoke-Checked 'meta-ads-adapter-materialization' $libraryRoot {
    Materialize-ProviderAdapter 'PYTHON_META_ADS_REPORTING_ADAPTER.md' $metaRoot
  }
  Invoke-Checked 'meta-ads-hash-lock-contract' $metaRoot {
    Test-HashLockedPythonRequirements (Join-Path $metaRoot 'meta_ads_reporting/requirements-windows-py314.lock') 18
    $sdk = Get-Content -LiteralPath (Join-Path $metaRoot 'meta_ads_reporting/sdk-artifact.lock.json') -Raw | ConvertFrom-Json -Depth 20
    if ($sdk.provider -ne 'Meta' -or $sdk.distribution -ne 'facebook-business' -or $sdk.version -ne '26.0.1') { Fail 'Meta Business SDK artifact identity mismatch' }
    if ($sdk.license_expression -ne 'LicenseRef-Meta-Platform' -or -not $sdk.distribution_notice) { Fail 'Meta restricted license/redistribution notice missing' }
  }
  Invoke-Checked 'meta-lead-reconciliation-materialization' $libraryRoot {
    Materialize-ProviderAdapter 'PYTHON_META_LEAD_RECONCILIATION_ADAPTER.md' $metaLeadRoot
  }
  Invoke-Checked 'meta-lead-reconciliation-hash-lock-contract' $metaLeadRoot {
    Test-HashLockedPythonRequirements (Join-Path $metaLeadRoot 'meta_lead_reconciliation/requirements-windows-py314.lock') 18
    $sdk = Get-Content -LiteralPath (Join-Path $metaLeadRoot 'meta_lead_reconciliation/sdk-artifact.lock.json') -Raw | ConvertFrom-Json -Depth 20
    if ($sdk.provider -ne 'Meta' -or $sdk.distribution -ne 'facebook-business' -or $sdk.version -ne '26.0.1') { Fail 'Meta Lead SDK artifact identity mismatch' }
    if ($sdk.wheel.sha256 -ne '41555ce87ac105617ed678a489a630108ac162fbdccfbbf65fa20ae1de7cd0ca' -or $sdk.license_expression -ne 'LicenseRef-Meta-Platform' -or -not $sdk.distribution_notice) { Fail 'Meta Lead SDK integrity/license contract mismatch' }
  }
  Invoke-Checked 'tiktok-ads-adapter-materialization' $libraryRoot {
    Materialize-ProviderAdapter 'PYTHON_TIKTOK_ADS_REPORTING_ADAPTER.md' $tiktokRoot
  }
  Invoke-Checked 'tiktok-ads-hash-lock-contract' $tiktokRoot {
    Test-HashLockedPythonRequirements (Join-Path $tiktokRoot 'tiktok_ads_reporting/requirements.lock') 5
    $sdk = Get-Content -LiteralPath (Join-Path $tiktokRoot 'tiktok_ads_reporting/sdk-artifact.lock.json') -Raw | ConvertFrom-Json -Depth 20
    if ($sdk.provider -ne 'TikTok Pte. Ltd.' -or $sdk.distribution -ne 'tiktok-business-api-sdk-official' -or $sdk.version -ne '1.1.3') { Fail 'TikTok official SDK artifact identity mismatch' }
    if ($sdk.wheel.sha256 -ne '663b4a1f2585c4b386144d968646f695befd733c464418bbaeeac41befff33b7' -or $sdk.wheel.bytes -ne 678182 -or $sdk.license_expression -ne 'MIT' -or $sdk.publisher_evidence.provider_readme_recommends_distribution -ne $true) { Fail 'TikTok wheel integrity/publisher/license contract mismatch' }
  }
  Invoke-Checked 'tiktok-lead-adapter-materialization' $libraryRoot {
    Materialize-ProviderAdapter 'PYTHON_TIKTOK_LEAD_ADAPTER.md' $tiktokLeadRoot
  }
  Invoke-Checked 'tiktok-lead-official-contract-and-fail-closed-profile' $tiktokLeadRoot {
    Test-HashLockedPythonRequirements (Join-Path $tiktokLeadRoot 'tiktok_lead_adapter/requirements.lock') 5
    $sdk = Get-Content -LiteralPath (Join-Path $tiktokLeadRoot 'tiktok_lead_adapter/sdk-artifact.lock.json') -Raw | ConvertFrom-Json -Depth 20
    if ($sdk.provider -ne 'TikTok Pte. Ltd.' -or $sdk.distribution -ne 'tiktok-business-api-sdk-official' -or $sdk.version -ne '1.1.3') { Fail 'TikTok Lead official SDK artifact identity mismatch' }
    if ($sdk.wheel.sha256 -ne '663b4a1f2585c4b386144d968646f695befd733c464418bbaeeac41befff33b7' -or $sdk.wheel.bytes -ne 678182 -or $sdk.license_expression -ne 'MIT') { Fail 'TikTok Lead wheel integrity/license contract mismatch' }
    if ($sdk.publisher_evidence.repository_commit -ne 'f809c396520df2d7b201a9ccc5378d822b728ed3' -or $sdk.publisher_evidence.lead_client_generated_by_provider -ne $false) { Fail 'TikTok Lead wrapper provenance boundary mismatch' }
    $contract = Get-Content -LiteralPath (Join-Path $tiktokLeadRoot 'tiktok_lead_adapter/official-contract.lock.json') -Raw | ConvertFrom-Json -Depth 20
    if ($contract.lead_get.path -ne '/open_api/v1.3/lead/get/' -or $contract.lead_get.method -ne 'GET' -or $contract.lead_get.header -ne 'Access-Token') { Fail 'TikTok Lead GET contract mismatch' }
    if ($contract.subscription_create.path -ne '/open_api/v1.3/subscription/subscribe/' -or $contract.subscription_create.method -ne 'POST' -or $contract.subscription_create.subscribe_entity -ne 'LEAD') { Fail 'TikTok Lead subscription contract mismatch' }
    if ($contract.webhook.delivery_semantics -ne 'AT_LEAST_ONCE' -or $contract.webhook.public_signature_mechanism_demonstrated -ne $false -or $contract.webhook.automatic_persistence_allowed -ne $false) { Fail 'TikTok Lead webhook safety boundary mismatch' }
    $profile = Get-Content -LiteralPath (Join-Path $tiktokLeadRoot 'tiktok_lead_adapter/provider-profile.template.json') -Raw | ConvertFrom-Json -Depth 20
    if ($profile.decision -ne 'BLOCKED_ACCESS_AND_LIVE_PROOF_REQUIRED' -or $profile.automatic_business_write -ne $false) { Fail 'TikTok Lead profile is not fail-closed' }
    if (@($profile.PSObject.Properties | Where-Object { $_.Name -match '(_proven|_approved)$' -and $_.Value -ne $false }).Count -ne 0) { Fail 'TikTok Lead profile contains pre-approved proof' }
  }
  Invoke-Checked 'meta-whatsapp-adapter-materialization' $libraryRoot {
    Materialize-ProviderAdapter 'PYTHON_META_WHATSAPP_CLOUD_ADAPTER.md' $whatsAppRoot
  }
  Invoke-Checked 'meta-whatsapp-source-and-provenance-contract' $whatsAppRoot {
    $source = Get-Content -LiteralPath (Join-Path $whatsAppRoot 'whatsapp_cloud/official-source.lock.json') -Raw | ConvertFrom-Json -Depth 20
    if ($source.source_id -ne 'meta-whatsapp-api-examples' -or $source.commit -ne 'de70ee908a67026e642aaee3703d20464e2a9466' -or $source.commit_signature_verified -ne $true) { Fail 'Meta WhatsApp official source identity mismatch' }
    if ($source.license_expression -ne 'LicenseRef-Meta-Platform-API-Only' -or @($source.files).Count -ne 4) { Fail 'Meta WhatsApp license/source file contract mismatch' }
    if (@($source.files | Where-Object provenance -eq 'VERBATIM_REFERENCE_ONLY').Count -ne 3) { Fail 'Meta WhatsApp verbatim code provenance count mismatch' }
  }
  $auditPython = Resolve-PythonTool
  if (-not $auditPython) { Fail 'Python is required for the dependency-free Meta WhatsApp adapter audit' }
  Invoke-Checked 'meta-whatsapp-offline-unit-contract' (Join-Path $whatsAppRoot 'whatsapp_cloud') {
    & $auditPython -m unittest -v test_whatsapp_cloud.py
  }
  Invoke-Checked 'firebase-push-adapter-materialization' $libraryRoot {
    Materialize-ProviderAdapter 'GO_FIREBASE_CLOUD_MESSAGING_ADAPTER.md' $firebaseRoot
  }
  Invoke-Checked 'firebase-push-source-and-module-contract' $firebaseRoot {
    $source = Get-Content -LiteralPath (Join-Path $firebaseRoot 'firebase_push/sdk-source.lock.json') -Raw | ConvertFrom-Json -Depth 20
    if ($source.provider -ne 'Google Firebase' -or $source.module -ne 'firebase.google.com/go/v4' -or $source.version -ne 'v4.21.0') { Fail 'Firebase Admin Go identity mismatch' }
    if ($source.commit -ne 'eebb06f2a643fbb59b1cb262874a943584475128' -or $source.license_expression -ne 'Apache-2.0') { Fail 'Firebase source commit/license mismatch' }
    if ($source.archive_sha256 -ne '8d7bbc5621ddf19ddf0ebd74df867690cd3ee1bed7499eabf973eba668284b5e') { Fail 'Firebase source archive hash mismatch' }
    $goMod = Get-Content -LiteralPath (Join-Path $firebaseRoot 'firebase_push/go.mod') -Raw
    if ($goMod -notmatch '(?m)^require firebase\.google\.com/go/v4 v4\.21\.0$') { Fail 'Firebase exact module pin missing' }
  }
  Invoke-Checked 'mercadolibre-marketplace-adapter-materialization' $libraryRoot {
    Materialize-ProviderAdapter 'GO_MERCADOLIBRE_MARKETPLACE_ADAPTER.md' $mercadoLibreRoot
  }
  Invoke-Checked 'mercadolibre-current-http-authority-contract' $mercadoLibreRoot {
    $authority = Get-Content -LiteralPath (Join-Path $mercadoLibreRoot 'mercadolibre_marketplace/authority.lock.json') -Raw | ConvertFrom-Json -Depth 20
    if ($authority.provider -ne 'Mercado Libre' -or $authority.api_base_url -ne 'https://api.mercadolibre.com') { Fail 'Mercado Libre HTTP authority mismatch' }
    if ($authority.official_sdk_status -ne 'REJECTED_ARCHIVED_NOT_FUNCTIONAL' -or $authority.archived_go_sdk.archived -ne $true) { Fail 'Mercado Libre archived SDK rejection missing' }
    if (@($authority.official_contracts).Count -ne 6) { Fail 'Mercado Libre official contract inventory mismatch' }
  }
  Invoke-Checked 'azure-blob-evidence-adapter-materialization' $libraryRoot {
    Materialize-ProviderAdapter 'PYTHON_AZURE_BLOB_IMMUTABLE_EVIDENCE_ADAPTER.md' $azureBlobRoot
  }
  Invoke-Checked 'azure-blob-evidence-lock-and-authority-contract' $azureBlobRoot {
    $directory = Join-Path $azureBlobRoot 'azure_blob_storage_adapter'
    Test-PinnedHashPythonRequirements (Join-Path $directory 'requirements-windows-py312.lock') 12
    $lock = Get-Content -LiteralPath (Join-Path $directory 'official-artifact-lock.json') -Raw | ConvertFrom-Json -Depth 20
    if ($lock.package -ne 'azure-storage-blob' -or $lock.version -ne '12.30.0' -or $lock.source.commit -ne '245d2024e6764bb4dd5186b5cd30a776c46993cb') { Fail 'Azure Blob SDK artifact identity mismatch' }
    if ($lock.pypi.wheel.sha256 -ne 'd415ac50b67a8da6b3ae7e9f1014b1b55cd7aafa0b8d4ca9b380568dc7360423' -or $lock.source.archive_sha256 -ne '7c4743ccfe1cb81687ddc949df8343901c0b9a9e3c4e782ffaf4ad5f68c4a8fd') { Fail 'Azure Blob SDK artifact hash mismatch' }
    if ($lock.license -ne 'MIT' -or $lock.source.signature_verified -ne $false -or $lock.source.signature_reason -ne 'unsigned') { Fail 'Azure Blob SDK license/signature condition mismatch' }
    $profile = Get-Content -LiteralPath (Join-Path $directory 'provider-profile.template.json') -Raw | ConvertFrom-Json -Depth 20
    if ($profile.schema_version -ne 'azure-blob-evidence-profile/v1' -or $profile.state -ne 'BLOCKED_CONFIGURATION' -or $profile.policy_mode -ne 'Locked' -or $profile.automatic_storage_authorized -ne $false) { Fail 'Azure Blob profile enabled storage without authority' }
    if (@($profile.approvals.PSObject.Properties | Where-Object Value -ne $false).Count -ne 0) { Fail 'Azure Blob profile contains a pre-approved authority' }
    $backend = Get-Content -LiteralPath (Join-Path $directory 'azureblob/azure_backend.py') -Raw
    foreach ($required in @('overwrite=False','validate_content=True','BlobImmutabilityPolicyMode.LOCKED','version_id=version_id','decompress=False')) {
      if (-not $backend.Contains($required)) { Fail "Azure Blob official provider contract missing: $required" }
    }
    $storage = Get-Content -LiteralPath (Join-Path $directory 'azureblob/storage.py') -Raw
    foreach ($required in @('immutable_storage_with_versioning_enabled','prevent_encryption_scope_override','UNKNOWN_EFFECT','PARTIAL_VERIFICATION','exact-version download SHA-256 mismatch')) {
      if (-not $storage.Contains($required)) { Fail "Azure Blob evidence contract missing: $required" }
    }
  }
  $script:providerAdapterCount = 16
}

function Invoke-OfflineProviderAdapterRuntimeGates([string] $Python) {
  $whatsAppDirectory = Join-Path $tempRoot 'meta-whatsapp-adapter/whatsapp_cloud'
  Invoke-Checked 'meta-whatsapp-official-source-and-adapter-tests' $whatsAppDirectory {
    & $Python -m unittest -v test_whatsapp_cloud.py
  }
}

function Invoke-GoProviderAdapterRuntimeGates([string] $Go) {
  $awsTextractDirectory = Join-Path $tempRoot 'aws-textract-document-runtime/aws_textract_runtime'
  $previousProxy = $env:GOPROXY
  try {
    if ($AllowNetwork) {
      Invoke-Checked 'aws-textract-go-module-download' $awsTextractDirectory { & $Go mod download }
    } else {
      $env:GOPROXY = 'off'
    }
    Invoke-Checked 'aws-textract-go-module-verify' $awsTextractDirectory { & $Go mod verify }
    Invoke-Checked 'aws-textract-profile-security-and-sdk-contract-tests' $awsTextractDirectory { & $Go test ./... }
    Invoke-Checked 'aws-textract-vet' $awsTextractDirectory { & $Go vet ./... }
    Invoke-Checked 'aws-textract-build' $awsTextractDirectory { & $Go build ./cmd/analyze }
    Invoke-Checked 'aws-textract-version-contract' $awsTextractDirectory {
      $version = & $Go list -m -f '{{if eq .Path "github.com/aws/aws-sdk-go-v2/service/textract"}}{{.Version}}{{end}}' all
      if ($LASTEXITCODE -ne 0 -or ($version | Out-String).Trim() -ne 'v1.45.0') { Fail 'AWS Textract resolved version mismatch' }
    }
  } finally { $env:GOPROXY = $previousProxy }
  $firebaseDirectory = Join-Path $tempRoot 'firebase-push-adapter/firebase_push'
  $previousProxy = $env:GOPROXY
  try {
    if ($AllowNetwork) {
      Invoke-Checked 'firebase-admin-go-module-download' $firebaseDirectory { & $Go mod download }
    } else {
      $env:GOPROXY = 'off'
    }
    Invoke-Checked 'firebase-admin-go-module-verify' $firebaseDirectory { & $Go mod verify }
    Invoke-Checked 'firebase-push-unit-and-official-sdk-contract' $firebaseDirectory { & $Go test ./... }
    Invoke-Checked 'firebase-push-vet' $firebaseDirectory { & $Go vet ./... }
    Invoke-Checked 'firebase-admin-go-version-contract' $firebaseDirectory {
      $version = & $Go list -m -f '{{if eq .Path "firebase.google.com/go/v4"}}{{.Version}}{{end}}' all
      if ($LASTEXITCODE -ne 0 -or ($version | Out-String).Trim() -ne 'v4.21.0') { Fail 'Firebase Admin Go resolved version mismatch' }
    }
  } finally { $env:GOPROXY = $previousProxy }
  $mercadoLibreDirectory = Join-Path $tempRoot 'mercadolibre-marketplace-adapter/mercadolibre_marketplace'
  Invoke-Checked 'mercadolibre-marketplace-module-verify' $mercadoLibreDirectory { & $Go mod verify }
  Invoke-Checked 'mercadolibre-marketplace-official-http-contract-tests' $mercadoLibreDirectory { & $Go test ./... }
  Invoke-Checked 'mercadolibre-marketplace-vet' $mercadoLibreDirectory { & $Go vet ./... }
}

function Get-VenvPython([string] $VenvRoot) {
  $candidate = if ($IsWindows) { Join-Path $VenvRoot 'Scripts/python.exe' } else { Join-Path $VenvRoot 'bin/python' }
  if (-not (Test-Path -LiteralPath $candidate -PathType Leaf)) { Fail "virtual environment Python missing: $candidate" }
  $candidate
}

function Invoke-OsvPythonVenvBatchAudit([string[]] $PythonExecutables) {
  # Official contract: https://google.github.io/osv.dev/post-v1-querybatch/
  $packages = @{}
  foreach ($python in $PythonExecutables) {
    $items = (& $python -m pip list --format=json | ConvertFrom-Json)
    if ($LASTEXITCODE -ne 0) { Fail "pip inventory failed for OSV audit: $python" }
    foreach ($item in @($items)) {
      $name = [string]$item.name
      if ($name -and $name -notin @('pip','setuptools','wheel')) { $packages[$name.ToLowerInvariant()] = [string]$item.version }
    }
  }
  $queries = @()
  foreach ($name in ($packages.Keys | Sort-Object)) {
    $queries += @{ package=@{ ecosystem='PyPI'; name=$name }; version=$packages[$name] }
  }
  if ($queries.Count -eq 0) { Fail 'OSV audit received no provider packages' }
  $body = @{ queries=$queries } | ConvertTo-Json -Depth 6 -Compress
  $response = Invoke-RestMethod -Method Post -Uri 'https://api.osv.dev/v1/querybatch' -ContentType 'application/json' -Body $body
  if (@($response.results).Count -ne $queries.Count) { Fail 'OSV querybatch response cardinality mismatch' }
  $findings = [Collections.Generic.List[string]]::new()
  for ($index = 0; $index -lt $queries.Count; $index++) {
    $result = $response.results[$index]
    if ($result.next_page_token) { Fail "OSV response is paginated for $($queries[$index].package.name); complete pagination before admission" }
    foreach ($vulnerability in @($result.vulns)) {
      if ($null -ne $vulnerability) { $findings.Add("$($queries[$index].package.name)@$($queries[$index].version):$($vulnerability.id)") }
    }
  }
  if ($findings.Count -gt 0) { Fail "OSV found provider vulnerabilities: $($findings -join ',')" }
  Write-Output "OSV_PROVIDER_LOCKS_PASS packages=$($queries.Count) vulnerabilities=0"
}

function Invoke-NetworkProviderAdapterGates([string] $Python) {
  if (-not $AllowNetwork) {
    Write-Output 'NETWORK_PROVIDER_ADAPTER_GATES_SKIPPED reason=AllowNetwork_not_authorized'
    return
  }

  $amazonRoot = Join-Path $tempRoot 'amazon-spapi-adapter'
  $amazonEasyShipRoot = Join-Path $tempRoot 'amazon-easyship-adapter'
  $amazonFulfillmentDeliveryEvidenceRoot = Join-Path $tempRoot 'amazon-fulfillment-delivery-evidence-adapter'
  $amazonSupplySourcesRoot = Join-Path $tempRoot 'amazon-supply-sources-adapter'
  $amazonMliInventoryRoot = Join-Path $tempRoot 'amazon-mli-inventory-adapter'
  $amazonExternalInventoryRoot = Join-Path $tempRoot 'amazon-external-inventory-adapter'
  $googleAdsRoot = Join-Path $tempRoot 'google-ads-adapter'
  $metaLeadRoot = Join-Path $tempRoot 'meta-lead-reconciliation-adapter'
  $amazonVenv = Join-Path $tempRoot 'amazon-spapi-venv'
  Invoke-Checked 'amazon-spapi-create-venv' $tempRoot { & $Python -m venv $amazonVenv }
  $amazonPython = Get-VenvPython $amazonVenv
  Invoke-Checked 'amazon-spapi-frozen-install' $amazonRoot {
    & $amazonPython -m pip install --disable-pip-version-check --require-hashes -r (Join-Path $amazonRoot 'amazon_spapi_catalog/requirements.lock')
  }
  Invoke-Checked 'amazon-spapi-pip-check' $amazonRoot { & $amazonPython -m pip check }
  Invoke-Checked 'amazon-spapi-unit-and-sdk-contract' (Join-Path $amazonRoot 'amazon_spapi_catalog') { & $amazonPython -m unittest -v test_catalog_query.py }
  Invoke-Checked 'amazon-easyship-pip-check' $amazonEasyShipRoot { & $amazonPython -m pip check }
  Invoke-Checked 'amazon-easyship-unit-and-sdk-contract' (Join-Path $amazonEasyShipRoot 'amazon_spapi_easyship_handover') { & $amazonPython -m unittest -v test_easyship_handover.py }
  Invoke-Checked 'amazon-fulfillment-delivery-evidence-pip-check' $amazonFulfillmentDeliveryEvidenceRoot { & $amazonPython -m pip check }
  Invoke-Checked 'amazon-fulfillment-delivery-evidence-unit-and-sdk-contract' (Join-Path $amazonFulfillmentDeliveryEvidenceRoot 'amazon_spapi_fulfillment_delivery_evidence') { & $amazonPython -m unittest -v test_fulfillment_delivery_evidence.py }
  Invoke-Checked 'amazon-supply-sources-pip-check' $amazonSupplySourcesRoot { & $amazonPython -m pip check }
  Invoke-Checked 'amazon-supply-sources-unit-and-sdk-contract' (Join-Path $amazonSupplySourcesRoot 'amazon_spapi_supply_sources') { & $amazonPython -m unittest -v test_supply_sources.py }
  Invoke-Checked 'amazon-mli-inventory-pip-check' $amazonMliInventoryRoot { & $amazonPython -m pip check }
  Invoke-Checked 'amazon-mli-inventory-unit-and-sdk-contract' (Join-Path $amazonMliInventoryRoot 'amazon_spapi_mli_inventory') { & $amazonPython -m unittest -v test_mli_inventory.py }
  Invoke-Checked 'amazon-external-inventory-pip-check' $amazonExternalInventoryRoot { & $amazonPython -m pip check }
  Invoke-Checked 'amazon-external-inventory-unit-and-sdk-contract' (Join-Path $amazonExternalInventoryRoot 'amazon_spapi_external_inventory') { & $amazonPython -m unittest -v test_external_inventory.py }

  $pythonDetails = & $Python -c "import json, platform, struct, sys; print(json.dumps({'implementation':platform.python_implementation(),'version':list(sys.version_info[:2]),'bits':struct.calcsize('P')*8,'platform':sys.platform}))"
  if ($LASTEXITCODE -ne 0) { Fail 'could not inspect Python runtime for Google Merchant lock' }
  $runtime = $pythonDetails | ConvertFrom-Json
  if (-not $IsWindows -or $runtime.implementation -ne 'CPython' -or $runtime.version[0] -ne 3 -or $runtime.version[1] -ne 14 -or $runtime.bits -ne 64) {
    Fail 'Google Merchant frozen lock requires CPython 3.14 Windows x86-64; generate and admit a separate exact target lock for another runtime'
  }
  $merchantRoot = Join-Path $tempRoot 'google-merchant-adapter'
  $merchantVenv = Join-Path $tempRoot 'google-merchant-venv'
  Invoke-Checked 'google-merchant-create-venv' $tempRoot { & $Python -m venv $merchantVenv }
  $merchantPython = Get-VenvPython $merchantVenv
  Invoke-Checked 'google-merchant-frozen-install' $merchantRoot {
    & $merchantPython -m pip install --disable-pip-version-check --require-hashes -r (Join-Path $merchantRoot 'google_merchant_product_sync/requirements-windows-py314.lock')
  }
  Invoke-Checked 'google-merchant-pip-check' $merchantRoot { & $merchantPython -m pip check }
  Invoke-Checked 'google-merchant-unit-and-sdk-contract' (Join-Path $merchantRoot 'google_merchant_product_sync') { & $merchantPython -m unittest -v test_sync_product.py }

  $googleAdsVenv = Join-Path $tempRoot 'google-ads-venv'
  Invoke-Checked 'google-ads-create-venv' $tempRoot { & $Python -m venv $googleAdsVenv }
  $googleAdsPython = Get-VenvPython $googleAdsVenv
  Invoke-Checked 'google-ads-frozen-install' $googleAdsRoot {
    & $googleAdsPython -m pip install --disable-pip-version-check --require-hashes --no-deps -r (Join-Path $googleAdsRoot 'google_ads_reporting/requirements-windows-py314.lock')
  }
  Invoke-Checked 'google-ads-pip-check' $googleAdsRoot { & $googleAdsPython -m pip check }
  Invoke-Checked 'google-ads-unit-and-sdk-contract' (Join-Path $googleAdsRoot 'google_ads_reporting') {
    & $googleAdsPython -m unittest -v test_reporting_query.py
    if ($LASTEXITCODE -ne 0) { return }
    & $googleAdsPython -c "import importlib.metadata as m; import google.ads.googleads.v25; assert m.version('google-ads') == '31.4.0'"
  }

  $metaRoot = Join-Path $tempRoot 'meta-ads-adapter'
  $metaVenv = Join-Path $tempRoot 'meta-ads-venv'
  Invoke-Checked 'meta-ads-create-venv' $tempRoot { & $Python -m venv $metaVenv }
  $metaPython = Get-VenvPython $metaVenv
  Invoke-Checked 'meta-ads-frozen-install' $metaRoot {
    & $metaPython -m pip install --disable-pip-version-check --require-hashes -r (Join-Path $metaRoot 'meta_ads_reporting/requirements-windows-py314.lock')
  }
  Invoke-Checked 'meta-ads-pip-check' $metaRoot { & $metaPython -m pip check }
  Invoke-Checked 'meta-ads-unit-and-sdk-contract' (Join-Path $metaRoot 'meta_ads_reporting') { & $metaPython -m unittest -v test_insights.py }

  $metaLeadVenv = Join-Path $tempRoot 'meta-lead-reconciliation-venv'
  Invoke-Checked 'meta-lead-reconciliation-create-venv' $tempRoot { & $Python -m venv $metaLeadVenv }
  $metaLeadPython = Get-VenvPython $metaLeadVenv
  Invoke-Checked 'meta-lead-reconciliation-frozen-install' $metaLeadRoot {
    & $metaLeadPython -m pip install --disable-pip-version-check --require-hashes --no-deps -r (Join-Path $metaLeadRoot 'meta_lead_reconciliation/requirements-windows-py314.lock')
  }
  Invoke-Checked 'meta-lead-reconciliation-pip-check' $metaLeadRoot { & $metaLeadPython -m pip check }
  Invoke-Checked 'meta-lead-reconciliation-unit-and-sdk-contract' (Join-Path $metaLeadRoot 'meta_lead_reconciliation') { & $metaLeadPython -m unittest -v test_fetch_leads.py }

  $tiktokRoot = Join-Path $tempRoot 'tiktok-ads-adapter'
  $tiktokLeadRoot = Join-Path $tempRoot 'tiktok-lead-adapter'
  $tiktokVenv = Join-Path $tempRoot 'tiktok-ads-venv'
  Invoke-Checked 'tiktok-ads-create-venv' $tempRoot { & $Python -m venv $tiktokVenv }
  $tiktokPython = Get-VenvPython $tiktokVenv
  Invoke-Checked 'tiktok-ads-frozen-runtime-install' $tiktokRoot {
    & $tiktokPython -m pip install --disable-pip-version-check --require-hashes --no-deps -r (Join-Path $tiktokRoot 'tiktok_ads_reporting/requirements.lock')
  }
  Invoke-Checked 'tiktok-ads-pip-check' $tiktokRoot { & $tiktokPython -m pip check }
  $adapterDirectory = Join-Path $tiktokRoot 'tiktok_ads_reporting'
  Invoke-Checked 'tiktok-official-wheel-import' $adapterDirectory {
    & $tiktokPython -c "import importlib.metadata as m; from business_api_client.api.reporting_api import ReportingApi; assert m.version('tiktok-business-api-sdk-official') == '1.1.3'; assert callable(ReportingApi.report_integrated_get); print('TIKTOK_OFFICIAL_WHEEL_IMPORT_PASS')"
  }
  Invoke-Checked 'tiktok-ads-unit-and-official-sdk-contract' $adapterDirectory { & $tiktokPython -m unittest -v test_reporting.py }
  Invoke-Checked 'tiktok-lead-unit-and-official-sdk-contract' (Join-Path $tiktokLeadRoot 'tiktok_lead_adapter') {
    & $tiktokPython -m unittest -v test_tiktok_lead.py
  }
  Invoke-Checked 'provider-python-osv-querybatch' $tempRoot {
    Invoke-OsvPythonVenvBatchAudit @($amazonPython,$merchantPython,$googleAdsPython,$metaPython,$metaLeadPython,$tiktokPython)
  }
}

function Invoke-PrometheusRuleRegression {
  # Selection/hash bind the test engine; its provenance/SCA admission is separate.
  if (-not $PromtoolExecutable -and -not $PromtoolSha256) {
    Write-Output 'PROMETHEUS_RULE_REGRESSION_NOT_EXECUTED reason=explicit_admitted_binary_and_sha256_required'
    return
  }
  if (-not $PromtoolExecutable -or $PromtoolSha256 -cnotmatch '^[0-9a-f]{64}$') { Fail 'PromtoolExecutable and lowercase PromtoolSha256 are required together' }
  $promtool = Resolve-Tool 'promtool' $PromtoolExecutable
  if (-not $promtool -or (Get-FileHash -LiteralPath $promtool -Algorithm SHA256).Hash.ToLowerInvariant() -cne $PromtoolSha256) { Fail 'promtool exact artifact mismatch' }
  $ruleRoot = Join-Path $tempRoot 'prometheus-rule-regression'
  Invoke-Checked 'prometheus-rule-pack-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/SECURE_OPERATIONS_DELIVERY_CORE.md') -Destination $ruleRoot
  }
  $rulePython = Resolve-PythonTool
  Invoke-Checked 'prometheus-rules-source-engine-regression' $ruleRoot {
    $code = @'
import hashlib, os, sys
from pathlib import Path
sys.path.insert(0, str(Path('production_admission_gate').resolve()))
from run_official_tool import bounded_process
exe, expected = sys.argv[1:]
if hashlib.sha256(Path(exe).read_bytes()).hexdigest() != expected:
    raise SystemExit('PROMTOOL_IDENTITY_REJECTED')
for arguments in (['check','rules','ops/prometheus/platform.rules.yml'], ['test','rules','ops/prometheus/platform.rules.test.yml']):
    result = bounded_process([exe,*arguments], Path.cwd(), dict(os.environ), 60)
    sys.stdout.buffer.write(result.stdout)
    sys.stderr.buffer.write(result.stderr)
    if result.returncode:
        raise SystemExit(result.returncode)
'@
    & $rulePython -X utf8 -B -c $code $promtool $PromtoolSha256
  }
  $script:prometheusRuleArtifact = [ordered]@{
    path=$promtool; sha256=$PromtoolSha256
    rule_sha256=(Get-FileHash -LiteralPath (Join-Path $ruleRoot 'ops/prometheus/platform.rules.yml') -Algorithm SHA256).Hash.ToLowerInvariant()
    fixture_sha256=(Get-FileHash -LiteralPath (Join-Path $ruleRoot 'ops/prometheus/platform.rules.test.yml') -Algorithm SHA256).Hash.ToLowerInvariant()
    production_admission=$false
  }
}

function Invoke-LibraryAudit {
  Invoke-Checked 'library-structural-and-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'VERIFY_LIBRARY.ps1')
  }
  $script:implementationPackCount = @(Get-ChildItem -LiteralPath (Join-Path $libraryRoot 'implementation_packs') -Filter '*.md' | Where-Object Name -ne 'README.md').Count
  Invoke-PrometheusRuleRegression

  $readinessRoot = Join-Path $tempRoot 'project-start-readiness-validator'
  Invoke-Checked 'project-readiness-pack-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/PROJECT_START_READINESS_VALIDATOR.md') -Destination $readinessRoot
  }
  $readinessDirectory = Join-Path $readinessRoot 'project_readiness_gate'
  $readinessPython = Resolve-PythonTool
  Invoke-Checked 'project-and-library-readiness-syntax-and-102-regressions' $readinessDirectory {
    & $readinessPython -m py_compile (Join-Path $readinessDirectory 'validate_project_readiness.py') (Join-Path $readinessDirectory 'render_consumer_profile.py') (Join-Path $readinessDirectory 'render_project_advisory.py') (Join-Path $readinessDirectory 'test_validate_project_readiness.py') (Join-Path $readinessDirectory 'test_render_project_advisory.py')
    if ($LASTEXITCODE -ne 0) { Fail 'project readiness py_compile failed' }
    $readinessTestCount = (& $readinessPython -c "import unittest; print(unittest.defaultTestLoader.discover('.', pattern='test_*.py').countTestCases())" | Out-String).Trim()
    if ($LASTEXITCODE -ne 0 -or $readinessTestCount -ne '102') { Fail "project readiness regression count drifted: $readinessTestCount" }
    & $readinessPython -m unittest -q test_render_project_advisory.py test_validate_project_readiness.py test_validate_library_readiness.py
  }

  $referenceTelemetryRoot = Join-Path $tempRoot 'windows-reference-telemetry'
  Invoke-Checked 'windows-reference-telemetry-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/WINDOWS_REFERENCE_TELEMETRY_RUNTIME.md') -Destination $referenceTelemetryRoot
  }
  Invoke-Checked 'windows-reference-telemetry-six-policy-regressions' $referenceTelemetryRoot {
    $referencePolicyReceipt = Join-Path $tempRoot 'reference-policy.json'
    & $readinessPython -B (Join-Path $referenceTelemetryRoot 'reference_telemetry/test_reference_profile.py') $referencePolicyReceipt
    if ($LASTEXITCODE -ne 0) { Fail 'reference telemetry policy regressions failed' }
    $referencePolicy = Get-Content -LiteralPath $referencePolicyReceipt -Raw | ConvertFrom-Json
    if ($referencePolicy.pass -ne $true -or $referencePolicy.tests -ne 6 -or $referencePolicy.skips -ne 0) { Fail 'reference telemetry policy receipt mismatch' }
  }

  $httpMetricsRoot = Join-Path $tempRoot 'http-metrics-reference'
  Invoke-Checked 'http-metrics-reference-eight-file-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/GO_HTTP_METRICS_REFERENCE.md') -Destination $httpMetricsRoot
    if ($LASTEXITCODE -ne 0) { Fail 'HTTP metrics reference materialization failed' }
    if (@(Get-ChildItem -LiteralPath $httpMetricsRoot -Recurse -File).Count -ne 8) { Fail 'HTTP metrics reference file count drifted' }
  }
  Invoke-Checked 'http-metrics-reference-harness-syntax-and-help' $httpMetricsRoot {
    & $readinessPython -B -c "import ast,pathlib; ast.parse(pathlib.Path('reference_http_metrics/run_reference.py').read_text(encoding='utf-8'))"
    if ($LASTEXITCODE -ne 0) { Fail 'HTTP metrics reference syntax failed' }
    & $readinessPython -B (Join-Path $httpMetricsRoot 'reference_http_metrics/run_reference.py') --help
  }

  $historyCompositorRoot = Join-Path $tempRoot 'history-training-compositor'
  Invoke-Checked 'history-training-compositor-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/MARKDOWN_COMPOSITOR_CORE.md') -Destination $historyCompositorRoot
  }
  $historyTrainingRoot = Join-Path $tempRoot 'history-training'
  Invoke-Checked 'history-training-21-file-composition' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $historyCompositorRoot 'tools/compose-markdown-project.ps1') -PlanFile (Join-Path $libraryRoot 'markdown_system/HISTORY_MODEL_TRAINING_PACK_PLAN.md') -LibraryRoot $libraryRoot -Destination $historyTrainingRoot
    if ($LASTEXITCODE -ne 0) { Fail 'history training composition failed' }
    $historyFiles = @(Get-ChildItem -LiteralPath $historyTrainingRoot -Recurse -File | Where-Object Name -ne 'MATERIALIZATION_RECORD.md')
    if ($historyFiles.Count -ne 21) { Fail 'history training file count drifted' }
  }
  Invoke-Checked 'history-training-27-policy-regressions' (Join-Path $historyTrainingRoot 'history_training') {
    $historyTestCount = (& $readinessPython -B -c "import unittest, test_governance; print(unittest.defaultTestLoader.loadTestsFromModule(test_governance).countTestCases())" | Out-String).Trim()
    if ($LASTEXITCODE -ne 0 -or $historyTestCount -ne '27') { Fail 'history training regression count drifted' }
    & $readinessPython -B test_governance.py
  }

  $executionRoot = Join-Path $tempRoot 'engineering-execution-validator'
  Invoke-Checked 'engineering-execution-pack-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/ENGINEERING_EXECUTION_VALIDATOR.md') -Destination $executionRoot
  }
  $executionDirectory = Join-Path $executionRoot 'engineering_execution_kit'
  $executionAuthorities = @(
    'SYSTEMS_ENGINEERING_MASTER_MAP.md',
    'AI_INFERENCE_PERFORMANCE_HARDWARE.md',
    'AI_SECURITY_GOVERNANCE_PRIVACY.md',
    'ALGORITHMS_DATA_STRUCTURES_PROBLEM_SOLVING.md',
    'COMPUTER_SYSTEMS_PERFORMANCE_LOW_LATENCY.md',
    'DATABASE_STORAGE_INTERNALS.md',
    'DEEP_LEARNING_ANDREW_NG.md',
    'FRONTEND_PRODUCT_ENGINEERING_UX.md',
    'MARKET_MICROSTRUCTURE_EXCHANGE_SYSTEMS.md',
    'ML_PRODUCTION_LLMOPS_EVALUATION.md',
    'NATIVE_MOBILE_DESKTOP_ENGINEERING.md',
    'NETWORKING_DISTRIBUTED_STREAMING.md',
    'SECURITY_SRE_CLOUD_INFRASTRUCTURE.md',
    'SOFTWARE_ARCHITECTURE_SYSTEM_DESIGN.md',
    'SOFTWARE_BACKEND_API_ENGINEERING.md',
    'TOOLCHAINS_BUILDS_PACKAGING_FFI.md'
  )
  Invoke-Checked 'engineering-execution-syntax-and-27-regressions' $executionRoot {
    foreach ($authority in $executionAuthorities) {
      Copy-Item -LiteralPath (Join-Path $libraryRoot $authority) -Destination (Join-Path $executionRoot $authority)
    }
    & $readinessPython -m compileall -q $executionDirectory
    if ($LASTEXITCODE -ne 0) { Fail 'engineering execution compileall failed' }
    $executionTestCount = (& $readinessPython -c "import unittest; print(unittest.defaultTestLoader.discover('engineering_execution_kit', pattern='test_*.py').countTestCases())" | Out-String).Trim()
    if ($LASTEXITCODE -ne 0 -or $executionTestCount -ne '27') { Fail "engineering execution regression count drifted: $executionTestCount" }
    & $readinessPython (Join-Path $executionDirectory 'test_execution_state.py')
    if ($LASTEXITCODE -ne 0) { Fail 'engineering execution state suite failed' }
    & $readinessPython (Join-Path $executionDirectory 'test_validate_project.py')
  }
  $script:executionStateControlCount = 1

  $capabilityGapRoot = Join-Path $tempRoot 'capability-gap-resolution-gate'
  Invoke-Checked 'capability-gap-resolution-pack-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/CAPABILITY_GAP_RESOLUTION_GATE.md') -Destination $capabilityGapRoot
  }
  $capabilityGapDirectory = Join-Path $capabilityGapRoot 'capability_gap_resolution'
  Invoke-Checked 'capability-gap-resolution-syntax-and-nine-regressions' $capabilityGapDirectory {
    & $readinessPython -m py_compile (Join-Path $capabilityGapDirectory 'resolve_capability_gap.py') (Join-Path $capabilityGapDirectory 'verify_pack.py')
    if ($LASTEXITCODE -ne 0) { Fail 'capability gap resolution py_compile failed' }
    & $readinessPython (Join-Path $capabilityGapDirectory 'verify_pack.py')
  }
  $script:capabilityGapResolutionGateCount = 1

  $pnpmSelectionRoot = Join-Path $tempRoot 'pnpm-artifact-selection'
  Invoke-Checked 'pnpm-artifact-selection-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/PNPM_ARTIFACT_SELECTION_GATE.md') -Destination $pnpmSelectionRoot
  }
  Invoke-Checked 'pnpm-artifact-selection-seventeen-regressions' (Join-Path $pnpmSelectionRoot 'pnpm_artifact_selection') {
    & $readinessPython -m unittest -q test_selection.py
  }

  Invoke-Checked 'pnpm-consumer-routing-seventy-three-regressions-and-current-profiles' (Join-Path $pnpmSelectionRoot 'pnpm_artifact_selection') {
    & $readinessPython -m unittest -q test_routing.py
    if ($LASTEXITCODE -ne 0) { Fail 'pnpm routing regression suite failed' }
    $pnpmCurrentWebRoot = Join-Path $tempRoot 'pnpm-current-web-consumers'
    & pwsh -NoProfile -File (Join-Path $historyCompositorRoot 'tools/compose-markdown-project.ps1') -PlanFile (Join-Path $libraryRoot 'markdown_system/ENTERPRISE_WEB_PACK_PLAN.md') -LibraryRoot $libraryRoot -Destination $pnpmCurrentWebRoot
    if ($LASTEXITCODE -ne 0) { Fail 'pnpm current consumer composition failed' }
    & $readinessPython -B -c 'from pathlib import Path; import sys; from plan_install import consumer_files; root=Path(sys.argv[1]); consumer_files("enterprise-web", root); consumer_files("playwright", root/"microsoft_playwright_browser_gate"); consumer_files("lighthouse", root/"google_lighthouse_web_quality_gate"); print("PNPM_CURRENT_CONSUMER_PROFILES_PASS count=3 executed=false")' $pnpmCurrentWebRoot
  }

  $devSkimRoot = Join-Path $tempRoot 'microsoft-devskim-adapted-sast-gate'
  Invoke-Checked 'microsoft-devskim-adapted-sast-gate-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/MICROSOFT_DEVSKIM_ADAPTED_SAST_GATE.md') -Destination $devSkimRoot
  }
  $devSkimDirectory = Join-Path $devSkimRoot 'microsoft_devskim_adapted_sast'
  Invoke-Checked 'microsoft-devskim-adapted-sast-static-and-negatives' $devSkimDirectory {
    & pwsh -NoProfile -File (Join-Path $devSkimDirectory 'verify_pack.ps1')
  }
  if ($AllowNetwork) {
    $devSkimDotNet = Resolve-Tool 'dotnet' $DotNetExecutable
    $devSkimDotNetVersion = Get-ToolVersion $devSkimDotNet @('--version')
    if (-not $devSkimDotNet -or $devSkimDotNetVersion -cne '10.0.400') { Fail "AllowNetwork DevSkim audit requires exact .NET SDK 10.0.400; actual=$devSkimDotNetVersion" }
    Invoke-Checked 'microsoft-devskim-adapted-sast-runtime' $devSkimDirectory {
      & pwsh -NoProfile -File (Join-Path $devSkimDirectory 'verify_pack.ps1') -AllowNetwork -DotNetExecutable $devSkimDotNet
    }
  } else {
    Write-Output 'MICROSOFT_DEVSKIM_ADAPTED_SAST_RUNTIME_SKIPPED reason=AllowNetwork_not_authorized'
  }
  $script:microsoftDevSkimAdaptedSastGateCount = 1

  $gitlabOpengrepRoot = Join-Path $tempRoot 'gitlab-opengrep-signed-sast-gate'
  Invoke-Checked 'gitlab-opengrep-signed-sast-gate-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/GITLAB_OPENGREP_SIGNED_SAST_GATE.md') -Destination $gitlabOpengrepRoot
  }
  $gitlabOpengrepDirectory = Join-Path $gitlabOpengrepRoot 'gitlab_opengrep_signed_sast'
  Invoke-Checked 'gitlab-opengrep-signed-sast-static-and-negative' $gitlabOpengrepDirectory {
    & pwsh -NoProfile -File (Join-Path $gitlabOpengrepDirectory 'verify_pack.ps1')
  }
  if ($AllowNetwork) {
    Invoke-Checked 'gitlab-opengrep-signed-sast-current-rejection' $gitlabOpengrepDirectory {
      & pwsh -NoProfile -File (Join-Path $gitlabOpengrepDirectory 'verify_pack.ps1') -AllowNetwork
    }
  } else {
    Write-Output 'GITLAB_OPENGREP_SIGNED_SAST_RUNTIME_BLOCKED reason=Cosign_3.1.3_SCA_rejected'
  }
  $script:gitlabOpengrepSignedSastGateCount = 1

  $portableReleaseRoot = Join-Path $tempRoot 'portable-signed-release-evidence-gate'
  Invoke-Checked 'portable-signed-release-evidence-gate-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/PORTABLE_SIGNED_RELEASE_EVIDENCE_GATE.md') -Destination $portableReleaseRoot
  }
  $portableReleaseDirectory = Join-Path $portableReleaseRoot 'portable_release_evidence_gate'
  Invoke-Checked 'portable-signed-release-evidence-gate-vector-and-tamper' $portableReleaseDirectory {
    & pwsh -NoProfile -File (Join-Path $portableReleaseDirectory 'verify_pack.ps1')
  }
  $script:portableSignedReleaseEvidenceGateCount = 1

  $packRoot = Join-Path $tempRoot 'official-upstream-pack'
  Invoke-Checked 'official-upstream-pack-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/OFFICIAL_UPSTREAM_ACQUISITION_CORE.md') -Destination $packRoot
  }
  Invoke-Checked 'official-upstream-lock-negative-suite' $packRoot {
    & pwsh -NoProfile -File (Join-Path $packRoot 'elite_sources/test_upstream_acquisition.ps1')
  }
  Invoke-Checked 'official-source-profile-suite' $packRoot {
    & pwsh -NoProfile -File (Join-Path $packRoot 'elite_sources/test_source_profiles.ps1')
  }
  Invoke-Checked 'official-opaque-artifact-transport-suite' $packRoot {
    & pwsh -NoProfile -File (Join-Path $packRoot 'elite_sources/test_opaque_artifact_transport.ps1')
  }
  Invoke-Checked 'official-quarantine-wheel-transport-suite' $packRoot {
    & pwsh -NoProfile -File (Join-Path $packRoot 'elite_sources/test_quarantine_wheel_transport.ps1')
  }
  $lock = Get-Content -LiteralPath (Join-Path $packRoot 'elite_sources/upstream-source-lock.json') -Raw | ConvertFrom-Json -Depth 100
  $script:upstreamSourceCount = @($lock.sources).Count
  if ($script:upstreamSourceCount -le 0) { Fail 'official upstream lock contains no sources' }
  Invoke-Checked "official-upstream-lock-$script:upstreamSourceCount" $packRoot {
    & pwsh -NoProfile -File (Join-Path $packRoot 'elite_sources/acquire_upstream_sources.ps1') -LockPath (Join-Path $packRoot 'elite_sources/upstream-source-lock.json') -Destination (Join-Path $tempRoot 'not-acquired') -ValidateOnly
  }

  $kiotaRoot = Join-Path $tempRoot 'microsoft-kiota-openapi-client-gate'
  Invoke-Checked 'microsoft-kiota-openapi-client-gate-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/MICROSOFT_KIOTA_OPENAPI_CLIENT_GATE.md') -Destination $kiotaRoot
  }
  $kiotaDirectory = Join-Path $kiotaRoot 'microsoft_kiota_openapi_gate'
  Invoke-Checked 'microsoft-kiota-openapi-client-gate-source-contract' $kiotaDirectory {
    $kiotaLock = Get-Content -LiteralPath (Join-Path $kiotaDirectory 'source-lock.json') -Raw | ConvertFrom-Json -Depth 30
    $kiotaAcquisition = @($lock.sources | Where-Object id -CEQ 'microsoft-kiota-1.35.0')
    if ($kiotaAcquisition.Count -ne 1) { Fail 'Microsoft Kiota 1.35.0 source missing or duplicated in acquisition lock' }
    if ($kiotaLock.source.commit -cne $kiotaAcquisition[0].commit -or $kiotaLock.source.archive_sha256 -cne $kiotaAcquisition[0].archive_sha256 -or $kiotaLock.windows_x64.sha256 -cne $kiotaAcquisition[0].artifact_evidence.asset_sha256 -or $kiotaLock.windows_x64.executable_sha256 -cne $kiotaAcquisition[0].artifact_evidence.executable_sha256) { Fail 'Microsoft Kiota pack and acquisition identities differ' }
    if ($kiotaLock.verified_lane.language -cne 'Go' -or $kiotaLock.verified_lane.generated_code_files -ne 6 -or $kiotaLock.verified_lane.osv_findings_at_verification -ne 0) { Fail 'Microsoft Kiota verified lane drifted' }
    foreach ($scriptName in @('acquire_kiota.ps1','generate_go_client.ps1','verify_contract.ps1')) {
      $tokens = $null
      $errors = $null
      [void][Management.Automation.Language.Parser]::ParseFile((Join-Path $kiotaDirectory $scriptName), [ref]$tokens, [ref]$errors)
      if (@($errors).Count -ne 0) { Fail "Microsoft Kiota script syntax failed: $scriptName" }
    }
  }
  if ($AllowNetwork) {
    $kiotaAuditGo = Resolve-Tool 'go' $GoExecutable
    $kiotaAuditGoVersion = Get-ToolVersion $kiotaAuditGo @('version')
    if (-not $kiotaAuditGo -or -not (Test-MinimumVersion $kiotaAuditGoVersion ([version]'1.26.7'))) { Fail 'AllowNetwork Kiota audit requires Go 1.26.7+' }
    $osvDirectory = Join-Path $tempRoot 'google-osv-scanner'
    [void](New-Item -ItemType Directory -Path $osvDirectory)
    $osvExecutable = Join-Path $osvDirectory 'osv-scanner_windows_amd64.exe'
    Invoke-WebRequest -Uri 'https://github.com/google/osv-scanner/releases/download/v2.5.1/osv-scanner_windows_amd64.exe' -OutFile $osvExecutable
    if ((Get-Item -LiteralPath $osvExecutable).Length -ne 58970112 -or (Get-FileHash -LiteralPath $osvExecutable -Algorithm SHA256).Hash.ToLowerInvariant() -cne '25e42f5ef6711fd8c0fb45390972205891dd44c6bd02ac93f0f63e8e98d9bfb6') { Fail 'Google OSV Scanner artifact identity mismatch for Kiota audit' }
    Invoke-Checked 'microsoft-kiota-openapi-client-gate-runtime' $kiotaDirectory {
      & pwsh -NoProfile -File (Join-Path $kiotaDirectory 'verify_contract.ps1') -GoExecutable $kiotaAuditGo -OsvScanner $osvExecutable
    }
  } else {
    Write-Output 'MICROSOFT_KIOTA_OPENAPI_RUNTIME_SKIPPED reason=network_not_enabled'
  }
  $script:microsoftKiotaOpenApiClientGateCount = 1

  $goFuzzRoot = Join-Path $tempRoot 'go-native-fuzz-gate'
  Invoke-Checked 'go-native-fuzz-gate-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/GO_NATIVE_FUZZ_GATE.md') -Destination $goFuzzRoot
  }
  $goFuzzDirectory = Join-Path $goFuzzRoot 'go_fuzz_gate'
  Invoke-Checked 'go-native-fuzz-gate-source-contract' $goFuzzDirectory {
    foreach ($scriptName in @('run_go_fuzz_gate.ps1','verify_pack.ps1')) {
      $tokens = $null; $errors = $null
      [void][Management.Automation.Language.Parser]::ParseFile((Join-Path $goFuzzDirectory $scriptName), [ref]$tokens, [ref]$errors)
      if (@($errors).Count -ne 0) { Fail "Go fuzz gate syntax failed: $scriptName" }
    }
    $fuzzTemplate = Get-Content -LiteralPath (Join-Path $goFuzzDirectory 'profile.example.json') -Raw | ConvertFrom-Json
    if ($fuzzTemplate.enabled -ne $false -or @($fuzzTemplate.targets).Count -ne 0) { Fail 'Go fuzz template must enable nothing' }
  }
  if ($Mode -eq 'Audit') {
    $goFuzzTool = Resolve-Tool 'go' $GoExecutable
    $goFuzzVersion = Get-ToolVersion $goFuzzTool @('version')
    if (-not $goFuzzTool -or -not (Test-MinimumVersion $goFuzzVersion ([version]'1.26.7'))) { Fail 'Audit Go fuzz gate requires Go 1.26.7+' }
    Invoke-Checked 'go-native-fuzz-gate-runtime' $goFuzzDirectory {
      & pwsh -NoProfile -File (Join-Path $goFuzzDirectory 'verify_pack.ps1') -GoExecutable $goFuzzTool
    }
  } else {
    Write-Output 'GO_NATIVE_FUZZ_RUNTIME_SKIPPED reason=preflight_mode'
  }
  $script:goNativeFuzzGateCount = 1

  $playwrightRoot = Join-Path $tempRoot 'microsoft-playwright-browser-gate'
  Invoke-Checked 'microsoft-playwright-browser-gate-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/MICROSOFT_PLAYWRIGHT_BROWSER_GATE.md') -Destination $playwrightRoot
  }
  $playwrightDirectory = Join-Path $playwrightRoot 'microsoft_playwright_browser_gate'
  Invoke-Checked 'microsoft-playwright-browser-gate-source-contract' $playwrightDirectory {
    $playwrightLock = Get-Content -LiteralPath (Join-Path $playwrightDirectory 'source-lock.json') -Raw | ConvertFrom-Json -Depth 30
    $playwrightAcquisition = @($lock.sources | Where-Object id -CEQ 'microsoft-playwright-1.62.1')
    if ($playwrightAcquisition.Count -ne 1) { Fail 'Microsoft Playwright source missing or duplicated in acquisition lock' }
    if ($playwrightLock.source.commit -cne $playwrightAcquisition[0].commit -or $playwrightLock.source.archiveSha256 -cne $playwrightAcquisition[0].archive_sha256 -or $playwrightLock.npm.tarballSha256 -cne '009534220efd98c0361d8c4ee7e3db1ed510ab88a23e98b5081ef1c8fed64965') { Fail 'Microsoft Playwright pack and acquisition identities differ' }
    if (@($playwrightLock.browsers).Count -ne 3 -or $playwrightLock.browsers[0].revision -cne '1234' -or $playwrightLock.browsers[1].revision -cne '1538' -or $playwrightLock.browsers[2].revision -cne '2336') { Fail 'Microsoft Playwright browser matrix drifted' }
    $playwrightLicense = Join-Path $playwrightDirectory 'LICENSE.playwright.txt'
    if ((Get-Item -LiteralPath $playwrightLicense).Length -ne 11399 -or (Get-FileHash -LiteralPath $playwrightLicense -Algorithm SHA256).Hash.ToLowerInvariant() -cne '7fab1461b41970ff376f1c9303a637076bfaaeb71cd12dd3a1c44aaf59a1a2b9') { Fail 'Microsoft Playwright packaged license identity drifted' }
  }
  $auditNode = Resolve-Tool 'node' $NodeExecutable
  $auditPnpm = Resolve-Tool 'pnpm' $PnpmExecutable
  if (-not $auditNode -or -not $auditPnpm) { Fail 'Node.js and pnpm are required for the Microsoft Playwright runtime gate' }
  Invoke-Checked 'microsoft-playwright-browser-gate-frozen-offline-install' $playwrightDirectory {
    & $auditPnpm install --ignore-workspace --frozen-lockfile --offline
  }
  Invoke-Checked 'microsoft-playwright-browser-gate-runtime' $playwrightDirectory {
    & pwsh -NoProfile -File (Join-Path $playwrightDirectory 'verify_contract.ps1')
  }
  $script:microsoftPlaywrightBrowserGateCount = 1

  $lighthouseRoot = Join-Path $tempRoot 'google-lighthouse-web-quality-gate'
  Invoke-Checked 'google-lighthouse-web-quality-gate-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/GOOGLE_LIGHTHOUSE_WEB_QUALITY_GATE.md') -Destination $lighthouseRoot
  }
  $lighthouseDirectory = Join-Path $lighthouseRoot 'google_lighthouse_web_quality_gate'
  Invoke-Checked 'google-lighthouse-web-quality-gate-source-contract' $lighthouseDirectory {
    $lighthouseLock = Get-Content -LiteralPath (Join-Path $lighthouseDirectory 'source-lock.json') -Raw | ConvertFrom-Json -Depth 30
    $lighthouseAcquisition = @($lock.sources | Where-Object id -CEQ 'google-lighthouse-13.4.1')
    if ($lighthouseAcquisition.Count -ne 1) { Fail 'Google Lighthouse source missing or duplicated in acquisition lock' }
    if ($lighthouseLock.source.commit -cne $lighthouseAcquisition[0].commit -or $lighthouseLock.source.archiveSha256 -cne $lighthouseAcquisition[0].archive_sha256 -or $lighthouseLock.npm.tarballSha256 -cne $lighthouseAcquisition[0].artifact_evidence.npm_tarball_sha256 -or $lighthouseLock.npm.attestationsSha256 -cne $lighthouseAcquisition[0].artifact_evidence.npm_attestation_sha256) { Fail 'Google Lighthouse pack and acquisition identities differ' }
    $lighthouseLicense = Join-Path $lighthouseDirectory 'LICENSE.lighthouse.txt'
    if ((Get-Item -LiteralPath $lighthouseLicense).Length -ne 11358 -or (Get-FileHash -LiteralPath $lighthouseLicense -Algorithm SHA256).Hash.ToLowerInvariant() -cne 'cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30') { Fail 'Google Lighthouse packaged license identity drifted' }
  }
  Invoke-Checked 'google-lighthouse-web-quality-gate-frozen-offline-install' $lighthouseDirectory {
    & $auditPnpm install --ignore-workspace --frozen-lockfile --offline
  }
  Invoke-Checked 'google-lighthouse-web-quality-gate-runtime-contract' $lighthouseDirectory {
    & pwsh -NoProfile -File (Join-Path $lighthouseDirectory 'verify_contract.ps1')
  }
  $script:googleLighthouseWebQualityGateCount = 1

  $powertoolsBatchRoot = Join-Path $tempRoot 'aws-powertools-idempotent-sqs-batch'
  Invoke-Checked 'aws-powertools-idempotent-sqs-batch-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/AWS_POWERTOOLS_IDEMPOTENT_SQS_BATCH_COMPONENT.md') -Destination $powertoolsBatchRoot
  }
  $powertoolsBatchDirectory = Join-Path $powertoolsBatchRoot 'aws_powertools_idempotent_sqs_batch'
  Invoke-Checked 'aws-powertools-idempotent-sqs-batch-contracts' $powertoolsBatchDirectory {
    & $readinessPython -m unittest -v test_contracts.py
  }
  $powertoolsBatchLock = Get-Content -LiteralPath (Join-Path $powertoolsBatchDirectory 'source-lock.json') -Raw | ConvertFrom-Json -Depth 20
  $powertoolsAcquisition = @($lock.sources | Where-Object id -CEQ 'aws-powertools-python-3.34.0')
  if ($powertoolsAcquisition.Count -ne 1) { Fail 'AWS Powertools Python source missing or duplicated in acquisition lock' }
  if ($powertoolsBatchLock.commit -cne $powertoolsAcquisition[0].commit -or $powertoolsBatchLock.archive_sha256 -cne $powertoolsAcquisition[0].archive_sha256 -or $powertoolsBatchLock.wheel_sha256 -cne $powertoolsAcquisition[0].wheel_sha256 -or $powertoolsBatchLock.license_expression -cne $powertoolsAcquisition[0].license_expression) { Fail 'AWS Powertools SQS component and acquisition source identities differ' }
  $script:awsPowertoolsBatchComponentCount = 1

  $durableExecutionRoot = Join-Path $tempRoot 'aws-lambda-durable-execution'
  Invoke-Checked 'aws-lambda-durable-execution-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/AWS_LAMBDA_DURABLE_EXECUTION_COMPONENT.md') -Destination $durableExecutionRoot
  }
  $durableExecutionDirectory = Join-Path $durableExecutionRoot 'aws_lambda_durable_execution'
  Invoke-Checked 'aws-lambda-durable-execution-contracts' $durableExecutionDirectory {
    & $readinessPython -m unittest -v test_contracts.py
  }
  $durableExecutionLock = Get-Content -LiteralPath (Join-Path $durableExecutionDirectory 'source-lock.json') -Raw | ConvertFrom-Json -Depth 20
  if ($durableExecutionLock.source_id -cne 'aws-lambda-durable-execution-python-1.7.0' -or $durableExecutionLock.commit -cne '075b65aacb80de8bb1507e3df2e53cc90cb3b874' -or $durableExecutionLock.archive_sha256 -cne 'd3e42423e61e258191fdd078a29d1fd58cfc55380b25a357f879b73665a02d11') { Fail 'AWS Durable Execution source identity drifted' }
  if ($durableExecutionLock.core_wheel.sha256 -cne '1482e1f439e36deb1fcf9a086499b3788c3facf3ac755e50376a739081cc72a0' -or @($durableExecutionLock.files).Count -ne 8) { Fail 'AWS Durable Execution wheel or verbatim inventory drifted' }
  if ($durableExecutionLock.testing_wheel_condition.classification -cne 'INCOMPATIBLE_WITH_SELECTED_TAG_TEST_TREE' -or $durableExecutionLock.testing_wheel_condition.source_execution_py_sha256 -ceq $durableExecutionLock.testing_wheel_condition.wheel_execution_py_sha256) { Fail 'AWS Durable Execution testing-wheel condition drifted' }
  $script:awsDurableExecutionComponentCount = 1

  $textractorOfficialRoot = Join-Path $tempRoot 'aws-textractor-official'
  Invoke-Checked 'aws-textractor-official-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/PYTHON_AWS_TEXTRACTOR_OFFICIAL_COMPONENT.md') -Destination $textractorOfficialRoot
  }
  $textractorOfficialDirectory = Join-Path $textractorOfficialRoot 'aws_textractor_official'
  Invoke-Checked 'aws-textractor-official-contracts' $textractorOfficialDirectory {
    & $readinessPython -m unittest -v test_contracts.py
  }
  $textractorOfficialLock = Get-Content -LiteralPath (Join-Path $textractorOfficialDirectory 'source-lock.json') -Raw | ConvertFrom-Json -Depth 20
  if ($textractorOfficialLock.repository -cne 'aws-samples/amazon-textract-textractor' -or $textractorOfficialLock.release -cne 'v1.10.0' -or $textractorOfficialLock.commit -cne '8ea5f9ae65dbcbb2b83b409e0007c395e59b0e92') { Fail 'AWS Textractor official source identity drifted' }
  if ($textractorOfficialLock.archive.sha256 -cne '4fa38999e355d977d2b46816b76c9d32ab6c5d8607987e086045470d101c7219' -or $textractorOfficialLock.wheel.sha256 -cne '8524224f07a776ca2959e0d80d5e031a30bf53938c74e8352848f3d0b6b3762c' -or @($textractorOfficialLock.files).Count -ne 7) { Fail 'AWS Textractor archive, wheel or official-file inventory drifted' }
  if ($textractorOfficialLock.verification.linux_deterministic_core.passed -ne 69 -or $textractorOfficialLock.verification.runtime_osv.packages -ne 15 -or $textractorOfficialLock.verification.runtime_osv.findings -ne 0 -or $textractorOfficialLock.verification.test_osv.packages -ne 26 -or $textractorOfficialLock.verification.test_osv.findings -ne 0) { Fail 'AWS Textractor deterministic-suite or SCA evidence drifted' }
  $script:awsTextractorOfficialComponentCount = 1

  $tesseraRoot = Join-Path $tempRoot 'transparency-dev-tessera-posix-evidence-log'
  Invoke-Checked 'transparency-dev-tessera-pack-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/TRANSPARENCY_DEV_TESSERA_POSIX_EVIDENCE_LOG.md') -Destination $tesseraRoot
  }
  $tesseraDirectory = Join-Path $tesseraRoot 'tessera_posix'
  $auditPython = Resolve-PythonTool
  if (-not $auditPython) { Fail 'python is required for Tessera static contract regressions' }
  Invoke-Checked 'transparency-dev-tessera-static-contracts' $tesseraDirectory {
    & $auditPython -m unittest -v test_contracts.py
  }
  $tesseraLock = Get-Content -LiteralPath (Join-Path $tesseraDirectory 'source-lock.json') -Raw | ConvertFrom-Json -Depth 20
  $tesseraAcquisition = $lock.sources | Where-Object id -CEQ 'transparency-dev-tessera-a8f33c5'
  if ($null -eq $tesseraAcquisition) { Fail 'Tessera source missing from official acquisition lock' }
  if ($tesseraLock.source.commit -cne $tesseraAcquisition.commit -or $tesseraLock.source.archive_sha256 -cne $tesseraAcquisition.archive_sha256 -or $tesseraLock.source.license_sha256 -cne $tesseraAcquisition.license_sha256) { Fail 'Tessera pack and acquisition source identities differ' }
  if ($tesseraLock.target.os -cne 'linux' -or $tesseraLock.target.architecture -cne 'amd64' -or $tesseraLock.evidence.govulncheck_reachable -ne 0 -or $tesseraLock.evidence.govulncheck_imported_packages -ne 0 -or $tesseraLock.evidence.osv_records -le 0) { Fail 'Tessera conditioned platform or vulnerability evidence drifted' }
  $script:evidenceLogCount = 1
  Write-Output 'TRANSPARENCY_DEV_TESSERA_LINUX_RUNTIME_GATES_SKIPPED reason=Audit_is_offline_cross_platform; exact dated Linux build evidence remains required'

  $documentPackRoot = Join-Path $tempRoot 'official-document-sdk-pack'
  Invoke-Checked 'official-document-sdk-pack-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/OFFICIAL_DOCUMENT_SDK_ARTIFACT_CORE.md') -Destination $documentPackRoot
  }
  Invoke-Checked 'official-document-sdk-acquisition-suite' $documentPackRoot {
    & pwsh -NoProfile -File (Join-Path $documentPackRoot 'official_document_sdks/test_artifact_acquisition.ps1')
  }
  $documentLock = Get-Content -LiteralPath (Join-Path $documentPackRoot 'official_document_sdks/artifact-lock.json') -Raw | ConvertFrom-Json -Depth 20
  $script:documentArtifactCount = @($documentLock.artifacts).Count
  if ($script:documentArtifactCount -le 0) { Fail 'official document SDK lock contains no artifacts' }
  Invoke-Checked "official-document-sdk-lock-$script:documentArtifactCount" $documentPackRoot {
    & pwsh -NoProfile -File (Join-Path $documentPackRoot 'official_document_sdks/acquire_document_sdk_artifacts.ps1') -ArtifactId 'azure-ai-contentunderstanding-1.1.0-sdist' -ValidateOnly
  }

  $officialInvoiceRoot = Join-Path $tempRoot 'microsoft-azure-document-intelligence-official-invoice'
  Invoke-Checked 'microsoft-azure-di-official-invoice-pack-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/MICROSOFT_AZURE_DOCUMENT_INTELLIGENCE_OFFICIAL_INVOICE_SAMPLE.md') -Destination $officialInvoiceRoot
  }
  $officialInvoiceDirectory = Join-Path $officialInvoiceRoot 'azure_document_intelligence_official_invoice'
  Invoke-Checked 'microsoft-azure-di-official-invoice-sample-regression' $officialInvoiceDirectory {
    & $auditPython -m unittest -v test_official_invoice_sample.py
  }
  Invoke-Checked 'microsoft-azure-di-official-invoice-fixture-regression' $officialInvoiceDirectory {
    & pwsh -NoProfile -File (Join-Path $officialInvoiceDirectory 'test_fixture_acquisition.ps1')
  }
  $officialSamplePath = Join-Path $officialInvoiceDirectory 'upstream/sample_analyze_invoices_from_bytes_source.py'
  $officialLicensePath = Join-Path $officialInvoiceDirectory 'upstream/LICENSE.txt'
  if ((Get-Item -LiteralPath $officialSamplePath).Length -ne 13084 -or (Get-FileHash -LiteralPath $officialSamplePath -Algorithm SHA256).Hash.ToLowerInvariant() -cne 'ebc32b0fc534e625b15e636c6b55376f01d28a83f35217dd141202073297a293') { Fail 'Microsoft Azure DI official invoice sample identity mismatch' }
  if ((Get-Item -LiteralPath $officialLicensePath).Length -ne 1074 -or (Get-FileHash -LiteralPath $officialLicensePath -Algorithm SHA256).Hash.ToLowerInvariant() -cne '7c77a44a8acd9b41fdc209864a8016b3d430b5d0e09309818d5b7444336df744') { Fail 'Microsoft Azure DI official invoice license identity mismatch' }
  $officialFixtureLock = Get-Content -LiteralPath (Join-Path $officialInvoiceDirectory 'fixture-lock.json') -Raw | ConvertFrom-Json -Depth 20
  if ($officialFixtureLock.source.commit -cne '8555d14532a9688b751d8408d822d1dd5feb47f6' -or @($officialFixtureLock.fixtures).Count -ne 1 -or $officialFixtureLock.fixtures[0].sha256 -cne '489f0c63b6a05e0ae0fcfbf29121299fc2ccacab35d5cd40a5d63d42764078cb') { Fail 'Microsoft Azure DI official invoice fixture lock identity mismatch' }
  $script:officialInvoiceSampleCount = 1

  $officialAzureCuRoot = Join-Path $tempRoot 'microsoft-azure-content-understanding-official-invoice'
  Invoke-Checked 'microsoft-azure-cu-official-invoice-pack-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/MICROSOFT_AZURE_CONTENT_UNDERSTANDING_OFFICIAL_INVOICE_SAMPLE.md') -Destination $officialAzureCuRoot
  }
  $officialAzureCuDirectory = Join-Path $officialAzureCuRoot 'azure_content_understanding_official_invoice'
  Invoke-Checked 'microsoft-azure-cu-official-invoice-sample-regression' $officialAzureCuDirectory {
    & $auditPython -m unittest -v test_official_sample.py
  }
  $officialAzureCuSamplePath = Join-Path $officialAzureCuDirectory 'upstream/sample_analyze_invoice.py'
  $officialAzureCuTestPath = Join-Path $officialAzureCuDirectory 'test_official_sample.py'
  $officialAzureCuSamplePyc = Join-Path $tempRoot 'azure-cu-official-sample.pyc'
  $officialAzureCuTestPyc = Join-Path $tempRoot 'azure-cu-official-test.pyc'
  Invoke-Checked 'microsoft-azure-cu-official-invoice-python-syntax' $officialAzureCuDirectory {
    & $auditPython -c 'import py_compile, sys; py_compile.compile(sys.argv[1], cfile=sys.argv[3], doraise=True); py_compile.compile(sys.argv[2], cfile=sys.argv[4], doraise=True)' $officialAzureCuSamplePath $officialAzureCuTestPath $officialAzureCuSamplePyc $officialAzureCuTestPyc
  }
  if (-not (Test-Path -LiteralPath $officialAzureCuSamplePyc -PathType Leaf) -or -not (Test-Path -LiteralPath $officialAzureCuTestPyc -PathType Leaf)) { Fail 'Microsoft Azure CU official invoice syntax bytecode missing' }
  $officialAzureCuLicensePath = Join-Path $officialAzureCuDirectory 'upstream/LICENSE.txt'
  $officialAzureCuLockPath = Join-Path $officialAzureCuDirectory 'source-lock.json'
  if ((Get-Item -LiteralPath $officialAzureCuSamplePath).Length -ne 10539 -or (Get-FileHash -LiteralPath $officialAzureCuSamplePath -Algorithm SHA256).Hash.ToLowerInvariant() -cne '0cb9d7b0e183cd1c3f677c79a9d8d5dd9a1ed30bf397f84b753eedd7819c53fd') { Fail 'Microsoft Azure CU official invoice sample identity mismatch' }
  if ((Get-Item -LiteralPath $officialAzureCuLicensePath).Length -ne 1074 -or (Get-FileHash -LiteralPath $officialAzureCuLicensePath -Algorithm SHA256).Hash.ToLowerInvariant() -cne '7c77a44a8acd9b41fdc209864a8016b3d430b5d0e09309818d5b7444336df744') { Fail 'Microsoft Azure CU official invoice license identity mismatch' }
  if ((Get-FileHash -LiteralPath $officialAzureCuLockPath -Algorithm SHA256).Hash.ToLowerInvariant() -cne '8ec2f600bbc996d126b19668c0c3f28701032b451b9f545ace50805be277814f') { Fail 'Microsoft Azure CU official invoice source lock identity mismatch' }
  $officialAzureCuLock = Get-Content -LiteralPath $officialAzureCuLockPath -Raw | ConvertFrom-Json -Depth 20
  $officialAzureCuFiles = @($officialAzureCuLock.files)
  $officialAzureCuSampleLock = @($officialAzureCuFiles | Where-Object path -CEQ 'upstream/sample_analyze_invoice.py')
  if ($officialAzureCuLock.commit -cne '129a5cbb06fba44e11ffc92f6bcdca14a2f5b6eb' -or $officialAzureCuFiles.Count -ne 2 -or $officialAzureCuSampleLock.Count -ne 1 -or $officialAzureCuSampleLock[0].sha256 -cne '0cb9d7b0e183cd1c3f677c79a9d8d5dd9a1ed30bf397f84b753eedd7819c53fd' -or $officialAzureCuLock.official_claim.analyzer_id -cne 'prebuilt-invoice') { Fail 'Microsoft Azure CU official invoice lock contract mismatch' }
  $script:officialAzureCuSampleCount = 1

  $officialAzureCuBinaryRoot = Join-Path $tempRoot 'microsoft-azure-content-understanding-official-binary'
  Invoke-Checked 'microsoft-azure-cu-official-binary-pack-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/MICROSOFT_AZURE_CONTENT_UNDERSTANDING_OFFICIAL_BINARY_DOCUMENT_SAMPLE.md') -Destination $officialAzureCuBinaryRoot
  }
  $officialAzureCuBinaryDirectory = Join-Path $officialAzureCuBinaryRoot 'azure_content_understanding_official_binary'
  Invoke-Checked 'microsoft-azure-cu-official-binary-regression' $officialAzureCuBinaryDirectory {
    & $auditPython test_official_sample.py
  }
  $officialAzureCuBinarySyntaxFiles = @(
    (Join-Path $officialAzureCuBinaryDirectory 'upstream/sample_analyze_binary.py'),
    (Join-Path $officialAzureCuBinaryDirectory 'upstream/test_sample_analyze_binary.py'),
    (Join-Path $officialAzureCuBinaryDirectory 'test_official_sample.py')
  )
  for ($binarySyntaxIndex = 0; $binarySyntaxIndex -lt $officialAzureCuBinarySyntaxFiles.Count; $binarySyntaxIndex++) {
    $binaryPyc = Join-Path $tempRoot ("azure-cu-binary-{0}.pyc" -f $binarySyntaxIndex)
    Invoke-Checked "microsoft-azure-cu-official-binary-syntax-$binarySyntaxIndex" $officialAzureCuBinaryDirectory {
      & $auditPython -c 'import py_compile, sys; py_compile.compile(sys.argv[1], cfile=sys.argv[2], doraise=True)' $officialAzureCuBinarySyntaxFiles[$binarySyntaxIndex] $binaryPyc
    }
    if (-not (Test-Path -LiteralPath $binaryPyc -PathType Leaf)) { Fail "Microsoft Azure CU binary syntax bytecode $binarySyntaxIndex missing" }
  }
  $officialAzureCuBinarySamplePath = Join-Path $officialAzureCuBinaryDirectory 'upstream/sample_analyze_binary.py'
  $officialAzureCuBinaryOfficialTestPath = Join-Path $officialAzureCuBinaryDirectory 'upstream/test_sample_analyze_binary.py'
  $officialAzureCuBinaryLicensePath = Join-Path $officialAzureCuBinaryDirectory 'upstream/LICENSE.txt'
  if ((Get-Item -LiteralPath $officialAzureCuBinarySamplePath).Length -ne 6979 -or (Get-FileHash -LiteralPath $officialAzureCuBinarySamplePath -Algorithm SHA256).Hash.ToLowerInvariant() -cne 'dce0d3c2684bb0d5bca01f1015e601d4f15407b8272c72b773e2cfb822588893') { Fail 'Microsoft Azure CU official binary sample identity mismatch' }
  if ((Get-Item -LiteralPath $officialAzureCuBinaryOfficialTestPath).Length -ne 20749 -or (Get-FileHash -LiteralPath $officialAzureCuBinaryOfficialTestPath -Algorithm SHA256).Hash.ToLowerInvariant() -cne '88f8cb38d39a351639c835f246da5c443d5c50cc068a363032ab43ee229b276f') { Fail 'Microsoft Azure CU official binary test identity mismatch' }
  if ((Get-Item -LiteralPath $officialAzureCuBinaryLicensePath).Length -ne 1074 -or (Get-FileHash -LiteralPath $officialAzureCuBinaryLicensePath -Algorithm SHA256).Hash.ToLowerInvariant() -cne '7c77a44a8acd9b41fdc209864a8016b3d430b5d0e09309818d5b7444336df744') { Fail 'Microsoft Azure CU official binary license identity mismatch' }
  $officialAzureCuBinaryLock = Get-Content -LiteralPath (Join-Path $officialAzureCuBinaryDirectory 'source-lock.json') -Raw | ConvertFrom-Json -Depth 20
  $officialAzureCuBinaryFiles = @($officialAzureCuBinaryLock.files)
  if ($officialAzureCuBinaryLock.revision -cne '129a5cbb06fba44e11ffc92f6bcdca14a2f5b6eb' -or $officialAzureCuBinaryLock.package_version -cne '1.1.0' -or $officialAzureCuBinaryLock.sdist.sha256 -cne '00f39c7cf2ba50c8d4586315f4295a45fad0881daaff07d13d33aea2bafc7ec8' -or $officialAzureCuBinaryFiles.Count -ne 3) { Fail 'Microsoft Azure CU official binary lock contract mismatch' }
  $script:officialAzureCuBinarySampleCount = 1

  $officialAzureCuCopyRoot = Join-Path $tempRoot 'microsoft-azure-content-understanding-official-analyzer-copy'
  Invoke-Checked 'microsoft-azure-cu-official-analyzer-copy-pack-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/MICROSOFT_AZURE_CONTENT_UNDERSTANDING_OFFICIAL_ANALYZER_COPY_SAMPLE.md') -Destination $officialAzureCuCopyRoot
  }
  $officialAzureCuCopyDirectory = Join-Path $officialAzureCuCopyRoot 'azure_content_understanding_official_analyzer_copy'
  Invoke-Checked 'microsoft-azure-cu-official-analyzer-copy-regression' $officialAzureCuCopyDirectory {
    & $auditPython -m unittest -v test_official_copy_sample.py
  }
  $officialAzureCuCopySyntaxFiles = @(
    (Join-Path $officialAzureCuCopyDirectory 'upstream/sample_copy_analyzer.py'),
    (Join-Path $officialAzureCuCopyDirectory 'upstream/test_sample_copy_analyzer.py'),
    (Join-Path $officialAzureCuCopyDirectory 'test_official_copy_sample.py')
  )
  for ($copySyntaxIndex = 0; $copySyntaxIndex -lt $officialAzureCuCopySyntaxFiles.Count; $copySyntaxIndex++) {
    $copyPyc = Join-Path $tempRoot ("azure-cu-copy-{0}.pyc" -f $copySyntaxIndex)
    Invoke-Checked "microsoft-azure-cu-official-analyzer-copy-syntax-$copySyntaxIndex" $officialAzureCuCopyDirectory {
      & $auditPython -c 'import py_compile, sys; py_compile.compile(sys.argv[1], cfile=sys.argv[2], doraise=True)' $officialAzureCuCopySyntaxFiles[$copySyntaxIndex] $copyPyc
    }
    if (-not (Test-Path -LiteralPath $copyPyc -PathType Leaf)) { Fail "Microsoft Azure CU analyzer copy syntax bytecode $copySyntaxIndex missing" }
  }
  $officialAzureCuCopyLockPath = Join-Path $officialAzureCuCopyDirectory 'source-lock.json'
  $officialAzureCuCopyLock = Get-Content -LiteralPath $officialAzureCuCopyLockPath -Raw | ConvertFrom-Json -Depth 20
  $officialAzureCuCopyFiles = @($officialAzureCuCopyLock.files)
  $officialAzureCuCopySample = @($officialAzureCuCopyFiles | Where-Object path -CEQ 'upstream/sample_copy_analyzer.py')
  $officialAzureCuCopyTest = @($officialAzureCuCopyFiles | Where-Object path -CEQ 'upstream/test_sample_copy_analyzer.py')
  if ($officialAzureCuCopyLock.commit -cne '129a5cbb06fba44e11ffc92f6bcdca14a2f5b6eb' -or $officialAzureCuCopyFiles.Count -ne 3 -or $officialAzureCuCopySample.Count -ne 1 -or $officialAzureCuCopySample[0].sha256 -cne '4b8aa70c941c9b0b8fcd54ddac588ef0d1a6dbd4cefe5e2c21d11fd8e44bdb71' -or $officialAzureCuCopyTest.Count -ne 1 -or $officialAzureCuCopyTest[0].sha256 -cne '227d59a71c0e1ecaa67a9eb406654c79bc6a1294334a1d1aca6eceeb58c18f8e' -or $officialAzureCuCopyLock.official_claim.source_base_analyzer_id -cne 'prebuilt-document' -or $officialAzureCuCopyLock.official_claim.automatic_execution -ne $false) { Fail 'Microsoft Azure CU official analyzer copy lock contract mismatch' }
  $script:officialAzureCuCopySampleCount = 1

  $officialGoogleRoot = Join-Path $tempRoot 'google-cloud-document-ai-official-process-sample'
  Invoke-Checked 'google-cloud-docai-official-process-pack-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/GOOGLE_CLOUD_DOCUMENT_AI_OFFICIAL_PROCESS_SAMPLE.md') -Destination $officialGoogleRoot
  }
  $officialGoogleDirectory = Join-Path $officialGoogleRoot 'google_document_ai_official_process_sample'
  Invoke-Checked 'google-cloud-docai-official-process-regression' $officialGoogleDirectory {
    & $auditPython -m unittest -v test_official_process_document_sample.py
  }
  $officialGoogleCustomSamplePyc = Join-Path $tempRoot 'google-docai-custom-sample.pyc'
  $officialGoogleCustomTestPyc = Join-Path $tempRoot 'google-docai-custom-test.pyc'
  $officialGoogleCustomContractPyc = Join-Path $tempRoot 'google-docai-custom-contract.pyc'
  Invoke-Checked 'google-cloud-docai-official-custom-extractor-syntax' $officialGoogleDirectory {
    & $auditPython -c 'import py_compile, sys; py_compile.compile(sys.argv[1], cfile=sys.argv[4], doraise=True); py_compile.compile(sys.argv[2], cfile=sys.argv[5], doraise=True); py_compile.compile(sys.argv[3], cfile=sys.argv[6], doraise=True)' 'upstream/google/documentai/snippets/handle_response_sample.py' 'upstream/google/documentai/snippets/handle_response_sample_test.py' 'test_official_custom_extractor_sample.py' $officialGoogleCustomSamplePyc $officialGoogleCustomTestPyc $officialGoogleCustomContractPyc
  }
  foreach ($compiled in @($officialGoogleCustomSamplePyc,$officialGoogleCustomTestPyc,$officialGoogleCustomContractPyc)) {
    if (-not (Test-Path -LiteralPath $compiled -PathType Leaf)) { Fail "Google Custom Document Extractor bytecode missing: $compiled" }
  }
  Invoke-Checked 'google-cloud-docai-official-fixture-regression' $officialGoogleDirectory {
    & pwsh -NoProfile -File (Join-Path $officialGoogleDirectory 'test_fixture_acquisition.ps1')
  }
  $officialGoogleValidatorPyc = Join-Path $tempRoot 'google-docai-validator.pyc'
  Invoke-Checked 'google-cloud-docai-official-validator-compile' $officialGoogleDirectory {
    & $auditPython -c 'import py_compile, sys; py_compile.compile(sys.argv[1], cfile=sys.argv[2], doraise=True)' 'validate_official_packing_list_output.py' $officialGoogleValidatorPyc
  }
  if (-not (Test-Path -LiteralPath $officialGoogleValidatorPyc -PathType Leaf)) { Fail 'Google official packing-list validator bytecode missing' }
  $officialGoogleSamplePath = Join-Path $officialGoogleDirectory 'upstream/google/documentai/snippets/process_document_sample.py'
  $officialGoogleTestPath = Join-Path $officialGoogleDirectory 'upstream/google/documentai/snippets/process_document_sample_test.py'
  $officialGoogleCustomSamplePath = Join-Path $officialGoogleDirectory 'upstream/google/documentai/snippets/handle_response_sample.py'
  $officialGoogleCustomTestPath = Join-Path $officialGoogleDirectory 'upstream/google/documentai/snippets/handle_response_sample_test.py'
  $officialGoogleLicensePath = Join-Path $officialGoogleDirectory 'upstream/google/LICENSE.txt'
  if ((Get-Item -LiteralPath $officialGoogleSamplePath).Length -ne 3596 -or (Get-FileHash -LiteralPath $officialGoogleSamplePath -Algorithm SHA256).Hash.ToLowerInvariant() -cne '7e384dc2c38ebdcbc2c268a63bb9be9794576c1408034e84d53fd8ac5dc504ee') { Fail 'Google official process_document sample identity mismatch' }
  if ((Get-Item -LiteralPath $officialGoogleTestPath).Length -ne 1697 -or (Get-FileHash -LiteralPath $officialGoogleTestPath -Algorithm SHA256).Hash.ToLowerInvariant() -cne '5b39bbd2d309efcc4791bd739fa38961639fa33d7a4e58b739050d3cfd16fcd7') { Fail 'Google official process_document upstream test identity mismatch' }
  if ((Get-Item -LiteralPath $officialGoogleCustomSamplePath).Length -ne 20533 -or (Get-FileHash -LiteralPath $officialGoogleCustomSamplePath -Algorithm SHA256).Hash.ToLowerInvariant() -cne '8fc1adaceba3ac8ad4d871f8003b2ce43fbadfd9cc7b80599cefcc6edc8bfb91') { Fail 'Google official Custom Document Extractor sample identity mismatch' }
  if ((Get-Item -LiteralPath $officialGoogleCustomTestPath).Length -ne 7859 -or (Get-FileHash -LiteralPath $officialGoogleCustomTestPath -Algorithm SHA256).Hash.ToLowerInvariant() -cne 'a60d4d9a53aa3c6383d60bb2cfc22489fde622d09ac994d68161ee771b1fea45') { Fail 'Google official Custom Document Extractor upstream test identity mismatch' }
  if ((Get-Item -LiteralPath $officialGoogleLicensePath).Length -ne 11357 -or (Get-FileHash -LiteralPath $officialGoogleLicensePath -Algorithm SHA256).Hash.ToLowerInvariant() -cne 'c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4') { Fail 'Google official process_document license identity mismatch' }
  $officialGoogleLock = Get-Content -LiteralPath (Join-Path $officialGoogleDirectory 'fixture-lock.json') -Raw | ConvertFrom-Json -Depth 20
  if (@($officialGoogleLock.sources).Count -ne 2 -or @($officialGoogleLock.fixtures).Count -ne 4 -or $officialGoogleLock.sources[0].commit -cne '841379c28828404c6d3944e3da9a8a0f46b4cd0e' -or $officialGoogleLock.sources[1].commit -cne '001ba391ab4a2f40d001cc0387618cb3c3699523') { Fail 'Google official process_document fixture lock identity mismatch' }
  $script:officialGoogleProcessSampleCount = 1

  $officialGoogleLifecycleRoot = Join-Path $tempRoot 'google-document-ai-official-lifecycle'
  Invoke-Checked 'google-document-ai-official-lifecycle-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/GOOGLE_DOCUMENT_AI_OFFICIAL_LIFECYCLE.md') -Destination $officialGoogleLifecycleRoot
  }
  $officialGoogleLifecycleDirectory = Join-Path $officialGoogleLifecycleRoot 'google_document_ai_official_lifecycle'
  Invoke-Checked 'google-document-ai-official-lifecycle-static-regression' $officialGoogleLifecycleDirectory {
    & $auditPython (Join-Path $officialGoogleLifecycleDirectory 'verify_official_lifecycle.py') --static-only
  }
  $script:officialGoogleLifecycleCount = 1

  $officialDaprRoot = Join-Path $tempRoot 'dapr-official-transactional-outbox'
  Invoke-Checked 'dapr-official-outbox-pack-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/DAPR_OFFICIAL_TRANSACTIONAL_OUTBOX.md') -Destination $officialDaprRoot
  }
  $officialDaprDirectory = Join-Path $officialDaprRoot 'dapr_official_transactional_outbox'
  Invoke-Checked 'dapr-official-outbox-source-regression' $officialDaprDirectory {
    & pwsh -NoProfile -File (Join-Path $officialDaprDirectory 'test_embedded_sources.ps1')
  }
  Invoke-Checked 'dapr-official-outbox-acquisition-regression' $officialDaprDirectory {
    & pwsh -NoProfile -File (Join-Path $officialDaprDirectory 'test_acquisition.ps1')
  }
  $officialDaprLock = Get-Content -LiteralPath (Join-Path $officialDaprDirectory 'official-source-lock.json') -Raw | ConvertFrom-Json -Depth 30
  if ($officialDaprLock.runtime.tag -cne 'v1.18.3' -or $officialDaprLock.runtime.commit -cne '2dcf2e3548faf3c740b3d29ee300e2132db65cff' -or $officialDaprLock.runtime.archiveSha256 -cne 'dd64436b01c06bdf59db6c7c949dce81eff36c293581980af93e12a80bfd6681') { Fail 'Dapr Runtime official outbox lock identity mismatch' }
  if ($officialDaprLock.runtime.focusedReachableVulnerabilities -ne 0 -or $officialDaprLock.rejectedGoSdk.admission -cne 'REJECTED_UNTIL_OFFICIAL_FIX_AND_REVALIDATION') { Fail 'Dapr Runtime scan or Go SDK rejection condition mismatch' }
  $officialDaprRuntimePath = Join-Path $officialDaprDirectory 'upstream/runtime/pkg/runtime/pubsub/outbox.go'
  $officialDaprDocsPath = Join-Path $officialDaprDirectory 'upstream/docs/howto-outbox.md'
  if ((Get-FileHash -LiteralPath $officialDaprRuntimePath -Algorithm SHA256).Hash.ToLowerInvariant() -cne '8a0a060481e107f3f0617fc1d6baf86c602105d3ef29c429a83abbc0c04923fc') { Fail 'Dapr Runtime official outbox source identity mismatch' }
  if ((Get-FileHash -LiteralPath $officialDaprDocsPath -Algorithm SHA256).Hash.ToLowerInvariant() -cne '60487ccdd57bd7c308c0fc87f5e48e7843fb09dbb723e1532a0d8d412672e385') { Fail 'Dapr Docs official outbox identity mismatch' }
  $script:officialDaprOutboxCount = 1

  $officialPgDurableRoot = Join-Path $tempRoot 'microsoft-pg-durable-human-handoff'
  Invoke-Checked 'microsoft-pg-durable-human-handoff-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/MICROSOFT_PG_DURABLE_HUMAN_HANDOFF.md') -Destination $officialPgDurableRoot
  }
  $officialPgDurableDirectory = Join-Path $officialPgDurableRoot 'microsoft_pg_durable_human_handoff/upstream'
  $officialPgDurableFiles = @(Get-ChildItem -LiteralPath $officialPgDurableDirectory -Recurse -File)
  if ($officialPgDurableFiles.Count -ne 26) { Fail "Microsoft pg_durable handoff expected 26 exact files, got $($officialPgDurableFiles.Count)" }
  $officialPgDurableIdentities = @{
    'LICENSE.txt' = 'aab465651982b811f0dcd1f6403eb8e89282b75d3bf6ac685ada7e67023a1009'
    'examples/invoice-approval/scripts/smoke_check.sh' = 'cbf8a30bc7bf25983efcaae16431d4819232c4880be9b7824a2824bddf43048e'
    'examples/invoice-approval/sql/05_start_workflow.sql' = '11bf11a2d8ecb6e22dbf24be31ffb636b72407232c010bfbd04cd06452f2505d'
    'examples/invoice-approval/sql/07_approve.sql' = 'eefb437da3d82c84f01b571bec6b222524c36f7ab60d1e62ef55177a0db3fa30'
  }
  foreach ($entry in $officialPgDurableIdentities.GetEnumerator()) {
    $identityPath = Join-Path $officialPgDurableDirectory $entry.Key
    if ((Get-FileHash -LiteralPath $identityPath -Algorithm SHA256).Hash.ToLowerInvariant() -cne $entry.Value) { Fail "Microsoft pg_durable official identity mismatch: $($entry.Key)" }
  }
  if (Test-Path -LiteralPath (Join-Path $officialPgDurableDirectory 'examples/invoice-approval/scripts/live_smoke_check.sh')) { Fail 'Microsoft pg_durable no-final-LF live smoke was incorrectly represented as VERBATIM' }
  $officialPgDurableRequirements = (Get-Content -LiteralPath (Join-Path $officialPgDurableDirectory 'examples/invoice-approval/function-app/requirements.txt') -Raw).Trim()
  if ($officialPgDurableRequirements -cne 'azure-functions>=1.20.0') { Fail 'Microsoft pg_durable upstream open dependency condition drifted' }
  [void](Get-Content -LiteralPath (Join-Path $officialPgDurableDirectory 'examples/invoice-approval/function-app/host.json') -Raw | ConvertFrom-Json)
  [void](Get-Content -LiteralPath (Join-Path $officialPgDurableDirectory 'examples/invoice-approval/function-app/classify_invoice/function.json') -Raw | ConvertFrom-Json)
  $officialPgDurablePyc = Join-Path $tempRoot 'pg-durable-classify-invoice.pyc'
  Invoke-Checked 'microsoft-pg-durable-python-syntax' $officialPgDurableDirectory {
    & $auditPython -c 'import py_compile, sys; py_compile.compile(sys.argv[1], cfile=sys.argv[2], doraise=True)' 'examples/invoice-approval/function-app/classify_invoice/__init__.py' $officialPgDurablePyc
  }
  if (-not (Test-Path -LiteralPath $officialPgDurablePyc -PathType Leaf)) { Fail 'Microsoft pg_durable classifier bytecode missing' }
  $script:officialPgDurableHandoffCount = 1
  Write-Output 'MICROSOFT_PG_DURABLE_LIVE_GATES_SKIPPED reason=Preview_runtime_PostgreSQL_Azure_identity_and_cost_authorization_required'

  $officialAvmSftpRoot = Join-Path $tempRoot 'microsoft-avm-secure-sftp-intake'
  Invoke-Checked 'microsoft-avm-secure-sftp-intake-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/MICROSOFT_AVM_SECURE_SFTP_INTAKE.md') -Destination $officialAvmSftpRoot
  }
  $officialAvmSftpDirectory = Join-Path $officialAvmSftpRoot 'azure_avm_secure_sftp_intake'
  $officialAvmSftpFiles = @(Get-ChildItem -LiteralPath $officialAvmSftpDirectory -Recurse -File)
  if ($officialAvmSftpFiles.Count -ne 8) { Fail "Microsoft AVM SFTP expected 8 files, got $($officialAvmSftpFiles.Count)" }
  Invoke-Checked 'microsoft-avm-secure-sftp-intake-contract' $officialAvmSftpDirectory {
    if ($AllowNetwork) {
      & pwsh -NoProfile -File (Join-Path $officialAvmSftpDirectory 'verify_contract.ps1') -AllowNetwork
    } else {
      & pwsh -NoProfile -File (Join-Path $officialAvmSftpDirectory 'verify_contract.ps1')
    }
  }
  $officialAvmSftpIdentities = @{
    'upstream/LICENSE' = 'c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383'
    'upstream/local-user/main.bicep' = 'efd933ceeb703f6d71b6a8f185ff41f38e0734014f8f35930c1b2afb9c059adc'
    'upstream/waf-aligned/main.test.bicep' = '8f4e0c67787dfbb829824121cb813415cafd85d9ed2f15740c2b5ad9de4eab24'
  }
  foreach ($entry in $officialAvmSftpIdentities.GetEnumerator()) {
    if ((Get-FileHash -LiteralPath (Join-Path $officialAvmSftpDirectory $entry.Key) -Algorithm SHA256).Hash.ToLowerInvariant() -cne $entry.Value) { Fail "Microsoft AVM SFTP official identity mismatch: $($entry.Key)" }
  }
  $script:officialAvmSftpIntakeCount = 1
  Write-Output 'MICROSOFT_AVM_SFTP_LIVE_GATES_SKIPPED reason=Azure_account_network_cost_key_transfer_and_retained_store_required'

  $secureEmailMimeRoot = Join-Path $tempRoot 'secure-email-mime-quarantine-core'
  Invoke-Checked 'secure-email-mime-core-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/SECURE_EMAIL_MIME_QUARANTINE_CORE.md') -Destination $secureEmailMimeRoot
  }
  $secureEmailMimeDirectory = Join-Path $secureEmailMimeRoot 'secure_email_mime_quarantine_core'
  if (@(Get-ChildItem -LiteralPath $secureEmailMimeDirectory -File).Count -ne 5) { Fail 'secure email MIME core expected 5 files' }
  Invoke-Checked 'secure-email-mime-core-contract' $secureEmailMimeDirectory {
    & pwsh -NoProfile -File (Join-Path $secureEmailMimeDirectory 'verify_contract.ps1') -PythonExecutable $auditPython
  }
  $secureEmailProfile = Get-Content -Raw -LiteralPath (Join-Path $secureEmailMimeDirectory 'provider-profile.template.json') | ConvertFrom-Json
  if ($secureEmailProfile.automatic_storage_authorized -ne $false -or @($secureEmailProfile.allowed_content_types).Count -ne 0) { Fail 'secure email MIME distributed profile must enable nothing' }
  $script:secureEmailMimeCoreCount = 1
  Write-Output 'SECURE_EMAIL_MIME_LIVE_GATES_SKIPPED reason=provider_DNS_KMS_retained_raw_queue_and_target_security_required'

  $awsSesReceiverRoot = Join-Path $tempRoot 'aws-ses-immutable-email-receiver'
  Invoke-Checked 'aws-ses-immutable-email-receiver-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/AWS_SES_IMMUTABLE_EMAIL_RECEIVER.md') -Destination $awsSesReceiverRoot
  }
  $awsSesReceiverDirectory = Join-Path $awsSesReceiverRoot 'aws_ses_immutable_email_receiver'
  if (@(Get-ChildItem -LiteralPath $awsSesReceiverDirectory -File).Count -ne 9) { Fail 'AWS SES immutable email receiver expected 9 files' }
  $awsSesComposedRoot = Join-Path $tempRoot 'aws-ses-email-composed-contract'
  [void](New-Item -ItemType Directory -Path $awsSesComposedRoot)
  Copy-Item -LiteralPath $awsSesReceiverDirectory -Destination $awsSesComposedRoot -Recurse
  Copy-Item -LiteralPath $secureEmailMimeDirectory -Destination $awsSesComposedRoot -Recurse
  Invoke-Checked 'aws-ses-immutable-email-receiver-contract' $awsSesComposedRoot {
    & pwsh -NoProfile -File (Join-Path $awsSesComposedRoot 'aws_ses_immutable_email_receiver/verify_contract.ps1') -ComposedRoot $awsSesComposedRoot -PythonExecutable $auditPython
  }
  $awsSesReceiverProfile = Get-Content -Raw -LiteralPath (Join-Path $awsSesReceiverDirectory 'provider-profile.template.json') | ConvertFrom-Json
  if ($awsSesReceiverProfile.approve_live_effects -ne $false -or $awsSesReceiverProfile.dns_mx_proven -ne $false -or $awsSesReceiverProfile.ses_identity_proven -ne $false -or @($awsSesReceiverProfile.allowed_content_types).Count -ne 0) { Fail 'AWS SES receiver distributed profile must enable nothing' }
  $script:awsSesImmutableEmailReceiverCount = 1
  Write-Output 'AWS_SES_RECEIVER_LIVE_GATES_SKIPPED reason=AWS_account_SES_identity_MX_cost_approval_canary_DLQ_replay_and_target_security_required'

  $awsGuardDutyReleaseRoot = Join-Path $tempRoot 'aws-guardduty-immutable-release-gate'
  Invoke-Checked 'aws-guardduty-immutable-release-gate-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/AWS_GUARDDUTY_IMMUTABLE_RELEASE_GATE.md') -Destination $awsGuardDutyReleaseRoot
  }
  $awsGuardDutyReleaseDirectory = Join-Path $awsGuardDutyReleaseRoot 'aws_guardduty_immutable_release_gate'
  if (@(Get-ChildItem -LiteralPath $awsGuardDutyReleaseDirectory -File).Count -ne 9) { Fail 'AWS GuardDuty immutable release gate expected 9 files' }
  Invoke-Checked 'aws-guardduty-immutable-release-gate-contract' $awsGuardDutyReleaseDirectory {
    & pwsh -NoProfile -File (Join-Path $awsGuardDutyReleaseDirectory 'verify_contract.ps1') -PythonExecutable $auditPython
  }
  $awsGuardDutyReleaseProfile = Get-Content -Raw -LiteralPath (Join-Path $awsGuardDutyReleaseDirectory 'provider-profile.template.json') | ConvertFrom-Json
  if ($awsGuardDutyReleaseProfile.approve_live_effects -ne $false -or $awsGuardDutyReleaseProfile.approve_guardduty_costs -ne $false) { Fail 'AWS GuardDuty release distributed profile must enable nothing' }
  $script:awsGuardDutyImmutableReleaseGateCount = 1
  Write-Output 'AWS_GUARDDUTY_RELEASE_LIVE_GATES_SKIPPED reason=AWS_account_GuardDuty_cost_KMS_isolation_SNS_canaries_DLQ_redrive_and_rollback_required'

  $awsGuardDutyIdpDispatchRoot = Join-Path $tempRoot 'aws-guardduty-magika-idp-dispatch-gate'
  Invoke-Checked 'aws-guardduty-magika-idp-dispatch-gate-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/AWS_GUARDDUTY_MAGIKA_IDP_DISPATCH_GATE.md') -Destination $awsGuardDutyIdpDispatchRoot
  }
  $awsGuardDutyIdpDispatchDirectory = Join-Path $awsGuardDutyIdpDispatchRoot 'aws_guardduty_magika_idp_dispatch'
  if (@(Get-ChildItem -LiteralPath $awsGuardDutyIdpDispatchDirectory -File).Count -ne 11) { Fail 'AWS GuardDuty Magika IDP dispatch gate expected 11 files' }
  Invoke-Checked 'aws-guardduty-magika-idp-dispatch-gate-contract' $awsGuardDutyIdpDispatchDirectory {
    & pwsh -NoProfile -File (Join-Path $awsGuardDutyIdpDispatchDirectory 'verify_contract.ps1') -PythonExecutable $auditPython
  }
  $awsGuardDutyIdpDispatchProfile = Get-Content -Raw -LiteralPath (Join-Path $awsGuardDutyIdpDispatchDirectory 'provider-profile.template.json') | ConvertFrom-Json -Depth 20
  if (
    $awsGuardDutyIdpDispatchProfile.approve_live_effects -ne $false -or
    $awsGuardDutyIdpDispatchProfile.approve_provider_costs -ne $false -or
    $awsGuardDutyIdpDispatchProfile.aws_idp_source_acquired_and_verified -ne $false -or
    $awsGuardDutyIdpDispatchProfile.idp_config_corpus_evaluated -ne $false -or
    @($awsGuardDutyIdpDispatchProfile.allowed_magika_types).Count -ne 0
  ) { Fail 'AWS GuardDuty Magika IDP dispatch distributed profile must enable nothing' }
  $script:awsGuardDutyMagikaIdpDispatchGateCount = 1
  Write-Output 'AWS_GUARDDUTY_MAGIKA_IDP_DISPATCH_LIVE_GATES_SKIPPED reason=AWS_IDP_stack_config_corpus_cost_canaries_and_downstream_idempotency_required'

  $awsIdpHandoffRoot = Join-Path $tempRoot 'aws-idp-immutable-evaluation-handoff'
  Invoke-Checked 'aws-idp-immutable-evaluation-handoff-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/AWS_IDP_IMMUTABLE_EVALUATION_HANDOFF.md') -Destination $awsIdpHandoffRoot
  }
  $awsIdpHandoffDirectory = Join-Path $awsIdpHandoffRoot 'aws_idp_immutable_evaluation_handoff'
  if (@(Get-ChildItem -LiteralPath $awsIdpHandoffDirectory -File).Count -ne 8) { Fail 'AWS IDP immutable evaluation handoff expected 8 files' }
  Invoke-Checked 'aws-idp-immutable-evaluation-handoff-contract' $awsIdpHandoffDirectory {
    & pwsh -NoProfile -File (Join-Path $awsIdpHandoffDirectory 'verify_contract.ps1') -ComponentRoot $awsIdpHandoffDirectory
  }
  $awsIdpHandoffProfile = Get-Content -Raw -LiteralPath (Join-Path $awsIdpHandoffDirectory 'provider-profile.template.json') | ConvertFrom-Json -Depth 20
  if ($awsIdpHandoffProfile.acknowledgeConditionedState -ne $false -or @($awsIdpHandoffProfile.allowedClasses).Count -ne 0 -or $awsIdpHandoffProfile.retainDays -ne 0 -or $awsIdpHandoffProfile.maxCanonicalBytes -ne 0) { Fail 'AWS IDP handoff distributed profile must enable nothing' }
  $script:awsIdpImmutableEvaluationHandoffCount = 1
  Write-Output 'AWS_IDP_IMMUTABLE_EVALUATION_HANDOFF_LIVE_GATES_SKIPPED reason=AWS_IDP_stack_config_HITL_corpus_KMS_ObjectLock_DLQ_restore_cost_and_downstream_evaluator_required'

  $awsIdpDecisionRoot = Join-Path $tempRoot 'aws-idp-evaluation-decision-worker'
  Invoke-Checked 'aws-idp-evaluation-decision-worker-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/AWS_IDP_EVALUATION_DECISION_WORKER.md') -Destination $awsIdpDecisionRoot
  }
  $awsIdpDecisionDirectory = Join-Path $awsIdpDecisionRoot 'aws_idp_evaluation_decision_worker'
  if (@(Get-ChildItem -LiteralPath $awsIdpDecisionDirectory -File).Count -ne 8) { Fail 'AWS IDP evaluation decision worker expected 8 files' }
  Invoke-Checked 'aws-idp-evaluation-decision-worker-contract' $awsIdpDecisionDirectory {
    & pwsh -NoProfile -File (Join-Path $awsIdpDecisionDirectory 'verify_contract.ps1') -ComponentRoot $awsIdpDecisionDirectory
  }
  $awsIdpDecisionProfile = Get-Content -Raw -LiteralPath (Join-Path $awsIdpDecisionDirectory 'provider-profile.template.json') | ConvertFrom-Json -Depth 40
  if ($awsIdpDecisionProfile.acknowledgeConditionedState -ne $false -or $awsIdpDecisionProfile.retainDays -ne 0 -or $awsIdpDecisionProfile.maxObjectBytes -ne 0 -or $awsIdpDecisionProfile.idempotencySeconds -ne 0) { Fail 'AWS IDP decision worker distributed profile must enable nothing' }
  $script:awsIdpEvaluationDecisionWorkerCount = 1
  Write-Output 'AWS_IDP_EVALUATION_DECISION_WORKER_LIVE_GATES_SKIPPED reason=AWS_account_policy_corpus_HITL_KMS_ObjectLock_DLQ_load_restore_cost_and_transactional_persistence_required'

  $awsIdpPersistenceRoot = Join-Path $tempRoot 'aws-idp-postgres-persistence-boundary'
  Invoke-Checked 'aws-idp-postgres-persistence-boundary-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/AWS_IDP_POSTGRES_PERSISTENCE_BOUNDARY.md') -Destination $awsIdpPersistenceRoot
  }
  $awsIdpPersistenceDirectory = Join-Path $awsIdpPersistenceRoot 'aws_idp_postgres_persistence_boundary'
  if (@(Get-ChildItem -LiteralPath $awsIdpPersistenceDirectory -File).Count -ne 11) { Fail 'AWS IDP PostgreSQL persistence boundary expected 11 files' }
  Invoke-Checked 'aws-idp-postgres-persistence-boundary-contract' $awsIdpPersistenceDirectory {
    & pwsh -NoProfile -File (Join-Path $awsIdpPersistenceDirectory 'verify_contract.ps1') -ComponentRoot $awsIdpPersistenceDirectory
  }
  $awsIdpPersistenceProfile = Get-Content -Raw -LiteralPath (Join-Path $awsIdpPersistenceDirectory 'project-profile.template.json') | ConvertFrom-Json -Depth 40
  $awsIdpPersistencePolicy = Get-Content -Raw -LiteralPath (Join-Path $awsIdpPersistenceDirectory 'persistence-policy.template.json') | ConvertFrom-Json -Depth 40
  if ($awsIdpPersistenceProfile.acknowledgeConditionedState -ne $false -or $awsIdpPersistenceProfile.maxObjectBytes -ne 0 -or $awsIdpPersistencePolicy.automaticPersistenceAuthorized -ne $false -or @($awsIdpPersistencePolicy.allowedProfiles.PSObject.Properties).Count -ne 0) { Fail 'AWS IDP persistence templates must enable nothing' }
  $script:awsIdpPostgresPersistenceBoundaryCount = 1
  Write-Output 'AWS_IDP_POSTGRES_PERSISTENCE_BOUNDARY_LIVE_GATES_SKIPPED reason=AWS_Aurora_policy_migration_RLS_duplicates_races_DLQ_outbox_recovery_restore_rollback_cost_and_domain_mapping_required'

  $debeziumPostgresOutboxRoot = Join-Path $tempRoot 'debezium-postgres-outbox-runtime'
  Invoke-Checked 'debezium-postgres-outbox-runtime-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/DEBEZIUM_POSTGRES_OUTBOX_RUNTIME.md') -Destination $debeziumPostgresOutboxRoot
  }
  $debeziumPostgresOutboxDirectory = Join-Path $debeziumPostgresOutboxRoot 'debezium_postgres_outbox_runtime'
  if (@(Get-ChildItem -LiteralPath $debeziumPostgresOutboxDirectory -File).Count -ne 10) { Fail 'Debezium PostgreSQL outbox runtime expected 10 files' }
  Invoke-Checked 'debezium-postgres-outbox-runtime-contract' $debeziumPostgresOutboxDirectory {
    & pwsh -NoProfile -File (Join-Path $debeziumPostgresOutboxDirectory 'verify_contract.ps1') -ComponentRoot $debeziumPostgresOutboxDirectory
  }
  $debeziumPostgresOutboxProfile = Get-Content -Raw -LiteralPath (Join-Path $debeziumPostgresOutboxDirectory 'project-profile.template.json') | ConvertFrom-Json -Depth 40
  $debeziumPostgresOutboxLock = Get-Content -Raw -LiteralPath (Join-Path $debeziumPostgresOutboxDirectory 'runtime-lock.json') | ConvertFrom-Json -Depth 40
  if ($debeziumPostgresOutboxProfile.acknowledgeConditionedState -ne $false -or $debeziumPostgresOutboxProfile.heartbeatIntervalMs -ne 0 -or $debeziumPostgresOutboxLock.image.indexDigest -ne 'sha256:76db18f20116557b2e550d5844d2bbebc927cd9c215ec2271c4ab713cc74f386') { Fail 'Debezium PostgreSQL outbox templates or lock drifted' }
  $script:debeziumPostgresOutboxRuntimeCount = 1
  Write-Output 'DEBEZIUM_POSTGRES_OUTBOX_RUNTIME_LIVE_GATES_SKIPPED reason=OCI_pull_PostgreSQL_logical_publication_CDC_role_Kafka_TLS_inbox_duplicates_restart_order_load_recovery_rollback_observability_and_cost_required'

  $debeziumPostgresInboxRoot = Join-Path $tempRoot 'debezium-postgres-inbox-consumer'
  Invoke-Checked 'debezium-postgres-inbox-consumer-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/DEBEZIUM_POSTGRES_INBOX_CONSUMER.md') -Destination $debeziumPostgresInboxRoot
  }
  $debeziumPostgresInboxDirectory = Join-Path $debeziumPostgresInboxRoot 'debezium_postgres_inbox_consumer'
  if (@(Get-ChildItem -LiteralPath $debeziumPostgresInboxDirectory -Recurse -File).Count -ne 24) { Fail 'Debezium PostgreSQL inbox consumer expected 24 files' }
  Invoke-Checked 'debezium-postgres-inbox-consumer-contract' $debeziumPostgresInboxDirectory {
    & pwsh -NoProfile -File (Join-Path $debeziumPostgresInboxDirectory 'verify_contract.ps1')
  }
  $debeziumPostgresInboxProfile = Get-Content -Raw -LiteralPath (Join-Path $debeziumPostgresInboxDirectory 'project-profile.template.json') | ConvertFrom-Json -Depth 40
  $debeziumPostgresInboxLock = Get-Content -Raw -LiteralPath (Join-Path $debeziumPostgresInboxDirectory 'source-lock.json') | ConvertFrom-Json -Depth 40
  if ($debeziumPostgresInboxProfile.acknowledgeConditionedState -ne $false -or $debeziumPostgresInboxProfile.mappingDomainType -ne '' -or $debeziumPostgresInboxProfile.mappingExpression -ne '' -or $debeziumPostgresInboxProfile.mappingOutputSchemaSha256 -ne '' -or $debeziumPostgresInboxLock.adaptationRuntime.quarkusBom -ne '3.39.1' -or $debeziumPostgresInboxLock.adaptationRuntime.jacksonDatabindOverride -ne '2.22.1' -or $debeziumPostgresInboxLock.adaptationRuntime.celJava -ne 'dev.cel:cel:0.13.0' -or $debeziumPostgresInboxLock.adaptationRuntime.components -ne 170 -or $debeziumPostgresInboxLock.adaptationRuntime.osvFindings -ne 0) { Fail 'Debezium PostgreSQL inbox profile/runtime lock drifted' }
  $readinessMappingJson = (& $readinessPython -c 'import json,sys; sys.path.insert(0,sys.argv[1]); import validate_project_readiness as v; print(json.dumps(sorted(v.MAPPING_FIELDS)))' $readinessDirectory | Out-String).Trim()
  if ($LASTEXITCODE -ne 0 -or -not $readinessMappingJson) { Fail 'Readiness mapping contract could not be loaded' }
  $readinessMappingFields = @($readinessMappingJson | ConvertFrom-Json)
  $consumerMappingFields = @($debeziumPostgresInboxProfile.PSObject.Properties.Name | Where-Object { $_ -match '^mapping' } | Sort-Object)
  $mappingContractDrift = @(Compare-Object ($readinessMappingFields | Sort-Object) $consumerMappingFields)
  if ($mappingContractDrift.Count -ne 0) { Fail "Readiness/consumer mapping fields drifted: $($mappingContractDrift | Out-String)" }
  $rendererConsumerJson = (& $readinessPython -c 'import json,sys; sys.path.insert(0,sys.argv[1]); import render_consumer_profile as r; print(json.dumps(sorted(r.CONSUMER_FIELDS)))' $readinessDirectory | Out-String).Trim()
  if ($LASTEXITCODE -ne 0 -or -not $rendererConsumerJson) { Fail 'Readiness renderer consumer schema could not be loaded' }
  $rendererConsumerFields = @($rendererConsumerJson | ConvertFrom-Json)
  $consumerProfileFields = @($debeziumPostgresInboxProfile.PSObject.Properties.Name | Sort-Object)
  $consumerSchemaDrift = @(Compare-Object ($rendererConsumerFields | Sort-Object) $consumerProfileFields)
  if ($consumerSchemaDrift.Count -ne 0) { Fail "Readiness renderer/consumer profile schema drifted: $($consumerSchemaDrift | Out-String)" }
  $script:debeziumPostgresInboxConsumerCount = 1
  Write-Output 'DEBEZIUM_POSTGRES_INBOX_CONSUMER_LIVE_GATES_SKIPPED reason=Kafka_PostgreSQL_TLS_auth_schema_registry_mapping_schema_corpus_approval_poison_restart_redelivery_order_load_backup_restore_observability_required'

  $googleCelMappingRoot = Join-Path $tempRoot 'google-cel-document-mapping'
  Invoke-Checked 'google-cel-document-mapping-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/GOOGLE_CEL_DOCUMENT_MAPPING.md') -Destination $googleCelMappingRoot
  }
  $googleCelMappingDirectory = Join-Path $googleCelMappingRoot 'google_cel_document_mapping'
  if (@(Get-ChildItem -LiteralPath $googleCelMappingDirectory -Recurse -File).Count -ne 7) { Fail 'Google CEL document mapping expected 7 files' }
  $googleCelMappingProfile = Get-Content -Raw -LiteralPath (Join-Path $googleCelMappingDirectory 'project-profile.template.json') | ConvertFrom-Json -Depth 40
  $googleCelMappingLock = Get-Content -Raw -LiteralPath (Join-Path $googleCelMappingDirectory 'source-lock.json') | ConvertFrom-Json -Depth 40
  $googleCelMappingSource = Get-Content -Raw -LiteralPath (Join-Path $googleCelMappingDirectory 'src/main/java/com/elite/mapping/CelDocumentMapping.java')
  if ($googleCelMappingProfile.acknowledgeConditionedState -ne $false -or $googleCelMappingProfile.expression -ne '' -or @($googleCelMappingProfile.exactOutputKeys).Count -ne 0) { Fail 'Google CEL mapping template must fail closed' }
  if ($googleCelMappingLock.maven.coordinate -ne 'dev.cel:cel:0.13.0' -or $googleCelMappingLock.datedEvidence.upstreamFailures -ne 5 -or $googleCelMappingLock.datedEvidence.knownVulnerabilityFindings -ne 0) { Fail 'Google CEL mapping source lock drifted' }
  foreach ($marker in @('.setStandardMacros()','maxExpressionCodePointSize(4096)','maxParseRecursionDepth(32)','mapping expression must be single-line','mapping expression hash mismatch','mapping result keys do not match approved schema','document contains a non-JSON or imprecise numeric value')) {
    if (-not $googleCelMappingSource.Contains($marker)) { Fail "Google CEL mapping contract missing $marker" }
  }
  $script:googleCelDocumentMappingCount = 1
  Write-Output 'GOOGLE_CEL_DOCUMENT_MAPPING_TARGET_GATES_SKIPPED reason=project_schemas_mapping_corpus_approvals_determinism_load_rollback_and_transactional_join_required'

  $businessCentralPackRoot = Join-Path $tempRoot 'official-business-central-pack'
  Invoke-Checked 'official-business-central-pack-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/MICROSOFT_BUSINESS_CENTRAL_EXECUTABLE_PLATFORM.md') -Destination $businessCentralPackRoot
  }
  Invoke-Checked 'official-business-central-runtime-suite' $businessCentralPackRoot {
    & pwsh -NoProfile -File (Join-Path $businessCentralPackRoot 'business_central/test_business_central_runtime.ps1')
  }
  $businessCentralLock = Get-Content -LiteralPath (Join-Path $businessCentralPackRoot 'business_central/official-bc-platform-lock.json') -Raw | ConvertFrom-Json -Depth 20
  if ($businessCentralLock.bccontainerhelper.version -ne '6.1.16' -or $businessCentralLock.business_central_artifact.version -ne '28.4.53241.0') { Fail 'official Business Central platform lock identity mismatch' }
  $script:businessCentralArtifactCount = 1

  $markItDownRoot = Join-Path $tempRoot 'microsoft-markitdown-local-runtime'
  Invoke-Checked 'microsoft-markitdown-local-pack-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/MICROSOFT_MARKITDOWN_LOCAL_RUNTIME.md') -Destination $markItDownRoot
  }
  $markItDownDirectory = Join-Path $markItDownRoot 'markitdown_local_runtime'
  Test-PinnedHashPythonRequirements (Join-Path $markItDownDirectory 'requirements-win-py312.lock') 44
  Invoke-Checked 'microsoft-markitdown-local-profile-contract' $markItDownDirectory {
    $profile = Get-Content -LiteralPath (Join-Path $markItDownDirectory 'conversion-profile.template.json') -Raw | ConvertFrom-Json -Depth 20
    if ($profile.schema -ne 'elite-markitdown-local-conversion-profile/v2' -or $profile.package -ne 'markitdown==0.1.7') { Fail 'MarkItDown profile identity mismatch' }
    if (@($profile.formats | Where-Object decision -eq 'REQUIRED').Count -ne 0) { Fail 'MarkItDown template enabled an unconfigured format' }
    foreach ($control in @('remote_urls','plugins','llm','archive_expansion','automatic_business_storage')) {
      if ($profile.$control -ne $false) { Fail "MarkItDown template enabled $control" }
    }
  }
  if (-not $auditPython) { Fail 'python is required for MarkItDown local wrapper regressions' }
  Invoke-Checked 'microsoft-markitdown-local-wrapper-unit' $markItDownDirectory {
    & $auditPython -m unittest -v test_convert_local_document.py
  }
  if ($AllowNetwork) {
    $pythonDetails = & $auditPython -c "import json, platform, struct, sys; print(json.dumps({'implementation':platform.python_implementation(),'version':list(sys.version_info[:2]),'bits':struct.calcsize('P')*8,'platform':sys.platform}))"
    if ($LASTEXITCODE -ne 0) { Fail 'could not inspect Python runtime for MarkItDown lock' }
    $runtime = $pythonDetails | ConvertFrom-Json
    if (-not $IsWindows -or $runtime.implementation -ne 'CPython' -or $runtime.version[0] -ne 3 -or $runtime.version[1] -ne 12 -or $runtime.bits -ne 64) {
      Fail 'MarkItDown frozen lock requires CPython 3.12 Windows x86-64; generate and admit a separate exact target lock for another runtime'
    }
    $markItDownVenv = Join-Path $tempRoot 'microsoft-markitdown-venv'
    Invoke-Checked 'microsoft-markitdown-create-venv' $tempRoot { & $auditPython -m venv $markItDownVenv }
    $markItDownPython = Join-Path $markItDownVenv 'Scripts/python.exe'
    Invoke-Checked 'microsoft-markitdown-frozen-install' $markItDownDirectory {
      & $markItDownPython -m pip install --disable-pip-version-check --only-binary=:all: --require-hashes -r requirements-win-py312.lock
    }
    Invoke-Checked 'microsoft-markitdown-frozen-environment' $markItDownDirectory {
      & $markItDownPython -m pip check
      if ($LASTEXITCODE -ne 0) { Fail 'MarkItDown pip check failed' }
      & $markItDownPython verify_win_py312_environment.py
      if ($LASTEXITCODE -ne 0) { Fail 'MarkItDown 44-distribution verifier failed' }
      & $markItDownPython -m unittest -v test_convert_local_document.py
      if ($LASTEXITCODE -ne 0) { Fail 'MarkItDown wrapper regressions failed in frozen environment' }
      & $markItDownPython -c "import importlib.metadata; from markitdown import MarkItDown; assert importlib.metadata.version('markitdown') == '0.1.7'; assert callable(MarkItDown(enable_plugins=False).convert_local)"
      if ($LASTEXITCODE -ne 0) { Fail 'MarkItDown official package import/signature probe failed' }
    }
  } else {
    Write-Output 'MICROSOFT_MARKITDOWN_RUNTIME_GATES_SKIPPED reason=AllowNetwork_not_authorized target=CPython_3.12_Windows_x86-64'
  }

  $documentRoutingRoot = Join-Path $tempRoot 'document-pipeline-routing-gate'
  Invoke-Checked 'document-pipeline-routing-pack-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/DOCUMENT_PIPELINE_ROUTING_GATE.md') -Destination $documentRoutingRoot
  }
  $documentRoutingDirectory = Join-Path $documentRoutingRoot 'document_pipeline_routing'
  Invoke-Checked 'document-pipeline-routing-unit' $documentRoutingDirectory {
    & $auditPython -m py_compile validate_document_routing.py
    if ($LASTEXITCODE -ne 0) { Fail 'document routing validator compile failed' }
    & $auditPython -m unittest -v test_validate_document_routing.py
    if ($LASTEXITCODE -ne 0) { Fail 'document routing validator regressions failed' }
  }

  $secureFileRoot = Join-Path $tempRoot 'secure-local-file-ingestion-gate'
  Invoke-Checked 'secure-local-file-ingestion-pack-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/SECURE_LOCAL_FILE_INGESTION_GATE.md') -Destination $secureFileRoot
  }
  $secureFileDirectory = Join-Path $secureFileRoot 'secure_file_gate'
  Invoke-Checked 'secure-local-file-ingestion-unit' $secureFileDirectory {
    & $auditPython -m unittest discover -s . -p 'test_*.py' -v
  }

  $durableDocumentRoot = Join-Path $tempRoot 'microsoft-durable-document-orchestration'
  Invoke-Checked 'microsoft-durable-document-orchestration-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/MICROSOFT_DURABLE_DOCUMENT_ORCHESTRATION.md') -Destination $durableDocumentRoot
  }
  $durableDocumentDirectory = Join-Path $durableDocumentRoot 'durable_document_orchestration'
  Invoke-Checked 'microsoft-durable-document-orchestration-static-contract' $durableDocumentDirectory {
    Test-MultilineHashPythonRequirements (Join-Path $durableDocumentDirectory 'requirements-windows-py312.lock') 6
    $lock = Get-Content -LiteralPath (Join-Path $durableDocumentDirectory 'official-artifact-lock.json') -Raw | ConvertFrom-Json -Depth 20
    if ($lock.runtime_target -ne 'CPython 3.12 / Windows x86-64' -or $lock.runtime_packages -ne 6 -or $lock.osv_findings -ne 0 -or @($lock.sources).Count -ne 2) { Fail 'Durable document artifact lock summary mismatch' }
    $sdk = $lock.sources | Where-Object id -CEQ 'microsoft-durabletask-python-1.9.0'
    $skill = $lock.sources | Where-Object id -CEQ 'azure-samples-durable-task-scheduler-5636c25'
    if ($sdk.commit -cne '6dfdbac521d9d59d31d7ea975fb669d66f86f4c4' -or $sdk.archive_sha256 -cne 'c8c941b25804ae9abcc0470f5594cab4929216a981982616344eb2b84a7916e8' -or $sdk.wheel_sha256 -cne 'a759af4ad8e6897922575886e93bf40c6b3af936eaed83798d492ee561607185' -or $sdk.sdist_sha256 -cne '55460fbfda8941e721096c4639e72111b03638708df1556af466160ee649478d') { Fail 'Durable Task exact source/package identity mismatch' }
    if ($skill.commit -cne '5636c25ffbdaabaca1062c3b6c6e77072f2b2eb2' -or $skill.skill_sha256 -cne '553f0dc278c63bb129df5557a7d6f7d645bf827efba4cb077588688f88a2e9e2' -or $skill.packaged_license_provenance -cne 'ADAPTED_FINAL_NEWLINE_ONLY') { Fail 'Durable Task official skill/provenance identity mismatch' }
    $actualSkillHash = (Get-FileHash -Algorithm SHA256 -LiteralPath (Join-Path $durableDocumentDirectory 'upstream/durable-task-python-SKILL.md')).Hash.ToLowerInvariant()
    $actualLicenseHash = (Get-FileHash -Algorithm SHA256 -LiteralPath (Join-Path $durableDocumentDirectory 'upstream/LICENSE.md')).Hash.ToLowerInvariant()
    if ($actualSkillHash -cne $skill.skill_sha256 -or $actualLicenseHash -cne $skill.packaged_license_sha256) { Fail 'Durable Task packaged upstream bytes drifted' }
    $profile = Get-Content -LiteralPath (Join-Path $durableDocumentDirectory 'pipeline-profile.template.json') -Raw | ConvertFrom-Json -Depth 20
    if ($profile.enabled -ne $false -or $profile.backend -cne 'AWAITING_SELECTION' -or $profile.storage_lane -cne 'AWAITING_SELECTION' -or $profile.automatic_business_persistence -ne $false -or $profile.live_effects_approved -ne $false -or @($profile.document_classes).Count -ne 0) { Fail 'Durable document profile enabled an unconfigured effect' }
    $orchestration = Get-Content -LiteralPath (Join-Path $durableDocumentDirectory 'documentflow/orchestration.py') -Raw
    foreach ($required in @('UNKNOWN_ORIGINAL_STORAGE_EFFECT','RECONCILE_STORAGE','UNKNOWN_BUSINESS_EFFECT','PARTIAL_EFFECT','validate_review_event','documentflow_retain_final')) {
      if (-not $orchestration.Contains($required)) { Fail "Durable document fail-closed contract missing: $required" }
    }
    $runtimeSource = Get-Content -LiteralPath (Join-Path $durableDocumentDirectory 'documentflow/runtime.py') -Raw
    foreach ($required in @('class ManagedTaskHubGrpcWorker','loop.shutdown_asyncgens()','loop.shutdown_default_executor()','loop.close()')) {
      if (-not $runtimeSource.Contains($required)) { Fail "Durable document managed worker cleanup missing: $required" }
    }
  }
  $script:documentOrchestratorCount = 1
  if ($AllowNetwork) {
    $pythonDetails = & $auditPython -c "import json, platform, struct, sys; print(json.dumps({'implementation':platform.python_implementation(),'version':list(sys.version_info[:2]),'bits':struct.calcsize('P')*8,'platform':sys.platform}))"
    if ($LASTEXITCODE -ne 0) { Fail 'could not inspect Python runtime for Durable Task lock' }
    $runtime = $pythonDetails | ConvertFrom-Json
    if (-not $IsWindows -or $runtime.implementation -ne 'CPython' -or $runtime.version[0] -ne 3 -or $runtime.version[1] -ne 12 -or $runtime.bits -ne 64) {
      Fail 'Durable Task frozen lock requires CPython 3.12 Windows x86-64; generate and admit a separate exact target lock for another runtime'
    }
    $durableVenv = Join-Path $tempRoot 'microsoft-durable-document-venv'
    Invoke-Checked 'microsoft-durable-document-create-venv' $tempRoot { & $auditPython -m venv $durableVenv }
    $durablePython = Get-VenvPython $durableVenv
    Invoke-Checked 'microsoft-durable-document-frozen-install' $durableDocumentDirectory {
      & $durablePython -m pip install --disable-pip-version-check --only-binary=:all: --require-hashes -r requirements-windows-py312.lock
    }
    Invoke-Checked 'microsoft-durable-document-pip-check' $durableDocumentDirectory { & $durablePython -m pip check }
    Invoke-Checked 'microsoft-durable-document-unit-and-in-memory-integration' $durableDocumentDirectory {
      & $durablePython -W error::ResourceWarning -m unittest discover -s . -p 'test_*.py' -v
      if ($LASTEXITCODE -ne 0) { Fail 'Durable document orchestration tests failed' }
      & $durablePython -m compileall -q documentflow test_contracts.py test_pipeline.py
    }
  } else {
    Write-Output 'MICROSOFT_DURABLE_DOCUMENT_RUNTIME_GATES_SKIPPED reason=AllowNetwork_not_authorized target=CPython_3.12_Windows_x86-64'
  }

  $awsTextractRoot = Join-Path $tempRoot 'aws-textract-document-runtime'
  Invoke-Checked 'aws-textract-document-runtime-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/GO_AWS_TEXTRACT_DOCUMENT_RUNTIME.md') -Destination $awsTextractRoot
  }
  $awsTextractDirectory = Join-Path $awsTextractRoot 'aws_textract_runtime'
  Invoke-Checked 'aws-textract-static-lock-and-fail-closed-profile' $awsTextractDirectory {
    $goMod = Get-Content -LiteralPath (Join-Path $awsTextractDirectory 'go.mod') -Raw
    if ($goMod -notmatch '(?m)^\s*github\.com/aws/aws-sdk-go-v2/service/textract v1\.45\.0$') { Fail 'AWS Textract exact module pin missing' }
    $profile = Get-Content -LiteralPath (Join-Path $awsTextractDirectory 'document-profile.template.json') -Raw | ConvertFrom-Json -Depth 20
    if ($profile.schema -ne 'elite-aws-textract-document-profile/v2' -or $profile.sdk -ne 'github.com/aws/aws-sdk-go-v2/service/textract@v1.45.0') { Fail 'AWS Textract profile identity mismatch' }
    if ($profile.region -ne 'CONFIGURATION_REQUIRED' -or @($profile.classes | Where-Object decision -eq 'REQUIRED').Count -ne 0) { Fail 'AWS Textract template enabled unconfigured execution' }
    if (@($profile.classes | Where-Object automatic_storage -ne $false).Count -ne 0) { Fail 'AWS Textract template enabled automatic storage' }
  }

  $awsEnterpriseRoot = Join-Path $tempRoot 'aws-enterprise-storage-email-adapters'
  Invoke-Checked 'aws-enterprise-storage-email-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/GO_AWS_ENTERPRISE_STORAGE_EMAIL_ADAPTERS.md') -Destination $awsEnterpriseRoot
  }
  $awsEnterpriseDirectory = Join-Path $awsEnterpriseRoot 'aws_enterprise_adapters'
  Invoke-Checked 'aws-enterprise-object-lock-contract' $awsEnterpriseDirectory {
    $goMod = Get-Content -LiteralPath (Join-Path $awsEnterpriseDirectory 'go.mod') -Raw
    if ($goMod -notmatch '(?m)^\s*github\.com/aws/aws-sdk-go-v2/service/s3 v1\.107\.3$') { Fail 'AWS S3 exact module pin missing' }
    $profile = Get-Content -LiteralPath (Join-Path $awsEnterpriseDirectory 'provider-profile.template.json') -Raw | ConvertFrom-Json -Depth 20
    if ($profile.schema -ne 'elite-aws-enterprise-adapters-profile/v1' -or $profile.object_storage.decision -ne 'BLOCKED_ACCESS_AND_POLICY_REQUIRED' -or $profile.object_storage.object_lock.decision -ne 'BLOCKED_OWNER_MODE_RETENTION_REQUIRED' -or $profile.object_storage.cost_approved -ne $false) { Fail 'AWS enterprise profile enabled storage without authority' }
    $storage = Get-Content -LiteralPath (Join-Path $awsEnterpriseDirectory 'awsenterprise/storage.go') -Raw
    foreach ($required in @('GetObjectLockConfiguration','ObjectLockRetainUntilDate','GetObjectRetention','ExpectedBucketOwner','RetentionVerified')) {
      if (-not $storage.Contains($required)) { Fail "AWS Object Lock contract missing: $required" }
    }
  }
  $auditGo = Resolve-Tool 'go' $GoExecutable
  $auditGoVersion = Get-ToolVersion $auditGo @('version')
  if ($auditGo -and (Test-MinimumVersion $auditGoVersion ([version]'1.26.7'))) {
    Invoke-Checked 'aws-enterprise-storage-email-go-gates' $awsEnterpriseDirectory {
      $previousProxy = $env:GOPROXY
      try {
        $env:GOPROXY = 'off'
        & $auditGo mod verify
        if ($LASTEXITCODE -ne 0) { Fail 'AWS enterprise module verification failed' }
        & $auditGo test -count=1 ./...
        if ($LASTEXITCODE -ne 0) { Fail 'AWS enterprise tests failed' }
        & $auditGo vet ./...
        if ($LASTEXITCODE -ne 0) { Fail 'AWS enterprise vet failed' }
        & $auditGo build ./...
        if ($LASTEXITCODE -ne 0) { Fail 'AWS enterprise build failed' }
      } finally { $env:GOPROXY = $previousProxy }
    }
  } else {
    Write-Output 'AWS_ENTERPRISE_STORAGE_EMAIL_GO_GATES_SKIPPED reason=Go_1.26.7_not_available'
  }

  $awsIntakeRoot = Join-Path $tempRoot 'aws-secure-quarantine-intake'
  Invoke-Checked 'aws-secure-quarantine-intake-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/GO_AWS_SECURE_QUARANTINE_INTAKE.md') -Destination $awsIntakeRoot
  }
  $awsIntakeDirectory = Join-Path $awsIntakeRoot 'aws_secure_quarantine_intake'
  Invoke-Checked 'aws-secure-quarantine-intake-fail-closed-contract' $awsIntakeDirectory {
    $goMod = Get-Content -LiteralPath (Join-Path $awsIntakeDirectory 'go.mod') -Raw
    if ($goMod -notmatch '(?m)^\s*github\.com/aws/aws-sdk-go-v2/service/s3 v1\.107\.3$') { Fail 'AWS intake S3 exact module pin missing' }
    $profile = Get-Content -LiteralPath (Join-Path $awsIntakeDirectory 'provider-profile.template.json') -Raw | ConvertFrom-Json -Depth 20
    if ($profile.schema -ne 'elite-aws-quarantine-intake-profile/v1' -or $profile.enabled -ne $false -or $profile.decision -ne 'BLOCKED_ACCESS_COST_AND_POLICY_REQUIRED' -or $profile.upload.automatic_business_storage_authorized -ne $false) { Fail 'AWS intake profile enabled an unapproved effect' }
    if ($profile.aws.temporary_credentials_verified -ne $false -or $profile.aws.signature_age_policy_verified -ne $false -or $profile.aws.quarantine_boundary_verified -ne $false) { Fail 'AWS intake profile pre-verified target guardrails' }
    $intake = Get-Content -LiteralPath (Join-Path $awsIntakeDirectory 'awsintake/intake.go') -Raw
    foreach ($required in @('PresignPutObject','IfNoneMatch','ChecksumSHA256','ExpectedBucketOwner','ChecksumModeEnabled','type SessionStore interface','s.Store.Create','s.Store.Get','s.Store.Complete')) {
      if (-not $intake.Contains($required)) { Fail "AWS intake contract missing: $required" }
    }
  }
  if ($auditGo -and (Test-MinimumVersion $auditGoVersion ([version]'1.26.7'))) {
    Invoke-Checked 'aws-secure-quarantine-intake-go-gates' $awsIntakeDirectory {
      $previousProxy = $env:GOPROXY
      try {
        $env:GOPROXY = 'off'
        & $auditGo mod verify
        if ($LASTEXITCODE -ne 0) { Fail 'AWS intake module verification failed' }
        & $auditGo test -count=1 ./...
        if ($LASTEXITCODE -ne 0) { Fail 'AWS intake tests failed' }
        & $auditGo vet ./...
        if ($LASTEXITCODE -ne 0) { Fail 'AWS intake vet failed' }
        & $auditGo build ./...
        if ($LASTEXITCODE -ne 0) { Fail 'AWS intake build failed' }
      } finally { $env:GOPROXY = $previousProxy }
    }
  } else {
    Write-Output 'AWS_SECURE_QUARANTINE_INTAKE_GO_GATES_SKIPPED reason=Go_1.26.7_not_available'
  }

  $googleStorageRoot = Join-Path $tempRoot 'google-cloud-storage-adapter'
  Invoke-Checked 'google-cloud-storage-adapter-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/GO_GOOGLE_CLOUD_STORAGE_ADAPTER.md') -Destination $googleStorageRoot
  }
  $googleStorageDirectory = Join-Path $googleStorageRoot 'google_cloud_storage_adapter'
  Invoke-Checked 'google-cloud-storage-fail-closed-contract' $googleStorageDirectory {
    $goMod = Get-Content -LiteralPath (Join-Path $googleStorageDirectory 'go.mod') -Raw
    if ($goMod -notmatch '(?m)^require cloud\.google\.com/go/storage v1\.65\.1$' -or $goMod -notmatch '(?m)^\s*golang\.org/x/text v0\.39\.0 // indirect$') { Fail 'Google Cloud Storage exact/fixed module pins missing' }
    $profile = Get-Content -LiteralPath (Join-Path $googleStorageDirectory 'provider-profile.template.json') -Raw | ConvertFrom-Json -Depth 20
    if ($profile.schema -ne 'elite-google-cloud-storage-profile/v1' -or $profile.object_storage.decision -ne 'BLOCKED_ACCESS_COST_AND_POLICY_REQUIRED' -or $profile.object_storage.automatic_storage_authorized -ne $false -or $profile.object_storage.retention.irreversible_locked_mode_approved -ne $false) { Fail 'Google Cloud Storage profile enabled storage without authority' }
    $lock = Get-Content -LiteralPath (Join-Path $googleStorageDirectory 'official-artifact-lock.json') -Raw | ConvertFrom-Json -Depth 20
    if ($lock.module -ne 'cloud.google.com/go/storage' -or $lock.version -ne 'v1.65.1' -or $lock.commit -ne '3ab7d1390bbbb40cf197a545b941f9df760a3269' -or $lock.commit_signature_verified -ne $true -or $lock.license -ne 'Apache-2.0') { Fail 'Google Cloud Storage official artifact identity mismatch' }
    $backend = Get-Content -LiteralPath (Join-Path $googleStorageDirectory 'gcsstorage/google_backend.go') -Raw
    foreach ($required in @('Conditions{DoesNotExist: true}','SendCRC32C = true','KMSKeyName','ObjectRetention','Generation(generation).Attrs')) {
      if (-not $backend.Contains($required)) { Fail "Google Cloud Storage provider contract missing: $required" }
    }
  }
  if ($auditGo -and (Test-MinimumVersion $auditGoVersion ([version]'1.26.7'))) {
    Invoke-Checked 'google-cloud-storage-go-gates' $googleStorageDirectory {
      $previousProxy = $env:GOPROXY
      try {
        $env:GOPROXY = 'off'
        & $auditGo mod verify
        if ($LASTEXITCODE -ne 0) { Fail 'Google Cloud Storage module verification failed' }
        & $auditGo test -count=1 ./...
        if ($LASTEXITCODE -ne 0) { Fail 'Google Cloud Storage tests failed' }
        & $auditGo vet ./...
        if ($LASTEXITCODE -ne 0) { Fail 'Google Cloud Storage vet failed' }
        & $auditGo build ./...
        if ($LASTEXITCODE -ne 0) { Fail 'Google Cloud Storage build failed' }
      } finally { $env:GOPROXY = $previousProxy }
    }
  } else {
    Write-Output 'GOOGLE_CLOUD_STORAGE_GO_GATES_SKIPPED reason=Go_1.26.7_not_available'
  }

  Invoke-ProviderAdapterAudit
  $azureBlobRoot = Join-Path $tempRoot 'azure-blob-evidence-adapter'
  $azureBlobDirectory = Join-Path $azureBlobRoot 'azure_blob_storage_adapter'
  if ($AllowNetwork) {
    $pythonDetails = & $auditPython -c "import json, platform, struct, sys; print(json.dumps({'implementation':platform.python_implementation(),'version':list(sys.version_info[:2]),'bits':struct.calcsize('P')*8,'platform':sys.platform}))"
    if ($LASTEXITCODE -ne 0) { Fail 'could not inspect Python runtime for Azure Blob adapter lock' }
    $runtime = $pythonDetails | ConvertFrom-Json
    if (-not $IsWindows -or $runtime.implementation -ne 'CPython' -or $runtime.version[0] -ne 3 -or $runtime.version[1] -ne 12 -or $runtime.bits -ne 64) {
      Fail 'Azure Blob adapter frozen lock requires CPython 3.12 Windows x86-64; generate and admit a separate exact target lock for another runtime'
    }
    $azureBlobVenv = Join-Path $tempRoot 'azure-blob-evidence-venv'
    Invoke-Checked 'azure-blob-evidence-create-venv' $tempRoot { & $auditPython -m venv $azureBlobVenv }
    $azureBlobPython = Get-VenvPython $azureBlobVenv
    Invoke-Checked 'azure-blob-evidence-frozen-install' $azureBlobDirectory {
      & $azureBlobPython -m pip install --disable-pip-version-check --only-binary=:all: --require-hashes -r requirements-windows-py312.lock
    }
    Invoke-Checked 'azure-blob-evidence-pip-check' $azureBlobDirectory { & $azureBlobPython -m pip check }
    Invoke-Checked 'azure-blob-evidence-unit-and-sdk-contract' $azureBlobDirectory {
      & $azureBlobPython -m unittest -v
      if ($LASTEXITCODE -ne 0) { Fail 'Azure Blob evidence adapter tests failed' }
      & $azureBlobPython -m compileall -q azureblob test_storage.py
    }
  } else {
    Write-Output 'AZURE_BLOB_EVIDENCE_RUNTIME_GATES_SKIPPED reason=AllowNetwork_not_authorized target=CPython_3.12_Windows_x86-64'
  }

  $azureDocumentRoot = Join-Path $tempRoot 'azure-content-understanding-document-runtime'
  Invoke-Checked 'azure-content-understanding-runtime-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/AZURE_CONTENT_UNDERSTANDING_DOCUMENT_RUNTIME.md') -Destination $azureDocumentRoot
  }
  $azureDocumentDirectory = Join-Path $azureDocumentRoot 'azure_document_runtime'
  Invoke-Checked 'azure-content-understanding-runtime-lock-contract' $azureDocumentDirectory {
    Test-PinnedHashPythonRequirements (Join-Path $azureDocumentDirectory 'requirements-windows-py312.lock') 9
  }
  if ($AllowNetwork) {
    $pythonDetails = & $auditPython -c "import json, platform, struct, sys; print(json.dumps({'implementation':platform.python_implementation(),'version':list(sys.version_info[:2]),'bits':struct.calcsize('P')*8,'platform':sys.platform}))"
    if ($LASTEXITCODE -ne 0) { Fail 'could not inspect Python runtime for Azure Content Understanding lock' }
    $runtime = $pythonDetails | ConvertFrom-Json
    if (-not $IsWindows -or $runtime.implementation -ne 'CPython' -or $runtime.version[0] -ne 3 -or $runtime.version[1] -ne 12 -or $runtime.bits -ne 64) {
      Fail 'Azure Content Understanding frozen lock requires CPython 3.12 Windows x86-64; generate and admit a separate exact target lock for another runtime'
    }
    $azureVenv = Join-Path $tempRoot 'azure-content-understanding-venv'
    Invoke-Checked 'azure-content-understanding-create-venv' $tempRoot { & $auditPython -m venv $azureVenv }
    $azurePython = Get-VenvPython $azureVenv
    Invoke-Checked 'azure-content-understanding-frozen-install' $azureDocumentDirectory {
      & $azurePython -m pip install --disable-pip-version-check --require-hashes -r (Join-Path $azureDocumentDirectory 'requirements-windows-py312.lock')
    }
    Invoke-Checked 'azure-content-understanding-pip-check' $azureDocumentDirectory { & $azurePython -m pip check }
    Invoke-Checked 'azure-content-understanding-wrapper-unit' $azureDocumentDirectory {
      & $azurePython -m unittest -v test_content_understanding_runtime.py
    }
  } else {
    Write-Output 'AZURE_CONTENT_UNDERSTANDING_RUNTIME_GATES_SKIPPED reason=AllowNetwork_not_authorized target=CPython_3.12_Windows_x86-64'
  }

  $googleDocumentRoot = Join-Path $tempRoot 'google-document-ai-runtime'
  Invoke-Checked 'google-document-ai-runtime-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/GOOGLE_DOCUMENT_AI_RUNTIME.md') -Destination $googleDocumentRoot
  }
  $googleDocumentDirectory = Join-Path $googleDocumentRoot 'google_document_runtime'
  Invoke-Checked 'google-document-ai-runtime-lock-contract' $googleDocumentDirectory {
    Test-PinnedHashPythonRequirements (Join-Path $googleDocumentDirectory 'requirements-windows-py312.lock') 23
  }
  if ($AllowNetwork) {
    $pythonDetails = & $auditPython -c "import json, platform, struct, sys; print(json.dumps({'implementation':platform.python_implementation(),'version':list(sys.version_info[:2]),'bits':struct.calcsize('P')*8,'platform':sys.platform}))"
    if ($LASTEXITCODE -ne 0) { Fail 'could not inspect Python runtime for Google Document AI lock' }
    $runtime = $pythonDetails | ConvertFrom-Json
    if (-not $IsWindows -or $runtime.implementation -ne 'CPython' -or $runtime.version[0] -ne 3 -or $runtime.version[1] -ne 12 -or $runtime.bits -ne 64) {
      Fail 'Google Document AI frozen lock requires CPython 3.12 Windows x86-64; generate and admit a separate exact target lock for another runtime'
    }
    $googleDocumentVenv = Join-Path $tempRoot 'google-document-ai-venv'
    Invoke-Checked 'google-document-ai-create-venv' $tempRoot { & $auditPython -m venv $googleDocumentVenv }
    $googleDocumentPython = Get-VenvPython $googleDocumentVenv
    Invoke-Checked 'google-document-ai-frozen-install' $googleDocumentDirectory {
      & $googleDocumentPython -m pip install --disable-pip-version-check --require-hashes -r (Join-Path $googleDocumentDirectory 'requirements-windows-py312.lock')
    }
    Invoke-Checked 'google-document-ai-pip-check' $googleDocumentDirectory { & $googleDocumentPython -m pip check }
    Invoke-Checked 'google-document-ai-wrapper-unit' $googleDocumentDirectory {
      & $googleDocumentPython -m unittest -v test_document_ai_runtime.py
    }
    Invoke-Checked 'google-cloud-docai-official-custom-extractor-regression' $officialGoogleDirectory {
      & $googleDocumentPython test_official_custom_extractor_sample.py
    }
    Invoke-Checked 'google-document-ai-official-lifecycle-sdk-regression' $officialGoogleLifecycleDirectory {
      & $googleDocumentPython (Join-Path $officialGoogleLifecycleDirectory 'verify_official_lifecycle.py')
    }
  } else {
    Write-Output 'GOOGLE_DOCUMENT_AI_RUNTIME_GATES_SKIPPED reason=AllowNetwork_not_authorized target=CPython_3.12_Windows_x86-64'
    Write-Output 'GOOGLE_CUSTOM_EXTRACTOR_SDK_RUNTIME_GATES_SKIPPED reason=AllowNetwork_not_authorized; exact hashes and Python syntax remain verified offline'
  }

  $strictFieldRoot = Join-Path $tempRoot 'strict-document-field-evaluation-gate'
  Invoke-Checked 'strict-document-field-evaluation-pack-materialization' $libraryRoot {
    & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/STRICT_DOCUMENT_FIELD_EVALUATION_GATE.md') -Destination $strictFieldRoot
  }
  Invoke-Checked 'strict-document-field-evaluation-python-syntax' $strictFieldRoot {
    & $auditPython -m py_compile (Join-Path $strictFieldRoot 'strict_field_eval/strict_field_evaluation.py') (Join-Path $strictFieldRoot 'strict_field_eval/test_strict_field_evaluation.py')
  }
}

function Get-ToolchainReport {
  $python = Resolve-PythonTool
  $pwsh = Resolve-Tool 'pwsh'
  $go = Resolve-Tool 'go' $GoExecutable
  $dotnet = Resolve-Tool 'dotnet' $DotNetExecutable
  $node = Resolve-Tool 'node' $NodeExecutable
  $pnpm = Resolve-Tool 'pnpm' $PnpmExecutable
  $docker = Resolve-Tool 'docker' $DockerExecutable
  $psql = Resolve-Tool 'psql' $PsqlExecutable
  $pwshVersion = Get-ToolVersion $pwsh @('--version')
  $pythonVersion = Get-ToolVersion $python @('--version')
  $goVersion = Get-ToolVersion $go @('version')
  $dotnetVersion = Get-ToolVersion $dotnet @('--version')
  $nodeVersion = Get-ToolVersion $node @('--version')
  # V401: version discovery cannot establish admission of the pnpm bundle.
  # The pinned published release includes GHSA-vwc7-r8mq-g2x9. The isolated
  # ZIP-disabled candidate has scoped evidence, not general runtime admission.
  # Do not execute a PATH-provided replacement to infer that admission.
  $pnpmVersion = $null
  $dockerVersion = Get-ToolVersion $docker @('--version')
  $psqlVersion = Get-ToolVersion $psql @('--version')
  @(
    [ordered]@{ id='pwsh-7'; required_for='library'; path=$pwsh; version=$pwshVersion; available=([bool]$pwsh -and (Test-MinimumVersion $pwshVersion ([version]'7.0'))) },
    [ordered]@{ id='python-3.12+'; required_for='quality, licenses, document intelligence'; path=$python; version=$pythonVersion; available=([bool]$python -and (Test-MinimumVersion $pythonVersion ([version]'3.12'))) },
    [ordered]@{ id='go-1.26.7'; required_for='portable backend foundation'; path=$go; version=$goVersion; available=([bool]$go -and (Test-MinimumVersion $goVersion ([version]'1.26.7'))) },
    [ordered]@{ id='dotnet-10+'; required_for='WCF/SOAP client generation lane'; path=$dotnet; version=$dotnetVersion; available=([bool]$dotnet -and (Test-MinimumVersion $dotnetVersion ([version]'10.0'))) },
    [ordered]@{ id='node-24+'; required_for='web/BFF foundation'; path=$node; version=$nodeVersion; available=([bool]$node -and (Test-MinimumVersion $nodeVersion ([version]'24.0'))) },
    [ordered]@{ id='pnpm-11.25.0'; required_for='web/BFF frozen install'; path=$pnpm; version=$pnpmVersion; available=$false; admission='BLOCKED'; reason='V401: published pnpm affected by GHSA-vwc7-r8mq-g2x9; ZIP-disabled candidate only has scoped evidence. Exact runtime admission required.' },
    [ordered]@{ id='docker'; required_for='PostgreSQL integration/recovery and optional Business Central Windows-container gates'; path=$docker; version=$dockerVersion; available=[bool]$dockerVersion },
    [ordered]@{ id='psql'; required_for='direct PostgreSQL migration/restore gates'; path=$psql; version=$psqlVersion; available=[bool]$psqlVersion }
  )
}

function Write-Evidence([string] $Status, [object[]] $Toolchains) {
  if (-not $EvidencePath) { return }
  $target = [IO.Path]::GetFullPath($EvidencePath)
  if (Test-Path -LiteralPath $target) { Fail "EvidencePath already exists: $target" }
  $parent = Split-Path -Parent $target
  if (-not (Test-Path -LiteralPath $parent -PathType Container)) { Fail "EvidencePath parent missing: $parent" }
  $record = [ordered]@{
    schema='elite-executable-library-verification/v1'
    status=$Status
    mode=$Mode
    library_audit_executed=($Mode -ne 'Toolchain')
    availability_is_admission=$false
    allow_network=[bool]$AllowNetwork
    verified_at=[DateTimeOffset]::UtcNow.ToString('O')
    toolchains=$Toolchains
    steps=@($steps)
    prometheus_rule_artifact=$prometheusRuleArtifact
    non_claims=@('no provider sandbox was contacted','no Business Central container was created','no production readiness was inferred','conditional upstreams remain conditional')
  }
  [IO.File]::WriteAllText($target, (($record | ConvertTo-Json -Depth 8) + "`n"), $utf8)
}

try {
  # Availability-only diagnostic: no materialization, audit, install or provider calls.
  # Keep its status distinct from the integrated Preflight/Audit/Foundation gates.
  if ($Mode -eq 'Toolchain') {
    $toolchains = @(Get-ToolchainReport)
    $missing = @($toolchains | Where-Object { -not $_.available })
    $status = if ($missing.Count -eq 0) { 'TOOLS_AVAILABLE' } else { 'TOOLS_MISSING' }
    Write-Evidence $status $toolchains
    foreach ($tool in $toolchains) {
      Write-Output ("TOOLCHAIN id={0} available={1} version={2}" -f $tool.id, ([string]$tool.available).ToLowerInvariant(), $tool.version)
    }
    Write-Output "TOOLCHAIN_PROBE_$status missing=$($missing.Count) library_audit_executed=false availability_is_admission=false"
    if ($RequireAll -and $missing.Count -gt 0) { Fail "missing toolchains: $(@($missing.id) -join ', ')" }
    return
  }
  [void](New-Item -ItemType Directory -Path $tempRoot)
  Invoke-LibraryAudit
  $toolchains = @(Get-ToolchainReport)

  if ($Mode -eq 'Audit') {
    Write-Evidence 'PASS' $toolchains
    Write-Output "VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit packs=$implementationPackCount upstream_sources=$upstreamSourceCount execution_state_controls=$executionStateControlCount capability_gap_resolution_gates=$capabilityGapResolutionGateCount microsoft_devskim_adapted_sast_gates=$microsoftDevSkimAdaptedSastGateCount gitlab_opengrep_signed_sast_gates=$gitlabOpengrepSignedSastGateCount portable_signed_release_evidence_gates=$portableSignedReleaseEvidenceGateCount microsoft_kiota_openapi_client_gates=$microsoftKiotaOpenApiClientGateCount go_native_fuzz_gates=$goNativeFuzzGateCount microsoft_playwright_browser_gates=$microsoftPlaywrightBrowserGateCount google_lighthouse_web_quality_gates=$googleLighthouseWebQualityGateCount aws_powertools_batch_components=$awsPowertoolsBatchComponentCount aws_durable_execution_components=$awsDurableExecutionComponentCount aws_textractor_official_components=$awsTextractorOfficialComponentCount secure_email_mime_cores=$secureEmailMimeCoreCount aws_ses_immutable_email_receivers=$awsSesImmutableEmailReceiverCount aws_guardduty_immutable_release_gates=$awsGuardDutyImmutableReleaseGateCount aws_guardduty_magika_idp_dispatch_gates=$awsGuardDutyMagikaIdpDispatchGateCount aws_idp_immutable_evaluation_handoffs=$awsIdpImmutableEvaluationHandoffCount aws_idp_evaluation_decision_workers=$awsIdpEvaluationDecisionWorkerCount aws_idp_postgres_persistence_boundaries=$awsIdpPostgresPersistenceBoundaryCount debezium_postgres_outbox_runtimes=$debeziumPostgresOutboxRuntimeCount debezium_postgres_inbox_consumers=$debeziumPostgresInboxConsumerCount google_cel_document_mappings=$googleCelDocumentMappingCount document_sdk_artifacts=$documentArtifactCount official_invoice_samples=$officialInvoiceSampleCount official_azure_cu_samples=$officialAzureCuSampleCount official_azure_cu_binary_samples=$officialAzureCuBinarySampleCount official_azure_cu_copy_samples=$officialAzureCuCopySampleCount official_google_process_samples=$officialGoogleProcessSampleCount official_google_lifecycles=$officialGoogleLifecycleCount official_dapr_outboxes=$officialDaprOutboxCount official_pg_durable_handoffs=$officialPgDurableHandoffCount official_avm_sftp_intakes=$officialAvmSftpIntakeCount business_central_artifacts=$businessCentralArtifactCount provider_adapters=$providerAdapterCount document_orchestrators=$documentOrchestratorCount evidence_logs=$evidenceLogCount"
    return
  }

  foreach ($tool in $toolchains) {
    Write-Output ("TOOLCHAIN id={0} available={1} required_for={2}" -f $tool.id, ([string]$tool.available).ToLowerInvariant(), $tool.required_for)
  }
  $missing = @($toolchains | Where-Object { -not $_.available })
  if ($Mode -eq 'Preflight') {
    $status = if ($missing.Count -eq 0) { 'PASS' } else { 'BLOCKED' }
    Write-Evidence $status $toolchains
    Write-Output "EXECUTABLE_PREFLIGHT_$status missing=$($missing.Count)"
    Write-Output 'PROJECT_ACCESS_REMAINS_REQUIRED provider sandboxes, identities/scopes, secrets manager references, document corpus/ground truth, deployment and recovery targets'
    if ($RequireAll -and $missing.Count -gt 0) { Fail "missing toolchains: $(@($missing.id) -join ', ')" }
    return
  }

  $foundationRequired = @('python-3.12+','go-1.26.7','node-24+','pnpm-11.25.0')
  $foundationMissing = @($toolchains | Where-Object { $_.id -in $foundationRequired -and -not $_.available })
  if ($foundationMissing.Count -gt 0) { Fail "foundation toolchains missing: $(@($foundationMissing.id) -join ', ')" }

  $python = ($toolchains | Where-Object id -eq 'python-3.12+').path
  $go = ($toolchains | Where-Object id -eq 'go-1.26.7').path
  $pnpm = ($toolchains | Where-Object id -eq 'pnpm-11.25.0').path
  $node = ($toolchains | Where-Object id -eq 'node-24+').path
  $compositorRoot = Join-Path $tempRoot 'compositor'
  & pwsh -NoProfile -File (Join-Path $libraryRoot 'materialize_markdown_pack.ps1') -PackFile (Join-Path $libraryRoot 'implementation_packs/MARKDOWN_COMPOSITOR_CORE.md') -Destination $compositorRoot | Out-Null
  if ($LASTEXITCODE -ne 0) { Fail "compositor materialization exited $LASTEXITCODE" }
  $compositor = Join-Path $compositorRoot 'tools/compose-markdown-project.ps1'
  $backend = Join-Path $tempRoot 'backend'
  $web = Join-Path $tempRoot 'web'
  & pwsh -NoProfile -File $compositor -PlanFile (Join-Path $libraryRoot 'markdown_system/ENTERPRISE_BACKEND_PACK_PLAN.md') -LibraryRoot $libraryRoot -Destination $backend | Out-Null
  if ($LASTEXITCODE -ne 0) { Fail "backend composition exited $LASTEXITCODE" }
  & pwsh -NoProfile -File $compositor -PlanFile (Join-Path $libraryRoot 'markdown_system/ENTERPRISE_WEB_PACK_PLAN.md') -LibraryRoot $libraryRoot -Destination $web | Out-Null
  if ($LASTEXITCODE -ne 0) { Fail "web composition exited $LASTEXITCODE" }

  Invoke-Checked 'backend-go-unit' $backend { & $go test ./... }
  Invoke-Checked 'backend-ci-runner-unit' (Join-Path $backend 'ci') { & $python test_run_quality_gates.py }
  Invoke-Checked 'backend-packaging-unit' (Join-Path $backend 'packaging') { & $python test_validate_packaging.py }
  Invoke-Checked 'backend-operational-readiness-unit' (Join-Path $backend 'ops/readiness') { & $python test_validate_operational_readiness.py }
  Invoke-Checked 'backend-production-admission-unit' (Join-Path $backend 'production_admission_gate') { & $python -B -m unittest discover -s . -p 'test_*.py' -v }
  Invoke-Checked 'backend-license-gate-unit' (Join-Path $backend 'supply_chain') { & $python test_license_gate.py }
  Invoke-Checked 'web-license-gate-unit' (Join-Path $web 'supply_chain') { & $python test_license_gate.py }
  Invoke-OfflineProviderAdapterRuntimeGates $python
  Invoke-GoProviderAdapterRuntimeGates $go
  Invoke-NetworkProviderAdapterGates $python

  Push-Location $web
  try {
    & $pnpm install --frozen-lockfile --offline
    if ($LASTEXITCODE -ne 0) {
      if (-not $AllowNetwork) { Fail 'web frozen offline install failed; rerun with -AllowNetwork only after network use is authorized' }
      & $pnpm install --frozen-lockfile
      if ($LASTEXITCODE -ne 0) { Fail "web frozen install exited $LASTEXITCODE" }
    }
    $steps.Add([ordered]@{ id='web-frozen-install'; status='PASS' })
    & $pnpm verify
    if ($LASTEXITCODE -ne 0) { Fail "web verify exited $LASTEXITCODE" }
    $steps.Add([ordered]@{ id='web-typecheck-test-build'; status='PASS' })
  } finally { Pop-Location }

  $composedBrowserGate = Join-Path $web 'microsoft_playwright_browser_gate'
  Invoke-Checked 'web-microsoft-playwright-frozen-offline-install' $composedBrowserGate {
    & $pnpm install --ignore-workspace --frozen-lockfile --offline
  }
  Invoke-Checked 'web-microsoft-playwright-runtime-and-target' $composedBrowserGate {
    & pwsh -NoProfile -File (Join-Path $composedBrowserGate 'verify_contract.ps1') -RunEnterpriseWeb -WebRoot $web
  }

  $composedLighthouseGate = Join-Path $web 'google_lighthouse_web_quality_gate'
  Invoke-Checked 'web-google-lighthouse-frozen-offline-install' $composedLighthouseGate {
    & $pnpm install --ignore-workspace --frozen-lockfile --offline
  }
  Invoke-Checked 'web-google-lighthouse-browser-resolution' $composedBrowserGate {
    $script:foundationChromePath = (& $node --input-type=module -e "import { chromium } from '@playwright/test'; process.stdout.write(chromium.executablePath());" | Out-String).Trim()
    if ($LASTEXITCODE -ne 0 -or -not (Test-Path -LiteralPath $script:foundationChromePath -PathType Leaf)) { Fail 'pinned Playwright Chromium executable could not be resolved for Lighthouse' }
  }
  $listener = [Net.Sockets.TcpListener]::new([Net.IPAddress]::Loopback, 0)
  $listener.Start()
  $lighthousePort = ([Net.IPEndPoint]$listener.LocalEndpoint).Port
  $listener.Stop()
  $serverOut = Join-Path $tempRoot 'lighthouse-web.stdout.log'
  $serverErr = Join-Path $tempRoot 'lighthouse-web.stderr.log'
  $start = @{
    FilePath = $pnpm
    ArgumentList = @('exec','next','start','-H','127.0.0.1','-p',[string]$lighthousePort)
    WorkingDirectory = $web
    PassThru = $true
    RedirectStandardOutput = $serverOut
    RedirectStandardError = $serverErr
  }
  if ($IsWindows) { $start.WindowStyle = 'Hidden' }
  $webServer = Start-Process @start
  try {
    $ready = $false
    $targetUrl = "http://127.0.0.1:$lighthousePort/"
    for ($attempt = 0; $attempt -lt 120; $attempt++) {
      if ($webServer.HasExited) { break }
      try {
        $response = Invoke-WebRequest -UseBasicParsing -Uri $targetUrl -TimeoutSec 2
        if ($response.StatusCode -eq 200) { $ready = $true; break }
      } catch {}
      Start-Sleep -Milliseconds 250
    }
    if (-not $ready) {
      $stdout = if (Test-Path -LiteralPath $serverOut) { [IO.File]::ReadAllText($serverOut) } else { '' }
      $stderr = if (Test-Path -LiteralPath $serverErr) { [IO.File]::ReadAllText($serverErr) } else { '' }
      Fail "enterprise web did not start for Lighthouse; stdout=$stdout stderr=$stderr"
    }
    Invoke-Checked 'web-google-lighthouse-five-run-target' $composedLighthouseGate {
      & pwsh -NoProfile -File (Join-Path $composedLighthouseGate 'verify_contract.ps1') -RunTarget -TargetUrl $targetUrl -ChromePath $script:foundationChromePath -OutputDirectory (Join-Path $tempRoot 'lighthouse-results')
    }
  } finally {
    Stop-OwnedProcessTree $webServer
  }

  Write-Evidence 'PASS' $toolchains
  Write-Output 'VERIFY_EXECUTABLE_LIBRARY_PASS mode=Foundation backend=go-test web=typecheck-test-build browser=playwright-chromium-firefox-webkit quality=lighthouse-five-run'
  Write-Output 'PROJECT_CONDITIONED provider, PostgreSQL integration/recovery, authenticated-role journeys, AT accessibility, security/load/deploy and business acceptance gates require selected real environments'
} finally {
  if (Test-Path -LiteralPath $tempRoot) {
    $resolved = [IO.Path]::GetFullPath($tempRoot)
    $systemTemp = [IO.Path]::GetFullPath([IO.Path]::GetTempPath())
    if ($resolved.StartsWith($systemTemp, [StringComparison]::OrdinalIgnoreCase) -and (Split-Path -Leaf $resolved) -like 'elite-executable-verification-*') {
      Remove-Item -LiteralPath $resolved -Recurse -Force
    }
  }
}
