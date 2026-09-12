# Microsoft-generated ARCA WSFE SOAP Adapter

## 1. Metadata

```yaml
pack_id: "MICROSOFT-ARCA-WSFE-SOAP-ADAPTER"
pack_version: "0.3.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Mapea el owner fiscal y `CbtesAsoc` exactos y adquiere siete clases de parámetros desde operaciones WSFEv1 generadas, con respuestas deterministas seguras, consulta 602 fail-closed y cero secretos retornados."
stacks: [".NET SDK 10.0.400", "dotnet-svcutil 8.0.0", "WCF 10.0.652802", "ARCA WSFEv1 4.6"]
compatible_with: ["MICROSOFT-ARCA-WSFE-GENERATED-CLIENT 0.2.x", "MICROSOFT-ARCA-WSAA-CREDENTIAL-CORE 0.1.x", "GO-ARCA-FISCAL-ISSUANCE-API 0.6.x"]
incompatible_with: []
license_expression: "LicenseRef-Workspace-Owner AND ARCA public specification AND Microsoft dependency licenses"
upstream_sources: ["https://www.arca.gob.ar/ws/documentacion/manuales/manual-desarrollador-ARCA-COMPG.pdf", "https://learn.microsoft.com/dotnet/core/additional-tools/dotnet-svcutil-guide", "https://github.com/dotnet/wcf", "https://dotnet.microsoft.com/download/dotnet/10.0"]
verified_at: "2026-08-31"
```

Los catorce bloques son `AUTHORED`: importan el proxy generado y el credential core sin incorporar WSDL, proxy, token, sign, certificado, clave privada, mensajes del proveedor o tablas fiscales inventadas. No se atribuyen a ARCA ni Microsoft.

## 2. Applicability

Use después de generar los clientes WSAA/WSFE de homologación y materializar el credential core y el owner fiscal 0.3.x. Rechácelo si falta condición IVA receptor, detalle IVA/tributos, autenticación válida o si se intenta usar producción. Este pack cierra el mapping SOAP .NET y la adquisición exacta de parámetros; no inventa valores, no define política fiscal ni autoriza una llamada live.

## 3. Architecture contract

El owner PostgreSQL conserva workflow, número, idempotencia y reconciliación. El adapter obtiene un ticket cacheado, construye `FEAuthRequest`, mapea un comprobante por request y llama las operaciones generadas exactas. Los importes nacen de unidades menores y las líneas ya verificadas; `CondicionIVAReceptorIdSpecified=true` evita omisión silenciosa. Error común 602 sin `ResultGet` es la única consulta negativa admitida; otros errores, identidad divergente o respuestas incompletas fallan cerrados. Las siete consultas de parámetros conservan códigos, descripciones, vigencias y atributos provistos, rechazan duplicados/fechas inválidas/clase divergente y producen hash determinista. Hashes y códigos seguros excluyen mensajes, Token y Sign.

## 4. Exact file manifest

```text
CREATE arca/fiscal/bridge/README.md
CREATE arca/fiscal/bridge/src/Elite.Arca.Wsfe.Bridge/Elite.Arca.Wsfe.Bridge.csproj
CREATE arca/fiscal/bridge/src/Elite.Arca.Wsfe.Bridge/packages.lock.json
CREATE arca/fiscal/bridge/src/Elite.Arca.Wsfe.Bridge/Contracts.cs
CREATE arca/fiscal/bridge/src/Elite.Arca.Wsfe.Bridge/GeneratedAdapters.cs
CREATE arca/fiscal/bridge/src/Elite.Arca.Wsfe.Bridge/WsfeBridge.cs
CREATE arca/fiscal/bridge/src/Elite.Arca.Wsfe.Bridge/WsfeParameterBridge.cs
CREATE arca/fiscal/bridge/tests/Elite.Arca.Wsfe.Bridge.Tests/Elite.Arca.Wsfe.Bridge.Tests.csproj
CREATE arca/fiscal/bridge/tests/Elite.Arca.Wsfe.Bridge.Tests/packages.lock.json
CREATE arca/fiscal/bridge/tests/Elite.Arca.Wsfe.Bridge.Tests/Program.cs
CREATE arca/fiscal/bridge/tests/Elite.Arca.Wsfe.Parameter.Tests/Elite.Arca.Wsfe.Parameter.Tests.csproj
CREATE arca/fiscal/bridge/tests/Elite.Arca.Wsfe.Parameter.Tests/packages.lock.json
CREATE arca/fiscal/bridge/tests/Elite.Arca.Wsfe.Parameter.Tests/Program.cs
CREATE tools/test-arca-wsfe-soap-adapter.ps1
```

## 5. Materialization blocks

### FILE: `arca/fiscal/bridge/README.md`
```yaml
block_id: "MS-ARCA-WSFE-ADAPTER:readme:v1"
operation: CREATE
path: "arca/fiscal/bridge/README.md"
sha256: "e2590f12a1086cb457e1574bda6de2e42695742a9d9b85f9244c829773bbbf3b"
provenance: AUTHORED
source: "local adapter governed by ARCA manual 4.6"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````text
# ARCA WSFEv1 generated SOAP adapter

This adapter converts the durable fiscal owner's exact integer/date model into the current ARCA-generated `FEAuthRequest`, `FECAERequest`, `FECompConsultaReq` and their safe responses. It also queries the generated WSFEv1 parameter operations for voucher types, concepts, document types, VAT rates, other taxes, points of sale and recipient VAT conditions. It imports the generated proxy and the credential core; it does not embed WSDL, generated proxy, token, sign, certificate, private key, provider messages or a locally invented fiscal table.

Run `tools/test-arca-wsfe-soap-adapter.ps1` with the official .NET 10.0.400 executable, a homologation output from `MICROSOFT-ARCA-WSFE-GENERATED-CLIENT`, the materialized credential core and an empty work directory. Twelve invoice contracts prove the manual's invoice example mapping, mandatory recipient VAT condition, line-level IVA/Tributos, error 602 consultation outcome, rejection of other errors/identity mismatches and zero secret-bearing return fields. Fifteen parameter contracts prove all seven query classes, deterministic response hashing, exact provider-code preservation, rejection of provider errors/duplicates/class mismatches/invalid dates and zero credential or provider-message leakage.

This is a transport and parameter-acquisition library, not a production deployment or a fiscal policy engine. Parameter values are acquired from ARCA at runtime and remain conditioned on successful authentication and live homologation. Persistence, refresh policy, responsible tax/accounting approval, real credentials, service association and production admission remain mandatory.
````

### FILE: `arca/fiscal/bridge/src/Elite.Arca.Wsfe.Bridge/Elite.Arca.Wsfe.Bridge.csproj`
```yaml
block_id: "MS-ARCA-WSFE-ADAPTER:project:v1"
operation: CREATE
path: "arca/fiscal/bridge/src/Elite.Arca.Wsfe.Bridge/Elite.Arca.Wsfe.Bridge.csproj"
sha256: "66fece54d613f5c5666112ca467ddf84a2d88f465458833ff07e49a62b0637cc"
provenance: AUTHORED
source: "local project governed by Microsoft dotnet-svcutil documentation"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````xml
<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <TargetFramework>net10.0</TargetFramework>
    <ImplicitUsings>enable</ImplicitUsings>
    <Nullable>enable</Nullable>
    <TreatWarningsAsErrors>true</TreatWarningsAsErrors>
    <RestorePackagesWithLockFile>true</RestorePackagesWithLockFile>
    <NuGetAudit>true</NuGetAudit>
    <NuGetAuditMode>all</NuGetAuditMode>
  </PropertyGroup>
  <ItemGroup>
    <ProjectReference Include="../../../generated/Elite.Arca.Wsfe.Generated.csproj" />
    <ProjectReference Include="../../../../credentials/src/Elite.Arca.Credentials/Elite.Arca.Credentials.csproj" />
  </ItemGroup>
</Project>
````

### FILE: `arca/fiscal/bridge/src/Elite.Arca.Wsfe.Bridge/packages.lock.json`
```yaml
block_id: "MS-ARCA-WSFE-ADAPTER:source-lock:v1"
operation: CREATE
path: "arca/fiscal/bridge/src/Elite.Arca.Wsfe.Bridge/packages.lock.json"
sha256: "d15ee5465734d7ef84f113dacba6d014b546ac772f16dbfcf8b566ebd282b202"
provenance: AUTHORED
source: "lock generated by Microsoft .NET SDK 10.0.400"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````json
{
  "version": 1,
  "dependencies": {
    "net10.0": {
      "Microsoft.Extensions.ObjectPool": {
        "type": "Transitive",
        "resolved": "10.0.0",
        "contentHash": "bpeCq0IYmVLACyEUMzFIOQX+zZUElG1t+nu1lSxthe7B+1oNYking7b91305+jNB6iwojp9fqTY9O+Nh7ULQxg=="
      },
      "System.Security.Cryptography.Pkcs": {
        "type": "Transitive",
        "resolved": "10.0.11",
        "contentHash": "8IV+rI3xN/Mkq9MsSX7VZTv9T9Wt+tyhHLwBX9VNZBk8m8kknEs1JeTSjI2x9M2L/fSja+c2WSE5n5Yy0GXdoQ=="
      },
      "System.Security.Cryptography.Xml": {
        "type": "Transitive",
        "resolved": "10.0.11",
        "contentHash": "TokfVsaU2fmcwvsU6ihLCkseVfwhLX6Tmbl3qh3nOBIqHSX4yqFllpslA+PuHERH/ENu3kwqo4deZiyT0Wvh3Q==",
        "dependencies": {
          "System.Security.Cryptography.Pkcs": "10.0.11"
        }
      },
      "System.ServiceModel.Http": {
        "type": "Transitive",
        "resolved": "10.0.652802",
        "contentHash": "G02XZvmccf42QCU5MjviBIg69MSMAVHwL1inVPsNSpfp5g+t5BkQM3DyvWRLN4qmeFDWSF/mA1rIYONIDu/6Dg==",
        "dependencies": {
          "System.ServiceModel.Primitives": "10.0.652802"
        }
      },
      "System.ServiceModel.Primitives": {
        "type": "Transitive",
        "resolved": "10.0.652802",
        "contentHash": "ULfGNl75BNXkpF42wNV2CDXJ64dUZZEa8xO2mBsc4tqbW9QjruxjEB6bAr4Z/T1rNU+leOztIjCJQYsBGFWYlw==",
        "dependencies": {
          "Microsoft.Extensions.ObjectPool": "10.0.0",
          "System.Security.Cryptography.Xml": "10.0.0"
        }
      },
      "elite.arca.credentials": {
        "type": "Project",
        "dependencies": {
          "System.Security.Cryptography.Pkcs": "[10.0.11, )"
        }
      },
      "elite.arca.wsfe.generated": {
        "type": "Project",
        "dependencies": {
          "System.Security.Cryptography.Pkcs": "[10.0.11, )",
          "System.Security.Cryptography.Xml": "[10.0.11, )",
          "System.ServiceModel.Http": "[10.0.652802, )"
        }
      }
    }
  }
}
````

### FILE: `arca/fiscal/bridge/src/Elite.Arca.Wsfe.Bridge/Contracts.cs`
```yaml
block_id: "MS-ARCA-WSFE-ADAPTER:contracts:v1"
operation: CREATE
path: "arca/fiscal/bridge/src/Elite.Arca.Wsfe.Bridge/Contracts.cs"
sha256: "8ae5d5f1d09f113b7903f10ec23d2f761ec188b10e135a5906a4a242e1b70496"
provenance: AUTHORED
source: "local contracts governed by ARCA manual 4.6"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````csharp
using Elite.Arca.Credentials;
using Generated = Elite.Arca.Wsfe.Generated;

namespace Elite.Arca.Wsfe.Bridge;

public sealed record VatLine(int Id, long BaseMinorUnits, long AmountMinorUnits);
public sealed record OtherTaxLine(short Id, string Description, long BaseMinorUnits, int RateBasisPoints, long AmountMinorUnits);
public sealed record AssociatedVoucher(string TaxpayerCuit, int VoucherType, int PointOfSale, long Number, DateOnly IssuedOn);

public sealed class FiscalInvoice
{
    public required string TaxpayerCuit { get; init; }
    public required int PointOfSale { get; init; }
    public required int VoucherType { get; init; }
    public required int Concept { get; init; }
    public required int RecipientDocumentType { get; init; }
    public required string RecipientDocument { get; init; }
    public required int RecipientVatConditionId { get; init; }
    public required string Currency { get; init; }
    public required long VoucherNumber { get; init; }
    public required long TotalMinorUnits { get; init; }
    public required long NetMinorUnits { get; init; }
    public required long VatMinorUnits { get; init; }
    public required long ExemptMinorUnits { get; init; }
    public required long NonTaxedMinorUnits { get; init; }
    public required long OtherTaxMinorUnits { get; init; }
    public required DateOnly IssuedOn { get; init; }
    public DateOnly? ServiceFrom { get; init; }
    public DateOnly? ServiceUntil { get; init; }
    public DateOnly? PaymentDueOn { get; init; }
    public IReadOnlyList<VatLine> VatLines { get; init; } = [];
    public IReadOnlyList<OtherTaxLine> OtherTaxLines { get; init; } = [];
    public IReadOnlyList<AssociatedVoucher> AssociatedVouchers { get; init; } = [];
}

public sealed record SafeWsfeResult(bool Found, bool Authorized, string? Cae, DateOnly? CaeExpiresOn, string ResponseHash, IReadOnlyList<string> ProviderCodes);
public sealed record LastAuthorizedResult(long Number, string ResponseHash, IReadOnlyList<string> ProviderCodes);
public sealed record SafeParameterItem(string Code, string? Description, DateOnly? ValidFrom, DateOnly? ValidUntil, string? VoucherClass, string? EmissionType, string? Blocked, DateOnly? DeregisteredOn);
public sealed record SafeParameterResult(string Kind, string? VoucherClass, IReadOnlyList<SafeParameterItem> Items, string ResponseHash, IReadOnlyList<string> ProviderCodes);

public interface IWsaaAccessSource
{
    Task<LoginTicketAccess> GetAsync(CancellationToken cancellationToken);
}

public interface IWsfeSoap
{
    Task<Generated.FERecuperaLastCbteResponse> LastAuthorizedAsync(Generated.FEAuthRequest auth, int pointOfSale, int voucherType, CancellationToken cancellationToken);
    Task<Generated.FECompConsultaResponse> ConsultAsync(Generated.FEAuthRequest auth, Generated.FECompConsultaReq request, CancellationToken cancellationToken);
    Task<Generated.FECAEResponse> AuthorizeAsync(Generated.FEAuthRequest auth, Generated.FECAERequest request, CancellationToken cancellationToken);
}

public interface IWsfeParameterSoap
{
    Task<Generated.CbteTipoResponse> VoucherTypesAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken);
    Task<Generated.ConceptoTipoResponse> ConceptsAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken);
    Task<Generated.DocTipoResponse> DocumentTypesAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken);
    Task<Generated.IvaTipoResponse> VatRatesAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken);
    Task<Generated.FETributoResponse> OtherTaxesAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken);
    Task<Generated.FEPtoVentaResponse> PointsOfSaleAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken);
    Task<Generated.CondicionIvaReceptorResponse> RecipientVatConditionsAsync(Generated.FEAuthRequest auth, string voucherClass, CancellationToken cancellationToken);
}

public sealed class WsfeResponseException(string message, string responseHash, IReadOnlyList<string> providerCodes) : IOException(message)
{
    public string ResponseHash { get; } = responseHash;
    public IReadOnlyList<string> ProviderCodes { get; } = providerCodes;
}
````

### FILE: `arca/fiscal/bridge/src/Elite.Arca.Wsfe.Bridge/GeneratedAdapters.cs`
```yaml
block_id: "MS-ARCA-WSFE-ADAPTER:generated:v1"
operation: CREATE
path: "arca/fiscal/bridge/src/Elite.Arca.Wsfe.Bridge/GeneratedAdapters.cs"
sha256: "b165ff87ea974037001a9e23871da53d0fc1c0f36698eacc8e2465a89b908cb2"
provenance: AUTHORED
source: "local adapter over Microsoft-generated interfaces"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````csharp
using Elite.Arca.Credentials;
using Generated = Elite.Arca.Wsfe.Generated;

namespace Elite.Arca.Wsfe.Bridge;

public sealed class CredentialAccessSource(WsaaCredentialProvider provider) : IWsaaAccessSource
{
    public Task<LoginTicketAccess> GetAsync(CancellationToken cancellationToken) => provider.GetAsync(cancellationToken);
}

public sealed class GeneratedWsfeSoap(Generated.ServiceSoap client) : IWsfeSoap
{
    public async Task<Generated.FERecuperaLastCbteResponse> LastAuthorizedAsync(Generated.FEAuthRequest auth, int pointOfSale, int voucherType, CancellationToken cancellationToken) =>
        await client.FECompUltimoAutorizadoAsync(auth, pointOfSale, voucherType).WaitAsync(cancellationToken).ConfigureAwait(false);

    public async Task<Generated.FECompConsultaResponse> ConsultAsync(Generated.FEAuthRequest auth, Generated.FECompConsultaReq request, CancellationToken cancellationToken) =>
        await client.FECompConsultarAsync(auth, request).WaitAsync(cancellationToken).ConfigureAwait(false);

    public async Task<Generated.FECAEResponse> AuthorizeAsync(Generated.FEAuthRequest auth, Generated.FECAERequest request, CancellationToken cancellationToken) =>
        await client.FECAESolicitarAsync(auth, request).WaitAsync(cancellationToken).ConfigureAwait(false);
}

public sealed class GeneratedWsfeParameterSoap(Generated.ServiceSoap client) : IWsfeParameterSoap
{
    public async Task<Generated.CbteTipoResponse> VoucherTypesAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken) =>
        await client.FEParamGetTiposCbteAsync(auth).WaitAsync(cancellationToken).ConfigureAwait(false);
    public async Task<Generated.ConceptoTipoResponse> ConceptsAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken) =>
        await client.FEParamGetTiposConceptoAsync(auth).WaitAsync(cancellationToken).ConfigureAwait(false);
    public async Task<Generated.DocTipoResponse> DocumentTypesAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken) =>
        await client.FEParamGetTiposDocAsync(auth).WaitAsync(cancellationToken).ConfigureAwait(false);
    public async Task<Generated.IvaTipoResponse> VatRatesAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken) =>
        await client.FEParamGetTiposIvaAsync(auth).WaitAsync(cancellationToken).ConfigureAwait(false);
    public async Task<Generated.FETributoResponse> OtherTaxesAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken) =>
        await client.FEParamGetTiposTributosAsync(auth).WaitAsync(cancellationToken).ConfigureAwait(false);
    public async Task<Generated.FEPtoVentaResponse> PointsOfSaleAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken) =>
        await client.FEParamGetPtosVentaAsync(auth).WaitAsync(cancellationToken).ConfigureAwait(false);
    public async Task<Generated.CondicionIvaReceptorResponse> RecipientVatConditionsAsync(Generated.FEAuthRequest auth, string voucherClass, CancellationToken cancellationToken) =>
        await client.FEParamGetCondicionIvaReceptorAsync(auth, voucherClass).WaitAsync(cancellationToken).ConfigureAwait(false);
}
````

### FILE: `arca/fiscal/bridge/src/Elite.Arca.Wsfe.Bridge/WsfeBridge.cs`
```yaml
block_id: "MS-ARCA-WSFE-ADAPTER:bridge:v1"
operation: CREATE
path: "arca/fiscal/bridge/src/Elite.Arca.Wsfe.Bridge/WsfeBridge.cs"
sha256: "948bd4164e673e72bb73dff9144ec307ef258613f130cc85cc0c125046de38a8"
provenance: AUTHORED
source: "local exact mapping governed by ARCA manual 4.6"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````csharp
using System.Globalization;
using System.Security.Cryptography;
using System.Text.Json;
using Generated = Elite.Arca.Wsfe.Generated;

namespace Elite.Arca.Wsfe.Bridge;

public sealed class WsfeBridge(IWsaaAccessSource credentials, IWsfeSoap soap)
{
    public async Task<LastAuthorizedResult> LastAuthorizedAsync(FiscalInvoice invoice, CancellationToken cancellationToken = default)
    {
        Validate(invoice, requireNumber: false);
        var response = await soap.LastAuthorizedAsync(await AuthAsync(invoice, cancellationToken), invoice.PointOfSale, invoice.VoucherType, cancellationToken).ConfigureAwait(false);
        var codes = Codes(response.Errors, response.Events, null);
        var hash = Hash(new { operation = "FECompUltimoAutorizado", response.PtoVta, response.CbteTipo, response.CbteNro, codes });
        if (response.Errors is { Length: > 0 } || response.PtoVta != invoice.PointOfSale || response.CbteTipo != invoice.VoucherType || response.CbteNro < 0)
            throw new WsfeResponseException("ARCA last-authorized response is invalid.", hash, codes);
        return new LastAuthorizedResult(response.CbteNro, hash, codes);
    }

    public async Task<SafeWsfeResult> ConsultAsync(FiscalInvoice invoice, CancellationToken cancellationToken = default)
    {
        Validate(invoice, requireNumber: true);
        var response = await soap.ConsultAsync(await AuthAsync(invoice, cancellationToken), new Generated.FECompConsultaReq { PtoVta = invoice.PointOfSale, CbteTipo = invoice.VoucherType, CbteNro = invoice.VoucherNumber }, cancellationToken).ConfigureAwait(false);
        var result = response.ResultGet;
        var codes = Codes(response.Errors, response.Events, result?.Observaciones);
        var safe = new { operation = "FECompConsultar", result?.PtoVta, result?.CbteTipo, result?.CbteDesde, result?.CbteHasta, result?.Resultado, result?.CodAutorizacion, result?.FchVto, codes };
        var hash = Hash(safe);
        if (response.Errors is { Length: > 0 })
        {
            if (response.Errors.All(error => error.Code == 602) && result is null) return new SafeWsfeResult(false, false, null, null, hash, codes);
            throw new WsfeResponseException("ARCA consultation returned errors.", hash, codes);
        }
        if (result is null) throw new WsfeResponseException("ARCA consultation omitted ResultGet without error 602.", hash, codes);
        if (result.PtoVta != invoice.PointOfSale || result.CbteTipo != invoice.VoucherType || result.CbteDesde != invoice.VoucherNumber || result.CbteHasta != invoice.VoucherNumber)
            throw new WsfeResponseException("ARCA consultation identity mismatch.", hash, codes);
        return Result(result.Resultado, result.CodAutorizacion, result.FchVto, hash, codes);
    }

    public async Task<SafeWsfeResult> AuthorizeAsync(FiscalInvoice invoice, CancellationToken cancellationToken = default)
    {
        Validate(invoice, requireNumber: true);
        var response = await soap.AuthorizeAsync(await AuthAsync(invoice, cancellationToken), Request(invoice), cancellationToken).ConfigureAwait(false);
        var detail = response.FeDetResp is { Length: 1 } ? response.FeDetResp[0] : null;
        var codes = Codes(response.Errors, response.Events, detail?.Observaciones);
        var safe = new { operation = "FECAESolicitar", response.FeCabResp?.PtoVta, response.FeCabResp?.CbteTipo, response.FeCabResp?.CantReg, HeaderResult = response.FeCabResp?.Resultado, detail?.CbteDesde, detail?.CbteHasta, DetailResult = detail?.Resultado, detail?.CAE, detail?.CAEFchVto, codes };
        var hash = Hash(safe);
        if (response.Errors is { Length: > 0 } || response.FeCabResp is null || detail is null)
            throw new WsfeResponseException("ARCA authorization returned a general error or incomplete response.", hash, codes);
        if (response.FeCabResp.PtoVta != invoice.PointOfSale || response.FeCabResp.CbteTipo != invoice.VoucherType || response.FeCabResp.CantReg != 1 || detail.CbteDesde != invoice.VoucherNumber || detail.CbteHasta != invoice.VoucherNumber)
            throw new WsfeResponseException("ARCA authorization identity mismatch.", hash, codes);
        return Result(detail.Resultado, detail.CAE, detail.CAEFchVto, hash, codes);
    }

    private async Task<Generated.FEAuthRequest> AuthAsync(FiscalInvoice invoice, CancellationToken cancellationToken)
    {
        var access = await credentials.GetAsync(cancellationToken).ConfigureAwait(false);
        return new Generated.FEAuthRequest { Token = access.Token, Sign = access.Sign, Cuit = long.Parse(invoice.TaxpayerCuit, CultureInfo.InvariantCulture) };
    }

    private static Generated.FECAERequest Request(FiscalInvoice invoice) => new()
    {
        FeCabReq = new Generated.FECAECabRequest { CantReg = 1, PtoVta = invoice.PointOfSale, CbteTipo = invoice.VoucherType },
        FeDetReq = [new Generated.FECAEDetRequest
        {
            Concepto = invoice.Concept,
            DocTipo = invoice.RecipientDocumentType,
            DocNro = long.Parse(invoice.RecipientDocument, CultureInfo.InvariantCulture),
            CbteDesde = invoice.VoucherNumber,
            CbteHasta = invoice.VoucherNumber,
            CbteFch = Date(invoice.IssuedOn),
            ImpTotal = Money(invoice.TotalMinorUnits),
            ImpTotConc = Money(invoice.NonTaxedMinorUnits),
            ImpNeto = Money(invoice.NetMinorUnits),
            ImpOpEx = Money(invoice.ExemptMinorUnits),
            ImpTrib = Money(invoice.OtherTaxMinorUnits),
            ImpIVA = Money(invoice.VatMinorUnits),
            FchServDesde = OptionalDate(invoice.ServiceFrom),
            FchServHasta = OptionalDate(invoice.ServiceUntil),
            FchVtoPago = OptionalDate(invoice.PaymentDueOn),
            MonId = "PES",
            MonCotiz = 1d,
            MonCotizSpecified = true,
            CondicionIVAReceptorId = invoice.RecipientVatConditionId,
            CondicionIVAReceptorIdSpecified = true,
            Iva = invoice.VatLines.Count == 0 ? null! : invoice.VatLines.Select(line => new Generated.AlicIva { Id = line.Id, BaseImp = Money(line.BaseMinorUnits), Importe = Money(line.AmountMinorUnits) }).ToArray(),
            Tributos = invoice.OtherTaxLines.Count == 0 ? null! : invoice.OtherTaxLines.Select(line => new Generated.Tributo { Id = line.Id, Desc = line.Description, BaseImp = Money(line.BaseMinorUnits), Alic = Rate(line.RateBasisPoints), Importe = Money(line.AmountMinorUnits) }).ToArray(),
            CbtesAsoc = invoice.AssociatedVouchers.Count == 0 ? null! : invoice.AssociatedVouchers.Select(value => new Generated.CbteAsoc { Tipo = value.VoucherType, PtoVta = value.PointOfSale, Nro = value.Number, Cuit = value.TaxpayerCuit, CbteFch = Date(value.IssuedOn) }).ToArray()
        }]
    };

    private static SafeWsfeResult Result(string? result, string? cae, string? expiry, string hash, IReadOnlyList<string> codes)
    {
        if (result == "R") return new SafeWsfeResult(true, false, null, null, hash, codes);
        if (result != "A" || string.IsNullOrWhiteSpace(cae) || cae.Any(c => !char.IsAsciiDigit(c)) || cae.Length is < 8 or > 20 || !DateOnly.TryParseExact(expiry, "yyyyMMdd", CultureInfo.InvariantCulture, DateTimeStyles.None, out var date))
            throw new WsfeResponseException("ARCA result is neither an exact approval nor rejection.", hash, codes);
        return new SafeWsfeResult(true, true, cae, date, hash, codes);
    }

    private static IReadOnlyList<string> Codes(Generated.Err[]? errors, Generated.Evt[]? events, Generated.Obs[]? observations) =>
        (errors ?? []).Select(value => $"E:{value.Code}").Concat((events ?? []).Select(value => $"V:{value.Code}")).Concat((observations ?? []).Select(value => $"O:{value.Code}")).ToArray();

    private static string Hash(object value) => Convert.ToHexStringLower(SHA256.HashData(JsonSerializer.SerializeToUtf8Bytes(value)));
    private static string Date(DateOnly value) => value.ToString("yyyyMMdd", CultureInfo.InvariantCulture);
    private static string OptionalDate(DateOnly? value) => value is null ? string.Empty : Date(value.Value);
    private static double Money(long minor) => double.Parse((minor / 100m).ToString("0.00", CultureInfo.InvariantCulture), CultureInfo.InvariantCulture);
    private static double Rate(int basisPoints) => double.Parse((basisPoints / 100m).ToString("0.00", CultureInfo.InvariantCulture), CultureInfo.InvariantCulture);

    private static void Validate(FiscalInvoice value, bool requireNumber)
    {
        if (value.Currency != "ARS" || value.TaxpayerCuit.Length != 11 || value.TaxpayerCuit.Any(c => !char.IsAsciiDigit(c)) || !long.TryParse(value.RecipientDocument, NumberStyles.None, CultureInfo.InvariantCulture, out _) || value.PointOfSale is < 1 or > 99999 || value.VoucherType is < 1 or > 999 || value.Concept is < 1 or > 3 || value.RecipientDocumentType is < 0 or > 999 || value.RecipientVatConditionId is < 1 or > 999 || (requireNumber && value.VoucherNumber <= 0))
            throw new ArgumentException("Fiscal invoice identity is invalid.", nameof(value));
        var total = checked(value.NetMinorUnits + value.VatMinorUnits + value.ExemptMinorUnits + value.NonTaxedMinorUnits + value.OtherTaxMinorUnits);
        if (total <= 0 || total != value.TotalMinorUnits || value.VatLines.Sum(line => line.AmountMinorUnits) != value.VatMinorUnits || value.VatLines.Sum(line => line.BaseMinorUnits) != value.NetMinorUnits || value.OtherTaxLines.Sum(line => line.AmountMinorUnits) != value.OtherTaxMinorUnits)
            throw new ArgumentException("Fiscal monetary detail is inconsistent.", nameof(value));
        if ((value.Concept == 1 && (value.ServiceFrom is not null || value.ServiceUntil is not null || value.PaymentDueOn is not null)) || (value.Concept != 1 && (value.ServiceFrom is null || value.ServiceUntil is null || value.PaymentDueOn is null)))
            throw new ArgumentException("Fiscal service dates are inconsistent with Concepto.", nameof(value));
        var originalType = value.VoucherType switch { 3 => 1, 8 => 6, 13 => 11, _ => 0 };
        if (originalType == 0)
        {
            if (value.VoucherType is not (1 or 6 or 11) || value.AssociatedVouchers.Count != 0)
                throw new ArgumentException("Unsupported voucher type or association.", nameof(value));
        }
        else if (value.AssociatedVouchers.Count != 1 || value.AssociatedVouchers[0] is not { Number: > 0, PointOfSale: >= 1 and <= 99999 } associated || associated.VoucherType != originalType || associated.TaxpayerCuit != value.TaxpayerCuit || associated.IssuedOn > value.IssuedOn)
            throw new ArgumentException("Credit note associated voucher is invalid.", nameof(value));
    }
}
````

### FILE: `arca/fiscal/bridge/tests/Elite.Arca.Wsfe.Bridge.Tests/Elite.Arca.Wsfe.Bridge.Tests.csproj`
```yaml
block_id: "MS-ARCA-WSFE-ADAPTER:test-project:v1"
operation: CREATE
path: "arca/fiscal/bridge/tests/Elite.Arca.Wsfe.Bridge.Tests/Elite.Arca.Wsfe.Bridge.Tests.csproj"
sha256: "0b6dcf81a025eba0407c80a60a65b662cf1c26d879dac9cc833b09e381f7b2cc"
provenance: AUTHORED
source: "local test project governed by Microsoft .NET 10"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````xml
<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <OutputType>Exe</OutputType>
    <TargetFramework>net10.0</TargetFramework>
    <ImplicitUsings>enable</ImplicitUsings>
    <Nullable>enable</Nullable>
    <TreatWarningsAsErrors>true</TreatWarningsAsErrors>
    <RestorePackagesWithLockFile>true</RestorePackagesWithLockFile>
    <NuGetAudit>true</NuGetAudit>
    <NuGetAuditMode>all</NuGetAuditMode>
  </PropertyGroup>
  <ItemGroup><ProjectReference Include="../../src/Elite.Arca.Wsfe.Bridge/Elite.Arca.Wsfe.Bridge.csproj" /></ItemGroup>
</Project>
````

### FILE: `arca/fiscal/bridge/tests/Elite.Arca.Wsfe.Bridge.Tests/packages.lock.json`
```yaml
block_id: "MS-ARCA-WSFE-ADAPTER:test-lock:v1"
operation: CREATE
path: "arca/fiscal/bridge/tests/Elite.Arca.Wsfe.Bridge.Tests/packages.lock.json"
sha256: "6c8d1c9bd5d747a995742ad43525d9753f08767660470e9fb185955fefc02263"
provenance: AUTHORED
source: "lock generated by Microsoft .NET SDK 10.0.400"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````json
{
  "version": 1,
  "dependencies": {
    "net10.0": {
      "Microsoft.Extensions.ObjectPool": {
        "type": "Transitive",
        "resolved": "10.0.0",
        "contentHash": "bpeCq0IYmVLACyEUMzFIOQX+zZUElG1t+nu1lSxthe7B+1oNYking7b91305+jNB6iwojp9fqTY9O+Nh7ULQxg=="
      },
      "System.Security.Cryptography.Pkcs": {
        "type": "Transitive",
        "resolved": "10.0.11",
        "contentHash": "8IV+rI3xN/Mkq9MsSX7VZTv9T9Wt+tyhHLwBX9VNZBk8m8kknEs1JeTSjI2x9M2L/fSja+c2WSE5n5Yy0GXdoQ=="
      },
      "System.Security.Cryptography.Xml": {
        "type": "Transitive",
        "resolved": "10.0.11",
        "contentHash": "TokfVsaU2fmcwvsU6ihLCkseVfwhLX6Tmbl3qh3nOBIqHSX4yqFllpslA+PuHERH/ENu3kwqo4deZiyT0Wvh3Q==",
        "dependencies": {
          "System.Security.Cryptography.Pkcs": "10.0.11"
        }
      },
      "System.ServiceModel.Http": {
        "type": "Transitive",
        "resolved": "10.0.652802",
        "contentHash": "G02XZvmccf42QCU5MjviBIg69MSMAVHwL1inVPsNSpfp5g+t5BkQM3DyvWRLN4qmeFDWSF/mA1rIYONIDu/6Dg==",
        "dependencies": {
          "System.ServiceModel.Primitives": "10.0.652802"
        }
      },
      "System.ServiceModel.Primitives": {
        "type": "Transitive",
        "resolved": "10.0.652802",
        "contentHash": "ULfGNl75BNXkpF42wNV2CDXJ64dUZZEa8xO2mBsc4tqbW9QjruxjEB6bAr4Z/T1rNU+leOztIjCJQYsBGFWYlw==",
        "dependencies": {
          "Microsoft.Extensions.ObjectPool": "10.0.0",
          "System.Security.Cryptography.Xml": "10.0.0"
        }
      },
      "elite.arca.credentials": {
        "type": "Project",
        "dependencies": {
          "System.Security.Cryptography.Pkcs": "[10.0.11, )"
        }
      },
      "elite.arca.wsfe.bridge": {
        "type": "Project",
        "dependencies": {
          "Elite.Arca.Credentials": "[1.0.0, )",
          "Elite.Arca.Wsfe.Generated": "[1.0.0, )"
        }
      },
      "elite.arca.wsfe.generated": {
        "type": "Project",
        "dependencies": {
          "System.Security.Cryptography.Pkcs": "[10.0.11, )",
          "System.Security.Cryptography.Xml": "[10.0.11, )",
          "System.ServiceModel.Http": "[10.0.652802, )"
        }
      }
    }
  }
}
````

### FILE: `arca/fiscal/bridge/tests/Elite.Arca.Wsfe.Bridge.Tests/Program.cs`
```yaml
block_id: "MS-ARCA-WSFE-ADAPTER:tests:v1"
operation: CREATE
path: "arca/fiscal/bridge/tests/Elite.Arca.Wsfe.Bridge.Tests/Program.cs"
sha256: "e7e7fa86684981929390ad5d5427847e7b72511f76bc18dfb43144e4af6c3212"
provenance: AUTHORED
source: "local contract tests governed by ARCA manual 4.6"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````csharp
using Elite.Arca.Credentials;
using Elite.Arca.Wsfe.Bridge;
using Generated = Elite.Arca.Wsfe.Generated;

var invoice = new FiscalInvoice
{
    TaxpayerCuit = "33693450239", PointOfSale = 12, VoucherType = 1, Concept = 1,
    RecipientDocumentType = 80, RecipientDocument = "20111111112", RecipientVatConditionId = 1,
    Currency = "ARS", VoucherNumber = 1, TotalMinorUnits = 18405, NetMinorUnits = 15000,
    VatMinorUnits = 2625, ExemptMinorUnits = 0, NonTaxedMinorUnits = 0, OtherTaxMinorUnits = 780,
    IssuedOn = new DateOnly(2010, 9, 3),
    VatLines = [new(5, 10000, 2100), new(4, 5000, 525)],
    OtherTaxLines = [new(99, "Impuesto Municipal Matanza", 15000, 520, 780)]
};

var soap = new FakeSoap();
var bridge = new WsfeBridge(new FakeAccess(), soap);
var last = await bridge.LastAuthorizedAsync(invoice);
Check(last.Number == 0 && last.ResponseHash.Length == 64, "last authorized");

var approved = await bridge.AuthorizeAsync(invoice);
Check(approved is { Found: true, Authorized: true, Cae: "41124578989845" }, "authorization result");
Check(approved.CaeExpiresOn == new DateOnly(2010, 9, 13), "CAE expiry");
Check(soap.LastAuth is { Cuit: 33693450239 } && soap.LastAuth.Token == "token" && soap.LastAuth.Sign == "sign", "WSAA auth mapping");
var detail = soap.LastRequest!.FeDetReq.Single();
Check(soap.LastRequest.FeCabReq is { CantReg: 1, PtoVta: 12, CbteTipo: 1 }, "header mapping");
Check(detail.CondicionIVAReceptorIdSpecified && detail.CondicionIVAReceptorId == 1, "recipient VAT condition mapping");
Check(detail.ImpTotal == 184.05 && detail.ImpNeto == 150 && detail.ImpIVA == 26.25 && detail.ImpTrib == 7.8, "amount mapping");
Check(detail.Iva.Length == 2 && detail.Iva[0].Id == 5 && detail.Iva[1].BaseImp == 50, "VAT lines mapping");
Check(detail.Tributos.Single() is { Id: 99, Alic: 5.2, Importe: 7.8 }, "other-tax mapping");
Check(detail.CbtesAsoc is null, "ordinary invoice association omission");
Check(!System.Text.Json.JsonSerializer.Serialize(approved).Contains("token", StringComparison.OrdinalIgnoreCase), "safe response redaction");

var credit = new FiscalInvoice
{
    TaxpayerCuit = invoice.TaxpayerCuit, PointOfSale = invoice.PointOfSale, VoucherType = 3, Concept = invoice.Concept,
    RecipientDocumentType = invoice.RecipientDocumentType, RecipientDocument = invoice.RecipientDocument, RecipientVatConditionId = invoice.RecipientVatConditionId,
    Currency = invoice.Currency, VoucherNumber = 2, TotalMinorUnits = invoice.TotalMinorUnits, NetMinorUnits = invoice.NetMinorUnits,
    VatMinorUnits = invoice.VatMinorUnits, ExemptMinorUnits = invoice.ExemptMinorUnits, NonTaxedMinorUnits = invoice.NonTaxedMinorUnits,
    OtherTaxMinorUnits = invoice.OtherTaxMinorUnits, IssuedOn = invoice.IssuedOn, VatLines = invoice.VatLines, OtherTaxLines = invoice.OtherTaxLines,
    AssociatedVouchers = [new(invoice.TaxpayerCuit, 1, 12, 1, invoice.IssuedOn)]
};
await bridge.AuthorizeAsync(credit);
Check(soap.LastRequest!.FeDetReq.Single().CbtesAsoc.Single() is { Tipo: 1, PtoVta: 12, Nro: 1, Cuit: "33693450239", CbteFch: "20100903" }, "associated voucher mapping");

soap.ConsultNotFound = true;
var missing = await bridge.ConsultAsync(invoice);
Check(!missing.Found && !missing.Authorized && missing.ProviderCodes.SequenceEqual(["E:602"]), "manual error 602 mapping");

soap.ConsultNotFound = false;
soap.ConsultGeneralError = true;
await ThrowsAsync<WsfeResponseException>(() => bridge.ConsultAsync(invoice), "non-602 consult error");
soap.ConsultGeneralError = false;
soap.ConsultMismatch = true;
await ThrowsAsync<WsfeResponseException>(() => bridge.ConsultAsync(invoice), "consult identity mismatch");

var missingCondition = new FiscalInvoice
{
    TaxpayerCuit = invoice.TaxpayerCuit, PointOfSale = invoice.PointOfSale, VoucherType = invoice.VoucherType,
    Concept = invoice.Concept, RecipientDocumentType = invoice.RecipientDocumentType, RecipientDocument = invoice.RecipientDocument,
    RecipientVatConditionId = 0, Currency = invoice.Currency, VoucherNumber = invoice.VoucherNumber,
    TotalMinorUnits = invoice.TotalMinorUnits, NetMinorUnits = invoice.NetMinorUnits, VatMinorUnits = invoice.VatMinorUnits,
    ExemptMinorUnits = 0, NonTaxedMinorUnits = 0, OtherTaxMinorUnits = invoice.OtherTaxMinorUnits,
    IssuedOn = invoice.IssuedOn, VatLines = invoice.VatLines, OtherTaxLines = invoice.OtherTaxLines
};
await ThrowsAsync<ArgumentException>(() => bridge.AuthorizeAsync(missingCondition), "missing recipient VAT condition");

var creditWithoutAssociation = new FiscalInvoice
{
    TaxpayerCuit = credit.TaxpayerCuit, PointOfSale = credit.PointOfSale, VoucherType = credit.VoucherType, Concept = credit.Concept,
    RecipientDocumentType = credit.RecipientDocumentType, RecipientDocument = credit.RecipientDocument, RecipientVatConditionId = credit.RecipientVatConditionId,
    Currency = credit.Currency, VoucherNumber = credit.VoucherNumber, TotalMinorUnits = credit.TotalMinorUnits, NetMinorUnits = credit.NetMinorUnits,
    VatMinorUnits = credit.VatMinorUnits, ExemptMinorUnits = credit.ExemptMinorUnits, NonTaxedMinorUnits = credit.NonTaxedMinorUnits,
    OtherTaxMinorUnits = credit.OtherTaxMinorUnits, IssuedOn = credit.IssuedOn, VatLines = credit.VatLines, OtherTaxLines = credit.OtherTaxLines
};
await ThrowsAsync<ArgumentException>(() => bridge.AuthorizeAsync(creditWithoutAssociation), "credit without association");

Console.WriteLine("ARCA_WSFE_BRIDGE_TEST_PASS cases=15 secrets_returned=0 production_admitted=false");

static void Check(bool condition, string name)
{
    if (!condition) throw new InvalidOperationException($"Failed: {name}");
}

static async Task ThrowsAsync<T>(Func<Task> action, string name) where T : Exception
{
    try { await action(); }
    catch (T) { return; }
    throw new InvalidOperationException($"Expected {typeof(T).Name}: {name}");
}

sealed class FakeAccess : IWsaaAccessSource
{
    public Task<LoginTicketAccess> GetAsync(CancellationToken cancellationToken) => Task.FromResult(new LoginTicketAccess("token", "sign", DateTimeOffset.UtcNow.AddHours(1)));
}

sealed class FakeSoap : IWsfeSoap
{
    public Generated.FEAuthRequest? LastAuth { get; private set; }
    public Generated.FECAERequest? LastRequest { get; private set; }
    public bool ConsultNotFound { get; set; }
    public bool ConsultGeneralError { get; set; }
    public bool ConsultMismatch { get; set; }

    public Task<Generated.FERecuperaLastCbteResponse> LastAuthorizedAsync(Generated.FEAuthRequest auth, int pointOfSale, int voucherType, CancellationToken cancellationToken)
    {
        LastAuth = auth;
        return Task.FromResult(new Generated.FERecuperaLastCbteResponse { PtoVta = pointOfSale, CbteTipo = voucherType, CbteNro = 0, Events = [new Generated.Evt { Code = 1, Msg = "event text must not cross boundary" }] });
    }

    public Task<Generated.FECompConsultaResponse> ConsultAsync(Generated.FEAuthRequest auth, Generated.FECompConsultaReq request, CancellationToken cancellationToken)
    {
        if (ConsultNotFound) return Task.FromResult(new Generated.FECompConsultaResponse { Errors = [new Generated.Err { Code = 602, Msg = "No existen datos" }] });
        if (ConsultGeneralError) return Task.FromResult(new Generated.FECompConsultaResponse { Errors = [new Generated.Err { Code = 600, Msg = "secret-bearing provider text" }] });
        return Task.FromResult(new Generated.FECompConsultaResponse
        {
            ResultGet = new Generated.FECompConsResponse { PtoVta = ConsultMismatch ? request.PtoVta + 1 : request.PtoVta, CbteTipo = request.CbteTipo, CbteDesde = request.CbteNro, CbteHasta = request.CbteNro, Resultado = "A", CodAutorizacion = "41124578989845", FchVto = "20100913" }
        });
    }

    public Task<Generated.FECAEResponse> AuthorizeAsync(Generated.FEAuthRequest auth, Generated.FECAERequest request, CancellationToken cancellationToken)
    {
        LastAuth = auth;
        LastRequest = request;
        return Task.FromResult(new Generated.FECAEResponse
        {
            FeCabResp = new Generated.FECAECabResponse { PtoVta = request.FeCabReq.PtoVta, CbteTipo = request.FeCabReq.CbteTipo, CantReg = 1, Resultado = "A" },
            FeDetResp = [new Generated.FECAEDetResponse { CbteDesde = request.FeDetReq[0].CbteDesde, CbteHasta = request.FeDetReq[0].CbteHasta, Resultado = "A", CAE = "41124578989845", CAEFchVto = "20100913" }]
        });
    }
}
````

### FILE: `arca/fiscal/bridge/src/Elite.Arca.Wsfe.Bridge/WsfeParameterBridge.cs`
```yaml
block_id: "MS-ARCA-WSFE-ADAPTER:parameter-bridge:v1"
operation: CREATE
path: "arca/fiscal/bridge/src/Elite.Arca.Wsfe.Bridge/WsfeParameterBridge.cs"
sha256: "ba3ef0524cafdf4f0d33133ee5bfa60d83787c3f7178202109dd5742a76715f5"
provenance: AUTHORED
source: "local parameter bridge governed by ARCA WSFEv1 manual 4.6"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````csharp
using System.Globalization;
using System.Security.Cryptography;
using System.Text.Json;
using Elite.Arca.Credentials;
using Generated = Elite.Arca.Wsfe.Generated;

namespace Elite.Arca.Wsfe.Bridge;

public sealed class WsfeParameterBridge(IWsaaAccessSource credentials, IWsfeParameterSoap soap)
{
    public async Task<SafeParameterResult> VoucherTypesAsync(string taxpayerCuit, CancellationToken cancellationToken = default)
    {
        var response = await soap.VoucherTypesAsync(await AuthAsync(taxpayerCuit, cancellationToken), cancellationToken).ConfigureAwait(false);
        return Result("voucher_type", null, (response.ResultGet ?? []).Select(value => Item(value.Id, value.Desc, value.FchDesde, value.FchHasta)), response.Errors, response.Events);
    }

    public async Task<SafeParameterResult> ConceptsAsync(string taxpayerCuit, CancellationToken cancellationToken = default)
    {
        var response = await soap.ConceptsAsync(await AuthAsync(taxpayerCuit, cancellationToken), cancellationToken).ConfigureAwait(false);
        return Result("concept", null, (response.ResultGet ?? []).Select(value => Item(value.Id, value.Desc, value.FchDesde, value.FchHasta)), response.Errors, response.Events);
    }

    public async Task<SafeParameterResult> DocumentTypesAsync(string taxpayerCuit, CancellationToken cancellationToken = default)
    {
        var response = await soap.DocumentTypesAsync(await AuthAsync(taxpayerCuit, cancellationToken), cancellationToken).ConfigureAwait(false);
        return Result("document_type", null, (response.ResultGet ?? []).Select(value => Item(value.Id, value.Desc, value.FchDesde, value.FchHasta)), response.Errors, response.Events);
    }

    public async Task<SafeParameterResult> VatRatesAsync(string taxpayerCuit, CancellationToken cancellationToken = default)
    {
        var response = await soap.VatRatesAsync(await AuthAsync(taxpayerCuit, cancellationToken), cancellationToken).ConfigureAwait(false);
        return Result("vat_rate", null, (response.ResultGet ?? []).Select(value => Item(value.Id, value.Desc, value.FchDesde, value.FchHasta)), response.Errors, response.Events);
    }

    public async Task<SafeParameterResult> OtherTaxesAsync(string taxpayerCuit, CancellationToken cancellationToken = default)
    {
        var response = await soap.OtherTaxesAsync(await AuthAsync(taxpayerCuit, cancellationToken), cancellationToken).ConfigureAwait(false);
        return Result("other_tax", null, (response.ResultGet ?? []).Select(value => Item(value.Id, value.Desc, value.FchDesde, value.FchHasta)), response.Errors, response.Events);
    }

    public async Task<SafeParameterResult> PointsOfSaleAsync(string taxpayerCuit, CancellationToken cancellationToken = default)
    {
        var response = await soap.PointsOfSaleAsync(await AuthAsync(taxpayerCuit, cancellationToken), cancellationToken).ConfigureAwait(false);
        return Result("point_of_sale", null, (response.ResultGet ?? []).Select(value => new SafeParameterItem(value.Nro.ToString(CultureInfo.InvariantCulture), null, null, null, null, Required(value.EmisionTipo, "emission type"), Required(value.Bloqueado, "blocked state"), Date(value.FchBaja))), response.Errors, response.Events);
    }

    public async Task<SafeParameterResult> RecipientVatConditionsAsync(string taxpayerCuit, string voucherClass, CancellationToken cancellationToken = default)
    {
        if (voucherClass is not ("A" or "B" or "C" or "M")) throw new ArgumentException("Voucher class is outside the admitted WSFEv1 classes.", nameof(voucherClass));
        var response = await soap.RecipientVatConditionsAsync(await AuthAsync(taxpayerCuit, cancellationToken), voucherClass, cancellationToken).ConfigureAwait(false);
        return Result("recipient_vat_condition", voucherClass, (response.ResultGet ?? []).Select(value =>
        {
            if (!string.Equals(value.Cmp_Clase, voucherClass, StringComparison.Ordinal)) throw new InvalidDataException("ARCA recipient VAT class mismatch.");
            return new SafeParameterItem(value.Id.ToString(CultureInfo.InvariantCulture), Required(value.Desc, "description"), null, null, value.Cmp_Clase, null, null, null);
        }), response.Errors, response.Events);
    }

    private async Task<Generated.FEAuthRequest> AuthAsync(string taxpayerCuit, CancellationToken cancellationToken)
    {
        if (taxpayerCuit.Length != 11 || taxpayerCuit.Any(character => !char.IsAsciiDigit(character))) throw new ArgumentException("Taxpayer CUIT is invalid.", nameof(taxpayerCuit));
        var access = await credentials.GetAsync(cancellationToken).ConfigureAwait(false);
        return new Generated.FEAuthRequest { Token = access.Token, Sign = access.Sign, Cuit = long.Parse(taxpayerCuit, CultureInfo.InvariantCulture) };
    }

    private static SafeParameterItem Item<T>(T code, string description, string from, string until) where T : IFormattable =>
        new(code.ToString(null, CultureInfo.InvariantCulture), Required(description, "description"), DateRequired(from), Date(until), null, null, null, null);
    private static SafeParameterItem Item(string code, string description, string from, string until) =>
        new(RequiredCode(code), Required(description, "description"), DateRequired(from), Date(until), null, null, null, null);

    private static SafeParameterResult Result(string kind, string? voucherClass, IEnumerable<SafeParameterItem> values, Generated.Err[]? errors, Generated.Evt[]? events)
    {
        var codes = (errors ?? []).Select(value => $"E:{value.Code}").Concat((events ?? []).Select(value => $"V:{value.Code}")).ToArray();
        var items = values.OrderBy(value => value.Code, StringComparer.Ordinal).ToArray();
        if (errors is { Length: > 0 } || items.Length == 0 || items.Select(value => value.Code).Distinct(StringComparer.Ordinal).Count() != items.Length)
            throw new WsfeResponseException("ARCA parameter response is invalid.", Hash(new { kind, voucherClass, items, codes }), codes);
        foreach (var item in items)
            if (item.ValidFrom is not null && item.ValidUntil is not null && item.ValidUntil < item.ValidFrom)
                throw new InvalidDataException("ARCA parameter validity interval is invalid.");
        return new SafeParameterResult(kind, voucherClass, items, Hash(new { kind, voucherClass, items, codes }), codes);
    }

    private static string Required(string value, string name) => !string.IsNullOrWhiteSpace(value) && value.Length <= 250 ? value : throw new InvalidDataException($"ARCA parameter {name} is invalid.");
    private static string RequiredCode(string value) => !string.IsNullOrWhiteSpace(value) && value.Length <= 16 ? value : throw new InvalidDataException("ARCA parameter code is invalid.");
    private static DateOnly DateRequired(string value) => Date(value) ?? throw new InvalidDataException("ARCA parameter start date is required.");
    private static DateOnly? Date(string? value)
    {
        if (string.IsNullOrEmpty(value)) return null;
        if (!DateOnly.TryParseExact(value, "yyyyMMdd", CultureInfo.InvariantCulture, DateTimeStyles.None, out var date)) throw new InvalidDataException("ARCA parameter date is invalid.");
        return date;
    }
    private static string Hash(object value) => Convert.ToHexStringLower(SHA256.HashData(JsonSerializer.SerializeToUtf8Bytes(value)));
}
````

### FILE: `arca/fiscal/bridge/tests/Elite.Arca.Wsfe.Parameter.Tests/Elite.Arca.Wsfe.Parameter.Tests.csproj`
```yaml
block_id: "MS-ARCA-WSFE-ADAPTER:parameter-test-project:v1"
operation: CREATE
path: "arca/fiscal/bridge/tests/Elite.Arca.Wsfe.Parameter.Tests/Elite.Arca.Wsfe.Parameter.Tests.csproj"
sha256: "558b775b8175a3b51bfef8cf5d7c3c77adab48f022e89d762f5a6abaf054c564"
provenance: AUTHORED
source: "local contract-test project governed by Microsoft .NET documentation"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````xml
<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <OutputType>Exe</OutputType>
    <TargetFramework>net10.0</TargetFramework>
    <ImplicitUsings>enable</ImplicitUsings>
    <Nullable>enable</Nullable>
    <TreatWarningsAsErrors>true</TreatWarningsAsErrors>
    <AnalysisLevel>latest</AnalysisLevel>
    <RestorePackagesWithLockFile>true</RestorePackagesWithLockFile>
    <NuGetAudit>true</NuGetAudit>
    <NuGetAuditMode>all</NuGetAuditMode>
  </PropertyGroup>
  <ItemGroup>
    <ProjectReference Include="../../src/Elite.Arca.Wsfe.Bridge/Elite.Arca.Wsfe.Bridge.csproj" />
  </ItemGroup>
</Project>
````

### FILE: `arca/fiscal/bridge/tests/Elite.Arca.Wsfe.Parameter.Tests/packages.lock.json`
```yaml
block_id: "MS-ARCA-WSFE-ADAPTER:parameter-test-lock:v1"
operation: CREATE
path: "arca/fiscal/bridge/tests/Elite.Arca.Wsfe.Parameter.Tests/packages.lock.json"
sha256: "6c8d1c9bd5d747a995742ad43525d9753f08767660470e9fb185955fefc02263"
provenance: AUTHORED
source: "lock generated by Microsoft .NET SDK 10.0.400"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````json
{
  "version": 1,
  "dependencies": {
    "net10.0": {
      "Microsoft.Extensions.ObjectPool": {
        "type": "Transitive",
        "resolved": "10.0.0",
        "contentHash": "bpeCq0IYmVLACyEUMzFIOQX+zZUElG1t+nu1lSxthe7B+1oNYking7b91305+jNB6iwojp9fqTY9O+Nh7ULQxg=="
      },
      "System.Security.Cryptography.Pkcs": {
        "type": "Transitive",
        "resolved": "10.0.11",
        "contentHash": "8IV+rI3xN/Mkq9MsSX7VZTv9T9Wt+tyhHLwBX9VNZBk8m8kknEs1JeTSjI2x9M2L/fSja+c2WSE5n5Yy0GXdoQ=="
      },
      "System.Security.Cryptography.Xml": {
        "type": "Transitive",
        "resolved": "10.0.11",
        "contentHash": "TokfVsaU2fmcwvsU6ihLCkseVfwhLX6Tmbl3qh3nOBIqHSX4yqFllpslA+PuHERH/ENu3kwqo4deZiyT0Wvh3Q==",
        "dependencies": {
          "System.Security.Cryptography.Pkcs": "10.0.11"
        }
      },
      "System.ServiceModel.Http": {
        "type": "Transitive",
        "resolved": "10.0.652802",
        "contentHash": "G02XZvmccf42QCU5MjviBIg69MSMAVHwL1inVPsNSpfp5g+t5BkQM3DyvWRLN4qmeFDWSF/mA1rIYONIDu/6Dg==",
        "dependencies": {
          "System.ServiceModel.Primitives": "10.0.652802"
        }
      },
      "System.ServiceModel.Primitives": {
        "type": "Transitive",
        "resolved": "10.0.652802",
        "contentHash": "ULfGNl75BNXkpF42wNV2CDXJ64dUZZEa8xO2mBsc4tqbW9QjruxjEB6bAr4Z/T1rNU+leOztIjCJQYsBGFWYlw==",
        "dependencies": {
          "Microsoft.Extensions.ObjectPool": "10.0.0",
          "System.Security.Cryptography.Xml": "10.0.0"
        }
      },
      "elite.arca.credentials": {
        "type": "Project",
        "dependencies": {
          "System.Security.Cryptography.Pkcs": "[10.0.11, )"
        }
      },
      "elite.arca.wsfe.bridge": {
        "type": "Project",
        "dependencies": {
          "Elite.Arca.Credentials": "[1.0.0, )",
          "Elite.Arca.Wsfe.Generated": "[1.0.0, )"
        }
      },
      "elite.arca.wsfe.generated": {
        "type": "Project",
        "dependencies": {
          "System.Security.Cryptography.Pkcs": "[10.0.11, )",
          "System.Security.Cryptography.Xml": "[10.0.11, )",
          "System.ServiceModel.Http": "[10.0.652802, )"
        }
      }
    }
  }
}
````

### FILE: `arca/fiscal/bridge/tests/Elite.Arca.Wsfe.Parameter.Tests/Program.cs`
```yaml
block_id: "MS-ARCA-WSFE-ADAPTER:parameter-tests:v1"
operation: CREATE
path: "arca/fiscal/bridge/tests/Elite.Arca.Wsfe.Parameter.Tests/Program.cs"
sha256: "48ba1c5853cf6a685dbba39e5fcd15f799370fc9c1723a5dcd6bc02ba7d7917e"
provenance: AUTHORED
source: "local contract tests governed by ARCA WSFEv1 manual 4.6"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````csharp
using Elite.Arca.Credentials;
using Elite.Arca.Wsfe.Bridge;
using Generated = Elite.Arca.Wsfe.Generated;

var soap = new FakeParameters();
var bridge = new WsfeParameterBridge(new FakeAccess(), soap);

var vouchers = await bridge.VoucherTypesAsync("33693450239");
Check(vouchers.Kind == "voucher_type" && vouchers.Items.Single().Code == "1" && vouchers.Items.Single().ValidFrom == new DateOnly(2010, 1, 1), "voucher types");
Check(vouchers.ResponseHash.Length == 64 && vouchers.ProviderCodes.SequenceEqual(["V:10"]), "safe voucher receipt");
Check(soap.LastAuth is { Cuit: 33693450239, Token: "secret-token", Sign: "secret-sign" }, "credential mapping");

Check((await bridge.ConceptsAsync("33693450239")).Items.Single().Code == "1", "concepts");
Check((await bridge.DocumentTypesAsync("33693450239")).Items.Single().Code == "80", "document types");
Check((await bridge.VatRatesAsync("33693450239")).Items.Single().Code == "5", "VAT rates");
Check((await bridge.OtherTaxesAsync("33693450239")).Items.Single().Code == "99", "other taxes");
var point = (await bridge.PointsOfSaleAsync("33693450239")).Items.Single();
Check(point is { Code: "12", Description: null, EmissionType: "CAE", Blocked: "N", DeregisteredOn: null }, "points of sale");
var condition = await bridge.RecipientVatConditionsAsync("33693450239", "A");
Check(condition.VoucherClass == "A" && condition.Items.Single() is { Code: "1", VoucherClass: "A" }, "recipient VAT conditions");

soap.Error = true;
await Throws<WsfeResponseException>(() => bridge.VatRatesAsync("33693450239"), "provider errors");
soap.Error = false;
soap.Duplicate = true;
await Throws<WsfeResponseException>(() => bridge.DocumentTypesAsync("33693450239"), "duplicate codes");
soap.Duplicate = false;
soap.MismatchedClass = true;
await Throws<InvalidDataException>(() => bridge.RecipientVatConditionsAsync("33693450239", "A"), "class mismatch");
soap.MismatchedClass = false;
await Throws<ArgumentException>(() => bridge.RecipientVatConditionsAsync("33693450239", "Z"), "unknown class");
await Throws<ArgumentException>(() => bridge.ConceptsAsync("invalid"), "invalid CUIT");

var serialized = System.Text.Json.JsonSerializer.Serialize(vouchers);
Check(!serialized.Contains("secret", StringComparison.OrdinalIgnoreCase) && !serialized.Contains("provider message", StringComparison.OrdinalIgnoreCase), "no secrets or messages in result");
Console.WriteLine("ARCA_WSFE_PARAMETER_BRIDGE_PASS cases=15 secrets_returned=0 production_admitted=false");

static void Check(bool condition, string name)
{
    if (!condition) throw new InvalidOperationException($"FAILED: {name}");
}
static async Task Throws<T>(Func<Task> action, string name) where T : Exception
{
    try { await action(); }
    catch (T) { return; }
    throw new InvalidOperationException($"Expected {typeof(T).Name}: {name}");
}

sealed class FakeAccess : IWsaaAccessSource
{
    public Task<LoginTicketAccess> GetAsync(CancellationToken cancellationToken) => Task.FromResult(new LoginTicketAccess("secret-token", "secret-sign", DateTimeOffset.UtcNow.AddHours(1)));
}

sealed class FakeParameters : IWsfeParameterSoap
{
    internal Generated.FEAuthRequest? LastAuth { get; private set; }
    internal bool Error { get; set; }
    internal bool Duplicate { get; set; }
    internal bool MismatchedClass { get; set; }
    private Generated.Err[]? Errors => Error ? [new Generated.Err { Code = 500, Msg = "provider message secret" }] : null;
    private static Generated.Evt[] Events => [new Generated.Evt { Code = 10, Msg = "provider message" }];
    private T Capture<T>(Generated.FEAuthRequest auth, T response) { LastAuth = auth; return response; }

    public Task<Generated.CbteTipoResponse> VoucherTypesAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken) => Task.FromResult(Capture(auth, new Generated.CbteTipoResponse { ResultGet = [new Generated.CbteTipo { Id = 1, Desc = "Factura A", FchDesde = "20100101", FchHasta = "" }], Errors = Errors, Events = Events }));
    public Task<Generated.ConceptoTipoResponse> ConceptsAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken) => Task.FromResult(Capture(auth, new Generated.ConceptoTipoResponse { ResultGet = [new Generated.ConceptoTipo { Id = 1, Desc = "Productos", FchDesde = "20100101", FchHasta = "" }], Errors = Errors }));
    public Task<Generated.DocTipoResponse> DocumentTypesAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken)
    {
        var items = Duplicate ? new[] { Doc(), Doc() } : [Doc()];
        return Task.FromResult(Capture(auth, new Generated.DocTipoResponse { ResultGet = items, Errors = Errors }));
    }
    public Task<Generated.IvaTipoResponse> VatRatesAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken) => Task.FromResult(Capture(auth, new Generated.IvaTipoResponse { ResultGet = [new Generated.IvaTipo { Id = "5", Desc = "21%", FchDesde = "20100101", FchHasta = "" }], Errors = Errors }));
    public Task<Generated.FETributoResponse> OtherTaxesAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken) => Task.FromResult(Capture(auth, new Generated.FETributoResponse { ResultGet = [new Generated.TributoTipo { Id = 99, Desc = "Otros", FchDesde = "20100101", FchHasta = "" }], Errors = Errors }));
    public Task<Generated.FEPtoVentaResponse> PointsOfSaleAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken) => Task.FromResult(Capture(auth, new Generated.FEPtoVentaResponse { ResultGet = [new Generated.PtoVenta { Nro = 12, EmisionTipo = "CAE", Bloqueado = "N", FchBaja = "" }], Errors = Errors }));
    public Task<Generated.CondicionIvaReceptorResponse> RecipientVatConditionsAsync(Generated.FEAuthRequest auth, string voucherClass, CancellationToken cancellationToken) => Task.FromResult(Capture(auth, new Generated.CondicionIvaReceptorResponse { ResultGet = [new Generated.CondicionIvaReceptor { Id = 1, Desc = "IVA Responsable Inscripto", Cmp_Clase = MismatchedClass ? "B" : voucherClass }], Errors = Errors }));
    private static Generated.DocTipo Doc() => new() { Id = 80, Desc = "CUIT", FchDesde = "20100101", FchHasta = "" };
}
````

### FILE: `tools/test-arca-wsfe-soap-adapter.ps1`
```yaml
block_id: "MS-ARCA-WSFE-ADAPTER:runner:v1"
operation: CREATE
path: "tools/test-arca-wsfe-soap-adapter.ps1"
sha256: "c633ea8c9a0dcedf676cdd22b4ebff2fd695f7fb9bd2434607a711088f5b5a9f"
provenance: AUTHORED
source: "local locked runner governed by Microsoft .NET documentation"
license: "LicenseRef-Workspace-Owner"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding()]
param(
  [Parameter(Mandatory=$true)][string]$DotnetExecutable,
  [Parameter(Mandatory=$true)][string]$GeneratedRoot,
  [Parameter(Mandatory=$true)][string]$CredentialRoot,
  [Parameter(Mandatory=$true)][string]$WorkDirectory
)
$ErrorActionPreference='Stop'
$dotnet=[IO.Path]::GetFullPath($DotnetExecutable)
$generated=[IO.Path]::GetFullPath($GeneratedRoot)
$credentials=[IO.Path]::GetFullPath($CredentialRoot)
$work=[IO.Path]::GetFullPath($WorkDirectory)
if(-not(Test-Path $dotnet -PathType Leaf)){throw 'dotnet executable missing'}
if((@(& $dotnet --version 2>&1)[-1]).Trim()-ne'10.0.400'){throw 'dotnet SDK must equal 10.0.400'}
if(-not(Test-Path "$generated\Generated\WsfeV1Reference.cs" -PathType Leaf)-or-not(Test-Path "$generated\GENERATION_RECEIPT.json" -PathType Leaf)){throw 'generated homologation client/receipt missing'}
$receipt=Get-Content "$generated\GENERATION_RECEIPT.json" -Raw|ConvertFrom-Json
if($receipt.environment-ne'Homologation'-or$receipt.production_admitted-ne$false){throw 'only non-production homologation generation is admitted by this gate'}
$credentialProject=Join-Path $credentials 'arca\credentials\src\Elite.Arca.Credentials'
if(-not(Test-Path "$credentialProject\Elite.Arca.Credentials.csproj" -PathType Leaf)){throw 'credential core project missing'}
if(Test-Path $work){if((Get-ChildItem $work -Force).Count-ne0){throw 'work directory must be absent or empty'}}else{[IO.Directory]::CreateDirectory($work)|Out-Null}
$assets=[IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..\arca\fiscal\bridge'))
$bridge=Join-Path $work 'arca\fiscal\bridge';$generatedOut=Join-Path $work 'arca\fiscal\generated';$credentialOut=Join-Path $work 'arca\credentials\src\Elite.Arca.Credentials'
$sourceOut=Join-Path $bridge 'src\Elite.Arca.Wsfe.Bridge';$testOut=Join-Path $bridge 'tests\Elite.Arca.Wsfe.Bridge.Tests';$parameterTestOut=Join-Path $bridge 'tests\Elite.Arca.Wsfe.Parameter.Tests'
foreach($directory in @($bridge,$sourceOut,$testOut,$parameterTestOut,$generatedOut,$credentialOut)){[IO.Directory]::CreateDirectory($directory)|Out-Null}
foreach($name in @('Elite.Arca.Wsfe.Bridge.csproj','packages.lock.json','Contracts.cs','GeneratedAdapters.cs','WsfeBridge.cs','WsfeParameterBridge.cs')){Copy-Item (Join-Path $assets "src\Elite.Arca.Wsfe.Bridge\$name") $sourceOut}
foreach($name in @('Elite.Arca.Wsfe.Bridge.Tests.csproj','packages.lock.json','Program.cs')){Copy-Item (Join-Path $assets "tests\Elite.Arca.Wsfe.Bridge.Tests\$name") $testOut}
foreach($name in @('Elite.Arca.Wsfe.Parameter.Tests.csproj','packages.lock.json','Program.cs')){Copy-Item (Join-Path $assets "tests\Elite.Arca.Wsfe.Parameter.Tests\$name") $parameterTestOut}
Copy-Item "$assets\README.md" $bridge
foreach($name in @('Elite.Arca.Wsfe.Generated.csproj','packages.lock.json','GENERATION_RECEIPT.json')){Copy-Item (Join-Path $generated $name) $generatedOut}
Copy-Item "$generated\Generated" $generatedOut -Recurse
foreach($name in @('Elite.Arca.Credentials.csproj','packages.lock.json','CertificateStoreLoader.cs','CmsRequestSigner.cs','LoginTicketRequestFactory.cs','WsaaContracts.cs','WsaaCredentialProvider.cs')){Copy-Item (Join-Path $credentialProject $name) $credentialOut}
$project=Join-Path $bridge 'tests\Elite.Arca.Wsfe.Bridge.Tests\Elite.Arca.Wsfe.Bridge.Tests.csproj';$parameterProject=Join-Path $bridge 'tests\Elite.Arca.Wsfe.Parameter.Tests\Elite.Arca.Wsfe.Parameter.Tests.csproj'
$nuget=Join-Path $work 'NuGet.Config';[IO.File]::WriteAllText($nuget,'<?xml version="1.0" encoding="utf-8"?><configuration><packageSources><clear /></packageSources></configuration>',[Text.UTF8Encoding]::new($false))
$old=$env:DOTNET_CLI_TELEMETRY_OPTOUT;try{$env:DOTNET_CLI_TELEMETRY_OPTOUT='1';foreach($candidate in @($project,$parameterProject)){& $dotnet restore $candidate --locked-mode --configfile $nuget;if($LASTEXITCODE){throw "locked offline restore failed: $candidate"};& $dotnet build $candidate -c Release --no-restore --warnaserror;if($LASTEXITCODE){throw "warning-free build failed: $candidate"}};$invoice=@(& $dotnet run --project $project -c Release --no-build 2>&1);if($LASTEXITCODE-ne0-or($invoice-join"`n")-notmatch'ARCA_WSFE_BRIDGE_TEST_PASS cases=15'){throw "invoice contract tests failed`n$($invoice-join[Environment]::NewLine)"};$parameters=@(& $dotnet run --project $parameterProject -c Release --no-build 2>&1);if($LASTEXITCODE-ne0-or($parameters-join"`n")-notmatch'ARCA_WSFE_PARAMETER_BRIDGE_PASS cases=15 secrets_returned=0 production_admitted=false'){throw "parameter contract tests failed`n$($parameters-join[Environment]::NewLine)"}}finally{$env:DOTNET_CLI_TELEMETRY_OPTOUT=$old}
"ARCA_WSFE_ADAPTER_GATE_PASS environment=Homologation invoice_cases=15 parameter_cases=15 production_admitted=false"
````

## 6. Configuration surface

| Input | Secret | Validation/effect |
|---|---|---|
| generated homologation root | no | receipt must say Homologation and production_admitted=false |
| credential core root | no | source only; certificate is never a pack input |
| runtime certificate/ticket | yes | external credential core only; never returned or hashed |
| fiscal invoice | no | exact owner 0.3.x fields, integer units and approved IDs |
| parameter class | no | seven explicit WSFEv1 operation families; recipient class is A/B/C/M |

## 7. Dependency bill

| Dependency | Pin | Use | License/source |
|---|---:|---|---|
| .NET SDK | 10.0.400 | build/test | MIT/.NET terms; Microsoft |
| dotnet-svcutil | 8.0.0 | upstream proxy generation | Microsoft |
| System.ServiceModel.Http | 10.0.652802 | generated SOAP transport | MIT; Microsoft |
| cryptography packages | 10.0.11 | inherited WSAA core | MIT; Microsoft |
| ARCA WSFEv1 | manual 4.6 | protocol authority | ARCA public specification |

## 8. Apply order

Materialize generator, generate homologation clients, materialize credential core, then this adapter. Run its isolated locked tests. Acquire current parameter values through the explicit generated operations before configuring fiscal policy; never substitute a local guessed table. Only afterward add the Go↔.NET IPC worker; never bypass the PostgreSQL fiscal processor or call production by changing an endpoint.

## 9. Verification

Require 14/14 hash-exact materialization, .NET 10.0.400, locked offline restore after dependencies are acquired, warning-free builds, fifteen invoice contracts, fifteen parameter contracts and current NuGet vulnerable scan. Tests must inspect recipient VAT condition, amounts, both IVA lines, Tributo, CAE/date, exact associated-voucher mapping and rejection of a credit note without it, error 602, non-602 rejection, identity mismatch, all seven parameter families, deterministic hash, duplicates/dates/classes and safe response redaction.

## 10. Reconstruction evidence

V128 closes the owner fields required by this adapter; V129 records invoice mapping and V131 records exact parameter acquisition. V145 adds `CbtesAsoc` from the Microsoft-generated ARCA contract. Pack hashes, builds/tests/scan and remaining live conditions are recorded in `reconstruction_evidence/ARCA_ASSOCIATED_VOUCHER_EXECUTION_INVENTORY_2026-08-31_V145.md`.
