# Official Document SDK Artifact Lock

Fecha de corte y verificación: 2026-08-27.

Este lock complementa —no reemplaza— `OFFICIAL_UPSTREAM_ACQUISITION_CORE.md`. Fija wheels productivos y source distributions oficiales de Content Understanding, Document Intelligence, Document AI y OpenAI que un proyecto Python puede adquirir inmediatamente. Los accelerators permanecen separados; los samples publicados dentro de cada sdist conservan su clasificación upstream.

## Artefactos directos

| ID | Propietario | Package/version | Python | Bytes | SHA-256 | Licencia declarada |
|---|---|---|---|---:|---|---|
| `azure-ai-documentintelligence-1.0.2-wheel` | Microsoft Azure | `azure-ai-documentintelligence==1.0.2` | `>=3.8` | 106,005 | `e1fb446abbdeccc9759d897898a0fe13141ed29f9ad11fc705f951925822ed59` | MIT según distribución/package metadata |
| `azure-ai-contentunderstanding-1.1.0-wheel` | Microsoft Azure | `azure-ai-contentunderstanding==1.1.0` | `>=3.9` | 101,987 | `d1d6bdeffe02f5c8cc5309f0e120a1338f51af04638605f9617b7b501e2d1dd6` | MIT, `License-Expression` del wheel |
| `google-cloud-documentai-3.15.0-wheel` | Google Cloud | `google-cloud-documentai==3.15.0` | `>=3.10` | 310,694 | `f040f4f9db43411184197a808b11fde52b580723a9fca336f43ea7c6a885bfd8` | Apache-2.0 según distribución/package metadata |
| `openai-3.3.1-wheel` | OpenAI | `openai==3.3.1` | `>=3.10` | 1,690,337 | `9652df7fdf8ee6f5bd58e0a12f2b1d414a18e0f06bb7a9a57c8643a5f5469bd3` | Apache-2.0, `License-Expression` del wheel |

## Source distributions oficiales exactos

Los source distributions correspondientes a cuatro SDKs estables quedaron descargados y verificados:

| ID | Artefacto | Bytes | SHA-256 | Contenido observado | Publicación |
|---|---|---:|---|---|---|
| `azure-ai-documentintelligence-1.0.2-sdist` | `azure_ai_documentintelligence-1.0.2.tar.gz` | 170,940 | `4d75a2513f2839365ebabc0e0e1772f5601b3a8c9a71e75da12440da13b63484` | 131 entries; 103 Python; 56 samples; 20 tests; cuatro samples invoice síncronos/asíncronos desde URL/bytes | PyPI Microsoft, 2025-03-27 UTC |
| `azure-ai-contentunderstanding-1.1.0-sdist` | `azure_ai_contentunderstanding-1.1.0.tar.gz` | 230,330 | `00f39c7cf2ba50c8d4586315f4295a45fad0881daaff07d13d33aea2bafc7ec8` | 138 entries; 105 Python; 68 samples; 43 tests; runtime sync/async, models, LRO, analyzers y diagnostics | PyPI Microsoft, 2026-04-20 |
| `google-cloud-documentai-3.15.0-sdist` | `google_cloud_documentai-3.15.0.tar.gz` | 362,736 | `d50b69a8a62aaf803b0d926f93b6579df440f5f779f45fec2c5972768a1cdbbd` | 109 entries; 71 Python; 8 tests; cliente y tipos generados oficiales | PyPI Google, 2026-06-03 UTC |
| `openai-3.3.1-sdist` | `openai-3.3.1.tar.gz` | 1,282,113 | `6f22807de1a976c932cecda620e8172a8c3fdbaeed29c7f21564e0c2410edf56` | 1,821 entries; 1,760 Python; 189 tests; SDK completo y tipos Responses/file input | PyPI OpenAI, 2026-08-19 UTC |

```text
https://files.pythonhosted.org/packages/b6/5a/6dbcb8278c8d3779fc03992d8a9f87a9ce5280437d76e72f8121f9a90d46/azure_ai_contentunderstanding-1.1.0.tar.gz
https://files.pythonhosted.org/packages/44/7b/8115cd713e2caa5e44def85f2b7ebd02a74ae74d7113ba20bdd41fd6dd80/azure_ai_documentintelligence-1.0.2.tar.gz
https://files.pythonhosted.org/packages/e0/63/a775dd454e0fa523b999780665f3117c0026eaa604dd7c47ad3c45c7624f/google_cloud_documentai-3.15.0.tar.gz
https://files.pythonhosted.org/packages/7d/9c/ba0c292b4032ede74c249ca314ad64eb1bb5a03a843f6e01facb02f80cd8/openai-3.3.1.tar.gz
```

Los cuatro tarballs pasaron path safety sin rutas absolutas, traversal ni symlinks; sus 1.934 archivos Python compilaron sin error. Azure DI 1.0.2 se construyó e instaló desde el sdist exacto en un venv nuevo, `pip check` pasó y los símbolos `DocumentIntelligenceClient`, `AnalyzeDocumentRequest` y `AnalyzeResult` importaron con versión 1.0.2. Es código oficial real, no reconstrucción Elite. PyPI no reporta Trusted Publishing para todos los uploads; la admisión se apoya en URL inmutable, bytes, hash y metadata del owner sin inventar attestation.

URLs inmutables:

```text
https://files.pythonhosted.org/packages/d9/75/c9ec040f23082f54ffb1977ff8f364c2d21c79a640a13d1c1809e7fd6b1a/azure_ai_documentintelligence-1.0.2-py3-none-any.whl
https://files.pythonhosted.org/packages/67/d6/6f60ac39b73a0b8a2bc86cd6895e26eccc20d9b7bde86814fe6a36964dff/azure_ai_contentunderstanding-1.1.0-py3-none-any.whl
https://files.pythonhosted.org/packages/55/d1/2a873f97cb08bb592f5b944e39f040161d5f7d2d4edbfae4c7515869b12a/google_cloud_documentai-3.15.0-py3-none-any.whl
https://files.pythonhosted.org/packages/6a/db/2b7a1b3de659bb82aef979116c74e809982b13e42c057759767552b5155f/openai-3.3.1-py3-none-any.whl
https://files.pythonhosted.org/packages/b6/5a/6dbcb8278c8d3779fc03992d8a9f87a9ce5280437d76e72f8121f9a90d46/azure_ai_contentunderstanding-1.1.0.tar.gz
https://files.pythonhosted.org/packages/44/7b/8115cd713e2caa5e44def85f2b7ebd02a74ae74d7113ba20bdd41fd6dd80/azure_ai_documentintelligence-1.0.2.tar.gz
https://files.pythonhosted.org/packages/e0/63/a775dd454e0fa523b999780665f3117c0026eaa604dd7c47ad3c45c7624f/google_cloud_documentai-3.15.0.tar.gz
https://files.pythonhosted.org/packages/7d/9c/ba0c292b4032ede74c249ca314ad64eb1bb5a03a843f6e01facb02f80cd8/openai-3.3.1.tar.gz
```

Google PyPI registra Trusted Publishing y una attestation de publicación para `3.15.0`. Azure PyPI no expuso una attestation equivalente en la consulta auditada; el wheel queda fijado por digest, no por una afirmación de provenance ausente.

## Instalación directa de prueba

En un virtual environment vacío:

```powershell
python -m pip install `
  'https://files.pythonhosted.org/packages/d9/75/c9ec040f23082f54ffb1977ff8f364c2d21c79a640a13d1c1809e7fd6b1a/azure_ai_documentintelligence-1.0.2-py3-none-any.whl#sha256=e1fb446abbdeccc9759d897898a0fe13141ed29f9ad11fc705f951925822ed59' `
  'https://files.pythonhosted.org/packages/67/d6/6f60ac39b73a0b8a2bc86cd6895e26eccc20d9b7bde86814fe6a36964dff/azure_ai_contentunderstanding-1.1.0-py3-none-any.whl#sha256=d1d6bdeffe02f5c8cc5309f0e120a1338f51af04638605f9617b7b501e2d1dd6' `
  'https://files.pythonhosted.org/packages/55/d1/2a873f97cb08bb592f5b944e39f040161d5f7d2d4edbfae4c7515869b12a/google_cloud_documentai-3.15.0-py3-none-any.whl#sha256=f040f4f9db43411184197a808b11fde52b580723a9fca336f43ea7c6a885bfd8' `
  'https://files.pythonhosted.org/packages/6a/db/2b7a1b3de659bb82aef979116c74e809982b13e42c057759767552b5155f/openai-3.3.1-py3-none-any.whl#sha256=9652df7fdf8ee6f5bd58e0a12f2b1d414a18e0f06bb7a9a57c8643a5f5469bd3'
```

Verificación mínima ejecutada:

```python
import importlib.metadata as metadata
from azure.ai.documentintelligence import DocumentIntelligenceClient
from azure.ai.contentunderstanding import ContentUnderstandingClient
from azure.ai.contentunderstanding.models import AnalysisInput, ContentAnalyzer
from google.cloud import documentai_v1
from openai import OpenAI, AsyncOpenAI
from openai.types.responses import ResponseInputFileParam

assert metadata.version("azure-ai-documentintelligence") == "1.0.2"
assert metadata.version("azure-ai-contentunderstanding") == "1.1.0"
assert metadata.version("google-cloud-documentai") == "3.15.0"
assert DocumentIntelligenceClient.__name__ == "DocumentIntelligenceClient"
assert ContentUnderstandingClient.__name__ == "ContentUnderstandingClient"
assert AnalysisInput.__name__ == "AnalysisInput"
assert ContentAnalyzer.__name__ == "ContentAnalyzer"
assert documentai_v1.DocumentProcessorServiceClient.__name__ == "DocumentProcessorServiceClient"
assert metadata.version("openai") == "3.3.1"
assert OpenAI.__name__ == "OpenAI"
assert AsyncOpenAI.__name__ == "AsyncOpenAI"
assert ResponseInputFileParam.__name__ == "ResponseInputFileParam"
```

Resultado observado: los cuatro imports y versiones pasaron en Python 3.14 aislado. El SDK Content Understanding expuso 70 símbolos públicos de modelo en el probe. El extra opcional `[aio]` no existe en `1.1.0`; no debe agregarse al requisito. La operación real sigue requiriendo la cuenta/endpoint/identidad y recursos de cada proveedor.

## Grafo resuelto observado, no lock productivo

La instalación de prueba resolvió:

```text
azure-ai-contentunderstanding==1.1.0
azure-core==1.41.0
certifi==2026.7.22
cffi==2.1.1
charset-normalizer==3.5.1
cryptography==50.0.1
google-api-core==2.34.0
google-auth==2.57.0
googleapis-common-protos==1.75.2
grpcio==1.83.0
grpcio-status==1.83.0
idna==3.19
isodate==0.7.2
proto-plus==1.28.4
protobuf==7.36.0
pyasn1==0.6.4
pyasn1-modules==0.4.2
pycparser==3.0
requests==2.34.2
typing-extensions==4.16.0
urllib3==2.7.0
```

Este listado registra el probe, pero no autoriza un deployment: faltan hashes/markers de todas las transitivas por OS/Python. El proyecto debe generar un lock completo con hashes, SBOM, license report y vulnerability gate para su plataforma exacta.

## Admisión

- clasificación: `INTEGRATION_ONLY`;
- implementación inmediata demostrada: `MICROSOFT_AZURE_DOCUMENT_INTELLIGENCE_OFFICIAL_INVOICE_SAMPLE.md` materializa byte-a-byte el sample sync `prebuilt-invoice` de Azure DI 1.0.2 y su licencia, adquiere por hash el fixture oficial del mismo commit y ejecuta la función upstream offline; no se hereda este PASS a los demás samples;
- uso inmediato: crear cliente, tipos y requests contra un endpoint sandbox ya provisto;
- no incluye: cuenta, endpoint, identidad, cuota, processor/model, documentos, ground truth o exactitud;
- no mezclar Azure y Google en el mismo runtime salvo que el benchmark del proyecto lo justifique;
- no registrar API keys en Markdown, código o receipts;
- fijar endpoint/region, processor/model/version y response schema en `PROJECT_DOCUMENT_INTELLIGENCE_DECISION.md`;
- ejecutar contract probe y corpus evaluation antes de almacenamiento automático.

AWS usa el módulo Go `service/textract v1.45.0` dentro de `aws-sdk-go-v2-textract-1.45.0`; el source combinado `aws-sdk-go-v2-2026-08-20` permanece para otros adapters AWS en sus versiones fijadas. Oracle usa `oracle-oci-python-sdk-2.185.0`. Su release asset oficial contiene el wheel PyPI byte-idéntico SHA-256 `93556f08c6270e2e174d59d94fcd49fa794ba3f4d60630af173bc9b6204b085c`; el import offline y un lock de auditoría Python 3.14 con 23 paquetes/0 advisories pasan. Eso no convierte el `requirements.txt` upstream con rangos en lock productivo ni aporta un sample OCI ejecutable. Todos están fijados por source archive/commit en el pack oficial, por lo que no se duplican aquí como wheels.
