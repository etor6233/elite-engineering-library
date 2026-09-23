using Elite.Arca.Credentials;
using Generated = Elite.Arca.Wsfe.Generated;

namespace Elite.Arca.Wsfe.Bridge;

public sealed record VatLine(int Id, long BaseMinorUnits, long AmountMinorUnits);
public sealed record OtherTaxLine(short Id, string Description, long BaseMinorUnits, int RateBasisPoints, long AmountMinorUnits);
public sealed record AssociatedVoucher(string TaxpayerCuit, int VoucherType, int PointOfSale, long Number, DateOnly IssuedOn);

public sealed class FiscalInvoice
{
    public required string TaxpayerCuit { get; init; }
    public required int PointOfSale { get; init; }
    public required int VoucherType { get; init; }
    public required int Concept { get; init; }
    public required int RecipientDocumentType { get; init; }
    public required string RecipientDocument { get; init; }
    public required int RecipientVatConditionId { get; init; }
    public required string Currency { get; init; }
    public required long VoucherNumber { get; init; }
    public required long TotalMinorUnits { get; init; }
    public required long NetMinorUnits { get; init; }
    public required long VatMinorUnits { get; init; }
    public required long ExemptMinorUnits { get; init; }
    public required long NonTaxedMinorUnits { get; init; }
    public required long OtherTaxMinorUnits { get; init; }
    public required DateOnly IssuedOn { get; init; }
    public DateOnly? ServiceFrom { get; init; }
    public DateOnly? ServiceUntil { get; init; }
    public DateOnly? PaymentDueOn { get; init; }
    public IReadOnlyList<VatLine> VatLines { get; init; } = [];
    public IReadOnlyList<OtherTaxLine> OtherTaxLines { get; init; } = [];
    public IReadOnlyList<AssociatedVoucher> AssociatedVouchers { get; init; } = [];
}

public sealed record SafeWsfeResult(bool Found, bool Authorized, string? Cae, DateOnly? CaeExpiresOn, string ResponseHash, IReadOnlyList<string> ProviderCodes);
public sealed record LastAuthorizedResult(long Number, string ResponseHash, IReadOnlyList<string> ProviderCodes);
public sealed record SafeParameterItem(string Code, string? Description, DateOnly? ValidFrom, DateOnly? ValidUntil, string? VoucherClass, string? EmissionType, string? Blocked, DateOnly? DeregisteredOn);
public sealed record SafeParameterResult(string Kind, string? VoucherClass, IReadOnlyList<SafeParameterItem> Items, string ResponseHash, IReadOnlyList<string> ProviderCodes);

public interface IWsaaAccessSource
{
    Task<LoginTicketAccess> GetAsync(CancellationToken cancellationToken);
}

public interface IWsfeSoap
{
    Task<Generated.FERecuperaLastCbteResponse> LastAuthorizedAsync(Generated.FEAuthRequest auth, int pointOfSale, int voucherType, CancellationToken cancellationToken);
    Task<Generated.FECompConsultaResponse> ConsultAsync(Generated.FEAuthRequest auth, Generated.FECompConsultaReq request, CancellationToken cancellationToken);
    Task<Generated.FECAEResponse> AuthorizeAsync(Generated.FEAuthRequest auth, Generated.FECAERequest request, CancellationToken cancellationToken);
}

public interface IWsfeParameterSoap
{
    Task<Generated.CbteTipoResponse> VoucherTypesAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken);
    Task<Generated.ConceptoTipoResponse> ConceptsAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken);
    Task<Generated.DocTipoResponse> DocumentTypesAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken);
    Task<Generated.IvaTipoResponse> VatRatesAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken);
    Task<Generated.FETributoResponse> OtherTaxesAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken);
    Task<Generated.FEPtoVentaResponse> PointsOfSaleAsync(Generated.FEAuthRequest auth, CancellationToken cancellationToken);
    Task<Generated.CondicionIvaReceptorResponse> RecipientVatConditionsAsync(Generated.FEAuthRequest auth, string voucherClass, CancellationToken cancellationToken);
}

public sealed class WsfeResponseException(string message, string responseHash, IReadOnlyList<string> providerCodes) : IOException(message)
{
    public string ResponseHash { get; } = responseHash;
    public IReadOnlyList<string> ProviderCodes { get; } = providerCodes;
}
