# Go AWS Enterprise Storage and Email Adapters — reconstruction evidence V1

## Alcance

Reconstrucción del pack `GO-AWS-ENTERPRISE-STORAGE-EMAIL-ADAPTERS` 0.1.0. Nueve archivos `AUTHORED` integran los SDK oficiales AWS; no copian código AWS ni realizaron efectos externos.

## Entorno y módulos

- Go oficial `go1.26.7 windows/amd64`;
- AWS config `v1.32.38`;
- AWS S3 `v1.107.3`;
- AWS SES v2 `v1.67.0`;
- source AWS fijado en commit `284a4846e7bb941926e2d15144b17c01ffc98c1d`, Apache-2.0 + NOTICE y archive SHA-256 registrado.

## Resultado limpio

```text
Materialized 9 files into a new directory
all modules verified
GOPROXY=off go test -count=1 ./...
ok example.com/elite/aws-enterprise-adapters/awsenterprise
github.com/aws/aws-sdk-go-v2/config v1.32.38
github.com/aws/aws-sdk-go-v2/service/s3 v1.107.3
github.com/aws/aws-sdk-go-v2/service/sesv2 v1.67.0
```

Cuatro tests cubren S3 `PutObjectInput` con SHA-256, `IfNoneMatch=*`, AES256/KMS, ETag/VersionId y mismatch; SES `SendEmailInput`, verified-address syntax candidate, outbox message key/tag, receipt sin PII, duplicados y response sin ID. `NewClients` usa `config.LoadDefaultConfig` con región explícita.

## Fallos retenidos

`LIB-FAIL-041` registra que el primer draft declaraba config 1.32.38 sin importarlo; tidy retiró el módulo y el probe falló. Al añadir la fábrica, una ejecución sin repin resolvió 1.32.39. Se fijó explícitamente 1.32.38 y se repitieron verify/test/list offline. `LIB-FAIL-032` recurrencia 5 conserva además el primer patch de inserción con contexto no vigente.

## Condiciones abiertas

Sin cuenta/región/IAM/costo aprobados, bucket/retention/malware/recovery, identity SES/configuration set/events/suppression ni outbox productivo. No hubo llamada sandbox. Sigue `CONDITIONED`, no `REUSABLE_PACK`.

## Fuentes oficiales

- AWS SDK Go v2: https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/welcome.html
- S3 checksums: https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/s3-checksums.html
- S3 PutObject examples: https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/go_s3_code_examples.html
- SES SendEmail: https://docs.aws.amazon.com/ses/latest/APIReference/API_SendEmail.html
