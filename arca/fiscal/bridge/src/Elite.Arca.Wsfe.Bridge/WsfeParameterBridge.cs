using System.Globalization;
using System.Security.Cryptography;
using System.Text.Json;
using Elite.Arca.Credentials;
using Generated = Elite.Arca.Wsfe.Generated;

namespace Elite.Arca.Wsfe.Bridge;

public sealed class WsfeParameterBridge(IWsaaAccessSource credentials, IWsfeParameterSoap soap)
{
    public async Task<SafeParameterResult> VoucherTypesAsync(string taxpayerCuit, CancellationToken cancellationToken = default)
    {
        var response = await soap.VoucherTypesAsync(await AuthAsync(taxpayerCuit, cancellationToken), cancellationToken).ConfigureAwait(false);
        return Result("voucher_type", null, (response.ResultGet ?? []).Select(value => Item(value.Id, value.Desc, value.FchDesde, value.FchHasta)), response.Errors, response.Events);
    }

    public async Task<SafeParameterResult> ConceptsAsync(string taxpayerCuit, CancellationToken cancellationToken = default)
    {
        var response = await soap.ConceptsAsync(await AuthAsync(taxpayerCuit, cancellationToken), cancellationToken).ConfigureAwait(false);
        return Result("concept", null, (response.ResultGet ?? []).Select(value => Item(value.Id, value.Desc, value.FchDesde, value.FchHasta)), response.Errors, response.Events);
    }

    public async Task<SafeParameterResult> DocumentTypesAsync(string taxpayerCuit, CancellationToken cancellationToken = default)
    {
        var response = await soap.DocumentTypesAsync(await AuthAsync(taxpayerCuit, cancellationToken), cancellationToken).ConfigureAwait(false);
        return Result("document_type", null, (response.ResultGet ?? []).Select(value => Item(value.Id, value.Desc, value.FchDesde, value.FchHasta)), response.Errors, response.Events);
    }

    public async Task<SafeParameterResult> VatRatesAsync(string taxpayerCuit, CancellationToken cancellationToken = default)
    {
        var response = await soap.VatRatesAsync(await AuthAsync(taxpayerCuit, cancellationToken), cancellationToken).ConfigureAwait(false);
        return Result("vat_rate", null, (response.ResultGet ?? []).Select(value => Item(value.Id, value.Desc, value.FchDesde, value.FchHasta)), response.Errors, response.Events);
    }

    public async Task<SafeParameterResult> OtherTaxesAsync(string taxpayerCuit, CancellationToken cancellationToken = default)
    {
        var response = await soap.OtherTaxesAsync(await AuthAsync(taxpayerCuit, cancellationToken), cancellationToken).ConfigureAwait(false);
        return Result("other_tax", null, (response.ResultGet ?? []).Select(value => Item(value.Id, value.Desc, value.FchDesde, value.FchHasta)), response.Errors, response.Events);
    }

    public async Task<SafeParameterResult> PointsOfSaleAsync(string taxpayerCuit, CancellationToken cancellationToken = default)
    {
        var response = await soap.PointsOfSaleAsync(await AuthAsync(taxpayerCuit, cancellationToken), cancellationToken).ConfigureAwait(false);
        return Result("point_of_sale", null, (response.ResultGet ?? []).Select(value => new SafeParameterItem(value.Nro.ToString(CultureInfo.InvariantCulture), null, null, null, null, Required(value.EmisionTipo, "emission type"), Required(value.Bloqueado, "blocked state"), Date(value.FchBaja))), response.Errors, response.Events);
    }

    public async Task<SafeParameterResult> RecipientVatConditionsAsync(string taxpayerCuit, string voucherClass, CancellationToken cancellationToken = default)
    {
        if (voucherClass is not ("A" or "B" or "C" or "M")) throw new ArgumentException("Voucher class is outside the admitted WSFEv1 classes.", nameof(voucherClass));
        var response = await soap.RecipientVatConditionsAsync(await AuthAsync(taxpayerCuit, cancellationToken), voucherClass, cancellationToken).ConfigureAwait(false);
        return Result("recipient_vat_condition", voucherClass, (response.ResultGet ?? []).Select(value =>
        {
            if (!string.Equals(value.Cmp_Clase, voucherClass, StringComparison.Ordinal)) throw new InvalidDataException("ARCA recipient VAT class mismatch.");
            return new SafeParameterItem(value.Id.ToString(CultureInfo.InvariantCulture), Required(value.Desc, "description"), null, null, value.Cmp_Clase, null, null, null);
        }), response.Errors, response.Events);
    }

    private async Task<Generated.FEAuthRequest> AuthAsync(string taxpayerCuit, CancellationToken cancellationToken)
    {
        if (taxpayerCuit.Length != 11 || taxpayerCuit.Any(character => !char.IsAsciiDigit(character))) throw new ArgumentException("Taxpayer CUIT is invalid.", nameof(taxpayerCuit));
        var access = await credentials.GetAsync(cancellationToken).ConfigureAwait(false);
        return new Generated.FEAuthRequest { Token = access.Token, Sign = access.Sign, Cuit = long.Parse(taxpayerCuit, CultureInfo.InvariantCulture) };
    }

    private static SafeParameterItem Item<T>(T code, string description, string from, string until) where T : IFormattable =>
        new(code.ToString(null, CultureInfo.InvariantCulture), Required(description, "description"), DateRequired(from), Date(until), null, null, null, null);
    private static SafeParameterItem Item(string code, string description, string from, string until) =>
        new(RequiredCode(code), Required(description, "description"), DateRequired(from), Date(until), null, null, null, null);

    private static SafeParameterResult Result(string kind, string? voucherClass, IEnumerable<SafeParameterItem> values, Generated.Err[]? errors, Generated.Evt[]? events)
    {
        var codes = (errors ?? []).Select(value => $"E:{value.Code}").Concat((events ?? []).Select(value => $"V:{value.Code}")).ToArray();
        var items = values.OrderBy(value => value.Code, StringComparer.Ordinal).ToArray();
        if (errors is { Length: > 0 } || items.Length == 0 || items.Select(value => value.Code).Distinct(StringComparer.Ordinal).Count() != items.Length)
            throw new WsfeResponseException("ARCA parameter response is invalid.", Hash(new { kind, voucherClass, items, codes }), codes);
        foreach (var item in items)
            if (item.ValidFrom is not null && item.ValidUntil is not null && item.ValidUntil < item.ValidFrom)
                throw new InvalidDataException("ARCA parameter validity interval is invalid.");
        return new SafeParameterResult(kind, voucherClass, items, Hash(new { kind, voucherClass, items, codes }), codes);
    }

    private static string Required(string value, string name) => !string.IsNullOrWhiteSpace(value) && value.Length <= 250 ? value : throw new InvalidDataException($"ARCA parameter {name} is invalid.");
    private static string RequiredCode(string value) => !string.IsNullOrWhiteSpace(value) && value.Length <= 16 ? value : throw new InvalidDataException("ARCA parameter code is invalid.");
    private static DateOnly DateRequired(string value) => Date(value) ?? throw new InvalidDataException("ARCA parameter start date is required.");
    private static DateOnly? Date(string? value)
    {
        if (string.IsNullOrEmpty(value)) return null;
        if (!DateOnly.TryParseExact(value, "yyyyMMdd", CultureInfo.InvariantCulture, DateTimeStyles.None, out var date)) throw new InvalidDataException("ARCA parameter date is invalid.");
        return date;
    }
    private static string Hash(object value) => Convert.ToHexStringLower(SHA256.HashData(JsonSerializer.SerializeToUtf8Bytes(value)));
}
