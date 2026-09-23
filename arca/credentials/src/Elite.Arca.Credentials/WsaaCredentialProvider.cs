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
