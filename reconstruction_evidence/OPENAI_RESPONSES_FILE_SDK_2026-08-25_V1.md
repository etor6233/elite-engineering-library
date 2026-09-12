# OpenAI Responses file SDK evidence V1

Verified on 2026-08-25 against [official OpenAI Responses API documentation](https://developers.openai.com/api/reference/cli/resources/responses/methods/create) and the exact Python wheel.

| Property | Exact value |
|---|---|
| package | `openai==3.3.1` |
| release upload | 2026-08-19 |
| wheel | `openai-3.3.1-py3-none-any.whl` |
| bytes | `1690337` |
| SHA-256 | `9652df7fdf8ee6f5bd58e0a12f2b1d414a18e0f06bb7a9a57c8643a5f5469bd3` |
| Python | `>=3.10` |
| license | Apache-2.0 from wheel `License-Expression` |

The exact wheel installed in an isolated Python 3.14 environment. `OpenAI`, `AsyncOpenAI` and `ResponseInputFileParam` imported, package version equality passed, and the Responses client surface was present.

Official documentation establishes that Responses accepts text, image or file inputs; supports structured JSON outputs; and can use built-in file search/code interpreter, MCP tools and typed function calls. These are real provider capabilities useful for agentic ingestion and secondary extraction/evaluation.

Admission is `INTEGRATION_ONLY / CONDITIONED`. The wheel is immediately installable, but live calls require an OpenAI project/API key, approved model/data-retention policy, cost/region review and evaluation against the project's documents. It is not a field-confidence/provenance engine and does not replace Azure/Google/AWS document processors or authorize automatic storage of invoice, proforma, packing-list, purchase-order, bill-of-lading or customs fields.
