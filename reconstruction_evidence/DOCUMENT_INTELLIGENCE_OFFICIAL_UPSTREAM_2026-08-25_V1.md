# Official document-intelligence upstream evidence V1

Verified at: 2026-08-25.

## Environment

- Microsoft Windows 11 Pro `10.0.26200`, 64-bit;
- PowerShell `7.6.4`;
- Python `3.12.13`;
- CPU execution; no GPU claim;
- all source checkouts were detached at the commits below and lived outside the library.

## Docling 2.122.0

- repository commit: `facdd9bae5882c24c37d4e4c3ecb6e5510667c4d`;
- source archive SHA-256: `0b7a2cc4bc876b8ff40c0de13b19dfccd437aadec0e42223b4dc8cc51855e078` over `147274332` bytes;
- MIT `LICENSE` SHA-256: `5548fe5439075583606b3ce553c5c18f2397fb0abf0ad6b9c37dd6db7f137e85`;
- upstream `uv.lock` SHA-256: `55cbe9d11c3b5cf518313baf9a81a40e27cae06dfad3ac52e146d99e093a19bf`.

The full official distribution was required. A slim installation imported most backends but `DocumentConverter` failed because `docling_ibm_models` was absent. The complete install resolved, among others, `docling 2.122.0`, `docling-core 2.92.0`, `docling-parse 7.16.0`, `docling-ibm-models 3.14.0`, `rapidocr 3.9.2`, `torch 2.13.0` and `torchvision 0.28.0`.

An official PDF fixture, `tests/data/pdf/sources/multi_page.pdf`, converted successfully to structured JSON and Markdown. First execution took 58.949 seconds including official OCR/layout model downloads; the cached second execution took 10.715 seconds. Both runs produced identical logical JSON and Markdown hashes. The file hashes after Windows newline serialization were:

```text
Markdown SHA-256 A83591D45CC15DFBCD7F000B0F3710B28FBC645E3664E5E1BFF99C1C748C4D2B
JSON     SHA-256 E597729113D57E15FE008941E327CBC88847FC7BA5CB50CB91139E2BD177FEA8
```

An official RapidOCR/Torch conversion of `tests/data/ocr/sources/ocr_test.pdf` passed Docling's own fuzzy ground-truth verifier and returned `ConversionStatus.SUCCESS`.

Selected official tests:

```text
backend CSV/Markdown/Docling JSON/invalid input: 34 passed, 1 failed, 1 skipped
PDFium backend:                              9 passed
remaining invalid-input tests:              15 passed, 1 deselected
```

The one retained failure was `test_convert_remote_too_large_filesize_limit_wout_exception`: the source URL was rejected by the current official URL policy before the mocked content-length assertion. No local workaround was applied.

## Docling Eval 1.4.2

- repository commit: `f998acba3d33fe381edf597de602f7a1d3174c1a`;
- source archive SHA-256: `1f2193141cd99bca83883d2650845f36b3d655e6c7a28aec9f52ca10e8c49479` over `24438525` bytes;
- MIT `LICENSE` SHA-256: `56f9bf70349606c6fe1076b2661b96e3fdaf07710edc7dd16b230120ab80abff`;
- upstream `uv.lock` SHA-256: `9fbf62373240731761a5a6238c42a1650a75850a23a7c4f26ce24e5b488bbe08`.

The official frozen lock resolved `docling 2.83.0`, `docling-core 2.71.0`, `docling-ibm-models 3.13.0` and `docling-parse 5.6.2`; it is not the same runtime as Docling 2.122.0. With that frozen graph, the selected local suite returned 18 passed, 2 failed, 1 skipped and 3 deliberately deselected tests whose referenced example directory was absent. The two retained failures were an evaluator signature mismatch and a paged-image count mismatch. Therefore this release is useful evaluation code but is not promoted as an interchangeable runtime with Docling 2.122.0.

## Amazon Textractor 1.10.0

- repository commit: `8ea5f9ae65dbcbb2b83b409e0007c395e59b0e92`;
- source archive SHA-256: `4fa38999e355d977d2b46816b76c9d32ab6c5d8607987e086045470d101c7219` over `79080629` bytes;
- Apache-2.0 `LICENSE` SHA-256: `88046bf22d5b4f4b8cc85079ae6aae5424a3a1999db952ed152828ff325b2c6d`;
- `NOTICE` SHA-256: `27b12a9984d8e44b5ebea762d8c42423e57c76eadda0c7313e2c4decd443694e`.

The runtime requirements were insufficient to collect the official suite; `numpy`, `lxml` and `pandas` from the development environment were also required. With those present, 30 tests passed before a Windows CRLF representation assertion failed. Excluding only that test, 68 passed, 16 skipped and a word-ordering assertion failed. No live AWS call was made because `CALL_TEXTRACT` was unset. This supports local response parsing only; real extraction remains an `INTEGRATION_ONLY` service capability requiring an AWS account, region, costs and corpus evaluation.

## Admission decision

Docling 2.122.0 is the first executable local document-conversion candidate. Docling Eval remains a separately locked evaluator candidate. Textractor remains an optional official AWS integration. None of these results proves field-level correctness on REVESTEX documents; a project cannot store extracted business fields automatically until its supplied corpus and ground truth pass the project's declared thresholds.
