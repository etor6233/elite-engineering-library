// AUTHORED exact extraction of the existing quote display, shared by private read views.
// Pinned Intl currency digits and BigInt parts preserve all safe minor-unit integers.
// No price, FX, tax, business rounding or acceptance policy is introduced.
export function minorAmountPresentation(minorUnits: number, currency: string, displayLocale = "es-AR") {
  let amountLabel = new Intl.Locale(displayLocale).language === "en" ? "Amount cannot be verified; contact support." : "Importe no verificable; consultá al soporte.";
  let amountValid = false;
  if (Number.isSafeInteger(minorUnits) && minorUnits >= 0 &&
      Intl.supportedValuesOf("currency").includes(currency)) {
    const formatter = new Intl.NumberFormat(displayLocale, { style: "currency", currency: currency, currencyDisplay: "code" });
    const digits = formatter.resolvedOptions().maximumFractionDigits;
    if (typeof digits === "number" && Number.isInteger(digits) && digits >= 0 && digits <= 20) {
      const minor = BigInt(minorUnits);
      const scale = 10n ** BigInt(digits);
      const fraction = (minor % scale).toString().padStart(digits, "0");
      amountLabel = formatter.formatToParts(minor / scale).map(part => part.type === "fraction" ? fraction : part.value).join("");
      amountValid = true;
    }
  }
  return { amountLabel, amountValid };
}
