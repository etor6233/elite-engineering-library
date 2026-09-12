# Go BC Exact Amount Adapter

## 1. Metadata

```yaml
pack_id: "GO-BC-EXACT-AMOUNT-ADAPTER"
pack_version: "0.1.0"
status:
  authority: ELITE_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: REUSABLE_PACK
claim: "Exact integral line amounts, single-currency journal balance and signed reversal normalized to the existing positive-column representation, translated from fixed Microsoft BCApps functions."
stacks: ["Go 1.26.8", "PostgreSQL 18.6"]
compatible_with: ["GO-COMMERCE-PRICING-PAYMENT-API 0.6.2", "GO-ENTERPRISE-ACCOUNTING-LEDGER-API 0.1.1"]
incompatible_with: ["fractional quantity or rounding precision other than one minor unit", "implicit tax or FX", "Business Central correction-turnover equivalence", "unscoped journal or external ERP authority"]
license_expression: "MIT AND LicenseRef-Workspace-Owner"
upstream_sources: ["https://github.com/microsoft/BCApps/tree/2eae56d704a1fd035d104f333602aea7091b7749"]
verified_at: "2026-09-11"
```

The narrow mathematical translation is admitted, not the whole business platform. Two files are ADAPTED, three are local test/document glue, and one is the verbatim Microsoft MIT license. No caller pack is relabeled. Microsoft did not author or approve this Go/PostgreSQL integration.

## 2. Applicability

Select with both existing Commerce and Accounting owners when amounts use integral minor units, one authoritative currency per journal and nonnegative one-sided debit/credit storage. Existing request validation, quote acceptance and posting authorization remain their owners. Reject fractional units, alternate rounding, FX, tax, automatic discount policy, statutory accounting, negative correction turnover or full BC equivalence. This pack has no external calls and requests no secrets.

## 3. Architecture contract

LineAmount translates the exact SalesLine.UpdateAmounts expression before the int64 storage boundary. JournalTotals translates exact accumulation and nonzero-balance rejection; representation guards reject negative/opposite sides and overflow before persistence. Reverse SQL negates the signed amount then normalizes to positive sides inside the existing serializable transaction. Original entry links, tenant/org, period, idempotence, outbox and one-reversal guards remain existing owners. No new tables, migration or worker. API and durable formats stay compatible; rollback restores prior caller+adapter revision, retaining the two newly diagnosed invalid-input regressions. Cost is linear in caller-provided journal lines and constant per line; integer operands are bounded int64. No state or telemetry is owned by the arithmetic package.

## 4. Exact file manifest

```text
CREATE internal/bcamounts/amounts.go
CREATE internal/bcamounts/amounts_test.go
CREATE internal/platform/postgres/bc_accounting_sql.go
CREATE internal/accounting/amount_safety_test.go
CREATE docs/provenance/BC_AMOUNT_DERIVATION.md
CREATE licenses/Microsoft-BCApps-MIT.txt
```

## 5. Materialization blocks

### FILE: `internal/bcamounts/amounts.go`

```yaml
block_id: "GO-BC-EXACT-AMOUNT-ADAPTER:file1:v1"
operation: CREATE
provenance: ADAPTED
source: "https://github.com/microsoft/BCApps/tree/2eae56d704a1fd035d104f333602aea7091b7749; functions/line/SHA mapping in docs/provenance/BC_AMOUNT_DERIVATION.md; local representation/transaction delta declared"
license: "MIT"
sha256: "c0a5ba4801fa4c7772335e0fb2519c69d4d20e03ff47c3fee1e959dbad016fd5"
variables: []
secrets_allowed: false
```

````go
// Copyright (c) Microsoft Corporation.
// SPDX-License-Identifier: MIT
// Adapted from microsoft/BCApps at 2eae56d704a1fd035d104f333602aea7091b7749.
// This Go translation and its representation guards are local modifications.
// See docs/provenance/BC_AMOUNT_DERIVATION.md and licenses/Microsoft-BCApps-MIT.txt.
package bcamounts

import (
	"errors"
	"math/big"
)

var ErrAmount = errors.New("amount is outside the supported exact representation")
var ErrBalance = errors.New("journal is out of balance")

// LineAmount translates SalesLine.Table.al UpdateAmounts, line 5883:
// Round(Quantity * UnitPrice, AmountRoundingPrecision) - LineDiscountAmount.
// This profile uses integral quantities, minor units and rounding precision 1;
// the product is already exact, so rounding is an identity. Tax, FX, fractional
// quantities and percentage discounts are outside this function's contract.
// big.Int preserves AL's exact arithmetic before the explicit int64 boundary.
func LineAmount(quantity, unitPriceMinor, lineDiscountMinor int64) (int64, error) {
	var amount, right big.Int
	amount.SetInt64(quantity)
	right.SetInt64(unitPriceMinor)
	amount.Mul(&amount, &right)
	right.SetInt64(lineDiscountMinor)
	amount.Sub(&amount, &right)
	if !amount.IsInt64() {
		return 0, ErrAmount
	}
	return amount.Int64(), nil
}

// Entry is the existing library representation, not a Business Central type.
// Exactly one nonnegative side is positive; source AL signed amounts are mapped
// to Debit-Credit. The representation validation is a declared local delta.
type Entry struct{ Debit, Credit int64 }

// JournalTotals translates GenJnlPostBatch.ProcessBalanceOfLines' accumulation
// and CheckBalance into the library's single currency, single journal boundary.
// Document/currency grouping is supplied and authorized by the caller.
func JournalTotals(entries []Entry) (int64, int64, error) {
	var debit, credit, value big.Int
	for _, entry := range entries {
		if entry.Debit < 0 || entry.Credit < 0 || (entry.Debit > 0) == (entry.Credit > 0) {
			return 0, 0, ErrAmount
		}
		value.SetInt64(entry.Debit)
		debit.Add(&debit, &value)
		value.SetInt64(entry.Credit)
		credit.Add(&credit, &value)
	}
	if !debit.IsInt64() || !credit.IsInt64() {
		return 0, 0, ErrAmount
	}
	d, c := debit.Int64(), credit.Int64()
	if err := CheckBalance(d, c); err != nil {
		return 0, 0, err
	}
	return d, c, nil
}

// CheckBalance is the restricted equivalent of GenJnlPostBatch.CheckBalance:
// nonzero signed balance fails. Positive turnover is an existing library rule;
// it is explicit here and is not attributed to Microsoft.
func CheckBalance(debit, credit int64) error {
	if debit <= 0 || credit <= 0 {
		return ErrAmount
	}
	if debit != credit {
		return ErrBalance
	}
	return nil
}
````

### FILE: `internal/bcamounts/amounts_test.go`

```yaml
block_id: "GO-BC-EXACT-AMOUNT-ADAPTER:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local verification/documentation glue for the declared fixed-source adaptation; no vendor authorship claim"
license: "LicenseRef-Workspace-Owner"
sha256: "ecf4173acd5e875b22504a8ca16b968093ff70c9a40be6320ceacf8ce4e84a22"
variables: []
secrets_allowed: false
```

````go
// SPDX-License-Identifier: LicenseRef-Workspace-Owner
// Local test harness for the declared BCApps translation; not a Microsoft suite.
package bcamounts

import (
	"errors"
	"math"
	"testing"
)

func TestSalesLineUpdateAmountsExactIntegerProfile(t *testing.T) {
	// Arithmetic oracle is the fixed SalesLine.UpdateAmounts expression. Inputs
	// are deterministic local fixtures, not purported Microsoft fixture values.
	cases := []struct {
		name                 string
		q, p, discount, want int64
	}{
		{"single", 1, 10001, 0, 10001},
		{"quantity", 3, 1200, 0, 3600},
		{"fixed discount", 3, 1200, 100, 3500},
		{"zero", 0, 1200, 0, 0},
		{"credit quantity", -2, 1200, 0, -2400},
		{"maximum", 1, math.MaxInt64, 0, math.MaxInt64},
		{"minimum", 1, math.MinInt64, 0, math.MinInt64},
		{"intermediate above range final exact", 2, math.MaxInt64, math.MaxInt64, math.MaxInt64},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := LineAmount(tc.q, tc.p, tc.discount)
			if err != nil || got != tc.want {
				t.Fatalf("got=%d err=%v want=%d", got, err, tc.want)
			}
		})
	}
	for _, args := range [][3]int64{{2, math.MaxInt64, 0}, {-1, math.MinInt64, 0}, {math.MaxInt64, math.MaxInt64, 0}, {1, math.MinInt64, 1}} {
		if _, err := LineAmount(args[0], args[1], args[2]); !errors.Is(err, ErrAmount) {
			t.Fatalf("accepted out-of-range %+v: %v", args, err)
		}
	}
}

func TestPostingBalanceMatchesSignedAmountConservation(t *testing.T) {
	for amount := int64(1); amount <= 200; amount++ {
		// Adapted invariant from GenJnlPostBatch.CheckBalance and the reverse
		// tests' out-of-balance oracle. Deterministic amounts replace random AL
		// setup and the BC runtime; no claim that the Microsoft suite ran.
		d, c, err := JournalTotals([]Entry{{Debit: amount}, {Credit: amount}})
		if err != nil || d != amount || c != amount {
			t.Fatalf("amount=%d totals=%d/%d err=%v", amount, d, c, err)
		}
		if err := CheckBalance(amount, amount+1); !errors.Is(err, ErrBalance) {
			t.Fatalf("unbalanced amount %d accepted", amount)
		}
	}
	if _, _, err := JournalTotals([]Entry{{Debit: math.MaxInt64}, {Credit: math.MaxInt64}}); err != nil {
		t.Fatal(err)
	}
	for _, entries := range [][]Entry{
		nil, {{Debit: 1, Credit: 1}}, {{Debit: -1}, {Credit: -1}},
		{{Debit: 10, Credit: -5}, {Debit: -5, Credit: 10}},
		{{Debit: math.MaxInt64}, {Debit: math.MaxInt64}, {Debit: 3}, {Credit: math.MaxInt64}, {Credit: math.MaxInt64}, {Credit: 3}},
	} {
		if _, _, err := JournalTotals(entries); err == nil {
			t.Fatalf("invalid representation accepted: %+v", entries)
		}
	}
}

func FuzzIntegralSalesLineConservation(f *testing.F) {
	for _, seed := range [][2]uint64{{0, 0}, {1, 10001}, {3, 1200}, {1, math.MaxInt64}, {2, math.MaxInt64}, {math.MaxInt64, math.MaxInt64}, {math.MaxInt64, 1}} {
		f.Add(seed[0], seed[1])
	}
	f.Fuzz(func(t *testing.T, rawQuantity, rawPrice uint64) {
		q, p := rawQuantity&math.MaxInt64, rawPrice&math.MaxInt64
		amount, err := LineAmount(int64(q), int64(p), 0)
		// Independent overflow oracle uses division; it does not repeat the
		// big-integer multiplication used by the translated implementation.
		if p != 0 && q > math.MaxInt64/p {
			if !errors.Is(err, ErrAmount) {
				t.Fatalf("overflow accepted: q=%d p=%d result=%d err=%v", q, p, amount, err)
			}
			return
		}
		if err != nil || amount != int64(q*p) {
			t.Fatalf("amount not conserved: q=%d p=%d result=%d err=%v", q, p, amount, err)
		}
	})
}
````

### FILE: `internal/platform/postgres/bc_accounting_sql.go`

```yaml
block_id: "GO-BC-EXACT-AMOUNT-ADAPTER:file3:v1"
operation: CREATE
provenance: ADAPTED
source: "https://github.com/microsoft/BCApps/tree/2eae56d704a1fd035d104f333602aea7091b7749; functions/line/SHA mapping in docs/provenance/BC_AMOUNT_DERIVATION.md; local representation/transaction delta declared"
license: "MIT"
sha256: "787836439f350f98e31c27ddc27029aa6fe375707c142692c4a19abb5f0843d7"
variables: []
secrets_allowed: false
```

````go
// Copyright (c) Microsoft Corporation.
// SPDX-License-Identifier: MIT
// Adapted from microsoft/BCApps at 2eae56d704a1fd035d104f333602aea7091b7749,
// GenJnlPostReverse.Codeunit.al:209-257. Local PostgreSQL representation delta.
// The source negates signed G/L amounts; this boundary maps the result to the
// library's existing nonnegative debit/credit columns. Original entry links,
// journal state, tenant authorization and transaction remain the caller's owner.
// Negative source correction columns and turnover are not claimed equivalent.
package postgres

const reverseJournalLinesBCSQL = `insert into accounting.journal_line(tenant_id,journal_id,line_no,account_code,description,debit_minor_units,credit_minor_units)
select tenant_id,$3,line_no,account_code,'Reversal: '||description,
       greatest(credit_minor_units-debit_minor_units,0),
       greatest(debit_minor_units-credit_minor_units,0)
from accounting.journal_line where tenant_id=$1 and journal_id=$2`
````

### FILE: `internal/accounting/amount_safety_test.go`

```yaml
block_id: "GO-BC-EXACT-AMOUNT-ADAPTER:file4:v1"
operation: CREATE
provenance: AUTHORED
source: "local verification/documentation glue for the declared fixed-source adaptation; no vendor authorship claim"
license: "LicenseRef-Workspace-Owner"
sha256: "01cc31a33b42ef217dfec6483a68a5dc8bbdf6475e2dca3df4481f8240105c15"
variables: []
secrets_allowed: false
```

````go
package accounting

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"
)

// Local representation boundary regression. These are not Microsoft fixtures.
func TestJournalRejectsNegativeSidesAndWrappedTotalsBeforeRepository(t *testing.T) {
	cases := map[string][]Line{
		"negative opposite sides": {
			{LineNo: 1, AccountCode: "CASH", DebitMinorUnits: 10, CreditMinorUnits: -5},
			{LineNo: 2, AccountCode: "REVENUE", DebitMinorUnits: -5, CreditMinorUnits: 10},
		},
		"int64 totals wrap to positive": {
			{LineNo: 1, AccountCode: "CASH", DebitMinorUnits: math.MaxInt64},
			{LineNo: 2, AccountCode: "CASH", DebitMinorUnits: math.MaxInt64},
			{LineNo: 3, AccountCode: "CASH", DebitMinorUnits: 3},
			{LineNo: 4, AccountCode: "REVENUE", CreditMinorUnits: math.MaxInt64},
			{LineNo: 5, AccountCode: "REVENUE", CreditMinorUnits: math.MaxInt64},
			{LineNo: 6, AccountCode: "REVENUE", CreditMinorUnits: 3},
		},
	}
	for name, lines := range cases {
		t.Run(name, func(t *testing.T) {
			repo := &fakeRepo{}
			_, err := NewService(repo, &ids{}).CreateJournal(context.Background(), "tenant", Journal{
				OrganizationID: "org", PeriodID: "period", SourceType: "SALE", SourceID: "order",
				Currency: "ARS", PostingDate: time.Date(2026, 9, 11, 0, 0, 0, 0, time.UTC), Lines: lines,
			})
			if !errors.Is(err, ErrInvalid) || repo.journal.ID != "" {
				t.Fatalf("invalid journal reached repository: error=%v totalDebit=%d totalCredit=%d persisted=%t", err, repo.journal.TotalDebitMinorUnits, repo.journal.TotalCreditMinorUnits, repo.journal.ID != "")
			}
		})
	}
}
````

### FILE: `docs/provenance/BC_AMOUNT_DERIVATION.md`

```yaml
block_id: "GO-BC-EXACT-AMOUNT-ADAPTER:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local verification/documentation glue for the declared fixed-source adaptation; no vendor authorship claim"
license: "LicenseRef-Workspace-Owner"
sha256: "8c62106a8cca5cae11dc6e892ea6132ddf50c408fae8cb36949ab9314622b110"
variables: []
secrets_allowed: false
```

````markdown
# BCApps exact amount adaptation — verified local contract

This is local Go/SQL translated from Microsoft BCApps MIT source, not Microsoft-authored Go or a Business Central platform certification. Commit: `2eae56d704a1fd035d104f333602aea7091b7749`.

| Source | SHA-256 | Source function | Destination and supported contract |
|---|---|---|---|
| `src/Layers/W1/BaseApp/Sales/Document/SalesLine.Table.al` | `ca65615dc05bc555a878ba2635c955897930332ef171e313a00103074fbc7278` | `UpdateAmounts`, line5883 | `internal/bcamounts/amounts.go:LineAmount`: integral quantities/minor unit precision1, exact multiplication then fixed discount subtraction. Commerce uses discount0. No fiscal/FX/percentage-discount behavior is selected. |
| `src/Layers/W1/BaseApp/Finance/GeneralLedger/Posting/GenJnlPostBatch.Codeunit.al` | `be8be19c361c4610f0c3894eb97cfe0554ab54c1de0aaa491a8a3903fee4b150` | `ProcessBalanceOfLines` line431 and `CheckBalance` line546 | `JournalTotals` and `CheckBalance`: sum exact signed balance within an already scoped single journal/currency; no nonzero imbalance. Nonnegative one-sided columns, positive turnover and int64 storage bounds are explicit local representation constraints. |
| `src/Layers/W1/BaseApp/Finance/GeneralLedger/Reversal/GenJnlPostReverse.Codeunit.al` | `97d9268c6be04de015c50595e94b935710cd62610f987a08d66c9aab6790cba5` | `ReverseGLEntry` line209, signed negation line225 | `internal/platform/postgres/bc_accounting_sql.go`: negate signed debit-minus-credit, normalize into existing positive columns; source correction turnover is not claimed equivalent. Origin links and single reversal remain caller transaction rules. |

Source receipt files are in the staging sibling `official-source-inspection`; raw bytes were matched to Git blobs of the fixed commit and SHA-256. The root MIT license is preserved verbatim in `licenses/Microsoft-BCApps-MIT.txt`, SHA `c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383`.

The amount/balance arithmetic and reverse SQL expression are ADAPTED. Scope, IDs, payload validation, authorization, retries, transaction/outbox wiring, local error mapping and storage representation guards remain local. Existing owner files remain AUTHORED glue plus calls to these translations; this does not automatically admit all remaining domain behavior within those files. No legacy entire pack is relabeled.

Counterexamples discovered before the change are retained in `business-red.log`: negative opposite columns yielded5/5 and int64 wrap yielded1/1 at the domain repository boundary. Local regression tests are clearly labeled local; the historical Microsoft AL suite has not been executed on this machine. `GenJnlPostBatch.CheckBalance` and `ERMReverseGLEntries.ReverseForceDocBalanceNo` supply the no-unbalanced-posting oracle; deterministic local vectors and storage boundary cases verify this translation.

The local reference passed amount/balance unit regressions, exact Go1.26.8 native fuzzing (7 seeds,573594 executions/3s),54 current PostgreSQL migrations, connected commerce/order/stock/payment-request and accounting concurrent posting/reversal/close, vet and complete module build. Reconstruction and the G0–G8 decision are recorded in reconstruction_evidence/BC_EXACT_AMOUNT_ADAPTATION_V402.md. This admits the narrow adapter only; wider business owners, payment providers, fiscality, production release and target acceptance retain their own gates.
````

### FILE: `licenses/Microsoft-BCApps-MIT.txt`

```yaml
block_id: "GO-BC-EXACT-AMOUNT-ADAPTER:file6:v1"
operation: CREATE
provenance: VERBATIM
source: "https://github.com/microsoft/BCApps/tree/2eae56d704a1fd035d104f333602aea7091b7749/LICENSE; exact upstream MIT bytes"
license: "MIT"
sha256: "c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383"
variables: []
secrets_allowed: false
```

````text
    MIT License

    Copyright (c) Microsoft Corporation.

    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:

    The above copyright notice and this permission notice shall be included in all
    copies or substantial portions of the Software.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE
````

## 6. Configuration surface

No variables, secrets, implicit provider selection or runtime switches. Quantity, unit price and line discount are explicit function inputs; the Commerce caller supplies discount zero. The caller owns currency and grouping. Missing/incompatible semantics require rejection before invoking this adapter.

## 7. Dependency bill

| Component | Exact identity | Purpose | License |
|---|---|---|---|
| Microsoft BCApps selected source functions | 2eae56d704a1fd035d104f333602aea7091b7749; three SHA-bound AL files in derivation | Source for adaptation only; no BC runtime | MIT |
| Go | 1.26.8; go.exe21761eceb9302062c9623fb699f332c8c7fe000f15f70efe8da01a2cfbbc16b9 | errors/math-big and tests | BSD-3-Clause |
| PostgreSQL | 18.6; existing reference runtime lock | Existing owner transactions; no new dependency | PostgreSQL |

No new module dependency, network, package manager or dynamic code execution.

## 8. Apply order

Compose this adapter with the exact compatible Commerce and Accounting owners. The SQL constant belongs to their existing postgres package; tests intentionally reuse the existing accounting fake repository. Do not deploy the standalone materialization as an application. Existing databases receive no migration. The profile must bind caller revisions and preserve the Microsoft notice. Keep the original business decisions outside this adapter's claim.

## 9. Verification

Run go test ./internal/bcamounts ./internal/accounting; native fuzz target FuzzIntegralSalesLineConservation uses a finite explicit budget. On a disposable loopback PostgreSQL reference run TestAccountingPostingConcurrencyReversalAndClose and TestCommercePriceOrderAllocationPaymentFlow, rejecting skips. Run vet/build on the composed module. Compare all six materialized hashes and caller deltas to evidence. These checks do not replace provider contracts, fiscal rules, target authorization, SCA, deployment or production acceptance.

## 10. Reconstruction evidence

reconstruction_evidence/BC_EXACT_AMOUNT_ADAPTATION_V402.md records source identity, exact before/after, FAIL808 red/green, Go unit/fuzz, fresh54-migration PostgreSQL integration, vet/build, six-file reconstruction and the G0–G8 gate decision. Claim and licensing are per source function; unrelated AUTHORED business remains explicitly pending.
