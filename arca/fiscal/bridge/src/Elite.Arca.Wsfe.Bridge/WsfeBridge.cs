using System.Globalization;
using System.Security.Cryptography;
using System.Text.Json;
using Generated = Elite.Arca.Wsfe.Generated;

namespace Elite.Arca.Wsfe.Bridge;

public sealed class WsfeBridge(IWsaaAccessSource credentials, IWsfeSoap soap)
{
    public async Task<LastAuthorizedResult> LastAuthorizedAsync(FiscalInvoice invoice, CancellationToken cancellationToken = default)
    {
        Validate(invoice, requireNumber: false);
        var response = await soap.LastAuthorizedAsync(await AuthAsync(invoice, cancellationToken), invoice.PointOfSale, invoice.VoucherType, cancellationToken).ConfigureAwait(false);
        var codes = Codes(response.Errors, response.Events, null);
        var hash = Hash(new { operation = "FECompUltimoAutorizado", response.PtoVta, response.CbteTipo, response.CbteNro, codes });
        if (response.Errors is { Length: > 0 } || response.PtoVta != invoice.PointOfSale || response.CbteTipo != invoice.VoucherType || response.CbteNro < 0)
            throw new WsfeResponseException("ARCA last-authorized response is invalid.", hash, codes);
        return new LastAuthorizedResult(response.CbteNro, hash, codes);
    }

    public async Task<SafeWsfeResult> ConsultAsync(FiscalInvoice invoice, CancellationToken cancellationToken = default)
    {
        Validate(invoice, requireNumber: true);
        var response = await soap.ConsultAsync(await AuthAsync(invoice, cancellationToken), new Generated.FECompConsultaReq { PtoVta = invoice.PointOfSale, CbteTipo = invoice.VoucherType, CbteNro = invoice.VoucherNumber }, cancellationToken).ConfigureAwait(false);
        var result = response.ResultGet;
        var codes = Codes(response.Errors, response.Events, result?.Observaciones);
        var safe = new { operation = "FECompConsultar", result?.PtoVta, result?.CbteTipo, result?.CbteDesde, result?.CbteHasta, result?.Resultado, result?.CodAutorizacion, result?.FchVto, codes };
        var hash = Hash(safe);
        if (response.Errors is { Length: > 0 })
        {
            if (response.Errors.All(error => error.Code == 602) && result is null) return new SafeWsfeResult(false, false, null, null, hash, codes);
            throw new WsfeResponseException("ARCA consultation returned errors.", hash, codes);
        }
        if (result is null) throw new WsfeResponseException("ARCA consultation omitted ResultGet without error 602.", hash, codes);
        if (result.PtoVta != invoice.PointOfSale || result.CbteTipo != invoice.VoucherType || result.CbteDesde != invoice.VoucherNumber || result.CbteHasta != invoice.VoucherNumber)
            throw new WsfeResponseException("ARCA consultation identity mismatch.", hash, codes);
        return Result(result.Resultado, result.CodAutorizacion, result.FchVto, hash, codes);
    }

    public async Task<SafeWsfeResult> AuthorizeAsync(FiscalInvoice invoice, CancellationToken cancellationToken = default)
    {
        Validate(invoice, requireNumber: true);
        var response = await soap.AuthorizeAsync(await AuthAsync(invoice, cancellationToken), Request(invoice), cancellationToken).ConfigureAwait(false);
        var detail = response.FeDetResp is { Length: 1 } ? response.FeDetResp[0] : null;
        var codes = Codes(response.Errors, response.Events, detail?.Observaciones);
        var safe = new { operation = "FECAESolicitar", response.FeCabResp?.PtoVta, response.FeCabResp?.CbteTipo, response.FeCabResp?.CantReg, HeaderResult = response.FeCabResp?.Resultado, detail?.CbteDesde, detail?.CbteHasta, DetailResult = detail?.Resultado, detail?.CAE, detail?.CAEFchVto, codes };
        var hash = Hash(safe);
        if (response.Errors is { Length: > 0 } || response.FeCabResp is null || detail is null)
            throw new WsfeResponseException("ARCA authorization returned a general error or incomplete response.", hash, codes);
        if (response.FeCabResp.PtoVta != invoice.PointOfSale || response.FeCabResp.CbteTipo != invoice.VoucherType || response.FeCabResp.CantReg != 1 || detail.CbteDesde != invoice.VoucherNumber || detail.CbteHasta != invoice.VoucherNumber)
            throw new WsfeResponseException("ARCA authorization identity mismatch.", hash, codes);
        return Result(detail.Resultado, detail.CAE, detail.CAEFchVto, hash, codes);
    }

    private async Task<Generated.FEAuthRequest> AuthAsync(FiscalInvoice invoice, CancellationToken cancellationToken)
    {
        var access = await credentials.GetAsync(cancellationToken).ConfigureAwait(false);
        return new Generated.FEAuthRequest { Token = access.Token, Sign = access.Sign, Cuit = long.Parse(invoice.TaxpayerCuit, CultureInfo.InvariantCulture) };
    }

    private static Generated.FECAERequest Request(FiscalInvoice invoice) => new()
    {
        FeCabReq = new Generated.FECAECabRequest { CantReg = 1, PtoVta = invoice.PointOfSale, CbteTipo = invoice.VoucherType },
        FeDetReq = [new Generated.FECAEDetRequest
        {
            Concepto = invoice.Concept,
            DocTipo = invoice.RecipientDocumentType,
            DocNro = long.Parse(invoice.RecipientDocument, CultureInfo.InvariantCulture),
            CbteDesde = invoice.VoucherNumber,
            CbteHasta = invoice.VoucherNumber,
            CbteFch = Date(invoice.IssuedOn),
            ImpTotal = Money(invoice.TotalMinorUnits),
            ImpTotConc = Money(invoice.NonTaxedMinorUnits),
            ImpNeto = Money(invoice.NetMinorUnits),
            ImpOpEx = Money(invoice.ExemptMinorUnits),
            ImpTrib = Money(invoice.OtherTaxMinorUnits),
            ImpIVA = Money(invoice.VatMinorUnits),
            FchServDesde = OptionalDate(invoice.ServiceFrom),
            FchServHasta = OptionalDate(invoice.ServiceUntil),
            FchVtoPago = OptionalDate(invoice.PaymentDueOn),
            MonId = "PES",
            MonCotiz = 1d,
            MonCotizSpecified = true,
            CondicionIVAReceptorId = invoice.RecipientVatConditionId,
            CondicionIVAReceptorIdSpecified = true,
            Iva = invoice.VatLines.Count == 0 ? null! : invoice.VatLines.Select(line => new Generated.AlicIva { Id = line.Id, BaseImp = Money(line.BaseMinorUnits), Importe = Money(line.AmountMinorUnits) }).ToArray(),
            Tributos = invoice.OtherTaxLines.Count == 0 ? null! : invoice.OtherTaxLines.Select(line => new Generated.Tributo { Id = line.Id, Desc = line.Description, BaseImp = Money(line.BaseMinorUnits), Alic = Rate(line.RateBasisPoints), Importe = Money(line.AmountMinorUnits) }).ToArray(),
            CbtesAsoc = invoice.AssociatedVouchers.Count == 0 ? null! : invoice.AssociatedVouchers.Select(value => new Generated.CbteAsoc { Tipo = value.VoucherType, PtoVta = value.PointOfSale, Nro = value.Number, Cuit = value.TaxpayerCuit, CbteFch = Date(value.IssuedOn) }).ToArray()
        }]
    };

    private static SafeWsfeResult Result(string? result, string? cae, string? expiry, string hash, IReadOnlyList<string> codes)
    {
        if (result == "R") return new SafeWsfeResult(true, false, null, null, hash, codes);
        if (result != "A" || string.IsNullOrWhiteSpace(cae) || cae.Any(c => !char.IsAsciiDigit(c)) || cae.Length is < 8 or > 20 || !DateOnly.TryParseExact(expiry, "yyyyMMdd", CultureInfo.InvariantCulture, DateTimeStyles.None, out var date))
            throw new WsfeResponseException("ARCA result is neither an exact approval nor rejection.", hash, codes);
        return new SafeWsfeResult(true, true, cae, date, hash, codes);
    }

    private static IReadOnlyList<string> Codes(Generated.Err[]? errors, Generated.Evt[]? events, Generated.Obs[]? observations) =>
        (errors ?? []).Select(value => $"E:{value.Code}").Concat((events ?? []).Select(value => $"V:{value.Code}")).Concat((observations ?? []).Select(value => $"O:{value.Code}")).ToArray();

    private static string Hash(object value) => Convert.ToHexStringLower(SHA256.HashData(JsonSerializer.SerializeToUtf8Bytes(value)));
    private static string Date(DateOnly value) => value.ToString("yyyyMMdd", CultureInfo.InvariantCulture);
    private static string OptionalDate(DateOnly? value) => value is null ? string.Empty : Date(value.Value);
    private static double Money(long minor) => double.Parse((minor / 100m).ToString("0.00", CultureInfo.InvariantCulture), CultureInfo.InvariantCulture);
    private static double Rate(int basisPoints) => double.Parse((basisPoints / 100m).ToString("0.00", CultureInfo.InvariantCulture), CultureInfo.InvariantCulture);

    private static void Validate(FiscalInvoice value, bool requireNumber)
    {
        if (value.VatLines is null || value.OtherTaxLines is null || value.AssociatedVouchers is null)
            throw new ArgumentException("Fiscal collections must be arrays, not null.", nameof(value));
        if (value.Currency != "ARS" || value.TaxpayerCuit.Length != 11 || value.TaxpayerCuit.Any(c => !char.IsAsciiDigit(c)) || !long.TryParse(value.RecipientDocument, NumberStyles.None, CultureInfo.InvariantCulture, out _) || value.PointOfSale is < 1 or > 99999 || value.VoucherType is < 1 or > 999 || value.Concept is < 1 or > 3 || value.RecipientDocumentType is < 0 or > 999 || value.RecipientVatConditionId is < 1 or > 999 || (requireNumber && value.VoucherNumber <= 0))
            throw new ArgumentException("Fiscal invoice identity is invalid.", nameof(value));
        var total = checked(value.NetMinorUnits + value.VatMinorUnits + value.ExemptMinorUnits + value.NonTaxedMinorUnits + value.OtherTaxMinorUnits);
        if (total <= 0 || total != value.TotalMinorUnits || value.VatLines.Sum(line => line.AmountMinorUnits) != value.VatMinorUnits || value.VatLines.Sum(line => line.BaseMinorUnits) != value.NetMinorUnits || value.OtherTaxLines.Sum(line => line.AmountMinorUnits) != value.OtherTaxMinorUnits)
            throw new ArgumentException("Fiscal monetary detail is inconsistent.", nameof(value));
        if ((value.Concept == 1 && (value.ServiceFrom is not null || value.ServiceUntil is not null || value.PaymentDueOn is not null)) || (value.Concept != 1 && (value.ServiceFrom is null || value.ServiceUntil is null || value.PaymentDueOn is null)))
            throw new ArgumentException("Fiscal service dates are inconsistent with Concepto.", nameof(value));
        var originalType = value.VoucherType switch { 3 => 1, 8 => 6, 13 => 11, _ => 0 };
        if (originalType == 0)
        {
            if (value.VoucherType is not (1 or 6 or 11) || value.AssociatedVouchers.Count != 0)
                throw new ArgumentException("Unsupported voucher type or association.", nameof(value));
        }
        else if (value.AssociatedVouchers.Count != 1 || value.AssociatedVouchers[0] is not { Number: > 0, PointOfSale: >= 1 and <= 99999 } associated || associated.VoucherType != originalType || associated.TaxpayerCuit != value.TaxpayerCuit || associated.IssuedOn > value.IssuedOn)
            throw new ArgumentException("Credit note associated voucher is invalid.", nameof(value));
    }
}
