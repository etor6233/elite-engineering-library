using System.Globalization;
using System.Text;
using System.Xml;

namespace Elite.Arca.Credentials;

public sealed class LoginTicketRequestFactory(WsaaOptions options, IUtcClock clock, ILoginTicketIdSource ids)
{
    public byte[] Create()
    {
        options.Validate();
        var id = ids.Next();
        if (id <= 0) throw new InvalidOperationException("Login ticket ID must be positive.");
        var now = clock.UtcNow.ToUniversalTime();
        var settings = new XmlWriterSettings { Encoding = new UTF8Encoding(false), OmitXmlDeclaration = false, Indent = false, NewLineHandling = NewLineHandling.None };
        using var stream = new MemoryStream();
        using (var writer = XmlWriter.Create(stream, settings))
        {
            writer.WriteStartDocument();
            writer.WriteStartElement("loginTicketRequest");
            writer.WriteAttributeString("version", "1.0");
            writer.WriteStartElement("header");
            writer.WriteElementString("uniqueId", id.ToString(CultureInfo.InvariantCulture));
            writer.WriteElementString("generationTime", XmlConvert.ToString(now - options.GenerationSkew));
            writer.WriteElementString("expirationTime", XmlConvert.ToString(now + options.TicketLifetime));
            writer.WriteEndElement();
            writer.WriteElementString("service", options.Service);
            writer.WriteEndElement();
            writer.WriteEndDocument();
        }
        return stream.ToArray();
    }
}
