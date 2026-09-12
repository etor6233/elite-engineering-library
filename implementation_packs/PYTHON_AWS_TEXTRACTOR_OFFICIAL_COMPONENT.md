# Python AWS Textractor Official Component

## 1. Metadata

```yaml
pack_id: "PYTHON-AWS-TEXTRACTOR-OFFICIAL-COMPONENT"
pack_version: "0.1.0"
status:
  authority: ELITE_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa el runtime oficial AWS Amazon Textractor 1.10.0 mediante wheel hash-locked y sus pruebas/fixtures oficiales de Queries, con locks separados, SCA fechada y límites de portabilidad preservados."
stacks: ["Python 3.14", "amazon-textract-textractor 1.10.0", "Amazon Textract"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "SECURE-LOCAL-FILE-INGESTION-GATE 0.1.x", "STRICT-DOCUMENT-FIELD-EVALUATION-GATE 0.1.x", "GO-AWS-TEXTRACT-DOCUMENT-RUNTIME 0.2.x"]
incompatible_with: ["CALL_TEXTRACT en gates offline", "persistencia automática sin corpus", "grafo sin hashes", "credenciales o región embebidas", "claim de exactitud perfecta"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0"
upstream_sources: ["https://github.com/aws-samples/amazon-textract-textractor/tree/8ea5f9ae65dbcbb2b83b409e0007c395e59b0e92", "https://pypi.org/project/amazon-textract-textractor/1.10.0/"]
verified_at: "2026-08-28"
```

## 2. Applicability

Use cuando el proyecto elija Amazon Textract y necesite el paquete oficial AWS que convierte respuestas en objetos navegables y expone texto, tablas, forms, queries, IDs, expenses, Markdown/HTML y CLI. Es una alternativa/complemento Python al runtime Go existente, no una reescritura Elite del parser.

Rechazar si falta una cuenta/región/IAM/costo aprobados, si el input no pasó seguridad, si se pretende almacenar campos sin evaluación por clase o si se exige que un fixture oficial demuestre exactitud sobre facturas, proformas, packing lists, BOL, aduana u otros documentos del negocio.

## 3. Architecture contract

El wheel oficial 1.10.0 se adquiere por nombre, bytes y SHA-256 desde PyPI; no se redistribuye dentro del Markdown. El runtime productivo tiene 15 pins con hashes para wheels CPython 3.14 Windows/Linux; tooling se separa en once pins adicionales. LICENSE/NOTICE y cinco archivos oficiales de prueba/fixtures se reconstruyen desde el commit fijado, con adaptaciones declaradas sólo cuando Markdown añade el LF final.

La prueba oficial empaquetada procesa dos respuestas API guardadas sin llamar AWS: una confirma dos respuestas Query y otra confirma queries sin resultado. `CALL_TEXTRACT` debe estar ausente. La suite determinista upstream completa se probó en Linux; dos divergencias de newline/word ordering en Windows y el fuzz aleatorio sin marker permanecen condiciones upstream, no parches locales.

## 4. Exact file manifest

```text
CREATE aws_textractor_official/README.md
CREATE aws_textractor_official/source-lock.json
CREATE aws_textractor_official/requirements-runtime.lock
CREATE aws_textractor_official/requirements-test.lock
CREATE aws_textractor_official/test_contracts.py
CREATE aws_textractor_official/LICENSE
CREATE aws_textractor_official/NOTICE
CREATE aws_textractor_official/upstream/tests/__init__.py
CREATE aws_textractor_official/upstream/tests/utils.py
CREATE aws_textractor_official/upstream/tests/test_queries.py
CREATE aws_textractor_official/upstream/tests/fixtures/saved_api_responses/test_queries_as_strings.json
CREATE aws_textractor_official/upstream/tests/fixtures/saved_api_responses/test_bad_queries_as_strings.json
```

## 5. Materialization blocks

### FILE: `aws_textractor_official/README.md`
```yaml
block_id: "PYTHON-AWS-TEXTRACTOR-OFFICIAL:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "Elite packaging and conditions around official AWS source"
license: "LicenseRef-Workspace-Owner"
sha256: "c13c626e5000140bcd398f8c9b473cd44c1b9854f7b18925aa170b49bf42a8a0"
variables: []
secrets_allowed: false
```
````markdown
# AWS Amazon Textractor official component

This component uses the official AWS `amazon-textract-textractor` 1.10.0 wheel and byte-preserved source tests from `aws-samples/amazon-textract-textractor` tag `v1.10.0`, commit `8ea5f9ae65dbcbb2b83b409e0007c395e59b0e92`.

Install the runtime from an acquired wheelhouse with:

```text
python -m pip install --no-index --find-links <wheelhouse> --require-hashes -r requirements-runtime.lock
```

The official package then exposes both the `textractor` CLI and the `Textractor` Python API for text, tables, forms, queries, identity documents and expenses. AWS credentials, region, IAM, service cost and live calls remain external project inputs. Never set `CALL_TEXTRACT` during offline verification.

Install the separate test graph and execute the packaged official query regression with:

```text
python -m pip install --no-index --find-links <wheelhouse> --require-hashes -r requirements-test.lock
python -m pytest -q upstream/tests/test_queries.py
```

Expected offline result: `2 passed, 1 skipped`. The skip is the official live-only asynchronous request test. The complete deterministic upstream core suite, excluding only the file whose own docstring calls it disabled/flaky fuzzing, produced `69 passed, 16 skipped` on Linux Python 3.14. On Windows it produced `67 passed, 16 skipped` plus two newline/word-order portability failures; those exact failures remain conditions 156 and 157 and are not patched here. Condition 158 preserves the unmarked random fuzz test.

OSV-Scanner 2.5.1 found zero known vulnerabilities in the dated 15-package runtime graph and zero in the 26-package test graph on 2026-08-28. This is a dated observation, not a future guarantee; every project must rescan its final graph.

This is an official parser/runtime component, not proof of perfect extraction. It does not provide ground truth, schema mapping, field confidence policy, automatic database writes, secure intake, tenancy, ordering, idempotency, human review, reconciliation or cloud authorization. Compose it only after the routing and security gates, then prove each document class against the project's representative corpus before storing business facts.
````

### FILE: `aws_textractor_official/source-lock.json`
```yaml
block_id: "PYTHON-AWS-TEXTRACTOR-OFFICIAL:source-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "GitHub/PyPI/archive/wheel/test/SCA evidence lock"
license: "LicenseRef-Workspace-Owner"
sha256: "7b638fed83b675ff6de9d4758b01e7ecdd48833a108415c57e459554966165ff"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-aws-textractor-official-source-lock/v1",
  "repository": "aws-samples/amazon-textract-textractor",
  "release": "v1.10.0",
  "commit": "8ea5f9ae65dbcbb2b83b409e0007c395e59b0e92",
  "archive": {
    "bytes": 79080629,
    "sha256": "4fa38999e355d977d2b46816b76c9d32ab6c5d8607987e086045470d101c7219"
  },
  "wheel": {
    "name": "amazon_textract_textractor-1.10.0-py3-none-any.whl",
    "bytes": 311287,
    "sha256": "8524224f07a776ca2959e0d80d5e031a30bf53938c74e8352848f3d0b6b3762c"
  },
  "license": "Apache-2.0",
  "files": [
    {"path":"LICENSE","provenance":"VERBATIM","source_bytes":10142,"source_sha256":"09e8a9bcec8067104652c168685ab0931e7868f9c8284b66f5ae6edae5f1130b","packaged_sha256":"09e8a9bcec8067104652c168685ab0931e7868f9c8284b66f5ae6edae5f1130b"},
    {"path":"NOTICE","provenance":"VERBATIM","source_bytes":100,"source_sha256":"494d0884edd2d9dee12e96c5be08873f24bbb66d1fca68ac5fa7b542a8361465","packaged_sha256":"494d0884edd2d9dee12e96c5be08873f24bbb66d1fca68ac5fa7b542a8361465"},
    {"path":"upstream/tests/__init__.py","provenance":"ADAPTED_FINAL_NEWLINE_ONLY","source_bytes":0,"source_sha256":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","packaged_sha256":"01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b"},
    {"path":"upstream/tests/utils.py","provenance":"ADAPTED_FINAL_NEWLINE_ONLY","source_bytes":699,"source_sha256":"9cccbff833592ee146de908d19f9b8cff39f962cffe38f046a989380cbd3e9d6","packaged_sha256":"8dacd36f27de6d7d367ee6f3b0d023abf5c9245f92abeb77fd7c2d8cdd0cf4cd"},
    {"path":"upstream/tests/test_queries.py","provenance":"VERBATIM","source_bytes":3543,"source_sha256":"ab7366659a68e555e43a868a67a1b6f2e267c15eeb3eff11065498daf04b12a2","packaged_sha256":"ab7366659a68e555e43a868a67a1b6f2e267c15eeb3eff11065498daf04b12a2"},
    {"path":"upstream/tests/fixtures/saved_api_responses/test_queries_as_strings.json","provenance":"ADAPTED_FINAL_NEWLINE_ONLY","source_bytes":44405,"source_sha256":"f3d443b5e9de5da8de08178dc139e9daa416badef83ab491c284e6ea2dcb2461","packaged_sha256":"f3786e60a207e740a4e3cb35696d5292506d31d258f072e78e07e3a6bf0215ca"},
    {"path":"upstream/tests/fixtures/saved_api_responses/test_bad_queries_as_strings.json","provenance":"ADAPTED_FINAL_NEWLINE_ONLY","source_bytes":43146,"source_sha256":"8a90dfcdacdecb263b086a9725edd5cab32a125ee832aea2594b9d37cf46227c","packaged_sha256":"12e2327459483c96a05dfb9babe7968f3d01c5b199caeb94296c43a2561927a7"}
  ],
  "verification": {
    "offline_query_tests": {"passed":2,"skipped":1},
    "linux_deterministic_core": {"passed":69,"skipped":16,"excluded":"tests/test_parse_no_fail.py declared disabled/flaky upstream"},
    "windows_deterministic_core": {"passed":67,"skipped":16,"failed":2,"conditions":[156,157]},
    "runtime_osv": {"packages":15,"findings":0,"report_sha256":"3083672b5fcd9f1d529fa9256f75adfc6ba1a9802d8e4049899d73d2860b7787"},
    "test_osv": {"packages":26,"findings":0,"report_sha256":"1b63bfab44bbad8970d06b52267aa65c4bf36ce1e2e6ca7af1d95e7dbd58475d"}
  },
  "verified_at": "2026-08-28"
}
````

### FILE: `aws_textractor_official/requirements-runtime.lock`
```yaml
block_id: "PYTHON-AWS-TEXTRACTOR-OFFICIAL:runtime-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "Resolved exact runtime graph for official wheel"
license: "LicenseRef-Workspace-Owner"
sha256: "cb7719c63ccebfaa283b56b8c22b3a75c2082615a128915d9bf4385563087cc0"
variables: []
secrets_allowed: false
```
````text
amazon-textract-caller==0.2.4 --hash=sha256:ec7dc3517f1cc9b37b41a74b2b5ea040d67be91e8559a8150f44af75bf7f5590
amazon-textract-response-parser==1.0.3 --hash=sha256:834ffcec01085b82565b3f55625572f3425885918f09380cbe9e7fb02a7c8de7
amazon-textract-textractor==1.10.0 --hash=sha256:8524224f07a776ca2959e0d80d5e031a30bf53938c74e8352848f3d0b6b3762c
boto3==1.43.82 --hash=sha256:65319e8bba6e30f74a6e2727a5688725222da2ab71c6069bc484ce5bfd101c73
botocore==1.43.82 --hash=sha256:97b3e89061decc91e7745d726dae595cbe2053894611c49ea91d6fdcb2ecc36b
jmespath==1.1.0 --hash=sha256:a5663118de4908c91729bea0acadca56526eb2698e83de10cd116ae0f4e97c64
marshmallow==3.26.2 --hash=sha256:013fa8a3c4c276c24d26d84ce934dc964e2aa794345a0f8c7e5a7191482c8a73
pillow==12.3.0 --hash=sha256:fdafc9cce40277e0f7a0feabce0ee50dd2fa1800f3b38015e51296b5e814048d --hash=sha256:251bf95b67017e27b13d82f5b326234ca62d70f9cf4c2b9032de2358a3b12c7b
python-dateutil==2.9.0.post0 --hash=sha256:a8b2bc7bffae282281c8140a97d3aa9c14da0b136dfe83f850eea9a5f7470427
RapidFuzz==3.14.5 --hash=sha256:5dfa89d78f22cd773054caff44827b846161a29f2dcf7e78b8f90d086621e502 --hash=sha256:4900143d82071bdda533b00300c40b14b963ff826b3642cc463b6dd0f036585e
s3transfer==0.19.2 --hash=sha256:d8168eccca828cbb2cd573675333f3bddd254313a9c42494b84c76b539e8ba25
six==1.17.0 --hash=sha256:4721f391ed90541fddacab5acf947aa0d3dc7d27b2e1e8eda2be8970586c3274
tabulate==0.9.0 --hash=sha256:024ca478df22e9340661486f85298cff5f6dcdba14f3813e8830015b9ed1948f
urllib3==2.7.0 --hash=sha256:9fb4c81ebbb1ce9531cce37674bbc6f1360472bc18ca9a553ede278ef7276897
xlsxwriter==3.2.9 --hash=sha256:9a5db42bc5dff014806c58a20b9eae7322a134abb6fce3c92c181bfb275ec5b3
````

### FILE: `aws_textractor_official/requirements-test.lock`
```yaml
block_id: "PYTHON-AWS-TEXTRACTOR-OFFICIAL:test-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "Resolved exact test-only graph over runtime lock"
license: "LicenseRef-Workspace-Owner"
sha256: "02a03aee3f955f5dc99b3f0331eb418fe06185a621107a498e87a9861e0392b9"
variables: []
secrets_allowed: false
```
````text
-r requirements-runtime.lock
colorama==0.4.6 --hash=sha256:4f1d9991f5acc0ca119f9d443620b77f9d6b33703e51011c16baf57afb285fc6
iniconfig==2.3.0 --hash=sha256:f631c04d2c48c52b84d0d0549c99ff3859c98df65b3101406327ecc7d53fbf12
lxml==6.1.2 --hash=sha256:87e9673cd8a3445024fe38e7f91b55fa3428437eec9b7a7ff7d81979520c0d2d --hash=sha256:47e367dfe341521426692819803e260d0673899c0ff611f14af978d725e2c999
numpy==2.5.2 --hash=sha256:7587f53dfbd5edc0f7b87c6217b4c6d2d1f2ef9c3da70bc1315e7db5f8d7ec9d --hash=sha256:318b9a4c845dbea06708a29c84ee429cc3065048db34cdb799047643492050ee
packaging==26.3 --hash=sha256:d7193f7c8e4e93f444fde0262bf90af30e16fa0ad0ad44cb553c87339b23cd1c
pandas==3.0.5 --hash=sha256:cd8f7c6dc98527058ee6264219343f5392240a6f1bfa654fc5d79023020d0c92 --hash=sha256:9e94c2c5ca43bd3ca32bf64d32308887b65e5f9bfd8023ea52755107a999f93b
pdf2image==1.16.3 --hash=sha256:b6154164af3677211c22cbb38b2bd778b43aca02758e962fe1e231f6d3b0e380
pluggy==1.6.0 --hash=sha256:e920276dd6813095e9377c0bc5566d94c932c33b27a3e3945d8389c374dd4746
Pygments==2.21.0 --hash=sha256:2363c69b61c4a97c838da3b130dcd6468f4848992b21a82f2a63ec34377137d9
pytest==9.0.3 --hash=sha256:2c5efc453d45394fdd706ade797c0a81091eccd1d6e4bccfcd476e2b8e0ab5d9
tzdata==2026.3 --hash=sha256:dc096730c87af6cab1b171c9d532be840741ff5d459015e7f6947bd7d7e54931
````

### FILE: `aws_textractor_official/test_contracts.py`
```yaml
block_id: "PYTHON-AWS-TEXTRACTOR-OFFICIAL:contracts:v1"
operation: CREATE
provenance: AUTHORED
source: "Offline packaging contracts; does not replace upstream tests"
license: "LicenseRef-Workspace-Owner"
sha256: "dc5415f6ab1fc8e4993847290d6ae2e99d1197fe534fffeb92c4bdfd46683e97"
variables: []
secrets_allowed: false
```
````python
import hashlib
import json
import re
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parent


def digest(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


class AwsTextractorOfficialContracts(unittest.TestCase):
    def setUp(self) -> None:
        self.lock = json.loads((ROOT / "source-lock.json").read_text(encoding="utf-8"))

    def test_source_identity_and_packaged_hashes(self) -> None:
        self.assertEqual(self.lock["repository"], "aws-samples/amazon-textract-textractor")
        self.assertEqual(self.lock["release"], "v1.10.0")
        self.assertEqual(self.lock["commit"], "8ea5f9ae65dbcbb2b83b409e0007c395e59b0e92")
        self.assertEqual(self.lock["license"], "Apache-2.0")
        for item in self.lock["files"]:
            self.assertEqual(digest(ROOT / item["path"]), item["packaged_sha256"])

    def test_runtime_and_test_locks_are_exact_and_separate(self) -> None:
        runtime = (ROOT / "requirements-runtime.lock").read_text(encoding="utf-8")
        test = (ROOT / "requirements-test.lock").read_text(encoding="utf-8")
        self.assertEqual(len([line for line in runtime.splitlines() if "==" in line]), 15)
        self.assertEqual(len([line for line in test.splitlines() if "==" in line]), 11)
        self.assertIn("amazon-textract-textractor==1.10.0", runtime)
        self.assertIn("pytest==9.0.3", test)
        self.assertNotIn("pytest", runtime)
        for line in runtime.splitlines() + test.splitlines():
            if "==" in line:
                self.assertRegex(line, r"--hash=sha256:[0-9a-f]{64}")

    def test_official_query_fixtures_and_tests_are_present(self) -> None:
        tests = (ROOT / "upstream/tests/test_queries.py").read_text(encoding="utf-8")
        self.assertIn("What is the name of the package?", tests)
        self.assertIn('self.assertEqual(document.queries[0].result.answer, "Textractor")', tests)
        for name in ("test_queries_as_strings.json", "test_bad_queries_as_strings.json"):
            payload = json.loads((ROOT / "upstream/tests/fixtures/saved_api_responses" / name).read_text(encoding="utf-8"))
            self.assertIsInstance(payload.get("Blocks"), list)
            self.assertGreater(len(payload["Blocks"]), 0)

    def test_conditions_and_sca_are_fail_closed(self) -> None:
        readme = (ROOT / "README.md").read_text(encoding="utf-8")
        self.assertIn("69 passed, 16 skipped", readme)
        self.assertIn("67 passed, 16 skipped", readme)
        self.assertIn("not proof of perfect extraction", readme)
        self.assertEqual(self.lock["verification"]["runtime_osv"]["findings"], 0)
        self.assertEqual(self.lock["verification"]["test_osv"]["findings"], 0)
        self.assertEqual(self.lock["verification"]["windows_deterministic_core"]["conditions"], [156, 157])


if __name__ == "__main__":
    unittest.main()
````

### FILE: `aws_textractor_official/LICENSE`
```yaml
block_id: "PYTHON-AWS-TEXTRACTOR-OFFICIAL:license:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/aws-samples/amazon-textract-textractor/8ea5f9ae65dbcbb2b83b409e0007c395e59b0e92/LICENSE"
license: "Apache-2.0"
sha256: "09e8a9bcec8067104652c168685ab0931e7868f9c8284b66f5ae6edae5f1130b"
variables: []
secrets_allowed: false
```
````text

                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.
````

### FILE: `aws_textractor_official/NOTICE`
```yaml
block_id: "PYTHON-AWS-TEXTRACTOR-OFFICIAL:notice:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/aws-samples/amazon-textract-textractor/8ea5f9ae65dbcbb2b83b409e0007c395e59b0e92/NOTICE"
license: "Apache-2.0"
sha256: "494d0884edd2d9dee12e96c5be08873f24bbb66d1fca68ac5fa7b542a8361465"
variables: []
secrets_allowed: false
```
````text
Amazon Textract Textractor
Copyright 2022 Amazon.com, Inc. or its affiliates. All Rights Reserved. 
````

### FILE: `aws_textractor_official/upstream/tests/__init__.py`
```yaml
block_id: "PYTHON-AWS-TEXTRACTOR-OFFICIAL:tests-init:v1"
operation: CREATE
provenance: ADAPTED
source: "https://raw.githubusercontent.com/aws-samples/amazon-textract-textractor/8ea5f9ae65dbcbb2b83b409e0007c395e59b0e92/tests/__init__.py"
license: "Apache-2.0"
sha256: "01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b"
variables: []
secrets_allowed: false
```
````python

````

### FILE: `aws_textractor_official/upstream/tests/utils.py`
```yaml
block_id: "PYTHON-AWS-TEXTRACTOR-OFFICIAL:tests-utils:v1"
operation: CREATE
provenance: ADAPTED
source: "https://raw.githubusercontent.com/aws-samples/amazon-textract-textractor/8ea5f9ae65dbcbb2b83b409e0007c395e59b0e92/tests/utils.py"
license: "Apache-2.0"
sha256: "8dacd36f27de6d7d367ee6f3b0d023abf5c9245f92abeb77fd7c2d8cdd0cf4cd"
variables: []
secrets_allowed: false
```
````python
import inspect
import json
import os

def get_fixture_path():
    """Uses reflection to get correct saved response file

    :return: Path to the saved response file for the calling function
    :rtype: str
    """
    return os.path.join(
        os.path.abspath(os.path.dirname(__file__)),
        f"fixtures/saved_api_responses/{inspect.currentframe().f_back.f_code.co_name}.json"
    )

def save_document_to_fixture_path(document):
    with open(
        os.path.join(
            os.path.abspath(os.path.dirname(__file__)),
            f"fixtures/saved_api_responses/{inspect.currentframe().f_back.f_code.co_name}.json"
        ),
        "w"
    ) as f:
        json.dump(document.response, f)
````

### FILE: `aws_textractor_official/upstream/tests/test_queries.py`
```yaml
block_id: "PYTHON-AWS-TEXTRACTOR-OFFICIAL:test-queries:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/aws-samples/amazon-textract-textractor/8ea5f9ae65dbcbb2b83b409e0007c395e59b0e92/tests/test_queries.py"
license: "Apache-2.0"
sha256: "ab7366659a68e555e43a868a67a1b6f2e267c15eeb3eff11065498daf04b12a2"
variables: []
secrets_allowed: false
```
````python
import os
import unittest
from tests.utils import get_fixture_path
from textractor import Textractor
from textractor.entities.document import Document
from textractor.exceptions import InputError, InvalidProfileNameError
from textractor.data.constants import TextractFeatures

from .utils import save_document_to_fixture_path

class QueriesTests(unittest.TestCase):
    def test_queries_as_strings(self):
        profile_name = "default"
        current_directory = os.path.abspath(os.path.dirname(__file__))

        if profile_name is None:
            raise InvalidProfileNameError(
                "Textractor could not be initialized. Populate profile_name with a valid input in tests/test_table.py."
            )

        if os.environ.get("CALL_TEXTRACT"):
            extractor = Textractor(profile_name=profile_name, kms_key_id="")
            document = extractor.analyze_document(
                file_source=os.path.join(current_directory, "fixtures/single-page-1.png"),
                features=[TextractFeatures.QUERIES],
                queries=[
                    "What is the name of the package?",
                    "What is the title of the document?",
                ],
            )
        else:
            document = Document.open(get_fixture_path())

        self.assertEqual(len(document.queries), 2)
        self.assertEqual(document.queries[0].result.answer, "Textractor")
        self.assertEqual(document.queries[1].result.answer, "Textractor Test Document")

    def test_bad_queries_as_strings(self):
        profile_name = "default"
        current_directory = os.path.abspath(os.path.dirname(__file__))

        if profile_name is None:
            raise InvalidProfileNameError(
                "Textractor could not be initialized. Populate profile_name with a valid input in tests/test_table.py."
            )

        if os.environ.get("CALL_TEXTRACT"):
            extractor = Textractor(profile_name=profile_name, kms_key_id="")
            document = extractor.analyze_document(
                file_source=os.path.join(current_directory, "fixtures/single-page-1.png"),
                features=[TextractFeatures.QUERIES],
                queries=[
                    "Lorem ipsum?",
                    "The quick brown fox jumps over the lazy dog?",
                ],
            )
        else:
            document = Document.open(get_fixture_path())

        self.assertEqual(len(document.queries), 2)
        self.assertEqual(document.queries[0].result, None)
        self.assertEqual(document.queries[1].result, None)

    @unittest.skipIf(not os.environ.get("CALL_TEXTRACT"), "Asynchronous requests can't be processed without calling Textract")
    def test_query_feature_without_queries(self):
        profile_name = "default"
        current_directory = os.path.abspath(os.path.dirname(__file__))

        if profile_name is None:
            raise InvalidProfileNameError(
                "Textractor could not be initialized. Populate profile_name with a valid input in tests/test_table.py."
            )

        extractor = Textractor(profile_name=profile_name, kms_key_id="")
        with self.assertRaises(InputError):
            document = extractor.analyze_document(
                file_source=os.path.join(current_directory, "fixtures/single-page-1.png"),
                features=[TextractFeatures.TABLES],
                queries=[
                    "Lorem ipsum?",
                    "The quick brown fox jumps over the lazy dog?",
                ],
            )
````

### FILE: `aws_textractor_official/upstream/tests/fixtures/saved_api_responses/test_queries_as_strings.json`
```yaml
block_id: "PYTHON-AWS-TEXTRACTOR-OFFICIAL:fixture-query:v1"
operation: CREATE
provenance: ADAPTED
source: "https://raw.githubusercontent.com/aws-samples/amazon-textract-textractor/8ea5f9ae65dbcbb2b83b409e0007c395e59b0e92/tests/fixtures/saved_api_responses/test_queries_as_strings.json"
license: "Apache-2.0"
sha256: "f3786e60a207e740a4e3cb35696d5292506d31d258f072e78e07e3a6bf0215ca"
variables: []
secrets_allowed: false
```
````json
{"DocumentMetadata": {"Pages": 1}, "Blocks": [{"BlockType": "PAGE", "Geometry": {"BoundingBox": {"Width": 1.0, "Height": 1.0, "Left": 0.0, "Top": 0.0}, "Polygon": [{"X": 1.5123288876258822e-16, "Y": 0.0}, {"X": 1.0, "Y": 9.916888369396957e-17}, {"X": 1.0, "Y": 1.0}, {"X": 0.0, "Y": 1.0}]}, "Id": "30b46029-8a47-4884-ad9d-0b6fb435735c", "Relationships": [{"Type": "CHILD", "Ids": ["29f01f14-f00c-47b9-bf46-d284a5a2ec90", "978923ae-ee3b-46d7-a6d3-a29695d8039e", "66a95274-f6a9-4acf-9c24-38a55799abe2", "c31fcec8-23a5-4626-bb01-4f7ff7e22a5f", "7729e43e-039e-470a-b0c4-91f94a0187c7", "c6bb32fb-c831-43cd-a1c2-c82a883317eb", "846d2dc8-9da6-4377-a466-6944019df070", "fc67b68e-1cd4-4486-a7d5-9f946a169d93", "a8a0b1f3-9120-43db-8fad-43207fa5ad4d", "ad4e292b-39c5-4825-a7a7-e7fd9bb67a9c", "17007e0c-a3ee-43b4-8376-446bc91b2ed9", "b30cf0e5-3cd4-4b46-ab7d-29adebf09459", "bb910e60-9781-468f-b90e-7152c1b50db1", "8149a609-fb2e-4f20-b182-0ebf321ce10f", "123e3c05-0221-4fcd-84b4-c16c2f6dd93f", "49345fe7-f32a-4a73-b6cf-e1c937e26e8e", "7d4b91be-ad5a-40d8-8508-d171d2dccb72", "8462b984-dcfe-4fb4-93a3-701ea6e4b83d", "a63eaf16-d783-4b07-914e-a273ed906efc", "37109986-e9d2-4591-b878-2158f347d821", "b259db36-7921-401b-a46b-5ba15d6f18b4", "f03c7eff-fb07-4e1c-a9c5-04fe9563770f", "7ab0fa34-6c44-40cc-8e40-c74b402078a1", "7bd750c9-8d59-4997-8bf8-27b2bba5a691", "f9a9a914-60be-4b6f-b08f-2e7a6f06fec3", "7f3cc0fb-761f-4d50-b4fe-a7654dad95ff", "c71b1a0f-754a-4393-8dc6-ec40cab251de", "295828e5-309f-491f-a647-1c60aac34d3d"]}]}, {"BlockType": "LINE", "Confidence": 99.69444274902344, "Text": "Textractor Test", "Geometry": {"BoundingBox": {"Width": 0.6742472052574158, "Height": 0.06476041674613953, "Left": 0.1625330150127411, "Top": 0.10424439609050751}, "Polygon": [{"X": 0.1625330150127411, "Y": 0.10424439609050751}, {"X": 0.8367802500724792, "Y": 0.10424439609050751}, {"X": 0.8367802500724792, "Y": 0.16900481283664703}, {"X": 0.1625330150127411, "Y": 0.16900481283664703}]}, "Id": "29f01f14-f00c-47b9-bf46-d284a5a2ec90", "Relationships": [{"Type": "CHILD", "Ids": ["480b12d6-6309-4af9-8000-7e4b1c30cd67", "acb132e8-1203-4d72-ab04-a1794758c1dc"]}]}, {"BlockType": "LINE", "Confidence": 99.75547790527344, "Text": "Document", "Geometry": {"BoundingBox": {"Width": 0.4699675440788269, "Height": 0.06192830950021744, "Left": 0.2708991765975952, "Top": 0.22141507267951965}, "Polygon": [{"X": 0.2708991765975952, "Y": 0.22141507267951965}, {"X": 0.7408667206764221, "Y": 0.22141507267951965}, {"X": 0.7408667206764221, "Y": 0.2833433747291565}, {"X": 0.2708991765975952, "Y": 0.2833433747291565}]}, "Id": "978923ae-ee3b-46d7-a6d3-a29695d8039e", "Relationships": [{"Type": "CHILD", "Ids": ["37b7684c-8e6c-4449-8717-5736be0ef822"]}]}, {"BlockType": "LINE", "Confidence": 99.35688018798828, "Text": "Page (1)", "Geometry": {"BoundingBox": {"Width": 0.18069839477539062, "Height": 0.044398076832294464, "Left": 0.11965136229991913, "Top": 0.33620190620422363}, "Polygon": [{"X": 0.11965136229991913, "Y": 0.33620190620422363}, {"X": 0.30034974217414856, "Y": 0.33620190620422363}, {"X": 0.30034974217414856, "Y": 0.3806000053882599}, {"X": 0.11965136229991913, "Y": 0.3806000053882599}]}, "Id": "66a95274-f6a9-4acf-9c24-38a55799abe2", "Relationships": [{"Type": "CHILD", "Ids": ["f2d444d1-0755-4596-be67-e2e02ea0d117", "a4ca5384-b3c9-4624-856f-91dc09d8c84b"]}]}, {"BlockType": "LINE", "Confidence": 99.55121612548828, "Text": "Key - Values", "Geometry": {"BoundingBox": {"Width": 0.17561770975589752, "Height": 0.027066431939601898, "Left": 0.11902542412281036, "Top": 0.4110378324985504}, "Polygon": [{"X": 0.11902542412281036, "Y": 0.4110378324985504}, {"X": 0.2946431338787079, "Y": 0.4110378324985504}, {"X": 0.2946431338787079, "Y": 0.4381042718887329}, {"X": 0.11902542412281036, "Y": 0.4381042718887329}]}, "Id": "c31fcec8-23a5-4626-bb01-4f7ff7e22a5f", "Relationships": [{"Type": "CHILD", "Ids": ["552d9836-72d7-4182-a768-bc7474c7276b", "3b8767a8-e390-4780-ae32-9224218344ea", "57f8e0e5-8e47-497b-af84-f017095b3b26"]}]}, {"BlockType": "LINE", "Confidence": 99.78067016601562, "Text": "Name of package: Textractor", "Geometry": {"BoundingBox": {"Width": 0.34047141671180725, "Height": 0.023087535053491592, "Left": 0.11843090504407883, "Top": 0.4672558903694153}, "Polygon": [{"X": 0.11843090504407883, "Y": 0.4672558903694153}, {"X": 0.4589023292064667, "Y": 0.4672558903694153}, {"X": 0.4589023292064667, "Y": 0.4903434216976166}, {"X": 0.11843090504407883, "Y": 0.4903434216976166}]}, "Id": "7729e43e-039e-470a-b0c4-91f94a0187c7", "Relationships": [{"Type": "CHILD", "Ids": ["8aa5df27-d813-4175-b1cc-b7d4e2bb8deb", "3ee141d6-066b-4c16-8400-ec942db23d5b", "c7daf41e-d64a-4152-aab8-2bcadeb45429", "4f2d93af-8c22-460c-8dd5-f825b69f47ce"]}]}, {"BlockType": "LINE", "Confidence": 98.81005859375, "Text": "Date : 08/14/2022", "Geometry": {"BoundingBox": {"Width": 0.19092203676700592, "Height": 0.019717741757631302, "Left": 0.11804135888814926, "Top": 0.49777790904045105}, "Polygon": [{"X": 0.11804135888814926, "Y": 0.49777790904045105}, {"X": 0.3089633882045746, "Y": 0.49777790904045105}, {"X": 0.3089633882045746, "Y": 0.5174956321716309}, {"X": 0.11804135888814926, "Y": 0.5174956321716309}]}, "Id": "c6bb32fb-c831-43cd-a1c2-c82a883317eb", "Relationships": [{"Type": "CHILD", "Ids": ["8dfba6a8-5353-41c4-9edc-7983f48a544d", "a8861aff-1b85-4043-9304-92c9a69d5753", "78a90c72-4cb6-4e85-acab-1a11bd3171f3"]}]}, {"BlockType": "LINE", "Confidence": 99.84102630615234, "Text": "Table 1", "Geometry": {"BoundingBox": {"Width": 0.10481191426515579, "Height": 0.02227928303182125, "Left": 0.11691198498010635, "Top": 0.564307451248169}, "Polygon": [{"X": 0.11691198498010635, "Y": 0.564307451248169}, {"X": 0.22172389924526215, "Y": 0.564307451248169}, {"X": 0.22172389924526215, "Y": 0.5865867733955383}, {"X": 0.11691198498010635, "Y": 0.5865867733955383}]}, "Id": "846d2dc8-9da6-4377-a466-6944019df070", "Relationships": [{"Type": "CHILD", "Ids": ["62bb7d20-6218-4d3d-b6ea-e2a117be1acf", "fe116077-9d89-4101-978a-d8b37e0dd105"]}]}, {"BlockType": "LINE", "Confidence": 98.48675537109375, "Text": "Cell 1", "Geometry": {"BoundingBox": {"Width": 0.042510487139225006, "Height": 0.0135804433375597, "Left": 0.12771032750606537, "Top": 0.6181319952011108}, "Polygon": [{"X": 0.12771032750606537, "Y": 0.6181319952011108}, {"X": 0.17022082209587097, "Y": 0.6181319952011108}, {"X": 0.17022082209587097, "Y": 0.6317124366760254}, {"X": 0.12771032750606537, "Y": 0.6317124366760254}]}, "Id": "fc67b68e-1cd4-4486-a7d5-9f946a169d93", "Relationships": [{"Type": "CHILD", "Ids": ["64f38e28-c7bf-42ef-b2b4-3de7cefc2bc4", "32d8446e-c503-4086-8952-3262399bb965"]}]}, {"BlockType": "LINE", "Confidence": 98.98334503173828, "Text": "Cell 2", "Geometry": {"BoundingBox": {"Width": 0.043671898543834686, "Height": 0.012698537670075893, "Left": 0.285184383392334, "Top": 0.6184539794921875}, "Polygon": [{"X": 0.285184383392334, "Y": 0.6184539794921875}, {"X": 0.32885628938674927, "Y": 0.6184539794921875}, {"X": 0.32885628938674927, "Y": 0.6311525106430054}, {"X": 0.285184383392334, "Y": 0.6311525106430054}]}, "Id": "a8a0b1f3-9120-43db-8fad-43207fa5ad4d", "Relationships": [{"Type": "CHILD", "Ids": ["ace0aa2d-c783-4d0b-a833-c7b6fd4f6546", "aab04d55-a7f0-4bde-bb02-acf550174be0"]}]}, {"BlockType": "LINE", "Confidence": 98.5614242553711, "Text": "Cell 4", "Geometry": {"BoundingBox": {"Width": 0.04424715414643288, "Height": 0.012623554095625877, "Left": 0.6017507314682007, "Top": 0.6185309290885925}, "Polygon": [{"X": 0.6017507314682007, "Y": 0.6185309290885925}, {"X": 0.6459978818893433, "Y": 0.6185309290885925}, {"X": 0.6459978818893433, "Y": 0.631154477596283}, {"X": 0.6017507314682007, "Y": 0.631154477596283}]}, "Id": "ad4e292b-39c5-4825-a7a7-e7fd9bb67a9c", "Relationships": [{"Type": "CHILD", "Ids": ["f677177d-b11f-429a-84eb-b4f7b06108cc", "d0fbb516-cc32-4a36-85dd-ddc09f30b4fd"]}]}, {"BlockType": "LINE", "Confidence": 99.19432067871094, "Text": "Cell 5", "Geometry": {"BoundingBox": {"Width": 0.043169427663087845, "Height": 0.01277545653283596, "Left": 0.7619276642799377, "Top": 0.6184312701225281}, "Polygon": [{"X": 0.7619276642799377, "Y": 0.6184312701225281}, {"X": 0.8050971031188965, "Y": 0.6184312701225281}, {"X": 0.8050971031188965, "Y": 0.6312066912651062}, {"X": 0.7619276642799377, "Y": 0.6312066912651062}]}, "Id": "17007e0c-a3ee-43b4-8376-446bc91b2ed9", "Relationships": [{"Type": "CHILD", "Ids": ["242ed01e-884c-471a-8a10-3fcc247332b9", "bf6b9584-692a-4c85-b82d-8842bd6cbd8d"]}]}, {"BlockType": "LINE", "Confidence": 98.72454071044922, "Text": "Cell 6", "Geometry": {"BoundingBox": {"Width": 0.04389132186770439, "Height": 0.012387072667479515, "Left": 0.12712755799293518, "Top": 0.6778441667556763}, "Polygon": [{"X": 0.12712755799293518, "Y": 0.6778441667556763}, {"X": 0.17101888358592987, "Y": 0.6778441667556763}, {"X": 0.17101888358592987, "Y": 0.6902312636375427}, {"X": 0.12712755799293518, "Y": 0.6902312636375427}]}, "Id": "b30cf0e5-3cd4-4b46-ab7d-29adebf09459", "Relationships": [{"Type": "CHILD", "Ids": ["3834f062-9432-43f7-a07b-1881f7d72395", "74975fb6-9c93-4d54-8341-1d635393755c"]}]}, {"BlockType": "LINE", "Confidence": 98.60836791992188, "Text": "Cell 7", "Geometry": {"BoundingBox": {"Width": 0.0437912791967392, "Height": 0.012690716423094273, "Left": 0.28514277935028076, "Top": 0.6777071952819824}, "Polygon": [{"X": 0.28514277935028076, "Y": 0.6777071952819824}, {"X": 0.32893407344818115, "Y": 0.6777071952819824}, {"X": 0.32893407344818115, "Y": 0.6903979182243347}, {"X": 0.28514277935028076, "Y": 0.6903979182243347}]}, "Id": "bb910e60-9781-468f-b90e-7152c1b50db1", "Relationships": [{"Type": "CHILD", "Ids": ["d1d2f4cd-6824-416c-a5da-75f94b58545e", "48d5a449-1857-459e-85ba-e452f742c60c"]}]}, {"BlockType": "LINE", "Confidence": 98.65631866455078, "Text": "Cell 8", "Geometry": {"BoundingBox": {"Width": 0.04380646347999573, "Height": 0.01273917406797409, "Left": 0.4434744715690613, "Top": 0.6777850389480591}, "Polygon": [{"X": 0.4434744715690613, "Y": 0.6777850389480591}, {"X": 0.487280935049057, "Y": 0.6777850389480591}, {"X": 0.487280935049057, "Y": 0.690524160861969}, {"X": 0.4434744715690613, "Y": 0.690524160861969}]}, "Id": "8149a609-fb2e-4f20-b182-0ebf321ce10f", "Relationships": [{"Type": "CHILD", "Ids": ["18414781-1d95-49cd-b64b-023c05293899", "fa847736-b075-4b0d-b68a-3e1c14dd4d40"]}]}, {"BlockType": "LINE", "Confidence": 97.89591217041016, "Text": "Cell 9", "Geometry": {"BoundingBox": {"Width": 0.04354380443692207, "Height": 0.012920425273478031, "Left": 0.6017784476280212, "Top": 0.6776207089424133}, "Polygon": [{"X": 0.6017784476280212, "Y": 0.6776207089424133}, {"X": 0.6453222632408142, "Y": 0.6776207089424133}, {"X": 0.6453222632408142, "Y": 0.6905410885810852}, {"X": 0.6017784476280212, "Y": 0.6905410885810852}]}, "Id": "123e3c05-0221-4fcd-84b4-c16c2f6dd93f", "Relationships": [{"Type": "CHILD", "Ids": ["9dd38e13-cdd7-4a9a-9fff-592782934506", "d9b47ec0-3af1-4e70-8d84-33c2ebd33de3"]}]}, {"BlockType": "LINE", "Confidence": 99.05433654785156, "Text": "Cell 10", "Geometry": {"BoundingBox": {"Width": 0.053810954093933105, "Height": 0.01283883024007082, "Left": 0.7615987062454224, "Top": 0.6777843832969666}, "Polygon": [{"X": 0.7615987062454224, "Y": 0.6777843832969666}, {"X": 0.8154096603393555, "Y": 0.6777843832969666}, {"X": 0.8154096603393555, "Y": 0.6906231641769409}, {"X": 0.7615987062454224, "Y": 0.6906231641769409}]}, "Id": "49345fe7-f32a-4a73-b6cf-e1c937e26e8e", "Relationships": [{"Type": "CHILD", "Ids": ["200a7229-7ff2-4e35-9c1d-0696274b796d", "248e3d1a-d551-49c4-b164-184409d5962e"]}]}, {"BlockType": "LINE", "Confidence": 98.55902862548828, "Text": "Cell 11", "Geometry": {"BoundingBox": {"Width": 0.05301501974463463, "Height": 0.012528574094176292, "Left": 0.12711328268051147, "Top": 0.7360577583312988}, "Polygon": [{"X": 0.12711328268051147, "Y": 0.7360577583312988}, {"X": 0.1801283061504364, "Y": 0.7360577583312988}, {"X": 0.1801283061504364, "Y": 0.7485863566398621}, {"X": 0.12711328268051147, "Y": 0.7485863566398621}]}, "Id": "7d4b91be-ad5a-40d8-8508-d171d2dccb72", "Relationships": [{"Type": "CHILD", "Ids": ["a4ca960c-d96a-4526-aa10-49757457ae69", "f4f8c4ad-baff-42a5-9072-07f6325bf8de"]}]}, {"BlockType": "LINE", "Confidence": 98.943115234375, "Text": "Cell 12", "Geometry": {"BoundingBox": {"Width": 0.05333862826228142, "Height": 0.01295428816229105, "Left": 0.2852138876914978, "Top": 0.7360171675682068}, "Polygon": [{"X": 0.2852138876914978, "Y": 0.7360171675682068}, {"X": 0.3385525047779083, "Y": 0.7360171675682068}, {"X": 0.3385525047779083, "Y": 0.7489714622497559}, {"X": 0.2852138876914978, "Y": 0.7489714622497559}]}, "Id": "8462b984-dcfe-4fb4-93a3-701ea6e4b83d", "Relationships": [{"Type": "CHILD", "Ids": ["d1cc5f8a-89d8-44f6-87dc-90c3e00d9ccc", "50e689c6-d97c-4901-ac56-be3a4b182276"]}]}, {"BlockType": "LINE", "Confidence": 98.80059051513672, "Text": "Cell 13", "Geometry": {"BoundingBox": {"Width": 0.053260792046785355, "Height": 0.012730708345770836, "Left": 0.44334903359413147, "Top": 0.7359303832054138}, "Polygon": [{"X": 0.44334903359413147, "Y": 0.7359303832054138}, {"X": 0.49660980701446533, "Y": 0.7359303832054138}, {"X": 0.49660980701446533, "Y": 0.7486611008644104}, {"X": 0.44334903359413147, "Y": 0.7486611008644104}]}, "Id": "a63eaf16-d783-4b07-914e-a273ed906efc", "Relationships": [{"Type": "CHILD", "Ids": ["6aeede61-b5ca-469c-a8f6-1811bf6b800e", "6e5dfc35-ac52-4406-b6ef-512f7b6e49df"]}]}, {"BlockType": "LINE", "Confidence": 98.6966552734375, "Text": "Cell 14", "Geometry": {"BoundingBox": {"Width": 0.05358792096376419, "Height": 0.012640646658837795, "Left": 0.6018773913383484, "Top": 0.7359676957130432}, "Polygon": [{"X": 0.6018773913383484, "Y": 0.7359676957130432}, {"X": 0.655465304851532, "Y": 0.7359676957130432}, {"X": 0.655465304851532, "Y": 0.7486083507537842}, {"X": 0.6018773913383484, "Y": 0.7486083507537842}]}, "Id": "37109986-e9d2-4591-b878-2158f347d821", "Relationships": [{"Type": "CHILD", "Ids": ["c50870b1-7818-4695-88ff-9b78eecb55e8", "38cd6b7b-0874-45f7-a772-23c14fae39e8"]}]}, {"BlockType": "LINE", "Confidence": 99.1570053100586, "Text": "Cell 15", "Geometry": {"BoundingBox": {"Width": 0.05318611487746239, "Height": 0.012859954498708248, "Left": 0.7616134285926819, "Top": 0.7359786033630371}, "Polygon": [{"X": 0.7616134285926819, "Y": 0.7359786033630371}, {"X": 0.8147995471954346, "Y": 0.7359786033630371}, {"X": 0.8147995471954346, "Y": 0.7488385438919067}, {"X": 0.7616134285926819, "Y": 0.7488385438919067}]}, "Id": "b259db36-7921-401b-a46b-5ba15d6f18b4", "Relationships": [{"Type": "CHILD", "Ids": ["b3cd0077-41a3-4497-a952-15047833a548", "8982a73a-c06c-4fbf-a44c-3cf5d38fbf17"]}]}, {"BlockType": "LINE", "Confidence": 99.71761322021484, "Text": "Selection Element", "Geometry": {"BoundingBox": {"Width": 0.263517826795578, "Height": 0.02159660868346691, "Left": 0.1177731454372406, "Top": 0.8170759677886963}, "Polygon": [{"X": 0.1177731454372406, "Y": 0.8170759677886963}, {"X": 0.3812909722328186, "Y": 0.8170759677886963}, {"X": 0.3812909722328186, "Y": 0.8386725783348083}, {"X": 0.1177731454372406, "Y": 0.8386725783348083}]}, "Id": "f03c7eff-fb07-4e1c-a9c5-04fe9563770f", "Relationships": [{"Type": "CHILD", "Ids": ["95362018-d0b7-4e89-9539-851b6d869c49", "5c22a862-7d39-4547-aa0d-ce06c69c786f"]}]}, {"BlockType": "LINE", "Confidence": 99.4539794921875, "Text": "Selected Checkbox", "Geometry": {"BoundingBox": {"Width": 0.1507243663072586, "Height": 0.012731917202472687, "Left": 0.19152064621448517, "Top": 0.8825798630714417}, "Polygon": [{"X": 0.19152064621448517, "Y": 0.8825798630714417}, {"X": 0.3422450125217438, "Y": 0.8825798630714417}, {"X": 0.3422450125217438, "Y": 0.8953118324279785}, {"X": 0.19152064621448517, "Y": 0.8953118324279785}]}, "Id": "7ab0fa34-6c44-40cc-8e40-c74b402078a1", "Relationships": [{"Type": "CHILD", "Ids": ["0d9f78b8-3ba4-4ad5-8310-4af1e5e82763", "d4e4859e-b688-40df-acee-a05d421a4006"]}]}, {"BlockType": "LINE", "Confidence": 99.18217468261719, "Text": "Un-Selected Checkbox", "Geometry": {"BoundingBox": {"Width": 0.17911794781684875, "Height": 0.012731756083667278, "Left": 0.1922188550233841, "Top": 0.9288517832756042}, "Polygon": [{"X": 0.1922188550233841, "Y": 0.9288517832756042}, {"X": 0.37133678793907166, "Y": 0.9288517832756042}, {"X": 0.37133678793907166, "Y": 0.9415835738182068}, {"X": 0.1922188550233841, "Y": 0.9415835738182068}]}, "Id": "7bd750c9-8d59-4997-8bf8-27b2bba5a691", "Relationships": [{"Type": "CHILD", "Ids": ["3e81643d-cfbd-4dd0-8764-936893a37f4a", "dd2a767e-24f7-492e-b0c2-0e6bd564ebc7"]}]}, {"BlockType": "WORD", "Confidence": 99.58283233642578, "Text": "Textractor", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.462296724319458, "Height": 0.06334945559501648, "Left": 0.1625330150127411, "Top": 0.10538586229085922}, "Polygon": [{"X": 0.1625330150127411, "Y": 0.10538586229085922}, {"X": 0.6248297095298767, "Y": 0.10538586229085922}, {"X": 0.6248297095298767, "Y": 0.1687353104352951}, {"X": 0.1625330150127411, "Y": 0.1687353104352951}]}, "Id": "480b12d6-6309-4af9-8000-7e4b1c30cd67"}, {"BlockType": "WORD", "Confidence": 99.80604553222656, "Text": "Test", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.18449851870536804, "Height": 0.06476041674613953, "Left": 0.6522817015647888, "Top": 0.10424439609050751}, "Polygon": [{"X": 0.6522817015647888, "Y": 0.10424439609050751}, {"X": 0.8367802500724792, "Y": 0.10424439609050751}, {"X": 0.8367802500724792, "Y": 0.16900481283664703}, {"X": 0.6522817015647888, "Y": 0.16900481283664703}]}, "Id": "acb132e8-1203-4d72-ab04-a1794758c1dc"}, {"BlockType": "WORD", "Confidence": 99.75547790527344, "Text": "Document", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.4699675440788269, "Height": 0.06192830950021744, "Left": 0.2708991765975952, "Top": 0.22141507267951965}, "Polygon": [{"X": 0.2708991765975952, "Y": 0.22141507267951965}, {"X": 0.7408667206764221, "Y": 0.22141507267951965}, {"X": 0.7408667206764221, "Y": 0.2833433747291565}, {"X": 0.2708991765975952, "Y": 0.2833433747291565}]}, "Id": "37b7684c-8e6c-4449-8717-5736be0ef822"}, {"BlockType": "WORD", "Confidence": 99.0321273803711, "Text": "Page", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.10762534290552139, "Height": 0.0416223481297493, "Left": 0.11965136229991913, "Top": 0.3389776349067688}, "Polygon": [{"X": 0.11965136229991913, "Y": 0.3389776349067688}, {"X": 0.22727669775485992, "Y": 0.3389776349067688}, {"X": 0.22727669775485992, "Y": 0.3806000053882599}, {"X": 0.11965136229991913, "Y": 0.3806000053882599}]}, "Id": "f2d444d1-0755-4596-be67-e2e02ea0d117"}, {"BlockType": "WORD", "Confidence": 99.68162536621094, "Text": "(1)", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.059845417737960815, "Height": 0.04347170516848564, "Left": 0.24050432443618774, "Top": 0.33620190620422363}, "Polygon": [{"X": 0.24050432443618774, "Y": 0.33620190620422363}, {"X": 0.30034974217414856, "Y": 0.33620190620422363}, {"X": 0.30034974217414856, "Y": 0.37967362999916077}, {"X": 0.24050432443618774, "Y": 0.37967362999916077}]}, "Id": "a4ca5384-b3c9-4624-856f-91dc09d8c84b"}, {"BlockType": "WORD", "Confidence": 99.45826721191406, "Text": "Key", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.05174700915813446, "Height": 0.025642264634370804, "Left": 0.11902542412281036, "Top": 0.4124619960784912}, "Polygon": [{"X": 0.11902542412281036, "Y": 0.4124619960784912}, {"X": 0.17077243328094482, "Y": 0.4124619960784912}, {"X": 0.17077243328094482, "Y": 0.4381042718887329}, {"X": 0.11902542412281036, "Y": 0.4381042718887329}]}, "Id": "552d9836-72d7-4182-a768-bc7474c7276b"}, {"BlockType": "WORD", "Confidence": 99.68306732177734, "Text": "-", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.011502363719046116, "Height": 0.004489026963710785, "Left": 0.17738384008407593, "Top": 0.42237207293510437}, "Polygon": [{"X": 0.17738384008407593, "Y": 0.42237207293510437}, {"X": 0.18888619542121887, "Y": 0.42237207293510437}, {"X": 0.18888619542121887, "Y": 0.42686110734939575}, {"X": 0.17738384008407593, "Y": 0.42686110734939575}]}, "Id": "3b8767a8-e390-4780-ae32-9224218344ea"}, {"BlockType": "WORD", "Confidence": 99.5123291015625, "Text": "Values", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.0982576310634613, "Height": 0.02232169359922409, "Left": 0.19638550281524658, "Top": 0.4110378324985504}, "Polygon": [{"X": 0.19638550281524658, "Y": 0.4110378324985504}, {"X": 0.2946431338787079, "Y": 0.4110378324985504}, {"X": 0.2946431338787079, "Y": 0.4333595335483551}, {"X": 0.19638550281524658, "Y": 0.4333595335483551}]}, "Id": "57f8e0e5-8e47-497b-af84-f017095b3b26"}, {"BlockType": "WORD", "Confidence": 99.92672729492188, "Text": "Name", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.06867816299200058, "Height": 0.0178519356995821, "Left": 0.11843090504407883, "Top": 0.4683231711387634}, "Polygon": [{"X": 0.11843090504407883, "Y": 0.4683231711387634}, {"X": 0.1871090680360794, "Y": 0.4683231711387634}, {"X": 0.1871090680360794, "Y": 0.4861750900745392}, {"X": 0.11843090504407883, "Y": 0.4861750900745392}]}, "Id": "8aa5df27-d813-4175-b1cc-b7d4e2bb8deb"}, {"BlockType": "WORD", "Confidence": 99.99512481689453, "Text": "of", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.025785064324736595, "Height": 0.01896621100604534, "Left": 0.19438794255256653, "Top": 0.4672558903694153}, "Polygon": [{"X": 0.19438794255256653, "Y": 0.4672558903694153}, {"X": 0.22017300128936768, "Y": 0.4672558903694153}, {"X": 0.22017300128936768, "Y": 0.4862220883369446}, {"X": 0.19438794255256653, "Y": 0.4862220883369446}]}, "Id": "3ee141d6-066b-4c16-8400-ec942db23d5b"}, {"BlockType": "WORD", "Confidence": 99.57273864746094, "Text": "package:", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.10313721746206284, "Height": 0.02287496067583561, "Left": 0.22502322494983673, "Top": 0.4674684703350067}, "Polygon": [{"X": 0.22502322494983673, "Y": 0.4674684703350067}, {"X": 0.32816043496131897, "Y": 0.4674684703350067}, {"X": 0.32816043496131897, "Y": 0.4903434216976166}, {"X": 0.22502322494983673, "Y": 0.4903434216976166}]}, "Id": "c7daf41e-d64a-4152-aab8-2bcadeb45429"}, {"BlockType": "WORD", "Confidence": 99.6280746459961, "Text": "Textractor", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.12362529337406158, "Height": 0.01760973036289215, "Left": 0.3352770507335663, "Top": 0.46857234835624695}, "Polygon": [{"X": 0.3352770507335663, "Y": 0.46857234835624695}, {"X": 0.4589023292064667, "Y": 0.46857234835624695}, {"X": 0.4589023292064667, "Y": 0.4861820638179779}, {"X": 0.3352770507335663, "Y": 0.4861820638179779}]}, "Id": "4f2d93af-8c22-460c-8dd5-f825b69f47ce"}, {"BlockType": "WORD", "Confidence": 99.77838134765625, "Text": "Date", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.05016987770795822, "Height": 0.016659211367368698, "Left": 0.11804135888814926, "Top": 0.4987580180168152}, "Polygon": [{"X": 0.11804135888814926, "Y": 0.4987580180168152}, {"X": 0.16821123659610748, "Y": 0.4987580180168152}, {"X": 0.16821123659610748, "Y": 0.515417218208313}, {"X": 0.11804135888814926, "Y": 0.515417218208313}]}, "Id": "8dfba6a8-5353-41c4-9edc-7983f48a544d"}, {"BlockType": "WORD", "Confidence": 97.09461975097656, "Text": ":", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.00551046384498477, "Height": 0.012295318767428398, "Left": 0.17386683821678162, "Top": 0.502734363079071}, "Polygon": [{"X": 0.17386683821678162, "Y": 0.502734363079071}, {"X": 0.17937730252742767, "Y": 0.502734363079071}, {"X": 0.17937730252742767, "Y": 0.5150296688079834}, {"X": 0.17386683821678162, "Y": 0.5150296688079834}]}, "Id": "a8861aff-1b85-4043-9304-92c9a69d5753"}, {"BlockType": "WORD", "Confidence": 99.55718231201172, "Text": "08/14/2022", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.12413974851369858, "Height": 0.019717741757631302, "Left": 0.1848236471414566, "Top": 0.49777790904045105}, "Polygon": [{"X": 0.1848236471414566, "Y": 0.49777790904045105}, {"X": 0.3089633882045746, "Y": 0.49777790904045105}, {"X": 0.3089633882045746, "Y": 0.5174956321716309}, {"X": 0.1848236471414566, "Y": 0.5174956321716309}]}, "Id": "78a90c72-4cb6-4e85-acab-1a11bd3171f3"}, {"BlockType": "WORD", "Confidence": 99.86399841308594, "Text": "Table", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.07890205085277557, "Height": 0.02227928303182125, "Left": 0.11691198498010635, "Top": 0.564307451248169}, "Polygon": [{"X": 0.11691198498010635, "Y": 0.564307451248169}, {"X": 0.19581404328346252, "Y": 0.564307451248169}, {"X": 0.19581404328346252, "Y": 0.5865867733955383}, {"X": 0.11691198498010635, "Y": 0.5865867733955383}]}, "Id": "62bb7d20-6218-4d3d-b6ea-e2a117be1acf"}, {"BlockType": "WORD", "Confidence": 99.81806182861328, "Text": "1", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.016153108328580856, "Height": 0.020929092541337013, "Left": 0.205570787191391, "Top": 0.5654847621917725}, "Polygon": [{"X": 0.205570787191391, "Y": 0.5654847621917725}, {"X": 0.22172389924526215, "Y": 0.5654847621917725}, {"X": 0.22172389924526215, "Y": 0.5864138603210449}, {"X": 0.205570787191391, "Y": 0.5864138603210449}]}, "Id": "fe116077-9d89-4101-978a-d8b37e0dd105"}, {"BlockType": "WORD", "Confidence": 97.53768920898438, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.02881036140024662, "Height": 0.0135804433375597, "Left": 0.12771032750606537, "Top": 0.6181319952011108}, "Polygon": [{"X": 0.12771032750606537, "Y": 0.6181319952011108}, {"X": 0.15652069449424744, "Y": 0.6181319952011108}, {"X": 0.15652069449424744, "Y": 0.6317124366760254}, {"X": 0.12771032750606537, "Y": 0.6317124366760254}]}, "Id": "64f38e28-c7bf-42ef-b2b4-3de7cefc2bc4"}, {"BlockType": "WORD", "Confidence": 99.43582153320312, "Text": "1", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.009144669398665428, "Height": 0.011633284389972687, "Left": 0.1610761433839798, "Top": 0.6190181374549866}, "Polygon": [{"X": 0.1610761433839798, "Y": 0.6190181374549866}, {"X": 0.17022082209587097, "Y": 0.6190181374549866}, {"X": 0.17022082209587097, "Y": 0.6306514143943787}, {"X": 0.1610761433839798, "Y": 0.6306514143943787}]}, "Id": "32d8446e-c503-4086-8952-3262399bb965"}, {"BlockType": "WORD", "Confidence": 98.71371459960938, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.029801635071635246, "Height": 0.012698537670075893, "Left": 0.285184383392334, "Top": 0.6184539794921875}, "Polygon": [{"X": 0.285184383392334, "Y": 0.6184539794921875}, {"X": 0.3149860203266144, "Y": 0.6184539794921875}, {"X": 0.3149860203266144, "Y": 0.6311525106430054}, {"X": 0.285184383392334, "Y": 0.6311525106430054}]}, "Id": "ace0aa2d-c783-4d0b-a833-c7b6fd4f6546"}, {"BlockType": "WORD", "Confidence": 99.25298309326172, "Text": "2", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.010345458984375, "Height": 0.012169862166047096, "Left": 0.31851083040237427, "Top": 0.6186932921409607}, "Polygon": [{"X": 0.31851083040237427, "Y": 0.6186932921409607}, {"X": 0.32885628938674927, "Y": 0.6186932921409607}, {"X": 0.32885628938674927, "Y": 0.6308631896972656}, {"X": 0.31851083040237427, "Y": 0.6308631896972656}]}, "Id": "aab04d55-a7f0-4bde-bb02-acf550174be0"}, {"BlockType": "WORD", "Confidence": 98.03536224365234, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.029655318707227707, "Height": 0.012623554095625877, "Left": 0.6017507314682007, "Top": 0.6185309290885925}, "Polygon": [{"X": 0.6017507314682007, "Y": 0.6185309290885925}, {"X": 0.6314060688018799, "Y": 0.6185309290885925}, {"X": 0.6314060688018799, "Y": 0.631154477596283}, {"X": 0.6017507314682007, "Y": 0.631154477596283}]}, "Id": "f677177d-b11f-429a-84eb-b4f7b06108cc"}, {"BlockType": "WORD", "Confidence": 99.08748626708984, "Text": "4", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.011248490773141384, "Height": 0.011445663869380951, "Left": 0.6347494125366211, "Top": 0.618995726108551}, "Polygon": [{"X": 0.6347494125366211, "Y": 0.618995726108551}, {"X": 0.6459978818893433, "Y": 0.618995726108551}, {"X": 0.6459978818893433, "Y": 0.630441427230835}, {"X": 0.6347494125366211, "Y": 0.630441427230835}]}, "Id": "d0fbb516-cc32-4a36-85dd-ddc09f30b4fd"}, {"BlockType": "WORD", "Confidence": 98.73052978515625, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.029450906440615654, "Height": 0.01277545653283596, "Left": 0.7619276642799377, "Top": 0.6184312701225281}, "Polygon": [{"X": 0.7619276642799377, "Y": 0.6184312701225281}, {"X": 0.7913786172866821, "Y": 0.6184312701225281}, {"X": 0.7913786172866821, "Y": 0.6312066912651062}, {"X": 0.7619276642799377, "Y": 0.6312066912651062}]}, "Id": "242ed01e-884c-471a-8a10-3fcc247332b9"}, {"BlockType": "WORD", "Confidence": 99.65811920166016, "Text": "5", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.01000399049371481, "Height": 0.012092862278223038, "Left": 0.7950931191444397, "Top": 0.6188316345214844}, "Polygon": [{"X": 0.7950931191444397, "Y": 0.6188316345214844}, {"X": 0.8050971031188965, "Y": 0.6188316345214844}, {"X": 0.8050971031188965, "Y": 0.6309244632720947}, {"X": 0.7950931191444397, "Y": 0.6309244632720947}]}, "Id": "bf6b9584-692a-4c85-b82d-8842bd6cbd8d"}, {"BlockType": "WORD", "Confidence": 97.86188507080078, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.028879188001155853, "Height": 0.012308783829212189, "Left": 0.12712755799293518, "Top": 0.6778441667556763}, "Polygon": [{"X": 0.12712755799293518, "Y": 0.6778441667556763}, {"X": 0.15600673854351044, "Y": 0.6778441667556763}, {"X": 0.15600673854351044, "Y": 0.6901530027389526}, {"X": 0.12712755799293518, "Y": 0.6901530027389526}]}, "Id": "3834f062-9432-43f7-a07b-1881f7d72395"}, {"BlockType": "WORD", "Confidence": 99.58719635009766, "Text": "6", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.010334332473576069, "Height": 0.011654086410999298, "Left": 0.16068454086780548, "Top": 0.678577184677124}, "Polygon": [{"X": 0.16068454086780548, "Y": 0.678577184677124}, {"X": 0.17101888358592987, "Y": 0.678577184677124}, {"X": 0.17101888358592987, "Y": 0.6902312636375427}, {"X": 0.16068454086780548, "Y": 0.6902312636375427}]}, "Id": "74975fb6-9c93-4d54-8341-1d635393755c"}, {"BlockType": "WORD", "Confidence": 98.22537994384766, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.030321579426527023, "Height": 0.0126646738499403, "Left": 0.28514277935028076, "Top": 0.6777071952819824}, "Polygon": [{"X": 0.28514277935028076, "Y": 0.6777071952819824}, {"X": 0.3154643476009369, "Y": 0.6777071952819824}, {"X": 0.3154643476009369, "Y": 0.6903718709945679}, {"X": 0.28514277935028076, "Y": 0.6903718709945679}]}, "Id": "d1d2f4cd-6824-416c-a5da-75f94b58545e"}, {"BlockType": "WORD", "Confidence": 98.9913558959961, "Text": "7", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.010197600349783897, "Height": 0.01210785936564207, "Left": 0.3187364637851715, "Top": 0.678290069103241}, "Polygon": [{"X": 0.3187364637851715, "Y": 0.678290069103241}, {"X": 0.32893407344818115, "Y": 0.678290069103241}, {"X": 0.32893407344818115, "Y": 0.6903979182243347}, {"X": 0.3187364637851715, "Y": 0.6903979182243347}]}, "Id": "48d5a449-1857-459e-85ba-e452f742c60c"}, {"BlockType": "WORD", "Confidence": 98.12083435058594, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.029959699138998985, "Height": 0.01273917406797409, "Left": 0.4434744715690613, "Top": 0.6777850389480591}, "Polygon": [{"X": 0.4434744715690613, "Y": 0.6777850389480591}, {"X": 0.473434180021286, "Y": 0.6777850389480591}, {"X": 0.473434180021286, "Y": 0.690524160861969}, {"X": 0.4434744715690613, "Y": 0.690524160861969}]}, "Id": "18414781-1d95-49cd-b64b-023c05293899"}, {"BlockType": "WORD", "Confidence": 99.19180297851562, "Text": "8", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.010379361920058727, "Height": 0.011977242305874825, "Left": 0.4769015908241272, "Top": 0.6784162521362305}, "Polygon": [{"X": 0.4769015908241272, "Y": 0.6784162521362305}, {"X": 0.487280935049057, "Y": 0.6784162521362305}, {"X": 0.487280935049057, "Y": 0.6903935074806213}, {"X": 0.4769015908241272, "Y": 0.6903935074806213}]}, "Id": "fa847736-b075-4b0d-b68a-3e1c14dd4d40"}, {"BlockType": "WORD", "Confidence": 96.94509887695312, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.030068129301071167, "Height": 0.012920425273478031, "Left": 0.6017784476280212, "Top": 0.6776207089424133}, "Polygon": [{"X": 0.6017784476280212, "Y": 0.6776207089424133}, {"X": 0.6318466067314148, "Y": 0.6776207089424133}, {"X": 0.6318466067314148, "Y": 0.6905410885810852}, {"X": 0.6017784476280212, "Y": 0.6905410885810852}]}, "Id": "9dd38e13-cdd7-4a9a-9fff-592782934506"}, {"BlockType": "WORD", "Confidence": 98.84673309326172, "Text": "9", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.010458518750965595, "Height": 0.011845496483147144, "Left": 0.6348637342453003, "Top": 0.6785191297531128}, "Polygon": [{"X": 0.6348637342453003, "Y": 0.6785191297531128}, {"X": 0.6453222632408142, "Y": 0.6785191297531128}, {"X": 0.6453222632408142, "Y": 0.6903645992279053}, {"X": 0.6348637342453003, "Y": 0.6903645992279053}]}, "Id": "d9b47ec0-3af1-4e70-8d84-33c2ebd33de3"}, {"BlockType": "WORD", "Confidence": 98.36264038085938, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.029403511434793472, "Height": 0.012695634737610817, "Left": 0.7615987062454224, "Top": 0.6777843832969666}, "Polygon": [{"X": 0.7615987062454224, "Y": 0.6777843832969666}, {"X": 0.7910022139549255, "Y": 0.6777843832969666}, {"X": 0.7910022139549255, "Y": 0.6904799938201904}, {"X": 0.7615987062454224, "Y": 0.6904799938201904}]}, "Id": "200a7229-7ff2-4e35-9c1d-0696274b796d"}, {"BlockType": "WORD", "Confidence": 99.74604034423828, "Text": "10", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.020036855712532997, "Height": 0.012323135510087013, "Left": 0.7953728437423706, "Top": 0.6783000826835632}, "Polygon": [{"X": 0.7953728437423706, "Y": 0.6783000826835632}, {"X": 0.8154096603393555, "Y": 0.6783000826835632}, {"X": 0.8154096603393555, "Y": 0.6906231641769409}, {"X": 0.7953728437423706, "Y": 0.6906231641769409}]}, "Id": "248e3d1a-d551-49c4-b164-184409d5962e"}, {"BlockType": "WORD", "Confidence": 97.57048797607422, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.029744284227490425, "Height": 0.012400537729263306, "Left": 0.12711328268051147, "Top": 0.7360577583312988}, "Polygon": [{"X": 0.12711328268051147, "Y": 0.7360577583312988}, {"X": 0.15685756504535675, "Y": 0.7360577583312988}, {"X": 0.15685756504535675, "Y": 0.7484583258628845}, {"X": 0.12711328268051147, "Y": 0.7484583258628845}]}, "Id": "a4ca960c-d96a-4526-aa10-49757457ae69"}, {"BlockType": "WORD", "Confidence": 99.54756164550781, "Text": "11", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.019168399274349213, "Height": 0.012011912651360035, "Left": 0.1609598994255066, "Top": 0.7365744709968567}, "Polygon": [{"X": 0.1609598994255066, "Y": 0.7365744709968567}, {"X": 0.1801283061504364, "Y": 0.7365744709968567}, {"X": 0.1801283061504364, "Y": 0.7485863566398621}, {"X": 0.1609598994255066, "Y": 0.7485863566398621}]}, "Id": "f4f8c4ad-baff-42a5-9072-07f6325bf8de"}, {"BlockType": "WORD", "Confidence": 98.26697540283203, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.02954833395779133, "Height": 0.01248890534043312, "Left": 0.2852138876914978, "Top": 0.7360171675682068}, "Polygon": [{"X": 0.2852138876914978, "Y": 0.7360171675682068}, {"X": 0.3147622048854828, "Y": 0.7360171675682068}, {"X": 0.3147622048854828, "Y": 0.7485060691833496}, {"X": 0.2852138876914978, "Y": 0.7485060691833496}]}, "Id": "d1cc5f8a-89d8-44f6-87dc-90c3e00d9ccc"}, {"BlockType": "WORD", "Confidence": 99.61925506591797, "Text": "12", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.01955912820994854, "Height": 0.012924294918775558, "Left": 0.31899335980415344, "Top": 0.7360471487045288}, "Polygon": [{"X": 0.31899335980415344, "Y": 0.7360471487045288}, {"X": 0.3385525047779083, "Y": 0.7360471487045288}, {"X": 0.3385525047779083, "Y": 0.7489714622497559}, {"X": 0.31899335980415344, "Y": 0.7489714622497559}]}, "Id": "50e689c6-d97c-4901-ac56-be3a4b182276"}, {"BlockType": "WORD", "Confidence": 97.85566711425781, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.029667267575860023, "Height": 0.01260742824524641, "Left": 0.44334903359413147, "Top": 0.7359303832054138}, "Polygon": [{"X": 0.44334903359413147, "Y": 0.7359303832054138}, {"X": 0.47301629185676575, "Y": 0.7359303832054138}, {"X": 0.47301629185676575, "Y": 0.7485378384590149}, {"X": 0.44334903359413147, "Y": 0.7485378384590149}]}, "Id": "6aeede61-b5ca-469c-a8f6-1811bf6b800e"}, {"BlockType": "WORD", "Confidence": 99.74552154541016, "Text": "13", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.019049739465117455, "Height": 0.012203322723507881, "Left": 0.4775600731372833, "Top": 0.7364577651023865}, "Polygon": [{"X": 0.4775600731372833, "Y": 0.7364577651023865}, {"X": 0.49660980701446533, "Y": 0.7364577651023865}, {"X": 0.49660980701446533, "Y": 0.7486611008644104}, {"X": 0.4775600731372833, "Y": 0.7486611008644104}]}, "Id": "6e5dfc35-ac52-4406-b6ef-512f7b6e49df"}, {"BlockType": "WORD", "Confidence": 97.64884185791016, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.02963610365986824, "Height": 0.012640646658837795, "Left": 0.6018773913383484, "Top": 0.7359676957130432}, "Polygon": [{"X": 0.6018773913383484, "Y": 0.7359676957130432}, {"X": 0.6315134763717651, "Y": 0.7359676957130432}, {"X": 0.6315134763717651, "Y": 0.7486083507537842}, {"X": 0.6018773913383484, "Y": 0.7486083507537842}]}, "Id": "c50870b1-7818-4695-88ff-9b78eecb55e8"}, {"BlockType": "WORD", "Confidence": 99.74446868896484, "Text": "14", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.019679507240653038, "Height": 0.011959020048379898, "Left": 0.6357858180999756, "Top": 0.7366405725479126}, "Polygon": [{"X": 0.6357858180999756, "Y": 0.7366405725479126}, {"X": 0.655465304851532, "Y": 0.7366405725479126}, {"X": 0.655465304851532, "Y": 0.7485995888710022}, {"X": 0.6357858180999756, "Y": 0.7485995888710022}]}, "Id": "38cd6b7b-0874-45f7-a772-23c14fae39e8"}, {"BlockType": "WORD", "Confidence": 98.42084503173828, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.029446326196193695, "Height": 0.012561631388962269, "Left": 0.7616134285926819, "Top": 0.7359786033630371}, "Polygon": [{"X": 0.7616134285926819, "Y": 0.7359786033630371}, {"X": 0.7910597324371338, "Y": 0.7359786033630371}, {"X": 0.7910597324371338, "Y": 0.7485402226448059}, {"X": 0.7616134285926819, "Y": 0.7485402226448059}]}, "Id": "b3cd0077-41a3-4497-a952-15047833a548"}, {"BlockType": "WORD", "Confidence": 99.8931655883789, "Text": "15", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.01951681263744831, "Height": 0.012637583538889885, "Left": 0.7952827215194702, "Top": 0.7362009882926941}, "Polygon": [{"X": 0.7952827215194702, "Y": 0.7362009882926941}, {"X": 0.8147995471954346, "Y": 0.7362009882926941}, {"X": 0.8147995471954346, "Y": 0.7488385438919067}, {"X": 0.7952827215194702, "Y": 0.7488385438919067}]}, "Id": "8982a73a-c06c-4fbf-a44c-3cf5d38fbf17"}, {"BlockType": "WORD", "Confidence": 99.82105255126953, "Text": "Selection", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.1316467672586441, "Height": 0.02155419811606407, "Left": 0.1177731454372406, "Top": 0.8171183466911316}, "Polygon": [{"X": 0.1177731454372406, "Y": 0.8171183466911316}, {"X": 0.2494199126958847, "Y": 0.8171183466911316}, {"X": 0.2494199126958847, "Y": 0.8386725783348083}, {"X": 0.1177731454372406, "Y": 0.8386725783348083}]}, "Id": "95362018-d0b7-4e89-9539-851b6d869c49"}, {"BlockType": "WORD", "Confidence": 99.61417388916016, "Text": "Element", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.12110535055398941, "Height": 0.021324409171938896, "Left": 0.2601856291294098, "Top": 0.8170759677886963}, "Polygon": [{"X": 0.2601856291294098, "Y": 0.8170759677886963}, {"X": 0.3812909722328186, "Y": 0.8170759677886963}, {"X": 0.3812909722328186, "Y": 0.8384003639221191}, {"X": 0.2601856291294098, "Y": 0.8384003639221191}]}, "Id": "5c22a862-7d39-4547-aa0d-ce06c69c786f"}, {"BlockType": "WORD", "Confidence": 99.92532348632812, "Text": "Selected", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.06861749291419983, "Height": 0.012586626224219799, "Left": 0.19152064621448517, "Top": 0.8826417922973633}, "Polygon": [{"X": 0.19152064621448517, "Y": 0.8826417922973633}, {"X": 0.2601381242275238, "Y": 0.8826417922973633}, {"X": 0.2601381242275238, "Y": 0.8952284455299377}, {"X": 0.19152064621448517, "Y": 0.8952284455299377}]}, "Id": "0d9f78b8-3ba4-4ad5-8310-4af1e5e82763"}, {"BlockType": "WORD", "Confidence": 98.98263549804688, "Text": "Checkbox", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.07755771279335022, "Height": 0.012731917202472687, "Left": 0.26468729972839355, "Top": 0.8825798630714417}, "Polygon": [{"X": 0.26468729972839355, "Y": 0.8825798630714417}, {"X": 0.3422450125217438, "Y": 0.8825798630714417}, {"X": 0.3422450125217438, "Y": 0.8953118324279785}, {"X": 0.26468729972839355, "Y": 0.8953118324279785}]}, "Id": "d4e4859e-b688-40df-acee-a05d421a4006"}, {"BlockType": "WORD", "Confidence": 99.56431579589844, "Text": "Un-Selected", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.09667694568634033, "Height": 0.012559858150780201, "Left": 0.1922188550233841, "Top": 0.9289750456809998}, "Polygon": [{"X": 0.1922188550233841, "Y": 0.9289750456809998}, {"X": 0.28889578580856323, "Y": 0.9289750456809998}, {"X": 0.28889578580856323, "Y": 0.9415349364280701}, {"X": 0.1922188550233841, "Y": 0.9415349364280701}]}, "Id": "3e81643d-cfbd-4dd0-8764-936893a37f4a"}, {"BlockType": "WORD", "Confidence": 98.800048828125, "Text": "Checkbox", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.07832921296358109, "Height": 0.012731756083667278, "Left": 0.29300758242607117, "Top": 0.9288517832756042}, "Polygon": [{"X": 0.29300758242607117, "Y": 0.9288517832756042}, {"X": 0.37133678793907166, "Y": 0.9288517832756042}, {"X": 0.37133678793907166, "Y": 0.9415835738182068}, {"X": 0.29300758242607117, "Y": 0.9415835738182068}]}, "Id": "dd2a767e-24f7-492e-b0c2-0e6bd564ebc7"}, {"BlockType": "QUERY", "Id": "f9a9a914-60be-4b6f-b08f-2e7a6f06fec3", "Relationships": [{"Type": "ANSWER", "Ids": ["7f3cc0fb-761f-4d50-b4fe-a7654dad95ff"]}], "Query": {"Text": "What is the name of the package?"}}, {"BlockType": "QUERY_RESULT", "Confidence": 92.0, "Text": "Textractor", "Geometry": {"BoundingBox": {"Width": 0.12316476553678513, "Height": 0.01783355325460434, "Left": 0.3352365493774414, "Top": 0.4682959020137787}, "Polygon": [{"X": 0.3352365493774414, "Y": 0.4682959020137787}, {"X": 0.45840129256248474, "Y": 0.4682959020137787}, {"X": 0.45840129256248474, "Y": 0.4861294627189636}, {"X": 0.3352365493774414, "Y": 0.4861294627189636}]}, "Id": "7f3cc0fb-761f-4d50-b4fe-a7654dad95ff"}, {"BlockType": "QUERY", "Id": "c71b1a0f-754a-4393-8dc6-ec40cab251de", "Relationships": [{"Type": "ANSWER", "Ids": ["295828e5-309f-491f-a647-1c60aac34d3d"]}], "Query": {"Text": "What is the title of the document?"}}, {"BlockType": "QUERY_RESULT", "Confidence": 87.0, "Text": "Textractor Test Document", "Geometry": {"BoundingBox": {"Width": 0.6737357378005981, "Height": 0.1789960414171219, "Left": 0.1623164713382721, "Top": 0.10369881242513657}, "Polygon": [{"X": 0.1623164713382721, "Y": 0.10369881242513657}, {"X": 0.8360521793365479, "Y": 0.10369881242513657}, {"X": 0.8360521793365479, "Y": 0.28269484639167786}, {"X": 0.1623164713382721, "Y": 0.28269484639167786}]}, "Id": "295828e5-309f-491f-a647-1c60aac34d3d"}], "AnalyzeDocumentModelVersion": "1.0", "ResponseMetadata": {"RequestId": "47833de2-c87d-4f1a-b92c-57db9d86e75b", "HTTPStatusCode": 200, "HTTPHeaders": {"x-amzn-requestid": "47833de2-c87d-4f1a-b92c-57db9d86e75b", "content-type": "application/x-amz-json-1.1", "content-length": "41060", "date": "Tue, 27 Sep 2022 18:12:54 GMT"}, "RetryAttempts": 0}}
````

### FILE: `aws_textractor_official/upstream/tests/fixtures/saved_api_responses/test_bad_queries_as_strings.json`
```yaml
block_id: "PYTHON-AWS-TEXTRACTOR-OFFICIAL:fixture-empty-query:v1"
operation: CREATE
provenance: ADAPTED
source: "https://raw.githubusercontent.com/aws-samples/amazon-textract-textractor/8ea5f9ae65dbcbb2b83b409e0007c395e59b0e92/tests/fixtures/saved_api_responses/test_bad_queries_as_strings.json"
license: "Apache-2.0"
sha256: "12e2327459483c96a05dfb9babe7968f3d01c5b199caeb94296c43a2561927a7"
variables: []
secrets_allowed: false
```
````json
{"DocumentMetadata": {"Pages": 1}, "Blocks": [{"BlockType": "PAGE", "Geometry": {"BoundingBox": {"Width": 1.0, "Height": 1.0, "Left": 0.0, "Top": 0.0}, "Polygon": [{"X": 1.5123288876258822e-16, "Y": 0.0}, {"X": 1.0, "Y": 9.916888369396957e-17}, {"X": 1.0, "Y": 1.0}, {"X": 0.0, "Y": 1.0}]}, "Id": "29a76e43-8d1f-43fc-b817-4053cc89a5c6", "Relationships": [{"Type": "CHILD", "Ids": ["cb448d96-1467-45f5-b36d-85f272a0e5c7", "e0338673-c792-4a55-b240-0eb4aadbb6c6", "3f4d8578-6505-48e6-aee6-e2c7e812f30e", "d0c53bc3-89be-429f-923f-204a33f59ce2", "79048f1b-c151-40f3-8d3e-695eb1be3f70", "8425abcd-d055-4d1f-9ae7-cc0645d8115f", "9ab057d0-7747-485a-97ab-960fa750af15", "c1020ebf-bed2-4212-b24c-5972272cc8d8", "69c4ecc9-d075-4b7d-a846-f254d1137582", "35a58188-40c3-4761-8ea7-f188b3732fd1", "b2f081a5-21c5-4763-b876-0313f11a7575", "96e1d611-3ffe-48ab-b1b9-c3d7086200e2", "12ff6f0b-7639-44ee-b197-d34cdb52e7aa", "c708dcc5-82ed-4c19-975e-d7ff860e6bdc", "feedc078-68c0-42b8-a87c-72e3bd80ee78", "1e3fb0bf-0f48-4e7a-9a3f-cc1b37480da6", "0769c3cf-e76e-4bae-bea1-8e3b1943a107", "4c9957c5-c3af-43a4-a1d4-f6c31ea4e32b", "876d2f64-150b-4aea-9b90-7cd9a3fd18d8", "edb9f59e-e615-4157-8949-22b8d3b1602c", "a8830a85-1f77-4e7f-9219-921c2da3abe0", "e8589aa0-8cfd-4a4d-8e32-49a25ccf37c8", "32ec6d5b-8c95-460e-a900-8c2998a194c1", "83c8002d-2a4d-4e64-a344-5139503120e9", "cea16e7e-6a41-4529-9290-452d6735b495", "4bbcd465-02fe-4745-bf44-a4822a09233d"]}]}, {"BlockType": "LINE", "Confidence": 99.69444274902344, "Text": "Textractor Test", "Geometry": {"BoundingBox": {"Width": 0.6742472052574158, "Height": 0.06476041674613953, "Left": 0.1625330150127411, "Top": 0.10424439609050751}, "Polygon": [{"X": 0.1625330150127411, "Y": 0.10424439609050751}, {"X": 0.8367802500724792, "Y": 0.10424439609050751}, {"X": 0.8367802500724792, "Y": 0.16900481283664703}, {"X": 0.1625330150127411, "Y": 0.16900481283664703}]}, "Id": "cb448d96-1467-45f5-b36d-85f272a0e5c7", "Relationships": [{"Type": "CHILD", "Ids": ["0b1cf325-b523-4b86-b01e-50ade0563079", "e6668699-4d9e-4f8b-a242-3b414a26d3fd"]}]}, {"BlockType": "LINE", "Confidence": 99.75547790527344, "Text": "Document", "Geometry": {"BoundingBox": {"Width": 0.4699675440788269, "Height": 0.06192830950021744, "Left": 0.2708991765975952, "Top": 0.22141507267951965}, "Polygon": [{"X": 0.2708991765975952, "Y": 0.22141507267951965}, {"X": 0.7408667206764221, "Y": 0.22141507267951965}, {"X": 0.7408667206764221, "Y": 0.2833433747291565}, {"X": 0.2708991765975952, "Y": 0.2833433747291565}]}, "Id": "e0338673-c792-4a55-b240-0eb4aadbb6c6", "Relationships": [{"Type": "CHILD", "Ids": ["b2ab72c0-0470-44bf-9fd6-9e84e2ae8eea"]}]}, {"BlockType": "LINE", "Confidence": 99.35688018798828, "Text": "Page (1)", "Geometry": {"BoundingBox": {"Width": 0.18069839477539062, "Height": 0.044398076832294464, "Left": 0.11965136229991913, "Top": 0.33620190620422363}, "Polygon": [{"X": 0.11965136229991913, "Y": 0.33620190620422363}, {"X": 0.30034974217414856, "Y": 0.33620190620422363}, {"X": 0.30034974217414856, "Y": 0.3806000053882599}, {"X": 0.11965136229991913, "Y": 0.3806000053882599}]}, "Id": "3f4d8578-6505-48e6-aee6-e2c7e812f30e", "Relationships": [{"Type": "CHILD", "Ids": ["dcde4029-984c-4795-9d95-2a42332046c2", "fd202ecd-c9d1-4080-925f-38775e904101"]}]}, {"BlockType": "LINE", "Confidence": 99.55121612548828, "Text": "Key - Values", "Geometry": {"BoundingBox": {"Width": 0.17561770975589752, "Height": 0.027066431939601898, "Left": 0.11902542412281036, "Top": 0.4110378324985504}, "Polygon": [{"X": 0.11902542412281036, "Y": 0.4110378324985504}, {"X": 0.2946431338787079, "Y": 0.4110378324985504}, {"X": 0.2946431338787079, "Y": 0.4381042718887329}, {"X": 0.11902542412281036, "Y": 0.4381042718887329}]}, "Id": "d0c53bc3-89be-429f-923f-204a33f59ce2", "Relationships": [{"Type": "CHILD", "Ids": ["773c791e-dee0-43a6-be2c-f7c7ee5ab8b3", "93dea5cc-f1ba-47d8-bf03-ddcbfdec16c2", "455c6f4d-e69d-4685-9296-6671c3013f41"]}]}, {"BlockType": "LINE", "Confidence": 99.78067016601562, "Text": "Name of package: Textractor", "Geometry": {"BoundingBox": {"Width": 0.34047141671180725, "Height": 0.023087535053491592, "Left": 0.11843090504407883, "Top": 0.4672558903694153}, "Polygon": [{"X": 0.11843090504407883, "Y": 0.4672558903694153}, {"X": 0.4589023292064667, "Y": 0.4672558903694153}, {"X": 0.4589023292064667, "Y": 0.4903434216976166}, {"X": 0.11843090504407883, "Y": 0.4903434216976166}]}, "Id": "79048f1b-c151-40f3-8d3e-695eb1be3f70", "Relationships": [{"Type": "CHILD", "Ids": ["5adabdea-adcd-4f4c-9000-95bdb979a1d6", "38788f7f-9475-4577-bd54-7bcd74393c59", "1fdff324-efd6-4038-81de-280812891135", "dd335722-cf8f-457c-8ebc-013a79baedd6"]}]}, {"BlockType": "LINE", "Confidence": 98.81005859375, "Text": "Date : 08/14/2022", "Geometry": {"BoundingBox": {"Width": 0.19092203676700592, "Height": 0.019717741757631302, "Left": 0.11804135888814926, "Top": 0.49777790904045105}, "Polygon": [{"X": 0.11804135888814926, "Y": 0.49777790904045105}, {"X": 0.3089633882045746, "Y": 0.49777790904045105}, {"X": 0.3089633882045746, "Y": 0.5174956321716309}, {"X": 0.11804135888814926, "Y": 0.5174956321716309}]}, "Id": "8425abcd-d055-4d1f-9ae7-cc0645d8115f", "Relationships": [{"Type": "CHILD", "Ids": ["9271e8b3-8ab5-4b9f-895f-4830af054ca8", "b9b96df2-7918-4cdc-bd0a-53ce974b0743", "5d235159-eed6-43e0-b017-788735b366d8"]}]}, {"BlockType": "LINE", "Confidence": 99.84102630615234, "Text": "Table 1", "Geometry": {"BoundingBox": {"Width": 0.10481191426515579, "Height": 0.02227928303182125, "Left": 0.11691198498010635, "Top": 0.564307451248169}, "Polygon": [{"X": 0.11691198498010635, "Y": 0.564307451248169}, {"X": 0.22172389924526215, "Y": 0.564307451248169}, {"X": 0.22172389924526215, "Y": 0.5865867733955383}, {"X": 0.11691198498010635, "Y": 0.5865867733955383}]}, "Id": "9ab057d0-7747-485a-97ab-960fa750af15", "Relationships": [{"Type": "CHILD", "Ids": ["47cef666-b63a-4a71-be40-860ae19834b7", "8aae665d-4845-4731-b962-15e5bdb69190"]}]}, {"BlockType": "LINE", "Confidence": 98.48675537109375, "Text": "Cell 1", "Geometry": {"BoundingBox": {"Width": 0.042510487139225006, "Height": 0.0135804433375597, "Left": 0.12771032750606537, "Top": 0.6181319952011108}, "Polygon": [{"X": 0.12771032750606537, "Y": 0.6181319952011108}, {"X": 0.17022082209587097, "Y": 0.6181319952011108}, {"X": 0.17022082209587097, "Y": 0.6317124366760254}, {"X": 0.12771032750606537, "Y": 0.6317124366760254}]}, "Id": "c1020ebf-bed2-4212-b24c-5972272cc8d8", "Relationships": [{"Type": "CHILD", "Ids": ["f5658bba-30d9-4320-bb4c-25a1df038473", "c302a70c-06ed-4bc4-9b53-177c22274189"]}]}, {"BlockType": "LINE", "Confidence": 98.98334503173828, "Text": "Cell 2", "Geometry": {"BoundingBox": {"Width": 0.043671898543834686, "Height": 0.012698537670075893, "Left": 0.285184383392334, "Top": 0.6184539794921875}, "Polygon": [{"X": 0.285184383392334, "Y": 0.6184539794921875}, {"X": 0.32885628938674927, "Y": 0.6184539794921875}, {"X": 0.32885628938674927, "Y": 0.6311525106430054}, {"X": 0.285184383392334, "Y": 0.6311525106430054}]}, "Id": "69c4ecc9-d075-4b7d-a846-f254d1137582", "Relationships": [{"Type": "CHILD", "Ids": ["7255c7cf-932f-4ee3-97de-c9f4878ca831", "cc7c7aa8-045f-4b33-b9db-f4a72a48f909"]}]}, {"BlockType": "LINE", "Confidence": 98.5614242553711, "Text": "Cell 4", "Geometry": {"BoundingBox": {"Width": 0.04424715414643288, "Height": 0.012623554095625877, "Left": 0.6017507314682007, "Top": 0.6185309290885925}, "Polygon": [{"X": 0.6017507314682007, "Y": 0.6185309290885925}, {"X": 0.6459978818893433, "Y": 0.6185309290885925}, {"X": 0.6459978818893433, "Y": 0.631154477596283}, {"X": 0.6017507314682007, "Y": 0.631154477596283}]}, "Id": "35a58188-40c3-4761-8ea7-f188b3732fd1", "Relationships": [{"Type": "CHILD", "Ids": ["8825bb47-6e06-4236-abbe-66d999d64fdd", "faf6fc1c-0e79-4b3b-aed3-b04937d83f22"]}]}, {"BlockType": "LINE", "Confidence": 99.19432067871094, "Text": "Cell 5", "Geometry": {"BoundingBox": {"Width": 0.043169427663087845, "Height": 0.01277545653283596, "Left": 0.7619276642799377, "Top": 0.6184312701225281}, "Polygon": [{"X": 0.7619276642799377, "Y": 0.6184312701225281}, {"X": 0.8050971031188965, "Y": 0.6184312701225281}, {"X": 0.8050971031188965, "Y": 0.6312066912651062}, {"X": 0.7619276642799377, "Y": 0.6312066912651062}]}, "Id": "b2f081a5-21c5-4763-b876-0313f11a7575", "Relationships": [{"Type": "CHILD", "Ids": ["10700bf6-b20c-4e16-8c55-822c47674698", "0cb5e8fd-62c0-46f2-9b5b-fc06e583bb96"]}]}, {"BlockType": "LINE", "Confidence": 98.72454071044922, "Text": "Cell 6", "Geometry": {"BoundingBox": {"Width": 0.04389132186770439, "Height": 0.012387072667479515, "Left": 0.12712755799293518, "Top": 0.6778441667556763}, "Polygon": [{"X": 0.12712755799293518, "Y": 0.6778441667556763}, {"X": 0.17101888358592987, "Y": 0.6778441667556763}, {"X": 0.17101888358592987, "Y": 0.6902312636375427}, {"X": 0.12712755799293518, "Y": 0.6902312636375427}]}, "Id": "96e1d611-3ffe-48ab-b1b9-c3d7086200e2", "Relationships": [{"Type": "CHILD", "Ids": ["5bef9078-00b2-4a7c-b23c-7f846953d664", "964d934e-caa2-4b2e-93c2-b42e0deb258f"]}]}, {"BlockType": "LINE", "Confidence": 98.60836791992188, "Text": "Cell 7", "Geometry": {"BoundingBox": {"Width": 0.0437912791967392, "Height": 0.012690716423094273, "Left": 0.28514277935028076, "Top": 0.6777071952819824}, "Polygon": [{"X": 0.28514277935028076, "Y": 0.6777071952819824}, {"X": 0.32893407344818115, "Y": 0.6777071952819824}, {"X": 0.32893407344818115, "Y": 0.6903979182243347}, {"X": 0.28514277935028076, "Y": 0.6903979182243347}]}, "Id": "12ff6f0b-7639-44ee-b197-d34cdb52e7aa", "Relationships": [{"Type": "CHILD", "Ids": ["35747715-068a-4488-b730-84cff9b3b4d9", "30ee0652-084a-44df-8d79-d7e90688e86d"]}]}, {"BlockType": "LINE", "Confidence": 98.65631866455078, "Text": "Cell 8", "Geometry": {"BoundingBox": {"Width": 0.04380646347999573, "Height": 0.01273917406797409, "Left": 0.4434744715690613, "Top": 0.6777850389480591}, "Polygon": [{"X": 0.4434744715690613, "Y": 0.6777850389480591}, {"X": 0.487280935049057, "Y": 0.6777850389480591}, {"X": 0.487280935049057, "Y": 0.690524160861969}, {"X": 0.4434744715690613, "Y": 0.690524160861969}]}, "Id": "c708dcc5-82ed-4c19-975e-d7ff860e6bdc", "Relationships": [{"Type": "CHILD", "Ids": ["a9bfcdf3-73c5-4aea-a1ce-7c290ee8c858", "c26c2104-e33c-494e-857d-2eb3d8f2e755"]}]}, {"BlockType": "LINE", "Confidence": 97.89591217041016, "Text": "Cell 9", "Geometry": {"BoundingBox": {"Width": 0.04354380443692207, "Height": 0.012920425273478031, "Left": 0.6017784476280212, "Top": 0.6776207089424133}, "Polygon": [{"X": 0.6017784476280212, "Y": 0.6776207089424133}, {"X": 0.6453222632408142, "Y": 0.6776207089424133}, {"X": 0.6453222632408142, "Y": 0.6905410885810852}, {"X": 0.6017784476280212, "Y": 0.6905410885810852}]}, "Id": "feedc078-68c0-42b8-a87c-72e3bd80ee78", "Relationships": [{"Type": "CHILD", "Ids": ["a276a37a-dff2-48f2-87fa-e12c6e063285", "e891dc91-d00e-4d68-85e2-4ca9745e5f72"]}]}, {"BlockType": "LINE", "Confidence": 99.05433654785156, "Text": "Cell 10", "Geometry": {"BoundingBox": {"Width": 0.053810954093933105, "Height": 0.01283883024007082, "Left": 0.7615987062454224, "Top": 0.6777843832969666}, "Polygon": [{"X": 0.7615987062454224, "Y": 0.6777843832969666}, {"X": 0.8154096603393555, "Y": 0.6777843832969666}, {"X": 0.8154096603393555, "Y": 0.6906231641769409}, {"X": 0.7615987062454224, "Y": 0.6906231641769409}]}, "Id": "1e3fb0bf-0f48-4e7a-9a3f-cc1b37480da6", "Relationships": [{"Type": "CHILD", "Ids": ["0f688657-d127-4e86-b8a2-c4fc4f86c332", "6ecfa4d0-39c2-4457-a840-8cae6e1ccbe8"]}]}, {"BlockType": "LINE", "Confidence": 98.55902862548828, "Text": "Cell 11", "Geometry": {"BoundingBox": {"Width": 0.05301501974463463, "Height": 0.012528574094176292, "Left": 0.12711328268051147, "Top": 0.7360577583312988}, "Polygon": [{"X": 0.12711328268051147, "Y": 0.7360577583312988}, {"X": 0.1801283061504364, "Y": 0.7360577583312988}, {"X": 0.1801283061504364, "Y": 0.7485863566398621}, {"X": 0.12711328268051147, "Y": 0.7485863566398621}]}, "Id": "0769c3cf-e76e-4bae-bea1-8e3b1943a107", "Relationships": [{"Type": "CHILD", "Ids": ["485c672f-e8e6-4591-b74d-b4f57f1772ba", "dceb932a-7425-408e-84e4-ac67f192b4f9"]}]}, {"BlockType": "LINE", "Confidence": 98.943115234375, "Text": "Cell 12", "Geometry": {"BoundingBox": {"Width": 0.05333862826228142, "Height": 0.01295428816229105, "Left": 0.2852138876914978, "Top": 0.7360171675682068}, "Polygon": [{"X": 0.2852138876914978, "Y": 0.7360171675682068}, {"X": 0.3385525047779083, "Y": 0.7360171675682068}, {"X": 0.3385525047779083, "Y": 0.7489714622497559}, {"X": 0.2852138876914978, "Y": 0.7489714622497559}]}, "Id": "4c9957c5-c3af-43a4-a1d4-f6c31ea4e32b", "Relationships": [{"Type": "CHILD", "Ids": ["81a1b5fa-b0b7-4440-8d3d-6ea705478f9d", "f5957cac-9a73-4597-8bad-7fc3e7a6a03a"]}]}, {"BlockType": "LINE", "Confidence": 98.80059051513672, "Text": "Cell 13", "Geometry": {"BoundingBox": {"Width": 0.053260792046785355, "Height": 0.012730708345770836, "Left": 0.44334903359413147, "Top": 0.7359303832054138}, "Polygon": [{"X": 0.44334903359413147, "Y": 0.7359303832054138}, {"X": 0.49660980701446533, "Y": 0.7359303832054138}, {"X": 0.49660980701446533, "Y": 0.7486611008644104}, {"X": 0.44334903359413147, "Y": 0.7486611008644104}]}, "Id": "876d2f64-150b-4aea-9b90-7cd9a3fd18d8", "Relationships": [{"Type": "CHILD", "Ids": ["56a86bc6-4456-47a8-abdc-c1d65c846df0", "3bc5e241-e142-4649-b319-07bca393276f"]}]}, {"BlockType": "LINE", "Confidence": 98.6966552734375, "Text": "Cell 14", "Geometry": {"BoundingBox": {"Width": 0.05358792096376419, "Height": 0.012640646658837795, "Left": 0.6018773913383484, "Top": 0.7359676957130432}, "Polygon": [{"X": 0.6018773913383484, "Y": 0.7359676957130432}, {"X": 0.655465304851532, "Y": 0.7359676957130432}, {"X": 0.655465304851532, "Y": 0.7486083507537842}, {"X": 0.6018773913383484, "Y": 0.7486083507537842}]}, "Id": "edb9f59e-e615-4157-8949-22b8d3b1602c", "Relationships": [{"Type": "CHILD", "Ids": ["2285d452-ec5e-45fe-998f-46cff551cd4c", "9ae5fa57-a73d-4dd4-b955-7dbe429bb777"]}]}, {"BlockType": "LINE", "Confidence": 99.1570053100586, "Text": "Cell 15", "Geometry": {"BoundingBox": {"Width": 0.05318611487746239, "Height": 0.012859954498708248, "Left": 0.7616134285926819, "Top": 0.7359786033630371}, "Polygon": [{"X": 0.7616134285926819, "Y": 0.7359786033630371}, {"X": 0.8147995471954346, "Y": 0.7359786033630371}, {"X": 0.8147995471954346, "Y": 0.7488385438919067}, {"X": 0.7616134285926819, "Y": 0.7488385438919067}]}, "Id": "a8830a85-1f77-4e7f-9219-921c2da3abe0", "Relationships": [{"Type": "CHILD", "Ids": ["c71075f8-4d19-4d4e-b880-fe85e6b89ffe", "74789197-ef54-4c3e-99ba-e71283b3a5fb"]}]}, {"BlockType": "LINE", "Confidence": 99.71761322021484, "Text": "Selection Element", "Geometry": {"BoundingBox": {"Width": 0.263517826795578, "Height": 0.02159660868346691, "Left": 0.1177731454372406, "Top": 0.8170759677886963}, "Polygon": [{"X": 0.1177731454372406, "Y": 0.8170759677886963}, {"X": 0.3812909722328186, "Y": 0.8170759677886963}, {"X": 0.3812909722328186, "Y": 0.8386725783348083}, {"X": 0.1177731454372406, "Y": 0.8386725783348083}]}, "Id": "e8589aa0-8cfd-4a4d-8e32-49a25ccf37c8", "Relationships": [{"Type": "CHILD", "Ids": ["bb7b7a98-2476-477b-9e18-665aefef175d", "cb9ba9aa-963d-44e3-bfa6-b9cfae5b4a95"]}]}, {"BlockType": "LINE", "Confidence": 99.4539794921875, "Text": "Selected Checkbox", "Geometry": {"BoundingBox": {"Width": 0.1507243663072586, "Height": 0.012731917202472687, "Left": 0.19152064621448517, "Top": 0.8825798630714417}, "Polygon": [{"X": 0.19152064621448517, "Y": 0.8825798630714417}, {"X": 0.3422450125217438, "Y": 0.8825798630714417}, {"X": 0.3422450125217438, "Y": 0.8953118324279785}, {"X": 0.19152064621448517, "Y": 0.8953118324279785}]}, "Id": "32ec6d5b-8c95-460e-a900-8c2998a194c1", "Relationships": [{"Type": "CHILD", "Ids": ["72840516-aeff-403e-8ca3-71264b7a685f", "59e49da0-58fd-43c0-902e-703db3d5fd2b"]}]}, {"BlockType": "LINE", "Confidence": 99.18217468261719, "Text": "Un-Selected Checkbox", "Geometry": {"BoundingBox": {"Width": 0.17911794781684875, "Height": 0.012731756083667278, "Left": 0.1922188550233841, "Top": 0.9288517832756042}, "Polygon": [{"X": 0.1922188550233841, "Y": 0.9288517832756042}, {"X": 0.37133678793907166, "Y": 0.9288517832756042}, {"X": 0.37133678793907166, "Y": 0.9415835738182068}, {"X": 0.1922188550233841, "Y": 0.9415835738182068}]}, "Id": "83c8002d-2a4d-4e64-a344-5139503120e9", "Relationships": [{"Type": "CHILD", "Ids": ["747fbc4b-3087-4b6d-b5f0-e06340c77051", "8eefa517-30a1-48b4-a4d1-7ae547631df4"]}]}, {"BlockType": "WORD", "Confidence": 99.58283233642578, "Text": "Textractor", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.462296724319458, "Height": 0.06334945559501648, "Left": 0.1625330150127411, "Top": 0.10538586229085922}, "Polygon": [{"X": 0.1625330150127411, "Y": 0.10538586229085922}, {"X": 0.6248297095298767, "Y": 0.10538586229085922}, {"X": 0.6248297095298767, "Y": 0.1687353104352951}, {"X": 0.1625330150127411, "Y": 0.1687353104352951}]}, "Id": "0b1cf325-b523-4b86-b01e-50ade0563079"}, {"BlockType": "WORD", "Confidence": 99.80604553222656, "Text": "Test", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.18449851870536804, "Height": 0.06476041674613953, "Left": 0.6522817015647888, "Top": 0.10424439609050751}, "Polygon": [{"X": 0.6522817015647888, "Y": 0.10424439609050751}, {"X": 0.8367802500724792, "Y": 0.10424439609050751}, {"X": 0.8367802500724792, "Y": 0.16900481283664703}, {"X": 0.6522817015647888, "Y": 0.16900481283664703}]}, "Id": "e6668699-4d9e-4f8b-a242-3b414a26d3fd"}, {"BlockType": "WORD", "Confidence": 99.75547790527344, "Text": "Document", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.4699675440788269, "Height": 0.06192830950021744, "Left": 0.2708991765975952, "Top": 0.22141507267951965}, "Polygon": [{"X": 0.2708991765975952, "Y": 0.22141507267951965}, {"X": 0.7408667206764221, "Y": 0.22141507267951965}, {"X": 0.7408667206764221, "Y": 0.2833433747291565}, {"X": 0.2708991765975952, "Y": 0.2833433747291565}]}, "Id": "b2ab72c0-0470-44bf-9fd6-9e84e2ae8eea"}, {"BlockType": "WORD", "Confidence": 99.0321273803711, "Text": "Page", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.10762534290552139, "Height": 0.0416223481297493, "Left": 0.11965136229991913, "Top": 0.3389776349067688}, "Polygon": [{"X": 0.11965136229991913, "Y": 0.3389776349067688}, {"X": 0.22727669775485992, "Y": 0.3389776349067688}, {"X": 0.22727669775485992, "Y": 0.3806000053882599}, {"X": 0.11965136229991913, "Y": 0.3806000053882599}]}, "Id": "dcde4029-984c-4795-9d95-2a42332046c2"}, {"BlockType": "WORD", "Confidence": 99.68162536621094, "Text": "(1)", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.059845417737960815, "Height": 0.04347170516848564, "Left": 0.24050432443618774, "Top": 0.33620190620422363}, "Polygon": [{"X": 0.24050432443618774, "Y": 0.33620190620422363}, {"X": 0.30034974217414856, "Y": 0.33620190620422363}, {"X": 0.30034974217414856, "Y": 0.37967362999916077}, {"X": 0.24050432443618774, "Y": 0.37967362999916077}]}, "Id": "fd202ecd-c9d1-4080-925f-38775e904101"}, {"BlockType": "WORD", "Confidence": 99.45826721191406, "Text": "Key", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.05174700915813446, "Height": 0.025642264634370804, "Left": 0.11902542412281036, "Top": 0.4124619960784912}, "Polygon": [{"X": 0.11902542412281036, "Y": 0.4124619960784912}, {"X": 0.17077243328094482, "Y": 0.4124619960784912}, {"X": 0.17077243328094482, "Y": 0.4381042718887329}, {"X": 0.11902542412281036, "Y": 0.4381042718887329}]}, "Id": "773c791e-dee0-43a6-be2c-f7c7ee5ab8b3"}, {"BlockType": "WORD", "Confidence": 99.68306732177734, "Text": "-", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.011502363719046116, "Height": 0.004489026963710785, "Left": 0.17738384008407593, "Top": 0.42237207293510437}, "Polygon": [{"X": 0.17738384008407593, "Y": 0.42237207293510437}, {"X": 0.18888619542121887, "Y": 0.42237207293510437}, {"X": 0.18888619542121887, "Y": 0.42686110734939575}, {"X": 0.17738384008407593, "Y": 0.42686110734939575}]}, "Id": "93dea5cc-f1ba-47d8-bf03-ddcbfdec16c2"}, {"BlockType": "WORD", "Confidence": 99.5123291015625, "Text": "Values", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.0982576310634613, "Height": 0.02232169359922409, "Left": 0.19638550281524658, "Top": 0.4110378324985504}, "Polygon": [{"X": 0.19638550281524658, "Y": 0.4110378324985504}, {"X": 0.2946431338787079, "Y": 0.4110378324985504}, {"X": 0.2946431338787079, "Y": 0.4333595335483551}, {"X": 0.19638550281524658, "Y": 0.4333595335483551}]}, "Id": "455c6f4d-e69d-4685-9296-6671c3013f41"}, {"BlockType": "WORD", "Confidence": 99.92672729492188, "Text": "Name", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.06867816299200058, "Height": 0.0178519356995821, "Left": 0.11843090504407883, "Top": 0.4683231711387634}, "Polygon": [{"X": 0.11843090504407883, "Y": 0.4683231711387634}, {"X": 0.1871090680360794, "Y": 0.4683231711387634}, {"X": 0.1871090680360794, "Y": 0.4861750900745392}, {"X": 0.11843090504407883, "Y": 0.4861750900745392}]}, "Id": "5adabdea-adcd-4f4c-9000-95bdb979a1d6"}, {"BlockType": "WORD", "Confidence": 99.99512481689453, "Text": "of", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.025785064324736595, "Height": 0.01896621100604534, "Left": 0.19438794255256653, "Top": 0.4672558903694153}, "Polygon": [{"X": 0.19438794255256653, "Y": 0.4672558903694153}, {"X": 0.22017300128936768, "Y": 0.4672558903694153}, {"X": 0.22017300128936768, "Y": 0.4862220883369446}, {"X": 0.19438794255256653, "Y": 0.4862220883369446}]}, "Id": "38788f7f-9475-4577-bd54-7bcd74393c59"}, {"BlockType": "WORD", "Confidence": 99.57273864746094, "Text": "package:", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.10313721746206284, "Height": 0.02287496067583561, "Left": 0.22502322494983673, "Top": 0.4674684703350067}, "Polygon": [{"X": 0.22502322494983673, "Y": 0.4674684703350067}, {"X": 0.32816043496131897, "Y": 0.4674684703350067}, {"X": 0.32816043496131897, "Y": 0.4903434216976166}, {"X": 0.22502322494983673, "Y": 0.4903434216976166}]}, "Id": "1fdff324-efd6-4038-81de-280812891135"}, {"BlockType": "WORD", "Confidence": 99.6280746459961, "Text": "Textractor", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.12362529337406158, "Height": 0.01760973036289215, "Left": 0.3352770507335663, "Top": 0.46857234835624695}, "Polygon": [{"X": 0.3352770507335663, "Y": 0.46857234835624695}, {"X": 0.4589023292064667, "Y": 0.46857234835624695}, {"X": 0.4589023292064667, "Y": 0.4861820638179779}, {"X": 0.3352770507335663, "Y": 0.4861820638179779}]}, "Id": "dd335722-cf8f-457c-8ebc-013a79baedd6"}, {"BlockType": "WORD", "Confidence": 99.77838134765625, "Text": "Date", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.05016987770795822, "Height": 0.016659211367368698, "Left": 0.11804135888814926, "Top": 0.4987580180168152}, "Polygon": [{"X": 0.11804135888814926, "Y": 0.4987580180168152}, {"X": 0.16821123659610748, "Y": 0.4987580180168152}, {"X": 0.16821123659610748, "Y": 0.515417218208313}, {"X": 0.11804135888814926, "Y": 0.515417218208313}]}, "Id": "9271e8b3-8ab5-4b9f-895f-4830af054ca8"}, {"BlockType": "WORD", "Confidence": 97.09461975097656, "Text": ":", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.00551046384498477, "Height": 0.012295318767428398, "Left": 0.17386683821678162, "Top": 0.502734363079071}, "Polygon": [{"X": 0.17386683821678162, "Y": 0.502734363079071}, {"X": 0.17937730252742767, "Y": 0.502734363079071}, {"X": 0.17937730252742767, "Y": 0.5150296688079834}, {"X": 0.17386683821678162, "Y": 0.5150296688079834}]}, "Id": "b9b96df2-7918-4cdc-bd0a-53ce974b0743"}, {"BlockType": "WORD", "Confidence": 99.55718231201172, "Text": "08/14/2022", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.12413974851369858, "Height": 0.019717741757631302, "Left": 0.1848236471414566, "Top": 0.49777790904045105}, "Polygon": [{"X": 0.1848236471414566, "Y": 0.49777790904045105}, {"X": 0.3089633882045746, "Y": 0.49777790904045105}, {"X": 0.3089633882045746, "Y": 0.5174956321716309}, {"X": 0.1848236471414566, "Y": 0.5174956321716309}]}, "Id": "5d235159-eed6-43e0-b017-788735b366d8"}, {"BlockType": "WORD", "Confidence": 99.86399841308594, "Text": "Table", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.07890205085277557, "Height": 0.02227928303182125, "Left": 0.11691198498010635, "Top": 0.564307451248169}, "Polygon": [{"X": 0.11691198498010635, "Y": 0.564307451248169}, {"X": 0.19581404328346252, "Y": 0.564307451248169}, {"X": 0.19581404328346252, "Y": 0.5865867733955383}, {"X": 0.11691198498010635, "Y": 0.5865867733955383}]}, "Id": "47cef666-b63a-4a71-be40-860ae19834b7"}, {"BlockType": "WORD", "Confidence": 99.81806182861328, "Text": "1", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.016153108328580856, "Height": 0.020929092541337013, "Left": 0.205570787191391, "Top": 0.5654847621917725}, "Polygon": [{"X": 0.205570787191391, "Y": 0.5654847621917725}, {"X": 0.22172389924526215, "Y": 0.5654847621917725}, {"X": 0.22172389924526215, "Y": 0.5864138603210449}, {"X": 0.205570787191391, "Y": 0.5864138603210449}]}, "Id": "8aae665d-4845-4731-b962-15e5bdb69190"}, {"BlockType": "WORD", "Confidence": 97.53768920898438, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.02881036140024662, "Height": 0.0135804433375597, "Left": 0.12771032750606537, "Top": 0.6181319952011108}, "Polygon": [{"X": 0.12771032750606537, "Y": 0.6181319952011108}, {"X": 0.15652069449424744, "Y": 0.6181319952011108}, {"X": 0.15652069449424744, "Y": 0.6317124366760254}, {"X": 0.12771032750606537, "Y": 0.6317124366760254}]}, "Id": "f5658bba-30d9-4320-bb4c-25a1df038473"}, {"BlockType": "WORD", "Confidence": 99.43582153320312, "Text": "1", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.009144669398665428, "Height": 0.011633284389972687, "Left": 0.1610761433839798, "Top": 0.6190181374549866}, "Polygon": [{"X": 0.1610761433839798, "Y": 0.6190181374549866}, {"X": 0.17022082209587097, "Y": 0.6190181374549866}, {"X": 0.17022082209587097, "Y": 0.6306514143943787}, {"X": 0.1610761433839798, "Y": 0.6306514143943787}]}, "Id": "c302a70c-06ed-4bc4-9b53-177c22274189"}, {"BlockType": "WORD", "Confidence": 98.71371459960938, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.029801635071635246, "Height": 0.012698537670075893, "Left": 0.285184383392334, "Top": 0.6184539794921875}, "Polygon": [{"X": 0.285184383392334, "Y": 0.6184539794921875}, {"X": 0.3149860203266144, "Y": 0.6184539794921875}, {"X": 0.3149860203266144, "Y": 0.6311525106430054}, {"X": 0.285184383392334, "Y": 0.6311525106430054}]}, "Id": "7255c7cf-932f-4ee3-97de-c9f4878ca831"}, {"BlockType": "WORD", "Confidence": 99.25298309326172, "Text": "2", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.010345458984375, "Height": 0.012169862166047096, "Left": 0.31851083040237427, "Top": 0.6186932921409607}, "Polygon": [{"X": 0.31851083040237427, "Y": 0.6186932921409607}, {"X": 0.32885628938674927, "Y": 0.6186932921409607}, {"X": 0.32885628938674927, "Y": 0.6308631896972656}, {"X": 0.31851083040237427, "Y": 0.6308631896972656}]}, "Id": "cc7c7aa8-045f-4b33-b9db-f4a72a48f909"}, {"BlockType": "WORD", "Confidence": 98.03536224365234, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.029655318707227707, "Height": 0.012623554095625877, "Left": 0.6017507314682007, "Top": 0.6185309290885925}, "Polygon": [{"X": 0.6017507314682007, "Y": 0.6185309290885925}, {"X": 0.6314060688018799, "Y": 0.6185309290885925}, {"X": 0.6314060688018799, "Y": 0.631154477596283}, {"X": 0.6017507314682007, "Y": 0.631154477596283}]}, "Id": "8825bb47-6e06-4236-abbe-66d999d64fdd"}, {"BlockType": "WORD", "Confidence": 99.08748626708984, "Text": "4", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.011248490773141384, "Height": 0.011445663869380951, "Left": 0.6347494125366211, "Top": 0.618995726108551}, "Polygon": [{"X": 0.6347494125366211, "Y": 0.618995726108551}, {"X": 0.6459978818893433, "Y": 0.618995726108551}, {"X": 0.6459978818893433, "Y": 0.630441427230835}, {"X": 0.6347494125366211, "Y": 0.630441427230835}]}, "Id": "faf6fc1c-0e79-4b3b-aed3-b04937d83f22"}, {"BlockType": "WORD", "Confidence": 98.73052978515625, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.029450906440615654, "Height": 0.01277545653283596, "Left": 0.7619276642799377, "Top": 0.6184312701225281}, "Polygon": [{"X": 0.7619276642799377, "Y": 0.6184312701225281}, {"X": 0.7913786172866821, "Y": 0.6184312701225281}, {"X": 0.7913786172866821, "Y": 0.6312066912651062}, {"X": 0.7619276642799377, "Y": 0.6312066912651062}]}, "Id": "10700bf6-b20c-4e16-8c55-822c47674698"}, {"BlockType": "WORD", "Confidence": 99.65811920166016, "Text": "5", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.01000399049371481, "Height": 0.012092862278223038, "Left": 0.7950931191444397, "Top": 0.6188316345214844}, "Polygon": [{"X": 0.7950931191444397, "Y": 0.6188316345214844}, {"X": 0.8050971031188965, "Y": 0.6188316345214844}, {"X": 0.8050971031188965, "Y": 0.6309244632720947}, {"X": 0.7950931191444397, "Y": 0.6309244632720947}]}, "Id": "0cb5e8fd-62c0-46f2-9b5b-fc06e583bb96"}, {"BlockType": "WORD", "Confidence": 97.86188507080078, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.028879188001155853, "Height": 0.012308783829212189, "Left": 0.12712755799293518, "Top": 0.6778441667556763}, "Polygon": [{"X": 0.12712755799293518, "Y": 0.6778441667556763}, {"X": 0.15600673854351044, "Y": 0.6778441667556763}, {"X": 0.15600673854351044, "Y": 0.6901530027389526}, {"X": 0.12712755799293518, "Y": 0.6901530027389526}]}, "Id": "5bef9078-00b2-4a7c-b23c-7f846953d664"}, {"BlockType": "WORD", "Confidence": 99.58719635009766, "Text": "6", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.010334332473576069, "Height": 0.011654086410999298, "Left": 0.16068454086780548, "Top": 0.678577184677124}, "Polygon": [{"X": 0.16068454086780548, "Y": 0.678577184677124}, {"X": 0.17101888358592987, "Y": 0.678577184677124}, {"X": 0.17101888358592987, "Y": 0.6902312636375427}, {"X": 0.16068454086780548, "Y": 0.6902312636375427}]}, "Id": "964d934e-caa2-4b2e-93c2-b42e0deb258f"}, {"BlockType": "WORD", "Confidence": 98.22537994384766, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.030321579426527023, "Height": 0.0126646738499403, "Left": 0.28514277935028076, "Top": 0.6777071952819824}, "Polygon": [{"X": 0.28514277935028076, "Y": 0.6777071952819824}, {"X": 0.3154643476009369, "Y": 0.6777071952819824}, {"X": 0.3154643476009369, "Y": 0.6903718709945679}, {"X": 0.28514277935028076, "Y": 0.6903718709945679}]}, "Id": "35747715-068a-4488-b730-84cff9b3b4d9"}, {"BlockType": "WORD", "Confidence": 98.9913558959961, "Text": "7", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.010197600349783897, "Height": 0.01210785936564207, "Left": 0.3187364637851715, "Top": 0.678290069103241}, "Polygon": [{"X": 0.3187364637851715, "Y": 0.678290069103241}, {"X": 0.32893407344818115, "Y": 0.678290069103241}, {"X": 0.32893407344818115, "Y": 0.6903979182243347}, {"X": 0.3187364637851715, "Y": 0.6903979182243347}]}, "Id": "30ee0652-084a-44df-8d79-d7e90688e86d"}, {"BlockType": "WORD", "Confidence": 98.12083435058594, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.029959699138998985, "Height": 0.01273917406797409, "Left": 0.4434744715690613, "Top": 0.6777850389480591}, "Polygon": [{"X": 0.4434744715690613, "Y": 0.6777850389480591}, {"X": 0.473434180021286, "Y": 0.6777850389480591}, {"X": 0.473434180021286, "Y": 0.690524160861969}, {"X": 0.4434744715690613, "Y": 0.690524160861969}]}, "Id": "a9bfcdf3-73c5-4aea-a1ce-7c290ee8c858"}, {"BlockType": "WORD", "Confidence": 99.19180297851562, "Text": "8", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.010379361920058727, "Height": 0.011977242305874825, "Left": 0.4769015908241272, "Top": 0.6784162521362305}, "Polygon": [{"X": 0.4769015908241272, "Y": 0.6784162521362305}, {"X": 0.487280935049057, "Y": 0.6784162521362305}, {"X": 0.487280935049057, "Y": 0.6903935074806213}, {"X": 0.4769015908241272, "Y": 0.6903935074806213}]}, "Id": "c26c2104-e33c-494e-857d-2eb3d8f2e755"}, {"BlockType": "WORD", "Confidence": 96.94509887695312, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.030068129301071167, "Height": 0.012920425273478031, "Left": 0.6017784476280212, "Top": 0.6776207089424133}, "Polygon": [{"X": 0.6017784476280212, "Y": 0.6776207089424133}, {"X": 0.6318466067314148, "Y": 0.6776207089424133}, {"X": 0.6318466067314148, "Y": 0.6905410885810852}, {"X": 0.6017784476280212, "Y": 0.6905410885810852}]}, "Id": "a276a37a-dff2-48f2-87fa-e12c6e063285"}, {"BlockType": "WORD", "Confidence": 98.84673309326172, "Text": "9", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.010458518750965595, "Height": 0.011845496483147144, "Left": 0.6348637342453003, "Top": 0.6785191297531128}, "Polygon": [{"X": 0.6348637342453003, "Y": 0.6785191297531128}, {"X": 0.6453222632408142, "Y": 0.6785191297531128}, {"X": 0.6453222632408142, "Y": 0.6903645992279053}, {"X": 0.6348637342453003, "Y": 0.6903645992279053}]}, "Id": "e891dc91-d00e-4d68-85e2-4ca9745e5f72"}, {"BlockType": "WORD", "Confidence": 98.36264038085938, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.029403511434793472, "Height": 0.012695634737610817, "Left": 0.7615987062454224, "Top": 0.6777843832969666}, "Polygon": [{"X": 0.7615987062454224, "Y": 0.6777843832969666}, {"X": 0.7910022139549255, "Y": 0.6777843832969666}, {"X": 0.7910022139549255, "Y": 0.6904799938201904}, {"X": 0.7615987062454224, "Y": 0.6904799938201904}]}, "Id": "0f688657-d127-4e86-b8a2-c4fc4f86c332"}, {"BlockType": "WORD", "Confidence": 99.74604034423828, "Text": "10", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.020036855712532997, "Height": 0.012323135510087013, "Left": 0.7953728437423706, "Top": 0.6783000826835632}, "Polygon": [{"X": 0.7953728437423706, "Y": 0.6783000826835632}, {"X": 0.8154096603393555, "Y": 0.6783000826835632}, {"X": 0.8154096603393555, "Y": 0.6906231641769409}, {"X": 0.7953728437423706, "Y": 0.6906231641769409}]}, "Id": "6ecfa4d0-39c2-4457-a840-8cae6e1ccbe8"}, {"BlockType": "WORD", "Confidence": 97.57048797607422, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.029744284227490425, "Height": 0.012400537729263306, "Left": 0.12711328268051147, "Top": 0.7360577583312988}, "Polygon": [{"X": 0.12711328268051147, "Y": 0.7360577583312988}, {"X": 0.15685756504535675, "Y": 0.7360577583312988}, {"X": 0.15685756504535675, "Y": 0.7484583258628845}, {"X": 0.12711328268051147, "Y": 0.7484583258628845}]}, "Id": "485c672f-e8e6-4591-b74d-b4f57f1772ba"}, {"BlockType": "WORD", "Confidence": 99.54756164550781, "Text": "11", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.019168399274349213, "Height": 0.012011912651360035, "Left": 0.1609598994255066, "Top": 0.7365744709968567}, "Polygon": [{"X": 0.1609598994255066, "Y": 0.7365744709968567}, {"X": 0.1801283061504364, "Y": 0.7365744709968567}, {"X": 0.1801283061504364, "Y": 0.7485863566398621}, {"X": 0.1609598994255066, "Y": 0.7485863566398621}]}, "Id": "dceb932a-7425-408e-84e4-ac67f192b4f9"}, {"BlockType": "WORD", "Confidence": 98.26697540283203, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.02954833395779133, "Height": 0.01248890534043312, "Left": 0.2852138876914978, "Top": 0.7360171675682068}, "Polygon": [{"X": 0.2852138876914978, "Y": 0.7360171675682068}, {"X": 0.3147622048854828, "Y": 0.7360171675682068}, {"X": 0.3147622048854828, "Y": 0.7485060691833496}, {"X": 0.2852138876914978, "Y": 0.7485060691833496}]}, "Id": "81a1b5fa-b0b7-4440-8d3d-6ea705478f9d"}, {"BlockType": "WORD", "Confidence": 99.61925506591797, "Text": "12", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.01955912820994854, "Height": 0.012924294918775558, "Left": 0.31899335980415344, "Top": 0.7360471487045288}, "Polygon": [{"X": 0.31899335980415344, "Y": 0.7360471487045288}, {"X": 0.3385525047779083, "Y": 0.7360471487045288}, {"X": 0.3385525047779083, "Y": 0.7489714622497559}, {"X": 0.31899335980415344, "Y": 0.7489714622497559}]}, "Id": "f5957cac-9a73-4597-8bad-7fc3e7a6a03a"}, {"BlockType": "WORD", "Confidence": 97.85566711425781, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.029667267575860023, "Height": 0.01260742824524641, "Left": 0.44334903359413147, "Top": 0.7359303832054138}, "Polygon": [{"X": 0.44334903359413147, "Y": 0.7359303832054138}, {"X": 0.47301629185676575, "Y": 0.7359303832054138}, {"X": 0.47301629185676575, "Y": 0.7485378384590149}, {"X": 0.44334903359413147, "Y": 0.7485378384590149}]}, "Id": "56a86bc6-4456-47a8-abdc-c1d65c846df0"}, {"BlockType": "WORD", "Confidence": 99.74552154541016, "Text": "13", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.019049739465117455, "Height": 0.012203322723507881, "Left": 0.4775600731372833, "Top": 0.7364577651023865}, "Polygon": [{"X": 0.4775600731372833, "Y": 0.7364577651023865}, {"X": 0.49660980701446533, "Y": 0.7364577651023865}, {"X": 0.49660980701446533, "Y": 0.7486611008644104}, {"X": 0.4775600731372833, "Y": 0.7486611008644104}]}, "Id": "3bc5e241-e142-4649-b319-07bca393276f"}, {"BlockType": "WORD", "Confidence": 97.64884185791016, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.02963610365986824, "Height": 0.012640646658837795, "Left": 0.6018773913383484, "Top": 0.7359676957130432}, "Polygon": [{"X": 0.6018773913383484, "Y": 0.7359676957130432}, {"X": 0.6315134763717651, "Y": 0.7359676957130432}, {"X": 0.6315134763717651, "Y": 0.7486083507537842}, {"X": 0.6018773913383484, "Y": 0.7486083507537842}]}, "Id": "2285d452-ec5e-45fe-998f-46cff551cd4c"}, {"BlockType": "WORD", "Confidence": 99.74446868896484, "Text": "14", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.019679507240653038, "Height": 0.011959020048379898, "Left": 0.6357858180999756, "Top": 0.7366405725479126}, "Polygon": [{"X": 0.6357858180999756, "Y": 0.7366405725479126}, {"X": 0.655465304851532, "Y": 0.7366405725479126}, {"X": 0.655465304851532, "Y": 0.7485995888710022}, {"X": 0.6357858180999756, "Y": 0.7485995888710022}]}, "Id": "9ae5fa57-a73d-4dd4-b955-7dbe429bb777"}, {"BlockType": "WORD", "Confidence": 98.42084503173828, "Text": "Cell", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.029446326196193695, "Height": 0.012561631388962269, "Left": 0.7616134285926819, "Top": 0.7359786033630371}, "Polygon": [{"X": 0.7616134285926819, "Y": 0.7359786033630371}, {"X": 0.7910597324371338, "Y": 0.7359786033630371}, {"X": 0.7910597324371338, "Y": 0.7485402226448059}, {"X": 0.7616134285926819, "Y": 0.7485402226448059}]}, "Id": "c71075f8-4d19-4d4e-b880-fe85e6b89ffe"}, {"BlockType": "WORD", "Confidence": 99.8931655883789, "Text": "15", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.01951681263744831, "Height": 0.012637583538889885, "Left": 0.7952827215194702, "Top": 0.7362009882926941}, "Polygon": [{"X": 0.7952827215194702, "Y": 0.7362009882926941}, {"X": 0.8147995471954346, "Y": 0.7362009882926941}, {"X": 0.8147995471954346, "Y": 0.7488385438919067}, {"X": 0.7952827215194702, "Y": 0.7488385438919067}]}, "Id": "74789197-ef54-4c3e-99ba-e71283b3a5fb"}, {"BlockType": "WORD", "Confidence": 99.82105255126953, "Text": "Selection", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.1316467672586441, "Height": 0.02155419811606407, "Left": 0.1177731454372406, "Top": 0.8171183466911316}, "Polygon": [{"X": 0.1177731454372406, "Y": 0.8171183466911316}, {"X": 0.2494199126958847, "Y": 0.8171183466911316}, {"X": 0.2494199126958847, "Y": 0.8386725783348083}, {"X": 0.1177731454372406, "Y": 0.8386725783348083}]}, "Id": "bb7b7a98-2476-477b-9e18-665aefef175d"}, {"BlockType": "WORD", "Confidence": 99.61417388916016, "Text": "Element", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.12110535055398941, "Height": 0.021324409171938896, "Left": 0.2601856291294098, "Top": 0.8170759677886963}, "Polygon": [{"X": 0.2601856291294098, "Y": 0.8170759677886963}, {"X": 0.3812909722328186, "Y": 0.8170759677886963}, {"X": 0.3812909722328186, "Y": 0.8384003639221191}, {"X": 0.2601856291294098, "Y": 0.8384003639221191}]}, "Id": "cb9ba9aa-963d-44e3-bfa6-b9cfae5b4a95"}, {"BlockType": "WORD", "Confidence": 99.92532348632812, "Text": "Selected", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.06861749291419983, "Height": 0.012586626224219799, "Left": 0.19152064621448517, "Top": 0.8826417922973633}, "Polygon": [{"X": 0.19152064621448517, "Y": 0.8826417922973633}, {"X": 0.2601381242275238, "Y": 0.8826417922973633}, {"X": 0.2601381242275238, "Y": 0.8952284455299377}, {"X": 0.19152064621448517, "Y": 0.8952284455299377}]}, "Id": "72840516-aeff-403e-8ca3-71264b7a685f"}, {"BlockType": "WORD", "Confidence": 98.98263549804688, "Text": "Checkbox", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.07755771279335022, "Height": 0.012731917202472687, "Left": 0.26468729972839355, "Top": 0.8825798630714417}, "Polygon": [{"X": 0.26468729972839355, "Y": 0.8825798630714417}, {"X": 0.3422450125217438, "Y": 0.8825798630714417}, {"X": 0.3422450125217438, "Y": 0.8953118324279785}, {"X": 0.26468729972839355, "Y": 0.8953118324279785}]}, "Id": "59e49da0-58fd-43c0-902e-703db3d5fd2b"}, {"BlockType": "WORD", "Confidence": 99.56431579589844, "Text": "Un-Selected", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.09667694568634033, "Height": 0.012559858150780201, "Left": 0.1922188550233841, "Top": 0.9289750456809998}, "Polygon": [{"X": 0.1922188550233841, "Y": 0.9289750456809998}, {"X": 0.28889578580856323, "Y": 0.9289750456809998}, {"X": 0.28889578580856323, "Y": 0.9415349364280701}, {"X": 0.1922188550233841, "Y": 0.9415349364280701}]}, "Id": "747fbc4b-3087-4b6d-b5f0-e06340c77051"}, {"BlockType": "WORD", "Confidence": 98.800048828125, "Text": "Checkbox", "TextType": "PRINTED", "Geometry": {"BoundingBox": {"Width": 0.07832921296358109, "Height": 0.012731756083667278, "Left": 0.29300758242607117, "Top": 0.9288517832756042}, "Polygon": [{"X": 0.29300758242607117, "Y": 0.9288517832756042}, {"X": 0.37133678793907166, "Y": 0.9288517832756042}, {"X": 0.37133678793907166, "Y": 0.9415835738182068}, {"X": 0.29300758242607117, "Y": 0.9415835738182068}]}, "Id": "8eefa517-30a1-48b4-a4d1-7ae547631df4"}, {"BlockType": "QUERY", "Id": "cea16e7e-6a41-4529-9290-452d6735b495", "Query": {"Text": "Lorem ipsum?"}}, {"BlockType": "QUERY", "Id": "4bbcd465-02fe-4745-bf44-a4822a09233d", "Query": {"Text": "The quick brown fox jumps over the lazy dog?"}}], "AnalyzeDocumentModelVersion": "1.0", "ResponseMetadata": {"RequestId": "26c6df3a-768f-4d7a-bae5-32108a8c880d", "HTTPStatusCode": 200, "HTTPHeaders": {"x-amzn-requestid": "26c6df3a-768f-4d7a-bae5-32108a8c880d", "content-type": "application/x-amz-json-1.1", "content-length": "39883", "date": "Tue, 27 Sep 2022 18:12:51 GMT"}, "RetryAttempts": 0}}
````

## 6. Configuration surface

| Entrada | Tipo | Default | Gate | Secreto | Impacto |
|---|---|---|---|---|---|
| AWS account/region | identidad + región | ninguno | acceso y servicio aprobados | no | servicio/costo |
| credential chain | referencia runtime | ninguna | least privilege; nunca archivo/CLI | sí | autenticación |
| document class | routing ID | ninguna | una clase `REQUIRED` y corpus | no | schema/evaluación |
| Textract operation/features | enum AWS | ninguna | perfil fijo; no input libre | no | resultado/costo |
| Queries/Adapter | lista + versión | vacío | evaluados por clase | no | precisión/schema |
| `CALL_TEXTRACT` | env test | ausente | prohibido en offline; explícito sólo live | no | red/costo |
| automatic business storage | boolean | false | permanece false hasta evaluación completa | no | integridad |

## 7. Dependency bill

| Paquete | Pin | Uso | Licencia | Evidencia |
|---|---|---|---|---|
| `amazon-textract-textractor` | 1.10.0 | parser/CLI oficial | Apache-2.0 | wheel 311.287 bytes/SHA-256 fijado |
| runtime transitivo | 14 pins | AWS caller/parser, modelos, formatos | Apache-2.0/MIT/BSD/PSF | hashes Windows/Linux; OSV 15/0 fechado |
| test-only | 11 pins adicionales | pytest + fixtures deterministas | licencias de cada wheel | separado; pytest 9.0.3; OSV 26/0 fechado |

## 8. Apply order

1. Materializar en destino vacío y ejecutar `test_contracts.py`.
2. Adquirir los wheels exactos mediante el core oficial y verificar todos los hashes.
3. Instalar runtime con `--require-hashes`; instalar tooling sólo en test.
4. Ejecutar la prueba oficial de Queries con `CALL_TEXTRACT` ausente.
5. Componer routing, seguridad, perfil AWS, strict-field evaluation y sólo después live/corpus.
6. Rollback: detener nuevas llamadas, conservar input/respuesta/receipts y volver al lock/adapter anterior evaluado.

## 9. Verification

```text
python -m unittest -v aws_textractor_official/test_contracts.py
expected: 4 tests, OK

python -m pytest -q aws_textractor_official/upstream/tests/test_queries.py
expected offline with test lock: 2 passed, 1 skipped
```

La suite determinista source-exact produjo Linux 69 PASS/16 skips y Windows 67 PASS/16 skips/2 FAIL de portabilidad. OSV 2.5.1 produjo runtime 15/0 y test 26/0 el 2026-08-28. Live, corpus, carga, costo, IAM y persistencia empresarial siguen condicionados.

## 10. Reconstruction evidence

- source: AWS tag `v1.10.0`, commit `8ea5f9a…`, archive 79.080.629 bytes/SHA-256 `4fa38999…c7219`;
- wheel: PyPI 311.287 bytes/SHA-256 `8524224f…3762c`;
- materialización: cinco archivos authored y siete AWS legal/test/fixture con hashes individuales;
- pruebas: official Queries 2/1; Linux deterministic core 69/16; Windows 67/16/2 condicionados;
- SCA: runtime 15/0, test 26/0, reportes hash-locked y fechados;
- límites: conditions upstream 156–158; cero modificación semántica del código/test oficial;
- fecha: 2026-08-28; agente: Codex.
