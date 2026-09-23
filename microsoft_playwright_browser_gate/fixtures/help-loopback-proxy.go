// AUTHORED test-only TLS bridge for production Next/JWE help tests. Never deploy.
package main

import (
	"fmt"
	"io"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 { panic("expected explicit loopback Next URL") }
	u, err := url.Parse(os.Args[1])
	if err != nil || u.Scheme != "http" || u.Hostname() != "127.0.0.1" || u.Port() == "" || u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" { panic("only loopback HTTP target allowed") }
	server := httptest.NewTLSServer(httputil.NewSingleHostReverseProxy(u))
	defer server.Close()
	fmt.Println(server.URL)
	done := make(chan struct{})
	go func() { _, _ = io.Copy(io.Discard, os.Stdin); close(done) }()
	select { case <-done: case <-time.After(10 * time.Minute): }
}
