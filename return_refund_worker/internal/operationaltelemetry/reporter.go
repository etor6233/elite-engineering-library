// Package operationaltelemetry is an AUTHORED, closed-outcome OTLP/HTTP JSON
// reporter for the refund host. It is not a general OpenTelemetry SDK, universal
// PII filter, business ledger, durable acknowledgement or automatic retry queue.
package operationaltelemetry

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"math"
	"mime"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var ErrConfig = errors.New("operational telemetry configuration rejected: TELEMETRY_CONFIG_REJECTED")
var ErrUnavailable = errors.New("operational telemetry delivery unknown: TELEMETRY_DELIVERY_UNKNOWN")
var ErrObservation = errors.New("operational telemetry observation rejected: TELEMETRY_OBSERVATION_REJECTED")

type Config struct {
	Endpoint          string `json:"endpoint"`
	CAFile            string `json:"ca_file"`
	CASHA256          string `json:"ca_sha256"`
	CertificateFile   string `json:"certificate_file"`
	CertificateSHA256 string `json:"certificate_sha256"`
	PrivateKeyFile    string `json:"private_key_file"`
	PrivateKeySHA256  string `json:"private_key_sha256"`
	TimeoutMS         int    `json:"timeout_ms"`
}

var outcomes = [...]string{"started", "idle", "handled", "error", "stopped"}
var bounds = [...]float64{0.001, 0.01, 0.1, 1, 10, 60}

type series struct {
	count   int64
	sum     float64
	buckets [7]int64
}
type Reporter struct {
	client   *http.Client
	base     url.URL
	timeout  time.Duration
	gate     chan struct{}
	closed   atomic.Bool
	instance string
	started  time.Time
	sequence int64
	values   [5]series
}

// strictObject preserves duplicate-key detection instead of relying on the
// last-value-wins behavior of encoding/json. All input is already byte-bounded.
func strictObject(data []byte, allowed map[string]bool) (map[string]json.RawMessage, error) {
	d := json.NewDecoder(bytes.NewReader(data))
	token, err := d.Token()
	if err != nil || token != json.Delim('{') {
		return nil, ErrConfig
	}
	values := map[string]json.RawMessage{}
	for d.More() {
		token, err = d.Token()
		if err != nil {
			return nil, ErrConfig
		}
		key, ok := token.(string)
		if !ok || !allowed[key] {
			return nil, ErrConfig
		}
		if _, exists := values[key]; exists {
			return nil, ErrConfig
		}
		var raw json.RawMessage
		if d.Decode(&raw) != nil {
			return nil, ErrConfig
		}
		values[key] = raw
	}
	token, err = d.Token()
	if err != nil || token != json.Delim('}') {
		return nil, ErrConfig
	}
	if _, err = d.Token(); err != io.EOF {
		return nil, ErrConfig
	}
	return values, nil
}
func fixedFile(path, digest string, limit int64) ([]byte, error) {
	expected, err := hex.DecodeString(digest)
	if err != nil || len(expected) != 32 || len(digest) != 64 {
		return nil, ErrConfig
	}
	abs, err := filepath.Abs(path)
	if err != nil || path == "" {
		return nil, ErrConfig
	}
	for p := abs; ; p = filepath.Dir(p) {
		st, e := os.Lstat(p)
		if e != nil || st.Mode()&os.ModeSymlink != 0 {
			return nil, ErrConfig
		}
		parent := filepath.Dir(p)
		if parent == p {
			break
		}
	}
	f, err := os.Open(abs)
	if err != nil {
		return nil, ErrConfig
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil || !st.Mode().IsRegular() || st.Size() > limit {
		return nil, ErrConfig
	}
	data, err := io.ReadAll(io.LimitReader(f, limit+1))
	if err != nil || int64(len(data)) > limit {
		return nil, ErrConfig
	}
	got := sha256.Sum256(data)
	if !bytes.Equal(got[:], expected) {
		return nil, ErrConfig
	}
	return data, nil
}

// FromFile binds the config and every credential file to caller-approved hashes.
// Private key bytes are parsed in memory and never added to diagnostics/telemetry.
func FromFile(path, digest string) (*Reporter, error) {
	data, err := fixedFile(path, digest, 8192)
	if err != nil {
		return nil, ErrConfig
	}
	allowed := map[string]bool{"endpoint": true, "ca_file": true, "ca_sha256": true, "certificate_file": true, "certificate_sha256": true, "private_key_file": true, "private_key_sha256": true, "timeout_ms": true}
	values, err := strictObject(data, allowed)
	if err != nil || len(values) != len(allowed) {
		return nil, ErrConfig
	}
	var cfg Config
	if json.Unmarshal(data, &cfg) != nil {
		return nil, ErrConfig
	}
	ca, e1 := fixedFile(cfg.CAFile, cfg.CASHA256, 1<<20)
	cert, e2 := fixedFile(cfg.CertificateFile, cfg.CertificateSHA256, 1<<20)
	key, e3 := fixedFile(cfg.PrivateKeyFile, cfg.PrivateKeySHA256, 1<<20)
	if e1 != nil || e2 != nil || e3 != nil {
		return nil, ErrConfig
	}
	if cfg.TimeoutMS < 10 || cfg.TimeoutMS > 3000 {
		return nil, ErrConfig
	}
	return New(cfg.Endpoint, time.Duration(cfg.TimeoutMS)*time.Millisecond, ca, cert, key)
}

func New(endpoint string, timeout time.Duration, caPEM, certPEM, keyPEM []byte) (*Reporter, error) {
	u, err := url.Parse(endpoint)
	if err != nil || len(endpoint) > 2048 || u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || u.Opaque != "" || u.RawPath != "" || (u.Path != "" && u.Path != "/") || timeout < 10*time.Millisecond || timeout > 3*time.Second {
		return nil, ErrConfig
	}
	ca := x509.NewCertPool()
	if !ca.AppendCertsFromPEM(caPEM) {
		return nil, ErrConfig
	}
	certificate, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, ErrConfig
	}
	instance := make([]byte, 16)
	if _, err = rand.Read(instance); err != nil {
		return nil, ErrConfig
	}
	transport := &http.Transport{Proxy: nil, DialContext: (&net.Dialer{Timeout: timeout}).DialContext, TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS13, RootCAs: ca, Certificates: []tls.Certificate{certificate}}, TLSHandshakeTimeout: timeout, ResponseHeaderTimeout: timeout, MaxResponseHeaderBytes: 8192, DisableKeepAlives: true, ForceAttemptHTTP2: false}
	return &Reporter{client: &http.Client{Transport: transport, Timeout: timeout, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}, base: *u, timeout: timeout, gate: make(chan struct{}, 1), instance: hex.EncodeToString(instance), started: time.Now()}, nil
}
func (r *Reporter) Close() {
	if r != nil {
		r.closed.Store(true)
		r.client.CloseIdleConnections()
	}
}
func textAttr(key, value string) any {
	return map[string]any{"key": key, "value": map[string]any{"stringValue": value}}
}
func numberAttr(key string, value int64) any {
	return map[string]any{"key": key, "value": map[string]any{"intValue": strconv.FormatInt(value, 10)}}
}
func stamp(t time.Time) string { return strconv.FormatInt(t.UnixNano(), 10) }
func (r *Reporter) resource() any {
	return map[string]any{"attributes": []any{textAttr("service.name", "elite-return-refund-worker"), textAttr("deployment.environment.name", "reference"), textAttr("service.instance.id", r.instance)}}
}
func scope() any { return map[string]any{"name": "elite.refund.host", "version": "1.0.0"} }

// Observe serializes calls within one shared deadline (including waiting for
// another observer). Input is a closed outcome and monotonic elapsed duration;
// it accepts no message, error, request, tenant, business ID or arbitrary fields.
// A failed or partial three-signal export returns unknown and is never retried.
func (r *Reporter) Observe(parent context.Context, outcome string, elapsed time.Duration) error {
	idx := -1
	for i, value := range outcomes {
		if value == outcome {
			idx = i
		}
	}
	if r == nil || parent == nil || idx < 0 || elapsed < 0 {
		return ErrObservation
	}
	ctx, cancel := context.WithTimeout(parent, r.timeout)
	defer cancel()
	select {
	case r.gate <- struct{}{}:
		defer func() { <-r.gate }()
	case <-ctx.Done():
		return ErrUnavailable
	}
	if r.closed.Load() || ctx.Err() != nil {
		return ErrUnavailable
	}
	if r.sequence == math.MaxInt64 || r.values[idx].count == math.MaxInt64 {
		return ErrObservation
	}
	seconds := elapsed.Seconds()
	sum := r.values[idx].sum + seconds
	if math.IsInf(sum, 0) || math.IsNaN(sum) {
		return ErrObservation
	}
	ids := make([]byte, 24)
	if _, err := rand.Read(ids); err != nil {
		return ErrUnavailable
	}
	trace, span := hex.EncodeToString(ids[:16]), hex.EncodeToString(ids[16:])
	end := time.Now()
	start := end.Add(-elapsed)
	r.sequence++
	r.values[idx].count++
	r.values[idx].sum = sum
	bucket := len(bounds)
	for i, b := range bounds {
		if seconds <= b {
			bucket = i
			break
		}
	}
	r.values[idx].buckets[bucket]++
	attrs := []any{textAttr("outcome", outcome), numberAttr("sequence", r.sequence)}
	severity := 9
	status := map[string]any{"code": 1}
	if outcome == "error" {
		severity = 17
		status = map[string]any{"code": 2}
	}
	log := map[string]any{"resourceLogs": []any{map[string]any{"resource": r.resource(), "scopeLogs": []any{map[string]any{"scope": scope(), "logRecords": []any{map[string]any{"timeUnixNano": stamp(end), "observedTimeUnixNano": stamp(end), "severityNumber": severity, "body": map[string]any{"stringValue": "refund_host_observation"}, "attributes": attrs, "traceId": trace, "spanId": span}}}}}}}
	spans := map[string]any{"resourceSpans": []any{map[string]any{"resource": r.resource(), "scopeSpans": []any{map[string]any{"scope": scope(), "spans": []any{map[string]any{"traceId": trace, "spanId": span, "name": "refund.host.observation", "kind": 1, "startTimeUnixNano": stamp(start), "endTimeUnixNano": stamp(end), "attributes": attrs, "status": status}}}}}}}
	counts := []any{}
	hist := []any{}
	for i, value := range r.values {
		a := []any{textAttr("outcome", outcomes[i])}
		counts = append(counts, map[string]any{"attributes": a, "startTimeUnixNano": stamp(r.started), "timeUnixNano": stamp(end), "asInt": strconv.FormatInt(value.count, 10)})
		buckets := []string{}
		for _, n := range value.buckets {
			buckets = append(buckets, strconv.FormatInt(n, 10))
		}
		hist = append(hist, map[string]any{"attributes": a, "startTimeUnixNano": stamp(r.started), "timeUnixNano": stamp(end), "count": strconv.FormatInt(value.count, 10), "sum": value.sum, "explicitBounds": bounds[:], "bucketCounts": buckets})
	}
	metrics := map[string]any{"resourceMetrics": []any{map[string]any{"resource": r.resource(), "scopeMetrics": []any{map[string]any{"scope": scope(), "metrics": []any{map[string]any{"name": "elite_refund_observations", "sum": map[string]any{"dataPoints": counts, "aggregationTemporality": 2, "isMonotonic": true}}, map[string]any{"name": "elite_refund_observation_duration", "unit": "s", "histogram": map[string]any{"dataPoints": hist, "aggregationTemporality": 2}}}}}}}}
	for _, signal := range []struct {
		path string
		data any
	}{{"logs", log}, {"metrics", metrics}, {"traces", spans}} {
		if r.send(ctx, signal.path, signal.data) != nil {
			return ErrUnavailable
		}
	}
	return nil
}
func (r *Reporter) send(ctx context.Context, signal string, value any) error {
	data, err := json.Marshal(value)
	if err != nil || len(data) > 16384 {
		return ErrUnavailable
	}
	u := r.base
	u.Path = "/v1/" + signal
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(data))
	if err != nil {
		return ErrUnavailable
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := r.client.Do(req)
	if err != nil {
		return ErrUnavailable
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return ErrUnavailable
	}
	typ, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if err != nil || typ != "application/json" {
		return ErrUnavailable
	}
	b, err := io.ReadAll(io.LimitReader(response.Body, 1025))
	if err != nil || len(b) > 1024 {
		return ErrUnavailable
	}
	fields, err := strictObject(b, map[string]bool{"partialSuccess": true})
	if err != nil {
		return ErrUnavailable
	}
	if p, ok := fields["partialSuccess"]; ok {
		name := map[string]string{"logs": "rejectedLogRecords", "metrics": "rejectedDataPoints", "traces": "rejectedSpans"}[signal]
		partial, e := strictObject(p, map[string]bool{name: true, "errorMessage": true})
		if e != nil {
			return ErrUnavailable
		}
		for key, raw := range partial {
			if key == "errorMessage" {
				var message string
				if json.Unmarshal(raw, &message) != nil || message != "" {
					return ErrUnavailable
				}
			} else if strings.TrimSpace(string(raw)) != "0" && strings.TrimSpace(string(raw)) != "\"0\"" {
				return ErrUnavailable
			}
		}
	}
	return nil
}
