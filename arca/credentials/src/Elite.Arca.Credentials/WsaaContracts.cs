using System.Security.Cryptography.X509Certificates;
using System.Text.RegularExpressions;

namespace Elite.Arca.Credentials;

public sealed record WsaaOptions(string Service, TimeSpan TicketLifetime, TimeSpan GenerationSkew, TimeSpan RefreshBefore)
{
    private static readonly Regex ServicePattern = new("^[a-z0-9_]{1,32}$", RegexOptions.CultureInvariant | RegexOptions.NonBacktracking);
    public static WsaaOptions ElectronicInvoice { get; } = new("wsfe", TimeSpan.FromHours(12), TimeSpan.FromMinutes(5), TimeSpan.FromMinutes(10));

    public void Validate()
    {
        if (!ServicePattern.IsMatch(Service)) throw new ArgumentException("Invalid ARCA service identifier.", nameof(Service));
        if (TicketLifetime <= TimeSpan.Zero || TicketLifetime > TimeSpan.FromHours(12)) throw new ArgumentOutOfRangeException(nameof(TicketLifetime));
        if (GenerationSkew < TimeSpan.Zero || GenerationSkew > TimeSpan.FromMinutes(10)) throw new ArgumentOutOfRangeException(nameof(GenerationSkew));
        if (RefreshBefore <= TimeSpan.Zero || RefreshBefore >= TicketLifetime) throw new ArgumentOutOfRangeException(nameof(RefreshBefore));
    }
}

public interface IUtcClock { DateTimeOffset UtcNow { get; } }
public interface ILoginTicketIdSource { long Next(); }
public interface IWsaaTransport { Task<string> LoginCmsAsync(string cmsBase64, CancellationToken cancellationToken); }
public interface ICmsRequestSigner { string Sign(X509Certificate2 certificate, ReadOnlySpan<byte> request); }

public sealed class LoginTicketAccess
{
    public LoginTicketAccess(string token, string sign, DateTimeOffset expiresAt)
    {
        Token = string.IsNullOrWhiteSpace(token) ? throw new ArgumentException("Token is required.", nameof(token)) : token;
        Sign = string.IsNullOrWhiteSpace(sign) ? throw new ArgumentException("Sign is required.", nameof(sign)) : sign;
        ExpiresAt = expiresAt;
    }
    public string Token { get; }
    public string Sign { get; }
    public DateTimeOffset ExpiresAt { get; }
    public override string ToString() => "LoginTicketAccess [REDACTED]";
}

public sealed class SystemUtcClock : IUtcClock { public DateTimeOffset UtcNow => DateTimeOffset.UtcNow; }
