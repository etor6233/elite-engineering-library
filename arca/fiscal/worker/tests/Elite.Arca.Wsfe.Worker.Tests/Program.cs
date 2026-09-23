using System.Net;
using System.Net.Http.Json;
using System.Net.Sockets;
using System.Text;
using Elite.Arca.Wsfe.Bridge;
using Elite.Arca.Wsfe.Worker;

const string InvoiceJson = """
{"taxpayer_cuit":"33693450239","point_of_sale":12,"voucher_type":1,"concept":1,"recipient_document_type":80,"recipient_document":"20111111112","recipient_vat_condition_id":1,"currency":"ARS","voucher_number":41,"total_minor_units":18405,"net_minor_units":15000,"vat_minor_units":2625,"exempt_minor_units":0,"non_taxed_minor_units":0,"other_tax_minor_units":780,"issued_on":"2010-09-03","vat_lines":[{"id":5,"base_minor_units":10000,"amount_minor_units":2100},{"id":4,"base_minor_units":5000,"amount_minor_units":525}],"other_tax_lines":[{"id":99,"description":"Impuesto Municipal Matanza","base_minor_units":15000,"rate_basis_points":520,"amount_minor_units":780}]}
""";

var root = Path.Combine(Path.GetTempPath(), $"elite-wsfe-worker-{Environment.ProcessId}");
Directory.CreateDirectory(root);
var socketPath = Path.Combine(root, "worker.sock");
var operations = new FakeOperations();
await using var app = WorkerHost.Build([], socketPath, operations);
try
{
    await app.StartAsync();
    WorkerHost.SecureSocket(socketPath);
    Assert(Path.Exists(socketPath), "Kestrel created Unix socket");
    Assert(app.Urls.Count > 0 && app.Urls.All(value => value.Contains("unix:", StringComparison.OrdinalIgnoreCase) || value.Contains(socketPath, StringComparison.OrdinalIgnoreCase)), "only the Unix socket address was configured");
    using var handler = new SocketsHttpHandler
    {
        ConnectCallback = async (_, cancellationToken) =>
        {
            var socket = new Socket(AddressFamily.Unix, SocketType.Stream, ProtocolType.Unspecified);
            try
            {
                await socket.ConnectAsync(new UnixDomainSocketEndPoint(socketPath), cancellationToken);
                return new NetworkStream(socket, ownsSocket: true);
            }
            catch { socket.Dispose(); throw; }
        }
    };
    using var client = new HttpClient(handler) { BaseAddress = new Uri("http://arca-wsfe") };

    var health = await client.GetStringAsync("/health/live");
    Assert(health == "{\"status\":\"alive\"}", "liveness response is narrow");

    using var last = await client.PostAsync("/v1/last-authorized", Json(InvoiceJson));
    Assert(last.StatusCode == HttpStatusCode.OK, "last-authorized status");
    Assert((await last.Content.ReadAsStringAsync()).Contains("\"number\":40", StringComparison.Ordinal), "last-authorized result");
    Assert(operations.LastInvoice is { RecipientVatConditionId: 1, VatLines.Count: 2, OtherTaxLines.Count: 1 }, "exact tax details crossed IPC");

    using var authorized = await client.PostAsync("/v1/authorize", Json(InvoiceJson));
    var authorizedBody = await authorized.Content.ReadAsStringAsync();
    Assert(authorized.StatusCode == HttpStatusCode.OK && authorizedBody.Contains("\"cae_expires_on\":\"2010-09-13\"", StringComparison.Ordinal), "authorization result uses strict snake case");

    using var parameters = await client.PostAsync("/v1/parameters", Json("""{"taxpayer_cuit":"33693450239","kind":"vat_rate","voucher_class":null}"""));
    var parameterBody = await parameters.Content.ReadAsStringAsync();
    Assert(parameters.StatusCode == HttpStatusCode.OK && parameterBody.Contains("\"kind\":\"vat_rate\"", StringComparison.Ordinal) && parameterBody.Contains("\"code\":\"5\"", StringComparison.Ordinal), "parameter snapshot crosses UDS");
    Assert(!parameterBody.Contains("secret", StringComparison.OrdinalIgnoreCase) && !parameterBody.Contains("message", StringComparison.OrdinalIgnoreCase), "parameter response is safe");

    using var invalidParameters = await client.PostAsync("/v1/parameters", Json("""{"taxpayer_cuit":"33693450239","kind":"invented","voucher_class":null}"""));
    Assert(invalidParameters.StatusCode == HttpStatusCode.BadRequest, "unknown parameter kind fails closed");

    using var unknown = await client.PostAsync("/v1/authorize", Json(InvoiceJson.Replace("\"taxpayer_cuit\"", "\"unknown\":true,\"taxpayer_cuit\"", StringComparison.Ordinal)));
    Assert(unknown.StatusCode == HttpStatusCode.BadRequest, "unknown request field fails closed");

    operations.FailConsult = true;
    using var failed = await client.PostAsync("/v1/consult", Json(InvoiceJson));
    var failedBody = await failed.Content.ReadAsStringAsync();
    Assert(failed.StatusCode == HttpStatusCode.BadGateway && failedBody.Contains("\"code\":\"provider_rejected\"", StringComparison.Ordinal), "provider failure is classified");
    Assert(!failedBody.Contains("secret-token-sign", StringComparison.Ordinal) && !failedBody.Contains("message", StringComparison.OrdinalIgnoreCase), "provider message and credentials are absent");

    using var oversized = await client.PostAsync("/v1/authorize", Json("{\"padding\":\"" + new string('x', 70_000) + "\"}"));
    Assert(oversized.StatusCode == HttpStatusCode.RequestEntityTooLarge, "oversized request is rejected by Kestrel");
    Console.WriteLine("ARCA_WSFE_WORKER_TESTS_PASS tests=12");
}
finally
{
    await app.StopAsync();
    if (Path.Exists(socketPath)) File.Delete(socketPath);
    Directory.Delete(root);
}

static StringContent Json(string value) => new(value, Encoding.UTF8, "application/json");
static void Assert(bool condition, string name)
{
    if (!condition) throw new InvalidOperationException($"FAILED: {name}");
}

sealed class FakeOperations : IFiscalOperations
{
    internal FiscalInvoice? LastInvoice { get; private set; }
    internal bool FailConsult { get; set; }
    public Task<LastAuthorizedResult> LastAuthorizedAsync(FiscalInvoice invoice, CancellationToken cancellationToken)
    {
        LastInvoice = invoice;
        return Task.FromResult(new LastAuthorizedResult(40, new string('a', 64), []));
    }
    public Task<SafeWsfeResult> ConsultAsync(FiscalInvoice invoice, CancellationToken cancellationToken)
    {
        if (FailConsult) throw new WsfeResponseException("secret-token-sign provider message", new string('b', 64), ["E:100"]);
        return Task.FromResult(new SafeWsfeResult(false, false, null, null, new string('b', 64), ["E:602"]));
    }
    public Task<SafeWsfeResult> AuthorizeAsync(FiscalInvoice invoice, CancellationToken cancellationToken) =>
        Task.FromResult(new SafeWsfeResult(true, true, "41124578989845", new DateOnly(2010, 9, 13), new string('c', 64), []));
    public Task<SafeParameterResult> ParametersAsync(ParameterRequest request, CancellationToken cancellationToken)
    {
        if (request.Kind != "vat_rate" || request.VoucherClass is not null) throw new ArgumentException("unsupported parameter fixture");
        return Task.FromResult(new SafeParameterResult(request.Kind, request.VoucherClass, [new SafeParameterItem("5", "21%", new DateOnly(2009, 2, 20), null, null, null, null, null)], new string('d', 64), []));
    }
}
