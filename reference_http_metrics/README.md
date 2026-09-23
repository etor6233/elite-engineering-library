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

## Current connected host (V402317)

The unchanged closed-attribute instrument now lives in core/metrics.go; metrics.go
retains the historical wrapper. The ordinary API enables it only with explicit
HTTP_METRICS_ENABLED=true and HTTP_METRICS_ADDRESS=127.0.0.1:<port>. Its real OIDC
verifier is retained. The historical reference main's synthetic principal is never
installed in that host. Use ci/run_observed_reference.py and
docs/LOCAL_REFERENCE_OPERATIONS.md for the current connected profile. Historical
source-lock observations remain under historical_v375_through_v400. The old
run_reference.py requires its original external schema lock and binary; it is
preserved for that historical cohort, not the current acceptance command.
