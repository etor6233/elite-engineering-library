using System.Security.Cryptography;
using System.ServiceModel;
using System.ServiceModel.Channels;
using Elite.Arca.Credentials;
using Elite.Arca.Wsfe.Bridge;
using WsaaGenerated = Elite.Arca.Wsaa.Generated;
using WsfeGenerated = Elite.Arca.Wsfe.Generated;

namespace Elite.Arca.Wsfe.Worker;

public sealed class MonotonicLoginTicketIdSource : ILoginTicketIdSource
{
    private long current = RandomNumberGenerator.GetInt32(1, int.MaxValue - 1_000_000);

    public long Next()
    {
        var value = Interlocked.Increment(ref current);
        if (value >= int.MaxValue) throw new InvalidOperationException("Login ticket ID space exhausted; restart is required.");
        return value;
    }
}

public sealed class WsaaSoapTransport(Uri endpoint) : IWsaaTransport
{
    public async Task<string> LoginCmsAsync(string cmsBase64, CancellationToken cancellationToken)
    {
        var client = new WsaaGenerated.LoginCMSClient(SoapBinding(), new EndpointAddress(endpoint));
        try
        {
            var response = await client.loginCmsAsync(cmsBase64).WaitAsync(cancellationToken).ConfigureAwait(false);
            if (string.IsNullOrWhiteSpace(response.loginCmsReturn)) throw new InvalidDataException("WSAA returned an empty response.");
            Close(client);
            return response.loginCmsReturn;
        }
        catch { Abort(client); throw; }
    }

    internal static BasicHttpBinding SoapBinding() => new(BasicHttpSecurityMode.Transport)
    {
        OpenTimeout = TimeSpan.FromSeconds(10),
        CloseTimeout = TimeSpan.FromSeconds(5),
        SendTimeout = TimeSpan.FromSeconds(30),
        ReceiveTimeout = TimeSpan.FromSeconds(30),
        MaxBufferSize = 1 << 20,
        MaxReceivedMessageSize = 1 << 20,
        TransferMode = TransferMode.Buffered,
        ReaderQuotas = new System.Xml.XmlDictionaryReaderQuotas
        {
            MaxDepth = 32,
            MaxStringContentLength = 1 << 20,
            MaxArrayLength = 65_536,
            MaxBytesPerRead = 4_096,
            MaxNameTableCharCount = 16_384
        }
    };

    internal static void Close(ICommunicationObject client)
    {
        try
        {
            if (client.State == CommunicationState.Faulted) client.Abort();
            else client.Close();
        }
        catch { client.Abort(); }
    }

    internal static void Abort(ICommunicationObject client)
    {
        try { client.Abort(); } catch { }
    }
}

public sealed class WsfeSoapTransport(Uri endpoint) : IWsfeSoap, IWsfeParameterSoap
{
    public Task<WsfeGenerated.FERecuperaLastCbteResponse> LastAuthorizedAsync(WsfeGenerated.FEAuthRequest auth, int pointOfSale, int voucherType, CancellationToken cancellationToken) =>
        Invoke(client => client.FECompUltimoAutorizadoAsync(auth, pointOfSale, voucherType), cancellationToken);

    public Task<WsfeGenerated.FECompConsultaResponse> ConsultAsync(WsfeGenerated.FEAuthRequest auth, WsfeGenerated.FECompConsultaReq request, CancellationToken cancellationToken) =>
        Invoke(client => client.FECompConsultarAsync(auth, request), cancellationToken);

    public Task<WsfeGenerated.FECAEResponse> AuthorizeAsync(WsfeGenerated.FEAuthRequest auth, WsfeGenerated.FECAERequest request, CancellationToken cancellationToken) =>
        Invoke(client => client.FECAESolicitarAsync(auth, request), cancellationToken);

    public Task<WsfeGenerated.CbteTipoResponse> VoucherTypesAsync(WsfeGenerated.FEAuthRequest auth, CancellationToken cancellationToken) =>
        Invoke(client => client.FEParamGetTiposCbteAsync(auth), cancellationToken);
    public Task<WsfeGenerated.ConceptoTipoResponse> ConceptsAsync(WsfeGenerated.FEAuthRequest auth, CancellationToken cancellationToken) =>
        Invoke(client => client.FEParamGetTiposConceptoAsync(auth), cancellationToken);
    public Task<WsfeGenerated.DocTipoResponse> DocumentTypesAsync(WsfeGenerated.FEAuthRequest auth, CancellationToken cancellationToken) =>
        Invoke(client => client.FEParamGetTiposDocAsync(auth), cancellationToken);
    public Task<WsfeGenerated.IvaTipoResponse> VatRatesAsync(WsfeGenerated.FEAuthRequest auth, CancellationToken cancellationToken) =>
        Invoke(client => client.FEParamGetTiposIvaAsync(auth), cancellationToken);
    public Task<WsfeGenerated.FETributoResponse> OtherTaxesAsync(WsfeGenerated.FEAuthRequest auth, CancellationToken cancellationToken) =>
        Invoke(client => client.FEParamGetTiposTributosAsync(auth), cancellationToken);
    public Task<WsfeGenerated.FEPtoVentaResponse> PointsOfSaleAsync(WsfeGenerated.FEAuthRequest auth, CancellationToken cancellationToken) =>
        Invoke(client => client.FEParamGetPtosVentaAsync(auth), cancellationToken);
    public Task<WsfeGenerated.CondicionIvaReceptorResponse> RecipientVatConditionsAsync(WsfeGenerated.FEAuthRequest auth, string voucherClass, CancellationToken cancellationToken) =>
        Invoke(client => client.FEParamGetCondicionIvaReceptorAsync(auth, voucherClass), cancellationToken);

    private async Task<T> Invoke<T>(Func<WsfeGenerated.ServiceSoapClient, Task<T>> call, CancellationToken cancellationToken)
    {
        var client = new WsfeGenerated.ServiceSoapClient(WsaaSoapTransport.SoapBinding(), new EndpointAddress(endpoint));
        try
        {
            var response = await call(client).WaitAsync(cancellationToken).ConfigureAwait(false);
            WsaaSoapTransport.Close(client);
            return response;
        }
        catch { WsaaSoapTransport.Abort(client); throw; }
    }
}
