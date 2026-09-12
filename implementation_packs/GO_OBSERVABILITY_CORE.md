# Go Observability Core

## 1. Metadata

```yaml
pack_id: "GO-OBSERVABILITY-CORE"
pack_version: "0.3.3"
status:
  authority: SUPPORTED_REFERENCE
  implementation: REBUILD_VERIFIED
  admission: CANDIDATE
claim: "Candidato local: logger de política cerrada, contador monotónico e histograma validado con overflow separado y snapshot coherente. No detector universal de PII, SDK OTel ni observabilidad integral."
stacks: ["Go 1.26.8"]
compatible_with: ["módulo aislado Go 1.26.7; integración con owners empresariales no admitida"]
incompatible_with: ["mensaje libre", "valores arbitrarios/anidados", "política no revisada", "migración automática desde logger 0.1.x"]
license_expression: "LicenseRef-Workspace-Owner"
upstream_sources: ["https://go.dev/ref/spec#Integer_overflow", "https://pkg.go.dev/log/slog#JSONHandler.Handle", "https://developer.android.com/privacy-and-security/risks/log-info-disclosure", "https://opentelemetry.io/docs/security/handling-sensitive-data/", "https://opentelemetry.io/docs/specs/otel/metrics/api/", "https://opentelemetry.io/docs/specs/otel/error-handling/"]
verified_at: "2026-09-10"
```

## Reapertura de admisión V280 — 2026-09-07

No incorporar como protección de PII. El snapshot 0.1.0 está preservado en la
evidencia V280/V291 para reproducir V219: sus cinco tests pasan, pero cuatro
pruebas adversariales V280 fallan (email, secreto anidado, mensaje libre y delta
negativo). `REBUILD_VERIFIED` describe esa reconstrucción histórica, no admisión
vigente. `CANDIDATE` bloquea composición incluso con acknowledgeConditions=true.
El claim original debe leerse como intención refutada en esas entradas.
Antes de readmitir: alternativa oficial compatible o corrección declarada,
pruebas de privacidad/monotonicidad/límites, integración y gates afectados.
Owner del expediente: `reconstruction_evidence/CONNECTED_COVERAGE_AUDIT_V280.md`.

V291 / 0.1.1 corrige únicamente el contador: TryInc rechaza delta negativo y
overflow sin mutación/panic y comunica aceptación; Inc conserva firma y devuelve
valor sin cambios ante rechazo. No es implementación de la API OTel ni código
copiado de Google. Privacidad FAIL-386 y admisión CANDIDATE siguen abiertas; las
otras clases no se readmiten por este fix. Evidencia: COUNTER_MONOTONIC_REPAIR_V291.md.

## 2. Applicability

V342 / 0.3.1 validates the count returned by the configured JSON writer before
reporting successful delivery. Zero/short/negative/oversized counts with no error
now produce ErrDeliveryUnknown through the existing static error boundary.
The wrapper performs exactly one write; it never retries a partial record.
A full write is not proof of flush, fsync, remote acknowledgement or retention.
A misbehaving writer violates io.Writer's contract; this local defensive guard is
not a patch to Go. Existing privacy, blocking and re-entry constraints remain.
The isolated reference test connects real HTTP requests to policy-approved JSON
in an actual file, with correlation/counters/histogram and a closed-file failure.
It is a test harness, not production middleware or automatic runtime wiring.
CANDIDATE remains in force. See OBSERVABILITY_DELIVERY_INTEGRATION_V342.md.


V301 / 0.3.0 corrige el histograma, todavía CANDIDATE y fuera del perfil integral.
NewCheckedHistogram valida1..1024 límites finitos no negativos estrictamente
crecientes. TryObserve rechaza negativos, NaN, infinitos y overflow sin mutación.
Snapshot copia límites/counts/suma/total bajo el mismo lock; counts no acumulados.
El último bucket ahora es separado y su límite superior es +Inf, no el último
límite finito. Percentile devuelve un límite de bucket, no percentil exacto;
p0 selecciona el primer bucket ocupado, entradas/config inválidas retornan NaN.
Vacío válido conserva0, por lo que se debe consultar Count antes de interpretarlo.
NewHistogram conserva firma pero una configuración inválida queda inerte;
Observe no informa rechazo: migrar a los métodos comprobables antes de adoptar.
Cambios de semántica explícitos: no upgrade silencioso ni serializar NaN/+Inf
como JSON numérico; un exportador debe representar esos estados según su contrato.
Rango exacto del índice usa math/big oficial Go sólo al consultar; observación
lineal acotada. No SLO/overhead ni exporter demostrado. Ver
reconstruction_evidence/HISTOGRAM_BOUNDARY_REPAIR_V301.md.
Fuentes: Google Cloud Monitoring Distribution (población finita y overflow) y
OpenTelemetry data-model (buckets superiores inclusivos) y metrics/api (no negativos).
Google usa inclusividad opuesta: esto no es su formato ni un adapter compatible.
Código, límite1024 y contrato local AUTHORED; no se copia código de esas páginas.

V292 / 0.2.0 reemplaza la redacción heurística por política explícita de eventos y
valores. NewJSONLogger usa Go log/slog sin copiar su implementación; la política,
validación y tests siguen AUTHORED. NewLogger sin política no emite; Redact retorna
cero campos. Es un cambio incompatible deliberado: migrar a NewJSONLogger o
NewPolicyLogger y comprobar TryEmit, nunca activar una actualización silenciosa.
La evidencia V292 prueba rechazo de las tres exposiciones V280 y eventos útiles
reales serializados, no redacción universal ni aprobación del despliegue.

Sólo investigación/reparación aislada mientras siga CANDIDATE. No usar como
protección de PII ni componer en producto. Para diagnosticar rechazos del contador
usar TryInc y comprobar su booleano; Inc conserva compatibilidad pero no informa
el rechazo por separado. No usar esta métrica como ledger financiero o stock.

## 3. Architecture contract

- **Ownership**: `internal/observability` emite/agrega; la exportación (OTel OTLP/Prometheus) es composición del runtime.
- **Invariantes**: contador no negativo/no decreciente; logger sólo admite eventos y campos exactos, enum cerrado/bool/int64 acotado. Rechazo completo previo al sink, sin formatear input. Histograma no negativo/finito, suma y count sin overflow, count igual a suma de buckets; privacidad universal no demostrada.
- **Data flow**: política revisada → TryEmit → snapshot de primitivas → Go JSONHandler.Handle → writer → resultado explícito; Counter.TryInc; Histogram.Observe/Percentile.
- **Failure modes**: ErrRejected/ErrFiltered/ErrUnconfigured no escriben; ErrDeliveryUnknown no incluye error/panic original ni autoriza retry automático. Sink síncrono puede bloquear, reentrada no admitida, entrega no transaccional.
- **Seguridad/privacidad**: política es configuración confiable; nombres/enums/rangos mal aprobados pueden divulgar datos. No almacenar chats/documentos crudos en logs ni confundirlos con archivo original. Aplicación debe revisar terceros, exportación, acceso y retención.
- **Performance budget**: límites locales 128 eventos, 32 campos/evento, 64 valores/enum, identificadores ASCII de 64 bytes; son decisiones AUTHORED. Emisión O(campos × enum), orden JSON O(campos log campos); sin SLO/benchmark validado.

## 4. Exact file manifest

```text
CREATE internal/observability/observability.go
CREATE internal/observability/observability_test.go
```

## 5. Materialization blocks

### FILE: `internal/observability/observability.go`
```yaml
block_id: "GO-OBSERVABILITY-CORE:internal/observability/observability.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "0425ec30fa22dbec4f579274281e2668e12638eb7252d74eaf6772a0c801b2c5"
variables: []
secrets_allowed: false
```
````go
// Package observability provides CANDIDATE instruments, not general PII detection.
// The policy logger only emits reviewed event IDs and closed, typed field values.
// Policy configuration is trusted. Integration and operational admission remain pending.
// All policy glue is AUTHORED; JSON encoding uses the official Go log/slog handler.
package observability

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"math"
	"math/big"
	"sort"
	"sync"
	"time"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

var (
	ErrPolicy          = errors.New("observability: invalid policy")
	ErrRejected        = errors.New("observability: event rejected")
	ErrUnconfigured    = errors.New("observability: policy required")
	ErrFiltered        = errors.New("observability: level filtered")
	ErrDeliveryUnknown = errors.New("observability: delivery unknown")
)

type FieldKind uint8

const (
	EnumField FieldKind = iota + 1
	BoolField
	Int64Field
)

// FieldRule accepts only exact built-in types. Enum values and numeric ranges
// must be reviewed for privacy by the application owner; this is not a PII oracle.
type FieldRule struct {
	Kind     FieldKind
	Values   []string
	Min, Max int64
}

// Policy maps approved event IDs to exact, required fields. No dynamic message,
// optional extra fields, arbitrary strings, nested objects or callbacks are allowed.
type Policy map[string]map[string]FieldRule

type Entry struct {
	Level  Level
	Msg    string
	Fields map[string]any
	At     time.Time
}

type Sink func(Entry)

type Logger struct {
	mu     sync.Mutex
	level  Level
	policy Policy
	sink   func(Entry) error
}

// Deprecated: without an explicit policy, redaction cannot prove a safe record.
// Redact now releases no fields. It never traverses or formats untrusted values.
func Redact(fields map[string]any) map[string]any { return map[string]any{} }

// Deprecated: use NewJSONLogger or NewPolicyLogger and check TryEmit's error.
// A logger constructed without a reviewed policy rejects every event.
func NewLogger(level Level, sink Sink) *Logger { return &Logger{level: level} }

func identifier(s string) bool {
	if len(s) == 0 || len(s) > 64 {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		letter := c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
		if !letter && (i == 0 || !(c >= '0' && c <= '9' || c == '_' || c == '.' || c == '-')) {
			return false
		}
	}
	return true
}

func copyPolicy(level Level, p Policy) (Policy, error) {
	if level < LevelDebug || level > LevelError || len(p) == 0 || len(p) > 128 {
		return nil, ErrPolicy
	}
	out := make(Policy, len(p))
	for event, fields := range p {
		if !identifier(event) || len(fields) > 32 {
			return nil, ErrPolicy
		}
		rules := make(map[string]FieldRule, len(fields))
		for key, rule := range fields {
			if !identifier(key) || key == "time" || key == "level" || key == "msg" || key == "source" {
				return nil, ErrPolicy
			}
			switch rule.Kind {
			case EnumField:
				if len(rule.Values) == 0 || len(rule.Values) > 64 || rule.Min != 0 || rule.Max != 0 {
					return nil, ErrPolicy
				}
				seen := make(map[string]bool, len(rule.Values))
				for _, value := range rule.Values {
					if !identifier(value) || seen[value] {
						return nil, ErrPolicy
					}
					seen[value] = true
				}
				rule.Values = append([]string(nil), rule.Values...)
			case BoolField:
				if len(rule.Values) != 0 || rule.Min != 0 || rule.Max != 0 {
					return nil, ErrPolicy
				}
			case Int64Field:
				if len(rule.Values) != 0 || rule.Min > rule.Max {
					return nil, ErrPolicy
				}
			default:
				return nil, ErrPolicy
			}
			rules[key] = rule
		}
		out[event] = rules
	}
	return out, nil
}

func NewPolicyLogger(level Level, sink Sink, policy Policy) (*Logger, error) {
	if sink == nil {
		return nil, ErrPolicy
	}
	p, err := copyPolicy(level, policy)
	if err != nil {
		return nil, err
	}
	return &Logger{level: level, policy: p, sink: func(e Entry) error { sink(e); return nil }}, nil
}

// NewJSONLogger uses Go's JSONHandler, with no source or arbitrary context attrs.
// TryEmit calls Handle directly so write failures are observable, unlike Logger.Info.
// The caller owns writer lifecycle. A typed-nil writer is invalid caller input;
// any resulting writer panic is contained as ambiguous delivery.
func NewJSONLogger(level Level, writer io.Writer, policy Policy) (*Logger, error) {
	if writer == nil {
		return nil, ErrPolicy
	}
	p, err := copyPolicy(level, policy)
	if err != nil {
		return nil, err
	}
	handler := slog.NewJSONHandler(checkedJSONWriter{writer: writer}, &slog.HandlerOptions{Level: slog.LevelDebug})
	return &Logger{level: level, policy: p, sink: func(e Entry) error {
		record := slog.NewRecord(e.At, slog.Level(int(e.Level)*4-4), e.Msg, 0)
		keys := make([]string, 0, len(e.Fields))
		for key := range e.Fields {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			switch value := e.Fields[key].(type) {
			case string:
				record.AddAttrs(slog.String(key, value))
			case bool:
				record.AddAttrs(slog.Bool(key, value))
			case int64:
				record.AddAttrs(slog.Int64(key, value))
			}
		}
		return handler.Handle(context.Background(), record)
	}}, nil
}

// checkedJSONWriter detects sinks that violate io.Writer's count/error contract.
// It makes one call only: partial writes are ambiguous and must not be retried.
// A full write is not a claim of flush, fsync, remote acknowledgement or retention.
type checkedJSONWriter struct{ writer io.Writer }

func (w checkedJSONWriter) Write(p []byte) (int, error) {
	n, err := w.writer.Write(p)
	if n < 0 || n > len(p) {
		return 0, io.ErrShortWrite
	}
	if n != len(p) && err == nil {
		return n, io.ErrShortWrite
	}
	return n, err
}

func accepts(rule FieldRule, value any) bool {
	switch rule.Kind {
	case EnumField:
		s, ok := value.(string)
		if !ok {
			return false
		}
		for _, allowed := range rule.Values {
			if s == allowed {
				return true
			}
		}
	case BoolField:
		_, ok := value.(bool)
		return ok
	case Int64Field:
		n, ok := value.(int64)
		return ok && n >= rule.Min && n <= rule.Max
	}
	return false
}

// TryEmit snapshots approved primitives and rejects whole invalid events before
// calling a sink. The caller must not concurrently mutate fields. Sink delivery
// is serialized, synchronous and NOT transactional: it may block or partially
// write. ErrDeliveryUnknown must not trigger blind retry (duplicates possible).
// Sinks must not re-enter this logger. Errors never contain rejected input.
func (l *Logger) TryEmit(level Level, event string, fields map[string]any) (err error) {
	if l == nil || l.policy == nil || l.sink == nil {
		return ErrUnconfigured
	}
	if level < LevelDebug || level > LevelError {
		return ErrRejected
	}
	if level < l.level {
		return ErrFiltered
	}
	rules, ok := l.policy[event]
	if !ok || len(fields) != len(rules) {
		return ErrRejected
	}
	approved := make(map[string]any, len(rules))
	for key, rule := range rules {
		value, exists := fields[key]
		if !exists || !accepts(rule, value) {
			return ErrRejected
		}
		approved[key] = value
	}
	e := Entry{Level: level, Msg: event, Fields: approved, At: time.Now().UTC()}
	l.mu.Lock()
	defer l.mu.Unlock()
	defer func() {
		if recover() != nil {
			err = ErrDeliveryUnknown
		}
	}()
	if l.sink(e) != nil {
		return ErrDeliveryUnknown
	}
	return nil
}

// Deprecated: these compatibility methods discard rejection/delivery errors.
// Operational consumers must migrate to TryEmit and handle its status explicitly.
func (l *Logger) Emit(level Level, msg string, fields map[string]any) {
	_ = l.TryEmit(level, msg, fields)
}
func (l *Logger) Debug(msg string, fields map[string]any) { l.Emit(LevelDebug, msg, fields) }
func (l *Logger) Info(msg string, fields map[string]any)  { l.Emit(LevelInfo, msg, fields) }
func (l *Logger) Warn(msg string, fields map[string]any)  { l.Emit(LevelWarn, msg, fields) }
func (l *Logger) Error(msg string, fields map[string]any) { l.Emit(LevelError, msg, fields) }

// Counter is a monotonic counter.
type Counter struct {
	mu sync.Mutex
	v  int64
}

// Inc adds a non-negative delta and returns the current value. A negative delta
// or int64 overflow is a no-op. Use TryInc when rejection must be diagnosed.
func (c *Counter) Inc(delta int64) int64 {
	value, _ := c.TryInc(delta)
	return value
}

// TryInc atomically adds delta only when the result fits a non-negative int64.
// It returns the current value and whether the update was accepted. Invalid
// instrumentation input never panics, wraps, saturates, or partially increments.
// Counter must not be copied after first use (it contains a mutex).
func (c *Counter) TryInc(delta int64) (int64, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if delta < 0 || delta > math.MaxInt64-c.v {
		return c.v, false
	}
	c.v += delta
	return c.v, true
}

// Value returns the current count.
func (c *Counter) Value() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.v
}

// Histogram aggregates finite non-negative measurements in upper-inclusive buckets.
// It is a local CANDIDATE, not an OTel SDK. Do not copy it after first use.
type Histogram struct {
	mu         sync.Mutex
	bounds     []float64
	counts     []int64 // len(bounds)+1, including the unbounded overflow bucket
	sum        float64
	total      int64
	configured bool
}

// MaxHistogramBounds is a local resource limit, not an upstream requirement.
const MaxHistogramBounds = 1024

// NewCheckedHistogram validates and copies strictly ascending finite non-negative
// bounds. Units and SLO boundaries remain the caller's reviewed configuration.
func NewCheckedHistogram(bounds []float64) (*Histogram, error) {
	if len(bounds) == 0 || len(bounds) > MaxHistogramBounds {
		return nil, ErrRejected
	}
	for i, b := range bounds {
		if math.IsNaN(b) || math.IsInf(b, 0) || b < 0 || (i > 0 && b <= bounds[i-1]) {
			return nil, ErrRejected
		}
	}
	return &Histogram{bounds: append([]float64(nil), bounds...), counts: make([]int64, len(bounds)+1), configured: true}, nil
}

// NewHistogram preserves the legacy signature. Invalid configuration yields an
// inert instrument. Prefer NewCheckedHistogram to detect configuration errors.
func NewHistogram(bounds []float64) *Histogram {
	h, err := NewCheckedHistogram(bounds)
	if err != nil {
		return &Histogram{}
	}
	return h
}

// TryObserve records one measurement or rejects it without changing any aggregate.
// Non-finite/negative values, unconfigured instruments and numeric overflow fail.
func (h *Histogram) TryObserve(v float64) bool {
	if math.IsNaN(v) || math.IsInf(v, 0) || v < 0 {
		return false
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if !h.configured || h.total == math.MaxInt64 {
		return false
	}
	sum := h.sum + v
	if math.IsNaN(sum) || math.IsInf(sum, 0) {
		return false
	}
	bucket := len(h.bounds)
	for i, b := range h.bounds {
		if v <= b {
			bucket = i
			break
		}
	}
	h.counts[bucket]++
	h.total++
	h.sum = sum
	return true
}

// Observe preserves the signature; use TryObserve when rejected data must be
// accounted for. Rejection never panics, emits data or retries automatically.
func (h *Histogram) Observe(v float64) { h.TryObserve(v) }

func (h *Histogram) Count() int64 {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.total
}
func (h *Histogram) Sum() float64 {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.sum
}

// HistogramSnapshot is a coherent detached view. Counts are non-cumulative;
// the last count represents values above the last finite bound. It is not OTLP.
type HistogramSnapshot struct {
	Configured bool
	Bounds     []float64
	Counts     []int64
	Count      int64
	Sum        float64
}

func (h *Histogram) Snapshot() HistogramSnapshot {
	h.mu.Lock()
	defer h.mu.Unlock()
	return HistogramSnapshot{Configured: h.configured, Bounds: append([]float64(nil), h.bounds...), Counts: append([]int64(nil), h.counts...), Count: h.total, Sum: h.sum}
}

// Percentile returns the upper bound of the occupied nearest-rank bucket, NOT
// an exact quantile or an interpolated value. The overflow bucket returns +Inf.
// p=0 selects the first occupied bucket. Invalid p/configuration returns NaN;
// an empty valid instrument returns 0 for compatibility (inspect Count first).
// Rank uses the exact binary value of p, avoiding int64 overflow/float count loss.
// Rational work happens only on queries, never on the observation path.
func (h *Histogram) Percentile(p float64) float64 {
	if math.IsNaN(p) || math.IsInf(p, 0) || p < 0 || p > 100 {
		return math.NaN()
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if !h.configured {
		return math.NaN()
	}
	if h.total == 0 {
		return 0
	}
	r := new(big.Rat).SetFloat64(p)
	r.Mul(r, new(big.Rat).SetInt64(h.total))
	r.Quo(r, big.NewRat(100, 1))
	rank, rem := new(big.Int), new(big.Int)
	rank.QuoRem(r.Num(), r.Denom(), rem)
	if rem.Sign() > 0 {
		rank.Add(rank, big.NewInt(1))
	}
	if rank.Sign() == 0 {
		rank.SetInt64(1)
	}
	target := rank.Int64()
	var cumulative int64
	for i, n := range h.counts {
		cumulative += n
		if cumulative >= target {
			if i == len(h.bounds) {
				return math.Inf(1)
			}
			return h.bounds[i]
		}
	}
	return math.NaN()
}
````

### FILE: `internal/observability/observability_test.go`
```yaml
block_id: "GO-OBSERVABILITY-CORE:internal/observability/observability_test.go:v1"
operation: CREATE
provenance: AUTHORED
source: "local"
license: "LicenseRef-Workspace-Owner"
sha256: "0e817dbcd2c97d6622416bc83038464b7708b5512955bc4f7605637e745a9ae5"
variables: []
secrets_allowed: false
```
````go
package observability

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func fixturePolicy() Policy {
	return Policy{"job.done": {
		"outcome": {Kind: EnumField, Values: []string{"ok", "rejected"}},
		"retry":   {Kind: BoolField},
		"count":   {Kind: Int64Field, Min: 0, Max: 10},
	}, "heartbeat": {}}
}
func fixtureFields() map[string]any {
	return map[string]any{"outcome": "ok", "retry": false, "count": int64(1)}
}

func TestLegacyUnconfiguredDoesNotDisclose(t *testing.T) {
	calls := 0
	logger := NewLogger(LevelDebug, func(Entry) { calls++ })
	fields := map[string]any{"email": "synthetic@example.invalid", "nested": map[string]any{"token": "secret-marker"}}
	logger.Error("secret-marker", fields)
	if calls != 0 || logger.TryEmit(LevelError, "secret-marker", fields) != ErrUnconfigured {
		t.Fatal("legacy disclosed")
	}
	if len(Redact(fields)) != 0 || fields["email"] != "synthetic@example.invalid" {
		t.Fatal("legacy redaction")
	}
	var zero Logger
	if zero.TryEmit(LevelInfo, "heartbeat", nil) != ErrUnconfigured {
		t.Fatal("zero logger")
	}
	var absent *Logger
	if absent.TryEmit(LevelInfo, "heartbeat", nil) != ErrUnconfigured {
		t.Fatal("nil logger")
	}
}

func TestJSONApprovedEventEndToEnd(t *testing.T) {
	var buf bytes.Buffer
	logger, err := NewJSONLogger(LevelInfo, &buf, fixturePolicy())
	if err != nil {
		t.Fatal(err)
	}
	fields := fixtureFields()
	if err = logger.TryEmit(LevelWarn, "job.done", fields); err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err = json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if len(got) != 6 || got["msg"] != "job.done" || got["level"] != "WARN" || got["outcome"] != "ok" || got["retry"] != false || got["count"] != float64(1) || got["time"] == "" {
		t.Fatal(got)
	}
	if fields["count"] != int64(1) {
		t.Fatal("mutated input")
	}
	buf.Reset()
	if logger.TryEmit(LevelInfo, "heartbeat", nil) != nil || !strings.Contains(buf.String(), "heartbeat") {
		t.Fatal("empty approved schema failed")
	}
}

type hostileValue struct{ calls *int }

func (v hostileValue) String() string { *v.calls++; panic("must not format untrusted values") }

func TestRejectWholeEventWithoutFormatting(t *testing.T) {
	var buf bytes.Buffer
	logger, _ := NewJSONLogger(LevelDebug, &buf, fixturePolicy())
	calls := 0
	cycle := map[string]any{}
	cycle["self"] = cycle
	values := []any{
		"synthetic@example.invalid", "secret-marker", map[string]any{"token": "secret-marker"}, cycle,
		[]string{"ok"}, errors.New("secret-marker"), hostileValue{&calls}, int(1), nil,
	}
	for _, value := range values {
		fields := fixtureFields()
		fields["outcome"] = value
		if logger.TryEmit(LevelInfo, "job.done", fields) != ErrRejected {
			t.Fatal("accepted arbitrary value")
		}
	}
	for _, event := range []string{"secret-marker", "email synthetic@example.invalid", "job.done\nsecret"} {
		if logger.TryEmit(LevelInfo, event, fixtureFields()) != ErrRejected {
			t.Fatal("accepted message")
		}
	}
	for _, change := range []func(map[string]any){
		func(m map[string]any) { m["email"] = "synthetic@example.invalid" },
		func(m map[string]any) { delete(m, "outcome") },
		func(m map[string]any) { m["count"] = int64(-1) },
		func(m map[string]any) { m["count"] = int64(11) },
		func(m map[string]any) { m["count"] = float64(1) },
		func(m map[string]any) { m["retry"] = "false" },
	} {
		fields := fixtureFields()
		change(fields)
		if logger.TryEmit(LevelInfo, "job.done", fields) != ErrRejected {
			t.Fatal("accepted invalid fields")
		}
	}
	if buf.Len() != 0 || calls != 0 {
		t.Fatal("untrusted data touched sink/formatter")
	}
}

func TestPolicyValidation(t *testing.T) {
	cases := []Policy{
		nil, {}, {"": {}}, {"bad\nevent": {}}, {"event": {"msg": {Kind: BoolField}}},
		{"event": {"x": {Kind: 0}}}, {"event": {"x": {Kind: EnumField}}},
		{"event": {"x": {Kind: EnumField, Values: []string{"ok", "ok"}}}},
		{"event": {"x": {Kind: EnumField, Values: []string{"synthetic@example.invalid"}}}},
		{"event": {"x": {Kind: Int64Field, Min: 2, Max: 1}}},
		{"event": {"x": {Kind: BoolField, Values: []string{"true"}}}},
		{"event": {"x": {Kind: Int64Field, Values: []string{"one"}}}},
		{"event": {"x": {Kind: EnumField, Values: []string{"ok"}, Max: 1}}},
		{"event": {"x": {Kind: BoolField, Min: 1}}},
	}
	tooManyEvents := Policy{}
	for i := 0; i < 129; i++ {
		tooManyEvents[fmt.Sprintf("event%d", i)] = nil
	}
	cases = append(cases, tooManyEvents)
	tooManyFields := map[string]FieldRule{}
	for i := 0; i < 33; i++ {
		tooManyFields[fmt.Sprintf("field%d", i)] = FieldRule{Kind: BoolField}
	}
	cases = append(cases, Policy{"event": tooManyFields})
	values := []string{}
	for i := 0; i < 65; i++ {
		values = append(values, fmt.Sprintf("value%d", i))
	}
	cases = append(cases, Policy{"event": {"x": {Kind: EnumField, Values: values}}})
	for i, p := range cases {
		if _, err := NewJSONLogger(LevelInfo, &bytes.Buffer{}, p); err != ErrPolicy {
			t.Fatalf("case %d", i)
		}
	}
	for _, level := range []Level{-1, 4} {
		if _, err := NewJSONLogger(level, &bytes.Buffer{}, fixturePolicy()); err != ErrPolicy {
			t.Fatal("level")
		}
	}
	if _, err := NewJSONLogger(LevelInfo, nil, fixturePolicy()); err != ErrPolicy {
		t.Fatal("nil writer")
	}
	if _, err := NewPolicyLogger(LevelInfo, nil, fixturePolicy()); err != ErrPolicy {
		t.Fatal("nil sink")
	}
}

func TestPolicySnapshotAndSinkSnapshot(t *testing.T) {
	p := fixturePolicy()
	var seen Entry
	logger, err := NewPolicyLogger(LevelInfo, func(e Entry) { seen = e; e.Fields["count"] = int64(9) }, p)
	if err != nil {
		t.Fatal(err)
	}
	p["job.done"]["outcome"].Values[0] = "evil"
	delete(p, "heartbeat")
	fields := fixtureFields()
	if logger.TryEmit(LevelInfo, "job.done", fields) != nil || seen.Msg != "job.done" || fields["count"] != int64(1) {
		t.Fatal("snapshot failure")
	}
	fields["outcome"] = "evil"
	if logger.TryEmit(LevelInfo, "job.done", fields) != ErrRejected {
		t.Fatal("mutable policy")
	}
	if logger.TryEmit(LevelInfo, "heartbeat", nil) != nil {
		t.Fatal("event snapshot")
	}
}

type failingWriter struct {
	calls    int
	panicNow bool
}

func (w *failingWriter) Write(p []byte) (int, error) {
	w.calls++
	if w.panicNow {
		panic("secret-marker")
	}
	return 0, errors.New("secret-marker")
}
func TestDeliveryFailureContainedWithoutRetry(t *testing.T) {
	for _, panics := range []bool{false, true} {
		w := &failingWriter{panicNow: panics}
		logger, _ := NewJSONLogger(LevelInfo, w, fixturePolicy())
		for i := 1; i <= 2; i++ {
			if err := logger.TryEmit(LevelInfo, "heartbeat", nil); err != ErrDeliveryUnknown || strings.Contains(err.Error(), "secret-marker") {
				t.Fatal("failure disclosure")
			}
			if w.calls != i {
				t.Fatal("retry or lock not released")
			}
		}
	}
}

func TestLevelFiltering(t *testing.T) {
	calls := 0
	logger, _ := NewPolicyLogger(LevelWarn, func(Entry) { calls++ }, fixturePolicy())
	if logger.TryEmit(LevelDebug, "heartbeat", nil) != ErrFiltered || logger.TryEmit(Level(5), "heartbeat", nil) != ErrRejected {
		t.Fatal("level")
	}
	if logger.TryEmit(LevelError, "heartbeat", nil) != nil || calls != 1 {
		t.Fatal("accepted level")
	}
}
func TestConcurrentJSONRecords(t *testing.T) {
	var buf bytes.Buffer
	logger, _ := NewJSONLogger(LevelInfo, &buf, fixturePolicy())
	var wg sync.WaitGroup
	var failures atomic.Int32
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := 0; n < 10; n++ {
				if logger.TryEmit(LevelInfo, "job.done", fixtureFields()) != nil {
					failures.Add(1)
				}
			}
		}()
	}
	wg.Wait()
	lines := bytes.Split(bytes.TrimSpace(buf.Bytes()), []byte("\n"))
	if len(lines) != 320 || failures.Load() != 0 {
		t.Fatal("record loss")
	}
	for _, line := range lines {
		var record map[string]any
		if json.Unmarshal(line, &record) != nil || record["msg"] != "job.done" {
			t.Fatal("record interleaving")
		}
	}
}
func FuzzPolicyLoggerClosedValues(f *testing.F) {
	f.Add("job.done", "outcome", "ok")
	f.Add("job.done", "outcome", "secret-marker")
	f.Add("email@example.invalid", "outcome", "ok")
	f.Add("job.done", "nested", "secret-marker")
	f.Fuzz(func(t *testing.T, event, key, value string) {
		var buf bytes.Buffer
		logger, _ := NewJSONLogger(LevelInfo, &buf, Policy{"job.done": {"outcome": {Kind: EnumField, Values: []string{"ok", "rejected"}}}})
		err := logger.TryEmit(LevelInfo, event, map[string]any{key: value})
		allowed := event == "job.done" && key == "outcome" && (value == "ok" || value == "rejected")
		if !allowed {
			if err != ErrRejected || buf.Len() != 0 {
				t.Fatal("escaped closed policy")
			}
			return
		}
		var record map[string]any
		if err != nil || json.Unmarshal(buf.Bytes(), &record) != nil || len(record) != 4 || record["msg"] != event || record["outcome"] != value {
			t.Fatal("valid event lost")
		}
	})
}

func TestCounter(t *testing.T) {
	var c Counter
	if c.Inc(2) != 2 || c.Inc(3) != 5 || c.Value() != 5 {
		t.Fatalf("counter wrong: %d", c.Value())
	}
}

func TestHistogramPercentile(t *testing.T) {
	h := NewHistogram([]float64{10, 50, 100})
	for _, v := range []float64{5, 20, 60, 200} {
		h.Observe(v)
	}
	if h.Count() != 4 || h.Sum() != 285 {
		t.Fatalf("histogram count/sum wrong: %d/%v", h.Count(), h.Sum())
	}
	// p50 = 2nd observation → bucket [<=50] = 50.
	if h.Percentile(50) != 50 {
		t.Fatalf("p50 wrong: %v", h.Percentile(50))
	}
	// Values above100 occupy a separate unbounded bucket; no false ceiling.
	if !math.IsInf(h.Percentile(100), 1) {
		t.Fatalf("p100 wrong: %v", h.Percentile(100))
	}
}

func TestHistogramEmpty(t *testing.T) {
	h := NewHistogram([]float64{10})
	if h.Percentile(50) != 0 || h.Count() != 0 {
		t.Fatal("empty histogram should be zero")
	}
}

func TestCounterRejectNegative(t *testing.T) {
	var c Counter
	c.Inc(5)
	if got := c.Inc(-1); got != 5 || c.Value() != 5 {
		t.Fatalf("negative delta changed counter: %d", got)
	}
}

func TestCounterRejectOverflow(t *testing.T) {
	const max int64 = 1<<63 - 1
	var c Counter
	c.Inc(max)
	if got := c.Inc(1); got != max || c.Value() != max {
		t.Fatalf("overflow changed counter: %d", got)
	}
}

func TestCounterTryIncStatusAndRecovery(t *testing.T) {
	var c Counter
	for _, tc := range []struct {
		delta, want int64
		accepted    bool
	}{
		{0, 0, true}, {-1, 0, false}, {5, 5, true},
		{-1 << 63, 5, false}, {1<<63 - 1, 5, false}, {2, 7, true},
	} {
		got, accepted := c.TryInc(tc.delta)
		if got != tc.want || accepted != tc.accepted || c.Value() != tc.want {
			t.Fatalf("delta %d: value=%d accepted=%v", tc.delta, got, accepted)
		}
	}
}

func TestCounterConcurrentAcceptedAndRejected(t *testing.T) {
	var c Counter
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 256; j++ {
				if _, ok := c.TryInc(1); !ok {
					t.Error("valid update rejected")
				}
				if _, ok := c.TryInc(-1); ok {
					t.Error("negative update accepted")
				}
			}
		}()
	}
	wg.Wait()
	if got := c.Value(); got != 32*256 {
		t.Fatalf("lost increments: %d", got)
	}
}

func TestCounterConcurrentOverflow(t *testing.T) {
	const max int64 = 1<<63 - 1
	var c Counter
	c.Inc(max - 5)
	var accepted atomic.Int64
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, ok := c.TryInc(1); ok {
				accepted.Add(1)
			}
		}()
	}
	wg.Wait()
	if accepted.Load() != 5 || c.Value() != max {
		t.Fatal("overflow admission not atomic")
	}
}

func FuzzCounterMonotonic(f *testing.F) {
	const max int64 = 1<<63 - 1
	for _, pair := range [][2]int64{{0, 0}, {5, -1}, {max, 1}, {max - 1, 1}, {0, -1 << 63}, {0, max}, {max, max}} {
		f.Add(pair[0], pair[1])
	}
	f.Fuzz(func(t *testing.T, initial, delta int64) {
		initial &= max
		var c Counter
		c.Inc(initial)
		// Independent arbitrary-precision oracle, not the implementation's guard.
		sum := new(big.Int).Add(big.NewInt(initial), big.NewInt(delta))
		accepted := delta >= 0 && sum.IsInt64()
		want := initial
		if accepted {
			want = sum.Int64()
		}
		got, ok := c.TryInc(delta)
		if ok != accepted || got != want || c.Value() != want || got < initial {
			t.Fatalf("initial=%d delta=%d value=%d accepted=%v", initial, delta, got, ok)
		}
	})
}

func TestHistogramRegressionOverflow(t *testing.T) {
	h := NewHistogram([]float64{10, 100})
	h.Observe(200)
	if got := h.Percentile(100); !math.IsInf(got, 1) {
		t.Fatalf("overflow hidden: p100=%v for value200", got)
	}
}
func TestHistogramRegressionZeroPercentile(t *testing.T) {
	h := NewHistogram([]float64{10, 100})
	h.Observe(50)
	if got := h.Percentile(0); got != 100 {
		t.Fatalf("empty bucket selected: %v", got)
	}
}
func TestHistogramRegressionNonFinite(t *testing.T) {
	h := NewHistogram([]float64{10})
	h.Observe(2)
	h.Observe(math.NaN())
	if h.Count() != 1 || h.Sum() != 2 {
		t.Fatalf("NaN poisoned count/sum: %d/%v", h.Count(), h.Sum())
	}
}
func TestHistogramRegressionCounterOverflow(t *testing.T) {
	h := NewHistogram([]float64{10})
	h.total = math.MaxInt64
	h.counts[0] = math.MaxInt64
	h.Observe(1)
	if h.Count() != math.MaxInt64 {
		t.Fatalf("count wrapped: %d", h.Count())
	}
}

func TestHistogramConfigurationAndSnapshot(t *testing.T) {
	for _, bounds := range [][]float64{nil, {}, {2, 1}, {1, 1}, {-1}, {math.NaN()}, {math.Inf(1)}, {math.Inf(-1)}, make([]float64, MaxHistogramBounds+1)} {
		if h, err := NewCheckedHistogram(bounds); err != ErrRejected || h != nil {
			t.Fatalf("invalid configuration admitted")
		}
		legacy := NewHistogram(bounds)
		if legacy.TryObserve(1) || legacy.Count() != 0 || !math.IsNaN(legacy.Percentile(50)) {
			t.Fatal("invalid legacy configuration active")
		}
	}
	bounds := []float64{0, 10, 100}
	h, err := NewCheckedHistogram(bounds)
	if err != nil {
		t.Fatal(err)
	}
	bounds[0] = 999
	for _, v := range []float64{0, 1, 10, 11, 100, 101} {
		if !h.TryObserve(v) {
			t.Fatal("valid rejected")
		}
	}
	s := h.Snapshot()
	if !s.Configured || s.Count != 6 || s.Sum != 223 || !reflect.DeepEqual(s.Counts, []int64{1, 2, 2, 1}) || !reflect.DeepEqual(s.Bounds, []float64{0, 10, 100}) {
		t.Fatalf("snapshot=%+v", s)
	}
	s.Bounds[0] = 999
	s.Counts[0] = 999
	if h.Snapshot().Bounds[0] != 0 || h.Snapshot().Counts[0] != 1 {
		t.Fatal("snapshot aliases live state")
	}
	var zero Histogram
	if zero.TryObserve(0) || zero.Snapshot().Configured || !math.IsNaN(zero.Percentile(50)) {
		t.Fatal("zero value instrument active")
	}
	max := make([]float64, MaxHistogramBounds)
	for i := range max {
		max[i] = float64(i)
	}
	if _, err := NewCheckedHistogram(max); err != nil {
		t.Fatal("limit boundary rejected")
	}
}
func TestHistogramRejectAndRecovery(t *testing.T) {
	h := NewHistogram([]float64{10})
	h.Observe(2)
	before := h.Snapshot()
	for _, v := range []float64{-1, math.NaN(), math.Inf(1), math.Inf(-1)} {
		if h.TryObserve(v) || !reflect.DeepEqual(before, h.Snapshot()) {
			t.Fatal("rejection changed aggregate")
		}
	}
	for _, p := range []float64{-1, 101, math.NaN(), math.Inf(1), math.Inf(-1)} {
		if !math.IsNaN(h.Percentile(p)) {
			t.Fatal("invalid percentile masquerades as data")
		}
	}
	if !h.TryObserve(3) || h.Count() != 2 || h.Sum() != 5 {
		t.Fatal("recovery failed")
	}
	large := NewHistogram([]float64{10})
	large.Observe(math.MaxFloat64)
	before = large.Snapshot()
	if large.TryObserve(math.MaxFloat64) || !reflect.DeepEqual(before, large.Snapshot()) {
		t.Fatal("sum overflow was not atomic")
	}
	if !large.TryObserve(0) {
		t.Fatal("zero after rejected overflow failed")
	}
}
func TestHistogramPercentileLargeCount(t *testing.T) {
	h := NewHistogram([]float64{10, 100})
	h.total = math.MaxInt64
	h.counts[0] = math.MaxInt64 - 1
	h.counts[2] = 1
	if !math.IsInf(h.Percentile(100), 1) || h.Percentile(99) != 10 || h.Percentile(0) != 10 {
		t.Fatal("large count rank overflow")
	}
	if h.TryObserve(0) {
		t.Fatal("count overflow admitted")
	}
}
func TestHistogramConcurrentSnapshots(t *testing.T) {
	h := NewHistogram([]float64{1, 2})
	var wg sync.WaitGroup
	for n := 0; n < 16; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 128; i++ {
				if !h.TryObserve(2) || h.TryObserve(-1) {
					t.Error("observation admission")
				}
				s := h.Snapshot()
				if len(s.Counts) != 3 || s.Counts[0]+s.Counts[1]+s.Counts[2] != s.Count || s.Sum != float64(s.Count*2) {
					t.Error("torn snapshot")
				}
				if h.Percentile(100) != 2 {
					t.Error("concurrent percentile")
				}
			}
		}()
	}
	wg.Wait()
	if h.Count() != 2048 || h.Sum() != 4096 {
		t.Fatal("lost updates")
	}
}
func FuzzHistogramAggregate(f *testing.F) {
	for _, s := range [][]byte{{}, {0}, {10, 11, 100, 101, 255}, {1, 2, 3}, {255, 0, 255}, {100, 100, 100}} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 128 {
			data = data[:128]
		}
		h := NewHistogram([]float64{0, 10, 100})
		var counts [4]int64
		var sum float64
		for i, b := range data {
			v := float64(b)
			before := h.Snapshot()
			if h.TryObserve(math.NaN()) || !reflect.DeepEqual(before, h.Snapshot()) {
				t.Fatal("NaN mutated")
			}
			if !h.TryObserve(v) {
				t.Fatal("valid rejected")
			}
			bucket := 3
			switch {
			case b == 0:
				bucket = 0
			case b <= 10:
				bucket = 1
			case b <= 100:
				bucket = 2
			}
			counts[bucket]++
			sum += v
			s := h.Snapshot()
			if s.Count != int64(i+1) || s.Sum != sum || !reflect.DeepEqual(s.Counts, counts[:]) {
				t.Fatal("model mismatch")
			}
		}
		for _, p := range []float64{0, 25, 50, 100} {
			got := h.Percentile(p)
			if len(data) == 0 {
				if got != 0 {
					t.Fatal("empty")
				}
				continue
			}
			rank := int(math.Ceil(p * float64(len(data)) / 100))
			if rank == 0 {
				rank = 1
			}
			accumulated := 0
			want := math.Inf(1)
			for i, c := range counts {
				accumulated += int(c)
				if accumulated >= rank {
					if i < 3 {
						want = []float64{0, 10, 100}[i]
					}
					break
				}
			}
			if got != want {
				t.Fatal("percentile bucket mismatch")
			}
		}
	})
}

type boundaryWriter struct {
	mode   string
	calls  int
	output bytes.Buffer
}

func (w *boundaryWriter) Write(p []byte) (int, error) {
	w.calls++
	switch w.mode {
	case "zero":
		return 0, nil
	case "short":
		n := len(p) / 2
		w.output.Write(p[:n])
		return n, nil
	case "negative":
		return -1, nil
	case "oversized":
		return len(p) + 1, nil
	case "full-error":
		w.output.Write(p)
		return len(p), errors.New("private-writer-detail")
	case "short-error":
		n := len(p) / 2
		w.output.Write(p[:n])
		return n, io.ErrShortWrite
	default:
		return w.output.Write(p)
	}
}

func TestJSONDeliveryBoundaryRejectsIncompleteWrite(t *testing.T) {
	for _, mode := range []string{"zero", "short", "negative", "oversized", "full-error", "short-error"} {
		t.Run(mode, func(t *testing.T) {
			w := &boundaryWriter{mode: mode}
			l, err := NewJSONLogger(LevelInfo, w, Policy{"heartbeat": {}})
			if err != nil {
				t.Fatal(err)
			}
			if err = l.TryEmit(LevelInfo, "heartbeat", nil); err != ErrDeliveryUnknown {
				t.Fatalf("incomplete/invalid write reported success: mode=%s got=%v", mode, err)
			}
			if w.calls != 1 {
				t.Fatalf("ambiguous event was retried: calls=%d", w.calls)
			}
			if mode == "short" || mode == "short-error" {
				var v any
				if json.Unmarshal(w.output.Bytes(), &v) == nil {
					t.Fatal("fixture did not truncate JSON")
				}
			}
			// The caller explicitly discards the failed sink and installs a fresh logger.
			// This is a new heartbeat, not a retry of the uncertain event.
			replacement := &boundaryWriter{}
			fresh, err := NewJSONLogger(LevelInfo, replacement, Policy{"heartbeat": {}})
			if err != nil {
				t.Fatal(err)
			}
			if err = fresh.TryEmit(LevelInfo, "heartbeat", nil); err != nil {
				t.Fatal(err)
			}
			var record map[string]any
			if json.Unmarshal(replacement.output.Bytes(), &record) != nil || record["msg"] != "heartbeat" || replacement.calls != 1 {
				t.Fatal("fresh sink not usable")
			}
		})
	}
}

func TestJSONDeliveryBoundaryFullWrite(t *testing.T) {
	w := &boundaryWriter{}
	l, err := NewJSONLogger(LevelInfo, w, Policy{"heartbeat": {}})
	if err != nil {
		t.Fatal(err)
	}
	if err = l.TryEmit(LevelInfo, "heartbeat", nil); err != nil || w.calls != 1 {
		t.Fatalf("valid delivery failed: %v", err)
	}
	var record map[string]any
	if json.Unmarshal(w.output.Bytes(), &record) != nil || record["msg"] != "heartbeat" {
		t.Fatal("not a complete JSON record")
	}
}

// This is an isolated instrument integration test, not production middleware.
// Numeric slots identify synthetic requests only and never become metric labels.
func TestReferenceHTTPFileTelemetry(t *testing.T) {
	const requests = 32
	p := filepath.Join(t.TempDir(), "events.jsonl")
	file, err := os.OpenFile(p, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	policy := Policy{
		"request.started":   {"slot": {Kind: Int64Field, Min: 1, Max: requests}},
		"request.completed": {"slot": {Kind: Int64Field, Min: 1, Max: requests}, "outcome": {Kind: EnumField, Values: []string{"ok"}}},
	}
	logger, err := NewJSONLogger(LevelInfo, file, policy)
	if err != nil {
		t.Fatal(err)
	}
	var total, failed Counter
	latency, err := NewCheckedHistogram([]float64{0, 1, 10, 1000})
	if err != nil {
		t.Fatal(err)
	}
	var next int64
	var assign sync.Mutex
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		assign.Lock()
		next++
		slot := next
		assign.Unlock()
		if err := logger.TryEmit(LevelInfo, "request.started", map[string]any{"slot": slot}); err != nil {
			failed.Inc(1)
			http.Error(w, "telemetry unavailable", http.StatusServiceUnavailable)
			return
		}
		if err := logger.TryEmit(LevelInfo, "request.completed", map[string]any{"slot": slot, "outcome": "ok"}); err != nil {
			failed.Inc(1)
			http.Error(w, "telemetry unavailable", http.StatusServiceUnavailable)
			return
		}
		total.Inc(1)
		if !latency.TryObserve(float64(time.Since(started)) / float64(time.Millisecond)) {
			failed.Inc(1)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	server := httptest.NewServer(handler)
	defer server.Close()
	client := server.Client()
	client.Timeout = 5 * time.Second
	var wg sync.WaitGroup
	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, err := http.NewRequest(http.MethodGet, server.URL+"/?email=private-marker@example.invalid", nil)
			if err != nil {
				t.Error(err)
				return
			}
			req.Header.Set("Authorization", "Bearer private-marker")
			resp, err := client.Do(req)
			if err != nil {
				t.Error(err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusNoContent {
				t.Errorf("status=%d", resp.StatusCode)
			}
		}()
	}
	wg.Wait()
	server.Close()
	if err := file.Sync(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(data, []byte("private-marker")) || bytes.Contains(data, []byte("Authorization")) {
		t.Fatal("request details reached telemetry")
	}
	records := map[int64][]string{}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	lines := 0
	for scanner.Scan() {
		var record map[string]json.RawMessage
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			t.Fatal(err)
		}
		var slot int64
		var event string
		if json.Unmarshal(record["slot"], &slot) != nil || json.Unmarshal(record["msg"], &event) != nil || slot < 1 || slot > requests {
			t.Fatal("invalid correlation")
		}
		wantKeys := 4
		if event == "request.completed" {
			wantKeys = 5
		}
		if len(record) != wantKeys {
			t.Fatal("unexpected fields")
		}
		records[slot] = append(records[slot], event)
		lines++
	}
	if scanner.Err() != nil || lines != 2*requests || len(records) != requests {
		t.Fatal("missing complete records")
	}
	for slot, events := range records {
		if len(events) != 2 || events[0] != "request.started" || events[1] != "request.completed" {
			t.Fatalf("correlation/order mismatch slot=%d", slot)
		}
	}
	if total.Inc(0) != requests || failed.Inc(0) != 0 {
		t.Fatal("counter mismatch")
	}
	snap := latency.Snapshot()
	if snap.Count != requests {
		t.Fatal("histogram mismatch")
	}
	var sum int64
	for _, count := range snap.Counts {
		sum += count
	}
	if sum != snap.Count {
		t.Fatal("histogram buckets mismatch")
	}
	t.Logf("REFERENCE_HTTP_FILE_PASS requests=%d records=%d correlations=%d counters=%d histogram=%d private_marker_absent=true", requests, lines, len(records), total.Inc(0), snap.Count)
}

func TestReferenceClosedFileAndFreshDestination(t *testing.T) {
	p := filepath.Join(t.TempDir(), "closed.jsonl")
	file, err := os.Create(p)
	if err != nil {
		t.Fatal(err)
	}
	l, err := NewJSONLogger(LevelInfo, file, Policy{"heartbeat": {}})
	if err != nil {
		t.Fatal(err)
	}
	if err = file.Close(); err != nil {
		t.Fatal(err)
	}
	if err = l.TryEmit(LevelInfo, "heartbeat", nil); err != ErrDeliveryUnknown {
		t.Fatalf("closed file not diagnosed: %v", err)
	}
	contents, err := os.ReadFile(p)
	if err != nil || len(contents) != 0 {
		t.Fatal("closed sink mutated")
	}
	fresh, err := os.Create(filepath.Join(t.TempDir(), "fresh.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	defer fresh.Close()
	healthy, err := NewJSONLogger(LevelInfo, fresh, Policy{"heartbeat": {}})
	if err != nil {
		t.Fatal(err)
	}
	if err = healthy.TryEmit(LevelInfo, "heartbeat", nil); err != nil {
		t.Fatal(err)
	}
	if err = fresh.Sync(); err != nil {
		t.Fatal(err)
	}
	if _, err = fresh.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	var record map[string]any
	if err = json.NewDecoder(fresh).Decode(&record); err != nil || record["msg"] != "heartbeat" {
		t.Fatal("new destination unavailable")
	}
	t.Log(fmt.Sprint("REFERENCE_SINK_RECOVERY_PASS closed_destination=delivery_unknown replacement=new_event old_destination_unchanged=true"))
}

// V374: compare accepted attributes with the fixed official Go JSONHandler.
// Equal accepted serialization does not make the policy wrapper an OTel SDK.
func TestSourceSlogAcceptedSerializationV374(t *testing.T) {
	var local, official bytes.Buffer
	logger, e := NewJSONLogger(LevelInfo, &local, fixturePolicy())
	if e != nil {
		t.Fatal(e)
	}
	if e = logger.TryEmit(LevelWarn, "job.done", fixtureFields()); e != nil {
		t.Fatal(e)
	}
	handler := slog.NewJSONHandler(&official, &slog.HandlerOptions{Level: slog.LevelDebug})
	record := slog.NewRecord(time.Now(), slog.LevelWarn, "job.done", 0)
	record.AddAttrs(slog.String("outcome", "ok"), slog.Bool("retry", false), slog.Int64("count", 1))
	if e = handler.Handle(context.Background(), record); e != nil {
		t.Fatal(e)
	}
	var a, b map[string]any
	if e = json.Unmarshal(local.Bytes(), &a); e != nil {
		t.Fatal(e)
	}
	if e = json.Unmarshal(official.Bytes(), &b); e != nil {
		t.Fatal(e)
	}
	delete(a, "time")
	delete(b, "time")
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("local %v official %v", a, b)
	}
}
````


## 6. Configuration surface

Nivel, writer/sink y Policy explícitos; ningún secreto. Campos son obligatorios,
sin extras. Cada string runtime debe pertenecer al enum revisado. Prohibidos
time/level/msg/source como nombres de campo. Constructor copia política y enums.
Caller no muta maps durante uso concurrente; sink no reentra en el logger.
Usar TryEmit y manejar error sin imprimir el input rechazado. Los wrappers legacy
descartan errores y no sirven para acreditar entrega. La política de logs no
reescribe ni destruye originales autorizados de chats/documentos.

## 7. Dependency bill

| Package/image/tool | Pin exacto | Uso | Licencia | Runtime/build | Fuente oficial |
|---|---|---|---|---|---|
| Go standard library | go 1.26.7 (toolchain) | observabilidad | BSD-3-Clause | runtime | https://go.dev |

## 8. Apply order

1. Mantener CANDIDATE: sólo materialización aislada para reparación y pruebas.
2. Colocar ambos archivos en módulo Go de prueba; revisar política y migración.
3. Ejecutar tests/vet/build y GO_NATIVE_FUZZ_GATE con ambos targets del pack.
4. Antes de incorporación, demostrar integración, privacidad del target,
   delivery/retención/performance y admisión. No componer por reputación o PASS local.
5. Rollback de candidato: restaurar snapshot verificado anterior sólo en staging;
   no reactivar logger vulnerable 0.1.x ni borrar observabilidad del producto.

## 9. Verification

Histórico V219/0.1.0: cinco tests, full de 41 paquetes y vet PASS. No se hereda
a 0.1.1 ni se presenta como comprobación de privacidad. V291 verifica el contador
aislado con negativos, overflow, concurrencia, oráculo big.Int y fuzz acotado;
integración completa/SCA/race/carga/operación del target permanecen pendientes.
V292 verifica 0.2.0 con política cerrada y Go JSONHandler; evidencia:
`reconstruction_evidence/PRIVACY_POLICY_LOGGER_REPAIR_V292.md`.

## 10. Reconstruction evidence

Véase `reconstruction_evidence/GO_OBSERVABILITY_CORE_2026-09-02_V219.md`.

## Revisión de compatibilidad V365

Metadata 0.3.2; las dos fuentes materializables conservan sus SHA y APIs.
Declaración anterior: `["GO-ENTERPRISE-BACKEND-CORE 0.1.x"]`. Se conserva como historial, no como
compatibilidad actual demostrada. El módulo importa sólo stdlib; no existe un
adapter a los owners citados en estos dos archivos. No se sustituye el rango por
un latest ni se admite un CANDIDATE. El análisis de versión/identidad y las suites
del conjunto están en reconstruction_evidence/CORE_COMPATIBILITY_AUDIT_V365.md.
La equivalencia funcional con el perfil integral sigue pendiente por los owners
de V293; esta corrección no agrega módulos ni reduce requisitos.

## Auditoría por claim V374

La revisión 0.3.3 se reconstruye y prueba en Go1.26.8/Windows con fuentes
oficiales fijadas y comparaciones por contrato. Ver
`reconstruction_evidence/CORE_CLAIM_ADMISSION_V374.md` desde la raíz: identidad,
licencia local, prueba por invariante, fuente semántica y dictamen de equivalencia
se mantienen separados. El lenguaje/stdlib no es fuente de un negocio empresarial.
No cambia este claim, API, permiso de composición ni las condiciones de integración.
La revisión 0.3.2 y sus bytes quedan preservados en el expediente anterior.
