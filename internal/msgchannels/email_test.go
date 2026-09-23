package msgchannels

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestEmailConfigValidate(t *testing.T) {
	cfg := EmailConfig{APIURL: "https://api.sendgrid.com/v3/mail/send", APIKey: "k", From: "ventas@empresa.com", Timeout: time.Second}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("valid config rejected: %v", err)
	}
	bad := cfg
	bad.From = "no-arroba"
	if err := bad.Validate(); err == nil {
		t.Fatal("bad from accepted")
	}
	noKey := cfg
	noKey.APIKey = ""
	if err := noKey.Validate(); err == nil {
		t.Fatal("missing key accepted")
	}
}

func TestSendEmailRoundTrip(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("authorization") != "Bearer k" {
			t.Errorf("missing bearer")
		}
		var p emailPayload
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			t.Fatal(err)
		}
		if p.From.Email != "ventas@empresa.com" || p.Subject != "Oferta" {
			t.Errorf("unexpected payload: %+v", p)
		}
		if len(p.Personalizations) != 1 || p.Personalizations[0].To[0].Email != "cliente@x.com" {
			t.Errorf("unexpected recipient: %+v", p.Personalizations)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	cfg := EmailConfig{APIURL: srv.URL, APIKey: "k", From: "ventas@empresa.com", Timeout: 5 * time.Second}
	if err := SendEmail(context.Background(), cfg, "cliente@x.com", "Oferta", "Hola"); err != nil {
		t.Fatal(err)
	}
}
