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
