# Debezium PostgreSQL Inbox Consumer

## 1. Metadata

```yaml
pack_id: "DEBEZIUM-POSTGRES-INBOX-CONSUMER"
pack_version: "0.2.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa el consumer Red Hat/Debezium outbox adaptado: valida headers/payload, ejecuta un mapping CEL-Java aprobado y aplica inbox más proyección mapeada en una transacción PostgreSQL antes del ack Kafka."
stacks: ["Red Hat Debezium examples 7b0d765", "Google-origin CEL-Java 0.13.0", "Quarkus 3.39.1", "Jackson 2.22.1", "Java 21", "Kafka", "PostgreSQL 18.6", "OSV Scanner 2.5.1", "PowerShell 7"]
compatible_with: ["DEBEZIUM-POSTGRES-OUTBOX-RUNTIME 0.1.x", "AWS-IDP-POSTGRES-PERSISTENCE-BOUNDARY 0.2.x", "PG-TX-FOUNDATION 0.1.x", "GOOGLE-CEL-DOCUMENT-MAPPING 0.1.x contract"]
incompatible_with: ["Quarkus 3.33.2", "Jackson 2.22.0", "multiline CEL", "CEL macros", "mapping/schema sin hash", "auto commit", "ack antes de commit", "consumer sin inbox", "payload abierto", "evento desconocido ignorado", "log de payload", "exactly-once no demostrado"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0"
upstream_sources: ["https://github.com/debezium/debezium-examples/tree/7b0d765a02cf66ef29de2a01b67ec763e9acedf5/outbox", "https://github.com/cel-expr/cel-java/tree/c57d9d02259795c761e82f8870303b93a7a672c8", "https://repo.maven.apache.org/maven2/dev/cel/cel/0.13.0/", "https://repo.maven.apache.org/maven2/io/quarkus/platform/quarkus-bom/3.39.1/", "https://repo.maven.apache.org/maven2/com/fasterxml/jackson/core/jackson-databind/2.22.1/", "https://debezium.io/documentation/reference/3.6/transformations/outbox-event-router.html"]
verified_at: "2026-08-28"
```

## 2. Applicability

Use después del boundary 0.2 y runtime Debezium 0.1 cuando el proyecto seleccione Kafka/PostgreSQL. Completa deduplicación, mapping configurable aprobado y proyección documental concreta; no inventa mappings ERP/CRM ni autoriza infraestructura live.

## 3. Architecture contract

Conserva seis archivos oficiales Red Hat como autoridad de patrón: cinco Java son byte-exactos y la licencia está adaptada únicamente por el LF final que falta en el original. La adaptación importa el artefacto oficial Google-origin CEL-Java 0.13.0 sin copiar source: expresión single-line, macros vacías, límites 4.096/32, clase documental exacta, mapping/schema ID+version+SHA y output exacto. Tras validar topic/key/headers/payload, abre una sola transacción, inserta inbox, exige que `domainType` coincida con la clase del mapping, lee el `document_payload` revisado bajo RLS, ejecuta CEL y persiste payload mapeado más receipt; sólo después del commit ocurre ACK. Redelivery recompone/verifica la misma autoridad; evento cross-class, mapping divergente, input ausente o fallo hace rollback sin ACK.

## 4. Exact file manifest

```text
CREATE debezium_postgres_inbox_consumer/pom.xml
CREATE debezium_postgres_inbox_consumer/project-profile.template.json
CREATE debezium_postgres_inbox_consumer/README.md
CREATE debezium_postgres_inbox_consumer/run_postgres_tests.ps1
CREATE debezium_postgres_inbox_consumer/run_sca.ps1
CREATE debezium_postgres_inbox_consumer/sca-receipt.json
CREATE debezium_postgres_inbox_consumer/source-lock.json
CREATE debezium_postgres_inbox_consumer/src/main/java/com/elite/documentinbox/CelDocumentProjectionMapping.java
CREATE debezium_postgres_inbox_consumer/src/main/java/com/elite/documentinbox/DocumentIntakeEnvelope.java
CREATE debezium_postgres_inbox_consumer/src/main/java/com/elite/documentinbox/DocumentIntakeEnvelopeValidator.java
CREATE debezium_postgres_inbox_consumer/src/main/java/com/elite/documentinbox/JdbcDocumentIntakeProcessor.java
CREATE debezium_postgres_inbox_consumer/src/main/java/com/elite/documentinbox/KafkaDocumentIntakeConsumer.java
CREATE debezium_postgres_inbox_consumer/src/main/resources/application.properties
CREATE debezium_postgres_inbox_consumer/src/main/resources/db/migration/V001__document_consumer_inbox.sql
CREATE debezium_postgres_inbox_consumer/src/test/java/com/elite/documentinbox/CelDocumentProjectionMappingTest.java
CREATE debezium_postgres_inbox_consumer/src/test/java/com/elite/documentinbox/DocumentIntakeEnvelopeValidatorTest.java
CREATE debezium_postgres_inbox_consumer/src/test/java/com/elite/documentinbox/JdbcDocumentIntakeProcessorTest.java
CREATE debezium_postgres_inbox_consumer/upstream/ConsumedMessage.java
CREATE debezium_postgres_inbox_consumer/upstream/KafkaEventConsumer.java
CREATE debezium_postgres_inbox_consumer/upstream/LICENSE.txt
CREATE debezium_postgres_inbox_consumer/upstream/MessageLog.java
CREATE debezium_postgres_inbox_consumer/upstream/OrderEventHandler.java
CREATE debezium_postgres_inbox_consumer/upstream/ShipmentService.java
CREATE debezium_postgres_inbox_consumer/verify_contract.ps1
```

## 5. Materialization blocks

### FILE: `debezium_postgres_inbox_consumer/pom.xml`
```yaml
block_id: "DEBEZIUM-POSTGRES-INBOX-CONSUMER:pom-xml:v1"
operation: CREATE
provenance: ADAPTED
source: "https://github.com/debezium/debezium-examples/tree/7b0d765a02cf66ef29de2a01b67ec763e9acedf5/outbox"
license: "LicenseRef-Workspace-Owner"
sha256: "b86f2f7ee758727535016e7fd6eea04319ea18f3557cd8306eb9f59e870b0b63"
variables: []
secrets_allowed: false
```
````xml
<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 https://maven.apache.org/xsd/maven-4.0.0.xsd">
  <modelVersion>4.0.0</modelVersion>
  <groupId>com.elite</groupId>
  <artifactId>debezium-postgres-inbox-consumer</artifactId>
  <version>0.2.0</version>
  <properties>
    <maven.compiler.release>21</maven.compiler.release>
    <project.build.sourceEncoding>UTF-8</project.build.sourceEncoding>
    <quarkus.platform.version>3.39.1</quarkus.platform.version>
    <jackson-databind.version>2.22.1</jackson-databind.version>
    <cel.version>0.13.0</cel.version>
    <compiler-plugin.version>3.15.0</compiler-plugin.version>
    <surefire-plugin.version>3.5.5</surefire-plugin.version>
  </properties>
  <dependencyManagement>
    <dependencies>
      <dependency>
        <groupId>io.quarkus.platform</groupId><artifactId>quarkus-bom</artifactId><version>${quarkus.platform.version}</version><type>pom</type><scope>import</scope>
      </dependency>
    </dependencies>
  </dependencyManagement>
  <dependencies>
    <dependency><groupId>io.quarkus</groupId><artifactId>quarkus-arc</artifactId></dependency>
    <dependency><groupId>io.quarkus</groupId><artifactId>quarkus-jackson</artifactId></dependency>
    <dependency><groupId>com.fasterxml.jackson.core</groupId><artifactId>jackson-databind</artifactId><version>${jackson-databind.version}</version></dependency>
    <dependency><groupId>io.quarkus</groupId><artifactId>quarkus-jdbc-postgresql</artifactId></dependency>
    <dependency><groupId>io.quarkus</groupId><artifactId>quarkus-messaging-kafka</artifactId></dependency>
    <dependency><groupId>io.quarkus</groupId><artifactId>quarkus-smallrye-health</artifactId></dependency>
    <dependency><groupId>dev.cel</groupId><artifactId>cel</artifactId><version>${cel.version}</version></dependency>
    <dependency><groupId>org.junit.jupiter</groupId><artifactId>junit-jupiter</artifactId><scope>test</scope></dependency>
  </dependencies>
  <build>
    <plugins>
      <plugin><groupId>io.quarkus</groupId><artifactId>quarkus-maven-plugin</artifactId><version>${quarkus.platform.version}</version><extensions>true</extensions><executions><execution><goals><goal>build</goal><goal>generate-code</goal><goal>generate-code-tests</goal></goals></execution></executions></plugin>
      <plugin><groupId>org.apache.maven.plugins</groupId><artifactId>maven-compiler-plugin</artifactId><version>${compiler-plugin.version}</version></plugin>
      <plugin><groupId>org.apache.maven.plugins</groupId><artifactId>maven-surefire-plugin</artifactId><version>${surefire-plugin.version}</version><configuration><useModulePath>false</useModulePath><systemPropertyVariables><java.util.logging.manager>org.jboss.logmanager.LogManager</java.util.logging.manager></systemPropertyVariables></configuration></plugin>
    </plugins>
  </build>
</project>
````

### FILE: `debezium_postgres_inbox_consumer/project-profile.template.json`
```yaml
block_id: "DEBEZIUM-POSTGRES-INBOX-CONSUMER:project-profile-template-json:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://debezium-postgres-inbox-consumer/project-profile.template.json"
license: "LicenseRef-Workspace-Owner"
sha256: "c37930a730631f997df56cd5b6da80b22a34ca4893ef66d102e297abf0c27445"
variables: []
secrets_allowed: false
```
````json
{
  "schemaVersion": 1,
  "acknowledgeConditionedState": false,
  "topic": "",
  "consumerGroup": "",
  "bootstrapServers": "",
  "databaseJdbcUrlSecretRef": "",
  "databaseUsernameSecretRef": "",
  "databasePasswordSecretRef": "",
  "tlsTruststoreSecretRef": "",
  "tlsKeystoreSecretRef": "",
  "schemaRegistryUrl": "",
  "mappingId": "",
  "mappingDomainType": "",
  "mappingVersion": "",
  "mappingExpression": "",
  "mappingExpressionSha256": "",
  "mappingExactOutputKeysCsv": "",
  "mappingInputSchemaId": "",
  "mappingInputSchemaSha256": "",
  "mappingOutputSchemaId": "",
  "mappingOutputSchemaSha256": "",
  "owners": [],
  "approvals": [],
  "proofs": {
    "migrationAppliedByOwner": false,
    "runtimeRoleIsNonOwner": false,
    "tlsAndAuthenticationProven": false,
    "duplicateRaceRollbackProven": false,
    "restartRedeliveryOrderingProven": false,
    "loadAndRecoveryProven": false,
    "observabilityRedactionProven": false,
    "mappingSchemasAndCorpusApproved": false,
    "mappingDeterminismAndLimitsProven": false,
    "rollbackProven": false
  }
}
````

### FILE: `debezium_postgres_inbox_consumer/README.md`
```yaml
block_id: "DEBEZIUM-POSTGRES-INBOX-CONSUMER:readme-md:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://debezium-postgres-inbox-consumer/README.md"
license: "LicenseRef-Workspace-Owner"
sha256: "4a27591ea8ad1fd8acc002748e4a8ada3f7aa23fab67b835d1326b58e7efb0d8"
variables: []
secrets_allowed: false
```
````markdown
# Debezium PostgreSQL inbox consumer

This component turns the official Debezium outbox consumer pattern into a fail-closed document-intake consumer. The exact Red Hat/Debezium demo sources are preserved under `upstream/`; production code is an explicit local adaptation and is never attributed to Red Hat.

The consumer requires the Debezium Event Router `id`, `tenant_id`, `schema_version` and `type` headers. It validates one header value each, the exact topic, Kafka key, JSON schema and hashes before opening the database transaction. Inside that same transaction it reads the authoritative reviewed `document_payload`, requires the event `domainType` to equal the mapping's configured class, executes one approved CEL-Java 0.13.0 mapping with macros disabled, binds mapping/schema identities and hashes, and inserts inbox plus mapped projection. A duplicate recomputes and verifies the exact mapping receipt without repeating the effect; a cross-class event, divergent replay, mapping/schema mismatch, unknown event, missing intake row or projection conflict rolls back and is not acknowledged.

Run `verify_contract.ps1` for static gates. Run `run_postgres_tests.ps1` with exact Java 21, Maven 3.9.16 and PostgreSQL 18.6 paths for compilation and 18 real tests. Run `run_sca.ps1` to regenerate a CycloneDX 1.6 SBOM and require OSV 2.5.1 to report zero known findings. The adapted runtime uses official Quarkus 3.39.1, Jackson 2.22.1 and Google-origin CEL-Java 0.13.0 artifacts; no CEL source is embedded or attributed to Red Hat/Google. Deploy one approved mapping/domain class per consumer configuration; Kafka/PostgreSQL live deployment, schemas/corpus/approval, TLS/auth, schema registry, poison policy, restart/redelivery/order/load, backup/restore and observability remain project gates.
````

### FILE: `debezium_postgres_inbox_consumer/run_postgres_tests.ps1`
```yaml
block_id: "DEBEZIUM-POSTGRES-INBOX-CONSUMER:run-postgres-tests-ps1:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://debezium-postgres-inbox-consumer/run_postgres_tests.ps1"
license: "LicenseRef-Workspace-Owner"
sha256: "4411659a3e536e7b3006f5c8089ed15bca64fb4211f11c2e47b7731e798c0c7b"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding()]
param(
  [Parameter(Mandatory)][string]$JavaHome,
  [Parameter(Mandatory)][string]$MavenCommand,
  [Parameter(Mandatory)][string]$PostgresBin
)
$ErrorActionPreference = 'Stop'
$java = Join-Path $JavaHome 'bin/java.exe'
$initdb = Join-Path $PostgresBin 'initdb.exe'
$pgCtl = Join-Path $PostgresBin 'pg_ctl.exe'
$pgConfig = Join-Path $PostgresBin 'pg_config.exe'
foreach ($path in @($java,$MavenCommand,$initdb,$pgCtl,$pgConfig)) { if (-not (Test-Path -LiteralPath $path -PathType Leaf)) { throw "Required tool missing: $path" } }
$javaVersion = (& $java -version 2>&1 | Out-String)
if ($javaVersion -notmatch '21\.0\.8') { throw 'Exact Java 21.0.8 is required by this dated gate' }
$pgVersion = (& $pgConfig --version | Out-String).Trim()
if ($pgVersion -ne 'PostgreSQL 18.6') { throw "Exact PostgreSQL 18.6 required; got $pgVersion" }
$oldJavaHome = $env:JAVA_HOME
$oldTestUrl = $env:ELITE_TEST_PG_URL
$testRoot = Join-Path $env:TEMP ('elite-inbox-pg-' + [guid]::NewGuid().ToString('N'))
$data = Join-Path $testRoot 'data'
$log = Join-Path $testRoot 'postgres.log'
$started = $false
Push-Location -LiteralPath $PSScriptRoot
try {
  New-Item -ItemType Directory -Path $testRoot | Out-Null
  & $initdb -D $data --no-locale --encoding=UTF8 --auth=trust --username=postgres | Out-Null
  if ($LASTEXITCODE -ne 0) { throw 'initdb failed' }
  $listener = [System.Net.Sockets.TcpListener]::new([System.Net.IPAddress]::Loopback, 0)
  $listener.Start(); $port = ([System.Net.IPEndPoint]$listener.LocalEndpoint).Port; $listener.Stop()
  $startArgs = @('-D',('"' + $data + '"'),'-l',('"' + $log + '"'),'-o',('"-p ' + $port + ' -h 127.0.0.1"'),'-w','start')
  $startProcess = Start-Process -FilePath $pgCtl -ArgumentList $startArgs -WindowStyle Hidden -PassThru
  $startProcess.WaitForExit()
  if ($startProcess.ExitCode -ne 0) { throw 'PostgreSQL start failed' }
  $started = $true
  $env:JAVA_HOME = $JavaHome
  $env:ELITE_TEST_PG_URL = "jdbc:postgresql://127.0.0.1:$port/postgres"
  & $MavenCommand -B -ntp clean test
  if ($LASTEXITCODE -ne 0) { throw 'Maven test failed' }
  Write-Output 'DEBEZIUM_POSTGRES_INBOX_CONSUMER_POSTGRES_PASS tests=18 concurrency=32 cross_class_rejected=1 mapping_transactional=1 rollback=1'
}
finally {
  if ($started) {
    $stopArgs = @('-D',('"' + $data + '"'),'-m','fast','-w','stop')
    $stopProcess = Start-Process -FilePath $pgCtl -ArgumentList $stopArgs -WindowStyle Hidden -PassThru
    $stopProcess.WaitForExit()
    if ($stopProcess.ExitCode -ne 0) { Write-Warning 'PostgreSQL stop returned non-zero' }
  }
  $env:JAVA_HOME = $oldJavaHome
  $env:ELITE_TEST_PG_URL = $oldTestUrl
  Pop-Location
  $resolved = [System.IO.Path]::GetFullPath($testRoot)
  $tempResolved = [System.IO.Path]::GetFullPath($env:TEMP).TrimEnd('\') + '\'
  if ($resolved.StartsWith($tempResolved, [System.StringComparison]::OrdinalIgnoreCase) -and (Test-Path -LiteralPath $resolved)) {
    [System.IO.Directory]::Delete($resolved, $true)
  }
}
````

### FILE: `debezium_postgres_inbox_consumer/run_sca.ps1`
```yaml
block_id: "DEBEZIUM-POSTGRES-INBOX-CONSUMER:run-sca-ps1:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://debezium-postgres-inbox-consumer/run_sca.ps1"
license: "LicenseRef-Workspace-Owner"
sha256: "41eb917976c219fa5575887ba696381ece3f7cd4287cb05e5a1a7404ea197ae2"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding()]
param(
  [Parameter(Mandatory)][string]$JavaHome,
  [Parameter(Mandatory)][string]$MavenCommand,
  [Parameter(Mandatory)][string]$OsvScanner
)
$ErrorActionPreference = 'Stop'
$java = Join-Path $JavaHome 'bin/java.exe'
foreach ($path in @($java,$MavenCommand,$OsvScanner)) { if (-not (Test-Path -LiteralPath $path -PathType Leaf)) { throw "Required tool missing: $path" } }
if ((& $java -version 2>&1 | Out-String) -notmatch '21\.0\.8') { throw 'Exact Java 21.0.8 required' }
if ((& $OsvScanner --version | Out-String) -notmatch '2\.5\.1') { throw 'Exact OSV Scanner 2.5.1 required' }
$oldJavaHome = $env:JAVA_HOME
$report = Join-Path $env:TEMP ('elite-inbox-osv-' + [guid]::NewGuid().ToString('N') + '.json')
Push-Location -LiteralPath $PSScriptRoot
try {
  $env:JAVA_HOME = $JavaHome
  & $MavenCommand -B -ntp org.cyclonedx:cyclonedx-maven-plugin:2.9.1:makeAggregateBom -DskipTests
  if ($LASTEXITCODE -ne 0) { throw 'CycloneDX generation failed' }
  & $OsvScanner scan -L (Join-Path $PSScriptRoot 'target/bom.json') --format json --all-packages --output-file $report
  $scanExit = $LASTEXITCODE
  $scan = Get-Content -LiteralPath $report -Raw | ConvertFrom-Json -Depth 100
  $packages = @($scan.results.packages)
  $findings = @($packages | ForEach-Object { @($_.vulnerabilities) } | Where-Object { $null -ne $_ })
  if ($scanExit -ne 0 -or $packages.Count -ne 170 -or $findings.Count -ne 0) { throw "SCA failed: exit=$scanExit components=$($packages.Count) findings=$($findings.Count)" }
  $bom = Get-Content -LiteralPath (Join-Path $PSScriptRoot 'target/bom.json') -Raw
  if (-not $bom.Contains('pkg:maven/io.quarkus/quarkus-core@3.39.1?type=jar') -or -not $bom.Contains('pkg:maven/com.fasterxml.jackson.core/jackson-databind@2.22.1?type=jar') -or -not $bom.Contains('pkg:maven/dev.cel/cel@0.13.0?type=jar')) { throw 'SBOM runtime identity mismatch' }
  Write-Output 'DEBEZIUM_POSTGRES_INBOX_CONSUMER_SCA_PASS components=170 findings=0 cel=0.13.0'
}
finally {
  $env:JAVA_HOME = $oldJavaHome
  Pop-Location
  if (Test-Path -LiteralPath $report) { [System.IO.File]::Delete($report) }
}
````

### FILE: `debezium_postgres_inbox_consumer/sca-receipt.json`
```yaml
block_id: "DEBEZIUM-POSTGRES-INBOX-CONSUMER:sca-receipt-json:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://debezium-postgres-inbox-consumer/sca-receipt.json"
license: "LicenseRef-Workspace-Owner"
sha256: "0468dcc91a59662c8a246b73ff05816cb118c16167e0ea35c8bc73e08ff76331"
variables: []
secrets_allowed: false
```
````json
{
  "schemaVersion": 1,
  "scannedAt": "2026-08-28",
  "java": "21.0.8",
  "maven": "3.9.16",
  "cycloneDxPlugin": "2.9.1",
  "cycloneDxSpec": "1.6",
  "components": 170,
  "quarkusBom": "3.39.1",
  "jacksonDatabind": "2.22.1",
  "celJava": "0.13.0",
  "osvScanner": "2.5.1",
  "osvCommit": "c84fa4568f2526d0333e9a914ea8a0a5f74ad68b",
  "findings": 0,
  "uniqueIds": 0,
  "bomSha256": "21b75db4fc6070015b78dae77af3007e0aed8eda022dc609c1ef7c274406e503",
  "reportSha256": "e9d653a5c4f1701ccb10ef56a5d473ed59b0e66b5dfaa7127bc7a6b7c7beea37"
}
````

### FILE: `debezium_postgres_inbox_consumer/source-lock.json`
```yaml
block_id: "DEBEZIUM-POSTGRES-INBOX-CONSUMER:source-lock-json:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://debezium-postgres-inbox-consumer/source-lock.json"
license: "LicenseRef-Workspace-Owner"
sha256: "07ee81cea47838afea2a2b7ca367f6de49a61927533205286a5f89aa2eb1d1e6"
variables: []
secrets_allowed: false
```
````json
{
  "schemaVersion": 1,
  "repository": "debezium/debezium-examples",
  "commit": "7b0d765a02cf66ef29de2a01b67ec763e9acedf5",
  "commitDate": "2026-08-25T13:59:32Z",
  "commitVerified": false,
  "commitVerificationReason": "unsigned",
  "archiveUrl": "https://github.com/debezium/debezium-examples/archive/7b0d765a02cf66ef29de2a01b67ec763e9acedf5.zip",
  "archiveBytes": 9292212,
  "archiveSha256": "0fba61c33bdbdb8f32eeab790335dcec412826dd98fa0bbab282afa2130f51e4",
  "license": "Apache-2.0",
  "officialBuild": {
    "java": "Microsoft OpenJDK 21.0.8",
    "maven": "3.9.16",
    "quarkus": "3.33.2",
    "compiledSources": 6,
    "testsRun": 0,
    "result": "BUILD_SUCCESS_CONDITIONED"
  },
  "adaptationRuntime": {
    "quarkusBom": "3.39.1",
    "jacksonDatabindOverride": "2.22.1",
    "celJava": "dev.cel:cel:0.13.0",
    "celJavaJarSha256": "a3bc22d78affe749f81b17cce2131a3807275436157ba4d6fe3f10d8f541ea94",
    "cycloneDxPlugin": "2.9.1",
    "components": 170,
    "osvScanner": "2.5.1",
    "osvFindings": 0,
    "scannedAt": "2026-08-28"
  },
  "upstreamFiles": [
    {"path":"upstream/LICENSE.txt","sourcePath":"LICENSE.txt","sourceBytes":11357,"sourceSha256":"58d1e17ffe5109a7ae296caafcadfdbe6a7d176f0bc4ab01e12a689b0499d8bd","packagedBytes":11358,"sha256":"cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30","change":"final LF only"},
    {"path":"upstream/KafkaEventConsumer.java","sourcePath":"outbox/shipment-service/src/main/java/io/debezium/examples/outbox/shipment/facade/KafkaEventConsumer.java","sha256":"783184dfe7d894a7f17676e78a53daa00c99de9f2fd9b2cabac86922b1d813a3"},
    {"path":"upstream/OrderEventHandler.java","sourcePath":"outbox/shipment-service/src/main/java/io/debezium/examples/outbox/shipment/facade/OrderEventHandler.java","sha256":"ccd1de2cb73a6630a08355fddcbba979bea7b8f3650e6ddb16d2590f0b1198ab"},
    {"path":"upstream/ConsumedMessage.java","sourcePath":"outbox/shipment-service/src/main/java/io/debezium/examples/outbox/shipment/messagelog/ConsumedMessage.java","sha256":"aa78b54720fe26268141ea79720c3b1d691c6d57319d40eaacbd80608c9c72e1"},
    {"path":"upstream/MessageLog.java","sourcePath":"outbox/shipment-service/src/main/java/io/debezium/examples/outbox/shipment/messagelog/MessageLog.java","sha256":"b60bc4cf5fd7a4dd0d783dad04e2595fa039f862af3e1f6dede931af9b74e55b"},
    {"path":"upstream/ShipmentService.java","sourcePath":"outbox/shipment-service/src/main/java/io/debezium/examples/outbox/shipment/service/ShipmentService.java","sha256":"47e8fb65169a3a562b4572b27e3903d7c6f2f3e81dd42cf546d816126c896cf2"}
  ]
}
````

### FILE: `debezium_postgres_inbox_consumer/src/main/java/com/elite/documentinbox/CelDocumentProjectionMapping.java`
```yaml
block_id: "DEBEZIUM-POSTGRES-INBOX-CONSUMER:cel-document-projection-mapping-java:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://debezium-postgres-inbox-consumer/src/main/java/com/elite/documentinbox/CelDocumentProjectionMapping.java"
license: "LicenseRef-Workspace-Owner"
sha256: "ec2f6a2c27d5cc1383553139daaf934e281e1733cf73e2f956910bbe6c60e356"
variables: []
secrets_allowed: false
```
````java
package com.elite.documentinbox;

import com.fasterxml.jackson.core.JsonFactory;
import com.fasterxml.jackson.core.StreamReadFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;
import dev.cel.common.CelAbstractSyntaxTree;
import dev.cel.common.CelOptions;
import dev.cel.common.CelValidationException;
import dev.cel.common.types.MapType;
import dev.cel.common.types.SimpleType;
import dev.cel.compiler.CelCompilerFactory;
import dev.cel.runtime.CelEvaluationException;
import dev.cel.runtime.CelRuntime;
import dev.cel.runtime.CelRuntimeFactory;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import org.eclipse.microprofile.config.inject.ConfigProperty;

import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.util.ArrayList;
import java.util.HexFormat;
import java.util.IdentityHashMap;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.TreeMap;
import java.util.regex.Pattern;

@ApplicationScoped
public final class CelDocumentProjectionMapping {
    private static final CelOptions OPTIONS = CelOptions.current()
            .maxExpressionCodePointSize(4096).maxParseRecursionDepth(32).build();
    private static final MapType STRING_DYN_MAP = MapType.create(SimpleType.STRING, SimpleType.DYN);
    private static final Pattern ID = Pattern.compile("^[a-z][a-z0-9-]{2,63}$");
    private static final Pattern DOMAIN = Pattern.compile("^[a-z][a-z0-9._-]{1,126}[a-z0-9]$");
    private static final Pattern VERSION = Pattern.compile("^[1-9][0-9]*\\.[0-9]+\\.[0-9]+$");
    private static final Pattern SHA256 = Pattern.compile("^[0-9a-f]{64}$");
    private static final ObjectMapper JSON = new ObjectMapper(JsonFactory.builder()
            .enable(StreamReadFeature.STRICT_DUPLICATE_DETECTION).build())
            .enable(SerializationFeature.ORDER_MAP_ENTRIES_BY_KEYS);

    private final CelRuntime.Program program;
    private final String mappingId;
    private final String mappingDomainType;
    private final String mappingVersion;
    private final String expressionSha256;
    private final Set<String> exactOutputKeys;
    private final String inputSchemaId;
    private final String inputSchemaSha256;
    private final String outputSchemaId;
    private final String outputSchemaSha256;

    @Inject
    public CelDocumentProjectionMapping(
            @ConfigProperty(name = "elite.document-inbox.mapping-id") String mappingId,
            @ConfigProperty(name = "elite.document-inbox.mapping-domain-type") String mappingDomainType,
            @ConfigProperty(name = "elite.document-inbox.mapping-version") String mappingVersion,
            @ConfigProperty(name = "elite.document-inbox.mapping-expression") String expression,
            @ConfigProperty(name = "elite.document-inbox.mapping-expression-sha256") String expressionSha256,
            @ConfigProperty(name = "elite.document-inbox.mapping-output-keys") String exactOutputKeysCsv,
            @ConfigProperty(name = "elite.document-inbox.mapping-input-schema-id") String inputSchemaId,
            @ConfigProperty(name = "elite.document-inbox.mapping-input-schema-sha256") String inputSchemaSha256,
            @ConfigProperty(name = "elite.document-inbox.mapping-output-schema-id") String outputSchemaId,
            @ConfigProperty(name = "elite.document-inbox.mapping-output-schema-sha256") String outputSchemaSha256)
            throws CelValidationException, CelEvaluationException {
        if (!ID.matcher(required(mappingId)).matches() || !VERSION.matcher(required(mappingVersion)).matches()) {
            throw new IllegalArgumentException("Canonical mapping identity is required");
        }
        if (!DOMAIN.matcher(required(mappingDomainType)).matches()) {
            throw new IllegalArgumentException("Canonical mapping domain type is required");
        }
        if (expression.indexOf('\r') >= 0 || expression.indexOf('\n') >= 0) {
            throw new IllegalArgumentException("Mapping expression must be single-line");
        }
        String actualExpressionSha256 = sha256(required(expression));
        if (!actualExpressionSha256.equals(expressionSha256)) {
            throw new IllegalArgumentException("Mapping expression hash mismatch");
        }
        this.exactOutputKeys = parseKeys(exactOutputKeysCsv);
        this.inputSchemaId = schemaId(inputSchemaId);
        this.outputSchemaId = schemaId(outputSchemaId);
        this.inputSchemaSha256 = schemaHash(inputSchemaSha256);
        this.outputSchemaSha256 = schemaHash(outputSchemaSha256);
        this.mappingId = mappingId;
        this.mappingDomainType = mappingDomainType;
        this.mappingVersion = mappingVersion;
        this.expressionSha256 = actualExpressionSha256;
        var compiler = CelCompilerFactory.standardCelCompilerBuilder()
                .setOptions(OPTIONS).setStandardMacros()
                .addVar("document", STRING_DYN_MAP).setResultType(STRING_DYN_MAP).build();
        CelAbstractSyntaxTree ast = compiler.compile(expression).getAst();
        CelRuntime runtime = CelRuntimeFactory.standardCelRuntimeBuilder().setOptions(OPTIONS).build();
        this.program = runtime.createProgram(ast);
    }

    public MappedDocument apply(String domainType, String documentJson, String sourcePayloadSha256) {
        if (!mappingDomainType.equals(required(domainType))) {
            throw new IllegalArgumentException("Event domain type does not match approved mapping");
        }
        if (!SHA256.matcher(required(sourcePayloadSha256)).matches()) {
            throw new IllegalArgumentException("Source payload hash is invalid");
        }
        try {
            Object rawInput = JSON.readValue(required(documentJson), Object.class);
            if (!(rawInput instanceof Map<?, ?> inputMap)) {
                throw new IllegalArgumentException("Document payload must be an object");
            }
            validateJsonValue(inputMap, 0, new IdentityHashMap<>());
            Object evaluated = program.eval(Map.of("document", inputMap));
            if (!(evaluated instanceof Map<?, ?> rawOutput)) {
                throw new IllegalStateException("Mapping result is not an object");
            }
            LinkedHashMap<String, Object> output = new LinkedHashMap<>();
            for (Map.Entry<?, ?> entry : rawOutput.entrySet()) {
                if (!(entry.getKey() instanceof String key) || entry.getValue() == null) {
                    throw new IllegalStateException("Mapping result contains invalid key or null value");
                }
                validateJsonValue(entry.getValue(), 0, new IdentityHashMap<>());
                output.put(key, entry.getValue());
            }
            if (!output.keySet().equals(exactOutputKeys)) {
                throw new IllegalStateException("Mapping result keys do not match approved schema");
            }
            String mappedJson = JSON.writeValueAsString(new TreeMap<>(output));
            return new MappedDocument(mappingId, mappingDomainType, mappingVersion, expressionSha256,
                    inputSchemaId, inputSchemaSha256, outputSchemaId, outputSchemaSha256,
                    sourcePayloadSha256, sha256(mappedJson), mappedJson);
        }
        catch (CelEvaluationException e) {
            throw new IllegalArgumentException("Approved mapping evaluation failed", e);
        }
        catch (IllegalArgumentException | IllegalStateException e) {
            throw e;
        }
        catch (Exception e) {
            throw new IllegalArgumentException("Document mapping failed", e);
        }
    }

    private static Set<String> parseKeys(String csv) {
        String value = required(csv);
        if (!value.matches("[a-z][a-z0-9_]{0,62}(,[a-z][a-z0-9_]{0,62})*")) {
            throw new IllegalArgumentException("Exact output keys CSV is not canonical");
        }
        List<String> keys = List.of(value.split(",", -1));
        if (Set.copyOf(keys).size() != keys.size()) {
            throw new IllegalArgumentException("Exact output keys contain duplicates");
        }
        return Set.copyOf(keys);
    }

    private static String schemaId(String value) {
        if (!required(value).matches("^[a-z][a-z0-9._-]{2,127}$")) {
            throw new IllegalArgumentException("Canonical schema id is required");
        }
        return value;
    }

    private static String schemaHash(String value) {
        if (!SHA256.matcher(required(value)).matches()) {
            throw new IllegalArgumentException("Schema SHA-256 is required");
        }
        return value;
    }

    private static String required(String value) {
        if (value == null || value.isBlank() || !value.equals(value.trim())) {
            throw new IllegalArgumentException("Required value is missing or not canonical");
        }
        return value;
    }

    private static String sha256(String value) {
        try {
            return HexFormat.of().formatHex(MessageDigest.getInstance("SHA-256")
                    .digest(value.getBytes(StandardCharsets.UTF_8)));
        }
        catch (Exception impossible) {
            throw new IllegalStateException("SHA-256 unavailable", impossible);
        }
    }

    private static void validateJsonValue(Object value, int depth, IdentityHashMap<Object, Boolean> active) {
        if (depth > 16) throw new IllegalArgumentException("Document value exceeds maximum depth");
        if (value instanceof String text) {
            if (text.length() > 1_048_576) throw new IllegalArgumentException("Document string exceeds maximum length");
            return;
        }
        if (value instanceof Boolean || value instanceof Byte || value instanceof Short
                || value instanceof Integer || value instanceof Long) return;
        if (value instanceof Map<?, ?> map) {
            if (map.size() > 10_000 || active.put(map, Boolean.TRUE) != null) {
                throw new IllegalArgumentException("Document map is oversized or cyclic");
            }
            try {
                for (Map.Entry<?, ?> entry : map.entrySet()) {
                    if (!(entry.getKey() instanceof String key) || key.isEmpty() || key.length() > 256
                            || entry.getValue() == null) {
                        throw new IllegalArgumentException("Document map key/value is invalid");
                    }
                    validateJsonValue(entry.getValue(), depth + 1, active);
                }
            }
            finally { active.remove(map); }
            return;
        }
        if (value instanceof List<?> list) {
            if (list.size() > 10_000 || active.put(list, Boolean.TRUE) != null) {
                throw new IllegalArgumentException("Document list is oversized or cyclic");
            }
            try {
                for (Object item : new ArrayList<>(list)) {
                    if (item == null) throw new IllegalArgumentException("Document list contains null");
                    validateJsonValue(item, depth + 1, active);
                }
            }
            finally { active.remove(list); }
            return;
        }
        throw new IllegalArgumentException("Document contains a non-JSON or imprecise numeric value");
    }

    public record MappedDocument(String mappingId, String mappingDomainType, String mappingVersion, String expressionSha256,
                                 String inputSchemaId, String inputSchemaSha256,
                                 String outputSchemaId, String outputSchemaSha256,
                                 String sourcePayloadSha256, String mappedPayloadSha256,
                                 String mappedPayloadJson) { }
}
````

### FILE: `debezium_postgres_inbox_consumer/src/main/java/com/elite/documentinbox/DocumentIntakeEnvelope.java`
```yaml
block_id: "DEBEZIUM-POSTGRES-INBOX-CONSUMER:src-main-java-com-elite-documentinbox-documentintakeenvelope-java:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://debezium-postgres-inbox-consumer/src/main/java/com/elite/documentinbox/DocumentIntakeEnvelope.java"
license: "LicenseRef-Workspace-Owner"
sha256: "0a2dfa5db81de63bd1b00b3e44619b5dd98b02851e17486a398f4b37b8d6205c"
variables: []
secrets_allowed: false
```
````java
package com.elite.documentinbox;

import java.util.UUID;

public record DocumentIntakeEnvelope(
        UUID eventId,
        UUID tenantId,
        int schemaVersion,
        String eventType,
        String aggregateId,
        String snapshotSha256,
        String domainType,
        String payloadJson,
        String payloadSha256,
        String topic,
        int partition,
        long offset) {
}
````

### FILE: `debezium_postgres_inbox_consumer/src/main/java/com/elite/documentinbox/DocumentIntakeEnvelopeValidator.java`
```yaml
block_id: "DEBEZIUM-POSTGRES-INBOX-CONSUMER:src-main-java-com-elite-documentinbox-documentintakeenvelopevalidator-java:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://debezium-postgres-inbox-consumer/src/main/java/com/elite/documentinbox/DocumentIntakeEnvelopeValidator.java"
license: "LicenseRef-Workspace-Owner"
sha256: "415ef1fed845530d0fdf58fe5d9c563c12063448db4ca533e65808292e838b76"
variables: []
secrets_allowed: false
```
````java
package com.elite.documentinbox;

import com.fasterxml.jackson.core.JsonFactory;
import com.fasterxml.jackson.core.StreamReadFeature;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.kafka.common.header.Header;
import org.apache.kafka.common.header.Headers;

import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.util.HashSet;
import java.util.Iterator;
import java.util.Set;
import java.util.UUID;
import java.util.regex.Pattern;

public final class DocumentIntakeEnvelopeValidator {
    static final String EVENT_TYPE = "document.intake.accepted";
    private static final Pattern SHA256 = Pattern.compile("^[0-9a-f]{64}$");
    private static final Pattern DOMAIN = Pattern.compile("^[a-z][a-z0-9._-]{0,126}[a-z0-9]$");
    private static final int MAX_PAYLOAD_BYTES = 65_536;
    private static final Set<String> PAYLOAD_KEYS = Set.of("evaluationId", "snapshotSha256", "domainType");
    private static final ObjectMapper JSON = new ObjectMapper(JsonFactory.builder()
            .enable(StreamReadFeature.STRICT_DUPLICATE_DETECTION).build());

    private DocumentIntakeEnvelopeValidator() {
    }

    public static DocumentIntakeEnvelope validate(String expectedTopic, String topic, int partition, long offset,
                                                    String aggregateKey, String payload, Headers headers) {
        if (expectedTopic == null || expectedTopic.isBlank() || !expectedTopic.equals(topic)) {
            throw new IllegalArgumentException("Unexpected Kafka topic");
        }
        if (partition < 0 || offset < 0 || aggregateKey == null || !SHA256.matcher(aggregateKey).matches()) {
            throw new IllegalArgumentException("Invalid Kafka identity");
        }
        if (payload == null || payload.isBlank() || payload.getBytes(StandardCharsets.UTF_8).length > MAX_PAYLOAD_BYTES) {
            throw new IllegalArgumentException("Invalid payload size");
        }

        UUID eventId = parseUuid(exactHeader(headers, "id"), "id");
        UUID tenantId = parseUuid(exactHeader(headers, "tenant_id"), "tenant_id");
        String eventType = exactHeader(headers, "type");
        if (!EVENT_TYPE.equals(eventType)) {
            throw new IllegalArgumentException("Unsupported event type");
        }
        int schemaVersion;
        try {
            schemaVersion = Integer.parseInt(exactHeader(headers, "schema_version"));
        }
        catch (NumberFormatException e) {
            throw new IllegalArgumentException("Invalid schema version", e);
        }
        if (schemaVersion != 1) {
            throw new IllegalArgumentException("Unsupported schema version");
        }

        JsonNode node;
        try {
            node = JSON.readTree(payload);
        }
        catch (Exception e) {
            throw new IllegalArgumentException("Invalid JSON payload", e);
        }
        if (node == null || !node.isObject() || !fieldNames(node).equals(PAYLOAD_KEYS)) {
            throw new IllegalArgumentException("Payload schema is not closed");
        }
        String evaluationId = requiredText(node, "evaluationId");
        String snapshotSha256 = requiredText(node, "snapshotSha256");
        String domainType = requiredText(node, "domainType");
        if (!aggregateKey.equals(evaluationId) || !SHA256.matcher(snapshotSha256).matches()
                || !DOMAIN.matcher(domainType).matches()) {
            throw new IllegalArgumentException("Payload identity mismatch");
        }
        return new DocumentIntakeEnvelope(eventId, tenantId, schemaVersion, eventType, evaluationId,
                snapshotSha256, domainType, payload, sha256(payload), topic, partition, offset);
    }

    static String exactHeader(Headers headers, String name) {
        Iterator<Header> iterator = headers.headers(name).iterator();
        if (!iterator.hasNext()) {
            throw new IllegalArgumentException("Missing header " + name);
        }
        Header header = iterator.next();
        if (iterator.hasNext() || header.value() == null) {
            throw new IllegalArgumentException("Header must occur exactly once: " + name);
        }
        String value = new String(header.value(), StandardCharsets.UTF_8);
        if (value.isBlank() || !value.equals(value.trim())) {
            throw new IllegalArgumentException("Invalid header " + name);
        }
        return value;
    }

    private static Set<String> fieldNames(JsonNode node) {
        Set<String> names = new HashSet<>();
        node.fieldNames().forEachRemaining(names::add);
        return names;
    }

    private static String requiredText(JsonNode node, String field) {
        JsonNode value = node.get(field);
        if (value == null || !value.isTextual() || value.textValue().isBlank()) {
            throw new IllegalArgumentException("Invalid field " + field);
        }
        return value.textValue();
    }

    private static UUID parseUuid(String value, String field) {
        try {
            UUID parsed = UUID.fromString(value);
            if (!parsed.toString().equals(value)) {
                throw new IllegalArgumentException("UUID is not canonical");
            }
            return parsed;
        }
        catch (IllegalArgumentException e) {
            throw new IllegalArgumentException("Invalid UUID header " + field, e);
        }
    }

    private static String sha256(String value) {
        try {
            return java.util.HexFormat.of().formatHex(MessageDigest.getInstance("SHA-256")
                    .digest(value.getBytes(StandardCharsets.UTF_8)));
        }
        catch (Exception e) {
            throw new IllegalStateException("SHA-256 unavailable", e);
        }
    }
}
````

### FILE: `debezium_postgres_inbox_consumer/src/main/java/com/elite/documentinbox/JdbcDocumentIntakeProcessor.java`
```yaml
block_id: "DEBEZIUM-POSTGRES-INBOX-CONSUMER:src-main-java-com-elite-documentinbox-jdbcdocumentintakeprocessor-java:v1"
operation: CREATE
provenance: ADAPTED
source: "https://github.com/debezium/debezium-examples/tree/7b0d765a02cf66ef29de2a01b67ec763e9acedf5/outbox"
license: "LicenseRef-Workspace-Owner"
sha256: "fb571ec0d38dd20ee0db2a10f07ff845aad9c320d0db5f08ca76728c3ea27400"
variables: []
secrets_allowed: false
```
````java
package com.elite.documentinbox;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import org.eclipse.microprofile.config.inject.ConfigProperty;

import javax.sql.DataSource;
import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.sql.SQLException;

@ApplicationScoped
public class JdbcDocumentIntakeProcessor {
    public enum Result { APPLIED, DUPLICATE }

    private final DataSource dataSource;
    private final String consumerName;
    private final CelDocumentProjectionMapping mapping;

    @Inject
    public JdbcDocumentIntakeProcessor(DataSource dataSource,
                                       @ConfigProperty(name = "elite.document-inbox.consumer-name") String consumerName,
                                       CelDocumentProjectionMapping mapping) {
        if (consumerName == null || !consumerName.matches("^[a-z][a-z0-9._-]{2,126}$")) {
            throw new IllegalArgumentException("Invalid consumer name");
        }
        this.dataSource = dataSource;
        this.consumerName = consumerName;
        this.mapping = mapping;
    }

    public Result process(DocumentIntakeEnvelope event) {
        try (Connection connection = dataSource.getConnection()) {
            connection.setAutoCommit(false);
            connection.setTransactionIsolation(Connection.TRANSACTION_READ_COMMITTED);
            try {
                setTenant(connection, event);
                int inserted = insertInbox(connection, event);
                if (inserted == 0) {
                    verifyDuplicate(connection, event);
                    applyProjection(connection, event);
                    connection.commit();
                    return Result.DUPLICATE;
                }
                applyProjection(connection, event);
                connection.commit();
                return Result.APPLIED;
            }
            catch (RuntimeException | SQLException e) {
                rollback(connection, e);
                throw e;
            }
        }
        catch (SQLException e) {
            throw new IllegalStateException("Database processing failed", e);
        }
    }

    private void setTenant(Connection connection, DocumentIntakeEnvelope event) throws SQLException {
        try (PreparedStatement statement = connection.prepareStatement("select set_config('app.tenant_id', ?, true)")) {
            statement.setString(1, event.tenantId().toString());
            statement.executeQuery().close();
        }
    }

    private int insertInbox(Connection connection, DocumentIntakeEnvelope event) throws SQLException {
        String sql = """
                insert into document_intelligence.consumer_inbox
                (consumer_name,tenant_id,event_id,payload_sha256_hex,topic,partition_id,offset_id,aggregate_id,event_type,schema_version)
                values (?,?,?,?,?,?,?,?,?,?) on conflict do nothing
                """;
        try (PreparedStatement statement = connection.prepareStatement(sql)) {
            statement.setString(1, consumerName);
            statement.setObject(2, event.tenantId());
            statement.setObject(3, event.eventId());
            statement.setString(4, event.payloadSha256());
            statement.setString(5, event.topic());
            statement.setInt(6, event.partition());
            statement.setLong(7, event.offset());
            statement.setString(8, event.aggregateId());
            statement.setString(9, event.eventType());
            statement.setInt(10, event.schemaVersion());
            return statement.executeUpdate();
        }
    }

    private void verifyDuplicate(Connection connection, DocumentIntakeEnvelope event) throws SQLException {
        String sql = """
                select payload_sha256_hex,topic,partition_id,offset_id,aggregate_id,event_type,schema_version
                from document_intelligence.consumer_inbox
                where consumer_name=? and tenant_id=? and event_id=?
                """;
        try (PreparedStatement statement = connection.prepareStatement(sql)) {
            statement.setString(1, consumerName);
            statement.setObject(2, event.tenantId());
            statement.setObject(3, event.eventId());
            try (ResultSet row = statement.executeQuery()) {
                if (!row.next() || !row.getString(1).equals(event.payloadSha256())
                        || !row.getString(2).equals(event.topic()) || row.getInt(3) != event.partition()
                        || row.getLong(4) != event.offset() || !row.getString(5).equals(event.aggregateId())
                        || !row.getString(6).equals(event.eventType()) || row.getInt(7) != event.schemaVersion()
                        || row.next()) {
                    throw new IllegalStateException("Divergent duplicate event");
                }
            }
        }
    }

    private void applyProjection(Connection connection, DocumentIntakeEnvelope event) throws SQLException {
        CelDocumentProjectionMapping.MappedDocument mapped = mapAuthoritativeDocument(connection, event);
        String insert = """
                insert into document_intelligence.accepted_document_projection
                (tenant_id,evaluation_id,snapshot_sha256_hex,domain_type,event_id,event_payload_sha256_hex,
                 source_document_payload_sha256_hex,mapping_id,mapping_version,mapping_expression_sha256_hex,
                 input_schema_id,input_schema_sha256_hex,output_schema_id,output_schema_sha256_hex,
                 mapped_payload_sha256_hex,mapped_payload)
                values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?::jsonb)
                on conflict (tenant_id,evaluation_id) do nothing
                """;
        int applied;
        try (PreparedStatement statement = connection.prepareStatement(insert)) {
            statement.setObject(1, event.tenantId());
            statement.setString(2, event.aggregateId());
            statement.setString(3, event.snapshotSha256());
            statement.setString(4, mapped.mappingDomainType());
            statement.setObject(5, event.eventId());
            statement.setString(6, event.payloadSha256());
            statement.setString(7, mapped.sourcePayloadSha256());
            statement.setString(8, mapped.mappingId());
            statement.setString(9, mapped.mappingVersion());
            statement.setString(10, mapped.expressionSha256());
            statement.setString(11, mapped.inputSchemaId());
            statement.setString(12, mapped.inputSchemaSha256());
            statement.setString(13, mapped.outputSchemaId());
            statement.setString(14, mapped.outputSchemaSha256());
            statement.setString(15, mapped.mappedPayloadSha256());
            statement.setString(16, mapped.mappedPayloadJson());
            applied = statement.executeUpdate();
        }
        if (applied == 1) {
            return;
        }
        String verify = """
                select snapshot_sha256_hex,domain_type,event_id,event_payload_sha256_hex,
                       source_document_payload_sha256_hex,mapping_id,mapping_version,mapping_expression_sha256_hex,
                       input_schema_id,input_schema_sha256_hex,output_schema_id,output_schema_sha256_hex,
                       mapped_payload_sha256_hex,mapped_payload::text
                from document_intelligence.accepted_document_projection
                where tenant_id=? and evaluation_id=?
                """;
        try (PreparedStatement statement = connection.prepareStatement(verify)) {
            statement.setObject(1, event.tenantId());
            statement.setString(2, event.aggregateId());
            try (ResultSet row = statement.executeQuery()) {
                if (!row.next() || !row.getString(1).equals(event.snapshotSha256())
                        || !row.getString(2).equals(event.domainType())
                        || !row.getObject(3).equals(event.eventId())
                        || !row.getString(4).equals(event.payloadSha256())
                        || !row.getString(5).equals(mapped.sourcePayloadSha256())
                        || !row.getString(6).equals(mapped.mappingId())
                        || !row.getString(7).equals(mapped.mappingVersion())
                        || !row.getString(8).equals(mapped.expressionSha256())
                        || !row.getString(9).equals(mapped.inputSchemaId())
                        || !row.getString(10).equals(mapped.inputSchemaSha256())
                        || !row.getString(11).equals(mapped.outputSchemaId())
                        || !row.getString(12).equals(mapped.outputSchemaSha256())
                        || !row.getString(13).equals(mapped.mappedPayloadSha256())
                        || !jsonEquals(row.getString(14), mapped.mappedPayloadJson()) || row.next()) {
                    throw new IllegalStateException("Authoritative intake missing or projection divergent");
                }
            }
        }
    }

    private CelDocumentProjectionMapping.MappedDocument mapAuthoritativeDocument(
            Connection connection, DocumentIntakeEnvelope event) throws SQLException {
        String sql = """
                select document_payload::text,document_payload_sha256_hex
                from document_intelligence.intake_document
                where tenant_id=? and evaluation_id=? and snapshot_sha256_hex=? and domain_type=?
                """;
        try (PreparedStatement statement = connection.prepareStatement(sql)) {
            statement.setObject(1, event.tenantId());
            statement.setString(2, event.aggregateId());
            statement.setString(3, event.snapshotSha256());
            statement.setString(4, event.domainType());
            try (ResultSet row = statement.executeQuery()) {
                if (!row.next()) throw new IllegalStateException("Authoritative intake missing");
                String documentJson = row.getString(1);
                String documentSha256 = row.getString(2);
                if (row.next()) throw new IllegalStateException("Authoritative intake is not unique");
                return mapping.apply(event.domainType(), documentJson, documentSha256);
            }
        }
    }

    private static boolean jsonEquals(String left, String right) {
        try {
            var mapper = new com.fasterxml.jackson.databind.ObjectMapper();
            return mapper.readTree(left).equals(mapper.readTree(right));
        }
        catch (Exception e) {
            throw new IllegalStateException("Stored mapped payload is invalid JSON", e);
        }
    }

    private static void rollback(Connection connection, Exception original) {
        try {
            connection.rollback();
        }
        catch (SQLException rollbackFailure) {
            original.addSuppressed(rollbackFailure);
        }
    }
}
````

### FILE: `debezium_postgres_inbox_consumer/src/main/java/com/elite/documentinbox/KafkaDocumentIntakeConsumer.java`
```yaml
block_id: "DEBEZIUM-POSTGRES-INBOX-CONSUMER:src-main-java-com-elite-documentinbox-kafkadocumentintakeconsumer-java:v1"
operation: CREATE
provenance: ADAPTED
source: "https://github.com/debezium/debezium-examples/tree/7b0d765a02cf66ef29de2a01b67ec763e9acedf5/outbox"
license: "LicenseRef-Workspace-Owner"
sha256: "df239f0a5dd7a147bb6e0752dd2a64dae232760a2db2245e834e41cbb6ca015f"
variables: []
secrets_allowed: false
```
````java
package com.elite.documentinbox;

import io.smallrye.reactive.messaging.kafka.KafkaRecord;
import io.smallrye.reactive.messaging.kafka.api.IncomingKafkaRecordMetadata;
import io.smallrye.reactive.messaging.annotations.Blocking;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import org.eclipse.microprofile.config.inject.ConfigProperty;
import org.eclipse.microprofile.reactive.messaging.Acknowledgment;
import org.eclipse.microprofile.reactive.messaging.Incoming;
import org.jboss.logging.Logger;

import java.util.concurrent.CompletionStage;

@ApplicationScoped
public class KafkaDocumentIntakeConsumer {
    private static final Logger LOG = Logger.getLogger(KafkaDocumentIntakeConsumer.class);

    @Inject
    JdbcDocumentIntakeProcessor processor;

    @ConfigProperty(name = "elite.document-inbox.topic")
    String expectedTopic;

    @Incoming("document-intake")
    @Acknowledgment(Acknowledgment.Strategy.MANUAL)
    @Blocking(ordered = true)
    public CompletionStage<Void> onMessage(KafkaRecord<String, String> message) {
        IncomingKafkaRecordMetadata<?, ?> metadata = message.getMetadata(IncomingKafkaRecordMetadata.class)
                .orElseThrow(() -> new IllegalArgumentException("Kafka metadata missing"));
        DocumentIntakeEnvelope event = DocumentIntakeEnvelopeValidator.validate(expectedTopic, message.getTopic(),
                message.getPartition(), metadata.getOffset(), message.getKey(), message.getPayload(), message.getHeaders());
        JdbcDocumentIntakeProcessor.Result result = processor.process(event);
        LOG.infof("document intake event outcome=%s partition=%d", result, event.partition());
        return message.ack();
    }
}
````

### FILE: `debezium_postgres_inbox_consumer/src/main/resources/application.properties`
```yaml
block_id: "DEBEZIUM-POSTGRES-INBOX-CONSUMER:src-main-resources-application-properties:v1"
operation: CREATE
provenance: ADAPTED
source: "https://github.com/debezium/debezium-examples/tree/7b0d765a02cf66ef29de2a01b67ec763e9acedf5/outbox"
license: "LicenseRef-Workspace-Owner"
sha256: "25863cae265cd6359c1435454c5a1c34a384105d63ec24af8055e569c1d1ff5c"
variables: []
secrets_allowed: false
```
````properties
quarkus.datasource.db-kind=postgresql
quarkus.datasource.jdbc.url=${ELITE_DATABASE_JDBC_URL}
quarkus.datasource.username=${ELITE_DATABASE_USERNAME}
quarkus.datasource.password=${ELITE_DATABASE_PASSWORD}
quarkus.datasource.jdbc.transactions=disabled

mp.messaging.incoming.document-intake.connector=smallrye-kafka
mp.messaging.incoming.document-intake.topic=${ELITE_KAFKA_TOPIC}
mp.messaging.incoming.document-intake.group.id=${ELITE_KAFKA_CONSUMER_GROUP}
mp.messaging.incoming.document-intake.bootstrap.servers=${ELITE_KAFKA_BOOTSTRAP_SERVERS}
mp.messaging.incoming.document-intake.enable.auto.commit=false
mp.messaging.incoming.document-intake.auto.offset.reset=earliest
mp.messaging.incoming.document-intake.isolation.level=read_committed
mp.messaging.incoming.document-intake.failure-strategy=fail
mp.messaging.incoming.document-intake.security.protocol=SSL
mp.messaging.incoming.document-intake.ssl.truststore.location=${ELITE_KAFKA_TRUSTSTORE_PATH}
mp.messaging.incoming.document-intake.ssl.truststore.password=${ELITE_KAFKA_TRUSTSTORE_PASSWORD}
mp.messaging.incoming.document-intake.ssl.keystore.location=${ELITE_KAFKA_KEYSTORE_PATH}
mp.messaging.incoming.document-intake.ssl.keystore.password=${ELITE_KAFKA_KEYSTORE_PASSWORD}
elite.document-inbox.topic=${ELITE_KAFKA_TOPIC}
elite.document-inbox.consumer-name=${ELITE_KAFKA_CONSUMER_GROUP}
elite.document-inbox.mapping-id=${ELITE_MAPPING_ID}
elite.document-inbox.mapping-domain-type=${ELITE_MAPPING_DOMAIN_TYPE}
elite.document-inbox.mapping-version=${ELITE_MAPPING_VERSION}
elite.document-inbox.mapping-expression=${ELITE_MAPPING_EXPRESSION}
elite.document-inbox.mapping-expression-sha256=${ELITE_MAPPING_EXPRESSION_SHA256}
elite.document-inbox.mapping-output-keys=${ELITE_MAPPING_EXACT_OUTPUT_KEYS_CSV}
elite.document-inbox.mapping-input-schema-id=${ELITE_MAPPING_INPUT_SCHEMA_ID}
elite.document-inbox.mapping-input-schema-sha256=${ELITE_MAPPING_INPUT_SCHEMA_SHA256}
elite.document-inbox.mapping-output-schema-id=${ELITE_MAPPING_OUTPUT_SCHEMA_ID}
elite.document-inbox.mapping-output-schema-sha256=${ELITE_MAPPING_OUTPUT_SCHEMA_SHA256}

quarkus.log.console.json.enabled=true
quarkus.log.category."com.elite.documentinbox".level=INFO
````

### FILE: `debezium_postgres_inbox_consumer/src/main/resources/db/migration/V001__document_consumer_inbox.sql`
```yaml
block_id: "DEBEZIUM-POSTGRES-INBOX-CONSUMER:src-main-resources-db-migration-v001-document-consumer-inbox-sql:v1"
operation: CREATE
provenance: ADAPTED
source: "https://github.com/debezium/debezium-examples/tree/7b0d765a02cf66ef29de2a01b67ec763e9acedf5/outbox"
license: "LicenseRef-Workspace-Owner"
sha256: "ec8e041bccd303361a2ef93722f2d6c7f8ad77314f793ec9d92d49df5061dd4c"
variables: []
secrets_allowed: false
```
````sql
begin;

create table if not exists document_intelligence.consumer_inbox (
  consumer_name text not null,
  tenant_id uuid not null references platform.tenant (tenant_id),
  event_id uuid not null,
  payload_sha256_hex text not null check (payload_sha256_hex ~ '^[0-9a-f]{64}$'),
  topic text not null,
  partition_id integer not null check (partition_id >= 0),
  offset_id bigint not null check (offset_id >= 0),
  aggregate_id text not null check (aggregate_id ~ '^[0-9a-f]{64}$'),
  event_type text not null check (event_type = 'document.intake.accepted'),
  schema_version integer not null check (schema_version = 1),
  processed_at timestamptz not null default clock_timestamp(),
  primary key (consumer_name, tenant_id, event_id),
  unique (consumer_name, topic, partition_id, offset_id)
);

create table if not exists document_intelligence.accepted_document_projection (
  tenant_id uuid not null,
  evaluation_id text not null,
  snapshot_sha256_hex text not null check (snapshot_sha256_hex ~ '^[0-9a-f]{64}$'),
  domain_type text not null,
  event_id uuid not null,
  event_payload_sha256_hex text not null check (event_payload_sha256_hex ~ '^[0-9a-f]{64}$'),
  source_document_payload_sha256_hex text not null check (source_document_payload_sha256_hex ~ '^[0-9a-f]{64}$'),
  mapping_id text not null,
  mapping_version text not null,
  mapping_expression_sha256_hex text not null check (mapping_expression_sha256_hex ~ '^[0-9a-f]{64}$'),
  input_schema_id text not null,
  input_schema_sha256_hex text not null check (input_schema_sha256_hex ~ '^[0-9a-f]{64}$'),
  output_schema_id text not null,
  output_schema_sha256_hex text not null check (output_schema_sha256_hex ~ '^[0-9a-f]{64}$'),
  mapped_payload_sha256_hex text not null check (mapped_payload_sha256_hex ~ '^[0-9a-f]{64}$'),
  mapped_payload jsonb not null check (jsonb_typeof(mapped_payload) = 'object'),
  accepted_at timestamptz not null default clock_timestamp(),
  primary key (tenant_id, evaluation_id),
  foreign key (tenant_id, evaluation_id) references document_intelligence.intake_document (tenant_id, evaluation_id),
  foreign key (tenant_id) references platform.tenant (tenant_id)
);

alter table document_intelligence.consumer_inbox enable row level security;
alter table document_intelligence.accepted_document_projection enable row level security;

do $$
begin
  if not exists (select 1 from pg_roles where rolname = 'elite_document_consumer') then
    create role elite_document_consumer nologin;
  end if;
end
$$;

grant usage on schema document_intelligence to elite_document_consumer;
grant select on document_intelligence.intake_document to elite_document_consumer;
grant select, insert on document_intelligence.consumer_inbox to elite_document_consumer;
grant select, insert on document_intelligence.accepted_document_projection to elite_document_consumer;

drop policy if exists consumer_inbox_tenant_isolation on document_intelligence.consumer_inbox;
create policy consumer_inbox_tenant_isolation on document_intelligence.consumer_inbox
  using (tenant_id = nullif(current_setting('app.tenant_id', true), '')::uuid)
  with check (tenant_id = nullif(current_setting('app.tenant_id', true), '')::uuid);

drop policy if exists accepted_projection_tenant_isolation on document_intelligence.accepted_document_projection;
create policy accepted_projection_tenant_isolation on document_intelligence.accepted_document_projection
  using (tenant_id = nullif(current_setting('app.tenant_id', true), '')::uuid)
  with check (tenant_id = nullif(current_setting('app.tenant_id', true), '')::uuid);

create index if not exists consumer_inbox_processed_idx
  on document_intelligence.consumer_inbox (tenant_id, processed_at, event_id);

commit;
````

### FILE: `debezium_postgres_inbox_consumer/src/test/java/com/elite/documentinbox/CelDocumentProjectionMappingTest.java`
```yaml
block_id: "DEBEZIUM-POSTGRES-INBOX-CONSUMER:cel-document-projection-mapping-test-java:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://debezium-postgres-inbox-consumer/src/test/java/com/elite/documentinbox/CelDocumentProjectionMappingTest.java"
license: "LicenseRef-Workspace-Owner"
sha256: "9686d84705e76d4f1c26b779a8357437cdd491bce8cdea648c3577b5191fe663"
variables: []
secrets_allowed: false
```
````java
package com.elite.documentinbox;

import dev.cel.common.CelValidationException;
import org.junit.jupiter.api.Test;

import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.util.HexFormat;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

final class CelDocumentProjectionMappingTest {
    @Test
    void mapsAndBindsSchemasIdentityAndHashes() throws Exception {
        String expression = "{'invoice_number': document['invoiceNumber'], 'total_minor': document['totalMinor']}";
        var mapping = mapping(expression, sha256(expression), "invoice_number,total_minor");
        var result = mapping.apply("invoice", "{\"invoiceNumber\":\"INV-42\",\"totalMinor\":129900,\"ignored\":true}", "a".repeat(64));
        assertEquals("invoice-to-ledger", result.mappingId());
        assertEquals("1.0.0", result.mappingVersion());
        assertEquals("invoice-reviewed-v3", result.inputSchemaId());
        assertEquals("ledger-invoice-v1", result.outputSchemaId());
        assertEquals("{\"invoice_number\":\"INV-42\",\"total_minor\":129900}", result.mappedPayloadJson());
    }

    @Test
    void rejectsHashMultilineAndEmptySchemaAuthority() throws Exception {
        String expression = "{'x': document['x']}";
        assertThrows(IllegalArgumentException.class, () -> mapping(expression, "0".repeat(64), "x"));
        String multiline = expression + "\n";
        assertThrows(IllegalArgumentException.class, () -> mapping(multiline, sha256(multiline), "x"));
        assertThrows(IllegalArgumentException.class, () -> new CelDocumentProjectionMapping(
                "invoice-to-ledger", "invoice", "1.0.0", expression, sha256(expression), "x", "", "b".repeat(64),
                "ledger-invoice-v1", "c".repeat(64)));
    }

    @Test
    void rejectsMacrosAtCompilation() throws Exception {
        String expression = "{'items': document['items'].map(x, x)}";
        assertThrows(CelValidationException.class, () -> mapping(expression, sha256(expression), "items"));
    }

    @Test
    void rejectsMissingInputAndExtraOutput() throws Exception {
        String missing = "{'x': document['x']}";
        var missingMapping = mapping(missing, sha256(missing), "x");
        assertThrows(IllegalArgumentException.class, () -> missingMapping.apply("invoice", "{}", "a".repeat(64)));
        String extra = "{'x': document['x'], 'secret': document['secret']}";
        var extraMapping = mapping(extra, sha256(extra), "x");
        assertThrows(IllegalStateException.class,
                () -> extraMapping.apply("invoice", "{\"x\":1,\"secret\":2}", "a".repeat(64)));
    }

    @Test
    void rejectsDuplicateJsonKeyFloatAndNull() throws Exception {
        String expression = "{'x': document['x']}";
        var mapping = mapping(expression, sha256(expression), "x");
        assertThrows(IllegalArgumentException.class, () -> mapping.apply("invoice", "{\"x\":1,\"x\":2}", "a".repeat(64)));
        assertThrows(IllegalArgumentException.class, () -> mapping.apply("invoice", "{\"x\":1.5}", "a".repeat(64)));
        assertThrows(IllegalArgumentException.class, () -> mapping.apply("invoice", "{\"x\":null}", "a".repeat(64)));
    }

    @Test
    void rejectsOversizedExpressionAndBadSourceHash() throws Exception {
        String oversized = "{'x':'" + "a".repeat(4097) + "'}";
        assertThrows(CelValidationException.class, () -> mapping(oversized, sha256(oversized), "x"));
        String expression = "{'x': document['x']}";
        var mapping = mapping(expression, sha256(expression), "x");
        assertThrows(IllegalArgumentException.class, () -> mapping.apply("invoice", "{\"x\":1}", "bad"));
        assertThrows(IllegalArgumentException.class,
                () -> mapping.apply("packing-list", "{\"x\":1}", "a".repeat(64)));
    }

    private static CelDocumentProjectionMapping mapping(String expression, String hash, String keys) throws Exception {
        return new CelDocumentProjectionMapping("invoice-to-ledger", "invoice", "1.0.0", expression, hash, keys,
                "invoice-reviewed-v3", "b".repeat(64), "ledger-invoice-v1", "c".repeat(64));
    }

    private static String sha256(String value) throws Exception {
        return HexFormat.of().formatHex(MessageDigest.getInstance("SHA-256")
                .digest(value.getBytes(StandardCharsets.UTF_8)));
    }
}
````

### FILE: `debezium_postgres_inbox_consumer/src/test/java/com/elite/documentinbox/DocumentIntakeEnvelopeValidatorTest.java`
```yaml
block_id: "DEBEZIUM-POSTGRES-INBOX-CONSUMER:src-test-java-com-elite-documentinbox-documentintakeenvelopevalidatortest-java:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://debezium-postgres-inbox-consumer/src/test/java/com/elite/documentinbox/DocumentIntakeEnvelopeValidatorTest.java"
license: "LicenseRef-Workspace-Owner"
sha256: "329895c02a89e81f8fbe29b8fb940b703e991276756f9e913b134fd9d9f20457"
variables: []
secrets_allowed: false
```
````java
package com.elite.documentinbox;

import org.apache.kafka.common.header.internals.RecordHeaders;
import org.junit.jupiter.api.Test;

import java.nio.charset.StandardCharsets;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

class DocumentIntakeEnvelopeValidatorTest {
    private static final String EVALUATION = "a".repeat(64);
    private static final String SNAPSHOT = "b".repeat(64);
    private static final String PAYLOAD = "{\"evaluationId\":\"" + EVALUATION
            + "\",\"snapshotSha256\":\"" + SNAPSHOT + "\",\"domainType\":\"invoice\"}";

    @Test
    void validatesClosedIdentityAndHashes() {
        var event = DocumentIntakeEnvelopeValidator.validate("document.events", "document.events", 2, 9,
                EVALUATION, PAYLOAD, headers());
        assertEquals(EVALUATION, event.aggregateId());
        assertEquals("invoice", event.domainType());
        assertEquals(64, event.payloadSha256().length());
    }

    @Test
    void rejectsDuplicateOrMissingHeaders() {
        var duplicate = headers().add("id", bytes("018f76c6-4f0d-7000-a000-000000000001"));
        assertThrows(IllegalArgumentException.class, () -> DocumentIntakeEnvelopeValidator.validate(
                "document.events", "document.events", 0, 0, EVALUATION, PAYLOAD, duplicate));
        var missing = new RecordHeaders().add("id", bytes("018f76c6-4f0d-7000-a000-000000000001"));
        assertThrows(IllegalArgumentException.class, () -> DocumentIntakeEnvelopeValidator.validate(
                "document.events", "document.events", 0, 0, EVALUATION, PAYLOAD, missing));
    }

    @Test
    void rejectsUnknownTypeSchemaAndTopic() {
        var type = headers("unknown.event", "1");
        assertThrows(IllegalArgumentException.class, () -> DocumentIntakeEnvelopeValidator.validate(
                "document.events", "document.events", 0, 0, EVALUATION, PAYLOAD, type));
        var schema = headers("document.intake.accepted", "2");
        assertThrows(IllegalArgumentException.class, () -> DocumentIntakeEnvelopeValidator.validate(
                "document.events", "document.events", 0, 0, EVALUATION, PAYLOAD, schema));
        assertThrows(IllegalArgumentException.class, () -> DocumentIntakeEnvelopeValidator.validate(
                "document.events", "other.events", 0, 0, EVALUATION, PAYLOAD, headers()));
    }

    @Test
    void rejectsDuplicateJsonKeysAndOpenSchema() {
        String duplicate = "{\"evaluationId\":\"" + EVALUATION + "\",\"evaluationId\":\"" + EVALUATION
                + "\",\"snapshotSha256\":\"" + SNAPSHOT + "\",\"domainType\":\"invoice\"}";
        assertThrows(IllegalArgumentException.class, () -> DocumentIntakeEnvelopeValidator.validate(
                "document.events", "document.events", 0, 0, EVALUATION, duplicate, headers()));
        String extra = PAYLOAD.substring(0, PAYLOAD.length() - 1) + ",\"extra\":true}";
        assertThrows(IllegalArgumentException.class, () -> DocumentIntakeEnvelopeValidator.validate(
                "document.events", "document.events", 0, 0, EVALUATION, extra, headers()));
    }

    @Test
    void rejectsAggregateSnapshotDomainAndOffsetMismatches() {
        assertThrows(IllegalArgumentException.class, () -> DocumentIntakeEnvelopeValidator.validate(
                "document.events", "document.events", 0, -1, EVALUATION, PAYLOAD, headers()));
        assertThrows(IllegalArgumentException.class, () -> DocumentIntakeEnvelopeValidator.validate(
                "document.events", "document.events", 0, 0, "c".repeat(64), PAYLOAD, headers()));
        String bad = PAYLOAD.replace(SNAPSHOT, "no");
        assertThrows(IllegalArgumentException.class, () -> DocumentIntakeEnvelopeValidator.validate(
                "document.events", "document.events", 0, 0, EVALUATION, bad, headers()));
    }

    private static RecordHeaders headers() {
        return headers("document.intake.accepted", "1");
    }

    private static RecordHeaders headers(String type, String schema) {
        RecordHeaders result = new RecordHeaders();
        result.add("id", bytes("018f76c6-4f0d-7000-a000-000000000001"));
        result.add("tenant_id", bytes("11111111-1111-4111-8111-111111111111"));
        result.add("type", bytes(type));
        result.add("schema_version", bytes(schema));
        return result;
    }

    private static byte[] bytes(String value) {
        return value.getBytes(StandardCharsets.UTF_8);
    }
}
````

### FILE: `debezium_postgres_inbox_consumer/src/test/java/com/elite/documentinbox/JdbcDocumentIntakeProcessorTest.java`
```yaml
block_id: "DEBEZIUM-POSTGRES-INBOX-CONSUMER:src-test-java-com-elite-documentinbox-jdbcdocumentintakeprocessortest-java:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://debezium-postgres-inbox-consumer/src/test/java/com/elite/documentinbox/JdbcDocumentIntakeProcessorTest.java"
license: "LicenseRef-Workspace-Owner"
sha256: "23b0f4d867ccaa9d88b9daceb51a2aecdb95e488ae2c49e3a2acdc2640541a86"
variables: []
secrets_allowed: false
```
````java
package com.elite.documentinbox;

import org.junit.jupiter.api.Assumptions;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.postgresql.ds.PGSimpleDataSource;

import javax.sql.DataSource;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.ResultSet;
import java.util.ArrayList;
import java.util.List;
import java.util.UUID;
import java.util.HexFormat;
import java.util.concurrent.Callable;
import java.util.concurrent.Executors;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

class JdbcDocumentIntakeProcessorTest {
    private static final UUID TENANT = UUID.fromString("11111111-1111-4111-8111-111111111111");
    private static String adminUrl;
    private static JdbcDocumentIntakeProcessor processor;
    private static final String MAPPING_EXPRESSION = "{'invoice_number': document['invoiceNumber'], 'total_minor': document['totalMinor']}";

    @BeforeAll
    static void prepareDatabase() throws Exception {
        adminUrl = System.getenv("ELITE_TEST_PG_URL");
        Assumptions.assumeTrue(adminUrl != null && adminUrl.startsWith("jdbc:postgresql://"),
                "ELITE_TEST_PG_URL is required for PostgreSQL integration tests");
        try (Connection connection = DriverManager.getConnection(adminUrl, "postgres", "")) {
            connection.createStatement().execute("""
                    create schema platform;
                    create table platform.tenant(tenant_id uuid primary key);
                    insert into platform.tenant values ('11111111-1111-4111-8111-111111111111');
                    create schema document_intelligence;
                    create table document_intelligence.intake_document(
                      tenant_id uuid not null references platform.tenant,
                      evaluation_id text not null,
                      snapshot_sha256_hex text not null,
                      domain_type text not null,
                      document_payload_sha256_hex text not null,
                      document_payload jsonb not null,
                      primary key(tenant_id,evaluation_id));
                    alter table document_intelligence.intake_document enable row level security;
                    create policy intake_test_policy on document_intelligence.intake_document
                      using (tenant_id=nullif(current_setting('app.tenant_id',true),'')::uuid);
                    """);
            String migration = Files.readString(Path.of("src/main/resources/db/migration/V001__document_consumer_inbox.sql"));
            connection.createStatement().execute(migration);
            connection.createStatement().execute("create role elite_test_consumer login; grant elite_document_consumer to elite_test_consumer");
        }
        PGSimpleDataSource source = new PGSimpleDataSource();
        source.setURL(adminUrl);
        source.setUser("elite_test_consumer");
        var mapping = new CelDocumentProjectionMapping("invoice-to-ledger", "invoice", "1.0.0", MAPPING_EXPRESSION,
                sha256(MAPPING_EXPRESSION), "invoice_number,total_minor", "invoice-reviewed-v3", "b".repeat(64),
                "ledger-invoice-v1", "c".repeat(64));
        processor = new JdbcDocumentIntakeProcessor(source, "document-projection-v1", mapping);
    }

    @Test
    void appliesOnceAndReconcilesIdenticalDuplicate() throws Exception {
        DocumentIntakeEnvelope event = event('a', 'b', "018f76c6-4f0d-7000-a000-000000000001", 1);
        insertIntake(event);
        assertEquals(JdbcDocumentIntakeProcessor.Result.APPLIED, processor.process(event));
        assertEquals(JdbcDocumentIntakeProcessor.Result.DUPLICATE, processor.process(event));
        assertEquals(1, count("consumer_inbox", event.eventId().toString()));
        assertEquals(1, count("accepted_document_projection", event.eventId().toString()));
        try (Connection connection = DriverManager.getConnection(adminUrl, "postgres", "")) {
            ResultSet result = connection.createStatement().executeQuery("select mapping_id,mapping_version,mapped_payload from document_intelligence.accepted_document_projection where event_id='" + event.eventId() + "'");
            result.next();
            assertEquals("invoice-to-ledger", result.getString(1));
            assertEquals("1.0.0", result.getString(2));
            assertEquals("{\"total_minor\": 129900, \"invoice_number\": \"INV-a\"}", result.getString(3));
        }
    }

    @Test
    void divergentDuplicateFailsAndPreservesOriginal() throws Exception {
        DocumentIntakeEnvelope event = event('c', 'd', "018f76c6-4f0d-7000-a000-000000000002", 2);
        insertIntake(event);
        assertEquals(JdbcDocumentIntakeProcessor.Result.APPLIED, processor.process(event));
        DocumentIntakeEnvelope divergent = new DocumentIntakeEnvelope(event.eventId(), event.tenantId(), 1,
                event.eventType(), event.aggregateId(), event.snapshotSha256(), event.domainType(),
                event.payloadJson() + " ", "e".repeat(64), event.topic(), event.partition(), event.offset() + 1);
        assertThrows(IllegalStateException.class, () -> processor.process(divergent));
        assertEquals(1, count("consumer_inbox", event.eventId().toString()));
    }

    @Test
    void missingAuthoritativeIntakeRollsBackInbox() throws Exception {
        DocumentIntakeEnvelope event = event('e', 'f', "018f76c6-4f0d-7000-a000-000000000003", 3);
        assertThrows(IllegalStateException.class, () -> processor.process(event));
        assertEquals(0, count("consumer_inbox", event.eventId().toString()));
    }

    @Test
    void concurrentRedeliveryAppliesOneEffect() throws Exception {
        var executor = Executors.newFixedThreadPool(8);
        try {
            for (int round = 0; round < 4; round++) {
                char identity = (char) ('1' + round);
                DocumentIntakeEnvelope event = event(identity, (char) ('5' + round),
                        "018f76c6-4f0d-7000-a000-00000000000" + (4 + round), 10 + round);
                insertIntake(event);
                List<Callable<JdbcDocumentIntakeProcessor.Result>> calls = new ArrayList<>();
                for (int i = 0; i < 8; i++) calls.add(() -> processor.process(event));
                var results = executor.invokeAll(calls).stream().map(future -> {
                    try { return future.get(); } catch (Exception e) { throw new RuntimeException(e); }
                }).toList();
                assertEquals(1, results.stream().filter(r -> r == JdbcDocumentIntakeProcessor.Result.APPLIED).count());
                assertEquals(7, results.stream().filter(r -> r == JdbcDocumentIntakeProcessor.Result.DUPLICATE).count());
                assertEquals(1, count("accepted_document_projection", event.eventId().toString()));
            }
        }
        finally {
            executor.shutdownNow();
        }
    }

    @Test
    void reusedKafkaOffsetWithDifferentEventIsRejected() throws Exception {
        DocumentIntakeEnvelope first = event('9', 'a', "018f76c6-4f0d-7000-a000-000000000009", 99);
        DocumentIntakeEnvelope second = event('8', 'b', "018f76c6-4f0d-7000-a000-000000000008", 99);
        insertIntake(first); insertIntake(second);
        assertEquals(JdbcDocumentIntakeProcessor.Result.APPLIED, processor.process(first));
        assertThrows(IllegalStateException.class, () -> processor.process(second));
        assertEquals(0, count("consumer_inbox", second.eventId().toString()));
    }

    @Test
    void mappingFailureRollsBackInboxAndProjection() throws Exception {
        DocumentIntakeEnvelope event = event('7', 'e', "018f76c6-4f0d-7000-a000-000000000017", 117);
        insertIntake(event, "{\"invoiceNumber\":\"INV-7\"}");
        assertThrows(IllegalArgumentException.class, () -> processor.process(event));
        assertEquals(0, count("consumer_inbox", event.eventId().toString()));
        assertEquals(0, count("accepted_document_projection", event.eventId().toString()));
    }

    @Test
    void crossClassEventRollsBackBeforeProjection() throws Exception {
        DocumentIntakeEnvelope base = event('6', 'f', "018f76c6-4f0d-7000-a000-000000000016", 116);
        DocumentIntakeEnvelope event = new DocumentIntakeEnvelope(base.eventId(), base.tenantId(), base.schemaVersion(),
                base.eventType(), base.aggregateId(), base.snapshotSha256(), "packing-list", base.payloadJson(),
                base.payloadSha256(), base.topic(), base.partition(), base.offset());
        insertIntake(event);
        assertThrows(IllegalArgumentException.class, () -> processor.process(event));
        assertEquals(0, count("consumer_inbox", event.eventId().toString()));
        assertEquals(0, count("accepted_document_projection", event.eventId().toString()));
    }

    private static DocumentIntakeEnvelope event(char evaluation, char snapshot, String eventId, long offset) {
        String aggregate = String.valueOf(evaluation).repeat(64);
        return new DocumentIntakeEnvelope(UUID.fromString(eventId), TENANT, 1, "document.intake.accepted",
                aggregate, String.valueOf(snapshot).repeat(64), "invoice", "{}", String.valueOf(evaluation).repeat(64),
                "document.events", 0, offset);
    }

    private static void insertIntake(DocumentIntakeEnvelope event) throws Exception {
        insertIntake(event, "{\"invoiceNumber\":\"INV-" + event.aggregateId().charAt(0) + "\",\"totalMinor\":129900}");
    }

    private static void insertIntake(DocumentIntakeEnvelope event, String documentJson) throws Exception {
        try (Connection connection = DriverManager.getConnection(adminUrl, "postgres", "")) {
            var statement = connection.prepareStatement("insert into document_intelligence.intake_document values (?,?,?,?,?,?::jsonb)");
            statement.setObject(1, TENANT); statement.setString(2, event.aggregateId());
            statement.setString(3, event.snapshotSha256()); statement.setString(4, event.domainType());
            statement.setString(5, sha256(documentJson)); statement.setString(6, documentJson);
            statement.executeUpdate();
        }
    }

    private static int count(String table, String eventId) throws Exception {
        try (Connection connection = DriverManager.getConnection(adminUrl, "postgres", "")) {
            ResultSet result = connection.createStatement().executeQuery("select count(*) from document_intelligence." + table
                    + " where event_id='" + UUID.fromString(eventId) + "'");
            result.next(); return result.getInt(1);
        }
    }

    private static String sha256(String value) throws Exception {
        return HexFormat.of().formatHex(MessageDigest.getInstance("SHA-256")
                .digest(value.getBytes(StandardCharsets.UTF_8)));
    }
}
````

### FILE: `debezium_postgres_inbox_consumer/upstream/ConsumedMessage.java`
```yaml
block_id: "DEBEZIUM-POSTGRES-INBOX-CONSUMER:upstream-consumedmessage-java:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/debezium/debezium-examples/7b0d765a02cf66ef29de2a01b67ec763e9acedf5/outbox/shipment-service/src/main/java/io/debezium/examples/outbox/shipment/messagelog/ConsumedMessage.java"
license: "Apache-2.0"
sha256: "aa78b54720fe26268141ea79720c3b1d691c6d57319d40eaacbd80608c9c72e1"
variables: []
secrets_allowed: false
```
````java
/*
 * Copyright Debezium Authors.
 *
 * Licensed under the Apache Software License version 2.0, available at http://www.apache.org/licenses/LICENSE-2.0
 */
package io.debezium.examples.outbox.shipment.messagelog;

import java.time.Instant;
import java.util.UUID;

import jakarta.persistence.Entity;
import jakarta.persistence.Id;

@Entity
public class ConsumedMessage {

    @Id
    private UUID eventId;

    private Instant timeOfReceiving;

    ConsumedMessage() {
    }

    public ConsumedMessage(UUID eventId, Instant timeOfReceiving) {
        this.eventId = eventId;
        this.timeOfReceiving = timeOfReceiving;
    }

    public UUID getEventId() {
        return eventId;
    }

    public void setEventId(UUID eventId) {
        this.eventId = eventId;
    }

    public Instant getTimeOfReceiving() {
        return timeOfReceiving;
    }

    public void setTimeOfReceiving(Instant timeOfReceiving) {
        this.timeOfReceiving = timeOfReceiving;
    }
}
````

### FILE: `debezium_postgres_inbox_consumer/upstream/KafkaEventConsumer.java`
```yaml
block_id: "DEBEZIUM-POSTGRES-INBOX-CONSUMER:upstream-kafkaeventconsumer-java:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/debezium/debezium-examples/7b0d765a02cf66ef29de2a01b67ec763e9acedf5/outbox/shipment-service/src/main/java/io/debezium/examples/outbox/shipment/facade/KafkaEventConsumer.java"
license: "Apache-2.0"
sha256: "783184dfe7d894a7f17676e78a53daa00c99de9f2fd9b2cabac86922b1d813a3"
variables: []
secrets_allowed: false
```
````java
/*
 * Copyright Debezium Authors.
 *
 * Licensed under the Apache Software License version 2.0, available at http://www.apache.org/licenses/LICENSE-2.0
 */
package io.debezium.examples.outbox.shipment.facade;

import java.io.IOException;
import java.nio.charset.Charset;
import java.util.UUID;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.CompletionStage;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;

import org.apache.kafka.common.header.Header;
import org.eclipse.microprofile.reactive.messaging.Acknowledgment;
import org.eclipse.microprofile.reactive.messaging.Incoming;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import io.opentelemetry.context.Scope;
import io.opentelemetry.api.trace.Tracer;
import io.opentelemetry.api.trace.Span;
import io.smallrye.reactive.messaging.kafka.KafkaRecord;

@ApplicationScoped
public class KafkaEventConsumer {

    private static final Logger LOG = LoggerFactory.getLogger(KafkaEventConsumer.class);

    @Inject
    OrderEventHandler orderEventHandler;

    @Inject
    Tracer tracer;

    @Inject
    Span span;

    @Incoming("orders")
    @Acknowledgment(Acknowledgment.Strategy.MANUAL)
    public CompletionStage<Void> onMessage(KafkaRecord<String, String> message) throws IOException {
        return CompletableFuture.runAsync(() -> {
                span = tracer.spanBuilder("onMessage").startSpan();
                try (final Scope scope = span.makeCurrent()){
                    LOG.debug("Kafka message with key = {} arrived", message.getKey());
                    span.addEvent("Kafka message with key = "+message.getKey()+" arrived");

                    String eventId = getHeaderAsString(message, "id");
                    String eventType = getHeaderAsString(message, "eventType");

                    orderEventHandler.onOrderEvent(
                            UUID.fromString(eventId),
                            eventType,
                            message.getKey(),
                            message.getPayload(),
                            message.getTimestamp()
                    );

                    message.ack();
                    span.addEvent("ack");
                }
                catch (Exception e) {
                    LOG.error("Error while preparing shipment");
                    span.addEvent("Error while preparing shipment");
                    throw e;
                }
                finally {
                    span.end();
                }
        });
    }

    private String getHeaderAsString(KafkaRecord<?, ?> record, String name) {
        Header header = record.getHeaders().lastHeader(name);
        if (header == null) {
            throw new IllegalArgumentException("Expected record header '" + name + "' not present");
        }

        return new String(header.value(), Charset.forName("UTF-8"));
    }
}
````

### FILE: `debezium_postgres_inbox_consumer/upstream/LICENSE.txt`
```yaml
block_id: "DEBEZIUM-POSTGRES-INBOX-CONSUMER:upstream-license-txt:v1"
operation: CREATE
provenance: ADAPTED
source: "https://raw.githubusercontent.com/debezium/debezium-examples/7b0d765a02cf66ef29de2a01b67ec763e9acedf5/LICENSE.txt"
license: "Apache-2.0"
sha256: "cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30"
variables: []
secrets_allowed: false
```
````text

                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "[]"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright [yyyy] [name of copyright owner]

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
````

### FILE: `debezium_postgres_inbox_consumer/upstream/MessageLog.java`
```yaml
block_id: "DEBEZIUM-POSTGRES-INBOX-CONSUMER:upstream-messagelog-java:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/debezium/debezium-examples/7b0d765a02cf66ef29de2a01b67ec763e9acedf5/outbox/shipment-service/src/main/java/io/debezium/examples/outbox/shipment/messagelog/MessageLog.java"
license: "Apache-2.0"
sha256: "b60bc4cf5fd7a4dd0d783dad04e2595fa039f862af3e1f6dede931af9b74e55b"
variables: []
secrets_allowed: false
```
````java
/*
 * Copyright Debezium Authors.
 *
 * Licensed under the Apache Software License version 2.0, available at http://www.apache.org/licenses/LICENSE-2.0
 */
package io.debezium.examples.outbox.shipment.messagelog;

import java.time.Instant;
import java.util.UUID;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.persistence.EntityManager;
import jakarta.persistence.PersistenceContext;
import jakarta.transaction.Transactional;
import jakarta.transaction.Transactional.TxType;
import jakarta.inject.Inject;

import io.opentelemetry.api.trace.Span;
import io.opentelemetry.api.trace.Tracer;
import io.opentelemetry.context.Scope;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@ApplicationScoped
public class MessageLog {
    private static final Logger LOG = LoggerFactory.getLogger(MessageLog.class);

    @Inject
    Tracer tracer;

    @Inject
    Span span;

    @PersistenceContext
    EntityManager entityManager;

    @Transactional(value=TxType.MANDATORY)
    public void processed(UUID eventId) {
        span = tracer.spanBuilder("processed").startSpan();
        try (final Scope scope = span.makeCurrent()){
            entityManager.persist(new ConsumedMessage(eventId, Instant.now()));
        } finally {
            span.end();
        }
    }

    @Transactional(value=TxType.MANDATORY)
    public boolean alreadyProcessed(UUID eventId) {
        LOG.debug("Looking for event with id {} in message log", eventId);
        span = tracer.spanBuilder("alreadyProcessed").startSpan();
        try (final Scope scope = span.makeCurrent()){
            return entityManager.find(ConsumedMessage.class, eventId) != null;
        } finally {
            span.end();
        }
    }
}
````

### FILE: `debezium_postgres_inbox_consumer/upstream/OrderEventHandler.java`
```yaml
block_id: "DEBEZIUM-POSTGRES-INBOX-CONSUMER:upstream-ordereventhandler-java:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/debezium/debezium-examples/7b0d765a02cf66ef29de2a01b67ec763e9acedf5/outbox/shipment-service/src/main/java/io/debezium/examples/outbox/shipment/facade/OrderEventHandler.java"
license: "Apache-2.0"
sha256: "ccd1de2cb73a6630a08355fddcbba979bea7b8f3650e6ddb16d2590f0b1198ab"
variables: []
secrets_allowed: false
```
````java
/*
 * Copyright Debezium Authors.
 *
 * Licensed under the Apache Software License version 2.0, available at http://www.apache.org/licenses/LICENSE-2.0
 */
package io.debezium.examples.outbox.shipment.facade;

import java.io.IOException;
import java.time.Instant;
import java.util.UUID;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import jakarta.transaction.Transactional;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;

import io.debezium.examples.outbox.shipment.messagelog.MessageLog;
import io.debezium.examples.outbox.shipment.service.ShipmentService;

@ApplicationScoped
public class OrderEventHandler {

    private static final Logger LOGGER = LoggerFactory.getLogger(OrderEventHandler.class);

    @Inject
    MessageLog log;

    @Inject
    ShipmentService shipmentService;

    private final ObjectMapper objectMapper = new ObjectMapper();

    @Transactional
    public void onOrderEvent(UUID eventId, String eventType, String key, String event, Instant ts) {
        if (log.alreadyProcessed(eventId)) {
            LOGGER.info("Event with UUID {} was already retrieved, ignoring it", eventId);
            return;
        }

        JsonNode eventPayload = deserialize(event);

        LOGGER.info("Received 'Order' event -- key: {}, event id: '{}', event type: '{}', ts: '{}'", key, eventId, eventType, ts);

        if (eventType.equals("OrderCreated")) {
            shipmentService.orderCreated(eventPayload);
        }
        else if (eventType.equals("OrderLineUpdated")) {
            shipmentService.orderLineUpdated(eventPayload);
        }
        else {
            LOGGER.warn("Unkown event type");
        }

        log.processed(eventId);
    }

    private JsonNode deserialize(String event) {
        JsonNode eventPayload;

        try {
            String unescaped = objectMapper.readValue(event, String.class);
            eventPayload = objectMapper.readTree(unescaped);
        }
        catch (IOException e) {
            throw new RuntimeException("Couldn't deserialize event", e);
        }

        return eventPayload.has("schema") ? eventPayload.get("payload") : eventPayload;
    }
}
````

### FILE: `debezium_postgres_inbox_consumer/upstream/ShipmentService.java`
```yaml
block_id: "DEBEZIUM-POSTGRES-INBOX-CONSUMER:upstream-shipmentservice-java:v1"
operation: CREATE
provenance: VERBATIM
source: "https://raw.githubusercontent.com/debezium/debezium-examples/7b0d765a02cf66ef29de2a01b67ec763e9acedf5/outbox/shipment-service/src/main/java/io/debezium/examples/outbox/shipment/service/ShipmentService.java"
license: "Apache-2.0"
sha256: "47e8fb65169a3a562b4572b27e3903d7c6f2f3e81dd42cf546d816126c896cf2"
variables: []
secrets_allowed: false
```
````java
/*
 * Copyright Debezium Authors.
 *
 * Licensed under the Apache Software License version 2.0, available at http://www.apache.org/licenses/LICENSE-2.0
 */
package io.debezium.examples.outbox.shipment.service;

import java.time.LocalDateTime;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.persistence.EntityManager;
import jakarta.persistence.PersistenceContext;
import jakarta.transaction.Transactional;
import jakarta.transaction.Transactional.TxType;
import jakarta.inject.Inject;

import io.opentelemetry.api.trace.Span;
import io.opentelemetry.api.trace.Tracer;
import io.opentelemetry.context.Scope;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.fasterxml.jackson.databind.JsonNode;

import io.debezium.examples.outbox.shipment.model.Shipment;

@ApplicationScoped
public class ShipmentService {

    private static final Logger LOGGER = LoggerFactory.getLogger(ShipmentService.class);

    @Inject
    Tracer tracer;

    @Inject
    Span span;

    @PersistenceContext
    EntityManager entityManager;

    @Transactional(value=TxType.MANDATORY)
    public void orderCreated(JsonNode event) {
        LOGGER.info("Processing 'OrderCreated' event: {}", event);

        span = tracer.spanBuilder("orderCreated").startSpan();
        try (final Scope scope = span.makeCurrent()){
            final long orderId = event.get("id").asLong();
            final long customerId = event.get("customerId").asLong();
            final LocalDateTime orderDate = LocalDateTime.parse(event.get("orderDate").asText());

            entityManager.persist(new Shipment(customerId, orderId, orderDate));
        } finally {
            span.end();
        }
    }

    @Transactional(value=TxType.MANDATORY)
    public void orderLineUpdated(JsonNode event) {
        LOGGER.info("Processing 'OrderLineUpdated' event: {}", event);
        span = tracer.spanBuilder("orderLineUpdated").startSpan();
        try (final Scope scope = span.makeCurrent()){
            span.addEvent("Processing 'OrderLineUpdated' event: "+event);
        } finally {
            span.end();
        }
    }
}
````

### FILE: `debezium_postgres_inbox_consumer/verify_contract.ps1`
```yaml
block_id: "DEBEZIUM-POSTGRES-INBOX-CONSUMER:verify-contract-ps1:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://debezium-postgres-inbox-consumer/verify_contract.ps1"
license: "LicenseRef-Workspace-Owner"
sha256: "20442bf14a07658682047f7deb47cd6012d3a7875fa498fa37ee11e2ac75d0b9"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding()]
param()
$ErrorActionPreference = 'Stop'
$root = $PSScriptRoot
$lock = Get-Content -LiteralPath (Join-Path $root 'source-lock.json') -Raw | ConvertFrom-Json -Depth 30
if ($lock.repository -ne 'debezium/debezium-examples' -or $lock.commit -ne '7b0d765a02cf66ef29de2a01b67ec763e9acedf5' -or $lock.archiveSha256 -ne '0fba61c33bdbdb8f32eeab790335dcec412826dd98fa0bbab282afa2130f51e4') { throw 'Source identity mismatch' }
if ($lock.license -ne 'Apache-2.0' -or $lock.officialBuild.testsRun -ne 0 -or $lock.officialBuild.result -ne 'BUILD_SUCCESS_CONDITIONED') { throw 'Official sample condition was weakened' }
if ($lock.adaptationRuntime.quarkusBom -ne '3.39.1' -or $lock.adaptationRuntime.jacksonDatabindOverride -ne '2.22.1' -or $lock.adaptationRuntime.celJava -ne 'dev.cel:cel:0.13.0' -or $lock.adaptationRuntime.components -ne 170 -or $lock.adaptationRuntime.osvFindings -ne 0) { throw 'Adaptation runtime lock mismatch' }
foreach ($file in $lock.upstreamFiles) {
  $path = Join-Path $root ([string]$file.path)
  if (-not (Test-Path -LiteralPath $path -PathType Leaf)) { throw "Missing upstream file $($file.path)" }
  $actual = (Get-FileHash -LiteralPath $path -Algorithm SHA256).Hash.ToLowerInvariant()
  if ($actual -ne [string]$file.sha256) { throw "Upstream hash mismatch $($file.path)" }
}
$profile = Get-Content -LiteralPath (Join-Path $root 'project-profile.template.json') -Raw | ConvertFrom-Json -Depth 20
if ($profile.acknowledgeConditionedState -ne $false -or @($profile.owners).Count -ne 0 -or @($profile.approvals).Count -ne 0) { throw 'Profile must start closed' }
foreach ($proof in $profile.proofs.psobject.Properties) { if ($proof.Value -ne $false) { throw "Proof must start false: $($proof.Name)" } }
$consumer = Get-Content -LiteralPath (Join-Path $root 'src/main/java/com/elite/documentinbox/KafkaDocumentIntakeConsumer.java') -Raw
$processor = Get-Content -LiteralPath (Join-Path $root 'src/main/java/com/elite/documentinbox/JdbcDocumentIntakeProcessor.java') -Raw
$mapping = Get-Content -LiteralPath (Join-Path $root 'src/main/java/com/elite/documentinbox/CelDocumentProjectionMapping.java') -Raw
$validator = Get-Content -LiteralPath (Join-Path $root 'src/main/java/com/elite/documentinbox/DocumentIntakeEnvelopeValidator.java') -Raw
$migration = Get-Content -LiteralPath (Join-Path $root 'src/main/resources/db/migration/V001__document_consumer_inbox.sql') -Raw
$properties = Get-Content -LiteralPath (Join-Path $root 'src/main/resources/application.properties') -Raw
$pom = Get-Content -LiteralPath (Join-Path $root 'pom.xml') -Raw
foreach ($marker in @('@Acknowledgment(Acknowledgment.Strategy.MANUAL)','@Blocking(ordered = true)','processor.process(event)','return message.ack()')) { if (-not $consumer.Contains($marker)) { throw "Consumer contract missing $marker" } }
foreach ($marker in @('connection.setAutoCommit(false)','insert into document_intelligence.consumer_inbox','on conflict do nothing','applyProjection(connection, event)','mapAuthoritativeDocument(connection, event)','mapping.apply(event.domainType(), documentJson, documentSha256)','mapped_payload','connection.commit()','connection.rollback()','Divergent duplicate event')) { if (-not $processor.Contains($marker)) { throw "Transaction contract missing $marker" } }
foreach ($marker in @('.setStandardMacros()','maxExpressionCodePointSize(4096)','maxParseRecursionDepth(32)','Mapping expression must be single-line','Mapping expression hash mismatch','Event domain type does not match approved mapping','Mapping result keys do not match approved schema','inputSchemaSha256','outputSchemaSha256')) { if (-not $mapping.Contains($marker)) { throw "Mapping contract missing $marker" } }
foreach ($marker in @('STRICT_DUPLICATE_DETECTION','Header must occur exactly once','Unsupported event type','Payload schema is not closed','MAX_PAYLOAD_BYTES')) { if (-not $validator.Contains($marker)) { throw "Validation contract missing $marker" } }
foreach ($marker in @('enable row level security','elite_document_consumer nologin','primary key (consumer_name, tenant_id, event_id)','unique (consumer_name, topic, partition_id, offset_id)','foreign key (tenant_id, evaluation_id) references document_intelligence.intake_document','mapping_expression_sha256_hex','mapped_payload_sha256_hex','mapped_payload jsonb')) { if (-not $migration.Contains($marker)) { throw "Migration contract missing $marker" } }
foreach ($marker in @('enable.auto.commit=false','isolation.level=read_committed','failure-strategy=fail','security.protocol=SSL')) { if (-not $properties.Contains($marker)) { throw "Kafka config missing $marker" } }
foreach ($marker in @('<quarkus.platform.version>3.39.1</quarkus.platform.version>','<jackson-databind.version>2.22.1</jackson-databind.version>','<cel.version>0.13.0</cel.version>')) { if (-not $pom.Contains($marker)) { throw "POM runtime lock missing $marker" } }
$astTokens = $null; $astErrors = $null
[System.Management.Automation.Language.Parser]::ParseFile((Join-Path $root 'run_postgres_tests.ps1'), [ref]$astTokens, [ref]$astErrors) | Out-Null
if (@($astErrors).Count -ne 0) { throw 'PostgreSQL runner does not parse' }
$astTokens = $null; $astErrors = $null
[System.Management.Automation.Language.Parser]::ParseFile((Join-Path $root 'run_sca.ps1'), [ref]$astTokens, [ref]$astErrors) | Out-Null
if (@($astErrors).Count -ne 0) { throw 'SCA runner does not parse' }
Write-Output 'DEBEZIUM_POSTGRES_INBOX_CONSUMER_STATIC_PASS upstream=6 cel_mapping=1 conditioned=1'
````

## 6. Configuration surface

Perfil inicia cerrado. Debe fijar topic/group/bootstrap, secrets JDBC/TLS, mapping ID/clase/version/expresión/hash, output keys, schemas/hashes, owners, approvals y diez pruebas live. Cada configuración sirve una sola clase documental exacta; ningún secreto ni mapping de negocio se inventa en Markdown.

## 7. Dependency bill

| Dependencia | Pin | Uso | Licencia | Gate |
|---|---|---|---|---|
| Red Hat/Debezium examples | commit 7b0d765, archive SHA-256 0fba61… | patrón oficial/source evidence | Apache-2.0 | seis hashes exactos |
| Quarkus BOM | 3.39.1 | runtime Kafka/JDBC/CDI | Apache-2.0 | Maven Central + build |
| Jackson Databind | 2.22.1 | JSON cerrado | Apache-2.0 | override oficial corregido |
| CEL-Java | 0.13.0 | mapping type-checked sin macros | Apache-2.0 | JAR hash exacto + 6 pruebas mapping |
| PostgreSQL | 18.6 | transacción/RLS/concurrency | PostgreSQL | 18 tests reales |
| OSV Scanner | 2.5.1 | SCA SBOM | Apache-2.0 | 170/0 fechado |

## 8. Apply order

1. Materializar y ejecutar static/build/PostgreSQL/SCA. 2. Fijar schemas, mapping, corpus, hashes y aprobaciones. 3. Aplicar foundation, boundary 0.2 y migración inbox como owner. 4. Crear login no-owner miembro de elite_document_consumer. 5. Completar profile y secrets externos. 6. Desplegar DEV con Kafka TLS. 7. Probar duplicate/mapping-change/race/poison/restart/redelivery/order/load/outage/recovery/rollback.

## 9. Verification

PASS offline requiere 24/24 hashes, cinco fuentes Java upstream exactas más licencia normalizada sólo por LF, compilación Java 21, 18/18 tests sobre PostgreSQL 18.6, mapping y receipt dentro de la transacción, cuatro carreras de 1 APPLIED/7 DUPLICATE, rollback ante mapping inválido o clase cruzada, reuse de offset rechazado, profile cerrado y SBOM 170 componentes/0 findings OSV. Live permanece condicionado.

## 10. Reconstruction evidence

V97 conserva commit/archive/licencia/build oficial 6 fuentes/0 tests, debilidades UP-FAIL-182 y upgrade Quarkus/Jackson. V99 añade CEL-Java 0.13.0 oficial, clase documental exacta, mapping transaccional, 24 archivos, 18 tests PostgreSQL y SCA 170/0; schemas/corpus/aprobación y live continúan condicionados.
