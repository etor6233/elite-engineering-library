# Microsoft-platform ARCA WSAA Credential Core

## 1. Metadata

```yaml
pack_id: "MICROSOFT-ARCA-WSAA-CREDENTIAL-CORE"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un núcleo WSAA compilable y su adapter exacto al cliente generado para request oficial, CMS/PKCS#7 Microsoft, certificate store, parser seguro y cache single-flight sin persistir secretos ni dejar glue manual."
stacks: [".NET SDK 10.0.400", "System.Security.Cryptography.Pkcs 10.0.11", "PowerShell 7"]
compatible_with: ["MICROSOFT-ARCA-WSFE-GENERATED-CLIENT@0.2.x", "ARCA-WSAA-OFFICIAL-CLIENT-SAMPLES@0.1.x"]
incompatible_with: []
license_expression: "LicenseRef-Workspace-Owner; Microsoft dependency retains MIT terms"
upstream_sources: ["https://www.arca.gob.ar/ws/documentacion/wsaa.asp", "https://www.arca.gob.ar/ws/WSAA/WSAAmanualDev.pdf", "https://www.arca.gob.ar/ws/documentacion/manuales/manual-desarrollador-ARCA-COMPG.pdf", "https://learn.microsoft.com/dotnet/api/system.security.cryptography.pkcs.signedcms.computesignature", "https://learn.microsoft.com/security/engineering/cryptographic-recommendations"]
verified_at: "2026-08-30"
```

All eighteen files are local `AUTHORED` integration code governed by ARCA protocol documents, the exact generated interface and Microsoft platform APIs. No file is represented as copied Microsoft or ARCA product source.

## 2. Applicability

Use after generating WSAA/WSFEv1 transport clients. It supplies reusable credential mechanics and the concrete `IWsaaTransport` mapper to the generated `loginCms` method. Do not use it as fiscal-rule, invoice-state or homologation evidence.

## 3. Architecture contract

Certificate material comes only from a certificate store by exact thumbprint; private keys/PFX/passwords/token/sign never enter configuration, logs or receipts. The request uses `service=wsfe`, positive external identity, UTC bounds and the official twelve-hour maximum. Microsoft `SignedCms` creates attached CMS. XML input rejects DTD and oversized/incomplete/expired responses. A semaphore makes refresh single-flight and the cache remains memory-only.

## 4. Exact file manifest

```text
CREATE arca/credentials/README.md
CREATE arca/credentials/src/Elite.Arca.Credentials/Elite.Arca.Credentials.csproj
CREATE arca/credentials/src/Elite.Arca.Credentials/packages.lock.json
CREATE arca/credentials/src/Elite.Arca.Credentials/WsaaContracts.cs
CREATE arca/credentials/src/Elite.Arca.Credentials/LoginTicketRequestFactory.cs
CREATE arca/credentials/src/Elite.Arca.Credentials/CmsRequestSigner.cs
CREATE arca/credentials/src/Elite.Arca.Credentials/CertificateStoreLoader.cs
CREATE arca/credentials/src/Elite.Arca.Credentials/WsaaCredentialProvider.cs
CREATE arca/credentials/integration/Elite.Arca.Wsaa.Transport/Elite.Arca.Wsaa.Transport.csproj
CREATE arca/credentials/integration/Elite.Arca.Wsaa.Transport/packages.lock.json
CREATE arca/credentials/integration/Elite.Arca.Wsaa.Transport/GeneratedWsaaTransport.cs
CREATE arca/credentials/integration-tests/Elite.Arca.Wsaa.Transport.Tests/Elite.Arca.Wsaa.Transport.Tests.csproj
CREATE arca/credentials/integration-tests/Elite.Arca.Wsaa.Transport.Tests/packages.lock.json
CREATE arca/credentials/integration-tests/Elite.Arca.Wsaa.Transport.Tests/Program.cs
CREATE arca/credentials/tests/Elite.Arca.Credentials.Tests/Elite.Arca.Credentials.Tests.csproj
CREATE arca/credentials/tests/Elite.Arca.Credentials.Tests/packages.lock.json
CREATE arca/credentials/tests/Elite.Arca.Credentials.Tests/Program.cs
CREATE tools/test-arca-wsaa-credential-core.ps1
```

## 5. Materialization blocks

### FILE: `arca/credentials/README.md`
```yaml
block_id: "ARCA-CREDENTIALS:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local boundary governed by ARCA/Microsoft official documentation"
license: "LicenseRef-Workspace-Owner"
sha256: "5bb8efabc91aad0a462e5f88f1b143cf79390e9e7d4d10425c1f068a0eff3ea7"
variables: []
secrets_allowed: false
```
````text
# Elite ARCA WSAA credential core

This library creates the official WSAA login ticket request, signs it as attached CMS/PKCS#7 with Microsoft platform cryptography, calls an injected `loginCms` transport and caches the returned token/sign in memory with single-flight refresh.

Production configuration must identify an ARCA-issued certificate already installed in an operating-system or managed certificate store. PFX bytes, passwords, token, sign and private keys are never configuration values, logs or receipts. The library never persists credentials.

The default invoice profile requests service `wsfe` for at most twelve hours. A generated WSAA proxy implements `IWsaaTransport`; fiscal issuance remains a separate durable boundary.
````

### FILE: `arca/credentials/src/Elite.Arca.Credentials/Elite.Arca.Credentials.csproj`
```yaml
block_id: "ARCA-CREDENTIALS:project:v1"
operation: CREATE
provenance: AUTHORED
source: "local exact Microsoft package selection"
license: "LicenseRef-Workspace-Owner"
sha256: "9cf2556a0f40eedab1b2f6d0b0109644cbdd3bd36e0bf06b0756d8217db496ce"
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
    <AnalysisLevel>latest</AnalysisLevel>
    <EnforceCodeStyleInBuild>true</EnforceCodeStyleInBuild>
    <RestorePackagesWithLockFile>true</RestorePackagesWithLockFile>
    <NuGetAudit>true</NuGetAudit>
    <NuGetAuditMode>all</NuGetAuditMode>
  </PropertyGroup>
  <ItemGroup>
    <PackageReference Include="System.Security.Cryptography.Pkcs" Version="10.0.11" />
  </ItemGroup>
</Project>
````

### FILE: `arca/credentials/src/Elite.Arca.Credentials/packages.lock.json`
```yaml
block_id: "ARCA-CREDENTIALS:source-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "NuGet generated dependency lock"
license: "LicenseRef-Workspace-Owner"
sha256: "c01c113a73dff4f219a24e0ab72987e0149d7367d55379dba7da3a5919fa5665"
variables: []
secrets_allowed: false
```
````json
{
  "version": 1,
  "dependencies": {
    "net10.0": {
      "System.Security.Cryptography.Pkcs": {
        "type": "Direct",
        "requested": "[10.0.11, )",
        "resolved": "10.0.11",
        "contentHash": "8IV+rI3xN/Mkq9MsSX7VZTv9T9Wt+tyhHLwBX9VNZBk8m8kknEs1JeTSjI2x9M2L/fSja+c2WSE5n5Yy0GXdoQ=="
      }
    }
  }
}
````

### FILE: `arca/credentials/src/Elite.Arca.Credentials/WsaaContracts.cs`
```yaml
block_id: "ARCA-CREDENTIALS:contracts:v1"
operation: CREATE
provenance: AUTHORED
source: "local typed contract governed by ARCA WSAA"
license: "LicenseRef-Workspace-Owner"
sha256: "7c24f75a02611bd2903b096881d11d91a83e30f911605bea2ebd7d5aa073b629"
variables: []
secrets_allowed: false
```
````csharp
using System.Security.Cryptography.X509Certificates;
using System.Text.RegularExpressions;

namespace Elite.Arca.Credentials;

public sealed record WsaaOptions(string Service, TimeSpan TicketLifetime, TimeSpan GenerationSkew, TimeSpan RefreshBefore)
{
    private static readonly Regex ServicePattern = new("^[a-z0-9_]{1,32}$", RegexOptions.CultureInvariant | RegexOptions.NonBacktracking);
    public static WsaaOptions ElectronicInvoice { get; } = new("wsfe", TimeSpan.FromHours(12), TimeSpan.FromMinutes(5), TimeSpan.FromMinutes(10));

    public void Validate()
    {
        if (!ServicePattern.IsMatch(Service)) throw new ArgumentException("Invalid ARCA service identifier.", nameof(Service));
        if (TicketLifetime <= TimeSpan.Zero || TicketLifetime > TimeSpan.FromHours(12)) throw new ArgumentOutOfRangeException(nameof(TicketLifetime));
        if (GenerationSkew < TimeSpan.Zero || GenerationSkew > TimeSpan.FromMinutes(10)) throw new ArgumentOutOfRangeException(nameof(GenerationSkew));
        if (RefreshBefore <= TimeSpan.Zero || RefreshBefore >= TicketLifetime) throw new ArgumentOutOfRangeException(nameof(RefreshBefore));
    }
}

public interface IUtcClock { DateTimeOffset UtcNow { get; } }
public interface ILoginTicketIdSource { long Next(); }
public interface IWsaaTransport { Task<string> LoginCmsAsync(string cmsBase64, CancellationToken cancellationToken); }
public interface ICmsRequestSigner { string Sign(X509Certificate2 certificate, ReadOnlySpan<byte> request); }

public sealed class LoginTicketAccess
{
    public LoginTicketAccess(string token, string sign, DateTimeOffset expiresAt)
    {
        Token = string.IsNullOrWhiteSpace(token) ? throw new ArgumentException("Token is required.", nameof(token)) : token;
        Sign = string.IsNullOrWhiteSpace(sign) ? throw new ArgumentException("Sign is required.", nameof(sign)) : sign;
        ExpiresAt = expiresAt;
    }
    public string Token { get; }
    public string Sign { get; }
    public DateTimeOffset ExpiresAt { get; }
    public override string ToString() => "LoginTicketAccess [REDACTED]";
}

public sealed class SystemUtcClock : IUtcClock { public DateTimeOffset UtcNow => DateTimeOffset.UtcNow; }
````

### FILE: `arca/credentials/src/Elite.Arca.Credentials/LoginTicketRequestFactory.cs`
```yaml
block_id: "ARCA-CREDENTIALS:request:v1"
operation: CREATE
provenance: AUTHORED
source: "local implementation of official loginTicketRequest shape"
license: "LicenseRef-Workspace-Owner"
sha256: "33b1bdb23ffd7d0d771b921a411cbca7ef17b28244c8ca3a5d619e905ce4ef65"
variables: []
secrets_allowed: false
```
````csharp
using System.Globalization;
using System.Text;
using System.Xml;

namespace Elite.Arca.Credentials;

public sealed class LoginTicketRequestFactory(WsaaOptions options, IUtcClock clock, ILoginTicketIdSource ids)
{
    public byte[] Create()
    {
        options.Validate();
        var id = ids.Next();
        if (id <= 0) throw new InvalidOperationException("Login ticket ID must be positive.");
        var now = clock.UtcNow.ToUniversalTime();
        var settings = new XmlWriterSettings { Encoding = new UTF8Encoding(false), OmitXmlDeclaration = false, Indent = false, NewLineHandling = NewLineHandling.None };
        using var stream = new MemoryStream();
        using (var writer = XmlWriter.Create(stream, settings))
        {
            writer.WriteStartDocument();
            writer.WriteStartElement("loginTicketRequest");
            writer.WriteAttributeString("version", "1.0");
            writer.WriteStartElement("header");
            writer.WriteElementString("uniqueId", id.ToString(CultureInfo.InvariantCulture));
            writer.WriteElementString("generationTime", XmlConvert.ToString(now - options.GenerationSkew));
            writer.WriteElementString("expirationTime", XmlConvert.ToString(now + options.TicketLifetime));
            writer.WriteEndElement();
            writer.WriteElementString("service", options.Service);
            writer.WriteEndElement();
            writer.WriteEndDocument();
        }
        return stream.ToArray();
    }
}
````

### FILE: `arca/credentials/src/Elite.Arca.Credentials/CmsRequestSigner.cs`
```yaml
block_id: "ARCA-CREDENTIALS:cms:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration with Microsoft SignedCms"
license: "LicenseRef-Workspace-Owner"
sha256: "e313ff3ad16ced8eedcaea59cc4424fe1bde3afc2047a0589b074895cf5739a2"
variables: []
secrets_allowed: false
```
````csharp
using System.Security.Cryptography;
using System.Security.Cryptography.Pkcs;
using System.Security.Cryptography.X509Certificates;

namespace Elite.Arca.Credentials;

public sealed class CmsRequestSigner : ICmsRequestSigner
{
    public string Sign(X509Certificate2 certificate, ReadOnlySpan<byte> request)
    {
        ArgumentNullException.ThrowIfNull(certificate);
        if (!certificate.HasPrivateKey) throw new CryptographicException("Certificate has no private key.");
        var now = DateTimeOffset.UtcNow;
        if (now < certificate.NotBefore || now >= certificate.NotAfter) throw new CryptographicException("Certificate is outside its validity interval.");
        using var rsa = certificate.GetRSAPrivateKey();
        if (rsa is null || rsa.KeySize < 2048) throw new CryptographicException("An RSA private key of at least 2048 bits is required.");
        var content = request.ToArray();
        try
        {
            var cms = new SignedCms(new ContentInfo(content), detached: false);
            var signer = new CmsSigner(SubjectIdentifierType.IssuerAndSerialNumber, certificate) { IncludeOption = X509IncludeOption.EndCertOnly };
            cms.ComputeSignature(signer);
            return Convert.ToBase64String(cms.Encode());
        }
        finally { CryptographicOperations.ZeroMemory(content); }
    }
}
````

### FILE: `arca/credentials/src/Elite.Arca.Credentials/CertificateStoreLoader.cs`
```yaml
block_id: "ARCA-CREDENTIALS:store:v1"
operation: CREATE
provenance: AUTHORED
source: "local Microsoft platform certificate-store boundary"
license: "LicenseRef-Workspace-Owner"
sha256: "c8a6c7e9b770d2c17a2b814e54d9bc266d1dc8193772777c7ef8740aeb3b5ada"
variables: []
secrets_allowed: false
```
````csharp
using System.Security.Cryptography;
using System.Security.Cryptography.X509Certificates;

namespace Elite.Arca.Credentials;

public static class CertificateStoreLoader
{
    public static X509Certificate2 LoadByThumbprint(string thumbprint, StoreLocation location = StoreLocation.CurrentUser)
    {
        var normalized = NormalizeThumbprint(thumbprint);
        using var store = new X509Store(StoreName.My, location, OpenFlags.ReadOnly | OpenFlags.OpenExistingOnly);
        var matches = store.Certificates.Find(X509FindType.FindByThumbprint, normalized, validOnly: false)
            .Where(c => string.Equals(NormalizeThumbprint(c.Thumbprint), normalized, StringComparison.Ordinal))
            .ToArray();
        if (matches.Length != 1)
        {
            foreach (var match in matches) match.Dispose();
            throw new CryptographicException("Certificate thumbprint must resolve to exactly one certificate.");
        }
        try { Validate(matches[0]); return matches[0]; }
        catch { matches[0].Dispose(); throw; }
    }

    public static string NormalizeThumbprint(string value)
    {
        ArgumentException.ThrowIfNullOrWhiteSpace(value);
        var normalized = string.Concat(value.Where(Uri.IsHexDigit)).ToUpperInvariant();
        if (normalized.Length is not (40 or 64)) throw new ArgumentException("Certificate thumbprint must be SHA-1 or SHA-256 hex.", nameof(value));
        return normalized;
    }

    public static void Validate(X509Certificate2 certificate)
    {
        if (!certificate.HasPrivateKey) throw new CryptographicException("Certificate has no private key.");
        var now = DateTimeOffset.UtcNow;
        if (now < certificate.NotBefore || now >= certificate.NotAfter) throw new CryptographicException("Certificate is outside its validity interval.");
        using var rsa = certificate.GetRSAPrivateKey();
        if (rsa is null || rsa.KeySize < 2048) throw new CryptographicException("An RSA private key of at least 2048 bits is required.");
    }
}
````

### FILE: `arca/credentials/src/Elite.Arca.Credentials/WsaaCredentialProvider.cs`
```yaml
block_id: "ARCA-CREDENTIALS:provider:v1"
operation: CREATE
provenance: AUTHORED
source: "local single-flight credential lifecycle"
license: "LicenseRef-Workspace-Owner"
sha256: "f7a7fac5e6ce7bf4dfca176dbd349544020c8be3945001c4e70b658e64478e4f"
variables: []
secrets_allowed: false
```
````csharp
using System.Security.Cryptography;
using System.Security.Cryptography.X509Certificates;
using System.Xml;
using System.Xml.Linq;

namespace Elite.Arca.Credentials;

public sealed class WsaaCredentialProvider(
    WsaaOptions options,
    IUtcClock clock,
    LoginTicketRequestFactory requests,
    ICmsRequestSigner signer,
    X509Certificate2 certificate,
    IWsaaTransport transport) : IDisposable
{
    private readonly SemaphoreSlim gate = new(1, 1);
    private LoginTicketAccess? cached;
    private bool disposed;

    public async Task<LoginTicketAccess> GetAsync(CancellationToken cancellationToken = default)
    {
        ObjectDisposedException.ThrowIf(disposed, this);
        var current = Volatile.Read(ref cached);
        if (IsFresh(current)) return current!;
        await gate.WaitAsync(cancellationToken).ConfigureAwait(false);
        try
        {
            current = cached;
            if (IsFresh(current)) return current!;
            var request = requests.Create();
            string cms;
            try { cms = signer.Sign(certificate, request); }
            finally { CryptographicOperations.ZeroMemory(request); }
            var response = await transport.LoginCmsAsync(cms, cancellationToken).ConfigureAwait(false);
            var parsed = WsaaLoginResponseParser.Parse(response, clock.UtcNow, options.TicketLifetime + TimeSpan.FromHours(1));
            Volatile.Write(ref cached, parsed);
            return parsed;
        }
        finally { gate.Release(); }
    }

    private bool IsFresh(LoginTicketAccess? value) => value is not null && value.ExpiresAt - options.RefreshBefore > clock.UtcNow;
    public void Dispose() { if (disposed) return; disposed = true; gate.Dispose(); certificate.Dispose(); }
}

public static class WsaaLoginResponseParser
{
    public static LoginTicketAccess Parse(string xml, DateTimeOffset now, TimeSpan maximumFuture)
    {
        ArgumentException.ThrowIfNullOrWhiteSpace(xml);
        if (xml.Length > 1_048_576) throw new InvalidDataException("WSAA response exceeds one MiB.");
        var settings = new XmlReaderSettings { DtdProcessing = DtdProcessing.Prohibit, XmlResolver = null, MaxCharactersInDocument = 1_048_576 };
        using var input = new StringReader(xml);
        using var reader = XmlReader.Create(input, settings);
        var document = XDocument.Load(reader, LoadOptions.None);
        string? Value(string name) => document.Descendants().FirstOrDefault(element => element.Name.LocalName == name)?.Value;
        var token = Value("token"); var sign = Value("sign"); var expiration = Value("expirationTime");
        if (token is null or { Length: > 16384 } || sign is null or { Length: > 16384 } || expiration is null) throw new InvalidDataException("WSAA response is incomplete.");
        var expiresAt = XmlConvert.ToDateTimeOffset(expiration).ToUniversalTime();
        if (expiresAt <= now || expiresAt > now + maximumFuture) throw new InvalidDataException("WSAA expiration is outside the admitted interval.");
        return new LoginTicketAccess(token, sign, expiresAt);
    }
}
````

### FILE: `arca/credentials/integration/Elite.Arca.Wsaa.Transport/Elite.Arca.Wsaa.Transport.csproj`
```yaml
block_id: "ARCA-CREDENTIALS:transport-project:v1"
operation: CREATE
provenance: AUTHORED
source: "local exact composition with generated Microsoft WCF client"
license: "LicenseRef-Workspace-Owner"
sha256: "1037cb0c1c4cab2e40c1a255f43a6c253f2fb27b36a22b7bcc4a88e9d0c09799"
variables: ["GeneratedClientRoot"]
secrets_allowed: false
```
````xml
<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
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
    <ProjectReference Include="..\..\src\Elite.Arca.Credentials\Elite.Arca.Credentials.csproj" />
    <PackageReference Include="System.Security.Cryptography.Pkcs" Version="10.0.11" />
    <PackageReference Include="System.Security.Cryptography.Xml" Version="10.0.11" />
    <PackageReference Include="System.ServiceModel.Http" Version="10.0.652802" />
    <Compile Include="$(GeneratedClientRoot)\Generated\WsaaReference.cs" Link="Generated\WsaaReference.cs" />
  </ItemGroup>
  <Target Name="ValidateGeneratedClient" BeforeTargets="PrepareForBuild">
    <Error Condition="'$(GeneratedClientRoot)' == ''" Text="GeneratedClientRoot is required." />
    <Error Condition="!Exists('$(GeneratedClientRoot)\Generated\WsaaReference.cs')" Text="Generated WSAA client is missing." />
  </Target>
</Project>
````

### FILE: `arca/credentials/integration/Elite.Arca.Wsaa.Transport/packages.lock.json`
```yaml
block_id: "ARCA-CREDENTIALS:transport-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "NuGet generated WCF/crypto dependency lock"
license: "LicenseRef-Workspace-Owner"
sha256: "e795876298ad858baa373a2a351132d2a2690a66c8ea80516df5b7b29885c317"
variables: []
secrets_allowed: false
```
````json
{
  "version": 1,
  "dependencies": {
    "net10.0": {
      "System.Security.Cryptography.Pkcs": {
        "type": "Direct",
        "requested": "[10.0.11, )",
        "resolved": "10.0.11",
        "contentHash": "8IV+rI3xN/Mkq9MsSX7VZTv9T9Wt+tyhHLwBX9VNZBk8m8kknEs1JeTSjI2x9M2L/fSja+c2WSE5n5Yy0GXdoQ=="
      },
      "System.Security.Cryptography.Xml": {
        "type": "Direct",
        "requested": "[10.0.11, )",
        "resolved": "10.0.11",
        "contentHash": "TokfVsaU2fmcwvsU6ihLCkseVfwhLX6Tmbl3qh3nOBIqHSX4yqFllpslA+PuHERH/ENu3kwqo4deZiyT0Wvh3Q==",
        "dependencies": {
          "System.Security.Cryptography.Pkcs": "10.0.11"
        }
      },
      "System.ServiceModel.Http": {
        "type": "Direct",
        "requested": "[10.0.652802, )",
        "resolved": "10.0.652802",
        "contentHash": "G02XZvmccf42QCU5MjviBIg69MSMAVHwL1inVPsNSpfp5g+t5BkQM3DyvWRLN4qmeFDWSF/mA1rIYONIDu/6Dg==",
        "dependencies": {
          "System.ServiceModel.Primitives": "10.0.652802"
        }
      },
      "Microsoft.Extensions.ObjectPool": {
        "type": "Transitive",
        "resolved": "10.0.0",
        "contentHash": "bpeCq0IYmVLACyEUMzFIOQX+zZUElG1t+nu1lSxthe7B+1oNYking7b91305+jNB6iwojp9fqTY9O+Nh7ULQxg=="
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
      }
    }
  }
}
````

### FILE: `arca/credentials/integration/Elite.Arca.Wsaa.Transport/GeneratedWsaaTransport.cs`
```yaml
block_id: "ARCA-CREDENTIALS:generated-transport:v1"
operation: CREATE
provenance: AUTHORED
source: "local adapter to exact dotnet-svcutil LoginCMS interface"
license: "LicenseRef-Workspace-Owner"
sha256: "000e55f7f32609540b08d72e3e099c289ce9c8a895cb3d223f734ea374979cfe"
variables: []
secrets_allowed: false
```
````csharp
using Elite.Arca.Credentials;
using Elite.Arca.Wsaa.Generated;

namespace Elite.Arca.Wsaa.Transport;

public sealed class GeneratedWsaaTransport(LoginCMS client) : IWsaaTransport
{
    public async Task<string> LoginCmsAsync(string cmsBase64, CancellationToken cancellationToken)
    {
        ArgumentException.ThrowIfNullOrWhiteSpace(cmsBase64);
        cancellationToken.ThrowIfCancellationRequested();
        var response = await client.loginCmsAsync(new loginCmsRequest(cmsBase64)).WaitAsync(cancellationToken).ConfigureAwait(false);
        if (string.IsNullOrWhiteSpace(response.loginCmsReturn)) throw new InvalidDataException("WSAA returned an empty login ticket response.");
        return response.loginCmsReturn;
    }
}
````

### FILE: `arca/credentials/integration-tests/Elite.Arca.Wsaa.Transport.Tests/Elite.Arca.Wsaa.Transport.Tests.csproj`
```yaml
block_id: "ARCA-CREDENTIALS:transport-test-project:v1"
operation: CREATE
provenance: AUTHORED
source: "local generated-interface integration harness"
license: "LicenseRef-Workspace-Owner"
sha256: "54d334daa3e64467bd8a4cf67d7cf7b804a1b759df5cbed2c17f3162abb744d9"
variables: ["GeneratedClientRoot"]
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
    <ProjectReference Include="..\..\integration\Elite.Arca.Wsaa.Transport\Elite.Arca.Wsaa.Transport.csproj" AdditionalProperties="GeneratedClientRoot=$(GeneratedClientRoot)" />
  </ItemGroup>
</Project>
````

### FILE: `arca/credentials/integration-tests/Elite.Arca.Wsaa.Transport.Tests/packages.lock.json`
```yaml
block_id: "ARCA-CREDENTIALS:transport-test-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "NuGet generated integration dependency lock"
license: "LicenseRef-Workspace-Owner"
sha256: "f12eaf3dc3498f4f24bacda497347bb4f70f4c2b865b14a7b79c0fbbed14d9a4"
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
      "elite.arca.wsaa.transport": {
        "type": "Project",
        "dependencies": {
          "Elite.Arca.Credentials": "[1.0.0, )",
          "System.Security.Cryptography.Pkcs": "[10.0.11, )",
          "System.Security.Cryptography.Xml": "[10.0.11, )",
          "System.ServiceModel.Http": "[10.0.652802, )"
        }
      }
    }
  }
}
````

### FILE: `arca/credentials/integration-tests/Elite.Arca.Wsaa.Transport.Tests/Program.cs`
```yaml
block_id: "ARCA-CREDENTIALS:transport-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local exact generated signature regression"
license: "LicenseRef-Workspace-Owner"
sha256: "ee28bbc0b56c98cb384aee1deec11cddce3cef2ccc78f610e6190bfdd5b50ca6"
variables: []
secrets_allowed: false
```
````csharp
using Elite.Arca.Wsaa.Generated;
using Elite.Arca.Wsaa.Transport;

var ok = new FakeClient("<loginTicketResponse/>");
var transport = new GeneratedWsaaTransport(ok);
var result = await transport.LoginCmsAsync("cms-value", CancellationToken.None);
if (result != "<loginTicketResponse/>" || ok.LastCms != "cms-value" || ok.Calls != 1) throw new InvalidOperationException("generated transport mapping failed");
try { _ = await new GeneratedWsaaTransport(new FakeClient(" ")).LoginCmsAsync("cms", CancellationToken.None); throw new InvalidOperationException("empty response accepted"); } catch (InvalidDataException) { }
try { _ = await transport.LoginCmsAsync(" ", CancellationToken.None); throw new InvalidOperationException("empty CMS accepted"); } catch (ArgumentException) { }
Console.WriteLine("ARCA_GENERATED_WSAA_TRANSPORT_TEST_PASS tests=3 manual_glue=0 production_admitted=false");

sealed class FakeClient(string response) : LoginCMS
{
    public int Calls { get; private set; }
    public string? LastCms { get; private set; }
    public Task<loginCmsResponse> loginCmsAsync(loginCmsRequest request)
    {
        Calls++; LastCms = request.in0; return Task.FromResult(new loginCmsResponse(response));
    }
}
````

### FILE: `arca/credentials/tests/Elite.Arca.Credentials.Tests/Elite.Arca.Credentials.Tests.csproj`
```yaml
block_id: "ARCA-CREDENTIALS:test-project:v1"
operation: CREATE
provenance: AUTHORED
source: "local dependency-free executable test harness"
license: "LicenseRef-Workspace-Owner"
sha256: "c2919dcb275e1224e9fd59c9ebf8300854de43a4196f7b12aebd9ca79fc817cb"
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
  <ItemGroup><ProjectReference Include="..\..\src\Elite.Arca.Credentials\Elite.Arca.Credentials.csproj" /></ItemGroup>
</Project>
````

### FILE: `arca/credentials/tests/Elite.Arca.Credentials.Tests/packages.lock.json`
```yaml
block_id: "ARCA-CREDENTIALS:test-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "NuGet generated transitive dependency lock"
license: "LicenseRef-Workspace-Owner"
sha256: "74bcd5d1fcdb61aafe47185107202d3e91391ce889876a5a3327950544c878cf"
variables: []
secrets_allowed: false
```
````json
{
  "version": 1,
  "dependencies": {
    "net10.0": {
      "System.Security.Cryptography.Pkcs": {
        "type": "Transitive",
        "resolved": "10.0.11",
        "contentHash": "8IV+rI3xN/Mkq9MsSX7VZTv9T9Wt+tyhHLwBX9VNZBk8m8kknEs1JeTSjI2x9M2L/fSja+c2WSE5n5Yy0GXdoQ=="
      },
      "elite.arca.credentials": {
        "type": "Project",
        "dependencies": {
          "System.Security.Cryptography.Pkcs": "[10.0.11, )"
        }
      }
    }
  }
}
````

### FILE: `arca/credentials/tests/Elite.Arca.Credentials.Tests/Program.cs`
```yaml
block_id: "ARCA-CREDENTIALS:tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local cryptographic, parser and concurrency regressions"
license: "LicenseRef-Workspace-Owner"
sha256: "917e63398272aada2621e42be6e453622ad949e17c9e230d4371b495ae9ffa0a"
variables: []
secrets_allowed: false
```
````csharp
using System.Security.Cryptography;
using System.Security.Cryptography.Pkcs;
using System.Security.Cryptography.X509Certificates;
using System.Text;
using System.Xml;
using Elite.Arca.Credentials;

var tests = new List<(string Name, Func<Task> Run)>
{
    ("options", () => { WsaaOptions.ElectronicInvoice.Validate(); Throws<ArgumentException>(() => new WsaaOptions("WS FE!", TimeSpan.FromHours(12), TimeSpan.Zero, TimeSpan.FromMinutes(1)).Validate()); return Task.CompletedTask; }),
    ("request", () => { var bytes = Factory(new FixedClock()).Create(); var xml = Encoding.UTF8.GetString(bytes); Check(xml.Contains("<service>wsfe</service>", StringComparison.Ordinal)); Check(xml.Contains("<uniqueId>7</uniqueId>", StringComparison.Ordinal)); return Task.CompletedTask; }),
    ("cms-attached", () => { using var cert = Certificate(); var request = Factory(new FixedClock()).Create(); var encoded = new CmsRequestSigner().Sign(cert, request); var cms = new SignedCms(); cms.Decode(Convert.FromBase64String(encoded)); cms.CheckSignature(verifySignatureOnly: true); Check(cms.ContentInfo.Content.SequenceEqual(request)); return Task.CompletedTask; }),
    ("thumbprint", () => { Check(CertificateStoreLoader.NormalizeThumbprint("AA bb-" + new string('c', 36)) == "AABB" + new string('C', 36)); return Task.CompletedTask; }),
    ("parser", () => { var now = FixedClock.Value; var ticket = WsaaLoginResponseParser.Parse(Response(now.AddHours(1)), now, TimeSpan.FromHours(13)); Check(ticket.Token == "token-secret" && ticket.ToString().Contains("REDACTED", StringComparison.Ordinal) && !ticket.ToString().Contains("token-secret", StringComparison.Ordinal)); return Task.CompletedTask; }),
    ("parser-dtd", () => { Throws<XmlException>(() => WsaaLoginResponseParser.Parse("<!DOCTYPE x [<!ENTITY e SYSTEM 'file:///x'>]><loginTicketResponse>&e;</loginTicketResponse>", FixedClock.Value, TimeSpan.FromHours(13))); return Task.CompletedTask; }),
    ("expired", () => { Throws<InvalidDataException>(() => WsaaLoginResponseParser.Parse(Response(FixedClock.Value.AddSeconds(-1)), FixedClock.Value, TimeSpan.FromHours(13))); return Task.CompletedTask; }),
    ("single-flight", SingleFlight),
    ("sign-rejects-no-key", () => { using var rsa = RSA.Create(2048); var req = new CertificateRequest("CN=NoKey", rsa, HashAlgorithmName.SHA256, RSASignaturePadding.Pkcs1); using var withKey = req.CreateSelfSigned(DateTimeOffset.UtcNow.AddDays(-1), DateTimeOffset.UtcNow.AddDays(1)); using var publicOnly = X509CertificateLoader.LoadCertificate(withKey.Export(X509ContentType.Cert)); Throws<CryptographicException>(() => new CmsRequestSigner().Sign(publicOnly, [1,2,3])); return Task.CompletedTask; })
};
foreach (var test in tests) { await test.Run(); Console.WriteLine($"PASS {test.Name}"); }
Console.WriteLine($"ARCA_WSAA_CREDENTIAL_TEST_PASS tests={tests.Count} secrets_logged=0");

static LoginTicketRequestFactory Factory(IUtcClock clock) => new(WsaaOptions.ElectronicInvoice, clock, new FixedIds());
static async Task SingleFlight()
{
    var clock = new FixedClock(); using var cert = Certificate(); var transport = new FakeTransport(clock);
    using var provider = new WsaaCredentialProvider(WsaaOptions.ElectronicInvoice, clock, Factory(clock), new CmsRequestSigner(), cert, transport);
    var results = await Task.WhenAll(Enumerable.Range(0, 32).Select(_ => provider.GetAsync()));
    Check(transport.Calls == 1); Check(results.All(x => ReferenceEquals(x, results[0]))); _ = await provider.GetAsync(); Check(transport.Calls == 1);
}
static X509Certificate2 Certificate() { using var rsa = RSA.Create(2048); var req = new CertificateRequest("CN=Elite ARCA Test", rsa, HashAlgorithmName.SHA256, RSASignaturePadding.Pkcs1); return req.CreateSelfSigned(DateTimeOffset.UtcNow.AddDays(-1), DateTimeOffset.UtcNow.AddDays(1)); }
static string Response(DateTimeOffset expiration) => $"<loginTicketResponse><header><expirationTime>{XmlConvert.ToString(expiration)}</expirationTime></header><credentials><token>token-secret</token><sign>sign-secret</sign></credentials></loginTicketResponse>";
static void Check(bool condition) { if (!condition) throw new InvalidOperationException("check failed"); }
static void Throws<T>(Action action) where T : Exception { try { action(); } catch (T) { return; } throw new InvalidOperationException($"expected {typeof(T).Name}"); }
sealed class FixedClock : IUtcClock { public static readonly DateTimeOffset Value = new(2026, 8, 30, 12, 0, 0, TimeSpan.Zero); public DateTimeOffset UtcNow => Value; }
sealed class FixedIds : ILoginTicketIdSource { public long Next() => 7; }
sealed class FakeTransport(FixedClock clock) : IWsaaTransport { public int Calls; public async Task<string> LoginCmsAsync(string cmsBase64, CancellationToken cancellationToken) { if (string.IsNullOrWhiteSpace(cmsBase64)) throw new InvalidOperationException("CMS missing"); Interlocked.Increment(ref Calls); await Task.Delay(20, cancellationToken); return $"<loginTicketResponse><header><expirationTime>{XmlConvert.ToString(clock.UtcNow.AddHours(12))}</expirationTime></header><credentials><token>token-secret</token><sign>sign-secret</sign></credentials></loginTicketResponse>"; } }
````

### FILE: `tools/test-arca-wsaa-credential-core.ps1`
```yaml
block_id: "ARCA-CREDENTIALS:runner:v1"
operation: CREATE
provenance: AUTHORED
source: "local locked restore/build/test/vulnerability runner"
license: "LicenseRef-Workspace-Owner"
sha256: "4ac7cefe2cf789c2abafb4d9c8b2f48bed3b25537fff295081d1ea5bedf973be"
variables: ["DotnetExecutable", "Root", "GeneratedClientRoot"]
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding()]
param(
  [Parameter(Mandatory=$true)][string]$DotnetExecutable,
  [Parameter(Mandatory=$true)][string]$Root,
  [Parameter(Mandatory=$true)][string]$GeneratedClientRoot
)
$ErrorActionPreference='Stop'
$dotnet=[IO.Path]::GetFullPath($DotnetExecutable);$base=[IO.Path]::GetFullPath($Root)
if(-not(Test-Path -LiteralPath $dotnet -PathType Leaf)){throw 'dotnet executable missing'}
function Run([string[]]$Arguments){$lines=@(& $dotnet @Arguments 2>&1);if($LASTEXITCODE-ne0){throw "dotnet failed ($LASTEXITCODE): $($Arguments-join' ')`n$($lines-join[Environment]::NewLine)"};$lines}
$version=(@(Run @('--version'))[-1]).Trim();if($version-ne'10.0.400'){throw "dotnet SDK mismatch: $version"}
$project=Join-Path $base 'arca\credentials\tests\Elite.Arca.Credentials.Tests\Elite.Arca.Credentials.Tests.csproj'
$integrationProject=Join-Path $base 'arca\credentials\integration-tests\Elite.Arca.Wsaa.Transport.Tests\Elite.Arca.Wsaa.Transport.Tests.csproj'
$generated=[IO.Path]::GetFullPath($GeneratedClientRoot);if(-not(Test-Path -LiteralPath (Join-Path $generated 'Generated\WsaaReference.cs') -PathType Leaf)){throw 'generated WSAA client missing'}
$generatedProperty="-p:GeneratedClientRoot=$generated"
$oldTelemetry=$env:DOTNET_CLI_TELEMETRY_OPTOUT
try{
  $env:DOTNET_CLI_TELEMETRY_OPTOUT='1'
  Run @('restore',$project,'--locked-mode','--warnaserror')|Out-Null
  Run @('build',$project,'--configuration','Release','--no-restore','--warnaserror')|Out-Null
  $tests=Run @('run','--project',$project,'--configuration','Release','--no-build')
  if(($tests-join"`n")-notmatch'ARCA_WSAA_CREDENTIAL_TEST_PASS tests=9 secrets_logged=0'){throw 'test proof missing'}
  Run @('restore',$integrationProject,$generatedProperty,'--locked-mode','--warnaserror')|Out-Null
  Run @('build',$integrationProject,$generatedProperty,'--configuration','Release','--no-restore','--warnaserror')|Out-Null
  $integrationTests=Run @('run','--project',$integrationProject,$generatedProperty,'--configuration','Release','--no-build')
  if(($integrationTests-join"`n")-notmatch'ARCA_GENERATED_WSAA_TRANSPORT_TEST_PASS tests=3 manual_glue=0'){throw 'generated transport proof missing'}
  $audit=Run @('list',$project,'package','--vulnerable','--include-transitive')
  if(($audit-join"`n")-notmatch'no vulnerable packages'){throw "NuGet vulnerability gate did not return clean`n$($audit-join[Environment]::NewLine)"}
}finally{$env:DOTNET_CLI_TELEMETRY_OPTOUT=$oldTelemetry}
'ARCA_WSAA_CREDENTIAL_CORE_PASS tests=12 build_warnings=0 vulnerable_packages=0 secrets_logged=0 manual_glue=0 production_admitted=false'
````

## 6. Configuration surface

Runtime inputs are service profile, external positive unique-ID source and exact certificate thumbprint/store location. Secret bytes and passwords are prohibited. The test runner accepts SDK path, materialized root and generated-client root.

## 7. Dependency bill

Dependencies are Microsoft `System.Security.Cryptography.Pkcs`/Xml 10.0.11 and WCF Http 10.0.652802, fixed by four lockfiles and audited through NuGet. XML, certificate store, concurrency and memory clearing use .NET platform APIs.

## 8. Apply order

Generate the WSAA/WSFE clients, materialize this core, compose the included mapper, provide a PostgreSQL-backed unique-ID source in the application, install the project-authorized certificate in its secret store, then run homologation. Never copy token/sign into configuration.

## 9. Verification

Require 18/18 reconstruction, four locked restores, zero-warning builds, nine credential tests plus three exact generated-interface tests, current NuGet vulnerability clean result, `secrets_logged=0` and `manual_glue=0`. Live login/homologation requires the project's ARCA test certificate and remains prohibited without authority.

## 10. Reconstruction evidence

Recorded in `reconstruction_evidence/MICROSOFT_ARCA_WSAA_CREDENTIAL_CORE_2026-08-30_V125.md`.
