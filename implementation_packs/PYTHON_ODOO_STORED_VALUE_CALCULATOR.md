# Exact Odoo-derived stored-value calculator

## 1. Metadata

```yaml
pack_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Hash-bound same-currency gift issuance/redemption and loyalty earn/burn calculations, derived from selected Odoo19 methods; complete corresponding module source and isolated JSON process"
stacks: ["CPython 3.14.4", "Go 1.26.8", "PostgreSQL 18.6", "Next.js 16.3.4 for selected portal"]
compatible_with: ["MARKDOWN-COMPOSITOR 0.3.0", "V402 connected franchise profile", "explicit hash-bound reference program"]
incompatible_with: ["public bearer-wallet inferred from assisted internal account IDs", "arbitrary financial or tax policy", "production certification inferred from fixtures"]
license_expression: "LGPL-3.0-only AND LicenseRef-Workspace-Owner"
upstream_sources: ["https://github.com/odoo/odoo/tree/99edb6dd82b7b560930c00b03b694ba700785370"]
verified_at: "2026-09-12"
```

## 2. Applicability

Select both stored-value packs only with the connected commerce/payment/handover,
shared human-approval and portal owners in FRANCHISE_COMPLETE_PACK_PLAN. A named
reference profile supplies rules explicitly; it does not decide the target program,
tax, wage or accounting policy. No public gift-card code service is claimed.

## 3. Architecture contract

The LGPL calculator owns only the exact mapped commercial calculations. CPython
runs -I -S -B and executes verified source snapshots with no filesystem package
imports. AUTHORED Go glue owns scoped principals, locked order/account snapshots,
immutable approvals, idempotency, outbox and PostgreSQL receipts. Separate reviewers
approve bound effects. Deferred commit constraints reject expired operations after
waits; reservations and order locks prevent double spend. Receipt history and
current eligibility are separate. Full coverage creates an immutable local funding
receipt, never a zero-value provider charge; handover explicitly selects profile3.
BFF transports int64 and exact payload/receipt bytes as strings, bounds body/time,
and recovers persisted request keys after uncertain responses.

## 4. Exact file manifest

```text
CREATE docs/provenance/ODOO_STORED_VALUE_NOTICES.md
CREATE docs/provenance/ODOO_STORED_VALUE_RUNTIME_LOCK.json
CREATE docs/provenance/ODOO_STORED_VALUE_SOURCE_LOCK.json
CREATE odoo_loyalty/COPYRIGHT
CREATE odoo_loyalty/DERIVATION.json
CREATE odoo_loyalty/DERIVATION_CONNECTED.json
CREATE odoo_loyalty/LICENSE
CREATE odoo_loyalty/engine-lock.json
CREATE odoo_loyalty/engine.py
CREATE odoo_loyalty/oracle_loader.py
CREATE odoo_loyalty/orm_contract.py
CREATE odoo_loyalty/protocol.py
CREATE odoo_loyalty/run.py
CREATE odoo_loyalty/upstream/COPYRIGHT
CREATE odoo_loyalty/upstream/LICENSE
CREATE odoo_loyalty/upstream/addons/loyalty/__manifest__.py
CREATE odoo_loyalty/upstream/addons/loyalty/models/loyalty_card.py
CREATE odoo_loyalty/upstream/addons/loyalty/models/loyalty_history.py
CREATE odoo_loyalty/upstream/addons/loyalty/models/loyalty_program.py
CREATE odoo_loyalty/upstream/addons/loyalty/models/loyalty_reward.py
CREATE odoo_loyalty/upstream/addons/loyalty/models/loyalty_rule.py
CREATE odoo_loyalty/upstream/addons/loyalty/tests/test_loyalty.py
CREATE odoo_loyalty/upstream/addons/sale_loyalty/__manifest__.py
CREATE odoo_loyalty/upstream/addons/sale_loyalty/models/loyalty_card.py
CREATE odoo_loyalty/upstream/addons/sale_loyalty/models/loyalty_program.py
CREATE odoo_loyalty/upstream/addons/sale_loyalty/models/sale_order.py
CREATE odoo_loyalty/upstream/addons/sale_loyalty/models/sale_order_coupon_points.py
CREATE odoo_loyalty/upstream/addons/sale_loyalty/models/sale_order_line.py
CREATE odoo_loyalty/upstream/addons/sale_loyalty/tests/common.py
CREATE odoo_loyalty/upstream/addons/sale_loyalty/tests/test_buy_gift_card.py
CREATE odoo_loyalty/upstream/addons/sale_loyalty/tests/test_loyalty.py
CREATE odoo_loyalty/upstream/addons/sale_loyalty/tests/test_loyalty_history.py
CREATE odoo_loyalty/upstream/addons/sale_loyalty/tests/test_pay_with_gift_card.py
CREATE odoo_loyalty/upstream/addons/sale_loyalty/tests/test_program_multi_company.py
CREATE odoo_loyalty/upstream/addons/sale_loyalty/tests/test_program_numbers.py
CREATE odoo_loyalty/upstream/addons/sale_loyalty/tests/test_program_rules.py
CREATE odoo_loyalty/upstream/odoo/tools/float_utils.py
CREATE test_odoo_loyalty_derived.py
CREATE test_stored_value_source_loader.py
CREATE tools/rebuild_stored_value_source.py
```

## 5. Materialization blocks

### FILE: `docs/provenance/ODOO_STORED_VALUE_NOTICES.md`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file1:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "863294718fca6559b257953df9040da06d9f1dc87cb24776c4f21ec2d900ccd0"
variables: []
secrets_allowed: false
```

````markdown
# Odoo-derived gift card and loyalty calculation

The isolated Python module contains code adapted from Odoo Community19.0 at commit99edb6dd82b7b560930c00b03b694ba700785370. Original copyright is preserved verbatim in odoo_loyalty/COPYRIGHT and upstream/COPYRIGHT. The complete LGPLv3 and incorporated GPLv3 texts are odoo_loyalty/LICENSE and upstream/LICENSE. The24complete upstream files retain their headers and exact bytes. ODOO_STORED_VALUE_SOURCE_LOCK.json identifies official URLs, Git blobs, sizes and SHA256; this is a selected-file snapshot, not the complete ERP or a signed release claim.

engine.py is ADAPTED, not VERBATIM. DERIVATION.json and DERIVATION_CONNECTED.json identify original functions/segments and the changes: bounded same-currency inputs, Decimal representation, ORM contract replacement, monetary rounding and append-only local reversal receipts instead of Odoo history deletion. Four original focal function ASTs remain unchanged between the oracle-tested derivation and the connected module. No tax, full ORM, eWallet, currency conversion, coupon code issuance or complete Odoo payment workflow is claimed.

The authored Python contract, launcher and recordset glue are shipped as part of this LGPL-3.0-only module with corresponding source. Go/SQL/HTTP/portal code is separately declared AUTHORED workspace glue; none of it is attributed to Odoo, Microsoft, Google, xAI, SpaceX or another external company. Test fixtures derived from the official Odoo tests are identified as ADAPTED and retain LGPL-3.0-only distribution. Process separation does not by itself decide the legal scope of a combined distribution.

Keep the license/copyright texts, all distributed corresponding module source, derivation records, fixture/oracle sources and tools when redistributing this module or a modified version. Identify your modifications and preserve the recipient's applicable LGPL/GPL rights, including modification and debugging of modifications. No contract here restricts those rights. Consult the included license for the conditions of the chosen distribution; this artifact is not a legal certification of a future product or its other components.

Source replacement is supported: modify a separate source copy, then run tools/rebuild_stored_value_source.py with that module's baseline lock SHA and an absent destination outside it. It emits a complete runtime source snapshot, new manifest and explicit USER_MODIFIED_SOURCE_NOT_ADMITTED receipt. Original license/copyright texts cannot be removed by this helper. Review and test the modification, then explicitly configure the new launcher and manifest hashes. Prior provenance, gates and upstream authorship are never inherited by a local modification. The source-loading regression proves that the original hash rejects a modified snapshot and its explicitly selected new hash executes it. The canonical library is never overwritten by this tool.
````

### FILE: `docs/provenance/ODOO_STORED_VALUE_RUNTIME_LOCK.json`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file2:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "c168b42c0e7aed14d8370b959147d22828d43e88861ee27c3c007debd99d61e9"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite.stored-value-activation-lock.v1",
  "source_commit": "99edb6dd82b7b560930c00b03b694ba700785370",
  "files": {
    "odoo_loyalty/engine-lock.json": "d39918d5c4479af0d961fa9b3b61f3fb8b00e997d03be958083debad627657da",
    "odoo_loyalty/run.py": "2101f3e889607a8da2fd22722d64c20fd481fa2c38882e1495aeffd72bb3387e",
    "deploy/stored-value/profile.reference.json": "db96a3558474e1dd57448d4fd25af9753cdef7238faabac13b1f17727994e727"
  },
  "python": "Project-admitted CPython runtime, absolute path; launcher uses -I -S -B; reference tested3.14.4",
  "activation": "Explicit STORED_VALUE_ENABLED=true; tenant and organization must match the selected profile; source/profile hashes are configuration bindings after verification, not proof of provenance by themselves"
}
````

### FILE: `docs/provenance/ODOO_STORED_VALUE_SOURCE_LOCK.json`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file3:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "8aecf2bf35bfc71a01b5db0d817855fbf7f3d40cb3a919276f62c6b775a03152"
variables: []
secrets_allowed: false
```

````json
{
  "observed_at": "2026-09-12T01:20:41.145419+00:00",
  "profile_sha256": "71fbfa3efb26f442888abead01835f60129e7d06c90241b4b7b525bc12e7d966",
  "lock_sha256": "6e6e46a02b2d9980a214a8474ec4204782189f7e2b9dfd32059d62a3e6008760",
  "commit_metadata_sha256": "71e9cc57d0af499b9422e899c391aff14cd95c2361f346019be279583b3b7c4b",
  "files": [
    {
      "path": "COPYRIGHT",
      "source_url": "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/COPYRIGHT",
      "revision": "99edb6dd82b7b560930c00b03b694ba700785370",
      "git_blob": "b0f9743d204484b1080d339ec55a638b1cde32dc",
      "sha256": "86e49232d2162708d05405ed5ff6dc5594b73ef4b0f25d2749cfae13493b620c",
      "bytes": 433
    },
    {
      "path": "LICENSE",
      "source_url": "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/LICENSE",
      "revision": "99edb6dd82b7b560930c00b03b694ba700785370",
      "git_blob": "12334e0fdd8d3412d5e1deed72220fcc03431ff2",
      "sha256": "abc09dad5f84a76e1b0279237053cae16c03228ab27d8d467677054c2bd17eeb",
      "bytes": 43529
    },
    {
      "path": "addons/loyalty/__manifest__.py",
      "source_url": "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/loyalty/__manifest__.py",
      "revision": "99edb6dd82b7b560930c00b03b694ba700785370",
      "git_blob": "fc737c84bb2c46336609a58a5ece87ad77432d65",
      "sha256": "03537054d1786db61c364262242eff900c6b104c2e1e5c8df3a60cd2b1e403f5",
      "bytes": 1754
    },
    {
      "path": "addons/loyalty/models/loyalty_card.py",
      "source_url": "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/loyalty/models/loyalty_card.py",
      "revision": "99edb6dd82b7b560930c00b03b694ba700785370",
      "git_blob": "f0baa0b8bcbe434ff5777c48d2c8ef8a68003a3d",
      "sha256": "9eead538d45b3ad32a0f2eddc6b146d2a0ab5cc4d40de19bdea58aa70908ba83",
      "bytes": 8800
    },
    {
      "path": "addons/loyalty/models/loyalty_history.py",
      "source_url": "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/loyalty/models/loyalty_history.py",
      "revision": "99edb6dd82b7b560930c00b03b694ba700785370",
      "git_blob": "5c70b44be15f623fd1e897e6dc30b50e35560f33",
      "sha256": "c480709def3ad815fe6b306dac534785f24cf10b7432276a0276dcd5646bb8a5",
      "bytes": 885
    },
    {
      "path": "addons/loyalty/models/loyalty_program.py",
      "source_url": "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/loyalty/models/loyalty_program.py",
      "revision": "99edb6dd82b7b560930c00b03b694ba700785370",
      "git_blob": "239f33d03fe89ab14885d6ef97b96fe8267bbeb3",
      "sha256": "3ed5484b4404120c73afe810653c2b25915d50d4211158005f684d9290e08096",
      "bytes": 27848
    },
    {
      "path": "addons/loyalty/models/loyalty_reward.py",
      "source_url": "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/loyalty/models/loyalty_reward.py",
      "revision": "99edb6dd82b7b560930c00b03b694ba700785370",
      "git_blob": "08ea7a9cf5c319931b32864413fecbbad2ec1427",
      "sha256": "ec0aa37a9a5f12c01e4a860738afcc1b6ede113e4a3f4aa9d08966194936955c",
      "bytes": 15597
    },
    {
      "path": "addons/loyalty/models/loyalty_rule.py",
      "source_url": "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/loyalty/models/loyalty_rule.py",
      "revision": "99edb6dd82b7b560930c00b03b694ba700785370",
      "git_blob": "6ee64ff2055640f7b9a012338b4e93faa1db8570",
      "sha256": "955eafbc11279154ba30b4986fa2198c683d4c4519253f5ee0fd84daa39f2a66",
      "bytes": 6779
    },
    {
      "path": "addons/loyalty/tests/test_loyalty.py",
      "source_url": "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/loyalty/tests/test_loyalty.py",
      "revision": "99edb6dd82b7b560930c00b03b694ba700785370",
      "git_blob": "b7acc32a1c1b8e9667a9c2d130b42475cb0e0248",
      "sha256": "39c7660a20b402f80fafab7e6534c9bbcd2e52abd847088c0263004fe33092c3",
      "bytes": 14565
    },
    {
      "path": "addons/sale_loyalty/__manifest__.py",
      "source_url": "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/__manifest__.py",
      "revision": "99edb6dd82b7b560930c00b03b694ba700785370",
      "git_blob": "bb0716a9637eca91dca26385f9bac21848479726",
      "sha256": "284a427b610067bac9c82ae6e3cac37cb1ad51c3dee5617fc5b93d62b13ee083",
      "bytes": 1039
    },
    {
      "path": "addons/sale_loyalty/models/loyalty_card.py",
      "source_url": "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/models/loyalty_card.py",
      "revision": "99edb6dd82b7b560930c00b03b694ba700785370",
      "git_blob": "9a52e86399df0af0cf44ea79046fd7c59558498c",
      "sha256": "67a9501daf1b48bffbf10e34d11b7503df920a0048545affd38e86f08f2be712",
      "bytes": 2147
    },
    {
      "path": "addons/sale_loyalty/models/loyalty_program.py",
      "source_url": "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/models/loyalty_program.py",
      "revision": "99edb6dd82b7b560930c00b03b694ba700785370",
      "git_blob": "d8d8e3c40fed4ea46e592d79dbadda8c62bae27a",
      "sha256": "2f1bbcfb4419d8bd6149eac02e963670c342470dd2007eaa02e31b095704e658",
      "bytes": 1030
    },
    {
      "path": "addons/sale_loyalty/models/sale_order.py",
      "source_url": "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/models/sale_order.py",
      "revision": "99edb6dd82b7b560930c00b03b694ba700785370",
      "git_blob": "876798076f76b3d4b6164ef02eb098d7ac7dd478",
      "sha256": "6de1aeed8f1fc80268025f13bff2a56b57211b6a82ce9df950a8b08a623e3a06",
      "bytes": 78182
    },
    {
      "path": "addons/sale_loyalty/models/sale_order_coupon_points.py",
      "source_url": "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/models/sale_order_coupon_points.py",
      "revision": "99edb6dd82b7b560930c00b03b694ba700785370",
      "git_blob": "b0b50be60e2c61182876e6d0fa9a924029ebbcf0",
      "sha256": "81a49606130beed878aaa93c503970c8e4e9b4645b63f2b097999ccb99c4eeb2",
      "bytes": 676
    },
    {
      "path": "addons/sale_loyalty/models/sale_order_line.py",
      "source_url": "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/models/sale_order_line.py",
      "revision": "99edb6dd82b7b560930c00b03b694ba700785370",
      "git_blob": "0e97738849be4abea95ea107c04c7566d5c58fe3",
      "sha256": "4aba6673de316b9cb78def3c1d4e31d456ce170524c84d5f8dca93b2a8f18d93",
      "bytes": 7181
    },
    {
      "path": "addons/sale_loyalty/tests/common.py",
      "source_url": "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/tests/common.py",
      "revision": "99edb6dd82b7b560930c00b03b694ba700785370",
      "git_blob": "41c5dca474739038482d3e2f77a25bd88b05e5a7",
      "sha256": "ab2f2ef3334a82a642178427f24e8316e3ee2d4f91c63dede8712304d3c53db5",
      "bytes": 12605
    },
    {
      "path": "addons/sale_loyalty/tests/test_buy_gift_card.py",
      "source_url": "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/tests/test_buy_gift_card.py",
      "revision": "99edb6dd82b7b560930c00b03b694ba700785370",
      "git_blob": "5818de14efc9fde3101f0714cb446c13692d3bd3",
      "sha256": "f724115af4e50776f152bb933b01999bef7e800aec4eea0c43ad38f9916bb3e8",
      "bytes": 3202
    },
    {
      "path": "addons/sale_loyalty/tests/test_loyalty.py",
      "source_url": "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/tests/test_loyalty.py",
      "revision": "99edb6dd82b7b560930c00b03b694ba700785370",
      "git_blob": "0eb9ddcd8e035bc61b355744ae66db1040bfe127",
      "sha256": "b2e97e780f49b8303b1930cdc7b6687338bcdf0cd48039e747144eddd14db47e",
      "bytes": 55159
    },
    {
      "path": "addons/sale_loyalty/tests/test_loyalty_history.py",
      "source_url": "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/tests/test_loyalty_history.py",
      "revision": "99edb6dd82b7b560930c00b03b694ba700785370",
      "git_blob": "567039fe0292fd2c056c0c82e97781ad15934286",
      "sha256": "291cee97246ffe07b39d77b4c0bdaa3c3addb7cec321aa2235827ab92515775d",
      "bytes": 8005
    },
    {
      "path": "addons/sale_loyalty/tests/test_pay_with_gift_card.py",
      "source_url": "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/tests/test_pay_with_gift_card.py",
      "revision": "99edb6dd82b7b560930c00b03b694ba700785370",
      "git_blob": "13c06bbf7ae046d55c2ac3c183663c31840f0d85",
      "sha256": "9ecd12f97ddb7dcdcb8147d38752a1155c21583977e2ce048e40e653b4468f47",
      "bytes": 10659
    },
    {
      "path": "addons/sale_loyalty/tests/test_program_multi_company.py",
      "source_url": "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/tests/test_program_multi_company.py",
      "revision": "99edb6dd82b7b560930c00b03b694ba700785370",
      "git_blob": "c06a430d127004821e968ab58ff3dd0666a5499c",
      "sha256": "d338ab072168bacfa342f7a5a0011063d47f41f19ef1221a4adf0a70cb90a32b",
      "bytes": 4963
    },
    {
      "path": "addons/sale_loyalty/tests/test_program_numbers.py",
      "source_url": "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/tests/test_program_numbers.py",
      "revision": "99edb6dd82b7b560930c00b03b694ba700785370",
      "git_blob": "45d08547f109e952daa2b268ff298727c0d1ddac",
      "sha256": "845f5b205b1ac941280905bfa308ed15d7142d3069997b6545f5cc0e46b0c918",
      "bytes": 97516
    },
    {
      "path": "addons/sale_loyalty/tests/test_program_rules.py",
      "source_url": "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/tests/test_program_rules.py",
      "revision": "99edb6dd82b7b560930c00b03b694ba700785370",
      "git_blob": "ce4e57595c3e04d307b5f04db22ff6746bda1dca",
      "sha256": "f28fa822aa5d7116816f20f2010c2c3f94f926477de024b768ed52d44c3cbce1",
      "bytes": 21582
    },
    {
      "path": "odoo/tools/float_utils.py",
      "source_url": "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/odoo/tools/float_utils.py",
      "revision": "99edb6dd82b7b560930c00b03b694ba700785370",
      "git_blob": "123a8b0409a64f48e5b45875c3420bff0f32d657",
      "sha256": "ab7d51a44033792414dc63f598517ee76b922f76b2de9fa3eeb977ab78d7a962",
      "bytes": 19196
    }
  ],
  "status": "INSPECTED_NOT_ADMITTED",
  "executed": false,
  "installed": false,
  "full_archive_acquired": false,
  "scope": "24 selected complete official files; immutable commit, not full Odoo runtime",
  "acquisition_receipt_sha256": "7e581370505c7b28debdb8da31c61a1963d8a4eeae42930ddee3b2c358e658b5",
  "source_root": "odoo_loyalty/upstream"
}
````

### FILE: `odoo_loyalty/COPYRIGHT`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file4:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/COPYRIGHT; revision 99edb6dd82b7b560930c00b03b694ba700785370; Git blob b0f9743d204484b1080d339ec55a638b1cde32dc"
license: "LGPL-3.0-only"
sha256: "86e49232d2162708d05405ed5ff6dc5594b73ef4b0f25d2749cfae13493b620c"
variables: []
secrets_allowed: false
```

````text

Most of the files are

  Copyright (c) 2004-2015 Odoo S.A.

Many files also contain contributions from third
parties. In this case the original copyright of
the contributions can be traced through the
history of the source version control system.

When that is not the case, the files contain a prominent
notice stating the original copyright and applicable
license, or come with their own dedicated COPYRIGHT
and/or LICENSE file.

````

### FILE: `odoo_loyalty/DERIVATION.json`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file5:v1"
operation: CREATE
provenance: AUTHORED
source: "local source manifest, isolated loader, ORM/protocol and oracle glue distributed as part of the modified LGPL module; no Odoo corporate authorship"
license: "LGPL-3.0-only"
sha256: "26d841355c73e08cd9e99fd855c9fba0a69b6681b53b23275621c4741fa58fe2"
variables: []
secrets_allowed: false
```

````json
{
  "revision": "99edb6dd82b7b560930c00b03b694ba700785370",
  "source_receipt_sha256": "7e581370505c7b28debdb8da31c61a1963d8a4eeae42930ddee3b2c358e658b5",
  "license": "LGPL-3.0-only",
  "functions": [
    {
      "name": "_get_point_changes",
      "source_path": "addons/sale_loyalty/models/sale_order.py",
      "source_file_sha256": "6de1aeed8f1fc80268025f13bff2a56b57211b6a82ce9df950a8b08a623e3a06",
      "source_lines": [
        768,
        781
      ],
      "source_segment_sha256": "9c813522e920d7880dd7eb5a5a0e8b831893ca48e6e6d9cd0c72c5bc6fe6728b",
      "adaptation": "Dedent methods; binary 0.0 literals become Decimal; same control flow. ORM protocol is supplied by declared local glue; float_round binding is exact Decimal rounding."
    },
    {
      "name": "_get_real_points_for_coupon",
      "source_path": "addons/sale_loyalty/models/sale_order.py",
      "source_file_sha256": "6de1aeed8f1fc80268025f13bff2a56b57211b6a82ce9df950a8b08a623e3a06",
      "source_lines": [
        783,
        799
      ],
      "source_segment_sha256": "308749a1cd0fa23fea6de627dcaa88f782ebffca21e69902cdc15c430edf8162",
      "adaptation": "Dedent methods; binary 0.0 literals become Decimal; same control flow. ORM protocol is supplied by declared local glue; float_round binding is exact Decimal rounding."
    },
    {
      "name": "_program_check_compute_points",
      "source_path": "addons/sale_loyalty/models/sale_order.py",
      "source_file_sha256": "6de1aeed8f1fc80268025f13bff2a56b57211b6a82ce9df950a8b08a623e3a06",
      "source_lines": [
        1238,
        1375
      ],
      "source_segment_sha256": "b5d16594d709a313ee56758a963818c34b070b39ba1e6cc8db95ab9c845f73ee",
      "adaptation": "Dedent methods; binary 0.0 literals become Decimal; same control flow. ORM protocol is supplied by declared local glue; float_round binding is exact Decimal rounding."
    },
    {
      "name": "reward_from_bound_amounts",
      "source_path": "addons/sale_loyalty/models/sale_order.py",
      "source_file_sha256": "6de1aeed8f1fc80268025f13bff2a56b57211b6a82ce9df950a8b08a623e3a06",
      "source_lines": [
        579,
        601
      ],
      "source_segment_sha256": "d26d2bdf390118da818ba9fcfc052d51ac6cf68f123bac93bf94e300b3effa17",
      "adaptation": "Same-currency preconditions and discountable amount are supplied by bound order owner; only upstream cap/cost block is adapted. Infinity becomes Decimal; return serializes existing values. No tax calculation, FX, discount allocation, random code or fiscal product mapping copied."
    }
  ],
  "runtime_glue": [
    "orm_contract.py",
    "protocol.py"
  ],
  "not_claimed": [
    "Odoo ORM/runtime integration suite",
    "fiscal rules",
    "foreign currency conversion",
    "full Odoo payment workflows",
    "immutable upstream history"
  ]
}
````

### FILE: `odoo_loyalty/DERIVATION_CONNECTED.json`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file6:v1"
operation: CREATE
provenance: AUTHORED
source: "local source manifest, isolated loader, ORM/protocol and oracle glue distributed as part of the modified LGPL module; no Odoo corporate authorship"
license: "LGPL-3.0-only"
sha256: "ff5774ce6a226582d179f12b69a04d96ca4eabfe4d7c48db6e9df643a37ae5ce"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite.source-derivation.v1",
  "classification": "MIXED; isolated ADAPTED commercial functions and AUTHORED IPC/representation",
  "source_commit": "99edb6dd82b7b560930c00b03b694ba700785370",
  "functions": [
    {
      "destination": "filtered_issuance",
      "source": "upstream/addons/sale_loyalty/models/sale_order.py",
      "start_line": 1138,
      "end_line": 1140,
      "file_sha256": "6de1aeed8f1fc80268025f13bff2a56b57211b6a82ce9df950a8b08a623e3a06",
      "segment_sha256": "4e4d39ae35bb2054a98105a0bfc4816aa6d5069217c844ec36830ecf104878d6",
      "delta": "ORM binding replaced by narrow arguments/returned immutable receipt; Decimal arithmetic; no history deletion; precondition binding explicit"
    },
    {
      "destination": "reverse_point_change",
      "source": "upstream/addons/sale_loyalty/models/sale_order.py",
      "start_line": 200,
      "end_line": 200,
      "file_sha256": "6de1aeed8f1fc80268025f13bff2a56b57211b6a82ce9df950a8b08a623e3a06",
      "segment_sha256": "81feb0138f881626972be740b99fa270b46c7b9e554e5d692adb73a0b2e6b620",
      "delta": "ORM binding replaced by narrow arguments/returned immutable receipt; Decimal arithmetic; no history deletion; precondition binding explicit"
    },
    {
      "destination": "reward_required_points",
      "source": "upstream/addons/sale_loyalty/models/sale_order.py",
      "start_line": 1034,
      "end_line": 1035,
      "file_sha256": "6de1aeed8f1fc80268025f13bff2a56b57211b6a82ce9df950a8b08a623e3a06",
      "segment_sha256": "da3963fd19636dd46beec8dfa76cef6e43c8194f996b37f15fda1c93762e8b66",
      "delta": "ORM binding replaced by narrow arguments/returned immutable receipt; Decimal arithmetic; no history deletion; precondition binding explicit"
    }
  ],
  "rounding_delta": "Returned monetary value normalized by the same currency HALF-UP rounding representation before exact minor-unit conversion; reference currency 2 digits. Non-integral/overflow points storage rejected. Not tax calculation.",
  "journal_delta": "Odoo mutates Float points and removes cancel history. Candidate uses Decimal and append-only inverse entries. Balance may become negative after reversing spent issuance; no cancellation-generated new credit."
}
````

### FILE: `odoo_loyalty/LICENSE`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file7:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/LICENSE; revision 99edb6dd82b7b560930c00b03b694ba700785370; Git blob 12334e0fdd8d3412d5e1deed72220fcc03431ff2"
license: "LGPL-3.0-only"
sha256: "abc09dad5f84a76e1b0279237053cae16c03228ab27d8d467677054c2bd17eeb"
variables: []
secrets_allowed: false
final_newline: false
```

````text

For copyright information, please see the COPYRIGHT file.

Odoo is published under the GNU LESSER GENERAL PUBLIC LICENSE, Version 3
(LGPLv3), as included below. Since the LGPL is a set of additional
permissions on top of the GPL, the text of the GPL is included at the
bottom as well.

Some external libraries and contributions bundled with Odoo may be published
under other GPL-compatible licenses. For these, please refer to the relevant
source files and/or license files, in the source code tree.

**************************************************************************

                   GNU LESSER GENERAL PUBLIC LICENSE
                       Version 3, 29 June 2007

 Copyright (C) 2007 Free Software Foundation, Inc. <http://fsf.org/>
 Everyone is permitted to copy and distribute verbatim copies
 of this license document, but changing it is not allowed.


  This version of the GNU Lesser General Public License incorporates
the terms and conditions of version 3 of the GNU General Public
License, supplemented by the additional permissions listed below.

  0. Additional Definitions.

  As used herein, "this License" refers to version 3 of the GNU Lesser
General Public License, and the "GNU GPL" refers to version 3 of the GNU
General Public License.

  "The Library" refers to a covered work governed by this License,
other than an Application or a Combined Work as defined below.

  An "Application" is any work that makes use of an interface provided
by the Library, but which is not otherwise based on the Library.
Defining a subclass of a class defined by the Library is deemed a mode
of using an interface provided by the Library.

  A "Combined Work" is a work produced by combining or linking an
Application with the Library.  The particular version of the Library
with which the Combined Work was made is also called the "Linked
Version".

  The "Minimal Corresponding Source" for a Combined Work means the
Corresponding Source for the Combined Work, excluding any source code
for portions of the Combined Work that, considered in isolation, are
based on the Application, and not on the Linked Version.

  The "Corresponding Application Code" for a Combined Work means the
object code and/or source code for the Application, including any data
and utility programs needed for reproducing the Combined Work from the
Application, but excluding the System Libraries of the Combined Work.

  1. Exception to Section 3 of the GNU GPL.

  You may convey a covered work under sections 3 and 4 of this License
without being bound by section 3 of the GNU GPL.

  2. Conveying Modified Versions.

  If you modify a copy of the Library, and, in your modifications, a
facility refers to a function or data to be supplied by an Application
that uses the facility (other than as an argument passed when the
facility is invoked), then you may convey a copy of the modified
version:

   a) under this License, provided that you make a good faith effort to
   ensure that, in the event an Application does not supply the
   function or data, the facility still operates, and performs
   whatever part of its purpose remains meaningful, or

   b) under the GNU GPL, with none of the additional permissions of
   this License applicable to that copy.

  3. Object Code Incorporating Material from Library Header Files.

  The object code form of an Application may incorporate material from
a header file that is part of the Library.  You may convey such object
code under terms of your choice, provided that, if the incorporated
material is not limited to numerical parameters, data structure
layouts and accessors, or small macros, inline functions and templates
(ten or fewer lines in length), you do both of the following:

   a) Give prominent notice with each copy of the object code that the
   Library is used in it and that the Library and its use are
   covered by this License.

   b) Accompany the object code with a copy of the GNU GPL and this license
   document.

  4. Combined Works.

  You may convey a Combined Work under terms of your choice that,
taken together, effectively do not restrict modification of the
portions of the Library contained in the Combined Work and reverse
engineering for debugging such modifications, if you also do each of
the following:

   a) Give prominent notice with each copy of the Combined Work that
   the Library is used in it and that the Library and its use are
   covered by this License.

   b) Accompany the Combined Work with a copy of the GNU GPL and this license
   document.

   c) For a Combined Work that displays copyright notices during
   execution, include the copyright notice for the Library among
   these notices, as well as a reference directing the user to the
   copies of the GNU GPL and this license document.

   d) Do one of the following:

       0) Convey the Minimal Corresponding Source under the terms of this
       License, and the Corresponding Application Code in a form
       suitable for, and under terms that permit, the user to
       recombine or relink the Application with a modified version of
       the Linked Version to produce a modified Combined Work, in the
       manner specified by section 6 of the GNU GPL for conveying
       Corresponding Source.

       1) Use a suitable shared library mechanism for linking with the
       Library.  A suitable mechanism is one that (a) uses at run time
       a copy of the Library already present on the user's computer
       system, and (b) will operate properly with a modified version
       of the Library that is interface-compatible with the Linked
       Version.

   e) Provide Installation Information, but only if you would otherwise
   be required to provide such information under section 6 of the
   GNU GPL, and only to the extent that such information is
   necessary to install and execute a modified version of the
   Combined Work produced by recombining or relinking the
   Application with a modified version of the Linked Version. (If
   you use option 4d0, the Installation Information must accompany
   the Minimal Corresponding Source and Corresponding Application
   Code. If you use option 4d1, you must provide the Installation
   Information in the manner specified by section 6 of the GNU GPL
   for conveying Corresponding Source.)

  5. Combined Libraries.

  You may place library facilities that are a work based on the
Library side by side in a single library together with other library
facilities that are not Applications and are not covered by this
License, and convey such a combined library under terms of your
choice, if you do both of the following:

   a) Accompany the combined library with a copy of the same work based
   on the Library, uncombined with any other library facilities,
   conveyed under the terms of this License.

   b) Give prominent notice with the combined library that part of it
   is a work based on the Library, and explaining where to find the
   accompanying uncombined form of the same work.

  6. Revised Versions of the GNU Lesser General Public License.

  The Free Software Foundation may publish revised and/or new versions
of the GNU Lesser General Public License from time to time. Such new
versions will be similar in spirit to the present version, but may
differ in detail to address new problems or concerns.

  Each version is given a distinguishing version number. If the
Library as you received it specifies that a certain numbered version
of the GNU Lesser General Public License "or any later version"
applies to it, you have the option of following the terms and
conditions either of that published version or of any later version
published by the Free Software Foundation. If the Library as you
received it does not specify a version number of the GNU Lesser
General Public License, you may choose any version of the GNU Lesser
General Public License ever published by the Free Software Foundation.

  If the Library as you received it specifies that a proxy can decide
whether future versions of the GNU Lesser General Public License shall
apply, that proxy's public statement of acceptance of any version is
permanent authorization for you to choose that version for the
Library.

**************************************************************************

                    GNU GENERAL PUBLIC LICENSE
                       Version 3, 29 June 2007

 Copyright (C) 2007 Free Software Foundation, Inc. <http://fsf.org/>
 Everyone is permitted to copy and distribute verbatim copies
 of this license document, but changing it is not allowed.

                            Preamble

  The GNU General Public License is a free, copyleft license for
software and other kinds of works.

  The licenses for most software and other practical works are designed
to take away your freedom to share and change the works.  By contrast,
the GNU General Public License is intended to guarantee your freedom to
share and change all versions of a program--to make sure it remains free
software for all its users.  We, the Free Software Foundation, use the
GNU General Public License for most of our software; it applies also to
any other work released this way by its authors.  You can apply it to
your programs, too.

  When we speak of free software, we are referring to freedom, not
price.  Our General Public Licenses are designed to make sure that you
have the freedom to distribute copies of free software (and charge for
them if you wish), that you receive source code or can get it if you
want it, that you can change the software or use pieces of it in new
free programs, and that you know you can do these things.

  To protect your rights, we need to prevent others from denying you
these rights or asking you to surrender the rights.  Therefore, you have
certain responsibilities if you distribute copies of the software, or if
you modify it: responsibilities to respect the freedom of others.

  For example, if you distribute copies of such a program, whether
gratis or for a fee, you must pass on to the recipients the same
freedoms that you received.  You must make sure that they, too, receive
or can get the source code.  And you must show them these terms so they
know their rights.

  Developers that use the GNU GPL protect your rights with two steps:
(1) assert copyright on the software, and (2) offer you this License
giving you legal permission to copy, distribute and/or modify it.

  For the developers' and authors' protection, the GPL clearly explains
that there is no warranty for this free software.  For both users' and
authors' sake, the GPL requires that modified versions be marked as
changed, so that their problems will not be attributed erroneously to
authors of previous versions.

  Some devices are designed to deny users access to install or run
modified versions of the software inside them, although the manufacturer
can do so.  This is fundamentally incompatible with the aim of
protecting users' freedom to change the software.  The systematic
pattern of such abuse occurs in the area of products for individuals to
use, which is precisely where it is most unacceptable.  Therefore, we
have designed this version of the GPL to prohibit the practice for those
products.  If such problems arise substantially in other domains, we
stand ready to extend this provision to those domains in future versions
of the GPL, as needed to protect the freedom of users.

  Finally, every program is threatened constantly by software patents.
States should not allow patents to restrict development and use of
software on general-purpose computers, but in those that do, we wish to
avoid the special danger that patents applied to a free program could
make it effectively proprietary.  To prevent this, the GPL assures that
patents cannot be used to render the program non-free.

  The precise terms and conditions for copying, distribution and
modification follow.

                       TERMS AND CONDITIONS

  0. Definitions.

  "This License" refers to version 3 of the GNU General Public License.

  "Copyright" also means copyright-like laws that apply to other kinds of
works, such as semiconductor masks.

  "The Program" refers to any copyrightable work licensed under this
License.  Each licensee is addressed as "you".  "Licensees" and
"recipients" may be individuals or organizations.

  To "modify" a work means to copy from or adapt all or part of the work
in a fashion requiring copyright permission, other than the making of an
exact copy.  The resulting work is called a "modified version" of the
earlier work or a work "based on" the earlier work.

  A "covered work" means either the unmodified Program or a work based
on the Program.

  To "propagate" a work means to do anything with it that, without
permission, would make you directly or secondarily liable for
infringement under applicable copyright law, except executing it on a
computer or modifying a private copy.  Propagation includes copying,
distribution (with or without modification), making available to the
public, and in some countries other activities as well.

  To "convey" a work means any kind of propagation that enables other
parties to make or receive copies.  Mere interaction with a user through
a computer network, with no transfer of a copy, is not conveying.

  An interactive user interface displays "Appropriate Legal Notices"
to the extent that it includes a convenient and prominently visible
feature that (1) displays an appropriate copyright notice, and (2)
tells the user that there is no warranty for the work (except to the
extent that warranties are provided), that licensees may convey the
work under this License, and how to view a copy of this License.  If
the interface presents a list of user commands or options, such as a
menu, a prominent item in the list meets this criterion.

  1. Source Code.

  The "source code" for a work means the preferred form of the work
for making modifications to it.  "Object code" means any non-source
form of a work.

  A "Standard Interface" means an interface that either is an official
standard defined by a recognized standards body, or, in the case of
interfaces specified for a particular programming language, one that
is widely used among developers working in that language.

  The "System Libraries" of an executable work include anything, other
than the work as a whole, that (a) is included in the normal form of
packaging a Major Component, but which is not part of that Major
Component, and (b) serves only to enable use of the work with that
Major Component, or to implement a Standard Interface for which an
implementation is available to the public in source code form.  A
"Major Component", in this context, means a major essential component
(kernel, window system, and so on) of the specific operating system
(if any) on which the executable work runs, or a compiler used to
produce the work, or an object code interpreter used to run it.

  The "Corresponding Source" for a work in object code form means all
the source code needed to generate, install, and (for an executable
work) run the object code and to modify the work, including scripts to
control those activities.  However, it does not include the work's
System Libraries, or general-purpose tools or generally available free
programs which are used unmodified in performing those activities but
which are not part of the work.  For example, Corresponding Source
includes interface definition files associated with source files for
the work, and the source code for shared libraries and dynamically
linked subprograms that the work is specifically designed to require,
such as by intimate data communication or control flow between those
subprograms and other parts of the work.

  The Corresponding Source need not include anything that users
can regenerate automatically from other parts of the Corresponding
Source.

  The Corresponding Source for a work in source code form is that
same work.

  2. Basic Permissions.

  All rights granted under this License are granted for the term of
copyright on the Program, and are irrevocable provided the stated
conditions are met.  This License explicitly affirms your unlimited
permission to run the unmodified Program.  The output from running a
covered work is covered by this License only if the output, given its
content, constitutes a covered work.  This License acknowledges your
rights of fair use or other equivalent, as provided by copyright law.

  You may make, run and propagate covered works that you do not
convey, without conditions so long as your license otherwise remains
in force.  You may convey covered works to others for the sole purpose
of having them make modifications exclusively for you, or provide you
with facilities for running those works, provided that you comply with
the terms of this License in conveying all material for which you do
not control copyright.  Those thus making or running the covered works
for you must do so exclusively on your behalf, under your direction
and control, on terms that prohibit them from making any copies of
your copyrighted material outside their relationship with you.

  Conveying under any other circumstances is permitted solely under
the conditions stated below.  Sublicensing is not allowed; section 10
makes it unnecessary.

  3. Protecting Users' Legal Rights From Anti-Circumvention Law.

  No covered work shall be deemed part of an effective technological
measure under any applicable law fulfilling obligations under article
11 of the WIPO copyright treaty adopted on 20 December 1996, or
similar laws prohibiting or restricting circumvention of such
measures.

  When you convey a covered work, you waive any legal power to forbid
circumvention of technological measures to the extent such circumvention
is effected by exercising rights under this License with respect to
the covered work, and you disclaim any intention to limit operation or
modification of the work as a means of enforcing, against the work's
users, your or third parties' legal rights to forbid circumvention of
technological measures.

  4. Conveying Verbatim Copies.

  You may convey verbatim copies of the Program's source code as you
receive it, in any medium, provided that you conspicuously and
appropriately publish on each copy an appropriate copyright notice;
keep intact all notices stating that this License and any
non-permissive terms added in accord with section 7 apply to the code;
keep intact all notices of the absence of any warranty; and give all
recipients a copy of this License along with the Program.

  You may charge any price or no price for each copy that you convey,
and you may offer support or warranty protection for a fee.

  5. Conveying Modified Source Versions.

  You may convey a work based on the Program, or the modifications to
produce it from the Program, in the form of source code under the
terms of section 4, provided that you also meet all of these conditions:

    a) The work must carry prominent notices stating that you modified
    it, and giving a relevant date.

    b) The work must carry prominent notices stating that it is
    released under this License and any conditions added under section
    7.  This requirement modifies the requirement in section 4 to
    "keep intact all notices".

    c) You must license the entire work, as a whole, under this
    License to anyone who comes into possession of a copy.  This
    License will therefore apply, along with any applicable section 7
    additional terms, to the whole of the work, and all its parts,
    regardless of how they are packaged.  This License gives no
    permission to license the work in any other way, but it does not
    invalidate such permission if you have separately received it.

    d) If the work has interactive user interfaces, each must display
    Appropriate Legal Notices; however, if the Program has interactive
    interfaces that do not display Appropriate Legal Notices, your
    work need not make them do so.

  A compilation of a covered work with other separate and independent
works, which are not by their nature extensions of the covered work,
and which are not combined with it such as to form a larger program,
in or on a volume of a storage or distribution medium, is called an
"aggregate" if the compilation and its resulting copyright are not
used to limit the access or legal rights of the compilation's users
beyond what the individual works permit.  Inclusion of a covered work
in an aggregate does not cause this License to apply to the other
parts of the aggregate.

  6. Conveying Non-Source Forms.

  You may convey a covered work in object code form under the terms
of sections 4 and 5, provided that you also convey the
machine-readable Corresponding Source under the terms of this License,
in one of these ways:

    a) Convey the object code in, or embodied in, a physical product
    (including a physical distribution medium), accompanied by the
    Corresponding Source fixed on a durable physical medium
    customarily used for software interchange.

    b) Convey the object code in, or embodied in, a physical product
    (including a physical distribution medium), accompanied by a
    written offer, valid for at least three years and valid for as
    long as you offer spare parts or customer support for that product
    model, to give anyone who possesses the object code either (1) a
    copy of the Corresponding Source for all the software in the
    product that is covered by this License, on a durable physical
    medium customarily used for software interchange, for a price no
    more than your reasonable cost of physically performing this
    conveying of source, or (2) access to copy the
    Corresponding Source from a network server at no charge.

    c) Convey individual copies of the object code with a copy of the
    written offer to provide the Corresponding Source.  This
    alternative is allowed only occasionally and noncommercially, and
    only if you received the object code with such an offer, in accord
    with subsection 6b.

    d) Convey the object code by offering access from a designated
    place (gratis or for a charge), and offer equivalent access to the
    Corresponding Source in the same way through the same place at no
    further charge.  You need not require recipients to copy the
    Corresponding Source along with the object code.  If the place to
    copy the object code is a network server, the Corresponding Source
    may be on a different server (operated by you or a third party)
    that supports equivalent copying facilities, provided you maintain
    clear directions next to the object code saying where to find the
    Corresponding Source.  Regardless of what server hosts the
    Corresponding Source, you remain obligated to ensure that it is
    available for as long as needed to satisfy these requirements.

    e) Convey the object code using peer-to-peer transmission, provided
    you inform other peers where the object code and Corresponding
    Source of the work are being offered to the general public at no
    charge under subsection 6d.

  A separable portion of the object code, whose source code is excluded
from the Corresponding Source as a System Library, need not be
included in conveying the object code work.

  A "User Product" is either (1) a "consumer product", which means any
tangible personal property which is normally used for personal, family,
or household purposes, or (2) anything designed or sold for incorporation
into a dwelling.  In determining whether a product is a consumer product,
doubtful cases shall be resolved in favor of coverage.  For a particular
product received by a particular user, "normally used" refers to a
typical or common use of that class of product, regardless of the status
of the particular user or of the way in which the particular user
actually uses, or expects or is expected to use, the product.  A product
is a consumer product regardless of whether the product has substantial
commercial, industrial or non-consumer uses, unless such uses represent
the only significant mode of use of the product.

  "Installation Information" for a User Product means any methods,
procedures, authorization keys, or other information required to install
and execute modified versions of a covered work in that User Product from
a modified version of its Corresponding Source.  The information must
suffice to ensure that the continued functioning of the modified object
code is in no case prevented or interfered with solely because
modification has been made.

  If you convey an object code work under this section in, or with, or
specifically for use in, a User Product, and the conveying occurs as
part of a transaction in which the right of possession and use of the
User Product is transferred to the recipient in perpetuity or for a
fixed term (regardless of how the transaction is characterized), the
Corresponding Source conveyed under this section must be accompanied
by the Installation Information.  But this requirement does not apply
if neither you nor any third party retains the ability to install
modified object code on the User Product (for example, the work has
been installed in ROM).

  The requirement to provide Installation Information does not include a
requirement to continue to provide support service, warranty, or updates
for a work that has been modified or installed by the recipient, or for
the User Product in which it has been modified or installed.  Access to a
network may be denied when the modification itself materially and
adversely affects the operation of the network or violates the rules and
protocols for communication across the network.

  Corresponding Source conveyed, and Installation Information provided,
in accord with this section must be in a format that is publicly
documented (and with an implementation available to the public in
source code form), and must require no special password or key for
unpacking, reading or copying.

  7. Additional Terms.

  "Additional permissions" are terms that supplement the terms of this
License by making exceptions from one or more of its conditions.
Additional permissions that are applicable to the entire Program shall
be treated as though they were included in this License, to the extent
that they are valid under applicable law.  If additional permissions
apply only to part of the Program, that part may be used separately
under those permissions, but the entire Program remains governed by
this License without regard to the additional permissions.

  When you convey a copy of a covered work, you may at your option
remove any additional permissions from that copy, or from any part of
it.  (Additional permissions may be written to require their own
removal in certain cases when you modify the work.)  You may place
additional permissions on material, added by you to a covered work,
for which you have or can give appropriate copyright permission.

  Notwithstanding any other provision of this License, for material you
add to a covered work, you may (if authorized by the copyright holders of
that material) supplement the terms of this License with terms:

    a) Disclaiming warranty or limiting liability differently from the
    terms of sections 15 and 16 of this License; or

    b) Requiring preservation of specified reasonable legal notices or
    author attributions in that material or in the Appropriate Legal
    Notices displayed by works containing it; or

    c) Prohibiting misrepresentation of the origin of that material, or
    requiring that modified versions of such material be marked in
    reasonable ways as different from the original version; or

    d) Limiting the use for publicity purposes of names of licensors or
    authors of the material; or

    e) Declining to grant rights under trademark law for use of some
    trade names, trademarks, or service marks; or

    f) Requiring indemnification of licensors and authors of that
    material by anyone who conveys the material (or modified versions of
    it) with contractual assumptions of liability to the recipient, for
    any liability that these contractual assumptions directly impose on
    those licensors and authors.

  All other non-permissive additional terms are considered "further
restrictions" within the meaning of section 10.  If the Program as you
received it, or any part of it, contains a notice stating that it is
governed by this License along with a term that is a further
restriction, you may remove that term.  If a license document contains
a further restriction but permits relicensing or conveying under this
License, you may add to a covered work material governed by the terms
of that license document, provided that the further restriction does
not survive such relicensing or conveying.

  If you add terms to a covered work in accord with this section, you
must place, in the relevant source files, a statement of the
additional terms that apply to those files, or a notice indicating
where to find the applicable terms.

  Additional terms, permissive or non-permissive, may be stated in the
form of a separately written license, or stated as exceptions;
the above requirements apply either way.

  8. Termination.

  You may not propagate or modify a covered work except as expressly
provided under this License.  Any attempt otherwise to propagate or
modify it is void, and will automatically terminate your rights under
this License (including any patent licenses granted under the third
paragraph of section 11).

  However, if you cease all violation of this License, then your
license from a particular copyright holder is reinstated (a)
provisionally, unless and until the copyright holder explicitly and
finally terminates your license, and (b) permanently, if the copyright
holder fails to notify you of the violation by some reasonable means
prior to 60 days after the cessation.

  Moreover, your license from a particular copyright holder is
reinstated permanently if the copyright holder notifies you of the
violation by some reasonable means, this is the first time you have
received notice of violation of this License (for any work) from that
copyright holder, and you cure the violation prior to 30 days after
your receipt of the notice.

  Termination of your rights under this section does not terminate the
licenses of parties who have received copies or rights from you under
this License.  If your rights have been terminated and not permanently
reinstated, you do not qualify to receive new licenses for the same
material under section 10.

  9. Acceptance Not Required for Having Copies.

  You are not required to accept this License in order to receive or
run a copy of the Program.  Ancillary propagation of a covered work
occurring solely as a consequence of using peer-to-peer transmission
to receive a copy likewise does not require acceptance.  However,
nothing other than this License grants you permission to propagate or
modify any covered work.  These actions infringe copyright if you do
not accept this License.  Therefore, by modifying or propagating a
covered work, you indicate your acceptance of this License to do so.

  10. Automatic Licensing of Downstream Recipients.

  Each time you convey a covered work, the recipient automatically
receives a license from the original licensors, to run, modify and
propagate that work, subject to this License.  You are not responsible
for enforcing compliance by third parties with this License.

  An "entity transaction" is a transaction transferring control of an
organization, or substantially all assets of one, or subdividing an
organization, or merging organizations.  If propagation of a covered
work results from an entity transaction, each party to that
transaction who receives a copy of the work also receives whatever
licenses to the work the party's predecessor in interest had or could
give under the previous paragraph, plus a right to possession of the
Corresponding Source of the work from the predecessor in interest, if
the predecessor has it or can get it with reasonable efforts.

  You may not impose any further restrictions on the exercise of the
rights granted or affirmed under this License.  For example, you may
not impose a license fee, royalty, or other charge for exercise of
rights granted under this License, and you may not initiate litigation
(including a cross-claim or counterclaim in a lawsuit) alleging that
any patent claim is infringed by making, using, selling, offering for
sale, or importing the Program or any portion of it.

  11. Patents.

  A "contributor" is a copyright holder who authorizes use under this
License of the Program or a work on which the Program is based.  The
work thus licensed is called the contributor's "contributor version".

  A contributor's "essential patent claims" are all patent claims
owned or controlled by the contributor, whether already acquired or
hereafter acquired, that would be infringed by some manner, permitted
by this License, of making, using, or selling its contributor version,
but do not include claims that would be infringed only as a
consequence of further modification of the contributor version.  For
purposes of this definition, "control" includes the right to grant
patent sublicenses in a manner consistent with the requirements of
this License.

  Each contributor grants you a non-exclusive, worldwide, royalty-free
patent license under the contributor's essential patent claims, to
make, use, sell, offer for sale, import and otherwise run, modify and
propagate the contents of its contributor version.

  In the following three paragraphs, a "patent license" is any express
agreement or commitment, however denominated, not to enforce a patent
(such as an express permission to practice a patent or covenant not to
sue for patent infringement).  To "grant" such a patent license to a
party means to make such an agreement or commitment not to enforce a
patent against the party.

  If you convey a covered work, knowingly relying on a patent license,
and the Corresponding Source of the work is not available for anyone
to copy, free of charge and under the terms of this License, through a
publicly available network server or other readily accessible means,
then you must either (1) cause the Corresponding Source to be so
available, or (2) arrange to deprive yourself of the benefit of the
patent license for this particular work, or (3) arrange, in a manner
consistent with the requirements of this License, to extend the patent
license to downstream recipients.  "Knowingly relying" means you have
actual knowledge that, but for the patent license, your conveying the
covered work in a country, or your recipient's use of the covered work
in a country, would infringe one or more identifiable patents in that
country that you have reason to believe are valid.

  If, pursuant to or in connection with a single transaction or
arrangement, you convey, or propagate by procuring conveyance of, a
covered work, and grant a patent license to some of the parties
receiving the covered work authorizing them to use, propagate, modify
or convey a specific copy of the covered work, then the patent license
you grant is automatically extended to all recipients of the covered
work and works based on it.

  A patent license is "discriminatory" if it does not include within
the scope of its coverage, prohibits the exercise of, or is
conditioned on the non-exercise of one or more of the rights that are
specifically granted under this License.  You may not convey a covered
work if you are a party to an arrangement with a third party that is
in the business of distributing software, under which you make payment
to the third party based on the extent of your activity of conveying
the work, and under which the third party grants, to any of the
parties who would receive the covered work from you, a discriminatory
patent license (a) in connection with copies of the covered work
conveyed by you (or copies made from those copies), or (b) primarily
for and in connection with specific products or compilations that
contain the covered work, unless you entered into that arrangement,
or that patent license was granted, prior to 28 March 2007.

  Nothing in this License shall be construed as excluding or limiting
any implied license or other defenses to infringement that may
otherwise be available to you under applicable patent law.

  12. No Surrender of Others' Freedom.

  If conditions are imposed on you (whether by court order, agreement or
otherwise) that contradict the conditions of this License, they do not
excuse you from the conditions of this License.  If you cannot convey a
covered work so as to satisfy simultaneously your obligations under this
License and any other pertinent obligations, then as a consequence you may
not convey it at all.  For example, if you agree to terms that obligate you
to collect a royalty for further conveying from those to whom you convey
the Program, the only way you could satisfy both those terms and this
License would be to refrain entirely from conveying the Program.

  13. Use with the GNU Affero General Public License.

  Notwithstanding any other provision of this License, you have
permission to link or combine any covered work with a work licensed
under version 3 of the GNU Affero General Public License into a single
combined work, and to convey the resulting work.  The terms of this
License will continue to apply to the part which is the covered work,
but the special requirements of the GNU Affero General Public License,
section 13, concerning interaction through a network will apply to the
combination as such.

  14. Revised Versions of this License.

  The Free Software Foundation may publish revised and/or new versions of
the GNU General Public License from time to time.  Such new versions will
be similar in spirit to the present version, but may differ in detail to
address new problems or concerns.

  Each version is given a distinguishing version number.  If the
Program specifies that a certain numbered version of the GNU General
Public License "or any later version" applies to it, you have the
option of following the terms and conditions either of that numbered
version or of any later version published by the Free Software
Foundation.  If the Program does not specify a version number of the
GNU General Public License, you may choose any version ever published
by the Free Software Foundation.

  If the Program specifies that a proxy can decide which future
versions of the GNU General Public License can be used, that proxy's
public statement of acceptance of a version permanently authorizes you
to choose that version for the Program.

  Later license versions may give you additional or different
permissions.  However, no additional obligations are imposed on any
author or copyright holder as a result of your choosing to follow a
later version.

  15. Disclaimer of Warranty.

  THERE IS NO WARRANTY FOR THE PROGRAM, TO THE EXTENT PERMITTED BY
APPLICABLE LAW.  EXCEPT WHEN OTHERWISE STATED IN WRITING THE COPYRIGHT
HOLDERS AND/OR OTHER PARTIES PROVIDE THE PROGRAM "AS IS" WITHOUT WARRANTY
OF ANY KIND, EITHER EXPRESSED OR IMPLIED, INCLUDING, BUT NOT LIMITED TO,
THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR
PURPOSE.  THE ENTIRE RISK AS TO THE QUALITY AND PERFORMANCE OF THE PROGRAM
IS WITH YOU.  SHOULD THE PROGRAM PROVE DEFECTIVE, YOU ASSUME THE COST OF
ALL NECESSARY SERVICING, REPAIR OR CORRECTION.

  16. Limitation of Liability.

  IN NO EVENT UNLESS REQUIRED BY APPLICABLE LAW OR AGREED TO IN WRITING
WILL ANY COPYRIGHT HOLDER, OR ANY OTHER PARTY WHO MODIFIES AND/OR CONVEYS
THE PROGRAM AS PERMITTED ABOVE, BE LIABLE TO YOU FOR DAMAGES, INCLUDING ANY
GENERAL, SPECIAL, INCIDENTAL OR CONSEQUENTIAL DAMAGES ARISING OUT OF THE
USE OR INABILITY TO USE THE PROGRAM (INCLUDING BUT NOT LIMITED TO LOSS OF
DATA OR DATA BEING RENDERED INACCURATE OR LOSSES SUSTAINED BY YOU OR THIRD
PARTIES OR A FAILURE OF THE PROGRAM TO OPERATE WITH ANY OTHER PROGRAMS),
EVEN IF SUCH HOLDER OR OTHER PARTY HAS BEEN ADVISED OF THE POSSIBILITY OF
SUCH DAMAGES.

  17. Interpretation of Sections 15 and 16.

  If the disclaimer of warranty and limitation of liability provided
above cannot be given local legal effect according to their terms,
reviewing courts shall apply local law that most closely approximates
an absolute waiver of all civil liability in connection with the
Program, unless a warranty or assumption of liability accompanies a
copy of the Program in return for a fee.

                     END OF TERMS AND CONDITIONS

            How to Apply These Terms to Your New Programs

  If you develop a new program, and you want it to be of the greatest
possible use to the public, the best way to achieve this is to make it
free software which everyone can redistribute and change under these terms.

  To do so, attach the following notices to the program.  It is safest
to attach them to the start of each source file to most effectively
state the exclusion of warranty; and each file should have at least
the "copyright" line and a pointer to where the full notice is found.

    <one line to give the program's name and a brief idea of what it does.>
    Copyright (C) <year>  <name of author>

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.

Also add information on how to contact you by electronic and paper mail.

  If the program does terminal interaction, make it output a short
notice like this when it starts in an interactive mode:

    <program>  Copyright (C) <year>  <name of author>
    This program comes with ABSOLUTELY NO WARRANTY; for details type `show w'.
    This is free software, and you are welcome to redistribute it
    under certain conditions; type `show c' for details.

The hypothetical commands `show w' and `show c' should show the appropriate
parts of the General Public License.  Of course, your program's commands
might be different; for a GUI interface, you would use an "about box".

  You should also get your employer (if you work as a programmer) or school,
if any, to sign a "copyright disclaimer" for the program, if necessary.
For more information on this, and how to apply and follow the GNU GPL, see
<http://www.gnu.org/licenses/>.

  The GNU General Public License does not permit incorporating your program
into proprietary programs.  If your program is a subroutine library, you
may consider it more useful to permit linking proprietary applications with
the library.  If this is what you want to do, use the GNU Lesser General
Public License instead of this License.  But first, please read
<http://www.gnu.org/philosophy/why-not-lgpl.html>.


**************************************************************************
````

### FILE: `odoo_loyalty/engine-lock.json`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file8:v1"
operation: CREATE
provenance: AUTHORED
source: "local source manifest, isolated loader, ORM/protocol and oracle glue distributed as part of the modified LGPL module; no Odoo corporate authorship"
license: "LGPL-3.0-only"
sha256: "d39918d5c4479af0d961fa9b3b61f3fb8b00e997d03be958083debad627657da"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite.odoo-derived-source-lock.v1",
  "files": [
    {
      "path": "COPYRIGHT",
      "bytes": 433,
      "sha256": "86e49232d2162708d05405ed5ff6dc5594b73ef4b0f25d2749cfae13493b620c"
    },
    {
      "path": "DERIVATION.json",
      "bytes": 2828,
      "sha256": "26d841355c73e08cd9e99fd855c9fba0a69b6681b53b23275621c4741fa58fe2"
    },
    {
      "path": "DERIVATION_CONNECTED.json",
      "bytes": 2245,
      "sha256": "ff5774ce6a226582d179f12b69a04d96ca4eabfe4d7c48db6e9df643a37ae5ce"
    },
    {
      "path": "LICENSE",
      "bytes": 43529,
      "sha256": "abc09dad5f84a76e1b0279237053cae16c03228ab27d8d467677054c2bd17eeb"
    },
    {
      "path": "engine.py",
      "bytes": 13177,
      "sha256": "7c0bfe5065d94eac389484b40fa3842b43504864655d475c4c053af0d9209e73"
    },
    {
      "path": "orm_contract.py",
      "bytes": 8004,
      "sha256": "7528cca6cdf442c838ec52393512e2cb8cf8e7382885035c1b2957da77a581b2"
    },
    {
      "path": "protocol.py",
      "bytes": 12914,
      "sha256": "5c0ae67ca00d31f41085b70258716043faac5e0a20953de3646b7e3559083941"
    },
    {
      "path": "run.py",
      "bytes": 1980,
      "sha256": "2101f3e889607a8da2fd22722d64c20fd481fa2c38882e1495aeffd72bb3387e"
    },
    {
      "path": "upstream/COPYRIGHT",
      "bytes": 433,
      "sha256": "86e49232d2162708d05405ed5ff6dc5594b73ef4b0f25d2749cfae13493b620c"
    },
    {
      "path": "upstream/LICENSE",
      "bytes": 43529,
      "sha256": "abc09dad5f84a76e1b0279237053cae16c03228ab27d8d467677054c2bd17eeb"
    },
    {
      "path": "upstream/addons/loyalty/__manifest__.py",
      "bytes": 1754,
      "sha256": "03537054d1786db61c364262242eff900c6b104c2e1e5c8df3a60cd2b1e403f5"
    },
    {
      "path": "upstream/addons/loyalty/models/loyalty_card.py",
      "bytes": 8800,
      "sha256": "9eead538d45b3ad32a0f2eddc6b146d2a0ab5cc4d40de19bdea58aa70908ba83"
    },
    {
      "path": "upstream/addons/loyalty/models/loyalty_history.py",
      "bytes": 885,
      "sha256": "c480709def3ad815fe6b306dac534785f24cf10b7432276a0276dcd5646bb8a5"
    },
    {
      "path": "upstream/addons/loyalty/models/loyalty_program.py",
      "bytes": 27848,
      "sha256": "3ed5484b4404120c73afe810653c2b25915d50d4211158005f684d9290e08096"
    },
    {
      "path": "upstream/addons/loyalty/models/loyalty_reward.py",
      "bytes": 15597,
      "sha256": "ec0aa37a9a5f12c01e4a860738afcc1b6ede113e4a3f4aa9d08966194936955c"
    },
    {
      "path": "upstream/addons/loyalty/models/loyalty_rule.py",
      "bytes": 6779,
      "sha256": "955eafbc11279154ba30b4986fa2198c683d4c4519253f5ee0fd84daa39f2a66"
    },
    {
      "path": "upstream/addons/loyalty/tests/test_loyalty.py",
      "bytes": 14565,
      "sha256": "39c7660a20b402f80fafab7e6534c9bbcd2e52abd847088c0263004fe33092c3"
    },
    {
      "path": "upstream/addons/sale_loyalty/__manifest__.py",
      "bytes": 1039,
      "sha256": "284a427b610067bac9c82ae6e3cac37cb1ad51c3dee5617fc5b93d62b13ee083"
    },
    {
      "path": "upstream/addons/sale_loyalty/models/loyalty_card.py",
      "bytes": 2147,
      "sha256": "67a9501daf1b48bffbf10e34d11b7503df920a0048545affd38e86f08f2be712"
    },
    {
      "path": "upstream/addons/sale_loyalty/models/loyalty_program.py",
      "bytes": 1030,
      "sha256": "2f1bbcfb4419d8bd6149eac02e963670c342470dd2007eaa02e31b095704e658"
    },
    {
      "path": "upstream/addons/sale_loyalty/models/sale_order.py",
      "bytes": 78182,
      "sha256": "6de1aeed8f1fc80268025f13bff2a56b57211b6a82ce9df950a8b08a623e3a06"
    },
    {
      "path": "upstream/addons/sale_loyalty/models/sale_order_coupon_points.py",
      "bytes": 676,
      "sha256": "81a49606130beed878aaa93c503970c8e4e9b4645b63f2b097999ccb99c4eeb2"
    },
    {
      "path": "upstream/addons/sale_loyalty/models/sale_order_line.py",
      "bytes": 7181,
      "sha256": "4aba6673de316b9cb78def3c1d4e31d456ce170524c84d5f8dca93b2a8f18d93"
    },
    {
      "path": "upstream/addons/sale_loyalty/tests/common.py",
      "bytes": 12605,
      "sha256": "ab2f2ef3334a82a642178427f24e8316e3ee2d4f91c63dede8712304d3c53db5"
    },
    {
      "path": "upstream/addons/sale_loyalty/tests/test_buy_gift_card.py",
      "bytes": 3202,
      "sha256": "f724115af4e50776f152bb933b01999bef7e800aec4eea0c43ad38f9916bb3e8"
    },
    {
      "path": "upstream/addons/sale_loyalty/tests/test_loyalty.py",
      "bytes": 55159,
      "sha256": "b2e97e780f49b8303b1930cdc7b6687338bcdf0cd48039e747144eddd14db47e"
    },
    {
      "path": "upstream/addons/sale_loyalty/tests/test_loyalty_history.py",
      "bytes": 8005,
      "sha256": "291cee97246ffe07b39d77b4c0bdaa3c3addb7cec321aa2235827ab92515775d"
    },
    {
      "path": "upstream/addons/sale_loyalty/tests/test_pay_with_gift_card.py",
      "bytes": 10659,
      "sha256": "9ecd12f97ddb7dcdcb8147d38752a1155c21583977e2ce048e40e653b4468f47"
    },
    {
      "path": "upstream/addons/sale_loyalty/tests/test_program_multi_company.py",
      "bytes": 4963,
      "sha256": "d338ab072168bacfa342f7a5a0011063d47f41f19ef1221a4adf0a70cb90a32b"
    },
    {
      "path": "upstream/addons/sale_loyalty/tests/test_program_numbers.py",
      "bytes": 97516,
      "sha256": "845f5b205b1ac941280905bfa308ed15d7142d3069997b6545f5cc0e46b0c918"
    },
    {
      "path": "upstream/addons/sale_loyalty/tests/test_program_rules.py",
      "bytes": 21582,
      "sha256": "f28fa822aa5d7116816f20f2010c2c3f94f926477de024b768ed52d44c3cbce1"
    },
    {
      "path": "upstream/odoo/tools/float_utils.py",
      "bytes": 19196,
      "sha256": "ab7d51a44033792414dc63f598517ee76b922f76b2de9fa3eeb977ab78d7a962"
    }
  ]
}
````

### FILE: `odoo_loyalty/engine.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file9:v1"
operation: CREATE
provenance: ADAPTED
source: "https://github.com/odoo/odoo/tree/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty; exact file/segment mappings in odoo_loyalty/DERIVATION.json and DERIVATION_CONNECTED.json; Decimal/IPC and immutable history adaptations declared; test expectations adapted, no full Odoo integration-suite claim"
license: "LGPL-3.0-only"
sha256: "7c0bfe5065d94eac389484b40fa3842b43504864655d475c4c053af0d9209e73"
variables: []
secrets_allowed: false
```

````python
# SPDX-License-Identifier: LGPL-3.0-only
# Derived from Odoo Community 19.0, commit 99edb6dd82b7b560930c00b03b694ba700785370.
# Original copyright and full LGPLv3/GPLv3 texts: COPYRIGHT and LICENSE.
# Modified for an isolated exact-decimal calculation module; see DERIVATION.json.
from collections import defaultdict
from decimal import Decimal, ROUND_DOWN, ROUND_HALF_UP, ROUND_HALF_EVEN

def _(message, **values):
    return message % values if values else message

class _Date:
    @staticmethod
    def today():
        # The contract rejects currency conversion; this argument is unused by same-currency adapter.
        return None

class fields:
    Date = _Date

def float_round(value, precision_digits=None, precision_rounding=None, rounding_method='HALF-UP'):
    # AUTHORED exact representation adapter; business rounding modes come from the mapped source.
    if (precision_digits is None) == (precision_rounding is None):
        raise ValueError('one rounding precision is required')
    step = Decimal(str(precision_rounding)) if precision_rounding is not None else Decimal(10) ** -precision_digits
    if step <= 0 or not step.is_finite():
        raise ValueError('invalid rounding precision')
    modes = {'DOWN': ROUND_DOWN, 'HALF-UP': ROUND_HALF_UP, 'HALF-EVEN': ROUND_HALF_EVEN}
    if rounding_method not in modes:
        raise ValueError('unsupported rounding mode')
    return (value / step).quantize(Decimal(1), rounding=modes[rounding_method]) * step

def _get_point_changes(self):
    """
    Returns the changes in points per coupon as a dict.

    Used when validating/cancelling an order
    """
    points_per_coupon = defaultdict(lambda: 0)
    for coupon_point in self.coupon_point_ids:
        points_per_coupon[coupon_point.coupon_id] += coupon_point.points
    for line in self.order_line:
        if not line.reward_id or not line.coupon_id:
            continue
        points_per_coupon[line.coupon_id] -= line.points_cost
    return points_per_coupon


def _get_real_points_for_coupon(self, coupon, post_confirm=False):
    """
    Returns the actual points usable for this coupon for this order. Set pos_confirm to True to include points for future orders.

    This is calculated by taking the points on the coupon, the points the order will give to the coupon (if applicable) and removing the points taken by already applied rewards.
    """
    self.ensure_one()
    points = coupon.points
    if self.state not in ('sale', 'done'):
        if coupon.program_id.applies_on != 'future':
            # Points that will be given by the order upon confirming the order
            points += self.coupon_point_ids.filtered(lambda p: p.coupon_id == coupon).points
        # Points already used by rewards
        points -= sum(self.order_line.filtered(lambda l: l.coupon_id == coupon).mapped('points_cost'))
    if any(rule.reward_point_mode == 'money' for rule in coupon.program_id.rule_ids):
        points = coupon.currency_id.round(points)
    return points


def _program_check_compute_points(self, programs):
    """
    Checks the program validity from the order lines aswell as computing the number of points to add.

    Returns a dict containing the error message or the points that will be given with the keys 'points'.
    """
    self.ensure_one()

    # Prepare quantities
    order_lines = self._get_not_rewarded_order_lines().filtered(
        lambda line: not line.combo_item_id
    )
    products = order_lines.product_id
    products_qties = dict.fromkeys(products, 0)
    for line in order_lines:
        product_qty = line.product_uom_id._compute_quantity(
            line.product_uom_qty, line.product_id.uom_id
        )
        products_qties[line.product_id] += product_qty
    # Contains the products that can be applied per rule
    products_per_rule = programs._get_valid_products(products)

    # Prepare amounts
    so_products_per_rule = programs._get_valid_products(self.order_line.product_id)
    lines_per_rule = defaultdict(lambda: self.env['sale.order.line'])
    # Skip lines that have no effect on the minimum amount to reach.
    for line in self.order_line - self._get_no_effect_on_threshold_lines():
        is_discount = line.reward_id.reward_type == 'discount'
        reward_program = line.reward_id.program_id
        # Skip lines for automatic discounts, as well as combo item lines.
        if (is_discount and reward_program.trigger == 'auto') or line.combo_item_id:
            continue
        for program in programs:
            # Skip lines for the current program's discounts.
            if is_discount and reward_program == program:
                continue
            for rule in program.rule_ids:
                # Skip lines to which the rule doesn't apply.
                if line.product_id in so_products_per_rule.get(rule, []):
                    lines_per_rule[rule] |= line._get_lines_with_price()

    result = {}
    for program in programs:
        # Used for error messages
        # By default False, but True if no rules and applies_on current -> misconfigured coupons program
        code_matched = not bool(program.rule_ids) and program.applies_on == 'current' # Stays false if all triggers have code and none have been activated
        minimum_amount_matched = code_matched
        product_qty_matched = code_matched
        points = 0
        # Some rules may split their points per unit / money spent
        #  (i.e. gift cards 2x50$ must result in two 50$ codes)
        rule_points = []
        program_result = result.setdefault(program, dict())
        for rule in program.rule_ids:
            # prevent bottomless ewallet spending
            if program.program_type == 'ewallet' and not program.trigger_product_ids:
                break
            if rule.mode == 'with_code' and rule not in self.code_enabled_rule_ids:
                continue
            code_matched = True
            rule_amount = rule._compute_amount(self.currency_id)
            untaxed_amount = sum(lines_per_rule[rule].mapped('price_subtotal'))
            tax_amount = sum(lines_per_rule[rule].mapped('price_tax'))
            if rule_amount > (rule.minimum_amount_tax_mode == 'incl' and (untaxed_amount + tax_amount) or untaxed_amount):
                continue
            minimum_amount_matched = True
            if not products_per_rule.get(rule):
                continue
            rule_products = products_per_rule[rule]
            ordered_rule_products_qty = sum(products_qties[product] for product in rule_products)
            if ordered_rule_products_qty < rule.minimum_qty or not rule_products:
                continue
            product_qty_matched = True
            if not rule.reward_point_amount:
                continue
            # Count all points separately if the order is for the future and the split option is enabled
            if program.applies_on == 'future' and rule.reward_point_split and rule.reward_point_mode != 'order':
                if rule.reward_point_mode == 'unit':
                    rule_points.extend(rule.reward_point_amount for _ in range(int(ordered_rule_products_qty)))
                elif rule.reward_point_mode == 'money':
                    for line in self.order_line:
                        if (
                            line.is_reward_line
                            or line.combo_item_id
                            or line.product_id not in rule_products
                            or line.product_uom_qty <= 0
                        ):
                            continue
                        line_price_total = self._get_order_line_price(line, 'price_total')
                        points_per_unit = float_round(
                            (rule.reward_point_amount * line_price_total / line.product_uom_qty),
                            precision_digits=2, rounding_method='DOWN')
                        if not points_per_unit:
                            continue
                        rule_points.extend([points_per_unit] * int(line.product_uom_qty))
            else:
                # All checks have been passed we can now compute the points to give
                if rule.reward_point_mode == 'order':
                    points += rule.reward_point_amount
                elif rule.reward_point_mode == 'money':
                    # Compute amount paid for rule
                    # NOTE: this accounts for discounts -> 1 point per $ * (100$ - 30%) will
                    # result in 70 points
                    amount_paid = Decimal('0.0')
                    rule_products = so_products_per_rule.get(rule, [])
                    for line in self.order_line - self._get_no_effect_on_threshold_lines():
                        if line.combo_item_id or line.reward_id.program_id.program_type in [
                            'ewallet', 'gift_card', program.program_type
                        ]:
                            continue
                        line_price_total = self._get_order_line_price(line, 'price_total')
                        amount_paid += (
                            line_price_total if line.product_id in rule_products
                            else Decimal('0.0')
                        )

                    points += float_round(rule.reward_point_amount * amount_paid, precision_digits=2, rounding_method='DOWN')
                elif rule.reward_point_mode == 'unit':
                    points += rule.reward_point_amount * ordered_rule_products_qty
        # NOTE: for programs that are nominative we always allow the program to be 'applied' on the order
        #  with 0 points so that `_get_claimable_rewards` returns the rewards associated with those programs
        if not program.is_nominative:
            if not code_matched:
                program_result['error'] = _("This program requires a code to be applied.")
            elif not minimum_amount_matched:
                program_result['error'] = _(
                    "To take advantage of this offer, your order must include at least %(amount)s %(currency)s of the eligible products.",
                    amount=min(program.rule_ids.mapped('minimum_amount')),
                    currency=program.currency_id.name,
                )
            elif not product_qty_matched:
                program_result['error'] = _("You don't have the required product quantities on your sales order.")
        elif self.partner_id.is_public and not self._allow_nominative_programs():
            program_result['error'] = _("This program is not available for public users.")
        if 'error' not in program_result:
            points_result = [points] + rule_points
            program_result['points'] = points_result
    return result


def reward_from_bound_amounts(self, reward, coupon, discountable):
    reward_currency = reward.currency_id
    reward_program = reward.program_id
    max_discount = reward_currency._convert(reward.discount_max_amount, self.currency_id, self.company_id, fields.Date.today()) or Decimal('Infinity')
    # discount should never surpass the order's current total amount
    max_discount = min(self.amount_total, max_discount)
    if reward.discount_mode == 'per_point':
        points = self._get_real_points_for_coupon(coupon)
        if not reward_program.is_payment_program:
            # Rewards cannot be partially offered to customers
            points = points // reward.required_points * reward.required_points
        max_discount = min(max_discount,
            reward_currency._convert(reward.discount * points,
                self.currency_id, self.company_id, fields.Date.today()))
    elif reward.discount_mode == 'per_order':
        max_discount = min(max_discount,
            reward_currency._convert(reward.discount, self.currency_id, self.company_id, fields.Date.today()))
    elif reward.discount_mode == 'percent':
        max_discount = min(max_discount, discountable * (reward.discount / 100))

    # Discount per taxes
    point_cost = reward.required_points if not reward.clear_wallet else self._get_real_points_for_coupon(coupon)
    if reward.discount_mode == 'per_point' and not reward.clear_wallet:
        # Calculate the actual point cost if the cost is per point
        converted_discount = self.currency_id._convert(min(max_discount, discountable), reward_currency, self.company_id, fields.Date.today())
        point_cost = coupon.currency_id.round(converted_discount / reward.discount)

    return {"value": min(max_discount, discountable), "points_cost": point_cost}



# ADAPTED sale_order.py:1138-1140. Container arguments replace ORM attributes.
def filtered_issuance(status, program):
    all_point_changes = [p for p in status['points'] if p]
    if not all_point_changes and program.is_nominative:
        all_point_changes = [0]
    return all_point_changes

# ADAPTED exact cancellation arithmetic from _action_cancel: coupon.points -= changes.
# Argument/return plumbing replaces an ORM property update; immutable journal
# delta is computed against zero rather than mutating or deleting historical rows.
def reverse_point_change(points, changes):
    points -= changes
    return points
````

### FILE: `odoo_loyalty/oracle_loader.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file10:v1"
operation: CREATE
provenance: AUTHORED
source: "local source manifest, isolated loader, ORM/protocol and oracle glue distributed as part of the modified LGPL module; no Odoo corporate authorship"
license: "LGPL-3.0-only"
sha256: "6c23b3b3e19cdd39a7d91de166722addd670ad2e903718ad95d0d3df1557bcc6"
variables: []
secrets_allowed: false
```

````python
# AUTHORED test-only loading glue. Exact source functions retain Odoo LGPL attribution.
import ast
import hashlib
import importlib.util
from pathlib import Path
from types import SimpleNamespace
from collections import defaultdict

ROOT = Path(__file__).parent
_float_path = ROOT/'upstream/odoo/tools/float_utils.py'
_spec = importlib.util.spec_from_file_location('odoo_original_float_utils', _float_path)
_float_module = importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(_float_module)
original_float_round = _float_module.float_round


def original_functions():
    path = ROOT/'upstream/addons/sale_loyalty/models/sale_order.py'
    tree = ast.parse(path.read_bytes())
    names = {'_get_point_changes', '_get_real_points_for_coupon', '_program_check_compute_points'}
    nodes = [node for node in ast.walk(tree) if isinstance(node, ast.FunctionDef) and node.name in names]
    if len(nodes) != len(names):
        raise ValueError('source oracle missing functions')
    scope = {'defaultdict':defaultdict,'float_round':original_float_round,
             '_':lambda text, **kwargs: text % kwargs if kwargs else text}
    exec(compile(ast.Module(body=nodes,type_ignores=[]), str(path), 'exec'), scope)
    return SimpleNamespace(**{name:scope[name] for name in names})
````

### FILE: `odoo_loyalty/orm_contract.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file11:v1"
operation: CREATE
provenance: AUTHORED
source: "local source manifest, isolated loader, ORM/protocol and oracle glue distributed as part of the modified LGPL module; no Odoo corporate authorship"
license: "LGPL-3.0-only"
sha256: "7528cca6cdf442c838ec52393512e2cb8cf8e7382885035c1b2957da77a581b2"
variables: []
secrets_allowed: false
```

````python
# AUTHORED interface glue for the isolated source-derived calculations.
# This is a bounded recordset protocol, not Odoo ORM or an implementation of its database.
from decimal import Decimal, ROUND_HALF_UP
from types import MethodType


class NullRecord:
    def __bool__(self):
        return False
    def __getattr__(self, _name):
        return self
    def __eq__(self, other):
        return isinstance(other, NullRecord)
    def __hash__(self):
        return 0


NULL = NullRecord()


class Record:
    def __init__(self, kind, identity, **values):
        self._kind = kind
        self.id = identity
        self.__dict__.update(values)
    def __hash__(self):
        return hash((self._kind, self.id))
    def __eq__(self, other):
        return isinstance(other, Record) and (self._kind, self.id) == (other._kind, other.id)
    def ensure_one(self):
        return self


class Records:
    def __init__(self, values=()):
        self.values = tuple(dict.fromkeys(values))
    def __iter__(self):
        return iter(self.values)
    def __bool__(self):
        return bool(self.values)
    def __len__(self):
        return len(self.values)
    def __contains__(self, value):
        return value in self.values
    def __or__(self, other):
        return Records(self.values + other.values)
    def __sub__(self, other):
        return Records(x for x in self.values if x not in other)
    def filtered(self, condition):
        if isinstance(condition, str):
            return Records(x for x in self.values if getattr(x, condition))
        return Records(x for x in self.values if condition(x))
    def mapped(self, field):
        values = []
        for record in self.values:
            value = record
            for name in field.split('.'):
                value = getattr(value, name)
            if isinstance(value, Records):
                values.extend(value.values)
            else:
                values.append(value)
        if all(isinstance(x, (Record, NullRecord)) for x in values):
            return Records(x for x in values if x)
        return values
    def __getattr__(self, field):
        if field == 'ids':
            return [x.id for x in self.values]
        result = self.mapped(field)
        if isinstance(result, Records):
            return result
        if not result:
            # Odoo scalar access on an empty recordset is False; numeric empty fields add as zero.
            return False
        if len(result) != 1:
            raise ValueError('recordset scalar requires one record')
        return result[0]
    def _get_valid_products(self, products):
        # Only explicit product-id filters are supported by the profile. Categories/tags/domains are rejected at IPC.
        return {rule: Records(p for p in products if not rule.product_ids or p.id in rule.product_ids)
                for program in self.values for rule in program.rule_ids}


class Currency(Record):
    def __init__(self, code, digits, number=Decimal):
        super().__init__('currency', code, name=code)
        self.digits = digits
        self.number = number
    def _convert(self, amount, other, _company, _date):
        if self.id != other.id:
            raise ValueError('currency conversion is outside this module')
        return amount
    def round(self, amount):
        if self.number is Decimal:
            return amount.quantize(Decimal(10) ** -self.digits, rounding=ROUND_HALF_UP)
        from .oracle_loader import original_float_round
        return original_float_round(amount, precision_digits=self.digits)


class Unit(Record):
    def _compute_quantity(self, quantity, other):
        if self.id != other.id:
            raise ValueError('UoM conversion is outside this module')
        return quantity


class Environment:
    def __getitem__(self, key):
        if key != 'sale.order.line':
            raise ValueError('unsupported ORM surface')
        return Records()


class Order(Record):
    def _get_not_rewarded_order_lines(self):
        return self.order_line.filtered(lambda line: line.product_id and not line.reward_id)
    def _get_no_effect_on_threshold_lines(self):
        return self.order_line.filtered(lambda line: line.threshold_excluded)
    def _get_order_line_price(self, line, price_type):
        return sum(line._get_lines_with_price().mapped(price_type))
    def _allow_nominative_programs(self):
        return False


def build(program_data, order_data, number=Decimal, functions=None):
    from . import engine
    functions = functions or engine
    currency = Currency(program_data['currency'], program_data['currency_digits'], number)
    company = Record('company', order_data['organization_id'])
    program = Record('program', program_data['id'], program_type=program_data['kind'],
                     applies_on=program_data['applies_on'], is_nominative=program_data['nominative'],
                     trigger=program_data['trigger'], currency_id=currency,
                     is_payment_program=program_data['kind'] in ('gift_card', 'ewallet'),
                     trigger_product_ids=program_data['trigger_product_ids'])
    rules = []
    for data in program_data['rules']:
        rule = Record('rule', data['id'], program_id=program,
                      reward_point_mode=data['mode'], reward_point_amount=number(data['points']),
                      reward_point_split=data['split'], minimum_qty=number(data['minimum_qty']),
                      minimum_amount=number(data['minimum_amount']),
                      minimum_amount_tax_mode=data['tax_mode'], product_ids=data['product_ids'],
                      mode='with_code' if data['code_required'] else 'auto')
        rule._compute_amount = lambda target, rule=rule: currency._convert(rule.minimum_amount, target, company, None)
        rules.append(rule)
    program.rule_ids = Records(rules)
    unit = Unit('unit', 'normalized-unit')
    products = {data['product_id']: Record('product', data['product_id'], uom_id=unit) for data in order_data['lines']}
    lines = []
    for data in order_data['lines']:
        reward = NULL
        if data['reward_program_type']:
            reward_program = program if data['reward_program_id'] == program.id else Record(
                'program', data['reward_program_id'], program_type=data['reward_program_type'], trigger=data['reward_trigger'])
            reward = Record('reward', data['id'], reward_type='discount', program_id=reward_program)
        line = Record('line', data['id'], product_id=products[data['product_id']],
                      product_uom_id=unit, product_uom_qty=number(data['quantity']),
                      price_subtotal=number(data['subtotal']), price_tax=number(data['tax']),
                      price_total=number(data['total']), combo_item_id=False,
                      is_reward_line=bool(reward), reward_id=reward,
                      threshold_excluded=data['threshold_excluded'],
                      coupon_id=NULL, points_cost=number('0'))
        line._get_lines_with_price = lambda line=line: Records([line])
        lines.append(line)
    order = Order('order', order_data['order_id'], state=order_data['state'],
                  partner_id=Record('partner', order_data['subject_id'], is_public=order_data['public_subject']),
                  company_id=company, currency_id=currency, amount_total=number(order_data['total']),
                  order_line=Records(lines), coupon_point_ids=Records(),
                  code_enabled_rule_ids=Records(r for r in rules if r.id in order_data['enabled_rule_ids']),
                  env=Environment())
    for name in ('_get_point_changes', '_get_real_points_for_coupon', '_program_check_compute_points'):
        setattr(order, name, MethodType(getattr(functions, name), order))
    return program, order


def coupon_for(program, identity, points, number=Decimal):
    return Record('coupon', identity, program_id=program, currency_id=program.currency_id, points=number(points))
````

### FILE: `odoo_loyalty/protocol.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file12:v1"
operation: CREATE
provenance: AUTHORED
source: "local source manifest, isolated loader, ORM/protocol and oracle glue distributed as part of the modified LGPL module; no Odoo corporate authorship"
license: "LGPL-3.0-only"
sha256: "5c0ae67ca00d31f41085b70258716043faac5e0a20953de3646b7e3559083941"
variables: []
secrets_allowed: false
```

````python
# AUTHORED bounded JSON/IPC and representation glue. Source-derived rules live in engine.py.
import hashlib
import json
import re
import sys
from decimal import Decimal, localcontext
from pathlib import Path

from . import engine
from .orm_contract import Record, Records, build, coupon_for

REVISION = '99edb6dd82b7b560930c00b03b694ba700785370'
MAX_BYTES = 131072
DECIMAL = re.compile(r'-?(?:0|[1-9][0-9]{0,17})(?:\.[0-9]{1,6})?\Z')
IDENTITY = re.compile(r'[A-Za-z0-9][A-Za-z0-9_.:-]{0,127}\Z')


def exact(value, fields):
    if type(value) is not dict or set(value) != set(fields.split()):
        raise ValueError('unexpected or missing fields')


def identity(value):
    if type(value) is not str or not IDENTITY.fullmatch(value):
        raise ValueError('invalid identity')
    return value


def decimal(value, minimum=None, positive=False):
    if type(value) is not str or not DECIMAL.fullmatch(value):
        raise ValueError('bounded canonical decimal string required')
    number = Decimal(value)
    if number == 0 and value.startswith('-'):
        raise ValueError('negative zero is not canonical')
    if minimum is not None and number < minimum or positive and number <= 0:
        raise ValueError('decimal outside permitted range')
    return number


def flag(value):
    if type(value) is not bool:
        raise ValueError('boolean required')


def identities(values, maximum=100):
    if type(values) is not list or len(values) > maximum:
        raise ValueError('identity list exceeds limit')
    for value in values:
        identity(value)
    if len(set(values)) != len(values):
        raise ValueError('duplicate identity')


def canonical(value):
    return json.dumps(value, sort_keys=True, separators=(',', ':'), ensure_ascii=False).encode('utf-8')


def profile(program):
    exact(program, 'id kind currency currency_digits applies_on nominative trigger trigger_product_ids rules')
    identity(program['id'])
    if program['kind'] not in ('gift_card', 'loyalty') or program['applies_on'] not in ('future', 'current', 'both'):
        raise ValueError('unsupported program kind/application')
    if not re.fullmatch('[A-Z]{3}', program['currency']) or type(program['currency_digits']) is not int or program['currency_digits'] not in range(5):
        raise ValueError('invalid currency contract')
    if program['trigger'] not in ('auto', 'with_code'):
        raise ValueError('unsupported trigger')
    flag(program['nominative'])
    identities(program['trigger_product_ids'])
    rules = program['rules']
    if type(rules) is not list or not 1 <= len(rules) <= 20:
        raise ValueError('bounded nonempty rules required')
    ids = []
    for rule in rules:
        exact(rule, 'id mode points split minimum_qty minimum_amount tax_mode product_ids code_required')
        ids.append(identity(rule['id']))
        if rule['mode'] not in ('order', 'money', 'unit') or rule['tax_mode'] not in ('incl', 'excl'):
            raise ValueError('unsupported rule mode')
        decimal(rule['points'], positive=True)
        decimal(rule['minimum_qty'], minimum=0)
        decimal(rule['minimum_amount'], minimum=0)
        flag(rule['split'])
        flag(rule['code_required'])
        identities(rule['product_ids'])
        # Directly reflects LoyaltyRule._constraint_trigger_multi; not a new policy.
        if rule['split'] and program['applies_on'] == 'both':
            raise ValueError('split per unit is not allowed when applies_on is both')
    if len(ids) != len(set(ids)):
        raise ValueError('duplicate rule')


def order_contract(order, program):
    exact(order, 'order_id organization_id subject_id state public_subject currency total enabled_rule_ids lines')
    for key in ('order_id', 'organization_id', 'subject_id'):
        identity(order[key])
    if order['state'] not in ('draft', 'sent', 'sale', 'done') or order['currency'] != program['currency']:
        raise ValueError('state or currency binding differs')
    flag(order['public_subject'])
    identities(order['enabled_rule_ids'], 20)
    if not set(order['enabled_rule_ids']) <= {r['id'] for r in program['rules']}:
        raise ValueError('unknown enabled rule')
    decimal(order['total'], minimum=0)
    if type(order['lines']) is not list or not 1 <= len(order['lines']) <= 100:
        raise ValueError('bounded order lines required')
    ids = []
    for line in order['lines']:
        exact(line, 'id product_id quantity subtotal tax total reward_program_type reward_program_id reward_trigger threshold_excluded')
        ids.append(identity(line['id']))
        identity(line['product_id'])
        quantity = decimal(line['quantity'], minimum=0)
        if quantity != quantity.to_integral_value() or quantity > 10000:
            raise ValueError('only bounded normalized integer quantities admitted')
        subtotal, tax, total = (decimal(line[k]) for k in ('subtotal', 'tax', 'total'))
        if subtotal + tax != total:
            raise ValueError('line amounts do not bind')
        flag(line['threshold_excluded'])
        if line['reward_program_type'] not in ('', 'loyalty', 'gift_card', 'ewallet', 'promotion'):
            raise ValueError('unknown reward program')
        if line['reward_program_type']:
            identity(line['reward_program_id'])
            if line['reward_trigger'] not in ('auto', 'with_code'):
                raise ValueError('invalid reward trigger')
        elif line['reward_program_id'] != '' or line['reward_trigger'] != '' or total < 0:
            raise ValueError('normal line reward binding differs')
    if len(set(ids)) != len(ids) or sum(decimal(line['total']) for line in order['lines']) != decimal(order['total']):
        raise ValueError('duplicate line or order total mismatch')


def stringify(value):
    if isinstance(value, Decimal):
        if not value.is_finite():
            raise ValueError('non-finite result')
        return format(value, 'f')
    if type(value) is dict:
        return {str(k): (v if k == 'applied_minor_units' and type(v) is int else stringify(v)) for k,v in value.items()}
    if type(value) is list:
        return [stringify(v) for v in value]
    if type(value) is int and type(value) is not bool:
        return str(value)
    return value


def calculate(request):
    exact(request, 'schema operation program program_sha256 order data')
    if request['schema'] != 'elite.odoo-loyalty-calc.v1' or request['operation'] not in ('evaluate', 'reward', 'changes', 'reverse'):
        raise ValueError('unsupported contract')
    profile(request['program'])
    expected = hashlib.sha256(canonical(request['program'])).hexdigest()
    if request['program_sha256'] != expected:
        raise ValueError('program hash mismatch')
    with localcontext() as context:
        context.prec = 80
        order_contract(request['order'], request['program'])
        program, order = build(request['program'], request['order'])
        data = request['data']
        if request['operation'] == 'evaluate':
            exact(data, '')
            result = order._program_check_compute_points(Records([program]))[program]
            if 'error' not in result:
                result['points'] = engine.filtered_issuance(result, program)
        elif request['operation'] == 'reward':
            exact(data, 'coupon_id balance pending_earned pending_cost discountable reward')
            identity(data['coupon_id'])
            for key in ('balance', 'pending_earned', 'pending_cost'):
                decimal(data[key])
            discountable = decimal(data['discountable'], minimum=0)
            if discountable > order.amount_total:
                raise ValueError('discountable exceeds bound order total')
            reward_data = data['reward']
            exact(reward_data, 'mode discount required_points max_amount clear_wallet')
            if reward_data['mode'] not in ('per_point', 'per_order', 'percent'):
                raise ValueError('unsupported reward mode')
            decimal(reward_data['discount'], positive=True)
            decimal(reward_data['required_points'], positive=True)
            decimal(reward_data['max_amount'], minimum=0)
            flag(reward_data['clear_wallet'])
            coupon = coupon_for(program, data['coupon_id'], data['balance'])
            order.coupon_point_ids = Records([Record('coupon_point', 'pending', coupon_id=coupon, points=Decimal(data['pending_earned']))])
            cost = Record('cost', 'pending', coupon_id=coupon, points_cost=Decimal(data['pending_cost']), reward_id=True)
            order.order_line = order.order_line | Records([cost])
            reward = Record('reward', 'requested', program_id=program, currency_id=program.currency_id,
                            discount_mode=reward_data['mode'], discount=Decimal(reward_data['discount']),
                            required_points=Decimal(reward_data['required_points']), discount_max_amount=Decimal(reward_data['max_amount']),
                            clear_wallet=reward_data['clear_wallet'])
            # Same eligibility precondition as _get_claimable_rewards before invoking discount calculation.
            if order._get_real_points_for_coupon(coupon) < reward.required_points or discountable == 0:
                raise ValueError('reward is not claimable')
            result = engine.reward_from_bound_amounts(order, reward, coupon, discountable)
            result['value'] = program.currency_id.round(result['value'])
        elif request['operation'] == 'reverse':
            exact(data, 'entries')
            if type(data['entries']) is not list or not 1 <= len(data['entries']) <= 40:
                raise ValueError('bounded reversal entries required')
            result = []
            ids = set()
            for entry in data['entries']:
                exact(entry, 'account_id points_delta applied_minor_units')
                identity(entry['account_id'])
                if entry['account_id'] in ids:
                    raise ValueError('duplicate account reversal')
                ids.add(entry['account_id'])
                value = decimal(entry['points_delta'])
                minor = entry['applied_minor_units']
                if type(minor) is not int or not -(2**63)+1 <= minor <= (2**63)-1:
                    raise ValueError('bounded exact minor units required')
                result.append({'account_id':entry['account_id'],
                               'points_delta':engine.reverse_point_change(Decimal('0'), value),
                               'applied_minor_units':-minor})
        else:
            exact(data, 'issued used')
            coupons = {}
            for key in ('issued', 'used'):
                if type(data[key]) is not list or len(data[key]) > 100:
                    raise ValueError('changes exceed limit')
                for item in data[key]:
                    exact(item, 'coupon_id points')
                    identity(item['coupon_id'])
                    decimal(item['points'], minimum=0)
                    coupons.setdefault(item['coupon_id'], coupon_for(program, item['coupon_id'], '0'))
            order.coupon_point_ids = Records(Record('earned', str(i), coupon_id=coupons[x['coupon_id']], points=Decimal(x['points'])) for i,x in enumerate(data['issued']))
            order.order_line = Records(Record('spent', str(i), coupon_id=coupons[x['coupon_id']], points_cost=Decimal(x['points']), reward_id=True) for i,x in enumerate(data['used']))
            result = {coupon.id: amount for coupon,amount in order._get_point_changes().items()}
        return {'schema':'elite.odoo-loyalty-result.v1','source_revision':REVISION,
                'program_sha256':expected,'request_sha256':hashlib.sha256(canonical(request)).hexdigest(),
                'result':stringify(result)}


def pairs(entries):
    result = {}
    for key,value in entries:
        if key in result or key.casefold() in (x.casefold() for x in result):
            raise ValueError('duplicate key')
        result[key] = value
    return result


def main():
    try:
        raw = sys.stdin.buffer.read(MAX_BYTES+1)
        if len(raw) > MAX_BYTES:
            raise ValueError('request exceeds limit')
        request = json.loads(raw.decode('utf-8'), object_pairs_hook=pairs,
                             parse_constant=lambda _value: (_ for _ in ()).throw(ValueError('nonfinite JSON')))
        output = canonical(calculate(request))
        if len(output) > MAX_BYTES:
            raise ValueError('response exceeds limit')
        sys.stdout.buffer.write(output+b'\n')
        return 0
    except (ValueError, KeyError, TypeError, ArithmeticError):
        # No request data or upstream internals appear in errors.
        sys.stdout.buffer.write(b'{"error":"invalid_or_unsupported_calculation"}\n')
        return 2


if __name__ == '__main__':
    raise SystemExit(main())
````

### FILE: `odoo_loyalty/run.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file13:v1"
operation: CREATE
provenance: AUTHORED
source: "local source manifest, isolated loader, ORM/protocol and oracle glue distributed as part of the modified LGPL module; no Odoo corporate authorship"
license: "LGPL-3.0-only"
sha256: "2101f3e889607a8da2fd22722d64c20fd481fa2c38882e1495aeffd72bb3387e"
variables: []
secrets_allowed: false
```

````python
# AUTHORED isolated launcher and exact source-manifest binding.
import hashlib
import json
from pathlib import Path
import sys
from types import ModuleType


def run():
    root = Path(__file__).resolve().parent
    lock_path = root/'engine-lock.json'
    raw = lock_path.read_bytes()
    if len(sys.argv) != 2 or len(raw) > 65536 or hashlib.sha256(raw).hexdigest() != sys.argv[1]:
        return 2
    lock = json.loads(raw)
    if lock['schema'] != 'elite.odoo-derived-source-lock.v1':
        return 2
    sources = {}
    for record in lock['files']:
        relative = Path(record['path'])
        if relative.is_absolute() or '..' in relative.parts:
            return 2
        path = (root/relative).resolve()
        if not path.is_relative_to(root) or path.is_symlink():
            return 2
        data = path.read_bytes()
        if len(data) != record['bytes'] or hashlib.sha256(data).hexdigest() != record['sha256']:
            return 2
        if record['path'] in sources:
            return 2
        sources[record['path']] = data
    # Do not use filesystem package/bytecode import after checking source. A
    # new __init__.py, .pyc or changed file must not replace verified bytes.
    package = ModuleType('odoo_loyalty')
    package.__path__ = []
    package.__package__ = 'odoo_loyalty'
    sys.modules['odoo_loyalty'] = package
    for name in ('engine', 'orm_contract', 'protocol'):
        relative = name + '.py'
        source = sources[relative]
        module = ModuleType('odoo_loyalty.' + name)
        module.__package__ = 'odoo_loyalty'
        module.__file__ = str(root/relative)
        sys.modules[module.__name__] = module
        setattr(package, name, module)
        exec(compile(source, module.__file__, 'exec', dont_inherit=True), module.__dict__)
    return package.protocol.main()


if __name__ == '__main__':
    try:
        raise SystemExit(run())
    except (OSError,ValueError,KeyError,TypeError):
        raise SystemExit(2)
````

### FILE: `odoo_loyalty/upstream/COPYRIGHT`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file14:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/COPYRIGHT; revision 99edb6dd82b7b560930c00b03b694ba700785370; Git blob b0f9743d204484b1080d339ec55a638b1cde32dc"
license: "LGPL-3.0-only"
sha256: "86e49232d2162708d05405ed5ff6dc5594b73ef4b0f25d2749cfae13493b620c"
variables: []
secrets_allowed: false
```

````text

Most of the files are

  Copyright (c) 2004-2015 Odoo S.A.

Many files also contain contributions from third
parties. In this case the original copyright of
the contributions can be traced through the
history of the source version control system.

When that is not the case, the files contain a prominent
notice stating the original copyright and applicable
license, or come with their own dedicated COPYRIGHT
and/or LICENSE file.

````

### FILE: `odoo_loyalty/upstream/LICENSE`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file15:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/LICENSE; revision 99edb6dd82b7b560930c00b03b694ba700785370; Git blob 12334e0fdd8d3412d5e1deed72220fcc03431ff2"
license: "LGPL-3.0-only"
sha256: "abc09dad5f84a76e1b0279237053cae16c03228ab27d8d467677054c2bd17eeb"
variables: []
secrets_allowed: false
final_newline: false
```

````text

For copyright information, please see the COPYRIGHT file.

Odoo is published under the GNU LESSER GENERAL PUBLIC LICENSE, Version 3
(LGPLv3), as included below. Since the LGPL is a set of additional
permissions on top of the GPL, the text of the GPL is included at the
bottom as well.

Some external libraries and contributions bundled with Odoo may be published
under other GPL-compatible licenses. For these, please refer to the relevant
source files and/or license files, in the source code tree.

**************************************************************************

                   GNU LESSER GENERAL PUBLIC LICENSE
                       Version 3, 29 June 2007

 Copyright (C) 2007 Free Software Foundation, Inc. <http://fsf.org/>
 Everyone is permitted to copy and distribute verbatim copies
 of this license document, but changing it is not allowed.


  This version of the GNU Lesser General Public License incorporates
the terms and conditions of version 3 of the GNU General Public
License, supplemented by the additional permissions listed below.

  0. Additional Definitions.

  As used herein, "this License" refers to version 3 of the GNU Lesser
General Public License, and the "GNU GPL" refers to version 3 of the GNU
General Public License.

  "The Library" refers to a covered work governed by this License,
other than an Application or a Combined Work as defined below.

  An "Application" is any work that makes use of an interface provided
by the Library, but which is not otherwise based on the Library.
Defining a subclass of a class defined by the Library is deemed a mode
of using an interface provided by the Library.

  A "Combined Work" is a work produced by combining or linking an
Application with the Library.  The particular version of the Library
with which the Combined Work was made is also called the "Linked
Version".

  The "Minimal Corresponding Source" for a Combined Work means the
Corresponding Source for the Combined Work, excluding any source code
for portions of the Combined Work that, considered in isolation, are
based on the Application, and not on the Linked Version.

  The "Corresponding Application Code" for a Combined Work means the
object code and/or source code for the Application, including any data
and utility programs needed for reproducing the Combined Work from the
Application, but excluding the System Libraries of the Combined Work.

  1. Exception to Section 3 of the GNU GPL.

  You may convey a covered work under sections 3 and 4 of this License
without being bound by section 3 of the GNU GPL.

  2. Conveying Modified Versions.

  If you modify a copy of the Library, and, in your modifications, a
facility refers to a function or data to be supplied by an Application
that uses the facility (other than as an argument passed when the
facility is invoked), then you may convey a copy of the modified
version:

   a) under this License, provided that you make a good faith effort to
   ensure that, in the event an Application does not supply the
   function or data, the facility still operates, and performs
   whatever part of its purpose remains meaningful, or

   b) under the GNU GPL, with none of the additional permissions of
   this License applicable to that copy.

  3. Object Code Incorporating Material from Library Header Files.

  The object code form of an Application may incorporate material from
a header file that is part of the Library.  You may convey such object
code under terms of your choice, provided that, if the incorporated
material is not limited to numerical parameters, data structure
layouts and accessors, or small macros, inline functions and templates
(ten or fewer lines in length), you do both of the following:

   a) Give prominent notice with each copy of the object code that the
   Library is used in it and that the Library and its use are
   covered by this License.

   b) Accompany the object code with a copy of the GNU GPL and this license
   document.

  4. Combined Works.

  You may convey a Combined Work under terms of your choice that,
taken together, effectively do not restrict modification of the
portions of the Library contained in the Combined Work and reverse
engineering for debugging such modifications, if you also do each of
the following:

   a) Give prominent notice with each copy of the Combined Work that
   the Library is used in it and that the Library and its use are
   covered by this License.

   b) Accompany the Combined Work with a copy of the GNU GPL and this license
   document.

   c) For a Combined Work that displays copyright notices during
   execution, include the copyright notice for the Library among
   these notices, as well as a reference directing the user to the
   copies of the GNU GPL and this license document.

   d) Do one of the following:

       0) Convey the Minimal Corresponding Source under the terms of this
       License, and the Corresponding Application Code in a form
       suitable for, and under terms that permit, the user to
       recombine or relink the Application with a modified version of
       the Linked Version to produce a modified Combined Work, in the
       manner specified by section 6 of the GNU GPL for conveying
       Corresponding Source.

       1) Use a suitable shared library mechanism for linking with the
       Library.  A suitable mechanism is one that (a) uses at run time
       a copy of the Library already present on the user's computer
       system, and (b) will operate properly with a modified version
       of the Library that is interface-compatible with the Linked
       Version.

   e) Provide Installation Information, but only if you would otherwise
   be required to provide such information under section 6 of the
   GNU GPL, and only to the extent that such information is
   necessary to install and execute a modified version of the
   Combined Work produced by recombining or relinking the
   Application with a modified version of the Linked Version. (If
   you use option 4d0, the Installation Information must accompany
   the Minimal Corresponding Source and Corresponding Application
   Code. If you use option 4d1, you must provide the Installation
   Information in the manner specified by section 6 of the GNU GPL
   for conveying Corresponding Source.)

  5. Combined Libraries.

  You may place library facilities that are a work based on the
Library side by side in a single library together with other library
facilities that are not Applications and are not covered by this
License, and convey such a combined library under terms of your
choice, if you do both of the following:

   a) Accompany the combined library with a copy of the same work based
   on the Library, uncombined with any other library facilities,
   conveyed under the terms of this License.

   b) Give prominent notice with the combined library that part of it
   is a work based on the Library, and explaining where to find the
   accompanying uncombined form of the same work.

  6. Revised Versions of the GNU Lesser General Public License.

  The Free Software Foundation may publish revised and/or new versions
of the GNU Lesser General Public License from time to time. Such new
versions will be similar in spirit to the present version, but may
differ in detail to address new problems or concerns.

  Each version is given a distinguishing version number. If the
Library as you received it specifies that a certain numbered version
of the GNU Lesser General Public License "or any later version"
applies to it, you have the option of following the terms and
conditions either of that published version or of any later version
published by the Free Software Foundation. If the Library as you
received it does not specify a version number of the GNU Lesser
General Public License, you may choose any version of the GNU Lesser
General Public License ever published by the Free Software Foundation.

  If the Library as you received it specifies that a proxy can decide
whether future versions of the GNU Lesser General Public License shall
apply, that proxy's public statement of acceptance of any version is
permanent authorization for you to choose that version for the
Library.

**************************************************************************

                    GNU GENERAL PUBLIC LICENSE
                       Version 3, 29 June 2007

 Copyright (C) 2007 Free Software Foundation, Inc. <http://fsf.org/>
 Everyone is permitted to copy and distribute verbatim copies
 of this license document, but changing it is not allowed.

                            Preamble

  The GNU General Public License is a free, copyleft license for
software and other kinds of works.

  The licenses for most software and other practical works are designed
to take away your freedom to share and change the works.  By contrast,
the GNU General Public License is intended to guarantee your freedom to
share and change all versions of a program--to make sure it remains free
software for all its users.  We, the Free Software Foundation, use the
GNU General Public License for most of our software; it applies also to
any other work released this way by its authors.  You can apply it to
your programs, too.

  When we speak of free software, we are referring to freedom, not
price.  Our General Public Licenses are designed to make sure that you
have the freedom to distribute copies of free software (and charge for
them if you wish), that you receive source code or can get it if you
want it, that you can change the software or use pieces of it in new
free programs, and that you know you can do these things.

  To protect your rights, we need to prevent others from denying you
these rights or asking you to surrender the rights.  Therefore, you have
certain responsibilities if you distribute copies of the software, or if
you modify it: responsibilities to respect the freedom of others.

  For example, if you distribute copies of such a program, whether
gratis or for a fee, you must pass on to the recipients the same
freedoms that you received.  You must make sure that they, too, receive
or can get the source code.  And you must show them these terms so they
know their rights.

  Developers that use the GNU GPL protect your rights with two steps:
(1) assert copyright on the software, and (2) offer you this License
giving you legal permission to copy, distribute and/or modify it.

  For the developers' and authors' protection, the GPL clearly explains
that there is no warranty for this free software.  For both users' and
authors' sake, the GPL requires that modified versions be marked as
changed, so that their problems will not be attributed erroneously to
authors of previous versions.

  Some devices are designed to deny users access to install or run
modified versions of the software inside them, although the manufacturer
can do so.  This is fundamentally incompatible with the aim of
protecting users' freedom to change the software.  The systematic
pattern of such abuse occurs in the area of products for individuals to
use, which is precisely where it is most unacceptable.  Therefore, we
have designed this version of the GPL to prohibit the practice for those
products.  If such problems arise substantially in other domains, we
stand ready to extend this provision to those domains in future versions
of the GPL, as needed to protect the freedom of users.

  Finally, every program is threatened constantly by software patents.
States should not allow patents to restrict development and use of
software on general-purpose computers, but in those that do, we wish to
avoid the special danger that patents applied to a free program could
make it effectively proprietary.  To prevent this, the GPL assures that
patents cannot be used to render the program non-free.

  The precise terms and conditions for copying, distribution and
modification follow.

                       TERMS AND CONDITIONS

  0. Definitions.

  "This License" refers to version 3 of the GNU General Public License.

  "Copyright" also means copyright-like laws that apply to other kinds of
works, such as semiconductor masks.

  "The Program" refers to any copyrightable work licensed under this
License.  Each licensee is addressed as "you".  "Licensees" and
"recipients" may be individuals or organizations.

  To "modify" a work means to copy from or adapt all or part of the work
in a fashion requiring copyright permission, other than the making of an
exact copy.  The resulting work is called a "modified version" of the
earlier work or a work "based on" the earlier work.

  A "covered work" means either the unmodified Program or a work based
on the Program.

  To "propagate" a work means to do anything with it that, without
permission, would make you directly or secondarily liable for
infringement under applicable copyright law, except executing it on a
computer or modifying a private copy.  Propagation includes copying,
distribution (with or without modification), making available to the
public, and in some countries other activities as well.

  To "convey" a work means any kind of propagation that enables other
parties to make or receive copies.  Mere interaction with a user through
a computer network, with no transfer of a copy, is not conveying.

  An interactive user interface displays "Appropriate Legal Notices"
to the extent that it includes a convenient and prominently visible
feature that (1) displays an appropriate copyright notice, and (2)
tells the user that there is no warranty for the work (except to the
extent that warranties are provided), that licensees may convey the
work under this License, and how to view a copy of this License.  If
the interface presents a list of user commands or options, such as a
menu, a prominent item in the list meets this criterion.

  1. Source Code.

  The "source code" for a work means the preferred form of the work
for making modifications to it.  "Object code" means any non-source
form of a work.

  A "Standard Interface" means an interface that either is an official
standard defined by a recognized standards body, or, in the case of
interfaces specified for a particular programming language, one that
is widely used among developers working in that language.

  The "System Libraries" of an executable work include anything, other
than the work as a whole, that (a) is included in the normal form of
packaging a Major Component, but which is not part of that Major
Component, and (b) serves only to enable use of the work with that
Major Component, or to implement a Standard Interface for which an
implementation is available to the public in source code form.  A
"Major Component", in this context, means a major essential component
(kernel, window system, and so on) of the specific operating system
(if any) on which the executable work runs, or a compiler used to
produce the work, or an object code interpreter used to run it.

  The "Corresponding Source" for a work in object code form means all
the source code needed to generate, install, and (for an executable
work) run the object code and to modify the work, including scripts to
control those activities.  However, it does not include the work's
System Libraries, or general-purpose tools or generally available free
programs which are used unmodified in performing those activities but
which are not part of the work.  For example, Corresponding Source
includes interface definition files associated with source files for
the work, and the source code for shared libraries and dynamically
linked subprograms that the work is specifically designed to require,
such as by intimate data communication or control flow between those
subprograms and other parts of the work.

  The Corresponding Source need not include anything that users
can regenerate automatically from other parts of the Corresponding
Source.

  The Corresponding Source for a work in source code form is that
same work.

  2. Basic Permissions.

  All rights granted under this License are granted for the term of
copyright on the Program, and are irrevocable provided the stated
conditions are met.  This License explicitly affirms your unlimited
permission to run the unmodified Program.  The output from running a
covered work is covered by this License only if the output, given its
content, constitutes a covered work.  This License acknowledges your
rights of fair use or other equivalent, as provided by copyright law.

  You may make, run and propagate covered works that you do not
convey, without conditions so long as your license otherwise remains
in force.  You may convey covered works to others for the sole purpose
of having them make modifications exclusively for you, or provide you
with facilities for running those works, provided that you comply with
the terms of this License in conveying all material for which you do
not control copyright.  Those thus making or running the covered works
for you must do so exclusively on your behalf, under your direction
and control, on terms that prohibit them from making any copies of
your copyrighted material outside their relationship with you.

  Conveying under any other circumstances is permitted solely under
the conditions stated below.  Sublicensing is not allowed; section 10
makes it unnecessary.

  3. Protecting Users' Legal Rights From Anti-Circumvention Law.

  No covered work shall be deemed part of an effective technological
measure under any applicable law fulfilling obligations under article
11 of the WIPO copyright treaty adopted on 20 December 1996, or
similar laws prohibiting or restricting circumvention of such
measures.

  When you convey a covered work, you waive any legal power to forbid
circumvention of technological measures to the extent such circumvention
is effected by exercising rights under this License with respect to
the covered work, and you disclaim any intention to limit operation or
modification of the work as a means of enforcing, against the work's
users, your or third parties' legal rights to forbid circumvention of
technological measures.

  4. Conveying Verbatim Copies.

  You may convey verbatim copies of the Program's source code as you
receive it, in any medium, provided that you conspicuously and
appropriately publish on each copy an appropriate copyright notice;
keep intact all notices stating that this License and any
non-permissive terms added in accord with section 7 apply to the code;
keep intact all notices of the absence of any warranty; and give all
recipients a copy of this License along with the Program.

  You may charge any price or no price for each copy that you convey,
and you may offer support or warranty protection for a fee.

  5. Conveying Modified Source Versions.

  You may convey a work based on the Program, or the modifications to
produce it from the Program, in the form of source code under the
terms of section 4, provided that you also meet all of these conditions:

    a) The work must carry prominent notices stating that you modified
    it, and giving a relevant date.

    b) The work must carry prominent notices stating that it is
    released under this License and any conditions added under section
    7.  This requirement modifies the requirement in section 4 to
    "keep intact all notices".

    c) You must license the entire work, as a whole, under this
    License to anyone who comes into possession of a copy.  This
    License will therefore apply, along with any applicable section 7
    additional terms, to the whole of the work, and all its parts,
    regardless of how they are packaged.  This License gives no
    permission to license the work in any other way, but it does not
    invalidate such permission if you have separately received it.

    d) If the work has interactive user interfaces, each must display
    Appropriate Legal Notices; however, if the Program has interactive
    interfaces that do not display Appropriate Legal Notices, your
    work need not make them do so.

  A compilation of a covered work with other separate and independent
works, which are not by their nature extensions of the covered work,
and which are not combined with it such as to form a larger program,
in or on a volume of a storage or distribution medium, is called an
"aggregate" if the compilation and its resulting copyright are not
used to limit the access or legal rights of the compilation's users
beyond what the individual works permit.  Inclusion of a covered work
in an aggregate does not cause this License to apply to the other
parts of the aggregate.

  6. Conveying Non-Source Forms.

  You may convey a covered work in object code form under the terms
of sections 4 and 5, provided that you also convey the
machine-readable Corresponding Source under the terms of this License,
in one of these ways:

    a) Convey the object code in, or embodied in, a physical product
    (including a physical distribution medium), accompanied by the
    Corresponding Source fixed on a durable physical medium
    customarily used for software interchange.

    b) Convey the object code in, or embodied in, a physical product
    (including a physical distribution medium), accompanied by a
    written offer, valid for at least three years and valid for as
    long as you offer spare parts or customer support for that product
    model, to give anyone who possesses the object code either (1) a
    copy of the Corresponding Source for all the software in the
    product that is covered by this License, on a durable physical
    medium customarily used for software interchange, for a price no
    more than your reasonable cost of physically performing this
    conveying of source, or (2) access to copy the
    Corresponding Source from a network server at no charge.

    c) Convey individual copies of the object code with a copy of the
    written offer to provide the Corresponding Source.  This
    alternative is allowed only occasionally and noncommercially, and
    only if you received the object code with such an offer, in accord
    with subsection 6b.

    d) Convey the object code by offering access from a designated
    place (gratis or for a charge), and offer equivalent access to the
    Corresponding Source in the same way through the same place at no
    further charge.  You need not require recipients to copy the
    Corresponding Source along with the object code.  If the place to
    copy the object code is a network server, the Corresponding Source
    may be on a different server (operated by you or a third party)
    that supports equivalent copying facilities, provided you maintain
    clear directions next to the object code saying where to find the
    Corresponding Source.  Regardless of what server hosts the
    Corresponding Source, you remain obligated to ensure that it is
    available for as long as needed to satisfy these requirements.

    e) Convey the object code using peer-to-peer transmission, provided
    you inform other peers where the object code and Corresponding
    Source of the work are being offered to the general public at no
    charge under subsection 6d.

  A separable portion of the object code, whose source code is excluded
from the Corresponding Source as a System Library, need not be
included in conveying the object code work.

  A "User Product" is either (1) a "consumer product", which means any
tangible personal property which is normally used for personal, family,
or household purposes, or (2) anything designed or sold for incorporation
into a dwelling.  In determining whether a product is a consumer product,
doubtful cases shall be resolved in favor of coverage.  For a particular
product received by a particular user, "normally used" refers to a
typical or common use of that class of product, regardless of the status
of the particular user or of the way in which the particular user
actually uses, or expects or is expected to use, the product.  A product
is a consumer product regardless of whether the product has substantial
commercial, industrial or non-consumer uses, unless such uses represent
the only significant mode of use of the product.

  "Installation Information" for a User Product means any methods,
procedures, authorization keys, or other information required to install
and execute modified versions of a covered work in that User Product from
a modified version of its Corresponding Source.  The information must
suffice to ensure that the continued functioning of the modified object
code is in no case prevented or interfered with solely because
modification has been made.

  If you convey an object code work under this section in, or with, or
specifically for use in, a User Product, and the conveying occurs as
part of a transaction in which the right of possession and use of the
User Product is transferred to the recipient in perpetuity or for a
fixed term (regardless of how the transaction is characterized), the
Corresponding Source conveyed under this section must be accompanied
by the Installation Information.  But this requirement does not apply
if neither you nor any third party retains the ability to install
modified object code on the User Product (for example, the work has
been installed in ROM).

  The requirement to provide Installation Information does not include a
requirement to continue to provide support service, warranty, or updates
for a work that has been modified or installed by the recipient, or for
the User Product in which it has been modified or installed.  Access to a
network may be denied when the modification itself materially and
adversely affects the operation of the network or violates the rules and
protocols for communication across the network.

  Corresponding Source conveyed, and Installation Information provided,
in accord with this section must be in a format that is publicly
documented (and with an implementation available to the public in
source code form), and must require no special password or key for
unpacking, reading or copying.

  7. Additional Terms.

  "Additional permissions" are terms that supplement the terms of this
License by making exceptions from one or more of its conditions.
Additional permissions that are applicable to the entire Program shall
be treated as though they were included in this License, to the extent
that they are valid under applicable law.  If additional permissions
apply only to part of the Program, that part may be used separately
under those permissions, but the entire Program remains governed by
this License without regard to the additional permissions.

  When you convey a copy of a covered work, you may at your option
remove any additional permissions from that copy, or from any part of
it.  (Additional permissions may be written to require their own
removal in certain cases when you modify the work.)  You may place
additional permissions on material, added by you to a covered work,
for which you have or can give appropriate copyright permission.

  Notwithstanding any other provision of this License, for material you
add to a covered work, you may (if authorized by the copyright holders of
that material) supplement the terms of this License with terms:

    a) Disclaiming warranty or limiting liability differently from the
    terms of sections 15 and 16 of this License; or

    b) Requiring preservation of specified reasonable legal notices or
    author attributions in that material or in the Appropriate Legal
    Notices displayed by works containing it; or

    c) Prohibiting misrepresentation of the origin of that material, or
    requiring that modified versions of such material be marked in
    reasonable ways as different from the original version; or

    d) Limiting the use for publicity purposes of names of licensors or
    authors of the material; or

    e) Declining to grant rights under trademark law for use of some
    trade names, trademarks, or service marks; or

    f) Requiring indemnification of licensors and authors of that
    material by anyone who conveys the material (or modified versions of
    it) with contractual assumptions of liability to the recipient, for
    any liability that these contractual assumptions directly impose on
    those licensors and authors.

  All other non-permissive additional terms are considered "further
restrictions" within the meaning of section 10.  If the Program as you
received it, or any part of it, contains a notice stating that it is
governed by this License along with a term that is a further
restriction, you may remove that term.  If a license document contains
a further restriction but permits relicensing or conveying under this
License, you may add to a covered work material governed by the terms
of that license document, provided that the further restriction does
not survive such relicensing or conveying.

  If you add terms to a covered work in accord with this section, you
must place, in the relevant source files, a statement of the
additional terms that apply to those files, or a notice indicating
where to find the applicable terms.

  Additional terms, permissive or non-permissive, may be stated in the
form of a separately written license, or stated as exceptions;
the above requirements apply either way.

  8. Termination.

  You may not propagate or modify a covered work except as expressly
provided under this License.  Any attempt otherwise to propagate or
modify it is void, and will automatically terminate your rights under
this License (including any patent licenses granted under the third
paragraph of section 11).

  However, if you cease all violation of this License, then your
license from a particular copyright holder is reinstated (a)
provisionally, unless and until the copyright holder explicitly and
finally terminates your license, and (b) permanently, if the copyright
holder fails to notify you of the violation by some reasonable means
prior to 60 days after the cessation.

  Moreover, your license from a particular copyright holder is
reinstated permanently if the copyright holder notifies you of the
violation by some reasonable means, this is the first time you have
received notice of violation of this License (for any work) from that
copyright holder, and you cure the violation prior to 30 days after
your receipt of the notice.

  Termination of your rights under this section does not terminate the
licenses of parties who have received copies or rights from you under
this License.  If your rights have been terminated and not permanently
reinstated, you do not qualify to receive new licenses for the same
material under section 10.

  9. Acceptance Not Required for Having Copies.

  You are not required to accept this License in order to receive or
run a copy of the Program.  Ancillary propagation of a covered work
occurring solely as a consequence of using peer-to-peer transmission
to receive a copy likewise does not require acceptance.  However,
nothing other than this License grants you permission to propagate or
modify any covered work.  These actions infringe copyright if you do
not accept this License.  Therefore, by modifying or propagating a
covered work, you indicate your acceptance of this License to do so.

  10. Automatic Licensing of Downstream Recipients.

  Each time you convey a covered work, the recipient automatically
receives a license from the original licensors, to run, modify and
propagate that work, subject to this License.  You are not responsible
for enforcing compliance by third parties with this License.

  An "entity transaction" is a transaction transferring control of an
organization, or substantially all assets of one, or subdividing an
organization, or merging organizations.  If propagation of a covered
work results from an entity transaction, each party to that
transaction who receives a copy of the work also receives whatever
licenses to the work the party's predecessor in interest had or could
give under the previous paragraph, plus a right to possession of the
Corresponding Source of the work from the predecessor in interest, if
the predecessor has it or can get it with reasonable efforts.

  You may not impose any further restrictions on the exercise of the
rights granted or affirmed under this License.  For example, you may
not impose a license fee, royalty, or other charge for exercise of
rights granted under this License, and you may not initiate litigation
(including a cross-claim or counterclaim in a lawsuit) alleging that
any patent claim is infringed by making, using, selling, offering for
sale, or importing the Program or any portion of it.

  11. Patents.

  A "contributor" is a copyright holder who authorizes use under this
License of the Program or a work on which the Program is based.  The
work thus licensed is called the contributor's "contributor version".

  A contributor's "essential patent claims" are all patent claims
owned or controlled by the contributor, whether already acquired or
hereafter acquired, that would be infringed by some manner, permitted
by this License, of making, using, or selling its contributor version,
but do not include claims that would be infringed only as a
consequence of further modification of the contributor version.  For
purposes of this definition, "control" includes the right to grant
patent sublicenses in a manner consistent with the requirements of
this License.

  Each contributor grants you a non-exclusive, worldwide, royalty-free
patent license under the contributor's essential patent claims, to
make, use, sell, offer for sale, import and otherwise run, modify and
propagate the contents of its contributor version.

  In the following three paragraphs, a "patent license" is any express
agreement or commitment, however denominated, not to enforce a patent
(such as an express permission to practice a patent or covenant not to
sue for patent infringement).  To "grant" such a patent license to a
party means to make such an agreement or commitment not to enforce a
patent against the party.

  If you convey a covered work, knowingly relying on a patent license,
and the Corresponding Source of the work is not available for anyone
to copy, free of charge and under the terms of this License, through a
publicly available network server or other readily accessible means,
then you must either (1) cause the Corresponding Source to be so
available, or (2) arrange to deprive yourself of the benefit of the
patent license for this particular work, or (3) arrange, in a manner
consistent with the requirements of this License, to extend the patent
license to downstream recipients.  "Knowingly relying" means you have
actual knowledge that, but for the patent license, your conveying the
covered work in a country, or your recipient's use of the covered work
in a country, would infringe one or more identifiable patents in that
country that you have reason to believe are valid.

  If, pursuant to or in connection with a single transaction or
arrangement, you convey, or propagate by procuring conveyance of, a
covered work, and grant a patent license to some of the parties
receiving the covered work authorizing them to use, propagate, modify
or convey a specific copy of the covered work, then the patent license
you grant is automatically extended to all recipients of the covered
work and works based on it.

  A patent license is "discriminatory" if it does not include within
the scope of its coverage, prohibits the exercise of, or is
conditioned on the non-exercise of one or more of the rights that are
specifically granted under this License.  You may not convey a covered
work if you are a party to an arrangement with a third party that is
in the business of distributing software, under which you make payment
to the third party based on the extent of your activity of conveying
the work, and under which the third party grants, to any of the
parties who would receive the covered work from you, a discriminatory
patent license (a) in connection with copies of the covered work
conveyed by you (or copies made from those copies), or (b) primarily
for and in connection with specific products or compilations that
contain the covered work, unless you entered into that arrangement,
or that patent license was granted, prior to 28 March 2007.

  Nothing in this License shall be construed as excluding or limiting
any implied license or other defenses to infringement that may
otherwise be available to you under applicable patent law.

  12. No Surrender of Others' Freedom.

  If conditions are imposed on you (whether by court order, agreement or
otherwise) that contradict the conditions of this License, they do not
excuse you from the conditions of this License.  If you cannot convey a
covered work so as to satisfy simultaneously your obligations under this
License and any other pertinent obligations, then as a consequence you may
not convey it at all.  For example, if you agree to terms that obligate you
to collect a royalty for further conveying from those to whom you convey
the Program, the only way you could satisfy both those terms and this
License would be to refrain entirely from conveying the Program.

  13. Use with the GNU Affero General Public License.

  Notwithstanding any other provision of this License, you have
permission to link or combine any covered work with a work licensed
under version 3 of the GNU Affero General Public License into a single
combined work, and to convey the resulting work.  The terms of this
License will continue to apply to the part which is the covered work,
but the special requirements of the GNU Affero General Public License,
section 13, concerning interaction through a network will apply to the
combination as such.

  14. Revised Versions of this License.

  The Free Software Foundation may publish revised and/or new versions of
the GNU General Public License from time to time.  Such new versions will
be similar in spirit to the present version, but may differ in detail to
address new problems or concerns.

  Each version is given a distinguishing version number.  If the
Program specifies that a certain numbered version of the GNU General
Public License "or any later version" applies to it, you have the
option of following the terms and conditions either of that numbered
version or of any later version published by the Free Software
Foundation.  If the Program does not specify a version number of the
GNU General Public License, you may choose any version ever published
by the Free Software Foundation.

  If the Program specifies that a proxy can decide which future
versions of the GNU General Public License can be used, that proxy's
public statement of acceptance of a version permanently authorizes you
to choose that version for the Program.

  Later license versions may give you additional or different
permissions.  However, no additional obligations are imposed on any
author or copyright holder as a result of your choosing to follow a
later version.

  15. Disclaimer of Warranty.

  THERE IS NO WARRANTY FOR THE PROGRAM, TO THE EXTENT PERMITTED BY
APPLICABLE LAW.  EXCEPT WHEN OTHERWISE STATED IN WRITING THE COPYRIGHT
HOLDERS AND/OR OTHER PARTIES PROVIDE THE PROGRAM "AS IS" WITHOUT WARRANTY
OF ANY KIND, EITHER EXPRESSED OR IMPLIED, INCLUDING, BUT NOT LIMITED TO,
THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR
PURPOSE.  THE ENTIRE RISK AS TO THE QUALITY AND PERFORMANCE OF THE PROGRAM
IS WITH YOU.  SHOULD THE PROGRAM PROVE DEFECTIVE, YOU ASSUME THE COST OF
ALL NECESSARY SERVICING, REPAIR OR CORRECTION.

  16. Limitation of Liability.

  IN NO EVENT UNLESS REQUIRED BY APPLICABLE LAW OR AGREED TO IN WRITING
WILL ANY COPYRIGHT HOLDER, OR ANY OTHER PARTY WHO MODIFIES AND/OR CONVEYS
THE PROGRAM AS PERMITTED ABOVE, BE LIABLE TO YOU FOR DAMAGES, INCLUDING ANY
GENERAL, SPECIAL, INCIDENTAL OR CONSEQUENTIAL DAMAGES ARISING OUT OF THE
USE OR INABILITY TO USE THE PROGRAM (INCLUDING BUT NOT LIMITED TO LOSS OF
DATA OR DATA BEING RENDERED INACCURATE OR LOSSES SUSTAINED BY YOU OR THIRD
PARTIES OR A FAILURE OF THE PROGRAM TO OPERATE WITH ANY OTHER PROGRAMS),
EVEN IF SUCH HOLDER OR OTHER PARTY HAS BEEN ADVISED OF THE POSSIBILITY OF
SUCH DAMAGES.

  17. Interpretation of Sections 15 and 16.

  If the disclaimer of warranty and limitation of liability provided
above cannot be given local legal effect according to their terms,
reviewing courts shall apply local law that most closely approximates
an absolute waiver of all civil liability in connection with the
Program, unless a warranty or assumption of liability accompanies a
copy of the Program in return for a fee.

                     END OF TERMS AND CONDITIONS

            How to Apply These Terms to Your New Programs

  If you develop a new program, and you want it to be of the greatest
possible use to the public, the best way to achieve this is to make it
free software which everyone can redistribute and change under these terms.

  To do so, attach the following notices to the program.  It is safest
to attach them to the start of each source file to most effectively
state the exclusion of warranty; and each file should have at least
the "copyright" line and a pointer to where the full notice is found.

    <one line to give the program's name and a brief idea of what it does.>
    Copyright (C) <year>  <name of author>

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.

Also add information on how to contact you by electronic and paper mail.

  If the program does terminal interaction, make it output a short
notice like this when it starts in an interactive mode:

    <program>  Copyright (C) <year>  <name of author>
    This program comes with ABSOLUTELY NO WARRANTY; for details type `show w'.
    This is free software, and you are welcome to redistribute it
    under certain conditions; type `show c' for details.

The hypothetical commands `show w' and `show c' should show the appropriate
parts of the General Public License.  Of course, your program's commands
might be different; for a GUI interface, you would use an "about box".

  You should also get your employer (if you work as a programmer) or school,
if any, to sign a "copyright disclaimer" for the program, if necessary.
For more information on this, and how to apply and follow the GNU GPL, see
<http://www.gnu.org/licenses/>.

  The GNU General Public License does not permit incorporating your program
into proprietary programs.  If your program is a subroutine library, you
may consider it more useful to permit linking proprietary applications with
the library.  If this is what you want to do, use the GNU Lesser General
Public License instead of this License.  But first, please read
<http://www.gnu.org/philosophy/why-not-lgpl.html>.


**************************************************************************
````

### FILE: `odoo_loyalty/upstream/addons/loyalty/__manifest__.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file16:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/loyalty/__manifest__.py; revision 99edb6dd82b7b560930c00b03b694ba700785370; Git blob fc737c84bb2c46336609a58a5ece87ad77432d65"
license: "LGPL-3.0-only"
sha256: "03537054d1786db61c364262242eff900c6b104c2e1e5c8df3a60cd2b1e403f5"
variables: []
secrets_allowed: false
```

````python
# Part of Odoo. See LICENSE file for full copyright and licensing details.

{
    'name': "Coupons & Loyalty",
    'summary': "Use discounts, gift card, eWallets and loyalty programs in different sales channels",
    'category': 'Sales',
    'version': '1.0',
    'depends': ['product', 'portal', 'account'],
    'data': [
        'security/ir.model.access.csv',
        'security/loyalty_security.xml',
        'report/loyalty_report_templates.xml',
        'report/loyalty_report.xml',
        'data/mail_template_data.xml',
        'data/loyalty_data.xml',
        'wizard/loyalty_card_update_balance_views.xml',
        'wizard/loyalty_generate_wizard_views.xml',
        'views/loyalty_card_views.xml',
        'views/loyalty_history_views.xml',
        'views/loyalty_mail_views.xml',
        'views/loyalty_program_views.xml',
        'views/loyalty_reward_views.xml',
        'views/loyalty_rule_views.xml',
        'views/portal_templates.xml',
        'views/res_partner_views.xml',
    ],
    'demo': [
        'data/loyalty_demo.xml',
    ],
    'assets': {
        'web.assets_backend': [
            'loyalty/static/src/js/**/*.js',
            'loyalty/static/src/scss/*.scss',
            'loyalty/static/src/xml/*.xml',

            ('remove', 'loyalty/static/src/js/portal/**/*'),
            # Don't include dark mode files in light mode
            ('remove', 'loyalty/static/src/scss/*.dark.scss'),
        ],
        "web.assets_web_dark": [
            'loyalty/static/src/scss/*.dark.scss',
        ],
        'web.assets_frontend': [
            'loyalty/static/src/js/portal/**/*',
            'loyalty/static/src/interactions/*',
        ],
    },
    'installable': True,
    'author': 'Odoo S.A.',
    'license': 'LGPL-3',
}
````

### FILE: `odoo_loyalty/upstream/addons/loyalty/models/loyalty_card.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file17:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/loyalty/models/loyalty_card.py; revision 99edb6dd82b7b560930c00b03b694ba700785370; Git blob f0baa0b8bcbe434ff5777c48d2c8ef8a68003a3d"
license: "LGPL-3.0-only"
sha256: "9eead538d45b3ad32a0f2eddc6b146d2a0ab5cc4d40de19bdea58aa70908ba83"
variables: []
secrets_allowed: false
```

````python
# Part of Odoo. See LICENSE file for full copyright and licensing details.

from uuid import uuid4

from odoo import _, api, fields, models
from odoo.exceptions import ValidationError
from odoo.tools import format_amount


class LoyaltyCard(models.Model):
    _name = 'loyalty.card'
    _inherit = ['mail.thread']
    _description = "Loyalty Coupon"
    _rec_name = 'code'

    @api.model
    def _generate_code(self):
        """
        Barcode identifiable codes.
        """
        return "044" + str(uuid4())[7:-18]

    @api.depends('program_id', 'code')
    def _compute_display_name(self):
        for card in self:
            card.display_name = f"{card.program_id.name}: {card.code}"

    program_id = fields.Many2one(
        comodel_name='loyalty.program',
        ondelete='restrict',
        index='btree_not_null',
        default=lambda self: self.env.context.get('active_id', None),
    )
    program_type = fields.Selection(related='program_id.program_type')
    # TODO probably isn't useful to store this company_id anymore
    company_id = fields.Many2one(related='program_id.company_id', store=True, precompute=True)
    currency_id = fields.Many2one(related='program_id.currency_id')
    # Reserved for this partner if non-empty
    partner_id = fields.Many2one(comodel_name='res.partner', index=True)
    points = fields.Float(tracking=True)
    point_name = fields.Char(related='program_id.portal_point_name', readonly=True)
    points_display = fields.Char(compute='_compute_points_display')

    code = fields.Char(required=True, default=lambda self: self._generate_code())
    expiration_date = fields.Date()

    use_count = fields.Integer(compute='_compute_use_count')
    active = fields.Boolean(default=True)
    history_ids = fields.One2many(
        comodel_name='loyalty.history',
        inverse_name='card_id',
        readonly=True,
    )

    _card_code_unique = models.Constraint(
        'UNIQUE(code)',
        "A coupon/loyalty card must have a unique code.",
    )

    @api.constrains('code')
    def _contrains_code(self):
        # Prevent a coupon from having the same code a program
        if self.env['loyalty.rule'].search_count([('mode', '=', 'with_code'), ('code', 'in', self.mapped('code'))]):
            raise ValidationError(_("A trigger with the same code as one of your coupon already exists."))

    @api.depends('points', 'point_name')
    def _compute_points_display(self):
        for card in self:
            card.points_display = card._format_points(card.points)

    @api.onchange('expiration_date')
    def _restrict_expiration_on_loyalty(self):
        for card in self:
            if card.program_type == 'loyalty' and card.expiration_date:
                raise ValidationError(_("Expiration date cannot be set on a loyalty card."))

    def _format_points(self, points):
        self.ensure_one()
        if self.program_id.currency_id and self.point_name == self.program_id.currency_id.symbol:
            return format_amount(self.env, points, self.program_id.currency_id)
        if points == int(points):
            return f"{int(points)} {self.point_name or ''}"
        return f"{points:.2f} {self.point_name or ''}"

    # Meant to be overriden
    def _compute_use_count(self):
        self.use_count = 0

    def _get_default_template(self):
        self.ensure_one()
        return self.program_id.communication_plan_ids.filtered(lambda m: m.trigger == 'create').mail_template_id[:1]

    def _get_mail_author(self):
        self.ensure_one()
        return (
            self.env.user._is_internal() and self.env.user or self.company_id or self.env.company
        ).partner_id

    def _get_signature(self):
        """To be overriden"""
        self.ensure_one()
        return None

    def _has_source_order(self):
        return False

    def action_coupon_send(self):
        """ Open a window to compose an email, with the default template returned by `_get_default_template`
            message loaded by default
        """
        self.ensure_one()
        default_template = self._get_default_template()
        compose_form = self.env.ref('mail.email_compose_message_wizard_form', False)
        ctx = dict(
            default_model='loyalty.card',
            default_res_ids=self.ids,
            default_template_id=default_template and default_template.id,
            default_composition_mode='comment',
            default_email_layout_xmlid='mail.mail_notification_light',
            force_email=True,
        )
        return {
            'name': _("Compose Email"),
            'type': 'ir.actions.act_window',
            'view_mode': 'form',
            'res_model': 'mail.compose.message',
            'views': [(compose_form.id, 'form')],
            'view_id': compose_form.id,
            'target': 'new',
            'context': ctx,
        }

    def _send_creation_communication(self, force_send=False):
        """
        Sends the 'At Creation' communication plan if it exist for the given coupons.
        """
        if self.env.context.get('loyalty_no_mail', False) or self.env.context.get('action_no_send_mail', False):
            return
        # Ideally one per program, but multiple is supported
        create_comm_per_program = dict()
        for program in self.program_id:
            create_comm_per_program[program] = program.communication_plan_ids.filtered(lambda c: c.trigger == 'create')
        for coupon in self:
            if not create_comm_per_program[coupon.program_id] or not coupon._mail_get_customer():
                continue
            for comm in create_comm_per_program[coupon.program_id]:
                mail_template = comm.mail_template_id
                email_values = {}
                if not mail_template.email_from:
                    # provide author_id & email_from values to ensure the email gets sent
                    author = coupon._get_mail_author()
                    email_values.update(author_id=author.id, email_from=author.email_formatted)
                mail_template.send_mail(
                    res_id=coupon.id,
                    force_send=force_send,
                    email_layout_xmlid='mail.mail_notification_light',
                    email_values=email_values,
                )

    def _send_points_reach_communication(self, points_changes):
        """
        Send the 'When Reaching' communicaton plans for the given coupons.

        If a coupons passes multiple milestones we will only send the one with the highest target.
        """
        if self.env.context.get('loyalty_no_mail', False):
            return
        milestones_per_program = dict()
        for program in self.program_id:
            milestones_per_program[program] = program.communication_plan_ids\
                .filtered(lambda c: c.trigger == 'points_reach')\
                .sorted('points', reverse=True)
        for coupon in self:
            if not coupon._mail_get_customer():
                continue
            coupon_change = points_changes[coupon]
            # Do nothing if coupon lost points or did not change
            if not milestones_per_program[coupon.program_id] or\
                not coupon.partner_id or\
                coupon_change['old'] >= coupon_change['new']:
                continue
            this_milestone = False
            for milestone in milestones_per_program[coupon.program_id]:
                if coupon_change['old'] < milestone.points and milestone.points <= coupon_change['new']:
                    this_milestone = milestone
                    break
            if not this_milestone:
                continue
            this_milestone.mail_template_id.send_mail(res_id=coupon.id, email_layout_xmlid='mail.mail_notification_light')


    @api.model_create_multi
    def create(self, vals_list):
        res = super().create(vals_list)
        res._send_creation_communication()
        return res

    def write(self, vals):
        if not self.env.context.get('loyalty_no_mail', False) and 'points' in vals:
            points_before = {coupon: coupon.points for coupon in self}
        res = super().write(vals)
        if not self.env.context.get('loyalty_no_mail', False) and 'points' in vals:
            points_changes = {coupon: {'old': points_before[coupon], 'new': coupon.points} for coupon in self}
            self._send_points_reach_communication(points_changes)
        return res

    def action_loyalty_update_balance(self):
        return {
            'name': _("Update Balance"),
            'type': 'ir.actions.act_window',
            'view_mode': 'form',
            'res_model': 'loyalty.card.update.balance',
            'target': 'new',
            'context': {
                'default_card_id': self.id,
            },
        }
````

### FILE: `odoo_loyalty/upstream/addons/loyalty/models/loyalty_history.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file18:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/loyalty/models/loyalty_history.py; revision 99edb6dd82b7b560930c00b03b694ba700785370; Git blob 5c70b44be15f623fd1e897e6dc30b50e35560f33"
license: "LGPL-3.0-only"
sha256: "c480709def3ad815fe6b306dac534785f24cf10b7432276a0276dcd5646bb8a5"
variables: []
secrets_allowed: false
```

````python
# Part of Odoo. See LICENSE file for full copyright and licensing details.

from odoo import fields, models


class LoyaltyHistory(models.Model):
    _name = 'loyalty.history'
    _description = "History for Loyalty cards and Ewallets"
    _order = 'id desc'

    card_id = fields.Many2one(comodel_name='loyalty.card', required=True, index=True, ondelete='cascade')
    company_id = fields.Many2one(related='card_id.company_id')

    description = fields.Text(required=True)

    issued = fields.Float()
    used = fields.Float()

    order_model = fields.Char(readonly=True)
    order_id = fields.Many2oneReference(model_field='order_model', readonly=True)

    def _get_order_portal_url(self):
        self.ensure_one()
        return False

    def _get_order_description(self):
        self.ensure_one()
        return self.env[self.order_model].browse(self.order_id).display_name
````

### FILE: `odoo_loyalty/upstream/addons/loyalty/models/loyalty_program.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file19:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/loyalty/models/loyalty_program.py; revision 99edb6dd82b7b560930c00b03b694ba700785370; Git blob 239f33d03fe89ab14885d6ef97b96fe8267bbeb3"
license: "LGPL-3.0-only"
sha256: "3ed5484b4404120c73afe810653c2b25915d50d4211158005f684d9290e08096"
variables: []
secrets_allowed: false
```

````python
# Part of Odoo. See LICENSE file for full copyright and licensing details.

from collections import defaultdict
from uuid import uuid4

from odoo import _, api, fields, models
from odoo.exceptions import UserError, ValidationError


class LoyaltyProgram(models.Model):
    _name = 'loyalty.program'
    _description = "Loyalty Program"
    _order = 'sequence'
    _rec_name = 'name'

    @api.model
    def default_get(self, fields):
        defaults = super().default_get(fields)
        program_type = defaults.get('program_type')
        if program_type:
            program_default_values = self._program_type_default_values()
            if program_type in program_default_values:
                default_values = program_default_values[program_type]
                defaults.update({k: v for k, v in default_values.items() if k in fields})
        return defaults

    name = fields.Char(string="Program Name", translate=True, required=True)
    active = fields.Boolean(default=True)
    sequence = fields.Integer(copy=False)
    company_id = fields.Many2one(
        string="Company", comodel_name='res.company', default=lambda self: self.env.company
    )
    currency_id = fields.Many2one(
        string="Currency",
        comodel_name='res.currency',
        compute='_compute_currency_id',
        precompute=True,
        store=True,
        readonly=False,
        required=True,
    )
    currency_symbol = fields.Char(related='currency_id.symbol')
    pricelist_ids = fields.Many2many(
        string="Pricelist",
        help="This program is specific to this pricelist set.",
        comodel_name='product.pricelist',
        domain="[('currency_id', '=', currency_id)]",
    )

    total_order_count = fields.Integer(
        string="Total Order Count", compute='_compute_total_order_count'
    )

    rule_ids = fields.One2many(
        string="Conditional rules",
        comodel_name='loyalty.rule',
        inverse_name='program_id',
        compute='_compute_from_program_type',
        store=True,
        readonly=False,
        copy=True,
    )
    reward_ids = fields.One2many(
        string="Rewards",
        comodel_name='loyalty.reward',
        inverse_name='program_id',
        compute='_compute_from_program_type',
        store=True,
        readonly=False,
        copy=True,
    )
    communication_plan_ids = fields.One2many(
        comodel_name='loyalty.mail',
        inverse_name='program_id',
        compute='_compute_from_program_type',
        store=True,
        readonly=False,
        copy=True,
    )

    # These fields are used for the simplified view of gift_card and ewallet
    mail_template_id = fields.Many2one(
        string="Email template",
        comodel_name='mail.template',
        compute='_compute_mail_template_id',
        inverse='_inverse_mail_template_id',
        readonly=False,
    )
    trigger_product_ids = fields.Many2many(related='rule_ids.product_ids', readonly=False)

    coupon_ids = fields.One2many(comodel_name='loyalty.card', inverse_name='program_id')
    coupon_count = fields.Integer(compute='_compute_coupon_count')
    coupon_count_display = fields.Char(string="Items", compute='_compute_coupon_count_display')

    program_type = fields.Selection(
        selection=[
            ('coupons', "Coupons"),
            ('gift_card', "Gift Card"),
            ('loyalty', "Loyalty Cards"),
            ('promotion', "Promotions"),
            ('ewallet', "eWallet"),
            ('promo_code', "Discount Code"),
            ('buy_x_get_y', "Buy X Get Y"),
            ('next_order_coupons', "Next Order Coupons"),
        ],
        required=True,
        default='promotion',
    )
    date_from = fields.Date(
        string="Start Date",
        help="The start date is included in the validity period of this program",
    )
    date_to = fields.Date(
        string="End date",
        help="The end date is included in the validity period of this program",
    )
    limit_usage = fields.Boolean(string="Limit Usage")
    max_usage = fields.Integer()
    # Dictates when the points can be used:
    # current: if the order gives enough points on that order, the reward may directly be claimed, points lost otherwise
    # future: if the order gives enough points on that order, a coupon is generated for a next order
    # both: points are accumulated on the coupon to claim rewards, the reward may directly be claimed
    applies_on = fields.Selection(
        selection=[
            ('current', "Current order"),
            ('future', "Future orders"),
            ('both', "Current & Future orders"),
        ],
        compute='_compute_from_program_type',
        store=True,
        readonly=False,
        required=True,
        default='current',
    )
    trigger = fields.Selection(
        help="""
        Automatic: Customers will be eligible for a reward automatically in their cart.
        Use a code: Customers will be eligible for a reward if they enter a code.
        """,
        selection=[('auto', "Automatic"), ('with_code', "Use a code")],
        compute='_compute_from_program_type',
        store=True,
        readonly=False,
    )
    portal_visible = fields.Boolean(
        help="""
        Show in web portal, PoS customer ticket, eCommerce checkout, the number of points available
         and used by reward.
        """,
        default=False,
    )
    portal_point_name = fields.Char(
        translate=True,
        compute='_compute_portal_point_name',
        store=True,
        readonly=False,
        default='Points',
    )
    is_nominative = fields.Boolean(compute='_compute_is_nominative')
    is_payment_program = fields.Boolean(compute='_compute_is_payment_program')

    payment_program_discount_product_id = fields.Many2one(
        string="Discount Product",
        help="Product used in the sales order to apply the discount.",
        comodel_name='product.product',
        compute='_compute_payment_program_discount_product_id',
        readonly=True,
    )

    # Technical field used for a label
    available_on = fields.Boolean(
        string="Available On",
        help="Manage where your program should be available for use.",
        store=False,
    )

    _check_max_usage = models.Constraint(
        'CHECK (limit_usage = False OR max_usage > 0)',
        "Max usage must be strictly positive if a limit is used.",
    )

    @api.constrains('currency_id', 'pricelist_ids')
    def _check_pricelist_currency(self):
        if any(
            pricelist.currency_id != program.currency_id
            for program in self
            for pricelist in program.pricelist_ids
        ):
            raise UserError(_(
                "The loyalty program's currency must be the same as all it's pricelists ones."
            ))

    @api.constrains('date_from', 'date_to')
    def _check_date_from_date_to(self):
        if any(p.date_to and p.date_from and p.date_from > p.date_to for p in self):
            raise UserError(_(
                "The validity period's start date must be anterior or equal to its end date."
            ))

    @api.constrains('reward_ids')
    def _constrains_reward_ids(self):
        if self.env.context.get('loyalty_skip_reward_check'):
            return
        if any(not program.reward_ids for program in self):
            raise ValidationError(_("A program must have at least one reward."))

    def _compute_total_order_count(self):
        self.total_order_count = 0

    @api.depends('coupon_count', 'program_type')
    def _compute_coupon_count_display(self):
        program_items_name = self._program_items_name()
        for program in self:
            program.coupon_count_display = "%i %s" % (program.coupon_count or 0, program_items_name[program.program_type] or '')

    @api.depends('communication_plan_ids.mail_template_id')
    def _compute_mail_template_id(self):
        for program in self:
            program.mail_template_id = program.communication_plan_ids.mail_template_id[:1]

    def _inverse_mail_template_id(self):
        for program in self:
            if program.program_type not in ('gift_card', 'ewallet'):
                continue
            if not program.mail_template_id:
                program.communication_plan_ids = [(5, 0, 0)]
            elif not program.communication_plan_ids:
                program.communication_plan_ids = self.env['loyalty.mail'].create({
                    'program_id': program.id,
                    'trigger': 'create',
                    'mail_template_id': program.mail_template_id.id,
                })
            else:
                program.communication_plan_ids.write({
                    'trigger': 'create',
                    'mail_template_id': program.mail_template_id.id,
                })

    @api.depends('company_id')
    def _compute_currency_id(self):
        for program in self:
            program.currency_id = program.company_id.currency_id or program.currency_id

    @api.depends('coupon_ids')
    def _compute_coupon_count(self):
        read_group_data = self.env['loyalty.card']._read_group([('program_id', 'in', self.ids)], ['program_id'], ['__count'])
        count_per_program = {program.id: count for program, count in read_group_data}
        for program in self:
            program.coupon_count = count_per_program.get(program.id, 0)

    @api.depends('program_type', 'applies_on')
    def _compute_is_nominative(self):
        for program in self:
            program.is_nominative = program.applies_on == 'both' or\
                (program.program_type in ('ewallet', 'loyalty') and program.applies_on == 'future')

    @api.depends('program_type')
    def _compute_is_payment_program(self):
        for program in self:
            program.is_payment_program = program.program_type in ('gift_card', 'ewallet')

    @api.depends('reward_ids.discount_line_product_id')
    def _compute_payment_program_discount_product_id(self):
        for program in self:
            if program.is_payment_program:
                program.payment_program_discount_product_id = program.reward_ids[:1].discount_line_product_id
            else:
                program.payment_program_discount_product_id = False

    @api.model
    def _program_items_name(self):
        return {
            'coupons': _("Coupons"),
            'promotion': _("Promos"),
            'gift_card': _("Gift Cards"),
            'loyalty': _("Loyalty Cards"),
            'ewallet': _("eWallets"),
            'promo_code': _("Discounts"),
            'buy_x_get_y': _("Promos"),
            'next_order_coupons': _("Coupons"),
        }

    @api.model
    def _program_type_default_values(self):
        # All values to change when program_type changes
        # NOTE: any field used in `rule_ids`, `reward_ids` and `communication_plan_ids` MUST be present in the kanban view for it to work properly.
        first_sale_product = self.env['product.product'].search([('company_id', 'in', [False, self.env.company.id]), ('sale_ok', '=', True)], limit=1)
        return {
            'coupons': {
                'applies_on': 'current',
                'trigger': 'with_code',
                'portal_visible': False,
                'portal_point_name': _("Coupon point(s)"),
                'rule_ids': [(5, 0, 0)],
                'reward_ids': [(5, 0, 0), (0, 0, {
                    'required_points': 1,
                    'discount': 10,
                })],
                'communication_plan_ids': [(5, 0, 0), (0, 0, {
                    'trigger': 'create',
                    'mail_template_id': (self.env.ref('loyalty.mail_template_loyalty_card', raise_if_not_found=False) or self.env['mail.template']).id,
                })],
            },
            'promotion': {
                'applies_on': 'current',
                'trigger': 'auto',
                'portal_visible': False,
                'portal_point_name': _("Promo point(s)"),
                'rule_ids': [(5, 0, 0), (0, 0, {
                    'reward_point_amount': 1,
                    'reward_point_mode': 'order',
                    'minimum_amount': 50,
                    'minimum_qty': 0,
                })],
                'reward_ids': [(5, 0, 0), (0, 0, {
                    'required_points': 1,
                    'discount': 10,
                })],
                'communication_plan_ids': [(5, 0, 0)],
            },
            'gift_card': {
                'applies_on': 'future',
                'trigger': 'auto',
                'portal_visible': True,
                'portal_point_name': self.env.company.currency_id.symbol,
                'rule_ids': [(5, 0, 0), (0, 0, {
                    'reward_point_amount': 1,
                    'reward_point_mode': 'money',
                    'reward_point_split': True,
                    'product_ids': self.env.ref('loyalty.gift_card_product_50', raise_if_not_found=False),
                    'minimum_qty': 0,
                })],
                'reward_ids': [(5, 0, 0), (0, 0, {
                    'reward_type': 'discount',
                    'discount_mode': 'per_point',
                    'discount': 1,
                    'discount_applicability': 'order',
                    'required_points': 1,
                    'description': _("Gift Card"),
                })],
                'communication_plan_ids': [(5, 0, 0), (0, 0, {
                    'trigger': 'create',
                    'mail_template_id': (self.env.ref('loyalty.mail_template_gift_card', raise_if_not_found=False) or self.env['mail.template']).id,
                })],
            },
            'loyalty': {
                'applies_on': 'both',
                'trigger': 'auto',
                'portal_visible': True,
                'portal_point_name': _("Loyalty point(s)"),
                'rule_ids': [(5, 0, 0), (0, 0, {
                    'reward_point_mode': 'money',
                })],
                'reward_ids': [(5, 0, 0), (0, 0, {
                    'discount': 5,
                    'required_points': 200,
                })],
                'communication_plan_ids': [(5, 0, 0)],
            },
            'ewallet': {
                'trigger': 'auto',
                'applies_on': 'future',
                'portal_visible': True,
                'portal_point_name': self.env.company.currency_id.symbol,
                'rule_ids': [(5, 0, 0), (0, 0, {
                    'reward_point_amount': '1',
                    'reward_point_mode': 'money',
                    'reward_point_split': False,
                    'product_ids': self.env.ref('loyalty.ewallet_product_50', raise_if_not_found=False),
                })],
                'reward_ids': [(5, 0, 0), (0, 0, {
                    'reward_type': 'discount',
                    'discount_mode': 'per_point',
                    'discount': 1,
                    'discount_applicability': 'order',
                    'required_points': 1,
                    'description': _("eWallet"),
                })],
                'communication_plan_ids': [(5, 0, 0)],
            },
            'promo_code': {
                'applies_on': 'current',
                'trigger': 'with_code',
                'portal_visible': False,
                'portal_point_name': _("Discount point(s)"),
                'rule_ids': [(5, 0, 0), (0, 0, {
                    'mode': 'with_code',
                    'code': 'PROMO_CODE_' + str(uuid4())[:4], # We should try not to trigger any unicity constraint
                    'minimum_qty': 0,
                })],
                'reward_ids': [(5, 0, 0), (0, 0, {
                    'discount_applicability': 'specific',
                    'discount_product_ids': first_sale_product,
                    'discount_mode': 'percent',
                    'discount': 10,
                })],
                'communication_plan_ids': [(5, 0, 0)],
            },
            'buy_x_get_y': {
                'applies_on': 'current',
                'trigger': 'auto',
                'portal_visible': False,
                'portal_point_name': _("Credit(s)"),
                'rule_ids': [(5, 0, 0), (0, 0, {
                    'reward_point_mode': 'unit',
                    'product_ids': first_sale_product,
                    'minimum_qty': 2,
                })],
                'reward_ids': [(5, 0, 0), (0, 0, {
                    'reward_type': 'product',
                    'reward_product_id': first_sale_product.id,
                    'required_points': 2,
                })],
                'communication_plan_ids': [(5, 0, 0)],
            },
            'next_order_coupons': {
                'applies_on': 'future',
                'trigger': 'auto',
                'portal_visible': True,
                'portal_point_name': _("Coupon point(s)"),
                'rule_ids': [(5, 0, 0), (0, 0, {
                    'minimum_amount': 100,
                    'minimum_qty': 0,
                })],
                'reward_ids': [(5, 0, 0), (0, 0, {
                    'reward_type': 'discount',
                    'discount_mode': 'percent',
                    'discount': 15,
                    'discount_applicability': 'order',
                })],
                'communication_plan_ids': [(5, 0, 0), (0, 0, {
                    'trigger': 'create',
                    'mail_template_id': (
                        self.env.ref('loyalty.mail_template_loyalty_card', raise_if_not_found=False)
                        or self.env['mail.template']
                    ).id,
                })],
            },
        }

    @api.depends('program_type')
    def _compute_from_program_type(self):
        program_type_defaults = self._program_type_default_values()
        grouped_programs = defaultdict(lambda: self.env['loyalty.program'])
        for program in self:
            grouped_programs[program.program_type] |= program
        for program_type, programs in grouped_programs.items():
            if program_type in program_type_defaults:
                programs.write(program_type_defaults[program_type])

    @api.depends('currency_id', 'program_type')
    def _compute_portal_point_name(self):
        for program in self:
            if program.program_type not in ('ewallet', 'gift_card'):
                continue
            program.portal_point_name = program.currency_id.symbol or ''

    def _get_valid_products(self, products):
        '''
        Returns a dict containing the products that match per rule of the program
        '''
        rule_products = dict()
        for rule in self.rule_ids:
            domain = rule._get_valid_product_domain()
            if domain:
                rule_products[rule] = products.filtered_domain(domain)
            elif not domain and rule.program_type != 'gift_card':
                rule_products[rule] = products
            else:
                continue
        return rule_products

    def action_open_loyalty_cards(self):
        self.ensure_one()
        action = self.env['ir.actions.act_window']._for_xml_id('loyalty.loyalty_card_action')
        action['name'] = self._program_items_name()[self.program_type]
        action['display_name'] = action['name']
        action['context'] = {
            'program_type': self.program_type,
            'program_item_name': self._program_items_name()[self.program_type],
            'default_program_id': self.id,
            # For the wizard
            'default_mode': self.program_type == 'ewallet' and 'selected' or 'anonymous',
        }
        return action

    @api.ondelete(at_uninstall=False)
    def _unlink_except_active(self):
        if any(program.active for program in self):
            raise UserError(_("You can not delete a program in an active state"))

    def write(self, vals):
        # There is an issue when we change the program type, since we clear the rewards and create new ones.
        # The orm actually does it in this order upon writing, triggering the constraint before creating the new rewards.
        # However we can check that the result of reward_ids would actually be empty or not, and if not, skip the constraint.
        if 'reward_ids' in vals and self._fields['reward_ids'].convert_to_cache(vals['reward_ids'], self):
            self = self.with_context(loyalty_skip_reward_check=True)
            # We need add the program type to the context to avoid getting the default value
            # ('discount') for reward type when calling the `default_get` method of
            #`loyalty.reward`.
            if 'program_type' in vals:
                self = self.with_context(program_type=vals['program_type'])
                res = super().write(vals)
            else:
                for program in self:
                    program = program.with_context(program_type=program.program_type)
                    super(LoyaltyProgram, program).write(vals)
                res = True
        else:
            res = super().write(vals)

        # Propagate active state to children
        if 'active' in vals:
            for program in self.with_context(active_test=False):
                program.rule_ids.active = program.active
                program.reward_ids.active = program.active
                program.communication_plan_ids.active = program.active
                program.reward_ids.with_context(active_test=True).discount_line_product_id.active = program.active

        return res

    @api.model
    def get_program_templates(self):
        '''
        Returns the templates to be used for promotional programs.
        '''
        ctx_menu_type = self.env.context.get('menu_type')
        if ctx_menu_type == 'gift_ewallet':
            return {
                'gift_card': {
                    'title': _("Gift Card"),
                    'description': _("Sell Gift Cards, that allows to purchase products"),
                    'icon': 'gift_card',
                },
                'ewallet': {
                    'title': _("eWallet"),
                    'description': _("Fill in your eWallet, to pay future orders"),
                    'icon': 'ewallet',
                },
            }
        return {
            'promotion': {
                'title': _("Promotional Program"),
                'description': _("Automatic promo: 10% off on orders higher than $50"),
                'icon': 'promotional_program',
            },
            'promo_code': {
                'title': _("Promo Code"),
                'description': _("Get 10% off on some products, with a code"),
                'icon': 'promo_code',
            },
            'buy_x_get_y': {
                'title': _("Buy X Get Y"),
                'description': _("Buy 2 products and get a third one for free"),
                'icon': '2_plus_1',
            },
            'next_order_coupons': {
                'title': _("Next Order Coupon"),
                'description': _("Send a coupon after an order, valid for next purchase"),
                'icon': 'coupons',
            },
            'loyalty': {
                'title': _("Loyalty Card"),
                'description': _("Win points with each purchase, and claim gifts"),
                'icon': 'loyalty_cards',
            },
            'coupons': {
                'title': _("Coupon"),
                'description': _("Generate and share unique coupons with your customers"),
                'icon': 'coupons',
            },
            'fidelity': {
                'title': _("Fidelity Card"),
                'description': _("Buy 10 products to get 10$ off on the 11th one"),
                'icon': 'fidelity_cards',
            },
        }

    @api.model
    def create_from_template(self, template_id):
        '''
        Creates the program from the template id defined in `get_program_templates`.

        Returns an action leading to that new record.
        '''
        template_values = self._get_template_values()
        if template_id not in template_values:
            return False
        program = self.create(template_values[template_id])
        action = {}
        if self.env.context.get('menu_type') == 'gift_ewallet':
            action = self.env['ir.actions.act_window']._for_xml_id('loyalty.loyalty_program_gift_ewallet_action')
            action['views'] = [[False, 'form']]
        else:
            action = self.env['ir.actions.act_window']._for_xml_id('loyalty.loyalty_program_discount_loyalty_action')
            view_id = self.env.ref('loyalty.loyalty_program_view_form').id
            action['views'] = [[view_id, 'form']]
        action['view_mode'] = 'form'
        action['res_id'] = program.id
        return action

    @api.model
    def _get_template_values(self):
        '''
        Returns the values to create a program using the template keys defined above.
        '''
        program_type_defaults = self._program_type_default_values()
        # For programs that require a product get the first sellable.
        product = self.env['product.product'].search([('sale_ok', '=', True)], limit=1)
        return {
            'gift_card': {
                'name': _("Gift Card"),
                'program_type': 'gift_card',
                **program_type_defaults['gift_card']
            },
            'ewallet': {
                'name': _("eWallet"),
                'program_type': 'ewallet',
                **program_type_defaults['ewallet'],
            },
            'loyalty': {
                'name': _("Loyalty Cards"),
                'program_type': 'loyalty',
                **program_type_defaults['loyalty'],
            },
            'coupons': {
                'name': _("Coupons"),
                'program_type': 'coupons',
                **program_type_defaults['coupons'],
            },
            'promotion': {
                'name': _("Promotional Program"),
                'program_type': 'promotion',
                **program_type_defaults['promotion'],
            },
            'promo_code': {
                'name': _("Discount code"),
                'program_type': 'promo_code',
                **program_type_defaults['promo_code'],
            },
            'buy_x_get_y': {
                'name': _("2+1 Free"),
                'program_type': 'buy_x_get_y',
                **program_type_defaults['buy_x_get_y'],
            },
            'next_order_coupons': {
                'name': _("Next Order Coupons"),
                'program_type': 'next_order_coupons',
                **program_type_defaults['next_order_coupons'],
            },
            'fidelity': {
                'name': _("Fidelity Cards"),
                'program_type': 'loyalty',
                'applies_on': 'both',
                'trigger': 'auto',
                'rule_ids': [(0, 0, {
                    'reward_point_mode': 'unit',
                    'product_ids': product,
                })],
                'reward_ids': [(0, 0, {
                    'discount_mode': 'per_order',
                    'required_points': 11,
                    'discount_applicability': 'specific',
                    'discount_product_ids': product,
                    'discount': 10,
                })]
            },
        }

    @api.model_create_multi
    def create(self, vals_list):
        """
        trigger_product_ids will overwrite product ids defined in a loyalty rule in certain instances. Thus, it should
        be explicitly removed from an incoming vals dict unless, of course, it was actually a visible field.
        """
        for vals in vals_list:
            if 'trigger_product_ids' in vals and vals.get('program_type') not in ['gift_card', 'ewallet']:
                del vals['trigger_product_ids']

        return super().create(vals_list)
````

### FILE: `odoo_loyalty/upstream/addons/loyalty/models/loyalty_reward.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file20:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/loyalty/models/loyalty_reward.py; revision 99edb6dd82b7b560930c00b03b694ba700785370; Git blob 08ea7a9cf5c319931b32864413fecbbad2ec1427"
license: "LGPL-3.0-only"
sha256: "ec0aa37a9a5f12c01e4a860738afcc1b6ede113e4a3f4aa9d08966194936955c"
variables: []
secrets_allowed: false
```

````python
# Part of Odoo. See LICENSE file for full copyright and licensing details.

import ast
import json

from odoo import _, api, fields, models
from odoo.exceptions import ValidationError
from odoo.fields import Domain


class LoyaltyReward(models.Model):
    _name = 'loyalty.reward'
    _description = "Loyalty Reward"
    _rec_name = 'description'
    _order = 'required_points asc'

    @api.model
    def default_get(self, fields):
        # Try to copy the values of the program types default's
        result = super().default_get(fields)
        if 'program_type' in self.env.context:
            program_type = self.env.context['program_type']
            program_default_values = self.env['loyalty.program']._program_type_default_values()
            if program_type in program_default_values and\
                len(program_default_values[program_type]['reward_ids']) == 2 and\
                isinstance(program_default_values[program_type]['reward_ids'][1][2], dict):
                result.update({
                    k: v for k, v in program_default_values[program_type]['reward_ids'][1][2].items() if k in fields
                })
        return result

    def _get_discount_mode_select(self):
        # The value is provided in the loyalty program's view since we may not have a program_id yet
        #  and makes sure to display the currency related to the program instead of the company's.
        symbol = self.env.context.get('currency_symbol', self.env.company.currency_id.symbol)
        return [
            ('percent', "%"),
            ('per_order', symbol),
            ('per_point', _("%s per point", symbol)),
        ]

    @api.depends('program_id', 'description')
    def _compute_display_name(self):
        for reward in self:
            reward.display_name = f"{reward.program_id.name} - {reward.description}"

    active = fields.Boolean(default=True)
    program_id = fields.Many2one(comodel_name='loyalty.program', ondelete='cascade', required=True, index=True)
    program_type = fields.Selection(related='program_id.program_type')
    # Stored for security rules
    company_id = fields.Many2one(related='program_id.company_id', store=True)
    currency_id = fields.Many2one(related='program_id.currency_id')

    description = fields.Char(
        translate=True,
        compute='_compute_description',
        precompute=True,
        store=True,
        readonly=False,
        required=True,
    )

    reward_type = fields.Selection(
        selection=[
            ('product', "Free Product"),
            ('discount', "Discount"),
        ],
        required=True,
        default='discount',
    )
    user_has_debug = fields.Boolean(compute='_compute_user_has_debug')

    # Discount rewards
    discount = fields.Float(string="Discount", default=10)
    discount_mode = fields.Selection(
        selection=_get_discount_mode_select, required=True, default='percent'
    )
    discount_applicability = fields.Selection(
        selection=[
            ('order', "Order"),
            ('cheapest', "Cheapest Product"),
            ('specific', "Specific Products"),
        ],
        default='order',
    )
    discount_product_domain = fields.Char(default="[]")
    discount_product_ids = fields.Many2many(
        string="Discounted Products", comodel_name='product.product'
    )
    discount_product_category_id = fields.Many2one(
        string="Discounted Prod. Categories", comodel_name='product.category'
    )
    discount_product_tag_id = fields.Many2one(
        string="Discounted Prod. Tag", comodel_name='product.tag'
    )
    all_discount_product_ids = fields.Many2many(
        comodel_name='product.product', compute='_compute_all_discount_product_ids'
    )
    reward_product_domain = fields.Char(compute='_compute_reward_product_domain', store=False)
    discount_max_amount = fields.Monetary(
        string="Max Discount",
        help="This is the max amount this reward may discount, leave to 0 for no limit.",
    )
    discount_line_product_id = fields.Many2one(
        help="Product used in the sales order to apply the discount. Each reward has its own"
             " product for reporting purpose",
        comodel_name='product.product',
        ondelete='restrict',
        copy=False,
    )
    is_global_discount = fields.Boolean(compute='_compute_is_global_discount')

    # Product rewards
    reward_product_id = fields.Many2one(
        string="Product", comodel_name='product.product', domain=[('type', '!=', 'combo')]
    )
    reward_product_tag_id = fields.Many2one(string="Product Tag", comodel_name='product.tag')
    multi_product = fields.Boolean(compute='_compute_multi_product')
    reward_product_ids = fields.Many2many(
        string="Reward Products",
        help="These are the products that can be claimed with this rule.",
        comodel_name='product.product',
        compute='_compute_multi_product',
        search='_search_reward_product_ids',
    )
    reward_product_qty = fields.Integer(default=1)
    reward_product_uom_id = fields.Many2one(
        comodel_name='uom.uom', compute='_compute_reward_product_uom_id'
    )

    required_points = fields.Float(string="Points needed", default=1)
    point_name = fields.Char(related='program_id.portal_point_name', readonly=True)
    clear_wallet = fields.Boolean(default=False)

    _required_points_positive = models.Constraint(
        'CHECK (required_points > 0)',
        "The required points for a reward must be strictly positive.",
    )
    _product_qty_positive = models.Constraint(
        "CHECK (reward_type != 'product' OR reward_product_qty > 0)",
        "The reward product quantity must be strictly positive.",
    )
    _discount_positive = models.Constraint(
        "CHECK (reward_type != 'discount' OR discount > 0)",
        "The discount must be strictly positive.",
    )

    @api.depends('reward_product_id.product_tmpl_id.uom_id', 'reward_product_tag_id')
    def _compute_reward_product_uom_id(self):
        for reward in self:
            reward.reward_product_uom_id = reward.reward_product_ids.product_tmpl_id.uom_id[:1]

    def _find_all_category_children(self, category_id, child_ids):
        if len(category_id.child_id) > 0:
            for child_id in category_id.child_id:
                child_ids.append(child_id.id)
                self._find_all_category_children(child_id, child_ids)
        return child_ids

    def _get_discount_product_domain(self):
        self.ensure_one()
        constrains = []
        if self.discount_product_ids:
            constrains.append([('id', 'in', self.discount_product_ids.ids)])
        if self.discount_product_category_id:
            product_category_ids = self._find_all_category_children(self.discount_product_category_id, [])
            product_category_ids.append(self.discount_product_category_id.id)
            constrains.append([('categ_id', 'in', product_category_ids)])
        if self.discount_product_tag_id:
            constrains.append([('all_product_tag_ids', 'in', self.discount_product_tag_id.id)])
        domain = Domain.OR(constrains) if constrains else Domain.TRUE
        if self.discount_product_domain and self.discount_product_domain != '[]':
            domain &= Domain(ast.literal_eval(self.discount_product_domain))
        return domain

    @api.model
    def _get_active_products_domain(self):
        return [
            '|',
                ('reward_type', '!=', 'product'),
                '&',
                    ('reward_type', '=', 'product'),
                    '|',
                        '&',
                            ('reward_product_tag_id', '=', False),
                            ('reward_product_id.active', '=', True),
                        '&',
                            ('reward_product_tag_id', '!=', False),
                            ('reward_product_ids.active', '=', True)
        ]

    @api.depends('discount_product_domain')
    def _compute_reward_product_domain(self):
        compute_all_discount_product = self.env['ir.config_parameter'].sudo().get_param('loyalty.compute_all_discount_product_ids', 'enabled')
        for reward in self:
            if compute_all_discount_product == 'enabled':
                reward.reward_product_domain = "null"
            else:
                reward.reward_product_domain = json.dumps(list(reward._get_discount_product_domain()))

    @api.depends('discount_product_ids', 'discount_product_category_id', 'discount_product_tag_id', 'discount_product_domain')
    def _compute_all_discount_product_ids(self):
        compute_all_discount_product = self.env['ir.config_parameter'].sudo().get_param('loyalty.compute_all_discount_product_ids', 'enabled')
        for reward in self:
            if compute_all_discount_product == 'enabled':
                reward.all_discount_product_ids = self.env['product.product'].search(reward._get_discount_product_domain())
            else:
                reward.all_discount_product_ids = self.env['product.product']

    @api.depends('reward_product_id', 'reward_product_tag_id', 'reward_type')
    def _compute_multi_product(self):
        for reward in self:
            products = reward.reward_product_id + reward.reward_product_tag_id.product_ids.filtered(
                lambda product: product.type != 'combo'
            )
            reward.multi_product = reward.reward_type == 'product' and len(products) > 1
            reward.reward_product_ids = reward.reward_type == 'product' and products or self.env['product.product']

    def _search_reward_product_ids(self, operator, value):
        if operator != 'in':
            return NotImplemented
        return [
            '&', ('reward_type', '=', 'product'),
            '|', ('reward_product_id', operator, value),
            ('reward_product_tag_id.product_ids', operator, value)
        ]

    @api.depends('reward_type', 'reward_product_id', 'discount_mode', 'reward_product_tag_id',
                 'discount', 'currency_id', 'discount_applicability', 'all_discount_product_ids')
    def _compute_description(self):
        for reward in self:
            reward_string = ""
            if reward.program_type == 'gift_card':
                reward_string = _("Gift Card")
            elif reward.program_type == 'ewallet':
                reward_string = _("eWallet")
            elif reward.reward_type == 'product':
                products = reward.reward_product_ids
                if len(products) == 0:
                    reward_string = _("Free Product")
                elif len(products) == 1:
                    reward_string = _("Free Product - %s", reward.reward_product_id.with_context(display_default_code=False).display_name)
                else:
                    reward_string = _("Free Product - [%s]", ', '.join(products.with_context(display_default_code=False).mapped('display_name')))
            elif reward.reward_type == 'discount':
                format_string = "%(amount)g %(symbol)s"
                if reward.currency_id.position == 'before':
                    format_string = "%(symbol)s %(amount)g"
                formatted_amount = format_string % {'amount': reward.discount, 'symbol': reward.currency_id.symbol}
                if reward.discount_mode == 'percent':
                    reward_string = _("%g%% on ", reward.discount)
                elif reward.discount_mode == 'per_point':
                    reward_string = _("%s per point on ", formatted_amount)
                elif reward.discount_mode == 'per_order':
                    reward_string = _("%s on ", formatted_amount)
                if reward.discount_applicability == 'order':
                    reward_string += _("your order")
                elif reward.discount_applicability == 'cheapest':
                    reward_string += _("the cheapest product")
                elif reward.discount_applicability == 'specific':
                    product_available = self.env['product.product'].search(reward._get_discount_product_domain(), limit=2)
                    if len(product_available) == 1:
                        reward_string += product_available.with_context(display_default_code=False).display_name
                    else:
                        reward_string += _("specific products")
                if reward.discount_max_amount:
                    format_string = "%(amount)g %(symbol)s"
                    if reward.currency_id.position == 'before':
                        format_string = "%(symbol)s %(amount)g"
                    formatted_amount = format_string % {'amount': reward.discount_max_amount, 'symbol': reward.currency_id.symbol}
                    reward_string += _(" (Max %s)", formatted_amount)
            reward.description = reward_string

    @api.depends('reward_type', 'discount_applicability', 'discount_mode')
    def _compute_is_global_discount(self):
        for reward in self:
            reward.is_global_discount = (
                reward.reward_type == 'discount'
                and reward.discount_applicability == 'order'
                and reward.discount_mode in ['per_order', 'percent']
            )

    @api.depends_context('uid')
    @api.depends('reward_type')
    def _compute_user_has_debug(self):
        self.user_has_debug = self.env.user.has_group('base.group_no_one')

    @api.constrains('reward_product_id')
    def _check_reward_product_id_no_combo(self):
        if any(reward.reward_product_id.type == 'combo' for reward in self):
            raise ValidationError(_("A reward product can't be of type \"combo\"."))

    def _create_missing_discount_line_products(self):
        # Make sure we create the product that will be used for our discounts
        rewards = self.filtered(lambda r: not r.discount_line_product_id)
        products = self.env['product.product'].create(rewards._get_discount_product_values())
        for reward, product in zip(rewards, products):
            reward.discount_line_product_id = product

    @api.model_create_multi
    def create(self, vals_list):
        res = super().create(vals_list)
        res._create_missing_discount_line_products()
        return res

    def write(self, vals):
        res = super().write(vals)
        if 'description' in vals:
            self._create_missing_discount_line_products()
            # Keep the name of our discount product up to date
            for reward in self:
                reward.discount_line_product_id.write({'name': reward.description})
        if 'active' in vals:
            if vals['active']:
                self.discount_line_product_id.action_unarchive()
            else:
                self.discount_line_product_id.action_archive()
        return res

    def update_field_translations(self, field_name, translations, source_lang=''):
        res = super().update_field_translations(field_name, translations, source_lang=source_lang)
        if field_name == 'description' and self.discount_line_product_id:
            self.discount_line_product_id.update_field_translations('name', translations, source_lang=source_lang)
        return res

    def unlink(self):
        programs = self.program_id
        res = super().unlink()
        # Not guaranteed to trigger the constraint
        programs._constrains_reward_ids()
        return res

    def _get_discount_product_values(self):
        return [{
            'name': reward.description,
            'type': 'service',
            'sale_ok': False,
            'purchase_ok': False,
            'lst_price': 0,
        } for reward in self]
````

### FILE: `odoo_loyalty/upstream/addons/loyalty/models/loyalty_rule.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file21:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/loyalty/models/loyalty_rule.py; revision 99edb6dd82b7b560930c00b03b694ba700785370; Git blob 6ee64ff2055640f7b9a012338b4e93faa1db8570"
license: "LGPL-3.0-only"
sha256: "955eafbc11279154ba30b4986fa2198c683d4c4519253f5ee0fd84daa39f2a66"
variables: []
secrets_allowed: false
```

````python
# Part of Odoo. See LICENSE file for full copyright and licensing details.

import ast

from odoo import _, api, fields, models
from odoo.exceptions import ValidationError
from odoo.fields import Domain


class LoyaltyRule(models.Model):
    _name = 'loyalty.rule'
    _description = "Loyalty Rule"

    @api.model
    def default_get(self, fields):
        # Try to copy the values of the program types default's
        result = super().default_get(fields)
        if 'program_type' in self.env.context:
            program_type = self.env.context['program_type']
            program_default_values = self.env['loyalty.program']._program_type_default_values()
            if program_type in program_default_values and\
                len(program_default_values[program_type]['rule_ids']) == 2 and\
                isinstance(program_default_values[program_type]['rule_ids'][1][2], dict):
                result.update({
                    k: v for k, v in program_default_values[program_type]['rule_ids'][1][2].items() if k in fields
                })
        return result

    def _get_reward_point_mode_selection(self):
        # The value is provided in the loyalty program's view since we may not have a program_id yet
        #  and makes sure to display the currency related to the program instead of the company's.
        symbol = self.env.context.get('currency_symbol', self.env.company.currency_id.symbol)
        return [
            ('order', _("per order")),
            ('money', _("per %s spent", symbol)),
            ('unit', _("per unit paid")),
        ]

    active = fields.Boolean(default=True)
    program_id = fields.Many2one(comodel_name='loyalty.program', ondelete='cascade', required=True, index=True)
    program_type = fields.Selection(related='program_id.program_type')
    # Stored for security rules
    company_id = fields.Many2one(related='program_id.company_id', store=True)
    currency_id = fields.Many2one(related='program_id.currency_id')

    # Only for dev mode
    user_has_debug = fields.Boolean(compute='_compute_user_has_debug')
    product_domain = fields.Char(default="[]")

    product_ids = fields.Many2many(string="Products", comodel_name='product.product')
    product_category_id = fields.Many2one(string="Categories", comodel_name='product.category')
    product_tag_id = fields.Many2one(string="Product Tag", comodel_name='product.tag')

    reward_point_amount = fields.Float(string="Reward", default=1)
    # Only used for program_id.applies_on == 'future'
    reward_point_split = fields.Boolean(
        string="Split per unit",
        help="Whether to separate reward coupons per matched unit, only applies to 'future' programs and trigger mode per money spent or unit paid...",
        default=False,
    )
    reward_point_name = fields.Char(related='program_id.portal_point_name', readonly=True)
    reward_point_mode = fields.Selection(
        selection=_get_reward_point_mode_selection, required=True, default='order'
    )

    minimum_qty = fields.Integer(string="Minimum Quantity", default=1)
    minimum_amount = fields.Monetary(string="Minimum Purchase")
    minimum_amount_tax_mode = fields.Selection(
        selection=[
            ('incl', "tax included"),
            ('excl', "tax excluded"),
        ],
        required=True,
        default='incl',
    )

    mode = fields.Selection(
        string="Application",
        selection=[
            ('auto', "Automatic"),
            ('with_code', "With a promotion code"),
        ],
        compute='_compute_mode',
        store=True,
        readonly=False,
    )
    code = fields.Char(string="Discount code", compute='_compute_code', store=True, readonly=False)

    _reward_point_amount_positive = models.Constraint(
        'CHECK (reward_point_amount > 0)',
        "Rule points reward must be strictly positive.",
    )

    @api.constrains('reward_point_split')
    def _constraint_trigger_multi(self):
        # Prevent setting trigger multi in case of nominative programs, it does not make sense to allow this
        for rule in self:
            if rule.reward_point_split and (rule.program_id.applies_on == 'both' or rule.program_id.program_type == 'ewallet'):
                raise ValidationError(_("Split per unit is not allowed for Loyalty and eWallet programs."))

    @api.constrains('code', 'active')
    def _constrains_code(self):
        mapped_codes = self.filtered(lambda r: r.code and r.active).mapped('code')
        # Program code must be unique
        if len(mapped_codes) != len(set(mapped_codes)) or\
            self.env['loyalty.rule'].search_count([
                ('mode', '=', 'with_code'),
                ('code', 'in', mapped_codes),
                ('id', 'not in', self.ids),
                ('active', '=', True),
            ]):
            raise ValidationError(_("The promo code must be unique."))
        # Prevent coupons and programs from sharing a code
        if self.env['loyalty.card'].search_count([
            ('code', 'in', mapped_codes), ('active', '=', True)
        ]):
            raise ValidationError(_("A coupon with the same code was found."))

    @api.depends('mode')
    def _compute_code(self):
        # Reset code when mode is set to auto
        for rule in self:
            if rule.mode == 'auto':
                rule.code = False

    @api.depends('code')
    def _compute_mode(self):
        for rule in self:
            if rule.code:
                rule.mode = 'with_code'
            else:
                rule.mode = 'auto'

    @api.depends_context('uid')
    @api.depends('mode')
    def _compute_user_has_debug(self):
        self.user_has_debug = self.env.user.has_group('base.group_no_one')

    def _get_valid_product_domain(self):
        self.ensure_one()
        constrains = []
        if self.product_ids:
            constrains.append([('id', 'in', self.product_ids.ids)])
        if self.product_category_id:
            constrains.append([('categ_id', 'child_of', self.product_category_id.id)])
        if self.product_tag_id:
            constrains.append([('all_product_tag_ids', 'in', self.product_tag_id.id)])
        domain = Domain.OR(constrains) if constrains else Domain.TRUE
        if self.product_domain and self.product_domain != '[]':
            domain &= Domain(ast.literal_eval(self.product_domain))
        return domain

    def _get_valid_products(self):
        self.ensure_one()
        return self.env['product.product'].search(self._get_valid_product_domain())

    def _compute_amount(self, currency_to):
        self.ensure_one()
        return self.currency_id._convert(
            self.minimum_amount,
            currency_to,
            self.company_id or self.env.company,
            fields.Date.today()
        )
````

### FILE: `odoo_loyalty/upstream/addons/loyalty/tests/test_loyalty.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file22:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/loyalty/tests/test_loyalty.py; revision 99edb6dd82b7b560930c00b03b694ba700785370; Git blob b7acc32a1c1b8e9667a9c2d130b42475cb0e0248"
license: "LGPL-3.0-only"
sha256: "39c7660a20b402f80fafab7e6534c9bbcd2e52abd847088c0263004fe33092c3"
variables: []
secrets_allowed: false
```

````python
# Part of Odoo. See LICENSE file for full copyright and licensing details.

from unittest.mock import patch

from psycopg2 import IntegrityError

from odoo.exceptions import UserError, ValidationError
from odoo.fields import Command
from odoo.tests import Form, TransactionCase, tagged
from odoo.tools import mute_logger


@tagged('post_install', '-at_install')
class TestLoyalty(TransactionCase):

    @classmethod
    def setUpClass(cls):
        super().setUpClass()

        cls.program = cls.env['loyalty.program'].create({
            'name': 'Test Program',
            'reward_ids': [(0, 0, {})],
        })
        cls.product = cls.env['product.product'].with_context(default_taxes_id=False).create({
            'name': "Test Product",
            'type': 'consu',
            'list_price': 20.0,
        })

    def create_program_with_code(self, code):
        return self.env['loyalty.program'].create({
            'name': "Discount delivery",
            'program_type': 'promo_code',
            'rule_ids': [Command.create({
                'code': code,
                'minimum_amount': 0,
            })],
        })

    def test_loyalty_program_default_values(self):
        # Test that the default values are correctly set when creating a new program
        program = self.env['loyalty.program'].create({'name': "Test"})
        self._check_promotion_default_values(program)

    def _check_promotion_default_values(self, program):
        self.assertEqual(program.program_type, 'promotion')
        self.assertEqual(program.trigger, 'auto')
        self.assertEqual(program.portal_visible, False)
        self.assertTrue(program.rule_ids)
        self.assertTrue(len(program.rule_ids) == 1)
        self.assertEqual(program.rule_ids.reward_point_amount, 1)
        self.assertEqual(program.rule_ids.reward_point_mode, 'order')
        self.assertEqual(program.rule_ids.minimum_amount, 50)
        self.assertEqual(program.rule_ids.minimum_qty, 0)
        self.assertTrue(program.reward_ids)
        self.assertTrue(len(program.reward_ids) == 1)
        self.assertEqual(program.reward_ids.required_points, 1)
        self.assertEqual(program.reward_ids.discount, 10)
        self.assertFalse(program.communication_plan_ids)

    def test_loyalty_program_default_values_in_form(self):
        # Test that the default values are correctly set when creating a new program in a form
        with Form(self.env['loyalty.program']) as program_form:
            program_form.name = 'Test'
            program = program_form.save()
        self._check_promotion_default_values(program)

    def test_discount_product_unlink(self):
        # Test that we can not unlink discount line product id
        with mute_logger('odoo.sql_db'), self.assertRaises(IntegrityError):
            self.program.reward_ids.discount_line_product_id.unlink()

    def test_loyalty_mail(self):
        # Test basic loyalty_mail functionalities
        loyalty_card_model_id = self.env.ref('loyalty.model_loyalty_card')
        create_tmpl, fifty_tmpl, hundred_tmpl = self.env['mail.template'].create([
            {
                'name': 'CREATE',
                'model_id': loyalty_card_model_id.id,
            },
            {
                'name': '50 points',
                'model_id': loyalty_card_model_id.id,
            },
            {
                'name': '100 points',
                'model_id': loyalty_card_model_id.id,
            },
        ])
        self.program.write({'communication_plan_ids': [
            (0, 0, {
                'program_id': self.program.id,
                'trigger': 'create',
                'mail_template_id': create_tmpl.id,
            }),
            (0, 0, {
                'program_id': self.program.id,
                'trigger': 'points_reach',
                'points': 50,
                'mail_template_id': fifty_tmpl.id,
            }),
            (0, 0, {
                'program_id': self.program.id,
                'trigger': 'points_reach',
                'points': 100,
                'mail_template_id': hundred_tmpl.id,
            }),
        ]})

        sent_mails = self.env['mail.template']

        def mock_send_mail(self, *args, **kwargs):
            nonlocal sent_mails
            sent_mails |= self

        partner = self.env['res.partner'].create({'name': 'Test Partner'})
        with patch('odoo.addons.mail.models.mail_template.MailTemplate.send_mail', new=mock_send_mail):
            # Send mail at creation
            coupon = self.env['loyalty.card'].create({
                'program_id': self.program.id,
                'partner_id': partner.id,
                'points': 0,
            })
            self.assertEqual(sent_mails, create_tmpl)
            sent_mails = self.env['mail.template']
            # 50 points mail
            coupon.points = 50
            self.assertEqual(sent_mails, fifty_tmpl)
            sent_mails = self.env['mail.template']
            # Check that it does not get sent again
            coupon.points = 99
            self.assertFalse(sent_mails)
            # 100 points mail
            coupon.points = 100
            self.assertEqual(sent_mails, hundred_tmpl)
            sent_mails = self.env['mail.template']
            # Reset and go straight to 100 points
            coupon.points = 0
            self.assertFalse(sent_mails)
            coupon.points = 100
            self.assertEqual(sent_mails, hundred_tmpl)

    def test_loyalty_program_preserve_reward_upon_writing(self):
        self.program.program_type = 'buy_x_get_y'
        # recompute of rewards
        self.program.flush_recordset(['reward_ids'])

        self.program.write({
            'reward_ids': [
                Command.create({
                    'description': 'Test Product',
                }),
            ],
        })
        self.assertTrue(all(r.reward_type == 'product' for r in self.program.reward_ids))

    def test_loyalty_program_preserve_reward_with_always_edit(self):
        with Form(self.env['loyalty.program']) as program_form:
            program_form.name = 'Test'
            program_form.program_type = 'buy_x_get_y'
            program_form.reward_ids.remove(0)
            with program_form.reward_ids.new() as new_reward:
                new_reward.reward_product_qty = 2
            program = program_form.save()
            self.assertEqual(program.reward_ids.reward_type, 'product')
            self.assertEqual(program.reward_ids.reward_product_qty, 2)

    def test_archiving_unarchiving(self):
        self.program.write({
            'reward_ids': [
                Command.create({
                    'description': 'Test Product',
                }),
            ],
        })
        before_archived_reward_ids = self.program.reward_ids
        self.program.action_archive()
        self.program.action_unarchive()
        after_archived_reward_ids = self.program.reward_ids
        self.assertEqual(before_archived_reward_ids, after_archived_reward_ids)

    def test_prevent_archive_pricelist_linked_to_program(self):
        self.program.pricelist_ids = demo_pricelist = self.env['product.pricelist'].create({
            'name': "Demo"
        })
        with self.assertRaises(UserError):
            demo_pricelist.action_archive()
        self.program.action_archive()
        demo_pricelist.action_archive()

    def test_prevent_archiving_product_linked_to_active_loyalty_reward(self):
        self.program.program_type = 'promotion'
        self.program.flush_recordset()
        reward = self.env['loyalty.reward'].create({
            'program_id': self.program.id,
            'discount_line_product_id': self.product.id,
        })
        self.program.write({
            'reward_ids': [Command.link(reward.id)],
        })
        with self.assertRaises(ValidationError):
            self.product.action_archive()
        self.program.action_archive()
        self.product.action_archive()

    def test_prevent_archiving_product_used_for_discount_reward(self):
        """
        Ensure products cannot be archived while they have a specific program active.
        """
        self.program.write({
            'name': f"50% Discount on {self.product.name}",
            'program_type': 'promotion',
            'reward_ids': [Command.create({
                'discount': 50.0,
                'discount_applicability': 'specific',
                'discount_product_ids': self.product.ids,
            })],
        })
        with self.assertRaises(ValidationError):
            self.product.action_archive()
        self.program.action_archive()
        self.product.action_archive()
        self.assertFalse(self.product.active)

    def test_prevent_archiving_product_when_archiving_program(self):
        """
        Test prevent archiving a product when archiving a "Buy X Get Y" program.
        We just have to archive the free product that has been created while creating
        the program itself not the product we already had before.
        """
        loyalty_program = self.env['loyalty.program'].create({
            'name': 'Test Program',
            'program_type': 'buy_x_get_y',
            'reward_ids': [
                Command.create({
                    'description': 'Test Product',
                    'reward_product_id': self.product.id,
                    'reward_type': 'product'
                }),
            ],
        })
        loyalty_program.action_archive()
        # Make sure that the main product didn't get archived
        self.assertTrue(self.product.active)

    def test_merge_loyalty_cards(self):
        """Test merging nominative loyalty cards from source partners to a destination partner
        when partners are merged.
        """
        program = self.env['loyalty.program'].create({
            'name': 'Test Program',
            'is_nominative': True,
            'applies_on': 'both',
        })

        partner_1, partner_2, dest_partner = self.env['res.partner'].create([
            {'name': 'Source Partner 1'},
            {'name': 'Source Partner 2'},
            {'name': 'Destination Partner'},
        ])
        self.env['loyalty.card'].create([
            {
                'partner_id': partner_1.id,
                'program_id': program.id,
                'points': 10
            }, {
                'partner_id': partner_2.id,
                'program_id': program.id,
                'points': 20
            }, {
                'partner_id': dest_partner.id,
                'program_id': program.id,
                'points': 30
            }
        ])

        self.env['base.partner.merge.automatic.wizard']._merge(
            [partner_1.id, partner_2.id, dest_partner.id], dest_partner
        )

        dest_partner_loyalty_cards = self.env['loyalty.card'].search([
            ('partner_id', '=', dest_partner.id),
            ('program_id', '=', program.id),
        ])

        self.assertEqual(len(dest_partner_loyalty_cards), 1)
        self.assertEqual(dest_partner_loyalty_cards.points, 60)
        self.assertFalse(self.env['loyalty.card'].search([
            ('partner_id', 'in', [partner_1.id, partner_2.id]),
        ]))

    def test_card_description_on_tag_change(self):
        product_tag = self.env['product.tag'].create({'name': 'Multiple Products'})
        product1 = self.product
        product1.product_tag_ids = product_tag
        self.env['product.product'].create({
            'name': 'Test Product 2',
            'list_price': 30.0,
            'product_tag_ids': product_tag,
        })
        reward = self.env['loyalty.reward'].create({
            'program_id': self.program.id,
            'reward_type': 'product',
            'reward_product_id': product1.id,
        })
        reward_description_single_product = reward.description
        reward.reward_product_tag_id = product_tag
        reward_description_product_tag = reward.description
        self.assertNotEqual(
            reward_description_single_product,
            reward_description_product_tag,
            "Reward description should be changed after adding a tag"
        )
        self.assertEqual(
            reward_description_product_tag,
            "Free Product - [Test Product, Test Product 2]",
            "Reward description for reward with tag should be 'Free Product - [Test Product, Test Product 2]'"
        )

    def test_prevent_unarchive_when_conflicting_active_program_exists(self):
        """Unarchiving a program should fail if another active program already has the same rule
           code."""
        program = self.create_program_with_code("FREE")
        program.action_archive()
        # create another active program with the same rule code
        self.create_program_with_code("FREE")
        # attempt to unarchive the first program
        with self.assertRaises(ValidationError):
            program.action_unarchive()

    def test_prevent_unarchive_when_batch_contains_duplicate_codes(self):
        """Unarchiving multiple programs at once should fail if they share the same rule code."""
        program1 = self.create_program_with_code("FREE")
        program1.action_archive()
        # create another program with the same rule code and archive it
        program2 = self.create_program_with_code("FREE")
        program2.action_archive()
        # attempt to unarchive both programs together
        with self.assertRaises(ValidationError):
            (program1 + program2).action_unarchive()

    def test_discount_description_translation(self):
        """A discount product's name field should automatically update for all languages for which changes
        are made on the reward's description"""
        self.env['res.lang']._activate_lang('fr_FR')
        program = self.env['loyalty.program'].create({
            'name': 'Test Program',
            'reward_ids': [(0, 0, {})],
        })
        reward = self.env['loyalty.reward'].with_context(lang='en_US').create({
            'program_id': program.id,
            'reward_type': 'discount',
            'description': 'My Discount'
        })
        product = reward.discount_line_product_id
        translations = {'en_US': 'Test Discount EN', 'fr_FR': 'Test Discount FR'}
        reward.update_field_translations('description', translations)
        product.invalidate_recordset(['name'])
        self.assertEqual(product.with_context(lang='en_US').name, 'Test Discount EN')
        self.assertEqual(product.with_context(lang='fr_FR').name, 'Test Discount FR')
````

### FILE: `odoo_loyalty/upstream/addons/sale_loyalty/__manifest__.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file23:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/__manifest__.py; revision 99edb6dd82b7b560930c00b03b694ba700785370; Git blob bb0716a9637eca91dca26385f9bac21848479726"
license: "LGPL-3.0-only"
sha256: "284a427b610067bac9c82ae6e3cac37cb1ad51c3dee5617fc5b93d62b13ee083"
variables: []
secrets_allowed: false
```

````python
# Part of Odoo. See LICENSE file for full copyright and licensing details.

{
    'name': 'Sale Loyalty',
    'summary': 'Use discounts and loyalty programs in sales orders',
    'description': 'Integrate discount and loyalty programs mechanisms in sales orders.',
    'category': 'Sales/Sales',
    'version': '1.0',
    'depends': ['sale', 'loyalty'],
    'auto_install': True,
    'data': [
        'security/ir.model.access.csv',

        'data/sale_loyalty_data.xml',

        'wizard/sale_loyalty_coupon_wizard_views.xml',
        'wizard/sale_loyalty_reward_wizard_views.xml',

        'views/loyalty_card_views.xml',
        'views/loyalty_program_views.xml',
        'views/sale_order_views.xml',
        'views/sale_portal_templates.xml',
        'views/res_partner_views.xml',
        'views/sale_loyalty_menus.xml',
    ],
    'assets': {
        'web.assets_backend': [
            'sale_loyalty/static/src/**/*',
        ],
    },
    'uninstall_hook': 'uninstall_hook',
    'author': 'Odoo S.A.',
    'license': 'LGPL-3',
}
````

### FILE: `odoo_loyalty/upstream/addons/sale_loyalty/models/loyalty_card.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file24:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/models/loyalty_card.py; revision 99edb6dd82b7b560930c00b03b694ba700785370; Git blob 9a52e86399df0af0cf44ea79046fd7c59558498c"
license: "LGPL-3.0-only"
sha256: "67a9501daf1b48bffbf10e34d11b7503df920a0048545affd38e86f08f2be712"
variables: []
secrets_allowed: false
```

````python
# Part of Odoo. See LICENSE file for full copyright and licensing details.

from odoo import fields, models


class LoyaltyCard(models.Model):
    _inherit = 'loyalty.card'

    order_id = fields.Many2one(
        string="Order Reference",
        help="The sales order from which coupon is generated",
        comodel_name='sale.order',
        readonly=True,
    )
    order_id_partner_id = fields.Many2one(
        string="Sale Order Customer", comodel_name='res.partner', related='order_id.partner_id'
    )

    def _get_default_template(self):
        default_template = super()._get_default_template()
        if not default_template:
            default_template = self.env.ref('loyalty.mail_template_loyalty_card', raise_if_not_found=False)
        return default_template

    def _mail_get_partner_fields(self, introspect_fields=False):
        return super()._mail_get_partner_fields(introspect_fields=introspect_fields) + ['order_id_partner_id']

    def _get_mail_author(self):
        # Default author is the order's salesperson if available, else the order's company.
        if not self.order_id or self.order_id.sudo().company_id not in self.env.companies:
            return super()._get_mail_author()
        self.ensure_one()
        return (self.order_id.user_id or self.order_id.company_id).partner_id

    def _get_signature(self):
        return self.order_id.user_id.signature or super()._get_signature()

    def _compute_use_count(self):
        super()._compute_use_count()
        read_group_res = self.env['sale.order.line']._read_group(
            [('coupon_id', 'in', self.ids)], ['coupon_id'], ['__count'])
        count_per_coupon = {coupon.id: count for coupon, count in read_group_res}
        for card in self:
            card.use_count += count_per_coupon.get(card.id, 0)

    def _has_source_order(self):
        return super()._has_source_order() or bool(self.order_id)

    def action_archive(self):
        self.env['sale.order.coupon.points'].search([
            ('coupon_id', 'in', self.ids),
            ('order_id.state', '=', 'draft'),
        ]).unlink()
        return super().action_archive()
````

### FILE: `odoo_loyalty/upstream/addons/sale_loyalty/models/loyalty_program.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file25:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/models/loyalty_program.py; revision 99edb6dd82b7b560930c00b03b694ba700785370; Git blob d8d8e3c40fed4ea46e592d79dbadda8c62bae27a"
license: "LGPL-3.0-only"
sha256: "2f1bbcfb4419d8bd6149eac02e963670c342470dd2007eaa02e31b095704e658"
variables: []
secrets_allowed: false
```

````python
# Part of Odoo. See LICENSE file for full copyright and licensing details.

from odoo import fields, models


class LoyaltyProgram(models.Model):
    _inherit = 'loyalty.program'

    order_count = fields.Integer(compute='_compute_order_count')
    sale_ok = fields.Boolean(string="Sales", default=True)

    def _compute_order_count(self):
        # An order should count only once PER program but may appear in multiple programs
        read_group_res = self.env['sale.order.line']._read_group(
            [('reward_id', 'in', self.reward_ids.ids)], ['order_id'], ['reward_id:array_agg'])
        for program in self:
            program_reward_ids = program.reward_ids.ids
            program.order_count = sum(
                any(id_ in reward_ids for id_ in program_reward_ids)
                for __, reward_ids in read_group_res
            )

    def _compute_total_order_count(self):
        super()._compute_total_order_count()
        for program in self:
            program.total_order_count += program.order_count
````

### FILE: `odoo_loyalty/upstream/addons/sale_loyalty/models/sale_order.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file26:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/models/sale_order.py; revision 99edb6dd82b7b560930c00b03b694ba700785370; Git blob 876798076f76b3d4b6164ef02eb098d7ac7dd478"
license: "LGPL-3.0-only"
sha256: "6de1aeed8f1fc80268025f13bff2a56b57211b6a82ce9df950a8b08a623e3a06"
variables: []
secrets_allowed: false
```

````python
# Part of Odoo. See LICENSE file for full copyright and licensing details.

import itertools
import random
from collections import defaultdict
from functools import partial

from pytz import timezone

from odoo import _, api, fields, models
from odoo.exceptions import UserError, ValidationError
from odoo.fields import Command, Domain
from odoo.tools import float_round, lazy, str2bool


def _generate_random_reward_code():
    return str(random.getrandbits(32))


class SaleOrder(models.Model):
    _inherit = 'sale.order'

    # Contains how much points should be given to a coupon upon validating the order
    applied_coupon_ids = fields.Many2many(
        string="Manually Applied Coupons", comodel_name='loyalty.card', copy=False
    )
    code_enabled_rule_ids = fields.Many2many(
        string="Manually Triggered Rules", comodel_name='loyalty.rule', copy=False
    )
    coupon_point_ids = fields.One2many(
        comodel_name='sale.order.coupon.points', inverse_name='order_id', copy=False
    )
    reward_amount = fields.Float(compute='_compute_reward_total')

    # Display Fields
    gift_card_count = fields.Integer(compute='_compute_gift_card_count')
    loyalty_data = fields.Json(compute='_compute_loyalty_data')

    @api.depends('order_line')
    def _compute_reward_total(self):
        for order in self:
            reward_amount = 0
            for line in order.order_line:
                if not line.reward_id:
                    continue
                if line.reward_id.reward_type != 'product':
                    reward_amount += line.price_subtotal
                else:
                    # Free product are 'regular' product lines with a price_unit of 0
                    reward_amount -= line.product_id.lst_price * line.product_uom_qty
            order.reward_amount = reward_amount

    def _compute_loyalty_data(self):
        self.loyalty_data = {}

        confirmed_so = self.filtered(lambda order: order.state == 'sale' and bool(order.id))
        if not confirmed_so:
            return

        loyalty_history_data = self.env['loyalty.history'].sudo()._read_group(
            domain=[
                ('order_id', 'in', confirmed_so.ids),
                ('order_model', '=', self._name),
            ],
            groupby=['order_id'],
            aggregates=['issued:sum', 'used:sum'],
        )
        loyalty_history_data_per_order = {
            order_id: {
                'total_issued': issued,
                'total_cost': cost,
            }
            for order_id, issued, cost in loyalty_history_data
        }
        for order in confirmed_so:
            if order.id not in loyalty_history_data_per_order:
                continue
            coupons = order.coupon_point_ids.coupon_id
            coupon_point_name = (len(coupons) == 1 and coupons.point_name) or _("Points")
            order.loyalty_data = {
                'point_name': coupon_point_name,
                'issued': loyalty_history_data_per_order[order.id]['total_issued'],
                'cost': loyalty_history_data_per_order[order.id]['total_cost'],
            }

    def _compute_gift_card_count(self):
        gift_card_data = dict(
            self.env['loyalty.card']._read_group(
                domain=[
                    ('order_id', 'in', self.ids),
                    ('program_type', '=', 'gift_card'),
                ],
                groupby=['order_id'],
                aggregates=['__count'],
            )
        )
        for order in self:
            order.gift_card_count = gift_card_data.get(order, 0)

    def _add_loyalty_history_lines(self):
        self.ensure_one()
        points_per_coupon = defaultdict(partial(defaultdict, int))
        for coupon_point in self.coupon_point_ids:
            points_per_coupon[coupon_point.coupon_id]['issued'] = coupon_point.points
        for line in self.order_line:
            if not line.coupon_id:
                continue
            points_per_coupon[line.coupon_id]['cost'] += line.points_cost

        create_values = []
        base_values = {
            'order_id': self.id,
            'order_model': self._name,
            'description': _("Order %s", self.display_name),
        }
        for coupon, point_dict in points_per_coupon.items():
            cost = point_dict.get('cost', 0.0)
            issued = point_dict.get('issued', 0.0)
            create_values.append({
                **base_values,
                'card_id': coupon.id,
                'used': cost,
                'issued': issued,
            })

        self.env['loyalty.history'].create(create_values)

    def _get_no_effect_on_threshold_lines(self):
        """Return the lines that have no effect on the minimum amount to reach."""
        self.ensure_one()
        return self.env['sale.order.line']

    def copy(self, default=None):
        new_orders = super().copy(default)
        reward_lines = new_orders.order_line.filtered('is_reward_line')
        if reward_lines:
            reward_lines.unlink()
        return new_orders

    def action_confirm(self):
        """
        Override to validate and update coupon rewards.

        If called with one SO, checks if there exists rewards that are available but not claimed,
        and if so returns a notification action.

        :raises ValidationError: A coupon gave a negative amount of points.
        :return: True or a notification action
        :rtype: bool | dict
        """
        for order in self:
            all_coupons = order.applied_coupon_ids | order.coupon_point_ids.coupon_id | order.order_line.coupon_id
            if any(order._get_real_points_for_coupon(coupon) < 0 for coupon in all_coupons):
                raise ValidationError(_("One or more rewards on the sale order is invalid. Please check them."))
            order._update_programs_and_rewards()
            order._add_loyalty_history_lines()
        has_claimable_rewards = len(self) == 1 and bool(self._get_claimable_rewards())

        # Remove any coupon from 'current' program that don't claim any reward.
        # This is to avoid ghost coupons that are lost forever.
        # Claiming a reward for that program will require either an automated check or a manual input again.
        reward_coupons = self.order_line.coupon_id
        self.coupon_point_ids.filtered(
            lambda pe: pe.coupon_id.program_id.applies_on == 'current' and pe.coupon_id not in reward_coupons
        ).coupon_id.sudo().unlink()
        # Add/remove the points to our coupons
        for coupon, change in self.filtered(lambda s: s.state != 'sale')._get_point_changes().items():
            coupon.points += change
        res = super().action_confirm()
        # Prioritize any action from super()
        if isinstance(res, bool) and has_claimable_rewards:
            res = {
                'type': 'ir.actions.client',
                'tag': 'display_notification',
                'params': {
                    'type': 'info',
                    'title': _("Rewards Available"),
                    'message': _("There are available rewards not added to this order."),
                    'next': {'type': 'ir.actions.act_window_close'},
                },
            }
        self._send_reward_coupon_mail()
        return res

    def _action_cancel(self):
        previously_confirmed = self.filtered(lambda s: s.state == 'sale')
        res = super()._action_cancel()

        order_history_lines = self.env['loyalty.history'].search([
            ('order_model', '=', self._name),
            ('order_id', 'in', previously_confirmed.ids),
        ])
        if order_history_lines:
            order_history_lines.sudo().unlink()

        # Add/remove the points to our coupons
        for coupon, changes in previously_confirmed.filtered(
            lambda s: s.state != 'sale'
        )._get_point_changes().items():
            coupon.points -= changes
        # Remove any rewards
        self.order_line.filtered(lambda l: l.is_reward_line).unlink()
        self.coupon_point_ids.coupon_id.sudo().filtered(
            lambda c: not c.program_id.is_nominative and c.order_id in self and not c.use_count)\
            .unlink()
        self.coupon_point_ids.sudo().unlink()
        return res

    def action_open_reward_wizard(self):
        self.ensure_one()
        self._update_programs_and_rewards()
        claimable_rewards = self._get_claimable_rewards()
        if len(claimable_rewards) == 1:
            coupon = next(iter(claimable_rewards))
            rewards = claimable_rewards[coupon]
            if len(rewards) == 1 and not rewards.multi_product:
                self._apply_program_reward(claimable_rewards[coupon], coupon)
                return True
        elif not claimable_rewards:
            return True
        return self.env['ir.actions.actions']._for_xml_id('sale_loyalty.sale_loyalty_reward_wizard_action')

    def action_view_gift_cards(self):
        self.ensure_one()
        return {
            'name': _("Gift Cards"),
            'type': 'ir.actions.act_window',
            'view_mode': 'list,form',
            'res_model': 'loyalty.card',
            'domain': [('order_id', '=', self.id), ('program_type', '=', 'gift_card')],
            'context': {'create': False},
        }

    def _send_reward_coupon_mail(self):
        coupons = self.env['loyalty.card']
        for order in self:
            coupons |= order._get_reward_coupons()
        if coupons:
            coupons._send_creation_communication(force_send=True)

    def _get_applied_global_discount_lines(self):
        """
        Returns the first line of the currently applied global discount or False
        """
        self.ensure_one()
        return self.order_line.filtered(lambda l: l.reward_id.is_global_discount)

    def _get_applied_global_discount(self):
        """
        Returns the currently applied global discount reward or False
        """
        return self._get_applied_global_discount_lines().reward_id

    def _get_reward_values_product(self, reward, coupon, product=None, **kwargs):
        """
        Returns an array of dict containing the values required for the reward lines
        """
        self.ensure_one()
        assert reward.reward_type == 'product'

        reward_products = reward.reward_product_ids
        product = product or reward_products[:1]
        if not product or product not in reward_products:
            raise UserError(_("Invalid product to claim."))
        taxes = self.fiscal_position_id.map_tax(product.taxes_id._filter_taxes_by_company(self.company_id))
        points = self._get_real_points_for_coupon(coupon)
        claimable_count = float_round(points / reward.required_points, precision_rounding=1, rounding_method='DOWN') if not reward.clear_wallet else 1
        cost = points if reward.clear_wallet else claimable_count * reward.required_points
        return [{
            'name': reward.description,
            'product_id': product.id,
            'discount': 100,
            'product_uom_qty': reward.reward_product_qty * claimable_count,
            'reward_id': reward.id,
            'coupon_id': coupon.id,
            'points_cost': cost,
            'reward_identifier_code': _generate_random_reward_code(),
            'sequence': max(self.order_line.filtered(lambda x: not x.is_reward_line).mapped('sequence'), default=10) + 1,
            'tax_ids': [Command.clear()] + [Command.link(tax.id) for tax in taxes],
        }]

    def _discountable_amount(self, rewards_to_ignore):
        """Compute the `discountable` amount for the current order, ignoring the provided rewards.

        :param rewards_to_ignore: the rewards to ignore from the total amount (if they were already
            applied on the order)
        :type rewards_to_ignore: `loyalty.reward` recordset

        :return: The discountable amount
        :rtype: float
        """
        self.ensure_one()

        discountable = 0

        for line in self.order_line - self._get_no_effect_on_threshold_lines():
            if rewards_to_ignore and line.reward_id in rewards_to_ignore:
                # Ignore the existing reward line if it was already applied
                continue
            if not line.product_uom_qty or not line.price_unit:
                # Ignore lines whose amount will be 0 (bc of empty qty or 0 price)
                continue
            tax_data = line.tax_ids.compute_all(
                line.price_unit,
                currency=line.currency_id,
                quantity=line.product_uom_qty,
                product=line.product_id,
                partner=line.order_partner_id,
            )
            # To compute the discountable amount we get the subtotal and add
            # non-fixed tax totals. This way fixed taxes will not be discounted
            taxes = line.tax_ids.filtered(lambda t: t.amount_type != 'fixed')
            discountable += tax_data['total_excluded'] + sum(
                tax['amount'] for tax in tax_data['taxes'] if tax['id'] in taxes.ids
            )
        return discountable

    def _discountable_order(self, reward):
        """Compute the `discountable` amount (and amounts per tax group) for the current order.

        :param reward: if provided, the reward whose discountable amounts must be computed.
            It must be applicable at the order level.
        :type reward: `loyalty.reward` record, can be empty to compute the amounts regardless of the
            program configuration

        :return: A tuple with the first element being the total discountable amount of the order,
            and the second a dictionary mapping each non-fixed taxes group to its corresponding
            total untaxed amount of the eligible order lines.
        :rtype: tuple(float, dict(account.tax: float))
        """
        self.ensure_one()
        reward.ensure_one()
        assert reward.discount_applicability == 'order'

        lines = self.order_line.filtered(lambda line: not line.display_type)
        if not reward.program_id.is_payment_program:
            # Gift cards and eWallets are applied on the total order amount
            # Other types of programs are not expected to apply on delivery lines
            lines -= self._get_no_effect_on_threshold_lines()
        else:
            # Prevent paying for a payment program's own top-up product using that
            # same program (e.g. topping up an eWallet by paying with the eWallet).
            top_up_products = reward.program_id.trigger_product_ids
            lines -= lines.filtered(lambda line: line.product_id in top_up_products)

        discountable = 0
        discountable_per_tax = defaultdict(float)

        AccountTax = self.env['account.tax']
        base_lines = []
        for line in lines:
            base_line = line._prepare_base_line_for_taxes_computation()
            taxes = base_line['tax_ids']
            discountable_taxes = base_line['tax_ids'].flatten_taxes_hierarchy()
            if not reward.program_id.is_payment_program:
                # To compute the discountable amount we get the subtotal and add
                # non-fixed tax totals. This way fixed taxes will not be discounted
                # This does not apply to Gift Cards and e-Wallet, where the total
                # order amount may be paid with the card balance
                taxes = taxes.filtered(lambda t: t.amount_type != 'fixed')
                discountable_taxes = discountable_taxes.filtered(lambda t: t.amount_type != 'fixed')
            base_line['discount_taxes'] = taxes
            base_line['discountable_taxes'] = discountable_taxes
            base_lines.append(base_line)
        AccountTax._add_tax_details_in_base_lines(base_lines, self.company_id)
        AccountTax._round_base_lines_tax_details(base_lines, self.company_id)

        def grouping_function(base_line, tax_data):
            if not tax_data:
                return None
            return {
                'taxes': base_line['discount_taxes'],
                'skip': (
                    tax_data['tax'] not in base_line['discountable_taxes']
                    or base_line['record'] not in lines
                ),
            }

        base_lines_aggregated_values = AccountTax._aggregate_base_lines_tax_details(base_lines, grouping_function)
        values_per_grouping_key = AccountTax._aggregate_base_lines_aggregated_values(base_lines_aggregated_values)
        for grouping_key, values in values_per_grouping_key.items():
            if grouping_key and grouping_key['skip']:
                continue

            taxes = grouping_key['taxes'] if grouping_key else self.env['account.tax']
            discountable += values['raw_base_amount_currency'] + values['raw_tax_amount_currency']
            discountable_per_tax[taxes] += (
                values['raw_base_amount_currency']
                + sum(
                    tax_data['raw_tax_amount_currency']
                    for base_line, taxes_data in values['base_line_x_taxes_data']
                    for tax_data in taxes_data
                    if tax_data['tax'].price_include
                )
            )
        return discountable, discountable_per_tax

    def _cheapest_line(self, reward):
        self.ensure_one()
        cheapest_line = False
        cheapest_line_price_unit = False
        domain = reward._get_discount_product_domain()
        for line in (self.order_line - self._get_no_effect_on_threshold_lines()):
            line_price_unit = self._get_order_line_price(line, 'price_unit')
            if (
                line.reward_id
                or line.combo_item_id
                or not line.product_uom_qty
                or not line_price_unit
                or not line.product_id.filtered_domain(domain)
            ):
                continue
            if not cheapest_line or cheapest_line_price_unit > line_price_unit:
                cheapest_line = line._get_lines_with_price()
                cheapest_line_price_unit = line_price_unit
        return cheapest_line

    def _discountable_cheapest(self, reward):
        """
        Returns the discountable and discountable_per_tax for a discount that applies to the cheapest line
        """
        self.ensure_one()
        assert reward.discount_applicability == 'cheapest'

        cheapest_line = self._cheapest_line(reward)
        if not cheapest_line:
            return False, False

        discountable = 0
        discountable_per_tax = defaultdict(int)
        for line in cheapest_line:
            discountable += line.price_total / line.product_uom_qty
            taxes = line.tax_ids.filtered(lambda t: t.amount_type != 'fixed')
            discountable_per_tax[taxes] += line.price_unit * (1 - (line.discount or 0) / 100)

        return discountable, discountable_per_tax

    def _get_specific_discountable_lines(self, reward):
        """
        Returns all lines to which `reward` can apply
        """
        self.ensure_one()
        assert reward.discount_applicability == 'specific'

        discountable_lines = self.env['sale.order.line']
        for line in (self.order_line - self._get_no_effect_on_threshold_lines()):
            domain = reward._get_discount_product_domain()
            if (
                not line.reward_id
                and not line.combo_item_id
                and line.product_id.filtered_domain(domain)
            ):
                discountable_lines |= line._get_lines_with_price()
        return discountable_lines

    def _discountable_specific(self, reward):
        """
        Special function to compute the discountable for 'specific' types of discount.
        The goal of this function is to make sure that applying a 5$ discount on an order with a
         5$ product and a 5% discount does not make the order go below 0.

        Returns the discountable and discountable_per_tax for a discount that only applies to specific products.
        """
        self.ensure_one()
        assert reward.discount_applicability == 'specific'

        lines_to_discount = self._get_specific_discountable_lines(reward).filtered(
            lambda line: bool(line.product_uom_qty and line.price_total)
        )
        discount_lines = defaultdict(lambda: self.env['sale.order.line'])
        order_lines = self.order_line - self._get_no_effect_on_threshold_lines()
        remaining_amount_per_line = defaultdict(int)
        for line in order_lines:
            if not line.product_uom_qty or not line.price_total:
                continue
            remaining_amount_per_line[line] = line.price_total
            if line.reward_id.reward_type == 'discount':
                discount_lines[line.reward_identifier_code] |= line

        order_lines -= self.order_line.filtered('reward_id')
        cheapest_line = False
        for lines in discount_lines.values():
            line_reward = lines.reward_id
            discounted_lines = order_lines
            if line_reward.discount_applicability == 'cheapest':
                # get the discounted cheapest line applicable for given reward domain
                cheapest_line = cheapest_line or self._cheapest_line(line_reward)
                discounted_lines = cheapest_line
            elif line_reward.discount_applicability == 'specific':
                discounted_lines = self._get_specific_discountable_lines(line_reward)
            if not discounted_lines:
                continue
            common_lines = discounted_lines & lines_to_discount
            if line_reward.discount_mode == 'percent':
                for line in discounted_lines:
                    if line_reward.discount_applicability == 'cheapest':
                        remaining_amount_per_line[line] *= (1 - line_reward.discount / 100 / line.product_uom_qty)
                    else:
                        remaining_amount_per_line[line] *= (1 - line_reward.discount / 100)
            else:
                non_common_lines = discounted_lines - lines_to_discount
                # Fixed prices are per tax
                discounted_amounts = defaultdict(int, {
                    sol.tax_ids.filtered(lambda t: t.amount_type != 'fixed'): abs(sol.price_total)
                    for sol in lines
                })
                for line in itertools.chain(non_common_lines, common_lines):
                    # For gift card and eWallet programs we have no tax but we can consume the amount completely
                    if lines.reward_id.program_id.is_payment_program:
                        discounted_amount = discounted_amounts[lines.tax_ids.filtered(lambda t: t.amount_type != 'fixed')]
                    else:
                        discounted_amount = discounted_amounts[line.tax_ids.filtered(lambda t: t.amount_type != 'fixed')]
                    if discounted_amount == 0:
                        continue
                    remaining = remaining_amount_per_line[line]
                    consumed = min(remaining, discounted_amount)
                    if lines.reward_id.program_id.is_payment_program:
                        discounted_amounts[lines.tax_ids.filtered(lambda t: t.amount_type != 'fixed')] -= consumed
                    else:
                        discounted_amounts[line.tax_ids.filtered(lambda t: t.amount_type != 'fixed')] -= consumed
                    remaining_amount_per_line[line] -= consumed

        discountable = 0
        discountable_per_tax = defaultdict(int)
        for line in lines_to_discount:
            discountable += remaining_amount_per_line[line]
            line_discountable = line.price_unit * line.product_uom_qty * (1 - (line.discount or 0.0) / 100.0)
            # line_discountable is the same as in a 'order' discount
            #  but first multiplied by a factor for the taxes to apply
            #  and then multiplied by another factor coming from the discountable
            taxes = line.tax_ids.filtered(lambda t: t.amount_type != 'fixed')
            discountable_per_tax[taxes] += line_discountable *\
                (remaining_amount_per_line[line] / line.price_total)
        return discountable, discountable_per_tax

    def _get_reward_values_discount(self, reward, coupon, **kwargs):
        self.ensure_one()
        assert reward.reward_type == 'discount'

        reward_applies_on = reward.discount_applicability
        reward_product = reward.discount_line_product_id
        reward_program = reward.program_id
        reward_currency = reward.currency_id
        sequence = max(
            self.order_line.filtered(lambda x: not x.is_reward_line).mapped('sequence'),
            default=10
        ) + 1
        base_reward_line_values = {
            'product_id': reward_product.id,
            'product_uom_qty': 1.0,
            'tax_ids': [Command.clear()],
            'name': reward.description,
            'reward_id': reward.id,
            'coupon_id': coupon.id,
            'sequence': sequence,
            'reward_identifier_code': _generate_random_reward_code(),
        }

        discountable = 0
        discountable_per_tax = defaultdict(int)
        if reward_applies_on == 'order':
            discountable, discountable_per_tax = self._discountable_order(reward)
        elif reward_applies_on == 'specific':
            discountable, discountable_per_tax = self._discountable_specific(reward)
        elif reward_applies_on == 'cheapest':
            discountable, discountable_per_tax = self._discountable_cheapest(reward)

        if not discountable:
            if not reward_program.is_payment_program and any(line.reward_id.program_id.is_payment_program for line in self.order_line):
                return [{
                    **base_reward_line_values,
                    'name': _("TEMPORARY DISCOUNT LINE"),
                    'price_unit': 0,
                    'product_uom_qty': 0,
                    'points_cost': 0,
                }]
            raise UserError(_("There is nothing to discount"))

        max_discount = reward_currency._convert(reward.discount_max_amount, self.currency_id, self.company_id, fields.Date.today()) or float('inf')
        # discount should never surpass the order's current total amount
        max_discount = min(self.amount_total, max_discount)
        if reward.discount_mode == 'per_point':
            points = self._get_real_points_for_coupon(coupon)
            if not reward_program.is_payment_program:
                # Rewards cannot be partially offered to customers
                points = points // reward.required_points * reward.required_points
            max_discount = min(max_discount,
                reward_currency._convert(reward.discount * points,
                    self.currency_id, self.company_id, fields.Date.today()))
        elif reward.discount_mode == 'per_order':
            max_discount = min(max_discount,
                reward_currency._convert(reward.discount, self.currency_id, self.company_id, fields.Date.today()))
        elif reward.discount_mode == 'percent':
            max_discount = min(max_discount, discountable * (reward.discount / 100))

        # Discount per taxes
        point_cost = reward.required_points if not reward.clear_wallet else self._get_real_points_for_coupon(coupon)
        if reward.discount_mode == 'per_point' and not reward.clear_wallet:
            # Calculate the actual point cost if the cost is per point
            converted_discount = self.currency_id._convert(min(max_discount, discountable), reward_currency, self.company_id, fields.Date.today())
            point_cost = coupon.currency_id.round(converted_discount / reward.discount)

        if reward_program.is_payment_program:  # Gift card / eWallet
            reward_line_values = {
                **base_reward_line_values,
                'price_unit': -min(max_discount, discountable),
                'points_cost': point_cost,
            }

            if reward_program.program_type == 'gift_card':
                # For gift cards, the SOL should consider the discount product taxes
                taxes_to_apply = reward_product.taxes_id._filter_taxes_by_company(self.company_id)
                if taxes_to_apply:
                    mapped_taxes = self.fiscal_position_id.map_tax(taxes_to_apply)
                    price_incl_taxes = mapped_taxes.filtered('price_include')
                    tax_res = mapped_taxes.with_context(
                        force_price_include=True,
                        round=False,
                        round_base=False,
                    ).compute_all(
                        reward_line_values['price_unit'],
                        currency=self.currency_id,
                    )
                    new_price = tax_res['total_excluded']
                    new_price += sum(
                        tax_data['amount']
                        for tax_data in tax_res['taxes']
                        if tax_data['id'] in price_incl_taxes.ids
                    )
                    reward_line_values.update({
                        'price_unit': new_price,
                        'tax_ids': [Command.set(mapped_taxes.ids)],
                    })
            return [reward_line_values]

        discount_factor = min(1, (max_discount / discountable)) if discountable else 1
        reward_dict = {}
        for tax, price in discountable_per_tax.items():
            if not price:
                continue
            mapped_taxes = self.fiscal_position_id.map_tax(tax)
            tax_desc = ''
            if len(discountable_per_tax) > 1 and any(t.name for t in mapped_taxes):
                tax_desc = _(
                    " - On products with the following taxes: %(taxes)s",
                    taxes=", ".join(mapped_taxes.mapped('name')),
                )
            reward_dict[tax] = {
                **base_reward_line_values,
                'name': _(
                    "Discount %(desc)s%(tax_str)s",
                    desc=reward.description,
                    tax_str=tax_desc,
                ) if mapped_taxes else reward.description,
                'price_unit': -(price * discount_factor),
                'points_cost': 0,
                'tax_ids': [Command.clear()] + [Command.link(tax.id) for tax in mapped_taxes]
            }
        # We only assign the point cost to one line to avoid counting the cost multiple times
        if reward_dict:
            reward_dict[next(iter(reward_dict))]['points_cost'] = point_cost
        # Returning .values() directly does not return a subscribable list
        return list(reward_dict.values())

    def _get_program_domain(self):
        """
        Returns the base domain that all programs have to comply to.
        """
        self.ensure_one()
        today = self._get_confirmed_tx_create_date()
        return [('active', '=', True), ('sale_ok', '=', True),
                *self.env['loyalty.program']._check_company_domain([self.company_id.id, self.company_id.parent_id.id]),
                '|', ('pricelist_ids', '=', False), ('pricelist_ids', 'in', [self.pricelist_id.id]),
                '|', ('date_from', '=', False), ('date_from', '<=', today),
                '|', ('date_to', '=', False), ('date_to', '>=', today)]

    def _get_trigger_domain(self):
        """
        Returns the base domain that all triggers have to comply to.
        """
        self.ensure_one()
        today = self._get_confirmed_tx_create_date()
        return [('active', '=', True), ('program_id.sale_ok', '=', True),
                *self.env['loyalty.program']._check_company_domain([self.company_id.id, self.company_id.parent_id.id]),
                '|', ('program_id.pricelist_ids', '=', False),
                     ('program_id.pricelist_ids', 'in', [self.pricelist_id.id]),
                '|', ('program_id.date_from', '=', False), ('program_id.date_from', '<=', today),
                '|', ('program_id.date_to', '=', False), ('program_id.date_to', '>=', today)]

    def _get_program_timezone(self):
        """Get the timezone to be used for loyalty date checking on the current order."""
        self.ensure_one()
        return (
            self.company_id.partner_id.tz
            or self.env['ir.config_parameter'].sudo().get_param('loyalty.timezone', 'UTC')
        )

    def _get_confirmed_tx_create_date(self):
        """Return the creation date of the earliest confirmed transaction to check which loyalty
        programs are applicable. If no transactions are confirmed, return the current day, using
        the company's time zone.
        """
        self.ensure_one()
        order_tz = self._get_program_timezone()
        confirmed_txs_dates = self.sudo().transaction_ids.filtered(
            lambda tx: tx.state in ('done', 'authorized'),
        ).mapped('create_date')
        if confirmed_txs_dates:
            # If order is getting confirmed, use the earliest finalized transaction's create date
            tx_date = min(confirmed_txs_dates)
            return tx_date.astimezone(timezone(order_tz)).date()
        return fields.Date.context_today(self.with_context(tz=order_tz))

    def _get_applicable_program_points(self, domain=None):
        """
        Returns a dict with the points per program for each (automatic) program that is applicable
        """
        self.ensure_one()
        if not domain:
            domain = [('trigger', '=', 'auto')]
        # Make sure domain always complies with the order's domain rules
        domain = Domain.AND([self._get_program_domain(), domain])
        # No other way than to test all programs to the order
        programs = self.env['loyalty.program'].search(domain)
        all_status = self._program_check_compute_points(programs)
        program_points = {p: status['points'][0] for p, status in all_status.items() if 'points' in status}
        return program_points

    def _get_points_programs(self):
        """
        Returns all programs that give points on the current order.
        """
        self.ensure_one()
        return self.coupon_point_ids.filtered('points').coupon_id.program_id

    def _get_reward_programs(self):
        """
        Returns all programs that are being used for rewards.
        """
        self.ensure_one()
        return self.order_line.reward_id.program_id

    def _get_reward_coupons(self):
        """
        Returns all coupons that are a reward.
        """
        self.ensure_one()
        return self.coupon_point_ids.filtered('points').coupon_id.filtered(
            lambda c: c.program_id.applies_on == 'future',
        )

    def _get_applied_programs(self):
        """
        Returns all applied programs on current order.

        Applied programs is the combination of both new points for your order and the programs linked to rewards.
        """
        self.ensure_one()
        return self._get_points_programs() | self._get_reward_programs()

    def _recompute_prices(self):
        """Recompute coupons/promotions after pricelist prices reset."""
        super()._recompute_prices()
        for order in self:
            if any(line.is_reward_line for line in order.order_line):
                order._update_programs_and_rewards()

    def _get_point_changes(self):
        """
        Returns the changes in points per coupon as a dict.

        Used when validating/cancelling an order
        """
        points_per_coupon = defaultdict(lambda: 0)
        for coupon_point in self.coupon_point_ids:
            points_per_coupon[coupon_point.coupon_id] += coupon_point.points
        for line in self.order_line:
            if not line.reward_id or not line.coupon_id:
                continue
            points_per_coupon[line.coupon_id] -= line.points_cost
        return points_per_coupon

    def _get_real_points_for_coupon(self, coupon, post_confirm=False):
        """
        Returns the actual points usable for this coupon for this order. Set pos_confirm to True to include points for future orders.

        This is calculated by taking the points on the coupon, the points the order will give to the coupon (if applicable) and removing the points taken by already applied rewards.
        """
        self.ensure_one()
        points = coupon.points
        if self.state not in ('sale', 'done'):
            if coupon.program_id.applies_on != 'future':
                # Points that will be given by the order upon confirming the order
                points += self.coupon_point_ids.filtered(lambda p: p.coupon_id == coupon).points
            # Points already used by rewards
            points -= sum(self.order_line.filtered(lambda l: l.coupon_id == coupon).mapped('points_cost'))
        if any(rule.reward_point_mode == 'money' for rule in coupon.program_id.rule_ids):
            points = coupon.currency_id.round(points)
        return points

    def _add_points_for_coupon(self, coupon_points):
        """
        Updates (or creates) an entry in coupon_point_ids for the given coupons.
        """
        self.ensure_one()
        if self.state == 'sale':
            for coupon, points in coupon_points.items():
                coupon.sudo().points += points
        for pe in self.coupon_point_ids.sudo():
            if pe.coupon_id in coupon_points:
                pe.points = coupon_points.pop(pe.coupon_id)
        if coupon_points:
            self.sudo().with_context(tracking_disable=True).write({
                'coupon_point_ids': [(0, 0, {
                    'coupon_id': coupon.id,
                    'points': points,
                }) for coupon, points in coupon_points.items()]
            })

    def _update_loyalty_history(self, coupon_id, points):
        self.ensure_one()
        order_coupon_history = self.env['loyalty.history'].search([
            ('card_id', '=', coupon_id.id),
            ('order_model', '=', self._name),
            ('order_id', '=', self.id),
        ], limit=1)
        if order_coupon_history:
            order_coupon_history.update({'used': order_coupon_history.used + points})
        else:
            issued = self.coupon_point_ids.filtered(lambda p: p.coupon_id == coupon_id).points
            self.env['loyalty.history'].create({
                'card_id': coupon_id.id,
                'order_model': self._name,
                'order_id': self.id,
                'description': _("Order %s", self.display_name),
                'issued': issued,
                'used': points,
            })

    def _remove_program_from_points(self, programs):
        self.coupon_point_ids.filtered(lambda p: p.coupon_id.program_id in programs).sudo().unlink()

    def _get_reward_line_values(self, reward, coupon, **kwargs):
        self.ensure_one()
        self = self.with_context(lang=self._get_lang())
        reward = reward.with_context(lang=self._get_lang())
        if reward.reward_type == 'discount':
            return self._get_reward_values_discount(reward, coupon, **kwargs)
        elif reward.reward_type == 'product':
            return self._get_reward_values_product(reward, coupon, **kwargs)

    def _write_vals_from_reward_vals(self, reward_vals, old_lines, delete=True):
        """
        Update, create new reward line and delete old lines in one write on `order_line`

        Returns the untouched old lines.
        """
        self.ensure_one()
        command_list = []
        for vals, line in zip(reward_vals, old_lines):
            if vals['product_id'] == line.product_id.id:
                vals['name'] = line.name  # Preserve custom description
            command_list.append((Command.UPDATE, line.id, vals))
        if len(reward_vals) > len(old_lines):
            command_list.extend((Command.CREATE, 0, vals) for vals in reward_vals[len(old_lines):])
        elif len(reward_vals) < len(old_lines) and delete:
            command_list.extend((Command.DELETE, line.id) for line in old_lines[len(reward_vals):])
        self.write({'order_line': command_list})
        return self.env['sale.order.line'] if delete else old_lines[len(reward_vals):]

    def _best_global_discount_already_applied(self, current_reward, new_reward, discountable=None):
        """Determine whether current_reward is better than new_reward.

        This function compares the discount amount of two rewards to determine whether the current
        one is better than another one.

        Notes
        -----

            If the discount amounts of both the current and the new rewards exceed the order total,
            the reward with the smaller discount amount is considered the best.
            This is to ensure that the most advantageous discount is applied for the customer,
            who will keep the most important voucher, having saved the same amount in the end.

        :param loyalty.reward current_reward: The reward currently applied on the sale order.
        :param loyalty.reward new_reward: The reward to compare with.
        :param float discountable: The total discountable amount of the sale order.
            If not provided, it will be calculated on the fly.
        :return: True if current_reward is considered better than new_reward.
        :rtype: bool
        """
        self.ensure_one()
        current_reward.ensure_one()
        new_reward.ensure_one()

        if current_reward == new_reward:
            return True

        if discountable is None:  # Only recompute if discountable is not given, not if its zero
            discountable = self._discountable_amount(current_reward)

        discount_current_reward = self._get_discount_amount(current_reward, discountable)
        discount_new_reward = self._get_discount_amount(new_reward, discountable)

        discount_current_bigger_than_discountable = self.currency_id.compare_amounts(
            amount1=discount_current_reward,
            amount2=discountable,
        ) >= 0
        discount_new_bigger_than_discountable = self.currency_id.compare_amounts(
            amount1=discount_new_reward,
            amount2=discountable,
        ) >= 0
        compare_current_and_new_reward = self.currency_id.compare_amounts(
            amount1=discount_current_reward,
            amount2=discount_new_reward,
        )

        if discount_current_bigger_than_discountable and discount_new_bigger_than_discountable:
            # If both discounts are greater than the discountable amount, the lower discount
            # is better as it reduces the discount amount 'spent' by the customer.
            return compare_current_and_new_reward <= 0

        # Return True only if the discount of the new reward is greater than the current reward
        # discount.
        return compare_current_and_new_reward >= 0

    def _get_discount_amount(self, reward, discountable):
        """Compute the discount amount for the given reward, w.r.t. the discountable amount.

        :param loyalty.reward reward: The reward for which to calculate the maximum discount.
        :param float discountable: The total discountable amount of the sale order.
        :return: The maximum discount amount.
        :rtype: float
        """
        if reward.discount_mode == 'per_order':
            return reward.currency_id._convert(
                from_amount=reward.discount,
                to_currency=self.currency_id,
                company=self.company_id,
                date=fields.Date.today(),
            )
        elif reward.discount_mode == 'percent':
            return discountable * (reward.discount / 100)

    def _apply_program_reward(self, reward, coupon, **kwargs):
        """
        Applies the reward to the order provided the given coupon has enough points.
        This method does not check for program rules.

        This method also assumes the points added by the program triggers have already been computed.
        The temporary points are used if the program is applicable to the current order.

        Returns a dict containing the error message or empty if everything went correctly.
        NOTE: A call to `_update_programs_and_rewards` is expected to reorder the discounts.
        """
        self.ensure_one()
        # Use the old lines before creating new ones. These should already be in a 'reset' state.
        old_reward_lines = kwargs.get('old_lines', self.env['sale.order.line'])
        if reward.is_global_discount:
            global_discount_reward_lines = self._get_applied_global_discount_lines()
            global_discount_reward = global_discount_reward_lines.reward_id
            if (
                global_discount_reward
                and global_discount_reward != reward
                and self._best_global_discount_already_applied(global_discount_reward, reward)
            ):
                return {'error': _("A better global discount is already applied.")}
            elif global_discount_reward and global_discount_reward != reward:
                # Invalidate the old global discount as it may impact the new discount to apply
                global_discount_reward_lines._reset_loyalty(True)
                old_reward_lines |= global_discount_reward_lines
        if not reward.program_id.is_nominative and reward.program_id.applies_on == 'future' and coupon in self.coupon_point_ids.coupon_id:
            return {'error': _("The coupon can only be claimed on future orders.")}
        elif self._get_real_points_for_coupon(coupon) < reward.required_points:
            return {'error': _("The coupon does not have enough points for the selected reward.")}
        reward_vals = self._get_reward_line_values(reward, coupon, **kwargs)
        self._write_vals_from_reward_vals(reward_vals, old_reward_lines)
        return {}

    def _get_claimable_rewards(self, forced_coupons=None):
        """
        Fetch all rewards that are currently claimable from all concerned coupons,
         meaning coupons from applied programs and applied rewards or the coupons given as parameter.

        Returns a dict containing the all the claimable rewards grouped by coupon.
        Coupons that can not claim any reward are not contained in the result.
        """
        self.ensure_one()
        result = defaultdict(lambda: self.env['loyalty.reward'])

        check_date = self._get_confirmed_tx_create_date()

        all_coupons = forced_coupons or (self.coupon_point_ids.coupon_id | self.order_line.coupon_id | self.applied_coupon_ids)
        if not all_coupons:
            return result

        has_payment_reward = any(line.reward_id.program_id.is_payment_program for line in self.order_line)
        global_discount_reward = self._get_applied_global_discount()
        active_products_domain = self.env['loyalty.reward']._get_active_products_domain()

        # Only evaluate discountable amount if needed
        discountable = lazy(lambda: self._discountable_amount(global_discount_reward))
        total_is_zero = lazy(lambda: self.currency_id.is_zero(discountable))

        for coupon in all_coupons:
            # Skip coupons generated by this order that only apply on future orders
            if coupon.program_id.applies_on == 'future' and coupon.order_id == self:
                continue
            if coupon.expiration_date and coupon.expiration_date < check_date:
                continue
            points = self._get_real_points_for_coupon(coupon)
            for reward in coupon.program_id.reward_ids:
                if (
                    reward.is_global_discount
                    and global_discount_reward
                    and self._best_global_discount_already_applied(
                        global_discount_reward, reward, discountable
                    )
                ):
                    continue
                # Discounts are not allowed if the total is zero unless there is a payment reward, in which case we allow discounts.
                # If the total is 0 again without the payment reward it will be removed.
                is_discount = reward.reward_type == 'discount'
                is_payment_program = reward.program_id.is_payment_program
                if is_discount and total_is_zero and (not has_payment_reward or is_payment_program):
                    continue
                # Skip discount that has already been applied if not part of a payment program
                if is_discount and not is_payment_program and reward in self.order_line.reward_id:
                    continue
                if reward.reward_type == 'product' and not reward.filtered_domain(
                    active_products_domain
                ):
                    continue
                if points >= reward.required_points:
                    result[coupon] |= reward
        return result

    def _allow_nominative_programs(self):
        """
        Whether or not this order may use nominative programs.
        """
        self.ensure_one()
        return True

    def _update_programs_and_rewards(self):
        """
        Updates applied programs's given points with the current state of the order.
        Checks automatic programs for applicability.
        Updates applied rewards using the new points and the current state of the order (for example with % discounts).
        """
        self.ensure_one()

        # +===================================================+
        # |       STEP 1: Retrieve all applicable programs    |
        # +===================================================+

        # Automatically load in eWallet and loyalty cards coupons with previously received points
        if self._allow_nominative_programs():
            loyalty_card = self.env['loyalty.card'].search([
                ('id', 'not in', self.applied_coupon_ids.ids),
                ('partner_id', '=', self.partner_id.id),
                ('points', '>', 0),
                '|', ('program_id.program_type', '=', 'ewallet'),
                     '&', ('program_id.program_type', '=', 'loyalty'),
                          ('program_id.applies_on', '!=', 'current'),
            ])
            if loyalty_card:
                self.applied_coupon_ids += loyalty_card
        # Programs that are applied to the order and count points
        points_programs = self._get_points_programs()
        # Coupon programs that require the program's rules to match but do not count for points
        coupon_programs = self.applied_coupon_ids.program_id
        # Programs that are automatic and not yet applied
        program_domain = self._get_program_domain()
        domain = Domain.AND([program_domain, [('id', 'not in', points_programs.ids), ('trigger', '=', 'auto'), ('rule_ids.mode', '=', 'auto')]])
        automatic_programs = self.env['loyalty.program'].search(domain).filtered(lambda p:
            not p.limit_usage or p.total_order_count < p.max_usage)

        all_programs_to_check = points_programs | coupon_programs | automatic_programs
        all_coupons = self.coupon_point_ids.coupon_id | self.applied_coupon_ids
        # First basic check using the program_domain -> for example if a program gets archived mid quotation
        domain_matching_programs = all_programs_to_check.filtered_domain(program_domain)
        all_programs_status = {p: {'error': 'error'} for p in all_programs_to_check - domain_matching_programs}
        # Compute applicability and points given for all programs that passed the domain check
        # Note that points are computed with reward lines present
        all_programs_status.update(self._program_check_compute_points(domain_matching_programs))
        # Delay any unlink to the end of the function since they cause a full cache invalidation
        lines_to_unlink = self.env['sale.order.line']
        coupons_to_unlink = self.env['loyalty.card']
        point_entries_to_unlink = self.env['sale.order.coupon.points']
        # Remove any coupons that are expired
        if initial_coupons := self.applied_coupon_ids:
            check_date = self._get_confirmed_tx_create_date()
            self.applied_coupon_ids = initial_coupons.filtered(
                lambda c: not c.expiration_date or c.expiration_date >= check_date,
            )
            removed = initial_coupons - self.applied_coupon_ids
            lines_to_unlink |= self.order_line.filtered(lambda sol: sol.coupon_id in removed)
        point_ids_per_program = defaultdict(lambda: self.env['sale.order.coupon.points'])
        for pe in self.coupon_point_ids:
            # Update coupons that were created for Public User
            if pe.coupon_id.partner_id.is_public and not self.partner_id.is_public:
                pe.coupon_id.partner_id = self.partner_id
            # Remove any point entry for a coupon that does not belong to the customer
            if pe.coupon_id.partner_id and pe.coupon_id.partner_id != self.partner_id:
                pe.points = 0
                point_entries_to_unlink |= pe
            else:
                point_ids_per_program[pe.coupon_id.program_id] |= pe

        # +==========================================+
        # |       STEP 2: Update applied programs    |
        # +==========================================+

        # Programs that were not applied via a coupon
        for program in points_programs:
            status = all_programs_status[program]
            program_point_entries = point_ids_per_program[program]
            if 'error' in status:
                # Program is not applicable anymore
                coupons_from_order = program_point_entries.coupon_id.filtered(lambda c: c.order_id == self)
                all_coupons -= coupons_from_order
                # Invalidate those lines so that they don't impact anything further down the line
                program_reward_lines = self.order_line.filtered(lambda l: l.coupon_id in coupons_from_order)
                program_reward_lines._reset_loyalty(True)
                lines_to_unlink |= program_reward_lines
                # Delete coupon created by this order for this program if it is not nominative
                if not program.is_nominative:
                    coupons_to_unlink |= coupons_from_order
                else:
                    # Only remove the coupon_point_id
                    point_entries_to_unlink |= program_point_entries
                    point_entries_to_unlink.points = 0
                # Remove the code activated rules
                self.code_enabled_rule_ids -= program.rule_ids
            else:
                # Program stays applicable, update our points
                all_point_changes = [p for p in status['points'] if p]
                if not all_point_changes and program.is_nominative:
                    all_point_changes = [0]
                for pe, points in zip(program_point_entries.sudo(), all_point_changes):
                    pe.points = points
                if len(program_point_entries) < len(all_point_changes):
                    new_coupon_points = all_point_changes[len(program_point_entries):]
                    # next_order_coupons should be linked to the order's partner
                    partner_id = program.program_type == 'next_order_coupons' and self.partner_id.id
                    # NOTE: Maybe we could batch the creation of coupons across multiple programs but this really only applies to gift cards
                    new_coupons = self.env['loyalty.card'].with_context(loyalty_no_mail=True, tracking_disable=True).create([{
                        'program_id': program.id,
                        'partner_id': partner_id,
                        'points': 0,
                        'order_id': self.id,
                    } for _ in new_coupon_points])
                    self._add_points_for_coupon({coupon: x for coupon, x in zip(new_coupons, new_coupon_points)})
                elif len(program_point_entries) > len(all_point_changes):
                    point_ids_to_unlink = program_point_entries[len(all_point_changes):]
                    all_coupons -= point_ids_to_unlink.coupon_id
                    coupons_to_unlink |= point_ids_to_unlink.coupon_id
                    point_ids_to_unlink.points = 0

        # Programs applied using a coupon
        applied_coupon_per_program = defaultdict(lambda: self.env['loyalty.card'])
        for coupon in self.applied_coupon_ids:
            applied_coupon_per_program[coupon.program_id] |= coupon
        for program in coupon_programs:
            if program not in domain_matching_programs or\
                (program.applies_on == 'current' and 'error' in all_programs_status[program]):
                program_reward_lines = self.order_line.filtered(lambda l: l.coupon_id in applied_coupon_per_program[program])
                program_reward_lines._reset_loyalty(True)
                lines_to_unlink |= program_reward_lines
                self.applied_coupon_ids -= applied_coupon_per_program[program]
                all_coupons -= applied_coupon_per_program[program]

        # +==========================================+
        # |       STEP 3: Update reward lines        |
        # +==========================================+

        # We will reuse these lines as much as possible, this resets the order in a reward-less state
        reward_line_pool = self.order_line.filtered(lambda l: l.reward_id and l.coupon_id)._reset_loyalty()
        seen_rewards = set()
        line_rewards = []
        payment_rewards = [] # gift_card and ewallet are considered as payments and should always be applied last
        for line in self.order_line:
            if line.reward_identifier_code in seen_rewards or not line.reward_id or\
                not line.coupon_id:
                continue
            seen_rewards.add(line.reward_identifier_code)
            if line.reward_id.program_id.is_payment_program:
                payment_rewards.append((line.reward_id, line.coupon_id, line.reward_identifier_code, line.product_id))
            else:
                line_rewards.append((line.reward_id, line.coupon_id, line.reward_identifier_code, line.product_id))

        for reward_key in itertools.chain(line_rewards, payment_rewards):
            coupon = reward_key[1]
            reward = reward_key[0]
            program = reward.program_id
            points = self._get_real_points_for_coupon(coupon)
            if coupon not in all_coupons or points < reward.required_points or program not in domain_matching_programs:
                # Reward is not applicable anymore, the reward lines will simply be removed at the end of this function
                continue
            try:
                values_list = self._get_reward_line_values(reward, coupon, product=reward_key[3])
            except UserError:
                # It could happen that we have nothing to discount after changing the order.
                values_list = []
            reward_line_pool = self._write_vals_from_reward_vals(values_list, reward_line_pool, delete=False)

        lines_to_unlink |= reward_line_pool

        # +==========================================+
        # |       STEP 4: Apply new programs         |
        # +==========================================+

        for program in automatic_programs:
            program_status = all_programs_status[program]
            if 'error' in program_status:
                continue
            self.__try_apply_program(program, False, program_status)

        # +==========================================+
        # |       STEP 5: Cleanup                    |
        # +==========================================+

        order_line_update = [(Command.DELETE, line.id) for line in lines_to_unlink]
        if order_line_update:
            self.write({'order_line': order_line_update})
        if coupons_to_unlink:
            coupons_to_unlink.sudo().unlink()
        if point_entries_to_unlink:
            point_entries_to_unlink.sudo().unlink()

    def _get_not_rewarded_order_lines(self):
        return self.order_line.filtered(lambda line: line.product_id and not line.reward_id)

    def _get_order_line_price(self, order_line, price_type):
        return sum(order_line._get_lines_with_price().mapped(price_type))

    def _program_check_compute_points(self, programs):
        """
        Checks the program validity from the order lines aswell as computing the number of points to add.

        Returns a dict containing the error message or the points that will be given with the keys 'points'.
        """
        self.ensure_one()

        # Prepare quantities
        order_lines = self._get_not_rewarded_order_lines().filtered(
            lambda line: not line.combo_item_id
        )
        products = order_lines.product_id
        products_qties = dict.fromkeys(products, 0)
        for line in order_lines:
            product_qty = line.product_uom_id._compute_quantity(
                line.product_uom_qty, line.product_id.uom_id
            )
            products_qties[line.product_id] += product_qty
        # Contains the products that can be applied per rule
        products_per_rule = programs._get_valid_products(products)

        # Prepare amounts
        so_products_per_rule = programs._get_valid_products(self.order_line.product_id)
        lines_per_rule = defaultdict(lambda: self.env['sale.order.line'])
        # Skip lines that have no effect on the minimum amount to reach.
        for line in self.order_line - self._get_no_effect_on_threshold_lines():
            is_discount = line.reward_id.reward_type == 'discount'
            reward_program = line.reward_id.program_id
            # Skip lines for automatic discounts, as well as combo item lines.
            if (is_discount and reward_program.trigger == 'auto') or line.combo_item_id:
                continue
            for program in programs:
                # Skip lines for the current program's discounts.
                if is_discount and reward_program == program:
                    continue
                for rule in program.rule_ids:
                    # Skip lines to which the rule doesn't apply.
                    if line.product_id in so_products_per_rule.get(rule, []):
                        lines_per_rule[rule] |= line._get_lines_with_price()

        result = {}
        for program in programs:
            # Used for error messages
            # By default False, but True if no rules and applies_on current -> misconfigured coupons program
            code_matched = not bool(program.rule_ids) and program.applies_on == 'current' # Stays false if all triggers have code and none have been activated
            minimum_amount_matched = code_matched
            product_qty_matched = code_matched
            points = 0
            # Some rules may split their points per unit / money spent
            #  (i.e. gift cards 2x50$ must result in two 50$ codes)
            rule_points = []
            program_result = result.setdefault(program, dict())
            for rule in program.rule_ids:
                # prevent bottomless ewallet spending
                if program.program_type == 'ewallet' and not program.trigger_product_ids:
                    break
                if rule.mode == 'with_code' and rule not in self.code_enabled_rule_ids:
                    continue
                code_matched = True
                rule_amount = rule._compute_amount(self.currency_id)
                untaxed_amount = sum(lines_per_rule[rule].mapped('price_subtotal'))
                tax_amount = sum(lines_per_rule[rule].mapped('price_tax'))
                if rule_amount > (rule.minimum_amount_tax_mode == 'incl' and (untaxed_amount + tax_amount) or untaxed_amount):
                    continue
                minimum_amount_matched = True
                if not products_per_rule.get(rule):
                    continue
                rule_products = products_per_rule[rule]
                ordered_rule_products_qty = sum(products_qties[product] for product in rule_products)
                if ordered_rule_products_qty < rule.minimum_qty or not rule_products:
                    continue
                product_qty_matched = True
                if not rule.reward_point_amount:
                    continue
                # Count all points separately if the order is for the future and the split option is enabled
                if program.applies_on == 'future' and rule.reward_point_split and rule.reward_point_mode != 'order':
                    if rule.reward_point_mode == 'unit':
                        rule_points.extend(rule.reward_point_amount for _ in range(int(ordered_rule_products_qty)))
                    elif rule.reward_point_mode == 'money':
                        for line in self.order_line:
                            if (
                                line.is_reward_line
                                or line.combo_item_id
                                or line.product_id not in rule_products
                                or line.product_uom_qty <= 0
                            ):
                                continue
                            line_price_total = self._get_order_line_price(line, 'price_total')
                            points_per_unit = float_round(
                                (rule.reward_point_amount * line_price_total / line.product_uom_qty),
                                precision_digits=2, rounding_method='DOWN')
                            if not points_per_unit:
                                continue
                            rule_points.extend([points_per_unit] * int(line.product_uom_qty))
                else:
                    # All checks have been passed we can now compute the points to give
                    if rule.reward_point_mode == 'order':
                        points += rule.reward_point_amount
                    elif rule.reward_point_mode == 'money':
                        # Compute amount paid for rule
                        # NOTE: this accounts for discounts -> 1 point per $ * (100$ - 30%) will
                        # result in 70 points
                        amount_paid = 0.0
                        rule_products = so_products_per_rule.get(rule, [])
                        for line in self.order_line - self._get_no_effect_on_threshold_lines():
                            if line.combo_item_id or line.reward_id.program_id.program_type in [
                                'ewallet', 'gift_card', program.program_type
                            ]:
                                continue
                            line_price_total = self._get_order_line_price(line, 'price_total')
                            amount_paid += (
                                line_price_total if line.product_id in rule_products
                                else 0.0
                            )

                        points += float_round(rule.reward_point_amount * amount_paid, precision_digits=2, rounding_method='DOWN')
                    elif rule.reward_point_mode == 'unit':
                        points += rule.reward_point_amount * ordered_rule_products_qty
            # NOTE: for programs that are nominative we always allow the program to be 'applied' on the order
            #  with 0 points so that `_get_claimable_rewards` returns the rewards associated with those programs
            if not program.is_nominative:
                if not code_matched:
                    program_result['error'] = _("This program requires a code to be applied.")
                elif not minimum_amount_matched:
                    program_result['error'] = _(
                        "To take advantage of this offer, your order must include at least %(amount)s %(currency)s of the eligible products.",
                        amount=min(program.rule_ids.mapped('minimum_amount')),
                        currency=program.currency_id.name,
                    )
                elif not product_qty_matched:
                    program_result['error'] = _("You don't have the required product quantities on your sales order.")
            elif self.partner_id.is_public and not self._allow_nominative_programs():
                program_result['error'] = _("This program is not available for public users.")
            if 'error' not in program_result:
                points_result = [points] + rule_points
                program_result['points'] = points_result
        return result

    def __try_apply_program(self, program, coupon, status):
        self.ensure_one()
        all_points = status['points']
        points = all_points[0]
        coupons = coupon or self.env['loyalty.card']
        if coupon:
            if program.is_nominative:
                self._add_points_for_coupon({coupon: points})
        elif not coupon:
            # If the program only applies on the current order it does not make sense to fetch already existing coupons
            if program.is_nominative:
                coupon = self.env['loyalty.card'].search(
                    [('partner_id', '=', self.partner_id.id), ('program_id', '=', program.id)], limit=1)
                # Do not apply 'nominative' programs if no point is given and no coupon exists
                if not points and not coupon:
                    return {'error': _("No card found for this loyalty program and no points will be given with this order.")}
                elif coupon:
                    self._add_points_for_coupon({coupon: points})
                coupons = coupon
            if not coupon:
                all_points = [p for p in all_points if p]
                partner = False
                # Loyalty programs and ewallets are nominative
                if program.is_nominative or program.program_type == 'next_order_coupons':
                    partner = self.partner_id.id
                coupons = self.env['loyalty.card'].sudo().with_context(loyalty_no_mail=True, tracking_disable=True).create([{
                    'program_id': program.id,
                    'partner_id': partner,
                    'points': 0,
                    'order_id': self.id,
                } for _ in all_points])
                self._add_points_for_coupon({coupon: x for coupon, x in zip(coupons, all_points)})
        return {'coupon': coupons}

    def _try_apply_program(self, program, coupon=None):
        """
        Tries to apply a program using the coupon if provided.

        This function provides the full routine to apply a program, it will check for applicability
        aswell as creating the necessary coupons and co-models to give the points to the customer.

        This function does not apply any reward to the order, rewards have to be given manually.

        Returns a dict containing the error message or containing the associated coupon(s).
        """
        self.ensure_one()
        # Basic checks
        if not program.filtered_domain(self._get_program_domain()):
            return {'error': _("The program is not available for this order.")}
        elif program in self._get_applied_programs():
            return {'error': _("This program is already applied to this order."), 'already_applied': True}
        elif program.reward_ids:
            global_rewards = program.reward_ids.filtered('is_global_discount')
            applied_global_reward = self._get_applied_global_discount()
            best_global_rewards = max(
                global_rewards,
                key=lambda reward: self._get_discount_amount(
                    reward, self._discountable_amount(applied_global_reward)
                )
            ) if len(global_rewards) > 1 else global_rewards
            if (
                best_global_rewards
                and applied_global_reward
                and self._best_global_discount_already_applied(applied_global_reward, best_global_rewards)
            ):
                return {'error': _(
                    "This discount (%(discount)s) is not compatible with \"%(other_discount)s\". "
                    "Please remove it in order to apply this one.",
                    discount=best_global_rewards.description,
                    other_discount=applied_global_reward.description
                )}
        # Check for applicability from the program's triggers/rules.
        # This step should also compute the amount of points to give for that program on that order.
        status = self._program_check_compute_points(program)[program]
        if 'error' in status:
            return status
        return self.__try_apply_program(program, coupon, status)

    def _try_apply_code(self, code):
        """
        Tries to apply a promotional code to the sales order.
        It can be either from a coupon or a program rule.

        Returns a dict with the following possible keys:
         - 'not_found': Populated with True if the code did not yield any result.
         - 'error': Any error message that could occur.
         OR The result of `_get_claimable_rewards` with the found or newly created coupon, it will be empty if the coupon was consumed completely.
        """
        self.ensure_one()

        base_domain = self._get_trigger_domain()
        domain = Domain.AND([base_domain, [('mode', '=', 'with_code'), ('code', '=', code)]])
        rule = self.env['loyalty.rule'].search(domain)
        program = rule.program_id
        coupon = False
        check_date = self._get_confirmed_tx_create_date()

        if (
            rule in self.code_enabled_rule_ids
            and program in self.order_line.filtered("is_reward_line").reward_id.program_id
        ):
            return {'error': _("This promo code is already applied.")}

        # No trigger was found from the code, try to find a coupon
        if not program:
            coupon = self.env['loyalty.card'].search([('code', '=', code)])
            if not coupon or\
                not coupon.program_id.active or\
                not coupon.program_id.reward_ids or\
                not coupon.program_id.filtered_domain(self._get_program_domain()):
                return {'error': _("This code is invalid (%s).", code), 'not_found': True}
            if coupon.expiration_date and coupon.expiration_date < check_date:
                return {'error': _("This coupon is expired.")}
            elif coupon.points < min(coupon.program_id.reward_ids.mapped('required_points')):
                return {'error': _("This coupon has already been used.")}
            program = coupon.program_id

        if not program or not program.active:
            return {'error': _("This code is invalid (%s).", code), 'not_found': True}
        elif program.program_type in ('loyalty', 'ewallet'):
            return {'error': _("This program cannot be applied with code.")}

        # Lock the loyalty program row to block several processes that try to
        # read it at the same time. We also use NOWAIT to make sure we trigger a
        # serialization error when the processes don't have the lock and thus,
        # trigger a retry of the transaction.
        self.env.cr.execute("""
            SELECT id FROM loyalty_program WHERE id=%s FOR UPDATE NOWAIT
        """, (program.id,))

        if (program.limit_usage and program.total_order_count >= program.max_usage):
            return {'error': _("This code is expired (%s).", code)}

        # Rule will count the next time the points are updated
        if rule:
            self.code_enabled_rule_ids |= rule
        program_is_applied = program in self._get_points_programs()
        # Condition that need to apply program (if not applied yet):
        # current -> always
        # future -> if no coupon
        # nominative -> non blocking if card exists with points
        if coupon:
            self.applied_coupon_ids += coupon
        if program_is_applied:
            # Update the points for our programs, this will take the new trigger in account
            self._update_programs_and_rewards()
        elif program.applies_on != 'future' or not coupon:
            apply_result = self._try_apply_program(program, coupon)
            if 'error' in apply_result and (not program.is_nominative or (program.is_nominative and not coupon)):
                if rule:
                    self.code_enabled_rule_ids -= rule
                if coupon and not apply_result.get('already_applied', False):
                    self.applied_coupon_ids -= coupon
                return apply_result
            coupon = apply_result.get('coupon', self.env['loyalty.card'])
        return self._get_claimable_rewards(forced_coupons=coupon)

    def _validate_order(self):
        """
        Override of sale to create invoice for zero amount order. If the order total is zero and
        automatic invoicing is enabled, it creates and posts an invoice.

        :return: None
        """
        super()._validate_order()
        if self.amount_total or not self.reward_amount:
            return
        auto_invoice = self.env['ir.config_parameter'].get_param('sale.automatic_invoice')
        if str2bool(auto_invoice):
            # create an invoice for order with zero total amount and automatic invoice enabled
            self._force_lines_to_invoice_policy_order()
            invoice = self._create_invoices(final=True)
            invoice.action_post()

            if invoice._is_ready_to_be_sent():
                invoice.is_move_sent = True  # Mark invoice as sent
                send_context = {'allow_raising': False, 'allow_fallback_pdf': True}

                default_template_param = (
                    self.env['ir.config_parameter']
                    .sudo()
                    .get_param('sale.default_invoice_email_template', False)
                )

                if default_template_param:
                    mail_template = self.env['mail.template'].sudo().browse(int(default_template_param))
                    if mail_template.exists():
                        send_context['mail_template'] = mail_template

                self.env['account.move.send']._generate_and_send_invoices(invoice, **send_context)
````

### FILE: `odoo_loyalty/upstream/addons/sale_loyalty/models/sale_order_coupon_points.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file27:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/models/sale_order_coupon_points.py; revision 99edb6dd82b7b560930c00b03b694ba700785370; Git blob b0b50be60e2c61182876e6d0fa9a924029ebbcf0"
license: "LGPL-3.0-only"
sha256: "81a49606130beed878aaa93c503970c8e4e9b4645b63f2b097999ccb99c4eeb2"
variables: []
secrets_allowed: false
```

````python
# Part of Odoo. See LICENSE file for full copyright and licensing details.

from odoo import fields, models


class SaleOrderCouponPoints(models.Model):
    _name = 'sale.order.coupon.points'
    _description = "Sale Order Coupon Points - Keeps track of how a sale order impacts a coupon"

    order_id = fields.Many2one(comodel_name='sale.order', ondelete='cascade', required=True, index=True)
    coupon_id = fields.Many2one(comodel_name='loyalty.card', ondelete='cascade', required=True)
    points = fields.Float(required=True)

    _order_coupon_unique = models.Constraint(
        'UNIQUE (order_id, coupon_id)',
        "The coupon points entry already exists.",
    )
````

### FILE: `odoo_loyalty/upstream/addons/sale_loyalty/models/sale_order_line.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file28:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/models/sale_order_line.py; revision 99edb6dd82b7b560930c00b03b694ba700785370; Git blob 0e97738849be4abea95ea107c04c7566d5c58fe3"
license: "LGPL-3.0-only"
sha256: "4aba6673de316b9cb78def3c1d4e31d456ce170524c84d5f8dca93b2a8f18d93"
variables: []
secrets_allowed: false
```

````python
# Part of Odoo. See LICENSE file for full copyright and licensing details.

from odoo import api, fields, models


class SaleOrderLine(models.Model):
    _inherit = 'sale.order.line'

    is_reward_line = fields.Boolean(
        string="Is a program reward line", compute='_compute_is_reward_line'
    )
    reward_id = fields.Many2one(
        comodel_name='loyalty.reward', ondelete='restrict', readonly=True
    )
    coupon_id = fields.Many2one(comodel_name='loyalty.card', ondelete='restrict', readonly=True)
    reward_identifier_code = fields.Char(
        help="Technical field used to link multiple reward lines from the same reward together."
    )
    points_cost = fields.Float(help="How much point this reward costs on the loyalty card.")

    def _compute_name(self):
        # Avoid computing the name for reward lines
        reward = self.filtered('reward_id')
        super(SaleOrderLine, self - reward)._compute_name()

    def _compute_discount(self):
        rewards = self.filtered('reward_id')
        return super(SaleOrderLine, self - rewards)._compute_discount()

    @api.depends('reward_id')
    def _compute_is_reward_line(self):
        for line in self:
            line.is_reward_line = bool(line.reward_id)

    def _compute_tax_ids(self):
        reward_lines = self.filtered('is_reward_line')
        super(SaleOrderLine, self - reward_lines)._compute_tax_ids()
        # Discount reward line is split per tax, the discount is set on the line but not on the product
        # as the product is the generic discount line.
        # In case of a free product, retrieving the tax on the line instead of the product won't affect the behavior.
        for line in reward_lines:
            line = line.with_company(line.company_id)
            fpos = line.order_id.fiscal_position_id or line.order_id.fiscal_position_id._get_fiscal_position(line.order_partner_id)
            # If company_id is set, always filter taxes by the company
            taxes = line.tax_ids.filtered(lambda r: not line.company_id or r.company_id == line.company_id)
            line.tax_ids = fpos.map_tax(taxes)

    def _get_display_price(self):
        # A product created from a promotion does not have a list_price.
        # The price_unit of a reward order line is computed by the promotion, so it can be used directly
        if self.is_reward_line and self.reward_id.reward_type != 'product':
            return self.price_unit
        return super()._get_display_price()

    def _can_be_invoiced_alone(self):
        return super()._can_be_invoiced_alone() and not self.is_reward_line

    def _is_discount_line(self):
        return super()._is_discount_line() or self.reward_id.reward_type == 'discount'

    def _reset_loyalty(self, complete=False):
        """
        Reset the line(s) to a state which does not impact reward computation.
        If complete is set to True we also remove the coupon and reward from the line(s).
            This option should be used when the line will be unlinked.

        Returns self
        """
        vals = {
            'points_cost': 0,
            'price_unit': 0,
            'technical_price_unit': 0,
        }
        if complete:
            vals.update({
                'coupon_id': False,
                'reward_id': False,
            })
        self.write(vals)
        return self

    @api.model_create_multi
    def create(self, vals_list):
        res = super().create(vals_list)
        # Update our coupon points if the order is in a confirmed state
        for line in res:
            if line.coupon_id and line.points_cost and line.state == 'sale':
                line.coupon_id.points -= line.points_cost
                line.order_id._update_loyalty_history(line.coupon_id, line.points_cost)
        return res

    def write(self, vals):
        cost_in_vals = 'points_cost' in vals
        if cost_in_vals:
            previous_vals = {line: (line.points_cost, line.coupon_id) for line in self}
        res = super().write(vals)
        if cost_in_vals:
            # Update our coupon points if the order is in a confirmed state
            for line, (previous_cost, previous_coupon) in previous_vals.items():
                if line.state != 'sale':
                    continue
                if line.points_cost != previous_cost or line.coupon_id != previous_coupon:
                    previous_coupon.points += previous_cost
                    line.coupon_id.points -= line.points_cost
                    # Same coupon: apply the delta in a single call to avoid a redundant search.
                    if line.coupon_id == previous_coupon:
                        line.order_id._update_loyalty_history(line.coupon_id, line.points_cost - previous_cost)
                    else:
                        if previous_coupon:
                            line.order_id._update_loyalty_history(previous_coupon, -previous_cost)
                        if line.coupon_id:
                            line.order_id._update_loyalty_history(line.coupon_id, line.points_cost)
        return res

    def unlink(self):
        # Remove related reward lines
        reward_coupon_set = {(l.reward_id, l.coupon_id, l.reward_identifier_code) for l in self if l.reward_id}
        related_lines = self.env['sale.order.line']
        related_lines |= self.order_id.order_line.filtered(lambda l: (l.reward_id, l.coupon_id, l.reward_identifier_code) in reward_coupon_set)
        # Remove the line's coupon from order if it is the last line using that coupon
        coupons_to_unlink = self.env['loyalty.card']
        for line in self:
            if line.coupon_id:
                # 2 cases:
                #  case 1: coupon has been applied directly
                #  case 2: coupon was created from a program
                if line.coupon_id in line.order_id.applied_coupon_ids:
                    line.order_id.applied_coupon_ids -= line.coupon_id
                elif line.coupon_id.order_id == line.order_id and line.coupon_id.program_id.applies_on == 'current' and\
                    not any(oLine.coupon_id == line.coupon_id and oLine not in related_lines for oLine in line.order_id.order_line):
                    # ondelete='restrict' would prevent deletion of the coupon unlink after unlinking lines
                    coupons_to_unlink |= line.coupon_id
                    line.order_id.code_enabled_rule_ids = line.order_id.code_enabled_rule_ids.filtered(lambda r: r.program_id != line.coupon_id.program_id)
        # Give back the points if the order is confirmed, points are given back if the order is cancelled but in this case we need to do it directly
        for line in related_lines:
            if line.state == 'sale':
                line.coupon_id.points += line.points_cost
        res = super(SaleOrderLine, self | related_lines).unlink()
        coupons_to_unlink.sudo().unlink()
        return res

    def _sellable_lines_domain(self):
        return super()._sellable_lines_domain() + [('reward_id', '=', False)]

    # === TOOLING ===#

    def _can_be_edited_on_portal(self):
        return super()._can_be_edited_on_portal() and not self.is_reward_line
````

### FILE: `odoo_loyalty/upstream/addons/sale_loyalty/tests/common.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file29:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/tests/common.py; revision 99edb6dd82b7b560930c00b03b694ba700785370; Git blob 41c5dca474739038482d3e2f77a25bd88b05e5a7"
license: "LGPL-3.0-only"
sha256: "ab2f2ef3334a82a642178427f24e8316e3ee2d4f91c63dede8712304d3c53db5"
variables: []
secrets_allowed: false
```

````python
# Part of Odoo. See LICENSE file for full copyright and licensing details.

from collections import defaultdict

from odoo.exceptions import ValidationError
from odoo.fields import Command

from odoo.addons.sale.tests.common import SaleCommon


class TestSaleCouponCommon(SaleCommon):

    @classmethod
    def setUpClass(cls):
        super().setUpClass()

        # set currency to not rely on demo data and avoid possible race condition
        cls.currency_ratio = 1.0

        # Set all the existing programs to active=False to avoid interference
        cls.env['loyalty.program'].search([]).sudo().write({'active': False})

        # Taxes
        tax_group_group = cls.env['account.tax.group'].create({
            'name': 'Test Account Tax Group'
        })
        cls.tax_15pc_excl = cls.env['account.tax'].create({
            'name': "Tax 15%",
            'amount_type': 'percent',
            'amount': 15,
            'type_tax_use': 'sale',
            'tax_group_id': tax_group_group.id,
        })

        cls.tax_10pc_incl = cls.env['account.tax'].create({
            'name': "10% Tax incl",
            'amount_type': 'percent',
            'amount': 10,
            'price_include_override': 'tax_included',
            'tax_group_id': tax_group_group.id,
        })

        cls.tax_10pc_base_incl = cls.env['account.tax'].create({
            'name': "10% Tax incl base amount",
            'amount_type': 'percent',
            'amount': 10,
            'price_include_override': 'tax_included',
            'include_base_amount': True,
            'tax_group_id': tax_group_group.id,
        })

        cls.tax_10pc_excl = cls.env['account.tax'].create({
            'name': "10% Tax excl",
            'amount_type': 'percent',
            'amount': 10,
            'price_include_override': 'tax_excluded',
            'tax_group_id': tax_group_group.id,
        })

        cls.tax_20pc_excl = cls.env['account.tax'].create({
            'name': "20% Tax excl",
            'amount_type': 'percent',
            'amount': 20,
            'price_include_override': 'tax_excluded',
            'tax_group_id': tax_group_group.id,
        })

        cls.tax_group = cls.env['account.tax'].create({
            'name': "tax_group",
            'amount_type': 'group',
            'children_tax_ids': [Command.set((cls.tax_10pc_incl + cls.tax_10pc_base_incl).ids)],
            'tax_group_id': tax_group_group.id,
        })

        #products
        cls.product_A = cls.env['product.product'].create({
            'name': 'Product A',
            'list_price': 100,
            'sale_ok': True,
            'taxes_id': [(6, 0, [cls.tax_15pc_excl.id])],
        })

        cls.product_B = cls.env['product.product'].create({
            'name': 'Product B',
            'list_price': 5,
            'sale_ok': True,
            'taxes_id': [(6, 0, [cls.tax_15pc_excl.id])],
        })

        cls.product_C = cls.env['product.product'].create({
            'name': 'Product C',
            'list_price': 100,
            'sale_ok': True,
            'taxes_id': [(6, 0, [])],
        })

        cls.product_D = cls.env['product.product'].create({
            'name': 'Product D',
            'list_price': 100,
            'sale_ok': True,
            'taxes_id': [(6, 0, [cls.tax_group.id])],
        })

        cls.product_gift_card = cls.env['product.product'].create({
            'name': 'Gift Card 50',
            'type': 'service',
            'list_price': 50,
            'sale_ok': True,
            'taxes_id': False,
        })

        # Immediate Program By A + B: get B free
        # No Conditions
        cls.program_gift_card = cls.env['loyalty.program'].create({
            'name': 'Gift Cards',
            'applies_on': 'future',
            'program_type': 'gift_card',
            'trigger': 'auto',
            'rule_ids': [(0, 0, {
                'product_ids': cls.product_gift_card,
                'reward_point_amount': 1,
                'reward_point_mode': 'money',
                'reward_point_split': True,
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount': 1,
                'discount_mode': 'per_point',
                'discount_applicability': 'order',
            })]
        })
        cls.immediate_promotion_program = cls.env['loyalty.program'].create({
            'name': 'Buy A + 1 B, 1 B are free',
            'program_type': 'promotion',
            'applies_on': 'current',
            'company_id': cls.env.company.id,
            'trigger': 'auto',
            'rule_ids': [(0, 0, {
                'product_ids': cls.product_A,
                'reward_point_amount': 1,
                'reward_point_mode': 'order',
                'minimum_qty': 1,
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'product',
                'reward_product_id': cls.product_B.id,
                'reward_product_qty': 1,
                'required_points': 1,
            })],
        })
        cls.code_promotion_program = cls.env['loyalty.program'].create({
            'name': 'Buy 1 A + Enter code, 1 A is free',
            'program_type': 'coupons',
            'trigger': 'with_code',
            'applies_on': 'current',
            'company_id': cls.env.company.id,
            'rule_ids': [(0, 0, {
                'product_ids': cls.product_A,
                'reward_point_amount': 1,
                'reward_point_mode': 'order',
                'minimum_qty': 1,
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'product',
                'reward_product_id': cls.product_A.id,
                'reward_product_qty': 1,
                'required_points': 1,
            })],
        })
        cls.code_promotion_program_with_discount = cls.env['loyalty.program'].create({
            'name': 'Buy 1 C + Enter code, 10 percent discount on C',
            'program_type': 'coupons',
            'trigger': 'with_code',
            'applies_on': 'current',
            'company_id': cls.env.company.id,
            'rule_ids': [(0, 0, {
                'mode': 'with_code',
                'code': 'promotion_code_disc',
                'product_ids': cls.product_C,
                'reward_point_amount': 1,
                'reward_point_mode': 'order',
                'minimum_qty': 1,
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount_mode': 'percent',
                'discount': 10,
                'discount_applicability': 'order',
                'required_points': 1,
            })],
        })

    def _extract_rewards_from_claimable(self, status):
        rewards = self.env['loyalty.reward']
        for info in status.values():
            for reward_count in info['rewards']:
                rewards |= reward_count[0]

    def _apply_promo_code(self, order, code, no_reward_fail=True):
        status = order._try_apply_code(code)
        if 'error' in status:
            raise ValidationError(status['error'])
        if not status and no_reward_fail:
            # Can happen if global discount got filtered out in `_get_claimable_rewards`
            raise ValidationError('No reward to claim with this coupon')
        coupons = self.env['loyalty.card']
        rewards = self.env['loyalty.reward']
        for coupon, coupon_rewards in status.items():
            coupons |= coupon
            rewards |= coupon_rewards
        if len(coupons) == 1 and len(rewards) == 1:
            status = order._apply_program_reward(rewards, coupons)
            if 'error' in status:
                raise ValidationError(status['error'])
        elif len(coupons) == 1 and len(rewards) > 1:
            return rewards

    def _claim_reward(self, order, program, coupon=False):
        if len(program.reward_ids) != 1:
            return False
        coupon = coupon or order.coupon_point_ids.coupon_id.filtered(lambda c: c.program_id == program)
        if len(coupon) != 1:
            return False
        status = order._apply_program_reward(program.reward_ids, coupon)
        return 'error' not in status

    def _auto_rewards(self, order, programs):
        order._update_programs_and_rewards()
        coupons_per_program = defaultdict(lambda: self.env['loyalty.card'])
        for coupon in order.coupon_point_ids.coupon_id:
            coupons_per_program[coupon.program_id] |= coupon
        for program in programs:
            if len(program.reward_ids) > 1 or len(coupons_per_program[program]) != 1 or not program.active:
                continue
            self._claim_reward(order, program, coupons_per_program[program])

    def _generate_coupons(self, loyality_program, coupon_qty=1):
        self.env['loyalty.generate.wizard'].with_context(active_id=loyality_program.id).create({
            'coupon_qty': coupon_qty,
        }).generate_coupons()
        return loyality_program.coupon_ids

class TestSaleCouponNumbersCommon(TestSaleCouponCommon):
    @classmethod
    def setUpClass(cls):
        super().setUpClass()

        cls.largeCabinet = cls.env['product.product'].create({
            'name': 'Large Cabinet',
            'list_price': 320.0,
            'taxes_id': False,
        })
        cls.conferenceChair = cls.env['product.product'].create({
            'name': 'Conference Chair',
            'list_price': 16.5,
            'taxes_id': False,
        })
        cls.pedalBin = cls.env['product.product'].create({
            'name': 'Pedal Bin',
            'list_price': 47.0,
            'taxes_id': False,
        })
        cls.drawerBlack = cls.env['product.product'].create({
            'name': 'Drawer Black',
            'list_price': 25.0,
            'taxes_id': False,
        })
        cls.largeMeetingTable = cls.env['product.product'].create({
            'name': 'Large Meeting Table',
            'list_price': 40000.0,
            'taxes_id': False,
        })

        cls.p1 = cls.env['loyalty.program'].create({
            'name': 'Code for 10% on orders',
            'trigger': 'with_code',
            'program_type': 'promotion',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {
                'mode': 'with_code',
                'code': 'test_10pc',
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount_mode': 'percent',
                'discount': 10,
                'discount_applicability': 'order',
                'required_points': 1,
            })],
        })
        cls.p2 = cls.env['loyalty.program'].create({
            'name': 'Buy 3 cabinets, get one for free',
            'trigger': 'auto',
            'program_type': 'promotion',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {
                'product_ids': cls.largeCabinet,
                'reward_point_mode': 'unit',
                'minimum_qty': 3,
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'product',
                'reward_product_id': cls.largeCabinet.id,
                'reward_product_qty': 1,
                'required_points': 3,
            })],
        })
        cls.p3 = cls.env['loyalty.program'].create({
            'name': 'Buy 1 drawer black, get a free Large Meeting Table',
            'trigger': 'auto',
            'program_type': 'promotion',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {
                'product_ids': cls.drawerBlack,
                'reward_point_mode': 'order',
                'minimum_qty': 1,
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'product',
                'reward_product_id': cls.largeMeetingTable.id,
                'reward_product_qty': 1,
                'required_points': 1,
            })],
        })
        cls.discount_coupon_program = cls.env['loyalty.program'].create({
            'name': '$100 coupon',
            'program_type': 'coupons',
            'trigger': 'with_code',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {
                'minimum_amount': 100,
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount_mode': 'per_point',
                'discount': 100,
                'discount_applicability': 'order',
                'required_points': 1,
            })],
        })
        cls.all_programs = cls.env['loyalty.program'].search([])
````

### FILE: `odoo_loyalty/upstream/addons/sale_loyalty/tests/test_buy_gift_card.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file30:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/tests/test_buy_gift_card.py; revision 99edb6dd82b7b560930c00b03b694ba700785370; Git blob 5818de14efc9fde3101f0714cb446c13692d3bd3"
license: "LGPL-3.0-only"
sha256: "f724115af4e50776f152bb933b01999bef7e800aec4eea0c43ad38f9916bb3e8"
variables: []
secrets_allowed: false
```

````python
# Part of Odoo. See LICENSE file for full copyright and licensing details.

from odoo import Command
from odoo.tests.common import tagged

from odoo.addons.sale_loyalty.tests.common import TestSaleCouponCommon


@tagged('-at_install', 'post_install')
class TestBuyGiftCard(TestSaleCouponCommon):

    def test_buying_gift_card(self):
        order = self.empty_order
        self.immediate_promotion_program.active = False
        order.write({'order_line': [
            (0, False, {
                'product_id': self.product_A.id,
                'name': 'Ordinary Product A',
                'product_uom_qty': 1.0,
            }),
            (0, False, {
                'product_id': self.product_gift_card.id,
                'name': 'Gift Card Product',
                'product_uom_qty': 1.0,
            })
        ]})
        self.assertEqual(len(order.order_line.ids), 2)
        self.assertEqual(len(order._get_reward_coupons()), 0)
        order._update_programs_and_rewards()
        self.assertEqual(len(order._get_reward_coupons()), 1)
        order.order_line[1].product_uom_qty = 2
        order._update_programs_and_rewards()
        self.assertEqual(len(order._get_reward_coupons()), 2)
        order.order_line[1].product_uom_qty = 1
        order._update_programs_and_rewards()
        self.assertEqual(len(order._get_reward_coupons()), 1)

    def test_gift_card_email_sender(self):
        """Ensure that sending gift card emails have a sender.
        Either the order's salesman if available, otherwise the order's company.
        """
        mail_template = self.env['mail.template'].create({
            'name': "Gift Card Mail",
            'model_id': self.env.ref('loyalty.model_loyalty_card').id,
            'auto_delete': False,
        })
        self.program_gift_card.communication_plan_ids = [Command.create({
            'trigger': 'create',
            'mail_template_id': mail_template.id,
        })]
        order = self.empty_order
        salesman = order.user_id.partner_id.ensure_one()
        salesman.email = "sales@company.co"
        company = order.company_id.partner_id
        company.email = "noreply@company.co"
        order.write({
            'order_line': [Command.create({'product_id': self.product_gift_card.id})],
        })
        order._update_programs_and_rewards()

        # Create an order without salesman to test company-based fallback
        orders = order + order.copy({'user_id': None})

        # Clear out the mailbox before sending mail
        self.env['mail.mail'].search([]).sudo().unlink()

        # Confirm order as Public User to trigger loyalty mail
        public_user = self.env.ref('base.public_user')
        orders.with_user(public_user).with_company(order.company_id).sudo().action_confirm()

        mails = self.env['mail.mail'].search([])
        self.assertEqual(len(mails), 2)
        salesman_mail = mails.filtered(lambda m: m.author_id == salesman).ensure_one()
        company_mail = mails.filtered(lambda m: m.author_id == company).ensure_one()
        self.assertEqual(salesman_mail.email_from, salesman.email_formatted)
        self.assertEqual(company_mail.email_from, company.email_formatted)
````

### FILE: `odoo_loyalty/upstream/addons/sale_loyalty/tests/test_loyalty.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file31:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/tests/test_loyalty.py; revision 99edb6dd82b7b560930c00b03b694ba700785370; Git blob 0eb9ddcd8e035bc61b355744ae66db1040bfe127"
license: "LGPL-3.0-only"
sha256: "b2e97e780f49b8303b1930cdc7b6687338bcdf0cd48039e747144eddd14db47e"
variables: []
secrets_allowed: false
```

````python
# Part of Odoo. See LICENSE file for full copyright and licensing details.

from freezegun import freeze_time

from odoo.exceptions import UserError, ValidationError
from odoo.fields import Command
from odoo.tests import new_test_user, tagged
from odoo.tools.float_utils import float_compare

from odoo.addons.sale_loyalty.tests.common import TestSaleCouponCommon


@tagged('post_install', '-at_install')
class TestLoyalty(TestSaleCouponCommon):

    @classmethod
    def setUpClass(cls):
        super().setUpClass()
        cls.env['loyalty.program'].search([]).write({'active': False})

        cls.product_a, cls.product_b = cls.env['product.product'].create([
            {
                'name': 'Product C',
                'list_price': 100,
                'sale_ok': True,
                'taxes_id': [Command.set([])],
            },
            {
                'name': "Product B",
                'sale_ok': True,
            }
        ])

        cls.ewallet_program = cls.env['loyalty.program'].create({
            'name': 'eWallet Program',
            'program_type': 'ewallet',
            'trigger': 'auto',
            'applies_on': 'future',
            'reward_ids': [Command.create({
                'reward_type': 'discount',
                'discount_mode': 'per_point',
                'discount': 1,
            })],
            'rule_ids': [Command.create({
                'reward_point_amount': '1',
                'reward_point_mode': 'money',
                'product_ids': cls.env.ref('loyalty.ewallet_product_50'),
            })],
            'trigger_product_ids': cls.env.ref('loyalty.ewallet_product_50'),
        })

        cls.ewallet = cls.env['loyalty.card'].create({
            'program_id': cls.ewallet_program.id,
            'partner_id': cls.partner.id,
            'points': 10,
        })
        cls.ewallet_program.coupon_ids = [Command.set([cls.ewallet.id])]

        cls.user_salemanager = new_test_user(cls.env, login='user_salemanager', groups='sales_team.group_sale_manager')

        cls.promotion_code_10pc = cls.env['loyalty.program'].create({
            'name': "Code for 10% on orders",
            'trigger': 'with_code',
            'program_type': 'promotion',
            'applies_on': 'current',
            'rule_ids': [Command.create({
                'mode': 'with_code',
                'code': 'test_10pc',
            })],
            'reward_ids': [Command.create({
                'reward_type': 'discount',
                'discount_mode': 'percent',
                'discount': 10,
                'discount_applicability': 'order',
                'required_points': 1,
            })],
        })

    def test_nominative_programs(self):
        loyalty_program = self.env['loyalty.program'].create({
            'name': 'Loyalty Program',
            'program_type': 'loyalty',
            'trigger': 'auto',
            'applies_on': 'both',
            'rule_ids': [(0, 0, {
                'reward_point_mode': 'unit',
                'reward_point_amount': 1,
                'product_ids': [self.product_a.id],
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount': 1.5,
                'discount_mode': 'per_point',
                'discount_applicability': 'order',
                'required_points': 3,
            })],
        })

        order = self.empty_order
        order._update_programs_and_rewards()
        claimable_rewards = order._get_claimable_rewards()
        # Should be empty since we do not have any coupon created yet
        self.assertFalse(claimable_rewards, "No program should be applicable")
        loyalty_card = self.env['loyalty.card'].create({
            'program_id': loyalty_program.id,
            'partner_id': self.partner.id,
            'points': 10,
        })
        self.ewallet.points = 0

        order.write({
            'order_line': [(0, 0, {
                'product_id': self.product_a.id,
                'product_uom_qty': 1,
            })]
        })
        order._update_programs_and_rewards()
        claimable_rewards = order._get_claimable_rewards()
        self.assertEqual(len(claimable_rewards), 1, "The ewallet program should not be applicable since the card has no points.")
        vals = order._get_reward_values_discount(loyalty_program.reward_ids[0], loyalty_card)
        self.assertEqual(
            vals[0]['points_cost'] % loyalty_program.reward_ids.required_points,
            0,
            "Can only use a whole number of required points",
        )
        self.assertEqual(vals[0]['points_cost'], 9, "Use maximum available points for the reward")
        self.ewallet.points = 50
        order._update_programs_and_rewards()
        claimable_rewards = order._get_claimable_rewards()
        self.assertEqual(len(claimable_rewards), 2, "Now that the ewallet has some points they should both be applicable.")

    def test_cancel_order_with_coupons(self):
        """This test ensure that creating an order with coupons will not
        raise an access error on POS line modele when canceling the order."""

        self.env['loyalty.program'].create({
            'name': '10% Discount',
            'program_type': 'coupons',
            'applies_on': 'current',
            'trigger': 'auto',
            'rule_ids': [(0, 0, {})],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount': 10,
                'discount_mode': 'percent',
                'discount_applicability': 'order',
            })]
        })

        order = self.env['sale.order'].with_user(self.user_salemanager).create({
            'partner_id': self.partner.id,
            'order_line': [
                (0, 0, {
                    'product_id': self.product_a.id,
                })
            ]
        })

        order._update_programs_and_rewards()
        self.assertTrue(order.coupon_point_ids)

        # Canceling the order should not raise an access error:
        # During the cancel process, we are trying to get `use_count` of the coupon,
        # and we call the `_compute_use_count` that is also in pos_loyalty.
        # This last one will try to find related POS lines while user have not access to POS.
        order._action_cancel()
        self.assertFalse(order.coupon_point_ids)

    def test_salesperson_can_cancel_order_with_coupons(self):
        """Test that a salesperson can cancel an order with coupon points without access error."""
        user_salesman = new_test_user(
            self.env, login="user_salesman", groups="sales_team.group_sale_salesman"
        )
        self.env["loyalty.program"].create({
            "name": "10% Discount",
            "program_type": "coupons",
            "trigger": "auto",
            "reward_ids": [Command.create({"reward_type": "discount", "discount": 10})],
        })
        order = (
            self
            .env["sale.order"]
            .with_user(user_salesman)
            .create({
                "partner_id": self.partner.id,
                "order_line": [Command.create({"product_id": self.product_a.id})],
            })
        )
        order.action_confirm()
        self.assertTrue(order.coupon_point_ids)
        order._action_cancel()
        self.assertFalse(order.coupon_point_ids)

    def test_distribution_amount_payment_programs(self):
        """
        Check how the amount of a payment reward is distributed.
        An ewallet should not be used to refund taxes.
        Its amount must be distributed between the products.
        """

        # Create two products
        product_a, product_b = self.env['product.product'].create([
            {
                'name': 'Product A',
                'list_price': 100,
                'sale_ok': True,
                'taxes_id': [Command.set(self.tax_15pc_excl.ids)],
            },
            {
                'name': 'Product B',
                'list_price': 100,
                'sale_ok': True,
                'taxes_id': [Command.set(self.tax_15pc_excl.ids)],
            },
        ])

        # Create a coupon and a ewallet
        coupon_program, ewallet_program = self.env['loyalty.program'].create([
            {
                'name': 'Coupon Program',
                'program_type': 'coupons',
                'trigger': 'with_code',
                'applies_on': 'both',
                'reward_ids': [Command.create({
                        'reward_type': 'discount',
                        'discount': 100.0,
                        'discount_applicability': 'specific',
                        'discount_product_domain': '[("name", "=", "Product A")]',
                })],
            },
            {
                'name': 'eWallet Program',
                'program_type': 'ewallet',
                'applies_on': 'future',
                'trigger': 'auto',
                'rule_ids': [Command.create({
                    'reward_point_mode': 'money',
                })],
                'reward_ids': [Command.create({
                    'discount_mode': 'per_point',
                    'discount': 1,
                    'discount_applicability': 'order',
                })],
            }
        ])

        coupon_partner, _ = self.env['loyalty.card'].create([
            {
                'program_id': coupon_program.id,
                'partner_id': self.partner.id,
                'points': 1,
                'code': '5555',
            },
            {
                'program_id': ewallet_program.id,
                'partner_id': self.partner.id,
                'points': 115,
            },
        ])

        # Create the order
        order = self.env['sale.order'].with_user(self.user_salemanager).create({
            'partner_id': self.partner.id,
            'order_line': [
                    Command.create({
                        'product_id': product_a.id,
                    }),
                    Command.create({
                        'product_id': product_b.id,
                    }),
            ]
        })

        self.assertEqual(order.amount_total, 230.0)
        self.assertEqual(order.amount_untaxed, 200.0)
        self.assertEqual(order.amount_tax, 30.0)

        # Apply the eWallet
        order._update_programs_and_rewards()
        self._claim_reward(order, ewallet_program)

        self.assertEqual(order.amount_total, 115.0)
        self.assertEqual(order.amount_untaxed, 85.0)
        self.assertEqual(order.amount_tax, 30.0)
        self.assertEqual(order.reward_amount, -115.0)

        # Apply the coupon
        self._apply_promo_code(order, coupon_partner.code)

        self.assertEqual(order.amount_total, 0.0)
        self.assertEqual(order.amount_untaxed, -15.0)
        self.assertEqual(order.amount_tax, 15.0)
        self.assertEqual(order.reward_amount, -215.0)

    def test_discount_max_amount_on_specific_product(self):
        product_a = self.product_A
        product_b = self.product_B
        product_a.write({'taxes_id': [Command.set(self.tax_20pc_excl.ids)]})
        product_b.write({'list_price': -20, 'taxes_id': [Command.set(self.tax_20pc_excl.ids)]})

        self.env['loyalty.program'].search([]).write({'active': False})
        promotion = self.env['loyalty.program'].create({
            'name': '10% Discount',
            'program_type': 'promotion',
            'trigger': 'auto',
            'rule_ids': [Command.create({'reward_point_amount': 1, 'reward_point_mode': 'unit'})],
            'reward_ids': [Command.create({
                'discount': 10.0,
                'discount_max_amount': 9,
                'discount_applicability': 'specific',
                'discount_product_ids': [product_a.id],
            })],
        })

        order = self.env['sale.order'].create({
            'partner_id': self.partner.id,
            'order_line': [Command.create({'product_id': product_a.id})],
        })
        self.assertEqual(order.reward_amount, 0)

        self._auto_rewards(order, promotion)
        reward_amount_tax_included = sum(l.price_total for l in order.order_line if l.reward_id)
        msg = "Max discount amount reached, the reward amount should be the max amount value."
        self.assertEqual(reward_amount_tax_included, -9, msg)

        order.order_line = [Command.clear(), Command.create({'product_id': product_b.id})]
        self._auto_rewards(order, promotion)
        reward_amount_tax_included = sum(l.price_total for l in order.order_line if l.reward_id)
        msg = "This product is not eligible to the discount."
        self.assertEqual(reward_amount_tax_included, 0, msg=msg)

        order.order_line = [
            Command.clear(),
            Command.create({'product_id': product_a.id}),  # price_total = 120
            Command.create({'product_id': product_b.id}),  # price_total = -20
        ]
        self._auto_rewards(order, promotion)
        reward_amount_tax_included = sum(l.price_total for l in order.order_line if l.reward_id)
        msg = "Reward amount above the max amount, the reward should be the max amount value."
        self.assertEqual(reward_amount_tax_included, -9, msg)

        order.order_line = [
            Command.clear(),
            Command.create({'product_id': product_a.id}),                     # price_total = 120
            Command.create({'product_id': product_b.id, 'price_unit': -95}),  # price_total = -114
        ]
        self._auto_rewards(order, promotion)
        reward_amount_tax_included = sum(l.price_total for l in order.order_line if l.reward_id)
        msg = "Reward amount should never surpass the order's current total amount."
        self.assertEqual(reward_amount_tax_included, -6, msg)

        order.order_line = [
            Command.clear(),
            Command.create({'product_id': product_a.id, 'price_unit': 50}),  # price_total = 60
            Command.create({'product_id': product_b.id, 'price_unit': -5}),  # price_total = -6
        ]
        self._auto_rewards(order, promotion)
        reward_amount_tax_included = sum(l.price_total for l in order.order_line if l.reward_id)
        msg = "Reward amount should be the percentage one if under the max amount discount."
        self.assertEqual(reward_amount_tax_included, -6, msg)

    def test_points_awarded_global_discount_code_no_domain_program(self):
        """
        Check the calculation for points awarded when there is a global discount applied and the
        loyalty program applies on all products (no domain).
        """
        LoyaltyProgram = self.env['loyalty.program']
        loyalty_program = LoyaltyProgram.create(LoyaltyProgram._get_template_values()['loyalty'])

        loyalty_card = self.env['loyalty.card'].create({
            'program_id': loyalty_program.id,
            'partner_id': self.partner.id,
            'points': 0,
        })

        order = self.env['sale.order'].with_user(self.user_salemanager).create({
            'partner_id': self.partner.id,
            'order_line': [
                Command.create({
                    'product_id': self.product_A.id,
                    'tax_ids': False,
                }),
            ]
        })

        promotion_program = self.env['loyalty.program'].create([{
            'name': "Coupon Program",
            'program_type': 'promotion',
            'trigger': 'auto',
            'applies_on': 'current',
            'rule_ids': [Command.create({
                    'reward_point_amount': 1,
                    'reward_point_mode': 'order',
                    'minimum_amount': 10,
                })],
            'reward_ids': [Command.create({
                    'reward_type': 'discount',
                    'discount': 10.0,
                    'discount_applicability': 'order',
                    'required_points': 1,
                })],
        }])

        self.assertEqual(order.amount_total, 100)
        self._auto_rewards(order, promotion_program)
        self.assertEqual(order.amount_total, 90)
        order.action_confirm()
        self.assertEqual(loyalty_card.points, 90)

    def test_multiple_rewards_after_confirm(self):
        """
        Check that multiple rewards from a loyalty promotion program are correctly applied to a SO
        after its confirmation by asserting that:
            - Both rewards are applied to the order lines.
            - The total points cost matches the rule's requirement.
            - The coupon's points are fully consumed after applying the rewards.
        """
        promo_program = self.env['loyalty.program'].create({
            'name': 'Multiple Rewards Promotion',
            'program_type': 'promotion',
            'applies_on': 'current',
            'company_id': self.env.company.id,
            'trigger': 'auto',
            'rule_ids': [
                Command.create({
                    'product_ids': self.product_A,
                    'reward_point_amount': 1,
                    'reward_point_mode': 'order',
                    'minimum_qty': 1,
                }),
            ],
            'reward_ids': [
                Command.create({
                    'discount': 10,
                    'discount_applicability': 'specific',
                    'discount_product_ids': [self.product_A.id],
                    'required_points': 0.5,
                }),
                Command.create({
                    'discount': 15,
                    'discount_applicability': 'specific',
                    'discount_product_ids': [self.product_B.id],
                    'required_points': 0.5,
                }),
            ],
        })

        order = self.empty_order
        order.order_line = [
            Command.create({'product_id': self.product_A.id, 'product_uom_qty': 1}),
            Command.create({'product_id': self.product_B.id, 'product_uom_qty': 1}),
        ]
        order.action_confirm()

        order._update_programs_and_rewards()
        coupon = order.coupon_point_ids.coupon_id.filtered(lambda c: c.program_id == promo_program)
        reward1, reward2 = rewards = promo_program.reward_ids
        order._apply_program_reward(reward1, coupon)
        order._apply_program_reward(reward2, coupon)

        self.assertEqual(order.order_line.reward_id, rewards, "All rewards should be applied")
        self.assertEqual(sum(order.order_line.mapped('points_cost')), 1)
        self.assertEqual(coupon.points, 0)

    def test_points_awarded_discount_code_no_domain_program(self):
        """
        Check the calculation for points awarded when there is a discount coupon applied and the
        loyalty program applies on all products (no domain).
        """
        LoyaltyProgram = self.env['loyalty.program']
        loyalty_program = LoyaltyProgram.create(LoyaltyProgram._get_template_values()['loyalty'])
        loyalty_card = self.env['loyalty.card'].create({
            'program_id': loyalty_program.id,
            'partner_id': self.partner.id,
            'points': 0,
        })

        order = self.env['sale.order'].with_user(self.user_salemanager).create({
            'partner_id': self.partner.id,
            'order_line': [
                Command.create({
                    'product_id': self.product_A.id,
                    'tax_ids': False,
                }),
            ]
        })

        self.assertEqual(order.amount_total, 100)
        self._apply_promo_code(order, "test_10pc")
        self.assertEqual(order.amount_total, 90)
        order.action_confirm()
        self.assertEqual(loyalty_card.points, 90)

    def test_points_awarded_general_discount_code_specific_domain_program(self):
        """
        Check the calculation for points awarded when there is a discount coupon applied and the
        loyalty program applies on a specific domain. The discount code has no domain. The product
        related to that discount is not in the domain of the loyalty program.
        Expected behavior: The discount is not included in the computation of points
        """
        product_category_food = self.env['product.category'].create({
            'name': "Food",
        })

        self.product_A.categ_id = product_category_food
        self.product_B.list_price = 50

        LoyaltyProgram = self.env['loyalty.program']
        loyalty_program = LoyaltyProgram.create(LoyaltyProgram._get_template_values()['loyalty'])
        loyalty_program.rule_ids.product_category_id = product_category_food.id
        loyalty_card = self.env['loyalty.card'].create({
            'program_id': loyalty_program.id,
            'partner_id': self.partner.id,
            'points': 0,
        })

        order = self.env['sale.order'].with_user(self.user_salemanager).create({
            'partner_id': self.partner.id,
            'order_line': [
                Command.create({
                    'product_id': self.product_A.id,
                    'tax_ids': False,
                }),
                Command.create({
                    'product_id': self.product_B.id,
                    'tax_ids': False,
                }),
            ]
        })

        self.assertEqual(order.amount_total, 150)
        self._apply_promo_code(order, "test_10pc")
        self.assertEqual(order.amount_total, 135)  # (product_A + product_B) * 0.9
        order.action_confirm()
        self.assertEqual(loyalty_card.points, 100)

    def test_points_awarded_specific_discount_code_specific_domain_program(self):
        """
        Check the calculation for points awarded when there is a discount coupon applied and the
        loyalty program applies on a specific domain. The discount code has the same domain as the
        loyalty program. The product related to that discount code is set up to be included in the
        domain of the loyalty program.
        Expected behavior: The discount is included in the computation of points
        """
        product_category_food = self.env['product.category'].create({
            'name': "Food",
        })

        self.product_A.categ_id = product_category_food
        self.product_B.list_price = 50

        LoyaltyProgram = self.env['loyalty.program']
        loyalty_program = LoyaltyProgram.create(LoyaltyProgram._get_template_values()['loyalty'])
        loyalty_program.rule_ids.product_category_id = product_category_food.id
        loyalty_card = self.env['loyalty.card'].create({
            'program_id': loyalty_program.id,
            'partner_id': self.partner.id,
            'points': 0,
        })

        self.promotion_code_10pc.rule_ids.product_category_id = product_category_food.id
        self.promotion_code_10pc.reward_ids.discount_applicability = 'specific'
        self.promotion_code_10pc.reward_ids.discount_product_category_id = product_category_food.id

        discount_product = self.env['product.product'].search([('id', '=', self.promotion_code_10pc.reward_ids.discount_line_product_id.id)])
        discount_product.categ_id = product_category_food.id

        order = self.env['sale.order'].with_user(self.user_salemanager).create({
            'partner_id': self.partner.id,
            'order_line': [
                Command.create({
                    'product_id': self.product_A.id,
                    'tax_ids': False,
                }),
                Command.create({
                    'product_id': self.product_B.id,
                    'tax_ids': False,
                }),
            ]
        })

        self.assertEqual(order.amount_total, 150)
        self._apply_promo_code(order, "test_10pc")
        self.assertEqual(order.amount_total, 140)  # (product_A * 0.9 ) + product_B
        order.action_confirm()
        self.assertEqual(loyalty_card.points, 90)

    def test_points_awarded_ewallet(self):
        """
        Check the calculation for point awarded when using ewallet
        """
        LoyaltyProgram = self.env['loyalty.program']
        loyalty_program = LoyaltyProgram.create(LoyaltyProgram._get_template_values()['loyalty'])
        loyalty_card = self.env['loyalty.card'].create({
            'program_id': loyalty_program.id,
            'partner_id': self.partner.id,
            'points': 0,
        })
        order = self.env['sale.order'].with_user(self.user_salemanager).create({
            'partner_id': self.partner.id,
            'order_line': [
                Command.create({
                    'product_id': self.product_A.id,
                    'tax_ids': False,
                }),
            ]
        })

        self.assertEqual(order.amount_total, 100)
        order._update_programs_and_rewards()
        self._claim_reward(order, self.ewallet_program, coupon=self.ewallet)
        self.assertEqual(order.amount_total, 90)
        order.action_confirm()
        self.assertEqual(loyalty_card.points, 100)

    def test_points_awarded_giftcard(self):
        """
        Check the calculation for point awarded when using a gift card
        """
        LoyaltyProgram = self.env['loyalty.program']
        loyalty_program = LoyaltyProgram.create(LoyaltyProgram._get_template_values()['loyalty'])
        loyalty_card = self.env['loyalty.card'].create({
            'program_id': loyalty_program.id,
            'partner_id': self.partner.id,
            'points': 0,
        })

        program_gift_card = self.env['loyalty.program'].create({
            'name': "Gift Cards",
            'applies_on': 'future',
            'program_type': 'gift_card',
            'trigger': 'auto',
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount': 1,
                'discount_mode': 'per_point',
                'discount_applicability': 'order',
            })]
        })

        self.env['loyalty.generate.wizard'].with_context(active_id=program_gift_card.id).create({
            'coupon_qty': 1,
            'points_granted': 50,
        }).generate_coupons()
        gift_card = program_gift_card.coupon_ids[0]

        order = self.env['sale.order'].with_user(self.user_salemanager).create({
            'partner_id': self.partner.id,
            'order_line': [
                Command.create({
                    'product_id': self.product_A.id,
                    'tax_ids': False,
                }),
            ]
        })

        self.assertEqual(order.amount_total, 100)
        self._apply_promo_code(order, gift_card.code)
        self.assertEqual(order.amount_total, 50)
        order.action_confirm()
        self.assertEqual(loyalty_card.points, 100)

    def test_multiple_discount_specific(self):
        """
        Check the discount calculation if it is based on the remaining amount
        """

        product_A = self.env['product.product'].create({
            'name': 'Product A',
            'list_price': 100,
            'sale_ok': True,
            'taxes_id': [],
        })

        coupon_program = self.env['loyalty.program'].create([{
            'name': 'Coupon Program',
            'program_type': 'promotion',
            'trigger': 'auto',
            'applies_on': 'current',
            'rule_ids': [Command.create({
                    'reward_point_amount': 1,
                    'reward_point_mode': 'unit',
                })],
            'reward_ids': [Command.create({
                    'reward_type': 'discount',
                    'discount': 10.0,
                    'discount_applicability': 'specific',
                    'required_points': 1,
                })],
        }])

        order = self.env['sale.order'].with_user(self.user_salemanager).create({
            'partner_id': self.partner.id,
            'order_line': [Command.create({
                    'product_id': product_A.id,
                    'product_uom_qty': 3,
                })]
        })

        self.assertEqual(float_compare(order.amount_total, 300, precision_rounding=3), 0)

        order._update_programs_and_rewards()
        self._claim_reward(order, coupon_program)
        self.assertEqual(float_compare(order.amount_total, 270, precision_rounding=3), 0, "300 * 0.9 = 270")

        order._update_programs_and_rewards()
        self._claim_reward(order, coupon_program)
        self.assertEqual(float_compare(order.amount_total, 243, precision_rounding=3), 0, "300 * 0.9 * 0.9 = 243")

        order._update_programs_and_rewards()
        self._claim_reward(order, coupon_program)
        self.assertEqual(float_compare(order.amount_total, 218.7, precision_rounding=3), 0, "300 * 0.9 * 0.9 * 0.9 = 218.7")

    def test_promotion_program_restricted_to_pricelists(self):
        self.env['product.pricelist'].search([]).action_archive()
        company_currency = self.env.company.currency_id
        pricelist_1, pricelist_2 = self.env['product.pricelist'].create([
            {'name': 'Basic company_currency pricelist', 'currency_id': company_currency.id},
            {'name': 'Other company_currency pricelist', 'currency_id': company_currency.id},
        ])
        self.immediate_promotion_program.active = True
        order = self.empty_order.copy()
        order.write({'order_line': [
            (0, False, {
                'product_id': self.product_A.id,
                'name': '1 Product A',
                'product_uom_qty': 1.0,
            }),
            (0, False, {
                'product_id': self.product_B.id,
                'name': '2 Product B',
                'product_uom_qty': 1.0,
            }),
        ]})

        applied_message = "The promo offer should have been applied."
        not_applied_message = "The promo offer should not have been applied because the order's " \
                              "pricelist is not eligible to this promotion."

        order.pricelist_id = self.env['product.pricelist']
        order._update_programs_and_rewards()
        self._claim_reward(order, self.immediate_promotion_program)
        self.assertEqual(len(order.order_line.ids), 3, applied_message)

        order.pricelist_id = pricelist_1
        order._update_programs_and_rewards()
        self._claim_reward(order, self.immediate_promotion_program)
        self.assertEqual(len(order.order_line.ids), 3, applied_message)

        self.immediate_promotion_program.pricelist_ids = [pricelist_1.id]
        order.pricelist_id = self.env['product.pricelist']
        order._update_programs_and_rewards()
        self._claim_reward(order, self.immediate_promotion_program)
        self.assertEqual(len(order.order_line.ids), 2, not_applied_message)

        order.pricelist_id = pricelist_1
        order._update_programs_and_rewards()
        self._claim_reward(order, self.immediate_promotion_program)
        self.assertEqual(len(order.order_line.ids), 3, applied_message)

        order.pricelist_id = pricelist_2
        order._update_programs_and_rewards()
        self._claim_reward(order, self.immediate_promotion_program)
        self.assertEqual(len(order.order_line.ids), 2, not_applied_message)

        self.immediate_promotion_program.pricelist_ids = [pricelist_1.id, pricelist_2.id]
        order.pricelist_id = self.env['product.pricelist']
        order._update_programs_and_rewards()
        self._claim_reward(order, self.immediate_promotion_program)
        self.assertEqual(len(order.order_line.ids), 2, not_applied_message)

        order.pricelist_id = pricelist_1
        order._update_programs_and_rewards()
        self._claim_reward(order, self.immediate_promotion_program)
        self.assertEqual(len(order.order_line.ids), 3, applied_message)

    def test_coupon_program_restricted_to_pricelists(self):
        self.env['product.pricelist'].search([]).action_archive()
        company_currency = self.env.company.currency_id
        pricelist_1, pricelist_2 = self.env['product.pricelist'].create([
            {'name': 'Basic company_currency pricelist', 'currency_id': company_currency.id},
            {'name': 'Other company_currency pricelist', 'currency_id': company_currency.id},
        ])

        self.code_promotion_program.active = True
        self.env['loyalty.generate.wizard'].with_context(
            active_id=self.code_promotion_program.id
        ).create({'coupon_qty': 7, 'points_granted': 1}).generate_coupons()
        coupons = self.code_promotion_program.coupon_ids

        order_no_pricelist = self.empty_order.copy()
        order_no_pricelist.write({'pricelist_id': None, 'order_line': [
            (0, False, {
                'product_id': self.product_A.id,
                'name': '1 Product A',
                'product_uom_qty': 1.0,
            }),
        ]})
        order_pricelist_1 = order_no_pricelist.copy()
        order_pricelist_1.pricelist_id = pricelist_1
        order_pricelist_2 = order_no_pricelist.copy()
        order_pricelist_2.pricelist_id = pricelist_2

        applied_message = "The coupon code should have been applied."
        not_applied_message = "The coupon code should not have been applied because the order's " \
                              "pricelist is not eligible to this promotion."

        order_0 = order_no_pricelist.copy()
        self._apply_promo_code(order_0, coupons[0].code)
        self.assertEqual(len(order_0.order_line.ids), 2, applied_message)

        order_1 = order_pricelist_1.copy()
        self._apply_promo_code(order_1, coupons[1].code)
        self.assertEqual(len(order_1.order_line.ids), 2, applied_message)

        self.code_promotion_program.pricelist_ids = [pricelist_1.id]
        order_2 = order_no_pricelist.copy()
        with self.assertRaises(ValidationError):
            self._apply_promo_code(order_2, coupons[2].code)
        self.assertEqual(len(order_2.order_line.ids), 1, not_applied_message)

        order_3 = order_pricelist_1.copy()
        self._apply_promo_code(order_3, coupons[3].code)
        self.assertEqual(len(order_3.order_line.ids), 2, applied_message)

        order_4 = order_pricelist_2.copy()
        with self.assertRaises(ValidationError):
            self._apply_promo_code(order_4, coupons[4].code)
        self.assertEqual(len(order_4.order_line.ids), 1, not_applied_message)

        self.code_promotion_program.pricelist_ids = [pricelist_1.id, pricelist_2.id]
        order_5 = order_no_pricelist.copy()
        with self.assertRaises(ValidationError):
            self._apply_promo_code(order_5, coupons[5].code)
        self.assertEqual(len(order_5.order_line.ids), 1, not_applied_message)

        order_6 = order_pricelist_1.copy()
        self._apply_promo_code(order_6, coupons[6].code)
        self.assertEqual(len(order_6.order_line.ids), 2, applied_message)

    def test_specific_promotion_on_free_product(self):

        product_A = self.env['product.product'].create({
            'name': 'Product A',
            'list_price': 100,
            'sale_ok': True,
            'taxes_id': [],
        })

        promotion_program = self.env['loyalty.program'].create([{
            'name': 'Promotion Program',
            'program_type': 'promotion',
            'trigger': 'auto',
            'applies_on': 'current',
            'rule_ids': [Command.create({
                'reward_point_amount': 1,
                'reward_point_mode': 'unit',
            })],
            'reward_ids': [Command.create({
                'reward_type': 'discount',
                'discount': 10.0,
                'discount_applicability': 'specific',
                'discount_product_ids': [product_A.id],
                'required_points': 1,
            })],
        }])

        order = self.env['sale.order'].with_user(self.user_salemanager).create({
            'partner_id': self.partner.id,
            'order_line': [
                Command.create({
                    'product_id': product_A.id,
                }),
                Command.create({
                    'product_id': product_A.id,
                    'discount': 100,
                }),
            ]
        })

        order._update_programs_and_rewards()
        self._claim_reward(order, promotion_program)
        self.assertEqual(order.amount_total, 90)

    def test_gift_card_program_without_product(self):
        product_A = self.env['product.product'].create({
            'name': 'Product A',
            'list_price': 100,
            'sale_ok': True,
            'taxes_id': [],
        })

        giftcard_program = self.env['loyalty.program'].create([{
            'name': 'Gift Card Program',
            'program_type': 'gift_card',
            'trigger': 'auto',
            'applies_on': 'current',
            'rule_ids': [Command.create({
                'reward_point_amount': 1,
                'reward_point_mode': 'unit',
            })],
        }])

        order = self.env['sale.order'].with_user(self.user_salemanager).create({
            'partner_id': self.partner.id,
            'order_line': [
                Command.create({
                    'product_id': product_A.id,
                }),
            ]
        })

        order._update_programs_and_rewards()
        self._claim_reward(order, giftcard_program)

        self.assertEqual(giftcard_program.coupon_count, 0)

    def test_ewallet_code_use_restriction(self):
        self.env['loyalty.generate.wizard'].with_context(active_id=self.ewallet_program.id).create({
            'coupon_qty': 1,
            'points_granted': 100,
        }).generate_coupons()

        order = self.env['sale.order'].with_user(self.user_salemanager).create({
            'partner_id': self.partner.id,
            'order_line': [
                Command.create({
                    'product_id': self.product_a.id,
                }),
            ],
        })

        with self.assertRaises(ValidationError):
            self._apply_promo_code(order, self.ewallet_program.coupon_ids[0].code)

    def test_100_percent_discount(self):
        """
        Check whether a program offering 100% discount on an order reduces the order's total amount
        to zero.

        Assumes global tax rounding, as there's no good way to ensure the tax of the reward product
        equals the sum of taxes of the lines when each of them gets rounded.
        """
        self.env.company.tax_calculation_rounding_method = 'round_globally'
        loyalty_program = self.env['loyalty.program'].create([{
            'name': 'Full Discount',
            'program_type': 'loyalty',
            'trigger': 'auto',
            'applies_on': 'both',
            'rule_ids': [(0, 0, {
                'reward_point_mode': 'unit',
                'reward_point_amount': 1,
                'product_ids': [self.product_a.id],
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount': 100,
                'discount_mode': 'percent',
                'discount_applicability': 'order',
                'required_points': 1,
            })],
        }])
        self.env['loyalty.card'].create({
            'program_id': loyalty_program.id, 'partner_id': self.partner.id, 'points': 2
        })
        order = self.env['sale.order'].create({
            'partner_id': self.partner.id,
            'order_line': [(0, 0, {
                'product_id': self.product_A.id, 'product_uom_qty': 1, 'price_unit': price
            }) for price in (5.60, 8.92, 44.91, 217.26, 2400.00)],
        })

        order._update_programs_and_rewards()
        self._claim_reward(order, loyalty_program)
        msg = "100% discount on order should reduce total amount to 0"
        self.assertEqual(order.amount_total, 0, msg=msg)

    def test_discount_on_taxes_with_child_tax(self):
        """
        Check whether a program discount properly apply when product contain group of tax.
        """
        self.env.company.tax_calculation_rounding_method = 'round_globally'
        loyalty_program = self.env['loyalty.program'].create([{
            'name': '90% Discount',
            'program_type': 'loyalty',
            'trigger': 'auto',
            'applies_on': 'both',
            'rule_ids': [(0, 0, {
                'reward_point_mode': 'unit',
                'reward_point_amount': 1,
                'product_ids': [self.product_a.id],
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount': 90,
                'discount_mode': 'percent',
                'discount_applicability': 'order',
                'required_points': 1,
            })],
        }])
        self.env['loyalty.card'].create({'program_id': loyalty_program.id, 'partner_id': self.partner.id, 'points': 2})
        order = self.env['sale.order'].create({
            'partner_id': self.partner.id,
            'order_line': [(0, 0, {'product_id': self.product_D.id, 'product_uom_qty': 1})],
        })

        order._update_programs_and_rewards()
        self._claim_reward(order, loyalty_program)
        msg = "Discountable should take child tax amount into account"
        self.assertEqual(order.amount_total, 10, msg=msg)

    def test_ewallet_program_without_trigger_product(self):
        self.ewallet_program.trigger_product_ids = [Command.clear()]
        self.ewallet.points = 1000

        order = self.env['sale.order'].create({
            'partner_id': self.partner.id,
            'order_line': [Command.create({
                'product_id': self.product_a.id,
                'points_cost': 100,
                'product_uom_qty': 1,
            })],
        })
        order._update_programs_and_rewards()
        self._claim_reward(order, self.ewallet_program, coupon=self.ewallet)
        order.action_confirm()

        self.assertEqual(self.ewallet.points, 900)

    def test_ewallet_applied_ewallet_topup_in_order(self):
        self.ewallet.points = 10
        ewallet_top_up = Command.create({
            'product_id': self.env.ref('loyalty.ewallet_product_50').id,
            'product_uom_qty': 1,
            'price_unit': 50,
        })
        order = self.env['sale.order'].create({
            'partner_id': self.partner.id,
            'order_line': [Command.create({
                'product_id': self.product_a.id,
                'points_cost': 100,
                'product_uom_qty': 1,
            }),
                ewallet_top_up
            ],
        })
        order._update_programs_and_rewards()
        self._claim_reward(order, self.ewallet_program, coupon=self.ewallet)
        order.action_confirm()

        self.assertEqual(self.ewallet.points, 50)

        # Case 2: eWallet top-up should be excluded from the discountable amount when paying with an eWallet
        order = self.env['sale.order'].create({
            'partner_id': self.partner.id,
            'order_line': [ewallet_top_up],
        })
        order._update_programs_and_rewards()
        with self.assertRaisesRegex(UserError, "There is nothing to discount"):
            self._claim_reward(order, self.ewallet_program, coupon=self.ewallet)

    def test_discount_reward_claimable_only_once(self):
        """
        Check that discount rewards already applied won't be shown in the claimable rewards anymore.
        """
        program = self.env['loyalty.program'].create({
            'name': "10% Discount & Gift",
            'applies_on': 'current',
            'trigger': 'with_code',
            'program_type': 'promotion',
            'rule_ids': [Command.create({'mode': 'with_code', 'code': "10PERCENT&GIFT"})],
            'reward_ids': [
                Command.create({
                    'reward_type': 'product',
                    'reward_product_id': self.product_B.id,
                    'reward_product_qty': 1,
                }),
                Command.create({
                    'reward_type': 'discount',
                    'discount': 10,
                    'discount_mode': 'percent',
                    'discount_applicability': 'specific',
                }),
            ],
        })

        coupon = self.env['loyalty.card'].create({
            'program_id': program.id, 'points': 20, 'code': 'GIFT_CARD'
        })

        order = self.env['sale.order'].create({
            'partner_id': self.partner.id,
            'order_line': [Command.create({'product_id': self.product_a.id})]
        })

        product_reward = program.reward_ids.filtered(lambda reward: reward.reward_type == 'product')
        discount_reward = program.reward_ids - product_reward
        order._apply_program_reward(discount_reward, coupon)
        rewards = order._get_claimable_rewards()[coupon]
        msg = "Only the free product should be applicable, as the discount was already applied."
        self.assertEqual(rewards, product_reward, msg)

    def test_archived_reward_products(self):
        """
        Check that we do not use loyalty rewards that have no active reward product.
        In the case where the reward is based on reward_product_tag_id we also check
        the case where at least one reward is  active.
        """

        LoyaltyProgram = self.env['loyalty.program']
        loyalty_program = LoyaltyProgram.create(LoyaltyProgram._get_template_values()['loyalty'])
        loyalty_program_tag = LoyaltyProgram.create(LoyaltyProgram._get_template_values()['loyalty'])

        free_product_tag = self.env['product.tag'].create({'name': 'Free Product'})
        self.product_b.write({'product_tag_ids': [(4, free_product_tag.id)]})
        product_c = self.env['product.template'].create(
            {
                'name': 'Free Product C',
                'list_price': 1,
                'product_tag_ids': [(4, free_product_tag.id)],
            }
        )

        loyalty_program.reward_ids[0].write({
            'reward_type': 'product',
            'required_points': 1,
            'reward_product_id': self.product_b,
        })
        loyalty_program_tag.reward_ids[0].write({
            'reward_type': 'product',
            'required_points': 1,
            'reward_product_tag_id': free_product_tag.id,
        })
        self.product_b.active = False
        product_c.active = False

        order = self.env['sale.order'].create({
            'partner_id': self.partner.id,
            'order_line': [
                Command.create({
                    'product_id': self.product_a.id,
                }),
            ]
        })

        order._update_programs_and_rewards()
        rewards = [value.ids for value in order._get_claimable_rewards().values()]
        self.assertTrue(all(loyalty_program.reward_ids[0].id not in r for r in rewards))
        self.assertTrue(all(loyalty_program_tag.reward_ids[0].id not in r for r in rewards))

        product_c.active = True
        order._update_programs_and_rewards()
        rewards = [value.ids for value in order._get_claimable_rewards().values()]
        self.assertTrue(any(loyalty_program_tag.reward_ids[0].id in r for r in rewards))

    def test_domain_on_cheapest_reward(self):
        product_tag = self.env['product.tag'].create({'name': "Discountable"})
        self.env['loyalty.program'].create({
            'name': "10% Discount",
            'program_type': 'promo_code',
            'rule_ids': [Command.create({'code': "10discount"})],
            'reward_ids': [
                Command.create({
                    'reward_type': 'discount',
                    'discount': 10,
                    'discount_mode': 'percent',
                    'discount_applicability': 'cheapest',
                    'discount_product_tag_id': product_tag.id,
                }),
            ],
        })
        self.product_A.product_tag_ids = product_tag
        order = self.empty_order
        order.write({
            'order_line':[
                # product_A: lst_price: 100, Tax included price: 115
                Command.create({'product_id': self.product_A.id}),
                # Product_B: lst_price: 5, Tax included price: 5.75
                Command.create({'product_id': self.product_B.id}),
            ]
        })
        self._apply_promo_code(order, '10discount')
        msg = "Discount should only be applied to the line with a correctly tagged product."
        self.assertEqual(order.order_line[2].price_total, -11.5, msg)

        self.product_C.write({
            'list_price': 50,
            'product_tag_ids': product_tag,
        })
        order.order_line[2:].unlink()
        order.write({
            'order_line':[
                # product_C: lst_price = Tax included price: 50
                Command.create({'product_id': self.product_C.id}),
            ]
        })
        self._apply_promo_code(order, '10discount')
        msg = "Discount should be applied to the line with the cheapest valid product."
        self.assertEqual(order.order_line[3].price_total, -5.0, msg)

    def test_sol_free_product_description_equals_reward_description(self):
        """
        Ensure that if a "Free Product" reward is added to a sale order,
        its line description matches the reward description.
        """
        loyalty_program = self.env['loyalty.program'].create(
            self.env['loyalty.program']._get_template_values()['buy_x_get_y']
        )
        reward = loyalty_program.reward_ids[0]
        updated_description = f"{reward.description} Adding manual description"
        reward.description = updated_description

        order = self.empty_order
        order.write({
            'order_line': [
                Command.create({
                    'product_id': reward.reward_product_id.id,
                    'name': '1 Product',
                    'product_uom_qty': 4.0,
                }),
            ]
        })
        order._update_programs_and_rewards()
        self._claim_reward(order, loyalty_program)

        self.assertEqual(len(order.order_line.ids), 2)
        self.assertEqual(order.order_line[1].name, updated_description)

    def test_archiving_loyalty_card_unlinks_draft_points_from_sale_order(self):
        """
        When a loyalty card has points accrued from a draft sale order, archiving the
        card should unlink those draft points so they are no longer claimable on that order
        """
        loyalty_program = self.env['loyalty.program'].create({
            'name': 'Loyalty Program',
            'program_type': 'loyalty',
            'trigger': 'auto',
            'applies_on': 'both',
            'rule_ids': [
                Command.create({
                    'reward_point_mode': 'unit',
                    'reward_point_amount': 100,
                    'product_ids': [self.product_a.id],
                }),
            ],
            'reward_ids': [
                Command.create({
                    'reward_type': 'discount',
                    'discount': 50,
                    'discount_mode': 'percent',
                    'discount_applicability': 'order',
                    'required_points': 10,
                }),
            ],
        })
        loyalty_card = self.env['loyalty.card'].create({
            'program_id': loyalty_program.id,
            'partner_id': self.partner.id,
            'points': 0,
        })
        sale_order = self.empty_order
        sale_order.write({
            'order_line': [
                Command.create({
                    'product_id': self.product_a.id,
                }),
            ]
        })
        sale_order._update_programs_and_rewards()
        claimable_rewards = sale_order._get_claimable_rewards()
        self.assertTrue(claimable_rewards[loyalty_card])
        loyalty_card.action_archive()
        claimable_rewards = sale_order._get_claimable_rewards()
        self.assertFalse(claimable_rewards.get(loyalty_card))

    def test_free_product_sol_is_zero_price(self):
        self.env['res.config.settings'].create({
            'group_discount_per_so_line': True,
        }).execute()
        loyalty_program = self.env['loyalty.program'].create({
            'name': 'Loyalty Program',
            'program_type': 'promotion',
            'trigger': 'auto',
            'applies_on': 'both',
            'rule_ids': [
                Command.create({
                    'reward_point_mode': 'unit',
                    'reward_point_amount': 1,
                    'product_ids': [self.product_a.id],
                }),
            ],
            'reward_ids': [
                Command.create({
                    'reward_type': 'product',
                    'reward_product_id': self.product_B.id,
                    'reward_product_qty': 1,
                    'required_points': 1,
                }),
            ],
        })
        sale_order = self.empty_order
        sale_order.write({
            'order_line': [
                Command.create({
                    'product_id': self.product_a.id,
                    'product_uom_qty': 1,
                }),
            ],
        })
        sale_order._update_programs_and_rewards()
        self._claim_reward(sale_order, loyalty_program)
        # In real use case, so.plan_id is set to False in _verify_cart_after_update in
        # sale_subscription module. Since discount depends on so.plan_id, this triggers
        # a recomputation of the discount.
        # Here we manually call the compute method to simulate the behavior
        sale_order.order_line._compute_discount()
        reward_line = sale_order.order_line.filtered('reward_id')
        self.assertEqual(reward_line.discount, 100)
        self.assertEqual(reward_line.price_total, 0)

    def test_reapplying_reward_keeps_reward_price_unit(self):
        """
        Ensure that re-applying a reward doesn't reset the existing reward line unit price to zero
        """
        self.immediate_promotion_program.active = True
        sale_order = self.empty_order
        sale_order.write({
            'order_line': [
                Command.create({
                    'product_id': self.product_A.id,
                    'product_uom_qty': 1,
                }),
            ],
        })
        sale_order._update_programs_and_rewards()
        self._claim_reward(sale_order, self.immediate_promotion_program)
        reward_line = sale_order.order_line.filtered('reward_id')
        reward_line_price_unit = reward_line.price_unit
        sale_order._update_programs_and_rewards()
        self._claim_reward(sale_order, self.immediate_promotion_program)
        self.assertEqual(reward_line.price_unit, reward_line_price_unit)

    @freeze_time("2026-01-10")
    def test_expired_ewallet_is_not_claimable(self):
        self.ewallet.expiration_date = '2026-01-01'
        sale_order = self.empty_order
        sale_order.write({
            'partner_id': self.partner.id,
            'order_line': [
                Command.create({
                    'product_id': self.product_a.id,
                }),
            ],
        })
        sale_order.action_open_reward_wizard()
        sale_order._update_programs_and_rewards()
        claimable_rewards = sale_order._get_claimable_rewards()
        self.assertFalse(claimable_rewards.get(self.ewallet))
````

### FILE: `odoo_loyalty/upstream/addons/sale_loyalty/tests/test_loyalty_history.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file32:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/tests/test_loyalty_history.py; revision 99edb6dd82b7b560930c00b03b694ba700785370; Git blob 567039fe0292fd2c056c0c82e97781ad15934286"
license: "LGPL-3.0-only"
sha256: "291cee97246ffe07b39d77b4c0bdaa3c3addb7cec321aa2235827ab92515775d"
variables: []
secrets_allowed: false
```

````python
# Part of Odoo. See LICENSE file for full copyright and licensing details.

from odoo.fields import Command
from odoo.tests import tagged

from odoo.addons.sale_loyalty.tests.common import TestSaleCouponCommon


@tagged('post_install', '-at_install')
class TestLoyaltyhistory(TestSaleCouponCommon):

    @classmethod
    def setUpClass(cls):
        super().setUpClass()
        cls.partner_a = cls.env['res.partner'].create({'name': 'Jean Jacques'})
        cls.loyalty_program = cls.env['loyalty.program'].create({
            'name': 'Full Discount',
            'program_type': 'loyalty',
            'trigger': 'auto',
            'applies_on': 'both',
            'rule_ids': [Command.create({
                'reward_point_mode': 'unit',
                'reward_point_amount': 1,
                'product_ids': [cls.product_A.id],
            })],
            'reward_ids': [
                Command.create({
                    'reward_type': 'discount',
                    'discount': 10,
                    'discount_mode': 'percent',
                    'discount_applicability': 'order',
                    'required_points': 1,
                }),
                Command.create({
                    'active': False,
                    'reward_type': 'product',
                    'reward_product_id': cls.product_B.id,
                    'required_points': 2,
                }),
            ],
        })
        cls.loyalty_card = cls.env['loyalty.card'].create({
            'program_id': cls.loyalty_program.id, 'partner_id': cls.partner_a.id, 'points': 2
        })

    def test_add_loyalty_history_line_with_reward(self):
        order = self.empty_order
        order.write({
            'order_line': [
                Command.create({
                'product_id': self.product_A.id,
                'name': 'Ordinary Product A',
                'product_uom_qty': 1.0,
                }),
            ],
        })
        order._update_programs_and_rewards()
        self._auto_rewards(order, self.immediate_promotion_program)

        order.action_confirm()
        coupon_applied = self.immediate_promotion_program.coupon_ids.filtered(lambda x: x.order_id == order)
        history_records = len(coupon_applied.history_ids.filtered(lambda history: history.order_id == order.id))
        self.assertEqual(history_records, 1, "A history line should be created on confirmation of order")

    def test_add_loyalty_history_line_without_reward(self):
        order = self.empty_order
        order.write({
            'partner_id': self.partner_a.id,
            'order_line': [
                Command.create({
                    'product_id': self.product_A.id,
                    'tax_ids': False,
                }),
            ]
        })
        order.action_confirm()
        order._update_programs_and_rewards()
        self._claim_reward(order, self.loyalty_program)
        history_records = self.loyalty_card.history_ids.filtered(lambda history: history.order_id == order.id)
        self.assertEqual(history_records.used, 1.0,
                        "The history line should be updated on change of order lines in a confirmed order")

    def test_delete_loyalty_history_line_on_cancel(self):
        order = self.empty_order
        order.write({
            'partner_id': self.partner_a.id,
            'order_line': [
                Command.create({
                    'product_id': self.product_A.id,
                    'tax_ids': False,
                }),
            ]
        })
        order._update_programs_and_rewards()
        self._claim_reward(order, self.loyalty_program)
        order.action_confirm()
        lines_before_cancel = len(self.loyalty_card.history_ids)
        order._action_cancel()
        self.assertEqual(lines_before_cancel - 1, len(self.loyalty_card.history_ids),
                         "History line should be deleted after order cancel")

    def test_loyalty_history_multi_reward(self):
        """Verify that applying multiple rewards sums up the total points cost."""
        self.loyalty_card.points = initial_points = 4
        self.loyalty_program.with_context(active_test=False).reward_ids.active = True
        order = self.empty_order
        order.write({
            'partner_id': self.partner_a.id,
            'order_line': [
                Command.create({
                    'product_id': self.product_A.id,
                    'tax_ids': False,
                }),
            ],
        })
        for reward in self.loyalty_program.reward_ids:
            order._apply_program_reward(reward, self.loyalty_card)
        self.assertEqual(len(order.order_line.filtered('reward_id')), 2)
        self.assertEqual(order.order_line.mapped('points_cost'), [0, 1, 2])

        order.action_confirm()
        loyalty_history = self.loyalty_card.history_ids
        self.assertEqual(loyalty_history.issued, 1, "1 point should be rewarded")
        self.assertEqual(loyalty_history.used, 3, "A total of 3 points should be used")
        self.assertEqual(
            self.loyalty_card.points,
            initial_points + loyalty_history.issued - loyalty_history.used,
            "Loyalty points should equal initial points + points issued - points used",
        )

    def test_loyalty_history_created_on_post_confirm_reward(self):
        """ History line must be created when a reward is claimed on a confirmed
            sale order where the card was not used on the order before confirmation.
        """
        order = self.env['sale.order'].create({
            'partner_id': self.partner_a.id,
            'order_line': [Command.create({
                'product_id': self.product_A.id,
                'tax_ids': False,
            })],
        })
        order.action_confirm()
        order._update_programs_and_rewards()
        coupon = order.coupon_point_ids.coupon_id.filtered(lambda c: c.program_id == self.immediate_promotion_program)
        self.assertFalse(
            coupon.history_ids.filtered(lambda h: h.order_id == order.id and h.order_model == 'sale.order'),
            "No history should exist before claiming reward",
        )
        self._claim_reward(order, self.immediate_promotion_program, coupon)
        history = coupon.history_ids.filtered(lambda h: h.order_id == order.id and h.order_model == 'sale.order')
        self.assertEqual(len(history), 1, "History line must be created when reward is claimed on confirmed order")
        self.assertEqual(history.used, 1.0, "History used should reflect the reward points cost")

    def test_loyalty_history_updated_on_points_cost_write(self):
        """ write() on sale.order.line must update history.used by the delta
            when points_cost changes on a confirmed order.
        """
        self.loyalty_card.points = 10
        self.loyalty_program.with_context(active_test=False).reward_ids.active = True
        product_reward = self.loyalty_program.reward_ids.filtered(
            lambda r: r.reward_type == 'product'
        )
        order = self.env['sale.order'].create({
            'partner_id': self.partner_a.id,
            'order_line': [Command.create({
                'product_id': self.product_A.id,
                'product_uom_qty': 2,
                'tax_ids': False,
            })],
        })
        order._update_programs_and_rewards()
        order._apply_program_reward(product_reward, self.loyalty_card)
        order.action_confirm()
        history = self.loyalty_card.history_ids.filtered(lambda h: h.order_id == order.id and h.order_model == 'sale.order')
        self.assertEqual(len(history), 1, "A history line should exist after confirmation")
        used_after_confirm = history.used
        reward_line = order.order_line.filtered('reward_id')
        reward_line.write({'points_cost': reward_line.points_cost + 1})
        self.assertEqual(history.used, used_after_confirm + 1, "history.used must increase by delta when points_cost is written on confirmed order")
````

### FILE: `odoo_loyalty/upstream/addons/sale_loyalty/tests/test_pay_with_gift_card.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file33:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/tests/test_pay_with_gift_card.py; revision 99edb6dd82b7b560930c00b03b694ba700785370; Git blob 13c06bbf7ae046d55c2ac3c183663c31840f0d85"
license: "LGPL-3.0-only"
sha256: "9ecd12f97ddb7dcdcb8147d38752a1155c21583977e2ce048e40e653b4468f47"
variables: []
secrets_allowed: false
```

````python
# Part of Odoo. See LICENSE file for full copyright and licensing details.

from odoo.fields import Command
from odoo.tests import tagged

from odoo.addons.sale_loyalty.tests.common import TestSaleCouponCommon


@tagged('-at_install', 'post_install')
class TestPayWithGiftCard(TestSaleCouponCommon):

    def test_paying_with_single_gift_card_over(self):
        self.env['loyalty.generate.wizard'].with_context(active_id=self.program_gift_card.id).create({
            'coupon_qty': 1,
            'points_granted': 100,
        }).generate_coupons()
        gift_card = self.program_gift_card.coupon_ids[0]
        order = self.empty_order
        order.write({'order_line': [
            Command.create({
                'product_id': self.product_A.id,
                'name': 'Ordinary Product A',
                'product_uom_qty': 1.0,
            })
        ]})
        before_gift_card_payment = order.amount_total
        self.assertNotEqual(before_gift_card_payment, 0)
        self._apply_promo_code(order, gift_card.code)
        order.action_confirm()
        self.assertEqual(before_gift_card_payment - order.amount_total, 100 - gift_card.points)

    def test_paying_with_single_gift_card_under(self):
        self.env['loyalty.generate.wizard'].with_context(active_id=self.program_gift_card.id).create({
            'coupon_qty': 1,
            'points_granted': 100,
        }).generate_coupons()
        gift_card = self.program_gift_card.coupon_ids[0]
        order = self.empty_order
        order.write({'order_line': [
            Command.create({
                'product_id': self.product_B.id,
                'name': 'Ordinary Product b',
                'product_uom_qty': 1.0,
            })
        ]})
        before_gift_card_payment = order.amount_total
        self.assertNotEqual(before_gift_card_payment, 0)
        self._apply_promo_code(order, gift_card.code)
        order.action_confirm()
        self.assertEqual(before_gift_card_payment - order.amount_total, 100 - gift_card.points)

    def test_paying_with_multiple_gift_card(self):
        self.env['loyalty.generate.wizard'].with_context(active_id=self.program_gift_card.id).create({
            'coupon_qty': 2,
            'points_granted': 100,
        }).generate_coupons()
        gift_card_1, gift_card_2 = self.program_gift_card.coupon_ids
        order = self.empty_order
        order.write({'order_line': [
            Command.create({
                'product_id': self.product_A.id,
                'name': 'Ordinary Product A',
                'product_uom_qty': 20.0,
            })
        ]})
        before_gift_card_payment = order.amount_total
        self._apply_promo_code(order, gift_card_1.code)
        self._apply_promo_code(order, gift_card_2.code)
        self.assertEqual(order.amount_total, before_gift_card_payment - 200)

    def test_paying_with_gift_card_and_discount(self):
        # Test that discounts take precedence on payment rewards
        self.env['loyalty.generate.wizard'].with_context(active_id=self.program_gift_card.id).create({
            'coupon_qty': 1,
            'points_granted': 50,
        }).generate_coupons()
        gift_card_1 = self.program_gift_card.coupon_ids
        order = self.empty_order
        order.write({'order_line': [
            Command.create({
                'product_id': self.product_C.id,
                'name': 'Ordinary Product C',
                'product_uom_qty': 1.0,
            })
        ]})
        self.env['loyalty.program'].create({
            'name': 'Code for 10% on orders',
            'trigger': 'with_code',
            'program_type': 'promotion',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {
                'mode': 'with_code',
                'code': 'test_10pc',
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount_mode': 'percent',
                'discount': 10,
                'discount_applicability': 'order',
                'required_points': 1,
            })],
        })
        self.assertEqual(order.amount_total, 100)
        self._apply_promo_code(order, gift_card_1.code)
        self.assertEqual(order.amount_total, 50)
        self._apply_promo_code(order, "test_10pc")
        # real flows also have to update the programs and rewards
        order._update_programs_and_rewards()
        self.assertEqual(order.amount_total, 40) # 100 - 10% - 50

    def test_paying_with_gift_card_blocking_discount(self):
        # Test that a payment program making the order total 0 still allows the user to claim discounts
        self.env['loyalty.generate.wizard'].with_context(active_id=self.program_gift_card.id).create({
            'coupon_qty': 1,
            'points_granted': 100,
        }).generate_coupons()
        gift_card_1 = self.program_gift_card.coupon_ids
        order = self.empty_order
        order.write({'order_line': [
            Command.create({
                'product_id': self.product_C.id,
                'name': 'Ordinary Product C',
                'product_uom_qty': 1.0,
            })
        ]})
        self.env['loyalty.program'].create({
            'name': 'Code for 10% on orders',
            'trigger': 'with_code',
            'program_type': 'promotion',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {
                'mode': 'with_code',
                'code': 'test_10pc',
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount_mode': 'percent',
                'discount': 10,
                'discount_applicability': 'order',
                'required_points': 1,
            })],
        })
        self.assertEqual(order.amount_total, 100)
        self._apply_promo_code(order, gift_card_1.code)
        self.assertEqual(order.amount_total, 0)
        self._apply_promo_code(order, "test_10pc")
        # real flows also have to update the programs and rewards
        order._update_programs_and_rewards()
        self.assertEqual(order.amount_total, 0) # 100 - 10% - 90

    def test_gift_card_product_has_no_taxes_on_creation(self):
        gift_card_program = self.env['loyalty.program'].create({
            'name': 'Gift Cards',
            'applies_on': 'future',
            'program_type': 'gift_card',
            'trigger': 'auto',
            'rule_ids': [Command.create({
                'product_ids': self.product_gift_card,
                'reward_point_amount': 1,
                'reward_point_mode': 'money',
                'reward_point_split': True,
            })],
            'reward_ids': [Command.create({
                'reward_type': 'discount',
                'discount': 1,
                'discount_mode': 'per_point',
                'discount_applicability': 'order',
            })]
        })
        self.assertFalse(gift_card_program.reward_ids.discount_line_product_id.taxes_id)

    def test_paying_with_gift_card_uses_gift_card_product_taxes(self):
        order = self.empty_order
        order.order_line = [
            Command.create({
                'product_id': self.product_B.id,
                'name': 'Ordinary Product b',
                'product_uom_qty': 1.0,
                'price_unit': 200.0,
            })
        ]
        sol = order.order_line
        before_gift_card_payment = order.amount_total
        self.assertNotEqual(before_gift_card_payment, 0)

        self.env['loyalty.generate.wizard'].with_context(active_id=self.program_gift_card.id).create({
            'coupon_qty': 1,
            'points_granted': 100,
        }).generate_coupons()
        gift_card = self.program_gift_card.coupon_ids[0]

        # TODO check amount total of gift_card_line

        # TAX EXCL
        self.program_gift_card.reward_ids.discount_line_product_id.taxes_id = [
            Command.link(self.tax_15pc_excl.id)
        ]
        self._apply_promo_code(order, gift_card.code)
        gift_card_line = order.order_line - sol
        self.assertAlmostEqual(gift_card_line.price_total, -100.0)
        self.assertAlmostEqual(order.amount_total, before_gift_card_payment - 100.0)
        self.assertTrue(all(line.tax_ids for line in order.order_line))
        self.assertEqual(order.order_line.tax_ids, self.tax_15pc_excl)

        # TAX INCL
        gift_card_line.unlink()  # Remove gift card
        self.program_gift_card.reward_ids.discount_line_product_id.taxes_id = [
            Command.set(self.tax_10pc_incl.ids)
        ]
        self._apply_promo_code(order, gift_card.code)
        gift_card_line = order.order_line - sol
        self.assertAlmostEqual(gift_card_line.price_total, -100.0)
        self.assertAlmostEqual(order.amount_total, before_gift_card_payment - 100.0)
        self.assertTrue(all(line.tax_ids for line in order.order_line))
        self.assertEqual(gift_card_line.tax_ids, self.tax_10pc_incl)

        # TAX INCL + TAX EXCL
        gift_card_line.unlink()  # Remove gift card
        self.program_gift_card.reward_ids.discount_line_product_id.taxes_id = [
            Command.link(self.tax_15pc_excl.id)
        ]
        self._apply_promo_code(order, gift_card.code)
        gift_card_line = order.order_line - sol
        self.assertAlmostEqual(gift_card_line.price_total, -100.0)
        self.assertAlmostEqual(order.amount_total, before_gift_card_payment - 100.0)
        self.assertTrue(all(line.tax_ids for line in order.order_line))
        self.assertEqual(gift_card_line.tax_ids, self.tax_10pc_incl + self.tax_15pc_excl)

    def test_paying_with_gift_card_fixed_tax(self):
        """ Test payment of sale order with fixed tax using gift card """
        self.env['loyalty.generate.wizard'].with_context(active_id=self.program_gift_card.id).create({
            'coupon_qty': 1,
            'points_granted': 100,
        }).generate_coupons()
        gift_card = self.program_gift_card.coupon_ids[0]

        tax_10_fixed = self.env['account.tax'].create({
            'name': "10$ Fixed tax",
            'amount_type': 'fixed',
            'amount': 10,
        })
        self.product_A.write({'list_price': 90})
        self.product_A.taxes_id = tax_10_fixed

        order = self.empty_order
        order.write({'order_line': [
            Command.create({
                'product_id': self.product_A.id,
                'name': "Ordinary Product A",
                'product_uom_qty': 1.0,
            })
        ]})
        self._apply_promo_code(order, gift_card.code)
        order.action_confirm()
        self.assertEqual(order.amount_total, 0, "The order should be totally paid")
````

### FILE: `odoo_loyalty/upstream/addons/sale_loyalty/tests/test_program_multi_company.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file34:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/tests/test_program_multi_company.py; revision 99edb6dd82b7b560930c00b03b694ba700785370; Git blob c06a430d127004821e968ab58ff3dd0666a5499c"
license: "LGPL-3.0-only"
sha256: "d338ab072168bacfa342f7a5a0011063d47f41f19ef1221a4adf0a70cb90a32b"
variables: []
secrets_allowed: false
```

````python
# Part of Odoo. See LICENSE file for full copyright and licensing details.

from odoo.fields import Command
from odoo.tests import tagged

from odoo.addons.sale_loyalty.tests.common import TestSaleCouponCommon


@tagged('post_install', '-at_install')
class TestSaleCouponMultiCompany(TestSaleCouponCommon):

    def setUp(self):
        super(TestSaleCouponMultiCompany, self).setUp()

        self.company_a = self.env.company
        self.company_b = self.env['res.company'].create(dict(name="TEST"))

        self.immediate_promotion_program_c2 = self.env['loyalty.program'].create({
            'name': 'Buy A + 1 B, 1 B are free',
            'trigger': 'auto',
            'program_type': 'promotion',
            'applies_on': 'current',
            'company_id': self.company_b.id,
            'rule_ids': [(0, 0, {
                'product_ids': self.product_A,
                'reward_point_amount': 1,
                'reward_point_mode': 'order',
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'product',
                'reward_product_id': self.product_B.id,
                'reward_product_qty': 1,
                'required_points': 1,
            })],
        })

    def _get_applicable_programs(self, order):
        return self.env['loyalty.program'].browse(p.id for p in order._get_applicable_program_points())

    def test_applicable_programs(self):

        order = self.empty_order
        order.write({'order_line': [
            (0, False, {
                'product_id': self.product_A.id,
                'name': '1 Product A',
                'product_uom_qty': 1.0,
            }),
            (0, False, {
                'product_id': self.product_B.id,
                'name': '2 Product B',
                'product_uom_qty': 1.0,
            })
        ]})
        order._update_programs_and_rewards()

        self.assertNotIn(self.immediate_promotion_program_c2, self._get_applicable_programs(order))
        self.assertNotIn(self.immediate_promotion_program_c2, order._get_applied_programs())

        order_b = self.env["sale.order"].create({
            'company_id': self.company_b.id,
            'partner_id': order.partner_id.id,
        })
        order_b.write({'order_line': [
            (0, False, {
                'product_id': self.product_A.id,
                'name': '1 Product A',
                'product_uom_qty': 1.0,
            }),
            (0, False, {
                'product_id': self.product_B.id,
                'name': '2 Product B',
                'product_uom_qty': 1.0,
            })
        ]})
        self.assertNotIn(self.immediate_promotion_program, self._get_applicable_programs(order_b))
        order_b._update_programs_and_rewards()
        self.assertIn(self.immediate_promotion_program_c2, order_b._get_applied_programs())
        self.assertNotIn(self.immediate_promotion_program, order_b._get_applied_programs())

    def test_applicable_programs_on_branch(self):
        # create a branch
        branch_a = self.env['res.company'].create(
            {'name': 'Branch A', 'parent_id': self.company_a.id}
        )

        # create an order
        order = self.env['sale.order'].create(
            {'order_line': [
                Command.create({
                    'product_id': self.product_A.id,
                    'name': '1 Product A',
                    'product_uom_qty': 1.0,
                }),
                Command.create({
                    'product_id': self.product_B.id,
                    'name': '2 Product B',
                    'product_uom_qty': 1.0,
                })
            ],
            'company_id': branch_a.id,
            'partner_id': self.partner.id
            }
        )

        order._update_programs_and_rewards()
        self.assertIn(self.immediate_promotion_program, order._get_applied_programs())

    def test_applicable_programs_confirm_on_branch(self):
        # create a branch
        self.env['loyalty.program'].search([]).write({'active': False})
        branch_a = self.env['res.company'].create(
            {'name': 'Branch A', 'parent_id': self.company_a.id}
        )

        LoyaltyProgram = self.env['loyalty.program']
        LoyaltyProgram.create(LoyaltyProgram._get_template_values()['loyalty'])

        self.sale_user.write({'company_ids': [Command.set((branch_a + self.company_a).ids)]})

        # create an order
        order = self.empty_order
        order.update(
            {
                'order_line': [
                    Command.create({
                        'product_id': self.product_A.id,
                    }),
                ],
                'company_id': branch_a.id,
                'partner_id': self.partner.id,
                'user_id': self.sale_user.id
            }
        )

        order.with_user(self.sale_user).with_company(branch_a.id).sudo(False).action_confirm()
        self.assertEqual(order.state, 'sale')
````

### FILE: `odoo_loyalty/upstream/addons/sale_loyalty/tests/test_program_numbers.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file35:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/tests/test_program_numbers.py; revision 99edb6dd82b7b560930c00b03b694ba700785370; Git blob 45d08547f109e952daa2b268ff298727c0d1ddac"
license: "LGPL-3.0-only"
sha256: "845f5b205b1ac941280905bfa308ed15d7142d3069997b6545f5cc0e46b0c918"
variables: []
secrets_allowed: false
```

````python
# Part of Odoo. See LICENSE file for full copyright and licensing details.

from odoo.exceptions import ValidationError
from odoo.fields import Command
from odoo.tests import tagged
from odoo.tools.float_utils import float_compare

from odoo.addons.sale_loyalty.tests.common import TestSaleCouponNumbersCommon


@tagged('post_install', '-at_install')
class TestSaleCouponProgramNumbers(TestSaleCouponNumbersCommon):

    def test_program_numbers_free_and_paid_product_qty(self):
        # These tests will focus on numbers (free product qty, SO total, reduction total..)
        order = self.empty_order
        sol1 = self.env['sale.order.line'].create({
            'product_id': self.largeCabinet.id,
            'name': 'Large Cabinet',
            'product_uom_qty': 3.0,
            'order_id': order.id,
        })

        # Check we correctly get a free product
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(len(order.order_line.ids), 2, "We should have 2 lines as we now have one 'Free Large Cabinet' line as we bought 3 of them")

        # Check free product's price is not added to total when applying reduction (Or the discount will also be applied on the free product's price)
        self._apply_promo_code(order, 'test_10pc')
        self.assertEqual(len(order.order_line.ids), 3, "We should have 3 lines as we should have a new line for promo code reduction")
        self.assertEqual(order.amount_total, 864, "Only paid product should have their price discounted")
        order.order_line.filtered(lambda x: 'Discount' in x.name).unlink()  # Remove Discount
        order._remove_program_from_points(self.p1)

        # Check free product is removed since we are below minimum required quantity
        sol1.product_uom_qty = 2
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(len(order.order_line.ids), 1, "Free Large Cabinet should have been removed")

        # Free product in cart will be considered as paid product when changing quantity of paid product, so the free product quantity computation will be wrong.
        # 75 Large Cabinet in cart, 25 free, set quantity to 6 Large Cabinet, you should have 2 free Large Cabinet but you get 8 because it add the 25 initial free Large Cabinet to the total paid Large Cabinet when computing (25+10 > 35 > /4 = 8 free Large Cabinet)
        sol1.product_uom_qty = 75
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(sum(order.order_line.filtered(lambda x: x.is_reward_line).mapped('product_uom_qty')), 25, "We should have 25 Free Large Cabinet")
        sol1.product_uom_qty = 6
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(sum(order.order_line.filtered(lambda x: x.is_reward_line).mapped('product_uom_qty')), 2, "We should have 2 Free Large Cabinet")

    def test_program_numbers_check_eligibility(self):
        # These tests will focus on numbers (free product qty, SO total, reduction total..)

        # Check if we have enough paid product to receive free product in case of a free product that is different from the paid product required
        # Buy A, get free b. (remember we need a paid B in cart to receive free b). If your cart is 4A 1B then you should receive 1b (you are eligible to receive 4 because you have 4A but since you dont have enought B in your cart, you are limited to the B quantity)
        order = self.empty_order
        sol1 = self.env['sale.order.line'].create({
            'product_id': self.drawerBlack.id,
            'name': 'drawer black',
            'product_uom_qty': 3.0,
            'order_id': order.id,
        })
        sol2 = self.env['sale.order.line'].create({
            'product_id': self.largeMeetingTable.id,
            'name': 'Large Meeting Table',
            'product_uom_qty': 1.0,
            'order_id': order.id,
        })
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(len(order.order_line.ids), 3, "We should have a 'Free Large Meeting Table' promotion line")
        self.assertEqual(sum(order.order_line.filtered(lambda x: x.is_reward_line).mapped('product_uom_qty')), 1, "We should receive one and only one free Large Meeting Table")

        # Check the required value amount to be eligible for the program is correctly computed (eg: it does not add negative value (from free product) to total)
        # A = free b | Have your cart with A 2B b | cart value should be A + 1B but in code it is only A (free b value is subsstract 2 times)
        # This is because _amount_all() is summing all SO lines (so + (-b.value)) and again in _check_promo_code() order.amount_untaxed + order.reward_amount | amount_untaxed has already free product value substracted (_amount_all)
        sol1.product_uom_qty = 1
        sol2.product_uom_qty = 2
        self.p1.rule_ids.minimum_amount = 5000
        self._auto_rewards(order, self.all_programs)
        self._apply_promo_code(order, 'test_10pc')
        self.assertEqual(len(order.order_line.ids), 4, "We should have 4 lines as we should have a new line for promo code reduction")

        # Check you can still have auto applied promotion if you have a promo code set to the order
        self.env['sale.order.line'].create({
            'product_id': self.largeCabinet.id,
            'name': 'Large Cabinet',
            'product_uom_qty': 4.0,
            'order_id': order.id,
        })
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(len(order.order_line.ids), 6, "We should have 2 more lines as we now have one 'Free Large Cabinet' line since we bought 4 of them")

    def test_program_numbers_taxes_and_rules(self):
        percent_tax = self.env['account.tax'].create({
            'name': "15% Tax",
            'amount_type': 'percent',
            'amount': 15,
            'price_include_override': 'tax_included',
        })
        p_specific_product = self.env['loyalty.program'].create({
            'name': '20% reduction on Large Cabinet in cart',
            'trigger': 'auto',
            'program_type': 'promotion',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {
                'reward_point_mode': 'order',
                'minimum_amount_tax_mode': 'excl',
                'minimum_amount': 320.00,
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount': 20,
                'discount_mode': 'percent',
                'discount_applicability': 'specific',
                'discount_product_ids': self.largeCabinet,
                'required_points': 1,
            })],
        })
        self.all_programs |= p_specific_product
        order = self.empty_order
        self.largeCabinet.taxes_id = percent_tax
        sol1 = self.env['sale.order.line'].create({
            'product_id': self.largeCabinet.id,
            'name': 'Large Cabinet',
            'product_uom_qty': 1.0,
            'order_id': order.id,
        })

        self._auto_rewards(order, self.all_programs)
        self.assertEqual(len(order.order_line.ids), 1, "We should not get the reduction line since we dont have 320$ tax excluded (cabinet is 320$ tax included)")
        sol1.tax_ids.price_include_override = 'tax_excluded'
        sol1._compute_tax_ids()
        self.env.flush_all()
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(len(order.order_line.ids), 2, "We should now get the reduction line since we have 320$ tax included (cabinet is 320$ tax included)")
        # Name                 | Qty | price_unit |  Tax     |  HTVA   |   TVAC  |  TVA  |
        # --------------------------------------------------------------------------------
        # Conference Chair     |  1  |    320.00  | 15% excl |  320.00 |  368.00 |   48.00
        # 20% discount on      |  1  |    -64.00  | 15% excl |  -64.00 |  -73.60 |   -9.60
        #        large cabinet |
        # --------------------------------------------------------------------------------
        # TOTAL                                              |  256.00 |  294.40 |   38.40
        self.assertAlmostEqual(order.amount_total, 294.4, 2, "Check discount has been applied correctly (eg: on taxes aswell)")

        # test coupon with code works the same as auto applied_programs
        p_specific_product.write({'trigger': 'with_code'})
        p_specific_product.rule_ids.write({'mode': 'with_code', 'code': '20pc'})
        order.order_line.filtered(lambda l: l.is_reward_line).unlink()
        order._remove_program_from_points(p_specific_product)
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(len(order.order_line.ids), 1, "Reduction should be removed since we deleted it and it is now a promo code usage, it shouldn't be automatically reapplied")

        self._apply_promo_code(order, '20pc')
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(len(order.order_line.ids), 2, "We should now get the reduction line since we have 320$ tax included (cabinet is 320$ tax included)")

        # check discount applied only on Large Cabinet
        self.env['sale.order.line'].create({
            'product_id': self.drawerBlack.id,
            'name': 'Drawer Black',
            'product_uom_qty': 10.0,
            'order_id': order.id,
        })
        self._auto_rewards(order, self.all_programs)
        # Name                 | Qty | price_unit |  Tax     |  HTVA   |   TVAC  |  TVA  |
        # --------------------------------------------------------------------------------
        # Drawer Black         | 10  |     25.00  |        / |  250.00 |  250.00 |       /
        # Large Cabinet        |  1  |    320.00  | 15% excl |  320.00 |  368.00 |   48.00
        # 20% discount on      |  1  |    -64.00  | 15% excl |  -64.00 |  -73.60 |   -9.60
        #        large cabinet |
        # --------------------------------------------------------------------------------
        # TOTAL                                              |  506.00 |  544.40 |   38.40
        self.assertEqual(order.amount_total, 544.4, "We should only get reduction on cabinet")
        sol1.product_uom_qty = 8
        self._auto_rewards(order, self.all_programs)
        # Note: Since we now have 2 free Large Cabinet, we should discount only 8 of the 10 Large Cabinet in carts since we don't want to discount free Large Cabinet
        # Name                 | Qty | price_unit |  Tax     |  HTVA   |   TVAC  |  TVA  |
        # --------------------------------------------------------------------------------
        # Drawer Black         | 10  |     25.00  |        / |  250.00 |  250.00 |       /
        # Large Cabinet        |  8  |    320.00  | 15% excl | 2560.00 | 2944.00 |  384.00
        # Free Large Cabinet   |  2  |      0.00  | 15% excl |    0.00 |    0.00 |    0.00
        # 20% discount on      |  1  |   -512.00  | 15% excl | -512.00 | -588.80 |  -78.80
        #        large cabinet |
        # --------------------------------------------------------------------------------
        # TOTAL                                              | 2298.00 | 2605.20 |  305.20
        self.assertAlmostEqual(order.amount_total, 2605.20, 2, "Changing cabinet quantity should change discount amount correctly")

        p_specific_product.reward_ids.discount_max_amount = 200
        self._auto_rewards(order, self.all_programs)
        # Name                 | Qty | price_unit |  Tax     |  HTVA   |   TVAC  |  TVA  |
        # --------------------------------------------------------------------------------
        # Drawer Black         | 10  |     25.00  |        / |  250.00 |  250.00 |       /
        # Large Cabinet        |  8  |    320.00  | 15% excl | 2560.00 | 2944.00 |  384.00
        # Free Large Cabinet   |  2  |      0.00  | 15% excl |    0.00 |    0.00 |    0.00
        # 20% discount on      |  1  |   -173.91  | 15% excl | -173.91 | -200.00 |  -26.09
        #        large cabinet |
        #  limited to 200 HTVA
        # --------------------------------------------------------------------------------
        # TOTAL                                              | 2636.09 | 2994.00 |  357.91
        self.assertEqual(order.amount_total, 2994.0, "The discount should be limited to $200 tax included")
        self.assertEqual(order.amount_untaxed, 2636.09, "The discount should be limited to $200 tax included (2)")

    def test_program_numbers_one_discount_line_per_tax(self):
        order = self.empty_order
        self.env['ir.config_parameter'].set_param('loyalty.compute_all_discount_product_ids', 'enabled')
        # Create taxes
        self.tax_15pc_excl = self.env['account.tax'].create({
            'name': "15% Tax excl",
            'amount_type': 'percent',
            'amount': 15,
        })
        self.tax_50pc_excl = self.env['account.tax'].create({
            'name': "50% Tax excl",
            'amount_type': 'percent',
            'amount': 50,
        })
        self.tax_35pc_incl = self.env['account.tax'].create({
            'name': "35% Tax incl",
            'amount_type': 'percent',
            'amount': 35,
            'price_include_override': 'tax_included',
        })

        # Set tax and prices on products as neeed for the test
        (self.product_A + self.largeCabinet + self.conferenceChair + self.pedalBin + self.drawerBlack).write({'list_price': 100})
        (self.largeCabinet + self.drawerBlack).write({'taxes_id': [(4, self.tax_15pc_excl.id, False)]})
        self.conferenceChair.taxes_id = self.tax_10pc_incl
        self.pedalBin.taxes_id = None
        self.product_A.taxes_id = (self.tax_35pc_incl + self.tax_50pc_excl)

        # Add products in order
        self.env['sale.order.line'].create({
            'product_id': self.largeCabinet.id,
            'name': 'Large Cabinet',
            'product_uom_qty': 4.0,
            'order_id': order.id,
        })
        sol2 = self.env['sale.order.line'].create({
            'product_id': self.conferenceChair.id,
            'name': 'Conference Chair',
            'product_uom_qty': 3.0,
            'order_id': order.id,
        })
        self.env['sale.order.line'].create({
            'product_id': self.pedalBin.id,
            'name': 'Pedal Bin',
            'product_uom_qty': 5.0,
            'order_id': order.id,
        })
        self.env['sale.order.line'].create({
            'product_id': self.product_A.id,
            'name': 'product A with multiple taxes',
            'product_uom_qty': 3.0,
            'order_id': order.id,
        })
        self.env['sale.order.line'].create({
            'product_id': self.drawerBlack.id,
            'name': 'Drawer Black',
            'product_uom_qty': 2.0,
            'order_id': order.id,
        })

        # Create needed programs
        self.immediate_promotion_program.active = False
        self.p2.active = False
        self.p3.active = False

        # NOTE: programs may not make much sense but they have been modified in order to validate the result since the change from coupon to loyalty.
        self.p_large_cabinet = self.env['loyalty.program'].create({
            'name': 'Buy 1 large cabinet, get 3/4 for free',
            'trigger': 'auto',
            'program_type': 'promotion',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {
                'product_ids': self.largeCabinet,
                'reward_point_mode': 'unit',
                'minimum_qty': 1,
                'reward_point_amount': 0.752,
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'product',
                'reward_product_id': self.largeCabinet.id,
                'reward_product_qty': 1,
                'required_points': 1,
            })],
        })
        self.p_conference_chair = self.env['loyalty.program'].create({
            'name': 'Buy 1 chair, get one for free',
            'trigger': 'auto',
            'program_type': 'promotion',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {
                'product_ids': self.conferenceChair,
                'reward_point_mode': 'unit',
                'minimum_qty': 1,
                'reward_point_amount': 0.4,
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'product',
                'reward_product_id': self.conferenceChair.id,
                'reward_product_qty': 1,
                'required_points': 1,
            })],
        })
        self.p_pedal_bin = self.env['loyalty.program'].create({
            'name': 'Buy 1 bin, get one for free',
            'trigger': 'auto',
            'program_type': 'promotion',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {
                'product_ids': self.pedalBin,
                'reward_point_mode': 'unit',
                'minimum_qty': 1,
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'product',
                'reward_product_id': self.pedalBin.id,
                'reward_product_qty': 1,
                'required_points': 1,
            })],
        })
        self.all_programs |= (self.p_large_cabinet | self.p_conference_chair | self.p_pedal_bin)
        # Name                 | Qty | price_unit |  Tax     |  HTVA   |   TVAC  |  TVA  |
        # --------------------------------------------------------------------------------
        # Conference Chair     |  5  |    100.00  | 10% incl |  454.55 |  500.00 |   45.45
        # Pedal bin            |  10 |    100.00  | /        | 1000.00 | 1000.00 |       /
        # Large Cabinet        |  7  |    100.00  | 15% excl |  700.00 |  805.00 |  105.00
        # Drawer Black         |  2  |    100.00  | 15% excl |  200.00 |  230.00 |   30.00
        # Product A            |  3  |    100.00  | 35% incl |  222.22 |  411.11 |  188.89
        #                                           50% excl
        # --------------------------------------------------------------------------------
        # TOTAL                                              | 2576.77 | 2946.11 |  369.34

        self.assertRecordValues(order, [{
            'amount_total': 1901.11,
            'amount_untaxed': 1594.95,
        }])
        self.assertEqual(len(order.order_line.ids), 5, "The order without any programs should have 5 lines")

        # Apply all the programs
        self._auto_rewards(order, self.all_programs)

        # Name                 | Qty | price_unit |  Tax     |  HTVA   |   TVAC  |  TVA  |
        # --------------------------------------------------------------------------------
        # Free ConferenceChair |  2  |   -100.00  | 10% incl | -181.82 | -200.00 |  -18.18
        # Free Pedal Bin       |  5  |   -100.00  | /        | -500.00 | -500.00 |       /
        # Free Large Cabinet   |  3  |   -100.00  | 15% excl | -300.00 | -345.00 |  -45.00
        # --------------------------------------------------------------------------------
        # TOTAL AFTER APPLYING FREE PRODUCT PROGRAMS         | 1594.95 | 1901.11 |  306.16

        self.assertRecordValues(order, [{
            'amount_total': 1901.11,
            'amount_untaxed': 1594.95,
        }])
        self.assertEqual(len(order.order_line.ids), 8, "Order should contains 5 regular product lines and 3 free product lines")

        # Apply 10% on top of everything
        self._apply_promo_code(order, 'test_10pc')

        # Name                 | Qty | price_unit |  Tax     |  HTVA   |   TVAC  |  TVA  |
        # --------------------------------------------------------------------------------
        # 10% on tax 10% incl  |  1  |    -30.00  | 10% incl | -27.27  | -30.00  |   -2.73
        # 10% on no tax        |  1  |    -50.00  | /        | -50.00  | -50.00  |       /
        # 10% on tax 15% excl  |  1  |    -60.00  | 15% excl | -60.00  | -69.00  |   -9.00
        # 10% on tax 35%+50%   |  1  |    -30.00  | 35% incl | -22.22  | -41.11  |  -18.89
        #                                           50% excl
        # --------------------------------------------------------------------------------
        # TOTAL AFTER APPLYING 10% GLOBAL PROGRAM            | 1435.46 | 1711.00 | 275.54

        self.assertRecordValues(order, [{
            'amount_total': 1711.0,
            'amount_untaxed': 1435.45,
        }])
        self.assertEqual(len(order.order_line.ids), 12, "Order should contains 5 regular product lines, 3 free product lines and 4 discount lines (one for every tax)")

        # -- This is a test inside the test
        order.order_line._compute_tax_ids()
        self.assertRecordValues(order, [{
            'amount_total': 1711.0,
            'amount_untaxed': 1435.45,
        }])
        self.assertEqual(len(order.order_line.ids), 12, "Recomputing tax on sale order lines should not change number of order line")
        self._auto_rewards(order, self.all_programs)
        self.assertRecordValues(order, [{
            'amount_total': 1711.0,
            'amount_untaxed': 1435.45,
        }])
        self.assertEqual(len(order.order_line.ids), 12, "Recomputing tax on sale order lines should not change number of order line")
        # -- End test inside the test

        # Now we want to apply a 20% discount only on Large Cabinet
        self.all_programs |= self.env['loyalty.program'].create({
            'name': '20% reduction on Large Cabinet in cart',
            'trigger': 'auto',
            'program_type': 'promotion',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {})],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount': 20,
                'discount_applicability': 'specific',
                'discount_product_ids': self.largeCabinet,
                'required_points': 1,
                'clear_wallet': 1,
            })],
        })
        self._auto_rewards(order, self.all_programs)

        # 20% on large cabinet which are already discounted by 10%
        # Name                 | Qty | price_unit |  Tax     |  HTVA   |   TVAC  |  TVA  |
        # --------------------------------------------------------------------------------
        # 20% on Large Cabinet |  1  |    -72.00  | 15% excl | -72.00  | -82.8  |  -10.80
        # --------------------------------------------------------------------------------
        # TOTAL AFTER APPLYING 20% ON LARGE CABINET          | 1363.46 | 1628.2 |  264.74

        self.assertRecordValues(order, [{
            'amount_total': 1628.2,
            'amount_untaxed': 1363.45,
        }])
        self.assertEqual(len(order.order_line.ids), 13, "Order should have a new discount line for 20% on Large Cabinet")

        # Check that if you delete one of the discount tax line, the others tax lines from the same promotion got deleted as well.
        order.order_line.filtered(lambda l: '10%' in l.name)[0].unlink()
        order._remove_program_from_points(self.p1)
        self.assertEqual(len(order.order_line.ids), 9, "All of the 10% discount line per tax should be removed")
        # At this point, removing the Conference Chair's discount line (split per tax) removed also the others discount lines
        # linked to the same program (eg: other taxes lines). So the coupon got removed from the SO since there were no discount lines left

        # Add back the coupon to continue the test flow
        self._apply_promo_code(order, 'test_10pc')
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(len(order.order_line.ids), 13, "The 10% discount line should be back")

        # Check that if you change a product qty, his discount tax line got updated
        self.p_conference_chair.rule_ids.reward_point_amount = 0.752
        sol2.product_uom_qty = 4
        self._auto_rewards(order, self.all_programs)
        # Name                 | Qty | price_unit |  Tax     |  HTVA   |   TVAC  |  TVA  |
        # --------------------------------------------------------------------------------
        # Large Cabinet        |  4  |    100.00  | 15% excl |  400.00 |  460.00 |   60.00
        # Conference Chair     |  4  |    100.00  | 10% incl |  363.63 |  400.00 |   36.36
        # Pedal Bins           |  5  |    100.00  | /        |  500.00 |  500.00 |       /
        # Drawer Black         |  2  |    100.00  | 15% excl |  200.00 |  230.00 |   30.00
        # Product A            |  3  |    100.00  | 35% incl |  222.22 |  411.11 |  188.89
        #                                           50% excl
        # Free - Large Cabinet |  3  |      0.00  | 15% excl |    0.00 |    0.00 |    0.00
        # Free - Conference Ch |  3  |      0.00  | 10% incl |    0.00 |    0.00 |    0.00
        # Free - Pedal Bins    |  5  |      0.00  | /        |    0.00 |    0.00 |       /
        # 20% on Large Cabinet |  1  |    -80.00  | 15% excl |  -80.00 |  -92.00 |  -12.00
        # 10% on tax 15% excl  |  1  |    -52.00  | 15% excl |  -52.00 |  -59.80 |   -7.80
        # 10% on tax 10% excl  |  1  |    -40.00  | 15% excl |  -36.36 |  -40.00 |   -3.64
        # 10% on no tax        |  1  |    -50.00  | /        |  -50.00 |  -50.00 |       /
        # 10% on tax 35+50%    |  1  |    -30.00  | 35% incl |  -22.22 |  -41.11 |  -18.89
        #                                           50% excl
        # --------------------------------------------------------------------------------
        # TOTAL                                              | 1445.27 | 1718.20 |  272.92

        self.assertEqual(order.amount_untaxed, 1445.27, "The order should have one more paid Conference Chair with 10% incl tax and discounted by 10%")

        # Check that if you remove a product, his reward lines got removed, especially the discount per tax one
        sol2.unlink()
        self._auto_rewards(order, self.all_programs)
        # Name                 | Qty | price_unit |  Tax     |  HTVA   |   TVAC  |  TVA  |
        # --------------------------------------------------------------------------------
        # Pedal Bins           |  5  |    100.00  | /        |  500.00 |  500.00 |       /
        # Large Cabinet        |  4  |    100.00  | 15% excl |  400.00 |  460.00 |   60.00
        # Drawer Black         |  2  |    100.00  | 15% excl |  200.00 |  230.00 |   30.00
        # Product A            |  3  |    100.00  | 35% incl |  222.22 |  411.11 |  188.89
        #                                           50% excl
        # Pedal Bins           |  5  |      0.00  | /        |    0.00 |    0.00 |       /
        # Large Cabinet        |  3  |      0.00  | 15% excl |    0.00 |    0.00 |    0.00
        # 20% on Large Cabinet |  1  |    -80.00  | 15% excl |  -80.00 |  -92.00 |  -12.00
        # 10% on tax 15% excl  |  1  |    -52.00  | 15% excl |  -52.00 |  -59.80 |   -7.80
        # 10% on no tax        |  1  |    -50.00  | /        |  -50.00 |  -50.00 |       /
        # 10% on tax 35+50%    |  1  |    -30.00  | 35% incl |  -22.22 |  -41.11 |  -18.89
        #                                           50% excl
        # --------------------------------------------------------------------------------
        # TOTAL                                              | 1118.00 | 1349.00 |  240.20

        self.assertRecordValues(order, [{
            'amount_total': 1358.2,
            'amount_untaxed': 1118.0,
        }])
        self.assertEqual(len(order.order_line.ids), 10, "Order should contains 10 lines: 4 products lines, 2 free products lines and 4 discount lines")

    def test_program_numbers_extras(self):
        # Check that you can't apply a global discount promo code if there is already an auto applied global discount
        p1_copy = self.p1.copy({'trigger': 'auto', 'name': 'Auto applied 10% global discount', 'rule_ids': [(0, 0, {})]})
        self.all_programs |= p1_copy
        order = self.empty_order
        self.env['sale.order.line'].create({
            'product_id': self.largeCabinet.id,
            'name': 'Large Cabinet',
            'product_uom_qty': 1.0,
            'order_id': order.id,
        })
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(len(order.order_line.ids), 2, "We should get 1 Large Cabinet line and 1 10% auto applied global discount line")
        self.assertEqual(order.amount_total, 288, "320$ - 10%")
        with self.assertRaises(ValidationError):
            # Can't apply a second global discount
            self._apply_promo_code(order, 'test_10pc')

    def test_program_fixed_price(self):
        # Check fixed amount discount
        order = self.empty_order
        self.p3.active = False
        fixed_amount_program = self.env['loyalty.program'].create({
            'name': '$249 discount',
            'trigger': 'auto',
            'program_type': 'promotion',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {
                'reward_point_mode': 'order',
                'reward_point_amount': 1,
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount': 249,
                'discount_mode': 'per_point',
                'discount_applicability': 'order',
                'required_points': 1,
            })],
        })
        self.all_programs |= fixed_amount_program
        self.tax_0pc_excl = self.env['account.tax'].create({
            'name': "0% Tax excl",
            'amount_type': 'percent',
            'amount': 0,
        })
        fixed_amount_program.reward_ids.discount_line_product_id.write({'taxes_id': [(4, self.tax_0pc_excl.id, False)]})
        sol1 = self.env['sale.order.line'].create({
            'product_id': self.drawerBlack.id,
            'name': 'Drawer Black',
            'product_uom_qty': 1.0,
            'order_id': order.id,
            'tax_ids': [(4, self.tax_0pc_excl.id)]
        })
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(order.amount_total, 0, "Total should be null. The fixed amount discount is higher than the SO total, it should be reduced to the SO total")
        self.assertEqual(len(order.order_line.ids), 2, "There should be the product line and the reward line")
        sol1.product_uom_qty = 17
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(order.amount_total, 176, "Fixed amount discount should be totally deduced")
        self.assertEqual(len(order.order_line.ids), 2, "Number of lines should be unchanged as we just recompute the reward line")
        fixed_amount_program.write({'active': False})  # Check archived product will remove discount lines on recompute
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(len(order.order_line.ids), 1, "Archiving the program should remove the program reward line")

    def test_program_next_order(self):
        order = self.empty_order
        self.all_programs |= self.env['loyalty.program'].create({
            'name': 'Free Pedal Bin if at least 1 article',
            'trigger': 'auto',
            'applies_on': 'future',
            'program_type': 'promotion',
            'rule_ids': [(0, 0, {
                'minimum_qty': 2,
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'product',
                'reward_product_id': self.pedalBin.id,
                'reward_product_qty': 1,
                'required_points': 1,
            })],
        })
        sol1 = self.env['sale.order.line'].create({
            'product_id': self.largeCabinet.id,
            'name': 'Large Cabinet',
            'product_uom_qty': 1.0,
            'order_id': order.id,
        })
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(len(order.order_line.ids), 1, "Nothing should be added to the cart")
        self.assertEqual(len(order._get_reward_coupons()), 0, "No coupon should have been generated yet")

        sol1.product_uom_qty = 2
        self._auto_rewards(order, self.all_programs)
        generated_coupon = order._get_reward_coupons()
        self.assertEqual(len(order.order_line.ids), 1, "Nothing should be added to the cart (2)")
        self.assertEqual(len(generated_coupon), 1, "A coupon should have been generated")
        self.assertEqual(generated_coupon.points, 0, "The coupon should not have it's points already.")

        sol1.product_uom_qty = 1
        self._auto_rewards(order, self.all_programs)
        generated_coupon = order._get_reward_coupons()
        self.assertEqual(len(order.order_line.ids), 1, "Nothing should be added to the cart (3)")
        self.assertEqual(len(generated_coupon), 0, "No more coupon should have been generated and the existing one should not have been deleted")

        sol1.product_uom_qty = 2
        self._auto_rewards(order, self.all_programs)
        generated_coupon = order._get_reward_coupons()
        self.assertEqual(len(generated_coupon), 1, "We should still have only 1 coupon as we now benefit again from the program but no need to create a new one (see next assert)")
        self.assertEqual(generated_coupon.points, 0, "The coupon should not have it's points already.")
        self.assertFalse(order._get_claimable_rewards(), "No rewards should be claimable")

        order.action_confirm()
        self.assertEqual(
            generated_coupon.points, 1,
            "The coupon should have 1 point after confirmation",
        )
        self.assertFalse(
            order._get_claimable_rewards(),
            "Next-order coupon rewards shouldn't be claimable on current order",
        )

    def test_coupon_rule_minimum_amount(self):
        """ Ensure coupon with minimum amount rule are correctly
            applied on orders
        """
        order = self.empty_order
        self.env['sale.order.line'].create({
            'product_id': self.conferenceChair.id,
            'name': 'Conference Chair',
            'product_uom_qty': 10.0,
            'order_id': order.id,
        })
        self.assertEqual(order.amount_total, 165.0, "The order amount is not correct")
        self.env['loyalty.generate.wizard'].with_context(active_id=self.discount_coupon_program.id).create({
            'coupon_qty': 1,
            'points_granted': 1,
        }).generate_coupons()
        coupon = self.discount_coupon_program.coupon_ids[0]
        self._apply_promo_code(order, coupon.code)
        self.assertEqual(order.amount_total, 65.0, "The coupon should be correctly applied")
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(order.amount_total, 65.0, "The coupon should not be removed from the order")

    def test_coupon_and_program_discount_fixed_amount(self):
        """ Ensure coupon and program discount both with
            minimum amount rule can cohexists without making
            the order go below 0
        """
        order = self.empty_order
        orderline = self.env['sale.order.line'].create({
            'product_id': self.conferenceChair.id,
            'name': 'Conference Chair',
            'product_uom_qty': 10.0,
            'order_id': order.id,
        })
        self.assertEqual(order.amount_total, 165.0, "The order amount is not correct")

        self.env['loyalty.program'].create({
            'name': '$100 promotion program',
            'program_type': 'promotion',
            'trigger': 'with_code',
            'rule_ids': [(0, 0, {
                'mode': 'with_code',
                'code': 'testpromo',
                'minimum_amount': 100,
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount': 100,
                'discount_mode': 'per_point',
                'discount_applicability': 'order',
            })],
        })

        self._apply_promo_code(order, 'testpromo')
        self.assertEqual(order.amount_total, 65.0, "The promotion program should be correctly applied")
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(order.amount_total, 65.0, "The promotion program should not be removed after recomputation")

        self.env['loyalty.generate.wizard'].with_context(active_id=self.discount_coupon_program.id).create({
            'coupon_qty': 1,
            'points_granted': 1,
        }).generate_coupons()
        coupon = self.discount_coupon_program.coupon_ids[0]
        with self.assertRaises(ValidationError):
            self._apply_promo_code(order, coupon.code)
        orderline.write({'product_uom_qty': 15})
        self._apply_promo_code(order, coupon.code)
        self.assertEqual(order.amount_total, 47.5, "The promotion program should now be correctly applied")

        orderline.write({'product_uom_qty': 5})
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(order.amount_total, 82.5, "The promotion programs should have been removed from the order to avoid negative amount")

    def test_coupon_and_coupon_discount_fixed_amount_tax_excl(self):
        """ Ensure multiple coupon can cohexists without making
            the order go below 0
            * Have an order of 300 (3 lines: 1 tax excl 15%, 2 notax)
            * Apply a coupon A of 10% discount, unconditioned
            * Apply a coupon B of 288.5 discount, unconditioned
            * Order should not go below 0
            * Even applying the coupon in reverse order should yield same result
        """

        self.immediate_promotion_program.active = False
        coupon_program = self.env['loyalty.program'].create({
            'name': '$288.5 coupon',
            'program_type': 'coupons',
            'trigger': 'with_code',
            'applies_on': 'current',
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount_mode': 'per_point',
                'discount': 288.5,
                'discount_applicability': 'order',
                'required_points': 1,
            })],
        })

        order = self.empty_order
        self.env['sale.order.line'].create([
        {
            'product_id': self.conferenceChair.id,
            'name': 'Conference Chair',
            'product_uom_qty': 1.0,
            'price_unit': 100.0,
            'order_id': order.id,
            'tax_ids': [(6, 0, (self.tax_15pc_excl.id,))],
        },
        {
            'product_id': self.pedalBin.id,
            'name': 'Computer Case',
            'product_uom_qty': 1.0,
            'price_unit': 100.0,
            'order_id': order.id,
            'tax_ids': [(6, 0, [])],
        },
        {
            'product_id': self.product_A.id,
            'name': 'Computer Case',
            'product_uom_qty': 1.0,
            'price_unit': 100.0,
            'order_id': order.id,
            'tax_ids': [(6, 0, [])],
        },
        ])

        self._apply_promo_code(order, 'test_10pc')
        self.assertEqual(order.amount_total, 283.5, "The promotion program should be correctly applied")

        self.env['loyalty.generate.wizard'].with_context(active_id=coupon_program.id).create({
            'coupon_qty': 1,
            'points_granted': 1,
        }).generate_coupons()
        coupon = coupon_program.coupon_ids
        self._apply_promo_code(order, coupon.code)
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(order.amount_tax, 0.0)
        self.assertEqual(order.amount_untaxed, 0.0, "The untaxed amount should not go below 0")
        self.assertEqual(order.amount_total, 0.0, "The promotion program should not make the order total go below 0")

        order.order_line[3:].unlink() #remove all coupon
        order._remove_program_from_points(coupon_program)
        order._remove_program_from_points(self.p1)

        self._auto_rewards(order, self.all_programs)
        self.assertEqual(len(order.order_line), 3, "The promotion program should be removed")
        self._apply_promo_code(order, coupon.code)
        self.assertEqual(order.amount_total, 26.5, "The promotion program should be correctly applied")
        self._auto_rewards(order, self.all_programs)
        self._apply_promo_code(order, 'test_10pc')
        self._auto_rewards(order, self.all_programs)
        self.assertAlmostEqual(order.amount_tax, 1.14, 2)
        self.assertEqual(order.amount_untaxed, 22.71)
        self.assertEqual(order.amount_total, 23.85, "The promotion program should not make the order total go below 0be altered after recomputation")
        # It should stay the same after a recompute, order matters
        self._auto_rewards(order, self.all_programs)
        self.assertAlmostEqual(order.amount_tax, 1.14, 2)
        self.assertEqual(order.amount_untaxed, 22.71)
        self.assertEqual(order.amount_total, 23.85, "The promotion program should not make the order total go below 0be altered after recomputation")

    def test_coupon_and_coupon_discount_fixed_amount_tax_incl(self):
        """ Ensure multiple coupon can cohexists without making
            the order go below 0
            * Have an order of 300 (3 lines: 1 tax incl 10%, 2 notax)
            * Apply a coupon A of 10% discount, unconditioned
            * Apply a coupon B of 290 discount, unconditioned
            * Order should not go below 0
            * Even applying the coupon in reverse order should yield same result
        """

        self.immediate_promotion_program.active = False
        coupon_program = self.env['loyalty.program'].create({
            'name': '$290 coupon',
            'program_type': 'coupons',
            'trigger': 'with_code',
            'applies_on': 'current',
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount_mode': 'per_point',
                'discount': 290,
                'discount_applicability': 'order',
                'required_points': 1,
            })],
        })

        order = self.empty_order
        self.env['sale.order.line'].create([
        {
            'product_id': self.conferenceChair.id,
            'name': 'Conference Chair',
            'product_uom_qty': 1.0,
            'price_unit': 100.0,
            'order_id': order.id,
            'tax_ids': [(6, 0, (self.tax_10pc_incl.id,))],
        },
        {
            'product_id': self.pedalBin.id,
            'name': 'Computer Case',
            'product_uom_qty': 1.0,
            'price_unit': 100.0,
            'order_id': order.id,
            'tax_ids': [(6, 0, [])],
        },
        {
            'product_id': self.product_A.id,
            'name': 'Computer Case',
            'product_uom_qty': 1.0,
            'price_unit': 100.0,
            'order_id': order.id,
            'tax_ids': [(6, 0, [])],
        },
        ])

        self._apply_promo_code(order, 'test_10pc')
        self.assertEqual(order.amount_total, 270.0, "The promotion program should be correctly applied")

        self.env['loyalty.generate.wizard'].with_context(active_id=coupon_program.id).create({
            'coupon_qty': 1,
            'points_granted': 1,
        }).generate_coupons()
        coupon = coupon_program.coupon_ids
        self._apply_promo_code(order, coupon.code)
        self.assertEqual(order.amount_total, 0, "The promotion program should not make the order total go below 0")
        self.assertEqual(order.amount_tax, 0)
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(order.amount_total, 0, "The promotion program should not be altered after recomputation")
        self.assertEqual(order.amount_tax, 0)

        order.order_line[3:].unlink() #remove all coupon
        order._remove_program_from_points(coupon_program)
        order._remove_program_from_points(self.p1)

        self._auto_rewards(order, self.all_programs)
        self.assertEqual(len(order.order_line), 3, "The promotion program should be removed")
        self._apply_promo_code(order, coupon.code)
        self.assertEqual(order.amount_total, 10.0, "The promotion program should be correctly applied")
        self._apply_promo_code(order, 'test_10pc')
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(order.amount_total, 9.0, "The promotion program should not make the order total go below 0")
        self.assertEqual(order.amount_tax, 0.27)
        self.assertEqual(order.amount_untaxed, 8.73)
        # It should stay the same after a recompute, order matters
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(order.amount_total, 9.0, "The promotion program should not make the order total go below 0")
        self.assertEqual(order.amount_tax, 0.27)
        self.assertEqual(order.amount_untaxed, 8.73)

    def test_program_discount_on_multiple_specific_products(self):
        """ Ensure a discount on multiple specific products is correctly computed.
            - Simple: Discount must be applied on all the products set on the promotion
            - Advanced: This discount must be split by different taxes
        """
        order = self.empty_order
        self.p3.active = False
        p_specific_products = self.env['loyalty.program'].create({
            'name': '20% reduction on Conference Chair and Drawer Black in cart',
            'program_type': 'promotion',
            'trigger': 'auto',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {})],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount_mode': 'percent',
                'discount': 25,
                'discount_applicability': 'specific',
                'discount_product_ids': [(6, 0, [self.conferenceChair.id, self.drawerBlack.id])],
                'required_points': 1,
            })],
        })
        self.all_programs |= p_specific_products

        self.env['sale.order.line'].create({
            'product_id': self.conferenceChair.id,
            'name': 'Conference Chair',
            'product_uom_qty': 4.0,
            'order_id': order.id,
        })
        sol2 = self.env['sale.order.line'].create({
            'product_id': self.drawerBlack.id,
            'name': 'Drawer Black',
            'product_uom_qty': 2.0,
            'order_id': order.id,
        })

        self._auto_rewards(order, self.all_programs)
        self.assertEqual(len(order.order_line.ids), 3, "Conference Chair + Drawer Black + 20% discount line")
        # Name                 | Qty | price_unit |  Tax     |  HTVA   |   TVAC  |  TVA  |
        # --------------------------------------------------------------------------------
        # Conference Chair     |  4  |     16.50  |       /  |   66.00 |   66.00 |   0.00
        # Drawer Black         |  2  |     25.00  |       /  |   50.00 |   50.00 |   0.00
        # 25% discount         |  1  |    -29.00  |       /  |  -29.00 |  -29.00 |   0.00
        # --------------------------------------------------------------------------------
        # TOTAL                                              |   87.00 |   87.00 |   0.00
        self.assertEqual(order.amount_total, 87.00, "Total should be 87.00, see above comment")

        # remove Drawer Black case from promotion
        p_specific_products.reward_ids.discount_product_ids = [(6, 0, [self.conferenceChair.id])]
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(len(order.order_line.ids), 3, "Should still be Conference Chair + Drawer Black + 20% discount line")
        # Name                 | Qty | price_unit |  Tax     |  HTVA   |   TVAC  |  TVA  |
        # --------------------------------------------------------------------------------
        # Conference Chair     |  4  |     16.50  |       /  |   66.00 |   66.00 |   0.00
        # Drawer Black         |  2  |     25.00  |       /  |   50.00 |   50.00 |   0.00
        # 25% discount         |  1  |    -16.50  |       /  |  -16.50 |  -16.50 |   0.00
        # --------------------------------------------------------------------------------
        # TOTAL                                              |   99.50 |   99.50 |   0.00
        self.assertEqual(order.amount_total, 99.50, "The 12.50 discount from the drawer black should be gone")

        # =========================================================================
        # PART 2: Same flow but with different taxes on products to ensure discount is split per VAT
        # Add back Drawer Black in promotion
        p_specific_products.reward_ids.discount_product_ids = [(6, 0, [self.conferenceChair.id, self.drawerBlack.id])]

        percent_tax = self.env['account.tax'].create({
            'name': "30% Tax",
            'amount_type': 'percent',
            'amount': 30,
            'price_include_override': 'tax_included',
        })
        sol2.tax_ids = percent_tax

        self._auto_rewards(order, self.all_programs)
        self.assertEqual(len(order.order_line.ids), 4, "Conference Chair + Drawer Black + 20% on no TVA product (Conference Chair) + 20% on 15% tva product (Drawer Black)")
        # Name                 | Qty | price_unit |  Tax     |  HTVA   |   TVAC  |  TVA  |
        # --------------------------------------------------------------------------------
        # Conference Chair     |  4  |     16.50  |       /  |   66.00 |   66.00 |   0.00
        # Drawer Black         |  2  |     25.00  | 30% incl |   38.46 |   50.00 |  11.54
        # 25% discount         |  1  |    -16.50  |       /  |  -16.50 |  -16.50 |   0.00
        # 25% discount         |  1  |    -12.50  | 30% incl |   -9.62 |  -12.50 |  -2.88
        # --------------------------------------------------------------------------------
        # TOTAL                                              |   78.35 |   87.00 |   8.66
        self.assertEqual(order.amount_total, 87.00, "Total untaxed should be as per above comment")
        self.assertEqual(order.amount_untaxed, 78.35, "Total with taxes should be as per above comment")

    def test_program_numbers_free_prod_with_min_amount_and_qty_on_same_prod(self):
        # This test focus on giving a free product based on both
        # minimum amount and quantity condition on an
        # auto applied promotion program

        order = self.empty_order
        self.p3.active = False
        self.all_programs |= self.env['loyalty.program'].create({
            'name': 'Buy 2 Chairs, get 1 free',
            'trigger': 'auto',
            'program_type': 'promotion',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {
                'product_ids': self.conferenceChair,
                'reward_point_mode': 'order',
                'minimum_qty': 2,
                'minimum_amount': self.conferenceChair.lst_price * 2,
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'product',
                'reward_product_id': self.conferenceChair.id,
                'reward_product_qty': 1,
                'required_points': 1,
            })],
        })
        sol1 = self.env['sale.order.line'].create({
            'product_id': self.conferenceChair.id,
            'name': 'Conf Chair',
            'product_uom_qty': 1.0,
            'order_id': order.id,
        })
        sol2 = self.env['sale.order.line'].create({
            'product_id': self.drawerBlack.id,
            'name': 'Drawer',
            'product_uom_qty': 1.0,
            'order_id': order.id,
        }) # dummy line

        self._auto_rewards(order, self.all_programs)
        self.assertEqual(len(order.order_line.ids), 2, "The promotion lines should not be applied")
        sol1.write({'product_uom_qty': 2.0})
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(len(order.order_line.ids), 3, "The promotion lines should have been added")
        self.assertEqual(order.amount_total, self.conferenceChair.lst_price * (sol1.product_uom_qty) + self.drawerBlack.lst_price * sol2.product_uom_qty, "The promotion line was not applied to the amount total")
        sol2.unlink()
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(len(order.order_line.ids), 2, "The other product should not affect the promotion")
        self.assertEqual(order.amount_total, self.conferenceChair.lst_price * (sol1.product_uom_qty), "The promotion line was not applied to the amount total")
        sol1.write({'product_uom_qty': 1.0})
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(len(order.order_line.ids), 1, "The promotion lines should have been removed")

    def test_program_step_percentages(self):
        # test step-like percentages increase over amount
        testprod = self.env['product.product'].create({
            'name': 'testprod',
            'lst_price': 118.0,
        })

        self.all_programs |= self.env['loyalty.program'].create({
            'name': '10% discount',
            'trigger': 'auto',
            'program_type': 'promotion',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {
                'reward_point_mode': 'order',
                'minimum_amount': 1500.00,
                'minimum_amount_tax_mode': 'incl',
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount': 10,
                'discount_mode': 'percent',
                'discount_applicability': 'order',
                'required_points': 1,
            })],
        })
        self.all_programs |= self.env['loyalty.program'].create({
            'name': '15% discount',
            'trigger': 'auto',
            'program_type': 'promotion',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {
                'reward_point_mode': 'order',
                'minimum_amount': 1750.00,
                'minimum_amount_tax_mode': 'incl',
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount': 15,
                'discount_mode': 'percent',
                'discount_applicability': 'order',
                'required_points': 1,
            })],
        })
        self.all_programs |= self.env['loyalty.program'].create({
            'name': '20% discount',
            'trigger': 'auto',
            'program_type': 'promotion',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {
                'reward_point_mode': 'order',
                'minimum_amount': 2000.00,
                'minimum_amount_tax_mode': 'incl',
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount': 20,
                'discount_mode': 'percent',
                'discount_applicability': 'order',
                'required_points': 1,
            })],
        })
        self.all_programs |= self.env['loyalty.program'].create({
            'name': '25% discount',
            'trigger': 'auto',
            'program_type': 'promotion',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {
                'reward_point_mode': 'order',
                'minimum_amount': 2500.00,
                'minimum_amount_tax_mode': 'incl',
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount': 25,
                'discount_mode': 'percent',
                'discount_applicability': 'order',
                'required_points': 1,
            })],
        })

        #apply 10%
        order = self.empty_order
        order_line = self.env['sale.order.line'].create({
            'product_id': testprod.id,
            'name': 'testprod',
            'product_uom_qty': 14.0,
            'price_unit': 118.0,
            'order_id': order.id,
            'tax_ids': False,
        })
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(order.amount_total, 1486.80, "10% discount should be applied")
        self.assertEqual(len(order.order_line.ids), 2, "discount should be applied")

        #switch to 15%
        order_line.write({'product_uom_qty': 15})
        self.assertEqual(order.amount_total, 1604.8, "Discount improperly applied")
        self.assertEqual(len(order.order_line.ids), 2, "No discount applied while it should")

        #switch to 20%
        order_line.write({'product_uom_qty': 17})
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(order.amount_total, 1604.8, "Discount improperly applied")
        self.assertEqual(len(order.order_line.ids), 2, "No discount applied while it should")

        #still 20%
        order_line.write({'product_uom_qty': 20})
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(order.amount_total, 1888.0, "Discount improperly applied")
        self.assertEqual(len(order.order_line.ids), 2, "No discount applied while it should")

        #back to 10%
        order_line.write({'product_uom_qty': 14})
        self._auto_rewards(order, self.all_programs)
        self.assertEqual(order.amount_total, 1486.80, "Discount improperly applied")
        self.assertEqual(len(order.order_line.ids), 2, "No discount applied while it should")

    def test_program_free_prods_with_min_qty_and_reward_qty_and_rule(self):
        order = self.empty_order
        coupon_program = self.env['loyalty.program'].create({
            'name': '2 free conference chair if at least 1 large cabinet',
            'trigger': 'with_code',
            'program_type': 'coupons',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {
                'product_ids': self.largeCabinet,
                'reward_point_mode': 'order',
                'minimum_qty': 1,
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount': 100,
                'discount_mode': 'percent',
                'discount_applicability': 'specific',
                'discount_product_ids': self.conferenceChair,
                'discount_max_amount': 200,
                'required_points': 1,
            })],
        })
        # set large cabinet and conference chair prices
        self.largeCabinet.write({'list_price': 500, 'sale_ok': True,})
        self.conferenceChair.write({'list_price': 100, 'sale_ok': True})

        # create SOL
        self.env['sale.order.line'].create({
            'product_id': self.largeCabinet.id,
            'name': 'Large Cabinet',
            'product_uom_qty': 1.0,
            'order_id': order.id,
        })
        sol2 = self.env['sale.order.line'].create({
            'product_id': self.conferenceChair.id,
            'name': 'Conference chair',
            'product_uom_qty': 2.0,
            'order_id': order.id,
        })

        self.assertEqual(len(order.order_line), 2, 'The order must contain 2 order lines since the coupon is not yet applied')
        self.assertEqual(order.amount_total, 700.0, 'The price must be 500.0 since the coupon is not yet applied')

        # generate and apply coupon
        self.env['loyalty.generate.wizard'].with_context(active_id=coupon_program.id).create({
            'coupon_qty': 1,
            'points_granted': 1,
        }).generate_coupons()
        coupon = coupon_program.coupon_ids
        self._apply_promo_code(order, coupon.code)

        # Name                  | Qty | price_unit |  Tax     |  HTVA   |   TVAC  |  TVA  |
        # --------------------------------------------------------------------------------
        # Conference Chair      |  2  |    100.00  | /        |  200.00 |  200.00 |       /
        # Large Cabinet         |  1  |    500.00  | /        |  500.00 |  500.00 |       /
        #
        # Free Conference Chair |  2  |   -100.00  | /        | -200.00 | -200.00 |       /
        # --------------------------------------------------------------------------------
        # TOTAL                                               |  500.00 |  500.00 |       /

        self.assertEqual(len(order.order_line), 3, 'The order must contain 3 order lines including one for free conference chair')
        self.assertEqual(order.amount_total, 500.0, 'The price must be 500.0 since two conference chairs are free')
        self.assertEqual(order.order_line[2].price_total, -200.0, 'The last order line should apply a reduction of 200.0 since there are two conference chairs that cost 100.0 each')

        # prevent user to get illicite discount by decreasing the to 1 the reward product qty after applying the coupon
        sol2.product_uom_qty = 1.0
        self._auto_rewards(order, self.all_programs)

        # in this case user should not have -200.0
        # Name                  | Qty | price_unit |  Tax     |  HTVA   |   TVAC  |  TVA  |
        # --------------------------------------------------------------------------------
        # Conference Chair      |  1  |    100.00  | /        |  100.00 |  100.00 |       /
        # Large Cabine          |  1  |    500.00  | /        |  500.00 |  500.00 |       /
        #
        # Free Conference Chair |  2  |   -100.00  | /        | -200.00 | -200.00 |       /
        # --------------------------------------------------------------------------------
        # TOTAL                                               |  400.00 |  400.00 |       /


        # he should rather have this one
        # Name                  | Qty | price_unit |  Tax     |  HTVA   |   TVAC  |  TVA  |
        # --------------------------------------------------------------------------------
        # Conference Chair      |  1  |    100.00  | /        |  100.00 |  100.00 |       /
        # Large Cabinet         |  1  |    500.00  | /        |  500.00 |  500.00 |       /
        #
        # Free Conference Chair |  1  |   -100.00  | /        | -100.00 | -100.00 |       /
        # --------------------------------------------------------------------------------
        # TOTAL                                               |  500.00 |  500.00 |       /

        self.assertEqual(order.amount_total, 500.0, 'The price must be 500.0 since two conference chairs are free and the user only bought one')
        self.assertEqual(order.order_line[2].price_total, -100.0, 'The last order line should apply a reduction of 100.0 since there is one conference chair that cost 100.0')

    def test_program_free_product_different_than_rule_product_with_multiple_application(self):
        order = self.empty_order

        self.p3.active = False
        self.all_programs |= self.env['loyalty.program'].create({
            'name': 'Buy 1 drawer black, get a free Large Meeting Table',
            'trigger': 'auto',
            'program_type': 'promotion',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {
                'product_ids': self.drawerBlack,
                'reward_point_mode': 'order',
                'minimum_qty': 1,
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount': 100,
                'discount_mode': 'percent',
                'discount_applicability': 'specific',
                'discount_product_ids': self.largeMeetingTable,
                'required_points': 1,
            })],
        })

        self.env['sale.order.line'].create({
            'product_id': self.drawerBlack.id,
            'product_uom_qty': 2.0,
            'order_id': order.id,
        })
        sol_B = self.env['sale.order.line'].create({
            'product_id': self.largeMeetingTable.id,
            'product_uom_qty': 1.0,
            'order_id': order.id,
        })

        self._auto_rewards(order, self.all_programs)

        self.assertEqual(len(order.order_line), 3, 'The order must contain 3 order lines: 1x for Black Drawer, 1x for Large Meeting Table and 1x for free Large Meeting Table')
        self.assertEqual(order.amount_total, self.drawerBlack.list_price * 2, 'The price must be 50.0 since the Large Meeting Table is free: 2*25.00 (Black Drawer) + 1*40000.00 (Large Meeting Table) - 1*40000.00 (free Large Meeting Table)')

        sol_B.product_uom_qty = 2

        self._auto_rewards(order, self.all_programs)

        self.assertEqual(len(order.order_line), 3, 'The order must contain 3 order lines: 1x for Black Drawer, 1x for Large Meeting Table and 1x for free Large Meeting Table')
        self.assertEqual(order.amount_total, self.drawerBlack.list_price * 2, 'The price must be 50.0 since the 2 Large Meeting Table are free: 2*25.00 (Black Drawer) + 2*40000.00 (Large Meeting Table) - 2*40000.00 (free Large Meeting Table)')

    def test_program_modify_reward_line_qty(self):
        order = self.empty_order
        product_F = self.env['product.product'].create({
            'name': 'Product F',
            'list_price': 100,
            'sale_ok': True,
            'taxes_id': [(6, 0, [])],
        })
        self.all_programs |= self.env['loyalty.program'].create({
            'name': '1 Product F = 5$ discount',
            'trigger': 'auto',
            'program_type': 'promotion',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {
                'product_ids': product_F,
                'reward_point_mode': 'order',
                'minimum_qty': 1,
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount': 5,
                'discount_mode': 'per_point',
                'required_points': 1,
            })],
        })

        self.env['sale.order.line'].create({
            'product_id': product_F.id,
            'product_uom_qty': 2.0,
            'order_id': order.id,
        })

        self._auto_rewards(order, self.all_programs)

        self.assertEqual(len(order.order_line), 2, 'The order must contain 2 order lines: 1x Product F and 1x 5$ discount')
        self.assertEqual(order.amount_total, 195.0, 'The price must be 195.0 since there is a 5$ discount and 2x Product F')
        self.assertEqual(sum(order.order_line.filtered(lambda x: x.is_reward_line).mapped('product_uom_qty')), 1, 'The reward line should have a quantity of 1 since Fixed Amount discounts apply only once per Sale Order')

        order.order_line[1].product_uom_qty = 2

        self.assertEqual(len(order.order_line), 2, 'The order must contain 2 order lines: 1x Product F and 1x 5$ discount')
        self.assertEqual(order.amount_total, 190.0, 'The price must be 190.0 since there is now 2x 5$ discount and 2x Product F')
        self.assertEqual(order.order_line.filtered(lambda x: x.is_reward_line).price_unit, -5, 'The discount unit price should still be -5 after the quantity was manually changed')

    def test_program_multi_product_max_discount(self):
        order = self.empty_order
        coupon_program = self.env['loyalty.program'].create({
            'name': "50% off for cheapest product(max $30)",
            'trigger': 'with_code',
            'program_type': 'coupons',
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount': 50,
                'discount_mode': 'percent',
                'discount_applicability': 'cheapest',
                'discount_max_amount': 30,
            })],
        })

        # create SOL
        self.env['sale.order.line'].create({
            'product_id': self.largeCabinet.id,
            'product_uom_qty': 2.0,
            'order_id': order.id,
        })

        # generate and apply coupon
        self.env['loyalty.generate.wizard'].with_context(active_id=coupon_program.id).create({
            'coupon_qty': 1,
            'points_granted': 1,
        }).generate_coupons()

        coupon = coupon_program.coupon_ids
        self._apply_promo_code(order, coupon.code)

        self.assertEqual(len(order.order_line), 2, "The order must contain 2 order lines")
        self.assertEqual(
            order.amount_total, 610.0, "The price must be 610.0 since the max discount is 30"
        )

    def test_specific_discount_product_group(self):
        # Tests the following:
        # 1 program: -5$ on [A, B]
        # 1 program: -10$ on A
        # Order with A (6$) B (4$) C (10$)
        # Apply both coupons -> order total should be 10$
        # Apply a 10% discount -> order total should be 9$
        # Redo the same process but discount first
        product_a, product_b, product_c = self.env['product.product'].create([
            {
                'name': 'Product A',
                'list_price': 6,
                'sale_ok': True,
                'taxes_id': [(6, 0, [])],
            },
            {
                'name': 'Product B',
                'list_price': 4,
                'sale_ok': True,
                'taxes_id': [(6, 0, [])],
            },
            {
                'name': 'Product C',
                'list_price': 10,
                'sale_ok': True,
                'taxes_id': [(6, 0, [])],
            },
        ])
        programs = self.env['loyalty.program'].create([
            {
                'name': '-5 USD on [A, B]',
                'trigger': 'auto',
                'program_type': 'promotion',
                'applies_on': 'current',
                'rule_ids': [(0, 0, {
                })],
                'reward_ids': [(0, 0, {
                    'reward_type': 'discount',
                    'discount': 5,
                    'discount_mode': 'per_point',
                    'discount_applicability': 'specific',
                    'discount_product_ids': product_a | product_b,
                    'required_points': 1,
                })],
            },
            {
                'name': '-10 USD on A',
                'trigger': 'auto',
                'program_type': 'promotion',
                'applies_on': 'current',
                'rule_ids': [(0, 0, {
                })],
                'reward_ids': [(0, 0, {
                    'reward_type': 'discount',
                    'discount': 10,
                    'discount_mode': 'per_point',
                    'discount_applicability': 'specific',
                    'discount_product_ids': product_a,
                    'required_points': 1,
                })],
            },
        ])
        order = self.empty_order
        self.env['sale.order.line'].create([
            {
                'product_id': product_a.id,
                'name': 'Product A',
                'product_uom_qty': 1,
                'order_id': order.id,
            },
            {
                'product_id': product_b.id,
                'name': 'Product B',
                'product_uom_qty': 1,
                'order_id': order.id,
            },
            {
                'product_id': product_c.id,
                'name': 'Product C',
                'product_uom_qty': 1,
                'order_id': order.id,
            },
        ])
        self._auto_rewards(order, programs)
        self.assertEqual(order.amount_total, 10, "The total should be 10$.")
        # Try to apply another 10%
        self._apply_promo_code(order, 'test_10pc')
        self.assertEqual(order.amount_total, 9, "The total should be 9$.")
        # Now the order way around
        order.order_line.filtered('reward_id').unlink()
        self._apply_promo_code(order, 'test_10pc')
        self.assertEqual(order.amount_total, 18, "The total should be 9$.")
        self._auto_rewards(order, programs)
        self.assertEqual(order.amount_total, 9, "The total should be 9$.")

    def test_specific_discount_multiple_taxes(self):
        # Check the following setup
        # Product A 10$ 10% tva excl
        # Product B 10$ 20% tva excl
        # Program A -100% on product A
        # Program B -5$ fixed on both products
        # Applying both programs in a different order should result in a different
        #  outcome since discountable amounts are computed per tax
        # Applying program A before B should yield a better final price
        product_a, product_b = self.env['product.product'].create([
            {
                'name': 'Product A',
                'list_price': 10,
                'sale_ok': True,
                'taxes_id': [(6, 0, [self.tax_10pc_excl.id])],
            },
            {
                'name': 'Product B',
                'list_price': 10,
                'sale_ok': True,
                'taxes_id': [(6, 0, [self.tax_20pc_excl.id])],
            },
        ])
        program_a, program_b = self.env['loyalty.program'].create([
            {
                'name': '-100% on A',
                'trigger': 'auto',
                'program_type': 'promotion',
                'applies_on': 'current',
                'rule_ids': [(0, 0, {
                })],
                'reward_ids': [(0, 0, {
                    'reward_type': 'discount',
                    'discount': 100,
                    'discount_mode': 'percent',
                    'discount_applicability': 'specific',
                    'discount_product_ids': product_a,
                    'required_points': 1,
                })],
            },
            {
                'name': '-5 USD on [A, B]',
                'trigger': 'auto',
                'program_type': 'promotion',
                'applies_on': 'current',
                'rule_ids': [(0, 0, {
                })],
                'reward_ids': [(0, 0, {
                    'reward_type': 'discount',
                    'discount': 5,
                    'discount_mode': 'per_point',
                    'discount_applicability': 'specific',
                    'discount_product_ids': product_a | product_b,
                    'required_points': 1,
                })],
            },
        ])

        order = self.empty_order
        self.env['sale.order.line'].create([
            {
                'product_id': product_a.id,
                'name': 'Product A',
                'product_uom_qty': 1,
                'order_id': order.id,
            },
            {
                'product_id': product_b.id,
                'name': 'Product B',
                'product_uom_qty': 1,
                'order_id': order.id,
            },
        ])
        self._auto_rewards(order, program_a)
        self.assertEqual(order.amount_total, 12, 'Total should be 12$')
        self._auto_rewards(order, program_b)
        self.assertAlmostEqual(order.amount_total, 7, 0, 'Total should be 7$')
        # Now the order way around
        order.order_line.filtered('reward_id').unlink()
        self._auto_rewards(order, program_b)
        self.assertAlmostEqual(order.amount_total, 18, 0, 'Total should be 18$')
        self._auto_rewards(order, program_a)
        # We essentially create a discount of -100% off of an already discounted product
        # (11 - 2.4) = 8.6$ discount ~
        self.assertAlmostEqual(order.amount_total, 9.4, 1, 'Total should be 9.4$')

    def test_fixed_amount_taxes_attribution(self):
        program = self.env['loyalty.program'].create({
            'name': '-5 USD',
            'trigger': 'auto',
            'program_type': 'promotion',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount': 5,
                'discount_mode': 'per_point',
                'discount_applicability': 'order',
                'required_points': 1,
            })],
        })

        order = self.empty_order
        sol = self.env['sale.order.line'].create({
            'product_id': self.drawerBlack.id,
            'price_unit': 10,
            'product_uom_qty': 1.0,
            'order_id': order.id,
        })

        self._auto_rewards(order, program)

        self.assertEqual(order.amount_total, 5, 'Price should be 10$ - 5$(discount) = 5$')
        self.assertEqual(order.amount_tax, 0, 'No taxes are applied yet')

        sol.tax_ids = self.tax_10pc_base_incl
        self._auto_rewards(order, program)

        self.assertEqual(order.amount_total, 5, 'Price should be 10$ - 5$(discount) = 5$')
        self.assertEqual(float_compare(order.amount_tax, 5 / 11, precision_rounding=3), 0, '10% Tax included in 5$')

        sol.tax_ids = self.tax_10pc_excl
        self._auto_rewards(order, program)

        # Value is 5.99 instead of 6 because you cannot have 6 with 10% tax excluded and a precision rounding of 2
        self.assertAlmostEqual(order.amount_total, 6, 1, msg='Price should be 11$ - 5$(discount) = 6$')
        self.assertEqual(float_compare(order.amount_tax, 6 / 11, precision_rounding=3), 0, '10% Tax included in 6$')

        sol.tax_ids = self.tax_20pc_excl
        self._auto_rewards(order, program)

        self.assertEqual(order.amount_total, 7, 'Price should be 12$ - 5$(discount) = 7$')
        self.assertEqual(float_compare(order.amount_tax, 7 / 12, precision_rounding=3), 0, '20% Tax included on 7$')

        sol.tax_ids = self.tax_10pc_base_incl + self.tax_10pc_excl
        self._auto_rewards(order, program)

        self.assertAlmostEqual(order.amount_total, 6, 1, msg='Price should be 11$ - 5$(discount) = 6$')
        self.assertEqual(float_compare(order.amount_tax, 6 / 12, precision_rounding=3), 0, '20% Tax included on 6$')

    def test_fixed_amount_taxes_attribution_multiline(self):

        program = self.env['loyalty.program'].create({
            'name': '-5 USD',
            'trigger': 'auto',
            'program_type': 'promotion',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount': 5,
                'discount_mode': 'per_point',
                'discount_applicability': 'order',
                'required_points': 1,
            })],
        })

        order = self.empty_order
        sol1 = self.env['sale.order.line'].create({
            'product_id': self.drawerBlack.id,
            'price_unit': 10,
            'product_uom_qty': 1.0,
            'order_id': order.id,
        })
        sol2 = self.env['sale.order.line'].create({
            'product_id': self.drawerBlack.id,
            'price_unit': 10,
            'product_uom_qty': 1.0,
            'order_id': order.id,
        })

        self._auto_rewards(order, program)

        self.assertAlmostEqual(order.amount_total, 15, 1, msg='Price should be 20$ - 5$(discount) = 15$')
        self.assertEqual(order.amount_tax, 0, 'No taxes are applied yet')

        sol1.tax_ids = self.tax_10pc_base_incl
        self._auto_rewards(order, program)

        self.assertAlmostEqual(order.amount_total, 15, 1, msg='Price should be 20$ - 5$(discount) = 15$')
        self.assertEqual(float_compare(order.amount_tax, 5 / 11 + 0, precision_rounding=3), 0,
                         '10% Tax included in 5$ in sol1 (highest cost) and 0 in sol2')

        sol2.tax_ids = self.tax_10pc_excl
        self._auto_rewards(order, program)

        self.assertAlmostEqual(order.amount_total, 16, 1, msg='Price should be 21$ - 5$(discount) = 16$')
        # Tax amount = 10% in 10$ + 10% in 11$ - 10% in 5$ (apply on excluded)
        self.assertEqual(float_compare(order.amount_tax, 5 / 11, precision_rounding=3), 0)

        sol2.tax_ids = self.tax_10pc_base_incl + self.tax_10pc_excl
        self._auto_rewards(order, program)

        self.assertAlmostEqual(order.amount_total, 16, 1, msg='Price should be 21$ - 5$(discount) = 16$')
        # Promo apply on line 2 (10% inc + 10% exc)
        # Tax amount = 10% in 10$ + 10% in 10$ + 10% in 11 - 10% in 5$ - 10% in 4.55$ (100/110*5)
        #            = 10/11 + 10/11 + 11/11 - 5/11 - 4.55/11
        #            = 21.45/11
        self.assertEqual(float_compare(order.amount_tax, 21.45 / 11, precision_rounding=3), 0)

        sol3 = self.env['sale.order.line'].create({
            'product_id': self.drawerBlack.id,
            'price_unit': 10,
            'product_uom_qty': 1.0,
            'order_id': order.id,
        })
        sol3.tax_ids = self.tax_10pc_excl
        self._auto_rewards(order, program)

        self.assertAlmostEqual(order.amount_total, 27, 1, msg='Price should be 32$ - 5$(discount) = 27$')
        # Promo apply on line 2 (10% inc + 10% exc)
        # Tax amount = 10% in 10$ + 10% in 10$ + 10% in 11$ + 10% in 11$ - 10% in 5$ - 10% in 4.55$ (100/110*5)
        #            = 10/11 + 10/11 + 11/11 + 11/11 - 5/11 - 4.55/11
        #            = 32.45/11
        self.assertEqual(float_compare(order.amount_tax, 32.45 / 11, precision_rounding=3), 0)

    def test_fixed_amount_with_negative_cost(self):
        program = self.env['loyalty.program'].create({
            'name': '-10 USD',
            'trigger': 'auto',
            'program_type': 'promotion',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount': 10,
                'discount_mode': 'per_point',
                'discount_applicability': 'order',
                'required_points': 1,
            })],
        })

        order = self.empty_order

        sol1 = self.env['sale.order.line'].create({
            'product_id': self.drawerBlack.id,
            'price_unit': 10,
            'product_uom_qty': 1.0,
            'order_id': order.id,
        })

        self.env['sale.order.line'].create({
            'product_id': self.drawerBlack.id,
            'name': 'hand discount',
            'price_unit': -5,
            'product_uom_qty': 1.0,
            'order_id': order.id,
        })

        self._auto_rewards(order, program)

        self.assertEqual(len(order.order_line), 3, 'Promotion should add 1 line')
        self.assertEqual(order.amount_total, 0, '10$ discount should cover the whole price')

        sol1.price_unit = 20
        self._auto_rewards(order, program)

        self.assertEqual(len(order.order_line), 3, 'Promotion should add 1 line')
        self.assertEqual(order.amount_total, 5, '10$ discount should be applied on top of the 15$ original price')

    def test_fixed_amount_change_promo_amount(self):
        program = self.env['loyalty.program'].create({
            'name': '-10 USD',
            'trigger': 'auto',
            'program_type': 'promotion',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount': 10,
                'discount_mode': 'per_point',
                'discount_applicability': 'order',
                'required_points': 1,
            })],
        })

        order = self.empty_order

        self.env['sale.order.line'].create({
            'product_id': self.drawerBlack.id,
            'price_unit': 10,
            'product_uom_qty': 1.0,
            'order_id': order.id,
        })

        self._auto_rewards(order, program)

        self.assertEqual(len(order.order_line), 2, 'Promotion should add 1 line')
        self.assertEqual(order.amount_total, 0, '10$ - 10$(discount) = 0$(total) ')

        program.reward_ids.discount = 5
        self._auto_rewards(order, program)

        self.assertEqual(len(order.order_line), 2, 'Promotion should add 1 line')
        self.assertEqual(order.amount_total, 5, '10$ - 5$(discount) = 5$(total) ')

    def test_fixed_tax_not_affected(self):
        program = self.env['loyalty.program'].create({
            'name': '50% discount',
            'program_type': 'promotion',
            'trigger': 'auto',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount_mode': 'percent',
                'discount': 50,
                'discount_applicability': 'order',
                'required_points': 1,
            })],
        })

        order = self.empty_order
        # Create taxes
        self.tax_15pc_excl = self.env['account.tax'].create({
            'name': "15% Tax excl",
            'amount_type': 'percent',
            'amount': 15,
        })
        self.tax_10_fixed = self.env['account.tax'].create({
            'name': "10$ Fixed tax",
            'amount_type': 'fixed',
            'amount': 10,
        })

        # Set tax and prices on products as neeed for the test
        self.product_A.write({'list_price': 100})
        self.product_A.taxes_id = (self.tax_15pc_excl + self.tax_10_fixed)

        # Add products in order
        self.env['sale.order.line'].create({
            'product_id': self.product_A.id,
            'name': 'product A',
            'product_uom_qty': 1.0,
            'order_id': order.id,
        })

        self._auto_rewards(order, program)

        self.assertEqual(len(order.order_line), 2, 'Promotion should add 1 line')
        self.assertEqual(order.amount_total, 67.5, '100$ + 15% tax + 10$ tax - 50%(discount) = 67.5$(total) ')
        self.assertEqual(order.amount_tax, 17.5, '15% tax + 10$ tax$ - 50%$(discount) = 17.5$(total) ')

    def test_fixed_tax_not_affected_2(self):
        program = self.env['loyalty.program'].create({
            'name': '50$ discount',
            'program_type': 'promotion',
            'trigger': 'auto',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount_mode': 'per_order',
                'discount': 50,
                'discount_applicability': 'order',
                'required_points': 1,
            })],
        })

        order = self.empty_order
        # Create taxes
        self.tax_15pc_excl = self.env['account.tax'].create({
            'name': "15% Tax excl",
            'amount_type': 'percent',
            'amount': 15,
        })
        self.tax_10_fixed = self.env['account.tax'].create({
            'name': "10$ Fixed tax",
            'amount_type': 'fixed',
            'amount': 10,
        })

        # Set tax and prices on products as neeed for the test
        self.product_A.write({'list_price': 100})
        self.product_A.taxes_id = (self.tax_15pc_excl + self.tax_10_fixed)

        # Add products in order
        self.env['sale.order.line'].create({
            'product_id': self.product_A.id,
            'name': 'product A',
            'product_uom_qty': 1.0,
            'order_id': order.id,
        })

        self._auto_rewards(order, program)

        self.assertEqual(len(order.order_line), 2, 'Promotion should add 1 line')
        self.assertEqual(order.amount_total, 75, '100$ + 15% tax + 10$ tax - 50$(discount) = 75$(total) ')

    def test_loyalty_card_tax_total(self):
        loyalty_program = self.env['loyalty.program'].create({
            'name': 'Test loyalty card',
            'program_type': 'loyalty',
            'trigger': 'auto',
            'applies_on': 'both',
            'rule_ids': [Command.create({
                'reward_point_mode': 'money',
                'reward_point_amount': 0.01,
            })],
            'reward_ids': [Command.create({
                'reward_type': 'discount',
                'discount_mode': 'per_point',
                'discount': 1,
                'discount_applicability': 'cheapest',
                'required_points': 1,
            })],
        })
        order = self.empty_order
        self.env['loyalty.card'].create([{
            'program_id': loyalty_program.id,
            'partner_id': order.partner_id.id,
            'points': 3.39,
        }])

        # Create taxes
        tax_15pc_excl = self.env['account.tax'].create({
            'name': "15% Tax excl",
            'amount_type': 'percent',
            'amount': 15,
        })

        # Set tax and prices on products as neeed for the test
        self.product_A.write({
            'list_price': 140.0,
            'taxes_id': [Command.set(tax_15pc_excl.ids)]
        })

        order.order_line = [
            Command.create({
                'product_id': self.product_A.id,
            }),
        ]

        self._auto_rewards(order, loyalty_program)

        self.assertEqual(len(order.order_line), 2, 'Promotion should add 1 line')
        self.assertEqual(order.order_line[0].tax_ids, tax_15pc_excl)
        self.assertEqual(order.order_line[1].tax_ids, tax_15pc_excl)
        self.assertEqual(order.amount_total, 156.0, '140$ + 15% - 5$ = 156$')

    def test_rounded_used_loyalty_points(self):
        """Check that the loyalty points used in a reward are rounded according to the currency."""
        loyalty_program = self.env['loyalty.program'].create({
            'name': 'Test loyalty card',
            'program_type': 'loyalty',
            'trigger': 'auto',
            'applies_on': 'both',
            'rule_ids': [Command.set([])],
            'reward_ids': [Command.create({
                'reward_type': 'discount',
                'discount_mode': 'per_point',
                'discount': 0.03,
                'discount_applicability': 'order',
                'required_points': 1,
            })],
        })
        order = self.empty_order
        self.env['loyalty.card'].create([{
            'program_id': loyalty_program.id,
            'partner_id': order.partner_id.id,
            'points': 3030,
        }])
        product_a = self._create_product(
            name='product_a',
            lst_price=3000.0,
            taxes_id=[Command.set([])],
        )
        order.order_line = [Command.create({'product_id': product_a.id})]

        coupon = loyalty_program.coupon_ids[0]
        order._apply_program_reward(loyalty_program.reward_ids[0], coupon)
        order.action_confirm()
        self.assertEqual(len(order.order_line), 2, 'Promotion should add 1 line')
        used_points = coupon.history_ids[0].used
        self.assertEqual(used_points, coupon.currency_id.round(used_points))

    def test_non_monetary_points_not_rounded_by_currency(self):
        """Check that points from non-monetary programs are not rounded using the currency rounding factor.

        When a loyalty program uses reward_point_mode='order' (dimensionless points, not money),
        its currency's rounding factor must not be applied to point values. A currency with a
        large rounding factor (e.g. 250) would otherwise collapse 1 point to 0, making rewards
        with required_points=1 permanently unclaimable.
        """
        large_rounding_currency = self.env['res.currency'].create({
            'name': 'TST',
            'symbol': 'T',
            'rounding': 250.0,
            'active': True,
        })
        loyalty_program = self.env['loyalty.program'].create({
            'name': 'Test 1 point per order',
            'program_type': 'loyalty',
            'trigger': 'auto',
            'applies_on': 'both',
            'currency_id': large_rounding_currency.id,
            'rule_ids': [Command.create({
                'reward_point_mode': 'order',
                'reward_point_amount': 1,
            })],
            'reward_ids': [Command.create({
                'reward_type': 'discount',
                'discount_mode': 'per_order',
                'discount': 10,
                'discount_applicability': 'order',
                'required_points': 1,
            })],
        })
        order = self.empty_order
        order.order_line = [Command.create({'product_id': self.largeCabinet.id})]
        coupon = self.env['loyalty.card'].create({
            'program_id': loyalty_program.id,
            'partner_id': order.partner_id.id,
            'points': 1,
        })

        points = order._get_real_points_for_coupon(coupon)
        self.assertEqual(points, 1, "1 point from an order-based program must not be rounded to 0 by the currency's rounding factor")

        claimable = order._get_claimable_rewards(forced_coupons=coupon)
        self.assertIn(coupon, claimable, "The reward must be claimable with 1 point")

    def test_rounding_program_application(self):
        """Check that the loyalty program is applied with the currency settings of the order."""
        JPY_company = self.env['res.company'].create({
            'name': 'Test',
            'currency_id': self.env.ref('base.JPY').id,
        })
        USD_pricelist = self.env['product.pricelist'].create({
            'name': 'USD pricelist',
            'company_id': JPY_company.id,
        })
        partner = self.env['res.partner'].create({
            'name': 'Test Partner',
            'company_id': JPY_company.id,
            'property_product_pricelist': USD_pricelist,
        })
        product_a = self._create_product(
            name='product_a',
            lst_price=0.4,
            taxes_id=[Command.set([])],
        )
        loyalty_program = self.env['loyalty.program'].create({
            'name': '10% promotion',
            'program_type': 'promotion',
            'trigger': 'auto',
            'applies_on': 'current',
            'rule_ids': [Command.create({})],
            'reward_ids': [Command.create({
                'reward_type': 'discount',
                'discount_mode': 'percent',
                'discount': 10,
                'discount_applicability': 'order',
                'required_points': 1,
            })],
            'company_id': JPY_company.id
        })
        loyalty_card = self.env['loyalty.card'].create({
            'program_id': loyalty_program.id,
            'partner_id': partner.id,
            'points': 10,
        })

        order = self.env['sale.order'].with_company(JPY_company).create({
            'partner_id': partner.id,
            'pricelist_id': USD_pricelist.id,
            'order_line': [
                Command.create({
                    'product_id': product_a.id,
                }),
            ],
        })

        self.assertEqual(len(order.order_line), 1, 'Promotion line should not be present')

        order._update_programs_and_rewards()
        rewards = order._get_claimable_rewards(loyalty_card)
        applicable_reward = rewards.get(loyalty_card)

        self.assertTrue(applicable_reward, 'Reward should be applicable')
        self.assertTrue(applicable_reward.program_id, 'Promotion should be applicable')
        self.assertEqual(applicable_reward.program_id.id, loyalty_program.id, '10% promotion should be applicable')

        self._claim_reward(order, loyalty_program, loyalty_card)

        self.assertEqual(len(order.order_line), 2, 'Promotion should add 1 line')
        reward_line = order.order_line.filtered(lambda x: x.is_reward_line)
        self.assertTrue(reward_line, 'Promotion should add 1 line')
        self.assertEqual(reward_line.reward_id.program_id.id, loyalty_program.id, 'Reward line should be 10% promotion')
        self.assertEqual(order.reward_amount, -0.04, '10% promotion should be applied')
        self.assertEqual(order.amount_total, 0.36, '10% promotion should be applied')

    def test_apply_order_and_specific_discounts(self):
        """Ensure you can apply a full-order discount, and then a product-specific discount."""
        order_program, specific_program = self.env['loyalty.program'].create([
            {
                'name': "$50 discount",
                'program_type': 'promotion',
                'trigger': 'auto',
                'applies_on': 'current',
                'rule_ids': [Command.create({})],
                'reward_ids': [Command.create({
                    'reward_type': 'discount',
                    'discount_mode': 'per_order',
                    'discount': 50,
                    'discount_applicability': 'order',
                    'required_points': 1,
                })],
            },
            {
                'name': "$10 discount on Pedal Bin",
                'program_type': 'promotion',
                'trigger': 'auto',
                'applies_on': 'current',
                'rule_ids': [Command.create({})],
                'reward_ids': [Command.create({
                    'reward_type': 'discount',
                    'discount_mode': 'per_order',
                    'discount': 10,
                    'discount_applicability': 'specific',
                    'discount_product_ids': self.pedalBin.ids,
                    'required_points': 1,
                })],
            },
        ])
        order = self.empty_order
        order.order_line = [Command.create({
            'product_id': self.pedalBin.id,
            'tax_ids': self.tax_20pc_excl.ids,
        })]

        self.assertAlmostEqual(
            order.amount_total,
            self.pedalBin.list_price * (1 + self.tax_20pc_excl.amount / 100),  # $56.4
            msg="Order total should equal product list price plus taxes",
        )

        self._auto_rewards(order, order_program)
        self.assertAlmostEqual(
            order.amount_total,
            self.pedalBin.list_price * (1 + self.tax_20pc_excl.amount / 100) - 50,  # $6.4
            msg="The order total should be $50 less than initially after the discount is applied.",
        )

        self._auto_rewards(order, specific_program)
        self.assertFalse(
            order.amount_total,
            "Order total should be 0, as a specific discount should have been applied.",
        )
````

### FILE: `odoo_loyalty/upstream/addons/sale_loyalty/tests/test_program_rules.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file36:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty/tests/test_program_rules.py; revision 99edb6dd82b7b560930c00b03b694ba700785370; Git blob ce4e57595c3e04d307b5f04db22ff6746bda1dca"
license: "LGPL-3.0-only"
sha256: "f28fa822aa5d7116816f20f2010c2c3f94f926477de024b768ed52d44c3cbce1"
variables: []
secrets_allowed: false
```

````python
# Part of Odoo. See LICENSE file for full copyright and licensing details.

from datetime import date, timedelta

from freezegun import freeze_time
from pytz import timezone

from odoo.exceptions import ValidationError
from odoo.fields import Command, Datetime

from odoo.addons.payment.tests.common import PaymentCommon
from odoo.addons.sale_loyalty.tests.common import TestSaleCouponCommon


class TestProgramRules(TestSaleCouponCommon, PaymentCommon):
    # Test all the validity rules to allow a customer to have a reward.
    # The check based on the products is already done in the basic operations test

    def test_program_rules_minimum_purchased_amount(self):
        # Test case: Based on the minimum purchased

        self.immediate_promotion_program.rule_ids.write({
            'product_ids': False,
            'minimum_amount': 1006,
            'minimum_amount_tax_mode': 'excl'
        })

        order = self.empty_order
        order.write({'order_line': [
            (0, False, {
                'product_id': self.product_A.id,
                'name': '1 Product A',
                'product_uom_qty': 1.0,
            }),
            (0, False, {
                'product_id': self.product_B.id,
                'name': '2 Product B',
                'product_uom_qty': 1.0,
            })
        ]})
        order._update_programs_and_rewards()
        self._claim_reward(order, self.immediate_promotion_program)
        self.assertEqual(len(order.order_line.ids), 2, "The promo offer shouldn't have been applied as the purchased amount is not enough")

        order = self.env['sale.order'].create({'partner_id': self.partner.id})
        order.write({'order_line': [
            (0, False, {
                'product_id': self.product_A.id,
                'name': '10 Product A',
                'product_uom_qty': 10.0,
            }),
            (0, False, {
                'product_id': self.product_B.id,
                'name': '2 Product B',
                'product_uom_qty': 1.0,
            })
        ]})
        order._update_programs_and_rewards()
        self._claim_reward(order, self.immediate_promotion_program)
        # 10*100 + 5 = 1005
        self.assertEqual(len(order.order_line.ids), 2, "The promo offer should not be applied as the purchased amount is not enough")

        self.immediate_promotion_program.rule_ids.minimum_amount = 1005
        order._update_programs_and_rewards()
        self._claim_reward(order, self.immediate_promotion_program)
        self.assertEqual(len(order.order_line.ids), 3, "The promo offer should be applied as the purchased amount is now enough")

        # 10*(100*1.15) + (5*1.15) = 10*115 + 5.75 = 1155.75
        self.immediate_promotion_program.rule_ids.minimum_amount = 1006
        self.immediate_promotion_program.rule_ids.minimum_amount_tax_mode = 'incl'
        order._update_programs_and_rewards()
        self._claim_reward(order, self.immediate_promotion_program)
        self.assertEqual(len(order.order_line.ids), 3, "The promo offer should be applied as the initial amount required is now tax included")

    def test_program_rules_min_amount_not_reached_and_specific_product(self):
        """
        Test that the discount isn't applied if the min amount isn't reached for the specified
        product.
        """
        self.env['loyalty.program'].search([]).active = False
        order = self.empty_order
        program = self.env['loyalty.program'].create({
            'name': "Discount on Product A",
            'program_type': 'promotion',
            'trigger': 'auto',
            'applies_on': 'current',
            'rule_ids': [Command.create({
                'minimum_amount': 110,
                'minimum_amount_tax_mode': 'excl',
                'product_ids': [Command.set(self.product_A.ids)],
            })],
            'reward_ids': [Command.create({
                'reward_type': 'discount',
                'discount': 10,
                'discount_mode': 'percent',
                'discount_applicability': 'specific',
                'discount_product_ids': [Command.set(self.product_A.ids)],
            })],
        })
        self.env['sale.order.line'].create([{
            'product_id': self.product_A.id,
            'product_uom_qty': 1.0,
            'order_id': order.id,
        }, {
            'product_id': self.product_B.id,
            'product_uom_qty': 40.0,
            'order_id': order.id,
        }])
        self.assertEqual(len(order.order_line), 2)
        self.assertEqual(order.amount_untaxed, 300)

        order._update_programs_and_rewards()
        self._claim_reward(order, program)

        self.assertEqual(len(order.order_line), 2)
        self.assertEqual(order.amount_untaxed, 300)

    def test_program_rules_min_amount_reached_and_specific_product(self):
        """
        Test that the discount is applied if the min amount is reached for the specified product.
        """
        self.env['loyalty.program'].search([]).active = False
        order = self.empty_order
        program = self.env['loyalty.program'].create({
            'name': "Discount on Product A",
            'program_type': 'promotion',
            'trigger': 'auto',
            'applies_on': 'current',
            'rule_ids': [Command.create({
                'minimum_amount': 110,
                'minimum_amount_tax_mode': 'excl',
                'product_ids': [Command.set(self.product_A.ids)],
            })],
            'reward_ids': [Command.create({
                'reward_type': 'discount',
                'discount': 10,
                'discount_mode': 'percent',
                'discount_applicability': 'specific',
                'discount_product_ids': [Command.set(self.product_A.ids)],
            })],
        })
        self.env['sale.order.line'].create([{
            'product_id': self.product_A.id,
            'product_uom_qty': 2.0,
            'order_id': order.id,
        }, {
            'product_id': self.product_B.id,
            'product_uom_qty': 20.0,
            'order_id': order.id,
        }])
        self.assertEqual(len(order.order_line), 2)
        self.assertEqual(order.amount_untaxed, 300)

        order._update_programs_and_rewards()
        self._claim_reward(order, program)

        self.assertEqual(len(order.order_line), 3)
        self.assertEqual(order.amount_untaxed, 280)

    def test_program_rules_coupon_qty_and_amount_remove_not_eligible(self):
        ''' This test will:
                * Check quantity and amount requirements works as expected (since it's slightly different from a promotion_program)
                * Ensure that if a reward from a coupon_program was allowed and the conditions are not met anymore,
                  the reward will be removed on recompute.
        '''
        self.immediate_promotion_program.active = False  # Avoid having this program to add rewards on this test
        order = self.empty_order

        program = self.env['loyalty.program'].create({
            'name': 'Get 10% discount if buy at least 4 Product A and $320',
            'program_type': 'coupons',
            'applies_on': 'current',
            'trigger': 'with_code',
            'rule_ids': [(0, 0, {
                'product_ids': self.product_A,
                'minimum_qty': 3,
                'minimum_amount': 320,
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount_mode': 'percent',
                'discount': 10,
                'discount_applicability': 'order',
            })],
        })

        sol1 = self.env['sale.order.line'].create({
            'product_id': self.product_A.id,
            'name': 'Product A',
            'product_uom_qty': 2.0,
            'order_id': order.id,
        })

        sol2 = self.env['sale.order.line'].create({
            'product_id': self.product_B.id,
            'name': 'Product B',
            'product_uom_qty': 4.0,
            'order_id': order.id,
        })

        # Default value for coupon generate wizard is generate by quantity and generate only one coupon
        self.env['loyalty.generate.wizard'].with_context(active_id=program.id).create({'coupon_qty': 1, 'points_granted': 1}).generate_coupons()
        coupon = program.coupon_ids[0]

        # Not enough amount since we only have 220 (100*2 + 5*4)
        with self.assertRaises(ValidationError):
            self._apply_promo_code(order, coupon.code)

        sol2.product_uom_qty = 24

        # Not enough qty since we only have 3 Product A (Amount is ok: 100*2 + 5*24 = 320)
        with self.assertRaises(ValidationError):
            self._apply_promo_code(order, coupon.code)

        sol1.product_uom_qty = 3

        self._apply_promo_code(order, coupon.code)
        self._claim_reward(order, program, coupon)

        self.assertEqual(len(order.order_line.ids), 3, "The order should contain the Product A line, the Product B line and the discount line")

        sol1.product_uom_qty = 2
        order._update_programs_and_rewards()

        self.assertEqual(len(order.order_line.ids), 2, "The discount line should have been removed as we don't meet the program requirements")

    def test_program_rules_promotion_use_best(self):
        ''' This test verifies that only the best global discount is applied.
        '''
        self.immediate_promotion_program.active = False  # Avoid having this program to add rewards on this test
        order = self.empty_order

        p1 = self.env['loyalty.program'].create({
            'name': 'Get 5% discount if buy at least 2 Product',
            'trigger': 'auto',
            'program_type': 'promotion',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {
                'reward_point_mode': 'order',
                'minimum_qty': 2,
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount': 5,
                'discount_mode': 'percent',
                'discount_applicability': 'order',
                'required_points': 1,
            })],
        })
        p2 = self.env['loyalty.program'].create({
            'name': 'Get 10% discount if buy at least 4 Product',
            'trigger': 'auto',
            'program_type': 'promotion',
            'applies_on': 'current',
            'rule_ids': [(0, 0, {
                'reward_point_mode': 'order',
                'minimum_qty': 4,
            })],
            'reward_ids': [(0, 0, {
                'reward_type': 'discount',
                'discount': 10,
                'discount_mode': 'percent',
                'discount_applicability': 'order',
                'required_points': 1,
            })],
        })
        sol = self.env['sale.order.line'].create({
            'product_id': self.product_A.id,
            'name': 'Product A',
            'product_uom_qty': 1.0,
            'order_id': order.id,
        })

        order._update_programs_and_rewards()
        self._claim_reward(order, p1)
        self._claim_reward(order, p2)
        self.assertEqual(len(order.order_line.ids), 1, "The order should only contains the Product A line")

        sol.product_uom_qty = 3
        order._update_programs_and_rewards()
        self._claim_reward(order, p1)
        self._claim_reward(order, p2)
        discounts = set(order.order_line.mapped('name')) - {'Product A'}
        self.assertEqual(len(discounts), 1, "The order should contains the Product A line and a discount")
        # The name of the discount is dynamically changed to smth looking like:
        # "Discount Get 5% discount if buy at least 2 Product - On product with following tax: Tax 15.00%"
        self.assertTrue('Discount 5% on your order' in discounts.pop(), "The discount should be a 5% discount")

        sol.product_uom_qty = 5
        order._update_programs_and_rewards()
        self._claim_reward(order, p1)
        self._claim_reward(order, p2)
        discounts = set(order.order_line.mapped('name')) - {'Product A'}
        self.assertEqual(len(discounts), 1, "The order should contains the Product A line and a discount")
        self.assertTrue('Discount 10% on your order' in discounts.pop(), "The discount should be a 10% discount")

    @freeze_time('2011-11-02 09:00:21')
    def test_program_rules_validity_dates(self):
        # Test date_to (no date_from)
        today = date.today()
        past_day = today - timedelta(days=2)
        future_day = today + timedelta(days=2)
        self.immediate_promotion_program.write({'date_to': past_day})
        order = self.empty_order
        order.write({'order_line': [
            Command.create({
                'product_id': self.product_A.id,
                'name': '1 Product A',
                'product_uom_qty': 1.0,
            }),
            Command.create({
                'product_id': self.product_B.id,
                'name': '2 Product B',
                'product_uom_qty': 1.0,
            })
        ]})
        self._auto_rewards(order, self.immediate_promotion_program)
        msg = "The promo shouldn't have been applied as it is expired."
        self.assertEqual(len(order.order_line.ids), 2, msg)

        self.immediate_promotion_program.write({'date_to': future_day})
        self._auto_rewards(order, self.immediate_promotion_program)
        msg = "The promo should have been applied we're between the validity dates."
        self.assertEqual(len(order.order_line.ids), 3, msg)

        # Test date_from (no date_to)
        self.immediate_promotion_program.write({
            'date_from': future_day, 'date_to': False,
        })
        self._auto_rewards(order, self.immediate_promotion_program)
        msg = "The promo shouldn't have been applied as it is not active yet."
        self.assertEqual(len(order.order_line.ids), 2, msg)

        self.immediate_promotion_program.write({'date_from': past_day})
        self._auto_rewards(order, self.immediate_promotion_program)
        msg = "The promo should have been applied we're between the validity dates."
        self.assertEqual(len(order.order_line.ids), 3, msg)

        # Test date_from and date_to
        self.immediate_promotion_program.write({'date_from': past_day, 'date_to': future_day})
        self._auto_rewards(order, self.immediate_promotion_program)
        msg = "The promo should have been applied as we're between the validity dates"
        self.assertEqual(len(order.order_line.ids), 3, msg)

        self.immediate_promotion_program.write({
            'date_from': today + timedelta(days=1),
            'date_to': future_day,
        })
        self._auto_rewards(order, self.immediate_promotion_program)
        msg = "The promo offer shouldn't have been applied as it is not active yet."
        self.assertEqual(len(order.order_line.ids), 2, msg)

        self.immediate_promotion_program.write({
            'date_from': past_day,
            'date_to': today - timedelta(days=1),
        })
        self._auto_rewards(order, self.immediate_promotion_program)
        msg = "The promo offer shouldn't have been applied as it is expired."
        self.assertEqual(len(order.order_line.ids), 2, msg)

        self.immediate_promotion_program.write({'date_from': today, 'date_to': today})
        self._auto_rewards(order, self.immediate_promotion_program)
        msg = "The promo should have been applied as today is a valid starting and ending date."
        self.assertEqual(len(order.order_line.ids), 3, msg)

    def test_program_rules_number_of_uses(self):
        # Test case: Based on the number of allowed uses
        self.immediate_promotion_program.write({
            'limit_usage': True,
            'max_usage': 1,
        })
        order = self.empty_order
        order.write({'order_line': [
            Command.create({
                'product_id': self.product_A.id,
                'name': '1 Product A',
                'product_uom_qty': 1.0,
            })
        ]})
        self._auto_rewards(order, self.immediate_promotion_program)
        self.assertEqual(len(order.order_line.ids), 2, "The promo offer should have been applied")

        order = self.env['sale.order'].create({
            'partner_id': self.env['res.partner'].create({'name': 'My Partner'}).id
        })

        order.write({'order_line': [
            Command.create({
                'product_id': self.product_B.id,
                'name': '2 Product B',
                'product_uom_qty': 1.0,
            })
        ]})
        # Invalidate total_order_count
        self.immediate_promotion_program.invalidate_recordset(['order_count', 'total_order_count'])
        self._auto_rewards(order, self.immediate_promotion_program)
        msg = "The promo offer shouldn't have been applied as the number of uses is exceeded"
        self.assertEqual(len(order.order_line.ids), 1, msg)

    def test_program_rules_validity_date_timezones(self):
        """Test that the validity dates are checked according to the company's time zone"""
        self.env.company.partner_id.tz = 'Europe/London'
        self.partner.tz = 'America/Los_Angeles'
        midnight = Datetime.today()
        yesterday = (midnight - timedelta(days=1)).date()
        self.immediate_promotion_program.update({
            'date_to': yesterday,
            'limit_usage': True,
            'max_usage': 1,
        })
        order = self.empty_order.with_context(tz=self.partner.tz)
        order.order_line = [
            Command.create({'product_id': self.product_A.id}),
            Command.create({'product_id': self.product_B.id}),
        ]

        with freeze_time(midnight):
            # Try apply reward at UTC midnight with LA time zone in context (expired)
            self._auto_rewards(order, self.immediate_promotion_program)
            self.assertFalse(
                order.order_line.filtered('is_reward_line'),
                "Promo should not be applied if only valid in the customer's time zone",
            )

        with freeze_time(timezone(self.env.company.partner_id.tz).localize(midnight)):
            # Try apply reward at London midnight (expired)
            self._auto_rewards(order, self.immediate_promotion_program)
            self.assertFalse(
                order.order_line.filtered('is_reward_line'),
                "Promo should not be applied if only valid in the customer's time zone",
            )

        self.partner.tz = 'Europe/Brussels'
        with freeze_time(timezone(self.partner.tz).localize(midnight)):
            # Apply reward at Brussels midnight (still valid in company's time zone)
            self._auto_rewards(order, self.immediate_promotion_program)
            self.assertTrue(
                order.order_line.filtered('is_reward_line'),
                "Promo should be applied if valid in the company's time zone",
            )

    def test_program_rules_validity_date_transactions(self):
        """Test that the validity dates are checked according to the time of transaction."""
        today = Datetime.today()
        tomorrow = today + timedelta(days=1)
        self.immediate_promotion_program.update({
            'date_to': today,
            'limit_usage': True,
            'max_usage': 1,
            'reward_ids': [Command.set(self.program_gift_card.reward_ids.ids)],
        })
        order = self.empty_order
        order.order_line = [
            Command.create({'product_id': self.product_A.id}),
            Command.create({'product_id': self.product_B.id}),
        ]

        intial_amount = order.amount_total
        self._auto_rewards(order, self.immediate_promotion_program)
        self.assertLess(order.amount_total, intial_amount, "A discount should be applied")

        tx = order.transaction_ids = self._create_transaction(
            flow='redirect',
            sale_order_ids=[order.id],
            state='pending',
            reference=order.name,
            amount=order.amount_total,
        )
        # Our slow provider only gets around to confirming the transaction the next day
        with freeze_time(tomorrow):
            tx._set_done()
            tx._post_process()
            self.assertAlmostEqual(
                order.amount_total, tx.amount,
                msg="Discount should still apply if transaction gets confirmed post-expiration",
            )

    def test_buy_x_get_y_free_applies_correctly_with_non_unit_uom(self):
        buy_x_get_y = self.env['loyalty.program'].create({
            'name': 'Buy 12 Take 6',
            'program_type': 'buy_x_get_y',
            'trigger': 'auto',
            'applies_on': 'current',
            'rule_ids': [Command.create({
                'reward_point_mode': 'unit',
                'product_ids': self.product_A.ids,
            })],
            'reward_ids': [Command.create({
                'reward_type': 'product',
                'reward_product_id': self.product_A.id,
                'required_points': 12,
                'reward_product_qty': 6,
            })],
        })
        order = self.empty_order
        order.order_line = [
            Command.create({
                'product_id': self.product_A.id,
                'product_uom_id': self.ref('uom.product_uom_dozen'),
                'product_uom_qty': 1,
            }),
        ]
        self._auto_rewards(order, buy_x_get_y)
        reward_line = order.order_line.filtered('is_reward_line')
        self.assertTrue(reward_line)
        self.assertEqual(reward_line.product_uom_qty, 6)
````

### FILE: `odoo_loyalty/upstream/odoo/tools/float_utils.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file37:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/odoo/odoo/99edb6dd82b7b560930c00b03b694ba700785370/odoo/tools/float_utils.py; revision 99edb6dd82b7b560930c00b03b694ba700785370; Git blob 123a8b0409a64f48e5b45875c3420bff0f32d657"
license: "LGPL-3.0-only"
sha256: "ab7d51a44033792414dc63f598517ee76b922f76b2de9fa3eeb977ab78d7a962"
variables: []
secrets_allowed: false
```

````python
# Part of Odoo. See LICENSE file for full copyright and licensing details.

from typing import Literal, overload

import builtins
import math

RoundingMethod = Literal['UP', 'DOWN', 'HALF-UP', 'HALF-DOWN', 'HALF-EVEN']

__all__ = [
    "float_compare",
    "float_div",
    "float_is_zero",
    "float_repr",
    "float_round",
    "float_split",
    "float_split_str",
]


def round(f: float) -> float:
    # P3's builtin round differs from P2 in the following manner:
    # * it rounds half to even rather than up (away from 0)
    # * round(-0.) loses the sign (it returns -0 rather than 0)
    # * round(x) returns an int rather than a float
    #
    # this compatibility shim implements Python 2's round in terms of
    # Python 3's so that important rounding error under P3 can be
    # trivially fixed, assuming the P2 behaviour to be debugged and
    # correct.
    roundf = builtins.round(f)
    if builtins.round(f + 1) - roundf != 1:
        return f + math.copysign(0.5, f)
    # copysign ensures round(-0.) -> -0 *and* result is a float
    return math.copysign(roundf, f)


def _float_check_precision(
    precision_digits: int | None = None,
    precision_rounding: float | None = None,
) -> float:
    if precision_rounding is not None and precision_digits is None:
        assert precision_rounding > 0,\
            f"precision_rounding must be positive, got {precision_rounding}"
    elif precision_digits is not None and precision_rounding is None:
        # TODO: `int`s will also get the `is_integer` method starting from python 3.12
        assert float(precision_digits).is_integer() and precision_digits >= 0,\
            f"precision_digits must be a non-negative integer, got {precision_digits}"
        precision_rounding = 10 ** -precision_digits
    else:
        msg = "exactly one of precision_digits and precision_rounding must be specified"
        raise AssertionError(msg)
    return precision_rounding


@overload
def float_round(
    value: float,
    precision_digits: int,
    rounding_method: RoundingMethod = ...,
) -> float: ...


@overload
def float_round(
    value: float,
    precision_rounding: float,
    rounding_method: RoundingMethod = ...,
) -> float: ...


def float_round(
    value: float,
    precision_digits: int | None = None,
    precision_rounding: float | None = None,
    rounding_method: RoundingMethod = 'HALF-UP',
) -> float:
    """Return ``value`` rounded to ``precision_digits`` decimal digits,
       minimizing IEEE-754 floating point representation errors, and applying
       the tie-breaking rule selected with ``rounding_method``, by default
       HALF-UP (away from zero).
       Precision must be given by ``precision_digits`` or ``precision_rounding``,
       not both!

       :param value: the value to round
       :param precision_digits: number of fractional digits to round to.
       :param precision_rounding: decimal number representing the minimum
           non-zero value at the desired precision (for example, 0.01 for a
           2-digit precision).
       :param rounding_method: the rounding method used:
           - 'HALF-UP' will round to the closest number with ties going away from zero.
           - 'HALF-DOWN' will round to the closest number with ties going towards zero.
           - 'HALF-EVEN' will round to the closest number with ties going to the closest
              even number.
           - 'UP' will always round away from 0.
           - 'DOWN' will always round towards 0.
       :return: rounded float
    """
    rounding_factor = _float_check_precision(precision_digits=precision_digits,
                                             precision_rounding=precision_rounding)
    if rounding_factor == 0 or value == 0:
        return 0.0

    # NORMALIZE - ROUND - DENORMALIZE
    # In order to easily support rounding to arbitrary 'steps' (e.g. coin values),
    # we normalize the value before rounding it as an integer, and de-normalize
    # after rounding: e.g. float_round(1.3, precision_rounding=.5) == 1.5
    def normalize(val):
        return val / rounding_factor

    def denormalize(val):
        return val * rounding_factor

    # inverting small rounding factors reduces rounding errors
    if rounding_factor < 1:
        rounding_factor = float_invert(rounding_factor)
        normalize, denormalize = denormalize, normalize

    normalized_value = normalize(value)

    # Due to IEEE-754 float/double representation limits, the approximation of the
    # real value may be slightly below the tie limit, resulting in an error of
    # 1 unit in the last place (ulp) after rounding.
    # For example 2.675 == 2.6749999999999998.
    # To correct this, we add a very small epsilon value, scaled to the
    # the order of magnitude of the value, to tip the tie-break in the right
    # direction.
    # Credit: discussion with OpenERP community members on bug 882036
    epsilon_magnitude = math.log2(abs(normalized_value))
    # `2**(epsilon_magnitude - 52)` would be the minimal size, but we increase it to be
    # more tolerant of inaccuracies accumulated after multiple floating point operations
    epsilon = 2**(epsilon_magnitude - 50)

    match rounding_method:
        case 'HALF-UP':  # 0.5 rounds away from 0
            result = round(normalized_value + math.copysign(epsilon, normalized_value))
        case 'HALF-EVEN':  # 0.5 rounds towards closest even number
            integral = math.floor(normalized_value)
            remainder = abs(normalized_value - integral)
            is_half = abs(0.5 - remainder) < epsilon
            # if is_half & integral is odd, add odd bit to make it even
            result = integral + (integral & 1) if is_half else round(normalized_value)
        case 'HALF-DOWN':  # 0.5 rounds towards 0
            result = round(normalized_value - math.copysign(epsilon, normalized_value))
        case 'UP':  # round to number furthest from zero
            result = math.trunc(normalized_value + math.copysign(1 - epsilon, normalized_value))
        case 'DOWN':  # round to number closest to zero
            result = math.trunc(normalized_value + math.copysign(epsilon, normalized_value))
        case _:
            msg = f"unknown rounding method: {rounding_method}"
            raise ValueError(msg)

    return denormalize(result)


@overload
def float_is_zero(
    value: float,
    precision_digits: int,
) -> bool: ...


@overload
def float_is_zero(
    value: float,
    precision_rounding: float,
) -> bool: ...


def float_is_zero(
    value: float,
    precision_digits: int | None = None,
    precision_rounding: float | None = None,
) -> bool:
    """Returns true if ``value`` is small enough to be treated as
       zero at the given precision (smaller than the corresponding *epsilon*).
       The precision (``10**-precision_digits`` or ``precision_rounding``)
       is used as the zero *epsilon*: values less than that are considered
       to be zero.
       Precision must be given by ``precision_digits`` or ``precision_rounding``,
       not both!

       Warning: ``float_is_zero(value1-value2)`` is not equivalent to
       ``float_compare(value1,value2) == 0``, as the former will round after
       computing the difference, while the latter will round before, giving
       different results for e.g. 0.006 and 0.002 at 2 digits precision.

       :param precision_digits: number of fractional digits to round to.
       :param precision_rounding: decimal number representing the minimum
           non-zero value at the desired precision (for example, 0.01 for a
           2-digit precision).
       :param value: value to compare with the precision's zero
       :return: True if ``value`` is considered zero
    """
    epsilon = _float_check_precision(precision_digits=precision_digits,
                                     precision_rounding=precision_rounding)
    return value == 0.0 or abs(float_round(value, precision_rounding=epsilon)) < epsilon


@overload
def float_compare(
    value1: float,
    value2: float,
    precision_digits: int,
) -> Literal[-1, 0, 1]: ...


@overload
def float_compare(
    value1: float,
    value2: float,
    precision_rounding: float,
) -> Literal[-1, 0, 1]: ...


def float_compare(
    value1: float,
    value2: float,
    precision_digits: int | None = None,
    precision_rounding: float | None = None,
) -> Literal[-1, 0, 1]:
    """Compare ``value1`` and ``value2`` after rounding them according to the
       given precision. A value is considered lower/greater than another value
       if their rounded value is different. This is not the same as having a
       non-zero difference!
       Precision must be given by ``precision_digits`` or ``precision_rounding``,
       not both!

       Example: 1.432 and 1.431 are equal at 2 digits precision,
       so this method would return 0
       However 0.006 and 0.002 are considered different (this method returns 1)
       because they respectively round to 0.01 and 0.0, even though
       0.006-0.002 = 0.004 which would be considered zero at 2 digits precision.

       Warning: ``float_is_zero(value1-value2)`` is not equivalent to
       ``float_compare(value1,value2) == 0``, as the former will round after
       computing the difference, while the latter will round before, giving
       different results for e.g. 0.006 and 0.002 at 2 digits precision.

       :param value1: first value to compare
       :param value2: second value to compare
       :param precision_digits: number of fractional digits to round to.
       :param precision_rounding: decimal number representing the minimum
           non-zero value at the desired precision (for example, 0.01 for a
           2-digit precision).
       :return: (resp.) -1, 0 or 1, if ``value1`` is (resp.) lower than,
           equal to, or greater than ``value2``, at the given precision.
    """
    rounding_factor = _float_check_precision(precision_digits=precision_digits,
                                             precision_rounding=precision_rounding)
    # equal numbers round equally, so we can skip that step
    # doing this after _float_check_precision to validate parameters first
    if value1 == value2:
        return 0
    value1 = float_round(value1, precision_rounding=rounding_factor)
    value2 = float_round(value2, precision_rounding=rounding_factor)
    delta = value1 - value2
    if float_is_zero(delta, precision_rounding=rounding_factor):
        return 0
    return -1 if delta < 0.0 else 1


@overload
def float_div(
    value1: float,
    value2: float,
    precision_digits: int,
) -> tuple[int, float]: ...


@overload
def float_div(
    value1: float,
    value2: float,
    precision_rounding: float,
) -> tuple[int, float]: ...


def float_div(
    value1: float,
    value2: float,
    precision_digits: int | None = None,
    precision_rounding: float | None = None,
) -> tuple[int, float]:
    """Return the euclidean division of ``value1`` by ``value2`` as the tuple
       ``(quotient, remainder)``, computed at the given precision so that the
       result is free of the IEEE-754 representation errors that affect the
       native ``int(value1 / value2)`` and ``value1 % value2`` operations.
       Both operands are rounded onto the precision grid and scaled to integers
       before dividing, so that ``value1`` rounds to ``quotient * value2 +
       remainder`` at the given precision, with ``quotient`` an integer.
       Precision must be given by ``precision_digits`` or ``precision_rounding``,
       not both!

       :param value1: the dividend
       :param value2: the divisor
       :param precision_digits: number of fractional digits to round to.
       :param precision_rounding: decimal number representing the minimum
           non-zero value at the desired precision (for example, 0.01 for a
           2-digit precision).
       :return: the tuple ``(quotient, remainder)``
    """
    rounding = _float_check_precision(precision_digits=precision_digits,
                                      precision_rounding=precision_rounding)
    # snap both operands onto the precision grid and scale them to exact
    # integers, so the euclidean division carries no representation error
    scaled_value1 = builtins.round(float_round(value1, precision_rounding=rounding) / rounding)
    scaled_value2 = builtins.round(float_round(value2, precision_rounding=rounding) / rounding)
    quotient, remainder = divmod(scaled_value1, scaled_value2)
    return quotient, float_round(remainder * rounding, precision_rounding=rounding)


def float_repr(value: float, precision_digits: int) -> str:
    """Returns a string representation of a float with the
       given number of fractional digits. This should not be
       used to perform a rounding operation (this is done via
       :func:`~.float_round`), but only to produce a suitable
       string representation for a float.

       :param value: the value to represent
       :param precision_digits: number of fractional digits to include in the output
       :return: the string representation of the value
    """
    # Can't use str() here because it seems to have an intrinsic
    # rounding to 12 significant digits, which causes a loss of
    # precision. e.g. str(123456789.1234) == str(123456789.123)!!
    if float_is_zero(value, precision_digits=precision_digits):
        value = 0.0
    return "%.*f" % (precision_digits, value)


def float_split_str(value: float, precision_digits: int) -> tuple[str, str]:
    """Splits the given float 'value' in its unitary and decimal parts,
       returning each of them as a string, rounding the value using
       the provided ``precision_digits`` argument.

       The length of the string returned for decimal places will always
       be equal to ``precision_digits``, adding zeros at the end if needed.

       In case ``precision_digits`` is zero, an empty string is returned for
       the decimal places.

       Examples:
           1.432 with precision 2 => ('1', '43')
           1.49  with precision 1 => ('1', '5')
           1.1   with precision 3 => ('1', '100')
           1.12  with precision 0 => ('1', '')

       :param value: value to split.
       :param precision_digits: number of fractional digits to round to.
       :return: returns the tuple(<unitary part>, <decimal part>) of the given value
    """
    value = float_round(value, precision_digits=precision_digits)
    value_repr = float_repr(value, precision_digits)
    return tuple(value_repr.split('.')) if precision_digits else (value_repr, '')


def float_split(value: float, precision_digits: int) -> tuple[int, int]:
    """ same as float_split_str() except that it returns the unitary and decimal
        parts as integers instead of strings. In case ``precision_digits`` is zero,
        0 is always returned as decimal part.
    """
    units, cents = float_split_str(value, precision_digits)
    if not cents:
        return int(units), 0
    return int(units), int(cents)


def json_float_round(
    value: float,
    precision_digits: int,
    rounding_method: RoundingMethod = 'HALF-UP',
) -> float:
    """Not suitable for float calculations! Similar to float_repr except that it
    returns a float suitable for json dump

    This may be necessary to produce "exact" representations of rounded float
    values during serialization, such as what is done by `json.dumps()`.
    Unfortunately `json.dumps` does not allow any form of custom float representation,
    nor any custom types, everything is serialized from the basic JSON types.

    :param precision_digits: number of fractional digits to round to.
    :param rounding_method: the rounding method used: 'HALF-UP', 'UP' or 'DOWN',
           the first one rounding up to the closest number with the rule that
           number>=0.5 is rounded up to 1, the second always rounding up and the
           latest one always rounding down.
    :return: a rounded float value that must not be used for calculations, but
             is ready to be serialized in JSON with minimal chances of
             representation errors.
    """
    rounded_value = float_round(value, precision_digits=precision_digits, rounding_method=rounding_method)
    rounded_repr = float_repr(rounded_value, precision_digits=precision_digits)
    # As of Python 3.1, rounded_repr should be the shortest representation for our
    # rounded float, so we create a new float whose repr is expected
    # to be the same value, or a value that is semantically identical
    # and will be used in the json serialization.
    # e.g. if rounded_repr is '3.1750', the new float repr could be 3.175
    # but not 3.174999999999322452.
    # Cfr. bpo-1580: https://bugs.python.org/issue1580
    return float(rounded_repr)


_INVERTDICT = {
    1e-1: 1e+1, 1e-2: 1e+2, 1e-3: 1e+3, 1e-4: 1e+4, 1e-5: 1e+5,
    1e-6: 1e+6, 1e-7: 1e+7, 1e-8: 1e+8, 1e-9: 1e+9, 1e-10: 1e+10,
    2e-1: 5e+0, 2e-2: 5e+1, 2e-3: 5e+2, 2e-4: 5e+3, 2e-5: 5e+4,
    2e-6: 5e+5, 2e-7: 5e+6, 2e-8: 5e+7, 2e-9: 5e+8, 2e-10: 5e+9,
    5e-1: 2e+0, 5e-2: 2e+1, 5e-3: 2e+2, 5e-4: 2e+3, 5e-5: 2e+4,
    5e-6: 2e+5, 5e-7: 2e+6, 5e-8: 2e+7, 5e-9: 2e+8, 5e-10: 2e+9,
}


def float_invert(value: float) -> float:
    """Inverts a floating point number with increased accuracy.

    :param value: value to invert.
    :return: inverted float.
    """
    result = _INVERTDICT.get(value)
    if result is None:
        coefficient, exponent = f'{value:.15e}'.split('e')
        # invert exponent by changing sign, and coefficient by dividing by its square
        result = float(f'{coefficient}e{-int(exponent)}') / float(coefficient)**2
    return result


if __name__ == "__main__":

    import time
    start = time.time()
    count = 0

    def try_round(amount, expected, precision_digits=3):
        result = float_repr(float_round(amount, precision_digits=precision_digits),
                            precision_digits=precision_digits)
        if result != expected:
            print('###!!! Rounding error: got %s , expected %s' % (result, expected))
            return complex(1, 1)
        return 1

    # Extended float range test, inspired by Cloves Almeida's test on bug #882036.
    fractions = [.0, .015, .01499, .675, .67499, .4555, .4555, .45555]
    expecteds = ['.00', '.02', '.01', '.68', '.67', '.46', '.456', '.4556']
    precisions = [2, 2, 2, 2, 2, 2, 3, 4]
    for magnitude in range(7):
        for frac, exp, prec in zip(fractions, expecteds, precisions):
            for sign in [-1, 1]:
                for x in range(0, 10000, 97):
                    n = x * 10**magnitude
                    f = sign * (n + frac)
                    f_exp = ('-' if f != 0 and sign == -1 else '') + str(n) + exp
                    count += try_round(f, f_exp, precision_digits=prec)

    stop = time.time()
    count, errors = int(count.real), int(count.imag)

    # Micro-bench results:
    # 47130 round calls in 0.422306060791 secs, with Python 2.6.7 on Core i3 x64
    # with decimal:
    # 47130 round calls in 6.612248100021 secs, with Python 2.6.7 on Core i3 x64
    print(count, " round calls, ", errors, "errors, done in ", (stop-start), 'secs')
````

### FILE: `test_odoo_loyalty_derived.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file38:v1"
operation: CREATE
provenance: ADAPTED
source: "https://github.com/odoo/odoo/tree/99edb6dd82b7b560930c00b03b694ba700785370/addons/sale_loyalty; exact file/segment mappings in odoo_loyalty/DERIVATION.json and DERIVATION_CONNECTED.json; Decimal/IPC and immutable history adaptations declared; test expectations adapted, no full Odoo integration-suite claim"
license: "LGPL-3.0-only"
sha256: "4b80a32026d9c34e55ad855bb2ca227290432da927fab5b8eb81f0a184a755cf"
variables: []
secrets_allowed: false
```

````python
import copy
import hashlib
import json
import random
import subprocess
import sys
import unittest
from decimal import Decimal
from pathlib import Path

from odoo_loyalty import engine
from odoo_loyalty.orm_contract import Records, Record, build, coupon_for
from odoo_loyalty.oracle_loader import original_functions
from odoo_loyalty.protocol import calculate, canonical, pairs


def program(kind='gift_card'):
    # ADAPTED fixture configuration: Odoo tests/common.py:104-130, source-only arithmetic.
    return {'id':'gift-1' if kind == 'gift_card' else 'loyalty-1','kind':kind,'currency':'ARS',
            'currency_digits':2,'applies_on':'future' if kind == 'gift_card' else 'both',
            'nominative':kind == 'loyalty','trigger':'auto','trigger_product_ids':['gift-50'],
            'rules':[{'id':'rule-1','mode':'money','points':'1','split':kind == 'gift_card',
                      'minimum_qty':'1','minimum_amount':'0','tax_mode':'incl',
                      'product_ids':['gift-50'] if kind == 'gift_card' else [],'code_required':False}]}


def order(quantity=2, total='100', product='gift-50'):
    return {'order_id':'order-1','organization_id':'org-1','subject_id':'subject-1','state':'draft',
            'public_subject':False,'currency':'ARS','total':total,'enabled_rule_ids':[],
            'lines':[{'id':'line-1','product_id':product,'quantity':str(quantity),'subtotal':total,'tax':'0','total':total,
                      'reward_program_type':'','reward_program_id':'','reward_trigger':'','threshold_excluded':False}]}


def request(p=None, o=None, operation='evaluate', data=None):
    p = p or program()
    return {'schema':'elite.odoo-loyalty-calc.v1','operation':operation,'program':p,
            'program_sha256':hashlib.sha256(canonical(p)).hexdigest(),'order':o or order(),'data':data or {}}


class DerivedTests(unittest.TestCase):
    def test_original_functions_are_source_ast(self):
        manifest=json.loads((Path(__file__).parent/'odoo_loyalty/DERIVATION.json').read_text())
        for mapping in manifest['functions']:
            path=Path(__file__).parent/'odoo_loyalty/upstream'/mapping['source_path']
            raw=path.read_bytes()
            self.assertEqual(hashlib.sha256(raw).hexdigest(),mapping['source_file_sha256'])
            lo,hi=mapping['source_lines']
            self.assertEqual(hashlib.sha256(b''.join(raw.splitlines(keepends=True)[lo-1:hi])).hexdigest(),mapping['source_segment_sha256'])

    def test_upstream_gift_card_quantity_expectations(self):
        # Actual expected card counts from test_buy_gift_card.py:12-36; full Odoo ORM suite is not claimed.
        for quantity in [1,2,1]:
            result=calculate(request(o=order(quantity,str(quantity*50))))['result']
            cards=[Decimal(p) for p in result['points'] if Decimal(p)]
            self.assertEqual(cards,[Decimal(50)]*quantity)

    def test_rule_modes_and_source_oracle(self):
        original=original_functions()
        rng=random.Random(40219)
        for mode in ['order','money','unit']:
            for _ in range(80):
                quantity=rng.randint(1,20)
                total=Decimal(rng.randint(1,100000))/100
                p=program('loyalty')
                p['rules'][0].update(mode=mode,points=str(Decimal(rng.randint(1,100))/10),minimum_amount=str(rng.randint(0,10)))
                o=order(quantity,format(total,'f'),'product-A')
                exact_p,exact_o=build(p,o)
                float_p,float_o=build(p,o,float,original)
                actual=exact_o._program_check_compute_points(Records([exact_p]))[exact_p]
                oracle=float_o._program_check_compute_points(Records([float_p]))[float_p]
                self.assertEqual(actual.keys(),oracle.keys())
                if 'points' in actual:
                    self.assertEqual([x.quantize(Decimal('.01')) if isinstance(x,Decimal) else Decimal(x) for x in actual['points']],
                                     [Decimal(str(x)).quantize(Decimal('.01')) for x in oracle['points']])
                else:
                    self.assertEqual(actual,oracle)

    def test_source_filters_discount_and_payment_lines(self):
        p=program('loyalty')
        o=order(1,'100','product-A')
        for kind,amount in [('promotion','-10'),('gift_card','-20')]:
            o['lines'].append({'id':kind,'product_id':'product-A','quantity':'1','subtotal':amount,'tax':'0','total':amount,
                               'reward_program_type':kind,'reward_program_id':kind,'reward_trigger':'with_code','threshold_excluded':False})
        o['total']='70'
        # Source excludes gift-card tender from money-spend rewards, but includes the discount line.
        self.assertEqual(calculate(request(p,o))['result']['points'],['90.00'])
        o['lines'][1].update(subtotal='-0.25',total='-0.25')
        o['total']='79.75'
        self.assertEqual(calculate(request(p,o))['result']['points'],['99.75'])

    def test_threshold_code_products_and_public_subject(self):
        p=program()
        p['rules'][0].update(code_required=True,minimum_amount='100',minimum_qty='2')
        r=request(p,order())
        self.assertIn('error',calculate(r)['result'])
        r['order']['enabled_rule_ids']=['rule-1']
        self.assertNotIn('error',calculate(r)['result'])
        for amount,quantity in [('99.99',2),('100',1)]:
            r['order']=order(quantity,amount)
            r['order']['enabled_rule_ids']=['rule-1']
            self.assertIn('error',calculate(r)['result'])
        r=request(program(),order(product='other-product'))
        self.assertIn('error',calculate(r)['result'])
        r=request(program('loyalty'))
        r['order']['public_subject']=True
        self.assertIn('error',calculate(r)['result'])

    def test_available_points_future_and_finalized(self):
        original=original_functions()
        for kind in ['gift_card','loyalty']:
            for state,expected in [('draft',Decimal(80) if kind=='gift_card' else Decimal(90)),('sale',Decimal(100))]:
                for number,functions in [(Decimal,None),(float,original)]:
                    p,o=build(program(kind),order(),number,functions)
                    o.state=state
                    card=coupon_for(p,'card-1','100',number)
                    o.coupon_point_ids=Records([Record('earned','e',coupon_id=card,points=number('10'))])
                    o.order_line=Records([Record('used','u',coupon_id=card,points_cost=number('20'))])
                    self.assertEqual(Decimal(str(o._get_real_points_for_coupon(card))),expected)

    def test_gift_card_over_under_and_multiple_bound_amounts(self):
        # Source test_pay_with_gift_card.py:12-69 expects total reduction == before/after balance delta.
        for order_amount,balance,expected in [('115','100','100'),('5.75','100','5.75'),('2300','200','200')]:
            data={'coupon_id':'card-1','balance':balance,'pending_earned':'0','pending_cost':'0','discountable':order_amount,
                  'reward':{'mode':'per_point','discount':'1','required_points':'1','max_amount':'0','clear_wallet':False}}
            result=calculate(request(o=order(1,order_amount),operation='reward',data=data))['result']
            self.assertEqual(Decimal(result['value']),Decimal(expected))
            self.assertEqual(Decimal(result['points_cost']),Decimal(expected))
            self.assertEqual(Decimal(balance)-Decimal(result['points_cost']),Decimal(balance)-Decimal(expected))

    def test_changes_preserves_card_identity_and_reversal_amount(self):
        data={'issued':[{'coupon_id':'card-A','points':'100'},{'coupon_id':'card-A','points':'5'},{'coupon_id':'card-B','points':'50'}],
              'used':[{'coupon_id':'card-A','points':'30'},{'coupon_id':'card-B','points':'10'}]}
        result=calculate(request(operation='changes',data=data))['result']
        self.assertEqual(result,{'card-A':'75','card-B':'40'})
        # The durable layer must record the negation of this same committed change for source cancel intent.
        self.assertEqual({k:-Decimal(v) for k,v in result.items()},{'card-A':Decimal(-75),'card-B':Decimal(-40)})

    def test_exact_decimal_large_value_does_not_inherit_float_loss(self):
        p=program('loyalty')
        o=order(1,'99999999999999999.99','product-A')
        self.assertEqual(calculate(request(p,o))['result']['points'],['99999999999999999.99'])
        self.assertNotEqual(Decimal(str(float('99999999999999999.99'))),Decimal('99999999999999999.99'))

    def test_invalid_hash_amount_binding_and_no_partial_contract(self):
        mutations=[lambda r:r.update(program_sha256='0'*64),
                   lambda r:r['program']['rules'][0].update(points='0'),
                   lambda r:r['program'].update(kind='ewallet'),
                   lambda r:r['order'].update(currency='USD'),
                   lambda r:r['order'].update(total='101'),
                   lambda r:r['order']['lines'][0].update(quantity='1.5'),
                   lambda r:r['order']['lines'][0].update(tax='1'),
                   lambda r:r['order'].update(public_subject=None),
                   lambda r:r['program'].update(currency_digits=True),
                   lambda r:r['order']['lines'][0].update(subtotal=100),
                   lambda r:r['order'].update(enabled_rule_ids=['unknown']),
                   lambda r:r['program'].update(tax_rate='21'),
                   lambda r:r['order'].update(lines=r['order']['lines']*2)]
        for mutate in mutations:
            r=request();mutate(r)
            with self.subTest(r=r),self.assertRaises((ValueError,TypeError)):
                calculate(r)

    def test_duplicate_json_and_exact_ipc_output(self):
        with self.assertRaises(ValueError):
            json.loads('{"schema":1,"Schema":2}',object_pairs_hook=pairs)
        r=request()
        process=subprocess.run([sys.executable,'-B','-m','odoo_loyalty.protocol'],input=canonical(r),capture_output=True,cwd=Path(__file__).parent,timeout=5)
        self.assertEqual(process.returncode,0,process.stderr)
        self.assertEqual(json.loads(process.stdout),calculate(r))
        bad=subprocess.run([sys.executable,'-B','-m','odoo_loyalty.protocol'],input=b'{"sensitive":"never-echo-this"}',capture_output=True,cwd=Path(__file__).parent,timeout=5)
        self.assertEqual(bad.returncode,2)
        self.assertNotIn(b'never-echo-this',bad.stdout+bad.stderr)


if __name__ == '__main__':
    unittest.main()
````

### FILE: `test_stored_value_source_loader.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file39:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "5461b32f352ce686168d8f4fed5feef0e1c565ca097b772c0eaeb01429716deb"
variables: []
secrets_allowed: false
```

````python
# AUTHORED source-closure regressions. Extra files contain only test exceptions;
# no provider calls, credentials, user data or harmful effects are involved.
import hashlib
import json
from pathlib import Path
import py_compile
import shutil
import subprocess
import sys
import tempfile
import unittest
from decimal import Decimal
from importlib.util import cache_from_source
from test_odoo_loyalty_derived import request

class SourceLoaderTests(unittest.TestCase):
    def run_copy(self, mutation=None):
        with tempfile.TemporaryDirectory(prefix='stored-value-source-fixture-') as temporary:
            root=Path(temporary)/'odoo_loyalty'
            shutil.copytree(Path(__file__).parent/'odoo_loyalty',root,ignore=shutil.ignore_patterns('__pycache__','*.pyc'))
            if mutation: mutation(root)
            lock=hashlib.sha256((root/'engine-lock.json').read_bytes()).hexdigest()
            return subprocess.run([sys.executable,'-I','-S','-B',str(root/'run.py'),lock],input=json.dumps(request()).encode(),capture_output=True,timeout=10)
    def test_ignores_unlocked_initializer(self):
        result=self.run_copy(lambda root:(root/'__init__.py').write_text("raise RuntimeError('UNLOCKED_INITIALIZER_FIXTURE')\n",encoding='utf-8'))
        self.assertEqual(result.returncode,0,result.stderr.decode())
        self.assertEqual([Decimal(p) for p in json.loads(result.stdout)['result']['points']],[Decimal(50),Decimal(50)])
    def test_ignores_unlocked_bytecode(self):
        def add(root):
            source=root/'fixture_cache_source.py'
            source.write_text("raise RuntimeError('UNLOCKED_BYTECODE_FIXTURE')\n",encoding='utf-8')
            cache=Path(cache_from_source(str(root/'protocol.py')));cache.parent.mkdir(exist_ok=True)
            py_compile.compile(str(source),cfile=str(cache),doraise=True,invalidation_mode=py_compile.PycInvalidationMode.UNCHECKED_HASH)
        result=self.run_copy(add)
        self.assertEqual(result.returncode,0,result.stderr.decode())
        self.assertEqual([Decimal(p) for p in json.loads(result.stdout)['result']['points']],[Decimal(50),Decimal(50)])
    def test_rejects_changed_locked_source(self):
        def add(root):
            p=root/'engine.py';p.write_bytes(p.read_bytes()+b'\n# changed fixture\n')
        result=self.run_copy(add)
        self.assertEqual(result.returncode,2)
        self.assertEqual(result.stdout,b'')


class SourceReplacementTests(unittest.TestCase):
    def test_exact_and_modified_source_rebuilds_preserve_licenses(self):
        from tools.rebuild_stored_value_source import rebuild
        original=Path(__file__).parent/'odoo_loyalty'
        baseline=hashlib.sha256((original/'engine-lock.json').read_bytes()).hexdigest()
        with tempfile.TemporaryDirectory(prefix='stored-value-replacement-fixture-') as temporary:
            root=Path(temporary)
            exact=root/'exact';receipt=rebuild(original,baseline,exact)
            self.assertEqual(receipt['classification'],'EXACT_REBUILD')
            self.assertEqual(receipt['replacement_lock_sha256'],baseline)
            (exact/'engine.py').write_bytes((exact/'engine.py').read_bytes()+b'\n# User modification fixture, no original attribution claimed.\n')
            modified=root/'modified';receipt=rebuild(exact,baseline,modified)
            self.assertEqual(receipt['classification'],'USER_MODIFIED_SOURCE_NOT_ADMITTED')
            self.assertEqual([r['path'] for r in receipt['changes']],['engine.py'])
            for name in ['LICENSE','COPYRIGHT']:
                self.assertEqual((original/name).read_bytes(),(modified/name).read_bytes())
            run=lambda digest:subprocess.run([sys.executable,'-I','-S','-B',str(modified/'run.py'),digest],input=json.dumps(request()).encode(),capture_output=True,timeout=10)
            self.assertEqual(run(baseline).returncode,2)
            self.assertEqual(run(receipt['replacement_lock_sha256']).returncode,0)
            with self.assertRaises(ValueError):rebuild(exact,baseline,modified)

if __name__=='__main__':unittest.main()
````

### FILE: `tools/rebuild_stored_value_source.py`

```yaml
block_id: "PYTHON-ODOO-STORED-VALUE-CALCULATOR:file40:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration/configuration/authorization/persistence/transport/UI glue; no upstream company authorship"
license: "LicenseRef-Workspace-Owner"
sha256: "efda04b47fa0dddf3c66516664dd207514171e6ae7b92840b55464b70bb5c02d"
variables: []
secrets_allowed: false
```

````python
"""AUTHORED rebuild/receipt glue for a user's separately modified LGPL module.
No source changes are made here and the result is not automatically admitted.
"""
import argparse
import hashlib
import json
from pathlib import Path
import re


def sha(raw):return hashlib.sha256(raw).hexdigest()


def rebuild(source,expected,output):
    source=source.resolve();output=output.resolve()
    if output.exists() or output.is_relative_to(source) or not re.fullmatch('[0-9a-f]{64}',expected):
        raise ValueError('an absent destination outside the source module is required')
    raw=(source/'engine-lock.json').read_bytes()
    if len(raw)>65536 or sha(raw)!=expected:raise ValueError('baseline source manifest mismatch')
    lock=json.loads(raw)
    if lock['schema']!='elite.odoo-derived-source-lock.v1' or not 1<=len(lock['files'])<=128:raise ValueError('source lock contract')
    snapshots={};changes=[];total=0
    for row in lock['files']:
        rel=row['path'];part=Path(rel)
        if not isinstance(rel,str) or part.is_absolute() or '..' in part.parts or rel in snapshots:raise ValueError('source path contract')
        p=(source/part).resolve()
        if not p.is_relative_to(source):raise ValueError('source path escapes module')
        data=p.read_bytes();total+=len(data)
        if total>8*1024*1024:raise ValueError('source module exceeds rebuild budget')
        after=sha(data)
        if after!=row['sha256'] or len(data)!=row['bytes']:changes.append({'path':rel,'before_sha256':row['sha256'],'after_sha256':after})
        snapshots[rel]=data;row.update(bytes=len(data),sha256=after)
    if not {'LICENSE','COPYRIGHT','run.py','engine.py','orm_contract.py','protocol.py'}<=snapshots.keys():raise ValueError('incomplete licensed module')
    # Preserve license texts/attribution when making a derivative snapshot.
    if any(r['path'] in ('LICENSE','COPYRIGHT','upstream/LICENSE','upstream/COPYRIGHT') for r in changes):raise ValueError('preserve original license and attribution texts')
    output.mkdir(parents=True,exist_ok=False)
    for rel,data in snapshots.items():
        p=output/rel;p.parent.mkdir(parents=True,exist_ok=True);p.write_bytes(data)
    encoded=(json.dumps(lock,indent=2)+'\n').encode('utf-8');(output/'engine-lock.json').write_bytes(encoded)
    receipt={'schema':'elite.user-source-replacement.v1','classification':'USER_MODIFIED_SOURCE_NOT_ADMITTED' if changes else 'EXACT_REBUILD','baseline_lock_sha256':expected,'replacement_lock_sha256':sha(encoded),'changes':changes,'next_step':'Retain LGPL sources/notices. Review and test the change, then explicitly bind its script/manifest hashes in the deployment. No upstream authorship or prior gate status is inherited.'}
    (output/'source-replacement-receipt.json').write_text(json.dumps(receipt,indent=2)+'\n',encoding='utf-8',newline='\n')
    return receipt


if __name__=='__main__':
    p=argparse.ArgumentParser(description=__doc__)
    p.add_argument('--module-dir',type=Path,required=True);p.add_argument('--expected-lock-sha256',required=True);p.add_argument('--destination',type=Path,required=True)
    a=p.parse_args()
    try:print(json.dumps(rebuild(a.module_dir,a.expected_lock_sha256,a.destination),indent=2))
    except (OSError,ValueError,KeyError,TypeError) as e:p.exit(2,str(e)+'\n')
````

## 6. Configuration surface

See deploy/stored-value/profile.reference.json and docs/stored-value-runtime.md in
the connected pair. Enable only a hash-bound tenant/org/profile and exact Python
source lock. Empty credentials do not authorize a live payment. Source overrides
write a new tree/receipt and are not automatically admitted.

## 7. Dependency bill

Official Odoo19.0 commit99edb6dd82b7b560930c00b03b694ba700785370, selected source only,
LGPL-3.0-only. No Odoo application installation or added pip dependency. Existing
locked CPython3.14.4/Go1.26.8/PG18.6/Next16.3.4 runtime owners are retained. Detailed
source/Gitblob/SHA256 and runtime manifests are under docs/provenance/ODOO_*.
Complete modified module source and exact LICENSE/COPYRIGHT ship with the Python
pack; IPC alone does not determine the legal boundary.

## 8. Apply order

Compose into an absent/empty destination using the canonical full profile and
MARKDOWN-COMPOSITOR0.3.0. Apply selected migrations in numeric order; shared0066,
stored-value0067 and handover0068 retain distinct owners. Never replace existing
history. Down migrations refuse populated immutable history; explicit dependency
drops avoid CASCADE. Enable the optional host module only after exact profile/source
validation. The updater preserves source LICENSE without adding a final LF.

## 9. Verification

Source oracle/expectation comparisons, loader/source-replacement regressions and
Go fuzz receipts are recorded in STORED_VALUE_SOURCE_CLOSURE_V402.md. Real PG/HTTP
partial gift and loyalty SDK payments, fully funded delivery, separate reviewer,
commit expiry rollback, concurrency, body/permission negatives, mobile browser and
lost-response recovery are in STORED_VALUE_TENDER_V402.md and
STORED_VALUE_OPERATOR_FLOW_V402.md. Changed-source checks must rerun their affected
gates. These receipts do not claim the full Odoo suite, SAST/DAST, target load or
production acceptance. Full composition SCA remains the separate T2803 owner.

## 10. Reconstruction evidence

Five profiles reconstruct exactly from source packs: full franchise84/1112,
backend40/567, web8/163, serverless57/793 and HTTP metrics85/1120. The two narrower
Go compositions compile changed owners and test packages with zero tests executed.
Runtime receipts from the identical candidate source are retained without replaying
unchanged payment/browser/domain suites. Canonical admission evidence is
reconstruction_evidence/STORED_VALUE_CONNECTED_RELEASE_V402.md and its JSON.

CONDITIONED selection is limited to the declared assisted reference program,
exact source/runtime/profile hashes, corresponding LGPL source/notices and the
existing scoped approval/payment/handover owners. The module's local functionality
is proven. Whole-composition SCA/T2803, operations/T2809 and target business/live
acceptance remain their explicit owners; these are not asserted by this pack.
The metadata-only change to REBUILD_VERIFIED preserves every materialized byte;
a final reference reconstruction binds the admitted pack bytes before publication.
