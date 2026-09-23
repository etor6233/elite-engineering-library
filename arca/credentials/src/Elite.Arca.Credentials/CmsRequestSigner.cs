using System.Security.Cryptography;
using System.Security.Cryptography.Pkcs;
using System.Security.Cryptography.X509Certificates;

namespace Elite.Arca.Credentials;

public sealed class CmsRequestSigner : ICmsRequestSigner
{
    public string Sign(X509Certificate2 certificate, ReadOnlySpan<byte> request)
    {
        ArgumentNullException.ThrowIfNull(certificate);
        if (!certificate.HasPrivateKey) throw new CryptographicException("Certificate has no private key.");
        var now = DateTimeOffset.UtcNow;
        if (now < certificate.NotBefore || now >= certificate.NotAfter) throw new CryptographicException("Certificate is outside its validity interval.");
        using var rsa = certificate.GetRSAPrivateKey();
        if (rsa is null || rsa.KeySize < 2048) throw new CryptographicException("An RSA private key of at least 2048 bits is required.");
        var content = request.ToArray();
        try
        {
            var cms = new SignedCms(new ContentInfo(content), detached: false);
            var signer = new CmsSigner(SubjectIdentifierType.IssuerAndSerialNumber, certificate) { IncludeOption = X509IncludeOption.EndCertOnly };
            cms.ComputeSignature(signer);
            return Convert.ToBase64String(cms.Encode());
        }
        finally { CryptographicOperations.ZeroMemory(content); }
    }
}
