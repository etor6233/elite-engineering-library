# Go HTTP Metrics Reference

## 1. Metadata

```yaml
pack_id: "GO-HTTP-METRICS-REFERENCE"
pack_version: "0.1.15"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa instrumentacion HTTP con atributos cerrados sobre middleware/SDK/exporter oficiales fijados y un ensayo finito de API real, PostgreSQL, regla Prometheus intacta, fallo y recuperacion; identidad sintetica sólo en el main de referencia. No despliega ni admite identidad, SLI, TLS o alertas de produccion."
stacks: ["Go 1.26.8", "CPython 3.14", "trusted local Windows x64", "PostgreSQL 18.6", "Prometheus 3.14.0"]
compatible_with: ["GO-ENTERPRISE-BACKEND 0.4.4", "SECURE-OPS-DELIVERY-CORE 1.1.4"]
incompatible_with: ["production deployment of the synthetic principal", "unadmitted runtime artifacts", "arbitrary consumer schema"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://github.com/open-telemetry/opentelemetry-go-contrib/tree/c8a87a60ba1b3374fd16df11fc3eeae6c41abbc9/instrumentation/net/http/otelhttp", "https://github.com/open-telemetry/opentelemetry-go/tree/93a693edeed0e07ce5ebd1dfe67af42d1e2055d8/exporters/prometheus", "https://opentelemetry.io/docs/specs/semconv/http/http-metrics/"]
verified_at: "2026-09-11"
```

## 2. Applicability

Opt-in synthetic qualification on a trusted local account. Resolve the exact external runtime lock and build/SCA/license conditions before use. The reference command is never a production identity adapter. See the materialized README and V375 evidence.

## 3. Architecture contract

- Reuse the real HTTP handler, order service, repository and migrations; no fake HTTP statuses or injected Prometheus samples.
- Keep one durable order/outbox/idempotency result across fault and replay.
- Export only method/status and duration histogram; private registry, no traces or exemplars.
- Preserve actual five-minute rate window, one-percent threshold, ten-minute hold and thirty-second evaluation.
- Bound the reference run, retain failures and stop owned processes. Production exposure, SLI selection and response remain target gates.

## 4. Exact file manifest

```text
CREATE reference_http_metrics/go.mod
CREATE reference_http_metrics/go.sum
CREATE reference_http_metrics/main.go
CREATE reference_http_metrics/metrics.go
CREATE reference_http_metrics/metrics_test.go
CREATE reference_http_metrics/README.md
CREATE reference_http_metrics/run_reference.py
CREATE reference_http_metrics/source-lock.json
```

## 5. Materialization blocks


### FILE: `reference_http_metrics/go.mod`

```yaml
block_id: "GO-HTTP-METRICS-REFERENCE:go.mod:v1"
operation: CREATE
provenance: AUTHORED
source: "local composition over pinned official dependency APIs; no upstream runtime source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "1817c036fec9d5bb51b6785dc15ad25bae0bedc4984fa9ecd198f6c4abe857f1"
variables: []
secrets_allowed: false
```

````text
module elite.local/enterprise/reference_http_metrics

go 1.26.0

require (
	elite.local/enterprise v0.0.0
	github.com/jackc/pgx/v5 v5.10.0
	github.com/prometheus/client_golang v1.24.1
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.70.0
	go.opentelemetry.io/otel v1.46.0
	go.opentelemetry.io/otel/exporters/prometheus v0.67.0
	go.opentelemetry.io/otel/sdk v1.46.0
	go.opentelemetry.io/otel/sdk/metric v1.46.0
	go.opentelemetry.io/otel/trace v1.46.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/coreos/go-oidc/v3 v3.20.0 // indirect
	github.com/felixge/httpsnoop v1.1.0 // indirect
	github.com/go-jose/go-jose/v4 v4.1.4 // indirect
	github.com/go-logr/logr v1.4.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.70.1 // indirect
	github.com/prometheus/otlptranslator v1.0.0 // indirect
	github.com/prometheus/procfs v0.21.1 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel/metric v1.46.0 // indirect
	golang.org/x/oauth2 v0.36.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace elite.local/enterprise => ..
````

### FILE: `reference_http_metrics/go.sum`

```yaml
block_id: "GO-HTTP-METRICS-REFERENCE:go.sum:v1"
operation: CREATE
provenance: AUTHORED
source: "local composition over pinned official dependency APIs; no upstream runtime source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "5779dd559c4a188490a8932427b5f84a152b162781b172533d2e666dbeff8ed2"
variables: []
secrets_allowed: false
```

````text
github.com/beorn7/perks v1.0.1 h1:VlbKKnNfV8bJzeqoa4cOKqO6bYr3WgKZxO8Z16+hsOM=
github.com/beorn7/perks v1.0.1/go.mod h1:G2ZrVWU2WbWT9wwq4/hrbKbnv/1ERSJQ0ibhJ6rlkpw=
github.com/cespare/xxhash/v2 v2.3.0 h1:UL815xU9SqsFlibzuggzjXhog7bL6oX9BbNZnL2UFvs=
github.com/cespare/xxhash/v2 v2.3.0/go.mod h1:VGX0DQ3Q6kWi7AoAeZDth3/j3BFtOZR5XLFGgcrjCOs=
github.com/coreos/go-oidc/v3 v3.20.0 h1:EtE0WIBHk03N+DqGkY4+UONzzZHk7amKt6IyNd7OsZE=
github.com/coreos/go-oidc/v3 v3.20.0/go.mod h1:DYCf24+ncYi+XkIH97GY1+dqoRlbaSI26KVTCI9SrY4=
github.com/davecgh/go-spew v1.1.0/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/felixge/httpsnoop v1.1.0 h1:3YtUj32ZZkqZtt3sZZsClsymw/QDuVfpNhoA31zeORc=
github.com/felixge/httpsnoop v1.1.0/go.mod h1:Zqxgdd+1Rkcz8euOqdr7lqgCRJztwr5hp9vDSi5UZCE=
github.com/go-jose/go-jose/v4 v4.1.4 h1:moDMcTHmvE6Groj34emNPLs/qtYXRVcd6S7NHbHz3kA=
github.com/go-jose/go-jose/v4 v4.1.4/go.mod h1:x4oUasVrzR7071A4TnHLGSPpNOm2a21K9Kf04k1rs08=
github.com/go-logr/logr v1.2.2/go.mod h1:jdQByPbusPIv2/zmleS9BjJVeZ6kBagPoEUsqbVz/1A=
github.com/go-logr/logr v1.4.4 h1:tG4xh9yMsRCAiodLVTxyrkzSZ9+o0L1Kg/+cPVcbP/8=
github.com/go-logr/logr v1.4.4/go.mod h1:9T104GzyrTigFIr8wt5mBrctHMim0Nb2HLGrmQ40KvY=
github.com/go-logr/stdr v1.2.2 h1:hSWxHoqTgW2S2qGc0LTAI563KZ5YKYRhT3MFKZMbjag=
github.com/go-logr/stdr v1.2.2/go.mod h1:mMo/vtBO5dYbehREoey6XUKy/eSumjCCveDpRre4VKE=
github.com/google/go-cmp v0.7.0 h1:wk8382ETsv4JYUZwIsn6YpYiWiBsYLSJiTsyBybVuN8=
github.com/google/go-cmp v0.7.0/go.mod h1:pXiqmnSA92OHEEa9HXL2W4E7lf9JzCmGVUdgjX3N/iU=
github.com/google/uuid v1.6.0 h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=
github.com/google/uuid v1.6.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
github.com/jackc/pgpassfile v1.0.0 h1:/6Hmqy13Ss2zCq62VdNG8tM1wchn8zjSGOBJ6icpsIM=
github.com/jackc/pgpassfile v1.0.0/go.mod h1:CEx0iS5ambNFdcRtxPj5JhEz+xB6uRky5eyVu/W2HEg=
github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 h1:iCEnooe7UlwOQYpKFhBabPMi4aNAfoODPEFNiAnClxo=
github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761/go.mod h1:5TJZWKEWniPve33vlWYSoGYefn3gLQRzjfDlhSJ9ZKM=
github.com/jackc/pgx/v5 v5.10.0 h1:VhSvgU2jSli8o3AqIEOTJr7rZwAEUVo4E4XhR94Zfr0=
github.com/jackc/pgx/v5 v5.10.0/go.mod h1:mal1tBGAFfLHvZzaYh77YS/eC6IX9OWbRV1QIIM0Jn4=
github.com/jackc/puddle/v2 v2.2.2 h1:PR8nw+E/1w0GLuRFSmiioY6UooMp6KJv0/61nB7icHo=
github.com/jackc/puddle/v2 v2.2.2/go.mod h1:vriiEXHvEE654aYKXXjOvZM39qJ0q+azkZFrfEOc3H4=
github.com/klauspost/compress v1.19.1 h1:VsB4HPswih7mmZ8WleSFQ75c/Ui1M4trX5oAsJnhSlk=
github.com/klauspost/compress v1.19.1/go.mod h1:cwPg85FWrGar70rWktvGQj8/hthj3wpl0PGDogxkrSQ=
github.com/kylelemons/godebug v1.1.0 h1:RPNrshWIDI6G2gRW9EHilWtl7Z6Sb1BR0xunSBf0SNc=
github.com/kylelemons/godebug v1.1.0/go.mod h1:9/0rRGxNHcop5bhtWyNeEfOS8JIWk580+fNqagV/RAw=
github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 h1:C3w9PqII01/Oq1c1nUAm88MOHcQC9l5mIlSMApZMrHA=
github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822/go.mod h1:+n7T8mK8HuQTcFwEeznm/DIxMOiR9yIdICNftLE1DvQ=
github.com/pmezard/go-difflib v1.0.0/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
github.com/prometheus/client_golang v1.24.1 h1:JnJkREXzWxUdCuPFpIWZiPispT9xVV59uiuyR2bPlnU=
github.com/prometheus/client_golang v1.24.1/go.mod h1:F+oSRECHg4sse5ucfYpYDeIv/hu68Zo0uoHKetWnzcE=
github.com/prometheus/client_model v0.6.2 h1:oBsgwpGs7iVziMvrGhE53c/GrLUsZdHnqNwqPLxwZyk=
github.com/prometheus/client_model v0.6.2/go.mod h1:y3m2F6Gdpfy6Ut/GBsUqTWZqCUvMVzSfMLjcu6wAwpE=
github.com/prometheus/common v0.70.1 h1:1HvjP4D5oL3t8RsPlwxA9onvvStjtIHYE5XuuwOi/PY=
github.com/prometheus/common v0.70.1/go.mod h1:VdFUQDMZK3VLkurFUVhia6uys/0suUp86TJz5qbJRhc=
github.com/prometheus/otlptranslator v1.0.0 h1:s0LJW/iN9dkIH+EnhiD3BlkkP5QVIUVEoIwkU+A6qos=
github.com/prometheus/otlptranslator v1.0.0/go.mod h1:vRYWnXvI6aWGpsdY/mOT/cbeVRBlPWtBNDb7kGR3uKM=
github.com/prometheus/procfs v0.21.1 h1:GljZCt+zSTS+NZq88cyQ1LjZ+RCHp3uVuabBWA5+OJI=
github.com/prometheus/procfs v0.21.1/go.mod h1:aB55Cww9pdSJVHk0hUf0inxWyyjPogFIjmHKYgMKmtY=
github.com/stretchr/objx v0.1.0/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
github.com/stretchr/testify v1.3.0/go.mod h1:M5WIy9Dh21IEIfnGCwXGc5bZfKNJtfHm1UVUgZn+9EI=
github.com/stretchr/testify v1.7.0/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
github.com/stretchr/testify v1.12.1 h1:EuwCh5fleGS7H32xRwO3wRGT7DxrDhLAT6FF8MpWDWE=
github.com/stretchr/testify v1.12.1/go.mod h1:MDEgiDPPsNp5cuIrHPPCyornHKgEVbtFUmoNlxoYthg=
go.opentelemetry.io/auto/sdk v1.2.1 h1:jXsnJ4Lmnqd11kwkBV2LgLoFMZKizbCi5fNZ/ipaZ64=
go.opentelemetry.io/auto/sdk v1.2.1/go.mod h1:KRTj+aOaElaLi+wW1kO/DZRXwkF4C5xPbEe3ZiIhN7Y=
go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.70.0 h1:LMuyCAyfalSjDyjdC65nK6N0zoTT63+E/u95X0JovZI=
go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.70.0/go.mod h1:085m8qbm4hgc8rZWGDEa4vmyyo2c3nPxUslYUKUIU04=
go.opentelemetry.io/otel v1.46.0 h1:FHt5/CDyVxi/8IM1CH7VE/rRgq3kLHa2mSTVMO8AWyc=
go.opentelemetry.io/otel v1.46.0/go.mod h1:Gj3SEScelsNC45tp4nSxRYlS+f5iez7W8XPMCt905kE=
go.opentelemetry.io/otel/exporters/prometheus v0.67.0 h1:7IefDa35e6V3NoiqIeLDMDxMFyZDk5qcoC0Ax4cC16E=
go.opentelemetry.io/otel/exporters/prometheus v0.67.0/go.mod h1:nsPI1awTg5Vmg1YrommL2mVarVGlqc4yXOoKAkPRD0c=
go.opentelemetry.io/otel/metric v1.46.0 h1:yBnkXvgV7AXFILZc5K6IZe/CBFF3OS7BJ8ov6/lj0K8=
go.opentelemetry.io/otel/metric v1.46.0/go.mod h1:iPmdWqifKUdzziPkvvzIJXITl56fQx2mGM/DHLB3/2o=
go.opentelemetry.io/otel/metric/x v0.68.0 h1:TA/cBT23D3MnxYPwHL7YFOdYGdx0A0v+s7Mzotpd1dU=
go.opentelemetry.io/otel/metric/x v0.68.0/go.mod h1:agudOmvWhwUTjgibWDzxD2PoWYnpw5Ht5jISYOD2Hd4=
go.opentelemetry.io/otel/sdk v1.46.0 h1:h5CNQQjEbuQXY/JfZtgt3i7HVFV3aHPO2OAwO2eTYPI=
go.opentelemetry.io/otel/sdk v1.46.0/go.mod h1:GAERFXFt5SYCEB+YiKUbMBeza6UaDH7GmGOZEfh2gSM=
go.opentelemetry.io/otel/sdk/metric v1.46.0 h1:0piZ26EG4RBfebb2jhDH6ERCYHoVWduc3kLgPCwSnSE=
go.opentelemetry.io/otel/sdk/metric v1.46.0/go.mod h1:I1PbKrdVc8Qu8HYVDNtqVIwLwjNrhsV/uFuxfwg8mO4=
go.opentelemetry.io/otel/trace v1.46.0 h1:OULy7ccdJnZtJ0UDYFOIGaCmiWzJ8Vi2G/Rsu60qs1c=
go.opentelemetry.io/otel/trace v1.46.0/go.mod h1:J7GAXweO77XSFkB/rmAqk9D6ihszhFjLU+d9WuUxDLI=
go.uber.org/goleak v1.3.0 h1:2K3zAYmnTNqV73imy9J1T3WC+gmCePx2hEGkimedGto=
go.uber.org/goleak v1.3.0/go.mod h1:CoHD4mav9JJNrW/WLlf7HGZPjdw8EucARQHekz1X6bE=
go.yaml.in/yaml/v2 v2.4.4 h1:tuyd0P+2Ont/d6e2rl3be67goVK4R6deVxCUX5vyPaQ=
go.yaml.in/yaml/v2 v2.4.4/go.mod h1:gMZqIpDtDqOfM0uNfy0SkpRhvUryYH0Z6wdMYcacYXQ=
go.yaml.in/yaml/v3 v3.0.5 h1:N6y/pJk8buWs9NY5ERU2HSMfm+IuD/OtfdAnq6kESPw=
go.yaml.in/yaml/v3 v3.0.5/go.mod h1:HVTZu1O7/Vkt2N+BFy8Zza+lnLsABggaTM2ZpNIGuKg=
golang.org/x/oauth2 v0.36.0 h1:peZ/1z27fi9hUOFCAZaHyrpWG5lwe0RJEEEeH0ThlIs=
golang.org/x/oauth2 v0.36.0/go.mod h1:YDBUJMTkDnJS+A4BP4eZBjCqtokkg1hODuPjwiGPO7Q=
golang.org/x/sync v0.22.0 h1:SZjpbeLmrCk4xhRSZFNZW5gFUeCeFgjekvI/+gfScek=
golang.org/x/sync v0.22.0/go.mod h1:9xrNwdLfx4jkKbNva9FpL6vEN7evnE43NNNJQ2LF3+0=
golang.org/x/sys v0.47.0 h1:o7XGOvZQCADBQQ4Y7VNq2dRWQR7JmOUW8Kxx4ZsNgWs=
golang.org/x/sys v0.47.0/go.mod h1:4GL1E5IUh+htKOUEOaiffhrAeqysfVGipDYzABqnCmw=
golang.org/x/text v0.40.0 h1:Ub2Z6/xjgF1WrYQz2nuITOEegKFtiIy+rieRJ5lHZKs=
golang.org/x/text v0.40.0/go.mod h1:hpnzDAfGV753zIKo+wk3u1bVKCGPbrnF7+7LBF/UHVY=
google.golang.org/protobuf v1.36.11 h1:fV6ZwhNocDyBLK0dj+fg8ektcVegBBuEolpbTQyBNVE=
google.golang.org/protobuf v1.36.11/go.mod h1:HTf+CrKn2C3g5S8VImy6tdcUvCska2kB7j23XfzDpco=
gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
````

### FILE: `reference_http_metrics/main.go`

```yaml
block_id: "GO-HTTP-METRICS-REFERENCE:main.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local composition over pinned official dependency APIs; no upstream runtime source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "fc7b6fab0107e502d48cbd0515233bdb21d7140f98374f799394603400f086eb"
variables: []
secrets_allowed: false
```

````go
// Finite synthetic qualification host; never a production identity provider.
package main

import (
 "context"
 "crypto/rand"
 "encoding/hex"
 "encoding/json"
 "errors"
 "flag"
 "fmt"
 "net"
 "net/http"
 "os"
 "time"

 "elite.local/enterprise/internal/order"
 "elite.local/enterprise/internal/platform/httpapi"
 "elite.local/enterprise/internal/platform/identity"
 "elite.local/enterprise/internal/platform/postgres"
 "github.com/jackc/pgx/v5/pgxpool"
)

type ids struct{}
func (ids) New() string { var b [16]byte; if _,err:=rand.Read(b[:]);err!=nil {panic(err)};return hex.EncodeToString(b[:]) }
type fixtureIdentity struct{}
func (fixtureIdentity) Verify(_ context.Context, raw string)(identity.Principal,error){
 if raw!="synthetic-local-reference-only" {return identity.Principal{},identity.ErrUnauthenticated}
 return identity.Principal{Subject:"synthetic-customer",TenantID:"018f4d4a-7b36-7a21-8d10-2f4c54c28b01",Permissions:map[string]struct{}{"order:create":{}},Organizations:map[string]struct{}{"integration-org":{}}},nil
}
func run() error {
 ttl:=flag.Duration("ttl",18*time.Minute,"finite reference lifetime, at most twenty minutes")
 flag.Parse()
 if *ttl<=0 || *ttl>20*time.Minute {return errors.New("invalid reference lifetime")}
 ctx,cancel:=context.WithTimeout(context.Background(),*ttl);defer cancel()
 pool,err:=pgxpool.New(ctx,os.Getenv("REFERENCE_DATABASE_URL"));if err!=nil{return errors.New("reference database configuration")}
 defer pool.Close()
 if err=pool.Ping(ctx);err!=nil{return errors.New("reference database unavailable")}
 app,metrics,shutdown,err:=instrument(httpapi.New(order.NewService(postgres.NewOrders(pool),ids{}),fixtureIdentity{}));if err!=nil{return errors.New("reference metrics configuration")}
 defer func(){c,done:=context.WithTimeout(context.Background(),5*time.Second);defer done();_ = shutdown(c)}()
 listener,err:=net.Listen("tcp","127.0.0.1:0");if err!=nil{return err}
 metricListener,err:=net.Listen("tcp","127.0.0.1:0");if err!=nil{listener.Close();return err}
 server:=&http.Server{Handler:app,ReadHeaderTimeout:5*time.Second,ReadTimeout:15*time.Second,WriteTimeout:15*time.Second,IdleTimeout:30*time.Second,MaxHeaderBytes:16384}
 metricServer:=&http.Server{Handler:metrics,ReadHeaderTimeout:3*time.Second,WriteTimeout:5*time.Second,IdleTimeout:10*time.Second,MaxHeaderBytes:4096}
 failures:=make(chan error,2)
 go func(){failures<-server.Serve(listener)}()
 go func(){failures<-metricServer.Serve(metricListener)}()
 if err=json.NewEncoder(os.Stdout).Encode(map[string]string{"api":listener.Addr().String(),"metrics":metricListener.Addr().String()});err!=nil{return err}
 select {case <-ctx.Done():case err=<-failures:if !errors.Is(err,http.ErrServerClosed){cancel()}}
 c,done:=context.WithTimeout(context.Background(),5*time.Second);defer done()
 a:=server.Shutdown(c);b:=metricServer.Shutdown(c)
 if a!=nil||b!=nil{return errors.New("reference shutdown incomplete")}
 return nil
}
func main(){if err:=run();err!=nil{fmt.Fprintln(os.Stderr,"HTTP_METRICS_REFERENCE_FAILED");os.Exit(1)}}

````

### FILE: `reference_http_metrics/metrics.go`

```yaml
block_id: "GO-HTTP-METRICS-REFERENCE:metrics.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local composition over pinned official dependency APIs; no upstream runtime source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "0c259b32e3cf15b371178aa8d92c658774913607f6d6435abf7b9d2d96377883"
variables: []
secrets_allowed: false
```

````go
// AUTHORED composition over pinned official OpenTelemetry and Prometheus APIs.
// Trusted local qualification only. It does not enable a production endpoint.
package main

import (
 "context"
 "net/http"

 "github.com/prometheus/client_golang/prometheus"
 "github.com/prometheus/client_golang/prometheus/promhttp"
 "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
 "go.opentelemetry.io/otel/attribute"
 exporter "go.opentelemetry.io/otel/exporters/prometheus"
 "go.opentelemetry.io/otel/propagation"
 sdk "go.opentelemetry.io/otel/sdk/metric"
 "go.opentelemetry.io/otel/sdk/metric/exemplar"
 "go.opentelemetry.io/otel/sdk/resource"
 "go.opentelemetry.io/otel/trace/noop"
)

func instrument(next http.Handler) (http.Handler, http.Handler, func(context.Context) error, error) {
 registry := prometheus.NewRegistry()
 reader, err := exporter.New(exporter.WithRegisterer(registry), exporter.WithoutScopeInfo(), exporter.WithoutTargetInfo())
 if err != nil { return nil, nil, nil, err }
 view := func(i sdk.Instrument) (sdk.Stream, bool) {
  if i.Name != "http.server.request.duration" {
   return sdk.Stream{Aggregation:sdk.AggregationDrop{}}, true
  }
  return sdk.Stream{
   Name:i.Name, Unit:"s",
   Aggregation:sdk.AggregationExplicitBucketHistogram{Boundaries:[]float64{0.001,0.005,0.01,0.05,0.1,0.5,1,5,10}},
   AttributeFilter:attribute.NewAllowKeysFilter(attribute.Key("http.request.method"),attribute.Key("http.response.status_code")),
  },true
 }
 provider:=sdk.NewMeterProvider(sdk.WithReader(reader),sdk.WithResource(resource.Empty()),sdk.WithView(view),sdk.WithCardinalityLimit(128),sdk.WithExemplarFilter(exemplar.AlwaysOffFilter))
 wrapped:=otelhttp.NewHandler(next,"reference-http",otelhttp.WithMeterProvider(provider),otelhttp.WithTracerProvider(noop.NewTracerProvider()),otelhttp.WithPropagators(propagation.NewCompositeTextMapPropagator()))
 return wrapped,promhttp.HandlerFor(registry,promhttp.HandlerOpts{ErrorHandling:promhttp.HTTPErrorOnError,Timeout:5e9}),provider.Shutdown,nil
}

````

### FILE: `reference_http_metrics/metrics_test.go`

```yaml
block_id: "GO-HTTP-METRICS-REFERENCE:metrics_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local composition over pinned official dependency APIs; no upstream runtime source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "f21368b97b4ccbc8ab653b573dd568a3f98ae9be1846e2952acabfca210bcebb"
variables: []
secrets_allowed: false
```

````go
package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRealHTTPStatusAndClosedLabels(t *testing.T) {
	handler, metrics, close, err := instrument(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/fail" {
			w.WriteHeader(503)
			return
		}
		w.WriteHeader(204)
	}))
	if err != nil {
		t.Fatal(err)
	}
	defer close(context.Background())
	server := httptest.NewServer(handler)
	defer server.Close()
	for _, path := range []string{"/ok?secret=private@example.invalid", "/fail"} {
		req, _ := http.NewRequest("GET", server.URL+path, strings.NewReader("PRIVATE_BODY"))
		req.Host = "PRIVATE_HOST.invalid"
		req.Header.Set("Authorization", "Bearer PRIVATE_TOKEN")
		req.Header.Set("Cookie", "PRIVATE_COOKIE")
		req.Header.Set("X-Private", "PRIVATE_HEADER")
		res, err := server.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		io.Copy(io.Discard, res.Body)
		res.Body.Close()
		if path == "/fail" && res.StatusCode != 503 {
			t.Fatal(res.StatusCode)
		}
	}
	w := httptest.NewRecorder()
	metrics.ServeHTTP(w, httptest.NewRequest("GET", "http://metrics/", nil))
	body := w.Body.String()
	for _, marker := range []string{"PRIVATE_", "private@example", "/fail", "/ok", "server_address", "url_", "http_route", "target_info", "otel_scope", "http_server_request_body"} {
		if strings.Contains(body, marker) {
			t.Fatalf("unexpected metric content %s", marker)
		}
	}
	for _, line := range []string{
		"http_server_request_duration_seconds_count{http_request_method=\"GET\",http_response_status_code=\"204\"} 1",
		"http_server_request_duration_seconds_count{http_request_method=\"GET\",http_response_status_code=\"503\"} 1",
	} {
		if !strings.Contains(body, line) {
			t.Fatalf("missing %s in %s", line, body)
		}
	}
}

func TestIndependentInstancesDoNotMixRequests(t *testing.T) {
	a, am, ac, err := instrument(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(201) }))
	if err != nil {
		t.Fatal(err)
	}
	defer ac(context.Background())
	b, bm, bc, err := instrument(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(503) }))
	if err != nil {
		t.Fatal(err)
	}
	defer bc(context.Background())
	a.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "http://reference/order", nil))
	b.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "http://reference/failure", nil))
	for _, c := range []struct {
		h            http.Handler
		want, reject string
	}{
		{am, "http_response_status_code=\"201\"", "http_response_status_code=\"503\""},
		{bm, "http_response_status_code=\"503\"", "http_response_status_code=\"201\""},
	} {
		w := httptest.NewRecorder()
		c.h.ServeHTTP(w, httptest.NewRequest("GET", "http://metrics/", nil))
		if !strings.Contains(w.Body.String(), c.want) || strings.Contains(w.Body.String(), c.reject) {
			t.Fatal("instance metrics mixed")
		}
		for _, forbidden := range []string{"go_gc", "go_goroutines", "process_", "target_info"} {
			if strings.Contains(w.Body.String(), forbidden) {
				t.Fatal("default collector leaked", forbidden)
			}
		}
	}
}
````

### FILE: `reference_http_metrics/README.md`

```yaml
block_id: "GO-HTTP-METRICS-REFERENCE:README.md:v1"
operation: CREATE
provenance: AUTHORED
source: "local composition over pinned official dependency APIs; no upstream runtime source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "473717b994f9e09f98d5cbf7d788c760319ed824d6a1bdd14025ee19314caf27"
variables: []
secrets_allowed: false
```

````markdown
# HTTP request metrics and canonical alert reference

Trusted local Windows x64 qualification over the actual enterprise HTTP handler,
order service and PostgreSQL repository. This command uses a fixed synthetic
principal exclusively in an isolated reference cluster; never deploy this main
or treat it as an OIDC adapter. The production main and authorization code stay
unchanged. API and metric listeners bind only to loopback.

The AUTHORED metrics factory composes pinned official otelhttp, OpenTelemetry SDK
and Prometheus exporter. It retains only request method/status and the duration
histogram; drops other instruments/attributes, disables exemplars and tracing,
uses a private registry/resource and bounds aggregation cardinality to128.
Separate instances do not share series. It does not cover outbound requests,
unrecovered panics, arbitrary new handlers, a production SLI policy or all PII.

Materialize the full reference plan into an absent destination. This nested Go
module resolves the materialized enterprise module from its parent. Use the
exact Go1.26.8 toolchain and the admitted/cache-verified dependency graph; no
floating updates. From reference_http_metrics run:
  go test -mod=readonly -count=1 ./...
  go vet -mod=readonly ./...
  go build -mod=readonly -trimpath -o http-metrics-reference.exe .

Build flags include GOTOOLCHAIN=local, GOWORK=off and CGO_ENABLED=0. Network
acquisition follows the consumer's source/dependency admission; the V375 run used
existing cached source artifacts. Record build info, source hashes and SCA for
the resulting binary. source-lock.json records the qualified reference inputs;
its binary digest is evidence of that build, not an approval of any new binary.

Run run_reference.py with the existing admitted V372 runtime lock and its exact
SHA, --consumer pointing at the composed root, --host and --host-sha256 identifying
the rebuilt reference command, --rules and --rules-sha256 identifying the actual
ops/prometheus/platform.rules.yml, and --work-root an existing private local
directory. The V372 telemetry-reference binary generates ephemeral test TLS
identities; the lock also binds PostgreSQL18.6, Prometheus3.14.0 and promtool.
Those runtime binaries remain external and are not bundled or installed here.

The harness creates a new private directory and a new owned PostgreSQL cluster,
checks data-directory identity, applies the exact lock's complete migrations and
seeds two synthetic tenant/organization rows. It proves order creation/replay,
401/400 handling and real PostgreSQL-fault500 responses. One failed request every
240seconds must make the unchanged1% /10minute /30second-evaluation rule fire;
the previous denominator floor must remain below its threshold on the same
real series. Repair and successful replays must clear the alert, with exactly
one order, one idempotency record and one order.created outbox row afterward.
The reference takes roughly15minutes because it preserves the actual hold.
It never compresses the alert window, changes system time or writes synthetic
samples into Prometheus. The native host expires after20minutes; the harness
uses a finite1140second budget and retains failed runs.

API/metrics traffic is loopback HTTP in this trusted-account test; Prometheus
query access uses ephemeral mutual TLS1.3. Bind real telemetry exposure,
transport security, identity, SLI selection (including health routes), alert
routing, response policy, quotas, storage and recovery to the deployment target.
No real customer input, private history, provider, payment, mail or fiscal call
is authorized by this reference. Current scope does not close TEST02/03/07.

Evidence: reconstruction_evidence/HTTP_SLI_ALERT_INTEGRATION_V375.md.
````

### FILE: `reference_http_metrics/run_reference.py`

```yaml
block_id: "GO-HTTP-METRICS-REFERENCE:run_reference.py:v1"
operation: CREATE
provenance: AUTHORED
source: "local composition over pinned official dependency APIs; no upstream runtime source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "e23057fc1a71ced1dddfbefc51886ce66508d483effe64e87a933527956f31ab"
variables: []
secrets_allowed: false
```

````python
"""Finite, synthetic HTTP -> official metrics -> canonical alert qualification."""
from pathlib import Path
import argparse,subprocess,os,json,hashlib,socket,urllib.request,urllib.error,urllib.parse,ssl,time,uuid,csv,re
p=argparse.ArgumentParser()
for n in ['runtime-lock','runtime-lock-sha256','consumer','host','host-sha256','rules','rules-sha256','work-root']:p.add_argument('--'+n,required=True)
a=p.parse_args();sha=lambda p:hashlib.sha256(p.read_bytes()).hexdigest()
lock_path=Path(a.runtime_lock);assert sha(lock_path)==a.runtime_lock_sha256
lock=json.loads(lock_path.read_text());bins={}
for n in ['postgres','initdb','pg_ctl','psql','prometheus','promtool','telemetry-reference']:
 pin=lock['binaries'][n];q=Path(pin['path']);assert q.is_file() and q.stat().st_size==pin['bytes'] and sha(q)==pin['sha256'];bins[n]=q
host=Path(a.host);rules=Path(a.rules);assert sha(host)==a.host_sha256 and sha(rules)==a.rules_sha256
consumer=Path(a.consumer);migrations=sorted((consumer/'db/migrations').glob('*.up.sql'))
assert {x.name:sha(x) for x in migrations}==lock['migrations']
root=Path(a.work_root).resolve(strict=True);W=root/('http-'+uuid.uuid4().hex);W.mkdir()
system=Path(os.environ['SYSTEMROOT'])/'System32'
sid=list(csv.reader(subprocess.check_output([str(system/'whoami.exe'),'/user','/fo','csv','/nh'],text=True).splitlines()))[0][1]
assert re.fullmatch(r'S-[0-9-]+',sid)
subprocess.run([str(system/'icacls.exe'),str(W),'/inheritance:r','/grant:r','*'+sid+':(OI)(CI)F','*S-1-5-18:(OI)(CI)F'],check=True,stdout=subprocess.DEVNULL,creationflags=subprocess.CREATE_NO_WINDOW)
print('HTTP_REFERENCE_STAGE '+str(W),flush=True)
env={k:os.environ[k] for k in ['SYSTEMROOT','WINDIR']};env.update(TEMP=str(W),TMP=str(W),PGCONNECT_TIMEOUT='3',COMSPEC=str(system/'cmd.exe'))
procs={};handles=[];events=[];result={'state':'RUNNING','rules_sha256':a.rules_sha256,'host_sha256':a.host_sha256,'requests':{},'synthetic_only':True}
start=time.monotonic();deadline=start+1140
def log(kind,**values):
 events.append({'seconds':round(time.monotonic()-start,3),'kind':kind,**values})
 (W/'progress.json').write_text(json.dumps(events,indent=2)+'\n',encoding='utf-8',newline='\n')
 print(json.dumps(events[-1]),flush=True)
def command(name,args,seconds=30):
 q=subprocess.run([str(v) for v in args],env=env,stdin=subprocess.DEVNULL,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,timeout=seconds,creationflags=subprocess.CREATE_NO_WINDOW)
 (W/(name+'.log')).write_bytes(q.stdout)
 (W/(name+'.process.json')).write_text(json.dumps({'exit_code':q.returncode,'output_bytes':len(q.stdout)}),encoding='utf-8')
 if q.returncode:raise RuntimeError('reference command failed: '+name+' exit='+str(q.returncode))
 return q.stdout.decode().strip()
def port():
 with socket.socket() as s:s.bind(('127.0.0.1',0));return s.getsockname()[1]
pgport=port();promport=port();DATA=W/'pgdata';DB='elite_http_test_'+uuid.uuid4().hex
def launch(name,args,extra=None):
 f=(W/(name+'.log')).open('xb');handles.append(f)
 q=subprocess.Popen([str(v) for v in args],cwd=W,env=env|dict(extra or {}),stdin=subprocess.DEVNULL,stdout=f,stderr=subprocess.STDOUT,creationflags=subprocess.CREATE_NO_WINDOW);procs[name]=q;return q
def wait(fn,label,seconds=30):
 end=min(time.monotonic()+seconds,deadline)
 while time.monotonic()<end:
  try:
   value=fn()
   if value:return value
  except (OSError,ValueError,urllib.error.URLError):pass
  time.sleep(.25)
 raise TimeoutError(label)
def sql(query,db=None):
 return command('sql-'+uuid.uuid4().hex,[bins['psql'],'-X','-w','-qAt','-h','127.0.0.1','-p',pgport,'-U','postgres','-d',db or DB,'-v','ON_ERROR_STOP=1','-c',query],20)
def pg_ready():
 if procs['postgres'].poll() is not None:raise RuntimeError('owned PostgreSQL exited')
 try:observed=sql('show data_directory','postgres')
 except RuntimeError:return False
 assert Path(observed).resolve()==DATA.resolve()
 return True
def get(url,context=None):
 with urllib.request.urlopen(url,context=context,timeout=5) as r:return r.read(2_000_001)
def api(path):return json.loads(get('https://127.0.0.1:'+str(promport)+path,tls))
def alert():
 return [x for x in api('/api/v1/alerts')['data']['alerts'] if x['labels']['alertname']=='PlatformHighErrorRate']
def firing():return any(x['state']=='firing' for x in alert())
def request(key,token='synthetic-local-reference-only',body=None):
 data=body or b'{"OrganizationID":"integration-org","Currency":"USD","TotalMinorUnits":100}'
 req=urllib.request.Request('http://'+address['api']+'/v1/orders?contact=PRIVATE_QUERY@example.invalid',data=data,headers={'Authorization':'Bearer '+token,'Content-Type':'application/json','Idempotency-Key':key,'Cookie':'PRIVATE_COOKIE','X-Private':'PRIVATE_HEADER'},method='POST')
 try:
  with urllib.request.urlopen(req,timeout=15) as r:status=r.status;payload=r.read(16384)
 except urllib.error.HTTPError as e:status=e.code;payload=e.read(16384)
 result['requests'][str(status)]=result['requests'].get(str(status),0)+1
 assert b'PRIVATE_' not in payload
 return status,payload
def verify_metrics(label):
 data=get('http://'+address['metrics']+'/metrics')
 assert len(data)<=2_000_000
 for token in [b'PRIVATE_',b'integration-org',b'synthetic-local-reference-only',b'customer_order',b'Authorization',b'http_route',b'server_address',b'url_',b'target_info',b'otel_scope']:
  assert token not in data,token
 (W/(label+'.prom')).write_bytes(data)
 return data
try:
 command('certificates',[bins['telemetry-reference'],'--certificates',W])
 tls=ssl.create_default_context(cafile=str(W/'ca.pem'));tls.minimum_version=ssl.TLSVersion.TLSv1_3;tls.load_cert_chain(str(W/'client.pem'),str(W/'client.key'))
 command('initdb',[bins['initdb'],'-D',DATA,'-U','postgres','-A','trust','--no-locale','-E','UTF8'])
 with (DATA/'postgresql.conf').open('a') as f:f.write(f"\nlisten_addresses='127.0.0.1'\nport={pgport}\nfsync=on\nsynchronous_commit=on\nfull_page_writes=on\n")
 launch('postgres',[bins['postgres'],'-D',DATA]);wait(pg_ready,'owned postgres')
 sql('create database '+DB,'postgres')
 for m in migrations:command('migration-'+m.stem,[bins['psql'],'-X','-w','-qAt','-h','127.0.0.1','-p',pgport,'-U','postgres','-d',DB,'-v','ON_ERROR_STOP=1','-f',m])
 sql("insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values('018f4d4a-7b36-7a21-8d10-2f4c54c28b01','http-reference','Synthetic Local Only','Synthetic Local Only'); insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values('018f4d4a-7b36-7a21-8d10-2f4c54c28b01','integration-org','reference','Synthetic Local Only','enterprise')")
 launch('host',[host,'--ttl','20m'],{'REFERENCE_DATABASE_URL':f'postgres://postgres@127.0.0.1:{pgport}/{DB}?sslmode=disable'})
 def host_ready():
  if procs['host'].poll() is not None:raise RuntimeError('reference host exited')
  text=(W/'host.log').read_text()
  if '\n' not in text:return None
  return json.loads(text.splitlines()[0])
 address=wait(host_ready,'host startup')
 assert set(address)=={'api','metrics'} and all(re.fullmatch(r'127\.0\.0\.1:[0-9]{1,5}',v) for v in address.values())
 first=request('reference-http-order-0001');assert first[0]==201,first
 again=request('reference-http-order-0001');assert again[0]==200 and json.loads(first[1])==json.loads(again[1])
 assert request('reference-http-unauth-0001','PRIVATE_TOKEN')[0]==401
 assert request('reference-http-invalid-0001',body=b'{"PRIVATE_BODY":true}')[0]==400
 assert sql("select count(*) from sales.customer_order")=='1'
 verify_metrics('healthy')
 (W/'prometheus.json').write_text(json.dumps({'global':{'scrape_interval':'1s','evaluation_interval':'30s'},'rule_files':[str(rules)],'scrape_configs':[{'job_name':'reference-http','static_configs':[{'targets':[address['metrics']]}]}]},indent=2))
 (W/'web.json').write_text(json.dumps({'tls_server_config':{'cert_file':str(W/'server.pem'),'key_file':str(W/'server.key'),'client_ca_file':str(W/'ca.pem'),'client_auth_type':'RequireAndVerifyClientCert','min_version':'TLS13'}}))
 command('rules-check',[bins['promtool'],'check','config',W/'prometheus.json'])
 launch('prometheus',[bins['prometheus'],'--config.file='+str(W/'prometheus.json'),'--web.listen-address=127.0.0.1:'+str(promport),'--web.config.file='+str(W/'web.json'),'--storage.tsdb.path='+str(W/'tsdb'),'--storage.tsdb.retention.time=1h','--storage.tsdb.retention.size=64MB','--log.level=error'])
 wait(lambda:api('/api/v1/status/config')['status']=='success','prometheus startup')
 wait(lambda:api('/api/v1/query?query='+urllib.parse.quote('up{job="reference-http"}'))['data']['result'],'first scrape')
 time.sleep(3)
 sql('alter table sales.customer_order rename to customer_order_reference_fault')
 fault=time.monotonic();next_error=fault;pending=None;fired=None;last_state=None;error_count=0
 log('fault_injected',method='rename owned synthetic order table',interval_seconds=240,rule_hold_seconds=600)
 while time.monotonic()<deadline-100:
  if time.monotonic()>=next_error:
   assert request('reference-http-failure-'+str(error_count).zfill(5))[0]==500
   error_count+=1;next_error+=240;verify_metrics('fault-'+str(error_count))
   log('actual_http_500',count=error_count)
  active=alert();state=active[0]['state'] if active else 'inactive'
  if state!=last_state:log('alert_state',state=state);last_state=state
  if state=='pending' and pending is None:pending=time.monotonic()
  if state=='firing':
   assert pending is not None and time.monotonic()-pending>=590
   fired=time.monotonic();result['firing_alert']=active;break
  assert procs['host'].poll() is None and procs['prometheus'].poll() is None
  time.sleep(1)
 assert fired is not None,'alert did not fire within finite budget'
 expr='sum(rate(http_server_request_duration_seconds_count{http_response_status_code=~"5.."}[5m])) / sum(rate(http_server_request_duration_seconds_count[5m]))'
 current=api('/api/v1/query?query='+urllib.parse.quote(expr));legacy=api('/api/v1/query?query='+urllib.parse.quote(expr.replace('/ sum(rate(http_server_request_duration_seconds_count[5m]))','/ clamp_min(sum(rate(http_server_request_duration_seconds_count[5m])), 1)')))
 result['observed_ratios']={'corrected':current,'legacy':legacy}
 assert float(current['data']['result'][0]['value'][1])>0.99
 assert float(legacy['data']['result'][0]['value'][1])<0.01, 'legacy window ratio was not below threshold'
 sql('alter table sales.customer_order_reference_fault rename to customer_order')
 for i in range(500):assert request('reference-http-order-0001')[0]==200
 wait(lambda:not firing(),'alert recovery',70)
 assert sql('select count(*) from sales.customer_order')=='1'
 assert sql('select count(*) from platform.idempotency_record')=='1'
 assert sql("select count(*) from platform.outbox_event where event_type='order.created'")=='1'
 verify_metrics('recovered')
 result.update(state='PASS',actual_sparse_failures=error_count,pending_to_firing_seconds=round(fired-pending,3),fault_to_firing_seconds=round(fired-fault,3),same_order_after_recovery=True,order_count=1,outbox_count=1,idempotency_count=1)
 log('recovered',order_count=1,outbox_count=1,idempotency_count=1)
except BaseException as e:
 result.update(state='FAIL',failure_type=type(e).__name__,failure=str(e))
 raise
finally:
 if 'host' in procs and procs['host'].poll() is None:procs['host'].terminate();procs['host'].wait(10)
 if 'prometheus' in procs and procs['prometheus'].poll() is None:procs['prometheus'].terminate();procs['prometheus'].wait(10)
 if 'postgres' in procs and procs['postgres'].poll() is None:command('postgres-stop',[bins['pg_ctl'],'-D',DATA,'-m','fast','-w','stop']);procs['postgres'].wait(10)
 for f in handles:f.close()
 result['children_exited']=all(q.poll() is not None for q in procs.values());result['elapsed_seconds']=round(time.monotonic()-start,3)
 result['event_log_sha256']=sha(W/'progress.json') if (W/'progress.json').exists() else None
 (W/'result.json').write_text(json.dumps(result,indent=2)+'\n',encoding='utf-8',newline='\n')
 print('HTTP_REFERENCE_RESULT '+json.dumps(result),flush=True)
````

### FILE: `reference_http_metrics/source-lock.json`

```yaml
block_id: "GO-HTTP-METRICS-REFERENCE:source-lock.json:v1"
operation: CREATE
provenance: AUTHORED
source: "local composition over pinned official dependency APIs; no upstream runtime source copied"
license: "LicenseRef-Workspace-Owner"
sha256: "08853405d56323b7f16905fe66ea4574f0ee17eb53b74a71b375b9a3d7d34ac2"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite-http-metrics-reference-source/v1",
  "toolchain": "go1.26.8",
  "upstreams": [
    {
      "module": "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp",
      "version": "v0.70.0",
      "origin": {
        "VCS": "git",
        "URL": "https://github.com/open-telemetry/opentelemetry-go-contrib",
        "Subdir": "instrumentation/net/http/otelhttp",
        "Hash": "c8a87a60ba1b3374fd16df11fc3eeae6c41abbc9",
        "Ref": "refs/tags/instrumentation/net/http/otelhttp/v0.70.0"
      },
      "zip_sha256": "5afae012c182ed95c9068c4b4de772677fbfecb9fbe41605046f907abba8bc92",
      "go_sum": "h1:LMuyCAyfalSjDyjdC65nK6N0zoTT63+E/u95X0JovZI=",
      "files_verified_against_zip": 36,
      "license_sha256": "1ae07514be1d7bb33f0698f8d91fb51b8b9fe1463157ec1c72081a49b9bc6f40",
      "license_expression": "Apache-2.0 AND BSD-3-Clause",
      "runtime_source_changed": false,
      "cached_source_reuse": true
    },
    {
      "module": "go.opentelemetry.io/otel/exporters/prometheus",
      "version": "v0.67.0",
      "origin": {
        "VCS": "git",
        "URL": "https://github.com/open-telemetry/opentelemetry-go",
        "Subdir": "exporters/prometheus",
        "Hash": "93a693edeed0e07ce5ebd1dfe67af42d1e2055d8",
        "Ref": "refs/tags/exporters/prometheus/v0.67.0"
      },
      "zip_sha256": "ce725796160333dcee033dc9e6a2a81ea8d6fb24846f9ddbdc714c5bc86d6d43",
      "go_sum": "h1:7IefDa35e6V3NoiqIeLDMDxMFyZDk5qcoC0Ax4cC16E=",
      "files_verified_against_zip": 63,
      "license_sha256": "1ae07514be1d7bb33f0698f8d91fb51b8b9fe1463157ec1c72081a49b9bc6f40",
      "license_expression": "Apache-2.0",
      "runtime_source_changed": false
    }
  ],
  "public_graph": [
    {
      "package": {
        "ecosystem": "Go",
        "name": "cloud.google.com/go/compute/metadata"
      },
      "version": "v0.3.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/alecthomas/kingpin/v2"
      },
      "version": "v2.4.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/alecthomas/units"
      },
      "version": "v0.0.0-20240927000941-0f3dac36c52b"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/beorn7/perks"
      },
      "version": "v1.0.1"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/cespare/xxhash/v2"
      },
      "version": "v2.3.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/coreos/go-oidc/v3"
      },
      "version": "v3.20.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/davecgh/go-spew"
      },
      "version": "v1.1.1"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/felixge/httpsnoop"
      },
      "version": "v1.1.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/go-jose/go-jose/v4"
      },
      "version": "v4.1.4"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/go-logr/logr"
      },
      "version": "v1.4.4"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/go-logr/stdr"
      },
      "version": "v1.2.2"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/golang-jwt/jwt/v5"
      },
      "version": "v5.3.1"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/golang/protobuf"
      },
      "version": "v1.5.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/google/go-cmp"
      },
      "version": "v0.7.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/google/uuid"
      },
      "version": "v1.6.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/jackc/pgpassfile"
      },
      "version": "v1.0.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/jackc/pgservicefile"
      },
      "version": "v0.0.0-20240606120523-5a60cdf6a761"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/jackc/pgx/v5"
      },
      "version": "v5.10.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/jackc/puddle/v2"
      },
      "version": "v2.2.2"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/jpillora/backoff"
      },
      "version": "v1.0.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/json-iterator/go"
      },
      "version": "v1.1.12"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/julienschmidt/httprouter"
      },
      "version": "v1.3.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/klauspost/compress"
      },
      "version": "v1.19.1"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/kr/pretty"
      },
      "version": "v0.3.1"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/kylelemons/godebug"
      },
      "version": "v1.1.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/modern-go/concurrent"
      },
      "version": "v0.0.0-20180306012644-bacd9c7ef1dd"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/modern-go/reflect2"
      },
      "version": "v1.0.2"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/munnerz/goautoneg"
      },
      "version": "v0.0.0-20191010083416-a7dc8b61c822"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/mwitkow/go-conntrack"
      },
      "version": "v0.0.0-20190716064945-2f068394615f"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/pmezard/go-difflib"
      },
      "version": "v1.0.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/prometheus/client_golang"
      },
      "version": "v1.24.1"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/prometheus/client_model"
      },
      "version": "v0.6.2"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/prometheus/common"
      },
      "version": "v0.70.1"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/prometheus/otlptranslator"
      },
      "version": "v1.0.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/prometheus/procfs"
      },
      "version": "v0.21.1"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/rogpeppe/go-internal"
      },
      "version": "v1.14.1"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/stretchr/objx"
      },
      "version": "v0.1.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/stretchr/testify"
      },
      "version": "v1.12.1"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "github.com/xhit/go-str2duration/v2"
      },
      "version": "v2.1.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "go.opentelemetry.io/auto/sdk"
      },
      "version": "v1.2.1"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
      },
      "version": "v0.70.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "go.opentelemetry.io/otel"
      },
      "version": "v1.46.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "go.opentelemetry.io/otel/exporters/prometheus"
      },
      "version": "v0.67.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "go.opentelemetry.io/otel/metric"
      },
      "version": "v1.46.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "go.opentelemetry.io/otel/metric/x"
      },
      "version": "v0.68.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "go.opentelemetry.io/otel/sdk"
      },
      "version": "v1.46.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "go.opentelemetry.io/otel/sdk/metric"
      },
      "version": "v1.46.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "go.opentelemetry.io/otel/trace"
      },
      "version": "v1.46.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "go.uber.org/goleak"
      },
      "version": "v1.3.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "go.yaml.in/yaml/v2"
      },
      "version": "v2.4.4"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "go.yaml.in/yaml/v3"
      },
      "version": "v3.0.5"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "golang.org/x/mod"
      },
      "version": "v0.40.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "golang.org/x/net"
      },
      "version": "v0.57.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "golang.org/x/oauth2"
      },
      "version": "v0.36.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "golang.org/x/sync"
      },
      "version": "v0.22.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "golang.org/x/sys"
      },
      "version": "v0.47.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "golang.org/x/text"
      },
      "version": "v0.40.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "golang.org/x/tools"
      },
      "version": "v0.47.0"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "google.golang.org/protobuf"
      },
      "version": "v1.36.11"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "gopkg.in/check.v1"
      },
      "version": "v1.0.0-20201130134442-10cb98267c6c"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "gopkg.in/yaml.v3"
      },
      "version": "v3.0.1"
    },
    {
      "package": {
        "ecosystem": "Go",
        "name": "stdlib"
      },
      "version": "1.26.8"
    }
  ],
  "compiled_license_inventory": [
    {
      "module": "elite.local/enterprise",
      "provenance": "AUTHORED",
      "license": "LicenseRef-Workspace-Owner"
    },
    {
      "module": "elite.local/enterprise/reference_http_metrics",
      "provenance": "AUTHORED",
      "license": "LicenseRef-Workspace-Owner"
    },
    {
      "module": "github.com/beorn7/perks",
      "version": "v1.0.1",
      "sum": "h1:VlbKKnNfV8bJzeqoa4cOKqO6bYr3WgKZxO8Z16+hsOM=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 1058,
          "sha256": "0db7c9ebb3717e526f34f87dd1ee8bc77d36846e29cb0cec9246f7138fbe962b",
          "observed_license_texts": [
            "MIT"
          ]
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "github.com/cespare/xxhash/v2",
      "version": "v2.3.0",
      "sum": "h1:UL815xU9SqsFlibzuggzjXhog7bL6oX9BbNZnL2UFvs=",
      "license_texts": [
        {
          "name": "LICENSE.txt",
          "bytes": 1068,
          "sha256": "f566a9f97bacdaf00d9f21dd991e81dc11201c4e016c86b470799429a1c9a79c",
          "observed_license_texts": [
            "MIT"
          ]
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "github.com/coreos/go-oidc/v3",
      "version": "v3.20.0",
      "sum": "h1:EtE0WIBHk03N+DqGkY4+UONzzZHk7amKt6IyNd7OsZE=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 11325,
          "sha256": "cb5e8e7e5f4a3988e1063c142c60dc2df75605f4c46515e776e3aca6df976e14",
          "observed_license_texts": [
            "Apache-2.0"
          ]
        },
        {
          "name": "NOTICE",
          "bytes": 126,
          "sha256": "dccd26c6fd9c296daf44d0bc56bb4efc566edd4880381b3331c9a63e6e471338",
          "observed_license_texts": []
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "github.com/felixge/httpsnoop",
      "version": "v1.1.0",
      "sum": "h1:3YtUj32ZZkqZtt3sZZsClsymw/QDuVfpNhoA31zeORc=",
      "license_texts": [
        {
          "name": "LICENSE.txt",
          "bytes": 1101,
          "sha256": "8108e14b5cb2b6bcd51583d3849ef0c2b5c4de61758256c20308bc70d1775dbb",
          "observed_license_texts": [
            "MIT"
          ]
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "github.com/go-jose/go-jose/v4",
      "version": "v4.1.4",
      "sum": "h1:moDMcTHmvE6Groj34emNPLs/qtYXRVcd6S7NHbHz3kA=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 11358,
          "sha256": "cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30",
          "observed_license_texts": [
            "Apache-2.0"
          ]
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "github.com/go-logr/logr",
      "version": "v1.4.4",
      "sum": "h1:tG4xh9yMsRCAiodLVTxyrkzSZ9+o0L1Kg/+cPVcbP/8=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 11357,
          "sha256": "b40930bbcf80744c86c46a12bc9da056641d722716c378f5659b9e555ef833e1",
          "observed_license_texts": [
            "Apache-2.0"
          ]
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "github.com/go-logr/stdr",
      "version": "v1.2.2",
      "sum": "h1:hSWxHoqTgW2S2qGc0LTAI563KZ5YKYRhT3MFKZMbjag=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 11357,
          "sha256": "c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4",
          "observed_license_texts": [
            "Apache-2.0"
          ]
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "github.com/google/uuid",
      "version": "v1.6.0",
      "sum": "h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 1480,
          "sha256": "0a8d61ed3cbfd5312326e8126c31ce9c627a283adc99131b56896d29ada04b2d",
          "observed_license_texts": [
            "BSD-family-text-retained"
          ]
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "github.com/jackc/pgpassfile",
      "version": "v1.0.0",
      "sum": "h1:/6Hmqy13Ss2zCq62VdNG8tM1wchn8zjSGOBJ6icpsIM=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 1073,
          "sha256": "adb1663fda031df8f4344aa68f299fd87d80353e31339406742ded21dae65702",
          "observed_license_texts": [
            "MIT"
          ]
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "github.com/jackc/pgservicefile",
      "version": "v0.0.0-20240606120523-5a60cdf6a761",
      "sum": "h1:iCEnooe7UlwOQYpKFhBabPMi4aNAfoODPEFNiAnClxo=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 1073,
          "sha256": "fc505773403fe869ed64cc2235cdd13988a427bb7e3a7e7004a3f4b27420f8fc",
          "observed_license_texts": [
            "MIT"
          ]
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "github.com/jackc/pgx/v5",
      "version": "v5.10.0",
      "sum": "h1:VhSvgU2jSli8o3AqIEOTJr7rZwAEUVo4E4XhR94Zfr0=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 1078,
          "sha256": "467f95e074fe23079a5623ed652619682692041b8551da27e3c2ddb9659a1507",
          "observed_license_texts": [
            "MIT"
          ]
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "github.com/jackc/puddle/v2",
      "version": "v2.2.2",
      "sum": "h1:PR8nw+E/1w0GLuRFSmiioY6UooMp6KJv0/61nB7icHo=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 1073,
          "sha256": "2d50e98a4900b4d6457a38d39c1432fdc156fc2f7b365f2e33ec9344acbb0057",
          "observed_license_texts": [
            "MIT"
          ]
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "github.com/munnerz/goautoneg",
      "version": "v0.0.0-20191010083416-a7dc8b61c822",
      "sum": "h1:C3w9PqII01/Oq1c1nUAm88MOHcQC9l5mIlSMApZMrHA=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 1546,
          "sha256": "aa1376b9bc5dea6f30cdefde40c176a254f247d2814d0a9929395138631b2ae0",
          "observed_license_texts": [
            "BSD-family-text-retained"
          ]
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "github.com/prometheus/client_golang",
      "version": "v1.24.1",
      "sum": "h1:JnJkREXzWxUdCuPFpIWZiPispT9xVV59uiuyR2bPlnU=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 11357,
          "sha256": "c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4",
          "observed_license_texts": [
            "Apache-2.0"
          ]
        },
        {
          "name": "NOTICE",
          "bytes": 631,
          "sha256": "8548bfb2bc1913396df38c548a6ec60b8dbcac3b0800485a40569c8bdd184471",
          "observed_license_texts": []
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "github.com/prometheus/client_model",
      "version": "v0.6.2",
      "sum": "h1:oBsgwpGs7iVziMvrGhE53c/GrLUsZdHnqNwqPLxwZyk=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 11357,
          "sha256": "c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4",
          "observed_license_texts": [
            "Apache-2.0"
          ]
        },
        {
          "name": "NOTICE",
          "bytes": 167,
          "sha256": "6c79faa15168885fb88a316ae0df18f486deafddefdc16826cdc56dfbd421e24",
          "observed_license_texts": []
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "github.com/prometheus/common",
      "version": "v0.70.1",
      "sum": "h1:1HvjP4D5oL3t8RsPlwxA9onvvStjtIHYE5XuuwOi/PY=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 11357,
          "sha256": "c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4",
          "observed_license_texts": [
            "Apache-2.0"
          ]
        },
        {
          "name": "NOTICE",
          "bytes": 178,
          "sha256": "600244c8052c8c1d307043fafacbc83b429c171dd8e08f5a537c0a54014585ee",
          "observed_license_texts": []
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "github.com/prometheus/otlptranslator",
      "version": "v1.0.0",
      "sum": "h1:s0LJW/iN9dkIH+EnhiD3BlkkP5QVIUVEoIwkU+A6qos=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 11357,
          "sha256": "c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4",
          "observed_license_texts": [
            "Apache-2.0"
          ]
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "go.opentelemetry.io/auto/sdk",
      "version": "v1.2.1",
      "sum": "h1:jXsnJ4Lmnqd11kwkBV2LgLoFMZKizbCi5fNZ/ipaZ64=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 11357,
          "sha256": "c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4",
          "observed_license_texts": [
            "Apache-2.0"
          ]
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp",
      "version": "v0.70.0",
      "sum": "h1:LMuyCAyfalSjDyjdC65nK6N0zoTT63+E/u95X0JovZI=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 12892,
          "sha256": "1ae07514be1d7bb33f0698f8d91fb51b8b9fe1463157ec1c72081a49b9bc6f40",
          "observed_license_texts": [
            "Apache-2.0",
            "BSD-family-text-retained"
          ]
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "go.opentelemetry.io/otel",
      "version": "v1.46.0",
      "sum": "h1:FHt5/CDyVxi/8IM1CH7VE/rRgq3kLHa2mSTVMO8AWyc=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 12892,
          "sha256": "1ae07514be1d7bb33f0698f8d91fb51b8b9fe1463157ec1c72081a49b9bc6f40",
          "observed_license_texts": [
            "Apache-2.0",
            "BSD-family-text-retained"
          ]
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "go.opentelemetry.io/otel/exporters/prometheus",
      "version": "v0.67.0",
      "sum": "h1:7IefDa35e6V3NoiqIeLDMDxMFyZDk5qcoC0Ax4cC16E=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 12892,
          "sha256": "1ae07514be1d7bb33f0698f8d91fb51b8b9fe1463157ec1c72081a49b9bc6f40",
          "observed_license_texts": [
            "Apache-2.0",
            "BSD-family-text-retained"
          ]
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "go.opentelemetry.io/otel/metric",
      "version": "v1.46.0",
      "sum": "h1:yBnkXvgV7AXFILZc5K6IZe/CBFF3OS7BJ8ov6/lj0K8=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 12892,
          "sha256": "1ae07514be1d7bb33f0698f8d91fb51b8b9fe1463157ec1c72081a49b9bc6f40",
          "observed_license_texts": [
            "Apache-2.0",
            "BSD-family-text-retained"
          ]
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "go.opentelemetry.io/otel/sdk",
      "version": "v1.46.0",
      "sum": "h1:h5CNQQjEbuQXY/JfZtgt3i7HVFV3aHPO2OAwO2eTYPI=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 12892,
          "sha256": "1ae07514be1d7bb33f0698f8d91fb51b8b9fe1463157ec1c72081a49b9bc6f40",
          "observed_license_texts": [
            "Apache-2.0",
            "BSD-family-text-retained"
          ]
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "go.opentelemetry.io/otel/sdk/metric",
      "version": "v1.46.0",
      "sum": "h1:0piZ26EG4RBfebb2jhDH6ERCYHoVWduc3kLgPCwSnSE=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 12892,
          "sha256": "1ae07514be1d7bb33f0698f8d91fb51b8b9fe1463157ec1c72081a49b9bc6f40",
          "observed_license_texts": [
            "Apache-2.0",
            "BSD-family-text-retained"
          ]
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "go.opentelemetry.io/otel/trace",
      "version": "v1.46.0",
      "sum": "h1:OULy7ccdJnZtJ0UDYFOIGaCmiWzJ8Vi2G/Rsu60qs1c=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 12892,
          "sha256": "1ae07514be1d7bb33f0698f8d91fb51b8b9fe1463157ec1c72081a49b9bc6f40",
          "observed_license_texts": [
            "Apache-2.0",
            "BSD-family-text-retained"
          ]
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "golang.org/x/oauth2",
      "version": "v0.36.0",
      "sum": "h1:peZ/1z27fi9hUOFCAZaHyrpWG5lwe0RJEEEeH0ThlIs=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 1453,
          "sha256": "911f8f5782931320f5b8d1160a76365b83aea6447ee6c04fa6d5591467db9dad",
          "observed_license_texts": [
            "BSD-family-text-retained"
          ]
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "golang.org/x/sync",
      "version": "v0.22.0",
      "sum": "h1:SZjpbeLmrCk4xhRSZFNZW5gFUeCeFgjekvI/+gfScek=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 1453,
          "sha256": "911f8f5782931320f5b8d1160a76365b83aea6447ee6c04fa6d5591467db9dad",
          "observed_license_texts": [
            "BSD-family-text-retained"
          ]
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "golang.org/x/sys",
      "version": "v0.47.0",
      "sum": "h1:o7XGOvZQCADBQQ4Y7VNq2dRWQR7JmOUW8Kxx4ZsNgWs=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 1453,
          "sha256": "911f8f5782931320f5b8d1160a76365b83aea6447ee6c04fa6d5591467db9dad",
          "observed_license_texts": [
            "BSD-family-text-retained"
          ]
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "golang.org/x/text",
      "version": "v0.40.0",
      "sum": "h1:Ub2Z6/xjgF1WrYQz2nuITOEegKFtiIy+rieRJ5lHZKs=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 1453,
          "sha256": "911f8f5782931320f5b8d1160a76365b83aea6447ee6c04fa6d5591467db9dad",
          "observed_license_texts": [
            "BSD-family-text-retained"
          ]
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    },
    {
      "module": "google.golang.org/protobuf",
      "version": "v1.36.11",
      "sum": "h1:fV6ZwhNocDyBLK0dj+fg8ektcVegBBuEolpbTQyBNVE=",
      "license_texts": [
        {
          "name": "LICENSE",
          "bytes": 1479,
          "sha256": "4835612df0098ca95f8e7d9e3bffcb02358d435dbb38057c844c99d7f725eb20",
          "observed_license_texts": [
            "BSD-family-text-retained"
          ]
        }
      ],
      "local_source_verified_by_go_mod_verify": true
    }
  ],
  "backend": {
    "pack_id": "GO-ENTERPRISE-BACKEND",
    "version": "0.4.4",
    "pack_sha256": "69ae4b58a4edc0e5b4b2a7d037624a637b585c8a38cb5c53c4ba3c5785021c5e"
  },
  "rules": {
    "pack_id": "SECURE-OPS-DELIVERY-CORE",
    "version": "1.1.4",
    "sha256": "cc718b7bbe35de3f230eeb47c854c270580dbb488101ee5f30c0b758e1d967b0"
  },
  "tested_binary_sha256": "4643654625c71aaf5b45367bd3a198381257588d073ae0349fc90c1bc1b5ab85",
  "authority_claim": "Official middleware/exporter unchanged; SDK1.46.0 selection and reference glue AUTHORED. No production identity, complete security or all-service SLO admission.",
  "parent_composition": {
    "plan": "markdown_system/FRANCHISE_COMPLETE_PACK_PLAN.md",
    "sha256": "8bb5eaac51e36e0a6ce8811415f0c61f4e1a744d76db6f40ada7b3cb0f91ccb2",
    "packs": 68,
    "files": 797,
    "compiler_input_parity_files": 441,
    "require_complete_composed_packages": true,
    "historical_v375_plan_sha256": "53ee063cc9660727f9a9d92c2d773c6a018b9492750e99f52d10e4a4ced01d66",
    "v376_change": "Only BFF package.json and pnpm-lock.yaml changed in the parent; all Go/SQL inputs remain byte-identical. Rebuild must match the recorded V375 executable before carrying its runtime evidence.",
    "historical_v376_plan_sha256": "96341c87a51729255d91b3480b9934753ff55015869d350e33b0fa556967b612",
    "v377_change": "BFF public indexing adds5files and modifies3page files. No dependency, Go, SQL or alert-rule bytes change; preserve the exact V375 executable and prior narrow runtime evidence.",
    "historical_v377_plan_sha256": "4bf15a615bf973f0dd54040f8eccaff728ea9253a72b704988c2560b14f52432",
    "v378_change": "Public presentation locale/timeZone integration under existing BFF/portals/browser owners:4newfiles,8changedfiles; all Go/SQL/dependency/rule source bytes unchanged.",
    "historical_v378_plan_sha256": "1534632f07eb43b706682573a1accf2fb602078196fe86baef473a363d0acdc3",
    "v379_change": "Read-only role-filtered versioned operational help:10newfiles,6changedfiles under existing BFF/portals/browser owners; domain Go/SQL/dependencies/rules unchanged.",
    "historical_v379_plan_sha256": "c216ecdf5863d591706790a993339650d13b20342451c5b3ee6d282f39a645c0",
    "v380_change": "Complete15existing versioned guide index and15role matrix:3newfiles under existing portal owner; domain Go/SQL/dependencies unchanged.",
    "historical_v380_plan_sha256": "982e031b3f44316c302eb4b1526da75240be9a4633ac74895b38a8f5cabd63c9",
    "v381_change": "Existing seven-file role workspace selected; permission-filtered BFF navigation and one browser test. Domain Go/SQL/dependency artifacts unchanged.",
    "historical_v381_plan_sha256": "f5531fbfbe9f4de27d95a8fa311872d8f56ad31e10bbbeb36ba0d7399604dc47",
    "v382_change": "Private admin/customer read permission and display repair; existing quote formatter shared unchanged; five new test/helper/document files. Production Go, SQL and dependency artifacts unchanged.",
    "historical_v382_plan_sha256": "40763eb2a2e3afaf861b774f6456b7d2dd372d9c96b3fe75bd0fb0313eda16fd",
    "v383_change": "Private read error recovery; five authored component/boundary/test files; Go, SQL and dependency artifacts unchanged.",
    "historical_v383_plan_sha256": "7343dbdeaabba64897e764c5af0beb8fce093c8850266f6c21209745f812218c",
    "v387_change": "Three existing customer/factory lists gain cursor navigation;2newAUTHORED files. HTTP runtime, migrations, business and dependencies unchanged.",
    "historical_v387_plan_sha256": "e705d88e6b03bd6e7b336b3b6df34cac27e28da8c11289b3d4bcfe718baec6c0",
    "v388_change": "Admin cursor traversal of existing orders/cases/leads using V387 helper; runtime/domain/SQL/dependency bytes unchanged.",
    "historical_v388_plan_sha256": "3461f321aec5c70da7573aa4eb095239e296812832ea9edf90abaa9eced0b437",
    "v389_change": "Customer appointment display uses existing configured locale/zone; domain/SQL/runtime/dependency bytes unchanged.",
    "historical_v389_plan_sha256": "96fba9178aee75a632d8753b9ef46bce1b5fcdc432676b3ffc14e13f7217f2ea",
    "v390_change": "Customer cancellation UI recovery and bound receipts; domain/SQL/runtime/dependency bytes unchanged.",
    "historical_v390_plan_sha256": "dae1449cea384482e0d439262207f9fe65b1455325364c52e15f6e7379a186bc",
    "v393_change": "Unused optional Sharp closure omitted and runtime image optimizer disabled; Go/domain/SQL bytes unchanged.",
    "historical_v393_plan_sha256": "c6c7b155c0f9e9f3a7e73b95478583e0488c1f0aa58feb7d952c0719d3ba9b5e",
    "v400_change": "Opt-in customer surveys:23new authored source files and application activation; existing metrics code unchanged.",
    "historical_v400_initial_plan_sha256": "e953f1fef5c064483d2ac666608d791fa82fedc7b8c7e52816d8a43dc7ee7c2e",
    "v400_atomic_migration": "Survey migrations now transactional; only2SQLsources differ from the initial connected reference."
  }
}
````

## 6. Configuration surface

Runtime lock and SHA, consumer root, rebuilt host and SHA, canonical rules and SHA, private existing work root. No consumer secrets or approvals bundled.

## 7. Dependency bill

Exact nested go.mod/go.sum and source-lock.json; SDK1.46.0, otelhttp0.70.0 and exporter0.67.0. Preserve all upstream notices on redistribution. Library ships authored glue/locks, not upstream runtime binaries. No floating updates.

## 8. Apply order

Reconstruct into an absent destination, compose the compatible parent module, run Go tests/vet/build with readonly dependencies, check exact artifact identity and execute run_reference.py. V375 records the actual evidence and failed predecessors. Structural materialization is not runtime PASS.

## 9. Verification

V375 records 238 middleware and 142 exporter official tests, two authored HTTP privacy/isolation tests, vet/build, module verification, 62 public-version SCA queries and the actual finite fault/alert/recovery run. The final 68-pack consumer rebuilds the tested executable byte-identically. Revalidate exact source/runtime/rule hashes and all target conditions before using the reference in another environment.

## 10. Reconstruction evidence

See `reconstruction_evidence/HTTP_SLI_ALERT_INTEGRATION_V375.md` for exact receipts, source origins, failed predecessors and limits. Final library Preflight is recorded there separately from runtime proof. This is conditioned synthetic qualification, not production acceptance.

Rollback and limits:

Remove this opt-in reference selection by returning to the prior source profile. No production process or database is touched by a reference run. Retain failed directories. Do not substitute this scope for integral security, all-service SLOs or TEST02 completion.
