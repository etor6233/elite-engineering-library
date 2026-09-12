# Google Document AI Official Lifecycle

## 1. Metadata

```yaml
pack_id: "GOOGLE-DOCUMENT-AI-OFFICIAL-LIFECYCLE"
pack_version: "0.2.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa el ciclo documental oficial Google para crear processors, inicializar/copiar datasets y schema, importar train/test/unassigned, entrenar, evaluar, desplegar, fijar default, batch-procesar uno o muchos objetos GCS, interpretar respuestas OCR/form/table/entity/split/layout/custom y undeploy como rollback, con source/tests exactos, hashes y gates offline."
stacks: ["PowerShell 7", "Python 3.10+", "google-cloud-documentai 3.15.0", "google-cloud-storage 3.13.1"]
compatible_with: ["OFFICIAL-DOCUMENT-SDK-ARTIFACT-CORE 0.3.x", "GOOGLE-CLOUD-DOCUMENT-AI-OFFICIAL-PROCESS-SAMPLE 0.2.x", "GOOGLE-DOCUMENT-AI-RUNTIME 0.1.x", "STRICT-DOCUMENT-FIELD-EVALUATION-GATE 0.1.x"]
incompatible_with: ["persistencia automática", "placeholders sin configurar", "ejecución live sin GCP/ADC/IAM/cuota/costo", "importar el script documental como módulo sin retirar su main configurado", "confundir tests mock con precisión de negocio"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0"
upstream_sources: ["https://github.com/GoogleCloudPlatform/python-docs-samples/tree/dc0eecc2187791fea70f48201241288e25fc9620/documentai/snippets", "https://docs.cloud.google.com/document-ai/docs/copy-processor-versions?hl=en", "https://developers.google.com/terms/site-policies"]
verified_at: "2026-08-28"
```

## 2. Applicability

Use este pack cuando Google Document AI haya sido seleccionado y el agente necesite código oficial inmediato para todo el lifecycle alrededor de la extracción: processor, dataset/schema, importación, training, evaluation, deploy, default, migración y rollback por undeploy. La clase documental y sus campos siguen siendo configurables; el lifecycle no se limita a invoice, packing list, BOL ni una taxonomía cerrada.

No lo ejecute live hasta completar proyecto, región, processor type/ID, dataset, buckets, IAM, ADC, CMEK/retención, cuota y costo. No convierte la salida del proveedor en dato empresarial confirmado: la persistencia continúa bloqueada por ground truth, evaluación estricta, revisión de ambigüedad y reconciliación.

## 3. Architecture contract

Los dieciséis archivos bajo `upstream/google/documentai/snippets` son `VERBATIM` desde `python-docs-samples` commit Google firmado `dc0eecc...`: create/train/evaluate/deploy/default/undeploy más batch/response y sus ocho módulos de test. El layout elimina sólo el segmento redundante del nombre de repositorio para conservar portabilidad Windows; source-lock retiene la ruta Git exacta. El bloque de 432 líneas bajo `cloud-docs` es código Google de la página oficial actual, extraído del único nodo DOM que contiene `def import_documents`; quitar markup y decodificar HTML no alteró lógica, pero se clasifica `ADAPTED` para no fingir identidad de archivo Git. Google licencia sus code samples bajo Apache-2.0.

El verificador local es sólo arnés `AUTHORED`: fija bytes/hashes, compila 17 Python, prueba por AST lifecycle/batch/response, ejercita helpers puros de respuesta y ejecuta la función Google set-default con un capture client sin red. `--static-only` conserva hashes, compilación y AST cuando los SDK aún no están instalados. Cinco tests Google son offline-safe y pasan con SDK Document AI 3.15.0; set-default, batch y response son live y quedan preservados, no adulterados. El script largo de documentación posee placeholders y llamada `main` superior: se configura una copia de proyecto antes de ejecutarlo, nunca se importa como módulo libre de efectos.

## 4. Exact file manifest

```text
CREATE google_document_ai_official_lifecycle/README.md
CREATE google_document_ai_official_lifecycle/source-lock.json
CREATE google_document_ai_official_lifecycle/upstream/google/cloud-docs/copy_processor_versions_and_datasets_official_docs.py
CREATE google_document_ai_official_lifecycle/upstream/google/documentai/snippets/create_processor_sample_test.py
CREATE google_document_ai_official_lifecycle/upstream/google/documentai/snippets/create_processor_sample.py
CREATE google_document_ai_official_lifecycle/upstream/google/documentai/snippets/deploy_processor_version_sample_test.py
CREATE google_document_ai_official_lifecycle/upstream/google/documentai/snippets/deploy_processor_version_sample.py
CREATE google_document_ai_official_lifecycle/upstream/google/documentai/snippets/evaluate_processor_version_sample_test.py
CREATE google_document_ai_official_lifecycle/upstream/google/documentai/snippets/evaluate_processor_version_sample.py
CREATE google_document_ai_official_lifecycle/upstream/google/documentai/snippets/batch_process_documents_sample.py
CREATE google_document_ai_official_lifecycle/upstream/google/documentai/snippets/batch_process_documents_sample_test.py
CREATE google_document_ai_official_lifecycle/upstream/google/documentai/snippets/handle_response_sample.py
CREATE google_document_ai_official_lifecycle/upstream/google/documentai/snippets/handle_response_sample_test.py
CREATE google_document_ai_official_lifecycle/upstream/google/documentai/snippets/set_default_processor_version_sample_test.py
CREATE google_document_ai_official_lifecycle/upstream/google/documentai/snippets/set_default_processor_version_sample.py
CREATE google_document_ai_official_lifecycle/upstream/google/documentai/snippets/train_processor_version_sample_test.py
CREATE google_document_ai_official_lifecycle/upstream/google/documentai/snippets/train_processor_version_sample.py
CREATE google_document_ai_official_lifecycle/upstream/google/documentai/snippets/undeploy_processor_version_sample_test.py
CREATE google_document_ai_official_lifecycle/upstream/google/documentai/snippets/undeploy_processor_version_sample.py
CREATE google_document_ai_official_lifecycle/upstream/google/LICENSE.txt
CREATE google_document_ai_official_lifecycle/verify_official_lifecycle.py
```

## 5. Materialization blocks
### FILE: `google_document_ai_official_lifecycle/README.md`
```yaml
block_id: "GOOGLE-DOCAI-LIFECYCLE:file-01:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace lifecycle guide/lock/verifier constrained by exact Google upstream"
license: "LicenseRef-Workspace-Owner"
sha256: "27d6dc0c38acf4a840a60a9db6846f344c503abea65473a3a4f015ab557fe554"
variables: []
secrets_allowed: false
```
````markdown
# Google Document AI official lifecycle

This directory materializes Google-authored code for the processor lifecycle around document extraction. Sixteen files are byte-exact source/test pairs from signed commit `dc0eecc2187791fea70f48201241288e25fc9620` of `GoogleCloudPlatform/python-docs-samples`: create, train, evaluate, deploy, set-default, undeploy, batch processing and response handling. Batch accepts one GCS document or a complete prefix, waits for the long-running operation, inspects each document status, lists JSON result shards and reconstructs Document objects. Response handling covers OCR quality/layout, forms, tables, entities, splitter output, layout chunks and custom extraction. `copy_processor_versions_and_datasets_official_docs.py` is the single Python block extracted from Google's current English Cloud documentation page; its logic was not rewritten. It initializes/copies dataset schema, imports train/test/unassigned documents, imports a processor version, deploys it and makes it default.

Run `verify_official_lifecycle.py` after installing the exact admitted `google-cloud-documentai==3.15.0` and `google-cloud-storage==3.13.1` artifacts. It verifies every upstream hash and byte length, compiles the official sources, checks the required lifecycle, batch and response functions/API calls by AST, exercises the pure response helpers and executes the unmodified set-default function against a capture client without network. Google's eight test modules remain unmodified. Five lifecycle tests are fully mocked and pass offline; set-default, batch and response tests call live Google resources and must only run with explicit GCP/ADC/storage/quota/cost authority.

Use order:

1. complete the project intake for GCP project, region, processor type, identities, Cloud Storage, IAM, CMEK/retention, quota and cost;
2. install only the hash-locked SDK artifact selected by `OFFICIAL-DOCUMENT-SDK-ARTIFACT-CORE`;
3. run `verify_official_lifecycle.py` and the five offline-safe Google tests;
4. configure the Google functions with project-specific IDs and storage URIs;
5. create/initialize the processor and dataset, import reviewed ground-truth documents, train, evaluate against project thresholds, deploy, set default, batch-process approved GCS objects/prefixes, consume every per-document JSON result and retain undeploy as rollback;
6. execute live only after approvals and preserve operation/evaluation receipts.

This is official provider code, not a universal accuracy guarantee. It never grants automatic business persistence. Extracted fields still pass the library's ground-truth, strict evaluation, ambiguity/review and reconciliation gates before storage. The Cloud documentation script contains placeholders and a top-level `main` call; configure a project copy before execution and never import it as a side-effect-free module.
````

### FILE: `google_document_ai_official_lifecycle/source-lock.json`
```yaml
block_id: "GOOGLE-DOCAI-LIFECYCLE:file-02:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace lifecycle guide/lock/verifier constrained by exact Google upstream"
license: "LicenseRef-Workspace-Owner"
sha256: "4da1dc2777024579998989fcf42622c71881a81609f9d6660f05f8cf9a568fb2"
variables: []
secrets_allowed: false
```
````json
{
  "schema_version": "1.0",
  "verified_at": "2026-08-28",
  "python_docs_samples": {
    "repository": "https://github.com/GoogleCloudPlatform/python-docs-samples",
    "commit": "dc0eecc2187791fea70f48201241288e25fc9620",
    "commit_signature_verified": true,
    "license": "Apache-2.0",
    "files": [
      {"path": "documentai/snippets/create_processor_sample.py", "bytes": 2075, "git_blob": "c74dd744f74e5593a6f5fff595212773f190768f", "sha256": "8d0c19e8a2757d09f588994bc36df45ed10105e267be8d2d94ace70330ee755b"},
      {"path": "documentai/snippets/create_processor_sample_test.py", "bytes": 1539, "git_blob": "6c71021c70dea18683f19d40001ce928c4bf4d03", "sha256": "298f0462b9e34203eb53bfcb3b95023c1b157b5bea0ad30356bb9c8bffd1afd0"},
      {"path": "documentai/snippets/train_processor_version_sample.py", "bytes": 3113, "git_blob": "ac35cecbf44cf0778daec372bdecf421751bc7a3", "sha256": "02cfa9ba766b97941dbf4c218dbad5118b462bc022a2534ab01cd548922148c8"},
      {"path": "documentai/snippets/train_processor_version_sample_test.py", "bytes": 2113, "git_blob": "c9dcf81565d69e74a754c8189c0213dc385a1c34", "sha256": "622c63b6856d7c91d572dc8ce368cf87114eee632f48af70c7c19ac55e774374"},
      {"path": "documentai/snippets/evaluate_processor_version_sample.py", "bytes": 3227, "git_blob": "6a0adfb093b6d0549d4e297f810623b9be6fe554", "sha256": "f51e3fa46b66c79d04a362731e0c2bcb9992cb623fea873ddb8b969ccbd71141"},
      {"path": "documentai/snippets/evaluate_processor_version_sample_test.py", "bytes": 1833, "git_blob": "f7cf503316e842d6966f2c503a035fbde509439a", "sha256": "fd56d4a55d1ee490f29da50529642ce940a754774f7c04bb1bfff4ad91712703"},
      {"path": "documentai/snippets/deploy_processor_version_sample.py", "bytes": 2141, "git_blob": "f3943a843f7aa5592396d3a3ef167666ad2f9c59", "sha256": "048c87f38fd74b4d9792db14f8d81d336e988cd6c7c3a1e789802bfe3d4bc0d3"},
      {"path": "documentai/snippets/deploy_processor_version_sample_test.py", "bytes": 1548, "git_blob": "58e151ad9afc1ae75e6a3699ca61630f9a50703d", "sha256": "467a79113d6d5cdd8ecc6c6f92333c851c067766ef079f5296d57c4519841aed"},
      {"path": "documentai/snippets/set_default_processor_version_sample.py", "bytes": 2451, "git_blob": "1b3a3198565e5cee469d4951d573c881ccff6980", "sha256": "70dfc4620ef793413ec1e31d2d5af2b28cbf3c888e71eed0c7cb76151ea93706"},
      {"path": "documentai/snippets/set_default_processor_version_sample_test.py", "bytes": 1536, "git_blob": "b142e81336bde9b5a861efdd4ec117f980941b9a", "sha256": "a728bac30d2f09297f8e400b203bbaaf53417e810e39a1d6c3422cf6c5ddf400"},
      {"path": "documentai/snippets/undeploy_processor_version_sample.py", "bytes": 2334, "git_blob": "2086661c44db560935747933077d80be01076541", "sha256": "8163efdc6d0dfb8e0fd48a2ca2ac1a64a9e69e9c9d3334da3ac32f81553998e8"},
      {"path": "documentai/snippets/undeploy_processor_version_sample_test.py", "bytes": 1564, "git_blob": "41e1f5e51024d68c6c56e357b3263e89b1769381", "sha256": "a29bf337d22a317e59bba27c7ce3d1365412bdffd77b91cda6192939e2fe9f16"},
      {"path": "documentai/snippets/batch_process_documents_sample.py", "bytes": 7031, "git_blob": "6d15d8713a7aa96d6012f8b36e2fea2333ac94d0", "sha256": "3ccd1f828172c76f9268451091a90225dbd173cd7979755d5a66adebcea41759"},
      {"path": "documentai/snippets/batch_process_documents_sample_test.py", "bytes": 2568, "git_blob": "168fe3ac3ee003580a63d2eb4c7aa1eb26eb6253", "sha256": "674a525239061de0a8d4c7075e846f69d91713aff27940f70d3f68ca3a19f6ec"},
      {"path": "documentai/snippets/handle_response_sample.py", "bytes": 20533, "git_blob": "58bbb1debe09c6fad3814789f3db78af4679e061", "sha256": "8fc1adaceba3ac8ad4d871f8003b2ce43fbadfd9cc7b80599cefcc6edc8bfb91"},
      {"path": "documentai/snippets/handle_response_sample_test.py", "bytes": 7859, "git_blob": "b7c65834ccafbcb8fbae6f46807d0a7184ff1048", "sha256": "a60d4d9a53aa3c6383d60bb2cfc22489fde622d09ac994d68161ee771b1fea45"}
    ]
  },
  "google_cloud_documentation": {
    "url": "https://docs.cloud.google.com/document-ai/docs/copy-processor-versions?hl=en",
    "last_modified": "Wed, 26 Aug 2026 23:31:34 GMT",
    "page_sha256": "21c4fd07e8931e3549a085e51f1294f130dca70ea35b4007f08e4b3c5687a625",
    "code_license": "Apache-2.0",
    "extraction": "Select the only pre/code DOM block whose decoded text contains def import_documents; remove markup tags, HTML-decode text, write UTF-8 without BOM; no code logic edited.",
    "matched_code_blocks": 1,
    "extracted_path": "copy_processor_versions_and_datasets_official_docs.py",
    "bytes": 18662,
    "sha256": "10897a366d7d7201f4ad4901e1b14d16472dc8f7e614b1411eb32b3c9361f232"
  },
  "license_file": {
    "path": "LICENSE.txt",
    "bytes": 11357,
    "sha256": "c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4"
  }
}
````

### FILE: `google_document_ai_official_lifecycle/upstream/google/cloud-docs/copy_processor_versions_and_datasets_official_docs.py`
```yaml
block_id: "GOOGLE-DOCAI-LIFECYCLE:file-03:v1"
operation: CREATE
provenance: ADAPTED
source: "https://docs.cloud.google.com/document-ai/docs/copy-processor-versions?hl=en (single decoded DOM code block; logic unchanged)"
license: "Apache-2.0"
sha256: "10897a366d7d7201f4ad4901e1b14d16472dc8f7e614b1411eb32b3c9361f232"
variables: []
secrets_allowed: false
```
````python
import time
from pathlib import Path
from typing import Optional, Tuple
from google.cloud.documentai_v1beta3.services.document_service import pagers
from google.api_core.client_options import ClientOptions
from google.api_core.operation import Operation
from google.cloud import documentai_v1beta3 as documentai
from google.cloud import storage
from tqdm import tqdm

source_project_id = "source-project-id"
source_location = "processor-location"
source_processor_id = "source-processor-id"
source_processor_version_to_import = "source-processor-version-id"
migrate_dataset = False # Either True or False
source_exported_gcs_path = (
    "gs://bucket/path/to/export_dataset/"
)
destination_project_id = "< destination-project-id >"
# Give empty string if you wish to create a new processor
destination_processor_id = ""

exported_bucket_name = source_exported_gcs_path.split("/")[2]
exported_bucket_path_prefix = "/".join(source_exported_gcs_path.split("/")[3:])
destination_location = source_location

def sample_get_processor(project_id: str, processor_id: str, location: str)->Tuple[str, str]:
    """
    This function returns Processor Display Name and Type of Processor from source project

    Args:
        project_id (str): Project ID
        processor_id (str): Document AI Processor ID
        location (str): Processor Location

    Returns:
        Tuple[str, str]: Returns Processor Display name and type
    """
    client = documentai.DocumentProcessorServiceClient()
    print(
        f"Fetching processor({processor_id}) details from source project ({project_id})"
    )
    name = f"projects/{project_id}/locations/{location}/processors/{processor_id}"
    request = documentai.GetProcessorRequest(
        name=name,
    )
    response = client.get_processor(request=request)
    print(f"Processor Name: {response.name}")
    print(f"Processor Display Name: {response.display_name}")
    print(f"Processor Type: {response.type_}")
    return response.display_name, response.type_

def sample_create_processor(project_id: str, location: str, display_name: str, processor_type: str)->documentai.Processor:
    """It will create Processor in Destination project

    Args:
        project_id (str): Project ID
        location (str): Location fo processor
        display_name (str): Processor Display Name
        processor_type (str): Google Cloud Document AI Processor type

    Returns:
        documentai.Processor: Returns details abouts newly created processor
    """

    client = documentai.DocumentProcessorServiceClient()
    request = documentai.CreateProcessorRequest(
        parent=f"projects/{project_id}/locations/{location}",
        processor={
            "type_": processor_type,
            "display_name": display_name,
        },
    )
    print(f"Creating Processor in project: {project_id} in location: {location}")
    print(f"Display Name: {display_name} & Processor Type: {processor_type}")
    res = client.create_processor(request=request)
    return res

def initialize_dataset(project_id: str, processor_id: str, location: str)-> Operation:
    """It will configure dataset for target processor in destination project

    Args:
        project_id (str): Project ID
        processor_id (str): DocuemntAI Processor ID
        location (str): Processor Location

    Returns:
        Operation: An object representing a long-running operation
    """

    # opts = ClientOptions(api_endpoint=f"{location}-documentai.googleapis.com")
    client = documentai.DocumentServiceClient()  # client_options=opts
    dataset = documentai.types.Dataset(
        name=f"projects/{project_id}/locations/{location}/processors/{processor_id}/dataset",
        state=3,
        unmanaged_dataset_config={},
        spanner_indexing_config={},
    )
    request = documentai.types.UpdateDatasetRequest(dataset=dataset)
    print(
        f"Configuring Dataset in  project: {project_id} for processor: {processor_id}"
    )
    response = client.update_dataset(request=request)
    return response

def get_dataset_schema(project_id: str, processor_id: str, location: str)->documentai.DatasetSchema:
    """It helps to fetch processor schema

    Args:
        project_id (str): Project ID
        processor_id (str): DocumentAI Processor ID
        location (str): Processor Location

    Returns:
        documentai.DatasetSchema: Return deails about Processor Dataset Schema
    """

    # Create a client
    processor_name = (
        f"projects/{project_id}/locations/{location}/processors/{processor_id}"
    )
    client = documentai.DocumentServiceClient()
    request = documentai.GetDatasetSchemaRequest(
        name=processor_name + "/dataset/datasetSchema"
    )
    # Make the request
    print(f"Fetching schema from source processor: {processor_id}")
    response = client.get_dataset_schema(request=request)
    return response

def upload_dataset_schema(schema: documentai.DatasetSchema)->documentai.DatasetSchema:
    """It helps to update the schema in destination processor

    Args:
        schema (documentai.DatasetSchema): Document AI Processor Schema details & Metadata

    Returns:
        documentai.DatasetSchema: Returns Dataset Schema object
    """

    client = documentai.DocumentServiceClient()
    request = documentai.UpdateDatasetSchemaRequest(dataset_schema=schema)
    print("Updating Schema in destination processor")
    res = client.update_dataset_schema(request=request)
    return res

def store_document_as_json(document: str, bucket_name: str, file_name: str)->None:
    """It helps to upload data to Cloud Storage and stores as a blob

    Args:
        document (str): Processor response in json string format
        bucket_name (str): Cloud Storage bucket name
        file_name (str): Cloud Storage blob uri
    """

    print(f"\tUploading file to Cloud Storage gs://{bucket_name}/{file_name}")
    storage_client = storage.Client()
    process_result_bucket = storage_client.get_bucket(bucket_name)
    document_blob = storage.Blob(
        name=str(Path(file_name)), bucket=process_result_bucket
    )
    document_blob.upload_from_string(document, content_type="application/json")

def list_documents(project_id: str, location: str, processor: str, page_size: Optional[int]=100, page_token: Optional[str]="")->pagers.ListDocumentsPager:
    """This function helps to list the samples present in processor dataset

    Args:
        project_id (str): Project ID
        location (str): Processor Location 
        processor (str): DocumentAI Processor ID
        page_size (Optional[int], optional): The maximum number of documents to return. Defaults to 100.
        page_token (Optional[str], optional): A page token, received from a previous ListDocuments call. Defaults to "".

    Returns:
        pagers.ListDocumentsPager: Returns all details about documents present in Processor Dataset
    """
    client = documentai.DocumentServiceClient()
    dataset = (
        f"projects/{project_id}/locations/{location}/processors/{processor}/dataset"
    )
    request = documentai.types.ListDocumentsRequest(
        dataset=dataset,
        page_token=page_token,
        page_size=page_size,
        return_total_size=True,
    )
    print(f"Listingll  documents/Samples present in processor: {processor}")
    operation = client.list_documents(request)
    return operation

def get_document(project_id: str, location: str, processor: str, doc_id: documentai.DocumentId)->documentai.GetDocumentResponse:
    """It will fetch data for individual sample/document present in dataset

    Args:
        project_id (str): Project ID
        location (str): Processor Location
        processor (str): Document AI Processor ID
        doc_id (documentai.DocumentId): Document identifier

    Returns:
        documentai.GetDocumentResponse: Returns data related to doc_id
    """

    client = documentai.DocumentServiceClient()
    dataset = (
        f"projects/{project_id}/locations/{location}/processors/{processor}/dataset"
    )
    request = documentai.GetDocumentRequest(dataset=dataset, document_id=doc_id)
    operation = client.get_document(request)
    return operation

def import_documents(project_id: str, processor_id: str, location: str, gcs_path: str)->Operation:
    """It helps to import samples/docuemnts from Cloud Storage path to processor via API call

    Args:
        project_id (str): Project ID
        processor_id (str): Document AI Processor ID
        location (str): Processor Location
        gcs_path (str): Cloud Storage path uri prefix 

    Returns:
        Operation: An object representing a long-running operation
    """

    client = documentai.DocumentServiceClient()
    dataset = (
        f"projects/{project_id}/locations/{location}/processors/{processor_id}/dataset"
    )
    request = documentai.ImportDocumentsRequest(
        dataset=dataset,
        batch_documents_import_configs=[
            {
                "dataset_split": "DATASET_SPLIT_TRAIN",
                "batch_input_config": {
                    "gcs_prefix": {"gcs_uri_prefix": gcs_path + "train/"}
                },
            },
            {
                "dataset_split": "DATASET_SPLIT_TEST",
                "batch_input_config": {
                    "gcs_prefix": {"gcs_uri_prefix": gcs_path + "test/"}
                },
            },
            {
                "dataset_split": "DATASET_SPLIT_UNASSIGNED",
                "batch_input_config": {
                    "gcs_prefix": {"gcs_uri_prefix": gcs_path + "unassigned/"}
                },
            },
        ],
    )
    print(
        f"Importing Documents/samples from {gcs_path} to corresponding tran_test_unassigned sections"
    )
    response = client.import_documents(request=request)
    return response

def import_processor_version(source_processor_version_name: str, destination_processor_name: str)->Operation:
    """It helps to import processor version from source processor to destanation processor

    Args:
        source_processor_version_name (str): source processor name in this format projects/{project}/locations/{location}/processors/{processor}
        destination_processor_name (str): destination processor name in this format projects/{project}/locations/{location}/processors/{processor}

    Returns:
        Operation: An object representing a long-running operation
    """

    from google.cloud import documentai_v1beta3

    # provide the source version(to copy) processor details in the following format
    client = documentai_v1beta3.DocumentProcessorServiceClient()
    # provide the new processor name in the parent variable in format 'projects/{project_number}/locations/{location}/processors/{new_processor_id}'
    import google.cloud.documentai_v1beta3 as documentai

    op_import_version_req = (
        documentai.types.document_processor_service.ImportProcessorVersionRequest(
            processor_version_source=source_processor_version_name,
            parent=destination_processor_name,
        )
    )

    print("Importing processor from source to destination")
    print(f"\tSource: {source_processor_version_name}")
    print(f"\tDestination: {destination_processor_name}")
    # copying the processor
    operation = client.import_processor_version(request=op_import_version_req)
    print(operation.metadata)
    print("Waitin for operation to complete...")
    operation.result()
    return operation

def deploy_and_set_default_processor_version(
    project_id: str, location: str, processor_id: str, processor_version_id: str
)->None:
    """It helps to deploy to imported processor version and set it as default version

    Args:
        project_id (str): Project ID
        location (str): Processor Location
        processor_id (str): Document AI Processor ID
        processor_version_id (str): Document AI Processor Version ID
    """

    # Construct the resource name of the processor version
    processor_name = (
        f"projects/{project_id}/locations/{location}/processors/{processor_id}"
    )
    default_processor_version_name = f"projects/{project_id}/locations/{location}/processors/{processor_id}/processorVersions/{processor_version_id}"

    # Initialize the Document AI client
    client_options = ClientOptions(api_endpoint=f"{location}-documentai.googleapis.com")
    client = documentai.DocumentProcessorServiceClient(client_options=client_options)

    # Deploy the processor version
    operation = client.deploy_processor_version(name=default_processor_version_name)
    print(f"Deploying processor version: {operation.operation.name}")
    print("Waiting for operation to complete...")
    result = operation.result()
    print("Processor version deployed")
    # Set the deployed version as the default version
    request = documentai.SetDefaultProcessorVersionRequest(
        processor=processor_name,
        default_processor_version=default_processor_version_name,
    )
    operation = client.set_default_processor_version(request=request)
    print(f"Setting default processor version: {operation.operation.name}")
    operation.result()
    print(f"Default processor version set {default_processor_version_name}")

def main(destination_processor_id: str, migrate_dataset: bool = False)->None:
    """Entry function to perform Processor Migration from Source Project to Destination project

    Args:
        destination_processor_id (str): Either empty string or processor id in desination project
    """
    # Checking processor id of destination project
    if destination_processor_id == "":
        # Fetching Processor Display Name and Type of Processor from source project
        display_name, processor_type = sample_get_processor(
            source_project_id, source_processor_id, source_location
        )
        # Creating Processor in Destination project
        des_processor = sample_create_processor(
            destination_project_id, destination_location, display_name, processor_type
        )
        print(des_processor)
        destination_processor_id = des_processor.name.split("/")[-1]
        # configuring dataset for target processor in destination project
        r = initialize_dataset(
            destination_project_id, destination_processor_id, destination_location
        )
    # fetching processor schema from source processor
    exported_schema = get_dataset_schema(
        source_project_id, source_processor_id, source_location
    )
    exported_schema.name = f"projects/{destination_project_id}/locations/{destination_location}/processors/{destination_processor_id}/dataset/datasetSchema"
    # Copying schema from source processor to desination processor 
    import_schema = upload_dataset_schema(exported_schema)
    if migrate_dataset == True: # to migrate dataset from source to destination processor
        print("Migrating Dataset from source to destination processor")
        # Fetching/listing the samples/JSONs present in source processor dataset
        results = list_documents(source_project_id, source_location, source_processor_id)
        document_list = results.document_metadata
        while len(document_list) != results.total_size:
            page_token = results.next_page_token
            results = list_documents(
                source_project_id,
                source_location,
                source_processor_id,
                page_token=page_token,
            )
            document_list.extend(results.document_metadata)
        print("Exporting Dataset...")
        for doc in tqdm(document_list):
            doc_id = doc.document_id
            split_type = doc.dataset_type
            if split_type == 3:
                split = "unassigned"
            elif split_type == 2:
                split = "test"
            elif split_type == 1:
                split = "train"
            else:
                split = "unknown"
            file_name = doc.display_name
            # fetching/downloading data for individual sample/document present in dataset
            res = get_document(
                source_project_id, source_location, source_processor_id, doc_id
            )
            output_file_name = (
                f"{exported_bucket_path_prefix.strip('/')}/{split}/{file_name}.json"
            )
            # Converting Document AI Proto object to JSON string
            json_data = documentai.Document.to_json(res.document)
            # Uploading JSON data to specified Cloud Storage path
            store_document_as_json(json_data, exported_bucket_name, output_file_name)

        print(f"Importing dataset to {destination_processor_id}")
        gcs_path = source_exported_gcs_path.strip("/") + "/"
        project = destination_project_id
        location = destination_location
        processor = destination_processor_id
        # importing samples/docuemnts from Cloud Storage path to destination processor
        res = import_documents(project, processor, location, gcs_path)
        print(f"Waiting for {len(document_list)*1.5} seconds")
        time.sleep(len(document_list) * 1.5)
    else:
        print("\tSkipping Dataset Migration actions like, exporting source dataset to Cloud Storage and importing dataset to destination processor")

    # Checking for source processor version, if id provided then it will be imported to destination processor
    if source_processor_version_to_import != "":
        print(f"Importing Processor Version {source_processor_version_to_import}")
        source_version = f"projects/{source_project_id}/locations/{source_location}/processors/{source_processor_id}/processorVersions/{source_processor_version_to_import}"
        destination_version = f"projects/{destination_project_id}/locations/{destination_location}/processors/{destination_processor_id}"
        # source_version = f"projects/{source_project_id}/locations/us/processors/a82fc086440d7ea1/processorVersions/f1eeed93aad5e317"  # Data for testing
        # Importing processor version from source processor to destanation processor
        operation = import_processor_version(source_version, destination_version)
        name = operation.metadata.common_metadata.resource
        destination_processor_version_id = name.split("/")[-1]
        # deploying newly imported processor version and set it as default version in desination project
        deploy_and_set_default_processor_version(
            destination_project_id,
            destination_location,
            destination_processor_id,
            destination_processor_version_id,
        )

main(destination_processor_id, migrate_dataset)
print("Process Completed!!!")

````

### FILE: `google_document_ai_official_lifecycle/upstream/google/documentai/snippets/create_processor_sample_test.py`
```yaml
block_id: "GOOGLE-DOCAI-LIFECYCLE:file-04:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/GoogleCloudPlatform/python-docs-samples/dc0eecc2187791fea70f48201241288e25fc9620/documentai/snippets/create_processor_sample_test.py"
license: "Apache-2.0"
sha256: "298f0462b9e34203eb53bfcb3b95023c1b157b5bea0ad30356bb9c8bffd1afd0"
variables: []
secrets_allowed: false
```
````python
# Copyright 2022 Google LLC
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
from unittest import mock
from uuid import uuid4

from documentai.snippets import create_processor_sample

location = "us"
project_id = os.environ["GOOGLE_CLOUD_PROJECT"]
processor_display_name = f"test-processor-{uuid4()}"
processor_type = "OCR_PROCESSOR"


@mock.patch("google.cloud.documentai.DocumentProcessorServiceClient.create_processor")
@mock.patch("google.cloud.documentai.Processor")
def test_create_processor(create_processor_mock, processor_mock, capsys):
    create_processor_mock.return_value = processor_mock

    create_processor_sample.create_processor_sample(
        project_id=project_id,
        location=location,
        processor_display_name=processor_display_name,
        processor_type=processor_type,
    )

    create_processor_mock.assert_called_once()

    out, _ = capsys.readouterr()

    assert "Processor Name:" in out
    assert "Processor Display Name:" in out
    assert "Processor Type:" in out
````

### FILE: `google_document_ai_official_lifecycle/upstream/google/documentai/snippets/create_processor_sample.py`
```yaml
block_id: "GOOGLE-DOCAI-LIFECYCLE:file-05:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/GoogleCloudPlatform/python-docs-samples/dc0eecc2187791fea70f48201241288e25fc9620/documentai/snippets/create_processor_sample.py"
license: "Apache-2.0"
sha256: "8d0c19e8a2757d09f588994bc36df45ed10105e267be8d2d94ace70330ee755b"
variables: []
secrets_allowed: false
```
````python
# Copyright 2022 Google LLC
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

# [START documentai_create_processor]

from google.api_core.client_options import ClientOptions
from google.cloud import documentai  # type: ignore

# TODO(developer): Uncomment these variables before running the sample.
# project_id = 'YOUR_PROJECT_ID'
# location = 'YOUR_PROCESSOR_LOCATION' # Format is 'us' or 'eu'
# processor_display_name = 'YOUR_PROCESSOR_DISPLAY_NAME' # Must be unique per project, e.g.: 'My Processor'
# processor_type = 'YOUR_PROCESSOR_TYPE' # Use fetch_processor_types to get available processor types


def create_processor_sample(
    project_id: str, location: str, processor_display_name: str, processor_type: str
) -> None:
    # You must set the api_endpoint if you use a location other than 'us'.
    opts = ClientOptions(api_endpoint=f"{location}-documentai.googleapis.com")

    client = documentai.DocumentProcessorServiceClient(client_options=opts)

    # The full resource name of the location
    # e.g.: projects/project_id/locations/location
    parent = client.common_location_path(project_id, location)

    # Create a processor
    processor = client.create_processor(
        parent=parent,
        processor=documentai.Processor(
            display_name=processor_display_name, type_=processor_type
        ),
    )

    # Print the processor information
    print(f"Processor Name: {processor.name}")
    print(f"Processor Display Name: {processor.display_name}")
    print(f"Processor Type: {processor.type_}")


# [END documentai_create_processor]
````

### FILE: `google_document_ai_official_lifecycle/upstream/google/documentai/snippets/deploy_processor_version_sample_test.py`
```yaml
block_id: "GOOGLE-DOCAI-LIFECYCLE:file-06:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/GoogleCloudPlatform/python-docs-samples/dc0eecc2187791fea70f48201241288e25fc9620/documentai/snippets/deploy_processor_version_sample_test.py"
license: "Apache-2.0"
sha256: "467a79113d6d5cdd8ecc6c6f92333c851c067766ef079f5296d57c4519841aed"
variables: []
secrets_allowed: false
```
````python
# Copyright 2022 Google LLC
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
from unittest import mock

from documentai.snippets import deploy_processor_version_sample

location = "us"
project_id = os.environ["GOOGLE_CLOUD_PROJECT"]
processor_id = "aaaaaaaaa"
processor_version_id = "xxxxxxxxxx"


# TODO: Switch to Real Endpoint when Deployable Versions are Available
@mock.patch(
    "google.cloud.documentai.DocumentProcessorServiceClient.deploy_processor_version"
)
@mock.patch("google.api_core.operation.Operation")
def test_deploy_processor_version(
    operation_mock, deploy_processor_version_mock, capsys
):
    deploy_processor_version_mock.return_value = operation_mock

    deploy_processor_version_sample.deploy_processor_version_sample(
        project_id=project_id,
        location=location,
        processor_id=processor_id,
        processor_version_id=processor_version_id,
    )

    deploy_processor_version_mock.assert_called_once()

    out, _ = capsys.readouterr()

    assert "operation" in out
````

### FILE: `google_document_ai_official_lifecycle/upstream/google/documentai/snippets/deploy_processor_version_sample.py`
```yaml
block_id: "GOOGLE-DOCAI-LIFECYCLE:file-07:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/GoogleCloudPlatform/python-docs-samples/dc0eecc2187791fea70f48201241288e25fc9620/documentai/snippets/deploy_processor_version_sample.py"
license: "Apache-2.0"
sha256: "048c87f38fd74b4d9792db14f8d81d336e988cd6c7c3a1e789802bfe3d4bc0d3"
variables: []
secrets_allowed: false
```
````python
# Copyright 2022 Google LLC
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

# [START documentai_deploy_processor_version]

from google.api_core.client_options import ClientOptions
from google.api_core.exceptions import FailedPrecondition
from google.cloud import documentai  # type: ignore

# TODO(developer): Uncomment these variables before running the sample.
# project_id = 'YOUR_PROJECT_ID'
# location = 'YOUR_PROCESSOR_LOCATION' # Format is 'us' or 'eu'
# processor_id = 'YOUR_PROCESSOR_ID'
# processor_version_id = 'YOUR_PROCESSOR_VERSION_ID'


def deploy_processor_version_sample(
    project_id: str, location: str, processor_id: str, processor_version_id: str
) -> None:
    # You must set the api_endpoint if you use a location other than 'us'.
    opts = ClientOptions(api_endpoint=f"{location}-documentai.googleapis.com")

    client = documentai.DocumentProcessorServiceClient(client_options=opts)

    # The full resource name of the processor version
    # e.g.: projects/project_id/locations/location/processors/processor_id/processorVersions/processor_version_id
    name = client.processor_version_path(
        project_id, location, processor_id, processor_version_id
    )

    # Make DeployProcessorVersion request
    try:
        operation = client.deploy_processor_version(name=name)
        # Print operation details
        print(operation.operation.name)
        # Wait for operation to complete
        operation.result()
    # Deploy request will fail if the
    # processor version is already deployed
    except FailedPrecondition as e:
        print(e.message)


# [END documentai_deploy_processor_version]
````

### FILE: `google_document_ai_official_lifecycle/upstream/google/documentai/snippets/evaluate_processor_version_sample_test.py`
```yaml
block_id: "GOOGLE-DOCAI-LIFECYCLE:file-08:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/GoogleCloudPlatform/python-docs-samples/dc0eecc2187791fea70f48201241288e25fc9620/documentai/snippets/evaluate_processor_version_sample_test.py"
license: "Apache-2.0"
sha256: "fd56d4a55d1ee490f29da50529642ce940a754774f7c04bb1bfff4ad91712703"
variables: []
secrets_allowed: false
```
````python
# Copyright 2022 Google LLC
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
from unittest import mock

from documentai.snippets import evaluate_processor_version_sample

location = "us"
project_id = os.environ["GOOGLE_CLOUD_PROJECT"]
processor_id = "aaaaaaaaa"
processor_version_id = "xxxxxxxxxx"
gcs_input_uri = "gs://bucket/directory/"


# Mocking request as evaluation can take a long time
@mock.patch(
    "google.cloud.documentai.DocumentProcessorServiceClient.evaluate_processor_version"
)
@mock.patch("google.cloud.documentai.EvaluateProcessorVersionResponse")
@mock.patch("google.api_core.operation.Operation")
def test_evaluate_processor_version(
    operation_mock,
    evaluate_processor_version_response_mock,
    evaluate_processor_version_mock,
    capsys,
):
    operation_mock.result.return_value = evaluate_processor_version_response_mock
    evaluate_processor_version_mock.return_value = operation_mock

    evaluate_processor_version_sample.evaluate_processor_version_sample(
        project_id=project_id,
        location=location,
        processor_id=processor_id,
        processor_version_id=processor_version_id,
        gcs_input_uri=gcs_input_uri,
    )

    evaluate_processor_version_mock.assert_called_once()

    out, _ = capsys.readouterr()

    assert "operation" in out
````

### FILE: `google_document_ai_official_lifecycle/upstream/google/documentai/snippets/evaluate_processor_version_sample.py`
```yaml
block_id: "GOOGLE-DOCAI-LIFECYCLE:file-09:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/GoogleCloudPlatform/python-docs-samples/dc0eecc2187791fea70f48201241288e25fc9620/documentai/snippets/evaluate_processor_version_sample.py"
license: "Apache-2.0"
sha256: "f51e3fa46b66c79d04a362731e0c2bcb9992cb623fea873ddb8b969ccbd71141"
variables: []
secrets_allowed: false
```
````python
# Copyright 2022 Google LLC
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

# [START documentai_evaluate_processor_version]

from google.api_core.client_options import ClientOptions
from google.cloud import documentai  # type: ignore

# TODO(developer): Uncomment these variables before running the sample.
# project_id = 'YOUR_PROJECT_ID'
# location = 'YOUR_PROCESSOR_LOCATION' # Format is 'us' or 'eu'
# processor_id = 'YOUR_PROCESSOR_ID'
# processor_version_id = 'YOUR_PROCESSOR_VERSION_ID'
# gcs_input_uri = # Format: gs://bucket/directory/


def evaluate_processor_version_sample(
    project_id: str,
    location: str,
    processor_id: str,
    processor_version_id: str,
    gcs_input_uri: str,
) -> None:
    # You must set the api_endpoint if you use a location other than 'us', e.g.:
    opts = ClientOptions(api_endpoint=f"{location}-documentai.googleapis.com")

    client = documentai.DocumentProcessorServiceClient(client_options=opts)

    # The full resource name of the processor version
    # e.g. `projects/{project_id}/locations/{location}/processors/{processor_id}/processorVersions/{processor_version_id}`
    name = client.processor_version_path(
        project_id, location, processor_id, processor_version_id
    )

    evaluation_documents = documentai.BatchDocumentsInputConfig(
        gcs_prefix=documentai.GcsPrefix(gcs_uri_prefix=gcs_input_uri)
    )

    # NOTE: Alternatively, specify a list of GCS Documents
    #
    # gcs_input_uri = "gs://bucket/directory/file.pdf"
    # input_mime_type = "application/pdf"
    #
    # gcs_document = documentai.GcsDocument(
    #     gcs_uri=gcs_input_uri, mime_type=input_mime_type
    # )
    # gcs_documents = [gcs_document]
    # evaluation_documents = documentai.BatchDocumentsInputConfig(
    #     gcs_documents=documentai.GcsDocuments(documents=gcs_documents)
    # )
    #

    request = documentai.EvaluateProcessorVersionRequest(
        processor_version=name,
        evaluation_documents=evaluation_documents,
    )

    # Make EvaluateProcessorVersion request
    # Continually polls the operation until it is complete.
    # This could take some time for larger files
    operation = client.evaluate_processor_version(request=request)
    # Print operation details
    # Format: projects/PROJECT_NUMBER/locations/LOCATION/operations/OPERATION_ID
    print(f"Waiting for operation {operation.operation.name} to complete...")
    # Wait for operation to complete
    response = documentai.EvaluateProcessorVersionResponse(operation.result())

    # After the operation is complete,
    # Print evaluation ID from operation response
    print(f"Evaluation Complete: {response.evaluation}")


# [END documentai_evaluate_processor_version]
````

### FILE: `google_document_ai_official_lifecycle/upstream/google/documentai/snippets/set_default_processor_version_sample_test.py`
```yaml
block_id: "GOOGLE-DOCAI-LIFECYCLE:file-10:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/GoogleCloudPlatform/python-docs-samples/dc0eecc2187791fea70f48201241288e25fc9620/documentai/snippets/set_default_processor_version_sample_test.py"
license: "Apache-2.0"
sha256: "a728bac30d2f09297f8e400b203bbaaf53417e810e39a1d6c3422cf6c5ddf400"
variables: []
secrets_allowed: false
```
````python
# Copyright 2022 Google LLC
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

from documentai.snippets import set_default_processor_version_sample

location = "us"
project_id = os.environ["GOOGLE_CLOUD_PROJECT"]
processor_id = "aeb8cea219b7c272"
current_default_processor_version = "pretrained-ocr-v1.0-2020-09-23"
new_default_processor_version = "pretrained-ocr-v1.1-2022-09-12"


def test_set_default_processor_version(capsys):
    set_default_processor_version_sample.set_default_processor_version_sample(
        project_id=project_id,
        location=location,
        processor_id=processor_id,
        processor_version_id=new_default_processor_version,
    )
    out, _ = capsys.readouterr()

    assert "operation" in out

    # Set back to previous default
    set_default_processor_version_sample.set_default_processor_version_sample(
        project_id=project_id,
        location=location,
        processor_id=processor_id,
        processor_version_id=current_default_processor_version,
    )
````

### FILE: `google_document_ai_official_lifecycle/upstream/google/documentai/snippets/set_default_processor_version_sample.py`
```yaml
block_id: "GOOGLE-DOCAI-LIFECYCLE:file-11:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/GoogleCloudPlatform/python-docs-samples/dc0eecc2187791fea70f48201241288e25fc9620/documentai/snippets/set_default_processor_version_sample.py"
license: "Apache-2.0"
sha256: "70dfc4620ef793413ec1e31d2d5af2b28cbf3c888e71eed0c7cb76151ea93706"
variables: []
secrets_allowed: false
```
````python
# Copyright 2022 Google LLC
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

# [START documentai_set_default_processor_version]

from google.api_core.client_options import ClientOptions
from google.api_core.exceptions import NotFound
from google.cloud import documentai  # type: ignore

# TODO(developer): Uncomment these variables before running the sample.
# project_id = 'YOUR_PROJECT_ID'
# location = 'YOUR_PROCESSOR_LOCATION' # Format is 'us' or 'eu'
# processor_id = 'YOUR_PROCESSOR_ID' # Create processor before running sample
# processor_version_id = 'YOUR_PROCESSOR_VERSION_ID'


def set_default_processor_version_sample(
    project_id: str, location: str, processor_id: str, processor_version_id: str
) -> None:
    # You must set the api_endpoint if you use a location other than 'us'.
    opts = ClientOptions(api_endpoint=f"{location}-documentai.googleapis.com")

    client = documentai.DocumentProcessorServiceClient(client_options=opts)

    # The full resource name of the processor
    # e.g.: projects/project_id/locations/location/processors/processor_id
    processor = client.processor_path(project_id, location, processor_id)

    # The full resource name of the processor version
    # e.g.: projects/project_id/locations/location/processors/processor_id/processorVersions/processor_version_id
    processor_version = client.processor_version_path(
        project_id, location, processor_id, processor_version_id
    )

    request = documentai.SetDefaultProcessorVersionRequest(
        processor=processor, default_processor_version=processor_version
    )

    # Make SetDefaultProcessorVersion request
    try:
        operation = client.set_default_processor_version(request)
        # Print operation details
        print(operation.operation.name)
        # Wait for operation to complete
        operation.result()
    except NotFound as e:
        print(e.message)


# [END documentai_set_default_processor_version]
````

### FILE: `google_document_ai_official_lifecycle/upstream/google/documentai/snippets/train_processor_version_sample_test.py`
```yaml
block_id: "GOOGLE-DOCAI-LIFECYCLE:file-12:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/GoogleCloudPlatform/python-docs-samples/dc0eecc2187791fea70f48201241288e25fc9620/documentai/snippets/train_processor_version_sample_test.py"
license: "Apache-2.0"
sha256: "622c63b6856d7c91d572dc8ce368cf87114eee632f48af70c7c19ac55e774374"
variables: []
secrets_allowed: false
```
````python
# Copyright 2022 Google LLC
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
from unittest import mock

from documentai.snippets import train_processor_version_sample

location = "us"
project_id = os.environ["GOOGLE_CLOUD_PROJECT"]
processor_id = "aaaaaaaaa"
processor_version_display_name = "new-processor-version"
train_data_uri = "gs://bucket/directory/"
test_data_uri = "gs://bucket/directory/"


# Mocking request as training can take a long time
@mock.patch(
    "google.cloud.documentai.DocumentProcessorServiceClient.train_processor_version"
)
@mock.patch("google.cloud.documentai.TrainProcessorVersionResponse")
@mock.patch("google.cloud.documentai.TrainProcessorVersionMetadata")
@mock.patch("google.api_core.operation.Operation")
def test_train_processor_version(
    operation_mock,
    train_processor_version_metadata_mock,
    train_processor_version_response_mock,
    train_processor_version_mock,
    capsys,
):
    operation_mock.result.return_value = train_processor_version_response_mock
    operation_mock.metadata.return_value = train_processor_version_metadata_mock
    train_processor_version_mock.return_value = operation_mock

    train_processor_version_sample.train_processor_version_sample(
        project_id=project_id,
        location=location,
        processor_id=processor_id,
        processor_version_display_name=processor_version_display_name,
        train_data_uri=train_data_uri,
        test_data_uri=test_data_uri,
    )

    train_processor_version_mock.assert_called_once()

    out, _ = capsys.readouterr()

    assert "operation" in out
````

### FILE: `google_document_ai_official_lifecycle/upstream/google/documentai/snippets/train_processor_version_sample.py`
```yaml
block_id: "GOOGLE-DOCAI-LIFECYCLE:file-13:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/GoogleCloudPlatform/python-docs-samples/dc0eecc2187791fea70f48201241288e25fc9620/documentai/snippets/train_processor_version_sample.py"
license: "Apache-2.0"
sha256: "02cfa9ba766b97941dbf4c218dbad5118b462bc022a2534ab01cd548922148c8"
variables: []
secrets_allowed: false
```
````python
# Copyright 2022 Google LLC
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

# [START documentai_train_processor_version]

from typing import Optional

from google.api_core.client_options import ClientOptions
from google.cloud import documentai  # type: ignore

# TODO(developer): Uncomment these variables before running the sample.
# project_id = 'YOUR_PROJECT_ID'
# location = 'YOUR_PROCESSOR_LOCATION' # Format is 'us' or 'eu'
# processor_id = 'YOUR_PROCESSOR_ID'
# processor_version_display_name = 'new-processor-version'
# train_data_uri = 'gs://bucket/directory/' # (Optional)
# test_data_uri = 'gs://bucket/directory/' # (Optional)


def train_processor_version_sample(
    project_id: str,
    location: str,
    processor_id: str,
    processor_version_display_name: str,
    train_data_uri: Optional[str] = None,
    test_data_uri: Optional[str] = None,
) -> None:
    # You must set the api_endpoint if you use a location other than 'us', e.g.:
    opts = ClientOptions(api_endpoint=f"{location}-documentai.googleapis.com")

    client = documentai.DocumentProcessorServiceClient(client_options=opts)

    # The full resource name of the processor
    # e.g. `projects/{project_id}/locations/{location}/processors/{processor_id}
    parent = client.processor_path(project_id, location, processor_id)

    processor_version = documentai.ProcessorVersion(
        display_name=processor_version_display_name
    )

    # If train/test data is not supplied, the default sets in the Cloud Console will be used
    input_data = documentai.TrainProcessorVersionRequest.InputData(
        training_documents=documentai.BatchDocumentsInputConfig(
            gcs_prefix=documentai.GcsPrefix(gcs_uri_prefix=train_data_uri)
        ),
        test_documents=documentai.BatchDocumentsInputConfig(
            gcs_prefix=documentai.GcsPrefix(gcs_uri_prefix=test_data_uri)
        ),
    )

    request = documentai.TrainProcessorVersionRequest(
        parent=parent, processor_version=processor_version, input_data=input_data
    )

    operation = client.train_processor_version(request=request)
    # Print operation details
    print(operation.operation.name)
    # Wait for operation to complete
    response = documentai.TrainProcessorVersionResponse(operation.result())

    metadata = documentai.TrainProcessorVersionMetadata(operation.metadata)

    print(f"New Processor Version:{response.processor_version}")
    print(f"Training Set Validation: {metadata.training_dataset_validation}")
    print(f"Test Set Validation: {metadata.test_dataset_validation}")


# [END documentai_train_processor_version]
````

### FILE: `google_document_ai_official_lifecycle/upstream/google/documentai/snippets/undeploy_processor_version_sample_test.py`
```yaml
block_id: "GOOGLE-DOCAI-LIFECYCLE:file-14:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/GoogleCloudPlatform/python-docs-samples/dc0eecc2187791fea70f48201241288e25fc9620/documentai/snippets/undeploy_processor_version_sample_test.py"
license: "Apache-2.0"
sha256: "a29bf337d22a317e59bba27c7ce3d1365412bdffd77b91cda6192939e2fe9f16"
variables: []
secrets_allowed: false
```
````python
# Copyright 2022 Google LLC
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
from unittest import mock

from documentai.snippets import undeploy_processor_version_sample

location = "us"
project_id = os.environ["GOOGLE_CLOUD_PROJECT"]
processor_id = "aaaaaaaaa"
processor_version_id = "xxxxxxxxxx"


# TODO: Switch to Real Endpoint when Deployable Versions are Available
@mock.patch(
    "google.cloud.documentai.DocumentProcessorServiceClient.undeploy_processor_version"
)
@mock.patch("google.api_core.operation.Operation")
def test_undeploy_processor_version(
    operation_mock, undeploy_processor_version_mock, capsys
):
    undeploy_processor_version_mock.return_value = operation_mock

    undeploy_processor_version_sample.undeploy_processor_version_sample(
        project_id=project_id,
        location=location,
        processor_id=processor_id,
        processor_version_id=processor_version_id,
    )

    undeploy_processor_version_mock.assert_called_once()

    out, _ = capsys.readouterr()

    assert "operation" in out
````

### FILE: `google_document_ai_official_lifecycle/upstream/google/documentai/snippets/undeploy_processor_version_sample.py`
```yaml
block_id: "GOOGLE-DOCAI-LIFECYCLE:file-15:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/GoogleCloudPlatform/python-docs-samples/dc0eecc2187791fea70f48201241288e25fc9620/documentai/snippets/undeploy_processor_version_sample.py"
license: "Apache-2.0"
sha256: "8163efdc6d0dfb8e0fd48a2ca2ac1a64a9e69e9c9d3334da3ac32f81553998e8"
variables: []
secrets_allowed: false
```
````python
# Copyright 2022 Google LLC
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

# [START documentai_undeploy_processor_version]

from google.api_core.client_options import ClientOptions
from google.api_core.exceptions import FailedPrecondition
from google.api_core.exceptions import InvalidArgument
from google.cloud import documentai  # type: ignore

# TODO(developer): Uncomment these variables before running the sample.
# project_id = 'YOUR_PROJECT_ID'
# location = 'YOUR_PROCESSOR_LOCATION' # Format is 'us' or 'eu'
# processor_id = 'YOUR_PROCESSOR_ID' # Create processor before running sample
# processor_version_id = 'YOUR_PROCESSOR_VERSION_ID'


def undeploy_processor_version_sample(
    project_id: str, location: str, processor_id: str, processor_version_id: str
) -> None:
    # You must set the api_endpoint if you use a location other than 'us'.
    opts = ClientOptions(api_endpoint=f"{location}-documentai.googleapis.com")

    client = documentai.DocumentProcessorServiceClient(client_options=opts)

    # The full resource name of the processor version
    # e.g.: projects/project_id/locations/location/processors/processor_id/processorVersions/processor_version_id
    name = client.processor_version_path(
        project_id, location, processor_id, processor_version_id
    )

    # Make UndeployProcessorVersion request
    try:
        operation = client.undeploy_processor_version(name=name)
        # Print operation details
        print(operation.operation.name)
        # Wait for operation to complete
        operation.result()
    # Undeploy request will fail if the
    # processor version is already undeployed
    # or if a request is made on a pretrained processor version
    except (FailedPrecondition, InvalidArgument) as e:
        print(e.message)


# [END documentai_undeploy_processor_version]
````

### FILE: `google_document_ai_official_lifecycle/upstream/google/documentai/snippets/batch_process_documents_sample.py`
```yaml
block_id: "GOOGLE-DOCAI-LIFECYCLE:file-18:v1"
operation: CREATE
provenance: VERBATIM
source: "https://github.com/GoogleCloudPlatform/python-docs-samples/blob/dc0eecc2187791fea70f48201241288e25fc9620/documentai/snippets/batch_process_documents_sample.py; git blob 6d15d8713a7aa96d6012f8b36e2fea2333ac94d0"
license: "Apache-2.0"
sha256: "3ccd1f828172c76f9268451091a90225dbd173cd7979755d5a66adebcea41759"
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
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# [START documentai_batch_process_document]
import re
from typing import Optional

from google.api_core.client_options import ClientOptions
from google.api_core.exceptions import InternalServerError
from google.api_core.exceptions import RetryError
from google.cloud import documentai  # type: ignore
from google.cloud import storage

# TODO(developer): Uncomment these variables before running the sample.
# project_id = "YOUR_PROJECT_ID"
# location = "YOUR_PROCESSOR_LOCATION" # Format is "us" or "eu"
# processor_id = "YOUR_PROCESSOR_ID" # Create processor before running sample
# gcs_output_uri = "YOUR_OUTPUT_URI" # Must end with a trailing slash `/`. Format: gs://bucket/directory/subdirectory/
# processor_version_id = "YOUR_PROCESSOR_VERSION_ID" # Optional. Example: pretrained-ocr-v1.0-2020-09-23

# TODO(developer): You must specify either `gcs_input_uri` and `mime_type` or `gcs_input_prefix`
# gcs_input_uri = "YOUR_INPUT_URI" # Format: gs://bucket/directory/file.pdf
# input_mime_type = "application/pdf"
# gcs_input_prefix = "YOUR_INPUT_URI_PREFIX" # Format: gs://bucket/directory/
# field_mask = "text,entities,pages.pageNumber"  # Optional. The fields to return in the Document object.


def batch_process_documents(
    project_id: str,
    location: str,
    processor_id: str,
    gcs_output_uri: str,
    processor_version_id: Optional[str] = None,
    gcs_input_uri: Optional[str] = None,
    input_mime_type: Optional[str] = None,
    gcs_input_prefix: Optional[str] = None,
    field_mask: Optional[str] = None,
    timeout: int = 400,
) -> None:
    # You must set the `api_endpoint` if you use a location other than "us".
    opts = ClientOptions(api_endpoint=f"{location}-documentai.googleapis.com")

    client = documentai.DocumentProcessorServiceClient(client_options=opts)

    if gcs_input_uri:
        # Specify specific GCS URIs to process individual documents
        gcs_document = documentai.GcsDocument(
            gcs_uri=gcs_input_uri, mime_type=input_mime_type
        )
        # Load GCS Input URI into a List of document files
        gcs_documents = documentai.GcsDocuments(documents=[gcs_document])
        input_config = documentai.BatchDocumentsInputConfig(gcs_documents=gcs_documents)
    else:
        # Specify a GCS URI Prefix to process an entire directory
        gcs_prefix = documentai.GcsPrefix(gcs_uri_prefix=gcs_input_prefix)
        input_config = documentai.BatchDocumentsInputConfig(gcs_prefix=gcs_prefix)

    # Cloud Storage URI for the Output Directory
    gcs_output_config = documentai.DocumentOutputConfig.GcsOutputConfig(
        gcs_uri=gcs_output_uri, field_mask=field_mask
    )

    # Where to write results
    output_config = documentai.DocumentOutputConfig(gcs_output_config=gcs_output_config)

    if processor_version_id:
        # The full resource name of the processor version, e.g.:
        # projects/{project_id}/locations/{location}/processors/{processor_id}/processorVersions/{processor_version_id}
        name = client.processor_version_path(
            project_id, location, processor_id, processor_version_id
        )
    else:
        # The full resource name of the processor, e.g.:
        # projects/{project_id}/locations/{location}/processors/{processor_id}
        name = client.processor_path(project_id, location, processor_id)

    request = documentai.BatchProcessRequest(
        name=name,
        input_documents=input_config,
        document_output_config=output_config,
    )

    # BatchProcess returns a Long Running Operation (LRO)
    operation = client.batch_process_documents(request)

    # Continually polls the operation until it is complete.
    # This could take some time for larger files
    # Format: projects/{project_id}/locations/{location}/operations/{operation_id}
    try:
        print(f"Waiting for operation {operation.operation.name} to complete...")
        operation.result(timeout=timeout)
    # Catch exception when operation doesn't finish before timeout
    except (RetryError, InternalServerError) as e:
        print(e.message)

    # NOTE: Can also use callbacks for asynchronous processing
    #
    # def my_callback(future):
    #   result = future.result()
    #
    # operation.add_done_callback(my_callback)

    # After the operation is complete,
    # get output document information from operation metadata
    metadata = documentai.BatchProcessMetadata(operation.metadata)

    if metadata.state != documentai.BatchProcessMetadata.State.SUCCEEDED:
        raise ValueError(f"Batch Process Failed: {metadata.state_message}")

    storage_client = storage.Client()

    print("Output files:")
    # One process per Input Document
    for process in list(metadata.individual_process_statuses):
        # output_gcs_destination format: gs://BUCKET/PREFIX/OPERATION_NUMBER/INPUT_FILE_NUMBER/
        # The Cloud Storage API requires the bucket name and URI prefix separately
        matches = re.match(r"gs://(.*?)/(.*)", process.output_gcs_destination)
        if not matches:
            print(
                "Could not parse output GCS destination:",
                process.output_gcs_destination,
            )
            continue

        output_bucket, output_prefix = matches.groups()

        # Get List of Document Objects from the Output Bucket
        output_blobs = storage_client.list_blobs(output_bucket, prefix=output_prefix)

        # Document AI may output multiple JSON files per source file
        for blob in output_blobs:
            # Document AI should only output JSON files to GCS
            if blob.content_type != "application/json":
                print(
                    f"Skipping non-supported file: {blob.name} - Mimetype: {blob.content_type}"
                )
                continue

            # Download JSON File as bytes object and convert to Document Object
            print(f"Fetching {blob.name}")
            document = documentai.Document.from_json(
                blob.download_as_bytes(), ignore_unknown_fields=True
            )

            # For a full list of Document object attributes, please reference this page:
            # https://cloud.google.com/python/docs/reference/documentai/latest/google.cloud.documentai_v1.types.Document

            # Read the text recognition output from the processor
            print("The document contains the following text:")
            print(document.text)


# [END documentai_batch_process_document]
````

### FILE: `google_document_ai_official_lifecycle/upstream/google/documentai/snippets/batch_process_documents_sample_test.py`
```yaml
block_id: "GOOGLE-DOCAI-LIFECYCLE:file-19:v1"
operation: CREATE
provenance: VERBATIM
source: "https://github.com/GoogleCloudPlatform/python-docs-samples/blob/dc0eecc2187791fea70f48201241288e25fc9620/documentai/snippets/batch_process_documents_sample_test.py; git blob 168fe3ac3ee003580a63d2eb4c7aa1eb26eb6253"
license: "Apache-2.0"
sha256: "674a525239061de0a8d4c7075e846f69d91713aff27940f70d3f68ca3a19f6ec"
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

# flake8: noqa

import os
from uuid import uuid4

from documentai.snippets import batch_process_documents_sample

location = "us"
project_id = os.environ["GOOGLE_CLOUD_PROJECT"]
processor_id = "90484cfdedb024f6"
processor_version_id = "pretrained-form-parser-v1.0-2020-09-23"
gcs_input_uri = "gs://cloud-samples-data/documentai/invoice.pdf"
gcs_input_prefix = "gs://cloud-samples-data/documentai/workflows/"
input_mime_type = "application/pdf"
gcs_output_uri = f"gs://document-ai-python/{uuid4()}/"
field_mask = "text,pages.pageNumber"


def test_batch_process_documents(capsys):
    batch_process_documents_sample.batch_process_documents(
        project_id=project_id,
        location=location,
        processor_id=processor_id,
        gcs_input_uri=gcs_input_uri,
        input_mime_type=input_mime_type,
        gcs_output_uri=gcs_output_uri,
        field_mask=field_mask,
    )
    out, _ = capsys.readouterr()

    assert "operation" in out
    assert "Fetching" in out
    assert "text:" in out


def test_batch_process_documents_processor_version(capsys):
    batch_process_documents_sample.batch_process_documents(
        project_id=project_id,
        location=location,
        processor_id=processor_id,
        processor_version_id=processor_version_id,
        gcs_input_uri=gcs_input_uri,
        input_mime_type=input_mime_type,
        gcs_output_uri=gcs_output_uri,
        field_mask=field_mask,
    )
    out, _ = capsys.readouterr()

    assert "operation" in out
    assert "Fetching" in out
    assert "text:" in out


def test_batch_process_documents_gcs_prefix(capsys):
    batch_process_documents_sample.batch_process_documents(
        project_id=project_id,
        location=location,
        processor_id=processor_id,
        gcs_input_prefix=gcs_input_prefix,
        gcs_output_uri=gcs_output_uri,
        field_mask=field_mask,
    )
    out, _ = capsys.readouterr()

    assert "operation" in out
    assert "Fetching" in out
    assert "text:" in out
````

### FILE: `google_document_ai_official_lifecycle/upstream/google/documentai/snippets/handle_response_sample.py`
```yaml
block_id: "GOOGLE-DOCAI-LIFECYCLE:file-20:v1"
operation: CREATE
provenance: VERBATIM
source: "https://github.com/GoogleCloudPlatform/python-docs-samples/blob/dc0eecc2187791fea70f48201241288e25fc9620/documentai/snippets/handle_response_sample.py; git blob 58bbb1debe09c6fad3814789f3db78af4679e061"
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

### FILE: `google_document_ai_official_lifecycle/upstream/google/documentai/snippets/handle_response_sample_test.py`
```yaml
block_id: "GOOGLE-DOCAI-LIFECYCLE:file-21:v1"
operation: CREATE
provenance: VERBATIM
source: "https://github.com/GoogleCloudPlatform/python-docs-samples/blob/dc0eecc2187791fea70f48201241288e25fc9620/documentai/snippets/handle_response_sample_test.py; git blob b7c65834ccafbcb8fbae6f46807d0a7184ff1048"
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

### FILE: `google_document_ai_official_lifecycle/upstream/google/LICENSE.txt`
```yaml
block_id: "GOOGLE-DOCAI-LIFECYCLE:file-16:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/GoogleCloudPlatform/python-docs-samples/dc0eecc2187791fea70f48201241288e25fc9620/LICENSE"
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

### FILE: `google_document_ai_official_lifecycle/verify_official_lifecycle.py`
```yaml
block_id: "GOOGLE-DOCAI-LIFECYCLE:file-17:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace lifecycle guide/lock/verifier constrained by exact Google upstream"
license: "LicenseRef-Workspace-Owner"
sha256: "9a50f81fe795194f0c411a7301549d0522a6d3a28765baced5cce306942f6900"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import ast
import hashlib
import importlib.metadata as metadata
import importlib.util
import json
import py_compile
import sys
import tempfile
from pathlib import Path


ROOT = Path(__file__).resolve().parent
UPSTREAM = ROOT / "upstream" / "google"
sys.dont_write_bytecode = True


def sha256(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def load_module(path: Path):
    spec = importlib.util.spec_from_file_location(path.stem, path)
    if spec is None or spec.loader is None:
        raise AssertionError(f"cannot load {path}")
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def called_attributes(path: Path) -> set[str]:
    tree = ast.parse(path.read_text(encoding="utf-8"), filename=str(path))
    return {
        node.func.attr
        for node in ast.walk(tree)
        if isinstance(node, ast.Call) and isinstance(node.func, ast.Attribute)
    }


def defined_functions(path: Path) -> set[str]:
    tree = ast.parse(path.read_text(encoding="utf-8"), filename=str(path))
    return {node.name for node in tree.body if isinstance(node, ast.FunctionDef)}


def main() -> None:
    static_only = "--static-only" in sys.argv[1:]
    lock = json.loads((ROOT / "source-lock.json").read_text(encoding="utf-8"))
    expected: list[tuple[Path, int, str]] = []
    for item in lock["python_docs_samples"]["files"]:
        expected.append(
            (
                UPSTREAM / item["path"],
                item["bytes"],
                item["sha256"],
            )
        )
    docs = lock["google_cloud_documentation"]
    docs_path = UPSTREAM / "cloud-docs" / docs["extracted_path"]
    expected.append((docs_path, docs["bytes"], docs["sha256"]))
    license_item = lock["license_file"]
    expected.append((UPSTREAM / license_item["path"], license_item["bytes"], license_item["sha256"]))

    for path, byte_count, digest in expected:
        assert path.is_file(), f"missing upstream file: {path}"
        assert path.stat().st_size == byte_count, f"byte mismatch: {path}"
        assert sha256(path) == digest, f"SHA-256 mismatch: {path}"

    python_files = [path for path, _, _ in expected if path.suffix == ".py"]
    with tempfile.TemporaryDirectory(prefix="google-docai-lifecycle-pyc-") as tmp:
        for index, path in enumerate(python_files):
            py_compile.compile(
                str(path), cfile=str(Path(tmp) / f"{index}.pyc"), doraise=True
            )

    required_calls = {
        "create_processor_sample.py": {"create_processor"},
        "train_processor_version_sample.py": {"train_processor_version"},
        "evaluate_processor_version_sample.py": {"evaluate_processor_version"},
        "deploy_processor_version_sample.py": {"deploy_processor_version"},
        "set_default_processor_version_sample.py": {"set_default_processor_version"},
        "undeploy_processor_version_sample.py": {"undeploy_processor_version"},
        "batch_process_documents_sample.py": {
            "batch_process_documents",
            "list_blobs",
            "download_as_bytes",
            "from_json",
        },
    }
    sample_root = UPSTREAM / "documentai" / "snippets"
    for filename, required in required_calls.items():
        missing = required - called_attributes(sample_root / filename)
        assert not missing, f"{filename} missing API calls: {sorted(missing)}"

    required_docs_functions = {
        "sample_create_processor",
        "initialize_dataset",
        "get_dataset_schema",
        "upload_dataset_schema",
        "list_documents",
        "get_document",
        "import_documents",
        "import_processor_version",
        "deploy_and_set_default_processor_version",
        "main",
    }
    assert required_docs_functions <= defined_functions(docs_path)
    assert {
        "update_dataset",
        "update_dataset_schema",
        "import_documents",
        "import_processor_version",
        "deploy_processor_version",
        "set_default_processor_version",
    } <= called_attributes(docs_path)

    response_path = sample_root / "handle_response_sample.py"
    required_response_functions = {
        "process_document_ocr_sample",
        "process_document_form_sample",
        "process_document_entity_extraction_sample",
        "process_document_splitter_sample",
        "process_document_layout_sample",
        "process_document_custom_extractor_sample",
        "process_document",
        "print_entity",
        "layout_to_text",
        "page_refs_to_string",
    }
    assert required_response_functions <= defined_functions(response_path)
    assert {
        "processor_version_path",
        "process_document",
    } <= called_attributes(response_path)

    if static_only:
        print(
            "GOOGLE_DOCUMENT_AI_OFFICIAL_LIFECYCLE_OFFLINE_PASS "
            f"upstream_files={len(expected)} python_files={len(python_files)} "
            f"functions={len(required_docs_functions)} static_only=1"
        )
        return

    assert metadata.version("google-cloud-documentai") == "3.15.0"
    assert metadata.version("google-cloud-storage") == "3.13.1"

    set_default = load_module(sample_root / "set_default_processor_version_sample.py")
    response = load_module(response_path)

    class Segment:
        def __init__(self, start: int, end: int) -> None:
            self.start_index = start
            self.end_index = end

    class Anchor:
        text_segments = [Segment(0, 5), Segment(6, 10)]

    class Layout:
        text_anchor = Anchor()

    class PageRef:
        def __init__(self, page: int) -> None:
            self.page = page

    assert response.layout_to_text(Layout(), "alpha beta") == "alphabeta"
    assert response.page_refs_to_string([PageRef(0)]) == "page 1 is"
    assert response.page_refs_to_string([PageRef(0), PageRef(2)]) == "pages 1, 3 are"

    class FakeOperation:
        class OperationName:
            name = "operations/offline-set-default"

        operation = OperationName()

        def __init__(self) -> None:
            self.results = 0

        def result(self) -> None:
            self.results += 1

    class FakeClient:
        instances: list["FakeClient"] = []

        def __init__(self, client_options) -> None:
            self.endpoint = client_options.api_endpoint
            self.requests = []
            self.operation = FakeOperation()
            self.__class__.instances.append(self)

        @staticmethod
        def processor_path(project: str, location: str, processor: str) -> str:
            return f"projects/{project}/locations/{location}/processors/{processor}"

        @staticmethod
        def processor_version_path(
            project: str, location: str, processor: str, version: str
        ) -> str:
            return (
                f"projects/{project}/locations/{location}/processors/{processor}"
                f"/processorVersions/{version}"
            )

        def set_default_processor_version(self, request):
            self.requests.append(request)
            return self.operation

    original_client = set_default.documentai.DocumentProcessorServiceClient
    set_default.documentai.DocumentProcessorServiceClient = FakeClient
    try:
        set_default.set_default_processor_version_sample(
            "offline-project", "eu", "processor-1", "version-2"
        )
    finally:
        set_default.documentai.DocumentProcessorServiceClient = original_client

    assert len(FakeClient.instances) == 1
    captured = FakeClient.instances[0]
    assert captured.endpoint == "eu-documentai.googleapis.com"
    assert len(captured.requests) == 1
    request = captured.requests[0]
    assert request.processor == "projects/offline-project/locations/eu/processors/processor-1"
    assert request.default_processor_version.endswith("/processorVersions/version-2")
    assert captured.operation.results == 1

    print(
        "GOOGLE_DOCUMENT_AI_OFFICIAL_LIFECYCLE_OFFLINE_PASS "
        f"upstream_files={len(expected)} python_files={len(python_files)} "
        f"functions={len(required_docs_functions)}"
    )


if __name__ == "__main__":
    main()
````

## 6. Configuration surface

Configuración obligatoria previa a live: `project_id`, `location`, processor type/ID/version, display names, train/test/unassigned GCS prefixes, destination project/processor si hay migración, identidad ADC, IAM mínimo, storage/CMEK/retención, cuotas, presupuesto y umbrales de evaluación del proyecto. Ningún secreto se materializa y ninguna variable se supone.

## 7. Dependency bill

Runtime: Python 3.10+ y artefactos oficiales `google-cloud-documentai==3.15.0` y `google-cloud-storage==3.13.1` fijados por `OFFICIAL-DOCUMENT-SDK-ARTIFACT-CORE`; la migración documental usa además `tqdm` declarado por la página oficial y debe quedar fijado por el proyecto antes de live. Verificación upstream: `pytest==8.4.1`. El pack no instala ni actualiza dependencias automáticamente.

## 8. Apply order

1. materializar artifact core, secure ingestion, process/custom-extractor, este lifecycle y strict evaluation;
2. ejecutar `verify_official_lifecycle.py --static-only`; repetir sin flag con el SDK exacto;
3. ejecutar los cinco tests Google offline-safe; mantener set-default, batch y response live bloqueados;
4. completar intake/approvals y generar configuración de proyecto sin modificar los archivos `VERBATIM`;
5. crear processor/dataset, importar ground truth revisado, entrenar, evaluar y promover sólo si alcanza los umbrales;
6. desplegar/cambiar default con receipt y conservar undeploy/versión anterior como rollback.

## 9. Verification

PASS observado 2026-08-28: 21/21 archivos allowlisted, 18 upstream (17 Python + licencia), todos los SHA/bytes exactos; 17/17 Python compilan; arnés estático emite `GOOGLE_DOCUMENT_AI_OFFICIAL_LIFECYCLE_OFFLINE_PASS upstream_files=18 python_files=17 functions=10 static_only=1`. En venv aislado con Document AI 3.15.0 + Storage 3.13.1, `pip check` y el arnés completo pasan; cinco tests oficiales Google offline-safe también pasan. Set-default produjo 403 con proyecto ficticio y batch/response usan proyecto, processor, GCS y fixtures live; se preservan sin presentarlos como PASS offline. La verificación live exige cuenta/ADC/IAM/storage/cuota/costo autorizados y métricas contra corpus target.

## 10. Reconstruction evidence

La reconstrucción se gobierna por `source-lock.json`, commit firmado `dc0eecc...`, hash del HTML oficial `21c4fd07...`, hash del bloque extraído `10897a36...`, licencia Apache-2.0, materialización 21/21, comparación SHA, compilación, AST, helpers de respuesta, cinco tests Google offline-safe y el ledger de fallos 941–968. Una actualización requiere commit/page/SDK nuevos, diff, hashes, SCA, mismas regresiones y canary/rollback live; nunca se hereda el PASS.
