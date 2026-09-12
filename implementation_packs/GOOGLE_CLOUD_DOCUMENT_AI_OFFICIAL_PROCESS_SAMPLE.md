# Google Cloud Document AI Official Process Sample

## 1. Metadata

```yaml
pack_id: "GOOGLE-CLOUD-DOCUMENT-AI-OFFICIAL-PROCESS-SAMPLE"
pack_version: "0.2.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa byte-a-byte los samples oficiales Google process_document y Custom Document Extractor con schema_override por solicitud, sus tests oficiales, licencia Apache-2.0, fixtures invoice/packing-list y contratos offline reproducibles contra el SDK fijado."
stacks: ["PowerShell 7", "Python 3.10+", "google-cloud-documentai 3.15.0"]
compatible_with: ["OFFICIAL-DOCUMENT-SDK-ARTIFACT-CORE 0.3.x", "GOOGLE-DOCUMENT-AI-RUNTIME 0.1.x", "SECURE-LOCAL-FILE-INGESTION-GATE 0.1.x", "STRICT-DOCUMENT-FIELD-EVALUATION-GATE 0.1.x"]
incompatible_with: ["persistencia automática", "ADC o IDs embebidos", "fixtures no fijados", "tratar el output packing-list como ground truth", "instalar requirements históricos sin nuevo lock"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0"
upstream_sources: ["https://github.com/GoogleCloudPlatform/python-docs-samples/tree/841379c28828404c6d3944e3da9a8a0f46b4cd0e/documentai/snippets", "https://github.com/GoogleCloudPlatform/python-docs-samples/tree/dc0eecc2187791fea70f48201241288e25fc9620/documentai/snippets", "https://github.com/GoogleCloudPlatform/document-ai-samples/tree/001ba391ab4a2f40d001cc0387618cb3c3699523/web-app-pix2info-python"]
verified_at: "2026-08-28"
```

## 2. Applicability

Use este pack cuando el blueprint seleccione Google Cloud Document AI y necesite caminos oficiales standalone para enviar bytes locales a un processor/version o configurar campos diferentes por solicitud mediante un Foundation Model Custom Document Extractor, más fixtures oficiales de invoice y packing-list. Requiere el SDK exacto admitido por el artifact core; la ejecución live requiere proyecto GCP, ADC, processor, región, cuota/costo y acceso autorizados.

No use el JSON packing-list como ground truth: Google lo publica como output de ejemplo del Invoice Processor y sus campos focales incluyen confidencias de 0,30–0,56. El pack demuestra integración, identidad de fixtures y comportamiento oficial; schema del negocio, corpus, evaluación, revisión y persistencia siguen fuera del sample.

## 3. Architecture contract

Los cinco archivos bajo `upstream/google` son `VERBATIM`: licencia, sample process_document y test live del commit Google firmado `841379c...`, más sample handle_response y test live del commit Google firmado actual `dc0eecc...`. El segundo sample implementa `DocumentSchema` + `ProcessOptions.schema_override` por solicitud sin entrenamiento adicional del código cliente. El layout corto evita el límite de ruta de Windows sin cambiar los bytes oficiales. El lock enlaza commits firmados y cuatro artefactos: invoice del test, packing-list PNG, output JSON cacheado y licencia del segundo repositorio. El adquiridor exige propiedades exactas, commit-pinned raw URLs, approval hash-linked, aceptación Apache-2.0, cache/red explícita, staging, tamaño/SHA y destino inexistente.

Las regresiones offline ejecutan las funciones Google reales: una usa dobles sólo para el SDK y prueba endpoint regional, processor/version path, bytes, MIME, field mask y page selector; la otra captura la llamada sin red y prueba que el sample CDE exacto entrega un schema override con clase, base type, tres propiedades y ocurrencias al SDK oficial 3.15.0. El validador real comprueba PNG, JSON, texto, siete entidades, barcode CODE_39 y menciones exactas; su receipt siempre declara `automatic_storage_authorized=false`. Los tests upstream live se conservan sin modificar y sólo se ejecutan con autoridad GCP.

Presupuesto offline menor a 10 segundos; adquisición menor a 0,6 MB más overhead. Rollback elimina sólo el destino nuevo después de preservar receipts. Actualizar exige nuevos commits, firmas, hashes, diff, SDK lock y suites.

## 4. Exact file manifest

```text
CREATE google_document_ai_official_process_sample/README.md
CREATE google_document_ai_official_process_sample/upstream/google/LICENSE.txt
CREATE google_document_ai_official_process_sample/upstream/google/documentai/snippets/process_document_sample.py
CREATE google_document_ai_official_process_sample/upstream/google/documentai/snippets/process_document_sample_test.py
CREATE google_document_ai_official_process_sample/upstream/google/documentai/snippets/handle_response_sample.py
CREATE google_document_ai_official_process_sample/upstream/google/documentai/snippets/handle_response_sample_test.py
CREATE google_document_ai_official_process_sample/fixture-lock.json
CREATE google_document_ai_official_process_sample/fixture-approval.template.json
CREATE google_document_ai_official_process_sample/acquire_official_fixtures.ps1
CREATE google_document_ai_official_process_sample/validate_official_packing_list_output.py
CREATE google_document_ai_official_process_sample/test_official_process_document_sample.py
CREATE google_document_ai_official_process_sample/test_official_custom_extractor_sample.py
CREATE google_document_ai_official_process_sample/test_fixture_acquisition.ps1
```

## 5. Materialization blocks

### FILE: `google_document_ai_official_process_sample/README.md`
```yaml
block_id: "GOOGLE-DOCAI-OFFICIAL-PROCESS:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local integration guide constrained by exact Google upstream"
license: "LicenseRef-Workspace-Owner"
sha256: "c4cc3ce596a0b213a09d7391ec8494f9aeb3fe62a15b3536abd3451cff641785"
variables: []
secrets_allowed: false
```
````markdown
# Google Cloud Document AI official process sample

This directory contains the exact Google Cloud `process_document_sample.py`, its exact upstream test and Apache-2.0 license from signed commit `841379c28828404c6d3944e3da9a8a0f46b4cd0e` of `GoogleCloudPlatform/python-docs-samples`. It also contains the exact current `handle_response_sample.py` and its exact upstream test from signed commit `dc0eecc2187791fea70f48201241288e25fc9620` of the same Google repository. That second official sample constructs `DocumentSchema` properties and sends them through `ProcessOptions.schema_override`, so a Foundation Model Custom Document Extractor can receive a different extraction schema for each request.

`fixture-lock.json` also fixes the invoice used by that official test plus an official packing-list PNG and its cached Document AI JSON from signed commit `001ba391ab4a2f40d001cc0387618cb3c3699523` of `GoogleCloudPlatform/document-ai-samples`. The packing-list output contains confidence values below 0.60; it is evidence of provider behavior, not authority to persist business fields.

Use order:

1. materialize `OFFICIAL-DOCUMENT-SDK-ARTIFACT-CORE` and this pack;
2. run `test_official_process_document_sample.py`, `test_official_custom_extractor_sample.py` and `test_fixture_acquisition.ps1` offline;
3. copy the approval template, bind its `lock_sha256`, identify the approver/date and accept Apache-2.0;
4. acquire all fixtures from an approved cache or with explicit network authorization;
5. run `validate_official_packing_list_output.py --fixture-root <destination> --output <new-receipt.json>`;
6. run Google's exact live tests only with an authorized GCP project, processors, ADC, quota and cost; the custom extractor test uses Google's published Foundation Model processor/version path and invoice fixture;
7. route live results through project schema, ground truth, evaluation, review and reconciliation. Never store directly from this sample.

No credentials or project IDs are embedded. The exact upstream tests contain Google's published processor IDs and read `GOOGLE_CLOUD_PROJECT`; they remain live samples, not production identities or an accuracy SLA. The offline custom-extractor contract never calls Google: it hash-verifies and executes the unmodified official function with a capture double, proving that SDK 3.15.0 receives the official per-request schema override.
````

### FILE: `google_document_ai_official_process_sample/upstream/google/LICENSE.txt`
```yaml
block_id: "GOOGLE-DOCAI-OFFICIAL-PROCESS:google-license:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/GoogleCloudPlatform/python-docs-samples/841379c28828404c6d3944e3da9a8a0f46b4cd0e/LICENSE"
license: "Apache-2.0"
sha256: "c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4"
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

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "[]"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright [yyyy] [name of copyright owner]

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
````

### FILE: `google_document_ai_official_process_sample/upstream/google/documentai/snippets/process_document_sample.py`
```yaml
block_id: "GOOGLE-DOCAI-OFFICIAL-PROCESS:google-sample:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/GoogleCloudPlatform/python-docs-samples/841379c28828404c6d3944e3da9a8a0f46b4cd0e/documentai/snippets/process_document_sample.py"
license: "Apache-2.0"
sha256: "7e384dc2c38ebdcbc2c268a63bb9be9794576c1408034e84d53fd8ac5dc504ee"
variables: []
secrets_allowed: false
```
````python
# Copyright 2020 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# [START documentai_process_document]
from typing import Optional

from google.api_core.client_options import ClientOptions
from google.cloud import documentai  # type: ignore

# TODO(developer): Uncomment these variables before running the sample.
# project_id = "YOUR_PROJECT_ID"
# location = "YOUR_PROCESSOR_LOCATION" # Format is "us" or "eu"
# processor_id = "YOUR_PROCESSOR_ID" # Create processor before running sample
# file_path = "/path/to/local/pdf"
# mime_type = "application/pdf" # Refer to https://cloud.google.com/document-ai/docs/file-types for supported file types
# field_mask = "text,entities,pages.pageNumber"  # Optional. The fields to return in the Document object.
# processor_version_id = "YOUR_PROCESSOR_VERSION_ID" # Optional. Processor version to use


def process_document_sample(
    project_id: str,
    location: str,
    processor_id: str,
    file_path: str,
    mime_type: str,
    field_mask: Optional[str] = None,
    processor_version_id: Optional[str] = None,
) -> None:
    # You must set the `api_endpoint` if you use a location other than "us".
    opts = ClientOptions(api_endpoint=f"{location}-documentai.googleapis.com")

    client = documentai.DocumentProcessorServiceClient(client_options=opts)

    if processor_version_id:
        # The full resource name of the processor version, e.g.:
        # `projects/{project_id}/locations/{location}/processors/{processor_id}/processorVersions/{processor_version_id}`
        name = client.processor_version_path(
            project_id, location, processor_id, processor_version_id
        )
    else:
        # The full resource name of the processor, e.g.:
        # `projects/{project_id}/locations/{location}/processors/{processor_id}`
        name = client.processor_path(project_id, location, processor_id)

    # Read the file into memory
    with open(file_path, "rb") as image:
        image_content = image.read()

    # Load binary data
    raw_document = documentai.RawDocument(content=image_content, mime_type=mime_type)

    # For more information: https://cloud.google.com/document-ai/docs/reference/rest/v1/ProcessOptions
    # Optional: Additional configurations for processing.
    process_options = documentai.ProcessOptions(
        # Process only specific pages
        individual_page_selector=documentai.ProcessOptions.IndividualPageSelector(
            pages=[1]
        )
    )

    # Configure the process request
    request = documentai.ProcessRequest(
        name=name,
        raw_document=raw_document,
        field_mask=field_mask,
        process_options=process_options,
    )

    result = client.process_document(request=request)

    # For a full list of `Document` object attributes, reference this page:
    # https://cloud.google.com/document-ai/docs/reference/rest/v1/Document
    document = result.document

    # Read the text recognition output from the processor
    print("The document contains the following text:")
    print(document.text)


# [END documentai_process_document]
````

### FILE: `google_document_ai_official_process_sample/upstream/google/documentai/snippets/process_document_sample_test.py`
```yaml
block_id: "GOOGLE-DOCAI-OFFICIAL-PROCESS:google-live-test:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/GoogleCloudPlatform/python-docs-samples/841379c28828404c6d3944e3da9a8a0f46b4cd0e/documentai/snippets/process_document_sample_test.py"
license: "Apache-2.0"
sha256: "5b39bbd2d309efcc4791bd739fa38961639fa33d7a4e58b739050d3cfd16fcd7"
variables: []
secrets_allowed: false
```
````python
# # Copyright 2020 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# flake8: noqa

import os

from documentai.snippets import process_document_sample

location = "us"
project_id = os.environ["GOOGLE_CLOUD_PROJECT"]
processor_id = "90484cfdedb024f6"
processor_version_id = "stable"
file_path = "resources/invoice.pdf"
mime_type = "application/pdf"
field_mask = "text,pages.pageNumber"


def test_process_document(capsys):
    process_document_sample.process_document_sample(
        project_id=project_id,
        location=location,
        processor_id=processor_id,
        file_path=file_path,
        mime_type=mime_type,
        field_mask=field_mask,
    )
    out, _ = capsys.readouterr()

    assert "text:" in out
    assert "Invoice" in out


def test_process_document_processor_version(capsys):
    process_document_sample.process_document_sample(
        project_id=project_id,
        location=location,
        processor_id=processor_id,
        processor_version_id=processor_version_id,
        file_path=file_path,
        mime_type=mime_type,
        field_mask=field_mask,
    )
    out, _ = capsys.readouterr()

    assert "text:" in out
    assert "Invoice" in out
````

### FILE: `google_document_ai_official_process_sample/upstream/google/documentai/snippets/handle_response_sample.py`
```yaml
block_id: "GOOGLE-DOCAI-OFFICIAL-PROCESS:handle-response-sample:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/GoogleCloudPlatform/python-docs-samples/dc0eecc2187791fea70f48201241288e25fc9620/documentai/snippets/handle_response_sample.py"
license: "Apache-2.0"
sha256: "8fc1adaceba3ac8ad4d871f8003b2ce43fbadfd9cc7b80599cefcc6edc8bfb91"
variables: []
secrets_allowed: false
```
````python
# Copyright 2020 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# [START documentai_process_ocr_document]
# [START documentai_process_form_document]
# [START documentai_process_specialized_document]
# [START documentai_process_splitter_document]
# [START documentai_process_layout_document]
# [START documentai_process_custom_extractor_document]

from typing import Optional, Sequence

from google.api_core.client_options import ClientOptions
from google.cloud import documentai

# TODO(developer): Uncomment these variables before running the sample.
# project_id = "YOUR_PROJECT_ID"
# location = "YOUR_PROCESSOR_LOCATION" # Format is "us" or "eu"
# processor_id = "YOUR_PROCESSOR_ID" # Create processor before running sample
# processor_version = "rc" # Refer to https://cloud.google.com/document-ai/docs/manage-processor-versions for more information
# file_path = "/path/to/local/pdf"
# mime_type = "application/pdf" # Refer to https://cloud.google.com/document-ai/docs/file-types for supported file types


# [END documentai_process_ocr_document]
# [END documentai_process_form_document]
# [END documentai_process_specialized_document]
# [END documentai_process_splitter_document]
# [END documentai_process_layout_document]
# [END documentai_process_custom_extractor_document]


# [START documentai_process_ocr_document]
def process_document_ocr_sample(
    project_id: str,
    location: str,
    processor_id: str,
    processor_version: str,
    file_path: str,
    mime_type: str,
) -> None:
    # Optional: Additional configurations for Document OCR Processor.
    # For more information: https://cloud.google.com/document-ai/docs/enterprise-document-ocr
    process_options = documentai.ProcessOptions(
        ocr_config=documentai.OcrConfig(
            enable_native_pdf_parsing=True,
            enable_image_quality_scores=True,
            enable_symbol=True,
            # OCR Add Ons https://cloud.google.com/document-ai/docs/ocr-add-ons
            premium_features=documentai.OcrConfig.PremiumFeatures(
                compute_style_info=True,
                enable_math_ocr=False,  # Enable to use Math OCR Model
                enable_selection_mark_detection=True,
            ),
        )
    )
    # Online processing request to Document AI
    document = process_document(
        project_id,
        location,
        processor_id,
        processor_version,
        file_path,
        mime_type,
        process_options=process_options,
    )

    text = document.text
    print(f"Full document text: {text}\n")
    print(f"There are {len(document.pages)} page(s) in this document.\n")

    for page in document.pages:
        print(f"Page {page.page_number}:")
        print_page_dimensions(page.dimension)
        print_detected_languages(page.detected_languages)

        print_blocks(page.blocks, text)
        print_paragraphs(page.paragraphs, text)
        print_lines(page.lines, text)
        print_tokens(page.tokens, text)

        if page.symbols:
            print_symbols(page.symbols, text)

        if page.image_quality_scores:
            print_image_quality_scores(page.image_quality_scores)

        if page.visual_elements:
            print_visual_elements(page.visual_elements, text)


def print_page_dimensions(dimension: documentai.Document.Page.Dimension) -> None:
    print(f"    Width: {str(dimension.width)}")
    print(f"    Height: {str(dimension.height)}")


def print_detected_languages(
    detected_languages: Sequence[documentai.Document.Page.DetectedLanguage],
) -> None:
    print("    Detected languages:")
    for lang in detected_languages:
        print(f"        {lang.language_code} ({lang.confidence:.1%} confidence)")


def print_blocks(blocks: Sequence[documentai.Document.Page.Block], text: str) -> None:
    print(f"    {len(blocks)} blocks detected:")
    first_block_text = layout_to_text(blocks[0].layout, text)
    print(f"        First text block: {repr(first_block_text)}")
    last_block_text = layout_to_text(blocks[-1].layout, text)
    print(f"        Last text block: {repr(last_block_text)}")


def print_paragraphs(
    paragraphs: Sequence[documentai.Document.Page.Paragraph], text: str
) -> None:
    print(f"    {len(paragraphs)} paragraphs detected:")
    first_paragraph_text = layout_to_text(paragraphs[0].layout, text)
    print(f"        First paragraph text: {repr(first_paragraph_text)}")
    last_paragraph_text = layout_to_text(paragraphs[-1].layout, text)
    print(f"        Last paragraph text: {repr(last_paragraph_text)}")


def print_lines(lines: Sequence[documentai.Document.Page.Line], text: str) -> None:
    print(f"    {len(lines)} lines detected:")
    first_line_text = layout_to_text(lines[0].layout, text)
    print(f"        First line text: {repr(first_line_text)}")
    last_line_text = layout_to_text(lines[-1].layout, text)
    print(f"        Last line text: {repr(last_line_text)}")


def print_tokens(tokens: Sequence[documentai.Document.Page.Token], text: str) -> None:
    print(f"    {len(tokens)} tokens detected:")
    first_token_text = layout_to_text(tokens[0].layout, text)
    first_token_break_type = tokens[0].detected_break.type_.name
    print(f"        First token text: {repr(first_token_text)}")
    print(f"        First token break type: {repr(first_token_break_type)}")
    if tokens[0].style_info:
        print_style_info(tokens[0].style_info)

    last_token_text = layout_to_text(tokens[-1].layout, text)
    last_token_break_type = tokens[-1].detected_break.type_.name
    print(f"        Last token text: {repr(last_token_text)}")
    print(f"        Last token break type: {repr(last_token_break_type)}")
    if tokens[-1].style_info:
        print_style_info(tokens[-1].style_info)


def print_symbols(
    symbols: Sequence[documentai.Document.Page.Symbol], text: str
) -> None:
    print(f"    {len(symbols)} symbols detected:")
    first_symbol_text = layout_to_text(symbols[0].layout, text)
    print(f"        First symbol text: {repr(first_symbol_text)}")
    last_symbol_text = layout_to_text(symbols[-1].layout, text)
    print(f"        Last symbol text: {repr(last_symbol_text)}")


def print_image_quality_scores(
    image_quality_scores: documentai.Document.Page.ImageQualityScores,
) -> None:
    print(f"    Quality score: {image_quality_scores.quality_score:.1%}")
    print("    Detected defects:")

    for detected_defect in image_quality_scores.detected_defects:
        print(f"        {detected_defect.type_}: {detected_defect.confidence:.1%}")


def print_style_info(style_info: documentai.Document.Page.Token.StyleInfo) -> None:
    """
    Only supported in version `pretrained-ocr-v2.0-2023-06-02`
    """
    print(f"           Font Size: {style_info.font_size}pt")
    print(f"           Font Type: {style_info.font_type}")
    print(f"           Bold: {style_info.bold}")
    print(f"           Italic: {style_info.italic}")
    print(f"           Underlined: {style_info.underlined}")
    print(f"           Handwritten: {style_info.handwritten}")
    print(
        f"           Text Color (RGBa): {style_info.text_color.red}, {style_info.text_color.green}, {style_info.text_color.blue}, {style_info.text_color.alpha}"
    )


def print_visual_elements(
    visual_elements: Sequence[documentai.Document.Page.VisualElement], text: str
) -> None:
    """
    Only supported in version `pretrained-ocr-v2.0-2023-06-02`
    """
    checkboxes = [x for x in visual_elements if "checkbox" in x.type]
    math_symbols = [x for x in visual_elements if x.type == "math_formula"]

    if checkboxes:
        print(f"    {len(checkboxes)} checkboxes detected:")
        print(f"        First checkbox: {repr(checkboxes[0].type)}")
        print(f"        Last checkbox: {repr(checkboxes[-1].type)}")

    if math_symbols:
        print(f"    {len(math_symbols)} math symbols detected:")
        first_math_symbol_text = layout_to_text(math_symbols[0].layout, text)
        print(f"        First math symbol: {repr(first_math_symbol_text)}")


# [END documentai_process_ocr_document]
# [START documentai_process_form_document]
def process_document_form_sample(
    project_id: str,
    location: str,
    processor_id: str,
    processor_version: str,
    file_path: str,
    mime_type: str,
) -> documentai.Document:
    # Online processing request to Document AI
    document = process_document(
        project_id, location, processor_id, processor_version, file_path, mime_type
    )

    # Read the table and form fields output from the processor
    # The form processor also contains OCR data. For more information
    # on how to parse OCR data please see the OCR sample.

    text = document.text
    print(f"Full document text: {repr(text)}\n")
    print(f"There are {len(document.pages)} page(s) in this document.")

    # Read the form fields and tables output from the processor
    for page in document.pages:
        print(f"\n\n**** Page {page.page_number} ****")

        print(f"\nFound {len(page.tables)} table(s):")
        for table in page.tables:
            num_columns = len(table.header_rows[0].cells)
            num_rows = len(table.body_rows)
            print(f"Table with {num_columns} columns and {num_rows} rows:")

            # Print header rows
            print("Columns:")
            print_table_rows(table.header_rows, text)
            # Print body rows
            print("Table body data:")
            print_table_rows(table.body_rows, text)

        print(f"\nFound {len(page.form_fields)} form field(s):")
        for field in page.form_fields:
            name = layout_to_text(field.field_name, text)
            value = layout_to_text(field.field_value, text)
            print(f"    * {repr(name.strip())}: {repr(value.strip())}")

    # Supported in version `pretrained-form-parser-v2.0-2022-11-10` and later.
    # For more information: https://cloud.google.com/document-ai/docs/form-parser
    if document.entities:
        print(f"Found {len(document.entities)} generic entities:")
        for entity in document.entities:
            print_entity(entity)
            # Print Nested Entities
            for prop in entity.properties:
                print_entity(prop)

    return document


def print_table_rows(
    table_rows: Sequence[documentai.Document.Page.Table.TableRow], text: str
) -> None:
    for table_row in table_rows:
        row_text = ""
        for cell in table_row.cells:
            cell_text = layout_to_text(cell.layout, text)
            row_text += f"{repr(cell_text.strip())} | "
        print(row_text)


# [END documentai_process_form_document]


# [START documentai_process_specialized_document]
def process_document_entity_extraction_sample(
    project_id: str,
    location: str,
    processor_id: str,
    processor_version: str,
    file_path: str,
    mime_type: str,
) -> None:
    # Online processing request to Document AI
    document = process_document(
        project_id, location, processor_id, processor_version, file_path, mime_type
    )

    # Print extracted entities from entity extraction processor output.
    # For a complete list of processors see:
    # https://cloud.google.com/document-ai/docs/processors-list
    #
    # OCR and other data is also present in the processor's response.
    # Refer to the OCR samples for how to parse other data in the response.

    print(f"Found {len(document.entities)} entities:")
    for entity in document.entities:
        print_entity(entity)
        # Print Nested Entities (if any)
        for prop in entity.properties:
            print_entity(prop)


# [END documentai_process_specialized_document]


# [START documentai_process_custom_extractor_document]


def process_document_custom_extractor_sample(
    project_id: str,
    location: str,
    processor_id: str,
    processor_version: str,
    file_path: str,
    mime_type: str,
) -> None:
    # Entities to extract from Foundation Model CDE
    properties = [
        documentai.DocumentSchema.EntityType.Property(
            name="invoice_id",
            value_type="string",
            occurrence_type=documentai.DocumentSchema.EntityType.Property.OccurrenceType.REQUIRED_ONCE,
        ),
        documentai.DocumentSchema.EntityType.Property(
            name="notes",
            value_type="string",
            occurrence_type=documentai.DocumentSchema.EntityType.Property.OccurrenceType.OPTIONAL_MULTIPLE,
        ),
        documentai.DocumentSchema.EntityType.Property(
            name="terms",
            value_type="string",
            occurrence_type=documentai.DocumentSchema.EntityType.Property.OccurrenceType.OPTIONAL_MULTIPLE,
        ),
    ]
    # Optional: For Generative AI processors, request different fields than the
    # schema for a processor version
    process_options = documentai.ProcessOptions(
        schema_override=documentai.DocumentSchema(
            display_name="CDE Schema",
            description="Document Schema for the CDE Processor",
            entity_types=[
                documentai.DocumentSchema.EntityType(
                    name="custom_extraction_document_type",
                    base_types=["document"],
                    properties=properties,
                )
            ],
        )
    )

    # Online processing request to Document AI
    document = process_document(
        project_id,
        location,
        processor_id,
        processor_version,
        file_path,
        mime_type,
        process_options=process_options,
    )

    for entity in document.entities:
        print_entity(entity)
        # Print Nested Entities (if any)
        for prop in entity.properties:
            print_entity(prop)


# [START documentai_process_form_document]
# [START documentai_process_specialized_document]


def print_entity(entity: documentai.Document.Entity) -> None:
    # Fields detected. For a full list of fields for each processor see
    # the processor documentation:
    # https://cloud.google.com/document-ai/docs/processors-list
    key = entity.type_

    # Some other value formats in addition to text are available
    # e.g. dates: `entity.normalized_value.date_value.year`
    text_value = entity.text_anchor.content or entity.mention_text
    confidence = entity.confidence
    normalized_value = entity.normalized_value.text
    print(f"    * {repr(key)}: {repr(text_value)} ({confidence:.1%} confident)")

    if normalized_value:
        print(f"    * Normalized Value: {repr(normalized_value)}")


# [END documentai_process_form_document]
# [END documentai_process_specialized_document]
# [END documentai_process_custom_extractor_document]


# [START documentai_process_splitter_document]
def process_document_splitter_sample(
    project_id: str,
    location: str,
    processor_id: str,
    processor_version: str,
    file_path: str,
    mime_type: str,
) -> None:
    # Online processing request to Document AI
    document = process_document(
        project_id, location, processor_id, processor_version, file_path, mime_type
    )

    # Read the splitter output from a document splitter/classifier processor:
    # e.g. https://cloud.google.com/document-ai/docs/processors-list#processor_procurement-document-splitter
    # This processor only provides text for the document and information on how
    # to split the document on logical boundaries. To identify and extract text,
    # form elements, and entities please see other processors like the OCR, form,
    # and specalized processors.

    print(f"Found {len(document.entities)} subdocuments:")
    for entity in document.entities:
        conf_percent = f"{entity.confidence:.1%}"
        pages_range = page_refs_to_string(entity.page_anchor.page_refs)

        # Print subdocument type information, if available
        if entity.type_:
            print(
                f"{conf_percent} confident that {pages_range} a '{entity.type_}' subdocument."
            )
        else:
            print(f"{conf_percent} confident that {pages_range} a subdocument.")


def page_refs_to_string(
    page_refs: Sequence[documentai.Document.PageAnchor.PageRef],
) -> str:
    """Converts a page ref to a string describing the page or page range."""
    pages = [str(int(page_ref.page) + 1) for page_ref in page_refs]
    if len(pages) == 1:
        return f"page {pages[0]} is"
    else:
        return f"pages {', '.join(pages)} are"


# [END documentai_process_splitter_document]


# [START documentai_process_layout_document]
def process_document_layout_sample(
    project_id: str,
    location: str,
    processor_id: str,
    processor_version: str,
    file_path: str,
    mime_type: str,
) -> documentai.Document:
    process_options = documentai.ProcessOptions(
        layout_config=documentai.ProcessOptions.LayoutConfig(
            chunking_config=documentai.ProcessOptions.LayoutConfig.ChunkingConfig(
                chunk_size=1000,
                include_ancestor_headings=True,
            )
        )
    )

    document = process_document(
        project_id,
        location,
        processor_id,
        processor_version,
        file_path,
        mime_type,
        process_options=process_options,
    )

    print("Document Layout Blocks")
    for block in document.document_layout.blocks:
        print(block)

    print("Document Chunks")
    for chunk in document.chunked_document.chunks:
        print(chunk)

    # [END documentai_process_layout_document]
    return document


# [START documentai_process_ocr_document]
# [START documentai_process_form_document]
# [START documentai_process_specialized_document]
# [START documentai_process_splitter_document]
# [START documentai_process_layout_document]
# [START documentai_process_custom_extractor_document]


def process_document(
    project_id: str,
    location: str,
    processor_id: str,
    processor_version: str,
    file_path: str,
    mime_type: str,
    process_options: Optional[documentai.ProcessOptions] = None,
) -> documentai.Document:
    # You must set the `api_endpoint` if you use a location other than "us".
    client = documentai.DocumentProcessorServiceClient(
        client_options=ClientOptions(
            api_endpoint=f"{location}-documentai.googleapis.com"
        )
    )

    # The full resource name of the processor version, e.g.:
    # `projects/{project_id}/locations/{location}/processors/{processor_id}/processorVersions/{processor_version_id}`
    # You must create a processor before running this sample.
    name = client.processor_version_path(
        project_id, location, processor_id, processor_version
    )

    # Read the file into memory
    with open(file_path, "rb") as image:
        image_content = image.read()

    # Configure the process request
    request = documentai.ProcessRequest(
        name=name,
        raw_document=documentai.RawDocument(content=image_content, mime_type=mime_type),
        # Only supported for Document OCR processor
        process_options=process_options,
    )

    result = client.process_document(request=request)

    # For a full list of `Document` object attributes, reference this page:
    # https://cloud.google.com/document-ai/docs/reference/rest/v1/Document
    return result.document


# [END documentai_process_specialized_document]
# [END documentai_process_splitter_document]
# [END documentai_process_layout_document]
# [END documentai_process_custom_extractor_document]


def layout_to_text(layout: documentai.Document.Page.Layout, text: str) -> str:
    """
    Document AI identifies text in different parts of the document by their
    offsets in the entirety of the document"s text. This function converts
    offsets to a string.
    """
    # If a text segment spans several lines, it will
    # be stored in different text segments.
    return "".join(
        text[int(segment.start_index) : int(segment.end_index)]
        for segment in layout.text_anchor.text_segments
    )


# [END documentai_process_form_document]
# [END documentai_process_ocr_document]
````

### FILE: `google_document_ai_official_process_sample/upstream/google/documentai/snippets/handle_response_sample_test.py`
```yaml
block_id: "GOOGLE-DOCAI-OFFICIAL-PROCESS:handle-response-test:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/GoogleCloudPlatform/python-docs-samples/dc0eecc2187791fea70f48201241288e25fc9620/documentai/snippets/handle_response_sample_test.py"
license: "Apache-2.0"
sha256: "a60d4d9a53aa3c6383d60bb2cfc22489fde622d09ac994d68161ee771b1fea45"
variables: []
secrets_allowed: false
```
````python
# # Copyright 2020 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# flake8: noqa

import os

from documentai.snippets import handle_response_sample
from documentai.snippets import handle_response_sample_v1beta3


def test_process_document_ocr(capsys):
    location = "us"
    project_id = os.environ["GOOGLE_CLOUD_PROJECT"]
    processor_id = "52a38e080c1a7296"
    processor_version = "pretrained-ocr-v2.0-2023-06-02"
    file_path = "resources/handwritten_form.pdf"
    mime_type = "application/pdf"

    handle_response_sample.process_document_ocr_sample(
        project_id=project_id,
        location=location,
        processor_id=processor_id,
        processor_version=processor_version,
        file_path=file_path,
        mime_type=mime_type,
    )
    out, _ = capsys.readouterr()

    assert "Page 1" in out
    assert "en" in out
    assert "FakeDoc" in out

    # Document Quality
    assert "Quality score" in out

    # Font Detection
    assert "Font Size" in out
    assert "Handwritten" in out

    # Checkbox Detection
    file_path = "resources/checkbox.png"
    mime_type = "image/png"

    handle_response_sample.process_document_ocr_sample(
        project_id=project_id,
        location=location,
        processor_id=processor_id,
        processor_version=processor_version,
        file_path=file_path,
        mime_type=mime_type,
    )
    out, _ = capsys.readouterr()

    assert "unfilled_checkbox" in out
    assert "filled_checkbox" in out


def test_process_document_ocr_checkbox(capsys):
    location = "us"
    project_id = os.environ["GOOGLE_CLOUD_PROJECT"]
    processor_id = "52a38e080c1a7296"
    processor_version = "pretrained-ocr-v2.0-2023-06-02"
    file_path = "resources/checkbox.png"
    mime_type = "image/png"

    handle_response_sample.process_document_ocr_sample(
        project_id=project_id,
        location=location,
        processor_id=processor_id,
        processor_version=processor_version,
        file_path=file_path,
        mime_type=mime_type,
    )
    out, _ = capsys.readouterr()

    assert "unfilled_checkbox" in out
    assert "filled_checkbox" in out


def test_process_document_form():
    location = "us"
    project_id = os.environ["GOOGLE_CLOUD_PROJECT"]
    processor_id = "90484cfdedb024f6"
    processor_version = "pretrained-form-parser-v2.0-2022-11-10"
    file_path = "resources/invoice.pdf"
    mime_type = "application/pdf"

    document = handle_response_sample.process_document_form_sample(
        project_id=project_id,
        location=location,
        processor_id=processor_id,
        processor_version=processor_version,
        file_path=file_path,
        mime_type=mime_type,
    )

    assert len(document.pages) == 1
    assert len(document.pages[0].tables[0].header_rows[0].cells) == 4
    assert len(document.pages[0].tables[0].body_rows) == 6
    assert len(document.entities) > 0


def test_process_document_quality(capsys):
    location = "us"
    project_id = os.environ["GOOGLE_CLOUD_PROJECT"]
    processor_id = "52a38e080c1a7296"
    processor_version = "pretrained-ocr-v2.0-2023-06-02"
    poor_quality_file_path = "resources/document_quality_poor.pdf"
    mime_type = "application/pdf"

    handle_response_sample.process_document_ocr_sample(
        project_id=project_id,
        location=location,
        processor_id=processor_id,
        processor_version=processor_version,
        file_path=poor_quality_file_path,
        mime_type=mime_type,
    )
    out, _ = capsys.readouterr()

    expected_strings = [
        "Quality score",
        "defect_blurry",
        "defect_noisy",
    ]
    for expected_string in expected_strings:
        assert expected_string in out


def test_process_document_entity_extraction(capsys):
    location = "us"
    project_id = os.environ["GOOGLE_CLOUD_PROJECT"]
    processor_id = "feacd98c28866ede"
    processor_version = "stable"
    file_path = "resources/us_driver_license.pdf"
    mime_type = "application/pdf"

    handle_response_sample.process_document_entity_extraction_sample(
        project_id=project_id,
        location=location,
        processor_id=processor_id,
        processor_version=processor_version,
        file_path=file_path,
        mime_type=mime_type,
    )
    out, _ = capsys.readouterr()

    expected_strings = [
        "Document Id",
        "97551579",
    ]
    for expected_string in expected_strings:
        assert expected_string in out


def test_process_document_splitter(capsys):
    location = "us"
    project_id = os.environ["GOOGLE_CLOUD_PROJECT"]
    processor_id = "ed55eeb2b276066f"
    processor_version = "stable"
    file_path = "resources/multi_document.pdf"
    mime_type = "application/pdf"

    handle_response_sample.process_document_splitter_sample(
        project_id=project_id,
        location=location,
        processor_id=processor_id,
        processor_version=processor_version,
        file_path=file_path,
        mime_type=mime_type,
    )
    out, _ = capsys.readouterr()

    # Remove newlines and quotes from output for easier comparison
    out = out.replace(' "" ', " ").replace("\n", "")

    expected_strings = [
        "Found 8 subdocuments",
        "confident that pages 1, 2 are a subdocument",
        "confident that page 10 is a subdocument",
    ]
    for expected_string in expected_strings:
        assert expected_string in out


def test_process_document_summarizer(capsys):
    location = "us"
    project_id = os.environ["GOOGLE_CLOUD_PROJECT"]
    processor_id = "a2ab373924245a07"
    processor_version = "pretrained-foundation-model-v1.1-2023-09-12"
    file_path = "resources/superconductivity.pdf"
    mime_type = "application/pdf"

    handle_response_sample_v1beta3.process_document_summarizer_sample(
        project_id=project_id,
        location=location,
        processor_id=processor_id,
        processor_version=processor_version,
        file_path=file_path,
        mime_type=mime_type,
    )
    out, _ = capsys.readouterr()

    expected_strings = [
        "Superconductivity",
    ]
    for expected_string in expected_strings:
        assert expected_string in out


def test_process_document_layout():
    document = handle_response_sample.process_document_layout_sample(
        project_id=os.environ["GOOGLE_CLOUD_PROJECT"],
        location="us",
        processor_id="85b02a52f356f564",
        processor_version="pretrained",
        file_path="resources/superconductivity.pdf",
        mime_type="application/pdf",
    )

    assert document
    assert document.document_layout
    assert document.chunked_document


def test_process_document_custom_extractor(capsys):
    location = "us"
    project_id = os.environ["GOOGLE_CLOUD_PROJECT"]
    processor_id = "295e41049f27a2fa"
    processor_version = "pretrained-foundation-model-v1.0-2023-08-22"
    file_path = "resources/invoice.pdf"
    mime_type = "application/pdf"

    handle_response_sample.process_document_custom_extractor_sample(
        project_id=project_id,
        location=location,
        processor_id=processor_id,
        processor_version=processor_version,
        file_path=file_path,
        mime_type=mime_type,
    )
    out, _ = capsys.readouterr()

    expected_strings = ["invoice_id", "001"]
    for expected_string in expected_strings:
        assert expected_string in out
````

### FILE: `google_document_ai_official_process_sample/test_official_custom_extractor_sample.py`
```yaml
block_id: "GOOGLE-DOCAI-OFFICIAL-PROCESS:custom-extractor-contract:v1"
operation: CREATE
provenance: AUTHORED
source: "local verification harness constrained to exact Google sample and SDK"
license: "LicenseRef-Workspace-Owner"
sha256: "930bf606c8cc50fc089f3388485c7a68a0e0f4b953961cd19e192bc807a70914"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import hashlib
import importlib.util
from pathlib import Path

from google.cloud import documentai


ROOT = Path(__file__).resolve().parent
SAMPLE = ROOT / "upstream/google/documentai/snippets/handle_response_sample.py"
EXPECTED_SHA256 = "8fc1adaceba3ac8ad4d871f8003b2ce43fbadfd9cc7b80599cefcc6edc8bfb91"


def load_sample():
    assert hashlib.sha256(SAMPLE.read_bytes()).hexdigest() == EXPECTED_SHA256
    spec = importlib.util.spec_from_file_location("google_handle_response_sample", SAMPLE)
    assert spec is not None and spec.loader is not None
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def main() -> None:
    sample = load_sample()
    captured: dict[str, object] = {}

    def fake_process_document(
        project_id: str,
        location: str,
        processor_id: str,
        processor_version: str,
        file_path: str,
        mime_type: str,
        process_options: documentai.ProcessOptions | None = None,
    ) -> documentai.Document:
        captured.update(
            project_id=project_id,
            location=location,
            processor_id=processor_id,
            processor_version=processor_version,
            file_path=file_path,
            mime_type=mime_type,
            process_options=process_options,
        )
        return documentai.Document()

    sample.process_document = fake_process_document
    sample.process_document_custom_extractor_sample(
        project_id="offline-project",
        location="us",
        processor_id="offline-processor",
        processor_version="pretrained-foundation-model",
        file_path="offline.pdf",
        mime_type="application/pdf",
    )

    options = captured["process_options"]
    assert isinstance(options, documentai.ProcessOptions)
    schema = options.schema_override
    assert schema.display_name == "CDE Schema"
    assert len(schema.entity_types) == 1
    entity_type = schema.entity_types[0]
    assert entity_type.name == "custom_extraction_document_type"
    assert list(entity_type.base_types) == ["document"]
    assert [prop.name for prop in entity_type.properties] == ["invoice_id", "notes", "terms"]
    assert (
        entity_type.properties[0].occurrence_type
        == documentai.DocumentSchema.EntityType.Property.OccurrenceType.REQUIRED_ONCE
    )
    assert captured["mime_type"] == "application/pdf"
    print("GOOGLE_CUSTOM_EXTRACTOR_OFFICIAL_SAMPLE_OFFLINE_PASS properties=3 schema_override=1")


if __name__ == "__main__":
    main()
````

### FILE: `google_document_ai_official_process_sample/fixture-lock.json`
```yaml
block_id: "GOOGLE-DOCAI-OFFICIAL-PROCESS:fixture-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "registry derived from two signed Google commits and exact public artifacts"
license: "LicenseRef-Workspace-Owner"
sha256: "2f2e544b4eecadc255d4a9ea10dfaf7c6096edf465cb99fd557fe6965f22c474"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-google-document-ai-official-fixture-lock/v1",
  "verified_at": "2026-08-28",
  "sources": [
    {
      "id": "google-python-docs-samples",
      "repository": "GoogleCloudPlatform/python-docs-samples",
      "commit": "841379c28828404c6d3944e3da9a8a0f46b4cd0e",
      "commit_signature_verified": true,
      "license_expression": "Apache-2.0",
      "license_url": "https://raw.githubusercontent.com/GoogleCloudPlatform/python-docs-samples/841379c28828404c6d3944e3da9a8a0f46b4cd0e/LICENSE",
      "license_bytes": 11357,
      "license_sha256": "c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4"
    },
    {
      "id": "google-document-ai-samples",
      "repository": "GoogleCloudPlatform/document-ai-samples",
      "commit": "001ba391ab4a2f40d001cc0387618cb3c3699523",
      "commit_signature_verified": true,
      "license_expression": "Apache-2.0",
      "license_url": "https://raw.githubusercontent.com/GoogleCloudPlatform/document-ai-samples/001ba391ab4a2f40d001cc0387618cb3c3699523/LICENSE",
      "license_bytes": 11358,
      "license_sha256": "cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30"
    }
  ],
  "fixtures": [
    {
      "id": "google-python-docs-documentai-invoice-pdf",
      "source_id": "google-python-docs-samples",
      "url": "https://raw.githubusercontent.com/GoogleCloudPlatform/python-docs-samples/841379c28828404c6d3944e3da9a8a0f46b4cd0e/documentai/snippets/resources/invoice.pdf",
      "bytes": 58980,
      "sha256": "8dd79a9bfc49c8bc79180686b19a223a52177221460f5279e1c2edcada214ce5",
      "cache_path": "google-python-docs-samples/invoice.pdf",
      "destination_path": "google/invoice.pdf",
      "mime_type": "application/pdf",
      "classification": "OFFICIAL_LIVE_TEST_FIXTURE"
    },
    {
      "id": "google-document-ai-packing-list-png",
      "source_id": "google-document-ai-samples",
      "url": "https://raw.githubusercontent.com/GoogleCloudPlatform/document-ai-samples/001ba391ab4a2f40d001cc0387618cb3c3699523/web-app-pix2info-python/src/samples/4-INVOICE_PROCESSOR/b.%20Packing%20list%20with%20barcode.png",
      "bytes": 69350,
      "sha256": "b45f86a82194a6210ec3f2b149a9c486934b01f6087cf7e6f859f32c5f709e7b",
      "cache_path": "google-document-ai-samples/packing-list.png",
      "destination_path": "google_document_ai_samples/packing-list.png",
      "mime_type": "image/png",
      "classification": "OFFICIAL_SMOKE_FIXTURE"
    },
    {
      "id": "google-document-ai-packing-list-output-json",
      "source_id": "google-document-ai-samples",
      "url": "https://raw.githubusercontent.com/GoogleCloudPlatform/document-ai-samples/001ba391ab4a2f40d001cc0387618cb3c3699523/web-app-pix2info-python/src/samples/4-INVOICE_PROCESSOR/b.%20Packing%20list%20with%20barcode.json",
      "bytes": 361762,
      "sha256": "509cdc9e763f4c37253c617fba800452473c7a7e2bc5545ad6f6c5db956b91fb",
      "cache_path": "google-document-ai-samples/packing-list.document.json",
      "destination_path": "google_document_ai_samples/packing-list.document.json",
      "mime_type": "application/json",
      "classification": "OFFICIAL_CACHED_PROVIDER_OUTPUT"
    },
    {
      "id": "google-document-ai-samples-license",
      "source_id": "google-document-ai-samples",
      "url": "https://raw.githubusercontent.com/GoogleCloudPlatform/document-ai-samples/001ba391ab4a2f40d001cc0387618cb3c3699523/LICENSE",
      "bytes": 11358,
      "sha256": "cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30",
      "cache_path": "google-document-ai-samples/LICENSE.txt",
      "destination_path": "google_document_ai_samples/LICENSE.txt",
      "mime_type": "text/plain",
      "classification": "UPSTREAM_LICENSE"
    }
  ]
}
````

### FILE: `google_document_ai_official_process_sample/fixture-approval.template.json`
```yaml
block_id: "GOOGLE-DOCAI-OFFICIAL-PROCESS:approval-template:v1"
operation: CREATE
provenance: AUTHORED
source: "local approval template"
license: "LicenseRef-Workspace-Owner"
sha256: "1b905cde57f16f937c507204ab5da6822b3d9c9908d2e243ddcc9e1824a4c562"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-google-document-ai-official-fixture-approval/v1",
  "lock_sha256": "REPLACE_WITH_FIXTURE_LOCK_SHA256",
  "approved_fixture_ids": [
    "google-python-docs-documentai-invoice-pdf",
    "google-document-ai-packing-list-png",
    "google-document-ai-packing-list-output-json",
    "google-document-ai-samples-license"
  ],
  "approved_by": "REPLACE_WITH_IDENTITY",
  "approved_at": "YYYY-MM-DD",
  "accept_google_apache_2_0": false
}
````

### FILE: `google_document_ai_official_process_sample/acquire_official_fixtures.ps1`
```yaml
block_id: "GOOGLE-DOCAI-OFFICIAL-PROCESS:fixture-acquirer:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed acquisition adapter"
license: "LicenseRef-Workspace-Owner"
sha256: "f0d6de833e1f90a52db951bb7eaaf88e5028f658fb737a1bc9b3ba09f98b3c4c"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding()]
param(
  [string] $LockPath = (Join-Path $PSScriptRoot 'fixture-lock.json'),
  [string] $ApprovalPath = '',
  [string[]] $FixtureId = @(),
  [string] $DestinationRoot = '',
  [string] $CacheRoot = '',
  [switch] $AllowNetwork,
  [switch] $ValidateOnly
)

$ErrorActionPreference = 'Stop'
$utf8 = [Text.UTF8Encoding]::new($false)
function Fail([string]$Message) { throw "GOOGLE_DOCAI_OFFICIAL_FIXTURE_FAILED: $Message" }
function Assert-ExactProperties($Object, [string[]]$Expected, [string]$Label) {
  $actual = @($Object.psobject.Properties.Name | Sort-Object)
  $wanted = @($Expected | Sort-Object)
  if (($actual -join '|') -cne ($wanted -join '|')) { Fail "$Label properties mismatch" }
}
function Assert-SafeRelativePath([string]$Value, [string]$Label) {
  $normalized = $Value.Replace('\','/')
  $segments = $normalized.Split('/', [StringSplitOptions]::None)
  if ([string]::IsNullOrWhiteSpace($Value) -or [IO.Path]::IsPathRooted($Value) -or $segments -contains '' -or $segments -contains '.' -or $segments -contains '..') { Fail "$Label unsafe: $Value" }
  foreach ($segment in $segments) { if ($segment.EndsWith('.') -or $segment.EndsWith(' ')) { Fail "$Label non-canonical: $Value" } }
  $normalized
}
function Sha256([string]$Path) { (Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant() }

$resolvedLock = (Resolve-Path -LiteralPath $LockPath).Path
$lockBytes = [IO.File]::ReadAllBytes($resolvedLock)
$lockHash = [Convert]::ToHexString([Security.Cryptography.SHA256]::HashData($lockBytes)).ToLowerInvariant()
$lock = [Text.Encoding]::UTF8.GetString($lockBytes) | ConvertFrom-Json -Depth 30
Assert-ExactProperties $lock @('schema','verified_at','sources','fixtures') 'lock'
if ($lock.schema -cne 'elite-google-document-ai-official-fixture-lock/v1') { Fail 'unsupported lock schema' }
if ($lock.verified_at -cnotmatch '^\d{4}-\d{2}-\d{2}$') { Fail 'invalid verified_at' }
$sourceIds = @($lock.sources | ForEach-Object { $_.id })
if ($sourceIds.Count -eq 0 -or ($sourceIds | Sort-Object -Unique).Count -ne $sourceIds.Count) { Fail 'source IDs must be non-empty and unique' }
foreach ($source in $lock.sources) {
  Assert-ExactProperties $source @('id','repository','commit','commit_signature_verified','license_expression','license_url','license_bytes','license_sha256') "source $($source.id)"
  if ($source.repository -cnotmatch '^GoogleCloudPlatform/[A-Za-z0-9_.-]+$' -or $source.commit -cnotmatch '^[0-9a-f]{40}$' -or $source.commit_signature_verified -ne $true -or $source.license_expression -cne 'Apache-2.0') { Fail "invalid source authority: $($source.id)" }
  $licenseUri = [uri]$source.license_url
  if ($licenseUri.Scheme -cne 'https' -or $licenseUri.Host -cne 'raw.githubusercontent.com' -or -not $licenseUri.AbsolutePath.Contains("/$($source.commit)/", [StringComparison]::Ordinal)) { Fail "invalid source license URL: $($source.id)" }
  if ([long]$source.license_bytes -le 0 -or $source.license_sha256 -cnotmatch '^[0-9a-f]{64}$') { Fail "invalid source license integrity: $($source.id)" }
}
$fixtureIds = @($lock.fixtures | ForEach-Object { $_.id })
$destinations = @($lock.fixtures | ForEach-Object { $_.destination_path })
if ($fixtureIds.Count -eq 0 -or ($fixtureIds | Sort-Object -Unique).Count -ne $fixtureIds.Count) { Fail 'fixture IDs must be non-empty and unique' }
if (($destinations | Sort-Object -Unique).Count -ne $destinations.Count) { Fail 'fixture destinations must be unique' }
foreach ($fixture in $lock.fixtures) {
  Assert-ExactProperties $fixture @('id','source_id','url','bytes','sha256','cache_path','destination_path','mime_type','classification') "fixture $($fixture.id)"
  $source = @($lock.sources | Where-Object id -CEQ $fixture.source_id)
  if ($source.Count -ne 1) { Fail "fixture source missing: $($fixture.id)" }
  $uri = [uri]$fixture.url
  if ($uri.Scheme -cne 'https' -or $uri.Host -cne 'raw.githubusercontent.com' -or -not $uri.AbsolutePath.Contains("/$($source[0].commit)/", [StringComparison]::Ordinal)) { Fail "fixture URL is not source/commit pinned: $($fixture.id)" }
  $null = Assert-SafeRelativePath $fixture.cache_path "fixture cache $($fixture.id)"
  $null = Assert-SafeRelativePath $fixture.destination_path "fixture destination $($fixture.id)"
  if ([long]$fixture.bytes -le 0 -or $fixture.sha256 -cnotmatch '^[0-9a-f]{64}$' -or [string]::IsNullOrWhiteSpace($fixture.mime_type)) { Fail "invalid fixture integrity: $($fixture.id)" }
}
$selectedIds = if ($FixtureId.Count -eq 0) { @($fixtureIds) } else { @($FixtureId) }
if (($selectedIds | Sort-Object -Unique).Count -ne $selectedIds.Count) { Fail 'selected fixture IDs must be unique' }
$selected = @()
foreach ($id in $selectedIds) {
  $item = @($lock.fixtures | Where-Object id -CEQ $id)
  if ($item.Count -ne 1) { Fail "unknown fixture id: $id" }
  $selected += $item[0]
}
Write-Output "GOOGLE_DOCAI_OFFICIAL_FIXTURE_LOCK_VALID sources=$($sourceIds.Count) fixtures=$($fixtureIds.Count) selected=$($selected.Count)"
if ($ValidateOnly) { return }

if ([string]::IsNullOrWhiteSpace($ApprovalPath) -or [string]::IsNullOrWhiteSpace($DestinationRoot)) { Fail 'ApprovalPath and DestinationRoot are required outside ValidateOnly' }
$approval = Get-Content -LiteralPath (Resolve-Path -LiteralPath $ApprovalPath) -Raw | ConvertFrom-Json -Depth 10
Assert-ExactProperties $approval @('schema','lock_sha256','approved_fixture_ids','approved_by','approved_at','accept_google_apache_2_0') 'approval'
if ($approval.schema -cne 'elite-google-document-ai-official-fixture-approval/v1' -or $approval.lock_sha256 -cne $lockHash) { Fail 'approval schema or lock hash mismatch' }
if ([string]::IsNullOrWhiteSpace($approval.approved_by) -or $approval.approved_at -cnotmatch '^\d{4}-\d{2}-\d{2}$' -or $approval.accept_google_apache_2_0 -ne $true) { Fail 'approval identity/date/license acceptance missing' }
$approvedSet = (@($approval.approved_fixture_ids) | Sort-Object) -join '|'
$selectedSet = ($selectedIds | Sort-Object) -join '|'
if ($approvedSet -cne $selectedSet) { Fail 'approval fixture set mismatch' }

$destination = [IO.Path]::GetFullPath($DestinationRoot)
if (Test-Path -LiteralPath $destination) { Fail "destination already exists: $destination" }
$parent = Split-Path -Parent $destination
if (-not (Test-Path -LiteralPath $parent -PathType Container)) { Fail "destination parent missing: $parent" }
$staging = Join-Path $parent ('.google-docai-fixtures-' + [guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Path $staging | Out-Null
try {
  $receiptItems = @()
  foreach ($fixture in $selected) {
    $relativeCache = Assert-SafeRelativePath $fixture.cache_path "fixture cache $($fixture.id)"
    $candidate = if ($CacheRoot) { Join-Path ([IO.Path]::GetFullPath($CacheRoot)) $relativeCache } else { '' }
    $download = $candidate
    $networkUsed = $false
    if (-not $download -or -not (Test-Path -LiteralPath $download -PathType Leaf)) {
      if (-not $AllowNetwork) { Fail "fixture absent from cache and network not approved: $($fixture.id)" }
      $download = Join-Path $staging ('.download-' + [guid]::NewGuid().ToString('N'))
      Invoke-WebRequest -Uri $fixture.url -OutFile $download -Headers @{'User-Agent'='Elite-Engineering-Library-Fixture-Acquirer'}
      $networkUsed = $true
    }
    $item = Get-Item -LiteralPath $download
    $hash = Sha256 $download
    if ($item.Length -ne [long]$fixture.bytes -or $hash -cne $fixture.sha256) { Fail "fixture bytes/hash mismatch: $($fixture.id)" }
    $relativeDestination = Assert-SafeRelativePath $fixture.destination_path "fixture destination $($fixture.id)"
    $target = [IO.Path]::GetFullPath((Join-Path $staging $relativeDestination))
    if (-not $target.StartsWith([IO.Path]::GetFullPath($staging) + [IO.Path]::DirectorySeparatorChar, [StringComparison]::OrdinalIgnoreCase)) { Fail "fixture escaped staging: $($fixture.id)" }
    $targetParent = Split-Path -Parent $target
    if (-not (Test-Path -LiteralPath $targetParent)) { New-Item -ItemType Directory -Path $targetParent -Force | Out-Null }
    Copy-Item -LiteralPath $download -Destination $target
    $receiptItems += [ordered]@{ id=$fixture.id; source_id=$fixture.source_id; destination_path=$relativeDestination; bytes=[long]$fixture.bytes; sha256=$fixture.sha256; network_used=$networkUsed }
  }
  $receipt = [ordered]@{ schema='elite-google-document-ai-official-fixture-receipt/v1'; lock_sha256=$lockHash; acquired_at=[DateTimeOffset]::UtcNow.ToString('O'); items=$receiptItems }
  [IO.File]::WriteAllText((Join-Path $staging 'FIXTURE_ACQUISITION_RECEIPT.json'), (($receipt | ConvertTo-Json -Depth 8) + "`n"), $utf8)
  Move-Item -LiteralPath $staging -Destination $destination
  Write-Output "GOOGLE_DOCAI_OFFICIAL_FIXTURE_ACQUISITION_PASS fixtures=$($selected.Count)"
} finally {
  if (Test-Path -LiteralPath $staging) { Remove-Item -LiteralPath $staging -Recurse -Force }
}
````

### FILE: `google_document_ai_official_process_sample/validate_official_packing_list_output.py`
```yaml
block_id: "GOOGLE-DOCAI-OFFICIAL-PROCESS:packing-list-validator:v1"
operation: CREATE
provenance: AUTHORED
source: "local identity/schema validator over exact Google cached output"
license: "LicenseRef-Workspace-Owner"
sha256: "e67dab1a1d9f1a14948be69d13b01e0bf925aee1cd73809dafd0c81f03e411c1"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import argparse
import hashlib
import json
import os
from pathlib import Path
import tempfile

ROOT = Path(__file__).resolve().parent
EXPECTED_TEXT_SHA256 = "f6d1fd3d86adab02f4251d6d5910f88906535cda60906c42ab4e170fd51aee30"
EXPECTED_ENTITY_TYPES = [
    "invoice_type",
    "line_item",
    "purchase_order",
    "ship_to_address",
    "ship_to_name",
    "supplier_name",
    "supplier_website",
]


def sha256(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def fail(message: str) -> None:
    raise SystemExit(f"GOOGLE_DOCAI_PACKING_LIST_VALIDATION_FAILED: {message}")


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--fixture-root", required=True, type=Path)
    parser.add_argument("--output", required=True, type=Path)
    args = parser.parse_args()
    fixture_root = args.fixture_root.resolve(strict=True)
    output = args.output.resolve(strict=False)
    if output.exists() or not output.parent.is_dir():
        fail("output must be new and its parent must exist")

    lock = json.loads((ROOT / "fixture-lock.json").read_text(encoding="utf-8"))
    fixtures = {item["id"]: item for item in lock["fixtures"]}
    required = [
        "google-document-ai-packing-list-png",
        "google-document-ai-packing-list-output-json",
        "google-document-ai-samples-license",
    ]
    paths: dict[str, Path] = {}
    for fixture_id in required:
        item = fixtures.get(fixture_id)
        if item is None:
            fail(f"lock fixture missing: {fixture_id}")
        target = (fixture_root / item["destination_path"]).resolve(strict=True)
        if fixture_root != target and fixture_root not in target.parents:
            fail(f"fixture escaped root: {fixture_id}")
        if target.stat().st_size != item["bytes"] or sha256(target) != item["sha256"]:
            fail(f"fixture identity mismatch: {fixture_id}")
        paths[fixture_id] = target

    png = paths["google-document-ai-packing-list-png"]
    if not png.read_bytes().startswith(b"\x89PNG\r\n\x1a\n"):
        fail("packing-list fixture is not PNG")
    data = json.loads(paths["google-document-ai-packing-list-output-json"].read_text(encoding="utf-8"))
    if set(data) != {"uri", "mime_type", "text", "pages", "entities"}:
        fail("cached Document top-level schema mismatch")
    if data["mime_type"] != "image/png" or len(data["pages"]) != 1 or len(data["entities"]) != 7:
        fail("cached Document cardinality/mime mismatch")
    text_hash = hashlib.sha256(data["text"].encode("utf-8")).hexdigest()
    if text_hash != EXPECTED_TEXT_SHA256:
        fail("cached Document text mismatch")
    entity_types = sorted(entity["type_"] for entity in data["entities"])
    if entity_types != EXPECTED_ENTITY_TYPES:
        fail("cached Document entity types mismatch")
    barcodes = data["pages"][0].get("detected_barcodes", [])
    if len(barcodes) != 1 or barcodes[0].get("barcode") != {"format_": "CODE_39", "value_format": "TEXT", "raw_value": "E-58702865"}:
        fail("cached Document barcode mismatch")
    entities = {entity["type_"]: entity for entity in data["entities"]}
    if entities["purchase_order"].get("mention_text") != "86976433539-01" or entities["line_item"].get("mention_text") != "910-005694 MX MASTER 3 ADVANCED WIRELESS 1":
        fail("cached Document focal mentions mismatch")
    confidence = {name: entities[name]["confidence"] for name in ("purchase_order", "ship_to_name", "supplier_name")}
    receipt = {
        "schema": "elite-google-document-ai-official-output-validation/v1",
        "source_commit": lock["sources"][1]["commit"],
        "input_sha256": fixtures["google-document-ai-packing-list-png"]["sha256"],
        "output_sha256": fixtures["google-document-ai-packing-list-output-json"]["sha256"],
        "text_sha256": text_hash,
        "pages": 1,
        "entities": 7,
        "entity_types": entity_types,
        "barcode": "E-58702865",
        "selected_confidences": confidence,
        "automatic_storage_authorized": False,
        "classification": "OFFICIAL_CACHED_PROVIDER_OUTPUT_SMOKE_ONLY",
    }
    fd, temp_name = tempfile.mkstemp(prefix=f".{output.name}.", dir=output.parent)
    try:
        with os.fdopen(fd, "w", encoding="utf-8", newline="\n") as stream:
            json.dump(receipt, stream, indent=2, sort_keys=True)
            stream.write("\n")
        os.replace(temp_name, output)
    finally:
        if os.path.exists(temp_name):
            os.unlink(temp_name)
    print("GOOGLE_DOCAI_PACKING_LIST_OUTPUT_VALIDATION_PASS pages=1 entities=7 barcodes=1 storage_authorized=false")


if __name__ == "__main__":
    main()
````

### FILE: `google_document_ai_official_process_sample/test_official_process_document_sample.py`
```yaml
block_id: "GOOGLE-DOCAI-OFFICIAL-PROCESS:offline-sample-regression:v1"
operation: CREATE
provenance: AUTHORED
source: "local offline regression over verbatim Google sample"
license: "LicenseRef-Workspace-Owner"
sha256: "acda959d2bc2697cf84526b54349e471b31c9c7f256eb9ef3562e2dbb747db2b"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import contextlib
import hashlib
import importlib.util
import io
import json
from pathlib import Path
import sys
import tempfile
import types
import unittest

ROOT = Path(__file__).resolve().parent
SAMPLE = ROOT / "upstream/google/documentai/snippets/process_document_sample.py"
UPSTREAM_TEST = ROOT / "upstream/google/documentai/snippets/process_document_sample_test.py"
LICENSE = ROOT / "upstream/google/LICENSE.txt"


class FakeClientOptions:
    instances: list["FakeClientOptions"] = []

    def __init__(self, **kwargs):
        self.kwargs = kwargs
        self.__class__.instances.append(self)


class Value:
    def __init__(self, **kwargs):
        self.__dict__.update(kwargs)


class IndividualPageSelector(Value):
    pass


class ProcessOptions(Value):
    IndividualPageSelector = IndividualPageSelector


class FakeClient:
    instances: list["FakeClient"] = []

    def __init__(self, **kwargs):
        self.kwargs = kwargs
        self.calls: list[Value] = []
        self.__class__.instances.append(self)

    def processor_path(self, project, location, processor):
        return f"projects/{project}/locations/{location}/processors/{processor}"

    def processor_version_path(self, project, location, processor, version):
        return f"projects/{project}/locations/{location}/processors/{processor}/processorVersions/{version}"

    def process_document(self, *, request):
        self.calls.append(request)
        return Value(document=Value(text="Packing List\nE-58702865"))


def load_upstream_module():
    for name in list(sys.modules):
        if name == "google" or name.startswith("google."):
            del sys.modules[name]
    google = types.ModuleType("google")
    api_core = types.ModuleType("google.api_core")
    client_options = types.ModuleType("google.api_core.client_options")
    client_options.ClientOptions = FakeClientOptions
    cloud = types.ModuleType("google.cloud")
    documentai = types.ModuleType("google.cloud.documentai")
    documentai.DocumentProcessorServiceClient = FakeClient
    documentai.RawDocument = Value
    documentai.ProcessOptions = ProcessOptions
    documentai.ProcessRequest = Value
    google.api_core = api_core
    google.cloud = cloud
    api_core.client_options = client_options
    cloud.documentai = documentai
    sys.modules.update({
        "google": google,
        "google.api_core": api_core,
        "google.api_core.client_options": client_options,
        "google.cloud": cloud,
        "google.cloud.documentai": documentai,
    })
    spec = importlib.util.spec_from_file_location("official_google_process_document_sample", SAMPLE)
    module = importlib.util.module_from_spec(spec)
    assert spec and spec.loader
    spec.loader.exec_module(module)
    return module


class OfficialGoogleProcessDocumentSampleTests(unittest.TestCase):
    def setUp(self):
        FakeClient.instances.clear()
        FakeClientOptions.instances.clear()

    def test_verbatim_hashes_and_lock(self):
        self.assertEqual(hashlib.sha256(SAMPLE.read_bytes()).hexdigest(), "7e384dc2c38ebdcbc2c268a63bb9be9794576c1408034e84d53fd8ac5dc504ee")
        self.assertEqual(hashlib.sha256(UPSTREAM_TEST.read_bytes()).hexdigest(), "5b39bbd2d309efcc4791bd739fa38961639fa33d7a4e58b739050d3cfd16fcd7")
        self.assertEqual(hashlib.sha256(LICENSE.read_bytes()).hexdigest(), "c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4")
        lock = json.loads((ROOT / "fixture-lock.json").read_text(encoding="utf-8"))
        self.assertEqual([source["commit"] for source in lock["sources"]], ["841379c28828404c6d3944e3da9a8a0f46b4cd0e", "001ba391ab4a2f40d001cc0387618cb3c3699523"])
        self.assertTrue(all(source["commit_signature_verified"] for source in lock["sources"]))
        self.assertEqual(len(lock["fixtures"]), 4)

    def test_upstream_function_builds_exact_processor_request(self):
        module = load_upstream_module()
        with tempfile.TemporaryDirectory() as directory:
            fixture = Path(directory) / "packing-list.png"
            fixture.write_bytes(b"official-google-fixture-contract")
            output = io.StringIO()
            with contextlib.redirect_stdout(output):
                module.process_document_sample("project-1", "eu", "processor-1", str(fixture), "image/png", "text,entities,pages.pageNumber")
        client = FakeClient.instances[-1]
        request = client.calls[-1]
        self.assertEqual(FakeClientOptions.instances[-1].kwargs, {"api_endpoint": "eu-documentai.googleapis.com"})
        self.assertEqual(request.name, "projects/project-1/locations/eu/processors/processor-1")
        self.assertEqual(request.raw_document.content, b"official-google-fixture-contract")
        self.assertEqual(request.raw_document.mime_type, "image/png")
        self.assertEqual(request.field_mask, "text,entities,pages.pageNumber")
        self.assertEqual(request.process_options.individual_page_selector.pages, [1])
        self.assertIn("Packing List", output.getvalue())

    def test_upstream_function_honors_processor_version(self):
        module = load_upstream_module()
        with tempfile.TemporaryDirectory() as directory:
            fixture = Path(directory) / "invoice.pdf"
            fixture.write_bytes(b"official-google-invoice-contract")
            with contextlib.redirect_stdout(io.StringIO()):
                module.process_document_sample("project-2", "us", "processor-2", str(fixture), "application/pdf", processor_version_id="stable")
        request = FakeClient.instances[-1].calls[-1]
        self.assertEqual(request.name, "projects/project-2/locations/us/processors/processor-2/processorVersions/stable")


if __name__ == "__main__":
    unittest.main()
````

### FILE: `google_document_ai_official_process_sample/test_fixture_acquisition.ps1`
```yaml
block_id: "GOOGLE-DOCAI-OFFICIAL-PROCESS:acquisition-regression:v1"
operation: CREATE
provenance: AUTHORED
source: "local acquisition regression"
license: "LicenseRef-Workspace-Owner"
sha256: "2716f415e26a5dff7f44ef4998fd18535518e409fe89c237820a50c0fa246a31"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
$ErrorActionPreference = 'Stop'
$runner = Join-Path $PSScriptRoot 'acquire_official_fixtures.ps1'
$realLock = Join-Path $PSScriptRoot 'fixture-lock.json'
& $runner -LockPath $realLock -ValidateOnly
if (-not $?) { throw 'real fixture lock validation failed' }
$root = Join-Path ([IO.Path]::GetTempPath()) ('elite-google-docai-fixture-test-' + [guid]::NewGuid().ToString('N'))
$utf8 = [Text.UTF8Encoding]::new($false)
function Write-Json([string]$Path, $Value) { [IO.File]::WriteAllText($Path, (($Value | ConvertTo-Json -Depth 20) + "`n"), $utf8) }
function Expect-Failure([scriptblock]$Action, [string]$Pattern) {
  try { & $Action; throw "expected failure matching $Pattern" } catch { if ($_.Exception.Message -notmatch $Pattern) { throw } }
}
New-Item -ItemType Directory -Path $root | Out-Null
try {
  $cache = Join-Path $root 'cache'
  $cacheFile = Join-Path $cache 'google/fixture.bin'
  New-Item -ItemType Directory -Path (Split-Path -Parent $cacheFile) -Force | Out-Null
  [IO.File]::WriteAllBytes($cacheFile, [byte[]](1,2,3,4,5))
  $hash = (Get-FileHash -LiteralPath $cacheFile -Algorithm SHA256).Hash.ToLowerInvariant()
  $commit = '1111111111111111111111111111111111111111'
  $lock = [ordered]@{
    schema='elite-google-document-ai-official-fixture-lock/v1'; verified_at='2026-08-28'
    sources=@([ordered]@{id='google-test';repository='GoogleCloudPlatform/test';commit=$commit;commit_signature_verified=$true;license_expression='Apache-2.0';license_url="https://raw.githubusercontent.com/GoogleCloudPlatform/test/$commit/LICENSE";license_bytes=5;license_sha256=$hash})
    fixtures=@([ordered]@{id='fixture';source_id='google-test';url="https://raw.githubusercontent.com/GoogleCloudPlatform/test/$commit/fixture.bin";bytes=5;sha256=$hash;cache_path='google/fixture.bin';destination_path='fixtures/fixture.bin';mime_type='application/octet-stream';classification='OFFICIAL_SMOKE_FIXTURE'})
  }
  $lockPath = Join-Path $root 'lock.json'; Write-Json $lockPath $lock
  $lockHash = (Get-FileHash -LiteralPath $lockPath -Algorithm SHA256).Hash.ToLowerInvariant()
  $approval = [ordered]@{schema='elite-google-document-ai-official-fixture-approval/v1';lock_sha256=$lockHash;approved_fixture_ids=@('fixture');approved_by='test';approved_at='2026-08-28';accept_google_apache_2_0=$true}
  $approvalPath = Join-Path $root 'approval.json'; Write-Json $approvalPath $approval
  $destination = Join-Path $root 'destination'
  & $runner -LockPath $lockPath -ApprovalPath $approvalPath -FixtureId fixture -DestinationRoot $destination -CacheRoot $cache
  if (-not $?) { throw 'offline acquisition failed' }
  $target = Join-Path $destination 'fixtures/fixture.bin'
  if ((Get-FileHash -LiteralPath $target -Algorithm SHA256).Hash.ToLowerInvariant() -cne $hash -or -not (Test-Path -LiteralPath (Join-Path $destination 'FIXTURE_ACQUISITION_RECEIPT.json'))) { throw 'positive acquisition output mismatch' }
  Expect-Failure { & $runner -LockPath $lockPath -FixtureId unknown -ValidateOnly } 'unknown fixture id'
  Expect-Failure { & $runner -LockPath $lockPath -ApprovalPath $approvalPath -FixtureId fixture -DestinationRoot $destination -CacheRoot $cache } 'destination already exists'
  [IO.File]::WriteAllBytes($cacheFile, [byte[]](9,9,9,9,9))
  Expect-Failure { & $runner -LockPath $lockPath -ApprovalPath $approvalPath -FixtureId fixture -DestinationRoot (Join-Path $root 'tampered') -CacheRoot $cache } 'bytes/hash mismatch'
  [IO.File]::WriteAllBytes($cacheFile, [byte[]](1,2,3,4,5))
  $approval.accept_google_apache_2_0 = $false; $badApproval = Join-Path $root 'bad-approval.json'; Write-Json $badApproval $approval
  Expect-Failure { & $runner -LockPath $lockPath -ApprovalPath $badApproval -FixtureId fixture -DestinationRoot (Join-Path $root 'unlicensed') -CacheRoot $cache } 'license acceptance missing'
  Write-Output 'GOOGLE_DOCAI_OFFICIAL_FIXTURE_TEST_PASS positives=2 negatives=4'
} finally {
  if (Test-Path -LiteralPath $root) { Remove-Item -LiteralPath $root -Recurse -Force }
}
````

## 6. Configuration surface

- `GOOGLE_CLOUD_PROJECT`, ADC, location, processor ID/version: sólo para los tests/live oficiales; nunca se guardan en Markdown, lock o receipts.
- `fixture-approval.json`: identidad, fecha, aceptación Apache-2.0, hash del lock y conjunto exacto.
- `-FixtureId`, `-CacheRoot`, `-AllowNetwork`, `-ValidateOnly`: frontera de adquisición; red deshabilitada por defecto.
- `--fixture-root` y `--output`: validación del output oficial y receipt create-only.
- El field mask, processor/version y schema productivos se fijan en la decisión del proyecto; el sample CDE demuestra la forma oficial de cambiarlos por solicitud, pero no decide los campos del negocio.

## 7. Dependency bill

| Dependencia | Identidad | Licencia | Uso | Condición |
|---|---|---|---|---|
| Google python-docs-samples | commits firmados `841379c28828404c6d3944e3da9a8a0f46b4cd0e` y `dc0eecc2187791fea70f48201241288e25fc9620` | Apache-2.0 | process_document, Custom Extractor schema override y tests live | bytes/SHA exactos |
| Google document-ai-samples | commit firmado `001ba391ab4a2f40d001cc0387618cb3c3699523` | Apache-2.0 | packing-list input/output | acquire por approval/hash |
| google-cloud-documentai | `3.15.0` | Apache-2.0 | runtime seleccionado | artifact core + lock transitivo target |
| PowerShell | 7+ | MIT | adquisición/regresión | local/CI autorizado |
| Python | 3.10+ | PSF-2.0 | sample/validación | entorno aislado |

El `requirements.txt` histórico de los samples fija 3.0.1 y no se materializa como lock vigente; no se instala por inferencia.

## 8. Apply order

1. Complete el blueprint y seleccione Google Document AI.
2. Materialice/adquiera el artifact core 0.3.x y genere lock transitivo del target.
3. Materialice este pack y ejecute ambas regresiones offline, incluida la construcción exacta del schema override CDE.
4. Complete approval con hash/identidad/fecha/Apache-2.0.
5. Adquiera los cuatro artefactos por cache o red explícita y preserve receipt/licencias.
6. Valide el output packing-list oficial; observe las confidencias sin convertirlas en reglas universales.
7. Con GCP/ADC/processors autorizados ejecute los tests live exactos o los samples contra corpus propio; para una clase nueva use la ruta oficial CDE por solicitud o una versión entrenada/evaluada.
8. Evalúe por campo/clase y sólo persista mediante evidence/review/reconciliation gates.

## 9. Verification

```powershell
python .\google_document_ai_official_process_sample\test_official_process_document_sample.py
python .\google_document_ai_official_process_sample\test_official_custom_extractor_sample.py
pwsh -NoProfile -File .\google_document_ai_official_process_sample\test_fixture_acquisition.ps1
pwsh -NoProfile -File .\google_document_ai_official_process_sample\acquire_official_fixtures.ps1 -LockPath .\google_document_ai_official_process_sample\fixture-lock.json -ApprovalPath <approval> -DestinationRoot <new-directory> -CacheRoot <approved-cache>
python .\google_document_ai_official_process_sample\validate_official_packing_list_output.py --fixture-root <new-directory> --output <new-receipt.json>
```

PASS offline exige tres tests sobre process_document, el contrato CDE `properties=3 schema_override=1`, SHA exacto de ambos nuevos archivos Google, lock real, 2 positivos/4 negativos de adquisición y compilación del validador. PASS de fixtures exige cuatro hashes, páginas=1, entidades=7, barcode=1 y storage=false. PASS live requiere además cuenta/processor/costo y evidencia separada.

## 10. Reconstruction evidence

- python-docs-samples commit firmado, sample 3.596 bytes SHA `7e384dc2...504ee`, test 1.697 bytes SHA `5b39bbd2...6fcd7`, licencia 11.357 bytes SHA `c71d239d...d0ab4`.
- python-docs-samples commit firmado actual `dc0eecc2187791fea70f48201241288e25fc9620`: `handle_response_sample.py` 20.533 bytes SHA `8fc1adac...bfb91` y test oficial 7.859 bytes SHA `a60d4d9a...fea45`; ambos compilan y la función CDE exacta construye `schema_override` con el SDK oficial 3.15.0 en la regresión offline.
- document-ai-samples commit firmado, packing-list PNG 69.350 bytes SHA `b45f86a8...09e7b`, JSON 361.762 bytes SHA `509cdc9e...b91fb`, licencia SHA `cfc7749b...23d30`.
- Recorrido real offline adquirido: 4 fixtures; JSON oficial páginas=1, entidades=7, barcode `E-58702865`; confidencias focales purchase order 0,537, ship-to 0,556 y supplier 0,302; storage false.
- No se realizó llamada live, no se atribuye exactitud productiva y ninguna confianza se transforma en autorización.
