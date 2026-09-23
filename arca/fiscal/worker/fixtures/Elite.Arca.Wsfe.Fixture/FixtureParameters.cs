using Elite.Arca.Wsfe.Bridge;
using Generated = Elite.Arca.Wsfe.Generated;

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
