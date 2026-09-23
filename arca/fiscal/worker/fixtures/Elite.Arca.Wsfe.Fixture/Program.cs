// AUTHORED local fixture composition. No live transport or user certificate.
using System.Security.Cryptography;
using System.Security.Cryptography.Pkcs;
using System.Security.Cryptography.X509Certificates;
using System.Xml;
using Elite.Arca.Credentials;
using Elite.Arca.Wsfe.Bridge;
using Elite.Arca.Wsfe.Worker;
using Generated = Elite.Arca.Wsfe.Generated;

if (args.Length != 1 || Environment.GetEnvironmentVariable("ELITE_ARCA_FIXTURE") != "1")
    throw new InvalidOperationException("Explicit local fixture mode and one Unix socket path required.");
using var rsa = RSA.Create(2048);
var request = new CertificateRequest("CN=Elite Local ARCA Fixture Only", rsa, HashAlgorithmName.SHA256, RSASignaturePadding.Pkcs1);
using var certificate = request.CreateSelfSigned(DateTimeOffset.UtcNow.AddMinutes(-1), DateTimeOffset.UtcNow.AddHours(2));
var wsaa = new FixtureWsaa();
using var credentials = new WsaaCredentialProvider(WsaaOptions.ElectronicInvoice, new SystemUtcClock(),
    new LoginTicketRequestFactory(WsaaOptions.ElectronicInvoice, new SystemUtcClock(), new MonotonicLoginTicketIdSource()),
    new CmsRequestSigner(), certificate, wsaa);
var access = new CredentialAccessSource(credentials);
var soap = new FixtureSoap();
var operations = new BridgeFiscalOperations(new WsfeBridge(access, soap), new WsfeParameterBridge(access, new FakeParameters()));
await using var app = WorkerHost.Build([], args[0], operations);
app.MapGet("/fixture/stats", () => new { scope = "LOCAL_FIXTURES", last = soap.LastCalls, authorize = soap.AuthorizeCalls, consult = soap.ConsultCalls, cms = wsaa.Calls });
app.MapPost("/fixture/stop", (IHostApplicationLifetime lifetime) => { lifetime.StopApplication(); return Results.Ok(new { stopping = true }); });
await app.StartAsync();
WorkerHost.SecureSocket(args[0]);
if (app.Urls.Count == 0 || app.Urls.Any(value => !value.Contains("unix:", StringComparison.OrdinalIgnoreCase) && !value.Contains(args[0], StringComparison.OrdinalIgnoreCase)))
    throw new InvalidOperationException("Fixture must expose only the Unix socket.");
Console.WriteLine("ARCA_CONNECTED_FIXTURE_READY uds_only=true live_transport=false");
await app.WaitForShutdownAsync();

sealed class FixtureWsaa : IWsaaTransport
{
    public int Calls;
    public Task<string> LoginCmsAsync(string cmsBase64, CancellationToken cancellationToken)
    {
        cancellationToken.ThrowIfCancellationRequested();
        var cms = new SignedCms(); cms.Decode(Convert.FromBase64String(cmsBase64)); cms.CheckSignature(verifySignatureOnly: true);
        if (cms.ContentInfo.Content.Length == 0) throw new InvalidDataException("Empty signed fixture request.");
        Interlocked.Increment(ref Calls);
        var expiration = XmlConvert.ToString(DateTimeOffset.UtcNow.AddHours(1));
        return Task.FromResult($"<loginTicketResponse><header><expirationTime>{expiration}</expirationTime></header><credentials><token>public-fixture-token</token><sign>public-fixture-sign</sign></credentials></loginTicketResponse>");
    }
}

sealed class FixtureSoap : IWsfeSoap
{
    private readonly object sync = new();
    private readonly Dictionary<(long Cuit, int Point, int Type, long Number), Generated.FECompConsResponse> issued = [];
    public int LastCalls, AuthorizeCalls, ConsultCalls;
    private static void Auth(Generated.FEAuthRequest auth)
    {
        if (auth.Token != "public-fixture-token" || auth.Sign != "public-fixture-sign") throw new InvalidDataException("Fixture credential mapping failed.");
    }
    public Task<Generated.FERecuperaLastCbteResponse> LastAuthorizedAsync(Generated.FEAuthRequest auth, int pointOfSale, int voucherType, CancellationToken cancellationToken)
    {
        Auth(auth); Interlocked.Increment(ref LastCalls);
        lock (sync) return Task.FromResult(new Generated.FERecuperaLastCbteResponse { PtoVta = pointOfSale, CbteTipo = voucherType, CbteNro = checked((int)issued.Keys.Where(k => k.Cuit == auth.Cuit && k.Point == pointOfSale && k.Type == voucherType).Select(k => k.Number).DefaultIfEmpty(0).Max()) });
    }
    public Task<Generated.FECompConsultaResponse> ConsultAsync(Generated.FEAuthRequest auth, Generated.FECompConsultaReq request, CancellationToken cancellationToken)
    {
        Auth(auth); Interlocked.Increment(ref ConsultCalls);
        lock (sync) return Task.FromResult(issued.TryGetValue((auth.Cuit, request.PtoVta, request.CbteTipo, request.CbteNro), out var result)
            ? new Generated.FECompConsultaResponse { ResultGet = result }
            : new Generated.FECompConsultaResponse { Errors = [new Generated.Err { Code = 602, Msg = "fixture absent" }] });
    }
    public Task<Generated.FECAEResponse> AuthorizeAsync(Generated.FEAuthRequest auth, Generated.FECAERequest request, CancellationToken cancellationToken)
    {
        Auth(auth); var count = Interlocked.Increment(ref AuthorizeCalls);
        var detail = request.FeDetReq.Single(); var header = request.FeCabReq;
        if (detail.ImpTotal != 121 || detail.ImpNeto != 100 || detail.ImpIVA != 21 || !detail.CondicionIVAReceptorIdSpecified || detail.MonId != "PES")
            throw new InvalidDataException("Exact fixture invoice mapping failed.");
        lock (sync)
        {
            if (header.CbteTipo == 8 && (detail.CbtesAsoc is not { Length: 1 } || !issued.ContainsKey((auth.Cuit, detail.CbtesAsoc[0].PtoVta, detail.CbtesAsoc[0].Tipo, detail.CbtesAsoc[0].Nro))))
                throw new InvalidDataException("Credit fixture must reference issued original.");
            issued.Add((auth.Cuit, header.PtoVta, header.CbteTipo, detail.CbteDesde), new Generated.FECompConsResponse { PtoVta = header.PtoVta, CbteTipo = header.CbteTipo, CbteDesde = detail.CbteDesde, CbteHasta = detail.CbteHasta, Resultado = "A", CodAutorizacion = "41124578989845", FchVto = "20300913" });
        }
        // Provider commits before losing its first response. Retry must consult.
        if (count == 1) throw new IOException("Deliberate fixture response loss after authorization.");
        return Task.FromResult(new Generated.FECAEResponse { FeCabResp = new Generated.FECAECabResponse { PtoVta = header.PtoVta, CbteTipo = header.CbteTipo, CantReg = 1, Resultado = "A" }, FeDetResp = [new Generated.FECAEDetResponse { CbteDesde = detail.CbteDesde, CbteHasta = detail.CbteHasta, Resultado = "A", CAE = "41124578989845", CAEFchVto = "20300913" }] });
    }
}
