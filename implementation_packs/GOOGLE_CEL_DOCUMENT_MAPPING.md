# Google CEL Document Mapping

## 1. Metadata

```yaml
pack_id: "GOOGLE-CEL-DOCUMENT-MAPPING"
pack_version: "0.1.0"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CONDITIONED
claim: "Materializa un mapper documental configurable sobre CEL-Java oficial: expresión single-line versionada/hash-bound, type-check previo, macros y funciones host deshabilitadas, output cerrado y receipt de identidad."
stacks: ["Google-origin CEL-Java 0.13.0", "Java 21", "Maven 3.9.16", "CycloneDX 2.9.1", "OSV Scanner 2.5.1", "PowerShell 7"]
compatible_with: ["DEBEZIUM-POSTGRES-INBOX-CONSUMER 0.1.x", "STRICT-DOCUMENT-FIELD-EVALUATION-GATE 0.1.x"]
incompatible_with: ["multiline CEL", "macros", "host functions", "floating point", "unversioned mapping", "unhashed mapping", "open output keys", "automatic business persistence"]
license_expression: "LicenseRef-Workspace-Owner AND Apache-2.0"
upstream_sources: ["https://github.com/cel-expr/cel-java/tree/c57d9d02259795c761e82f8870303b93a7a672c8", "https://repo.maven.apache.org/maven2/dev/cel/cel/0.13.0/"]
verified_at: "2026-08-28"
```

## 2. Applicability

Use para convertir el JSON documental ya asegurado/evaluado en un objeto de staging configurable por negocio sin ejecutar JavaScript ni código arbitrario. No sustituye schema input/output, ground truth, revisión, reglas legales ni autorización de persistencia.

## 3. Architecture contract

Depende del artefacto oficial CEL-Java, pero los siete archivos son integración local `AUTHORED`, no código Google. La expresión se compila una vez con tipo `map<string,dyn>` de entrada/salida, sin macros ni funciones host, máximo 4.096 code points y profundidad 32. Debe ser single-line y coincidir con SHA-256/version/ID aprobados. Input/output rechazan null, floating point, ciclos, profundidad mayor a 16, colecciones mayores a 10.000 y valores no JSON; la salida exige el conjunto exacto de keys. El resultado conserva mapping ID/version/hash. El caller debe validar schemas hash-locked antes/después y no puede persistir sólo porque CEL evaluó.

## 4. Exact file manifest

```text
CREATE google_cel_document_mapping/README.md
CREATE google_cel_document_mapping/pom.xml
CREATE google_cel_document_mapping/project-profile.template.json
CREATE google_cel_document_mapping/source-lock.json
CREATE google_cel_document_mapping/src/main/java/com/elite/mapping/CelDocumentMapping.java
CREATE google_cel_document_mapping/src/test/java/com/elite/mapping/CelDocumentMappingTest.java
CREATE google_cel_document_mapping/verify.ps1
```

## 5. Materialization blocks

### FILE: `google_cel_document_mapping/README.md`
```yaml
block_id: "GOOGLE-CEL-DOCUMENT-MAPPING:readme-md:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://google-cel-document-mapping/README.md"
license: "LicenseRef-Workspace-Owner"
sha256: "f917de509c8c4d535b680b49ec171df6c3d651d1c55c664385f7dc81c024d714"
variables: []
secrets_allowed: false
```
````markdown
# Google-origin CEL document mapping

This component compiles an approved, hash-bound CEL expression once and evaluates it without host functions or macros. It accepts one JSON-like `document` map and returns only the exact approved output keys together with mapping identity, semantic version and expression SHA-256.

The integration is local `AUTHORED` code depending on the official Apache-2.0 `dev.cel:cel:0.13.0` artifact. It is not Google source. The official release has five Windows-only newline test failures recorded in `source-lock.json`; this wrapper rejects every multiline expression before CEL compilation. It also limits expression size and parse depth, rejects nulls, floating-point values, cycles, oversized collections and non-JSON values, and disables every macro.

This mapper does not prove extraction accuracy, schema correctness or business persistence. Before an output can reach the Debezium boundary, the project must validate the exact input and output schemas with their locked hashes, evaluate its own ground-truth corpus and approve the mapping profile. Keep `acknowledgeConditionedState=false` until every project proof is real.
````

### FILE: `google_cel_document_mapping/pom.xml`
```yaml
block_id: "GOOGLE-CEL-DOCUMENT-MAPPING:pom-xml:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://google-cel-document-mapping/pom.xml"
license: "LicenseRef-Workspace-Owner"
sha256: "2f287b0439789190b60b7cc1ab450ae79033e73cd42d7f52bbb0c34a910df12b"
variables: []
secrets_allowed: false
```
````xml
<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 https://maven.apache.org/xsd/maven-4.0.0.xsd">
  <modelVersion>4.0.0</modelVersion>
  <groupId>com.elite.mapping</groupId>
  <artifactId>google-cel-document-mapping</artifactId>
  <version>0.1.0</version>
  <properties>
    <maven.compiler.release>21</maven.compiler.release>
    <project.build.sourceEncoding>UTF-8</project.build.sourceEncoding>
    <cel.version>0.13.0</cel.version>
    <junit.version>5.13.4</junit.version>
  </properties>
  <dependencies>
    <dependency><groupId>dev.cel</groupId><artifactId>cel</artifactId><version>${cel.version}</version></dependency>
    <dependency><groupId>org.junit.jupiter</groupId><artifactId>junit-jupiter</artifactId><version>${junit.version}</version><scope>test</scope></dependency>
  </dependencies>
  <build><plugins>
    <plugin><groupId>org.apache.maven.plugins</groupId><artifactId>maven-compiler-plugin</artifactId><version>3.15.0</version></plugin>
    <plugin><groupId>org.apache.maven.plugins</groupId><artifactId>maven-surefire-plugin</artifactId><version>3.5.5</version></plugin>
    <plugin><groupId>org.cyclonedx</groupId><artifactId>cyclonedx-maven-plugin</artifactId><version>2.9.1</version></plugin>
  </plugins></build>
</project>
````

### FILE: `google_cel_document_mapping/project-profile.template.json`
```yaml
block_id: "GOOGLE-CEL-DOCUMENT-MAPPING:project-profile-template-json:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://google-cel-document-mapping/project-profile.template.json"
license: "LicenseRef-Workspace-Owner"
sha256: "004946e779cbc42df9640f65a128a4aa6496adeb8729c9675604918d7766d69b"
variables: []
secrets_allowed: false
```
````json
{
  "schemaVersion": 1,
  "acknowledgeConditionedState": false,
  "mappingId": "",
  "mappingVersion": "",
  "expression": "",
  "expressionSha256": "",
  "exactOutputKeys": [],
  "inputSchemaId": "",
  "inputSchemaSha256": "",
  "outputSchemaId": "",
  "outputSchemaSha256": "",
  "owners": [],
  "approvals": [],
  "proofs": {
    "inputSchemaValidated": false,
    "outputSchemaValidated": false,
    "mappingCorpusEvaluated": false,
    "duplicateDeterminismProven": false,
    "loadAndResourceLimitsProven": false,
    "rollbackProven": false
  }
}
````

### FILE: `google_cel_document_mapping/source-lock.json`
```yaml
block_id: "GOOGLE-CEL-DOCUMENT-MAPPING:source-lock-json:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://google-cel-document-mapping/source-lock.json"
license: "LicenseRef-Workspace-Owner"
sha256: "8907da4545163e09fd36bde44ceb6cc4d3c9c9dd9df8fec25790e0de54c14a53"
variables: []
secrets_allowed: false
```
````json
{
  "schemaVersion": 1,
  "upstream": {
    "project": "cel-expr/cel-java",
    "origin": "Google Common Expression Language",
    "release": "v0.13.0",
    "commit": "c57d9d02259795c761e82f8870303b93a7a672c8",
    "commitVerified": false,
    "commitVerificationReason": "unsigned",
    "license": "Apache-2.0",
    "archiveBytes": 1892506,
    "archiveSha256": "088d90e7e178fe09fd664d5dcb5c60a7d946fbe9d27cadb7400f9868e69b2cc8"
  },
  "maven": {
    "coordinate": "dev.cel:cel:0.13.0",
    "jarBytes": 499701,
    "jarSha256": "a3bc22d78affe749f81b17cce2131a3807275436157ba4d6fe3f10d8f541ea94",
    "pomBytes": 4259,
    "pomSha256": "3a6eccb1ce1aececfec7e90f56073e87f1bc14491b6d00104bf84beeedaed0f6"
  },
  "datedEvidence": {
    "verifiedAt": "2026-08-28",
    "upstreamFocalTests": 2836,
    "upstreamFailures": 5,
    "failureScope": "Windows LF versus System.lineSeparator in ErrorsTest",
    "sbomComponents": 18,
    "knownVulnerabilityFindings": 0
  }
}
````

### FILE: `google_cel_document_mapping/src/main/java/com/elite/mapping/CelDocumentMapping.java`
```yaml
block_id: "GOOGLE-CEL-DOCUMENT-MAPPING:cel-document-mapping-java:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://google-cel-document-mapping/src/main/java/com/elite/mapping/CelDocumentMapping.java"
license: "LicenseRef-Workspace-Owner"
sha256: "6d4b423e34ba999596cc7c0d2d92fe1b1b274d0c52908adf3033e2a37f26cf43"
variables: []
secrets_allowed: false
```
````java
package com.elite.mapping;

import dev.cel.common.CelAbstractSyntaxTree;
import dev.cel.common.CelOptions;
import dev.cel.common.CelValidationException;
import dev.cel.common.types.MapType;
import dev.cel.common.types.SimpleType;
import dev.cel.compiler.CelCompiler;
import dev.cel.compiler.CelCompilerFactory;
import dev.cel.runtime.CelEvaluationException;
import dev.cel.runtime.CelRuntime;
import dev.cel.runtime.CelRuntimeFactory;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.Set;
import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.HexFormat;
import java.util.IdentityHashMap;
import java.util.List;

public final class CelDocumentMapping {
  private static final CelOptions OPTIONS =
      CelOptions.current().maxExpressionCodePointSize(4096).maxParseRecursionDepth(32).build();
  private static final MapType STRING_DYN_MAP = MapType.create(SimpleType.STRING, SimpleType.DYN);

  private final CelRuntime.Program program;
  private final Set<String> exactOutputKeys;
  private final String mappingId;
  private final String mappingVersion;
  private final String expressionSha256;

  public CelDocumentMapping(String mappingId, String mappingVersion, String expression,
      String expressionSha256, Set<String> exactOutputKeys)
      throws CelValidationException, CelEvaluationException {
    if (expression == null || expression.isBlank()) {
      throw new IllegalArgumentException("mapping expression is required");
    }
    if (exactOutputKeys == null || exactOutputKeys.isEmpty()) {
      throw new IllegalArgumentException("exact output keys are required");
    }
    if (mappingId == null || !mappingId.matches("[a-z][a-z0-9-]{2,63}")) {
      throw new IllegalArgumentException("canonical mapping id is required");
    }
    if (mappingVersion == null || !mappingVersion.matches("[1-9][0-9]*\\.[0-9]+\\.[0-9]+")) {
      throw new IllegalArgumentException("semantic mapping version is required");
    }
    if (expression.indexOf('\r') >= 0 || expression.indexOf('\n') >= 0) {
      throw new IllegalArgumentException("mapping expression must be single-line");
    }
    String actualHash = sha256(expression);
    if (!actualHash.equals(expressionSha256)) {
      throw new IllegalArgumentException("mapping expression hash mismatch");
    }
    this.mappingId = mappingId;
    this.mappingVersion = mappingVersion;
    this.expressionSha256 = actualHash;
    this.exactOutputKeys = Set.copyOf(exactOutputKeys);
    CelCompiler compiler =
        CelCompilerFactory.standardCelCompilerBuilder()
            .setOptions(OPTIONS)
            .setStandardMacros()
            .addVar("document", STRING_DYN_MAP)
            .setResultType(STRING_DYN_MAP)
            .build();
    CelAbstractSyntaxTree ast = compiler.compile(expression).getAst();
    CelRuntime runtime = CelRuntimeFactory.standardCelRuntimeBuilder().setOptions(OPTIONS).build();
    this.program = runtime.createProgram(ast);
  }

  public MappingResult apply(Map<String, Object> document) throws CelEvaluationException {
    if (document == null) {
      throw new IllegalArgumentException("document is required");
    }
    validateJsonValue(document, 0, new IdentityHashMap<>());
    Object result = program.eval(Map.of("document", Map.copyOf(document)));
    if (!(result instanceof Map<?, ?> raw)) {
      throw new IllegalStateException("mapping result is not an object");
    }
    LinkedHashMap<String, Object> output = new LinkedHashMap<>();
    for (Map.Entry<?, ?> entry : raw.entrySet()) {
      if (!(entry.getKey() instanceof String key) || entry.getValue() == null) {
        throw new IllegalStateException("mapping result contains invalid key or null value");
      }
      validateJsonValue(entry.getValue(), 0, new IdentityHashMap<>());
      output.put(key, entry.getValue());
    }
    if (!output.keySet().equals(exactOutputKeys)) {
      throw new IllegalStateException("mapping result keys do not match approved schema");
    }
    return new MappingResult(mappingId, mappingVersion, expressionSha256, Map.copyOf(output));
  }

  private static String sha256(String value) {
    try {
      return HexFormat.of().formatHex(
          MessageDigest.getInstance("SHA-256").digest(value.getBytes(StandardCharsets.UTF_8)));
    } catch (NoSuchAlgorithmException impossible) {
      throw new IllegalStateException("SHA-256 unavailable", impossible);
    }
  }

  private static void validateJsonValue(
      Object value, int depth, IdentityHashMap<Object, Boolean> active) {
    if (depth > 16) {
      throw new IllegalArgumentException("document value exceeds maximum depth");
    }
    if (value instanceof String text) {
      if (text.length() > 1_048_576) {
        throw new IllegalArgumentException("document string exceeds maximum length");
      }
      return;
    }
    if (value instanceof Boolean || value instanceof Byte || value instanceof Short
        || value instanceof Integer || value instanceof Long) {
      return;
    }
    if (value instanceof Map<?, ?> map) {
      if (map.size() > 10_000 || active.put(map, Boolean.TRUE) != null) {
        throw new IllegalArgumentException("document map is oversized or cyclic");
      }
      try {
        for (Map.Entry<?, ?> entry : map.entrySet()) {
          if (!(entry.getKey() instanceof String key) || key.isEmpty() || key.length() > 256
              || entry.getValue() == null) {
            throw new IllegalArgumentException("document map key/value is invalid");
          }
          validateJsonValue(entry.getValue(), depth + 1, active);
        }
      } finally {
        active.remove(map);
      }
      return;
    }
    if (value instanceof List<?> list) {
      if (list.size() > 10_000 || active.put(list, Boolean.TRUE) != null) {
        throw new IllegalArgumentException("document list is oversized or cyclic");
      }
      try {
        for (Object item : list) {
          if (item == null) {
            throw new IllegalArgumentException("document list contains null");
          }
          validateJsonValue(item, depth + 1, active);
        }
      } finally {
        active.remove(list);
      }
      return;
    }
    throw new IllegalArgumentException("document contains a non-JSON or imprecise numeric value");
  }

  public record MappingResult(
      String mappingId, String mappingVersion, String expressionSha256, Map<String, Object> values) {}
}
````

### FILE: `google_cel_document_mapping/src/test/java/com/elite/mapping/CelDocumentMappingTest.java`
```yaml
block_id: "GOOGLE-CEL-DOCUMENT-MAPPING:cel-document-mapping-test-java:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://google-cel-document-mapping/src/test/java/com/elite/mapping/CelDocumentMappingTest.java"
license: "LicenseRef-Workspace-Owner"
sha256: "3e31c39660b23105318a65b0a6f9d7e6680527b91f0be995f992538446eab415"
variables: []
secrets_allowed: false
```
````java
package com.elite.mapping;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

import dev.cel.common.CelValidationException;
import dev.cel.runtime.CelEvaluationException;
import java.util.Map;
import java.util.Set;
import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.util.HexFormat;
import java.util.HashMap;
import org.junit.jupiter.api.Test;

final class CelDocumentMappingTest {
  @Test
  void mapsApprovedFields() throws Exception {
    String expression = "{'supplier_id': document['supplier']['id'], 'total_minor': document['total_minor']}";
    CelDocumentMapping mapping = new CelDocumentMapping(
        "supplier-invoice", "1.0.0", expression, sha256(expression),
        Set.of("supplier_id", "total_minor"));
    CelDocumentMapping.MappingResult result = mapping.apply(Map.of(
        "supplier", Map.of("id", "SUP-42"),
        "total_minor", 129900L,
        "unapproved", "not copied"));
    assertEquals(Map.of("supplier_id", "SUP-42", "total_minor", 129900L), result.values());
    assertEquals("supplier-invoice", result.mappingId());
    assertEquals("1.0.0", result.mappingVersion());
  }

  @Test
  void rejectsUnapprovedOutputKey() throws Exception {
    String expression = "{'supplier_id': document['supplier_id'], 'secret': document['secret']}";
    CelDocumentMapping mapping = new CelDocumentMapping(
        "supplier-invoice", "1.0.0", expression, sha256(expression),
        Set.of("supplier_id"));
    assertThrows(IllegalStateException.class,
        () -> mapping.apply(Map.of("supplier_id", "SUP-42", "secret", "x")));
  }

  @Test
  void missingInputFailsClosed() throws Exception {
    String expression = "{'supplier_id': document['supplier_id']}";
    CelDocumentMapping mapping = new CelDocumentMapping(
        "supplier-invoice", "1.0.0", expression, sha256(expression), Set.of("supplier_id"));
    assertThrows(CelEvaluationException.class, () -> mapping.apply(Map.of()));
  }

  @Test
  void macrosAreNotAvailable() {
    assertThrows(CelValidationException.class,
        () -> { String e = "{'items': document['items'].map(x, x)}"; new CelDocumentMapping(
            "supplier-invoice", "1.0.0", e, sha256(e), Set.of("items")); });
  }

  @Test
  void oversizedExpressionIsRejected() {
    String expression = " ".repeat(4097) + "{}";
    assertThrows(CelValidationException.class,
        () -> new CelDocumentMapping("supplier-invoice", "1.0.0", expression,
            sha256(expression), Set.of("x")));
  }

  @Test
  void hashAndMultilineAreRejectedBeforeCompilation() {
    assertThrows(IllegalArgumentException.class,
        () -> new CelDocumentMapping("supplier-invoice", "1.0.0", "{}", "0".repeat(64), Set.of("x")));
    String multiline = "{'x': 1}\n";
    assertThrows(IllegalArgumentException.class,
        () -> new CelDocumentMapping("supplier-invoice", "1.0.0", multiline,
            sha256(multiline), Set.of("x")));
  }

  @Test
  void impreciseNumericAndCyclicInputAreRejected() throws Exception {
    String expression = "{'value': document['value']}";
    CelDocumentMapping mapping = new CelDocumentMapping(
        "supplier-invoice", "1.0.0", expression, sha256(expression), Set.of("value"));
    assertThrows(IllegalArgumentException.class, () -> mapping.apply(Map.of("value", 1.5d)));
    Map<String, Object> cyclic = new HashMap<>();
    cyclic.put("value", cyclic);
    assertThrows(IllegalArgumentException.class, () -> mapping.apply(cyclic));
  }

  private static String sha256(String value) throws Exception {
    return HexFormat.of().formatHex(
        MessageDigest.getInstance("SHA-256").digest(value.getBytes(StandardCharsets.UTF_8)));
  }
}
````

### FILE: `google_cel_document_mapping/verify.ps1`
```yaml
block_id: "GOOGLE-CEL-DOCUMENT-MAPPING:verify-ps1:v1"
operation: CREATE
provenance: AUTHORED
source: "workspace-owner://google-cel-document-mapping/verify.ps1"
license: "LicenseRef-Workspace-Owner"
sha256: "0b295f7c521baec508ad8cc732fc9ae5f620397d5687295e811c5d006bbbd9de"
variables: []
secrets_allowed: false
```
````powershell
#requires -Version 7.0
[CmdletBinding()]
param(
  [Parameter(Mandatory)][string]$JavaHome,
  [Parameter(Mandatory)][string]$MavenExecutable,
  [Parameter(Mandatory)][string]$OsvScanner,
  [switch]$AllowNetwork
)
$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $MyInvocation.MyCommand.Path
$java = Join-Path $JavaHome 'bin/java.exe'
if (-not (Test-Path -LiteralPath $java -PathType Leaf)) { throw 'Java executable missing' }
if ((& $java -version 2>&1 | Out-String) -notmatch 'version "21\.') { throw 'Java 21 required' }
if ((& $MavenExecutable -version | Out-String) -notmatch 'Apache Maven 3\.9\.16') { throw 'Maven 3.9.16 required' }
if ((& $OsvScanner --version | Out-String) -notmatch 'osv-scanner version: 2\.5\.1') { throw 'OSV Scanner 2.5.1 required' }
if ((Get-FileHash -LiteralPath $OsvScanner -Algorithm SHA256).Hash.ToLowerInvariant() -ne '25e42f5ef6711fd8c0fb45390972205891dd44c6bd02ac93f0f63e8e98d9bfb6') { throw 'OSV Scanner hash mismatch' }
$env:JAVA_HOME = $JavaHome
$env:Path = (Join-Path $JavaHome 'bin') + [IO.Path]::PathSeparator + $env:Path
& $MavenExecutable -B -ntp -f (Join-Path $root 'pom.xml') -o clean test
if ($LASTEXITCODE -ne 0) { throw "Offline Maven tests failed: $LASTEXITCODE" }
if (-not $AllowNetwork) { throw 'AllowNetwork is required to refresh the CycloneDX/OSV evidence' }
& $MavenExecutable -B -ntp -f (Join-Path $root 'pom.xml') org.cyclonedx:cyclonedx-maven-plugin:2.9.1:makeAggregateBom
if ($LASTEXITCODE -ne 0) { throw "CycloneDX generation failed: $LASTEXITCODE" }
$report = Join-Path ([IO.Path]::GetTempPath()) ('elite-cel-mapping-osv-' + [guid]::NewGuid().ToString('N') + '.json')
try {
  & $OsvScanner scan source --sbom (Join-Path $root 'target/bom.json') --format json --output-file $report
  if ($LASTEXITCODE -ne 0) { throw "OSV findings or scan failure: $LASTEXITCODE" }
  if (-not (Test-Path -LiteralPath $report -PathType Leaf)) { throw 'OSV report missing' }
  $sbom = Get-Content -LiteralPath (Join-Path $root 'target/bom.json') -Raw | ConvertFrom-Json -Depth 100
  if (@($sbom.components).Count -ne 18) { throw 'Unexpected runtime component count' }
  'GOOGLE_CEL_DOCUMENT_MAPPING_PASS tests=7 sbom=18 osv=0'
} finally {
  if (Test-Path -LiteralPath $report -PathType Leaf) { [IO.File]::Delete($report) }
}
````

## 6. Configuration surface

Completar `project-profile.template.json`: mapping ID/version/expression/SHA, keys exactas, IDs+hashes de schemas, owners/approvals y pruebas. `acknowledgeConditionedState` permanece false hasta evidencia real. No admite secretos en expresión o perfil; runtime secrets pertenecen al componente que consume el mapping.

## 7. Dependency bill

| Dependencia | Pin | Uso | Licencia | Evidencia |
|---|---:|---|---|---|
| CEL-Java | 0.13.0 | compiler/runtime no Turing-completo | Apache-2.0 | JAR/POM/archive hashes; 2.836 upstream focales, 5 Windows failures visibles |
| Java | 21 | compilación/runtime | GPL-2.0-with-classpath-exception según distribución | runner exige major 21 |
| Maven | 3.9.16 | build/tests/SBOM | Apache-2.0 | runner exige versión exacta |
| CycloneDX Maven | 2.9.1 | SBOM 1.6 | Apache-2.0 | 18 componentes |
| Google OSV Scanner | 2.5.1 | SCA | Apache-2.0 | binario/hash/versión exactos; 18/0 fechado |

## 8. Apply order

1. Materializar. 2. Ejecutar tests offline y SCA con red autorizada. 3. Definir schemas input/output cerrados y sus hashes. 4. Escribir mapping single-line y calcular SHA. 5. Evaluar corpus/ground truth, duplicados y límites. 6. Aprobar profile. 7. Conectar después de evaluación documental y antes de staging/outbox. 8. Mantener persistencia bloqueada hasta que el consumer valide todos los receipts.

## 9. Verification

`verify.ps1` exige Java 21, Maven 3.9.16 y OSV 2.5.1 exacto. Tests Maven se ejecutan offline; `-AllowNetwork` es obligatorio para regenerar CycloneDX/OSV. El gate fechado exige 7/7 tests, 18 componentes y cero findings. Además ejecutar compositor/global y round-trip byte a byte. Live target debe probar schemas, corpus, determinismo, carga, rollback y unión transaccional.

## 10. Reconstruction evidence

Evidencia gobernante: `reconstruction_evidence/CEL_JAVA_DOCUMENT_MAPPING_CANDIDATE_2026-08-28_V98.md`, a elevar a V98 final sólo después de materialización canónica, runner desde round-trip, perfil y gates globales. Fallos `UP-FAIL-185` y `LIB-FAIL-1082..1088` permanecen visibles.
