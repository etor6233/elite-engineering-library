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
