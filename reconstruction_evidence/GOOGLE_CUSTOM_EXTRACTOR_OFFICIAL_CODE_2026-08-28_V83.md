# Google Custom Document Extractor official code — V83

Date: 2026-08-28  
Scope: replace document-type-by-document-type sample hunting with an exact official implementation that accepts an extraction schema per request, while preserving the existing AWS and Microsoft configurable engines.

## Official authority and identity

- Repository: [GoogleCloudPlatform/python-docs-samples](https://github.com/GoogleCloudPlatform/python-docs-samples/tree/dc0eecc2187791fea70f48201241288e25fc9620/documentai/snippets).
- Default-branch head observed: `dc0eecc2187791fea70f48201241288e25fc9620`.
- GitHub verification: `verified=true`, reason `valid`.
- Repository license: Apache-2.0; exact `LICENSE` is 11,357 bytes, SHA-256 `c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4`.
- `documentai/snippets/handle_response_sample.py`: Git blob `58bbb1debe09c6fad3814789f3db78af4679e061`, 20,533 bytes, SHA-256 `8fc1adaceba3ac8ad4d871f8003b2ce43fbadfd9cc7b80599cefcc6edc8bfb91`.
- `documentai/snippets/handle_response_sample_test.py`: Git blob `b7c65834ccafbcb8fbae6f46807d0a7184ff1048`, 7,859 bytes, SHA-256 `a60d4d9a53aa3c6383d60bb2cfc22489fde622d09ac994d68161ee771b1fea45`.
- Google documents this same path in [Handle processing response](https://cloud.google.com/document-ai/docs/handle-response), whose code samples are Apache-2.0 under the Google Developers Site Policies.

## What the exact Google code provides

The unmodified `process_document_custom_extractor_sample` constructs:

1. `DocumentSchema.EntityType.Property` values, including required/optional occurrence;
2. one `DocumentSchema.EntityType` based on `document`;
3. `ProcessOptions(schema_override=DocumentSchema(...))`;
4. a normal `ProcessRequest` through `DocumentProcessorServiceClient`;
5. entity type, text/mention, confidence, normalized value and nested properties from the response.

The schema is attached to each request. Therefore a project is not limited to invoice, packing-list, BOL, proforma or customs code branches: the selected fields/classes are processor input. Google explicitly describes Custom Extractor as the path for document types without a pretrained processor, and lists bill of lading among variable-layout examples in the [Custom extractor overview](https://cloud.google.com/document-ai/docs/custom-extractor-overview).

This is not an Elite-authored extraction algorithm. The two runtime files are `VERBATIM` Google code. The local contract is only a hash-bound verification harness and is declared `AUTHORED`.

## Rebuild and execution evidence

Pack `GOOGLE-CLOUD-DOCUMENT-AI-OFFICIAL-PROCESS-SAMPLE` was promoted from 0.1.0 to 0.2.0:

- 13 files materialized from Markdown;
- five Google `VERBATIM` blocks: license, two official samples and two official live tests;
- eight local guide/approval/acquisition/validation/test files with no Google authorship claim;
- exact materialized hashes for the two new Google files equal the official raw hashes above;
- both official files and the local contract compile under Python;
- with exact installed `google-cloud-documentai==3.15.0`, the unmodified function executed against a capture boundary and emitted `GOOGLE_CUSTOM_EXTRACTOR_OFFICIAL_SAMPLE_OFFLINE_PASS properties=3 schema_override=1`;
- existing process_document regression remained 3/3 PASS;
- fixture acquisition remained 2 positive/4 negative PASS.

Profiles after composition:

- minimal Google official-code profile: 2 packs / 19 files;
- complete Google document lane: 8 packs / 81 files;
- all files are available immediately to the agent after composition; live network effects remain approval-gated.

## Current leader-source comparison

- Google `document-ai-samples` remains at verified head `001ba391ab4a2f40d001cc0387618cb3c3699523`; the new generic CDE code comes from the current signed `python-docs-samples` head instead of the deprecated split workflow.
- AWS GenAI IDP Accelerator latest stable release remains v0.6.5 at `1b5fd74454e593de233342a02ee911af8ee38359`, already exact-locked. Its official code supplies multi-document discovery, clustering, JSON schema generation, evaluation/ground truth, confidence, validation and review. Post-release `main` was observed at unsigned `720f3052a57c22ae0d1e50886965c065ee86d60b` and was not substituted for the stable release.
- Microsoft Content Processing Solution Accelerator remains at verified head `659eaa1f503dd08b1e1aea1c72eab11c7c191d00`, already exact-locked and previously audited for Extract→Map→Evaluate→Save.

No unsigned or unreleased revision replaced a signed/released source.

## Preserved condition

The upstream Google test requires `GOOGLE_CLOUD_PROJECT`, ADC, an authorized Custom Extractor processor/version, quota and cost. Its repository requirements currently pin `google-cloud-documentai==3.0.1`; the library does not install that historical manifest. The code instead passed the local compatibility contract with the exact official 3.15.0 artifact already admitted. Audit mode verifies source identity and syntax without network; the executable runtime gate runs inside the frozen Google venv only when network/cache use is authorized.

This condition is recorded as `UP-FAIL-172`; it does not convert the official source into local code and does not assert a live result that was not executed.

## Governing gates

After the change and before this evidence file:

- `VERIFY_LIBRARY_PASS`: 59 packs, 551 materializable files, 393 Markdown files;
- composition: Google runtime 81, Google official sample profile 19;
- `VERIFY_EXECUTABLE_LIBRARY_PASS mode=Audit`: 59 packs, 111 upstream sources, 12 document SDK artifacts;
- exact Google Custom Extractor hashes and three Python syntax outputs PASS offline;
- independent SDK 3.15.0 runtime contract PASS;
- failure memory: 934 local rows + 172 upstream conditions = 1,106 unique IDs, zero local rows open.

The next project can select the Google lane and define arbitrary business document schemas without waiting for a separate code sample per document name. Provider access, business schema and project evidence remain explicit inputs rather than invented defaults.
