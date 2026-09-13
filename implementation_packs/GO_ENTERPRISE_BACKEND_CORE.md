# Go Enterprise Backend Core

V316: outbox usa attempts monotónico como generación; MarkPublished/Release
exigen generación y lease vigente después del row lock. RemainingLease revalida
al publisher y limita su contexto por evento, incluyendo entradas tardías de lote.
Core0.4.4 y workers0.3.1 se reconstruyen juntos: cambian firmas internas de store;
no mezclar con consumers anteriores ni reiniciar attempts. Sin migración nueva.
Tests PostgreSQL requieren OUTBOX_TEST_DATABASE_URL y
OUTBOX_PROCESSOR_TEST_DATABASE_URL en DB loopback distintas elite_outbox_test_<unique>.
TEST_DATABASE_URL genérica no habilita esas pruebas; no aceptar SKIP como PASS.
Evidencia OUTBOX_GENERATION_FENCING_V316.md: seis rojos DB,tres publisher,
33 tests raíz por corte conectado y4/4 rebuild. Entrega externa sigue al menos
una vez: publisher debe honrar contexto/idempotency key y reconciliar ambigüedad.
No exactly-once remoto, supervisor/alerta, dead-letter policy o producción probados.

V314: ambos hosts esperan el drenaje HTTP antes de retornar. Shutdown bloquea
esta ruta hasta completar requests o deadline15s; al vencer cierra conexiones
y devuelve el fallo. Listener/startup errors conservan causa. Cuatro regresiones
HTTP TCP y una de ciclo con seis negativos; evidencia HTTP_HOST_SHUTDOWN_V314.md.
No acredita señales OS, target, WebSockets/hijacked ni handlers que ignoren
cancelación. Sin dependencia ni regla de negocio nueva; no promoción integral.

V313: actualización de seguridad por GO-2026-5970/CVE-2026-56852.
x/text0.29.0 ->0.39.0, x/sync0.17.0 ->0.21.0 y piso explícito
x/mod0.40.0 por GO-2026-6179/6180; conservar override al mantener locks;
artefactos oficiales Go con sumdb, commit y BSD-3-Clause verificados.
Sin backport local ni cambio de negocio. El fixture de reembolso exige DB
descartable exclusiva; su seed no prueba constraints del journey de origen.
Evidencia reconstruction_evidence/COMPOSITION_DEPENDENCY_SECURITY_V313.md.
No promoción integral, runtime monitor ni aprobación de producción.

V313: actualización de seguridad por GO-2026-5970/CVE-2026-56852.
x/text0.29.0 ->0.39.0 y mínimo transitivo x/sync0.17.0 ->0.21.0;
artefactos oficiales Go con sumdb, commit y BSD-3-Clause verificados.
Sin backport local ni cambio de negocio. El fixture de reembolso exige DB
descartable exclusiva; su seed no prueba constraints del journey de origen.
Evidencia reconstruction_evidence/COMPOSITION_DEPENDENCY_SECURITY_V313.md.
No promoción integral, runtime monitor ni aprobación de producción.

## 1. Metadata

```yaml
pack_id: "GO-ENTERPRISE-BACKEND"
pack_version: "0.4.8"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Backend compilado multi-tenant para dominio transaccional, OIDC estricto, autorización, HTTP, PostgreSQL, optimistic concurrency, idempotencia y outbox con claim/lease concurrente."
stacks: ["Go 1.26.7", "PostgreSQL 18.6", "pgx 5.10.0", "go-oidc 3.20.0"]
compatible_with: ["PBC-CORE", "PG-TX-FOUNDATION"]
incompatible_with: ["uso de gen_random_uuid sin extensión pgcrypto o adaptación a UUID generado por aplicación"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0 AND MIT dependencies"
upstream_sources: ["https://go.dev", "https://github.com/jackc/pgx", "https://github.com/coreos/go-oidc"]
verified_at: "2026-09-08"
```

## 2. Applicability

Backend de referencia para servicios empresariales nuevos. No es obligatorio: el blueprint puede elegir JVM o Rust con evidencia superior. TypeScript no forma parte de este core. La revisión 0.4 añade `organization_ids` al principal verificado y autorización fail-closed por recurso; las migrations y el issuer real todavía deben probarse al componer un proyecto.

## 3. Architecture contract

Dominio puro → application service → ports → PostgreSQL adapter / HTTP adapter. Invariantes y permisos viven en dominio; la persistencia usa optimistic concurrency y escribe outbox en la misma transacción. HTTP limita body, rechaza campos desconocidos, usa idempotency key y problemas estables.

## 4. Exact file manifest

```text
CREATE cmd/api/main.go
CREATE go.mod
CREATE go.sum
CREATE internal/order/model_test.go
CREATE internal/order/model.go
CREATE internal/order/service.go
CREATE internal/platform/httpapi/context.go
CREATE internal/platform/httpapi/server_test.go
CREATE internal/platform/httpapi/server.go
CREATE internal/platform/identity/oidc.go
CREATE internal/platform/identity/principal_test.go
CREATE internal/platform/postgres/outbox_integration_test.go
CREATE internal/platform/postgres/outbox.go
CREATE internal/platform/postgres/orders_integration_test.go
CREATE internal/platform/postgres/orders.go
CREATE db/migrations/0002_enterprise_order_core.up.sql
CREATE db/migrations/0002_enterprise_order_core.down.sql
CREATE db/tests/0002_enterprise_order_core.test.sql
```

## 5. Materialization blocks

### FILE: `cmd/api/main.go`

```yaml
block_id: "GO-ENTERPRISE-BACKEND:cmd-api-main-go:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "98f2277d4ddafb998c3500c663856f7af524d61d4c9bdd9b2ee60db55954fbae"
variables: []
secrets_allowed: false
```

````go
package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"elite.local/enterprise/internal/order"
	"elite.local/enterprise/internal/platform/httpapi"
	"elite.local/enterprise/internal/platform/identity"
	"elite.local/enterprise/internal/platform/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ids struct{}

func (ids) New() string {
	var value [16]byte
	if _, err := rand.Read(value[:]); err != nil {
		panic(err)
	}
	return hex.EncodeToString(value[:])
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		slog.Error("DATABASE_URL is required")
		os.Exit(2)
	}
	issuer, audience := os.Getenv("OIDC_ISSUER"), os.Getenv("OIDC_AUDIENCE")
	if issuer == "" || audience == "" {
		slog.Error("OIDC_ISSUER and OIDC_AUDIENCE are required")
		os.Exit(2)
	}
	verifier, err := identity.NewOIDCVerifier(ctx, issuer, audience)
	if err != nil {
		slog.Error("OIDC discovery failed")
		os.Exit(1)
	}
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		slog.Error("database configuration failed")
		os.Exit(2)
	}
	defer pool.Close()
	if err = pool.Ping(ctx); err != nil {
		slog.Error("database unavailable")
		os.Exit(1)
	}
	handler := httpapi.New(order.NewService(postgres.NewOrders(pool), ids{}), verifier)
	server := &http.Server{Addr: ":8080", Handler: handler, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 1 << 20}
	slog.Info("api starting", "address", server.Addr)
	if err = httpapi.ServeUntilShutdown(ctx, server, 15*time.Second); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("server failed")
		os.Exit(1)
	}
}
````

### FILE: `go.mod`

```yaml
block_id: "GO-ENTERPRISE-BACKEND:go-mod:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "d2cf3b03eeaaee0113944651327731d6dae7d294d0bc1448f34609c4a6a2d810"
variables: []
secrets_allowed: false
```

````text
module elite.local/enterprise

go 1.26.0

toolchain go1.26.8

require (
	github.com/coreos/go-oidc/v3 v3.20.0
	github.com/jackc/pgx/v5 v5.10.0
)

require (
	github.com/go-jose/go-jose/v4 v4.1.4
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	golang.org/x/oauth2 v0.36.0
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/text v0.40.0 // indirect
)

require (
	example.com/elite/official-payment-webhooks v0.0.0
	github.com/mercadopago/sdk-go v1.14.0
	github.com/stripe/stripe-go/v86 v86.3.0
)

require (
	github.com/aws/aws-sdk-go-v2/credentials v1.19.37 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.38 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.40 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.40 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.39 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.38 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.5.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.33.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.38.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.45.7 // indirect
	github.com/aws/smithy-go v1.28.1 // indirect
	github.com/google/uuid v1.6.0
)

replace example.com/elite/official-payment-webhooks => ./official_payment_webhooks

require (
	example.com/elite/aws-textract-document-runtime v0.0.0
	github.com/aws/aws-sdk-go-v2 v1.44.0
	github.com/aws/aws-sdk-go-v2/config v1.32.38
	github.com/aws/aws-sdk-go-v2/service/textract v1.45.0
)

replace example.com/elite/aws-textract-document-runtime => ./aws_textract_runtime

require elite.local/enterprise/reference_http_metrics v0.0.0

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/felixge/httpsnoop v1.1.0 // indirect
	github.com/go-logr/logr v1.4.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_golang v1.24.1 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.70.1 // indirect
	github.com/prometheus/otlptranslator v1.0.0 // indirect
	github.com/prometheus/procfs v0.21.1 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.70.0 // indirect
	go.opentelemetry.io/otel v1.46.0 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.67.0 // indirect
	go.opentelemetry.io/otel/metric v1.46.0 // indirect
	go.opentelemetry.io/otel/sdk v1.46.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.46.0 // indirect
	go.opentelemetry.io/otel/trace v1.46.0 // indirect
	golang.org/x/mod v0.40.0
	golang.org/x/sys v0.47.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace elite.local/enterprise/reference_http_metrics => ./reference_http_metrics
````

### FILE: `go.sum`

```yaml
block_id: "GO-ENTERPRISE-BACKEND:go-sum:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "10196842152d7c38ab4fdf0f3173dbe22a62e962639ac3c8cab73320731bb715"
variables: []
secrets_allowed: false
```

````text
github.com/aws/aws-sdk-go-v2 v1.44.0 h1:4IbaHhtzy+4h37z4JQyO9a2QsiCml3CNYHtq5hIHigo=
github.com/aws/aws-sdk-go-v2 v1.44.0/go.mod h1:bttEH6JqnUL8LepvDVfdrds/fZ5bCIxzpe3abyUrhDU=
github.com/aws/aws-sdk-go-v2/config v1.32.38 h1:n4yPHBjtQ3BrIIUyk0/LAqf/BL2iv0Tw6XZcMRzM0ps=
github.com/aws/aws-sdk-go-v2/config v1.32.38/go.mod h1:dencYsOS1R7rBy8zehCvwBYzdxxL4Q/nRK7In03wjN8=
github.com/aws/aws-sdk-go-v2/credentials v1.19.37 h1:FJ8Iz4/xISMB/rwLlgfWujfGDFWr0oneQgtA6KPcYLY=
github.com/aws/aws-sdk-go-v2/credentials v1.19.37/go.mod h1:Q6pWOgVUp49x4g5QVi29wHofUoICnZ+Zq4jHbRN/7ec=
github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.38 h1:Nqo2jU1wz5rnBM9XQyXfVD1RP8txkbP3EDx8hR/hbCE=
github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.38/go.mod h1:PzJFHhjR2vWFKHe8HmY5Lxhvwyxnr5MERtk0nDxWNbk=
github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.40 h1:UIXlbijuB2XK1Kr57fo8iIxCuaSHJzwZ1uo+2tbEYIk=
github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.40/go.mod h1:wcEsL6jscjZjVUinb0Q5qD/GXOG1yT3GNfmT9HuDwzU=
github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.40 h1:xLQVRDs2NddDmK9BEyh5KSlJ1Gpy5/GIJXrV6WcVGAE=
github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.40/go.mod h1:XRXnpFVFGLaEVK+olDdFIM1vNa04ETW452oFGEPUxAo=
github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.39 h1:vo4xvMRs/F6h1E52qsgLqCQgWIQXgIJUauG6rlZEh4U=
github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.39/go.mod h1:jB03R1ij/A+OE2e1dz6vgj076gd7vlYcfstAzj3HcnU=
github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.17 h1:OvYZOB3qA6zvfdRFiRFRzVSiElMYrz3GdntkXZxlp1o=
github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.17/go.mod h1:JgR/2Ew50ACfIWau1oeMRX59tMtC0kM+PYQGEaT04cY=
github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.38 h1:H/5TI1jqaHsNoDQ60UwvPvJBg4GURkinXI3Qga29t2w=
github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.38/go.mod h1:PTVFf+XH++7NJOky+RLBYQx0QA5NcaeEYFQ2fsi0nwo=
github.com/aws/aws-sdk-go-v2/service/signin v1.5.7 h1:YcczQ6zNH/ojIzD/ikDrO+RfW06wmdMp18d4NH5hXY4=
github.com/aws/aws-sdk-go-v2/service/signin v1.5.7/go.mod h1:nl9RVnb9ulgAYzOkjLq1NyFxmWcnH2maCUEuOdESy98=
github.com/aws/aws-sdk-go-v2/service/sso v1.33.7 h1:P+bMNiA93gyuYT3Oh+4dWtvrnGcu2bd9Uy5hRJM8BNo=
github.com/aws/aws-sdk-go-v2/service/sso v1.33.7/go.mod h1:zy+397isDFLvleg9H18Zq2MGzMso7uKyJyzR7DWSgFk=
github.com/aws/aws-sdk-go-v2/service/ssooidc v1.38.7 h1:WWkehGZ4nWtOKLMy0yi8+RqzzVqAGe60hGaxwF06JAw=
github.com/aws/aws-sdk-go-v2/service/ssooidc v1.38.7/go.mod h1:T8AI4SbQYm9ybcVmki2T3n7Qg1g3kfWoeQlNwNYOyO8=
github.com/aws/aws-sdk-go-v2/service/sts v1.45.7 h1:yU/9y2r7s9kSUPbHXbpQTa4LA8kt+CMgpu1OBrhx8p4=
github.com/aws/aws-sdk-go-v2/service/sts v1.45.7/go.mod h1:0lQTDEBArMevQXpxu443LVGjKxxEeSsSnrw9n8YiTMg=
github.com/aws/aws-sdk-go-v2/service/textract v1.45.0 h1:es1kFEIsARdt9US5dt5C6N4Nk8YXax5b75UTpDrrnIY=
github.com/aws/aws-sdk-go-v2/service/textract v1.45.0/go.mod h1:xFfDl2uBgGT5AoRgDjSzbZg1BZfqk8NlkIMOJaHulZU=
github.com/aws/smithy-go v1.28.1 h1:R/nXH00c8qcfCzQVELtRw+eLQWtzv+VAIEFJ1/xxXlQ=
github.com/aws/smithy-go v1.28.1/go.mod h1:YE2RhdIuDbA5E5bTdciG9KrW3+TiEONeUWCqxX9i1Fc=
github.com/beorn7/perks v1.0.1 h1:VlbKKnNfV8bJzeqoa4cOKqO6bYr3WgKZxO8Z16+hsOM=
github.com/beorn7/perks v1.0.1/go.mod h1:G2ZrVWU2WbWT9wwq4/hrbKbnv/1ERSJQ0ibhJ6rlkpw=
github.com/cespare/xxhash/v2 v2.3.0 h1:UL815xU9SqsFlibzuggzjXhog7bL6oX9BbNZnL2UFvs=
github.com/cespare/xxhash/v2 v2.3.0/go.mod h1:VGX0DQ3Q6kWi7AoAeZDth3/j3BFtOZR5XLFGgcrjCOs=
github.com/coreos/go-oidc/v3 v3.20.0 h1:EtE0WIBHk03N+DqGkY4+UONzzZHk7amKt6IyNd7OsZE=
github.com/coreos/go-oidc/v3 v3.20.0/go.mod h1:DYCf24+ncYi+XkIH97GY1+dqoRlbaSI26KVTCI9SrY4=
github.com/davecgh/go-spew v1.1.0/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/davecgh/go-spew v1.1.1 h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=
github.com/davecgh/go-spew v1.1.1/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
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
github.com/mercadopago/sdk-go v1.14.0 h1:3PYp9GPa+iysx2lcaKbpBEkXgEw4IIpbw3T6Jhl0IzI=
github.com/mercadopago/sdk-go v1.14.0/go.mod h1:hvQlOYb3MuYPGfjox7jeGBFbeL0nS+iPwUnxAv6yGWI=
github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 h1:C3w9PqII01/Oq1c1nUAm88MOHcQC9l5mIlSMApZMrHA=
github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822/go.mod h1:+n7T8mK8HuQTcFwEeznm/DIxMOiR9yIdICNftLE1DvQ=
github.com/pmezard/go-difflib v1.0.0 h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=
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
github.com/stretchr/testify v1.11.1 h1:7s2iGBzp5EwR7/aIZr8ao5+dra3wiQyKjjFuvgVKu7U=
github.com/stretchr/testify v1.11.1/go.mod h1:wZwfW3scLgRK+23gO65QZefKpKQRnfz6sD981Nm4B6U=
github.com/stretchr/testify v1.12.1 h1:EuwCh5fleGS7H32xRwO3wRGT7DxrDhLAT6FF8MpWDWE=
github.com/stretchr/testify v1.12.1/go.mod h1:MDEgiDPPsNp5cuIrHPPCyornHKgEVbtFUmoNlxoYthg=
github.com/stretchr/testify v1.3.0/go.mod h1:M5WIy9Dh21IEIfnGCwXGc5bZfKNJtfHm1UVUgZn+9EI=
github.com/stretchr/testify v1.7.0/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
github.com/stripe/stripe-go/v86 v86.3.0 h1:BKtYc3NtRa4EGzKAmp4jvl5q7kk2rwMZ+llF18N5vHI=
github.com/stripe/stripe-go/v86 v86.3.0/go.mod h1:Co7QRXCKGNOPTugAdvjgRo+KcMtd9hxy+pZMN0yThsQ=
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
golang.org/x/mod v0.40.0 h1:hUv+3cXcdRHz08UmSiOob7sadHig73uo5bkXxQ/tvUs=
golang.org/x/mod v0.40.0/go.mod h1:0/weTWkPWGBikyTWAX3dkjVztMmBA5hM0DH6BElSupE=
golang.org/x/oauth2 v0.36.0 h1:peZ/1z27fi9hUOFCAZaHyrpWG5lwe0RJEEEeH0ThlIs=
golang.org/x/oauth2 v0.36.0/go.mod h1:YDBUJMTkDnJS+A4BP4eZBjCqtokkg1hODuPjwiGPO7Q=
golang.org/x/sync v0.21.0 h1:HLII4xRRTtCRkxYp4HNFF0Js/Og6q2i++KXbg0gHCwM=
golang.org/x/sync v0.21.0/go.mod h1:9xrNwdLfx4jkKbNva9FpL6vEN7evnE43NNNJQ2LF3+0=
golang.org/x/sync v0.22.0 h1:SZjpbeLmrCk4xhRSZFNZW5gFUeCeFgjekvI/+gfScek=
golang.org/x/sync v0.22.0/go.mod h1:9xrNwdLfx4jkKbNva9FpL6vEN7evnE43NNNJQ2LF3+0=
golang.org/x/sys v0.47.0 h1:o7XGOvZQCADBQQ4Y7VNq2dRWQR7JmOUW8Kxx4ZsNgWs=
golang.org/x/sys v0.47.0/go.mod h1:4GL1E5IUh+htKOUEOaiffhrAeqysfVGipDYzABqnCmw=
golang.org/x/text v0.39.0 h1:UbZz4pLOvn600D6Oh6GGEI6VAmndrEBLv8/6BEXzyus=
golang.org/x/text v0.39.0/go.mod h1:3UwRclnC2g0TU9x8PZiyfOajCd1zaUNHF9cvqcQZ+ZM=
golang.org/x/text v0.40.0 h1:Ub2Z6/xjgF1WrYQz2nuITOEegKFtiIy+rieRJ5lHZKs=
golang.org/x/text v0.40.0/go.mod h1:hpnzDAfGV753zIKo+wk3u1bVKCGPbrnF7+7LBF/UHVY=
google.golang.org/protobuf v1.36.11 h1:fV6ZwhNocDyBLK0dj+fg8ektcVegBBuEolpbTQyBNVE=
google.golang.org/protobuf v1.36.11/go.mod h1:HTf+CrKn2C3g5S8VImy6tdcUvCska2kB7j23XfzDpco=
gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
````

### FILE: `internal/order/model_test.go`

```yaml
block_id: "GO-ENTERPRISE-BACKEND:internal-order-model_test-go:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "698bd5a301f4e1aba622a30b16eff2235c19de9fb7df51f3f58fad29602fa206"
variables: []
secrets_allowed: false
```

````go
package order

import (
	"errors"
	"testing"
)

func TestTransitionRequiresPermissionAndIncrementsVersion(t *testing.T) {
	o, err := New("order-1", "018f4d4a-7b36-7a21-8d10-2f4c54c28a01", "org-1", "customer-1", "USD", 100)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = o.Transition(Placed, map[string]struct{}{}); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected conflict, got %v", err)
	}
	next, err := o.Transition(Placed, map[string]struct{}{"order:create": {}})
	if err != nil {
		t.Fatal(err)
	}
	if next.State != Placed || next.Version != 2 {
		t.Fatalf("unexpected order: %+v", next)
	}
}

func TestTransitionRejectsInvalidStateJump(t *testing.T) {
	o, _ := New("order-1", "018f4d4a-7b36-7a21-8d10-2f4c54c28a01", "org-1", "customer-1", "USD", 100)
	if _, err := o.Transition(Delivered, map[string]struct{}{"*": {}}); !errors.Is(err, ErrConflict) {
		t.Fatalf("expected conflict, got %v", err)
	}
}
````

### FILE: `internal/order/model.go`

```yaml
block_id: "GO-ENTERPRISE-BACKEND:internal-order-model-go:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "8c2e7960d67ab292045c3ac45be1a11b2fa25d385a206e650fc120ac5f1ce8b8"
variables: []
secrets_allowed: false
```

````go
package order

import (
	"errors"
	"fmt"
)

type State string

const (
	Draft     State = "draft"
	Placed    State = "placed"
	Confirmed State = "confirmed"
	Paid      State = "paid"
	Allocated State = "allocated"
	Delivered State = "delivered"
	Cancelled State = "cancelled"
)

var ErrConflict = errors.New("order conflict")

type Order struct {
	ID                string `json:"id"`
	TenantID          string `json:"tenantId"`
	OrganizationID    string `json:"organizationId"`
	CustomerPrincipal string `json:"customerPrincipal"`
	State             State  `json:"state"`
	Currency          string `json:"currency"`
	TotalMinorUnits   int64  `json:"totalMinorUnits"`
	Version           int64  `json:"version"`
}

type transition struct {
	from       State
	to         State
	permission string
}

var transitions = []transition{
	{Draft, Placed, "order:create"},
	{Placed, Confirmed, "order:transition"},
	{Confirmed, Paid, "payment:reconcile"},
	{Paid, Allocated, "inventory:reserve"},
	{Allocated, Delivered, "order:transition"},
	{Draft, Cancelled, "order:transition"},
	{Placed, Cancelled, "order:transition"},
}

func New(id, tenantID, organizationID, customerPrincipal, currency string, totalMinorUnits int64) (Order, error) {
	if id == "" || tenantID == "" || organizationID == "" || customerPrincipal == "" {
		return Order{}, fmt.Errorf("%w: identifiers are required", ErrConflict)
	}
	if len(currency) != 3 || totalMinorUnits < 0 {
		return Order{}, fmt.Errorf("%w: invalid money", ErrConflict)
	}
	return Order{ID: id, TenantID: tenantID, OrganizationID: organizationID, CustomerPrincipal: customerPrincipal, State: Draft, Currency: currency, TotalMinorUnits: totalMinorUnits, Version: 1}, nil
}

func (o Order) Transition(target State, permissions map[string]struct{}) (Order, error) {
	for _, candidate := range transitions {
		if candidate.from != o.State || candidate.to != target {
			continue
		}
		if _, all := permissions["*"]; !all {
			if _, allowed := permissions[candidate.permission]; !allowed {
				return Order{}, fmt.Errorf("%w: permission %s required", ErrConflict, candidate.permission)
			}
		}
		o.State = target
		o.Version++
		return o, nil
	}
	return Order{}, fmt.Errorf("%w: transition %s -> %s is not configured", ErrConflict, o.State, target)
}
````

### FILE: `internal/order/service.go`

```yaml
block_id: "GO-ENTERPRISE-BACKEND:internal-order-service-go:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "26b9f1c8757a7d390ef9f80b330a0f94b2fcb1ee703b930c80d9281b1c39a9bf"
variables: []
secrets_allowed: false
```

````go
package order

import (
	"context"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("order not found")

type Repository interface {
	Create(context.Context, Order, string, string) (Order, bool, error)
	Get(context.Context, string, string, string) (Order, error)
	SaveTransition(context.Context, Order, int64) error
}

type IDGenerator interface{ New() string }

type Service struct {
	repository Repository
	ids        IDGenerator
}

func NewService(repository Repository, ids IDGenerator) *Service {
	return &Service{repository: repository, ids: ids}
}

type Create struct {
	TenantID          string
	OrganizationID    string
	CustomerPrincipal string
	Currency          string
	TotalMinorUnits   int64
	IdempotencyKey    string
	RequestHash       string
}

func (s *Service) Create(ctx context.Context, command Create) (Order, bool, error) {
	if len(command.IdempotencyKey) < 16 || len(command.IdempotencyKey) > 128 || len(command.RequestHash) != 64 {
		return Order{}, false, fmt.Errorf("%w: invalid idempotency metadata", ErrConflict)
	}
	entity, err := New(s.ids.New(), command.TenantID, command.OrganizationID, command.CustomerPrincipal, command.Currency, command.TotalMinorUnits)
	if err != nil {
		return Order{}, false, err
	}
	return s.repository.Create(ctx, entity, command.IdempotencyKey, command.RequestHash)
}

func (s *Service) Transition(ctx context.Context, tenantID, organizationID, id string, target State, permissions map[string]struct{}) (Order, error) {
	current, err := s.repository.Get(ctx, tenantID, organizationID, id)
	if err != nil {
		return Order{}, err
	}
	next, err := current.Transition(target, permissions)
	if err != nil {
		return Order{}, err
	}
	if err := s.repository.SaveTransition(ctx, next, current.Version); err != nil {
		return Order{}, err
	}
	return next, nil
}
````

### FILE: `internal/platform/httpapi/context.go`

```yaml
block_id: "GO-ENTERPRISE-BACKEND:internal-platform-httpapi-context-go:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "a0d52c102729cfe34e07fddb11b2264cacef54dc6eddd21567f079e7247d2342"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"context"
	"net/http"
	"time"
)

func contextWithTimeout(r *http.Request, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), timeout)
}
````

### FILE: `internal/platform/httpapi/server_test.go`

```yaml
block_id: "GO-ENTERPRISE-BACKEND:internal-platform-httpapi-server_test-go:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "27ffaf9b6fcb4366ae7152a3950888484a2c10e89d93db96674a0d175d3237ce"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"elite.local/enterprise/internal/order"
	"elite.local/enterprise/internal/platform/identity"
)

type fixedIDs struct{}

func (fixedIDs) New() string { return "order-1" }

type repository struct{ created order.Order }

type verifier struct {
	principal identity.Principal
	err       error
}

func (v verifier) Verify(context.Context, string) (identity.Principal, error) {
	return v.principal, v.err
}

func allowedVerifier() verifier {
	return verifier{principal: identity.Principal{
		Subject: "customer-1", TenantID: "018f4d4a-7b36-7a21-8d10-2f4c54c28a01",
		Permissions:   map[string]struct{}{"order:create": {}},
		Organizations: map[string]struct{}{"org-1": {}},
	}}
}

func (r *repository) Create(_ context.Context, entity order.Order, _, _ string) (order.Order, bool, error) {
	r.created = entity
	return entity, false, nil
}
func (r *repository) Get(context.Context, string, string, string) (order.Order, error) {
	return order.Order{}, order.ErrNotFound
}
func (r *repository) SaveTransition(context.Context, order.Order, int64) error { return nil }

func TestCreateOrderContract(t *testing.T) {
	repo := &repository{}
	handler := New(order.NewService(repo, fixedIDs{}), allowedVerifier())
	body := `{"OrganizationID":"org-1","Currency":"USD","TotalMinorUnits":100}`
	request := httptest.NewRequest(http.MethodPost, "/v1/orders", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Idempotency-Key", "order-create-000001")
	request.Header.Set("Authorization", "Bearer verified-by-test-double")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", response.Code, response.Body.String())
	}
	if repo.created.State != order.Draft {
		t.Fatalf("unexpected state %s", repo.created.State)
	}
}

func TestProblemUsesProblemJSON(t *testing.T) {
	handler := New(order.NewService(&repository{}, fixedIDs{}), allowedVerifier())
	request := httptest.NewRequest(http.MethodPost, "/v1/orders", strings.NewReader(`{}`))
	request.Header.Set("Authorization", "Bearer verified-by-test-double")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("expected 415, got %d", response.Code)
	}
	if got := response.Header().Get("Content-Type"); got != "application/problem+json" {
		t.Fatalf("unexpected content type %q", got)
	}
}

func TestCreateOrderRequiresAuthentication(t *testing.T) {
	handler := New(order.NewService(&repository{}, fixedIDs{}), allowedVerifier())
	request := httptest.NewRequest(http.MethodPost, "/v1/orders", strings.NewReader(`{}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", response.Code)
	}
}

func TestCreateOrderRejectsUnauthorizedOrganization(t *testing.T) {
	repo := &repository{}
	handler := New(order.NewService(repo, fixedIDs{}), allowedVerifier())
	request := httptest.NewRequest(http.MethodPost, "/v1/orders", strings.NewReader(`{"OrganizationID":"org-other","Currency":"USD","TotalMinorUnits":100}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Idempotency-Key", "order-create-000002")
	request.Header.Set("Authorization", "Bearer verified-by-test-double")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d: %s", response.Code, response.Body.String())
	}
	if repo.created.ID != "" {
		t.Fatal("repository called for unauthorized organization")
	}
}

func TestHTTPShutdownDrainsActiveRequest(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	started, release := make(chan struct{}), make(chan struct{})
	var releaseOnce sync.Once
	defer releaseOnce.Do(func() { close(release) })
	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-release
		_, _ = io.WriteString(w, "completed")
	})}
	defer server.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	acceptClosed := make(chan struct{})
	go func() {
		done <- serveUntilShutdown(ctx, server, 2*time.Second, func() error { err := server.Serve(listener); close(acceptClosed); return err })
	}()
	type response struct {
		status int
		body   string
		err    error
	}
	received := make(chan response, 1)
	client := &http.Client{Timeout: 3 * time.Second}
	defer client.CloseIdleConnections()
	go func() {
		r, e := client.Get("http://" + listener.Addr().String())
		if e != nil {
			received <- response{err: e}
			return
		}
		defer r.Body.Close()
		body, e := io.ReadAll(r.Body)
		received <- response{r.StatusCode, string(body), e}
	}()
	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("request did not start")
	}
	cancel()
	select {
	case <-acceptClosed:
	case <-time.After(2 * time.Second):
		t.Fatal("accept loop did not stop")
	}
	connection, err := net.DialTimeout("tcp", listener.Addr().String(), 200*time.Millisecond)
	if err == nil {
		connection.Close()
		t.Error("accepted a connection while draining")
	}
	early := false
	var serveErr error
	select {
	case serveErr = <-done:
		early = true
		t.Error("host returned before the active request drained")
	case <-time.After(150 * time.Millisecond):
	}
	releaseOnce.Do(func() { close(release) })
	if !early {
		select {
		case serveErr = <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("drain did not complete")
		}
	}
	if serveErr != nil {
		t.Error(serveErr)
	}
	select {
	case r := <-received:
		if r.err != nil || r.status != 200 || r.body != "completed" {
			t.Fatal("active response lost", r)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("response missing")
	}
	if !t.Failed() {
		t.Log("HTTP_GRACEFUL_DRAIN_PASS accept_closed=1 request_completed=1 host_waited=1")
	}
}

func TestHTTPShutdownDeadlineClosesActiveConnection(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	started, canceled := make(chan struct{}), make(chan struct{})
	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-r.Context().Done()
		close(canceled)
	})}
	defer server.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() {
		done <- serveUntilShutdown(ctx, server, 80*time.Millisecond, func() error { return server.Serve(listener) })
	}()
	client := &http.Client{Timeout: 2 * time.Second}
	defer client.CloseIdleConnections()
	requestDone := make(chan error, 1)
	go func() {
		r, e := client.Get("http://" + listener.Addr().String())
		if r != nil {
			r.Body.Close()
		}
		requestDone <- e
	}()
	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("request did not start")
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("missing drain failure: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("shutdown exceeded bound")
	}
	select {
	case <-canceled:
	case <-time.After(time.Second):
		t.Fatal("active connection was not canceled")
	}
	select {
	case err := <-requestDone:
		if err == nil {
			t.Fatal("truncated request reported success")
		}
	case <-time.After(time.Second):
		t.Fatal("client did not terminate")
	}
	t.Log("HTTP_SHUTDOWN_DEADLINE_PASS failure_reported=1 active_connection_closed=1")
}

func TestHTTPShutdownStartupFailure(t *testing.T) {
	occupied, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer occupied.Close()
	server := &http.Server{Addr: occupied.Addr().String()}
	defer server.Close()
	done := make(chan error, 1)
	go func() { done <- ServeUntilShutdown(context.Background(), server, time.Second) }()
	select {
	case err := <-done:
		var op *net.OpError
		if !errors.As(err, &op) || op.Op != "listen" {
			t.Fatalf("startup cause lost: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("startup failure waits for cancellation")
	}
	t.Log("HTTP_STARTUP_FAILURE_PASS cause_preserved=1 returned_without_signal=1")
}

func TestHTTPShutdownListenerFailureDrainsRequest(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	started, release := make(chan struct{}), make(chan struct{})
	var releaseOnce sync.Once
	defer releaseOnce.Do(func() { close(release) })
	server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-release
		_, _ = io.WriteString(w, "committed")
	})}
	defer server.Close()
	acceptClosed := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		done <- serveUntilShutdown(context.Background(), server, 2*time.Second, func() error {
			err := server.Serve(listener)
			close(acceptClosed)
			return err
		})
	}()
	response := make(chan string, 1)
	client := &http.Client{Timeout: 3 * time.Second}
	defer client.CloseIdleConnections()
	go func() {
		r, err := client.Get("http://" + listener.Addr().String())
		if err != nil {
			response <- "failed"
			return
		}
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			response <- "failed"
			return
		}
		response <- string(body)
	}()
	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("request did not start")
	}
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case <-acceptClosed:
	case <-time.After(time.Second):
		t.Fatal("listener failure not observed")
	}
	select {
	case err := <-done:
		t.Fatalf("returned before request drained: %v", err)
	case <-time.After(150 * time.Millisecond):
	}
	releaseOnce.Do(func() { close(release) })
	select {
	case err := <-done:
		if !errors.Is(err, net.ErrClosed) {
			t.Fatalf("listener cause lost: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("drain did not complete")
	}
	select {
	case body := <-response:
		if body != "committed" {
			t.Fatalf("active response lost: %s", body)
		}
	case <-time.After(time.Second):
		t.Fatal("response did not complete")
	}
	t.Log("HTTP_LISTENER_FAILURE_DRAIN_PASS cause_preserved=1 request_completed=1")
}

func TestHTTPShutdownRejectsInvalidLifecycle(t *testing.T) {
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	for _, tc := range []struct {
		name     string
		ctx      context.Context
		server   *http.Server
		timeout  time.Duration
		nilServe bool
	}{
		{"nil_context", nil, &http.Server{}, time.Second, false},
		{"nil_server", context.Background(), nil, time.Second, false},
		{"zero_timeout", context.Background(), &http.Server{}, 0, false},
		{"negative_timeout", context.Background(), &http.Server{}, -time.Second, false},
		{"nil_serve", context.Background(), &http.Server{}, time.Second, true},
		{"already_canceled", canceled, &http.Server{}, time.Second, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			var serve func() error = func() error { called = true; return nil }
			if tc.nilServe {
				serve = nil
			}
			err := serveUntilShutdown(tc.ctx, tc.server, tc.timeout, serve)
			if err == nil || called {
				t.Fatalf("invalid lifecycle started: called=%v error=%v", called, err)
			}
			if tc.name == "already_canceled" && !errors.Is(err, context.Canceled) {
				t.Fatal("cancellation cause lost")
			}
		})
	}
	if err := ServeUntilShutdown(context.Background(), nil, time.Second); err == nil {
		t.Fatal("nil public server accepted")
	}
}
````

### FILE: `internal/platform/httpapi/server.go`

```yaml
block_id: "GO-ENTERPRISE-BACKEND:internal-platform-httpapi-server-go:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "66f182b2e02a6ef076f27e39a06d4cf540f22eb230920eafcc5a35b572fdb6e1"
variables: []
secrets_allowed: false
```

````go
package httpapi

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"elite.local/enterprise/internal/order"
	"elite.local/enterprise/internal/platform/identity"
)

type Server struct {
	orders   *order.Service
	verifier identity.Verifier
}

func New(orders *order.Service, verifier identity.Verifier) http.Handler {
	s := &Server{orders: orders, verifier: verifier}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health/live", s.live)
	mux.HandleFunc("POST /v1/orders", s.createOrder)
	return recoverMiddleware(mux)
}

func (s *Server) live(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) createOrder(w http.ResponseWriter, r *http.Request) {
	principal, err := authenticate(r.Context(), r.Header.Get("Authorization"), s.verifier)
	if err != nil {
		writeProblem(w, 401, "UNAUTHENTICATED", "a valid bearer token is required")
		return
	}
	if !principal.Allowed("order:create") {
		writeProblem(w, 403, "FORBIDDEN", "order:create permission is required")
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		writeProblem(w, 415, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return
	}
	key := r.Header.Get("Idempotency-Key")
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
	if err != nil {
		writeProblem(w, 400, "INVALID_BODY", "body is invalid or too large")
		return
	}
	var input struct {
		OrganizationID, Currency string
		TotalMinorUnits          int64
	}
	decoder := json.NewDecoder(strings.NewReader(string(body)))
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&input); err != nil {
		writeProblem(w, 400, "INVALID_BODY", "body does not match the contract")
		return
	}
	if !principal.AllowedOrganization(input.OrganizationID) {
		writeProblem(w, 403, "ORGANIZATION_FORBIDDEN", "token is not authorized for this organization")
		return
	}
	hash := sha256.Sum256(body)
	entity, replayed, err := s.orders.Create(r.Context(), order.Create{TenantID: principal.TenantID, OrganizationID: input.OrganizationID, CustomerPrincipal: principal.Subject, Currency: input.Currency, TotalMinorUnits: input.TotalMinorUnits, IdempotencyKey: key, RequestHash: hex.EncodeToString(hash[:])})
	if err != nil {
		if errors.Is(err, order.ErrConflict) {
			writeProblem(w, 409, "ORDER_CONFLICT", err.Error())
			return
		}
		writeProblem(w, 500, "INTERNAL_ERROR", "request failed")
		return
	}
	if replayed {
		writeJSON(w, 200, entity)
	} else {
		writeJSON(w, 201, entity)
	}
}

func authenticate(ctx context.Context, authorization string, verifier identity.Verifier) (identity.Principal, error) {
	const prefix = "Bearer "
	if !strings.HasPrefix(authorization, prefix) || len(authorization) == len(prefix) {
		return identity.Principal{}, identity.ErrUnauthenticated
	}
	return verifier.Verify(ctx, strings.TrimSpace(strings.TrimPrefix(authorization, prefix)))
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeProblem(w http.ResponseWriter, status int, code, detail string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"type": "urn:elite:problem:" + code, "status": status, "code": code, "detail": detail})
}
func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				writeProblem(w, 500, "INTERNAL_ERROR", "request failed")
			}
		}()
		ctx, cancel := contextWithTimeout(r, 10*time.Second)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ServeUntilShutdown stops admission on cancellation and waits for active HTTP
// requests to drain before the caller can close shared resources. A drain timeout
// force-closes connections and is returned to the caller. Hijacked connections
// require a separate owner; net/http Shutdown and Close do not wait for them.
func ServeUntilShutdown(ctx context.Context, server *http.Server, timeout time.Duration) error {
	if server == nil {
		return errors.New("HTTP server is required")
	}
	return serveUntilShutdown(ctx, server, timeout, server.ListenAndServe)
}
func serveUntilShutdown(ctx context.Context, server *http.Server, timeout time.Duration, serve func() error) error {
	if ctx == nil || server == nil || timeout <= 0 || serve == nil {
		return errors.New("HTTP shutdown requires context, server, positive timeout and serve callback")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	done := make(chan error, 1)
	go func() { done <- serve() }()
	var serveErr error
	serveFinished := false
	select {
	case <-ctx.Done():
	case serveErr = <-done:
		serveFinished = true
	}
	// ListenAndServe/Serve returns as soon as listeners close, before requests
	// drain. Shutdown must therefore be awaited on this caller's path, including
	// when an unexpected listener failure initiated the shutdown.
	shutdown, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	shutdownErr := server.Shutdown(shutdown)
	var closeErr error
	if shutdownErr != nil {
		closeErr = server.Close()
	}
	if !serveFinished {
		serveErr = <-done
	}
	if ctx.Err() != nil && errors.Is(serveErr, http.ErrServerClosed) {
		serveErr = nil
	}
	return errors.Join(serveErr, shutdownErr, closeErr)
}
````

### FILE: `internal/platform/identity/oidc.go`

```yaml
block_id: "GO-ENTERPRISE-BACKEND:internal-platform-identity-oidc-go:v1"
operation: CREATE
provenance: AUTHORED
source: "local composition using coreos/go-oidc public API"
license: "LicenseRef-Workspace-Owner"
sha256: "cfc53d009d4395644f65f36c23f1b7b1a0109f535f36b739b7278f661b8f90fa"
variables: []
secrets_allowed: false
```

````go
package identity

import (
	"context"
	"errors"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
)

var ErrUnauthenticated = errors.New("unauthenticated")

type Principal struct {
	Subject       string
	TenantID      string
	Permissions   map[string]struct{}
	Organizations map[string]struct{}
}

func (p Principal) Allowed(permission string) bool {
	if _, all := p.Permissions["*"]; all {
		return true
	}
	_, allowed := p.Permissions[permission]
	return allowed
}

func (p Principal) AllowedOrganization(organizationID string) bool {
	if organizationID == "" {
		return false
	}
	if _, all := p.Permissions["*"]; all {
		return true
	}
	_, allowed := p.Organizations[organizationID]
	return allowed
}

type Verifier interface {
	Verify(context.Context, string) (Principal, error)
}

type OIDCVerifier struct{ verifier *oidc.IDTokenVerifier }

func NewOIDCVerifier(ctx context.Context, issuer, audience string) (*OIDCVerifier, error) {
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, err
	}
	return &OIDCVerifier{verifier: provider.Verifier(&oidc.Config{
		ClientID:             audience,
		SupportedSigningAlgs: []string{oidc.RS256},
	})}, nil
}

func (v *OIDCVerifier) Verify(ctx context.Context, raw string) (Principal, error) {
	token, err := v.verifier.Verify(ctx, raw)
	if err != nil {
		return Principal{}, fmt.Errorf("%w: token rejected", ErrUnauthenticated)
	}
	var claims struct {
		TenantID        string   `json:"tenant_id"`
		Permissions     []string `json:"permissions"`
		OrganizationIDs []string `json:"organization_ids"`
	}
	if err := token.Claims(&claims); err != nil || token.Subject == "" || claims.TenantID == "" {
		return Principal{}, fmt.Errorf("%w: required claims missing", ErrUnauthenticated)
	}
	permissions := make(map[string]struct{}, len(claims.Permissions))
	for _, permission := range claims.Permissions {
		if permission != "" {
			permissions[permission] = struct{}{}
		}
	}
	organizations := make(map[string]struct{}, len(claims.OrganizationIDs))
	for _, organizationID := range claims.OrganizationIDs {
		if organizationID != "" {
			organizations[organizationID] = struct{}{}
		}
	}
	return Principal{Subject: token.Subject, TenantID: claims.TenantID, Permissions: permissions, Organizations: organizations}, nil
}
````

### FILE: `internal/platform/postgres/outbox.go`

```yaml
block_id: "GO-ENTERPRISE-BACKEND:internal-platform-postgres-outbox-go:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "ac68c30fdc76b1b9d9fb2a3c6014a9d805604ffa3bf6625ee8b939b399d72328"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboxEvent struct {
	TenantID         string
	EventID          string
	AggregateType    string
	AggregateID      string
	AggregateVersion int64
	EventType        string
	SchemaVersion    int
	OccurredAt       time.Time
	Payload          json.RawMessage
	Headers          json.RawMessage
	// Attempts is the monotonically increasing claim generation. Never reset it.
	Attempts int
}

var ErrOutboxClaimLost = errors.New("outbox event claim lost")

type Outbox struct{ pool *pgxpool.Pool }

func NewOutbox(pool *pgxpool.Pool) *Outbox { return &Outbox{pool: pool} }

func (o *Outbox) Claim(ctx context.Context, workerID string, lease time.Duration, batchSize int) ([]OutboxEvent, error) {
	if o == nil || o.pool == nil || workerID == "" || lease < time.Microsecond || batchSize < 1 || batchSize > 1000 {
		return nil, fmt.Errorf("invalid outbox claim parameters")
	}
	rows, err := o.pool.Query(ctx, `
with candidates as (
  select tenant_id,event_id
    from platform.outbox_event
   where published_at is null
     and available_at <= clock_timestamp()
     and (claimed_until is null or claimed_until < clock_timestamp())
   order by available_at,occurred_at,event_id
   for update skip locked
   limit $1
)
update platform.outbox_event as event
   set claimed_by=$2,
       claimed_until=clock_timestamp()+$3::interval,
       attempts=attempts+1
  from candidates
 where event.tenant_id=candidates.tenant_id
   and event.event_id=candidates.event_id
returning event.tenant_id,event.event_id,event.aggregate_type,event.aggregate_id,
          event.aggregate_version,event.event_type,event.schema_version,
          event.occurred_at,event.payload,event.headers,event.attempts`, batchSize, workerID, lease.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	events := make([]OutboxEvent, 0, batchSize)
	for rows.Next() {
		var event OutboxEvent
		if err := rows.Scan(&event.TenantID, &event.EventID, &event.AggregateType, &event.AggregateID, &event.AggregateVersion, &event.EventType, &event.SchemaVersion, &event.OccurredAt, &event.Payload, &event.Headers, &event.Attempts); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, rows.Err()
}

// Lock before checking the DB clock: a wait on this row may outlive the lease.
// Worker identity alone cannot distinguish a restarted or concurrent attempt.
const lockedOutbox = `with owned as materialized (
 select tenant_id,event_id,claimed_by,claimed_until,attempts
 from platform.outbox_event where tenant_id=$1 and event_id=$2
 and published_at is null for update
) `

func (o *Outbox) MarkPublished(ctx context.Context, tenantID, eventID, workerID string, attempt int) error {
	if o == nil || o.pool == nil || workerID == "" || attempt < 1 {
		return ErrOutboxClaimLost
	}
	result, err := o.pool.Exec(ctx, lockedOutbox+`update platform.outbox_event e set published_at=clock_timestamp(),claimed_by=null,claimed_until=null,last_error_code=null
 from owned o where e.tenant_id=o.tenant_id and e.event_id=o.event_id
 and o.claimed_by=$3 and o.attempts=$4 and o.claimed_until>clock_timestamp()`, tenantID, eventID, workerID, attempt)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return ErrOutboxClaimLost
	}
	return nil
}

func (o *Outbox) Release(ctx context.Context, tenantID, eventID, workerID, errorCode string, attempt int, retryAfter time.Duration) error {
	if o == nil || o.pool == nil || workerID == "" || attempt < 1 {
		return ErrOutboxClaimLost
	}
	if retryAfter < 0 || errorCode == "" {
		return fmt.Errorf("invalid outbox release parameters")
	}
	result, err := o.pool.Exec(ctx, lockedOutbox+`update platform.outbox_event e set claimed_by=null,claimed_until=null,last_error_code=$4,available_at=clock_timestamp()+$5::interval
 from owned o where e.tenant_id=o.tenant_id and e.event_id=o.event_id
 and o.claimed_by=$3 and o.attempts=$6 and o.claimed_until>clock_timestamp()`, tenantID, eventID, workerID, errorCode, retryAfter.String(), attempt)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return ErrOutboxClaimLost
	}
	return nil
}

// RemainingLease is a conservative preflight, not a remote-effect fence. The
// publisher still needs stable idempotency and reconciliation after ambiguity.
func (o *Outbox) RemainingLease(ctx context.Context, tenantID, eventID, workerID string, attempt int) (time.Duration, error) {
	if o == nil || o.pool == nil || workerID == "" || attempt < 1 {
		return 0, ErrOutboxClaimLost
	}
	var micros int64
	err := o.pool.QueryRow(ctx, `select floor(extract(epoch from (claimed_until-clock_timestamp()))*1000000)::bigint
 from platform.outbox_event where tenant_id=$1 and event_id=$2 and claimed_by=$3 and attempts=$4
 and published_at is null and claimed_until>clock_timestamp()`, tenantID, eventID, workerID, attempt).Scan(&micros)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrOutboxClaimLost
	}
	if err != nil {
		return 0, err
	}
	if micros <= 0 || micros > int64((time.Duration(1<<63-1))/time.Microsecond) {
		return 0, ErrOutboxClaimLost
	}
	return time.Duration(micros) * time.Microsecond, nil
}
````

### FILE: `internal/platform/postgres/outbox_integration_test.go`

```yaml
block_id: "GO-ENTERPRISE-BACKEND:internal-platform-postgres-outbox-integration-test-go:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "f4f8ea1921f2981e5ea93a9e0b22b6dc504978e7abd996808a37647c97e129aa"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestOutboxClaimsAreDisjointAndOwned(t *testing.T) {
	ctx, pool := outboxTestPool(t)
	const tenantID = "018f4d4a-7b36-7a21-8d10-2f4c54c28b02"
	cleanup := func() {
		_, _ = pool.Exec(ctx, `delete from platform.outbox_event where tenant_id=$1`, tenantID)
		_, _ = pool.Exec(ctx, `delete from platform.tenant where tenant_id=$1`, tenantID)
	}
	cleanup()
	defer cleanup()
	if _, err := pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'outbox-test','Outbox Test S.A.','Outbox Test')`, tenantID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values ($1,'018f4d4a-7b36-7a21-8d10-2f4c54c28b10','test','a',1,'test.created',1,clock_timestamp(),'{}'),($1,'018f4d4a-7b36-7a21-8d10-2f4c54c28b11','test','b',1,'test.created',1,clock_timestamp(),'{}')`, tenantID); err != nil {
		t.Fatal(err)
	}
	store := NewOutbox(pool)
	first, err := store.Claim(ctx, "worker-a", time.Minute, 1)
	if err != nil || len(first) != 1 {
		t.Fatalf("first claim failed: events=%+v err=%v", first, err)
	}
	second, err := store.Claim(ctx, "worker-b", time.Minute, 1)
	if err != nil || len(second) != 1 || second[0].EventID == first[0].EventID {
		t.Fatalf("second claim was not disjoint: first=%+v second=%+v err=%v", first, second, err)
	}
	if err := store.MarkPublished(ctx, first[0].TenantID, first[0].EventID, "worker-b", first[0].Attempts); err == nil {
		t.Fatal("wrong worker marked an event")
	}
	if err := store.MarkPublished(ctx, first[0].TenantID, first[0].EventID, "worker-a", first[0].Attempts); err != nil {
		t.Fatal(err)
	}
	if err := store.Release(ctx, second[0].TenantID, second[0].EventID, "worker-b", "PROVIDER_TIMEOUT", second[0].Attempts, 0); err != nil {
		t.Fatal(err)
	}
	reclaimed, err := store.Claim(ctx, "worker-c", time.Minute, 1)
	if err != nil || len(reclaimed) != 1 || reclaimed[0].EventID != second[0].EventID {
		t.Fatalf("released event was not reclaimed: events=%+v err=%v", reclaimed, err)
	}
}

func outboxTestConfig(raw string) (*pgxpool.Config, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, errors.New("explicit disposable database required")
	}
	cfg, err := pgxpool.ParseConfig(raw)
	if err != nil || cfg.ConnConfig.Host != "127.0.0.1" || !strings.HasPrefix(cfg.ConnConfig.Database, "elite_outbox_test_") || len(strings.TrimPrefix(cfg.ConnConfig.Database, "elite_outbox_test_")) < 16 {
		return nil, errors.New("requires dedicated loopback elite_outbox_test_<unique> database")
	}
	for _, fallback := range cfg.ConnConfig.Fallbacks {
		if fallback.Host != "127.0.0.1" {
			return nil, errors.New("non-loopback fallback forbidden")
		}
	}
	cfg.MaxConns = 6
	cfg.ConnConfig.ConnectTimeout = 5 * time.Second
	return cfg, nil
}

func outboxTestPool(t *testing.T) (context.Context, *pgxpool.Pool) {
	t.Helper()
	raw := os.Getenv("OUTBOX_TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("OUTBOX_TEST_DATABASE_URL is not set; isolated database required")
	}
	cfg, err := outboxTestConfig(raw)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	if err = pool.Ping(ctx); err != nil {
		t.Fatal(err)
	}
	return ctx, pool
}

func outboxLeaseFixture(t *testing.T) (context.Context, *pgxpool.Pool, *Outbox) {
	t.Helper()
	ctx, pool := outboxTestPool(t)
	const tenant = "018f4d4a-7b36-7a21-8d10-2f4c54c28504"
	cleanup := func() {
		clean, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := pool.Exec(clean, `delete from platform.outbox_event where tenant_id=$1`, tenant); err != nil {
			t.Error(err)
		}
		if _, err := pool.Exec(clean, `delete from platform.tenant where tenant_id=$1`, tenant); err != nil {
			t.Error(err)
		}
	}
	cleanup()
	t.Cleanup(cleanup)
	if _, err := pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'outbox-v316','Synthetic','Synthetic')`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,'018f4d4a-7b36-7a21-8d10-2f4c54c28505','test','v316',1,'test.created',1,clock_timestamp(),'{}')`, tenant); err != nil {
		t.Fatal(err)
	}
	return ctx, pool, NewOutbox(pool)
}

func finishOutboxClaim(ctx context.Context, store *Outbox, event OutboxEvent, release bool) error {
	if release {
		return store.Release(ctx, event.TenantID, event.EventID, "same-worker", "PUBLISH_FAILED", event.Attempts, 0)
	}
	return store.MarkPublished(ctx, event.TenantID, event.EventID, "same-worker", event.Attempts)
}

func TestOutboxExpiredClaimCannotFinalize(t *testing.T) {
	for _, release := range []bool{false, true} {
		t.Run(map[bool]string{false: "publish", true: "release"}[release], func(t *testing.T) {
			ctx, pool, store := outboxLeaseFixture(t)
			events, err := store.Claim(ctx, "same-worker", time.Minute, 1)
			if err != nil || len(events) != 1 {
				t.Fatalf("claim: %v %d", err, len(events))
			}
			e := events[0]
			if _, err = pool.Exec(ctx, `update platform.outbox_event set claimed_until=clock_timestamp()-interval '1 second' where tenant_id=$1 and event_id=$2`, e.TenantID, e.EventID); err != nil {
				t.Fatal(err)
			}
			if err = finishOutboxClaim(ctx, store, e, release); err == nil {
				t.Fatal("expired claim finalized")
			}
		})
	}
}

func TestOutboxReclaimedGenerationCannotFinalize(t *testing.T) {
	for _, release := range []bool{false, true} {
		t.Run(map[bool]string{false: "publish", true: "release"}[release], func(t *testing.T) {
			ctx, pool, store := outboxLeaseFixture(t)
			first, err := store.Claim(ctx, "same-worker", time.Minute, 1)
			if err != nil || len(first) != 1 {
				t.Fatal("first claim", err)
			}
			e := first[0]
			if _, err = pool.Exec(ctx, `update platform.outbox_event set claimed_until=clock_timestamp()-interval '1 second' where tenant_id=$1 and event_id=$2`, e.TenantID, e.EventID); err != nil {
				t.Fatal(err)
			}
			second, err := store.Claim(ctx, "same-worker", time.Minute, 1)
			if err != nil || len(second) != 1 || second[0].EventID != e.EventID || second[0].Attempts != e.Attempts+1 {
				t.Fatal("reclaim", err)
			}
			if err = finishOutboxClaim(ctx, store, e, release); err == nil {
				t.Fatal("stale generation finalized reassigned work")
			}
			if err = finishOutboxClaim(ctx, store, second[0], release); err != nil {
				t.Fatal("current generation failed", err)
			}
		})
	}
}

func TestOutboxFinalizationChecksClockAfterLock(t *testing.T) {
	for _, release := range []bool{false, true} {
		t.Run(map[bool]string{false: "publish", true: "release"}[release], func(t *testing.T) {
			ctx, pool, store := outboxLeaseFixture(t)
			events, err := store.Claim(ctx, "same-worker", 500*time.Millisecond, 1)
			if err != nil || len(events) != 1 {
				t.Fatal("claim", err)
			}
			e := events[0]
			tx, err := pool.Begin(ctx)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback(context.Background())
			if _, err = tx.Exec(ctx, `select event_id from platform.outbox_event where tenant_id=$1 and event_id=$2 for update`, e.TenantID, e.EventID); err != nil {
				t.Fatal(err)
			}
			done := make(chan error, 1)
			go func() { done <- finishOutboxClaim(ctx, store, e, release) }()
			deadline := time.Now().Add(2 * time.Second)
			blocked := false
			for time.Now().Before(deadline) {
				// Observe from fresh transactions: pg_stat_activity can retain a
				// statistics snapshot for the lifetime of the lock-holding tx.
				if err = pool.QueryRow(ctx, `select exists(select 1 from pg_stat_activity where datname=current_database() and wait_event_type='Lock' and pid<>pg_backend_pid())`).Scan(&blocked); err != nil {
					t.Fatal(err)
				}
				if blocked {
					break
				}
				time.Sleep(5 * time.Millisecond)
			}
			if !blocked {
				t.Fatal("finalization did not wait on row lock")
			}
			expired := false
			for time.Now().Before(deadline) {
				if err = tx.QueryRow(ctx, `select claimed_until<=clock_timestamp() from platform.outbox_event where tenant_id=$1 and event_id=$2`, e.TenantID, e.EventID).Scan(&expired); err != nil {
					t.Fatal(err)
				}
				if expired {
					break
				}
				time.Sleep(10 * time.Millisecond)
			}
			if !expired {
				t.Fatal("lease did not expire while blocked")
			}
			if err = tx.Commit(ctx); err != nil {
				t.Fatal(err)
			}
			select {
			case err = <-done:
				if err == nil {
					t.Fatal("claim finalized after lock wait exceeded lease")
				}
			case <-time.After(2 * time.Second):
				t.Fatal("finalization did not return")
			}
		})
	}
}

func TestOutboxDatabaseGuard(t *testing.T) {
	for _, raw := range []string{"", "postgres://postgres@127.0.0.1/main?sslmode=disable", "postgres://postgres@localhost/elite_outbox_test_0123456789abcdef?sslmode=disable", "postgres://postgres@192.0.2.1/elite_outbox_test_0123456789abcdef?sslmode=disable", "postgres://postgres@127.0.0.1/elite_outbox_test_short?sslmode=disable", "host=127.0.0.1,192.0.2.1 dbname=elite_outbox_test_0123456789abcdef sslmode=disable"} {
		if _, err := outboxTestConfig(raw); err == nil {
			t.Fatal("unsafe fixture target accepted")
		}
	}
	if _, err := outboxTestConfig("postgres://postgres@127.0.0.1/elite_outbox_test_0123456789abcdef?sslmode=disable"); err != nil {
		t.Fatal(err)
	}
}

func TestOutboxRemainingLeaseRequiresCurrentOwner(t *testing.T) {
	ctx, pool, store := outboxLeaseFixture(t)
	events, err := store.Claim(ctx, "same-worker", time.Minute, 1)
	if err != nil || len(events) != 1 {
		t.Fatal("claim", err)
	}
	e := events[0]
	if remaining, err := store.RemainingLease(ctx, e.TenantID, e.EventID, "same-worker", e.Attempts); err != nil || remaining <= 0 || remaining > time.Minute {
		t.Fatal("invalid remaining lease", remaining, err)
	}
	for _, tc := range []struct {
		worker  string
		attempt int
	}{{"other-worker", e.Attempts}, {"same-worker", 0}, {"same-worker", e.Attempts + 1}} {
		if _, err = store.RemainingLease(ctx, e.TenantID, e.EventID, tc.worker, tc.attempt); !errors.Is(err, ErrOutboxClaimLost) {
			t.Fatal("invalid owner/generation retained lease", err)
		}
	}
	if _, err = pool.Exec(ctx, `update platform.outbox_event set claimed_until=clock_timestamp()-interval '1 second' where tenant_id=$1 and event_id=$2`, e.TenantID, e.EventID); err != nil {
		t.Fatal(err)
	}
	if _, err = store.RemainingLease(ctx, e.TenantID, e.EventID, "same-worker", e.Attempts); !errors.Is(err, ErrOutboxClaimLost) {
		t.Fatal("expired lease retained budget", err)
	}
	second, err := store.Claim(ctx, "same-worker", time.Minute, 1)
	if err != nil || len(second) != 1 {
		t.Fatal("reclaim", err)
	}
	if _, err = store.RemainingLease(ctx, e.TenantID, e.EventID, "same-worker", e.Attempts); !errors.Is(err, ErrOutboxClaimLost) {
		t.Fatal("stale generation retained budget", err)
	}
	if err = finishOutboxClaim(ctx, store, second[0], false); err != nil {
		t.Fatal(err)
	}
	if _, err = store.RemainingLease(ctx, e.TenantID, e.EventID, "same-worker", second[0].Attempts); !errors.Is(err, ErrOutboxClaimLost) {
		t.Fatal("published event retained budget", err)
	}
}
````

### FILE: `internal/platform/postgres/orders_integration_test.go`

```yaml
block_id: "GO-ENTERPRISE-BACKEND:internal-platform-postgres-orders-integration-test-go:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "abb00147702e3dd2150a356a5c1b6a5ef5d3c411de8d7d509df9a31a4ad383d2"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"errors"
	"os"
	"testing"

	"elite.local/enterprise/internal/order"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestOrdersRepositoryIntegration(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	const tenantID = "018f4d4a-7b36-7a21-8d10-2f4c54c28b01"
	const organizationID = "integration-org"
	const orderID = "integration-order"
	cleanup := func() {
		_, _ = pool.Exec(ctx, `delete from platform.outbox_event where tenant_id=$1`, tenantID)
		_, _ = pool.Exec(ctx, `delete from platform.idempotency_record where tenant_id=$1`, tenantID)
		_, _ = pool.Exec(ctx, `delete from sales.customer_order where tenant_id=$1`, tenantID)
		_, _ = pool.Exec(ctx, `delete from org.organization where tenant_id=$1`, tenantID)
		_, _ = pool.Exec(ctx, `delete from platform.tenant where tenant_id=$1`, tenantID)
	}
	cleanup()
	defer cleanup()
	if _, err := pool.Exec(ctx, `insert into platform.tenant(tenant_id,tenant_code,legal_name,display_name) values($1,'integration-tenant','Integration Tenant S.A.','Integration Tenant')`, tenantID); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `insert into org.organization(tenant_id,organization_id,organization_code,display_name,organization_type) values($1,$2,'integration-org','Integration Org','enterprise')`, tenantID, organizationID); err != nil {
		t.Fatal(err)
	}

	repository := NewOrders(pool)
	created := order.Order{ID: orderID, TenantID: tenantID, OrganizationID: organizationID, CustomerPrincipal: "customer-1", State: order.Draft, Currency: "USD", TotalMinorUnits: 100, Version: 1}
	actual, replayed, err := repository.Create(ctx, created, "integration-idem-0001", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	if err != nil || replayed || actual.ID != orderID {
		t.Fatalf("create failed: order=%+v replayed=%v err=%v", actual, replayed, err)
	}
	replay, replayed, err := repository.Create(ctx, created, "integration-idem-0001", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	if err != nil || !replayed || replay.ID != orderID {
		t.Fatalf("replay failed: order=%+v replayed=%v err=%v", replay, replayed, err)
	}
	if _, _, err := repository.Create(ctx, created, "integration-idem-0001", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"); !errors.Is(err, order.ErrConflict) {
		t.Fatalf("expected hash conflict, got %v", err)
	}
	loaded, err := repository.Get(ctx, tenantID, organizationID, orderID)
	if err != nil || loaded.TenantID != tenantID {
		t.Fatalf("tenant-scoped get failed: order=%+v err=%v", loaded, err)
	}
	next, err := loaded.Transition(order.Placed, map[string]struct{}{"order:create": {}})
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.SaveTransition(ctx, next, loaded.Version); err != nil {
		t.Fatal(err)
	}
	loaded, err = repository.Get(ctx, tenantID, organizationID, orderID)
	if err != nil || loaded.State != order.Placed || loaded.Version != 2 {
		t.Fatalf("transition persistence failed: order=%+v err=%v", loaded, err)
	}
}
````

### FILE: `internal/platform/postgres/orders.go`

```yaml
block_id: "GO-ENTERPRISE-BACKEND:internal-platform-postgres-orders-go:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "41845e84ee338adb9ab871b25a53511d293aeff028c74ce918ee7fa2b4f01a24"
variables: []
secrets_allowed: false
```

````go
package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"elite.local/enterprise/internal/order"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Orders struct{ pool *pgxpool.Pool }

func NewOrders(pool *pgxpool.Pool) *Orders { return &Orders{pool: pool} }

func (r *Orders) Create(ctx context.Context, entity order.Order, key, requestHash string) (order.Order, bool, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return order.Order{}, false, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	claim, err := tx.Exec(ctx, `insert into platform.idempotency_record(tenant_id,scope,idempotency_key,request_sha256_hex,status,locked_until,expires_at) values($1,'order:create',$2,$3,'processing',clock_timestamp()+interval '30 seconds',clock_timestamp()+interval '24 hours') on conflict (tenant_id,scope,idempotency_key) do nothing`, entity.TenantID, key, requestHash)
	if err != nil {
		return order.Order{}, false, err
	}
	if claim.RowsAffected() == 0 {
		var existingHash, status string
		var response []byte
		err = tx.QueryRow(ctx, `select request_sha256_hex,status,response_body from platform.idempotency_record where tenant_id=$1 and scope='order:create' and idempotency_key=$2 for update`, entity.TenantID, key).Scan(&existingHash, &status, &response)
		if err != nil {
			return order.Order{}, false, err
		}
		if existingHash != requestHash {
			return order.Order{}, false, fmt.Errorf("%w: idempotency key reused", order.ErrConflict)
		}
		if status != "completed" {
			return order.Order{}, false, fmt.Errorf("%w: request in progress", order.ErrConflict)
		}
		var replay order.Order
		if err := json.Unmarshal(response, &replay); err != nil {
			return order.Order{}, false, err
		}
		return replay, true, tx.Commit(ctx)
	}
	_, err = tx.Exec(ctx, `insert into sales.customer_order(tenant_id,order_id,organization_id,customer_principal_id,state,currency,total_minor_units,version) values($1,$2,$3,$4,$5,$6,$7,$8)`, entity.TenantID, entity.ID, entity.OrganizationID, entity.CustomerPrincipal, entity.State, entity.Currency, entity.TotalMinorUnits, entity.Version)
	if err != nil {
		return order.Order{}, false, err
	}
	payload, _ := json.Marshal(entity)
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,gen_random_uuid(),'order',$2,$3,'order.created',1,clock_timestamp(),$4)`, entity.TenantID, entity.ID, entity.Version, payload)
	if err != nil {
		return order.Order{}, false, err
	}
	_, err = tx.Exec(ctx, `update platform.idempotency_record set status='completed',response_code=201,response_body=$3,resource_type='order',resource_id=$4,locked_until=null where tenant_id=$1 and scope='order:create' and idempotency_key=$2`, entity.TenantID, key, payload, entity.ID)
	if err != nil {
		return order.Order{}, false, err
	}
	return entity, false, tx.Commit(ctx)
}

func (r *Orders) Get(ctx context.Context, tenantID, organizationID, id string) (order.Order, error) {
	var entity order.Order
	err := r.pool.QueryRow(ctx, `select order_id,tenant_id,organization_id,customer_principal_id,state,currency,total_minor_units,version from sales.customer_order where tenant_id=$1 and organization_id=$2 and order_id=$3`, tenantID, organizationID, id).Scan(&entity.ID, &entity.TenantID, &entity.OrganizationID, &entity.CustomerPrincipal, &entity.State, &entity.Currency, &entity.TotalMinorUnits, &entity.Version)
	if errors.Is(err, pgx.ErrNoRows) {
		return order.Order{}, order.ErrNotFound
	}
	return entity, err
}

func (r *Orders) SaveTransition(ctx context.Context, entity order.Order, expectedVersion int64) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err := tx.Exec(ctx, `update sales.customer_order set state=$4,version=$5,updated_at=clock_timestamp() where tenant_id=$1 and organization_id=$2 and order_id=$3 and version=$6`, entity.TenantID, entity.OrganizationID, entity.ID, entity.State, entity.Version, expectedVersion)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return fmt.Errorf("%w: concurrent update", order.ErrConflict)
	}
	payload, _ := json.Marshal(entity)
	_, err = tx.Exec(ctx, `insert into platform.outbox_event(tenant_id,event_id,aggregate_type,aggregate_id,aggregate_version,event_type,schema_version,occurred_at,payload) values($1,gen_random_uuid(),'order',$2,$3,$4,1,clock_timestamp(),$5)`, entity.TenantID, entity.ID, entity.Version, "order."+string(entity.State), payload)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
````

### FILE: `db/migrations/0002_enterprise_order_core.up.sql`

```yaml
block_id: "GO-ENTERPRISE-BACKEND:enterprise-order-migration-up:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "b37ddc36366b93eda630e466f7398645bde118160e4fd2dd5760b0d3ddf5664d"
variables: []
secrets_allowed: false
```

````sql
begin;

create schema org;
create schema sales;

create table org.organization (
  tenant_id uuid not null,
  organization_id text not null,
  parent_organization_id text,
  organization_code text not null,
  display_name text not null,
  organization_type text not null
    check (organization_type in ('enterprise', 'franchisor', 'franchisee', 'factory', 'warehouse', 'store', 'service_center')),
  status text not null default 'active'
    check (status in ('provisioning', 'active', 'suspended', 'closed')),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, organization_id),
  foreign key (tenant_id) references platform.tenant (tenant_id),
  foreign key (tenant_id, parent_organization_id)
    references org.organization (tenant_id, organization_id),
  unique (tenant_id, organization_code),
  check (length(organization_id) between 1 and 128),
  check (organization_code ~ '^[a-z][a-z0-9]*(-[a-z0-9]+)*$'),
  check (parent_organization_id is null or parent_organization_id <> organization_id)
);

create table sales.customer_order (
  tenant_id uuid not null,
  order_id text not null,
  organization_id text not null,
  customer_principal_id text not null,
  state text not null
    check (state in ('draft', 'placed', 'confirmed', 'paid', 'allocated', 'delivered', 'cancelled')),
  currency text not null check (currency ~ '^[A-Z]{3}$'),
  total_minor_units bigint not null check (total_minor_units >= 0),
  version bigint not null check (version > 0),
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, order_id),
  foreign key (tenant_id, organization_id)
    references org.organization (tenant_id, organization_id),
  check (length(order_id) between 1 and 128),
  check (length(customer_principal_id) between 1 and 256),
  check (updated_at >= created_at)
);

create index customer_order_org_state_idx
  on sales.customer_order (tenant_id, organization_id, state, created_at desc, order_id);

create index customer_order_customer_idx
  on sales.customer_order (tenant_id, customer_principal_id, created_at desc, order_id);

commit;
````

### FILE: `db/migrations/0002_enterprise_order_core.down.sql`

```yaml
block_id: "GO-ENTERPRISE-BACKEND:enterprise-order-migration-down:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "5a8160b8cdc05a5b41595463f5bd5d335b0ea76953c62750c8d651b5429d87e3"
variables: []
secrets_allowed: false
```

````sql
begin;

drop schema sales cascade;
drop schema org cascade;

commit;
````

### FILE: `db/tests/0002_enterprise_order_core.test.sql`

```yaml
block_id: "GO-ENTERPRISE-BACKEND:enterprise-order-migration-test:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "0aa043944b41371cd94080e4d8b6353e999e218004b4d87b2384eb7cd0d36f80"
variables: []
secrets_allowed: false
```

````sql
begin;

insert into platform.tenant (tenant_id, tenant_code, legal_name, display_name)
values
  ('018f4d4a-7b36-7a21-8d10-2f4c54c28a01', 'tenant-a', 'Tenant A S.A.', 'Tenant A'),
  ('018f4d4a-7b36-7a21-8d10-2f4c54c28a02', 'tenant-b', 'Tenant B S.A.', 'Tenant B');

insert into org.organization (
  tenant_id, organization_id, organization_code, display_name, organization_type
) values
  ('018f4d4a-7b36-7a21-8d10-2f4c54c28a01', 'org-main', 'main', 'Main A', 'enterprise'),
  ('018f4d4a-7b36-7a21-8d10-2f4c54c28a02', 'org-main', 'main', 'Main B', 'enterprise');

insert into sales.customer_order (
  tenant_id, order_id, organization_id, customer_principal_id,
  state, currency, total_minor_units, version
) values
  ('018f4d4a-7b36-7a21-8d10-2f4c54c28a01', 'order-1', 'org-main', 'customer-a', 'draft', 'USD', 100, 1),
  ('018f4d4a-7b36-7a21-8d10-2f4c54c28a02', 'order-1', 'org-main', 'customer-b', 'draft', 'USD', 200, 1);

do $test$
declare
  tenant_a_count bigint;
  tenant_b_total bigint;
begin
  select count(*) into tenant_a_count
    from sales.customer_order
   where tenant_id = '018f4d4a-7b36-7a21-8d10-2f4c54c28a01';
  select sum(total_minor_units) into tenant_b_total
    from sales.customer_order
   where tenant_id = '018f4d4a-7b36-7a21-8d10-2f4c54c28a02';
  if tenant_a_count <> 1 or tenant_b_total <> 200 then
    raise exception 'tenant isolation fixture failed';
  end if;
end;
$test$;

do $test$
begin
  begin
    insert into sales.customer_order (
      tenant_id, order_id, organization_id, customer_principal_id,
      state, currency, total_minor_units, version
    ) values (
      '018f4d4a-7b36-7a21-8d10-2f4c54c28a01', 'invalid-order', 'org-main',
      'customer-a', 'invented', 'USD', 0, 1
    );
    raise exception 'expected state check violation';
  exception
    when check_violation then null;
  end;
end;
$test$;

do $test$
declare
  changed bigint;
begin
  update sales.customer_order
     set state = 'placed', version = 2, updated_at = clock_timestamp()
   where tenant_id = '018f4d4a-7b36-7a21-8d10-2f4c54c28a01'
     and order_id = 'order-1'
     and version = 1;
  get diagnostics changed = row_count;
  if changed <> 1 then
    raise exception 'expected one optimistic update, got %', changed;
  end if;

  update sales.customer_order
     set state = 'confirmed', version = 3, updated_at = clock_timestamp()
   where tenant_id = '018f4d4a-7b36-7a21-8d10-2f4c54c28a01'
     and order_id = 'order-1'
     and version = 1;
  get diagnostics changed = row_count;
  if changed <> 0 then
    raise exception 'stale optimistic update was accepted';
  end if;
end;
$test$;

rollback;
````

### FILE: `internal/platform/identity/principal_test.go`

```yaml
block_id: "GO-ENTERPRISE-BACKEND:identity-principal-test-go:v1"
operation: CREATE
provenance: AUTHORED
source: "local verified composition"
license: "LicenseRef-Workspace-Owner"
sha256: "54fa990f8267e9250951c8d3e277d6aaa1c7bee04a19ebb2100fd96df0c59ac5"
variables: []
secrets_allowed: false
```

````go
package identity

import "testing"

func TestPrincipalOrganizationScope(t *testing.T) {
	p := Principal{
		Permissions:   map[string]struct{}{"order:create": {}},
		Organizations: map[string]struct{}{"org-a": {}},
	}
	if !p.AllowedOrganization("org-a") {
		t.Fatal("assigned organization rejected")
	}
	if p.AllowedOrganization("org-b") || p.AllowedOrganization("") {
		t.Fatal("unassigned or empty organization accepted")
	}
	admin := Principal{Permissions: map[string]struct{}{"*": {}}}
	if !admin.AllowedOrganization("org-any") {
		t.Fatal("wildcard principal did not receive global organization access")
	}
}
````

## 6. Configuration surface

`DATABASE_URL`, `OIDC_ISSUER` y `OIDC_AUDIENCE` son obligatorios; bind address, timeouts, pool, claims permitidos y shutdown se externalizan al componer el deployment. Tenant y subject se derivan del token verificado, nunca del body. No se aceptan defaults productivos silenciosos.

## 7. Dependency bill

Go 1.26.7 se verificó con el SHA-256 oficial del archive Windows amd64. pgx v5.10.0, coreos/go-oidc v3.20.0 y transitive modules están fijados por go.mod/go.sum. `go-oidc` es Apache-2.0; pgx y dependencias aplicables son MIT/BSD/Apache según módulo. El producto debe generar SBOM, conservar licencias/notices y reauditar cada upgrade.

## 8. Apply order

Materializar, verificar hashes, instalar Go 1.26.7, ejecutar `go mod download`, `go test ./...`, `go vet ./...` y `go build ./cmd/api`. Componer migrations antes de ejecutar contra PostgreSQL.

## 9. Verification

Gates: gofmt sin diff, tests, vet, build, race test en CI, integration test PostgreSQL 18.6, auth negative tests, workload, SBOM y restore.

## 10. Reconstruction evidence

Estado `REBUILD_VERIFIED / CONDITIONED`, evidencia `GO-RESOURCE-AUTHZ-20260824-V1`:

- 18 archivos materializados con SHA-256 verificado;
- toolchain oficial Go 1.26.7 archive SHA-256 `f4f534a486e4bc3387fa18f08208f2f854b7aaea8a08f2a2d829a914a05abb11`;
- `gofmt`, `go mod download`, `go test ./...`, `go vet ./...` y `go build ./cmd/api`: PASS;
- tests de dominio y contrato HTTP, incluidos rechazo sin bearer y organización no asignada: PASS;
- OIDC discovery/verifier materializado con issuer, audience y RS256 fijados; tenant/subject ya no se aceptan desde el body y `organization_ids` limita recursos;
- PostgreSQL 18.6: foundation + migration + SQL tests + down/up: PASS;
- repository integration: create, replay, hash conflict, tenant-scoped get y transición/outbox: PASS.
- outbox integration: claims disjuntos con `SKIP LOCKED`, ownership al confirmar, release y reclaim: PASS.

Condiciones: probar el verifier contra un issuer/JWKS real y rotación/revocación, versionar permisos y membresías, elegir IDs sortable, añadir publisher/provider real con retry budget, carreras de idempotencia, race/load/security, roles mínimos y deployment. La evidencia de autorización por recurso autoriza reconstrucción/build/unit/HTTP/repository/outbox integration en PostgreSQL real; no autoriza todavía producción.

V402 composed delta: Payment/initial-handover composition. New behavior and tests belong to the explicit runtime/portal packs; source and library release claims remain bounded to their evidence. Existing source provenance is preserved.

Canonical V402 integration: selected by the current profile with exact dependencies and caller overlays. Metadata promotion records byte reconstruction, not closure of every admission/release gate. Payload provenance is unchanged.

V402 composed delta: COMPOSITION_SECURITY_RELEASE_V402.md/json: preserve admitted Go x/mod0.40 graph floor and original DevSkim short-error diagnostics; no new corporate authorship.

V402 composed delta: T2806 connected document reference. Optional host hook, explicit bound document-review kind, exact AWS module closure and preserved security-floor checksums; no unchanged business policy modified. DOCUMENT_REFERENCE_RELEASE_V402.md/json.

V402 composed delta: V402317 connected local API/Next telemetry, current OIDC, fixed official middleware, finite supervised alert/fault/load/WAL recovery. Historical lock kept separate; docs/LOCAL_REFERENCE_OPERATIONS.md. No production admission.
