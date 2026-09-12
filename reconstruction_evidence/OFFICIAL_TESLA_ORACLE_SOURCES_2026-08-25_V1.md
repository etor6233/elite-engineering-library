# Official Tesla and Oracle Sources — 2026-08-25 V1

## Exact sources

| Source | Release / commit | Archive evidence | License / lock evidence |
|---|---|---|---|
| `teslamotors/vehicle-command` | `v0.4.1` / `49977a18fd68567501d59e16a6c9e4a8b9348544` | 432,966 bytes; SHA-256 `c768826cb323543ab6c7867b43ec8db34444853ea6f0f31e20274e905cf020ec` | Apache-2.0 `LICENSE` `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30`; `go.sum` `9493a17b77c3bc9368b4c5cb0beab4511fb1afd448695b850944095dba676874` |
| `teslamotors/fleet-telemetry` | `v0.9.4` / `d64c73ab65e7c5fb5fc12b35fe507e2c6054227b` | 942,798 bytes; SHA-256 `c8166638de82a55f17b9f92d76335d95e7f111a80b11fdee7d6c7aa1e35ca344` | Apache-2.0 `LICENSE` `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30`; `go.sum` `9128da99915057f407b0bd24c9255ba885dde74d9334334f8000fc38083294ce` |
| `oracle-samples/db-sample-schemas` | `v23.3` / `e3325a83e56c516815844025418a96ecaf219751` | 12,195,946 bytes; SHA-256 `d8a8e07524c0ad99406f0a3822cfc6c5f8cafce9b687cef6fd974f16b609d3f2` | UPL-1.0 `LICENSE.txt` `2da4f8e1f04662e5db9b224a20dfd13db8bc396398271d607bda0343212fbce3` |

The exact archives contained 168, 158 and 310 files respectively. No Git repository was created.

## Executed verification

An official Go 1.26.7 distribution verified by its published SHA-256 executed the Tesla source tests.

- `vehicle-command`: `go test ./...` passed, including authentication, dispatcher, Schnorr, account, cache, CLI, Internet connector, protocol, proxy and vehicle packages. Packages without tests were reported honestly.
- `fleet-telemetry`: MQTT, Redis, simple datastore/transformers, logger, metrics adapters, Airbrake and telemetry packages passed. The full build failed on Windows because Kafka (`librdkafka`) and ZeroMQ native symbols were unavailable; dependent command/streaming/integration packages consequently failed. The failure is retained, not converted to a pass.
- Oracle schemas were inspected and hashed but not installed because the audit environment did not provide Oracle Database 19c+ or a disposable privileged schema. Oracle documents them as sample schemas.

## Admission

Tesla code is immediately acquirable official integration code for Tesla Fleet API; it is not a generic electric-vehicle franchise backend, factory system or telemetry protocol for other manufacturers. Fleet Telemetry remains conditioned on its documented native/runtime and real-service requirements.

Oracle Customer Orders is `SAMPLE_ONLY`: it can provide Oracle-specific fixtures and learning material, but it cannot govern the portable production domain, security or transaction design.
