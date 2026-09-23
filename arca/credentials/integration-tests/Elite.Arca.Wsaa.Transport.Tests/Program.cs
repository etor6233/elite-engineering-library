using Elite.Arca.Wsaa.Generated;
using Elite.Arca.Wsaa.Transport;

var ok = new FakeClient("<loginTicketResponse/>");
var transport = new GeneratedWsaaTransport(ok);
var result = await transport.LoginCmsAsync("cms-value", CancellationToken.None);
if (result != "<loginTicketResponse/>" || ok.LastCms != "cms-value" || ok.Calls != 1) throw new InvalidOperationException("generated transport mapping failed");
try { _ = await new GeneratedWsaaTransport(new FakeClient(" ")).LoginCmsAsync("cms", CancellationToken.None); throw new InvalidOperationException("empty response accepted"); } catch (InvalidDataException) { }
try { _ = await transport.LoginCmsAsync(" ", CancellationToken.None); throw new InvalidOperationException("empty CMS accepted"); } catch (ArgumentException) { }
Console.WriteLine("ARCA_GENERATED_WSAA_TRANSPORT_TEST_PASS tests=3 manual_glue=0 production_admitted=false");

sealed class FakeClient(string response) : LoginCMS
{
    public int Calls { get; private set; }
    public string? LastCms { get; private set; }
    public Task<loginCmsResponse> loginCmsAsync(loginCmsRequest request)
    {
        Calls++; LastCms = request.in0; return Task.FromResult(new loginCmsResponse(response));
    }
}
