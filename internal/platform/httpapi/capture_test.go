package httpapi

import (
	"bufio"
	"context"
	cap "elite.local/enterprise/internal/capture"
	"elite.local/enterprise/internal/platform/identity"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type noReadBody struct{ t *testing.T }

func (r noReadBody) Read([]byte) (int, error) {
	r.t.Fatal("unauthorized capture read body")
	return 0, nil
}

func TestCaptureSlowFrameBounded(t *testing.T) {
	spy := &captureSpy{}
	mux := http.NewServeMux()
	CaptureModule{Decoder: spy}.Register(mux, captureVerifier{identity.Principal{Subject: "operator", TenantID: "018f4d4a-7b36-7a21-8d10-2f4c54c28901", Permissions: map[string]struct{}{"inventory:read": {}}, Organizations: map[string]struct{}{"org": {}}}})
	server := httptest.NewServer(mux)
	defer server.Close()
	conn, e := net.Dial("tcp", strings.TrimPrefix(server.URL, "http://"))
	if e != nil {
		t.Fatal(e)
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(7 * time.Second))
	_, e = fmt.Fprintf(conn, "POST /v1/capture/decode HTTP/1.1\r\nHost: fixture\r\nAuthorization: Bearer fixture\r\nContent-Type: application/octet-stream\r\nContent-Length: 1000\r\n\r\nE")
	if e != nil {
		t.Fatal(e)
	}
	r, e := http.ReadResponse(bufio.NewReader(conn), nil)
	if e != nil {
		t.Fatal(e)
	}
	defer r.Body.Close()
	if r.StatusCode != 408 || spy.calls != 0 {
		t.Fatal(r.StatusCode, spy.calls)
	}
}
func (noReadBody) Close() error { return nil }

type captureVerifier struct{ p identity.Principal }

func (v captureVerifier) Verify(context.Context, string) (identity.Principal, error) { return v.p, nil }

type captureSpy struct{ calls int }

func (s *captureSpy) Decode(context.Context, []byte) cap.Result { s.calls++; return cap.Result{} }
func TestCaptureAuthenticatesBeforeBodyAndDecoder(t *testing.T) {
	spy := &captureSpy{}
	mux := http.NewServeMux()
	CaptureModule{Decoder: spy}.Register(mux, captureVerifier{identity.Principal{Subject: "no-privilege", TenantID: "018f4d4a-7b36-7a21-8d10-2f4c54c28901", Organizations: map[string]struct{}{"org": {}}}})
	r := httptest.NewRequest("POST", "/v1/capture/decode", nil)
	r.Body = noReadBody{t}
	r.Header.Set("Authorization", "Bearer fixture")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, r)
	if rec.Code != 403 || spy.calls != 0 {
		t.Fatal(rec.Code, spy.calls)
	}
}
