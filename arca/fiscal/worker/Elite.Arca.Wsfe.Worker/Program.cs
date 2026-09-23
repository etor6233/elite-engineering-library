using System.Security.Cryptography.X509Certificates;
using Elite.Arca.Credentials;
using Elite.Arca.Wsfe.Bridge;

namespace Elite.Arca.Wsfe.Worker;

public static class Program
{
    private static readonly Uri WsaaHomologation = new("https://wsaahomo.afip.gov.ar/ws/services/LoginCms");
    private static readonly Uri WsfeHomologation = new("https://wswhomo.afip.gov.ar/wsfev1/service.asmx");

    public static async Task<int> Main(string[] args)
    {
        try
        {
            if (!string.Equals(Environment.GetEnvironmentVariable("ARCA_ENVIRONMENT"), "homologation", StringComparison.Ordinal))
                throw new InvalidOperationException("Only the ARCA homologation environment is admitted by this worker.");
            var socketPath = Required("ARCA_WSFE_SOCKET");
            var thumbprint = Required("ARCA_CERTIFICATE_THUMBPRINT");
            var storeLocation = Required("ARCA_CERTIFICATE_STORE") switch
            {
                "CurrentUser" => StoreLocation.CurrentUser,
                "LocalMachine" => StoreLocation.LocalMachine,
                _ => throw new InvalidOperationException("ARCA_CERTIFICATE_STORE must be CurrentUser or LocalMachine.")
            };
            using var certificate = CertificateStoreLoader.LoadByThumbprint(thumbprint, storeLocation);
            using var credentials = new WsaaCredentialProvider(
                WsaaOptions.ElectronicInvoice,
                new SystemUtcClock(),
                new LoginTicketRequestFactory(WsaaOptions.ElectronicInvoice, new SystemUtcClock(), new MonotonicLoginTicketIdSource()),
                new CmsRequestSigner(),
                certificate,
                new WsaaSoapTransport(WsaaHomologation));
            var access = new CredentialAccessSource(credentials);
            var transport = new WsfeSoapTransport(WsfeHomologation);
            var operations = new BridgeFiscalOperations(new WsfeBridge(access, transport), new WsfeParameterBridge(access, transport));
            await using var app = WorkerHost.Build(args, socketPath, operations);
            await app.StartAsync().ConfigureAwait(false);
            WorkerHost.SecureSocket(socketPath);
            await app.WaitForShutdownAsync().ConfigureAwait(false);
            return 0;
        }
        catch
        {
            Console.Error.WriteLine("ARCA WSFE worker startup or runtime failed; inspect controlled platform diagnostics.");
            return 1;
        }
    }

    private static string Required(string name) => Environment.GetEnvironmentVariable(name) is { Length: > 0 } value ? value : throw new InvalidOperationException($"{name} is required.");
}
