using System.Security.Cryptography;
using System.Security.Cryptography.Pkcs;
using System.Security.Cryptography.X509Certificates;
using System.Text;
using System.Xml;
using Elite.Arca.Credentials;

var tests = new List<(string Name, Func<Task> Run)>
{
    ("options", () => { WsaaOptions.ElectronicInvoice.Validate(); Throws<ArgumentException>(() => new WsaaOptions("WS FE!", TimeSpan.FromHours(12), TimeSpan.Zero, TimeSpan.FromMinutes(1)).Validate()); return Task.CompletedTask; }),
    ("request", () => { var bytes = Factory(new FixedClock()).Create(); var xml = Encoding.UTF8.GetString(bytes); Check(xml.Contains("<service>wsfe</service>", StringComparison.Ordinal)); Check(xml.Contains("<uniqueId>7</uniqueId>", StringComparison.Ordinal)); return Task.CompletedTask; }),
    ("cms-attached", () => { using var cert = Certificate(); var request = Factory(new FixedClock()).Create(); var encoded = new CmsRequestSigner().Sign(cert, request); var cms = new SignedCms(); cms.Decode(Convert.FromBase64String(encoded)); cms.CheckSignature(verifySignatureOnly: true); Check(cms.ContentInfo.Content.SequenceEqual(request)); return Task.CompletedTask; }),
    ("thumbprint", () => { Check(CertificateStoreLoader.NormalizeThumbprint("AA bb-" + new string('c', 36)) == "AABB" + new string('C', 36)); return Task.CompletedTask; }),
    ("parser", () => { var now = FixedClock.Value; var ticket = WsaaLoginResponseParser.Parse(Response(now.AddHours(1)), now, TimeSpan.FromHours(13)); Check(ticket.Token == "token-secret" && ticket.ToString().Contains("REDACTED", StringComparison.Ordinal) && !ticket.ToString().Contains("token-secret", StringComparison.Ordinal)); return Task.CompletedTask; }),
    ("parser-dtd", () => { Throws<XmlException>(() => WsaaLoginResponseParser.Parse("<!DOCTYPE x [<!ENTITY e SYSTEM 'file:///x'>]><loginTicketResponse>&e;</loginTicketResponse>", FixedClock.Value, TimeSpan.FromHours(13))); return Task.CompletedTask; }),
    ("expired", () => { Throws<InvalidDataException>(() => WsaaLoginResponseParser.Parse(Response(FixedClock.Value.AddSeconds(-1)), FixedClock.Value, TimeSpan.FromHours(13))); return Task.CompletedTask; }),
    ("single-flight", SingleFlight),
    ("sign-rejects-no-key", () => { using var rsa = RSA.Create(2048); var req = new CertificateRequest("CN=NoKey", rsa, HashAlgorithmName.SHA256, RSASignaturePadding.Pkcs1); using var withKey = req.CreateSelfSigned(DateTimeOffset.UtcNow.AddDays(-1), DateTimeOffset.UtcNow.AddDays(1)); using var publicOnly = X509CertificateLoader.LoadCertificate(withKey.Export(X509ContentType.Cert)); Throws<CryptographicException>(() => new CmsRequestSigner().Sign(publicOnly, [1,2,3])); return Task.CompletedTask; })
};
foreach (var test in tests) { await test.Run(); Console.WriteLine($"PASS {test.Name}"); }
Console.WriteLine($"ARCA_WSAA_CREDENTIAL_TEST_PASS tests={tests.Count} secrets_logged=0");

static LoginTicketRequestFactory Factory(IUtcClock clock) => new(WsaaOptions.ElectronicInvoice, clock, new FixedIds());
static async Task SingleFlight()
{
    var clock = new FixedClock(); using var cert = Certificate(); var transport = new FakeTransport(clock);
    using var provider = new WsaaCredentialProvider(WsaaOptions.ElectronicInvoice, clock, Factory(clock), new CmsRequestSigner(), cert, transport);
    var results = await Task.WhenAll(Enumerable.Range(0, 32).Select(_ => provider.GetAsync()));
    Check(transport.Calls == 1); Check(results.All(x => ReferenceEquals(x, results[0]))); _ = await provider.GetAsync(); Check(transport.Calls == 1);
}
static X509Certificate2 Certificate() { using var rsa = RSA.Create(2048); var req = new CertificateRequest("CN=Elite ARCA Test", rsa, HashAlgorithmName.SHA256, RSASignaturePadding.Pkcs1); return req.CreateSelfSigned(DateTimeOffset.UtcNow.AddDays(-1), DateTimeOffset.UtcNow.AddDays(1)); }
static string Response(DateTimeOffset expiration) => $"<loginTicketResponse><header><expirationTime>{XmlConvert.ToString(expiration)}</expirationTime></header><credentials><token>token-secret</token><sign>sign-secret</sign></credentials></loginTicketResponse>";
static void Check(bool condition) { if (!condition) throw new InvalidOperationException("check failed"); }
static void Throws<T>(Action action) where T : Exception { try { action(); } catch (T) { return; } throw new InvalidOperationException($"expected {typeof(T).Name}"); }
sealed class FixedClock : IUtcClock { public static readonly DateTimeOffset Value = new(2026, 8, 30, 12, 0, 0, TimeSpan.Zero); public DateTimeOffset UtcNow => Value; }
sealed class FixedIds : ILoginTicketIdSource { public long Next() => 7; }
sealed class FakeTransport(FixedClock clock) : IWsaaTransport { public int Calls; public async Task<string> LoginCmsAsync(string cmsBase64, CancellationToken cancellationToken) { if (string.IsNullOrWhiteSpace(cmsBase64)) throw new InvalidOperationException("CMS missing"); Interlocked.Increment(ref Calls); await Task.Delay(20, cancellationToken); return $"<loginTicketResponse><header><expirationTime>{XmlConvert.ToString(clock.UtcNow.AddHours(12))}</expirationTime></header><credentials><token>token-secret</token><sign>sign-secret</sign></credentials></loginTicketResponse>"; } }
