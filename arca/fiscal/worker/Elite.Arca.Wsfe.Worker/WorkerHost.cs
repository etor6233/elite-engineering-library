using System.Text.Json.Serialization;
using Elite.Arca.Wsfe.Bridge;
using Microsoft.AspNetCore.Http.Json;

namespace Elite.Arca.Wsfe.Worker;

public interface IFiscalOperations
{
    Task<LastAuthorizedResult> LastAuthorizedAsync(FiscalInvoice invoice, CancellationToken cancellationToken);
    Task<SafeWsfeResult> ConsultAsync(FiscalInvoice invoice, CancellationToken cancellationToken);
    Task<SafeWsfeResult> AuthorizeAsync(FiscalInvoice invoice, CancellationToken cancellationToken);
    Task<SafeParameterResult> ParametersAsync(ParameterRequest request, CancellationToken cancellationToken);
}

public sealed record ParameterRequest(string TaxpayerCuit, string Kind, string? VoucherClass);

public sealed class BridgeFiscalOperations(WsfeBridge bridge, WsfeParameterBridge parameters) : IFiscalOperations
{
    public Task<LastAuthorizedResult> LastAuthorizedAsync(FiscalInvoice invoice, CancellationToken cancellationToken) => bridge.LastAuthorizedAsync(invoice, cancellationToken);
    public Task<SafeWsfeResult> ConsultAsync(FiscalInvoice invoice, CancellationToken cancellationToken) => bridge.ConsultAsync(invoice, cancellationToken);
    public Task<SafeWsfeResult> AuthorizeAsync(FiscalInvoice invoice, CancellationToken cancellationToken) => bridge.AuthorizeAsync(invoice, cancellationToken);
    public Task<SafeParameterResult> ParametersAsync(ParameterRequest request, CancellationToken cancellationToken) => request.Kind switch
    {
        "voucher_type" when request.VoucherClass is null => parameters.VoucherTypesAsync(request.TaxpayerCuit, cancellationToken),
        "concept" when request.VoucherClass is null => parameters.ConceptsAsync(request.TaxpayerCuit, cancellationToken),
        "document_type" when request.VoucherClass is null => parameters.DocumentTypesAsync(request.TaxpayerCuit, cancellationToken),
        "vat_rate" when request.VoucherClass is null => parameters.VatRatesAsync(request.TaxpayerCuit, cancellationToken),
        "other_tax" when request.VoucherClass is null => parameters.OtherTaxesAsync(request.TaxpayerCuit, cancellationToken),
        "point_of_sale" when request.VoucherClass is null => parameters.PointsOfSaleAsync(request.TaxpayerCuit, cancellationToken),
        "recipient_vat_condition" when request.VoucherClass is not null => parameters.RecipientVatConditionsAsync(request.TaxpayerCuit, request.VoucherClass, cancellationToken),
        _ => throw new ArgumentException("Unsupported parameter request.", nameof(request))
    };
}

public static class WorkerHost
{
    public static WebApplication Build(string[] args, string socketPath, IFiscalOperations operations)
    {
        ArgumentNullException.ThrowIfNull(operations);
        ValidateSocket(socketPath);
        var builder = WebApplication.CreateSlimBuilder(args);
        builder.Logging.ClearProviders();
        builder.Services.Configure<JsonOptions>(options =>
        {
            options.SerializerOptions.PropertyNamingPolicy = System.Text.Json.JsonNamingPolicy.SnakeCaseLower;
            options.SerializerOptions.DictionaryKeyPolicy = System.Text.Json.JsonNamingPolicy.SnakeCaseLower;
            options.SerializerOptions.UnmappedMemberHandling = JsonUnmappedMemberHandling.Disallow;
        });
        builder.WebHost.ConfigureKestrel(options =>
        {
            options.AddServerHeader = false;
            options.Limits.MaxRequestBodySize = 64 << 10;
            options.Limits.MaxRequestHeaderCount = 32;
            options.Limits.MaxRequestHeadersTotalSize = 16 << 10;
            options.Limits.RequestHeadersTimeout = TimeSpan.FromSeconds(5);
            options.Limits.KeepAliveTimeout = TimeSpan.FromSeconds(30);
            options.ListenUnixSocket(socketPath);
        });
        var app = builder.Build();
        app.MapGet("/health/live", static () => Results.Ok(new { status = "alive" }));
        app.MapPost("/v1/last-authorized", (FiscalInvoice invoice, CancellationToken token) => Execute(() => operations.LastAuthorizedAsync(invoice, token)));
        app.MapPost("/v1/consult", (FiscalInvoice invoice, CancellationToken token) => Execute(() => operations.ConsultAsync(invoice, token)));
        app.MapPost("/v1/authorize", (FiscalInvoice invoice, CancellationToken token) => Execute(() => operations.AuthorizeAsync(invoice, token)));
        app.MapPost("/v1/parameters", (ParameterRequest request, CancellationToken token) => Execute(() => operations.ParametersAsync(request, token)));
        return app;
    }

    public static void SecureSocket(string socketPath)
    {
        if (!Path.Exists(socketPath)) throw new IOException("Kestrel did not create the Unix socket.");
        if (!OperatingSystem.IsWindows())
            File.SetUnixFileMode(socketPath, UnixFileMode.UserRead | UnixFileMode.UserWrite | UnixFileMode.GroupRead | UnixFileMode.GroupWrite);
    }

    private static async Task<IResult> Execute<T>(Func<Task<T>> operation)
    {
        try { return Results.Ok(await operation().ConfigureAwait(false)); }
        catch (WsfeResponseException error)
        {
            return Results.Json(new SafeFailure("provider_rejected", error.ResponseHash, error.ProviderCodes), statusCode: StatusCodes.Status502BadGateway);
        }
        catch (ArgumentException) { return Results.Json(new SafeFailure("invalid_request", null, []), statusCode: StatusCodes.Status400BadRequest); }
        catch (OperationCanceledException) { return Results.Json(new SafeFailure("cancelled", null, []), statusCode: 499); }
        catch { return Results.Json(new SafeFailure("provider_unavailable", null, []), statusCode: StatusCodes.Status503ServiceUnavailable); }
    }

    private static void ValidateSocket(string socketPath)
    {
        if (string.IsNullOrWhiteSpace(socketPath) || !Path.IsPathFullyQualified(socketPath)) throw new ArgumentException("An absolute Unix socket path is required.", nameof(socketPath));
        var parent = Path.GetDirectoryName(socketPath);
        if (parent is null || !Directory.Exists(parent)) throw new DirectoryNotFoundException("The Unix socket parent directory must be provisioned before startup.");
        if (Path.Exists(socketPath)) throw new IOException("The Unix socket path already exists; stale sockets must be investigated before startup.");
        if ((File.GetAttributes(parent) & FileAttributes.ReparsePoint) != 0) throw new IOException("The Unix socket parent directory cannot be a reparse point.");
    }

    private sealed record SafeFailure(string Code, string? ResponseHash, IReadOnlyList<string> ProviderCodes);
}
