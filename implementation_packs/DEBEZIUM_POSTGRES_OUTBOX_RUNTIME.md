# Debezium PostgreSQL Outbox Runtime

## 1. Metadata

```yaml
pack_id: "DEBEZIUM-POSTGRES-OUTBOX-RUNTIME"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa lock OCI por plataforma, configuración Event Router fail-closed, preflight PostgreSQL, contrato consumer inbox y registro HTTPS para Debezium 3.6.1.Final sobre la outbox documental append-only."
stacks: ["Red Hat Debezium 3.6.1.Final", "Kafka Connect", "PostgreSQL logical replication", "Quay OCI", "Python 3.12+", "PowerShell 7"]
compatible_with: ["AWS-IDP-POSTGRES-PERSISTENCE-BOUNDARY 0.2.x", "OFFICIAL-UPSTREAM-ACQUISITION-CORE 0.4.x", "PG-TX-FOUNDATION 0.1.x"]
incompatible_with: ["tag OCI sin digest", "publication automática", "wildcard table include", "TLS deshabilitado", "secret inline", "outbox actualizable", "consumer sin inbox", "exactly-once no demostrado"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0 AND MIT"
upstream_sources: ["https://github.com/debezium/debezium/tree/63371830ceff75437f95c775b91eeeaf55a056c1", "https://github.com/debezium/container-images", "https://quay.io/repository/debezium/connect", "https://debezium.io/documentation/reference/3.6/transformations/outbox-event-router.html"]
verified_at: "2026-08-28"
```

## 2. Applicability

Use con boundary 0.2.0, PostgreSQL logical WAL, publicación exacta, Kafka Connect/broker seguros, secret/CA read-only mounts y consumer inbox idempotente. Es configuración ejecutable sobre runtime oficial fijado; no incluye broker/database, no hace pull sin autorización y no promete exactly-once.

## 3. Architecture contract

La imagen `quay.io/debezium/connect:3.6.1.Final` se identifica exclusivamente por digest OCI index y manifest linux/amd64 o linux/arm64. El conector usa pgoutput, publicación precreada, snapshot initial, slot no drop, tabla única `document_intelligence.outbox_event`, TLS verify-full, secrets file-provider y SMT oficial EventRouter con campos default. Invalid UPDATE es fatal. Tenant/schema/eventType van en headers, aggregateid es key. Register falla ante configuración divergente salvo update explícitamente aprobado; consumer contract exige at-least-once/inbox/hash.

## 4. Exact file manifest

```text
CREATE debezium_postgres_outbox_runtime/README.md
CREATE debezium_postgres_outbox_runtime/runtime-lock.json
CREATE debezium_postgres_outbox_runtime/connector.template.json
CREATE debezium_postgres_outbox_runtime/project-profile.template.json
CREATE debezium_postgres_outbox_runtime/consumer-contract.json
CREATE debezium_postgres_outbox_runtime/postgres_preflight.sql
CREATE debezium_postgres_outbox_runtime/validate_config.py
CREATE debezium_postgres_outbox_runtime/test_validate_config.py
CREATE debezium_postgres_outbox_runtime/register_connector.ps1
CREATE debezium_postgres_outbox_runtime/verify_contract.ps1
```

## 5. Materialization blocks

### FILE: `debezium_postgres_outbox_runtime/README.md`
```yaml
block_id: "DEBEZIUM-POSTGRES-OUTBOX-RUNTIME:readme:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://debezium-postgres-outbox-runtime/README.md"
license: "LicenseRef-Workspace-Owner"
sha256: "ea5b4e9dc3ed15acce69abd6aa4485d8918cbdab6dddb2d6538431f6be189d3b"
variables: []
secrets_allowed: false
```
````markdown
# Debezium PostgreSQL outbox runtime

This conditioned component configures Red Hat/Debezium 3.6.1.Final to capture only `document_intelligence.outbox_event`, whose columns match the official outbox Event Router contract. The official Quay image is fixed by OCI index plus linux/amd64 and linux/arm64 manifest/config digests; the tag is documentary and is never the execution identity. Debezium source release/commit/archive and 48 focused Event Router tests remain fixed by the library acquisition core.

The connector uses `pgoutput`, an exact pre-created publication, an exact table include list, TLS verify-full, file-provider secret references, initial snapshot, retained slot, fatal invalid-operation behavior, JSON payload expansion and deterministic routing. It adds `tenantId`, `schemaVersion` and `eventType` headers and uses `aggregateid` as the Kafka key for partition ordering.

Before registration:

1. run the boundary 0.2 migration as DBA; grant the connector login membership in NOLOGIN `elite_outbox_cdc` and no write privilege;
2. enable logical WAL, create the exact publication/table and durable slot policy, then run `postgres_preflight.sql`;
3. run the locked image digest on the selected architecture with Kafka Connect file ConfigProvider, read-only secret/CA mounts, TLS and broker authentication;
4. complete every profile approval/evidence field and render with `validate_config.py`;
5. register through HTTPS using `register_connector.ps1`; auth comes only from a named environment variable and divergent existing config fails unless an approved update is explicit;
6. prove initial snapshot, duplicate/redelivery, restart, slot retention, ordering, load, Kafka outage/recovery, schema evolution, monitoring and rollback.

Delivery is at-least-once. Every consumer must use header `id` plus tenant-scoped inbox and payload hash in the same transaction as its business effect. This pack does not claim exactly-once, does not delete/mark outbox rows, and is not an immutable audit ledger. Live PostgreSQL/Kafka/Quay pulls and costs remain project-authorized operations.
````

### FILE: `debezium_postgres_outbox_runtime/runtime-lock.json`
```yaml
block_id: "DEBEZIUM-POSTGRES-OUTBOX-RUNTIME:runtime-lock:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://debezium-postgres-outbox-runtime/runtime-lock.json"
license: "LicenseRef-Workspace-Owner"
sha256: "85694132ffe5452b2d846c067c56be46b76f759cd4f98c67c335a52f1e46c139"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-debezium-runtime-lock/v1",
  "verifiedAt": "2026-08-28",
  "image": {
    "repository": "quay.io/debezium/connect",
    "tagObserved": "3.6.1.Final",
    "indexDigest": "sha256:76db18f20116557b2e550d5844d2bbebc927cd9c215ec2271c4ab713cc74f386",
    "mediaType": "application/vnd.oci.image.index.v1+json",
    "platforms": {
      "linux/amd64": {
        "manifestDigest": "sha256:1a39c97202ef3da294f9775a935d113e7a2ec7c0ab150026507682d014e99990",
        "configDigest": "sha256:f70c6880d091a2533bb7012b8a498af23720e4d6e7df732f080d5234673faaef",
        "compressedLayerBytes": 1012256126
      },
      "linux/arm64": {
        "manifestDigest": "sha256:b988179c42d05339aa37c6eaab3439481a6ad11b3531760f3228e4e7725f9e96",
        "configDigest": "sha256:0e2152faa0211397c52af862f27099da60428dd81cff5a9b6605a60b943b8270",
        "compressedLayerBytes": 991323851
      }
    }
  },
  "source": {
    "repository": "debezium/debezium",
    "release": "v3.6.1.Final",
    "commit": "63371830ceff75437f95c775b91eeeaf55a056c1",
    "archiveSha256": "3cc98e463d57b86029d991480893df4adf05577925137b4f465315ee0c94a938",
    "license": "Apache-2.0",
    "focusedTests": 48
  },
  "officialDocumentation": "https://debezium.io/documentation/reference/3.6/transformations/outbox-event-router.html"
}
````

### FILE: `debezium_postgres_outbox_runtime/connector.template.json`
```yaml
block_id: "DEBEZIUM-POSTGRES-OUTBOX-RUNTIME:connector-template:v1"
operation: CREATE
provenance: ADAPTED
source: "https://debezium.io/documentation/reference/3.6/transformations/outbox-event-router.html"
license: "Apache-2.0 AND LicenseRef-Workspace-Owner"
sha256: "dc7032f377c2b1239a471dfb163cd2a8edc88a2ca21fcc2e1318d1a07df20dd6"
variables: []
secrets_allowed: false
```
````json
{
  "name": "REQUIRED",
  "config": {
    "connector.class": "io.debezium.connector.postgresql.PostgresConnector",
    "plugin.name": "pgoutput",
    "database.hostname": "${file:/run/secrets/elite.properties:db.hostname}",
    "database.port": "${file:/run/secrets/elite.properties:db.port}",
    "database.user": "${file:/run/secrets/elite.properties:db.user}",
    "database.password": "${file:/run/secrets/elite.properties:db.password}",
    "database.dbname": "${file:/run/secrets/elite.properties:db.name}",
    "database.sslmode": "verify-full",
    "database.sslrootcert": "/run/secrets/ca.pem",
    "topic.prefix": "REQUIRED",
    "slot.name": "REQUIRED",
    "slot.drop.on.stop": "false",
    "publication.name": "REQUIRED",
    "publication.autocreate.mode": "disabled",
    "snapshot.mode": "initial",
    "schema.include.list": "document_intelligence",
    "table.include.list": "document_intelligence.outbox_event",
    "tombstones.on.delete": "false",
    "provide.transaction.metadata": "true",
    "heartbeat.interval.ms": "REQUIRED",
    "transforms": "outbox",
    "transforms.outbox.type": "io.debezium.transforms.outbox.EventRouter",
    "transforms.outbox.table.field.event.id": "id",
    "transforms.outbox.table.field.event.key": "aggregateid",
    "transforms.outbox.table.field.event.payload": "payload",
    "transforms.outbox.table.field.event.timestamp": "occurred_at",
    "transforms.outbox.route.by.field": "aggregatetype",
    "transforms.outbox.route.topic.regex": "(?<routedByValue>.*)",
    "transforms.outbox.route.topic.replacement": "REQUIRED",
    "transforms.outbox.table.fields.additional.placement": "tenant_id:header:tenantId,schema_version:header:schemaVersion,type:header:eventType",
    "transforms.outbox.table.expand.json.payload": "true",
    "transforms.outbox.table.op.invalid.behavior": "fatal",
    "transforms.outbox.route.tombstone.on.empty.payload": "false"
  }
}
````

### FILE: `debezium_postgres_outbox_runtime/project-profile.template.json`
```yaml
block_id: "DEBEZIUM-POSTGRES-OUTBOX-RUNTIME:project-profile:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://debezium-postgres-outbox-runtime/project-profile.template.json"
license: "LicenseRef-Workspace-Owner"
sha256: "4530db3e031b6ca9eab9c541c4062061d0c8d5bc04bd6ddc85053a77b1634f6e"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-debezium-postgres-outbox-profile/v1",
  "acknowledgeConditionedState": false,
  "approvedBy": "REQUIRED",
  "approvedAt": "REQUIRED",
  "platform": "REQUIRED",
  "imageDigest": "REQUIRED",
  "connectUrl": "REQUIRED",
  "connectorName": "REQUIRED",
  "topicPrefix": "REQUIRED",
  "routeTopicReplacement": "REQUIRED",
  "slotName": "REQUIRED",
  "publicationName": "REQUIRED",
  "heartbeatIntervalMs": 0,
  "secretFileMountedEvidence": "REQUIRED",
  "tlsEvidence": "REQUIRED",
  "kafkaAuthenticationEvidence": "REQUIRED",
  "postgresLogicalReplicationEvidence": "REQUIRED",
  "publicationExactTableEvidence": "REQUIRED",
  "cdcRoleMembershipEvidence": "REQUIRED",
  "consumerInboxEvidence": "REQUIRED",
  "duplicateReplayEvidence": "REQUIRED",
  "restartOrderingLoadEvidence": "REQUIRED",
  "recoveryEvidence": "REQUIRED",
  "rollbackEvidence": "REQUIRED",
  "observabilityEvidence": "REQUIRED",
  "costApproval": "REQUIRED"
}
````

### FILE: `debezium_postgres_outbox_runtime/consumer-contract.json`
```yaml
block_id: "DEBEZIUM-POSTGRES-OUTBOX-RUNTIME:consumer-contract:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://debezium-postgres-outbox-runtime/consumer-contract.json"
license: "LicenseRef-Workspace-Owner"
sha256: "40defc2adc0354febfbdedeac8b9d1ecaa9cb6ad7a9eec0ec86364776be9b8d4"
variables: []
secrets_allowed: false
```
````json
{
  "schema": "elite-debezium-outbox-consumer-contract/v1",
  "delivery": "AT_LEAST_ONCE",
  "requiredHeaders": ["id", "tenantId", "schemaVersion", "eventType"],
  "messageKey": "aggregateid",
  "requiredBehavior": [
    "insert header id into tenant-scoped consumer inbox before business effect in the same transaction",
    "same id and same payload hash returns the stored result",
    "same id and different payload hash is a terminal conflict",
    "commit inbox and business effect atomically",
    "retry transient infrastructure failures",
    "dead-letter terminal poison messages with redacted evidence",
    "reconcile unknown outcomes before retrying an external effect"
  ],
  "exactlyOnceClaimAllowed": false
}
````

### FILE: `debezium_postgres_outbox_runtime/postgres_preflight.sql`
```yaml
block_id: "DEBEZIUM-POSTGRES-OUTBOX-RUNTIME:postgres-preflight:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://debezium-postgres-outbox-runtime/postgres_preflight.sql"
license: "LicenseRef-Workspace-Owner"
sha256: "bd33872fda044f9ddd7ac8698b21374aba486dbf892e31161b04269bdb03ca6b"
variables: []
secrets_allowed: false
```
````sql
select current_setting('server_version_num')::integer >= 150000 as supported_version;
select current_setting('wal_level') = 'logical' as logical_wal;
select exists(select 1 from pg_publication where pubname = :'publication_name') as publication_exists;
select exists(
  select 1
  from pg_publication_tables
  where pubname = :'publication_name'
    and schemaname = 'document_intelligence'
    and tablename = 'outbox_event'
) as exact_outbox_published;
select pg_has_role(:'connector_role', 'elite_outbox_cdc', 'member') as cdc_membership;
select has_table_privilege(:'connector_role', 'document_intelligence.outbox_event', 'SELECT') as outbox_select;
````

### FILE: `debezium_postgres_outbox_runtime/validate_config.py`
```yaml
block_id: "DEBEZIUM-POSTGRES-OUTBOX-RUNTIME:validate-config:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://debezium-postgres-outbox-runtime/validate_config.py"
license: "LicenseRef-Workspace-Owner"
sha256: "54e12314265da2447aa4954595f259905d71179e236ebc398ba4433a9c7299cc"
variables: []
secrets_allowed: false
```
````python
from __future__ import annotations

import argparse
import copy
import json
import re
from datetime import datetime
from pathlib import Path
from typing import Any


DIGEST = re.compile(r"^sha256:[0-9a-f]{64}$")
TOKEN = re.compile(r"^[a-z][a-z0-9._-]{0,254}$")
INDEX_DIGEST = "sha256:76db18f20116557b2e550d5844d2bbebc927cd9c215ec2271c4ab713cc74f386"
PLATFORM_DIGESTS = {
    "linux/amd64": ("sha256:1a39c97202ef3da294f9775a935d113e7a2ec7c0ab150026507682d014e99990", "sha256:f70c6880d091a2533bb7012b8a498af23720e4d6e7df732f080d5234673faaef"),
    "linux/arm64": ("sha256:b988179c42d05339aa37c6eaab3439481a6ad11b3531760f3228e4e7725f9e96", "sha256:0e2152faa0211397c52af862f27099da60428dd81cff5a9b6605a60b943b8270"),
}
REQUIRED_PROFILE = {
    "approvedBy", "approvedAt", "platform", "imageDigest", "connectUrl", "connectorName",
    "topicPrefix", "routeTopicReplacement", "slotName", "publicationName",
    "secretFileMountedEvidence", "tlsEvidence", "kafkaAuthenticationEvidence",
    "postgresLogicalReplicationEvidence", "publicationExactTableEvidence",
    "cdcRoleMembershipEvidence", "consumerInboxEvidence", "duplicateReplayEvidence",
    "restartOrderingLoadEvidence", "recoveryEvidence", "rollbackEvidence",
    "observabilityEvidence", "costApproval",
}


class Rejected(ValueError):
    pass


def load(path: Path) -> dict[str, Any]:
    def pairs(items: list[tuple[str, Any]]) -> dict[str, Any]:
        value: dict[str, Any] = {}
        for key, item in items:
            if key in value:
                raise Rejected(f"duplicate JSON key: {key}")
            value[key] = item
        return value
    value = json.loads(path.read_text(encoding="utf-8"), object_pairs_hook=pairs)
    if not isinstance(value, dict):
        raise Rejected(f"not an object: {path}")
    return value


def validate_lock(lock: dict[str, Any]) -> None:
    image = lock.get("image")
    source = lock.get("source")
    if lock.get("schema") != "elite-debezium-runtime-lock/v1" or not isinstance(image, dict) or not isinstance(source, dict):
        raise Rejected("runtime lock identity mismatch")
    if image.get("repository") != "quay.io/debezium/connect" or image.get("tagObserved") != "3.6.1.Final" or image.get("indexDigest") != INDEX_DIGEST:
        raise Rejected("runtime image lock is invalid")
    platforms = image.get("platforms")
    if not isinstance(platforms, dict) or set(platforms) != set(PLATFORM_DIGESTS):
        raise Rejected("runtime platforms are not exact")
    for name, platform in platforms.items():
        expected_manifest, expected_config = PLATFORM_DIGESTS[name]
        if not isinstance(platform, dict) or platform.get("manifestDigest") != expected_manifest or platform.get("configDigest") != expected_config or int(platform.get("compressedLayerBytes", 0)) <= 0:
            raise Rejected("runtime platform lock is invalid")
    if source.get("repository") != "debezium/debezium" or source.get("release") != "v3.6.1.Final" or source.get("commit") != "63371830ceff75437f95c775b91eeeaf55a056c1" or source.get("focusedTests") != 48:
        raise Rejected("source authority drifted")


def validate_template(template: dict[str, Any]) -> None:
    config = template.get("config")
    if template.get("name") != "REQUIRED" or not isinstance(config, dict):
        raise Rejected("connector template must begin unresolved")
    exact = {
        "connector.class": "io.debezium.connector.postgresql.PostgresConnector",
        "plugin.name": "pgoutput", "database.sslmode": "verify-full",
        "publication.autocreate.mode": "disabled", "snapshot.mode": "initial",
        "schema.include.list": "document_intelligence",
        "table.include.list": "document_intelligence.outbox_event",
        "slot.drop.on.stop": "false", "tombstones.on.delete": "false",
        "transforms": "outbox", "transforms.outbox.type": "io.debezium.transforms.outbox.EventRouter",
        "transforms.outbox.table.field.event.id": "id",
        "transforms.outbox.table.field.event.key": "aggregateid",
        "transforms.outbox.table.field.event.payload": "payload",
        "transforms.outbox.table.field.event.timestamp": "occurred_at",
        "transforms.outbox.route.by.field": "aggregatetype",
        "transforms.outbox.table.expand.json.payload": "true",
        "transforms.outbox.table.op.invalid.behavior": "fatal",
        "transforms.outbox.route.tombstone.on.empty.payload": "false",
    }
    for key, expected in exact.items():
        if config.get(key) != expected:
            raise Rejected(f"connector invariant drifted: {key}")
    secret_keys = ("database.hostname", "database.port", "database.user", "database.password", "database.dbname")
    if any(not str(config.get(key, "")).startswith("${file:/run/secrets/elite.properties:") for key in secret_keys):
        raise Rejected("database secrets must remain file-provider references")


def render(template: dict[str, Any], profile: dict[str, Any], lock: dict[str, Any]) -> dict[str, Any]:
    validate_lock(lock)
    validate_template(template)
    if profile.get("schema") != "elite-debezium-postgres-outbox-profile/v1" or profile.get("acknowledgeConditionedState") is not True:
        raise Rejected("project profile is not approved")
    for key in REQUIRED_PROFILE:
        value = profile.get(key)
        if not isinstance(value, str) or not value.strip() or value == "REQUIRED":
            raise Rejected(f"unresolved profile field: {key}")
    try:
        approved = datetime.fromisoformat(profile["approvedAt"].replace("Z", "+00:00"))
    except ValueError as error:
        raise Rejected("approval time is invalid") from error
    if approved.tzinfo is None:
        raise Rejected("approval time has no timezone")
    platform = profile["platform"]
    platforms = lock["image"]["platforms"]
    if platform not in platforms or profile["imageDigest"] != platforms[platform]["manifestDigest"]:
        raise Rejected("runtime platform digest mismatch")
    for key in ("connectorName", "topicPrefix", "slotName", "publicationName"):
        if not TOKEN.fullmatch(profile[key]):
            raise Rejected(f"invalid connector token: {key}")
    heartbeat = profile.get("heartbeatIntervalMs")
    if not isinstance(heartbeat, int) or not 1000 <= heartbeat <= 300000:
        raise Rejected("heartbeat interval is invalid")
    route = profile["routeTopicReplacement"]
    if route != "outbox.event.${routedByValue}":
        raise Rejected("route topic replacement is not admitted")
    output = copy.deepcopy(template)
    output["name"] = profile["connectorName"]
    config = output["config"]
    config["topic.prefix"] = profile["topicPrefix"]
    config["slot.name"] = profile["slotName"]
    config["publication.name"] = profile["publicationName"]
    config["heartbeat.interval.ms"] = str(heartbeat)
    config["transforms.outbox.route.topic.replacement"] = route
    return output


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--component-root", type=Path, default=Path(__file__).parent)
    parser.add_argument("--profile", type=Path, required=True)
    parser.add_argument("--output", type=Path)
    args = parser.parse_args()
    root = args.component_root.resolve()
    result = render(load(root / "connector.template.json"), load(args.profile.resolve()), load(root / "runtime-lock.json"))
    data = json.dumps(result, sort_keys=True, separators=(",", ":"), ensure_ascii=False) + "\n"
    if args.output:
        target = args.output.resolve()
        if target.exists():
            raise Rejected("output already exists")
        target.write_text(data, encoding="utf-8", newline="\n")
    else:
        print(data, end="")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
````

### FILE: `debezium_postgres_outbox_runtime/test_validate_config.py`
```yaml
block_id: "DEBEZIUM-POSTGRES-OUTBOX-RUNTIME:test-validate-config:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://debezium-postgres-outbox-runtime/test_validate_config.py"
license: "LicenseRef-Workspace-Owner"
sha256: "6500b6e8fcdb2bf0bcb1c539107e3f99c7b634f4a55d5121f6d62a79d2d2d5b8"
variables: []
secrets_allowed: false
```
````python
import copy,json,tempfile,unittest
from pathlib import Path
import validate_config as subject


ROOT=Path(__file__).parent


def approved():
    value=json.loads((ROOT/"project-profile.template.json").read_text())
    value.update({"acknowledgeConditionedState":True,"approvedBy":"owner","approvedAt":"2026-08-28T10:00:00Z","platform":"linux/amd64","imageDigest":"sha256:1a39c97202ef3da294f9775a935d113e7a2ec7c0ab150026507682d014e99990","connectUrl":"https://connect.internal:8083","connectorName":"elite-document-outbox","topicPrefix":"elite-doc","routeTopicReplacement":"outbox.event.${routedByValue}","slotName":"elite_document_outbox","publicationName":"elite_document_outbox","heartbeatIntervalMs":10000})
    for key,value0 in list(value.items()):
        if value0=="REQUIRED":value[key]="evidence://"+key
    return value


class ConfigTests(unittest.TestCase):
    def setUp(self):
        self.lock=subject.load(ROOT/"runtime-lock.json");self.template=subject.load(ROOT/"connector.template.json")

    def test_approved_profile_renders_exact_outbox_connector(self):
        result=subject.render(self.template,approved(),self.lock);config=result["config"]
        self.assertEqual("document_intelligence.outbox_event",config["table.include.list"]);self.assertEqual("fatal",config["transforms.outbox.table.op.invalid.behavior"]);self.assertEqual("elite_document_outbox",config["slot.name"])

    def test_template_is_fail_closed_and_secret_refs_survive(self):
        with self.assertRaises(subject.Rejected):subject.render(self.template,json.loads((ROOT/"project-profile.template.json").read_text()),self.lock)
        result=subject.render(self.template,approved(),self.lock);self.assertTrue(result["config"]["database.password"].startswith("${file:"))

    def test_platform_digest_must_match(self):
        for mutation in ({"platform":"windows/amd64"},{"imageDigest":"sha256:"+"0"*64}):
            value=approved()|mutation
            with self.assertRaises(subject.Rejected):subject.render(self.template,value,self.lock)

    def test_identity_tokens_route_and_heartbeat_are_closed(self):
        for mutation in ({"slotName":"Bad Slot"},{"routeTopicReplacement":"${routedByValue}"},{"heartbeatIntervalMs":0}):
            value=approved()|mutation
            with self.assertRaises(subject.Rejected):subject.render(self.template,value,self.lock)

    def test_connector_invariants_cannot_be_weakened(self):
        for key,bad in (("database.sslmode","disable"),("publication.autocreate.mode","all_tables"),("table.include.list","public.*"),("transforms.outbox.table.op.invalid.behavior","warn"),("slot.drop.on.stop","true")):
            template=copy.deepcopy(self.template);template["config"][key]=bad
            with self.assertRaises(subject.Rejected):subject.render(template,approved(),self.lock)

    def test_lock_identity_and_platforms_are_exact(self):
        for mutation in (("indexDigest","sha256:"+"0"*64),("tagObserved","latest")):
            lock=copy.deepcopy(self.lock);lock["image"][mutation[0]]=mutation[1]
            with self.assertRaises(subject.Rejected):subject.validate_lock(lock)

    def test_duplicate_json_key_is_rejected(self):
        with tempfile.TemporaryDirectory() as directory:
            path=Path(directory)/"bad.json";path.write_text('{"a":1,"a":2}',encoding="utf-8")
            with self.assertRaises(subject.Rejected):subject.load(path)


if __name__=="__main__":unittest.main()
````

### FILE: `debezium_postgres_outbox_runtime/register_connector.ps1`
```yaml
block_id: "DEBEZIUM-POSTGRES-OUTBOX-RUNTIME:register-connector:v1"
operation: CREATE
provenance: ADAPTED
source: "https://debezium.io/documentation/reference/3.6/connectors/postgresql.html"
license: "Apache-2.0 AND LicenseRef-Workspace-Owner"
sha256: "8334f10623a82aef5243d676f6bfee17b6c6eb8e0e976c6605da81e8f73db222"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding(SupportsShouldProcess)]param(
  [Parameter(Mandatory)][string]$Profile,
  [string]$AuthorizationHeaderEnvironmentVariable,
  [switch]$AllowUpdate,
  [string]$ComponentRoot=$PSScriptRoot
)
$ErrorActionPreference='Stop'
$profilePath=(Resolve-Path -LiteralPath $Profile).Path
$profileObject=Get-Content -LiteralPath $profilePath -Raw|ConvertFrom-Json -Depth 64
$temp=Join-Path ([IO.Path]::GetTempPath())('elite-debezium-config-'+[guid]::NewGuid().ToString('N')+'.json')
try{
  python -B (Join-Path $ComponentRoot 'validate_config.py') --component-root $ComponentRoot --profile $profilePath --output $temp
  if($LASTEXITCODE-ne0){throw 'connector configuration validation failed'}
  $request=Get-Content -LiteralPath $temp -Raw|ConvertFrom-Json -Depth 64
  $base=([string]$profileObject.connectUrl).TrimEnd('/')
  if($base-notmatch'^https://'){throw 'Kafka Connect URL must use HTTPS.'}
  $headers=@{}
  if($AuthorizationHeaderEnvironmentVariable){
    if($AuthorizationHeaderEnvironmentVariable-notmatch'^[A-Za-z_][A-Za-z0-9_]{0,127}$'){throw 'Authorization environment variable name is invalid.'}
    $value=[Environment]::GetEnvironmentVariable($AuthorizationHeaderEnvironmentVariable)
    if([string]::IsNullOrWhiteSpace($value)){throw 'Authorization environment variable is empty.'}
    $headers.Authorization=$value
  }
  $name=[Uri]::EscapeDataString([string]$request.name)
  $existing=Invoke-WebRequest -Uri "$base/connectors/$name/config" -Headers $headers -Method Get -SkipHttpErrorCheck
  $jsonOptions=@{Depth=64;Compress=$true}
  if($existing.StatusCode-eq404){
    if(-not$PSCmdlet.ShouldProcess([string]$request.name,'create Debezium connector')){return}
    $body=$request|ConvertTo-Json @jsonOptions
    $response=Invoke-WebRequest -Uri "$base/connectors" -Headers $headers -Method Post -ContentType 'application/json' -Body $body -SkipHttpErrorCheck
    if($response.StatusCode-notin@(200,201)){throw "Connector create failed HTTP $($response.StatusCode)."}
  }elseif($existing.StatusCode-eq200){
    $current=$existing.Content|ConvertFrom-Json -Depth 64
    $expected=$request.config|ConvertTo-Json @jsonOptions
    $actual=$current|ConvertTo-Json @jsonOptions
    if($actual-ne$expected){
      if(-not$AllowUpdate){throw 'Existing connector configuration is divergent; use an approved update.'}
      if(-not$PSCmdlet.ShouldProcess([string]$request.name,'update Debezium connector')){return}
      $response=Invoke-WebRequest -Uri "$base/connectors/$name/config" -Headers $headers -Method Put -ContentType 'application/json' -Body $expected -SkipHttpErrorCheck
      if($response.StatusCode-ne200){throw "Connector update failed HTTP $($response.StatusCode)."}
    }
  }else{throw "Connector lookup failed HTTP $($existing.StatusCode)."}
  $status=Invoke-WebRequest -Uri "$base/connectors/$name/status" -Headers $headers -Method Get -SkipHttpErrorCheck
  if($status.StatusCode-ne200){throw "Connector status failed HTTP $($status.StatusCode)."}
  $state=($status.Content|ConvertFrom-Json -Depth 64).connector.state
  if($state-ne'RUNNING'){throw "Connector is not RUNNING: $state"}
  'DEBEZIUM_POSTGRES_OUTBOX_REGISTRATION_PASS'
}finally{
  if(Test-Path -LiteralPath $temp){[IO.File]::Delete($temp)}
}
````

### FILE: `debezium_postgres_outbox_runtime/verify_contract.ps1`
```yaml
block_id: "DEBEZIUM-POSTGRES-OUTBOX-RUNTIME:verify-contract:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://debezium-postgres-outbox-runtime/verify_contract.ps1"
license: "LicenseRef-Workspace-Owner"
sha256: "171c9640f2766d561864ff90483f4f9e8db911dda2e43ebb17af8a75a9cecddd"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding()]param([string]$ComponentRoot=$PSScriptRoot)
$ErrorActionPreference='Stop'
$files=@('README.md','runtime-lock.json','connector.template.json','project-profile.template.json','consumer-contract.json','postgres_preflight.sql','validate_config.py','test_validate_config.py','register_connector.ps1','verify_contract.ps1')
foreach($file in $files){if(-not(Test-Path -LiteralPath (Join-Path $ComponentRoot $file)-PathType Leaf)){throw "Missing $file"}}
$lock=Get-Content -LiteralPath (Join-Path $ComponentRoot 'runtime-lock.json')-Raw|ConvertFrom-Json -Depth 64
if($lock.image.indexDigest-ne'sha256:76db18f20116557b2e550d5844d2bbebc927cd9c215ec2271c4ab713cc74f386'-or@($lock.image.platforms.PSObject.Properties).Count-ne2-or$lock.source.focusedTests-ne48){throw 'Runtime lock drifted.'}
$profile=Get-Content -LiteralPath (Join-Path $ComponentRoot 'project-profile.template.json')-Raw|ConvertFrom-Json -Depth 64
if($profile.acknowledgeConditionedState-ne$false-or$profile.heartbeatIntervalMs-ne0){throw 'Profile template must begin fail closed.'}
$connector=Get-Content -LiteralPath (Join-Path $ComponentRoot 'connector.template.json')-Raw
foreach($marker in @('document_intelligence.outbox_event','io.debezium.transforms.outbox.EventRouter','publication.autocreate.mode','"disabled"','table.op.invalid.behavior','"fatal"','${file:/run/secrets/elite.properties:db.password}')){if(-not$connector.Contains($marker)){throw "Connector marker missing: $marker"}}
$consumer=Get-Content -LiteralPath (Join-Path $ComponentRoot 'consumer-contract.json')-Raw|ConvertFrom-Json -Depth 64
if($consumer.delivery-ne'AT_LEAST_ONCE'-or$consumer.exactlyOnceClaimAllowed-ne$false-or@($consumer.requiredHeaders).Count-ne4){throw 'Consumer contract drifted.'}
$preflight=Get-Content -LiteralPath (Join-Path $ComponentRoot 'postgres_preflight.sql')-Raw
foreach($marker in @("current_setting('wal_level') = 'logical'",'pg_publication_tables','elite_outbox_cdc','has_table_privilege')){if(-not$preflight.Contains($marker)){throw "Preflight marker missing: $marker"}}
foreach($script in @('register_connector.ps1','verify_contract.ps1')){$tokens=$null;$errors=$null;[void][Management.Automation.Language.Parser]::ParseFile((Join-Path $ComponentRoot $script),[ref]$tokens,[ref]$errors);if(@($errors).Count-ne0){throw "PowerShell parse failed: $script"}}
Push-Location $ComponentRoot
try{
  python -B -c "from pathlib import Path;[compile(Path(p).read_text(encoding='utf-8'),p,'exec') for p in ('validate_config.py','test_validate_config.py')]"
  if($LASTEXITCODE-ne0){throw 'compile failed'}
  python -B -m unittest -v test_validate_config.py
  if($LASTEXITCODE-ne0){throw 'tests failed'}
}finally{Pop-Location}
'DEBEZIUM_POSTGRES_OUTBOX_RUNTIME_PASS'
````

## 6. Configuration surface

Profile inicia cerrado y exige owner/time, plataforma+digest, Connect HTTPS, name/topic/route/slot/publication/heartbeat y evidencias de secrets/TLS/Kafka/logical replication/publication/CDC role/inbox/duplicates/restart-order-load/recovery/rollback/observability/costo. Credenciales DB viven sólo en file provider montado read-only.

## 7. Dependency bill

- Debezium source 3.6.1.Final commit/archive exactos, Apache-2.0, 48 tests focales históricos PASS.
- Quay official connect image: OCI index `76db18…`; linux amd64 `1a39c9…`; arm64 `b98817…`.
- Kafka Connect/broker y PostgreSQL 15+ con logical replication son target y no se incluyen.
- Ocho archivos `AUTHORED`, dos `ADAPTED`, cero `VERBATIM`.

## 8. Apply order

1. Aplicar boundary 0.2 migration; crear publicación exacta y membership CDC read-only.
2. Ejecutar preflight y validar plataforma/digest.
3. Ejecutar runtime por digest con file ConfigProvider, secrets/CA read-only y Kafka TLS/auth.
4. Completar/aprobar profile; correr validator.
5. Registrar por HTTPS; divergencia exige update explícito.
6. Probar snapshot, duplicate/replay, restart/slot, ordering/load, broker outage/recovery, schema, observabilidad y rollback.
7. Admitir consumers sólo con inbox/payload-hash+business effect atómico.

## 9. Verification

El verifier exige 10 archivos, lock exacto, templates cerrados, invariantes EventRouter/secret/TLS/publication/table, preflight, parse PowerShell, compile y 7 tests. Pull/launch/registration/PostgreSQL/Kafka live permanecen condicionados.

## 10. Reconstruction evidence

V96 debe fijar hashes 10/10, boundary 0.2 Debezium-shaped 11/11, round-trips, 12+7 tests, cfn-lint y gates globales. No se afirma live ni exactly-once.
