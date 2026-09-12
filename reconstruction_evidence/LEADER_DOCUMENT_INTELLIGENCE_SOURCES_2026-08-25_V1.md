# Leader document-intelligence sources evidence V1

Verified at: 2026-08-25.

## Scope

This evidence extends the earlier IBM Docling and AWS Textractor audit with exact public code from Microsoft Azure, Google Cloud, Amazon Web Services, NVIDIA and Microsoft Research. It proves archive identity, license identity and local static parsing where feasible. It does not claim a live provider call, production deployment or accuracy on private REVESTEX documents.

## Exact official archives

| Source ID | Commit/release | Bytes | Archive SHA-256 | License and SHA-256 |
|---|---|---:|---|---|
| `azure-document-intelligence-code-samples` | `8ac994c2712f3104d00ca4224b0ded6fd1ecbf70` | 57,555,887 | `018915d5ef37cba5fa3421a8716eff293127d5bc7492cafbf8cf141e6484f086` | MIT `d9a1b1e30d633d5732ea18e3cba9538d293ebc53e1a9e4e96ab739e0c5c4f1cb` |
| `azure-ai-document-processing-samples` | `9988025446845089204e03b4ec277e102bc8b970` | 8,643,533 | `afba8dd47fd2dd0687e86a53c357ef99c05592e39525fea013458d77dd80e532` | MIT `d9a1b1e30d633d5732ea18e3cba9538d293ebc53e1a9e4e96ab739e0c5c4f1cb` |
| `azure-document-intelligence-in-a-box` | `0ab9966b954e790bf17da6f6b8cc3652a480223e` | 8,071,866 | `a0d1296507d1ecb5c288af3b677c315de2706a3cccbec03a9d199c4080475625` | MIT `d9a1b1e30d633d5732ea18e3cba9538d293ebc53e1a9e4e96ab739e0c5c4f1cb` |
| `google-document-ai-samples` | `001ba391ab4a2f40d001cc0387618cb3c3699523` | 155,461,201 | `3131ff0604685967932cad1b689a64f7e50af1dc6b45524d99c5389342d34e8f` | Apache-2.0 `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30` |
| `aws-textract-document-extraction-platform` | `c7fed5353abcae184fbd705d5cf6a8512af0e9c4` | 1,124,136 | `52755b28d8afd2c8cee23ecba76b11c364edd3ecad066b2e4a51eb2f77219bb1` | Apache-2.0 `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30` |
| `aws-textract-queries-example` | `c7ee158086bbd2cad3d0b92d2a83895d1cb712c3` | 261,447 | `f50fde66f3340377fdb163af85fe1bc7da5c2a83c4ac28a81830032044f5d5f8` | MIT-0 `ef47b4ae2a1a8d38ef87b38dc5927967e7b472ffa3f7d125e89dbc8353a1e7af` |
| `nvidia-nemo-retriever-26.5.0` | release `26.5.0`, `0a6cb709b1ee4ca1d3f8c2fe0bbd1a0247f0aa09` | 10,209,364 | `b1c6856c83554cb5afda0b5bbe7c65d6f6866e40e365ee83c2dffc3c0c65d43d` | Apache-2.0 `c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4` |
| `microsoft-table-transformer-1.0.0` | release `v1.0.0`, `324edb4b9c99f0baf418f4c78ff263f93024b952` | 119,837 | `16200f2eb7eb9170812b3b95a890cf44fe4f297b3a5026880f0ccb981a0e282e` | MIT `c2cfccb812fe482101a8f04597dfc5a9991a6b2748266c47ac91b6a5aae15383` |
| `sap-btp-dox-invoice-validation` | `990f4d8650c12d39c9a84974091c1b248a549ab6` | 3,666,469 | `4dc23ccc7c1eef78513bf287cce7be2ea5de4a1a2ea38406b22ea02f8b472765` | Apache-2.0 `c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4` |
| `oracle-oci-python-sdk-2.185.0` | release `v2.185.0`, `e988c91dcc9963718454cb1e215bd44540524881` | 50,664,498 | `7584935da31a4af1898b709500bfa17bec96e68eb232785d999f93bd6f73ddda` | UPL-1.0 OR Apache-2.0 `8922f6a4bf38ae164a077bd7a294d690cb9e60826d1651d3b36cc7358379bf96` |
| `oracle-ai-invoice-handling` | `f1f52bb99cdb91c4e4e8332caba2c051fccd638f` | 728,626 | `c3875ae121442d5f3972b5c335a42f5c03ace18037b7ce8b0e88cad0e6eab4a5` | UPL-1.0 `c48d4b7187285146c57bf0b035658806be91b1ad65397e5e53ffa9f0f5acd6d9` |

## Local verification

The archives were downloaded from their official GitHub organizations and verified before inspection. Python AST parsing produced:

```text
Azure Document Intelligence code samples:      38 passed
Azure AI document processing samples:          23 passed
Azure Document Intelligence in a Box:            2 passed
Google Cloud Document AI samples:              117 passed
AWS document data extraction platform:          63 passed
AWS Textract Queries example:                   18 passed
Total:                                         261 passed
```

The source sets contained respectively 367, 145, 148, 1,021, 281 and 32 files. Static parsing proves readable Python source only; it is not a substitute for cloud credentials, deployed resources, model training or provider tests.

The materialized acquisition pack was also exercised on 2026-08-25 against two newly locked sources. It emitted `UPSTREAM_SOURCE_ACQUIRED` and exact `elite-official-source-receipt/v1` records for:

```text
aws-textract-queries-example
  commit  c7ee158086bbd2cad3d0b92d2a83895d1cb712c3
  archive f50fde66f3340377fdb163af85fe1bc7da5c2a83c4ac28a81830032044f5d5f8
  license MIT-0

microsoft-table-transformer-1.0.0
  commit  324edb4b9c99f0baf418f4c78ff263f93024b952
  archive 16200f2eb7eb9170812b3b95a890cf44fe4f297b3a5026880f0ccb981a0e282e
  license MIT
```

The wrapper shell subsequently reported a stale nonzero native exit code even though both PowerShell script operations completed and both receipts were parsed. This wrapper artifact is not recorded as an upstream failure; the two exact success records are the evidence.

After extending the canonical lock, the same materialized runner acquired all three enterprise additions and emitted exact receipts without a wrapper anomaly:

```text
UPSTREAM_SOURCE_ACQUIRED id=sap-btp-dox-invoice-validation
UPSTREAM_SOURCE_ACQUIRED id=oracle-oci-python-sdk-2.185.0
UPSTREAM_SOURCE_ACQUIRED id=oracle-ai-invoice-handling
```

Their receipt commits, archive hashes and license expressions exactly matched the table above. At this V1 snapshot, `UPSTREAM_ACQUISITION_TEST_PASS` and `UPSTREAM_LOCK_VALID sources=39 selected=39` passed from a clean pack materialization. The later BCQuality admission raises the global lock to 40 without changing these eleven document-source identities.

The AWS extraction platform's exact `pnpm-lock.yaml` SHA-256 is `240755c9c46f7d8e41985f50d141c7f2f93d15bbae461c595f21cd655b147e22`. Azure processing sample `requirements.txt` SHA-256 is `f483550b266285db7667653836720ebd68bdba70359fa8a89c99df3745e77f60`; AWS Queries example `requirements.txt` SHA-256 is `b31d5894495033bc8efc8e71677ad4a9a2643fff29568d678a3ba13cd0822990`.

## Preserved platform condition

NVIDIA NeMo Retriever could not be fully extracted with Windows archive tooling because its release archive contains a Unix-oriented `.agents` symlink/entry. Its intended execution also requires the supported Linux, GPU and NVIDIA service/model environment. The archive and license hashes passed; no Windows compatibility or local execution claim is made.

Microsoft Table Transformer is an exact public release and specialized table component. It was not promoted to a universal document extractor or a maintained cloud service.

Oracle OCI Python SDK release 2.185.0 contains 115 files under `src/oci/ai_document`, including `AIServiceDocumentClient.analyze_document`, invoice processor/model types and custom model operations. Of 115 AI Document Python files, 114 passed local AST parsing; one could not be opened from the deeply nested Windows temp path because of the platform path limit, not a Python syntax error. The full SDK contains 34,932 files and was not relabeled as a document-only package.

SAP's current invoice-validation sample contains 143 files, including 55 TypeScript files, CAP service/auth models, line-item/custom-schema extraction, correction/approval UI and Terraform. Its root and nested package manifests do not include an npm lockfile, so no reproducible install/build claim is made. It remains `SAMPLE_ONLY` despite providing a materially complete reference flow.

Oracle's invoice-handling archive contains 102 files and an importable Oracle Integration project. Thirty-two XML/XSD/WSDL files parsed successfully. Two `.json`-suffixed project artifacts use Oracle-specific/non-JSON content and were not accepted by a generic JSON parser; they are not presented as portable JSON.

The already locked `aws/aws-sdk-go-v2@284a4846e7bb941926e2d15144b17c01ffc98c1d` archive was reacquired through the materialized runner. Its `service/textract` module is release `v1.44.2` and contains 147 files, including `AnalyzeExpense`, asynchronous expense/document analysis, `QueriesConfig`, generated schemas and model-version fields. `go.mod` SHA-256 is `5b756eaf1b5931c7dab00afcab9b365041e9c4e052ca15cab6302d028b7f3d9c`; `go.sum` SHA-256 is `0c350d80c4b68d97f362d390afb476b621ed53550c9a88e3b9cb8e322c1e61f2`. The current project preflight reports Go unavailable, so the module was inspected but not compiled in this pass.

Official PyPI wheels `azure-ai-documentintelligence==1.0.2` (`e1fb446abbdeccc9759d897898a0fe13141ed29f9ad11fc705f951925822ed59`) and `google-cloud-documentai==3.15.0` (`f040f4f9db43411184197a808b11fde52b580723a9fca336f43ea7c6a885bfd8`) were installed into a new Python virtual environment. `DocumentIntelligenceClient` and `DocumentProcessorServiceClient` imported with their exact expected distribution versions. This proves package installation/import, not account access or extraction accuracy. The full observed transitive graph is preserved in `markdown_system/OFFICIAL_DOCUMENT_SDK_ARTIFACT_LOCK.md` and must be hash-locked per target before deployment.

## Admission decision

At this V1 snapshot the source lock contained 39 exact archives; it now contains 40 after a non-document BCQuality addition. These eleven document additions provide official code for prebuilt invoice/receipt integration, custom extraction, classification, document splitting, query adapters, review/metrics, structured content, tables and invoice-to-ERP workflows. They are immediately acquirable and inspectable. Live extraction becomes executable only after the chosen provider account, region, identity/scopes and project resources exist. Proforma invoices, packing lists, purchase orders, bills of lading and other business-specific forms remain custom-model candidates until representative documents, schemas and approved ground truth exist.
