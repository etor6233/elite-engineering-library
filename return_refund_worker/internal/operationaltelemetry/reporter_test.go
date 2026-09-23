package operationaltelemetry

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"math"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type certs struct {
	ca, client, key []byte
	server          tls.Certificate
	pool            *x509.CertPool
}

func certificates(t *testing.T) certs {
	t.Helper()
	now := time.Now()
	pub, priv, e := ed25519.GenerateKey(rand.Reader)
	if e != nil {
		t.Fatal(e)
	}
	root := &x509.Certificate{SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "reference-test-ca"}, NotBefore: now.Add(-time.Minute), NotAfter: now.Add(time.Hour), IsCA: true, BasicConstraintsValid: true, KeyUsage: x509.KeyUsageCertSign}
	raw, e := x509.CreateCertificate(rand.Reader, root, root, pub, priv)
	if e != nil {
		t.Fatal(e)
	}
	ca := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: raw})
	makeLeaf := func(serial int64, server bool) ([]byte, []byte) {
		p, k, e := ed25519.GenerateKey(rand.Reader)
		if e != nil {
			t.Fatal(e)
		}
		usage := x509.ExtKeyUsageClientAuth
		if server {
			usage = x509.ExtKeyUsageServerAuth
		}
		c := &x509.Certificate{SerialNumber: big.NewInt(serial), Subject: pkix.Name{CommonName: "reference-test-peer"}, NotBefore: now.Add(-time.Minute), NotAfter: now.Add(time.Hour), KeyUsage: x509.KeyUsageDigitalSignature, ExtKeyUsage: []x509.ExtKeyUsage{usage}, IPAddresses: []net.IP{net.ParseIP("127.0.0.1")}}
		b, e := x509.CreateCertificate(rand.Reader, c, root, p, priv)
		if e != nil {
			t.Fatal(e)
		}
		pk, e := x509.MarshalPKCS8PrivateKey(k)
		if e != nil {
			t.Fatal(e)
		}
		return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: b}), pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: pk})
	}
	client, key := makeLeaf(2, false)
	sc, sk := makeLeaf(3, true)
	pair, e := tls.X509KeyPair(sc, sk)
	if e != nil {
		t.Fatal(e)
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(ca)
	return certs{ca, client, key, pair, pool}
}
func serverFor(t *testing.T, c certs, h http.HandlerFunc) *httptest.Server {
	t.Helper()
	s := httptest.NewUnstartedServer(h)
	s.TLS = &tls.Config{MinVersion: tls.VersionTLS13, Certificates: []tls.Certificate{c.server}, ClientAuth: tls.RequireAndVerifyClientCert, ClientCAs: c.pool}
	s.StartTLS()
	t.Cleanup(s.Close)
	return s
}
func newReporter(t *testing.T, c certs, s *httptest.Server) *Reporter {
	t.Helper()
	r, e := New(s.URL, 3*time.Second, c.ca, c.client, c.key)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(r.Close)
	return r
}
func accepted(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = io.WriteString(w, "{}")
}
func decodeMap(t *testing.T, b []byte) map[string]any {
	t.Helper()
	var j map[string]any
	if json.Unmarshal(b, &j) != nil {
		t.Fatalf("invalid JSON: %q", b)
	}
	return j
}
func child(v any, key string) any {
	if m, ok := v.(map[string]any); ok {
		return m[key]
	}
	return nil
}
func first(v any) any {
	if a, ok := v.([]any); ok && len(a) > 0 {
		return a[0]
	}
	return nil
}
func TestThreeSignalsPrivacyCorrelationAndMonotonicity(t *testing.T) {
	c := certificates(t)
	var mu sync.Mutex
	var paths []string
	var payloads [][]byte
	s := serverFor(t, c, func(w http.ResponseWriter, q *http.Request) {
		if q.Method != "POST" || q.Header.Get("Content-Type") != "application/json" || len(q.TLS.VerifiedChains) != 1 {
			t.Error("transport contract")
		}
		b, e := io.ReadAll(io.LimitReader(q.Body, 16385))
		if e != nil || len(b) > 16384 {
			t.Error("body budget")
		}
		mu.Lock()
		paths = append(paths, q.URL.Path)
		payloads = append(payloads, b)
		mu.Unlock()
		accepted(w)
	})
	r := newReporter(t, c, s)
	for _, o := range []string{"started", "idle", "handled", "error", "stopped"} {
		if e := r.Observe(context.Background(), o, 10*time.Millisecond); e != nil {
			t.Fatal(e)
		}
	}
	if len(paths) != 15 {
		t.Fatal(paths)
	}
	for n := 0; n < 5; n++ {
		if strings.Join(paths[n*3:n*3+3], ",") != "/v1/logs,/v1/metrics,/v1/traces" {
			t.Fatal(paths)
		}
		l := first(child(first(child(first(child(decodeMap(t, payloads[n*3]), "resourceLogs")), "scopeLogs")), "logRecords"))
		tr := first(child(first(child(first(child(decodeMap(t, payloads[n*3+2]), "resourceSpans")), "scopeSpans")), "spans"))
		if child(l, "traceId") != child(tr, "traceId") || child(l, "spanId") != child(tr, "spanId") || len(child(l, "traceId").(string)) != 32 {
			t.Fatal("correlation mismatch")
		}
		if child(child(l, "body"), "stringValue") != "refund_host_observation" {
			t.Fatal("unapproved message")
		}
	}
	if r.sequence != 5 {
		t.Fatal(r.sequence)
	}
	for _, v := range r.values {
		if v.count != 1 || v.buckets[1] != 1 {
			t.Fatal(v)
		}
	}
	before := len(paths)
	for _, o := range []string{"private@example.invalid", "password=private", "STARTED", ""} {
		if e := r.Observe(context.Background(), o, 0); !errors.Is(e, ErrObservation) {
			t.Fatal(e)
		}
	}
	if e := r.Observe(context.Background(), "idle", -1); !errors.Is(e, ErrObservation) {
		t.Fatal(e)
	}
	if len(paths) != before {
		t.Fatal("rejected input touched sink")
	}
	for _, b := range payloads {
		if strings.Contains(string(b), "private") || strings.Contains(string(b), "password") {
			t.Fatal("private marker")
		}
	}
}
func TestRemoteRejectionNeverRetriesOrLeaksBody(t *testing.T) {
	for _, scenario := range []string{"redirect", "failure", "rate", "partial", "warning", "oversized", "invalid", "duplicate", "wrong-type"} {
		t.Run(scenario, func(t *testing.T) {
			c := certificates(t)
			var calls atomic.Int64
			s := serverFor(t, c, func(w http.ResponseWriter, q *http.Request) {
				calls.Add(1)
				w.Header().Set("Content-Type", "application/json")
				switch scenario {
				case "redirect":
					w.Header().Set("Location", "https://invalid.example.invalid/private")
					w.WriteHeader(302)
				case "failure":
					w.WriteHeader(500)
				case "rate":
					w.WriteHeader(429)
				case "partial":
					_, _ = io.WriteString(w, `{"partialSuccess":{"rejectedLogRecords":"1"}}`)
				case "warning":
					_, _ = io.WriteString(w, `{"partialSuccess":{"errorMessage":"private@example.invalid"}}`)
				case "oversized":
					_, _ = io.WriteString(w, strings.Repeat("x", 1025))
				case "invalid":
					_, _ = io.WriteString(w, "private@example.invalid")
				case "duplicate":
					_, _ = io.WriteString(w, `{"partialSuccess":{},"partialSuccess":{}}`)
				case "wrong-type":
					w.Header().Set("Content-Type", "text/plain")
					_, _ = io.WriteString(w, "{}")
				}
			})
			r := newReporter(t, c, s)
			e := r.Observe(context.Background(), "error", 0)
			if !errors.Is(e, ErrUnavailable) || strings.Contains(e.Error(), "private") {
				t.Fatal(e)
			}
			if calls.Load() != 1 {
				t.Fatal("automatic retry/extra signal", calls.Load())
			}
		})
	}
}
func TestMidExportFailureIsUnknownWithoutReplay(t *testing.T) {
	c := certificates(t)
	var n atomic.Int64
	s := serverFor(t, c, func(w http.ResponseWriter, q *http.Request) {
		if n.Add(1) == 2 {
			w.WriteHeader(503)
			return
		}
		accepted(w)
	})
	r := newReporter(t, c, s)
	if !errors.Is(r.Observe(context.Background(), "handled", 0), ErrUnavailable) || n.Load() != 2 || r.values[2].count != 1 {
		t.Fatal("partial export not retained as uncertain")
	}
}
func TestCancellationBoundsSinkAndQueuedObserver(t *testing.T) {
	c := certificates(t)
	entered := make(chan struct{})
	release := make(chan struct{})
	var once sync.Once
	s := serverFor(t, c, func(w http.ResponseWriter, q *http.Request) {
		once.Do(func() { close(entered) })
		select {
		case <-release:
		case <-q.Context().Done():
		}
	})
	defer close(release)
	r := newReporter(t, c, s)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- r.Observe(ctx, "idle", 0) }()
	select {
	case <-entered:
	case <-time.After(time.Second):
		cancel()
		t.Fatal("handler did not receive request")
	}
	qctx, qcancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer qcancel()
	start := time.Now()
	if !errors.Is(r.Observe(qctx, "idle", 0), ErrUnavailable) || time.Since(start) > time.Second {
		t.Fatal("waiter deadline")
	}
	cancel()
	select {
	case e := <-done:
		if !errors.Is(e, ErrUnavailable) {
			t.Fatal(e)
		}
	case <-time.After(time.Second):
		t.Fatal("active cancellation")
	}
}
func TestParallelObserversAndOverflow(t *testing.T) {
	c := certificates(t)
	var n atomic.Int64
	s := serverFor(t, c, func(w http.ResponseWriter, q *http.Request) { n.Add(1); accepted(w) })
	r := newReporter(t, c, s)
	var wg sync.WaitGroup
	for i := 0; i < 24; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if e := r.Observe(context.Background(), "handled", time.Millisecond); e != nil {
				t.Error(e)
			}
		}()
	}
	wg.Wait()
	if r.sequence != 24 || r.values[2].count != 24 || n.Load() != 72 {
		t.Fatal(r.sequence, n.Load())
	}
	r.sequence = math.MaxInt64
	if !errors.Is(r.Observe(context.Background(), "idle", 0), ErrObservation) || n.Load() != 72 {
		t.Fatal("sequence overflow")
	}
	r.sequence = 24
	r.values[1].count = math.MaxInt64
	if !errors.Is(r.Observe(context.Background(), "idle", 0), ErrObservation) || n.Load() != 72 {
		t.Fatal("counter overflow")
	}
	r.Close()
	if !errors.Is(r.Observe(context.Background(), "idle", 0), ErrUnavailable) {
		t.Fatal("closed observer")
	}
}
func TestTLSIdentityAndEndpointValidation(t *testing.T) {
	c := certificates(t)
	other := certificates(t)
	s := serverFor(t, c, func(w http.ResponseWriter, q *http.Request) { accepted(w) })
	for _, b := range []struct{ ca, cert, key []byte }{{other.ca, c.client, c.key}, {c.ca, other.client, other.key}} {
		r, e := New(s.URL, time.Second, b.ca, b.cert, b.key)
		if e != nil {
			t.Fatal(e)
		}
		if !errors.Is(r.Observe(context.Background(), "started", 0), ErrUnavailable) {
			t.Fatal("foreign identity accepted")
		}
		r.Close()
	}
	for _, u := range []string{"http://localhost", "https://user:password@localhost", "https://localhost/path", "https://localhost?x=y", "https://localhost#fragment", "https://"} {
		if _, e := New(u, time.Second, c.ca, c.client, c.key); !errors.Is(e, ErrConfig) {
			t.Fatal(u, e)
		}
	}
	for _, d := range []time.Duration{0, 9 * time.Millisecond, 4 * time.Second, time.Duration(math.MaxInt64)} {
		if _, e := New(s.URL, d, c.ca, c.client, c.key); !errors.Is(e, ErrConfig) {
			t.Fatal(d, e)
		}
	}
}
func TestExactConfigCredentialHashesAndNoDiagnostics(t *testing.T) {
	c := certificates(t)
	dir := t.TempDir()
	put := func(name string, b []byte) (string, string) {
		p := filepath.Join(dir, name)
		if e := os.WriteFile(p, b, 0600); e != nil {
			t.Fatal(e)
		}
		sum := sha256.Sum256(b)
		return p, hex.EncodeToString(sum[:])
	}
	caf, cah := put("ca.pem", c.ca)
	cf, ch := put("client.pem", c.client)
	kf, kh := put("client.key", c.key)
	cfg := Config{Endpoint: "https://127.0.0.1:1", CAFile: caf, CASHA256: cah, CertificateFile: cf, CertificateSHA256: ch, PrivateKeyFile: kf, PrivateKeySHA256: kh, TimeoutMS: 1000}
	b, _ := json.Marshal(cfg)
	p, h := put("config.json", b)
	r, e := FromFile(p, h)
	if e != nil {
		t.Fatal(e)
	}
	r.Close()
	for i, raw := range [][]byte{append(b[:len(b)-1], []byte(`,"endpoint":"https://other.invalid"}`)...), []byte(`{"endpoint":null}`), []byte(`{"Endpoint":"https://other.invalid"}`), append(b, []byte(` {}`)...)} {
		p, h = put("bad"+string(rune('a'+i)), raw)
		if _, e = FromFile(p, h); !errors.Is(e, ErrConfig) {
			t.Fatal(e)
		}
	}
	p, h = put("config2.json", b)
	if e = os.WriteFile(kf, []byte("private@example.invalid"), 0600); e != nil {
		t.Fatal(e)
	}
	_, e = FromFile(p, h)
	if !errors.Is(e, ErrConfig) || strings.Contains(e.Error(), "private") {
		t.Fatal(e)
	}
	if _, e = FromFile(dir, h); !errors.Is(e, ErrConfig) {
		t.Fatal("directory accepted")
	}
}

// Protocol seeds include exact duplicate keys, trailing values and wrong types.
func FuzzStrictOTLPResponseObject(f *testing.F) {
	for _, seed := range []string{`{}`, `{"partialSuccess":{}}`, `{"partialSuccess":{"rejectedLogRecords":"1"}}`, `{"partialSuccess":{},"partialSuccess":{}}`, `[]`, `{"unexpected":1}`, `{} {}`, `{"partialSuccess":null}`} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, raw string) {
		if len(raw) > 2048 {
			return
		}
		values, err := strictObject([]byte(raw), map[string]bool{"partialSuccess": true})
		if err != nil {
			if err.Error() != ErrConfig.Error() {
				t.Fatal("unbounded diagnostic")
			}
			return
		}
		if !json.Valid([]byte(raw)) {
			t.Fatal("invalid JSON accepted")
		}
		if len(values) > 1 {
			t.Fatal("unknown or duplicate field accepted")
		}
		for k := range values {
			if k != "partialSuccess" {
				t.Fatal("unknown key")
			}
		}
		if strings.HasPrefix(strings.TrimSpace(raw), "[") {
			t.Fatal("array accepted")
		}
	})
}
func TestConfigTimeoutCannotOverflowIntoAcceptedDuration(t *testing.T) {
	c := certificates(t)
	dir := t.TempDir()
	put := func(n string, b []byte) (string, string) {
		p := filepath.Join(dir, n)
		if err := os.WriteFile(p, b, 0600); err != nil {
			t.Fatal(err)
		}
		h := sha256.Sum256(b)
		return p, hex.EncodeToString(h[:])
	}
	caf, cah := put("ca", c.ca)
	cf, ch := put("cert", c.client)
	kf, kh := put("key", c.key)
	for _, n := range []int{0, 9, 3001, math.MaxInt64, 9223372036855} {
		cfg := Config{Endpoint: "https://127.0.0.1:1", CAFile: caf, CASHA256: cah, CertificateFile: cf, CertificateSHA256: ch, PrivateKeyFile: kf, PrivateKeySHA256: kh, TimeoutMS: n}
		b, _ := json.Marshal(cfg)
		p, h := put("cfg", b)
		if _, err := FromFile(p, h); !errors.Is(err, ErrConfig) {
			t.Fatal("timeout accepted", n)
		}
	}
}
