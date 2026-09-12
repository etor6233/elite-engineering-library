# Official structured logging audit — V281

Registro de admisión: USE_CONDITIONED_PACK / ready=false; alcance estrecho, no cierre de observabilidad.

Fecha: 2026-09-07. Clasificación: mantenimiento / investigación T2809.
Procedencia de este expediente y probes: AUTHORED. Código probado: distribución oficial Go, sin cambios.
No es un pack nuevo, una actualización del perfil de producto ni el cierre de T2809.

## Registro ejecutable del candidato

El resolver materializado emitió exit 0: `CAPABILITY_GAP_RESOLUTION_PASS capability=OBSERVABILITY-STRUCTURED-LOGGING state=USE_CONDITIONED_PACK ready=false`.
Este PASS valida el expediente condicionado; explícitamente no habilita implementación. Los FAIL G4–G8 significan evidencia obligatoria de admisión no completada, no cinco defectos nuevos upstream.
Para reproducir, colocar los dos manuales autoridad, este expediente, la distribución exacta y su LICENSE en un staging con los paths del registro; usar el resolver de CAPABILITY_GAP_RESOLUTION_GATE.

```json
{
  "schema_version": 1,
  "capability_id": "OBSERVABILITY-STRUCTURED-LOGGING",
  "requirement_ref": "OFFICIAL_STRUCTURED_LOGGING_AUDIT_V281.md#claim-y-decision-acotados",
  "observed_at": "2026-09-07T14:28:57Z",
  "max_age_days": 30,
  "authority_refs": [
    "PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD.md",
    "SECURITY_SRE_CLOUD_INFRASTRUCTURE.md"
  ],
  "source_searches": [
    {
      "search_id": "official-go-and-exporter-comparison",
      "query": "Go slog exact runtime 1.26.8 source license tests; Google and Microsoft telemetry scope",
      "official_domains": [
        "go.dev",
        "go.googlesource.com",
        "github.com",
        "microsoft.github.io",
        "docs.cloud.google.com",
        "opentelemetry.io"
      ],
      "executed_at": "2026-09-07T14:28:57Z",
      "evidence_refs": [
        "OFFICIAL_STRUCTURED_LOGGING_AUDIT_V281.md"
      ]
    }
  ],
  "candidates": [
    {
      "candidate_id": "GO-SLOG-1.26.8",
      "origin_kind": "OFFICIAL_CODE",
      "provenance": "VERBATIM",
      "source_url": "https://go.dev/dl/go1.26.8.windows-amd64.zip",
      "claim": "Structured Logger/Handler/Writer primitive, explicit LogValuer and ReplaceAttr; tested locally without a product integration.",
      "non_claims": [
        "Universal PII detection",
        "Full observability or monotonic metrics",
        "Production integration or toolchain promotion"
      ],
      "status": "CONDITIONED",
      "immutable_revision": "c293dd49cbe25e1fe8d97d94a5cb618e7b6d831e",
      "artifact_ref": "go1.26.8.windows-amd64.zip",
      "artifact_sha256": "b92c3b2adae85a11ba71fe7216daf0d84e82af4c8ab6c5625807f28622043a59",
      "license_expression": "BSD-3-Clause",
      "license_evidence_ref": "current-toolchain/go/LICENSE",
      "gates": {
        "G0": "PASS",
        "G1": "PASS",
        "G2": "PASS",
        "G3": "PASS",
        "G4": "FAIL",
        "G5": "FAIL",
        "G6": "FAIL",
        "G7": "FAIL",
        "G8": "FAIL"
      }
    }
  ],
  "decision": {
    "state": "USE_CONDITIONED_PACK",
    "selected_candidate_id": "GO-SLOG-1.26.8",
    "reason": "Conditioned reference for the narrow logging primitive only; full required observability remains unresolved. FAIL means mandatory admission evidence not complete, not an assertion of an upstream vulnerability.",
    "canonical_updates": [
      "OFFICIAL_STRUCTURED_LOGGING_AUDIT_V281.md"
    ],
    "blockers": [
      "T2801 maintenance readiness and implementation assurance",
      "Product graph SCA, contract, event/attribute policy and connected journey",
      "Metrics, sink health, load, supervision, retention, recovery and rollback",
      "No replacement pack admitted; GO-OBSERVABILITY-CORE remains CANDIDATE"
    ],
    "exhaustion": null
  }
}
```

## Resultado demostrado

- Go 1.26.7 íntegro: 278 casos oficiales PASS, 0 FAIL, 0 SKIP, dos paquetes.
- Go 1.26.8 íntegro: los mismos 278 casos PASS, 0 FAIL, 0 SKIP.
- Probes propios de alcance: 11 casos PASS en cada runtime. Incluyen subtests; no son 11 suites independientes.
- Los 37 archivos de src/log/slog son idénticos entre ambas distribuciones.
- En el toolchain temporal anterior, 33 de esos archivos son idénticos y faltan cuatro archivos de benchmarks; no se encontraron diferencias de bytes en los presentes.
- La primera ejecución sobre el toolchain recortado falló antes de ejecutar tests: faltan internal/testenv y testing/slogtest. No es un defecto de upstream. FAIL-20260907-388 conserva el fallo.
- Ningún runtime, dependency lock, perfil integral ni código de producto fue actualizado. Los Go completos viven sólo en staging aislado.

## Identidad y fuentes oficiales

| Artefacto | Identidad exacta / SHA-256 |
|---|---|
| go1.26.7.windows-amd64.zip | f4f534a486e4bc3387fa18f08208f2f854b7aaea8a08f2a2d829a914a05abb11 |
| go1.26.7 commit | e3336a22ad3f0a90bd252c95d8b5544e02674205 |
| go1.26.8.windows-amd64.zip | b92c3b2adae85a11ba71fe7216daf0d84e82af4c8ab6c5625807f28622043a59 |
| go1.26.8 commit | c293dd49cbe25e1fe8d97d94a5cb618e7b6d831e |
| LICENSE leído de la distribución 1.26.7 | 911f8f5782931320f5b8d1160a76365b83aea6447ee6c04fa6d5591467db9dad |

Descargas y digests: [API oficial](https://go.dev/dl/?mode=json&include=all).
Revisión actual: [Go 1.26.8 en Git at Google](https://go.googlesource.com/go/+/refs/tags/go1.26.8).
Licencia: BSD-3-Clause, Copyright Go Authors; conservar licencia/atribución y no usar Google LLC como aval del derivado.
[LICENSE fijada](https://github.com/golang/go/blob/c293dd49cbe25e1fe8d97d94a5cb618e7b6d831e/LICENSE).

La API oficial devolvió go1.27.1 y go1.26.8 como descargas estables actuales.
Se probó 1.26.8 como candidato patch de la rama ya usada; no se migró de major.
El [historial oficial](https://go.dev/doc/devel/release) describe 1.26.8 como correcciones de cgo, compilador, runtime, debug/elf y os.
Esto no reemplaza un análisis de vulnerabilidades del artefacto y grafo de producto ni autoriza actualizarlo por antigüedad solamente.

## Comparación: no confundir responsabilidades

| Fuente oficial | Qué aporta | Por qué no cierra por sí sola el problema |
|---|---|---|
| Go log/slog | logger estructurado, handlers, grupos, ReplaceAttr y LogValuer | no detector universal de datos sensibles, métricas, retención, alertas ni delivery durable |
| GoogleCloudPlatform/opentelemetry-operations-go | exportadores a servicios de Google Cloud | cuenta, IAM, destino y costo del target; no política completa de datos permitidos |
| Microsoft waza / OTEL | tracing de evaluaciones y control explícito de exportación de payload | no es el host operativo de la franquicia; fallo de exporter es best-effort |
| OpenTelemetry security guidance | minimización y procesamiento de datos sensibles | guía y componentes requieren configuración y pruebas; no garantia de exactitud universal |

Fuentes revisadas: [Google instrumentación](https://docs.cloud.google.com/stackdriver/docs/instrumentation/choose-approach),
[exportadores Google](https://github.com/GoogleCloudPlatform/opentelemetry-operations-go),
[Microsoft waza](https://microsoft.github.io/waza/guides/otel/),
[OpenTelemetry sensitive data](https://opentelemetry.io/docs/security/handling-sensitive-data/).
Google/Microsoft no fueron copiados, fijados ni compilados en esta auditoría: comparación documental, no admisión de sus repositorios.
OpenTelemetry no se reatribuye a Google o Microsoft.

## Claim y decisión acotados

Candidato Go: CONDITIONED para referencia de logging estructurado.
No es REUSABLE_PACK ni reemplazo admitido de GO-OBSERVABILITY-CORE.

Autoridad interna: SECURITY_SRE_CLOUD_INFRASTRUCTURE.md §15; PUBLIC_CODE_ARCHITECTURE_ADMISSION_STANDARD.md G0–G8.
Alineación: separación de señales, atributos explícitos, minimización previa a exportación y salud del sink.
No se eleva AUTHORED por usar una API de Go.

| Gate | Estado de evidencia |
|---|---|
| G0 identidad | PASS para archivos/distribuciones exactos |
| G1 licencia | PASS para estudio y prueba local del código BSD-3-Clause; producto conserva notices |
| G2 alineación | PASS sólo para primitiva estructurada y límites declarados |
| G3 arquitectura | PASS para Logger → Handler → Writer, no para arquitectura de servicio completa |
| G4 correctitud | suites focales PASS; grafo/toolchain/contrato de producto todavía no admitidos |
| G5 seguridad | pendiente política de eventos/atributos, trust boundaries, SCA y pruebas del journey |
| G6 rendimiento/resiliencia | pendiente carga, cardinalidad, sink lento/caído, recursos y comportamiento del host |
| G7 operación | pendiente alertas, retención/acceso, recuperación, supervisión y rollback |
| G8 pack integrado | pendiente implementación y evidencia conectada; no se compone automáticamente |

## Límites reproducidos

1. Sin política adicional, mensaje, atributo password, mapa anidado y error conservan el marcador sensible sintético.
2. LogValuer protege el valor cuando se usa como atributo directo o dentro de slog.Group; no transforma un valor encapsulado en un mapa opaco serializado por JSON.
3. ReplaceAttr visita el atributo password de un grupo estructurado, no las claves interiores de un mapa arbitrario.
4. Logger.Info no devuelve el error del Writer: detectar pérdida de exportación requiere un owner explícito.
5. Nada de esto reproduce ni corrige el contador local decreciente. FAIL-386/387 siguen abiertos y el pack permanece CANDIDATE.

Estos casos PASS caracterizan limitaciones; no significan que una filtración esté aprobada.
El ejemplo oficial ExampleLogValuer_secret, ejecutado en la suite, sigue siendo ejemplo, no wrapper productivo.

## Reproducción y evidencia

Staging local: elite-v281-30cf0957089d499db7bf9b2d1669b7db bajo el temporal del sistema.
No se distribuyen binarios, datos de usuario ni aprobaciones del mantenimiento.

En una distribución oficial completa, verificada por digest:
```text
go test -json -count=1 log/slog testing/slogtest
```
Para los probes, crear módulo aislado elite.local/slog-audit con go 1.26.7 y ejecutar go test -json -count=1 ./...
El siguiente bloque conserva exactamente el archivo de prueba local, sin afirmar procedencia upstream.

| Log / archivo | SHA-256 |
|---|---|
| slog-official-tests.jsonl (primer intento FAIL) | f24038c85f35d7d44315a2874dd5804e64e2916247889d49488750603afe95d3 |
| slog-complete-tests.jsonl | cdeda85998b4afb1c20f15e4eb8ea2c5198eea0f05f74273e4b67028806e4e65 |
| slog-current-tests.jsonl | e00a3c0768b04a37bfed5ba30347c774a708e8be5b91e6e4616fa623511ae2f5 |
| slog-scope-tests.jsonl | 814e2d03e9699be084ef3d86e2b0bfb07f07c6f4e4614caca416d479b2736f9e |
| slog-current-scope-tests.jsonl | d5ba97c476a3d3b55d8937cd0f8b962d5b5e1ffc1c2265372c499019ae8d6d27 |
| slog_scope_test.go | 5d672b6d4dd9df562c500111adc7829d814ba28edc8914a342ddb5dafb2adc09 |

```go
// AUTHORED audit probes, not upstream code or a production logging adapter.
package slogaudit

import (
	"bytes"
	"errors"
	"log/slog"
	"strings"
	"testing"
)

type secret string
func (secret) LogValue() slog.Value { return slog.StringValue("REDACTED") }

func TestDefaultDoesNotDetectSensitiveData(t *testing.T) {
	for _, location := range []string{"message", "attribute", "map", "error"} {
		t.Run(location, func(t *testing.T) {
			var output bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&output, nil))
			const marker = "synthetic-sensitive-marker"
			switch location {
			case "message": logger.Info(marker)
			case "attribute": logger.Info("event", "password", marker)
			case "map": logger.Info("event", "metadata", map[string]string{"password": marker})
			case "error": logger.Info("event", "error", errors.New(marker))
			}
			if !strings.Contains(output.String(), marker) { t.Fatal("baseline behavior changed; re-audit privacy assumptions") }
		})
	}
}

func TestLogValuerScope(t *testing.T) {
	for _, location := range []string{"attribute", "group", "opaque_map"} {
		t.Run(location, func(t *testing.T) {
			var output bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&output, nil))
			const marker = "synthetic-token-marker"
			switch location {
			case "attribute": logger.Info("event", "token", secret(marker))
			case "group": logger.Info("event", slog.Group("auth", "token", secret(marker)))
			case "opaque_map": logger.Info("event", "auth", map[string]any{"token": secret(marker)})
			}
			leaked := strings.Contains(output.String(), marker)
			if leaked != (location == "opaque_map") { t.Fatalf("unexpected LogValuer scope: %s", output.String()) }
			if !leaked && !strings.Contains(output.String(), "REDACTED") { t.Fatal("missing explicit redaction") }
		})
	}
}

func TestReplaceAttrDoesNotWalkOpaqueMap(t *testing.T) {
	var output bytes.Buffer
	visited := 0
	logger := slog.New(slog.NewJSONHandler(&output, &slog.HandlerOptions{ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == "password" { visited++; return slog.String("password", "REDACTED") }
		return a
	}}))
	logger.Info("event", slog.Group("structured", "password", "group-marker"), "opaque", map[string]string{"password": "map-marker"})
	if visited != 1 || strings.Contains(output.String(), "group-marker") || !strings.Contains(output.String(), "map-marker") {
		t.Fatalf("unexpected traversal: visited=%d output=%s", visited, output.String())
	}
}

type failedSink struct { calls int }
func (s *failedSink) Write(p []byte) (int, error) { s.calls++; return 0, errors.New("synthetic sink unavailable") }

func TestLoggerDoesNotPropagateSinkFailure(t *testing.T) {
	sink := new(failedSink)
	logger := slog.New(slog.NewJSONHandler(sink, nil))
	logger.Info("event") // No returned error: delivery health needs its own owner.
	if sink.calls != 1 { t.Fatalf("sink calls = %d", sink.calls) }
}
```
