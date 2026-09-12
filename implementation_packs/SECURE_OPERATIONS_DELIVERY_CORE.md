# Secure Operations and Delivery Core

V374: PlatformHighErrorRate now divides failed-request rate by actual total-request rate. The previous denominator floor of1request/second hid low-traffic outages. The existing1% threshold and10minute hold remain unchanged. Nine source-engine test scenarios /13alert assertions cover sparse outages, counter reset, exact/above threshold, healthy/zero traffic, pending hold, short spike and recovery. Run the target-admitted Prometheus promtool with `promtool check rules ops/prometheus/platform.rules.yml` and `promtool test rules ops/prometheus/platform.rules.test.yml`. Reference evidence uses the exact V372 Prometheus3.14.0 binary; production must bind its own artifact, labels, traffic/SLI policy and runbooks. This error-ratio alert does not implement multi-window error-budget burn alerts. Evidence: reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md.

V345: el ejecutor limita stdout/stderr durante la lectura, conserva el plazo
aunque un descendiente retenga la tubería y recoge sólo el proceso directo.
Errores de ejecución estáticos; sin traceback con argv privado. No instala
supervisor de servicios, reinicios, colector ni retención. Evidencia:
reconstruction_evidence/OPERATIONAL_PROCESS_BOUNDARY_V345.md.

V315: el gate operativo exige las referencias SBOM/provenance en ambos modos;
ninguna referencia vacía/nula/malformada desactiva el chequeo de archivo.
Errores de inspección retornan FAIL con el campo afectado. Doce tests,28
negativos de referencias, control completo, escape/directorio/ausencia y fallo
de acceso inyectado; reconstrucción2/2. Evidencia OPERATIONAL_EVIDENCE_GATE_V315.md.
La presencia de archivos no acredita contenido, hashes, firmas ni producción;
el gate de admisión productiva separado conserva toda su autoridad.

## 1. Metadata

```yaml
pack_id: "SECURE-OPS-DELIVERY-CORE"
pack_version: "1.1.4"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa contratos ejecutables separados para readiness operativo y admisión productiva hash-bound, una frontera genérica para herramientas oficiales y los ocho adaptadores semánticos exigidos por el gate final: proveedores, carga, recovery, ofensiva, edge, identidad, deploy/rollback y aceptación empresarial. Business acceptance fija GitHub Spec Kit 1.0.1 y Microsoft Playwright 1.62.1, exige escenarios reales con efectos, aprobaciones independientes, defectos gobernados y owner signoff. No convierte specs generadas, fixtures, exit 0 ni aprobación sintética en evidencia productiva."
stacks:
  - "Python 3.14+ stdlib validator"
  - "OpenTelemetry Collector 0.158 configuration contract"
  - "Prometheus rule contract"
  - "Exact Stripe/Mercado Pago/Amazon/Google/Meta/TikTok/Firebase provider sources"
  - "Kubernetes kubectl 1.37.0 exact artifact identity"
  - "GitHub Spec Kit 1.0.1 + Microsoft Playwright 1.62.1 acceptance authority"
compatible_with: ["PORTABLE-CI-GATE-RUNNER 0.1.x", "DEPENDENCY-LICENSE-EVIDENCE-CORE 0.1.x"]
incompatible_with: []
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources:
  - "https://opentelemetry.io/docs/collector/configuration/"
  - "https://prometheus.io/docs/practices/alerting/"
  - "https://cyclonedx.org/specification/overview/"
  - "https://slsa.dev/spec/v1.2/"
  - "https://owasp.org/www-project-application-security-verification-standard/"
  - "https://csrc.nist.gov/pubs/sp/800/218/final"
  - "https://developers.cloudflare.com/api/resources/rulesets/methods/list/"
  - "https://openid.net/certification/about-conformance-suite/"
  - "https://gitlab.com/openid/conformance-suite/-/releases/release-v5.2.4"
  - "https://docs.stripe.com/webhooks"
  - "https://docs.stripe.com/error-low-level"
  - "https://www.mercadopago.com.ar/developers/es/docs/checkout-bricks/additional-content/your-integrations/notifications/webhooks"
  - "https://docs.aws.amazon.com/sdkref/latest/guide/feature-retry-behavior.html"
  - "https://developers.google.com/merchant/api/guides/quotas-limits/quotas"
  - "https://www.postgresql.org/docs/18/app-pgverifybackup.html"
  - "https://grafana.com/docs/k6/latest/using-k6/thresholds/"
  - "https://github.com/grafana/k6/releases/tag/v2.2.0"
  - "https://www.zaproxy.org/docs/automate/automation-framework/"
  - "https://github.com/cloudflare/cloudflare-go/releases/tag/v7.9.0"
  - "https://kubernetes.io/docs/concepts/workloads/controllers/deployment/"
  - "https://github.com/kubernetes/kubernetes/releases/tag/v1.37.0"
  - "https://dl.k8s.io/release/v1.37.0/bin/windows/amd64/kubectl.exe.sha256"
  - "https://github.com/github/spec-kit/releases/tag/v1.0.1"
  - "https://github.com/microsoft/playwright/releases/tag/v1.62.1"
verified_at: "2026-09-10"
```

Los bloques son `AUTHORED`; las fuentes públicas definen estándares y prácticas, no aportan código copiado. El pack no instala un proveedor ni afirma conformidad SLSA/ASVS: convierte decisiones y evidencias mínimas en un gate fail-closed.

## 2. Applicability

Use for every deployable service that must declare secrets, telemetry, SLOs, supply-chain evidence, release/rollback and recovery readiness. Reject the example as production evidence until every placeholder is replaced by target-owned values and executable proof.

## 3. Architecture contract

- valores secretos quedan fuera de repositorio; sólo se admiten referencias `secret-manager://`, `vault://`, `workload-identity://` o `kms://`;
- logs, métricas y trazas son obligatorios y deben declarar redacción, correlación, retención y backend;
- cada SLO incluye SLI, objetivo, ventana, owner, runbook y alertas burn-rate;
- el release referencia artefacto inmutable `sha256:<64 hex>`, SBOM CycloneDX 1.7 y provenance SLSA 1.2;
- críticos conocidos bloquean release salvo excepción con caducidad y owner;
- rollback y restore probado son datos obligatorios, no frases;
- OTel y Prometheus son configuraciones base que deben validarse contra la distribución elegida.
- `READY_TO_BUILD` nunca equivale a `PRODUCTION_ADMITTED`; la segunda autoridad exige ocho receipts frescos ligados al mismo proyecto, entorno y digest de release.

## 4. Exact file manifest

```text
CREATE ops/readiness/operational-readiness.example.json
CREATE ops/readiness/validate_operational_readiness.py
CREATE ops/readiness/test_validate_operational_readiness.py
CREATE ops/otel/collector.yaml
CREATE ops/prometheus/platform.rules.yml
CREATE ops/prometheus/platform.rules.test.yml
CREATE ops/runbooks/release-rollback.md
CREATE production_admission_gate/production-admission.template.json
CREATE production_admission_gate/validate_production_admission.py
CREATE production_admission_gate/test_validate_production_admission.py
CREATE production_admission_gate/README.md
CREATE production_admission_gate/official-tool-run.template.json
CREATE production_admission_gate/run_official_tool.py
CREATE production_admission_gate/test_run_official_tool.py
CREATE production_admission_gate/k6-load-admission.template.json
CREATE production_admission_gate/validate_k6_load_admission.py
CREATE production_admission_gate/test_validate_k6_load_admission.py
CREATE production_admission_gate/postgres-recovery-admission.template.json
CREATE production_admission_gate/validate_postgres_recovery_admission.py
CREATE production_admission_gate/test_validate_postgres_recovery_admission.py
CREATE production_admission_gate/zap-offensive-admission.template.json
CREATE production_admission_gate/validate_zap_offensive_admission.py
CREATE production_admission_gate/test_validate_zap_offensive_admission.py
CREATE production_admission_gate/cloudflare-edge-admission.template.json
CREATE production_admission_gate/validate_cloudflare_edge_admission.py
CREATE production_admission_gate/test_validate_cloudflare_edge_admission.py
CREATE production_admission_gate/identity-authorization-admission.template.json
CREATE production_admission_gate/validate_identity_authorization_admission.py
CREATE production_admission_gate/test_validate_identity_authorization_admission.py
CREATE production_admission_gate/deploy-rollback-admission.template.json
CREATE production_admission_gate/validate_deploy_rollback_admission.py
CREATE production_admission_gate/test_validate_deploy_rollback_admission.py
CREATE production_admission_gate/business-acceptance-admission.template.json
CREATE production_admission_gate/validate_business_acceptance_admission.py
CREATE production_admission_gate/test_validate_business_acceptance_admission.py
CREATE production_admission_gate/provider-admission.template.json
CREATE production_admission_gate/validate_provider_admission.py
CREATE production_admission_gate/test_validate_provider_admission.py
```

## 5. Materialization blocks

### FILE: `ops/readiness/operational-readiness.example.json`

```yaml
block_id: "SECURE-OPS:readiness-example:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "5e5a96dd23de6d5cd0804e6996339101cbcd5cc6d343461aa7a5527c500ea6ac"
variables: []
secrets_allowed: false
```

````json
{
  "schema_version": "elite.operational-readiness.v1",
  "service": {"name": "replace-me", "owner": "team:replace-me", "tier": 1},
  "environments": ["staging", "production"],
  "secrets": {
    "provider": "replace-me",
    "references": ["secret-manager://replace-me/database-credential"],
    "rotation_days": 90,
    "workload_identity": true,
    "break_glass_runbook": "ops/runbooks/break-glass.md"
  },
  "telemetry": {
    "signals": ["logs", "metrics", "traces"],
    "protocol": "otlp",
    "backend": "replace-me",
    "correlation_fields": ["service.name", "deployment.environment.name", "trace_id", "tenant.id"],
    "redaction_test": "evidence/redaction-test.json",
    "retention_policy": "policy/telemetry-retention.md"
  },
  "slos": [{
    "name": "api-availability",
    "objective": 0.999,
    "window_days": 30,
    "sli": "1 - (sum(rate(http_server_request_duration_seconds_count{http_response_status_code=~\"5..\"}[5m])) / sum(rate(http_server_request_duration_seconds_count[5m])))",
    "owner": "team:replace-me",
    "runbook": "ops/runbooks/api-availability.md",
    "page_burn_rates": [14.4, 6.0],
    "ticket_burn_rates": [3.0, 1.0]
  }],
  "supply_chain": {
    "lockfiles": ["replace-with-selected-runtime-lockfile"],
    "sbom": {"format": "CycloneDX", "spec_version": "1.7", "artifact": "evidence/bom.cdx.json"},
    "provenance": {"format": "SLSA", "spec_version": "1.2", "artifact": "evidence/provenance.json"},
    "vulnerability_policy": {"max_critical": 0, "max_high": 0, "exception_register": "security/exceptions.json"},
    "license_policy": "security/allowed-licenses.json",
    "secret_scan_evidence": "evidence/secret-scan.json"
  },
  "release": {
    "artifact_digest": "sha256:0000000000000000000000000000000000000000000000000000000000000000",
    "deployment_strategy": "canary",
    "migration_strategy": "expand-migrate-contract",
    "health_gates": ["availability", "latency", "error-budget", "queue-lag"],
    "rollback_command": ["replace-me-with-deployer", "rollback", "--digest", "replace-me"],
    "rollback_test": "evidence/rollback-test.json",
    "approval_record": "evidence/release-approval.json"
  },
  "recovery": {
    "rpo_minutes": 15,
    "rto_minutes": 60,
    "backup_policy": "ops/runbooks/backup-policy.md",
    "restore_evidence": "reconstruction_evidence/POSTGRES_BACKUP_RESTORE_CORE_2026-08-24_V1.md",
    "dependency_order": ["identity", "database", "messaging", "backend", "web"],
    "last_drill_utc": "2026-08-24T22:00:13Z"
  }
}
````

### FILE: `ops/readiness/validate_operational_readiness.py`

```yaml
block_id: "SECURE-OPS:validator:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "950e8ff3aa63dea22b10c6bccf0501a5b57b3139fb0df51c8e7d65d79ce799ac"
variables: []
secrets_allowed: false
```

````python
from __future__ import annotations

import argparse
import json
import re
import sys
from pathlib import Path
from typing import Any

SECRET_KEYS = re.compile(r"(^|_)(password|passwd|token|secret|private_key|api_key|client_secret)($|_)", re.I)
SECRET_MARKERS = ("-----BEGIN PRIVATE KEY-----", "-----BEGIN RSA PRIVATE KEY-----", "AKIA")
SECRET_REF = re.compile(r"^(secret-manager|vault|workload-identity|kms)://[^\s]+$")
DIGEST = re.compile(r"^sha256:[0-9a-f]{64}$")


def require(condition: bool, message: str, errors: list[str]) -> None:
    if not condition:
        errors.append(message)


def reject_embedded_secrets(value: Any, path: str, errors: list[str]) -> None:
    if isinstance(value, dict):
        for key, child in value.items():
            child_path = f"{path}.{key}"
            if SECRET_KEYS.search(key) and key not in {"secret_scan_evidence", "secrets"}:
                errors.append(f"{child_path}: secret-shaped key is forbidden")
            reject_embedded_secrets(child, child_path, errors)
    elif isinstance(value, list):
        for index, child in enumerate(value):
            reject_embedded_secrets(child, f"{path}[{index}]", errors)
    elif isinstance(value, str) and any(marker in value for marker in SECRET_MARKERS):
        errors.append(f"{path}: embedded credential marker is forbidden")


def nonempty_string(value: Any) -> bool:
    return isinstance(value, str) and bool(value.strip()) and "replace-me" not in value


def validate(document: dict[str, Any], require_evidence_files: bool = False, root: Path | None = None) -> list[str]:
    errors: list[str] = []
    require(document.get("schema_version") == "elite.operational-readiness.v1", "schema_version: unsupported", errors)
    reject_embedded_secrets(document, "$", errors)

    service = document.get("service", {})
    require(nonempty_string(service.get("name")), "service.name: unresolved", errors)
    require(nonempty_string(service.get("owner")), "service.owner: unresolved", errors)
    require(service.get("tier") in {1, 2, 3}, "service.tier: must be 1, 2 or 3", errors)
    require(set(document.get("environments", [])) >= {"staging", "production"}, "environments: staging and production required", errors)

    secrets = document.get("secrets", {})
    require(nonempty_string(secrets.get("provider")), "secrets.provider: unresolved", errors)
    refs = secrets.get("references", [])
    require(bool(refs) and all(isinstance(ref, str) and SECRET_REF.fullmatch(ref) for ref in refs), "secrets.references: invalid secret reference", errors)
    require(secrets.get("workload_identity") is True, "secrets.workload_identity: must be true", errors)
    require(isinstance(secrets.get("rotation_days"), int) and 1 <= secrets["rotation_days"] <= 365, "secrets.rotation_days: must be 1..365", errors)
    require(nonempty_string(secrets.get("break_glass_runbook")), "secrets.break_glass_runbook: unresolved", errors)

    telemetry = document.get("telemetry", {})
    require(set(telemetry.get("signals", [])) == {"logs", "metrics", "traces"}, "telemetry.signals: logs, metrics and traces required", errors)
    require(telemetry.get("protocol") == "otlp", "telemetry.protocol: must be otlp", errors)
    for key in ("backend", "redaction_test", "retention_policy"):
        require(nonempty_string(telemetry.get(key)), f"telemetry.{key}: unresolved", errors)
    require({"service.name", "deployment.environment.name", "trace_id"} <= set(telemetry.get("correlation_fields", [])), "telemetry.correlation_fields: required fields missing", errors)

    slos = document.get("slos", [])
    require(isinstance(slos, list) and len(slos) > 0, "slos: at least one required", errors)
    for index, slo in enumerate(slos if isinstance(slos, list) else []):
        prefix = f"slos[{index}]"
        for key in ("name", "sli", "owner", "runbook"):
            require(nonempty_string(slo.get(key)), f"{prefix}.{key}: unresolved", errors)
        objective = slo.get("objective")
        require(isinstance(objective, (int, float)) and 0.9 <= objective < 1, f"{prefix}.objective: must be >=0.9 and <1", errors)
        require(slo.get("window_days") in {7, 28, 30, 90}, f"{prefix}.window_days: unsupported", errors)
        require(len(slo.get("page_burn_rates", [])) >= 2, f"{prefix}.page_burn_rates: two windows required", errors)
        require(len(slo.get("ticket_burn_rates", [])) >= 2, f"{prefix}.ticket_burn_rates: two windows required", errors)

    chain = document.get("supply_chain", {})
    require(bool(chain.get("lockfiles")) and all(nonempty_string(v) for v in chain.get("lockfiles", [])), "supply_chain.lockfiles: unresolved", errors)
    require(chain.get("sbom", {}).get("format") == "CycloneDX" and chain.get("sbom", {}).get("spec_version") == "1.7", "supply_chain.sbom: CycloneDX 1.7 required", errors)
    require(chain.get("provenance", {}).get("format") == "SLSA" and chain.get("provenance", {}).get("spec_version") == "1.2", "supply_chain.provenance: SLSA 1.2 required", errors)
    for artifact in ("sbom", "provenance"):
        require(nonempty_string(chain.get(artifact, {}).get("artifact")), f"supply_chain.{artifact}.artifact: evidence reference required", errors)
    policy = chain.get("vulnerability_policy", {})
    require(policy.get("max_critical") == 0, "supply_chain.vulnerability_policy.max_critical: must be zero", errors)
    for key in ("license_policy", "secret_scan_evidence"):
        require(nonempty_string(chain.get(key)), f"supply_chain.{key}: unresolved", errors)

    release = document.get("release", {})
    require(bool(DIGEST.fullmatch(str(release.get("artifact_digest", "")))) and set(str(release.get("artifact_digest", "")).split(":", 1)[1]) != {"0"}, "release.artifact_digest: real immutable digest required", errors)
    require(release.get("deployment_strategy") in {"canary", "blue-green", "rolling"}, "release.deployment_strategy: unsupported", errors)
    require(release.get("migration_strategy") == "expand-migrate-contract", "release.migration_strategy: expand-migrate-contract required", errors)
    require(len(release.get("health_gates", [])) >= 3, "release.health_gates: at least three required", errors)
    rollback = release.get("rollback_command", [])
    require(isinstance(rollback, list) and len(rollback) >= 2 and all(nonempty_string(v) for v in rollback), "release.rollback_command: argv array required", errors)
    for key in ("rollback_test", "approval_record"):
        require(nonempty_string(release.get(key)), f"release.{key}: unresolved", errors)

    recovery = document.get("recovery", {})
    require(isinstance(recovery.get("rpo_minutes"), int) and recovery["rpo_minutes"] > 0, "recovery.rpo_minutes: positive integer required", errors)
    require(isinstance(recovery.get("rto_minutes"), int) and recovery["rto_minutes"] > 0, "recovery.rto_minutes: positive integer required", errors)
    for key in ("backup_policy", "restore_evidence", "last_drill_utc"):
        require(nonempty_string(recovery.get(key)), f"recovery.{key}: unresolved", errors)
    require(len(recovery.get("dependency_order", [])) >= 3, "recovery.dependency_order: incomplete", errors)

    if require_evidence_files:
        base = root or Path.cwd()
        evidence_paths = {
            "telemetry.redaction_test": telemetry.get("redaction_test"),
            "supply_chain.sbom.artifact": chain.get("sbom", {}).get("artifact"),
            "supply_chain.provenance.artifact": chain.get("provenance", {}).get("artifact"),
            "supply_chain.secret_scan_evidence": chain.get("secret_scan_evidence"),
            "release.rollback_test": release.get("rollback_test"),
            "release.approval_record": release.get("approval_record"),
            "recovery.restore_evidence": recovery.get("restore_evidence"),
        }
        for field, item in evidence_paths.items():
            if not nonempty_string(item):
                errors.append(f"{field}: evidence reference required")
                continue
            try:
                candidate = (base / item).resolve()
                try:
                    candidate.relative_to(base.resolve())
                except ValueError:
                    errors.append(f"{field}: evidence path escapes root")
                    continue
                require(candidate.is_file(), f"{field}: evidence file not found", errors)
            except (OSError, ValueError, RuntimeError):
                errors.append(f"{field}: cannot inspect file")
    return sorted(set(errors))


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("document", type=Path)
    parser.add_argument("--require-evidence-files", action="store_true")
    parser.add_argument("--root", type=Path, default=Path.cwd())
    args = parser.parse_args()
    document = json.loads(args.document.read_text(encoding="utf-8"))
    errors = validate(document, args.require_evidence_files, args.root)
    print(json.dumps({"status": "PASS" if not errors else "FAIL", "errors": errors}, indent=2))
    return 0 if not errors else 2


if __name__ == "__main__":
    sys.exit(main())
````

### FILE: `ops/readiness/test_validate_operational_readiness.py`

```yaml
block_id: "SECURE-OPS:validator-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "7030fa5a1294d237eb8ffdafdc97fb6d092bc260198c7ecb82b9e2fa41907098"
variables: []
secrets_allowed: false
```

````python
from __future__ import annotations

import copy
import json
import tempfile
import unittest
from unittest.mock import patch
from pathlib import Path

from validate_operational_readiness import validate


class ReadinessTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        cls.template = json.loads((Path(__file__).parent / "operational-readiness.example.json").read_text(encoding="utf-8"))

    def valid(self) -> dict:
        value = copy.deepcopy(self.template)
        value["service"] = {"name": "orders", "owner": "team:commerce", "tier": 1}
        value["secrets"]["provider"] = "production-vault"
        value["secrets"]["references"] = ["vault://production/orders/database"]
        value["secrets"]["break_glass_runbook"] = "ops/runbooks/break-glass.md"
        value["telemetry"]["backend"] = "regional-otel-gateway"
        value["slos"][0]["owner"] = "team:commerce"
        value["supply_chain"]["lockfiles"] = ["go.sum"]
        value["release"]["artifact_digest"] = "sha256:" + "a" * 64
        value["release"]["rollback_command"] = ["deployer", "rollback", "--digest", "sha256:" + "b" * 64]
        return value

    def test_valid_contract_passes(self) -> None:
        self.assertEqual([], validate(self.valid()))

    def test_inline_secret_key_fails(self) -> None:
        value = self.valid()
        value["database_password"] = "do-not-store-this"
        self.assertTrue(any("secret-shaped" in error for error in validate(value)))

    def test_missing_signal_fails(self) -> None:
        value = self.valid()
        value["telemetry"]["signals"] = ["logs", "metrics"]
        self.assertTrue(any("telemetry.signals" in error for error in validate(value)))

    def test_placeholder_digest_fails(self) -> None:
        value = self.valid()
        value["release"]["artifact_digest"] = "sha256:" + "0" * 64
        self.assertTrue(any("artifact_digest" in error for error in validate(value)))

    def test_slo_without_runbook_fails(self) -> None:
        value = self.valid()
        value["slos"][0]["runbook"] = "replace-me"
        self.assertTrue(any("runbook" in error for error in validate(value)))

    def test_missing_restore_evidence_fails_closed(self) -> None:
        value = self.valid()
        value["recovery"]["restore_evidence"] = ""
        self.assertTrue(any("restore_evidence" in error for error in validate(value)))

    def write_evidence(self, root: Path, value: dict) -> None:
        # These files test presence only; they are not operational proofs.
        refs = [value["telemetry"]["redaction_test"], value["supply_chain"]["sbom"]["artifact"],
                value["supply_chain"]["provenance"]["artifact"], value["supply_chain"]["secret_scan_evidence"],
                value["release"]["rollback_test"], value["release"]["approval_record"], value["recovery"]["restore_evidence"]]
        for ref in refs:
            path = root / ref
            path.parent.mkdir(parents=True, exist_ok=True)
            path.write_text("{}\n", encoding="utf-8")

    def test_complete_evidence_file_set_passes_presence_gate(self) -> None:
        with tempfile.TemporaryDirectory() as folder:
            value = self.valid(); root = Path(folder)
            self.write_evidence(root, value)
            self.assertEqual([], validate(value, True, root))

    def test_artifact_references_cannot_disable_evidence_check(self) -> None:
        missing = object()
        for artifact in ("sbom", "provenance"):
            for mode in (False, True):
                for bad in (missing, None, "", "replace-me", 42, [], {}):
                    with self.subTest(artifact=artifact, require_files=mode, bad_type=type(bad).__name__):
                        with tempfile.TemporaryDirectory() as folder:
                            root = Path(folder); value = self.valid()
                            self.write_evidence(root, value)
                            if bad is missing:
                                del value["supply_chain"][artifact]["artifact"]
                            else:
                                value["supply_chain"][artifact]["artifact"] = bad
                            errors = validate(value, mode, root)
                            self.assertTrue(any(f"supply_chain.{artifact}.artifact" in e for e in errors), errors)

    def test_missing_or_directory_evidence_is_rejected(self) -> None:
        with tempfile.TemporaryDirectory() as folder:
            root = Path(folder); value = self.valid(); self.write_evidence(root, value)
            for artifact in ("sbom", "provenance"):
                path = root / value["supply_chain"][artifact]["artifact"]
                contents = path.read_bytes(); path.unlink()
                self.assertTrue(validate(value, True, root))
                path.mkdir()
                self.assertTrue(validate(value, True, root))
                path.rmdir(); path.write_bytes(contents)

    def test_evidence_escape_is_rejected_even_when_file_exists(self) -> None:
        with tempfile.TemporaryDirectory() as folder:
            root = Path(folder) / "project"; root.mkdir()
            value = self.valid(); self.write_evidence(root, value)
            outside = root.parent / "outside.json"; outside.write_text("{}", encoding="utf-8")
            value["supply_chain"]["sbom"]["artifact"] = "../outside.json"
            self.assertTrue(any("escapes root" in e for e in validate(value, True, root)))

    def test_invalid_filesystem_evidence_is_reported(self) -> None:
        with tempfile.TemporaryDirectory() as folder:
            root = Path(folder); value = self.valid(); self.write_evidence(root, value)
            value["supply_chain"]["sbom"]["artifact"] = "evidence/invalid\x00.json"
            errors = validate(value, True, root)
            self.assertTrue(errors)

    def test_unreadable_evidence_returns_validation_failure(self) -> None:
        with tempfile.TemporaryDirectory() as folder:
            root = Path(folder); value = self.valid(); self.write_evidence(root, value)
            with patch.object(Path, "is_file", side_effect=PermissionError("synthetic access failure")):
                errors = validate(value, True, root)
            self.assertTrue(any("cannot inspect file" in e for e in errors), errors)


if __name__ == "__main__":
    unittest.main()
````

### FILE: `ops/otel/collector.yaml`

```yaml
block_id: "SECURE-OPS:otel:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "d477634639f5d3cc119aac13e84c976733caa2e4788690172ce5ff6b38628562"
variables: []
secrets_allowed: false
```

````yaml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 127.0.0.1:4317
      http:
        endpoint: 127.0.0.1:4318
processors:
  memory_limiter:
    check_interval: 1s
    limit_mib: 512
    spike_limit_mib: 128
  batch:
    send_batch_size: 1024
    timeout: 5s
exporters:
  otlphttp/platform:
    endpoint: ${env:OTEL_EXPORTER_OTLP_ENDPOINT}
    headers:
      Authorization: ${env:OTEL_EXPORTER_AUTHORIZATION}
    sending_queue:
      enabled: true
      queue_size: 10000
    retry_on_failure:
      enabled: true
extensions:
  health_check:
    endpoint: 127.0.0.1:13133
service:
  extensions: [health_check]
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [otlphttp/platform]
    metrics:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [otlphttp/platform]
    logs:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [otlphttp/platform]
````

### FILE: `ops/prometheus/platform.rules.yml`

```yaml
block_id: "SECURE-OPS:prom-rules:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "cc718b7bbe35de3f230eeb47c854c270580dbb488101ee5f30c0b758e1d967b0"
variables: []
secrets_allowed: false
```

````yaml
groups:
  - name: platform-symptoms
    interval: 30s
    rules:
      - alert: PlatformHighErrorRate
        expr: sum(rate(http_server_request_duration_seconds_count{http_response_status_code=~"5.."}[5m])) / sum(rate(http_server_request_duration_seconds_count[5m])) > 0.01
        for: 10m
        labels:
          severity: page
        annotations:
          summary: User-visible server error rate exceeds one percent
          runbook_url: https://replace.invalid/runbooks/api-availability
      - alert: PlatformOutboxBacklog
        expr: max(platform_outbox_oldest_unpublished_age_seconds) > 300
        for: 10m
        labels:
          severity: ticket
        annotations:
          summary: Transactional outbox has unpublished events older than five minutes
          runbook_url: https://replace.invalid/runbooks/outbox-backlog
      - alert: PlatformBackupStale
        expr: time() - max(platform_backup_last_success_unixtime) > 90000
        for: 15m
        labels:
          severity: page
        annotations:
          summary: No successful platform backup has been observed within twenty-five hours
          runbook_url: https://replace.invalid/runbooks/backup-stale
````

### FILE: `ops/runbooks/release-rollback.md`

```yaml
block_id: "SECURE-OPS:release-runbook:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "7a13d5f65fff61ed249fe421b6cd5bf2ff94a04c32b388bf0338aaaaa205f09a"
variables: []
secrets_allowed: false
```

````markdown
# Release and rollback runbook

Before release, record the source revision, immutable artifact digest, builder identity, SBOM, provenance, vulnerability/license/secret scans, migration compatibility, restore evidence, approver and planned health gates. Never deploy a mutable tag as the release identity.

Use expand-migrate-contract: deploy additive schema first, backfill with resumable jobs, shift reads/writes, observe, and remove old schema only in a later release. A rollback must not require reversing a destructive migration.

During canary, compare availability, tail latency, error-budget burn, saturation, queue/outbox lag and business invariants with the stable population. Stop promotion on any failed gate. Roll back to the previously recorded digest using the exact argv command in the readiness record; do not rebuild it.

After rollback, verify traffic, database compatibility, job ownership, event publication and external side effects. Preserve telemetry and an incident timeline. A rollback test is complete only when it proves the old artifact can serve against the current compatible schema.
````

### FILE: `production_admission_gate/production-admission.template.json`

```yaml
block_id: "SECURE-OPS:production-admission-template:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://production-admission/template"
license: "LicenseRef-Workspace-Owner"
sha256: "3cb8cf605f48cd40d716d332bd2f450205f761fd26f5fc66608acbd6e5b24c78"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite-production-admission/v1",
  "project_id": "",
  "environment": "production",
  "release_digest": "",
  "evaluated_at": "",
  "controls": []
}
````

### FILE: `production_admission_gate/validate_production_admission.py`

```yaml
block_id: "SECURE-OPS:production-admission-validator:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://production-admission/validator"
license: "LicenseRef-Workspace-Owner"
sha256: "95064dc965cd8399d3f649261abd566286fed15d738b52f108ff4da2d2f1a74f"
variables: []
secrets_allowed: false
```

````python
from __future__ import annotations

import argparse
from datetime import datetime, timezone
import hashlib
import json
from pathlib import Path, PurePosixPath
import re
import sys

SCHEMA = "elite-production-admission/v1"
CONTROLS = {
    "EDGE_CDN_WAF", "IDENTITY_AUTHORIZATION", "PROVIDERS", "POSTGRES_RECOVERY",
    "LOAD_RESILIENCE", "OFFENSIVE_SECURITY", "DEPLOY_ROLLBACK", "BUSINESS_ACCEPTANCE",
}
REQUIRED_ASSERTIONS = {
    "EDGE_CDN_WAF": {"dns_proxy_active", "waf_rulesets_active", "tls_policy_pass", "origin_bypass_blocked", "cache_behavior_pass"},
    "IDENTITY_AUTHORIZATION": {"oidc_conformance_pass", "role_positive_journeys_pass", "unauthorized_negative_journeys_pass", "tenant_isolation_pass", "session_revocation_rotation_pass"},
    "PROVIDERS": {"all_required_sandboxes_pass", "webhook_auth_pass", "idempotency_reconciliation_pass", "rate_limit_retry_pass", "terms_cost_scope_approved"},
    "POSTGRES_RECOVERY": {"pg_verifybackup_pass", "restore_boot_pass", "data_integrity_pass", "pitr_pass", "rpo_met", "rto_met"},
    "LOAD_RESILIENCE": {"k6_thresholds_pass", "expected_workload_represented", "soak_pass", "saturation_recovery_pass", "no_error_budget_breach"},
    "OFFENSIVE_SECURITY": {"zap_automation_pass", "authenticated_scope_pass", "api_scope_pass", "zero_unaccepted_high_critical", "manual_review_pass"},
    "DEPLOY_ROLLBACK": {"immutable_digest_deployed", "canary_health_pass", "rollout_status_pass", "rollback_executed", "previous_digest_restored", "migration_compatibility_pass"},
    "BUSINESS_ACCEPTANCE": {"acceptance_scenarios_pass", "finance_operations_security_approved", "no_open_sev1_sev2", "approvers_distinct", "owner_signoff"},
}
DIGEST = re.compile(r"^sha256:[0-9a-f]{64}$")
HEX = re.compile(r"^[0-9a-f]{64}$")


def utc(value: object, label: str, errors: list[str]) -> datetime | None:
    try:
        if not isinstance(value, str) or not value.endswith("Z"):
            raise ValueError
        return datetime.fromisoformat(value[:-1] + "+00:00")
    except ValueError:
        errors.append(f"{label} must be ISO-8601 UTC ending Z")
        return None


def evidence(root: Path, item: object, label: str, errors: list[str]) -> None:
    if not isinstance(item, dict) or set(item) != {"path", "sha256", "bytes"}:
        errors.append(f"{label} must contain exactly path, sha256 and bytes")
        return
    raw = item.get("path")
    if not isinstance(raw, str) or "\\" in raw:
        errors.append(f"{label}.path must be canonical project-relative")
        return
    rel = PurePosixPath(raw)
    if rel.is_absolute() or any(part in {"", ".", ".."} for part in rel.parts):
        errors.append(f"{label}.path must be canonical project-relative")
        return
    candidate = root.joinpath(*rel.parts)
    try:
        resolved = candidate.resolve(strict=True)
        resolved.relative_to(root)
    except (FileNotFoundError, ValueError):
        errors.append(f"{label}.path is missing or escapes project root")
        return
    if candidate.is_symlink() or not resolved.is_file():
        errors.append(f"{label}.path must be a regular non-symlink file")
        return
    data = resolved.read_bytes()
    if item.get("bytes") != len(data) or not HEX.fullmatch(str(item.get("sha256", ""))) or item["sha256"] != hashlib.sha256(data).hexdigest():
        errors.append(f"{label} bytes/SHA-256 mismatch")


def validate(record: object, root: Path, now: datetime | None = None) -> list[str]:
    errors: list[str] = []
    if not isinstance(record, dict):
        return ["record must be an object"]
    if set(record) != {"schema", "project_id", "environment", "release_digest", "evaluated_at", "controls"}:
        errors.append("record fields must be exact")
    if record.get("schema") != SCHEMA:
        errors.append(f"schema must be {SCHEMA}")
    if not isinstance(record.get("project_id"), str) or not record["project_id"].strip():
        errors.append("project_id is required")
    if record.get("environment") != "production":
        errors.append("environment must be production")
    if not DIGEST.fullmatch(str(record.get("release_digest", ""))) or set(str(record.get("release_digest", "")).split(":")[-1]) == {"0"}:
        errors.append("release_digest must be a real sha256 digest")
    evaluated = utc(record.get("evaluated_at"), "evaluated_at", errors)
    clock = now or datetime.now(timezone.utc)
    if evaluated and (evaluated > clock or (clock - evaluated).total_seconds() > 3600):
        errors.append("evaluated_at must be within the last hour")
    controls = record.get("controls")
    if not isinstance(controls, list):
        return sorted(set(errors + ["controls must be an array"]))
    ids = [item.get("id") for item in controls if isinstance(item, dict)]
    if len(controls) != 8 or set(ids) != CONTROLS or len(ids) != len(set(ids)):
        errors.append("controls must contain the eight unique production controls")
    for index, item in enumerate(controls):
        label = f"controls[{index}]"
        if not isinstance(item, dict):
            errors.append(f"{label} must be an object")
            continue
        required = {"id", "project_id", "environment", "release_digest", "result", "executed_at", "expires_at", "executor_ref", "tool", "target", "assertions", "evidence"}
        if set(item) != required:
            errors.append(f"{label} fields must be exact")
        for field in ("project_id", "environment", "release_digest"):
            if item.get(field) != record.get(field):
                errors.append(f"{label}.{field} must match admission record")
        if item.get("result") != "PASS":
            errors.append(f"{label}.result must be PASS")
        start = utc(item.get("executed_at"), f"{label}.executed_at", errors)
        end = utc(item.get("expires_at"), f"{label}.expires_at", errors)
        if start and end and not (start <= clock <= end and (end - start).total_seconds() <= 30 * 86400):
            errors.append(f"{label} evidence is stale, future or valid for more than 30 days")
        if not isinstance(item.get("executor_ref"), str) or not item["executor_ref"].strip():
            errors.append(f"{label}.executor_ref is required")
        tool = item.get("tool")
        if not isinstance(tool, dict) or set(tool) != {"name", "version", "source", "digest"} or not all(isinstance(tool.get(k), str) and tool[k].strip() for k in tool):
            errors.append(f"{label}.tool identity is incomplete")
        elif not DIGEST.fullmatch(tool["digest"]):
            errors.append(f"{label}.tool.digest must be sha256")
        if not isinstance(item.get("target"), dict) or not item["target"]:
            errors.append(f"{label}.target must identify the real target")
        assertions = item.get("assertions")
        expected_assertions = REQUIRED_ASSERTIONS.get(str(item.get("id")), set())
        if not isinstance(assertions, dict) or set(assertions) != expected_assertions or any(value is not True for value in assertions.values()):
            errors.append(f"{label}.assertions must contain exactly the required true semantic assertions")
        files = item.get("evidence")
        if not isinstance(files, list) or not files:
            errors.append(f"{label}.evidence must be non-empty")
        else:
            for number, proof in enumerate(files):
                evidence(root, proof, f"{label}.evidence[{number}]", errors)
    return sorted(set(errors))


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--project-root", required=True, type=Path)
    parser.add_argument("--record", default="PRODUCTION_ADMISSION.json")
    parser.add_argument("--report", default="PRODUCTION_ADMISSION_REPORT.json")
    args = parser.parse_args()
    root = args.project_root.resolve(strict=True)
    record_path = (root / args.record).resolve(strict=True)
    if record_path.parent != root or record_path.is_symlink():
        raise SystemExit("record must be a regular file directly under project root")
    record = json.loads(record_path.read_text(encoding="utf-8"))
    errors = validate(record, root)
    report = {"schema": "elite-production-admission-report/v1", "status": "PRODUCTION_ADMITTED" if not errors else "BLOCKED", "errors": errors}
    output = root / args.report
    if output.exists():
        raise SystemExit("report already exists")
    output.write_text(json.dumps(report, indent=2, sort_keys=True) + "\n", encoding="utf-8", newline="\n")
    print(json.dumps(report, indent=2, sort_keys=True))
    return 0 if not errors else 2


if __name__ == "__main__":
    sys.exit(main())
````

### FILE: `production_admission_gate/test_validate_production_admission.py`

```yaml
block_id: "SECURE-OPS:production-admission-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://production-admission/tests"
license: "LicenseRef-Workspace-Owner"
sha256: "9f47069b683a0f6bef24ce537e74253d95cd3c56231f2982339b592d87d9c7d8"
variables: []
secrets_allowed: false
```

````python
from __future__ import annotations

import copy
from datetime import datetime, timedelta, timezone
import hashlib
import tempfile
from pathlib import Path
import unittest

from validate_production_admission import CONTROLS, REQUIRED_ASSERTIONS, validate


class ProductionAdmissionTests(unittest.TestCase):
    def setUp(self) -> None:
        self.temp = tempfile.TemporaryDirectory()
        self.root = Path(self.temp.name).resolve()
        proof = self.root / "evidence.json"
        proof.write_bytes(b'{"pass":true}\n')
        now = datetime(2026, 8, 29, 18, tzinfo=timezone.utc)
        self.now = now
        stamp = lambda value: value.isoformat().replace("+00:00", "Z")
        evidence = {"path": "evidence.json", "bytes": proof.stat().st_size, "sha256": hashlib.sha256(proof.read_bytes()).hexdigest()}
        digest = "sha256:" + "a" * 64
        self.record = {
            "schema": "elite-production-admission/v1", "project_id": "revestex", "environment": "production",
            "release_digest": digest, "evaluated_at": stamp(now), "controls": [
                {"id": control, "project_id": "revestex", "environment": "production", "release_digest": digest,
                 "result": "PASS", "executed_at": stamp(now - timedelta(minutes=5)), "expires_at": stamp(now + timedelta(days=7)),
                 "executor_ref": "workload-identity://ci/production", "tool": {"name": "official-tool", "version": "1.0.0", "source": "https://example.invalid/official", "digest": "sha256:" + "b" * 64},
                 "target": {"id": f"production:{control}"}, "assertions": {key: True for key in sorted(REQUIRED_ASSERTIONS[control])}, "evidence": [evidence]}
                for control in sorted(CONTROLS)
            ]}

    def tearDown(self) -> None:
        self.temp.cleanup()

    def test_complete_fresh_hash_bound_record_passes(self) -> None:
        self.assertEqual([], validate(self.record, self.root, self.now))

    def test_missing_control_fails(self) -> None:
        value = copy.deepcopy(self.record); value["controls"].pop()
        self.assertTrue(any("eight unique" in error for error in validate(value, self.root, self.now)))

    def test_release_mismatch_fails(self) -> None:
        value = copy.deepcopy(self.record); value["controls"][0]["release_digest"] = "sha256:" + "c" * 64
        self.assertTrue(any("release_digest must match" in error for error in validate(value, self.root, self.now)))

    def test_stale_receipt_fails(self) -> None:
        value = copy.deepcopy(self.record); value["controls"][0]["expires_at"] = "2026-08-29T17:00:00Z"
        self.assertTrue(any("stale" in error for error in validate(value, self.root, self.now)))

    def test_tampered_evidence_fails(self) -> None:
        value = copy.deepcopy(self.record); (self.root / "evidence.json").write_text("tampered", encoding="utf-8")
        self.assertTrue(any("SHA-256 mismatch" in error for error in validate(value, self.root, self.now)))

    def test_false_assertion_fails(self) -> None:
        value = copy.deepcopy(self.record); key = next(iter(value["controls"][0]["assertions"])); value["controls"][0]["assertions"][key] = False
        self.assertTrue(any("semantic assertions" in error for error in validate(value, self.root, self.now)))

    def test_generic_exit_zero_assertion_cannot_replace_semantics(self) -> None:
        value = copy.deepcopy(self.record); value["controls"][0]["assertions"] = {"official_process_exit_zero": True}
        self.assertTrue(any("semantic assertions" in error for error in validate(value, self.root, self.now)))

    def test_duplicate_control_fails(self) -> None:
        value = copy.deepcopy(self.record); value["controls"][1]["id"] = value["controls"][0]["id"]
        self.assertTrue(any("eight unique" in error for error in validate(value, self.root, self.now)))

    def test_symlink_evidence_fails_when_supported(self) -> None:
        link = self.root / "link.json"
        try:
            link.symlink_to(self.root / "evidence.json")
        except OSError:
            self.skipTest("symlinks unavailable")
        value = copy.deepcopy(self.record); value["controls"][0]["evidence"][0]["path"] = "link.json"
        self.assertTrue(any("non-symlink" in error for error in validate(value, self.root, self.now)))


if __name__ == "__main__":
    unittest.main()
````

### FILE: `production_admission_gate/README.md`

```yaml
block_id: "SECURE-OPS:production-admission-readme:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://production-admission/readme"
license: "LicenseRef-Workspace-Owner"
sha256: "86edd9f61e6fab952bdc9e99f8e461b6450412e2b46e9689fdb2c0eed4020dfb"
variables: []
secrets_allowed: false
```

````markdown
# Executable production admission gate

This is Elite-authored orchestration around evidence produced by selected official tools. It does not claim that Stripe, Mercado Pago, Amazon, Google, Meta, TikTok, Mercado Libre, Cloudflare, OpenID Foundation, PostgreSQL, Grafana, ZAP, Kubernetes or GitHub authored these validators.

`PRODUCTION_ADMITTED` requires exactly eight fresh receipts bound to the same project, production environment and immutable release digest: edge/CDN/WAF, identity/authorization, selected providers, PostgreSQL restore/PITR, load/resilience, authorized offensive security, deploy/canary/rollback, and business acceptance. Every receipt identifies its tool/version/source/digest, real target, executor reference, true assertions and local evidence bytes+SHA-256. Missing, stale, duplicated, mismatched, symlinked or tampered evidence fails closed.

The receipt producer remains provider-specific. Recommended authorities are Cloudflare Rulesets API, OpenID Foundation Conformance Suite, PostgreSQL `pg_verifybackup` plus a real restore, Grafana k6 thresholds, ZAP Automation Framework, Kubernetes rollout/rollback and GitHub deployment environments where plan capabilities permit. Their actual outputs—not a narrative—must be stored and hashed by the project.

`run_official_tool.py` is the shared execution boundary. It requires explicit authorization, verifies the selected binary SHA-256, uses an argv array with `shell=False`, applies a timeout, passes only named environment variables plus minimal OS variables, captures bounded stdout/stderr and emits a hash-bound execution receipt. Capture uses fair nonblocking pipe reads (maximum64KiB per stream per pass), at most10MiB retained per stream and a monotonic deadline through both EOFs and child exit. Excess output stops the direct child before a receipt is written. Stdin is DEVNULL; Windows execution is hidden. Read/timeout/cleanup errors use fixed messages without private argv/output causes. Cleanup closes raw pipes and waits up to5seconds for the direct child after kill. This timeout cannot interrupt OS process creation, guarantee scheduling latency, or kill descendants; there are no reader threads, restart attempts or unbounded communicate calls. A descendant holding a write end cannot extend capture beyond its deadline; it still requires an external job/cgroup owner for termination. This is finite tool execution, not a durable service manager, heartbeat monitor, collector or sandbox. Python3.14.4/Windows is locally verified; other OS/runtime targets need their own execution evidence.

Successful tool output is intentionally preserved verbatim in evidence files. Access control, redaction review, storage quota and retention for those files remain target conditions; this patch does not classify arbitrary stdout as safe telemetry.

That receipt is not a production control receipt: a provider-specific adapter must still validate the exact semantic assertions required by `validate_production_admission.py`.

For `LOAD_RESILIENCE`, use the exact Grafana k6 2.2.0 identity in the official-source lock. Run `BASELINE`, `SOAK` and `SATURATION_RECOVERY` through the shared boundary, preserving every threshold exit, canonical argv hash, real target and workload-script hash. `validate_k6_load_admission.py` then verifies the approved error budget, all three runs, k6 v2 summaries and metrics, and emits the semantic receipt consumed by the final gate. One isolated k6 PASS is never promoted.

For `POSTGRES_RECOVERY`, use PostgreSQL 18.6 tools from an explicitly selected distribution and record each executable hash. The semantic adapter requires five bound runs: `pg_verifybackup`, restore start, restored-cluster readiness, application data invariants and PITR verification. It also requires WAL identity, an approved recovery policy and measured RPO/RTO. PostgreSQL explicitly states that `pg_verifybackup` does not replace a test restore; a logical dump/restore or zero exit alone is never promoted.

For `OFFENSIVE_SECURITY`, use the exact OWASP ZAP 2.17.0 identity and obtain explicit authorization for the target before active scanning. Run separate authenticated web and API Automation Framework plans through the shared boundary. The ZAP adapter binds plans, argv, reports and targets; rejects unaccepted high/critical alerts; and requires a two-person manual review covering authorization, authenticated scopes and business logic. A zero-alert report does not prove absence of vulnerabilities.

For `EDGE_CDN_WAF`, acquire the exact official Cloudflare `cloudflare-go` 7.9.0 source and build a project-owned probe whose binary is hash-bound by the shared runner. Four runs are mandatory: Cloudflare API configuration snapshot, direct-origin bypass negative, two-request public-cache behavior and authenticated private-cache bypass. `validate_cloudflare_edge_admission.py` requires the same binary, argv, zone, hostname and observation hashes; proves proxied DNS, managed/custom WAF, Full (strict) TLS, origin protection and cache separation; and emits the semantic receipt consumed by the final gate. The library supplies no account, token, zone, origin or production observation.

For `IDENTITY_AUTHORIZATION`, acquire OpenID Foundation Conformance Suite 5.2.4 and Microsoft Playwright 1.62.1 from the exact source lock. Five ordered runs bind one real issuer, client, audience, conformance profile, application, immutable release, declared roles and at least two tenants. `validate_identity_authorization_admission.py` rejects failed/interrupted conformance results, unreviewed non-passed cases, missing role effects, unauthorized disclosure/mutation, an incomplete ordered cross-tenant matrix, accepted revoked/expired/replayed credentials and incomplete break-glass governance. The Foundation suite is free to run, while formal certification is separate and may cost money. Source or synthetic fixture availability never admits a production IdP.

For `PROVIDERS`, select only the operations actually required by the project and bind them to the exact admitted Stripe, Mercado Pago, Amazon SP-API, Google Ads, Meta Ads, TikTok Ads, Meta WhatsApp, Firebase or Google Merchant source identity. Mercado Libre is explicitly different: its active path is Elite-authored against current official HTTP documentation because its archived SDK is not represented as reusable official code. Every selected provider needs four ordered real executions: sandbox/test-account contract for every operation, webhook authenticity and durable replay/reconciliation or a two-person-approved non-applicability decision, idempotency/deduplication plus uncertain-outcome reconciliation, and bounded rate-limit/retry behavior. `validate_provider_admission.py` also requires current terms, least scopes, jurisdiction/data region/retention, independent finance/security/reconciliation ownership and approved quota/cost. Fixtures and SDK presence never admit an account or external effect.

For `DEPLOY_ROLLBACK`, acquire the exact official Kubernetes v1.37.0 source and kubectl artifact from the locked profile. Seven ordered target-bound runs are mandatory: cluster/context/authorization preflight, apply of an immutable image digest, canary health, completed Deployment rollout, an actually executed rollback, verification that the prior digest and service/data health were restored, and forward/backward migration compatibility. `validate_deploy_rollback_admission.py` fixes the published Windows amd64 kubectl SHA-256, restricts semantics to Kubernetes `Deployment`, requires two-person policy approval and rejects mutable tags, target drift, incomplete replicas, old replicas, SLO breach, unsafe migration or a rollback that did not restore the previous digest. A local version command, dry run, synthetic observation or successful `kubectl rollout status` cannot admit production.

For `BUSINESS_ACCEPTANCE`, use GitHub Spec Kit 1.0.1 to preserve the approved constitution, specification, plan, tasks and convergence artifacts, then execute every approved capability scenario with Microsoft Playwright 1.62.1 against the selected production target. `validate_business_acceptance_admission.py` requires unique role/tenant-bound scenarios, observed expected outcomes and backend effects, explicit financial-effect disposition, hash-bound Playwright execution, independent finance/operations/security approvals, no accepted-open Sev1/Sev2 defect and a distinct business-owner signoff over the exact specification, observation, approval and defect hashes. Subject hashes are evidence references, not invented digital signatures. Fixtures, generated specs, a converged task list or browser exit zero cannot substitute for real target results and human authority.

```powershell
python production_admission_gate/validate_production_admission.py --project-root .
```

Exit 0 and `status=PRODUCTION_ADMITTED` are necessary but do not override a missing provider-specific semantic assertion. The project must define those assertions before executing each official tool; a generic `required_gate_executed` fixture exists only in tests and is not a production profile.
````

### FILE: `production_admission_gate/official-tool-run.template.json`

```yaml
block_id: "SECURE-OPS:official-tool-run-template:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://production-admission/official-tool-run-template"
license: "LicenseRef-Workspace-Owner"
sha256: "464108e56fa28dea2ca47512dedd039a90303d575f16ae1c3cfdaa8085c92885"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite-official-tool-run/v1",
  "project_id": "",
  "environment": "production",
  "release_digest": "",
  "control_id": "",
  "tool": {"path":"", "sha256":"", "name":"", "version":"", "source":""},
  "arguments": [],
  "working_directory": ".",
  "environment_variable_names": [],
  "timeout_seconds": 300,
  "target": {}
}
````

### FILE: `production_admission_gate/run_official_tool.py`

```yaml
block_id: "SECURE-OPS:official-tool-runner:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://production-admission/official-tool-runner"
license: "LicenseRef-Workspace-Owner"
sha256: "41e73b0245c20f5056caf8ae818f0ce61dc502fedf4b28caa245cd3b0f09c5e9"
variables: []
secrets_allowed: false
```

````python
from __future__ import annotations

import argparse
from datetime import datetime, timezone
import hashlib
import json
import os
from pathlib import Path, PurePosixPath
import re
import subprocess
import sys
import tempfile
import time

from validate_production_admission import CONTROLS, DIGEST, HEX

SCHEMA = "elite-official-tool-run/v1"
ENV_NAME = re.compile(r"^[A-Z][A-Z0-9_]{0,63}$")
MAX_OUTPUT = 10 * 1024 * 1024


def safe_project_path(root: Path, value: object, label: str, *, file: bool) -> Path:
    if not isinstance(value, str) or "\\" in value:
        raise ValueError(f"{label} must be canonical project-relative")
    relative = PurePosixPath(value)
    if relative.is_absolute() or any(part in {"", ".."} for part in relative.parts):
        raise ValueError(f"{label} must be canonical project-relative")
    candidate = root.joinpath(*relative.parts)
    resolved = candidate.resolve(strict=True)
    resolved.relative_to(root)
    if candidate.is_symlink() or (file and not resolved.is_file()) or (not file and not resolved.is_dir()):
        raise ValueError(f"{label} has invalid type or symlink")
    return resolved


def atomic(path: Path, data: bytes) -> None:
    if path.exists():
        raise ValueError(f"output already exists: {path.name}")
    with tempfile.NamedTemporaryFile(dir=path.parent, prefix=f".{path.name}.", delete=False) as handle:
        temporary = Path(handle.name); handle.write(data); handle.flush(); os.fsync(handle.fileno())
    try:
        temporary.replace(path)
    except BaseException:
        temporary.unlink(missing_ok=True); raise


def bounded_process(arguments: list[str], cwd: Path, environment: dict[str, str], timeout: int) -> subprocess.CompletedProcess:
    """Bound each captured stream; supervise only the direct child, never restart.

    The deadline includes pipe EOF after child exit. No communicate(), reader
    threads, inherited stdin or output retry. OS process creation itself is not
    interruptible here. Descendants require a target-owned job/cgroup supervisor.
    """
    process = None
    streams = []
    buffers = [bytearray(), bytearray()]
    failure = None
    completed = None
    deadline = time.monotonic() + timeout
    try:
        flags = subprocess.CREATE_NO_WINDOW if os.name == "nt" else 0
        process = subprocess.Popen(arguments, cwd=cwd, env=environment, shell=False,
                                   stdin=subprocess.DEVNULL, stdout=subprocess.PIPE,
                                   stderr=subprocess.PIPE, bufsize=0, creationflags=flags)
        streams = [process.stdout, process.stderr]
        for stream in streams:
            os.set_blocking(stream.fileno(), False)
        pending = {0, 1}
        while True:
            if time.monotonic() >= deadline:
                failure = "official tool timed out"
                break
            progressed = False
            for index in tuple(pending):
                # At most one bounded read per stream per pass prevents starvation.
                size = min(65536, MAX_OUTPUT - len(buffers[index]) + 1)
                try:
                    chunk = os.read(streams[index].fileno(), size)
                except BlockingIOError:
                    continue
                if not chunk:
                    pending.remove(index)
                    continue
                progressed = True
                if len(buffers[index]) + len(chunk) > MAX_OUTPUT:
                    failure = "official tool output exceeded 10 MiB"
                    break
                buffers[index].extend(chunk)
            if failure:
                break
            code = process.poll()
            if code is not None and not pending:
                completed = subprocess.CompletedProcess(arguments, code, bytes(buffers[0]), bytes(buffers[1]))
                break
            if not progressed:
                time.sleep(min(.005, max(0, deadline - time.monotonic())))
    except (OSError, ValueError, subprocess.SubprocessError):
        failure = "official tool execution unavailable"
    finally:
        # Raw unbuffered handles: closing cannot wait on a buffered reader lock.
        # Close the pipes even when a descendant keeps its inherited write end.
        for stream in streams:
            try:
                stream.close()
            except OSError:
                failure = "official tool cleanup unavailable"
        if process is not None:
            try:
                if process.poll() is None:
                    process.kill()
                process.wait(timeout=5)
            except (OSError, subprocess.SubprocessError):
                failure = "official tool cleanup unavailable"
    if failure:
        # Do not chain TimeoutExpired/OSError: argv and tool output may be private.
        raise ValueError(failure) from None
    if completed is None:
        raise ValueError("official tool execution unavailable") from None
    return completed


def run(profile: object, root: Path, output: Path, authorized: bool) -> dict[str, object]:
    if not authorized:
        raise ValueError("explicit --authorize-execution is required")
    if not isinstance(profile, dict) or set(profile) != {"schema", "project_id", "environment", "release_digest", "control_id", "tool", "arguments", "working_directory", "environment_variable_names", "timeout_seconds", "target"}:
        raise ValueError("profile fields must be exact")
    if profile["schema"] != SCHEMA or profile["environment"] != "production" or profile["control_id"] not in CONTROLS:
        raise ValueError("profile identity is invalid")
    if not isinstance(profile["project_id"], str) or not profile["project_id"].strip() or not DIGEST.fullmatch(str(profile["release_digest"])):
        raise ValueError("project/release identity is invalid")
    tool = profile["tool"]
    if not isinstance(tool, dict) or set(tool) != {"path", "sha256", "name", "version", "source"} or not HEX.fullmatch(str(tool.get("sha256", ""))):
        raise ValueError("tool identity is invalid")
    binary = Path(str(tool["path"])).resolve(strict=True)
    if Path(str(tool["path"])).is_symlink() or not binary.is_file() or hashlib.sha256(binary.read_bytes()).hexdigest() != tool["sha256"]:
        raise ValueError("tool binary SHA-256 mismatch or invalid file")
    arguments = profile["arguments"]
    if not isinstance(arguments, list) or len(arguments) > 128 or any(not isinstance(v, str) or "\x00" in v or len(v) > 4096 for v in arguments):
        raise ValueError("arguments must be a bounded argv array")
    cwd = safe_project_path(root, profile["working_directory"], "working_directory", file=False)
    names = profile["environment_variable_names"]
    if not isinstance(names, list) or len(names) != len(set(names)) or any(not isinstance(v, str) or not ENV_NAME.fullmatch(v) for v in names):
        raise ValueError("environment_variable_names are invalid")
    missing = [name for name in names if name not in os.environ]
    if missing:
        raise ValueError(f"required environment variables are absent: {missing}")
    timeout = profile["timeout_seconds"]
    if type(timeout) is not int or not 1 <= timeout <= 7200:
        raise ValueError("timeout_seconds must be 1..7200")
    base_names = {"PATH", "SYSTEMROOT", "WINDIR", "TEMP", "TMP", "LANG", "LC_ALL", "SSL_CERT_FILE", "SSL_CERT_DIR"}
    environment = {name: value for name, value in os.environ.items() if name in base_names or name in names}
    started = datetime.now(timezone.utc)
    completed = bounded_process([str(binary), *arguments], cwd, environment, timeout)
    stdout = output.with_suffix(".stdout.bin"); stderr = output.with_suffix(".stderr.bin")
    atomic(stdout, completed.stdout); atomic(stderr, completed.stderr)
    proof = lambda path, data: {"path": path.relative_to(root).as_posix(), "bytes": len(data), "sha256": hashlib.sha256(data).hexdigest()}
    arguments_sha256 = hashlib.sha256(json.dumps(arguments, ensure_ascii=False, separators=(",", ":")).encode()).hexdigest()
    receipt = {"schema": "elite-official-tool-execution/v1", "project_id": profile["project_id"], "environment": "production", "release_digest": profile["release_digest"], "control_id": profile["control_id"], "executed_at": started.isoformat().replace("+00:00", "Z"), "exit_code": completed.returncode, "tool": {key: tool[key] for key in ("name", "version", "source", "sha256")}, "arguments_sha256": arguments_sha256, "working_directory": Path(profile["working_directory"]).as_posix(), "target": profile["target"], "environment_variable_names": names, "stdout": proof(stdout, completed.stdout), "stderr": proof(stderr, completed.stderr)}
    atomic(output, (json.dumps(receipt, indent=2, sort_keys=True) + "\n").encode())
    return receipt


def main() -> int:
    parser = argparse.ArgumentParser(); parser.add_argument("--project-root", required=True, type=Path); parser.add_argument("--profile", required=True); parser.add_argument("--output", required=True); parser.add_argument("--authorize-execution", action="store_true"); args = parser.parse_args()
    root = args.project_root.resolve(strict=True); profile_path = safe_project_path(root, args.profile, "profile", file=True)
    output = root.joinpath(*PurePosixPath(args.output).parts).resolve(strict=False); output.parent.resolve(strict=True).relative_to(root)
    receipt = run(json.loads(profile_path.read_text(encoding="utf-8")), root, output, args.authorize_execution)
    print(json.dumps(receipt, indent=2, sort_keys=True)); return 0 if receipt["exit_code"] == 0 else 2


if __name__ == "__main__": sys.exit(main())
````

### FILE: `production_admission_gate/test_run_official_tool.py`

```yaml
block_id: "SECURE-OPS:official-tool-runner-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://production-admission/official-tool-runner-tests"
license: "LicenseRef-Workspace-Owner"
sha256: "71d20612e8cd3150db8d064675cccba6044f784092e19a59b296ff385534a6ce"
variables: []
secrets_allowed: false
```

````python
from __future__ import annotations

import copy
import hashlib
import json
import os
from pathlib import Path
import sys
import tempfile
import unittest
import time
import traceback
import threading
from unittest.mock import patch
import run_official_tool as runner

from run_official_tool import run


class OfficialToolRunnerTests(unittest.TestCase):
    def setUp(self) -> None:
        self.temp = tempfile.TemporaryDirectory(); self.root = Path(self.temp.name).resolve(); (self.root / "work").mkdir()
        binary = Path(sys.executable).resolve()
        self.profile = {"schema":"elite-official-tool-run/v1", "project_id":"revestex", "environment":"production", "release_digest":"sha256:" + "a" * 64, "control_id":"LOAD_RESILIENCE", "tool":{"path":str(binary), "sha256":hashlib.sha256(binary.read_bytes()).hexdigest(), "name":"CPython fixture", "version":sys.version.split()[0], "source":"https://python.org"}, "arguments":["-I", "-c", "print('official-output')"], "working_directory":"work", "environment_variable_names":[], "timeout_seconds":10, "target":{"id":"production:test"}}

    def tearDown(self) -> None: self.temp.cleanup()

    def test_exact_binary_runs_without_shell_and_hashes_outputs(self) -> None:
        receipt = run(self.profile, self.root, self.root / "receipt.json", True)
        self.assertEqual(0, receipt["exit_code"]); self.assertEqual(b"official-output\r\n" if os.name == "nt" else b"official-output\n", (self.root / "receipt.stdout.bin").read_bytes())
        self.assertEqual(hashlib.sha256((self.root / "receipt.stdout.bin").read_bytes()).hexdigest(), receipt["stdout"]["sha256"])
        expected = hashlib.sha256(json.dumps(self.profile["arguments"], ensure_ascii=False, separators=(",", ":")).encode()).hexdigest()
        self.assertEqual(expected, receipt["arguments_sha256"]); self.assertEqual("work", receipt["working_directory"])

    def test_explicit_authorization_is_required(self) -> None:
        with self.assertRaisesRegex(ValueError, "authorize"): run(self.profile, self.root, self.root / "receipt.json", False)

    def test_binary_hash_mismatch_fails_before_execution(self) -> None:
        value = copy.deepcopy(self.profile); value["tool"]["sha256"] = "0" * 64
        with self.assertRaisesRegex(ValueError, "SHA-256"): run(value, self.root, self.root / "receipt.json", True)

    def test_missing_environment_reference_fails_without_value_leak(self) -> None:
        value = copy.deepcopy(self.profile); value["environment_variable_names"] = ["ELITE_MISSING_TEST_SECRET"]
        with self.assertRaisesRegex(ValueError, "ELITE_MISSING_TEST_SECRET"): run(value, self.root, self.root / "receipt.json", True)

    def test_nonzero_exit_is_preserved_not_promoted(self) -> None:
        value = copy.deepcopy(self.profile); value["arguments"] = ["-I", "-c", "raise SystemExit(7)"]
        receipt = run(value, self.root, self.root / "receipt.json", True); self.assertEqual(7, receipt["exit_code"])


    def assert_no_receipt(self):
        self.assertEqual([], list(self.root.glob("receipt*")))

    def test_excess_is_stopped_before_later_effect(self):
        self.profile["arguments"] = ["-I", "-c", "import os,time,pathlib;os.write(1,b'x'*65536);time.sleep(.3);pathlib.Path('continued').write_text('bad')"]
        with patch.object(runner, "MAX_OUTPUT", 1024), self.assertRaisesRegex(ValueError, "output exceeded"):
            run(self.profile, self.root, self.root / "receipt.json", True)
        self.assertFalse((self.root / "work/continued").exists())
        self.assert_no_receipt()

    def test_timeout_has_no_private_argument_and_child_is_reaped(self):
        self.profile["timeout_seconds"] = 1
        self.profile["arguments"] = ["-I", "-c", "import time;time.sleep(4)", "PRIVATE_MARKER_345"]
        children = []
        popen = runner.subprocess.Popen
        def launch(*args, **kwargs):
            child = popen(*args, **kwargs); children.append(child); return child
        before = {t.ident for t in threading.enumerate()}
        start = time.monotonic()
        with patch.object(runner.subprocess, "Popen", side_effect=launch):
            try: run(self.profile, self.root, self.root / "receipt.json", True)
            except ValueError as error:
                self.assertEqual("official tool timed out", str(error))
                self.assertNotIn("PRIVATE_MARKER_345", traceback.format_exc())
            else: self.fail("timeout accepted")
        self.assertLess(time.monotonic() - start, 3)
        self.assertEqual(1, len(children))
        self.assertIsNotNone(children[0].poll())
        self.assertTrue(children[0].stdout.closed and children[0].stderr.closed)
        self.assertEqual(before, {t.ident for t in threading.enumerate()})
        self.assert_no_receipt()

    def test_both_streams_at_exact_limit_preserve_binary_content(self):
        limit = 1024 * 1024
        self.profile["arguments"] = ["-I", "-c", "import os;[(os.write(1,bytes(range(256))*16),os.write(2,bytes(reversed(range(256)))*16)) for _ in range(256)]"]
        with patch.object(runner, "MAX_OUTPUT", limit):
            receipt = run(self.profile, self.root, self.root / "receipt.json", True)
        self.assertEqual(0, receipt["exit_code"])
        for kind, data in [("stdout", bytes(range(256))*4096), ("stderr", bytes(reversed(range(256)))*4096)]:
            self.assertEqual(data, (self.root / receipt[kind]["path"]).read_bytes())
            self.assertEqual(limit, receipt[kind]["bytes"])
            self.assertEqual(hashlib.sha256(data).hexdigest(), receipt[kind]["sha256"])

    def test_each_stream_one_byte_over_limit_rejected(self):
        for fd in (1, 2):
            with self.subTest(fd=fd):
                self.profile["arguments"] = ["-I", "-c", f"import os;os.write({fd},b'x'*1025)"]
                with patch.object(runner, "MAX_OUTPUT", 1024), self.assertRaisesRegex(ValueError, "output exceeded"):
                    run(self.profile, self.root, self.root / "receipt.json", True)
                self.assert_no_receipt()

    def test_closed_streams_do_not_make_live_child_successful(self):
        self.profile["timeout_seconds"] = 1
        self.profile["arguments"] = ["-I", "-c", "import os,time;os.close(1);os.close(2);time.sleep(4)"]
        with self.assertRaisesRegex(ValueError, "timed out"):
            run(self.profile, self.root, self.root / "receipt.json", True)
        self.assert_no_receipt()

    def test_inherited_pipe_has_deadline_without_waiting_for_descendant(self):
        # The descendant is finite and exits itself. This explicitly does NOT
        # claim tree termination; a service manager must supply that boundary.
        self.profile["timeout_seconds"] = 1
        child_code = "import time,pathlib;time.sleep(2.5);pathlib.Path('descendant-finished').write_text('done')"
        self.profile["arguments"] = ["-I", "-c", f"import subprocess,sys;subprocess.Popen([sys.executable,'-I','-c',{child_code!r}])"]
        start = time.monotonic()
        try:
            with self.assertRaisesRegex(ValueError, "timed out"):
                run(self.profile, self.root, self.root / "receipt.json", True)
            self.assertLess(time.monotonic() - start, 2)
            self.assert_no_receipt()
        finally:
            until = time.monotonic() + 4
            while not (self.root / "work/descendant-finished").exists() and time.monotonic() < until:
                time.sleep(.02)
            self.assertTrue((self.root / "work/descendant-finished").is_file())

    def test_stdin_is_closed_and_unselected_environment_is_absent(self):
        self.profile["arguments"] = ["-I", "-c", "import sys,os;assert sys.stdin.buffer.read()==b'';assert 'ELITE_UNSELECTED_SECRET' not in os.environ;print('closed')"]
        with patch.dict(os.environ, {"ELITE_UNSELECTED_SECRET": "PRIVATE_MARKER_345"}):
            receipt = run(self.profile, self.root, self.root / "receipt.json", True)
        self.assertEqual(0, receipt["exit_code"])

    def test_boolean_timeout_rejected_before_process_creation(self):
        self.profile["timeout_seconds"] = True
        with patch.object(runner.subprocess, "Popen") as launch, self.assertRaisesRegex(ValueError, "timeout_seconds"):
            run(self.profile, self.root, self.root / "receipt.json", True)
        launch.assert_not_called()

    def test_read_error_is_static_and_cleans_up_direct_child(self):
        self.profile["arguments"] = ["-I", "-c", "import time;time.sleep(4)"]
        children = []; popen = runner.subprocess.Popen; read = runner.os.read
        def launch(*args, **kwargs):
            child = popen(*args, **kwargs); children.append(child); return child
        def fail_read(fd, count):
            if children and fd in (children[0].stdout.fileno(), children[0].stderr.fileno()):
                raise OSError("PRIVATE_MARKER_345")
            return read(fd, count)
        with patch.object(runner.subprocess, "Popen", side_effect=launch), patch.object(runner.os, "read", side_effect=fail_read):
            try: run(self.profile, self.root, self.root / "receipt.json", True)
            except ValueError as error:
                self.assertEqual("official tool execution unavailable", str(error))
                self.assertNotIn("PRIVATE_MARKER_345", traceback.format_exc())
            else: self.fail("read failure accepted")
        self.assertIsNotNone(children[0].poll())
        self.assert_no_receipt()



if __name__ == "__main__": unittest.main()
````

### FILE: `production_admission_gate/k6-load-admission.template.json`

```yaml
block_id: "SECURE-OPS:k6-load-admission-template:v1"
operation: CREATE
provenance: AUTHORED
source: "local semantic adapter over exact official Grafana k6 2.2.0 execution evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "c3f66aeba0d003c6843aafab3de71c9dd7630bf2e619fe1a1f97f96af7e3a744"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite-k6-load-admission/v1",
  "project_id": "",
  "environment": "production",
  "release_digest": "",
  "evaluated_at": "",
  "expires_at": "",
  "executor_ref": "",
  "tool": {
    "name": "Grafana k6",
    "version": "2.2.0",
    "source": "https://github.com/grafana/k6/releases/tag/v2.2.0",
    "sha256": "87dfa91bc3e47bc4bd77911d59d7ff79f25cd76fb8322c97072f08f08a0da5ed"
  },
  "target": {},
  "error_budget": "",
  "runs": [
    {"kind":"BASELINE", "execution_receipt":"", "summary":"", "script":"", "expected_arguments":[]},
    {"kind":"SOAK", "execution_receipt":"", "summary":"", "script":"", "expected_arguments":[]},
    {"kind":"SATURATION_RECOVERY", "execution_receipt":"", "summary":"", "script":"", "expected_arguments":[]}
  ],
  "output": ""
}
````

### FILE: `production_admission_gate/validate_k6_load_admission.py`

```yaml
block_id: "SECURE-OPS:k6-load-admission-validator:v1"
operation: CREATE
provenance: AUTHORED
source: "local semantic adapter over exact official Grafana k6 2.2.0 execution evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "801dc18eef947fc2291960a353e1dcef2fd03af3bf3cc05e14d6357681987201"
variables: []
secrets_allowed: false
```

````python
from __future__ import annotations

import argparse
from datetime import datetime, timezone
import hashlib
import json
from pathlib import Path, PurePosixPath
import re
import sys
import tempfile

DIGEST = re.compile(r"sha256:[0-9a-f]{64}")
SHA = re.compile(r"[0-9a-f]{64}")
KINDS = ("BASELINE", "SOAK", "SATURATION_RECOVERY")
PROFILE_FIELDS = {"schema", "project_id", "environment", "release_digest", "evaluated_at", "expires_at", "executor_ref", "tool", "target", "error_budget", "runs", "output"}
RUN_FIELDS = {"kind", "execution_receipt", "summary", "script", "expected_arguments"}
K6_SHA256 = "87dfa91bc3e47bc4bd77911d59d7ff79f25cd76fb8322c97072f08f08a0da5ed"


def canonical_sha(value: object) -> str:
    return hashlib.sha256(json.dumps(value, ensure_ascii=False, separators=(",", ":")).encode()).hexdigest()


def safe(root: Path, raw: object, label: str, *, exists: bool = True) -> Path:
    if not isinstance(raw, str) or not raw or "\\" in raw:
        raise ValueError(f"{label} must be a non-empty POSIX relative path")
    pure = PurePosixPath(raw)
    if pure.is_absolute() or any(part in ("", ".", "..") for part in pure.parts):
        raise ValueError(f"{label} is unsafe")
    path = root.joinpath(*pure.parts).resolve(strict=exists)
    path.relative_to(root)
    if exists and not path.is_file():
        raise ValueError(f"{label} must be a file")
    return path


def read_json(path: Path, label: str) -> dict[str, object]:
    try:
        value = json.loads(path.read_text(encoding="utf-8"))
    except (OSError, UnicodeError, json.JSONDecodeError) as error:
        raise ValueError(f"{label} is not valid UTF-8 JSON") from error
    if not isinstance(value, dict):
        raise ValueError(f"{label} must be an object")
    return value


def proof(root: Path, path: Path) -> dict[str, object]:
    data = path.read_bytes()
    return {"path": path.relative_to(root).as_posix(), "bytes": len(data), "sha256": hashlib.sha256(data).hexdigest()}


def verify_proof(root: Path, item: object, label: str) -> Path:
    if not isinstance(item, dict) or set(item) != {"path", "bytes", "sha256"}:
        raise ValueError(f"{label} proof fields must be exact")
    path = safe(root, item["path"], label)
    actual = proof(root, path)
    if actual != item:
        raise ValueError(f"{label} proof mismatch")
    return path


def utc(raw: object, label: str) -> datetime:
    if not isinstance(raw, str) or not raw.endswith("Z"):
        raise ValueError(f"{label} must be ISO-8601 UTC ending Z")
    try:
        return datetime.fromisoformat(raw[:-1] + "+00:00").astimezone(timezone.utc)
    except ValueError as error:
        raise ValueError(f"{label} must be ISO-8601 UTC ending Z") from error


def metric(summary: dict[str, object], name: str) -> dict[str, object]:
    results = summary.get("results")
    metrics = results.get("metrics") if isinstance(results, dict) else None
    if not isinstance(metrics, list):
        raise ValueError("k6 summary results.metrics must be a list")
    matches = [item for item in metrics if isinstance(item, dict) and item.get("name") == name]
    if len(matches) != 1 or not isinstance(matches[0].get("values"), dict):
        raise ValueError(f"k6 summary requires exactly one {name} metric")
    return matches[0]["values"]


def number(values: dict[str, object], key: str, label: str) -> float:
    value = values.get(key)
    if isinstance(value, bool) or not isinstance(value, (int, float)) or value < 0:
        raise ValueError(f"{label} must be a non-negative number")
    return float(value)


def atomic(path: Path, data: bytes) -> None:
    if path.exists():
        raise ValueError("output already exists")
    path.parent.mkdir(parents=False, exist_ok=True)
    with tempfile.NamedTemporaryFile(dir=path.parent, prefix=f".{path.name}.", delete=False) as handle:
        temporary = Path(handle.name); handle.write(data); handle.flush()
    try:
        temporary.replace(path)
    except Exception:
        temporary.unlink(missing_ok=True); raise


def validate(profile: object, root: Path) -> dict[str, object]:
    if not isinstance(profile, dict) or set(profile) != PROFILE_FIELDS:
        raise ValueError("profile fields must be exact")
    if profile["schema"] != "elite-k6-load-admission/v1" or profile["environment"] != "production":
        raise ValueError("profile identity is invalid")
    if not isinstance(profile["project_id"], str) or not profile["project_id"].strip() or not DIGEST.fullmatch(str(profile["release_digest"])):
        raise ValueError("project/release identity is invalid")
    start = utc(profile["evaluated_at"], "evaluated_at"); end = utc(profile["expires_at"], "expires_at")
    if not start < end or (end - start).total_seconds() > 30 * 86400:
        raise ValueError("admission validity must be positive and at most 30 days")
    if not isinstance(profile["executor_ref"], str) or not profile["executor_ref"].strip() or not isinstance(profile["target"], dict) or not profile["target"]:
        raise ValueError("executor_ref and real target are required")
    tool = profile["tool"]
    expected_tool = {"name":"Grafana k6", "version":"2.2.0", "source":"https://github.com/grafana/k6/releases/tag/v2.2.0", "sha256":K6_SHA256}
    if tool != expected_tool:
        raise ValueError("exact Grafana k6 2.2.0 tool identity is required")

    budget_path = safe(root, profile["error_budget"], "error_budget")
    budget = read_json(budget_path, "error_budget")
    budget_fields = {"schema", "project_id", "environment", "release_digest", "approved_at", "approvers", "runs"}
    if set(budget) != budget_fields or budget.get("schema") != "elite-load-error-budget/v1":
        raise ValueError("error budget fields/schema must be exact")
    for key in ("project_id", "environment", "release_digest"):
        if budget.get(key) != profile[key]: raise ValueError(f"error budget {key} mismatch")
    utc(budget["approved_at"], "error_budget.approved_at")
    approvers = budget["approvers"]
    if not isinstance(approvers, list) or len(approvers) < 2 or any(not isinstance(x, str) or not x.strip() for x in approvers) or len(set(approvers)) != len(approvers):
        raise ValueError("error budget requires at least two distinct approvers")
    limits = budget["runs"]
    if not isinstance(limits, dict) or set(limits) != set(KINDS):
        raise ValueError("error budget must define exactly three run kinds")

    runs = profile["runs"]
    if not isinstance(runs, list) or len(runs) != 3 or any(not isinstance(run, dict) or set(run) != RUN_FIELDS for run in runs):
        raise ValueError("runs must contain three exact records")
    if [run["kind"] for run in runs] != list(KINDS):
        raise ValueError("runs must be ordered BASELINE, SOAK, SATURATION_RECOVERY")
    evidence = [proof(root, budget_path)]
    for run in runs:
        kind = run["kind"]; limit = limits[kind]
        required_limits = {"minimum_iterations", "minimum_duration_seconds", "maximum_http_req_failed_rate", "maximum_http_req_duration_p95_ms"}
        if not isinstance(limit, dict) or set(limit) != required_limits:
            raise ValueError(f"{kind} budget fields must be exact")
        minimum_iterations = number(limit, "minimum_iterations", f"{kind}.minimum_iterations")
        minimum_duration = number(limit, "minimum_duration_seconds", f"{kind}.minimum_duration_seconds")
        maximum_failed = number(limit, "maximum_http_req_failed_rate", f"{kind}.maximum_http_req_failed_rate")
        maximum_p95 = number(limit, "maximum_http_req_duration_p95_ms", f"{kind}.maximum_http_req_duration_p95_ms")
        if minimum_iterations < 1 or minimum_duration < 1 or maximum_failed > 1 or maximum_p95 <= 0:
            raise ValueError(f"{kind} budget values are invalid")
        script_path = safe(root, run["script"], f"{kind}.script"); script_proof = proof(root, script_path)
        receipt_path = safe(root, run["execution_receipt"], f"{kind}.execution_receipt"); receipt = read_json(receipt_path, f"{kind}.execution_receipt")
        if receipt.get("schema") != "elite-official-tool-execution/v1" or receipt.get("project_id") != profile["project_id"] or receipt.get("environment") != "production" or receipt.get("release_digest") != profile["release_digest"] or receipt.get("control_id") != "LOAD_RESILIENCE" or receipt.get("exit_code") != 0:
            raise ValueError(f"{kind} execution identity/exit mismatch")
        receipt_tool = receipt.get("tool")
        if receipt_tool != {"name":tool["name"], "version":tool["version"], "source":tool["source"], "sha256":tool["sha256"]}:
            raise ValueError(f"{kind} tool identity mismatch")
        arguments = run["expected_arguments"]
        if not isinstance(arguments, list) or any(not isinstance(value, str) or not value for value in arguments) or receipt.get("arguments_sha256") != canonical_sha(arguments):
            raise ValueError(f"{kind} arguments are not hash-bound")
        expected_target = dict(profile["target"]); expected_target.update({"run_kind":kind, "workload_script_sha256":script_proof["sha256"], "minimum_duration_seconds":minimum_duration})
        if receipt.get("target") != expected_target:
            raise ValueError(f"{kind} target/script/duration binding mismatch")
        for channel in ("stdout", "stderr"):
            verify_proof(root, receipt.get(channel), f"{kind}.{channel}")
        summary_path = safe(root, run["summary"], f"{kind}.summary"); summary = read_json(summary_path, f"{kind}.summary")
        metadata = summary.get("metadata")
        if summary.get("version") != "1.0.0" or not isinstance(metadata, dict) or metadata.get("k6_version") != "2.2.0":
            raise ValueError(f"{kind} requires k6 v2 machine-readable summary")
        if number(metric(summary, "iterations"), "count", f"{kind}.iterations") < minimum_iterations:
            raise ValueError(f"{kind} iterations below approved minimum")
        if number(metric(summary, "http_req_failed"), "rate", f"{kind}.http_req_failed") > maximum_failed:
            raise ValueError(f"{kind} failure rate exceeds approved budget")
        if number(metric(summary, "http_req_duration"), "p(95)", f"{kind}.http_req_duration") > maximum_p95:
            raise ValueError(f"{kind} p95 exceeds approved budget")
        evidence.extend([proof(root, receipt_path), proof(root, summary_path), script_proof, receipt["stdout"], receipt["stderr"]])

    return {"id":"LOAD_RESILIENCE", "project_id":profile["project_id"], "environment":"production", "release_digest":profile["release_digest"], "result":"PASS", "executed_at":profile["evaluated_at"], "expires_at":profile["expires_at"], "executor_ref":profile["executor_ref"], "tool":{"name":tool["name"], "version":tool["version"], "source":tool["source"], "digest":"sha256:" + tool["sha256"]}, "target":profile["target"], "assertions":{"k6_thresholds_pass":True, "expected_workload_represented":True, "soak_pass":True, "saturation_recovery_pass":True, "no_error_budget_breach":True}, "evidence":evidence}


def main() -> int:
    parser = argparse.ArgumentParser(); parser.add_argument("--project-root", required=True, type=Path); parser.add_argument("--profile", required=True); args = parser.parse_args()
    try:
        root = args.project_root.resolve(strict=True); profile_path = safe(root, args.profile, "profile"); profile = read_json(profile_path, "profile")
        output = safe(root, profile.get("output"), "output", exists=False); output.parent.resolve(strict=True).relative_to(root)
        receipt = validate(profile, root); atomic(output, (json.dumps(receipt, indent=2, sort_keys=True) + "\n").encode())
        print(f"K6_LOAD_ADMISSION_PASS project={receipt['project_id']} release={receipt['release_digest']}"); return 0
    except (ValueError, OSError) as error:
        print(f"K6_LOAD_ADMISSION_FAILED: {error}", file=sys.stderr); return 1


if __name__ == "__main__": sys.exit(main())
````

### FILE: `production_admission_gate/test_validate_k6_load_admission.py`

```yaml
block_id: "SECURE-OPS:k6-load-admission-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local regression tests for official Grafana k6 execution interpretation"
license: "LicenseRef-Workspace-Owner"
sha256: "70d7fb9fa8c39736c2a13e71abbb71ec7571a50b013b460472180306403f95de"
variables: []
secrets_allowed: false
```

````python
from __future__ import annotations

import copy
from datetime import datetime, timedelta, timezone
import hashlib
import json
from pathlib import Path
import tempfile
import unittest

from validate_k6_load_admission import canonical_sha, validate


class K6LoadAdmissionTests(unittest.TestCase):
    def setUp(self) -> None:
        self.temp = tempfile.TemporaryDirectory(); self.root = Path(self.temp.name).resolve(); (self.root / "evidence").mkdir()
        now = datetime.now(timezone.utc); stamp = lambda value: value.isoformat().replace("+00:00", "Z")
        self.profile = {"schema":"elite-k6-load-admission/v1", "project_id":"revestex", "environment":"production", "release_digest":"sha256:" + "a" * 64, "evaluated_at":stamp(now), "expires_at":stamp(now + timedelta(days=7)), "executor_ref":"ci:run/42", "tool":{"name":"Grafana k6", "version":"2.2.0", "source":"https://github.com/grafana/k6/releases/tag/v2.2.0", "sha256":"87dfa91bc3e47bc4bd77911d59d7ff79f25cd76fb8322c97072f08f08a0da5ed"}, "target":{"service":"https://production.invalid", "region":"primary"}, "error_budget":"evidence/budget.json", "runs":[], "output":"evidence/load-control.json"}
        budget_runs = {}
        for kind in ("BASELINE", "SOAK", "SATURATION_RECOVERY"):
            budget_runs[kind] = {"minimum_iterations":2, "minimum_duration_seconds":60, "maximum_http_req_failed_rate":0.01, "maximum_http_req_duration_p95_ms":500}
        self.write("evidence/budget.json", {"schema":"elite-load-error-budget/v1", "project_id":"revestex", "environment":"production", "release_digest":self.profile["release_digest"], "approved_at":stamp(now), "approvers":["performance-owner", "service-owner"], "runs":budget_runs})
        for kind in ("BASELINE", "SOAK", "SATURATION_RECOVERY"):
            stem = kind.lower(); script = self.root / "evidence" / f"{stem}.js"; script.write_text("export default function() {}\n", encoding="utf-8")
            args = ["run", "--summary-export", f"evidence/{stem}-summary.json", f"evidence/{stem}.js"]
            target = dict(self.profile["target"]); target.update({"run_kind":kind, "workload_script_sha256":hashlib.sha256(script.read_bytes()).hexdigest(), "minimum_duration_seconds":60.0})
            stdout = self.root / "evidence" / f"{stem}.stdout.bin"; stderr = self.root / "evidence" / f"{stem}.stderr.bin"; stdout.write_bytes(b"pass\n"); stderr.write_bytes(b"")
            pf = lambda path: {"path":path.relative_to(self.root).as_posix(), "bytes":path.stat().st_size, "sha256":hashlib.sha256(path.read_bytes()).hexdigest()}
            receipt = {"schema":"elite-official-tool-execution/v1", "project_id":"revestex", "environment":"production", "release_digest":self.profile["release_digest"], "control_id":"LOAD_RESILIENCE", "executed_at":stamp(now), "exit_code":0, "tool":dict(self.profile["tool"]), "arguments_sha256":canonical_sha(args), "working_directory":".", "target":target, "environment_variable_names":[], "stdout":pf(stdout), "stderr":pf(stderr)}
            self.write(f"evidence/{stem}-receipt.json", receipt)
            metrics = [{"name":"iterations", "values":{"count":2}}, {"name":"http_req_failed", "values":{"rate":0.0}}, {"name":"http_req_duration", "values":{"p(95)":120.0}}]
            self.write(f"evidence/{stem}-summary.json", {"version":"1.0.0", "metadata":{"k6_version":"2.2.0"}, "results":{"metrics":metrics}})
            self.profile["runs"].append({"kind":kind, "execution_receipt":f"evidence/{stem}-receipt.json", "summary":f"evidence/{stem}-summary.json", "script":f"evidence/{stem}.js", "expected_arguments":args})

    def tearDown(self) -> None: self.temp.cleanup()
    def write(self, relative: str, value: object) -> None: (self.root / relative).write_text(json.dumps(value), encoding="utf-8")

    def test_three_distinct_runs_emit_semantic_load_receipt(self) -> None:
        result = validate(self.profile, self.root); self.assertEqual("PASS", result["result"]); self.assertEqual(16, len(result["evidence"])); self.assertTrue(all(result["assertions"].values()))

    def test_nonzero_k6_exit_is_rejected(self) -> None:
        value = copy.deepcopy(self.profile); path = self.root / value["runs"][0]["execution_receipt"]; receipt = json.loads(path.read_text()); receipt["exit_code"] = 99; self.write(value["runs"][0]["execution_receipt"], receipt)
        with self.assertRaisesRegex(ValueError, "identity/exit"): validate(value, self.root)

    def test_script_tamper_is_rejected(self) -> None:
        (self.root / self.profile["runs"][1]["script"]).write_text("tampered\n", encoding="utf-8")
        with self.assertRaisesRegex(ValueError, "binding mismatch"): validate(self.profile, self.root)

    def test_argument_drift_is_rejected(self) -> None:
        value = copy.deepcopy(self.profile); value["runs"][0]["expected_arguments"].append("--insecure-skip-tls-verify")
        with self.assertRaisesRegex(ValueError, "hash-bound"): validate(value, self.root)

    def test_budget_breach_is_rejected(self) -> None:
        path = self.root / self.profile["runs"][2]["summary"]; summary = json.loads(path.read_text()); summary["results"]["metrics"][1]["values"]["rate"] = 0.02; self.write(self.profile["runs"][2]["summary"], summary)
        with self.assertRaisesRegex(ValueError, "failure rate"): validate(self.profile, self.root)

    def test_missing_soak_is_rejected(self) -> None:
        value = copy.deepcopy(self.profile); value["runs"].pop(1)
        with self.assertRaisesRegex(ValueError, "three exact"): validate(value, self.root)

    def test_single_approver_is_rejected(self) -> None:
        path = self.root / self.profile["error_budget"]; budget = json.loads(path.read_text()); budget["approvers"] = ["one"]; self.write(self.profile["error_budget"], budget)
        with self.assertRaisesRegex(ValueError, "two distinct"): validate(self.profile, self.root)

    def test_legacy_summary_is_rejected(self) -> None:
        path = self.root / self.profile["runs"][0]["summary"]; summary = json.loads(path.read_text()); summary["version"] = "legacy"; self.write(self.profile["runs"][0]["summary"], summary)
        with self.assertRaisesRegex(ValueError, "machine-readable"): validate(self.profile, self.root)


if __name__ == "__main__": unittest.main()
````

### FILE: `production_admission_gate/postgres-recovery-admission.template.json`

```yaml
block_id: "SECURE-OPS:postgres-recovery-admission-template:v1"
operation: CREATE
provenance: AUTHORED
source: "local semantic adapter over exact official PostgreSQL 18.6 recovery tool evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "d1be51eedde7ea2cafa297ce5c6c74540a243ddde5ccc8176e073f178c3c1751"
variables: []
secrets_allowed: false
```

````json
{
  "schema":"elite-postgres-recovery-admission/v1",
  "project_id":"",
  "environment":"production",
  "release_digest":"",
  "evaluated_at":"",
  "expires_at":"",
  "executor_ref":"",
  "target":{},
  "policy":"",
  "drill":"",
  "runs":[
    {"kind":"VERIFY_BACKUP","execution_receipt":"","expected_arguments":[],"input_files":[]},
    {"kind":"RESTORE_START","execution_receipt":"","expected_arguments":[],"input_files":[]},
    {"kind":"RESTORE_BOOT","execution_receipt":"","expected_arguments":[],"input_files":[]},
    {"kind":"DATA_INTEGRITY","execution_receipt":"","expected_arguments":[],"input_files":[]},
    {"kind":"PITR","execution_receipt":"","expected_arguments":[],"input_files":[]}
  ],
  "output":""
}
````

### FILE: `production_admission_gate/validate_postgres_recovery_admission.py`

```yaml
block_id: "SECURE-OPS:postgres-recovery-admission-validator:v1"
operation: CREATE
provenance: AUTHORED
source: "local semantic adapter over exact official PostgreSQL 18.6 recovery tool evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "cdb4ccf494cd59c7bab42ae9a07521e810ff226aabe07d55302d60a77e2e08c9"
variables: []
secrets_allowed: false
```

````python
from __future__ import annotations

import argparse
from datetime import datetime
import json
from pathlib import Path
import re
import sys

from validate_k6_load_admission import atomic, canonical_sha, proof, read_json, safe, utc, verify_proof

KINDS = ("VERIFY_BACKUP", "RESTORE_START", "RESTORE_BOOT", "DATA_INTEGRITY", "PITR")
TOOLS = {"VERIFY_BACKUP":"PostgreSQL pg_verifybackup", "RESTORE_START":"PostgreSQL pg_ctl", "RESTORE_BOOT":"PostgreSQL pg_isready", "DATA_INTEGRITY":"PostgreSQL psql", "PITR":"PostgreSQL psql"}
SOURCE = "https://ftp.postgresql.org/pub/source/v18.6/"
PROFILE_FIELDS = {"schema","project_id","environment","release_digest","evaluated_at","expires_at","executor_ref","target","policy","drill","runs","output"}
RUN_FIELDS = {"kind","execution_receipt","expected_arguments","input_files"}
LSN = re.compile(r"[0-9A-F]+/[0-9A-F]+")


def seconds(later: datetime, earlier: datetime, label: str) -> float:
    value = (later - earlier).total_seconds()
    if value < 0: raise ValueError(f"{label} cannot be negative")
    return value


def validate(profile: object, root: Path) -> dict[str, object]:
    if not isinstance(profile, dict) or set(profile) != PROFILE_FIELDS: raise ValueError("profile fields must be exact")
    if profile["schema"] != "elite-postgres-recovery-admission/v1" or profile["environment"] != "production": raise ValueError("profile identity is invalid")
    if not isinstance(profile["project_id"], str) or not profile["project_id"].strip() or not re.fullmatch(r"sha256:[0-9a-f]{64}", str(profile["release_digest"])): raise ValueError("project/release identity is invalid")
    evaluated = utc(profile["evaluated_at"], "evaluated_at"); expires = utc(profile["expires_at"], "expires_at")
    if not evaluated < expires or seconds(expires, evaluated, "validity") > 30 * 86400: raise ValueError("admission validity must be positive and at most 30 days")
    if not isinstance(profile["executor_ref"], str) or not profile["executor_ref"].strip() or not isinstance(profile["target"], dict) or set(profile["target"]) != {"cluster","backup_id","recovery_target_time","restore_cluster"}: raise ValueError("exact real recovery target is required")
    target = profile["target"]
    if any(not isinstance(target[key], str) or not target[key].strip() for key in target): raise ValueError("recovery target values are required")
    recovery_target = utc(target["recovery_target_time"], "target.recovery_target_time")

    policy_path = safe(root, profile["policy"], "policy"); policy = read_json(policy_path, "policy")
    if set(policy) != {"schema","project_id","environment","release_digest","approved_at","approvers","maximum_rpo_seconds","maximum_rto_seconds"} or policy.get("schema") != "elite-postgres-recovery-policy/v1": raise ValueError("recovery policy fields/schema must be exact")
    for key in ("project_id","environment","release_digest"):
        if policy.get(key) != profile[key]: raise ValueError(f"policy {key} mismatch")
    utc(policy["approved_at"], "policy.approved_at")
    approvers = policy["approvers"]
    if not isinstance(approvers, list) or len(approvers) < 2 or len(set(approvers)) != len(approvers) or any(not isinstance(x, str) or not x.strip() for x in approvers): raise ValueError("policy requires two distinct approvers")
    max_rpo = policy["maximum_rpo_seconds"]; max_rto = policy["maximum_rto_seconds"]
    if isinstance(max_rpo, bool) or isinstance(max_rto, bool) or not isinstance(max_rpo, (int,float)) or not isinstance(max_rto, (int,float)) or max_rpo < 0 or max_rto <= 0: raise ValueError("RPO/RTO policy values are invalid")

    drill_path = safe(root, profile["drill"], "drill"); drill = read_json(drill_path, "drill")
    drill_fields = {"schema","project_id","environment","release_digest","backup_id","recovery_target_time","recovered_through","restore_started_at","restore_ready_at","wal_archive","invariants"}
    if set(drill) != drill_fields or drill.get("schema") != "elite-postgres-recovery-drill/v1": raise ValueError("drill fields/schema must be exact")
    for key in ("project_id","environment","release_digest"):
        if drill.get(key) != profile[key]: raise ValueError(f"drill {key} mismatch")
    if drill.get("backup_id") != target["backup_id"] or drill.get("recovery_target_time") != target["recovery_target_time"]: raise ValueError("drill target identity mismatch")
    recovered = utc(drill["recovered_through"], "drill.recovered_through"); started = utc(drill["restore_started_at"], "drill.restore_started_at"); ready = utc(drill["restore_ready_at"], "drill.restore_ready_at")
    rpo = seconds(recovery_target, recovered, "achieved RPO"); rto = seconds(ready, started, "achieved RTO")
    if rpo > max_rpo: raise ValueError("achieved RPO exceeds approved maximum")
    if rto > max_rto: raise ValueError("achieved RTO exceeds approved maximum")
    wal = drill["wal_archive"]
    if not isinstance(wal, dict) or set(wal) != {"timeline","start_lsn","end_lsn","location_digest"} or not isinstance(wal["timeline"], int) or wal["timeline"] < 1 or not LSN.fullmatch(str(wal["start_lsn"])) or not LSN.fullmatch(str(wal["end_lsn"])) or not re.fullmatch(r"sha256:[0-9a-f]{64}", str(wal["location_digest"])): raise ValueError("WAL archive identity is invalid")
    invariants = drill["invariants"]
    if not isinstance(invariants, list) or not invariants or any(not isinstance(x, dict) or set(x) != {"name","pass"} or not isinstance(x["name"], str) or not x["name"].strip() or x["pass"] is not True for x in invariants): raise ValueError("all named data invariants must pass")

    runs = profile["runs"]
    if not isinstance(runs, list) or len(runs) != 5 or any(not isinstance(x, dict) or set(x) != RUN_FIELDS for x in runs) or [x["kind"] for x in runs] != list(KINDS): raise ValueError("five ordered recovery runs are required")
    evidence = [proof(root, policy_path), proof(root, drill_path)]
    for run in runs:
        kind = run["kind"]; receipt_path = safe(root, run["execution_receipt"], f"{kind}.receipt"); receipt = read_json(receipt_path, f"{kind}.receipt")
        if receipt.get("schema") != "elite-official-tool-execution/v1" or receipt.get("project_id") != profile["project_id"] or receipt.get("environment") != "production" or receipt.get("release_digest") != profile["release_digest"] or receipt.get("control_id") != "POSTGRES_RECOVERY" or receipt.get("exit_code") != 0: raise ValueError(f"{kind} execution identity/exit mismatch")
        tool = receipt.get("tool")
        if not isinstance(tool, dict) or set(tool) != {"name","version","source","sha256"} or tool.get("name") != TOOLS[kind] or tool.get("version") != "18.6" or tool.get("source") != SOURCE or not re.fullmatch(r"[0-9a-f]{64}", str(tool.get("sha256"))): raise ValueError(f"{kind} exact PostgreSQL tool identity mismatch")
        args = run["expected_arguments"]
        if not isinstance(args, list) or any(not isinstance(x, str) or not x for x in args) or receipt.get("arguments_sha256") != canonical_sha(args): raise ValueError(f"{kind} arguments are not hash-bound")
        inputs = run["input_files"]
        if not isinstance(inputs, list) or len(inputs) != len(set(inputs)): raise ValueError(f"{kind} input_files are invalid")
        input_proofs = [proof(root, safe(root, item, f"{kind}.input")) for item in inputs]
        expected_target = dict(target); expected_target.update({"run_kind":kind, "input_sha256s":[item["sha256"] for item in input_proofs]})
        if receipt.get("target") != expected_target: raise ValueError(f"{kind} target/input binding mismatch")
        for channel in ("stdout","stderr"): verify_proof(root, receipt.get(channel), f"{kind}.{channel}")
        evidence.extend([proof(root, receipt_path), receipt["stdout"], receipt["stderr"], *input_proofs])

    return {"id":"POSTGRES_RECOVERY", "project_id":profile["project_id"], "environment":"production", "release_digest":profile["release_digest"], "result":"PASS", "executed_at":profile["evaluated_at"], "expires_at":profile["expires_at"], "executor_ref":profile["executor_ref"], "tool":{"name":"PostgreSQL recovery toolchain", "version":"18.6", "source":SOURCE, "digest":"sha256:" + canonical_sha([run["kind"] for run in runs])}, "target":target, "assertions":{"pg_verifybackup_pass":True,"restore_boot_pass":True,"data_integrity_pass":True,"pitr_pass":True,"rpo_met":True,"rto_met":True}, "evidence":evidence}


def main() -> int:
    parser=argparse.ArgumentParser(); parser.add_argument("--project-root",required=True,type=Path); parser.add_argument("--profile",required=True); args=parser.parse_args()
    try:
        root=args.project_root.resolve(strict=True); profile=read_json(safe(root,args.profile,"profile"),"profile"); output=safe(root,profile.get("output"),"output",exists=False); output.parent.resolve(strict=True).relative_to(root); receipt=validate(profile,root); atomic(output,(json.dumps(receipt,indent=2,sort_keys=True)+"\n").encode()); print(f"POSTGRES_RECOVERY_ADMISSION_PASS project={receipt['project_id']} release={receipt['release_digest']}"); return 0
    except (ValueError,OSError) as error: print(f"POSTGRES_RECOVERY_ADMISSION_FAILED: {error}",file=sys.stderr); return 1


if __name__ == "__main__": sys.exit(main())
````

### FILE: `production_admission_gate/test_validate_postgres_recovery_admission.py`

```yaml
block_id: "SECURE-OPS:postgres-recovery-admission-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local regressions for exact PostgreSQL verify/restore/PITR evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "8679bd82b7a2c24e5fb59709ba1fb5ead77ed8cf0d2f25e8a7a439cc982e5bf5"
variables: []
secrets_allowed: false
```

````python
from __future__ import annotations
import copy, hashlib, json
from datetime import datetime, timedelta, timezone
from pathlib import Path
import tempfile, unittest
from validate_k6_load_admission import canonical_sha
from validate_postgres_recovery_admission import KINDS, SOURCE, TOOLS, validate

class PostgresRecoveryAdmissionTests(unittest.TestCase):
    def setUp(self):
        self.temp=tempfile.TemporaryDirectory(); self.root=Path(self.temp.name).resolve(); (self.root/"evidence").mkdir(); now=datetime.now(timezone.utc); z=lambda d:d.isoformat().replace("+00:00","Z"); digest="sha256:"+"b"*64
        self.profile={"schema":"elite-postgres-recovery-admission/v1","project_id":"revestex","environment":"production","release_digest":digest,"evaluated_at":z(now),"expires_at":z(now+timedelta(days=7)),"executor_ref":"ci:recovery/42","target":{"cluster":"primary","backup_id":"basebackup-42","recovery_target_time":z(now-timedelta(seconds=10)),"restore_cluster":"isolated-drill-42"},"policy":"evidence/policy.json","drill":"evidence/drill.json","runs":[],"output":"evidence/postgres-control.json"}
        self.write("evidence/policy.json",{"schema":"elite-postgres-recovery-policy/v1","project_id":"revestex","environment":"production","release_digest":digest,"approved_at":z(now-timedelta(days=1)),"approvers":["database-owner","business-owner"],"maximum_rpo_seconds":60,"maximum_rto_seconds":300})
        self.write("evidence/drill.json",{"schema":"elite-postgres-recovery-drill/v1","project_id":"revestex","environment":"production","release_digest":digest,"backup_id":"basebackup-42","recovery_target_time":self.profile["target"]["recovery_target_time"],"recovered_through":z(now-timedelta(seconds=20)),"restore_started_at":z(now-timedelta(seconds=200)),"restore_ready_at":z(now-timedelta(seconds=20)),"wal_archive":{"timeline":1,"start_lsn":"0/1000000","end_lsn":"0/2000000","location_digest":"sha256:"+"c"*64},"invariants":[{"name":"orders-total","pass":True},{"name":"ledger-balance","pass":True}]})
        for kind in KINDS:
            stem=kind.lower(); inp=self.root/"evidence"/f"{stem}.input"; inp.write_text(kind+"\n"); args=["--elite-fixture",kind]; out=self.root/"evidence"/f"{stem}.stdout.bin"; err=self.root/"evidence"/f"{stem}.stderr.bin"; out.write_bytes(b"PASS\n"); err.write_bytes(b"")
            pf=lambda p:{"path":p.relative_to(self.root).as_posix(),"bytes":p.stat().st_size,"sha256":hashlib.sha256(p.read_bytes()).hexdigest()}; target=dict(self.profile["target"]); target.update({"run_kind":kind,"input_sha256s":[pf(inp)["sha256"]]})
            receipt={"schema":"elite-official-tool-execution/v1","project_id":"revestex","environment":"production","release_digest":digest,"control_id":"POSTGRES_RECOVERY","executed_at":z(now),"exit_code":0,"tool":{"name":TOOLS[kind],"version":"18.6","source":SOURCE,"sha256":"d"*64},"arguments_sha256":canonical_sha(args),"working_directory":".","target":target,"environment_variable_names":[],"stdout":pf(out),"stderr":pf(err)}; self.write(f"evidence/{stem}-receipt.json",receipt); self.profile["runs"].append({"kind":kind,"execution_receipt":f"evidence/{stem}-receipt.json","expected_arguments":args,"input_files":[f"evidence/{stem}.input"]})
    def tearDown(self): self.temp.cleanup()
    def write(self,path,value): (self.root/path).write_text(json.dumps(value),encoding="utf-8")
    def test_full_restore_pitr_and_objectives_emit_semantic_receipt(self): self.assertTrue(all(validate(self.profile,self.root)["assertions"].values()))
    def test_verifybackup_nonzero_is_rejected(self):
        p=self.root/self.profile["runs"][0]["execution_receipt"]; v=json.loads(p.read_text()); v["exit_code"]=1; self.write(self.profile["runs"][0]["execution_receipt"],v)
        with self.assertRaisesRegex(ValueError,"identity/exit"): validate(self.profile,self.root)
    def test_logical_or_verify_only_without_five_runs_is_rejected(self):
        v=copy.deepcopy(self.profile); v["runs"]=v["runs"][:1]
        with self.assertRaisesRegex(ValueError,"five ordered"): validate(v,self.root)
    def test_rpo_breach_is_rejected(self):
        p=self.root/self.profile["drill"]; v=json.loads(p.read_text()); v["recovered_through"]="2020-01-01T00:00:00Z"; self.write(self.profile["drill"],v)
        with self.assertRaisesRegex(ValueError,"RPO exceeds"): validate(self.profile,self.root)
    def test_rto_breach_is_rejected(self):
        p=self.root/self.profile["drill"]; v=json.loads(p.read_text()); v["restore_started_at"]="2020-01-01T00:00:00Z"; self.write(self.profile["drill"],v)
        with self.assertRaisesRegex(ValueError,"RTO exceeds"): validate(self.profile,self.root)
    def test_failed_invariant_is_rejected(self):
        p=self.root/self.profile["drill"]; v=json.loads(p.read_text()); v["invariants"][0]["pass"]=False; self.write(self.profile["drill"],v)
        with self.assertRaisesRegex(ValueError,"invariants"): validate(self.profile,self.root)
    def test_input_tamper_is_rejected(self):
        (self.root/self.profile["runs"][3]["input_files"][0]).write_text("tampered")
        with self.assertRaisesRegex(ValueError,"binding mismatch"): validate(self.profile,self.root)
    def test_tool_version_drift_is_rejected(self):
        p=self.root/self.profile["runs"][4]["execution_receipt"]; v=json.loads(p.read_text()); v["tool"]["version"]="18.5"; self.write(self.profile["runs"][4]["execution_receipt"],v)
        with self.assertRaisesRegex(ValueError,"tool identity"): validate(self.profile,self.root)

if __name__=="__main__": unittest.main()
````

### FILE: `production_admission_gate/cloudflare-edge-admission.template.json`

```yaml
block_id: "SECURE-OPS:cloudflare-edge-admission-template:v1"
operation: CREATE
provenance: AUTHORED
source: "local schema binding official Cloudflare evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "6d991aba14235d83127cf7475ed5902604900326960d5d21642f7730b52fca78"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite-cloudflare-edge-admission/v1",
  "project_id": "",
  "environment": "production",
  "release_digest": "",
  "evaluated_at": "",
  "expires_at": "",
  "executor_ref": "",
  "tool": {
    "name": "Cloudflare cloudflare-go",
    "version": "7.9.0",
    "source": "https://github.com/cloudflare/cloudflare-go/releases/tag/v7.9.0"
  },
  "target": {
    "zone_id": "",
    "hostname": "",
    "origin_address": "",
    "public_url": "",
    "private_url": ""
  },
  "policy": "",
  "runs": [
    {"kind":"CLOUDFLARE_CONFIG_SNAPSHOT", "execution_receipt":"", "observation":"", "expected_arguments":[]},
    {"kind":"ORIGIN_BYPASS_NEGATIVE", "execution_receipt":"", "observation":"", "expected_arguments":[]},
    {"kind":"PUBLIC_CACHE_BEHAVIOR", "execution_receipt":"", "observation":"", "expected_arguments":[]},
    {"kind":"PRIVATE_CACHE_BYPASS", "execution_receipt":"", "observation":"", "expected_arguments":[]}
  ],
  "output": ""
}
````

### FILE: `production_admission_gate/validate_cloudflare_edge_admission.py`

```yaml
block_id: "SECURE-OPS:cloudflare-edge-admission-validator:v1"
operation: CREATE
provenance: AUTHORED
source: "local semantic validator over exact cloudflare-go and target evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "c35b5c5759ca62cb5732578d607f1b355b7eecc6bb4a1880a5913d5358e65c7d"
variables: []
secrets_allowed: false
```

````python
from __future__ import annotations
import argparse,json,re,sys
from pathlib import Path
from urllib.parse import urlparse
from validate_k6_load_admission import atomic,canonical_sha,proof,read_json,safe,utc,verify_proof

KINDS=("CLOUDFLARE_CONFIG_SNAPSHOT","ORIGIN_BYPASS_NEGATIVE","PUBLIC_CACHE_BEHAVIOR","PRIVATE_CACHE_BYPASS")
FIELDS={"schema","project_id","environment","release_digest","evaluated_at","expires_at","executor_ref","tool","target","policy","runs","output"}
RUN_FIELDS={"kind","execution_receipt","observation","expected_arguments"}
TARGET_FIELDS={"zone_id","hostname","origin_address","public_url","private_url"}
SOURCE="https://github.com/cloudflare/cloudflare-go/releases/tag/v7.9.0"

def nonempty(value:object,label:str)->str:
 if not isinstance(value,str) or not value.strip(): raise ValueError(f"{label} is required")
 return value

def https(value:object,label:str)->None:
 value=nonempty(value,label); parsed=urlparse(value)
 if parsed.scheme!="https" or not parsed.hostname or parsed.username or parsed.password or parsed.fragment: raise ValueError(f"{label} must be credential-free HTTPS")

def bind(value:dict[str,object],profile:dict[str,object],label:str)->None:
 for key in ("project_id","environment","release_digest"):
  if value.get(key)!=profile[key]: raise ValueError(f"{label} {key} mismatch")

def validate_snapshot(value:dict[str,object],profile:dict[str,object],policy:dict[str,object])->None:
 required={"schema","project_id","environment","release_digest","captured_at","zone_id","hostname","dns_records","waf_rulesets","tls"}
 if set(value)!=required or value.get("schema")!="elite-cloudflare-edge-snapshot/v1": raise ValueError("Cloudflare snapshot fields/schema must be exact")
 bind(value,profile,"snapshot"); utc(value["captured_at"],"snapshot.captured_at")
 target=profile["target"]
 if value["zone_id"]!=target["zone_id"] or value["hostname"]!=target["hostname"]: raise ValueError("snapshot target mismatch")
 records=value["dns_records"]
 if not isinstance(records,list) or not records: raise ValueError("DNS records are required")
 address=[]
 for item in records:
  if not isinstance(item,dict) or set(item)!={"id","name","type","content","proxied"}: raise ValueError("DNS record fields must be exact")
  if item["name"]==target["hostname"] and item["type"] in ("A","AAAA","CNAME"): address.append(item)
 if not address or any(item["proxied"] is not True or not nonempty(item["id"],"DNS id") or not nonempty(item["content"],"DNS content") for item in address): raise ValueError("all hostname address records must be proxied")
 rules=value["waf_rulesets"]
 if not isinstance(rules,list) or not rules: raise ValueError("WAF rulesets are required")
 normalized={}
 for item in rules:
  if not isinstance(item,dict) or set(item)!={"id","name","kind","phase","version","status","enabled_actions"}: raise ValueError("WAF ruleset fields must be exact")
  if item["status"]!="active" or not isinstance(item["enabled_actions"],list) or not item["enabled_actions"]: raise ValueError("WAF rulesets must be active with enabled actions")
  normalized[item["id"]]=item
 managed=policy["managed_ruleset_ids"]; custom=policy["custom_ruleset_ids"]
 if not managed or not custom or any(i not in normalized for i in managed+custom): raise ValueError("approved managed and custom WAF rulesets are required")
 if any(normalized[i]["phase"]!="http_request_firewall_managed" or "execute" not in normalized[i]["enabled_actions"] for i in managed): raise ValueError("managed WAF execute ruleset is not active")
 defensive={"block","challenge","managed_challenge","js_challenge"}
 if any(normalized[i]["phase"]!="http_request_firewall_custom" or not defensive.intersection(normalized[i]["enabled_actions"]) for i in custom): raise ValueError("custom WAF defensive ruleset is not active")
 tls=value["tls"]
 if not isinstance(tls,dict) or set(tls)!={"mode","minimum_version","always_use_https","certificate_status"}: raise ValueError("TLS fields must be exact")
 if tls!={"mode":"full_strict","minimum_version":policy["minimum_tls_version"],"always_use_https":True,"certificate_status":"active"} or tls["minimum_version"] not in ("1.2","1.3"): raise ValueError("Full (strict) TLS policy is not proven")

def validate_observation(kind:str,value:dict[str,object],profile:dict[str,object],policy:dict[str,object])->None:
 fields={"schema","project_id","environment","release_digest","captured_at","kind","target","result"}
 if set(value)!=fields or value.get("schema")!="elite-cloudflare-edge-observation/v1" or value.get("kind")!=kind: raise ValueError(f"{kind} observation fields/schema mismatch")
 bind(value,profile,kind); utc(value["captured_at"],f"{kind}.captured_at")
 result=value["result"]
 if not isinstance(result,dict): raise ValueError(f"{kind} result must be an object")
 if kind=="ORIGIN_BYPASS_NEGATIVE":
  if value["target"]!=profile["target"]["origin_address"] or set(result)!={"connection_succeeded","http_status","blocked","protection_mode"}: raise ValueError("origin bypass observation mismatch")
  if result["blocked"] is not True or result["protection_mode"]!=policy["origin_protection_mode"] or (result["connection_succeeded"] is True and result["http_status"] not in (400,401,403,421,525,526)): raise ValueError("origin bypass was not blocked")
 elif kind=="PUBLIC_CACHE_BEHAVIOR":
  if value["target"]!=profile["target"]["public_url"] or set(result)!={"authorization_sent","first_status","second_status","first_cf_cache_status","second_cf_cache_status","cache_control","first_body_sha256","second_body_sha256"}: raise ValueError("public cache observation mismatch")
  if result["authorization_sent"] is not False or result["first_status"]!=200 or result["second_status"]!=200 or result["second_cf_cache_status"]!="HIT" or result["first_body_sha256"]!=result["second_body_sha256"] or not re.fullmatch(r"[0-9a-f]{64}",str(result["first_body_sha256"])) or re.search(r"(?:private|no-store)",str(result["cache_control"]),re.I): raise ValueError("public cache behavior is not proven")
 else:
  if value["target"]!=profile["target"]["private_url"] or set(result)!={"authorization_sent","status","cf_cache_status","cache_control","set_cookie","shared_response_reused"}: raise ValueError("private cache observation mismatch")
  if result["authorization_sent"] is not True or not isinstance(result["status"],int) or not 200<=result["status"]<400 or result["cf_cache_status"] not in ("BYPASS","DYNAMIC") or not re.search(r"(?:private|no-store)",str(result["cache_control"]),re.I) or result["shared_response_reused"] is not False: raise ValueError("private response cache isolation is not proven")

def validate(profile:object,root:Path)->dict[str,object]:
 if not isinstance(profile,dict) or set(profile)!=FIELDS: raise ValueError("profile fields must be exact")
 if profile["schema"]!="elite-cloudflare-edge-admission/v1" or profile["environment"]!="production" or not nonempty(profile["project_id"],"project_id") or not re.fullmatch(r"sha256:[0-9a-f]{64}",str(profile["release_digest"])): raise ValueError("profile identity is invalid")
 start=utc(profile["evaluated_at"],"evaluated_at"); end=utc(profile["expires_at"],"expires_at")
 if not start<end or (end-start).total_seconds()>30*86400: raise ValueError("admission validity must be positive and at most 30 days")
 nonempty(profile["executor_ref"],"executor_ref")
 if profile["tool"]!={"name":"Cloudflare cloudflare-go","version":"7.9.0","source":SOURCE}: raise ValueError("exact Cloudflare cloudflare-go v7.9.0 identity is required")
 target=profile["target"]
 if not isinstance(target,dict) or set(target)!=TARGET_FIELDS: raise ValueError("target fields must be exact")
 for key in ("zone_id","hostname","origin_address"): nonempty(target[key],f"target.{key}")
 if target["origin_address"]==target["hostname"]: raise ValueError("origin must be distinct from public hostname")
 https(target["public_url"],"target.public_url"); https(target["private_url"],"target.private_url")
 if urlparse(target["public_url"]).hostname!=target["hostname"] or urlparse(target["private_url"]).hostname!=target["hostname"]: raise ValueError("edge URLs must use the admitted hostname")
 policy_path=safe(root,profile["policy"],"policy"); policy=read_json(policy_path,"policy")
 policy_fields={"schema","project_id","environment","release_digest","approved_at","approvers","zone_id","hostname","managed_ruleset_ids","custom_ruleset_ids","minimum_tls_version","origin_protection_mode","public_cache_path","private_cache_path"}
 if set(policy)!=policy_fields or policy.get("schema")!="elite-cloudflare-edge-policy/v1": raise ValueError("edge policy fields/schema must be exact")
 bind(policy,profile,"policy"); utc(policy["approved_at"],"policy.approved_at")
 approvers=policy["approvers"]
 if not isinstance(approvers,list) or len(approvers)<2 or len(set(approvers))!=len(approvers) or any(not isinstance(x,str) or not x.strip() for x in approvers): raise ValueError("edge policy requires two distinct approvers")
 if policy["zone_id"]!=target["zone_id"] or policy["hostname"]!=target["hostname"] or policy["origin_protection_mode"] not in ("authenticated_origin_pull","mtls","cloudflare_tunnel","private_network"): raise ValueError("edge policy target/origin protection mismatch")
 for key in ("managed_ruleset_ids","custom_ruleset_ids"):
  if not isinstance(policy[key],list) or len(set(policy[key]))!=len(policy[key]) or any(not isinstance(x,str) or not x for x in policy[key]): raise ValueError(f"{key} must contain unique IDs")
 if urlparse(target["public_url"]).path!=policy["public_cache_path"] or urlparse(target["private_url"]).path!=policy["private_cache_path"]: raise ValueError("approved cache paths mismatch")
 runs=profile["runs"]
 if not isinstance(runs,list) or len(runs)!=4 or any(not isinstance(x,dict) or set(x)!=RUN_FIELDS for x in runs) or [x["kind"] for x in runs]!=list(KINDS): raise ValueError("four ordered edge runs are required")
 evidence=[proof(root,policy_path)]; digest=""
 for run in runs:
  kind=run["kind"]; observation_path=safe(root,run["observation"],f"{kind}.observation"); observation=read_json(observation_path,f"{kind}.observation")
  if kind=="CLOUDFLARE_CONFIG_SNAPSHOT": validate_snapshot(observation,profile,policy)
  else: validate_observation(kind,observation,profile,policy)
  observation_proof=proof(root,observation_path); receipt_path=safe(root,run["execution_receipt"],f"{kind}.receipt"); receipt=read_json(receipt_path,f"{kind}.receipt")
  if receipt.get("schema")!="elite-official-tool-execution/v1" or receipt.get("project_id")!=profile["project_id"] or receipt.get("environment")!="production" or receipt.get("release_digest")!=profile["release_digest"] or receipt.get("control_id")!="EDGE_CDN_WAF" or receipt.get("exit_code")!=0: raise ValueError(f"{kind} execution identity/exit mismatch")
  tool=receipt.get("tool")
  if not isinstance(tool,dict) or set(tool)!={"name","version","source","sha256"} or tool.get("name")!="Project Cloudflare edge probe" or tool.get("version")!="cloudflare-go/7.9.0" or tool.get("source")!=SOURCE or not re.fullmatch(r"[0-9a-f]{64}",str(tool.get("sha256"))): raise ValueError(f"{kind} exact probe identity mismatch")
  digest=tool["sha256"] if not digest else digest
  if tool["sha256"]!=digest: raise ValueError("all edge runs must use the same probe binary")
  args=run["expected_arguments"]
  if not isinstance(args,list) or any(not isinstance(x,str) or not x for x in args) or receipt.get("arguments_sha256")!=canonical_sha(args): raise ValueError(f"{kind} arguments are not hash-bound")
  expected_target=dict(target); expected_target.update({"run_kind":kind,"observation_sha256":observation_proof["sha256"]})
  if receipt.get("target")!=expected_target: raise ValueError(f"{kind} target/observation binding mismatch")
  for channel in ("stdout","stderr"): verify_proof(root,receipt.get(channel),f"{kind}.{channel}")
  evidence.extend([proof(root,receipt_path),observation_proof,receipt["stdout"],receipt["stderr"]])
 return {"id":"EDGE_CDN_WAF","project_id":profile["project_id"],"environment":"production","release_digest":profile["release_digest"],"result":"PASS","executed_at":profile["evaluated_at"],"expires_at":profile["expires_at"],"executor_ref":profile["executor_ref"],"tool":{"name":"Cloudflare cloudflare-go + project edge probe","version":"7.9.0","source":SOURCE,"digest":"sha256:"+digest},"target":target,"assertions":{"dns_proxy_active":True,"waf_rulesets_active":True,"tls_policy_pass":True,"origin_bypass_blocked":True,"cache_behavior_pass":True},"evidence":evidence}

def main()->int:
 p=argparse.ArgumentParser(); p.add_argument("--project-root",required=True,type=Path); p.add_argument("--profile",required=True); a=p.parse_args()
 try:
  root=a.project_root.resolve(strict=True); profile=read_json(safe(root,a.profile,"profile"),"profile"); output=safe(root,profile.get("output"),"output",exists=False); output.parent.resolve(strict=True).relative_to(root); receipt=validate(profile,root); atomic(output,(json.dumps(receipt,indent=2,sort_keys=True)+"\n").encode()); print(f"CLOUDFLARE_EDGE_ADMISSION_PASS project={receipt['project_id']} release={receipt['release_digest']}"); return 0
 except (ValueError,OSError) as error: print(f"CLOUDFLARE_EDGE_ADMISSION_FAILED: {error}",file=sys.stderr); return 1
if __name__=="__main__": sys.exit(main())
````

### FILE: `production_admission_gate/test_validate_cloudflare_edge_admission.py`

```yaml
block_id: "SECURE-OPS:cloudflare-edge-admission-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local regressions for Cloudflare edge semantic evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "c4392d8959ed1630e1a4cca7121e8e038dac3a99ee65aba9eaa3cc0f36b939f9"
variables: []
secrets_allowed: false
```

````python
from __future__ import annotations
import copy,hashlib,json,tempfile,unittest
from datetime import datetime,timedelta,timezone
from pathlib import Path
from validate_cloudflare_edge_admission import KINDS,SOURCE,validate
from validate_k6_load_admission import canonical_sha

class CloudflareEdgeAdmissionTests(unittest.TestCase):
 def setUp(self):
  self.temp=tempfile.TemporaryDirectory(); self.root=Path(self.temp.name); (self.root/"evidence").mkdir(); now=datetime.now(timezone.utc); z=lambda d:d.isoformat().replace("+00:00","Z"); digest="sha256:"+"e"*64
  self.target={"zone_id":"zone-123","hostname":"app.example.test","origin_address":"203.0.113.10","public_url":"https://app.example.test/assets/app.js","private_url":"https://app.example.test/api/me"}
  self.profile={"schema":"elite-cloudflare-edge-admission/v1","project_id":"revestex","environment":"production","release_digest":digest,"evaluated_at":z(now),"expires_at":z(now+timedelta(days=7)),"executor_ref":"edge:change/42","tool":{"name":"Cloudflare cloudflare-go","version":"7.9.0","source":SOURCE},"target":self.target,"policy":"evidence/policy.json","runs":[],"output":"evidence/edge-control.json"}
  self.policy={"schema":"elite-cloudflare-edge-policy/v1","project_id":"revestex","environment":"production","release_digest":digest,"approved_at":z(now),"approvers":["security-owner","platform-owner"],"zone_id":"zone-123","hostname":"app.example.test","managed_ruleset_ids":["managed-1"],"custom_ruleset_ids":["custom-1"],"minimum_tls_version":"1.2","origin_protection_mode":"authenticated_origin_pull","public_cache_path":"/assets/app.js","private_cache_path":"/api/me"}; self.write("evidence/policy.json",self.policy)
  values={
   "CLOUDFLARE_CONFIG_SNAPSHOT":{"schema":"elite-cloudflare-edge-snapshot/v1","project_id":"revestex","environment":"production","release_digest":digest,"captured_at":z(now),"zone_id":"zone-123","hostname":"app.example.test","dns_records":[{"id":"dns-1","name":"app.example.test","type":"A","content":"198.51.100.9","proxied":True}],"waf_rulesets":[{"id":"managed-1","name":"Cloudflare Managed","kind":"zone","phase":"http_request_firewall_managed","version":"9","status":"active","enabled_actions":["execute"]},{"id":"custom-1","name":"Project Custom","kind":"zone","phase":"http_request_firewall_custom","version":"3","status":"active","enabled_actions":["block"]}],"tls":{"mode":"full_strict","minimum_version":"1.2","always_use_https":True,"certificate_status":"active"}},
   "ORIGIN_BYPASS_NEGATIVE":{"schema":"elite-cloudflare-edge-observation/v1","project_id":"revestex","environment":"production","release_digest":digest,"captured_at":z(now),"kind":"ORIGIN_BYPASS_NEGATIVE","target":"203.0.113.10","result":{"connection_succeeded":True,"http_status":403,"blocked":True,"protection_mode":"authenticated_origin_pull"}},
   "PUBLIC_CACHE_BEHAVIOR":{"schema":"elite-cloudflare-edge-observation/v1","project_id":"revestex","environment":"production","release_digest":digest,"captured_at":z(now),"kind":"PUBLIC_CACHE_BEHAVIOR","target":"https://app.example.test/assets/app.js","result":{"authorization_sent":False,"first_status":200,"second_status":200,"first_cf_cache_status":"MISS","second_cf_cache_status":"HIT","cache_control":"public, max-age=3600","first_body_sha256":"a"*64,"second_body_sha256":"a"*64}},
   "PRIVATE_CACHE_BYPASS":{"schema":"elite-cloudflare-edge-observation/v1","project_id":"revestex","environment":"production","release_digest":digest,"captured_at":z(now),"kind":"PRIVATE_CACHE_BYPASS","target":"https://app.example.test/api/me","result":{"authorization_sent":True,"status":200,"cf_cache_status":"BYPASS","cache_control":"private, no-store","set_cookie":True,"shared_response_reused":False}}}
  for kind in KINDS:
   stem=kind.lower(); obs=self.write(f"evidence/{stem}.json",values[kind]); args=["--mode",stem,"--zone-id","zone-123","--hostname","app.example.test"]
   out=self.file(f"evidence/{stem}.stdout",b"edge probe pass\n"); err=self.file(f"evidence/{stem}.stderr",b""); target=dict(self.target); target.update({"run_kind":kind,"observation_sha256":self.pf(obs)["sha256"]})
   receipt={"schema":"elite-official-tool-execution/v1","project_id":"revestex","environment":"production","release_digest":digest,"control_id":"EDGE_CDN_WAF","executed_at":z(now),"exit_code":0,"tool":{"name":"Project Cloudflare edge probe","version":"cloudflare-go/7.9.0","source":SOURCE,"sha256":"f"*64},"arguments_sha256":canonical_sha(args),"working_directory":".","target":target,"environment_variable_names":["CLOUDFLARE_API_TOKEN"],"stdout":self.pf(out),"stderr":self.pf(err)}; self.write(f"evidence/{stem}-receipt.json",receipt); self.profile["runs"].append({"kind":kind,"execution_receipt":f"evidence/{stem}-receipt.json","observation":f"evidence/{stem}.json","expected_arguments":args})
 def tearDown(self): self.temp.cleanup()
 def file(self,p,data): q=self.root/p; q.write_bytes(data); return q
 def write(self,p,v): q=self.root/p; q.write_text(json.dumps(v),encoding="utf-8"); return q
 def pf(self,p): data=p.read_bytes(); return {"path":p.relative_to(self.root).as_posix(),"bytes":len(data),"sha256":hashlib.sha256(data).hexdigest()}
 def mutate_observation(self,index,change): p=self.root/self.profile["runs"][index]["observation"]; v=json.loads(p.read_text()); change(v); self.write(self.profile["runs"][index]["observation"],v)
 def test_complete_edge_evidence_emits_semantic_receipt(self): self.assertTrue(all(validate(self.profile,self.root)["assertions"].values()))
 def test_unproxied_dns_is_rejected(self):
  self.mutate_observation(0,lambda v:v["dns_records"][0].update(proxied=False))
  with self.assertRaisesRegex(ValueError,"proxied"): validate(self.profile,self.root)
 def test_missing_managed_waf_is_rejected(self):
  self.mutate_observation(0,lambda v:v.update(waf_rulesets=v["waf_rulesets"][1:]))
  with self.assertRaisesRegex(ValueError,"managed and custom"): validate(self.profile,self.root)
 def test_weak_tls_is_rejected(self):
  self.mutate_observation(0,lambda v:v["tls"].update(mode="flexible"))
  with self.assertRaisesRegex(ValueError,"strict"): validate(self.profile,self.root)
 def test_origin_bypass_success_is_rejected(self):
  self.mutate_observation(1,lambda v:v["result"].update(blocked=False,http_status=200))
  with self.assertRaisesRegex(ValueError,"bypass was not blocked"): validate(self.profile,self.root)
 def test_public_cache_miss_twice_is_rejected(self):
  self.mutate_observation(2,lambda v:v["result"].update(second_cf_cache_status="MISS"))
  with self.assertRaisesRegex(ValueError,"public cache behavior"): validate(self.profile,self.root)
 def test_private_cache_hit_is_rejected(self):
  self.mutate_observation(3,lambda v:v["result"].update(cf_cache_status="HIT",shared_response_reused=True))
  with self.assertRaisesRegex(ValueError,"cache isolation"): validate(self.profile,self.root)
 def test_nonzero_probe_exit_is_rejected(self):
  p=self.root/self.profile["runs"][0]["execution_receipt"]; v=json.loads(p.read_text()); v["exit_code"]=2; self.write(self.profile["runs"][0]["execution_receipt"],v)
  with self.assertRaisesRegex(ValueError,"identity/exit"): validate(self.profile,self.root)
 def test_observation_tamper_is_rejected_by_receipt_binding(self):
  self.mutate_observation(2,lambda v:v["result"].update(cache_control="public, max-age=7200"))
  with self.assertRaisesRegex(ValueError,"target/observation binding"): validate(self.profile,self.root)

if __name__=="__main__": unittest.main()
````

### FILE: `production_admission_gate/zap-offensive-admission.template.json`

```yaml
block_id: "SECURE-OPS:zap-offensive-admission-template:v1"
operation: CREATE
provenance: AUTHORED
source: "local semantic adapter over exact official OWASP ZAP 2.17.0 evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "43f2af0bc0cc9e3320c59da7412e84c4df62d6ab959c04f42d5a1c66ddc6fe68"
variables: []
secrets_allowed: false
```

````json
{
  "schema":"elite-zap-offensive-admission/v1",
  "project_id":"",
  "environment":"production",
  "release_digest":"",
  "evaluated_at":"",
  "expires_at":"",
  "executor_ref":"",
  "target":{"base_url":"","api_base_url":"","authorization_ticket":""},
  "policy":"",
  "manual_review":"",
  "runs":[
    {"kind":"WEB_AUTHENTICATED","execution_receipt":"","plan":"","report":"","expected_arguments":[]},
    {"kind":"API_AUTHENTICATED","execution_receipt":"","plan":"","report":"","expected_arguments":[]}
  ],
  "output":""
}
````

### FILE: `production_admission_gate/validate_zap_offensive_admission.py`

```yaml
block_id: "SECURE-OPS:zap-offensive-admission-validator:v1"
operation: CREATE
provenance: AUTHORED
source: "local semantic adapter over exact official OWASP ZAP 2.17.0 evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "1467f9982a647e8444d78e6698bc3c466a2546aa4e843b0a15338eda335bd244"
variables: []
secrets_allowed: false
```

````python
from __future__ import annotations
import argparse, json, re, sys
from pathlib import Path
from urllib.parse import urlparse
from validate_k6_load_admission import atomic, canonical_sha, proof, read_json, safe, utc, verify_proof

KINDS=("WEB_AUTHENTICATED","API_AUTHENTICATED")
FIELDS={"schema","project_id","environment","release_digest","evaluated_at","expires_at","executor_ref","target","policy","manual_review","runs","output"}
RUN_FIELDS={"kind","execution_receipt","plan","report","expected_arguments"}
SOURCE="https://github.com/zaproxy/zaproxy/releases/tag/v2.17.0"

def https_url(value: object, label: str) -> None:
    if not isinstance(value,str): raise ValueError(f"{label} must be HTTPS")
    parsed=urlparse(value)
    if parsed.scheme!="https" or not parsed.hostname or parsed.username or parsed.password or parsed.fragment: raise ValueError(f"{label} must be credential-free HTTPS")

def alerts(report: dict[str,object]) -> list[dict[str,object]]:
    if not isinstance(report.get("site"),list) or not report["site"]: raise ValueError("ZAP report requires non-empty site list")
    found=[]
    for site in report["site"]:
        if not isinstance(site,dict) or not isinstance(site.get("@name"),str) or not isinstance(site.get("alerts"),list): raise ValueError("ZAP site/alerts shape is invalid")
        for alert in site["alerts"]:
            if not isinstance(alert,dict) or not {"pluginid","riskcode","name","instances"}.issubset(alert) or not isinstance(alert["instances"],list): raise ValueError("ZAP alert shape is invalid")
            try: risk=int(alert["riskcode"])
            except (TypeError,ValueError) as error: raise ValueError("ZAP riskcode is invalid") from error
            if risk not in (0,1,2,3,4): raise ValueError("ZAP riskcode is outside admitted range")
            found.append({"pluginid":str(alert["pluginid"]),"riskcode":risk,"name":str(alert["name"])})
    return found

def validate(profile: object, root: Path) -> dict[str,object]:
    if not isinstance(profile,dict) or set(profile)!=FIELDS: raise ValueError("profile fields must be exact")
    if profile["schema"]!="elite-zap-offensive-admission/v1" or profile["environment"]!="production" or not isinstance(profile["project_id"],str) or not profile["project_id"].strip() or not re.fullmatch(r"sha256:[0-9a-f]{64}",str(profile["release_digest"])): raise ValueError("profile identity is invalid")
    start=utc(profile["evaluated_at"],"evaluated_at"); end=utc(profile["expires_at"],"expires_at")
    if not start<end or (end-start).total_seconds()>30*86400: raise ValueError("admission validity must be positive and at most 30 days")
    if not isinstance(profile["executor_ref"],str) or not profile["executor_ref"].strip() or not isinstance(profile["target"],dict) or set(profile["target"])!={"base_url","api_base_url","authorization_ticket"}: raise ValueError("exact authorized target is required")
    target=profile["target"]; https_url(target["base_url"],"base_url"); https_url(target["api_base_url"],"api_base_url")
    if not isinstance(target["authorization_ticket"],str) or not target["authorization_ticket"].strip(): raise ValueError("authorization ticket is required")

    policy_path=safe(root,profile["policy"],"policy"); policy=read_json(policy_path,"policy")
    if set(policy)!={"schema","project_id","environment","release_digest","approved_at","approvers","accepted_alerts"} or policy.get("schema")!="elite-offensive-security-policy/v1": raise ValueError("policy fields/schema must be exact")
    for key in ("project_id","environment","release_digest"):
        if policy.get(key)!=profile[key]: raise ValueError(f"policy {key} mismatch")
    utc(policy["approved_at"],"policy.approved_at"); approvers=policy["approvers"]
    if not isinstance(approvers,list) or len(approvers)<2 or len(set(approvers))!=len(approvers) or any(not isinstance(x,str) or not x.strip() for x in approvers): raise ValueError("policy requires two distinct approvers")
    accepted=policy["accepted_alerts"]
    if not isinstance(accepted,list): raise ValueError("accepted_alerts must be a list")
    accepted_ids=set()
    for item in accepted:
        if not isinstance(item,dict) or set(item)!={"pluginid","riskcode","reason","owner","expires_at"} or int(item.get("riskcode",-1))<3 or not all(isinstance(item.get(k),str) and item[k].strip() for k in ("pluginid","reason","owner")): raise ValueError("accepted alert record is invalid")
        if utc(item["expires_at"],"accepted_alert.expires_at")<end: raise ValueError("accepted alert expires before admission")
        key=(item["pluginid"],int(item["riskcode"]));
        if key in accepted_ids: raise ValueError("duplicate accepted alert")
        accepted_ids.add(key)

    review_path=safe(root,profile["manual_review"],"manual_review"); review=read_json(review_path,"manual_review")
    if set(review)!={"schema","project_id","environment","release_digest","reviewed_at","reviewers","authorization_confirmed","authenticated_web_reviewed","authenticated_api_reviewed","business_logic_reviewed","notes"} or review.get("schema")!="elite-offensive-manual-review/v1": raise ValueError("manual review fields/schema must be exact")
    for key in ("project_id","environment","release_digest"):
        if review.get(key)!=profile[key]: raise ValueError(f"manual review {key} mismatch")
    utc(review["reviewed_at"],"manual_review.reviewed_at"); reviewers=review["reviewers"]
    if not isinstance(reviewers,list) or len(reviewers)<2 or len(set(reviewers))!=len(reviewers) or any(not isinstance(x,str) or not x.strip() for x in reviewers): raise ValueError("manual review requires two distinct reviewers")
    if any(review.get(k) is not True for k in ("authorization_confirmed","authenticated_web_reviewed","authenticated_api_reviewed","business_logic_reviewed")): raise ValueError("manual review assertions must all pass")
    if not isinstance(review["notes"],str) or not review["notes"].strip(): raise ValueError("manual review notes are required")

    runs=profile["runs"]
    if not isinstance(runs,list) or len(runs)!=2 or any(not isinstance(x,dict) or set(x)!=RUN_FIELDS for x in runs) or [x["kind"] for x in runs]!=list(KINDS): raise ValueError("authenticated web and API runs are required")
    evidence=[proof(root,policy_path),proof(root,review_path)]; unaccepted=[]
    for run in runs:
        kind=run["kind"]; plan_path=safe(root,run["plan"],f"{kind}.plan"); plan_proof=proof(root,plan_path)
        receipt_path=safe(root,run["execution_receipt"],f"{kind}.receipt"); receipt=read_json(receipt_path,f"{kind}.receipt")
        if receipt.get("schema")!="elite-official-tool-execution/v1" or receipt.get("project_id")!=profile["project_id"] or receipt.get("environment")!="production" or receipt.get("release_digest")!=profile["release_digest"] or receipt.get("control_id")!="OFFENSIVE_SECURITY" or receipt.get("exit_code")!=0: raise ValueError(f"{kind} execution identity/exit mismatch")
        tool=receipt.get("tool")
        if not isinstance(tool,dict) or set(tool)!={"name","version","source","sha256"} or tool.get("name")!="OWASP ZAP" or tool.get("version")!="2.17.0" or tool.get("source")!=SOURCE or not re.fullmatch(r"[0-9a-f]{64}",str(tool.get("sha256"))): raise ValueError(f"{kind} exact ZAP identity mismatch")
        args=run["expected_arguments"]
        if not isinstance(args,list) or any(not isinstance(x,str) or not x for x in args) or receipt.get("arguments_sha256")!=canonical_sha(args): raise ValueError(f"{kind} arguments are not hash-bound")
        expected_target=dict(target); expected_target.update({"run_kind":kind,"automation_plan_sha256":plan_proof["sha256"]})
        if receipt.get("target")!=expected_target: raise ValueError(f"{kind} target/plan binding mismatch")
        for channel in ("stdout","stderr"): verify_proof(root,receipt.get(channel),f"{kind}.{channel}")
        report_path=safe(root,run["report"],f"{kind}.report"); report=read_json(report_path,f"{kind}.report")
        if str(report.get("@version"))!="2.17.0": raise ValueError(f"{kind} ZAP report version mismatch")
        for alert in alerts(report):
            if alert["riskcode"]>=3 and (alert["pluginid"],alert["riskcode"]) not in accepted_ids: unaccepted.append((kind,alert["pluginid"],alert["riskcode"]))
        evidence.extend([proof(root,receipt_path),plan_proof,proof(root,report_path),receipt["stdout"],receipt["stderr"]])
    if unaccepted: raise ValueError(f"unaccepted high/critical ZAP alerts: {unaccepted}")
    return {"id":"OFFENSIVE_SECURITY","project_id":profile["project_id"],"environment":"production","release_digest":profile["release_digest"],"result":"PASS","executed_at":profile["evaluated_at"],"expires_at":profile["expires_at"],"executor_ref":profile["executor_ref"],"tool":{"name":"OWASP ZAP","version":"2.17.0","source":SOURCE,"digest":"sha256:"+canonical_sha([run["kind"] for run in runs])},"target":target,"assertions":{"zap_automation_pass":True,"authenticated_scope_pass":True,"api_scope_pass":True,"zero_unaccepted_high_critical":True,"manual_review_pass":True},"evidence":evidence}

def main()->int:
    p=argparse.ArgumentParser(); p.add_argument("--project-root",required=True,type=Path); p.add_argument("--profile",required=True); a=p.parse_args()
    try:
        root=a.project_root.resolve(strict=True); profile=read_json(safe(root,a.profile,"profile"),"profile"); output=safe(root,profile.get("output"),"output",exists=False); output.parent.resolve(strict=True).relative_to(root); receipt=validate(profile,root); atomic(output,(json.dumps(receipt,indent=2,sort_keys=True)+"\n").encode()); print(f"ZAP_OFFENSIVE_ADMISSION_PASS project={receipt['project_id']} release={receipt['release_digest']}"); return 0
    except (ValueError,OSError) as error: print(f"ZAP_OFFENSIVE_ADMISSION_FAILED: {error}",file=sys.stderr); return 1
if __name__=="__main__": sys.exit(main())
````

### FILE: `production_admission_gate/test_validate_zap_offensive_admission.py`

```yaml
block_id: "SECURE-OPS:zap-offensive-admission-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local regressions for authorized OWASP ZAP web/API evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "09d9b0fc6b266f11de9f3b8a06da4df8c0448ce549f9dd34bf9f09324f24871c"
variables: []
secrets_allowed: false
```

````python
from __future__ import annotations
import copy,hashlib,json
from datetime import datetime,timedelta,timezone
from pathlib import Path
import tempfile,unittest
from validate_k6_load_admission import canonical_sha
from validate_zap_offensive_admission import KINDS,SOURCE,validate
class ZapAdmissionTests(unittest.TestCase):
 def setUp(self):
  self.temp=tempfile.TemporaryDirectory(); self.root=Path(self.temp.name).resolve(); (self.root/"evidence").mkdir(); now=datetime.now(timezone.utc); z=lambda d:d.isoformat().replace("+00:00","Z"); digest="sha256:"+"e"*64
  self.profile={"schema":"elite-zap-offensive-admission/v1","project_id":"revestex","environment":"production","release_digest":digest,"evaluated_at":z(now),"expires_at":z(now+timedelta(days=7)),"executor_ref":"security:engagement/42","target":{"base_url":"https://app.example.test","api_base_url":"https://api.example.test","authorization_ticket":"SEC-42"},"policy":"evidence/policy.json","manual_review":"evidence/review.json","runs":[],"output":"evidence/offensive-control.json"}
  self.write("evidence/policy.json",{"schema":"elite-offensive-security-policy/v1","project_id":"revestex","environment":"production","release_digest":digest,"approved_at":z(now),"approvers":["security-owner","system-owner"],"accepted_alerts":[]})
  self.write("evidence/review.json",{"schema":"elite-offensive-manual-review/v1","project_id":"revestex","environment":"production","release_digest":digest,"reviewed_at":z(now),"reviewers":["security-reviewer","business-reviewer"],"authorization_confirmed":True,"authenticated_web_reviewed":True,"authenticated_api_reviewed":True,"business_logic_reviewed":True,"notes":"Authorized scope and residual risk reviewed."})
  for kind in KINDS:
   stem=kind.lower(); plan=self.root/"evidence"/f"{stem}.yaml"; plan.write_text("env:\n  contexts: []\n",encoding="utf-8"); args=["-cmd","-autorun",f"evidence/{stem}.yaml"]
   out=self.root/"evidence"/f"{stem}.stdout.bin"; err=self.root/"evidence"/f"{stem}.stderr.bin"; out.write_bytes(b"Automation plan succeeded\n"); err.write_bytes(b""); pf=lambda p:{"path":p.relative_to(self.root).as_posix(),"bytes":p.stat().st_size,"sha256":hashlib.sha256(p.read_bytes()).hexdigest()}; target=dict(self.profile["target"]); target.update({"run_kind":kind,"automation_plan_sha256":pf(plan)["sha256"]})
   receipt={"schema":"elite-official-tool-execution/v1","project_id":"revestex","environment":"production","release_digest":digest,"control_id":"OFFENSIVE_SECURITY","executed_at":z(now),"exit_code":0,"tool":{"name":"OWASP ZAP","version":"2.17.0","source":SOURCE,"sha256":"f"*64},"arguments_sha256":canonical_sha(args),"working_directory":".","target":target,"environment_variable_names":[],"stdout":pf(out),"stderr":pf(err)}; self.write(f"evidence/{stem}-receipt.json",receipt); self.write(f"evidence/{stem}-report.json",{"@version":"2.17.0","site":[{"@name":self.profile["target"]["base_url" if kind=="WEB_AUTHENTICATED" else "api_base_url"],"alerts":[{"pluginid":"10021","riskcode":"1","name":"Low fixture","instances":[]}]}]}); self.profile["runs"].append({"kind":kind,"execution_receipt":f"evidence/{stem}-receipt.json","plan":f"evidence/{stem}.yaml","report":f"evidence/{stem}-report.json","expected_arguments":args})
 def tearDown(self): self.temp.cleanup()
 def write(self,p,v): (self.root/p).write_text(json.dumps(v),encoding="utf-8")
 def test_web_api_and_manual_review_emit_semantic_receipt(self): self.assertTrue(all(validate(self.profile,self.root)["assertions"].values()))
 def test_unauthorized_target_is_rejected(self):
  v=copy.deepcopy(self.profile); v["target"]["authorization_ticket"]=""
  with self.assertRaisesRegex(ValueError,"authorization ticket"): validate(v,self.root)
 def test_high_alert_is_rejected(self):
  p=self.root/self.profile["runs"][0]["report"]; v=json.loads(p.read_text()); v["site"][0]["alerts"][0]["riskcode"]="3"; self.write(self.profile["runs"][0]["report"],v)
  with self.assertRaisesRegex(ValueError,"unaccepted high"): validate(self.profile,self.root)
 def test_expired_exception_is_rejected(self):
  p=self.root/self.profile["policy"]; v=json.loads(p.read_text()); v["accepted_alerts"]=[{"pluginid":"42","riskcode":3,"reason":"temporary","owner":"security","expires_at":"2020-01-01T00:00:00Z"}]; self.write(self.profile["policy"],v)
  with self.assertRaisesRegex(ValueError,"expires before"): validate(self.profile,self.root)
 def test_plan_tamper_is_rejected(self):
  (self.root/self.profile["runs"][0]["plan"]).write_text("tampered")
  with self.assertRaisesRegex(ValueError,"binding mismatch"): validate(self.profile,self.root)
 def test_missing_api_run_is_rejected(self):
  v=copy.deepcopy(self.profile); v["runs"]=v["runs"][:1]
  with self.assertRaisesRegex(ValueError,"web and API"): validate(v,self.root)
 def test_manual_review_false_is_rejected(self):
  p=self.root/self.profile["manual_review"]; v=json.loads(p.read_text()); v["business_logic_reviewed"]=False; self.write(self.profile["manual_review"],v)
  with self.assertRaisesRegex(ValueError,"manual review assertions"): validate(self.profile,self.root)
 def test_nonzero_zap_exit_is_rejected(self):
  p=self.root/self.profile["runs"][1]["execution_receipt"]; v=json.loads(p.read_text()); v["exit_code"]=2; self.write(self.profile["runs"][1]["execution_receipt"],v)
  with self.assertRaisesRegex(ValueError,"identity/exit"): validate(self.profile,self.root)
if __name__=="__main__": unittest.main()
````

### FILE: `production_admission_gate/identity-authorization-admission.template.json`

```yaml
block_id: "SECURE-OPS:identity-authorization-admission-template:v1"
operation: CREATE
provenance: AUTHORED
source: "local configuration contract bound to OpenID Foundation and Microsoft test evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "84adaf6b5c6fc5e49c327b0356a85882c57ec570eb934bd2a7990e7bdcdb32d7"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite-identity-authorization-admission/v1",
  "project_id": "",
  "environment": "production",
  "release_digest": "",
  "evaluated_at": "",
  "expires_at": "",
  "executor_ref": "",
  "tool": {
    "name": "OpenID Foundation Conformance Suite",
    "version": "5.2.4",
    "source": "https://gitlab.com/openid/conformance-suite/-/releases/release-v5.2.4"
  },
  "target": {
    "issuer_url": "",
    "client_id": "",
    "audience": "",
    "conformance_profile": "",
    "application_url": "",
    "tenant_ids": [],
    "role_ids": []
  },
  "policy": "",
  "runs": [
    {"kind":"OIDC_CONFORMANCE", "execution_receipt":"", "observation":"", "expected_arguments":[]},
    {"kind":"ROLE_POSITIVE_JOURNEYS", "execution_receipt":"", "observation":"", "expected_arguments":[]},
    {"kind":"UNAUTHORIZED_NEGATIVE_JOURNEYS", "execution_receipt":"", "observation":"", "expected_arguments":[]},
    {"kind":"TENANT_ISOLATION", "execution_receipt":"", "observation":"", "expected_arguments":[]},
    {"kind":"SESSION_REVOCATION_ROTATION", "execution_receipt":"", "observation":"", "expected_arguments":[]}
  ],
  "output": ""
}
````

### FILE: `production_admission_gate/validate_identity_authorization_admission.py`

```yaml
block_id: "SECURE-OPS:identity-authorization-admission-validator:v1"
operation: CREATE
provenance: AUTHORED
source: "local fail-closed semantic validator for exact OpenID conformance and project authorization evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "fced616670e10cff96e42d4cb55df93fc4be8f7dad9a57290a1a8b0bc8328cb7"
variables: []
secrets_allowed: false
```

````python
from __future__ import annotations
import argparse,json,re,sys
from itertools import permutations
from pathlib import Path
from urllib.parse import urlparse
from validate_k6_load_admission import atomic,canonical_sha,proof,read_json,safe,utc,verify_proof

KINDS=("OIDC_CONFORMANCE","ROLE_POSITIVE_JOURNEYS","UNAUTHORIZED_NEGATIVE_JOURNEYS","TENANT_ISOLATION","SESSION_REVOCATION_ROTATION")
FIELDS={"schema","project_id","environment","release_digest","evaluated_at","expires_at","executor_ref","tool","target","policy","runs","output"}
RUN_FIELDS={"kind","execution_receipt","observation","expected_arguments"}
TARGET_FIELDS={"issuer_url","client_id","audience","conformance_profile","application_url","tenant_ids","role_ids"}
OIDF_SOURCE="https://gitlab.com/openid/conformance-suite/-/releases/release-v5.2.4"
PLAYWRIGHT_SOURCE="https://github.com/microsoft/playwright/releases/tag/v1.62.1"

def nonempty(value:object,label:str)->str:
 if not isinstance(value,str) or not value.strip(): raise ValueError(f"{label} is required")
 return value

def https(value:object,label:str)->None:
 value=nonempty(value,label); parsed=urlparse(value)
 if parsed.scheme!="https" or not parsed.hostname or parsed.username or parsed.password or parsed.fragment: raise ValueError(f"{label} must be credential-free HTTPS")

def unique_strings(value:object,label:str,minimum:int=1)->list[str]:
 if not isinstance(value,list) or len(value)<minimum or len(set(value))!=len(value) or any(not isinstance(x,str) or not x.strip() for x in value): raise ValueError(f"{label} must contain at least {minimum} unique non-empty values")
 return value

def bind(value:dict[str,object],profile:dict[str,object],label:str)->None:
 for key in ("project_id","environment","release_digest"):
  if value.get(key)!=profile[key]: raise ValueError(f"{label} {key} mismatch")

def validate_oidc(value:dict[str,object],profile:dict[str,object])->None:
 fields={"schema","project_id","environment","release_digest","captured_at","kind","issuer_url","client_id","conformance_profile","suite_release","plan_id","summary","non_passed"}
 if set(value)!=fields or value.get("schema")!="elite-identity-observation/v1" or value.get("kind")!="OIDC_CONFORMANCE": raise ValueError("OIDC conformance observation fields/schema mismatch")
 bind(value,profile,"OIDC conformance"); utc(value["captured_at"],"OIDC captured_at"); target=profile["target"]
 if value["issuer_url"]!=target["issuer_url"] or value["client_id"]!=target["client_id"] or value["conformance_profile"]!=target["conformance_profile"] or value["suite_release"]!="release-v5.2.4" or not nonempty(value["plan_id"],"OIDC plan_id"): raise ValueError("OIDC exact target/profile/release mismatch")
 summary=value["summary"]
 if not isinstance(summary,dict) or set(summary)!={"passed","review","warning","skipped","failed","interrupted","total"} or any(not isinstance(v,int) or v<0 for v in summary.values()): raise ValueError("OIDC summary must contain exact non-negative counters")
 if summary["total"]<=0 or summary["total"]!=sum(summary[k] for k in ("passed","review","warning","skipped","failed","interrupted")) or summary["failed"] or summary["interrupted"]: raise ValueError("OIDC conformance contains failed/interrupted or inconsistent results")
 exceptions=value["non_passed"]
 if not isinstance(exceptions,list) or len(exceptions)!=summary["review"]+summary["warning"]+summary["skipped"]: raise ValueError("every non-passed OIDC result requires an approved disposition")
 seen=set()
 for item in exceptions:
  if not isinstance(item,dict) or set(item)!={"test_id","status","reason","reviewers"} or item["status"] not in ("REVIEW","WARNING","SKIPPED") or not nonempty(item["test_id"],"OIDC test_id") or item["test_id"] in seen or not nonempty(item["reason"],"OIDC disposition reason") or len(unique_strings(item["reviewers"],"OIDC disposition reviewers",2))<2: raise ValueError("OIDC non-passed disposition is invalid")
  seen.add(item["test_id"])

def validate_roles(value:dict[str,object],profile:dict[str,object])->None:
 fields={"schema","project_id","environment","release_digest","captured_at","kind","application_url","journeys"}
 if set(value)!=fields or value.get("schema")!="elite-identity-observation/v1" or value.get("kind")!="ROLE_POSITIVE_JOURNEYS": raise ValueError("role observation fields/schema mismatch")
 bind(value,profile,"roles"); utc(value["captured_at"],"roles captured_at")
 if value["application_url"]!=profile["target"]["application_url"]: raise ValueError("role target mismatch")
 journeys=value["journeys"]
 if not isinstance(journeys,list) or {x.get("role_id") for x in journeys if isinstance(x,dict)}!=set(profile["target"]["role_ids"]): raise ValueError("every declared role requires a positive journey")
 for item in journeys:
  if set(item)!={"role_id","tenant_id","journey_id","authenticated","status","expected_backend_effect","effect_observed"} or item["tenant_id"] not in profile["target"]["tenant_ids"] or not nonempty(item["journey_id"],"role journey_id") or item["authenticated"] is not True or not isinstance(item["status"],int) or not 200<=item["status"]<300 or item["expected_backend_effect"] is not True or item["effect_observed"] is not True: raise ValueError("positive role journey did not prove its backend effect")

def validate_unauthorized(value:dict[str,object],profile:dict[str,object])->None:
 required={"unauthenticated","expired_token","wrong_issuer","wrong_audience","insufficient_role","object_ownership"}
 fields={"schema","project_id","environment","release_digest","captured_at","kind","application_url","scenarios"}
 if set(value)!=fields or value.get("schema")!="elite-identity-observation/v1" or value.get("kind")!="UNAUTHORIZED_NEGATIVE_JOURNEYS": raise ValueError("unauthorized observation fields/schema mismatch")
 bind(value,profile,"unauthorized"); utc(value["captured_at"],"unauthorized captured_at")
 if value["application_url"]!=profile["target"]["application_url"]: raise ValueError("unauthorized target mismatch")
 scenarios=value["scenarios"]
 if not isinstance(scenarios,list) or {x.get("scenario") for x in scenarios if isinstance(x,dict)}!=required: raise ValueError("all mandatory unauthorized scenarios are required")
 for item in scenarios:
  if set(item)!={"scenario","status","data_disclosed","mutation_observed"} or item["status"] not in (401,403,404) or item["data_disclosed"] is not False or item["mutation_observed"] is not False: raise ValueError("unauthorized journey leaked data or changed state")

def validate_tenants(value:dict[str,object],profile:dict[str,object])->None:
 fields={"schema","project_id","environment","release_digest","captured_at","kind","tenant_ids","scenarios"}
 if set(value)!=fields or value.get("schema")!="elite-identity-observation/v1" or value.get("kind")!="TENANT_ISOLATION": raise ValueError("tenant observation fields/schema mismatch")
 bind(value,profile,"tenants"); utc(value["captured_at"],"tenants captured_at"); tenants=profile["target"]["tenant_ids"]
 if value["tenant_ids"]!=tenants: raise ValueError("tenant set/order mismatch")
 scenarios=value["scenarios"]; expected={(a,b,action) for a,b in permutations(tenants,2) for action in ("read","write","list")}
 if not isinstance(scenarios,list) or len(scenarios)!=len(expected): raise ValueError("complete ordered cross-tenant matrix is required")
 actual=set()
 for item in scenarios:
  if not isinstance(item,dict) or set(item)!={"actor_tenant","target_tenant","action","resource_key","status","data_disclosed","mutation_observed"} or not nonempty(item["resource_key"],"tenant resource_key"): raise ValueError("tenant scenario fields are invalid")
  key=(item["actor_tenant"],item["target_tenant"],item["action"]); actual.add(key)
  if item["actor_tenant"]==item["target_tenant"] or item["status"] not in (403,404) or item["data_disclosed"] is not False or item["mutation_observed"] is not False: raise ValueError("cross-tenant isolation failed")
 if actual!=expected: raise ValueError("complete ordered cross-tenant matrix is required")

def validate_sessions(value:dict[str,object],profile:dict[str,object])->None:
 required={"logout_revocation","token_revocation","signing_key_rotation","session_expiry","replay_rejection"}
 fields={"schema","project_id","environment","release_digest","captured_at","kind","issuer_url","cases"}
 if set(value)!=fields or value.get("schema")!="elite-identity-observation/v1" or value.get("kind")!="SESSION_REVOCATION_ROTATION": raise ValueError("session observation fields/schema mismatch")
 bind(value,profile,"sessions"); utc(value["captured_at"],"sessions captured_at")
 if value["issuer_url"]!=profile["target"]["issuer_url"]: raise ValueError("session issuer mismatch")
 cases=value["cases"]
 if not isinstance(cases,list) or {x.get("case") for x in cases if isinstance(x,dict)}!=required: raise ValueError("all revocation/rotation/session cases are required")
 for item in cases:
  if set(item)!={"case","status","old_credential_accepted","rejected","audit_event_observed"} or item["status"] not in (401,403) or item["old_credential_accepted"] is not False or item["rejected"] is not True or item["audit_event_observed"] is not True: raise ValueError("session revocation/rotation assertion failed")

def validate(profile:object,root:Path)->dict[str,object]:
 if not isinstance(profile,dict) or set(profile)!=FIELDS: raise ValueError("profile fields must be exact")
 if profile["schema"]!="elite-identity-authorization-admission/v1" or profile["environment"]!="production" or not nonempty(profile["project_id"],"project_id") or not re.fullmatch(r"sha256:[0-9a-f]{64}",str(profile["release_digest"])): raise ValueError("profile identity is invalid")
 start=utc(profile["evaluated_at"],"evaluated_at"); end=utc(profile["expires_at"],"expires_at")
 if not start<end or (end-start).total_seconds()>30*86400: raise ValueError("admission validity must be positive and at most 30 days")
 nonempty(profile["executor_ref"],"executor_ref")
 if profile["tool"]!={"name":"OpenID Foundation Conformance Suite","version":"5.2.4","source":OIDF_SOURCE}: raise ValueError("exact OpenID Foundation Conformance Suite 5.2.4 identity is required")
 target=profile["target"]
 if not isinstance(target,dict) or set(target)!=TARGET_FIELDS: raise ValueError("target fields must be exact")
 https(target["issuer_url"],"target.issuer_url"); https(target["application_url"],"target.application_url")
 for key in ("client_id","audience","conformance_profile"): nonempty(target[key],f"target.{key}")
 unique_strings(target["tenant_ids"],"target.tenant_ids",2); unique_strings(target["role_ids"],"target.role_ids",1)
 policy_path=safe(root,profile["policy"],"policy"); policy=read_json(policy_path,"policy")
 policy_fields={"schema","project_id","environment","release_digest","approved_at","approvers","issuer_url","client_id","audience","conformance_profile","tenant_ids","role_ids","break_glass"}
 if not isinstance(policy,dict) or set(policy)!=policy_fields or policy.get("schema")!="elite-identity-authorization-policy/v1": raise ValueError("identity policy fields/schema must be exact")
 bind(policy,profile,"policy"); utc(policy["approved_at"],"policy.approved_at"); unique_strings(policy["approvers"],"policy.approvers",2)
 for key in ("issuer_url","client_id","audience","conformance_profile","tenant_ids","role_ids"):
  if policy[key]!=target[key]: raise ValueError(f"identity policy {key} mismatch")
 bg=policy["break_glass"]
 if not isinstance(bg,dict) or set(bg)!={"configured","max_minutes","two_person_approval","audit_required","post_use_review_required"} or bg["configured"] is not True or not isinstance(bg["max_minutes"],int) or not 1<=bg["max_minutes"]<=60 or any(bg[k] is not True for k in ("two_person_approval","audit_required","post_use_review_required")): raise ValueError("break-glass governance is incomplete")
 runs=profile["runs"]
 if not isinstance(runs,list) or len(runs)!=5 or any(not isinstance(x,dict) or set(x)!=RUN_FIELDS for x in runs) or [x["kind"] for x in runs]!=list(KINDS): raise ValueError("five ordered identity runs are required")
 evidence=[proof(root,policy_path)]; tool_hashes={}
 validators={"OIDC_CONFORMANCE":validate_oidc,"ROLE_POSITIVE_JOURNEYS":validate_roles,"UNAUTHORIZED_NEGATIVE_JOURNEYS":validate_unauthorized,"TENANT_ISOLATION":validate_tenants,"SESSION_REVOCATION_ROTATION":validate_sessions}
 for run in runs:
  kind=run["kind"]; observation_path=safe(root,run["observation"],f"{kind}.observation"); observation=read_json(observation_path,f"{kind}.observation"); validators[kind](observation,profile); observation_proof=proof(root,observation_path)
  receipt_path=safe(root,run["execution_receipt"],f"{kind}.receipt"); receipt=read_json(receipt_path,f"{kind}.receipt")
  if receipt.get("schema")!="elite-official-tool-execution/v1" or receipt.get("project_id")!=profile["project_id"] or receipt.get("environment")!="production" or receipt.get("release_digest")!=profile["release_digest"] or receipt.get("control_id")!="IDENTITY_AUTHORIZATION" or receipt.get("exit_code")!=0: raise ValueError(f"{kind} execution identity/exit mismatch")
  tool=receipt.get("tool"); expected=("OpenID Foundation Conformance Runner","5.2.4",OIDF_SOURCE) if kind=="OIDC_CONFORMANCE" else ("Microsoft Playwright","1.62.1",PLAYWRIGHT_SOURCE)
  if not isinstance(tool,dict) or set(tool)!={"name","version","source","sha256"} or (tool.get("name"),tool.get("version"),tool.get("source"))!=expected or not re.fullmatch(r"[0-9a-f]{64}",str(tool.get("sha256"))): raise ValueError(f"{kind} exact official tool identity mismatch")
  family="oidf" if kind=="OIDC_CONFORMANCE" else "playwright"; tool_hashes.setdefault(family,tool["sha256"])
  if tool_hashes[family]!=tool["sha256"]: raise ValueError(f"all {family} runs must use the same binary/source hash")
  args=run["expected_arguments"]
  if not isinstance(args,list) or not args or any(not isinstance(x,str) or not x for x in args) or receipt.get("arguments_sha256")!=canonical_sha(args): raise ValueError(f"{kind} arguments are not hash-bound")
  expected_target=dict(target); expected_target.update({"run_kind":kind,"observation_sha256":observation_proof["sha256"]})
  if receipt.get("target")!=expected_target: raise ValueError(f"{kind} target/observation binding mismatch")
  if not isinstance(receipt.get("environment_variable_names"),list) or any(re.search(r"(?:secret|password|token|key|credential)",str(x),re.I) and not str(x).endswith("_REF") for x in receipt["environment_variable_names"]): raise ValueError(f"{kind} environment must expose references, not secret-bearing names")
  for channel in ("stdout","stderr"): verify_proof(root,receipt.get(channel),f"{kind}.{channel}")
  evidence.extend([proof(root,receipt_path),observation_proof,receipt["stdout"],receipt["stderr"]])
 return {"id":"IDENTITY_AUTHORIZATION","project_id":profile["project_id"],"environment":"production","release_digest":profile["release_digest"],"result":"PASS","executed_at":profile["evaluated_at"],"expires_at":profile["expires_at"],"executor_ref":profile["executor_ref"],"tool":{"name":"OpenID Foundation Conformance Suite + Microsoft Playwright","version":"5.2.4 + 1.62.1","source":OIDF_SOURCE,"digest":"sha256:"+canonical_sha([tool_hashes["oidf"],tool_hashes["playwright"]])},"target":target,"assertions":{"oidc_conformance_pass":True,"role_positive_journeys_pass":True,"unauthorized_negative_journeys_pass":True,"tenant_isolation_pass":True,"session_revocation_rotation_pass":True},"evidence":evidence}

def main()->int:
 p=argparse.ArgumentParser(); p.add_argument("--project-root",required=True,type=Path); p.add_argument("--profile",required=True); a=p.parse_args()
 try:
  root=a.project_root.resolve(strict=True); profile=read_json(safe(root,a.profile,"profile"),"profile"); output=safe(root,profile.get("output"),"output",exists=False); output.parent.resolve(strict=True).relative_to(root); receipt=validate(profile,root); atomic(output,(json.dumps(receipt,indent=2,sort_keys=True)+"\n").encode()); print(f"IDENTITY_AUTHORIZATION_ADMISSION_PASS project={receipt['project_id']} release={receipt['release_digest']}"); return 0
 except (ValueError,OSError) as error: print(f"IDENTITY_AUTHORIZATION_ADMISSION_FAILED: {error}",file=sys.stderr); return 1
if __name__=="__main__": sys.exit(main())
````

### FILE: `production_admission_gate/test_validate_identity_authorization_admission.py`

```yaml
block_id: "SECURE-OPS:identity-authorization-admission-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local regressions for OpenID, roles, unauthorized, cross-tenant, session and break-glass evidence"
license: "LicenseRef-Workspace-Owner"
sha256: "8e6e878af1c2f84d6294556c6b4fdfb97e9f07093ab667d4c01e415c0af2a417"
variables: []
secrets_allowed: false
```

````python
from __future__ import annotations
import copy,hashlib,json,tempfile,unittest
from datetime import datetime,timedelta,timezone
from itertools import permutations
from pathlib import Path
from validate_identity_authorization_admission import KINDS,OIDF_SOURCE,PLAYWRIGHT_SOURCE,validate
from validate_k6_load_admission import canonical_sha

class IdentityAuthorizationAdmissionTests(unittest.TestCase):
 def setUp(self):
  self.temp=tempfile.TemporaryDirectory(); self.root=Path(self.temp.name); (self.root/"evidence").mkdir(); now=datetime.now(timezone.utc); z=lambda d:d.isoformat().replace("+00:00","Z"); digest="sha256:"+"e"*64
  self.target={"issuer_url":"https://id.example.test/tenant","client_id":"web-client","audience":"revestex-api","conformance_profile":"oidcc-basic-certification-test-plan","application_url":"https://app.example.test","tenant_ids":["north","south"],"role_ids":["admin","seller"]}
  self.profile={"schema":"elite-identity-authorization-admission/v1","project_id":"revestex","environment":"production","release_digest":digest,"evaluated_at":z(now),"expires_at":z(now+timedelta(days=7)),"executor_ref":"identity:change/42","tool":{"name":"OpenID Foundation Conformance Suite","version":"5.2.4","source":OIDF_SOURCE},"target":self.target,"policy":"evidence/policy.json","runs":[],"output":"evidence/identity-control.json"}
  self.write("evidence/policy.json",{"schema":"elite-identity-authorization-policy/v1","project_id":"revestex","environment":"production","release_digest":digest,"approved_at":z(now),"approvers":["security-owner","business-owner"],"issuer_url":self.target["issuer_url"],"client_id":"web-client","audience":"revestex-api","conformance_profile":self.target["conformance_profile"],"tenant_ids":["north","south"],"role_ids":["admin","seller"],"break_glass":{"configured":True,"max_minutes":30,"two_person_approval":True,"audit_required":True,"post_use_review_required":True}})
  observations={
   "OIDC_CONFORMANCE":{"schema":"elite-identity-observation/v1","project_id":"revestex","environment":"production","release_digest":digest,"captured_at":z(now),"kind":"OIDC_CONFORMANCE","issuer_url":self.target["issuer_url"],"client_id":"web-client","conformance_profile":self.target["conformance_profile"],"suite_release":"release-v5.2.4","plan_id":"plan-42","summary":{"passed":3,"review":0,"warning":0,"skipped":1,"failed":0,"interrupted":0,"total":4},"non_passed":[{"test_id":"optional-encryption","status":"SKIPPED","reason":"feature excluded from approved profile","reviewers":["security-owner","identity-owner"]}]},
   "ROLE_POSITIVE_JOURNEYS":{"schema":"elite-identity-observation/v1","project_id":"revestex","environment":"production","release_digest":digest,"captured_at":z(now),"kind":"ROLE_POSITIVE_JOURNEYS","application_url":self.target["application_url"],"journeys":[{"role_id":role,"tenant_id":"north","journey_id":f"{role}-write","authenticated":True,"status":200,"expected_backend_effect":True,"effect_observed":True} for role in self.target["role_ids"]]},
   "UNAUTHORIZED_NEGATIVE_JOURNEYS":{"schema":"elite-identity-observation/v1","project_id":"revestex","environment":"production","release_digest":digest,"captured_at":z(now),"kind":"UNAUTHORIZED_NEGATIVE_JOURNEYS","application_url":self.target["application_url"],"scenarios":[{"scenario":case,"status":401 if case in ("unauthenticated","expired_token","wrong_issuer","wrong_audience") else 403,"data_disclosed":False,"mutation_observed":False} for case in ("unauthenticated","expired_token","wrong_issuer","wrong_audience","insufficient_role","object_ownership")]},
   "TENANT_ISOLATION":{"schema":"elite-identity-observation/v1","project_id":"revestex","environment":"production","release_digest":digest,"captured_at":z(now),"kind":"TENANT_ISOLATION","tenant_ids":["north","south"],"scenarios":[{"actor_tenant":a,"target_tenant":b,"action":action,"resource_key":"shared-id-42","status":403,"data_disclosed":False,"mutation_observed":False} for a,b in permutations(self.target["tenant_ids"],2) for action in ("read","write","list")]},
   "SESSION_REVOCATION_ROTATION":{"schema":"elite-identity-observation/v1","project_id":"revestex","environment":"production","release_digest":digest,"captured_at":z(now),"kind":"SESSION_REVOCATION_ROTATION","issuer_url":self.target["issuer_url"],"cases":[{"case":case,"status":401,"old_credential_accepted":False,"rejected":True,"audit_event_observed":True} for case in ("logout_revocation","token_revocation","signing_key_rotation","session_expiry","replay_rejection")]}}
  for kind in KINDS:
   stem=kind.lower(); obs=self.write(f"evidence/{stem}.json",observations[kind]); args=["--plan","plan-42"] if kind=="OIDC_CONFORMANCE" else ["test",stem]; out=self.file(f"evidence/{stem}.stdout",b"pass\n"); err=self.file(f"evidence/{stem}.stderr",b""); tool={"name":"OpenID Foundation Conformance Runner","version":"5.2.4","source":OIDF_SOURCE,"sha256":"a"*64} if kind=="OIDC_CONFORMANCE" else {"name":"Microsoft Playwright","version":"1.62.1","source":PLAYWRIGHT_SOURCE,"sha256":"b"*64}; target=dict(self.target); target.update({"run_kind":kind,"observation_sha256":self.pf(obs)["sha256"]})
   receipt={"schema":"elite-official-tool-execution/v1","project_id":"revestex","environment":"production","release_digest":digest,"control_id":"IDENTITY_AUTHORIZATION","executed_at":z(now),"exit_code":0,"tool":tool,"arguments_sha256":canonical_sha(args),"working_directory":".","target":target,"environment_variable_names":["OIDC_CLIENT_SECRET_REF"],"stdout":self.pf(out),"stderr":self.pf(err)}; self.write(f"evidence/{stem}-receipt.json",receipt); self.profile["runs"].append({"kind":kind,"execution_receipt":f"evidence/{stem}-receipt.json","observation":f"evidence/{stem}.json","expected_arguments":args})
 def tearDown(self): self.temp.cleanup()
 def file(self,p,data): q=self.root/p; q.write_bytes(data); return q
 def write(self,p,v): q=self.root/p; q.write_text(json.dumps(v),encoding="utf-8"); return q
 def pf(self,p): data=p.read_bytes(); return {"path":p.relative_to(self.root).as_posix(),"bytes":len(data),"sha256":hashlib.sha256(data).hexdigest()}
 def mutate(self,index,change): p=self.root/self.profile["runs"][index]["observation"]; v=json.loads(p.read_text()); change(v); self.write(self.profile["runs"][index]["observation"],v)
 def test_complete_identity_evidence_emits_semantic_receipt(self): self.assertTrue(all(validate(self.profile,self.root)["assertions"].values()))
 def test_failed_oidc_case_is_rejected(self):
  self.mutate(0,lambda v:v["summary"].update(passed=2,failed=1))
  with self.assertRaisesRegex(ValueError,"failed/interrupted"): validate(self.profile,self.root)
 def test_unreviewed_oidc_warning_is_rejected(self):
  self.mutate(0,lambda v:(v["summary"].update(passed=2,warning=1),v.update(non_passed=v["non_passed"])))
  with self.assertRaisesRegex(ValueError,"every non-passed"): validate(self.profile,self.root)
 def test_missing_role_journey_is_rejected(self):
  self.mutate(1,lambda v:v.update(journeys=v["journeys"][:1]))
  with self.assertRaisesRegex(ValueError,"every declared role"): validate(self.profile,self.root)
 def test_unauthorized_mutation_is_rejected(self):
  self.mutate(2,lambda v:v["scenarios"][0].update(mutation_observed=True))
  with self.assertRaisesRegex(ValueError,"changed state"): validate(self.profile,self.root)
 def test_incomplete_tenant_matrix_is_rejected(self):
  self.mutate(3,lambda v:v.update(scenarios=v["scenarios"][:-1]))
  with self.assertRaisesRegex(ValueError,"complete ordered"): validate(self.profile,self.root)
 def test_cross_tenant_disclosure_is_rejected(self):
  self.mutate(3,lambda v:v["scenarios"][0].update(data_disclosed=True))
  with self.assertRaisesRegex(ValueError,"isolation failed"): validate(self.profile,self.root)
 def test_old_session_acceptance_is_rejected(self):
  self.mutate(4,lambda v:v["cases"][0].update(old_credential_accepted=True,rejected=False,status=200))
  with self.assertRaisesRegex(ValueError,"assertion failed"): validate(self.profile,self.root)
 def test_nonzero_official_tool_exit_is_rejected(self):
  p=self.root/self.profile["runs"][0]["execution_receipt"]; v=json.loads(p.read_text()); v["exit_code"]=2; self.write(self.profile["runs"][0]["execution_receipt"],v)
  with self.assertRaisesRegex(ValueError,"identity/exit"): validate(self.profile,self.root)
 def test_observation_tamper_is_rejected_by_receipt_binding(self):
  self.mutate(1,lambda v:v["journeys"][0].update(journey_id="renamed"))
  with self.assertRaisesRegex(ValueError,"target/observation binding"): validate(self.profile,self.root)
 def test_secret_bearing_environment_name_is_rejected(self):
  p=self.root/self.profile["runs"][1]["execution_receipt"]; v=json.loads(p.read_text()); v["environment_variable_names"]=["OIDC_CLIENT_SECRET"]; self.write(self.profile["runs"][1]["execution_receipt"],v)
  with self.assertRaisesRegex(ValueError,"references, not secret"): validate(self.profile,self.root)
 def test_incomplete_break_glass_is_rejected(self):
  p=self.root/self.profile["policy"]; v=json.loads(p.read_text()); v["break_glass"]["two_person_approval"]=False; self.write(self.profile["policy"],v)
  with self.assertRaisesRegex(ValueError,"break-glass"): validate(self.profile,self.root)

if __name__=="__main__": unittest.main()
````

### FILE: `production_admission_gate/deploy-rollback-admission.template.json`

```yaml
block_id: "SECURE-OPS:deploy-rollback-admission-template:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://production-admission/kubernetes-deploy-rollback-template"
license: "LicenseRef-Workspace-Owner"
sha256: "82ed615f9701521b1fbdf704e805db8a2638c781b618c0d8371a34b789f5ef62"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite-deploy-rollback-admission/v1",
  "project_id": "",
  "environment": "production",
  "release_digest": "",
  "evaluated_at": "",
  "expires_at": "",
  "executor_ref": "",
  "tool": {
    "name": "Kubernetes kubectl",
    "version": "1.37.0",
    "source": "https://github.com/kubernetes/kubernetes/releases/tag/v1.37.0",
    "sha256": "4721b614a67bb4932a0369e61f4a323d8c6ca00943d3a2ff14837c124da06f0e"
  },
  "target": {},
  "policy": "",
  "runs": [],
  "output": ""
}
````

### FILE: `production_admission_gate/validate_deploy_rollback_admission.py`

```yaml
block_id: "SECURE-OPS:deploy-rollback-admission-validator:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://production-admission/kubernetes-deploy-rollback-validator"
license: "LicenseRef-Workspace-Owner"
sha256: "650c61f3130815267b70afeea3db05ae280822338e449f3a3b27c1ba0cd664ac"
variables: []
secrets_allowed: false
```

````python
from __future__ import annotations
import argparse,json,re,sys
from pathlib import Path
from urllib.parse import urlparse
from validate_k6_load_admission import atomic,canonical_sha,proof,read_json,safe,utc,verify_proof

KINDS=("PREFLIGHT","APPLY_IMMUTABLE_DIGEST","CANARY_HEALTH","ROLLOUT_STATUS","ROLLBACK","POST_ROLLBACK_VERIFY","MIGRATION_COMPATIBILITY")
FIELDS={"schema","project_id","environment","release_digest","evaluated_at","expires_at","executor_ref","tool","target","policy","runs","output"}
RUN_FIELDS={"kind","execution_receipt","observation","expected_arguments"}
TARGET_FIELDS={"cluster_uid_sha256","context","namespace","workload_kind","workload_name","container_name","service_url"}
KUBECTL_SOURCE="https://github.com/kubernetes/kubernetes/releases/tag/v1.37.0"
KUBECTL_SHA="4721b614a67bb4932a0369e61f4a323d8c6ca00943d3a2ff14837c124da06f0e"

def nonempty(value:object,label:str)->str:
 if not isinstance(value,str) or not value.strip(): raise ValueError(f"{label} is required")
 return value

def https(value:object,label:str)->None:
 value=nonempty(value,label); parsed=urlparse(value)
 if parsed.scheme!="https" or not parsed.hostname or parsed.username or parsed.password or parsed.fragment: raise ValueError(f"{label} must be credential-free HTTPS")

def bind(value:dict[str,object],profile:dict[str,object],label:str)->None:
 for key in ("project_id","environment","release_digest"):
  if value.get(key)!=profile[key]: raise ValueError(f"{label} {key} mismatch")

def base_observation(value:object,profile:dict[str,object],kind:str,fields:set[str])->dict[str,object]:
 required={"schema","project_id","environment","release_digest","captured_at","kind"}|fields
 if not isinstance(value,dict) or set(value)!=required or value.get("schema")!="elite-deploy-rollback-observation/v1" or value.get("kind")!=kind: raise ValueError(f"{kind} observation fields/schema mismatch")
 bind(value,profile,kind); utc(value["captured_at"],f"{kind}.captured_at")
 return value

def validate_observation(value:object,profile:dict[str,object],policy:dict[str,object],kind:str)->None:
 target=profile["target"]
 if kind=="PREFLIGHT":
  v=base_observation(value,profile,kind,{"client_version","server_version","cluster_uid_sha256","context","namespace","authorized_verbs","workload_found"})
  if v["client_version"]!="v1.37.0" or not re.fullmatch(r"v1\.(3[4-7])\.[0-9]+",str(v["server_version"])): raise ValueError("kubectl/client-server version or supported skew mismatch")
  if v["cluster_uid_sha256"]!=target["cluster_uid_sha256"] or v["context"]!=target["context"] or v["namespace"]!=target["namespace"]: raise ValueError("preflight target mismatch")
  if v["authorized_verbs"]!=["get","list","patch","watch"] or v["workload_found"] is not True: raise ValueError("preflight authorization/workload check failed")
 elif kind=="APPLY_IMMUTABLE_DIGEST":
  v=base_observation(value,profile,kind,{"workload_kind","workload_name","container_name","requested_image_digest","observed_image_digest","revision_before","revision_after","generation","observed_generation"})
  if (v["workload_kind"],v["workload_name"],v["container_name"])!=(target["workload_kind"],target["workload_name"],target["container_name"]): raise ValueError("apply target mismatch")
  if v["requested_image_digest"]!=policy["desired_image_digest"] or v["observed_image_digest"]!=policy["desired_image_digest"]: raise ValueError("immutable desired image digest was not observed")
  if not all(isinstance(v[x],int) and v[x]>=1 for x in ("revision_before","revision_after","generation","observed_generation")) or v["revision_after"]<=v["revision_before"] or v["generation"]!=v["observed_generation"]: raise ValueError("apply revision/generation did not converge")
 elif kind=="CANARY_HEALTH":
  v=base_observation(value,profile,kind,{"service_url","samples","healthy_samples","error_rate","p95_latency_ms","probes_pass","alerts_firing","error_budget_breached"})
  if v["service_url"]!=target["service_url"] or not isinstance(v["samples"],int) or v["samples"]<policy["canary"]["minimum_samples"] or v["healthy_samples"]!=v["samples"]: raise ValueError("canary sample policy failed")
  if not isinstance(v["error_rate"],(int,float)) or v["error_rate"]>policy["canary"]["max_error_rate"] or not isinstance(v["p95_latency_ms"],(int,float)) or v["p95_latency_ms"]>policy["canary"]["max_p95_latency_ms"] or v["probes_pass"] is not True or v["alerts_firing"]!=[] or v["error_budget_breached"] is not False: raise ValueError("canary health/SLO policy failed")
 elif kind=="ROLLOUT_STATUS":
  v=base_observation(value,profile,kind,{"workload_kind","workload_name","progressing","progressing_reason","available","desired_replicas","updated_replicas","available_replicas","old_replicas","application_probes_pass"})
  if (v["workload_kind"],v["workload_name"])!=(target["workload_kind"],target["workload_name"]): raise ValueError("rollout target mismatch")
  if v["progressing"] is not True or v["progressing_reason"]!="NewReplicaSetAvailable" or v["available"] is not True or not isinstance(v["desired_replicas"],int) or v["desired_replicas"]<1 or v["updated_replicas"]!=v["desired_replicas"] or v["available_replicas"]!=v["desired_replicas"] or v["old_replicas"]!=0 or v["application_probes_pass"] is not True: raise ValueError("rollout status/application probes failed")
 elif kind=="ROLLBACK":
  v=base_observation(value,profile,kind,{"workload_kind","workload_name","from_revision","to_revision","rollback_command_issued","rollout_complete"})
  if (v["workload_kind"],v["workload_name"])!=(target["workload_kind"],target["workload_name"]) or v["from_revision"]!=policy["desired_revision"] or v["to_revision"]!=policy["previous_revision"] or v["rollback_command_issued"] is not True or v["rollout_complete"] is not True: raise ValueError("actual rollback execution was not proven")
 elif kind=="POST_ROLLBACK_VERIFY":
  v=base_observation(value,profile,kind,{"observed_image_digest","desired_replicas","available_replicas","service_health_pass","data_invariants_pass","alerts_firing"})
  if v["observed_image_digest"]!=policy["previous_image_digest"] or not isinstance(v["desired_replicas"],int) or v["desired_replicas"]<1 or v["available_replicas"]!=v["desired_replicas"] or v["service_health_pass"] is not True or v["data_invariants_pass"] is not True or v["alerts_firing"]!=[]: raise ValueError("previous digest/health was not restored")
 else:
  v=base_observation(value,profile,kind,{"migration_id","expand_contract_used","forward_compatible","backward_compatible","old_new_versions_coexisted","rollback_data_safe","integrity_checks_pass","destructive_change_before_convergence"})
  if v["migration_id"]!=policy["migration"]["migration_id"] or any(v[x] is not True for x in ("expand_contract_used","forward_compatible","backward_compatible","old_new_versions_coexisted","rollback_data_safe","integrity_checks_pass")) or v["destructive_change_before_convergence"] is not False: raise ValueError("migration compatibility/rollback safety failed")

def validate(profile:object,root:Path)->dict[str,object]:
 if not isinstance(profile,dict) or set(profile)!=FIELDS: raise ValueError("profile fields must be exact")
 if profile["schema"]!="elite-deploy-rollback-admission/v1" or profile["environment"]!="production" or not nonempty(profile["project_id"],"project_id") or not re.fullmatch(r"sha256:[0-9a-f]{64}",str(profile["release_digest"])) or set(str(profile["release_digest"]).split(":")[-1])=={"0"}: raise ValueError("profile identity is invalid")
 start=utc(profile["evaluated_at"],"evaluated_at"); end=utc(profile["expires_at"],"expires_at")
 if not start<end or (end-start).total_seconds()>30*86400: raise ValueError("admission validity must be positive and at most 30 days")
 nonempty(profile["executor_ref"],"executor_ref")
 if profile["tool"]!={"name":"Kubernetes kubectl","version":"1.37.0","source":KUBECTL_SOURCE,"sha256":KUBECTL_SHA}: raise ValueError("exact official kubectl 1.37.0 identity is required")
 target=profile["target"]
 if not isinstance(target,dict) or set(target)!=TARGET_FIELDS: raise ValueError("target fields must be exact")
 if not re.fullmatch(r"sha256:[0-9a-f]{64}",str(target["cluster_uid_sha256"])): raise ValueError("cluster_uid_sha256 must be sha256")
 for key in ("context","namespace","workload_name","container_name"): nonempty(target[key],f"target.{key}")
 if target["workload_kind"]!="Deployment": raise ValueError("only Kubernetes Deployment rollout semantics are supported")
 https(target["service_url"],"target.service_url")
 policy_path=safe(root,profile["policy"],"policy"); policy=read_json(policy_path,"policy")
 pfields={"schema","project_id","environment","release_digest","approved_at","approvers","target","desired_image_digest","previous_image_digest","desired_revision","previous_revision","timeout_seconds","canary","migration"}
 if not isinstance(policy,dict) or set(policy)!=pfields or policy.get("schema")!="elite-deploy-rollback-policy/v1": raise ValueError("deploy policy fields/schema must be exact")
 bind(policy,profile,"policy"); utc(policy["approved_at"],"policy.approved_at")
 if not isinstance(policy["approvers"],list) or len(policy["approvers"])<2 or len(set(policy["approvers"]))!=len(policy["approvers"]) or any(not nonempty(x,"policy approver") for x in policy["approvers"]): raise ValueError("deploy policy requires at least two distinct approvers")
 if policy["target"]!=target: raise ValueError("deploy policy target mismatch")
 for key in ("desired_image_digest","previous_image_digest"):
  if not re.fullmatch(r"[^\s@]+@sha256:[0-9a-f]{64}",str(policy[key])): raise ValueError(f"{key} must be an immutable registry digest")
 if policy["desired_image_digest"]==policy["previous_image_digest"]: raise ValueError("desired and previous image digests must differ")
 if not all(isinstance(policy[x],int) and policy[x]>=1 for x in ("desired_revision","previous_revision","timeout_seconds")) or policy["desired_revision"]<=policy["previous_revision"] or policy["timeout_seconds"]>3600: raise ValueError("revision/timeout policy is invalid")
 canary=policy["canary"]
 if not isinstance(canary,dict) or set(canary)!={"minimum_samples","max_error_rate","max_p95_latency_ms","automatic_rollback"} or not isinstance(canary["minimum_samples"],int) or canary["minimum_samples"]<3 or not isinstance(canary["max_error_rate"],(int,float)) or not 0<=canary["max_error_rate"]<=0.05 or not isinstance(canary["max_p95_latency_ms"],int) or canary["max_p95_latency_ms"]<1 or canary["automatic_rollback"] is not True: raise ValueError("canary policy is invalid")
 migration=policy["migration"]
 if not isinstance(migration,dict) or set(migration)!={"migration_id","expand_contract_required","forward_backward_required","rollback_data_recovery_ref"} or not nonempty(migration["migration_id"],"migration_id") or migration["expand_contract_required"] is not True or migration["forward_backward_required"] is not True or not nonempty(migration["rollback_data_recovery_ref"],"rollback_data_recovery_ref"): raise ValueError("migration policy is incomplete")
 runs=profile["runs"]
 if not isinstance(runs,list) or len(runs)!=len(KINDS) or any(not isinstance(x,dict) or set(x)!=RUN_FIELDS for x in runs) or [x["kind"] for x in runs]!=list(KINDS): raise ValueError("seven ordered deploy/rollback runs are required")
 evidence=[proof(root,policy_path)]; tool_hash=None
 for run in runs:
  kind=run["kind"]; observation_path=safe(root,run["observation"],f"{kind}.observation"); observation=read_json(observation_path,f"{kind}.observation"); validate_observation(observation,profile,policy,kind); observation_proof=proof(root,observation_path)
  receipt_path=safe(root,run["execution_receipt"],f"{kind}.receipt"); receipt=read_json(receipt_path,f"{kind}.receipt")
  if receipt.get("schema")!="elite-official-tool-execution/v1" or receipt.get("project_id")!=profile["project_id"] or receipt.get("environment")!="production" or receipt.get("release_digest")!=profile["release_digest"] or receipt.get("control_id")!="DEPLOY_ROLLBACK" or receipt.get("exit_code")!=0: raise ValueError(f"{kind} execution identity/exit mismatch")
  tool=receipt.get("tool")
  if tool!={"name":"Kubernetes kubectl","version":"1.37.0","source":KUBECTL_SOURCE,"sha256":KUBECTL_SHA}: raise ValueError(f"{kind} exact kubectl identity mismatch")
  tool_hash=tool_hash or tool["sha256"]
  if tool["sha256"]!=tool_hash: raise ValueError("all deploy/rollback runs must use the same kubectl hash")
  args=run["expected_arguments"]
  if not isinstance(args,list) or not args or any(not isinstance(x,str) or not x for x in args) or receipt.get("arguments_sha256")!=canonical_sha(args): raise ValueError(f"{kind} arguments are not hash-bound")
  expected_target=dict(target); expected_target.update({"run_kind":kind,"observation_sha256":observation_proof["sha256"]})
  if receipt.get("target")!=expected_target: raise ValueError(f"{kind} target/observation binding mismatch")
  env=receipt.get("environment_variable_names")
  if not isinstance(env,list) or any(re.search(r"(?:secret|password|token|key|credential|kubeconfig)",str(x),re.I) and not str(x).endswith("_REF") for x in env): raise ValueError(f"{kind} environment must expose references, not secrets")
  for channel in ("stdout","stderr"): verify_proof(root,receipt.get(channel),f"{kind}.{channel}")
  evidence.extend([proof(root,receipt_path),observation_proof,receipt["stdout"],receipt["stderr"]])
 return {"id":"DEPLOY_ROLLBACK","project_id":profile["project_id"],"environment":"production","release_digest":profile["release_digest"],"result":"PASS","executed_at":profile["evaluated_at"],"expires_at":profile["expires_at"],"executor_ref":profile["executor_ref"],"tool":{"name":"Kubernetes kubectl","version":"1.37.0","source":KUBECTL_SOURCE,"digest":"sha256:"+KUBECTL_SHA},"target":target,"assertions":{"immutable_digest_deployed":True,"canary_health_pass":True,"rollout_status_pass":True,"rollback_executed":True,"previous_digest_restored":True,"migration_compatibility_pass":True},"evidence":evidence}

def main()->int:
 p=argparse.ArgumentParser(); p.add_argument("--project-root",required=True,type=Path); p.add_argument("--profile",required=True); a=p.parse_args()
 try:
  root=a.project_root.resolve(strict=True); profile=read_json(safe(root,a.profile,"profile"),"profile"); output=safe(root,profile.get("output"),"output",exists=False); output.parent.resolve(strict=True).relative_to(root); receipt=validate(profile,root); atomic(output,(json.dumps(receipt,indent=2,sort_keys=True)+"\n").encode()); print(f"DEPLOY_ROLLBACK_ADMISSION_PASS project={receipt['project_id']} release={receipt['release_digest']}"); return 0
 except (ValueError,OSError) as error: print(f"DEPLOY_ROLLBACK_ADMISSION_FAILED: {error}",file=sys.stderr); return 1
if __name__=="__main__": sys.exit(main())
````

### FILE: `production_admission_gate/test_validate_deploy_rollback_admission.py`

```yaml
block_id: "SECURE-OPS:deploy-rollback-admission-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://production-admission/kubernetes-deploy-rollback-tests"
license: "LicenseRef-Workspace-Owner"
sha256: "405ed7a9a35647e890c804eda34904d0e9fd23a2d8b6c1ca8123dea3b23dd99d"
variables: []
secrets_allowed: false
```

````python
from __future__ import annotations
import copy,hashlib,json,tempfile,unittest
from datetime import datetime,timedelta,timezone
from pathlib import Path
from validate_deploy_rollback_admission import KINDS,KUBECTL_SHA,KUBECTL_SOURCE,validate
from validate_k6_load_admission import canonical_sha

class DeployRollbackAdmissionTests(unittest.TestCase):
 def setUp(self):
  self.temp=tempfile.TemporaryDirectory(); self.root=Path(self.temp.name); (self.root/"evidence").mkdir(); now=datetime.now(timezone.utc); z=lambda d:d.isoformat().replace("+00:00","Z"); digest="sha256:"+"e"*64; desired="registry.example.test/revestex@sha256:"+"a"*64; previous="registry.example.test/revestex@sha256:"+"b"*64
  self.target={"cluster_uid_sha256":"sha256:"+"c"*64,"context":"prod-ar","namespace":"revestex","workload_kind":"Deployment","workload_name":"api","container_name":"api","service_url":"https://api.example.test/health"}
  self.profile={"schema":"elite-deploy-rollback-admission/v1","project_id":"revestex","environment":"production","release_digest":digest,"evaluated_at":z(now),"expires_at":z(now+timedelta(days=7)),"executor_ref":"deploy:change/112","tool":{"name":"Kubernetes kubectl","version":"1.37.0","source":KUBECTL_SOURCE,"sha256":KUBECTL_SHA},"target":self.target,"policy":"evidence/policy.json","runs":[],"output":"evidence/deploy-control.json"}
  self.policy={"schema":"elite-deploy-rollback-policy/v1","project_id":"revestex","environment":"production","release_digest":digest,"approved_at":z(now),"approvers":["platform-owner","business-owner"],"target":self.target,"desired_image_digest":desired,"previous_image_digest":previous,"desired_revision":42,"previous_revision":41,"timeout_seconds":600,"canary":{"minimum_samples":3,"max_error_rate":0.01,"max_p95_latency_ms":500,"automatic_rollback":True},"migration":{"migration_id":"expand-112","expand_contract_required":True,"forward_backward_required":True,"rollback_data_recovery_ref":"recovery:112"}}
  self.write("evidence/policy.json",self.policy)
  common={"schema":"elite-deploy-rollback-observation/v1","project_id":"revestex","environment":"production","release_digest":digest,"captured_at":z(now)}
  observations={
   "PREFLIGHT":dict(common,kind="PREFLIGHT",client_version="v1.37.0",server_version="v1.37.0",cluster_uid_sha256=self.target["cluster_uid_sha256"],context="prod-ar",namespace="revestex",authorized_verbs=["get","list","patch","watch"],workload_found=True),
   "APPLY_IMMUTABLE_DIGEST":dict(common,kind="APPLY_IMMUTABLE_DIGEST",workload_kind="Deployment",workload_name="api",container_name="api",requested_image_digest=desired,observed_image_digest=desired,revision_before=41,revision_after=42,generation=12,observed_generation=12),
   "CANARY_HEALTH":dict(common,kind="CANARY_HEALTH",service_url=self.target["service_url"],samples=5,healthy_samples=5,error_rate=0.0,p95_latency_ms=120,probes_pass=True,alerts_firing=[],error_budget_breached=False),
   "ROLLOUT_STATUS":dict(common,kind="ROLLOUT_STATUS",workload_kind="Deployment",workload_name="api",progressing=True,progressing_reason="NewReplicaSetAvailable",available=True,desired_replicas=3,updated_replicas=3,available_replicas=3,old_replicas=0,application_probes_pass=True),
   "ROLLBACK":dict(common,kind="ROLLBACK",workload_kind="Deployment",workload_name="api",from_revision=42,to_revision=41,rollback_command_issued=True,rollout_complete=True),
   "POST_ROLLBACK_VERIFY":dict(common,kind="POST_ROLLBACK_VERIFY",observed_image_digest=previous,desired_replicas=3,available_replicas=3,service_health_pass=True,data_invariants_pass=True,alerts_firing=[]),
   "MIGRATION_COMPATIBILITY":dict(common,kind="MIGRATION_COMPATIBILITY",migration_id="expand-112",expand_contract_used=True,forward_compatible=True,backward_compatible=True,old_new_versions_coexisted=True,rollback_data_safe=True,integrity_checks_pass=True,destructive_change_before_convergence=False)}
  for kind in KINDS:
   stem=kind.lower(); obs=self.write(f"evidence/{stem}.json",observations[kind]); args=["--context","prod-ar","--namespace","revestex",kind.lower()]; out=self.file(f"evidence/{stem}.stdout",b"pass\n"); err=self.file(f"evidence/{stem}.stderr",b""); target=dict(self.target); target.update({"run_kind":kind,"observation_sha256":self.pf(obs)["sha256"]}); receipt={"schema":"elite-official-tool-execution/v1","project_id":"revestex","environment":"production","release_digest":digest,"control_id":"DEPLOY_ROLLBACK","executed_at":z(now),"exit_code":0,"tool":{"name":"Kubernetes kubectl","version":"1.37.0","source":KUBECTL_SOURCE,"sha256":KUBECTL_SHA},"arguments_sha256":canonical_sha(args),"working_directory":".","target":target,"environment_variable_names":["KUBECONFIG_REF"],"stdout":self.pf(out),"stderr":self.pf(err)}; self.write(f"evidence/{stem}-receipt.json",receipt); self.profile["runs"].append({"kind":kind,"execution_receipt":f"evidence/{stem}-receipt.json","observation":f"evidence/{stem}.json","expected_arguments":args})
 def tearDown(self): self.temp.cleanup()
 def file(self,p,data): q=self.root/p; q.write_bytes(data); return q
 def write(self,p,v): q=self.root/p; q.write_text(json.dumps(v),encoding="utf-8"); return q
 def pf(self,p): data=p.read_bytes(); return {"path":p.relative_to(self.root).as_posix(),"bytes":len(data),"sha256":hashlib.sha256(data).hexdigest()}
 def mutate(self,index,change): p=self.root/self.profile["runs"][index]["observation"]; v=json.loads(p.read_text()); change(v); self.write(self.profile["runs"][index]["observation"],v)
 def test_complete_realistic_evidence_emits_semantic_receipt(self): self.assertTrue(all(validate(self.profile,self.root)["assertions"].values()))
 def test_mutable_image_is_rejected(self):
  p=self.root/self.profile["policy"]; v=json.loads(p.read_text()); v["desired_image_digest"]="registry.example.test/revestex:latest"; self.write(self.profile["policy"],v)
  with self.assertRaisesRegex(ValueError,"immutable registry digest"): validate(self.profile,self.root)
 def test_apply_digest_drift_is_rejected(self):
  self.mutate(1,lambda v:v.update(observed_image_digest="registry.example.test/revestex@sha256:"+"d"*64))
  with self.assertRaisesRegex(ValueError,"desired image digest"): validate(self.profile,self.root)
 def test_canary_budget_breach_is_rejected(self):
  self.mutate(2,lambda v:v.update(error_budget_breached=True))
  with self.assertRaisesRegex(ValueError,"canary health"): validate(self.profile,self.root)
 def test_rollout_old_replicas_are_rejected(self):
  self.mutate(3,lambda v:v.update(old_replicas=1))
  with self.assertRaisesRegex(ValueError,"rollout status"): validate(self.profile,self.root)
 def test_rollback_not_executed_is_rejected(self):
  self.mutate(4,lambda v:v.update(rollback_command_issued=False))
  with self.assertRaisesRegex(ValueError,"actual rollback"): validate(self.profile,self.root)
 def test_previous_digest_not_restored_is_rejected(self):
  self.mutate(5,lambda v:v.update(observed_image_digest=self.policy["desired_image_digest"]))
  with self.assertRaisesRegex(ValueError,"previous digest"): validate(self.profile,self.root)
 def test_unsafe_migration_is_rejected(self):
  self.mutate(6,lambda v:v.update(backward_compatible=False))
  with self.assertRaisesRegex(ValueError,"migration compatibility"): validate(self.profile,self.root)
 def test_wrong_cluster_is_rejected(self):
  self.mutate(0,lambda v:v.update(cluster_uid_sha256="sha256:"+"f"*64))
  with self.assertRaisesRegex(ValueError,"preflight target"): validate(self.profile,self.root)
 def test_nonzero_kubectl_exit_is_rejected(self):
  p=self.root/self.profile["runs"][3]["execution_receipt"]; v=json.loads(p.read_text()); v["exit_code"]=1; self.write(self.profile["runs"][3]["execution_receipt"],v)
  with self.assertRaisesRegex(ValueError,"identity/exit"): validate(self.profile,self.root)
 def test_observation_tamper_breaks_receipt_binding(self):
  self.mutate(3,lambda v:v.update(desired_replicas=4,updated_replicas=4,available_replicas=4))
  with self.assertRaisesRegex(ValueError,"target/observation binding"): validate(self.profile,self.root)
 def test_wrong_kubectl_hash_is_rejected(self):
  p=self.root/self.profile["runs"][0]["execution_receipt"]; v=json.loads(p.read_text()); v["tool"]["sha256"]="0"*64; self.write(self.profile["runs"][0]["execution_receipt"],v)
  with self.assertRaisesRegex(ValueError,"exact kubectl identity"): validate(self.profile,self.root)
 def test_secret_kubeconfig_name_is_rejected(self):
  p=self.root/self.profile["runs"][0]["execution_receipt"]; v=json.loads(p.read_text()); v["environment_variable_names"]=["KUBECONFIG"]; self.write(self.profile["runs"][0]["execution_receipt"],v)
  with self.assertRaisesRegex(ValueError,"references, not secrets"): validate(self.profile,self.root)
 def test_missing_ordered_run_is_rejected(self):
  self.profile["runs"]=self.profile["runs"][:-1]
  with self.assertRaisesRegex(ValueError,"seven ordered"): validate(self.profile,self.root)

if __name__=="__main__": unittest.main()
````

### FILE: `production_admission_gate/business-acceptance-admission.template.json`

```yaml
block_id: "SECURE-OPS:business-acceptance-admission-template:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://production-admission/business-acceptance-template"
license: "LicenseRef-Workspace-Owner"
sha256: "6ac4cbfa13035f11454506e26f5cf6c0ac19ab51ba798c960a6fcc8ec70141b4"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite-business-acceptance-admission/v1",
  "project_id": "",
  "environment": "production",
  "release_digest": "",
  "evaluated_at": "",
  "expires_at": "",
  "executor_ref": "",
  "tool": {
    "name": "Microsoft Playwright",
    "version": "1.62.1",
    "source": "https://github.com/microsoft/playwright/releases/tag/v1.62.1"
  },
  "target": {},
  "specification": "",
  "playwright_run": {},
  "approvals": "",
  "defects": "",
  "owner_signoff": "",
  "output": ""
}
````

### FILE: `production_admission_gate/validate_business_acceptance_admission.py`

```yaml
block_id: "SECURE-OPS:business-acceptance-admission-validator:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://production-admission/business-acceptance-validator"
license: "LicenseRef-Workspace-Owner"
sha256: "530f6927bcf0916d527087670266d9227e0af6558f89efbd9743b345091f43ea"
variables: []
secrets_allowed: false
```

````python
from __future__ import annotations
import argparse,json,re,sys
from datetime import datetime
from pathlib import Path
from urllib.parse import urlparse
from validate_k6_load_admission import atomic,canonical_sha,proof,read_json,safe,utc,verify_proof

FIELDS={"schema","project_id","environment","release_digest","evaluated_at","expires_at","executor_ref","tool","target","specification","playwright_run","approvals","defects","owner_signoff","output"}
TARGET_FIELDS={"application_url","business_scope_id","tenant_ids","role_ids"}
RUN_FIELDS={"execution_receipt","observation","expected_arguments"}
PLAYWRIGHT_SOURCE="https://github.com/microsoft/playwright/releases/tag/v1.62.1"
SPEC_KIT_SOURCE="https://github.com/github/spec-kit/releases/tag/v1.0.1"
SPEC_KIT_COMMIT="9118ed15a0ba65053469a94c560ea5d233f75884"

def nonempty(value:object,label:str)->str:
 if not isinstance(value,str) or not value.strip(): raise ValueError(f"{label} is required")
 return value

def unique(value:object,label:str,minimum:int=1)->list[str]:
 if not isinstance(value,list) or len(value)<minimum or len(set(value))!=len(value) or any(not isinstance(x,str) or not x.strip() for x in value): raise ValueError(f"{label} must contain at least {minimum} unique non-empty values")
 return value

def https(value:object,label:str)->None:
 value=nonempty(value,label); p=urlparse(value)
 if p.scheme!="https" or not p.hostname or p.username or p.password or p.fragment: raise ValueError(f"{label} must be credential-free HTTPS")

def bind(value:dict[str,object],profile:dict[str,object],label:str)->None:
 for key in ("project_id","environment","release_digest"):
  if value.get(key)!=profile[key]: raise ValueError(f"{label} {key} mismatch")

def validate(profile:object,root:Path)->dict[str,object]:
 if not isinstance(profile,dict) or set(profile)!=FIELDS: raise ValueError("profile fields must be exact")
 if profile["schema"]!="elite-business-acceptance-admission/v1" or profile["environment"]!="production" or not nonempty(profile["project_id"],"project_id") or not re.fullmatch(r"sha256:[0-9a-f]{64}",str(profile["release_digest"])) or set(str(profile["release_digest"]).split(":")[-1])=={"0"}: raise ValueError("profile identity is invalid")
 start=utc(profile["evaluated_at"],"evaluated_at"); end=utc(profile["expires_at"],"expires_at")
 if not start<end or (end-start).total_seconds()>30*86400: raise ValueError("admission validity must be positive and at most 30 days")
 nonempty(profile["executor_ref"],"executor_ref")
 if profile["tool"]!={"name":"Microsoft Playwright","version":"1.62.1","source":PLAYWRIGHT_SOURCE}: raise ValueError("exact Microsoft Playwright 1.62.1 identity is required")
 target=profile["target"]
 if not isinstance(target,dict) or set(target)!=TARGET_FIELDS: raise ValueError("target fields must be exact")
 https(target["application_url"],"target.application_url"); nonempty(target["business_scope_id"],"target.business_scope_id"); unique(target["tenant_ids"],"target.tenant_ids",1); unique(target["role_ids"],"target.role_ids",1)

 spec_path=safe(root,profile["specification"],"specification"); spec=read_json(spec_path,"specification"); spec_proof=proof(root,spec_path)
 sfields={"schema","project_id","environment","release_digest","approved_at","spec_kit","artifacts","capabilities","scenarios"}
 if not isinstance(spec,dict) or set(spec)!=sfields or spec.get("schema")!="elite-business-acceptance-specification/v1": raise ValueError("specification fields/schema must be exact")
 bind(spec,profile,"specification"); utc(spec["approved_at"],"specification.approved_at")
 if spec["spec_kit"]!={"name":"GitHub Spec Kit","version":"1.0.1","source":SPEC_KIT_SOURCE,"commit":SPEC_KIT_COMMIT}: raise ValueError("exact GitHub Spec Kit 1.0.1 identity is required")
 artifacts=spec["artifacts"]
 required_artifacts={"constitution","spec","plan","tasks","convergence"}
 if not isinstance(artifacts,list) or {x.get("kind") for x in artifacts if isinstance(x,dict)}!=required_artifacts or any(not isinstance(x,dict) or set(x)!={"kind","evidence"} for x in artifacts): raise ValueError("five exact Spec Kit artifacts are required")
 evidence=[spec_proof]
 for item in artifacts:
  verify_proof(root,item["evidence"],f"specification artifact {item['kind']}"); evidence.append(item["evidence"])
 capabilities=unique(spec["capabilities"],"specification.capabilities",1); scenarios=spec["scenarios"]
 sf={"scenario_id","capability_id","actor_role","tenant_id","preconditions","expected_outcomes"}
 if not isinstance(scenarios,list) or not scenarios or any(not isinstance(x,dict) or set(x)!=sf for x in scenarios): raise ValueError("acceptance scenarios must use the exact schema")
 scenario_ids=[]
 for item in scenarios:
  scenario_ids.append(nonempty(item["scenario_id"],"scenario_id"))
  if item["capability_id"] not in capabilities or item["actor_role"] not in target["role_ids"] or item["tenant_id"] not in target["tenant_ids"] or not unique(item["preconditions"],"scenario.preconditions") or not unique(item["expected_outcomes"],"scenario.expected_outcomes"): raise ValueError("scenario is not bound to declared capability/role/tenant/outcomes")
 if len(set(scenario_ids))!=len(scenario_ids) or {x["capability_id"] for x in scenarios}!=set(capabilities): raise ValueError("every capability requires scenarios with unique IDs")

 run=profile["playwright_run"]
 if not isinstance(run,dict) or set(run)!=RUN_FIELDS: raise ValueError("playwright_run fields must be exact")
 obs_path=safe(root,run["observation"],"playwright observation"); obs=read_json(obs_path,"playwright observation"); obs_proof=proof(root,obs_path)
 ofields={"schema","project_id","environment","release_digest","captured_at","application_url","business_scope_id","specification_sha256","results"}
 if not isinstance(obs,dict) or set(obs)!=ofields or obs.get("schema")!="elite-business-acceptance-observation/v1": raise ValueError("business observation fields/schema mismatch")
 bind(obs,profile,"observation"); utc(obs["captured_at"],"observation.captured_at")
 if obs["application_url"]!=target["application_url"] or obs["business_scope_id"]!=target["business_scope_id"] or obs["specification_sha256"]!=spec_proof["sha256"]: raise ValueError("business observation target/specification mismatch")
 results=obs["results"]; rf={"scenario_id","status","expected_outcomes_observed","backend_effect_verified","financial_effect","evidence"}
 if not isinstance(results,list) or {x.get("scenario_id") for x in results if isinstance(x,dict)}!=set(scenario_ids) or len(results)!=len(scenario_ids): raise ValueError("every approved scenario requires exactly one result")
 for item in results:
  if not isinstance(item,dict) or set(item)!=rf or item["status"]!="PASS" or item["expected_outcomes_observed"] is not True or item["backend_effect_verified"] is not True or item["financial_effect"] not in ("VERIFIED","NOT_APPLICABLE"): raise ValueError("acceptance scenario result did not prove its outcomes/effects")
  verify_proof(root,item["evidence"],f"scenario {item.get('scenario_id')} evidence"); evidence.append(item["evidence"])
 receipt_path=safe(root,run["execution_receipt"],"playwright receipt"); receipt=read_json(receipt_path,"playwright receipt")
 if receipt.get("schema")!="elite-official-tool-execution/v1" or receipt.get("project_id")!=profile["project_id"] or receipt.get("environment")!="production" or receipt.get("release_digest")!=profile["release_digest"] or receipt.get("control_id")!="BUSINESS_ACCEPTANCE" or receipt.get("exit_code")!=0: raise ValueError("Playwright execution identity/exit mismatch")
 tool=receipt.get("tool")
 if not isinstance(tool,dict) or set(tool)!={"name","version","source","sha256"} or (tool.get("name"),tool.get("version"),tool.get("source"))!=("Microsoft Playwright","1.62.1",PLAYWRIGHT_SOURCE) or not re.fullmatch(r"[0-9a-f]{64}",str(tool.get("sha256"))) or set(str(tool.get("sha256")))=={"0"}: raise ValueError("exact Microsoft Playwright execution identity is required")
 args=run["expected_arguments"]
 if not isinstance(args,list) or not args or any(not isinstance(x,str) or not x for x in args) or receipt.get("arguments_sha256")!=canonical_sha(args): raise ValueError("Playwright arguments are not hash-bound")
 expected_target=dict(target); expected_target.update({"specification_sha256":spec_proof["sha256"],"observation_sha256":obs_proof["sha256"]})
 if receipt.get("target")!=expected_target: raise ValueError("Playwright target/observation binding mismatch")
 env=receipt.get("environment_variable_names")
 if not isinstance(env,list) or any(re.search(r"(?:secret|password|token|key|credential)",str(x),re.I) and not str(x).endswith("_REF") for x in env): raise ValueError("Playwright environment must expose references, not secrets")
 for channel in ("stdout","stderr"): verify_proof(root,receipt.get(channel),f"playwright.{channel}")
 evidence.extend([obs_proof,proof(root,receipt_path),receipt["stdout"],receipt["stderr"]])

 approvals_path=safe(root,profile["approvals"],"approvals"); approvals=read_json(approvals_path,"approvals"); approvals_proof=proof(root,approvals_path)
 afields={"schema","project_id","environment","release_digest","captured_at","specification_sha256","observation_sha256","approvals"}
 if not isinstance(approvals,dict) or set(approvals)!=afields or approvals.get("schema")!="elite-business-functional-approvals/v1": raise ValueError("functional approvals fields/schema mismatch")
 bind(approvals,profile,"approvals"); utc(approvals["captured_at"],"approvals.captured_at")
 if approvals["specification_sha256"]!=spec_proof["sha256"] or approvals["observation_sha256"]!=obs_proof["sha256"]: raise ValueError("functional approvals evidence binding mismatch")
 required_roles={"finance","operations","security"}; decisions=approvals["approvals"]; df={"role","subject_sha256","decision","approved_at","evidence"}
 if not isinstance(decisions,list) or {x.get("role") for x in decisions if isinstance(x,dict)}!=required_roles or len(decisions)!=3: raise ValueError("finance, operations and security approvals are required")
 subjects=[]
 for item in decisions:
  if not isinstance(item,dict) or set(item)!=df or not re.fullmatch(r"sha256:[0-9a-f]{64}",str(item["subject_sha256"])) or item["decision"]!="APPROVED": raise ValueError("functional approval is invalid")
  utc(item["approved_at"],f"{item.get('role')} approved_at"); verify_proof(root,item["evidence"],f"{item.get('role')} approval evidence"); subjects.append(item["subject_sha256"]); evidence.append(item["evidence"])
 if len(set(subjects))!=3: raise ValueError("finance, operations and security approvers must be distinct")
 evidence.append(approvals_proof)

 defects_path=safe(root,profile["defects"],"defects"); defects=read_json(defects_path,"defects"); defects_proof=proof(root,defects_path)
 dfields={"schema","project_id","environment","release_digest","captured_at","specification_sha256","defects"}
 if not isinstance(defects,dict) or set(defects)!=dfields or defects.get("schema")!="elite-business-defect-register/v1": raise ValueError("defect register fields/schema mismatch")
 bind(defects,profile,"defects"); utc(defects["captured_at"],"defects.captured_at")
 if defects["specification_sha256"]!=spec_proof["sha256"] or not isinstance(defects["defects"],list): raise ValueError("defect register binding/list mismatch")
 seen=set()
 for item in defects["defects"]:
  fields={"defect_id","severity","status","reason","reviewers","expires_at"}
  if not isinstance(item,dict) or set(item)!=fields or not nonempty(item["defect_id"],"defect_id") or item["defect_id"] in seen or item["severity"] not in (1,2,3,4) or item["status"] not in ("CLOSED","ACCEPTED"): raise ValueError("open, duplicate or malformed defect is forbidden")
  seen.add(item["defect_id"])
  if item["severity"] in (1,2) and item["status"]!="CLOSED": raise ValueError("Sev1/Sev2 defects cannot be accepted open")
  if item["status"]=="ACCEPTED":
   nonempty(item["reason"],"accepted defect reason"); unique(item["reviewers"],"accepted defect reviewers",2); expiry=utc(item["expires_at"],"accepted defect expires_at")
   if expiry<=start: raise ValueError("accepted defect expiry must be after evaluation")
 evidence.append(defects_proof)

 signoff_path=safe(root,profile["owner_signoff"],"owner_signoff"); signoff=read_json(signoff_path,"owner_signoff"); signoff_proof=proof(root,signoff_path)
 xfields={"schema","project_id","environment","release_digest","approved_at","role","subject_sha256","decision","specification_sha256","observation_sha256","approvals_sha256","defects_sha256","evidence"}
 if not isinstance(signoff,dict) or set(signoff)!=xfields or signoff.get("schema")!="elite-business-owner-signoff/v1": raise ValueError("owner signoff fields/schema mismatch")
 bind(signoff,profile,"owner signoff"); utc(signoff["approved_at"],"owner signoff approved_at")
 if signoff["role"]!="business_owner" or signoff["decision"]!="APPROVED" or not re.fullmatch(r"sha256:[0-9a-f]{64}",str(signoff["subject_sha256"])) or signoff["subject_sha256"] in subjects: raise ValueError("business owner signoff must be approved and independently owned")
 if (signoff["specification_sha256"],signoff["observation_sha256"],signoff["approvals_sha256"],signoff["defects_sha256"])!=(spec_proof["sha256"],obs_proof["sha256"],approvals_proof["sha256"],defects_proof["sha256"]): raise ValueError("business owner signoff evidence binding mismatch")
 verify_proof(root,signoff["evidence"],"business owner signoff evidence"); evidence.extend([signoff_proof,signoff["evidence"]])
 return {"id":"BUSINESS_ACCEPTANCE","project_id":profile["project_id"],"environment":"production","release_digest":profile["release_digest"],"result":"PASS","executed_at":profile["evaluated_at"],"expires_at":profile["expires_at"],"executor_ref":profile["executor_ref"],"tool":{"name":"GitHub Spec Kit + Microsoft Playwright","version":"1.0.1 + 1.62.1","source":SPEC_KIT_SOURCE,"digest":"sha256:"+canonical_sha([SPEC_KIT_COMMIT,tool["sha256"]])},"target":target,"assertions":{"acceptance_scenarios_pass":True,"finance_operations_security_approved":True,"no_open_sev1_sev2":True,"approvers_distinct":True,"owner_signoff":True},"evidence":evidence}

def main()->int:
 p=argparse.ArgumentParser(); p.add_argument("--project-root",required=True,type=Path); p.add_argument("--profile",required=True); a=p.parse_args()
 try:
  root=a.project_root.resolve(strict=True); profile=read_json(safe(root,a.profile,"profile"),"profile"); output=safe(root,profile.get("output"),"output",exists=False); output.parent.resolve(strict=True).relative_to(root); receipt=validate(profile,root); atomic(output,(json.dumps(receipt,indent=2,sort_keys=True)+"\n").encode()); print(f"BUSINESS_ACCEPTANCE_ADMISSION_PASS project={receipt['project_id']} release={receipt['release_digest']}"); return 0
 except (ValueError,OSError) as error: print(f"BUSINESS_ACCEPTANCE_ADMISSION_FAILED: {error}",file=sys.stderr); return 1
if __name__=="__main__": sys.exit(main())
````

### FILE: `production_admission_gate/test_validate_business_acceptance_admission.py`

```yaml
block_id: "SECURE-OPS:business-acceptance-admission-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://production-admission/business-acceptance-tests"
license: "LicenseRef-Workspace-Owner"
sha256: "405c28e3958471684cb3963635fd9ad6bf4b2f5d06915f2fde010561a076ef2b"
variables: []
secrets_allowed: false
```

````python
from __future__ import annotations
import hashlib,json,tempfile,unittest
from datetime import datetime,timedelta,timezone
from pathlib import Path
from validate_business_acceptance_admission import PLAYWRIGHT_SOURCE,SPEC_KIT_COMMIT,SPEC_KIT_SOURCE,validate
from validate_k6_load_admission import canonical_sha

class BusinessAcceptanceAdmissionTests(unittest.TestCase):
 def setUp(self):
  self.temp=tempfile.TemporaryDirectory(); self.root=Path(self.temp.name); (self.root/"evidence").mkdir(); now=datetime.now(timezone.utc); self.now=now; z=lambda d:d.isoformat().replace("+00:00","Z"); digest="sha256:"+"e"*64; self.target={"application_url":"https://app.example.test","business_scope_id":"launch-112","tenant_ids":["north","south"],"role_ids":["admin","seller"]}
  self.profile={"schema":"elite-business-acceptance-admission/v1","project_id":"revestex","environment":"production","release_digest":digest,"evaluated_at":z(now),"expires_at":z(now+timedelta(days=7)),"executor_ref":"acceptance:change/113","tool":{"name":"Microsoft Playwright","version":"1.62.1","source":PLAYWRIGHT_SOURCE},"target":self.target,"specification":"evidence/specification.json","playwright_run":{},"approvals":"evidence/approvals.json","defects":"evidence/defects.json","owner_signoff":"evidence/owner.json","output":"evidence/business-control.json"}
  artifacts=[]
  for kind in ("constitution","spec","plan","tasks","convergence"):
   p=self.file(f"evidence/{kind}.md",f"# {kind}\n".encode()); artifacts.append({"kind":kind,"evidence":self.pf(p)})
  scenarios=[{"scenario_id":"sale-order","capability_id":"commerce","actor_role":"seller","tenant_id":"north","preconditions":["catalog active"],"expected_outcomes":["order committed","receipt visible"]},{"scenario_id":"admin-audit","capability_id":"administration","actor_role":"admin","tenant_id":"south","preconditions":["admin authenticated"],"expected_outcomes":["audit visible"]}]
  spec={"schema":"elite-business-acceptance-specification/v1","project_id":"revestex","environment":"production","release_digest":digest,"approved_at":z(now),"spec_kit":{"name":"GitHub Spec Kit","version":"1.0.1","source":SPEC_KIT_SOURCE,"commit":SPEC_KIT_COMMIT},"artifacts":artifacts,"capabilities":["commerce","administration"],"scenarios":scenarios}; sp=self.write(self.profile["specification"],spec); spec_hash=self.pf(sp)["sha256"]
  results=[]
  for item in scenarios:
   ep=self.file(f"evidence/{item['scenario_id']}.json",b'{"verified":true}\n'); results.append({"scenario_id":item["scenario_id"],"status":"PASS","expected_outcomes_observed":True,"backend_effect_verified":True,"financial_effect":"VERIFIED" if item["capability_id"]=="commerce" else "NOT_APPLICABLE","evidence":self.pf(ep)})
  obs={"schema":"elite-business-acceptance-observation/v1","project_id":"revestex","environment":"production","release_digest":digest,"captured_at":z(now),"application_url":self.target["application_url"],"business_scope_id":"launch-112","specification_sha256":spec_hash,"results":results}; op=self.write("evidence/observation.json",obs); obs_hash=self.pf(op)["sha256"]
  args=["test","acceptance","--project=chromium-desktop"]; out=self.file("evidence/playwright.stdout",b"2 passed\n"); err=self.file("evidence/playwright.stderr",b""); receipt={"schema":"elite-official-tool-execution/v1","project_id":"revestex","environment":"production","release_digest":digest,"control_id":"BUSINESS_ACCEPTANCE","executed_at":z(now),"exit_code":0,"tool":{"name":"Microsoft Playwright","version":"1.62.1","source":PLAYWRIGHT_SOURCE,"sha256":"a"*64},"arguments_sha256":canonical_sha(args),"working_directory":".","target":dict(self.target,specification_sha256=spec_hash,observation_sha256=obs_hash),"environment_variable_names":["TEST_USER_TOKEN_REF"],"stdout":self.pf(out),"stderr":self.pf(err)}; self.write("evidence/playwright-receipt.json",receipt); self.profile["playwright_run"]={"execution_receipt":"evidence/playwright-receipt.json","observation":"evidence/observation.json","expected_arguments":args}
  approvals=[]
  for idx,role in enumerate(("finance","operations","security"),1):
   ep=self.file(f"evidence/{role}-approval.txt",f"approved {role}\n".encode()); approvals.append({"role":role,"subject_sha256":"sha256:"+str(idx)*64,"decision":"APPROVED","approved_at":z(now),"evidence":self.pf(ep)})
  ap=self.write(self.profile["approvals"],{"schema":"elite-business-functional-approvals/v1","project_id":"revestex","environment":"production","release_digest":digest,"captured_at":z(now),"specification_sha256":spec_hash,"observation_sha256":obs_hash,"approvals":approvals})
  dp=self.write(self.profile["defects"],{"schema":"elite-business-defect-register/v1","project_id":"revestex","environment":"production","release_digest":digest,"captured_at":z(now),"specification_sha256":spec_hash,"defects":[{"defect_id":"UI-LOW-1","severity":4,"status":"ACCEPTED","reason":"cosmetic only","reviewers":["product-owner","ux-owner"],"expires_at":z(now+timedelta(days=5))}]})
  ep=self.file("evidence/owner-approval.txt",b"approved owner\n"); self.write(self.profile["owner_signoff"],{"schema":"elite-business-owner-signoff/v1","project_id":"revestex","environment":"production","release_digest":digest,"approved_at":z(now),"role":"business_owner","subject_sha256":"sha256:"+"4"*64,"decision":"APPROVED","specification_sha256":spec_hash,"observation_sha256":obs_hash,"approvals_sha256":self.pf(ap)["sha256"],"defects_sha256":self.pf(dp)["sha256"],"evidence":self.pf(ep)})
 def tearDown(self): self.temp.cleanup()
 def file(self,p,data): q=self.root/p; q.write_bytes(data); return q
 def write(self,p,v): q=self.root/p; q.write_text(json.dumps(v),encoding="utf-8"); return q
 def pf(self,p): data=p.read_bytes(); return {"path":p.relative_to(self.root).as_posix(),"bytes":len(data),"sha256":hashlib.sha256(data).hexdigest()}
 def mutate(self,path,change): p=self.root/path; v=json.loads(p.read_text()); change(v); self.write(path,v)
 def test_complete_business_evidence_emits_semantic_receipt(self): self.assertTrue(all(validate(self.profile,self.root)["assertions"].values()))
 def test_missing_capability_scenario_is_rejected(self):
  self.mutate(self.profile["specification"],lambda v:v.update(scenarios=v["scenarios"][:1]))
  with self.assertRaisesRegex(ValueError,"every capability"): validate(self.profile,self.root)
 def test_failed_scenario_is_rejected(self):
  self.mutate("evidence/observation.json",lambda v:v["results"][0].update(status="FAIL"))
  with self.assertRaisesRegex(ValueError,"did not prove"): validate(self.profile,self.root)
 def test_backend_effect_is_required(self):
  self.mutate("evidence/observation.json",lambda v:v["results"][0].update(backend_effect_verified=False))
  with self.assertRaisesRegex(ValueError,"did not prove"): validate(self.profile,self.root)
 def test_observation_tamper_breaks_receipt_binding(self):
  self.mutate("evidence/observation.json",lambda v:v["results"][0].update(financial_effect="NOT_APPLICABLE"))
  with self.assertRaisesRegex(ValueError,"target/observation binding"): validate(self.profile,self.root)
 def test_nonzero_playwright_exit_is_rejected(self):
  self.mutate("evidence/playwright-receipt.json",lambda v:v.update(exit_code=1))
  with self.assertRaisesRegex(ValueError,"identity/exit"): validate(self.profile,self.root)
 def test_wrong_playwright_version_is_rejected(self):
  self.mutate("evidence/playwright-receipt.json",lambda v:v["tool"].update(version="latest"))
  with self.assertRaisesRegex(ValueError,"exact Microsoft Playwright execution"): validate(self.profile,self.root)
 def test_secret_environment_name_is_rejected(self):
  self.mutate("evidence/playwright-receipt.json",lambda v:v.update(environment_variable_names=["TEST_USER_TOKEN"]))
  with self.assertRaisesRegex(ValueError,"references, not secrets"): validate(self.profile,self.root)
 def test_missing_security_approval_is_rejected(self):
  self.mutate(self.profile["approvals"],lambda v:v.update(approvals=v["approvals"][:2]))
  with self.assertRaisesRegex(ValueError,"finance, operations and security"): validate(self.profile,self.root)
 def test_same_functional_approver_is_rejected(self):
  self.mutate(self.profile["approvals"],lambda v:v["approvals"][1].update(subject_sha256=v["approvals"][0]["subject_sha256"]))
  with self.assertRaisesRegex(ValueError,"must be distinct"): validate(self.profile,self.root)
 def test_open_sev1_is_rejected(self):
  self.mutate(self.profile["defects"],lambda v:v["defects"].append({"defect_id":"BLOCKER","severity":1,"status":"ACCEPTED","reason":"later","reviewers":["a","b"],"expires_at":(self.now+timedelta(days=1)).isoformat().replace("+00:00","Z")}))
  with self.assertRaisesRegex(ValueError,"Sev1/Sev2"): validate(self.profile,self.root)
 def test_unreviewed_accepted_defect_is_rejected(self):
  self.mutate(self.profile["defects"],lambda v:v["defects"][0].update(reviewers=["one"]))
  with self.assertRaisesRegex(ValueError,"at least 2"): validate(self.profile,self.root)
 def test_owner_must_be_independent(self):
  self.mutate(self.profile["owner_signoff"],lambda v:v.update(subject_sha256="sha256:"+"1"*64))
  with self.assertRaisesRegex(ValueError,"independently owned"): validate(self.profile,self.root)
 def test_owner_signoff_tamper_is_rejected(self):
  self.mutate(self.profile["owner_signoff"],lambda v:v.update(defects_sha256="0"*64))
  with self.assertRaisesRegex(ValueError,"signoff evidence binding"): validate(self.profile,self.root)

if __name__=="__main__": unittest.main()
````

### FILE: `production_admission_gate/provider-admission.template.json`

```yaml
block_id: "SECURE-OPS:provider-admission-template:v1"
operation: CREATE
provenance: AUTHORED
source: "local blocked template for selected-provider production evidence"
license: LicenseRef-Workspace-Owner
sha256: "d8c060845a376621926682b7a3cad61f2b25ec3bec298227ad8b485b4d1a81b9"
variables: []
secrets_allowed: false
```

````json
{
  "schema": "elite-provider-admission/v1",
  "project_id": "replace-me",
  "environment": "production",
  "release_digest": "sha256:replace-me",
  "evaluated_at": "replace-me",
  "expires_at": "replace-me",
  "executor_ref": "replace-me",
  "target": {
    "system_url": "https://replace-me.invalid",
    "jurisdiction": "replace-me"
  },
  "providers": [],
  "output": "evidence/providers-control.json"
}
````

### FILE: `production_admission_gate/validate_provider_admission.py`

```yaml
block_id: "SECURE-OPS:provider-admission-validator:v1"
operation: CREATE
provenance: AUTHORED
source: "local semantic validator over exact provider SDK and official-contract evidence"
license: LicenseRef-Workspace-Owner
sha256: "4854522ac58c3133d0fabfeb38771b056fa312ee225e8ad6039586c811a576b2"
variables: []
secrets_allowed: false
```

````python
from __future__ import annotations
import argparse,json,re,sys
from pathlib import Path
from urllib.parse import urlparse
from validate_k6_load_admission import atomic,canonical_sha,proof,read_json,safe,utc,verify_proof

KINDS=("SANDBOX_CONTRACT","WEBHOOK_AUTH","IDEMPOTENCY_RECONCILIATION","RATE_LIMIT_RETRY")
FIELDS={"schema","project_id","environment","release_digest","evaluated_at","expires_at","executor_ref","target","providers","output"}
TARGET_FIELDS={"system_url","jurisdiction"}
PROVIDER_FIELDS={"provider_id","policy","runs"}
RUN_FIELDS={"kind","execution_receipt","observation","expected_arguments"}
POLICY_FIELDS={"schema","project_id","environment","release_digest","provider_id","approved_at","approvers","adapter_id","source_identity","operations","scopes","account_reference","sandbox_reference","webhook_mode","mutation_mode","terms_url","terms_version","terms_reviewed_at","data_regions","retention_days","quota_cost_approved","finance_owner","security_owner","reconciliation_owner"}
TOOLS={
 "stripe":{"name":"Stripe Go SDK","version":"86.3.0","source":"https://github.com/stripe/stripe-go/tree/a2df585a800a97fe8ec4ebf551b4449bdb3d90a1","source_identity":"stripe-go-86.3.0"},
 "mercadopago":{"name":"Mercado Pago Go SDK","version":"1.14.0","source":"https://github.com/mercadopago/sdk-go/tree/f910ee53fbb6819e435eaf3d0f800cb1fe74ae09","source_identity":"mercadopago-sdk-go-1.14.0"},
 "amazon-spapi":{"name":"Amazon Selling Partner API Python SDK","version":"1.11.1","source":"https://github.com/amzn/selling-partner-api-sdk/tree/8e792ae345a8d334ccdbdd03181f05f040e6a4fc","source_identity":"amazon-selling-partner-api-sdk-python-1.11.1"},
 "google-ads":{"name":"Google Ads Python SDK","version":"31.3.0","source":"https://github.com/googleads/google-ads-python/tree/f7bf312d26904ea50ad9ef02e395edccfaa1ebb7","source_identity":"google-ads-python-31.3.0"},
 "meta-ads":{"name":"Meta Business Python SDK","version":"26.0.1","source":"https://github.com/facebook/facebook-python-business-sdk/tree/788f363d15b1269ab5efb7cd00fb5e3b133cd99b","source_identity":"meta-business-sdk-python-26.0.1"},
 "tiktok-ads":{"name":"TikTok Business API Python SDK","version":"f809c396520df2d7b201a9ccc5378d822b728ed3","source":"https://github.com/tiktok/tiktok-business-api-sdk/tree/f809c396520df2d7b201a9ccc5378d822b728ed3","source_identity":"tiktok-business-api-sdk"},
 "meta-whatsapp":{"name":"Meta WhatsApp Cloud API examples","version":"de70ee908a67026e642aaee3703d20464e2a9466","source":"https://github.com/fbsamples/whatsapp-api-examples/tree/de70ee908a67026e642aaee3703d20464e2a9466","source_identity":"meta-whatsapp-api-examples"},
 "firebase-fcm":{"name":"Firebase Admin Go SDK","version":"4.21.0","source":"https://github.com/firebase/firebase-admin-go/tree/eebb06f2a643fbb59b1cb262874a943584475128","source_identity":"firebase-admin-go-4.21.0"},
 "google-merchant":{"name":"Google Merchant Products Python SDK","version":"1.8.0","source":"https://github.com/googleapis/google-cloud-python/tree/97d7b42cd74b41211f5ec8871cc0dd15debdb1a0/packages/google-shopping-merchant-products","source_identity":"google-shopping-merchant-products-python-1.8.0"},
 "mercadolibre":{"name":"Elite Mercado Libre official-HTTP-contract adapter","version":"0.2.0","source":"https://developers.mercadolibre.com.ar/en_us/categories-and-attributes/manage-questions-and-answers","source_identity":"GO-MERCADOLIBRE-MARKETPLACE-ADAPTER@0.2.0:AUTHORED"},
}

def nonempty(value:object,label:str)->str:
 if not isinstance(value,str) or not value.strip() or "replace-me" in value.lower(): raise ValueError(f"{label} is required")
 return value

def https(value:object,label:str)->None:
 value=nonempty(value,label); parsed=urlparse(value)
 if parsed.scheme!="https" or not parsed.hostname or parsed.username or parsed.password or parsed.fragment: raise ValueError(f"{label} must be credential-free HTTPS")

def strings(value:object,label:str,minimum:int=1)->list[str]:
 if not isinstance(value,list) or len(value)<minimum or len(set(value))!=len(value) or any(not isinstance(x,str) or not x.strip() for x in value): raise ValueError(f"{label} must contain at least {minimum} unique values")
 return value

def bind(value:dict[str,object],profile:dict[str,object],label:str)->None:
 for key in ("project_id","environment","release_digest"):
  if value.get(key)!=profile[key]: raise ValueError(f"{label} {key} mismatch")

def common_observation(value:object,profile:dict[str,object],policy:dict[str,object],kind:str)->dict[str,object]:
 if not isinstance(value,dict): raise ValueError(f"{kind} observation must be an object")
 required={"schema","project_id","environment","release_digest","captured_at","provider_id","kind","account_reference","operations"}
 if not required.issubset(value) or value.get("schema")!="elite-provider-observation/v1" or value.get("kind")!=kind: raise ValueError(f"{kind} observation identity mismatch")
 bind(value,profile,kind); utc(value["captured_at"],f"{kind}.captured_at")
 if value.get("provider_id")!=policy["provider_id"] or value.get("account_reference")!=policy["account_reference"] or value.get("operations")!=policy["operations"]: raise ValueError(f"{kind} provider/account/operations mismatch")
 return value

def sandbox(value:object,profile:dict[str,object],policy:dict[str,object])->None:
 value=common_observation(value,profile,policy,"SANDBOX_CONTRACT"); fields={"schema","project_id","environment","release_digest","captured_at","provider_id","kind","account_reference","operations","sandbox_reference","cases"}
 if set(value)!=fields or value["sandbox_reference"]!=policy["sandbox_reference"]: raise ValueError("sandbox observation fields/reference mismatch")
 cases=value["cases"]
 if not isinstance(cases,list) or {x.get("operation") for x in cases if isinstance(x,dict)}!=set(policy["operations"]): raise ValueError("every approved provider operation requires a sandbox case")
 for item in cases:
  if set(item)!={"operation","contract_pass","provider_response_observed","no_unapproved_production_effect","response_schema_sha256"} or item["contract_pass"] is not True or item["provider_response_observed"] is not True or item["no_unapproved_production_effect"] is not True or not re.fullmatch(r"[0-9a-f]{64}",str(item["response_schema_sha256"])): raise ValueError("sandbox contract case failed")

def webhook(value:object,profile:dict[str,object],policy:dict[str,object])->None:
 value=common_observation(value,profile,policy,"WEBHOOK_AUTH"); fields={"schema","project_id","environment","release_digest","captured_at","provider_id","kind","account_reference","operations","mode","positive_signature_verified","altered_payload_rejected","stale_or_replayed_rejected","duplicate_converged","durable_inbox","missed_event_reconciled","not_applicable_reason","reviewers"}
 if set(value)!=fields or value["mode"]!=policy["webhook_mode"]: raise ValueError("webhook observation fields/mode mismatch")
 if value["mode"]=="REQUIRED":
  if any(value[k] is not True for k in ("positive_signature_verified","altered_payload_rejected","stale_or_replayed_rejected","duplicate_converged","durable_inbox","missed_event_reconciled")) or value["not_applicable_reason"]!="" or value["reviewers"]!=[]: raise ValueError("required webhook authenticity/durability evidence failed")
 elif value["mode"]=="NOT_APPLICABLE":
  if any(value[k] is not False for k in ("positive_signature_verified","altered_payload_rejected","stale_or_replayed_rejected","duplicate_converged","durable_inbox","missed_event_reconciled")) or not nonempty(value["not_applicable_reason"],"webhook not_applicable_reason") or len(strings(value["reviewers"],"webhook reviewers",2))<2: raise ValueError("webhook N/A requires reason and two reviewers")
 else: raise ValueError("webhook mode must be REQUIRED or NOT_APPLICABLE")

def idempotency(value:object,profile:dict[str,object],policy:dict[str,object])->None:
 value=common_observation(value,profile,policy,"IDEMPOTENCY_RECONCILIATION"); fields={"schema","project_id","environment","release_digest","captured_at","provider_id","kind","account_reference","operations","mutation_mode","idempotency_key_or_dedup_used","duplicate_same_remote_identity","remote_effect_count","uncertain_outcome_reconciled","local_state_matches_provider","recovery_or_compensation_pass"}
 if set(value)!=fields or value["mutation_mode"]!=policy["mutation_mode"]: raise ValueError("idempotency observation fields/mode mismatch")
 if value["idempotency_key_or_dedup_used"] is not True or value["duplicate_same_remote_identity"] is not True or value["remote_effect_count"]!=1 or value["uncertain_outcome_reconciled"] is not True or value["local_state_matches_provider"] is not True or value["recovery_or_compensation_pass"] is not True: raise ValueError("idempotency/reconciliation evidence failed")

def rate_limit(value:object,profile:dict[str,object],policy:dict[str,object])->None:
 value=common_observation(value,profile,policy,"RATE_LIMIT_RETRY"); fields={"schema","project_id","environment","release_digest","captured_at","provider_id","kind","account_reference","operations","throttle_observed","retry_after_or_backoff_honored","jitter_or_provider_sdk_retry","max_attempts","attempts_observed","permanent_error_not_retried","exhaustion_durable","within_approved_quota","retry_storm_absent"}
 if set(value)!=fields or value["throttle_observed"] is not True or value["retry_after_or_backoff_honored"] is not True or value["jitter_or_provider_sdk_retry"] is not True or not isinstance(value["max_attempts"],int) or not 2<=value["max_attempts"]<=10 or not isinstance(value["attempts_observed"],int) or not 2<=value["attempts_observed"]<=value["max_attempts"] or any(value[k] is not True for k in ("permanent_error_not_retried","exhaustion_durable","within_approved_quota","retry_storm_absent")): raise ValueError("bounded provider rate-limit/retry evidence failed")

def validate(profile:object,root:Path)->dict[str,object]:
 if not isinstance(profile,dict) or set(profile)!=FIELDS: raise ValueError("profile fields must be exact")
 if profile["schema"]!="elite-provider-admission/v1" or profile["environment"]!="production" or not nonempty(profile["project_id"],"project_id") or not re.fullmatch(r"sha256:[0-9a-f]{64}",str(profile["release_digest"])): raise ValueError("profile identity is invalid")
 start=utc(profile["evaluated_at"],"evaluated_at"); end=utc(profile["expires_at"],"expires_at")
 if not start<end or (end-start).total_seconds()>30*86400: raise ValueError("admission validity must be positive and at most 30 days")
 nonempty(profile["executor_ref"],"executor_ref"); target=profile["target"]
 if not isinstance(target,dict) or set(target)!=TARGET_FIELDS: raise ValueError("target fields must be exact")
 https(target["system_url"],"target.system_url"); nonempty(target["jurisdiction"],"target.jurisdiction")
 providers=profile["providers"]
 if not isinstance(providers,list) or not providers or any(not isinstance(x,dict) or set(x)!=PROVIDER_FIELDS for x in providers): raise ValueError("at least one exact provider entry is required")
 ids=[x["provider_id"] for x in providers]
 if len(ids)!=len(set(ids)) or any(x not in TOOLS for x in ids): raise ValueError("provider IDs must be unique and admitted")
 evidence=[]; tool_hashes=[]; summaries=[]
 validators={"SANDBOX_CONTRACT":sandbox,"WEBHOOK_AUTH":webhook,"IDEMPOTENCY_RECONCILIATION":idempotency,"RATE_LIMIT_RETRY":rate_limit}
 for entry in providers:
  provider_id=entry["provider_id"]; exact=TOOLS[provider_id]; policy_path=safe(root,entry["policy"],f"{provider_id}.policy"); policy=read_json(policy_path,f"{provider_id}.policy")
  if not isinstance(policy,dict) or set(policy)!=POLICY_FIELDS or policy.get("schema")!="elite-provider-policy/v1" or policy.get("provider_id")!=provider_id: raise ValueError(f"{provider_id} policy fields/schema mismatch")
  bind(policy,profile,f"{provider_id}.policy"); approved=utc(policy["approved_at"],f"{provider_id}.approved_at"); reviewed=utc(policy["terms_reviewed_at"],f"{provider_id}.terms_reviewed_at")
  if reviewed>approved or (approved-reviewed).total_seconds()>90*86400: raise ValueError(f"{provider_id} terms review must be no older than 90 days at approval")
  approvers=strings(policy["approvers"],f"{provider_id}.approvers",2)
  for key in ("adapter_id","account_reference","sandbox_reference","terms_version","finance_owner","security_owner","reconciliation_owner"): nonempty(policy[key],f"{provider_id}.{key}")
  if len({policy["finance_owner"],policy["security_owner"],policy["reconciliation_owner"]})<3 or not set((policy["finance_owner"],policy["security_owner"])).issubset(set(approvers)): raise ValueError(f"{provider_id} finance/security/reconciliation ownership is not independent")
  strings(policy["operations"],f"{provider_id}.operations"); strings(policy["scopes"],f"{provider_id}.scopes"); strings(policy["data_regions"],f"{provider_id}.data_regions")
  https(policy["terms_url"],f"{provider_id}.terms_url")
  if policy["source_identity"]!=exact["source_identity"] or policy["webhook_mode"] not in ("REQUIRED","NOT_APPLICABLE") or policy["mutation_mode"] not in ("READ_ONLY","MUTATING") or policy["quota_cost_approved"] is not True or not isinstance(policy["retention_days"],int) or not 1<=policy["retention_days"]<=3650: raise ValueError(f"{provider_id} source/modes/terms/cost policy invalid")
  runs=entry["runs"]
  if not isinstance(runs,list) or len(runs)!=4 or any(not isinstance(x,dict) or set(x)!=RUN_FIELDS for x in runs) or [x["kind"] for x in runs]!=list(KINDS): raise ValueError(f"{provider_id} requires four ordered provider runs")
  evidence.append(proof(root,policy_path)); family_hash=None
  for run in runs:
   kind=run["kind"]; observation_path=safe(root,run["observation"],f"{provider_id}.{kind}.observation"); observation=read_json(observation_path,f"{provider_id}.{kind}.observation"); validators[kind](observation,profile,policy); observation_proof=proof(root,observation_path)
   receipt_path=safe(root,run["execution_receipt"],f"{provider_id}.{kind}.receipt"); receipt=read_json(receipt_path,f"{provider_id}.{kind}.receipt")
   if receipt.get("schema")!="elite-official-tool-execution/v1" or receipt.get("project_id")!=profile["project_id"] or receipt.get("environment")!="production" or receipt.get("release_digest")!=profile["release_digest"] or receipt.get("control_id")!="PROVIDERS" or receipt.get("exit_code")!=0: raise ValueError(f"{provider_id} {kind} execution identity/exit mismatch")
   tool=receipt.get("tool")
   if not isinstance(tool,dict) or set(tool)!={"name","version","source","sha256"} or any(tool.get(k)!=exact[k] for k in ("name","version","source")) or not re.fullmatch(r"[0-9a-f]{64}",str(tool.get("sha256"))): raise ValueError(f"{provider_id} exact official tool identity mismatch")
   if family_hash is None: family_hash=tool["sha256"]
   if family_hash!=tool["sha256"]: raise ValueError(f"{provider_id} runs must use one exact tool/source hash")
   args=run["expected_arguments"]
   if not isinstance(args,list) or not args or any(not isinstance(x,str) or not x for x in args) or receipt.get("arguments_sha256")!=canonical_sha(args): raise ValueError(f"{provider_id} {kind} arguments are not hash-bound")
   expected_target={"system_url":target["system_url"],"jurisdiction":target["jurisdiction"],"provider_id":provider_id,"account_reference":policy["account_reference"],"run_kind":kind,"observation_sha256":observation_proof["sha256"]}
   if receipt.get("target")!=expected_target: raise ValueError(f"{provider_id} {kind} target/observation binding mismatch")
   if not isinstance(receipt.get("environment_variable_names"),list) or any(re.search(r"(?:secret|password|token|key|credential)",str(x),re.I) and not str(x).endswith("_REF") for x in receipt["environment_variable_names"]): raise ValueError(f"{provider_id} environment must expose references, not secret-bearing names")
   for channel in ("stdout","stderr"): verify_proof(root,receipt.get(channel),f"{provider_id}.{kind}.{channel}")
   evidence.extend([proof(root,receipt_path),observation_proof,receipt["stdout"],receipt["stderr"]])
  tool_hashes.append(family_hash); summaries.append({"provider_id":provider_id,"adapter_id":policy["adapter_id"],"source_identity":policy["source_identity"],"account_reference":policy["account_reference"],"operations":policy["operations"],"webhook_mode":policy["webhook_mode"],"mutation_mode":policy["mutation_mode"]})
 return {"id":"PROVIDERS","project_id":profile["project_id"],"environment":"production","release_digest":profile["release_digest"],"result":"PASS","executed_at":profile["evaluated_at"],"expires_at":profile["expires_at"],"executor_ref":profile["executor_ref"],"tool":{"name":"Selected exact provider SDKs and official-contract adapters","version":"policy-locked","source":"OFFICIAL-UPSTREAM-ACQUISITION-CORE:commerce-communications-leaders","digest":"sha256:"+canonical_sha(tool_hashes)},"target":{"system_url":target["system_url"],"jurisdiction":target["jurisdiction"],"providers":summaries},"assertions":{"all_required_sandboxes_pass":True,"webhook_auth_pass":True,"idempotency_reconciliation_pass":True,"rate_limit_retry_pass":True,"terms_cost_scope_approved":True},"evidence":evidence}

def main()->int:
 p=argparse.ArgumentParser(); p.add_argument("--project-root",required=True,type=Path); p.add_argument("--profile",required=True); a=p.parse_args()
 try:
  root=a.project_root.resolve(strict=True); profile=read_json(safe(root,a.profile,"profile"),"profile"); output=safe(root,profile.get("output"),"output",exists=False); output.parent.resolve(strict=True).relative_to(root); receipt=validate(profile,root); atomic(output,(json.dumps(receipt,indent=2,sort_keys=True)+"\n").encode()); print(f"PROVIDER_ADMISSION_PASS project={receipt['project_id']} release={receipt['release_digest']} providers={len(receipt['target']['providers'])}"); return 0
 except (ValueError,OSError) as error: print(f"PROVIDER_ADMISSION_FAILED: {error}",file=sys.stderr); return 1
if __name__=="__main__": sys.exit(main())
````

### FILE: `production_admission_gate/test_validate_provider_admission.py`

```yaml
block_id: "SECURE-OPS:provider-admission-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local positive and fail-closed provider admission regressions"
license: LicenseRef-Workspace-Owner
sha256: "3d0210c4282c0fa9a221ee1b588949fa095ee73dac062c4caf67d00095d32e82"
variables: []
secrets_allowed: false
```

````python
from __future__ import annotations
import copy,hashlib,json,tempfile,unittest
from datetime import datetime,timedelta,timezone
from pathlib import Path
from validate_k6_load_admission import canonical_sha
from validate_provider_admission import KINDS,TOOLS,validate

class ProviderAdmissionTests(unittest.TestCase):
 def setUp(self):
  self.temp=tempfile.TemporaryDirectory(); self.root=Path(self.temp.name); (self.root/"evidence").mkdir(); now=datetime.now(timezone.utc); self.now=now; z=lambda d:d.isoformat().replace("+00:00","Z"); digest="sha256:"+"c"*64
  self.profile={"schema":"elite-provider-admission/v1","project_id":"revestex","environment":"production","release_digest":digest,"evaluated_at":z(now),"expires_at":z(now+timedelta(days=7)),"executor_ref":"providers:change/77","target":{"system_url":"https://app.example.test","jurisdiction":"AR"},"providers":[],"output":"evidence/providers-control.json"}
  self.add_provider("stripe","REQUIRED","MUTATING",["create_payment_intent","receive_payment_event"],"stripe-account-hash")
  self.add_provider("google-merchant","NOT_APPLICABLE","MUTATING",["insert_product","get_product_status"],"merchant-account-hash")
 def tearDown(self): self.temp.cleanup()
 def file(self,p,data): q=self.root/p; q.parent.mkdir(parents=True,exist_ok=True); q.write_bytes(data); return q
 def write(self,p,v): q=self.root/p; q.parent.mkdir(parents=True,exist_ok=True); q.write_text(json.dumps(v),encoding="utf-8"); return q
 def pf(self,p): data=p.read_bytes(); return {"path":p.relative_to(self.root).as_posix(),"bytes":len(data),"sha256":hashlib.sha256(data).hexdigest()}
 def add_provider(self,pid,webhook_mode,mutation_mode,operations,account):
  z=lambda d:d.isoformat().replace("+00:00","Z"); exact=TOOLS[pid]; stem=f"evidence/{pid}"; policy={"schema":"elite-provider-policy/v1","project_id":"revestex","environment":"production","release_digest":self.profile["release_digest"],"provider_id":pid,"approved_at":z(self.now),"approvers":["finance-owner","security-owner"],"adapter_id":f"elite-{pid}-adapter","source_identity":exact["source_identity"],"operations":operations,"scopes":["minimum-required"],"account_reference":account,"sandbox_reference":f"{pid}-sandbox-hash","webhook_mode":webhook_mode,"mutation_mode":mutation_mode,"terms_url":"https://example.test/provider-terms","terms_version":"2026-08","terms_reviewed_at":z(self.now-timedelta(days=1)),"data_regions":["approved-region"],"retention_days":365,"quota_cost_approved":True,"finance_owner":"finance-owner","security_owner":"security-owner","reconciliation_owner":"operations-owner"}; self.write(f"{stem}-policy.json",policy); runs=[]
  base={"schema":"elite-provider-observation/v1","project_id":"revestex","environment":"production","release_digest":self.profile["release_digest"],"captured_at":z(self.now),"provider_id":pid,"account_reference":account,"operations":operations}
  observations={
   "SANDBOX_CONTRACT":dict(base,kind="SANDBOX_CONTRACT",sandbox_reference=policy["sandbox_reference"],cases=[{"operation":op,"contract_pass":True,"provider_response_observed":True,"no_unapproved_production_effect":True,"response_schema_sha256":"d"*64} for op in operations]),
   "WEBHOOK_AUTH":dict(base,kind="WEBHOOK_AUTH",mode=webhook_mode,positive_signature_verified=webhook_mode=="REQUIRED",altered_payload_rejected=webhook_mode=="REQUIRED",stale_or_replayed_rejected=webhook_mode=="REQUIRED",duplicate_converged=webhook_mode=="REQUIRED",durable_inbox=webhook_mode=="REQUIRED",missed_event_reconciled=webhook_mode=="REQUIRED",not_applicable_reason="provider operations expose no callback" if webhook_mode=="NOT_APPLICABLE" else "",reviewers=["security-owner","operations-owner"] if webhook_mode=="NOT_APPLICABLE" else []),
   "IDEMPOTENCY_RECONCILIATION":dict(base,kind="IDEMPOTENCY_RECONCILIATION",mutation_mode=mutation_mode,idempotency_key_or_dedup_used=True,duplicate_same_remote_identity=True,remote_effect_count=1,uncertain_outcome_reconciled=True,local_state_matches_provider=True,recovery_or_compensation_pass=True),
   "RATE_LIMIT_RETRY":dict(base,kind="RATE_LIMIT_RETRY",throttle_observed=True,retry_after_or_backoff_honored=True,jitter_or_provider_sdk_retry=True,max_attempts=4,attempts_observed=3,permanent_error_not_retried=True,exhaustion_durable=True,within_approved_quota=True,retry_storm_absent=True)}
  for kind in KINDS:
   key=kind.lower(); obs=self.write(f"{stem}-{key}.json",observations[kind]); args=["provider-check",pid,key]; out=self.file(f"{stem}-{key}.stdout",b"pass\n"); err=self.file(f"{stem}-{key}.stderr",b""); target=dict(self.profile["target"],provider_id=pid,account_reference=account,run_kind=kind,observation_sha256=self.pf(obs)["sha256"]); receipt={"schema":"elite-official-tool-execution/v1","project_id":"revestex","environment":"production","release_digest":self.profile["release_digest"],"control_id":"PROVIDERS","executed_at":z(self.now),"exit_code":0,"tool":{"name":exact["name"],"version":exact["version"],"source":exact["source"],"sha256":("a" if pid=="stripe" else "b")*64},"arguments_sha256":canonical_sha(args),"working_directory":".","target":target,"environment_variable_names":["PROVIDER_CREDENTIAL_REF"],"stdout":self.pf(out),"stderr":self.pf(err)}; self.write(f"{stem}-{key}-receipt.json",receipt); runs.append({"kind":kind,"execution_receipt":f"{stem}-{key}-receipt.json","observation":f"{stem}-{key}.json","expected_arguments":args})
  self.profile["providers"].append({"provider_id":pid,"policy":f"{stem}-policy.json","runs":runs})
 def observation(self,pindex,rindex): return self.root/self.profile["providers"][pindex]["runs"][rindex]["observation"]
 def receipt(self,pindex,rindex): return self.root/self.profile["providers"][pindex]["runs"][rindex]["execution_receipt"]
 def mutate(self,path,change): v=json.loads(path.read_text()); change(v); self.write(path.relative_to(self.root).as_posix(),v)
 def test_complete_selected_provider_matrix_emits_semantic_receipt(self): self.assertTrue(all(validate(self.profile,self.root)["assertions"].values()))
 def test_missing_selected_provider_is_rejected(self):
  self.profile["providers"]=[]
  with self.assertRaisesRegex(ValueError,"at least one"): validate(self.profile,self.root)
 def test_exact_official_tool_identity_is_required(self):
  self.mutate(self.receipt(0,0),lambda v:v["tool"].update(version="latest"))
  with self.assertRaisesRegex(ValueError,"exact official tool"): validate(self.profile,self.root)
 def test_every_operation_needs_sandbox_case(self):
  self.mutate(self.observation(0,0),lambda v:v.update(cases=v["cases"][:1]))
  with self.assertRaisesRegex(ValueError,"every approved"): validate(self.profile,self.root)
 def test_sandbox_cannot_hide_production_effect(self):
  self.mutate(self.observation(0,0),lambda v:v["cases"][0].update(no_unapproved_production_effect=False))
  with self.assertRaisesRegex(ValueError,"sandbox contract case failed"): validate(self.profile,self.root)
 def test_webhook_tamper_acceptance_is_rejected(self):
  self.mutate(self.observation(0,1),lambda v:v.update(altered_payload_rejected=False))
  with self.assertRaisesRegex(ValueError,"webhook authenticity"): validate(self.profile,self.root)
 def test_webhook_not_applicable_needs_two_reviewers(self):
  self.mutate(self.observation(1,1),lambda v:v.update(reviewers=["security-owner"]))
  with self.assertRaisesRegex(ValueError,"at least 2"): validate(self.profile,self.root)
 def test_duplicate_remote_effect_is_rejected(self):
  self.mutate(self.observation(0,2),lambda v:v.update(remote_effect_count=2))
  with self.assertRaisesRegex(ValueError,"idempotency/reconciliation"): validate(self.profile,self.root)
 def test_uncertain_outcome_must_reconcile(self):
  self.mutate(self.observation(0,2),lambda v:v.update(uncertain_outcome_reconciled=False))
  with self.assertRaisesRegex(ValueError,"idempotency/reconciliation"): validate(self.profile,self.root)
 def test_unbounded_retry_is_rejected(self):
  self.mutate(self.observation(0,3),lambda v:v.update(max_attempts=99,attempts_observed=99))
  with self.assertRaisesRegex(ValueError,"bounded provider"): validate(self.profile,self.root)
 def test_permanent_error_retry_is_rejected(self):
  self.mutate(self.observation(0,3),lambda v:v.update(permanent_error_not_retried=False))
  with self.assertRaisesRegex(ValueError,"bounded provider"): validate(self.profile,self.root)
 def test_stale_terms_review_is_rejected(self):
  p=self.root/self.profile["providers"][0]["policy"]; self.mutate(p,lambda v:v.update(terms_reviewed_at=(self.now-timedelta(days=91)).isoformat().replace("+00:00","Z")))
  with self.assertRaisesRegex(ValueError,"no older than 90"): validate(self.profile,self.root)
 def test_nonzero_execution_is_rejected(self):
  self.mutate(self.receipt(0,0),lambda v:v.update(exit_code=3))
  with self.assertRaisesRegex(ValueError,"identity/exit"): validate(self.profile,self.root)
 def test_secret_bearing_environment_name_is_rejected(self):
  self.mutate(self.receipt(0,0),lambda v:v.update(environment_variable_names=["PROVIDER_ACCESS_TOKEN"]))
  with self.assertRaisesRegex(ValueError,"references, not secret"): validate(self.profile,self.root)

if __name__=="__main__": unittest.main()
````


### FILE: `ops/prometheus/platform.rules.test.yml`

```yaml
block_id: "SECURE-OPS:prom-rule-tests:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "3c5ed3e760f8a53cc9f6e4037e275177e1b6a3976e11e6dc6ca8d53485e3c643"
variables: []
secrets_allowed: false
```

````yaml
{
  "rule_files": [
    "platform.rules.yml"
  ],
  "evaluation_interval": "30s",
  "tests": [
    {
      "name": "low-traffic-all-errors",
      "interval": "1m",
      "input_series": [
        {
          "series": "http_server_request_duration_seconds_count{http_response_status_code=\"500\"}",
          "values": "0 0 0 0 1 1 1 1 2 2 2 2 3 3 3 3 4 4 4 4 5 5 5 5 6 6 6 6 7 7 7 7 8 8 8 8 9 9 9 9 10"
        }
      ],
      "alert_rule_test": [
        {
          "eval_time": "10m",
          "alertname": "PlatformHighErrorRate",
          "exp_alerts": []
        },
        {
          "eval_time": "20m",
          "alertname": "PlatformHighErrorRate",
          "exp_alerts": [
            {
              "exp_labels": {
                "severity": "page"
              },
              "exp_annotations": {
                "summary": "User-visible server error rate exceeds one percent",
                "runbook_url": "https://replace.invalid/runbooks/api-availability"
              }
            }
          ]
        }
      ]
    },
    {
      "name": "high-traffic-all-errors",
      "interval": "1m",
      "input_series": [
        {
          "series": "http_server_request_duration_seconds_count{http_response_status_code=\"500\"}",
          "values": "0+600x40"
        }
      ],
      "alert_rule_test": [
        {
          "eval_time": "4m",
          "alertname": "PlatformHighErrorRate",
          "exp_alerts": []
        },
        {
          "eval_time": "20m",
          "alertname": "PlatformHighErrorRate",
          "exp_alerts": [
            {
              "exp_labels": {
                "severity": "page"
              },
              "exp_annotations": {
                "summary": "User-visible server error rate exceeds one percent",
                "runbook_url": "https://replace.invalid/runbooks/api-availability"
              }
            }
          ]
        }
      ]
    },
    {
      "name": "low-traffic-healthy",
      "interval": "1m",
      "input_series": [
        {
          "series": "http_server_request_duration_seconds_count{http_response_status_code=\"200\"}",
          "values": "0 0 0 0 1 1 1 1 2 2 2 2 3 3 3 3 4 4 4 4 5 5 5 5 6 6 6 6 7 7 7 7 8 8 8 8 9 9 9 9 10"
        }
      ],
      "alert_rule_test": [
        {
          "eval_time": "20m",
          "alertname": "PlatformHighErrorRate",
          "exp_alerts": []
        }
      ]
    },
    {
      "name": "no-traffic",
      "interval": "1m",
      "input_series": [
        {
          "series": "http_server_request_duration_seconds_count{http_response_status_code=\"500\"}",
          "values": "0+0x40"
        },
        {
          "series": "http_server_request_duration_seconds_count{http_response_status_code=\"200\"}",
          "values": "0+0x40"
        }
      ],
      "alert_rule_test": [
        {
          "eval_time": "20m",
          "alertname": "PlatformHighErrorRate",
          "exp_alerts": []
        }
      ]
    },
    {
      "name": "exact-one-percent",
      "interval": "1m",
      "input_series": [
        {
          "series": "http_server_request_duration_seconds_count{http_response_status_code=\"500\"}",
          "values": "0+1x40"
        },
        {
          "series": "http_server_request_duration_seconds_count{http_response_status_code=\"200\"}",
          "values": "0+99x40"
        }
      ],
      "alert_rule_test": [
        {
          "eval_time": "20m",
          "alertname": "PlatformHighErrorRate",
          "exp_alerts": []
        }
      ]
    },
    {
      "name": "above-one-percent",
      "interval": "1m",
      "input_series": [
        {
          "series": "http_server_request_duration_seconds_count{http_response_status_code=\"500\"}",
          "values": "0+2x40"
        },
        {
          "series": "http_server_request_duration_seconds_count{http_response_status_code=\"200\"}",
          "values": "0+98x40"
        }
      ],
      "alert_rule_test": [
        {
          "eval_time": "20m",
          "alertname": "PlatformHighErrorRate",
          "exp_alerts": [
            {
              "exp_labels": {
                "severity": "page"
              },
              "exp_annotations": {
                "summary": "User-visible server error rate exceeds one percent",
                "runbook_url": "https://replace.invalid/runbooks/api-availability"
              }
            }
          ]
        }
      ]
    },
    {
      "name": "counter-reset-still-all-errors",
      "interval": "1m",
      "input_series": [
        {
          "series": "http_server_request_duration_seconds_count{http_response_status_code=\"500\"}",
          "values": "0 0 0 0 1 1 1 1 2 2 2 2 3 3 3 0 0 0 0 1 1 1 1 2 2 2 2 3 3 3 3 4 4 4 4 5 5 5 5 6 6"
        }
      ],
      "alert_rule_test": [
        {
          "eval_time": "35m",
          "alertname": "PlatformHighErrorRate",
          "exp_alerts": [
            {
              "exp_labels": {
                "severity": "page"
              },
              "exp_annotations": {
                "summary": "User-visible server error rate exceeds one percent",
                "runbook_url": "https://replace.invalid/runbooks/api-availability"
              }
            }
          ]
        }
      ]
    },
    {
      "name": "short-spike-and-recovery",
      "interval": "1m",
      "input_series": [
        {
          "series": "http_server_request_duration_seconds_count{http_response_status_code=\"500\"}",
          "values": "0 20 40 60 80 100 100 100 100 100 100 100 100 100 100 100 100 100 100 100 100 100 100 100 100 100 100 100 100 100 100 100 100 100 100 100 100 100 100 100 100"
        },
        {
          "series": "http_server_request_duration_seconds_count{http_response_status_code=\"200\"}",
          "values": "0+80x40"
        }
      ],
      "alert_rule_test": [
        {
          "eval_time": "4m",
          "alertname": "PlatformHighErrorRate",
          "exp_alerts": []
        },
        {
          "eval_time": "20m",
          "alertname": "PlatformHighErrorRate",
          "exp_alerts": []
        }
      ]
    },
    {
      "name": "recovery-after-firing",
      "interval": "1m",
      "input_series": [
        {
          "series": "http_server_request_duration_seconds_count{http_response_status_code=\"500\"}",
          "values": "0 0 0 0 1 1 1 1 2 2 2 2 3 3 3 3 4 4 4 4 5 5 5 5 5 5 5 5 5 5 5 5 5 5 5 5 5 5 5 5 5 5 5 5 5 5 5 5 5 5 5 5 5 5 5 5 5 5 5 5 5"
        },
        {
          "series": "http_server_request_duration_seconds_count{http_response_status_code=\"200\"}",
          "values": "0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 60 120 180 240 300 360 420 480 540 600 660 720 780 840 900 960 1020 1080 1140 1200 1260 1320 1380 1440 1500 1560 1620 1680 1740 1800 1860 1920 1980 2040 2100 2160 2220 2280 2340 2400"
        }
      ],
      "alert_rule_test": [
        {
          "eval_time": "20m",
          "alertname": "PlatformHighErrorRate",
          "exp_alerts": [
            {
              "exp_labels": {
                "severity": "page"
              },
              "exp_annotations": {
                "summary": "User-visible server error rate exceeds one percent",
                "runbook_url": "https://replace.invalid/runbooks/api-availability"
              }
            }
          ]
        },
        {
          "eval_time": "30m",
          "alertname": "PlatformHighErrorRate",
          "exp_alerts": []
        }
      ]
    }
  ]
}
````

## 6. Configuration surface

The operational readiness JSON and the separate production-admission JSON are typed configuration surfaces. The latter contains exactly eight project/environment/release-bound controls and accepts only fresh PASS receipts with explicit tool/target/assertions plus local evidence bytes and SHA-256. Provider adds four ordered runs per selection; deploy adds seven ordered kubectl-bound runs; business acceptance adds exact Spec Kit artifacts, scenario results, three functional approvals, defect register and independent owner signoff. Their distributed templates omit all target authority and cannot pass. Placeholders, zero or mutable digests, plaintext secrets, missing owners, stale evidence, target drift, duplicate controls or tampering fail validation. Credentials and raw identities never enter these files.

## 7. Dependency bill

| Tool/component | Pin/baseline | Use | License | Scope | Source |
|---|---|---|---|---|---|
| Python | `3.14+` stdlib | readiness validator/tests | PSF-2.0 | verification | `python.org` |
| OpenTelemetry Collector | `0.158.0` evidence baseline | telemetry config validation | Apache-2.0 | operations | `opentelemetry.io` |
| Prometheus `promtool` | `3.13.1` evidence baseline | rules validation | Apache-2.0 | operations | `prometheus.io` |
| CycloneDX / SLSA | `1.7` / `1.2` contracts | evidence formats | specification terms | release | official specifications |
| Stripe/Mercado Pago/Amazon/Google/Meta/TikTok/Firebase/Mercado Libre | exact identities in provider registry and acquisition lock | provider sandbox/webhook/idempotency/retry evidence; source not embedded | upstream terms plus local glue license | provider admission | official repos/docs listed in metadata and implementation packs |
| Cloudflare/OpenID/PostgreSQL/k6/ZAP | project-selected exact identities | produce target receipts; not embedded | upstream terms | production admission | official APIs/docs/repos listed in metadata |
| Kubernetes kubectl | `1.37.0`; Windows amd64 SHA-256 `4721b614...` | execute and bind real Deployment apply/rollout/rollback evidence | Apache-2.0 | target cluster; binary not embedded | official release/source/checksum listed in metadata and acquisition lock |
| GitHub Spec Kit + Microsoft Playwright | `1.0.1` / `9118ed15...` + `1.62.1` / `26a9e470...` | bind approved artifacts and execute business acceptance scenarios | MIT + Apache-2.0 | specification/test; sources not embedded here | exact identities in acquisition lock |

## 8. Apply order

Materialize with the application, replace examples through project-owned configuration, validate readiness and telemetry rules, select only required providers, bind each to its exact source/adapter identity and execute its four real evidence lanes. Then generate actual SBOM/provenance/scans, prove migration/restore/rollback and promote the immutable digest through canary gates. After target execution, build the eight receipts and run the production validator; `READY_TO_BUILD` is never promoted by inference. Existing records are merged, never overwritten blindly. Rollback follows the exact prior digest and compatible schema plan.

## 9. Verification

The six operational unit tests and 101 production-admission/official-runner/semantic-adapter cases execute from clean materialization; on this Windows host 100 passed and the symlink-only case skipped because creating symlinks is not authorized. Business-acceptance negatives cover missing capability scenarios, failed or effect-free journeys, observation/receipt tamper, nonzero or wrong Playwright execution, secret-bearing environment names, missing or duplicate functional authority, accepted-open Sev1/Sev2, unreviewed exceptions and non-independent/tampered owner signoff. Deploy negatives from V112 remain intact. These are contract regressions only: no project user, financial effect, production cluster or owner decision was exercised. Examples intentionally remain blocked; an agent must never edit a validator to admit a project.

Production remains conditioned on actual secret provider, PKI/mTLS boundary decision, telemetry backend, executable alerts/runbooks, build platform attestations, scanner results, deployment adapter and recovery evidence.

## 10. Reconstruction evidence

The V1 evidence remains the operational baseline. V113 extends V112 with the three-file business-acceptance adapter. The final V113 audit must recheck version 1.1.0, 37/37 materialized files, 101 production tests and exact Spec Kit/Playwright identities. All eight production controls now have semantic receipt producers. This proves reusable admission machinery, not user acceptance or a production-ready absent project.

V112 extended V111 with the three-file Kubernetes deploy/rollback adapter. Its audit rechecked version 1.0.0, 34/34 materialized files, 87 production tests and exact Kubernetes/source identities.

V111 extended the original six-file pack with four production-admission files, three generic official-tool files and six three-file semantic adapters for Grafana k6, PostgreSQL, OWASP ZAP, Cloudflare, identity/authorization and selected providers. Its audit rechecked version 0.9.0, 31/31 materialized files, 73 production tests and exact provider/source identities.
