# Microsoft BC item UOM foundation — V167

## Result

`GO-SUPPLY-FACTORY-INVENTORY-API` 0.11.0 materializes 66 files and the enterprise backend composes 29 packs / 358 implementation files. V167 adds eight files: six `ADAPTED` migration/domain/PostgreSQL/test files and two local `AUTHORED` derivation/HTTP files. No Go or SQL file is represented as Microsoft-authored code, and this evidence does not claim complete Business Central breakbulk or WMS equivalence.

## Fixed official authority

Microsoft BCApps commit `2eae56d704a1fd035d104f333602aea7091b7749`, tree `f9846fb1254c6311985c6131bb2af11c7f169e1b`, MIT:

| File | Bytes | SHA-256 |
|---|---:|---|
| `ItemUnitofMeasure.Table.al` | 20,689 | `ab785d08d4bbc9160c2e068b5660ada732f752d35ae7c94a6677ef0fb79b2009` |
| `CreatePutaway.Codeunit.al` | 77,370 | `b004c0663d237585d8b57a06356431acb2302cdfbbd30f7b82154d282f325b2b` |
| `WhseChangeUnitofMeasure.Report.al` | 8,008 | `d1fca83ae830a8b615c558e15cb6ec3734560614ba66e86d801f3e9794b06186` |
| `Replenishment.Codeunit.al` | 19,064 | `e6d5e3873b4828db8dca49755670bd9815faa7ad1b0396189014af57380b7c4e` |
| `WhseCreatePick.Codeunit.al` | 1,584 | `43ffaded138adb24d0cfc927c276d63ef3dfb2760f83f518bcbc8bbc217163bc` |
| `SCMMovement.Codeunit.al` | 89,432 | `0e63996a0816d9379bdf76acc5712990f10624de81d875f2e499e6944b2955ed` |
| `SCMWMSItemUnitofMeasure.Codeunit.al` | 125,924 | `8faec658fad625a12c83b251bab88553319bc1cf564e7671104f6a19cadf4713` |

Product authority is fixed in Microsoft Learn for item UOM setup, warehouse management and automatic breakbulk. The admitted V167 subset is base factor `1`, positive/aligned alternate factor, explicit precision, exact base conversion and fail-closed mutation/residual. Same-UOM-first selection and Take/Place packaging evidence remain outside V167.

## Executed proof

- Official Go `1.26.7` archive: 74,955,002 bytes, SHA-256 `f4f534a486e4bc3387fa18f08208f2f854b7aaea8a08f2a2d829a914a05abb11`.
- Official PostgreSQL `18.6` archive: 343,808,005 bytes, SHA-256 `fbe23da234ee31547bf8a36d29dfd81e82b849df2d2b78d2eecb43d360252f8c`.
- Clean PostgreSQL applied migrations `0001`–`0033`.
- Focal PostgreSQL test proved `2.5 BOX × 12 = 30 EA`, rejected a residual `0.1 BOX`, rejected factor `2.5`, rejected a duplicate and rejected mutation with SQLSTATE `55000`.
- The first full run exposed five cleanup regressions through the new FK and one resulting outbox contamination. All affected teardown owners were corrected; the database was discarded and recreated.
- From the clean database: `go test ./... -count=1`, `go vet ./...`, `go build ./cmd/...`, migration 0033 down/up and the focal after down/up all passed offline.
- Canonical Markdown rebuilt 66/66 pack files with zero missing, extra or hash differences.
- `VERIFY_LIBRARY_PASS`: 94 packs, 1,075 materializable files, 523 Markdown after this evidence, backend 29/358 and web 6/84.
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: 94 packs, 121 fixed upstream sources and every available offline lane passed; live/provider-dependent lanes remained explicitly skipped.
- Canonical pack: 451,909 bytes, SHA-256 `5fe413deaf47996f3647ba01f06dfea1cec18a90c093e7da8c0d9a8b60931d4f`.
- Failure memory after the complete V167 cycle: 1,611 local lessons plus 202 upstream conditions = 1,813 unique IDs, zero open and zero duplicates. PostgreSQL was stopped; eight `elite-v167-*` temporary test paths remain outside the library because deletion was blocked by execution policy.

## Admission boundary

This foundation is immediately usable for immutable UOM configuration and exact conversion. It is not yet a packaging-stock owner: full breakbulk must connect package identity and UOM-specific balances to receipt, movement, pick, registration, cancellation, lot tracking and concurrency before promotion. No project is production-ready until its actual WMS/ERP ownership, load, recovery, security and business acceptance are proven.
