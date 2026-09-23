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
