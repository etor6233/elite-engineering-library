using System.Security.Cryptography;
using System.Security.Cryptography.X509Certificates;

namespace Elite.Arca.Credentials;

public static class CertificateStoreLoader
{
    public static X509Certificate2 LoadByThumbprint(string thumbprint, StoreLocation location = StoreLocation.CurrentUser)
    {
        var normalized = NormalizeThumbprint(thumbprint);
        using var store = new X509Store(StoreName.My, location, OpenFlags.ReadOnly | OpenFlags.OpenExistingOnly);
        var matches = store.Certificates.Find(X509FindType.FindByThumbprint, normalized, validOnly: false)
            .Where(c => string.Equals(NormalizeThumbprint(c.Thumbprint), normalized, StringComparison.Ordinal))
            .ToArray();
        if (matches.Length != 1)
        {
            foreach (var match in matches) match.Dispose();
            throw new CryptographicException("Certificate thumbprint must resolve to exactly one certificate.");
        }
        try { Validate(matches[0]); return matches[0]; }
        catch { matches[0].Dispose(); throw; }
    }

    public static string NormalizeThumbprint(string value)
    {
        ArgumentException.ThrowIfNullOrWhiteSpace(value);
        var normalized = string.Concat(value.Where(Uri.IsHexDigit)).ToUpperInvariant();
        if (normalized.Length is not (40 or 64)) throw new ArgumentException("Certificate thumbprint must be SHA-1 or SHA-256 hex.", nameof(value));
        return normalized;
    }

    public static void Validate(X509Certificate2 certificate)
    {
        if (!certificate.HasPrivateKey) throw new CryptographicException("Certificate has no private key.");
        var now = DateTimeOffset.UtcNow;
        if (now < certificate.NotBefore || now >= certificate.NotAfter) throw new CryptographicException("Certificate is outside its validity interval.");
        using var rsa = certificate.GetRSAPrivateKey();
        if (rsa is null || rsa.KeySize < 2048) throw new CryptographicException("An RSA private key of at least 2048 bits is required.");
    }
}
