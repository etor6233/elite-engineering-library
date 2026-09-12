# Go AWS Textract Document Runtime — reconstruction evidence V1

## Alcance

Reconstrucción del pack `GO-AWS-TEXTRACT-DOCUMENT-RUNTIME` 0.1.0. El runtime y tests son glue `AUTHORED`; llaman los tipos/métodos públicos de módulos Apache-2.0 oficiales AWS y no copian su código fuente. No se realizó llamada AWS ni se afirmó precisión productiva.

## Entorno y dependencias oficiales

- fecha: 2026-08-26;
- Go oficial: `go1.26.7 windows/amd64`;
- `github.com/aws/aws-sdk-go-v2/config v1.32.38`;
- `github.com/aws/aws-sdk-go-v2/service/textract v1.44.2`;
- source AWS fijado en lock: commit `284a4846e7bb941926e2d15144b17c01ffc98c1d`, archive SHA-256 `b7622b29fa7fb852cc0ea1c8ae23f22be64874b8869e63c37efd7efb581d0747`, LICENSE/NOTICE verificados.

## Reconstrucción y pruebas

El materializador reconstruyó 7/7 archivos y verificó cada SHA-256. Después de una adquisición online exacta para poblar el module cache, la repetición limpia desde el Markdown se ejecutó con red de módulos desactivada:

```text
all modules verified
GOPROXY=off go test -count=1 ./...
? example.com/elite/aws-textract-document-runtime/cmd/analyze [no test files]
ok example.com/elite/aws-textract-document-runtime/textractruntime
github.com/aws/aws-sdk-go-v2/config v1.32.38
github.com/aws/aws-sdk-go-v2/service/textract v1.44.2
```

Cuatro tests cubren `AnalyzeExpenseInput`/`Document.Bytes`, output de expense, `AnalyzeDocumentInput`, FeatureTypes, QueriesConfig, AdaptersConfig, respuesta/receipt, SHA y commit atómico. Los negativos rechazan response vacío, Queries incompletas, adapter sin QUERIES, features en expense y extensión no soportada. El template mantiene proforma, purchase order, packing list, BOL, delivery note, certificados y aduana en `BLOCKED_*`, siempre `automatic_storage=false`.

## Condiciones abiertas

No existían cuenta AWS, IAM, presupuesto, endpoint/región autorizada, documentos representativos, ground truth, Queries aprobadas ni adapter entrenado. Faltan llamada sandbox, métricas por campo/clase, revisión, privacidad/residencia, costo/cuota, throttling/retry, load/concurrency, integración de storage validado y rollback real. La admisión permanece `CONDITIONED`.

## Fuentes oficiales

- AWS SDK for Go v2 Textract source: https://github.com/aws/aws-sdk-go-v2/tree/284a4846e7bb941926e2d15144b17c01ffc98c1d/service/textract
- AnalyzeExpense API: https://docs.aws.amazon.com/textract/latest/APIReference/API_AnalyzeExpense.html
- AnalyzeDocument API: https://docs.aws.amazon.com/textract/latest/APIReference/API_AnalyzeDocument.html
- Textract Queries adapters: https://docs.aws.amazon.com/textract/latest/dg/textract-using-adapters.html
