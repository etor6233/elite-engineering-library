# Secure Operations and Delivery Core — Reconstruction Evidence V1

## Scope

- Date: 2026-08-24
- Pack: `SECURE-OPS-DELIVERY-CORE` 0.1.0
- Python: stdlib unit/semantic validator
- OpenTelemetry Collector: official core distribution 0.158.0 for Windows amd64
- Prometheus: official 3.13.1 for Windows amd64

This evidence authorizes clean reconstruction of the operational contract and syntax/component validation of the supplied OTel/Prometheus baselines. It does not claim that a telemetry backend, alert destination, secret provider, deployer or build attestation service has been selected.

## Results

| Gate | Result | Detail |
|---|---:|---|
| clean materialization | PASS | 6 files; every embedded SHA-256 verified |
| semantic validator tests | PASS | 6/6: valid contract, inline secret, missing signal, placeholder digest, missing runbook, missing restore evidence |
| unresolved example | PASS (negative) | validator returned `FAIL` and exit 2 for placeholders/zero digest |
| OTel archive integrity | PASS | official SHA-256 `72dda387d6b3d51bb3f5c40694d8262458589cb6342cabc96255e2560218a282` |
| OTel component/config validation | PASS | `otelcol 0.158.0 validate --config collector.yaml` |
| Prometheus archive integrity | PASS | official SHA-256 `5409abdcac847984ab7869d7814e6e8cff65b4411d62e7477b960b92eadfa08a` |
| Prometheus rule validation | PASS | `promtool 3.13.1 check rules`; 3 rules found |

## Standards fixed by the contract

- OpenTelemetry Collector 0.158 configuration model; OTLP logs, metrics and traces;
- Prometheus symptom-oriented alert rules;
- CycloneDX 1.7 SBOM declaration;
- SLSA 1.2 provenance declaration;
- OWASP ASVS 5.0.0 and NIST SSDF 1.1 as project verification inputs, not claimed certifications.

## Remaining conditions

Actual projects must replace every placeholder, use a nonzero immutable artifact digest, connect a real secret/telemetry/deployment provider, generate and verify the declared evidence, exercise alert delivery, redaction, rotation, rollback and restore, and document any time-bounded security exception.
