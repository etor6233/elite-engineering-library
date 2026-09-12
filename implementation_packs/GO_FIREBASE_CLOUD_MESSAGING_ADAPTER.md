# Go Firebase Cloud Messaging Adapter

## 1. Metadata

```yaml
pack_id: "GO-FIREBASE-CLOUD-MESSAGING-ADAPTER"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un adaptador Go autónomo sobre Firebase Admin Go 4.21.0 para construir mensajes FCM, ejecutar dry-run oficial por defecto y permitir un envío real únicamente después de aprobaciones explícitas, con evidencia atómica sin tokens ni identificadores en claro."
stacks: ["Go 1.26.7", "Firebase Admin Go 4.21.0", "Firebase Cloud Messaging"]
compatible_with: ["OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "GO-PROVIDER-INTEGRATION-CORE 0.1.x"]
incompatible_with: ["credenciales o tokens dentro de Markdown/JSON/código", "broadcast no consentido", "envío real implícito", "persistencia automática de estado de negocio"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0"
upstream_sources: ["https://github.com/firebase/firebase-admin-go/tree/eebb06f2a643fbb59b1cb262874a943584475128", "https://pkg.go.dev/firebase.google.com/go/v4/messaging@v4.21.0"]
verified_at: "2026-08-26"
```

## 2. Applicability

Use cuando el proyecto haya seleccionado Firebase Cloud Messaging, demostrado un proyecto Firebase y necesite notificar a un dispositivo de prueba consentido desde Go. Adóptese sólo después de aceptar la licencia oficial, probar Application Default Credentials, habilitar FCM, aprobar contenido/datos, TTL, cuota/costo, retención y reconciliación. Rechácese para marketing masivo, topics, APNs/WebPush específicos, tokens no consentidos o escritura automática de estado de negocio. El claim cubre construcción del `messaging.Message`, `SendDryRun`/`Send`, fail-closed y recibo redactado; no demuestra entrega, visualización, opt-out, facturación ni producción.

## 3. Architecture contract

El CLI carga un perfil y una solicitud JSON con campos desconocidos prohibidos. El token se resuelve exclusivamente desde la variable de entorno nombrada en el perfil; Application Default Credentials queda a cargo del runtime oficial de Google. La frontera valida todas las aprobaciones, proyecto, modo, token, allowlist/reservas de data, contenido y TTL antes de construir el tipo oficial `messaging.Message`. `DRY_RUN` invoca `SendDryRun`; `SEND` exige además `real_send_approved=true`. Un error del proveedor o una respuesta sin message ID no deja directorio parcial. El recibo conserva hashes y no contiene proyecto, token ni message ID en claro. No escribe dominio automáticamente. Performance y cuotas deben medirse en el target aprobado. Rollback: retirar el wiring, revocar credenciales/tokens y conservar los recibos de reconciliación; toda actualización exige lock, licencia, tests y evidencia nuevos.

## 4. Exact file manifest

```text
CREATE firebase_push/go.mod
CREATE firebase_push/go.sum
CREATE firebase_push/sdk-source.lock.json
CREATE firebase_push/provider-profile.template.json
CREATE firebase_push/message.template.json
CREATE firebase_push/firebase_push.go
CREATE firebase_push/firebase_push_test.go
CREATE firebase_push/cmd/firebase-push/main.go
CREATE firebase_push/README.md
```

## 5. Materialization blocks

### FILE: `firebase_push/go.mod`
```yaml
block_id: "GO-FIREBASE-CLOUD-MESSAGING:gomod:v1"
operation: CREATE
provenance: AUTHORED
source: "dependency pin for official Firebase Admin Go"
license: "LicenseRef-Workspace-Owner"
sha256: "9543e144c5c213983232f03756ef73da52a40e21811112ae7f28496b8d2d0c2c"
variables: []
secrets_allowed: false
```
````go
module elite.local/firebasepush

go 1.26.0

require firebase.google.com/go/v4 v4.21.0

require (
	cel.dev/expr v0.25.2 // indirect
	cloud.google.com/go v0.123.0 // indirect
	cloud.google.com/go/auth v0.20.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.8 // indirect
	cloud.google.com/go/compute/metadata v0.9.0 // indirect
	cloud.google.com/go/firestore v1.22.0 // indirect
	cloud.google.com/go/iam v1.11.0 // indirect
	cloud.google.com/go/longrunning v1.0.0 // indirect
	cloud.google.com/go/monitoring v1.29.0 // indirect
	cloud.google.com/go/storage v1.62.1 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.32.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric v0.56.0 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.56.0 // indirect
	github.com/MicahParks/keyfunc v1.9.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cncf/xds/go v0.0.0-20260202195803-dba9d589def2 // indirect
	github.com/envoyproxy/go-control-plane/envoy v1.37.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.3.3 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-jose/go-jose/v4 v4.1.4 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.15 // indirect
	github.com/googleapis/gax-go/v2 v2.22.0 // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 // indirect
	github.com/spiffe/go-spiffe/v2 v2.6.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/detectors/gcp v1.43.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.68.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.68.0 // indirect
	go.opentelemetry.io/otel v1.43.0 // indirect
	go.opentelemetry.io/otel/metric v1.43.0 // indirect
	go.opentelemetry.io/otel/sdk v1.43.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.43.0 // indirect
	go.opentelemetry.io/otel/trace v1.43.0 // indirect
	golang.org/x/crypto v0.52.0 // indirect
	golang.org/x/net v0.55.0 // indirect
	golang.org/x/oauth2 v0.36.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	golang.org/x/time v0.15.0 // indirect
	google.golang.org/api v0.279.0 // indirect
	google.golang.org/appengine/v2 v2.0.6 // indirect
	google.golang.org/genproto v0.0.0-20260511170946-3700d4141b60 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260511170946-3700d4141b60 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260511170946-3700d4141b60 // indirect
	google.golang.org/grpc v1.81.1 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
````

### FILE: `firebase_push/go.sum`
```yaml
block_id: "GO-FIREBASE-CLOUD-MESSAGING:gosum:v1"
operation: CREATE
provenance: AUTHORED
source: "Go checksum database resolution for exact module graph"
license: "LicenseRef-Workspace-Owner"
sha256: "ad8189a073109a42337ef084ed5e2d7762c27e14e9cddc13d809ec9104a8516a"
variables: []
secrets_allowed: false
```
````text
cel.dev/expr v0.25.2 h1:K6j46C81hXtZQfuX60cVWQFBJahKSE2gfRbNuvr5bFs=
cel.dev/expr v0.25.2/go.mod h1:hrXvqGP6G6gyx8UAHSHJ5RGk//1Oj5nXQ2NI02Nrsg4=
cloud.google.com/go v0.123.0 h1:2NAUJwPR47q+E35uaJeYoNhuNEM9kM8SjgRgdeOJUSE=
cloud.google.com/go v0.123.0/go.mod h1:xBoMV08QcqUGuPW65Qfm1o9Y4zKZBpGS+7bImXLTAZU=
cloud.google.com/go/auth v0.20.0 h1:kXTssoVb4azsVDoUiF8KvxAqrsQcQtB53DcSgta74CA=
cloud.google.com/go/auth v0.20.0/go.mod h1:942/yi/itH1SsmpyrbnTMDgGfdy2BUqIKyd0cyYLc5Q=
cloud.google.com/go/auth/oauth2adapt v0.2.8 h1:keo8NaayQZ6wimpNSmW5OPc283g65QNIiLpZnkHRbnc=
cloud.google.com/go/auth/oauth2adapt v0.2.8/go.mod h1:XQ9y31RkqZCcwJWNSx2Xvric3RrU88hAYYbjDWYDL+c=
cloud.google.com/go/compute/metadata v0.9.0 h1:pDUj4QMoPejqq20dK0Pg2N4yG9zIkYGdBtwLoEkH9Zs=
cloud.google.com/go/compute/metadata v0.9.0/go.mod h1:E0bWwX5wTnLPedCKqk3pJmVgCBSM6qQI1yTBdEb3C10=
cloud.google.com/go/firestore v1.22.0 h1:avooeboIq37vKXobrbPUFhFBxS/c3FqmWoX0xs8dO6E=
cloud.google.com/go/firestore v1.22.0/go.mod h1:PaM4i7i7ruALSKmlpHXXZaPObcZw0W7ie5UOPr72iTU=
cloud.google.com/go/iam v1.11.0 h1:KieQ9Pb+LLPak1O3Rv3GgCxhnmkYf7Xyh0P5HfF1jFM=
cloud.google.com/go/iam v1.11.0/go.mod h1:KP+nKGugNJW4LcLx1uEZcq1ok5sQHFaQehQNl4QDgV4=
cloud.google.com/go/logging v1.18.0 h1:KhzZq+1cSkPH9YUaKLLhLtQxIHitVayBmk0sGfoM9+k=
cloud.google.com/go/logging v1.18.0/go.mod h1:ZGKnpBaURITh+g/uom2VhbiFoFWvejcrHPDhxFtU/gI=
cloud.google.com/go/longrunning v1.0.0 h1:lwzWEYD8+NkYV7dhexOz6kmlvajZA70+bW/xMhRVVdY=
cloud.google.com/go/longrunning v1.0.0/go.mod h1:8nqFBPOO1U/XkhWl0I19AMZEphrHi73VNABIpKYaTwM=
cloud.google.com/go/monitoring v1.29.0 h1:AHhDsFaSax1/4k+qlIDX/SDGe6hggnfXJ9dkgD9qBPY=
cloud.google.com/go/monitoring v1.29.0/go.mod h1:72NOVjJXHY/HBfoLT0+qlCZBT059+9VXLeAnL2PeeVM=
cloud.google.com/go/storage v1.62.1 h1:Os0G3XbUbjZumkpDUf2Y0rLoXJTCF1kU2kWUujKYXD8=
cloud.google.com/go/storage v1.62.1/go.mod h1:cpYz/kRVZ+UQAF1uHeea10/9ewcRbxGoGNKsS9daSXA=
cloud.google.com/go/trace v1.16.0 h1:GmQovzFc5F0CNfl0VLgL64aoTtu7xsM0YajW2GlG9+E=
cloud.google.com/go/trace v1.16.0/go.mod h1:r+bdAn16dKLSV1G2D5v3e58IlQlizfxWrUfjx7kM7X0=
firebase.google.com/go/v4 v4.21.0 h1:HBZV4jrLtFYj8EwWyqEZOuRLfkfkV2bpnfyyXHOhPxY=
firebase.google.com/go/v4 v4.21.0/go.mod h1:CDumIdA5oTiyDpLNVcQoW8ZrB5CTgyE2D45DuENIABg=
github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.32.0 h1:rIkQfkCOVKc1OiRCNcSDD8ml5RJlZbH/Xsq7lbpynwc=
github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.32.0/go.mod h1:RD2SsorTmYhF6HkTmDw7KmPYQk8OBYwTkuasChwv7R4=
github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric v0.56.0 h1:O2sXMyJh8b7devAGdE+163xtRurt0RVpB6DIzX5vGfg=
github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric v0.56.0/go.mod h1:hEpiGU18xf70qb3jbTcIggWAiEfX/cOIVc2OTe4OegA=
github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/cloudmock v0.56.0 h1:ZIT85vKP7LBS84XJ0WdJ3dPOX3iz4j3c0+lpajGQMyo=
github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/cloudmock v0.56.0/go.mod h1:rqP9UEhOXv9WhQ7Gjz+G5y/pf8+BJZW5/Ts0AhE0PwE=
github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.56.0 h1:0YP0+/ixwu+Uqeu/FGiBZNQ19huiUxxiPXIc9WsLKuQ=
github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.56.0/go.mod h1:6ZZMQhZKDvUvkJw2rc+oDP90tMMzuU/J+5HG1ZmPOmE=
github.com/MicahParks/keyfunc v1.9.0 h1:lhKd5xrFHLNOWrDc4Tyb/Q1AJ4LCzQ48GVJyVIID3+o=
github.com/MicahParks/keyfunc v1.9.0/go.mod h1:IdnCilugA0O/99dW+/MkvlyrsX8+L8+x95xuVNtM5jw=
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
github.com/golang-jwt/jwt/v4 v4.4.2/go.mod h1:m21LjoU+eqJr34lmDMbreY2eSTRJ1cv77w39/MY0Ch0=
github.com/golang-jwt/jwt/v4 v4.5.2 h1:YtQM7lnr8iZ+j5q71MGKkNw9Mn7AjHM68uc9g5fXeUI=
github.com/golang-jwt/jwt/v4 v4.5.2/go.mod h1:m21LjoU+eqJr34lmDMbreY2eSTRJ1cv77w39/MY0Ch0=
github.com/golang/protobuf v1.5.0/go.mod h1:FsONVRAS9T7sI+LIUmWTfcYkHO4aIWwzhcaSAoJOfIk=
github.com/golang/protobuf v1.5.4 h1:i7eJL8qZTpSEXOPTxNKhASYpMn+8e5Q6AdndVa1dWek=
github.com/golang/protobuf v1.5.4/go.mod h1:lnTiLA8Wa4RWRcIUkrtSVa5nRhsEGBg48fD6rSs7xps=
github.com/google/go-cmp v0.5.5/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/gNBxE=
github.com/google/go-cmp v0.7.0 h1:wk8382ETsv4JYUZwIsn6YpYiWiBsYLSJiTsyBybVuN8=
github.com/google/go-cmp v0.7.0/go.mod h1:pXiqmnSA92OHEEa9HXL2W4E7lf9JzCmGVUdgjX3N/iU=
github.com/google/martian/v3 v3.3.3 h1:DIhPTQrbPkgs2yJYdXU/eNACCG5DVQjySNRNlflZ9Fc=
github.com/google/martian/v3 v3.3.3/go.mod h1:iEPrYcgCF7jA9OtScMFQyAlZZ4YXTKEtJ1E6RWzmBA0=
github.com/google/s2a-go v0.1.9 h1:LGD7gtMgezd8a/Xak7mEWL0PjoTQFvpRudN895yqKW0=
github.com/google/s2a-go v0.1.9/go.mod h1:YA0Ei2ZQL3acow2O62kdp9UlnvMmU7kA6Eutn0dXayM=
github.com/google/uuid v1.6.0 h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=
github.com/google/uuid v1.6.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
github.com/googleapis/enterprise-certificate-proxy v0.3.15 h1:xolVQTEXusUcAA5UgtyRLjelpFFHWlPQ4XfWGc7MBas=
github.com/googleapis/enterprise-certificate-proxy v0.3.15/go.mod h1:vqVt9yG9480NtzREnTlmGSBmFrA+bzb0yl0TxoBQXOg=
github.com/googleapis/gax-go/v2 v2.22.0 h1:PjIWBpgGIVKGoCXuiCoP64altEJCj3/Ei+kSU5vlZD4=
github.com/googleapis/gax-go/v2 v2.22.0/go.mod h1:irWBbALSr0Sk3qlqb9SyJ1h68WjgeFuiOzI4Rqw5+aY=
github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10 h1:GFCKgmp0tecUJ0sJuv4pzYCqS9+RGSn52M3FUwPs+uo=
github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10/go.mod h1:t/avpk3KcrXxUnYOhZhMXJlSEyie6gQbtLq5NM3loB8=
github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 h1:Jamvg5psRIccs7FGNTlIRMkT8wgtp5eCXdBlqhYGL6U=
github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
github.com/spiffe/go-spiffe/v2 v2.6.0 h1:l+DolpxNWYgruGQVV0xsfeya3CsC7m8iBzDnMpsbLuo=
github.com/spiffe/go-spiffe/v2 v2.6.0/go.mod h1:gm2SeUoMZEtpnzPNs2Csc0D/gX33k1xIx7lEzqblHEs=
github.com/stretchr/testify v1.11.1 h1:7s2iGBzp5EwR7/aIZr8ao5+dra3wiQyKjjFuvgVKu7U=
github.com/stretchr/testify v1.11.1/go.mod h1:wZwfW3scLgRK+23gO65QZefKpKQRnfz6sD981Nm4B6U=
github.com/yuin/goldmark v1.4.13/go.mod h1:6yULJ656Px+3vBD8DxQVa3kxgyrAnzto9xy5taEt/CY=
go.opentelemetry.io/auto/sdk v1.2.1 h1:jXsnJ4Lmnqd11kwkBV2LgLoFMZKizbCi5fNZ/ipaZ64=
go.opentelemetry.io/auto/sdk v1.2.1/go.mod h1:KRTj+aOaElaLi+wW1kO/DZRXwkF4C5xPbEe3ZiIhN7Y=
go.opentelemetry.io/contrib/detectors/gcp v1.43.0 h1:62yY3dT7/ShwOxzA0RsKRgshBmfElKI4d/Myu2OxDFU=
go.opentelemetry.io/contrib/detectors/gcp v1.43.0/go.mod h1:RyaZMFY7yi1kAs45S6mbFGz8O8rqB0dTY14uzvG4LCs=
go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.68.0 h1:0Qx7VGBacMm9ZENQ7TnNObTYI4ShC+lHI16seduaxZo=
go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.68.0/go.mod h1:Sje3i3MjSPKTSPvVWCaL8ugBzJwik3u4smCjUeuupqg=
go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.68.0 h1:CqXxU8VOmDefoh0+ztfGaymYbhdB/tT3zs79QaZTNGY=
go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.68.0/go.mod h1:BuhAPThV8PBHBvg8ZzZ/Ok3idOdhWIodywz2xEcRbJo=
go.opentelemetry.io/otel v1.43.0 h1:mYIM03dnh5zfN7HautFE4ieIig9amkNANT+xcVxAj9I=
go.opentelemetry.io/otel v1.43.0/go.mod h1:JuG+u74mvjvcm8vj8pI5XiHy1zDeoCS2LB1spIq7Ay0=
go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.43.0 h1:TC+BewnDpeiAmcscXbGMfxkO+mwYUwE/VySwvw88PfA=
go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.43.0/go.mod h1:J/ZyF4vfPwsSr9xJSPyQ4LqtcTPULFR64KwTikGLe+A=
go.opentelemetry.io/otel/metric v1.43.0 h1:d7638QeInOnuwOONPp4JAOGfbCEpYb+K6DVWvdxGzgM=
go.opentelemetry.io/otel/metric v1.43.0/go.mod h1:RDnPtIxvqlgO8GRW18W6Z/4P462ldprJtfxHxyKd2PY=
go.opentelemetry.io/otel/sdk v1.43.0 h1:pi5mE86i5rTeLXqoF/hhiBtUNcrAGHLKQdhg4h4V9Dg=
go.opentelemetry.io/otel/sdk v1.43.0/go.mod h1:P+IkVU3iWukmiit/Yf9AWvpyRDlUeBaRg6Y+C58QHzg=
go.opentelemetry.io/otel/sdk/metric v1.43.0 h1:S88dyqXjJkuBNLeMcVPRFXpRw2fuwdvfCGLEo89fDkw=
go.opentelemetry.io/otel/sdk/metric v1.43.0/go.mod h1:C/RJtwSEJ5hzTiUz5pXF1kILHStzb9zFlIEe85bhj6A=
go.opentelemetry.io/otel/trace v1.43.0 h1:BkNrHpup+4k4w+ZZ86CZoHHEkohws8AY+WTX09nk+3A=
go.opentelemetry.io/otel/trace v1.43.0/go.mod h1:/QJhyVBUUswCphDVxq+8mld+AvhXZLhe+8WVFxiFff0=
golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2/go.mod h1:djNgcEr1/C05ACkg1iLfiJU5Ep61QUkGW8qpdssI0+w=
golang.org/x/crypto v0.0.0-20210921155107-089bfa567519/go.mod h1:GvvjBRRGRdwPK5ydBHafDWAxML/pGHZbMvKqRZ5+Abc=
golang.org/x/crypto v0.52.0 h1:RMs7fP2rXdep0CftQlK8Uf+kibLm7qkCcradZWYz988=
golang.org/x/crypto v0.52.0/go.mod h1:1QgfPxDqh0T2M/elOJtp9RvuR95kVjir0e6/BvEmGbc=
golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4/go.mod h1:jJ57K6gSWd91VN4djpZkiMVwK6gcyfeH4XE8wZrZaV4=
golang.org/x/net v0.0.0-20190620200207-3b0461eec859/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
golang.org/x/net v0.0.0-20210226172049-e18ecbb05110/go.mod h1:m0MpNAwzfU5UDzcl9v0D8zg8gWTRqZa9RBIspLL5mdg=
golang.org/x/net v0.0.0-20220722155237-a158d28d115b/go.mod h1:XRhObCWvk6IyKnWLug+ECip1KBveYUHfp+8e9klMJ9c=
golang.org/x/net v0.55.0 h1:bcvxaJn3e1U6InsFWt1JUq1aSjnRxLzT2rtD2KfkDF8=
golang.org/x/net v0.55.0/go.mod h1:L5U2KuzuOe1lY7Z+aWVIKK6qEeJXnXV9yzGA+WCHJww=
golang.org/x/oauth2 v0.36.0 h1:peZ/1z27fi9hUOFCAZaHyrpWG5lwe0RJEEEeH0ThlIs=
golang.org/x/oauth2 v0.36.0/go.mod h1:YDBUJMTkDnJS+A4BP4eZBjCqtokkg1hODuPjwiGPO7Q=
golang.org/x/sync v0.0.0-20190423024810-112230192c58/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
golang.org/x/sync v0.20.0 h1:e0PTpb7pjO8GAtTs2dQ6jYa5BWYlMuX047Dco/pItO4=
golang.org/x/sync v0.20.0/go.mod h1:9xrNwdLfx4jkKbNva9FpL6vEN7evnE43NNNJQ2LF3+0=
golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
golang.org/x/sys v0.0.0-20201119102817-f84b799fce68/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
golang.org/x/sys v0.45.0 h1:dO4czNzziLiiXplLQgBCEpCvXQ3dnkn0SdaZSYdQ+FY=
golang.org/x/sys v0.45.0/go.mod h1:4GL1E5IUh+htKOUEOaiffhrAeqysfVGipDYzABqnCmw=
golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1/go.mod h1:bj7SfCRtBDWHUb9snDiAeCFNEtKQo2Wmx5Cou7ajbmo=
golang.org/x/term v0.0.0-20210927222741-03fcf44c2211/go.mod h1:jbD1KX2456YbFQfuXm/mYQcufACuNUgVhRMnK/tPxf8=
golang.org/x/text v0.3.0/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
golang.org/x/text v0.3.3/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
golang.org/x/text v0.3.7/go.mod h1:u+2+/6zg+i71rQMx5EYifcz6MCKuco9NR6JIITiCfzQ=
golang.org/x/text v0.3.8/go.mod h1:E6s5w1FMmriuDzIBO73fBruAKo1PCIq6d2Q6DHfQ8WQ=
golang.org/x/text v0.37.0 h1:Cqjiwd9eSg8e0QAkyCaQTNHFIIzWtidPahFWR83rTrc=
golang.org/x/text v0.37.0/go.mod h1:a5sjxXGs9hsn/AJVwuElvCAo9v8QYLzvavO5z2PiM38=
golang.org/x/time v0.15.0 h1:bbrp8t3bGUeFOx08pvsMYRTCVSMk89u4tKbNOZbp88U=
golang.org/x/time v0.15.0/go.mod h1:Y4YMaQmXwGQZoFaVFk4YpCt4FLQMYKZe9oeV/f4MSno=
golang.org/x/tools v0.0.0-20180917221912-90fa682c2a6e/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
golang.org/x/tools v0.0.0-20191119224855-298f0cb1881e/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
golang.org/x/tools v0.1.12/go.mod h1:hNGJHUnrk76NpqgfD5Aqm5Crs+Hm0VOH/i9J2+nxYbc=
golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQqLJ2OPfmY0=
golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQqLJ2OPfmY0=
gonum.org/v1/gonum v0.17.0 h1:VbpOemQlsSMrYmn7T2OUvQ4dqxQXU+ouZFQsZOx50z4=
gonum.org/v1/gonum v0.17.0/go.mod h1:El3tOrEuMpv2UdMrbNlKEh9vd86bmQ6vqIcDwxEOc1E=
google.golang.org/api v0.279.0 h1:hsx2M2OaRcaKtVYK6vXEUnQvdjnend7ZYES+lYaot74=
google.golang.org/api v0.279.0/go.mod h1:B9TqLBwJqVjp1mtt7WeoQwWRwvu/400y5lETOql+giQ=
google.golang.org/appengine/v2 v2.0.6 h1:LvPZLGuchSBslPBp+LAhihBeGSiRh1myRoYK4NtuBIw=
google.golang.org/appengine/v2 v2.0.6/go.mod h1:WoEXGoXNfa0mLvaH5sV3ZSGXwVmy8yf7Z1JKf3J3wLI=
google.golang.org/genproto v0.0.0-20260511170946-3700d4141b60 h1:rhBdfmsOlOZIvz3Y5/BdUzPg2CkO8L7QQPKj96B8554=
google.golang.org/genproto v0.0.0-20260511170946-3700d4141b60/go.mod h1:8xo2Pj1b20ZOCpzlU3B9qieMwVIAXx1QVZWLMlPL6sM=
google.golang.org/genproto/googleapis/api v0.0.0-20260511170946-3700d4141b60 h1:3WsB1FAbiRIf2tOxscWKs3pQBD9he1NsrnbhMuWfekc=
google.golang.org/genproto/googleapis/api v0.0.0-20260511170946-3700d4141b60/go.mod h1:7yoXV7RIh5gblj/xVYoogxAWvA9wUeVbpsK/M694l00=
google.golang.org/genproto/googleapis/rpc v0.0.0-20260511170946-3700d4141b60 h1:seT2EwLWM78plQ7wcDfuWBc/4FAEAXDDiaSol4ku4qo=
google.golang.org/genproto/googleapis/rpc v0.0.0-20260511170946-3700d4141b60/go.mod h1:4Hqkh8ycfw05ld/3BWL7rJOSfebL2Q+DVDeRgYgxUU8=
google.golang.org/grpc v1.81.1 h1:VnnIIZ88UzOOKLukQi+ImGz8O1Wdp8nAGGnvOfEIWQQ=
google.golang.org/grpc v1.81.1/go.mod h1:xGH9GfzOyMTGIOXBJmXt+BX/V0kcdQbdcuwQ/zNw42I=
google.golang.org/protobuf v1.26.0-rc.1/go.mod h1:jlhhOSvTdKEhbULTjvd4ARK9grFBp09yW+WbY/TyQbw=
google.golang.org/protobuf v1.30.0/go.mod h1:HV8QOd/L58Z+nl8r43ehVNZIU/HEI6OcFqwMG9pJV4I=
google.golang.org/protobuf v1.36.11 h1:fV6ZwhNocDyBLK0dj+fg8ektcVegBBuEolpbTQyBNVE=
google.golang.org/protobuf v1.36.11/go.mod h1:HTf+CrKn2C3g5S8VImy6tdcUvCska2kB7j23XfzDpco=
gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
````

### FILE: `firebase_push/sdk-source.lock.json`
```yaml
block_id: "GO-FIREBASE-CLOUD-MESSAGING:sourcelock:v1"
operation: CREATE
provenance: AUTHORED
source: "official Firebase Admin Go release, commit, archive, license and critical-file audit"
license: "LicenseRef-Workspace-Owner"
sha256: "0b2dabefdbcf1e044881ebd24303223c553cc205b29ee9ad420fe942ca3a1f50"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-go-sdk-source-lock/v1",
  "provider": "Google Firebase",
  "module": "firebase.google.com/go/v4",
  "version": "v4.21.0",
  "repository": "firebase/firebase-admin-go",
  "commit": "eebb06f2a643fbb59b1cb262874a943584475128",
  "commit_signature_verified": false,
  "archive_bytes": 322131,
  "archive_sha256": "8d7bbc5621ddf19ddf0ebd74df867690cd3ee1bed7499eabf973eba668284b5e",
  "license_expression": "Apache-2.0",
  "license_sha256": "c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4",
  "upstream_go_mod_sha256": "a4baa15771d86ad3f97e729437a5971d168f138440fe742dad595c9d728a1db1",
  "upstream_go_sum_sha256": "9cd89415c44406855db8406b673cef99ca8964d630bd77f052e1071aac33e925",
  "messaging_go_sha256": "b80b2ab57062a50dce0cf3060ae8857e7b9145a1661895a146277705c30210d9",
  "verified_at": "2026-08-26"
}
````

### FILE: `firebase_push/provider-profile.template.json`
```yaml
block_id: "GO-FIREBASE-CLOUD-MESSAGING:profile:v1"
operation: CREATE
provenance: AUTHORED
source: "fail-closed project/provider approval profile"
license: "LicenseRef-Workspace-Owner"
sha256: "fe862ef2788ed1cb368159e546dcba585c3a98246656f0c525d2f6311e8922fc"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-firebase-push-profile/v1",
  "provider": "Google Firebase Cloud Messaging",
  "sdk_module": "firebase.google.com/go/v4",
  "sdk_version": "v4.21.0",
  "source_commit": "eebb06f2a643fbb59b1cb262874a943584475128",
  "decision": "BLOCKED_PROJECT_ACCESS_CONSENT_AND_MESSAGE_POLICY_REQUIRED",
  "official_source_license_accepted": false,
  "firebase_project_proven": false,
  "application_default_credentials_proven": false,
  "fcm_api_enabled_proven": false,
  "test_device_consent_proven": false,
  "message_policy_approved": false,
  "quota_and_cost_approved": false,
  "reconciliation_approved": false,
  "project_id": "",
  "mode": "DRY_RUN",
  "real_send_approved": false,
  "device_token_environment_variable": "FIREBASE_TEST_DEVICE_TOKEN",
  "approved_data_keys": [],
  "notification_content_approved": false,
  "maximum_ttl_seconds": 0,
  "data_retention": "",
  "automatic_business_write": false
}
````

### FILE: `firebase_push/message.template.json`
```yaml
block_id: "GO-FIREBASE-CLOUD-MESSAGING:message:v1"
operation: CREATE
provenance: AUTHORED
source: "minimal FCM message request template"
license: "LicenseRef-Workspace-Owner"
sha256: "18be24fdff42998690cae145990c35422bc3ae3bf3e9883a4710ef604e575d8c"
variables: []
secrets_allowed: false
```
````json
{
  "notification": null,
  "data": {},
  "android_ttl_seconds": 0
}
````

### FILE: `firebase_push/firebase_push.go`
```yaml
block_id: "GO-FIREBASE-CLOUD-MESSAGING:adapter:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed adapter invoking public Firebase Admin Go messaging APIs"
license: "LicenseRef-Workspace-Owner"
sha256: "4a6c7c65e5c272e65a37f37f85a633cc9311a323448e294730fec60c1cf34ffe"
variables: []
secrets_allowed: false
```
````go
package firebasepush

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"firebase.google.com/go/v4/messaging"
)

const (
	SDKVersion   = "v4.21.0"
	SourceCommit = "eebb06f2a643fbb59b1cb262874a943584475128"
)

type ApprovedTemplate struct{}

type Profile struct {
	Schema                              string   `json:"schema"`
	Provider                            string   `json:"provider"`
	SDKModule                           string   `json:"sdk_module"`
	SDKVersion                          string   `json:"sdk_version"`
	SourceCommit                        string   `json:"source_commit"`
	Decision                            string   `json:"decision"`
	OfficialSourceLicenseAccepted       bool     `json:"official_source_license_accepted"`
	FirebaseProjectProven               bool     `json:"firebase_project_proven"`
	ApplicationDefaultCredentialsProven bool     `json:"application_default_credentials_proven"`
	FCMAPIEnabledProven                 bool     `json:"fcm_api_enabled_proven"`
	TestDeviceConsentProven             bool     `json:"test_device_consent_proven"`
	MessagePolicyApproved               bool     `json:"message_policy_approved"`
	QuotaAndCostApproved                bool     `json:"quota_and_cost_approved"`
	ReconciliationApproved              bool     `json:"reconciliation_approved"`
	ProjectID                           string   `json:"project_id"`
	Mode                                string   `json:"mode"`
	RealSendApproved                    bool     `json:"real_send_approved"`
	DeviceTokenEnvironmentVariable      string   `json:"device_token_environment_variable"`
	ApprovedDataKeys                    []string `json:"approved_data_keys"`
	NotificationContentApproved         bool     `json:"notification_content_approved"`
	MaximumTTLSeconds                   int64    `json:"maximum_ttl_seconds"`
	DataRetention                       string   `json:"data_retention"`
	AutomaticBusinessWrite              bool     `json:"automatic_business_write"`
}

type Notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type Request struct {
	Notification      *Notification     `json:"notification"`
	Data              map[string]string `json:"data"`
	AndroidTTLSeconds int64             `json:"android_ttl_seconds"`
}

type Receipt struct {
	Schema               string `json:"schema"`
	CreatedAt            string `json:"created_at"`
	SDKVersion           string `json:"sdk_version"`
	SourceCommit         string `json:"source_commit"`
	Mode                 string `json:"mode"`
	ProjectIDSHA256      string `json:"project_id_sha256"`
	DeviceTokenSHA256    string `json:"device_token_sha256"`
	RequestSHA256        string `json:"request_sha256"`
	ResponseSHA256       string `json:"response_sha256"`
	MessageIDSHA256      string `json:"message_id_sha256"`
	AutomaticDomainWrite bool   `json:"automatic_domain_write"`
}

type Client interface {
	Send(context.Context, *messaging.Message) (string, error)
	SendDryRun(context.Context, *messaging.Message) (string, error)
}

var (
	projectPattern = regexp.MustCompile(`^[a-z][a-z0-9-]{4,61}[a-z0-9]$`)
	envPattern     = regexp.MustCompile(`^[A-Z][A-Z0-9_]{2,80}$`)
	keyPattern     = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_.-]{0,127}$`)
)

func hash(value []byte) string {
	sum := sha256.Sum256(value)
	return hex.EncodeToString(sum[:])
}

func LoadJSON(path string, target any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return errors.New("JSON must contain exactly one value")
	}
	return nil
}

func Validate(profile Profile, request Request, token string) error {
	if profile.Schema != "elite-firebase-push-profile/v1" || profile.Provider != "Google Firebase Cloud Messaging" || profile.SDKModule != "firebase.google.com/go/v4" || profile.SDKVersion != SDKVersion || profile.SourceCommit != SourceCommit {
		return errors.New("Firebase profile authority mismatch")
	}
	if profile.Decision != "PROVEN" || !profile.OfficialSourceLicenseAccepted || !profile.FirebaseProjectProven || !profile.ApplicationDefaultCredentialsProven || !profile.FCMAPIEnabledProven || !profile.TestDeviceConsentProven || !profile.MessagePolicyApproved || !profile.QuotaAndCostApproved || !profile.ReconciliationApproved {
		return errors.New("Firebase project, access, consent or message policy remains blocked")
	}
	if profile.AutomaticBusinessWrite || strings.TrimSpace(profile.DataRetention) == "" {
		return errors.New("retention is required and automatic_business_write must remain false")
	}
	if !projectPattern.MatchString(profile.ProjectID) || !envPattern.MatchString(profile.DeviceTokenEnvironmentVariable) {
		return errors.New("project_id or token environment reference is invalid")
	}
	if profile.Mode != "DRY_RUN" && profile.Mode != "SEND" {
		return errors.New("mode must be DRY_RUN or SEND")
	}
	if profile.Mode == "SEND" && !profile.RealSendApproved {
		return errors.New("real Firebase send is not approved")
	}
	if profile.Mode == "DRY_RUN" && profile.RealSendApproved {
		return errors.New("real_send_approved must remain false in DRY_RUN")
	}
	if len(token) < 20 || len(token) > 4096 || strings.IndexFunc(token, func(r rune) bool { return r <= ' ' }) >= 0 {
		return errors.New("device token is missing or malformed")
	}
	if profile.MaximumTTLSeconds < 0 || profile.MaximumTTLSeconds > 2419200 || request.AndroidTTLSeconds < 0 || request.AndroidTTLSeconds > profile.MaximumTTLSeconds {
		return errors.New("Android TTL exceeds approved bounds")
	}
	approved := make(map[string]struct{}, len(profile.ApprovedDataKeys))
	for _, key := range profile.ApprovedDataKeys {
		if !keyPattern.MatchString(key) || reservedKey(key) {
			return errors.New("approved_data_keys contains an invalid or reserved key")
		}
		if _, exists := approved[key]; exists {
			return errors.New("approved_data_keys contains a duplicate")
		}
		approved[key] = struct{}{}
	}
	for key, value := range request.Data {
		if _, ok := approved[key]; !ok {
			return fmt.Errorf("data key is not approved: %s", key)
		}
		if len(value) > 4096 {
			return fmt.Errorf("data value is too large: %s", key)
		}
	}
	if request.Notification != nil {
		if !profile.NotificationContentApproved {
			return errors.New("notification content is not approved")
		}
		if strings.TrimSpace(request.Notification.Title) == "" || strings.TrimSpace(request.Notification.Body) == "" || len(request.Notification.Title) > 200 || len(request.Notification.Body) > 2000 {
			return errors.New("notification title/body is missing or too large")
		}
	}
	if request.Notification == nil && len(request.Data) == 0 {
		return errors.New("message must contain approved notification or data")
	}
	return nil
}

func reservedKey(key string) bool {
	lower := strings.ToLower(key)
	return lower == "from" || strings.HasPrefix(lower, "google.") || strings.HasPrefix(lower, "gcm.")
}

func BuildMessage(request Request, token string) *messaging.Message {
	message := &messaging.Message{Token: token, Data: request.Data}
	if request.Notification != nil {
		message.Notification = &messaging.Notification{Title: request.Notification.Title, Body: request.Notification.Body}
	}
	if request.AndroidTTLSeconds > 0 {
		ttl := time.Duration(request.AndroidTTLSeconds) * time.Second
		message.Android = &messaging.AndroidConfig{TTL: &ttl}
	}
	return message
}

func Execute(ctx context.Context, client Client, profile Profile, request Request, token, outputDirectory string) (Receipt, error) {
	if err := Validate(profile, request, token); err != nil {
		return Receipt{}, err
	}
	abs, err := filepath.Abs(outputDirectory)
	if err != nil {
		return Receipt{}, err
	}
	if _, err := os.Stat(abs); !os.IsNotExist(err) {
		if err == nil {
			return Receipt{}, errors.New("output directory must not exist")
		}
		return Receipt{}, err
	}
	parent := filepath.Dir(abs)
	if err := os.MkdirAll(parent, 0o700); err != nil {
		return Receipt{}, err
	}
	stage, err := os.MkdirTemp(parent, ".firebase-push-stage-")
	if err != nil {
		return Receipt{}, err
	}
	defer os.RemoveAll(stage)

	message := BuildMessage(request, token)
	var messageID string
	if profile.Mode == "DRY_RUN" {
		messageID, err = client.SendDryRun(ctx, message)
	} else {
		messageID, err = client.Send(ctx, message)
	}
	if err != nil {
		return Receipt{}, fmt.Errorf("Firebase provider failed: %w", err)
	}
	if strings.TrimSpace(messageID) == "" {
		return Receipt{}, errors.New("Firebase response is missing message id")
	}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return Receipt{}, err
	}
	responseBytes, err := json.Marshal(map[string]string{"message_id": messageID})
	if err != nil {
		return Receipt{}, err
	}
	receipt := Receipt{
		Schema: "elite-firebase-push-receipt/v1", CreatedAt: time.Now().UTC().Format(time.RFC3339Nano),
		SDKVersion: SDKVersion, SourceCommit: SourceCommit, Mode: profile.Mode,
		ProjectIDSHA256: hash([]byte(profile.ProjectID)), DeviceTokenSHA256: hash([]byte(token)),
		RequestSHA256: hash(requestBytes), ResponseSHA256: hash(responseBytes), MessageIDSHA256: hash([]byte(messageID)), AutomaticDomainWrite: false,
	}
	formattedResponse := append(responseBytes, '\n')
	receiptBytes, err := json.MarshalIndent(receipt, "", "  ")
	if err != nil {
		return Receipt{}, err
	}
	receiptBytes = append(receiptBytes, '\n')
	if err := os.WriteFile(filepath.Join(stage, "provider-response.json"), formattedResponse, 0o600); err != nil {
		return Receipt{}, err
	}
	if err := os.WriteFile(filepath.Join(stage, "PUSH_RECEIPT.json"), receiptBytes, 0o600); err != nil {
		return Receipt{}, err
	}
	if err := os.Rename(stage, abs); err != nil {
		return Receipt{}, err
	}
	return receipt, nil
}

func SortedDataKeys(data map[string]string) []string {
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
````

### FILE: `firebase_push/firebase_push_test.go`
```yaml
block_id: "GO-FIREBASE-CLOUD-MESSAGING:tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local regression and negative tests against official messaging types"
license: "LicenseRef-Workspace-Owner"
sha256: "b6c143462a7cc724942d9aff81587ba4498c680dd9d4cd91635efb1d4e0fbf81"
variables: []
secrets_allowed: false
```
````go
package firebasepush

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"firebase.google.com/go/v4/messaging"
)

type fakeClient struct {
	sent, dry int
	message   *messaging.Message
	id        string
	err       error
}

func (f *fakeClient) Send(_ context.Context, message *messaging.Message) (string, error) {
	f.sent++
	f.message = message
	return f.id, f.err
}
func (f *fakeClient) SendDryRun(_ context.Context, message *messaging.Message) (string, error) {
	f.dry++
	f.message = message
	return f.id, f.err
}

func provenProfile() Profile {
	return Profile{Schema: "elite-firebase-push-profile/v1", Provider: "Google Firebase Cloud Messaging", SDKModule: "firebase.google.com/go/v4", SDKVersion: SDKVersion, SourceCommit: SourceCommit, Decision: "PROVEN", OfficialSourceLicenseAccepted: true, FirebaseProjectProven: true, ApplicationDefaultCredentialsProven: true, FCMAPIEnabledProven: true, TestDeviceConsentProven: true, MessagePolicyApproved: true, QuotaAndCostApproved: true, ReconciliationApproved: true, ProjectID: "elite-test-project", Mode: "DRY_RUN", DeviceTokenEnvironmentVariable: "FIREBASE_TEST_DEVICE_TOKEN", ApprovedDataKeys: []string{"order_id", "status"}, NotificationContentApproved: true, MaximumTTLSeconds: 3600, DataRetention: "30 days", AutomaticBusinessWrite: false}
}

func validRequest() Request {
	return Request{Notification: &Notification{Title: "Order update", Body: "Your order changed"}, Data: map[string]string{"order_id": "A-1", "status": "ready"}, AndroidTTLSeconds: 600}
}

const testToken = "verified-test-device-token-value-123456789"

func TestDryRunUsesOfficialMessageAndRedactsReceipt(t *testing.T) {
	client := &fakeClient{id: "projects/elite/messages/secret-provider-id"}
	output := filepath.Join(t.TempDir(), "evidence")
	receipt, err := Execute(context.Background(), client, provenProfile(), validRequest(), testToken, output)
	if err != nil {
		t.Fatal(err)
	}
	if client.dry != 1 || client.sent != 0 || client.message.Token != testToken || client.message.Notification.Title != "Order update" {
		t.Fatalf("unexpected official message call: %#v", client)
	}
	serialized, _ := json.Marshal(receipt)
	if strings.Contains(string(serialized), testToken) || strings.Contains(string(serialized), "elite-test-project") || strings.Contains(string(serialized), "secret-provider-id") {
		t.Fatal("receipt leaked an identifier")
	}
	if _, err := os.Stat(filepath.Join(output, "PUSH_RECEIPT.json")); err != nil {
		t.Fatal(err)
	}
}

func TestRealSendRequiresExplicitApproval(t *testing.T) {
	profile := provenProfile()
	profile.Mode = "SEND"
	if err := Validate(profile, validRequest(), testToken); err == nil || !strings.Contains(err.Error(), "not approved") {
		t.Fatalf("expected real-send block, got %v", err)
	}
	profile.RealSendApproved = true
	client := &fakeClient{id: "message-id"}
	if _, err := Execute(context.Background(), client, profile, validRequest(), testToken, filepath.Join(t.TempDir(), "evidence")); err != nil {
		t.Fatal(err)
	}
	if client.sent != 1 || client.dry != 0 {
		t.Fatalf("wrong send mode: %#v", client)
	}
}

func TestProfileAndMessageFailClosed(t *testing.T) {
	profile := provenProfile()
	profile.Decision = "BLOCKED"
	if err := Validate(profile, validRequest(), testToken); err == nil {
		t.Fatal("blocked profile passed")
	}
	profile = provenProfile()
	request := validRequest()
	request.Data["unapproved"] = "x"
	if err := Validate(profile, request, testToken); err == nil {
		t.Fatal("unapproved key passed")
	}
	profile = provenProfile()
	profile.ApprovedDataKeys = []string{"google.reserved"}
	if err := Validate(profile, validRequest(), testToken); err == nil {
		t.Fatal("reserved key passed")
	}
}

func TestProviderFailureIsAtomic(t *testing.T) {
	parent := t.TempDir()
	output := filepath.Join(parent, "evidence")
	_, err := Execute(context.Background(), &fakeClient{err: errors.New("provider unavailable")}, provenProfile(), validRequest(), testToken, output)
	if err == nil || !strings.Contains(err.Error(), "provider unavailable") {
		t.Fatalf("expected provider error, got %v", err)
	}
	if _, err := os.Stat(output); !os.IsNotExist(err) {
		t.Fatalf("partial output exists: %v", err)
	}
	entries, err := os.ReadDir(parent)
	if err != nil || len(entries) != 0 {
		t.Fatalf("staging remains: %v %v", entries, err)
	}
}

func TestMissingMessageIDIsAtomic(t *testing.T) {
	parent := t.TempDir()
	output := filepath.Join(parent, "evidence")
	if _, err := Execute(context.Background(), &fakeClient{}, provenProfile(), validRequest(), testToken, output); err == nil || !strings.Contains(err.Error(), "missing message id") {
		t.Fatalf("expected message-id error, got %v", err)
	}
	if _, err := os.Stat(output); !os.IsNotExist(err) {
		t.Fatal("partial output exists")
	}
}

func TestTTLAndTokenBounds(t *testing.T) {
	request := validRequest()
	request.AndroidTTLSeconds = 3601
	if err := Validate(provenProfile(), request, testToken); err == nil {
		t.Fatal("excessive TTL passed")
	}
	if err := Validate(provenProfile(), validRequest(), "short"); err == nil {
		t.Fatal("short token passed")
	}
}

func TestLoadJSONRejectsUnknownFieldsAndTrailingValues(t *testing.T) {
	directory := t.TempDir()
	unknown := filepath.Join(directory, "unknown.json")
	if err := os.WriteFile(unknown, []byte(`{"schema":"elite-firebase-push-request/v1","unknown":true}`), 0o600); err != nil {
		t.Fatal(err)
	}
	var request Request
	if err := LoadJSON(unknown, &request); err == nil {
		t.Fatal("unknown field passed")
	}

	trailing := filepath.Join(directory, "trailing.json")
	if err := os.WriteFile(trailing, []byte(`{"data":{}} {"data":{}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := LoadJSON(trailing, &request); err == nil {
		t.Fatal("trailing JSON value passed")
	}
}
````

### FILE: `firebase_push/cmd/firebase-push/main.go`
```yaml
block_id: "GO-FIREBASE-CLOUD-MESSAGING:cli:v1"
operation: CREATE
provenance: AUTHORED
source: "local CLI using Firebase Admin Go and Application Default Credentials"
license: "LicenseRef-Workspace-Owner"
sha256: "90202057ecf0522beb731adf60cdb0556be507a7f02fc40098558309ad6866cd"
variables: []
secrets_allowed: false
```
````go
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"

	firebasepush "elite.local/firebasepush"
)

func main() {
	profilePath := flag.String("profile", "", "approved provider profile JSON")
	messagePath := flag.String("message", "", "approved message JSON")
	output := flag.String("output", "", "new evidence directory")
	flag.Parse()
	if *profilePath == "" || *messagePath == "" || *output == "" {
		fatal("profile, message and output are required")
	}
	var profile firebasepush.Profile
	var request firebasepush.Request
	if err := firebasepush.LoadJSON(*profilePath, &profile); err != nil {
		fatal(err.Error())
	}
	if err := firebasepush.LoadJSON(*messagePath, &request); err != nil {
		fatal(err.Error())
	}
	token := os.Getenv(profile.DeviceTokenEnvironmentVariable)
	if err := firebasepush.Validate(profile, request, token); err != nil {
		fatal(err.Error())
	}
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: profile.ProjectID})
	if err != nil {
		fatal(err.Error())
	}
	client, err := app.Messaging(ctx)
	if err != nil {
		fatal(err.Error())
	}
	receipt, err := firebasepush.Execute(ctx, client, profile, request, token, *output)
	if err != nil {
		fatal(err.Error())
	}
	fmt.Printf("FIREBASE_PUSH_PASS mode=%s receipt=%s\n", receipt.Mode, *output)
}

func fatal(message string) {
	fmt.Fprintln(os.Stderr, "FIREBASE_PUSH_FAILED:", message)
	os.Exit(1)
}
````

### FILE: `firebase_push/README.md`
```yaml
block_id: "GO-FIREBASE-CLOUD-MESSAGING:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "operator guide and bounded claims"
license: "LicenseRef-Workspace-Owner"
sha256: "3ed5935c2a8c2e79e1b186c0fe7a4f64411ffa55c1941d871f88be6d466eb68d"
variables: []
secrets_allowed: false
```
````markdown
# Google Firebase Cloud Messaging adapter

This adapter imports the official `firebase.google.com/go/v4` module at v4.21.0. It supports one approved test-device message per invocation and defaults to official `SendDryRun`. Real `Send` requires `mode=SEND` and a separate `real_send_approved=true`; the template starts blocked.

Before use, prove the Firebase project, Application Default Credentials, FCM API, device consent, message/data policy, TTL, quota/cost and reconciliation. Put the registration token only in the named environment variable. Do not store tokens in JSON, Markdown, logs or receipts. ADC is supplied by the selected Google environment/secret mechanism, not embedded in this pack.

```powershell
go mod download
go mod verify
go test ./...
go run ./cmd/firebase-push -profile .\provider-profile.json -message .\message.json -output .\evidence\push-001
```

Tests exercise dry-run selection, explicit real-send approval, official `messaging.Message`, allowlisted data, reserved keys, token/TTL bounds, provider failure, missing message ID, atomic output and receipt redaction. They do not prove credentials, device delivery, foreground/background behavior, APNs/WebPush configuration, consent/opt-out, quotas, billing or reconciliation.
````

## 6. Configuration surface

| Variable/campo | Tipo/default seguro | Validación | Secreto | Mutabilidad/efecto |
|---|---|---|---|---|
| `decision` | enum/`BLOCKED` | sólo `PROVEN` habilita | no | nueva aprobación |
| pruebas de proyecto/ADC/API/consentimiento/políticas/cuota/reconciliación | bool/`false` | todas `true` | no | nueva evidencia |
| `project_id` | string/vacío | patrón Firebase | identificador sensible | cambio de proyecto |
| `mode` | enum/`DRY_RUN` | `DRY_RUN` o `SEND` | no | por campaña controlada |
| `real_send_approved` | bool/`false` | requerido para `SEND` | no | aprobación separada |
| `device_token_environment_variable` | string/vacío | nombre de entorno válido | no; su valor sí | runtime |
| token resuelto | string/sin default | longitud acotada; nunca archivo | sí | runtime/rotación |
| `approved_data_keys` | lista/vacía | patrón y claves reservadas prohibidas | no | política |
| `maximum_ttl_seconds` | entero/`0` | 1..2.419.200 | no | política |
| `automatic_business_write` | bool/`false` | debe permanecer `false` | no | inmutable en este pack |

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| `firebase.google.com/go/v4` | `v4.21.0`, commit `eebb06f2a643fbb59b1cb262874a943584475128` | Admin SDK y FCM | Apache-2.0 | runtime | https://github.com/firebase/firebase-admin-go |
| Go | `1.26.7` | format, módulos, tests, vet, build | BSD-3-Clause | build/runtime | https://go.dev/dl/ |
| grafo transitivo | exacto en `go.sum` | dependencias del SDK oficial | según módulo; inventariar en proyecto | runtime/build | Go module proxy + repos oficiales |

## 8. Apply order

1. Componer y aprobar `OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x`; adquirir/verificar Firebase sólo después del expediente.
2. Materializar los nueve archivos en un workspace vacío o comprobar colisiones idénticas en uno existente.
3. Ejecutar `go mod download`, `go mod verify`, `go test ./...`, `go vet ./...` y build con Go 1.26.7.
4. Copiar las plantillas fuera del pack, completar pruebas/aprobaciones y resolver ADC/token desde el secret manager del target.
5. Empezar en `DRY_RUN`; demostrar recibo, redacción y reconciliación. Habilitar `SEND` sólo con aprobación separada y dispositivo consentido.
6. Rollback: detener wiring, revocar credenciales/tokens y conservar evidencia; no modificar outputs canónicos.

## 9. Verification

- Materialización: 9/9 paths, hashes y manifest exactos; reconstrucción byte-idéntica.
- Dependencias: `go mod verify` debe emitir `all modules verified`; el módulo oficial debe resolver exactamente `v4.21.0`.
- Código: `go test ./...`, `go vet ./...` y `go build ./cmd/firebase-push` deben terminar cero.
- Negativos: perfil bloqueado, data no aprobada/reservada, TTL/token inválidos, JSON desconocido/trailing, real send no aprobado, fallo del proveedor y message ID ausente deben fallar cerrados.
- Seguridad/privacidad: recibo sin proyecto, token ni message ID en claro; ningún secreto en pack; output atómico.
- Integración real queda condicionada a proyecto/ADC/API/device consent y un dry-run autorizado; entrega, carga, cuota, billing, APNs/WebPush y reconciliación productiva siguen gates del proyecto.
- Licencias/SBOM: conservar lock Apache-2.0, `go.mod`, `go.sum` y generar inventario transitivo en la adopción.

## 10. Reconstruction evidence

Evidencia canónica: `reconstruction_evidence/FIREBASE_CLOUD_MESSAGING_ADAPTER_2026-08-26_V1.md`. Entorno limpio temporal, Go 1.26.7 oficial verificado por hash, Firebase Admin Go 4.21.0, nueve archivos reconstruidos, siete tests, module verify, vet y build. No se usaron credenciales ni se envió una notificación real; por ello la admisión permanece `CONDITIONED`.
