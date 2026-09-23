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

foreach (var field in new[] { "VatLines", "OtherTaxLines", "AssociatedVouchers" })
{
    var malformed = System.Text.Json.Nodes.JsonNode.Parse(System.Text.Json.JsonSerializer.Serialize(invoice))!;
    malformed[field] = null;
    var nullCollections = System.Text.Json.JsonSerializer.Deserialize<FiscalInvoice>(malformed.ToJsonString())!;
    await ThrowsAsync<ArgumentException>(() => bridge.AuthorizeAsync(nullCollections), "null collection " + field);
}
Console.WriteLine("ARCA_WSFE_BRIDGE_TEST_PASS cases=18 secrets_returned=0 production_admitted=false");

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
