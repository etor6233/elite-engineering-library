# Official Agent and Merchant Sources — 2026-08-25 V1

## Scope

This evidence covers three official source archives added to `OFFICIAL-UPSTREAM-ACQUISITION-CORE` 0.1.1. Acquisition used immutable GitHub commit archives; no Git repository was created.

## Exact sources

| Source | Commit | Archive bytes | Archive SHA-256 | License evidence |
|---|---:|---:|---|---|
| `aws-samples/sample-well-architected-skills-and-steering` v6.3.1 | `4b58ba02670abcd67458eeafb44fb6118b6af2ef` | 8,078,075 | `933c3ce6e467f6c102eb1e92228abbbc0b7f0091690af48251f3312d3fc1b097` | MIT-0 `LICENSE` `6f829c66e5130120d139a1a9011cd0952bb7197ef78baa8c495fc99acbe37685`; `NOTICE` `1db5530fc81aa3762e76516f3f446f85804a309b5ea64bf79130a649abb9d40b` |
| `google/skills` | `7b596c6e315d8661999c2b91e8139171c6c6330d` | 2,430,289 | `fe9ce520456d32357e76b169dc1da7f83d127a9202813cd7309071f5ddd2fba9` | Apache-2.0 `LICENSE` `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30` |
| `google/merchant-api-samples` | `e622d99e9f61504cff537853321d5faf51e13a69` | 1,367,243 | `525245e193d53b7160cdf7e638954735b7e75e94de248e668420507027938830` | Apache-2.0 `LICENSE` `6dc0e068dcf3a5bc8e054205b85b7720e1d49265bbc64bf515d2cf79197df69a` |

The archives contained 2,059, 928 and 773 files respectively. Hashes were computed over the downloaded source ZIPs and root license/notice files before admission.

## Executed verification

- The canonical acquisition lock validated all 25 entries and all negative lock tests passed.
- Google Merchant's 206 Python sample files passed Python 3.12 bytecode compilation after using a short Windows audit path. The first attempt retained a Windows long-path write failure affecting `__pycache__`; it was an audit-path failure rather than a syntax failure.
- The local environment did not provide Go, so the Go samples were not built. No live Merchant API call was attempted because no Google Cloud project, Merchant Center account, developer registration or credential was supplied.
- The Google Merchant archive contains official samples for Java, Python, PHP, Go, Node.js, .NET and Apps Script plus the official `mapi-developer-assistant` agent skill.

## Admission result

- AWS and Google agent assets: `PINNED_CANDIDATE`. They may guide reviews or provider workflows from the exact source, but they do not prove the generated product.
- Google Merchant samples: `INTEGRATION_ONLY`. They provide official executable integration starting points; production use remains blocked on real account access, contract tests, reconciliation and project-specific policy.
- None of these sources supplies REVESTEX business rules, a ledger, Argentine fiscal behavior or a guarantee of document-field accuracy.
