package main

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"elite.local/enterprise/internal/whatsappbridge"
)

func whatsappTLSFixture(t *testing.T, linesPerConnection ...int) (whatsappHostConfig, map[string]string, <-chan []byte) {
	t.Helper()
	key, e := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if e != nil {
		t.Fatal(e)
	}
	now := time.Now()
	template := &x509.Certificate{SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "local-fixture"}, DNSNames: []string{"collector.local"}, NotBefore: now.Add(-time.Hour), NotAfter: now.Add(time.Hour), IsCA: true, BasicConstraintsValid: true, KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign, ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}}
	der, e := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if e != nil {
		t.Fatal(e)
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyDER, e := x509.MarshalPKCS8PrivateKey(key)
	if e != nil {
		t.Fatal(e)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER})
	pair, e := tls.X509KeyPair(certPEM, keyPEM)
	if e != nil {
		t.Fatal(e)
	}
	roots := x509.NewCertPool()
	roots.AppendCertsFromPEM(certPEM)
	listener, e := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{MinVersion: tls.VersionTLS13, Certificates: []tls.Certificate{pair}, ClientAuth: tls.RequireAndVerifyClientCert, ClientCAs: roots})
	if e != nil {
		t.Fatal(e)
	}
	stop := make(chan struct{})
	output := make(chan []byte, 4)
	go func() {
		defer close(stop)
		for {
			conn, e := listener.Accept()
			if e != nil {
				return
			}
			func(conn net.Conn) {
				defer conn.Close()
				_ = conn.SetDeadline(time.Now().Add(3 * time.Second))
				lines := 1
				if len(linesPerConnection) == 1 {
					lines = linesPerConnection[0]
				}
				reader := bufio.NewReader(conn)
				for i := 0; i < lines; i++ {
					line, e := reader.ReadBytes('\n')
					if e != nil {
						return
					}
					output <- line
				}
			}(conn)
		}
	}()
	t.Cleanup(func() { _ = listener.Close(); <-stop })
	dir := t.TempDir()
	ca := filepath.Join(dir, "ca.pem")
	cert := filepath.Join(dir, "client.pem")
	private := filepath.Join(dir, "client-key.pem")
	for p, b := range map[string][]byte{ca: certPEM, cert: certPEM, private: keyPEM} {
		if e := os.WriteFile(p, b, 0600); e != nil {
			t.Fatal(e)
		}
	}
	var c whatsappHostConfig
	c.Collector.Address = listener.Addr().String()
	c.Collector.ServerName = "collector.local"
	c.Collector.CAFile = ca
	h := sha256.Sum256(certPEM)
	c.Collector.CASHA256 = hex.EncodeToString(h[:])
	return c, map[string]string{"WHATSAPP_COLLECTOR_CLIENT_CERT_FILE": cert, "WHATSAPP_COLLECTOR_CLIENT_KEY_FILE": private}, output
}

func TestWhatsAppHostAuthenticatedReporter(t *testing.T) {
	c, values, out := whatsappTLSFixture(t)
	lookup := func(k string) string { return values[k] }
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	r, e := dialWhatsAppReporter(ctx, c, lookup)
	if e != nil {
		t.Fatal(e)
	}
	defer r.Close()
	if e = r.Report(ctx, whatsappbridge.StatusPollReport{Outcome: "COMPLETED", Claimed: true, InsertedObservations: 1, NextDelay: time.Second}); e != nil {
		t.Fatal(e)
	}
	select {
	case line := <-out:
		var v map[string]any
		if json.Unmarshal(line, &v) != nil || v["outcome"] != "COMPLETED" || len(v) != 6 {
			t.Fatal(string(line))
		}
	case <-ctx.Done():
		t.Fatal("no authenticated report")
	}
	c.Collector.ServerName = "wrong.local"
	if r, e := dialWhatsAppReporter(ctx, c, lookup); e == nil {
		r.Close()
		t.Fatal("untrusted collector identity")
	}
	c.Collector.ServerName = "collector.local"
	values["WHATSAPP_COLLECTOR_CLIENT_KEY_FILE"] = ""
	if r, e := dialWhatsAppReporter(ctx, c, lookup); e == nil {
		r.Close()
		t.Fatal("missing client key")
	}
}
