# Business Central FX derivation

Algorithm source: Microsoft BCApps commit `2eae56d704a1fd035d104f333602aea7091b7749`, `src/Layers/W1/BaseApp/Finance/Currency/CurrencyExchangeRate.Table.al`, SHA256 `7f25d2996d4f2098cdfff16dfc8cb35be6101ed5052e12fbfebffab826759ac6`. Exact Git blob, byte length and SHA256 were independently checked before this connected implementation. `BC_FX_SOURCE_LOCK.json` records that file, Currency.Table.al, the enum and three inspected upstream test codeunits.

| Upstream method | Narrow adaptation | Local boundary |
|---|---|---|
| ExchangeAmtFCYToFCY, lines371–430 | ExchangeExact divides/multiplies exact Exchange Rate Amount and Relational Exch. Rate Amount through the explicit local base. | Only empty Relational Currency Code is accepted. Relational chains, event subscribers, adjustment factors and BC runtime are not implemented. |
| FindCurrency2, lines459–469 | FindLast chooses the last effective Starting Date at or before the explicit request date. | No WorkDate fallback. Duplicate identities/dates and unsupported rows fail closed. |
| ExchangeAmount, lines329–346 | Calls explicit destination precision and preserves same-currency/zero bypass. | Currency minor-unit decimals and rounding quantum are explicit configuration. Same-currency and zero inputs preserve the upstream early return. |

`internal/bcfx/exchange.go` and `selection.go` are ADAPTED, retain Microsoft copyright/MIT and declare Go representation, bounds and error handling. `rounding.go` is separately AUTHORED representation/rounding glue, not an adapted AL built-in implementation. The surrounding snapshot, hash/identity/expiry checks, API, persistence and host are AUTHORED integration glue. The upstream test codeunits are inspected source, not an executed AL suite; in particular CurrencyUT constrains currency codes, ERMCurrencyFactor verifies division by the currency factor and explicit precision, and ERMChangeExchangeRate exercises editable factor behavior that this immutable profile does not expose.

The initial raw GitHub test-source request failed with HTTP503. The same fixed Git blobs were subsequently obtained through the official GitHub blob API and verified; the original failure remains in the admission history. No branch or revision changed.

The profile accepts a supplied immutable source document, with source ID/revision and rate identities, whose exact bytes are SHA256-bound by the separately activated profile. This verifies configured bytes and identities; it does not certify the economic truth or legal acceptability of a supplied rate. Synthetic fixture rates are explicitly marked synthetic. A future materialized project selects its authorized source and dates without requiring live credentials to exercise this infrastructure.

The accounting conversion receipt is historical conversion evidence and creates no journal, posting, settlement, tax calculation or movement of money. All such effects remain with their existing owners. No signature, financial recommendation, live feed or complete Business Central runtime equivalence is claimed.

License: `licenses/Microsoft-BCApps-MIT.txt` is the exact upstream MIT notice. No new external dependency is added by this adapter. The obsolete binary64 in-memory FX pack remains outside this profile and is not relabeled.

Source correction FX-DOC-ROUND-001: the official System.Round page inspected for this candidate (https://learn.microsoft.com/en-us/dynamics365/business-central/dev-itpro/developer/methods-auto/system/system-round-method, last-updated2025-01-28) has an inconsistent negative example: -1234.56789 with precision1 is printed as -1234, while its nearest rule and other negative examples imply -1235. The root candidate comment asserting AL built-in equality was corrected before publication. This profile explicitly selects nearest/ties-away-from-zero and verifies it with an independent Python Decimal oracle; it makes no AL runtime-equivalence claim for that contradiction. No arbitrary exclusion of negative amounts is introduced.
