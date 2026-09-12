# Go Google Cloud Storage Evidence Adapter

## 1. Metadata

```yaml
pack_id: "GO-GOOGLE-CLOUD-STORAGE-ADAPTER"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa integración Go sobre el SDK oficial Google Cloud Storage 1.65.1 para creación única por generación, CRC32C, Cloud KMS, retención por objeto y verificación posterior de la generación exacta."
stacks: ["Go 1.26.7", "Google Cloud Storage Go 1.65.1", "Cloud KMS", "Object Retention Lock"]
compatible_with: ["GOOGLE-DOCUMENT-AI-RUNTIME 0.1.x", "SECURE-LOCAL-FILE-INGESTION-GATE 0.1.x", "STRICT-DOCUMENT-FIELD-EVALUATION-GATE 0.1.x"]
incompatible_with: ["credenciales en archivos", "overwrite silencioso", "bucket no aprobado", "retención no aprobada", "auto-storage documental", "habilitación o bloqueo automático de políticas"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0"
upstream_sources: ["https://github.com/googleapis/google-cloud-go/tree/3ab7d1390bbbb40cf197a545b941f9df760a3269/storage", "https://pkg.go.dev/cloud.google.com/go/storage@v1.65.1", "https://docs.cloud.google.com/storage/docs/object-lock", "https://cloud.google.com/storage/docs/generations-preconditions"]
verified_at: "2026-08-27"
```

## 2. Applicability

Use sólo cuando el blueprint seleccione Google Cloud Storage para evidencia documental y el usuario haya aprobado proyecto/cuenta, IAM, costo, bucket, Cloud KMS y política de retención. El código guarda bytes ya admitidos por seguridad y evaluación; no extrae, clasifica ni autoriza persistencia automática. Rechazar si falta el número de proyecto esperado, si el bucket no reporta `ObjectRetentionMode=Enabled`, si la retención o el costo no están aprobados, o si el proyecto necesita cambiar IAM/bucket/policy mediante esta operación.

El SDK es código oficial Google Apache-2.0. `gcsstorage` es glue `AUTHORED` local compilado contra el SDK: no se atribuye a Google. No existe claim de llamada cloud, corpus, IAM ni compliance real sin evidencia del proyecto.

## 3. Architecture contract

`Store.Put` valida localmente approvals, límites, nombres, KMS y retención. Consulta owner/configuración del bucket antes del efecto. `GoogleCloudBackend.Create` usa el SDK oficial con `storage.Conditions{DoesNotExist:true}`, CRC32C Castagnoli provisto por caller, `KMSKeyName` y `ObjectRetention`; GCS debe rechazar checksum distinto o nombre ya existente. Tras `Close`, el adapter consulta `Object(name).Generation(generation).Attrs` y compara generación, tamaño, CRC32C, KMS, modo/fecha y ETag. Sólo entonces emite `retention_verified=true`.

Si la creación pudo ocurrir y falla la verificación, devuelve `PartialVerificationError` junto con receipt `created=true`, generación y `retention_verified=false`: el caller reconcilia esa generación y no reintenta a ciegas. El receipt hashea bucket, object name y KMS resource; no contiene bytes, credenciales ni metadata libre. El adapter nunca crea/reconfigura buckets, habilita Object Retention Lock, bloquea políticas, cambia IAM, elimina objetos ni reduce retención. `Locked` requiere aprobación irreversible separada. Performance, quota, cost, VPC/service perimeter, logging, malware, recuperación y regulación pertenecen al target.

## 4. Exact file manifest

```text
CREATE google_cloud_storage_adapter/go.mod
CREATE google_cloud_storage_adapter/go.sum
CREATE google_cloud_storage_adapter/provider-profile.template.json
CREATE google_cloud_storage_adapter/official-artifact-lock.json
CREATE google_cloud_storage_adapter/gcsstorage/storage.go
CREATE google_cloud_storage_adapter/gcsstorage/google_backend.go
CREATE google_cloud_storage_adapter/gcsstorage/storage_test.go
CREATE google_cloud_storage_adapter/gcsstorage/crc_test.go
CREATE google_cloud_storage_adapter/README.md
```

## 5. Materialization blocks

### FILE: `google_cloud_storage_adapter/go.mod`
```yaml
block_id: "GO-GOOGLE-CLOUD-STORAGE-ADAPTER:gomod:v1"
operation: CREATE
provenance: AUTHORED
source: "exact official Google Cloud Storage Go module plus resolved graph"
license: "LicenseRef-Workspace-Owner"
sha256: "3a0d72638310d1ef66969dd6525df7b447815dc14b0c24833e66a7092582e09a"
variables: []
secrets_allowed: false
```
````text
module example.com/elite/google-cloud-storage-adapter

go 1.26.0

require cloud.google.com/go/storage v1.65.1

require (
	cel.dev/expr v0.25.1 // indirect
	cloud.google.com/go v0.123.0 // indirect
	cloud.google.com/go/auth v0.20.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.8 // indirect
	cloud.google.com/go/compute/metadata v0.9.0 // indirect
	cloud.google.com/go/iam v1.11.0 // indirect
	cloud.google.com/go/monitoring v1.29.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.32.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric v0.57.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.57.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cncf/xds/go v0.0.0-20260202195803-dba9d589def2 // indirect
	github.com/envoyproxy/go-control-plane/envoy v1.37.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.3.3 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-jose/go-jose/v4 v4.1.4 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.17 // indirect
	github.com/googleapis/gax-go/v2 v2.23.0 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/spiffe/go-spiffe/v2 v2.6.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/detectors/gcp v1.43.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.68.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.67.0 // indirect
	go.opentelemetry.io/otel v1.44.0 // indirect
	go.opentelemetry.io/otel/metric v1.44.0 // indirect
	go.opentelemetry.io/otel/sdk v1.44.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.44.0 // indirect
	go.opentelemetry.io/otel/trace v1.44.0 // indirect
	golang.org/x/crypto v0.53.0 // indirect
	golang.org/x/net v0.56.0 // indirect
	golang.org/x/oauth2 v0.36.0 // indirect
	golang.org/x/sync v0.21.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/text v0.39.0 // indirect
	golang.org/x/time v0.15.0 // indirect
	google.golang.org/api v0.287.1 // indirect
	google.golang.org/genproto v0.0.0-20260519071638-aa98bba5eb94 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260630182238-925bb5da69e7 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260630182238-925bb5da69e7 // indirect
	google.golang.org/grpc v1.82.1 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
````

### FILE: `google_cloud_storage_adapter/go.sum`
```yaml
block_id: "GO-GOOGLE-CLOUD-STORAGE-ADAPTER:gosum:v1"
operation: CREATE
provenance: AUTHORED
source: "Go checksum database entries from exact graph"
license: "LicenseRef-Workspace-Owner"
sha256: "374651b97a68537d2b0e43bd0747aec1974eec33ecb110b457378812b50dec70"
variables: []
secrets_allowed: false
```
````text
cel.dev/expr v0.25.1 h1:1KrZg61W6TWSxuNZ37Xy49ps13NUovb66QLprthtwi4=
cel.dev/expr v0.25.1/go.mod h1:hrXvqGP6G6gyx8UAHSHJ5RGk//1Oj5nXQ2NI02Nrsg4=
cloud.google.com/go v0.123.0 h1:2NAUJwPR47q+E35uaJeYoNhuNEM9kM8SjgRgdeOJUSE=
cloud.google.com/go v0.123.0/go.mod h1:xBoMV08QcqUGuPW65Qfm1o9Y4zKZBpGS+7bImXLTAZU=
cloud.google.com/go/auth v0.20.0 h1:kXTssoVb4azsVDoUiF8KvxAqrsQcQtB53DcSgta74CA=
cloud.google.com/go/auth v0.20.0/go.mod h1:942/yi/itH1SsmpyrbnTMDgGfdy2BUqIKyd0cyYLc5Q=
cloud.google.com/go/auth/oauth2adapt v0.2.8 h1:keo8NaayQZ6wimpNSmW5OPc283g65QNIiLpZnkHRbnc=
cloud.google.com/go/auth/oauth2adapt v0.2.8/go.mod h1:XQ9y31RkqZCcwJWNSx2Xvric3RrU88hAYYbjDWYDL+c=
cloud.google.com/go/compute/metadata v0.9.0 h1:pDUj4QMoPejqq20dK0Pg2N4yG9zIkYGdBtwLoEkH9Zs=
cloud.google.com/go/compute/metadata v0.9.0/go.mod h1:E0bWwX5wTnLPedCKqk3pJmVgCBSM6qQI1yTBdEb3C10=
cloud.google.com/go/iam v1.11.0 h1:KieQ9Pb+LLPak1O3Rv3GgCxhnmkYf7Xyh0P5HfF1jFM=
cloud.google.com/go/iam v1.11.0/go.mod h1:KP+nKGugNJW4LcLx1uEZcq1ok5sQHFaQehQNl4QDgV4=
cloud.google.com/go/logging v1.18.0 h1:KhzZq+1cSkPH9YUaKLLhLtQxIHitVayBmk0sGfoM9+k=
cloud.google.com/go/logging v1.18.0/go.mod h1:ZGKnpBaURITh+g/uom2VhbiFoFWvejcrHPDhxFtU/gI=
cloud.google.com/go/longrunning v1.2.0 h1:WjYH3YHBGCxGJP9M4dWGHBfXr/cFIjMkNgWcJj7/iMM=
cloud.google.com/go/longrunning v1.2.0/go.mod h1:5KMQALFGOCtFoi2xSOA1u3H7WKlhmckgiyFw7+LGQp0=
cloud.google.com/go/monitoring v1.29.0 h1:AHhDsFaSax1/4k+qlIDX/SDGe6hggnfXJ9dkgD9qBPY=
cloud.google.com/go/monitoring v1.29.0/go.mod h1:72NOVjJXHY/HBfoLT0+qlCZBT059+9VXLeAnL2PeeVM=
cloud.google.com/go/storage v1.65.1 h1:LRRpBJUTf+OXDPX9jZUKZ3mSLIsz3htG+qUpeNZovyA=
cloud.google.com/go/storage v1.65.1/go.mod h1:UsS9OgFg/XHOSYakQ8ZtLWWeyGkk1WnmD/GsGfN0BHM=
cloud.google.com/go/trace v1.16.0 h1:GmQovzFc5F0CNfl0VLgL64aoTtu7xsM0YajW2GlG9+E=
cloud.google.com/go/trace v1.16.0/go.mod h1:r+bdAn16dKLSV1G2D5v3e58IlQlizfxWrUfjx7kM7X0=
github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.32.0 h1:rIkQfkCOVKc1OiRCNcSDD8ml5RJlZbH/Xsq7lbpynwc=
github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.32.0/go.mod h1:RD2SsorTmYhF6HkTmDw7KmPYQk8OBYwTkuasChwv7R4=
github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric v0.57.0 h1:jLdiS1vO+XJFyDSWRHBx56r4s/NNtcl5J6KyCcWUX/w=
github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric v0.57.0/go.mod h1:8lmpHY+1VRoteiOwyrQMDt1YGXOrFKCz+1wJW7n3ODY=
github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/cloudmock v0.57.0 h1:cSjUzZ7KU8hicTgzaSv9NmSyM9fTVK3y5lsBUl3wOis=
github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/cloudmock v0.57.0/go.mod h1:dzcEjy1WJ0Q4u9twNR3LcLhNoYMRCrMCMafpxa0TjPQ=
github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.57.0 h1:RoO5+d7uCmDqovLrHCr2/BuViUXvdcrNxyNM1pN9dDQ=
github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.57.0/go.mod h1:YqwkQPrWSC7+byyc1VlKbWLBF5JsW5IoL6xUkemYSXk=
github.com/cespare/xxhash/v2 v2.3.0 h1:UL815xU9SqsFlibzuggzjXhog7bL6oX9BbNZnL2UFvs=
github.com/cespare/xxhash/v2 v2.3.0/go.mod h1:VGX0DQ3Q6kWi7AoAeZDth3/j3BFtOZR5XLFGgcrjCOs=
github.com/cncf/xds/go v0.0.0-20260202195803-dba9d589def2 h1:aBangftG7EVZoUb69Os8IaYg++6uMOdKK83QtkkvJik=
github.com/cncf/xds/go v0.0.0-20260202195803-dba9d589def2/go.mod h1:qwXFYgsP6T7XnJtbKlf1HP8AjxZZyzxMmc+Lq5GjlU4=
github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc h1:U9qPSI2PIWSS1VwoXQT9A3Wy9MM3WgvqSxFWenqJduM=
github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/envoyproxy/go-control-plane v0.14.0 h1:hbG2kr4RuFj222B6+7T83thSPqLjwBIfQawTkC++2HA=
github.com/envoyproxy/go-control-plane v0.14.0/go.mod h1:NcS5X47pLl/hfqxU70yPwL9ZMkUlwlKxtAohpi2wBEU=
github.com/envoyproxy/go-control-plane/envoy v1.37.0 h1:u3riX6BoYRfF4Dr7dwSOroNfdSbEPe9Yyl09/B6wBrQ=
github.com/envoyproxy/go-control-plane/envoy v1.37.0/go.mod h1:DReE9MMrmecPy+YvQOAOHNYMALuowAnbjjEMkkWOi6A=
github.com/envoyproxy/go-control-plane/ratelimit v0.1.0 h1:/G9QYbddjL25KvtKTv3an9lx6VBE2cnb8wp1vEGNYGI=
github.com/envoyproxy/go-control-plane/ratelimit v0.1.0/go.mod h1:Wk+tMFAFbCXaJPzVVHnPgRKdUdwW/KdbRt94AzgRee4=
github.com/envoyproxy/protoc-gen-validate v1.3.3 h1:MVQghNeW+LZcmXe7SY1V36Z+WFMDjpqGAGacLe2T0ds=
github.com/envoyproxy/protoc-gen-validate v1.3.3/go.mod h1:TsndJ/ngyIdQRhMcVVGDDHINPLWB7C82oDArY51KfB0=
github.com/felixge/httpsnoop v1.0.4 h1:NFTV2Zj1bL4mc9sqWACXbQFVBBg2W3GPvqp8/ESS2Wg=
github.com/felixge/httpsnoop v1.0.4/go.mod h1:m8KPJKqk1gH5J9DgRY2ASl2lWCfGKXixSwevea8zH2U=
github.com/go-jose/go-jose/v4 v4.1.4 h1:moDMcTHmvE6Groj34emNPLs/qtYXRVcd6S7NHbHz3kA=
github.com/go-jose/go-jose/v4 v4.1.4/go.mod h1:x4oUasVrzR7071A4TnHLGSPpNOm2a21K9Kf04k1rs08=
github.com/go-logr/logr v1.2.2/go.mod h1:jdQByPbusPIv2/zmleS9BjJVeZ6kBagPoEUsqbVz/1A=
github.com/go-logr/logr v1.4.3 h1:CjnDlHq8ikf6E492q6eKboGOC0T8CDaOvkHCIg8idEI=
github.com/go-logr/logr v1.4.3/go.mod h1:9T104GzyrTigFIr8wt5mBrctHMim0Nb2HLGrmQ40KvY=
github.com/go-logr/stdr v1.2.2 h1:hSWxHoqTgW2S2qGc0LTAI563KZ5YKYRhT3MFKZMbjag=
github.com/go-logr/stdr v1.2.2/go.mod h1:mMo/vtBO5dYbehREoey6XUKy/eSumjCCveDpRre4VKE=
github.com/golang/protobuf v1.5.4 h1:i7eJL8qZTpSEXOPTxNKhASYpMn+8e5Q6AdndVa1dWek=
github.com/golang/protobuf v1.5.4/go.mod h1:lnTiLA8Wa4RWRcIUkrtSVa5nRhsEGBg48fD6rSs7xps=
github.com/google/go-cmp v0.7.0 h1:wk8382ETsv4JYUZwIsn6YpYiWiBsYLSJiTsyBybVuN8=
github.com/google/go-cmp v0.7.0/go.mod h1:pXiqmnSA92OHEEa9HXL2W4E7lf9JzCmGVUdgjX3N/iU=
github.com/google/martian/v3 v3.3.3 h1:DIhPTQrbPkgs2yJYdXU/eNACCG5DVQjySNRNlflZ9Fc=
github.com/google/martian/v3 v3.3.3/go.mod h1:iEPrYcgCF7jA9OtScMFQyAlZZ4YXTKEtJ1E6RWzmBA0=
github.com/google/s2a-go v0.1.9 h1:LGD7gtMgezd8a/Xak7mEWL0PjoTQFvpRudN895yqKW0=
github.com/google/s2a-go v0.1.9/go.mod h1:YA0Ei2ZQL3acow2O62kdp9UlnvMmU7kA6Eutn0dXayM=
github.com/google/uuid v1.6.0 h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=
github.com/google/uuid v1.6.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
github.com/googleapis/enterprise-certificate-proxy v0.3.17 h1:73NfMHdiqo9JFU9+7a5ExpVa10/R29pXfZIaW559nrg=
github.com/googleapis/enterprise-certificate-proxy v0.3.17/go.mod h1:rSEsBUemEBZEexP2y6jPp16LUmUbjmSbcPMQizR0o4k=
github.com/googleapis/gax-go/v2 v2.23.0 h1:Tchl7qkvE7Ip3y+ztvNufYFvkfqTe7NfLTYGIdJRLuE=
github.com/googleapis/gax-go/v2 v2.23.0/go.mod h1:rBQKOVJCdb8IFEzg+FCwlt1LP/xMDGuqUXhUG+XMXEg=
github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 h1:GFCKgmp0tecUJ0sJuv4pzYCqS9+RGSn52M3FUwPs+uo=
github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10/go.mod h1:t/avpk3KcrXxUnYOhZhMXJlSEyie6gQbtLq5NM3loB8=
github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 h1:Jamvg5psRIccs7FGNTlIRMkT8wgtp5eCXdBlqhYGL6U=
github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
github.com/spiffe/go-spiffe/v2 v2.6.0 h1:l+DolpxNWYgruGQVV0xsfeya3CsC7m8iBzDnMpsbLuo=
github.com/spiffe/go-spiffe/v2 v2.6.0/go.mod h1:gm2SeUoMZEtpnzPNs2Csc0D/gX33k1xIx7lEzqblHEs=
github.com/stretchr/testify v1.11.1 h1:7s2iGBzp5EwR7/aIZr8ao5+dra3wiQyKjjFuvgVKu7U=
github.com/stretchr/testify v1.11.1/go.mod h1:wZwfW3scLgRK+23gO65QZefKpKQRnfz6sD981Nm4B6U=
go.opentelemetry.io/auto/sdk v1.2.1 h1:jXsnJ4Lmnqd11kwkBV2LgLoFMZKizbCi5fNZ/ipaZ64=
go.opentelemetry.io/auto/sdk v1.2.1/go.mod h1:KRTj+aOaElaLi+wW1kO/DZRXwkF4C5xPbEe3ZiIhN7Y=
go.opentelemetry.io/contrib/detectors/gcp v1.43.0 h1:62yY3dT7/ShwOxzA0RsKRgshBmfElKI4d/Myu2OxDFU=
go.opentelemetry.io/contrib/detectors/gcp v1.43.0/go.mod h1:RyaZMFY7yi1kAs45S6mbFGz8O8rqB0dTY14uzvG4LCs=
go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.68.0 h1:0Qx7VGBacMm9ZENQ7TnNObTYI4ShC+lHI16seduaxZo=
go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.68.0/go.mod h1:Sje3i3MjSPKTSPvVWCaL8ugBzJwik3u4smCjUeuupqg=
go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.67.0 h1:OyrsyzuttWTSur2qN/Lm0m2a8yqyIjUVBZcxFPuXq2o=
go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.67.0/go.mod h1:C2NGBr+kAB4bk3xtMXfZ94gqFDtg/GkI7e9zqGh5Beg=
go.opentelemetry.io/otel v1.44.0 h1:JjwHmHpA4iZ3wBxluu2fbbE7j4kqlE8jXyAyPXH7HqU=
go.opentelemetry.io/otel v1.44.0/go.mod h1:BMgjTHL9WPRlRjL2oZCBTL4whCGtXch2H4BhOPIAyYc=
go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.44.0 h1:hqxVTu/GtBF+vJ8d1fzW7fRxZFvgoDjWcxwwCaFDYpU=
go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.44.0/go.mod h1:z5fVEF4X5v0ESvlJqBrrFlBVoj5EQuefZpzsu7R+x5Q=
go.opentelemetry.io/otel/metric v1.44.0 h1:1w0gILTcHdr3YI+ixLyjemwrVnsMURbTZFrSYCdDdmc=
go.opentelemetry.io/otel/metric v1.44.0/go.mod h1:8O7hanEPBNgEMmybD3s2VBKcgWOCsA6tzHBPODAiquo=
go.opentelemetry.io/otel/metric/x v0.66.0 h1:YkCrx1zLOChi9ZcZ6euupOcsgzbVlec7D/xoEU1+cTA=
go.opentelemetry.io/otel/metric/x v0.66.0/go.mod h1:d1+BDj9t96do0/1LoU1ayfCv79ZgNE41qbhBvnMOBZk=
go.opentelemetry.io/otel/sdk v1.44.0 h1:nHYwb9lK+fJPU/dnT6s7W7Z8itMWyqrnVfbheVYrZ58=
go.opentelemetry.io/otel/sdk v1.44.0/go.mod h1:Osuydd3Se74nqjAKxid74N5eC+jfEqfTegHRnq58oK0=
go.opentelemetry.io/otel/sdk/metric v1.44.0 h1:3LlKgI+VjbVsjNRFZJZAJ30WjXC5VkNRks6si09iEfI=
go.opentelemetry.io/otel/sdk/metric v1.44.0/go.mod h1:5B5pMARnXxKhltooO4xUuCBorl65a4EpnTalObqOigA=
go.opentelemetry.io/otel/trace v1.44.0 h1:jxF5CsGYCe74MCRx2X4g7WsY/VBKRqqpNvXlX/6gtIk=
go.opentelemetry.io/otel/trace v1.44.0/go.mod h1:oLl1jrMQAVo6v3GAggN+1VH9VIz9iUSvW53sW1Q8PIE=
golang.org/x/crypto v0.53.0 h1:QZ4Muo8THX6CizN2vPPd5fBGHyogrdK9fG4wLPFUsto=
golang.org/x/crypto v0.53.0/go.mod h1:DNLU434OwVakk9PzuwV8w62mAJpRJL3vsgcfp4Qnsio=
golang.org/x/net v0.56.0 h1:Rw8j/hFzGvJUZwNBXnAtf5sVDVt+65SK2C7IxCxZt5o=
golang.org/x/net v0.56.0/go.mod h1:D3Ku6r+V6JROoZK144D2XfMHFcMq/0zSfLelVTCFKec=
golang.org/x/oauth2 v0.36.0 h1:peZ/1z27fi9hUOFCAZaHyrpWG5lwe0RJEEEeH0ThlIs=
golang.org/x/oauth2 v0.36.0/go.mod h1:YDBUJMTkDnJS+A4BP4eZBjCqtokkg1hODuPjwiGPO7Q=
golang.org/x/sync v0.21.0 h1:HLII4xRRTtCRkxYp4HNFF0Js/Og6q2i++KXbg0gHCwM=
golang.org/x/sync v0.21.0/go.mod h1:9xrNwdLfx4jkKbNva9FpL6vEN7evnE43NNNJQ2LF3+0=
golang.org/x/sys v0.46.0 h1:noSf2Fq6F8DBgS+LysIkx7rIExoNHJsxOAtPp4rthXw=
golang.org/x/sys v0.46.0/go.mod h1:4GL1E5IUh+htKOUEOaiffhrAeqysfVGipDYzABqnCmw=
golang.org/x/text v0.39.0 h1:UbZz4pLOvn600D6Oh6GGEI6VAmndrEBLv8/6BEXzyus=
golang.org/x/text v0.39.0/go.mod h1:3UwRclnC2g0TU9x8PZiyfOajCd1zaUNHF9cvqcQZ+ZM=
golang.org/x/time v0.15.0 h1:bbrp8t3bGUeFOx08pvsMYRTCVSMk89u4tKbNOZbp88U=
golang.org/x/time v0.15.0/go.mod h1:Y4YMaQmXwGQZoFaVFk4YpCt4FLQMYKZe9oeV/f4MSno=
gonum.org/v1/gonum v0.17.0 h1:VbpOemQlsSMrYmn7T2OUvQ4dqxQXU+ouZFQsZOx50z4=
gonum.org/v1/gonum v0.17.0/go.mod h1:El3tOrEuMpv2UdMrbNlKEh9vd86bmQ6vqIcDwxEOc1E=
google.golang.org/api v0.287.1 h1:LiyJx32VU3cwQfLchn/513qKhc25hq0pEANYJoWNnnI=
google.golang.org/api v0.287.1/go.mod h1:lM2kYRzYUCBY91P9h6VF1PYmvhxii3O5hji37qRvIcY=
google.golang.org/genproto v0.0.0-20260519071638-aa98bba5eb94 h1:YJjbgu+dkp5kUJLfpMyCLfBIWZb/FcJyuLeo1gVBOuo=
google.golang.org/genproto v0.0.0-20260519071638-aa98bba5eb94/go.mod h1:RRHjglSYABVCWpQ7USCpdfhcd9t4PkajvVwyynZizTc=
google.golang.org/genproto/googleapis/api v0.0.0-20260630182238-925bb5da69e7 h1:jQ9p21COKWjP3VwuFrNRiiOTMh3mPpN45R7SLrH/HUU=
google.golang.org/genproto/googleapis/api v0.0.0-20260630182238-925bb5da69e7/go.mod h1:KqHwBx2upmfa1XSi1WuRvC+2VGCLtooKkfmyvRbUmqA=
google.golang.org/genproto/googleapis/rpc v0.0.0-20260630182238-925bb5da69e7 h1:eM/YSd5bBFagF51o1E745Ta7RwzpW0h+z+QDNZOgmQ8=
google.golang.org/genproto/googleapis/rpc v0.0.0-20260630182238-925bb5da69e7/go.mod h1:4Hqkh8ycfw05ld/3BWL7rJOSfebL2Q+DVDeRgYgxUU8=
google.golang.org/grpc v1.82.1 h1:NnAxzGRA0677vCa4BUkOAnO5+FfQqVl9iUXeD0IqcGE=
google.golang.org/grpc v1.82.1/go.mod h1:yzTZ1TB1Z3SG+LIYaI+WiE8D5+PZ3ArnrSp8zF3+/ZA=
google.golang.org/protobuf v1.36.11 h1:fV6ZwhNocDyBLK0dj+fg8ektcVegBBuEolpbTQyBNVE=
google.golang.org/protobuf v1.36.11/go.mod h1:HTf+CrKn2C3g5S8VImy6tdcUvCska2kB7j23XfzDpco=
gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
````

### FILE: `google_cloud_storage_adapter/provider-profile.template.json`
```yaml
block_id: "GO-GOOGLE-CLOUD-STORAGE-ADAPTER:profile:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed Google Cloud storage authority profile"
license: "LicenseRef-Workspace-Owner"
sha256: "3d3bca488f2fb5febcd67ffde5667139c857f7dac7dd91e30b7d606c54a87aba"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-google-cloud-storage-profile/v1",
  "provider": "Google Cloud Storage",
  "sdk": "cloud.google.com/go/storage@v1.65.1",
  "credentials": {
    "mechanism": "Application Default Credentials",
    "secret_reference_required": true
  },
  "object_storage": {
    "decision": "BLOCKED_ACCESS_COST_AND_POLICY_REQUIRED",
    "project_number": "",
    "bucket": "",
    "object_prefix": "",
    "max_bytes": 0,
    "kms_key_resource": "",
    "object_retention_mode_required": "Enabled",
    "retention": {
      "mode": "",
      "retain_until_source": "",
      "policy_approved": false,
      "irreversible_locked_mode_approved": false
    },
    "access_approved": false,
    "cost_approved": false,
    "automatic_storage_authorized": false
  }
}
````

### FILE: `google_cloud_storage_adapter/official-artifact-lock.json`
```yaml
block_id: "GO-GOOGLE-CLOUD-STORAGE-ADAPTER:lock:v1"
operation: CREATE
provenance: AUTHORED
source: "official Go proxy and GitHub identity/license audit"
license: "LicenseRef-Workspace-Owner"
sha256: "2c33e9aa7420835a5427fb055324e2000ed0a5fee50f798ef9550e4d0f3f010a"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-official-go-module-lock/v1",
  "module": "cloud.google.com/go/storage",
  "version": "v1.65.1",
  "tag": "storage/v1.65.1",
  "commit": "3ab7d1390bbbb40cf197a545b941f9df760a3269",
  "commit_signature_verified": true,
  "published_at": "2026-08-26T17:25:54Z",
  "go_proxy_zip": "https://proxy.golang.org/cloud.google.com/go/storage/@v/v1.65.1.zip",
  "zip_bytes": 998644,
  "zip_sha256": "0030613d00242a7b2cd0133e86e67c12c066925e88909f185d56706c0bdb6d64",
  "license": "Apache-2.0",
  "license_sha256": "cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30"
}
````

### FILE: `google_cloud_storage_adapter/gcsstorage/storage.go`
```yaml
block_id: "GO-GOOGLE-CLOUD-STORAGE-ADAPTER:store:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed orchestration over official SDK primitives"
license: "LicenseRef-Workspace-Owner"
sha256: "4925eb9ae4400adb161964e068d836a47afe188a9bba51f4b34c1edb4af7e0a6"
variables: []
secrets_allowed: false
```
````go
package gcsstorage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/crc32"
	"regexp"
	"strings"
	"time"
)

const (
	ModeLocked   = "Locked"
	ModeUnlocked = "Unlocked"

	maxRetention = time.Duration(3155760000) * time.Second
)

var (
	bucketPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{1,61}[a-z0-9]$`)
	kmsKeyPattern = regexp.MustCompile(`^projects/[a-z][a-z0-9-]{4,28}[a-z0-9]/locations/[a-z0-9-]+/keyRings/[A-Za-z0-9_-]{1,63}/cryptoKeys/[A-Za-z0-9_-]{1,63}$`)
)

type Config struct {
	Bucket                   string
	ExpectedProjectNumber    uint64
	KMSKeyName               string
	RetentionMode            string
	RetainUntil              time.Time
	MaxBytes                 int64
	AccessApproved           bool
	CostApproved             bool
	RetentionPolicyApproved  bool
	IrreversibleLockApproved bool
}

type Request struct {
	ObjectName  string
	ContentType string
	Data        []byte
}

type BucketState struct {
	ProjectNumber       uint64
	ObjectRetentionMode string
}

type CreateRequest struct {
	Bucket        string
	ObjectName    string
	ContentType   string
	Data          []byte
	CRC32C        uint32
	KMSKeyName    string
	RetentionMode string
	RetainUntil   time.Time
	Metadata      map[string]string
}

type ObjectState struct {
	Generation    int64
	Size          int64
	CRC32C        uint32
	KMSKeyName    string
	RetentionMode string
	RetainUntil   time.Time
	ETag          string
}

type Backend interface {
	BucketState(context.Context, string) (BucketState, error)
	Create(context.Context, CreateRequest) (ObjectState, error)
	Inspect(context.Context, string, string, int64) (ObjectState, error)
}

type Receipt struct {
	Schema            string `json:"schema"`
	Provider          string `json:"provider"`
	ProjectNumber     uint64 `json:"project_number"`
	BucketSHA256      string `json:"bucket_sha256"`
	ObjectNameSHA256  string `json:"object_name_sha256"`
	Generation        int64  `json:"generation"`
	Size              int64  `json:"size"`
	CRC32C            uint32 `json:"crc32c"`
	ContentSHA256     string `json:"content_sha256"`
	KMSKeyNameSHA256  string `json:"kms_key_name_sha256"`
	RetentionMode     string `json:"retention_mode"`
	RetainUntil       string `json:"retain_until"`
	Created           bool   `json:"created"`
	RetentionVerified bool   `json:"retention_verified"`
}

type PartialVerificationError struct {
	Generation int64
	Cause      error
}

func (e *PartialVerificationError) Error() string {
	return fmt.Sprintf("object generation %d was created but post-write verification failed: %v", e.Generation, e.Cause)
}

func (e *PartialVerificationError) Unwrap() error { return e.Cause }

type Store struct {
	backend Backend
	now     func() time.Time
}

func NewStore(backend Backend, now func() time.Time) (*Store, error) {
	if backend == nil {
		return nil, errors.New("backend is required")
	}
	if now == nil {
		return nil, errors.New("clock is required")
	}
	return &Store{backend: backend, now: now}, nil
}

func (s *Store) Put(ctx context.Context, cfg Config, req Request) (Receipt, error) {
	if ctx == nil {
		return Receipt{}, errors.New("context is required")
	}
	now := s.now().UTC()
	retainUntil, err := validate(cfg, req, now)
	if err != nil {
		return Receipt{}, err
	}

	bucketState, err := s.backend.BucketState(ctx, cfg.Bucket)
	if err != nil {
		return Receipt{}, fmt.Errorf("read bucket state: %w", err)
	}
	if bucketState.ProjectNumber != cfg.ExpectedProjectNumber {
		return Receipt{}, fmt.Errorf("bucket project number mismatch: got %d", bucketState.ProjectNumber)
	}
	if bucketState.ObjectRetentionMode != "Enabled" {
		return Receipt{}, fmt.Errorf("bucket object retention is not Enabled")
	}

	contentHash := sha256Hex(req.Data)
	crc := crc32.Checksum(req.Data, crc32.MakeTable(crc32.Castagnoli))
	created, err := s.backend.Create(ctx, CreateRequest{
		Bucket:        cfg.Bucket,
		ObjectName:    req.ObjectName,
		ContentType:   req.ContentType,
		Data:          req.Data,
		CRC32C:        crc,
		KMSKeyName:    cfg.KMSKeyName,
		RetentionMode: cfg.RetentionMode,
		RetainUntil:   retainUntil,
		Metadata: map[string]string{
			"elite-content-sha256": contentHash,
			"elite-retention-mode": strings.ToLower(cfg.RetentionMode),
		},
	})
	if err != nil {
		return Receipt{}, fmt.Errorf("create object: %w", err)
	}

	receipt := newReceipt(cfg, req, created, contentHash, retainUntil)
	if created.Generation <= 0 {
		return receipt, &PartialVerificationError{Generation: created.Generation, Cause: errors.New("provider returned no generation")}
	}

	observed, err := s.backend.Inspect(ctx, cfg.Bucket, req.ObjectName, created.Generation)
	if err != nil {
		return receipt, &PartialVerificationError{Generation: created.Generation, Cause: fmt.Errorf("inspect exact generation: %w", err)}
	}
	if err := verifyObserved(observed, created.Generation, int64(len(req.Data)), crc, cfg.KMSKeyName, cfg.RetentionMode, retainUntil); err != nil {
		return receipt, &PartialVerificationError{Generation: created.Generation, Cause: err}
	}
	receipt.RetentionVerified = true
	return receipt, nil
}

func validate(cfg Config, req Request, now time.Time) (time.Time, error) {
	if !cfg.AccessApproved || !cfg.CostApproved || !cfg.RetentionPolicyApproved {
		return time.Time{}, errors.New("access, cost and retention policy approvals are required")
	}
	if !bucketPattern.MatchString(cfg.Bucket) || strings.Contains(cfg.Bucket, "..") {
		return time.Time{}, errors.New("invalid bucket name")
	}
	if cfg.ExpectedProjectNumber == 0 {
		return time.Time{}, errors.New("expected project number is required")
	}
	if !kmsKeyPattern.MatchString(cfg.KMSKeyName) {
		return time.Time{}, errors.New("invalid Cloud KMS key resource name")
	}
	if cfg.RetentionMode != ModeLocked && cfg.RetentionMode != ModeUnlocked {
		return time.Time{}, errors.New("retention mode must be Locked or Unlocked")
	}
	if cfg.RetentionMode == ModeLocked && !cfg.IrreversibleLockApproved {
		return time.Time{}, errors.New("Locked retention requires irreversible lock approval")
	}
	retainUntil := cfg.RetainUntil.UTC().Truncate(time.Second)
	if !retainUntil.After(now.Add(time.Minute)) {
		return time.Time{}, errors.New("retain-until must be more than one minute in the future")
	}
	if retainUntil.Sub(now) > maxRetention {
		return time.Time{}, errors.New("retain-until exceeds Google Cloud Storage 100-year maximum")
	}
	if cfg.MaxBytes <= 0 || int64(len(req.Data)) > cfg.MaxBytes {
		return time.Time{}, errors.New("content exceeds configured byte limit")
	}
	if len(req.Data) == 0 {
		return time.Time{}, errors.New("content is empty")
	}
	if err := validateObjectName(req.ObjectName); err != nil {
		return time.Time{}, err
	}
	if req.ContentType == "" || len(req.ContentType) > 255 || strings.ContainsAny(req.ContentType, "\r\n") {
		return time.Time{}, errors.New("invalid content type")
	}
	return retainUntil, nil
}

func validateObjectName(name string) error {
	if name == "" || len(name) > 1024 || strings.HasPrefix(name, "/") || strings.HasSuffix(name, "/") || strings.Contains(name, `\`) {
		return errors.New("invalid object name")
	}
	for _, r := range name {
		if r < 0x20 || r == 0x7f {
			return errors.New("object name contains control characters")
		}
	}
	for _, segment := range strings.Split(name, "/") {
		if segment == "" || segment == "." || segment == ".." {
			return errors.New("object name contains an unsafe path segment")
		}
	}
	return nil
}

func verifyObserved(got ObjectState, generation, size int64, crc uint32, kmsKey, mode string, retainUntil time.Time) error {
	switch {
	case got.Generation != generation:
		return fmt.Errorf("generation mismatch: got %d", got.Generation)
	case got.Size != size:
		return fmt.Errorf("size mismatch: got %d", got.Size)
	case got.CRC32C != crc:
		return fmt.Errorf("CRC32C mismatch: got %d", got.CRC32C)
	case got.KMSKeyName != kmsKey:
		return errors.New("Cloud KMS key mismatch")
	case got.RetentionMode != mode:
		return fmt.Errorf("retention mode mismatch: got %q", got.RetentionMode)
	case !got.RetainUntil.Equal(retainUntil):
		return fmt.Errorf("retain-until mismatch: got %s", got.RetainUntil.UTC().Format(time.RFC3339))
	case got.ETag == "":
		return errors.New("provider returned no ETag")
	default:
		return nil
	}
}

func newReceipt(cfg Config, req Request, created ObjectState, contentHash string, retainUntil time.Time) Receipt {
	return Receipt{
		Schema:            "elite-google-cloud-storage-receipt/v1",
		Provider:          "Google Cloud Storage",
		ProjectNumber:     cfg.ExpectedProjectNumber,
		BucketSHA256:      sha256Hex([]byte(cfg.Bucket)),
		ObjectNameSHA256:  sha256Hex([]byte(req.ObjectName)),
		Generation:        created.Generation,
		Size:              int64(len(req.Data)),
		CRC32C:            crc32.Checksum(req.Data, crc32.MakeTable(crc32.Castagnoli)),
		ContentSHA256:     contentHash,
		KMSKeyNameSHA256:  sha256Hex([]byte(cfg.KMSKeyName)),
		RetentionMode:     cfg.RetentionMode,
		RetainUntil:       retainUntil.Format(time.RFC3339),
		Created:           true,
		RetentionVerified: false,
	}
}

func sha256Hex(value []byte) string {
	sum := sha256.Sum256(value)
	return hex.EncodeToString(sum[:])
}
````

### FILE: `google_cloud_storage_adapter/gcsstorage/google_backend.go`
```yaml
block_id: "GO-GOOGLE-CLOUD-STORAGE-ADAPTER:backend:v1"
operation: CREATE
provenance: AUTHORED
source: "integration against cloud.google.com/go/storage v1.65.1 public API"
license: "LicenseRef-Workspace-Owner"
sha256: "067c9668555f1b40488400026e1c98ee5f92997cc3e042cade4e2115db5e6242"
variables: []
secrets_allowed: false
```
````go
package gcsstorage

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
)

type GoogleCloudBackend struct {
	client *storage.Client
}

func NewGoogleCloudBackend(client *storage.Client) (*GoogleCloudBackend, error) {
	if client == nil {
		return nil, fmt.Errorf("Google Cloud Storage client is required")
	}
	return &GoogleCloudBackend{client: client}, nil
}

func (b *GoogleCloudBackend) BucketState(ctx context.Context, bucket string) (BucketState, error) {
	attrs, err := b.client.Bucket(bucket).Attrs(ctx)
	if err != nil {
		return BucketState{}, err
	}
	return BucketState{ProjectNumber: attrs.ProjectNumber, ObjectRetentionMode: attrs.ObjectRetentionMode}, nil
}

func (b *GoogleCloudBackend) Create(ctx context.Context, req CreateRequest) (ObjectState, error) {
	object := b.client.Bucket(req.Bucket).Object(req.ObjectName).If(storage.Conditions{DoesNotExist: true})
	writer := object.NewWriter(ctx)
	writer.ContentType = req.ContentType
	writer.Metadata = cloneMetadata(req.Metadata)
	writer.KMSKeyName = req.KMSKeyName
	writer.CRC32C = req.CRC32C
	writer.SendCRC32C = true
	writer.Retention = &storage.ObjectRetention{Mode: req.RetentionMode, RetainUntil: req.RetainUntil}

	n, err := writer.Write(req.Data)
	if err != nil {
		_ = writer.CloseWithError(err)
		return ObjectState{}, err
	}
	if n != len(req.Data) {
		err := fmt.Errorf("short write: wrote %d of %d bytes", n, len(req.Data))
		_ = writer.CloseWithError(err)
		return ObjectState{}, err
	}
	if err := writer.Close(); err != nil {
		return ObjectState{}, err
	}
	return stateFromAttrs(writer.Attrs())
}

func (b *GoogleCloudBackend) Inspect(ctx context.Context, bucket, name string, generation int64) (ObjectState, error) {
	attrs, err := b.client.Bucket(bucket).Object(name).Generation(generation).Attrs(ctx)
	if err != nil {
		return ObjectState{}, err
	}
	return stateFromAttrs(attrs)
}

func stateFromAttrs(attrs *storage.ObjectAttrs) (ObjectState, error) {
	if attrs == nil {
		return ObjectState{}, fmt.Errorf("Google Cloud Storage returned no object attributes")
	}
	state := ObjectState{
		Generation: attrs.Generation,
		Size:       attrs.Size,
		CRC32C:     attrs.CRC32C,
		KMSKeyName: attrs.KMSKeyName,
		ETag:       attrs.Etag,
	}
	if attrs.Retention != nil {
		state.RetentionMode = attrs.Retention.Mode
		state.RetainUntil = attrs.Retention.RetainUntil.UTC()
	}
	return state, nil
}

func cloneMetadata(src map[string]string) map[string]string {
	dst := make(map[string]string, len(src))
	for key, value := range src {
		dst[key] = value
	}
	return dst
}
````

### FILE: `google_cloud_storage_adapter/gcsstorage/storage_test.go`
```yaml
block_id: "GO-GOOGLE-CLOUD-STORAGE-ADAPTER:tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local deterministic negative and partial-effect tests"
license: "LicenseRef-Workspace-Owner"
sha256: "d30537bbc7adab99d5c1f1feff3268d92ddc0689d8b925c7b82b1102850df946"
variables: []
secrets_allowed: false
```
````go
package gcsstorage

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

var fixedNow = time.Date(2026, 8, 27, 15, 0, 0, 0, time.UTC)

type fakeBackend struct {
	bucket     BucketState
	bucketErr  error
	create     ObjectState
	createErr  error
	inspect    ObjectState
	inspectErr error
	createReq  CreateRequest
	inspectGen int64
}

func (f *fakeBackend) BucketState(context.Context, string) (BucketState, error) {
	return f.bucket, f.bucketErr
}
func (f *fakeBackend) Create(_ context.Context, req CreateRequest) (ObjectState, error) {
	f.createReq = req
	return f.create, f.createErr
}
func (f *fakeBackend) Inspect(_ context.Context, _, _ string, generation int64) (ObjectState, error) {
	f.inspectGen = generation
	return f.inspect, f.inspectErr
}

func validConfig() Config {
	return Config{
		Bucket:                   "elite-document-evidence",
		ExpectedProjectNumber:    123456789012,
		KMSKeyName:               "projects/elite1/locations/us/keyRings/documents/cryptoKeys/evidence",
		RetentionMode:            ModeLocked,
		RetainUntil:              fixedNow.Add(24 * time.Hour),
		MaxBytes:                 1 << 20,
		AccessApproved:           true,
		CostApproved:             true,
		RetentionPolicyApproved:  true,
		IrreversibleLockApproved: true,
	}
}

func validRequest() Request {
	return Request{ObjectName: "tenant-hash/source-hash.pdf", ContentType: "application/pdf", Data: []byte("document bytes")}
}

func validBackend() *fakeBackend {
	data := validRequest().Data
	crc := crcForTest(data)
	state := ObjectState{
		Generation:    42,
		Size:          int64(len(data)),
		CRC32C:        crc,
		KMSKeyName:    validConfig().KMSKeyName,
		RetentionMode: ModeLocked,
		RetainUntil:   validConfig().RetainUntil,
		ETag:          "etag-42",
	}
	return &fakeBackend{bucket: BucketState{ProjectNumber: 123456789012, ObjectRetentionMode: "Enabled"}, create: state, inspect: state}
}

func TestPutVerified(t *testing.T) {
	backend := validBackend()
	store, err := NewStore(backend, func() time.Time { return fixedNow })
	if err != nil {
		t.Fatal(err)
	}
	receipt, err := store.Put(context.Background(), validConfig(), validRequest())
	if err != nil {
		t.Fatal(err)
	}
	if !receipt.Created || !receipt.RetentionVerified || receipt.Generation != 42 {
		t.Fatalf("unexpected receipt: %+v", receipt)
	}
	if receipt.BucketSHA256 == validConfig().Bucket || receipt.ObjectNameSHA256 == validRequest().ObjectName {
		t.Fatal("receipt leaked raw identifiers")
	}
	if backend.inspectGen != 42 {
		t.Fatalf("inspected generation %d", backend.inspectGen)
	}
	if backend.createReq.CRC32C != crcForTest(validRequest().Data) {
		t.Fatal("CRC32C not propagated")
	}
	if backend.createReq.Metadata["elite-content-sha256"] != receipt.ContentSHA256 {
		t.Fatal("content hash metadata mismatch")
	}
}

func TestConfigRejectsUnsafeInputs(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Config, *Request)
	}{
		{"approvals", func(c *Config, _ *Request) { c.CostApproved = false }},
		{"bucket", func(c *Config, _ *Request) { c.Bucket = "Bad Bucket" }},
		{"project", func(c *Config, _ *Request) { c.ExpectedProjectNumber = 0 }},
		{"kms", func(c *Config, _ *Request) { c.KMSKeyName = "alias/key" }},
		{"mode", func(c *Config, _ *Request) { c.RetentionMode = "Governance" }},
		{"locked-approval", func(c *Config, _ *Request) { c.IrreversibleLockApproved = false }},
		{"expired", func(c *Config, _ *Request) { c.RetainUntil = fixedNow }},
		{"too-long", func(c *Config, _ *Request) { c.RetainUntil = fixedNow.Add(maxRetention + time.Second) }},
		{"limit", func(c *Config, _ *Request) { c.MaxBytes = 1 }},
		{"empty", func(_ *Config, r *Request) { r.Data = nil }},
		{"traversal", func(_ *Config, r *Request) { r.ObjectName = "a/../b.pdf" }},
		{"content-type", func(_ *Config, r *Request) { r.ContentType = "application/pdf\r\nx: y" }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, req := validConfig(), validRequest()
			tt.mutate(&cfg, &req)
			backend := validBackend()
			store, _ := NewStore(backend, func() time.Time { return fixedNow })
			if _, err := store.Put(context.Background(), cfg, req); err == nil {
				t.Fatal("expected rejection")
			}
			if !reflect.DeepEqual(backend.createReq, CreateRequest{}) {
				t.Fatal("provider called after validation rejection")
			}
		})
	}
}

func TestBucketAuthorityAndRetentionRequired(t *testing.T) {
	for _, state := range []BucketState{
		{ProjectNumber: 999, ObjectRetentionMode: "Enabled"},
		{ProjectNumber: 123456789012, ObjectRetentionMode: ""},
	} {
		backend := validBackend()
		backend.bucket = state
		store, _ := NewStore(backend, func() time.Time { return fixedNow })
		if _, err := store.Put(context.Background(), validConfig(), validRequest()); err == nil {
			t.Fatal("expected bucket rejection")
		}
		if !reflect.DeepEqual(backend.createReq, CreateRequest{}) {
			t.Fatal("provider write occurred")
		}
	}
}

func TestCreateFailureIsNotReportedAsCreated(t *testing.T) {
	backend := validBackend()
	backend.createErr = errors.New("precondition failed")
	store, _ := NewStore(backend, func() time.Time { return fixedNow })
	receipt, err := store.Put(context.Background(), validConfig(), validRequest())
	if err == nil || receipt.Created {
		t.Fatalf("receipt=%+v err=%v", receipt, err)
	}
}

func TestPostWriteFailureReturnsPartialReceipt(t *testing.T) {
	backend := validBackend()
	backend.inspectErr = errors.New("read unavailable")
	store, _ := NewStore(backend, func() time.Time { return fixedNow })
	receipt, err := store.Put(context.Background(), validConfig(), validRequest())
	var partial *PartialVerificationError
	if !errors.As(err, &partial) {
		t.Fatalf("expected partial error, got %v", err)
	}
	if !receipt.Created || receipt.RetentionVerified || receipt.Generation != 42 {
		t.Fatalf("unexpected receipt: %+v", receipt)
	}
}

func TestObservedMismatchFailsClosed(t *testing.T) {
	mutations := []func(*ObjectState){
		func(s *ObjectState) { s.Generation++ },
		func(s *ObjectState) { s.Size++ },
		func(s *ObjectState) { s.CRC32C++ },
		func(s *ObjectState) { s.KMSKeyName = "projects/other1/locations/us/keyRings/r/cryptoKeys/k" },
		func(s *ObjectState) { s.RetentionMode = ModeUnlocked },
		func(s *ObjectState) { s.RetainUntil = s.RetainUntil.Add(time.Second) },
		func(s *ObjectState) { s.ETag = "" },
	}
	for index, mutate := range mutations {
		backend := validBackend()
		mutate(&backend.inspect)
		store, _ := NewStore(backend, func() time.Time { return fixedNow })
		receipt, err := store.Put(context.Background(), validConfig(), validRequest())
		if err == nil || !receipt.Created || receipt.RetentionVerified {
			t.Fatalf("case %d receipt=%+v err=%v", index, receipt, err)
		}
	}
}

func TestUnlockedDoesNotRequireIrreversibleApproval(t *testing.T) {
	cfg := validConfig()
	cfg.RetentionMode = ModeUnlocked
	cfg.IrreversibleLockApproved = false
	backend := validBackend()
	backend.create.RetentionMode, backend.inspect.RetentionMode = ModeUnlocked, ModeUnlocked
	store, _ := NewStore(backend, func() time.Time { return fixedNow })
	if _, err := store.Put(context.Background(), cfg, validRequest()); err != nil {
		t.Fatal(err)
	}
}

func TestConstructorsRejectNil(t *testing.T) {
	if _, err := NewStore(nil, func() time.Time { return fixedNow }); err == nil {
		t.Fatal("nil backend accepted")
	}
	if _, err := NewStore(validBackend(), nil); err == nil {
		t.Fatal("nil clock accepted")
	}
	if _, err := NewGoogleCloudBackend(nil); err == nil {
		t.Fatal("nil Google client accepted")
	}
	if _, err := stateFromAttrs(nil); err == nil {
		t.Fatal("nil provider attributes accepted")
	}
	store, _ := NewStore(validBackend(), func() time.Time { return fixedNow })
	if _, err := store.Put(nil, validConfig(), validRequest()); err == nil {
		t.Fatal("nil context accepted")
	}
}

func crcForTest(data []byte) uint32 {
	return crc32Checksum(data)
}
````

### FILE: `google_cloud_storage_adapter/gcsstorage/crc_test.go`
```yaml
block_id: "GO-GOOGLE-CLOUD-STORAGE-ADAPTER:crc-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local CRC32C test helper"
license: "LicenseRef-Workspace-Owner"
sha256: "5088535a35d1ef378a202e123617dba278bfc7bcb873e4b651beb26c8b762e4d"
variables: []
secrets_allowed: false
```
````go
package gcsstorage

import "hash/crc32"

func crc32Checksum(data []byte) uint32 {
	return crc32.Checksum(data, crc32.MakeTable(crc32.Castagnoli))
}
````

### FILE: `google_cloud_storage_adapter/README.md`
```yaml
block_id: "GO-GOOGLE-CLOUD-STORAGE-ADAPTER:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "local usage and non-claim boundary"
license: "LicenseRef-Workspace-Owner"
sha256: "88ddf649dfd4b069543a93948599754587fa182a7fcf4ed494d92ef753283979"
variables: []
secrets_allowed: false
```
````markdown
# Google Cloud Storage evidence adapter

This module is local `AUTHORED` integration glue compiled against Google's official `cloud.google.com/go/storage@v1.65.1`; it is not Google-authored application code.

It performs one immutable-by-name object create with the official `DoesNotExist` generation precondition, caller-calculated CRC32C validation, an exact Cloud KMS key and per-object retention. Before writing, the configured bucket must report the expected numeric project owner and `ObjectRetentionMode=Enabled`. After writing, the adapter reads the exact returned generation and verifies generation, size, CRC32C, KMS key, retention mode/date and ETag. A post-write mismatch returns a non-empty partial receipt plus an error so callers reconcile that generation instead of retrying blindly.

The adapter never creates or reconfigures buckets, enables Object Retention Lock, locks an existing policy, changes IAM, deletes objects or makes a cloud call without explicit access, cost and retention approvals. `Locked` retention additionally requires explicit irreversible-lock approval. Application Default Credentials remain outside files and receipts.

Run in a materialized tree with the exact toolchain and downloaded module graph:

```text
go mod verify
go test ./...
go vet ./...
go build ./...
```

These offline gates do not prove target IAM, service perimeter, KMS grants, quota/cost, retention compliance, audit logging, malware policy, recovery or live Google Cloud behavior. Those remain project gates before enabling `automatic_storage_authorized`.
````

## 6. Configuration surface

| Variable | Tipo | Default seguro | Validación | Secreto | Mutabilidad/efecto |
|---|---|---|---|---|---|
| bucket | string | vacío/bloqueado | nombre GCS estricto; owner se verifica por número | no | cambio exige nueva aprobación |
| expected_project_number | uint64 | 0/bloqueado | mayor que cero; igualdad con BucketAttrs | no | autoridad |
| kms_key_resource | string | vacío/bloqueado | resource name Cloud KMS completo | no | nueva aprobación/reconciliación |
| retention_mode | enum | vacío/bloqueado | `Locked` o `Unlocked` | no | `Locked` irreversible |
| retain_until | RFC3339 | vacío/bloqueado | futuro >1 minuto y máximo oficial 100 años | no | nunca reducir sin proceso externo |
| max_bytes | int64 | 0/bloqueado | positivo y aplicado antes de provider | no | policy |
| access/cost/policy approvals | boolean | false | los tres true antes de leer bucket | no | decisión humana |
| irreversible_locked_mode_approved | boolean | false | true sólo si mode `Locked` | no | irreversible |
| ADC | credential chain | ausente | mecanismo externo; nunca serializado | sí | rotación externa |

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Google Cloud Storage Go | `cloud.google.com/go/storage v1.65.1`, commit `3ab7d139…` | API GCS | Apache-2.0 | runtime | Go proxy + googleapis/google-cloud-go |
| Go | `1.26.7 windows/amd64`, archive SHA-256 `f4f534a4…` | compiler/test | BSD-3-Clause | build | go.dev |
| govulncheck | `golang.org/x/vuln/cmd/govulncheck v1.7.0`, commit `617f44b7…` | reachability scan | BSD-3-Clause | verification | go.googlesource.com/vuln |
| grafo transitivo | hashes completos en `go.sum` | SDK runtime | licencias de cada módulo | runtime | proxy/checksum DB oficial |

El build list publicado incluye `x/crypto/openpgp` a nivel módulo con `GO-2026-5932`, sin fix; `govulncheck` verificó que el paquete/símbolos no se importan ni alcanzan. `x/text` se elevó al fix oficial 0.39.0. `x/mod` aparece sólo en el grafo de módulos de tooling, no entre los 47 módulos runtime escaneados; Go 1.26.7 ya corrige la vulnerabilidad homóloga de `cmd/go`. Una actualización reinicia batch OSV y reachability.

## 8. Apply order

1. Completar blueprint, routing documental y profile Google; no materializar por auto-selección.
2. Adquirir/verificar Go, SDK/tag/commit/archive/licencia y grafo exactos.
3. Materializar los nueve archivos en destino vacío o sin colisiones.
4. Completar profile sin secretos y obtener approvals humanos; mantener `automatic_storage_authorized=false`.
5. Ejecutar verify/test/vet/build y scans. Crear un client oficial con ADC sólo en el composition root del proyecto.
6. Ejecutar sandbox GCP con bucket no productivo y probar owner, precondition, CRC, KMS, ambos modos aprobados, efecto parcial, audit logs, cuota/costo y cleanup permitido.
7. Habilitar una clase sólo después de seguridad, evaluación exacta, retention/legal y recuperación. Rollback: deshabilitar llamadas nuevas y reconciliar generaciones; nunca intentar borrar/aflojar objetos retenidos como rollback.

## 9. Verification

- `go mod verify` → todos los módulos verificados.
- `go test ./...` → 8 pruebas principales y subcasos de inputs, owner/config, precondition contract, post-write mismatch y receipt parcial.
- `go vet ./...` y `go build ./...` → exit 0 offline.
- `govulncheck v1.7.0 ./...` → 0 alcanzables, 0 en paquetes importados; un advisory de módulo no llamado retenido.
- OSV query batch del build list → denominador explícito; advisories/fixes y no-reachability quedan documentados, nunca ocultos.
- Materializador/compositor → manifest↔blocks↔SHA exactos; profile Google produce record sin colisión.
- Pendiente del proyecto: sandbox real, IAM/KMS/perimeter, cost/quota, audit, concurrency/load, malware, corpus, backup/recovery y retención legal. Sin esos gates, `CONDITIONED` permanece.

## 10. Reconstruction evidence

Reconstrucción limpia Windows 11 / PowerShell 7 / Go 1.26.7 el 2026-08-27. SDK oficial `v1.65.1`, tag/commit verificado `3ab7d139…`; ZIP proxy 998644 bytes SHA-256 `0030613d…`; licencia Apache-2.0 SHA-256 `cfc7749b…`. Los nueve archivos se materializan con SHA-256 individuales y ejecutan module verify, pruebas, vet, build y reachability scan sin credenciales ni cloud. Evidencia: `reconstruction_evidence/GO_GOOGLE_CLOUD_STORAGE_ADAPTER_2026-08-27_V1.md`. No hay live GCP ni autorización de almacenamiento inferida.
